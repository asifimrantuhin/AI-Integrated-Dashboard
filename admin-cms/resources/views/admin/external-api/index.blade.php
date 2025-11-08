@extends('layouts.admin')

@section('title', 'External API Management')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <div>
                <h2>External API Management</h2>
                <p class="text-muted mb-0"> Configure data feeds per module and monitor automated synchronisation.</p>
            </div>
            <div class="d-flex gap-2">
                <a href="{{ route('admin.external-apis.specs') }}" class="btn btn-outline-secondary">
                    <i class="bi bi-journal-text"></i> API Format Guide
                </a>
                <a href="{{ route('admin.external-apis.create') }}" class="btn btn-primary">
                    <i class="bi bi-plus-circle"></i> Add New API
                </a>
            </div>
        </div>

        <div class="card">
            <div class="card-body p-0">
                <div class="table-responsive">
                    <table class="table table-striped mb-0">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Module</th>
                                <th>Endpoint</th>
                                <th>Method</th>
                                <th>Status</th>
                                <th>Sync Interval</th>
                                <th>Last Sync</th>
                                <th class="text-end">Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            @forelse($apis as $api)
                                @php
                                    $lastSync = $api->syncs->first();
                                @endphp
                                <tr>
                                    <td class="align-middle">
                                        <strong>{{ $api->name }}</strong>
                                    </td>
                                    <td class="align-middle">
                                        <span class="badge bg-light text-dark border">{{ strtoupper($api->module) }}</span>
                                    </td>
                                    <td class="align-middle">
                                        <small class="text-break">{{ $api->url }}</small>
                                    </td>
                                    <td class="align-middle"><span class="badge bg-info">{{ $api->method }}</span></td>
                                    <td class="align-middle">
                                        @if($api->is_active)
                                            <span class="badge bg-success">Active</span>
                                        @else
                                            <span class="badge bg-secondary">Paused</span>
                                        @endif
                                    </td>
                                    <td class="align-middle">
                                        Every {{ $api->sync_interval }} min
                                    </td>
                                    <td class="align-middle">
                                        {{ optional($api->last_sync_at)->format('Y-m-d H:i') ?? 'Never' }}
                                    </td>
                                    <td class="align-middle text-end">
                                        <div class="d-flex flex-wrap gap-1 justify-content-end">
                                            <a href="{{ route('admin.external-apis.edit', $api) }}" class="btn btn-sm btn-outline-secondary" title="Edit">
                                                <i class="bi bi-pencil-square"></i>
                                            </a>
                                            <a href="{{ route('admin.external-apis.sync-history', $api) }}" class="btn btn-sm btn-outline-secondary" title="History">
                                                <i class="bi bi-clock-history"></i>
                                            </a>
                                            <a href="{{ route('admin.external-apis.test', $api) }}" target="_blank" class="btn btn-sm btn-outline-info" title="Test (opens JSON)">
                                                <i class="bi bi-bug"></i>
                                            </a>
                                            <form action="{{ route('admin.external-apis.sync', $api) }}" method="POST" class="d-inline" onsubmit="return confirm('Run manual sync for {{ $api->name }}?');">
                                                @csrf
                                                <button type="submit" class="btn btn-sm btn-outline-primary" title="Run Sync Now">
                                                    <i class="bi bi-arrow-repeat"></i>
                                                </button>
                                            </form>
                                            <form action="{{ route('admin.external-apis.destroy', $api) }}" method="POST" class="d-inline" onsubmit="return confirm('Delete this API configuration?');">
                                                @csrf
                                                @method('DELETE')
                                                <button type="submit" class="btn btn-sm btn-outline-danger" title="Delete">
                                                    <i class="bi bi-trash"></i>
                                                </button>
                                            </form>
                                        </div>
                                        <div class="small text-muted mt-1">
                                            Last run: {{ $lastSync?->completed_at?->diffForHumans() ?? 'N/A' }}
                                        </div>
                                    </td>
                                </tr>
                            @empty
                                <tr>
                                    <td colspan="8" class="text-center py-4">No external APIs configured yet.</td>
                                </tr>
                            @endforelse
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</div>
@endsection

