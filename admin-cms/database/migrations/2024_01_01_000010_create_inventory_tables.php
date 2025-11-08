<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Inventory Raw Data
        Schema::create('inventory_raw_datas', function (Blueprint $table) {
            $table->id();
            $table->integer('company_id');
            $table->integer('gl_id');
            $table->date('month');
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();

            $table->index(['company_id', 'month']);
        });

        // Inventory GL Accounts
        Schema::create('inventory_gl_accounts', function (Blueprint $table) {
            $table->id();
            $table->string('gl_account');
            $table->string('gl_account_name');
            $table->text('description')->nullable();
            $table->timestamps();
        });

        // COGS GP
        Schema::create('cogs_gps', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->integer('company_id');
            $table->decimal('cogs', 15, 2)->default(0);
            $table->decimal('gp', 15, 2)->default(0);
            $table->decimal('gp_percentage', 10, 2)->default(0);
            $table->timestamps();

            $table->index(['month', 'company_id']);
        });

        // Inventory SAP Data
        Schema::create('inventroy_sap_datas', function (Blueprint $table) {
            $table->id();
            $table->integer('company');
            $table->integer('year');
            $table->integer('month');
            $table->timestamps();

            $table->index(['company', 'year', 'month']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('inventroy_sap_datas');
        Schema::dropIfExists('cogs_gps');
        Schema::dropIfExists('inventory_gl_accounts');
        Schema::dropIfExists('inventory_raw_datas');
    }
};

