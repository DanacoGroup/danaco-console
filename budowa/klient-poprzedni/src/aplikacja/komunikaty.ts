import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Waga komunikatu rozstrzyga o barwie kreski, o ikonie dymka oraz o roli
 * dostępnościowej elementu, dlatego jest jedynym pokrętłem odróżniającym
 * powiadomienia od siebie.
 */
export type WagaKomunikatu = 'info' | 'sukces' | 'ostrz' | 'blad';

/**
 * Treść komunikatu pokazywanego Operatorowi: tytuł oraz zdanie objaśniające,
 * przy czym waga pozostaje polem nieobowiązkowym, ponieważ jej brak znaczy
 * komunikat informacyjny.
 */
export interface Komunikat {
  tytul: string;
  tresc: string;
  waga?: WagaKomunikatu;
}

/**
 * Ikona przypisana wadze komunikatu; stan nigdy nie opiera się na samej barwie,
 * więc dymek niesie obok kreski znak graficzny czytelny także wtedy, gdy różnicy
 * barw czytający nie rozpoznaje.
 */
const ZNAKI: Readonly<Record<WagaKomunikatu, NazwaIkony>> = {
  info: 'info',
  sukces: 'ptaszek-kolo',
  ostrz: 'ostrzezenie',
  blad: 'blad',
};

/**
 * Wariant dymka w słowniku biblioteki komponentów
 * (`komponenty/powiadomienie.css`). Skróty wagi zostają w API pliku — pełne
 * nazwy niesie wyłącznie klasa CSS.
 */
const WARIANTY: Readonly<Record<WagaKomunikatu, string>> = {
  info: 'informacja',
  sukces: 'sukces',
  ostrz: 'ostrzezenie',
  blad: 'blad',
};

/**
 * Czas pozostawania dymka na ekranie, podany w milisekundach; po jego upływie
 * dymek znika samoczynnie, ponieważ komunikat jest doniesieniem o zdarzeniu,
 * a nie stanem wymagającym zamknięcia ręką.
 */
const CZAS_ZYCIA = 6000;

/**
 * Stos dymków powstaje przy pierwszym komunikacie, nie przy uruchomieniu
 * aplikacji, więc dokument bez komunikatu nie niesie pustego pojemnika
 * ogłaszanego technologiom wspomagającym.
 */
let stos: HTMLElement | null = null;

/**
 * Komunikaty aplikacji są jedynym miejscem, w którym aplikacja mówi Operatorowi,
 * co się właśnie stało; plik pokazuje dymek z biblioteki komponentów i nie zna
 * ani kontraktu, ani żadnego widoku.
 */
export function pokazKomunikat(komunikat: Komunikat): void {
  const waga = komunikat.waga ?? 'info';

  const element = document.createElement('div');
  element.className = `dn-toast dn-toast--${WARIANTY[waga]}`;
  element.setAttribute('role', waga === 'blad' ? 'alert' : 'status');

  const tresc = document.createElement('div');

  const tytul = document.createElement('p');
  tytul.className = 'dn-toast-tytul';
  tytul.textContent = komunikat.tytul;

  const opis = document.createElement('p');
  opis.className = 'dn-toast-tresc';
  opis.textContent = komunikat.tresc;

  // Ikonę stylizuje biblioteka selektorem `.dn-toast > svg` — bez klasy własnej.
  tresc.append(tytul, opis);
  element.append(elementIkony(ZNAKI[waga], { rozmiar: 20 }), tresc);

  gospodarzKomunikatow().append(element);
  window.setTimeout(() => element.remove(), CZAS_ZYCIA);
}

/**
 * Stos dymków przypięty do ciała dokumentu powstaje raz i zostaje na resztę
 * pracy aplikacji; kolejne komunikaty dokładają się do niego, zamiast tworzyć
 * własne pojemniki.
 */
function gospodarzKomunikatow(): HTMLElement {
  if (stos !== null) return stos;

  stos = document.createElement('div');
  stos.className = 'dn-toasty';
  stos.setAttribute('aria-live', 'polite');
  document.body.append(stos);
  return stos;
}
