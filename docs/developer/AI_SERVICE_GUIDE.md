# AI Service Developer Guide

The FastAPI-based AI service powers predictive and prescriptive analytics for iDash. This document explains the data flow, extensibility points, and operational considerations.

---

## 1. Technology Stack

- **Framework**: FastAPI (ASGI)
- **Modeling Libraries**: Prophet (time series), scikit-learn (Isolation Forest), NumPy, pandas
- **Data Access**: Custom MySQL helper (`database/database.py`)
- **Prescriptive Engine**: Heuristic recommendations in `services/prescriptive_service.py`
- **Runtime**: Uvicorn/Gunicorn

Directory layout:
```
ai-service/
├── main.py                     # FastAPI app + routes
├── database/database.py        # MySQL access helpers
├── services/
│   ├── forecast_service.py     # Forecast generation (sales, production, finance, inventory)
│   ├── analysis_service.py     # Trend, anomaly, correlation analysis
│   └── prescriptive_service.py # Recommendation & scenario logic
├── requirements.txt            # Python dependencies
└── README.md
```

---

## 2. Endpoint Lifecycle (Route → Response)

1. **Request models** (Pydantic)
   - `ForecastRequest`, `PredictRequest`, `PrescriptiveRequest`, `ScenarioRequest`, `AnalysisRequest` define validated payloads.

2. **Route handlers** (`main.py`)
   - `/api/forecast/<type>`: fetch historical data, run forecasts, save metadata.
   - `/api/predict/sales/summary`: summarise grouped sales + recommend actions.
   - `/api/prescribe/<module>`: combine forecasts with heuristics for inventory/finance.
   - `/api/analyze/anomaly/enriched`: anomaly detection with severity + recommended actions.
   - `/api/scenario/whatif`: simulate price/volume/cost adjustments.

3. **Service logic** (`services/*.py`)
   - `ForecastService`: wraps Prophet pipeline, calculates confidence, reorder points, etc.
   - `AnalysisService`: IsolationForest anomalies, trend slope calculations.
   - `PrescriptiveService`: heuristics for sales, inventory, finance and scenario modelling.

4. **Database helper** (`database/database.py`)
   - Provides scoped queries: sales breakdown, inventory history, financial summaries.
   - Manages connection reuse; ensure credentials set via env vars.

---

## 3. Adding a New Forecast Type

1. Extend `ForecastService` with `generate_<domain>_forecast` method.
2. Add corresponding database fetch helper (e.g., `get_<domain>_data`).
3. Create FastAPI route in `main.py` mirroring existing endpoints.
4. Update backend API to proxy and persist results if required.
5. Document new route in README and relevant developer guide.

---

## 4. Coding Standards

- Format code with `black` (`black .`).
- Lint using `flake8` or `ruff`.
- Keep long-running operations async-friendly; heavy CPU tasks may require background workers or Celery if concurrency increases.
- Guard against empty datasets; raise HTTP 404 with descriptive message.
- Return ISO timestamps via `datetime.now().isoformat()` for consistency.

---

## 5. Testing

| Type           | Tool            | Command |
|----------------|-----------------|---------|
| Unit           | pytest          | `pytest` |
| Type checking  | mypy (optional) | `mypy services/` |

Mock database queries using `monkeypatch` or dependency injection. For Prophet-heavy tests, seed deterministic data to avoid flaky outputs.

---

## 6. Performance Considerations

- Prophet models can be expensive; reuse trained models when forecasting similar series (future enhancement).
- Limit forecast horizon to practical window (default 30 days) unless business requires longer.
- Cache expensive grouped results using Redis or in-process caching if request volume spikes.
- Use async DB driver (`aiomysql`) in future iterations if concurrency becomes bottleneck.

---

## 7. Deployment Checklist

1. Build container image with slim Python base.
2. Install system dependencies (gcc, prophet requirements).
3. Run database migrations indirectly via Laravel if schema changes.
4. Configure environment variables (`DB_HOST`, `DB_USER`, `AI_LOG_LEVEL`).
5. Start using Gunicorn + Uvicorn workers (`gunicorn -k uvicorn.workers.UvicornWorker main:app -w 2 --threads 4`).
6. Attach health probe to `/health`.

---

## 8. Troubleshooting

| Issue                               | Resolution |
|-------------------------------------|------------|
| `Sales forecast generation failed`  | Verify historical data availability, inspect Prophet warnings |
| Connection errors                   | Check DB credentials, network security groups |
| High latency                        | Scale worker count, profile heavy endpoints, review MySQL indexes |
| Inconsistent recommendations        | Log intermediate values, tweak heuristics in `prescriptive_service.py` |

Keep this guide in sync with future additions such as MLflow integration, model versioning, or GPU acceleration.
