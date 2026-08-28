/**
 * Cztery warstwy widoczności modułu Translate rozstrzygają, co stoi na ekranie bez interakcji,
 * a co dopiero po wywołaniu; nośnikiem rozwinięcia jest element details.
 */

/** Warstwa widoczności elementu jest podziałem z rozdziału o warstwach modułu, przyjmującym cztery wartości liczbowe. */
export type WarstwaWidocznosci = 1 | 2 | 3 | 4;

/** Nazwa warstwy i sposób dostępu do jej elementów opisują, jak operator dociera do funkcji tej warstwy w oknie. */
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
    dostep: 'Znacznik kontekstowy albo selektor ▼; element zwija się po użyciu.',
  },
  3: {
    nazwa: 'rozwinięcie kontekstowe',
    dostep: 'Menu ⋮ przy elemencie albo menu ☰ obszaru roboczego.',
  },
  4: {
    nazwa: 'funkcje eksperckie',
    dostep:
      'Skrót klawiszowy, wyszukiwarka funkcji modułu albo okno konfiguracji; ' +
      'polecenie języka naturalnego w Chat Window prowadzi do tego samego miejsca.',
  },
};

/** Znakuje element warstwą, żeby arkusz stylów i sprawdzian mogli pytać o warstwę elementu przez atrybut danych. */
export function oznaczWarstwe(element: HTMLElement, warstwa: WarstwaWidocznosci): void {
  element.dataset['warstwa'] = String(warstwa);
}

/** Element warstwy od drugiej do czwartej ma zapowiedź zawsze widoczną, a treść pokazuje się dopiero na wywołanie. */
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
  /** Znacznik wywołania z wykazu ikonografii: `⋮`, `☰` albo `▼`. */
  znacznik: string;
}

export function utworzRozwiniecie(opis: OpisRozwiniecia): Rozwiniecie {
  const element = document.createElement('details');
  element.className = 'mt-rozwiniecie';
  oznaczWarstwe(element, opis.warstwa);

  const znacznik = document.createElement('span');
  znacznik.className = 'mt-rozwiniecie__znacznik';
  znacznik.textContent = opis.znacznik;
  // Znacznik jest ozdobą uchwytu, nie jego nazwą: czytnik ekranu odczyta nazwę elementu, nie kropki.
  znacznik.setAttribute('aria-hidden', 'true');

  const nazwa = document.createElement('span');
  nazwa.className = 'mt-rozwiniecie__nazwa';
  nazwa.textContent = opis.nazwa;

  const warstwa = document.createElement('span');
  warstwa.className = 'mt-rozwiniecie__warstwa';
  warstwa.textContent = WARSTWY[opis.warstwa].nazwa;

  const uchwyt = document.createElement('summary');
  uchwyt.className = 'mt-rozwiniecie__uchwyt';
  uchwyt.append(znacznik, nazwa, warstwa);

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'mt-rozwiniecie__wyjasnienie';
  wyjasnienie.textContent = `${opis.wyjasnienie} ${WARSTWY[opis.warstwa].dostep}`;

  const tresc = document.createElement('div');
  tresc.className = 'mt-rozwiniecie__tresc';

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
