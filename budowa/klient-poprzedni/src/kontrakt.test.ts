import { readFileSync, readdirSync } from 'node:fs';
import { join } from 'node:path';

import { describe, expect, it } from 'vitest';

import {
  KOMENDY,
  PROTOCOL_VERSION,
  ZDARZENIA,
  ZDARZENIA_NIEZNANEJ,
  czyKomenda,
  czyZdarzenie,
} from '../../shared/contract';

/**
 * Zgodność klienta z kontraktem.
 *
 * Kontrakt jest jedynym źródłem prawdy nazw, a `contract.ts` jego wytworem.
 * Trzy rzeczy mogą się tu rozejść po cichu i żadnej nie wychwyci
 * ani kompilator, ani przegląd:
 *
 *  1. Generat starszy od źródła — `contract.json` zmienione, generator
 *     niepuszczony. Kompilacja przechodzi, bo stała nadal istnieje.
 *  2. Literał nazwy powielony w kodzie klienta zamiast wzięty z generatu.
 *     Zmiana nazwy w kontrakcie zostawia wtedy w interfejsie martwe wywołanie,
 *     które rdzeń odbije jako `*.unknown`.
 *  3. Port rdzenia zaszyty w kliencie rozjechany z portem domyślnym rdzenia —
 *     klient szuka gniazda tam, gdzie nikt nie nasłuchuje.
 */

const KATALOG_ZRODEL = new URL('.', import.meta.url).pathname;
const KATALOG_BUDOWY = join(KATALOG_ZRODEL, '..', '..');

/**
 * Obszary komend bez własnego zdarzenia `*.unknown`.
 *
 * Wykaz jest pusty: każdy obszar wnoszący komendy ma dziś własne zdarzenie
 * odmowy, więc nierozpoznana komenda wraca pod nazwą swojego obszaru, a nie
 * pod nazwą połączenia. Odmowa poczty przedstawia się jako sprawa poczty.
 *
 * Wykaz stoi tu jako zapora, nie jako zgoda: obszar wniesiony bez zdarzenia
 * zapasowego ma być decyzją, a nie przeoczeniem, więc sprawdzian wypada
 * niepomyślnie zarówno wtedy, gdy taki obszar się pojawi, jak i wtedy, gdy
 * wiersz zostanie tu po obszarze już domkniętym.
 */
const OBSZARY_NA_ZDARZENIU_POLACZENIA: readonly string[] = [];

