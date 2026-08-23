import { pokazKomunikat } from '../aplikacja/komunikaty';
import { przyciskAkcji } from '../modele/kontrolki-formularza';
import type { Kanal } from '../protokol/kanal';
import {
  naglowekWykazu,
  podepnijDoWykazu,
  powodKomendy,
  stanKomendy,
  stanPozycji,
  zapewnijOdczyt,
  zapomnijWykazKomend,
  zdanieOPozycji,
  type StanPokrycia,
} from './wykaz-komend-rdzenia';

// Moduły mają jedno wejście: kto woła `utworzPokrycieKomend`, ten bierze stąd
// także unieważnienie wykazu, zdanie stanu nieustalonego i typ rozstrzygnięcia,
// zamiast sięgać po nie warstwę niżej.
export { POKRYCIE_W_ODCZYCIE } from './wykaz-komend-rdzenia';
export { zapomnijWykazKomend };
export type { StanPokrycia };

/**
 * Pokrycie komend — jedno źródło prawdy o tym, czego rdzeń nie obsługuje,
 * do użytku wszystkich modułów.
 *
 * Zdanie wpisane w moduł na sztywno („kontrakt nie ma komendy X") przestaje
 * być prawdą w dniu, w którym zmieni się rdzeń albo kontrakt, i nikt go nie
 * zdejmuje, bo nic go z rdzeniem nie łączy. Taki napis myli też stronę braku:
 * komenda bywa w kontrakcie, a nie ma uchwytu w złożonym rdzeniu. Pokrycie
 * bierze rozstrzygnięcie z odczytu, więc nazywa brak tam, gdzie jest.
 *
 * Byt stoi w korzeniu `moduly/`, a nie w module, bo tę samą potrzebę ma każdy
 * moduł; przepisywanie dałoby tyle samo rozjeżdżających się zdań o jednym
 * stanie produktu. Sam fakt i zdania o nim leżą warstwę niżej,
 * w `wykaz-komend-rdzenia.ts`.
 *
 * Kiedy tego nie używać: pozycja, dla której żadna komenda nie jest nawet
 * pomyślana (czynność wyłącznie okienna), zostaje przy `przyciskBezKomendy`
 * z `modele/kontrolki-formularza` — bez nazwy komendy nie ma czego sprawdzać
 * u rdzenia, a byt nie zgaduje.
 *
 * Użycie:
 *
 *   const pokrycie = utworzPokrycieKomend(kanal);          // raz na moduł
 *   panel.append(pokrycie.przycisk('Historia zatwierdzeń', 'developer.git.log',
 *     'Odczyt historii zatwierdzeń'));
 *   void pokrycie.odczytaj();                              // raz po montażu
 *   // w zamknij() okna:  pokrycie.zamknij();
 *
 * Powitanie idzie raz na połączenie, nie raz na moduł — wszystkie wywołania
 * `odczytaj()` czekają na jedną odpowiedź (`wykaz-komend-rdzenia.ts`).
 */

/** Pozycje bez czynności wraz z odczytem ich pokrycia w rdzeniu. */
export interface PokrycieKomend {
  /**
   * Przycisk pozycji, której okno nie wykonuje. Widoczny i w pełni klikalny,
   * bez `disabled`: po naciśnięciu nazywa powód, zamiast milczeć.
   *
   * @param etykieta nazwa pozycji w panelu akcji
   * @param komenda komenda, która tę pozycję by wykonała — nazwa z kontraktu
   *   albo z wykazu okien operacyjnych, także taka, której kontrakt nie ma
   * @param czynnosc czym ta pozycja jest dla czytającego — wchodzi w zdanie powodu
   */
  przycisk(etykieta: string, komenda: string, czynnosc: string): HTMLButtonElement;
  /** To samo dla pozycji, która potrzebuje kilku komend naraz. */
  przyciskWielu(etykieta: string, komendy: readonly string[], czynnosc: string): HTMLButtonElement;
  /** Wykaz pokrycia komend okna jako element informacyjny (`details`). */
  wykaz(komendy: readonly string[], klasa?: string): HTMLElement;
  /** Samo zdanie powodu — dla okien, które budują kontrolkę własną. */
  zdanie(komenda: string, czynnosc: string): string;
  /** Rozstrzygnięcie o komendzie bez budowania czegokolwiek. */
  stan(komenda: string): StanPokrycia;
  /** Przerysowanie własnej kontrolki po każdej zmianie wykazu. Woła się od razu. */
  naOdczyt(przerysuj: () => void): void;
  /** Pyta rdzeń o wykaz jego komend i wypełnia powody. Woła się raz, po montażu. */
  odczytaj(): Promise<void>;
  /** Odczyt wymuszony — po ponowieniu połączenia albo wymianie rdzenia. */
  odswiez(): Promise<void>;
  /** Odpina kontrolki tego modułu od wspólnego wykazu. Wołane z `zamknij()` okna. */
  zamknij(): void;
}

