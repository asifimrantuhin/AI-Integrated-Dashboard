<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        Schema::create('external_apis', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('url');
            $table->string('method')->default('GET'); // GET, POST, PUT, etc.
            $table->text('headers')->nullable(); // JSON encoded headers
            $table->text('body')->nullable(); // JSON encoded request body
            $table->text('authentication')->nullable(); // JSON encoded auth config
            $table->string('module'); // sales, production, finance, etc.
            $table->string('data_type'); // Type of data this API provides
            $table->text('data_mapping')->nullable(); // JSON encoded field mappings
            $table->boolean('is_active')->default(true);
            $table->integer('sync_interval')->default(60); // Minutes
            $table->timestamp('last_sync_at')->nullable();
            $table->timestamps();

            $table->index(['module', 'is_active']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('external_apis');
    }
};

