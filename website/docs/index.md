---
hide:
  - navigation
  - toc
  - path
---

<!-- markdownlint-disable MD041 -->

![OpenSLO Logo](assets/images/openslo_logo.svg)
/// caption
<p style="text-align: center; color: gray; font-size: 1.0rem; line-height: 1.2;">
  Open Service Level Objective (SLO) Specification.<br>
  Designed to make SLOs ergonomic to modern<br>
  developer Git workflow.
</p>
///

# What is OpenSLO?

OpenSLO is a service level objective (SLO) language that declaratively defines
reliability and performance targets using a simple YAML specification.<br>
It is released under Apache 2.0 and we welcome contributions
from the reliability engineering ecosystem.

SLOs are reliability targets for services that allow organizations to make
better decisions in how to create, operate,
and run cloud services and applications.<br>
To learn more about SLOs, check out [SLOconf.com](https://sloconf.com/).

<div class="grid cards" markdown>

- <p align="center">
    <img src="assets/images/illustrations/icons/brackets.svg" width="50%" alt="">
  </p>

    <p style="text-align: center;">
      __Define SLOs as Code__
    </p>

    ---

    Create declarative definitions of SLOs that describe thresholds, metrics, and goals for your applications and infrastructure.

    [:octicons-arrow-right-24: Getting started](specification.md)

- <p align="center">
    <img src="assets/images/illustrations/icons/terminal.svg" width="50%" alt="">
  </p>

    <p style="text-align: center;">
      __Oslo CLI__
    </p>

    ---

    Validate OpenSLO definitions in your CI/CD workflows with the Oslo CLI.

    [:octicons-arrow-right-24: Reference](tools.md#oslo-cli)

- <p align="center">
    <img src="assets/images/illustrations/icons/nodes.svg" width="50%" alt="">
  </p>

    <p style="text-align: center;">
      __Vendor Agnostic Metrics__
    </p>

    ---

    OpenSLO is designed to be implementation neutral and allow multiple vendors and other projects to describe and share SLOs in a well-defined format.

    [:octicons-arrow-right-24: Customization](schema/v1/datasource.md)

</div>
