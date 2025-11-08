<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ExternalApi extends Model
{
    use HasFactory;

    protected $fillable = [
        'name',
        'module',
        'url',
        'method',
        'headers',
        'query_params',
        'body',
        'authentication',
        'data_type',
        'items_path',
        'data_mapping',
        'is_active',
        'sync_interval',
        'last_sync_at',
    ];

    protected $casts = [
        'headers' => 'array',
        'query_params' => 'array',
        'body' => 'array',
        'authentication' => 'array',
        'data_mapping' => 'array',
        'is_active' => 'boolean',
        'last_sync_at' => 'datetime',
        'created_at' => 'datetime',
        'updated_at' => 'datetime',
    ];

    public function syncs()
    {
        return $this->hasMany(ApiSync::class, 'external_api_id');
    }

    public function scopeActive($query)
    {
        return $query->where('is_active', true);
    }
}

