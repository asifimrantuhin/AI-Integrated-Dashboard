import AssessmentIcon from '@mui/icons-material/Assessment'
import LocalShippingIcon from '@mui/icons-material/LocalShipping'
import PendingActionsIcon from '@mui/icons-material/PendingActions'
import PercentIcon from '@mui/icons-material/Percent'
import { Box, Card, CardContent, Chip, Divider, Grid, LinearProgress, List, ListItem, ListItemIcon, ListItemText, Stack, Typography } from '@mui/material'

const formatCurrency = (value) =>
  typeof value === 'number'
    ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
    : value

const formatPercent = (value) =>
  typeof value === 'number' ? `${value.toFixed(1)}%` : value

const PipelineSnapshotCard = ({ pipeline, title = 'Pipeline Overview', dense = false }) => {
  if (!pipeline) {
    return null
  }

  const stages = pipeline.stages || []
  const hasData = pipeline.total_orders || pipeline.total_value || stages.length
  if (!hasData) {
    return null
  }

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={1} sx={{ mb: 2 }}>
          <Typography variant="h6">{title}</Typography>
          {pipeline.conversion_rate !== undefined && (
            <Chip
              size="small"
              color={pipeline.conversion_rate >= 70 ? 'success' : pipeline.conversion_rate >= 40 ? 'warning' : 'error'}
              label={`Conversion ${formatPercent(pipeline.conversion_rate)}`}
            />
          )}
        </Stack>

        <Grid container spacing={2} sx={{ mb: 2 }}>
          <Grid item xs={6} sm={3}>
            <Stack spacing={0.5} alignItems="flex-start">
              <Typography variant="overline" color="text.secondary">
                Total Orders
              </Typography>
              <Typography variant="subtitle1">{pipeline.total_orders?.toLocaleString() || '—'}</Typography>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack spacing={0.5} alignItems="flex-start">
              <Typography variant="overline" color="text.secondary">
                Order Book
              </Typography>
              <Stack direction="row" spacing={1} alignItems="center">
                <AssessmentIcon fontSize="small" color="primary" />
                <Typography variant="subtitle1">{formatCurrency(pipeline.total_value || 0)}</Typography>
              </Stack>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack spacing={0.5} alignItems="flex-start">
              <Typography variant="overline" color="text.secondary">
                Delivered
              </Typography>
              <Stack direction="row" spacing={1} alignItems="center">
                <LocalShippingIcon fontSize="small" color="success" />
                <Typography variant="subtitle1">{formatCurrency(pipeline.delivered_value || 0)}</Typography>
              </Stack>
            </Stack>
          </Grid>
          <Grid item xs={6} sm={3}>
            <Stack spacing={0.5} alignItems="flex-start">
              <Typography variant="overline" color="text.secondary">
                Pending
              </Typography>
              <Stack direction="row" spacing={1} alignItems="center">
                <PendingActionsIcon fontSize="small" color="warning" />
                <Typography variant="subtitle1">{formatCurrency(pipeline.pending_value || 0)}</Typography>
              </Stack>
            </Stack>
          </Grid>
        </Grid>

        <Divider sx={{ mb: 2 }} />

        <Typography variant="subtitle2" gutterBottom>
          Stage Breakdown
        </Typography>
        <List dense={dense} disablePadding>
          {stages.map((stage, index) => {
            const achievement = stage.value ? (stage.delivered_value / stage.value) * 100 : 0
            return (
              <ListItem key={`${stage.status}-${index}`} sx={{ alignItems: 'flex-start', py: dense ? 0.75 : 1 }}>
                <ListItemIcon sx={{ minWidth: 36 }}>
                  <PercentIcon fontSize="small" color={achievement >= 90 ? 'success' : achievement >= 60 ? 'warning' : 'error'} />
                </ListItemIcon>
                <ListItemText
                  primary={
                    <Stack direction="row" justifyContent="space-between" alignItems="center">
                      <Typography variant="subtitle2" textTransform="capitalize">
                        {stage.status}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        {stage.orders?.toLocaleString() || 0} orders
                      </Typography>
                    </Stack>
                  }
                  secondary={
                    <Stack spacing={0.75} mt={0.75}>
                      <Stack direction="row" spacing={2} flexWrap="wrap">
                        <Typography variant="caption" color="text.secondary">
                          Value: {formatCurrency(stage.value || 0)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Pending: {formatCurrency(stage.pending_value || 0)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Discount: {formatPercent(stage.avg_discount || 0)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Margin: {formatPercent(stage.avg_margin || 0)}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          Age: {(stage.avg_age_days || 0).toFixed(1)} days
                        </Typography>
                      </Stack>
                      <Box>
                        <LinearProgress
                          variant="determinate"
                          value={Math.max(0, Math.min(100, achievement))}
                          color={achievement >= 90 ? 'success' : achievement >= 60 ? 'warning' : 'error'}
                          sx={{ height: 6, borderRadius: 999 }}
                        />
                      </Box>
                    </Stack>
                  }
                />
              </ListItem>
            )
          })}
          {!stages.length && (
            <ListItem>
              <ListItemText primary="No pipeline stages available." primaryTypographyProps={{ variant: 'body2', color: 'text.secondary' }} />
            </ListItem>
          )}
        </List>
      </CardContent>
    </Card>
  )
}

export default PipelineSnapshotCard
