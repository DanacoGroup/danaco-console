import type { PanelSection } from '../../../shared/contract';
import { przyciskAkcji } from '../modele/kontrolki-formularza';
import type { OpisSekcji } from './uklad-sekcji';

/**
 * Kształt jednej sekcji panelu — sam widok, bez reguły układu i bez komendy.
 *
 * Jedna odpowiedzialność: zamiana pary (opis sekcji, jej stan w układzie) na
 * węzły drzewa dokumentu. Kto stoi przed kim, rozstrzyga `uklad-sekcji.ts`;
 * czym się to zapisuje — `zrodlo-sekcji-paneli.ts`. Ten plik nie wie ani jednego,
 * ani drugiego: cztery czynności oddaje wywołaniem zwrotnym, więc da się go
 * czytać i sprawdzać bez rdzenia.
 *
 * Ani jeden przycisk nie jest wygaszany i ani jeden nie pyta „czy na pewno?" —
 * zdjęcie sekcji z widoku wykonuje się od razu, a wraca ją pas sekcji zdjętych
 * stojący pod panelem. Skutek jest odwracalny jednym naciśnięciem, więc
 * potwierdzanie go byłoby przeszkodą bez treści.
 *
 * Stan idzie atrybutem, nie samą barwą. Zwinięcie niesie `aria-expanded` na
 * przycisku i `hidden` na ciele, a nie klasa zmieniająca wygląd — inaczej
 * czytnik ekranu ogłaszałby treść, której na ekranie nie ma.
 *
 * Wygląd w całości z biblioteki (`komponenty/karta.css`, `przycisk.css`) —
 * plik nie zna ani jednej barwy i ani jednego odstępu.
 */

/** Cztery czynności, które Operator wykonuje na sekcji panelu. */
export interface CzynnosciSekcji {
  przyZwinieciu(id: string): void;
  przyPrzesunieciu(id: string, kierunek: -1 | 1): void;
  przyZdjeciu(id: string): void;
  przyPrzywroceniu(id: string): void;
}

/** Przycisk czynności sekcji — jeden kształt dla wszystkich czterech. */
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
 * Pas sekcji zdjętych z widoku — jedyna droga powrotu.
 *
 * Bez niego zdjęcie byłoby skutkiem nieodwracalnym z poziomu widoku, a to
 * dopiero kazałoby pytać „czy na pewno?". Pas znika, gdy nie ma czego wracać:
 * pusty pas z napisem „brak" zabierałby miejsce, nic nie mówiąc.
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
