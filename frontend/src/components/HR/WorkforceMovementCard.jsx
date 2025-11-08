import { Card, CardContent, Typography, Grid } from '@mui/material'

const formatNumber = (value, fraction = 0) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: fraction }) : value

const WorkforceMovementCard = ({ movements = [] }) => {
  if (!movements.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Workforce Movements
        </Typography>
        <Grid container spacing={2}>
          {movements.map((movement) => (
            <Grid item xs={12} sm={6} md={4} key={movement.type}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ textTransform: 'capitalize' }}>
                {movement.type}
              </Typography>
              <Typography variant="h5">
                {formatNumber(movement.total_count)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Hours/Amount: {formatNumber(movement.total_amount, 1)}
              </Typography>
            </Grid>
          ))}
        </Grid>
      </CardContent>
    </Card>
  )
}

export default WorkforceMovementCard
