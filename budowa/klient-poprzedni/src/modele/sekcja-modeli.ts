import './modele.css';

import type { Kanal } from '../protokol/kanal';
import { utworzPanelKont } from './panel-kont';
import { utworzPanelTozsamosci } from './panel-tozsamosci';
import { utworzPodgladPromptu } from './podglad-promptu';
import { utworzStanKont } from './stan-kont';
import { utworzStanyOdczytu } from './stany-odczytu';
import { utworzStanTozsamosci } from './stan-tozsamosci';
import { utworzUstawieniaBytu } from './ustawienia-bytu';
import { utworzWyborOsi, type PodpowiedzKonta } from './wybor-osi';
import { utworzZakladkiSekcji } from './zakladki-sekcji';
import { utworzZrodloModeli } from './zrodlo-modeli';

/**
 * Sekcja modeli, kont i tożsamości — cztery obszary związane jedną osią.
 *
 * Oś jest wspólna dla całej sekcji. Ustawienia, tożsamość i podgląd promptu
 * mówią o tym samym bycie: wskazanym modelu albo wskazanym koncie. Trzy
 * osobne wybory osi w trzech panelach dałyby trzy rozbieżne wskazania,
 * dlatego pasek osi jest jeden, nad zakładkami.
 *
 * Wykaz kont zasila wybór osi: konta znane rdzeniowi stają się podpowiedziami
 * bytu osi `account`, a identyfikatory modeli składa osobne źródło z rejestru
 * kanałów i z modeli domyślnych kont. Podpowiedź nie zamyka pola — byt spoza
 * wykazu wolno wpisać wprost.
 *
 * Sekcja otwiera się przed odpowiedzią rdzenia. Każdy z czterech obszarów ma
 * własny komunikat na wypadek braku danych i żaden nie blokuje pozostałych.
 */
export interface SekcjaModeli {
  /** Sekcja osadzana w oknie albo w widoku gospodarza. */
  element: HTMLElement;
  /** Wczytuje rejestr kont, katalogi i nakładkę obowiązującą. */
  odswiez(): void;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

export function utworzSekcjeModeli(kanal: Kanal): SekcjaModeli {
  const stanKont = utworzStanKont(kanal);
  const stanTozsamosci = utworzStanTozsamosci(kanal);
  const zrodloModeli = utworzZrodloModeli(kanal);

  const os = utworzWyborOsi();
  const panelKont = utworzPanelKont(stanKont);
  const ustawienia = utworzUstawieniaBytu(kanal);
  const panelTozsamosci = utworzPanelTozsamosci(stanTozsamosci);
  const podglad = utworzPodgladPromptu(stanTozsamosci);

  const zakladki = utworzZakladkiSekcji([
    { kod: 'konta', nazwa: 'Konta modeli i code CLI', element: panelKont.element },
    { kod: 'ustawienia', nazwa: 'Ustawienia osi', element: ustawienia.element },
    { kod: 'tozsamosc', nazwa: 'Zasady i tożsamość', element: panelTozsamosci.element },
    { kod: 'prompt', nazwa: 'Podgląd złożonego promptu', element: podglad.element },
  ]);

  const stany = utworzStanyOdczytu(stanKont, () => void stanKont.odswiez());

  const pasekOsi = document.createElement('div');
  pasekOsi.className = 'dm-sekcja__os';
  pasekOsi.append(os.element);

  const element = document.createElement('section');
  element.className = 'dm-sekcja';
  element.append(pasekOsi, stany.element, zakladki.element);

  /** Pasek osi nie dotyczy rejestru kont — rejestr jest wspólny dla osi. */
  const ubierzPasek = (): void => {
    pasekOsi.hidden = zakladki.czynna() === 'konta';
  };

  const rozeslijOs = (): void => {
    const wskazanie = os.wskazanie();
    ustawienia.ustawWskazanie(wskazanie);
    void stanTozsamosci.ustawWskazanie(wskazanie);
  };

  /**
   * Identyfikatory modeli trzymamy między odczytami, bo pochodzą z rejestru
   * kanałów, a ten czyta się osobną komendą. Zmiana rejestru kont odświeża
   * podpowiedzi kont i zostawia modele takimi, jakie były — podanie pustej
   * listy skasowałoby podpowiedź bez powodu.
   */
  let modele: readonly string[] = [];

  const podpowiedziKont = (): PodpowiedzKonta[] =>
    stanKont.konta().map((konto) => ({ identyfikator: konto.id, nazwa: konto.name }));

  const naniesPodpowiedzi = (): void => os.ustawPodpowiedzi(modele, podpowiedziKont());

  /** Podpowiedzi modeli składamy z rejestru kanałów i z modeli domyślnych kont. */
  const odswiezPodpowiedzi = (): void => {
    void zrodloModeli.identyfikatory(stanKont.konta()).then((odczytane) => {
      modele = odczytane;
      naniesPodpowiedzi();
    });
  };

  stanKont.naZmiane(() => {
    stany.odswiez();
    panelKont.odswiez();
    naniesPodpowiedzi();
  });

  stanTozsamosci.naZmiane(() => {
    panelTozsamosci.odswiez();
    podglad.odswiez();
  });

  os.naZmiane(rozeslijOs);
  zakladki.naZmiane(ubierzPasek);

  ubierzPasek();
  panelKont.odswiez();
  panelTozsamosci.odswiez();
  podglad.odswiez();

  return {
    element,

    /**
     * Odczyt nie rozsyła osi ponownie. Oś zmienia wyłącznie pasek u góry,
     * a jego zmiana już ją rozesłała; powtórzenie tutaj kazałoby rdzeniowi
     * przeczytać te same zapisy i tę samą nakładkę drugi raz pod rząd.
     */
    odswiez() {
      void stanKont.odswiez().then(odswiezPodpowiedzi);
      void stanTozsamosci.odswiez();
      ustawienia.odswiez();
    },

    rozlacz() {
      stanKont.rozlacz();
      stanTozsamosci.rozlacz();
      ustawienia.rozlacz();
    },
  };
}
