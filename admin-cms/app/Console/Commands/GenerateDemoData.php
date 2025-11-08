<?php

namespace App\Console\Commands;

use Carbon\Carbon;
use Carbon\CarbonPeriod;
use Faker\Factory as Faker;
use Illuminate\Console\Command;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Schema;

class GenerateDemoData extends Command
{
    protected $signature = 'idash:generate-demo-data {--years=2 : Number of years to generate (1-2)}';

    protected $description = 'Reset module data (except users) and generate demo data for the last N years.';

    public function handle(): int
    {
        $years = (int) $this->option('years');
        if ($years < 1 || $years > 5) {
            $this->warn('Years option must be 1 to 5. Defaulting to 1.');
            $years = 1;
        }

        $faker = Faker::create();
        $this->info("Generating {$years} year(s) of demo data...");

        $this->resetTables();

        $startMonth = Carbon::now()->subYears($years)->startOfMonth();
        $endMonth = Carbon::now()->subMonth()->endOfMonth();
        $monthSeries = $this->createMonthSeries($startMonth, $endMonth);
        $daySeries = $this->createDaySeries(Carbon::now()->subMonths(6)->startOfDay(), Carbon::now()->startOfDay());

        DB::transaction(function () use ($faker, $years, $monthSeries, $daySeries) {
            $companies = $this->seedCompanies();
            $channels = $this->seedChannels();
            $glAccounts = $this->seedInventoryAccounts();
            $hrDepartments = $this->seedDepartments($companies);
            $financeStructure = $this->seedFinanceStructure();
            $loanHeads = $this->seedBankLoanHeads();
            $suppliers = $this->seedSuppliers($companies);

            $this->seedSalesData($companies, $channels, $monthSeries, $daySeries, $faker, $years);
            $this->seedProductionData($companies, $monthSeries, $daySeries, $faker);
            $this->seedFinanceData(
                $companies,
                $financeStructure['categories'],
                $financeStructure['departments'],
                $financeStructure['expenses'],
                $loanHeads,
                $monthSeries
            );
            $this->seedInventoryData($companies, $glAccounts, $monthSeries, $faker);
            $this->seedHRData($companies, $hrDepartments, $monthSeries, $daySeries, $faker);
            $this->seedSupplyChainData($companies, $monthSeries, $faker);
            $this->seedManufacturingData($companies, $suppliers, $monthSeries, $daySeries, $faker);
            $this->seedNPDData($companies, $faker);
            $this->seedForecasts($companies, $monthSeries, $faker);
        });

        $this->info('Demo data generation complete.');

        return Command::SUCCESS;
    }

    private function resetTables(): void
    {
        $tables = [
            // Sales
            'channelwise_monthly_report', 'channelwise_lic_data', 'best_selling_products', 'best_selling_pgs',
            'top_channel_d_bs', 'order_delivery_summaries', 'top_retailers', 'sales_orders', 'sales_deliveries',
            // Production & Manufacturing
            'production_analyses', 'wastage_datas', 'cost_analyses', 'production_efficiency', 'machine_maintenances',
            'quality_control_checks', 'production_plans', 'material_requirements', 'energy_consumption', 'supplier_performance',
            // Finance
            'bank_loan_heads', 'bank_loan_status_raw_data', 'financial_expense_raw_data',
            'budget_summaries', 'budget_monthlies',
            'bdgt_expenses', 'bdgt_sub_heads', 'bdgt_expense_groups', 'bdgt_departments', 'bdgt_categories',
            // Inventory
            'inventory_raw_datas', 'inventory_gl_accounts', 'cogs_gps', 'inventroy_sap_datas',
            // HR
            'employee_basic_infos', 'employee_attendances', 'employee_tran_overs', 'yearly_employee_promotions',
            'hris_companies', 'hris_departments', 'hris_promotion_break_downs',
            // Supply Chain
            'supply_chain_master_datas', 'supply_chain_grn_datas', 'supply_chain_invoice_datas', 'supply_chain_pos', 'supply_chain_raw_datas', 'purchase_requisitions',
            // NPD (if present)
            'npd_projects', 'projects_deliverables', 'projects_sub_deliverables',
            // Forecast/cache/report config
            'ai_forecasts', 'api_cache',
            // Masters
            'companies', 'channels'
        ];

        $this->info('Resetting existing data (excluding users)...');
        DB::statement('SET FOREIGN_KEY_CHECKS=0');
        foreach ($tables as $table) {
            if (Schema::hasTable($table)) {
                DB::table($table)->truncate();
            }
        }
        DB::statement('SET FOREIGN_KEY_CHECKS=1');
    }

    private function createMonthSeries(Carbon $start, Carbon $end): array
    {
        $period = CarbonPeriod::create($start, '1 month', $end);
        return collect($period)->map(fn ($date) => $date->copy())->values()->all();
    }

    private function createDaySeries(Carbon $start, Carbon $end): array
    {
        $period = CarbonPeriod::create($start, '1 day', $end);
        return collect($period)->map(fn ($date) => $date->copy())->values()->all();
    }

    private function seedCompanies(): array
    {
        $companies = [
            ['code' => 'ID-EL', 'name' => 'iDash Electronics Ltd', 'description' => 'Consumer durables and smart appliances'],
            ['code' => 'ID-EN', 'name' => 'iDash Energy Solutions', 'description' => 'Industrial power systems and automation'],
            ['code' => 'ID-DI', 'name' => 'iDash Digital Industries', 'description' => 'Smart manufacturing and IoT services'],
        ];

        $records = [];
        foreach ($companies as $company) {
            $id = DB::table('companies')->insertGetId([
                'name' => $company['name'],
                'code' => $company['code'],
                'description' => $company['description'],
                'is_active' => true,
                'created_at' => now(),
                'updated_at' => now(),
            ]);
            $records[] = array_merge($company, ['id' => $id]);
        }

        return $records;
    }

