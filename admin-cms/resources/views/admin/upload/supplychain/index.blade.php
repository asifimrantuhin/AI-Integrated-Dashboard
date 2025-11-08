@extends('layouts.admin')

@section('title', 'Supply Chain Data Upload')
@section('page-title', 'Supply Chain Data Upload Management')

@section('content')
<div class="d-flex justify-content-between align-items-center mb-4">
    <h3>Supply Chain Data Uploads</h3>
    <a href="{{ route('admin.upload.supplychain.create') }}" class="btn btn-primary">
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
                    <th>Data Type</th>
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
                        <td><span class="badge bg-info">{{ ucfirst(str_replace('_', ' ', $upload->data_type)) }}</span></td>
                        <td>
                            @if($upload->status == 'completed')
                                <span class="badge bg-success">Completed</span>
                            @elseif($upload->status == 'failed')
                                <span class="badge bg-danger">Failed</span>
                            @else
                                <span class="badge bg-secondary">{{ ucfirst($upload->status) }}</span>
                            @endif
                        </td>
                        <td>{{ $upload->records_processed ?? 0 }}</td>
                        <td>{{ $upload->user->name ?? 'N/A' }}</td>
                        <td>{{ $upload->created_at->format('Y-m-d H:i:s') }}</td>
                        <td>
                            <a href="{{ route('admin.upload.supplychain.show', $upload->id) }}" class="btn btn-sm btn-info">
                                <i class="bi bi-eye"></i> View
                            </a>
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
@endsection

