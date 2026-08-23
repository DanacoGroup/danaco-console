import type { Window } from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import type { KomunikatZmiany } from '../sterowanie/komunikat-zmiany';
import { utworzStanSterowania } from '../sterowanie/stan-sterowania';
import { utworzZmianeOkna } from '../sterowanie/zmiana-okna';
import { utworzZmianeUstawienia } from '../sterowanie/zmiana-ustawienia';
import type { StanPrzelacznika, ZaleznosciPrzelacznika } from './przelacznik-srodowiska';

/**
 * Port przełącznika środowiska do kompletu sterowania okna.
 *
 * Przełącznik (`przelacznik-srodowiska.ts`) nie zna kanału ani stanu globalnego —
 * żąda migawki, subskrypcji zmian i wysyłki. Ten plik podpina te trzy rzeczy pod
 * stan rdzenia, tak samo jak `katalog-akcji.ts` zamienia kanał na źródło pozycji
 * panelu akcji.
 *
 * Zasięg wykonania ma w kliencie dwa widoki: listę wyboru w szufladzie ustawień
 * (`sterowanie/srodowisko-wykonania.ts`, pokazuje nazwy wartości wyliczenia) oraz
 * ten przełącznik nad polem wypowiedzi (pokazuje maszyny). Reguła protokołu nie
 * jest tu pisana po raz drugi: odczyt idzie przez `StanSterowania`, zapis przez
 * `ZmianaOkna` — te same byty katalogu `sterowanie/`, ta sama komenda
 * `window.update` z tym samym `windowId`, ta sama nazwa zmiany w komunikacie.
 * Oba widoki przyjmują wyłącznie stan potwierdzony przez rdzeń: odpowiedź na
 * komendę oraz zdarzenia `window.changed` i `config.changed` rozgłaszane do
 * wszystkich urządzeń konta.
 *
 * Nazwa maszyny zdalnej jest ustawieniem poziomu okna `host_wykonania`, a nie
 * polem `window.update` — dlatego port wczytuje ustawienia zasięgu okna. Bez tego
 * przełącznik pokazywałby „nie wskazano" przy hoście zapisanym w bazie.
 */
export interface ZrodloSrodowiska extends ZaleznosciPrzelacznika {
  /** Odłącza subskrypcje kanału i stanu. */
  rozlacz(): void;
}

/**
 * Nazwa zmiany w komunikacie o jej losie — ta sama, którą wysyła lista wyboru
 * w szufladzie (`sterowanie/srodowisko-wykonania.ts`). Odmowa rdzenia dociera
 * więc do operatora tym samym zdaniem, niezależnie od tego, którym z dwóch
 * widoków jej dotknął.
 */
const NAZWA = 'Środowisko wykonania';

export function utworzZrodloSrodowiska(kanal: Kanal, okno: Window): ZrodloSrodowiska {
  const stan = utworzStanSterowania(okno);

  /**
   * Kolejka wywołań czekających na los swojej zmiany.
   *
   * `ZmianaOkna.zastosuj` melduje wynik wywołaniem zwrotnym, a przełącznik żąda
   * obietnicy — kolejka jest całym przekładem między jednym a drugim. Odbiorca
   * kompletu dostaje dokładnie jeden komunikat na jedną wysyłkę, więc kolejka
   * nie rośnie: każdy meldunek zdejmuje z niej najstarsze oczekiwanie.
   */
  const oczekujace: ((komunikat: KomunikatZmiany) => void)[] = [];

  const zmiana = utworzZmianeOkna(kanal, stan, (komunikat) => {
    oczekujace.shift()?.(komunikat);
  });

  // Ustawienia poziomu okna czyta port, ale ich nie zapisuje — nazwę hosta
  // nadaje pole w szufladzie sterowania. Meldunek o losie zapisu nie ma tu
  // odbiorcy, bo ten port nigdy nie woła `zapisz`.
  const ustawienia = utworzZmianeUstawienia(kanal, stan, () => undefined);
  ustawienia.wczytaj();

  const odsubskrybuj: Odsubskrybuj[] = [];

  return {
    migawka(): StanPrzelacznika {
      const biezaca = stan.migawka();
      return {
        srodowisko: biezaca.okno.executionEnv,
        host: biezaca.ustawienia.hostWykonania,
      };
    },

    naZmiane(sluchacz) {
      odsubskrybuj.push(stan.naZmiane(() => sluchacz()));
    },

    /**
     * Odmowa rdzenia dochodzi dosłownie. Komunikat niepowodzenia składa
     * `sterowanie/komunikat-zmiany.ts` — kod i zdanie błędu wprost z kontraktu,
     * bez parafrazy — i tym samym zdaniem odzywa się pasek komunikatów kompletu
     * przy zmianie z szuflady. Obietnica odrzucona tym zdaniem trafia do pola
     * odmowy przełącznika.
     */
    zastosuj(srodowisko) {
      return new Promise<void>((spelnij, odrzuc) => {
        oczekujace.push((komunikat) => {
          if (komunikat.udany) spelnij();
          else odrzuc(new Error(komunikat.tresc));
        });
        zmiana.zastosuj(NAZWA, { executionEnv: srodowisko });
      });
    },

    rozlacz() {
      for (const zdejmij of odsubskrybuj.splice(0)) zdejmij();
      zmiana.rozlacz();
      ustawienia.rozlacz();
    },
  };
}
