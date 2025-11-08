import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import AssessmentIcon from '@mui/icons-material/Assessment'

const SalesAlertsCard = ({ alerts = [], forecast }) => (
  <Card sx={{ height: '100%' }}>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        Alerts & Forecast Summary
      </Typography>
      <List dense>
        {alerts.length > 0 ? (
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
            <ListItemText primary="All clear. No critical alerts." />
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
                primary={`Forecast (Total): ৳ ${forecast.total_forecast?.toLocaleString(undefined, { maximumFractionDigits: 0 })}`}
                secondary={
                  forecast.average_daily_forecast
                    ? `Avg Daily Forecast: ৳ ${forecast.average_daily_forecast.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
                    : undefined
                }
              />
            </ListItem>
            <ListItem>
              <ListItemText
                primary={`Confidence Level: ${forecast.confidence_level?.toFixed(1)}%`}
                secondary={forecast.model_used ? `Model: ${forecast.model_used}` : undefined}
              />
            </ListItem>
          </List>
        </>
      )}
    </CardContent>
  </Card>
)

export default SalesAlertsCard
