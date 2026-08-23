import { beforeEach, describe, expect, it, vi } from 'vitest';

import {
  Command,
  TelemetryFormat,
  UsageDimension,
  type UsageAggregate,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzStanDiagnostyki } from './stan-diagnostyki';
import { utworzZakladkeZuzycia } from './zuzycie-zakladka';
import { utworzZrodloDiagnostics } from './zrodlo-diagnostics';
import { utworzZrodloZuzycia } from './zuzycie-zrodlo';

/**
 * Zużycie i koszt — zakładka Usage & Cost kontenera Observability Tools.
 *
 * Sprawdzian pilnuje czterech rzeczy stanowiących o odbiorze: że zestawienie ma
 * drogę z okna i idzie jednym wymiarem, że PUSTY okres ma zdanie, a nie puste
 * miejsce, że koszt niepełny mówi o sobie, i że raport rozliczeniowy WYTWARZA
 * plik nazwany, a nie ciszę po naciśnięciu.
 */

/** Pozycja zestawienia w postaci, w której rdzeń ją oddaje. */
function pozycja(zmiany: Partial<UsageAggregate> = {}): UsageAggregate {
  return {
    dimension: UsageDimension.Channel,
    dimensionId: 'kanal-1',
    dimensionLabel: 'Kanał główny',
    requests: 40,
    promptTokens: 1000,
    completionTokens: 500,
    totalTokens: 1500,
    ...zmiany,
  };
}

