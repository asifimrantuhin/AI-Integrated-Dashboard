import { Card, CardContent, Typography, Grid, Chip } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const formatCurrencyIfNeeded = (label, value) => {
  if (typeof value !== 'number') return value
  if (label.toLowerCase().includes('attrition')) {
    return value.toFixed(1)
  }
  return formatCurrencyCrore(value)
}

const BIForecastGrid = ({ forecasts = [] }) => {
  if (!forecasts.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          AI Forecast Outlook
        </Typography>
        <Grid container spacing={2}>
          {forecasts.map((forecast, index) => (
            <Grid item xs={12} sm={6} md={4} key={`${forecast.type}-${index}`}>
              <Card variant="outlined">
                <CardContent>
                  <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                    {forecast.label}
                  </Typography>
                  <Typography variant="h5">
                    {formatCurrencyIfNeeded(forecast.label, forecast.value)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Model: {forecast.model_used || 'N/A'}
                  </Typography>
                  <Chip
                    label={`Confidence ${forecast.confidence?.toFixed(1) ?? 0}%`}
                    size="small"
                    color={forecast.confidence >= 80 ? 'success' : forecast.confidence >= 60 ? 'primary' : 'warning'}
                  />
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </CardContent>
    </Card>
  )
}

export default BIForecastGrid
