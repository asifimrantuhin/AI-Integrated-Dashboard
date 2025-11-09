import LightbulbIcon from '@mui/icons-material/Lightbulb'
import { Card, CardContent, List, ListItem, ListItemIcon, ListItemText, Typography } from '@mui/material'

const InsightListCard = ({
  title = 'Insights',
  insights = [],
  emptyMessage = 'Insights will appear once data is available.',
  iconColor = 'warning',
  dense = false,
  maxItems,
}) => {
  const items = Array.isArray(insights) ? insights.filter(Boolean) : []
  const limited = maxItems ? items.slice(0, maxItems) : items

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          {title}
        </Typography>
        <List dense={dense} disablePadding>
          {limited.length ? (
            limited.map((insight, index) => (
              <ListItem key={`${insight}-${index}`} sx={{ py: dense ? 0.5 : 1 }}>
                <ListItemIcon sx={{ minWidth: 36 }}>
                  <LightbulbIcon color={iconColor} fontSize={dense ? 'small' : 'medium'} />
                </ListItemIcon>
                <ListItemText
                  primaryTypographyProps={{ variant: dense ? 'body2' : 'body1' }}
                  primary={insight}
                />
              </ListItem>
            ))
          ) : (
            <ListItem>
              <ListItemText primary={emptyMessage} primaryTypographyProps={{ variant: 'body2', color: 'text.secondary' }} />
            </ListItem>
          )}
        </List>
      </CardContent>
    </Card>
  )
}

export default InsightListCard
