@extends('layouts.admin')

@section('title', 'Upload HR Data')
@section('page-title', 'Upload HR Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header d-flex justify-content-between align-items-center">
                <h4>Upload HR Data File</h4>
                <a href="{{ route('admin.upload.hr.index') }}" class="btn btn-sm btn-outline-secondary">
                    <i class="bi bi-clock-history"></i> View History
                </a>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.hr.store') }}" method="POST" enctype="multipart/form-data">
                    @csrf

                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <optgroup label="Snapshots">
                                <option value="headcount_snapshot">Headcount Snapshot (totals per category)</option>
                                <option value="attendance_summary">Attendance Summary (present/absent/leave)</option>
                                <option value="movement_summary">Movement & Turnover (monthly)</option>
                                <option value="promotion_summary">Promotion Summary (yearly)</option>
                            </optgroup>
                            <optgroup label="Legacy Aliases">
                                <option value="employee_basic_info">Employee Basic Info (legacy format)</option>
                                <option value="employee_attendance">Employee Attendance (legacy)</option>
                                <option value="employee_tran_over">Employee Transfer/Overtime (legacy)</option>
                                <option value="employee_promotions">Employee Promotions (legacy)</option>
                            </optgroup>
                        </select>
                        <small class="form-text text-muted">Snapshots feed the new AI enrichment layer. Legacy options remain for backward compatibility.</small>
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB.</small>
                    </div>

                    <div class="alert alert-info" role="alert">
                        <strong>Need a template?</strong> Download ready-made sample files from the <a href="{{ route('admin.upload.hr.index') }}#samples">HR upload samples</a> section.
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

