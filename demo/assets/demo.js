/* global Chart */
(function () {
  const readyCallbacks = []
  let readyFired = false

  function runCallbacks() {
    if (readyFired) return
    readyFired = true
    while (readyCallbacks.length) {
      const fn = readyCallbacks.shift()
      try {
        fn()
      } catch (err) {
        console.error('iDash demo callback failed', err)
      }
    }
  }

  function onReady(callback) {
    if (document.readyState === 'complete' || document.readyState === 'interactive') {
      setTimeout(callback, 0)
    } else {
      readyCallbacks.push(callback)
    }
  }

  document.addEventListener('DOMContentLoaded', runCallbacks)

  function createChart(el, config) {
    if (!window.Chart || !el) return null
    return new Chart(el, config)
  }

  function formatCurrency(value) {
    if (value === null || value === undefined || Number.isNaN(Number(value))) {
      return '—'
    }
    return `৳${Number(value).toLocaleString('en-US', { maximumFractionDigits: 1 })}`
  }

  function formatMillions(value, decimals = 0) {
    if (value === null || value === undefined || Number.isNaN(Number(value))) {
      return '—'
    }
    const scaled = Number(value) / 1_000_000
    return `৳${scaled.toLocaleString('en-US', {
      minimumFractionDigits: 0,
      maximumFractionDigits: decimals,
    })}M`
  }

  function formatMillionsLabel(value, decimals = 0) {
    if (value === null || value === undefined || Number.isNaN(Number(value))) {
      return '—'
    }
    return `৳${Number(value).toLocaleString('en-US', {
      minimumFractionDigits: 0,
      maximumFractionDigits: decimals,
    })}M`
  }

  function populateTable(id, rows) {
    const table = document.getElementById(id)
    if (!table) return
    const tbody = table.querySelector('tbody')
    if (!tbody) return
    tbody.innerHTML = rows
      .map(
        (row) => `
        <tr>
          ${row.map((col) => `<td>${col}</td>`).join('')}
        </tr>`
      )
      .join('')
  }

  window.iDashDemo = {
    createChart,
    populateTable,
    formatCurrency,
    formatMillions,
    formatMillionsLabel,
    onReady,
  }

  if (!window.__iDashDemoReady) {
    window.__iDashDemoReady = true
    document.dispatchEvent(new Event('iDashDemoReady'))
  }
})()
