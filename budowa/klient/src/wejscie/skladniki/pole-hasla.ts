/**
 * Składnik — pole hasła z odsłonięciem. Kontrolka odsłonięcia zmienia typ
 * pola i własną etykietę, więc czytnik ekranu wie, w którym stanie stoi.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciHasla {
  /** Klucz katalogu — etykieta. */
  etykieta: string;
  /** Identyfikator kontrolki. */
  id: string;
  /** Wartość podpowiedzi przeglądarki. */
  uzupelnij?: string | null;
  /** Klucz katalogu — zdanie pod polem. */
  opis?: string | null;
  /** Kontrolka niesie `aria-invalid`; obwódka bierze się z tego stanu. */
  bledne?: boolean;
  /** Identyfikator komunikatu opisującego usterkę. */
  opisuje?: string | null;
  /** Pole zakłada nowe hasło — kontrolka ogłasza wtedy wymagania menedżerom haseł. */
  nowe?: boolean;
}

/**
 * Wymagania hasła ogłoszone przeglądarce i menedżerom haseł. Bez nich menedżer
 * Google, Microsoftu czy Apple układa hasło według własnej reguły i podpowiada
 * takie, którego rejestracja nie przyjmie — Operator dowiaduje się o tym
 * dopiero po odmowie. Zapis `passwordrules` jest wspólną składnią dla tych
 * menedżerów; `minlength` i `pattern` mówią to samo samej przeglądarce.
 */
const WYMAGANIA_MENEDZERA =
  'minlength: 12; required: upper; required: lower; required: digit; allowed: special;';

/** Trzy warunki wymagane w jednym wyrażeniu: mała, wielka, cyfra, dwanaście znaków. */
const WZORZEC_HASLA = '(?=.*[a-ząćęłńóśźż])(?=.*[A-ZĄĆĘŁŃÓŚŹŻ])(?=.*\\d).{12,}';

export function poleHasla(w: WlasciwosciHasla): HTMLElement {
  const odkryte = zeZnacznika(ikony.okoOdkryte);
  odkryte.setAttribute('class', 'ik-odkryte');
  odkryte.setAttribute('aria-hidden', 'true');
  const zakryte = zeZnacznika(ikony.okoZakryte);
  zakryte.setAttribute('class', 'ik-zakryte');
  zakryte.setAttribute('aria-hidden', 'true');

  return el('div', { klasa: 'dn-pole' }, [
    el('label', { klasa: 'dn-pole-etykieta', for: w.id, tekst: tekst(w.etykieta) }),
    el('span', { klasa: 'au-haslo' }, [
      el('input', {
        klasa: 'dn-pole-kontrolka',
        type: 'password',
        id: w.id,
        autocomplete: w.uzupelnij ?? 'current-password',
        passwordrules: w.nowe === true ? WYMAGANIA_MENEDZERA : null,
        minlength: w.nowe === true ? '12' : null,
        pattern: w.nowe === true ? WZORZEC_HASLA : null,
        'aria-invalid': w.bledne === true ? 'true' : null,
        'aria-describedby': w.opisuje ?? null,
      }),
      el(
        'button',
        {
          klasa: 'au-haslo-oko',
          type: 'button',
          'aria-pressed': 'false',
          // Etykieta kontrolki jest tekstem dla użytkownika, więc stoi w katalogu treści, nie w tym pliku.
          'aria-label': tekst('dostep.haslo.pokaz'),
          dane: { odsloniecie: w.id },
        },
        [odkryte, zakryte],
      ),
    ]),
    w.opis ? el('span', { klasa: 'dn-pole-opis', id: `${w.id}-opis`, tekst: tekst(w.opis) }) : null,
  ]);
}
