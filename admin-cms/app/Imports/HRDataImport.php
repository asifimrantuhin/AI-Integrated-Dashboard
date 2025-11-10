<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;

class HRDataImport implements ToCollection, WithHeadingRow
{
    protected FileUpload $upload;
    protected string $dataType;

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

                if ($mappedData !== null) {
                    $this->insertData($mappedData);
                    $recordsProcessed++;
                }
            }

            $this->upload->update([
                'records_processed' => $recordsProcessed,
                'status' => 'completed',
            ]);
        } catch (\Throwable $e) {
            Log::error('HR import error: ' . $e->getMessage(), [
                'upload_id' => $this->upload->id,
                'type' => $this->dataType,
            ]);

            $this->upload->update([
                'status' => 'failed',
                'error_message' => $e->getMessage(),
            ]);

            throw $e;
        }
    }

    private function mapData(array $row): ?array
    {
        switch ($this->dataType) {
            case 'headcount_snapshot':
            case 'employee_basic_info':
                return [
                    'report_date' => $row['report_date'] ?? $row['date'] ?? now()->toDateString(),
                    'total_active_staff' => $row['total_active_staff'] ?? $row['active_staff'] ?? $row['active_total'] ?? 0,
                    'total_active_worker' => $row['total_active_worker'] ?? $row['active_worker'] ?? 0,
                    'total_contractual_employee' => $row['total_contractual_employee'] ?? $row['contractual'] ?? 0,
                    'total_probationary_employee' => $row['total_probationary_employee'] ?? $row['probationary'] ?? 0,
                    'total_permanent_employee' => $row['total_permanent_employee'] ?? $row['permanent'] ?? 0,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'attendance_summary':
            case 'employee_attendance':
                return [
                    'date' => $row['date'] ?? now()->toDateString(),
                    'total_present' => $row['total_present'] ?? $row['present'] ?? 0,
                    'total_absent' => $row['total_absent'] ?? $row['absent'] ?? 0,
                    'total_leave' => $row['total_leave'] ?? $row['leave'] ?? $row['leave_count'] ?? 0,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'movement_summary':
            case 'employee_tran_over':
                $month = $row['month'] ?? (isset($row['date']) ? date('m', strtotime($row['date'])) : now()->format('m'));
                $year = $row['year'] ?? (isset($row['date']) ? date('Y', strtotime($row['date'])) : now()->format('Y'));

                return [
                    'year' => (int) $year,
                    'month' => (string) $month,
                    'job_type' => $row['job_type'] ?? $row['type'] ?? 'staff',
                    'new_employee_no' => $row['new_employee_no'] ?? $row['new'] ?? $row['new_hires'] ?? 0,
                    'resigned_employee' => $row['resigned_employee'] ?? $row['resigned'] ?? $row['attrition'] ?? 0,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'promotion_summary':
            case 'employee_promotions':
                return [
                    'year' => (int) ($row['year'] ?? date('Y')), 
                    'promoted_count' => $row['promoted_count'] ?? $row['count'] ?? 0,
                    'details' => $row['details'] ?? $row['note'] ?? null,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        switch ($this->dataType) {
            case 'headcount_snapshot':
            case 'employee_basic_info':
                DB::table('employee_basic_infos')->insert($data);
                break;

            case 'attendance_summary':
            case 'employee_attendance':
                DB::table('employee_attendances')->insert($data);
                break;

            case 'movement_summary':
            case 'employee_tran_over':
                DB::table('employee_tran_overs')->insert($data);
                break;

            case 'promotion_summary':
            case 'employee_promotions':
                DB::table('yearly_employee_promotions')->insert($data);
                break;
        }
    }
}

