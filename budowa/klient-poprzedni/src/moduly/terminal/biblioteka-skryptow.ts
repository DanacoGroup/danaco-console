import { TerminalShell } from '../../../../shared/contract';

/**
 * Biblioteka skryptów okna Script Library — skrypty, snippety, ich parametry
 * i kontrola treści przed uruchomieniem.
 *
 * Biblioteka jest WIDOKIEM na dziennik rdzenia: wykaz czyta `terminal.script.list`,
 * a zapis idzie przez `terminal.script.save`, który zakłada kolejną wersję
 * pozycji. Ten byt trzyma pozycje w pamięci wyłącznie po to, żeby je pokazać
 * i sprawdzić przed uruchomieniem — po odświeżeniu strony wykaz wraca z rdzenia.
 * Wywóz do pliku zostaje jako droga wyniesienia treści POZA rdzeń.
 *
 * Uruchomienie idzie tą samą drogą co każde polecenie karty: rdzeń podaje treść
 * jako pojedynczy argument programowi powłoki (`bash -c`, `powershell -Command`,
 * `python -u -c`, `node -e`), więc skrypt wieloliniowy wykonuje się bez
 * zapisywania go do pliku. Stąd wynika warunek, o którym okno przypomina:
 * powłoka karty musi być tą, w której skrypt napisano — treść Basha podana
 * programowi `python` nie jest skryptem, tylko błędem składni.
 *
 * Parametry deklaruje się nagłówkiem w komentarzu na początku treści:
 *
 *   # args:
 *   #   srodowisko: nazwa środowiska wdrożenia
 *   #   wersja: numer wydania
 *
 * a w treści stoją jako `{{srodowisko}}`. Podstawienie dzieje się w kliencie
 * przed wysłaniem, więc Operator widzi w potwierdzeniu polecenie, które
 * naprawdę poszło do rdzenia — nie szablon.
 */

/** Rodzaj pozycji biblioteki: skrypt wieloliniowy albo snippet jednego polecenia. */
export type RodzajPozycji = 'skrypt' | 'snippet';

/** Jedna pozycja biblioteki. */
export interface PozycjaBiblioteki {
  /**
   * Identyfikator pozycji nadany przez rdzeń. Pusty znaczy pozycję, która nie
   * ma jeszcze wiersza w dzienniku — tak wchodzi pozycja wczytana z pliku
   * wywozu, zanim zostanie zapisana.
   */
  id?: string;
  nazwa: string;
  rodzaj: RodzajPozycji;
  /** Powłoka, w której treść ma sens — sprawdzana przy uruchomieniu wobec karty. */
  powloka: TerminalShell;
  tresc: string;
  /** Znaczniki porządkujące wykaz; pusty znaczy „bez znaczników”. */
  tagi: string;
  /** Numer wersji rosnący przy każdym zapisie treści. */
  wersja: number;
  /** Czas ostatniego uruchomienia w milisekundach epoki; zero znaczy „nie uruchamiano”. */
  ostatnieUruchomienie: number;
}

/** Jeden zadeklarowany parametr pozycji. */
export interface ParametrSkryptu {
  nazwa: string;
  opis: string;
}

export interface Biblioteka {
  pozycje(): readonly PozycjaBiblioteki[];
  /** Dokłada albo podmienia pozycję o tej nazwie; oddaje prawdę, gdy pozycja była nowa. */
  zapisz(pozycja: PozycjaBiblioteki): boolean;
  usun(nazwa: string): boolean;
  znajdz(nazwa: string): PozycjaBiblioteki | null;
  /** Odnotowuje uruchomienie pozycji — czas ostatniego biegu stoi w wykazie. */
  odnotujUruchomienie(nazwa: string, chwila: number): void;
  /**
   * Zastępuje całą zawartość wykazem z rdzenia.
   *
   * Zastąpienie, nie scalenie: prawdą o bibliotece jest dziennik rdzenia, więc
   * pozycja, której rdzeń nie oddał, przestaje istnieć także na ekranie.
   */
  zastap(nowe: readonly PozycjaBiblioteki[]): void;
}

