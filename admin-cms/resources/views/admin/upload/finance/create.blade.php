@extends('layouts.admin')

@section('title', 'Upload Finance Data')
@section('page-title', 'Upload Finance Data')

@section('content')
<div class="row">
    <div class="col-md-8">
        <div class="card">
            <div class="card-header">
                <h4>Upload Finance Data File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.upload.finance.store') }}" method="POST" enctype="multipart/form-data">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="data_type" class="form-label">Data Type <span class="text-danger">*</span></label>
                        <select name="data_type" id="data_type" class="form-select" required>
                            <option value="">Select Data Type</option>
                            <option value="bank_loan">Bank Loan</option>
                            <option value="bank_loan_status">Bank Loan Status</option>
                            <option value="budget">Budget</option>
                            <option value="expense">Expense</option>
                            <option value="financial_statement">Financial Statement</option>
                            <option value="budget_summary">Budget vs Actual Summary</option>
                            <option value="expense_summary">Expense Breakdown Summary</option>
                        </select>
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv" required>
                        <small class="form-text text-muted">Supported formats: Excel (.xlsx, .xls), CSV (.csv). Max size: 10MB</small>
                    </div>

                    <div class="mb-3">
                        <a href="#" class="btn btn-outline-secondary" id="sampleLink" style="display:none;">
                            <i class="bi bi-download"></i> Download Sample File
                        </a>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.upload.finance.index') }}" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary">
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
                <h5>Data Type Guidance</h5>
            </div>
            <div class="card-body">
                <ul class="list-unstyled">
                    <li><strong>Bank Loan:</strong> Monthly net sales & financial expense per company</li>
                    <li><strong>Bank Loan Status:</strong> Closing exposure by loan head</li>
                    <li><strong>Budget / Expense:</strong> Raw budget allocation and invoices</li>
                    <li><strong>Budget Summary:</strong> Aggregated budget vs actual per cost center</li>
                    <li><strong>Expense Summary:</strong> Monthly expense totals by category/department</li>
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
        sampleLink.href = '{{ route("admin.upload.finance.sample", ":type") }}'.replace(':type', dataType);
        sampleLink.style.display = 'inline-block';
    } else {
        sampleLink.style.display = 'none';
        sampleLink.removeAttribute('href');
    }
});
</script>
@endpush

