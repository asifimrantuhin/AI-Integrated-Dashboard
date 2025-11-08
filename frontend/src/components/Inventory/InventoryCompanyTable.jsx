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

const InventoryCompanyTable = ({ companies = [] }) => {
  if (!companies.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Inventory by Company
        </Typography>
        <TableContainer>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>Company</TableCell>
                <TableCell align="right">Inventory Value</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {companies.map((company, index) => (
                <TableRow key={`${company.company_id}-${index}`}>
                  <TableCell>{company.company_name || `Company #${company.company_id}`}</TableCell>
                  <TableCell align="right">{formatCurrency(company.amount)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </CardContent>
    </Card>
  )
}

export default InventoryCompanyTable
