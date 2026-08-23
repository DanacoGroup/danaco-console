/**
 * Teksty modułu Browser widoczne dla Operatora — kody katalogu, objaśnienia
 * dymków [?] i zdania nazywające braki kontraktu.
 *
 * Jedna odpowiedzialność: słowo mówione do Operatora. Pliki budujące elementy
 * nie trzymają ani jednego zdania, bo wtedy zmiana brzmienia wymagałaby wejścia
 * w widok, a te same zdania powtarzałyby się w trzech oknach naraz (wzór:
 * `sterowanie/etykiety-sterowania.ts`).
 *
 * Zdania o brakach stoją tu razem z resztą: brak drogi w kontrakcie jest treścią
 * widoku tak samo jak nazwa przycisku — okno ma go wypowiedzieć, nie przemilczeć.
 */

import { Command } from '../../../../shared/contract';
import type { NarzedzieAdnotacji } from './slady-adnotacji';

/** Kod modułu z kolumny `modul.kod` rdzenia (`migracja_007_zaczyn_slownikow.sql`). */
export const KOD_MODULU = 'browser';

/**
 * Kody okien operacyjnych modułu. Kod jest bezmodułowy —
 * `migracja_030_rejestr_okien_operacyjnych.sql` nadaje definicjom okien kody bez
 * przedrostka modułu, więc kody przychodzące z `Module.operationalWindowCodes`
 * zestawiają się wprost z tymi wartościami.
 *
 * Rejestr rdzenia niesie dziś wiersze trzech pierwszych okien. Automation Studio
 * i Capture & Monitor Panel opracowanie modułu wymienia na równi z nimi, więc
 * moduł je buduje i podaje ich kody katalogowi — pasek uczciwości wypowiada
 * wtedy rozjazd („kody budowane, których rdzeń modułowi nie przypisuje")
 * zamiast go przemilczeć.
 */
export const KODY_OKIEN = {
  przegladarka: 'browser-window',
  zrodla: 'sources-panel',
  notatki: 'notes-panel',
  automatyzacja: 'automation-studio',
  materialy: 'capture-monitor-panel',
} as const;

/**
 * Kody paneli warstwy czwartej — kluczy stanu odsłonięcia i znaczników `data-panel`.
 *
 * Panele nie są oknami operacyjnymi: rejestr rdzenia ich nie zna i znać nie ma
 * powodu, bo są częścią okien, w których stoją. Kod jest tu wyłącznie nazwą
 * pozycji w stanie warstw widoczności.
 */
export const KODY_PANELI = {
  edytorScenariusza: 'edytor-scenariusza',
  limityPrzebiegu: 'limity-przebiegu',
  narzedziaInspekcyjne: 'narzedzia-inspekcyjne',
  macierzIzolacji: 'macierz-izolacji',
} as const;

/**
 * Klasy własne dymka [?] podawane bibliotecznej fabryce `utworzDymekObjasnienia`.
 *
 * Stoją w jednym miejscu, bo sięga po nie każdy dymek modułu: rozpisane
 * literałem rozjechałyby się przy pierwszej zmianie nazwy klasy w `browser.css`.
 * Fabryka jest wspólna, wygląd znaku pozostaje modułu — pierścień o wymiarze
 * pola wyboru zamiast bibliotecznego kwadratu ikony.
 */
export const KLASY_DYMKA = { powloka: 'mb-dymek', znak: 'mb-dymek__znak' } as const;

