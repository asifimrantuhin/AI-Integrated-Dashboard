<?php

namespace App\Services\ExternalApi;

use App\Models\ExternalApi;
use Illuminate\Support\Arr;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Str;
use Illuminate\Support\Carbon;

class DataMapper
{
    public function mapAndPersist(ExternalApi $api, array $payload): int
    {
        $mapping = $api->data_mapping ?? [];
        if (empty($mapping)) {
            return 0;
        }

        $mappings = $mapping['mappings'] ?? [$mapping];
        $itemsPathDefault = $mapping['items_path'] ?? $api->items_path;
        $processed = 0;

        foreach ($mappings as $map) {
            $itemsPath = $map['items_path'] ?? $itemsPathDefault;
            $items = $this->resolveItems($payload, $itemsPath);
            if (!is_iterable($items)) {
                continue;
            }

            $targetTable = $map['target_table'] ?? null;
            $fieldMap = $map['fields'] ?? [];
            $transforms = $map['transforms'] ?? [];
            $upsertKeys = $map['upsert_keys'] ?? [];

            if (!$targetTable || empty($fieldMap)) {
                continue;
            }

            foreach ($items as $item) {
                if (!is_array($item)) {
                    continue;
                }

                $row = [];
                foreach ($fieldMap as $column => $path) {
                    $value = data_get($item, $path);
                    if (array_key_exists($column, $transforms)) {
                        $value = $this->applyTransform($value, $transforms[$column]);
                    }
                    $row[$column] = $value;
                }

                $timestamp = now();
                if (!array_key_exists('updated_at', $row)) {
                    $row['updated_at'] = $timestamp;
                }
                if (!array_key_exists('created_at', $row)) {
                    $row['created_at'] = $timestamp;
                }

                if (!empty($upsertKeys)) {
                    $keys = Arr::only($row, $upsertKeys);
                    DB::table($targetTable)->updateOrInsert($keys, $row);
                } else {
                    DB::table($targetTable)->insert($row);
                }

                $processed++;
            }
        }

        return $processed;
    }

    private function resolveItems(array $payload, ?string $itemsPath)
    {
        if (!$itemsPath) {
            return $payload;
        }

        if ($itemsPath === '*') {
            return $payload;
        }

        $items = data_get($payload, $itemsPath);
        return $items ?? [];
    }

    private function applyTransform($value, string $instruction)
    {
        if ($value === null) {
            return null;
        }

        if ($instruction === 'int') {
            return (int) $value;
        }
        if ($instruction === 'float') {
            return (float) $value;
        }
        if ($instruction === 'string') {
            return (string) $value;
        }
        if (Str::startsWith($instruction, 'date:')) {
            $format = Str::after($instruction, 'date:');
            return $this->formatDate($value, $format ?: 'Y-m-d');
        }
        if ($instruction === 'bool') {
            return filter_var($value, FILTER_VALIDATE_BOOLEAN, FILTER_NULL_ON_FAILURE);
        }

        return $value;
    }

    private function formatDate($value, string $format)
    {
        try {
            return Carbon::parse($value)->format($format);
        } catch (\Throwable $e) {
            return $value;
        }
    }
}
