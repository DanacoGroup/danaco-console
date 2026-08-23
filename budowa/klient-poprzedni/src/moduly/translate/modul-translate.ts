import './translate.css';

import type { Kanal } from '../../protokol/kanal';
import { utworzKatalogOkien } from '../katalog-okien';
import { KOD_MODULU, KODY_OKIEN } from './katalog-okien-translate';
import { utworzMacierzIzolacji, type MacierzIzolacji } from './macierz-izolacji';
import { utworzOknoFormatStudio, type OknoFormatStudio } from './okno-format-studio';
import { utworzOknoGlossaryManager, type OknoGlossaryManager } from './okno-glossary-manager';
import { utworzOknoQaReview, type OknoQaReview } from './okno-qa-review';
import { utworzOknoSourcePanel, type OknoSourcePanel } from './okno-source-panel';
import {
  utworzOknoWarsztatTranslate,
  type OknoWarsztatTranslate,
} from './okno-warsztat-translate';
import {
  utworzOknoTranslationMemory,
  type OknoTranslationMemory,
} from './okno-translation-memory';
import {
  utworzOknoTranslationPanels,
  type OknoTranslationPanels,
} from './okno-translation-panels';
import { utworzPasekKontekstu, type PasekKontekstu } from './pasek-kontekstu';
import { podepnijSkroty, utworzWykazSkrotow } from './skroty-translate';
import { utworzStanTranslate, type StanTranslate } from './stan-translate';
import { utworzWystawienieOperacji } from './wystawienie-operacji';
import { utworzWyszukiwarkeFunkcji, type WyszukiwarkaFunkcji } from './wyszukiwarka-funkcji';
import { utworzWyzwalaczeOkien, type WyzwalaczeOkien } from './wyzwalacze-okien';

/**
 * Moduł Translate — sześć okien operacyjnych w jednym układzie.
 *
 * Układ wynika z warstw widoczności opracowania, nie z upodobania. Warstwa
 * pierwsza stoi na ekranie bez interakcji: Source Panel i Translation Panels
 * obok siebie, bo zapis źródła aktualizuje wszystkie panele naraz i skutek
 * musi być widoczny w tej samej chwili. Format Studio jest kolumną sąsiadującą
 * warstwy drugiej, otwieraną przy pracy z dokumentem. Glossary Manager,
 * Translation Memory Panel i QA & Review Center są zarządcami warstwy trzeciej
 * — stoją w kolumnie bocznej i są zwinięte, dopóki ich nie wywołać.
 *
 * Chat Window i Execution Loop Window są oknami wspólnymi platformy i moduł ich
 * nie buduje: pierwsze montuje scena sesji, drugie należy do pętli wykonawczej.
 * Katalog okien mówi obok, ile okien modułu rdzeń zna, a ile moduł buduje —
 * liczba jest liczona z odpowiedzi rdzenia, nie wpisana.
 *
 * Cały moduł ma jeden stan. Tekst źródłowy, panele języków i identyfikator okna
 * są wspólne wszystkim oknom — gdyby każde prowadziło własną kopię, zapis
 * źródła odświeżałby jedno okno, a pozostałe zostałyby przy poprzedniej treści.
 *
 * Drogi warstwy czwartej są trzy i wszystkie prowadzą do tych samych rzeczy:
 * skrót klawiszowy, wyszukiwarka funkcji z pełnym katalogiem opracowania oraz
 * rozwinięcia przy oknach, których dotyczą. Zasada jednego kliknięcia jest przez
 * to spełniona także dla pozycji, które nie mają własnego przycisku w oknie.
 */
export interface ModulTranslate {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Odczytuje okno modułu z rdzenia dla wskazanej sesji. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia i nasłuch skrótów. */
  rozlacz(): void;
}

