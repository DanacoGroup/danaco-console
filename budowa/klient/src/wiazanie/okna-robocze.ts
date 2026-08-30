/**
 * Okna robocze aplikacji. Okno robocze jest kontenerem jednego działania
 * Operatora: otwiera się na Centrum dowodzenia, trzyma własne karty i własny
 * widok bieżący. Okien bywa wiele naraz, a zamknięcie jednego wraca do
 * poprzedniego — praca w tamtym nie jest przerywana ani odtwarzana.
 */

/** Okno robocze widziane przez wiązanie: nazwa dla przełącznika i karta bieżąca. */
export interface OknoRobocze {
  /** Oznaczenie okna nadane przy otwarciu; nie wychodzi poza klienta. */
  id: string;
  /** Nazwa pokazywana w przełączniku; okno bez pracy nazywa się Centrum dowodzenia. */
  nazwa: string;
  /** Kod modułu karty bieżącej; pustka znaczy kartę główną — Centrum. */
  kartaBiezaca: string;
}

const OKNA: OknoRobocze[] = [];
let biezace = '';
let licznik = 0;

/** Nazwa okna roboczego bez podjętej pracy. */
const NAZWA_POCZATKOWA = 'Centrum dowodzenia';

/** Wykaz okien roboczych w kolejności otwarcia. */
export function oknaRobocze(): OknoRobocze[] {
  return OKNA;
}

/** Okno robocze bieżące; przy pierwszym pytaniu otwiera pierwsze okno. */
export function oknoBiezace(): OknoRobocze {
  if (OKNA.length === 0) otworzOkno();
  return OKNA.find((okno) => okno.id === biezace) ?? OKNA[0];
}

/** Otwiera okno robocze na Centrum dowodzenia i czyni je bieżącym. */
export function otworzOkno(): OknoRobocze {
  licznik += 1;
  const okno: OknoRobocze = { id: 'okr-' + String(licznik), nazwa: NAZWA_POCZATKOWA, kartaBiezaca: '' };
  OKNA.push(okno);
  biezace = okno.id;
  return okno;
}

/** Czyni wskazane okno bieżącym; okno nieznane zostawia stan bez zmiany. */
export function przelaczOkno(id: string): OknoRobocze | undefined {
  const okno = OKNA.find((kandydat) => kandydat.id === id);
  if (okno === undefined) return undefined;
  biezace = okno.id;
  return okno;
}

/**
 * Zamyka okno bieżące i wraca do poprzedniego. Ostatnie okno nie znika —
 * aplikacja bez okna roboczego nie miałaby czego pokazać; wraca na Centrum.
 */
export function zamknijOkno(): OknoRobocze {
  const numer = OKNA.findIndex((okno) => okno.id === biezace);
  if (OKNA.length === 1 || numer < 0) {
    const jedyne = OKNA[0] ?? otworzOkno();
    jedyne.nazwa = NAZWA_POCZATKOWA;
    jedyne.kartaBiezaca = '';
    biezace = jedyne.id;
    return jedyne;
  }
  OKNA.splice(numer, 1);
  const poprzednie = OKNA[Math.max(0, numer - 1)];
  biezace = poprzednie.id;
  return poprzednie;
}

/** Zapisuje na oknie bieżącym kartę, na której Operator stanął. */
export function zapiszKarte(kod: string, nazwa: string): void {
  const okno = oknoBiezace();
  okno.kartaBiezaca = kod;
  okno.nazwa = kod === '' ? NAZWA_POCZATKOWA : nazwa;
}
