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
 * Dyktowanie łączy sześć warstw obsługi mowy — urządzenia, nagrywanie, przytrzymanie, dostępność, dostarczenie i transkrypcję — w jeden interfejs dla paska poleceń, bez rysowania własnego widoku.
 */
export interface Dyktowanie {
  /** Czy dyktowanie da się wykonać tu i teraz — pytanie przed narysowaniem ikony. */
  dostepnosc(): Promise<DostepnoscDyktowania>;
  /** Ponowne pytanie o dostępność, na przykład po zmianie ustawień silnika. */
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
  /** Rozpoczyna nagrywanie — dla trybu przełącznika, bez przytrzymania. */
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

  // Wskazanie urządzenia przetrwa odświeżenie wykazu, by nie przełączyć mikrofonu bez wiedzy Operatora.
  const odsubskrybujSprzet = naZmianeUrzadzen(() => dostepnosc.odswiez());

  /** Zakończenie nagrania: odmowa i powodzenie idą tą samą magistralą, pasek ma jeden punkt nasłuchu. */
  function domknij(): void {
    void nagrywanie
      .zakoncz()
      .then(async (nagranie) => {
        // Puszczenie przycisku bez próbki nie jest błędem — brak nagrania nie generuje wyniku ani zgłoszenia.
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
      // Zwolnienie mikrofonu przy zamknięciu okna jest obowiązkowe, aby nie została zapalona zbędna lampka.
      nagrywanie.przerwij();
    },
  };
}
