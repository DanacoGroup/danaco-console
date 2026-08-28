import type { PanelSection } from '../../../shared/contract';
import { przyciskAkcji } from '../modele/kontrolki-formularza';
import type { OpisSekcji } from './uklad-sekcji';

// Kształt jednej sekcji panelu — sam widok, bez reguły układu, oddający czynności wywołaniem zwrotnym.

/** Cztery czynności, które operator wykonuje na sekcji panelu bez wygaszania przycisku ani pytania o potwierdzenie. */
export interface CzynnosciSekcji {
  przyZwinieciu(id: string): void;
  przyPrzesunieciu(id: string, kierunek: -1 | 1): void;
  przyZdjeciu(id: string): void;
  przyPrzywroceniu(id: string): void;
}

/** Przycisk czynności sekcji — jeden kształt dla wszystkich czterech czynności dostępnych na panelu bocznym. */
function przyciskSekcji(etykieta: string, opis: string): HTMLButtonElement {
  const kontrolka = przyciskAkcji(etykieta, 'dn-btn dn-btn--sm dn-btn--zarys');
  kontrolka.title = opis;
  kontrolka.setAttribute('aria-label', opis);
  return kontrolka;
}

/**
 * Jedna sekcja panelu wraz z nagłówkiem czynności.
 *
 * `stan` pochodzi z układu obowiązującego (czyli z odpowiedzi rdzenia), a nie
 * z zamiaru widoku — dlatego zwinięcie widać dopiero wtedy, gdy rdzeń je
 * przyjął.
 */
export function zlozSekcje(
  opis: OpisSekcji,
  stan: PanelSection,
  przedrostek: string,
  czynnosci: CzynnosciSekcji,
): HTMLElement {
  const tytul = document.createElement('h4');
  tytul.className = `dn-karta-tytul ${przedrostek}-sekcja__tytul`;
  tytul.textContent = opis.tytul;

  const zwin = przyciskSekcji(
    stan.collapsed ? 'Rozwiń' : 'Zwiń',
    `${stan.collapsed ? 'Rozwiń' : 'Zwiń'} sekcję ${opis.tytul}`,
  );
  zwin.setAttribute('aria-expanded', String(!stan.collapsed));
  zwin.addEventListener('click', () => czynnosci.przyZwinieciu(opis.id));

  const wGore = przyciskSekcji('▲', `Przesuń sekcję ${opis.tytul} wyżej w panelu`);
  wGore.addEventListener('click', () => czynnosci.przyPrzesunieciu(opis.id, -1));

  const wDol = przyciskSekcji('▼', `Przesuń sekcję ${opis.tytul} niżej w panelu`);
  wDol.addEventListener('click', () => czynnosci.przyPrzesunieciu(opis.id, 1));

  const zdejmij = przyciskSekcji(
    'Zdejmij',
    `Zdejmij sekcję ${opis.tytul} z widoku — wraca z pasa sekcji zdjętych`,
  );
  zdejmij.addEventListener('click', () => czynnosci.przyZdjeciu(opis.id));

  const naglowek = document.createElement('header');
  naglowek.className = `dn-karta-naglowek ${przedrostek}-sekcja__naglowek`;
  naglowek.append(tytul, zwin, wGore, wDol, zdejmij);

  const cialo = document.createElement('div');
  cialo.className = `dn-karta-cialo ${przedrostek}-sekcja__cialo`;
  cialo.hidden = stan.collapsed;
  cialo.append(opis.tresc);

  const element = document.createElement('section');
  element.className = `dn-karta ${przedrostek}-sekcja`;
  element.dataset['sekcja'] = opis.id;
  element.dataset['kolejnosc'] = String(stan.order);
  element.setAttribute('aria-label', `Sekcja panelu: ${opis.tytul}`);
  element.append(naglowek, cialo);
  return element;
}

/**
 * Pas sekcji zdjętych z widoku jest jedyną drogą powrotu, bo bez niego zdjęcie byłoby skutkiem nieodwracalnym z poziomu widoku.
 */
export function zlozPasZdjetych(
  zdjete: ReadonlyArray<{ id: string; tytul: string }>,
  przedrostek: string,
  przyPrzywroceniu: (id: string) => void,
): HTMLElement {
  const element = document.createElement('div');
  element.className = `${przedrostek}-sekcje__zdjete`;
  element.hidden = zdjete.length === 0;
  if (zdjete.length === 0) return element;

  const podpis = document.createElement('span');
  podpis.className = 'dn-pole-opis';
  podpis.textContent = 'Sekcje zdjęte z widoku:';
  element.append(podpis);

  for (const sekcja of zdjete) {
    const kontrolka = przyciskSekcji(sekcja.tytul, `Przywróć sekcję ${sekcja.tytul} do panelu`);
    kontrolka.addEventListener('click', () => przyPrzywroceniu(sekcja.id));
    element.append(kontrolka);
  }
  return element;
}
