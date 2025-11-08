<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // NPD Projects
        Schema::create('npd_projects', function (Blueprint $table) {
            $table->id();
            $table->integer('p_id')->unique();
            $table->string('indent_no', 100);
            $table->string('name', 150);
            $table->string('pmo', 150)->comment('project management office');
            $table->string('project_manager', 150);
            $table->string('location', 80);
            $table->string('reason', 250);
            $table->string('details', 250);
            $table->string('sponsors', 100);
            $table->date('start_date');
            $table->date('end_date');
            $table->string('type', 50);
            $table->string('progress', 20);
            $table->decimal('budget', 15, 2)->default(0);
            $table->string('lead_department', 70);
            $table->string('reponsible_departments', 70);
            $table->string('status', 20);
            $table->string('status_name', 80);
            $table->string('status_background', 40);
            $table->timestamps();

            $table->index(['p_id', 'status']);
        });

        // Project Deliverables
        Schema::create('projects_deliverables', function (Blueprint $table) {
            $table->id();
            $table->integer('d_id')->unique();
            $table->string('name', 100);
            $table->float('weightage', 150);
            $table->date('start_date');
            $table->date('end_date');
            $table->string('acknowledges', 100);
            $table->decimal('budget', 15, 2)->default(0);
            $table->decimal('progress', 5, 2)->default(0);
            $table->unsignedBigInteger('npd_project_id');
            $table->timestamps();

            $table->foreign('npd_project_id')->references('id')->on('npd_projects')->onDelete('cascade');
            $table->index('npd_project_id');
        });

        // Project Sub Deliverables
        Schema::create('projects_sub_deliverables', function (Blueprint $table) {
            $table->id();
            $table->integer('sd_id')->unique();
            $table->string('name', 100);
            $table->float('weightage', 150);
            $table->date('start_date');
            $table->date('end_date');
            $table->string('acknowledges', 100);
            $table->decimal('budget', 15, 2)->default(0);
            $table->decimal('progress', 5, 2)->default(0);
            $table->unsignedBigInteger('deliverable_id');
            $table->timestamps();

            $table->foreign('deliverable_id')->references('id')->on('projects_deliverables')->onDelete('cascade');
            $table->index('deliverable_id');
        });
    }

    public function down()
    {
        Schema::dropIfExists('projects_sub_deliverables');
        Schema::dropIfExists('projects_deliverables');
        Schema::dropIfExists('npd_projects');
    }
};

