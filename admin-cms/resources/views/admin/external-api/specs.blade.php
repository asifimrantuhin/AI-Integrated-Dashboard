@extends('layouts.admin')

@section('title', 'External API Format Guide')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <h2>External API Format Guide</h2>
            <a href="{{ route('admin.external-apis.index') }}" class="btn btn-secondary">
                <i class="bi bi-arrow-left"></i> Back to APIs
            </a>
        </div>

        @foreach($modules as $key => $spec)
            <div class="card mb-4">
                <div class="card-header bg-light d-flex justify-content-between align-items-center">
                    <div>
                        <h5 class="mb-0">{{ strtoupper($key) }} - {{ $spec['label'] }}</h5>
                        <small class="text-muted">{{ $spec['description'] ?? '' }}</small>
                    </div>
                    <span class="badge bg-dark-subtle text-dark">Module Key: {{ $key }}</span>
                </div>
                <div class="card-body">
                    <div class="row g-4">
                        <div class="col-md-4">
                            <h6>Request Expectations</h6>
                            <ul class="mb-0">
                                <li>Method: <strong>{{ $spec['defaults']['method'] ?? 'GET' }}</strong></li>
                                <li>Response Type: <strong>JSON</strong> (XML supported when configured)</li>
                                @if(!empty($spec['defaults']['items_path']))
                                    <li>Items Path: <code>{{ $spec['defaults']['items_path'] }}</code></li>
                                @endif
                            </ul>
                            @if(!empty($spec['defaults']['query_params']))
                                <div class="mt-2">
                                    <p class="fw-semibold mb-1">Query Parameters</p>
                                    <pre class="bg-light p-2 rounded small">{{ json_encode($spec['defaults']['query_params'], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES) }}</pre>
                                </div>
                            @endif
                            @if(!empty($spec['defaults']['body']))
                                <div class="mt-2">
                                    <p class="fw-semibold mb-1">Request Body</p>
                                    <pre class="bg-light p-2 rounded small">{{ json_encode($spec['defaults']['body'], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES) }}</pre>
                                </div>
                            @endif
                        </div>
                        <div class="col-md-5">
                            <h6>Response Mapping Template</h6>
                            <pre class="bg-light p-2 rounded small">{{ json_encode($spec['defaults']['data_mapping'], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES) }}</pre>
                        </div>
                        <div class="col-md-3">
                            <h6>Sample Payload</h6>
                            <pre class="bg-light p-2 rounded small">{{ json_encode($spec['defaults']['sample_response'] ?? [], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES) }}</pre>
                        </div>
                    </div>
                </div>
            </div>
        @endforeach

        <div class="card">
            <div class="card-header">
                <h5 class="mb-0">Placeholder Reference</h5>
            </div>
            <div class="card-body">
                <div class="row g-3">
                    @foreach($placeholders as $placeholder => $description)
                        <div class="col-md-4">
                            <div class="border rounded p-3 h-100 bg-light">
                                <code>{{ $placeholder }}</code>
                                <p class="mb-0 small text-muted">{{ $description }}</p>
                            </div>
                        </div>
                    @endforeach
                </div>
            </div>
        </div>
    </div>
</div>
@endsection
