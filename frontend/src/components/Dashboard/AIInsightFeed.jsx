import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText } from '@mui/material'
import LightbulbIcon from '@mui/icons-material/Lightbulb'

const AIInsightFeed = ({ insights = [] }) => {
  if (!insights.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          AI Insights
        </Typography>
        <List>
          {insights.map((insight, index) => (
            <ListItem key={`${insight}-${index}`} disablePadding>
              <ListItemIcon>
                <LightbulbIcon color="warning" />
              </ListItemIcon>
              <ListItemText primary={insight} />
            </ListItem>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default AIInsightFeed
