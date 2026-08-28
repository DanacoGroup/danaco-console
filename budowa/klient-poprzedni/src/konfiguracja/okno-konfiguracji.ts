import { oznaczOknoAplikacji } from '../komponenty/okno-aplikacji';
import './konfiguracja.css';

import { elementIkony } from '../ikony/ikony';
import { otworzOknoPunktowIzolacji } from '../punkty-izolacji/indeks';
import type { Kanal } from '../protokol/kanal';
import { utworzRejestrKanalow } from '../sterowanie/rejestr-kanalow';
import { utworzPanelKanalow, type PanelKanalow } from '../sterowanie/panel-kanalow';
import { utworzNawigacjeKategorii } from './nawigacja-kategorii';
import { utworzPanelKategorii } from './panel-kategorii';
import { utworzPanelObszarowSesji, type PanelObszarowSesji } from './panel-obszarow-sesji';
import { utworzPanelZaczepow, type PanelZaczepow } from './panel-zaczepow';
import { utworzStanObszarowSesji, type StanObszarowSesji } from './stan-obszarow-sesji';
import { utworzPasekPunktuWidzenia } from './pasek-punktu-widzenia';
import { utworzStanKonfiguracji } from './stan-konfiguracji';
import { utworzStanyOdczytu } from './stany-odczytu';

/**
 * Okno konfiguracji — szkielet: pasek punktu widzenia, kolumna kategorii,
 * formularz kategorii. Stoi na natywnym `<dialog>`, otwiera się przed
 * odpowiedzią rdzenia; katalog, który nie dotarł, zostawia komunikat
 * w formularzu, okno pozostaje czynne.
 */
export interface OknoKonfiguracji {
  /** Element `<dialog>` osadzony w dokumencie. */
  element: HTMLDialogElement;
  /** Otwiera okno i wczytuje katalog wraz z wpisami. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje kanału i usuwa okno z dokumentu. */
  rozlacz(): void;
}

/** Pas rejestrów okna konfiguracji: obszary sesji i kanały modelu, obie strony pobierane raz i współdzielone przez panele okna. */
interface Rejestry {
  element: HTMLElement;
  obszary: StanObszarowSesji;
  panelObszarow: PanelObszarowSesji;
  panelZaczepow: PanelZaczepow;
  panelKanalow: PanelKanalow;
}

/**
 * Składa pas dwóch rejestrów globalnych: obszarów sesji, rozstrzyganych przez
 * rdzeń, oraz kanałów modelu, utrzymywanych jako czytająca pamięć podręczna
 * nad `channel.list` i odświeżanych po każdym udanym zapisie.
 */
function zlozRejestry(kanal: Kanal): Rejestry {
  const obszary = utworzStanObszarowSesji(kanal);
  const panelObszarow = utworzPanelObszarowSesji(obszary, () => void obszary.odswiez());
  // Zaczepy stoją na tym samym stanie obszarów — jeden punkt widzenia, jeden adres zapisu.
  const panelZaczepow = utworzPanelZaczepow(obszary);
  const rejestrKanalow = utworzRejestrKanalow(kanal);
  const panelKanalow = utworzPanelKanalow({ kanal, rejestrKanalow });

  const element = document.createElement('div');
  element.className = 'dk-okno__rejestry';
  element.append(panelObszarow.element, panelZaczepow.element, panelKanalow.element);

  return { element, obszary, panelObszarow, panelZaczepow, panelKanalow };
}

export function utworzOknoKonfiguracji(kanal: Kanal): OknoKonfiguracji {
  const stan = utworzStanKonfiguracji(kanal);
  const punkt = utworzPasekPunktuWidzenia();
  const kategorie = utworzNawigacjeKategorii();
  // Pozycja „Izolacja" prowadzi do jedynego okna punktów izolacji, otwierając ten sam egzemplarz.
  const panel = utworzPanelKategorii(stan, () => void otworzOknoPunktowIzolacji(kanal));
  const stany = utworzStanyOdczytu(stan, () => wczytaj());

  const { obszary, panelObszarow, panelZaczepow, panelKanalow, element: rejestry } =
    zlozRejestry(kanal);

  const element = oznaczOknoAplikacji({
    element: document.createElement('dialog'),
    kod: 'okno-konfiguracji',
    nazwa: 'Okno konfiguracji',
  });
  element.className = 'dn-modal dk-okno';

  const cialo = document.createElement('div');
  cialo.className = 'dk-okno__cialo';
  cialo.append(kategorie.element, panel.element);

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka dk-okno__stopka';

  const odswiez = przycisk('Odczytaj katalog ponownie', 'dn-btn dn-btn--zarys');
  const zamknij = przycisk('Zamknij', 'dn-btn dn-btn--atrament');
  stopka.append(odswiez, zamknij);

  element.append(
    naglowek(() => element.close()),
    punkt.element,
    cialo,
    rejestry,
    stany.element,
    stopka,
  );

  kategorie.naWybor((kategoria) => panel.pokaz(kategoria));

  // Jeden pasek zasięgu na okno: punkt widzenia rządzi katalogiem kluczy i obszarami sesji.
  punkt.naZmiane((wybrany) => {
    stan.ustawPunkt(wybrany);
    obszary.ustawPunkt(wybrany);
    void obszary.odswiez();
  });

  obszary.naZmiane(() => {
    panelObszarow.odswiez();
    panelZaczepow.odswiez();
  });

  /** Zmiana stanu nanosi wartości na pola zbudowane; nie przebudowuje formularza, zmienia tylko wartość. */
  stan.naZmiane(() => {
    stany.odswiez();
    panel.odswiez();
  });

  /** Pierwsze wczytanie i każde ponowne idą tą samą drogą: katalog, kategorie, formularz czynny. */
  function wczytaj(): void {
    void stan.odswiez().then(() => {
      kategorie.odswiez(stan.kategorie());
      panel.pokaz(kategorie.wybrana());
    });
  }

  odswiez.addEventListener('click', () => {
    wczytaj();
    void obszary.odswiez();
    panelKanalow.odswiez();
  });
  zamknij.addEventListener('click', () => element.close());

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      wczytaj();
      void obszary.odswiez();
      panelKanalow.odswiez();
    },

    zamknij: () => element.close(),

    rozlacz() {
      stan.rozlacz();
      element.remove();
    },
  };
}

/** Nagłówek okna konfiguracji: ikona modułu, tytuł okna oraz przycisk zamknięcia osadzony po prawej stronie paska. */
function naglowek(naZamkniecie: () => void): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-modal-naglowek dk-okno__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Konfiguracja';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dk-okno__rozpychacz';

  const zamkniecie = document.createElement('button');
  zamkniecie.type = 'button';
  zamkniecie.className = 'dn-btn-ikona';
  zamkniecie.setAttribute('aria-label', 'Zamknij okno konfiguracji');
  zamkniecie.title = 'Zamknij okno konfiguracji';
  zamkniecie.append(elementIkony('zamknij', { rozmiar: 18 }));
  zamkniecie.addEventListener('click', naZamkniecie);

  element.append(elementIkony('ustawienia', { rozmiar: 20 }), tytul, rozpychacz, zamkniecie);
  return element;
}

function przycisk(tresc: string, klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = tresc;
  return element;
}
