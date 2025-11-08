import { ArrowDropDown, ArrowDropUp } from '@mui/icons-material'
import { Card, CardContent, Typography, Box, Chip, Stack } from '@mui/material'

const formatValue = (value, unit, prefix = '', suffix = '') => {
  if (value === null || value === undefined) return '—'
  const resolvedPrefix = prefix || (unit && !unit.endsWith('%') ? `${unit} ` : '')
  const resolvedSuffix = suffix || (unit && unit.endsWith('%') ? unit : '')
  if (typeof value === 'number') {
    return `${resolvedPrefix}${value.toLocaleString(undefined, { maximumFractionDigits: 2 })}${resolvedSuffix}`
  }
  return value
}

const KPIWidget = ({ title, value, change, prefix, suffix, unit, trend }) => {
  const formattedValue = formatValue(value, unit, prefix, suffix)
  const hasChange = typeof change === 'number' && !Number.isNaN(change)
  const hasTrend = trend && (trend.direction || trend.label)
  const trendIsPositive = trend?.direction ? trend.direction === 'up' : change >= 0

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Typography variant="subtitle2" color="text.secondary" gutterBottom>
          {title}
        </Typography>
        <Stack direction="row" spacing={1.5} alignItems="center">
          <Typography variant="h5" component="div">
            {formattedValue}
          </Typography>
          {hasChange && (
            <Chip
              label={`${change > 0 ? '+' : ''}${change.toFixed(1)}%`}
              color={change >= 0 ? 'success' : 'error'}
              size="small"
            />
          )}
        </Stack>
        {hasTrend && (
          <Stack direction="row" spacing={0.5} alignItems="center" mt={1}>
            {trend?.direction && (
              trendIsPositive ? (
                <ArrowDropUp fontSize="small" color="success" />
              ) : (
                <ArrowDropDown fontSize="small" color="error" />
              )
            )}
            <Typography variant="caption" color="text.secondary">
              {trend?.label || (trendIsPositive ? 'Improving' : 'Declining')}
            </Typography>
          </Stack>
        )}
      </CardContent>
    </Card>
  )
}

export default KPIWidget
