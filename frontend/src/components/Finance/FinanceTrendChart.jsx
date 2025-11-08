import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import { Card, CardContent, Typography } from '@mui/material'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const FinanceTrendChart = ({ trend = [] }) => {
  if (!trend.length) {
    return null
  }

  const labels = trend.map((item) => item.month)

  const data = {
    labels,
    datasets: [
      {
        label: 'Budget',
        data: trend.map((item) => item.budget ?? 0),
        borderColor: 'rgb(54, 162, 235)',
        backgroundColor: 'rgba(54, 162, 235, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Actual',
        data: trend.map((item) => item.actual ?? 0),
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        tension: 0.4,
        fill: true,
      },
    ],
  }

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Budget vs Actual Trend',
      },
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  }

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent sx={{ height: 400 }}>
        <Typography variant="h6" gutterBottom>
          Financial Trend
        </Typography>
        <Line data={data} options={options} />
      </CardContent>
    </Card>
  )
}

export default FinanceTrendChart
