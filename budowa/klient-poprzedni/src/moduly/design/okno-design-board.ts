import { Command, type DesignAsset } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { OKNO_DESIGN_BOARD } from './etykiety-designu';
import { utworzInspektorWarstwy, type InspektorWarstwy } from './inspektor-warstwy';
import { utworzNarzedziaPlanszy, type NarzedziaPlanszy } from './narzedzia-planszy';
import { utworzPanelWarstw, type PanelWarstw } from './panel-warstw';
import { utworzPasekKontekstu, type PasekKontekstu } from './pasek-kontekstu';
import { utworzTrybyKanwy, type TrybyKanwy } from './tryby-kanwy';
import { utworzWersjeKompozycji, type WersjeKompozycji } from './wersje-kompozycji';
import { utworzAdnotacjeKompozycji, type AdnotacjeKompozycji } from './adnotacje-kompozycji';
import { utworzPlansze, type PlanszaKompozycji } from './plansza-kompozycji';
import { przytnijKod } from './przyciecie-pol';
import { skutekZapisuKompozycji } from './skutek-designu';
import { naglowekOkna, utworzStanOkna, type StanOkna } from './stan-okna';
import { utworzStanKompozycji, type StanKompozycji } from './stan-kompozycji';
import { utworzTabliczkeRozmowy, type TabliczkaRozmowy } from './tabliczka-rozmowy';
import { powodZTorem } from './tor-komendy';
import { zdanieOSkutkuRozmowy, type SkutekRozmowy } from './wejscie-rozmowy';
import type { StanDesignu } from './stan-designu';

/**
 * Design Board — okno wiodące modułu Design: kanwa, panel warstw i przybornik w jednym oknie wraz
 * z zapisem kompozycji do rdzenia jednym polem warstw komendy zbiorczej.
 */
export interface OknoDesignBoard {
  element: HTMLElement;
  odswiez(): void;
  /** Odpina subskrypcje okna — dziś jedną: bez niej kursory współpracy trafiałyby do panelu odłączonego. */
  rozlacz(): void;
  /** Dokłada zasób na kanwę — droga z Assets Panel. */
  przyjmijZasob(zasob: DesignAsset): void;
  /** Przyjmuje skutek rozmowy: zasób nowy ląduje warstwą, każdy inny skutek jest nazwany na tabliczce. */
  przyjmijZRozmowy(skutek: SkutekRozmowy, zasob: DesignAsset): void;
  /** Prowadzi ognisko do adnotacji warstwy zaznaczonej; bez zaznaczenia okno mówi wprost, czego brakuje. */
  opiszWarstwe(): void;
}

