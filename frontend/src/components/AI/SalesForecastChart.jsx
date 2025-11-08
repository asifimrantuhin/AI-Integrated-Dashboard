import React, { useEffect, useState } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import { Box, Card, CardContent, Typography, CircularProgress, Alert } from '@mui/material';
import api from '../../services/api';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

const SalesForecastChart = ({ startDate, endDate, channelId, days = 30, forecastData, historicalData = [] }) => {
  const [forecast, setForecast] = useState(forecastData || null);
  const [loading, setLoading] = useState(!forecastData);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (forecastData) {
      setForecast(forecastData);
      setLoading(false);
      return;
    }

    const fetchForecast = async () => {
      try {
        setLoading(true);
        const response = await api.post('/ai/forecast/sales', {
          forecast_type: 'sales',
          start_date: startDate,
          end_date: endDate,
          days,
          channel_id: channelId,
        });
        setForecast(response.data);
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load forecast');
      } finally {
        setLoading(false);
      }
    };

    fetchForecast();
  }, [forecastData, startDate, endDate, channelId, days]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={300}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  if (!forecast || !forecast.forecast_data?.length) {
    return null;
  }

  const historicalLabels = historicalData.map((item) => item.month || item.date);
  const historicalValues = historicalData.map((item) => item.billed || item.value || 0);

  const newLabels = forecast.forecast_data.map((item) => item.date);
  const forecastValues = forecast.forecast_data.map((item) => item.forecast);
  const upperBound = forecast.forecast_data.map((item) => item.upper_bound);
  const lowerBound = forecast.forecast_data.map((item) => item.lower_bound);

  const labels = [...historicalLabels, ...newLabels];
  const historicalDataset = [...historicalValues, ...new Array(newLabels.length).fill(null)];
  const forecastDataset = [...new Array(historicalLabels.length).fill(null), ...forecastValues];
  const upperBoundDataset = [...new Array(historicalLabels.length).fill(null), ...upperBound];
  const lowerBoundDataset = [...new Array(historicalLabels.length).fill(null), ...lowerBound];

  const data = {
    labels,
    datasets: [
      {
        label: 'Historical Sales',
        data: historicalDataset,
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Forecast',
        data: forecastDataset,
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        borderDash: [6, 6],
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Forecast Upper',
        data: upperBoundDataset,
        borderColor: 'rgba(255, 99, 132, 0.3)',
        backgroundColor: 'rgba(255, 99, 132, 0.1)',
        borderDash: [2, 2],
        fill: '-1',
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Forecast Lower',
        data: lowerBoundDataset,
        borderColor: 'rgba(255, 99, 132, 0.3)',
        backgroundColor: 'rgba(255, 99, 132, 0.1)',
        borderDash: [2, 2],
        fill: '-1',
        tension: 0.4,
        pointRadius: 0,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Sales Forecast',
      },
    },
    interaction: {
      mode: 'nearest',
      intersect: false,
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  return (
    <Card>
      <CardContent sx={{ height: 400 }}>
        <Typography variant="h6" gutterBottom>
          AI Sales Forecast
        </Typography>
        <Line data={data} options={options} />
      </CardContent>
    </Card>
  );
};

export default SalesForecastChart;

