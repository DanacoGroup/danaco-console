import './design.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzOknoAssetsPanel, type OknoAssetsPanel } from './okno-assets-panel';
import { utworzOknoDesignBoard, type OknoDesignBoard } from './okno-design-board';
import {
  utworzOknoPreviewWindow,
  type OknoPreviewWindow,
} from '../studio/okno-preview-window';
import { utworzOknoPromptBuilder, type OknoPromptBuilder } from './okno-prompt-builder';
import {
  utworzOknoTokensSystemPanel,
  type OknoTokensSystemPanel,
} from './okno-tokens-system-panel';
import { utworzPasekUczciwosci, type PasekUczciwosci } from './pasek-uczciwosci';
import { utworzWarsztatyDesignu, type WarsztatyDesignu } from './okna-warsztatow-designu';
import { podepnijSkroty, utworzWykazSkrotow } from './skroty-designu';
import { utworzStanDesignu, type StanDesignu } from './stan-designu';
import { nasluchujWejsciaRozmowy } from './wejscie-rozmowy';
import {
  utworzWyszukiwarkeFunkcji,
  type WyszukiwarkaFunkcji,
} from './wyszukiwarka-funkcji-designu';

/**
 * Moduł Design — pięć okien operacyjnych osadzonych w jednym układzie, złożonych wedle warstwy
 * widoczności i roli okna, z jednym wspólnym stanem zasobów.
 */
export interface ModulDesign {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Ustala okno modułu, czyta zaplecze i zleca odczyt zasobów. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzModulDesign(kanal: Kanal): ModulDesign {
  const stan: StanDesignu = utworzStanDesignu(kanal);

  const plansza: OknoDesignBoard = utworzOknoDesignBoard(stan);
  const zasoby: OknoAssetsPanel = utworzOknoAssetsPanel(stan, (zasob) =>
    plansza.przyjmijZasob(zasob),
  );
  const podglad: OknoPreviewWindow = utworzOknoPreviewWindow({ rodzaj: 'design', stan });
  const kreator: OknoPromptBuilder = utworzOknoPromptBuilder(stan);

  const zetony: OknoTokensSystemPanel = utworzOknoTokensSystemPanel(stan);
  // Pięć warsztatów: fotografia, wektor, druk, bazy zdjęciowe, publikacja — 75 komend obszaru.
  const warsztaty: WarsztatyDesignu = utworzWarsztatyDesignu(stan, kanal);
  const wyszukiwarka: WyszukiwarkaFunkcji = utworzWyszukiwarkeFunkcji();
  const uczciwosc: PasekUczciwosci = utworzPasekUczciwosci(kanal);

  const pasDrugi = document.createElement('div');
  pasDrugi.className = 'md-modul__pas md-modul__pas--para';
  pasDrugi.append(zasoby.element, podglad.element, kreator.element);

  // Warsztaty stoją w osobnym pasie pod oknami warstwy pierwszej: praca nad materiałem, nie koncepcją.
  const pasWarsztatow = document.createElement('div');
  pasWarsztatow.className = 'md-modul__pas md-modul__pas--warsztaty';
  pasWarsztatow.append(...warsztaty.okna.map((okno) => okno.element));

  const pasTrzeci = document.createElement('div');
  pasTrzeci.className = 'md-modul__pas md-modul__pas--rozwiniecia';
  pasTrzeci.append(zetony.element, wyszukiwarka.element, utworzWykazSkrotow());

  const element = document.createElement('div');
  element.className = 'md-modul';
  element.dataset['modul'] = 'design';
  element.setAttribute('aria-label', 'Moduł Design — okna operacyjne');
  // Element modułu przyjmuje ognisko, by skrót klawiszowy działał, zanim Operator w cokolwiek kliknie.
  element.tabIndex = -1;
  element.append(plansza.element, pasDrugi, pasWarsztatow, pasTrzeci, uczciwosc.element);

  const odepnijSkroty = podepnijSkroty(element, {
    otworzWyszukiwarke: () => wyszukiwarka.otworz(),
    opiszWarstwe: () => plansza.opiszWarstwe(),
    otworzZetony: () => zetony.otworz(),
  });

  function odswiez(): void {
    plansza.odswiez();
    zasoby.odswiez();
    podglad.odswiez();
    kreator.odswiez();
    for (const okno of warsztaty.okna) okno.odswiez();
  }

  const odsubskrybujStan = stan.obserwuj(odswiez);
  const odsubskrybujPostep = stan.zaplecze.naPostep((tresc) => kreator.przyjmijPostep(tresc));
  const odsubskrybujRozmowe = nasluchujWejsciaRozmowy(
    stan.zrodlo,
    () => stan.idOkna(),
    (skutek, zasob) => plansza.przyjmijZRozmowy(skutek, zasob),
  );

  odswiez();

  return {
    element,

    async wczytaj(idSesji) {
      // Okno modułu musi być znane przed odczytem zasobów, bo odczyt dotyczyłby innego okna.
      await Promise.all([
        stan.ustalOkno(idSesji),
        stan.odswiezZaplecze(),
        // Katalog okien i wykaz komend idą razem z resztą: pas uczciwości ma mierzyć rdzeń, nie stać.
        uczciwosc.odczytaj(),
      ]);
      await zasoby.wczytaj();
      // Warsztaty czytają ten sam magazyn materiału, dopiero po ustaleniu okna modułu.
      await warsztaty.wczytaj();
    },

    rozlacz() {
      odsubskrybujStan();
      odsubskrybujPostep();
      odsubskrybujRozmowe();
      odepnijSkroty();
      // Plansza odpina subskrypcję obecności, bo kursory trafiałyby do panelu odłączonego od dokumentu.
      plansza.rozlacz();
      zetony.zamknij();
      uczciwosc.zamknij();
      stan.rozlacz();
    },
  };
}
