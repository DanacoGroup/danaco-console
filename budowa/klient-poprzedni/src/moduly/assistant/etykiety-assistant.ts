import {
  AssistantActionStatus,
  AssistantActivityKind,
  AssistantOrigin,
  ConfigScope,
  KnowledgeScope,
  MemoryEntryOrigin,
  SlashEntryKind,
} from '../../../../shared/contract';

/** Teksty modułu Assistant widoczne dla Operatora — jedna odpowiedzialność: nazwy i zdania. */

/**
 * Oprawa znaku pomocy modułu: pierścień w rzędzie kontrolek, nie kwadratowy przycisk ikonowy
 * biblioteki. Stała stoi w tym pliku, bo sięgają po nią wszystkie widoki modułu.
 */
export const KLASY_DYMKA = { powloka: 'ma-dymek', znak: 'ma-dymek__znak' } as const;

/** Nazwa stanu zlecenia w kolumnie „Stan” Actions Monitora — jeden napis dla każdej wartości wyliczenia stanu. */
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
  // Anulowane nie ma własnego wariantu w bibliotece — plakietka bazowa niesie je bez barwy stanu.
  [AssistantActionStatus.Cancelled]: 'dn-plakietka',
};

/** Rodzaj wpisu dziennika w Activity Feed — nazwa czytelna dla Operatora, jedna na każdą wartość wyliczenia rodzaju. */
export const NAZWY_RODZAJOW: Readonly<Record<AssistantActivityKind, string>> = {
  [AssistantActivityKind.Command]: 'Polecenie',
  [AssistantActivityKind.Result]: 'Wynik',
  [AssistantActivityKind.Note]: 'Notatka',
};

/** Droga wydania polecenia — głosem albo tekstem; jeden zapis dla obu, wspólny dziennikowi i konsoli poleceń. */
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

/** Pochodzenie wpisu pamięci — ręka Operatora albo propozycja modelu, nazwane osobno w wykazie pamięci. */
export const NAZWY_POCHODZEN: Readonly<Record<MemoryEntryOrigin, string>> = {
  [MemoryEntryOrigin.Operator]: 'Operator',
  [MemoryEntryOrigin.Model]: 'model (propozycja)',
};

/**
 * Poziomy zasięgu pamięci przestawiane w zakładce kontekstów — węższy wykaz niż pełny zasięg
 * konfiguracji, bo obejmuje tylko poziomy widoczności pamięci karty sesji z treścią.
 */
export const POZIOMY_PAMIECI: readonly ConfigScope[] = [
  ConfigScope.Global,
  ConfigScope.Project,
  ConfigScope.Session,
  ConfigScope.Window,
];

/** Grupy filtra statusu Actions Monitora — nośnik wyboru, nie napis; nazwę każdej grupy niesie osobny słownik. */
export type KodFiltraZlecen = 'wszystkie' | 'wToku' | 'zakonczone' | 'nieudane' | 'anulowane';

/** Nazwa grupy filtra widoczna w kontrolce wyboru statusu; wchodzi także w zdanie pustki po zawężeniu wykazu. */
export const NAZWY_FILTRA: Readonly<Record<KodFiltraZlecen, string>> = {
  wszystkie: 'wszystkie',
  wToku: 'w toku',
  zakonczone: 'zakończone',
  nieudane: 'nieudane',
  anulowane: 'anulowane',
};

/** Pozycje listy wyboru filtra statusu zleceń, ułożone w kolejności pokazywanej na ekranie kontrolki wyboru. */
export const POZYCJE_FILTRA: ReadonlyArray<readonly [KodFiltraZlecen, string]> = [
  ['wszystkie', 'Filtr statusu: wszystkie'],
  ['wToku', 'Filtr statusu: w toku'],
  ['zakonczone', 'Filtr statusu: zakończone'],
  ['nieudane', 'Filtr statusu: nieudane'],
  ['anulowane', 'Filtr statusu: anulowane'],
];

