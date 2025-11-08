@extends('layouts.admin')

@section('title', 'File Upload Management')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <h2>File Upload Management</h2>
            <a href="{{ route('admin.file-upload.create') }}" class="btn btn-primary">
                <i class="bi bi-upload"></i> Upload New File
            </a>
        </div>

        <div class="card">
            <div class="card-body">
                <table class="table table-striped">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>File Name</th>
                            <th>File Type</th>
                            <th>Status</th>
                            <th>Records Processed</th>
                            <th>Uploaded By</th>
                            <th>Uploaded At</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        @forelse($uploads as $upload)
                            <tr>
                                <td>{{ $upload->id }}</td>
                                <td>{{ $upload->file_name }}</td>
                                <td>
                                    <span class="badge bg-info">{{ $upload->file_type }}</span>
                                </td>
                                <td>
                                    @if($upload->status == 'completed')
                                        <span class="badge bg-success">Completed</span>
                                    @elseif($upload->status == 'processing')
                                        <span class="badge bg-warning">Processing</span>
                                    @elseif($upload->status == 'failed')
                                        <span class="badge bg-danger">Failed</span>
                                    @else
                                        <span class="badge bg-secondary">Pending</span>
                                    @endif
                                </td>
                                <td>{{ $upload->records_processed ?? 0 }}</td>
                                <td>{{ $upload->user->name ?? 'N/A' }}</td>
                                <td>{{ $upload->created_at->format('Y-m-d H:i:s') }}</td>
                                <td>
                                    <a href="{{ route('admin.file-upload.show', $upload->id) }}" class="btn btn-sm btn-info">
                                        <i class="bi bi-eye"></i> View
                                    </a>
                                    @if($upload->status == 'failed')
                                        <button class="btn btn-sm btn-warning" onclick="retryUpload({{ $upload->id }})">
                                            <i class="bi bi-arrow-clockwise"></i> Retry
                                        </button>
                                    @endif
                                </td>
                            </tr>
                        @empty
                            <tr>
                                <td colspan="8" class="text-center">No uploads found</td>
                            </tr>
                        @endforelse
                    </tbody>
                </table>

                {{ $uploads->links() }}
            </div>
        </div>
    </div>
</div>
@endsection

@push('scripts')
<script>
function retryUpload(id) {
    if(confirm('Are you sure you want to retry this upload?')) {
        // Implement retry logic
        window.location.href = '{{ route("admin.file-upload.index") }}/' + id + '/retry';
    }
}
</script>
@endpush

