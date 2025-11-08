<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        Schema::create('api_syncs', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('external_api_id');
            $table->enum('status', ['pending', 'processing', 'completed', 'failed'])->default('pending');
            $table->integer('records_synced')->default(0);
            $table->text('response_data')->nullable();
            $table->text('error_message')->nullable();
            $table->timestamp('started_at')->nullable();
            $table->timestamp('completed_at')->nullable();
            $table->timestamps();

            $table->foreign('external_api_id')->references('id')->on('external_apis')->onDelete('cascade');
            $table->index(['external_api_id', 'status']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('api_syncs');
    }
};

