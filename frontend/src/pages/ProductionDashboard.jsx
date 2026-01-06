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
import ProductionTrendChart from '../components/Production/ProductionTrendChart'
import ProductionLineTable from '../components/Production/ProductionLineTable'
import WastageCostCard from '../components/Production/WastageCostCard'
import MaintenanceIssuesCard from '../components/Production/MaintenanceIssuesCard'
import ProductionAlertsCard from '../components/Production/ProductionAlertsCard'
import ProductionForecastChart from '../components/AI/ProductionForecastChart'
import useDashboardData from '../hooks/useDashboardData'
import { useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import RecommendationList from '../components/Dashboard/RecommendationList'
import InsightListCard from '../components/Dashboard/InsightListCard'

const ProductionDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale, load } = useDashboardData('/production/overview')

  const location = useLocation()

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

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
            Production Excellence
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Monitor factory throughput, efficiency, and predictive alerts to keep every line running at peak performance.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh production data">
                <Refresh />
              </IconButton>
            </span>
          </Tooltip>
        </Stack>
      </Stack>

      <Grid container spacing={3}>
        {data.kpis?.map((kpi) => (
          <Grid item xs={12} sm={6} md={4} lg={3} key={kpi.label}>
            <KPIWidget title={kpi.label} value={kpi.value} change={kpi.change} unit={kpi.unit} trend={kpi.trend} />
          </Grid>
        ))}

        {data.trend && (
          <Grid item xs={12} md={8}>
            <ProductionTrendChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <ProductionAlertsCard alerts={data.alerts} forecast={data.forecast} />
        </Grid>

        <Grid item xs={12} md={6}>
          <ProductionLineTable lines={data.lines} />
        </Grid>
        <Grid item xs={12} md={6}>
          <WastageCostCard wastage={data.wastage} />
        </Grid>

        {data.recommendations?.length ? (
          <Grid item xs={12} md={6}>
            <RecommendationList title="AI Recommended Actions" items={data.recommendations} />
          </Grid>
        ) : null}

        {data.insights?.length ? (
          <Grid item xs={12} md={6}>
            <InsightListCard title="Operations Insights" insights={data.insights} />
          </Grid>
        ) : null}

        <Grid item xs={12} md={6}>
          <MaintenanceIssuesCard maintenance={data.maintenance} />
        </Grid>
        <Grid item xs={12} md={6}>
          <ProductionForecastChart
            forecastData={data.forecast}
            historicalData={data.trend}
          />
        </Grid>
      </Grid>
    </Container>
  )
}

export default ProductionDashboard
