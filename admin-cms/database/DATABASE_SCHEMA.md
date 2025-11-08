# Database Schema Documentation

This document describes the database schema for the Admin CMS system.

## Core Tables

### users
- **Purpose**: User authentication and authorization
- **Fields**:
  - `id` (primary key)
  - `name` (string)
  - `email` (string, unique)
  - `email_verified_at` (timestamp, nullable)
  - `password` (string)
  - `role` (string, default: 'admin')
  - `remember_token` (string)
  - `created_at`, `updated_at` (timestamps)

### file_uploads
- **Purpose**: Track file uploads for all modules
- **Fields**:
  - `id` (primary key)
  - `file_name` (string)
  - `file_path` (string)
  - `module` (string) - sales, production, finance, inventory, hr, supplychain, npd
  - `data_type` (string) - Type of data within the module
  - `uploaded_by` (foreign key -> users.id)
  - `status` (enum: pending, processing, completed, failed)
  - `records_processed` (integer)
  - `error_message` (text, nullable)
  - `created_at`, `updated_at` (timestamps)

### external_apis
- **Purpose**: Store external API configurations
- **Fields**:
  - `id` (primary key)
  - `name` (string)
  - `url` (string)
  - `method` (string, default: 'GET')
  - `headers` (text, JSON)
  - `body` (text, JSON)
  - `authentication` (text, JSON)
  - `module` (string)
  - `data_type` (string)
  - `data_mapping` (text, JSON)
  - `is_active` (boolean)
  - `sync_interval` (integer, minutes)
  - `last_sync_at` (timestamp, nullable)
  - `created_at`, `updated_at` (timestamps)

### api_syncs
- **Purpose**: Track API synchronization history
- **Fields**:
  - `id` (primary key)
  - `external_api_id` (foreign key -> external_apis.id)
  - `status` (enum: pending, processing, completed, failed)
  - `records_synced` (integer)
  - `response_data` (text, JSON)
  - `error_message` (text, nullable)
  - `started_at` (timestamp, nullable)
  - `completed_at` (timestamp, nullable)
  - `created_at`, `updated_at` (timestamps)

### companies
- **Purpose**: Company master data
- **Fields**:
  - `id` (primary key)
  - `name` (string)
  - `code` (string, unique)
  - `description` (text, nullable)
  - `is_active` (boolean)
  - `created_at`, `updated_at` (timestamps)

### channels
- **Purpose**: Sales channel master data
- **Fields**:
  - `id` (primary key)
  - `name` (string)
  - `status` (integer)
  - `created_at`, `updated_at` (timestamps)

## Sales Module Tables

### channelwise_monthly_report
- Monthly sales report by channel
- Key fields: `data_month`, `channel_id`, `lifting_target`, `billed`, `delivered`, `primary_collection`, `ims_target`, `ims`, etc.

### channelwise_lic_data
- Daily sales report by channel
- Key fields: `data_date`, `channel_id`, `billed`, `delivery`, `ims`, etc.

### best_selling_products
- Best selling products data
- Key fields: `year_month`, `channel_id`, `product_id`, `product_name`, `qty`, `value`

### best_selling_pgs
- Best selling product groups
- Key fields: `year_month`, `channel_id`, `category_id`, `category_name`, `qty`, `value`

### top_channel_d_bs
- Top distributors and retailers
- Key fields: `db_name`, `amount`, `type` (0=distributor, 1=retailer), `date`

### order_delivery_summaries
- Order vs delivery summary
- Key fields: `months`, `channel_id`, `amounts`, `types` (0=order, 1=delivery)

### top_retailers
- Top retailers data
- Key fields: `date`, `db_name`, `amount`

### sales_orders
- Sales orders data
- Key fields: `so_number`, `customer_code`, `document_date`, `product_code`, `product_name`, `so_qty`, `total_price`

### sales_deliveries
- Sales deliveries data
- Key fields: `sap_sales_order_no`, `sap_chalan_no`, `item_delivered_date`, `delivered_quantity`, `delivered_value`

## Production Module Tables

### production_analyses
- Production analysis data
- Key fields: `month`, `factory`, `summary_group`, `cmonthly_amount`, `pmonthly_amount`, `tmonthly_amount`, `amonthly_amount`

### wastage_datas
- Wastage data
- Key fields: `month`, `factory`, `group_name`, `wastage`, `month_wastage`, `avg`, `amount`

### cost_analyses
- Cost analysis data
- Key fields: `month`, `factory`, `cost_type`, `amount`

## Finance Module Tables

### bdgt_categories
- Budget categories
- Key fields: `name`, `status`

### bdgt_departments
- Budget departments (linked to categories)
- Key fields: `name`, `category_id`, `status`

### bdgt_expense_groups
- Budget expense groups (linked to departments)
- Key fields: `name`, `department_id`, `status`

