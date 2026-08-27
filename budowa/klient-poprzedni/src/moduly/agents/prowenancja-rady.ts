import { MessageStatus } from '../../../../shared/contract';
import type { RadaDoradcy } from './zrodlo-doradcy';

/** Prowenancja rady doradcy, czyli zapis tego, kto zapytał, kogo i co doradca odpowiedział. */

/**
 * Jedna pozycja wykazu konsultacji: rada doradcy wraz ze zdaniem nagłówkowym
 * o tym, kto pytał i kogo, oraz ze zdaniem o drodze żądania przez rdzeń.
 * Zapis powstaje raz, a widok bierze treść i jej pochodzenie z jednego miejsca.
 */
export interface ZapisProwenancji {
  /** Rada wraz z pochodzeniem — źródło wszystkich pól poniżej. */
  rada: RadaDoradcy;
  /** Zdanie nagłówkowe: kto pytał, kogo i kiedy. */
  naglowek: string;
  /** Zdanie o drodze przez rdzeń — komendy i okno konsultacji. */
  droga: string;
}

/**
 * Etykieta, którą nosi każda rada i która nigdy nie pojawia się bez treści
 * rady. Odróżnia radę doradcy od własnej odpowiedzi eksperta wszędzie tam,
 * gdzie treść rady zostaje pokazana albo przeniesiona dalej.
 */
export const ETYKIETA_RADY = 'RADA DORADCY — nie jest odpowiedzią eksperta';

export function zapiszProwenancje(rada: RadaDoradcy): ZapisProwenancji {
  const chwila = rada.chwila > 0 ? new Date(rada.chwila).toLocaleString('pl') : 'bez czasu';
  return {
    rada,
    naglowek:
      `Ekspert „${rada.nazwaEksperta}" zapytał doradcę „${rada.nazwaDoradcy}" ` +
      `(kanał ${rada.kanalDoradcy}) — ${chwila}.`,
    droga:
      `Droga przez rdzeń: window.create (okno ${rada.idOkna}, kanał doradcy) → ` +
      'message.send → message.changed → window.close. Rada pochodzi z modelu ' +
      'doradcy, nie z modelu bazowego eksperta.',
  };
}

/**
 * Zdanie o stanie wiadomości oddanym przez rdzeń.
 *
 * Stany `stopped` i `error` wracają tą samą drogą co odpowiedź udana — z pustą
 * albo urwaną treścią. Bez tego zdania pusta rada wyglądałaby na milczenie
 * doradcy.
 */
export function zdanieOStanie(stan: MessageStatus): string {
  if (stan === MessageStatus.Complete) return '';
  if (stan === MessageStatus.Stopped) {
    return 'Doradca został zatrzymany w połowie tury — rada jest niepełna.';
  }
  if (stan === MessageStatus.Error) {
    return 'Rdzeń zakończył turę doradcy błędem — poniższa treść nie jest pełną radą.';
  }
  return `Rdzeń oddał turę doradcy w stanie „${stan}" — rada może być niepełna.`;
}

/**
 * Treść rady gotowa do przeniesienia do instrukcji eksperta.
 *
 * Przeniesiona treść niesie ze sobą nagłówek prowenancji, więc po wklejeniu do
 * instrukcji nadal widać, skąd akapit pochodzi.
 */
export function radaDoPrzeniesienia(zapis: ZapisProwenancji): string {
  return `# ${ETYKIETA_RADY}\n# ${zapis.naglowek}\n# ${zapis.droga}\n\n${zapis.rada.tresc}`;
}
