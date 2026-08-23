import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Panel rodzin rdzenia modułu Browser — droga z okna do tych komend obszaru
 * `browser.*`, które nie mają własnej kontrolki w oknach opisanych
 * opracowaniem: kart i przestrzeni roboczych, monitorów, kanałów, kolejki
 * czytania, zakładek, pobrań, wytworów, zrzutów, narzędzi inspekcyjnych,
 * nagrywarki makr i granic Wykonawcy.
 *
 * Jedna odpowiedzialność: postawić czynność i pokazać to, co rdzeń naprawdę
 * oddał. Panel nie pamięta niczego między naciśnięciami — stan modułu stoi
 * w `stan-przegladania.ts`, a wykaz wyniku jest odczytem, nie kopią.
 *
 * Dlaczego jeden panel, a nie kontrolka przy każdej pozycji: rodzin jest
 * kilkanaście, a każda ma dwa–trzy pola. Rozsypane po oknach dałyby kilkadziesiąt
 * kontrolek w miejscach, w których Operator ich nie szuka; zebrane w sekcje
 * przy oknie, do którego opracowanie je przypisuje, zostają w zasięgu jednego
 * kliknięcia i nie zasłaniają pracy podstawowej.
 *
 * Panel mówi prawdę o odmowie: kod i treść odmowy rdzenia idą wprost do wiersza
 * odpowiedzi. Cisza po naciśnięciu — albo zdanie „gotowe" bez pokrycia — byłaby
 * tą samą szkodą, przed którą stoją sprawdziany skutku po stronie rdzenia.
 */
export interface PanelRodzin {
  element: HTMLElement;
}

/** Opis jednej czynności panelu: nazwa, pola i wywołanie rdzenia. */
interface Czynnosc {
  nazwa: string;
  /** Pola wypełniane przez Operatora; wartości idą do wywołania w kolejności. */
  pola?: readonly { klucz: string; etykieta: string }[];
  /** Wywołanie rdzenia; oddaje zdanie o skutku albo rzuca odmowę zdaniem. */
  wykonaj(wartosci: Record<string, string>): Promise<string>;
}

/** Sekcja panelu — rodzina komend wraz z jej czynnościami. */
export interface SekcjaRodzin {
  tytul: string;
  opis: string;
  czynnosci: readonly Czynnosc[];
}

/**
 * Składa panel z sekcji. Każda sekcja dostaje własny wykaz pól i rząd
 * przycisków; wynik czynności ląduje we wspólnym wierszu odpowiedzi panelu.
 */
export function utworzPanelRodzin(sekcje: readonly SekcjaRodzin[]): PanelRodzin {
  const odpowiedz = utworzWierszOdpowiedzi();
  const element = document.createElement('section');
  element.className = 'mb-rodziny';
  element.setAttribute('aria-label', 'Rodziny komend modułu Browser');

  for (const sekcja of sekcje) {
    const blok = document.createElement('div');
    blok.className = 'mb-rodziny__sekcja';

    const naglowek = document.createElement('h4');
    naglowek.className = 'mb-rodziny__tytul';
    naglowek.textContent = sekcja.tytul;

    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent = sekcja.opis;

    const pola = new Map<string, HTMLInputElement>();
    const wiersz = document.createElement('div');
    wiersz.className = 'mb-rodziny__pola';
    for (const czynnosc of sekcja.czynnosci) {
      for (const pole of czynnosc.pola ?? []) {
        if (pola.has(pole.klucz)) continue;
        const kontrolka = document.createElement('input');
        kontrolka.type = 'text';
        kontrolka.className = 'dn-pole-kontrolka mb-rodziny__pole';
        kontrolka.placeholder = pole.etykieta;
        kontrolka.setAttribute('aria-label', pole.etykieta);
        pola.set(pole.klucz, kontrolka);
        wiersz.append(kontrolka);
      }
    }

    const pasek = document.createElement('div');
    pasek.className = 'mb-rodziny__pasek';
    for (const czynnosc of sekcja.czynnosci) {
      const guzik = przycisk(czynnosc.nazwa, 'dn-btn dn-btn--zarys dn-btn--sm');
      guzik.dataset['rodzina'] = sekcja.tytul;
      guzik.addEventListener('click', () => {
        const wartosci: Record<string, string> = {};
        for (const [klucz, kontrolka] of pola) wartosci[klucz] = kontrolka.value.trim();
        odpowiedz.pokaz(`${czynnosc.nazwa}: w toku…`, true);
        void czynnosc
          .wykonaj(wartosci)
          .then((zdanie) => odpowiedz.pokaz(zdanie, true))
          .catch((powod: unknown) => odpowiedz.pokaz(String(powod), false));
      });
      pasek.append(guzik);
    }

    blok.append(naglowek, opis, wiersz, pasek);
    element.append(blok);
  }

  element.append(odpowiedz.element);
  return { element };
}

/**
 * Odmowa rdzenia zamieniona na zdanie, którym panel kończy czynność.
 * Rzucona, a nie zwrócona: czynność, która się nie udała, nie ma prawa
 * wyglądać jak udana.
 */
export function odrzuc(czynnosc: string, blad?: { code?: string; message?: string }): never {
  throw new Error(opisOdmowy(czynnosc, blad?.code, blad?.message));
}

/** Okno przeglądarki albo odmowa nazywająca jego brak. */
export function oknoAlboOdmowa(stan: StanPrzegladania): string {
  const idOkna = stan.idOkna();
  if (idOkna === '') throw new Error(stan.powod());
  return idOkna;
}

/** Wartość pola albo odmowa nazywająca, czego brakuje. */
export function wymagane(wartosci: Record<string, string>, klucz: string, nazwa: string): string {
  const wartosc = wartosci[klucz] ?? '';
  if (wartosc === '') throw new Error(`Podaj ${nazwa} — bez tego rdzeń nie ma czego wykonać.`);
  return wartosc;
}

/** Wartość pola albo `undefined` dla pola pustego. */
export function opcjonalne(wartosci: Record<string, string>, klucz: string): string | undefined {
  const wartosc = wartosci[klucz] ?? '';
  return wartosc === '' ? undefined : wartosc;
}

/** Liczba z pola albo `undefined`; wartość nieliczbowa jest odmową. */
export function liczbaPola(wartosci: Record<string, string>, klucz: string): number | undefined {
  const wartosc = wartosci[klucz] ?? '';
  if (wartosc === '') return undefined;
  const liczba = Number(wartosc);
  if (!Number.isFinite(liczba)) throw new Error(`Wartość „${wartosc}" nie jest liczbą.`);
  return liczba;
}
