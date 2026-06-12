import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import { createRequire } from 'node:module';
import { resolve } from 'node:path';
import { afterEach, beforeEach, describe, it } from 'node:test';

const require = createRequire(import.meta.url);
const ts = require('../../frontend/node_modules/typescript');
const productHandlers = await loadProductHandlers();
const originalFetch = globalThis.fetch;
const originalConsoleError = console.error;

async function loadProductHandlers() {
  const sourcePath = resolve('frontend/src/catalog/productHandlers.ts');
  const source = await readFile(sourcePath, 'utf8');
  const { outputText } = ts.transpileModule(source, {
    compilerOptions: {
      module: ts.ModuleKind.ES2022,
      target: ts.ScriptTarget.ES2022,
    },
  });

  const moduleUrl = `data:text/javascript;base64,${Buffer.from(outputText).toString('base64')}`;
  return import(moduleUrl);
}

function response(payload, ok = true) {
  return {
    ok,
    async json() {
      return payload;
    },
  };
}

function mockFetchWith(result) {
  const calls = [];
  globalThis.fetch = async (...args) => {
    calls.push(args);
    if (result instanceof Error) {
      throw result;
    }
    return result;
  };
  return calls;
}

describe('catalog product handlers', () => {
  beforeEach(() => {
    console.error = () => {};
  });

  afterEach(() => {
    globalThis.fetch = originalFetch;
    console.error = originalConsoleError;
  });

  it('searches KLS products with an encoded query and requested limit', async () => {
    const payload = [
      {
        code: '12-345-67-89',
        eng: 'Straight forceps',
        viet: 'Mo ta',
        brand: 'Martin',
      },
    ];
    const calls = mockFetchWith(response(payload));

    const result = await productHandlers.searchKLSProduct('straight forceps/large', 25);

    assert.deepEqual(calls, [
      ['http://localhost:8000/catalog/kls/?q=straight%20forceps%2Flarge&limit=25'],
    ]);
    assert.deepEqual(result, payload);
  });

  it('returns an empty KLS result when the backend response is not ok', async () => {
    mockFetchWith(response({ detail: 'server error' }, false));

    const result = await productHandlers.searchKLSProduct('forceps', 50);

    assert.deepEqual(result, []);
  });

  it('searches Aesculap products and returns parsed JSON', async () => {
    const payload = [
      {
        code: 'AB123',
        eng: 'Bone lever',
        viet: 'Mo ta',
        image: 'AB123.jpg',
        alternative: 'ALT-123',
        brand: 'B-Braun',
      },
    ];
    const calls = mockFetchWith(response(payload));

    const result = await productHandlers.searchAesculapProduct('AB 123', 10);

    assert.deepEqual(calls, [
      ['http://localhost:8000/catalog/aes/?q=AB%20123&limit=10'],
    ]);
    assert.deepEqual(result, payload);
  });

  it('returns an empty Aesculap result and logs fetch failures', async () => {
    const errors = [];
    console.error = (...args) => errors.push(args);
    mockFetchWith(new Error('network unavailable'));

    const result = await productHandlers.searchAesculapProduct('AB123', 10);

    assert.deepEqual(result, []);
    assert.equal(errors.length, 1);
    assert.equal(errors[0][0], 'Error searching Aesculap product:');
    assert.match(errors[0][1].message, /network unavailable/);
  });

  it('fetches KLS image metadata for an encoded product code', async () => {
    const payload = {
      code: '12/345',
      img1_url: 'one.jpg',
      img2_url: null,
      img3_url: null,
    };
    const calls = mockFetchWith(response(payload));

    const result = await productHandlers.getKLSImages('12/345');

    assert.deepEqual(calls, [
      ['http://localhost:8000/catalog/kls/images/12%2F345'],
    ]);
    assert.deepEqual(result, payload);
  });

  it('returns null when KLS image lookup fails', async () => {
    mockFetchWith(response({ detail: 'not found' }, false));

    const result = await productHandlers.getKLSImages('missing');

    assert.equal(result, null);
  });
});
