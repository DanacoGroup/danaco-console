import { LICZBA_MAX } from '../../okna-rownolegle/identyfikatory';

/** Sufit uczestnika debaty — opis stanu, nie zakaz: liczba gniazd sceny okien równoległych, sufit sceny, nie składu debaty. */
export const GNIAZD_SCENY = LICZBA_MAX;

/** Stan sufitu widziany z liczby uczestników, których okno dziś zna, wraz z gotowymi zdaniami opisowymi. */
export interface StanSufitu {
  /** Uczestnicy znani oknu w tej chwili. */
  uczestnicy: number;
  /** Ilu uczestników nie dostałoby własnego gniazda sceny; zero, gdy starcza. */
  ponadGniazda: number;
  /** Zdania opisujące stan — w kolejności czytania. */
  zdania: string[];
}

/**
 * Składa opis stanu. Żadne ze zdań nie orzeka zakazu: mówią, ilu uczestników
 * jest, czego nie liczy rdzeń ani kontrakt, gdzie leży sufit sceny i czego
 * jeszcze nie ustalono.
 */
export function stanSufitu(iluUczestnikow: number): StanSufitu {
  const uczestnicy = Number.isFinite(iluUczestnikow) && iluUczestnikow > 0 ? Math.trunc(iluUczestnikow) : 0;
  const ponadGniazda = uczestnicy > GNIAZD_SCENY ? uczestnicy - GNIAZD_SCENY : 0;

  const zdania = [
    uczestnicy === 0
      ? 'Debata nie ma dziś w tym oknie ani jednego uczestnika.'
      : `Debata ma dziś w tym oknie ${uczestnicy} ${odmianaUczestnikow(uczestnicy)}.`,
    'Liczby uczestników nie ogranicza ani rdzeń, ani kontrakt: tabela debata_uczestnik nie ma więzu na liczbę wierszy, roundtable.model.add nie ma pola z granicą, a jedyną granicą obszaru jest granica_tur — granica liczby tur, nie składu.',
    `Sufit gniazd należy do czego innego: scena okien równoległych ma ${GNIAZD_SCENY} gniazda (okna-rownolegle/identyfikatory.ts → LICZBA_MAX, gniazda okno-1…okno-${GNIAZD_SCENY}). Ten moduł mieści cały skład w jednym oknie, w panelu na uczestnika, więc gniazd sceny nie zużywa.`,
  ];

  if (ponadGniazda > 0) {
    zdania.push(
      `Skład jest o ${ponadGniazda} większy niż liczba gniazd sceny — rdzeń go przyjął i okno go rysuje. Własnego gniazda na scenie okien równoległych nie starczyłoby dla ${ponadGniazda} z nich, gdyby każdy uczestnik miał stać w osobnym oknie.`,
    );
  }

  zdania.push(
    'Czy debata ma prawo wyjść ponad liczbę gniazd sceny, nie jest ustalone. Do czasu ustalenia okno nie zgaduje w żadną stronę.',
  );

  return { uczestnicy, ponadGniazda, zdania };
}

/** Nota o suficie w miejscu, w którym Operator dokłada uczestników — akapit opisowy, nie ostrzeżenie i nie bramka. */
export function utworzNoteSufitu(iluUczestnikow: number): HTMLElement {
  const stan = stanSufitu(iluUczestnikow);

  const nota = document.createElement('div');
  nota.className = 'dr-sufit';
  nota.dataset['ponadGniazda'] = stan.ponadGniazda > 0 ? 'tak' : 'nie';

  const tytul = document.createElement('span');
  tytul.className = 'dr-sufit__tytul';
  tytul.textContent = 'Ilu uczestników mieści debata';
  nota.append(tytul);

  for (const zdanie of stan.zdania) {
    const akapit = document.createElement('p');
    akapit.className = 'dr-sufit__zdanie';
    akapit.textContent = zdanie;
    nota.append(akapit);
  }
  return nota;
}

/**
 * Odmiana rzeczownika przy liczbie — polszczyzna, nie „1 uczestników".
 * Rzeczownik jest męskoosobowy, więc poza jedynką zostaje dopełniacz mnogi
 * („ma 2 uczestników", „ma 7 uczestników") — jedna forma, bez wyjątków.
 */
function odmianaUczestnikow(ile: number): string {
  return ile === 1 ? 'uczestnika' : 'uczestników';
}
