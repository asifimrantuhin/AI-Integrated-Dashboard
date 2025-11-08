import { Container } from '@mui/material'
import DashboardView from '../components/Dashboard/DashboardView'

function Dashboard() {
  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <DashboardView title="Executive Dashboard" endpoint="/dashboard/executive" />
    </Container>
  )
}

export default Dashboard
