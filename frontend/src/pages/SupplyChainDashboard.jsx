import {
  Alert,
  Box,
  Chip,
  Container,
  Grid,
  IconButton,
  Skeleton,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material'
import { Refresh, WarningAmber } from '@mui/icons-material'
import KPIWidget from '../components/Dashboard/KPIWidget'
import SupplyTrendChart from '../components/SupplyChain/SupplyTrendChart'
import SupplierPerformanceTable from '../components/SupplyChain/SupplierPerformanceTable'
import PendingOrdersCard from '../components/SupplyChain/PendingOrdersCard'
import SupplyAlertsCard from '../components/SupplyChain/SupplyAlertsCard'
import SupplyForecastChart from '../components/AI/SupplyForecastChart'
import useDashboardData from '../hooks/useDashboardData'

const SupplyChainDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale } = useDashboardData('/supplychain/overview')

  if (loading) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Skeleton variant="text" height={48} width={320} sx={{ mb: 2 }} />
        <Grid container spacing={3}>
          {Array.from({ length: 6 }).map((_, index) => (
            <Grid item xs={12} sm={6} md={4} key={index}>
              <Skeleton variant="rounded" height={120} animation="wave" />
            </Grid>
          ))}
          <Grid item xs={12} md={8}>
            <Skeleton variant="rounded" height={320} animation="wave" />
          </Grid>
          <Grid item xs={12} md={4}>
            <Skeleton variant="rounded" height={320} animation="wave" />
          </Grid>
        </Grid>
      </Container>
    )
  }

  if (error) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Alert severity="error" action={<Tooltip title="Retry"><span><IconButton color="inherit" onClick={refresh}><Refresh /></IconButton></span></Tooltip>}>
          {error}
        </Alert>
      </Container>
    )
  }

  if (!data) {
    return null
  }

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Stack direction={{ xs: 'column', md: 'row' }} spacing={2} justifyContent="space-between" alignItems={{ xs: 'flex-start', md: 'center' }} sx={{ mb: 3 }}>
        <Box>
          <Typography variant="h4" gutterBottom>
            Supply Chain Nerve Center
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Anticipate supplier risk, delivery bottlenecks, and working capital swings in real time.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh supply chain data">
                <Refresh />
              </IconButton>
            </span>
          </Tooltip>
        </Stack>
      </Stack>

      <Grid container spacing={3}>
        {data.kpis?.map((kpi) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={kpi.label}>
            <KPIWidget
              title={kpi.label}
              value={kpi.value}
              change={kpi.change}
              unit={kpi.unit}
              trend={kpi.trend}
            />
          </Grid>
        ))}

        {data.trend && (
          <Grid item xs={12} md={8}>
            <SupplyTrendChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <SupplyAlertsCard alerts={data.alerts} forecast={data.forecast} />
        </Grid>

        <Grid item xs={12} md={6}>
          <SupplierPerformanceTable suppliers={data.suppliers} />
        </Grid>
        <Grid item xs={12} md={6}>
          <PendingOrdersCard pendingOrders={data.pending_orders} />
        </Grid>

        <Grid item xs={12}>
          <SupplyForecastChart
            forecastData={data.forecast}
            historicalData={data.trend}
          />
        </Grid>
      </Grid>
    </Container>
  )
}

export default SupplyChainDashboard
