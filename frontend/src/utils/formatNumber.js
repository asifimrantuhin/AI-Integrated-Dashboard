export function formatCurrencyCrore(value, symbol = '৳') {
  if (value === null || value === undefined) return '—'
  if (typeof value !== 'number') return value
  const abs = Math.abs(value)
  // 1 crore = 10,000,000
  if (abs >= 1e7) {
    const v = value / 1e7
    return `${symbol} ${v.toLocaleString(undefined, { maximumFractionDigits: 2 })} Cr`
  }
  return `${symbol} ${value.toLocaleString(undefined, { maximumFractionDigits: 0 })}`
}

export default formatCurrencyCrore
