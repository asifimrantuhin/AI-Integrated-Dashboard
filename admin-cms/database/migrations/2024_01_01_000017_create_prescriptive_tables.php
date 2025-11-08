<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        Schema::create('scenario_simulations', function (Blueprint $table) {
            $table->id();
            $table->string('module');
            $table->string('scenario_key')->nullable();
            $table->json('base_metrics');
            $table->json('adjustments');
            $table->json('results');
            $table->unsignedBigInteger('created_by')->nullable();
            $table->timestamps();

            $table->index(['module', 'created_at']);
        });

        Schema::create('prescriptive_recommendations', function (Blueprint $table) {
            $table->id();
            $table->string('module');
            $table->string('entity_type')->nullable();
            $table->string('entity_id')->nullable();
            $table->string('risk_level')->default('medium');
            $table->json('recommendation');
            $table->decimal('impact_score', 8, 2)->default(0);
            $table->json('metadata')->nullable();
            $table->timestamps();

            $table->index(['module', 'entity_type']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('prescriptive_recommendations');
        Schema::dropIfExists('scenario_simulations');
    }
};
