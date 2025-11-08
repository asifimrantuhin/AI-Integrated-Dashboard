import { Card, CardContent, CardHeader, Grid, Typography } from '@mui/material'

const ScenarioImpactCard = ({ title = 'Scenario Impact', scenario }) => {
  if (!scenario) {
    return null
  }

  const metrics = [
    { label: 'Projected Sales', value: scenario.projected_sales },
    { label: 'Projected Margin', value: scenario.projected_margin },
    { label: 'Incremental Profit', value: scenario.incremental_profit },
  ]

  return (
    <Card>
      <CardHeader title={title} subheader={scenario.narrative} />
      <CardContent>
        <Grid container spacing={2}>
          {metrics.map((metric) => (
            <Grid item xs={12} md={4} key={metric.label}>
              <Typography variant="subtitle2" color="text.secondary">
                {metric.label}
              </Typography>
              <Typography variant="h6">
                {metric.value !== undefined ? metric.value.toLocaleString(undefined, { maximumFractionDigits: 2 }) : '—'}
              </Typography>
            </Grid>
          ))}
        </Grid>
      </CardContent>
    </Card>
  )
}

export default ScenarioImpactCard
