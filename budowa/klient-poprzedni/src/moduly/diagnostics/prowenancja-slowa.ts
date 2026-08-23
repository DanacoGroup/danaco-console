import {
  ModelCallQuality,
  ModelCallSpanKind,
  ModelCallStatus,
  TelemetryFormat,
  type ModelCallTrace,
} from '../../../../shared/contract';

/**
 * Słowa prowenancji widoczne dla Operatora.
 *
 * Wartości kontraktu są angielskie i techniczne (`ok`, `otlp`, `unrated`).
 * W oknie stoi pełna nazwa polska, bo Operator czyta okno, nie kontrakt —
 * a zakaz numeracji i kodów w produkcie znaczy również zakaz pokazywania
 * wartości pola zamiast jej nazwy. Przekład jest jednostronny: do rdzenia jedzie
 * wyłącznie wartość kontraktu, nigdy napis z tego pliku.
 *
 * Wykazy są pełne wobec kontraktu i kompilator tego pilnuje: `Record` po typie
 * wartości nie skompiluje się, gdy kontrakt dołoży stan albo format, więc nowa
 * wartość nie przemknie do okna jako pusty napis.
 */

/** Stan wywołania — pełnym zdaniem, nie wartością pola. */
export const STAN_WYWOLANIA: Record<ModelCallStatus, string> = {
  [ModelCallStatus.Running]: 'w biegu',
  [ModelCallStatus.Ok]: 'zakończone odpowiedzią modelu',
  [ModelCallStatus.Failed]: 'nieudane — błąd dostawcy albo kanału',
  [ModelCallStatus.Timeout]: 'przerwane przekroczeniem limitu czasu',
  [ModelCallStatus.Cancelled]: 'przerwane przez Operatora',
};

/** Ocena trafności nadawana ręcznie przez Operatora. */
export const OCENA_WYWOLANIA: Record<ModelCallQuality, string> = {
  [ModelCallQuality.Unrated]: 'bez oceny',
  [ModelCallQuality.Accurate]: 'odpowiedź trafna',
  [ModelCallQuality.Partial]: 'odpowiedź częściowo trafna',
  [ModelCallQuality.Inaccurate]: 'odpowiedź nietrafna',
};

/** Rodzaj odcinka drzewa śladu: jeden odcinek to jeden krok wywołania. */
export const RODZAJ_ODCINKA: Record<ModelCallSpanKind, string> = {
  [ModelCallSpanKind.Prompt]: 'złożenie i wysłanie promptu',
  [ModelCallSpanKind.Completion]: 'wytwarzanie odpowiedzi',
  [ModelCallSpanKind.Tool]: 'wywołanie narzędzia',
  [ModelCallSpanKind.Subagent]: 'wywołanie podagenta',
  [ModelCallSpanKind.Retrieval]: 'sięgnięcie po kontekst',
  [ModelCallSpanKind.Cache]: 'pamięć podręczna promptu',
};

/**
 * Format wydania śladu wraz z tym, czym jest.
 *
 * Opracowanie modułu przypisuje śladom OpenTelemetry Protocol oraz JSON,
 * a JSON Lines i wartości rozdzielone przecinkiem wylicza przy dziennikach,
 * błędach i koszcie. Kontrakt przyjmuje dla śladu wszystkie cztery, więc okno
 * daje wszystkie cztery i mówi, który z nich opracowanie śladom przypisuje —
 * odjęcie formatu, który rdzeń przyjmuje, byłoby brakiem funkcji zrobionym
 * w oknie.
 */
export const FORMAT_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'OpenTelemetry Protocol w postaci JSON — format śladu z opracowania',
  [TelemetryFormat.Json]: 'JSON platformy, bez przekładu na format zewnętrzny',
  [TelemetryFormat.Jsonl]: 'JSON Lines, jeden byt w wierszu — w opracowaniu format dziennika',
  [TelemetryFormat.Csv]: 'wartości rozdzielone przecinkiem — w opracowaniu format kosztu i błędów',
};

/** Rozszerzenie pliku wydania; nazwa pliku ma mówić, co w nim jest. */
export const ROZSZERZENIE_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'json',
  [TelemetryFormat.Json]: 'json',
  [TelemetryFormat.Jsonl]: 'jsonl',
  [TelemetryFormat.Csv]: 'csv',
};

/** Rodzaj treści pliku wydania — po nim przeglądarka wie, co dostała. */
export const RODZAJ_TRESCI_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'application/json',
  [TelemetryFormat.Json]: 'application/json',
  [TelemetryFormat.Jsonl]: 'application/x-ndjson',
  [TelemetryFormat.Csv]: 'text/csv',
};

/** Znacznik czasu w zapisie lokalnym Operatora. */
export function czas(znacznik: number): string {
  return new Date(znacznik).toLocaleString('pl-PL');
}

/**
 * Zdanie o jednym wywołaniu złożone WYŁĄCZNIE z pól, które rdzeń oddał.
 *
 * Pole nieoddane nie staje się zerem: „koszt 0" i „kanał bez cennika" to dwa
 * różne zdania o instalacji, a wywołanie w biegu nie ma jeszcze opóźnienia ani
 * liczby tokenów i nie jest to usterka.
 */
export function opisWywolania(wywolanie: ModelCallTrace): string {
  const czesci: string[] = [
    STAN_WYWOLANIA[wywolanie.status],
    `początek ${czas(wywolanie.startedAt)}`,
  ];
  if (wywolanie.latencyMs !== undefined) {
    czesci.push(`opóźnienie ${String(wywolanie.latencyMs)} ms`);
  }
  if (wywolanie.slow === true) czesci.push('POWYŻEJ progu opóźnienia z ustawienia rdzenia');
  if (wywolanie.totalTokens !== undefined) czesci.push(`tokenów ${String(wywolanie.totalTokens)}`);
  if (wywolanie.cachedTokens !== undefined) {
    czesci.push(`z pamięci podręcznej promptu ${String(wywolanie.cachedTokens)}`);
  }
  czesci.push(
    wywolanie.cost === undefined
      ? 'kosztu rdzeń nie podał — kanał bez cennika, a to nie znaczy zero'
      : `koszt ${String(wywolanie.cost)} ${wywolanie.currency ?? '(waluty rdzeń nie podał)'}`,
  );
  czesci.push(OCENA_WYWOLANIA[wywolanie.quality ?? ModelCallQuality.Unrated]);
  czesci.push(
    wywolanie.contentStored
      ? 'treść promptu i odpowiedzi jest zapisana'
      : 'treści promptu i odpowiedzi rdzeń NIE zapisał — zapis treści był wyłączony ustawieniem w chwili wywołania',
  );
  if (wywolanie.redacted === true) czesci.push('treść poddana redakcji danych wrażliwych');
  return czesci.join(' · ');
}

/** Tytuł pozycji wykazu: model i kanał, a gdy ich nie ma — powiedziane wprost. */
export function tytulWywolania(wywolanie: ModelCallTrace): string {
  const model = wywolanie.model ?? '(modelu rdzeń nie podał)';
  const kanal = wywolanie.channelId ?? '(kanału rdzeń nie podał)';
  const dostawca = wywolanie.provider === undefined ? '' : `, dostawca ${wywolanie.provider}`;
  return `${model} — kanał ${kanal}${dostawca}`;
}
