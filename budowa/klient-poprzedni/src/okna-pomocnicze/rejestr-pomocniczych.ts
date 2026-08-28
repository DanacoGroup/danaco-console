/**
 * Rejestr okien pomocniczych modułów pasa inżynierskiego — Developer, Diagnostics, Apps i Roundtable — nazywa też pozycje niezbudowane wraz z powodem, bo stany pozycji są rozłączne i mieszanie ich zaciera prawdę o oknie.
 */
export type StanPomocniczego =
  /** Okno stoi w pasie i woła rdzeń. */
  | 'zbudowane'
  /** Okno modułu już to niesie — pas go nie powiela. */
  | 'stoi-w-module'
  /** Buduje je inna praca; pas zostawia miejsce i nie zajmuje go. */
  | 'w-innej-pracy'
  /** Komendy nie ma w kontrakcie — zmierzone, nie założone. */
  | 'bez-komendy-rdzenia'
  /** Komenda jest, ale nikt nie rozstrzygnął, czy okno ma tu stanąć. */
  | 'do-rozstrzygniecia';

/** Jedna pozycja spisu okien pomocniczych modułu: kod, nazwa widoczna, przeznaczenie, stan i wyjaśnienie. */
export interface OpisPomocniczego {
  /** Kod pozycji — trafia w `data-okno` wiersza spisu. */
  kod: string;
  /** Nazwa widziana przez Operatora. */
  nazwa: string;
  /** Po co Operatorowi to okno. */
  przeznaczenie: string;
  stan: StanPomocniczego;
  /** Zdanie mówiące wprost, co z tym oknem jest, czytane przez Operatora w pasie inżynierskim. */
  wyjasnienie: string;
}

/** Podgląd w tle bash — pozycja wspólna wszystkim modułom pasa inżynierskiego korzystającym ze wspólnej powłoki. */
const PODGLAD_BASH: OpisPomocniczego = {
  kod: 'podglad-bash',
  nazwa: 'Podgląd w tle (bash)',
  przeznaczenie: 'Zbiorcze wyjście wszystkich otwartych kart powłoki, bez otwierania karty.',
  stan: 'zbudowane',
  wyjasnienie:
    'Stoi na komendzie `terminal.output.stream`, którą rdzeń rejestruje ' +
    '(`handlers_terminal_wyjscie.go`); nowe wiersze dochodzą zdarzeniem `stream.chunk`.',
};

/**
 * Terminal — okno pomocnicze Developera, Diagnostics i Apps, którego kartę powłoki stawia wyłącznie wytwórnia paneli.
 */
const TERMINAL: OpisPomocniczego = {
  kod: 'terminal',
  nazwa: 'Terminal',
  przeznaczenie: 'Karta powłoki wewnątrz modułu — polecenia, procesy, wyjście na żywo.',
  stan: 'zbudowane',
  wyjasnienie:
    'Terminal nie jest samodzielnym modułem (rozstrzygnięcie Właściciela) i stoi jako okno ' +
    'pomocnicze Developera, Diagnostics i Apps. Karta bierze tryb uprawnień i katalog roboczy ' +
    'z okna GOSPODARZA: rdzeń wiąże kartę powłoki z oknem, nie z modułem ' +
    '(`adapter_modul_terminal_okno.go` żąda wyłącznie `windowId` okna otwartego).',
};

/**
 * Historia rozmowy jest pozycją wspólną wszystkim spisom, bo rdzeń wiąże historię z oknem, nie z modułem, więc żaden moduł pasa nie jest dla niej bezprzedmiotowy.
 */
const HISTORIA_ROZMOWY: OpisPomocniczego = {
  kod: 'historia-rozmowy',
  nazwa: 'Historia rozmowy',
  przeznaczenie:
    'Trwały zapis wypowiedzi okna od najnowszej wstecz — usuwanie pozycji i zasada przechowywania.',
  stan: 'zbudowane',
  wyjasnienie:
    'Stoi na trzech komendach z uchwytem w rdzeniu — `history.load`, `history.delete`, ' +
    '`retention.set` (`handlers_historia.go`) — oraz na zdarzeniu `history.changed`, którym ' +
    'rdzeń zgłasza kasowanie cudzą ręką. Zasadę przechowywania rdzeń EGZEKWUJE: przy starcie ' +
    '(`trwalosc_kosza.go`), przy każdym odczycie i zaraz po jej zapisaniu.',
};

