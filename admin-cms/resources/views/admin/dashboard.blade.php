@extends('layouts.admin')

@section('title', 'Admin Dashboard')
@section('page-title', 'Admin Dashboard')

@section('content')
<div class="row">
    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-cart" style="font-size: 3rem; color: #0d6efd;"></i>
                <h5 class="card-title mt-3">Sales</h5>
                <p class="card-text">Upload and manage sales data</p>
                <a href="{{ route('admin.upload.sales.index') }}" class="btn btn-primary">
                    Manage Sales
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-gear" style="font-size: 3rem; color: #198754;"></i>
                <h5 class="card-title mt-3">Production</h5>
                <p class="card-text">Upload and manage production data</p>
                <a href="{{ route('admin.upload.production.index') }}" class="btn btn-success">
                    Manage Production
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-currency-dollar" style="font-size: 3rem; color: #dc3545;"></i>
                <h5 class="card-title mt-3">Finance</h5>
                <p class="card-text">Upload and manage finance data</p>
                <a href="{{ route('admin.upload.finance.index') }}" class="btn btn-danger">
                    Manage Finance
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-box" style="font-size: 3rem; color: #ffc107;"></i>
                <h5 class="card-title mt-3">Inventory</h5>
                <p class="card-text">Upload and manage inventory data</p>
                <a href="{{ route('admin.upload.inventory.index') }}" class="btn btn-warning">
                    Manage Inventory
                </a>
            </div>
        </div>
    </div>
</div>

<div class="row mt-4">
    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-people" style="font-size: 3rem; color: #6f42c1;"></i>
                <h5 class="card-title mt-3">HR</h5>
                <p class="card-text">Upload and manage HR data</p>
                <a href="{{ route('admin.upload.hr.index') }}" class="btn btn-secondary" style="background-color: #6f42c1; border-color: #6f42c1;">
                    Manage HR
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-truck" style="font-size: 3rem; color: #20c997;"></i>
                <h5 class="card-title mt-3">Supply Chain</h5>
                <p class="card-text">Upload and manage supply chain data</p>
                <a href="{{ route('admin.upload.supplychain.index') }}" class="btn btn-info">
                    Manage Supply Chain
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-lightbulb" style="font-size: 3rem; color: #fd7e14;"></i>
                <h5 class="card-title mt-3">NPD</h5>
                <p class="card-text">Upload and manage NPD data</p>
                <a href="{{ route('admin.upload.npd.index') }}" class="btn btn-warning" style="background-color: #fd7e14; border-color: #fd7e14;">
                    Manage NPD
                </a>
            </div>
        </div>
    </div>

    <div class="col-md-3">
        <div class="card text-center">
            <div class="card-body">
                <i class="bi bi-link-45deg" style="font-size: 3rem; color: #0dcaf0;"></i>
                <h5 class="card-title mt-3">External API</h5>
                <p class="card-text">Configure and sync external APIs</p>
                <a href="{{ route('admin.external-apis.index') }}" class="btn btn-info">
                    Manage APIs
                </a>
            </div>
        </div>
    </div>
</div>

<div class="row mt-4">
    <div class="col-12">
        <div class="card">
            <div class="card-header">
                <h5>Quick Actions</h5>
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-2">
                        <a href="{{ route('admin.upload.sales.create') }}" class="btn btn-outline-primary w-100 mb-2">
                            <i class="bi bi-upload"></i> Upload Sales
                        </a>
                    </div>
                    <div class="col-md-2">
                        <a href="{{ route('admin.upload.production.create') }}" class="btn btn-outline-success w-100 mb-2">
                            <i class="bi bi-upload"></i> Upload Production
                        </a>
                    </div>
                    <div class="col-md-2">
                        <a href="{{ route('admin.upload.finance.create') }}" class="btn btn-outline-danger w-100 mb-2">
                            <i class="bi bi-upload"></i> Upload Finance
                        </a>
                    </div>
                    <div class="col-md-2">
                        <a href="{{ route('admin.upload.inventory.create') }}" class="btn btn-outline-warning w-100 mb-2">
                            <i class="bi bi-upload"></i> Upload Inventory
                        </a>
                    </div>
                    <div class="col-md-2">
                        <a href="{{ route('admin.upload.hr.create') }}" class="btn btn-outline-secondary w-100 mb-2">
                            <i class="bi bi-upload"></i> Upload HR
                        </a>
                    </div>
                    <div class="col-md-2">
                        <a href="{{ route('admin.data.viewer') }}" class="btn btn-outline-info w-100 mb-2">
                            <i class="bi bi-database"></i> View Data
                        </a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
@endsection
