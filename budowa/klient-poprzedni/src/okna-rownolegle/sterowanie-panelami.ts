import { MENU_ETYKIETA } from './etykiety-paneli';
import { ikonaPanelu } from './ikony-paneli';
import { utworzMenuPaneli, type PozycjaMenu } from './menu-paneli';
import { utworzMenuRozwijane } from './menu-rozwijane';
import { utworzSkrotyPaneli } from './skroty-paneli';

/**
 * Komplet sterowania panelami osadzany w nagłówku okna rozmowy składa rząd skrótów i przycisk menu w jeden byt, przebudowując wnętrze przy zmianie modułu, bez znajomości stanu paneli i z obowiązkowym zamknięciem nasłuchu menu.
 */
export interface SterowaniePanelami {
  /** Element osadzany w nagłówku okna rozmowy; stały przez życie gniazda. */
  element: HTMLElement;
  /** Przestawia wykaz pozycji — po zmianie modułu gniazda. */
  ustawPozycje(pozycje: readonly PozycjaMenu[], nieotwieralne: number): void;
  /** Odświeża znaczniki „otwarty/zamknięty" w menu i na skrótach. */
  odswiez(): void;
  /** Zdejmuje nasłuchy dokumentu założone przez menu rozwijane. */
  zamknij(): void;
}

export interface OpcjeSterowania {
  pozycje: readonly PozycjaMenu[];
  /** Ile pozycji spisu nie da się otworzyć — do zdania pod wykazem. */
  nieotwieralne: number;
  czyOtwarty(kod: string): boolean;
  naWybor(kod: string): void;
  /** Sekcje doklejane w menu za kreską przechodzą nietknięte, bez własności sterowania panelami. */
  sekcjeDalsze?: readonly HTMLElement[];
}

export function utworzSterowaniePanelami(opcje: OpcjeSterowania): SterowaniePanelami {
  const element = document.createElement('div');
  element.className = 'dn-okna__sterowanie-paneli';

  /** Byty wnętrza; przebudowywane razem przy każdej zmianie wykazu. */
  let skroty = zbudujSkroty(opcje.pozycje, opcje);
  let menu = zbudujMenu(opcje.pozycje, opcje.nieotwieralne, opcje);
  element.append(skroty.element, menu.uchwyt.element);

  return {
    element,

    ustawPozycje(pozycje, nieotwieralne) {
      // Menu poprzedniego wykazu ma własny nasłuch dokumentu — bez zamknięcia przeżyłoby podmianę wnętrza.
      menu.uchwyt.zamknij();
      menu.tresc.zamknij();

      skroty = zbudujSkroty(pozycje, opcje);
      menu = zbudujMenu(pozycje, nieotwieralne, opcje);
      element.replaceChildren(skroty.element, menu.uchwyt.element);
    },

    odswiez() {
      skroty.odswiez();
      menu.tresc.odswiez();
    },

    zamknij() {
      menu.uchwyt.zamknij();
      menu.tresc.zamknij();
    },
  };
}

/** Rząd skrótów ikonowych panelu; pusty wykaz pozycji daje rząd, który wtedy nie zajmuje żadnego miejsca. */
function zbudujSkroty(pozycje: readonly PozycjaMenu[], opcje: OpcjeSterowania) {
  return utworzSkrotyPaneli({
    pozycje,
    ikonaPozycji: ikonaPanelu,
    czyOtwarty: opcje.czyOtwarty,
    naWybor: opcje.naWybor,
  });
}

/**
 * Menu `⋮` wraz z jego treścią.
 *
 * Uchwyt (byt rozwijany) i treść (wykaz paneli) są rozdzielone: uchwyt nie wie,
 * co pokazuje, a wykaz nie wie, że siedzi w menu. Dzięki temu to samo menu
 * przyjmuje czynności sesji, których ten pakiet nie buduje.
 */
function zbudujMenu(
  pozycje: readonly PozycjaMenu[],
  nieotwieralne: number,
  opcje: OpcjeSterowania,
) {
  const uchwyt = utworzMenuRozwijane({ ikona: 'wiecej', etykieta: MENU_ETYKIETA });
  const tresc = utworzMenuPaneli({
    pozycje,
    nieotwieralne,
    czyOtwarty: opcje.czyOtwarty,
    naWybor: opcje.naWybor,
    // Ten sam element wraca przy przebudowie, bo dołączenie przenosi go, nie kopiuje.
    ...(opcje.sekcjeDalsze === undefined ? {} : { sekcjeDalsze: opcje.sekcjeDalsze }),
  });
  uchwyt.ustawTresc([tresc.element]);
  return { uchwyt, tresc };
}
