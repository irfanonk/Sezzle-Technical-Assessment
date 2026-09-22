import type { Operation } from '../api/types';

type OperationSelectProps = {
  operations: Operation[];
  value: string;
  onChange: (operationName: string) => void;
};

export function OperationSelect({ operations, value, onChange }: OperationSelectProps) {
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor="operation" className="text-sm font-medium text-slate-700">
        Operation
      </label>
      <select
        id="operation"
        name="operation"
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-base text-slate-900 focus:border-slate-500 focus:outline-none"
      >
        {operations.map((operation) => (
          <option key={operation.name} value={operation.name}>
            {operation.label} ({operation.symbol})
          </option>
        ))}
      </select>
    </div>
  );
}
