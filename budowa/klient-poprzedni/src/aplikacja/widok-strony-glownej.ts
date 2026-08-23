import { utworzFormularzKomponentu } from '../strona-glowna/formularz-komponentu';
import { pozycjaZeSrodowiska } from '../strona-glowna/pozycje-srodowisk';
import type { Kanal } from '../protokol/kanal';
import type { TozsamoscKlienta } from '../protokol/tozsamosc-klienta';
import { zadajWejscieNaStroneGlowna } from '../protokol/wejscie-strony-glownej';
import {
  utworzStroneGlowna,
  wepnijKomponenty,
  wepnijModuly,
  wepnijSrodowiska,
  zasilSesjeStronyGlownej,
  type KodSrodowiska,
} from '../strona-glowna/indeks';
import { wykonajZamiarUstawien } from './akcje-ustawien';
import { utworzPasekAplikacji } from './pasek-aplikacji';
import type { WidokTrasy } from './router';
import { Trasa } from './trasy';

/** Zależności widoku strony głównej. */
export interface ZaleznosciStronyGlownej {
  /** Nazwa projektu na pasku aplikacji. */
  projekt: string;
  /** Kanał kontraktu — źródło sesji w tle i katalogu ustawień. */
  kanal: Kanal;
  /**
   * Tożsamość klienta z powitania (`uzgodnienie.klient`) — bez niej strefa
   * sesji w tle pokazuje wykaz bez czynności powrotu (`session.bind` żąda
   * `clientId` z powitania, nie tożsamości nadanej po raz drugi).
   */
  klient?: TozsamoscKlienta;
  /** Przejście na inną trasę. */
  naTrase(trasa: Trasa): void;
  /** Wejście do środowiska; drugi argument wskazuje pozycję nawigacji. */
  naSrodowisko(kod: KodSrodowiska, modul?: string): void;
}

/** Widok trasy „Centrum dowodzenia". */
export interface WidokStronyGlownej extends WidokTrasy {
  /** Oznacza środowisko, w którym trwa praca; `null` zdejmuje oznaczenie. */
  ustawSrodowiskoCzynne(kod: KodSrodowiska | null): void;
  /**
   * Pierwszy odczyt strony domknięty — wykaz środowisk stoi albo rdzeń odmówił.
   *
   * Obietnica spełnia się w obu przypadkach, bo jest sygnałem „strona ma czym
   * stanąć", nie sygnałem powodzenia. Czeka na nią scena wejścia
   * (`ladowanie/`); odmowa, która nigdy nie domyka obietnicy, zostawiłaby scenę
   * nad gotowym produktem.
   *
   * Bez tożsamości klienta `home.enter` nie idzie i obietnica jest spełniona od
   * razu — stanowisko podglądu buduje stronę bez uzgodnienia i nie ma na co
   * czekać.
   */
  gotowa: Promise<void>;
}

/**
 * Widok trasy Centrum dowodzenia.
 *
 * Jedna odpowiedzialność: związanie strony głównej z resztą aplikacji.
 * Strona sama nie otwiera środowiska i nie wysyła komendy — zgłasza wybór,
 * a skutek należy do tej warstwy.
 *
 * Zmiana środowiska prowadzi przez tę stronę, dlatego jest ona trasą
 * początkową: uruchomienie aplikacji pokazuje przedpokój pracy, a nie okno
 * komunikacji wyrwane z kontekstu.
 */