/** Artefakty — pozycja bez pokrycia w kontrakcie, wpisana do spisu z zamysłu jako zapowiedź przyszłej zdolności. */
const ARTEFAKTY: OpisPomocniczego = {
  kod: 'artefakty',
  nazwa: 'Artefakty',
  przeznaczenie: 'Wytwory przebiegu: paczki, raporty, zrzuty, pliki wyjściowe budowania.',
  stan: 'bez-komendy-rdzenia',
  wyjasnienie:
    'Kontrakt nie ma ani jednej komendy artefaktów (zmierzone przeglądem wykazu `Command` ' +
    'w `shared/contract.ts`). `DeveloperBuild.logRef` odsyła do logu, lecz komendy rozwijającej ' +
    'ten odnośnik też nie ma — okno artefaktów nie miałoby czego zapytać.',
};

/** Pliki środowiska — pozycja bez pokrycia w kontrakcie, czekająca na komendę rdzenia obsługującą ten zasób. */
const PLIKI_SRODOWISKA: OpisPomocniczego = {
  kod: 'pliki-srodowiska',
  nazwa: 'Pliki środowiska',
  przeznaczenie: 'Zmienne i pliki konfiguracji uruchomienia (`.env` i pokrewne).',
  stan: 'bez-komendy-rdzenia',
  wyjasnienie:
    'Kontrakt ma `environment.list` i `environment.enter`, ale ŚRODOWISKO znaczy tam co innego: ' +
    'kartę wejścia platformy z wykazem modułów (`Environment.navigationKind`, `moduleCodes`), ' +
    'a nie pliki konfiguracji uruchomienia. Podłożenie jednego pod drugie byłoby zgadywaniem.',
};

/** Przeglądarka — komendy potrzebne w kontrakcie już są, decyzji o budowie samego panelu jeszcze nie podjęto. */
const PRZEGLADARKA: OpisPomocniczego = {
  kod: 'przegladarka',
  nazwa: 'Przeglądarka',
  przeznaczenie: 'Podgląd uruchomionej aplikacji i stron obok kodu.',
  stan: 'do-rozstrzygniecia',
  wyjasnienie:
    'Komendy `browser.navigate` i `browser.snapshot.get` w kontrakcie SĄ, ale należą do modułu ' +
    'Browser. Czy okno Browsera ma stanąć wewnątrz Developera i Diagnostics jako pomocnicze — ' +
    'tego nikt nie rozstrzygnął, a wstawienie go tu byłoby decyzją, nie robotą.',
};

/**
 * Zasoby Designu to panel wędrujący do modułu, który z nich korzysta — dostępny w Studio Editorze i module Apps, czytający te same dane co panel operacyjny Designu, nie drugie źródło prawdy.
 */
const ZASOBY_DESIGNU: OpisPomocniczego = {
  kod: 'zasoby-designu',
  nazwa: 'Zasoby (Design)',
  przeznaczenie: 'Wykaz zasobów graficznych Designu obok rozmowy, z zawężaniem po nazwie.',
  stan: 'zbudowane',
  wyjasnienie:
    'Stoi na komendzie `design.asset.list` i zdarzeniu `design.asset.changed` — obie z uchwytem ' +
    'w rdzeniu. Panel czyta zbiór PEŁNY: `windowId` filtru zasobów jest opcjonalne, a okno, ' +
    'które panel dostaje, jest oknem GOSPODARZA (Apps), nie oknem Designu — podstawienie go ' +
    'zawęziłoby wykaz do zasobów cudzego okna, czyli najczęściej do pustki.',
};