export function utworzModulTranslate(kanal: Kanal): ModulTranslate {
  const stan: StanTranslate = utworzStanTranslate(kanal);
  // Sesja ostatniego wejścia — ponowienie odczytu kontekstu musi wiedzieć,
  // o czyje okna pytać, a moduł nie sięga po sesję sam (dostaje ją z powłoki).
  let ostatniaSesja = '';

  const kontekst: PasekKontekstu = utworzPasekKontekstu(stan, () => {
    if (ostatniaSesja === '') return;
    void stan.wczytajKontekst(ostatniaSesja);
  });

  const zrodlo: OknoSourcePanel = utworzOknoSourcePanel(stan);
  const panele: OknoTranslationPanels = utworzOknoTranslationPanels(stan);
  const formaty: OknoFormatStudio = utworzOknoFormatStudio(stan, stan.dokument, (tekst) =>
    zrodlo.wstawZrodlo(tekst),
  );
  const glosariusz: OknoGlossaryManager = utworzOknoGlossaryManager(stan);
  const pamiec: OknoTranslationMemory = utworzOknoTranslationMemory(stan);
  const jakosc: OknoQaReview = utworzOknoQaReview(stan);
  // Warsztat jest oknem zarządcy warstwy trzeciej, tak jak trzy okna obok:
  // stoi zwinięty, dopóki Operator go nie wywoła, bo prowadzi czynności
  // wsadowe i wymianę z otoczeniem, a nie bieżący przekład.
  const warsztat: OknoWarsztatTranslate = utworzOknoWarsztatTranslate(stan);

  const wyszukiwarka: WyszukiwarkaFunkcji = utworzWyszukiwarkeFunkcji();
  const izolacja: MacierzIzolacji = utworzMacierzIzolacji(kanal);

  const wyzwalacze: WyzwalaczeOkien = utworzWyzwalaczeOkien([
    {
      kod: KODY_OKIEN.formaty,
      nazwa: 'Format Studio',
      warstwa: 2,
      znacznik: '▼',
      okno: formaty.element,
    },
    {
      kod: KODY_OKIEN.glosariusz,
      nazwa: 'Glossary Manager',
      warstwa: 3,
      znacznik: '☰',
      okno: glosariusz.element,
      skrot: 'Ctrl/Cmd + G',
    },
    {
      kod: KODY_OKIEN.pamiec,
      nazwa: 'Translation Memory Panel',
      warstwa: 3,
      znacznik: '☰',
      okno: pamiec.element,
      skrot: 'Ctrl/Cmd + M',
    },
    {
      kod: KODY_OKIEN.jakosc,
      nazwa: 'QA & Review Center',
      warstwa: 3,
      znacznik: '☰',
      okno: jakosc.element,
      skrot: 'Ctrl/Cmd + Shift + Q',
    },
    {
      kod: KODY_OKIEN.warsztat,
      nazwa: 'Warsztat tłumaczenia',
      warstwa: 3,
      znacznik: '☰',
      okno: warsztat.element,
    },
  ]);

  const pasWiodacych = document.createElement('div');
  pasWiodacych.className = 'mt-modul__pas mt-modul__pas--wiodace';
  pasWiodacych.append(zrodlo.element, panele.element);

  const pasRoboczy = document.createElement('div');
  pasRoboczy.className = 'mt-modul__pas mt-modul__pas--robocze';
  pasRoboczy.append(formaty.element);

  const boczna = document.createElement('div');
  boczna.className = 'mt-modul__boczna';
  boczna.append(glosariusz.element, pamiec.element, jakosc.element, warsztat.element);

  // Pasek uczciwości bierze zdanie z bytu wspólnego (`moduly/katalog-okien.ts`),
  // tego samego dla wszystkich modułów, więc zdanie o rozjeździe nie rozjedzie
  // się między modułami po cichu. Byt wypowiada obie strony: okna katalogu
  // rdzenia, których moduł nie buduje, oraz okna budowane spoza katalogu.
  const katalog = utworzKatalogOkien(kanal, KOD_MODULU, [
    KODY_OKIEN.zrodlo,
    KODY_OKIEN.panele,
    KODY_OKIEN.glosariusz,
    KODY_OKIEN.pamiec,
    KODY_OKIEN.formaty,
    KODY_OKIEN.jakosc,
    KODY_OKIEN.warsztat,
  ]);
  const uczciwosc = katalog.zdanieElement('dn-pole-opis mt-uczciwosc');

  const zaplecze = document.createElement('div');
  zaplecze.className = 'mt-modul__zaplecze';
  zaplecze.append(
    uczciwosc,
    wyszukiwarka.element,
    utworzWykazSkrotow(),
    izolacja.element,
    utworzWystawienieOperacji(kanal).element,
  );

  const element = document.createElement('div');
  element.className = 'mt-modul';
  element.dataset['modul'] = 'translate';
  element.setAttribute('aria-label', 'Moduł Translate — okna operacyjne');
  element.append(kontekst.element, wyzwalacze.element, pasWiodacych, pasRoboczy, boczna, zaplecze);

  const odepnijSkroty = podepnijSkroty(element, {
    dodajJezyk: () => panele.ogniskujDodanieJezyka(),
    poprzedniDoUwagi: () => panele.przejdzDoUwagi(true),
    nastepnyDoUwagi: () => panele.przejdzDoUwagi(false),
    otworzGlosariusz: () => wyzwalacze.otworz(KODY_OKIEN.glosariusz),
    otworzPamiec: () => wyzwalacze.otworz(KODY_OKIEN.pamiec),
    otworzJakosc: () => wyzwalacze.otworz(KODY_OKIEN.jakosc),
    otworzWyszukiwarke: () => wyszukiwarka.otworz(),
  });

  const odsubskrybuj = stan.obserwuj(odswiezWszystko);

  function odswiezWszystko(): void {
    kontekst.odswiez();
    zrodlo.odswiez();
    panele.odswiez();
    formaty.odswiez();
    glosariusz.odswiez();
    pamiec.odswiez();
    jakosc.odswiez();
    warsztat.odswiez();
  }

  odswiezWszystko();

  // Jeden odczyt na cały kanał: jeśli o katalog okien zapytał już inny moduł
  // albo powłoka, to wywołanie nie wyśle drugiego zapytania.
  void katalog.odczytaj();

  return {
    element,

    async wczytaj(idSesji) {
      ostatniaSesja = idSesji;
      // Jedyny odczyt wykonywany bez czynności Operatora w oknach warstwy
      // pierwszej. Macierz izolacji idzie obok, bo dotyczy procesu sesji i jej
      // odmowa nie zatrzymuje modułu — powód zapisuje sama macierz.
      await stan.wczytajKontekst(idSesji);
      void izolacja.odczytaj();
    },

    rozlacz() {
      odsubskrybuj();
      odepnijSkroty();
      // Pasek uczciwości odpina się od wspólnego katalogu: bez tego wpis kanału
      // trzymałby przerysowanie elementu zdjętego już z drzewa.
      katalog.zamknij();
      // Stery kanału zwijają się przed zejściem modułu: rozwinięte trzymają
      // nasłuch na dokumencie, którego zniknięcie elementu nie zdejmuje.
      panele.rozlacz();
      stan.rozlacz();
    },
  };
}
