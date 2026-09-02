// Wspólny stan i pomocniki okna Library. Znacznik Właściciela nie zna stanu
// „brak danych", więc miejsce po wzorze niesie jedno zdanie niegotowości.
import type { LibraryFile } from '../../../shared/contract.ts';
import type { Wynik } from '../protokol/kanal.ts';

export const ZDANIE_PUSTO = 'Rdzeń nie podał jeszcze danych dla tego panelu.';

export interface Kontekst {
  plik: () => LibraryFile | null;
  kolekcja: () => string;
  ustawPlik: (plik: LibraryFile | null) => void;
  ustawKolekcje: (id: string) => void;
  przyPliku: (uchwyt: (plik: LibraryFile | null) => void) => void;
  przyKolekcji: (uchwyt: (id: string) => void) => void;
  zglosOdswiezenie: (uchwyt: () => void) => void;
  odswiez: () => void;
}

export function utworzKontekst(): Kontekst {
  let wybrany: LibraryFile | null = null;
  let kolekcja = '';
  const przyPliku: ((plik: LibraryFile | null) => void)[] = [];
  const przyKolekcji: ((id: string) => void)[] = [];
  const odswiezenia: (() => void)[] = [];
  return {
    plik: () => wybrany,
    kolekcja: () => kolekcja,
    ustawPlik: (plik) => {
      wybrany = plik;
      for (const uchwyt of przyPliku) uchwyt(plik);
    },
    ustawKolekcje: (id) => {
      kolekcja = id;
      for (const uchwyt of przyKolekcji) uchwyt(id);
    },
    przyPliku: (uchwyt) => przyPliku.push(uchwyt),
    przyKolekcji: (uchwyt) => przyKolekcji.push(uchwyt),
    zglosOdswiezenie: (uchwyt) => odswiezenia.push(uchwyt),
    odswiez: () => {
      for (const uchwyt of odswiezenia) uchwyt();
    },
  };
}

export function wartosc<T>(odpowiedz: Wynik<T>): T | null {
  return odpowiedz.udany && odpowiedz.wynik !== undefined ? odpowiedz.wynik : null;
}

export function zdejmijDzieci(wezel: Element | null): void {
  wezel?.replaceChildren();
}

export function pustka(wezel: Element | null, zdanie: string = ZDANIE_PUSTO): void {
  if (wezel === null) return;
  wezel.replaceChildren();
  const nota = wezel.ownerDocument.createElement('div');
  nota.className = 'dn-meta';
  nota.textContent = zdanie;
  wezel.appendChild(nota);
}

export function pustkaWiersza(cialo: Element | null, kolumn: number, zdanie: string): void {
  if (cialo === null) return;
  cialo.replaceChildren();
  const wiersz = cialo.ownerDocument.createElement('tr');
  const komorka = cialo.ownerDocument.createElement('td');
  komorka.colSpan = kolumn;
  komorka.className = 'dn-meta';
  komorka.textContent = zdanie;
  wiersz.appendChild(komorka);
  cialo.appendChild(wiersz);
}

const JEDNOSTKI = ['B', 'kB', 'MB', 'GB', 'TB'];

export function rozmiar(bajty: number | undefined): string {
  if (bajty === undefined) return '';
  let wartoscMiary = bajty;
  let rzad = 0;
  while (wartoscMiary >= 1000 && rzad < JEDNOSTKI.length - 1) {
    wartoscMiary /= 1000;
    rzad += 1;
  }
  return `${wartoscMiary.toLocaleString('pl-PL', { maximumFractionDigits: 1 })} ${JEDNOSTKI[rzad]}`;
}

export function dzien(znacznik: number | undefined): string {
  if (znacznik === undefined) return '';
  const chwila = new Date(znacznik);
  const miesiac = String(chwila.getMonth() + 1).padStart(2, '0');
  const dzienMiesiaca = String(chwila.getDate()).padStart(2, '0');
  return `${String(chwila.getFullYear())}-${miesiac}-${dzienMiesiaca}`;
}

export function skrocony(znacznik: number | undefined): string {
  const pelny = dzien(znacznik);
  return pelny === '' ? '' : pelny.slice(5);
}

export function wzorZ(rodzic: Element | null, wybor: string): HTMLElement | null {
  const wezel = rodzic?.querySelector(wybor) ?? null;
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

export function panel(korzen: ParentNode, identyfikator: string): HTMLElement | null {
  const wezel = korzen.querySelector(`#${identyfikator} .sta-okno-tresc`);
  return wezel instanceof HTMLElement ? wezel : null;
}

export function wpisz(wezel: Element | null, tekst: string): void {
  if (wezel === null) return;
  if (tekst === '') {
    wezel.remove();
    return;
  }
  wezel.textContent = tekst;
}

export function plakietka(dokument: Document, tresc: string, odmiana = ''): HTMLElement {
  const znak = dokument.createElement('span');
  znak.className = odmiana === '' ? 'dn-plakietka' : `dn-plakietka dn-plakietka--${odmiana}`;
  znak.textContent = tresc;
  return znak;
}

export function pozycja(dokument: Document, tresc: string, miara: string): HTMLElement {
  const wiersz = dokument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  wiersz.append(tresc);
  if (miara !== '') {
    const meta = dokument.createElement('span');
    meta.className = 'dn-meta';
    meta.textContent = miara;
    wiersz.appendChild(meta);
  }
  return wiersz;
}

export function etykietaGrupy(dokument: Document, tresc: string): HTMLElement {
  const wezel = dokument.createElement('div');
  wezel.className = 'pt-etykieta';
  wezel.textContent = tresc;
  return wezel;
}
