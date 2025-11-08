<?php

namespace App\Http\Controllers\Admin;

use App\Http\Controllers\Controller;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Artisan;
use Illuminate\Support\Facades\Log;

class DemoDataController extends Controller
{
    public function index()
    {
        return view('admin.demo-data.index');
    }

    public function store(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'years' => 'nullable|integer|min:1|max:2',
        ]);

        $years = $validated['years'] ?? 2;

        try {
            Artisan::call('idash:generate-demo-data', [
                '--years' => $years,
            ]);

            $output = Artisan::output();

            return redirect()
                ->route('admin.demo-data.index')
                ->with('success', 'Demo data generation completed successfully.')
                ->with('command_output', $output);
        } catch (\Throwable $e) {
            Log::error('Demo data generation failed', ['error' => $e->getMessage()]);

            return redirect()
                ->route('admin.demo-data.index')
                ->with('error', 'Failed to generate demo data. Check logs for details.');
        }
    }
}
