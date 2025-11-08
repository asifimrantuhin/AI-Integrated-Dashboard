<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\FileUploadController;
use App\Http\Controllers\ExternalApiController;
use App\Http\Controllers\DataViewController;
use App\Http\Controllers\SalesUploadController;
use App\Http\Controllers\ProductionUploadController;
use App\Http\Controllers\FinanceUploadController;
use App\Http\Controllers\InventoryUploadController;
use App\Http\Controllers\HRUploadController;
use App\Http\Controllers\SupplyChainUploadController;
use App\Http\Controllers\NPDUploadController;
use App\Http\Controllers\Admin\DemoDataController;

/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/

Route::get('/', function () {
    return redirect()->route('admin.dashboard');
});

Route::prefix('admin')->middleware(['auth'])->group(function () {
    // Admin Dashboard
    Route::get('/dashboard', function () {
        return view('admin.dashboard');
    })->name('admin.dashboard');

    // Module-wise File Upload Routes - Sales
    Route::prefix('upload/sales')->name('admin.upload.sales.')->group(function () {
        Route::get('/', [SalesUploadController::class, 'index'])->name('index');
        Route::get('/create', [SalesUploadController::class, 'create'])->name('create');
        Route::post('/store', [SalesUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [SalesUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [SalesUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - Production
    Route::prefix('upload/production')->name('admin.upload.production.')->group(function () {
        Route::get('/', [ProductionUploadController::class, 'index'])->name('index');
        Route::get('/create', [ProductionUploadController::class, 'create'])->name('create');
        Route::post('/store', [ProductionUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [ProductionUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [ProductionUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - Finance
    Route::prefix('upload/finance')->name('admin.upload.finance.')->group(function () {
        Route::get('/', [FinanceUploadController::class, 'index'])->name('index');
        Route::get('/create', [FinanceUploadController::class, 'create'])->name('create');
        Route::post('/store', [FinanceUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [FinanceUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [FinanceUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - Inventory
    Route::prefix('upload/inventory')->name('admin.upload.inventory.')->group(function () {
        Route::get('/', [InventoryUploadController::class, 'index'])->name('index');
        Route::get('/create', [InventoryUploadController::class, 'create'])->name('create');
        Route::post('/store', [InventoryUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [InventoryUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [InventoryUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - HR
    Route::prefix('upload/hr')->name('admin.upload.hr.')->group(function () {
        Route::get('/', [HRUploadController::class, 'index'])->name('index');
        Route::get('/create', [HRUploadController::class, 'create'])->name('create');
        Route::post('/store', [HRUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [HRUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [HRUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - Supply Chain
    Route::prefix('upload/supplychain')->name('admin.upload.supplychain.')->group(function () {
        Route::get('/', [SupplyChainUploadController::class, 'index'])->name('index');
        Route::get('/create', [SupplyChainUploadController::class, 'create'])->name('create');
        Route::post('/store', [SupplyChainUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [SupplyChainUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [SupplyChainUploadController::class, 'downloadSample'])->name('sample');
    });

    // Module-wise File Upload Routes - NPD
    Route::prefix('upload/npd')->name('admin.upload.npd.')->group(function () {
        Route::get('/', [NPDUploadController::class, 'index'])->name('index');
        Route::get('/create', [NPDUploadController::class, 'create'])->name('create');
        Route::post('/store', [NPDUploadController::class, 'store'])->name('store');
        Route::get('/{id}', [NPDUploadController::class, 'show'])->name('show');
        Route::get('/sample/{dataType}', [NPDUploadController::class, 'downloadSample'])->name('sample');
    });

    // External API Integration
    Route::prefix('external-apis')->group(function () {
        Route::get('specs', [ExternalApiController::class, 'specs'])->name('admin.external-apis.specs');
        Route::get('', [ExternalApiController::class, 'index'])->name('admin.external-apis.index');
        Route::get('create', [ExternalApiController::class, 'create'])->name('admin.external-apis.create');
        Route::post('', [ExternalApiController::class, 'store'])->name('admin.external-apis.store');
        Route::get('{externalApi}/edit', [ExternalApiController::class, 'edit'])->name('admin.external-apis.edit');
        Route::put('{externalApi}', [ExternalApiController::class, 'update'])->name('admin.external-apis.update');
        Route::delete('{externalApi}', [ExternalApiController::class, 'destroy'])->name('admin.external-apis.destroy');
        Route::post('{externalApi}/sync', [ExternalApiController::class, 'syncData'])->name('admin.external-apis.sync');
        Route::get('{externalApi}/sync-history', [ExternalApiController::class, 'syncHistory'])->name('admin.external-apis.sync-history');
        Route::get('{externalApi}/test', [ExternalApiController::class, 'testApi'])->name('admin.external-apis.test');
    });

    // Database Data Viewer
    Route::get('/data-viewer', [DataViewController::class, 'index'])->name('admin.data.viewer');
    Route::get('/data-viewer/{table}', [DataViewController::class, 'show'])->name('admin.data.viewer.show');

    // Demo Data Generator
    Route::get('/demo-data', [DemoDataController::class, 'index'])->name('admin.demo-data.index');
    Route::post('/demo-data', [DemoDataController::class, 'store'])->name('admin.demo-data.store');
});

// Authentication Routes (Laravel Breeze/UI would typically handle these)
require __DIR__.'/auth.php';
