type OperandFieldsProps = {
  values: string[];
  errors: (string | null)[];
  onChange: (index: number, value: string) => void;
};

/** One input per operand; the count comes from the operation's arity. */
export function OperandFields({ values, errors, onChange }: OperandFieldsProps) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:gap-4">
      {values.map((value, index) => {
        const inputId = `operand-${index}`;
        const error = errors[index] ?? null;
        const errorId = `${inputId}-error`;

        return (
          <div key={inputId} className="flex flex-1 flex-col gap-1.5">
            <label htmlFor={inputId} className="text-sm font-medium text-slate-700">
              Value {index + 1}
            </label>
            <input
              id={inputId}
              name={inputId}
              type="text"
              inputMode="decimal"
              autoComplete="off"
              value={value}
              onChange={(event) => onChange(index, event.target.value)}
              aria-invalid={error !== null}
              aria-describedby={error ? errorId : undefined}
              className={`w-full rounded-lg border px-3 py-2.5 text-base text-slate-900 focus:outline-none ${
                error
                  ? 'border-red-400 focus:border-red-500'
                  : 'border-slate-300 focus:border-slate-500'
              }`}
            />
            {error && (
              <p id={errorId} className="text-sm text-red-700">
                {error}
              </p>
            )}
          </div>
        );
      })}
    </div>
  );
}
