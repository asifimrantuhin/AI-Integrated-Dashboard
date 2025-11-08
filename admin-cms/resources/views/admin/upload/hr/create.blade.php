@extends('layouts.admin')

@section('title', 'Upload HR Data')
@section('page-title', 'Upload HR Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header">
                <h4>Upload HR Data File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.hr.store') }}" method="POST" enctype="multipart/form-data">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <option value="employee_basic_info">Employee Basic Info</option>
                            <option value="employee_attendance">Employee Attendance</option>
                            <option value="employee_promotions">Employee Promotions</option>
                            <option value="employee_tran_over">Employee Transfer/Overtime</option>
                        </select>
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB</small>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.upload.hr.index') }}" class="btn btn-secondary">Cancel</a>
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

