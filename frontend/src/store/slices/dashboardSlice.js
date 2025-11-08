import { createSlice } from '@reduxjs/toolkit'

const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState: {
    salesData: null,
    productionData: null,
    financeData: null,
    inventoryData: null,
    hrData: null,
    supplyChainData: null,
    loading: false,
    error: null,
  },
  reducers: {
    setSalesData: (state, action) => {
      state.salesData = action.payload
    },
    setProductionData: (state, action) => {
      state.productionData = action.payload
    },
    setFinanceData: (state, action) => {
      state.financeData = action.payload
    },
    setInventoryData: (state, action) => {
      state.inventoryData = action.payload
    },
    setHRData: (state, action) => {
      state.hrData = action.payload
    },
    setSupplyChainData: (state, action) => {
      state.supplyChainData = action.payload
    },
    setLoading: (state, action) => {
      state.loading = action.payload
    },
    setError: (state, action) => {
      state.error = action.payload
    },
  },
})

export const {
  setSalesData,
  setProductionData,
  setFinanceData,
  setInventoryData,
  setHRData,
  setSupplyChainData,
  setLoading,
  setError,
} = dashboardSlice.actions

export default dashboardSlice.reducer

