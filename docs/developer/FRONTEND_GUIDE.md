# Frontend Developer Guide

This guide describes the development workflow for the React (Vite) frontend, from routing to rendering, with conventions for contributing features safely.

---

## 1. Project Overview

- **Framework**: React 18 + Vite
- **Routing**: `react-router-dom`
- **State**: Redux Toolkit (`src/store`)
- **UI**: Material UI v5, custom components under `src/components`
- **Data Access**: Axios client (`src/services/api.js`) with sessionStorage caching via `useDashboardData`

Directory highlights:
```
src/
├── App.jsx                # Route definitions
├── components/            # Reusable UI blocks
├── hooks/                 # Shared hooks (data fetching, utilities)
├── pages/                 # Route-level screens
├── services/api.js        # Axios instance & interceptors
├── store/                 # Redux slices & store setup
└── theme/                 # MUI theme overrides
```

---

## 2. Request → View Workflow

1. **Route registration** (`src/App.jsx`)
   - Add a `<Route path="/new" element={<NewPage />} />` entry.
   - Protect routes with `<PrivateRoute roles={["executive"]} />` if necessary.

2. **Page component** (`src/pages/NewPage.jsx`)
   - Import `useDashboardData('/api/endpoint')` to fetch backend data.
   - Compose UI using layout components (`Container`, `Grid`, `Stack`).
   - Handle loading (`Skeleton`), error (`Alert`), and stale states using the hook outputs.

3. **Hooks & Data**
   - `useDashboardData` caches GET requests in `sessionStorage`, refreshes stale data, and exposes `refresh/stale/loading` flags.
   - For bespoke logic, create dedicated hooks in `src/hooks/`.

4. **Components**
   - Place new widgets under relevant module folder (e.g., `components/Sales/`).
   - Favor composable dashboard building blocks (e.g., `KPIWidget`, `RecommendationList`, `ScenarioImpactCard`).

5. **Styling**
   - Adhere to MUI system props, avoid inline styles where possible.
   - Extend theme tokens in `src/theme/theme.js` for shared colors/spacing.

---

## 3. State Management Rules

- Use Redux slices **only** for cross-page state (auth, user profile, feature flags).
- Local component state (`useState`, `useReducer`) is preferred for UI concerns.
- Async logic should stay in hooks; avoid dispatching thunks from inside components when `useDashboardData` covers the use case.

---

## 4. API Conventions

- All calls should use the shared Axios instance (`import api from '../services/api'`).
- Backend routes live under `/api/...`; pass query params via `{ params: { key: value } }`.
- Ensure endpoints exist or request addition in backend developer guide.
- Include `x-request-id` header (already added in interceptor) for traceability.

---

## 5. Coding Standards

- Use functional components with hooks; avoid legacy class components.
- Favor TypeScript-like JSDoc comments for complex props/hook return types.
- Apply ESLint & Prettier (run `npm run lint` if configured).
- Keep components focused: prefer smaller composable pieces over monolithic files.
- Place mock/sample data under `src/test/` for storybook/unit tests.

---

## 6. Testing Strategy

| Layer           | Tooling                | Command            |
|-----------------|------------------------|--------------------|
| Unit            | Vitest + React Testing Library | `npm test` |
| Integration     | Cypress (optional)     | `npm run cypress`  |
| Accessibility   | `@axe-core/react` (manual) | run in dev tools |

- For dashboards, create snapshot tests around `useDashboardData` mocks to ensure layout regressions are caught.

---

## 7. Performance Tips

- Memoize expensive derived values with `useMemo`.
- Lazy-load heavy charts via `React.lazy` + `Suspense` when necessary.
- Use `Skeleton` placeholders to keep UX responsive.
- Respect cache TTLs and avoid forced refresh unless user requests.

---

## 8. Release Checklist

1. `npm install` (ensure lockfile updated if deps change).
2. `npm run lint && npm run build` must succeed.
3. Update changelog / release notes with user-facing improvements.
4. Verify KPIs/charts on staging, including stale indicators and refresh button behaviour.
5. Tag commit and push; CI pipeline will publish static assets.

---

## 9. Useful References

- [Material UI docs](https://mui.com/material-ui/getting-started/installation/)
- [Redux Toolkit patterns](https://redux-toolkit.js.org/tutorials/quick-start)
- [Axios cancellation](https://axios-http.com/docs/cancellation) (already handled in `useDashboardData`).

Keep this guide updated when adding new architectural patterns (e.g., micro-frontends, module federation).
