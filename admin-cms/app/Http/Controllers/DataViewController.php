<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

class DataViewController extends Controller
{
    public function index()
    {
        $tables = [
            // Core Tables
            'users' => 'Users',
            'companies' => 'Companies',
            'channels' => 'Channels',
            'file_uploads' => 'File Uploads',
            'external_apis' => 'External APIs',
            'api_syncs' => 'API Syncs',
            
            // Sales Module
            'channelwise_monthly_report' => 'Sales - Monthly Report',
            'channelwise_lic_data' => 'Sales - Daily Report',
            'best_selling_products' => 'Sales - Best Selling Products',
            'best_selling_pgs' => 'Sales - Best Selling PGs',
            'top_channel_d_bs' => 'Sales - Top Distributors/Retailers',
            'order_delivery_summaries' => 'Sales - Order/Delivery Summary',
            'top_retailers' => 'Sales - Top Retailers',
            'sales_orders' => 'Sales - Orders',
            'sales_deliveries' => 'Sales - Deliveries',
            
            // Production Module
            'production_analyses' => 'Production - Analysis',
            'wastage_datas' => 'Production - Wastage Data',
            'cost_analyses' => 'Production - Cost Analysis',
            
            // Finance Module
            'bdgt_categories' => 'Finance - Budget Categories',
            'bdgt_departments' => 'Finance - Budget Departments',
            'bdgt_expense_groups' => 'Finance - Budget Expense Groups',
            'bdgt_sub_heads' => 'Finance - Budget Sub Heads',
            'bdgt_expenses' => 'Finance - Budget Expenses',
            'budget_summaries' => 'Finance - Budget Summaries',
            'budget_monthlies' => 'Finance - Budget Monthlies',
            'bank_loan_heads' => 'Finance - Bank Loan Heads',
            'bank_loan_status_raw_data' => 'Finance - Bank Loan Status',
            'financial_expense_raw_data' => 'Finance - Expense Raw Data',
            
            // Inventory Module
            'inventory_raw_datas' => 'Inventory - Raw Data',
            'inventory_gl_accounts' => 'Inventory - GL Accounts',
            'cogs_gps' => 'Inventory - COGS & GP',
            'inventroy_sap_datas' => 'Inventory - SAP Data',
            
            // HR Module
            'employee_basic_infos' => 'HR - Employee Basic Info',
            'employee_attendances' => 'HR - Employee Attendance',
            'employee_tran_overs' => 'HR - Employee Turnover',
            'yearly_employee_promotions' => 'HR - Yearly Promotions',
            'hris_companies' => 'HR - HRIS Companies',
            'hris_departments' => 'HR - HRIS Departments',
            'hris_promotion_break_downs' => 'HR - Promotion Breakdowns',
            
            // Supply Chain Module
            'supply_chain_raw_datas' => 'Supply Chain - Raw Data',
            'supply_chain_pos' => 'Supply Chain - Purchase Orders',
            'purchase_requisitions' => 'Supply Chain - Purchase Requisitions',
            
            // NPD Module
            'npd_projects' => 'NPD - Projects',
            'projects_deliverables' => 'NPD - Project Deliverables',
            'projects_sub_deliverables' => 'NPD - Project Sub Deliverables',
        ];
        
        return view('admin.data.viewer', compact('tables'));
    }
    
