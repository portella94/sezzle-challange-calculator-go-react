import { afterEach, describe, expect, it, vi } from 'vitest';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ThemeProvider } from '@mui/material';
import { CalculatorForm } from './CalculatorForm';
import { theme } from '../theme';
import { calculate } from '../api/calculator';
import { CalculatorApiError } from '../api/client';

// Mock the API boundary so the form is tested in isolation from the network.
vi.mock('../api/calculator', () => ({ calculate: vi.fn() }));
const calculateMock = vi.mocked(calculate);

function renderForm() {
  return render(
    <ThemeProvider theme={theme}>
      <CalculatorForm />
    </ThemeProvider>,
  );
}

async function selectOperation(user: ReturnType<typeof userEvent.setup>, label: string) {
  await user.click(screen.getByRole('combobox', { name: /operation/i }));
  await user.click(await screen.findByRole('option', { name: label }));
}

afterEach(() => {
  vi.clearAllMocks();
});

describe('CalculatorForm', () => {
  it('renders two operand inputs for a binary operation by default', () => {
    renderForm();
    expect(screen.getByLabelText('A')).toBeInTheDocument();
    expect(screen.getByLabelText('B')).toBeInTheDocument();
  });

  it('shows a single operand input when a unary operation is selected', async () => {
    const user = userEvent.setup();
    renderForm();

    await selectOperation(user, 'Square root (√)');

    expect(screen.getByLabelText('Value')).toBeInTheDocument();
    expect(screen.queryByLabelText('B')).not.toBeInTheDocument();
  });

  it('validates required operands and does not call the API', async () => {
    const user = userEvent.setup();
    renderForm();

    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(await screen.findAllByText('This field is required')).toHaveLength(2);
    expect(calculateMock).not.toHaveBeenCalled();
  });

  it('rejects non-numeric input client-side', async () => {
    const user = userEvent.setup();
    renderForm();

    await user.type(screen.getByLabelText('A'), 'abc');
    await user.type(screen.getByLabelText('B'), '2');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    expect(await screen.findByText('Enter a valid number')).toBeInTheDocument();
    expect(calculateMock).not.toHaveBeenCalled();
  });

  it('submits a valid calculation and renders the result', async () => {
    const user = userEvent.setup();
    calculateMock.mockResolvedValue({ operation: 'add', operands: [2, 3], result: 5 });
    renderForm();

    await user.type(screen.getByLabelText('A'), '2');
    await user.type(screen.getByLabelText('B'), '3');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    const result = await screen.findByTestId('result');
    expect(result).toHaveTextContent('5');
    expect(calculateMock).toHaveBeenCalledWith({ operation: 'add', operands: [2, 3] });
  });

  it('surfaces a backend error in an alert', async () => {
    const user = userEvent.setup();
    calculateMock.mockRejectedValue(new CalculatorApiError('division_by_zero', 'division by zero'));
    renderForm();

    await selectOperation(user, 'Divide (÷)');
    await user.type(screen.getByLabelText('A'), '1');
    await user.type(screen.getByLabelText('B'), '0');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    const alert = await screen.findByRole('alert');
    expect(within(alert).getByText('division by zero')).toBeInTheDocument();
    expect(screen.queryByTestId('result')).not.toBeInTheDocument();
  });

  it('sends only the single operand for a unary operation', async () => {
    const user = userEvent.setup();
    calculateMock.mockResolvedValue({ operation: 'sqrt', operands: [16], result: 4 });
    renderForm();

    await selectOperation(user, 'Square root (√)');
    await user.type(screen.getByLabelText('Value'), '16');
    await user.click(screen.getByRole('button', { name: /calculate/i }));

    await waitFor(() => expect(calculateMock).toHaveBeenCalledWith({ operation: 'sqrt', operands: [16] }));
  });
});