    private function seedChannels(): array
    {
        $channels = ['Corporate', 'Distribution', 'Retail', 'Online'];
        $records = [];
        foreach ($channels as $name) {
            $id = DB::table('channels')->insertGetId([
                'name' => $name,
                'status' => 1,
                'created_at' => now(),
                'updated_at' => now(),
            ]);
            $records[] = ['id' => $id, 'name' => $name];
        }

        return $records;
    }

    private function seedInventoryAccounts(): array
    {
        $accounts = [
            ['gl_account' => '110000', 'gl_account_name' => 'Finished Goods'],
            ['gl_account' => '120000', 'gl_account_name' => 'Raw Materials'],
            ['gl_account' => '130000', 'gl_account_name' => 'Work In Progress'],
        ];

        $records = [];
        foreach ($accounts as $account) {
            $id = DB::table('inventory_gl_accounts')->insertGetId([
                'gl_account' => $account['gl_account'],
                'gl_account_name' => $account['gl_account_name'],
                'description' => 'Demo account',
                'created_at' => now(),
                'updated_at' => now(),
            ]);
            $records[] = array_merge($account, ['id' => $id]);
        }

        return $records;
    }

    private function seedDepartments(array $companies): array
    {
        $departments = collect(['Sales', 'Production', 'Finance', 'Supply Chain', 'HR']);
        $records = [];

        foreach ($companies as $company) {
            $companyId = DB::table('hris_companies')->insertGetId([
                'name' => $company['name'],
                'code' => $company['code'],
                'created_at' => now(),
                'updated_at' => now(),
            ]);

            foreach ($departments as $name) {
                $id = DB::table('hris_departments')->insertGetId([
                    'name' => $name,
                    'company_id' => $companyId,
                    'created_at' => now(),
                    'updated_at' => now(),
                ]);
                $records[] = ['id' => $id, 'company_id' => $company['id'], 'name' => $name];
            }
        }

        return $records;
    }

    private function seedFinanceStructure(): array
    {
        $structure = [
            'Operations' => [
                [
                    'department' => 'Manufacturing Finance',
                    'group' => 'Plant Operations',
                    'sub_heads' => [
                        'Utilities' => ['Electricity Costs', 'Water & Steam'],
                        'Maintenance' => ['Spare Parts', 'Service Contracts'],
                    ],
                ],
                [
                    'department' => 'Logistics Finance',
                    'group' => 'Supply Chain Support',
                    'sub_heads' => [
                        'Transportation' => ['Freight Charges', 'Fuel Usage'],
                        'Warehousing' => ['Storage Rent', 'Handling Equipment'],
                    ],
                ],
            ],
            'Sales & Marketing' => [
                [
                    'department' => 'Sales Finance',
                    'group' => 'Revenue Enablement',
                    'sub_heads' => [
                        'Incentives' => ['Sales Commission', 'Distributor Bonus'],
                        'Field Operations' => ['Travel & Lodging', 'Channel Promotions'],
                    ],
                ],
                [
                    'department' => 'Marketing Finance',
                    'group' => 'Brand Growth',
                    'sub_heads' => [
                        'Campaigns' => ['Digital Ads', 'Events & Exhibitions'],
                        'Research' => ['Market Research', 'Customer Surveys'],
                    ],
                ],
            ],
            'Innovation & R&D' => [
                [
                    'department' => 'Product Development Finance',
                    'group' => 'Innovation Projects',
                    'sub_heads' => [
                        'Prototype' => ['Materials', 'Testing Fees'],
                        'Talent' => ['Consultants', 'Training'],
                    ],
                ],
            ],
        ];

        $categoryRecords = [];
        $departmentRecords = [];
        $expenseRecords = [];

        foreach ($structure as $categoryName => $departmentConfigs) {
            $categoryId = DB::table('bdgt_categories')->insertGetId([
                'name' => $categoryName,
                'status' => 1,
                'created_at' => now(),
                'updated_at' => now(),
            ]);
            $categoryRecords[] = ['id' => $categoryId, 'name' => $categoryName];

            foreach ($departmentConfigs as $config) {
                $departmentId = DB::table('bdgt_departments')->insertGetId([
                    'name' => $config['department'],
                    'category_id' => $categoryId,
                    'status' => 1,
                    'created_at' => now(),
                    'updated_at' => now(),
                ]);
                $departmentRecords[] = [
                    'id' => $departmentId,
                    'category_id' => $categoryId,
                    'name' => $config['department'],
                ];

                $groupId = DB::table('bdgt_expense_groups')->insertGetId([
                    'name' => $config['group'],
                    'department_id' => $departmentId,
                    'status' => 1,
                    'created_at' => now(),
                    'updated_at' => now(),
                ]);

                foreach ($config['sub_heads'] as $subHeadName => $expenses) {
                    $subHeadId = DB::table('bdgt_sub_heads')->insertGetId([
                        'name' => $subHeadName,
                        'expense_group_id' => $groupId,
                        'status' => 1,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ]);

                    foreach ($expenses as $expenseName) {
                        $expenseId = DB::table('bdgt_expenses')->insertGetId([
                            'name' => $expenseName,
                            'sub_head_id' => $subHeadId,
                            'created_at' => now(),
                            'updated_at' => now(),
                        ]);

                        $expenseRecords[] = [
                            'id' => $expenseId,
                            'sub_head_id' => $subHeadId,
                            'name' => $expenseName,
                        ];
                    }
                }
            }
        }

        return [
            'categories' => $categoryRecords,
            'departments' => $departmentRecords,
            'expenses' => $expenseRecords,
        ];
    }

