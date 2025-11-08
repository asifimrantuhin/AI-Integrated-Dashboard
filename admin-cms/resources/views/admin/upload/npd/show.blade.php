@extends('layouts.admin')

@section('title', 'NPD Upload Details')
@section('page-title', 'NPD Upload Details')

@section('content')
<div class="row">
    <div class="col-12">
        <div class="card mb-4">
            <div class="card-header">
                <h4>Upload Details</h4>
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-md-6">
                        <p><strong>File Name:</strong> {{ $upload->file_name }}</p>
                        <p><strong>Data Type:</strong> {{ ucfirst(str_replace('_', ' ', $upload->data_type)) }}</p>
                        <p><strong>Status:</strong> 
                            @if($upload->status == 'completed')
                                <span class="badge bg-success">Completed</span>
                            @elseif($upload->status == 'failed')
                                <span class="badge bg-danger">Failed</span>
                            @else
                                <span class="badge bg-secondary">{{ ucfirst($upload->status) }}</span>
                            @endif
                        </p>
                        <p><strong>Records Processed:</strong> {{ $upload->records_processed ?? 0 }}</p>
                    </div>
                    <div class="col-md-6">
                        <p><strong>Uploaded By:</strong> {{ $upload->user->name ?? 'N/A' }}</p>
                        <p><strong>Uploaded At:</strong> {{ $upload->created_at->format('Y-m-d H:i:s') }}</p>
                        @if($upload->error_message)
                            <p><strong>Error:</strong> <span class="text-danger">{{ $upload->error_message }}</span></p>
                        @endif
                    </div>
                </div>
            </div>
        </div>

        @if(count($data) > 0)
            <div class="card">
                <div class="card-header">
                    <h4>Uploaded Data Preview</h4>
                </div>
                <div class="card-body">
                    <div class="table-responsive">
                        <table class="table table-striped table-bordered">
                            <thead>
                                <tr>
                                    @if(count($data) > 0 && is_array($data[0]))
                                        @foreach(array_keys($data[0]) as $key)
                                            <th>{{ ucfirst(str_replace('_', ' ', $key)) }}</th>
                                        @endforeach
                                    @endif
                                </tr>
                            </thead>
                            <tbody>
                                @foreach($data as $row)
                                    <tr>
                                        @if(is_array($row))
                                            @foreach($row as $value)
                                                <td>{{ $value }}</td>
                                            @endforeach
                                        @endif
                                    </tr>
                                @endforeach
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        @endif

        <div class="mt-3">
            <a href="{{ route('admin.upload.npd.index') }}" class="btn btn-secondary">
                <i class="bi bi-arrow-left"></i> Back to List
            </a>
        </div>
    </div>
</div>
@endsection