export function utworzWidokStronyGlownej(
  zaleznosci: ZaleznosciStronyGlownej,
): WidokStronyGlownej {
  const strona = utworzStroneGlowna();

  // Wpięcie źródła danych: strefa sesji w tle żyje z `session.list` i zdarzeń
  // rdzenia, a powrót do sesji prowadzi przez to samo przejście do
  // środowiska, którym idzie wybór karty.
  zasilSesjeStronyGlownej({
    strona,
    kanal: zaleznosci.kanal,
    klient: zaleznosci.klient,
    naPrzejscie: (kod, modul) => zaleznosci.naSrodowisko(kod, modul),
  });

  /**
   * Wejście na stronę główną idzie komendą `home.enter`.
   *
   * Kontrakt daje nawigacji jeden ciąg: `home.enter` → `environment.list` →
   * `environment.enter` → `module.list` → `workspace.enter`.
   *
   * Ponad `environment.list` wejście daje: środowiska z kodami modułów
   * (`adapterNawigacji.StronaGlowna` woła `srodowiska(ctx, true)`), sesje
   * czynne konta, sesję ostatnio ogniskowaną na tym kliencie oraz żywy stan
   * sesji trwających w tle. `environment.list` bez `includeModules` nie niesie
   * żadnej z tych rzeczy i o żadną nie da się dopytać bez `clientId`.
   *
   * Sesje z odpowiedzi zostają nieużyte: strefa sesji ma źródło ciągłe —
   * `session.list`, zdarzenia rdzenia i archiwum (`zasilSesjeStronyGlownej`) —
   * a `home.enter` oddaje jedynie migawkę sesji czynnych. Zasilenie strefy
   * z obu naraz dałoby dwie rozjeżdżające się prawdy o tej samej rzeczy.
   *
   * Bez tożsamości klienta wejścia nie ma: `clientId` jest w żądaniu polem
   * obowiązkowym i musi pochodzić z powitania (drugie wywołanie
   * `tozsamoscKlienta()` nadałoby identyfikator nowy i rozdzieliło ognisko od
   * połączenia). Stanowisko podglądu buduje stronę bez uzgodnienia, więc dla
   * niego zostaje odczyt samego wykazu — brak tożsamości nie gasi ekranu.
   */
  const klient = zaleznosci.klient;
  let gotowa: Promise<void>;
  if (klient === undefined) {
    wepnijSrodowiska(strona, zaleznosci.kanal);
    gotowa = Promise.resolve();
  } else {
    gotowa = zadajWejscieNaStroneGlowna(zaleznosci.kanal, { clientId: klient.id }).then(
      (wynik) => {
        const srodowiska = wynik.wynik?.environments;
        // Odmowa nie gasi strefy: karty zastane, po których da się wejść do pracy,
        // są prawdziwsze niż pusty ekran.
        if (!wynik.udany || srodowiska === undefined || srodowiska.length === 0) return;
        strona.srodowiska.ustawWykaz(srodowiska.map(pozycjaZeSrodowiska));
      },
    );
  }

  const pasek = utworzPasekAplikacji({
    naUstawienie: (pozycja) => wykonajZamiarUstawien(pozycja, zaleznosci.kanal),
    projekt: zaleznosci.projekt,
    naTrase: zaleznosci.naTrase,
  });

  const element = document.createElement('div');
  element.className = 'dn-widok dn-widok--strona-glowna';
  element.append(pasek.element, strona.element);

  strona.naWyborSrodowiska((pozycja) => zaleznosci.naSrodowisko(pozycja.kod));

  // Kafel komponentu nie jest bramą — prowadzi tam, gdzie komponent pracuje,
  // a miejsce wskazuje macierz widoczności `srodowisko_modul`. Cel przychodzi
  // z `environment.list`, więc dopisanie modułu do środowiska w bazie zmienia
  // cel kafla bez zmiany po tej stronie.
  const komponenty = wepnijKomponenty({
    kanal: zaleznosci.kanal,
    strefa: strona.komponenty,
    naPrzejscie: (kod, modul) => zaleznosci.naSrodowisko(kod, modul),
  });

  strona.naWyborKomponentu((pozycja) => komponenty.wybierz(pozycja.kod));

  // Moduły bez pozycji w nawigacji mają na tej stronie swoją jedyną drogę:
  // rdzeń oddaje `automations` i `terminal` z pustym wykazem środowisk, więc
  // boczna nawigacja ich nie pokazuje. Cel przejścia wyznacza ta sama macierz
  // co dla kafli komponentów — `komponenty.wybierz` — żeby nie liczyć jej
  // drugi raz.
  wepnijModuly({
    kanal: zaleznosci.kanal,
    strefa: strona.moduly,
    naPrzejscie: (kodModulu) => komponenty.wybierz(kodModulu),
  });

  // Zakładanie komponentu ze strony głównej. Formularz wchodzi w przybornik
  // strefy komponentów, bo kanał ma warstwa widoku — strefa rdzenia nie zna.
  // Po założeniu wykaz dociąga się z rdzenia na nowo.
  const formularz = utworzFormularzKomponentu({
    kanal: zaleznosci.kanal,
    naZalozenie: () => komponenty.odswiez(),
  });
  strona.ustawienia.przybornik.append(formularz.element);
  strona.ustawienia.naDodanie(() => formularz.przelacz());

  // Listwa ustawień prowadzi do okna konfiguracji; dwie pozostałe pozycje
  // czekają na własne ekrany. Wykaz skutków mieszka w `akcje-ustawien`.
  strona.naWyborUstawienia((pozycja) => wykonajZamiarUstawien(pozycja, zaleznosci.kanal));

  return {
    element,
    przyWejsciu: () => pasek.trasy.ustawBiezaca(Trasa.StronaGlowna),
    ustawSrodowiskoCzynne: strona.ustawSrodowiskoCzynne,
    gotowa,
  };
}
