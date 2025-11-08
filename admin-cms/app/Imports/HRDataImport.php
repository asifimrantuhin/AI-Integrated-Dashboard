<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class HRDataImport implements ToCollection, WithHeadingRow
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
            Log::error('HR import error: ' . $e->getMessage());
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
            case 'employee_basic_info':
                return [
                    'employee_id' => $row['employee_id'] ?? null,
                    'name' => $row['name'] ?? null,
                    'department' => $row['department'] ?? null,
                    'company_id' => $row['company_id'] ?? null,
                    'status' => $row['status'] ?? 1,
                ];
            case 'employee_attendance':
                return [
                    'employee_id' => $row['employee_id'] ?? null,
                    'date' => $row['date'] ?? null,
                    'status' => $row['status'] ?? 'present',
                ];
            case 'employee_promotions':
                return [
                    'employee_id' => $row['employee_id'] ?? null,
                    'year' => $row['year'] ?? date('Y'),
                    'promotion_date' => $row['promotion_date'] ?? null,
                    'new_designation' => $row['new_designation'] ?? null,
                ];
            case 'employee_tran_over':
                return [
                    'employee_id' => $row['employee_id'] ?? null,
                    'date' => $row['date'] ?? null,
                    'type' => $row['type'] ?? null,
                    'amount' => $row['amount'] ?? 0,
                ];
            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'employee_basic_info':
                    DB::table('employee_basic_infos')->insert($data);
                    break;
                case 'employee_attendance':
                    DB::table('employee_attendances')->insert($data);
                    break;
                case 'employee_promotions':
                    DB::table('yearly_employee_promotions')->insert($data);
                    break;
                case 'employee_tran_over':
                    DB::table('employee_tran_overs')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

