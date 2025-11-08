@extends('layouts.admin')

@section('title', 'Edit External API')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="card">
            <div class="card-header d-flex justify-content-between align-items-center">
                <h4 class="mb-0">Edit External API - {{ $api->name }}</h4>
                <a href="{{ route('admin.external-apis.specs') }}" class="btn btn-sm btn-outline-secondary">
                    <i class="bi bi-journal-text"></i> Module Format Guide
                </a>
            </div>
            <div class="card-body">
                <form action="{{ route('admin.external-apis.update', $api) }}" method="POST" id="apiForm">
                    @csrf
                    @method('PUT')
                    <div class="row g-3">
                        <div class="col-md-6">
                            <label for="name" class="form-label">API Name <span class="text-danger">*</span></label>
                            <input type="text" name="name" id="name" class="form-control" value="{{ old('name', $api->name) }}" required>
                            @error('name')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-6">
                            <label for="module" class="form-label">Module <span class="text-danger">*</span></label>
                            <select name="module" id="module" class="form-select" required>
                                @foreach($modules as $key => $spec)
                                    <option value="{{ $key }}" {{ old('module', $api->module) === $key ? 'selected' : '' }}>{{ strtoupper($key) }} - {{ $spec['label'] }}</option>
                                @endforeach
                            </select>
                            @error('module')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-8">
                            <label for="url" class="form-label">API Endpoint URL <span class="text-danger">*</span></label>
                            <input type="url" name="url" id="url" class="form-control" value="{{ old('url', $api->url) }}" required>
                            @error('url')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-2">
                            <label for="method" class="form-label">HTTP Method <span class="text-danger">*</span></label>
                            <select name="method" id="method" class="form-select" required>
                                @foreach(['GET','POST','PUT','PATCH','DELETE'] as $method)
                                    <option value="{{ $method }}" {{ old('method', $api->method) === $method ? 'selected' : '' }}>{{ $method }}</option>
                                @endforeach
                            </select>
                            @error('method')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-2">
                            <label for="data_type" class="form-label">Response Type <span class="text-danger">*</span></label>
                            <select name="data_type" id="data_type" class="form-select" required>
                                <option value="json" {{ old('data_type', $api->data_type) === 'json' ? 'selected' : '' }}>JSON</option>
                                <option value="xml" {{ old('data_type', $api->data_type) === 'xml' ? 'selected' : '' }}>XML</option>
                            </select>
                            @error('data_type')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-4">
                            <label class="form-label">Authentication</label>
                            <div class="input-group">
                                <select name="auth_type" id="auth_type" class="form-select">
                                    @php $auth = $api->authentication ?? []; @endphp
                                    <option value="none">None</option>
                                    <option value="bearer" {{ old('auth_type', $auth['type'] ?? null) === 'bearer' ? 'selected' : '' }}>Bearer Token</option>
                                    <option value="api_key" {{ old('auth_type', $auth['type'] ?? null) === 'api_key' ? 'selected' : '' }}>API Key</option>
                                    <option value="basic" {{ old('auth_type', $auth['type'] ?? null) === 'basic' ? 'selected' : '' }}>Basic</option>
                                </select>
                                <input type="text" name="auth_token" id="auth_token" class="form-control" placeholder="token / key" value="{{ old('auth_token', $auth['token'] ?? null) }}">
                            </div>
                            @error('auth_type')<div class="text-danger">{{ $message }}</div>@enderror
                            @error('auth_token')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-2" id="authHeaderWrapper" style="display: none;">
                            <label for="auth_header" class="form-label">API Key Header</label>
                            <input type="text" name="auth_header" id="auth_header" class="form-control" value="{{ old('auth_header', $auth['header'] ?? null) }}">
                            @error('auth_header')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-2">
                            <label for="sync_interval" class="form-label">Sync Interval (min)</label>
                            <select name="sync_interval" id="sync_interval" class="form-select">
                                @foreach([15, 30, 60, 120, 240, 1440] as $interval)
                                    <option value="{{ $interval }}" {{ old('sync_interval', $api->sync_interval) == $interval ? 'selected' : '' }}>Every {{ $interval }}</option>
                                @endforeach
                            </select>
                            @error('sync_interval')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-2 d-flex align-items-end">
                            <div class="form-check form-switch">
                                <input class="form-check-input" type="checkbox" role="switch" id="is_active" name="is_active" value="1" {{ old('is_active', $api->is_active) ? 'checked' : '' }}>
                                <label class="form-check-label" for="is_active">Active</label>
                            </div>
                        </div>
                        <div class="col-12">
                            <label for="headers" class="form-label">Custom Headers (JSON)</label>
                            <textarea name="headers" id="headers" class="form-control font-monospace" rows="3">{{ old('headers', $api->headers ? json_encode($api->headers, JSON_PRETTY_PRINT) : null) }}</textarea>
                            @error('headers')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-6">
                            <label for="query_params" class="form-label">Query Parameters (JSON)</label>
                            <textarea name="query_params" id="query_params" class="form-control font-monospace" rows="3">{{ old('query_params', $api->query_params ? json_encode($api->query_params, JSON_PRETTY_PRINT) : null) }}</textarea>
                            @error('query_params')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-6">
                            <label for="body" class="form-label">Request Body (JSON)</label>
                            <textarea name="body" id="body" class="form-control font-monospace" rows="3">{{ old('body', $api->body ? json_encode($api->body, JSON_PRETTY_PRINT) : null) }}</textarea>
                            @error('body')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-md-6">
                            <label for="items_path" class="form-label">Items Path (dot notation)</label>
                            <input type="text" name="items_path" id="items_path" class="form-control" value="{{ old('items_path', $api->items_path) }}">
                            @error('items_path')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                        <div class="col-12">
                            <label for="data_mapping" class="form-label">Data Mapping (JSON) <span class="text-danger">*</span></label>
                            <textarea name="data_mapping" id="data_mapping" class="form-control font-monospace" rows="10" required>{{ old('data_mapping', $api->data_mapping ? json_encode($api->data_mapping, JSON_PRETTY_PRINT) : null) }}</textarea>
                            @error('data_mapping')<div class="text-danger">{{ $message }}</div>@enderror
                        </div>
                    </div>

                    <div class="d-flex justify-content-between mt-4">
                        <a href="{{ route('admin.external-apis.index') }}" class="btn btn-secondary">Cancel</a>
                        <button type="submit" class="btn btn-primary">
                            <i class="bi bi-save"></i> Update API
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
@endsection

