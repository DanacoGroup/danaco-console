import type { BrowserSnapshot, WindowListResponse } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { utworzMaterialSesji, type MaterialSesji } from './material-sesji';
import { utworzOdczytMigawki, type StanMigawki } from './odczyt-migawki';
import { utworzWarstwyWidocznosci, type WarstwyWidocznosci } from './warstwy-widocznosci';
import {
  utworzWykazyZebranego,
  type StanZaciagniecia,
  type WykazyZebranego,
} from './wykazy-zebranego';
import { utworzZapisyPrzekazania, type ZapisyPrzekazania } from './zapisy-przekazania';
import { utworzZebraneWSesji, type ZebraneWSesji } from './zebrane-w-sesji';
import { utworzZrodloAutomatyk, type ZrodloAutomatyk } from './zrodlo-automatyk';
import { utworzZrodloBrowser, type ZrodloBrowser } from './zrodlo-browser';
import { utworzZrodloIzolacji, type ZrodloIzolacji } from './zrodlo-izolacji';
import { oknoModulu, opisOkna, utworzZrodloOkien, type ZrodloOkien } from './zrodlo-okien';

/**
 * Jedno źródło prawdy modułu Browser. Pięć okien modułu pracuje na tym samym
 * oknie przeglądarki, tej samej migawce i tym samym zbiorze zebranym; osobne
 * stany rozjechałyby notatkę ze źródłem, a zaznaczenie ze stroną.
 */
