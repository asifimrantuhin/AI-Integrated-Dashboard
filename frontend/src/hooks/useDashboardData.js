import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import api from '../services/api'

const CACHE_PREFIX = 'idash:dashboard:'
const CACHE_TTL = 1000 * 60 * 5 // 5 minutes

let safeSession = null
try {
  if (typeof window !== 'undefined') {
    safeSession = window.sessionStorage
  }
} catch (error) {
  safeSession = null
}

const readCache = (key) => {
  if (!safeSession) return null
  try {
    const raw = safeSession.getItem(key)
    if (!raw) return null
    const parsed = JSON.parse(raw)
    if (!parsed?.timestamp || Date.now() - parsed.timestamp > CACHE_TTL) {
      safeSession.removeItem(key)
      return null
    }
    return parsed
  } catch (error) {
    safeSession?.removeItem(key)
    return null
  }
}

const writeCache = (key, value) => {
  if (!safeSession) return
  try {
    safeSession.setItem(
      key,
      JSON.stringify({
        payload: value,
        timestamp: Date.now(),
      })
    )
  } catch (error) {
    safeSession?.removeItem(key)
  }
}

export const useDashboardData = (endpoint) => {
  const cacheKey = useMemo(() => `${CACHE_PREFIX}${endpoint}`, [endpoint])
  const [data, setData] = useState({ kpis: [], charts: {}, ai_insights: [], last_updated: null })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [stale, setStale] = useState(false)
  const [refreshing, setRefreshing] = useState(false)
  const abortRef = useRef(null)

  const fetchData = useCallback(
    async ({ force = false } = {}) => {
      if (!force && abortRef.current) return

      if (abortRef.current) {
        abortRef.current.abort()
      }

      const controller = new AbortController()
      abortRef.current = controller

      try {
        if (force) {
          setRefreshing(true)
        } else {
          setLoading(true)
        }
        setError(null)

        const response = await api.get(endpoint, { signal: controller.signal })
        const payload = response.data
        setData(payload)
        writeCache(cacheKey, payload)
        setStale(false)
      } catch (err) {
        if (err.name === 'CanceledError' || err.name === 'AbortError') {
          return
        }
        setError(err.response?.data?.error || 'Unable to load dashboard data')
      } finally {
        setLoading(false)
        setRefreshing(false)
        abortRef.current = null
      }
    },
    [endpoint, cacheKey]
  )

  useEffect(() => {
    const cached = readCache(cacheKey)
    if (cached?.payload) {
      setData(cached.payload)
      setLoading(false)
      const age = Date.now() - cached.timestamp
      // Consider cache stale if older than half the TTL
      if (age > CACHE_TTL / 2) {
        setStale(true)
        fetchData({ force: true })
      }
    } else {
      fetchData()
    }

    return () => {
      if (abortRef.current) {
        abortRef.current.abort()
      }
    }
  }, [cacheKey, fetchData])

  return {
    data,
    loading,
    error,
    stale,
    refreshing,
    refresh: () => fetchData({ force: true }),
  }
}

export default useDashboardData
