---
title: "OpenSLO Community Meeting: June 2021"
slug: openslo-project-meeting-june-2021
date: 2021-06-17
draft: false
authors:
  - Mateusz Hawrus
categories:
  - Meetings
---

The first OpenSLO community meeting explored the roadmap, reusable SLO definitions, and validation tooling.

<!-- more -->

## Summary

1. **Separate SLIs and SLOs**:
   Participants proposed separating data queries from reliability objectives. Teams could then change a data source without rewriting their objectives. The discussion also covered different ownership and access needs for the two definitions. [Discussion at 2:35](https://www.youtube.com/watch?v=86T_eL-aGmc&t=155s).

2. **Production and delivery use cases**:
   The group compared production error budgets with performance tests and deployment checks. Weighted objectives, shorter evaluation windows, and regression detection were candidates for the roadmap. [Discussion at 12:10](https://www.youtube.com/watch?v=86T_eL-aGmc&t=730s).

3. **Validation and starter definitions**:
   Participants discussed the boundary between manifest validation and provider-specific query checks. They also proposed an interactive generator and shared examples to help users create their first SLOs. These were ideas for future work. [Discussion at 29:24](https://www.youtube.com/watch?v=86T_eL-aGmc&t=1764s).

## Action Items

1. Collect the proposed use cases and tooling ideas in the project roadmap.
2. Clarify how implementations can extend the specification with additional fields.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/86T_eL-aGmc" title="OpenSLO community meeting, June 17, 2021" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