/** Objaśnienia dymków [?] przy elementach konfiguracji okien modułu. */
export const OBJASNIENIA = {
  adres:
    'Adres strony przekazywany rdzeniowi komendą browser.navigate. Rdzeń otwiera stronę ' +
    'w oknie przeglądarki tej sesji i odsyła migawkę jej treści.',
  nowaKarta:
    'Zaznaczone przejście otwiera stronę w nowej karcie okna przeglądarki zamiast ' +
    'zastępować kartę bieżącą (pole newTab komendy browser.navigate).',
  zrodloStrony:
    'Migawka pobrana z polem includeHtml niesie źródło strony obok treści renderowanej. ' +
    'Bez niego rdzeń odsyła wyłącznie tekst widoczny.',
  zrzut:
    'Migawka pobrana z polem includeScreenshot niesie odnośnik do zrzutu ekranu strony. ' +
    'Odnośnik wskazuje zasób rdzenia — klient go nie tworzy. Osobny zrzut w czterech trybach ' +
    'wykonuje komenda browser.screenshot.capture; dopóki rdzeń nie ma jej uchwytu, pole ' +
    'screenshotRef wraca puste i okno mówi o tym po każdym pobraniu.',
  kluczowe:
    'Źródło kluczowe idzie do rdzenia z polem key komendy browser.source.add. ' +
    'Kontrakt nie ma komendy zmieniającej istniejące źródło, więc oznaczenie zapisuje się ' +
    'ponownym dodaniem tego samego adresu.',
  powiazanie:
    'Notatka wiąże się ze źródłem polem sourceId komendy browser.note.add. Wykaz do wyboru ' +
    'składa się ze źródeł dodanych w tej sesji przeglądania.',
  cytat:
    'Zaznaczony fragment strony idzie w polu quote komendy browser.note.add — obok treści ' +
    'notatki, nie zamiast niej.',
  przekazanie:
    'Przekazanie idzie komendą context.transfer: rdzeń przenosi komplet kontekstu do modułu ' +
    'docelowego jedną czynnością, zamiast przepisywania treści ręcznie.',
  trybCzytnika:
    'Tryb czytnika pokazuje wyłącznie tekst migawki, bez źródła strony i bez metadanych. ' +
    'Rzecz dzieje się w kliencie — rdzeń nie zna tego przełącznika.',
  podzialEkranu:
    'Podział ekranu ustawia podgląd strony i wyodrębnioną treść obok siebie zamiast jedno ' +
    'pod drugim. Rzecz dzieje się w kliencie — rdzeń nie zna tego przełącznika.',
  adnotacja:
    'Tryb adnotacji kładzie płótno NAD podglądem strony — DOM podglądu zostaje nietknięty. ' +
    'Rysunek nie zmienia ani jednego węzła strony; „Dodaj do rozmowy" spłaszcza go do PNG ' +
    'i wysyła komendą message.send do okna rozmowy tej sesji.',
  automatyka:
    'Scenariusz przeglądania zapisuje się jako automatyka platformy: komenda ' +
    'automation.workflow.save niesie dokładnie to, z czego scenariusz się składa — nazwę, opis ' +
    'i kroki w kolejności wykonania. Nagrywarka kroków ma własną komendę browser.macro.record ' +
    'i wejdzie tą drogą, gdy rdzeń dostanie jej uchwyt.',
  krokZeStrony:
    'Krok powstaje z bieżącej migawki: rodzaj „komenda", komenda browser.navigate, w treści ' +
    'żądania adres strony. Tak zbudowany scenariusz odtwarza obejście stron w tej samej ' +
    'kolejności, w której Operator je otwierał.',
  harmonogram:
    'Cykliczność idzie do rdzenia komendą automation.schedule.set w notacji cron, a odczyt ' +
    'wraca komendą schedule.get. Uruchomienie jednorazowe i po zdarzeniu opisuje wyzwalacz — ' +
    'kontrakt niesie rodzaje cron, webhook, file i condition.',
  edytorScenariusza:
    'Widok tekstowy pokazuje kroki dwiema notacjami. Zapis przyjmuje wyłącznie JSON, bo klient ' +
    'nie niesie czytnika YAML; widok YAML jest odczytem kroków wpisanych w JSON albo ' +
    'odczytanych z rdzenia.',
  limitKrokow:
    'Limit obowiązuje w kliencie przy zapisie: scenariusz dłuższy niż limit nie wychodzi ' +
    'do rdzenia. Utrwalenie granic ma komendy browser.executor.limits.set i .get — do czasu ' +
    'ich obsługi w rdzeniu wartość żyje w tej karcie i ginie wraz z nią.',
  przechwycenie:
    'Przechwycenie pobiera migawkę komendą browser.snapshot.get z polami includeScreenshot ' +
    'i includeHtml. Rodzaj pozycji („zrzut", „archiwum", „treść") bierze się z tego, co rdzeń ' +
    'naprawdę oddał, a nie z tego, o co okno prosiło.',
  monitorZmian:
    'Monitor zapamiętuje treść strony z chwili założenia. Sprawdzenie przechodzi pod jego adres ' +
    'komendą browser.navigate — a więc PRZESTAWIA wspólny podgląd — i zestawia obie treści ' +
    'w kliencie. Sprawdzanie cykliczne prowadzi rdzeń komendą browser.monitor.add; do czasu ' +
    'jej obsługi każde sprawdzenie jest czynnością Operatora.',
  macierzIzolacji:
    'Polityka izolacji obowiązująca dla okna przeglądarki, czytana komendą ' +
    'isolation.policy.preview po rozstrzygnięciu poziomów zasięgu. Panel jest odczytem: ' +
    'zmiany zapisuje okno konfiguracji punktów izolacji.',
  inspekcja:
    'Podgląd źródła i porównanie pracują na tym, co niesie migawka: źródło przychodzi z polem ' +
    'includeHtml, a porównanie zestawia dwie migawki tego samego okna. Drzewo DOM, ruch ' +
    'sieciowy, konsola i emulacja urządzenia mają własne komendy prowadzone protokołem CDP ' +
    '(browser.dom.inspect, browser.network.har, browser.console.read, browser.device.emulate) ' +
    'i wejdą tą drogą, gdy rdzeń dostanie ich uchwyty.',
  szukanieNaStronie:
    'Wyszukiwanie przegląda treść migawki — całość tego, co moduł dostał od rdzenia. ' +
    'Komendy szukania wewnątrz strony kontrakt nie niesie, więc liczone są wiersze treści, ' +
    'a nie trafienia w wyrenderowanym dokumencie.',
  obecnoscWykonawcy:
    'Wykonawca odbiera dokładnie tę migawkę, którą widać w podglądzie. Wskaźnik mówi, z której ' +
    'chwili ona pochodzi — model nie widzi zmian dokonanych na stronie po jej pobraniu.',
  samoczynneZrodla:
    'Przy zaznaczeniu każde nowe przejście dopisuje źródło komendą browser.source.add. ' +
    'Strona odwiedzona ponownie nie zakłada drugiego wpisu: wykaz okna ma już jej adres, ' +
    'więc wpis zostaje scalony przez pominięcie.',
  eksportZrodel:
    'Wypis powstaje w kliencie z wykazu okna i pobiera się jako plik. Kontrakt nie ma komendy ' +
    'eksportu bibliografii — treść jest już w oknie, więc komenda nie jest do tego potrzebna.',
  klasyfikacjaNotatki:
    'Klasyfikacja, wątek i przypięcie mają pola w kontrakcie (BrowserNote) oraz komendę zapisu ' +
    'browser.note.update. Dopóki rdzeń nie ma jej uchwytu, oznaczenie żyje w tej karcie ' +
    'i po jej przeładowaniu notatki wracają bez klasyfikacji.',
  warstwyWidocznosci:
    'Moduł ujawnia funkcje stopniowo: w spoczynku widoczne jest okno przeglądarki, ' +
    'a rozszerzenia boczne otwierają się znacznikiem. Tryb administracyjny odsłania wyzwalacze ' +
    'warstwy czwartej — narzędzia inspekcyjne, macierz izolacji i widok tekstowy scenariusza.',
} as const;

