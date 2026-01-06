import { Fragment, isValidElement, useMemo, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import {
  Alert,
  Box,
  Button,
  Chip,
  Grid,
  IconButton,
  Skeleton,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material'
import { Refresh, WarningAmber } from '@mui/icons-material'
import KPIWidget from './KPIWidget'
import AIInsightFeed from './AIInsightFeed'
import useDashboardData from '../../hooks/useDashboardData'

const SkeletonGrid = () => (
  <Grid container spacing={3} sx={{ mt: 1 }}>
    {Array.from({ length: 6 }).map((_, index) => (
      <Grid item xs={12} sm={6} md={4} lg={3} key={index}>
        <Skeleton variant="rounded" height={120} animation="wave" />
      </Grid>
    ))}
    <Grid item xs={12} md={6} lg={4}>
      <Skeleton variant="rounded" height={280} animation="wave" />
    </Grid>
  </Grid>
)

const DashboardView = ({ title, endpoint, description }) => {
  const { data, loading, error, stale, refreshing, refresh, load } = useDashboardData(endpoint)

  const location = useLocation()

  // Trigger a non-forced load when the route changes. `load()` will request
  // fresh data unless a fetch is already in-flight; the manual `refresh()`
  // button still forces an immediate update when the user requests it.
  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  const lastUpdatedLabel = useMemo(() => {
    if (!data?.last_updated) return null
    const date = new Date(data.last_updated)
    if (Number.isNaN(date.getTime())) return null
    return date.toLocaleString()
  }, [data?.last_updated])

  const chartElements = useMemo(() => {
    if (!data?.charts) return []
    if (Array.isArray(data.charts)) {
      return data.charts
        .map((item, index) => {
          const element = isValidElement(item) ? item : item?.component
          if (!isValidElement(element)) return null
          return {
            key: item?.key || `chart-${index}`,
            element,
          }
        })
        .filter(Boolean)
    }
    return Object.entries(data.charts)
      .map(([key, value]) => {
        const element = isValidElement(value) ? value : value?.component
        if (!isValidElement(element)) return null
        return {
          key,
          element,
        }
      })
      .filter(Boolean)
  }, [data?.charts])

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" gutterBottom>
          {title}
        </Typography>
        <Skeleton variant="text" width={260} />
        <SkeletonGrid />
      </Box>
    )
  }

  if (error) {
    return (
      <Stack spacing={2}>
        <Alert severity="error" action={<Button color="inherit" onClick={refresh}>Retry</Button>}>
          {error}
        </Alert>
      </Stack>
    )
  }

  return (
    <Stack spacing={3}>
      <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ xs: 'flex-start', sm: 'center' }} justifyContent="space-between">
        <Box>
          <Typography variant="h4" gutterBottom>
            {title}
          </Typography>
          {description && (
            <Typography variant="body2" color="text.secondary">
              {description}
            </Typography>
          )}
          {lastUpdatedLabel && (
            <Stack direction="row" spacing={1} alignItems="center" mt={0.5}>
              <Typography variant="caption" color="text.secondary">
                Last updated: {lastUpdatedLabel}
              </Typography>
              {stale && (
                <Tooltip title="Information may be outdated. Refresh for the latest metrics.">
                  <WarningAmber fontSize="inherit" color="warning" />
                </Tooltip>
              )}
            </Stack>
          )}
        </Box>
        <Stack direction="row" spacing={1} alignItems="center">
          {stale && <Chip color="warning" size="small" label="Stale" />}
          {refreshing && <Chip color="info" size="small" label="Refreshing" />}
          <Tooltip title="Refresh data">
            <span>
              <IconButton color="primary" onClick={refresh} disabled={refreshing} aria-label={`Refresh ${title} data`}>
                <Refresh />
              </IconButton>
            </span>
          </Tooltip>
        </Stack>
      </Stack>

      {data?.kpis?.length ? (
        <Grid container spacing={3}>
          {data.kpis.map((kpi) => (
            <Grid item xs={12} sm={6} md={4} lg={3} key={kpi.label}>
              <KPIWidget title={kpi.label} value={kpi.value} change={kpi.change} unit={kpi.unit} trend={kpi.trend} />
            </Grid>
          ))}
        </Grid>
      ) : (
        <SkeletonGrid />
      )}

      {(data?.ai_insights?.length || chartElements.length) ? (
        <Grid container spacing={3}>
          {data?.ai_insights?.length ? (
            <Grid item xs={12} md={6} lg={4}>
              <AIInsightFeed insights={data.ai_insights} />
            </Grid>
          ) : null}
          {chartElements.map((item) => (
            <Grid item xs={12} md={6} lg={4} key={item.key}>
              <Fragment>{item.element}</Fragment>
            </Grid>
          ))}
        </Grid>
      ) : null}
    </Stack>
  )
}

export default DashboardView
