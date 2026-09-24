---
title: "OpenSLO Community Meeting: October 2021"
slug: openslo-project-meeting-october-2021
date: 2021-10-21
draft: false
authors:
  - Mateusz Hawrus
categories:
  - Meetings
---

The October meeting narrowed the first alerting proposal and discussed error metrics, examples, and reusable SLIs.

<!-- more -->

## Summary

1. **A smaller alerting proposal**:
   Weyert and Mike simplified the proposal around burn-rate alerts. Notification targets would use names and metadata, with implementation details left to the receiving system. Other alert types could follow later. [Discussion at 1:23](https://www.youtube.com/watch?v=dhZwnlb8ffY&t=83s).

2. **Good and bad event counts**:
   A proposal would allow bad events and a total count as an alternative to good events and a total count. Participants favored separate guidance and examples to explain the tradeoffs. [Discussion at 7:50](https://www.youtube.com/watch?v=dhZwnlb8ffY&t=470s).

3. **Reusable definitions**:
   Separating SLIs from SLOs could support shared templates and changes of data source. The group also revisited a demo that users could run themselves. [Discussion at 47:02](https://www.youtube.com/watch?v=dhZwnlb8ffY&t=2822s).

## Action Items

1. Give the alerting proposal a final review before merging it.
2. Write guidance on good and bad event counts in a separate repository document.
3. Investigate a shareable demo and version the specification for substantial changes.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/dhZwnlb8ffY" title="OpenSLO community meeting, October 21, 2021" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
