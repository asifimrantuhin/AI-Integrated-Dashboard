<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;

class SupplyChainDataImport implements ToCollection, WithHeadingRow
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
            Log::error('Supply Chain import error: ' . $e->getMessage(), [
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
            case 'po_data':
                return [
                    'company' => (string) ($row['company'] ?? $row['company_id'] ?? ''),
                    'po_number' => (string) ($row['po_number'] ?? $row['po'] ?? ''),
                    'po_date' => $row['po_date'] ?? $row['date'] ?? null,
                    'po_value' => (float) ($row['po_value'] ?? $row['po_amount'] ?? $row['amount'] ?? 0),
                    'pr_amount' => (float) ($row['pr_amount'] ?? $row['pr_value'] ?? 0),
                    'purchase_org' => $row['purchase_org'] ?? $row['purchase_organization'] ?? null,
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'grn_data':
                return [
                    'company' => (string) ($row['company'] ?? $row['company_id'] ?? ''),
                    'po_number' => (string) ($row['po_number'] ?? $row['po'] ?? ''),
                    'grn_date' => $row['grn_date'] ?? $row['date'] ?? null,
                    'grn_amount' => (float) ($row['grn_amount'] ?? $row['amount'] ?? 0),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'invoice_data':
                return [
                    'company' => (string) ($row['company'] ?? $row['company_id'] ?? ''),
                    'invoice_number' => (string) ($row['invoice_number'] ?? $row['invoice'] ?? ''),
                    'iv_date' => $row['iv_date'] ?? $row['invoice_date'] ?? $row['date'] ?? null,
                    'total_invoice' => (float) ($row['total_invoice'] ?? $row['amount'] ?? 0),
                    'created_at' => now(),
                    'updated_at' => now(),
                ];

            case 'supplier_performance':
                $overall = (float) ($row['overall_score'] ?? $row['overall'] ?? 0);
                $rating = $row['rating'] ?? ($overall >= 90 ? 'excellent' : ($overall >= 75 ? 'good' : ($overall >= 60 ? 'average' : 'poor')));

                return [
                    'company_id' => (int) ($row['company_id'] ?? $row['company'] ?? 0),
                    'supplier_code' => (string) ($row['supplier_code'] ?? $row['supplier_id'] ?? ''),
                    'supplier_name' => $row['supplier_name'] ?? 'Unknown supplier',
                    'evaluation_date' => $row['evaluation_date'] ?? $row['date'] ?? now()->toDateString(),
                    'total_orders' => (int) ($row['total_orders'] ?? $row['orders'] ?? 0),
                    'on_time_deliveries' => (int) ($row['on_time_deliveries'] ?? $row['on_time'] ?? 0),
                    'on_time_percentage' => (float) ($row['on_time_percentage'] ?? $row['on_time_pct'] ?? 0),
                    'quality_issues' => (int) ($row['quality_issues'] ?? $row['defects'] ?? 0),
                    'quality_score' => (float) ($row['quality_score'] ?? $row['quality'] ?? 0),
                    'cost_score' => (float) ($row['cost_score'] ?? $row['cost'] ?? 0),
                    'overall_score' => $overall,
                    'rating' => $rating,
                    'comments' => $row['comments'] ?? null,
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
            case 'po_data':
                DB::table('supply_chain_master_datas')->insert($data);
                break;

            case 'grn_data':
                DB::table('supply_chain_grn_datas')->insert($data);
                break;

            case 'invoice_data':
                DB::table('supply_chain_invoice_datas')->insert($data);
                break;

            case 'supplier_performance':
                DB::table('supplier_performance')->insert($data);
                break;
        }
    }
}

