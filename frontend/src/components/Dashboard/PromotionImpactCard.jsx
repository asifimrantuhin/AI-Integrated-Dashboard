import CampaignIcon from '@mui/icons-material/Campaign'
import InsightsIcon from '@mui/icons-material/Insights'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import { Avatar, Card, CardContent, List, ListItem, ListItemAvatar, ListItemText, Stack, Typography } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const formatPercent = (value) =>
  typeof value === 'number' ? `${value.toFixed(1)}%` : value

const PromotionImpactCard = ({ promotions = [], title = 'Promotion Impact', maxItems = 5 }) => {
  if (!promotions?.length) {
    return null
  }

  const sorted = [...promotions]
    .sort((a, b) => (b.roi || 0) - (a.roi || 0))
    .slice(0, maxItems)

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          {title}
        </Typography>
        <List dense disablePadding>
          {sorted.map((promo, index) => (
            <ListItem key={`${promo.campaign_code}-${index}`} sx={{ alignItems: 'flex-start', py: 1 }}>
              <ListItemAvatar>
                <Avatar sx={{ bgcolor: promo.roi >= 120 ? 'success.light' : promo.roi >= 90 ? 'warning.light' : 'error.light', color: 'text.primary' }}>
                  {promo.roi >= 120 ? <InsightsIcon /> : <CampaignIcon />}
                </Avatar>
              </ListItemAvatar>
              <ListItemText
                primary={
                  <Stack spacing={0.25}>
                    <Typography variant="subtitle2">{promo.campaign_name || promo.campaign_code}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {promo.channel_name ? `Channel: ${promo.channel_name}` : '—'}
                    </Typography>
                  </Stack>
                }
                secondary={
                  <Stack spacing={0.5} mt={0.5}>
                    <Stack direction="row" spacing={1} alignItems="center">
                      <TrendingUpIcon fontSize="small" color="primary" />
                      <Typography variant="caption" color="text.secondary">
                        Uplift: {formatCurrencyCrore(promo.revenue_uplift || 0)} ({formatPercent(promo.uplift_percentage || 0)})
                      </Typography>
                    </Stack>
                    <Stack direction="row" spacing={2} flexWrap="wrap">
                      <Typography variant="caption" color="text.secondary">
                        Spend: {formatCurrencyCrore(promo.spend_amount || 0)}
                      </Typography>
                      <Typography variant="caption" color={promo.roi >= 100 ? 'success.main' : 'warning.main'}>
                        ROI: {formatPercent(promo.roi || 0)}
                      </Typography>
                      {promo.audience_tags?.length ? (
                        <Typography variant="caption" color="text.secondary" noWrap>
                          Audience: {promo.audience_tags.join(', ')}
                        </Typography>
                      ) : null}
                    </Stack>
                  </Stack>
                }
              />
            </ListItem>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default PromotionImpactCard
