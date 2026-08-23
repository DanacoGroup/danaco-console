import './aod.css';

import type { Kanal } from '../protokol/kanal';
import { StanAwatara, utworzAwatarAod } from './awatar-aod';
import { utworzDymekSugestii } from './dymek-sugestii';
import type { WpisKolejkiDecyzji } from './kolejka-decyzji';
import { utworzKolumneAod, type KolumnaAod, type OpisKolumnyAod } from './kolumna-aod';
import { WagaUjawnienia } from './rodzaje-sugestii';
import {
  TrybObecnosci,
  utworzStanObecnosci,
  type OpisUjawnienia,
  type StanObecnosci,
} from './tryb-obecnosci';
import { KLASA_POWODU, type MagazynWyciszen } from './wyciszenie-aod';
import { utworzKontekstWyciszenia, type KontekstWyciszenia } from './wyciszenie-kontekst';
import { utworzMenuWyciszenia } from './wyciszenie-menu';

/**
 * Warstwa Always On Display — trzy elementy funkcji globalnej złożone w jedną
 * powierzchnię: pływający awatar (warstwa 1), dymek kontekstowy sugestii
 * i kolumna boczna powierzchni interakcji (warstwa 2).
 *
 * Warstwa leży na `--dn-z-aod` (1200 wg `design/KANON.md`) i osadza się poza
 * obszarem podmienianym przez moduł, więc przełączenie środowiska ani modułu
 * nie zabiera awatara i nie przeładowuje kolumny (rozdz. 2.4 opracowania).
 *
 * To jedyne miejsce, które zna naraz cztery rzeczy i wiąże je regułą
 * opracowania:
 *   • kolejkę decyzji — ile sugestii czeka i o jakiej wadze;
 *   • stan obecności — tryb pełny/cichy/ukryty i trwające wyciszenie;
 *   • progi ujawniania — waga, limit godzinowy, odstęp między dymkami;
 *   • układ — szerokość otwartej kolumny, o którą odsuwa się awatar i dymek.
 *
 * Skróty klawiszowe załącznika A.1 są zaczepione tutaj, bo tylko tutaj widać
 * wszystkie trzy elementy naraz.
 */

export interface WarstwaAod {
  /** Element do osadzenia w powłoce, poza obszarem roboczym modułu. */
  element: HTMLElement;
  /** Otwiera powierzchnię interakcji — listwa ustawień i skrót klawiszowy. */
  otworzPowierzchnie(): void;
  /** Stan obecności — do sterowania z okna Ustawień, gdy takie wejście powstanie. */
  obecnosc: StanObecnosci;
  /** Kolumna powierzchni interakcji. */
  kolumna: KolumnaAod;
  rozlacz(): void;
}

/**
 * Szerokość obszaru roboczego, poniżej której kolumna nie mieści się obok pracy.
 *
 * Rozdz. 2.5 opracowania: przy kolumnach zwężonych poniżej progu czytelności
 * powierzchnia interakcji zwija się do dymka kontekstowego, a dymek skraca
 * treść do jednego zdania z działaniem „Rozwiń".
 */
const PROG_CZYTELNOSCI = 640;

/**
 * @param magazyn magazyn stanu wyciszeń i trybu obecności — PODAWANY, nie brany
 *   na sztywno. Pominięty znaczy zapis miejscowy przeglądarki, `null` znaczy
 *   „bez zapisu". Gdy kontrakt poniesie wyciszenie nakładki jako byt rdzenia,
 *   podmiana magazynu w tym jednym wywołaniu przełoży zapis na rdzeń bez zmiany
 *   ani jednego wołacza.
 */
