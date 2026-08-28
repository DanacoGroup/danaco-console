/**
 * Składnik — baner komunikatu. Głowa nazywa rzecz, treść mówi, co z niej
 * wynika albo co zrobić; wygląd wnosi składnik biblioteki w wariancie ze
 * wstęgą.
 */

import { ikony, type NazwaZnaku } from '../ikony.ts';
import { el, tekst, zeZnacznika, type DanePodstawienia, type Dziecko } from '../narzedzia.ts';

/** Rodzaj banera; rozstrzyga barwę wstęgi i rolę dla czytnika ekranu, dopasowaną do wagi komunikatu wyświetlanego. */
export type RodzajBanera = 'informacja' | 'ostrzezenie' | 'blad' | 'sukces';

/** Postać licznika: zapis minut i sekund albo sama liczba sekund pozostałych do końca danego odliczania. */
export type PostacOdliczania = 'zegar' | 'sekundy';

export interface WlasciwosciBanera {
  rodzaj: RodzajBanera;
  ikona: NazwaZnaku;
  /** Klucz katalogu — głowa. */
  glowa: string;
  /** Klucz katalogu — treść; pomijalna, gdy głowa mówi wszystko. */
  tresc?: string | null;
  /** Węzły zamiast tekstu treści — gdy w zdaniu stoi wyróżniony fragment. */
  trescWezly?: Dziecko[] | null;
  dane?: DanePodstawienia;
  daneGlowy?: DanePodstawienia;
  /** Czas w sekundach — licznik w GŁOWIE. */
  odliczanie?: number | null;
  /** Czas w sekundach — licznik w TREŚCI. */
  odliczanieTresci?: number | null;
  postacOdliczania?: PostacOdliczania;
  id?: string | null;
}

/** Klasa wariantu biblioteki dla każdego rodzaju banera, przypisana zgodnie z jego wagą oraz znaczeniem komunikatu. */
const WARIANTY: Record<RodzajBanera, string> = {
  informacja: 'info',
  ostrzezenie: 'ostrzezenie',
  blad: 'blad',
  sukces: 'sukces',
};

function zLicznikiem(
  znacznik: string,
  zdanie: string,
  sekundy: number,
  postac: PostacOdliczania,
): HTMLElement {
  const czesci = zdanie.split('{odliczanie}');
  return el(znacznik, {}, [
    czesci[0] ?? '',
    el('span', {
      dane: {
        odliczanie: sekundy,
        // Wartość początkowa zostaje przy węźle, żeby licznik po dojściu do zera miał od czego ruszyć od nowa.
        'odliczanie-poczatek': sekundy,
        'odliczanie-postac': postac,
      },
    }),
    czesci[1] ?? '',
  ]);
}

export function baner(w: WlasciwosciBanera): HTMLElement {
  const znak = zeZnacznika(ikony[w.ikona]);
  znak.setAttribute('aria-hidden', 'true');
  const postac = w.postacOdliczania ?? 'zegar';
  const glowa = tekst(w.glowa, w.daneGlowy);
  const tresc = w.tresc ? tekst(w.tresc, w.dane) : '';

  // Licznik zerowy nie jest brakiem licznika: wartość z pomiaru dochodzi dopiero po odmowie rdzenia.
  const czesciTresci: Dziecko[] = [
    w.odliczanie === null || w.odliczanie === undefined
      ? el('b', { tekst: glowa })
      : zLicznikiem('b', glowa, w.odliczanie, postac),
  ];
  if (w.trescWezly) czesciTresci.push(el('span', {}, w.trescWezly));
  else if (w.odliczanieTresci !== null && w.odliczanieTresci !== undefined) {
    czesciTresci.push(zLicznikiem('span', tresc, w.odliczanieTresci, postac));
  } else if (w.tresc) czesciTresci.push(tresc);

  return el(
    'div',
    {
      klasa: `dn-alert dn-alert--wstega dn-alert--${WARIANTY[w.rodzaj]}`,
      id: w.id ?? null,
      // Baner z usterką przerywa czytnikowi ekranu; informacyjny czeka na przerwę — rola idzie za rodzajem.
      role: w.rodzaj === 'blad' || w.rodzaj === 'ostrzezenie' ? 'alert' : 'status',
    },
    [
      el('span', { klasa: 'dn-alert-znak', 'aria-hidden': 'true' }, [znak]),
      el('span', { klasa: 'dn-alert-tresc' }, czesciTresci),
    ],
  );
}
