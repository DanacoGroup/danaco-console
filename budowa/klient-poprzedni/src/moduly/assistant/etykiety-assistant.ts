import {
  AssistantActionStatus,
  AssistantActivityKind,
  AssistantOrigin,
  ConfigScope,
  KnowledgeScope,
  MemoryEntryOrigin,
  SlashEntryKind,
} from '../../../../shared/contract';

/**
 * Teksty modułu Assistant widoczne dla Operatora.
 *
 * Jedna odpowiedzialność: nazwy i zdania. Pliki budujące elementy nie noszą
 * napisów, bo napis zmienia się z innego powodu niż układ — i zmieniany
 * w dziesięciu miejscach rozjeżdża się między oknami (wzór:
 * `sterowanie/etykiety-sterowania.ts`).
 *
 * Wartości słowników są kluczowane stałymi kontraktu, więc dopisanie stanu
 * w kontrakcie przerywa kompilację tutaj, zamiast wypuścić na ekran pusty
 * napis.
 */

/**
 * Oprawa znaku [?] modułu: pierścień w rzędzie kontrolek, nie kwadratowy
 * przycisk ikonowy biblioteki (`komponenty/dymek.ts`, pole `KlasyDymka`).
 *
 * Stała stoi w tym pliku, bo sięgają po nią wszystkie widoki modułu, a jest
 * jedną nastawą jednego pojęcia — dokładnie jak słowniki poniżej. Kopia
 * w każdym widoku rozjechałaby oprawę dymków między oknami po pierwszej
 * poprawce arkusza.
 */
export const KLASY_DYMKA = { powloka: 'ma-dymek', znak: 'ma-dymek__znak' } as const;

/** Nazwa stanu zlecenia w kolumnie „Stan" Actions Monitora. */
export const NAZWY_STANOW: Readonly<Record<AssistantActionStatus, string>> = {
  [AssistantActionStatus.Queued]: 'W kolejce',
  [AssistantActionStatus.Running]: 'W toku',
  [AssistantActionStatus.Paused]: 'Wstrzymane',
  [AssistantActionStatus.Done]: 'Wykonane',
  [AssistantActionStatus.Failed]: 'Błąd',
  [AssistantActionStatus.Cancelled]: 'Anulowane',
};

/**
 * Plakietka stanu z biblioteki komponentów — pełna nazwa klasy, nie sklejka.
 *
 * Nazwa składana w czasie działania (`dn-plakietka--${stan}`) byłaby poza
 * zasięgiem kontroli klas CSS bramy; wykaz jawny daje jej wszystkie warianty
 * wprost.
 */
export const PLAKIETKI_STANOW: Readonly<Record<AssistantActionStatus, string>> = {
  [AssistantActionStatus.Queued]: 'dn-plakietka dn-plakietka--informacja',
  [AssistantActionStatus.Running]: 'dn-plakietka dn-plakietka--sygnal',
  [AssistantActionStatus.Paused]: 'dn-plakietka dn-plakietka--ostrzezenie',
  [AssistantActionStatus.Done]: 'dn-plakietka dn-plakietka--sukces',
  [AssistantActionStatus.Failed]: 'dn-plakietka dn-plakietka--blad',
  // Anulowane nie ma własnego wariantu w bibliotece — plakietka bazowa niesie
  // je bez barwy stanu, a znaczenie niesie napis, nie kolor.
  [AssistantActionStatus.Cancelled]: 'dn-plakietka',
};

/** Rodzaj wpisu dziennika w Activity Feed. */
export const NAZWY_RODZAJOW: Readonly<Record<AssistantActivityKind, string>> = {
  [AssistantActivityKind.Command]: 'Polecenie',
  [AssistantActivityKind.Result]: 'Wynik',
  [AssistantActivityKind.Note]: 'Notatka',
};

