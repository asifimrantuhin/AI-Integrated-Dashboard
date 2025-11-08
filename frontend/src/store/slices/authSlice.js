import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import api from '../../services/api'

export const login = createAsyncThunk(
  'auth/login',
  async ({ email, password, companyId }, { rejectWithValue }) => {
    try {
      const response = await api.post('/auth/login', { email, password, company_id: companyId })
      const { token, user, roles, permissions, default_company_id, company_ids } = response.data
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify({ user, roles, permissions, default_company_id, company_ids }))
      return response.data
    } catch (error) {
      return rejectWithValue(error.response?.data?.error || 'Login failed')
    }
  }
)

export const checkAuth = createAsyncThunk(
  'auth/checkAuth',
  async (_, { rejectWithValue }) => {
    try {
      const response = await api.get('/auth/user')
      return response.data
    } catch (error) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      return rejectWithValue('Not authenticated')
    }
  }
)

export const logout = createAsyncThunk('auth/logout', async () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
})

const storedUser = JSON.parse(localStorage.getItem('user') || 'null')

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    user: storedUser?.user || null,
    roles: storedUser?.roles || [],
    permissions: storedUser?.permissions || [],
    defaultCompanyId: storedUser?.default_company_id || null,
    companyIds: storedUser?.company_ids || [],
    token: localStorage.getItem('token') || null,
    isAuthenticated: !!localStorage.getItem('token'),
    loading: false,
    error: null,
  },
  reducers: {
    clearError: (state) => {
      state.error = null
    },
    setDefaultCompany: (state, action) => {
      state.defaultCompanyId = action.payload
      const stored = JSON.parse(localStorage.getItem('user') || 'null')
      if (stored) {
        stored.default_company_id = action.payload
        localStorage.setItem('user', JSON.stringify(stored))
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(login.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading = false
        state.isAuthenticated = true
        state.user = action.payload.user
        state.roles = action.payload.roles || []
        state.permissions = action.payload.permissions || []
        state.defaultCompanyId = action.payload.default_company_id || null
        state.companyIds = action.payload.company_ids || []
        state.token = action.payload.token
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload
      })
      .addCase(checkAuth.fulfilled, (state, action) => {
        state.isAuthenticated = true
        state.user = action.payload.user
        state.roles = action.payload.roles || []
        state.permissions = action.payload.permissions || []
        state.defaultCompanyId = action.payload.default_company_id || null
        state.companyIds = action.payload.company_ids || []
        const stored = {
          user: state.user,
          roles: state.roles,
          permissions: state.permissions,
          default_company_id: state.defaultCompanyId,
          company_ids: state.companyIds,
        }
        localStorage.setItem('user', JSON.stringify(stored))
      })
      .addCase(checkAuth.rejected, (state) => {
        state.isAuthenticated = false
        state.user = null
        state.roles = []
        state.permissions = []
        state.defaultCompanyId = null
        state.companyIds = []
        state.token = null
      })
      .addCase(logout.fulfilled, (state) => {
        state.isAuthenticated = false
        state.user = null
        state.roles = []
        state.permissions = []
        state.defaultCompanyId = null
        state.companyIds = []
        state.token = null
      })
  },
})

export const { clearError, setDefaultCompany } = authSlice.actions
export default authSlice.reducer

