import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'

const formatCurrency = (value) =>
  typeof value === 'number' ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : value

const PendingOrdersCard = ({ pendingOrders = [] }) => {
  if (!pendingOrders.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Pending Purchase Orders
        </Typography>
        <List dense>
          {pendingOrders.map((order, index) => (
            <>
              <ListItem key={`${order.po_number}-${index}`}>
                <ListItemText
                  primary={`${order.po_number} • Pending ${formatCurrency(order.pending_amount)}`}
                  secondary={`PO Date: ${new Date(order.po_date).toLocaleDateString()} | GRN: ${formatCurrency(order.grn_amount)} | Days Pending: ${order.days_pending}`}
                />
              </ListItem>
              {index < pendingOrders.length - 1 && <Divider component="li" />}
            </>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default PendingOrdersCard
