<?php

namespace App\Http\Controllers;

use App\Models\ExternalApi;
use App\Models\ApiSync;
use App\Services\ExternalApi\SyncService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;

class ExternalApiController extends Controller
{
    public function __construct(private readonly SyncService $syncService)
    {
    }

    public function index()
    {
        $apis = ExternalApi::with(['syncs' => function ($query) {
            $query->latest()->limit(1);
        }])->orderBy('module')->orderBy('name')->get();

        return view('admin.external-api.index', [
            'apis' => $apis,
            'moduleSpecs' => config('external_apis.modules'),
        ]);
    }

    public function specs()
    {
        return view('admin.external-api.specs', [
            'modules' => config('external_apis.modules'),
            'placeholders' => config('external_apis.placeholder_help'),
        ]);
    }

    public function create()
    {
        return view('admin.external-api.create', [
            'modules' => config('external_apis.modules'),
        ]);
    }

    public function store(Request $request)
    {
        $data = $this->validatedPayload($request);
        ExternalApi::create($data);

        return redirect()->route('admin.external-apis.index')
            ->with('success', 'External API configured successfully.');
    }

    public function edit(ExternalApi $externalApi)
    {
        return view('admin.external-api.edit', [
            'api' => $externalApi,
            'modules' => config('external_apis.modules'),
        ]);
    }

    public function update(Request $request, ExternalApi $externalApi)
    {
        $data = $this->validatedPayload($request, $externalApi);
        $externalApi->update($data);

        return redirect()->route('admin.external-apis.index')
            ->with('success', 'External API updated successfully.');
    }

    public function destroy(ExternalApi $externalApi)
    {
        $externalApi->delete();

        return redirect()->route('admin.external-apis.index')
            ->with('success', 'External API deleted successfully.');
    }

    public function testApi(ExternalApi $externalApi)
    {
        try {
            $payload = $this->syncService->test($externalApi);

            return response()->json([
                'status' => 'success',
                'preview' => $payload['data'],
                'raw' => $payload['raw'],
            ]);
        } catch (\Throwable $e) {
            Log::error('External API test failed', [
                'external_api_id' => $externalApi->id,
                'error' => $e->getMessage(),
            ]);

            return response()->json([
                'status' => 'error',
                'message' => $e->getMessage(),
            ], 500);
        }
    }

    public function syncData(ExternalApi $externalApi)
    {
        try {
            $sync = $this->syncService->sync($externalApi);

            return redirect()->route('admin.external-apis.index')
                ->with('success', sprintf('Data synchronized successfully (%d records).', $sync->records_synced));
        } catch (\Throwable $e) {
            Log::error('External API manual sync failed', [
                'external_api_id' => $externalApi->id,
                'error' => $e->getMessage(),
            ]);

            return redirect()->back()->with('error', 'Error syncing data: ' . $e->getMessage());
        }
    }

    public function syncHistory(ExternalApi $externalApi)
    {
        $syncs = ApiSync::where('external_api_id', $externalApi->id)
            ->orderByDesc('created_at')
            ->paginate(20);

        return view('admin.external-api.sync-history', [
            'api' => $externalApi,
            'syncs' => $syncs,
        ]);
    }

    private function validatedPayload(Request $request, ?ExternalApi $existing = null): array
    {
        $modules = array_keys(config('external_apis.modules'));

        $data = $request->validate([
            'name' => 'required|string|max:255',
            'module' => 'required|string|in:' . implode(',', $modules),
            'url' => 'required|url',
            'method' => 'required|in:GET,POST,PUT,PATCH,DELETE',
            'headers' => 'nullable|string',
            'query_params' => 'nullable|string',
            'body' => 'nullable|string',
            'auth_type' => 'nullable|in:none,bearer,api_key,basic',
            'auth_token' => 'nullable|string',
            'auth_header' => 'nullable|string|max:255',
            'data_type' => 'required|in:json,xml',
            'items_path' => 'nullable|string',
            'data_mapping' => 'required|string',
            'sync_interval' => 'nullable|integer|min:5|max:1440',
            'is_active' => 'nullable|boolean',
        ]);

        $headers = $this->decodeJson($data['headers'] ?? null, 'headers');
        $queryParams = $this->decodeJson($data['query_params'] ?? null, 'query_params');
        $body = $this->decodeJson($data['body'] ?? null, 'body');
        $dataMapping = $this->decodeJson($data['data_mapping'] ?? null, 'data_mapping');

        $authentication = null;
        $authType = $data['auth_type'] ?? null;
        if ($authType && $authType !== 'none') {
            $authentication = array_filter([
                'type' => $authType,
                'token' => $data['auth_token'] ?? null,
                'header' => $authType === 'api_key' ? ($data['auth_header'] ?? 'X-API-Key') : null,
            ]);
        }

        return [
            'name' => $data['name'],
            'module' => $data['module'],
            'url' => $data['url'],
            'method' => $data['method'],
            'headers' => $headers ?? [],
            'query_params' => $queryParams ?? [],
            'body' => $body ?: null,
            'authentication' => $authentication,
            'data_type' => $data['data_type'],
            'items_path' => $data['items_path'] ?? null,
            'data_mapping' => $dataMapping ?? [],
            'sync_interval' => (int) ($data['sync_interval'] ?? ($existing->sync_interval ?? 60)),
            'is_active' => $request->boolean('is_active', true),
        ];
    }

    private function decodeJson(?string $json, string $field): ?array
    {
        if ($json === null || $json === '') {
            return null;
        }

        try {
            $decoded = json_decode($json, true, 512, JSON_THROW_ON_ERROR);
            return $decoded;
        } catch (\JsonException $e) {
            throw ValidationException::withMessages([
                $field => 'Invalid JSON: ' . $e->getMessage(),
            ]);
        }
    }
}
