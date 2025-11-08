<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class NPDDataImport implements ToCollection, WithHeadingRow
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
            Log::error('NPD import error: ' . $e->getMessage());
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
            case 'npd_projects':
                return [
                    'project_name' => $row['project_name'] ?? null,
                    'start_date' => $row['start_date'] ?? null,
                    'end_date' => $row['end_date'] ?? null,
                    'status' => $row['status'] ?? 'pending',
                    'company_id' => $row['company_id'] ?? null,
                ];
            case 'project_deliverables':
                return [
                    'project_id' => $row['project_id'] ?? null,
                    'name' => $row['name'] ?? null,
                    'due_date' => $row['due_date'] ?? null,
                    'status' => $row['status'] ?? 'pending',
                ];
            case 'project_sub_deliverables':
                return [
                    'deliverable_id' => $row['deliverable_id'] ?? null,
                    'name' => $row['name'] ?? null,
                    'due_date' => $row['due_date'] ?? null,
                    'status' => $row['status'] ?? 'pending',
                ];
            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'npd_projects':
                    DB::table('npd_projects')->insert($data);
                    break;
                case 'project_deliverables':
                    DB::table('projects_deliverables')->insert($data);
                    break;
                case 'project_sub_deliverables':
                    DB::table('projects_sub_deliverables')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

