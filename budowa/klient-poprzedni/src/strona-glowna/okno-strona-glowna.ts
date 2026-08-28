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

/** Strona główna, centrum dowodzenia, składa strefy w jeden widok i wystawia ich zdarzenia na zewnątrz, nie otwierając środowiska samodzielnie ani nie budując żadnej karty czy kafla. */
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
  /** Wykaz środowisk strony aktualizuje się przy każdym ustawieniu wykazu strefy. */
  wykazSrodowisk: WykazSrodowisk;
  /** Strefa komponentów: czwórka rodzajów stoi od razu z kontraktu, personalizowane dokłada wpięcie. */
  komponenty: StrefaKomponentow;
  /** Kafle modułów bez pozycji w bocznej nawigacji; bez wykazu strefa zostaje ukryta. */
  moduly: StrefaModulow;
  /** Listwa ustawień — niesie przybornik i segment „Dodaj nowy". */
  ustawienia: ListwaUstawien;
}

export function utworzStroneGlowna(): StronaGlowna {
  const wykazSrodowisk = utworzWykazSrodowisk();

  // Karty środowisk i wiersze sesji w tle czytają tę samą listę, bo ustawienie wykazu jest opakowane.
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
  // Linie kafli idą razem, sesje w tle za nimi: strefy działania widoczne naraz, bez przewijania.
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
