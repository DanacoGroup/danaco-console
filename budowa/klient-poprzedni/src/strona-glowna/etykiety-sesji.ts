import { SessionStatus, type SessionPresence } from '../../../shared/contract';

/** Plik przekłada dane kontraktu sesji na polskie napisy karty sesji, bez elementów i bez stylu. */

/** Odmiana plakietki biblioteki komponentów niosąca stan sesji przenosi status sesji na wygląd znacznika w karcie. */
export type OdmianaStanu = 'sukces' | 'ostrzezenie' | 'neutralna';

export interface EtykietaStanu {
  tekst: string;
  odmiana: OdmianaStanu;
}

/** Komplet stanów kontraktu udostępnia etykietę i odmianę plakietki dla każdego stanu sesji zdefiniowanego w kontrakcie. */
const ETYKIETY_STANU: Readonly<Record<SessionStatus, EtykietaStanu>> = {
  active: { tekst: 'czynna', odmiana: 'sukces' },
  paused: { tekst: 'wstrzymana', odmiana: 'ostrzezenie' },
  finished: { tekst: 'zakończona', odmiana: 'neutralna' },
  archived: { tekst: 'archiwalna', odmiana: 'neutralna' },
};

export function etykietaStanuSesji(stan: SessionStatus): EtykietaStanu {
  return ETYKIETY_STANU[stan];
}

/** Opis liczby okien z żywego odpisu obecności buduje napis w postaci liczby okien i liczby okien w strumieniu. */
export function opisOkien(obecnosc: SessionPresence): string {
  const czesci = [liczebnikOkien(obecnosc.openWindowCount)];
  if (obecnosc.streamingWindowCount > 0) {
    czesci.push(`${obecnosc.streamingWindowCount} w strumieniu`);
  }
  return czesci.join(' · ');
}

/** Opis chwili ostatniej czynności podaje czas względem teraz, a dla chwili starszej niż doba oddaje pełną datę. */
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

/** Funkcja oddaje polską liczbę mnogą rzeczownika okno zgodnie z regułami odmiany liczebnika w języku polskim. */
function liczebnikOkien(liczba: number): string {
  if (liczba === 1) return '1 okno';
  const jednosc = liczba % 10;
  const setka = liczba % 100;
  const kilka = jednosc >= 2 && jednosc <= 4 && (setka < 12 || setka > 14);
  return `${liczba} ${kilka ? 'okna' : 'okien'}`;
}