    private function seedBankLoanHeads(): array
    {
        if (!Schema::hasTable('bank_loan_heads')) {
            return [];
        }

        $heads = ['Term Loan', 'Working Capital', 'Equipment Finance', 'Project Finance'];
        $records = [];

        foreach ($heads as $head) {
            $id = DB::table('bank_loan_heads')->insertGetId([
                'loan_head' => $head,
                'created_at' => now(),
                'updated_at' => now(),
            ]);
            $records[] = ['id' => $id, 'loan_head' => $head];
        }

        return $records;
    }

    private function seedSuppliers(array $companies): array
    {
        $supplierSets = [];
        foreach ($companies as $company) {
            $supplierSets[$company['id']] = collect(['Global Metals', 'Prime Plastics', 'Nimbus Logistics', 'ElectroParts', 'BrightChem'])->map(function ($name, $idx) use ($company) {
                $code = strtoupper(substr($name, 0, 3)) . sprintf('%02d', $idx + 1);
                return ['company_id' => $company['id'], 'code' => $code, 'name' => $name];
            })->all();
        }
        return $supplierSets;
    }

    private function seedSalesData(array $companies, array $channels, array $months, array $days, $faker, int $years): void
    {
        $monthlyRows = [];
        $dailyRows = [];
        $bestProducts = [];
        $bestGroups = [];
        $orderDelivery = [];

        $baseProducts = ['Solar Inverter', 'Smart Switch', 'Power Battery', 'Automation Hub', 'Cooling Unit'];
        $productGroups = ['Energy', 'Electronics', 'Automation'];

        foreach ($months as $month) {
            foreach ($companies as $company) {
                foreach ($channels as $channel) {
                    $billed = mt_rand(60, 120) * 10_000_000;
                    $liftingTarget = $billed * mt_rand(105, 120) / 100;
                    $primary = $billed * mt_rand(85, 95) / 100;
                    $imsTarget = $billed * mt_rand(60, 75) / 100;
                    $ims = $imsTarget * mt_rand(90, 110) / 100;
                    $memoTarget = $billed * 0.1;
                    $memoQty = $memoTarget * mt_rand(90, 110) / 100;
                    $pgTarget = $billed * 0.2;
                    $pgCover = $pgTarget * mt_rand(95, 105) / 100;
                    $retailers = mt_rand(200, 500);

                    $monthlyRows[] = [
                        'data_month' => $month->toDateString(),
                        'channel_id' => $channel['id'],
                        'channel_name' => $channel['name'],
                        'lifting_target' => round($liftingTarget, 2),
                        'billed' => round($billed, 2),
                        'delivered' => round($billed * mt_rand(92, 105) / 100, 2),
                        'primary_collection' => round($primary, 2),
                        'ims_target' => round($imsTarget, 2),
                        'ims' => round($ims, 2),
                        'market_collection' => round($ims * mt_rand(90, 110) / 100, 2),
                        'memo_target' => round($memoTarget, 2),
                        'memo_qty' => round($memoQty, 2),
                        'pg_target' => round($pgTarget, 2),
                        'pg_cover' => round($pgCover, 2),
                        'total_retailer' => $retailers,
                        'business_retailer' => (int) round($retailers * mt_rand(70, 90) / 100),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];

                    $orderDelivery[] = [
                        'months' => $month->toDateString(),
                        'channel_id' => $channel['id'],
                        'amounts' => round($billed * mt_rand(105, 120) / 100, 2),
                        'types' => 0,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                    $orderDelivery[] = [
                        'months' => $month->toDateString(),
                        'channel_id' => $channel['id'],
                        'amounts' => round($billed * mt_rand(92, 101) / 100, 2),
                        'types' => 1,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];

                    foreach ($faker->randomElements($baseProducts, 3) as $product) {
                        $bestProducts[] = [
                            'year_month' => $month->toDateString(),
                            'channel_id' => $channel['id'],
                            'product_id' => $faker->numberBetween(1000, 9999),
                            'product_name' => $product,
                            'qty' => round(mt_rand(2, 6) * 1_000 + mt_rand(0, 500), 2),
                            'value' => round($billed * mt_rand(5, 12) / 100, 2),
                            'cat_id' => mt_rand(1, 5),
                            'created_at' => now(),
                            'updated_at' => now(),
                        ];
                    }

                    foreach ($productGroups as $groupName) {
                        $bestGroups[] = [
                            'year_month' => $month->toDateString(),
                            'channel_id' => $channel['id'],
                            'category_id' => mt_rand(1, 10),
                            'category_name' => $groupName,
                            'qty' => round(mt_rand(5, 10) * 1_000, 2),
                            'value' => round($billed * mt_rand(15, 30) / 100, 2),
                            'created_at' => now(),
                            'updated_at' => now(),
                        ];
                    }
                }
            }
        }

        foreach ($days as $day) {
            foreach ($channels as $channel) {
                $dailyRows[] = [
                    'data_date' => $day->toDateString(),
                    'channel_id' => $channel['id'],
                    'channel_name' => $channel['name'],
                    'lifting_target' => round(mt_rand(15, 30) * 1_000_000, 2),
                    'billed' => round(mt_rand(12, 24) * 1_000_000, 2),
                    'delivery' => round(mt_rand(10, 22) * 1_000_000, 2),
                    'lifting_collection' => round(mt_rand(9, 20) * 1_000_000, 2),
                    'ims_target' => round(mt_rand(8, 16) * 1_000_000, 2),
                    'ims' => round(mt_rand(7, 15) * 1_000_000, 2),
                    'ims_collection' => round(mt_rand(6, 14) * 1_000_000, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        DB::table('channelwise_monthly_report')->insert($monthlyRows);
        DB::table('order_delivery_summaries')->insert($orderDelivery);
        foreach (array_chunk($dailyRows, 1000) as $chunk) {
            DB::table('channelwise_lic_data')->insert($chunk);
        }
        DB::table('best_selling_products')->insert($bestProducts);
        DB::table('best_selling_pgs')->insert($bestGroups);

        $topDistributors = [];
        $topRetailers = [];
        for ($i = 0; $i < 30; $i++) {
            $day = Carbon::now()->subDays($i);
            $topDistributors[] = [
                'db_name' => 'Distributor ' . ($i + 1),
                'amount' => round(mt_rand(20, 60) * 10_000_000, 2),
                'type' => 0,
                'date' => $day->toDateString(),
                'created_at' => now(),
                'updated_at' => now(),
            ];
            $topRetailers[] = [
                'date' => $day->toDateString(),
                'db_name' => 'Retailer ' . ($i + 1),
                'amount' => round(mt_rand(8, 20) * 10_000_000, 2),
                'created_at' => now(),
                'updated_at' => now(),
            ];
        }
        DB::table('top_channel_d_bs')->insert($topDistributors);
        DB::table('top_retailers')->insert($topRetailers);
    }

    private function seedProductionData(array $companies, array $months, array $days, $faker): void
    {
        $analysisRows = [];
        $wastageRows = [];
        $costRows = [];
        $efficiencyRows = [];
        $maintenanceRows = [];
        $energyRows = [];

        $factories = ['North Plant', 'Central Plant', 'South Plant'];
        foreach ($months as $month) {
            foreach ($companies as $company) {
                foreach ($factories as $factory) {
                    $baseline = mt_rand(40, 80) * 10_000_000;
                    $analysisRows[] = [
                        'month' => $month->toDateString(),
                        'factory' => $factory,
                        'summary_group' => 'Gross Profit',
                        'cmonthly_amount' => round($baseline * 1.1, 2),
                        'cavg_amount' => round($baseline * 1.05, 2),
                        'cmonthly_per' => mt_rand(45, 60),
                        'cavg_per' => mt_rand(40, 55),
                        'pmonthly_amount' => round($baseline * 0.9, 2),
                        'pavg_amount' => round($baseline * 0.88, 2),
                        'pmonthly_per' => mt_rand(35, 50),
                        'pavg_per' => mt_rand(30, 45),
                        'tmonthly_amount' => round($baseline * 1.5, 2),
                        'tmonthly_per' => mt_rand(60, 80),
                        'amonthly_amount' => round($baseline * 1.3, 2),
                        'amonthly_per' => mt_rand(55, 75),
                        'aavg_amount' => round($baseline * 1.2, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];

                    $wastageRows[] = [
                        'month' => $month->toDateString(),
                        'factory' => $factory,
                        'group_name' => 'Material',
                        'std' => mt_rand(2, 5),
                        'wastage' => round(mt_rand(150, 350), 2),
                        'month_wastage' => round(mt_rand(200, 400), 2),
                        'avg' => round(mt_rand(180, 360), 2),
                        'amount' => round(mt_rand(5, 12) * 10_000_000, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];

                    foreach (['Labor', 'Utilities', 'Maintenance'] as $costType) {
                        $costRows[] = [
                            'month' => $month->toDateString(),
                            'factory' => $factory,
                            'cost_type' => $costType,
                            'amount' => round(mt_rand(8, 20) * 10_000_000, 2),
                            'created_at' => now(),
                            'updated_at' => now(),
                        ];
                    }
                }
            }
        }

        foreach ($days as $day) {
            foreach ($companies as $company) {
                foreach ($factories as $factory) {
                    $planned = mt_rand(180, 240);
                    $actual = mt_rand(160, 230);
                    $efficiency = $planned > 0 ? round(($actual / $planned) * 100, 2) : 0;
                    $efficiencyRows[] = [
                        'company_id' => $company['id'],
                        'factory_id' => $company['id'],
                        'production_line_id' => null,
                        'production_date' => $day->toDateString(),
                        'shift' => Arr::random(['A', 'B', 'C']),
                        'planned_output' => $planned,
                        'actual_output' => $actual,
                        'efficiency_percentage' => $efficiency,
                        'planned_hours' => 8,
                        'actual_hours' => mt_rand(7, 9),
                        'downtime_minutes' => mt_rand(15, 60),
                        'oee' => round($efficiency * mt_rand(70, 90) / 100, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }
            }
        }

        for ($i = 0; $i < 120; $i++) {
            $day = Carbon::now()->subDays($i);
            $maintenanceRows[] = [
                'company_id' => Arr::random($companies)['id'],
                'factory_id' => Arr::random($companies)['id'],
                'machine_code' => 'MC-' . mt_rand(100, 999),
                'machine_name' => 'Machine #' . mt_rand(1, 50),
                'maintenance_type' => Arr::random(['preventive', 'corrective', 'breakdown']),
                'maintenance_date' => $day->toDateString(),
                'start_time' => '08:00:00',
                'end_time' => '12:00:00',
                'downtime_minutes' => mt_rand(30, 180),
                'description' => 'Routine service',
                'actions_taken' => 'Adjusted components and lubricated parts',
                'cost' => round(mt_rand(5, 15) * 1_000_000, 2),
                'status' => 'completed',
                'technician_name' => 'Technician ' . mt_rand(1, 50),
                'created_at' => now(),
                'updated_at' => now(),
            ];
        }

        foreach ($months as $month) {
            foreach ($companies as $company) {
                foreach (['electricity', 'gas', 'water'] as $energyType) {
                    $energyRows[] = [
                        'company_id' => $company['id'],
                        'factory_id' => $company['id'],
                        'consumption_date' => $month->toDateString(),
                        'energy_type' => $energyType,
                        'consumption_amount' => round(mt_rand(500, 1200), 2),
                        'unit' => 'MWh',
                        'cost' => round(mt_rand(5, 12) * 1_000_000, 2),
                        'meter_reading' => 'MR-' . mt_rand(1000, 9999),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }
            }
        }

        DB::table('production_analyses')->insert($analysisRows);
        DB::table('wastage_datas')->insert($wastageRows);
        DB::table('cost_analyses')->insert($costRows);
        DB::table('production_efficiency')->insert($efficiencyRows);
        DB::table('machine_maintenances')->insert($maintenanceRows);
        DB::table('energy_consumption')->insert($energyRows);
    }

    private function seedFinanceData(
        array $companies,
        array $categories,
        array $departments,
        array $expenses,
        array $loanHeads,
        array $months
    ): void {
        if (empty($categories) || empty($departments)) {
            return;
        }

        $budgetSummaries = [];
        $budgetMonthlies = [];
        $financialExpenseRows = [];
        $loanStatusRows = [];

        $departmentsByCategory = collect($departments)->groupBy('category_id');

        foreach ($months as $month) {
            foreach ($categories as $category) {
                $deptList = $departmentsByCategory->get($category['id'], collect());
                foreach ($deptList as $department) {
                    $budget = mt_rand(5, 15) * 10_000_000;
                    $actual = $budget * mt_rand(80, 120) / 100;
                    $budgetSummaries[] = [
                        'month' => $month->toDateString(),
                        'category_id' => $category['id'],
                        'department_id' => $department['id'],
                        'budget_amount' => round($budget, 2),
                        'actual_amount' => round($actual, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }
            }

            foreach ($expenses as $expense) {
                $budget = mt_rand(1, 5) * 1_000_000;
                $actual = $budget * mt_rand(80, 120) / 100;
                $budgetMonthlies[] = [
                    'month' => $month->toDateString(),
                    'expense_id' => $expense['id'],
                    'budget_amount' => round($budget, 2),
                    'actual_amount' => round($actual, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
                $financialExpenseRows[] = [
                    'month' => $month->toDateString(),
                    'expense_id' => $expense['id'],
                    'amount' => round($actual, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }

            foreach ($companies as $company) {
                foreach ($loanHeads as $head) {
                    $loanAmount = mt_rand(100, 220) * 10_000_000;
                    $loanStatusRows[] = [
                        'month' => $month->toDateString(),
                        'loan_head' => $head['loan_head'] ?? $head['name'] ?? 'Loan',
                        'company_id' => $company['code'],
                        'amount' => round($loanAmount, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }
            }
        }

        if (Schema::hasTable('budget_summaries') && !empty($budgetSummaries)) {
            DB::table('budget_summaries')->insert($budgetSummaries);
        }
        if (Schema::hasTable('budget_monthlies') && !empty($budgetMonthlies)) {
            DB::table('budget_monthlies')->insert($budgetMonthlies);
        }
        if (Schema::hasTable('financial_expense_raw_data') && !empty($financialExpenseRows)) {
            DB::table('financial_expense_raw_data')->insert($financialExpenseRows);
        }
        if (Schema::hasTable('bank_loan_status_raw_data') && !empty($loanStatusRows)) {
            DB::table('bank_loan_status_raw_data')->insert($loanStatusRows);
        }
    }

    private function seedInventoryData(array $companies, array $glAccounts, array $months, $faker): void
    {
        $inventoryRows = [];
        $cogsRows = [];
        $sapRows = [];

        foreach ($months as $month) {
            foreach ($companies as $company) {
                foreach ($glAccounts as $account) {
                    $inventoryRows[] = [
                        'company_id' => $company['id'],
                        'gl_id' => $account['id'],
                        'month' => $month->toDateString(),
                        'amount' => round(mt_rand(40, 90) * 10_000_000, 2),
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }

                $cogs = mt_rand(80, 140) * 10_000_000;
                $gp = mt_rand(30, 60) * 10_000_000;
                $cogsRows[] = [
                    'month' => $month->toDateString(),
                    'company_id' => $company['id'],
                    'cogs' => round($cogs, 2),
                    'gp' => round($gp, 2),
                    'gp_percentage' => round($gp / ($cogs + $gp) * 100, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $sapRows[] = [
                    'company' => $company['id'],
                    'year' => $month->year,
                    'month' => $month->month,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        DB::table('inventory_raw_datas')->insert($inventoryRows);
        DB::table('cogs_gps')->insert($cogsRows);
        DB::table('inventroy_sap_datas')->insert($sapRows);
    }

    private function seedHRData(array $companies, array $departments, array $months, array $days, $faker): void
    {
        $basicRows = [];
        $attendanceRows = [];
        $turnoverRows = [];
        $promotionRows = [];
        $promotionBreakRows = [];

        foreach ($months as $month) {
            $year = $month->year;
            foreach ($companies as $company) {
                $staff = mt_rand(500, 750);
                $workers = mt_rand(400, 600);
                $permanent = (int) round(($staff + $workers) * 0.7);
                $contract = (int) round(($staff + $workers) * 0.15);

                $basicRows[] = [
                    'total_active_staff' => $staff,
                    'total_active_worker' => $workers,
                    'total_contractual_employee' => $contract,
                    'total_probationary_employee' => mt_rand(30, 60),
                    'total_permanent_employee' => $permanent,
                    'report_date' => $month->toDateString(),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $turnoverRows[] = [
                    'job_type' => 'staff',
                    'month' => $month->format('F'),
                    'year' => $year,
                    'new_employee_no' => mt_rand(10, 25),
                    'resigned_employee' => mt_rand(5, 15),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
                $turnoverRows[] = [
                    'job_type' => 'worker',
                    'month' => $month->format('F'),
                    'year' => $year,
                    'new_employee_no' => mt_rand(12, 30),
                    'resigned_employee' => mt_rand(6, 18),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        foreach ($days as $day) {
            $attendanceRows[] = [
                'date' => $day->toDateString(),
                'total_present' => mt_rand(950, 1200),
                'total_absent' => mt_rand(40, 80),
                'total_leave' => mt_rand(15, 40),
                'created_at' => now(),
                'updated_at' => now(),
            ];
        }

        foreach (range(Carbon::now()->year - 4, Carbon::now()->year) as $year) {
            $promotionCount = mt_rand(30, 60);
            $promotionRows[] = [
                'year' => $year,
                'promoted_count' => $promotionCount,
                'details' => 'Yearly promotion summary',
                'created_at' => now(),
                'updated_at' => now(),
            ];

            foreach ($departments as $department) {
                $promotionBreakRows[] = [
                    'year' => $year,
                    'department_id' => $department['id'],
                    'promoted_count' => mt_rand(3, 12),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        foreach (array_chunk($basicRows, 1000) as $chunk) {
            DB::table('employee_basic_infos')->insert($chunk);
        }
        foreach (array_chunk($attendanceRows, 2000) as $chunk) {
            DB::table('employee_attendances')->insert($chunk);
        }
        DB::table('employee_tran_overs')->insert($turnoverRows);
        DB::table('yearly_employee_promotions')->insert($promotionRows);

        if (Schema::hasTable('hris_promotion_break_downs')) {
            DB::table('hris_promotion_break_downs')->insert($promotionBreakRows);
        }
    }

    private function seedSupplyChainData(array $companies, array $months, $faker): void
    {
        $poRows = [];
        $masterRows = [];
        $grnRows = [];
        $invoiceRows = [];
        $rawRows = [];
        $purchaseRows = [];

        foreach ($months as $month) {
            foreach ($companies as $company) {
                $poNumber = mt_rand(100000, 999999);
                $poValue = mt_rand(25, 60) * 10_000_000;
                $grnValue = $poValue * mt_rand(80, 100) / 100;
                $invoiceValue = $grnValue * mt_rand(85, 105) / 100;

                $poRows[] = [
                    'company_id' => $company['id'],
                    'plant' => mt_rand(100, 199),
                    'vendor_code' => mt_rand(5000, 5999),
                    'vandor_name' => 'Vendor ' . mt_rand(1, 50),
                    'material_code' => 'MAT-' . mt_rand(1000, 9999),
                    'material_name' => 'Raw Material ' . mt_rand(1, 20),
                    'material_group' => 'Group ' . mt_rand(1, 5),
                    'material_group_description' => 'Material group description',
                    'purchasing_organization' => Arr::random(['North', 'South', 'Central']),
                    'purchasing_group' => 'GRP-' . mt_rand(10, 99),
                    'pr_number' => mt_rand(100000, 199999),
                    'po_number' => $poNumber,
                    'po_item_number' => mt_rand(10, 99),
                    'po_date' => $month->toDateString(),
                    'po_qty' => mt_rand(100, 500),
                    'uom' => 'EA',
                    'po_currency' => 'INR',
                    'po_amount' => round($poValue, 2),
                    'master_po' => round($poValue * 1.05, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $masterRows[] = [
                    'company' => (string) $company['id'],
                    'po_number' => (string) $poNumber,
                    'po_date' => $month->toDateString(),
                    'po_value' => round($poValue, 2),
                    'pr_amount' => round($poValue * 0.6, 2),
                    'purchase_org' => Arr::random(['North', 'South', 'Central']),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $grnRows[] = [
                    'company' => $company['id'],
                    'po_number' => (string) $poNumber,
                    'grn_date' => $month->copy()->addDays(mt_rand(5, 15))->toDateString(),
                    'grn_amount' => round($grnValue, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $invoiceRows[] = [
                    'company' => $company['id'],
                    'invoice_number' => 'INV-' . $poNumber,
                    'iv_date' => $month->copy()->addDays(mt_rand(15, 25))->toDateString(),
                    'total_invoice' => round($invoiceValue, 2),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $rawRows[] = [
                    'plant' => mt_rand(100, 199),
                    'pr_id' => mt_rand(1000, 9999),
                    'pr_item' => mt_rand(1, 20),
                    'pr_date' => $month->copy()->subDays(mt_rand(10, 20))->toDateString(),
                    'pr_material' => 'PR-MAT-' . mt_rand(100, 999),
                    'pr_material_text' => 'Material request',
                    'pr_qty' => mt_rand(100, 250),
                    'prunit_measure' => 'EA',
                    'prrelease_date' => $month->copy()->subDays(mt_rand(5, 10))->toDateString(),
                    'po_id' => $poNumber,
                    'po_item' => mt_rand(10, 99),
                    'po_date' => $month->toDateString(),
                    'vendor_id' => mt_rand(5000, 5999),
                    'vendor_details' => 'Vendor description',
                    'po_qty' => mt_rand(100, 250),
                    'pounit_measure' => 'EA',
                    'po_amount' => round($poValue, 2),
                    'po_currency' => 'INR',
                    'po_delivery_date' => $month->copy()->addDays(mt_rand(7, 20))->toDateString(),
                    'po_released_date' => $month->copy()->addDays(mt_rand(1, 5))->toDateString(),
                    'lc_number' => mt_rand(10000, 99999),
                    'actual_po1' => mt_rand(1000, 9999),
                    'actual_po2' => mt_rand(1000, 9999),
                    'actual_po1item' => mt_rand(1, 20),
                    'actual_po2item' => mt_rand(1, 20),
                    'grn1_id' => mt_rand(1000, 9999),
                    'grn1_item' => mt_rand(1, 20),
                    'grn1_date' => $month->copy()->addDays(mt_rand(5, 15))->toDateString(),
                    'grn1_qtn' => mt_rand(80, 200),
                    'grn1_amount' => round($grnValue, 2),
                    'invoice1_date' => $month->copy()->addDays(mt_rand(20, 30))->toDateString(),
                    'invoice1_id' => 'INV-' . mt_rand(1000, 9999),
                    'invoice1_vendor_date' => $month->copy()->addDays(mt_rand(18, 28))->toDateString(),
                    'invoice1_qty' => mt_rand(80, 200),
                    'invoice1_amount' => round($invoiceValue, 2),
                    'invoice1_dn' => mt_rand(1000, 9999),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
                $purchaseRows[] = [
                    'pr_id' => mt_rand(100000, 199999),
                    'pr_item' => mt_rand(1, 20),
                    'pr_date' => $month->copy()->subDays(mt_rand(15, 30))->toDateString(),
                    'material_code' => 'MAT-' . mt_rand(1000, 9999),
                    'material_name' => 'Material ' . mt_rand(1, 25),
                    'quantity' => mt_rand(80, 200),
                    'unit' => 'EA',
                    'plant' => mt_rand(100, 199),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        if (Schema::hasTable('supply_chain_master_datas') && !empty($masterRows)) {
            DB::table('supply_chain_master_datas')->insert($masterRows);
        }
        if (Schema::hasTable('supply_chain_pos') && !empty($poRows)) {
            DB::table('supply_chain_pos')->insert($poRows);
        }
        if (Schema::hasTable('supply_chain_grn_datas') && !empty($grnRows)) {
            DB::table('supply_chain_grn_datas')->insert($grnRows);
        }
        if (Schema::hasTable('supply_chain_invoice_datas') && !empty($invoiceRows)) {
            DB::table('supply_chain_invoice_datas')->insert($invoiceRows);
        }
        if (Schema::hasTable('supply_chain_raw_datas') && !empty($rawRows)) {
            DB::table('supply_chain_raw_datas')->insert($rawRows);
        }
        if (Schema::hasTable('purchase_requisitions') && !empty($purchaseRows)) {
            DB::table('purchase_requisitions')->insert($purchaseRows);
        }
    }

    private function seedManufacturingData(array $companies, array $suppliers, array $months, array $days, $faker): void
    {
        $qcRows = [];
        $planRows = [];
        $materialRows = [];
        $supplierRows = [];

        foreach ($days as $day) {
            foreach ($companies as $company) {
                $qcRows[] = [
                    'company_id' => $company['id'],
                    'factory_id' => $company['id'],
                    'production_line_id' => null,
                    'product_code' => 'PR-' . mt_rand(100, 999),
                    'product_name' => 'Product ' . mt_rand(1, 50),
                    'check_date' => $day->toDateString(),
                    'check_type' => Arr::random(['incoming', 'in_process', 'final']),
                    'status' => Arr::random(['passed', 'passed', 'failed']),
                    'total_checked' => mt_rand(80, 140),
                    'passed_count' => mt_rand(70, 120),
                    'failed_count' => mt_rand(0, 10),
                    'defect_rate' => round(mt_rand(1, 8), 2),
                    'defects' => $faker->sentence(6),
                    'remarks' => $faker->sentence(8),
                    'inspector_name' => $faker->name(),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        foreach ($months as $month) {
            foreach ($companies as $company) {
                $planRows[] = [
                    'company_id' => $company['id'],
                    'factory_id' => $company['id'],
                    'plan_number' => 'PLAN-' . $month->format('Ym') . '-' . $company['id'],
                    'product_code' => 'PRD-' . mt_rand(100, 999),
                    'product_name' => 'Planned Product',
                    'plan_date' => $month->toDateString(),
                    'start_date' => $month->copy()->addDays(1)->toDateString(),
                    'end_date' => $month->copy()->addDays(25)->toDateString(),
                    'planned_quantity' => round(mt_rand(400, 600), 2),
                    'actual_quantity' => round(mt_rand(350, 580), 2),
                    'status' => Arr::random(['approved', 'completed']),
                    'priority' => Arr::random(['medium', 'high', 'urgent']),
                    'notes' => 'Demo plan',
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                $materialRows[] = [
                    'company_id' => $company['id'],
                    'production_plan_id' => count($planRows),
                    'material_code' => 'MAT-' . mt_rand(1000, 9999),
                    'material_name' => 'Material ' . mt_rand(1, 50),
                    'material_type' => Arr::random(['raw', 'component', 'packaging']),
                    'unit' => 'kg',
                    'required_quantity' => round(mt_rand(200, 350), 2),
                    'available_quantity' => round(mt_rand(150, 320), 2),
                    'shortage_quantity' => round(mt_rand(0, 40), 2),
                    'required_date' => $month->copy()->addDays(10)->toDateString(),
                    'status' => Arr::random(['pending', 'ordered', 'received']),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                foreach ($suppliers[$company['id']] as $supplier) {
                    $orders = mt_rand(15, 30);
                    $onTime = mt_rand((int) ($orders * 0.7), $orders);
                    $qualityScore = mt_rand(80, 98) / 1.0;
                    $supplierRows[] = [
                        'company_id' => $company['id'],
                        'supplier_code' => $supplier['code'],
                        'supplier_name' => $supplier['name'],
                        'evaluation_date' => $month->toDateString(),
                        'total_orders' => $orders,
                        'on_time_deliveries' => $onTime,
                        'on_time_percentage' => round(($onTime / $orders) * 100, 2),
                        'quality_issues' => mt_rand(0, 3),
                        'quality_score' => round($qualityScore, 2),
                        'cost_score' => round(mt_rand(80, 95), 2),
                        'overall_score' => round(($qualityScore + mt_rand(80, 95)) / 2, 2),
                        'rating' => Arr::random(['excellent', 'good', 'average']),
                        'comments' => 'Supplier performance review',
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];
                }
            }
        }

        if (Schema::hasTable('quality_control_checks')) {
            foreach (array_chunk($qcRows, 1000) as $chunk) {
                DB::table('quality_control_checks')->insert($chunk);
            }
        }
        if (Schema::hasTable('production_plans')) {
            DB::table('production_plans')->insert($planRows);
        }
        if (Schema::hasTable('material_requirements')) {
            DB::table('material_requirements')->insert($materialRows);
        }
        if (Schema::hasTable('supplier_performance')) {
            DB::table('supplier_performance')->insert($supplierRows);
        }
    }

    private function seedNPDData(array $companies, $faker): void
    {
        if (!Schema::hasTable('npd_projects')) {
            return;
        }

        $projectRows = [];
        $deliverableRows = [];
        $subDeliverableRows = [];

        foreach ($companies as $company) {
            foreach (range(1, 3) as $index) {
                $start = Carbon::now()->subMonths(mt_rand(6, 18));
                $end = $start->copy()->addMonths(mt_rand(4, 10));

                $projectId = count($projectRows) + 1;
                $projectRows[] = [
                    'p_id' => 1000 + $projectId,
                    'indent_no' => 'IND-' . mt_rand(10000, 99999),
                    'name' => 'Project ' . strtoupper($faker->lexify('???')),
                    'pmo' => 'PMO Office',
                    'project_manager' => $faker->name(),
                    'location' => $faker->city(),
                    'reason' => 'Portfolio expansion',
                    'details' => 'Pilot project for new product line',
                    'sponsors' => 'Executive Board',
                    'start_date' => $start->toDateString(),
                    'end_date' => $end->toDateString(),
                    'type' => Arr::random(['growth', 'sustaining', 'productivity']),
                    'progress' => mt_rand(50, 95) . '%',
                    'budget' => round(mt_rand(5, 10) * 1_000_000, 2),
                    'lead_department' => Arr::random(['R&D', 'Engineering', 'Innovation']),
                    'reponsible_departments' => 'R&D, Engineering',
                    'status' => Arr::random(['on_track', 'at_risk', 'delayed']),
                    'status_name' => Arr::random(['On Track', 'Slightly Behind', 'Needs Attention']),
                    'status_background' => Arr::random(['success', 'warning', 'danger']),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

                foreach (range(1, 3) as $deliverableIndex) {
                    $deliverableId = count($deliverableRows) + 1;
                    $deliverableRows[] = [
                        'd_id' => 2000 + $deliverableId,
                        'name' => 'Deliverable ' . $deliverableId,
                        'weightage' => mt_rand(10, 30),
                        'start_date' => $start->copy()->addWeeks($deliverableIndex)->toDateString(),
                        'end_date' => $start->copy()->addWeeks($deliverableIndex * 2)->toDateString(),
                        'acknowledges' => 'PMO',
                        'budget' => round(mt_rand(1, 3) * 500_000, 2),
                        'progress' => mt_rand(40, 90),
                        'npd_project_id' => $projectId,
                        'created_at' => now(),
                        'updated_at' => now(),
                    ];

                    foreach (range(1, 2) as $subIndex) {
                        $subDeliverableRows[] = [
                            'sd_id' => 3000 + count($subDeliverableRows) + 1,
                            'name' => 'Task ' . $deliverableId . '-' . $subIndex,
                            'weightage' => mt_rand(5, 15),
                            'start_date' => $start->copy()->addWeeks($subIndex)->toDateString(),
                            'end_date' => $start->copy()->addWeeks($subIndex + 1)->toDateString(),
                            'acknowledges' => 'Team Lead',
                            'budget' => round(mt_rand(1, 2) * 200_000, 2),
                            'progress' => mt_rand(30, 90),
                            'deliverable_id' => $deliverableId,
                            'created_at' => now(),
                            'updated_at' => now(),
                        ];
                    }
                }
            }
        }

        DB::table('npd_projects')->insert($projectRows);
        if (Schema::hasTable('projects_deliverables')) {
            DB::table('projects_deliverables')->insert($deliverableRows);
        }
        if (Schema::hasTable('projects_sub_deliverables')) {
            DB::table('projects_sub_deliverables')->insert($subDeliverableRows);
        }
    }

    private function seedForecasts(array $companies, array $months, $faker): void
    {
        $forecastTypes = ['sales', 'production', 'finance', 'inventory', 'hr', 'supplychain'];
        $forecastRows = [];
        $forecastDetails = [];

        $futureMonths = collect(CarbonPeriod::create(Carbon::now()->startOfMonth(), '1 month', Carbon::now()->addMonths(6)->startOfMonth()))
            ->map(fn ($date) => $date->copy())
            ->all();
        foreach ($forecastTypes as $type) {
            foreach ($futureMonths as $future) {
                $value = mt_rand(70, 140) * 10_000_000;
                $upper = $value * 1.1;
                $lower = $value * 0.9;

                $details = [];
                foreach (range(1, 6) as $offset) {
                    $date = $future->copy()->addMonths($offset - 1);
                    $details[] = [
                        'date' => $date->toDateString(),
                        'forecast' => round($value * (1 + ($offset - 3) / 20), 2),
                        'upper_bound' => round($upper, 2),
                        'lower_bound' => round($lower, 2),
                    ];
                }

                $forecastRows[] = [
                    'forecast_type' => $type,
                    'entity_type' => null,
                    'entity_id' => null,
                    'forecast_date' => $future->toDateString(),
                    'forecasted_value' => round($value, 2),
                    'confidence_level' => mt_rand(80, 95),
                    'upper_bound' => round($upper, 2),
                    'lower_bound' => round($lower, 2),
                    'model_used' => Arr::random(['Prophet', 'ARIMA', 'LSTM']),
                    'forecast_details' => json_encode($details),
                    'status' => 'active',
                    'created_at' => now(),
                    'updated_at' => now(),
                ];
            }
        }

        DB::table('ai_forecasts')->insert($forecastRows);
    }
}
