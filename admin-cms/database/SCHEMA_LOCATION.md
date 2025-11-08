# Database Schema Location

## Schema Files Location

The database schema for the Admin CMS is located in the following places:

### 1. Migrations
**Location**: `database/migrations/`

All database migrations are stored here:
- `2024_01_01_000001_create_users_table.php` - Users table
- `2024_01_01_000002_create_file_uploads_table.php` - File uploads tracking
- `2024_01_01_000003_create_external_apis_table.php` - External API configurations
- `2024_01_01_000004_create_api_syncs_table.php` - API synchronization history
- `2024_01_01_000005_create_companies_table.php` - Companies master data
- `2024_01_01_000006_create_channels_table.php` - Sales channels
- `2024_01_01_000007_create_sales_tables.php` - All sales module tables
- `2024_01_01_000008_create_production_tables.php` - Production module tables
- `2024_01_01_000009_create_finance_tables.php` - Finance module tables
- `2024_01_01_000010_create_inventory_tables.php` - Inventory module tables
- `2024_01_01_000011_create_hr_tables.php` - HR module tables
- `2024_01_01_000012_create_supply_chain_tables.php` - Supply chain module tables
- `2024_01_01_000013_create_npd_tables.php` - NPD module tables

### 2. Schema Documentation
**Location**: `database/DATABASE_SCHEMA.md`

Complete documentation of all database tables, fields, relationships, and indexes.

### 3. Models
**Location**: `app/Models/`

Eloquent models that represent database tables:
- `User.php` - User model
- `FileUpload.php` - File upload model
- `ExternalApi.php` - External API model
- `ApiSync.php` - API sync model

### 4. Backend API Models (Go)
**Location**: `../backend-api/models/`

Go models that represent the same database structure:
- `user.go` - User model
- `sales.go` - Sales models
- `production.go` - Production models
- `finance.go` - Finance models
- `inventory.go` - Inventory models
- `hr.go` - HR models
- `supplychain.go` - Supply chain models
- `npd.go` - NPD models

## Running Migrations

To create the database schema:

```bash
cd admin-cms
php artisan migrate
```

To rollback migrations:

```bash
php artisan migrate:rollback
```

To see migration status:

```bash
php artisan migrate:status
```

## Database Configuration

The database configuration is in `.env` file:

```env
DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=idash_db
DB_USERNAME=your_username
DB_PASSWORD=your_password
```

## Viewing Schema

### Using Laravel Tinker
```bash
php artisan tinker
Schema::getTableListing();
Schema::getColumnListing('table_name');
```

### Using Database Client
Connect to MySQL database and use:
```sql
SHOW TABLES;
DESCRIBE table_name;
SHOW CREATE TABLE table_name;
```

### Using Data Viewer
Navigate to Admin CMS > Data Viewer to see all tables and their data.

## Schema Updates

When updating the schema:

1. Create a new migration:
   ```bash
   php artisan make:migration update_table_name
   ```

2. Edit the migration file in `database/migrations/`

3. Run the migration:
   ```bash
   php artisan migrate
   ```

4. Update the schema documentation in `database/DATABASE_SCHEMA.md`

5. Update models if needed in `app/Models/`

6. Update Go models in `../backend-api/models/` if needed

## Related Files

- **Backend API Schema**: See `../backend-api/models/` for Go models
- **Frontend Types**: See `../frontend/src/types/` for TypeScript types (if available)
- **API Documentation**: See main project README for API endpoints

