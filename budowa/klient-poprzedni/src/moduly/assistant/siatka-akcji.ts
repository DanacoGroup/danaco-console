import type { Action } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { BRAKI } from './braki-kontraktu';
import { KLASY_DYMKA, ODCZYTY, PUSTE } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import { KOD_MODULU, type StanAssistant } from './stan-assistant';

/**
 * Siatka szybkich akcji Voice Console. Plik odpowiada wyłącznie za widok katalogu
 * akcji zasięgu modułu, a zestaw pochodzi z komendy `action.list`, nie z listy
 * zaszytej w kliencie: wykaz pozycji należy do rdzenia.
 */
export interface SiatkaAkcji {
  element: HTMLElement;
  /** Odczyt katalogu akcji modułu z rdzenia. */
  wczytaj(): Promise<void>;
}

/**
 * Wstawia treść akcji do pola polecenia. Akcja jest zaczynem polecenia, a nie jego
 * wysyłką, więc naciśnięcie kafla niczego nie wykonuje; polecenie wychodzi dopiero
 * z paska, po przeczytaniu przez Operatora.
 */
export type NaAkcje = (tresc: string) => void;

export function utworzSiatkeAkcji(stan: StanAssistant, naAkcje: NaAkcje): SiatkaAkcji {
  const okno: StanOkna = utworzStanOkna();

  const lista = document.createElement('ul');
  lista.className = 'ma-akcje';

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis ma-brak';
  granica.textContent = BRAKI.bibliotekaSkrotow;

  okno.tresc.append(lista, granica);
  okno.puste(PUSTE.akcjeSpoczynek);

  const tytul = document.createElement('h4');
  tytul.className = 'ma-panel__tytul';
  tytul.append(
    'Szybkie akcje',
    utworzDymekObjasnienia(
      'Katalog akcji zasięgu modułu z rdzenia (action.list). Naciśnięcie ' +
        'wstawia treść akcji do pola polecenia — wysyłkę rozstrzyga Operator.',
      KLASY_DYMKA,
    ),
  );

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'szybkie-akcje';
  element.append(tytul, okno.element);

  async function wczytaj(): Promise<void> {
    okno.ladowanie(ODCZYTY.akcje);
    const wynik = await stan.zaplecze.akcje(KOD_MODULU);
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt katalogu akcji', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const akcje = [...wynik.wynik.actions].sort((pierwsza, druga) => pierwsza.order - druga.order);
    lista.replaceChildren(...akcje.map((akcja) => kafel(akcja, naAkcje)));
    if (akcje.length === 0) {
      okno.puste(PUSTE.akcje);
      return;
    }
    okno.gotowe();
  }

  return { element, wczytaj };
}

/**
 * Jeden kafel katalogu; nazwa akcji jest zaczynem treści polecenia. Kafel wstawia do
 * pola nazwę wiersza, a nie jego opis, ponieważ opis jest zdaniem o komendzie rdzenia,
 * nie treścią polecenia dla asystenta.
 */
function kafel(akcja: Action, naAkcje: NaAkcje): HTMLElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-kafel ma-akcje__kafel';
  przycisk.dataset['akcja'] = akcja.id;

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-kafel-etykieta';
  nazwa.textContent = akcja.name;

  const opis = document.createElement('span');
  opis.className = 'dn-kafel-opis';
  opis.textContent = opisWiersza(akcja);

  przycisk.append(nazwa, opis);
  przycisk.addEventListener('click', () => {
    // Wiersz wskazuje komendę kontraktu; wysyłanie jej wprost byłoby zgadywaniem.
    naAkcje(akcja.name);
  });

  const pozycja = document.createElement('li');
  pozycja.append(przycisk);
  return pozycja;
}

/**
 * Opis wiersza katalogu: czym wiersz jest, jaką komendę rdzenia wskazuje i jakiego bytu
 * wymaga. Warunek z pola `requires` stoi w opisie kafla, a nie w oknie po naciśnięciu,
 * ponieważ kafel niczego nie wykonuje i sam wypełnia pole.
 */
function opisWiersza(akcja: Action): string {
  const czesci = [akcja.description ?? '', `komenda rdzenia: ${akcja.command}`];
  if (akcja.requires !== undefined && akcja.requires !== '') {
    czesci.push(`wymaga bytu: ${akcja.requires}`);
  }
  return czesci.filter((czesc) => czesc !== '').join(' · ');
}
