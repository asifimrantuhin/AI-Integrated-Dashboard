import { Card, CardContent, CardHeader, List, ListItem, ListItemText, Chip, Stack, Typography } from '@mui/material'

const formatValue = (value) => {
  if (value === null || value === undefined) return '—'
  if (typeof value === 'number') {
    return value.toLocaleString(undefined, { maximumFractionDigits: 1 })
  }
  return value
}

const RecommendationList = ({ title, items = [], subtitleKey = 'rationale' }) => {
  return (
    <Card sx={{ height: '100%' }}>
      <CardHeader title={title} subheader={`${items.length} recommendations`} />
      <CardContent>
        {items.length === 0 ? (
          <Typography variant="body2" color="text.secondary">
            No recommendations available.
          </Typography>
        ) : (
          <List disablePadding dense>
            {items.map((item, index) => (
              <ListItem key={index} alignItems="flex-start" sx={{ mb: 1, borderRadius: 1, backgroundColor: 'background.default' }}>
                <ListItemText
                  primary={
                    <Stack direction="row" spacing={1} alignItems="center">
                      <Typography variant="subtitle2">
                        {item.entity_name || item.entity_type || `Item ${index + 1}`}
                      </Typography>
                      {item.risk_level && (
                        <Chip size="small" label={item.risk_level} color={item.risk_level === 'high' ? 'error' : item.risk_level === 'low' ? 'success' : 'warning'} />
                      )}
                    </Stack>
                  }
                  secondary={
                    <>
                      <Typography variant="body2" color="text.primary">
                        {item.recommended_action || item.action}
                      </Typography>
                      {item[subtitleKey] && (
                        <Typography variant="caption" display="block" color="text.secondary">
                          {item[subtitleKey]}
                        </Typography>
                      )}
                      <Stack direction="row" spacing={2} mt={0.5}>
                        {item.growth_rate !== undefined && (
                          <Typography variant="caption" color="text.secondary">
                            Growth: {formatValue(item.growth_rate)}%
                          </Typography>
                        )}
                        {item.lift_vs_baseline !== undefined && (
                          <Typography variant="caption" color="text.secondary">
                            Lift vs baseline: {formatValue(item.lift_vs_baseline)}%
                          </Typography>
                        )}
                        {item.recommended_order_qty !== undefined && (
                          <Typography variant="caption" color="text.secondary">
                            Order qty: {formatValue(item.recommended_order_qty)}
                          </Typography>
                        )}
                        {item.average_variance !== undefined && (
                          <Typography variant="caption" color="text.secondary">
                            Variance: {formatValue(item.average_variance)}
                          </Typography>
                        )}
                      </Stack>
                    </>
                  }
                />
              </ListItem>
            ))}
          </List>
        )}
      </CardContent>
    </Card>
  )
}

export default RecommendationList
