<?php

namespace App\Http\Controllers;

use App\Models\FileUpload;
use Illuminate\Http\Request;
use Maatwebsite\Excel\Facades\Excel;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\DB;

class FileUploadController extends Controller
{
    public function index()
    {
        $uploads = FileUpload::orderBy('created_at', 'desc')->paginate(20);
        return view('admin.file-upload.index', compact('uploads'));
    }

    public function create()
    {
        return view('admin.file-upload.create');
    }

    public function store(Request $request)
    {
        $request->validate([
            'file' => 'required|mimes:xlsx,xls,csv,json|max:10240',
            'type' => 'required|in:sales,production,finance,inventory,hr,supplychain',
        ]);

        try {
            $file = $request->file('file');
            $fileName = time() . '_' . $file->getClientOriginalName();
            $filePath = $file->storeAs('uploads', $fileName);

            // Save upload record
            $upload = FileUpload::create([
                'file_name' => $fileName,
                'file_path' => $filePath,
                'file_type' => $request->type,
                'status' => 'pending',
                'uploaded_by' => auth()->id(),
            ]);

            // Process file
            $this->processFile($upload, $filePath, $request->type);

            return redirect()->route('admin.file-upload.index')
                ->with('success', 'File uploaded successfully');
        } catch (\Exception $e) {
            return redirect()->back()
                ->with('error', 'Error uploading file: ' . $e->getMessage());
        }
    }

    private function processFile($upload, $filePath, $type)
    {
        try {
            $fullPath = storage_path('app/' . $filePath);
            $extension = pathinfo($fullPath, PATHINFO_EXTENSION);

            if (in_array($extension, ['xlsx', 'xls'])) {
                $data = Excel::toArray([], $fullPath);
                $this->importExcelData($data[0], $type, $upload->id);
            } elseif ($extension === 'csv') {
                $data = $this->parseCsv($fullPath);
                $this->importCsvData($data, $type, $upload->id);
            } elseif ($extension === 'json') {
                $data = json_decode(file_get_contents($fullPath), true);
                $this->importJsonData($data, $type, $upload->id);
            }

            $upload->update(['status' => 'completed']);
        } catch (\Exception $e) {
            $upload->update([
                'status' => 'failed',
                'error_message' => $e->getMessage()
            ]);
            throw $e;
        }
    }

    private function importExcelData($data, $type, $uploadId)
    {
        $recordsProcessed = 0;
        
        try {
            foreach ($data as $row) {
                $mappedData = $this->mapDataToModel($row, $type);
                
                if ($mappedData) {
                    $this->insertData($mappedData, $type);
                    $recordsProcessed++;
                }
            }
            
            FileUpload::where('id', $uploadId)->update([
                'records_processed' => $recordsProcessed
            ]);
        } catch (\Exception $e) {
            throw $e;
        }
    }

    private function importCsvData($data, $type, $uploadId)
    {
        $this->importExcelData($data, $type, $uploadId);
    }

    private function importJsonData($data, $type, $uploadId)
    {
        $recordsProcessed = 0;
        
        try {
            if (is_array($data)) {
                foreach ($data as $row) {
                    $mappedData = $this->mapDataToModel($row, $type);
                    
                    if ($mappedData) {
                        $this->insertData($mappedData, $type);
                        $recordsProcessed++;
                    }
                }
            }
            
            FileUpload::where('id', $uploadId)->update([
                'records_processed' => $recordsProcessed
            ]);
        } catch (\Exception $e) {
            throw $e;
        }
    }

    private function mapDataToModel($row, $type)
    {
        // Map data based on type
        // This is a simplified version - implement full mapping based on your requirements
        switch ($type) {
            case 'sales':
                return $this->mapSalesData($row);
            case 'production':
                return $this->mapProductionData($row);
            case 'finance':
                return $this->mapFinanceData($row);
            case 'inventory':
                return $this->mapInventoryData($row);
            case 'hr':
                return $this->mapHRData($row);
            case 'supplychain':
                return $this->mapSupplyChainData($row);
            default:
                return null;
        }
    }

    private function mapSalesData($row)
    {
        // Map sales data to database model
        // Implement based on your sales table structure
        return $row; // Placeholder
    }

    private function mapProductionData($row)
    {
        // Map production data
        return $row; // Placeholder
    }

    private function mapFinanceData($row)
    {
        // Map finance data
        return $row; // Placeholder
    }

    private function mapInventoryData($row)
    {
        // Map inventory data
        return $row; // Placeholder
    }

    private function mapHRData($row)
    {
        // Map HR data
        return $row; // Placeholder
    }

    private function mapSupplyChainData($row)
    {
        // Map supply chain data
        return $row; // Placeholder
    }

    private function insertData($data, $type)
    {
        // Insert data into appropriate table based on type
        // This is a placeholder - implement based on your models
        try {
            switch ($type) {
                case 'sales':
                    // DB::table('channelwise_monthly_report')->insert($data);
                    break;
                case 'production':
                    // DB::table('production_analyses')->insert($data);
                    break;
                // Add other cases
            }
        } catch (\Exception $e) {
            throw $e;
        }
    }

    private function parseCsv($filePath)
    {
        $data = [];
        if (($handle = fopen($filePath, "r")) !== false) {
            $header = fgetcsv($handle);
            while (($row = fgetcsv($handle)) !== false) {
                $data[] = array_combine($header, $row);
            }
            fclose($handle);
        }
        return $data;
    }

    public function show($id)
    {
        $upload = FileUpload::findOrFail($id);
        $data = $this->getUploadData($upload);
        return view('admin.file-upload.show', compact('upload', 'data'));
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
            } elseif ($extension === 'json') {
                $data = json_decode(file_get_contents($fullPath), true);
                if (!is_array($data)) {
                    $data = [];
                }
            }

            return array_slice($data, 0, 100); // Return first 100 rows for preview
        } catch (\Exception $e) {
            return [];
        }
    }

    public function downloadSample($type)
    {
        $samplePath = storage_path('app/samples/' . $type . '_sample.xlsx');
        if (file_exists($samplePath)) {
            return response()->download($samplePath);
        }
        return redirect()->back()->with('error', 'Sample file not found');
    }
}

