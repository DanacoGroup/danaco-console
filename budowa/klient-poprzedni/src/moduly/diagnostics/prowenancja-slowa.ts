import {
  ModelCallQuality,
  ModelCallSpanKind,
  ModelCallStatus,
  TelemetryFormat,
  type ModelCallTrace,
} from '../../../../shared/contract';

// Słowa prowenancji tłumaczą wartości kontraktu na polskie nazwy pełnym zdaniem, nigdy odwrotnie.

/** Stan wywołania modelu wyrażony pełnym zdaniem po polsku, nigdy samą wartością pola kontraktu, którą czytelnik musiałby sam tłumaczyć. */
export const STAN_WYWOLANIA: Record<ModelCallStatus, string> = {
  [ModelCallStatus.Running]: 'w biegu',
  [ModelCallStatus.Ok]: 'zakończone odpowiedzią modelu',
  [ModelCallStatus.Failed]: 'nieudane — błąd dostawcy albo kanału',
  [ModelCallStatus.Timeout]: 'przerwane przekroczeniem limitu czasu',
  [ModelCallStatus.Cancelled]: 'przerwane przez Operatora',
};

/** Ocena trafności odpowiedzi modelu nadawana ręcznie przez osobę przeglądającą ślad, nie wyliczana automatycznie przez rdzeń. */
export const OCENA_WYWOLANIA: Record<ModelCallQuality, string> = {
  [ModelCallQuality.Unrated]: 'bez oceny',
  [ModelCallQuality.Accurate]: 'odpowiedź trafna',
  [ModelCallQuality.Partial]: 'odpowiedź częściowo trafna',
  [ModelCallQuality.Inaccurate]: 'odpowiedź nietrafna',
};

/** Rodzaj odcinka drzewa śladu wywołania: jeden odcinek odpowiada dokładnie jednemu krokowi całego wywołania modelu. */
export const RODZAJ_ODCINKA: Record<ModelCallSpanKind, string> = {
  [ModelCallSpanKind.Prompt]: 'złożenie i wysłanie promptu',
  [ModelCallSpanKind.Completion]: 'wytwarzanie odpowiedzi',
  [ModelCallSpanKind.Tool]: 'wywołanie narzędzia',
  [ModelCallSpanKind.Subagent]: 'wywołanie podagenta',
  [ModelCallSpanKind.Retrieval]: 'sięgnięcie po kontekst',
  [ModelCallSpanKind.Cache]: 'pamięć podręczna promptu',
};

/** Format wydania śladu: kontrakt przyjmuje cztery formaty telemetrii i okno udostępnia wszystkie, nie odejmując żadnego. */
export const FORMAT_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'OpenTelemetry Protocol w postaci JSON — format śladu z opracowania',
  [TelemetryFormat.Json]: 'JSON platformy, bez przekładu na format zewnętrzny',
  [TelemetryFormat.Jsonl]: 'JSON Lines, jeden byt w wierszu — w opracowaniu format dziennika',
  [TelemetryFormat.Csv]: 'wartości rozdzielone przecinkiem — w opracowaniu format kosztu i błędów',
};

/** Rozszerzenie pliku wydania śladu dobrane tak, aby sama nazwa pliku mówiła odbiorcy, jaką treść on niesie. */
export const ROZSZERZENIE_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'json',
  [TelemetryFormat.Json]: 'json',
  [TelemetryFormat.Jsonl]: 'jsonl',
  [TelemetryFormat.Csv]: 'csv',
};

/** Rodzaj treści pliku wydania w zapisie MIME, po którym przeglądarka rozpoznaje, jaki plik dostała do pobrania. */
export const RODZAJ_TRESCI_WYDANIA: Record<TelemetryFormat, string> = {
  [TelemetryFormat.Otlp]: 'application/json',
  [TelemetryFormat.Json]: 'application/json',
  [TelemetryFormat.Jsonl]: 'application/x-ndjson',
  [TelemetryFormat.Csv]: 'text/csv',
};

/** Znacznik czasu wywołania modelu przeliczony na zapis w strefie i formacie lokalnym przeglądającego ślad. */
export function czas(znacznik: number): string {
  return new Date(znacznik).toLocaleString('pl-PL');
}

/** Zdanie opisujące jedno wywołanie modelu, złożone wyłącznie z pól, które rdzeń rzeczywiście oddał w odpowiedzi. */
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

/** Tytuł pozycji wykazu złożony z modelu i kanału wywołania, a gdy rdzeń któregoś z nich nie podał — powiedziane wprost. */
export function tytulWywolania(wywolanie: ModelCallTrace): string {
  const model = wywolanie.model ?? '(modelu rdzeń nie podał)';
  const kanal = wywolanie.channelId ?? '(kanału rdzeń nie podał)';
  const dostawca = wywolanie.provider === undefined ? '' : `, dostawca ${wywolanie.provider}`;
  return `${model} — kanał ${kanal}${dostawca}`;
}
