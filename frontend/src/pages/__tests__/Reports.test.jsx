import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import Reports from '../Reports'
import api from '../../services/api'

vi.mock('../../services/api', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
    },
  }
})

describe('Reports page', () => {
  beforeEach(() => {
    api.get.mockReset()
    api.post.mockReset()
  })

  it('fetches and displays available reports', async () => {
    api.get.mockResolvedValueOnce({
      data: [
        {
          id: 'sales-monthly',
          name: 'Sales Monthly Summary',
          description: 'Channel and product sales by month',
          filters: ['yearMonth'],
          formats: ['json', 'csv'],
        },
      ],
    })

    render(<Reports />)

    await waitFor(() => expect(screen.getByText('Sales Monthly Summary')).toBeInTheDocument())
  })

  it('runs report and renders preview data', async () => {
    api.get
      .mockResolvedValueOnce({
        data: [
          {
            id: 'sales-monthly',
            name: 'Sales Monthly Summary',
            description: 'Channel and product sales by month',
            filters: ['yearMonth'],
            formats: ['json'],
          },
        ],
      })
      .mockResolvedValueOnce({
        data: [
          { channel_name: 'Retail', billed: 1000 },
          { channel_name: 'Wholesale', billed: 2000 },
        ],
      })

    render(<Reports />)

    const selectButton = await screen.findByRole('button', { name: 'Select' })
    fireEvent.click(selectButton)

    const filterInput = screen.getByLabelText('yearMonth')
    fireEvent.change(filterInput, { target: { value: '2024-01' } })

    const runButton = screen.getByRole('button', { name: 'Run Report' })
    fireEvent.click(runButton)

    await waitFor(() => expect(screen.getByText('Retail')).toBeInTheDocument())
    expect(screen.getByText('Wholesale')).toBeInTheDocument()
  })
})
