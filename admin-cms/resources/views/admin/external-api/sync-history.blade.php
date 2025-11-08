@extends('layouts.admin')

@section('title', 'Sync History')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <h2>Sync History - {{ $api->name }}</h2>
            <a href="{{ route('admin.external-apis.index') }}" class="btn btn-secondary">
                <i class="bi bi-arrow-left"></i> Back to APIs
            </a>
        </div>

        <div class="card">
            <div class="card-body">
                <table class="table table-striped">
                    <thead>
                        <tr>
                            <th>ID</th>
                            <th>Status</th>
                            <th>Records Synced</th>
                            <th>Started</th>
                            <th>Completed</th>
                            <th>Duration</th>
                            <th>Response Snippet</th>
                            <th>Error</th>
                        </tr>
                    </thead>
                    <tbody>
                        @forelse($syncs as $sync)
                            <tr>
                                <td>{{ $sync->id }}</td>
                                <td>
                                    @if($sync->status == 'completed')
                                        <span class="badge bg-success">Completed</span>
                                    @elseif($sync->status == 'processing')
                                        <span class="badge bg-warning">Processing</span>
                                    @elseif($sync->status == 'failed')
                                        <span class="badge bg-danger">Failed</span>
                                    @else
                                        <span class="badge bg-secondary">Pending</span>
                                    @endif
                                </td>
                                <td>{{ $sync->records_synced }}</td>
                                <td>{{ $sync->started_at?->format('Y-m-d H:i:s') ?? 'N/A' }}</td>
                                <td>{{ $sync->completed_at?->format('Y-m-d H:i:s') ?? 'N/A' }}</td>
                                <td>
                                    @if($sync->started_at && $sync->completed_at)
                                        {{ $sync->started_at->diffInSeconds($sync->completed_at) }}s
                                    @else
                                        N/A
                                    @endif
                                </td>
                                <td>
                                    @if($sync->response_data)
                                        <pre class="small bg-light p-2 rounded mb-0">{{ \Illuminate\Support\Str::limit(json_encode($sync->response_data), 120) }}</pre>
                                    @else
                                        <span class="text-muted">—</span>
                                    @endif
                                </td>
                                <td>
                                    @if($sync->error_message)
                                        <span class="text-danger">{{ $sync->error_message }}</span>
                                    @else
                                        <span class="text-muted">—</span>
                                    @endif
                                </td>
                            </tr>
                        @empty
                            <tr>
                                <td colspan="8" class="text-center">No sync history found</td>
                            </tr>
                        @endforelse
                    </tbody>
                </table>

                {{ $syncs->links() }}
            </div>
        </div>
    </div>
</div>
@endsection