@push('scripts')
<script>
const moduleSpecs = @json($modules);
const moduleSelect = document.getElementById('module');
const itemsPathField = document.getElementById('items_path');
const queryParamsField = document.getElementById('query_params');
const bodyField = document.getElementById('body');
const mappingField = document.getElementById('data_mapping');
const authTypeField = document.getElementById('auth_type');
const authHeaderWrapper = document.getElementById('authHeaderWrapper');

function maybePopulateDefaults(moduleKey) {
    const spec = moduleSpecs[moduleKey];
    if (!spec) {
        return;
    }
    const defaults = spec.defaults || {};

    if (!itemsPathField.value && defaults.items_path) {
        itemsPathField.value = defaults.items_path;
    }
    if (!queryParamsField.value && defaults.query_params) {
        queryParamsField.value = JSON.stringify(defaults.query_params, null, 2);
    }
    if (!bodyField.value && defaults.body) {
        bodyField.value = JSON.stringify(defaults.body, null, 2);
    }
    if (!mappingField.value && defaults.data_mapping) {
        mappingField.value = JSON.stringify(defaults.data_mapping, null, 2);
    }
}

moduleSelect?.addEventListener('change', function () {
    maybePopulateDefaults(this.value);
});

authTypeField?.addEventListener('change', function () {
    if (this.value === 'api_key') {
        authHeaderWrapper.style.display = 'block';
    } else {
        authHeaderWrapper.style.display = 'none';
    }
});

if (authTypeField?.value === 'api_key') {
    authHeaderWrapper.style.display = 'block';
}

maybePopulateDefaults(moduleSelect?.value);
</script>
@endpush

