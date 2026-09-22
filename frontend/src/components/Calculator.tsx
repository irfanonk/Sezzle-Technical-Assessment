import { useState } from 'react';
import type { FormEvent } from 'react';
import { useCalculateMutation } from '../queries/useCalculateMutation';
import { useOperationsQuery } from '../queries/useOperationsQuery';
import { validateOperands } from '../validation/operands';
import { ErrorNotice } from './ErrorNotice';
import { OperandFields } from './OperandFields';
import { OperationSelect } from './OperationSelect';
import { ResultPanel } from './ResultPanel';

export function Calculator() {
  const operationsQuery = useOperationsQuery();
  const calculation = useCalculateMutation();

  const [selectedName, setSelectedName] = useState<string | null>(null);
  const [operandValues, setOperandValues] = useState<string[]>([]);
  const [operandErrors, setOperandErrors] = useState<(string | null)[]>([]);

  const operations = operationsQuery.data ?? [];
  const selectedOperation =
    operations.find((operation) => operation.name === selectedName) ?? operations[0] ?? null;
  const arity = selectedOperation?.arity ?? 0;
  const values = sized(operandValues, arity);

  function handleOperationChange(operationName: string) {
    setSelectedName(operationName);
    setOperandValues([]);
    setOperandErrors([]);
    calculation.reset();
  }

  function handleOperandChange(index: number, value: string) {
    setOperandValues((current) => {
      const next = sized(current, arity);
      next[index] = value;
      return next;
    });

    // A result or error from previous inputs would be misleading once they change.
    if (calculation.isSuccess || calculation.isError) {
      calculation.reset();
    }
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!selectedOperation || calculation.isPending) {
      return;
    }

    const validation = validateOperands(values);
    setOperandErrors(validation.errors);

    if (!validation.isValid) {
      calculation.reset();
      return;
    }

    calculation.mutate({
      operation: selectedOperation.name,
      operands: validation.operands,
    });
  }

  if (operationsQuery.isPending) {
    return (
      <p role="status" className="text-sm text-slate-600">
        Loading operations…
      </p>
    );
  }

  if (operationsQuery.isError) {
    return <ErrorNotice error={operationsQuery.error} onRetry={() => operationsQuery.refetch()} />;
  }

  if (!selectedOperation) {
    return <p className="text-sm text-slate-600">No operations are available right now.</p>;
  }

  return (
    <div className="flex flex-col gap-5">
      <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-4">
        <OperationSelect
          operations={operations}
          value={selectedOperation.name}
          onChange={handleOperationChange}
        />

        <OperandFields values={values} errors={operandErrors} onChange={handleOperandChange} />

        <button
          type="submit"
          disabled={calculation.isPending}
          className="w-full rounded-lg bg-slate-900 px-4 py-3 text-base font-semibold text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-400 sm:w-auto sm:self-start sm:px-6"
        >
          {calculation.isPending ? 'Calculating…' : 'Calculate'}
        </button>
      </form>

      {calculation.isPending && (
        <p role="status" className="text-sm text-slate-600">
          Calculating result…
        </p>
      )}

      {calculation.isError && <ErrorNotice error={calculation.error} />}

      {calculation.isSuccess && <ResultPanel result={calculation.data} />}
    </div>
  );
}

function sized(values: string[], length: number): string[] {
  return Array.from({ length }, (_, index) => values[index] ?? '');
}
