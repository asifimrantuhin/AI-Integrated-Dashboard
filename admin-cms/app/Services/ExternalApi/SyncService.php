<?php

namespace App\Services\ExternalApi;

use App\Models\ApiSync;
use App\Models\ExternalApi;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Str;

class SyncService
{
    public function __construct(private readonly DataMapper $mapper)
    {
    }

    public function test(ExternalApi $api): array
    {
        return $this->fetchPayload($api);
    }

    public function sync(ExternalApi $api): ApiSync
    {
        $sync = ApiSync::create([
            'external_api_id' => $api->id,
            'status' => 'processing',
            'started_at' => now(),
        ]);

        try {
            $payload = $this->fetchPayload($api);
            $records = $this->mapper->mapAndPersist($api, $payload['data']);

            $sync->update([
                'status' => 'completed',
                'completed_at' => now(),
                'records_synced' => $records,
                'response_data' => $this->summariseResponse($payload['raw']),
            ]);

            $api->last_sync_at = now();
            $api->save();

            return $sync;
        } catch (\Throwable $e) {
            $sync->update([
                'status' => 'failed',
                'completed_at' => now(),
                'error_message' => $e->getMessage(),
            ]);

            throw $e;
        }
    }

    private function fetchPayload(ExternalApi $api): array
    {
        $headers = $this->buildHeaders($api);
        $query = $this->resolvePlaceholders($api, $api->query_params ?? []);
        $body = $this->resolvePlaceholders($api, $api->body ?? []);

        $options = [];
        if (!empty($query)) {
            $options['query'] = $query;
        }
        if (!empty($body) && in_array($api->method, ['POST', 'PUT', 'PATCH'])) {
            $options['json'] = $body;
        }

        $response = Http::withHeaders($headers)->send($api->method ?? 'GET', $api->url, $options);

        if (!$response->successful()) {
            Log::warning('External API non-success response', [
                'api_id' => $api->id,
                'status' => $response->status(),
                'body' => $response->body(),
            ]);
            $response->throw();
        }

        $raw = $this->parseResponse($api, $response->body());
        $data = $this->extractData($raw, $api->items_path);

        return [
            'raw' => $raw,
            'data' => $data,
        ];
    }

    private function buildHeaders(ExternalApi $api): array
    {
        $headers = $api->headers ?? [];
        $auth = $api->authentication ?? [];
        $type = $auth['type'] ?? ($auth['auth_type'] ?? null);

        if ($type) {
            $token = $this->resolvePlaceholderValue($api, $auth['token'] ?? $auth['value'] ?? null);
            switch ($type) {
                case 'bearer':
                    $headers['Authorization'] = 'Bearer ' . $token;
                    break;
                case 'api_key':
                    $keyName = $auth['header'] ?? 'X-API-Key';
                    $headers[$keyName] = $token;
                    break;
                case 'basic':
                    $headers['Authorization'] = 'Basic ' . base64_encode($token);
                    break;
            }
        }

        return $headers;
    }

    private function parseResponse(ExternalApi $api, string $body): mixed
    {
        $type = $api->data_type ?? 'json';

        if ($type === 'xml') {
            $xml = simplexml_load_string($body, 'SimpleXMLElement', LIBXML_NOCDATA);
            return json_decode(json_encode($xml), true);
        }

        $decoded = json_decode($body, true);
        return $decoded ?? $body;
    }

    private function extractData(mixed $raw, ?string $itemsPath): array
    {
        if ($itemsPath === null || $itemsPath === '*' || $itemsPath === '') {
            return is_array($raw) ? $raw : [$raw];
        }

        $data = data_get($raw, $itemsPath, []);
        if (!is_array($data)) {
            return [$data];
        }

        return $data;
    }

    private function summariseResponse(mixed $raw): array
    {
        if (!is_array($raw)) {
            return ['raw' => Str::limit((string) $raw, 500)];
        }

        $collection = collect($raw);
        if ($collection->isAssoc()) {
            return $collection->map(function ($value) {
                return is_array($value) ? Arr::first($value) : $value;
            })->take(10)->toArray();
        }

        return $collection->take(5)->toArray();
    }

    private function resolvePlaceholders(ExternalApi $api, $value)
    {
        if (is_array($value)) {
            return collect($value)->map(function ($item) use ($api) {
                return $this->resolvePlaceholders($api, $item);
            })->toArray();
        }

        return $this->resolvePlaceholderValue($api, $value);
    }

    private function resolvePlaceholderValue(ExternalApi $api, $value)
    {
        if (!is_string($value)) {
            return $value;
        }

        $replacements = [
            '{{today}}' => now()->toDateString(),
            '{{yesterday}}' => now()->subDay()->toDateString(),
            '{{month_start}}' => now()->startOfMonth()->toDateString(),
            '{{month_end}}' => now()->endOfMonth()->toDateString(),
            '{{now_iso}}' => now()->toIso8601String(),
            '{{last_sync_at}}' => optional($api->last_sync_at)->toIso8601String(),
        ];

        if (Str::contains($value, '{{company_id}}')) {
            $companyId = config('idash.default_company_id');
            if ($companyId === null) {
                $companyId = config('app.default_company_id');
            }
            if ($companyId !== null) {
                $value = str_replace('{{company_id}}', $companyId, $value);
            }
        }

        foreach ($replacements as $placeholder => $replacement) {
            if ($replacement) {
                $value = str_replace($placeholder, $replacement, $value);
            }
        }

        if (preg_match('/^\{\{env:([^}]+)\}\}$/', $value, $matches)) {
            return env($matches[1], null);
        }

        return $value;
    }
}
