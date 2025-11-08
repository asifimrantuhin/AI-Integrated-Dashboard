import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText } from '@mui/material'
import WarningAmberIcon from '@mui/icons-material/WarningAmber'

const BIAlertsCard = ({ alerts = [] }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        Executive Alerts
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
            <ListItemText primary="No critical alerts" />
          </ListItem>
        )}
      </List>
    </CardContent>
  </Card>
)

export default BIAlertsCard
