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
import { formatCurrencyCrore } from '../../utils/formatNumber'

const InventoryCompanyTable = ({ companies } = {}) => {
  // Make destructuring safe if `props` is null and ensure companies is an array
  const safeCompanies = Array.isArray(companies) ? companies : (companies ? [companies] : [])
  if (!safeCompanies || !safeCompanies.length) {
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
              {safeCompanies.map((company, index) => (
                <TableRow key={`${company.company_id}-${index}`}>
                  <TableCell>{company.company_name || `Company #${company.company_id}`}</TableCell>
                  <TableCell align="right">{formatCurrencyCrore(company.amount)}</TableCell>
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
