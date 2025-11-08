import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import WarningAmberIcon from '@mui/icons-material/WarningAmber'
import Inventory2Icon from '@mui/icons-material/Inventory2'

const InventoryAlertsCard = ({ alerts = [], forecast }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        Inventory Alerts & Forecast
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
            <ListItemText primary="Inventory levels within optimal range." />
          </ListItem>
        )}
      </List>
      {forecast && (
        <>
          <Divider sx={{ my: 2 }} />
          <List dense>
            <ListItem>
              <ListItemIcon>
                <Inventory2Icon color="primary" />
              </ListItemIcon>
              <ListItemText
                primary={`Forecast Inventory: ${forecast.total_forecast?.toLocaleString(undefined, { maximumFractionDigits: 0 })}`}
                secondary={
                  forecast.average_daily_forecast
                    ? `Avg Daily Forecast: ${forecast.average_daily_forecast.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
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

export default InventoryAlertsCard
