import { fireEvent, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';
import { Calculator } from './Calculator';
import { calculateCalls, deferred, jsonResponse, renderWithClient, stubFetch } from '../test/utils';

const OPERATIONS = [
  { name: 'add', label: 'Addition', symbol: '+', arity: 2 },
  { name: 'divide', label: 'Division', symbol: '÷', arity: 2 },
  { name: 'square_root', label: 'Square Root', symbol: '√', arity: 1 },
];

const operationsEnvelope = { data: { operations: OPERATIONS }, error: null };

function isOperationsRequest(input: RequestInfo | URL) {
  return String(input).endsWith('/operations');
}

/** Operations always succeed; only the calculation response varies per test. */
function stubBackend(calculate: (init?: RequestInit) => Promise<Response>) {
  return stubFetch(async (input, init) => {
    if (isOperationsRequest(input)) {
      return jsonResponse(operationsEnvelope);
    }
    return calculate(init);
  });
}

async function fillOperands(user: ReturnType<typeof userEvent.setup>, values: string[]) {
  const inputs = screen.getAllByRole('textbox');
  for (const [index, value] of values.entries()) {
    const input = inputs[index];
    if (input) {
      await user.type(input, value);
    }
  }
}

describe('operations loading', () => {
  it('shows a loading state while operations are loading', () => {
    stubFetch(() => new Promise<Response>(() => {}));

    renderWithClient(<Calculator />);

    expect(screen.getByRole('status')).toHaveTextContent('Loading operations');
  });

  it('renders the operations returned by the backend', async () => {
    stubFetch(async () => jsonResponse(operationsEnvelope));

    renderWithClient(<Calculator />);

    const select = await screen.findByLabelText('Operation');
    expect(within(select).getAllByRole('option').map((option) => option.textContent)).toEqual([
      'Addition (+)',
      'Division (÷)',
      'Square Root (√)',
    ]);
  });

  it('generates one operand input per arity of the selected operation', async () => {
    const user = userEvent.setup();
    stubFetch(async () => jsonResponse(operationsEnvelope));

    renderWithClient(<Calculator />);

    await screen.findByLabelText('Operation');
    expect(screen.getAllByRole('textbox')).toHaveLength(2);

    await user.selectOptions(screen.getByLabelText('Operation'), 'square_root');
    expect(screen.getAllByRole('textbox')).toHaveLength(1);
  });

  it('shows an error and recovers when retrying', async () => {
    const user = userEvent.setup();
    const fetchMock = stubFetch(async () => {
      throw new TypeError('Failed to fetch');
    });

    renderWithClient(<Calculator />);

    expect(await screen.findByRole('alert')).toHaveTextContent(
      'Cannot reach the calculator service',
    );

    fetchMock.mockImplementation(async () => jsonResponse(operationsEnvelope));
    await user.click(screen.getByRole('button', { name: 'Retry' }));

    expect(await screen.findByLabelText('Operation')).toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
  });
});

describe('calculation', () => {
  it('blocks the request when operands are missing or not numbers', async () => {
    const user = userEvent.setup();
    const fetchMock = stubBackend(async () => jsonResponse({ data: { result: 0 }, error: null }));

    renderWithClient(<Calculator />);
    await screen.findByLabelText('Operation');

    await user.click(screen.getByRole('button', { name: 'Calculate' }));
    expect(screen.getAllByText('Enter a value.')).toHaveLength(2);
    expect(calculateCalls(fetchMock)).toHaveLength(0);

    await fillOperands(user, ['abc', '2']);
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(screen.getByText('Enter a valid number.')).toBeInTheDocument();
    expect(calculateCalls(fetchMock)).toHaveLength(0);
  });

  it('shows a pending state while the calculation is in flight', async () => {
    const user = userEvent.setup();
    const pending = deferred<Response>();
    stubBackend(() => pending.promise);

    renderWithClient(<Calculator />);
    await screen.findByLabelText('Operation');

    await fillOperands(user, ['2', '3']);
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Calculating result');
    expect(screen.getByRole('button', { name: 'Calculating…' })).toBeDisabled();

    pending.resolve(jsonResponse({ data: { result: 5 }, error: null }));
    expect(await screen.findByText('5')).toBeInTheDocument();
  });

  it('sends the calculation and displays the result', async () => {
    const user = userEvent.setup();
    const fetchMock = stubBackend(async () =>
      jsonResponse({ data: { result: 5 }, error: null }),
    );

    renderWithClient(<Calculator />);
    await screen.findByLabelText('Operation');

    await fillOperands(user, ['2', '3']);
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByText('5')).toBeInTheDocument();

    const [, init] = calculateCalls(fetchMock)[0] as [string, RequestInit];
    expect(JSON.parse(String(init.body))).toEqual({ operation: 'add', operands: [2, 3] });
  });

  it('displays the error returned by the backend', async () => {
    const user = userEvent.setup();
    stubBackend(async () =>
      jsonResponse(
        { data: null, error: { code: 'division_by_zero', message: 'cannot divide by zero' } },
        400,
      ),
    );

    renderWithClient(<Calculator />);
    await screen.findByLabelText('Operation');

    await user.selectOptions(screen.getByLabelText('Operation'), 'divide');
    await fillOperands(user, ['1', '0']);
    await user.click(screen.getByRole('button', { name: 'Calculate' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('cannot divide by zero');
  });

  it('prevents duplicate submissions while a calculation is pending', async () => {
    const user = userEvent.setup();
    const pending = deferred<Response>();
    const fetchMock = stubBackend(() => pending.promise);

    renderWithClient(<Calculator />);
    await screen.findByLabelText('Operation');

    await fillOperands(user, ['2', '3']);
    const submitButton = screen.getByRole('button', { name: 'Calculate' });
    await user.click(submitButton);

    await waitFor(() => expect(calculateCalls(fetchMock)).toHaveLength(1));

    await user.click(screen.getByRole('button', { name: 'Calculating…' }));
    const form = screen.getByRole('button', { name: 'Calculating…' }).closest('form');
    if (form) {
      fireEvent.submit(form);
    }

    expect(calculateCalls(fetchMock)).toHaveLength(1);

    pending.resolve(jsonResponse({ data: { result: 5 }, error: null }));
    expect(await screen.findByText('5')).toBeInTheDocument();
  });
});
