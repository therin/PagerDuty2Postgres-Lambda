
Pagerduty2Postgres imports data from the PagerDuty API into a Postgres database for easy querying and analysis.

This is an almost identical copy of Stripe's amazing [pd2pg](https://github.com/stripe/pd2pg) solution. Unfortunately, they wrote it in Ruby (yikes) and thus it's undeployable as AWS Lambda (yet?). Also Go is amazing.

This app is designed to be run on AWS Lambda service, is written in GO and deployed as a single binary.
Cloudformation deployment code is included in this repo.

### Main features
- Collect summary statistics about on-call activity.
- Calculate per-user, per-service, per-escalation-policy on-call metrics.
- Determine the frequency of on-hours vs. off-hours pages.
- Produce custom on-call reports with incident-level detail.
- Back-test proposed on-call changes.
- Perform one-off queries against historical pager data.

### Solution infrastructure overview:

![pagerduty2postgres lambda](https://user-images.githubusercontent.com/2115124/47610311-a90a7680-daae-11e8-8a5b-1259091caf16.jpeg)
