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
} from '@mui/material'

const formatCurrency = (value) =>
  typeof value === 'number' ? `৳ ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}` : value

const DepartmentScoreTable = ({ departments = [] }) => {
  if (!departments.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Department Scorecard
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Department</TableCell>
                <TableCell align="right">Revenue</TableCell>
                <TableCell align="right">Cost</TableCell>
                <TableCell align="right">Margin</TableCell>
                <TableCell align="right">Attendance %</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {departments.map((dept, index) => (
                <TableRow key={`${dept.department}-${index}`}>
                  <TableCell>{dept.department || 'N/A'}</TableCell>
                  <TableCell align="right">{formatCurrency(dept.revenue)}</TableCell>
                  <TableCell align="right">{formatCurrency(dept.cost)}</TableCell>
                  <TableCell align="right">{formatCurrency(dept.margin)}</TableCell>
                  <TableCell align="right">{dept.attendance?.toFixed(1) ?? '0.0'}%</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default DepartmentScoreTable