/**
 * Pięć narzędzi paska adnotacji — kod i nazwa mówiona Operatorowi. Kolejność na
 * pasku jest kolejnością tego wykazu.
 */
export const NARZEDZIA_ADNOTACJI: readonly { kod: NarzedzieAdnotacji; nazwa: string }[] = [
  { kod: 'olowek', nazwa: 'Ołówek odręczny' },
  { kod: 'linia', nazwa: 'Linia' },
  { kod: 'prostokat', nazwa: 'Prostokąt' },
  { kod: 'elipsa', nazwa: 'Elipsa' },
  { kod: 'tekst', nazwa: 'Tekst' },
];

/**
 * Cztery barwy adnotacji — żeton motywu i klasa próbki na przycisku.
 *
 * Żeton i klasa stoją w jednym wierszu, bo opisują tę samą barwę dwiema drogami:
 * arkusz maluje próbkę na przycisku regułą `.mb-adnotacja__barwa--…`, a płótno 2D
 * nie zna `var()` i musi dostać żeton po nazwie, żeby rozwiązać go w chwili
 * rysowania. Rozdzielone na dwa wykazy rozjechałyby się przy zmianie palety.
 *
 * Żetony są prymitywne, nie semantyczne: barwa stanu (`--dn-blad-tekst`)
 * przełącza się wraz z motywem, a tusz adnotacji ma być ten sam w PNG wysłanym
 * z motywu jasnego i z ciemnego, bo załącznik ogląda się poza motywem.
 */
