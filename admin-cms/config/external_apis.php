<?php

return [
    'modules' => [
        'sales' => [
            'label' => 'Sales',
            'description' => 'Channel, pipeline, target and promotion metrics aggregated by period.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.records',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'items_path' => 'data.records',
                    'mappings' => [
                        [
                            'target_table' => 'channelwise_monthly_report',
                            'upsert_keys' => ['channel_id', 'data_month'],
                            'fields' => [
                                'channel_id' => 'channel_id',
                                'data_month' => 'period',
                                'channel_name' => 'channel_name',
                                'lifting_target' => 'metrics.lifting_target',
                                'billed' => 'metrics.billed',
                                'primary_collection' => 'metrics.primary_collection',
                                'ims_target' => 'metrics.ims_target',
                                'ims' => 'metrics.ims',
                                'market_collection' => 'metrics.market_collection',
                            ],
                            'transforms' => [
                                'data_month' => 'date:Y-m-d',
                                'lifting_target' => 'float',
                                'billed' => 'float',
                                'primary_collection' => 'float',
                                'ims_target' => 'float',
                                'ims' => 'float',
                                'market_collection' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.targets',
                            'target_table' => 'sales_channel_targets',
                            'upsert_keys' => ['channel_id', 'data_month'],
                            'fields' => [
                                'data_month' => 'period',
                                'channel_id' => 'channel_id',
                                'channel_name' => 'channel_name',
                                'revenue_target' => 'targets.revenue',
                                'volume_target' => 'targets.volume',
                                'promotion_budget' => 'targets.promotion_budget',
                                'gross_margin_target' => 'targets.gross_margin',
                                'new_customer_target' => 'targets.new_customers',
                                'owner' => 'owner',
                            ],
                            'transforms' => [
                                'data_month' => 'date:Y-m-d',
                                'revenue_target' => 'float',
                                'volume_target' => 'float',
                                'promotion_budget' => 'float',
                                'gross_margin_target' => 'float',
                                'new_customer_target' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.orders',
                            'target_table' => 'sales_order_book',
                            'upsert_keys' => ['order_number'],
                            'fields' => [
                                'order_number' => 'order_number',
                                'order_date' => 'dates.order_date',
                                'channel_id' => 'channel.id',
                                'channel_name' => 'channel.name',
                                'customer_code' => 'customer.code',
                                'customer_name' => 'customer.name',
                                'region' => 'customer.region',
                                'status' => 'status',
                                'order_amount' => 'amounts.total',
                                'discount_amount' => 'amounts.discount',
                                'gross_margin' => 'amounts.margin',
                                'fulfilled_at' => 'dates.fulfilled_at',
                            ],
                            'transforms' => [
                                'order_date' => 'date:Y-m-d',
                                'fulfilled_at' => 'date:Y-m-d',
                                'order_amount' => 'float',
                                'discount_amount' => 'float',
                                'gross_margin' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.promotions',
                            'target_table' => 'sales_promotion_performance',
                            'upsert_keys' => ['campaign_code'],
                            'fields' => [
                                'campaign_code' => 'campaign_code',
                                'campaign_name' => 'campaign_name',
                                'channel_id' => 'channel.id',
                                'channel_name' => 'channel.name',
                                'start_date' => 'window.start',
                                'end_date' => 'window.end',
                                'spend_amount' => 'metrics.spend',
                                'revenue_uplift' => 'metrics.uplift_value',
                                'uplift_percentage' => 'metrics.uplift_percent',
                                'roi' => 'metrics.roi',
                                'audience_tags' => 'audience.tags',
                            ],
                            'transforms' => [
                                'start_date' => 'date:Y-m-d',
                                'end_date' => 'date:Y-m-d',
                                'spend_amount' => 'float',
                                'revenue_uplift' => 'float',
                                'uplift_percentage' => 'float',
                                'roi' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'records' => [
                            [
                                'channel_id' => 12,
                                'channel_name' => 'Modern Trade',
                                'period' => '2025-01-01',
                                'metrics' => [
                                    'lifting_target' => 1250000.0,
                                    'billed' => 1175000.0,
                                    'primary_collection' => 985000.0,
                                    'ims_target' => 890000.0,
                                    'ims' => 865000.0,
                                    'market_collection' => 905000.0,
                                ],
                            ],
                        ],
                        'targets' => [
                            [
                                'channel_id' => 12,
                                'channel_name' => 'Modern Trade',
                                'period' => '2025-01-01',
                                'targets' => [
                                    'revenue' => 1500000,
                                    'volume' => 48000,
                                    'promotion_budget' => 95000,
                                    'gross_margin' => 320000,
                                    'new_customers' => 420,
                                ],
                                'owner' => 'Regional Sales Lead',
                            ],
                        ],
                        'orders' => [
                            [
                                'order_number' => 'SO-10045',
                                'channel' => ['id' => 4, 'name' => 'Retail'],
                                'dates' => [
                                    'order_date' => '2025-01-11',
                                    'fulfilled_at' => '2025-01-14',
                                ],
                                'customer' => [
                                    'code' => 'C-882',
                                    'name' => 'Best Mart BD',
                                    'region' => 'Dhaka',
                                ],
                                'amounts' => [
                                    'total' => 218000,
                                    'discount' => 8500,
                                    'margin' => 42000,
                                ],
                                'status' => 'delivered',
                            ],
                        ],
                        'promotions' => [
                            [
                                'campaign_code' => 'DIP-2025-Q1',
                                'campaign_name' => 'January Bundle',
                                'channel' => ['id' => 3, 'name' => 'Distributor'],
                                'window' => [
                                    'start' => '2025-01-01',
                                    'end' => '2025-01-10',
                                ],
                                'metrics' => [
                                    'spend' => 35000,
                                    'uplift_value' => 82000,
                                    'uplift_percent' => 18.5,
                                    'roi' => 235,
                                ],
                                'audience' => [
                                    'tags' => ['bundle', 'new_launch'],
                                ],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'inventory' => [
            'label' => 'Inventory',
            'description' => 'Inventory valuation, turnover and COGS/GP inputs.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'payload.inventory',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                    'as_of' => '{{today}}',
                ],
                'data_mapping' => [
                    'items_path' => 'payload.inventory',
                    'mappings' => [
                        [
                            'target_table' => 'inventory_raw_datas',
                            'upsert_keys' => ['gl_id', 'month', 'company_id'],
                            'fields' => [
                                'gl_id' => 'gl.id',
                                'company_id' => 'company_id',
                                'month' => 'period',
                                'amount' => 'values.on_hand',
                                'status' => 'values.status',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m',
                                'amount' => 'float',
                                'status' => 'int',
                            ],
                        ],
                        [
                            'items_path' => 'payload.cogs',
                            'target_table' => 'cogs_gps',
                            'upsert_keys' => ['company_id', 'month'],
                            'fields' => [
                                'company_id' => 'company_id',
                                'month' => 'period',
                                'cogs' => 'metrics.cogs',
                                'gp' => 'metrics.gp',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m-d',
                                'cogs' => 'float',
                                'gp' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'payload' => [
                        'inventory' => [
                            [
                                'gl' => ['id' => 501, 'code' => 'FG-1001'],
                                'company_id' => 1,
                                'period' => '2025-01',
                                'values' => [
                                    'on_hand' => 2345000.45,
                                    'status' => 1,
                                ],
                            ],
                        ],
                        'cogs' => [
                            [
                                'company_id' => 1,
                                'period' => '2025-01-31',
                                'metrics' => [
                                    'cogs' => 950000,
                                    'gp' => 285000,
                                ],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'finance' => [
            'label' => 'Finance',
            'description' => 'Budget, expense, loan exposure and financial statements.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.summary',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{month_end}}',
                ],
                'data_mapping' => [
                    'items_path' => 'data.summary',
                    'mappings' => [
                        [
                            'target_table' => 'budget_summaries',
                            'upsert_keys' => ['month', 'category_id', 'department_id'],
                            'fields' => [
                                'month' => 'month',
                                'category_id' => 'category_id',
                                'department_id' => 'department_id',
                                'budget_amount' => 'totals.budget',
                                'actual_amount' => 'totals.actual',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m',
                                'budget_amount' => 'float',
                                'actual_amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.expense_summary',
                            'target_table' => 'budget_monthlies',
                            'upsert_keys' => ['expense_id', 'month'],
                            'fields' => [
                                'expense_id' => 'expense_id',
                                'month' => 'month',
                                'budget_amount' => 'metrics.budget',
                                'actual_amount' => 'metrics.actual',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m',
                                'budget_amount' => 'float',
                                'actual_amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.loan_status',
                            'target_table' => 'bank_loan_status_raw_data',
                            'upsert_keys' => ['company_id', 'month', 'head'],
                            'fields' => [
                                'company_id' => 'company_id',
                                'month' => 'month',
                                'head' => 'head',
                                'amount' => 'amount',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m-d',
                                'amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.financial_expense',
                            'target_table' => 'financial_expense_raw_data',
                            'upsert_keys' => ['expense_id', 'month'],
                            'fields' => [
                                'expense_id' => 'expense_id',
                                'month' => 'month',
                                'amount' => 'amount',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m-d',
                                'amount' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'summary' => [
                            [
                                'month' => '2025-01',
                                'category_id' => 4,
                                'department_id' => 15,
                                'totals' => [
                                    'budget' => 12500000,
                                    'actual' => 11850000,
                                ],
                            ],
                        ],
                        'expense_summary' => [
                            [
                                'expense_id' => 244,
                                'month' => '2025-01',
                                'metrics' => [
                                    'budget' => 2200000,
                                    'actual' => 2455000,
                                ],
                            ],
                        ],
                        'loan_status' => [
                            [
                                'company_id' => 1,
                                'month' => '2025-01-31',
                                'head' => 'Short Term Loan',
                                'amount' => 18500000,
                            ],
                        ],
                        'financial_expense' => [
                            [
                                'expense_id' => 244,
                                'month' => '2025-01-31',
                                'amount' => 425000,
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'production' => [
            'label' => 'Production',
            'description' => 'Efficiency, wastage, cost and maintenance metrics per factory/line.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'results.lines',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'items_path' => 'results.lines',
                    'mappings' => [
                        [
                            'target_table' => 'production_efficiency',
                            'upsert_keys' => ['factory_id', 'production_line_id', 'production_date'],
                            'fields' => [
                                'factory_id' => 'factory.id',
                                'production_line_id' => 'line.id',
                                'production_date' => 'date',
                                'planned_output' => 'plan',
                                'actual_output' => 'actual',
                                'efficiency_percentage' => 'metrics.efficiency',
                                'downtime_minutes' => 'metrics.downtime_minutes',
                                'oee' => 'metrics.oee',
                            ],
                            'transforms' => [
                                'production_date' => 'date:Y-m-d',
                                'planned_output' => 'float',
                                'actual_output' => 'float',
                                'efficiency_percentage' => 'float',
                                'downtime_minutes' => 'float',
                                'oee' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'results.wastage',
                            'target_table' => 'wastage_datas',
                            'upsert_keys' => ['factory', 'month'],
                            'fields' => [
                                'factory' => 'factory',
                                'month' => 'period',
                                'wastage' => 'metrics.wastage_qty',
                                'amount' => 'metrics.wastage_value',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m',
                                'wastage' => 'float',
                                'amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'results.costs',
                            'target_table' => 'cost_analyses',
                            'upsert_keys' => ['factory', 'month', 'cost_type'],
                            'fields' => [
                                'factory' => 'factory',
                                'month' => 'period',
                                'cost_type' => 'type',
                                'amount' => 'amount',
                            ],
                            'transforms' => [
                                'month' => 'date:Y-m',
                                'amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'results.maintenance',
                            'target_table' => 'machine_maintenances',
                            'upsert_keys' => ['maintenance_date', 'machine_code'],
                            'fields' => [
                                'maintenance_date' => 'date',
                                'factory_id' => 'factory.id',
                                'machine_code' => 'machine.code',
                                'machine_name' => 'machine.name',
                                'downtime_minutes' => 'metrics.downtime',
                                'events' => 'metrics.events',
                                'cost' => 'metrics.cost',
                            ],
                            'transforms' => [
                                'maintenance_date' => 'date:Y-m-d',
                                'downtime_minutes' => 'float',
                                'cost' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'results' => [
                        'lines' => [
                            [
                                'factory' => ['id' => 3],
                                'line' => ['id' => 21],
                                'date' => '2025-01-15',
                                'plan' => 1500,
                                'actual' => 1425,
                                'metrics' => [
                                    'efficiency' => 0.95,
                                    'downtime_minutes' => 120,
                                    'oee' => 0.87,
                                ],
                            ],
                        ],
                        'wastage' => [
                            [
                                'factory' => 'Plant A',
                                'period' => '2025-01',
                                'metrics' => [
                                    'wastage_qty' => 420,
                                    'wastage_value' => 18500,
                                ],
                            ],
                        ],
                        'costs' => [
                            [
                                'factory' => 'Plant A',
                                'period' => '2025-01',
                                'type' => 'energy',
                                'amount' => 72000,
                            ],
                        ],
                        'maintenance' => [
                            [
                                'date' => '2025-01-14',
                                'factory' => ['id' => 3],
                                'machine' => ['code' => 'FILL-04', 'name' => 'Filling Line 4'],
                                'metrics' => [
                                    'downtime' => 95,
                                    'events' => 2,
                                    'cost' => 12500,
                                ],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'supplychain' => [
            'label' => 'Supply Chain',
            'description' => 'Purchase orders, GRN, invoice and supplier performance feeds.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.purchase_orders',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                    'from' => '{{month_start}}',
                    'to' => '{{month_end}}',
                ],
                'data_mapping' => [
                    'items_path' => 'data.purchase_orders',
                    'mappings' => [
                        [
                            'target_table' => 'supply_chain_master_datas',
                            'upsert_keys' => ['po_number'],
                            'fields' => [
                                'po_number' => 'po_number',
                                'company' => 'company_id',
                                'supplier_id' => 'supplier.id',
                                'purchase_org' => 'purchase_org',
                                'po_date' => 'dates.po_date',
                                'delivery_date' => 'dates.delivery_date',
                                'po_value' => 'amounts.total',
                                'status' => 'status.code',
                            ],
                            'transforms' => [
                                'po_date' => 'date:Y-m-d',
                                'delivery_date' => 'date:Y-m-d',
                                'po_value' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.grn',
                            'target_table' => 'supply_chain_grn_datas',
                            'upsert_keys' => ['grn_number'],
                            'fields' => [
                                'grn_number' => 'grn_number',
                                'po_number' => 'po_number',
                                'company' => 'company_id',
                                'grn_date' => 'grn_date',
                                'grn_amount' => 'grn_amount',
                            ],
                            'transforms' => [
                                'grn_date' => 'date:Y-m-d',
                                'grn_amount' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.invoices',
                            'target_table' => 'supply_chain_invoice_datas',
                            'upsert_keys' => ['invoice_number'],
                            'fields' => [
                                'invoice_number' => 'invoice_number',
                                'po_number' => 'po_number',
                                'company' => 'company_id',
                                'iv_date' => 'invoice_date',
                                'total_invoice' => 'amounts.total',
                            ],
                            'transforms' => [
                                'iv_date' => 'date:Y-m-d',
                                'total_invoice' => 'float',
                            ],
                        ],
                        [
                            'items_path' => 'data.suppliers',
                            'target_table' => 'supplier_performance',
                            'upsert_keys' => ['supplier_name', 'evaluation_date'],
                            'fields' => [
                                'supplier_name' => 'supplier_name',
                                'company_id' => 'company_id',
                                'evaluation_date' => 'evaluation_date',
                                'overall_score' => 'scores.overall',
                                'on_time_percentage' => 'scores.on_time',
                                'quality_score' => 'scores.quality',
                                'cost_score' => 'scores.cost',
                                'rating' => 'rating',
                            ],
                            'transforms' => [
                                'evaluation_date' => 'date:Y-m-d',
                                'overall_score' => 'float',
                                'on_time_percentage' => 'float',
                                'quality_score' => 'float',
                                'cost_score' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'purchase_orders' => [
                            [
                                'po_number' => 'PO-2025-0012',
                                'company_id' => 1,
                                'supplier' => ['id' => 21],
                                'purchase_org' => 'DH-PLANT',
                                'dates' => [
                                    'po_date' => '2025-01-05',
                                    'delivery_date' => '2025-01-17',
                                ],
                                'amounts' => ['total' => 820000.0],
                                'status' => ['code' => 'open'],
                            ],
                        ],
                        'grn' => [
                            [
                                'grn_number' => 'GRN-5521',
                                'po_number' => 'PO-2025-0012',
                                'company_id' => 1,
                                'grn_date' => '2025-01-16',
                                'grn_amount' => 640000.0,
                            ],
                        ],
                        'invoices' => [
                            [
                                'invoice_number' => 'INV-8820',
                                'po_number' => 'PO-2025-0012',
                                'company_id' => 1,
                                'invoice_date' => '2025-01-20',
                                'amounts' => ['total' => 620000.0],
                            ],
                        ],
                        'suppliers' => [
                            [
                                'supplier_name' => 'Global Packaging Co.',
                                'company_id' => 1,
                                'evaluation_date' => '2025-01-31',
                                'scores' => [
                                    'overall' => 87.5,
                                    'on_time' => 82.0,
                                    'quality' => 92.0,
                                    'cost' => 80.0,
                                ],
                                'rating' => 'B+',
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'hr' => [
            'label' => 'Human Resources',
            'description' => 'Headcount, attendance, turnover and promotion feeds.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.headcount',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'items_path' => 'data.headcount',
                    'mappings' => [
                        [
                            'target_table' => 'employee_basic_infos',
                            'upsert_keys' => ['report_date'],
                            'fields' => [
                                'report_date' => 'report_date',
                                'total_active_staff' => 'metrics.active_staff',
                                'total_active_worker' => 'metrics.active_worker',
                                'total_contractual_employee' => 'metrics.contractual',
                                'total_probationary_employee' => 'metrics.probationary',
                                'total_permanent_employee' => 'metrics.permanent',
                            ],
                            'transforms' => [
                                'report_date' => 'date:Y-m-d',
                                'total_active_staff' => 'int',
                                'total_active_worker' => 'int',
                                'total_contractual_employee' => 'int',
                                'total_probationary_employee' => 'int',
                                'total_permanent_employee' => 'int',
                            ],
                        ],
                        [
                            'items_path' => 'data.attendance',
                            'target_table' => 'employee_attendances',
                            'upsert_keys' => ['date'],
                            'fields' => [
                                'date' => 'date',
                                'total_present' => 'counts.present',
                                'total_absent' => 'counts.absent',
                                'total_leave' => 'counts.leave',
                            ],
                            'transforms' => [
                                'date' => 'date:Y-m-d',
                                'total_present' => 'int',
                                'total_absent' => 'int',
                                'total_leave' => 'int',
                            ],
                        ],
                        [
                            'items_path' => 'data.turnover',
                            'target_table' => 'employee_tran_overs',
                            'upsert_keys' => ['year', 'month', 'job_type'],
                            'fields' => [
                                'year' => 'year',
                                'month' => 'month',
                                'job_type' => 'job_type',
                                'new_employee_no' => 'counts.new',
                                'resigned_employee' => 'counts.resigned',
                            ],
                            'transforms' => [
                                'year' => 'int',
                                'new_employee_no' => 'int',
                                'resigned_employee' => 'int',
                            ],
                        ],
                        [
                            'items_path' => 'data.promotions',
                            'target_table' => 'yearly_employee_promotions',
                            'upsert_keys' => ['year'],
                            'fields' => [
                                'year' => 'year',
                                'promoted_count' => 'promoted_count',
                                'details' => 'details',
                            ],
                            'transforms' => [
                                'year' => 'int',
                                'promoted_count' => 'int',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'headcount' => [
                            [
                                'report_date' => '2025-01-31',
                                'metrics' => [
                                    'active_staff' => 1840,
                                    'active_worker' => 1220,
                                    'contractual' => 210,
                                    'probationary' => 95,
                                    'permanent' => 1535,
                                ],
                            ],
                        ],
                        'attendance' => [
                            [
                                'date' => '2025-01-29',
                                'counts' => [
                                    'present' => 1760,
                                    'absent' => 55,
                                    'leave' => 25,
                                ],
                            ],
                        ],
                        'turnover' => [
                            [
                                'year' => 2025,
                                'month' => '01',
                                'job_type' => 'staff',
                                'counts' => [
                                    'new' => 12,
                                    'resigned' => 8,
                                ],
                            ],
                        ],
                        'promotions' => [
                            [
                                'year' => 2025,
                                'promoted_count' => 42,
                                'details' => 'Leadership & staff promotions published Feb-2025',
                            ],
                        ],
                    ],
                ],
            ],
        ],
    ],

    'placeholder_help' => [
        '{{company_id}}' => 'Company identifier chosen during configuration (optional).',
        '{{today}}' => 'Current date in YYYY-MM-DD.',
        '{{yesterday}}' => 'Yesterday in YYYY-MM-DD.',
        '{{month_start}}' => 'First day of current month.',
        '{{month_end}}' => 'Last day of current month.',
        '{{now_iso}}' => 'Current timestamp in ISO-8601.',
        '{{last_sync_at}}' => 'Timestamp of the previous successful sync (ISO-8601).',
    ],
];