/** Droga wydania polecenia — głosem albo tekstem; jeden zapis dla obu. */
export const NAZWY_DROG: Readonly<Record<AssistantOrigin, string>> = {
  [AssistantOrigin.Voice]: 'głosem',
  [AssistantOrigin.Text]: 'tekstem',
};

/**
 * Nazwa poziomu zasięgu w wykazie pamięci i w selektorze zapisu.
 *
 * Wykaz jest kompletny wobec `ConfigScope`, więc dopisanie poziomu w kontrakcie
 * przerywa kompilację tutaj, zamiast wypuścić na ekran pusty napis.
 */
export const NAZWY_ZASIEGOW: Readonly<Record<ConfigScope, string>> = {
  [ConfigScope.Application]: 'aplikacja',
  [ConfigScope.Global]: 'globalny',
  [ConfigScope.Environment]: 'środowisko',
  [ConfigScope.Module]: 'moduł',
  [ConfigScope.ModulePair]: 'para modułów',
  [ConfigScope.Project]: 'projekt',
  [ConfigScope.Session]: 'karta sesji',
  [ConfigScope.Role]: 'rola',
  [ConfigScope.Window]: 'okno komunikacji',
};

/** Pochodzenie wpisu pamięci — ręka Operatora albo propozycja modelu. */
export const NAZWY_POCHODZEN: Readonly<Record<MemoryEntryOrigin, string>> = {
  [MemoryEntryOrigin.Operator]: 'Operator',
  [MemoryEntryOrigin.Model]: 'model (propozycja)',
};

/**
 * Poziomy zasięgu pamięci przestawiane w zakładce kontekstów.
 *
 * Wykaz jest krótszy niż `ConfigScope`, bo `memory.toggle` przestawia poziomy
 * WIDOCZNOŚCI pamięci karty sesji, a nie wszystkie poziomy konfiguracji
 * platformy. Cztery poziomy niżej to te, na których pamięć asystenta ma treść:
 * globalny (ustalenia wspólne), projekt, karta sesji i okno komunikacji.
 */
export const POZIOMY_PAMIECI: readonly ConfigScope[] = [
  ConfigScope.Global,
  ConfigScope.Project,
  ConfigScope.Session,
  ConfigScope.Window,
];

/** Grupy filtra statusu Actions Monitora — nośnik wyboru, nie napis. */
export type KodFiltraZlecen = 'wszystkie' | 'wToku' | 'zakonczone' | 'nieudane' | 'anulowane';

/** Nazwa grupy filtra; wchodzi także w zdanie pustki po zawężeniu. */
export const NAZWY_FILTRA: Readonly<Record<KodFiltraZlecen, string>> = {
  wszystkie: 'wszystkie',
  wToku: 'w toku',
  zakonczone: 'zakończone',
  nieudane: 'nieudane',
  anulowane: 'anulowane',
};

/** Pozycje listy wyboru filtra, w kolejności pokazywanej na ekranie. */
export const POZYCJE_FILTRA: ReadonlyArray<readonly [KodFiltraZlecen, string]> = [
  ['wszystkie', 'Filtr statusu: wszystkie'],
  ['wToku', 'Filtr statusu: w toku'],
  ['zakonczone', 'Filtr statusu: zakończone'],
  ['nieudane', 'Filtr statusu: nieudane'],
  ['anulowane', 'Filtr statusu: anulowane'],
];

/** Zakresy wskaźnika wiedzy w zakładkach pamięci semantycznej i bazy wiedzy. */
export const POZYCJE_ZAKRESU_WIEDZY: ReadonlyArray<readonly [KnowledgeScope, string]> = [
  [KnowledgeScope.All, 'Zakres wiedzy: wszystko'],
  [KnowledgeScope.Library, 'Zakres wiedzy: biblioteka'],
  [KnowledgeScope.History, 'Zakres wiedzy: historia rozmów'],
  [KnowledgeScope.Workspace, 'Zakres wiedzy: pliki przestrzeni roboczej'],
];

