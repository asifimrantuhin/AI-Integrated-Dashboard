@extends('layouts.admin')

@section('title', 'Data Viewer')
@section('page-title', 'Database Data Viewer')

@section('content')
<div class="card">
    <div class="card-header">
        <h4>Select Table to View</h4>
    </div>
    <div class="card-body">
        <div class="row">
            @foreach($tables as $tableKey => $tableName)
                <div class="col-md-3 mb-3">
                    <a href="{{ route('admin.data.viewer.show', $tableKey) }}" class="btn btn-outline-primary w-100">
                        <i class="bi bi-table"></i> {{ $tableName }}
                    </a>
                </div>
            @endforeach
        </div>
    </div>
</div>
@endsection

