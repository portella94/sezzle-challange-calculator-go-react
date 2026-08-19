import { afterEach, describe, expect, it, vi } from 'vitest';
import { CalculatorApiError, postJson } from './client';

function mockFetch(response: Partial<Response> & { json: () => Promise<unknown> }) {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue(response as Response));
}

afterEach(() => {
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

describe('postJson', () => {
  it('returns parsed JSON on a 2xx response', async () => {
    mockFetch({ ok: true, json: () => Promise.resolve({ result: 5 }) });

    const data = await postJson<{ a: number }, { result: number }>('/v1/calculate', { a: 1 });
    expect(data).toEqual({ result: 5 });
  });

  it('sends the body as JSON with the correct headers', async () => {
    const fetchMock = vi.fn().mockResolvedValue({ ok: true, json: () => Promise.resolve({}) } as Response);
    vi.stubGlobal('fetch', fetchMock);

    await postJson('/v1/calculate', { operation: 'add', operands: [1, 2] });

    expect(fetchMock).toHaveBeenCalledOnce();
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit];
    expect(url).toContain('/v1/calculate');
    expect(init.method).toBe('POST');
    expect(JSON.parse(init.body as string)).toEqual({ operation: 'add', operands: [1, 2] });
  });

  it('throws CalculatorApiError with the backend error code on non-2xx', async () => {
    mockFetch({
      ok: false,
      json: () => Promise.resolve({ error: { code: 'division_by_zero', message: 'division by zero' } }),
    });

    await expect(postJson('/v1/calculate', {})).rejects.toMatchObject({
      code: 'division_by_zero',
      message: 'division by zero',
    });
  });

  it('throws a network_error when fetch rejects', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')));

    await expect(postJson('/v1/calculate', {})).rejects.toBeInstanceOf(CalculatorApiError);
    await expect(postJson('/v1/calculate', {})).rejects.toMatchObject({ code: 'network_error' });
  });
});
