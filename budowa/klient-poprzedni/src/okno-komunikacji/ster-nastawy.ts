import type { NazwaIkony } from '../ikony/ikony';
import {
  utworzMenuDrzewo,
  type PozycjaMenu,
  type StopkaMenu,
} from '../komponenty/menu-drzewo';

/**
 * Ster nastawy jest obudową jednego komponentu paska zlecenia: bierze gotowy mechanizm menu i dodaje wysyłkę do rdzenia, potwierdzenie i pokazanie odmowy dosłownie, bez parafrazy.
 */
export interface OpcjeSteru {
  /** Nazwa rodzajowa nastawy idzie do aria-label, bo czytnik ekranu musi wiedzieć, czego wartość dotyczy. */
  nastawa: string;
  /** Ikona na uchwycie przed wartością; pominięta znaczy „bez ikony". */
  ikona?: NazwaIkony;
  /** Droga do rejestru na dole menu; pominięta = menu jej nie ma. */
  stopka?: StopkaMenu;
  /** Od ilu liści menu stawia pole szukania; pominięty = próg mechanizmu. */
  progSzukania?: number;
  /**
   * Wykonanie wyboru. Odrzucenie niesie zdanie odmowy rdzenia, pokazywane
   * dosłownie, bez parafrazy.
   */
  wykonaj(klucz: string): Promise<void>;
  /** Przerysowanie steru ze stanu potwierdzonego — wołane po każdej wysyłce. */
  odswiez(): void;
}

/**
 * Waga zdania pod uchwytem: odmowa wygląda na błąd, a spokojne znaczy, że praca się odbyła i wynik jest taki, jak nagranie bez mowy albo wykaz z nazwami zastępczymi.
 */
export type WagaZdania = 'odmowa' | 'spokojne';

export interface SterNastawy {
  /** Element montowany w pasku zlecenia. */
  element: HTMLElement;
  /** Podaje wartość na uchwyt i całe drzewo naraz. */
  ustaw(wartosc: string, drzewo: readonly PozycjaMenu[]): void;
  /** Zdanie pod uchwytem bez źródła w menu: mikrofon melduje nagranie, katalog melduje odczyt rozszerzeń. */
  zdanie(tresc: string, waga?: WagaZdania): void;
}

/**
 * Jeden obsadzony komponent paska — tyle, ile pasek o nim wie.
 *
 * Pasek nie zna ani wartości, ani drzewa swoich sterów: podaje im miejsce
 * i mówi „przerysuj się ze stanu". Zawartość jest własnością steru.
 */
export interface SterPaska {
  /** Element montowany w pasku zlecenia. */
  element: HTMLElement;
  /** Przerysowuje uchwyt i drzewo ze stanu potwierdzonego przez rdzeń. */
  odswiez(): void;
  /** Większość sterów nie subskrybuje niczego samo — mikrofon i rozszerzenia mają subskrypcje do zdjęcia. */
  rozlacz?(): void;
}

export function utworzSterNastawy(opcje: OpcjeSteru): SterNastawy {
  const element = document.createElement('div');
  element.className = 'dc-ster-zlecenia';

  const odmowa = document.createElement('p');
  odmowa.className = 'dn-plakietka dc-ster-zlecenia__odmowa';
  odmowa.setAttribute('role', 'alert');
  odmowa.hidden = true;

  /** Jedyne miejsce, w którym zdanie pod uchwytem powstaje i znika. */
  function pokazZdanie(tresc: string, waga: WagaZdania): void {
    odmowa.textContent = tresc;
    odmowa.hidden = tresc === '';
    odmowa.classList.toggle('dn-plakietka--blad', waga === 'odmowa');
    odmowa.classList.toggle('dn-plakietka--informacja', waga === 'spokojne');
  }

  pokazZdanie('', 'odmowa');

  const menu = utworzMenuDrzewo({
    nastawa: opcje.nastawa,
    ...(opcje.ikona === undefined ? {} : { ikona: opcje.ikona }),
    ...(opcje.stopka === undefined ? {} : { stopka: opcje.stopka }),
    ...(opcje.progSzukania === undefined ? {} : { progSzukania: opcje.progSzukania }),
    naWybor: (klucz) => {
      void wybierz(klucz);
    },
  });

  element.append(menu.element, odmowa);

  async function wybierz(klucz: string): Promise<void> {
    pokazZdanie('', 'odmowa');
    element.setAttribute('aria-busy', 'true');
    try {
      await opcje.wykonaj(klucz);
    } catch (blad) {
      pokazZdanie(blad instanceof Error ? blad.message : String(blad), 'odmowa');
    } finally {
      element.removeAttribute('aria-busy');
      opcje.odswiez();
    }
  }

  return {
    element,
    ustaw: (wartosc, drzewo) => menu.ustaw(wartosc, drzewo),
    zdanie: (tresc, waga = 'odmowa') => pokazZdanie(tresc, waga),
  };
}
