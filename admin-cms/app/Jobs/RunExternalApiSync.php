<?php

namespace App\Jobs;

use App\Models\ExternalApi;
use App\Services\ExternalApi\SyncService;
use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;

class RunExternalApiSync implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    public int $tries = 3;
    public int $timeout = 120;

    public function __construct(private readonly int $externalApiId)
    {
    }

    public static function fromApi(ExternalApi $api): self
    {
        return new self($api->id);
    }

    public function handle(SyncService $service): void
    {
        $api = ExternalApi::find($this->externalApiId);

        if (!$api || !$api->is_active) {
            return;
        }

        try {
            $service->sync($api);
        } catch (\Throwable $e) {
            Log::error('External API sync job failed', [
                'external_api_id' => $this->externalApiId,
                'error' => $e->getMessage(),
            ]);
            throw $e;
        }
    }
}
