import { Card, CardContent, Typography, List, ListItem, ListItemIcon, ListItemText, Divider } from '@mui/material'
import ScienceIcon from '@mui/icons-material/Science'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const formatNumber = (value, fraction = 0) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: fraction }) : value

const WastageCostCard = ({ wastage = [] }) => {
  if (!wastage.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Wastage & Scrap Analysis
        </Typography>
        <List dense>
          {wastage.map((item, index) => (
            <>
              <ListItem key={`${item.factory}-${index}`}>
                <ListItemIcon>
                  <ScienceIcon color={item.rate > 5 ? 'error' : 'primary'} />
                </ListItemIcon>
                <ListItemText
                  primary={item.factory || 'Unknown Factory'}
                  secondary={`Wastage: ${formatNumber(item.wastage, 1)} | Cost: ${formatCurrencyCrore(item.amount)} | Rate: ${formatNumber(item.rate, 2)}%`}
                />
              </ListItem>
              {index < wastage.length - 1 && <Divider component="li" />}
            </>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default WastageCostCard
