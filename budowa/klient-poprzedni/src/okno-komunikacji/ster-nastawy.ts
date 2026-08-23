import type { NazwaIkony } from '../ikony/ikony';
import {
  utworzMenuDrzewo,
  type PozycjaMenu,
  type StopkaMenu,
} from '../komponenty/menu-drzewo';

/**
 * Ster nastawy — obudowa jednego komponentu paska zlecenia.
 *
 * Nie jest drugim mechanizmem menu: rozwijanie, haczyk, opisy, gałęzie, grupy,
 * pole szukania i stopka należą do `komponenty/menu-drzewo.ts`, a ten plik
 * bierze ten mechanizm gotowy.
 *
 * Istnieje, bo mechanizm drzewa nie wykonuje wyboru — oddaje klucz wołającemu
 * i na tym kończy. Każdy ster paska musi zaś zrobić z tym kluczem to samo:
 * pójść do rdzenia, poczekać na potwierdzenie i pokazać odmowę dosłownie.
 * Bez tej obudowy ten sam kawałek — czyszczenie zdania odmowy, `aria-busy` na
 * czas wysyłki, wyświetlenie błędu i powrót do stanu potwierdzonego — stałby
 * osobno w każdym sterze.
 *
 * Wyróżnienie idzie wyłącznie z migawki stanu: ster nie zapisuje wyboru
 * u siebie, tylko po wysyłce woła `odswiez` wołającego, a ten czyta stan
 * potwierdzony przez rdzeń. Nieudana zmiana nie zostawia więc mylącej etykiety
 * na uchwycie.
 *
 * Ster nie traci klikalności ani na czas wysyłki, ani po odmowie — `aria-busy`
 * mówi o pracy, nie odbiera możliwości działania.
 */
export interface OpcjeSteru {
  /**
   * Nazwa rodzajowa nastawy — „Model", „Wysiłek", „Katalog roboczy".
   *
   * Nie trafia na ekran: na uchwycie stoi wartość. Idzie do `aria-label`
   * uchwytu, bo czytnik ekranu musi wiedzieć, czego dotyczy wartość, której
   * nazwa sama tego nie mówi.
   */
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
 * Waga zdania stojącego pod uchwytem. `odmowa` znaczy „coś się nie udało"
 * i wygląda na błąd; `spokojne` znaczy „praca się odbyła, a wynik jest taki" —
 * nagranie bez mowy albo wykaz z nazwami zastępczymi. Jeden wspólny ton
 * zamieniłby poprawny wynik w fałszywy alarm.
 */
export type WagaZdania = 'odmowa' | 'spokojne';

export interface SterNastawy {
  /** Element montowany w pasku zlecenia. */
  element: HTMLElement;
  /** Podaje wartość na uchwyt i całe drzewo naraz. */
  ustaw(wartosc: string, drzewo: readonly PozycjaMenu[]): void;
  /**
   * Zdanie pod uchwytem, którego źródłem nie jest wybór w menu. Mikrofon
   * melduje wynik nagrania, którego nikt z menu nie zamawiał, a katalog
   * rozszerzeń bywa nieodczytany, zanim padnie pierwsze kliknięcie. Bez tej
   * drogi obie sytuacje musiałyby zbudować własny wiersz zdania obok
   * istniejącego.
   *
   * Treść pusta chowa zdanie; `hidden` zdejmuje je z układu, więc ster bez
   * zdania nie zostawia pustego pasa w rzędzie.
   */
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
  /**
   * Zdejmuje własne nasłuchy steru; pominięte znaczy „ten ster ich nie ma".
   *
   * Większość sterów paska nie subskrybuje niczego samodzielnie — czytają
   * jedną migawkę, a jedyną subskrypcję trzyma pasek. Mikrofon i rozszerzenia
   * mają własne: sprzęt dźwiękowy i `extension.changed` żyją poza oknem i bez
   * zdjęcia przeżyłyby je, czyli wyciekły.
   */
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
