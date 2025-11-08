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
import FinanceTrendChart from '../components/Finance/FinanceTrendChart'
import DepartmentPerformanceTable from '../components/Finance/DepartmentPerformanceTable'
import ExpenseCategoryCard from '../components/Finance/ExpenseCategoryCard'
import LoanExposureCard from '../components/Finance/LoanExposureCard'
import FinanceAlertsCard from '../components/Finance/FinanceAlertsCard'
import FinanceForecastChart from '../components/AI/FinanceForecastChart'
import RecommendationList from '../components/Dashboard/RecommendationList'
import ScenarioImpactCard from '../components/Dashboard/ScenarioImpactCard'
import useDashboardData from '../hooks/useDashboardData'

const FinanceDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale } = useDashboardData('/finance/overview')

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
            Finance Command Center
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Stay ahead of profitability, liquidity, and risk exposure with AI-backed variance insights.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh finance data">
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
            <FinanceTrendChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <FinanceAlertsCard alerts={data.alerts} forecast={data.forecast} />
        </Grid>

        {data.prescriptions?.financial_actions && (
          <Grid item xs={12} md={6}>
            <RecommendationList title="Financial Recommendations" items={data.prescriptions.financial_actions} subtitleKey="rationale" />
          </Grid>
        )}

        <Grid item xs={12} md={6}>
          <DepartmentPerformanceTable departments={data.departments} />
        </Grid>
        <Grid item xs={12} md={6}>
          <ExpenseCategoryCard categories={data.categories} />
        </Grid>

        <Grid item xs={12} md={6}>
          <LoanExposureCard loans={data.loans} />
        </Grid>
        <Grid item xs={12} md={6}>
          <FinanceForecastChart
            forecastData={data.forecast}
            historicalData={data.trend}
          />
        </Grid>
        {data.scenario && (
          <Grid item xs={12}>
            <ScenarioImpactCard scenario={data.scenario} title="Finance What-if" />
          </Grid>
        )}
      </Grid>
    </Container>
  )
}

export default FinanceDashboard