/** Nazwa rodzaju pozycji katalogu w wierszu wykazu narzędzi. */
export const NAZWY_RODZAJOW_KATALOGU: Readonly<Record<SlashEntryKind, string>> = {
  [SlashEntryKind.Tool]: 'narzędzie albo umiejętność',
  [SlashEntryKind.Action]: 'komenda akcji',
};

/** Rodzaj pozycji katalogu po ukośniku w filtrze Command & Tools Hub. */
export type KodFiltraKatalogu = 'wszystkie' | SlashEntryKind;

/** Pozycje listy wyboru rodzaju pozycji katalogu narzędzi. */
export const POZYCJE_FILTRA_KATALOGU: ReadonlyArray<readonly [KodFiltraKatalogu, string]> = [
  ['wszystkie', 'Rodzaj pozycji: wszystkie'],
  [SlashEntryKind.Tool, 'Rodzaj pozycji: narzędzia i umiejętności'],
  [SlashEntryKind.Action, 'Rodzaj pozycji: komendy akcji'],
];

/** Rodzaje wpisu dziennika w filtrze Activity Feed; `wszystkie` bez zawężenia. */
export type KodFiltraDziennika = 'wszystkie' | AssistantActivityKind;

/** Pozycje listy wyboru filtra rodzaju wpisu dziennika. */
export const POZYCJE_FILTRA_DZIENNIKA: ReadonlyArray<readonly [KodFiltraDziennika, string]> = [
  ['wszystkie', 'Rodzaj wpisu: wszystkie'],
  [AssistantActivityKind.Command, 'Rodzaj wpisu: polecenia'],
  [AssistantActivityKind.Result, 'Rodzaj wpisu: wyniki'],
  [AssistantActivityKind.Note, 'Rodzaj wpisu: notatki'],
];

/**
 * Zdania stanu pustego.
 *
 * Stan pusty tłumaczy, czym okno jest i czym się je zapełnia. Zdanie mówiące
 * wyłącznie, czego nie ma („Brak danych", „Brak działań w historii"), zostawia
 * Operatora przed oknem, o którym nie wie ani po co ono stoi, ani co miałby
 * zrobić, żeby coś w nim zobaczyć. Każde zdanie poniżej ma więc dwie części:
 * czym okno jest i co je zapełni.
 *
 * Każde okno ma swoje zdanie spoczynku, bo jedno wspólne nie nazywa okna,
 * przed którym Operator stoi. Rozróżnienie „przed pytaniem" i „po odpowiedzi"
 * zostaje: to dwa różne fakty i mają dwa różne zdania.
 */