/** Usuwa komentarze blokowe i wierszowe, zostawiając sam kod. */
function bezKomentarzy(tresc: string): string {
  return tresc.replace(/\/\*[\s\S]*?\*\//g, '').replace(/(^|[^:])\/\/[^\n]*/g, '$1');
}

/** Zbiera ścieżki wszystkich plików o wskazanych rozszerzeniach w drzewie. */
function plikiDrzewa(katalog: string, rozszerzenia: readonly string[]): string[] {
  const znalezione: string[] = [];
  for (const wpis of readdirSync(katalog, { withFileTypes: true })) {
    const sciezka = join(katalog, wpis.name);
    if (wpis.isDirectory()) {
      znalezione.push(...plikiDrzewa(sciezka, rozszerzenia));
      continue;
    }
    if (rozszerzenia.some((rozszerzenie) => wpis.name.endsWith(rozszerzenie))) {
      znalezione.push(sciezka);
    }
  }
  return znalezione;
}

describe('generat kontraktu odpowiada źródłu prawdy', () => {
  const zrodlo = JSON.parse(
    readFileSync(join(KATALOG_BUDOWY, 'shared', 'contract.json'), 'utf8'),
  ) as {
    kontrakt: { protokol: string };
    // Pozycja kontraktu nazywa się `typ`, nie `nazwa` — to samo pole niesie
    // nazwę komendy i nazwę zdarzenia, bo obie jadą jedną kopertą.
    komendy: { typ: string }[];
    zdarzenia: { typ: string }[];
    // Zdarzeń zapasowych jest jedno na obszar, więc suma wykazu liczy się
    // z obszarów — bez nich sprawdzian nie ma jak rozpoznać, którą z dwóch grup
    // generator zgubił.
    obszary: string[];
    zdarzeniaNieznane: Record<string, unknown>;
  };

  it('niesie tyle komend, ile deklaruje kontrakt', () => {
    expect(KOMENDY).toHaveLength(zrodlo.komendy.length);
  });

  it('niesie zdarzenia nazwane wraz ze zdarzeniami zapasowymi obszarów', () => {
    const wygenerowane = new Set<string>(ZDARZENIA);
    for (const zdarzenie of zrodlo.zdarzenia) {
      expect(wygenerowane.has(zdarzenie.typ), `zdarzenie ${zdarzenie.typ}`).toBe(true);
    }
    // Wykaz zdarzeń to zdarzenia nazwane plus jedno zdarzenie zapasowe na
    // obszar — rozjazd tej sumy znaczy, że generator zgubił jedną z dwóch grup.
    expect(ZDARZENIA).toHaveLength(zrodlo.zdarzenia.length + zrodlo.obszary.length);
    expect(Object.keys(ZDARZENIA_NIEZNANEJ)).toHaveLength(zrodlo.obszary.length);
  });

  it('podaje wersję protokołu ze źródła', () => {
    expect(PROTOCOL_VERSION).toBe(zrodlo.kontrakt.protokol);
  });

  it('rozpoznaje każdą nazwę ze źródła i żadnej spoza niego', () => {
    for (const komenda of KOMENDY) {
      expect(czyKomenda(komenda), `komenda ${komenda}`).toBe(true);
    }
    for (const zdarzenie of ZDARZENIA) {
      expect(czyZdarzenie(zdarzenie), `zdarzenie ${zdarzenie}`).toBe(true);
    }
    expect(czyKomenda('session.wymyslona')).toBe(false);
    expect(czyZdarzenie('session.wymyslone')).toBe(false);
  });

  it('nie miesza wykazu komend z wykazem zdarzeń', () => {
    const zdarzenia = new Set<string>(ZDARZENIA);
    const wspolne = (KOMENDY as readonly string[]).filter((nazwa) => zdarzenia.has(nazwa));
    expect(wspolne, 'nazwy stojące w obu wykazach').toEqual([]);
  });

  it('kieruje obszar bez własnego zdarzenia na zdarzenie połączenia', () => {
    const obszaryZeZdarzeniem = new Set(Object.keys(ZDARZENIA_NIEZNANEJ));
    const bezZdarzenia = new Set<string>();
    for (const komenda of KOMENDY as readonly string[]) {
      const obszar = komenda.split('.')[0];
      if (!obszaryZeZdarzeniem.has(obszar)) bezZdarzenia.add(obszar);
    }
    expect([...bezZdarzenia].sort()).toEqual(OBSZARY_NA_ZDARZENIU_POLACZENIA);
  });
});

/**
 * Pliki klienta, w których nazwa kontraktu stoi dziś wprost, zamiast pochodzić
 * z `contract.ts`.
 *
 * Kontrakt nazywa to błędem wstrzymującym etap i ma rację: zmiana nazwy
 * w `contract.json` przechodzi wtedy przez kompilację obu stron i zostawia
 * w interfejsie martwe wywołanie, które rdzeń odbije zdarzeniem `*.unknown`.
 *
 * Wykaz jest zaporą, nie zgodą. Sprawdzian wypada niepomyślnie zarówno wtedy,
 * gdy powielenie pojawi się w pliku spoza wykazu, jak i wtedy, gdy plik
 * z wykazu zostanie posprzątany, a wiersz zostanie. Dług nie rośnie i nie znika
 * po cichu.
 */
const PLIKI_Z_POWIELONYM_LITERALEM = [
  'aod/rozpoznanie-decyzji.ts',
  'mobile/pozycje-decyzji.ts',
  'moduly/agents/archiwum-ekspertow.ts',
  'moduly/agents/okno-connectors-manager.ts',
  'moduly/apps/etykiety-apps.ts',
  'moduly/multitasking/braki-kontraktu.ts',
  'moduly/multitasking/okno-coordinator-chat.ts',
  'moduly/multitasking/okno-executor-chat.ts',
  'moduly/multitasking/okno-results-analyzer.ts',
  'moduly/multitasking/panel-subagent-network.ts',
  'moduly/multitasking/plan-etapow.ts',
  'moduly/multitasking/sekcja-monitor.ts',
  'moduly/roundtable/widok-panelu-debaty.ts',
  'okna-pomocnicze/rejestr-pomocniczych.ts',
];

describe('klient nie powiela literałów nazw kontraktu', () => {
  it('nie dokłada powielenia w pliku spoza wykazu ani nie zostawia wiersza po sprzątnięciu', () => {
    const nazwyKontraktu = new Set<string>([
      ...(KOMENDY as readonly string[]),
      ...(ZDARZENIA as readonly string[]),
    ]);
    const literal = /(['"`])([a-z][a-zA-Z]*(?:\.[a-z][a-zA-Z]*)+)\1/g;

    const zPowieleniem = new Set<string>();
    for (const plik of plikiDrzewa(KATALOG_ZRODEL, ['.ts'])) {
      if (plik.endsWith('.test.ts')) continue;
      // Komentarze odpadają przed przeglądem. Opisy w tym drzewie cytują nazwy
      // kontraktu w odwrotnych apostrofach i mają to robić — cytat w komentarzu
      // nie jest powieleniem, bo nic z niego nie idzie na drut.
      const tresc = bezKomentarzy(readFileSync(plik, 'utf8'));
      for (const trafienie of tresc.matchAll(literal)) {
        if (!nazwyKontraktu.has(trafienie[2])) continue;
        zPowieleniem.add(plik.slice(KATALOG_ZRODEL.length));
      }
    }

    expect(
      [...zPowieleniem].sort(),
      'pliki z nazwą kontraktu zapisaną wprost zamiast wziętą z contract.ts',
    ).toEqual(PLIKI_Z_POWIELONYM_LITERALEM);
  });
});

describe('port rdzenia w kliencie zgadza się z portem rdzenia', () => {
  it('stała klienta równa się PortDomyslny z konfiguracji rdzenia', () => {
    const zrodloGo = readFileSync(
      join(KATALOG_BUDOWY, 'server', 'internal', 'konfiguracja', 'ustawienia.go'),
      'utf8',
    );
    const zRdzenia = /PortDomyslny\s*=\s*(\d+)/.exec(zrodloGo);
    expect(zRdzenia, 'stała PortDomyslny w źródle rdzenia').not.toBeNull();

    const zrodloKlienta = readFileSync(
      join(KATALOG_ZRODEL, 'polaczenie', 'adres-rdzenia.ts'),
      'utf8',
    );
    const zKlienta = /PORT_RDZENIA_LOKALNEGO\s*=\s*'(\d+)'/.exec(zrodloKlienta);
    expect(zKlienta, 'stała PORT_RDZENIA_LOKALNEGO w źródle klienta').not.toBeNull();

    expect(zKlienta?.[1]).toBe(zRdzenia?.[1]);
  });
});
