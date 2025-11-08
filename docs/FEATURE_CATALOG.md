# Feature Catalog

A high-level inventory of capabilities delivered by iDash across functional modules.

---

## 1. Platform Foundations

- **Unified Authentication & RBAC**
  - JWT SSO across frontend, backend API, admin CMS
  - Role hierarchy (Executive, Manager, Analyst) with granular permissions
  - Multi-company support with contextual switching

- **AI Orchestration Layer**
  - FastAPI service running Prophet-based forecasts
  - Prescriptive recommendations (sales, inventory, finance)
  - Scenario simulation APIs for what-if analysis

- **Data Fabric & Governance**
  - Laravel migrations covering 1,000+ crore-scale datasets
  - ETL imports (Excel/CSV) for Sales, Production, Finance, Inventory, HR, Supply Chain, NPD
  - Demo data generator with configurable history (1–2 years)
  - Audit logging in `api_syncs`, `prescriptive_recommendations`, `scenario_simulations`

---

## 2. Executive Intelligence

| Capability                     | Details |
|--------------------------------|---------|
| Executive Overview Dashboard   | Cross-module KPIs, AI insights, alerts |
| Conversational Assistant       | Natural language queries (planned) |
| Scenario Planning              | Price/volume/cost impact simulations |
| AI Alerts                      | Consolidated anomaly feed per department |

---

## 3. Sales & Distribution

- Channel and product-level dashboards with KPI tiles, charts, and insights.
- Predictive sales summaries by channel/product (next month & quarter).
- Prescriptive recommendations (recovery plans, scaling guidance).
- Anomaly detection for unusual dips, revenue shocks.
- What-if panel illustrating incremental revenue & margin shifts.
- Order vs delivery monitoring, top distributors/retailers, cumulative trends.

---

## 4. Production & Manufacturing

- OEE analytics, efficiency, wastage, maintenance trackers.
- Production forecasts aligned with capacity constraints.
- Machine downtime, quality control, energy consumption dashboards.
- Integration with manufacturing KPIs via AI service (planned prescriptive add-ons).

---

## 5. Finance & Treasury

- Budget vs actual dashboards with variance analytics.
- Financial forecasts (cash flow, category budgets).
- Prescriptive guidance on expense control vs investment opportunities.
- Loan exposure analytics, net sales tracking, department & category breakouts.
- Scenario cards for margin optimization and surplus capital deployment.

---

## 6. Inventory & Supply Chain

- Stock valuation, turnover, GMROI, category/company breakdowns.
- Inventory forecasts with reorder point & safety stock indicators.
- Prescriptive order recommendations + slow mover identification.
- Supply chain module captures PO/GRN/invoice lifecycle, supplier performance, lead times.
- AI risk prediction for delayed orders (roadmap).

---

## 7. HR Analytics

- Workforce trends (headcount, attendance, promotions, departmental analytics).
- Attrition prediction with AI-assisted recommendations (roadmap to integrate prescriptive actions).
- Employee movement tracking, talent pipeline metrics.

---

## 8. Data Integration & Automation

- External API connectors configurable in Admin CMS with module-specific request/response templates, mapping engine, and auto-scheduling.
- Asynchronous sync status dashboard with retry management.
- Queue-based ingestion for large files and scheduled integrations.

---

## 9. Reporting & Exports

- Parameterized operational reports (Sales, IMS, Collections, etc.).
- Download formats: JSON/CSV, Excel via Laravel Excel.
- Role-based access with per-company visibility.

---

## 10. Security & Compliance

- Security headers, rate limiting, CORS policies applied at API layer.
- Audit logs for admin actions and AI recommendations.
- Configurable password policies, multi-factor support (roadmap).

---

## 11. UX Enhancements

- `useDashboardData` hook with caching, stale indicators, refresh controls.
- Responsive layout with active navigation cues.
- Skeleton loaders, Chip statuses, trend arrows for KPIs.
- Request tracing header (`x-request-id`) for cross-service observability.

This catalog should evolve alongside new features; update the relevant sections when modules gain additional functionality or move from roadmap to GA.