export function utworzBiblioteke(): Biblioteka {
  const pozycje: PozycjaBiblioteki[] = [];

  function miejsce(nazwa: string): number {
    return pozycje.findIndex((pozycja) => pozycja.nazwa === nazwa);
  }

  return {
    pozycje: () => pozycje,

    zastap(nowe) {
      pozycje.splice(0, pozycje.length, ...nowe);
    },

    zapisz(pozycja) {
      const stare = miejsce(pozycja.nazwa);
      if (stare < 0) {
        pozycje.push(pozycja);
        return true;
      }
      pozycje.splice(stare, 1, pozycja);
      return false;
    },

    usun(nazwa) {
      const stare = miejsce(nazwa);
      if (stare < 0) return false;
      pozycje.splice(stare, 1);
      return true;
    },

    znajdz: (nazwa) => pozycje.find((pozycja) => pozycja.nazwa === nazwa) ?? null,

    odnotujUruchomienie(nazwa, chwila) {
      const pozycja = pozycje[miejsce(nazwa)];
      if (pozycja !== undefined) pozycja.ostatnieUruchomienie = chwila;
    },
  };
}

/** Wzorzec odwołania do parametru w treści skryptu. */
const ODWOLANIE_PARAMETRU = /\{\{([A-Za-z_][A-Za-z0-9_]*)\}\}/g;

/**
 * Czyta deklarację parametrów z nagłówka treści.
 *
 * Nagłówek stoi w komentarzu, bo skrypt ma pozostać uruchamialny także poza
 * platformą — plik z deklaracją w składni obcej powłoce nie wykonałby się
 * nigdzie indziej. Przyjmowane są oba znaki komentarza wykazu powłok: kratka
 * (Bash, PowerShell, Python) i podwójny ukośnik (Node.js).
 */
