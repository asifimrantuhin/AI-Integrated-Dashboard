<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SupplyChainDataImport implements ToCollection, WithHeadingRow
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
            Log::error('Supply Chain import error: ' . $e->getMessage());
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
            case 'po_data':
                return [
                    'company' => $row['company'] ?? $row['company_id'] ?? null,
                    'po_number' => $row['po_number'] ?? null,
                    'po_date' => $row['po_date'] ?? null,
                    'po_value' => $row['po_value'] ?? 0,
                    'pr_amount' => $row['pr_amount'] ?? 0,
                    'purchase_org' => $row['purchase_org'] ?? null,
                ];
            case 'grn_data':
                return [
                    'company' => $row['company'] ?? $row['company_id'] ?? null,
                    'po_number' => $row['po_number'] ?? null,
                    'grn_date' => $row['grn_date'] ?? null,
                    'grn_amount' => $row['grn_amount'] ?? 0,
                ];
            case 'invoice_data':
                return [
                    'company' => $row['company'] ?? $row['company_id'] ?? null,
                    'iv_date' => $row['iv_date'] ?? $row['invoice_date'] ?? null,
                    'total_invoice' => $row['total_invoice'] ?? $row['amount'] ?? 0,
                ];
            default:
                return null;
        }
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'po_data':
                    DB::table('supply_chain_master_datas')->insert($data);
                    break;
                case 'grn_data':
                    DB::table('supply_chain_grn_datas')->insert($data);
                    break;
                case 'invoice_data':
                    DB::table('supply_chain_invoice_datas')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

