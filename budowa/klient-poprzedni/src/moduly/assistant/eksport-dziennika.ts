import type { AssistantActivityEntry } from '../../../../shared/contract';
import { chwila, dzien, NAZWY_RODZAJOW } from './etykiety-assistant';

/**
 * Dziennik działań asystenta jako log tekstowy.
 *
 * Eksport nie potrzebuje komendy kontraktu i mieć jej nie musi: wpisy są już
 * w oknie, a plik powstaje z tego, co Operator na ekranie widzi. Tę samą drogę
 * ma eksport pamięci projektu (`moduly/workspace/pamiec-pozycja.ts`)
 * i pomocnik `pobierzPlik` z `modele/kontrolki-formularza-braki.ts`.
 *
 * Log obejmuje wpisy po zawężeniu, nie cały zapis rdzenia. Plik ma odpowiadać
 * temu, co widać: eksport szerszy niż widok kazałby zgadywać, skąd wzięły się
 * wiersze, których na ekranie nie było. Nagłówek pliku nazywa więc zarówno
 * chwilę pobrania, jak i liczbę wpisów, które do niego weszły.
 *
 * Nagrania log nie niesie. Wpis ma wyłącznie odnośnik (`audioRef`), a kontrakt
 * nie ma komendy pobrania dźwięku — odnośnik idzie do pliku wprost, bo jest
 * tym, co rdzeń naprawdę oddał.
 */

/** Nazwa pliku dziennika; chwila w nazwie rozróżnia kolejne pobrania. */
export function nazwaPlikuDziennika(znacznik: number): string {
  const data = new Date(znacznik);
  const dwa = (liczba: number): string => String(liczba).padStart(2, '0');
  return [
    'dziennik-asystenta-',
    String(data.getFullYear()),
    dwa(data.getMonth() + 1),
    dwa(data.getDate()),
    '-',
    dwa(data.getHours()),
    dwa(data.getMinutes()),
    '.txt',
  ].join('');
}

/** Treść logu: nagłówek pochodzenia, a pod nim wpisy grupowane dniami. */
export function logDziennika(
  wpisy: readonly AssistantActivityEntry[],
  zakres: string,
  pobrano: number,
): string {
  const wiersze: string[] = [
    'Danaco Console — moduł Assistant, dziennik działań',
    `Pobrano: ${chwila(pobrano)}`,
    `Zakres widoku: ${zakres}`,
    `Liczba wpisów: ${String(wpisy.length)}`,
    '',
  ];

  let ostatniDzien = '';
  for (const wpis of wpisy) {
    const biezacy = dzien(wpis.createdAt);
    if (biezacy !== ostatniDzien) {
      if (ostatniDzien !== '') wiersze.push('');
      wiersze.push(`── ${biezacy} ──`);
      ostatniDzien = biezacy;
    }
    wiersze.push(`${chwila(wpis.createdAt)}  ${NAZWY_RODZAJOW[wpis.kind]}: ${wpis.content}`);
    if (wpis.actionId !== undefined && wpis.actionId !== '') {
      wiersze.push(`    zlecenie: ${wpis.actionId}`);
    }
    if (wpis.audioRef !== undefined && wpis.audioRef !== '') {
      wiersze.push(`    odnośnik nagrania: ${wpis.audioRef}`);
    }
  }
  wiersze.push('');
  return wiersze.join('\n');
}
