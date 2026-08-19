import { afterEach, describe, expect, it, vi } from 'vitest';
import { calculate } from './calculator';

afterEach(() => {
  vi.unstubAllGlobals();
});

describe('calculate', () => {
  it('POSTs to the calculate endpoint and returns the response', async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ operation: 'add', operands: [2, 3], result: 5 }),
    } as Response);
    vi.stubGlobal('fetch', fetchMock);

    const res = await calculate({ operation: 'add', operands: [2, 3] });

    expect(res.result).toBe(5);
    expect(fetchMock.mock.calls[0][0]).toContain('/v1/calculate');
  });
});
