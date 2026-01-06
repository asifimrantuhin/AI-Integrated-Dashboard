import React from 'react'
import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const ExpenseCategoryCard = ({ categories } = {}) => {
  // Ensure `categories` is always an array to avoid null/undefined access
  const safeCategories = Array.isArray(categories) ? categories : (categories ? [categories] : [])
  if (!safeCategories || !safeCategories.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Top Expense Categories
        </Typography>
        <List dense>
          {safeCategories.map((category, index) => (
            <React.Fragment key={`${category?.category_id || index}-${index}`}>
              <ListItem>
                <ListItemText
                  primary={category?.category_name || 'Unassigned'}
                  secondary={formatCurrencyCrore(category?.actual)}
                />
              </ListItem>
              {index < safeCategories.length - 1 && <Divider component="li" />}
            </React.Fragment>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default ExpenseCategoryCard
