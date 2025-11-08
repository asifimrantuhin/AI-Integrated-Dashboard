<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;
use Illuminate\Support\Facades\DB;

return new class extends Migration
{
    public function up(): void
    {
        Schema::table('external_apis', function (Blueprint $table) {
            if (!Schema::hasColumn('external_apis', 'module')) {
                $table->string('module')->default('sales')->after('name');
            }

            if (!Schema::hasColumn('external_apis', 'url')) {
                $table->string('url')->nullable()->after('module');
            }

            if (!Schema::hasColumn('external_apis', 'query_params')) {
                $table->text('query_params')->nullable()->after('headers');
            }

            if (!Schema::hasColumn('external_apis', 'items_path')) {
                $table->string('items_path')->nullable()->after('data_type');
            }

            if (!Schema::hasColumn('external_apis', 'is_active')) {
                $table->boolean('is_active')->default(true)->after('data_mapping');
            }

            if (!Schema::hasColumn('external_apis', 'sync_interval')) {
                $table->integer('sync_interval')->default(60)->after('is_active');
            }

            if (!Schema::hasColumn('external_apis', 'last_sync_at')) {
                $table->timestamp('last_sync_at')->nullable()->after('sync_interval');
            }
        });

        if (Schema::hasColumn('external_apis', 'endpoint')) {
            DB::table('external_apis')->whereNull('url')->update(['url' => DB::raw('endpoint')]);
        }
    }

    public function down(): void
    {
        Schema::table('external_apis', function (Blueprint $table) {
            if (Schema::hasColumn('external_apis', 'query_params')) {
                $table->dropColumn('query_params');
            }
            if (Schema::hasColumn('external_apis', 'items_path')) {
                $table->dropColumn('items_path');
            }
            if (Schema::hasColumn('external_apis', 'is_active')) {
                $table->dropColumn('is_active');
            }
            if (Schema::hasColumn('external_apis', 'sync_interval')) {
                $table->dropColumn('sync_interval');
            }
            if (Schema::hasColumn('external_apis', 'last_sync_at')) {
                $table->dropColumn('last_sync_at');
            }
            if (Schema::hasColumn('external_apis', 'url')) {
                $table->dropColumn('url');
            }
            if (Schema::hasColumn('external_apis', 'module')) {
                $table->dropColumn('module');
            }
        });
    }
};
