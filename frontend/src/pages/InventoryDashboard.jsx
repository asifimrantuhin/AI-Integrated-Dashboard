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
import InventoryValueChart from '../components/Inventory/InventoryValueChart'
import InventoryCategoryCard from '../components/Inventory/InventoryCategoryCard'
import InventoryCompanyTable from '../components/Inventory/InventoryCompanyTable'
import InventoryTurnoverCard from '../components/Inventory/InventoryTurnoverCard'
import InventoryAlertsCard from '../components/Inventory/InventoryAlertsCard'
import InventoryForecastChart from '../components/AI/InventoryForecastChart'
import RecommendationList from '../components/Dashboard/RecommendationList'
import ScenarioImpactCard from '../components/Dashboard/ScenarioImpactCard'
import useDashboardData from '../hooks/useDashboardData'
import InsightListCard from '../components/Dashboard/InsightListCard'

const InventoryDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale } = useDashboardData('/inventory/overview')

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
            Inventory Control Tower
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Balance stock, cash, and service levels with predictive coverage insights across every warehouse.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh inventory data">
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
            <InventoryValueChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <InventoryAlertsCard alerts={data.alerts} forecast={data.forecast} />
        </Grid>

        {data.prescriptions?.inventory_actions && (
          <Grid item xs={12} md={6}>
            <RecommendationList title="Inventory Actions" items={data.prescriptions.inventory_actions} subtitleKey="action" />
          </Grid>
        )}
        {data.slow_movers?.length ? (
          <Grid item xs={12} md={6}>
            <RecommendationList title="Slow Movers" items={data.slow_movers} subtitleKey="suggestion" />
          </Grid>
        ) : null}

        {(data.executive_summary?.length || data.insights?.length) ? (
          <Grid item xs={12} md={6}>
            <InsightListCard
              title="Inventory Insights"
              insights={[...(data.executive_summary || []), ...(data.insights || [])]}
              iconColor="warning"
            />
          </Grid>
        ) : null}

        <Grid item xs={12} md={6}>
          <InventoryCategoryCard categories={data.categories} />
        </Grid>
        <Grid item xs={12} md={6}>
          <InventoryCompanyTable companies={data.companies} />
        </Grid>

        <Grid item xs={12} md={6}>
          <InventoryTurnoverCard turnover={data.turnover} />
        </Grid>
        <Grid item xs={12} md={6}>
          <InventoryForecastChart
            forecastData={data.forecast}
            historicalData={data.trend}
          />
        </Grid>
        {data.scenario && (
          <Grid item xs={12}>
            <ScenarioImpactCard scenario={data.scenario} title="Inventory What-if" />
          </Grid>
        )}

      </Grid>
    </Container>
  )
}

export default InventoryDashboard