export const PUSTE = {
  konsolaSpoczynek:
    'Voice Console jest punktem wejścia modułu Assistant: stąd wydajesz asystentowi ' +
    'polecenia — głosem albo tekstem, jednym zapisem. Okno nie pytało jeszcze rdzenia ' +
    'o przydział okna modułu; odczyt rusza z wejściem do modułu.',
  konsola:
    'Voice Console jest punktem wejścia modułu Assistant: stąd wydajesz asystentowi ' +
    'polecenia — głosem albo tekstem, jednym zapisem. Wpisz treść w polu transkrypcji ' +
    'powyżej i naciśnij „Wyślij polecenie": rdzeń założy zlecenie, którego przebieg ' +
    'pokaże Actions Monitor, a ślad — Activity Feed.',
  monitorSpoczynek:
    'Actions Monitor prowadzi zlecenia wieloetapowe asystenta: pokazuje ich stan i etap ' +
    'oraz pozwala je wstrzymać, wznowić, anulować i przestawić im priorytet. Okno nie ' +
    'pytało jeszcze rdzenia — odczyt rusza z wejściem do modułu.',
  zlecenia:
    'Actions Monitor prowadzi zlecenia wieloetapowe asystenta: pokazuje ich stan i etap ' +
    'oraz pozwala nimi sterować. Rdzeń nie prowadzi dziś ani jednego — wydaj polecenie ' +
    'w Voice Console, a zlecenie z niego stanie w tej tabeli.',
  dziennikSpoczynek:
    'Activity Feed prowadzi chronologiczny zapis działań asystenta — polecenia, wyniki ' +
    'i notatki jednym ciągiem, wspólnym dla drogi głosowej i tekstowej. Okno nie pytało ' +
    'jeszcze rdzenia — odczyt rusza z wejściem do modułu.',
  dziennik:
    'Activity Feed prowadzi chronologiczny zapis działań asystenta — polecenia, wyniki ' +
    'i notatki jednym ciągiem, wspólnym dla drogi głosowej i tekstowej. Zapis jest ' +
    'jeszcze pusty; pierwsze polecenie wydane w Voice Console założy pierwszy wpis.',
  dziennikZlecenia:
    'Zapis jest zawężony do jednego zlecenia, a rdzeń nie ma dla niego ani jednego ' +
    'wpisu. Naciśnij „Cały dziennik", aby wrócić do całości zapisu modułu.',
  dziennikSzukanie:
    'Rdzeń oddał wpisy dziennika, ale żaden nie odpowiada bieżącemu zawężeniu. ' +
    'Wyczyść pole wyszukiwania albo przestaw rodzaj wpisu, aby zobaczyć zapis w całości.',
  pamiec:
    'Memory & Context Manager pokazuje pamięć asystenta w całości: fakty, wyszukiwanie ' +
    'po znaczeniu, poziomy pamięci karty sesji i wskaźnik wiedzy. Rdzeń nie ma dziś ani ' +
    'jednego wpisu w tym zasięgu — pierwszy zapiszesz polem powyżej.',
  pamiecSpoczynek:
    'Memory & Context Manager pokazuje pamięć asystenta w całości: fakty, wyszukiwanie ' +
    'po znaczeniu, poziomy pamięci karty sesji i wskaźnik wiedzy. Okno nie pytało jeszcze ' +
    'rdzenia — odczyt rusza z wejściem do modułu.',
  wiedzaSpoczynek:
    'Wyszukiwanie po znaczeniu sięga wskaźnika wiedzy Operatora (knowledge.search) — ' +
    'biblioteki, historii rozmów i plików przestrzeni roboczej. Wpisz pytanie i naciśnij ' +
    '„Szukaj po znaczeniu"; wskaźnik zbudujesz przyciskiem w zakładce bazy wiedzy.',
  wiedza:
    'Wskaźnik nie ma fragmentu odpowiadającego temu pytaniu. Zbuduj wskaźnik w zakładce ' +
    'bazy wiedzy albo zapytaj inaczej — szukanie idzie po znaczeniu, nie po słowach.',
  narzedziaSpoczynek:
    'Command & Tools Hub prowadzi katalog pozycji po ukośniku (tools.catalog.list): ' +
    'narzędzia i umiejętności do dołożenia oraz komendy akcji. Okno nie pytało jeszcze ' +
    'rdzenia — odczyt rusza z wejściem do modułu.',
  narzedzia:
    'Rdzeń nie zna pozycji katalogu spełniającej to zawężenie. Wyczyść pole szukania albo ' +
    'przestaw rodzaj i grupę — katalog liczy setki pozycji i domyślnie jest przycięty.',
  rutynySpoczynek:
    'Rutyny to automatyki modułu Automations wraz z ich harmonogramem ' +
    '(automation.workflow.list, schedule.get). Okno nie pytało jeszcze rdzenia — ' +
    'odczyt rusza przyciskiem „Odczytaj rutyny".',
  rutyny:
    'Rdzeń nie prowadzi ani jednej automatyki. Makro złożone w zakładce obok wysłane ' +
    'przyciskiem „→ Wyślij do Automations" założy pierwszą i stanie w tym wykazie.',
  akcjeSpoczynek:
    'Szybkie akcje pokazują katalog akcji zasięgu modułu pochodzący z rdzenia ' +
    '(action.list); naciśnięcie kafla wstawia treść do pola polecenia. Okno ' +
    'nie pytało jeszcze rdzenia — odczyt rusza z wejściem do modułu.',
  akcje:
    'Szybkie akcje pokazują katalog akcji zasięgu modułu pochodzący z rdzenia ' +
    '(action.list). Rdzeń nie zna dziś ani jednej akcji tego zasięgu — katalog ' +
    'zapełnia się wierszami rejestru akcji rdzenia, nie z tego okna.',
} as const;

