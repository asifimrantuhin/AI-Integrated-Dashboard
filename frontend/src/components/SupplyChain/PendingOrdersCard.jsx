import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

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
                  primary={`${order.po_number} • Pending ${formatCurrencyCrore(order.pending_amount)}`}
                  secondary={`PO Date: ${new Date(order.po_date).toLocaleDateString()} | GRN: ${formatCurrencyCrore(order.grn_amount)} | Days Pending: ${order.days_pending}`}
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
