@extends('layouts.admin')

@section('title', 'Demo Data Generator')
@section('page-title', 'Demo Data Generator')

@section('content')
<div class="card">
    <div class="card-body">
        <p class="text-muted">Generate sample data for all analytical modules (sales, production, finance, inventory, HR, supply chain, manufacturing) across the last one or two years. Existing records (except users) will be purged before new data is created.</p>

        <form method="POST" action="{{ route('admin.demo-data.store') }}" onsubmit="return confirm('This will erase existing module data and replace it with demo data. Continue?')">
            @csrf
            <div class="row mb-3">
                <div class="col-md-4">
                    <label class="form-label fw-semibold">Years of historical data</label>
                    <select name="years" class="form-select">
                        <option value="1">1 Year</option>
                        <option value="2" selected>2 Years</option>
                    </select>
                    <div class="form-text">Defaults to 2 years if left unchanged.</div>
                </div>
            </div>
            <button type="submit" class="btn btn-primary">Generate Demo Data</button>
        </form>

        @if(session('command_output'))
            <div class="alert alert-secondary mt-4" role="alert">
                <pre class="mb-0 small">{{ trim(session('command_output')) }}</pre>
            </div>
        @endif
    </div>
</div>
@endsection
