import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText } from '@mui/material'
import LightbulbIcon from '@mui/icons-material/Lightbulb'

const BIInsightsPanel = ({ insights = [] }) => (
  <Card>
    <CardContent>
      <Typography variant="h6" gutterBottom>
        AI Insights
      </Typography>
      <List dense>
        {insights.length ? (
          insights.map((insight, index) => (
            <ListItem key={`${insight}-${index}`}>
              <ListItemIcon>
                <LightbulbIcon color="warning" />
              </ListItemIcon>
              <ListItemText primary={insight} />
            </ListItem>
          ))
        ) : (
          <ListItem>
            <ListItemText primary="Insights will appear as data syncs in." />
          </ListItem>
        )}
      </List>
    </CardContent>
  </Card>
)

export default BIInsightsPanel
