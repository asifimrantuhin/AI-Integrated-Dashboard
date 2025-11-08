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
import SalesPerformanceChart from '../components/Sales/SalesPerformanceChart'
import SalesChannelTable from '../components/Sales/SalesChannelTable'
import ProductPerformanceCard from '../components/Sales/ProductPerformanceCard'
import SalesAlertsCard from '../components/Sales/SalesAlertsCard'
import SalesForecastChart from '../components/AI/SalesForecastChart'
import useDashboardData from '../hooks/useDashboardData'
import RecommendationList from '../components/Dashboard/RecommendationList'
import ScenarioImpactCard from '../components/Dashboard/ScenarioImpactCard'
import AnomalyAlertList from '../components/Dashboard/AnomalyAlertList'

const SalesDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale } = useDashboardData('/sales/overview')

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
            Sales Intelligence
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Track revenue health, channel execution, and AI-driven demand signals across the network.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && (
            <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />
          )}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh sales data">
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
            <SalesPerformanceChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <SalesAlertsCard alerts={data.alerts} forecast={data.forecast} />
        </Grid>

        {data.predictions?.length ? (
          <Grid item xs={12} md={6}>
            <RecommendationList title="Channel Recommendations" items={data.recommendations || []} />
          </Grid>
        ) : null}
        {data.anomalies?.anomalies?.length ? (
          <Grid item xs={12} md={6}>
            <AnomalyAlertList anomalies={data.anomalies} />
          </Grid>
        ) : null}

        <Grid item xs={12} md={6}>
          <SalesChannelTable channels={data.channels} />
        </Grid>
        <Grid item xs={12} md={6}>
          <ProductPerformanceCard products={data.products} />
        </Grid>

        <Grid item xs={12}>
          <SalesForecastChart
            forecastData={data.forecast}
            historicalData={data.trend}
          />
        </Grid>
        {data.scenario && (
          <Grid item xs={12}>
            <ScenarioImpactCard scenario={data.scenario} title="What-if Impact" />
          </Grid>
        )}
      </Grid>
    </Container>
  )
}

export default SalesDashboard
