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

const SalesPerformanceChart = ({ trend = [] }) => {
  if (!trend.length) {
    return null
  }

  const labels = trend.map((item) => item.month)

  const data = {
    labels,
    datasets: [
      {
        label: 'Sales',
        data: trend.map((item) => item.billed || 0),
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Target',
        data: trend.map((item) => item.target || 0),
        borderColor: 'rgb(255, 159, 64)',
        backgroundColor: 'rgba(255, 159, 64, 0.2)',
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
        text: 'Sales vs Target Trend',
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
          Sales Trend
        </Typography>
        <Line data={data} options={options} />
      </CardContent>
    </Card>
  )
}

export default SalesPerformanceChart
