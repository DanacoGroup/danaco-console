import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/** Waga komunikatu — rozstrzyga o kresce i ikonie dymka. */
export type WagaKomunikatu = 'info' | 'sukces' | 'ostrz' | 'blad';

/** Treść komunikatu pokazywanego Operatorowi. */
export interface Komunikat {
  tytul: string;
  tresc: string;
  waga?: WagaKomunikatu;
}

/** Ikona wagi; stan nigdy nie opiera się na samej barwie. */
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

/** Ile czasu dymek zostaje na ekranie, w milisekundach. */
const CZAS_ZYCIA = 6000;

/** Stos dymków; powstaje przy pierwszym komunikacie, nie przy uruchomieniu. */
let stos: HTMLElement | null = null;

/**
 * Komunikaty aplikacji — jedno miejsce, w którym mówi się Operatorowi,
 * co się właśnie stało.
 *
 * Jedna odpowiedzialność: pokazanie dymka z biblioteki `komponenty/`.
 * Plik nie zna ani kontraktu, ani żadnego widoku.
 *
 * Żaden przycisk w aplikacji nie jest wyszarzany, więc każde naciśnięcie musi
 * dać odpowiedź. Gdy zamiar nie ma odpowiednika w komendzie kontraktu albo
 * rdzeń odmawia, dymek mówi to wprost — nie udaje wykonanej czynności.
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

/** Stos dymków przypięty do ciała dokumentu; powstaje raz. */
function gospodarzKomunikatow(): HTMLElement {
  if (stos !== null) return stos;

  stos = document.createElement('div');
  stos.className = 'dn-toasty';
  stos.setAttribute('aria-live', 'polite');
  document.body.append(stos);
  return stos;
}
