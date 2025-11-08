import { Card, CardContent, CardHeader, List, ListItem, ListItemAvatar, Avatar, ListItemText, Typography } from '@mui/material'
import { Warning, Info } from '@mui/icons-material'

const severityPalette = {
  high: 'error.main',
  medium: 'warning.main',
  low: 'success.main',
}

const AnomalyAlertList = ({ title = 'Anomaly Alerts', anomalies }) => {
  const items = anomalies?.anomalies || []

  if (!items.length) {
    return null
  }

  return (
    <Card>
      <CardHeader title={title} subheader={`${items.length} anomalies detected`} />
      <CardContent>
        <List dense disablePadding>
          {items.map((anomaly, index) => {
            const severity = anomaly.severity || 'medium'
            return (
              <ListItem key={index} alignItems="flex-start" sx={{ mb: 1 }}>
                <ListItemAvatar>
                  <Avatar sx={{ bgcolor: severityPalette[severity] || 'warning.main' }}>
                    {severity === 'high' ? <Warning fontSize="small" /> : <Info fontSize="small" />}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={
                    <Typography variant="subtitle2" color="text.primary">
                      {`Value ${anomaly.value?.toLocaleString?.() ?? anomaly.value} (${severity})`}
                    </Typography>
                  }
                  secondary={
                    <Typography variant="caption" color="text.secondary">
                      {anomaly.recommended_action}
                    </Typography>
                  }
                />
              </ListItem>
            )
          })}
        </List>
      </CardContent>
    </Card>
  )
}

export default AnomalyAlertList
