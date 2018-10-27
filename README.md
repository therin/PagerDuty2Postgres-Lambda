# PagerDuty2Postgres Lambda

Pagerduty2Postgres Lambda imports data from the PagerDuty API into a Postgres database for easy querying and analysis.

Lambda is written in GO and deployed as a single binary.


Main features include:

Collect summary statistics about on-call activity.
Calculate per-user, per-service, per-escalation-policy on-call metrics.
Determine the frequency of on-hours vs. off-hours pages.
Produce custom on-call reports with incident-level detail.
Back-test proposed on-call changes.
Perform one-off queries against historical pager data.


Solution overview

