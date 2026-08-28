import { beforeEach, describe, expect, it, vi } from 'vitest';

import {
  Command,
  ModelCallQuality,
  ModelCallSpanKind,
  ModelCallStatus,
  TelemetryFormat,
  type ModelCallTrace,
} from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzZrodloProwenancji } from './prowenancja-zrodlo';
import { utworzStanDiagnostyki } from './stan-diagnostyki';
import { utworzZakladkeProwenancji } from './zakladka-prowenancji';
import { utworzZrodloDiagnostics } from './zrodlo-diagnostics';

// Sprawdzian ujawnia pracę modeli w oknie: wykaz wywołań, odczyt śladu, ocenę i wydanie do pliku.

/** Funkcja tworzy jedno wywołanie modelu w postaci, w której rdzeń je oddaje, gotowe do podmiany wybranych pól w pojedynczym sprawdzianie. */
function wywolanie(zmiany: Partial<ModelCallTrace> = {}): ModelCallTrace {
  return {
    id: 'wywolanie-1',
    status: ModelCallStatus.Ok,
    startedAt: 1_700_000_000_000,
    contentStored: true,
    model: 'model-przykladowy',
    channelId: 'kanal-1',
    latencyMs: 1234,
    totalTokens: 900,
    ...zmiany,
  };
}

/** Funkcja tworzy kanał próbny, który zapamiętuje wysłane żądania i oddaje odpowiedź wskazaną dla każdej komendy z osobna. */
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
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane };
}

/** Funkcja składa zakładkę prowenancji wraz ze stanem diagnostyki i kanałem próbnym, gotową do sprawdzenia jej zachowania. */
function zakladka(odpowiedzi: Record<string, unknown>, zeZrodlem = true) {
  const { kanal, wyslane } = kanalProbny(odpowiedzi);
  const diagnostyka = utworzZrodloDiagnostics(kanal);
  const stan = utworzStanDiagnostyki(diagnostyka);
  const widok = utworzZakladkeProwenancji(
    diagnostyka,
    stan,
    zeZrodlem ? utworzZrodloProwenancji(kanal) : undefined,
  );
  return { widok, wyslane, stan };
}

/** Funkcja znajduje w oknie przycisk o dokładnie podanej etykiecie; brak takiego przycisku sprawdzian zgłasza jako brak drogi. */
function przycisk(gdzie: HTMLElement, etykieta: string): HTMLButtonElement {
  const znaleziony = [...gdzie.querySelectorAll('button')].find(
    (kontrolka) => kontrolka.textContent === etykieta,
  );
  expect(znaleziony, `przycisk „${etykieta}" musi być w oknie`).toBeDefined();
  return znaleziony as HTMLButtonElement;
}

/** Funkcja zwraca treść zakładki jako jeden napis tekstowy, przydatny do sprawdzania obecności zdań, a nie układu ekranu. */
function napis(gdzie: HTMLElement): string {
  return gdzie.textContent ?? '';
}

/** Funkcja oddaje sterowanie kolejce mikrozadań trzykrotnie, aby łańcuch obietnic wywołań zdążył się rozstrzygnąć przed sprawdzeniem wyniku. */
async function przemiel(): Promise<void> {
  await Promise.resolve();
  await Promise.resolve();
  await Promise.resolve();
}

beforeEach(() => {
  // Pobranie pliku sięga po URL obiektowy, którego jsdom nie ma; wydanie śladu musi być sprawdzalne.
  URL.createObjectURL = vi.fn(() => 'blob:slad');
  URL.revokeObjectURL = vi.fn();
});

