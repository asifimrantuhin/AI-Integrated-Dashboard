<?php

namespace App\Console;

use App\Console\Commands\GenerateDemoData;
use App\Console\Commands\SyncDueExternalApis;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;

class Kernel extends ConsoleKernel
{
    /**
     * The Artisan commands provided by your application.
     *
     * @var array
     */
    protected $commands = [
        GenerateDemoData::class,
        SyncDueExternalApis::class,
    ];

    protected function schedule(Schedule $schedule): void
    {
        $schedule->command('external-apis:sync-due')->everyFifteenMinutes();
    }

    protected function commands(): void
    {
        $this->load(__DIR__ . '/Commands');

        if (file_exists(base_path('routes/console.php'))) {
            require base_path('routes/console.php');
        }
    }
}
