<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up()
    {
        // Supply Chain Raw Data
        Schema::create('supply_chain_raw_datas', function (Blueprint $table) {
            $table->id();
            $table->integer('plant');
            $table->integer('pr_id');
            $table->integer('pr_item');
            $table->date('pr_date');
            $table->string('pr_material');
            $table->string('pr_material_text');
            $table->integer('pr_qty');
            $table->string('prunit_measure');
            $table->date('prrelease_date');
            $table->integer('po_id');
            $table->integer('po_item');
            $table->date('po_date');
            $table->integer('vendor_id');
            $table->string('vendor_details');
            $table->integer('po_qty');
            $table->string('pounit_measure');
            $table->decimal('po_amount', 15, 2)->default(0);
            $table->string('po_currency', 25);
            $table->date('po_delivery_date');
            $table->date('po_released_date');
            $table->integer('lc_number');
            $table->integer('actual_po1');
            $table->integer('actual_po2');
            $table->integer('actual_po1item');
            $table->integer('actual_po2item');
            $table->integer('grn1_id');
            $table->integer('grn1_item');
            $table->date('grn1_date');
            $table->integer('grn1_qtn');
            $table->decimal('grn1_amount', 15, 2)->default(0);
            $table->date('invoice1_date');
            $table->string('invoice1_id');
            $table->date('invoice1_vendor_date');
            $table->integer('invoice1_qty');
            $table->decimal('invoice1_amount', 15, 2)->default(0);
            $table->integer('invoice1_dn');
            $table->timestamps();

            $table->index(['plant', 'po_date']);
        });

        // Supply Chain PO
        Schema::create('supply_chain_pos', function (Blueprint $table) {
            $table->id();
            $table->unsignedBigInteger('company_id');
            $table->integer('plant');
            $table->integer('vendor_code');
            $table->string('vandor_name', 200);
            $table->string('material_code', 200);
            $table->string('material_name', 200);
            $table->string('material_group', 200);
            $table->string('material_group_description', 200);
            $table->string('purchasing_organization', 200);
            $table->string('purchasing_group', 200);
            $table->integer('pr_number');
            $table->integer('po_number');
            $table->integer('po_item_number');
            $table->date('po_date');
            $table->integer('po_qty');
            $table->string('uom', 50);
            $table->string('po_currency', 25);
            $table->decimal('po_amount', 15, 2)->default(0);
            $table->decimal('master_po', 15, 2)->default(0);
            $table->timestamps();

            $table->foreign('company_id')->references('id')->on('companies')->onDelete('cascade');
            $table->index(['company_id', 'po_date']);
        });

        // Purchase Requisitions
        Schema::create('purchase_requisitions', function (Blueprint $table) {
            $table->id();
            $table->integer('pr_id');
            $table->integer('pr_item');
            $table->date('pr_date');
            $table->string('material_code');
            $table->string('material_name');
            $table->integer('quantity');
            $table->string('unit');
            $table->integer('plant');
            $table->timestamps();

            $table->index(['pr_id', 'pr_date']);
        });
    }

    public function down()
    {
        Schema::dropIfExists('purchase_requisitions');
        Schema::dropIfExists('supply_chain_pos');
        Schema::dropIfExists('supply_chain_raw_datas');
    }
};

