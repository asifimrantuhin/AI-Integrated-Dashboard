<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class ProductionDataImport implements ToCollection, WithHeadingRow
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
            Log::error('Production import error: ' . $e->getMessage());
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
            case 'production_analysis':
                return $this->mapProductionAnalysis($row);
            case 'wastage_data':
                return $this->mapWastageData($row);
            case 'cost_analysis':
                return $this->mapCostAnalysis($row);
            case 'line_efficiency':
                return $this->mapLineEfficiency($row);
            case 'maintenance_events':
                return $this->mapMaintenanceEvent($row);
            default:
                return null;
        }
    }

    private function mapProductionAnalysis(array $row): ?array
    {
        return [
            'month' => $row['month'] ?? null,
            'factory' => $row['factory'] ?? null,
            'summary_group' => $row['summary_group'] ?? null,
            'amonthly_amount' => $row['amonthly_amount'] ?? 0,
            'amonthly_per' => $row['amonthly_per'] ?? 0,
            'cmonthly_amount' => $row['cmonthly_amount'] ?? 0,
            'cmonthly_per' => $row['cmonthly_per'] ?? 0,
        ];
    }

    private function mapWastageData(array $row): ?array
    {
        return [
            'month' => $row['month'] ?? null,
            'factory' => $row['factory'] ?? null,
            'wastage' => $row['wastage'] ?? 0,
            'amount' => $row['amount'] ?? 0,
        ];
    }

    private function mapCostAnalysis(array $row): ?array
    {
        return [
            'month' => $row['month'] ?? null,
            'factory' => $row['factory'] ?? null,
            'cost_type' => $row['cost_type'] ?? null,
            'amount' => $row['amount'] ?? 0,
        ];
    }

    private function mapLineEfficiency(array $row): ?array
    {
        return [
            'production_date' => $row['production_date'] ?? $row['date'] ?? null,
            'factory_id' => $row['factory_id'] ?? $row['factory'] ?? null,
            'production_line_id' => $row['production_line_id'] ?? $row['line_id'] ?? null,
            'actual_output' => $row['actual_output'] ?? $row['actual'] ?? 0,
            'planned_output' => $row['planned_output'] ?? $row['planned'] ?? 0,
            'efficiency_percentage' => $row['efficiency_percentage'] ?? $row['efficiency'] ?? 0,
            'downtime_minutes' => $row['downtime_minutes'] ?? $row['downtime'] ?? 0,
            'oee' => $row['oee'] ?? $row['oee_percentage'] ?? 0,
        ];
    }

    private function mapMaintenanceEvent(array $row): ?array
    {
        return [
            'maintenance_date' => $row['maintenance_date'] ?? $row['date'] ?? null,
            'factory_id' => $row['factory_id'] ?? $row['factory'] ?? null,
            'machine_code' => $row['machine_code'] ?? null,
            'machine_name' => $row['machine_name'] ?? null,
            'downtime_minutes' => $row['downtime_minutes'] ?? $row['downtime'] ?? 0,
            'cost' => $row['cost'] ?? 0,
            'description' => $row['description'] ?? $row['remarks'] ?? null,
        ];
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'production_analysis':
                    DB::table('production_analyses')->insert($data);
                    break;
                case 'wastage_data':
                    DB::table('wastage_datas')->insert($data);
                    break;
                case 'cost_analysis':
                    DB::table('cost_analyses')->insert($data);
                    break;
                case 'line_efficiency':
                    DB::table('production_efficiency')->insert($data);
                    break;
                case 'maintenance_events':
                    DB::table('machine_maintenances')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

