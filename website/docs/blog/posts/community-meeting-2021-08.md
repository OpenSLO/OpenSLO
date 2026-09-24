---
title: "OpenSLO Community Meeting: August 2021"
slug: openslo-project-meeting-august-2021
date: 2021-08-19
draft: false
authors:
  - Mateusz Hawrus
categories:
  - Meetings
---

The August meeting focused on reusable alerting resources and metadata for OpenSLO implementations.

<!-- more -->

## Summary

1. **Reusable alerting resources**:
   Weyert presented a proposal to separate alert policies, conditions, and notification targets. The group discussed reusing these building blocks across SLOs and keeping notification settings flexible for each implementation. [Discussion at 13:31](https://www.youtube.com/watch?v=88GrJRi6rM8&t=811s).

2. **Alert behavior**:
   Participants compared burn-rate alerts, traffic spikes, low sample counts, and recovery signals. Multiple conditions and windows remained design questions for the proposal. [Discussion at 18:14](https://www.youtube.com/watch?v=88GrJRi6rM8&t=1094s).

3. **Metadata and labels**:
   Proposed metadata would carry extra context for generators and other tools. Labels could support filtering, grouping, and generated metrics. The distinction between these fields still needed discussion. [Discussion at 45:41](https://www.youtube.com/watch?v=88GrJRi6rM8&t=2741s).

## Action Items

1. Review the alerting proposal and suggest changes in the pull request.
2. Make notification parameters optional and refine recovery behavior.
3. Follow up with a labels proposal and feedback on the metadata fields.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/88GrJRi6rM8" title="OpenSLO community meeting, August 19, 2021" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