export const BARWY_ADNOTACJI: readonly {
  kod: string;
  nazwa: string;
  zeton: string;
  klasa: string;
}[] = [
  { kod: 'czerwien', nazwa: 'Czerwony', zeton: '--dn-czerwien-400', klasa: 'mb-adnotacja__barwa--czerwien' },
  { kod: 'bursztyn', nazwa: 'Bursztynowy', zeton: '--dn-bursztyn-400', klasa: 'mb-adnotacja__barwa--bursztyn' },
  { kod: 'zielen', nazwa: 'Zielony', zeton: '--dn-zielen-400', klasa: 'mb-adnotacja__barwa--zielen' },
  { kod: 'sygnal', nazwa: 'Błękit sygnałowy', zeton: '--dn-sygnal-500', klasa: 'mb-adnotacja__barwa--sygnal' },
];

/**
 * Zdania trybu adnotacji, w tym zdanie nazywające brak tła.
 *
 * `bezTla` stoi w widoku na stałe, nie w dymku: pod adnotacją miałby leżeć zrzut
 * strony, którego rdzeń nie oddaje. Operator ma to wiedzieć przed wysłaniem,
 * a nie domyślić się z pustego obrazka po fakcie.
 */
export const ADNOTACJA = {
  bezTla:
    'Płótno jest przezroczyste — rdzeń nie oddaje zrzutu strony (pole screenshotRef zostaje ' +
    'puste), więc w załączniku pójdzie sam rysunek, BEZ obrazu strony pod spodem.',

  // Rysunek idzie jako URI danych, rdzeń materializuje go do pliku na własnym
  // nośniku (`core/adapter_rozmowa_zalaczniki.go`) i podaje modelowi ścieżkę
  // w treści zapytania, wraz z prośbą o sięgnięcie po nią narzędziem odczytu.
  // Obrazu wklejonego w wiadomość model nie dostaje — kanał CLI nie ma na to
  // pola — więc zobaczy rysunek dopiero po otwarciu pliku.
  modelDostajeSciezke:
    'Rysunek trafi na nośnik rdzenia, a model dostanie ŚCIEŻKĘ do pliku — zobaczy go dopiero, ' +
    'gdy sięgnie po plik narzędziem odczytu. Obrazu wklejonego w wiadomość model nie dostaje.',

  pusteBezWysylki:
    'Nie ma czego dodać do rozmowy — na płótnie nie ma ani jednego śladu. Narysuj oznaczenie ' +
    'jednym z narzędzi paska.',
  brakSplaszczenia:
    'Spłaszczenie do PNG nic nie oddało (canvas.toDataURL zwrócił pustkę) — obraz adnotacji ' +
    'nie powstał, więc nie ma czego wysłać do rozmowy.',
  wysylkaWToku: 'Spłaszczona adnotacja idzie do okna rozmowy komendą message.send…',
} as const;

/**
 * Cztery klasyfikacje notatki z opracowania modułu wraz z pozycją „bez
 * klasyfikacji", od której zaczyna każda notatka.
 *
 * Kod jest kluczem pamięci widoku, nazwa — napisem dla Operatora. Kontrakt pola
 * klasyfikacji nie niesie, więc oznaczenie żyje w karcie (`zebrane-w-sesji.ts`)
 * i panel mówi o tym wprost.
 */
export const KLASYFIKACJE: readonly { kod: string; nazwa: string }[] = [
  { kod: '', nazwa: 'bez klasyfikacji' },
  { kod: 'obserwacja', nazwa: 'obserwacja' },
  { kod: 'cytat', nazwa: 'cytat' },
  { kod: 'pytanie-otwarte', nazwa: 'pytanie otwarte' },
  { kod: 'wniosek', nazwa: 'wniosek' },
];

/** Nazwa klasyfikacji po kodzie; kod nieznany przedstawia się sam sobą. */
export function nazwaKlasyfikacji(kod: string): string {
  return KLASYFIKACJE.find((pozycja) => pozycja.kod === kod)?.nazwa ?? kod;
}

