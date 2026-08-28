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

// Zgodność klienta z kontraktem: generat, powielone literały nazw i port rdzenia mogą się rozjechać.

const KATALOG_ZRODEL = new URL('.', import.meta.url).pathname;
const KATALOG_BUDOWY = join(KATALOG_ZRODEL, '..', '..');

/**
 * Obszary komend bez własnego zdarzenia `*.unknown`. Wykaz jest pusty: każdy
 * obszar wnoszący komendy ma dziś własne zdarzenie odmowy. Wykaz stoi jako
 * zapora, nie zgoda: obszar bez zdarzenia zapasowego ma być decyzją, nie
 * przeoczeniem.
 */
const OBSZARY_NA_ZDARZENIU_POLACZENIA: readonly string[] = [];

/** Usuwa komentarze blokowe i wierszowe z treści pliku źródłowego, zostawiając wyłącznie sam kod do przeglądu. */
function bezKomentarzy(tresc: string): string {
  return tresc.replace(/\/\*[\s\S]*?\*\//g, '').replace(/(^|[^:])\/\/[^\n]*/g, '$1');
}

/** Zbiera ścieżki wszystkich plików o wskazanych rozszerzeniach w całym przeszukiwanym drzewie katalogów. */
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
    // Pozycja kontraktu nazywa się `typ`, nie `nazwa` — pole niesie nazwę komendy i zdarzenia.
    komendy: { typ: string }[];
    zdarzenia: { typ: string }[];
    // Zdarzeń zapasowych jest jedno na obszar, więc suma wykazu liczy się z liczby obszarów.
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
    // Wykaz zdarzeń to zdarzenia nazwane plus jedno zdarzenie zapasowe na obszar.
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
 * Pliki klienta, w których nazwa kontraktu stoi dziś wprost, zamiast
 * pochodzić z `contract.ts`. Wykaz jest zaporą, nie zgodą — sprawdzian
 * wypada niepomyślnie przy powieleniu spoza wykazu i przy wierszu bez
 * pokrycia.
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
      // Komentarze odpadają przed przeglądem; cytat nazwy kontraktu w komentarzu nie jest powieleniem.
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
