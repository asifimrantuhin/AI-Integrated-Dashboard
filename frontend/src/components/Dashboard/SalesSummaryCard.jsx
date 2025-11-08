import { useEffect, useState } from 'react'
import { Card, CardContent, Typography, Grid, Box, CircularProgress } from '@mui/material'
import api from '../../services/api'
import { SalesLineChart, SalesBarChart } from '../Charts/SalesChart'

function SalesSummaryCard() {
  const [loading, setLoading] = useState(true)
  const [salesData, setSalesData] = useState(null)
  const [chartData, setChartData] = useState(null)

  useEffect(() => {
    fetchSalesData()
  }, [])

  const fetchSalesData = async () => {
    try {
      setLoading(true)
      const response = await api.get('/sales/summary')
      setSalesData(response.data)
      
      // Fetch chart data
      const chartResponse = await api.get('/sales/cumulative?yearMonth=2024')
      setChartData(chartResponse.data)
    } catch (error) {
      console.error('Error fetching sales data:', error)
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="center" alignItems="center" minHeight={200}>
            <CircularProgress />
          </Box>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Sales Summary
        </Typography>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={6} md={3}>
            <Typography variant="body2" color="textSecondary">
              Lifting
            </Typography>
            <Typography variant="h6">
              {salesData?.lifting?.toLocaleString() || 0}
            </Typography>
          </Grid>
          <Grid item xs={6} md={3}>
            <Typography variant="body2" color="textSecondary">
              IMS
            </Typography>
            <Typography variant="h6">
              {salesData?.ims?.toLocaleString() || 0}
            </Typography>
          </Grid>
          <Grid item xs={6} md={3}>
            <Typography variant="body2" color="textSecondary">
              Primary Collection
            </Typography>
            <Typography variant="h6">
              {salesData?.primary_collection?.toLocaleString() || 0}
            </Typography>
          </Grid>
          <Grid item xs={6} md={3}>
            <Typography variant="body2" color="textSecondary">
              Market Collection
            </Typography>
            <Typography variant="h6">
              {salesData?.market_collection?.toLocaleString() || 0}
            </Typography>
          </Grid>
        </Grid>
        {chartData && chartData.length > 0 && (
          <Box sx={{ mt: 3 }}>
            <SalesLineChart data={chartData} title="Sales Trend" />
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

export default SalesSummaryCard

