# Frontend (React + Vite)

The React application renders role-based dashboards backed by predictive and prescriptive analytics from the Go API and FastAPI services.

---

## ✨ Highlights

- React 18 + Vite with fast HMR and code splitting.
- Material UI design system with custom theming and skeleton loaders.
- Redux Toolkit for auth/session state; server data via `useDashboardData` hook (sessionStorage caching, stale refresh indicators, abortable requests).
- KPI widgets, AI insight feeds, anomaly lists, scenario impact cards.
- Axios client with request tracing headers and JWT propagation.

---

## 🔧 Setup

```bash
cd frontend
npm install
npm run dev -- --host
```

The dev server defaults to `http://localhost:5173`. Configure the backend URL via `.env`:

```
VITE_API_URL=http://localhost:8080
```

---

## 📂 Project Structure

```
src/
├── components/        # Reusable widgets (Dashboard, Sales, Finance, etc.)
├── hooks/             # Shared hooks (useDashboardData, etc.)
├── pages/             # Route-level pages (SalesDashboard, InventoryDashboard)
├── services/api.js    # Axios instance & interceptors
├── store/             # Redux slices (auth, layout)
├── theme/             # Material UI theme overrides
└── App.jsx            # Route definitions
```

---

## 🧪 Scripts

| Command           | Purpose |
|-------------------|---------|
| `npm run dev`     | Start Vite dev server |
| `npm run build`   | Production build (outputs `dist/`) |
| `npm run preview` | Preview production bundle |
| `npm test`        | Run unit tests (Vitest + RTL) |

Linting is handled via ESLint/Prettier (`npm run lint`).

---

## 🧭 Development Flow

1. Define new route in `App.jsx`.
2. Create dashboard/page under `src/pages` and fetch data with `useDashboardData('/module/endpoint')`.
3. Compose UI using existing dashboard components (KPIWidget, RecommendationList, ScenarioImpactCard).
4. Add module-specific widgets in `src/components/<Module>/` as needed.
5. Write tests for hooks/components with mocked API responses.

Refer to the [Frontend Developer Guide](../docs/developer/FRONTEND_GUIDE.md) for deeper workflow details, testing strategy, and release checklist.

---

## 🚀 Deployment

1. `npm run build`
2. Serve `dist/` via CDN, static hosting (S3 + CloudFront, Netlify, Vercel) or behind Nginx.
3. Ensure `VITE_API_URL` points to the production API domain.

---

## 📎 Resources

- [Setup Guide](../docs/SETUP_GUIDE.md)
- [Feature Catalog](../docs/FEATURE_CATALOG.md)
- [SRE Runbook](../docs/SRE_RUNBOOK.md)

