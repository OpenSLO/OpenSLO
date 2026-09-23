import argparse
import json
from pathlib import Path
import posixpath
import re
from urllib.parse import urlsplit, urlunsplit

from jinja2 import Environment, FileSystemLoader, StrictUndefined
from markdown.extensions import Extension
from markdown.treeprocessors import Treeprocessor
from markupsafe import Markup, escape
from pydantic import BaseModel, Field, TypeAdapter


ROOT = Path(__file__).resolve().parent


class Rule(BaseModel):
    description: str
    errorCode: str = ""
    details: str = ""
    examples: list[str] = Field(default_factory=list)
    conditions: list[str] = Field(default_factory=list)


class TypeInfo(BaseModel):
    name: str
    kind: str
    package: str = ""


class PropertyPlan(BaseModel):
    path: str
    typeInfo: TypeInfo
    rules: list[Rule] = Field(default_factory=list)
    examples: list[str] = Field(default_factory=list)
    values: list[str] = Field(default_factory=list)


class Property(PropertyPlan):
    typeDoc: str = ""
    fieldDoc: str = ""
    deprecatedDoc: str = ""
    componentPlans: list[PropertyPlan] = Field(default_factory=list)

    @property
    def display_path(self):
        return self.path.removeprefix("$.")

    @property
    def anchor(self):
        # Keep map keys, map values, and array elements distinct in fragment links.
        path = self.display_path.replace(".*~", "-keys").replace(".*", "-values")
        path = path.replace("[*]", "-items")
        return path.replace(".", "-").lower()

    @property
    def required(self):
        return any(
            rule.errorCode == "required" and not rule.conditions for rule in self.rules
        )

    @property
    def conditionally_required(self):
        return not self.required and any(
            rule.errorCode == "required" for rule in self.rules
        )

    @property
    def conditional_values(self):
        return any(
            rule.conditions and rule.errorCode in {"equal_to", "one_of"}
            for rule in self.rules
        )


class ObjectSchema(BaseModel):
    properties: list[Property]
    doc: str = ""


APIDocs = dict[str, dict[str, ObjectSchema]]


def parse_api(data: bytes) -> APIDocs:
    api = TypeAdapter(APIDocs).validate_json(data)
    if not api:
        raise ValueError("The schema manifest contains no versions")
    slugs = set()
    for version, objects in api.items():
        slug = SchemaDocumentation.version_slug(version)
        if slug in slugs:
            raise ValueError(f"Duplicate schema version path: {slug}")
        slugs.add(slug)
        if not objects:
            raise ValueError(f"The schema manifest contains no objects for {version}")
        for kind, schema in objects.items():
            if not re.fullmatch(r"[A-Za-z][A-Za-z0-9]*", kind):
                raise ValueError(f"Invalid schema object kind: {kind}")
            paths = [prop.path for prop in schema.properties]
            if paths.count("$") != 1 or len(paths) != len(set(paths)):
                raise ValueError(
                    f"Invalid or duplicate property paths for {version}/{kind}"
                )
    return dict(
        sorted(api.items(), key=lambda item: SchemaDocumentation.version_order(item[0]))
    )


def rule_text(value: str) -> Markup:
    """Escape rule text while retaining quoted literals and HTTP links."""
    pattern = r"(?<!\w)'([^'\n]+)'(?!\w)|(https?://[^\s<>]+)"
    parts = []
    offset = 0
    for match in re.finditer(pattern, value):
        parts.append(escape(value[offset : match.start()]))
        if match.group(1) is not None:
            parts.append(Markup("<code>{}</code>").format(match.group(1)))
        else:
            url = match.group(2).rstrip(".,)")
            parts.append(Markup('<a href="{}">{}</a>').format(url, url))
            parts.append(escape(match.group(2)[len(url) :]))
        offset = match.end()
    parts.append(escape(value[offset:]))
    return Markup("").join(parts)


def table_cell(value: str) -> Markup:
    # Markdown still parses punctuation inside inline HTML code elements.
    entities = {ord(char): f"&#{ord(char)};" for char in "\\*_[]`|"}
    return Markup(str(escape(value)).translate(entities).replace("\n", "<br>"))


def code_block(value: str) -> str:
    fence = "`" * max(
        3, max((len(run) for run in re.findall(r"`+", value)), default=0) + 1
    )
    return f"{fence}text\n{value}\n{fence}"


def property_path(value: str) -> Markup:
    return Markup(".<wbr>").join(table_cell(part) for part in value.split("."))


