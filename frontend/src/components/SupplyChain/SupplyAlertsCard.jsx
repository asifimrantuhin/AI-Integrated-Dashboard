import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import AssessmentIcon from '@mui/icons-material/Assessment'

const formatCurrency = (value) =>
  typeof value === 'number' ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : value

const SupplyAlertsCard = ({ alerts = [], forecast }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        Supply Chain Alerts & Forecast
      </Typography>
      <List dense>
        {alerts.length ? (
          alerts.map((alert, index) => (
            <ListItem key={`${alert}-${index}`}>
              <ListItemIcon>
                <WarningAmberIcon color="error" />
              </ListItemIcon>
              <ListItemText primary={alert} />
            </ListItem>
          ))
        ) : (
          <ListItem>
            <ListItemText primary="No critical supply chain alerts." />
          </ListItem>
        )}
      </List>
      {forecast && (
        <>
          <Divider sx={{ my: 2 }} />
          <List dense>
            <ListItem>
              <ListItemIcon>
                <AssessmentIcon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary={`Forecast Spend: ${formatCurrency(forecast.total_forecast)}`}
                secondary={
                  forecast.average_daily_forecast
                    ? `Avg Daily Forecast: ${formatCurrency(forecast.average_daily_forecast)}`
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

export default SupplyAlertsCard
