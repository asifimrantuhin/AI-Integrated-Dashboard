import { useState } from 'react'
import { Card, CardContent, Typography, TextField, Button, Box, CircularProgress, Alert } from '@mui/material'
import api from '../../services/api'

const BIAssistant = () => {
  const [question, setQuestion] = useState('')
  const [answer, setAnswer] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (!question.trim()) return

    try {
      setLoading(true)
      setError(null)
      setAnswer('')
      const response = await api.post('/bi/assistant', {
        question,
        context: {},
      })
      setAnswer(response.data?.answer || 'No response')
    } catch (err) {
      setError(err.response?.data?.error || 'Assistant is unavailable')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Ask iDash Assistant
        </Typography>
        <Box component="form" onSubmit={handleSubmit} display="flex" gap={2} flexDirection={{ xs: 'column', sm: 'row' }}>
          <TextField
            fullWidth
            placeholder="Ask about business performance..."
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
          />
          <Button type="submit" variant="contained" disabled={loading}>
            {loading ? 'Thinking...' : 'Ask'}
          </Button>
        </Box>
        {loading && (
          <Box display="flex" justifyContent="center" mt={2}>
            <CircularProgress size={24} />
          </Box>
        )}
        {error && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {error}
          </Alert>
        )}
        {answer && !loading && (
          <Alert severity="info" sx={{ mt: 2 }}>
            {answer}
          </Alert>
        )}
      </CardContent>
    </Card>
  )
}

export default BIAssistant
