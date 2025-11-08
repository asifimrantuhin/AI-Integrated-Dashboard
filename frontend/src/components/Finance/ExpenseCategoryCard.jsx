import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'

const formatCurrency = (value) =>
  typeof value === 'number' ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : value

const ExpenseCategoryCard = ({ categories = [] }) => {
  if (!categories.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Top Expense Categories
        </Typography>
        <List dense>
          {categories.map((category, index) => (
            <>
              <ListItem key={`${category.category_id}-${index}`}>
                <ListItemText
                  primary={category.category_name || 'Unassigned'}
                  secondary={formatCurrency(category.actual)}
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

export default ExpenseCategoryCard
