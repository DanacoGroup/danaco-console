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
 * Moduł Design — pięć okien operacyjnych osadzonych w jednym układzie.
 *
 * Jedna odpowiedzialność: złożenie okien modułu i rozdanie im jednego stanu.
 *
 * Układ idzie warstwami widoczności opracowania, nie kolejnością plików.
 * W pasie pierwszym i drugim stoją okna warstwy pierwszej — te, które są
 * widoczne bez interakcji. W pasie trzecim stoją rozwinięcia warstw wyższych:
 * Tokens & System Panel (warstwa trzecia, wywoływany menu kebab), wyszukiwarka
 * funkcji i wykaz skrótów (warstwa czwarta). Zwinięte nie znaczy ukryte —
 * zapowiedź nad każdym mówi, co jest pod spodem.
 *
 * Układ wynika z roli okna. Design Board jest wiodące i stoi w pasie pierwszym
 * na całą szerokość — na nim odbywa się praca koncepcyjna. W pasie drugim stoją
 * trzy pozostałe w kolejności katalogu rdzenia: Assets Panel (zarządca) wskazuje
 * zasób, Preview Window (pomocnicze) pokazuje zasób wskazany, Prompt Builder
 * (kreator) zleca nowy. Podgląd stoi między nimi, bo patrzy i na to, co zarządca
 * wskazał, i na to, co kreator dopiero przyniósł.
 *
 * Jeden zbiór zasobów na cały moduł: wynik generowania z kreatora wchodzi do
 * wykazu zarządcy, stamtąd na kanwę wiodącego, a podgląd czyta ten sam wybór,
 * bo `stan-designu` jest jeden. Zdarzenie `design.asset.changed` wciąga zasób
 * tą samą drogą także wtedy, gdy zlecenie przyszło z obcego połączenia.
 *
 * Druga droga na kanwę prowadzi z rozmowy. Zasób zlecony spoza okien modułu
 * wchodzi zdarzeniem `design.asset.changed`: stan wciąga go do wykazu zarządcy,
 * a to złożenie kładzie go warstwą na kanwie wiodącego (`wejscie-rozmowy.ts`).
 * Wiązanie mieszka tutaj, bo wiąże dwa okna i nie jest sprawą żadnego z nich
 * z osobna.
 *
 * Moduł nie osadza się sam — oddaje element; gdzie stanie, rozstrzyga warstwa
 * składająca.
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
  // Pięć warsztatów: fotografia, wektor, druk, bazy zdjęciowe, publikacja.
  // Komenda bez okna jest funkcją, której Operator nie ma — a te pięć okien
  // niesie siedemdziesiąt pięć komend obszaru.
  const warsztaty: WarsztatyDesignu = utworzWarsztatyDesignu(stan, kanal);
  const wyszukiwarka: WyszukiwarkaFunkcji = utworzWyszukiwarkeFunkcji();
  const uczciwosc: PasekUczciwosci = utworzPasekUczciwosci(kanal);

  const pasDrugi = document.createElement('div');
  pasDrugi.className = 'md-modul__pas md-modul__pas--para';
  pasDrugi.append(zasoby.element, podglad.element, kreator.element);

  // Warsztaty stoją w osobnym pasie, pod oknami warstwy pierwszej: są miejscem
  // pracy nad materiałem, a nie nad koncepcją, więc Operator wchodzi w nie po
  // tym, jak koncepcja stoi na kanwie.
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
  // Element modułu przyjmuje ognisko, żeby skrót klawiszowy działał także wtedy,
  // gdy Operator jeszcze w nic w module nie kliknął.
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
      // Okno modułu musi być znane przed odczytem zasobów: `design.asset.list`
      // przyjmuje `windowId`, a bez niego odczyt dotyczyłby czegoś innego niż
      // to okno. Zaplecze idzie równolegle — jest niezależne.
      await Promise.all([
        stan.ustalOkno(idSesji),
        stan.odswiezZaplecze(),
        // Katalog okien i wykaz komend idą razem z resztą: pas uczciwości ma
        // mierzyć rdzeń, a nie stać na zdaniu „odczyt w toku" do końca sesji.
        uczciwosc.odczytaj(),
      ]);
      await zasoby.wczytaj();
      // Warsztaty czytają ten sam magazyn materiału. Odczyt idzie PO ustaleniu
      // okna modułu, bo bez niego wykaz dotyczyłby czegoś innego niż to okno.
      await warsztaty.wczytaj();
    },

    rozlacz() {
      odsubskrybujStan();
      odsubskrybujPostep();
      odsubskrybujRozmowe();
      odepnijSkroty();
      // Plansza odpina subskrypcję `design.board.presence`: kursory współpracy
      // trafiałyby inaczej do panelu odłączonego już od dokumentu.
      plansza.rozlacz();
      zetony.zamknij();
      uczciwosc.zamknij();
      stan.rozlacz();
    },
  };
}
