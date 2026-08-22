# Danaco Console — Architektura techniczna

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Architektura techniczna |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (źródło prawdy techniczne) · projektant (model konfiguracji i warstwy widoczności) |
| **Przeznaczenie** | Ustala model wdrożenia, stos technologiczny, warstwy systemu i wszystkie mechanizmy przekrojowe (komunikacja, sesje, obliczenia, dane, izolacja, rozszerzenia, uwierzytelnianie, synchronizacja), na których opierają się pozostałe opracowania zbioru |
| **Zakres** | architektura techniczna rdzenia i klienta; nie obejmuje wyglądu interfejsu ani zachowania pojedynczych okien |
| **Poza zakresem** | wygląd i zachowanie interfejsu — [Elementy okien](../interfejs-uzytkownika/elementy-okien.md), [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md); pełny model danych — [Model danych](model-danych.md); pełne kontrakty komunikacji — [Kontrakty komunikacji](kontrakty-komunikacji.md) |
| **Dokument nadrzędny** | [Koncepcja platformy](koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Model danych](model-danych.md) · [Model konfiguracji](model-konfiguracji.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Izolacja i zależności](izolacja-i-zaleznosci.md) · [Integracja modeli](integracja-modeli.md) · [Rozszerzenia](rozszerzenia.md) · [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) |
| **Źródła normatywne** | `budowa/shared/contract.json` · `budowa/go.mod` · `budowa/client/package.json` |
| **Zasada nadrzędna** | Każdy mechanizm przekrojowy platformy ma dokładnie jedno miejsce odniesienia w tym dokumencie; opracowania modułowe i środowiskowe go przywołują, nie powielają |

Dokument opisuje architekturę techniczną produktu Danaco Console: model wdrożenia, stos technologiczny, warstwy systemu, warstwę komunikacji operacyjnej, warstwy widoczności interfejsu, model sesji i procesów, realizację zadań obliczeniowych, integrację modeli, warstwę danych, komunikację, izolację, konfigurację, rozszerzenia, uwierzytelnianie oraz synchronizację.

---

## Spis treści

