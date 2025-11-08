@extends('layouts.admin')

@section('title', 'Upload File')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="card">
            <div class="card-header">
                <h4>Upload File</h4>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.file-upload.store') }}" method="POST" enctype="multipart/form-data" id="uploadForm">
                    @csrf
                    
                    <div class="mb-3">
                        <label for="file_type" class="form-label">File Type <span class="text-danger">*</span></label>
                        <select name="type" id="file_type" class="form-select" required>
                            <option value="">Select File Type</option>
                            <option value="sales">Sales</option>
                            <option value="production">Production</option>
                            <option value="finance">Finance</option>
                            <option value="inventory">Inventory</option>
                            <option value="hr">HR</option>
                            <option value="supplychain">Supply Chain</option>
                        </select>
                        @error('type')
                            <div class="text-danger">{{ $message }}</div>
                        @enderror
                    </div>

                    <div class="mb-3">
                        <label for="file" class="form-label">Select File <span class="text-danger">*</span></label>
                        <input type="file" name="file" id="file" class="form-control" accept=".xlsx,.xls,.csv,.json" required>
                        <small class="form-text text-muted">
                            Supported formats: Excel (.xlsx, .xls), CSV (.csv), JSON (.json). Max size: 10MB
                        </small>
                        @error('file')
                            <div class="text-danger">{{ $message }}</div>
                        @enderror
                    </div>

                    <div class="mb-3">
                        <a href="{{ route('admin.file-upload.sample', 'sales') }}" class="btn btn-outline-secondary" id="sampleLink" style="display: none;">
                            <i class="bi bi-download"></i> Download Sample File
                        </a>
                    </div>

                    <div class="mb-3">
                        <div class="progress" id="uploadProgress" style="display: none;">
                            <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: 0%"></div>
                        </div>
                    </div>

                    <div class="d-flex justify-content-between">
                        <a href="{{ route('admin.file-upload.index') }}" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary" id="submitBtn">
                            <i class="bi bi-upload"></i> Upload File
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <!-- Sample File Instructions -->
        <div class="card mt-4">
            <div class="card-header">
                <h5>Sample File Instructions</h5>
            </div>
            <div class="card-body">
                <ul>
                    <li>Download the sample file for your selected file type</li>
                    <li>Fill in the data according to the sample format</li>
                    <li>Ensure all required columns are present</li>
                    <li>Save the file and upload it</li>
                    <li>The system will validate and process the data automatically</li>
                </ul>
            </div>
        </div>
    </div>
</div>
@endsection

@push('scripts')
<script>
document.getElementById('file_type').addEventListener('change', function() {
    const fileType = this.value;
    const sampleLink = document.getElementById('sampleLink');
    if (fileType) {
        sampleLink.href = '{{ route("admin.file-upload.sample", ":type") }}'.replace(':type', fileType);
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

    // Show progress bar
    document.getElementById('uploadProgress').style.display = 'block';
    document.getElementById('submitBtn').disabled = true;
});
</script>
@endpush

