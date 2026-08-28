/**
 * Wspólny kształt paneli dopełniających okna modułu Automations, bo cztery panele mają tę samą
 * budowę: tytuł, zdanie o przeznaczeniu, pola i pas przycisków.
 */
import { przyciskAkcji, wiersz } from '../../modele/kontrolki-formularza';
import type { Wynik } from '../../protokol/kanal';
import { utworzStanTresci, type StanTresci } from './stany-okna';

/** Panel dopełniający jedno okno operacyjne, osadzany wewnątrz niego wraz z polami, przyciskami i nośnikiem stanu treści. */
export interface PanelDobudowy {
  /** Element osadzany wewnątrz okna. */
  element: HTMLElement;
  /** Nośnik stanu treści panelu — ładowanie, odmowa, potwierdzenie. */
  tresc: StanTresci;
  /** Dokłada wiersz pola do siatki panelu i oddaje kontrolkę. */
  dodajPole<T extends HTMLElement>(etykieta: string, kontrolka: T, objasnienie?: string): T;
  /** Dokłada przycisk czynności do pasa akcji panelu. */
  dodajCzynnosc(etykieta: string, czynnosc: () => void): HTMLButtonElement;
}

/** Składa panel o podanym tytule i przeznaczeniu, gotowy do dołożenia pól i czynności przez okno operacyjne. */
export function utworzPanelDobudowy(tytul: string, przeznaczenie: string): PanelDobudowy {
  const naglowek = document.createElement('h4');
  naglowek.className = 'da-panel__tytul';
  naglowek.textContent = tytul;

  const opis = document.createElement('p');
  opis.className = 'da-panel__opis';
  opis.textContent = przeznaczenie;

  const pola = document.createElement('div');
  pola.className = 'da-panel__pola';

  const akcje = document.createElement('div');
  akcje.className = 'da-panel__akcje';

  const tresc = utworzStanTresci();

  const element = document.createElement('section');
  element.className = 'da-panel';
  element.setAttribute('aria-label', tytul);
  element.append(naglowek, opis, pola, akcje, tresc.element);

  return {
    element,
    tresc,
    dodajPole(etykieta, kontrolka, objasnienie) {
      const opcje: { klasa: string; objasnienie?: string } = { klasa: 'da-panel__pole' };
      if (objasnienie !== undefined) opcje.objasnienie = objasnienie;
      pola.append(wiersz(etykieta, kontrolka, opcje));
      return kontrolka;
    },
    dodajCzynnosc(etykieta, czynnosc) {
      const przycisk = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
      przycisk.addEventListener('click', czynnosc);
      akcje.append(przycisk);
      return przycisk;
    },
  };
}

/**
 * Wykonanie czynności panelu kończy się ładowaniem, potem odmową nazwaną powodem rdzenia albo
 * potwierdzeniem z zapisem odpowiedzi.
 */
export function wykonajCzynnoscPanelu<T>(
  panel: PanelDobudowy,
  opisPracy: string,
  obietnica: Promise<Wynik<T>>,
  zdanieOdmowy: string,
  zdaniePowodzenia: string,
): void {
  panel.tresc.ladowanie(opisPracy);
  void obietnica.then((wynik) => {
    if (!wynik.udany || wynik.wynik === undefined) {
      panel.tresc.blad(zdanieOdmowy, wynik.blad);
      return;
    }
    panel.tresc.potwierdzenie(zdaniePowodzenia, true);
    panel.tresc.tresc().replaceChildren(zapisOdpowiedzi(wynik.wynik));
  });
}

/**
 * Zapis odpowiedzi rdzenia w postaci czytelnej dla człowieka; wartości wrażliwych tu nie ma
 * i być nie może.
 */
function zapisOdpowiedzi(wynik: unknown): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'da-panel__odpowiedz';
  element.textContent = JSON.stringify(wynik, null, 2);
  return element;
}

/**
 * Odczyt zapisu strukturalnego z pola tekstowego: pole puste daje undefined, zapis nieczytelny
 * daje null.
 */
export function odczytajZapis(tekst: string): unknown | null | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  try {
    return JSON.parse(oczyszczony) as unknown;
  } catch {
    return null;
  }
}

/** Liczba z pola liczbowego; pole puste albo nieczytelne daje undefined zamiast rzucać wyjątkiem parsowania. */
export function odczytajLiczbe(tekst: string): number | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  const liczba = Number(oczyszczony);
  return Number.isFinite(liczba) ? liczba : undefined;
}

/** Chwila z pola tekstowego w postaci daty; puste pole daje undefined, nieczytelny zapis też daje undefined. */
export function odczytajChwile(tekst: string): number | undefined {
  const oczyszczony = tekst.trim();
  if (oczyszczony === '') return undefined;
  const chwila = Date.parse(oczyszczony);
  return Number.isFinite(chwila) ? chwila : undefined;
}

/** Wykaz z pola tekstowego rozdzielonego przecinkami; puste pole daje wykaz pusty, nie undefined ani null. */
export function odczytajWykaz(tekst: string): string[] {
  return tekst
    .split(',')
    .map((pozycja) => pozycja.trim())
    .filter((pozycja) => pozycja !== '');
}
