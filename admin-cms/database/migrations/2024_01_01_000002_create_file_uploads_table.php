<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        Schema::create('file_uploads', function (Blueprint $table) {
            $table->id();
            $table->string('file_name');
            $table->string('file_path');
            $table->string('module'); // sales, production, finance, etc.
            $table->string('data_type'); // monthly_report, daily_report, etc.
            $table->unsignedBigInteger('uploaded_by');
            $table->enum('status', ['pending', 'processing', 'completed', 'failed'])->default('pending');
            $table->integer('records_processed')->default(0);
            $table->text('error_message')->nullable();
            $table->timestamps();

            $table->foreign('uploaded_by')->references('id')->on('users')->onDelete('cascade');
            $table->index(['module', 'data_type']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('file_uploads');
    }
};

