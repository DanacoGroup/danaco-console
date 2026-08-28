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

/** Zależności widoku strony głównej: projekt, kanał, tożsamość klienta oraz przejścia do innych tras aplikacji. */
export interface ZaleznosciStronyGlownej {
  /** Nazwa projektu na pasku aplikacji. */
  projekt: string;
  /** Kanał kontraktu — źródło sesji w tle i katalogu ustawień. */
  kanal: Kanal;
  // Tożsamość klienta z powitania; bez niej strefa sesji w tle pokazuje wykaz bez czynności powrotu.
  klient?: TozsamoscKlienta;
  /** Przejście na inną trasę. */
  naTrase(trasa: Trasa): void;
  /** Wejście do środowiska; drugi argument wskazuje pozycję nawigacji. */
  naSrodowisko(kod: KodSrodowiska, modul?: string): void;
}

/** Widok trasy „Centrum dowodzenia" — strefy środowisk, komponentów, modułów oraz sesji trwających w tle. */
export interface WidokStronyGlownej extends WidokTrasy {
  /** Oznacza środowisko, w którym trwa praca; `null` zdejmuje oznaczenie. */
  ustawSrodowiskoCzynne(kod: KodSrodowiska | null): void;
  // Pierwszy odczyt strony domknięty — sygnał „strona ma czym stanąć", nie sygnał powodzenia.
  gotowa: Promise<void>;
}

/** Widok trasy Centrum dowodzenia, wiążący stronę główną z resztą aplikacji; strona zgłasza wybór, nie skutek. */
export function utworzWidokStronyGlownej(
  zaleznosci: ZaleznosciStronyGlownej,
): WidokStronyGlownej {
  const strona = utworzStroneGlowna();

  // Wpięcie źródła danych — strefa sesji w tle żyje z `session.list` i zdarzeń rdzenia.
  zasilSesjeStronyGlownej({
    strona,
    kanal: zaleznosci.kanal,
    klient: zaleznosci.klient,
    naPrzejscie: (kod, modul) => zaleznosci.naSrodowisko(kod, modul),
  });

  // Wejście na stronę główną idzie komendą `home.enter`, dającą jeden ciąg aż po `workspace.enter`.
  const klient = zaleznosci.klient;
  let gotowa: Promise<void>;
  if (klient === undefined) {
    wepnijSrodowiska(strona, zaleznosci.kanal);
    gotowa = Promise.resolve();
  } else {
    gotowa = zadajWejscieNaStroneGlowna(zaleznosci.kanal, { clientId: klient.id }).then(
      (wynik) => {
        const srodowiska = wynik.wynik?.environments;
        // Odmowa nie gasi strefy — karty zastane, po których da się wejść, są prawdziwsze niż pusty widok.
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

  // Kafel komponentu nie jest bramą — prowadzi tam, gdzie komponent pracuje wg macierzy widoczności.
  const komponenty = wepnijKomponenty({
    kanal: zaleznosci.kanal,
    strefa: strona.komponenty,
    naPrzejscie: (kod, modul) => zaleznosci.naSrodowisko(kod, modul),
  });

  strona.naWyborKomponentu((pozycja) => komponenty.wybierz(pozycja.kod));

  // Moduły bez pozycji w nawigacji mają tu jedyną drogę wejścia — przez `komponenty.wybierz`.
  wepnijModuly({
    kanal: zaleznosci.kanal,
    strefa: strona.moduly,
    naPrzejscie: (kodModulu) => komponenty.wybierz(kodModulu),
  });

  // Zakładanie komponentu ze strony głównej — formularz wchodzi w przybornik strefy komponentów.
  const formularz = utworzFormularzKomponentu({
    kanal: zaleznosci.kanal,
    naZalozenie: () => komponenty.odswiez(),
  });
  strona.ustawienia.przybornik.append(formularz.element);
  strona.ustawienia.naDodanie(() => formularz.przelacz());

  // Listwa ustawień prowadzi do okna konfiguracji — wykaz skutków mieszka w `akcje-ustawien`.
  strona.naWyborUstawienia((pozycja) => wykonajZamiarUstawien(pozycja, zaleznosci.kanal));

  return {
    element,
    przyWejsciu: () => pasek.trasy.ustawBiezaca(Trasa.StronaGlowna),
    ustawSrodowiskoCzynne: strona.ustawSrodowiskoCzynne,
    gotowa,
  };
}