/**
 * Pozycje, których okno jeszcze nie wykonuje: nazwa czynności i komenda, która
 * ją wykona.
 *
 * Zdania o powodzie nie ma tutaj ani jednego. Powód rozstrzyga się przy oknie,
 * z odczytu wykazu komend rdzenia (`moduly/pokrycie-komend.ts`): kontrakt
 * komendę niesie, a rdzeń może mieć albo nie mieć jej uchwytu — i to się zmienia
 * wraz z rdzeniem, nie wraz z tym plikiem. Zdanie wpisane tu na sztywno
 * przestałoby być prawdziwe w dniu dobudowy obsługi i nikt by go nie zdjął.
 *
 * Nazwa komendy pochodzi wyłącznie ze stałych kontraktu. Napis powielony
 * w module przeżyłby zmianę nazwy w `contract.json` i zostawiłby w oknie
 * zdanie o komendzie, której już nie ma.
 */
export interface PozycjaBezObslugi {
  /** Czym pozycja jest dla czytającego — wchodzi w zdanie powodu. */
  czynnosc: string;
  /** Komenda, która tę pozycję wykona. */
  komenda: string;
}

export const POZYCJE_BEZ_OBSLUGI = {
  usuniecieZrodla: { czynnosc: 'Usunięcie źródła', komenda: Command.BrowserSourceRemove },
  zmianaNotatki: { czynnosc: 'Zmiana zapisanej notatki', komenda: Command.BrowserNoteUpdate },
  zakladka: { czynnosc: 'Zakładka strony', komenda: Command.BrowserBookmarkAdd },
  makro: { czynnosc: 'Nagrywanie makra przeglądania', komenda: Command.BrowserMacroRecord },
  pobrania: { czynnosc: 'Menedżer pobrań', komenda: Command.BrowserDownloadList },
  przewijanie: {
    czynnosc: 'Przewinięcie strony po stronie rdzenia',
    komenda: Command.BrowserScroll,
  },
  // „Dodaj do rozmowy" ma dać dwie rzeczy naraz: załącznik rozmowy oraz wpis
  // w wytworach sesji. Pierwszą niesie `message.send` polem `attachments`,
  // drugą `browser.artifact.add` — obie komendy kontrakt ma, więc okno pyta
  // o ich pokrycie w rdzeniu, zamiast orzekać o braku.
  wytworAdnotacji: {
    czynnosc: 'Wpis adnotacji w wytworach sesji',
    komenda: Command.BrowserArtifactAdd,
  },
  monitorCykliczny: {
    czynnosc: 'Cykliczne sprawdzanie monitora po stronie rdzenia',
    komenda: Command.BrowserMonitorAdd,
  },
  kanalyRss: { czynnosc: 'Subskrypcja kanału RSS albo Atom', komenda: Command.BrowserFeedSubscribe },
  kolejkaCzytania: {
    czynnosc: 'Kolejka czytania z przypomnieniem',
    komenda: Command.BrowserReadlistAdd,
  },
  drzewoDom: { czynnosc: 'Podgląd drzewa DOM i stylów', komenda: Command.BrowserDomInspect },
  ruchSieciowy: {
    czynnosc: 'Rejestr żądań sieciowych i eksport HAR',
    komenda: Command.BrowserNetworkHar,
  },
  konsolaStrony: { czynnosc: 'Odczyt konsoli i błędów strony', komenda: Command.BrowserConsoleRead },
  emulacjaUrzadzenia: {
    czynnosc: 'Emulacja urządzenia w podglądzie',
    komenda: Command.BrowserDeviceEmulate,
  },
  limityWykonawcy: {
    czynnosc: 'Utrwalenie limitów kroków Wykonawcy',
    komenda: Command.BrowserExecutorLimitsSet,
  },
  zrzutStrony: { czynnosc: 'Zrzut ekranu strony', komenda: Command.BrowserScreenshotCapture },
} as const satisfies Record<string, PozycjaBezObslugi>;

/**
 * Stany puste trzech okien modułu — tytuł i zdanie mówiące, czym okno jest i jak
 * je zapełnić.
 *
 * Stan pusty opisuje sytuację oczekiwaną, czyli pierwsze użycie, a nie awarię
 * odczytu; napis „Brak danych" nie mówi ani czym okno jest, ani co Operator ma
 * zrobić. Dlatego opisu żąda już fabryka `utworzStanOkna` — okna bez zdania
 * o sobie nie da się zbudować.
 */
