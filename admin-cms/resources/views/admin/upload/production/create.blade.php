@extends('layouts.admin')

@section('title', 'Upload Production Data')
@section('page-title', 'Upload Production Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header">
                <h4>Upload Production Data File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.production.store') }}" method="POST" enctype="multipart/form-data" id="uploadForm">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <option value="production_analysis">Production Analysis</option>
                            <option value="wastage_data">Wastage Data</option>
                            <option value="cost_analysis">Cost Analysis</option>
                            <option value="line_efficiency">Line Efficiency (Output)</option>
                            <option value="maintenance_events">Maintenance Events</option>
                        </select>
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB</small>
                    </div>

                    <div class="mb-3">
                        <a href="#" class="btn btn-outline-secondary" id="sampleLink" style="display: none;">
                            <i class="bi bi-download"></i> Download Sample File
                        </a>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.upload.production.index') }}" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary" id="submitBtn">
                            <i class="bi bi-upload"></i> Upload File
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>

    <div class="col-md-4">
        <div class="card">
            <div class="card-header">
                <h5>Data Types</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled">
                    <li><strong>Production Analysis:</strong> Monthly production analysis data</li>
                    <li><strong>Wastage Data:</strong> Production wastage tracking</li>
                    <li><strong>Cost Analysis:</strong> Production cost analysis</li>
                    <li><strong>Line Efficiency:</strong> Daily/shift throughput, downtime, OEE per production line</li>
                    <li><strong>Maintenance Events:</strong> Machine maintenance logs with cost and downtime</li>
                </ul>
            </div>
        </div>
    </div>
</div>
@endsection

@push('scripts')
<script>
document.getElementById('data_type').addEventListener('change', function() {
    const dataType = this.value;
    const sampleLink = document.getElementById('sampleLink');
    if (dataType) {
        sampleLink.href = '{{ route("admin.upload.production.sample", ":type") }}'.replace(':type', dataType);
        sampleLink.style.display = 'inline-block';
    }
});
</script>
@endpush

