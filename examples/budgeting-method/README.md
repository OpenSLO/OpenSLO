# Budgeting method

OpenSLO specifies three types of budgeting methods:

- Occurrences
- Timeslices
- RatioTimeslices

Here's a brief description of them and examples of use cases for each method.

## Occurrences

Method uses a ratio of good to total counts of events. This is the most common and easy to understand method,
all occurrences of events have the same impact on error budget, in a lot of cases it is used with SLIs that measure
things that are connected to traffic in your services. Method automatically weights impact by the total number
of requests served, so it will give an accurate reflection of availability of service.

## Timeslices

Method uses a ratio of good time slices to total time slices. Whole timewindow is divided into timeslices in size
specified by user. Each timeslice have calculated ratio of good to total events. When the ratio is below given
timeSliceTarget we count it as bad minute and subtract it from error budget. This method is useful to reflect SLA which
often says "Service needs to be available 99% of time". That can be understood as 99% of minutes have 90% of successful
requests.

## RatioTimeslices

Method uses an average of all time slices' success ratios. Here as in timeslices method timewindow is divided into parts.
The average ratio of success to total counts determines if the error budget is burned. This method is suitable when the
monitored service is exposed to an unpredictable environment and needs to be safe from rare situations when the whole error
budget could be burned in minutes. For example let's imagine some service that works well and has on average 1 million of
requests weekly. If some malicious group decided to prepare a DDoS attack and overflow our service with thousands of requests
for a few minutes we would burn the whole error budget but service was inaccessible only for a few minutes and users
didn't feel the reduction of quality. This method is something between occurrences and timeslices we see attack on SLO
it is counted as more impactful than just few minutes of poor quality of service but not that impactful to burn whole
error budget.

### Quick calculations

How would DDoS attack impact your service for different budgeting methods.

Assumptions:

- Timewindow: 2 weeks rolling
- Average daily requests: 24000
- Your service by average works with: 99% of success responses.
- Objective is: 95% of successful requests.
- Timeslice is 1 minute (if it exists).
- By request, I mean here anything that can end with success or failure. There is no connection with real indicators.
- Your service was under attack by a malicious group that decided to DDoS your service making it work with 1% of availability.

#### Occurrences burned budget

In two weeks you have 336000 requests. If a malicious group decides to attack your service with DDoS it will take 16632
failed requests to burn your whole error budget. Depending on the attack it can take different time but shouldn't be long.

#### Timeslices burned budget

Two weeks timewindow have 20160 timeslices, to burn your budget attack needs to last for 16 hours and 48 minutes.

#### RatioTimeslices burned budget

To burn your budget you need to last for about 13 hours and 45 minutes.
