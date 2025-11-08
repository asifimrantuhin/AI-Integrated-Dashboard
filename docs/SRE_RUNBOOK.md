# Site Reliability Engineering Runbook

This runbook captures operational expectations for keeping iDash reliable, observable, and secure across environments.

---

## 1. Service Topology

| Service        | Stack                     | Ports | Dependencies |
|----------------|---------------------------|-------|--------------|
| Frontend       | React + Vite (static)     | 5173  | Backend API  |
| Backend API    | Go + Echo                 | 8080  | MySQL, AI API, Redis (opt.) |
| Admin CMS      | Laravel + Blade           | 8001  | MySQL, Redis (queues) |
| AI Service     | FastAPI + Prophet         | 8000  | MySQL (read-only) |
| MySQL Cluster  | Managed (Aurora/MySQL)    | 3306  | n/a          |
| Redis (cache)  | ElastiCache/Redis OSS     | 6379  | n/a          |

All ingress traffic terminates at an HTTPS reverse proxy (Nginx / ALB) which routes to the respective services.

---

## 2. Monitoring & Telemetry

### Metrics
- **Backend API**: request rate, latency p50/p95, error rate, goroutine count, DB connection pool usage.
- **AI Service**: forecast latency, job success rate, queue depth, memory footprint.
- **Admin CMS**: queue length, job failures, artisan command duration.
- **Database**: CPU, buffer pool hit ratio, replication lag, slow queries.
- **Synthetic checks**: `/api/health`, `/health` endpoints every 60 seconds.

Instrument via Prometheus or Datadog exporters:
- Go: `promhttp` middleware
- FastAPI: `prometheus-fastapi-instrumentator`
- Laravel: `spatie/laravel-dashboard` or custom metrics route

### Logs
- Ship structured JSON logs to centralized store (ELK / OpenSearch / Datadog logs).
- Tag logs with environment, service, request ID (`x-request-id` header already injected by frontend and backend).

### Traces
- Adopt OpenTelemetry SDK in Go, Python, and Node to trace cross-service calls. Export to Jaeger/Tempo.

---

## 3. Alerting Policy

| Alert                       | Threshold                               | Action |
|-----------------------------|-----------------------------------------|--------|
| API 5xx rate                | >1% over 5 minutes                      | Page on-call, investigate release |
| Forecast job duration       | >120s p95 for 3 consecutive intervals   | Scale AI service, review dataset size |
| DB replication lag          | >30s for primary replica                | Check network, failover if primary unstable |
| Queue backlog               | >100 pending jobs for >10 minutes       | Scale workers, inspect stuck jobs |
| SSL certificate expiry      | <14 days remaining                      | Rotate certificate |

Alerting toolset: PagerDuty (critical), Slack Ops channel (warnings), email (informational).

---

## 4. Incident Response Workflow

1. **Triage**
   - Confirm alert validity (metrics + logs).
   - Determine blast radius (tenants, modules, environment).
2. **Mitigate**
   - Roll back recent deploy via Git tag / artifact redeploy.
   - Scale affected service horizontally if resource saturation.
   - Toggle feature flags if specific module causing errors.
3. **Communicate**
   - Update incident channel (#idash-incident) with status, ETA.
   - Notify stakeholders (product, support) every 30 minutes.
4. **Resolve**
   - Verify metrics back to baseline for 2 consecutive intervals.
   - Close incident in PagerDuty with resolution notes.
5. **Post-Incident Review**
   - Schedule blameless PIR within 48 hours.
   - Capture root cause, contributing factors, permanent fixes.

---

## 5. Scaling & Capacity

- **Backend API**: Auto-scale based on CPU >60% or requests/sec > 500. Pods target 2 vCPU / 4GB RAM.
- **AI Service**: CPU-heavy; use HPA on CPU + queue length (`ai_forecast_jobs`). Consider GPU-enabled node pool for large forecasts.
- **Admin CMS**: Scale horizontally for ETL spikes; offload heavy tasks to queues + worker pool.
- **Database**: Enable read replicas for analytics and AI; use connection pooling (max 30 per service).

Perform quarterly load tests (k6) covering peak workloads (dashboard refresh, forecast generation, data sync).

---

## 6. Backups & Disaster Recovery

| Asset         | Strategy                                      | RPO | RTO |
|---------------|-----------------------------------------------|-----|-----|
| MySQL         | Automated daily snapshots + PITR; replicate to secondary region | 5 min | 60 min |
| Redis         | AOF persistence with 15-minute snapshots      | 15 min | 30 min |
| Object Storage| Versioned buckets for file uploads            | 1 hour | 4 hours |
| Config        | Git + Terraform state stored in remote backend| 5 min | 30 min |

Disaster recovery runbooks: rehearse semi-annually. Ensure infrastructure as code can recreate environment in alternate region.

---

## 7. Security & Compliance

- Enforce HTTPS everywhere; HSTS via reverse proxy.
- Rotate secrets every 90 days; integrate with IAM roles where possible.
- Apply database encryption at rest and in transit (TLS).
- Log-admin actions in Admin CMS; ship to SIEM.
- Run SCA/DAST weekly (GitHub Advanced Security, OWASP ZAP pipeline).

---

## 8. Change Management

- All production deploys require:
  1. Passing automated tests (unit + e2e smoke).
  2. Staging soak time ≥2 hours.
  3. Change record (Jira ticket) referencing migration/feature.
- Enable feature flags for risky modules (e.g., prescriptive engine) to allow canary rollout.

---

## 9. Operational Runbooks

### Backend API Outage
1. Verify `/api/health` failing.
2. Check Go service logs for panic; inspect recent deploy.
3. Restart pods/process via orchestrator.
4. If persistent, roll back to previous artifact.

### AI Forecast Failures
1. Check AI logs for `Sales forecast generation failed` messages.
2. Validate database connectivity (read-only user).
3. Run smoke forecast via `POST /api/forecast/sales` with known payload.
4. If Prophet exceptions, purge corrupted cache, consider rerunning demo data import.

### Admin CMS Queue Jam
1. Inspect `jobs` table for failing entries.
2. Restart queue workers (`php artisan queue:restart`).
3. Check external API credentials if sync tasks failing.

---

## 10. Glossary

- **RPO**: Recovery Point Objective – acceptable data loss window.
- **RTO**: Recovery Time Objective – acceptable downtime window.
- **PIR**: Post-Incident Review.
- **HPA**: Horizontal Pod Autoscaler.

Maintain and iterate on this runbook as architecture evolves or new modules are introduced.
