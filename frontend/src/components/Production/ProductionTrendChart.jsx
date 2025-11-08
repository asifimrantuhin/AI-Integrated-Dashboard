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

const getNumber = (value, fallback) => {
  if (typeof value === 'number') return value
  if (typeof fallback === 'number') return fallback
  return 0
}

const ProductionTrendChart = ({ trend = [] }) => {
  if (!trend.length) {
    return null
  }

  const labels = trend.map((item) => item.date)

  const data = {
    labels,
    datasets: [
      {
        label: 'Actual Output',
        data: trend.map((item) => getNumber(item.actual_output, item.actual)),
        borderColor: 'rgb(54, 162, 235)',
        backgroundColor: 'rgba(54, 162, 235, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Planned Output',
        data: trend.map((item) => getNumber(item.planned_output, item.planned)),
        borderColor: 'rgb(255, 206, 86)',
        backgroundColor: 'rgba(255, 206, 86, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Downtime (mins)',
        data: trend.map((item) => getNumber(item.downtime_minutes, item.downtime)),
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        tension: 0.4,
        fill: false,
        yAxisID: 'y1',
      },
    ],
  }

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      mode: 'index',
      intersect: false,
    },
    stacked: false,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Production Trend',
      },
    },
    scales: {
      y: {
        type: 'linear',
        display: true,
        position: 'left',
        beginAtZero: true,
      },
      y1: {
        type: 'linear',
        display: true,
        position: 'right',
        beginAtZero: true,
        grid: {
          drawOnChartArea: false,
        },
      },
    },
  }

  return (
    <Card sx={{ height: '100%' }}>
      <CardContent sx={{ height: 400 }}>
        <Typography variant="h6" gutterBottom>
          Production Output vs Plan
        </Typography>
        <Line data={data} options={options} />
      </CardContent>
    </Card>
  )
}

export default ProductionTrendChart