/** Przebieg debaty — panel modułu Roundtable pokazywany obok rozmowy, niosący zapis wypowiedzi wielu modeli. */
const PRZEBIEG_DEBATY: OpisPomocniczego = {
  kod: 'przebieg-debaty',
  nazwa: 'Przebieg debaty',
  przeznaczenie:
    'Tura bieżąca, skład i wypowiedzi wielu modeli obok rozmowy — wypowiedź rośnie na żywo.',
  stan: 'zbudowane',
  wyjasnienie:
    'Stoi na zdarzeniach `roundtable.debate.changed` (tura i wypowiedzi w całości) oraz ' +
    '`stream.chunk` (fragmenty wypowiedzi na żywo, nadawane przez ' +
    '`adapter_modul_roundtable_glos.go`). Komendy odczytu przebiegu obszar `roundtable.*` nie ma, ' +
    'więc panel pokazuje wyłącznie to, co usłyszał od otwarcia.',
};

/** Spisy okien pomocniczych dla poszczególnych modułów pasa inżynierskiego, indeksowane kodem każdego modułu. */
const SPISY: ReadonlyMap<string, readonly OpisPomocniczego[]> = new Map([
  [
    'developer',
    [
      PODGLAD_BASH,
      TERMINAL,
      HISTORIA_ROZMOWY,
      {
        kod: 'pliki',
        nazwa: 'Pliki',
        przeznaczenie: 'Drzewo plików repozytorium wraz z otwarciem pliku w edytorze.',
        stan: 'stoi-w-module',
        wyjasnienie:
          'Niesie je okno Project Tree złożenia modułu (`developer.tree.get`, `developer.file.open`). ' +
          'Pas okien pomocniczych go nie powiela — drugie drzewo tych samych plików byłoby drugą ' +
          'drogą do tej samej treści.',
      },
      PRZEGLADARKA,
      ARTEFAKTY,
      PLIKI_SRODOWISKA,
    ],
  ],
  [
    'diagnostics',
    [
      PODGLAD_BASH,
      TERMINAL,
      HISTORIA_ROZMOWY,
      {
        kod: 'pliki',
        nazwa: 'Pliki',
        przeznaczenie: 'Pliki obszaru roboczego oglądane obok logów i błędów.',
        stan: 'bez-komendy-rdzenia',
        wyjasnienie:
          'Obszar `diagnostics.*` nie ma ani jednej komendy plikowej, a `developer.tree.get` wymaga ' +
          'kontraktem `windowId` OKNA DEVELOPERA — podanie mu okna Diagnostics kierowałoby odczyt ' +
          'do cudzego katalogu roboczego.',
      },
      PRZEGLADARKA,
      ARTEFAKTY,
      PLIKI_SRODOWISKA,
    ],
  ],
  [
    // Apps — trzeci moduł gospodarza Terminala, obok Developera i Diagnostics.
    'apps',
    [
      PODGLAD_BASH,
      TERMINAL,
      HISTORIA_ROZMOWY,
      ZASOBY_DESIGNU,
      {
        kod: 'pliki',
        nazwa: 'Pliki',
        przeznaczenie: 'Pliki obszaru roboczego oglądane obok warsztatu.',
        stan: 'bez-komendy-rdzenia',
        wyjasnienie:
          'Obszar `apps.*` nie ma ani jednej komendy plikowej, a `developer.tree.get` wymaga ' +
          'kontraktem `windowId` OKNA DEVELOPERA — podanie mu okna Apps kierowałoby odczyt do ' +
          'cudzego katalogu roboczego. Ta sama granica co w Diagnostics.',
      },
      PRZEGLADARKA,
      ARTEFAKTY,
      PLIKI_SRODOWISKA,
    ],
  ],
  [
    // Roundtable ma spis krótki z zamysłu: moduł prowadzi debatę modeli w jednym oknie, nie komplet okien.
    'roundtable',
    [PRZEBIEG_DEBATY, HISTORIA_ROZMOWY],
  ],
]);

/**
 * Spis okien pomocniczych modułu; moduł spoza spisu dostaje wykaz pusty, a pas mówi wprost, że okien pomocniczych dla niego nie spisano.
 */
export function oknaPomocnicze(kodModulu: string): readonly OpisPomocniczego[] {
  return SPISY.get(kodModulu) ?? [];
}
