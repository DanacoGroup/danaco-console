/**
 * Port kolejki decyzji — miejsce wpięcia wykazu przepływów wstrzymanych z pytaniem
 * do człowieka. Kontrakt nie niesie ani komendy odczytu takiej kolejki, ani zdarzenia,
 * ani struktury, więc port stoi pusty i czeka na źródło.
 */

import { Droga, type Kwit } from './kwit-decyzji';
import type { PozycjaDecyzji } from './pozycje-decyzji';

/**
 * Odpowiedź portu — wykaz albo nazwany brak źródła. Oba stany wracają tą samą drogą,
 * więc ekran nie musi rozróżniać wywołania udanego od nieudanego, żeby wiedzieć,
 * czy wolno mu pytać o pozycje.
 */
export interface WynikKolejkiDecyzji {
  /** Czy źródło w ogóle istnieje. `false` znaczy: nie pytaj o pozycje. */
  dostepna: boolean;
  pozycje: readonly PozycjaDecyzji[];
  /** Zdanie dla Operatora, gdy źródła nie ma. Wypełniane wyłącznie przy `false`. */
  powodBraku?: string;
}

/**
 * Rozstrzygnięcie jednej pozycji eskalacji podjęte z telefonu. Niesie wskazanie
 * pozycji, samą decyzję i nieobowiązkowe słowo, które idzie do rdzenia razem
 * z rozstrzygnięciem.
 */
export interface RozstrzygniecieEskalacji {
  idPozycji: string;
  /** Co Operator postanowił: puścić dalej czy odrzucić. */
  decyzja: 'zatwierdz' | 'odrzuc';
  /** Słowo Operatora doklejane do rozstrzygnięcia; puste jest dopuszczalne. */
  uwaga?: string;
}

/**
 * Kontrakt portu widziany przez warstwę mobilną. Źródło podaje własną nazwę, odczyt
 * wykazu i drogę rozstrzygnięcia, więc ekran pracuje tak samo niezależnie od tego,
 * co zostało wpięte.
 */
export interface KolejkaDecyzji {
  /** Nazwa źródła wypisywana na ekranie — Operator ma wiedzieć, kto to mówi. */
  nazwa: string;
  odczytaj(): Promise<WynikKolejkiDecyzji>;
  rozstrzygnij(rozstrzygniecie: RozstrzygniecieEskalacji): Promise<Kwit>;
}

/**
 * Zdanie o braku źródła — jedno, żeby ekran i sprawdzian mówiły tak samo. Trzymanie
 * go w jednym miejscu chroni przed rozejściem się treści widzianej przez człowieka
 * i treści sprawdzanej w próbie.
 */
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

/**
 * Wpina źródło kolejki decyzji. Wołane raz, przy montażu warstwy mobilnej, bo ekran
 * czyta port przy każdym odświeżeniu i podmiana źródła w trakcie pracy zmieniałaby
 * wykaz bez wiedzy człowieka.
 */
export function zainstalujKolejkeDecyzji(kolejka: KolejkaDecyzji): void {
  zainstalowana = kolejka;
}

/**
 * Zdejmuje wpięcie i przywraca brak źródła. Istnieje dla sprawdzianów, a nie dla
 * widoku: ekran nie ma czynności, która odpinałaby kolejkę decyzji w trakcie pracy
 * warstwy mobilnej.
 */
export function odepnijKolejkeDecyzji(): void {
  zainstalowana = BRAK_ZRODLA_KOLEJKI;
}

/**
 * Źródło kolejki decyzji obowiązujące teraz. Dopóki nikt nic nie wpiął, oddaje
 * implementację domyślną, która nazywa brak źródła zamiast udawać kolejkę pustą.
 */
export function kolejkaDecyzji(): KolejkaDecyzji {
  return zainstalowana;
}
