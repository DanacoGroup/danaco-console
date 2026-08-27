/**
 * SKŁADNIK — POLE TEKSTOWE.
 *
 * Etykieta stoi nad kontrolką, opis pod nią. Etykieta jest elementem `label`
 * związanym identyfikatorem, nie samym tekstem obok — bez tego wskazanie
 * etykiety nie ustawia kursora w polu, a czytnik ekranu nie wie, co czyta.
 */

import { el, tekst } from '../narzedzia.ts';

export interface WlasciwosciPola {
  /** Klucz katalogu — etykieta. */
  etykieta: string;
  /** Identyfikator kontrolki; wiąże ją z etykietą i z komunikatem usterki. */
  id: string;
  /** Typ pola; domyślnie zwykły tekst. */
  typ?: string | null;
  /** Wartość podpowiedzi przeglądarki. */
  uzupelnij?: string | null;
  /** Klucz katalogu — zdanie pod polem. */
  opis?: string | null;
  /**
   * Kontrolka niesie `aria-invalid`; obwódka bierze się z tego stanu, nie
   * z osobnej klasy — inaczej wygląd i ogłoszenie czytnika rozjeżdżają się.
   */
  bledne?: boolean;
  /** Identyfikator komunikatu opisującego usterkę. */
  opisuje?: string | null;
}

export function poleTekstowe(w: WlasciwosciPola): HTMLElement {
  return el('label', { klasa: 'dn-pole' }, [
    el('span', { klasa: 'dn-pole-etykieta', tekst: tekst(w.etykieta) }),
    el('input', {
      klasa: 'dn-pole-kontrolka',
      type: w.typ ?? 'text',
      id: w.id,
      autocomplete: w.uzupelnij ?? null,
      'aria-invalid': w.bledne === true ? 'true' : null,
      'aria-describedby': w.opisuje ?? null,
    }),
    w.opis ? el('span', { klasa: 'dn-pole-opis', tekst: tekst(w.opis) }) : null,
  ]);
}
