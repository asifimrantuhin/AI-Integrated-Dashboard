<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;
use Illuminate\Support\Facades\DB;

return new class extends Migration
{
    public function up()
    {
        // Add composite indexes for faster queries
        // Sales tables
        DB::statement('ALTER TABLE channelwise_monthly_report ADD INDEX idx_month_channel (data_month, channel_id)');
        DB::statement('ALTER TABLE channelwise_monthly_report ADD INDEX idx_month_billed (data_month, billed)');
        DB::statement('ALTER TABLE channelwise_lic_data ADD INDEX idx_date_channel (data_date, channel_id)');
        DB::statement('ALTER TABLE best_selling_products ADD INDEX idx_month_product (year_month, product_id)');
        DB::statement('ALTER TABLE sales_orders ADD INDEX idx_date_customer (document_date, customer_code)');
        DB::statement('ALTER TABLE sales_orders ADD INDEX idx_so_product (so_number, product_code)');
        
        // Production tables
        DB::statement('ALTER TABLE production_analyses ADD INDEX idx_month_factory (month, factory)');
        DB::statement('ALTER TABLE wastage_datas ADD INDEX idx_month_factory (month, factory)');
        DB::statement('ALTER TABLE cost_analyses ADD INDEX idx_month_factory_type (month, factory, cost_type)');
        
        // Finance tables
        DB::statement('ALTER TABLE budget_summaries ADD INDEX idx_month_category (month, category_id)');
        DB::statement('ALTER TABLE budget_monthlies ADD INDEX idx_month_expense (month, expense_id)');
        DB::statement('ALTER TABLE bank_loan_status_raw_data ADD INDEX idx_month_company (month, company_id)');
        DB::statement('ALTER TABLE financial_expense_raw_data ADD INDEX idx_month_expense (month, expense_id)');
        
        // Inventory tables
        DB::statement('ALTER TABLE inventory_raw_datas ADD INDEX idx_month_company (month, company_id)');
        DB::statement('ALTER TABLE cogs_gps ADD INDEX idx_month_company (month, company_id)');
        
        // HR tables
        DB::statement('ALTER TABLE employee_attendances ADD INDEX idx_date_present (date, total_present)');
        DB::statement('ALTER TABLE employee_tran_overs ADD INDEX idx_year_month (year, month)');
        
        // Supply Chain tables
        DB::statement('ALTER TABLE supply_chain_raw_datas ADD INDEX idx_po_date (po_date, plant)');
        DB::statement('ALTER TABLE supply_chain_pos ADD INDEX idx_po_date_company (po_date, company_id)');
        
        // NPD tables
        DB::statement('ALTER TABLE npd_projects ADD INDEX idx_status_date (status, start_date)');
        DB::statement('ALTER TABLE projects_deliverables ADD INDEX idx_project_status (npd_project_id, progress)');
    }

    public function down()
    {
        // Note: Dropping indexes in MySQL requires ALTER TABLE DROP INDEX
        // This is a simplified version - in production, you'd need to drop each index individually
        DB::statement('ALTER TABLE channelwise_monthly_report DROP INDEX idx_month_channel');
        DB::statement('ALTER TABLE channelwise_monthly_report DROP INDEX idx_month_billed');
        DB::statement('ALTER TABLE channelwise_lic_data DROP INDEX idx_date_channel');
        DB::statement('ALTER TABLE best_selling_products DROP INDEX idx_month_product');
        DB::statement('ALTER TABLE sales_orders DROP INDEX idx_date_customer');
        DB::statement('ALTER TABLE sales_orders DROP INDEX idx_so_product');
        DB::statement('ALTER TABLE production_analyses DROP INDEX idx_month_factory');
        DB::statement('ALTER TABLE wastage_datas DROP INDEX idx_month_factory');
        DB::statement('ALTER TABLE cost_analyses DROP INDEX idx_month_factory_type');
        DB::statement('ALTER TABLE budget_summaries DROP INDEX idx_month_category');
        DB::statement('ALTER TABLE budget_monthlies DROP INDEX idx_month_expense');
        DB::statement('ALTER TABLE bank_loan_status_raw_data DROP INDEX idx_month_company');
        DB::statement('ALTER TABLE financial_expense_raw_data DROP INDEX idx_month_expense');
        DB::statement('ALTER TABLE inventory_raw_datas DROP INDEX idx_month_company');
        DB::statement('ALTER TABLE cogs_gps DROP INDEX idx_month_company');
        DB::statement('ALTER TABLE employee_attendances DROP INDEX idx_date_present');
        DB::statement('ALTER TABLE employee_tran_overs DROP INDEX idx_year_month');
        DB::statement('ALTER TABLE supply_chain_raw_datas DROP INDEX idx_po_date');
        DB::statement('ALTER TABLE supply_chain_pos DROP INDEX idx_po_date_company');
        DB::statement('ALTER TABLE npd_projects DROP INDEX idx_status_date');
        DB::statement('ALTER TABLE projects_deliverables DROP INDEX idx_project_status');
    }
};

