import { formatResult } from '../lib/formatResult';

type ResultPanelProps = {
  result: number;
};

export function ResultPanel({ result }: ResultPanelProps) {
  return (
    <div className="rounded-lg border border-emerald-300 bg-emerald-50 p-4">
      <p className="text-sm font-medium text-emerald-800">Result</p>
      <output className="block break-words text-2xl font-semibold text-emerald-900 sm:text-3xl">
        {formatResult(result)}
      </output>
    </div>
  );
}