1. [Zasady architektoniczne](#1-zasady-architektoniczne)
2. [Model wdrożenia](#2-model-wdrożenia)
3. [Stos technologiczny](#3-stos-technologiczny)
4. [Warstwy systemu](#4-warstwy-systemu)
5. [Warstwa komunikacji operacyjnej](#5-warstwa-komunikacji-operacyjnej)
   - [5.1 Kanał pierwszy — Chat Window (Użytkownik ↔ Wykonawca)](#51-kanał-pierwszy--chat-window-użytkownik--wykonawca)
   - [5.2 Kanał drugi — Execution Loop Window (Koordynator ↔ Wykonawca)](#52-kanał-drugi--execution-loop-window-koordynator--wykonawca)
   - [5.3 Komponenty realizujące kanały](#53-komponenty-realizujące-kanały)
   - [5.4 Przepływ komunikatów](#54-przepływ-komunikatów)
   - [5.5 Odwzorowanie w kontraktach komunikacji i modelu danych](#55-odwzorowanie-w-kontraktach-komunikacji-i-modelu-danych)
6. [Warstwy widoczności interfejsu jako zasada architektoniczna](#6-warstwy-widoczności-interfejsu-jako-zasada-architektoniczna)
7. [Model sesji i procesów](#7-model-sesji-i-procesów)
8. [Realizacja zadań obliczeniowych i granice wykonania](#8-realizacja-zadań-obliczeniowych-i-granice-wykonania)
9. [Integracja modeli](#9-integracja-modeli)
   - [9.1 Wykaz komend integracji modeli](#91-wykaz-komend-integracji-modeli)
   - [9.2 Diagram sekwencji — wywołanie modelu przez kanał komunikacji operacyjnej](#92-diagram-sekwencji--wywołanie-modelu-przez-kanał-komunikacji-operacyjnej)
10. [Warstwa danych](#10-warstwa-danych)
   - [10.1 Schemat encji i relacji (skrócony)](#101-schemat-encji-i-relacji-skrócony)
11. [Komunikacja klient–serwer](#11-komunikacja-klientserwer)
12. [Izolacja i konfigurowalność zależności](#12-izolacja-i-konfigurowalność-zależności)
   - [12.1 Diagram sekwencji — przypisanie profilu izolacji do sesji](#121-diagram-sekwencji--przypisanie-profilu-izolacji-do-sesji)
   - [12.2 Wykaz komend obszaru `isolation`](#122-wykaz-komend-obszaru-isolation)
13. [Konfiguracja](#13-konfiguracja)
14. [Rozszerzenia](#14-rozszerzenia)
   - [14.1 Obszar kontraktu `extension`](#141-obszar-kontraktu-extension)
   - [14.2 Diagram sekwencji — instalacja rozszerzenia źródła Personal](#142-diagram-sekwencji--instalacja-rozszerzenia-źródła-personal)
15. [Uwierzytelnianie](#15-uwierzytelnianie)
   - [15.1 Diagram sekwencji — logowanie i nawiązanie połączenia](#151-diagram-sekwencji--logowanie-i-nawiązanie-połączenia)
   - [15.2 Obszar kontraktu `auth`](#152-obszar-kontraktu-auth)
16. [Synchronizacja](#16-synchronizacja)
17. [Widok wdrożeniowy](#17-widok-wdrożeniowy)
   - [17.1 Drzewo katalogów wdrożenia](#171-drzewo-katalogów-wdrożenia)
18. [Wykaz obszarów kontraktu](#18-wykaz-obszarów-kontraktu)
19. [Kody błędów kontraktu](#19-kody-błędów-kontraktu)
20. [Kryteria odbioru](#20-kryteria-odbioru)
21. [Załącznik A. Scenariusze eksploatacji](#załącznik-a-scenariusze-eksploatacji)
   - [A.1. Przetrwanie sesji po rozłączeniu klienta](#a1-przetrwanie-sesji-po-rozłączeniu-klienta)
   - [A.2. Przebieg zlecenia przez oba kanały komunikacji operacyjnej](#a2-przebieg-zlecenia-przez-oba-kanały-komunikacji-operacyjnej)
   - [A.3. Synchronizacja stanu między urządzeniami](#a3-synchronizacja-stanu-między-urządzeniami)
   - [A.4. Dodanie nowego kanału modelu](#a4-dodanie-nowego-kanału-modelu)
   - [A.5. Wykonanie zadania obliczeniowego wymagającego akceleracji](#a5-wykonanie-zadania-obliczeniowego-wymagającego-akceleracji)
   - [A.6. Zastosowanie profilu izolacji do sesji](#a6-zastosowanie-profilu-izolacji-do-sesji)
   - [A.7. Instalacja rozszerzenia z zasobów Operatora](#a7-instalacja-rozszerzenia-z-zasobów-operatora)
   - [A.8. Dotarcie do funkcji eksperckiej](#a8-dotarcie-do-funkcji-eksperckiej)
   - [A.9. Ponowienie żądania po błędzie ponawialnym](#a9-ponowienie-żądania-po-błędzie-ponawialnym)
22. [Załącznik B. Słownik terminów architektonicznych](#załącznik-b-słownik-terminów-architektonicznych)

---

## 1. Zasady architektoniczne

Architektura opiera się na ośmiu zasadach nadrzędnych. Każda z nich odzwierciedla się w konkretnych rozstrzygnięciach dalszych rozdziałów.

| Zasada | Treść | Odzwierciedlenie w architekturze |
|---|---|---|
| Prymat funkcji | Architektura umożliwia wykonanie działania; nie tworzy ograniczeń bez uzasadnienia produktowego. | Izolacja jako możliwość, nie wymóg (rozdz. 12) |
| Cienki punkt wejścia | Punkt startowy aplikacji wyłącznie komponuje moduły; logika mieszka w modułach. | Warstwy systemu (rozdz. 4) |
| Pojedyncza odpowiedzialność | Każdy moduł kodu odpowiada za jeden obszar. | Warstwy systemu (rozdz. 4), integracja modeli (rozdz. 9) |
| Rozdzielenie przez kontrakty | Moduły łączą się zdefiniowanymi interfejsami, nie bezpośrednimi zależnościami. | Komunikacja (rozdz. 11), integracja modeli (rozdz. 9), rozszerzenia (rozdz. 14) |
| Konfigurowalność zamiast wartości wpisanych na stałe | Parametry środowiska i zachowania pochodzą z konfiguracji; brak ustawienia oznacza wartość domyślną. | Konfiguracja (rozdz. 13), izolacja (rozdz. 12) |
| Pojedyncze źródło prawdy | Stan systemu przechowuje serwer. | Warstwa danych (rozdz. 10), synchronizacja (rozdz. 16) |
| Ciężar wykonania po stronie serwera | Praca obliczeniowa, dostęp do plików, procesy i wywołania modeli wykonują się w rdzeniu serwera; klient prezentuje obraz. | Realizacja zadań obliczeniowych (rozdz. 8), model wdrożenia (rozdz. 2) |
| Sterowanie przez kanały komunikacji operacyjnej | Procesy platformy uruchamiane są i nadzorowane przez dwa kanały komunikacji, nie przez rozproszone panele sterowania. | Warstwa komunikacji operacyjnej (rozdz. 5) |

---

## 2. Model wdrożenia

System działa w modelu hybrydowym: przetwarzanie po stronie serwera napisanego w języku Go, dostęp przez cienkiego klienta uruchamiającego natywne okno na urządzeniu. Stanowisko robocze Operatora nie zawiera lokalnego układu GPU i nie jest to wymagane do pracy platformy w pełnym zakresie.

```
MODEL HYBRYDOWY

   Przetwarzanie, logika, stan   ─────►   SERWER WYKONAWCZY (Go)
                                          (serwer wirtualny)
   Dostęp, prezentacja, obraz    ─────►   CIENKI KLIENT URZĄDZENIA
                                          (natywne okno na urządzeniu Operatora)
```

**Podział odpowiedzialności między warstwy wdrożenia:**

| Cecha | Serwer wykonawczy | Cienki klient urządzenia |
|---|---|---|
| Zawartość | Rdzeń aplikacji w języku Go oraz procesy sesji | Lekki pakiet startowy uruchamiający natywne okno aplikacji |
| Umiejscowienie | Serwer wirtualny | Urządzenie Operatora |
| Logika biznesowa | Tak | Nie |
| Praca obliczeniowa | Tak | Nie |
| Stan trwały | Tak (jedyne źródło prawdy) | Nie |
| Rola w połączeniu | Przyjmuje połączenia klientów | Łączy się z serwerem |

**Charakterystyka wdrożenia:**

| Cecha | Wartość |
|---|---|
| Użytkownik | Jeden (Operator). |
| Urządzenia | Kilka: komputery, telefony, tablety. |
| Wymagania sprzętowe stanowiska | Brak wymogu lokalnego układu GPU; wystarcza urządzenie zdolne wyświetlić natywne okno klienta. |
| Instalacja klienta | Jednorazowa na urządzenie. |
| Aktualizacja rdzenia | Centralna, po stronie serwera. |

---

## 3. Stos technologiczny

Każdy komponent stosu przypisany jest do warstwy systemu, w której działa (rozdz. 4).

| Komponent | Warstwa systemu | Technologia | Uzasadnienie |
|---|---|---|---|
| Rdzeń serwera | Rdzeń | Go | Orkiestracja, współbieżność, zarządzanie procesami sesji. |
| Okno klienta | Klient | Tauri (warstwa natywna w języku Rust) | Natywne okno o niskiej wadze, wieloplatformowość. |
| Silnik prezentacji | Klient | Chromium (webview platformy) | Standard renderowania osadzany przez platformę. |
| Interfejs | Klient | TypeScript (aplikacja webowa) | Wspólny kod prezentacji dla wszystkich platform. |
| Kanał komunikacji | Komunikacja | WebSocket, komunikaty w formacie JSON | Stałe, dwukierunkowe połączenie klient–serwer ze strumieniem na żywo (rozdz. 11). |
| Przechowywanie danych | Dane | SQLite | Trwałość transakcyjna w pojedynczym pliku. |

**Wersje wiążące stosu (odczyt z plików budowy, na dzień 2026-08-20):**

| Komponent | Plik źródłowy | Wersja |
|---|---|---|
| Moduł Go rdzenia | `budowa/go.mod` | `go 1.26.4` |
| TypeScript klienta | `budowa/client/package.json` | `5.8.3` |

**Silnik webview zależnie od platformy:**

| Platforma | Silnik webview |
|---|---|
| Windows | WebView2 lub osadzony Chromium. |
| Android | Systemowy komponent WebView. |
| iOS oraz iPadOS | Komponent WKWebView. |

---

## 4. Warstwy systemu

System dzieli się na cztery warstwy techniczne. Warstwy komunikują się w sąsiedztwie: klient przez warstwę komunikacji do rdzenia, rdzeń do warstwy danych; klient nie sięga bezpośrednio do danych.

```
┌──────────────────────────────────────────────────────────────────────┐
│ KLIENT          natywne okno · interfejs · interakcja Operatora      │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │  WebSocket (JSON, dwukierunkowo, strumień na żywo)
┌─────────────────────────────────┴────────────────────────────────────┐
│ KOMUNIKACJA     kanał dwukierunkowy klient–serwer · strumień na żywo │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
┌─────────────────────────────────┴────────────────────────────────────┐
│ RDZEŃ           orkiestracja sesji · Koordynator · kolejka zadań ·   │
│                 procesy · integracja modeli · polityki izolacji      │
└─────────────────────────────────┬────────────────────────────────────┘
                                  │
┌─────────────────────────────────┴────────────────────────────────────┐
│ DANE            SQLite · pliki treści · jedno źródło prawdy          │
└──────────────────────────────────────────────────────────────────────┘
        (klient nie sięga bezpośrednio do warstwy danych)
```

| Warstwa | Odpowiedzialność |
|---|---|
| Klient | Prezentacja, interakcja Operatora, natywne okno, obsługa okien komunikacji operacyjnej. |
| Komunikacja | Dwukierunkowy kanał klient–serwer, strumień na żywo, przenoszenie komunikatów obu kanałów operacyjnych. |
| Rdzeń | Orkiestracja sesji, Koordynator, kolejka zadań, kontrola jakości, procesy, integracja modeli, polityki izolacji. |
| Dane | Trwałe przechowywanie stanu, w tym zleceń, zadań i komunikatów obu kanałów. |

**Realizacja trzech warstw koncepcji platformy w warstwach technicznych:**

| Warstwa platformy | Realizacja w architekturze technicznej |
|---|---|
| Warstwa środowisk | Rdzeń — zarządzanie sesjami środowisk. |
| Warstwa modułów | Rdzeń — zarządzanie sesjami modułów; klient — budowa przestrzeni roboczej właściwej aktualnemu modułowi. |
| Warstwa funkcji globalnych | Realizowana w obrębie warstwy klienta i rdzenia, ponad podziałem na środowiska i moduły. |

---

## 5. Warstwa komunikacji operacyjnej

Platforma prowadzi dwa kanały komunikacji operacyjnej. Oba stanowią elementy pierwszoplanowe architektury: kanał pierwszy jest podstawowym mechanizmem sterowania procesami platformy, kanał drugi jest kanałem pętli wykonawczej. Każdy kanał ma własne okno w interfejsie, własny zestaw komunikatów w kontrakcie komunikacji oraz własne odwzorowanie w modelu danych.

**Role uczestniczące:**

| Rola | Charakter | Zakres odpowiedzialności |
|---|---|---|
| Użytkownik | Operator platformy | Zleca działania, zatwierdza, przerywa, koryguje zlecenie. |
| Koordynator | Komponent orkiestrujący rdzenia | Dekomponuje zlecenie na zadania, przydziela je, nadzoruje realizację, ocenia wynik, decyduje o ponowieniu. |
| Wykonawca | AI, agent lub system wykonawczy | Realizuje przydzielone zadania i zwraca wyniki. |

### 5.1. Kanał pierwszy — Chat Window (Użytkownik ↔ Wykonawca)

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą oraz centralnym punktem pracy. Przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie i przerywanie działań oraz wyjaśnianie wyniku i kontekstu. Każdy moduł i każde środowisko udostępnia to okno w tym samym miejscu układu — w lewej kolumnie obszaru roboczego, stałej, o pełnej wysokości.

### 5.2. Kanał drugi — Execution Loop Window (Koordynator ↔ Wykonawca)

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą i odpowiada za prowadzenie pętli wykonawczej: koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Okno zawiera bieżące zlecenie wraz z jego dekompozycją na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Okno otwierane jest jako kolumna sąsiadująca z Chat Window.

### 5.3. Komponenty realizujące kanały

| Komponent | Warstwa systemu | Rola w kanałach |
|---|---|---|
| Chat Window | Klient | Prezentuje kanał Użytkownik ↔ Wykonawca; przyjmuje polecenia języka naturalnego i strumień odpowiedzi. |
| Execution Loop Window | Klient | Prezentuje kanał Koordynator ↔ Wykonawca; pokazuje zlecenie, zadania, kontrolę jakości i sterowanie przebiegiem. |
| Koordynator | Rdzeń | Dekomponuje zlecenie, przydziela zadania, nadzoruje pętlę wykonawczą, ocenia wyniki, zarządza ponowieniami. |
| Kolejka zadań | Rdzeń | Przechowuje zadania oczekujące i w toku, ich priorytety, zależności oraz stany przejścia. |
| Kontrola jakości | Rdzeń | Weryfikuje wynik zadania wobec kryteriów zlecenia i zwraca ocenę do Koordynatora. |
| Mechanizm ponowień | Rdzeń | Realizuje ponowienie zadania z narastającym odstępem oraz przekazanie zadania nierozstrzygalnego do decyzji Użytkownika. Liczba prób nie ma granicy narzuconej z góry — ponowienia trwają do powodzenia albo do przerwania na polecenie Operatora, które mechanizm przyjmuje także w trakcie odczekiwania odstępu ([Integracja modeli](integracja-modeli.md) rozdz. 9.2). |
| Procesy wykonawcze | Rdzeń | Realizują zadania jako procesy sesji (rozdz. 7). |
| Kanał WebSocket | Komunikacja | Przenosi komunikaty obu kanałów w obie strony, w tym strumienie na żywo. |
| Rejestr zleceń, zadań i komunikatów | Dane | Przechowuje trwale zlecenia, zadania, komunikaty, wyniki kontroli jakości i historię ponowień. |

### 5.4. Przepływ komunikatów

```
UŻYTKOWNIK
   │  polecenie w języku naturalnym (Chat Window)
   ▼
RDZEŃ — KOORDYNATOR
   │  dekompozycja zlecenia na zadania  ──►  KOLEJKA ZADAŃ
   │                                            │  przydział zadania
   │                                            ▼
   │                                        WYKONAWCA (proces sesji)
   │                                            │  wynik zadania
   │                                            ▼
   │                                        KONTROLA JAKOŚCI
   │        ocena negatywna ──► ponowienie ──► KOLEJKA ZADAŃ
   │        ocena pozytywna ──► domknięcie zadania
   ▼
STRUMIEŃ ZDARZEŃ (WebSocket)
   ├──►  Chat Window            — odpowiedzi i wyniki dla Użytkownika
   └──►  Execution Loop Window  — stan zlecenia, zadań, kontroli jakości, ponowień
```

Sterowanie zwrotne przebiega tą samą drogą: polecenia wstrzymania, wznowienia, przerwania i korekty zlecenia trafiają z okna klienta do Koordynatora, który przekłada je na zmianę stanu zadań w kolejce oraz na sygnały dla procesów wykonawczych.

### 5.5. Odwzorowanie w kontraktach komunikacji i modelu danych

| Element kanału | Kategoria komunikatu (rozdz. 11) | Odwzorowanie w warstwie danych (rozdz. 10) |
|---|---|---|
| Polecenie Użytkownika | Polecenie (klient → serwer) | Wpis zlecenia powiązany z sesją. |
| Odpowiedź Wykonawcy | Strumień oraz zdarzenie (serwer → klient) | Wpis komunikatu kanału Użytkownik ↔ Wykonawca. |
| Dekompozycja zlecenia | Zdarzenie (serwer → klient) | Zbiór zadań powiązanych ze zleceniem. |
| Zmiana stanu zadania | Zdarzenie (serwer → klient) | Aktualizacja stanu zadania w kolejce. |
| Komunikat sterujący Koordynatora | Zdarzenie (serwer → klient) | Wpis komunikatu kanału Koordynator ↔ Wykonawca. |
| Wynik kontroli jakości | Zdarzenie (serwer → klient) | Wpis oceny powiązany z zadaniem. |
| Decyzja o ponowieniu | Zdarzenie (serwer → klient) | Wpis historii ponowień zadania. |
| Sterowanie przebiegiem | Polecenie (klient → serwer) | Aktualizacja stanu zlecenia i zadań. |

Serwer pozostaje jedynym źródłem prawdy dla obu kanałów; komplet zleceń, zadań i komunikatów utrwalany jest po stronie serwera i przekazywany na wszystkie urządzenia (rozdz. 16).

---

## 6. Warstwy widoczności interfejsu jako zasada architektoniczna

Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności. Jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

Zasada ma charakter architektoniczny: warstwa klienta buduje przestrzeń roboczą wyłącznie z elementów należących do warstw widoczności właściwych bieżącemu kontekstowi, a rdzeń dostarcza definicje elementów wraz z przypisaną warstwą i sposobem wywołania.

**Cztery warstwy widoczności:**

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu. | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

**Komponenty systemu realizujące zasadę:**

| Komponent | Warstwa systemu | Rola |
|---|---|---|
| Rejestr elementów interfejsu | Rdzeń | Przechowuje definicję każdego elementu wraz z przypisaną warstwą widoczności i sposobem wywołania. |
| Konfiguracja warstw wg roli użytkownika | Rdzeń, konfiguracja (rozdz. 13) | Wyznacza zakres warstw udostępnianych danej roli; rola podstawowa nie otrzymuje elementów warstwy 4. |
| Wyszukiwarka funkcji | Klient | Udostępnia każdą funkcję platformy po nazwie i opisie, niezależnie od warstwy, i uruchamia ją bezpośrednio z wyniku wyszukiwania. |
| Kompozytor przestrzeni roboczej | Klient | Buduje widok z elementów warstwy 1 oraz zwiniętych wyzwalaczy warstw 2–3 właściwych kontekstowi. |
| Kanał poleceń języka naturalnego | Komunikacja, rdzeń | Uruchamia funkcje warstwy 4 poleceniem wydanym w Chat Window. |

**Mechanizmy ukrywania funkcjonalności:** menu progresywne (zbiór jednorodnych wyborów jako jeden element zwinięty, `Agent ▼`), panele wysuwane znikające całkowicie po zamknięciu, grupowanie logiczne akcji (`Operacje ▼`), znaczniki kontekstowe w pasku kontekstu w postaci `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`, których kliknięcie otwiera selektor wartości danego znacznika.

**Zasada jednego kliknięcia.** Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

**Układ obszaru roboczego** jest wyłącznie pionowy, w podziale lewa–prawa; regulacji podlega wyłącznie szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Lewa, stała, pełna wysokość | Chat Window — kanał Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca (otwierana) | Execution Loop Window — kanał Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu (okna edycyjne, podglądu, monitory) |
| Kolejne kolumny boczne | Okna pomocnicze i panele — otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego |

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │                          │  boczne)
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ════════════════════════════════════════════════════════════════════════════
```

---

## 7. Model sesji i procesów

Każda sesja działa jako odrębny, izolowany proces po stronie serwera, niezależny od obecności klienta.

| Właściwość procesu sesji | Zachowanie |
|---|---|
| Model wykonania | Każda sesja działa jako odrębny proces po stronie serwera. |
| Izolacja procesu | Własny katalog roboczy i środowisko procesu. |
| Rejestr procesów | Rdzeń prowadzi rejestr aktywnych procesów, kluczowany identyfikatorem uruchomienia procesu. |
| Rozłączenie klienta | Nie kończy sesji — proces pracuje dalej po stronie serwera. |
| Ponowne połączenie | Podłącza klienta do aktywnej sesji. |
| Trwałość stanu | Stan sesji jest trwały; restart serwera go nie traci. |
| Zakończenie procesu | Pozwala wznowić sesję z zapisu. |
| Współbieżność | Sesje wykonują się równolegle i niezależnie. |

**Cykl życia sesji:**

```
   UTWORZENIE SESJI
        │  rdzeń uruchamia odrębny proces (własny katalog roboczy, środowisko procesu)
        │  wpis do rejestru aktywnych procesów — klucz: identyfikator uruchomienia procesu
        ▼
   ┌──────────────────┐   rozłączenie klienta     ┌──────────────────────────────┐
   │  AKTYWNA         │ ───────────────────────►  │  AKTYWNA BEZ KLIENTA         │
   │  klient          │   proces pracuje dalej    │  proces po stronie serwera   │
   │  podłączony      │ ◄───────────────────────  │  pracuje dalej               │
   └────────┬─────────┘   ponowne połączenie      └──────────────┬───────────────┘
            │           (podłączenie klienta do aktywnej sesji)  │
            │                                                    │
            │  zakończenie procesu / restart serwera             │
            ▼                                                    ▼
   ┌───────────────────────────────────────────────────────────────────────┐
   │  ZAPIS TRWAŁY   stan sesji zachowany (SQLite + pliki treści)          │
   └───────────────────────────────────┬───────────────────────────────────┘
                                       │  wznowienie sesji z zapisu
                                       ▼
                                (powrót do stanu AKTYWNA)
```

**Szablon profilu procesu sesji (pola konfiguracji):**

```
proces_sesji:
  identyfikator_uruchomienia:   <klucz w rejestrze aktywnych procesów>
  katalog_roboczy:              <ścieżka — izolacja procesu>
  srodowisko_procesu:           <zmienne środowiskowe procesu>
  model_procesu:                odrebny | wspoldzielony
  trwalosc_stanu:               wlaczona          # restart serwera nie traci stanu
  wznowienie_z_zapisu:          dozwolone
```

Sesje odpowiadają kartom sesji w oknie środowiska. Środowisko MultitaskingAI prowadzi wiele procesów jednocześnie — osobny proces dla każdej roli wykonawczej oraz procesy podagentów uruchamiane przez rolę Executor.

```
Środowisko MultitaskingAI — wiele procesów jednocześnie
   proces: Executor 1                (rola wykonawcza)
   proces: Executor 2                (rola wykonawcza)
   proces: Coordinator               (rola)
   proces: Executor 3 / Validator    (rola)
      └── procesy podagentów uruchamiane przez rolę Executor (Subagent Network)
```

---

## 8. Realizacja zadań obliczeniowych i granice wykonania

Architektura hybrydowa wyznacza jednoznaczne granice wykonania. Rdzeń w języku Go wykonuje całą pracę obliczeniową; cienki klient odpowiada wyłącznie za prezentację i interakcję.

**Granice wykonania:**

| Rodzaj pracy | Miejsce wykonania |
|---|---|
| Dostęp do plików, procesy systemowe, operacje na bazie SQLite | Rdzeń serwera. |
| Wywołania modeli i usług zewnętrznych, kolejki, harmonogramy | Rdzeń serwera. |
| Przetwarzanie dokumentów, konwersje, ekstrakcja tekstu, indeksowanie, wyszukiwanie | Rdzeń serwera. |
| Renderowanie interfejsu, interakcja wskaźnikiem i klawiaturą, edytory oraz podglądy w oknie | Cienki klient. |
| Stan trwały | Wyłącznie rdzeń serwera. |

**Zasady realizacji zadań obliczeniowych:**

| Zasada | Treść |
|---|---|
| Brak wymogu GPU po stronie stanowiska | Żadna funkcja platformy nie wymaga lokalnego układu GPU na urządzeniu Operatora. |
| Obliczenia neuronowe przez kanał modelu | Generowanie i przekształcanie obrazu, upscaling, usuwanie tła, segmentacja, synteza i klonowanie głosu, diaryzacja oraz transkrypcja realizowane są jako wywołania kanału modelu (rozdz. 9) albo zdalnego punktu wykonania wyposażonego w GPU. |
| Przetwarzanie lokalne w trybie CPU | Zadania wykonywane w procesie rdzenia — rozpoznawanie tekstu, transkrypcja lekka, transkodowanie, indeksowanie — realizowane są wyłącznie na procesorze. |
| Brak ukrytych zależności sprzętowych | Proces rdzenia nie osadza w sobie inferencji wymagającej GPU; zależność od akceleracji sprzętowej występuje wyłącznie po stronie zdalnego punktu wykonania. |
| Nadzór nad zasobami zewnętrznymi | Dane dostępowe zdalnych punktów wykonania i usług modelowych pochodzą z konfiguracji (rozdz. 13). |

**Zasady doboru bibliotek:**

| Zasada | Treść |
|---|---|
| Pierwszeństwo bibliotek w języku Go | Funkcję realizuje biblioteka w języku Go działająca w procesie rdzenia. |
| Zależności natywne w trybie CPU | Biblioteki z wiązaniem natywnym oraz narzędzia uruchamiane jako procesy zewnętrzne działają na procesorze i wchodzą w skład obrazu środowiska serwera. |
| Praca w webview po stronie klienta | Biblioteki prezentacji działają w webview klienta i obejmują wyłącznie renderowanie oraz interakcję, bez logiki biznesowej i bez dostępu do stanu trwałego. |
| Granica integracji systemowej | Funkcje wymagające przechwytywania zdarzeń systemu operacyjnego poza oknem aplikacji realizowane są w granicach platformy — wewnątrz jej okien i pól. |
| Jednolitość silnika prezentacji | Funkcje sterowania stroną i przechwytywania jej zawartości realizuje rdzeń przez protokół sterowania przeglądarką, po stronie serwera. |

---

## 9. Integracja modeli

Modele udostępniane są przez jednolitą warstwę dostawcy modelu. Rdzeń operuje na abstrakcji „wyślij zapytanie, odbierz odpowiedź” i nie zależy od kanału; każdy kanał integracji realizuje osobny adapter tej warstwy.

```
RDZEŃ
  │  abstrakcja: „wyślij zapytanie, odbierz odpowiedź”  (rdzeń nie zależy od kanału)
  ▼
WARSTWA DOSTAWCY MODELU  (jednolita)
  ├── adapter kanału API   ──►  klucz dostępu · żądanie HTTP
  ├── adapter kanału CLI   ──►  token narzędzia wiersza polecenia
  ├── adapter kanału SSH   ──►  host zdalny · własne modele pomocnicze
  └── adapter kanału HTTP  ──►  modele przeglądarkowe i zdalne punkty wykonania

  Dodanie kanału = dostarczenie nowego adaptera, bez zmian w rdzeniu.
```

**Kanały integracji modeli:**

| Kanał | Sposób połączenia | Rodzaj modeli |
|---|---|---|
| API | Klucz dostępu, żądanie HTTP. | Modele dostępne przez interfejs programistyczny. |
| CLI | Token narzędzia wiersza polecenia. | Modele obsługiwane narzędziem wiersza polecenia. |
| SSH | Połączenie z hostem zdalnym. | Własne modele pomocnicze na hoście zdalnym. |
| HTTP | Żądanie HTTP. | Modele przeglądarkowe oraz zdalne punkty wykonania wyposażone w GPU. |

**Zakres konfiguracji:**

| Element konfiguracji | Zakres wyboru |
|---|---|
| Model i kanał | Per sesja lub per rola. |
| Tożsamość modelu | Konfigurowalna: nazwa, persona, prompt systemowy. |

**Szablon kanału i tożsamości modelu (pola konfiguracji):**

```
kanal_modelu:
  kanal:            API | CLI | SSH | HTTP
  dane_dostepowe:   <klucz dostępu | token CLI | dane hosta SSH | adres HTTP>
  zasieg_wyboru:    per_sesja | per_rola
tozsamosc_modelu:
  nazwa:            <nazwa modelu>
  persona:          <persona>
  prompt_systemowy: <prompt systemowy>
```

### 9.1. Wykaz komend integracji modeli

Trzy obszary kontraktu komunikacji niosą komendy warstwy dostawcy modelu: `channel` (rejestr kanałów), `identity` (tożsamość i persona) oraz `model` (przypisanie kanału do okna albo karty sesji). Poniższy wykaz jest pełny dla tych trzech obszarów.

| Komenda | Obszar | Przeznaczenie |
|---|---|---|
| `channel.add` | `channel` | Dodaje wiersz do rejestru kanałów modelu |
| `channel.update` | `channel` | Zmienia wiersz rejestru kanałów |
| `channel.remove` | `channel` | Usuwa wiersz rejestru kanałów |
| `channel.list` | `channel` | Zwraca rejestr kanałów modelu |
| `channel.check` | `channel` | Sprawdza, czy kanał modelu odpowiada; narzędzie pomocnicze, wynik nie warunkuje zapisu ani wysłania |
| `channel.credential.status` | `channel` | Zwraca stan poświadczenia kanału — czy jest ustawione i kiedy zmienione — nigdy treść poświadczenia |
| `identity.category.list` | `identity` | Zwraca katalog kategorii zasad i tożsamości modelu |
| `identity.document.get` | `identity` | Odczytuje treść kategorii zasad zapisaną dla wskazanej osi |
| `identity.document.set` | `identity` | Zapisuje treść kategorii zasad dla wskazanej osi wraz z trybem podania |
| `identity.document.remove` | `identity` | Usuwa treść kategorii zasad; brak zapisu oznacza dziedziczenie z osi szerszej |
| `identity.effective.get` | `identity` | Zwraca nakładkę obowiązującą: tryb, warstwy w kolejności krytyczności i złożony prompt systemowy |
| `model.channel.set` | `model` | Ustala kanał modelu obsługujący okno albo kartę sesji |

Pełny opis czterech kanałów integracji, mechanizmu adaptera, tożsamości i persony modelu oraz przypisania per sesja i per rola zawiera [Integracja modeli](integracja-modeli.md).

### 9.2. Diagram sekwencji — wywołanie modelu przez kanał komunikacji operacyjnej

```
Wykonawca             Warstwa dostawcy modelu        Adapter kanału           Dostawca
   │                          │                            │                     │
   │  jednolite zapytanie     │                            │                     │
   ├─────────────────────────►│                            │                     │
   │                          │  wybór adaptera wg         │                     │
   │                          │  model.channel.set         │                     │
   │                          ├───────────────────────────►│                     │
   │                          │                            │  protokół kanału    │
   │                          │                            │  + dane dostępowe   │
   │                          │                            ├────────────────────►│
   │                          │                            │  odpowiedź          │
   │                          │                            │  w formacie dostawcy│
   │                          │                            │◄────────────────────┤
   │                          │  jednolity format          │                     │
   │                          │  strumienia                │                     │
   │                          │◄───────────────────────────┤                     │
   │  strumień WebSocket      │                            │                     │
   │◄─────────────────────────┤                            │                     │
```

Legenda: sekwencja identyczna dla każdego z czterech kanałów (API, CLI, SSH, HTTP); adapter tłumaczy jednolite zapytanie na protokół właściwy kanałowi i przekształca odpowiedź dostawcy z powrotem do jednolitego formatu strumienia rdzenia (rozdz. 5, rozdz. 11).

---

## 10. Warstwa danych

Serwer stanowi jedyne źródło prawdy. Dane strukturalne i treści obszerne rozdzielone są między bazę i pliki, powiązane odwołaniami.

```
SERWER — jedyne źródło prawdy
  ┌────────────────────────────┐        ┌────────────────────────────┐
  │ SQLite (jeden plik)        │◄──────►│ Pliki treści               │
  │ dane strukturalne,         │ odwoła-│ kanwa, artefakty,          │
  │ transakcje                 │  nia   │ pliki biblioteki           │
  └────────────────────────────┘        └────────────────────────────┘

  Urządzenia: brak stanu trwałego (prezentują wyłącznie obraz z serwera)
```

| Rodzaj danych | Sposób przechowywania | Lokalizacja |
|---|---|---|
| Dane strukturalne | Baza SQLite w jednym pliku, transakcyjnie. | Serwer. |
| Zlecenia, zadania, komunikaty obu kanałów operacyjnych, wyniki kontroli jakości, historia ponowień | Baza SQLite, transakcyjnie. | Serwer. |
| Treści obszerne (kanwa, artefakty, pliki biblioteki) | Pliki, z odwołaniem w bazie. | Serwer (pliki treści). |
| Stan trwały urządzeń | Brak — urządzenia nie przechowują stanu trwałego. | — |

**Operacje na danych** realizuje interfejs rdzenia; wszystkie są dostępne z okna aplikacji.

| Operacja | Zakres | Dostępność |
|---|---|---|
| Odczyt | Pojedyncze sesje, wpisy pamięci, historia, zlecenia i zadania. | Z okna aplikacji, przez interfejs rdzenia. |
| Edycja | Pojedyncze sesje, wpisy pamięci, historia, zlecenia i zadania. | Z okna aplikacji, przez interfejs rdzenia. |
| Usuwanie | Pojedyncze sesje, wpisy pamięci, historia, zlecenia i zadania. | Z okna aplikacji, przez interfejs rdzenia. |

### 10.1. Schemat encji i relacji (skrócony)

Model danych pełny dzieli się na osiemnaście grup encji (rozdz. 3–20 dokumentu [Model danych](model-danych.md)). Poniższy schemat pokazuje wyłącznie relację nadrzędności między grupami, od konta po pojedynczy zapis sesji — bez pól ani kluczy, które niesie dokument źródłowy.

```
Konto (właściciel, rozdz. 3 Modelu danych)
  │
  ├── Urządzenia (rozdz. 3)
  │
  └── Środowiska (rozdz. 4)
        │
        └── Moduły (rozdz. 4)
              │
              └── Karty sesji (rozdz. 8)
                    │
                    ├── Sesje / procesy (rozdz. 8; rozdz. 7 niniejszego dokumentu)
                    │     │
                    │     ├── Kanały modeli przypisane (rozdz. 11 Modelu danych; rozdz. 9 niniejszego)
                    │     ├── Wiadomości obu kanałów operacyjnych (rozdz. 19 Modelu danych; rozdz. 5)
                    │     └── Zadania, kolejki, ponowienia (rozdz. 10 Modelu danych; rozdz. 5.3)
                    │
                    ├── Pamięć wielopoziomowa (rozdz. 9 Modelu danych)
                    ├── Profile i punkty izolacji (rozdz. 15 Modelu danych; rozdz. 12 niniejszego)
                    └── Konfiguracja warstwowa (rozdz. 17 Modelu danych; rozdz. 13 niniejszego)

Byty poprzeczne, niezależne od gałęzi Środowisko → Moduł → Karta sesji:
  Projekty (rozdz. 7) · Agenci (rozdz. 12) · Role MultitaskingAI (rozdz. 13)
  Rozszerzenia (rozdz. 14) · Artefakty i biblioteka (rozdz. 16) · Komponenty własne (rozdz. 6)
```

Pełny schemat — z nazwami pól, typami, kluczami obcymi i ograniczeniami — niesie wyłącznie [Model danych](model-danych.md); niniejszy dokument nie powiela go i przywołuje wyłącznie relację nadrzędności między grupami.

Szczegółowy model danych opisuje dokument [Model danych](model-danych.md).

---

## 11. Komunikacja klient–serwer

Komunikacja przebiega jednym, stałym kanałem WebSocket, w formacie JSON, w trzech kategoriach komunikatów. Tym samym kanałem przenoszone są komunikaty obu kanałów komunikacji operacyjnej (rozdz. 5).

```
        polecenia   (klient → serwer)
KLIENT ─────────────────────────────────────►  SERWER
       ◄─────────────────────────────────────
        zdarzenia   (serwer → klient)
       ◄─────────────────────────────────────
        strumienie  (na żywo: odpowiedź modelu przekazywana w miarę generowania)

  Kanał: WebSocket — połączenie stałe i dwukierunkowe · format komunikatów: JSON
  Serwer inicjuje przekazanie zmian do wszystkich urządzeń.
```

| Kategoria komunikatu | Kierunek | Charakterystyka |
|---|---|---|
| Polecenia | Klient → serwer | Żądania działań inicjowane przez klienta, w tym polecenia Użytkownika i sterowanie przebiegiem pętli wykonawczej. |
| Zdarzenia | Serwer → klient | Powiadomienia inicjowane przez serwer, w tym zmiany stanu zleceń, zadań, kontroli jakości i ponowień. |
| Strumienie | Serwer → klient (na żywo) | Strumień odpowiedzi modelu przekazywany w miarę generowania. |

Pełny wykaz typów komunikatów opisuje dokument [Kontrakty komunikacji](kontrakty-komunikacji.md).

---

## 12. Izolacja i konfigurowalność zależności

Izolacja jest funkcją konfigurowalną, wybieraną przez Operatora. Aplikacja przewiduje pełny zbiór zakresów; domyślnie żaden nie jest aktywny. Zakresy stosuje się niezależnie i łącznie; zbiór aktywnych zakresów tworzy profil izolacji przypisywalny do sesji, roli lub projektu.

**Dwa rodzaje izolacji** schodzą się w jednym oknie konfiguracji punktów izolacji.

| Rodzaj izolacji | Obejmuje | Stan domyślny |
|---|---|---|
| Izolacja kontekstu | Historia, pamięć, kontekst. | Odrębny kontekst dla każdej karty sesji. |
| Izolacja techniczna | Osiem zakresów procesu sesji (poniżej). | Żaden zakres nie jest aktywny. |

**Zakresy izolacji technicznej:**

| Zakres | Opis |
|---|---|
| Katalog roboczy sesji | Ograniczenie pracy do wskazanego katalogu. |
| Środowisko procesu | Zakres zmiennych środowiskowych procesu. |
| Katalog danych i konfiguracji modelu | Odrębny katalog konfiguracji narzędzia modelu. |
| Dostęp sieciowy | Dopuszczenie lub odcięcie sieci. |
| Zakres odczytu i zapisu plików | Granice dostępu do plików. |
| Konto i token per sesja | Odrębne dane dostępowe dla sesji. |
| Model procesu | Proces odrębny lub współdzielony. |
| Serwer wykonania | Wykonanie lokalne, zdalne lub izolowane. |

**Poziomy zasięgu reguł izolacji.** Komenda `isolation.scope.list` (obszar `isolation`, `budowa/shared/contract.json`) zwraca „osiem poziomów zasięgu izolacji w kolejności rozstrzygania” — struktura wyniku `IsolationScopeLevel[]` typuje pole `scope` wyliczeniem `ConfigScope`. Spośród dziewięciu wartości tego wyliczenia poziom `application` nie uczestniczy w rozstrzyganiu izolacji — obowiązuje wyłącznie w konfiguracji aplikacji (rozdz. 13) — co daje osiem poziomów właściwych izolacji, z regułą pierwszeństwa poziomu najbardziej szczegółowego.

| Poziom zasięgu | Wartość `ConfigScope` | Przykład zastosowania | Pierwszeństwo |
|---|---|---|---|
| Globalny (domyślny) | `global` | Domyślna polityka izolacji obowiązująca w całej platformie. | Najniższe — warstwa bazowa |
| Środowisko | `environment` | Odrębna polityka dla wybranego środowiska. | rośnie ↓ |
| Moduł | `module` | Odrębna pamięć przypisana modułowi. | rośnie ↓ |
| Para modułów | `modulePair` | Współdzielenie wybranego mechanizmu między dwoma modułami. | rośnie ↓ |
| Projekt | `project` | Rozdzielenie kontekstu między projektami w module Workspace. | rośnie ↓ |
| Karta sesji | `session` | Jednorazowe współdzielenie historii między otwartymi kartami. | rośnie ↓ |
| Rola (środowisko MultitaskingAI) | `role` | Profil izolacji przypisany roli wykonawczej. | rośnie ↓ |
| Okno komunikacji | `window` | Ustawienie izolacji właściwe pojedynczemu oknu Chat Window albo Execution Loop Window (rozdz. 5); parametr opcjonalny `windowId` komendy `isolation.scope.list` wylicza byty tego poziomu. | Najwyższe — poziom najwęższy |

```
Zbiór aktywnych zakresów  ──►  PROFIL IZOLACJI  ──►  przypisanie do: sesji | roli | okna | projektu

Reguła pierwszeństwa — poziom najbardziej szczegółowy wygrywa:
   globalny ◄ środowisko ◄ moduł ◄ para modułów ◄ projekt ◄ karta sesji ◄ rola ◄ okno
   (brak ustawienia na poziomie = dziedziczenie z poziomu szerszego, aż do globalnego)
```

**Szablon profilu izolacji (pola konfiguracji):**

```
profil_izolacji:
  nazwa:            <nazwa profilu>
  poziom_zasiegu:   globalny | srodowisko | modul | para_modulow | projekt | karta_sesji | rola | okno
  izolacja_kontekstu:
    historia:       wspoldzielona | odrebna
    pamiec:         wspoldzielona | odrebna
    kontekst:       wspoldzielony | odrebny
  izolacja_techniczna:               # osiem zakresów — każdy: wlaczony | wylaczony
    katalog_roboczy_sesji:                wylaczony
    srodowisko_procesu:                   wylaczony
    katalog_danych_i_konfiguracji_modelu: wylaczony
    dostep_sieciowy:                      wylaczony
    zakres_odczytu_i_zapisu_plikow:       wylaczony
    konto_i_token_per_sesja:              wylaczony
    model_procesu:                        wylaczony
    serwer_wykonania:                     wylaczony
  przypisanie:      sesja | rola | projekt
  wartosc_domyslna: dziedziczenie z poziomu szerszego; brak aktywnej izolacji technicznej
```

Stanem wyjściowym jest odrębny kontekst dla każdej karty sesji oraz brak aktywnej izolacji technicznej — okno nigdy nie wymusza izolacji, wyłącznie ją udostępnia.

### 12.1. Diagram sekwencji — przypisanie profilu izolacji do sesji

Poniższa sekwencja pokazuje przebieg zapisu profilu izolacji w oknie konfiguracji punktów izolacji i jego przypisania do sesji, aż do chwili, w której proces sesji (rozdz. 7) zaczyna działać pod nową polityką. Nazwy komend pochodzą z obszaru `isolation` kontraktu komunikacji (rozdz. 18).

```
Operator            Klient (okno konfiguracji)           Rdzeń                    Proces sesji
   │                       │                                │                          │
   │  otwarcie okna        │                                │                          │
   │  punktów izolacji     │                                │                          │
   ├──────────────────────►│                                │                          │
   │                       │  isolation.scope.list          │                          │
   │                       ├───────────────────────────────►│                          │
   │                       │  scopes: IsolationScopeLevel[] │                          │
   │                       │◄───────────────────────────────┤                          │
   │  ustawienie           │                                │                          │
   │  przełączników        │                                │                          │
   ├──────────────────────►│  isolation.technical.set       │                          │
   │                       ├───────────────────────────────►│                          │
   │  nazwanie i zapis     │                                │                          │
   │  profilu              │                                │                          │
   ├──────────────────────►│  isolation.profile.save        │                          │
   │                       ├───────────────────────────────►│                          │
   │                       │  profileId                     │                          │
   │                       │◄───────────────────────────────┤                          │
   │  przypisanie profilu  │                                │                          │
   │  do sesji             │                                │                          │
   ├──────────────────────►│  isolation.profile.assign      │                          │
   │                       │  (scope: session)              │                          │
   │                       ├───────────────────────────────►│                          │
   │                       │                                │  rozstrzygnięcie polityki│
   │                       │                                │  wg IsolationScopeLevel  │
   │                       │                                ├─────────────────────────►│
   │                       │                                │                          │  proces działa
   │                       │                                │                          │  pod nową polityką
   │                       │  zdarzenie potwierdzające      │                          │
   │                       │◄───────────────────────────────┤                          │
   │  potwierdzenie        │                                │                          │
   │◄──────────────────────┤                                │                          │
```

Legenda: `isolation.scope.list`, `isolation.technical.set`, `isolation.profile.save`, `isolation.profile.assign` — komendy obszaru `isolation`; `isolation.policy.preview` (nieujęta w sekwencji) zwraca politykę obowiązującą po rozstrzygnięciu ośmiu poziomów zasięgu bez trwałego zapisu i służy do podglądu skutku przed zatwierdzeniem.

### 12.2. Wykaz komend obszaru `isolation`

Obszar `isolation` kontraktu komunikacji liczy dwanaście komend. Poniższy wykaz jest pełny dla tego obszaru.

| Komenda | Przeznaczenie |
|---|---|
| `isolation.scope.list` | Zwraca osiem poziomów zasięgu izolacji w kolejności rozstrzygania (rozdz. 12) |
| `isolation.context.get` | Odczytuje przełączniki izolacji kontekstu zapisane na wskazanym poziomie i warstwie |
| `isolation.context.set` | Zapisuje przełączniki izolacji kontekstu na wskazanym poziomie i warstwie |
| `isolation.technical.get` | Odczytuje przełączniki ośmiu zakresów technicznych zapisane na wskazanym poziomie i warstwie |
| `isolation.technical.set` | Zapisuje przełączniki ośmiu zakresów technicznych na wskazanym poziomie i warstwie |
| `isolation.profile.save` | Zapisuje profil izolacji; puste `profileId` zakłada nowy, podane zmienia istniejący |
| `isolation.profile.list` | Zwraca profile izolacji |
| `isolation.profile.load` | Wczytuje profil izolacji do panelu bez przypisywania go do poziomu |
| `isolation.profile.assign` | Przypisuje profil izolacji do wskazanego poziomu zasięgu i warstwy |
| `isolation.profile.delete` | Usuwa profil izolacji |
| `isolation.layer.set` | Przełącza warstwę izolacji między domyślną platformy a warstwą karty sesji |
| `isolation.policy.preview` | Zwraca politykę izolacji obowiązującą po rozstrzygnięciu ośmiu poziomów zasięgu, bez zapisu |

Pełny opis mechanizmu, macierz izolacji i szablony konfiguracji zawiera [Izolacja i zależności](izolacja-i-zaleznosci.md).

---

## 13. Konfiguracja

Konfiguracja ma strukturę warstwową. Wartości nakładają się od poziomu najszerszego do najwęższego; sesja dziedziczy wartości domyślne i nakłada własne zmiany, nie modyfikując poziomów wyższych.

```
Globalny  ──►  Środowisko  ──►  Projekt  ──►  Sesja
(najszerszy)                                  (najwęższy)

Nakładanie wartości: od poziomu najszerszego do najwęższego.
Sesja nakłada własne zmiany, nie modyfikując poziomów wyższych.
```

| Warstwa konfiguracji | Zakres | Nakładanie |
|---|---|---|
| Globalny | Cała platforma. | Warstwa bazowa (najszersza). |
| Środowisko | Wybrane środowisko. | Nakłada się na globalny. |
| Projekt | Wybrany projekt. | Nakłada się na środowisko. |
| Sesja | Bieżąca sesja. | Nakłada się na projekt; nie modyfikuje warstw wyższych. |

**Zasady okna konfiguracji:**

| Cecha | Zasada |
|---|---|
| Objaśnienie kontekstowe | Każde ustawienie posiada objaśnienie opisujące jego działanie i wpływ na aplikację; jest częścią definicji ustawienia i prezentowane przy elemencie w oknie konfiguracji. |
| Kompletność dostępu | Nie istnieje ustawienie niedostępne z okna konfiguracji. |
| Zasięg warstw widoczności | Ustawienie zakresu warstw widoczności przypisanego roli użytkownika (rozdz. 6) jest częścią konfiguracji. |

**Wyliczenie `ConfigScope`.** Cztery warstwy ogólne opisane wyżej (Globalny, Środowisko, Projekt, Sesja) porządkują konfigurację na poziomie produktowym. Kontrakt komunikacji rozstrzyga zasięg zapisu i odczytu wielu obszarów — nie tylko ustawień okna konfiguracji, lecz również izolacji (rozdz. 12), przypisania komponentu własnego (Rozszerzenia, rozdz. 4, obszar `component`) i innych mechanizmów — wyliczeniem `ConfigScope`, liczącym dziewięć wartości, drobniejszym niż cztery warstwy ogólne:

| # | Wartość | Zasięg | Pierwszeństwo |
|---|---|---|---|
| 1 | `application` | Aplikacja jako całość — nastawy samego programu, nie treści w nim prowadzonej: wymóg logowania, adres i postać nasłuchu. Zapisu na tym poziomie nie ma czym zawęzić — bytu poziomu nie ma (`scopeId` puste) | Najszersze — przegrywa z każdym węższym zapisem |
| 2 | `global` | Poziom globalny platformy | rośnie ↓ |
| 3 | `environment` | Środowisko | rośnie ↓ |
| 4 | `module` | Moduł | rośnie ↓ |
| 5 | `modulePair` | Para modułów (relacja) | rośnie ↓ |
| 6 | `project` | Projekt | rośnie ↓ |
| 7 | `session` | Karta sesji | rośnie ↓ |
| 8 | `role` | Rola (MultitaskingAI) | rośnie ↓ |
| 9 | `window` | Okno komunikacji | Najwęższe — wygrywa rozstrzyganie |

Nie każdy obszar kontraktu wykorzystuje wszystkie dziewięć wartości `ConfigScope` — obszar `isolation` udostępnia osiem z nich w rozdziale 12 niniejszego dokumentu, zgodnie z wykazem ośmiu poziomów zasięgu w [Izolacji i zależnościach](izolacja-i-zaleznosci.md) rozdz. 6; wartość `application` jest właściwa wyłącznie nastawom programu jako całości — wymogowi logowania oraz adresowi i postaci nasłuchu — nie treści prowadzonej w platformie. Deweloper wprowadzający nowy obszar korzystający z zasięgu konfigurowalnego wybiera podzbiór `ConfigScope` właściwy temu obszarowi i dokumentuje go w opracowaniu tego obszaru, nie definiuje własnego wyliczenia równoległego — zgodnie z regułą jednego źródła prawdy ([Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md#74-reguła-jednego-źródła-prawdy)).

---

## 14. Rozszerzenia

Rozszerzenia realizują jednolity kontrakt integracji i dzielą się na dwa źródła. Rdzeń nie rozróżnia źródła rozszerzenia.

```
Rodzaje rozszerzeń:  wtyczki · umiejętności · konektory · serwery MCP

   Danaco Plugin (wbudowane) ──┐
                               ├──►  JEDNOLITY KONTRAKT INTEGRACJI  ──►  RDZEŃ
   Personal (Operatora)      ──┘          (rdzeń nie rozróżnia źródła)
```

| Źródło | Opis | Przechowywanie |
|---|---|---|
| Danaco Plugin | Zestaw wbudowany, dostarczany z pakietem aplikacyjnym serwera. | Pakiet aplikacyjny serwera. |
| Personal | Rozszerzenia instalowane przez Operatora z urządzenia; przesyłane na serwer. | Katalog użytkownika. |

Rozszerzenia wykonują się w procesie rdzenia albo w piaskownicy nadzorowanej przez rdzeń, zgodnie z granicami wykonania (rozdz. 8).

### 14.1. Obszar kontraktu `extension`

Obszar `extension` liczy trzydzieści siedem komend, pogrupowanych rodzinami wokół cyklu życia rozszerzenia. Poniższa tabela zestawia liczebność rodzin; pełny wykaz trzydziestu siedmiu komend wraz z polami żądania i wyniku niesie [Rozszerzenia](rozszerzenia.md).

| Rodzina komend | Liczba komend | Zakres |
|---|---|---|
| `extension.*` (podstawa) | 6 | katalog, instalacja, konfiguracja, włączenie, odinstalowanie, wyszukiwanie |
| `extension.detail.*` | 1 | pełna metryka pozycji katalogu |
| `extension.collection.*` | 3 | kolekcje kuratorskie — zestawy instalowane grupowo |
| `extension.registry.*` | 1 | rejestr prywatny organizacji |
| `extension.update.*`, `extension.version.*` | 3 | sprawdzenie, przypięcie i cofnięcie wersji |
| `extension.bundle.*`, `extension.history.*`, `extension.admin.*` | 3 | zestawy definicji, chronologia zdarzeń, czynności zbiorcze |
| `extension.tool.*`, `extension.protocol.*`, `extension.sandbox.*` | 4 | narzędzia serwera MCP, dziennik protokołu, uruchomienie w piaskownicy izolacyjnej |
| `extension.definition.*`, `extension.transport.*`, `extension.credential.*` | 3 | budowa integracji z opisu API, transport, powiązanie z poświadczeniem |
| `extension.webhook.*`, `extension.mapping.*`, `extension.usage.*`, `extension.health.*` | 4 | webhooki, odwzorowanie pól, liczniki użycia, kondycja integracji |
| `extension.secret.*`, `extension.permission.*` | 4 | referencje sekretów i uprawnienia deklarowane w manifeście |
| `extension.signature.*`, `extension.manifest.*`, `extension.audit.*` | 3 | podpis cyfrowy, przegląd manifestu, zestawienie audytowe |

### 14.2. Diagram sekwencji — instalacja rozszerzenia źródła Personal

```
Operator              Klient                        Rdzeń                        Katalog użytkownika
   │                     │                             │                            │
   │  wybór paczki       │                             │                            │
   │  na urządzeniu      │                             │                            │
   ├────────────────────►│                             │                            │
   │                     │  extension.package.upload   │                            │
   │                     ├────────────────────────────►│                            │
   │                     │                             │  zapis paczki              │
   │                     │                             ├───────────────────────────►│
   │                     │  extension.signature.verify │                            │
   │                     ├────────────────────────────►│                            │
   │                     │  poziom zaufania wydawcy    │                            │
   │                     │◄────────────────────────────┤                            │
   │                     │  extension.install          │                            │
   │                     ├────────────────────────────►│                            │
   │                     │                             │  rejestracja pozycji       │
   │                     │                             │  jednolitym kontraktem     │
   │                     │  potwierdzenie              │                            │
   │                     │◄────────────────────────────┤                            │
   │  rozszerzenie       │                             │                            │
   │  dostępne w katalogu│                             │                            │
   │◄────────────────────┤                             │                            │
```

Legenda: `extension.package.upload` przyjmuje treść pliku przesłaną z urządzenia Operatora; `extension.signature.verify` sprawdza podpis cyfrowy i sumę kontrolną przed instalacją. Rdzeń obsługuje rozszerzenie źródła Personal tym samym jednolitym kontraktem integracji co rozszerzenia wbudowane Danaco Plugin (rozdz. A.7).

---

## 15. Uwierzytelnianie

Uwierzytelnianie jest jedynym mechanizmem kontroli dostępu w systemie. Kontrola dostępu jest włączona, a rdzeń działa na serwerze wirtualnym.

| Metoda | Dane | Uwagi |
|---|---|---|
| Rejestracja (pierwsze uruchomienie) | Login, adres e-mail uwierzytelniający, hasło. | Potwierdzenie przez e-mail. |
| Logowanie | Login i hasło. | — |
| Metody dodatkowe (okno konfiguracji) | PIN, Windows Hello, e-mail uwierzytelniający. | Konfigurowane w oknie konfiguracji. |
| Odzyskiwanie konta | Adres e-mail. | Ustanowienie nowego hasła. |

**Token i połączenie:**

```
Uwierzytelnienie (login + hasło)  ──►  TOKEN  ──►  nawiązanie połączenia WebSocket

  Dane dostępowe przechowywane są poza bazą.
  Repozytorium zawiera wyłącznie plik przykładowy z nazwami pól.
```

Adres połączenia klienta pochodzi z konfiguracji.

### 15.1. Diagram sekwencji — logowanie i nawiązanie połączenia

```
Operator              Klient                    Rdzeń (bramka)
   │                     │                         │
   │  login + hasło      │                         │
   ├────────────────────►│                         │
   │                     │  auth.login             │
   │                     ├────────────────────────►│
   │                     │                         │  weryfikacja skrótu hasła
   │                     │                         │  (dane dostępowe poza bazą)
   │                     │  token sesji bramki     │
   │                     │◄────────────────────────┤
   │                     │  connection.hello       │
   │                     │  (token: token sesji)   │
   │                     ├────────────────────────►│
   │                     │  authenticated: true    │
   │                     │  gatewayConfigured      │
   │                     │◄────────────────────────┤
   │  połączenie         │                         │
   │  WebSocket czynne   │                         │
   │◄────────────────────┤                         │
```

Legenda: `auth.login` (obszar `auth`) przyjmuje login i hasło albo aktywną metodę dodatkową — PIN, Windows Hello, adres e-mail uwierzytelniający; `connection.hello` (obszar `connection`) wiąże nawiązywane połączenie WebSocket z sesją bramki przez pole `token` i zwraca `authenticated`, informujące klienta, czy połączenie jest związane z ważną sesją. Pole `loginRequired` wyniku `connection.hello`, odczytane z nastawy zasięgu `application` (klucz `gateway.requireLogin`), kieruje klienta do okna logowania albo pomija je.

### 15.2. Obszar kontraktu `auth`

Obszar `auth` liczy szesnaście komend obejmujących rejestrację, logowanie, metody dodatkowe, odzyskiwanie konta i zarządzanie profilem. Pełny wykaz wraz z polami żądania i wyniku niesie [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md), Załącznik z pełnym wykazem komend kontraktu uwierzytelniania i dostępu; macierz metod dodatkowych — Załącznik B tamtego opracowania, a ważność tokenu i przebieg jego odnowienia — rozdz. 8 tamtego opracowania.

---

## 16. Synchronizacja

Synchronizacja odbywa się na żywo, kanałem WebSocket. Zmiana na jednym urządzeniu przekazywana jest przez serwer do pozostałych.

```
Urządzenie A ──(zmiana)──►  SERWER  ──(przekaz na żywo, WebSocket)──►  Urządzenie B
                              │                                        Urządzenie C
                              └──────────────────────────────────────►  ...

  Urządzenia prezentują wyłącznie obraz pochodzący z serwera.
  Telefon i komputer widzą ten sam stan w tej samej chwili.
```

| Cecha synchronizacji | Wartość |
|---|---|
| Kanał | WebSocket. |
| Tryb | Na żywo. |
| Kierunek | Zmiana na jednym urządzeniu → serwer → pozostałe urządzenia. |
| Zakres | Sesje, rozmowy, historia, pamięć, konfiguracja, stan procesów, zlecenia, zadania i komunikaty obu kanałów operacyjnych. |
| Źródło obrazu | Wyłącznie serwer (urządzenia nie przechowują stanu trwałego). |
| Spójność | Telefon i komputer widzą ten sam stan w tej samej chwili. |

---

## 17. Widok wdrożeniowy

```
        URZĄDZENIA OPERATORA                          SERWER WIRTUALNY
   ┌───────────────────────┐                    ┌────────────────────────────────┐
   │ Komputer              │                    │ RDZEŃ (Go)                     │
   │   okno Tauri (Rust)   │◄──── WebSocket ───►│   orkiestracja sesji           │
   ├───────────────────────┤                    │   Koordynator · kolejka zadań  │
   │ Telefon               │◄──── WebSocket ───►│   kontrola jakości · ponowienia│
   ├───────────────────────┤                    │   procesy sesji (izolowane)    │
   │ Tablet                │◄──── WebSocket ───►│   adaptery modeli              │
   └───────────────────────┘                    │   polityki izolacji            │
     cienki klient:                             │   serwer rozszerzeń            │
     natywne okno,                              ├────────────────────────────────┤
     brak stanu trwałego,                       │ DANE: SQLite + pliki treści    │
     brak wymogu GPU                            │        (jedno źródło prawdy)   │
                                                └────────────────────────────────┘

   Kanały modeli:  API (klucz, HTTP) · CLI (token) · SSH (modele pomocnicze) · HTTP (modele przeglądarkowe, zdalne punkty wykonania GPU)
```

### 17.1. Drzewo katalogów wdrożenia

Poniższe drzewo odzwierciedla podział katalogu `budowa/` na komponent serwerowy (rdzeń Go), komponent kliencki (powłoka TypeScript), warstwę współdzieloną kontraktu oraz komponent natywnego okna.

```
budowa/
├── server/                        RDZEŃ — serwer wykonawczy w języku Go
│   ├── cmd/
│   │   ├── danaco-console/        punkt startowy procesu serwera
│   │   └── danaco-narzedzia/      narzędzia pomocnicze wywoływane z wiersza poleceń
│   └── internal/
│       ├── core/                  orkiestracja sesji, Koordynator, kolejka zadań
│       ├── session/               model procesów sesji (rozdz. 7)
│       ├── dane/                  warstwa danych, dostęp do SQLite (rozdz. 10)
│       ├── store/                 przechowywanie i odwołania do plików treści
│       ├── repozytorium/          repozytoria danych encji modelu
│       ├── protocol/              obsługa ramek WebSocket i kontraktu (rozdz. 11)
│       ├── transport/             warstwa transportowa połączeń
│       ├── models/                integracja modeli — rejestr kanałów (rozdz. 9)
│       ├── zdalne/                adaptery kanałów zdalnych (SSH, HTTP)
│       ├── zewnetrzne/            integracje z systemami zewnętrznymi
│       ├── podagenci/             Subagent Network (środowisko MultitaskingAI)
│       ├── injection/             mechanizm rozszerzeń i wstrzykiwania zależności
│       ├── konfig/, konfiguracja/ warstwa konfiguracji (rozdz. 13)
│       ├── mowa/                  rozpoznawanie mowy, silnik lokalny CPU-only
│       ├── tokenizator/           liczenie tokenów wywołań modelu
│       ├── wiedza/                indeks znaczeniowy i wyszukiwanie
│       ├── poczta/                integracja poczty elektronicznej
│       └── narzedzia/             narzędzia współdzielone rdzenia
├── client/                        KLIENT — powłoka w języku TypeScript
│   └── src/
│       ├── main.ts                punkt startowy aplikacji klienckiej
│       ├── polaczenie/            nawiązanie i utrzymanie kanału WebSocket
│       ├── protokol/              kodowanie i dekodowanie komunikatów kontraktu
│       ├── powloka/               struktura okna, boczna nawigacja modułów
│       ├── okno-komunikacji/      Chat Window i Execution Loop Window (rozdz. 5)
│       ├── srodowiska/            warstwa środowisk (rozdz. 4 Koncepcji platformy)
│       ├── moduly/                warstwa modułów platformy
│       ├── strona-glowna/         strona główna i jej trzy strefy
│       ├── rozmowa/               prezentacja wiadomości i strumienia odpowiedzi
│       ├── uwierzytelnienie/      okna logowania i rejestracji (rozdz. 15)
│       ├── konfiguracja/          okno konfiguracji (rozdz. 13)
│       ├── ustawienia/            katalog ustawień sterowany danymi
│       ├── punkty-izolacji/       okno konfiguracji punktów izolacji (rozdz. 12)
│       ├── dostepy/               punkty dostępu i nadania dostępu
│       ├── komponenty/            komponenty własne wielokrotnego użycia
│       ├── okna-pomocnicze/       panele wysuwane i okna nakładkowe
│       ├── okna-rownolegle/       okna otwierane równolegle do obszaru roboczego
│       ├── sterowanie/            skróty klawiszowe, wyszukiwarka funkcji
│       ├── motyw/                 żetony i tryb jasny / ciemny
│       ├── ikony/                 zestaw ikon interfejsu
│       ├── powiadomienia/         centrum powiadomień
│       ├── mobile/                warstwa mobilna, mobilne centrum dowodzenia
│       ├── aod/                   nakładka Always On Display
│       ├── aktualizacja/          mechanizm aktualizacji klienta
│       └── ladowanie/             ekran startowy i stan ładowania
├── desktop/
│   └── src-tauri/                 okno natywne (Rust, warstwa Tauri — rozdz. 3)
└── shared/                        WARSTWA WSPÓŁDZIELONA — kontrakt komunikacji
    ├── contract.json              źródło normatywne kontraktu (rozdz. 18)
    ├── contract.go                typy kontraktu wygenerowane dla rdzenia
    ├── contract.ts                typy kontraktu wygenerowane dla klienta
    └── kontrasty-progi.json       progi kontrastu dostępności
```

Katalogi `budowa/poczta/`, `budowa/pomocniki/`, `budowa/scripts/`, `budowa/packaging/`, `budowa/witryna/` i `budowa/wydania/` należą do warstwy narzędzi budowy, dystrybucji i witryny produktowej — nie są częścią warstw architektury opisanych w rozdz. 4 i nie mają odpowiednika w modelu wdrożenia (rozdz. 2).

---

## 18. Wykaz obszarów kontraktu

Kontrakt komunikacji (`budowa/shared/contract.json`) dzieli komendy i zdarzenia na
**68 obszarów**. Każdy obszar grupuje rodzinę komend wokół jednego bytu albo jednego
modułu platformy; przypisanie komendy do obszaru jest stałe i nie zmienia się między
wydaniami. Poniższy wykaz jest źródłem prawdy dla pola „obszar” przywoływanego przez
opracowania modułowe, środowiskowe i interfejsu — żadne opracowanie nie definiuje
własnej listy obszarów. Kolumna „Komend” podaje liczbę komend policzoną wprost z pola
`typ` każdej pozycji tablicy `komendy` kontraktu (1119 komend łącznie); suma kolumny
odtwarza tę liczbę dokładnie.

| Obszar | Komend | Przeznaczenie |
|---|---|---|
| `connection` | 1 | Nawiązanie połączenia |
| `home` | 1 | Strona główna — jedyna droga wejścia do środowiska |
| `environment` | 2 | Środowiska platformy |
| `module` | 1 | Moduły platformy |
| `workspace` | 59 | Przestrzeń robocza karty sesji |
| `session` | 19 | Sesje |
| `window` | 7 | Okna komunikacji |
| `message` | 3 | Wiadomości |
| `config` | 9 | Konfiguracja |
| `settings` | 2 | Katalog ustawień sterowany danymi |
| `access` | 9 | Punkty dostępu i nadania dostępu okna rozmowy |
| `account` | 6 | Konta modeli i konta programu code CLI |
| `identity` | 5 | Tożsamość modelu — kategorie zasad i ich treść |
| `channel` | 6 | Kanały modeli |
| `queue` | 16 | Kolejki |
| `context` | 2 | Przekazanie kontekstu |
| `action` | 1 | Katalog akcji sterowany danymi |
| `stream` | 0 | Strumień fragmentów (typ komunikatu bez odrębnych komend — rozdz. 11) |
| `progress` | 0 | Telemetria postępu (typ komunikatu bez odrębnych komend — rozdz. 11) |
| `studio` | 179 | Moduł Studio — redakcja dokumentu |
| `automation` | 32 | Moduł Automations — automatyki, harmonogramy i przebiegi |
| `browser` | 47 | Moduł Browser — przeglądanie z podglądem współdzielonym |
| `research` | 76 | Moduł Research — źródła, ustalenia i raport |
| `library` | 48 | Moduł Library — repozytorium plików i wiedzy |
| `translate` | 64 | Moduł Translate — tłumaczenia równoległe i glosariusz |
| `roundtable` | 46 | Moduł Roundtable — debata wielu modeli |
| `design` | 106 | Moduł Design — generowanie i kompozycja zasobów wizualnych |
| `assistant` | 4 | Moduł Assistant — polecenia głosowe i zlecenia wieloetapowe |
| `terminal` | 27 | Moduł Terminal — powłoki i procesy |
| `developer` | 50 | Moduł Developer — kod, repozytorium i budowanie |
| `diagnostics` | 9 | Moduł Diagnostics — dziennik zdarzeń, błędy i rekomendacje |
| `apps` | 41 | Moduł Apps — architektura, warsztaty i wdrożenia |
| `agent` | 30 | Moduł Agents — eksperci, ich model, umiejętności, konektory i uprawnienia |
| `memory` | 13 | Pamięć projektu — Context Memory |
| `orchestration` | 8 | Orkiestracja — zależności między krokami układu |
| `isolation` | 12 | Punkty izolacji — zasięg, przełączniki i profile izolacji |
| `component` | 5 | Komponenty własne strony głównej — warstwa jednorodna nad magazynami modułowymi |
| `extension` | 37 | Katalog rozszerzeń — konektory, wtyczki, serwery MCP i umiejętności |
| `role` | 4 | Role pętli MultitaskingAI i ich wcielenia |
| `subagent` | 4 | Podagenci okna wykonawcy — panel Subagent Network |
| `advisor` | 1 | Konsultacja u doradcy — model pyta inny model o radę w trakcie tury |
| `monitor` | 2 | Monitor procesów — telemetria postępu poza modułem Automations |
| `schedule` | 6 | Harmonogramy — odczyt cykliczności ustalonej automatykom |
| `model` | 1 | Wybór kanału modelu dla okna albo karty sesji |
| `mobile` | 3 | Warstwa mobilna — mobilne centrum dowodzenia |
| `aod` | 12 | Always On Display — nakładka podglądu, rozmowy i obserwacji procesu |
| `auth` | 16 | Uwierzytelnianie Operatora — adres e-mail i hasło, PIN albo Windows Hello |
| `speech` | 8 | Rozpoznanie mowy — silnik lokalny, wyłącznie na procesorze |
| `tools` | 3 | Wykaz pozycji po ukośniku — narzędzia do dołożenia i komendy akcji |
| `archive` | 2 | Archiwa — spakowanie zasobów i wydobycie zawartości paczki |
| `device` | 8 | Urządzenia konta — wykaz maszyn z dostępem i unieważnienie dostępu |
| `document` | 2 | Dokumenty — zamiana formatu i wydobycie tekstu z zasobu |
| `history` | 2 | Historia rozmów — wczytanie zapisu i jego usunięcie |
| `image` | 9 | Obrazy — rozpoznanie, przekształcenie, korekta i zamiana formatu zasobu graficznego |
| `knowledge` | 2 | Wiedza — indeks znaczeniowy i wyszukiwanie po znaczeniu |
| `mail` | 10 | Poczta — konta nadawcze, foldery, wiadomości i wysyłka |
| `media` | 2 | Nagrania dźwiękowe i filmowe — rozpoznanie zawartości i zamiana formatu |
| `panel` | 2 | Sekcje paneli — odczyt i zapis układu sekcji okna |
| `retention` | 1 | Retencja — okres przechowywania zapisu rozmów i pamięci |
| `team` | 5 | Zespoły wykonawców — skład zapisywany, odczytywany i powielany |
| `alert` | 5 | Obszar alert — rodzina komend wniesiona wraz z modułami |
| `clipboard` | 4 | Obszar clipboard — rodzina komend wniesiona wraz z modułami |
| `health` | 6 | Obszar health — rodzina komend wniesiona wraz z modułami |
| `launcher` | 2 | Obszar launcher — rodzina komend wniesiona wraz z modułami |
| `provenance` | 5 | Obszar provenance — rodzina komend wniesiona wraz z modułami |
| `snippet` | 3 | Obszar snippet — rodzina komend wniesiona wraz z modułami |
| `usage` | 2 | Obszar usage — rodzina komend wniesiona wraz z modułami |
| `notification` | 4 | Centrum powiadomień — jedyny mechanizm powiadamiania platformy; trwały rejestr zdarzeń wymagających wiedzy albo decyzji Operatora |

Obszary `alert`, `clipboard`, `health`, `launcher`, `provenance`, `snippet`, `usage`
zostały wniesione do kontraktu razem z modułami, które z nich korzystają, i nie mają
odrębnego opracowania architektonicznego — ich komendy są opisane przy module
konsumującym (rozdz. 4.2 standardu redakcyjnego: [Standard redakcyjny i językowy](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md#42-wykazy-obowiązkowe)).

Dziesięć obszarów o największej liczbie komend zestawia poniższa tabela — grupują się
one wokół modułów o najbardziej rozbudowanym kontrakcie roboczym.

| Miejsce | Obszar | Komend | Moduł albo mechanizm |
|---|---|---|---|
| 1 | `studio` | 179 | Moduł Studio |
| 2 | `design` | 106 | Moduł Design |
| 3 | `research` | 76 | Moduł Research |
| 4 | `translate` | 64 | Moduł Translate |
| 5 | `workspace` | 59 | Przestrzeń robocza karty sesji |
| 6 | `developer` | 50 | Moduł Developer |
| 7 | `library` | 48 | Moduł Library |
| 8 | `browser` | 47 | Moduł Browser |
| 9 | `roundtable` | 46 | Moduł Roundtable |
| 10 | `apps` | 41 | Moduł Apps |

---

## 19. Kody błędów kontraktu

Kontrakt komunikacji definiuje dziewięć kodów błędów, wspólnych dla wszystkich 68 obszarów i 1119 komend. Kod błędu jest jednym z dwóch możliwych wyników każdej komendy — obok wyniku pozytywnego — i nie jest definiowany odrębnie przez żaden obszar ani moduł: każdy kod błędu specyficzny dla domeny — kod `not_authenticated` zwracany przy nieudanej próbie logowania — jest wystąpieniem jednego z dziewięciu kodów ogólnych, nie nowym kodem. Pełny format obiektu błędu i jego zastosowanie w kopercie protokołu opisuje [Kontrakty komunikacji](kontrakty-komunikacji.md) rozdz. 19; niniejszy rozdział jest jedynym miejscem odniesienia dla samego wykazu dziewięciu kodów, zgodnie z zasadą nadrzędną tego dokumentu.

Pole `kodyBledow` w `budowa/shared/contract.json` zawiera dziś osiem pierwszych kodów wykazu; dziewiąty kod `command_not_understood` jest postacią docelową, a jego dopisanie w kodzie stanowi osobne zadanie.

| Kod | Ponawialny | Znaczenie |
|---|---|---|
| `validation_failed` | nie | Treść żądania niezgodna z kontraktem — pole wymagane nieobecne, typ pola niezgodny albo wartość poza wyliczeniem |
| `not_found` | nie | Wskazany byt nie istnieje — identyfikator nie odpowiada żadnemu rekordowi osiągalnemu dla wywołującego |
| `not_authenticated` | nie | Brak uwierzytelnienia — połączenie nie przedstawiło ważnego tokenu dostępu (Bezpieczeństwo i uwierzytelnianie, rozdz. 8) |
| `permission_denied` | nie | Uwierzytelniony, lecz bez uprawnienia do czynności — rozróżnienie od `not_authenticated` jest istotne: pierwsze pytanie brzmi „kim jesteś”, drugie „co wolno tobie” |
| `conflict` | nie | Stan bytu wyklucza czynność — powtórna rejestracja konta już istniejącego (Bezpieczeństwo i uwierzytelnianie, rozdz. 4) |
| `channel_unavailable` | tak | Kanał modelu niedostępny; błąd dotyczy wyłącznie bieżącego wywołania, nie unieważnia kanału trwale (Integracja modeli, rozdz. 6) |
| `rate_limited` | tak | Ograniczenie tempa po stronie kanału modelu lub rdzenia platformy |
| `internal_error` | tak | Błąd wewnętrzny rdzenia, niezwiązany z treścią żądania |
| `command_not_understood` | nie | Polecenie języka naturalnego nie zostało dopasowane do żadnej funkcji — kod odrębny od `validation_failed`, aby interfejs odpowiadał inaczej na niezrozumiane polecenie niż na niezgodność treści żądania z kontraktem (Kontrakty komunikacji, rozdz. 7.2) |

Kolumna „Ponawialny” przenosi pole `retryable` źródła: kod ponawialny uzasadnia automatyczne powtórzenie tego samego żądania przez klienta po odczekaniu; kod nieponawialny wymaga zmiany treści żądania albo działania Operatora przed ponowną próbą. Zasada zero blokad (rozdz. 12; Koncepcja platformy, rozdz. 14) nie unieważnia `permission_denied` ani `not_authenticated` — oba kody rozstrzygają dostęp do platformy i do jej zasobów, co jest odrębne od blokowania czynności operacyjnej wewnątrz zasobu już dostępnego.

---

## 20. Kryteria odbioru

Warstwa architektoniczna jest zgodna z niniejszym opracowaniem, gdy spełnione są
łącznie poniższe warunki. Każdy warunek jest sprawdzalny bez odwołania do intencji
autora — wyłącznie przez odczyt kodu i konfiguracji.

| Warunek | Sposób sprawdzenia |
|---|---|
| Rdzeń zbudowany modułem Go w wersji zadeklarowanej w `budowa/go.mod` | `go version` na środowisku budowy zwraca wersję nie niższą niż zadeklarowana |
| Klient zbudowany w wersji TypeScript zadeklarowanej w `budowa/client/package.json` | `tsc --version` w katalogu `budowa/client` |
| Każda komenda kontraktu przypisana do dokładnie jednego obszaru z rozdz. 18 | odczyt `budowa/shared/contract.json`, pole `obszar` każdej pozycji `komendy` |
| Kanał komunikacji operacyjnej działa wyłącznie przez pojedyncze połączenie WebSocket (rozdz. 5, 11) | brama nie otwiera drugiego równoległego kanału do tego samego klienta |
| Warstwy widoczności interfejsu (rozdz. 6) nie mają wyjątków niezgodnych z tabelą warstw | przegląd komponentów `.dn-*` względem przypisanej warstwy w `design/zasoby/css/komponenty.css` |
| Izolacja i konfigurowalność zależności (rozdz. 12) zachowuje zasadę „brak ustawienia = wartość domyślna” | test integracyjny nieustawionego poziomu zasięgu zwraca wartość odziedziczoną, nie błąd |
| Rozszerzenia (rozdz. 14) wchodzą do platformy wyłącznie przez jednolity kontrakt rozszerzenia | manifest rozszerzenia waliduje się względem schematu wspólnego dla Danaco Plugin i Personal |
| Synchronizacja (rozdz. 16) nie pozostawia stanu rozbieżnego między urządzeniami po scenariuszu Załącznika A.3 | powtórzenie scenariusza A.3 na dwóch urządzeniach kończy się identycznym stanem sesji |

---

## Załącznik A. Scenariusze eksploatacji

Scenariusze ilustrują współdziałanie mechanizmów opisanych w rozdziałach 1–17. Nie wprowadzają nowych funkcji — pokazują istniejące zachowania architektury jako ponumerowane przepływy.

### A.1. Przetrwanie sesji po rozłączeniu klienta

1. Operator uruchamia sesję; rdzeń tworzy odrębny proces i wpisuje go do rejestru aktywnych procesów pod identyfikatorem uruchomienia procesu (rozdz. 7).
2. Klient rozłącza się (zamknięcie okna, utrata sieci); proces sesji pracuje dalej po stronie serwera.
3. Operator łączy się ponownie z dowolnego urządzenia; klient zostaje podłączony do aktywnej sesji.
4. Serwer zostaje zrestartowany; stan sesji zapisany trwale (SQLite oraz pliki treści) nie ginie.
5. Sesja zostaje wznowiona z zapisu i wraca do stanu aktywnego.

### A.2. Przebieg zlecenia przez oba kanały komunikacji operacyjnej

1. Użytkownik wydaje polecenie w języku naturalnym w Chat Window; klient przesyła je jako polecenie kanałem WebSocket (rozdz. 11).
2. Koordynator przyjmuje zlecenie, dokonuje jego dekompozycji na zadania i umieszcza je w kolejce zadań (rozdz. 5).
3. Wykonawca realizuje przydzielone zadanie w procesie sesji; wynik trafia do kontroli jakości.
4. Ocena negatywna uruchamia ponowienie zadania; ocena pozytywna domyka zadanie.
5. Execution Loop Window prezentuje stan zlecenia, zadań, wyniki kontroli jakości i decyzje o ponowieniu; Chat Window prezentuje odpowiedzi i wyniki dla Użytkownika.
6. Użytkownik steruje przebiegiem — wstrzymuje, wznawia, przerywa albo koryguje zlecenie; Koordynator przekłada polecenie na zmianę stanu zadań.

### A.3. Synchronizacja stanu między urządzeniami

1. Operator wprowadza zmianę na komputerze — edycję wpisu pamięci.
2. Klient przesyła polecenie do serwera kanałem WebSocket (rozdz. 11).
3. Serwer zapisuje zmianę jako jedyne źródło prawdy i inicjuje przekazanie jej do pozostałych urządzeń.
4. Telefon i tablet otrzymują zdarzenie na żywo i prezentują ten sam stan w tej samej chwili (rozdz. 16).

### A.4. Dodanie nowego kanału modelu

1. Powstaje potrzeba podłączenia modelu przez kanał dotąd nieobsługiwany.
2. Deweloper dostarcza nowy adapter warstwy dostawcy modelu dla tego kanału (rozdz. 9).
3. Rdzeń, operujący na abstrakcji „wyślij zapytanie, odbierz odpowiedź”, pozostaje bez zmian.
4. W konfiguracji ustala się dane dostępowe oraz zasięg wyboru kanału — per sesja lub per rola.

### A.5. Wykonanie zadania obliczeniowego wymagającego akceleracji

1. Zadanie wymaga przetwarzania neuronowego o wysokim koszcie obliczeniowym (rozdz. 8).
2. Rdzeń kieruje je przez kanał modelu do zdalnego punktu wykonania wyposażonego w GPU; dane dostępowe pochodzą z konfiguracji (rozdz. 13).
3. Urządzenie Operatora nie uczestniczy w obliczeniach — cienki klient prezentuje wyłącznie strumień postępu i wynik.
4. Zadania lekkie rdzeń wykonuje samodzielnie na procesorze, bez zależności od akceleracji sprzętowej.

### A.6. Zastosowanie profilu izolacji do sesji

1. Operator otwiera okno konfiguracji punktów izolacji i wybiera poziom zasięgu (rozdz. 12).
2. W macierzy izolacji ustala izolację kontekstu (historia, pamięć, kontekst) oraz wybrane zakresy izolacji technicznej.
3. Zestaw ustawień zostaje zapisany jako nazwany profil izolacji i przypisany do sesji, roli lub projektu.
4. Reguła z poziomu bardziej szczegółowego ma pierwszeństwo; brak ustawienia oznacza dziedziczenie z poziomu szerszego, aż do globalnego domyślnego.

### A.7. Instalacja rozszerzenia z zasobów Operatora

1. Operator instaluje rozszerzenie (wtyczkę, umiejętność, konektor lub serwer MCP) z urządzenia, jako źródło Personal (rozdz. 14).
2. Rozszerzenie zostaje przesłane na serwer i zapisane w katalogu użytkownika.
3. Rdzeń obsługuje je przez ten sam, jednolity kontrakt integracji co rozszerzenia wbudowane Danaco Plugin, nie rozróżniając źródła.

### A.8. Dotarcie do funkcji eksperckiej

1. Interfejs w stanie spoczynku prezentuje wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (rozdz. 6).
2. Operator wywołuje funkcję warstwy 4 poleceniem języka naturalnego w Chat Window, skrótem klawiszowym albo przez wyszukiwarkę funkcji.
3. Funkcja uruchamia się jednym działaniem; po zakończeniu panel zamyka się i znika z przestrzeni roboczej.
4. Zakres warstw udostępnianych roli użytkownika wynika z konfiguracji (rozdz. 13).

### A.9. Ponowienie żądania po błędzie ponawialnym

1. Klient wysyła komendę obszaru `channel` żądającą odpowiedzi modelu przez kanał chwilowo przeciążony.
2. Serwer zwraca `status: "error"` z kodem `channel_unavailable` i polem `retryable: true` (rozdz. 19).
3. Klient odczekuje interwał rosnący i wysyła to samo żądanie ponownie, bez zmiany treści ładunku ani nowego `id` koperty odrębnego od reguły ponowień właściwej danemu kanałowi.
4. Powodzenie zwraca wynik pozytywny. Liczba prób nie ma granicy narzuconej z góry, więc przebieg nie ma stanu wyczerpania prób: ponowienia trwają do powodzenia albo do przerwania na polecenie Operatora, wydanego w Chat Window albo Execution Loop Window zależnie od kanału komunikacji operacyjnej, który zainicjował żądanie (rozdz. 5). Polecenie przerwania działa także w trakcie odczekiwania narastającego odstępu; po przerwaniu zlecenie pozostaje w stanie oczekującym, prezentowanym Operatorowi w tym samym oknie.
5. Kod nieponawialny — `validation_failed`, `not_found`, `not_authenticated`, `permission_denied`, `conflict` (rozdz. 19) — nie uruchamia tego przebiegu automatycznego — klient prezentuje treść błędu i wymaga zmiany żądania przed kolejną próbą wysłaną przez Operatora.

---

## Załącznik B. Słownik terminów architektonicznych

Słownik zestawia nazwy bytów architektonicznych wprowadzonych w niniejszym dokumencie, wraz z rozdziałem, w którym każdy byt jest definiowany po raz pierwszy. Opracowania modułowe, środowiskowe i interfejsu przywołują te nazwy dosłownie i nie definiują ich ponownie (rozdz. 7.4 [Standardu redakcyjnego i językowego](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md#74-reguła-jednego-źródła-prawdy)).

| Termin | Definicja skrócona | Rozdział źródłowy |
|---|---|---|
| Rdzeń | Serwer wykonawczy w języku Go; jedyne miejsce logiki biznesowej, pracy obliczeniowej i stanu trwałego | rozdz. 2, 4 |
| Klient | Cienki klient uruchamiający natywne okno na urządzeniu Operatora; wyłącznie prezentacja i interakcja | rozdz. 2, 4 |
| Warstwa dostawcy modelu | Jednolita abstrakcja rdzenia „wyślij zapytanie, odbierz odpowiedź”, niezależna od kanału integracji | rozdz. 9 |
| Adapter kanału | Komponent tłumaczący jednolite zapytanie rdzenia na protokół i uwierzytelnienie właściwe jednemu z czterech kanałów (API, CLI, SSH, HTTP) | rozdz. 9 |
| Kanał modelu | Skonfigurowane połączenie z konkretnym modelem, przypisywalne per sesja albo per rola | rozdz. 9 |
| Proces sesji | Odrębny, izolowany proces po stronie serwera, działający niezależnie od obecności klienta | rozdz. 7 |
| Koordynator | Komponent orkiestrujący rdzenia; dekomponuje zlecenie na zadania, przydziela je i ocenia wyniki | rozdz. 5 |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania przydzielone przez Koordynatora albo polecenia Użytkownika | rozdz. 5 |
| Chat Window | Okno kanału pierwszego — komunikacja Użytkownik ↔ Wykonawca | rozdz. 5.1 |
| Execution Loop Window | Okno kanału drugiego — komunikacja Koordynator ↔ Wykonawca, pętla wykonawcza | rozdz. 5.2 |
| Warstwa widoczności interfejsu | Jeden z czterech poziomów ujawniania funkcji — od zawsze widocznej po funkcje eksperckie | rozdz. 6 |
| Profil izolacji | Nazwany, zapisany zestaw aktywnych zakresów izolacji kontekstu i izolacji technicznej, przypisywalny do poziomu zasięgu | rozdz. 12 |
| Poziom zasięgu izolacji | Jeden z ośmiu poziomów rozstrzygania reguł izolacji — od globalnego po okno komunikacji | rozdz. 12 |
| Warstwa konfiguracji | Jedna z czterech warstw ogólnych porządkujących nakładanie wartości ustawień na poziomie produktowym: Globalny, Środowisko, Projekt, Sesja; odrębna od wyliczenia `ConfigScope` | rozdz. 13 |
| Rozszerzenie | Wtyczka, umiejętność, konektor albo serwer MCP, realizujące jednolity kontrakt integracji | rozdz. 14 |
| Sesja bramki | Token wydany przez `auth.login`, wiążący połączenie WebSocket z uwierzytelnionym kontem Operatora | rozdz. 15 |
| Widok wdrożeniowy | Rozmieszczenie komponentów architektury na urządzeniach Operatora i serwerze wirtualnym | rozdz. 17 |
| `ConfigScope` | Wyliczenie dziewięciu poziomów zasięgu kontraktu — od `application` (najszerszy) po `window` (najwęższy) — wykorzystywane wybiórczo przez obszary rozstrzygające zasięg zapisu, odrębne od czterech warstw ogólnych konfiguracji | rozdz. 13 |
| Kod błędu | Jeden z dziewięciu kodów wspólnych całemu kontraktowi, niosący pole `retryable` rozstrzygające, czy ponowienie żądania bez zmian ma sens | rozdz. 19 |

---

*Koniec dokumentu. Architektura techniczna — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
