<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Production Analysis
        Schema::create('production_analyses', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->string('factory', 255);
            $table->string('summary_group', 255);
            $table->decimal('cmonthly_amount', 15, 2)->default(0);
            $table->decimal('cavg_amount', 15, 2)->default(0);
            $table->decimal('cmonthly_per', 10, 2)->default(0);
            $table->decimal('cavg_per', 10, 2)->default(0);
            $table->decimal('pmonthly_amount', 15, 2)->default(0);
            $table->decimal('pavg_amount', 15, 2)->default(0);
            $table->decimal('pmonthly_per', 10, 2)->default(0);
            $table->decimal('pavg_per', 10, 2)->default(0);
            $table->decimal('tmonthly_amount', 15, 2)->default(0);
            $table->decimal('tmonthly_per', 10, 2)->default(0);
            $table->decimal('amonthly_amount', 15, 2)->default(0);
            $table->decimal('amonthly_per', 10, 2)->default(0);
            $table->decimal('aavg_amount', 15, 2)->default(0);
            $table->timestamps();
            $table->softDeletes();

            $table->index(['month', 'factory']);
        });

        // Wastage Data
        Schema::create('wastage_datas', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->string('factory', 255);
            $table->string('group_name', 300)->nullable();
            $table->float('std')->default(0);
            $table->decimal('wastage', 15, 2)->default(0);
            $table->decimal('month_wastage', 15, 2)->default(0);
            $table->decimal('avg', 15, 2)->default(0);
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();
            $table->softDeletes();

            $table->index(['month', 'factory']);
        });

        // Cost Analysis
        Schema::create('cost_analyses', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->string('factory', 255);
            $table->string('cost_type', 255);
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();
            $table->softDeletes();

            $table->index(['month', 'factory', 'cost_type']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('cost_analyses');
        Schema::dropIfExists('wastage_datas');
        Schema::dropIfExists('production_analyses');
    }
};