    public function show($table)
    {
        $allowedTables = [
            'users', 'companies', 'channels', 'file_uploads', 'external_apis', 'api_syncs',
            'channelwise_monthly_report', 'channelwise_lic_data', 'best_selling_products',
            'best_selling_pgs', 'top_channel_d_bs', 'order_delivery_summaries',
            'top_retailers', 'sales_orders', 'sales_deliveries',
            'production_analyses', 'wastage_datas', 'cost_analyses',
            'bdgt_categories', 'bdgt_departments', 'bdgt_expense_groups',
            'bdgt_sub_heads', 'bdgt_expenses', 'budget_summaries',
            'budget_monthlies', 'bank_loan_heads', 'bank_loan_status_raw_data',
            'financial_expense_raw_data',
            'inventory_raw_datas', 'inventory_gl_accounts', 'cogs_gps',
            'inventroy_sap_datas',
            'employee_basic_infos', 'employee_attendances', 'employee_tran_overs',
            'yearly_employee_promotions', 'hris_companies', 'hris_departments',
            'hris_promotion_break_downs',
            'supply_chain_raw_datas', 'supply_chain_pos', 'purchase_requisitions',
            'npd_projects', 'projects_deliverables', 'projects_sub_deliverables',
        ];
        
        if (!in_array($table, $allowedTables)) {
            return redirect()->route('admin.data.viewer')
                ->with('error', 'Table not found or inaccessible');
        }
        
        try {
            if (!Schema::hasTable($table)) {
                return redirect()->route('admin.data.viewer')
                    ->with('error', 'Table does not exist');
            }
            
            $data = DB::table($table)->orderBy('id', 'desc')->paginate(50);
            $columns = Schema::getColumnListing($table);
            $tableLabel = $this->getTableLabel($table);
            
            return view('admin.data.view', compact('table', 'data', 'columns', 'tableLabel'));
        } catch (\Exception $e) {
            return redirect()->route('admin.data.viewer')
                ->with('error', 'Error loading data: ' . $e->getMessage());
        }
    }
    
    private function getTableLabel($table)
    {
        $labels = [
            'users' => 'Users',
            'companies' => 'Companies',
            'channels' => 'Channels',
            'file_uploads' => 'File Uploads',
            'external_apis' => 'External APIs',
            'api_syncs' => 'API Syncs',
            'channelwise_monthly_report' => 'Sales - Monthly Report',
            'channelwise_lic_data' => 'Sales - Daily Report',
            'best_selling_products' => 'Sales - Best Selling Products',
            'best_selling_pgs' => 'Sales - Best Selling PGs',
            'top_channel_d_bs' => 'Sales - Top Distributors/Retailers',
            'order_delivery_summaries' => 'Sales - Order/Delivery Summary',
            'top_retailers' => 'Sales - Top Retailers',
            'sales_orders' => 'Sales - Orders',
            'sales_deliveries' => 'Sales - Deliveries',
            'production_analyses' => 'Production - Analysis',
            'wastage_datas' => 'Production - Wastage Data',
            'cost_analyses' => 'Production - Cost Analysis',
            'bdgt_categories' => 'Finance - Budget Categories',
            'bdgt_departments' => 'Finance - Budget Departments',
            'bdgt_expense_groups' => 'Finance - Budget Expense Groups',
            'bdgt_sub_heads' => 'Finance - Budget Sub Heads',
            'bdgt_expenses' => 'Finance - Budget Expenses',
            'budget_summaries' => 'Finance - Budget Summaries',
            'budget_monthlies' => 'Finance - Budget Monthlies',
            'bank_loan_heads' => 'Finance - Bank Loan Heads',
            'bank_loan_status_raw_data' => 'Finance - Bank Loan Status',
            'financial_expense_raw_data' => 'Finance - Expense Raw Data',
            'inventory_raw_datas' => 'Inventory - Raw Data',
            'inventory_gl_accounts' => 'Inventory - GL Accounts',
            'cogs_gps' => 'Inventory - COGS & GP',
            'inventroy_sap_datas' => 'Inventory - SAP Data',
            'employee_basic_infos' => 'HR - Employee Basic Info',
            'employee_attendances' => 'HR - Employee Attendance',
            'employee_tran_overs' => 'HR - Employee Turnover',
            'yearly_employee_promotions' => 'HR - Yearly Promotions',
            'hris_companies' => 'HR - HRIS Companies',
            'hris_departments' => 'HR - HRIS Departments',
            'hris_promotion_break_downs' => 'HR - Promotion Breakdowns',
            'supply_chain_raw_datas' => 'Supply Chain - Raw Data',
            'supply_chain_pos' => 'Supply Chain - Purchase Orders',
            'purchase_requisitions' => 'Supply Chain - Purchase Requisitions',
            'npd_projects' => 'NPD - Projects',
            'projects_deliverables' => 'NPD - Project Deliverables',
            'projects_sub_deliverables' => 'NPD - Project Sub Deliverables',
        ];
        
        return $labels[$table] ?? $table;
    }
}
