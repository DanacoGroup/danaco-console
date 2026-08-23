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
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Przeznaczenie** | Źródło prawdy architektonicznej dla dewelopera |

Dokument opisuje architekturę techniczną produktu Danaco Console: model wdrożenia, stos technologiczny, warstwy systemu, warstwę komunikacji operacyjnej, warstwy widoczności interfejsu, model sesji i procesów, realizację zadań obliczeniowych, integrację modeli, warstwę danych, komunikację, izolację, konfigurację, rozszerzenia, uwierzytelnianie oraz synchronizację.

---

## Spis treści

1. Zasady architektoniczne
2. Model wdrożenia
3. Stos technologiczny
4. Warstwy systemu
5. Warstwa komunikacji operacyjnej
6. Warstwy widoczności interfejsu jako zasada architektoniczna
7. Model sesji i procesów
8. Realizacja zadań obliczeniowych i granice wykonania
9. Integracja modeli
10. Warstwa danych
11. Komunikacja klient–serwer
12. Izolacja i konfigurowalność zależności
13. Konfiguracja
14. Rozszerzenia
15. Uwierzytelnianie
16. Synchronizacja
17. Widok wdrożeniowy
- Załącznik A. Scenariusze eksploatacji

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
┌───────────────────────────────────────────────────────────────────┐
│ KLIENT          natywne okno · interfejs · interakcja Operatora     │
└─────────────────────────────────┬─────────────────────────────────┘
                                  │  WebSocket (JSON, dwukierunkowo, strumień na żywo)
┌─────────────────────────────────┴─────────────────────────────────┐
│ KOMUNIKACJA     kanał dwukierunkowy klient–serwer · strumień na żywo │
└─────────────────────────────────┬─────────────────────────────────┘
                                  │
┌─────────────────────────────────┴─────────────────────────────────┐
│ RDZEŃ           orkiestracja sesji · Koordynator · kolejka zadań ·  │
│                 procesy · integracja modeli · polityki izolacji     │
└─────────────────────────────────┬─────────────────────────────────┘
                                  │
┌─────────────────────────────────┴─────────────────────────────────┐
│ DANE            SQLite · pliki treści · jedno źródło prawdy         │
└───────────────────────────────────────────────────────────────────┘
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
| Mechanizm ponowień | Rdzeń | Realizuje ponowienie zadania z narastającym odstępem oraz przekazanie zadania nierozstrzygalnego do decyzji Użytkownika. |
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

**Mechanizmy ukrywania funkcjonalności:** menu progresywne (zbiór jednorodnych wyborów jako jeden element zwinięty, `Agent ▼`), panele wysuwane znikające całkowicie po zamknięciu, grupowanie logiczne akcji (`Operacje ▼`), znaczniki kontekstowe w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`, których kliknięcie otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

**Układ obszaru roboczego** jest wyłącznie pionowy, w podziale lewa–prawa; regulacji podlega wyłącznie szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Lewa, stała, pełna wysokość | Chat Window — kanał Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca (otwierana) | Execution Loop Window — kanał Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu (okna edycyjne, podglądu, monitory) |
| Kolejne kolumny boczne | Okna pomocnicze i panele — otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego |

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │                          │  boczne)
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
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
   │  klient          │   proces pracuje dalej     │  proces po stronie serwera   │
   │  podłączony      │ ◄───────────────────────   │  pracuje dalej               │
   └────────┬─────────┘   ponowne połączenie       └──────────────┬───────────────┘
            │           (podłączenie klienta do aktywnej sesji)   │
            │                                                     │
            │  zakończenie procesu / restart serwera              │
            ▼                                                     ▼
   ┌───────────────────────────────────────────────────────────────────────┐
   │  ZAPIS TRWAŁY   stan sesji zachowany (SQLite + pliki treści)           │
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

Modele udostępniane są przez jednolitą warstwę dostawcy modelu. Rdzeń operuje na abstrakcji „wyślij zapytanie, odbierz odpowiedź" i nie zależy od kanału; każdy kanał integracji realizuje osobny adapter tej warstwy.

```
RDZEŃ
  │  abstrakcja: „wyślij zapytanie, odbierz odpowiedź"  (rdzeń nie zależy od kanału)
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

---

## 10. Warstwa danych

Serwer stanowi jedyne źródło prawdy. Dane strukturalne i treści obszerne rozdzielone są między bazę i pliki, powiązane odwołaniami.

```
SERWER — jedyne źródło prawdy
  ┌────────────────────────────┐        ┌────────────────────────────┐
  │ SQLite (jeden plik)         │◄──────►│ Pliki treści               │
  │ dane strukturalne,          │ odwoła-│ kanwa, artefakty,          │
  │ transakcje                  │  nia   │ pliki biblioteki           │
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

Szczegółowy model danych opisuje dokument „Model danych".

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

Pełny wykaz typów komunikatów opisuje dokument „Kontrakty komunikacji".

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

**Poziomy zasięgu reguł izolacji.** Interfejs konfiguracji działa na siedmiu poziomach, od globalnego do roli, z regułą pierwszeństwa poziomu najbardziej szczegółowego.

| Poziom zasięgu | Przykład zastosowania | Pierwszeństwo |
|---|---|---|
| Globalny (domyślny) | Domyślna polityka izolacji obowiązująca w całej platformie. | Najniższe — warstwa bazowa |
| Środowisko | Odrębna polityka dla wybranego środowiska. | rośnie ↓ |
| Moduł | Odrębna pamięć przypisana modułowi. | rośnie ↓ |
| Para modułów | Współdzielenie wybranego mechanizmu między dwoma modułami. | rośnie ↓ |
| Projekt | Rozdzielenie kontekstu między projektami w module Workspace. | rośnie ↓ |
| Karta sesji | Jednorazowe współdzielenie historii między otwartymi kartami. | rośnie ↓ |
| Rola (środowisko MultitaskingAI) | Profil izolacji przypisany roli wykonawczej. | Najwyższe |

```
Zbiór aktywnych zakresów  ──►  PROFIL IZOLACJI  ──►  przypisanie do: sesji | roli | projektu

Reguła pierwszeństwa — poziom najbardziej szczegółowy wygrywa:
   globalny ◄ środowisko ◄ moduł ◄ para modułów ◄ projekt ◄ karta sesji ◄ rola
   (brak ustawienia na poziomie = dziedziczenie z poziomu szerszego, aż do globalnego)
```

**Szablon profilu izolacji (pola konfiguracji):**

```
profil_izolacji:
  nazwa:            <nazwa profilu>
  poziom_zasiegu:   globalny | srodowisko | modul | para_modulow | projekt | karta_sesji | rola
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

1. Operator wprowadza zmianę na komputerze (na przykład edycja wpisu pamięci).
2. Klient przesyła polecenie do serwera kanałem WebSocket (rozdz. 11).
3. Serwer zapisuje zmianę jako jedyne źródło prawdy i inicjuje przekazanie jej do pozostałych urządzeń.
4. Telefon i tablet otrzymują zdarzenie na żywo i prezentują ten sam stan w tej samej chwili (rozdz. 16).

### A.4. Dodanie nowego kanału modelu

1. Powstaje potrzeba podłączenia modelu przez kanał dotąd nieobsługiwany.
2. Deweloper dostarcza nowy adapter warstwy dostawcy modelu dla tego kanału (rozdz. 9).
3. Rdzeń, operujący na abstrakcji „wyślij zapytanie, odbierz odpowiedź", pozostaje bez zmian.
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

---

*Koniec dokumentu. Danaco Console — Architektura techniczna, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
