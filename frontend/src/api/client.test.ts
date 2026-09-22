import { describe, expect, it } from 'vitest';
import { apiRequest } from './client';
import { ApiError } from './errors';
import { jsonResponse, stubFetch } from '../test/utils';

describe('apiRequest', () => {
  it('unwraps the data field of the response envelope', async () => {
    stubFetch(async () => jsonResponse({ data: { result: 7 }, error: null }));

    await expect(apiRequest('/calculate')).resolves.toEqual({ result: 7 });
  });

  it('turns an error envelope into an ApiError with the backend code', async () => {
    stubFetch(async () =>
      jsonResponse(
        { data: null, error: { code: 'invalid_operand_count', message: 'needs 2 operands' } },
        400,
      ),
    );

    const error = await apiRequest('/calculate').catch((caught: unknown) => caught);

    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      code: 'invalid_operand_count',
      message: 'needs 2 operands',
      status: 400,
    });
  });

  it('reports a network error when the request cannot be sent', async () => {
    stubFetch(async () => {
      throw new TypeError('Failed to fetch');
    });

    const error = await apiRequest('/operations').catch((caught: unknown) => caught);

    expect(error).toMatchObject({ code: 'network_error', status: null });
  });

  it('reports an unexpected response when the body is not a valid envelope', async () => {
    stubFetch(
      async () =>
        ({
          ok: false,
          status: 502,
          json: async () => {
            throw new SyntaxError('Unexpected token');
          },
        }) as unknown as Response,
    );

    const error = await apiRequest('/operations').catch((caught: unknown) => caught);

    expect(error).toMatchObject({ code: 'unexpected_response', status: 502 });
  });
});