export function utworzPokrycieKomend(kanal: Kanal): PokrycieKomend {
  /** Kontrolki tego modułu; wspólny wykaz zna je jako jedno przerysowanie. */
  const moje: Array<() => void> = [];
  const odepnij = podepnijDoWykazu(kanal, () => {
    for (const przerysuj of moje) przerysuj();
  });

  /** Podpięcie kontrolki: przerysowanie teraz i po każdej zmianie wykazu. */
  function podepnij(przerysuj: () => void): void {
    moje.push(przerysuj);
    przerysuj();
  }

  return {
    przycisk: (etykieta, komenda, czynnosc) =>
      zbudujPrzycisk(kanal, podepnij, etykieta, [komenda], czynnosc),
    przyciskWielu: (etykieta, komendy, czynnosc) =>
      zbudujPrzycisk(kanal, podepnij, etykieta, komendy, czynnosc),
    wykaz: (komendy, klasa) => zbudujWykaz(kanal, podepnij, komendy, klasa),
    zdanie: (komenda, czynnosc) => zdanieOPozycji(kanal, [komenda], czynnosc),
    stan: (komenda) => stanKomendy(kanal, komenda),
    naOdczyt: podepnij,
    odczytaj: () => zapewnijOdczyt(kanal),
    odswiez: () => {
      zapomnijWykazKomend(kanal);
      return zapewnijOdczyt(kanal);
    },
    zamknij: () => {
      odepnij();
      moje.length = 0;
    },
  };
}

/** Przycisk pozycji bez czynności; powód czytany z `title` w chwili kliknięcia. */
function zbudujPrzycisk(
  kanal: Kanal,
  podepnij: (przerysuj: () => void) => void,
  etykieta: string,
  komendy: readonly string[],
  czynnosc: string,
): HTMLButtonElement {
  const kontrolka = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
  kontrolka.dataset['brakKomendy'] = 'tak';
  // Domknięcie na powodzie z chwili budowy mówiłoby „odczyt w toku” także długo
  // po odpowiedzi rdzenia — dlatego dymek bierze `title` dopiero przy kliknięciu.
  kontrolka.addEventListener('click', () => {
    pokazKomunikat({
      tytul: `${etykieta} — pokrycie w rdzeniu`,
      tresc: kontrolka.title,
      waga: 'ostrz',
    });
  });
  podepnij(() => {
    const powod = zdanieOPozycji(kanal, komendy, czynnosc);
    kontrolka.title = powod;
    kontrolka.setAttribute('aria-description', powod);
    kontrolka.dataset['pokrycie'] = stanPozycji(kanal, komendy);
  });
  return kontrolka;
}

/** Wykaz pokrycia komend okna; klasa rodziny modułu przychodzi z zewnątrz. */
function zbudujWykaz(
  kanal: Kanal,
  podepnij: (przerysuj: () => void) => void,
  komendy: readonly string[],
  klasa = '',
): HTMLElement {
  const element = document.createElement('details');
  if (klasa !== '') element.className = klasa;
  const naglowek = document.createElement('summary');
  const lista = document.createElement('ul');
  element.append(naglowek, lista);
  podepnij(() => {
    naglowek.textContent = naglowekWykazu(kanal, komendy);
    element.dataset['pokrycie'] = stanPozycji(kanal, komendy);
    lista.replaceChildren(...komendy.map((komenda) => wierszWykazu(kanal, komenda)));
  });
  return element;
}

/** Wiersz wykazu: nazwa komendy w danych, powód wprost w treści. */
function wierszWykazu(kanal: Kanal, komenda: string): HTMLElement {
  const wiersz = document.createElement('li');
  wiersz.dataset['komenda'] = komenda;
  wiersz.dataset['pokrycie'] = stanKomendy(kanal, komenda);
  wiersz.textContent = powodKomendy(kanal, komenda);
  return wiersz;
}
