import { Card, CardContent, Typography, Grid, List, ListItem, ListItemText } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const formatValue = (value) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: 0 }) : value

const ProductList = ({ title, items = [] }) => (
  <Card variant="outlined" sx={{ height: '100%' }}>
    <CardContent>
      <Typography variant="subtitle1" gutterBottom>
        {title}
      </Typography>
      <List dense>
        {items.slice(0, 7).map((item, index) => (
          <ListItem key={`${title}-${item.name}-${index}`} disablePadding>
              <ListItemText
              primary={item.name}
              secondary={`Sales: ${formatCurrencyCrore(item.value)} | Qty: ${formatValue(item.quantity)}`}
            />
          </ListItem>
        ))}
        {!items.length && (
          <Typography variant="body2" color="text.secondary">
            No data available
          </Typography>
        )}
      </List>
    </CardContent>
  </Card>
)

const ProductPerformanceCard = ({ products }) => {
  if (!products) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Product Performance
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <ProductList title="Top Products" items={products.top_products} />
          </Grid>
          <Grid item xs={12} md={4}>
            <ProductList title="Slow Movers" items={products.slow_movers} />
          </Grid>
          <Grid item xs={12} md={4}>
            <ProductList title="Best Product Groups" items={products.best_product_groups} />
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  )
}

export default ProductPerformanceCard
