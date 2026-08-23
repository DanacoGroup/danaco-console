import type { SettingOption } from '../../../shared/contract';
import {
  lista,
  napis,
  utworzZbiornikZamiarow,
  type Kontrolka,
  type ZaleznosciKontrolki,
} from './kontrolka';

/**
 * Kontrolki wyboru: przełącznik wartości logicznej, lista jednokrotna
 * i lista wielokrotna.
 *
 * Dopuszczalne wartości pochodzą wyłącznie z pola `options` katalogu.
 * Definicja bez opcji nie daje pustego ekranu bez wyjaśnienia: obie listy
 * niosą wtedy ostrzeżenie przy polu, a lista jednokrotna zostaje przy samej
 * pozycji pustej.
 */

/** Przełącznik wartości logicznej — stan, nie działanie. */
export function utworzKontrolkePrzelacznika(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const zbiornik = utworzZbiornikZamiarow();

  const pole = document.createElement('input');
  pole.id = zaleznosci.identyfikator;
  pole.type = 'checkbox';
  pole.className = 'dn-przelacznik';
  pole.addEventListener('change', () => zbiornik.zglos());

  const koszyk = document.createElement('label');
  koszyk.className = 'dn-wybor dk-kontrolka__przelacznik';
  koszyk.htmlFor = zaleznosci.identyfikator;

  const opis = document.createElement('span');
  opis.className = 'dk-kontrolka__stan';
  opis.textContent = 'nie';

  pole.addEventListener('change', () => {
    opis.textContent = pole.checked ? 'tak' : 'nie';
  });

  koszyk.append(pole, opis);

  return {
    element: koszyk,
    odczytaj: () => pole.checked,
    ustaw: (wartosc) => {
      pole.checked = prawda(wartosc);
      opis.textContent = pole.checked ? 'tak' : 'nie';
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
  };
}

/** Lista jednokrotna — jedna wartość z katalogu opcji. */
export function utworzKontrolkeListy(zaleznosci: ZaleznosciKontrolki): Kontrolka {
  const { definicja, identyfikator } = zaleznosci;
  const opcje = uporzadkujOpcje(definicja.options);
  const zbiornik = utworzZbiornikZamiarow();

  const pole = document.createElement('select');
  pole.id = identyfikator;
  pole.className = 'dn-pole-kontrolka dk-kontrolka--lista';
  pole.addEventListener('change', () => zbiornik.zglos());

  const pusta = document.createElement('option');
  pusta.value = '';
  pusta.textContent = definicja.required ? '— wybierz —' : '— bez wartości —';
  pole.append(pusta);

  for (const opcja of opcje) pole.append(elementOpcji(opcja));

  return {
    element: pole,
    odczytaj: () => (pole.value === '' ? null : pole.value),
    ustaw: (wartosc) => {
      pole.value = napis(wartosc);
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
    ostrzezenie: opcje.length === 0 ? OSTRZEZENIE_BEZ_OPCJI : undefined,
  };
}

/** Lista wielokrotna — podzbiór wartości z katalogu opcji. */
export function utworzKontrolkeListyWielokrotnej(
  zaleznosci: ZaleznosciKontrolki,
): Kontrolka {
  const { definicja, identyfikator } = zaleznosci;
  const opcje = uporzadkujOpcje(definicja.options);
  const zbiornik = utworzZbiornikZamiarow();

  const koszyk = document.createElement('div');
  koszyk.className = 'dk-kontrolka__wielokrotna';
  koszyk.id = identyfikator;
  koszyk.setAttribute('role', 'group');

  const pola: HTMLInputElement[] = [];

  for (const opcja of opcje) {
    const etykieta = document.createElement('label');
    etykieta.className = 'dn-wybor';

    const pole = document.createElement('input');
    pole.type = 'checkbox';
    pole.className = 'dn-check';
    pole.value = opcja.value;
    pole.addEventListener('change', () => zbiornik.zglos());
    pola.push(pole);

    const nazwa = document.createElement('span');
    nazwa.textContent = opcja.label;
    if (opcja.description !== undefined) etykieta.title = opcja.description;

    etykieta.append(pole, nazwa);
    koszyk.append(etykieta);
  }

  return {
    element: koszyk,
    odczytaj: () => pola.filter((pole) => pole.checked).map((pole) => pole.value),
    ustaw: (wartosc) => {
      const wybrane = new Set(lista(wartosc));
      for (const pole of pola) pole.checked = wybrane.has(pole.value);
    },
    naZatwierdzenie: zbiornik.naZatwierdzenie,
    ostrzezenie: opcje.length === 0 ? OSTRZEZENIE_BEZ_OPCJI : undefined,
  };
}

const OSTRZEZENIE_BEZ_OPCJI =
  'Katalog nie podał dopuszczalnych wartości tej pozycji — lista pozostaje pusta.';

/** Opcje w kolejności z katalogu; wiersz bez kolejności trafia na koniec. */
function uporzadkujOpcje(opcje: readonly SettingOption[] | undefined): SettingOption[] {
  return [...(opcje ?? [])].sort((pierwsza, druga) => miejsce(pierwsza) - miejsce(druga));
}

function miejsce(opcja: SettingOption): number {
  return Number.isFinite(opcja.order) ? opcja.order : Number.MAX_SAFE_INTEGER;
}

function elementOpcji(opcja: SettingOption): HTMLOptionElement {
  const element = document.createElement('option');
  element.value = opcja.value;
  element.textContent = opcja.label;
  if (opcja.description !== undefined) element.title = opcja.description;
  return element;
}

/** Wartość logiczna odczytana z kształtu nieznanego; napis „true" też liczy się jako prawda. */
function prawda(wartosc: unknown): boolean {
  if (typeof wartosc === 'boolean') return wartosc;
  if (typeof wartosc === 'number') return wartosc !== 0;
  if (typeof wartosc === 'string') return wartosc === 'true' || wartosc === '1';
  return false;
}
