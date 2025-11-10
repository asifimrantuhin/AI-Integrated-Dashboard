<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;
use Maatwebsite\Excel\Facades\Excel;
use App\Models\FileUpload;
use App\Imports\SupplyChainDataImport;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class SupplyChainUploadController extends Controller
{
    private const DATA_TYPES = [
        'po_data',
        'grn_data',
        'invoice_data',
        'supplier_performance',
    ];

    public function index()
    {
        $uploads = FileUpload::where('module', 'supplychain')
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return view('admin.upload.supplychain.index', compact('uploads'));
    }

    public function create()
    {
        return view('admin.upload.supplychain.create');
    }

    public function store(Request $request)
    {
        $request->validate([
            'file' => 'required|mimes:xlsx,xls,csv|max:10240',
            'data_type' => 'required|in:' . implode(',', self::DATA_TYPES),
        ]);

        try {
            $file = $request->file('file');
            $dataType = $request->input('data_type');
            $path = $file->store('uploads/supplychain/' . $dataType);

            $upload = FileUpload::create([
                'file_name' => $file->getClientOriginalName(),
                'file_path' => $path,
                'module' => 'supplychain',
                'data_type' => $dataType,
                'uploaded_by' => auth()->id(),
                'status' => 'processing',
            ]);

            Excel::import(new SupplyChainDataImport($upload, $dataType), $file);

            $upload->update(['status' => 'completed']);

            return redirect()->route('admin.upload.supplychain.index')
                ->with('success', 'Supply Chain data uploaded and processed successfully!');
        } catch (\Throwable $e) {
            Log::error('Supply Chain upload error: ' . $e->getMessage(), [
                'user_id' => auth()->id(),
            ]);

            if (isset($upload)) {
                $upload->update([
                    'status' => 'failed',
                    'error_message' => $e->getMessage(),
                ]);
            }

            return redirect()->back()
                ->with('error', 'File processing failed: ' . $e->getMessage())
                ->withInput();
        }
    }

    public function show($id)
    {
        $upload = FileUpload::findOrFail($id);
        $data = $this->getUploadData($upload);

        return view('admin.upload.supplychain.show', compact('upload', 'data'));
    }

    public function downloadSample($dataType)
    {
        $samplePath = storage_path('app/samples/supplychain/' . $dataType . '_sample.xlsx');

        if (file_exists($samplePath)) {
            return response()->download($samplePath);
        }

        return redirect()->back()->with('error', 'Sample file not found');
    }

    private function getUploadData($upload)
    {
        try {
            $fullPath = storage_path('app/' . $upload->file_path);
            $extension = pathinfo($fullPath, PATHINFO_EXTENSION);
            $data = [];

            if (in_array($extension, ['xlsx', 'xls'])) {
                $excelData = Excel::toArray([], $fullPath);
                $data = $excelData[0] ?? [];
            } elseif ($extension === 'csv') {
                $data = $this->parseCsv($fullPath);
            }

            return array_slice($data, 0, 100);
        } catch (\Throwable $e) {
            return [];
        }
    }

    private function parseCsv($filePath)
    {
        $data = [];
        if (($handle = fopen($filePath, 'r')) !== false) {
            while (($row = fgetcsv($handle, 1000, ',')) !== false) {
                $data[] = $row;
            }
            fclose($handle);
        }

        return $data;
    }
}

