/**
 * Dziennik działań asystenta jako log tekstowy. Plik powstaje w oknie z wpisów
 * po zawężeniu, bez komendy kontraktu, więc odpowiada temu, co widać na
 * ekranie. Nagrania log nie niesie, bo wpis ma wyłącznie odnośnik do niego.
 */
import type { AssistantActivityEntry } from '../../../../shared/contract';
import { chwila, dzien, NAZWY_RODZAJOW } from './etykiety-assistant';

/**
 * Nazwa pliku dziennika złożona z daty i godziny pobrania; chwila w nazwie
 * rozróżnia kolejne pobrania, więc plik pobrany ponownie nie zastępuje
 * poprzedniego.
 */
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

/**
 * Treść logu: nagłówek pochodzenia z chwilą pobrania, zakresem widoku i liczbą
 * wpisów, a pod nim wpisy grupowane dniami wraz z odnośnikiem zlecenia
 * i nagrania, gdy wpis je niesie.
 */
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
