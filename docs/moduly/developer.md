# Danaco Console — Moduł Developer

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
| **Tytuł** | Moduł Developer |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant (co, gdzie, w jakiej formie) · deweloper (co zbudować) |
| **Przeznaczenie** | Ustala interfejs modułu Developer: pełny katalog funkcji, narzędzi, okien operacyjnych i punktów sterowania modułu dostępnego wyłącznie w CodeStudio |
| **Zakres** | okna operacyjne modułu, katalog funkcji i narzędzi, katalog elementów interfejsu, komendy kontraktu obszaru `developer` |
| **Poza zakresem** | izolacja techniczna procesu jako mechanizm rdzenia — [Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Model danych](../architektura/model-danych.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) · [Rozszerzenia](../architektura/rozszerzenia.md) · [Integracja modeli](../architektura/integracja-modeli.md) · [Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/developer.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `developer`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność ([Koncepcja platformy](../architektura/koncepcja-platformy.md#14-zasady-nadrzędne-funkcjonalności), rozdz. 14); zero blokad w interfejsie; klucze konfiguracji jawne; izolacja techniczna procesu i zakres uprawnień pozostają w gestii Operatora, nie są wymogiem; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
   - [1.1 Definicja](#11-definicja)
   - [1.2 Dla kogo](#12-dla-kogo)
   - [1.3 Po co — wartość modułu](#13-po-co--wartość-modułu)
   - [1.4 Zakres tematyczny modułu](#14-zakres-tematyczny-modułu)
   - [1.5 Rozgraniczenie z innymi modułami](#15-rozgraniczenie-z-innymi-modułami)
   - [1.6 Miejsce w architekturze platformy](#16-miejsce-w-architekturze-platformy)
   - [1.7 Dostępność i forma udostępnienia](#17-dostępność-i-forma-udostępnienia)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - [2.1 Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
   - [3.1 Chat Window (kanał Użytkownik ↔ Wykonawca)](#31-chat-window-kanał-użytkownik--wykonawca)
   - [3.2 Execution Loop Window (kanał Koordynator ↔ Wykonawca)](#32-execution-loop-window-kanał-koordynator--wykonawca)
   - [3.3 Code Editor](#33-code-editor)
   - [3.4 Project Tree](#34-project-tree)
   - [3.5 Git Panel](#35-git-panel)
   - [3.6 Build Output i Run & Debug](#36-build-output-i-run--debug)
   - [3.7 Dev Tools](#37-dev-tools)
4. [Przepływy pracy](#4-przepływy-pracy)
   - [4.1 Przepływ podstawowy — edycja ze wsparciem AI i zatwierdzenie zmian](#41-przepływ-podstawowy--edycja-ze-wsparciem-ai-i-zatwierdzenie-zmian)
   - [4.2 Przepływ rozszerzony — budowanie, niepowodzenie testu, poprawka](#42-przepływ-rozszerzony--budowanie-niepowodzenie-testu-poprawka)
   - [4.3 Przepływ pętli wykonawczej — zadanie wielokrokowe](#43-przepływ-pętli-wykonawczej--zadanie-wielokrokowe)
   - [4.4 Przepływ rozszerzony — debugowanie zatrzymanego wykonania](#44-przepływ-rozszerzony--debugowanie-zatrzymanego-wykonania)
   - [4.5 Przepływ rozszerzony — repozytorium jako komponent budowy produktu w Apps](#45-przepływ-rozszerzony--repozytorium-jako-komponent-budowy-produktu-w-apps)
   - [4.6 Przepływ zbiorczy — równoległa praca nad wieloma repozytoriami](#46-przepływ-zbiorczy--równoległa-praca-nad-wieloma-repozytoriami)
5. [Stany, dane i powiązania](#5-stany-dane-i-powiązania)
   - [5.1 Model stanów sesji modułu](#51-model-stanów-sesji-modułu)
   - [5.2 Model danych wykorzystywany przez moduł](#52-model-danych-wykorzystywany-przez-moduł)
   - [5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Developer](#53-izolacja-i-konfigurowalność--punkty-właściwe-modułowi-developer)
   - [5.4 Powiązania z innymi modułami](#54-powiązania-z-innymi-modułami)
6. [Scenariusze użycia](#6-scenariusze-użycia)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
   - [7.1 Rdzeń edytora kodu (Code Editor)](#71-rdzeń-edytora-kodu-code-editor)
   - [7.2 Nawigacja po projekcie (Project Tree)](#72-nawigacja-po-projekcie-project-tree)
   - [7.3 Kontrola wersji (Git Panel)](#73-kontrola-wersji-git-panel)
   - [7.4 Uruchamianie i budowanie (Build Output)](#74-uruchamianie-i-budowanie-build-output)
   - [7.5 Debugowanie (Run & Debug)](#75-debugowanie-run--debug)
   - [7.6 Klient API i HTTP (API Client)](#76-klient-api-i-http-api-client)
   - [7.7 Bazy danych i dane (Data Console)](#77-bazy-danych-i-dane-data-console)
   - [7.8 Kontenery i środowiska wykonawcze (Containers)](#78-kontenery-i-środowiska-wykonawcze-containers)
   - [7.9 Zależności, bezpieczeństwo, jakość](#79-zależności-bezpieczeństwo-jakość)
   - [7.10 Operacje kontekstowe AI](#710-operacje-kontekstowe-ai)
   - [7.11 Zestawienie zależności bibliotecznych (warstwa serwerowa, Go)](#711-zestawienie-zależności-bibliotecznych-warstwa-serwerowa-go)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
   - [8.1 Edytor i jakość kodu](#81-edytor-i-jakość-kodu)
   - [8.2 Git i wersjonowanie](#82-git-i-wersjonowanie)
   - [8.3 Budowanie, uruchamianie, debugowanie](#83-budowanie-uruchamianie-debugowanie)
   - [8.4 Integracje deweloperskie — widoczność zakładek Dev Tools](#84-integracje-deweloperskie--widoczność-zakładek-dev-tools)
   - [8.5 AI i model](#85-ai-i-model)
   - [8.6 Izolacja i zależności — powiązania jawne](#86-izolacja-i-zależności--powiązania-jawne)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)
10. [Punkty łamania i kryteria odbioru](#10-punkty-łamania-i-kryteria-odbioru)
   - [10.1 Punkty łamania](#101-punkty-łamania)
   - [10.2 Kryteria odbioru](#102-kryteria-odbioru)
11. [Załącznik — pełny wykaz komend kontraktu modułu Developer](#załącznik--pełny-wykaz-komend-kontraktu-modułu-developer)
   - [Obszar `developer` — 50 komend](#obszar-developer--50-komend)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja

Developer jest modułem tworzenia i rozwoju kodu — kompletnym środowiskiem wytwarzania oprogramowania zintegrowanym ze wsparciem AI działającym bezpośrednio w kontekście danego repozytorium. Moduł nie jest edytorem tekstu z dopiętym oknem rozmowy, lecz jedną przestrzenią roboczą, w której struktura projektu, treść kodu, historia zmian, przebieg debugowania i wynik budowania współdzielą ten sam kontekst przekazywany modelowi — Wykonawca zna repozytorium, nad którym trwa praca, jego strukturę oraz stan ostatniego budowania, bez ręcznego wklejania tych informacji do rozmowy.

Moduł stanowi jedną przestrzeń zastępującą edytor i zintegrowane środowisko programistyczne, klienta Git, klienta API, przeglądarkę baz danych, menedżera kontenerów, debugger, profiler, przeglądarkę zależności, skaner bezpieczeństwa, generator dokumentacji API oraz narzędzia towarzyszące. Deweloper pracujący w Danaco Console realizuje pełny cykl — nawigacja → edycja → uruchomienie → debugowanie → test → przegląd → wersjonowanie → wydanie — bez uruchamiania zewnętrznego programu deweloperskiego.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Developer |
|---|---|
| Deweloperzy warstwy serwerowej i warstwy klienckiej | Pisanie, refaktoryzacja i przegląd kodu z bezpośrednim wsparciem AI |
| Liderzy techniczni i architekci | Analiza architektury repozytorium, przegląd zmian przed scaleniem gałęzi |
| Inżynierowie jakości | Generowanie i uruchamianie testów, obserwacja pokrycia kodu, testy mutacyjne |
| Zespoły utrzymaniowe | Refaktoryzacja istniejącego kodu, dokumentowanie modułów bez własnej dokumentacji |
| Nowi członkowie zespołu | Nawigacja i zrozumienie nieznanego repozytorium przy wsparciu wyjaśnień AI |

### 1.3. Po co — wartość modułu

| Problem klasycznego rozproszenia narzędzi | Rozwiązanie w Developer |
|---|---|
| Edytor kodu, klient Git i okno rozmowy z Wykonawcą to trzy osobne aplikacje | Code Editor, Git Panel i Chat Window w jednej karcie sesji, nad tym samym repozytorium |
| Wykonawca nie ma dostępu do struktury projektu bez ręcznego wklejania plików | Project Tree i zawartość otwartych plików stanowią naturalny kontekst poleceń |
| Wynik budowania sprawdzany w osobnym terminalu, oderwanym od edytora | Build Output aktualizuje się na żywo, z bezpośrednim przejściem z komunikatu błędu do linii kodu |
| Wiadomości commitów pisane pospiesznie, bez podsumowania zmian | Generowanie opisu commitu na podstawie rzeczywistego różnicowania zmian |
| Debugowanie, zapytania HTTP, konsola bazy danych i kontenery w czterech odrębnych programach | Dev Tools i Run & Debug w tym samym oknie modułowym, na tym samym kontekście repozytorium |
| Ręczne przełączanie się do zewnętrznego terminala przy budowaniu i uruchamianiu | Warstwa wykonawcza dostarczana przez moduł Terminal, powiązana konfiguracyjnie |

### 1.4. Zakres tematyczny modułu

| Obszar | Zakres |
|---|---|
| Edycja i nawigacja po kodzie | Wieloplikowy edytor z LSP, wielojęzyczne podświetlanie, refaktoryzacje semantyczne, nawigacja po symbolach, wyszukiwanie globalne |
| Kontrola wersji | Pełny klient Git (lokalny i zdalny), przeglądy różnic, gałęzie, konflikty, integracja z hostingiem repozytoriów |
| Uruchamianie i budowanie | Zadania budowania, uruchomienia, tryb obserwacji zmian; warstwa wykonawcza dostarczana przez moduł Terminal |
| Debugowanie | Debugger krokowy oparty na DAP, punkty przerwania, podgląd zmiennych, stos wywołań |
| Testowanie i jakość | Uruchamianie testów, pokrycie, mutacje, lint, analiza statyczna, formatowanie |
| Integracje deweloperskie | Klient HTTP/API, przeglądarka i konsola baz danych, kontenery i obrazy, zmienne środowiskowe i sekrety |
| Wsparcie AI w kontekście repozytorium | Generowanie, refaktoryzacja, przegląd, dokumentacja, wyjaśnianie, konwersja, agentowe zadania wielokrokowe |

### 1.5. Rozgraniczenie z innymi modułami

| Poza zakresem modułu | Właściwy moduł i powód |
|---|---|
| Surowa powłoka systemowa i sesje terminala | **Terminal** — Developer korzysta z niego jako warstwy wykonawczej (powiązanie konfiguracyjne) |
| Analiza logów produkcyjnych, agregacja stanu systemu, rekomendacje naprawcze | **Diagnostics** — Developer przekazuje mu błędy jako materiał źródłowy |
| Złożenie kompletnego produktu z wielu komponentów (warstwa kliencka, warstwa serwerowa, infrastruktura) | **Apps** — Developer jest jednym z komponentów procesu Product Builder |
| Projekt graficzny interfejsu, makiety, tokeny wizualne | **Design** |
| Orkiestracja agentów i harmonogramy automatyzacji | **Agents**, **Automations** |
| Zaawansowana edycja dokumentów tekstowych (PDF, DOCX, redakcja) | **Studio** |

Rozgraniczenie jest miękkie i konfigurowalne — konsola bazodanowa Developera współdzieli połączenie z modułem analitycznym, a klient API eksportuje kolekcje do Automations, gdy powiązanie zostanie skonfigurowane. Powiązania te są jawne w oknie konfiguracji (rozdz. 8), nigdy wbudowane na stałe.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: CodeStudio
        ▼
ŚRODOWISKO: CODESTUDIO ── boczna nawigacja modułów
        │        (Workspace · Roundtable · Design · Terminal ·
        │         Developer · Diagnostics · Apps · Agents)
        ▼
MODUŁ: DEVELOPER ───────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Code Editor · Project Tree ·
  Git Panel · Build Output i Run & Debug · Dev Tools
        │
        ▼
KARTA SESJI — repozytorium powiązane z sesją, własny kontekst
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 5.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | CodeStudio (wyłącznie) |
| Środowiska TalkIn i WorkSpace | Moduł niedostępny — praca z kodem źródłowym należy wyłącznie do trybu programistycznego |
| Komponent własny | Developer nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 7 (łącznie z Chat Window i Execution Loop Window) |
| Typ pracy | Sesyjna, repozytoryjna — jedna karta sesji pracuje nad jednym repozytorium lub jego wskazanym fragmentem |
| Układ okien | Układ pionowy (podział lewa–prawa); regulacji podlega wyłącznie szerokość kolumn |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Rola w module | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Główne okno komunikacji i podstawowy mechanizm sterowania procesami modułu | 1 | Widoczne bez interakcji po wejściu do modułu |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość | Pętla wykonawcza zadań programistycznych: dekompozycja zlecenia, kolejka zadań, kontrola jakości, sterowanie przebiegiem | 1 (przy zleceniu w toku) · 2 (poza zleceniem) | Otwiera się samoczynnie z chwilą przyjęcia zlecenia wielokrokowego; poza zleceniem wywoływane znacznikiem pętli w pasku kontekstu |
| 3 | Code Editor | Okno edycyjne | Kolumna dominująca obszaru roboczego | Edycja treści plików źródłowych | 1 | Widoczne bez interakcji; treść wypełnia otwarcie pliku z Project Tree |
| 4 | Project Tree | Okno list i źródeł | Kolumna nawigacyjna obszaru roboczego, wąska | Nawigacja po strukturze projektu | 1 | Widoczne bez interakcji |
| 5 | Git Panel | Okno zarządcy (managera) i podgląd | Kolumna boczna, otwierana jako rozszerzenie boczne | Zarządzanie zmianami i kontrola wersji | 2 | Znacznik gałęzi w pasku kontekstu, ikona kontroli wersji, polecenie języka naturalnego; panel zwija się po zatwierdzeniu zmian |
| 6 | Build Output i Run & Debug | Monitor procesu | Kolumna monitora obszaru roboczego, na żywo | Wynik kompilacji, budowania i testów; debugger krokowy | 2 (Build Output) · 3 (Run & Debug) | Build Output otwiera się samoczynnie z chwilą uruchomienia budowania lub testów; Run & Debug wywoływane zakładką monitora, punktem przerwania na marginesie Code Editor lub poleceniem języka naturalnego |
| 7 | Dev Tools | Okno integracji wielofunkcyjne | Kolumna boczna, otwierana jako rozszerzenie boczne | Klient API, konsola bazodanowa, kontenery, zależności i bezpieczeństwo | 3 | Menu ☰ obszaru roboczego, wyszukiwarka funkcji, polecenie języka naturalnego; zakładki niewłączone w konfiguracji pozostają w warstwie 4 |

```
 Makieta zbiorcza — Moduł Developer, stan spoczynku     Dostępność: CodeStudio
 ═══════════════════════════════════════════════════════════════════════════════
  [Danaco Console] [CodeStudio] [Worktree] [feature/oplaty ▼] [Model ▼]     [☰]
 ═══════════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window          │ Project  │ Code Editor
  nawigacja  │ Użytkownik ↔         │ Tree     │ treść pliku źródłowego
  modułów    │ Wykonawca            │ struk-   │
  (poza      │                      │ tura     │
  zakresem   │ historia rozmowy     │ katalo-  │
  tego       │                      │ gów      │
  dokumentu) │                      │ i plików │
             │ ──────────────────   │          │
             │ [📎] pole poleceń    │      [⋮] │ TypeScript · UTF-8       [⋮]
             │           [Wyślij ▶] │          │
 ═══════════════════════════════════════════════════════════════════════════════

 Warstwa 1 zajmuje całą powierzchnię: Chat Window, Project Tree, Code Editor
 oraz pasek kontekstu. Warstwy 2–3 pozostają zwinięte — znaczniki kontekstowe
 ([feature/oplaty ▼], [Model ▼]), menu ⋮ kolumn, menu ☰ obszaru roboczego.
 Execution Loop Window, Build Output i Run & Debug, Git Panel oraz Dev Tools
 nie zajmują powierzchni do chwili wywołania.
```

### 2.1. Warstwy widoczności w module

Moduł Developer stosuje zasadę nadrzędną interfejsu Danaco Console: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Pełny arsenał narzędzi programistycznych — edytor, kontrola wersji, budowanie, debugger, klient API, konsola bazodanowa, kontenery, analiza zależności i bezpieczeństwa — istnieje w architekturze modułu i pozostaje niewidoczny w interfejsie do chwili wystąpienia potrzeby użycia. Liczba narzędzi modułu nie wpływa na postrzeganą prostotę jego przestrzeni roboczej.

**Warstwa 1 — zawsze widoczna.** Chat Window jako kanał Użytkownik ↔ Wykonawca, Project Tree jako nawigacja po strukturze repozytorium, Code Editor jako okno wiodące z treścią pliku źródłowego, pasek kontekstu pracy (platforma, środowisko, katalog roboczy, gałąź, model) oraz wskaźniki stanu wykonania: plakietka statusu budowania, licznik zadań pętli wykonawczej, liczba błędów analizy statycznej w pasku statusu edytora. Elementy warstwy 1 zajmują ponad 80% powierzchni interfejsu modułu. Execution Loop Window należy do warstwy 1 na czas prowadzenia zlecenia wielokrokowego — otwiera się samoczynnie z chwilą przyjęcia zlecenia i zamyka po jego zakończeniu.

**Warstwa 2 — widoczna na żądanie, zwijana po użyciu.** Znaczniki kontekstowe paska kontekstu i wywoływane nimi selektory: selektor gałęzi (`[feature/oplaty ▼]`), selektor kanału modelu (`[Model ▼]`), selektor profilu budowania, selektor konfiguracji uruchomienia, selektor połączenia bazodanowego, przełącznik widoku Project Tree (drzewo / zmiany), wskaźnik i lista plików włączonych do kontekstu polecenia, wskaźnik pobierz/wyślij repozytorium zdalnego. Do warstwy 2 należą również Git Panel oraz Build Output — otwierane jako rozszerzenia boczne odpowiednio przy pracy nad zmianami i przy uruchomieniu budowania lub testów, znikające z przestrzeni roboczej po zamknięciu. Selektor po dokonaniu wyboru zwija się samoczynnie do znacznika.

**Warstwa 3 — rozwinięcia kontekstowe.** Menu kebab (⋮) kolumn i pozycji: menu operacji pliku i folderu w Project Tree, menu operacji zdalnych repozytorium w Git Panel, menu ustawień sesji w Chat Window, menu punktu przerwania. Menu hamburger (☰) obszaru roboczego otwierające kolumnę Dev Tools z zakładkami API Client, Data Console, Containers oraz Dependencies & Security. Panele popover i listy rozwijane: pasek pływający zaznaczenia AI nad zaznaczonym fragmentem kodu, panel wyników wyszukiwania Grep, panel podglądu różnicy zbiorczej pętli wykonawczej, panel rozwiązywania konfliktu scalania, panel zmiennych i stosu debugowania, rozbicie pokrycia kodu na pliki, szczegóły pozycji podatności. Zestawy akcji prezentowane są jako jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę pozycji.

**Warstwa 4 — funkcje eksperckie.** Dostępne wyłącznie poleceniem języka naturalnego w Chat Window, skrótem klawiszowym, wyszukiwarką funkcji (paleta poleceń `Ctrl/Cmd + K`) albo w trybie administracyjnym; użytkownik podstawowy nie widzi tych elementów w interfejsie. Należą tu: wysłanie wymuszone (force push), zmiana bazy (rebase), przeniesienie commitu (cherry-pick), poprawienie ostatniego commitu (amend), podpisywanie commitów i konfiguracja tożsamości commitującej, zarządzanie zdalnymi repozytoriami, zamiana masowa wyników Grep w całym repozytorium, testy mutacyjne, migracje schematu i edycja danych w siatce Data Console, wypchnięcie obrazu do rejestru i uruchomienie kontenera deweloperskiego, konfiguracja adapterów debugowania i punktów logujących, próg pokrycia kodu, konsola debugowania REPL, snippet manager, widok szesnastkowy plików binarnych oraz punkty sterowania z okna konfiguracji (rozdz. 8). Zakładki Dev Tools niewłączone ustawieniem konfiguracyjnym z rozdz. 8.4 pozostają w tej warstwie i nie są prezentowane jako zablokowany element interfejsu.

**Zasada jednego kliknięcia.** Każda funkcja modułu ukryta w warstwach 2–4 jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego skierowanym do Wykonawcy. Zagnieżdżanie funkcji głęboko w hierarchii menu jest w module zabronione — ukrycie zmniejsza chaos wizualny przestrzeni programistycznej, nie utrudnia dostępu do narzędzia.

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy). Stanowi centralny punkt pracy użytkownika w module Developer i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez moduł: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu.

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |

**Zawartość i pełny arsenał funkcji.**

- Historia rozmowy z Wykonawcą w kontekście bieżącego repozytorium — bieżącego pliku otwartego w Code Editor oraz jego sąsiedztwa w Project Tree.
- Operacje kontekstowe wywoływane wprost z poziomu okna: generowanie kodu, refaktoryzacja, wyjaśnienie fragmentu, dokumentacja, generowanie testów, przegląd kodu, wykrywanie błędów, optymalizacja wydajności, konwersja między językami lub bibliotekami.
- Odwołanie do zaznaczonego fragmentu w Code Editor jako przedmiotu polecenia („zastosuj do zaznaczenia”).
- Odwołanie do całego pliku lub wielu plików naraz jako kontekstu, wskazywanych z poziomu Project Tree (wielokrotny wybór).
- Wstawienie sugestii Wykonawcy do Code Editor z możliwością akceptacji, odrzucenia lub edycji przed wstawieniem; wynik prezentowany jako różnica do zatwierdzenia.
- Generowanie treści commitu na podstawie rzeczywistego różnicowania zmian widocznego w Git Panel oraz noty wydania dla zakresu commitów.
- Wyjaśnienie komunikatu błędu lub niepowodzenia testu widocznego w Build Output oraz stanu zatrzymania debugowania, z sugestią poprawki.
- Zlecenie zadania wielokrokowego obejmującego wiele plików — plan, edycje, uruchomienie testów, korekta — którego przebieg prowadzony jest w Execution Loop Window.
- Cytowanie fragmentów kodu w rozmowie blokiem `.dn-karta` w typografii `--dn-ff-mono`, z zachowaniem podświetlania składni języka.
- Załączanie plików (zrzutu błędu, pliku konfiguracyjnego) do wiadomości.
- Historia poleceń (nawigacja strzałkami), edycja własnego polecenia i ponowne wysłanie, rozgałęzianie wątku odpowiedzi.
- Wskaźnik długości kontekstu bieżącej rozmowy oraz liczby plików włączonych do kontekstu, z ręcznym dodaniem i usunięciem pozycji.
- Szybkie menu kontekstowe konfiguracji sesji: kanał modelu, zakres pamięci, izolacja techniczna procesu repozytorium.
- Zatwierdzanie i przerywanie działań prowadzonych przez Wykonawcę, w tym przerwanie trwającego zadania wielokrokowego.

**Makieta tekstowa.**

```
┌─ Chat Window ───────────────┐
│[Model ▼] [kontekst: 3 ▼] [⋮]│
├─────────────────────────────┤
│ Historia rozmowy            │
│ (przewijana)                │
│ ┌─ Użytkownik ────────────┐ │
│ │ „napisz test jednostkowy│ │
│ │ dla zaznaczonej funkcji”│ │
│ └─────────────────────────┘ │
│ ┌─ Wykonawca (strumień) ──┐ │
│ │ ```ts                   │ │
│ │ describe("obliczOplate",│ │
│ │   () => { … })          │ │
│ │ ```                     │ │
│ │           [ Operacje ▼ ]│ │
│ └─────────────────────────┘ │
├─────────────────────────────┤
│ [📎] Pole poleceń ……………………  │
│                [ Wyślij ▶ ] │
└─────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Wskaźnik kontekstu plików | Etykieta tekstowa w nagłówku | Pokazuje, które pliki są włączone do kontekstu bieżącego polecenia | Mały tekst `--dn-tekst-3` z licznikiem | domyślny · rozwinięty (lista plików) | Kliknięcie otwiera listę plików kontekstu z możliwością usunięcia pozycji | Nagłówek Chat Window | 2 | Znacznik kontekstowy w nagłówku kolumny; kliknięcie rozwija listę plików, po zamknięciu element zwija się do licznika |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia do Wykonawcy | Duże pole w stopce kolumny, rozciągliwe w pionie | domyślny · fokus · z załącznikiem · błąd (pusta wysyłka) | Enter wysyła, Shift+Enter nowa linia | Stopka Chat Window | 1 | Widoczne bez interakcji |
| Blok kodu w odpowiedzi | `.dn-karta`, typografia `--dn-ff-mono`, podświetlenie składni | Prezentacja wygenerowanego lub zmodyfikowanego kodu | Średni lub duży blok w treści odpowiedzi | strumieniowanie · kompletny | Przycisk kopiowania w rogu bloku przy wskazaniu kursorem | Treść odpowiedzi Wykonawcy | 1 | Widoczny w strumieniu odpowiedzi Wykonawcy |
| Przycisk „Wstaw do edytora” | Przycisk `--zarys`, mały | Przeniesienie fragmentu kodu do Code Editor | Mały przycisk tekstowy w pasku akcji bloku | domyślny · wskazanie kursorem | Otwiera podgląd wstawienia (różnica przed wstawieniem) w miejscu kursora lub zaznaczenia | Pod blokiem kodu odpowiedzi | 2 | Pasek akcji bloku kodu; widoczny wyłącznie przy bloku zawierającym kod |
| Przycisk „Uruchom testy” | Przycisk `.dn-btn--sygnal`, mały | Wykonanie wygenerowanego testu w warstwie wykonawczej Terminal | Mały przycisk wyróżniony akcentem sygnałowym | domyślny · ładowanie | Przekazuje polecenie do Terminala (gdy powiązanie skonfigurowane), wynik trafia do Build Output | Pod blokiem kodu odpowiedzi zawierającym test | 2 | Pasek akcji bloku kodu; widoczny wyłącznie przy bloku zawierającym test |
| Przycisk ikony ustawień sesji | Ikona (ustawienia) | Szybka zmiana konfiguracji warstwy sesji | Mała ikona 36×36 px | domyślny · rozwinięte menu | Otwiera uproszczone menu kontekstowe: kanał modelu, pamięć, izolacja techniczna | Prawy róg nagłówka | 3 | Ikona ustawień otwiera menu kontekstowe (kanał modelu, pamięć, izolacja techniczna); menu zamyka się po wyborze |
| Przycisk załącznika | Ikona (spinacz) | Dołączenie pliku do wiadomości | Mała ikona | domyślny · aktywny (plik dołączony) | Otwiera okno wyboru pliku lub akceptuje przeciągnięcie | Lewa strona pola poleceń | 2 | Ikona spinacza przy polu poleceń albo przeciągnięcie pliku na okno |

---

### 3.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą i odpowiada za prowadzenie pętli wykonawczej modułu Developer: dekompozycję zlecenia programistycznego na zadania, koordynację ich realizacji, nadzór nad przebiegiem, orkiestrację działań oraz kontrolę realizacji procesów. Okno otwierane jest jako kolumna sąsiadująca z Chat Window.

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Pętla wykonawcza właściwa karcie sesji i jej repozytorium |

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie programistyczne i jego dekompozycja na zadania właściwe modułowi: analiza repozytorium, edycja wskazanych plików, refaktoryzacja symbolu w wielu miejscach, wygenerowanie i uruchomienie testów, uruchomienie budowania, uruchomienie sesji debugowania, przygotowanie i zatwierdzenie zmian w Git Panel.
- Kolejka i stan zadań: oczekujące, w realizacji, zakończone powodzeniem, zakończone niepowodzeniem, ponowione.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą — przydział zadania, raport wykonania, żądanie uzupełnienia kontekstu, zgłoszenie kolizji z bieżącym stanem katalogu roboczego.
- Wyniki kontroli jakości przypisane do zadań: diagnostyka analizy statycznej, wynik zestawu testów, pokrycie kodu, wynik przeglądu różnicy, wynik skanu zależności i sekretów; decyzja o ponowieniu zadania podejmowana na podstawie tych wyników.
- Wskaźniki przebiegu pętli: liczba zadań w kolejce, liczba ponowień, czas trwania zlecenia, liczba plików objętych zmianą, liczba przebiegów budowania.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia w toku oraz ponowienie wskazanego zadania.
- Podgląd zbiorczej różnicy wielo-plikowej wypracowanej przez pętlę, z akceptacją całości lub pojedynczych plików przed zapisem do katalogu roboczego.
- Przejście z pozycji zadania do miejsca jego skutku: pliku w Code Editor, wpisu w Build Output, pozycji w Git Panel.

**Makieta tekstowa.**

```
┌─ Execution Loop ────────────┐
│ Koordynator ↔ Wykonawca     │
│ Zlecenie: „wydziel warstwę  │
│ walidacji opłat”            │
├─────────────────────────────┤
│ Zadania (4/7)               │
│  ✔ analiza modułu opłat     │
│  ✔ wyodrębnienie funkcji    │
│  ✔ aktualizacja importów    │
│  ◐ generowanie testów       │
│  ○ uruchomienie budowania   │
│  ○ przegląd różnicy         │
│  ○ przygotowanie commitu    │
├─────────────────────────────┤
│ zadania 7 · ponowienia 1    │
│ [ Kontrola jakości ▼ ]      │
│ [ Komunikaty ▼ ]            │
├─────────────────────────────┤
│              [ Operacje ▼ ] │
└─────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Etykieta z treścią bieżącego zlecenia | Identyfikuje zlecenie prowadzone przez pętlę | Średni tekst z ikoną pętli | brak zlecenia · zlecenie aktywne · zlecenie zakończone | Kliknięcie rozwija pełną treść zlecenia i jego kontekst plikowy | Nagłówek kolumny | 1 | Widoczny bez interakcji przez czas prowadzenia zlecenia |
| Pozycja zadania | Wiersz `.dn-listwa-pozycja` z ikoną stanu | Reprezentuje pojedyncze zadanie pętli | Mały wiersz z etykietą i znacznikiem stanu | oczekujące (○) · w realizacji (◐) · powodzenie (✔) · niepowodzenie (✖) · ponowione | Kliknięcie otwiera skutek zadania w Code Editor, Build Output lub Git Panel | Lista zadań | 1 | Widoczna na liście zadań pętli w toku |
| Strumień komunikatów sterujących | Przewijana lista wiadomości Koordynator ↔ Wykonawca | Prezentuje przydziały zadań i raporty wykonania | Średni blok przewijany, tekst pomocniczy | strumieniowanie · kompletny | Kliknięcie pozycji rozwija pełną treść komunikatu | Środek kolumny | 2 | Element zbiorczy `Komunikaty ▼` w kolumnie pętli; zwija się po zamknięciu |
| Blok kontroli jakości | Zestawienie wyników analizy statycznej, testów i przeglądu | Podstawa decyzji o ponowieniu zadania | Mały blok z plakietkami wyników | brak wyników · wyniki częściowe · komplet wyników | Kliknięcie pozycji otwiera pełny wynik w Build Output | Kolumna, pod listą zadań | 2 | Znacznik zbiorczy wyników rozwijany kliknięciem; po zamknięciu wraca do postaci plakietki |
| Wskaźniki przebiegu pętli | Zestaw liczników | Liczba zadań, ponowień, czas trwania, liczba plików | Małe etykiety liczbowe | domyślny | — (informacyjne, aktualizowane na żywo) | Kolumna, przy bloku kontroli jakości | 1 | Widoczne bez interakcji jako wskaźniki stanu wykonania |
| Przyciski sterowania przebiegiem | Zestaw przycisków sterujących | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Średnie przyciski w rzędzie, `--zarys` i `--blad` | aktywny (domyślny) · aktywny z ostrzeżeniem (brak pętli w toku) · ładowanie | Wykonuje operację natychmiast; przy braku pętli w toku przycisk pozostaje aktywny i wyświetla ostrzeżenie inline zamiast blokady | Stopka kolumny | 3 | Element zbiorczy `Operacje ▼` w stopce kolumny; rozwinięcie zawiera wstrzymanie, wznowienie, przerwanie i korektę zlecenia |
| Podgląd różnicy zbiorczej | Panel różnicy wielo-plikowej | Akceptacja wyniku pętli przed zapisem | Duży panel rozwijany w kolumnie | ukryty · rozwinięty | Akceptacja całości lub pojedynczych plików zapisuje zmiany do katalogu roboczego | Wywoływany z pozycji zadania edycyjnego | 3 | Panel popover wywoływany z pozycji zadania edycyjnego; po zamknięciu znika z przestrzeni roboczej |

---

### 3.3. Code Editor

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne |
| Waga wizualna | Kolumna dominująca obszaru roboczego (największa powierzchnia modułu) |
| Izolacja domyślna | Treść plików odzwierciedla stan repozytorium w katalogu roboczym karty sesji |

**Zawartość i pełny arsenał funkcji.**

*Podstawa edycyjna.*
- Podświetlanie składni tree-sitter dla wielu języków programowania i formatów konfiguracyjnych, inkrementalne i składniowo poprawne.
- Numeracja linii, zwijanie i rozwijanie bloków kodu.
- Autouzupełnianie z podpowiedziami sygnatur funkcji i typów, dostarczane przez serwer języka (LSP), wraz z podglądem typu przy wskazaniu kursorem i diagnostyką na żywo.
- Wielokrotne karty otwartych plików, z podziałem widoku na kolejne kolumny.
- Wielokursorowość, zaznaczanie kolumnowe i zaznaczanie kolejnych wystąpień tego samego wzorca.
- Minimapa pliku dla szybkiej orientacji w długich plikach.

*Formatowanie i jakość kodu.*
- Automatyczne formatowanie całego dokumentu lub zaznaczenia, według reguł stylu przypisanych repozytorium, także przy zapisie pliku.
- Wybór stylu wcięć (spacje/tabulacje) i szerokości wcięcia, z odczytem z `.editorconfig`.
- Linter i analiza statyczna — podkreślenia błędów składniowych i ostrzeżeń na żywo, z szybką poprawką i regułami przypisanymi repozytorium.
- Adnotacja autora linii — kto i kiedy ostatnio zmienił daną linię, z przejściem do commitu.

*Nawigacja i wyszukiwanie.*
- Przejdź do definicji i implementacji, znajdź wszystkie wystąpienia symbolu, przejdź do symbolu w pliku, ścieżka symbolu w pasku okruchów.
- Znajdź i zamień w bieżącym pliku, z obsługą wyrażeń regularnych.
- Grep — znajdź w plikach (wyszukiwanie globalne w całym repozytorium) z podglądem wyników, zamianą masową i poszanowaniem reguł ignorowania repozytorium.
- Paleta poleceń — klawiaturowy dostęp do wszystkich funkcji edytora.
- Szybkie otwarcie pliku po fragmencie nazwy, bez opuszczania klawiatury.
- Zakładki w kodzie do szybkiego powrotu.

*Operacje kontekstowe AI.*
- Generuj — wygenerowanie nowego fragmentu kodu na podstawie opisu lub sąsiedniego kontekstu.
- Refaktoryzuj — przekształcenie zaznaczonego kodu bez zmiany zachowania.
- Wyjaśnij — opis działania zaznaczonego fragmentu w języku naturalnym.
- Udokumentuj — wygenerowanie komentarzy dokumentacyjnych zgodnych z konwencją języka.
- Napraw — sugestia poprawki na podstawie zgłoszonego błędu lub ostrzeżenia analizy statycznej.
- Napisz test — wygenerowanie testu jednostkowego dla zaznaczonej funkcji lub klasy.
- Zoptymalizuj — wydajniejsza implementacja zaznaczonego fragmentu.
- Konwertuj — przepisanie fragmentu między językami lub bibliotekami, ze wskazaniem źródła i celu.
- Refaktoryzacje semantyczne oparte na LSP: zmiana nazwy symbolu globalnie, wyodrębnienie funkcji lub zmiennej, wstawienie w miejscu, przeniesienie symbolu, organizacja importów — z podglądem różnicy wielo-plikowej.

*Wersjonowanie i porównanie z poziomu edytora.*
- Podgląd historii bieżącego pliku, przywrócenie wcześniejszej wersji fragmentu.
- Porównanie dwóch dowolnych plików niezależnie od Git Panel (Diff Panel wewnętrzny edytora), w trybie inline i dwukolumnowym.
- Wskaźniki zmian na marginesie linii (dodano/zmieniono/usunięto względem ostatniego commitu), z cofnięciem zmiany linii.

*Personalizacja i pomocnicze.*
- Snippet manager — własne, parametryzowane fragmenty kodu wielokrotnego użytku, wyzwalane prefiksem.
- Zmiana motywu kolorystycznego kodu, rozmiaru czcionki, zawijania wierszy.
- Podgląd na żywo dla plików renderowalnych: Markdown, HTML, Mermaid, SVG — w kolumnie obok źródła.
- Nakładka pokrycia testami na treści pliku.
- Eksport, pobranie pliku, kopiowanie ścieżki, kopiowanie pełnej zawartości.

**Makieta tekstowa.**

```
┌─ Code Editor ───────────────────────────────────────────────────────┐
│ orders.ts ×│ orders.test.ts ×│ +                        [🔍]  [⋮]   │
├─────────────────────────────────────────────────────────────────────┤
│  1  export function obliczOplate(kwota: number): number {           │
│  2 +   if (kwota < 0) throw new Error("kwota ujemna");   ← margines │
│  3      return kwota * STAWKA;                              zmian   │
│  4    }                                                             │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│TypeScript · UTF-8 · LF · wcięcia: 2 spacje │0 błędów · 1 ostrzeżenie│
└─────────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Karta otwartego pliku | Zakładka `.dn-zakladki` z nazwą pliku | Reprezentuje jeden otwarty plik | Mała zakładka z etykietą i kropką zmian niezapisanych | aktywna · w tle · niezapisane zmiany (kropka) · tylko odczyt | Kliknięcie przełącza plik; „×” zamyka kartę | Nagłówek kolumny edytora | 1 | Widoczna bez interakcji |
| Margines zmian | Wąski pasek przy numeracji linii | Sygnalizuje linie dodane, zmienione, usunięte względem ostatniego commitu | Bardzo wąski pasek kolorowy | bez zmian · dodano (zielony) · zmieniono (żółty) · usunięto (znacznik czerwony) | Wskazanie kursorem pokazuje mini-podgląd poprzedniej treści linii; kliknięcie cofa zmianę linii | Lewa krawędź obszaru treści | 1 | Widoczny bez interakcji |
| Obszar treści | Edytowalny obszar kodu | Główna przestrzeń pisania i odczytu kodu | Duży, dominujący blok, przewijany | pusty stan (brak otwartego pliku) · edycja · tylko odczyt | Wpisywanie, zaznaczanie, wielokursorowość | Środek kolumny | 1 | Widoczny bez interakcji |
| Pasek pływający zaznaczenia AI | Mały kontekstowy pasek narzędzi | Szybki dostęp do operacji kontekstowych AI na zaznaczonym kodzie | Mały, unoszący się nad zaznaczeniem | ukryty · widoczny (po zaznaczeniu ≥ 1 znaku) | Kliknięcie operacji wysyła zaznaczenie do Chat Window z gotowym poleceniem | Nad zaznaczonym fragmentem kodu | 3 | Panel popover pojawiający się nad zaznaczeniem; element zbiorczy `Więcej ▼` zawiera pełną listę operacji kontekstowych |
| Pole Grep w plikach | Pole wyszukiwania globalnego | Wyszukiwanie wzorca w całym repozytorium, z zamianą masową | Pole `.dn-pole-kontrolka` rozwijane w panel wyników | zwinięte · rozwinięte · wyniki listowane | Enter uruchamia wyszukiwanie; kliknięcie wyniku otwiera plik w miejscu trafienia | Pasek narzędzi edytora | 2 | Ikona wyszukiwania w pasku narzędzi albo skrót klawiszowy; pole zwija się po otwarciu wyniku. Zamiana masowa w całym repozytorium należy do warstwy 4 |
| Paleta poleceń | Okno nakładkowe wyszukiwarki poleceń | Klawiaturowy dostęp do wszystkich funkcji edytora | Wąskie okno nakładkowe wyśrodkowane, pole `.dn-pole-kontrolka` w nagłówku | zamknięta · otwarta (`Ctrl/Cmd + K`) | Wpisanie frazy filtruje listę poleceń; Enter wykonuje wybrane | Nakładka nad edytorem | 4 | Wyszukiwarka funkcji wywoływana skrótem `Ctrl/Cmd + K`; brak reprezentacji wizualnej w stanie spoczynku |
| Pasek statusu edytora | Wiersz informacyjny | Język pliku, kodowanie, znaki końca linii, wcięcia, liczba błędów analizy statycznej, gałąź Git | Mały, niski pasek, tekst pomocniczy | brak błędów · błędy/ostrzeżenia (plakietka z liczbą) | Kliknięcie liczby błędów otwiera listę problemów pliku | Stopka kolumny edytora | 1 | Widoczny bez interakcji |
| Wskaźnik analizy statycznej w linii | Podkreślenie faliste pod kodem | Sygnalizuje błąd lub ostrzeżenie statycznej analizy | Bardzo mały, wewnątrz linii kodu | brak · ostrzeżenie (żółte) · błąd (czerwone) | Wskazanie kursorem pokazuje treść komunikatu w `.dn-tooltip`, z podpowiedzią poprawki | W treści kodu | 1 | Widoczny w treści kodu; treść komunikatu w `.dn-tooltip` przy wskazaniu kursorem |

---

### 3.4. Project Tree

| Aspekt | Wartość |
|---|---|
| Typologia | Okno list i źródeł |
| Waga wizualna | Kolumna nawigacyjna obszaru roboczego, wąska, stale widoczna |
| Izolacja domyślna | Odzwierciedla stan systemu plików katalogu roboczego karty sesji |

**Zawartość i pełny arsenał funkcji.**

- Drzewo katalogów i plików repozytorium, z ikonami rozpoznającymi typ pliku.
- Status Git per plik i folder — plakietki: zmodyfikowany, nowy, usunięty, skonfliktowany, ignorowany; status agregowany na foldery.
- Wyszukiwanie i szybkie otwarcie pliku po fragmencie nazwy (bez rozwijania całego drzewa).
- Filtrowanie widoku według wzorca, zgodnie z regułami pliku ignorowania repozytorium.
- Zwijanie i rozwijanie gałęzi drzewa, zwinięcie wszystkiego jednym poleceniem.
- Menu kontekstowe pozycji: nowy plik, nowy folder, zmień nazwę, usuń, duplikuj, wytnij, kopiuj, wklej, otwórz w nowej karcie, otwórz obok (podział widoku), pokaż w eksploratorze systemowym, skopiuj ścieżkę.
- Przeciąganie i upuszczanie — przenoszenie plików i folderów w drzewie.
- Wielokrotny wybór plików — do operacji zbiorczych lub jako zestaw kontekstu przekazywany do Chat Window.
- Wskaźnik plików aktualnie otwartych w Code Editor (podświetlenie pozycji w drzewie).
- Odświeżenie drzewa — ręczna synchronizacja ze stanem systemu plików, obok automatycznej aktualizacji na żywo opartej na obserwacji systemu plików.
- Przełącznik widoku: pełne drzewo projektu albo lista wyłącznie zmienionych plików (widok Git).
- Etykiety własne i zakładki plików — oznaczanie pozycji wymagających uwagi lub przeglądu i szybki powrót do nich.
- Podgląd plików nietekstowych: obrazy, czcionki oraz widok szesnastkowy dla plików binarnych.

**Makieta tekstowa.**

```
┌─ Project Tree ───────────────────┐
│ [🔍]              [ Widok ▼ ]    │
├──────────────────────────────────┤
│ ▾ 📁 api                         │
│    ▾ 📁 routes                   │
│       📄 orders.ts        ● M    │
│       📄 orders.test.ts   ● N    │
│    📁 middleware                 │
│ ▾ 📁 tests                       │
│    📄 setup.ts                   │
│ 📄 package.json           ● M    │
│ 📄 .env.example                  │
├──────────────────────────────────┤
│                            [ ⋮ ] │
└──────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Węzeł folderu | Wiersz drzewa z ikoną folderu | Reprezentuje katalog repozytorium | Mały wiersz z trójkątem rozwijania | zwinięty · rozwinięty | Kliknięcie rozwija/zwija zawartość | Ciało drzewa | 1 | Widoczny bez interakcji |
| Węzeł pliku | Wiersz drzewa z ikoną typu pliku | Reprezentuje pojedynczy plik | Mały wiersz z etykietą nazwy i plakietką statusu | domyślny · otwarty w edytorze (wyróżnienie) · zaznaczony (wielokrotny wybór) | Kliknięcie otwiera plik w Code Editor; podwójne kliknięcie otwiera w trybie trwałym (nie podglądowym) | Ciało drzewa | 1 | Widoczny bez interakcji |
| Plakietka statusu Git | Mała litera w kółku (M/N/U/C) | Sygnalizuje status pliku w repozytorium | Bardzo mała plakietka `.dn-plakietka` | zmodyfikowany (M) · nowy (N) · usunięty (U) · skonfliktowany (C, kolor błędu) | Kliknięcie otwiera plik od razu w widoku różnicowym Git Panel | Prawa krawędź wiersza pliku | 1 | Widoczna bez interakcji przy pozycji pliku |
| Pole szybkiego otwarcia | Pole wyszukiwania nazwy pliku | Odnalezienie i otwarcie pliku bez przeglądania drzewa | Pole `.dn-pole-kontrolka`, nagłówek kolumny | puste · z wynikami listowanymi poniżej | Strzałki nawigują wyniki, Enter otwiera wybrany plik | Nagłówek kolumny | 2 | Ikona lupy w nagłówku kolumny albo skrót klawiszowy; pole zwija się po otwarciu pliku |
| Przełącznik widoku | Para przycisków radiowych | Pełne drzewo kontra wyłącznie zmienione pliki | Mały segmentowany przełącznik | drzewo (domyślny) · zmiany | Przelicza widoczną listę pozycji | Nagłówek kolumny, pod polem wyszukiwania | 2 | Znacznik kontekstowy `Widok ▼`; zwija się po wyborze |
| Menu kontekstowe pozycji | Menu rozwijane prawym klawiszem lub ikoną „⋯” | Operacje na pliku/folderze: nowy, usuń, zmień nazwę, kopiuj ścieżkę | Lista pozycji tekstowych w `.dn-modal`-podobnym kontenerze | zamknięte · otwarte | Wybór pozycji wykonuje operację natychmiast; „Usuń” domyślnie usuwa od razu, z dostępnym „Cofnij” w powiadomieniu po wykonaniu — okno nakładkowe potwierdzenia przed usunięciem włącza ustawienie `developer.git.confirm_delete` | Wywoływane z dowolnego węzła | 3 | Prawy klawisz na węźle albo ikona ⋮ pozycji |
| Przycisk „Odśwież” | Ikona (odswiez) | Ręczna synchronizacja drzewa ze stanem systemu plików | Mała ikona, stopka kolumny | domyślny · ładowanie (obrót ikony) | Ponownie odczytuje strukturę katalogu roboczego | Stopka kolumny | 3 | Pozycja menu ⋮ w stopce kolumny |

---

### 3.5. Git Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządcy (managera) z elementami podglądu i porównania |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne obszaru roboczego |
| Izolacja domyślna | Odzwierciedla stan repozytorium Git w katalogu roboczym karty sesji |

**Zawartość i pełny arsenał funkcji.**

*Status i przygotowanie zmian.*
- Lista zmian: niezatwierdzone, przygotowane do commitu, nieśledzone — z ikonami statusu.
- Dodanie do przygotowania (stage) pojedynczego pliku, zaznaczonych plików albo wszystkich naraz; cofnięcie przygotowania (unstage).
- Widok różnicowy per plik — inline albo dwukolumnowy, z podświetlaniem składni i porównaniem słowo po słowie, z zaznaczaniem pojedynczych fragmentów (hunków) do przygotowania częściowego.

*Commit.*
- Pole treści commitu z licznikiem znaków pierwszej linii.
- Generowanie opisu commitu na podstawie rzeczywistego różnicowania przygotowanych zmian.
- Zatwierdzenie zmian (commit), poprawienie ostatniego commitu (amend), cofnięcie ostatniego commitu (revert) z zachowaniem historii.
- Podpisywanie commitów (GPG/SSH) oraz konfiguracja tożsamości commitującej dla repozytorium.

*Historia i gałęzie.*
- Log commitów — lista z autorem, datą, skrótem, treścią, filtrami po autorze i zakresie dat, z przejściem do plików commitu.
- Wykres gałęzi (branch graph) — wizualizacja rozgałęzień i scaleń w czasie.
- Utworzenie nowej gałęzi, przełączenie (checkout), scalenie (merge), zmiana bazy (rebase), przeniesienie commitu (cherry-pick), usunięcie gałęzi.
- Porównanie dwóch dowolnych gałęzi lub commitów w Diff Panel wywoływanym z poziomu Git Panel.
- Tagowanie wersji, przegląd i usuwanie tagów.

*Zdalne repozytorium.*
- Pobranie (fetch), aktualizacja (pull), wysłanie (push); zarządzanie zdalnymi repozytoriami (remotes).
- Wysłanie wymuszone (force push) — operacja jawna, sygnalizowana wizualnie jako nieodwracalna.
- Odłożenie zmian roboczych (stash) i przywrócenie stasha; lista odłożonych zestawów zmian.
- Integracja z hostingiem repozytoriów: żądania scalenia (PR/MR), przypisania, statusy potoków CI i komentarze — obsługiwane bez opuszczania okna, przez integrację z API GitHub, GitLab lub Bitbucket.

*Konflikty i integracja z AI.*
- Rozwiązywanie konfliktów scalania — widok trójstronny (wersja bieżąca / wersja przychodząca / wynik) z wyborem fragmentu lub ręczną edycją.
- Podsumowanie zmian (changelog) i nota wydania generowane na podstawie zakresu commitów, zgodnie z konwencją Conventional Commits.
- Przegląd zmian przed zatwierdzeniem (przegląd kodu na różnicy) wykonywany przez Wykonawcę, z uwagami przypisanymi do konkretnych linii.

**Makieta tekstowa.**

```
┌─ Git Panel ─────────────────────────────────────────────────────────┐
│ Gałąź: [ feature/oplaty ▼ ]     ↓ 0  ↑ 2                      [ ⋮ ] │
├─────────────────────────────────────────────────────────────────────┤
│ Przygotowane do commitu (2)          Niezatwierdzone (1)            │
│  ☑ orders.ts               M          ☐ .env.example        M       │
│  ☑ orders.test.ts          N                                        │
├─────────────────────────────────────────────────────────────────────┤
│ Wiadomość commitu:                                                  │
│ [ dodaj walidację ujemnej kwoty opłaty          ] [ ✨ ]            │
│                                       [ Zatwierdź zmiany ]          │
├─────────────────────────────────────────────────────────────────────┤
│ [ Historia ▼ ]                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor gałęzi | Rozwijana lista z etykietą bieżącej gałęzi | Wybór lub utworzenie gałęzi | Mały przycisk z etykietą | domyślny · rozwinięty · tworzenie nowej gałęzi | Wybór wykonuje przełączenie (checkout); wpisanie nowej nazwy proponuje utworzenie gałęzi | Nagłówek kolumny | 2 | Znacznik kontekstowy `[feature/oplaty ▼]` w pasku kontekstu; zwija się po wyborze gałęzi |
| Wskaźnik pobierz/wyślij | Dwie ikony z licznikiem commitów | Sygnalizuje różnicę względem zdalnej gałęzi | Małe ikony ze strzałkami i liczbą | zsynchronizowane (ukryte liczniki) · rozbieżność (widoczna liczba) | Kliknięcie wykonuje odpowiednio pobranie lub wysłanie | Nagłówek kolumny | 1 | Widoczny bez interakcji w obrębie otwartego Git Panel |
| Lista zmian ze znacznikami wyboru | Lista `.dn-tabela`-podobna z polami `.dn-check` | Wybór plików do przygotowania w commicie | Średnia lista wierszy z plakietką statusu | pusta (brak zmian) · z pozycjami · częściowo zaznaczona | Zaznaczenie przenosi plik do sekcji „przygotowane”; kliknięcie nazwy otwiera widok różnicowy | Ciało kolumny | 1 | Widoczna bez interakcji w obrębie otwartego Git Panel |
| Blok różnicy (widok inline) | Fragment kodu z oznaczeniem koloru | Prezentacja zmiany liniowej w pliku | Średni blok, tło zależne od typu zmiany | dodanie (zielone) · usunięcie (czerwone) · kontekst (neutralne) | Zaznaczenie fragmentu pozwala przygotować częściowo (hunk) zamiast całego pliku | Widok różnicowy pliku | 2 | Kliknięcie nazwy pliku na liście zmian; widok zwija się po powrocie do listy |
| Pole wiadomości commitu | Pole tekstowe jednowierszowe z rozwinięciem | Treść opisu zatwierdzanych zmian | Średnie pole `.dn-pole-kontrolka` | puste · wypełnione · błąd (pusta wysyłka) | Enter lub przycisk zatwierdza commit | Środek kolumny | 1 | Widoczne bez interakcji w obrębie otwartego Git Panel |
| Przycisk „generuj AI” | Mały przycisk z ikoną iskry | Wygenerowanie treści commitu na podstawie różnicy | Mały przycisk `--zarys` obok pola wiadomości | domyślny · ładowanie | Wypełnia pole wiadomości opisem do edycji przed zatwierdzeniem | Obok pola wiadomości commitu | 2 | Ikona iskry obok pola wiadomości albo polecenie języka naturalnego w Chat Window |
| Przycisk „Zatwierdź zmiany” | Główny przycisk akcji kolumny | Wykonanie commitu przygotowanych zmian | Duży przycisk `--zloty`, pełna szerokość kolumny | aktywny (domyślny) · aktywny z ostrzeżeniem (brak przygotowanych zmian lub pustej wiadomości) · ładowanie | Tworzy commit, czyści sekcję przygotowanych zmian, dodaje wpis do historii; przy braku przygotowanych zmian lub pustej wiadomości przycisk pozostaje aktywny i wyświetla ostrzeżenie inline zamiast blokady | Sekcja commitu | 1 | Widoczny bez interakcji w obrębie otwartego Git Panel |
| Pozycja historii commitów | Wiersz `.dn-listwa-pozycja` z kropką i skrótem | Reprezentuje jeden commit w historii | Mały wiersz z awatarem autora, skrótem i czasem względnym | domyślny · wskazanie kursorem · rozwinięty (szczegóły) | Kliknięcie pokazuje pełną listę zmienionych plików tego commitu | Sekcja historii | 2 | Element zbiorczy `Historia ▼` w stopce kolumny; rozwinięcie prezentuje listę commitów |
| Przycisk „Wyślij wymuszony” | Przycisk `.dn-btn--niebezpieczny`, mały, w menu zaawansowanym | Nadpisanie zdalnej gałęzi historią lokalną | Mały przycisk tekstowy w rozwijanym menu operacji zdalnych | domyślny · ładowanie | Domyślnie wykonuje wysłanie natychmiast, sygnalizując inline ostrzeżenie o nieodwracalności dla współpracowników; okno nakładkowe potwierdzenia przed wysłaniem włącza ustawienie `developer.git.confirm_force_push` | Menu operacji zdalnych repozytorium | 4 | Wyłącznie polecenie języka naturalnego, wyszukiwarka funkcji albo menu operacji zdalnych w trybie administracyjnym; użytkownik podstawowy nie widzi tego elementu |
| Widok rozwiązywania konfliktu | Panel trójstronny | Wybór wersji fragmentu przy konflikcie scalania | Duży panel z trzema kolumnami (bieżąca / przychodząca / wynik) | brak konfliktu (ukryty) · konflikt aktywny | Wybór „przyjmij bieżącą”/„przyjmij przychodzącą”/edycja ręczna oznacza fragment jako rozwiązany | Zastępuje widok różnicowy w trakcie scalania | 3 | Panel zastępujący widok różnicy, wywoływany z pozycji pliku skonfliktowanego |

---

### 3.6. Build Output i Run & Debug

Okno monitora procesu łączy dwie części: **Build Output** — strumień budowania, testów i pokrycia — oraz **Run & Debug** — konfiguracje uruchomień i debugger krokowy oparty na protokole DAP. Obie części zajmują tę samą kolumnę monitora i przełączane są zakładkami w jej nagłówku.

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor procesu |
| Waga wizualna | Kolumna monitora obszaru roboczego, na żywo, o regulowanej szerokości |
| Izolacja domyślna | Log budowania i sesja debugowania właściwe karcie sesji; historia przebiegów narasta w toku sesji |

**Zawartość i pełny arsenał funkcji — Build Output.**

- Log procesu budowania aktualizowany na żywo kanałem WebSocket.
- Plakietka statusu budowania: w toku, sukces, błąd, przerwane.
- Wynik testów — lista z rozróżnieniem przeszedł/nie przeszedł/pominięty, pasek postępu wykonania zestawu.
- Pokrycie kodu testami (code coverage) — podsumowanie procentowe, z rozbiciem na pliki i nakładką na treść w Code Editor.
- Filtrowanie logu: błędy, ostrzeżenia, informacje.
- Wyszukiwanie pełnotekstowe w logu budowania.
- Przejście z komunikatu błędu wprost do odpowiedniej linii w Code Editor („kliknij, aby otworzyć”), na podstawie rozpoznania wzorca `plik:linia:kolumna`.
- Historia poprzednich przebiegów budowania — lista z czasem trwania, wynikiem i znacznikiem commitu, dla którego przebieg wykonano, z ponownym uruchomieniem wskazanego przebiegu.
- Definicje zadań budowania, uruchomienia, testów i trybu obserwacji zmian; wybór profilu budowania (deweloperski/produkcyjny) oraz docelowej platformy budowania.
- Tryb obserwacji zmian — automatyczne przebudowanie i przeładowanie po zapisie pliku.
- Uruchomienie budowania, zatrzymanie budowania w toku, ponowne uruchomienie ostatniego przebiegu.
- Eksport raportu budowania i wyniku testów do pliku (Markdown, JUnit XML, HTML).
- Przekazanie błędu budowania jako materiału źródłowego do modułu Diagnostics.
- Testy mutacyjne oraz próg pokrycia prezentowany jako ostrzeżenie, nie blokada.

**Zawartość i pełny arsenał funkcji — Run & Debug.**

- Konfiguracje uruchomień powiązane z repozytorium, z wyborem adaptera debugowania właściwego językowi.
- Debugger krokowy: uruchomienie, zatrzymanie, wejście do funkcji, przejście, wyjście, kontynuacja, ponowne uruchomienie.
- Punkty przerwania zwykłe i warunkowe, punkty logujące, przerwanie na wyjątku, licznik trafień; punkty ustawiane na marginesie Code Editor.
- Podgląd zmiennych w zakresach lokalnym i globalnym, wyrażenia obserwowane, edycja wartości w trakcie zatrzymania.
- Stos wywołań i lista wątków, nawigacja po ramkach i przełączanie wątków.
- Konsola debugowania (REPL) wykonująca wyrażenia w kontekście zatrzymanej ramki.
- Wyjaśnienie stanu zatrzymania — opis przyczyny na podstawie stosu i zmiennych oraz sugestia poprawki, kierowany do Chat Window.

**Makieta tekstowa.**

```
┌─ Build Output │ Run & Debug ─────────────────────────────────────────┐
│ [ Profil ▼ ]                                  [ Operacje ▼ ]   [ ⋮ ] │
├──────────────────────────────────────────────────────────────────────┤
│ Status: ● w toku          Testy: ▓▓▓▓▓▓▓░░░ 71 / 100                 │
│ 12:09:01  kompilacja modułu api…                                     │
│ 12:09:04  ⚠ ostrzeżenie: nieużywany import w orders.ts:3             │
│ 12:09:07  ✖ test nie przeszedł: obliczOplate › kwota ujemna          │
├──────────────────────────────────────────────────────────────────────┤
│ Pokrycie: 78%                                     [ Historia ▼ ]     │
└──────────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor profilu budowania | Rozwijana lista | Wybór konfiguracji budowania (deweloperski/produkcyjny) | Mały przycisk z etykietą | domyślny · rozwinięty | Zmienia parametry kolejnego uruchomienia budowania | Nagłówek kolumny | 2 | Znacznik kontekstowy `[Profil ▼]` w nagłówku monitora; zwija się po wyborze |
| Przycisk „Uruchom budowanie” | Przycisk `--zloty` | Rozpoczęcie procesu budowania | Średni przycisk wyróżniony akcentem sygnałowym | aktywny (domyślny) · aktywny z ostrzeżeniem (budowanie już w toku) · ładowanie | Uruchamia proces, przełącza plakietkę statusu na „w toku”; gdy budowanie już trwa, przycisk pozostaje aktywny i wyświetla ostrzeżenie inline zamiast blokady | Nagłówek kolumny | 2 | Pozycja elementu zbiorczego `Operacje ▼` w nagłówku monitora albo polecenie języka naturalnego |
| Przycisk „Zatrzymaj” | Przycisk `--blad`, mały | Przerwanie budowania w toku | Mały przycisk | aktywny (domyślny) · aktywny z ostrzeżeniem (brak procesu w toku) · ładowanie | Gdy budowanie trwa, kończy proces i log oznacza przebieg jako „przerwane”; gdy nic się nie buduje, przycisk pozostaje aktywny i wyświetla ostrzeżenie inline zamiast blokady | Nagłówek kolumny, obok „Uruchom” | 2 | Widoczny w nagłówku monitora wyłącznie w trakcie trwania procesu |
| Plakietka statusu | `.dn-plakietka--ostrzezenie` / `.dn-plakietka--sukces` / `.dn-plakietka--blad` / `.dn-plakietka` z kropką | Bieżący stan procesu budowania | Mała pigułka ze stanem | w toku (żółta, pulsująca) · sukces (zielona) · błąd (czerwona) · przerwane (szara) | — (informacyjna, aktualizowana na żywo) | Nagłówkowa część kolumny | 1 | Widoczna bez interakcji jako wskaźnik stanu wykonania |
| Pasek postępu testów | Pasek wypełnienia z licznikiem | Postęp wykonania zestawu testów | Średni pasek poziomy z etykietą liczbową | w toku · zakończony | — (informacyjny, aktualizowany na żywo) | Nagłówkowa część kolumny | 1 | Widoczny bez interakcji jako wskaźnik stanu wykonania |
| Linia logu błędu | Wiersz tekstu mono z ikoną | Pojedynczy komunikat błędu lub niepowodzenia testu | Tekst `--dn-ff-mono`, czerwone obramowanie, ikona `blad` | domyślny | Kliknięcie „Otwórz w edytorze” przenosi do pliku i linii źródła błędu | Ciało logu | 1 | Widoczna bez interakcji w obrębie otwartego monitora |
| Przycisk „Wyjaśnij niepowodzenie” | Mały przycisk tekstowy | Wysłanie komunikatu błędu do Chat Window z prośbą o wyjaśnienie i poprawkę | Mały przycisk `--zarys` | domyślny | Otwiera Chat Window z wypełnionym kontekstem błędu | Pod linią błędu w logu | 2 | Pojawia się pod linią błędu przy wskazaniu kursorem; znika po opuszczeniu wiersza |
| Punkt przerwania | Znacznik na marginesie linii kodu | Zatrzymanie wykonania w wskazanej linii | Mała kropka na marginesie Code Editor | zwykły · warunkowy · logujący · trafiony (podświetlony) | Kliknięcie ustawia lub usuwa punkt; menu kontekstowe definiuje warunek i licznik trafień | Margines Code Editor, sterowany z Run & Debug | 2 | Kliknięcie marginesu Code Editor; warunek i licznik trafień ustawiane z menu kontekstowego warstwy 3 |
| Panel zmiennych i stosu | Lista zakresów, zmiennych i ramek | Podgląd stanu programu w chwili zatrzymania | Średni blok list rozwijanych | brak sesji · zatrzymany · wykonanie w toku | Rozwinięcie zakresu pokazuje wartości; edycja wartości zapisuje ją w kontekście ramki; wybór ramki przełącza kontekst | Część Run & Debug | 3 | Zakładka Run & Debug monitora, otwierana samoczynnie przy zatrzymaniu wykonania |
| Przyciski kroku debugowania | Zestaw przycisków sterowania | Krok, wejście, wyjście, kontynuacja, restart, zatrzymanie | Rząd małych przycisków | aktywne (domyślne) · aktywne z ostrzeżeniem (brak sesji debugowania) | Wykonuje operację protokołu DAP; przy braku sesji przycisk pozostaje aktywny i wyświetla ostrzeżenie inline zamiast blokady | Część Run & Debug | 3 | Rząd przycisków w zakładce Run & Debug oraz skróty klawiszowe |
| Pozycja historii przebiegów | Mała karta w rzędzie poziomym | Podsumowanie jednego wcześniejszego przebiegu budowania | Mała karta z ikoną wyniku, czasem trwania, skrótem commitu | sukces (✔) · błąd (✖) | Kliknięcie otwiera pełny log tego przebiegu w trybie tylko do odczytu | Stopka kolumny | 3 | Element zbiorczy `Historia ▼` w stopce monitora |
| Wskaźnik pokrycia kodu | Etykieta procentowa | Podsumowanie pokrycia testami | Mały tekst z plakietką koloru zależnego od progu | poniżej progu (ostrzeżenie) · powyżej progu (sukces) | Kliknięcie otwiera rozbicie pokrycia na poszczególne pliki | Stopka kolumny | 2 | Plakietka w stopce monitora; rozbicie pokrycia na pliki otwiera się jako panel popover warstwy 3 |

---

### 3.7. Dev Tools

Dev Tools jest kolumną boczną gromadzącą okna integracji deweloperskich, przełączane zakładkami: **API Client**, **Data Console**, **Containers** oraz **Dependencies & Security**. Widoczność każdej zakładki steruje ustawieniem konfiguracyjnym z rozdz. 8.4; zakładka niewłączona pozostaje ukryta, nigdy nie jest prezentowana jako zablokowany element interfejsu.

| Aspekt | Wartość |
|---|---|
| Typologia | Okno integracji wielofunkcyjne |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne obszaru roboczego |
| Izolacja domyślna | Poświadczenia zdalnych usług, rejestrów kontenerów i baz danych podlegają zakresowi „konto i token per sesja” (rozdz. 5.3) |

**Zawartość i pełny arsenał funkcji — API Client.**

- Zapytania HTTP/REST: metoda, adres, nagłówki, treść (JSON, formularz, multipart), podgląd odpowiedzi, czas i status.
- Kolekcje i środowiska: grupowanie zapytań, zmienne środowiskowe w postaci `{{placeholder}}`, przełączanie środowisk.
- Import i eksport OpenAPI: generowanie zapytań z kontraktu OpenAPI/Swagger, walidacja odpowiedzi względem schematu.
- Zapytania GraphQL z introspekcją, wywołania gRPC na podstawie `.proto` oraz testy połączeń WebSocket.
- Generowanie kodu klienta lub testów integracyjnych z definicji zapytania.
- Zapytania wersjonowane w repozytorium jako plik `.http`, uruchamiane wprost z Code Editor.

**Zawartość i pełny arsenał funkcji — Data Console.**

- Przeglądarka połączeń i schematu: drzewo baz, tabel, kolumn i indeksów dla wielu połączeń jednocześnie.
- Konsola SQL: edytor z autouzupełnianiem schematu, wykonanie zapytania, siatka wyników, eksport CSV i JSON.
- Edycja danych w siatce: edycja wierszy w miejscu, wstawianie i usuwanie z generowaniem instrukcji DML, transakcja z podglądem przed zapisem.
- Migracje schematu: podgląd, tworzenie i uruchamianie migracji.
- Generowanie SQL z opisu w języku naturalnym na podstawie schematu połączenia.
- Podgląd planu zapytania (EXPLAIN/ANALYZE) z wizualizacją kosztów.

**Zawartość i pełny arsenał funkcji — Containers.**

- Lista kontenerów i obrazów, uruchomienie, zatrzymanie, restart, usunięcie, logi na żywo i statystyki zużycia zasobów.
- Budowanie obrazu z `Dockerfile`, tagowanie, wypchnięcie do rejestru.
- Uruchomienie stosu wielu usług z `docker-compose.yml` i podgląd stanu usług.
- Powłoka w kontenerze udostępniana przez warstwę wykonawczą modułu Terminal.
- Otwarcie repozytorium w kontenerze deweloperskim zgodnie z `devcontainer.json`.

**Zawartość i pełny arsenał funkcji — Dependencies & Security.**

- Przeglądarka zależności: drzewo zależności, wersje, pakiety przestarzałe, aktualizacja.
- Skan podatności zależności (SCA) z wykryciem znanych CVE i wskazaniem bezpiecznej wersji.
- Skan sekretów wykrywający klucze i tokeny w kodzie przed zatwierdzeniem zmian.
- Analiza bezpieczeństwa kodu (SAST) wykrywająca wzorce podatności.
- Analiza licencji zależności z flagowaniem pozycji niezgodnych z polityką.

**Makieta tekstowa.**

```
┌─ Dev Tools ─────────────────┐
│ [API][Data][Kontenery][Sec] │
├─────────────────────────────┤
│ POST ▼  {{host}}/v1/orders  │
│ [ Nagłówki ▼ ]              │
│ ┌─ Treść (JSON) ──────────┐ │
│ │ { "kwota": -10 }        │ │
│ └─────────────────────────┘ │
│           [ Wyślij ▶ ]      │
├─────────────────────────────┤
│ 400 Bad Request · 84 ms     │
│ { "blad": "kwota ujemna" }  │
│              [ Operacje ▼ ] │
└─────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Zakładki narzędzi | Pasek `.dn-zakladki` | Przełączanie między API Client, Data Console, Containers, Dependencies & Security | Mały pasek zakładek | aktywna · w tle · ukryta (narzędzie niewłączone w konfiguracji) | Kliknięcie przełącza zawartość kolumny | Nagłówek kolumny | 2 | Widoczne w obrębie otwartej kolumny Dev Tools; zakładka niewłączona ustawieniem konfiguracyjnym pozostaje w warstwie 4 |
| Pasek zapytania | Selektor metody i pole adresu | Definicja zapytania HTTP | Średni rząd: mały selektor i szerokie pole `.dn-pole-kontrolka` | domyślny · fokus · wysyłanie | Enter lub „Wyślij” wykonuje zapytanie i wypełnia sekcję odpowiedzi | API Client | 1 | Widoczny bez interakcji w obrębie zakładki API Client |
| Podgląd odpowiedzi | Blok składanego JSON z nagłówkami | Prezentacja treści, statusu i czasu odpowiedzi | Duży blok przewijany, krój mono | pusty · wypełniony · błąd transportu | Rozwijanie węzłów JSON; kopiowanie fragmentu | API Client | 1 | Widoczny bez interakcji w obrębie zakładki API Client |
| Drzewo schematu bazy | Drzewo baz, tabel, kolumn | Nawigacja po strukturze połączenia | Wąskie drzewo w kolumnie | zwinięte · rozwinięte · brak połączenia | Kliknięcie tabeli otwiera jej dane w siatce wyników | Data Console | 2 | Znacznik kontekstowy `[Połączenie ▼]` w zakładce Data Console; drzewo zwija się po otwarciu tabeli |
| Siatka wyników | Tabela `.dn-tabela` z wynikami zapytania | Prezentacja i edycja danych | Duża tabela przewijana | pusta · z wynikami · edycja wiersza | Edycja komórki generuje instrukcję DML prezentowaną przed zapisem | Data Console | 1 | Widoczna bez interakcji w obrębie zakładki Data Console; edycja danych w siatce należy do warstwy 4 |
| Wiersz kontenera | Pozycja listy z ikoną stanu | Reprezentuje kontener lub obraz | Mały wiersz z etykietą i plakietką stanu | uruchomiony · zatrzymany · błąd | Menu wiersza uruchamia, zatrzymuje, restartuje lub usuwa kontener; kliknięcie otwiera logi na żywo | Containers | 1 | Widoczny bez interakcji w obrębie zakładki Containers; operacje wiersza w menu ⋮ warstwy 3 |
| Pozycja podatności | Wiersz z identyfikatorem CVE i wagą | Zgłoszenie znanej podatności zależności | Mały wiersz z plakietką wagi | informacja · ostrzeżenie · krytyczna | Kliknięcie pokazuje opis, wersję naprawiającą i ścieżkę zależności | Dependencies & Security | 1 | Widoczna bez interakcji w obrębie zakładki Dependencies & Security |

---

## 4. Przepływy pracy

### 4.1. Przepływ podstawowy — edycja ze wsparciem AI i zatwierdzenie zmian

```
 Project Tree          Code Editor              Chat Window            Git Panel
 ─────────────         ────────────              ───────────            ─────────
 1. otwarcie
    pliku       ───────►
 2. zaznaczenie
    fragmentu   ──────────────► operacja kontekstowa
                                  (Refaktoryzuj / Napraw)
                                        │
                                        ▼
                                 strumień odpowiedzi
                                  z fragmentem kodu
                                        │
 3. podgląd i akceptacja  ◄─────────────┘
    wstawienia do pliku
        │
        ▼
 4. przygotowanie zmiany  ─────────────────────────────────────────►  zmiany
                                                                    przygotowane
        │
        ▼
 5. generowanie opisu commitu  ────────────────────────────────────►  commit
    i zatwierdzenie                                                  zapisany
```

### 4.2. Przepływ rozszerzony — budowanie, niepowodzenie testu, poprawka

```
Build Output                          Code Editor                    Chat Window
─────────────                         ────────────                    ───────────
uruchomienie budowania
        │
        ▼
test „obliczOplate” nie przechodzi
        │
        ├── [ Otwórz w edytorze ] ──────────► kursor w linii testu
        │
        └── [ Wyjaśnij niepowodzenie ] ──────────────────────────────► kontekst błędu
                                                                         w rozmowie
                                                                              │
                                                                              ▼
                                                                     sugestia poprawki
                                                                              │
                                        ◄─────────────────────────────────────┘
                                  [ Wstaw do edytora ]
        │
        ▼
ponowne uruchomienie budowania (warstwa wykonawcza: Terminal,
gdy powiązanie skonfigurowane, albo bezpośrednio z Build Output)
        │
        ▼
status: sukces
```

### 4.3. Przepływ pętli wykonawczej — zadanie wielokrokowe

```
 Chat Window                Execution Loop Window                 Okna robocze
 ───────────                ─────────────────────                 ────────────
 zlecenie Użytkownika
 („wydziel warstwę
  walidacji opłat”)  ─────► Koordynator dekomponuje
                            zlecenie na zadania
                                     │
                                     ▼
                            przydział zadania
                            Wykonawcy          ──────────────►  edycja plików
                                     │                           (Code Editor)
                                     ▼
                            raport wykonania    ◄─────────────  różnica
                                     │                           wielo-plikowa
                                     ▼
                            kontrola jakości:   ──────────────►  uruchomienie
                            lint · testy ·                       budowania i testów
                            przegląd różnicy    ◄─────────────  (Build Output)
                                     │
                        ┌────────────┴────────────┐
                        ▼                         ▼
                 wynik negatywny            wynik pozytywny
                 ponowienie zadania          przygotowanie commitu
                        │                         (Git Panel)
                        └────────► kolejny obieg pętli
                                     │
                                     ▼
 potwierdzenie      ◄──────── zlecenie zakończone
 Użytkownika
```

### 4.4. Przepływ rozszerzony — debugowanie zatrzymanego wykonania

```
Run & Debug                         Code Editor                    Chat Window
────────────                        ────────────                    ───────────
ustawienie punktu przerwania ◄──────  margines linii
        │
        ▼
uruchomienie konfiguracji
debugowania (adapter DAP)
        │
        ▼
zatrzymanie w punkcie przerwania
  zmienne · stos · wątki
        │
        ├── konsola debugowania (REPL): wyrażenie w kontekście ramki
        │
        └── [ Wyjaśnij stan ] ───────────────────────────────────────► opis przyczyny
                                                                       i sugestia
                                                                       poprawki
                                                                       │
                                        ◄──────────────────────────────┘
                                  [ Wstaw do edytora ]
        │
        ▼
kontynuacja wykonania · zakończenie sesji debugowania
```

### 4.5. Przepływ rozszerzony — repozytorium jako komponent budowy produktu w Apps

```
Developer › Code Editor, Git Panel        Apps › Frontend Workspace / Backend Workspace
──────────────────────────────────        ────────────────────────────────────────────
 praca nad wydzielonym modułem
 repozytorium — usługą warstwy serwerowej
        │
        │  moduł Developer stanowi jeden z komponentów
        │  wykorzystywanych przy budowie kompletnego produktu
        ▼
 zatwierdzone zmiany w gałęzi                 integracja modułu w ramach
 funkcyjnej                    ─────────────► ustrukturyzowanego procesu
                                               budowy produktu (Product Builder)
```

### 4.6. Przepływ zbiorczy — równoległa praca nad wieloma repozytoriami

```
Karta sesji A · Developer        Karta sesji B · Developer        Karta sesji C · Developer
  repozytorium: api                repozytorium: frontend            repozytorium: infrastruktura
  (kontekst odrębny)                (kontekst odrębny)                (kontekst odrębny)

           stan wyjściowy: pełna izolacja kontekstu między kartami (rozdz. 5.3)
           współdzielona pamięć projektu obowiązuje po skonfigurowaniu
                        z okna konfiguracji punktów izolacji — przy
                        pracy nad wspólnym kontraktem API
```

---

## 5. Stany, dane i powiązania

### 5.1. Model stanów sesji modułu

```
                    ┌────────────────┐
                    │  Brak repo     │  (karta sesji bez powiązanego katalogu)
                    └───────┬────────┘
                            │ powiązanie repozytorium z kartą sesji
                            ▼
                    ┌────────────────┐
        ┌──────────►│  Przegląd      │◄───────────────┐
        │           │  (Project Tree │                │
        │           │  / Code Editor)│                │
        │           └───────┬────────┘                │
        │                   │ zlecenie operacji AI    │ odrzucenie sugestii
        │                   ▼                         │
        │           ┌────────────────┐                │
        │           │  Pętla         │  (Execution    │
        │           │  wykonawcza    │   Loop Window) │
        │           └───────┬────────┘                │
        │                   │ wynik gotowy            │
        │                   ▼                         │
        │           ┌────────────────┐                │
        │           │  Podgląd zmiany│────────────────┘
        │           │  w edytorze    │
        │           └───────┬────────┘
        │                   │ akceptacja i zapis
        │                   ▼
        │           ┌────────────────┐
        │           │  Zmiana        │──► Git Panel
        │           │ niezatwierdzona│
        │           └───────┬────────┘
        │                   │ przygotowanie i commit
        │                   ▼
        │           ┌────────────────┐
        └───────────│  Commit        │──► Build Output (budowanie uruchamiane
                    │  zapisany      │     automatycznie, gdy ustawienie
                    └────────────────┘     `developer.build.auto_after_commit`
                                            jest włączone)
```

Stan procesu sesji po stronie serwera jest trwały niezależnie od stanu połączenia klienta — rozłączenie nie przerywa trwającego budowania, sesji debugowania ani pętli wykonawczej i nie zamyka otwartych plików; po ponownym połączeniu Code Editor, Git Panel i Execution Loop Window przywracają repozytorium i przebieg w stanie, w jakim zostały pozostawione.

### 5.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Developer |
|---|---|
| `karta_sesji`, `sesja` | Nośnik kontekstu bieżącej pracy nad repozytorium |
| `proces_sesji` | Proces serwera obsługujący sesję; pole `katalog_roboczy` odpowiada lokalizacji repozytorium wykorzystywanego przez Code Editor, Project Tree i Git Panel; kanał WebSocket przenosi log budowania i strumień pętli wykonawczej |
| `wiadomosc`, `zalacznik_wiadomosci` | Historia Chat Window, w tym sugestie kodu i wyjaśnienia błędów budowania |
| `zadanie` | Uruchomienie budowania, zestawu testów lub pojedynczego kroku pętli wykonawczej jako jednostka pracy; archiwum przebiegów budowania |
| `kanal_modelu` | Model bazowy przypisany karcie sesji, generujący kod, refaktoryzacje, opisy commitów i przegląd różnicy |
| `ustawienie` | Parametry modułu podlegające zasadzie „brak ustawienia = wartość domyślna” — klucze zestawione w rozdz. 8 |
| `profil_izolacji`, `regula_izolacji_technicznej` | Konfiguracja izolacji technicznej procesu repozytorium (rozdz. 5.3), w szczególności zakresu odczytu i zapisu plików oraz kont i tokenów dostępu do zdalnego repozytorium, rejestrów kontenerów i baz danych |
| `powiazanie_komponentu` | Jawne powiązania Developer ◄──► Terminal, Developer ◄──► Diagnostics, Developer ───► Apps |

Stan repozytorium — zmiany, historia commitów, gałęzie — odczytywany jest na żywo z systemu kontroli wersji w katalogu roboczym procesu sesji, a nie duplikowany jako odrębna encja w warstwie danych platformy; Git Panel jest oknem prezentującym ten stan, nie jego właścicielem.

### 5.3. Izolacja i konfigurowalność — punkty właściwe modułowi Developer

| Zakres izolacji | Stan wyjściowy | Zastosowanie w module Developer | Typowy powód włączenia |
|---|---|---|---|
| Katalog roboczy sesji | Wyłączony (współdzielony) | Lokalizacja repozytorium widocznego w Project Tree i Code Editor | Odrębna kopia robocza repozytorium dla karty prowadzącej eksperymentalną refaktoryzację |
| Zakres odczytu i zapisu plików | Wyłączony (pełny dostęp) | Uprawnienia edytora i operacji AI do plików poza katalogiem repozytorium | Sesja pracująca wyłącznie na wskazanym podkatalogu dużego monorepozytorium |
| Konto i token per sesja | Wyłączony (współdzielone dane dostępowe) | Dane uwierzytelniające do zdalnego repozytorium (Git Panel), rejestrów kontenerów (Containers) i baz danych (Data Console) | Odrębny token dostępu dla karty pracującej na repozytorium klienta |
| Dostęp sieciowy | Wyłączony (współdzielony) | Połączenia wychodzące przy pobieraniu i wysyłaniu zmian, instalacji zależności, wywołaniach klienta API i skanach podatności | Ograniczenie sesji do wyłącznie wskazanych adresów zdalnego repozytorium |
| Model procesu | Wyłączony (współdzielony) | Instancja procesu wykonawczego modelu obsługującego operacje kontekstowe AI i pętlę wykonawczą | Odrębny proces modelu dla karty o wysokiej intensywności generowania kodu |
| Katalog danych i konfiguracji modelu | Wyłączony (współdzielony) | Dane pomocnicze kanału modelu wykorzystywanego w Chat Window | Odrębna konfiguracja modelu przy pracy nad kodem objętym poufnością umowną |

Historia i pamięć Chat Window pozostają domyślnie odrębne per karta sesji (izolacja kontekstu — [Koncepcja platformy](../architektura/koncepcja-platformy.md#6-konfigurowalność-zależności-i-punktów-izolacji), rozdz. 6), niezależnie od powyższych zakresów technicznych; współdzielenie pamięci między kartami pracującymi nad powiązanymi repozytoriami jest osobną, jawną decyzją konfigurowaną na poziomie zasięgu „projekt” lub „para modułów”. Zgodnie z zasadą nadrzędną platformy żaden z powyższych punktów nie jest wymuszony — moduł działa w pełni funkcjonalnie bez jakiejkolwiek konfiguracji izolacji.

### 5.4. Powiązania z innymi modułami

```
                     ┌────────────────────────────────┐
    Terminal  ◄──────┤            DEVELOPER           │
   (warstwa          │ generowanie, refaktoryzacja,   │
    wykonawcza       │ analiza architektury,          │
    poleceń Git      │ dokumentacja, testowanie,      │
    i budowania)     │ debugowanie, integracje        │
                     └────────┬────────────┬──────────┘
                              │            │
                              │            └─────────────────────────┐
                              ▼                                      ▼
                        Diagnostics                                Apps
                 (analiza błędów wykrytych              (komponent wykorzystywany
                  w toku pracy)                          przy budowie kompletnych
                                                           produktów cyfrowych)
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Terminal | Developer korzysta z Terminala jako warstwy wykonawczej dla poleceń budowania, instalacji zależności, uruchamiania, serwerów języka i powłoki w kontenerze | Konfiguracyjne | Build Output / Git Panel / Dev Tools → Terminal Tabs |
| Diagnostics | Developer korzysta z modułu Diagnostics przy analizie błędów wykrytych w toku pracy | Konfiguracyjne | Build Output → Errors Panel |
| Apps | Developer stanowi jeden z komponentów wykorzystywanych przy budowie kompletnych produktów cyfrowych | Konfiguracyjne | Code Editor / Git Panel → Frontend Workspace / Backend Workspace |

---

## 6. Scenariusze użycia

**Scenariusz 1 — refaktoryzacja z przeglądem różnicy przed zatwierdzeniem.**
Deweloper zaznacza funkcję w Code Editor budzącą wątpliwości co do czytelności, z paska pływającego wybiera „Refaktoryzuj”. Wynik pojawia się w Chat Window wraz z blokiem kodu; deweloper wstawia go do pliku, weryfikuje margines zmian, po czym w Git Panel przegląda pełny widok różnicowy przed przygotowaniem commitu.

**Scenariusz 2 — niepowodzenie testu i poprawka prowadzona przez Wykonawcę.**
Uruchomione w Build Output budowanie kończy się niepowodzeniem jednego testu. Deweloper klika „Otwórz w edytorze”, następnie „Wyjaśnij niepowodzenie” — Chat Window otrzymuje pełny kontekst błędu i wskazuje poprawkę. Po wstawieniu poprawki deweloper uruchamia budowanie ponownie bezpośrednio z Build Output.

**Scenariusz 3 — generowanie opisu commitu i zatwierdzenie zmian.**
Po serii drobnych poprawek w kilku plikach deweloper przygotowuje zmiany w Git Panel, zaznaczając wybrane pliki. Zamiast ręcznie formułować opis, wybiera „generuj AI” — model analizuje rzeczywistą różnicę i zwraca zwięzły opis zmian, który deweloper koryguje i zatwierdza.

**Scenariusz 4 — rozwiązywanie konfliktu scalania.**
Po próbie scalenia gałęzi funkcyjnej z gałęzią główną Git Panel sygnalizuje konflikt w dwóch plikach. Deweloper otwiera widok trójstronny, dla jednego pliku przyjmuje wersję przychodzącą, dla drugiego edytuje ręcznie łącząc oba fragmenty, po czym oznacza konflikt jako rozwiązany i kończy scalenie.

**Scenariusz 5 — zadanie wielokrokowe prowadzone w pętli wykonawczej.**
Deweloper zleca w Chat Window wydzielenie warstwy walidacji opłat. Koordynator dekomponuje zlecenie na siedem zadań widocznych w Execution Loop Window. Po edycjach plików pętla uruchamia testy; jedno zadanie kończy się niepowodzeniem i zostaje ponowione z uwzględnieniem komunikatu analizy statycznej. Deweloper śledzi liczniki przebiegu, akceptuje zbiorczą różnicę wielo-plikową, po czym pętla przygotowuje commit w Git Panel.

**Scenariusz 6 — diagnoza błędu w toku debugowania.**
Deweloper ustawia warunkowy punkt przerwania na marginesie Code Editor i uruchamia konfigurację debugowania. Po zatrzymaniu przegląda zmienne i stos wywołań, wykonuje wyrażenie kontrolne w konsoli debugowania, a przyciskiem „Wyjaśnij stan” otrzymuje w Chat Window opis przyczyny zatrzymania wraz z poprawką wstawianą jako różnica.

**Scenariusz 7 — weryfikacja kontraktu API i danych.**
Deweloper importuje kontrakt OpenAPI do API Client, wykonuje zapytanie testowe i otrzymuje odpowiedź niezgodną z oczekiwaniem. W Data Console sprawdza zawartość tabeli powiązanej z żądaniem, koryguje migrację schematu, po czym ponawia zapytanie i generuje z niego test integracyjny zapisywany w repozytorium.

**Scenariusz 8 — praca nad wydzielonym modułem repozytorium jako częścią większego produktu.**
Zespół budujący aplikację w module Apps deleguje pracę nad usługą warstwy serwerowej do karty sesji modułu Developer. Po zatwierdzeniu zmian w gałęzi funkcyjnej rezultat trafia — zgodnie z powiązaniem skonfigurowanym w oknie konfiguracji — do Backend Workspace modułu Apps jako gotowy komponent integrowanego produktu.

---

## 7. Katalog funkcji i narzędzi

Każda pozycja: nazwa funkcji, opis działania oraz zależności techniczne (biblioteki, formaty, protokoły, integracje). Biblioteki warstwy serwerowej podano w wariancie Go (rdzeń platformy), warstwa kliencka edytora w wariancie webowym.

### 7.1. Rdzeń edytora kodu (Code Editor)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Silnik edytora wieloplikowy** | Karty, podział widoku (split) na kolejne kolumny, wielokursorowość, zaznaczanie kolumnowe, zwijanie bloków, minimapa | Monaco Editor lub CodeMirror 6 (warstwa kliencka); model dokumentu przez WebSocket |
| **Podświetlanie składni tree-sitter** | Inkrementalne, składniowo poprawne podświetlanie i zwijanie dla 40+ języków | `tree-sitter` + gramatyki (go, ts, js, python, rust, java, c/cpp, sql, yaml, json…); bindingi `smacker/go-tree-sitter` |
| **Serwer języka (LSP)** | Autouzupełnianie sygnatur, podgląd typu przy wskazaniu kursorem, przejdź do definicji i implementacji, znajdź wystąpienia, diagnostyka na żywo, akcje kodu | Protokół **LSP**; serwery: `gopls`, `typescript-language-server`, `pyright`, `rust-analyzer`, `clangd`, `jdtls`; most LSP↔WebSocket |
| **Refaktoryzacje semantyczne** | Zmiana nazwy symbolu globalnie, wyodrębnienie funkcji lub zmiennej, wstawienie w miejscu, przeniesienie symbolu, organizacja importów | LSP `rename`/`codeAction`; podgląd zmian jako różnica wielo-plikowa |
| **Formatowanie i styl** | Formatowanie dokumentu lub zaznaczenia wg reguł repozytorium, wybór wcięć, format przy zapisie | `gofmt`/`goimports`, `prettier`, `black`, `rustfmt`; konfiguracja z `.editorconfig` |
| **Linter i analiza statyczna w linii** | Podkreślenia błędów i ostrzeżeń na żywo, szybka poprawka, reguły per repozytorium | `golangci-lint`, `eslint`, `ruff`/`pylint`, `staticcheck`; format wyjścia SARIF |
| **Nawigacja klawiaturowa** | Paleta poleceń, szybkie otwarcie pliku (dopasowanie rozmyte), przejdź do symbolu, ścieżka symbolu w pasku okruchów | dopasowanie rozmyte `sahilm/fuzzy`; indeks symboli z LSP `documentSymbol` |
| **Wyszukiwanie i zamiana** | W pliku (wyrażenia regularne), globalnie w repozytorium z podglądem i zamianą masową, filtry glob i reguły ignorowania | silnik `ripgrep`; RE2 (Go `regexp`); poszanowanie `.gitignore` |
| **Snippet manager** | Własne, parametryzowane fragmenty wielokrotnego użytku, wyzwalane prefiksem | format snippetów LSP (placeholdery `$1`); przechowywanie w `ustawienie` (zasięg projekt lub globalny) |
| **Adnotacja autora w linii** | Kto i kiedy zmienił linię, przejście do commitu | `go-git` log per linia; powiązanie z Git Panel |
| **Margines zmian** | Znaczniki dodano/zmieniono/usunięto względem HEAD, cofnięcie zmiany linii | różnica katalog roboczy↔HEAD: `sergi/go-diff` / `hexops/gotextdiff` |
| **Podgląd renderowalny na żywo** | Podgląd Markdown, HTML, Mermaid i SVG w kolumnie obok źródła | renderer MD (`goldmark`), Mermaid (klient), sanityzacja HTML |
| **Diff wewnętrzny edytora** | Porównanie dwóch dowolnych plików lub wersji niezależnie od Git | algorytm Myers (`go-diff`); tryb inline i dwukolumnowy |

### 7.2. Nawigacja po projekcie (Project Tree)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Drzewo plików z ikonami typów** | Struktura katalogów, rozpoznanie typu, otwieranie, przeciąganie i upuszczanie, operacje plikowe | odczyt systemu plików procesu sesji; `fsnotify` do zmian na żywo |
| **Status Git per pozycja** | Plakietki M/N/U/C i ignorowany na plikach i folderach | `go-git` status; agregacja na foldery |
| **Widok „tylko zmiany”** | Przełącznik: pełne drzewo ↔ lista zmienionych plików | wspólne źródło ze statusem Git |
| **Szybkie otwarcie (dopasowanie rozmyte)** | Otwarcie pliku po fragmencie nazwy bez rozwijania drzewa | indeks plików + `sahilm/fuzzy` |
| **Operacje plikowe zbiorcze** | Nowy, zmień nazwę, usuń, duplikuj, wytnij, kopiuj, wklej; wielokrotny wybór jako zestaw kontekstu AI | operacje systemu plików z „Cofnij” w powiadomieniu (bez modali-bramek) |
| **Etykiety własne i zakładki plików** | Oznaczanie plików do przeglądu, szybki powrót | przechowywanie w `ustawienie` zasięgu projekt |
| **Podgląd binarny i obrazów** | Podgląd obrazów, czcionek oraz plików binarnych w widoku szesnastkowym | dekodery obrazów; widok szesnastkowy dla binariów |

### 7.3. Kontrola wersji (Git Panel)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Status i przygotowanie zmian** | Zmiany niezatwierdzone, przygotowane i nieśledzone; przygotowanie i wycofanie przygotowania pliku, zaznaczenia lub wszystkiego, częściowe przygotowanie fragmentów różnicy | `go-git` lub `git2go` (libgit2) dla operacji na poziomie hunków |
| **Widok różnicowy** | Inline i dwukolumnowy, per plik i per hunk, porównanie słowo po słowie | `go-diff`; podświetlanie składni w różnicy (tree-sitter) |
| **Commit z opisem AI** | Pole treści z licznikiem, generowanie opisu z rzeczywistej różnicy, amend, revert, podpis | `kanal_modelu` (opis), podpis GPG/SSH, `git2go` |
| **Historia i wykres gałęzi** | Log z autorem, datą i skrótem, filtry, wizualny wykres gałęzi, przejście do plików commitu | `go-git` log i rev-walk; render grafu (klient) |
| **Operacje gałęziowe** | Utwórz, checkout, merge, rebase, cherry-pick, usuń; porównanie gałęzi i commitów | `git2go` (merge, rebase), `go-git` (odczyt) |
| **Zdalne repozytorium** | Fetch, pull, push, zarządzanie remotes, force-push jawnie sygnalizowany jako nieodwracalny | transport SSH/HTTPS; poświadczenia z `profil_izolacji` (token per sesja) |
| **Stash** | Odłożenie i przywrócenie zmian roboczych, lista stashy | `git2go` stash API |
| **Rozwiązywanie konfliktów 3-way** | Widok bieżąca/przychodząca/wynik, wybór fragmentu lub edycja ręczna | `git2go` merge; renderer trójpanelowy |
| **Przegląd kodu na różnicy (AI)** | Uwagi przypięte do konkretnych linii różnicy przed scaleniem | `kanal_modelu`; format komentarza per linia |
| **Changelog i nota wydania (AI)** | Podsumowanie zakresu commitów w notę wydania | `kanal_modelu`; konwencja Conventional Commits |
| **Integracja z hostingiem** | Żądania scalenia (PR/MR), przypisania, statusy CI, komentarze — bez opuszczania okna | integracja wg rozdz. 5.7 Modelu konfiguracji: API GitHub, GitLab, Bitbucket |

### 7.4. Uruchamianie i budowanie (Build Output)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Zadania budowania i uruchomień** | Definicje zadań (budowanie, uruchomienie, test, obserwacja zmian), profile deweloperski i produkcyjny, wybór platformy docelowej | `zadanie` (model danych); wykonanie przez **Terminal** (warstwa wykonawcza) |
| **Log na żywo (WebSocket)** | Strumień budowania, plakietka statusu, filtry błędów, ostrzeżeń i informacji, wyszukiwanie w logu | kanał WebSocket `proces_sesji`; parser komunikatów kompilatora |
| **Przejście błąd → linia kodu** | Kliknięcie komunikatu otwiera plik w miejscu błędu | parser ścieżek `plik:linia:kolumna`; mapowanie do Code Editor |
| **Wynik testów i pokrycie** | Przeszedł/nie przeszedł/pominięty, pasek postępu, pokrycie procentowe z rozbiciem na pliki i podświetleniem linii | format `go test -json`, JUnit XML, `lcov`/`cobertura`; nakładka pokrycia w edytorze |
| **Historia przebiegów** | Lista poprzednich budowań z czasem, wynikiem i skrótem commitu; ponowne uruchomienie | `zadanie` archiwum; powiązanie ze skrótem Git |
| **Tryb obserwacji zmian** | Automatyczne przebudowanie i przeładowanie przy zmianie plików | `fsnotify`; sygnał do warstwy wykonawczej |
| **Eksport raportu** | Raport budowania i testów do pliku (Markdown, JUnit, HTML) | serializatory; format JUnit XML i HTML |
| **Przekazanie błędu do Diagnostics** | Zgłoszenie niepowodzenia jako materiału źródłowego | `powiazanie_komponentu` Developer→Diagnostics |

### 7.5. Debugowanie (Run & Debug)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Debugger krokowy (DAP)** | Uruchomienie, zatrzymanie, wejście, przejście, wyjście, kontynuacja, restart | Protokół **DAP**; adaptery: `delve` (Go), `debugpy` (Python), `js-debug` (Node), `lldb`/`gdb` |
| **Punkty przerwania** | Zwykłe, warunkowe, logujące, przerwanie na wyjątku, licznik trafień | DAP `setBreakpoints`; margines Code Editor |
| **Podgląd zmiennych i wyrażeń** | Zakresy lokalne i globalne, wyrażenia obserwowane, edycja wartości w trakcie zatrzymania | DAP `scopes`/`variables`/`evaluate` |
| **Stos wywołań i wątki** | Nawigacja po ramkach, przełączanie wątków i goroutines | DAP `stackTrace`/`threads` |
| **Konsola debugowania (REPL)** | Wykonywanie wyrażeń w kontekście zatrzymania | DAP `evaluate` w kontekście ramki |
| **Wyjaśnij stan (AI)** | Opis przyczyny zatrzymania na podstawie stosu i zmiennych wraz z poprawką | `kanal_modelu` + migawka stanu DAP |

### 7.6. Klient API i HTTP (API Client)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Zapytania HTTP/REST** | Metoda, adres, nagłówki, treść (JSON, formularz, multipart), podgląd odpowiedzi, czas, status | klient HTTP Go (`net/http`); podgląd JSON składany |
| **Kolekcje i środowiska** | Grupowanie zapytań, zmienne środowiskowe `{{placeholder}}`, przełączanie środowisk | format kolekcji (własny, `.http`, import OpenAPI); zmienne z `ustawienie` |
| **Import i eksport OpenAPI** | Generowanie zapytań z kontraktu OpenAPI/Swagger, walidacja odpowiedzi ze schematem | `getkin/kin-openapi`; formaty OpenAPI 3.x i Swagger 2 |
| **GraphQL, gRPC, WebSocket** | Zapytania GraphQL z introspekcją, wywołania gRPC z `.proto`, testy WS | `graphql-go`, `grpcurl`/reflection, klient WS |
| **Generowanie klienta i testów (AI)** | Tworzenie kodu klienta lub testów integracyjnych z definicji zapytania | `kanal_modelu`; szablony per język |
| **Plik `.http` w repozytorium** | Zapytania wersjonowane jako plik repozytorium, uruchamiane z edytora | parser `.http`; wykonanie z Code Editor |

### 7.7. Bazy danych i dane (Data Console)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Przeglądarka połączeń i schematu** | Drzewo baz, tabel, kolumn i indeksów, wiele połączeń | `database/sql` + sterowniki `pgx`, `go-sql-driver/mysql`, `mattn/go-sqlite3` |
| **Konsola SQL** | Edytor SQL z autouzupełnianiem schematu, wykonanie, siatka wyników, eksport CSV i JSON | LSP SQL (`sqls`); eksport `encoding/csv` |
| **Edycja danych w siatce** | Edycja wierszy w miejscu, wstawianie i usuwanie z generowaniem SQL | generator DML; transakcja z podglądem przed zapisem |
| **Migracje** | Podgląd, tworzenie i uruchamianie migracji schematu | `golang-migrate`; formaty migracji `.sql` i `.go` |
| **Zapytanie z języka naturalnego (AI)** | Generowanie SQL z opisu na podstawie schematu połączenia | `kanal_modelu` + wprowadzony schemat jako kontekst |
| **Podgląd planu zapytania** | EXPLAIN/ANALYZE z wizualizacją kosztów | parser planu per silnik (Postgres, MySQL) |

### 7.8. Kontenery i środowiska wykonawcze (Containers)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Lista kontenerów i obrazów** | Podgląd, uruchomienie, zatrzymanie, restart, usunięcie, logi na żywo, statystyki | Docker SDK Go (`github.com/docker/docker/client`); API Docker i Podman |
| **Budowanie obrazu** | Budowanie z `Dockerfile`, tagowanie, wypchnięcie do rejestru | Docker SDK; format `Dockerfile`, `.dockerignore` |
| **Compose i wiele usług** | Uruchomienie stosu `docker-compose.yml`, podgląd usług | parser `docker-compose.yml`; API Compose |
| **Powłoka w kontenerze** | Wejście do kontenera przez warstwę Terminal | Docker exec API → **Terminal** |
| **Kontener deweloperski** | Otwarcie repozytorium w kontenerze wg `devcontainer.json` | specyfikacja `devcontainer.json`; Docker SDK |

### 7.9. Zależności, bezpieczeństwo, jakość

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Przeglądarka zależności** | Drzewo zależności, wersje, pakiety przestarzałe, aktualizacja | parsery manifestów: `go.mod`, `package.json`, `requirements.txt`, `Cargo.toml` |
| **Skan podatności (SCA)** | Wykrycie znanych CVE w zależnościach, wskazanie bezpiecznej wersji | `grype` / OSV API; format SBOM (CycloneDX, SPDX przez `syft`) |
| **Skan sekretów** | Wykrycie kluczy i tokenów w kodzie przed zatwierdzeniem zmian | `gitleaks`; reguły wyrażeń regularnych; hak pre-commit |
| **SAST (analiza bezpieczeństwa kodu)** | Wykrycie wzorców podatności w kodzie | `semgrep`; reguły; wyjście SARIF |
| **Analiza licencji** | Zestawienie licencji zależności, flagowanie niezgodnych | dane licencji z SBOM; polityka z `ustawienie` |
| **Mutacje i pokrycie zaawansowane** | Testy mutacyjne, próg pokrycia prezentowany jako ostrzeżenie, nie blokada | narzędzia mutacyjne per język; próg w `ustawienie` |

### 7.10. Operacje kontekstowe AI

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Generuj / Refaktoryzuj / Wyjaśnij / Udokumentuj / Napraw / Napisz test / Zoptymalizuj** | Operacje na zaznaczeniu lub pliku z paska pływającego, wynik jako różnica do akceptacji | `kanal_modelu`; kontekst z Project Tree, Code Editor i Git |
| **Konwersja język i biblioteka** | Przepisanie fragmentu między językami lub frameworkami | `kanal_modelu`; wskazanie źródła i celu |
| **Tryb agentowy (wielokrokowy)** | Realizacja zadania obejmującego wiele plików: plan → edycje → uruchomienie testów → korekta, prowadzona w Execution Loop Window | `zadanie`; pętla Koordynator↔Wykonawca z narzędziami (edycja, testy, Git); zgodne ze **SPECYFIKACJĄ AGENTÓW** |
| **Przegląd całego repozytorium (architektura)** | Opis struktury, zależności i punktów ryzyka nieznanego repozytorium | indeks symboli + graf zależności jako kontekst |
| **Wyjaśnij błąd budowania, testu lub debugowania** | Kontekst błędu → poprawka wstawiana jako różnica | powiązanie z Build Output i Run & Debug (DAP) |
| **Wskaźnik i zarządzanie kontekstem** | Licznik długości kontekstu i liczby włączonych plików, ręczne dodanie i usunięcie | metryka tokenów `kanal_modelu` |

### 7.11. Zestawienie zależności bibliotecznych (warstwa serwerowa, Go)

| Obszar | Rozwiązanie stosowane | Rozwiązanie zamienne | Uwaga |
|---|---|---|---|
| Git odczyt, status, log | `go-git` (czysty Go) | `git2go` (libgit2, cgo) | `go-git` działa bez cgo; merge, rebase i staging hunków obsługuje `git2go` |
| Diff tekstowy | `hexops/gotextdiff` | `sergi/go-diff` | oba realizują algorytm Myers; `gotextdiff` daje czytelniejszy format unified |
| Podświetlanie składni | `smacker/go-tree-sitter` | Chroma (`alecthomas/chroma`) | tree-sitter jest dokładniejszy; Chroma obsługuje języki bez gramatyki |
| Dopasowanie rozmyte | `sahilm/fuzzy` | `lithammer/fuzzysearch` | biblioteki wymienne |
| Obserwacja systemu plików | `fsnotify/fsnotify` | odpytywanie cykliczne | fsnotify jest rozwiązaniem podstawowym |
| Kontenery | Docker SDK Go | Podman REST API | wspólne API OCI |
| OpenAPI | `getkin/kin-openapi` | `pb33f/libopenapi` | walidacja 3.x |
| SBOM i podatności | `syft` + `grype` | OSV-Scanner | formaty CycloneDX i SPDX, wyjście SARIF |
| Sterowniki baz danych | `pgx`, `go-sql-driver/mysql`, `mattn/go-sqlite3` | `database/sql` generyczny | dobór per silnik |
| Debugowanie | adaptery **DAP** | natywne API (`delve`) | DAP ujednolica obsługę wielu języków |
| Serwery języka | most do serwerów LSP | wbudowane parsery | serwery działają jako procesy zewnętrzne w warstwie Terminal |

---

## 8. Punkty sterowania z okna konfiguracji

Wszystkie punkty są jawne, opatrzone objaśnieniem `[?]`, przechowywane jako `ustawienie` z kluczem ([Model danych](../architektura/model-danych.md#171-ustawienie-ustawienie), rozdz. 17.1). Klucze są jawne, wartości podlegają zasadzie „brak ustawienia = wartość domyślna”. Konfiguracja steruje widocznością i wartościami domyślnymi, nie wyłącza działania modułu; niewłączona integracja ukrywa zakładkę zamiast prezentować zablokowany element interfejsu.

### 8.1. Edytor i jakość kodu

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> **kluczami konfiguracji**, nie komendami kontraktu. Klucz wskazuje miejsce wartości w modelu
> konfiguracji; komendy kontraktu, którymi się go odczytuje i zapisuje, to `config.get` i
> `config.set` (obszar `config` w `budowa/shared/contract.json`).


| Klucz | Co Operator personalizuje | Wartość domyślna |
|---|---|---|
| `developer.editor.theme` | Motyw kolorystyczny kodu | z motywu platformy |
| `developer.editor.indentation` | Spacje lub tabulacje i szerokość wcięcia | 2 spacje / z `.editorconfig` |
| `developer.editor.word_wrap` | Zawijanie wierszy | wyłączone |
| `developer.editor.format_on_save` | Automatyczne formatowanie przy zapisie | włączone |
| `developer.lint.servers` | Aktywne narzędzia analizy statycznej dla poszczególnych języków | z konfiguracji repozytorium |
| `developer.lsp.servers` | Mapowanie język → serwer LSP | automatyczne wykrycie |
| `developer.editor.snippets` | Zestaw snippet (zasięg projekt lub globalny) | wbudowane |

### 8.2. Git i wersjonowanie

| Klucz | Co Operator personalizuje | Wartość domyślna |
|---|---|---|
| `developer.git.identity` | Nazwa i adres e-mail commitującego | z profilu Operatora |
| `developer.git.signing` | Podpis GPG/SSH commitów | wyłączone |
| `developer.git.confirm_delete` | Okno nakładkowe potwierdzenia przed usunięciem pliku | wyłączony (usunięcie od razu z „Cofnij”) |
| `developer.git.confirm_force_push` | Okno nakładkowe potwierdzenia przed wysłaniem wymuszonym | wyłączony (ostrzeżenie przy kontrolce) |
| `developer.git.commit_ai_default` | Automatyczne generowanie opisu commitu | na żądanie |
| `developer.git.hosting` | Integracja hostingu (GitHub, GitLab, Bitbucket) | brak |

### 8.3. Budowanie, uruchamianie, debugowanie

| Klucz | Co Operator personalizuje | Wartość domyślna |
|---|---|---|
| `developer.build.profile` | Profil domyślny (deweloperski lub produkcyjny) | deweloperski |
| `developer.build.watch` | Tryb obserwacji zmian i przeładowanie | wyłączony |
| `developer.test.coverage_threshold` | Próg pokrycia prezentowany jako ostrzeżenie | 0% (bez ostrzeżenia) |
| `developer.build.auto_after_commit` | Automatyczne budowanie po commicie | wyłączone |
| `developer.debug.adapters` | Mapowanie język → adapter DAP | automatyczne wykrycie |

### 8.4. Integracje deweloperskie — widoczność zakładek Dev Tools

| Klucz | Co Operator personalizuje | Wartość domyślna |
|---|---|---|
| `developer.tools.api_client` | Widoczność klienta API | ukryty |
| `developer.tools.data_console` | Widoczność konsoli bazodanowej i połączeń | ukryta |
| `developer.tools.containers` | Widoczność panelu kontenerów | ukryty |
| `developer.tools.security` | Skany SCA, SAST i sekretów oraz polityka | ukryte |
| `developer.data.connections` | Zdefiniowane połączenia baz danych (poświadczenia przez izolację) | brak |

### 8.5. AI i model

| Klucz | Co Operator personalizuje | Wartość domyślna |
|---|---|---|
| `developer.ai.model_channel` | Kanał modelu dla operacji kontekstowych | domyślny kanał sesji |
| `developer.ai.agent_mode` | Dopuszczenie zadań wielokrokowych prowadzonych w pętli wykonawczej | włączony |
| `developer.ai.auto_context_files` | Automatyczne włączanie otwartych plików do kontekstu | włączone |
| `developer.ai.diff_before_insert` | Podgląd różnicy przed wstawieniem sugestii | włączony |

### 8.6. Izolacja i zależności — powiązania jawne

Zgodnie z rozdz. 5.3 oraz oknem konfiguracji punktów izolacji Operator steruje zakresami: katalog roboczy sesji, zakres odczytu i zapisu plików, konto i token per sesja (dane dostępowe zdalnego repozytorium, rejestrów kontenerów i baz danych), dostęp sieciowy, model procesu, katalog danych modelu. Powiązania `Developer→Terminal` (warstwa wykonawcza), `Developer→Diagnostics` (błędy) i `Developer→Apps` (komponent) konfigurowane są jako `powiazanie_komponentu`.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + S` | Zapis pliku | Code Editor |
| `Ctrl/Cmd + P` | Szybkie otwarcie pliku | Code Editor, Project Tree |
| `Ctrl/Cmd + K` | Paleta poleceń | Code Editor |
| `Ctrl/Cmd + F` / `Ctrl/Cmd + H` | Znajdź / Zamień w pliku | Code Editor |
| `Ctrl/Cmd + Shift + F` | Grep — znajdź w plikach | Code Editor |
| `F12` / `Shift + F12` | Przejdź do definicji / znajdź wystąpienia | Code Editor |
| `Alt/Option + Shift + F` | Formatuj dokument | Code Editor |
| `Ctrl/Cmd + D` | Zaznacz kolejne wystąpienie (wielokursorowość) | Code Editor |
| `Ctrl/Cmd + /` | Komentarz / odkomentowanie linii | Code Editor |
| `F5` / `Shift + F5` | Uruchom / zakończ sesję debugowania | Run & Debug |
| `F9` | Ustaw lub usuń punkt przerwania | Code Editor, Run & Debug |
| `F10` / `F11` | Przejdź / wejdź do funkcji | Run & Debug |
| `Ctrl/Cmd + Enter` | Wysłanie polecenia do Wykonawcy | Chat Window |
| `Ctrl/Cmd + Shift + L` | Otwarcie kolumny pętli wykonawczej | Execution Loop Window |
| Ikona `kod` | Developer, Code Editor, model procesu (izolacja techniczna) | Code Editor, Chat Window |
| Ikona `plik` / `folder` | Węzły Project Tree | Project Tree |
| Ikona `strzalka-lewo` / `strzalka-prawo` | Nawigacja różnicy, historia | Git Panel, Code Editor |
| Ikona `uruchom` / `zatrzymaj` | Start i zatrzymanie budowania lub debugowania | Build Output i Run & Debug |
| Ikona `petla` | Pętla wykonawcza, ponowienie zadania | Execution Loop Window |
| Ikona `ptaszek` / `ptaszek-kolo` | Test przeszedł, walidacja | Build Output i Run & Debug |
| Ikona `blad` | Błąd budowania, konflikt, test nieudany | Build Output i Run & Debug, Git Panel |
| Ikona `klodka` | Konto i token per sesja (izolacja techniczna) | okno konfiguracji punktów izolacji |
| Ikona `oko` | Podgląd różnicy, podgląd zmiany | Git Panel, Code Editor |

## 10. Punkty łamania i kryteria odbioru

### 10.1. Punkty łamania

Moduł Developer dzieli obszar roboczy na kolumny sąsiadujące poziomo (rozdz. 2). Poniższe progi, wspólne całej platformie, rozstrzygają zachowanie układu tych kolumn wraz ze zmianą szerokości okna aplikacji.

| Żeton | Próg szerokości | Zachowanie układu w module Developer |
|---|---|---|
| `--dn-bp-w1` | 640 px | Telefon poziomo — widok mobilny; okna operacyjne modułu prezentowane pojedynczo, pełny ekran, nawigacja powrotna zastępuje układ kolumnowy |
| `--dn-bp-w2` | 960 px | Tablet — boczna nawigacja modułów zwija się do samych ikon; kolumny modułu Developer zachowują układ, lecz kolumny boczne (rozszerzenia warstwy 2–3) otwierają się jako nakładka zamiast stałej kolumny |
| `--dn-bp-w3` | 1280 px | Biurko — pełny kokpit; wszystkie okna operacyjne modułu Developer wymienione w rozdz. 2 dostępne jednocześnie w układzie kolumnowym opisanym w tym rozdziale |
| `--dn-bp-w4` | 1600 px | Szerokie biurko — para Chat Window i Execution Loop Window prezentowana jednocześnie obok okna wiodącego modułu, bez wzajemnego przesłaniania |

Poniżej progu `--dn-bp-w1` układ kolumnowy modułu Developer nie jest dostępny — zachowanie na urządzeniach mobilnych ustala [Mobile](../funkcje-globalne/mobile.md).

### 10.2. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Moduł jest osiągalny wyłącznie z bocznej nawigacji CodeStudio | Próba odnalezienia Developer w bocznej nawigacji TalkIn i WorkSpace — brak pozycji |
| Żaden przycisk akcji modułu nie występuje w stanie zablokowanym | Przegląd kontrolek akcji katalogu elementów rozdz. 3 pod kątem obecności `disabled`/`not-allowed` |
| Wysłanie wymuszone i usunięcie pliku pozostają domyślnie natychmiastowe, z „Cofnij” dostępnym po wykonaniu | Test: „Wyślij wymuszony” bez włączonego ustawienia `developer.git.confirm_force_push` wykonuje się od razu |
| Wszystkie 50 komend obszaru `developer` mają pokrycie w interfejsie modułu albo jawne wskazanie miejsca wywołania | Zestawienie Załącznika z katalogiem elementów rozdz. 3 |
| Paleta poleceń otwiera się skrótem `Ctrl/Cmd + K` z dowolnego miejsca edytora | Wywołanie skrótu w Code Editor |
| Test wygenerowany w odpowiedzi Wykonawcy uruchamia się w warstwie wykonawczej Terminal | Kliknięcie „Uruchom testy” pod blokiem kodu zawierającym test |
| Stany kontrolek z katalogu elementów rozdz. 3 są zaimplementowane w komplecie wymaganym rozdz. 4.2 standardu redakcyjnego | Przegląd katalogu elementów rozdz. 3 pod kątem kompletu stanów każdej kontrolki: spoczynek, wskazanie kursorem, wciśnięcie, ognisko, nieaktywny, ładowanie, pusty, błąd |

---

## Załącznik — pełny wykaz komend kontraktu modułu Developer

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `developer` — 50 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `developer.api.collection.list` | Zwraca kolekcje zapytań okna | `windowId:string` (wym)<br>`collectionId:string` (opc) | `collections:ApiCollection[]` (wym) |
| `developer.api.collection.save` | Zapisuje kolekcję zapytań wraz ze środowiskami | `windowId:string` (wym)<br>`collectionId:string` (opc)<br>`name:string` (wym)<br>`requests:json` (wym)<br>`environments:json` (opc) | `collection:ApiCollection` (wym) |
| `developer.api.openapi.import` | Buduje kolekcję zapytań z kontraktu OpenAPI | `windowId:string` (wym)<br>`path:string` (opc)<br>`url:string` (opc)<br>`collectionName:string` (opc) | `collection:ApiCollection` (wym)<br>`requestCount:int` (wym) |
| `developer.api.request` | Wykonuje zapytanie HTTP w sieci okna modułu | `windowId:string` (wym)<br>`method:string` (wym)<br>`url:string` (wym)<br>`headers:json` (opc)<br>`body:string` (opc)<br>`bodyKind:string` (opc)<br>`environmentId:string` (opc)<br>`timeoutMs:int` (opc) | `response:ApiResponse` (wym) |
| `developer.breakpoint.set` | Zakłada albo zdejmuje punkt przerwania na marginesie Code Editora | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`kind:BreakpointKind` (opc)<br>`condition:string` (opc)<br>`hitCondition:string` (opc)<br>`logMessage:string` (opc)<br>`remove:bool` (opc) | `breakpoints:Breakpoint[]` (wym) |
| `developer.build.list` | Zwraca przebiegi budowania z dziennika okna | `windowId:string` (wym)<br>`status:BuildStatus` (opc)<br>`limit:int` (opc) | `builds:DeveloperBuild[]` (wym)<br>`total:int` (opc) |
| `developer.build.log.get` | Rozwija odnośnik logu przebiegu budowania | `buildId:string` (wym)<br>`tail:int` (opc)<br>`fromLine:int` (opc) | `lines:string[]` (wym)<br>`truncated:bool` (wym)<br>`truncatedFrom:int` (opc) |
| `developer.build.run` | Uruchamia albo przerywa budowanie; log idzie na żywo do Build Output | `windowId:string` (wym)<br>`task:string` (wym)<br>`arguments:string[]` (opc)<br>`stop:bool` (opc) | `build:DeveloperBuild` (wym) |
| `developer.compose.up` | Podnosi albo zatrzymuje stos usług opisany plikiem compose | `windowId:string` (wym)<br>`file:string` (wym)<br>`services:string[]` (opc)<br>`down:bool` (opc) | `services:ContainerInfo[]` (wym)<br>`output:string` (opc) |
| `developer.container.action` | Wykonuje czynność cyklu życia kontenera | `containerId:string` (wym)<br>`action:ContainerActionKind` (wym)<br>`tail:int` (opc) | `container:ContainerInfo` (wym)<br>`output:string` (opc) |
| `developer.container.list` | Zwraca kontenery i obrazy silnika kontenerów | `windowId:string` (wym)<br>`all:bool` (opc)<br>`includeImages:bool` (opc) | `containers:ContainerInfo[]` (wym)<br>`images:ImageInfo[]` (opc)<br>`engineAvailable:bool` (wym) |
| `developer.contextual.op` | Wykonuje operację kontekstową modelu nad treścią repozytorium | `windowId:string` (wym)<br>`operation:ContextualOpKind` (wym)<br>`path:string` (opc)<br>`selection:string` (opc)<br>`contextPaths:string[]` (opc)<br>`instruction:string` (opc)<br>`targetLanguage:string` (opc) | `result:string` (wym)<br>`edits:DeveloperTextEdit[]` (opc)<br>`messageId:string` (opc) |
| `developer.coverage.get` | Zwraca pokrycie kodu testami dla przebiegu budowania | `buildId:string` (wym)<br>`path:string` (opc) | `files:DeveloperCoverage[]` (wym)<br>`percent:int` (wym)<br>`threshold:int` (opc) |
| `developer.data.connection.list` | Zwraca połączenia bazodanowe okna | `windowId:string` (wym) | `connections:DataConnection[]` (wym) |
| `developer.data.connection.set` | Zakłada albo zmienia definicję połączenia bazodanowego | `windowId:string` (wym)<br>`connectionId:string` (opc)<br>`name:string` (wym)<br>`engine:DataEngine` (wym)<br>`host:string` (opc)<br>`port:int` (opc)<br>`database:string` (wym)<br>`user:string` (opc)<br>`credentialRef:string` (opc)<br>`readOnly:bool` (opc) | `connection:DataConnection` (wym) |
| `developer.data.migration.run` | Uruchamia migracje schematu z katalogu roboczego okna | `connectionId:string` (wym)<br>`windowId:string` (wym)<br>`direction:string` (wym)<br>`target:string` (opc)<br>`dryRun:bool` (opc) | `applied:string[]` (wym)<br>`pending:string[]` (wym)<br>`output:string` (opc) |
| `developer.data.query.run` | Wykonuje zapytanie na połączeniu albo oddaje jego plan | `connectionId:string` (wym)<br>`sql:string` (wym)<br>`limit:int` (opc)<br>`explain:bool` (opc)<br>`transaction:bool` (opc) | `result:DataQueryResult` (wym) |
| `developer.data.schema.get` | Zwraca drzewo schematu połączenia w liście płaskiej | `connectionId:string` (wym)<br>`path:string` (opc)<br>`depth:int` (opc) | `nodes:DataSchemaNode[]` (wym) |
| `developer.debug.evaluate` | Wykonuje wyrażenie w kontekście zatrzymanej ramki | `sessionId:string` (wym)<br>`frameId:string` (wym)<br>`expression:string` (wym)<br>`assignTo:string` (opc) | `value:string` (wym)<br>`type:string` (opc)<br>`variablesRef:string` (opc) |
| `developer.debug.scope.get` | Zwraca stos wywołań, zakresy i zmienne zatrzymanego procesu | `sessionId:string` (wym)<br>`frameId:string` (opc)<br>`variablesRef:string` (opc) | `frames:DebugFrame[]` (wym)<br>`scopes:DebugScope[]` (wym)<br>`variables:DebugVariable[]` (wym) |
| `developer.debug.session.control` | Steruje przebiegiem sesji debugowania | `sessionId:string` (wym)<br>`step:DebugStepKind` (wym)<br>`threadId:string` (opc) | `session:DebugSession` (wym) |
| `developer.debug.session.start` | Rozpoczyna sesję debugowania pod adapterem właściwym językowi | `windowId:string` (wym)<br>`configurationId:string` (opc)<br>`program:string` (opc)<br>`arguments:string[]` (opc)<br>`adapter:string` (opc)<br>`stopOnEntry:bool` (opc) | `session:DebugSession` (wym) |
| `developer.dependency.list` | Zwraca drzewo zależności z manifestu repozytorium | `windowId:string` (wym)<br>`manifest:string` (opc)<br>`outdatedOnly:bool` (opc)<br>`depth:int` (opc) | `dependencies:DependencyNode[]` (wym)<br>`manifest:string` (wym) |
| `developer.file.create` | Zakłada plik albo katalog w katalogu roboczym okna | `windowId:string` (wym)<br>`path:string` (wym)<br>`kind:TreeNodeKind` (wym)<br>`content:string` (opc) | `node:DeveloperTreeNode` (wym) |
| `developer.file.delete` | Usuwa wskazane węzły drzewa projektu | `windowId:string` (wym)<br>`paths:string[]` (wym)<br>`recursive:bool` (opc) | `deletedPaths:string[]` (wym)<br>`undoToken:string` (opc) |
| `developer.file.move` | Przenosi węzły drzewa do wskazanego katalogu | `windowId:string` (wym)<br>`paths:string[]` (wym)<br>`targetPath:string` (wym) | `nodes:DeveloperTreeNode[]` (wym) |
| `developer.file.open` | Wczytuje plik repozytorium do Code Editor | `windowId:string` (wym)<br>`path:string` (wym) | `file:DeveloperFile` (wym) |
| `developer.file.rename` | Zmienia nazwę pliku albo katalogu bez zmiany jego miejsca | `windowId:string` (wym)<br>`path:string` (wym)<br>`newName:string` (wym) | `node:DeveloperTreeNode` (wym) |
| `developer.file.save` | Zapisuje treść pliku repozytorium | `windowId:string` (wym)<br>`path:string` (wym)<br>`content:string` (wym)<br>`createVersion:bool` (opc) | `file:DeveloperFile` (wym) |
| `developer.file.version.list` | Zwraca wersje pliku założone przy zapisie z polem `createVersion` | `windowId:string` (wym)<br>`path:string` (wym)<br>`limit:int` (opc) | `versions:DeveloperFileVersion[]` (wym)<br>`total:int` (opc) |
| `developer.file.version.restore` | Przywraca treść pliku z zapisanej wersji | `windowId:string` (wym)<br>`versionId:string` (wym)<br>`writeToDisk:bool` (opc) | `file:DeveloperFile` (wym) |
| `developer.format.run` | Formatuje dokument albo zaznaczenie według reguł stylu repozytorium | `windowId:string` (wym)<br>`path:string` (wym)<br>`content:string` (opc)<br>`startLine:int` (opc)<br>`endLine:int` (opc)<br>`writeToDisk:bool` (opc) | `file:DeveloperFile` (wym)<br>`changed:bool` (wym)<br>`formatter:string` (opc) |
| `developer.git.action` | Wykonuje czynność na repozytorium powiązanym z sesją | `windowId:string` (wym)<br>`action:GitActionKind` (wym)<br>`paths:string[]` (opc)<br>`message:string` (opc)<br>`branch:string` (opc)<br>`remote:string` (opc)<br>`force:bool` (opc) | `result:GitActionResult` (wym) |
| `developer.git.branch.list` | Zwraca gałęzie repozytorium wraz z rozbieżnością wobec gałęzi zdalnej | `windowId:string` (wym)<br>`includeRemote:bool` (opc) | `branches:GitBranch[]` (wym)<br>`current:string` (opc) |
| `developer.git.conflict.get` | Oddaje trzy wersje pliku skonfliktowanego do widoku trójstronnego | `windowId:string` (wym)<br>`path:string` (wym) | `conflict:GitConflict` (wym) |
| `developer.git.conflict.resolve` | Oznacza konflikt jako rozstrzygnięty wskazaną wersją albo treścią własną | `windowId:string` (wym)<br>`path:string` (wym)<br>`resolution:ConflictResolutionKind` (wym)<br>`content:string` (opc) | `resolved:bool` (wym)<br>`remainingPaths:string[]` (wym) |
| `developer.git.diff` | Zwraca różnice między wskazanymi stanami repozytorium | `windowId:string` (wym)<br>`paths:string[]` (opc)<br>`staged:bool` (opc)<br>`fromRef:string` (opc)<br>`toRef:string` (opc)<br>`contextLines:int` (opc) | `hunks:GitDiffHunk[]` (wym)<br>`binaryPaths:string[]` (opc) |
| `developer.git.log` | Zwraca historię zatwierdzeń repozytorium | `windowId:string` (wym)<br>`branch:string` (opc)<br>`path:string` (opc)<br>`author:string` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`limit:int` (opc) | `commits:GitCommit[]` (wym)<br>`total:int` (opc) |
| `developer.git.status` | Oddaje stan repozytorium katalogu roboczego okna bez wykonywania czynności | `windowId:string` (wym) | `status:DeveloperGitStatus` (wym) |
| `developer.grep.replace` | Zamienia trafienia wzorca w repozytorium; warstwa funkcji eksperckich | `windowId:string` (wym)<br>`pattern:string` (wym)<br>`replacement:string` (wym)<br>`regex:bool` (opc)<br>`paths:string[]` (opc)<br>`preview:bool` (opc) | `edits:DeveloperTextEdit[]` (wym)<br>`changedPaths:string[]` (wym)<br>`applied:bool` (wym) |
| `developer.grep.search` | Szuka wzorca w całym repozytorium katalogu roboczego okna | `windowId:string` (wym)<br>`pattern:string` (wym)<br>`regex:bool` (opc)<br>`caseSensitive:bool` (opc)<br>`include:string[]` (opc)<br>`exclude:string[]` (opc)<br>`limit:int` (opc) | `matches:DeveloperGrepMatch[]` (wym)<br>`total:int` (opc)<br>`truncated:bool` (wym) |
| `developer.image.build` | Buduje obraz kontenera z pliku Dockerfile repozytorium | `windowId:string` (wym)<br>`dockerfile:string` (wym)<br>`tag:string` (wym)<br>`buildArgs:json` (opc)<br>`push:bool` (opc)<br>`registry:string` (opc) | `imageId:string` (wym)<br>`logRef:string` (opc) |
| `developer.lint.get` | Zwraca zgłoszenia analizy statycznej dla wskazanych plików | `windowId:string` (wym)<br>`paths:string[]` (opc)<br>`limit:int` (opc) | `diagnostics:DeveloperDiagnostic[]` (wym)<br>`linterAvailable:bool` (wym)<br>`truncated:bool` (opc) |
| `developer.refactor.apply` | Wykonuje refaktoryzację semantyczną obejmującą wiele plików | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`column:int` (wym)<br>`kind:RefactorKind` (wym)<br>`newName:string` (opc)<br>`preview:bool` (opc) | `edits:DeveloperTextEdit[]` (wym)<br>`applied:bool` (wym)<br>`changedPaths:string[]` (opc) |
| `developer.scan.result.list` | Zwraca zgłoszenia skanu bezpieczeństwa | `scanId:string` (opc)<br>`windowId:string` (opc)<br>`kind:ScanKind` (opc)<br>`severity:ProblemSeverity` (opc)<br>`limit:int` (opc) | `findings:ScanFinding[]` (wym)<br>`total:int` (opc) |
| `developer.scan.run` | Uruchamia skan zależności, sekretów, kodu albo licencji | `windowId:string` (wym)<br>`kinds:ScanKind[]` (wym)<br>`paths:string[]` (opc) | `scan:ScanRun` (wym) |
| `developer.symbol.navigate` | Zwraca miejsca symbolu wskazane przez serwer języka | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`column:int` (wym)<br>`kind:SymbolNavigationKind` (wym) | `symbols:DeveloperSymbol[]` (wym)<br>`serverAvailable:bool` (wym) |
| `developer.test.result.get` | Zwraca wynik zestawu testów przebiegu budowania | `buildId:string` (wym)<br>`status:TestStatus` (opc) | `results:DeveloperTestResult[]` (wym)<br>`passed:int` (wym)<br>`failed:int` (wym)<br>`skipped:int` (wym)<br>`durationMs:int64` (opc) |
| `developer.toolchain.check` | Sprawdza obecność programów zewnętrznych w ścieżce wykonywalnej rdzenia | `programs:string[]` (wym) | `programs:ToolchainProgram[]` (wym) |
| `developer.tree.get` | Zwraca drzewo projektu odzwierciedlające stan repozytorium | `windowId:string` (wym)<br>`path:string` (opc)<br>`depth:int` (opc)<br>`includeHidden:bool` (opc) | `root:string` (wym)<br>`nodes:DeveloperTreeNode[]` (wym) |

**Zdarzenia obszaru `developer` — 2:**

| Zdarzenie | Przeznaczenie | Pola ładunku |
|---|---|---|
| `developer.build.changed` | Przyrost budowania; log narasta w toku przebiegu | `change:ChangeKind` (wym)<br>`build:DeveloperBuild` (wym)<br>`logLine:string` (opc) |
| `developer.debug.changed` | Zmiana stanu sesji debugowania: zatrzymanie, wznowienie, zakończenie | `change:ChangeKind` (wym)<br>`session:DebugSession` (wym)<br>`frame:DebugFrame` (opc)<br>`reason:string` (opc) |

Razem w wykazie: **50 komend** z jednego obszaru kontraktu.

---

*Koniec dokumentu. Moduł Developer — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
