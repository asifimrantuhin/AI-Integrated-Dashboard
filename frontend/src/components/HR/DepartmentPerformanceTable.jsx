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

const DepartmentPerformanceTable = ({ departments = [] }) => {
  if (!departments.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Department Performance
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Department</TableCell>
                <TableCell align="right">Employees</TableCell>
                <TableCell align="right">Present</TableCell>
                <TableCell align="right">Attendance %</TableCell>
                <TableCell align="right">Promotions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {departments.map((dept, index) => (
                <TableRow key={`${dept.department}-${index}`}>
                  <TableCell>{dept.department || 'N/A'}</TableCell>
                  <TableCell align="right">{formatNumber(dept.total_employees)}</TableCell>
                  <TableCell align="right">{formatNumber(dept.present_count)}</TableCell>
                  <TableCell align="right">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Box width="100%">
                        <LinearProgress
                          variant="determinate"
                          value={Math.min(100, Math.max(0, dept.attendance_rate || 0))}
                          color={dept.attendance_rate >= 90 ? 'success' : dept.attendance_rate >= 80 ? 'primary' : 'warning'}
                        />
                      </Box>
                      <Typography variant="caption">
                        {formatNumber(dept.attendance_rate, 1)}%
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell align="right">{formatNumber(dept.promotions)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default DepartmentPerformanceTable
