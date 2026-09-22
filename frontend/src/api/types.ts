/** Shared response envelope returned by every backend endpoint. */
export type ApiResponse<T> = {
  data: T | null;
  error: ApiErrorBody | null;
};

export type ApiErrorBody = {
  code: string;
  message: string;
};

/** Operation metadata advertised by GET /api/operations. */
export type Operation = {
  name: string;
  label: string;
  symbol: string;
  arity: number;
};

export type OperationsData = {
  operations: Operation[];
};

export type CalculationData = {
  result: number;
};

export type CalculationInput = {
  operation: string;
  operands: number[];
};
