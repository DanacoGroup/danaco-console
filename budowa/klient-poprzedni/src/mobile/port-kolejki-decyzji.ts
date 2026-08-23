import { Droga, type Kwit } from './kwit-decyzji';
import type { PozycjaDecyzji } from './pozycje-decyzji';

/**
 * Port kolejki decyzji — miejsce wpięcia wykazu przepływów wstrzymanych
 * z pytaniem do człowieka.
 *
 * Kontrakt nie niesie ani komendy odczytu takiej kolejki, ani zdarzenia, ani
 * struktury, więc port stoi pusty i czeka na źródło.
 *
 * Implementacja domyślna nie udaje kolejki pustej. Pusty wykaz znaczy „nic nie
 * czeka” i jest zdaniem uspokajającym; bez źródła byłoby ono nieprawdziwe
 * dokładnie w tej sytuacji, dla której ekran powstał. Dlatego
 * `BRAK_ZRODLA_KOLEJKI` oddaje `dostepna: false` wraz z powodem, a ekran pisze
 * ten powód wprost — tak samo jak pulpit w `mission-control/pas-decyzji.ts`.
 *
 * Wpięcie źródła to jedno wywołanie `zainstalujKolejkeDecyzji` przy montażu
 * warstwy mobilnej. Ekran się przez nie nie zmienia: pozycje portu mają ten sam
 * kształt co pozycje własne (`PozycjaDecyzji`), więc wchodzą do tego samego
 * wykazu i tego samego arkusza dróg.
 */

/** Odpowiedź portu — wykaz albo nazwany brak źródła. */
export interface WynikKolejkiDecyzji {
  /** Czy źródło w ogóle istnieje. `false` znaczy: nie pytaj o pozycje. */
  dostepna: boolean;
  pozycje: readonly PozycjaDecyzji[];
  /** Zdanie dla Operatora, gdy źródła nie ma. Wypełniane wyłącznie przy `false`. */
  powodBraku?: string;
}

/** Rozstrzygnięcie jednej pozycji eskalacji podjęte z telefonu. */
export interface RozstrzygniecieEskalacji {
  idPozycji: string;
  /** Co Operator postanowił: puścić dalej czy odrzucić. */
  decyzja: 'zatwierdz' | 'odrzuc';
  /** Słowo Operatora doklejane do rozstrzygnięcia; puste jest dopuszczalne. */
  uwaga?: string;
}

/** Kontrakt portu widziany przez warstwę mobilną. */
export interface KolejkaDecyzji {
  /** Nazwa źródła wypisywana na ekranie — Operator ma wiedzieć, kto to mówi. */
  nazwa: string;
  odczytaj(): Promise<WynikKolejkiDecyzji>;
  rozstrzygnij(rozstrzygniecie: RozstrzygniecieEskalacji): Promise<Kwit>;
}

/** Zdanie o braku źródła — jedno, żeby ekran i sprawdzian mówiły tak samo. */
export const ZDANIE_BRAKU_KOLEJKI =
  'Kolejki eskalacji rdzeń dziś nie wystawia: kontrakt nie niesie ani komendy jej ' +
  'odczytu, ani zdarzenia. Wykaz poniżej pokazuje wyłącznie to, co rdzeń umie ' +
  'udowodnić: kolejki wstrzymane, pętle zatrzymane i procesy, które stanęły.';

/**
 * Implementacja domyślna — uczciwy brak źródła.
 *
 * `rozstrzygnij` istnieje, bo port jest interfejsem; ekran nie rysuje
 * rozstrzygnięcia, dopóki `dostepna` jest `false`. Wywołana mimo to, oddaje
 * kwit nieudany z nazwanym powodem, nigdy potwierdzenie.
 */
export const BRAK_ZRODLA_KOLEJKI: KolejkaDecyzji = {
  nazwa: 'brak źródła kolejki decyzji',

  odczytaj() {
    return Promise.resolve({ dostepna: false, pozycje: [], powodBraku: ZDANIE_BRAKU_KOLEJKI });
  },

  rozstrzygnij(rozstrzygniecie) {
    return Promise.resolve({
      id: `kwit-brak-zrodla-${rozstrzygniecie.idPozycji}`,
      droga: Droga.Zatwierdz,
      pozycja: rozstrzygniecie.idPozycji,
      kroki: [],
      udany: false,
      zdanie:
        'Rozstrzygnięcia nie wysłano: kolejki eskalacji nie ma w kontrakcie, ' +
        'więc nie ma komendy, którą decyzja miałaby dojechać do rdzenia.',
      o: Date.now(),
    });
  },
};

let zainstalowana: KolejkaDecyzji = BRAK_ZRODLA_KOLEJKI;

/** Wpina źródło kolejki decyzji. Wołane raz, przy montażu. */
export function zainstalujKolejkeDecyzji(kolejka: KolejkaDecyzji): void {
  zainstalowana = kolejka;
}

/** Zdejmuje wpięcie — istnieje dla sprawdzianów, nie dla widoku. */
export function odepnijKolejkeDecyzji(): void {
  zainstalowana = BRAK_ZRODLA_KOLEJKI;
}

/** Źródło kolejki decyzji obowiązujące teraz. */
export function kolejkaDecyzji(): KolejkaDecyzji {
  return zainstalowana;
}
