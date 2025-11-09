<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('sales_channel_targets', function (Blueprint $table) {
            $table->id();
            $table->date('data_month');
            $table->unsignedBigInteger('channel_id')->nullable();
            $table->string('channel_name')->nullable();
            $table->decimal('revenue_target', 18, 2)->default(0);
            $table->decimal('volume_target', 18, 2)->default(0);
            $table->decimal('promotion_budget', 18, 2)->default(0);
            $table->decimal('gross_margin_target', 18, 2)->default(0);
            $table->decimal('new_customer_target', 18, 2)->default(0);
            $table->string('owner')->nullable();
            $table->timestamps();

            $table->index(['data_month', 'channel_id']);
        });

        Schema::create('sales_order_book', function (Blueprint $table) {
            $table->id();
            $table->string('order_number')->index();
            $table->date('order_date');
            $table->unsignedBigInteger('channel_id')->nullable();
            $table->string('channel_name')->nullable();
            $table->string('customer_code')->nullable();
            $table->string('customer_name')->nullable();
            $table->string('region')->nullable();
            $table->enum('status', ['draft', 'confirmed', 'dispatching', 'delivered', 'cancelled'])->default('confirmed');
            $table->decimal('order_amount', 18, 2)->default(0);
            $table->date('fulfilled_at')->nullable();
            $table->decimal('discount_amount', 18, 2)->default(0);
            $table->decimal('gross_margin', 18, 2)->default(0);
            $table->timestamps();

            $table->index(['order_date', 'channel_id']);
            $table->index(['status', 'fulfilled_at']);
        });

        Schema::create('sales_promotion_performance', function (Blueprint $table) {
            $table->id();
            $table->string('campaign_code')->index();
            $table->string('campaign_name')->nullable();
            $table->unsignedBigInteger('channel_id')->nullable();
            $table->string('channel_name')->nullable();
            $table->date('start_date')->nullable();
            $table->date('end_date')->nullable();
            $table->decimal('spend_amount', 18, 2)->default(0);
            $table->decimal('revenue_uplift', 18, 2)->default(0);
            $table->decimal('uplift_percentage', 8, 2)->default(0);
            $table->decimal('roi', 8, 2)->default(0);
            $table->json('audience_tags')->nullable();
            $table->timestamps();

            $table->index(['start_date', 'end_date']);
            $table->index(['channel_id']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('sales_promotion_performance');
        Schema::dropIfExists('sales_order_book');
        Schema::dropIfExists('sales_channel_targets');
    }
};