### bdgt_sub_heads
- Budget sub heads (linked to expense groups)
- Key fields: `name`, `expense_group_id`, `status`

### bdgt_expenses
- Budget expenses (linked to sub heads)
- Key fields: `name`, `sub_head_id`

### budget_summaries
- Budget summaries
- Key fields: `month`, `category_id`, `department_id`, `budget_amount`, `actual_amount`

### budget_monthlies
- Monthly budget data
- Key fields: `month`, `expense_id`, `budget_amount`, `actual_amount`

### bank_loan_heads
- Bank loan heads
- Key fields: `loan_head`

### bank_loan_status_raw_data
- Bank loan status data
- Key fields: `month`, `loan_head`, `company_id`, `amount`

### financial_expense_raw_data
- Financial expense raw data
- Key fields: `month`, `expense_id`, `amount`

## Inventory Module Tables

### inventory_raw_datas
- Inventory raw data
- Key fields: `company_id`, `gl_id`, `month`, `amount`

### inventory_gl_accounts
- Inventory GL accounts
- Key fields: `gl_account`, `gl_account_name`, `description`

### cogs_gps
- COGS and GP data
- Key fields: `month`, `company_id`, `cogs`, `gp`, `gp_percentage`

### inventroy_sap_datas
- Inventory SAP data
- Key fields: `company`, `year`, `month`

## HR Module Tables

### employee_basic_infos
- Employee basic information
- Key fields: `total_active_staff`, `total_active_worker`, `total_contractual_employee`, `total_probationary_employee`, `total_permanent_employee`, `report_date`

### employee_attendances
- Employee attendance data
- Key fields: `date`, `total_absent`, `total_present`, `total_leave`

### employee_tran_overs
- Employee turnover data
- Key fields: `job_type`, `month`, `year`, `new_employee_no`, `resigned_employee`

### yearly_employee_promotions
- Yearly employee promotions
- Key fields: `year`, `promoted_count`, `details`

### hris_companies
- HRIS companies
- Key fields: `name`, `code`

### hris_departments
- HRIS departments (linked to companies)
- Key fields: `name`, `company_id`

### hris_promotion_break_downs
- HRIS promotion breakdowns (linked to departments)
- Key fields: `year`, `department_id`, `promoted_count`

## Supply Chain Module Tables

### supply_chain_raw_datas
- Supply chain raw data
- Key fields: `plant`, `pr_id`, `po_id`, `po_date`, `vendor_id`, `po_amount`, `grn1_id`, `invoice1_id`

### supply_chain_pos
- Supply chain purchase orders
- Key fields: `company_id`, `plant`, `vendor_code`, `material_code`, `po_number`, `po_date`, `po_amount`

### purchase_requisitions
- Purchase requisitions
- Key fields: `pr_id`, `pr_item`, `pr_date`, `material_code`, `quantity`, `plant`

## NPD Module Tables

### npd_projects
- NPD projects
- Key fields: `p_id`, `indent_no`, `name`, `pmo`, `project_manager`, `start_date`, `end_date`, `budget`, `status`, `progress`

### projects_deliverables
- Project deliverables (linked to projects)
- Key fields: `d_id`, `name`, `weightage`, `start_date`, `end_date`, `budget`, `progress`, `npd_project_id`

### projects_sub_deliverables
- Project sub deliverables (linked to deliverables)
- Key fields: `sd_id`, `name`, `weightage`, `start_date`, `end_date`, `budget`, `progress`, `deliverable_id`

## Indexes

Most tables have indexes on:
- Date fields (for time-based queries)
- Foreign key fields
- Commonly queried fields (e.g., `channel_id`, `company_id`)

## Relationships

- `file_uploads.uploaded_by` -> `users.id`
- `api_syncs.external_api_id` -> `external_apis.id`
- `bdgt_departments.category_id` -> `bdgt_categories.id`
- `bdgt_expense_groups.department_id` -> `bdgt_departments.id`
- `bdgt_sub_heads.expense_group_id` -> `bdgt_expense_groups.id`
- `bdgt_expenses.sub_head_id` -> `bdgt_sub_heads.id`
- `budget_summaries.category_id` -> `bdgt_categories.id`
- `budget_summaries.department_id` -> `bdgt_departments.id`
- `budget_monthlies.expense_id` -> `bdgt_expenses.id`
- `financial_expense_raw_data.expense_id` -> `bdgt_expenses.id`
- `supply_chain_pos.company_id` -> `companies.id`
- `projects_deliverables.npd_project_id` -> `npd_projects.id`
- `projects_sub_deliverables.deliverable_id` -> `projects_deliverables.id`
- `hris_departments.company_id` -> `hris_companies.id`
- `hris_promotion_break_downs.department_id` -> `hris_departments.id`