export type FazaOdczytu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface StanPrzegladania {
  zrodlo: ZrodloBrowser;
  okna: ZrodloOkien;
  zapisy: ZapisyPrzekazania;
  zebrane: ZebraneWSesji;
  /** Komendy automatyk — zaplecze Automation Studio. */
  automatyki: ZrodloAutomatyk;
  /** Odczyty izolacji — zaplecze macierzy izolacji sesji (warstwa czwarta). */
  izolacja: ZrodloIzolacji;
  /** Materiał przechwycony w sesji — zaplecze Capture & Monitor Panel. */
  material: MaterialSesji;
  /** Pokrycie komend modułu w rdzeniu — jedno na moduł, odczyt idzie raz na połączenie. */
  pokrycie: PokrycieKomend;
  /** Stan stopniowego ujawniania okien i paneli modułu. */
  warstwy: WarstwyWidocznosci;
  /** Okno przeglądarki; pusty napis znaczy „rdzeń go jeszcze nie wskazał". */
  idOkna(): string;
  /** Faza odczytu okna przeglądarki. */
  faza(): FazaOdczytu;
  /** Zdanie o stanie okna przeglądarki — nazwa okna albo powód jego braku. */
  powod(): string;
  /** Migawka strony widoczna jednocześnie Operatorowi i modelowi. */
  migawka(): BrowserSnapshot | null;
  /** Stan ostatniego odczytu migawki — pustka i odmowa to dwie różne rzeczy. */
  fazaMigawki(): StanMigawki;
  /** Zdanie o ostatnim odczycie migawki; puste, gdy przyszła bez uwag. */
  powodMigawki(): string;
  /** Fragment zaznaczony w podglądzie strony. */
  zaznaczenie(): string;
  ustawZaznaczenie(fragment: string): void;
  wchlonMigawke(migawka: BrowserSnapshot): void;
  /** Ustala okno przeglądarki sesji (`window.list`) i jego stan. */
  ustalOkno(idSesji: string): Promise<void>;
  /** Zaciąga z rdzenia treść strony okna (`browser.snapshot.get`). */
  zaciagnijMigawke(zeZrodlem?: boolean): Promise<void>;
  /** Zaciąga z rdzenia wykaz źródeł i notatek okna (dwie komendy `*.list`). */
  zaciagnijZebrane(): Promise<void>;
  /** Zdanie o ostatnim zaciągnięciu wykazów — puste, gdy przyszły bez uwag; osobno od powodu okna. */
  powodZebranego(): string;
  /** Stan ostatniego zaciągnięcia wykazów — rozstrzyga, czy panel pokazuje czekanie, pustkę czy odmowę. */
  stanZebranego(): StanZaciagniecia;
  /** Powiadamia widoki o każdej zmianie stanu. */
  obserwuj(sluchacz: () => void): () => void;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzStanPrzegladania(kanal: Kanal): StanPrzegladania {
  const zrodlo = utworzZrodloBrowser(kanal);
  const okna = utworzZrodloOkien(kanal);
  const zapisy = utworzZapisyPrzekazania(kanal);
  const rozglos = utworzOgloszenia();
  const oglos = rozglos.oglos;
  // Zbiór zebranego ogłasza zmianę tym samym kanałem co odczyt okna, żeby okna zobaczyły ją naraz.
  const zebrane = utworzZebraneWSesji(oglos);
  const odczyt = utworzOdczytOkna(okna, utworzWykazyZebranego(zrodlo, zebrane), oglos);
  const migawki = utworzOdczytMigawki(zrodlo, oglos);
  const material = utworzMaterialSesji(oglos);
  const warstwy = utworzWarstwyWidocznosci(oglos);
  const pokrycie = utworzPokrycieKomend(kanal);

  let fragment = '';

  // Zdarzenie cudzego okna jest odrzucane, a przed ustaleniem okna odrzucane jest każde zdarzenie.
  const odsubskrybuj = zrodlo.naZmianeStrony((tresc) => {
    if (odczyt.idOkna() === '' || tresc.windowId !== odczyt.idOkna()) return;
    migawki.wchlon(tresc.snapshot);
  });

  /** Odczyt migawki adresowany do okna ustalonego przez `odczyt`. */
  const zaciagnijMigawke = (zeZrodlem = false): Promise<void> =>
    migawki.zaciagnij(odczyt.idOkna(), zeZrodlem);

  return {
    zrodlo,
    okna,
    zapisy,
    zebrane,
    automatyki: utworzZrodloAutomatyk(kanal),
    izolacja: utworzZrodloIzolacji(kanal),
    material,
    warstwy,
    pokrycie,
    idOkna: odczyt.idOkna,
    faza: odczyt.faza,
    powod: odczyt.powod,
    migawka: migawki.wartosc,
    fazaMigawki: migawki.stan,
    powodMigawki: migawki.powod,
    zaznaczenie: () => fragment,

    ustawZaznaczenie(nowy) {
      if (nowy === fragment) return;
      fragment = nowy;
      oglos();
    },

    wchlonMigawke: migawki.wchlon,

    // Migawka idzie zaraz po ustaleniu okna, tak samo jak wykazy — rdzeń zna treść strony zawsze.
    async ustalOkno(idSesji) {
      await odczyt.ustalOkno(idSesji);
      await zaciagnijMigawke();
    },

    zaciagnijMigawke,
    zaciagnijZebrane: odczyt.zaciagnijZebrane,
    powodZebranego: odczyt.powodZebranego,
    stanZebranego: odczyt.stanZebranego,
    obserwuj: rozglos.obserwuj,

    rozlacz() {
      odsubskrybuj();
      // Pozycje pokrycia odpinają się od wspólnego wykazu z tego samego powodu co pasek uczciwości.
      pokrycie.zamknij();
      rozglos.wyczysc();
    },
  };
}

/**
 * Rejestr widoków modułu: jedno ogłoszenie dociera do wszystkich okien naraz.
 *
 * Osobno od stanu, bo powiadamianie nie zależy od tego, co się zmieniło —
 * i dzięki temu każdy kawałek stanu ogłasza się dokładnie tak samo.
 */
interface Ogloszenia {
  oglos(): void;
  obserwuj(sluchacz: () => void): () => void;
  /** Zdejmuje wszystkich słuchaczy — moduł kończy pracę. */
  wyczysc(): void;
}

function utworzOgloszenia(): Ogloszenia {
  const sluchacze = new Set<() => void>();
  return {
    oglos() {
      for (const sluchacz of [...sluchacze]) sluchacz();
    },
    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },
    wyczysc: () => sluchacze.clear(),
  };
}

/**
 * Okno przeglądarki modułu: jego ustalenie, faza odczytu i wykazy okna.
 *
 * Osobno od migawki i zaznaczenia: tutaj mieszka rozmowa z rdzeniem o tym, do
 * którego okna moduł adresuje komendy, tam — to, co Operator i model widzą
 * na stronie.
 */
interface OdczytOkna {
  idOkna(): string;
  faza(): FazaOdczytu;
  powod(): string;
  powodZebranego(): string;
  stanZebranego(): StanZaciagniecia;
  ustalOkno(idSesji: string): Promise<void>;
  zaciagnijZebrane(): Promise<void>;
}

