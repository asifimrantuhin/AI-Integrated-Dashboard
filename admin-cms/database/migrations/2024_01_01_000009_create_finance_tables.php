<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Budget Categories
        Schema::create('bdgt_categories', function (Blueprint $table) {
            $table->id();
            $table->string('name', 50);
            $table->integer('status')->default(1);
            $table->timestamps();
        });

        // Budget Departments
        Schema::create('bdgt_departments', function (Blueprint $table) {
            $table->id();
            $table->string('name', 50);
            $table->unsignedBigInteger('category_id');
            $table->integer('status')->default(1);
            $table->timestamps();

            $table->foreign('category_id')->references('id')->on('bdgt_categories')->onDelete('cascade');
        });

        // Budget Expense Groups
        Schema::create('bdgt_expense_groups', function (Blueprint $table) {
            $table->id();
            $table->string('name', 50);
            $table->unsignedBigInteger('department_id');
            $table->integer('status')->default(1);
            $table->timestamps();

            $table->foreign('department_id')->references('id')->on('bdgt_departments')->onDelete('cascade');
        });

        // Budget Sub Heads
        Schema::create('bdgt_sub_heads', function (Blueprint $table) {
            $table->id();
            $table->string('name', 50);
            $table->unsignedBigInteger('expense_group_id');
            $table->integer('status')->default(1);
            $table->timestamps();

            $table->foreign('expense_group_id')->references('id')->on('bdgt_expense_groups')->onDelete('cascade');
        });

        // Budget Expenses
        Schema::create('bdgt_expenses', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->unsignedBigInteger('sub_head_id');
            $table->timestamps();

            $table->foreign('sub_head_id')->references('id')->on('bdgt_sub_heads')->onDelete('cascade');
        });

        // Budget Summaries
        Schema::create('budget_summaries', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->unsignedBigInteger('category_id');
            $table->unsignedBigInteger('department_id');
            $table->decimal('budget_amount', 15, 2)->default(0);
            $table->decimal('actual_amount', 15, 2)->default(0);
            $table->timestamps();

            $table->foreign('category_id')->references('id')->on('bdgt_categories')->onDelete('cascade');
            $table->foreign('department_id')->references('id')->on('bdgt_departments')->onDelete('cascade');
            $table->index(['month', 'category_id']);
        });

        // Budget Monthlies
        Schema::create('budget_monthlies', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->unsignedBigInteger('expense_id');
            $table->decimal('budget_amount', 15, 2)->default(0);
            $table->decimal('actual_amount', 15, 2)->default(0);
            $table->timestamps();

            $table->foreign('expense_id')->references('id')->on('bdgt_expenses')->onDelete('cascade');
            $table->index(['month', 'expense_id']);
        });

        // Bank Loan Heads
        Schema::create('bank_loan_heads', function (Blueprint $table) {
            $table->id();
            $table->string('loan_head');
            $table->timestamps();
        });

        // Bank Loan Status Raw Data
        Schema::create('bank_loan_status_raw_data', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->string('loan_head');
            $table->string('company_id', 25);
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();

            $table->index(['month', 'company_id']);
        });

        // Financial Expense Raw Data
        Schema::create('financial_expense_raw_data', function (Blueprint $table) {
            $table->id();
            $table->date('month');
            $table->unsignedBigInteger('expense_id');
            $table->decimal('amount', 15, 2)->default(0);
            $table->timestamps();

            $table->foreign('expense_id')->references('id')->on('bdgt_expenses')->onDelete('cascade');
            $table->index(['month', 'expense_id']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('financial_expense_raw_data');
        Schema::dropIfExists('bank_loan_status_raw_data');
        Schema::dropIfExists('bank_loan_heads');
        Schema::dropIfExists('budget_monthlies');
        Schema::dropIfExists('budget_summaries');
        Schema::dropIfExists('bdgt_expenses');
        Schema::dropIfExists('bdgt_sub_heads');
        Schema::dropIfExists('bdgt_expense_groups');
        Schema::dropIfExists('bdgt_departments');
        Schema::dropIfExists('bdgt_categories');
    }
};

