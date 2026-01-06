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
import { formatCurrencyCrore } from '../../utils/formatNumber'

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
                <TableCell align="right">Budget</TableCell>
                <TableCell align="right">Actual</TableCell>
                <TableCell align="right">Variance</TableCell>
                <TableCell align="right">Variance %</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {departments.map((dept, index) => (
                <TableRow key={`${dept.department_id}-${index}`}>
                  <TableCell>{dept.department_name || 'N/A'}</TableCell>
                  <TableCell align="right">{formatCurrencyCrore(dept.budget)}</TableCell>
                  <TableCell align="right">{formatCurrencyCrore(dept.actual)}</TableCell>
                  <TableCell align="right">{formatCurrencyCrore(dept.variance)}</TableCell>
                  <TableCell align="right">
                    <Box display="flex" alignItems="center" gap={1}>
                      <Box width="100%">
                        <LinearProgress
                          variant="determinate"
                          value={Math.min(100, Math.max(0, 100 - Math.abs(dept.variance_percent || 0)))}
                          color={dept.variance_percent >= 0 ? 'success' : 'error'}
                        />
                      </Box>
                      <Typography variant="caption">
                        {dept.variance_percent?.toFixed(1)}%
                      </Typography>
                    </Box>
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

export default DepartmentPerformanceTable
