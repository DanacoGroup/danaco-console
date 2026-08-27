/**
 * SKŁADNIK — POLE HASŁA Z ODSŁONIĘCIEM.
 *
 * Kontrolka odsłonięcia zmienia typ pola i własną etykietę, więc czytnik
 * ekranu wie, w którym stanie stoi. Dwa znaki leżą w przycisku, a widoczność
 * rozstrzyga arkusz po `aria-pressed` — przełączanie znaków skryptem
 * rozjeżdżałoby się ze stanem kontrolki.
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
}

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
        'aria-invalid': w.bledne === true ? 'true' : null,
        'aria-describedby': w.opisuje ?? null,
      }),
      el(
        'button',
        {
          klasa: 'au-haslo-oko',
          type: 'button',
          'aria-pressed': 'false',
          // Etykieta kontrolki jest tekstem dla użytkownika, więc stoi
          // w katalogu treści, nie w tym pliku — tak samo jak jej odpowiednik
          // po odsłonięciu hasła, który nakłada mechanika okna.
          'aria-label': tekst('dostep.haslo.pokaz'),
          dane: { odsloniecie: w.id },
        },
        [odkryte, zakryte],
      ),
    ]),
    w.opis ? el('span', { klasa: 'dn-pole-opis', id: `${w.id}-opis`, tekst: tekst(w.opis) }) : null,
  ]);
}
