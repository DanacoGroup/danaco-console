import { utworzDymekObjasnienia } from '../../komponenty/dymek';

/**
 * Dymek objaśnienia [?] przy elemencie konfiguracji okna Research.
 *
 * Nośnik zasady „każdy element konfiguracji zawiera objaśnienie kontekstowe".
 * Zachowanie w całości z biblioteki (`komponenty/dymek`); modułowi zostają
 * dwie klasy własne. Znak [?] Research jest pierścieniem 14–16 px ze
 * wskaźnikiem `help`, a nie bibliotecznym kwadratem `.dn-btn-ikona`, więc klasa
 * znaku zastępuje klasę biblioteczną, a klasa powłoki dokłada się do
 * `.dn-tooltip`.
 *
 * Trzy własności, których nie wolno ruszyć, niesie już wspólna fabryka:
 * pokazanie na `:hover` albo `:focus-within` (bez kliknięcia i bez zamykania) ·
 * znak jest przyciskiem, więc naciśnięcie prowadzi ognisko i każde naciśnięcie
 * daje odpowiedź · treść w `aria-label`, więc dymek nie potrzebuje
 * identyfikatora i nie zderza się między oknami.
 */
export function utworzDymekBadania(objasnienie: string): HTMLElement {
  return utworzDymekObjasnienia(objasnienie, { powloka: 'mr-dymek', znak: 'mr-dymek__znak' });
}

/** Pole formularza wzbogacone o dymek [?] — jedna linia zamiast trzech w oknie. */
export function zDymkiem(pole: HTMLElement, objasnienie: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mr-pole-z-dymkiem';
  element.append(pole, utworzDymekBadania(objasnienie));
  return element;
}
