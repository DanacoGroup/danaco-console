import { liczba, listaTekstow, obiekt, tekst } from './odczyt-fragmentu';

/**
 * Prowenancja wywołania — odpowiedź na pytanie „co poszło do modelu".
 *
 * Rdzeń nadaje ją pierwszym fragmentem strumienia, przed jakąkolwiek treścią
 * odpowiedzi, także wtedy, gdy proces kanału w ogóle nie wystartuje
 * (server/internal/injection/przebieg.go). Kształt odpowiada strukturze
 * `injection.Prowenancja`; nazwy pól po stronie sieci są angielskie, po stronie
 * modułu polskie.
 *
 * Prowenancja jest wyłącznie opisem wywołania — niczego nie dopuszcza ani nie
 * wstrzymuje.
 */
export interface Prowenancja {
  /** Plik wykonywalny kanału. */
  program: string;
  /** Pełny wiersz wywołania procesu. */
  argv: string[];
  /** Prompt systemowy złożony z warstw nakładki. */
  promptSystemowy: string;
  /** Tryb nakładki: dopisanie albo zastąpienie promptu. */
  trybNakladki: string;
  /** Nazwy warstw nakładki w kolejności złożenia. */
  warstwyNakladki: string[];
  /** Skrót nakładki — wyłącznie diagnostyka. */
  skrotNakladki: string;
  /** Model główny przekazany kanałowi. */
  model: string;
  /** Model zapasowy. */
  modelZapasowy: string;
  /** Nakład rozumowania. */
  naklad: string;
  /** Tryb uprawnień okna. */
  trybUprawnien: string;
  /** Katalogi udostępnione modelowi. */
  katalogi: string[];
  /** Katalog roboczy procesu. */
  katalogRoboczy: string;
  /** Plik ustawień przekazany kanałowi. */
  plikUstawien: string;
  /** Wskazane konfiguracje MCP. */
  konfiguracjaMCP: string[];
  /** Identyfikator wznawianej rozmowy kanału. */
  wznowienie: string;
  /** Konto — kod, nigdy poświadczenie. */
  konto: string;
  /** Skąd wzięła się tożsamość wywołania: wskazanie, kolejność puli, rotacja, otoczenie. */
  zrodloKonta: string;
  /** Katalog konfiguracji konta. */
  katalogKonta: string;
  /** Numer próby; druga i dalsze to skutek rotacji konta. */
  proba: number;
  /** Powód próby — wypełniony przy rotacji konta. */
  powodProby: string;
  /** Chwila złożenia prowenancji, zapis ISO 8601. */
  chwila: string;

  // ── Pola dodatkowe kanałów sieciowych/echo ─────────────────────────────
  // Serwer emituje je tylko dla kanałów bez procesu (api/echo); dla kanału
  // głównego CLI są nieobecne i zostają puste. Pokazujemy je wyłącznie, gdy
  // niosą wartość (models/prowenancja.go, znaczniki `,omitempty`).
  /** Kod kanału z rejestru, który wykonuje wywołanie. */
  kanal: string;
  /** Klucz adaptera obsługującego kanał (echo, cli, api, …). */
  adapter: string;
  /** Punkt końcowy wywołania sieciowego; puste dla kanałów bez sieci. */
  adres: string;
  /** Środowisko wykonania — gdzie model pracuje. */
  srodowiskoWykonania: string;
}

/**
 * Odczytuje prowenancję z ładunku fragmentu.
 * Ładunek nieczytelny daje `null` — fragment zostaje pokazany jako sam fakt
 * wywołania, bez szczegółów, zamiast przerywać strumień.
 */
export function odczytajProwenancje(dane: unknown): Prowenancja | null {
  const zrodlo = obiekt(dane);
  if (zrodlo === null) return null;
  return {
    program: tekst(zrodlo, 'program'),
    argv: listaTekstow(zrodlo, 'argv'),
    promptSystemowy: tekst(zrodlo, 'systemPrompt'),
    trybNakladki: tekst(zrodlo, 'overlayMode'),
    warstwyNakladki: listaTekstow(zrodlo, 'overlayLayers'),
    skrotNakladki: tekst(zrodlo, 'overlayHash'),
    model: tekst(zrodlo, 'model'),
    modelZapasowy: tekst(zrodlo, 'fallbackModel'),
    naklad: tekst(zrodlo, 'effort'),
    trybUprawnien: tekst(zrodlo, 'permissionMode'),
    katalogi: listaTekstow(zrodlo, 'directories'),
    katalogRoboczy: tekst(zrodlo, 'workingDirectory'),
    plikUstawien: tekst(zrodlo, 'settings'),
    konfiguracjaMCP: listaTekstow(zrodlo, 'mcpConfig'),
    wznowienie: tekst(zrodlo, 'resume'),
    konto: tekst(zrodlo, 'account'),
    zrodloKonta: tekst(zrodlo, 'accountOrigin'),
    katalogKonta: tekst(zrodlo, 'accountConfigDir'),
    proba: liczba(zrodlo, 'attempt'),
    powodProby: tekst(zrodlo, 'attemptReason'),
    chwila: tekst(zrodlo, 'at'),
    kanal: tekst(zrodlo, 'channel'),
    adapter: tekst(zrodlo, 'adapter'),
    adres: tekst(zrodlo, 'endpoint'),
    srodowiskoWykonania: tekst(zrodlo, 'executionEnv'),
  };
}

/** Wiersz wywołania w postaci jednej linii — do nagłówka bloku prowenancji. */
export function wierszWywolania(p: Prowenancja): string {
  return [p.program, ...p.argv].filter((czesc) => czesc.length > 0).join(' ');
}
