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
import DepartmentScoreTable from '../components/ExecutiveBI/DepartmentScoreTable'
import BIForecastGrid from '../components/ExecutiveBI/BIForecastGrid'
import BIAlertsCard from '../components/ExecutiveBI/BIAlertsCard'
import BIAssistant from '../components/ExecutiveBI/BIAssistant'
import SupplyTrendChart from '../components/SupplyChain/SupplyTrendChart'
import useDashboardData from '../hooks/useDashboardData'
import PipelineSnapshotCard from '../components/Dashboard/PipelineSnapshotCard'
import TargetVarianceCard from '../components/Dashboard/TargetVarianceCard'
import PromotionImpactCard from '../components/Dashboard/PromotionImpactCard'
import InsightListCard from '../components/Dashboard/InsightListCard'

const ExecutiveBIDashboard = () => {
  const { data, loading, error, refresh, refreshing, stale } = useDashboardData('/bi/overview')

  if (loading) {
    return (
      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Skeleton variant="text" height={48} width={360} sx={{ mb: 2 }} />
        <Grid container spacing={3}>
          {Array.from({ length: 6 }).map((_, index) => (
            <Grid item xs={12} sm={6} md={3} key={index}>
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
            Executive Intelligence
          </Typography>
          <Typography variant="body2" color="text.secondary">
            A consolidated view of revenue, profitability, and operational health powered by AI cognition.
          </Typography>
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip size="small" color="warning" label="Stale" icon={<WarningAmber />} />}
          {refreshing && <Chip size="small" color="info" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label="Refresh BI data">
                <Refresh />
              </IconButton>
            </span>
          </Tooltip>
        </Stack>
      </Stack>

      <Grid container spacing={3}>
        {data.kpis?.map((kpi) => (
          <Grid item xs={12} sm={6} md={3} key={kpi.label}>
            <KPIWidget title={kpi.label} value={kpi.value} change={kpi.change} unit={kpi.unit} trend={kpi.trend} />
          </Grid>
        ))}

        {data.trend && (
          <Grid item xs={12} md={8}>
            <SupplyTrendChart trend={data.trend} />
          </Grid>
        )}
        <Grid item xs={12} md={4}>
          <Stack spacing={2}>
            <BIAlertsCard alerts={data.alerts} />
            <InsightListCard
              title="Executive Summary"
              insights={data.executive_summary}
              emptyMessage="Summary arrives once AI reviews new data."
              dense
            />
            <InsightListCard title="AI Insights" insights={data.ai_insights} iconColor="warning" dense />
          </Stack>
        </Grid>

        {data.pipeline ? (
          <Grid item xs={12} lg={4}>
            <PipelineSnapshotCard pipeline={data.pipeline} title="Sales Pipeline" dense />
          </Grid>
        ) : null}

        {data.sales_targets?.length ? (
          <Grid item xs={12} lg={4}>
            <TargetVarianceCard targets={data.sales_targets} title="Channel Target Gaps" maxRows={4} />
          </Grid>
        ) : null}

        {data.promotions?.length ? (
          <Grid item xs={12} lg={4}>
            <PromotionImpactCard promotions={data.promotions} title="Top Promotions" maxItems={4} />
          </Grid>
        ) : null}

        <Grid item xs={12} md={6}>
          <DepartmentScoreTable departments={data.departments} />
        </Grid>
        <Grid item xs={12} md={6}>
          <BIForecastGrid forecasts={data.forecasts} />
        </Grid>

        <Grid item xs={12}>
          <BIAssistant />
        </Grid>
      </Grid>
    </Container>
  )
}

export default ExecutiveBIDashboard