export function utworzOknoDesignBoard(stan: StanDesignu): OknoDesignBoard {
  const okno: StanOkna = utworzStanOkna();
  const kompozycja: StanKompozycji = utworzStanKompozycji();
  const plansza: PlanszaKompozycji = utworzPlansze(kompozycja);
  const warstwy: PanelWarstw = utworzPanelWarstw(kompozycja);
  const narzedzia: NarzedziaPlanszy = utworzNarzedziaPlanszy(kompozycja);
  const tabliczka: TabliczkaRozmowy = utworzTabliczkeRozmowy();
  const inspektor: InspektorWarstwy = utworzInspektorWarstwy(kompozycja);
  // Wersjonowanie, wyrys, adnotacje i obecność dotyczą kompozycji jako całości, więc stoją pod kanwą.
  const wersje: WersjeKompozycji = utworzWersjeKompozycji(stan, kompozycja);
  const adnotacje: AdnotacjeKompozycji = utworzAdnotacjeKompozycji(stan, kompozycja);

  // Tryb kanwy znakuje planszę atrybutem danych, widocznym w sprawdzianie i arkuszu, nie klasą składaną.
  const tryby: TrybyKanwy = utworzTrybyKanwy((tryb) => {
    plansza.element.dataset['tryb'] = tryb;
  });

  const kontekst: PasekKontekstu = utworzPasekKontekstu({
    okno: () => stan.idOkna(),
    opisOkna: () => stan.opisOkna(),
    silnikow: () => stan.silniki().length,
    tryb: () => tryby.tryb(),
    powiekszenie: () => kompozycja.widok().powiekszenie,
  });

  const nazwa = poleTekstowe({
    etykieta: 'Nazwa kompozycji',
    podpowiedz: 'np. tablica nastroju — wariant zimowy',
    opis: 'Pole name żądania design.board.update. Puste zostawia kompozycję bez nazwy.',
  });
  nazwa.kontrolka.addEventListener('change', () => kompozycja.ustawNazwe(nazwa.kontrolka.value));

  const zapisz = przycisk('Zapisz kompozycję w rdzeniu', 'dn-btn dn-btn--sm dn-btn--atrament');
  // Droga powrotna kompozycji: bez odczytu okno po odświeżeniu pokazywałoby pustą kanwę.
  const odczytaj = przycisk('Odczytaj kompozycję z rdzenia', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const bok = document.createElement('div');
  bok.className = 'md-plansza__bok';
  bok.append(warstwy.element, inspektor.element);

  const cialo = document.createElement('div');
  cialo.className = 'md-plansza__cialo';
  cialo.append(plansza.element, bok);

  // Pasek kontekstu stoi najwyżej, a tabliczka wejścia przed przybornikiem.
  okno.tresc.append(
    kontekst.element,
    tabliczka.element,
    tryby.element,
    narzedzia.element,
    nazwa.element,
    zapisz,
    odczytaj,
    odpowiedz.element,
    cialo,
    wersje.element,
    adnotacje.element,
  );

  const element = document.createElement('section');
  element.className = 'md-okno md-okno--wiodace';
  element.dataset['okno'] = OKNO_DESIGN_BOARD.kod;
  element.append(naglowekOkna(OKNO_DESIGN_BOARD.nazwa, OKNO_DESIGN_BOARD.rola), okno.element);

  async function zapiszKompozycje(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(
        `Komenda ${Command.DesignBoardUpdate} wymaga okna modułu. ${stan.opisOkna()}`,
        false,
      );
      return;
    }
    okno.ladowanie('Zapis kompozycji w rdzeniu…');
    odpowiedz.pokaz('Zapis kompozycji…', true);
    const nazwaZamowiona = przytnijKod(kompozycja.nazwa());
    const warstwZamowionych = kompozycja.warstwy().length;
    // Zapis pod czuwaniem: po zerwanym gnieździe okno nie ma prawa stać w zapisie bez końca.
    const wynik = await stan.czuwanie.prowadz(
      'zapis kompozycji',
      stan.zrodlo.zapiszKompozycje({
        idOkna: stan.idOkna(),
        idKompozycji: kompozycja.idKompozycji(),
        nazwa: kompozycja.nazwa(),
        warstwy: kompozycja.warstwy(),
      }),
      {
        wToku: (zdanie) => okno.ladowanie(zdanie),
        cisza: (zdanie) => {
          okno.blad(zdanie);
          odpowiedz.pokaz(zdanie, false);
        },
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      const powod = opisOdmowyBledu('Zapis kompozycji', wynik.blad);
      okno.blad(powodZTorem(powod, Command.DesignBoardUpdate));
      odpowiedz.pokaz(powod, false);
      return;
    }
    kompozycja.ustawIdKompozycji(wynik.wynik.board.id);
    const skutek = skutekZapisuKompozycji(nazwaZamowiona, warstwZamowionych, wynik.wynik.board);
    odpowiedz.pokaz(skutek.zdanie, skutek.udany);
    // Faza wraca z ładowania tutaj, bo odświeżenie samo z niej nie wychodzi.
    okno.gotowe();
    odswiez();
  }

  zapisz.addEventListener('click', () => void zapiszKompozycje());
  odczytaj.addEventListener('click', () => void odczytajKompozycje());

  /** Odczyt kompozycji okna z rdzenia; gdy rdzeń odda kilka, okno bierze najświeższą po dacie. */
  async function odczytajKompozycje(): Promise<void> {
    if (stan.idOkna() === '') {
      odpowiedz.pokaz(`Odczyt kompozycji wymaga okna modułu. ${stan.opisOkna()}`, false);
      return;
    }
    okno.ladowanie('Odczyt kompozycji z rdzenia…');
    odpowiedz.pokaz('Odczyt kompozycji…', true);
    const wynik = await stan.czuwanie.prowadz(
      'odczyt kompozycji',
      stan.zrodlo.kompozycje(stan.idOkna()),
      {
        wToku: (zdanie) => okno.ladowanie(zdanie),
        cisza: (zdanie) => {
          okno.blad(zdanie);
          odpowiedz.pokaz(zdanie, false);
        },
        powrot: (zdanie) => odpowiedz.pokaz(zdanie, false),
        spozniona: (zdanie) => odpowiedz.pokaz(zdanie, false),
      },
    );
    if (wynik === null) return;
    if (!wynik.udany || wynik.wynik === undefined) {
      const powod = opisOdmowyBledu('Odczyt kompozycji', wynik.blad);
      okno.blad(powodZTorem(powod, Command.DesignBoardList));
      odpowiedz.pokaz(powod, false);
      return;
    }
    const wykaz = wynik.wynik.boards;
    okno.gotowe();
    if (wykaz.length === 0) {
      odswiez();
      odpowiedz.pokaz(
        'Rdzeń nie zna żadnej kompozycji tego okna — kanwa została bez zmiany. To jest ' +
          'pustka POTWIERDZONA, nie nieudany odczyt.',
        true,
      );
      return;
    }
    const najswiezsza = [...wykaz].sort((pierwsza, druga) => druga.updatedAt - pierwsza.updatedAt)[0];
    if (najswiezsza === undefined) return;
    const warstwy = najswiezsza.layers ?? [];
    kompozycja.wczytaj(najswiezsza.id, najswiezsza.name ?? '', warstwy);
    nazwa.kontrolka.value = najswiezsza.name ?? '';
    odswiez();
    odpowiedz.pokaz(
      `Kanwa pokazuje kompozycję ${najswiezsza.id} z rdzenia — warstw ${warstwy.length}` +
        (wykaz.length === 1
          ? '. Układ na kanwie został ZASTĄPIONY tym z bazy.'
          : `. Rdzeń oddał ${wykaz.length} kompozycji tego okna; wzięto najświeższą ` +
            'po polu updatedAt, bo okno prowadzi jedną kanwę.'),
      true,
    );
  }

  function odswiez(): void {
    plansza.odswiez();
    warstwy.odswiez();
    narzedzia.odswiez();
    inspektor.odswiez();
    kontekst.odswiez();
    if (okno.faza() === 'ladowanie') return;
    if (kompozycja.warstwy().length === 0) {
      // Stan pusty wymienia obie drogi na kanwę: zdanie o samym Assets Panel przemilczałoby rozmowę.
      okno.puste(
        'Plansza jest bez zestawionych jeszcze zasobów — przenieś zasób z Assets Panel. ' +
          'Drugą drogę na kanwę, wejście rozmowy, opisuje tabliczka nad planszą.',
      );
      return;
    }
    okno.gotowe();
  }

  kompozycja.obserwuj(odswiez);
  odswiez();

  /** Czy kanwa niesie już warstwę tego zasobu — jedna praca, jedna warstwa. */
  function czyWarstwaZasobu(idZasobu: string): boolean {
    return kompozycja.warstwy().some((warstwa) => warstwa.assetId === idZasobu);
  }

  return {
    element,
    odswiez,
    rozlacz: () => adnotacje.rozlacz(),

    przyjmijZasob(zasob) {
      kompozycja.dolozZasob(zasob);
    },

    opiszWarstwe() {
      if (inspektor.ogniskujAdnotacje()) return;
      odpowiedz.pokaz(
        'Adnotację nadaje się jednej warstwie — żadna nie jest zaznaczona. Wskaż warstwę ' +
          'na kanwie albo w panelu warstw i naciśnij skrót ponownie.',
        false,
      );
    },

    przyjmijZRozmowy(skutek, zasob) {
      // Warstwa powstaje wyłącznie z zasobu nowego i wyłącznie raz; tabliczka mówi, co okno zrobiło.
      const naKanwie = skutek === 'zasob-nowy' && !czyWarstwaZasobu(zasob.id);
      if (naKanwie) kompozycja.dolozZasob(zasob);
      tabliczka.zanotuj(zdanieOSkutkuRozmowy(skutek, zasob, naKanwie));
    },
  };
}
