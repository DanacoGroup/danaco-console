import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

/**
 * Ustalenie adresu gniazda.
 *
 * To jest miejsce, w którym klient albo trafia w rdzeń, albo szuka go tam,
 * gdzie nikt nie nasłuchuje — a wtedy ponawianie milczy bez końca i wygląda
 * jak awaria rdzenia, choć jest pomyłką adresu. Rozstrzygnięcie ma trzy
 * wejścia (wskazanie z budowania, adres z powłoki, lokalizacja dokumentu)
 * i jedno pierwszeństwo, więc sprawdzane są wszystkie drogi.
 *
 * Adres z powłoki jest stanem modułu, dlatego każdy przypadek wczytuje moduł
 * na nowo — inaczej mierzyłby ślad po przypadku poprzednim.
 */

const PORT_DOMYSLNY = '17870';

/** Podstawia lokalizację dokumentu widzianą przez moduł. */
function zLokalizacja(lokalizacja: Partial<Location> | undefined): void {
  if (lokalizacja === undefined) {
    vi.stubGlobal('location', undefined);
    return;
  }
  vi.stubGlobal('location', lokalizacja);
}

/** Wczytuje moduł od nowa, żeby stan adresu z powłoki startował pusty. */
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
    // Strona bez portu stoi na 80 albo 443, a tam rdzeń nie nasłuchuje — port
    // ze strony byłby wtedy adresem ślepym.
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
