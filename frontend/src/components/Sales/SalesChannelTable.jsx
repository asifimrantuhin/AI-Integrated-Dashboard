import {
  Card,
  CardContent,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  LinearProgress,
  Box,
} from '@mui/material'

const formatNumber = (value) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: 0 }) : value

const SalesChannelTable = ({ channels = [] }) => {
  if (!channels.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Channel Performance
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Channel</TableCell>
                <TableCell align="right">Sales</TableCell>
                <TableCell align="right">Target</TableCell>
                <TableCell align="right">Achievement</TableCell>
                <TableCell align="right">Contribution</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {channels.map((channel) => (
                <TableRow key={channel.channel_id || channel.channel_name}>
                  <TableCell>{channel.channel_name}</TableCell>
                  <TableCell align="right">৳ {formatNumber(channel.billed)}</TableCell>
                  <TableCell align="right">৳ {formatNumber(channel.target)}</TableCell>
                  <TableCell align="right">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Box width="100%">
                        <LinearProgress
                          variant="determinate"
                          value={Math.min(100, Math.max(0, channel.achievement || 0))}
                          color={channel.achievement >= 100 ? 'success' : 'primary'}
                        />
                      </Box>
                      <Typography variant="caption">
                        {channel.achievement?.toFixed(1)}%
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell align="right">
                    {channel.contribution?.toFixed(1)}%
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default SalesChannelTable
