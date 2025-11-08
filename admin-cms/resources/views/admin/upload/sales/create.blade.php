@extends('layouts.admin')

@section('title', 'Upload Sales Data')
@section('page-title', 'Upload Sales Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header">
                <h4>Upload Sales Data File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.sales.store') }}" method="POST" enctype="multipart/form-data" id="uploadForm">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <option value="monthly_report">Monthly Report</option>
                            <option value="daily_report">Daily Report</option>
                            <option value="best_selling">Best Selling Products</option>
                            <option value="top_distributors">Top Distributors</option>
                            <option value="top_retailers">Top Retailers</option>
                            <option value="order_delivery">Order vs Delivery</option>
                        </select>
                        @error('data_type')
                            <div class="text-danger">{{ $message }}</div>
                        @enderror
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">
                            Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB
                        </small>
                        @error('file')
                            <div class="text-danger">{{ $message }}</div>
                        @enderror
                    </div>

                    <div class="mb-3">
                        <a href="#" class="btn btn-outline-secondary" id="sampleLink" style="display: none;">
                            <i class="bi bi-download"></i> Download Sample File
                        </a>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.upload.sales.index') }}" class="btn btn-secondary">Cancel</a>
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
                <h5>Upload Instructions</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled">
                    <li><i class="bi bi-check-circle text-success"></i> Select the appropriate data type</li>
                    <li><i class="bi bi-check-circle text-success"></i> Download sample file for reference</li>
                    <li><i class="bi bi-check-circle text-success"></i> Ensure file format matches sample</li>
                    <li><i class="bi bi-check-circle text-success"></i> Maximum file size: 10MB</li>
                    <li><i class="bi bi-check-circle text-success"></i> Supported formats: Excel, CSV</li>
                </ul>
            </div>
        </div>

        <div class="card mt-3">
            <div class="card-header">
                <h5>Data Types</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled">
                    <li><strong>Monthly Report:</strong> Channel-wise monthly sales data</li>
                    <li><strong>Daily Report:</strong> Daily sales tracking data</li>
                    <li><strong>Best Selling:</strong> Top selling products data</li>
                    <li><strong>Top Distributors:</strong> Distributor performance data</li>
                    <li><strong>Top Retailers:</strong> Retailer performance data</li>
                    <li><strong>Order vs Delivery:</strong> Order and delivery comparison</li>
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
        sampleLink.href = '{{ route("admin.upload.sales.sample", ":type") }}'.replace(':type', dataType);
        sampleLink.style.display = 'inline-block';
    } else {
        sampleLink.style.display = 'none';
    }
});

document.getElementById('uploadForm').addEventListener('submit', function(e) {
    const fileInput = document.getElementById('file');
    const file = fileInput.files[0];
    
    if (file && file.size > 10 * 1024 * 1024) {
        e.preventDefault();
        alert('File size must be less than 10MB');
        return false;
    }
    
    document.getElementById('submitBtn').disabled = true;
    document.getElementById('submitBtn').innerHTML = '<span class="spinner-border spinner-border-sm"></span> Uploading...';
});
</script>
@endpush

