import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText } from '@mui/material'
import BuildCircleIcon from '@mui/icons-material/BuildCircle'
import AccessTimeIcon from '@mui/icons-material/AccessTime'

const formatNumber = (value, fraction = 0) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: fraction }) : value

const MaintenanceIssuesCard = ({ maintenance = [] }) => {
  if (!maintenance.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Maintenance Hotspots
        </Typography>
        <List dense>
          {maintenance.map((item, index) => (
            <ListItem key={`${item.machine_code}-${index}`}>
              <ListItemIcon>
                <BuildCircleIcon color={item.downtime_minutes > 60 ? 'error' : 'primary'} />
              </ListItemIcon>
              <ListItemText
                primary={`${item.machine_code} - ${item.machine_name || 'Machine'}`}
                secondary={`Downtime: ${formatNumber(item.downtime_minutes)} mins | Events: ${item.events} | Cost: ৳ ${formatNumber(item.cost, 0)} | Last: ${item.last_date ? new Date(item.last_date).toLocaleDateString() : 'N/A'}`}
              />
              <AccessTimeIcon fontSize="small" color="action" sx={{ ml: 1 }} />
            </ListItem>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default MaintenanceIssuesCard
