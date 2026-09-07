/**
 * Wspólne dla czynności okna Translate: pas działań nad panelem, pytanie
 * o wartość, wskazanie panelu języka, potwierdzenie czynności nieodwracalnej
 * i wydanie ścieżki pliku do schowka.
 */

import type { ErrorInfo, TranslationPanel } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { oglos } from './ogloszenie.ts';
import { cialoPanelu, wykazPanelu } from './okno-modulu.ts';

export const NAGLOWEK = 'Tłumaczenie';

const CZAS_POTWIERDZENIA = 6000;

/** Stan okna tłumaczenia dzielony przez wszystkie czynności karty. */
export interface StanTlumaczenia {
  kanal: Kanal;
  korzen: Element;
  idOkna: () => string;
  odswiez: () => Promise<void>;
  panele: TranslationPanel[];
  idDokumentu: string;
  idZasobu: string;
  idPakietu: string;
}

/** Czynności jednej rodziny: kod znacznika przycisku wskazuje wykonanie. */
export type CzynnosciTlumaczenia = Record<string, () => Promise<void>>;

export function zapytaj(pytanie: string, domyslna = ''): string {
  return (globalThis.prompt?.(pytanie, domyslna) ?? '').trim();
}

export function zapytajOLiczbe(pytanie: string, domyslna = ''): number | null {
  const wpis = zapytaj(pytanie, domyslna);
  if (wpis === '') return null;
  const wartosc = Number(wpis);
  if (!Number.isFinite(wartosc)) {
    oglos(NAGLOWEK, `Wartość „${wpis}” nie jest liczbą.`, 'ostrzezenie');
    return null;
  }
  return wartosc;
}

/* Kontrakt przyjmuje w tych polach wyłącznie wartości ze swojego zbioru, więc
   wpis spoza zbioru zatrzymuje się w oknie zamiast wracać odmową rdzenia. */
export function zapytajOWartosc<T extends string>(
  zbior: Record<string, T>,
  pytanie: string,
  domyslna: T,
): T | null {
  const dopuszczalne = Object.values(zbior);
  const wpis = zapytaj(`${pytanie} (${dopuszczalne.join(', ')})`, domyslna);
  if (wpis === '') return null;
  const wskazana = dopuszczalne.find((wartosc) => wartosc === wpis);
  if (wskazana === undefined) {
    oglos(NAGLOWEK, `Rdzeń nie zna wartości „${wpis}”.`, 'ostrzezenie');
    return null;
  }
  return wskazana;
}

export function zapytajOWartosci<T extends string>(
  zbior: Record<string, T>,
  pytanie: string,
  domyslne: readonly T[],
): T[] {
  const dopuszczalne = Object.values(zbior);
  const wpis = zapytaj(
    `${pytanie} (${dopuszczalne.join(', ')}) — po przecinku`,
    domyslne.join(','),
  );
  const wskazane = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  const przyjete = wskazane.flatMap((czesc) => dopuszczalne.filter((wartosc) => wartosc === czesc));
  if (przyjete.length !== wskazane.length) {
    oglos(NAGLOWEK, 'Wśród wskazanych wartości jest taka, której rdzeń nie zna.', 'ostrzezenie');
    return [];
  }
  return przyjete;
}

export function odmowa(blad: ErrorInfo | undefined, zdanie: string): void {
  oglos(NAGLOWEK, blad?.message ?? zdanie, 'ostrzezenie');
}

export function powiedz(zdanie: string): void {
  oglos(NAGLOWEK, zdanie);
}

export function oknoWskazane(stan: StanTlumaczenia): string {
  const idOkna = stan.idOkna();
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Okno tłumaczenia jeszcze nie powstało, więc rdzeń nie ma czego wskazać.',
      'ostrzezenie');
  }
  return idOkna;
}

/* Kontrakt nie zna wykazu paneli języków: panel staje się znany oknu dopiero
   z odpowiedzi, która go niesie. Stąd wskazanie idzie po zapamiętanych. */
