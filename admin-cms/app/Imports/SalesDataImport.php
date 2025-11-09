<?php

namespace App\Imports;

use App\Models\FileUpload;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Carbon;
use Illuminate\Support\Str;
use Maatwebsite\Excel\Concerns\ToCollection;
use Maatwebsite\Excel\Concerns\WithHeadingRow;

class SalesDataImport implements ToCollection, WithHeadingRow
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

                if ($mappedData) {
                    $this->insertData($mappedData);
                    $recordsProcessed++;
                }
            }

            $this->upload->update([
                'records_processed' => $recordsProcessed,
                'status' => 'completed'
            ]);
        } catch (\Throwable $e) {
            Log::error('Sales import error: ' . $e->getMessage(), [
                'upload_id' => $this->upload->id,
                'type' => $this->dataType,
            ]);

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
            case 'channel_targets':
                return $this->mapChannelTargets($row);
            case 'order_book':
                return $this->mapOrderBook($row);
            case 'promotion_performance':
                return $this->mapPromotionPerformance($row);
            default:
                return null;
        }
    }

    private function mapMonthlyReport(array $row): ?array
    {
        return [
            'data_month' => $this->parseDate($row['data_month'] ?? $row['month'] ?? null, 'Y-m-d'),
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'lifting_target' => $this->toDecimal($row['lifting_target'] ?? 0),
            'billed' => $this->toDecimal($row['billed'] ?? 0),
            'delivered' => $this->toDecimal($row['delivered'] ?? 0),
            'primary_collection' => $this->toDecimal($row['primary_collection'] ?? 0),
            'ims_target' => $this->toDecimal($row['ims_target'] ?? 0),
            'ims' => $this->toDecimal($row['ims'] ?? 0),
            'market_collection' => $this->toDecimal($row['market_collection'] ?? 0),
        ];
    }

    private function mapDailyReport(array $row): ?array
    {
        return [
            'data_date' => $this->parseDate($row['data_date'] ?? $row['date'] ?? null),
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'billed' => $this->toDecimal($row['billed'] ?? 0),
            'delivery' => $this->toDecimal($row['delivery'] ?? 0),
            'ims' => $this->toDecimal($row['ims'] ?? 0),
        ];
    }

    private function mapBestSelling(array $row): ?array
    {
        return [
            'year_month' => $this->parseDate($row['year_month'] ?? $row['month'] ?? null, 'Y-m-d'),
            'channel_id' => $row['channel_id'] ?? null,
            'product_id' => $row['product_id'] ?? null,
            'product_name' => $row['product_name'] ?? null,
            'qty' => $this->toDecimal($row['qty'] ?? $row['quantity'] ?? 0),
            'value' => $this->toDecimal($row['value'] ?? $row['amount'] ?? 0),
        ];
    }

    private function mapTopDistributors(array $row): ?array
    {
        return [
            'db_name' => $row['db_name'] ?? $row['distributor_name'] ?? null,
            'amount' => $this->toDecimal($row['amount'] ?? $row['value'] ?? 0),
            'type' => 0,
            'date' => $this->parseDate($row['date'] ?? null),
        ];
    }

    private function mapTopRetailers(array $row): ?array
    {
        return [
            'db_name' => $row['db_name'] ?? $row['retailer_name'] ?? null,
            'amount' => $this->toDecimal($row['amount'] ?? $row['value'] ?? 0),
            'type' => 1,
            'date' => $this->parseDate($row['date'] ?? null),
        ];
    }

    private function mapOrderDelivery(array $row): ?array
    {
        return [
            'months' => $row['months'] ?? $row['month'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'amounts' => $this->toDecimal($row['amounts'] ?? $row['amount'] ?? 0),
            'types' => $row['types'] ?? $row['type'] ?? 0,
        ];
    }

    private function mapChannelTargets(array $row): ?array
    {
        return [
            'data_month' => $this->parseDate($row['data_month'] ?? $row['month'] ?? null, 'Y-m-d'),
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'revenue_target' => $this->toDecimal($row['revenue_target'] ?? $row['revenue'] ?? 0),
            'volume_target' => $this->toDecimal($row['volume_target'] ?? $row['volume'] ?? 0),
            'promotion_budget' => $this->toDecimal($row['promotion_budget'] ?? $row['promo_budget'] ?? 0),
            'gross_margin_target' => $this->toDecimal($row['gross_margin_target'] ?? $row['gm_target'] ?? 0),
            'new_customer_target' => $this->toDecimal($row['new_customer_target'] ?? $row['new_customers'] ?? 0),
            'owner' => $row['owner'] ?? $row['owner_name'] ?? null,
            'created_at' => now(),
            'updated_at' => now(),
        ];
    }

    private function mapOrderBook(array $row): ?array
    {
        return [
            'order_number' => $row['order_number'] ?? $row['order_id'] ?? Str::uuid()->toString(),
            'order_date' => $this->parseDate($row['order_date'] ?? $row['date'] ?? null),
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'customer_code' => $row['customer_code'] ?? $row['customer_id'] ?? null,
            'customer_name' => $row['customer_name'] ?? null,
            'region' => $row['region'] ?? $row['territory'] ?? null,
            'status' => $this->normaliseStatus($row['status'] ?? 'confirmed'),
            'order_amount' => $this->toDecimal($row['order_amount'] ?? $row['amount'] ?? 0),
            'fulfilled_at' => $this->parseDate($row['fulfilled_at'] ?? $row['delivery_date'] ?? null),
            'discount_amount' => $this->toDecimal($row['discount_amount'] ?? $row['discount'] ?? 0),
            'gross_margin' => $this->toDecimal($row['gross_margin'] ?? $row['margin'] ?? 0),
            'created_at' => now(),
            'updated_at' => now(),
        ];
    }

    private function mapPromotionPerformance(array $row): ?array
    {
        return [
            'campaign_code' => $row['campaign_code'] ?? $row['campaign_id'] ?? Str::uuid()->toString(),
            'campaign_name' => $row['campaign_name'] ?? null,
            'channel_id' => $row['channel_id'] ?? null,
            'channel_name' => $row['channel_name'] ?? null,
            'start_date' => $this->parseDate($row['start_date'] ?? null),
            'end_date' => $this->parseDate($row['end_date'] ?? null),
            'spend_amount' => $this->toDecimal($row['spend_amount'] ?? $row['spend'] ?? 0),
            'revenue_uplift' => $this->toDecimal($row['revenue_uplift'] ?? $row['uplift_revenue'] ?? 0),
            'uplift_percentage' => $this->toDecimal($row['uplift_percentage'] ?? $row['uplift_pct'] ?? 0),
            'roi' => $this->toDecimal($row['roi'] ?? 0),
            'audience_tags' => $this->normaliseTags($row['audience_tags'] ?? $row['audience'] ?? null),
            'created_at' => now(),
            'updated_at' => now(),
        ];
    }

    private function insertData(array $data): void
    {
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
            case 'channel_targets':
                DB::table('sales_channel_targets')->insert($data);
                break;
            case 'order_book':
                DB::table('sales_order_book')->insert($data);
                break;
            case 'promotion_performance':
                DB::table('sales_promotion_performance')->insert($data);
                break;
        }
    }

    private function parseDate($value, string $format = 'Y-m-d')
    {
        if (empty($value)) {
            return null;
        }

        try {
            return Carbon::parse($value)->format($format);
        } catch (\Throwable $e) {
            return null;
        }
    }

    private function toDecimal($value): float
    {
        if ($value === null || $value === '') {
            return 0.0;
        }

        return (float) str_replace([','], '', $value);
    }

    private function normaliseStatus(string $status): string
    {
        $status = Str::of($status)->lower()->value();

        return match ($status) {
            'draft', 'pending' => 'draft',
            'confirmed', 'booked' => 'confirmed',
            'dispatch', 'dispatching', 'shipped' => 'dispatching',
            'delivered', 'completed' => 'delivered',
            'cancel', 'cancelled', 'void' => 'cancelled',
            default => 'confirmed',
        };
    }

    private function normaliseTags($value): ?string
    {
        if (empty($value)) {
            return null;
        }

        if (is_array($value)) {
            return json_encode(array_values(array_filter($value)));
        }

        $tags = array_filter(array_map('trim', preg_split('/[,;|]/', (string) $value)));

        return empty($tags) ? null : json_encode(array_values($tags));
    }
}

