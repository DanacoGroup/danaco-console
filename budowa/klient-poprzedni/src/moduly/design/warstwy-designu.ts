/**
 * Cztery warstwy widoczności modułu Design — jedno miejsce na regułę, która
 * rozstrzyga, co stoi na ekranie bez interakcji, a co dopiero po wywołaniu.
 * Funkcja niepotrzebna do bieżącego zadania nie jest widoczna.
 */

/**
 * Warstwa widoczności elementu. Numer warstwy jest jedyną nastawą, jaką element
 * dostaje w tej sprawie, a wszystko pozostałe — zwinięcie, znacznik i sposób
 * otwarcia — wynika z wykazu warstw.
 */
export type WarstwaWidocznosci = 1 | 2 | 3 | 4;

/**
 * Nazwa warstwy i sposób dostępu do jej elementów. Oba zdania są widoczne dla
 * Operatora, ponieważ warstwa zwinięta ma zapowiadać nie tylko swoją zawartość,
 * lecz także drogę, którą się ją otwiera.
 */
export interface OpisWarstwy {
  /** Nazwa warstwy widoczna dla Operatora. */
  readonly nazwa: string;
  /** Czym element tej warstwy się otwiera. */
  readonly dostep: string;
}

export const WARSTWY: Record<WarstwaWidocznosci, OpisWarstwy> = {
  1: {
    nazwa: 'zawsze widoczna',
    dostep: 'Bez interakcji — element stoi na ekranie od wejścia do modułu.',
  },
  2: {
    nazwa: 'widoczna na żądanie',
    dostep: 'Znacznik kontekstowy ▼ albo przycisk okna; element zwija się po użyciu.',
  },
  3: {
    nazwa: 'rozwinięcie kontekstowe',
    dostep: 'Menu ⋮ przy elemencie, menu ☰ kanwy albo menu kontekstowe zaznaczenia.',
  },
  4: {
    nazwa: 'funkcje eksperckie',
    dostep:
      'Skrót klawiszowy albo wyszukiwarka funkcji modułu; polecenie języka naturalnego ' +
      'w Chat Window prowadzi do tego samego miejsca.',
  },
};

/**
 * Znakuje element numerem warstwy. Arkusz stylów i sprawdzian pytają o atrybut
 * `data-warstwa`, więc jest on jedynym nośnikiem przynależności do warstwy i nie
 * dubluje się w nazwie klasy.
 */
export function oznaczWarstwe(element: HTMLElement, warstwa: WarstwaWidocznosci): void {
  element.dataset['warstwa'] = String(warstwa);
}

/**
 * Element warstwy zwiniętej: zapowiedź widoczna zawsze, treść dopiero na wywołanie.
 * Nośnikiem rozwinięcia jest znacznik `details`, którego postać trzyma przeglądarka,
 * więc element działa klawiaturą bez ani jednego nasłuchu.
 */
export interface Rozwiniecie {
  /** Element osadzany w oknie. */
  element: HTMLDetailsElement;
  /** Miejsce na treść — wypełnia je wywołujący. */
  tresc: HTMLElement;
  /** Rozwija element z zewnątrz — używa tego skrót klawiszowy. */
  rozwin(): void;
  /** Czy element jest w tej chwili rozwinięty. */
  rozwiniete(): boolean;
}

export interface OpisRozwiniecia {
  /** Warstwa, do której element należy. */
  warstwa: WarstwaWidocznosci;
  /** Nazwa elementu — ta sama, którą nosi w opracowaniu modułu. */
  nazwa: string;
  /** Jedno zdanie: po co element stoi w oknie. */
  wyjasnienie: string;
  /** Znacznik wywołania z opracowania: `⋮`, `☰` albo `▼`. */
  znacznik: string;
}

export function utworzRozwiniecie(opis: OpisRozwiniecia): Rozwiniecie {
  const element = document.createElement('details');
  element.className = 'md-rozwiniecie';
  oznaczWarstwe(element, opis.warstwa);

  const znacznik = document.createElement('span');
  znacznik.className = 'md-rozwiniecie__znacznik';
  znacznik.textContent = opis.znacznik;
  // Znacznik jest ozdobą uchwytu, nie jego nazwą; czytnik ekranu odczyta nazwę.
  znacznik.setAttribute('aria-hidden', 'true');

  const nazwa = document.createElement('span');
  nazwa.className = 'md-rozwiniecie__nazwa';
  nazwa.textContent = opis.nazwa;

  const warstwa = document.createElement('span');
  warstwa.className = 'md-rozwiniecie__warstwa';
  warstwa.textContent = WARSTWY[opis.warstwa].nazwa;

  const uchwyt = document.createElement('summary');
  uchwyt.className = 'md-rozwiniecie__uchwyt';
  uchwyt.append(znacznik, nazwa, warstwa);

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'md-rozwiniecie__wyjasnienie';
  wyjasnienie.textContent = `${opis.wyjasnienie} ${WARSTWY[opis.warstwa].dostep}`;

  const tresc = document.createElement('div');
  tresc.className = 'md-rozwiniecie__tresc';

  element.append(uchwyt, wyjasnienie, tresc);

  return {
    element,
    tresc,
    rozwin: () => {
      element.open = true;
    },
    rozwiniete: () => element.open,
  };
}