export const STANY_PUSTE = {
  przegladarka: {
    tytul: 'Podgląd strony jest pusty',
    opis:
      'Browser Window pokazuje treść strony pobranej przez rdzeń — Operatorowi i modelowi ' +
      'naraz. Wpisz adres w rzędzie nawigacji powyżej, a migawka stanie w tym miejscu.',
  },
  zrodla: {
    tytul: 'Wykaz źródeł jest pusty',
    opis:
      'Sources Panel gromadzi adresy, na które powołuje się sesja, i zasila Sources Manager ' +
      'modułu Research. Zapełnij go przyciskiem „Dodaj z bieżącej strony" albo wpisz adres ' +
      'w polach powyżej.',
  },
  notatki: {
    tytul: 'Wykaz notatek jest pusty',
    opis:
      'Notes Panel trzyma ustalenia spisane przy oglądanych stronach i zasila Library ' +
      'Explorer. Zaznacz fragment w Browser Window i naciśnij „Notatka" albo wypełnij ' +
      'formularz powyżej.',
  },
  automatyzacja: {
    tytul: 'Wykaz scenariuszy jest pusty',
    opis:
      'Automation Studio prowadzi powtarzalne obejścia stron: kroki, harmonogram i przebiegi. ' +
      'Zbuduj scenariusz przyciskiem „Krok z bieżącej strony" i zapisz go pod nazwą albo ' +
      'odczytaj scenariusze zapisane wcześniej.',
  },
  materialy: {
    tytul: 'Materiał sesji jest pusty',
    opis:
      'Capture & Monitor Panel gromadzi zrzuty, archiwa stron i monitory zmian tej sesji ' +
      'przeglądania. Naciśnij „Przechwyć bieżącą stronę", a pozycja stanie w tym wykazie.',
  },
} as const;

/**
 * Zdania wypowiadane przy wykazach zaciąganych z rdzenia.
 *
 * Wykaz ma powiedzieć, skąd pochodzi. Panel pokazujący pozycje bez słowa
 * o pochodzeniu wygląda tak samo, gdy czyta rdzeń, i wtedy, gdy pokazuje własną
 * pamięć — a to dwie różne obietnice wobec Operatora.
 */
export const WYKAZY = {
  zrodla:
    'Wykaz zaciągany z rdzenia komendą browser.source.list — pokazuje źródła okna, ' +
    'nie tylko dodane w tej karcie.',
  notatki:
    'Wykaz zaciągany z rdzenia komendą browser.note.list — pokazuje notatki okna, ' +
    'nie tylko dodane w tej karcie.',
  automatyki:
    'Wykaz zaciągany z rdzenia komendą automation.workflow.list — pokazuje automatyki ' +
    'Operatora, także te założone poza modułem Browser.',
  material:
    'Materiał gromadzi się w karcie: pozycją jest migawka, którą oddał rdzeń. Zapis pozycji ' +
    'jako wytworu sesji ma komendę browser.artifact.add; komendy odczytu wykazu wytworów okna ' +
    'kontrakt nie niesie, więc po przeładowaniu karty wykaz zaczyna się od nowa, a same ' +
    'migawki zostają w rdzeniu.',
} as const;

/**
 * Nazwy wyzwalaczy paska kontekstu wraz z warstwą, na której stoją.
 *
 * Brzmienie jest brzmieniem opracowania modułu (rozdz. 3.1): znaczniki `Źródła ▼`
 * i `Notatki ▼` należą do warstwy drugiej, menu `Operacje ▼` do trzeciej,
 * a narzędzia warstwy czwartej nie mają w pasku nazwy, dopóki Operator nie
 * włączy trybu administracyjnego.
 */
export const WYZWALACZE = {
  zrodla: 'Źródła',
  notatki: 'Notatki',
  automatyzacja: 'Automation Studio',
  materialy: 'Capture & Monitor Panel',
  inspekcja: 'Narzędzia inspekcyjne',
  izolacja: 'Macierz izolacji sesji',
  trybAdministracyjny: 'Tryb administracyjny',
  operacje: 'Operacje',
} as const;
