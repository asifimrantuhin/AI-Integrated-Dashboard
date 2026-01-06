import TrendingDownIcon from '@mui/icons-material/TrendingDown'
import TrendingUpIcon from '@mui/icons-material/TrendingUp'
import { Card, CardContent, LinearProgress, Stack, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Typography } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const TargetVarianceCard = ({ targets = [], title = 'Target vs Actual', maxRows = 5 }) => {
  if (!targets?.length) {
    return null
  }

  const sorted = [...targets]
    .sort((a, b) => (b.revenue_gap || 0) - (a.revenue_gap || 0))
    .slice(0, maxRows)

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          {title}
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Channel</TableCell>
                <TableCell align="right">Target</TableCell>
                <TableCell align="right">Actual</TableCell>
                <TableCell align="right">Gap</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {sorted.map((row) => {
                const achievement = Math.max(0, Math.min(100, row.achievement || 0))
                const gapPositive = (row.revenue_gap || 0) > 0
                return (
                  <TableRow key={`${row.channel_id}-${row.channel_name}`} hover>
                    <TableCell>
                      <Stack spacing={0.5}>
                        <Typography variant="subtitle2">{row.channel_name || 'Channel'}</Typography>
                        <Stack direction="row" spacing={0.5} alignItems="center">
                          {gapPositive ? (
                            <TrendingDownIcon fontSize="small" color="error" />
                          ) : (
                            <TrendingUpIcon fontSize="small" color="success" />
                          )}
                          <Typography variant="caption" color="text.secondary">
                            {row.owner ? `Owner: ${row.owner}` : `Achievement: ${achievement.toFixed(1)}%`}
                          </Typography>
                        </Stack>
                        <LinearProgress
                          variant="determinate"
                          value={achievement}
                          color={achievement >= 100 ? 'success' : achievement >= 80 ? 'primary' : achievement >= 60 ? 'warning' : 'error'}
                          sx={{ height: 6, borderRadius: 999, maxWidth: 160 }}
                        />
                      </Stack>
                    </TableCell>
                    <TableCell align="right">{formatCurrencyCrore(row.revenue_target || 0)}</TableCell>
                    <TableCell align="right">{formatCurrencyCrore(row.actual_revenue || 0)}</TableCell>
                    <TableCell align="right" sx={{ color: gapPositive ? 'error.main' : 'success.main', fontWeight: 600 }}>
                      {formatCurrencyCrore(row.revenue_gap || 0)}
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default TargetVarianceCard
