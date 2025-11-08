import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import InsightsIcon from '@mui/icons-material/Insights'

const HRAlertsCard = ({ alerts = [], forecast }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        HR Alerts & Forecast
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
            <ListItemText primary="No critical HR alerts." />
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
                primary={`Forecast Attrition: ${forecast.total_forecast?.toLocaleString(undefined, { maximumFractionDigits: 0 })}`}
                secondary={
                  forecast.average_daily_forecast
                    ? `Avg Daily: ${forecast.average_daily_forecast.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
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

export default HRAlertsCard
