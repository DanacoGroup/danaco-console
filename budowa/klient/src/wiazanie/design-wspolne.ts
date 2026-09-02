// Warstwa wspólna wiązania modułu Design: kontekst karty, dostęp do węzłów
// prototypu i jednolite zamknięcie wywołania kontraktu.

import type { Command, RequestOf, ResponseOf } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

export const KOD_MODULU_DESIGN = 'design';

export const NIEGOTOWE = 'Rdzeń nie podaje jeszcze danych dla tego widoku.';

export interface StanDesignu {
  idOkna: string;
  idSesji: string;
  idModulu: string;
  idPlanszy: string;
  idZasobu: string;
  idZestawuZetonow: string;
  idRamki: string;
  szukanie: string;
  odswiezenia: Map<string, () => void>;
}

export interface Kontekst {
  kanal: Kanal;
  korzen: HTMLElement;
  idKarty: string;
  przy: AddEventListenerOptions;
  stan: StanDesignu;
}

export function stanPusty(): StanDesignu {
  return {
    idOkna: '',
    idSesji: '',
    idModulu: '',
    idPlanszy: '',
    idZasobu: '',
    idZestawuZetonow: '',
    idRamki: '',
    szukanie: '',
    odswiezenia: new Map(),
  };
}

export function panel(korzen: Element, identyfikator: string): HTMLElement | null {
  return korzen.querySelector<HTMLElement>(`#${identyfikator}`);
}

export function tresc(wezel: Element | null): HTMLElement | null {
  return wezel?.querySelector<HTMLElement>('.sta-okno-tresc') ?? null;
}

export function znacznik(wezel: Element | null, wartosc: string): void {
  const cel = wezel?.querySelector<HTMLElement>('.sta-okno-znacznik');
  if (cel === null || cel === undefined) return;
  if (wartosc === '') {
    cel.remove();
    return;
  }
  cel.textContent = wartosc;
}

export function tytul(wezel: Element | null, wartosc: string): void {
  const cel = wezel?.querySelector<HTMLElement>('.sta-okno-tytul b');
  if (cel === null || cel === undefined || wartosc === '') return;
  cel.textContent = wartosc;
}

export function pusto(wezel: Element | null): void {
  wezel?.replaceChildren();
}

// Prototyp nie ma węzła pustki, więc zdanie niegotowości wchodzi w akapit meta.
export function nieGotowe(wezel: Element | null, zdanie: string = NIEGOTOWE): void {
  if (wezel === null) return;
  const napis = wezel.ownerDocument.createElement('div');
  napis.className = 'dn-meta';
  napis.textContent = zdanie;
  wezel.replaceChildren(napis);
}

export function chip(zakres: Element | null, podpis: string): HTMLElement | null {
  if (zakres === null) return null;
  const chipy = [...zakres.querySelectorAll<HTMLElement>('.sta-chip')];
  return chipy.find((kandydat) => (kandydat.textContent ?? '').trim().startsWith(podpis)) ?? null;
}

export function przycisk(zakres: Element | null, podpis: string): HTMLElement | null {
  if (zakres === null) return null;
  const przyciski = [...zakres.querySelectorAll<HTMLElement>('button')];
  return przyciski.find((kandydat) => (kandydat.textContent ?? '').trim().startsWith(podpis))
    ?? null;
}

export function naKlik(
  kontekst: Kontekst,
  wezel: Element | null,
  czyn: () => void | Promise<void>,
): void {
  if (wezel === null) return;
  if (wezel instanceof HTMLElement) wezel.setAttribute('role', 'button');
  wezel.addEventListener('click', (zdarzenie) => {
    zdarzenie.preventDefault();
    void czyn();
  }, kontekst.przy);
}

export function zdejmij(wezel: Element | null): void {
  wezel?.remove();
}

export function podpisRodzaju(rodzaj: string): string {
  const nazwy: Record<string, string> = {
    image: 'obraz',
    vector: 'wektor',
    composition: 'kompozycja',
    document: 'dokument',
    audio: 'dźwięk',
    video: 'wideo',
    archive: 'archiwum',
  };
  return nazwy[rodzaj] ?? rodzaj;
}

export function rozmiar(bajty: number): string {
  if (bajty < 1024) return `${String(bajty)} B`;
  if (bajty < 1024 * 1024) return `${(bajty / 1024).toFixed(0)} kB`;
  return `${(bajty / (1024 * 1024)).toFixed(1)} MB`;
}

export function chwila(znacznikCzasu: number): string {
  return new Date(znacznikCzasu).toLocaleString('pl-PL', {
    day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit',
  });
}

export async function poproszony<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<ResponseOf<K> | null> {
  const wynik: Wynik<ResponseOf<K>> = await wywolaj(kanal, komenda, zadanie);
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return wynik.wynik;
}

export async function wykonany<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
  podpis: string,
): Promise<ResponseOf<K> | null> {
  const wynik: Wynik<ResponseOf<K>> = await wywolaj(kanal, komenda, zadanie);
  if (!wynik.udany) {
    oglos('Design', `${podpis}: ${wynik.blad?.message ?? 'rdzeń odmówił bez opisu.'}`, 'blad');
    return null;
  }
  oglos('Design', `${podpis} — wykonane.`, 'informacja');
  return wynik.wynik ?? null;
}

export function zapowiedzBraku(podpis: string): void {
  oglos('Design', `${podpis} czeka na wskazanie w oknie.`, 'ostrzezenie');
}

export interface Wiersz {
  tekst: string;
  meta?: string;
  plakietka?: string;
  kropka?: 'sukces' | 'neutralna' | 'sygnal';
  wybrany?: boolean;
  naKlik?: () => void;
}

export function wypelnijWykaz(
  kontekst: Kontekst,
  lista: HTMLElement | null,
  wiersze: Wiersz[],
  zdanieBraku = NIEGOTOWE,
): void {
  if (lista === null) return;
  if (wiersze.length === 0) {
    nieGotowe(lista, zdanieBraku);
    return;
  }
  const dokument = lista.ownerDocument;
  const zbudowane = wiersze.map((wiersz) => {
    const poz = dokument.createElement('div');
    poz.className = 'dn-wykaz-modulu-poz';
    if (wiersz.kropka !== undefined) {
      const kropka = dokument.createElement('span');
      kropka.className = `dn-kropka dn-kropka--${wiersz.kropka}`;
      kropka.setAttribute('aria-hidden', 'true');
      poz.append(kropka, ' ');
    }
    poz.append(wiersz.tekst);
    if (wiersz.meta !== undefined && wiersz.meta !== '') {
      const meta = dokument.createElement('span');
      meta.className = 'dn-meta';
      meta.textContent = wiersz.meta;
      poz.append(' ', meta);
    }
    if (wiersz.plakietka !== undefined && wiersz.plakietka !== '') {
      const plakietka = dokument.createElement('span');
      plakietka.className = 'dn-plakietka dn-na-koniec';
      plakietka.textContent = wiersz.plakietka;
      poz.append(' ', plakietka);
    }
    if (wiersz.wybrany === true) poz.setAttribute('aria-current', 'true');
    if (wiersz.naKlik !== undefined) naKlik(kontekst, poz, wiersz.naKlik);
    return poz;
  });
  lista.replaceChildren(...zbudowane);
}
