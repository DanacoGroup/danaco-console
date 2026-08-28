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

// Pokrycie komend jest jednym źródłem prawdy o tym, czego rdzeń nie obsługuje.

/** Pozycje bez czynności wraz z odczytem ich pokrycia w rdzeniu, gotowe do wstawienia w dowolne okno modułu. */
export interface PokrycieKomend {
  // Przycisk pozycji, której okno nie wykonuje; po naciśnięciu nazywa powód, zamiast milczeć.
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

/** Przycisk pozycji bez czynności; powód czytany z tytułu kontrolki w chwili kliknięcia przez czytelnika. */
function zbudujPrzycisk(
  kanal: Kanal,
  podepnij: (przerysuj: () => void) => void,
  etykieta: string,
  komendy: readonly string[],
  czynnosc: string,
): HTMLButtonElement {
  const kontrolka = przyciskAkcji(etykieta, 'dn-btn dn-btn--zarys');
  kontrolka.dataset['brakKomendy'] = 'tak';
  // Domknięcie na powodzie z chwili budowy mówiłoby o odczycie w toku długo po odpowiedzi rdzenia.
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

/** Wykaz pokrycia komend okna; klasa rodziny modułu przychodzi z zewnątrz jako parametr tego wywołania. */
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

/** Wiersz wykazu: nazwa komendy stoi w danych elementu, a powód wprost w treści widocznej dla czytelnika. */
function wierszWykazu(kanal: Kanal, komenda: string): HTMLElement {
  const wiersz = document.createElement('li');
  wiersz.dataset['komenda'] = komenda;
  wiersz.dataset['pokrycie'] = stanKomendy(kanal, komenda);
  wiersz.textContent = powodKomendy(kanal, komenda);
  return wiersz;
}
