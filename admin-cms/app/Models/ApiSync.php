<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ApiSync extends Model
{
    use HasFactory;

    protected $fillable = [
        'external_api_id',
        'status',
        'records_synced',
        'response_data',
        'error_message',
        'started_at',
        'completed_at',
    ];

    protected $casts = [
        'started_at' => 'datetime',
        'completed_at' => 'datetime',
        'response_data' => 'array',
        'created_at' => 'datetime',
        'updated_at' => 'datetime',
    ];

    public function externalApi()
    {
        return $this->belongsTo(ExternalApi::class, 'external_api_id');
    }
    
    public function api()
    {
        return $this->externalApi();
    }
}

