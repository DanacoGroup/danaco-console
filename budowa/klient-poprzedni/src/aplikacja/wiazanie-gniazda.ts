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

/** Zależności powiązania gniazda sceny z oknem otwartym w rdzeniu. */
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

/** Gniazdo związane z oknem rdzenia: rozmowa i komplet sterowania. */
export interface WiazanieGniazda {
  /** Identyfikator okna nadany przez rdzeń. */
  idOkna: string;
  /**
   * Sesja, do której należy okno.
   *
   * Okno jest bytem podrzędnym wobec sesji i należy do dokładnie jednej,
   * więc wiązanie niesie ją ze sobą. Bez tego scena nie umiałaby powiedzieć,
   * czyje okna pokazuje — a po powiązaniu połączenia z inną sesją
   * (`session.bind`) to przestaje być oczywiste.
   */
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

/**
 * Powiązanie jednego gniazda sceny z jednym oknem rdzenia.
 *
 * Jedna odpowiedzialność: doprowadzenie do gniazda dwóch warstw, które okno
 * czynią użytecznym — rozmowy i kompletu sterowania. Ten plik nie buduje ani
 * jednego pola, ani jednej kontrolki: bierze gotowe z `rozmowa/`
 * i `widok-sterowania/`.
 *
 * Rozmowa wstawia się wprost do gniazda: na scenie sesji gniazdo nie buduje
 * własnego wbudowanego okna komunikacji
 * (`GniazdoOkna.OpcjeGniazda.wbudowanaRozmowa` zostaje wyłączone — patrz
 * komentarz przy tym ustawieniu w `okna-rownolegle/gniazdo-okna.ts`), więc
 * miejsce pod nagłówkiem gniazda ma dokładnie jednego mieszkańca: tury rozmowy
 * z modelem prowadzone komendą `message.send`, z podziałem na nadawców,
 * prowenancją i podsumowaniem tury.
 *
 * `zamontujWidokSterowania` wymaga miejsca „akcji paska". Sceną sesji rządzi
 * kilka okien naraz, więc jednakowe uchwyty na wspólnym pasku górnym byłyby
 * nie do rozróżnienia; uchwyt trafia do nagłówka własnego gniazda — obok
 * numeru, roli i stanu pętli.
 */
export function zwiazGniazdoZOknem(zaleznosci: ZaleznosciWiazania): WiazanieGniazda {
  const { kanal, rejestrKanalow, gniazdo, okno, persona, panel } = zaleznosci;

  // Pamięć rozmowy idzie za modułem okna: moduł bez pamięci sesyjnej — dziś
  // Agents, którego czat jest środowiskiem testowania agenta — nie odtwarza
  // wątku z rdzenia i czyści go przy zmianie testowanego eksperta. Polityka
  // wchodzi tutaj, bo tylko powłoka widzi naraz okno sceny i widok modułu;
  // składa ją `ulotnosc-okna`, a rozstrzyga profil modułu, nie ten plik.
  //
  // Przełącznik środowiska wchodzi tu, bo to jedyne miejsce mające oba końce:
  // wybór maszyny dokonuje się w oknie komunikacji jednym przełącznikiem,
  // a przełącznik żąda kompletu sterowania okna — rozmowa go nie widzi, komplet
  // nie widzi rozmowy. Kontrolki ten plik nie buduje: bierze gotową
  // z `okno-komunikacji/` i gotowy port stanu z `zrodlo-srodowiska`,
  // a rozmowie podaje wyłącznie element do postawienia nad polem wypowiedzi.
  const zrodloSrodowiska = utworzZrodloSrodowiska(kanal, okno);
  const przelacznikSrodowiska = utworzPrzelacznikSrodowiska(zrodloSrodowiska);

  // Pasek zlecenia wchodzi tu z tego samego powodu, co przełącznik: to jedyne
  // miejsce mające oba końce. Stery katalogu roboczego, modelu, wysiłku i trybu
  // zatwierdzania żądają kompletu sterowania okna, a rozmowa
  // go nie widzi. Pasek nie buduje ani jednego menu (robi to
  // `komponenty/menu-drzewo.ts`) i nie zna kontraktu (zna go
  // `okno-komunikacji/zrodlo-zlecenia.ts`); ten plik podaje mu tylko port,
  // rejestr kanałów i drogę do kolumny sterowania.
  const zrodloZlecenia = utworzZrodloZlecenia(kanal, okno);
  const pasekZlecenia = utworzPasekZlecenia({
    kanal,
    zrodlo: zrodloZlecenia,
    rejestrKanalow,
    sterowanieSrodowiska: przelacznikSrodowiska.element,
    // Stopka steru katalogów prowadzi do kolumny sterowania — jedynego miejsca
    // z polem ścieżki. Domknięcie, a nie odwołanie wprost: kolumna powstaje
    // niżej w tym pliku, a stopka woła to dopiero pod ręką Operatora.
    otworzSterowanie: () => sterowanie.otworz(),
  });

  const rozmowa = zamontujRozmowe(gniazdo.element, kanal, okno.id, {
    persona,
    rolaOkna: okno.windowRole,
    modul: okno.moduleId,
    politykaModulu: politykaOknaModulu,
    pasekZlecenia: pasekZlecenia.element,
  });

  // Port widoku transkryptu łączy menu `⋮` w nagłówku gniazda z rozmową: menu
  // pokazuje tryby zapisu i je przestawia, a tryby są własnością rozmowy, nie
  // układu okien. Wykaz idzie z `rozmowa/`, odczyt i zapis to wywołania
  // `ZamontowanaRozmowa`; `rozpoznajWidokZapisu` jest przekładem napisu na kod
  // trybu z tamtego katalogu, a nie drugim wykazem. Port podpinamy dopiero
  // tutaj, bo dopiero tutaj rozmowa istnieje.
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
      // Port zabrany przed zdjęciem rozmowy: wiersz `Widok transkryptu ›`
      // przeżyłby ją w menu gniazda, które zostaje na scenie, i wołałby tryby
      // rozmowy już rozłączonej.
      gniazdo.ustawWidokZapisu(null);
      odsubskrybujStanPary();
      sterowanie.rozlacz();
      // Porty zdejmowane przed rozmową: ich subskrypcje `window.changed`
      // i `config.changed` przeżyłyby widok, w którym stoi pasek zlecenia,
      // i przerysowywałyby stery już zdjęte z dokumentu.
      zrodloSrodowiska.rozlacz();
      zrodloZlecenia.rozlacz();
      rozmowa.rozlacz();
    },
  };
}
