import { opisOdmowy } from '../../komponenty/odmowa';
import { ADNOTACJA, BARWY_ADNOTACJI, POZYCJE_BEZ_OBSLUGI } from './etykiety-browser';
import { utworzPasekAdnotacji } from './pasek-adnotacji';
import { utworzPlotnoAdnotacji } from './plotno-adnotacji';
import { skutekAdnotacji } from './skutek-zapisu';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Tryb adnotacji — płótno nad sceną podglądu wraz z pływającym paskiem
 * i drogą „Dodaj do rozmowy". Plik składa warstwę i prowadzi jej rozmowę
 * z rdzeniem; warstwa leży nad sceną i nie dotyka podglądu.
 */
export interface WarstwaAdnotacji {
  /** Warstwa osadzana nad sceną podglądu strony. */
  element: HTMLElement;
  wlaczona(): boolean;
  /** Włącza albo wyłącza tryb; zwraca stan po zmianie. */
  ustaw(wlaczona: boolean): boolean;
  /** Spłaszcza płótno do PNG i wysyła je do okna rozmowy. */
  dodajDoRozmowy(): Promise<void>;
  /** Odpina obserwatora rozmiaru płótna — moduł kończy pracę. */
  rozlacz(): void;
}

/** Czego warstwa adnotacji potrzebuje od ramy okna: ujścia zdarzeń trybu, wysyłki oraz zamknięcia pracy. */
export interface UjsciaWarstwy {
  /** Odpowiedź pokazywana Operatorowi po każdym naciśnięciu. */
  powiedz(tresc: string, powodzenie: boolean): void;
  /** Zgłasza zmianę trybu ramie okna — przełącznik paska dolnego jest lustrem. */
  naZmianeTrybu(wlaczona: boolean): void;
}

export function utworzWarstweAdnotacji(
  stan: StanPrzegladania,
  ujscia: UjsciaWarstwy,
): WarstwaAdnotacji {
  // Barwa początkowa bierze się z pierwszej pozycji wykazu, a nie z osobnego literału.
  const plotno = utworzPlotnoAdnotacji(BARWY_ADNOTACJI[0]?.zeton ?? '');

  const bezTla = document.createElement('p');
  bezTla.className = 'mb-adnotacja__bez-tla';
  // Najpierw stoi to, czym model dysponuje, bo od tego zależy dopisane zdanie; potem brak tła i wpisu.

  // Zdanie o wpisie w wytworach przerysowuje się po każdej zmianie wykazu komend rdzenia.
  stan.pokrycie.naOdczyt(() => {
    const oWytworze = stan.pokrycie.zdanie(
      POZYCJE_BEZ_OBSLUGI.wytworAdnotacji.komenda,
      POZYCJE_BEZ_OBSLUGI.wytworAdnotacji.czynnosc,
    );
    bezTla.textContent = `${ADNOTACJA.modelDostajeSciezke} ${ADNOTACJA.bezTla} ${oWytworze}`;
  });

  const pasek = utworzPasekAdnotacji({
    narzedzie: plotno.narzedzie,
    ustawNarzedzie: plotno.ustawNarzedzie,
    zeton: plotno.zeton,
    ustawZeton: plotno.ustawZeton,
    ustawNapis: plotno.ustawNapis,
    wyczysc() {
      const bylo = plotno.liczbaSladow();
      plotno.wyczysc();
      // Zdanie podaje liczbę zdjętych śladów, żeby nie potwierdzać czynności, która się nie odbyła.
      ujscia.powiedz(
        bylo === 0
          ? 'Płótno adnotacji było już puste — nie było czego czyścić.'
          : `Zdjęto ${bylo} śladów adnotacji; płótno jest puste.`,
        bylo > 0,
      );
    },
    dodajDoRozmowy: () => void dodajDoRozmowy(),
    zamknij: () => ustaw(false),
  });

  const element = document.createElement('div');
  element.className = 'mb-adnotacja';
  element.hidden = true;
  element.setAttribute('aria-label', 'Warstwa adnotacji nad podglądem strony');
  element.append(plotno.element, bezTla, pasek.element);

  function ustaw(wlaczona: boolean): boolean {
    element.hidden = !wlaczona;
    // Bufor płótna ma rozmiar dopiero, gdy warstwa jest widoczna — schowany element mierzy zero.
    if (wlaczona) plotno.dopasuj();
    ujscia.naZmianeTrybu(wlaczona);
    return wlaczona;
  }

  // Treść opisuje rysunek słowami, bo model dostaje ścieżkę do pliku, nie sam obraz.
  function trescWiadomosci(): string {
    const migawka = stan.migawka();
    const skad = migawka === null ? 'strony bez odczytanej migawki' : migawka.url;
    return (
      `Adnotacja Operatora na podglądzie ${skad} — ${plotno.liczbaSladow()} śladów, ` +
      `spłaszczona do PNG i dołączona jako obraz. ${ADNOTACJA.modelDostajeSciezke} ` +
      `${ADNOTACJA.bezTla}`
    );
  }

  async function dodajDoRozmowy(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      ujscia.powiedz(stan.powod(), false);
      return;
    }
    if (plotno.liczbaSladow() === 0) {
      ujscia.powiedz(ADNOTACJA.pusteBezWysylki, false);
      return;
    }
    // Spłaszczenie sprawdzane przed wysyłką — pusty załącznik byłby obietnicą obrazu, którego nie ma.
    const obraz = plotno.doPng();
    if (!obraz.startsWith('data:image/png')) {
      ujscia.powiedz(ADNOTACJA.brakSplaszczenia, false);
      return;
    }
    ujscia.powiedz(ADNOTACJA.wysylkaWToku, true);
    const wynik = await stan.zapisy.zapytaj(idOkna, trescWiadomosci(), [obraz]);
    if (!wynik.udany || wynik.wynik === undefined) {
      ujscia.powiedz(
        opisOdmowy('Dodanie adnotacji do rozmowy', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    const skutek = skutekAdnotacji(wynik.wynik.message, idOkna, obraz);
    ujscia.powiedz(skutek.zdanie, skutek.udany);
  }

  return {
    element,
    wlaczona: () => !element.hidden,
    ustaw,
    dodajDoRozmowy,
    rozlacz: plotno.rozlacz,
  };
}
