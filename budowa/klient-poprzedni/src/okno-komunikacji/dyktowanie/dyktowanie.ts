import { utworzMagistrale, type Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import { utworzDostarczenie } from './dostarczenie-nagrania';
import {
  utworzZrodloDostepnosci,
  type DostepnoscDyktowania,
  type ZrodloDostepnosci,
} from './dostepnosc-dyktowania';
import { opiszOdmoweNagrywania, utworzNagrywanie, type StanNagrywania } from './nagrywanie';
import { utworzPrzytrzymanie } from './przytrzymanie';
import { utworzTranskrypcje } from './transkrypcja-nagrania';
import { odczytajUrzadzenia, naZmianeUrzadzen, type WykazUrzadzen } from './urzadzenia-dzwieku';
import { wynikNieprzetworzony, type WynikDyktowania } from './wynik-dyktowania';

/**
 * Dyktowanie — jeden byt spinający sześć warstw ogniwa mowa→tekst.
 *
 * Warstw jest sześć: wykaz urządzeń, nagrywanie, gest przytrzymania,
 * dostępność, dostarczenie nagrania i transkrypcja. Pasek polecenia ma z nich
 * zrobić jeden przycisk z jednym menu, a nie sześć zależności do posklejania
 * u siebie. Ten plik jest szwem: pasek zna wyłącznie jego, a wnętrze może się
 * przestawiać bez ruszania widoku.
 *
 * Plik nie rysuje niczego — nie tworzy ikony, menu ani paska i nie zna klas CSS.
 * Ikona mikrofonu, wykaz urządzeń z ptaszkiem i przełącznik „Przytrzymaj, aby
 * nagrać" należą do obszaru polecenia; gdyby ten plik rysował własny przycisk,
 * w pasku stanęłyby dwa mikrofony.
 *
 * Mikrofon, którego nie ma czym obsłużyć, nie pojawia się wcale: `dostepnosc()`
 * jest pytaniem, które pasek zadaje zanim narysuje ikonę, nie po naciśnięciu.
 * Wyszarzony mikrofon albo mikrofon odmawiający po kliknięciu byłby bramą;
 * krótszy pasek nią nie jest.
 */
export interface Dyktowanie {
  /** Czy dyktowanie da się wykonać tu i teraz — pytanie przed narysowaniem ikony. */
  dostepnosc(): Promise<DostepnoscDyktowania>;
  /** Ponowne pytanie o dostępność, np. po zmianie ustawień silnika. */
  odswiezDostepnosc(): void;
  /** Wykaz mikrofonów do menu; niesie powód, gdy wykaz jest pusty albo niepełny. */
  urzadzenia(): Promise<WykazUrzadzen>;
  /** Urządzenie wskazane przez Operatora; puste znaczy domyślne systemu. */
  urzadzenie(): string;
  /** Wskazanie urządzenia z menu. */
  wybierzUrzadzenie(id: string): void;
  /** Stan nagrywania — pasek zmienia nim postać ikony. */
  stan(): StanNagrywania;
  /** Subskrypcja stanu nagrywania. */
  naStan(sluchacz: (stan: StanNagrywania) => void): Odsubskrybuj;
  /** Subskrypcja wyniku: rozpoznano tekst, mowy brak albo nie przetworzono. */
  naWynik(sluchacz: (wynik: WynikDyktowania) => void): Odsubskrybuj;
  /** Podpina gest przytrzymania do przycisku paska; oddaje odpięcie. */
  podepnijPrzycisk(cel: HTMLElement): Odsubskrybuj;
  /** Rozpoczyna nagrywanie — dla trybu zatrzaskowego, bez przytrzymania. */
  rozpocznij(): void;
  /** Kończy nagrywanie i uruchamia transkrypcję. */
  zakoncz(): void;
  /** Porzuca nagranie bez transkrypcji — Escape w trakcie. */
  porzuc(): void;
  /** Odpina subskrypcje sprzętu i zwalnia mikrofon, jeśli jeszcze pracuje. */
  rozlacz(): void;
}

export function utworzDyktowanie(kanal: Kanal): Dyktowanie {
  const stany = utworzMagistrale<StanNagrywania>();
  const wyniki = utworzMagistrale<WynikDyktowania>();

  const dostarczenie = utworzDostarczenie(kanal);
  const dostepnosc: ZrodloDostepnosci = utworzZrodloDostepnosci(kanal, dostarczenie);
  const transkrypcja = utworzTranskrypcje(kanal, dostarczenie);
  const nagrywanie = utworzNagrywanie();

  let wybrane = '';

  const odsubskrybujStan = nagrywanie.naStan((stan) => stany.oglos(stan));

  // Zmiana sprzętu w trakcie pracy jest zwykłą rzeczą — Operator wpina zestaw
  // słuchawkowy w środku dnia. Wykaz odświeża się sam, ale wskazanie zostaje:
  // przestawienie go za Operatora znaczyłoby, że mówi do innego mikrofonu, niż
  // wybrał, i dowiedziałby się o tym dopiero po pustej transkrypcji.
  const odsubskrybujSprzet = naZmianeUrzadzen(() => dostepnosc.odswiez());

  /**
   * Zakończenie nagrania i cała droga do tekstu.
   *
   * Odmowa na każdym kroku wychodzi tą samą magistralą co powodzenie — pasek
   * ma jedno miejsce nasłuchu i nie musi rozstrzygać, czy zawiodło nagrywanie,
   * dostarczenie czy silnik. Rozróżnienie niesie sam wynik.
   */
  function domknij(): void {
    void nagrywanie
      .zakoncz()
      .then(async (nagranie) => {
        // Puszczenie przycisku przed pierwszą próbką nie jest błędem i nie ma
        // czego meldować: Operator nacisnął i rozmyślił się. Cisza w odpowiedzi
        // na brak nagrania jest tu właściwa — wynik pusty byłby zdaniem o niczym.
        if (nagranie === null) return;
        const wynik = await transkrypcja.wykonaj(
          nagranie.bajty,
          nagranie.rodzajTresci,
        );
        wyniki.oglos(wynik);
      })
      .catch((blad: unknown) => {
        wyniki.oglos(wynikNieprzetworzony(opiszOdmoweNagrywania(blad)));
      });
  }

  function rozpocznij(): void {
    void nagrywanie.rozpocznij(wybrane).catch((blad: unknown) => {
      wyniki.oglos(wynikNieprzetworzony(opiszOdmoweNagrywania(blad)));
    });
  }

  const przytrzymanie = utworzPrzytrzymanie({
    naStart: rozpocznij,
    naKoniec: domknij,
    naPorzucenie: () => nagrywanie.przerwij(),
  });

  return {
    dostepnosc: () => dostepnosc.odczytaj(),
    odswiezDostepnosc: dostepnosc.odswiez,
    urzadzenia: odczytajUrzadzenia,
    urzadzenie: () => wybrane,

    wybierzUrzadzenie(id) {
      wybrane = id;
    },

    stan: nagrywanie.stan,
    naStan: (sluchacz) => stany.subskrybuj(sluchacz),
    naWynik: (sluchacz) => wyniki.subskrybuj(sluchacz),
    podepnijPrzycisk: (cel) => przytrzymanie.podepnij(cel),
    rozpocznij,
    zakoncz: domknij,
    porzuc: () => nagrywanie.przerwij(),

    rozlacz() {
      odsubskrybujStan();
      odsubskrybujSprzet();
      // Zwolnienie mikrofonu przy zejściu okna jest obowiązkowe, nie porządkowe:
      // niezwolniony strumień zostawia zapaloną lampkę mikrofonu, a Operator ma
      // prawo wiedzieć, kiedy go nie słychać.
      nagrywanie.przerwij();
    },
  };
}