describe('Provenance Explorer — wykaz wywołań', () => {
  it('woła provenance.call.list i pokazuje wywołanie wraz z kosztem nieoddanym', async () => {
    const { widok, wyslane } = zakladka({
      [Command.ProvenanceCallList]: { calls: [wywolanie()], total: 7, truncated: true },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([Command.ProvenanceCallList]);
    expect(napis(widok.element)).toContain('model-przykladowy');
    // Koszt nieoddany nie staje się zerem.
    expect(napis(widok.element)).toContain('kosztu rdzeń nie podał');
    // Przycięcie wyniku mówi się wprost.
    expect(napis(widok.element)).toContain('WYNIK PRZYCIĘTY');
    expect(napis(widok.element)).toContain('z 7 spełniających warunki');
  });

  it('pustka rejestru ma zdanie „nie było jeszcze wywołań”, nie puste miejsce', async () => {
    const { widok } = zakladka({ [Command.ProvenanceCallList]: { calls: [] } });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    expect(napis(widok.element)).toContain('nie było jeszcze wywołań');
  });

  it('pustka przy założonym zawężeniu mówi o zawężeniu, nie o braku wywołań', async () => {
    const { widok } = zakladka({ [Command.ProvenanceCallList]: { calls: [] } });
    const wzorzec = widok.element.querySelector<HTMLInputElement>(
      'input[aria-label="Wzorzec w treści promptu i odpowiedzi"]',
    );
    expect(wzorzec).not.toBeNull();
    (wzorzec as HTMLInputElement).value = 'faktura';
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    expect(napis(widok.element)).toContain('wzorzec „faktura"');
    expect(napis(widok.element)).not.toContain('nie było jeszcze wywołań');
  });

  it('odmowę rdzenia nazywa odmową, a nie pustką', async () => {
    const { widok } = zakladka({});
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    expect(napis(widok.element)).toContain('Rdzeń odmówił odczytu rejestru wywołań modeli');
    expect(napis(widok.element)).toContain('not_found');
  });

  it('bez źródła prowenancji z montażu mówi, czego brakuje i po czyjej stronie', async () => {
    const { widok, wyslane } = zakladka({}, false);
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    expect(wyslane).toHaveLength(0);
    expect(napis(widok.element)).toContain('PO STRONIE KLIENTA');
    expect(napis(widok.element)).toContain('indeks.ts');
  });
});

describe('Provenance Explorer — odczyt śladu jednego wywołania', () => {
  it('woła provenance.call.get i pokazuje drzewo odcinków wraz z treścią', async () => {
    const { widok, wyslane } = zakladka({
      [Command.ProvenanceCallList]: { calls: [wywolanie()] },
      [Command.ProvenanceCallGet]: {
        call: wywolanie(),
        spans: [
          {
            id: 'odcinek-1',
            callId: 'wywolanie-1',
            name: 'prompt',
            kind: ModelCallSpanKind.Prompt,
            startedAt: 1_700_000_000_000,
            status: ModelCallStatus.Ok,
          },
          {
            id: 'odcinek-2',
            callId: 'wywolanie-1',
            parentSpanId: 'odcinek-1',
            name: 'szukaj-w-repozytorium',
            kind: ModelCallSpanKind.Tool,
            startedAt: 1_700_000_000_500,
            status: ModelCallStatus.Ok,
            durationMs: 40,
          },
        ],
        prompt: 'pytanie do modelu',
        response: 'odpowiedź modelu',
        redacted: false,
      },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();
    przycisk(widok.element, 'Odczytaj ślad wywołania').click();
    await przemiel();

    expect(wyslane.map((wpis) => wpis.komenda)).toContain(Command.ProvenanceCallGet);
    const drzewo = widok.element.querySelector('[data-drzewo-sladu]');
    expect(drzewo).not.toBeNull();
    // Odcinek podrzędny ma wcięcie — kto kogo wywołał, zostaje widoczne.
    expect(drzewo?.textContent).toContain('  szukaj-w-repozytorium');
    expect(napis(widok.element)).toContain('pytanie do modelu');
    expect(napis(widok.element)).toContain('bez redakcji danych wrażliwych');
  });

  it('wywołanie bez zapisanej treści mówi, że zapis był wyłączony ustawieniem', async () => {
    const bezTresci = wywolanie({ contentStored: false });
    const { widok } = zakladka({
      [Command.ProvenanceCallList]: { calls: [bezTresci] },
      [Command.ProvenanceCallGet]: { call: bezTresci, spans: [], redacted: false },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();
    przycisk(widok.element, 'Odczytaj ślad wywołania').click();
    await przemiel();

    expect(napis(widok.element)).toContain('zapis treści promptu i odpowiedzi był w chwili wywołania wyłączony');
    expect(napis(widok.element)).toContain('bez ani jednego odcinka');
  });
});

describe('Provenance Explorer — ocena wywołania przez Operatora', () => {
  it('ocena ma jawny chwyt i wysyła provenance.call.rate z oceną i uzasadnieniem', async () => {
    const { widok, wyslane } = zakladka({
      [Command.ProvenanceCallList]: { calls: [wywolanie()] },
      [Command.ProvenanceCallRate]: {
        call: wywolanie({ quality: ModelCallQuality.Partial, qualityNote: 'pominęła trzecie źródło' }),
      },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    // Chwyt jest widocznym przyciskiem, nie skrótem i nie kliknięciem wiersza.
    przycisk(widok.element, 'Oceń odpowiedź modelu').click();
    const trafnosc = widok.element.querySelector<HTMLSelectElement>(
      'select[aria-label="Trafność odpowiedzi modelu"]',
    );
    const uzasadnienie = widok.element.querySelector<HTMLTextAreaElement>(
      'textarea[aria-label="Uzasadnienie oceny wywołania"]',
    );
    expect(trafnosc).not.toBeNull();
    (trafnosc as HTMLSelectElement).value = ModelCallQuality.Partial;
    (uzasadnienie as HTMLTextAreaElement).value = 'pominęła trzecie źródło';

    przycisk(widok.element, 'Zapisz ocenę wywołania').click();
    await przemiel();

    const zapis = wyslane.find((wpis) => wpis.komenda === Command.ProvenanceCallRate);
    expect(zapis?.zadanie).toEqual({
      callId: 'wywolanie-1',
      quality: ModelCallQuality.Partial,
      note: 'pominęła trzecie źródło',
    });
    expect(napis(widok.element)).toContain('Rdzeń zapisał ocenę wywołania');
  });

  it('ocenę podmienioną przez rdzeń nazywa podmianą, a nie powodzeniem', async () => {
    const { widok } = zakladka({
      [Command.ProvenanceCallList]: { calls: [wywolanie()] },
      [Command.ProvenanceCallRate]: { call: wywolanie({ quality: ModelCallQuality.Unrated }) },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();
    przycisk(widok.element, 'Oceń odpowiedź modelu').click();
    const trafnosc = widok.element.querySelector<HTMLSelectElement>(
      'select[aria-label="Trafność odpowiedzi modelu"]',
    );
    (trafnosc as HTMLSelectElement).value = ModelCallQuality.Accurate;
    przycisk(widok.element, 'Zapisz ocenę wywołania').click();
    await przemiel();

    expect(napis(widok.element)).toContain('zapisał ocenę INNĄ niż wysłana');
  });
});

describe('Provenance Explorer — wydanie śladu', () => {
  it('mówi, co plik poniesie, zanim plik powstanie', () => {
    const { widok } = zakladka({});
    const zapowiedz = widok.element.querySelector('[data-zapowiedz-wydania]')?.textContent ?? '';
    expect(zapowiedz).toContain('OpenTelemetry Protocol');
    // Treść promptu nie wychodzi z instalacji przez samo naciśnięcie przycisku.
    expect(zapowiedz).toContain('bez treści promptów i odpowiedzi');
  });

  it('wydaje ślad wskazanych wywołań i mówi, czy treść jest po redakcji', async () => {
    const { widok, wyslane } = zakladka({
      [Command.ProvenanceCallList]: { calls: [wywolanie()] },
      [Command.ProvenanceTraceExport]: {
        content: '{"slad":true}',
        format: TelemetryFormat.Otlp,
        callCount: 1,
        redacted: false,
      },
    });
    przycisk(widok.element, 'Odczytaj rejestr wywołań').click();
    await przemiel();

    const wskazanie = widok.element.querySelector<HTMLInputElement>('input[data-do-wydania]');
    expect(wskazanie).not.toBeNull();
    (wskazanie as HTMLInputElement).checked = true;
    const zTrescia = widok.element.querySelector<HTMLInputElement>(
      'input[aria-label="Wydaj wraz z treścią promptu i odpowiedzi"]',
    );
    (zTrescia as HTMLInputElement).checked = true;

    przycisk(widok.element, 'Wydaj ślad wywołań do pliku').click();
    await przemiel();

    const wydanie = wyslane.find((wpis) => wpis.komenda === Command.ProvenanceTraceExport);
    expect(wydanie?.zadanie).toEqual({
      format: TelemetryFormat.Otlp,
      callIds: ['wywolanie-1'],
      includeContent: true,
    });
    expect(URL.createObjectURL).toHaveBeenCalledTimes(1);
    expect(napis(widok.element)).toContain('BEZ redakcji danych wrażliwych');
  });

  it('wydanie zerowe nie zapisuje pliku i mówi o tym wprost', async () => {
    const { widok } = zakladka({
      [Command.ProvenanceTraceExport]: {
        content: '',
        format: TelemetryFormat.Csv,
        callCount: 0,
      },
    });
    przycisk(widok.element, 'Wydaj ślad wywołań do pliku').click();
    await przemiel();

    expect(URL.createObjectURL).not.toHaveBeenCalled();
    expect(napis(widok.element)).toContain('ZERO wywołań');
  });
});

describe('Provenance Explorer — droga do wszystkich czterech komend', () => {
  it('źródło prowenancji wysyła każdą z czterech komend rodziny', async () => {
    const { kanal, wyslane } = kanalProbny({});
    const zrodlo = utworzZrodloProwenancji(kanal);
    await zrodlo.wykazWywolan({});
    await zrodlo.odczytajWywolanie({ callId: 'wywolanie-1' });
    await zrodlo.ocenWywolanie({ callId: 'wywolanie-1', quality: ModelCallQuality.Accurate });
    await zrodlo.wydajSlad({ format: TelemetryFormat.Otlp });

    expect(wyslane.map((wpis) => wpis.komenda)).toEqual([
      Command.ProvenanceCallList,
      Command.ProvenanceCallGet,
      Command.ProvenanceCallRate,
      Command.ProvenanceTraceExport,
    ]);
  });
});
