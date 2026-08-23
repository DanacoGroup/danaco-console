import { LICZBA_MAX } from '../../okna-rownolegle/identyfikatory';

/**
 * Sufit instancji Translation Panels — opis stanu, nie bramka.
 *
 * Górnej granicy liczby paneli nie stawia nikt i plik tego nie zmienia:
 * wpisanie progu z palca byłoby prawem wymyślonym po stronie klienta, a rdzeń
 * i tak przyjąłby panel ponad nim.
 *
 * Skąd wiadomo, że granicy nie ma:
 *  1. `store/migracja_053_tlumaczenie.sql` nie ma więzu na liczbę wierszy —
 *     `UNIQUE` stoi tylko na `identyfikator_zewnetrzny`, a `CHECK` tylko na
 *     kolumnie `stan`;
 *  2. `TranslateTargetAddRequest` niesie `windowId`, `language`, `tone`
 *     i `channelId` — ani jednego licznika, a `maxItems` nie pada w całym
 *     `shared/contract.ts`;
 *  3. `DodajPanel` (`adapter_modul_tlumaczenie_panele.go`) sprawdza `WindowId`,
 *     `Language` i istnienie okna, po czym woła model i zapisuje panel — bez
 *     licznika i bez odmowy przy jakiejkolwiek liczbie paneli okna;
 *  4. `.mt-panele` układa się jako
 *     `grid-template-columns: repeat(auto-fit, minmax(280px, 1fr))`, więc
 *     panele dokładają się bez końca.
 * Każdy panel kosztuje przy tym jedno wywołanie modelu: `DodajPanel` liczy
 * przekład przed zapisem, żeby odmowa modelu nie zostawiła pustego wiersza
 * w bazie. Rachunek rośnie więc liniowo z liczbą paneli.
 *
 * Dwie liczby z tego drzewa, które odpowiedzią na to pytanie nie są, i nota ma
 * je rozróżnić:
 *  — `LICZBA_MAX` (`okna-rownolegle/identyfikatory.ts`) to sufit gniazd sceny
 *    okien równoległych. Cały moduł Translate mieści się w jednym oknie,
 *    a panele stoją w jego siatce, więc gniazd sceny nie zużywają;
 *  — `granicaOkien` z profilu Translate (`okno-komunikacji/profile-wiedza.ts`)
 *    to granica okien komunikacji tego profilu. Okno komunikacji a panel języka
 *    to dwa różne byty.
 *
 * Forma idzie za `roundtable/sufit-uczestnikow.ts`. Sufit sceny jest
 * importowany, nie przepisywany: liczba jedzie z pliku, w którym stoi.
 */

/** Liczba gniazd sceny okien równoległych — sufit sceny, nie liczby paneli. */
export const GNIAZD_SCENY = LICZBA_MAX;

/** Granica okien komunikacji profilu Translate — także nie jest limitem paneli. */
export const OKIEN_KOMUNIKACJI_PROFILU = 2;

/** Stan sufitu widziany z liczby paneli, które okno zna w tej chwili. */
export interface StanSufituPaneli {
  /** Panele znane oknu w tej chwili. */
  panele: number;
  /** Zdania opisujące stan — w kolejności czytania. */
  zdania: string[];
}

/**
 * Składa opis stanu. Żadne ze zdań nie orzeka zakazu: mówią, ile paneli jest,
 * ile kosztowały i czego nie liczy rdzeń ani kontrakt.
 */
export function stanSufituPaneli(ilePaneli: number): StanSufituPaneli {
  const panele = Number.isFinite(ilePaneli) && ilePaneli > 0 ? Math.trunc(ilePaneli) : 0;

  const zdania = [
    panele === 0
      ? 'To okno nie zna dziś ani jednego panelu języka.'
      : `To okno zna dziś ${panele} ${odmianaPaneli(panele)} języka.`,
    'Liczby paneli nie ogranicza nikt: migracja 053 nie ma więzu na liczbę wierszy panelu, ' +
      'translate.target.add nie ma pola z granicą (w całym kontrakcie nie pada maxItems), ' +
      'a uchwyt DodajPanel sprawdza okno i język, po czym zapisuje panel bez licznika.',
    'KAŻDY PANEL KOSZTUJE JEDNO WYWOŁANIE MODELU: translate.target.add przekłada tekst źródłowy ' +
      'przed zapisem panelu, więc rachunek rośnie liniowo z liczbą języków, a nie z liczbą zapisów źródła.',
    `Sufit czterech należy do czego innego: ${GNIAZD_SCENY} gniazda ma scena okien równoległych ` +
      '(okna-rownolegle/identyfikatory.ts → LICZBA_MAX), a cały ten moduł mieści się w jednym oknie ' +
      `i gniazd sceny nie zużywa. Granica ${OKIEN_KOMUNIKACJI_PROFILU} z profilu Translate liczy OKNA ` +
      'KOMUNIKACJI, nie panele języka — to trzecia różna rzecz.',
    'Czy Translation Panels mają dostać górny próg liczbowy — nierozstrzygnięte. ' +
      'Do rozstrzygnięcia okno nie zgaduje i nie odbiera dodania panelu.',
  ];

  return { panele, zdania };
}

/**
 * Nota o suficie w miejscu, w którym dokłada się języki.
 *
 * Nota jest akapitem, nie ostrzeżeniem i nie bramką — nie odbiera żadnej
 * czynności i nie wygasza przycisku „+ Dodaj język". `data-panele` daje
 * sprawdzianowi odczytać liczbę bez parsowania zdania.
 */
export function utworzNoteSufituPaneli(ilePaneli: number): HTMLElement {
  const nota = document.createElement('div');
  nota.className = 'mt-sufit';
  odswiezNoteSufituPaneli(nota, ilePaneli);
  return nota;
}

/** Przepisuje notę bez wymiany elementu — liczba paneli zmienia się co dodanie. */
export function odswiezNoteSufituPaneli(nota: HTMLElement, ilePaneli: number): void {
  const stan = stanSufituPaneli(ilePaneli);
  nota.dataset['panele'] = String(stan.panele);

  const tytul = document.createElement('span');
  tytul.className = 'mt-sufit__tytul';
  tytul.textContent = 'Ile paneli mieści to okno';

  const zdania = stan.zdania.map((zdanie) => {
    const akapit = document.createElement('p');
    akapit.className = 'mt-sufit__zdanie';
    akapit.textContent = zdanie;
    return akapit;
  });

  nota.replaceChildren(tytul, ...zdania);
}

/**
 * Odmiana rzeczownika przy liczbie — polszczyzna, nie „2 paneli".
 * Rzeczownik męski nieżywotny: 1 panel, 2-4 panele, 5+ paneli, a nastoletnie
 * (12, 13, 14) wracają do dopełniacza mnogiego.
 */
function odmianaPaneli(ile: number): string {
  if (ile === 1) return 'panel';
  const jednosci = ile % 10;
  const dziesiatki = ile % 100;
  if (jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14)) return 'panele';
  return 'paneli';
}
