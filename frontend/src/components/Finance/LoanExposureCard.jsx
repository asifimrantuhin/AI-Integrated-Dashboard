import { Card, CardContent, Typography, List, ListItem, ListItemText, Divider } from '@mui/material'
import { formatCurrencyCrore } from '../../utils/formatNumber'

const LoanExposureCard = ({ loans = [] }) => {
  if (!loans.length) {
    return null
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Loan Exposure by Head
        </Typography>
        <List dense>
          {loans.map((loan, index) => (
            <>
              <ListItem key={`${loan.head}-${index}`}>
                <ListItemText primary={loan.head || 'Loan Head'} secondary={formatCurrencyCrore(loan.amount)} />
              </ListItem>
              {index < loans.length - 1 && <Divider component="li" />}
            </>
          ))}
        </List>
      </CardContent>
    </Card>
  )
}

export default LoanExposureCard
