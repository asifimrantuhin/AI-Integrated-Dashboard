import { Card, CardContent, Typography, Grid } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const formatNumber = (value, fraction = 1) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: fraction }) : value

const InventoryTurnoverCard = ({ turnover } = {}) => {
  if (!turnover) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Turnover & Efficiency
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6}>
            <Typography variant="body2" color="text.secondary">
              Average Inventory
            </Typography>
              <Typography variant="h6">
                {formatCurrencyCrore(turnover.average_inventory)}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="body2" color="text.secondary">
              COGS
            </Typography>
              <Typography variant="h6">
                {formatCurrencyCrore(turnover.cogs)}
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="body2" color="text.secondary">
              Turnover Days
            </Typography>
            <Typography variant="h6">
              {formatNumber(turnover.turnover_days)} days
            </Typography>
          </Grid>
          <Grid item xs={12} sm={6}>
            <Typography variant="body2" color="text.secondary">
              GMROI
            </Typography>
            <Typography variant="h6">
              {formatNumber(turnover.gmroi)}%
            </Typography>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default InventoryTurnoverCard
