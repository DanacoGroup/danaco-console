import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Panel rodzin rdzenia modułu Browser — droga z okna do komend przeglądania,
 * które nie mają własnej kontrolki w innych oknach modułu. Jedna
 * odpowiedzialność: postawić czynność i pokazać to, co rdzeń naprawdę oddał.
 */
export interface PanelRodzin {
  element: HTMLElement;
}

/** Opis jednej czynności panelu: nazwa, pola wypełniane przez operatora i wywołanie właściwego polecenia rdzenia. */
interface Czynnosc {
  nazwa: string;
  /** Pola wypełniane przez Operatora; wartości idą do wywołania w kolejności. */
  pola?: readonly { klucz: string; etykieta: string }[];
  /** Wywołanie rdzenia; oddaje zdanie o skutku albo rzuca odmowę zdaniem. */
  wykonaj(wartosci: Record<string, string>): Promise<string>;
}

/** Sekcja panelu rodzin — jedna rodzina komend rdzenia wraz z wykazem jej czynności oraz tytułem tej sekcji. */
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

/** Okno przeglądarki bieżącej sesji przeglądania albo odmowa nazywająca wprost jego brak dla operatora. */
export function oknoAlboOdmowa(stan: StanPrzegladania): string {
  const idOkna = stan.idOkna();
  if (idOkna === '') throw new Error(stan.powod());
  return idOkna;
}

/** Wartość pola formularza panelu albo odmowa nazywająca wprost, jakiego dokładnie pola brakuje operatorowi. */
export function wymagane(wartosci: Record<string, string>, klucz: string, nazwa: string): string {
  const wartosc = wartosci[klucz] ?? '';
  if (wartosc === '') throw new Error(`Podaj ${nazwa} — bez tego rdzeń nie ma czego wykonać.`);
  return wartosc;
}

/** Wartość pola formularza albo wartość pusta dla pola pozostawionego całkiem bez wypełnienia przez operatora. */
export function opcjonalne(wartosci: Record<string, string>, klucz: string): string | undefined {
  const wartosc = wartosci[klucz] ?? '';
  return wartosc === '' ? undefined : wartosc;
}

/** Liczba odczytana z pola formularza panelu albo wartość pusta; wartość nieliczbowa kończy się odmową. */
export function liczbaPola(wartosci: Record<string, string>, klucz: string): number | undefined {
  const wartosc = wartosci[klucz] ?? '';
  if (wartosc === '') return undefined;
  const liczba = Number(wartosc);
  if (!Number.isFinite(liczba)) throw new Error(`Wartość „${wartosc}" nie jest liczbą.`);
  return liczba;
}
