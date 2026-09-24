---
title: "OpenSLO Community Meeting: June 2024"
slug: openslo-project-meeting-june-2024
date: 2024-06-04
draft: false
authors:
  - Mateusz Hawrus
categories:
  - Meetings
---

The June meeting focused on multiple alert conditions, calendar-aligned windows, and Kubernetes compatibility.

<!-- more -->

## Summary

1. **Multiple alert conditions**:
   The group agreed to propose lifting the single-condition restriction, with AND as the default relationship. More complex expressions and vendor-specific extensions would remain a separate discussion. [Discussion at 26:11](https://www.youtube.com/watch?v=fpaw3tTw_pE&t=1571s).

2. **Calendar-aligned windows**:
   Participants examined week boundaries and ambiguous month durations. They discussed clearer definitions and possible SDK helpers, but left the design open for further review. [Discussion at 33:44](https://www.youtube.com/watch?v=fpaw3tTw_pE&t=2024s).

3. **Kubernetes and the SDK**:
   The v2 draft included an API version change for Kubernetes compatibility. Labels still needed attention. The group reaffirmed the plan to use Go definitions as the source for validation and generated schemas. [Status at 4:48](https://www.youtube.com/watch?v=fpaw3tTw_pE&t=288s), [SDK discussion at 46:14](https://www.youtube.com/watch?v=fpaw3tTw_pE&t=2774s).

## Action Items

1. Open a separate pull request for multiple alert conditions with AND semantics.
2. Continue exploring expressions and vendor-specific alert settings in the existing proposal.
3. Review the calendar-window proposal and share implementation concerns in the issue.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/fpaw3tTw_pE" title="OpenSLO community meeting, June 4, 2024" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
