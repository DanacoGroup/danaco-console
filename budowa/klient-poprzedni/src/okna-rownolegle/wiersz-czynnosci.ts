import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import './czynnosci.css';

/**
 * Czynność gotowa do postawienia jako wiersz sekcji czynności sesji w menu, odrębny od
 * przełącznika panelu, bo naciśnięcie wykonuje działanie, a nie włącza stan.
 */
export interface PozycjaCzynnosciMenu {
  /** Klucz techniczny — po nim idzie wybór skrótu i porządek. */
  klucz: string;
  /** Nazwa widoczna w menu. */
  nazwa: string;
  /** Zdanie pod nazwą: co czynność naprawdę robi. */
  przeznaczenie: string;
  ikona: NazwaIkony;
  /** Pominięty znaczy brak klawisza; wiersz bez skrótu nie dostaje pustego znacznika. */
  skrot?: string;
  /** Czy wiersz idzie kolorem ostrzegawczym (utrata danych bez odwrotu). */
  grozna?: boolean;
  wykonaj(): void;
}

/**
 * Buduje wiersz czynności menu: ikonę, nazwę, przeznaczenie i opcjonalny skrót
 * klawiszowy, z barwą ostrzegawczą dla czynności nieodwracalnych.
 */
export function zbudujWierszCzynnosci(pozycja: PozycjaCzynnosciMenu): HTMLElement {
  const wiersz = document.createElement('button');
  wiersz.type = 'button';
  wiersz.className = 'dn-czynnosci-sesji__pozycja';
  if (pozycja.grozna === true) wiersz.classList.add('dn-czynnosci-sesji__pozycja--grozna');
  wiersz.setAttribute('role', 'menuitem');
  wiersz.dataset.czynnosc = pozycja.klucz;
  wiersz.title = pozycja.przeznaczenie;

  const ikona = document.createElement('span');
  ikona.className = 'dn-czynnosci-sesji__ikona';
  ikona.append(elementIkony(pozycja.ikona, { rozmiar: 16 }));

  const tresc = document.createElement('span');
  tresc.className = 'dn-czynnosci-sesji__tresc';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-czynnosci-sesji__nazwa';
  nazwa.textContent = pozycja.nazwa;

  const przeznaczenie = document.createElement('span');
  przeznaczenie.className = 'dn-czynnosci-sesji__przeznaczenie';
  przeznaczenie.textContent = pozycja.przeznaczenie;

  tresc.append(nazwa, przeznaczenie);

  wiersz.append(ikona, tresc);

  if (pozycja.skrot !== undefined) {
    const skrot = document.createElement('kbd');
    skrot.className = 'dn-czynnosci-sesji__skrot';
    skrot.textContent = pozycja.skrot;
    wiersz.append(skrot);
  }

  wiersz.addEventListener('click', () => pozycja.wykonaj());

  return wiersz;
}
