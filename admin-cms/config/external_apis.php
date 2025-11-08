<?php

return [
    'modules' => [
        'sales' => [
            'label' => 'Sales',
            'description' => 'Channel or product level sales metrics aggregated by period.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.records',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'mappings' => [
                        [
                            'target_table' => 'channelwise_monthly_report',
                            'upsert_keys' => ['channel_id', 'data_month'],
                            'fields' => [
                                'channel_id' => 'channel_id',
                                'data_month' => 'period',
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
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'records' => [
                            [
                                'channel_id' => 12,
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
                    ],
                ],
            ],
        ],
        'inventory' => [
            'label' => 'Inventory',
            'description' => 'Inventory valuation and turnover metrics per GL or SKU.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'payload.items',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                    'as_of' => '{{today}}',
                ],
                'data_mapping' => [
                    'mappings' => [
                        [
                            'target_table' => 'inventory_raw_datas',
                            'upsert_keys' => ['gl_id', 'month'],
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
                    ],
                ],
                'sample_response' => [
                    'payload' => [
                        'items' => [
                            [
                                'gl' => [
                                    'id' => 501,
                                    'code' => 'FG-1001',
                                ],
                                'company_id' => 1,
                                'period' => '2025-01',
                                'values' => [
                                    'on_hand' => 2345000.45,
                                    'status' => 1,
                                ],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'finance' => [
            'label' => 'Finance',
            'description' => 'Budget vs actual summaries by category and department.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.rows',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{month_end}}',
                ],
                'data_mapping' => [
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
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'rows' => [
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
                    ],
                ],
            ],
        ],
        'production' => [
            'label' => 'Production',
            'description' => 'Production efficiency and cost metrics per factory.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'results.entries',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'mappings' => [
                        [
                            'target_table' => 'production_efficiency',
                            'upsert_keys' => ['factory_id', 'production_date'],
                            'fields' => [
                                'factory_id' => 'factory.id',
                                'production_date' => 'date',
                                'planned_output' => 'plan.planned_output',
                                'actual_output' => 'actual.actual_output',
                                'utilization' => 'metrics.utilization',
                            ],
                            'transforms' => [
                                'production_date' => 'date:Y-m-d',
                                'planned_output' => 'float',
                                'actual_output' => 'float',
                                'utilization' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'results' => [
                        'entries' => [
                            [
                                'factory' => ['id' => 3],
                                'date' => '2025-01-15',
                                'plan' => ['planned_output' => 1500],
                                'actual' => ['actual_output' => 1425],
                                'metrics' => ['utilization' => 0.95],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'supplychain' => [
            'label' => 'Supply Chain',
            'description' => 'Open purchase orders, GRN, Invoice cycles.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.purchase_orders',
                'query_params' => [
                    'company_id' => '{{company_id}}',
                ],
                'data_mapping' => [
                    'mappings' => [
                        [
                            'target_table' => 'purchase_orders',
                            'upsert_keys' => ['po_number'],
                            'fields' => [
                                'po_number' => 'po_number',
                                'supplier_id' => 'supplier.id',
                                'po_date' => 'dates.po_date',
                                'delivery_date' => 'dates.delivery_date',
                                'amount' => 'amounts.total',
                                'status' => 'status.code',
                            ],
                            'transforms' => [
                                'po_date' => 'date:Y-m-d',
                                'delivery_date' => 'date:Y-m-d',
                                'amount' => 'float',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'purchase_orders' => [
                            [
                                'po_number' => 'PO-2025-0012',
                                'supplier' => ['id' => 21],
                                'dates' => [
                                    'po_date' => '2025-01-05',
                                    'delivery_date' => '2025-01-17',
                                ],
                                'amounts' => [
                                    'total' => 820000.0,
                                ],
                                'status' => ['code' => 'open'],
                            ],
                        ],
                    ],
                ],
            ],
        ],
        'hr' => [
            'label' => 'Human Resources',
            'description' => 'Headcount, attrition, attendance and promotion records.',
            'defaults' => [
                'method' => 'GET',
                'items_path' => 'data.employees',
                'query_params' => [
                    'from' => '{{month_start}}',
                    'to' => '{{today}}',
                ],
                'data_mapping' => [
                    'mappings' => [
                        [
                            'target_table' => 'hr_employee_movements',
                            'upsert_keys' => ['employee_id', 'movement_date'],
                            'fields' => [
                                'employee_id' => 'employee_id',
                                'movement_date' => 'movement.date',
                                'movement_type' => 'movement.type',
                                'department_id' => 'department.id',
                                'designation' => 'designation',
                            ],
                            'transforms' => [
                                'movement_date' => 'date:Y-m-d',
                            ],
                        ],
                    ],
                ],
                'sample_response' => [
                    'data' => [
                        'employees' => [
                            [
                                'employee_id' => 4412,
                                'movement' => ['date' => '2025-01-20', 'type' => 'promotion'],
                                'department' => ['id' => 9],
                                'designation' => 'Senior Analyst',
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
