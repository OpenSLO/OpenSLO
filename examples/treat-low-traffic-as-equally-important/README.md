# Threat low traffic as equally important

Let's imagine application that performance degration even during low traffic periods (a handful of users are not satisfied)
should count equally with degration during peak hours when hundreds of users are not satisfied with application.
For instance, a system for reporting incidents or calling for help should operate in such a way. One person than can't get
is important too.

## Which OpenSLO budgeting method should be used?

Timeslices fits perfectly. In this method what we measure is how many good minutes (minutes where the system is operating
within defined boundaries) were observed, compared to the total number of minutes in the window. With this approach, a bad
minute that occurs during a low-traffic period will have the same effect on your SLO as a bad minute caused by the fact that
the platform is overloaded with traffic.

## Concrete example

Assumptions:

- Timewindow: 1 week calendar (of course rolling can be used too)
- Objective is: 99.95% of successful requests.
- Timeslice allowance (target) is 95%.
- Timeslice is 1 minute.

Number of minutes in the Timewindow: `1 week = 7 days`, `7 days = 168 hours`, `168 hours = 10 080 minutes`
Error Budget (how many minutes can be considered as bad): `(100% - 99.95%) * 10 080 minutes = 5 minutes 2.4 seconds`

Minute will be marked as bad, when just only `5` request fails and total number of requests is less than `100` in that minute, so it's very strict.
Having over five incidents like that during one Timewindow will violate SLO and should trigger meaningful conversation in organization about
reliability of the product.
