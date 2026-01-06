import React from 'react'
import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const InventoryCategoryCard = ({ categories }) => {
  const safeCategories = categories || []
  if (!safeCategories.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Inventory by GL Account
        </Typography>
        <List dense>
          {safeCategories.map((category, index) => (
            <React.Fragment key={`${category?.gl_id || index}-${index}`}>
              <ListItem>
                <ListItemText
                  primary={`${category?.gl_code || 'GL'} - ${category?.gl_name || 'Account'}`}
                  secondary={formatCurrencyCrore(category?.amount)}
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

export default InventoryCategoryCard