export function przyjmijPanele(stan: StanTlumaczenia, panele: readonly TranslationPanel[]): void {
  for (const panel of panele) {
    const miejsce = stan.panele.findIndex((znany) => znany.id === panel.id);
    if (miejsce === -1) stan.panele.push(panel);
    else stan.panele[miejsce] = panel;
  }
}

export function panelDocelowy(stan: StanTlumaczenia): string {
  const pierwszy = stan.panele[0];
  if (pierwszy === undefined) {
    oglos(NAGLOWEK, 'Okno nie zna jeszcze żadnego panelu języka — dołóż panel albo wczytaj źródło.',
      'ostrzezenie');
    return '';
  }
  if (stan.panele.length === 1) return pierwszy.id;
  const jezyki = stan.panele.map((panel) => panel.language).join(', ');
  const jezyk = zapytaj(`Język panelu (${jezyki})`, pierwszy.language);
  if (jezyk === '') return '';
  const wskazany = stan.panele.find(
    (panel) => panel.language.toLowerCase() === jezyk.toLowerCase(),
  );
  if (wskazany === undefined) {
    oglos(NAGLOWEK, `Okno nie zna panelu języka „${jezyk}”.`, 'ostrzezenie');
    return '';
  }
  return wskazany.id;
}

const OSTRZEZONE = new Set<string>();

/* Czynność nieodwracalna wchodzi dopiero za drugim naciśnięciem; pierwsze
   nazywa skutek. Zgoda gaśnie sama, żeby nie przeszła na późniejszy gest. */
export function potwierdzNieodwracalna(kod: string, zdanie: string): boolean {
  if (OSTRZEZONE.delete(kod)) return true;
  OSTRZEZONE.add(kod);
  globalThis.setTimeout(() => OSTRZEZONE.delete(kod), CZAS_POTWIERDZENIA);
  oglos(NAGLOWEK, `${zdanie} Naciśnij ponownie, aby wykonać.`, 'ostrzezenie');
  return false;
}

/* Okno nie ma miejsca, w którym zostawi plik: rdzeń zapisuje go po swojej
   stronie, a Operator dostaje potwierdzoną ścieżkę do schowka. */
export async function wydajSciezke(sciezka: string, opis: string): Promise<void> {
  if (sciezka === '') {
    oglos(NAGLOWEK, `${opis} Rdzeń nie wytworzył pliku.`, 'ostrzezenie');
    return;
  }
  try {
    await navigator.clipboard.writeText(sciezka);
    oglos(NAGLOWEK, `${opis} Ścieżka pliku jest w schowku: „${sciezka}”.`);
  } catch {
    oglos(NAGLOWEK,
      `${opis} Plik stoi pod „${sciezka}”, ale przeglądarka nie dała dostępu do schowka.`,
      'ostrzezenie');
  }
}

export function pokazWynik(
  stan: StanTlumaczenia,
  panelId: string,
  pozycje: readonly (readonly [string, string])[],
  pusty: string,
): void {
  wykazPanelu(cialoPanelu(stan.korzen, panelId), true, undefined, [...pozycje], pusty, (p) => p);
}

/* Pas działań stoi nad ciałem panelu, bo ciało jest wymieniane w całości przy
   każdym odświeżeniu wykazu i wszystko w nim postawione znikłoby razem z nim. */
export function pasDzialan(
  korzen: Element,
  panelId: string,
  podpis: string,
  czynnosci: readonly (readonly [string, string])[],
): void {
  const cialo = cialoPanelu(korzen, panelId);
  const sekcja = cialo?.parentElement ?? null;
  if (cialo === null || sekcja === null) return;
  const dokument = cialo.ownerDocument;
  const pas = dokument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const nazwa = dokument.createElement('span');
  nazwa.className = 'dn-meta';
  nazwa.textContent = podpis;
  pas.append(nazwa);
  for (const [kod, etykieta] of czynnosci) {
    const przycisk = dokument.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
    przycisk.dataset.czynnoscTlumaczenia = kod;
    przycisk.textContent = etykieta;
    pas.append(przycisk);
  }
  sekcja.insertBefore(pas, cialo);
}
