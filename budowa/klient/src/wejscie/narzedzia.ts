/**
 * Narzędzia okien — warstwa pod składnikami. Trzy rzeczy, których potrzebuje
 * każde okno: budowanie węzła, sięganie po łańcuch z katalogu treści
 * i podstawianie danych w łańcuch.
 */

import { tresci, type WezelTresci } from './tresci.ts';

/** Wartość atrybutu węzła: null i false pomijają atrybut całkowicie, true stawia go pusty, bez wartości. */
export type WartoscAtrybutu = string | number | boolean | null | undefined;

/** Atrybuty węzła przekazywane budowie; nazwy własne — klasa, tekst, dane — reszta idzie wprost jako atrybut HTML. */
export interface AtrybutyWezla {
  klasa?: string | null;
  tekst?: string | null;
  dane?: Record<string, WartoscAtrybutu> | null;
  [nazwa: string]: WartoscAtrybutu | Record<string, WartoscAtrybutu> | null | undefined;
}

/** Dziecko węzła: sam węzeł, łańcuch tekstu albo nic, gdy warunek składnika nie chce go wcale wstawiać. */
export type Dziecko = Node | string | null | undefined | false;

/**
 * Budowa węzła. Atrybuty rozpoznawane po nazwie: klasa, tekst, dane, reszta
 * wprost jako atrybut. Wartość null albo false pomija atrybut — dzięki temu
 * warianty składnika pisze się warunkiem, nie rozgałęzieniem.
 */
export function el(znacznik: string, atrybuty?: AtrybutyWezla, dzieci?: Dziecko[]): HTMLElement {
  const wezel = document.createElement(znacznik);
  for (const [nazwa, wartosc] of Object.entries(atrybuty ?? {})) {
    if (wartosc === null || wartosc === undefined || wartosc === false) continue;
    if (nazwa === 'klasa') wezel.className = String(wartosc);
    else if (nazwa === 'tekst') wezel.textContent = String(wartosc);
    else if (nazwa === 'dane') {
      for (const [dana, wartoscDanej] of Object.entries(wartosc as Record<string, WartoscAtrybutu>)) {
        if (wartoscDanej === null || wartoscDanej === undefined || wartoscDanej === false) continue;
        wezel.setAttribute(`data-${dana}`, wartoscDanej === true ? '' : String(wartoscDanej));
      }
    } else if (wartosc === true) wezel.setAttribute(nazwa, '');
    else wezel.setAttribute(nazwa, String(wartosc));
  }
  for (const dziecko of dzieci ?? []) {
    if (dziecko === null || dziecko === undefined || dziecko === false) continue;
    wezel.appendChild(typeof dziecko === 'string' ? document.createTextNode(dziecko) : dziecko);
  }
  return wezel;
}

/**
 * Węzeł zbudowany ze znacznika. Używany wyłącznie dla znaków z zestawu — nigdy
 * dla danych z zewnątrz, bo znacznik z zewnątrz wykonałby się jak treść okna.
 */
export function zeZnacznika(znacznik: string): SVGElement {
  const szablon = document.createElement('template');
  szablon.innerHTML = znacznik.trim();
  const pierwszy = szablon.content.firstElementChild;
  // Wpis diagnostyczny dla wykonawcy, nie tekst okna: pusty znacznik znaku jest usterką zestawu.
  if (pierwszy === null) throw new Error('[wejscie] pusty znacznik znaku');
  return pierwszy as SVGElement;
}

/** Dane podstawiane w miejsca ujęte w nawiasy klamrowe wewnątrz łańcucha wziętego z tego katalogu treści. */
export type DanePodstawienia = Record<string, string | number>;

/**
 * Podstawienie danych w łańcuch. Nazwy w nawiasach są częścią kontraktu
 * z tłumaczem — nie tłumaczy się ich, a brak danej zostawia nawias nietknięty,
 * żeby luka rzucała się w oczy zamiast zniknąć.
 */
export function podstaw(wzor: string, dane?: DanePodstawienia): string {
  return wzor.replace(/\{(\w+)\}/g, (calosc, nazwa: string) => {
    const wartosc = dane?.[nazwa];
    return wartosc === undefined ? calosc : String(wartosc);
  });
}

/**
 * Sięgnięcie po treść ścieżką kluczy: `tekst('dostep.logowanie.tytul')`.
 *
 * Brak klucza nie jest sytuacją do obsłużenia po cichu — wraca sama ścieżka
 * w nawiasach kątowych, żeby usterka katalogu rzucała się w oczy w oknie.
 */
export function wezel(sciezka: string): WezelTresci {
  let biezacy: unknown = tresci;
  for (const czlon of sciezka.split('.')) {
    if (biezacy === null || typeof biezacy !== 'object') return `⟨${sciezka}⟩`;
    biezacy = (biezacy as Record<string, unknown>)[czlon];
  }
  if (biezacy === undefined || biezacy === null) return `⟨${sciezka}⟩`;
  return biezacy as WezelTresci;
}

/** Łańcuch z katalogu, wraz z podstawieniem danych, gdy je podano; brak klucza zwraca samą jego ścieżkę. */
export function tekst(sciezka: string, dane?: DanePodstawienia): string {
  const znaleziony = wezel(sciezka);
  if (typeof znaleziony !== 'string') return `⟨${sciezka}⟩`;
  return dane === undefined ? znaleziony : podstaw(znaleziony, dane);
}

/** Wykaz łańcuchów z katalogu, wskazany ścieżką kluczy prowadzącą prosto do tablicy wewnątrz drzewa treści. */
export function wykaz(sciezka: string): string[] {
  const znaleziony = wezel(sciezka);
  if (!Array.isArray(znaleziony)) return [`⟨${sciezka}⟩`];
  return znaleziony.map((pozycja) => (typeof pozycja === 'string' ? pozycja : `⟨${sciezka}⟩`));
}

/** Wykaz par głowa–treść z katalogu; używa go kolumna tożsamości do pokazania kolejnych zdań o programie. */
export function wykazZalet(sciezka: string): { glowa: string; tresc: string }[] {
  const znaleziony = wezel(sciezka);
  if (!Array.isArray(znaleziony)) return [];
  return znaleziony.map((pozycja) => {
    const zaleta = pozycja as { glowa?: WezelTresci; tresc?: WezelTresci };
    return {
      glowa: typeof zaleta.glowa === 'string' ? zaleta.glowa : `⟨${sciezka}⟩`,
      tresc: typeof zaleta.tresc === 'string' ? zaleta.tresc : `⟨${sciezka}⟩`,
    };
  });
}

/** Dane podstawienia z katalogu — na przykład adres nadawcy listu — złożone z par klucza i jego wartości. */
export function daneZKatalogu(sciezka: string): DanePodstawienia {
  const znaleziony = wezel(sciezka);
  if (znaleziony === null || typeof znaleziony !== 'object' || Array.isArray(znaleziony)) return {};
  const dane: DanePodstawienia = {};
  for (const [klucz, wartosc] of Object.entries(znaleziony)) {
    if (typeof wartosc === 'string' || typeof wartosc === 'number') dane[klucz] = wartosc;
  }
  return dane;
}

/** Czas w postaci minut i sekund oddzielonych dwukropkiem; postać należy do widoku, nie do katalogu treści. */
export function naZegar(sekundy: number): string {
  const minuty = Math.floor(Math.max(0, sekundy) / 60);
  const reszta = Math.max(0, sekundy) % 60;
  return `${minuty < 10 ? '0' : ''}${minuty}:${reszta < 10 ? '0' : ''}${reszta}`;
}