/**
 * Zdania wykazu zleceń spoza okna modułu (Actions Monitor, druga tabela).
 *
 * Nie są to zdania stanu pustego i dlatego nie stoją w `PUSTE`: sekcja jest
 * ukryta, dopóki nie ma czego pokazać, więc pustki nie ma czym nazywać.
 *
 * Zdanie `zasieg` mówi wprost, czego w wykazie nie będzie — zleceń zamkniętych
 * przed otwarciem modułu. Wykaz podany bez tego zastrzeżenia wyglądałby na
 * komplet pracy asystenta w sesji, a nim nie jest: kontrakt nie ma odczytu
 * zleceń po sesji.
 */
export const SPOZA_OKNA = {
  tytul: 'Zlecenia asystenta spoza tego okna',
  zasieg:
    'Ta sama sesja, inne okno — tak siada zlecenie wydane głosem w nakładce AOD, ' +
    'bo nakładka dobiera okno sama. Wykaz zbiera się ze zdarzeń rdzenia od chwili ' +
    'wejścia do modułu: kontrakt nie ma odczytu zleceń po sesji, więc zlecenia ' +
    'zamknięte wcześniej nie mają jak się tu znaleźć. Sterowanie i podgląd działają ' +
    'tak samo jak wyżej — rdzeń przyjmuje je po identyfikatorze zlecenia, nie po oknie.',
} as const;

/** Zdania stanu ładowania — nazywają wywołanie, nie samą czynność. */
export const ODCZYTY = {
  zlecenia: 'Czytam stan zleceń asystenta…',
  dziennik: 'Czytam dziennik działań…',
  akcje: 'Czytam katalog akcji modułu…',
  polecenie: 'Rdzeń prowadzi polecenie…',
  przerwanie: 'Zamawiam przerwanie zlecenia…',
  pamiec: 'Czytam pamięć asystenta…',
  wiedza: 'Szukam we wskaźniku wiedzy po znaczeniu…',
  wskaznik: 'Buduję wskaźnik znaczenia…',
  poziomy: 'Przestawiam poziomy pamięci karty sesji…',
  narzedzia: 'Czytam katalog narzędzi i komend akcji…',
  dolozenia: 'Czytam dołożenia narzędzi karty sesji…',
  dolozenie: 'Zamawiam zmianę dołożenia narzędzia…',
  rutyny: 'Czytam automatyki i ich harmonogramy…',
  wysylkaMakra: 'Zapisuję makro jako automatykę…',
  harmonogram: 'Zapisuję harmonogram automatyki…',
  silnikMowy: 'Sprawdzam dostępność silnika mowy…',
  transkrypcja: 'Silnik rozpoznaje nagranie…',
} as const;

/** Chwila zapisana w dzienniku, w strefie Operatora. */
export function chwila(znacznik: number): string {
  return new Date(znacznik).toLocaleString('pl-PL');
}

/**
 * Dzień wpisu — nagłówek grupy chronologicznej Activity Feed.
 *
 * Napis jest jednocześnie kluczem grupowania, więc dwa wpisy z tej samej doby
 * dają dokładnie ten sam łańcuch. Data liczy się w strefie Operatora, bo to
 * jego doba rozdziela „dziś" od „wczoraj", a nie doba serwera.
 */
export function dzien(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}
