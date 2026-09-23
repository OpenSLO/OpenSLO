(() => {
  const tocSidebarVisible = window.matchMedia("(min-width: 60em)");
  let propertyControls;

  function openProperty(hash) {
    let id;
    try {
      id = decodeURIComponent(hash.slice(1));
    } catch {
      return;
    }
    const heading = document.getElementById(id);
    const panel = heading?.nextElementSibling;
    if (heading?.matches("h3, h4") && panel?.matches("details.property")) {
      panel.open = true;
    }
  }

  function placePropertyControls() {
    if (!propertyControls || propertyControls.hidden) {
      return;
    }
    const tocTitle = tocSidebarVisible.matches
      ? document.querySelector(".md-sidebar--secondary:not([hidden]) .md-nav__title")
      : null;
    if (tocTitle) {
      let heading = tocTitle.closest(".property-toc-heading");
      if (!heading) {
        heading = document.createElement("div");
        heading.className = "property-toc-heading";
        tocTitle.before(heading);
        heading.append(tocTitle);
      }
      heading.append(propertyControls);
    } else {
      document.querySelector(".md-header__title")?.after(propertyControls);
    }
  }

  document$.subscribe(() => {
    // Controls live outside the replaced content during instant navigation.
    propertyControls?.remove();
    propertyControls = document.querySelector(".md-content .property-controls");
    const properties = [...document.querySelectorAll("details.property")];
    if (propertyControls) {
      propertyControls.hidden = properties.length === 0;
      const toggle = propertyControls.querySelector('[data-property-action="toggle"]');
      function updateToggle() {
        const expanded = properties.every(property => property.open);
        const label = expanded ? "Collapse all properties" : "Expand all properties";
        toggle.setAttribute("aria-expanded", String(expanded));
        toggle.setAttribute("aria-label", label);
        toggle.title = label;
      }
      for (const property of properties) {
        property.addEventListener("toggle", updateToggle);
      }
      for (const button of propertyControls.querySelectorAll("button[data-property-action]")) {
        button.onclick = () => {
          const expand = button.dataset.propertyAction === "expand" ||
            (button.dataset.propertyAction === "toggle" && !properties.every(property => property.open));
          for (const property of properties) {
            property.open = expand;
          }
          updateToggle();
        };
      }
      updateToggle();
      placePropertyControls();
    }
    openProperty(location.hash);
  });

  tocSidebarVisible.addEventListener("change", placePropertyControls);

  window.addEventListener("hashchange", () => openProperty(location.hash));
  document.addEventListener("click", (event) => {
    if (event.button !== 0 || event.ctrlKey || event.metaKey || event.shiftKey || event.altKey) {
      return;
    }
    const link = event.target.closest("a[href]");
    if (!link || link.target || link.hasAttribute("download")) {
      return;
    }
    const url = new URL(link.href);
    if (url.origin === location.origin && url.pathname === location.pathname && url.search === location.search) {
      openProperty(url.hash);
    }
  }, true);
})();
