import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';

import { utworzScenaWejscia } from './scena-wejscia';

/**
 * Sprawdziany sceny wejścia pilnują trzech rzeczy, z których każda ma cenę
 * w produkcie: sceny, która nie staje (Operator patrzy na pustą stronę), sceny,
 * która nie schodzi (zasłona nad gotowym produktem) i sceny, która miga
 * (czyta się jak usterka obrazu).
 */
describe('scena wejścia', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    document.body.replaceChildren();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  /** Przesuwa czas i domyka mikrozadania obietnic. */
  async function uplyw(ms: number): Promise<void> {
    await vi.advanceTimersByTimeAsync(ms);
  }

  it('staje nad dokumentem dopiero po uruchomieniu', () => {
    const scena = utworzScenaWejscia({ gotowosc: new Promise(() => undefined) });

    expect(document.querySelector('.la-scena')).toBeNull();

    scena.uruchom();

    expect(document.querySelector('.la-scena')).not.toBeNull();
  });

  it('nie miga: gotowość natychmiastowa zostawia scenę na próg widoczności', async () => {
    utworzScenaWejscia({ gotowosc: Promise.resolve() }).uruchom();

    await uplyw(100);
    expect(document.querySelector('.la-scena')).not.toBeNull();

    await uplyw(1000);
    expect(document.querySelector('.la-scena')).toBeNull();
  });

  it('schodzi po domknięciu pierwszego odczytu strony', async () => {
    let domknij = (): void => undefined;
    const gotowosc = new Promise<void>((spelnij) => {
      domknij = spelnij;
    });
    utworzScenaWejscia({ gotowosc }).uruchom();

    await uplyw(2000);
    expect(document.querySelector('.la-scena')).not.toBeNull();

    domknij();
    await uplyw(1000);
    expect(document.querySelector('.la-scena')).toBeNull();
  });

  it('odmowa pierwszego odczytu zdejmuje scenę tak samo jak powodzenie', async () => {
    utworzScenaWejscia({ gotowosc: Promise.reject(new Error('rdzeń odmówił')) }).uruchom();

    await uplyw(1000);

    expect(document.querySelector('.la-scena')).toBeNull();
  });

  it('gotowość, która nie nadchodzi, nie zamienia sceny w zasłonę', async () => {
    utworzScenaWejscia({ gotowosc: new Promise(() => undefined) }).uruchom();

    await uplyw(7000);
    expect(document.querySelector('.la-scena')).not.toBeNull();

    await uplyw(2000);
    expect(document.querySelector('.la-scena')).toBeNull();
  });

  it('zdejmij kończy scenę bez czekania na gotowość', async () => {
    const scena = utworzScenaWejscia({ gotowosc: new Promise(() => undefined) });
    scena.uruchom();

    scena.zdejmij();
    await uplyw(1000);

    expect(document.querySelector('.la-scena')).toBeNull();
  });

  it('bryła niesie sześć ścian ze znakiem marki', () => {
    const scena = utworzScenaWejscia({ gotowosc: new Promise(() => undefined) });
    scena.uruchom();

    const sciany = document.querySelectorAll('.la-sciana');

    expect(sciany).toHaveLength(6);
    for (const sciana of sciany) expect(sciana.querySelector('svg')).not.toBeNull();
  });
});
