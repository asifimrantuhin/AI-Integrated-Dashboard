import { render, screen } from '@testing-library/react'
import KPIWidget from '../KPIWidget'

describe('KPIWidget', () => {
  it('renders title and value', () => {
    render(<KPIWidget title="Revenue" value={12345} change={12.5} prefix="৳ " />)

    expect(screen.getByText('Revenue')).toBeInTheDocument()
    expect(screen.getByText('৳ 12,345')).toBeInTheDocument()
    expect(screen.getByText('+12.5%')).toBeInTheDocument()
  })
})
