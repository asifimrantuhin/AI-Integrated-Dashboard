<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class FinanceDataImport implements ToCollection, WithHeadingRow
{
    protected $upload;
    protected $dataType;

    public function __construct(FileUpload $upload, string $dataType)
    {
        $this->upload = $upload;
        $this->dataType = $dataType;
    }

    public function collection(Collection $rows)
    {
        $recordsProcessed = 0;

        try {
            foreach ($rows as $row) {
                $mappedData = $this->mapData($row->toArray());
                
                if ($mappedData) {
                    $this->insertData($mappedData);
                    $recordsProcessed++;
                }
            }

            $this->upload->update([
                'records_processed' => $recordsProcessed,
                'status' => 'completed'
            ]);
        } catch (\Exception $e) {
            Log::error('Finance import error: ' . $e->getMessage());
            $this->upload->update([
                'status' => 'failed',
                'error_message' => $e->getMessage()
            ]);
            throw $e;
        }
    }

    private function mapData(array $row): ?array
    {
        switch ($this->dataType) {
            case 'bank_loan':
                return [
                    'year_month' => $row['year_month'] ?? null,
                    'company_code' => $row['company_code'] ?? null,
                    'net_sales' => $row['net_sales'] ?? 0,
                    'financial_expense' => $row['financial_expense'] ?? 0,
                ];
            case 'bank_loan_status':
                return [
                    'company_id' => $row['company_id'] ?? null,
                    'month' => $row['month'] ?? null,
                    'head' => $row['head'] ?? null,
                    'amount' => $row['amount'] ?? 0,
                ];
            case 'budget':
                return [
                    'month' => $row['month'] ?? null,
                    'category_id' => $row['category_id'] ?? null,
                    'department_id' => $row['department_id'] ?? null,
                    'amount' => $row['amount'] ?? 0,
                ];
            case 'expense':
                return [
                    'budget_id' => $row['budget_id'] ?? null,
                    'month' => $row['month'] ?? null,
                    'amount' => $row['amount'] ?? 0,
                ];
            case 'budget_summary':
                return [
                    'month' => $row['month'] ?? null,
                    'category_id' => $row['category_id'] ?? null,
                    'department_id' => $row['department_id'] ?? null,
                    'budget_amount' => $row['budget_amount'] ?? $row['budget'] ?? 0,
                    'actual_amount' => $row['actual_amount'] ?? $row['actual'] ?? 0,
                ];
            case 'expense_summary':
                return [
                    'month' => $row['month'] ?? null,
                    'expense_id' => $row['expense_id'] ?? $row['budget_id'] ?? null,
                    'budget_amount' => $row['budget_amount'] ?? $row['budget'] ?? 0,
                    'actual_amount' => $row['actual_amount'] ?? $row['actual'] ?? 0,
                ];
            case 'financial_statement':
                return [
                    'month' => $row['month'] ?? null,
                    'expense_id' => $row['expense_id'] ?? null,
                    'amount' => $row['amount'] ?? $row['actual_amount'] ?? 0,
                ];
            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'bank_loan':
                    DB::table('bank_loan')->insert($data);
                    break;
                case 'bank_loan_status':
                    DB::table('bank_loan_status_raw_data')->insert($data);
                    break;
                case 'budget':
                    DB::table('bdgt_budgets')->insert($data);
                    break;
                case 'expense':
                    DB::table('bdgt_expenses')->insert($data);
                    break;
                case 'budget_summary':
                    DB::table('budget_summaries')->insert($data);
                    break;
                case 'expense_summary':
                    DB::table('budget_monthlies')->insert($data);
                    break;
                case 'financial_statement':
                    DB::table('financial_expense_raw_data')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

