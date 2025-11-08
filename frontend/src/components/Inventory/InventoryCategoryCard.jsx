import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'

const formatCurrency = (value) =>
  typeof value === 'number' ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : value

const InventoryCategoryCard = ({ categories = [] }) => {
  if (!categories.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Inventory by GL Account
        </Typography>
        <List dense>
          {categories.map((category, index) => (
            <>
              <ListItem key={`${category.gl_id}-${index}`}>
                <ListItemText
                  primary={`${category.gl_code || 'GL'} - ${category.gl_name || 'Account'}`}
                  secondary={formatCurrency(category.amount)}
                />
              </ListItem>
              {index < categories.length - 1 && <Divider component="li" />}
            </>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default InventoryCategoryCard