class SchemaDocumentation:
    def __init__(self, api_path: Path = ROOT / "api.json"):
        self.api = parse_api(api_path.read_bytes())
        self.links = json.loads(
            (ROOT / "property-links.json").read_text(encoding="utf-8")
        )
        self.symbol_links = self.build_symbol_links()
        self.templates = Environment(
            loader=FileSystemLoader(ROOT / "templates"),
            undefined=StrictUndefined,
            keep_trailing_newline=True,
            trim_blocks=True,
            lstrip_blocks=True,
        )
        self.templates.filters["rule_text"] = rule_text
        self.templates.filters["table_cell"] = table_cell
        self.templates.filters["code_block"] = code_block
        self.templates.filters["property_path"] = property_path

    @staticmethod
    def version_slug(version: str):
        if not re.fullmatch(r"openslo(?:\.com)?/v\d+(?:(?:alpha|beta)\d*)?", version):
            raise ValueError(f"Invalid schema API version: {version}")
        return version.rsplit("/", 1)[1]

    @staticmethod
    def version_order(version: str):
        slug = SchemaDocumentation.version_slug(version)
        match = re.fullmatch(r"v(\d+)(?:(alpha|beta)(\d*))?", slug)
        major, stage, revision = match.groups()
        return (
            stage is not None,
            -int(major),
            {None: 0, "beta": 1, "alpha": 2}[stage],
            -int(revision or 0),
        )

    def object_schema(self, version: str, kind: str) -> ObjectSchema:
        try:
            return self.api[version][kind]
        except KeyError as error:
            raise ValueError(f"Unknown schema object: {version}/{kind}") from error

    def build_symbol_links(self):
        targets = {}
        for version, objects in self.api.items():
            slug = self.version_slug(version)
            version_links = self.links.get(version, {})
            packages = set()
            symbols = {}
            for kind, schema in objects.items():
                root = next(prop for prop in schema.properties if prop.path == "$")
                if root.typeInfo.package:
                    packages.add(root.typeInfo.package)
                    symbols[f"{root.typeInfo.package}#{root.typeInfo.name}"] = (
                        f"{kind.lower()}.md"
                    )
                for prop in schema.properties:
                    reference = version_links.get("_types", {}).get(prop.typeInfo.name)
                    if reference and prop.typeInfo.package:
                        symbols[f"{prop.typeInfo.package}#{prop.typeInfo.name}"] = (
                            reference["link"]
                        )
            for symbol, link in version_links.get("_symbols", {}).items():
                if "#" in symbol:
                    symbols[symbol] = link
                else:
                    for package in packages:
                        symbols[f"{package}#{symbol}"] = link
            for symbol, link in symbols.items():
                target = urlsplit(link)
                targets.setdefault(symbol, {})[slug] = target._replace(
                    path=posixpath.normpath(f"schema/{slug}/{target.path}")
                )
        return targets

    def symbol_link(self, url: str, source_path: str):
        parts = urlsplit(url)
        if parts.netloc != "pkg.go.dev":
            return url
        targets = self.symbol_links.get(
            f"{parts.path.lstrip('/')}#{parts.fragment}", {}
        )
        slug = next(
            (
                slug
                for slug in targets
                if source_path == f"schema/{slug}.md"
                or source_path.startswith(f"schema/{slug}/")
            ),
            None,
        )
        target = targets.get(slug)
        if target is None and len(set(targets.values())) == 1:
            target = next(iter(targets.values()))
        if target is None:
            return url
        return urlunsplit(
            target._replace(
                path=posixpath.relpath(target.path, posixpath.dirname(source_path))
            )
        )

    def object_description(self, version: str, kind: str):
        schema = self.object_schema(version, kind)
        root = next(prop for prop in schema.properties if prop.path == "$")
        return root.typeDoc or schema.doc

    def render_properties(self, properties, links=None):
        return self.templates.get_template("properties.md.j2").render(
            properties=properties, links=links or {}
        )

    def object_properties(self, version: str, kind: str):
        schema = self.object_schema(version, kind)
        version_links = self.links.get(version, {})
        path_links = {**version_links.get("_common", {}), **version_links.get(kind, {})}
        type_links = version_links.get("_types", {})
        rendered_links = {}
        for prop in schema.properties:
            link = path_links.get(prop.display_path, type_links.get(prop.typeInfo.name))
            if not link:
                continue
            target = urlsplit(link["link"])
            # Expand the canonical definition on its own page.
            if target.path == f"{kind.lower()}.md" and target.fragment == prop.anchor:
                continue
            rendered_links[prop.path] = {
                "url": link["link"],
                "description": self.templates.from_string(link["template"]).render(
                    path=prop.display_path,
                    typeInfo=prop.typeInfo,
                    required=prop.required,
                    link=link["link"],
                ),
            }
        properties = [
            prop
            for prop in schema.properties
            if not any(
                prop.path.startswith(path + ".") or prop.path.startswith(path + "[")
                for path in rendered_links
            )
        ]
        return self.render_properties(properties, rendered_links)

    def metadata_properties(self, version: str):
        objects = self.api[version]
        schema = objects.get("Service", next(iter(objects.values())))
        properties = [
            prop
            for prop in schema.properties
            if prop.path == "$.metadata" or prop.path.startswith("$.metadata.")
        ]
        if not properties:
            raise ValueError(f"No metadata properties found for {version}")
        return self.render_properties(properties)

    def version_overview(self, version: str):
        return self.templates.get_template("version.md.j2").render(
            version=version,
            slug=self.version_slug(version),
            objects=self.api[version],
            metadata=self.metadata_properties(version),
        )

    def version_links(self):
        return "\n".join(
            f"- [`{version}`](schema/{self.version_slug(version)}.md)"
            for version in self.api
        )


class SchemaLinksExtension(Extension):
    def __init__(self, schema: SchemaDocumentation):
        super().__init__()
        self.schema = schema
        self.source_path = ""

    def extendMarkdown(self, md):
        # Resolve symbols before MkDocs validates and rewrites relative links.
        md.treeprocessors.register(
            SchemaLinksTreeprocessor(md, self.schema, self.source_path),
            "schema_links",
            1,
        )


class SchemaLinksTreeprocessor(Treeprocessor):
    def __init__(self, md, schema: SchemaDocumentation, source_path: str):
        super().__init__(md)
        self.schema = schema
        self.source_path = source_path

    def run(self, root):
        for element in root.iter("a"):
            if (url := element.get("href")) is not None:
                element.set("href", self.schema.symbol_link(url, self.source_path))


def main():
    parser = argparse.ArgumentParser(description="Import an SDK schema manifest.")
    parser.add_argument("manifest", type=Path)
    args = parser.parse_args()
    data = args.manifest.read_bytes()
    parse_api(data)
    (ROOT / "api.json").write_bytes(data)


if __name__ == "__main__":
    main()
