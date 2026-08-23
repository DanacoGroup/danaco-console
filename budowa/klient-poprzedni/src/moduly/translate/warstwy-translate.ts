/**
 * Cztery warstwy widoczności modułu Translate — jedno miejsce na regułę, która
 * rozstrzyga, co stoi na ekranie bez interakcji, a co dopiero po wywołaniu.
 *
 * Warstwa nie jest ozdobą opisu: rozstrzyga postać elementu przy wejściu do
 * modułu. Warstwa pierwsza jest rozwinięta i pozostaje taka; warstwy druga,
 * trzecia i czwarta stoją zwinięte, a zapowiedź nad nimi mówi, co jest pod
 * spodem — zwinięte nie znaczy ukryte.
 *
 * Nośnikiem rozwinięcia jest `<details>`: postać trzyma przeglądarka, więc
 * element działa klawiaturą i ma poprawną semantykę bez ani jednego nasłuchu.
 * Druga kopia stanu w klasie CSS mogłaby się z atrybutem `open` wyłącznie
 * rozminąć.
 *
 * Znacznik wywołania (`⋮`, `☰`, `▼`) idzie z wykazu ikonografii opracowania
 * i stoi przy nazwie, żeby droga do elementu była widoczna, zanim się go
 * otworzy.
 */

/** Warstwa widoczności elementu — podział z rozdziału o warstwach modułu. */
export type WarstwaWidocznosci = 1 | 2 | 3 | 4;

/** Nazwa warstwy i sposób dostępu do jej elementów. */
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

/** Znakuje element warstwą — arkusz i sprawdzian pytają o `data-warstwa`. */
export function oznaczWarstwe(element: HTMLElement, warstwa: WarstwaWidocznosci): void {
  element.dataset['warstwa'] = String(warstwa);
}

/** Element warstwy 2–4: zapowiedź zawsze widoczna, treść na wywołanie. */
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
  // Znacznik jest ozdobą uchwytu, nie jego nazwą: czytnik ekranu odczyta nazwę
  // elementu, a nie pionowe trzy kropki.
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
