<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Employee Basic Info
        Schema::create('employee_basic_infos', function (Blueprint $table) {
            $table->id();
            $table->integer('total_active_staff')->default(0);
            $table->integer('total_active_worker')->default(0);
            $table->integer('total_contractual_employee')->default(0);
            $table->integer('total_probationary_employee')->default(0);
            $table->integer('total_permanent_employee')->default(0);
            $table->date('report_date')->nullable();
            $table->timestamps();

            $table->index('report_date');
        });

        // Employee Attendance
        Schema::create('employee_attendances', function (Blueprint $table) {
            $table->id();
            $table->date('date');
            $table->integer('total_absent')->default(0);
            $table->integer('total_present')->default(0);
            $table->integer('total_leave')->default(0);
            $table->timestamps();

            $table->index('date');
        });

        // Employee Turnover
        Schema::create('employee_tran_overs', function (Blueprint $table) {
            $table->id();
            $table->string('job_type');
            $table->string('month');
            $table->integer('year');
            $table->integer('new_employee_no')->default(0);
            $table->integer('resigned_employee')->default(0);
            $table->timestamps();

            $table->index(['year', 'month']);
        });

        // Yearly Employee Promotion
        Schema::create('yearly_employee_promotions', function (Blueprint $table) {
            $table->id();
            $table->integer('year');
            $table->integer('promoted_count')->default(0);
            $table->text('details')->nullable();
            $table->timestamps();

            $table->index('year');
        });

        // HRIS Companies
        Schema::create('hris_companies', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('code')->unique();
            $table->timestamps();
        });

        // HRIS Departments
        Schema::create('hris_departments', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->unsignedBigInteger('company_id');
            $table->timestamps();

            $table->foreign('company_id')->references('id')->on('hris_companies')->onDelete('cascade');
        });

        // HRIS Promotion Break Downs
        Schema::create('hris_promotion_break_downs', function (Blueprint $table) {
            $table->id();
            $table->integer('year');
            $table->unsignedBigInteger('department_id');
            $table->integer('promoted_count')->default(0);
            $table->timestamps();

            $table->foreign('department_id')->references('id')->on('hris_departments')->onDelete('cascade');
            $table->index('year');
        });
    }

    public function down()
    {
        Schema::dropIfExists('hris_promotion_break_downs');
        Schema::dropIfExists('hris_departments');
        Schema::dropIfExists('hris_companies');
        Schema::dropIfExists('yearly_employee_promotions');
        Schema::dropIfExists('employee_tran_overs');
        Schema::dropIfExists('employee_attendances');
        Schema::dropIfExists('employee_basic_infos');
    }
};

