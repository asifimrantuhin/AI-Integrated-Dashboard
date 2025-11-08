<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Monthly Report Table
        Schema::create('channelwise_monthly_report', function (Blueprint $table) {
            $table->id();
            $table->date('data_month');
            $table->integer('channel_id');
            $table->string('channel_name');
            $table->decimal('lifting_target', 15, 2)->default(0);
            $table->decimal('billed', 15, 2)->default(0);
            $table->decimal('delivered', 15, 2)->default(0);
            $table->decimal('primary_collection', 15, 2)->default(0);
            $table->decimal('ims_target', 15, 2)->default(0);
            $table->decimal('ims', 15, 2)->default(0);
            $table->decimal('market_collection', 15, 2)->default(0);
            $table->decimal('memo_target', 15, 2)->default(0);
            $table->decimal('memo_qty', 15, 2)->default(0);
            $table->decimal('pg_target', 15, 2)->default(0);
            $table->decimal('pg_cover', 15, 2)->default(0);
            $table->integer('total_retailer')->default(0);
            $table->integer('business_retailer')->default(0);
            $table->timestamps();

            $table->index(['data_month', 'channel_id']);
        });

        // Daily Report Table
        Schema::create('channelwise_lic_data', function (Blueprint $table) {
            $table->id();
            $table->date('data_date');
            $table->integer('channel_id');
            $table->string('channel_name');
            $table->decimal('lifting_target', 15, 2)->default(0);
            $table->decimal('billed', 15, 2)->default(0);
            $table->decimal('delivery', 15, 2)->default(0);
            $table->decimal('lifting_collection', 15, 2)->default(0);
            $table->decimal('ims_target', 15, 2)->default(0);
            $table->decimal('ims', 15, 2)->default(0);
            $table->decimal('ims_collection', 15, 2)->default(0);
            $table->timestamps();

            $table->index(['data_date', 'channel_id']);
        });

        // Best Selling Products
        Schema::create('best_selling_products', function (Blueprint $table) {
            $table->id();
            $table->date('year_month');
            $table->integer('channel_id');
            $table->integer('product_id');
            $table->string('product_name');
            $table->decimal('qty', 15, 2)->default(0);
            $table->decimal('value', 15, 2)->default(0);
            $table->integer('cat_id')->nullable();
            $table->timestamps();

            $table->index(['year_month', 'channel_id']);
        });

        // Best Selling PGs
        Schema::create('best_selling_pgs', function (Blueprint $table) {
            $table->id();
            $table->date('year_month');
            $table->integer('channel_id');
            $table->integer('category_id');
            $table->string('category_name');
            $table->decimal('qty', 15, 2)->default(0);
            $table->decimal('value', 15, 2)->default(0);
            $table->timestamps();

            $table->index(['year_month', 'channel_id']);
        });

        // Top Channel Distributors/Retailers
        Schema::create('top_channel_d_bs', function (Blueprint $table) {
            $table->id();
            $table->string('db_name');
            $table->decimal('amount', 15, 2)->default(0);
            $table->integer('type')->default(0); // 0 for distributor, 1 for retailer
            $table->date('date');
            $table->timestamps();

            $table->index(['date', 'type']);
        });

        // Order Delivery Summary
        Schema::create('order_delivery_summaries', function (Blueprint $table) {
            $table->id();
            $table->date('months');
            $table->integer('channel_id');
            $table->decimal('amounts', 15, 2)->default(0);
            $table->integer('types')->default(0); // 0 for order, 1 for delivery
            $table->timestamps();

            $table->index(['months', 'channel_id']);
        });

        // Top Retailers
        Schema::create('top_retailers', function (Blueprint $table) {
            $table->id();
            $table->date('date');
            $table->string('db_name');
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();

            $table->index('date');
        });

        // Sales Orders
        Schema::create('sales_orders', function (Blueprint $table) {
            $table->id();
            $table->integer('sales_organization');
            $table->string('dcument_type', 100);
            $table->string('so_number', 100);
            $table->string('customer_code', 100);
            $table->date('document_date');
            $table->string('line_item', 25);
            $table->string('product_code', 100);
            $table->string('product_name', 200);
            $table->string('product_group', 200);
            $table->string('so_qty', 50);
            $table->string('unit_price', 30);
            $table->string('total_price', 30);
            $table->string('currency', 25);
            $table->timestamps();

            $table->index(['so_number', 'document_date']);
        });

        // Sales Deliveries
        Schema::create('sales_deliveries', function (Blueprint $table) {
            $table->id();
            $table->integer('recference_no');
            $table->integer('factory_sap_code');
            $table->integer('delivered_cat_sap_code');
            $table->string('delivered_item_sap_code');
            $table->integer('delivered_quantity');
            $table->string('delivered_value', 100);
            $table->date('item_delivered_date');
            $table->integer('delivery_staus');
            $table->integer('no_of_partial_delivery');
            $table->string('sap_sales_order_no', 100);
            $table->string('sap_chalan_no', 100);
            $table->string('order_type', 10);
            $table->timestamps();

            $table->index(['sap_sales_order_no', 'item_delivered_date']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('sales_deliveries');
        Schema::dropIfExists('sales_orders');
        Schema::dropIfExists('top_retailers');
        Schema::dropIfExists('order_delivery_summaries');
        Schema::dropIfExists('top_channel_d_bs');
        Schema::dropIfExists('best_selling_pgs');
        Schema::dropIfExists('best_selling_products');
        Schema::dropIfExists('channelwise_lic_data');
        Schema::dropIfExists('channelwise_monthly_report');
    }
};

