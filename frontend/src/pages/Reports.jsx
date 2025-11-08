import { useEffect, useMemo, useState } from 'react'
import {
  Container,
  Card,
  CardContent,
  Typography,
  Grid,
  Button,
  MenuItem,
  TextField,
  Box,
  CircularProgress,
  Alert,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
} from '@mui/material'
import api from '../services/api'

function Reports() {
  const [reports, setReports] = useState([])
  const [selectedReport, setSelectedReport] = useState(null)
  const [format, setFormat] = useState('json')
  const [filters, setFilters] = useState({})
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [previewData, setPreviewData] = useState([])

  useEffect(() => {
    const fetchReports = async () => {
      try {
        const response = await api.get('/reports')
        setReports(response.data)
      } catch (err) {
        setError(err.response?.data?.error || 'Failed to load reports')
      }
    }

    fetchReports()
  }, [])

  const handleRun = async () => {
    if (!selectedReport) return
    setLoading(true)
    setError(null)
    try {
      const params = new URLSearchParams({ format })
      Object.entries(filters).forEach(([key, value]) => {
        if (value) params.append(key, value)
      })

      if (format === 'csv') {
        const response = await api.get(`/reports/${selectedReport.id}?${params.toString()}`, {
          responseType: 'blob',
        })
        const blob = new Blob([response.data], { type: 'text/csv' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.setAttribute('download', `${selectedReport.id}.csv`)
        document.body.appendChild(link)
        link.click()
        link.parentNode.removeChild(link)
      } else {
        const response = await api.get(`/reports/${selectedReport.id}?${params.toString()}`)
        setPreviewData(Array.isArray(response.data) ? response.data : [])
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Failed to generate report')
    } finally {
      setLoading(false)
    }
  }

  const availableFormats = useMemo(() => selectedReport?.formats || ['json'], [selectedReport])

  return (
    <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Reports
              </Typography>
              {reports.map((report) => (
                <Box key={report.id} mb={2}>
                  <Typography variant="subtitle1">{report.name}</Typography>
                  <Typography variant="body2" color="text.secondary">
                    {report.description}
                  </Typography>
                  <Button
                    size="small"
                    sx={{ mt: 1 }}
                    variant={selectedReport?.id === report.id ? 'contained' : 'outlined'}
                    onClick={() => {
                      setSelectedReport(report)
                      setPreviewData([])
                      setFilters({})
                    }}
                  >
                    {selectedReport?.id === report.id ? 'Selected' : 'Select'}
                  </Button>
                </Box>
              ))}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={8}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Generate Report
              </Typography>
              {!selectedReport && (
                <Alert severity="info">Select a report on the left to begin.</Alert>
              )}
              {selectedReport && (
                <Box>
                  <Typography variant="subtitle1" gutterBottom>
                    {selectedReport.name}
                  </Typography>
                  <Grid container spacing={2} mb={2}>
                    {selectedReport.filters.map((filterKey) => (
                      <Grid item xs={12} sm={6} key={filterKey}>
                        <TextField
                          fullWidth
                          label={filterKey}
                          value={filters[filterKey] || ''}
                          onChange={(e) => setFilters((prev) => ({ ...prev, [filterKey]: e.target.value }))}
                        />
                      </Grid>
                    ))}
                    <Grid item xs={12} sm={4}>
                      <TextField
                        select
                        fullWidth
                        label="Format"
                        value={format}
                        onChange={(e) => setFormat(e.target.value)}
                      >
                        {availableFormats.map((fmt) => (
                          <MenuItem key={fmt} value={fmt}>
                            {fmt.toUpperCase()}
                          </MenuItem>
                        ))}
                      </TextField>
                    </Grid>
                  </Grid>
                  <Button variant="contained" onClick={handleRun} disabled={loading}>
                    {loading ? 'Generating...' : 'Run Report'}
                  </Button>
                </Box>
              )}
              {error && (
                <Alert severity="error" sx={{ mt: 2 }}>
                  {error}
                </Alert>
              )}
              {loading && (
                <Box display="flex" justifyContent="center" mt={2}>
                  <CircularProgress size={24} />
                </Box>
              )}
              {format === 'json' && previewData.length > 0 && (
                <Box mt={3}>
                  <Typography variant="subtitle1" gutterBottom>
                    Preview (first 100 rows)
                  </Typography>
                  <Table size="small">
                    <TableHead>
                      <TableRow>
                        {Object.keys(previewData[0] || {}).map((key) => (
                          <TableCell key={key}>{key}</TableCell>
                        ))}
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {previewData.slice(0, 100).map((row, idx) => (
                        <TableRow key={idx}>
                          {Object.keys(previewData[0] || {}).map((key) => (
                            <TableCell key={key}>{row[key]}</TableCell>
                          ))}
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Container>
  )
}

export default Reports
