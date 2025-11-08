import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import WarningIcon from '@mui/icons-material/Warning'
import InsightsIcon from '@mui/icons-material/Insights'

const ProductionAlertsCard = ({ alerts = [], forecast }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        Operational Alerts & Forecast
      </Typography>
      <List dense>
        {alerts.length > 0 ? (
          alerts.map((alert, index) => (
            <ListItem key={`${alert}-${index}`}
              secondaryAction={<Typography variant="caption">Alert</Typography>}
            >
              <ListItemIcon>
                <WarningIcon color="error" />
              </ListItemIcon>
              <ListItemText primary={alert} />
            </ListItem>
          ))
        ) : (
          <ListItem>
            <ListItemText primary="All systems stable. No critical alerts." />
          </ListItem>
        )}
      </List>
      {forecast && (
        <>
          <Divider sx={{ my: 2 }} />
          <List dense>
            <ListItem>
              <ListItemIcon>
                <InsightsIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary={`Forecast Output: ${forecast.total_forecast?.toLocaleString(undefined, { maximumFractionDigits: 0 })}`}
                secondary={
                  forecast.average_daily_forecast
                    ? `Avg Daily Output: ${forecast.average_daily_forecast.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
                    : undefined
                }
              />
            </ListItem>
            <ListItem>
              <ListItemText
                primary={`Confidence: ${forecast.confidence_level?.toFixed(1)}%`}
                secondary={forecast.model_used ? `Model: ${forecast.model_used}` : undefined}
              />
            </ListItem>
          </List>
        </>
      )}
    </CardContent>
  </Card>
)

export default ProductionAlertsCard