export function utworzWarstweAod(
  kanal: Kanal,
  opis: OpisKolumnyAod = {},
  magazyn?: MagazynWyciszen | null,
): WarstwaAod {
  const obecnosc = utworzStanObecnosci(magazyn === undefined ? {} : { magazyn });
  const teraz = (): number => Date.now();

  const element = document.createElement('div');
  element.className = 'ao-warstwa';

  /**
   * Kontekst wyciszenia kontekstowego — jedna kopia na warstwę.
   *
   * Czyta bieżący moduł i bieżącą kartę sesji z rdzenia. Menu przy awatarze
   * i menu w nagłówku powierzchni sięgają po ten sam odczyt, więc obie kopie
   * menu nazywają ten sam byt tą samą nazwą.
   */
  const kontekstWyciszenia: KontekstWyciszenia = utworzKontekstWyciszenia(kanal);

  const kolumna = utworzKolumneAod(kanal, obecnosc, opis, kontekstWyciszenia);

  const dymek = utworzDymekSugestii({
    stery: kolumna.stery,
    naOdloz: () => {
      // „Odłóż": pozycja wraca do listy oczekujących, status pozostaje `nowa`.
      dymek.schowaj();
      odswiezWyglad();
    },
    naOdrzuc: (wpis) => {
      kolumna.kolejka.odrzuc(wpis.decyzja.klucz);
      kolumna.przerysujDecyzje();
      dymek.schowaj();
      odswiezWyglad();
    },
    naRozwin: () => {
      dymek.schowaj();
      kolumna.ogniskujSugestie();
      odswiezWyglad();
    },
    naZamkniecie: () => {
      dymek.schowaj();
      odswiezWyglad();
    },
  });

  const awatar = utworzAwatarAod({
    naKlikniecie: () => przelaczDymek(),
    naKlikniecePodwojne: () => {
      dymek.schowaj();
      kolumna.przelacz();
      odswiezWyglad();
    },
  });

  /**
   * Menu kebab NA POWIERZCHNI INTERAKCJI AWATARA — wyciszenie od ręki.
   *
   * Rozstrzygnięcie Właściciela z 17.08.2026: „dodatkowa akcja wyciszenia
   * bezpośrednio z pozycji awatara". Opracowanie wskazuje tę samą drogę —
   * rozdz. 2.6 i 8.3 nazywają menu kebab (⋮) miejscem wyciszania — więc menu
   * jest tym samym komponentem, który stoi w nagłówku kolumny, a nie wzorem
   * nowym.
   *
   * Kebab jest OSOBNYM wyzwalaczem: kliknięcie pojedyncze awatara nadal otwiera
   * dymek, podwójne nadal otwiera powierzchnię interakcji. Wyciszenie nie
   * wymaga przejścia do kolumny.
   */
  const menuAwatara = utworzMenuWyciszenia({
    stan: obecnosc,
    teraz,
    kontekst: kontekstWyciszenia,
    naOtwarcie: () => {
      void kontekstWyciszenia.odswiez().then(() => menuAwatara.odswiez(teraz()));
    },
    etykietaZnaku: 'Menu wyciszania Always On Display — wyciszenie od ręki z pozycji awatara',
  });
  menuAwatara.element.classList.add('ao-wyciszenie--awatar');

  element.append(kolumna.element, dymek.element, awatar.element, menuAwatara.element);

  /**
   * Klucz decyzji, którą ostatnio ujawniono samoczynnie.
   *
   * Bez tego ta sama sugestia otwierałaby dymek przy każdym zdarzeniu telemetrii
   * dotyczącym tego samego procesu — a rozdz. 3.1 mówi, że funkcja nie komentuje
   * pracy w sposób ciągły.
   */
  let ostatnioUjawniona: string | null = null;

  /** Najpilniejszy wpis kolejki albo `null`, gdy nic nie czeka. */
  function najpilniejszy(): WpisKolejkiDecyzji | null {
    const wpisy = kolumna.kolejka.wykaz(teraz());
    if (wpisy.length === 0) return null;
    // Wykaz idzie od najdłużej czekającego; wagę wysoką przepuszczamy przed nią.
    const wysoka = wpisy.find((wpis) => wpis.decyzja.wagaUjawnienia === WagaUjawnienia.Wysoka);
    return wysoka ?? wpisy[0] ?? null;
  }

  function pokazDymek(wpis: WpisKolejkiDecyzji): void {
    dymek.ustawWaski(globalThis.innerWidth - kolumna.szerokosc() < PROG_CZYTELNOSCI);
    dymek.pokaz(wpis, kolumna.dopasujKolejke(wpis.decyzja.idOkna, wpis.decyzja.idSesji));
    odswiezWyglad();
  }

  function przelaczDymek(): void {
    if (dymek.czyWidoczny()) {
      dymek.schowaj();
      odswiezWyglad();
      return;
    }

    const wpis = najpilniejszy();
    if (wpis === null) {
      // Rozdz. 3.6: przy braku istotnych zdarzeń funkcja nie wypełnia ciszy
      // treścią zastępczą. Kliknięcie awatara przy pustej kolejce prowadzi więc
      // wprost do powierzchni interakcji, a nie do pustego dymka.
      kolumna.otworz();
      odswiezWyglad();
      return;
    }

    pokazDymek(wpis);
  }

  /**
   * Reguła samoczynnego ujawnienia — rozdz. 3.4 i 3.5 razem.
   *
   * Dymek otwiera się sam wyłącznie przy wadze wysokiej, przy nienaruszonym
   * limicie godzinowym i odstępie, poza trybem cichym i poza wyciszeniem.
   * Sugestia krytyczna wyciszona ujawnia się plakietką, bez dymka — i to robi
   * `odswiezWyglad`, nie ta funkcja.
   */
  function rozwazUjawnienie(): void {
    const wpis = najpilniejszy();
    if (wpis === null) return;
    if (wpis.decyzja.klucz === ostatnioUjawniona) return;
    if (dymek.czyWidoczny()) return;

    const chwila = teraz();
    if (!obecnosc.czyOtworzycDymek(opisUjawnienia(wpis), chwila)) {
      return;
    }

    ostatnioUjawniona = wpis.decyzja.klucz;
    obecnosc.odnotujDymek(chwila);
    pokazDymek(wpis);
  }

  /**
   * Decyzja opisana tym, co reguła wyciszenia musi o niej wiedzieć.
   *
   * Klasa zdarzenia wychodzi z powodu rozpoznania (rozdz. 3.2), karta sesji —
   * z telemetrii, moduł — z okrężnego odczytu `window.list` w kontekście
   * wyciszenia. Modułu, którego nakładka nie rozpoznaje, nie podstawiamy:
   * sugestia bez modułu nie wpada w wyciszenie modułu.
   */
  function opisUjawnienia(wpis: WpisKolejkiDecyzji): OpisUjawnienia {
    const idModulu = kontekstWyciszenia.modulOkna(wpis.decyzja.idOkna);
    return {
      waga: wpis.decyzja.wagaUjawnienia,
      krytyczna: wpis.decyzja.krytyczna,
      klasa: KLASA_POWODU[wpis.decyzja.powod],
      ...(idModulu === undefined ? {} : { idModulu }),
      ...(wpis.decyzja.idSesji === undefined ? {} : { idSesji: wpis.decyzja.idSesji }),
    };
  }

  /** Stan awatara wyliczony ze stanu obecności i z kolejki — rozdz. 2.4 i 9.1. */
  function stanAwatara(liczba: number, wysoka: boolean, chwila: number): StanAwatara {
    if (obecnosc.tryb() === TrybObecnosci.Ukryty) return StanAwatara.Ukryty;
    // Rozdz. 9.1: stan „Wyciszony" zachodzi przy wyciszeniu czasowym,
    // kontekstowym ALBO klasy zdarzeń — nie tylko przy czasowym.
    if (obecnosc.czyJakiekolwiekWyciszenie(chwila)) return StanAwatara.Wyciszony;
    if (liczba === 0) return StanAwatara.Spoczynek;
    return wysoka ? StanAwatara.WagaWysoka : StanAwatara.SugestiaOczekujaca;
  }

  function odswiezWyglad(): void {
    const chwila = teraz();
    const wpisy = kolumna.kolejka.wykaz(chwila);

    /**
     * Sugestie, które wyciszenie wybiórcze wstrzymuje — kontekstowe i klasy
     * zdarzeń (rozdz. 3.5: „pozostałe zachowują pełne działanie").
     *
     * Plakietka liczy więc sugestie ujawniane, a nie wszystkie: liczba
     * obejmująca wstrzymane obiecywałaby Operatorowi coś, czego wyciszenie mu
     * nie pokaże. Że coś milczy, mówi stan awatara i menu wyciszania.
     */
    const ujawniane = wpisy.filter((wpis) => !obecnosc.czySugestiaWstrzymana(opisUjawnienia(wpis), chwila));
    const wysoka = ujawniane.some((wpis) => wpis.decyzja.wagaUjawnienia === WagaUjawnienia.Wysoka);

    awatar.ustawStan(stanAwatara(ujawniane.length, wysoka, chwila));

    // Wyciszenie czasowe chowa plakietkę (rozdz. 3.5) — z wyjątkiem sugestii
    // krytycznej, która ujawnia się plakietką mimo KAŻDEGO wyciszenia. Wtedy
    // plakietka liczy same sugestie krytyczne: liczba obejmująca wyciszone
    // obiecywałaby Operatorowi coś, czego wyciszenie mu nie pokaże.
    const krytyczne = wpisy.filter((wpis) => wpis.decyzja.krytyczna);
    const plakietkaOgolna = obecnosc.czyPlakietkaWidoczna(chwila);
    awatar.ustawLiczbe(plakietkaOgolna ? ujawniane.length : krytyczne.length);
    awatar.ustawPlakietkeWidoczna(plakietkaOgolna || krytyczne.length > 0);
    menuAwatara.odswiez(chwila);

    const odsuniecie = kolumna.szerokosc();
    awatar.ustawOdsuniecie(odsuniecie);
    // Kebab awatara jedzie z awatarem: stoi pod nim, przy tej samej krawędzi.
    menuAwatara.element.style.setProperty('--ao-odsuniecie', `${Math.max(0, odsuniecie)}px`);
    dymek.ustawOdsuniecie(odsuniecie);
    dymek.ustawWaski(globalThis.innerWidth - odsuniecie < PROG_CZYTELNOSCI);

    element.dataset['kolumna'] = kolumna.czyOtwarta() ? 'otwarta' : 'zamknieta';
    zwezObszarRoboczy(odsuniecie);
  }

  /**
   * Zwężenie obszaru roboczego o szerokość otwartej kolumny.
   *
   * Rozdz. 2.1 opracowania: „Otwarcie powierzchni interakcji ZWĘŻA kolumny
   * obszaru roboczego, NIE PRZESŁANIA ich". Sama kolumna leży w warstwie
   * pozycjonowanej stale, więc bez tego kroku kładłaby się na pracy zamiast
   * ustąpić jej miejsca.
   *
   * Zwężenie idzie wyściółką rodzica warstwy — powłoki, w której warstwa
   * siedzi — a nie zmianą w katalogu powłoki: funkcja globalna dokłada się do
   * układu, a nie przepisuje go. Klasa `ao-zwezenie` i szerokość jadą razem,
   * więc powłoka bez otwartej kolumny nie nosi po niej śladu.
   */
  function zwezObszarRoboczy(szerokosc: number): void {
    const gospodarz = element.parentElement;
    if (gospodarz === null) return;

    if (szerokosc <= 0) {
      gospodarz.classList.remove('ao-zwezenie');
      gospodarz.style.removeProperty('--ao-zwezenie');
      return;
    }

    gospodarz.classList.add('ao-zwezenie');
    gospodarz.style.setProperty('--ao-zwezenie', `${szerokosc}px`);
  }

  // --- skróty klawiszowe załącznika A.1 -----------------------------------

  function naKlawisz(zdarzenie: KeyboardEvent): void {
    const modyfikator = (zdarzenie.ctrlKey || zdarzenie.metaKey) && zdarzenie.shiftKey;
    if (!modyfikator) return;

    switch (zdarzenie.key.toLowerCase()) {
      case 'a':
        kolumna.przelacz();
        break;
      case 'v':
        kolumna.ogniskujGlos();
        break;
      case 'q':
        kolumna.ogniskujSugestie();
        break;
      case 'm':
        obecnosc.przelaczWyciszenieKwadransem(teraz());
        break;
      case 'h':
        obecnosc.przelaczUkrycie();
        break;
      default:
        return;
    }

    zdarzenie.preventDefault();
    odswiezWyglad();
  }

  document.addEventListener('keydown', naKlawisz);

  const przestanObserwowacObecnosc = obecnosc.obserwuj(() => {
    if (obecnosc.tryb() === TrybObecnosci.Ukryty) dymek.schowaj();
    odswiezWyglad();
  });

  const przestanObserwowacUklad = kolumna.obserwujUklad(() => {
    odswiezWyglad();
    rozwazUjawnienie();
  });

  odswiezWyglad();

  // Odczyt nadrabiający idzie od razu, bez otwierania kolumny: plakietka ma
  // mówić prawdę, zanim Operator cokolwiek kliknie.
  void kolumna.odswiezDecyzje().then(() => {
    odswiezWyglad();
    rozwazUjawnienie();
  });

  // Kontekst wyciszenia idzie tym samym torem: menu przy awatarze ma nazywać
  // bieżący moduł i bieżącą kartę sesji pełną nazwą, zanim Operator je kliknie.
  // Odmowa odczytu nie wywraca warstwy — pozycje kontekstowe mówią wtedy, czego
  // brakuje i po czyjej stronie.
  void kontekstWyciszenia.odswiez().then(() => odswiezWyglad());

  return {
    element,

    otworzPowierzchnie() {
      kolumna.otworz();
      odswiezWyglad();
    },

    obecnosc,
    kolumna,

    rozlacz() {
      document.removeEventListener('keydown', naKlawisz);
      przestanObserwowacObecnosc();
      przestanObserwowacUklad();
      zwezObszarRoboczy(0);
      menuAwatara.rozlacz();
      awatar.rozlacz();
      dymek.rozlacz();
      kolumna.rozlacz();
      element.remove();
    },
  };
}
