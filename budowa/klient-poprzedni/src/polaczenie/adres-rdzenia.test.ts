import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

// Ustalenie adresu gniazda sprawdza trzy wejścia: budowanie, powłokę i lokalizację dokumentu.

const PORT_DOMYSLNY = '17870';

/** Podstawia lokalizację dokumentu widzianą przez moduł podczas testu ustalania adresu gniazda rdzenia. */
function zLokalizacja(lokalizacja: Partial<Location> | undefined): void {
  if (lokalizacja === undefined) {
    vi.stubGlobal('location', undefined);
    return;
  }
  vi.stubGlobal('location', lokalizacja);
}

/** Wczytuje moduł adresu rdzenia od nowa, żeby stan adresu z powłoki startował pusty w każdym przypadku. */
async function swiezyModul(): Promise<typeof import('./adres-rdzenia')> {
  vi.resetModules();
  return import('./adres-rdzenia');
}

describe('adres rdzenia wyliczony z lokalizacji dokumentu', () => {
  beforeEach(() => {
    vi.unstubAllGlobals();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('bierze host i port ze strony, którą serwuje rdzeń', async () => {
    zLokalizacja({ protocol: 'http:', hostname: '127.0.0.1', port: '19000' });
    const { adresRdzenia } = await swiezyModul();
    expect(adresRdzenia()).toBe('ws://127.0.0.1:19000/ws');
  });

  it('przechodzi na warstwę szyfrowaną razem ze stroną', async () => {
    zLokalizacja({ protocol: 'https:', hostname: 'konsola.example', port: '8443' });
    const { adresRdzenia } = await swiezyModul();
    expect(adresRdzenia()).toBe('wss://konsola.example:8443/ws');
  });

  it('bierze port domyślny rdzenia, gdy strona go nie ma', async () => {
    // Strona bez portu stoi na 80 albo 443, gdzie rdzeń nie nasłuchuje — port strony byłby adresem ślepym.
    zLokalizacja({ protocol: 'http:', hostname: 'konsola.example', port: '' });
    const { adresRdzenia } = await swiezyModul();
    expect(adresRdzenia()).toBe(`ws://konsola.example:${PORT_DOMYSLNY}/ws`);
  });

  it('schodzi na pętlę zwrotną, gdy strona nie ma źródła sieciowego', async () => {
    // Tak wygląda pakiet osadzony w powłoce natywnej.
    zLokalizacja({ protocol: 'file:', hostname: '', port: '' });
    const { adresRdzenia } = await swiezyModul();
    expect(adresRdzenia()).toBe(`ws://127.0.0.1:${PORT_DOMYSLNY}/ws`);
  });

  it('schodzi na pętlę zwrotną przy braku lokalizacji', async () => {
    zLokalizacja(undefined);
    const { adresRdzenia } = await swiezyModul();
    expect(adresRdzenia()).toBe(`ws://127.0.0.1:${PORT_DOMYSLNY}/ws`);
  });
});

describe('adres podany przez powłokę natywną', () => {
  beforeEach(() => {
    vi.unstubAllGlobals();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('przekłada adres HTTP powłoki na adres gniazda', async () => {
    zLokalizacja({ protocol: 'http:', hostname: 'tauri.localhost', port: '' });
    const { adresRdzenia, przyjmijAdresZPowloki } = await swiezyModul();

    expect(przyjmijAdresZPowloki('http://127.0.0.1:19000')).toBe(true);
    expect(adresRdzenia()).toBe('ws://127.0.0.1:19000/ws');
  });

  it('przekłada adres szyfrowany na gniazdo szyfrowane', async () => {
    zLokalizacja({ protocol: 'file:', hostname: '', port: '' });
    const { adresRdzenia, przyjmijAdresZPowloki } = await swiezyModul();

    expect(przyjmijAdresZPowloki('https://konsola.example:8443')).toBe(true);
    expect(adresRdzenia()).toBe('wss://konsola.example:8443/ws');
  });

  it('odmawia adresu spoza HTTP i zostaje przy wyliczeniu', async () => {
    zLokalizacja({ protocol: 'http:', hostname: '127.0.0.1', port: '19000' });
    const { adresRdzenia, przyjmijAdresZPowloki } = await swiezyModul();

    expect(przyjmijAdresZPowloki('ftp://127.0.0.1:19000')).toBe(false);
    expect(adresRdzenia(), 'adres nie do rozłożenia wstrzymał wyliczenie').toBe(
      'ws://127.0.0.1:19000/ws',
    );
  });

  it('odmawia zapisu nie do rozłożenia i zostaje przy wyliczeniu', async () => {
    zLokalizacja({ protocol: 'http:', hostname: '127.0.0.1', port: '19000' });
    const { adresRdzenia, przyjmijAdresZPowloki } = await swiezyModul();

    expect(przyjmijAdresZPowloki('to nie jest adres')).toBe(false);
    expect(adresRdzenia()).toBe('ws://127.0.0.1:19000/ws');
  });
});
