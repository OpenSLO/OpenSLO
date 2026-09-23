from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[1]


class SpecificationImportTests(unittest.TestCase):
    def test_import_adapts_version_links_and_rejects_an_incomplete_source(self):
        with tempfile.TemporaryDirectory(prefix="openslo-specification-test-") as name:
            project = Path(name)
            (project / "docs").mkdir()
            script = project / "specification.py"
            script.write_bytes((ROOT / "specification.py").read_bytes())
            source = project / "README.md"
            source.write_text(
                "# OpenSLO\n\n## Introduction\n\nIntroduction.\n\n## Specification\n\n"
                + "\n\n".join(
                    f"> [Check work in progress for v2.](enhancements/v2alpha.md#{kind})"
                    for kind in ("datasource", "slo", "sli")
                )
            )
            result = subprocess.run(
                [sys.executable, str(script), str(source)],
                capture_output=True,
                text=True,
                timeout=10,
            )
            self.assertEqual(result.returncode, 0, result.stderr)
            output = project / "docs/specification.md"
            content = output.read_text()
            for kind in ("datasource", "slo", "sli"):
                self.assertIn(f"(schema/v1/{kind}.md)", content)
                self.assertIn(f"/enhancements/v2alpha.md#{kind})", content)
                self.assertNotIn(f"(schema/v2alpha/{kind}.md)", content)
            self.assertIn("The v1 specification makes `timeWindow` optional.", content)
            self.assertIn(
                "The Go SDK currently requires exactly one time window.", content
            )
            source.write_text("# Unrelated README\n")
            result = subprocess.run(
                [sys.executable, str(script), str(source)],
                capture_output=True,
                text=True,
                timeout=10,
            )
            self.assertNotEqual(result.returncode, 0)
            self.assertIn("must contain Introduction and Specification", result.stderr)
            self.assertEqual(output.read_text(), content)
