import { SessionStatus, type SessionPresence } from '../../../shared/contract';

/**
 * Etykiety sekcji sesji w tle — przekład danych kontraktu na polskie napisy.
 *
 * Zakres pliku to treść napisów karty sesji — bez elementów i bez stylu.
 *
 * Żadna funkcja nie wymyśla danych: każda przekłada wartość oddaną przez
 * rdzeń, a wartość nieobecną oddaje jako `undefined` — karta pomija wtedy
 * cały fragment zamiast pokazać zmyślony.
 *
 * Nazw środowisk tu nie ma: daje je `wykaz-srodowisk.ts`, zasilany wykazem
 * z rdzenia, żeby na jednym ekranie nie stały dwa źródła tej samej nazwy.
 */

/** Odmiana plakietki biblioteki komponentów niosąca stan sesji. */
export type OdmianaStanu = 'sukces' | 'ostrzezenie' | 'neutralna';

export interface EtykietaStanu {
  tekst: string;
  odmiana: OdmianaStanu;
}

/** Komplet stanów kontraktu — nowy stan w `contract.json` przerywa kompilację. */
const ETYKIETY_STANU: Readonly<Record<SessionStatus, EtykietaStanu>> = {
  active: { tekst: 'czynna', odmiana: 'sukces' },
  paused: { tekst: 'wstrzymana', odmiana: 'ostrzezenie' },
  finished: { tekst: 'zakończona', odmiana: 'neutralna' },
  archived: { tekst: 'archiwalna', odmiana: 'neutralna' },
};

export function etykietaStanuSesji(stan: SessionStatus): EtykietaStanu {
  return ETYKIETY_STANU[stan];
}

/** Liczby okien z żywego odpisu: „2 okna · 1 w strumieniu". */
export function opisOkien(obecnosc: SessionPresence): string {
  const czesci = [liczebnikOkien(obecnosc.openWindowCount)];
  if (obecnosc.streamingWindowCount > 0) {
    czesci.push(`${obecnosc.streamingWindowCount} w strumieniu`);
  }
  return czesci.join(' · ');
}

/** Chwila ostatniej czynności względem teraz; starsza niż doba — datą. */
export function opisCzasu(chwila: number, teraz: number = Date.now()): string {
  const minuty = Math.floor(Math.max(0, teraz - chwila) / 60_000);
  if (minuty < 1) return 'przed chwilą';
  if (minuty < 60) return `${minuty} min temu`;
  const godziny = Math.floor(minuty / 60);
  if (godziny < 24) return `${godziny} godz. temu`;
  return new Intl.DateTimeFormat('pl', { dateStyle: 'medium', timeStyle: 'short' }).format(
    new Date(chwila),
  );
}

/** Polska liczba mnoga: 1 okno, 2–4 okna, 5+ okien (z wyjątkiem 12–14). */
function liczebnikOkien(liczba: number): string {
  if (liczba === 1) return '1 okno';
  const jednosc = liczba % 10;
  const setka = liczba % 100;
  const kilka = jednosc >= 2 && jednosc <= 4 && (setka < 12 || setka > 14);
  return `${liczba} ${kilka ? 'okna' : 'okien'}`;
}