/** Zakresy wskaźnika wiedzy w zakładkach pamięci semantycznej i bazy wiedzy — pozycje gotowe do kontrolki wyboru. */
export const POZYCJE_ZAKRESU_WIEDZY: ReadonlyArray<readonly [KnowledgeScope, string]> = [
  [KnowledgeScope.All, 'Zakres wiedzy: wszystko'],
  [KnowledgeScope.Library, 'Zakres wiedzy: biblioteka'],
  [KnowledgeScope.History, 'Zakres wiedzy: historia rozmów'],
  [KnowledgeScope.Workspace, 'Zakres wiedzy: pliki przestrzeni roboczej'],
];

/** Nazwa rodzaju pozycji katalogu w wierszu wykazu narzędzi — narzędzie albo umiejętność, albo komenda akcji. */
export const NAZWY_RODZAJOW_KATALOGU: Readonly<Record<SlashEntryKind, string>> = {
  [SlashEntryKind.Tool]: 'narzędzie albo umiejętność',
  [SlashEntryKind.Action]: 'komenda akcji',
};

/** Rodzaj pozycji katalogu po ukośniku w filtrze Command & Tools Hub; wartość „wszystkie” znaczy brak zawężenia. */
export type KodFiltraKatalogu = 'wszystkie' | SlashEntryKind;

/** Pozycje listy wyboru rodzaju pozycji katalogu narzędzi, gotowe do osadzenia wprost w kontrolce filtra. */
export const POZYCJE_FILTRA_KATALOGU: ReadonlyArray<readonly [KodFiltraKatalogu, string]> = [
  ['wszystkie', 'Rodzaj pozycji: wszystkie'],
  [SlashEntryKind.Tool, 'Rodzaj pozycji: narzędzia i umiejętności'],
  [SlashEntryKind.Action, 'Rodzaj pozycji: komendy akcji'],
];

/** Rodzaje wpisu dziennika w filtrze Activity Feed; wartość „wszystkie” znaczy brak zawężenia całego wykazu. */
export type KodFiltraDziennika = 'wszystkie' | AssistantActivityKind;

/** Pozycje listy wyboru filtra rodzaju wpisu dziennika, gotowe do osadzenia wprost w kontrolce wyboru filtra. */
export const POZYCJE_FILTRA_DZIENNIKA: ReadonlyArray<readonly [KodFiltraDziennika, string]> = [
  ['wszystkie', 'Rodzaj wpisu: wszystkie'],
  [AssistantActivityKind.Command, 'Rodzaj wpisu: polecenia'],
  [AssistantActivityKind.Result, 'Rodzaj wpisu: wyniki'],
  [AssistantActivityKind.Note, 'Rodzaj wpisu: notatki'],
];

/**
 * Zdania stanu pustego, z których każde tłumaczy czym okno jest i czym się je zapełnia,
 * a nie tylko czego w nim nie ma — każde okno ma swoje zdanie spoczynku osobne od zdania
 * po odpowiedzi rdzenia.
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
 * Zdania wykazu zleceń spoza okna modułu w drugiej tabeli Actions Monitora — sekcja ukryta,
 * dopóki nie ma czego pokazać, więc zdania stanu pustego jej nie dotyczą.
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

/** Zdania stanu ładowania — nazywają wywołanie skierowane do rdzenia, nie samą czynność wykonaną przez Operatora. */
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

/** Chwila zapisana w dzienniku, sformatowana w strefie czasowej Operatora, a nie w strefie czasowej serwera. */
export function chwila(znacznik: number): string {
  return new Date(znacznik).toLocaleString('pl-PL');
}

/**
 * Dzień wpisu — nagłówek grupy chronologicznej Activity Feed, jednocześnie klucz grupowania
 * liczony w strefie Operatora, nie w strefie serwera.
 */
export function dzien(znacznik: number): string {
  return new Date(znacznik).toLocaleDateString('pl-PL', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  });
}
