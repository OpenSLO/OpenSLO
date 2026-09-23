---
title: "OpenSLO Community Meeting: February 2024"
slug: openslo-project-meeting-february-2024
date: 2024-02-21
draft: false
authors:
  - Pawel Rusakiewicz
categories:
  - Meetings
---

First community meeting in 2024 :tada:

<!-- more -->

## Summary

During the meeting, several agenda items were discussed and decisions were made regarding the OpenSLO project.

1. **Source of Truth for OpenSLO manifest**:
   The group decided to use Go structs as the source of truth for the OpenSLO manifest. This decision was made to ensure more control over validation. The plan is to generate JSON Schema/CRDs based on Go structs.

2. **Labels Support**:
   Labels support is to be extended for all objects' metadata, including adding labels to Objective.

3. **Oslo Conversion Removal**:
   The decision was made to remove the conversion entirely and extend [oslo](https://github.com/openslo/oslo) documentation on the usage of annotations for vendor-specific conversions. The focus will be on improving the SDK instead of expanding [oslo's](https://github.com/openslo/oslo) responsibilities.

## Action Items

Several action items were noted, including:

1. Documenting new use cases of OpenSLO.
2. Ensuring full Kubernetes compatibility.
3. Considering bringing OpenSLO under CNCF's umbrella.

## Recording

<div class="video-wrapper">
  <iframe width="560" height="315" src="https://www.youtube.com/embed/fpaw3tTw_pE?si=5MpyKo6xmnyIBOgJ" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
</div>
