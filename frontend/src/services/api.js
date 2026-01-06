import axios from 'axios'

// Default to relative '/api' so Vite dev server proxy (configured in vite.config.js)
// forwards requests to the backend at http://localhost:8080. If you need to
// target the backend directly (e.g. when serving built files), set
// `VITE_API_BASE` in environment.
const defaultBase = typeof import.meta !== 'undefined' && import.meta.env
  ? (import.meta.env.VITE_API_BASE || '/api')
  : '/api'

const api = axios.create({
  baseURL: defaultBase,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Attach a lightweight request id to help with tracing
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    const requestId = typeof window !== 'undefined' && window.crypto?.randomUUID
      ? window.crypto.randomUUID()
      : `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
    config.headers['x-request-id'] = requestId
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    if (error.response?.status === 429) {
      console.warn('Rate limit exceeded. Please retry shortly.')
    }
    return Promise.reject(error)
  }
)

export default api

