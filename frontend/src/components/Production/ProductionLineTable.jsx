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

const formatNumber = (value, fraction = 0) =>
  typeof value === 'number' ? value.toLocaleString(undefined, { maximumFractionDigits: fraction }) : value

const ProductionLineTable = ({ lines = [] }) => {
  if (!lines.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Line Performance
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Factory</TableCell>
                <TableCell>Line</TableCell>
                <TableCell align="right">Actual Output</TableCell>
                <TableCell align="right">Planned Output</TableCell>
                <TableCell align="right">Efficiency</TableCell>
                <TableCell align="right">OEE</TableCell>
                <TableCell align="right">Downtime (mins)</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {lines.map((line, index) => (
                <TableRow key={`${line.factory_id}-${line.production_line_id}-${index}`}>
                  <TableCell>{line.factory_id ? `Factory #${line.factory_id}` : 'N/A'}</TableCell>
                  <TableCell>{line.production_line_id ? `Line #${line.production_line_id}` : 'N/A'}</TableCell>
                  <TableCell align="right">{formatNumber(line.actual_output)}</TableCell>
                  <TableCell align="right">{formatNumber(line.planned_output)}</TableCell>
                  <TableCell align="right">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Box width="100%">
                        <LinearProgress
                          variant="determinate"
                          value={Math.min(100, Math.max(0, line.efficiency_percentage || 0))}
                          color={line.efficiency_percentage >= 90 ? 'success' : 'primary'}
                        />
                      </Box>
                      <Typography variant="caption">
                        {formatNumber(line.efficiency_percentage, 1)}%
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell align="right">{formatNumber(line.oee, 1)}%</TableCell>
                  <TableCell align="right">{formatNumber(line.downtime_minutes)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default ProductionLineTable
