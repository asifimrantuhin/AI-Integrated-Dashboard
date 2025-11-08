<?php

namespace App\Console\Commands;

use App\Jobs\RunExternalApiSync;
use App\Models\ExternalApi;
use Illuminate\Console\Command;

class SyncDueExternalApis extends Command
{
    protected $signature = 'external-apis:sync-due {--module=}';
    protected $description = 'Dispatch sync jobs for external APIs that are due based on their interval.';

    public function handle(): int
    {
        $module = $this->option('module');
        $now = now();

        $query = ExternalApi::active();
        if ($module) {
            $query->where('module', $module);
        }

        $apis = $query->get()->filter(function (ExternalApi $api) use ($now) {
            if (!$api->sync_interval || $api->sync_interval <= 0) {
                return true;
            }

            if (!$api->last_sync_at) {
                return true;
            }

            return $api->last_sync_at->addMinutes($api->sync_interval)->lte($now);
        });

        if ($apis->isEmpty()) {
            $this->info('No external APIs are due for synchronization.');
            return self::SUCCESS;
        }

        foreach ($apis as $api) {
            RunExternalApiSync::dispatch($api);
            $this->line(sprintf('Queued sync for [%s] (%s)', $api->name, $api->module));
        }

        $this->info('Dispatch completed.');
        return self::SUCCESS;
    }
}