function utworzOdczytOkna(
  okna: ZrodloOkien,
  wykazy: WykazyZebranego,
  oglos: () => void,
): OdczytOkna {
  let identyfikator = '';
  let faza: FazaOdczytu = 'spoczynek';
  let powod = 'Okno przeglądarki nie zostało jeszcze odczytane z rdzenia.';
  let powodZebranego = '';
  let stanZebranego: StanZaciagniecia = 'nietkniete';

  /** Wpisuje rozstrzygnięcie do stanu i ogłasza je oknom modułu. */
  function przyjmij(nowe: RozstrzygniecieOkna): void {
    faza = nowe.faza;
    powod = nowe.powod;
    if (nowe.identyfikator !== null) identyfikator = nowe.identyfikator;
    oglos();
  }

  // Odmowa wykazu nie wywraca okna w stan błędu — wynik nieudanego odczytu idzie osobnym stanem.

  // Zapowiedź odczytu idzie przed wywołaniem komend, żeby panele nie stały cały czas w stanie pustym.
  async function zaciagnijZebrane(): Promise<void> {
    stanZebranego = 'odczyt';
    powodZebranego = 'Odczyt wykazów okna z rdzenia w toku…';
    oglos();
    const wynik = await wykazy.zaciagnij(identyfikator);
    stanZebranego = wynik.stan;
    powodZebranego = wynik.powod;
    oglos();
  }

  return {
    idOkna: () => identyfikator,
    faza: () => faza,
    powod: () => powod,
    powodZebranego: () => powodZebranego,
    stanZebranego: () => stanZebranego,
    zaciagnijZebrane,

    async ustalOkno(idSesji) {
      if (idSesji === '') {
        przyjmij(brakSesji());
        return;
      }
      przyjmij(odczytWToku());
      const rozstrzygniecie = rozstrzygnijOkno(await okna.okna(idSesji));
      przyjmij(rozstrzygniecie);
      // Wykaz zaciąga się dopiero z identyfikatorem okna, tak samo jak odczyt stanu okna.
      if (rozstrzygniecie.zaciagacWykazy) await zaciagnijZebrane();
    },
  };
}

/**
 * Co odpowiedź rdzenia znaczy dla stanu okna — wyliczone osobno od zapisu.
 *
 * `identyfikator: null` znaczy „zostaw dotychczasowy": odmowa jednego odczytu
 * nie unieważnia okna raz ustalonego, więc nie wolno jej wyczyścić numeru.
 */
interface RozstrzygniecieOkna {
  faza: FazaOdczytu;
  powod: string;
  identyfikator: string | null;
  /** Czy iść po wykazy okna — idzie się wyłącznie z jego identyfikatorem. */
  zaciagacWykazy: boolean;
}

/** Stan sprzed otwarcia sesji: modułowi nie ma kto wskazać okna przeglądarki, więc odczyt jeszcze nie ruszył. */
function brakSesji(): RozstrzygniecieOkna {
  return {
    faza: 'spoczynek',
    powod: 'Sesja nie jest jeszcze otwarta — moduł nie ma okna, do którego adresuje komendy.',
    identyfikator: null,
    zaciagacWykazy: false,
  };
}

/** Zapowiedź odczytu okna sesji: okna modułu mają pokazać czekanie, a nie pustkę, dopóki rdzeń nie odpowie. */
function odczytWToku(): RozstrzygniecieOkna {
  return {
    faza: 'odczyt',
    powod: 'Odczyt okien sesji w toku…',
    identyfikator: null,
    zaciagacWykazy: false,
  };
}

/** Odpowiedź na odczyt okien sesji przełożona na stan okna modułu, wraz z pełnym rozstrzygnięciem odmowy. */
function rozstrzygnijOkno(wynik: Wynik<WindowListResponse>): RozstrzygniecieOkna {
  if (!wynik.udany || wynik.wynik === undefined) {
    return {
      faza: 'blad',
      powod: opisOdmowy('Odczyt okien sesji', wynik.blad?.code, wynik.blad?.message),
      identyfikator: null,
      zaciagacWykazy: false,
    };
  }
  const okno = oknoModulu(wynik.wynik.windows);
  if (okno === null) {
    return {
      faza: 'gotowe',
      powod:
        'Rdzeń nie przypisał tej sesji okna modułu Browser — komendy obszaru nie mają adresata.',
      identyfikator: '',
      zaciagacWykazy: false,
    };
  }
  return {
    faza: 'gotowe',
    powod: `Okno przeglądarki: ${opisOkna(okno)}.`,
    identyfikator: okno.id,
    zaciagacWykazy: true,
  };
}
