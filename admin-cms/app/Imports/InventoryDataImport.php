<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class InventoryDataImport implements ToCollection, WithHeadingRow
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
            Log::error('Inventory import error: ' . $e->getMessage());
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
            case 'inventory_raw_data':
                return [
                    'company_id' => $row['company_id'] ?? null,
                    'gl_id' => $row['gl_id'] ?? null,
                    'month' => $row['month'] ?? null,
                    'amount' => $row['amount'] ?? 0,
                    'status' => $row['status'] ?? 1,
                ];
            case 'cogs_gp':
                return [
                    'company_id' => $row['company_id'] ?? null,
                    'month' => $row['month'] ?? null,
                    'cogs' => $row['cogs'] ?? 0,
                    'gp' => $row['gp'] ?? 0,
                ];
            case 'inventory_gl_accounts':
                return [
                    'gl_code' => $row['gl_code'] ?? null,
                    'gl_name' => $row['gl_name'] ?? null,
                    'status' => $row['status'] ?? 1,
                ];
            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'inventory_raw_data':
                    DB::table('inventory_raw_datas')->insert($data);
                    break;
                case 'cogs_gp':
                    DB::table('cogs_gps')->insert($data);
                    break;
                case 'inventory_gl_accounts':
                    DB::table('inventory_gl_accounts')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

