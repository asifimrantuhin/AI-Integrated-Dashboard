import { Line, Bar } from 'react-chartjs-2'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

export function SalesLineChart({ data, title }) {
  const chartData = {
    labels: data?.map(item => item.month || item.date) || [],
    datasets: [
      {
        label: 'Sales',
        data: data?.map(item => item.lifting || item.ims || item.value) || [],
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        fill: true,
        tension: 0.4
      }
    ]
  }

  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: title || 'Sales Chart'
      }
    },
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }

  return <Line data={chartData} options={options} />
}

export function SalesBarChart({ data, title }) {
  const chartData = {
    labels: data?.map(item => item.label || item.name) || [],
    datasets: [
      {
        label: 'Sales',
        data: data?.map(item => item.value || item.amount) || [],
        backgroundColor: 'rgba(54, 162, 235, 0.6)',
        borderColor: 'rgba(54, 162, 235, 1)',
        borderWidth: 1
      }
    ]
  }

  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: title || 'Sales Chart'
      }
    },
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }

  return <Bar data={chartData} options={options} />
}

export function SalesComparisonChart({ data1, data2, label1, label2, title }) {
  const labels = data1?.map((_, index) => `Month ${index + 1}`) || []
  
  const chartData = {
    labels,
    datasets: [
      {
        label: label1 || 'Current',
        data: data1?.map(item => item.value || item.amount) || [],
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        fill: false
      },
      {
        label: label2 || 'Previous',
        data: data2?.map(item => item.value || item.amount) || [],
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        fill: false
      }
    ]
  }

  const options = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: title || 'Sales Comparison'
      }
    },
    scales: {
      y: {
        beginAtZero: true
      }
    }
  }

  return <Line data={chartData} options={options} />
}

