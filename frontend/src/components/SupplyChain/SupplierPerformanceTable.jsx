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
  Chip,
} from '@mui/material'

const formatScore = (value) =>
  typeof value === 'number' ? value.toFixed(1) : value

const SupplierPerformanceTable = ({ suppliers = [] }) => {
  if (!suppliers.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Supplier Performance
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Supplier</TableCell>
                <TableCell align="right">Overall</TableCell>
                <TableCell align="right">On-Time %</TableCell>
                <TableCell align="right">Quality</TableCell>
                <TableCell align="right">Cost</TableCell>
                <TableCell align="center">Rating</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {suppliers.map((supplier, index) => (
                <TableRow key={`${supplier.supplier_name}-${index}`}>
                  <TableCell>{supplier.supplier_name}</TableCell>
                  <TableCell align="right">{formatScore(supplier.overall_score)}</TableCell>
                  <TableCell align="right">{formatScore(supplier.on_time_percentage)}%</TableCell>
                  <TableCell align="right">{formatScore(supplier.quality_score)}</TableCell>
                  <TableCell align="right">{formatScore(supplier.cost_score)}</TableCell>
                  <TableCell align="center">
                    <Chip
                      label={supplier.rating || 'N/A'}
                      size="small"
                      color={supplier.rating === 'excellent' ? 'success' : supplier.rating === 'good' ? 'primary' : supplier.rating === 'average' ? 'warning' : 'default'}
                      sx={{ textTransform: 'capitalize' }}
                    />
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

export default SupplierPerformanceTable
