# Glossary

SLO community created over the time specific terminology, not every user coming to projects like OpenSLO knows exactly what
everything means. This glossary should help newcomers to understand things we discuss here and make contribution easier.
In addition, I created section with literature that I used to write these definitions.
Another thing is that we need a place where we can write down definitions to have easier discussions about
subjects that can be understood differently.

## Definitions

Service Level Indicator(SLI) – a carefully defined quantitative measure of some aspect of the level of service that is provided.

Service Level Objective(SLO) – a target value or range of values for a service level that is measured by an SLI.

Service Level Agreements(SLA) – an explicit or implicit contract with your users that includes consequences of
meeting (or missing) the SLOs they contain.

Time window – time in which reliability of service is measured and calculated.

Error budget – a measurement that allows us to see how much SLO can be violated in a given period of time.

Occurrences - error budget calculation method. It is calculated by counting the number of good events in the total count
of the events. Example: Count if total number of measured latencies, measured in last two weeks, was below 100 ms for 95% of
measurements. 

TimeSlices - error budget calculation method. This method uses ratio of good slices to total slices in budgeting
period. Example: How many minutes, in last two weeks, your service had latency below 100ms in measurements in these minutes. 

RatioTimeSlices - error budget calculation method. It is calculated by taking average of good events to total
events in all slices. Example: What is average percent of good events to total events in every minute for last two weeks.

-- Next iteration

Reliability Burn Down

Ratio Metrics

Threshold Metrics

Composite SLO

Aggregate SLO

## Literature

Books that are great to learn about SLO concept:

- [Google's SRE book](https://sre.google/sre-book/table-of-contents/)
- [Implementing Service Level Objectives](https://www.oreilly.com/library/view/implementing-service-level/9781492076803/)
