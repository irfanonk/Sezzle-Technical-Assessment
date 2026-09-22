const DISPLAY_PRECISION = 12;

/** Trims IEEE-754 noise such as 0.30000000000000004 for display. */
export function formatResult(value: number): string {
  return String(Number.parseFloat(value.toPrecision(DISPLAY_PRECISION)));
}
