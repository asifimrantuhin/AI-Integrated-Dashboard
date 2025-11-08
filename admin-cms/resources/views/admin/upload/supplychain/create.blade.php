@extends('layouts.admin')

@section('title', 'Upload Supply Chain Data')
@section('page-title', 'Upload Supply Chain Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header">
                <h4>Upload Supply Chain Data File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.supplychain.store') }}" method="POST" enctype="multipart/form-data">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <option value="po_data">PO Data</option>
                            <option value="grn_data">GRN Data</option>
                            <option value="invoice_data">Invoice Data</option>
                            <option value="raw_data">Raw Data</option>
                        </select>
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB</small>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.upload.supplychain.index') }}" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary">
                            <i class="bi bi-upload"></i> Upload File
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
@endsection

