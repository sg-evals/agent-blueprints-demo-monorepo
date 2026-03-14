# Incident Response Runbook

## Alerts

| Alert Name                  | Severity | Source         | Condition                                  |
|-----------------------------|----------|----------------|--------------------------------------------|
| HighErrorRate               | P1       | api-gateway    | 5xx rate > 5% for 5 minutes               |
| PaymentProcessingFailure    | P1       | billing-service| Payment failures > 10 in 2 minutes         |
| OrderQueueBacklog           | P2       | orders-service | Queue depth > 1000 for 10 minutes          |
| NotificationDeliveryDelay   | P2       | notifications  | Delivery latency p99 > 30s                |
| WorkerIngestLag             | P3       | worker-ingest  | Consumer lag > 5000 messages               |
| SearchIndexStaleness        | P3       | search-service | Index age > 15 minutes                     |

## Triage Steps

1. **Identify the affected service.** Check the alert source and confirm which service is impacted using the observability dashboards.
2. **Check recent deployments.** Review the last 30 minutes of deployments via CI/CD logs. If a deploy correlates with the incident, consider a rollback.
3. **Inspect logs.** Use the centralized logging system to filter by service name and time window. Look for stack traces, connection errors, or upstream timeouts.
4. **Verify dependencies.** Confirm that downstream services and databases are healthy. Check the service dependency diagram in `docs/architecture/overview.md`.
5. **Reproduce if possible.** For non-P1 incidents, attempt to reproduce in staging before applying fixes.

## Escalation

| Level   | Criteria                                   | Contact              | Response SLA |
|---------|--------------------------------------------|----------------------|--------------|
| Level 1 | On-call engineer                           | PagerDuty rotation   | 5 minutes    |
| Level 2 | Service owner / team lead                  | Slack #incidents     | 15 minutes   |
| Level 3 | Engineering manager + SRE lead             | Phone escalation     | 30 minutes   |
| Level 4 | VP Engineering (customer-facing P1 only)   | Direct phone         | 1 hour       |

### Escalation Rules

- **P1**: Start at Level 1. Escalate to Level 2 if not acknowledged within 10 minutes. Escalate to Level 3 if not mitigated within 30 minutes.
- **P2**: Start at Level 1. Escalate to Level 2 if not mitigated within 1 hour.
- **P3**: Handle during business hours. No automatic escalation.

## Recovery

1. **Apply the fix.** Deploy the fix or roll back the offending change.
2. **Verify recovery.** Confirm that error rates return to baseline and alerts resolve.
3. **Drain queues if needed.** For worker services, check that any backlogged messages are being processed.
4. **Communicate.** Post a summary in #incidents with: root cause, impact duration, and fix applied.
5. **Schedule a post-mortem.** For P1 and P2 incidents, schedule a blameless post-mortem within 48 hours. Use the post-mortem template in the wiki.
6. **Create follow-up tickets.** Track any hardening work (improved alerts, circuit breakers, tests) as follow-up issues.
