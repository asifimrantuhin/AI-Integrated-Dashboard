import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { Box, Card, CardContent, Typography, CircularProgress, Alert, Chip, Grid } from '@mui/material';
import api from '../../services/api';

const ProductionForecastChart = ({ startDate, endDate, factoryId, days = 30, forecastData, historicalData }) => {
  const [forecast, setForecast] = useState(forecastData || null);
  const [historical, setHistorical] = useState(historicalData || []);
  const [loading, setLoading] = useState(!forecastData);
  const [error, setError] = useState(null);

  useEffect(() => {
    if (forecastData) {
      setForecast(forecastData);
    }
    if (historicalData) {
      setHistorical(historicalData);
    }
  }, [forecastData, historicalData]);

  useEffect(() => {
    if (forecastData && historicalData) {
      return;
    }

    const fetchForecast = async () => {
      try {
        setLoading(true);
        setError(null);

        if (!historicalData) {
          const histResponse = await api.get('/production/analysis', {
            params: { month: startDate?.slice(0, 7), factory: factoryId },
          });
          setHistorical(histResponse.data || []);
        }

        if (!forecastData) {
          const forecastResponse = await api.post('/ai/forecast/production', {
            forecast_type: 'production',
            start_date: startDate,
            end_date: endDate,
            days,
            factory_id: factoryId,
          });
          setForecast(forecastResponse.data);
        }
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load production forecast');
      } finally {
        setLoading(false);
      }
    };

    fetchForecast();
  }, [forecastData, historicalData, startDate, endDate, factoryId, days]);

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

  const historicalLabels = (historical || []).map((item) =>
    item.month ? new Date(item.month).toLocaleDateString() : item.production_date ? new Date(item.production_date).toLocaleDateString() : ''
  );
  const historicalValues = (historical || []).map((item) =>
    item.cmonthly_amount ?? item.actual_output ?? item.actual_output ?? 0
  );

  const forecastLabels = forecast.forecast_data?.map((item) => new Date(item.date).toLocaleDateString()) || [];
  const forecastValues = forecast.forecast_data?.map((item) => item.forecast) || [];
  const upperBound = forecast.forecast_data?.map((item) => item.upper_bound) || [];
  const lowerBound = forecast.forecast_data?.map((item) => item.lower_bound) || [];

  const labels = [...historicalLabels, ...forecastLabels];
  const historicalDataset = [...historicalValues, ...new Array(forecastLabels.length).fill(null)];
  const forecastDataset = [...new Array(historicalLabels.length).fill(null), ...forecastValues];

  const chartData = {
    labels,
    datasets: [
      {
        label: 'Historical Production',
        data: historicalDataset,
        borderColor: 'rgb(54, 162, 235)',
        backgroundColor: 'rgba(54, 162, 235, 0.2)',
        fill: true,
        tension: 0.4,
      },
      {
        label: 'AI Forecast',
        data: forecastDataset,
        borderColor: 'rgb(255, 159, 64)',
        backgroundColor: 'rgba(255, 159, 64, 0.2)',
        borderDash: [5, 5],
        fill: true,
        tension: 0.4,
      },
      {
        label: 'Upper Bound',
        data: [...new Array(historicalLabels.length).fill(null), ...upperBound],
        borderColor: 'rgba(255, 159, 64, 0.3)',
        backgroundColor: 'rgba(255, 159, 64, 0.1)',
        borderDash: [2, 2],
        fill: '-1',
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Lower Bound',
        data: [...new Array(historicalLabels.length).fill(null), ...lowerBound],
        borderColor: 'rgba(255, 159, 64, 0.3)',
        backgroundColor: 'rgba(255, 159, 64, 0.1)',
        borderDash: [2, 2],
        fill: '-1',
        tension: 0.4,
        pointRadius: 0,
      },
    ],
  };

  return (
    <Card>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6">Production Forecast</Typography>
          <Chip
            label={`Confidence: ${forecast.confidence_level?.toFixed(1)}%`}
            color="primary"
            size="small"
          />
        </Box>

        <Box height={400}>
          <Line data={chartData} options={{ responsive: true, maintainAspectRatio: false }} />
        </Box>

        <Grid container spacing={2} mt={2}>
          <Grid item xs={4}>
            <Typography variant="body2" color="text.secondary">
              Total Forecast: {forecast.total_forecast?.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <Typography variant="body2" color="text.secondary">
              Avg Daily: {forecast.average_daily_forecast?.toLocaleString()}
            </Typography>
          </Grid>
          <Grid item xs={4}>
            <Typography variant="body2" color="text.secondary">
              Model: {forecast.model_used}
            </Typography>
          </Grid>
        </Grid>
      </CardContent>
    </Card>
  );
};

export default ProductionForecastChart;

