<?php

namespace App\Http\Controllers;

use App\Imports\HRDataImport;
use App\Models\FileUpload;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Maatwebsite\Excel\Facades\Excel;

class HRUploadController extends Controller
{
    private const CANONICAL_TYPES = [
        'headcount_snapshot',
        'attendance_summary',
        'movement_summary',
        'promotion_summary',
    ];

    private const TYPE_ALIASES = [
        'employee_basic_info' => 'headcount_snapshot',
        'employee_attendance' => 'attendance_summary',
        'employee_tran_over' => 'movement_summary',
        'employee_promotions' => 'promotion_summary',
    ];

    public function index()
    {
        $uploads = FileUpload::where('module', 'hr')
            ->orderBy('created_at', 'desc')
            ->paginate(20);

        return view('admin.upload.hr.index', compact('uploads'));
    }

    public function create()
    {
        return view('admin.upload.hr.create');
    }

    public function store(Request $request)
    {
        $accepted = array_merge(self::CANONICAL_TYPES, array_keys(self::TYPE_ALIASES));

        $request->validate([
            'file' => 'required|mimes:xlsx,xls,csv|max:10240',
            'data_type' => 'required|in:' . implode(',', $accepted),
        ]);

        try {
            $file = $request->file('file');
            $selectedType = $request->input('data_type');
            $dataType = $this->normalizeType($selectedType);
            $path = $file->store('uploads/hr/' . $dataType);

            $upload = FileUpload::create([
                'file_name' => $file->getClientOriginalName(),
                'file_path' => $path,
                'module' => 'hr',
                'data_type' => $dataType,
                'uploaded_by' => auth()->id(),
                'status' => 'processing',
            ]);

            Excel::import(new HRDataImport($upload, $dataType), $file);

            $upload->update(['status' => 'completed']);

            return redirect()->route('admin.upload.hr.index')
                ->with('success', 'HR data uploaded and processed successfully!');
        } catch (\Throwable $e) {
            Log::error('HR upload error: ' . $e->getMessage(), [
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

        return view('admin.upload.hr.show', compact('upload', 'data'));
    }

    public function downloadSample($dataType)
    {
        $samplePath = storage_path('app/samples/hr/' . $dataType . '_sample.xlsx');

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

    private function normalizeType(string $selected): string
    {
        return self::TYPE_ALIASES[$selected] ?? $selected;
    }
}

