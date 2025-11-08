import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { Box, Card, CardContent, Typography, CircularProgress, Alert, Chip, Grid } from '@mui/material';
import api from '../../services/api';

const HRAttritionForecastChart = ({ startDate, endDate, department, days = 30, forecastData, historicalData }) => {
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
          const histResponse = await api.get('/hr/overview', {
            params: { yearMonth: startDate?.slice(0, 7), department },
          });
          setHistorical(histResponse.data?.trend || []);
        }

        if (!forecastData) {
          const response = await api.post('/ai/forecast/hr', {
            forecast_type: 'hr',
            start_date: startDate,
            end_date: endDate,
            days,
            department,
          });
          setForecast(response.data);
        }
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load HR forecast');
      } finally {
        setLoading(false);
      }
    };

    fetchForecast();
  }, [forecastData, historicalData, startDate, endDate, department, days]);

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

  const historicalLabels = (historical || []).map((item) => item.month || item.date || '');
  const historicalValues = (historical || []).map((item) => item.attrition_count ?? item.AttritionCount ?? 0);

  const forecastLabels = forecast.forecast_data.map((item) => item.date);
  const forecastValues = forecast.forecast_data.map((item) => item.forecast);
  const upperBound = forecast.forecast_data.map((item) => item.upper_bound);
  const lowerBound = forecast.forecast_data.map((item) => item.lower_bound);

  const labels = [...historicalLabels, ...forecastLabels];
  const historicalDataset = [...historicalValues, ...new Array(forecastLabels.length).fill(null)];
  const forecastDataset = [...new Array(historicalLabels.length).fill(null), ...forecastValues];

  const data = {
    labels,
    datasets: [
      {
        label: 'Historical Attrition',
        data: historicalDataset,
        borderColor: 'rgb(255, 99, 132)',
        backgroundColor: 'rgba(255, 99, 132, 0.2)',
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Forecast',
        data: forecastDataset,
        borderColor: 'rgb(54, 162, 235)',
        backgroundColor: 'rgba(54, 162, 235, 0.2)',
        borderDash: [6, 6],
        tension: 0.4,
        fill: true,
      },
      {
        label: 'Upper Bound',
        data: [...new Array(historicalLabels.length).fill(null), ...upperBound],
        borderColor: 'rgba(54, 162, 235, 0.3)',
        backgroundColor: 'rgba(54, 162, 235, 0.1)',
        borderDash: [2, 2],
        fill: '-1',
        tension: 0.4,
        pointRadius: 0,
      },
      {
        label: 'Lower Bound',
        data: [...new Array(historicalLabels.length).fill(null), ...lowerBound],
        borderColor: 'rgba(54, 162, 235, 0.3)',
        backgroundColor: 'rgba(54, 162, 235, 0.1)',
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
          <Typography variant="h6">Attrition Forecast</Typography>
          <Chip
            label={`Confidence: ${forecast.confidence_level?.toFixed(1)}%`}
            color="primary"
            size="small"
          />
        </Box>
        <Box height={400}>
          <Line data={data} options={{ responsive: true, maintainAspectRatio: false }} />
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

export default HRAttritionForecastChart;
