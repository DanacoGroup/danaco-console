import './strona-glowna.css';
import './karty-strony-glownej.css';
import './sesje-strony-glownej.css';

import { utworzListweUstawien, type ListwaUstawien } from './listwa-ustawien';
import { utworzNaglowekStrony } from './naglowek-strony';
import type { PozycjaKomponentu } from './pozycje-komponentow';
import type { KodSrodowiska, PozycjaSrodowiska } from './pozycje-srodowisk';
import type { PozycjaUstawienia } from './pozycje-ustawien';
import { utworzStrefeKomponentow, type StrefaKomponentow } from './strefa-komponentow';
import { utworzStrefeModulow, type StrefaModulow } from './strefa-modulow';
import { utworzStrefeSesji, type StrefaSesji } from './strefa-sesji';
import { utworzStrefeSrodowisk, type StrefaSrodowisk } from './strefa-srodowisk';
import type { SluchaczWyboru } from './sygnal-wyboru';
import { utworzWykazSrodowisk, type WykazSrodowisk } from './wykaz-srodowisk';

/**
 * Strona główna — centrum dowodzenia.
 *
 * Jedna odpowiedzialność: złożenie stref w jeden widok i wystawienie ich
 * zdarzeń na zewnątrz. Plik nie buduje żadnej karty ani kafla — jest
 * kompozycją, nie logiką widoku.
 *
 * Strefy stoją w kolejności malejącej masy wizualnej:
 *   karty środowisk    — środek ciężkości, krój nagłówkowy 24/30
 *   kafle komponentów  — ta sama karta, krój bazowy półgruby
 *   kafle modułów spoza nawigacji — ta sama forma; strefa ukryta, dopóki rdzeń
 *                        nie poda modułu stojącego poza nawigacją
 *   listwa ustawień    — masa najniższa, tło `--dn-powierzchnia-2`
 *   sesje w tle        — wykaz z danych rdzenia, zasilany z zewnątrz przez
 *                        `wpiecie-sesji`; sama strona danych nie pobiera
 *
 * Strona nie otwiera środowiska samodzielnie. Zgłasza wybór, a skutek należy do
 * warstwy, która stronę zamontowała — przejście przez stronę główną ma być
 * świadome i widoczne.
 */
export interface StronaGlowna {
  /** Element widoku gotowy do osadzenia w powłoce aplikacji. */
  element: HTMLElement;
  /** Zdarzenie wyboru środowiska — strefa kart środowisk. */
  naWyborSrodowiska(sluchacz: SluchaczWyboru<PozycjaSrodowiska>): void;
  /** Zdarzenie wyboru komponentu własnego — strefa kafli komponentów. */
  naWyborKomponentu(sluchacz: SluchaczWyboru<PozycjaKomponentu>): void;
  /** Zdarzenie wyboru pozycji ustawień — listwa ustawień. */
  naWyborUstawienia(sluchacz: SluchaczWyboru<PozycjaUstawienia>): void;
  /** Oznacza środowisko czynne; `null` zdejmuje oznaczenie. */
  ustawSrodowiskoCzynne(kod: KodSrodowiska | null): void;
  /** Strefa sesji w tle — zasilana danymi rdzenia przez `wpiecie-sesji`. */
  sesje: StrefaSesji;
  /** Strefa środowisk — zasilana wykazem rdzenia przez `wpiecie-srodowisk`. */
  srodowiska: StrefaSrodowisk;
  /**
   * Wykaz środowisk strony — jedno źródło nazwy środowiska na tym ekranie.
   *
   * Aktualizuje się przy każdym `srodowiska.ustawWykaz`, którąkolwiek drogą
   * wykaz przyszedł: `home.enter` albo `environment.list`. Zasilenia kart bez
   * zasilenia wykazu nie da się tu napisać.
   */
  wykazSrodowisk: WykazSrodowisk;
  /**
   * Strefa komponentów własnych — czwórka rodzajów stoi od razu z kontraktu,
   * kafle personalizowane dokłada `wpiecie-komponentow` z `component.list`.
   */
  komponenty: StrefaKomponentow;
  /**
   * Kafle modułów, których rdzeń nie pokazuje w bocznej nawigacji żadnego
   * środowiska — jedyna droga do ich okien. Wykaz dokłada `wpiecie-modulow`
   * z `module.list`; bez niego strefa zostaje ukryta.
   */
  moduly: StrefaModulow;
  /** Listwa ustawień — niesie przybornik i segment „Dodaj nowy". */
  ustawienia: ListwaUstawien;
}

export function utworzStroneGlowna(): StronaGlowna {
  const wykazSrodowisk = utworzWykazSrodowisk();

  // Karty środowisk i wiersze sesji w tle czytają tę samą listę, bo
  // `ustawWykaz` jest opakowane: nie ma drogi, którą karty dostałyby nowe nazwy,
  // a wiersze sesji zostały przy starych. Wpięcia wołają dalej
  // `strona.srodowiska.ustawWykaz` i o wykazie nie muszą wiedzieć.
  const strefaSrodowisk = utworzStrefeSrodowisk();
  const srodowiska: StrefaSrodowisk = {
    ...strefaSrodowisk,
    ustawWykaz(pozycje) {
      wykazSrodowisk.ustaw(pozycje);
      strefaSrodowisk.ustawWykaz(pozycje);
    },
  };

  const komponenty = utworzStrefeKomponentow();
  const moduly = utworzStrefeModulow();
  const sesje = utworzStrefeSesji(wykazSrodowisk);
  const ustawienia = utworzListweUstawien();

  const element = document.createElement('main');
  element.className = 'dn-strona';

  const tresc = document.createElement('div');
  tresc.className = 'dn-strona__tresc';
  // Linie kafli idą razem, sesje w tle za nimi: strefy działania mają być
  // widoczne naraz, bez przewijania, a sesje w tle są wglądem w to, co już
  // biegnie, i należą do drugiego planu.
  tresc.append(
    utworzNaglowekStrony(),
    srodowiska.element,
    komponenty.element,
    moduly.element,
    ustawienia.element,
    sesje.element,
  );

  element.append(tresc);

  return {
    element,
    srodowiska,
    wykazSrodowisk,
    komponenty,
    moduly,

    naWyborSrodowiska: srodowiska.naWybor,
    naWyborKomponentu: komponenty.naWybor,
    naWyborUstawienia: ustawienia.naWybor,
    ustawienia,
    ustawSrodowiskoCzynne: srodowiska.ustawCzynne,
    sesje,
  };
}
