import type { Wynik } from './kanal';

/**
 * Przeniesienie udanego wyniku komendy na jego wycinek.
 *
 * Komenda kontraktu oddaje kopertę (`{ workflows: [...] }`, `{ files: [...] }`),
 * a widok potrzebuje samej zawartości. Pomocnik wycina pole, przepuszczając
 * odmowę nietkniętą: okna mają obowiązkowy stan błędu i muszą odróżnić „nic nie
 * ma" od „nie udało się zapytać". Wycięcie pola przez `?? []` albo `?.` w miejscu
 * wywołania zgubiłoby powód odmowy, a `permission_denied` i `not_found` są przy
 * sięganiu po cudze okna komunikacji zjawiskiem zwykłym, nie awarią.
 *
 * Odpowiedź udana, lecz bez treści, jest tu niepowodzeniem bez pola `blad`:
 * rdzeń nie podał powodu, więc pomocnik żadnego nie dopisuje. Pole jest wtedy
 * nieobecne w obiekcie, nie ustawione na `undefined` — stąd rozwinięcie warunkowe.
 */
export function przenies<Z, W>(wynik: Wynik<Z>, wybierz: (tresc: Z) => W): Wynik<W> {
  if (!wynik.udany || wynik.wynik === undefined) {
    return { udany: false, ...(wynik.blad === undefined ? {} : { blad: wynik.blad }) };
  }
  return { udany: true, wynik: wybierz(wynik.wynik) };
}
