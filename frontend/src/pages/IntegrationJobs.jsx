import { useEffect, useState } from 'react'
import { Container, Card, CardContent, Typography, Button, Box, Alert, Chip, Grid } from '@mui/material'
import api from '../services/api'

function IntegrationJobs() {
  const [status, setStatus] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const fetchStatus = async () => {
    try {
      const response = await api.get('/integration/status')
      setStatus(response.data)
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to load status')
    }
  }

  useEffect(() => {
    fetchStatus()
    const interval = setInterval(fetchStatus, 15_000)
    return () => clearInterval(interval)
  }, [])

  const triggerSync = async () => {
    try {
      setLoading(true)
      setError(null)
      await api.post('/integration/sync')
      await fetchStatus()
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to start sync')
    } finally {
      setLoading(false)
    }
  }

  const statusColor = (statusValue) => {
    switch (statusValue) {
      case 'running':
        return 'warning'
      case 'completed':
        return 'success'
      case 'completed_with_errors':
        return 'error'
      default:
        return 'default'
    }
  }

  return (
    <Container maxWidth="md" sx={{ mt: 4, mb: 4 }}>
      <Card>
        <CardContent>
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
            <Typography variant="h5">Data Integration Jobs</Typography>
            <Button variant="contained" onClick={triggerSync} disabled={loading || status?.status === 'running'}>
              {status?.status === 'running' ? 'Sync In Progress' : 'Run Full Sync'}
            </Button>
          </Box>
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}
          {status && (
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Current Status
                </Typography>
                <Chip label={status.status} color={statusColor(status.status)} sx={{ mt: 1 }} />
              </Grid>
              <Grid item xs={12} sm={6}>
                <Typography variant="subtitle2" color="text.secondary">
                  Last Updated
                </Typography>
                <Typography variant="body1" sx={{ mt: 1 }}>
                  {status.last_updated ? new Date(status.last_updated).toLocaleString() : 'N/A'}
                </Typography>
              </Grid>
              {status.started_at && (
                <Grid item xs={12} sm={6}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Started At
                  </Typography>
                  <Typography variant="body1" sx={{ mt: 1 }}>
                    {new Date(status.started_at).toLocaleString()}
                  </Typography>
                </Grid>
              )}
              {status.completed_at && (
                <Grid item xs={12} sm={6}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Completed At
                  </Typography>
                  <Typography variant="body1" sx={{ mt: 1 }}>
                    {new Date(status.completed_at).toLocaleString()}
                  </Typography>
                </Grid>
              )}
              {status.duration && (
                <Grid item xs={12} sm={6}>
                  <Typography variant="subtitle2" color="text.secondary">
                    Duration
                  </Typography>
                  <Typography variant="body1" sx={{ mt: 1 }}>
                    {status.duration}
                  </Typography>
                </Grid>
              )}
              {status.message && (
                <Grid item xs={12}>
                  <Alert severity={status.errors?.length ? 'warning' : 'info'}>{status.message}</Alert>
                </Grid>
              )}
              {status.errors?.length > 0 && (
                <Grid item xs={12}>
                  <Alert severity="warning">
                    {status.errors.map((errMsg, idx) => (
                      <div key={idx}>{errMsg}</div>
                    ))}
                  </Alert>
                </Grid>
              )}
            </Grid>
          )}
        </CardContent>
      </Card>
    </Container>
  )
}

export default IntegrationJobs
