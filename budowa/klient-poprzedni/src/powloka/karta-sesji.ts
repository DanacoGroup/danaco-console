import { ProgressStatus } from '../../../shared/contract';
import { ubierzCzynnosciKarty, zlozCzesciKarty } from './czesci-karty';

// Pojedyncza karta sesji w pasie zakładek — stan i cykl życia karty, budowanej z osobnych węzłów.

/** Opis karty potrzebny do jej zbudowania, dostarczany przez pas zakładek przy tworzeniu nowej karty sesji. */
export interface DaneKarty {
  id: string;
  tytul: string;
  stan: ProgressStatus;
}

/** Obsługa zdarzeń karty przekazywana przez pas zakładek, obejmująca zamknięcie oraz trwałe usunięcie sesji. */
export interface ObslugaKarty {
  przyWyborze(id: string): void;
  przyZamknieciu(id: string): void;
  /** Trwałe usunięcie sesji karty jest czynnością inną niż zamknięcie — dwa skutki, dwie różne nazwy. */
  przyUsunieciu(id: string): void;
  przyKlawiszu(id: string, zdarzenie: KeyboardEvent): void;
}

/** Karta sesji gotowa do osadzenia w pasie zakładek, niosąca element, sterowanie stanem oraz obsługę zdarzeń. */
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

/** Wygląd i etykieta stanu procesu karty sesji; stan nigdy nie idzie samą barwą kropki wskaźnika pracy. */
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