/** Kanał próbny: zapamiętuje żądania i oddaje odpowiedź wskazaną per komenda. */
function kanalProbny(odpowiedzi: Record<string, unknown>): {
  kanal: Kanal;
  wyslane: { komenda: string; zadanie: unknown }[];
} {
  const wyslane: { komenda: string; zadanie: unknown }[] = [];
  const kanal = {
    wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push({ komenda, zadanie });
      const tresc = odpowiedzi[komenda];
      przyWyniku?.(
        tresc === undefined
          ? { udany: false, blad: { code: 'not_found', message: 'brak wykonawcy', retryable: false } }
          : { udany: true, wynik: tresc },
      );
      return `zadanie-${String(wyslane.length)}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({ id: () => 'sesja-1' }) as unknown as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Zakładka złożona wraz z kanałem próbnym. */
function zakladka(odpowiedzi: Record<string, unknown>) {
  const { kanal, wyslane } = kanalProbny(odpowiedzi);
  const stan = utworzStanDiagnostyki(utworzZrodloDiagnostics(kanal), {
    zakres: { od: 1_700_000_000_000, do: 1_700_003_600_000 },
  });
  const zbudowana = utworzZakladkeZuzycia(utworzZrodloZuzycia(kanal), stan);
  return { zakladka: zbudowana, element: zbudowana.pozycja.element, wyslane };
}

/** Czeka na rozstrzygnięcie obietnic w kolejce zadań mikro. */
async function przemiel(): Promise<void> {
  await Promise.resolve();
  await Promise.resolve();
  await Promise.resolve();
}

describe('zakładka Usage & Cost', () => {
  it('stoi w kontenerze pod nazwą z opracowania', () => {
    const { zakladka: zbudowana } = zakladka({});
    expect(zbudowana.pozycja.nazwa).toBe('Usage & Cost');
    expect(zbudowana.pozycja.kod).toBe('usage-cost');
  });

  it('nie pyta rdzenia przed pierwszym odczytem', () => {
    const { zakladka: zbudowana, wyslane } = zakladka({});
    expect(zbudowana.czytano()).toBe(false);
    expect(wyslane).toHaveLength(0);
  });

  it('woła usage.summary.get jednym wymiarem wraz z okresem modułu', async () => {
    const { zakladka: zbudowana, wyslane } = zakladka({
      [Command.UsageSummaryGet]: {
        aggregates: [pozycja()],
        fromTime: 1_700_000_000_000,
        toTime: 1_700_003_600_000,
      },
    });

    zbudowana.odswiez();
    await przemiel();

    expect(wyslane).toHaveLength(1);
    expect(wyslane[0]?.komenda).toBe(Command.UsageSummaryGet);
    expect(wyslane[0]?.zadanie).toEqual({
      dimension: UsageDimension.Channel,
      fromTime: 1_700_000_000_000,
      toTime: 1_700_003_600_000,
    });
    expect(zbudowana.czytano()).toBe(true);
  });

  it('pokazuje pozycję wraz z tokenami, a brak cennika nazywa, nie zeruje', async () => {
    const { element, zakladka: zbudowana } = zakladka({
      [Command.UsageSummaryGet]: {
        aggregates: [pozycja({ costWithoutPrice: 7 })],
        fromTime: 1_700_000_000_000,
        toTime: 1_700_003_600_000,
        priceCoverage: 82,
      },
    });

    zbudowana.odswiez();
    await przemiel();

    const tresc = element.textContent ?? '';
    expect(tresc).toContain('Kanał główny');
    expect(tresc).toContain('tokenów łącznie 1500');
    expect(tresc).toContain('kosztu nie ma, bo kanał nie ma cennika');
    // Koszt niepełny musi to powiedzieć — inaczej liczba wygląda na rachunek.
    expect(tresc).toContain('NIEPEŁNY');
    expect(tresc).toContain('Poza kosztem zostało wywołań: 7');
  });

  it('pusty okres ma zdanie, które odróżnia go od braku pomiaru', async () => {
    const { element, zakladka: zbudowana } = zakladka({
      [Command.UsageSummaryGet]: {
        aggregates: [],
        fromTime: 1_700_000_000_000,
        toTime: 1_700_003_600_000,
      },
    });

    zbudowana.odswiez();
    await przemiel();

    const tresc = element.textContent ?? '';
    expect(tresc).toContain('Zużycia w tym okresie nie było');
    expect(tresc).toContain('pusty okres, nie brak pomiaru');
  });

  it('odmowę rdzenia pokazuje wprost', async () => {
    const { element, zakladka: zbudowana } = zakladka({});

    zbudowana.odswiez();
    await przemiel();

    expect(element.textContent ?? '').toContain('Zestawienia zużycia');
  });
});

describe('raport rozliczeniowy (usage.report.build)', () => {
  beforeEach(() => {
    // Pobranie pliku sięga po API przeglądarki, którego środowisko sprawdzianu
    // nie ma w całości; podstawiamy wyłącznie te dwa punkty styku.
    globalThis.URL.createObjectURL = vi.fn(() => 'blob:raport');
    globalThis.URL.revokeObjectURL = vi.fn();
  });

  it('wytwarza plik nazwany i mówi, ile wywołań objął', async () => {
    const { element, wyslane } = zakladka({
      [Command.UsageReportBuild]: {
        content: 'wymiar,zadania\nkanal-1,40\n',
        format: TelemetryFormat.Csv,
        generatedAt: 1_700_003_600_000,
        requests: 40,
      },
    });

    const przycisk = [...element.querySelectorAll('button')].find(
      (pozycja) => pozycja.textContent === 'Zbuduj raport rozliczeniowy',
    );
    expect(przycisk, 'raport musi mieć przycisk').toBeDefined();
    przycisk?.click();
    await przemiel();

    expect(wyslane[0]?.komenda).toBe(Command.UsageReportBuild);
    expect(wyslane[0]?.zadanie).toEqual({
      fromTime: 1_700_000_000_000,
      toTime: 1_700_003_600_000,
      format: TelemetryFormat.Csv,
    });
    const tresc = element.textContent ?? '';
    expect(tresc).toContain('zuzycie-raport-1700003600000.csv');
    expect(tresc).toContain('wywołań objętych 40');
    expect(globalThis.URL.createObjectURL).toHaveBeenCalled();
  });

  it('raport bez wywołań nazywa pustkę pliku', async () => {
    const { element } = zakladka({
      [Command.UsageReportBuild]: {
        content: 'wymiar,zadania\n',
        format: TelemetryFormat.Csv,
        generatedAt: 1_700_003_600_000,
        requests: 0,
      },
    });

    const przycisk = [...element.querySelectorAll('button')].find(
      (pozycja) => pozycja.textContent === 'Zbuduj raport rozliczeniowy',
    );
    przycisk?.click();
    await przemiel();

    expect(element.textContent ?? '').toContain('plik niesie sam nagłówek');
  });
});
