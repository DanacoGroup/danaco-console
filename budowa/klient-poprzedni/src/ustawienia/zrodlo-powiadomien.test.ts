import { describe, expect, it } from 'vitest';

import { Command } from '../../../shared/contract';
import {
  KLASY_ZDARZEN,
  KLUCZ_GLOWNY,
  kluczKanalow,
  kluczKlasy,
  utworzZrodloPowiadomien,
} from './zrodlo-powiadomien';

/**
 * Sprawdziany źródła sekcji Powiadomień. Mierzą trzy rzeczy, z których każda ma
 * cenę w oknie: świeżą bazę czytaną jako „wszystko wyłączone", zapis idący na
 * poziom, którego katalog nie dopuszcza, i piętnaście par pytań tam, gdzie
 * wystarczą dwa.
 */

/** Kanał podstawiony: zapisuje wywołania i oddaje przygotowane odpowiedzi. */
function kanalPodstawiony(odpowiedzi: Record<string, unknown>) {
  const wywolania: { komenda: string; tresc: Record<string, unknown> }[] = [];
  const kanal = {
    wyslij(komenda: string, tresc: Record<string, unknown>, odbierz: (wynik: unknown) => void) {
      wywolania.push({ komenda, tresc });
      const wynik = odpowiedzi[komenda];
      odbierz({ udany: wynik !== undefined, wynik, blad: undefined });
    },
    naZdarzenie: () => () => undefined,
    sesja: () => ({ id: () => 'sesja' }),
  };
  return { kanal: kanal as never, wywolania };
}

/** Definicja katalogu w kształcie kontraktu, w minimalnym wypełnieniu. */
function definicja(klucz: string, domyslna: unknown, opcje: { value: string; label: string }[] = []) {
  return {
    key: klucz,
    name: klucz,
    valueType: 'bool',
    defaultValue: domyslna,
    allowedScopes: ['global'],
    allowedAxes: [],
    options: opcje,
    order: 0,
    enabled: true,
    requiresRestart: false,
  };
}

/** Katalog wszystkich piętnastu kluczy sekcji. */
function katalogPelny() {
  const kanaly = [
    { value: 'mobile', label: 'Mobile' },
    { value: 'email', label: 'E-mail' },
  ];
  return {
    definitions: [
      definicja(KLUCZ_GLOWNY, true),
      ...KLASY_ZDARZEN.map((klasa) => definicja(kluczKlasy(klasa), true)),
      ...KLASY_ZDARZEN.map((klasa) =>
        definicja(kluczKanalow(klasa), klasa === 'system' ? ['email'] : [], kanaly),
      ),
    ],
  };
}

describe('źródło sekcji Powiadomienia', () => {
  it('świeża baza czyta się wartościami domyślnymi katalogu, nie zerami', async () => {
    // Zapisu nie ma ani jednego — tak wygląda platforma przed pierwszą zmianą.
    const { kanal } = kanalPodstawiony({
      [Command.SettingsDefinitionList]: katalogPelny(),
      [Command.ConfigGet]: { entries: [] },
    });

    const stan = await utworzZrodloPowiadomien(kanal).odczytaj();

    expect(stan.odmowa).toBe('');
    expect(stan.wlaczone).toBe(true);
    expect(stan.klasy).toHaveLength(7);
    expect(stan.klasy.every((klasa) => klasa.czynna)).toBe(true);
    expect(stan.klasy.find((k) => k.klasa === 'system')?.kanaly).toEqual(['email']);
    expect(stan.klasy.find((k) => k.klasa === 'termin')?.kanaly).toEqual([]);
  });

  it('zapis Operatora przesłania wartość domyślną', async () => {
    const { kanal } = kanalPodstawiony({
      [Command.SettingsDefinitionList]: katalogPelny(),
      [Command.ConfigGet]: {
        entries: [
          { key: kluczKlasy('termin'), value: false, scope: 'global' },
          { key: kluczKanalow('termin'), value: ['mobile'], scope: 'global' },
        ],
      },
    });

    const stan = await utworzZrodloPowiadomien(kanal).odczytaj();

    expect(stan.klasy.find((k) => k.klasa === 'termin')?.czynna).toBe(false);
    expect(stan.klasy.find((k) => k.klasa === 'termin')?.kanaly).toEqual(['mobile']);
  });

  it('pyta rdzeń dwa razy, nie piętnaście', async () => {
    const { kanal, wywolania } = kanalPodstawiony({
      [Command.SettingsDefinitionList]: katalogPelny(),
      [Command.ConfigGet]: { entries: [] },
    });

    await utworzZrodloPowiadomien(kanal).odczytaj();

    expect(wywolania).toHaveLength(2);
    expect(wywolania.map((w) => w.komenda)).toEqual([
      Command.SettingsDefinitionList,
      Command.ConfigGet,
    ]);
  });

  it('katalog bez kategorii mówi wprost, że nie ma czym sterować', async () => {
    const { kanal } = kanalPodstawiony({
      [Command.SettingsDefinitionList]: { definitions: [] },
      [Command.ConfigGet]: { entries: [] },
    });

    const stan = await utworzZrodloPowiadomien(kanal).odczytaj();

    expect(stan.odmowa).toContain('migracji 377');
    expect(stan.klasy).toHaveLength(0);
  });

  it('zapis idzie na poziom globalny — jedyny, który okno Ustawień umie wskazać', async () => {
    const { kanal, wywolania } = kanalPodstawiony({
      [Command.ConfigSet]: { applied: true },
    });

    await utworzZrodloPowiadomien(kanal).ustawKanaly('blad', ['mobile', 'email']);

    expect(wywolania).toHaveLength(1);
    expect(wywolania[0]?.tresc).toMatchObject({
      key: kluczKanalow('blad'),
      value: ['mobile', 'email'],
      scope: 'global',
    });
  });
});
