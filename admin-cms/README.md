# Admin CMS (Laravel)

The Laravel admin portal orchestrates data ingestion, configuration, and AI enablement for iDash.

---

## ✨ Capabilities

- Guided file uploads (Excel/CSV/JSON) with validation and sample templates.
- External API configuration with dynamic request/response mapping, interval scheduling, and module documentation.
- Demo data generator (`idash:demo-data`) seeding 1–2 years of cross-module history sized for 1000 crore turnover.
- Admin dashboards for monitoring imports, integrations, and prescriptive artifacts.
- Blade-powered UI with RBAC (roles/permissions synced with platform).

---

## ⚙️ Setup

```bash
cd admin-cms
composer install
cp .env.example .env
php artisan key:generate

php artisan migrate
php artisan db:seed --class=RolePermissionSeeder
php artisan serve --host=0.0.0.0 --port=8001
```

The admin portal runs at `http://localhost:8001` by default.

---

## 🔁 Demo & ETL Commands

| Command | Purpose |
|---------|---------|
| `php artisan idash:demo-data --years=2 --refresh` | Reset module tables (excluding users) and load demo history |
| `php artisan idash:sync-external --job=sales` | Trigger external API sync for a specific job |
| `php artisan queue:work` | Process queued imports/integration tasks |

Refer to the [Sample Data Guide](../docs/SAMPLE_DATA_GUIDE.md) for customization tips.

---

## 📂 Structure

```
app/
├── Console/Commands/    # Demo data + integration commands
├── Http/Controllers/    # Admin modules (data, integrations, demo UI)
├── Models/              # Eloquent models (FileUpload, ExternalApi, ApiSync, ...)
resources/views/admin/   # Blade templates
routes/web.php           # Admin routes (auth + role middleware)
```

Migrations under `database/migrations/` define the canonical schema for the entire platform.

---

## 🧭 Development Workflow

1. Define routes in `routes/web.php` guarded by `auth` + role middleware.
2. Implement controller actions retrieving data via Eloquent/query builder.
3. Render Blade views using `resources/views/admin/...` layouts.
4. For heavy ETL, dispatch jobs/commands instead of blocking web requests.

More patterns documented in the [Admin CMS Developer Guide](../docs/developer/ADMIN_CMS_GUIDE.md).

---

## 🧪 Testing & Maintenance

- Run tests: `php artisan test`
- Clear caches: `php artisan optimize:clear`
- Monitor `jobs` / `failed_jobs` tables; restart workers with `php artisan queue:restart`.

---

## 📎 Resources

- [Setup Guide](../docs/SETUP_GUIDE.md)
- [Feature Catalog](../docs/FEATURE_CATALOG.md)
- [SRE Runbook](../docs/SRE_RUNBOOK.md)

