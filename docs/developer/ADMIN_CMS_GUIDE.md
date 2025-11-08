# Admin CMS Developer Guide

The Laravel-based Admin CMS powers data ingestion, configuration, and demo tooling. Use this reference to extend features safely.

---

## 1. Stack Overview

- **Framework**: Laravel 10 (PHP 8.1)
- **Views**: Blade templates under `resources/views`
- **Authentication**: Laravel Sanctum / session guards
- **Queues**: Database queue driver (Redis recommended for production)
- **ETL**: Artisan commands (`app/Console/Commands`) orchestrate imports & demo data generation

Directory cheat sheet:
```
admin-cms/
├── app/
│   ├── Console/Commands/      # Data generation & sync commands
│   ├── Http/Controllers/      # Admin modules (Sales, Finance, etc.)
│   ├── Models/                # Eloquent models
│   └── Providers/             # Route + Auth service providers
├── database/
│   ├── migrations/            # Canonical schema for platform
│   ├── seeders/               # RBAC + sample data seeds
├── resources/views/           # Blade templates for admin UI
├── routes/web.php             # Admin routes (web guard)
├── routes/auth.php            # Auth scaffolding
└── README.md
```

---

## 2. Request Flow (Route → View)

1. **Routing** (`routes/web.php`)
   - Group routes with `Route::middleware(['auth', 'role:admin'])` for admin-only pages.
   - Example: `Route::get('/admin/data/sales', [SalesController::class, 'index'])->name('admin.sales.index');`

2. **Controller** (`app/Http/Controllers/Admin/*Controller.php`)
   - Fetch data via Eloquent or query builders.
   - For long-running imports, dispatch queued jobs or trigger artisan commands.
   - Return view with data: `return view('admin.sales.index', compact('channels', 'summary'));`

3. **Blade View** (`resources/views/admin/...`)
   - Use the shared layout `layouts/admin.blade.php`.
   - Write lean templates; push heavy logic into view composers or controllers.

4. **Assets**
   - Use Laravel Mix/Vite if front-end compilation required. Keep admin UI lightweight.

---

## 3. Demo Data & ETL Workflows

- Primary command: `php artisan idash:demo-data --years=2 --company=1`
  - Resets module tables (sales, production, finance, inventory, HR, supply chain, NPD).
  - Preserves user accounts for login continuity.
  - Generates AI-ready datasets sized to 1000 crore turnover assumptions.

- Integration sync: `php artisan idash:sync-external --job=sales`
  - Reads credentials from `external_apis` table.
  - Updates `api_syncs` records for observability.

Ensure commands are idempotent and log progress using `info()` to surface in queue worker logs.

---

## 4. Coding Standards

- Follow PSR-12 formatting (`composer exec phpcs`).
- Utilize form requests for validation where applicable.
- Keep controllers thin; move business logic into service classes or jobs.
- Prefer queued jobs for heavy file processing (`Excel` imports) to keep web requests responsive.
- Use policies/gates for permission checks beyond role middleware.

---

## 5. Database & Migrations

- Migrations define the canonical schema for the entire platform; backend (Go) models mirror these tables.
- When introducing new tables/columns, add migration + seed data + update documentation (sample data guide).
- For large tables, add indexes in the same migration (see `2024_01_01_000015_add_optimized_indexes.php`).

---

## 6. Testing

| Type           | Command |
|----------------|---------|
| Feature tests  | `php artisan test` |
| Static analysis| `composer exec phpstan analyse` |

Use database transactions (`RefreshDatabase` trait) to reset state per test.

---

## 7. Queues & Scheduling

- Queue worker start: `php artisan queue:work --tries=3`
- Scheduler: configure cron `* * * * * php /path/to/artisan schedule:run`
- Monitor `jobs` and `failed_jobs` tables; SRE runbook covers alerting.

---

## 8. Deployment Checklist

1. `composer install --no-dev --optimize-autoloader`
2. `php artisan config:cache && php artisan route:cache`
3. Run database migrations.
4. Start queue workers + scheduler.
5. Verify admin dashboards and demo data console UI (under `admin/demo-data`).

---

## 9. Troubleshooting Tips

| Issue                               | Resolution |
|-------------------------------------|------------|
| Demo data command hangs             | Check queue workers; ensure `QUEUE_CONNECTION` configured |
| File upload fails                   | Verify disk permissions, increase PHP upload limits |
| Slow Blade pages                    | Optimize queries with eager loading (`with()`) |
| Permissions not applied             | Run `php artisan db:seed --class=RolePermissionSeeder` and clear caches |

Keep this guide updated when new admin modules, workflows, or integration endpoints are introduced.