export function czytajParametry(tresc: string): ParametrSkryptu[] {
  const parametry: ParametrSkryptu[] = [];
  let wNaglowku = false;
  for (const wiersz of tresc.split('\n')) {
    const bezKomentarza = /^\s*(?:#|\/\/)\s?(.*)$/.exec(wiersz)?.[1];
    if (bezKomentarza === undefined) {
      // Pierwszy wiersz spoza komentarza kończy nagłówek: deklaracja stoi
      // wyłącznie na początku treści, żeby nie trzeba było czytać całości.
      if (wNaglowku) break;
      if (wiersz.trim() !== '') break;
      continue;
    }
    if (/^args\s*:\s*$/.test(bezKomentarza.trim())) {
      wNaglowku = true;
      continue;
    }
    if (!wNaglowku) continue;
    const trafienie = /^\s+([A-Za-z_][A-Za-z0-9_]*)\s*:\s*(.*)$/.exec(bezKomentarza);
    if (trafienie === null) break;
    parametry.push({ nazwa: trafienie[1] ?? '', opis: (trafienie[2] ?? '').trim() });
  }
  return parametry;
}

/** Wynik podstawienia parametrów w treści. */
export interface PodstawienieParametrow {
  tresc: string;
  /** Odwołania, dla których nie podano wartości — treść idzie z nimi nietknięta. */
  brakujace: readonly string[];
}

/**
 * Podstawia wartości parametrów w treści.
 *
 * Odwołanie bez wartości nie znika i nie zamienia się w pusty napis: skrypt
 * z pustą ścieżką w miejscu parametru wykonałby się na katalogu głównym zamiast
 * odmówić. Zostaje w treści widoczne, a jego nazwa wraca wołającemu.
 */
export function podstawParametry(
  tresc: string,
  wartosci: ReadonlyMap<string, string>,
): PodstawienieParametrow {
  const brakujace = new Set<string>();
  const podstawiona = tresc.replace(ODWOLANIE_PARAMETRU, (calosc, nazwa: string) => {
    const wartosc = wartosci.get(nazwa);
    if (wartosc === undefined || wartosc === '') {
      brakujace.add(nazwa);
      return calosc;
    }
    return wartosc;
  });
  return { tresc: podstawiona, brakujace: [...brakujace] };
}

/**
 * Kontrola wstępna treści — uwagi wykrywalne bez uruchomienia i bez programu
 * zewnętrznego.
 *
 * To nie jest linter i okno tak jej nie nazywa. Analizę statyczną prowadzą
 * osobne programy (`shellcheck`, `shfmt`, `PSScriptAnalyzer`, `ruff`), których
 * instalka Danaco Console nie niesie; komenda uruchamiająca je nad treścią stoi
 * w kontrakcie i okno ją wywołuje osobną czynnością. Kontrola wstępna działa
 * niezależnie od nich i wychwytuje to, co daje się rozstrzygnąć
 * pewnie: pustkę, znaki końca wiersza rodem z Windows, znacznik kolejności
 * bajtów i rozjazd między deklaracją parametrów a ich użyciem.
 */
export function kontrolaWstepna(tresc: string): string[] {
  const uwagi: string[] = [];
  if (tresc.trim() === '') {
    uwagi.push('Treść jest pusta — nie ma czego uruchomić.');
    return uwagi;
  }
  if (tresc.includes('\r')) {
    uwagi.push(
      'Treść ma znaki powrotu karetki. Powłoki uniksowe czytają je jako część polecenia i kończą przebieg ' +
        'komunikatem o nieznanym poleceniu; usuń zakończenia wierszy rodem z Windows.',
    );
  }
  if (tresc.charCodeAt(0) === 0xfeff) {
    uwagi.push('Treść zaczyna się znacznikiem kolejności bajtów — powłoka policzy go jako pierwszy znak polecenia.');
  }

  const zadeklarowane = new Set(czytajParametry(tresc).map((parametr) => parametr.nazwa));
  const uzyte = new Set<string>();
  for (const trafienie of tresc.matchAll(ODWOLANIE_PARAMETRU)) {
    const nazwa = trafienie[1];
    if (nazwa !== undefined) uzyte.add(nazwa);
  }
  const bezDeklaracji = [...uzyte].filter((nazwa) => !zadeklarowane.has(nazwa));
  if (bezDeklaracji.length > 0) {
    uwagi.push(
      `Treść odwołuje się do parametrów, których nagłówek nie deklaruje: ${bezDeklaracji.join(', ')}. ` +
        'Formularz ich nie pokaże, a odwołanie pojedzie do powłoki dosłownie.',
    );
  }
  const bezUzycia = [...zadeklarowane].filter((nazwa) => !uzyte.has(nazwa));
  if (bezUzycia.length > 0) {
    uwagi.push(
      `Nagłówek deklaruje parametry, których treść nie używa: ${bezUzycia.join(', ')}. ` +
        'Wartość wpisana w formularzu nigdzie nie trafi.',
    );
  }
  return uwagi;
}

/** Rozstrzygnięcie o sprawdzeniu składni: polecenie do wysłania albo powód, dla którego go nie ma. */
export interface SprawdzenieSkladni {
  /** Polecenie sprawdzające; puste znaczy „tej powłoki nie da się sprawdzić bez uruchomienia”. */
  polecenie: string;
  /** Zdanie o tym, co polecenie robi albo dlaczego go nie ma. */
  powod: string;
}

/** Znacznik domykający treść przekazywaną powłoce do sprawdzenia. */
const GRANICA_TRESCI = 'DANACO_KONIEC_SKRYPTU';

/**
 * Polecenie sprawdzające składnię skryptu bez jego wykonania.
 *
 * Bash ma na to własny przełącznik (`bash -n`), więc sprawdzenie nie wymaga
 * ani jednego programu spoza tego, który powłokę już uruchamia. Pozostałe
 * powłoki wykazu kontraktu takiego przełącznika nie mają albo jego użycie
 * wymagałoby przeniesienia treści przez plik, którego moduł nie ma czym
 * zapisać — okno mówi to wprost, zamiast pokazywać przycisk bez skutku.
 */
export function polecenieSprawdzeniaSkladni(
  powloka: TerminalShell,
  tresc: string,
): SprawdzenieSkladni {
  if (powloka !== TerminalShell.Bash && powloka !== TerminalShell.Ssh) {
    return {
      polecenie: '',
      powod:
        `Powłoka ${powloka} nie ma sprawdzenia składni bez uruchomienia treści. Analizę statyczną prowadzi ` +
        'osobna czynność okna („Analiza statyczna treści”, terminal.script.lint) programami leżącymi ' +
        'na maszynie rdzenia — jej odpowiedź mówi wprost, czy program analizy tam jest.',
    };
  }
  if (tresc.includes(GRANICA_TRESCI)) {
    return {
      polecenie: '',
      powod:
        `Treść zawiera znacznik ${GRANICA_TRESCI}, którym okno domyka przekazanie skryptu powłoce — ` +
        'sprawdzenie rozdzieliłoby treść w złym miejscu. Usuń ten napis albo sprawdź skrypt poza platformą.',
    };
  }
  return {
    polecenie: `bash -n <<'${GRANICA_TRESCI}'\n${tresc}\n${GRANICA_TRESCI}`,
    powod:
      'Powłoka czyta treść i sprawdza jej składnię bez wykonania ani jednego polecenia (bash -n). ' +
      'Zerowy kod wyjścia znaczy „składnia poprawna”, nie „skrypt zrobi to, co trzeba”.',
  };
}
