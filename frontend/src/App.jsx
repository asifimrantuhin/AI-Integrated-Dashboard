import { useEffect } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux'

import { checkAuth } from './store/slices/authSlice'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import SalesDashboard from './pages/SalesDashboard'
import FinanceDashboard from './pages/FinanceDashboard'
import ProductionDashboard from './pages/ProductionDashboard'
import HRDashboard from './pages/HRDashboard'
import SupplyChainDashboard from './pages/SupplyChainDashboard'
import InventoryDashboard from './pages/InventoryDashboard'
import Layout from './components/Layout/Layout'
import ExecutiveBIDashboard from './pages/ExecutiveBIDashboard'
import IntegrationJobs from './pages/IntegrationJobs'
import Reports from './pages/Reports'

function App() {
  const dispatch = useDispatch()
  const { isAuthenticated, roles } = useSelector((state) => state.auth)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (token) {
      dispatch(checkAuth())
    }
  }, [dispatch])

  if (!isAuthenticated) {
    return (
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    )
  }

  const hasRole = (role) => roles?.map((r) => r.toLowerCase()).includes(role)

  const getDefaultRoute = () => {
    if (hasRole('executive')) return '/bi'
    if (hasRole('finance_manager')) return '/finance'
    if (hasRole('sales_manager')) return '/sales'
    if (hasRole('production_manager')) return '/production'
    if (hasRole('hr_manager')) return '/hr'
    if (hasRole('supply_chain_manager')) return '/supplychain'
    if (hasRole('inventory_manager')) return '/inventory'
    if (hasRole('analyst')) return '/dashboard'
    return '/dashboard'
  }

  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Navigate to={getDefaultRoute()} replace />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/sales" element={<SalesDashboard />} />
        <Route path="/finance" element={<FinanceDashboard />} />
        <Route path="/production" element={<ProductionDashboard />} />
        <Route path="/hr" element={<HRDashboard />} />
        <Route path="/supplychain" element={<SupplyChainDashboard />} />
        <Route path="/inventory" element={<InventoryDashboard />} />
        <Route path="/bi" element={<ExecutiveBIDashboard />} />
        <Route path="/integration" element={<IntegrationJobs />} />
        <Route path="/reports" element={<Reports />} />
        <Route path="*" element={<Navigate to={getDefaultRoute()} replace />} />
      </Routes>
    </Layout>
  )
}

export default App

