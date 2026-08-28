import type { Window } from '../../../shared/contract';
import type { GniazdoOkna } from '../okna-rownolegle/indeks';
import { utworzPasekZlecenia } from '../okno-komunikacji/pasek-zlecenia';
import { utworzPrzelacznikSrodowiska } from '../okno-komunikacji/przelacznik-srodowiska';
import { utworzZrodloSrodowiska } from '../okno-komunikacji/zrodlo-srodowiska';
import { utworzZrodloZlecenia } from '../okno-komunikacji/zrodlo-zlecenia';
import type { Kanal } from '../protokol/kanal';
import {
  rozpoznajWidokZapisu,
  WIDOKI_ZAPISU,
  zamontujRozmowe,
  type ZamontowanaRozmowa,
} from '../rozmowa/indeks';
import type { RejestrKanalow } from '../sterowanie/indeks';
import { zamontujWidokSterowania, type WidokSterowania } from '../widok-sterowania/indeks';
import { zwiazStanPary } from './stan-pary-gniazda';
import { politykaOknaModulu } from './ulotnosc-okna';

/** Zależności powiązania gniazda sceny z oknem otwartym w rdzeniu: kanał, gniazdo, okno i kolumna sterowania. */
export interface ZaleznosciWiazania {
  kanal: Kanal;
  /** Wspólny katalog kanałów modelu — jeden na klienta. */
  rejestrKanalow: RejestrKanalow;
  /** Gniazdo sceny, w którym okno ma pracować. */
  gniazdo: GniazdoOkna;
  /** Okno potwierdzone przez rdzeń. */
  okno: Window;
  /** Tożsamość mówiącego po stronie modelu — nazwa kanału. */
  persona: string;
  /** Kolumna sterowania obok sceny. */
  panel: HTMLElement;
}

/** Gniazdo związane z oknem rdzenia: rozmowa, komplet sterowania oraz identyfikatory okna i sesji rdzenia. */
export interface WiazanieGniazda {
  /** Identyfikator okna nadany przez rdzeń. */
  idOkna: string;
  // Sesja, do której należy okno; bez niej scena nie umiałaby powiedzieć, czyje okna pokazuje.
  idSesji: string;
  /** Rozmowa osadzona w gnieździe. */
  rozmowa: ZamontowanaRozmowa;
  /** Kolumna sterowania tego okna. */
  sterowanie: WidokSterowania;
  /** Gniazdo poza sceną chowa swoją kolumnę sterowania, zamiast ją porzucać. */
  ustawWidocznosc(widoczne: boolean): void;
  /** Odłącza subskrypcje i zdejmuje obie warstwy. */
  rozlacz(): void;
}

/** Powiązanie jednego gniazda sceny z jednym oknem rdzenia: rozmowa i komplet sterowania osadzone razem. */
export function zwiazGniazdoZOknem(zaleznosci: ZaleznosciWiazania): WiazanieGniazda {
  const { kanal, rejestrKanalow, gniazdo, okno, persona, panel } = zaleznosci;

  // Pamięć rozmowy idzie za modułem okna; politykę składa `ulotnosc-okna`.

  // Przełącznik środowiska wchodzi tu, bo to jedyne miejsce mające oba końce sterowania i rozmowy.
  const zrodloSrodowiska = utworzZrodloSrodowiska(kanal, okno);
  const przelacznikSrodowiska = utworzPrzelacznikSrodowiska(zrodloSrodowiska);

  // Pasek zlecenia wchodzi tu z tego samego powodu, co przełącznik — jedyne miejsce z oboma końcami.
  const zrodloZlecenia = utworzZrodloZlecenia(kanal, okno);
  const pasekZlecenia = utworzPasekZlecenia({
    kanal,
    zrodlo: zrodloZlecenia,
    rejestrKanalow,
    sterowanieSrodowiska: przelacznikSrodowiska.element,
    // Stopka steru katalogów prowadzi do kolumny sterowania — jedynego miejsca z polem ścieżki.
    otworzSterowanie: () => sterowanie.otworz(),
  });

  const rozmowa = zamontujRozmowe(gniazdo.element, kanal, okno.id, {
    persona,
    rolaOkna: okno.windowRole,
    modul: okno.moduleId,
    politykaModulu: politykaOknaModulu,
    pasekZlecenia: pasekZlecenia.element,
  });

  // Port widoku transkryptu łączy menu w nagłówku gniazda z rozmową — tryby są własnością rozmowy.
  gniazdo.ustawWidokZapisu({
    tryby: WIDOKI_ZAPISU,
    biezacy: () => rozmowa.widokZapisu(),
    ustaw: (kod) => rozmowa.ustawWidokZapisu(rozpoznajWidokZapisu(kod)),
  });

  // Plakietka stanu pętli w nagłówku gniazda ożywa z kolejki rdzenia.
  const odsubskrybujStanPary = zwiazStanPary(kanal, gniazdo, okno.id);

  const sterowanie = zamontujWidokSterowania({
    kanal,
    okno,
    rejestrKanalow,
    panel,
    akcjePaska: gniazdo.naglowek.element,
  });

  return {
    idOkna: okno.id,
    idSesji: okno.sessionId,
    rozmowa,
    sterowanie,

    ustawWidocznosc(widoczne) {
      sterowanie.element.hidden = !widoczne;
    },

    rozlacz() {
      // Port zabrany przed zdjęciem rozmowy — inaczej przeżyłby ją w menu i wołał tryby rozłączonej.
      gniazdo.ustawWidokZapisu(null);
      odsubskrybujStanPary();
      sterowanie.rozlacz();
      // Porty zdejmowane przed rozmową — inaczej ich subskrypcje przerysowywałyby stery zdjęte z dokumentu.
      zrodloSrodowiska.rozlacz();
      zrodloZlecenia.rozlacz();
      rozmowa.rozlacz();
    },
  };
}
