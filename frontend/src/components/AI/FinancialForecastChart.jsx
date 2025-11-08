import React, { useEffect, useState } from 'react';
import { Line, Bar } from 'react-chartjs-2';
import { Box, Card, CardContent, Typography, CircularProgress, Alert, Chip, Tabs, Tab } from '@mui/material';
import { api } from '../../services/api';

const FinancialForecastChart = ({ startDate, endDate, budgetCategory, days = 30 }) => {
  const [forecastData, setForecastData] = useState(null);
  const [historicalData, setHistoricalData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [tabValue, setTabValue] = useState(0);

  useEffect(() => {
    fetchForecast();
  }, [startDate, endDate, budgetCategory, days]);

  const fetchForecast = async () => {
    setLoading(true);
    setError(null);
    try {
      const histResponse = await api.get('/finance/budget-vs-actual', {
        params: { startDate, endDate, budgetCategory }
      });
      setHistoricalData(histResponse.data);

      const forecastResponse = await api.post('/ai/forecast/finance', {
        forecast_type: 'finance',
        start_date: startDate,
        end_date: endDate,
        days: days,
        budget_category: budgetCategory
      });
      setForecastData(forecastResponse.data);
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to fetch forecast');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight={400}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  if (!forecastData) {
    return <Alert severity="info">No forecast data available</Alert>;
  }

  const historicalLabels = historicalData?.map(item => 
    new Date(item.month).toLocaleDateString()
  ) || [];
  const budgetValues = historicalData?.map(item => item.budget_amount) || [];
  const actualValues = historicalData?.map(item => item.actual_amount) || [];

  const forecastLabels = forecastData.forecast_data?.map(item => 
    new Date(item.date).toLocaleDateString()
  ) || [];
  const forecastValues = forecastData.forecast_data?.map(item => item.forecast) || [];

  const allLabels = [...historicalLabels, ...forecastLabels];
  const allBudget = [...budgetValues, ...new Array(forecastLabels.length).fill(null)];
  const allActual = [...actualValues, ...new Array(forecastLabels.length).fill(null)];
  const allForecast = [...new Array(historicalLabels.length).fill(null), ...forecastValues];

  const lineChartData = {
    labels: allLabels,
    datasets: [
      {
        label: 'Budget',
        data: allBudget,
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        fill: true,
      },
      {
        label: 'Actual',
        data: allActual,
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        fill: true,
      },
      {
        label: 'AI Forecast',
        data: allForecast,
        borderColor: 'rgb(153, 102, 255)',
        backgroundColor: 'rgba(153, 102, 255, 0.2)',
        borderDash: [5, 5],
        fill: true,
      },
    ],
  };

  const barChartData = {
    labels: forecastLabels,
    datasets: [
      {
        label: 'Forecasted Amount',
        data: forecastValues,
        backgroundColor: 'rgba(153, 102, 255, 0.6)',
        borderColor: 'rgb(153, 102, 255)',
        borderWidth: 1,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top',
      },
      title: {
        display: true,
        text: 'Financial Forecast with AI Predictions',
        font: {
          size: 16,
          weight: 'bold',
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: 'Amount',
        },
      },
    },
  };

  return (
    <Card>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6">Financial Forecast</Typography>
          <Chip 
            label={`Confidence: ${forecastData.confidence_level?.toFixed(1)}%`}
            color="primary"
            size="small"
          />
        </Box>

        <Tabs value={tabValue} onChange={(e, newValue) => setTabValue(newValue)}>
          <Tab label="Trend View" />
          <Tab label="Forecast View" />
        </Tabs>

        <Box height={400} mt={2}>
          {tabValue === 0 ? (
            <Line data={lineChartData} options={chartOptions} />
          ) : (
            <Bar data={barChartData} options={chartOptions} />
          )}
        </Box>

        <Box mt={2}>
          <Typography variant="body2" color="text.secondary">
            Total Forecast: {forecastData.total_forecast?.toLocaleString()} | 
            Avg Daily: {forecastData.average_daily_forecast?.toLocaleString()} | 
            Model: {forecastData.model_used}
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};

export default FinancialForecastChart;

