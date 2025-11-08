# Sample Data Guide

Use this guide to generate, customize, and reset demo data across all iDash modules.

---

## 1. Overview

The `idash:demo-data` artisan command seeds realistic historical records (sales, production, finance, inventory, HR, supply chain, NPD) sized for a 1000 crore turnover enterprise. User accounts remain untouched to preserve authentication flows.

Command location: `admin-cms/app/Console/Commands/GenerateDemoData.php`

---

## 2. Running the Generator

```bash
cd admin-cms
php artisan idash:demo-data --years=2 --company=1 --refresh
```

### Flags

| Flag        | Default | Description |
|-------------|---------|-------------|
| `--years`   | `1`     | Number of historical years to generate (1 or 2) |
| `--company` | `1`     | Target company ID (supports multi-company setups) |
| `--refresh` | `false` | When present, truncates module tables before seeding |
| `--channels`| `*`     | Optional comma-separated channel IDs to prioritize |

> Running with `--refresh` rebuilds dependent tables while excluding `users`, `roles`, `permissions`.

---

## 3. Data Coverage by Module

| Module        | Entities Populated                                                    |
|---------------|-----------------------------------------------------------------------|
| Sales         | Channelwise reports, best-selling products/PGs, orders, IMS, collections |
| Production    | Efficiency metrics, maintenance logs, OEE, production plans             |
| Finance       | Budgets, expenses, bank loan exposure, cash flow series                |
| Inventory     | Raw data values, COGS/GP summaries, GL accounts                        |
| Supply Chain  | Purchase orders, GRNs, invoices, supplier scorecards                   |
| HR            | Headcount, attendance, promotions, attrition probabilities             |
| NPD           | Projects, deliverables, milestones                                     |
| AI Forecasts  | Seeded entries in `ai_forecasts`, `prescriptive_recommendations`, `scenario_simulations` |

Generated data aligns with forecast expectations (e.g., monthly seasonality, growth rates) to showcase predictive and prescriptive insights.

---

## 4. Resetting the Environment

1. **Backup** existing data if needed.
2. Run `php artisan idash:demo-data --refresh --years=2`.
3. Clear caches: `php artisan optimize:clear`.
4. Verify dashboards by hitting `/sales/overview`, `/inventory/overview`, `/finance/overview`.

The command logs progress and summary statistics to console and `storage/logs/laravel.log`.

---

## 5. Customization Tips

- Modify distribution parameters (growth, seasonality) in the command to tailor industry scenarios.
- Adjust channel/product lists via seed data before running the generator for bespoke demos.
- Include additional `--segments` option (extend command) to target geography or business units.

---

## 6. Data Dictionary Resources

- Refer to `admin-cms/database/DATABASE_SCHEMA.md` for column definitions.
- AI service queries map to helper methods in `ai-service/database/database.py` (e.g., `get_sales_breakdown`).
- Backend models in `backend-api/models/*.go` mirror table structures and aggregate outputs.

Keep this guide current when adding new modules or extending the data generation command.
