<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Quality Control Tables
        Schema::create('quality_control_checks', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('factory_id')->index();
            $table->unsignedBigInteger('production_line_id')->nullable()->index();
            $table->string('product_code', 100)->index();
            $table->string('product_name', 255);
            $table->date('check_date')->index();
            $table->enum('check_type', ['incoming', 'in_process', 'final', 'sampling'])->index();
            $table->enum('status', ['passed', 'failed', 'pending', 'rework'])->index();
            $table->integer('total_checked')->default(0);
            $table->integer('passed_count')->default(0);
            $table->integer('failed_count')->default(0);
            $table->decimal('defect_rate', 5, 2)->default(0);
            $table->text('defects')->nullable();
            $table->text('remarks')->nullable();
            $table->string('inspector_name', 255)->nullable();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['check_date', 'status']);
            $table->index(['company_id', 'factory_id', 'check_date']);
        });

        // Production Planning
        Schema::create('production_plans', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('factory_id')->index();
            $table->string('plan_number', 100)->unique()->index();
            $table->string('product_code', 100)->index();
            $table->string('product_name', 255);
            $table->date('plan_date')->index();
            $table->date('start_date')->index();
            $table->date('end_date')->index();
            $table->decimal('planned_quantity', 15, 2)->default(0);
            $table->decimal('actual_quantity', 15, 2)->default(0);
            $table->enum('status', ['draft', 'approved', 'in_progress', 'completed', 'cancelled'])->index();
            $table->enum('priority', ['low', 'medium', 'high', 'urgent'])->default('medium')->index();
            $table->text('notes')->nullable();
            $table->unsignedBigInteger('created_by')->nullable();
            $table->unsignedBigInteger('approved_by')->nullable();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['plan_date', 'status']);
            $table->index(['factory_id', 'start_date', 'end_date']);
        });

        // Machine Maintenance
        Schema::create('machine_maintenances', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('factory_id')->index();
            $table->string('machine_code', 100)->index();
            $table->string('machine_name', 255);
            $table->enum('maintenance_type', ['preventive', 'corrective', 'breakdown', 'scheduled'])->index();
            $table->date('maintenance_date')->index();
            $table->time('start_time')->nullable();
            $table->time('end_time')->nullable();
            $table->integer('downtime_minutes')->default(0);
            $table->text('description')->nullable();
            $table->text('actions_taken')->nullable();
            $table->decimal('cost', 15, 2)->default(0);
            $table->enum('status', ['scheduled', 'in_progress', 'completed', 'cancelled'])->index();
            $table->string('technician_name', 255)->nullable();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['maintenance_date', 'status']);
            $table->index(['machine_code', 'maintenance_date']);
        });

        // Material Requirements Planning (MRP)
        Schema::create('material_requirements', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('production_plan_id')->nullable()->index();
            $table->string('material_code', 100)->index();
            $table->string('material_name', 255);
            $table->string('material_type', 100)->index();
            $table->string('unit', 50);
            $table->decimal('required_quantity', 15, 2)->default(0);
            $table->decimal('available_quantity', 15, 2)->default(0);
            $table->decimal('shortage_quantity', 15, 2)->default(0);
            $table->date('required_date')->index();
            $table->enum('status', ['pending', 'ordered', 'received', 'shortage'])->index();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['required_date', 'status']);
            $table->index(['material_code', 'required_date']);
        });

        // Production Efficiency Metrics
        Schema::create('production_efficiency', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('factory_id')->index();
            $table->unsignedBigInteger('production_line_id')->nullable()->index();
            $table->date('production_date')->index();
            $table->string('shift', 50)->nullable()->index();
            $table->decimal('planned_output', 15, 2)->default(0);
            $table->decimal('actual_output', 15, 2)->default(0);
            $table->decimal('efficiency_percentage', 5, 2)->default(0)->index();
            $table->integer('planned_hours')->default(0);
            $table->integer('actual_hours')->default(0);
            $table->integer('downtime_minutes')->default(0);
            $table->decimal('oee', 5, 2)->default(0)->comment('Overall Equipment Effectiveness');
            $table->timestamps();

            $table->index(['production_date', 'factory_id']);
            $table->index(['production_date', 'efficiency_percentage']);
        });

        // Supplier Performance
        Schema::create('supplier_performance', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->string('supplier_code', 100)->index();
            $table->string('supplier_name', 255);
            $table->date('evaluation_date')->index();
            $table->integer('total_orders')->default(0);
            $table->integer('on_time_deliveries')->default(0);
            $table->decimal('on_time_percentage', 5, 2)->default(0)->index();
            $table->integer('quality_issues')->default(0);
            $table->decimal('quality_score', 5, 2)->default(0)->index();
            $table->decimal('cost_score', 5, 2)->default(0);
            $table->decimal('overall_score', 5, 2)->default(0)->index();
            $table->enum('rating', ['excellent', 'good', 'average', 'poor'])->index();
            $table->text('comments')->nullable();
            $table->timestamps();

            $table->index(['evaluation_date', 'rating']);
            $table->index(['supplier_code', 'evaluation_date']);
        });

        // Energy Consumption
        Schema::create('energy_consumption', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id')->index();
            $table->unsignedBigInteger('factory_id')->index();
            $table->date('consumption_date')->index();
            $table->enum('energy_type', ['electricity', 'gas', 'water', 'steam', 'compressed_air'])->index();
            $table->decimal('consumption_amount', 15, 2)->default(0);
            $table->string('unit', 50);
            $table->decimal('cost', 15, 2)->default(0);
            $table->string('meter_reading', 100)->nullable();
            $table->timestamps();

            $table->index(['consumption_date', 'energy_type']);
            $table->index(['factory_id', 'consumption_date']);
        });

        // AI Forecast Results Storage
        Schema::create('ai_forecasts', function (Blueprint $table) {
            $table->id();
            $table->string('forecast_type', 50)->index()->comment('sales, production, finance, inventory');
            $table->string('entity_type', 50)->nullable()->index()->comment('product, channel, factory, etc.');
            $table->string('entity_id', 100)->nullable()->index();
            $table->date('forecast_date')->index();
            $table->decimal('forecasted_value', 15, 2)->default(0);
            $table->decimal('confidence_level', 5, 2)->default(0);
            $table->decimal('upper_bound', 15, 2)->nullable();
            $table->decimal('lower_bound', 15, 2)->nullable();
            $table->string('model_used', 100)->nullable();
            $table->json('forecast_details')->nullable();
            $table->enum('status', ['pending', 'active', 'expired'])->default('active')->index();
            $table->timestamps();
            $table->softDeletes();

            $table->index(['forecast_type', 'forecast_date', 'status']);
            $table->index(['entity_type', 'entity_id', 'forecast_date']);
        });

        // Report Filters Configuration
        Schema::create('report_filters', function (Blueprint $table) {
            $table->id();
            $table->string('report_type', 100)->index();
            $table->string('filter_name', 100)->index();
            $table->string('filter_type', 50)->comment('date_range, dropdown, multi_select, text');
            $table->json('filter_options')->nullable();
            $table->boolean('is_required')->default(false);
            $table->integer('sort_order')->default(0);
            $table->boolean('is_active')->default(true);
            $table->timestamps();

            $table->index(['report_type', 'is_active']);
        });

        // Cache for API Responses
        Schema::create('api_cache', function (Blueprint $table) {
            $table->id();
            $table->string('cache_key', 255)->unique()->index();
            $table->text('cache_value');
            $table->string('endpoint', 255)->index();
            $table->text('parameters')->nullable();
            $table->timestamp('expires_at')->index();
            $table->timestamps();

            $table->index(['endpoint', 'expires_at']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('api_cache');
        Schema::dropIfExists('report_filters');
        Schema::dropIfExists('ai_forecasts');
        Schema::dropIfExists('energy_consumption');
        Schema::dropIfExists('supplier_performance');
        Schema::dropIfExists('production_efficiency');
        Schema::dropIfExists('material_requirements');
        Schema::dropIfExists('machine_maintenances');
        Schema::dropIfExists('production_plans');
        Schema::dropIfExists('quality_control_checks');
    }
};

