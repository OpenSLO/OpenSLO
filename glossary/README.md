# Glossary

Overtime, the SLO community has created specific terminology. Since not every user coming to project like OpenSLO knows
exactly what everything means we have created a glossary. This glossary should help newcomers understand things that we
discuss here and make contributions easier.

In addition, I have created a section with literature that I used to write these definitions. There is a need for a place
where we can write down definitions to facilitate easier discussions about subjects adn share differing point of view.

## Definitions

Service Level Indicator(SLI) – a carefully defined quantitative measure of some aspect of the level of service that is
provided.

Service Level Objective(SLO) – a target value or range of values for a service level that is measured by an SLI.

Service Level Agreements(SLA) – an explicit or implicit contract with your users that includes consequences of
meeting (or missing) the SLOs they contain.

Time window – time in which reliability of service is measured and calculated.

Error budget – a measurement that allows us to see how much SLO can be violated in a given period of time.

Occurrences - error budget calculation method. It is calculated by counting the number of good events in the total count
of the events. Example: Count if total number of measured latencies, measured in last two weeks, was below 100 ms for 95%
of measurements.

TimeSlices - error budget calculation method. This method uses ratio of good slices to total slices in budgeting
period. Example: How many minutes, in last two weeks, your service had latency below 100ms in measurements in these minutes.

RatioTimeSlices - error budget calculation method. It is calculated by taking average of good events to total
events in all slices. Example: What is average percent of good events to total events in every minute for last two weeks.

Ratio Metrics - is an SLI composed of a relation of two metrics. Most commonly called good (numerator) and total
events (denominator). When picking objective for your service you specify how big part of total events should be good events.

Threshold Metrics - SLI where metrics returned form a query are compared to an arbitrarily chosen threshold. A comparator
can be lt(<), le(<=), gt(>), ge(>=). Based on the result of a comparison is determined if the error budget is burnt or not.

Composite SLO - it is SLO that is composed of few different objectives that can be seen as SLOs as well.
Each objective that is part of composite SLO can have its own threshold, query and data source. Composite SLO burns its
error budget when any of SLOs beneath burns their error budget. User should be able to specify weight that will decide
if composite SLO is burning faster or slower depending on which inner SLO is burning its error budget.

## Literature

Books that are great to learn about SLO concept:

- [Google's SRE book](https://sre.google/sre-book/table-of-contents/)
- [Implementing Service Level Objectives](https://www.oreilly.com/library/view/implementing-service-level/9781492076803/)
