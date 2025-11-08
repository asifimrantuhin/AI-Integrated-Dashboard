@extends('layouts.admin')

@section('title', 'Database Data - ' . ucfirst($table))

@section('content')
<div class="row">
    <div class="col-12">
        <div class="d-flex justify-content-between align-items-center mb-4">
            <h2>Database Data - {{ $tableLabel ?? ucfirst(str_replace('_', ' ', $table)) }}</h2>
            <div>
                <a href="{{ route('admin.data.viewer') }}" class="btn btn-secondary">
                    <i class="bi bi-arrow-left"></i> Back to Table List
                </a>
                <a href="{{ route('admin.dashboard') }}" class="btn btn-outline-secondary">
                    <i class="bi bi-house"></i> Dashboard
                </a>
            </div>
        </div>

        <div class="card">
            <div class="card-body">
                @if(count($columns) > 0)
                    <div class="table-responsive">
                        <table class="table table-striped table-bordered">
                            <thead>
                                <tr>
                                    @foreach($columns as $column)
                                        <th>{{ ucfirst(str_replace('_', ' ', $column)) }}</th>
                                    @endforeach
                                </tr>
                            </thead>
                            <tbody>
                                @forelse($data as $row)
                                    <tr>
                                        @foreach($columns as $column)
                                            <td>
                                                @if(isset($row->$column))
                                                    @if(is_object($row->$column) || is_array($row->$column))
                                                        {{ json_encode($row->$column) }}
                                                    @else
                                                        {{ $row->$column }}
                                                    @endif
                                                @else
                                                    -
                                                @endif
                                            </td>
                                        @endforeach
                                    </tr>
                                @empty
                                    <tr>
                                        <td colspan="{{ count($columns) }}" class="text-center">No data found</td>
                                    </tr>
                                @endforelse
                            </tbody>
                        </table>
                    </div>

                    {{ $data->links() }}
                @else
                    <div class="alert alert-warning">
                        No columns found for this table
                    </div>
                @endif
            </div>
        </div>
    </div>
</div>
@endsection

