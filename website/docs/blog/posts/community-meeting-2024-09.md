---
title: "OpenSLO Community Meeting: September 2024"
slug: openslo-project-meeting-september-2024
date: 2024-09-01
draft: false
authors:
  - Mateusz Hawrus
categories:
  - Meetings
---

The September meeting reviewed SDK progress and resolved several questions for the v2 draft.

<!-- more -->

## Summary

1. **Go SDK**:
   Resource definitions were in place, with validation and Kubernetes integration underway. The SDK had no official release yet. [Discussion at 15:47](https://www.youtube.com/watch?v=Di98rlDoVno&t=947s).

2. **Time windows**:
   Existing start times could define weekly calendar boundaries. The group also agreed to exclude months and larger units from rolling windows in v2. [Calendar discussion at 19:15](https://www.youtube.com/watch?v=Di98rlDoVno&t=1155s), [rolling windows at 42:00](https://www.youtube.com/watch?v=Di98rlDoVno&t=2520s).

3. **Specification changes**:
   Participants supported multiple threshold objectives and the name `serviceRef` for v2. They kept the proposed alert-condition expansion out of v1. The `total` versus `valid` terminology question remained open. [Objectives at 51:19](https://www.youtube.com/watch?v=Di98rlDoVno&t=3079s), [references at 1:08:36](https://www.youtube.com/watch?v=Di98rlDoVno&t=4116s).

## Action Items

1. Add an example of choosing a calendar week's start day.
2. Update the v2 draft and SDK for the agreed changes.
3. Clarify name validation and continue discussing SLI denominator terminology in GitHub.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/Di98rlDoVno" title="OpenSLO community meeting, September 2024" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
