import { ProgressStatus } from '../../../shared/contract';
import { ubierzCzynnosciKarty, zlozCzesciKarty } from './czesci-karty';

/**
 * Pojedyncza karta sesji w pasie zakładek.
 *
 * Jedna odpowiedzialność: stan i cykl życia jednej karty. Budowa jej węzłów
 * należy do `czesci-karty.ts`, mechanika pasa — wybór, kolejność, wędrujący
 * fokus — do `karty-sesji.ts`.
 *
 * Karta niesie trzy rzeczy: wskaźnik pracy w tle, tytuł równy nazwie otwartego
 * modułu oraz zamknięcie; obok zamknięcia stoi usunięcie trwałe sesji, czynność
 * o innym skutku i innej etykiecie. Wskaźnik pracy pokazuje, że proces biegnie,
 * choć Operator patrzy gdzie indziej.
 *
 * Nazwy stanów pochodzą z `shared/contract` — powłoka nie zakłada
 * własnego słownika stanu procesu.
 */

/** Opis karty potrzebny do jej zbudowania. */
export interface DaneKarty {
  id: string;
  tytul: string;
  stan: ProgressStatus;
}

/** Obsługa zdarzeń karty przekazywana przez pas zakładek. */
export interface ObslugaKarty {
  przyWyborze(id: string): void;
  przyZamknieciu(id: string): void;
  /**
   * Trwałe usunięcie sesji karty (`session.delete`) — czynność inna niż
   * `przyZamknieciu`. Zamknięcie zmienia stan sesji i zostawia zapis; usunięcie
   * wyprowadza zapis do kosza rdzenia, skąd po terminie znika trwale.
   * Dwie nazwy, bo dwa skutki.
   */
  przyUsunieciu(id: string): void;
  przyKlawiszu(id: string, zdarzenie: KeyboardEvent): void;
}

/** Karta sesji gotowa do osadzenia w pasie. */
export interface KartaSesji {
  readonly id: string;
  readonly element: HTMLElement;
  tytul(): string;
  ustawTytul(nazwa: string): void;
  stan(): ProgressStatus;
  ustawStan(nowy: ProgressStatus): void;
  ustawCzynna(czynna: boolean): void;
  czyCzynna(): boolean;
  ustawOgnisko(): void;
}

/** Wygląd i etykieta stanu procesu. Stan nigdy nie idzie samą barwą. */
interface WygladStanu {
  /** Wariant kropki z biblioteki komponentów. */
  klasa: string;
  /** Etykieta czytana przez technologie wspomagające i widoczna w dymku. */
  etykieta: string;
  /** Czy stan wymaga widocznej plakietki obok tytułu. */
  wyrozniony: boolean;
}

const WYGLAD_STANU: Readonly<Record<ProgressStatus, WygladStanu>> = {
  [ProgressStatus.Pending]: { klasa: '', etykieta: 'Oczekuje', wyrozniony: false },
  [ProgressStatus.Running]: { klasa: 'dn-kropka--tetno', etykieta: 'Praca w tle', wyrozniony: false },
  [ProgressStatus.Paused]: { klasa: 'dn-kropka--ostrzezenie', etykieta: 'Wstrzymana', wyrozniony: true },
  [ProgressStatus.Stopped]: { klasa: 'dn-kropka--neutralna', etykieta: 'Zatrzymana', wyrozniony: true },
  [ProgressStatus.Done]: { klasa: 'dn-kropka--sukces', etykieta: 'Zakończona', wyrozniony: false },
  [ProgressStatus.Failed]: { klasa: 'dn-kropka--blad', etykieta: 'Błąd', wyrozniony: true },
};

export function utworzKarteSesji(dane: DaneKarty, obsluga: ObslugaKarty): KartaSesji {
  let stan = dane.stan;
  let tytul = dane.tytul;

  const czesci = zlozCzesciKarty(dane.id);
  const { element, kropka, napis, plakietka, usun, zamknij } = czesci;

  element.addEventListener('click', () => obsluga.przyWyborze(dane.id));
  element.addEventListener('keydown', (zdarzenie) => obsluga.przyKlawiszu(dane.id, zdarzenie));

  // Zamknięcie nie jest wyborem karty — zdarzenie zatrzymuje się na przycisku.
  zamknij.addEventListener('click', (zdarzenie) => {
    zdarzenie.stopPropagation();
    obsluga.przyZamknieciu(dane.id);
  });

  usun.addEventListener('click', (zdarzenie) => {
    zdarzenie.stopPropagation();
    obsluga.przyUsunieciu(dane.id);
  });

  ubierzTytul();
  ubierzStan();

  function ubierzTytul(): void {
    napis.textContent = tytul;
    element.title = `${tytul} · ${WYGLAD_STANU[stan].etykieta}`;
    ubierzCzynnosciKarty(czesci, tytul);
  }

  function ubierzStan(): void {
    const wyglad = WYGLAD_STANU[stan];
    kropka.className = ['dn-kropka', 'dn-sesja__kropka', wyglad.klasa].filter(Boolean).join(' ');
    kropka.setAttribute('aria-label', `Stan pracy: ${wyglad.etykieta}`);
    plakietka.textContent = wyglad.etykieta;
    plakietka.hidden = !wyglad.wyrozniony;
    element.dataset.stan = stan;
    element.title = `${tytul} · ${wyglad.etykieta}`;
  }

  return {
    id: dane.id,
    element,

    tytul: () => tytul,

    ustawTytul(nazwa) {
      tytul = nazwa;
      ubierzTytul();
    },

    stan: () => stan,

    ustawStan(nowy) {
      stan = nowy;
      ubierzStan();
    },

    ustawCzynna(czynna) {
      element.setAttribute('aria-selected', String(czynna));
      element.tabIndex = czynna ? 0 : -1;
    },

    czyCzynna: () => element.getAttribute('aria-selected') === 'true',

    ustawOgnisko() {
      element.focus();
    },
  };
}
