<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>@yield('title', 'Admin Panel') - iDash</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/css/bootstrap.min.css" rel="stylesheet">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.0/font/bootstrap-icons.css">
    <style>
        .sidebar {
            min-height: 100vh;
            background-color: #f8f9fa;
            padding: 20px 0;
        }
        .sidebar .nav-link {
            color: #333;
            padding: 10px 20px;
            border-radius: 5px;
            margin: 5px 10px;
        }
        .sidebar .nav-link:hover, .sidebar .nav-link.active {
            background-color: #007bff;
            color: white;
        }
        .module-header {
            font-weight: bold;
            color: #666;
            padding: 10px 20px;
            margin-top: 20px;
            text-transform: uppercase;
            font-size: 0.85rem;
        }
    </style>
    @stack('styles')
</head>
<body>
    <div class="container-fluid">
        <div class="row">
            <!-- Sidebar -->
            <nav class="col-md-2 sidebar">
                <div class="position-sticky pt-3">
                    <div class="text-center mb-4">
                        <h4><i class="bi bi-speedometer2"></i> iDash Admin</h4>
                    </div>
                    
                    <ul class="nav flex-column">
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.dashboard') ? 'active' : '' }}" href="{{ route('admin.dashboard') }}">
                                <i class="bi bi-house"></i> Dashboard
                            </a>
                        </li>

                        <div class="module-header">Data Upload</div>
                        
                        <!-- Sales Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.sales.*') ? 'active' : '' }}" href="{{ route('admin.upload.sales.index') }}">
                                <i class="bi bi-cart"></i> Sales
                            </a>
                        </li>

                        <!-- Production Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.production.*') ? 'active' : '' }}" href="{{ route('admin.upload.production.index') }}">
                                <i class="bi bi-gear"></i> Production
                            </a>
                        </li>

                        <!-- Finance Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.finance.*') ? 'active' : '' }}" href="{{ route('admin.upload.finance.index') }}">
                                <i class="bi bi-currency-dollar"></i> Finance
                            </a>
                        </li>

                        <!-- Inventory Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.inventory.*') ? 'active' : '' }}" href="{{ route('admin.upload.inventory.index') }}">
                                <i class="bi bi-box"></i> Inventory
                            </a>
                        </li>

                        <!-- HR Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.hr.*') ? 'active' : '' }}" href="{{ route('admin.upload.hr.index') }}">
                                <i class="bi bi-people"></i> HR
                            </a>
                        </li>

                        <!-- Supply Chain Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.supplychain.*') ? 'active' : '' }}" href="{{ route('admin.upload.supplychain.index') }}">
                                <i class="bi bi-truck"></i> Supply Chain
                            </a>
                        </li>

                        <!-- NPD Module -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.upload.npd.*') ? 'active' : '' }}" href="{{ route('admin.upload.npd.index') }}">
                                <i class="bi bi-lightbulb"></i> NPD
                            </a>
                        </li>

                        <div class="module-header">Management</div>

                        <!-- External API -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.external-apis.*') ? 'active' : '' }}" href="{{ route('admin.external-apis.index') }}">
                                <i class="bi bi-link-45deg"></i> External API
                            </a>
                        </li>

                        <!-- Data Viewer -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.data.*') ? 'active' : '' }}" href="{{ route('admin.data.viewer') }}">
                                <i class="bi bi-database"></i> Data Viewer
                            </a>
                        </li>

                        <!-- Demo Data Generator -->
                        <li class="nav-item">
                            <a class="nav-link {{ request()->routeIs('admin.demo-data.*') ? 'active' : '' }}" href="{{ route('admin.demo-data.index') }}">
                                <i class="bi bi-bootstrap-reboot"></i> Demo Data
                            </a>
                        </li>
                    </ul>
                </div>
            </nav>

            <!-- Main Content -->
            <main class="col-md-10 ms-sm-auto px-md-4">
                <div class="pt-3 pb-2 mb-3 border-bottom">
                    <h1 class="h2">@yield('page-title', 'Admin Panel')</h1>
                </div>

                @if(session('success'))
                    <div class="alert alert-success alert-dismissible fade show" role="alert">
                        {{ session('success') }}
                        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
                    </div>
                @endif

                @if(session('error'))
                    <div class="alert alert-danger alert-dismissible fade show" role="alert">
                        {{ session('error') }}
                        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
                    </div>
                @endif

                @yield('content')
            </main>
        </div>
    </div>

    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.0/dist/js/bootstrap.bundle.min.js"></script>
    @stack('scripts')
</body>
</html>
