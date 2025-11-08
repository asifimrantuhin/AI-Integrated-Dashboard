<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Carbon\Carbon;

class RolePermissionSeeder extends Seeder
{
    public function run(): void
    {
        $timestamp = Carbon::now();

        $roles = [
            ['name' => 'executive', 'display_name' => 'Executive Admin', 'description' => 'Full access to all modules'],
            ['name' => 'finance_manager', 'display_name' => 'Finance Manager', 'description' => 'Access to finance and executive dashboards'],
            ['name' => 'sales_manager', 'display_name' => 'Sales Manager', 'description' => 'Access to sales dashboards and analytics'],
            ['name' => 'hr_manager', 'display_name' => 'HR Manager', 'description' => 'Access to HR analytics and workforce insights'],
            ['name' => 'production_manager', 'display_name' => 'Production Manager', 'description' => 'Access to production and factory dashboards'],
            ['name' => 'supply_chain_manager', 'display_name' => 'Supply Chain Manager', 'description' => 'Access to supply chain and procurement dashboards'],
            ['name' => 'inventory_manager', 'display_name' => 'Inventory Manager', 'description' => 'Access to inventory analytics and forecasting'],
            ['name' => 'analyst', 'display_name' => 'Business Analyst', 'description' => 'Access to BI and reporting tools'],
        ];

        foreach ($roles as $role) {
            DB::table('roles')->updateOrInsert(
                ['name' => $role['name']],
                [
                    'display_name' => $role['display_name'],
                    'description' => $role['description'],
                    'is_system' => true,
                    'created_at' => $timestamp,
                    'updated_at' => $timestamp,
                ]
            );
        }

        $permissions = [
            ['name' => 'view_dashboard', 'display_name' => 'View Dashboard', 'module' => 'core'],
            ['name' => 'manage_users', 'display_name' => 'Manage Users', 'module' => 'core'],
            ['name' => 'manage_roles', 'display_name' => 'Manage Roles', 'module' => 'core'],
            ['name' => 'view_sales', 'display_name' => 'View Sales Analytics', 'module' => 'sales'],
            ['name' => 'manage_sales', 'display_name' => 'Manage Sales Data', 'module' => 'sales'],
            ['name' => 'view_production', 'display_name' => 'View Production Analytics', 'module' => 'production'],
            ['name' => 'manage_production', 'display_name' => 'Manage Production Data', 'module' => 'production'],
            ['name' => 'view_finance', 'display_name' => 'View Finance Analytics', 'module' => 'finance'],
            ['name' => 'manage_finance', 'display_name' => 'Manage Finance Data', 'module' => 'finance'],
            ['name' => 'view_inventory', 'display_name' => 'View Inventory Analytics', 'module' => 'inventory'],
            ['name' => 'manage_inventory', 'display_name' => 'Manage Inventory Data', 'module' => 'inventory'],
            ['name' => 'view_hr', 'display_name' => 'View HR Analytics', 'module' => 'hr'],
            ['name' => 'manage_hr', 'display_name' => 'Manage HR Data', 'module' => 'hr'],
            ['name' => 'view_supply_chain', 'display_name' => 'View Supply Chain Analytics', 'module' => 'supplychain'],
            ['name' => 'manage_supply_chain', 'display_name' => 'Manage Supply Chain Data', 'module' => 'supplychain'],
            ['name' => 'view_bi', 'display_name' => 'View BI Insights', 'module' => 'bi'],
            ['name' => 'manage_reports', 'display_name' => 'Manage Reports', 'module' => 'reporting'],
        ];

        foreach ($permissions as $permission) {
            DB::table('permissions')->updateOrInsert(
                ['name' => $permission['name']],
                [
                    'display_name' => $permission['display_name'],
                    'module' => $permission['module'],
                    'description' => $permission['display_name'],
                    'created_at' => $timestamp,
                    'updated_at' => $timestamp,
                ]
            );
        }

        $rolePermissions = [
            'executive' => ['view_dashboard', 'manage_users', 'manage_roles', 'view_sales', 'manage_sales', 'view_production', 'manage_production', 'view_finance', 'manage_finance', 'view_inventory', 'manage_inventory', 'view_hr', 'manage_hr', 'view_supply_chain', 'manage_supply_chain', 'view_bi', 'manage_reports'],
            'finance_manager' => ['view_dashboard', 'view_finance', 'manage_finance', 'view_sales', 'view_inventory', 'view_bi', 'manage_reports'],
            'sales_manager' => ['view_dashboard', 'view_sales', 'manage_sales', 'view_finance', 'view_inventory', 'view_bi'],
            'hr_manager' => ['view_dashboard', 'view_hr', 'manage_hr', 'view_finance', 'view_bi'],
            'production_manager' => ['view_dashboard', 'view_production', 'manage_production', 'view_inventory', 'view_supply_chain', 'view_bi'],
            'supply_chain_manager' => ['view_dashboard', 'view_supply_chain', 'manage_supply_chain', 'view_inventory', 'view_bi'],
            'inventory_manager' => ['view_dashboard', 'view_inventory', 'manage_inventory', 'view_supply_chain', 'view_bi'],
            'analyst' => ['view_dashboard', 'view_sales', 'view_production', 'view_finance', 'view_inventory', 'view_hr', 'view_supply_chain', 'view_bi', 'manage_reports'],
        ];

        $roleIdMap = DB::table('roles')->pluck('id', 'name');
        $permissionIdMap = DB::table('permissions')->pluck('id', 'name');

        foreach ($rolePermissions as $roleName => $permissionList) {
            $roleId = $roleIdMap[$roleName] ?? null;
            if (!$roleId) {
                continue;
            }

            foreach ($permissionList as $permissionName) {
                $permissionId = $permissionIdMap[$permissionName] ?? null;
                if (!$permissionId) {
                    continue;
                }

                DB::table('role_permissions')->updateOrInsert(
                    ['role_id' => $roleId, 'permission_id' => $permissionId],
                    ['created_at' => $timestamp, 'updated_at' => $timestamp]
                );
            }
        }
    }
}
