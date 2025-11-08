<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SalesDataImport implements ToCollection, WithHeadingRow
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
            Log::error('Sales import error: ' . $e->getMessage());
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
            case 'monthly_report':
                return $this->mapMonthlyReport($row);
            case 'daily_report':
                return $this->mapDailyReport($row);
            case 'best_selling':
                return $this->mapBestSelling($row);
            case 'top_distributors':
                return $this->mapTopDistributors($row);
            case 'top_retailers':
                return $this->mapTopRetailers($row);
            case 'order_delivery':
                return $this->mapOrderDelivery($row);
            default:
                return null;
        }
    }

    private function mapMonthlyReport(array $row): ?array
    {
        // Map monthly report data
        return [
            'data_month' => $row['data_month'] ?? $row['month'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'lifting_target' => $row['lifting_target'] ?? 0,
            'billed' => $row['billed'] ?? 0,
            'delivered' => $row['delivered'] ?? 0,
            'primary_collection' => $row['primary_collection'] ?? 0,
            'ims_target' => $row['ims_target'] ?? 0,
            'ims' => $row['ims'] ?? 0,
            'market_collection' => $row['market_collection'] ?? 0,
        ];
    }

    private function mapDailyReport(array $row): ?array
    {
        // Map daily report data
        return [
            'data_date' => $row['data_date'] ?? $row['date'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'billed' => $row['billed'] ?? 0,
            'delivery' => $row['delivery'] ?? 0,
            'ims' => $row['ims'] ?? 0,
        ];
    }

    private function mapBestSelling(array $row): ?array
    {
        // Map best selling products data
        return [
            'year_month' => $row['year_month'] ?? $row['month'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'product_id' => $row['product_id'] ?? null,
            'product_name' => $row['product_name'] ?? null,
            'qty' => $row['qty'] ?? $row['quantity'] ?? 0,
            'value' => $row['value'] ?? $row['amount'] ?? 0,
        ];
    }

    private function mapTopDistributors(array $row): ?array
    {
        // Map top distributors data
        return [
            'db_name' => $row['db_name'] ?? $row['distributor_name'] ?? null,
            'amount' => $row['amount'] ?? $row['value'] ?? 0,
            'type' => 0, // 0 for distributor
            'date' => $row['date'] ?? now(),
        ];
    }

    private function mapTopRetailers(array $row): ?array
    {
        // Map top retailers data
        return [
            'db_name' => $row['db_name'] ?? $row['retailer_name'] ?? null,
            'amount' => $row['amount'] ?? $row['value'] ?? 0,
            'type' => 1, // 1 for retailer
            'date' => $row['date'] ?? now(),
        ];
    }

    private function mapOrderDelivery(array $row): ?array
    {
        // Map order vs delivery data
        return [
            'months' => $row['months'] ?? $row['month'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'amounts' => $row['amounts'] ?? $row['amount'] ?? 0,
            'types' => $row['types'] ?? $row['type'] ?? 0,
        ];
    }

    private function insertData(array $data): void
    {
        try {
            switch ($this->dataType) {
                case 'monthly_report':
                    DB::table('channelwise_monthly_report')->insert($data);
                    break;
                case 'daily_report':
                    DB::table('channelwise_lic_data')->insert($data);
                    break;
                case 'best_selling':
                    DB::table('best_selling_products')->insert($data);
                    break;
                case 'top_distributors':
                case 'top_retailers':
                    DB::table('top_channel_d_bs')->insert($data);
                    break;
                case 'order_delivery':
                    DB::table('order_delivery_summaries')->insert($data);
                    break;
            }
        } catch (\Exception $e) {
            Log::error('Data insert error: ' . $e->getMessage());
            throw $e;
        }
    }
}

