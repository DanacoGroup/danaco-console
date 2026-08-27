*Dokument opisuje środowisko MultitaskingAI platformy Danaco Console: warstwę centralną, role zespołu, silnik kolejek, orkiestrację, powłokę oraz katalog elementów interfejsu.*

# Danaco Console — Środowisko MultitaskingAI: specyfikacja projektowa interfejsu

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
| **Tytuł** | Środowisko MultitaskingAI — warstwa centralna, kanały komunikacji, role, silnik kolejek, panel orkiestracji, powłoka środowiska, integracja z Automations, makiety i katalog interfejsu |
| **Przeznaczenie** | Pełnozakresowy, gęsty wizualnie materiał źródłowy do projektowania i budowy interfejsu środowiska MultitaskingAI: co ma powstać, gdzie ma leżeć, w jakiej formie i do czego służy — dla każdego elementu, każdej roli, każdego okna |
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Zakres** | Środowisko w całości: warstwa centralna, kanały komunikacji, role, silnik kolejek, orkiestracja, integracja z Automations oraz powłoka środowiska — nawigacja, przestrzeń robocza okien ról, karty sesji, panel stanu procesu, skróty i paleta poleceń, współdzielenie kontekstu, funkcje wspólne i punkty sterowania z okna konfiguracji |
| **Środowisko** | MultitaskingAI — czwarte, najbardziej zaawansowane konfiguracyjnie środowisko platformy (Koncepcja platformy, rozdz. 9.4 i 13) |
| **Autorytatywne źródło** | architektura/koncepcja-platformy.md |
| **Źródła pomocnicze** | specyfikacje/specyfikacja-modulow.md, specyfikacje/specyfikacja-okien-operacyjnych.md, interfejs-uzytkownika/strona-glowna-i-nawigacja.md, interfejs-uzytkownika/system-wizualny.md, architektura/model-konfiguracji.md, architektura/model-danych.md, architektura/architektura.md, architektura/integracja-modeli.md, specyfikacje/specyfikacja-agentow.md, architektura/izolacja-i-zaleznosci.md, architektura/rozszerzenia.md, `interfejs-uzytkownika/przeplyw-okien.md`, `interfejs-uzytkownika/elementy-okien.md`, `funkcje-globalne/always-on-display.md` |
| **Opracowanie** | Danaco Console — dokumentacja projektu UI |

Dokument zestawia w jednym miejscu, w formie gęstej wizualnie (tabele, schematy, diagramy przepływu, makiety tekstowe), wszystko, co potrzebne do zaprojektowania i zbudowania interfejsu środowiska MultitaskingAI: warstwę centralną Always On Display, dwa kanały komunikacji operacyjnej (Chat Window oraz Execution Loop Window), pięć pozycji zespołu (Executor 1, Subagent Network, Coordinator, Executor 2, Executor 3 / Validator) z pełnym oprzyrządowaniem każdej z nich, silnik kolejek z jedenastoma akcjami, panel orkiestracji jako boczną nawigację środowiska, integrację z modułem Automations dającą pracę ciągłą 24 godziny na dobę, 7 dni w tygodniu, 365 dni w roku, powłokę środowiska wraz z jej funkcjami i punktami sterowania, makiety tekstowe wszystkich okien środowiska, zbiorczy katalog elementów interfejsu oraz katalog stanów. Dokument jest jedynym opracowaniem środowiska MultitaskingAI w zbiorze: łączy warstwę systemową (definicje mechanizmów, role, silnik kolejek, orkiestracja, hierarchia decyzji, zgodność ze źródłami) z warstwą projektową interfejsu (makiety, katalogi elementów, stany, powłoka środowiska).

**Zasada nadrzędna obowiązująca cały dokument.** Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód. Domyślne zachowanie systemu to wykonanie polecenia. Izolacja techniczna, izolacja kontekstu i wszelkie ograniczenia uprawnień są ustawieniami konfiguracyjnymi sterowanymi przez Operatora z okna konfiguracji — nigdy wymogiem stawianym przez platformę. Wszędzie tam, gdzie w dokumencie pojawia się słowo „ograniczenie”, „profil izolacji” lub „uprawnienia”, czyta się je jako ustawienie do świadomego włączenia, ze stanem wyjściowym „wyłączone / pełny dostęp”.

**Zasada układu.** Wszystkie okna środowiska rozmieszczone są w układzie pionowym, w kolumnach sąsiadujących poziomo. Chat Window zajmuje lewą kolumnę o pełnej wysokości obszaru roboczego, Execution Loop Window otwiera się jako kolumna sąsiadująca, obszar roboczy ról zajmuje kolumnę dominującą po prawej, a okna pomocnicze i panele otwierają się jako kolejne rozszerzenia boczne. Regulacji podlega wyłącznie szerokość kolumn.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
1. [Miejsce środowiska w platformie i architektura ogólna](#1-miejsce-środowiska-w-platformie-i-architektura-ogólna)
2. [Warstwa centralna — Always On Display](#2-warstwa-centralna--always-on-display)
3. [Role zespołu środowiska MultitaskingAI](#3-role-zespołu-środowiska-multitaskingai)
4. [Silnik kolejek](#4-silnik-kolejek)
5. [Orkiestracja i zależności](#5-orkiestracja-i-zależności)
6. [Panel orkiestracji — boczna nawigacja środowiska](#6-panel-orkiestracji--boczna-nawigacja-środowiska)
7. [Integracja z modułem Automations — praca ciągła 24/7/365](#7-integracja-z-modułem-automations--praca-ciągła-247365)
8. [Hierarchia decyzji](#8-hierarchia-decyzji)
9. [Okna środowiska — makiety tekstowe](#9-okna-środowiska--makiety-tekstowe)
10. [Katalog elementów interfejsu](#10-katalog-elementów-interfejsu)
11. [Diagramy przepływu — orkiestracja i kolejkowanie](#11-diagramy-przepływu--orkiestracja-i-kolejkowanie)
12. [Stany](#12-stany)
13. [Konfigurowalność i punkty izolacji na poziomie roli](#13-konfigurowalność-i-punkty-izolacji-na-poziomie-roli)
14. [Szablony konfiguracji zespołu](#14-szablony-konfiguracji-zespołu)
15. [Scenariusze operacyjne](#15-scenariusze-operacyjne)
16. [Zgodność z zasadami nadrzędnymi platformy](#16-zgodność-z-zasadami-nadrzędnymi-platformy)
17. [Powłoka środowiska — nawigacja i przestrzeń robocza okien](#17-powłoka-środowiska--nawigacja-i-przestrzeń-robocza-okien)
18. [Karty sesji — procesy orkiestracji i trwałość stanu](#18-karty-sesji--procesy-orkiestracji-i-trwałość-stanu)
19. [Panel stanu procesu, skróty klawiszowe i paleta poleceń](#19-panel-stanu-procesu-skróty-klawiszowe-i-paleta-poleceń)
20. [Współdzielenie kontekstu między rolami, procesami i modułami](#20-współdzielenie-kontekstu-między-rolami-procesami-i-modułami)
21. [Funkcje wspólne środowiska](#21-funkcje-wspólne-środowiska)
22. [Punkty sterowania z okna konfiguracji](#22-punkty-sterowania-z-okna-konfiguracji)
23. [Słownik pojęć](#23-słownik-pojęć)
- [Załącznik A. Katalog ikon i komponentów systemu wizualnego wykorzystanych w środowisku](#załącznik-a-katalog-ikon-i-komponentów-systemu-wizualnego-wykorzystanych-w-środowisku)
- [Załącznik B. Macierz zgodności ze źródłami](#załącznik-b-macierz-zgodności-ze-źródłami)

---

## Wprowadzenie

MultitaskingAI jest jedynym środowiskiem platformy, w którym jednostką organizacji pracy nie jest pojedynczy użytkownik wspierany przez AI, lecz zespół modeli i agentów realizujących wspólny proces. Zamiast bocznej nawigacji modułów środowisko udostępnia panel orkiestracji — sześć sekcji sterowania zespołem — a obok Chat Window, głównego okna komunikacji Użytkownik ↔ Wykonawca obecnego w każdej przestrzeni roboczej platformy, prowadzi Execution Loop Window (Koordynator ↔ Wykonawca), cztery okna robocze ról, mechanizm Subagent Network uruchamiany przez wykonawcę, silnik kolejek o jedenastu akcjach oraz warstwę orkiestracji egzekwującą zależności między nimi. Po spięciu z modułem Automations całość działa jako autonomiczna, wieloagentowa pętla pracująca w sposób ciągły bez stałej obecności użytkownika przy stanowisku roboczym.

Środowisko obejmuje również własną **powłokę** — ramę operacyjną orkiestracji wieloagentowej. Powłoka spina pasek nawigacji, karty sesji, panel orkiestracji, wielookienną przestrzeń roboczą ról, panel stanu procesu oraz warstwę centralną Always On Display w jedną przestrzeń, w której użytkownik prowadzi równolegle wiele procesów zespołowych. W odróżnieniu od trzech środowisk modułowych (TalkIn, WorkSpace, CodeStudio), których powłoka jest wspólna, MultitaskingAI ma powłokę własną: nawiguje nie po modułach, lecz po rolach, kolejkach i orkiestracji (rozdz. 6.1). Powłoka udostępnia okna ról i grupuje akcje silnika kolejek, przypisuje do ról gotowych agentów zbudowanych w module Agents i wiąże proces z automatykami modułu Automations; wewnętrzna logika ról (rozdz. 3), semantyka akcji kolejek (rozdz. 4), budowa agentów oraz silnik harmonogramów należą do mechanizmów opisanych w odpowiednich rozdziałach i modułach. Granica powłoki jest kompozycyjna: wiąże się ona z modułami i mechanizmami przez jawne, konfigurowalne powiązania (rozdz. 20, 22).

Dokument czyta się w pięciu warstwach nałożonych na siebie:

| Warstwa dokumentu | Co dostarcza | Rozdziały |
|---|---|---|
| Mechanizm | Definicja, cel, wejście/wyjście każdej roli i każdego mechanizmu | 1–8 |
| Forma | Makiety tekstowe okien — co gdzie leży w przestrzeni roboczej | 9 |
| Element | Katalog każdego pojedynczego elementu interfejsu — co to jest, forma, warstwa widoczności, stany, zachowanie | 10 |
| Zachowanie w czasie | Diagramy przepływu i katalog stanów — jak system zmienia się w toku działania | 11–12 |
| Powłoka | Funkcje powłoki środowiska: nawigacja, przestrzeń robocza okien, karty sesji, panel stanu, skróty, kontekst, funkcje wspólne | 17–21 |

Rozdziały 13–16 oraz 22–23 i załączniki domykają dokument o konfigurowalność, gotowe szablony, scenariusze, zgodność z zasadami platformy, punkty sterowania z okna konfiguracji i słownik. Wszystkie nazwy własne — ról, okien, akcji, sekcji — są przejęte bez zmian ze źródeł wymienionych w metryczce; dokument nie wprowadza skrótów ani kodów odsyłających w miejsce pełnych nazw.

---

## 1. Miejsce środowiska w platformie i architektura ogólna

### 1.1. Definicja i odróżnienie od pozostałych trzech środowisk

MultitaskingAI jest zaawansowanym środowiskiem wielomodelowym, umożliwiającym równoczesną pracę wielu instancji AI w ramach jednego projektu, procesu lub zadania — od prostych układów współpracy dwóch modeli po złożone, wieloetapowe systemy autonomiczne.

| Cecha | TalkIn / WorkSpace / CodeStudio | MultitaskingAI |
|---|---|---|
| Jednostka organizacji pracy | Pojedynczy użytkownik wspierany przez AI, w ramach jednego zadania | Zespół modeli i agentów realizujących wspólny proces |
| Zawartość bocznej nawigacji | Lista modułów dostępnych w środowisku | Panel orkiestracji — sześć sekcji sterowania zespołem (rozdz. 6) |
| Podstawowa jednostka pracy | Moduł → okno operacyjne | Rola → akcja silnika kolejek → orkiestracja |
| Obecność w macierzy dostępności modułów | Tak — jako kolumna macierzy | Nie — środowisko nie udostępnia modułów w bocznej nawigacji |
| Domyślny tryb pracy | Interakcja konwersacyjna, jednorazowa lub sesyjna | Zdolność do pracy ciągłej po skonfigurowaniu integracji z Automations (rozdz. 7) |
| Rola Always On Display | Doradztwo kontekstowe | Doradztwo kontekstowe oraz funkcja obserwatora lub operatora procesu (rozdz. 2) |

### 1.2. Odróżnienie od modułu Roundtable

| Cecha | Moduł Roundtable | Środowisko MultitaskingAI |
|---|---|---|
| Forma wielomodelowości | Debata i wypracowanie konsensusu | Zorganizowany podział pracy między wyspecjalizowane role |
| Sposób pracy modeli | Kilka modeli odpowiada równolegle na to samo zagadnienie | Role realizują rozdzielone zakresy wspólnego procesu |
| Mechanizmy porządkujące | Moderator Panel, Consensus Panel | Podział ról, silnik kolejek, orkiestracja zależności |
| Punkt ciężkości | Debata | Zorganizowany podział pracy |

### 1.3. Typowe zastosowania

| Zastosowanie | Charakterystyka |
|---|---|
| Wieloetapowe projekty w module Apps | Podział na Executor 1 i Executor 2 przy równoległej pracy nad backendem i frontendem |
| Procesy badawcze, redakcyjne, diagnostyczne | Dowolny proces wymagający podziału pracy między wiele modeli |
| Pętla pracy ciągłej 24/7/365 | Po spięciu z modułem Automations — realizacja projektów bez stałego nadzoru |

Środowisko nie jest ograniczone do modułu Apps — role środowiska MultitaskingAI mogą obsługiwać dowolny proces wymagający podziału pracy między wiele modeli.

### 1.4. Położenie w hierarchii platformy

```
STRONA GŁÓWNA
    │  wybór środowiska w strefie 1
    ▼
MultitaskingAI
    │
    ├── Panel orkiestracji (boczna nawigacja — rozdz. 6)
    │       Zespoły · Role · Kolejki · Orkiestracja · Harmonogram i automatyki · Monitor procesu
    │
    ├── Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca
    │       lewa kolumna, stała, pełna wysokość obszaru roboczego
    │
    ├── Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca
    │       kolumna sąsiadująca z Chat Window
    │
    ├── Cztery okna robocze (role)
    │       Executor 1 · Coordinator · Executor 2 · Executor 3 / Validator
    │       └── Subagent Network (mechanizm uruchamiany przez wykonawcę, do 15 podagentów)
    │
    └── Karty sesji
            każda karta = odrębny proces orkiestracji, z własnym zespołem ról

PONAD CAŁĄ STRUKTURĄ
    Always On Display — warstwa centralna (rozdz. 2)
    Mobile — zatwierdzanie i interwencja zdalna (rozdz. 2.4, rozdz. 7.5)
    Automations — po spięciu: harmonogram i trwałe zarządzanie wykonaniem cyklicznym (rozdz. 7)
```

### 1.5. Schemat warstw wewnętrznych środowiska

Poniższy schemat porządkuje sześć warstw funkcjonalnych, z których zbudowane jest środowisko — każda odpowiada rozdziałowi niniejszego dokumentu.

```
╔═══════════════════════════════════════════════════════════════════════╗
║  WARSTWA CENTRALNA · Always On Display                                  ║
║  nadzór nad całością — obserwator albo operator procesu (rozdz. 2)      ║
╠═══════════════════════════════════════════════════════════════════════╣
║  WARSTWA RÓL · pięć pozycji zespołu                                     ║
║  Executor 1 · Subagent Network · Coordinator · Executor 2 ·             ║
║  Executor 3 / Validator (rozdz. 3)                                      ║
╠═══════════════════════════════════════════════════════════════════════╣
║  WARSTWA SILNIKA KOLEJEK · jedenaście akcji, pięć zasięgów (rozdz. 4)   ║
╠═══════════════════════════════════════════════════════════════════════╣
║  WARSTWA ORKIESTRACJI · zależności sekwencyjne, warunkowe, równoległe   ║
║  (rozdz. 5)                                                              ║
╠═══════════════════════════════════════════════════════════════════════╣
║  WARSTWA PANELU ORKIESTRACJI · boczna nawigacja, sześć sekcji (rozdz. 6)║
╠═══════════════════════════════════════════════════════════════════════╣
║  WARSTWA INTEGRACJI Z AUTOMATIONS · praca ciągła 24/7/365 (rozdz. 7)    ║
╚═══════════════════════════════════════════════════════════════════════╝
```

### 1.6. Warstwy widoczności w środowisku

Interfejs Danaco Console ujawnia możliwości systemu stopniowo: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. W środowisku MultitaskingAI reguła ta ma znaczenie kluczowe — środowisko prowadzi wiele równoległych wątków (cztery role, dwie niezależne instancje Subagent Network po 15 podagentów, kolejki pięciu zasięgów, zależności orkiestracji, harmonogram pracy ciągłej), a mimo tej skali interfejs pozostaje wizualnie prosty. Złożoność środowiska istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba równolegle prowadzonych procesów, ról, podagentów, kolejek i zależności nie wpływa na postrzeganą prostotę przestrzeni roboczej.

Każdy element interfejsu środowiska należy do dokładnie jednej z czterech warstw widoczności. Katalogi elementów (rozdz. 10) oraz tabele okien (rozdz. 9) podają warstwę i sposób wywołania każdego elementu.

| Warstwa | Nazwa | Zawartość w środowisku MultitaskingAI | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, Execution Loop Window, aktywne okno wiodące roli, pasek kontekstu procesu, podstawowa nawigacja panelu orkiestracji, wskaźniki stanu wykonania ról i kolejek. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór agenta lub modelu bazowego dla roli, wybór kanału modelu, wybór trybu współpracy wykonawców, wybór wcielenia Executora 3 / Validatora, wybór zasięgu kolejki, tryb Always On Display | Znacznik kontekstowy, ikona, przycisk, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestaw jedenastu akcji silnika kolejek, ustawienia szybkie karty sesji, warianty operacji na sekcjach panelu, panel Subagent Network, panel porównania wyników | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Macierz punktów izolacji na poziomie roli, edytor skrótów klawiszowych, dziennik zdarzeń środowiska, tryb administracyjny, narzędzia diagnostyczne kanału modelu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, paleta poleceń, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Mechanizmy ukrywania funkcjonalności stosowane w środowisku:

- **Menu progresywne.** Jednorodne zbiory wyborów prezentowane są jako jeden element zwinięty (`Rola ▼`, `Agent ▼`, `Wcielenie ▼`), a lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Panel Subagent Network, panel porównania wyników, kreator reguły warunkowej i macierz izolacji otwierają się jako kolumny boczne i panele popover; po zamknięciu znikają całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Jedenaście akcji silnika kolejek prezentowanych jest jako jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji; sterowanie procesem prezentowane jest jako `Sterowanie ▼`.
- **Znaczniki kontekstowe.** Środowisko, proces, zespół, rola, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [MultitaskingAI] [Zespół „Budowa aplikacji"] [Executor 1] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja środowiska jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

Makiety w niniejszym dokumencie przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki kontekstowe, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są w tabelach elementów okien z podaniem warstwy i sposobu wywołania.

### 1.7. Dwa kanały komunikacji operacyjnej środowiska

Środowisko prowadzi dwa kanały komunikacji operacyjnej. Oba są elementami pierwszoplanowymi architektury środowiska.

| Kanał | Okno | Uczestnicy | Zakres | Rozmieszczenie | Warstwa |
|---|---|---|---|---|---|
| Kanał pierwszy | **Chat Window** | Użytkownik ↔ Wykonawca | Główne okno komunikacji: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie i przerywanie działań oraz wyjaśnianie wyniku i kontekstu. Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami środowiska | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 |
| Kanał drugi | **Execution Loop Window** | Koordynator ↔ Wykonawca | Okno pętli wykonawczej: prowadzenie pętli, koordynacja zadań, nadzór nad realizacją, orkiestracja działań i kontrola realizacji procesów przy pracy wielowątkowej. Zawiera bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem | Kolumna sąsiadująca z Chat Window | 1 |

**Chat Window w środowisku MultitaskingAI.** Okno jest obecne w każdej przestrzeni roboczej środowiska — w każdej karcie sesji, przy każdym układzie okien ról i w każdej sekcji panelu orkiestracji — zawsze w tym samym miejscu układu, w lewej kolumnie obszaru roboczego. Polecenie wydane w Chat Window steruje całym procesem: uruchamia i zatrzymuje pętlę, przypisuje role, dodaje zadania do kolejki, wywołuje dowolną z jedenastu akcji silnika kolejek i otwiera dowolne okno środowiska.

**Execution Loop Window w środowisku MultitaskingAI.** Okno jest pierwszoplanowym elementem tego środowiska, ponieważ MultitaskingAI z definicji prowadzi pracę wielowątkową. Execution Loop Window prezentuje pętlę wykonawczą prowadzoną przez Koordynatora wobec wykonawców: dekompozycję zlecenia na zadania, przydział zadań rolom, stan kolejek pięciu zasięgów, komunikaty sterujące (`route`, `retry`, `pause`, `resume`, `branch`, `condition`), wyniki kontroli jakości Executora 3 / Validatora, wskaźniki przebiegu pętli oraz sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Rola Koordynatora w kanale drugim odpowiada roli Coordinatora w składzie zespołu (rozdz. 3.3); Wykonawcą są Executor 1, Executor 2 oraz — pośrednio, przez wykonawcę macierzystego — podagenci Subagent Network.

Role: **Użytkownik** — zleca i zatwierdza; **Koordynator** — komponent orkiestrujący platformy, dekomponuje zlecenie, przydziela i nadzoruje zadania; **Wykonawca** — AI, agent lub system wykonawczy realizujący zadania.

**Odpowiedniki ról kanonicznych w środowisku.** Nazwy ról kanonicznych platformy odpowiadają rolom operacyjnym środowiska MultitaskingAI w następujący sposób.

| Rola kanoniczna | Odpowiednik w środowisku MultitaskingAI |
|---|---|
| Użytkownik | Użytkownik (Operator) — poziom 1 hierarchii decyzji (rozdz. 8) |
| Koordynator | Coordinator — komponent orkiestrujący, dekomponuje zlecenie, przydziela i nadzoruje zadania (rozdz. 3.3) |
| Wykonawca | Executor 1, Executor 2, Executor 3 / Validator oraz podagenci Subagent Network (rozdz. 3.1, 3.2, 3.4, 3.5) |

**Prowadzenie pętli wykonawczej.** Pętla wykonawcza jest podstawowym cyklem pracy środowiska. Koordynator prowadzi ją w powtarzalnym porządku, a Execution Loop Window przedstawia każdy krok tego porządku w czasie rzeczywistym.

| Krok pętli | Działanie Koordynatora | Widok w Execution Loop Window |
|---|---|---|
| 1. Przyjęcie zlecenia | Odebranie celu procesu ustalonego przez Użytkownika w Chat Window | Zlecenie w nagłówku okna |
| 2. Dekompozycja | Podział zlecenia na zadania, ustalenie zależności sekwencyjnych, warunkowych i równoległych (rozdz. 5.2) | Drzewo zadań z oznaczeniem zależności |
| 3. Przydział | Skierowanie zadań do ról i kolejek (`enqueue`, `route`), budowa promptu dla każdego Wykonawcy | Pozycje kolejki z przypisaniem roli |
| 4. Nadzór wykonania | Śledzenie pobrań (`dequeue`), postępu i wyników cząstkowych wszystkich równoległych wątków | Wskaźniki przebiegu i strumień komunikatów |
| 5. Kontrola jakości | Skierowanie wyniku do Executora 3 / Validatora przed uznaniem etapu za zakończony | Wynik oceny przy zadaniu |
| 6. Decyzja | Kontynuacja, `retry`, `branch` na inną ścieżkę albo eskalacja do Użytkownika | Decyzja zapisana w przebiegu |
| 7. Zamknięcie przebiegu | Scalenie wyników (`merge`), przekazanie rezultatu i uruchomienie kolejnego przebiegu | Podsumowanie przebiegu i licznik kolejnego |

**Koordynacja wielu równoległych wątków.** Środowisko prowadzi jednocześnie dwa tory wykonawcze, z których każdy uruchamia do 15 podagentów, co daje do 30 jednostek wykonawczych w jednym procesie; kart sesji, a więc i procesów, działa wiele równolegle (rozdz. 1.4, 18). Execution Loop Window utrzymuje czytelność tego obrazu przez następujące mechanizmy.

| Mechanizm | Działanie |
|---|---|
| Ścieżki wątków | Każdy tor wykonawczy prezentowany jest jako odrębna ścieżka pętli, z własnym licznikiem etapów i własnym stanem |
| Zwijanie wątków podrzędnych | Praca Subagent Network prezentowana jest jako jedna pozycja zbiorcza wykonawcy macierzystego; rozwinięcie ujawnia listę podagentów |
| Filtr zakresu | Widok zawężany do wskazanego wątku, roli, kolejki lub karty sesji |
| Sygnalizacja wyjątków | Wątki wymagające decyzji Użytkownika wysuwają się na początek listy; pozostałe pozostają zwinięte |
| Agregacja stanu | Nagłówek okna prezentuje jeden zbiorczy stan procesu: liczbę wątków aktywnych, wstrzymanych, oczekujących na zatwierdzenie i zakończonych błędem |

Nadzór nad całością procesu z perspektywy Użytkownika sprawuje Always On Display w trybie obserwatora albo operatora (rozdz. 2.2), a kanał interwencji z dowolnego miejsca udostępnia funkcja Mobile (rozdz. 2.4).

---

## 2. Warstwa centralna — Always On Display

### 2.1. Charakterystyka i umiejscowienie w środowisku

Always On Display jest funkcją globalną platformy — globalnym agentem towarzyszącym, nieposiadającym własnego środowiska ani modułu, obecnym jednocześnie we wszystkich częściach platformy. Pełną dokumentację funkcji — postać wizualną, reguły wyzwalania proaktywnych sugestii, katalog rodzajów sugestii, tor głosowy i jego granicę wobec modułu Assistant, warstwy widoczności oraz punkty sterowania — zawiera `funkcje-globalne/always-on-display.md`. Niniejszy rozdział opisuje wyłącznie specyfikę środowiska MultitaskingAI: rolę warstwy centralnej wobec zespołu ról, kolejek i przebiegu procesu.

W środowisku MultitaskingAI funkcja pełni rolę warstwy centralnej — nadrzędnej wobec wszystkich pięciu pozycji zespołu, z dostępem do wszystkich środowisk, projektów, agentów, sesji, historii i procesów. Ponieważ nie jest przypisana do żadnego pojedynczego modułu i ma dostęp do pełnego kontekstu działania użytkownika, jest naturalnym kandydatem do pełnienia roli nadzorczej nad procesem angażującym wiele modeli i wiele ról jednocześnie.

### 2.2. Dwa tryby względem procesu MultitaskingAI

| Tryb | Zakres działania | Typowe zastosowanie |
|---|---|---|
| Obserwator | Podgląd przebiegu pętli, statusów przebiegów i hierarchii decyzji (sekcja Monitor procesu); zgłaszanie sugestii i ostrzeżeń jako proaktywne doradztwo, bez ingerencji w przebieg procesu | Nadzór nad procesem o niskim ryzyku, decyzje pozostają w gestii Coordinatora i użytkownika |
| Operator | Wszystkie uprawnienia obserwatora, rozszerzone o zdolność interwencji: zatwierdzanie, wstrzymywanie kroków procesu, interwencja z dowolnego miejsca | Nadzór nad procesem działającym w trybie ciągłym (rozdz. 7), gdzie brak stałej obecności użytkownika przy stanowisku wymaga zdolności do samodzielnej interwencji |

Wybór trybu, jak wszystkie pozostałe zależności platformy, należy do użytkownika i pozostaje odwracalny w dowolnym momencie — nie jest to bramka bezpieczeństwa, lecz przełącznik zakresu działania. Dołączenie i odłączenie funkcji od procesu realizują polecenia `aod.observe.attach` i `aod.observe.detach`.

### 2.3. Zdolności właściwe nadzorowi nad procesem środowiska

Poniższa tabela obejmuje zdolności swoiste dla nadzoru nad zespołem ról, kolejkami i przebiegiem pętli. Zdolności wspólne całej platformie — proaktywne doradztwo, pomoc kontekstowa, komunikacja tekstowa i głosowa, pełny dostęp do środowisk, projektów, sesji i historii — opisuje `funkcje-globalne/always-on-display.md` (rozdz. 1.3, 3–5).

| Zdolność | Opis | Dostępna w trybie |
|---|---|---|
| Podgląd Monitora procesu | Stan wszystkich czterech ról, zawartość kolejek, zdefiniowane zależności orkiestracji, reguły harmonogramu, historia przebiegów | Obserwator, Operator |
| Sugestie wobec zespołu | Sugestie konfiguracji zespołu, wskazywanie powtarzających się niepowodzeń ról, rekomendacje kolejnego kroku procesu | Obserwator, Operator |
| Zatwierdzenie kroku | Akceptacja etapu procesu oczekującego na decyzję (np. wyniku Executora 3 / Validatora) | Operator |
| Wstrzymanie procesu | Wywołanie akcji `pause` na wskazanej kolejce lub całym procesie | Operator |
| Wznowienie procesu | Wywołanie akcji `resume` | Operator |
| Interwencja z dowolnego miejsca | Uruchomienie dowolnej z powyższych zdolności operatora przez kanał funkcji Mobile, niezależnie od lokalizacji użytkownika | Operator (przez Mobile) |

### 2.4. Współdziałanie z funkcją Mobile

| Element | Rola w interwencji |
|---|---|
| Always On Display (tryb operatora) | Obserwuje przebieg procesu i inicjuje interwencję (np. zatwierdzenie kroku, wstrzymanie procesu) |
| Mobile | Udostępnia kanał, przez który interwencja jest faktycznie wykonywana z dowolnego urządzenia, niezależnie od lokalizacji |

```
 Always On Display (operator)  ──inicjuje──►  decyzja interwencyjna
         │
         │  kanał wykonania
         ▼
      Mobile  ──►  zatwierdzenie / wstrzymanie / modyfikacja procesu
                   z dowolnego urządzenia, poza stanowiskiem roboczym
```

### 2.5. Zakres dostępu — pokrycie sekcji Monitor procesu

| Element kontekstu procesu | Zakres podglądu Always On Display |
|---|---|
| Role | Stan wszystkich czterech ról (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) |
| Kolejki | Zawartość kolejek wszystkich pięciu zasięgów |
| Orkiestracja | Zdefiniowane zależności sekwencyjne, warunkowe, równoległe |
| Harmonogram | Reguły czasowe i cykliczność pracy ciągłej |
| Historia | Historia przebiegów procesu |

### 2.6. Elementy interfejsu swoiste dla środowiska

Awatar funkcji, dymek kontekstowy sugestii, powierzchnia interakcji, plakietka powiadomień i przycisk mikrofonu są elementami wspólnymi całej platformy — ich formę, stany, warstwy widoczności i sposób wywołania podaje `funkcje-globalne/always-on-display.md` (rozdz. 2.6, 8). Poniższa tabela obejmuje elementy występujące wyłącznie w środowisku MultitaskingAI.

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Przełącznik trybu (obserwator / operator) | Segmentowana kontrolka dwustanowa | Wybór zakresu działania funkcji wobec procesu (rozdz. 2.2) | Przełącznik segmentowy / suwak (`.dn-suwak`) | Domyślny „obserwator” · zaznaczony „operator” | Zmiana trybu natychmiast rozszerza lub zawęża dostępne akcje interwencyjne; ustawienie „operator” pozostaje dostępne również przy braku aktywnego procesu — komunikat kontekstowy informuje wtedy, że nadzór obejmie najbliższe uruchomienie | Sekcja Monitor procesu panelu orkiestracji | 2 | Znacznik `AOD ▼` w panelu stanu procesu |
| Wskaźnik roli nadzoru | Plakietka przy nazwie funkcji w sekcji Monitor procesu | Sygnalizacja, czy funkcja nadzoruje bieżący przebieg jako obserwator, czy jako operator | Mała plakietka (`.dn-plakietka`) | Obserwator · operator · nieprzypisana do procesu | Kliknięcie otwiera przełącznik trybu | Sekcja Monitor procesu panelu orkiestracji | 2 | Widoczny po otwarciu sekcji Monitor procesu |
| Przycisk zatwierdzenia kroku | Kontrolka trybu operatora | Zatwierdzenie etapu procesu oczekującego na decyzję | Przycisk (`.dn-btn--zloty`) | Domyślny · ładowanie (przetwarzanie decyzji) · zatwierdzony | Kliknięcie w trybie obserwatora wyświetla komunikat kontekstowy ze wskazaniem przełączenia w tryb operatora, zamiast blokować przycisk; kliknięcie w trybie operatora zatwierdza krok i przekazuje proces do kolejnego kroku orkiestracji | Dymek funkcji, sekcja Monitor procesu | 2 | Przycisk akcji w dymku funkcji |

---

## 3. Role zespołu środowiska MultitaskingAI

MultitaskingAI udostępnia cztery okna robocze odpowiadające czterem rolom oraz jeden mechanizm — Subagent Network — uruchamiany przez rolę wykonawczą. Podział na rolę wykonującą (Executor) i rolę zarządzającą (Coordinator) jest zasadą fundamentalną środowiska: obie odpowiedzialności są rozdzielone, aby żaden pojedynczy model nie musiał jednocześnie planować i realizować.

### 3.0. Tabela zbiorcza — cel, wejście, wyjście

| Rola | Cel roli | Wejście | Wyjście |
|---|---|---|---|
| Executor 1 — główny wykonawca | Faktyczna realizacja pracy: tworzenie dokumentów, analiza danych, programowanie, projektowanie, budowa aplikacji, przetwarzanie materiałów; nie zarządza procesem | Zadanie z kolejki przypisanej roli, polecenie użytkownika, prompt zbudowany przez Coordinatora | Rezultat pracy, status wykonania do silnika kolejek, podzadania do Subagent Network przy złożonych zleceniach |
| Subagent Network | Mechanizm uruchamiany przez wykonawcę, nie odrębna rola | Zakres zadania wydzielony przez wykonawcę (`split`) | Wynik cząstkowy agregowany przez wykonawcę (`merge`) |
| Coordinator — koordynator | Nie tworzy końcowego produktu; projektuje i kontroluje sposób jego wytwarzania | Cel procesu ustalony przez użytkownika, wyniki cząstkowe wykonawców, raporty Executora 3 / Validatora | Plan pracy, prompty dla wykonawców, polecenia sterujące kolejką i orkiestracją |
| Executor 2 — równoległy wykonawca | Drugi niezależny model wykonawczy, procesy równoległe (backend/frontend, badania/raport, kod/dokumentacja) | Jak Executor 1; dodatkowo wyniki pośrednie Executora 1, zależnie od trybu współpracy | Jak Executor 1; wynik przekazywany zależnie od trybu współpracy |
| Executor 3 / Validator — czwarty model | Funkcja kontrolna lub doradcza domykająca zespół; odpowiada za jakość, zgodność lub rozstrzyganie rozbieżności | Rezultat pracy wykonawców, wynik pracy Subagent Network po agregacji | Ocena, raport, decyzja o przekazaniu dalej lub zwrocie do Coordinatora (`retry`) |

### 3.1. Executor 1 — główny wykonawca

**Pełne oprzyrządowanie.** Poniższa tabela zestawia wszystkie zdolności i narzędzia dostępne roli — od kanału połączenia z modelem po mechanizmy uruchamiane w toku pracy.

| Grupa narzędzi | Zawartość |
|---|---|
| Przypisanie wykonawcy | Agent skonfigurowany w module Agents (siedem komponentów definicji: model bazowy, tożsamość, instrukcje systemowe, skille, pluginy i konektory, pamięć, uprawnienia) albo model bazowy podłączony bezpośrednio jednym z czterech kanałów |
| Kanał modelu | API · CLI · SSH · HTTP — wybór dokonywany per rola, nadpisujący kanał domyślny zapisany w definicji agenta |
| Rozszerzenia dostępne wykonawcy | Wtyczki i umiejętności (Skills Manager) rozszerzające sposób wykonania zadania; konektory i serwery MCP (Connectors Manager) otwierające dostęp do systemów poza platformą |
| Pamięć | Poziomy: globalna, projekt, sesja, środowisko — korzystanie z pamięci projektu jako pole wejścia/wyjścia roli |
| Akcje silnika kolejek dostępne roli | `dequeue` — pobranie zadania; `split` — podział złożonego zadania przy uruchamianiu Subagent Network; `merge` — scalenie wyników podagentów |
| Uruchamianie Subagent Network | Do 15 równoczesnych podagentów (rozdz. 3.2) |
| Profil izolacji | Poziom zasięgu „Rola” — najwyższe pierwszeństwo spośród siedmiu poziomów; konfigurowalny, domyślnie brak aktywnej izolacji technicznej |
| Okno robocze | Executor Chat (rozdz. 9) |

**Charakterystyka.** Rola podstawowa — punkt wyjścia dla każdej konfiguracji zespołu; odpowiada za faktyczne wytwarzanie rezultatu pracy niezależnie od tego, czy w danym procesie uczestniczą pozostałe role.

### 3.2. Subagent Network

| Aspekt | Zawartość |
|---|---|
| Charakter | Mechanizm, nie odrębna rola zespołu — uruchamiany przez wykonawcę (Executor 1 lub Executor 2) |
| Zakres | Do 15 wyspecjalizowanych podagentów jednocześnie |
| Przykładowe wcielenia | Agent UI, Agent Backend, Agent API, Agent Security, Agent Database, Agent Testing — zestaw przykładowy, nie zamknięty; dowolne inne wcielenie jest przedmiotem konfiguracji |
| Oprzyrządowanie podagenta | Każdy podagent może nieść pełną definicję komponentu własnego rodzaju „agent” (model, tożsamość, instrukcje, skille, konektory, pamięć, uprawnienia) albo działać jako wąsko sprofilowany model bazowy dedykowany jednemu zakresowi zadania |
| Mechanizm działania | Wykonawca dzieli złożone zadanie na wąsko zdefiniowane zakresy (`split`), przydziela każdy zakres jednemu podagentowi, a po zakończeniu ich pracy agreguje wyniki (`merge`) |
| Punkty izolacji | Brak odrębnego poziomu zasięgu dla podagenta — Subagent Network dziedziczy profil izolacji roli macierzystej, chyba że użytkownik świadomie ustanowi odrębną konfigurację na poziomie karty sesji |

**Schemat mechanizmu.**

```
                     Executor 1  (albo Executor 2) — wykonawca macierzysty
                              │
                       split — podział złożonego zadania na wąskie zakresy
                              │
      ┌──────────┬───────────┼───────────┬────────────┬────────────┐
      ▼          ▼           ▼           ▼            ▼            ▼
   Agent UI   Agent      Agent API   Agent       Agent        Agent
              Backend                Security    Database     Testing     …
      │          │           │           │            │            │
      │          │      (do 15 wyspecjalizowanych podagentów jednocześnie)
      │          │           │           │            │            │
      └──────────┴───────────┼───────────┴────────────┴────────────┘
                              │
                       merge — agregacja wyników przez wykonawcę macierzystego
                              │
                    rezultat cząstkowy ──► silnik kolejek (rozdz. 4)

   Profil izolacji: dziedziczony z roli macierzystej (rozdz. 13)
```

### 3.3. Coordinator — koordynator

**Zakres odpowiedzialności.**

| Obszar | Zawartość |
|---|---|
| Planowanie | Etapy projektu, zależności, harmonogramy, logika przepływu pracy |
| Podział pracy | Przypisywanie ról, zadań, odpowiedzialności, modeli wykonawczych |
| Budowa promptów | Dynamiczne tworzenie i optymalizacja promptów dla wykonawców |
| Sterowanie procesem | Start, stop, pauza, wznowienie, przekazanie, powtórzenie, walidacja |
| Zarządzanie kolejką | Budowa zaawansowanych struktur kolejkujących |

**Sterowanie procesem → akcje silnika kolejek.**

| Operacja sterująca | Odpowiadająca akcja lub mechanizm |
|---|---|
| Start | Uruchomienie przepływu zadań |
| Stop | Trwałe zakończenie przepływu |
| Pauza | `pause` |
| Wznowienie | `resume` |
| Przekazanie | `route` — skierowanie zadania do innej roli |
| Powtórzenie | `retry` |
| Walidacja | Skierowanie wyniku do Executora 3 / Validatora przed uznaniem etapu za zakończony |

**Pełne oprzyrządowanie Coordinatora.**

| Grupa narzędzi | Zawartość |
|---|---|
| Przypisanie | Agent z modułu Agents lub model bazowy, jak pozostałe role |
| Widok planu | Panel etapów projektu z zależnościami i harmonogramem |
| Kreator promptów | Wewnętrzny mechanizm budowy i optymalizacji promptów przekazywanych wykonawcom, na podstawie celu procesu i uwag Executora 3 / Validatora |
| Akcje silnika kolejek dostępne roli | `enqueue`, `delay`, `retry`, `pause`, `resume`, `route`, `branch`, `condition` — pełny zestaw poza `dequeue`, `split`, `merge`, właściwymi wykonawcom |
| Widok stanu kolejki | Podgląd zawartości kolejek przypisanych procesowi, bez konieczności przełączania się do sekcji Kolejki panelu orkiestracji |
| Okno robocze | Coordinator Chat (rozdz. 9) |

### 3.4. Executor 2 — równoległy wykonawca

Oprzyrządowanie tożsame z Executorem 1 (rozdz. 3.1): ten sam zestaw kanałów modelu, rozszerzeń, poziomów pamięci, akcji `dequeue`/`split`/`merge` oraz własna, niezależna instancja Subagent Network (do 15 podagentów, niezależnie od podagentów uruchomionych przez Executora 1).

**Cztery tryby współpracy z Executorem 1.**

| Tryb | Opis | Kierunek przepływu danych | Przykład zastosowania |
|---|---|---|---|
| Praca niezależna | Executor 1 i Executor 2 realizują odrębne zadania bez wzajemnej zależności; synchronizacja ogranicza się do wspólnej kolejki lub wspólnego celu procesu | Brak przepływu bezpośredniego — tory zbiegają się w Coordinatorze lub Executorze 3 / Validatorze | Backend i frontend budowane równolegle w module Apps |
| Przekazywanie wyników | Wynik pracy jednego wykonawcy stanowi wejście dla drugiego; przepływ jednokierunkowy | Executor 1 → Executor 2 (lub odwrotnie), przez `route` | Badania (Executor 1) przekazane do redakcji raportu (Executor 2) |
| Praca naprzemienna | Wykonawcy przejmują zadanie kolejno, każdy realizując kolejny etap tego samego procesu | Executor 1 → Executor 2 → Executor 1 → …, sterowane przez Coordinatora | Kod i dokumentacja techniczna aktualizowane naprzemiennie |
| Praca iteracyjna | Wykonawcy poprawiają wzajemnie swoje rezultaty w powtarzających się cyklach, aż do osiągnięcia założonej jakości | Dwukierunkowa pętla Executor 1 ⇄ Executor 2, zamykana decyzją Coordinatora lub oceną Executora 3 / Validatora | Iteracyjne dopracowywanie kodu i testów do przejścia walidacji |

```
 Praca niezależna         Przekazywanie wyników      Praca naprzemienna        Praca iteracyjna

  Executor 1               Executor 1                 Executor 1 ──►            Executor 1 ⇄ Executor 2
      │                        │ route                Executor 2 ──►                 │
  Executor 2                   ▼                       Executor 1 ──► …          zamknięcie: Coordinator
      │                   Executor 2                  (sterowane przez               lub Executor 3/Validator
      ▼                                                 Coordinatora)
  Coordinator / Executor 3 (zbiegają się na końcu procesu, niezależnie od trybu)
```

Wybór trybu jest polem szablonu roli i pozostaje w gestii użytkownika lub Coordinatora — nie jest wymuszony i może zostać zmieniony w dowolnym momencie trwania procesu.

### 3.5. Executor 3 / Validator — czwarty model

| Wcielenie | Funkcja | Ikona odniesienia |
|---|---|---|
| Validator | Kontrola jakości pracy pozostałych modeli przed uznaniem etapu za zakończony | `ptaszek-kolo` |
| Reviewer | Recenzja wyników pod kątem poprawności i kompletności | `oko` |
| Security Auditor | Ocena bezpieczeństwa rozwiązania wytworzonego przez wykonawców | `tarcza` |
| Architect | Ocena zgodności rozwiązania z założeniami architektonicznymi | `kod` |
| Product Owner | Ocena zgodności rezultatu z wymaganiami biznesowymi | `dokument` |
| QA Lead | Przygotowanie przypadków testowych i raportów jakości | `ptaszek` |
| Arbitrator | Rozstrzyganie konfliktów powstałych między wynikami Executora 1 i Executora 2 | `waga` |

Konkretne wcielenie roli jest w pełni zależne od charakteru realizowanego procesu i pozostaje przedmiotem konfiguracji — powyższe zestawienie ma charakter przykładowy, nie wyczerpujący. Zgodnie z hierarchią decyzji (rozdz. 8) udział tej roli w procesie jest sterowany ustawieniem konfiguracyjnym: prostszy zespół pomija tę rolę, a użytkownik może w każdej chwili ominąć jej ocenę i zainterweniować bezpośrednio — funkcja kontrolna nie jest bramą blokującą przepływ pracy, lecz konfigurowalnym etapem, który można włączyć, wyłączyć lub pominąć.

**Pełne oprzyrządowanie.**

| Grupa narzędzi | Zawartość |
|---|---|
| Przypisanie | Agent z modułu Agents lub model bazowy |
| Wejście oceny | Rezultat pracy Executora 1 i Executora 2, wynik Subagent Network po agregacji |
| Akcje silnika kolejek dostępne roli | `retry` (zwrot etapu do Coordinatora), udział w `condition`/`branch` jako źródło warunku |
| Narzędzie rozstrzygania konfliktów | Panel porównania wyników dwóch wykonawców, dostępny we wcieleniu Arbitrator |
| Okno robocze | Results Analyzer (rozdz. 9) |

### 3.6. Zestawienie zbiorcze ról

| Rola | Tworzy produkt końcowy | Zarządza procesem | Uruchamia Subagent Network | Przypisanie wykonawcy | Sekcja panelu orkiestracji | Okno robocze |
|---|---|---|---|---|---|---|
| Executor 1 | Tak | Nie | Tak (do 15) | Agent lub model bazowy | Role | Executor Chat |
| Subagent Network | Nie — wynik agreguje wykonawca | Nie | — | Podagenci uruchamiani przez Executora | Role (rozwinięcie karty roli) | — (zagnieżdżone w Executor Chat) |
| Coordinator | Nie | Tak | Nie | Agent lub model bazowy | Role | Coordinator Chat |
| Executor 2 | Tak | Nie | Tak (do 15) | Agent lub model bazowy | Role | Executor Chat |
| Executor 3 / Validator | Zależnie od wcielenia (zwykle nie) | Nie | Nie | Agent lub model bazowy | Role | Results Analyzer |

Warstwą centralną nadrzędną wobec wszystkich pięciu pozycji tabeli jest Always On Display (rozdz. 2).

---

## 4. Silnik kolejek

Silnik kolejek jest mechanizmem technicznym leżącym u podstaw zarządzania kolejką realizowanego przez Coordinatora. Udostępnia operacje, za pomocą których zadania są dodawane, wstrzymywane, dzielone, łączone lub kierowane warunkowo między wykonawcami i podagentami. W warstwie komunikacji jest realizowany poleceniem `queue.action`, przyjmującym jako parametr jedną z jedenastu akcji.

### 4.1. Zasięgi kolejek

| Zasięg | Opis | Przykład zastosowania |
|---|---|---|
| Globalna | Obejmuje wszystkie procesy MultitaskingAI użytkownika | Kolejka wspólna dla wszystkich aktywnych zespołów |
| Lokalna | Obejmuje jeden proces orkiestracji (jedną kartę sesji) | Kolejka etapów pojedynczego projektu realizowanego w module Apps |
| Dla modelu | Przypisana konkretnemu modelowi bazowemu | Kolejka zadań przetwarzanych przez jeden, wskazany model |
| Dla agenta | Przypisana konkretnemu agentowi z modułu Agents | Kolejka zadań przypisanych agentowi pełniącemu rolę Executora 1 |
| Dla projektu | Przypisana projektowi prowadzonemu w module Workspace | Kolejka zadań realizowanych w ramach jednego, izolowanego projektu |

### 4.2. Jedenaście akcji silnika kolejek

| Akcja | Grupa funkcjonalna | Działanie | Typowy inicjator | Reprezentacja w interfejsie |
|---|---|---|---|---|
| `enqueue` | Obieg zadań | Dodanie nowego zadania do wskazanej kolejki | Coordinator, Automations (po integracji) | Przycisk ikonowy w wierszu tabeli kolejki, ikona `plus` |
| `dequeue` | Obieg zadań | Pobranie kolejnego zadania z kolejki do wykonania | Executor 1, Executor 2 | Automatyczne wywołanie przy wolnym slocie wykonawcy; widoczne jako etykieta w Executor Chat |
| `delay` | Sterowanie czasem | Odroczenie wykonania zadania o zadany czas lub do spełnienia warunku | Coordinator | Pole czasu w wierszu zadania, ikona `zegar` |
| `retry` | Ponawianie | Ponowienie zadania po niepowodzeniu wykonania lub po negatywnej ocenie walidacji | Coordinator, Executor 3 / Validator | Przycisk ikonowy, ikona `odswiez` |
| `pause` | Wstrzymywanie | Wstrzymanie przetwarzania wskazanej kolejki | Coordinator, Always On Display (operator) | Przycisk ikonowy, ikona `zatrzymaj` |
| `resume` | Wstrzymywanie | Wznowienie wstrzymanej kolejki | Coordinator, Always On Display (operator) | Przycisk ikonowy, ikona `uruchom` |
| `split` | Podział i scalanie | Podział zadania na mniejsze podzadania | Executor 1 lub Executor 2, przy uruchamianiu Subagent Network | Przycisk „Uruchom Subagent Network” w Executor Chat |
| `merge` | Podział i scalanie | Scalenie wyników wielu zadań lub podzadań w jeden rezultat | Executor 1 lub Executor 2, po agregacji wyników Subagent Network | Przycisk „Scal wyniki” w panelu Subagent Network |
| `route` | Kierowanie warunkowe | Skierowanie zadania do wskazanej roli lub kolejki na podstawie reguły | Coordinator, orkiestracja | Selektor roli docelowej w wierszu zadania, ikona `strzalka-prawo` |
| `branch` | Kierowanie warunkowe | Rozgałęzienie przepływu pracy na alternatywne ścieżki | Coordinator | Kreator reguły rozgałęzienia w sekcji Orkiestracja |
| `condition` | Kierowanie warunkowe | Warunkowe wykonanie kolejnego kroku procesu w zależności od wyniku poprzedniego | Coordinator, orkiestracja | Kreator warunku w sekcji Orkiestracja |

### 4.3. Cykl życia zadania w kolejce

```
Coordinator ── enqueue ──► KOLEJKA (globalna / lokalna / modelu / agenta / projektu)
                                │
                    (pause / resume — wg decyzji Coordinatora lub AOD)
                                │
                          dequeue ──► Executor 1 / Executor 2
                                            │
                              split (wg decyzji wykonawcy) ──► Subagent Network
                                            │                  │
                                            │            (do 15 podagentów)
                                            │                  │
                                            ◄──── merge ───────┘
                                            │
                                     wynik pracy wykonawcy
                                            │
                          route / branch / condition (wg reguł orkiestracji)
                                            │
                              ┌─────────────┴─────────────┐
                              ▼                            ▼
                    Executor 3 / Validator          kolejny etap procesu
                    (ocena jakości — sterowana
                     ustawieniem konfiguracyjnym)
                              │
                    retry (przy niepowodzeniu) ──► z powrotem do enqueue
```

### 4.4. Powiązanie z modułem Automations

Silnik kolejek środowiska MultitaskingAI może zostać powiązany z silnikiem kolejek modułu Automations (Queue Manager) — połączenie nie jest domyślne i wymaga decyzji użytkownika podjętej w oknie konfiguracji oraz w ustawieniach okna modułu Automations (rozdz. 7).

### 4.5. Katalog elementów interfejsu sekcji Kolejki

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Tabela kolejek | Lista wszystkich zdefiniowanych kolejek z ich zasięgiem, priorytetem i liczbą zadań | Przegląd i zarządzanie wszystkimi kolejkami procesu | Duża tabela (`.dn-tabela`) | Domyślna · pusty stan (`.dn-pusty-stan` przed pierwszą konfiguracją) · ładowanie · błąd pobrania danych | Kliknięcie wiersza rozwija szczegóły kolejki i zadania w niej zawarte | Sekcja Kolejki panelu orkiestracji | 1 | Widoczna bez interakcji |
| Selektor zasięgu kolejki | Pole wyboru jednego z pięciu zasięgów | Określenie, do jakiego poziomu należy tworzona lub edytowana kolejka | Pole wyboru (`.dn-select`) | Domyślny (globalna) · rozwinięty · wybrany | Zmiana zasięgu przypisuje kolejkę do wskazanego poziomu | Formularz nowej kolejki, wiersz tabeli kolejek | 2 | Lista rozwijana `Zasięg ▼` |
| Przycisk „Nowa kolejka” | Kontrolka dodania kolejki | Utworzenie nowej definicji kolejki | Przycisk (`.dn-btn--zloty`), ikona `plus` | Domyślny · najechanie · wciśnięty | Otwiera formularz definicji nowej kolejki | Nagłówek sekcji Kolejki | 1 | Widoczny bez interakcji |
| Zestaw przycisków akcji | Jedenaście przycisków ikonowych odpowiadających akcjom silnika kolejek | Wykonanie akcji na wskazanym zadaniu lub kolejce (rozdz. 4.2) | Mały przycisk ikonowy (`.dn-btn-ikona`), zestaw w kolumnie akcji wiersza | Domyślny · ładowanie (akcja w toku) · sukces (krótkie potwierdzenie) | Każdy przycisk pozostaje klikalny niezależnie od stanu zadania i wywołuje odpowiednią akcję `queue.action`; gdy akcja nie pasuje do bieżącego stanu zadania, wyświetla komunikat kontekstowy zamiast się blokować; wynik odświeża wiersz | Kolumna akcji tabeli kolejek | 3 | Rozwinięcie `Operacje ▼` w wierszu zadania |
| Pole priorytetu | Wartość liczbowa porządkująca kolejność obsługi | Ustalenie, które zadanie w kolejce obsłużyć w pierwszej kolejności | Pole liczbowe (`.dn-input`) | Domyślny · edytowany · ostrzeżenie (wartość nietypowa) | Zmiana wartości przelicza kolejność wyświetlania zadań w tabeli; wartość nietypowa jest sygnalizowana ostrzeżeniem, które nie blokuje zapisu | Wiersz zadania w tabeli kolejek | 2 | Kliknięcie wartości w wierszu |
| Selektor obsługi błędów | Pole wyboru zachowania przy niepowodzeniu zadania | Ustalenie domyślnej reakcji kolejki na błąd: `retry`, `route` lub `pause` | Pole wyboru (`.dn-select`) | Domyślny (`retry`) · rozwinięty · wybrany | Zmiana ustawia regułę stosowaną automatycznie przy kolejnym niepowodzeniu | Formularz definicji kolejki | 3 | Rozwinięcie formularza definicji kolejki |
| Kreator reguły warunkowej | Panel formularza budowy warunku `branch` lub `condition` | Zdefiniowanie reguły kierującej przepływ pracy na podstawie wyniku poprzedniego kroku | Panel formularza średniej wielkości | Domyślny (pusty) · w edycji · zapisany · ostrzeżenie walidacji reguły | Zapis dodaje regułę do listy zależności orkiestracji (rozdz. 5); reguła wewnętrznie sprzeczna zostaje zapisana z oznaczeniem ostrzegawczym zamiast odrzucenia zapisu | Sekcja Kolejki i sekcja Orkiestracja | 3 | Menu kebab (⋮) wiersza, sekcja Orkiestracja |

### 4.6. Makieta — sekcja Kolejki panelu orkiestracji

Makieta przedstawia stan spoczynku: widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

```
 Makieta — Panel orkiestracji ▸ Kolejki, stan spoczynku
 ══════════════════════════════════════════════════════════════════════════
  Panel        │ Chat Window     │ Kolejki                    [+ Nowa]  [⋮]
  orkiestracji │ Użytkownik ↔    │ [Zasięg ▼] [szukaj…]
  (boczna      │ Wykonawca       │ ───────────────────────────────────────
   nawigacja)  │                 │ Nazwa           Zasięg      Zadań  Prio
               │ [MultitaskingAI]│ Kolejka główna  lokalna       4      1 ⋮
  ▸ Zespoły    │ [Kolejki]  [⋮]  │ Kolejka QA      dla agenta    2      2 ⋮
  ▸ Role       │                 │ Kolejka Backend dla projektu  6      1 ⋮
  ▸ Kolejki ●  │ strumień        │ ───────────────────────────────────────
  ▸ Orkiestr.  │ poleceń         │ Zadanie #128 · Executor 1 · w toku
  ▸ Harmono-   │ i wyników       │ [ Operacje ▼ ]  enqueue · dequeue ·
    gram i     │                 │                 delay · retry · pause ·
    automatyki │                 │                 resume · split · merge ·
  ▸ Monitor    │                 │                 route · branch · condition
    procesu    │ [📎] …       [➤]│
      [☰]      │                 │
 ══════════════════════════════════════════════════════════════════════════
   Jedenaście akcji silnika kolejek jest zwiniętych w element zbiorczy
   `Operacje ▼` (warstwa 3); rozwinięcie zawiera pełną listę akcji
```

---

## 5. Orkiestracja i zależności

### 5.1. Funkcja orkiestracji

Orkiestracja pozwala definiować zależności pomiędzy modelami, agentami, zadaniami, kolejkami, automatyzacjami oraz projektami. W warstwie komunikacji definicja zależności jest realizowana poleceniem `orchestration.define`. Warstwa orkiestracji nadaje sens funkcjonowaniu poszczególnych ról i kolejek jako spójnej całości — odpowiada za to, że zadanie Executora 2 nie rozpocznie się przed zakończeniem etapu przypisanego Executorowi 1, jeśli taka zależność została ustalona, oraz że wynik pracy trafi do właściwej kolejki lub roli w kolejnym kroku procesu.

### 5.2. Rodzaje zależności

| Rodzaj zależności | Opis | Mechanizm realizujący |
|---|---|---|
| Sekwencyjna | Etap kolejnej roli rozpoczyna się dopiero po zakończeniu etapu roli poprzedzającej | `route` w połączeniu z regułą zależności zdefiniowaną przez Coordinatora |
| Warunkowa | Dalszy przebieg procesu zależy od wyniku poprzedniego etapu (np. oceny Executora 3 / Validatora) | `condition`, `branch` |
| Równoległa niezależna | Dwa lub więcej tory pracy przebiegają bez wzajemnego oczekiwania | Brak zależności zdefiniowanej w orkiestracji — odpowiednik trybu pracy niezależnej (rozdz. 3.4) |

### 5.3. Odpowiedzialność za orkiestrację

Zależności ustala Coordinator w ramach obszaru „Zarządzanie kolejką” i „Planowanie” (rozdz. 3.3); ich respektowanie w toku wykonania zapewnia warstwa orkiestracji działająca w powiązaniu z silnikiem kolejek. Żadna zależność nie jest wbudowana na stałe — każda jest możliwością świadomie ustanowioną przez użytkownika lub Coordinatora, w każdej chwili odwracalną.

### 5.4. Katalog elementów interfejsu sekcji Orkiestracja

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Lista zależności | Zestawienie wszystkich zdefiniowanych zależności procesu | Przegląd reguł porządkujących przepływ pracy między rolami i kolejkami | Lista pozycji (`.dn-karta--pozycja`) | Domyślna · pusty stan · z zaznaczoną pozycją | Kliknięcie pozycji otwiera szczegóły reguły do edycji | Sekcja Orkiestracja panelu orkiestracji | 1 | Widoczna bez interakcji |
| Widok grafu zależności | Wizualna mapa powiązań między rolami, kolejkami i etapami | Szybkie zrozumienie topologii procesu bez czytania listy tekstowej | Duży panel wizualny | Domyślny · powiększony (pełny ekran) · zaznaczony węzeł | Kliknięcie węzła podświetla powiązane zależności i otwiera panel szczegółów | Sekcja Orkiestracja, zakładka „Mapa” | 2 | Zakładka „Mapa” sekcji Orkiestracja |
| Przycisk „Dodaj zależność” | Kontrolka dodania nowej reguły | Utworzenie nowej zależności sekwencyjnej, warunkowej lub równoległej | Przycisk (`.dn-btn--zloty`), ikona `plus` | Domyślny · najechanie · wciśnięty | Otwiera kreator zależności (wybór rodzaju z rozdz. 5.2) | Nagłówek sekcji Orkiestracja | 1 | Widoczny bez interakcji |
| Plakietka statusu zgodności | Wskaźnik, czy bieżący przebieg respektuje wszystkie zdefiniowane zależności | Szybka sygnalizacja konfliktu lub naruszenia reguły | Mała plakietka (`.dn-plakietka--stan`) | Zgodny (sukces) · oczekujący (info) · naruszony (błąd) | Kliknięcie przy stanie „naruszony” otwiera szczegóły konfliktu | Lista zależności, widok grafu | 1 | Widoczna bez interakcji |

### 5.5. Diagram — przykładowa mapa zależności

```
                         ┌────────────┐
                         │ Coordinator│
                         └─────┬──────┘
                sekwencyjna    │    sekwencyjna
              ┌─────────────────┼─────────────────┐
              ▼                                    ▼
       ┌────────────┐                       ┌────────────┐
       │ Executor 1 │  równoległa niezależna │ Executor 2 │
       │ (backend)  │◄───────────────────────►  (frontend)│
       └─────┬──────┘                       └──────┬─────┘
             │                                       │
             └───────────────┬───────────────────────┘
                              ▼
                    ┌───────────────────┐
                    │ Executor 3 /       │  warunkowa (condition)
                    │ Validator          │──────────────┐
                    └─────────┬──────────┘              │
                    ocena pozytywna              ocena negatywna
                              ▼                          ▼
                     kolejny etap procesu         retry → Coordinator
```

---

## 6. Panel orkiestracji — boczna nawigacja środowiska

### 6.1. Zasada organizacji bocznej nawigacji

Trzy środowiska modułowe — TalkIn, WorkSpace i CodeStudio — udostępniają w bocznej nawigacji listę modułów dostępnych w danym środowisku. Środowisko MultitaskingAI organizuje jednak nie zadania modułowe, lecz zespół modeli i agentów realizujących wspólny proces, dlatego jego boczna nawigacja nie zawiera listy modułów, lecz panel orkiestracji — zestaw sekcji odpowiadających własnym warstwom sterowania tego środowiska. Panel orkiestracji jest odpowiednikiem bocznej nawigacji modułów w pozostałych środowiskach: nawiguje po tym, czym w tym środowisku faktycznie się steruje — po rolach, kolejkach i orkiestracji — a nie po modułach.

### 6.2. Sześć sekcji panelu

| Sekcja | Zawartość | Typowe działania | Powiązane mechanizmy i okna | Rozdział niniejszego dokumentu | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Zespoły | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek | Zapis nowego zespołu, wczytanie zapisanego zespołu, duplikowanie, usunięcie | Szablony konfiguracji zespołu | 14 | 1 | Pozycja panelu orkiestracji, widoczna bez interakcji; zestaw akcji zespołu w menu kebab (⋮) — warstwa 3 |
| Role | Cztery okna robocze: Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów; uruchomienie i podgląd Subagent Network | Przypisanie agenta lub modelu bazowego, wybór kanału modelu, ustalenie trybu współpracy, aktywacja Subagent Network | Agent Builder, Permissions Center | 3, 9 | 1 | Pozycja panelu orkiestracji; wybór wykonawcy i trybu współpracy jako znaczniki `▼` — warstwa 2 |
| Kolejki | Definicje kolejek pięciu zasięgów oraz ich akcje | Utworzenie kolejki, przypisanie zasięgu, konfiguracja reguł warunkowych, podgląd zawartości | Silnik kolejek, Queue Manager | 4 | 1 | Pozycja panelu orkiestracji; akcje kolejki w menu kebab (⋮) — warstwa 3 |
| Orkiestracja | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami | Definicja zależności sekwencyjnych, warunkowych, równoległych | Orchestrator | 5 | 2 | Pozycja panelu orkiestracji; edycja reguł zależności — warstwa 4 |
| Harmonogram i automatyki | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations | Ustalenie reguły czasowej i cykliczności, wpięcie lub odpięcie automatyki | Scheduler, Execution Monitor | 7 | 2 | Pozycja panelu orkiestracji; wpięcie automatyki — warstwa 4 |
| Monitor procesu | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display | Podgląd stanu ról i kolejek, historia przebiegów, przełączenie AOD między trybami | Always On Display, Mobile | 2, 8, 12 | 1 | Pozycja panelu orkiestracji, otwierana jako panel wysuwany; przełączenie trybu Always On Display — warstwa 4 |

Zawartość sekcji w skrócie: Zespoły — presety ról, powiązań i kolejek; Role — Executor 1, Executor 2, Coordinator, Executor 3 / Validator wraz z Subagent Network każdego wykonawcy; Kolejki — zasięgi globalny, lokalny, modelu, agenta i projektu; Orkiestracja — zależności sekwencyjne, warunkowe i równoległe; Harmonogram i automatyki — reguła czasowa pracy ciągłej oraz wpięte automatyki; Monitor procesu — statusy przebiegów, hierarchia decyzji i nadzór Always On Display.

### 6.3. Schemat układu panelu

```
┌─────────────────────────┐
│  PANEL ORKIESTRACJI     │
├─────────────────────────┤
│ ▸ Zespoły               │  presety ról, powiązań, kolejek
│ ▸ Role                  │  Executor 1 · Executor 2 · Coordinator · Executor 3/Validator
│                         │  (+ Subagent Network per wykonawca)
│ ▸ Kolejki               │  globalne · lokalne · modelu · agenta · projektu
│ ▸ Orkiestracja          │  zależności sekwencyjne · warunkowe · równoległe
│ ▸ Harmonogram           │  harmonogram pracy ciągłej + automatyki wpięte
│   i automatyki          │
│ ▸ Monitor procesu       │  statusy, hierarchia decyzji, nadzór AOD
└─────────────────────────┘
```

### 6.4. Katalog elementów interfejsu panelu

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pozycja nawigacyjna sekcji | Wiersz listy pionowej z ikoną i etykietą jednej z sześciu sekcji | Przełączenie widoku głównego obszaru roboczego na wskazaną sekcję | Pozycja listy średniej wagi (`.dn-karta--pozycja`) | Domyślna · najechanie · aktywna (zaznaczona, tło `--dn-akcent-tlo`) · z plakietką liczbową (np. liczba oczekujących zadań) | Kliknięcie ładuje zawartość sekcji w głównym obszarze roboczym | Panel orkiestracji, sześć wystąpień | 1 | Widoczna bez interakcji |
| Uchwyt przeciągania kolejności | Mały uchwyt przy pozycji sekcji | Zmiana kolejności wyświetlania sekcji w panelu | Mała ikonka | Domyślny (widoczny na najechanie) · przeciągany | Przeciągnięcie zmienia kolejność zapisaną w konfiguracji panelu | Panel orkiestracji, tryb edycji układu | 3 | Menu kontekstowe panelu, tryb edycji układu |
| Przełącznik widoczności sekcji | Kontrolka ukrycia lub przywrócenia sekcji | Dostosowanie zestawu widocznych sekcji do potrzeb procesu, bez usuwania funkcji | Mały przełącznik (`.dn-suwak`) | Domyślny (widoczna) · ukryta | Przełączenie natychmiast ukrywa lub przywraca pozycję w panelu | Okno konfiguracji panelu orkiestracji | 4 | Okno konfiguracji panelu orkiestracji |
| Przycisk zwinięcia panelu | Kontrolka zwężenia panelu do samych ikon | Zwiększenie przestrzeni głównego obszaru roboczego | Mała ikonka (`menu`) | Domyślny (rozwinięty) · zwinięty | Zwinięcie chowa etykiety tekstowe, pozostawiając same ikony sekcji | Nagłówek panelu orkiestracji | 1 | Ikona `☰` w nagłówku panelu |
| Przycisk „Przywróć domyślne” | Kontrolka resetu kolejności i widoczności | Powrót do zestawu i kolejności sześciu sekcji domyślnych | Przycisk tekstowy (`.dn-btn--duch`) | Domyślny · najechanie | Przywraca kolejność i widoczność ustaloną jako zestaw domyślny | Okno konfiguracji panelu orkiestracji | 4 | Okno konfiguracji panelu orkiestracji |

### 6.5. Konfigurowalność kolejności i widoczności

Kolejność i widoczność sekcji panelu orkiestracji podlegają konfiguracji z poziomu okna konfiguracji; przedstawiony w rozdziale 6.2 zestaw sześciu sekcji jest zestawem domyślnym, obowiązującym przy braku odmiennego ustawienia. Panel nigdy nie wymusza ukrycia żadnej sekcji ani nie blokuje możliwości przywrócenia zestawu domyślnego — brak ustawienia oznacza wartość domyślną, nigdy brak dostępności.

```
panel_orkiestracji:
  sekcje:
    - nazwa: "Zespoły"                    widoczna: tak   kolejnosc: 1
    - nazwa: "Role"                       widoczna: tak   kolejnosc: 2
    - nazwa: "Kolejki"                    widoczna: tak   kolejnosc: 3
    - nazwa: "Orkiestracja"               widoczna: tak   kolejnosc: 4
    - nazwa: "Harmonogram i automatyki"   widoczna: tak   kolejnosc: 5
    - nazwa: "Monitor procesu"            widoczna: tak   kolejnosc: 6
  przywrocenie_domyslnych: dostepne         # panel nigdy nie wymusza ukrycia sekcji
```

---

## 7. Integracja z modułem Automations — praca ciągła 24/7/365

### 7.1. Warunek uruchomienia

Po skonfigurowaniu połączenia między środowiskiem MultitaskingAI a modułem Automations — decyzją użytkownika podjętą w oknie konfiguracji oraz w ustawieniach okna modułu Automations — możliwe jest zbudowanie pełnego autonomicznego systemu realizacji projektów, zdolnego do działania w sposób ciągły. Integracja nie jest domyślna: platforma jej nie wymusza, lecz udostępnia jako możliwość skonfigurowania w dowolnym momencie.

### 7.2. Mapowanie mechanizmów

| Mechanizm MultitaskingAI | Okno modułu Automations | Funkcja po spięciu |
|---|---|---|
| Silnik kolejek (rozdz. 4) | Queue Manager | Trwałe zarządzanie kolejkami zadań poza pojedynczą sesją |
| Harmonogram pracy ciągłej (sekcja Harmonogram i automatyki) | Scheduler | Cykliczne uruchamianie procesu według reguły czasowej |
| Orkiestracja (rozdz. 5) | Orchestrator | Trwałe egzekwowanie zależności między kolejnymi przebiegami procesu |
| Monitor procesu (rozdz. 8, 12) | Execution Monitor | Status każdego przebiegu oraz sygnalizacja błędów wymagających interwencji |

### 7.3. Schemat spięcia

```
   ŚRODOWISKO MultitaskingAI                        MODUŁ Automations
   ──────────────────────────                       ─────────────────
   Silnik kolejek (rozdz. 4)        ──── spięcie ──► Queue Manager
   Harmonogram pracy ciągłej        ──── spięcie ──► Scheduler
   Orkiestracja (rozdz. 5)          ──── spięcie ──► Orchestrator
   Monitor procesu (rozdz. 8, 12)   ──── spięcie ──► Execution Monitor
                                          │
                                  praca ciągła 24/7/365
                                          │
                  Always On Display ◄── interwencja zdalna ──► Mobile
                            (rozdz. 2.4, rozdz. 7.5)
```

### 7.4. Efekt integracji — podział odpowiedzialności w pętli ciągłej

| Uczestnik pętli ciągłej | Odpowiedzialność |
|---|---|
| Coordinator | Planowanie procesu |
| Moduł Automations | Zarządzanie harmonogramami i kolejkami |
| Executor 1, Executor 2 (wraz z ich Subagent Network) | Realizacja zadań |
| Executor 3 / Validator | Kontrola jakości |
| Always On Display | Nadzór całości z perspektywy użytkownika — w trybie obserwatora lub operatora |

### 7.5. Rola funkcji Mobile w pracy ciągłej

Praca w trybie 24/7/365 z natury przebiega poza stałą obecnością użytkownika przy stanowisku roboczym. Funkcja Mobile udostępnia w takim przypadku kanał zatwierdzania, wstrzymywania i modyfikowania uruchomionych procesów z dowolnego miejsca — mechanizm komplementarny wobec Always On Display działającego w trybie operatora (rozdz. 2.4).

### 7.6. Katalog elementów interfejsu integracji

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Przełącznik „Powiąż z Automations” | Kontrolka spięcia silnika kolejek środowiska z Queue Managerem modułu Automations | Włączenie operacyjnego zaplecza harmonogramów i trwałego wykonania cyklicznego | Duży przełącznik (`.dn-suwak`) z etykietą | Domyślny (wyłączony) · włączony · ładowanie (w trakcie spinania) | Włączenie otwiera panel wyboru automatyki i harmonogramu | Sekcja Harmonogram i automatyki | 2 | Sekcja Harmonogram i automatyki |
| Plakietka statusu połączenia | Wskaźnik, czy integracja jest aktywna | Szybka sygnalizacja, czy proces działa w trybie ciągłym | Mała plakietka (`.dn-plakietka--stan`) | Niespięte (neutralny) · spięte (sukces) · błąd połączenia | Kliknięcie otwiera szczegóły połączenia z modułem Automations | Sekcja Harmonogram i automatyki, pasek stanu procesu | 1 | Widoczna bez interakcji |
| Pole harmonogramu | Formularz reguły czasowej i cykliczności | Ustalenie, kiedy i jak często proces ma być uruchamiany ponownie | Pole formularza (`.dn-pole` + selektor cykliczności) | Domyślny (brak harmonogramu) · skonfigurowany · ostrzeżenie (reguła sprzeczna) | Zapis aktywuje regułę w powiązanym Scheduler; reguła sprzeczna zostaje zapisana z ostrzeżeniem zamiast blokady zapisu | Sekcja Harmonogram i automatyki | 3 | Rozwinięcie panelu po włączeniu powiązania |
| Lista wpiętych automatyk | Zestawienie automatyk z modułu Automations powiązanych z bieżącym procesem | Przegląd i zarządzanie automatykami zasilającymi pętlę ciągłą | Lista pozycji (`.dn-karta--pozycja`) | Domyślna · pusty stan · z pozycją aktywną | Kliknięcie pozycji otwiera szczegóły automatyki | Sekcja Harmonogram i automatyki | 1 | Widoczna bez interakcji |
| Przycisk zatwierdzenia zdalnego | Kontrolka dostępna z poziomu funkcji Mobile | Zatwierdzenie, wstrzymanie lub modyfikacja procesu spoza stanowiska roboczego | Przycisk pełnej szerokości (kontekst mobilny) | Domyślny · ładowanie · potwierdzony | Wykonanie akcji natychmiast odzwierciedla się w Monitorze procesu | Widok Mobile procesu MultitaskingAI | 1 | Widoczny bez interakcji w widoku Mobile |
| Powiadomienie push interwencji | Komunikat wysyłany na urządzenie mobilne | Poinformowanie użytkownika o zdarzeniu wymagającym decyzji (np. powtarzające się niepowodzenie walidacji) | Toast / powiadomienie systemowe | Wysłane · odczytane · z akcją szybkiej odpowiedzi | Dotknięcie otwiera widok Mobile z kontekstem zdarzenia | Urządzenie mobilne, inicjowane przez Always On Display w trybie operatora | 1 | Dostarczane na urządzenie mobilne |

---

## 8. Hierarchia decyzji

### 8.1. Poziomy hierarchii

| Poziom | Podmiot | Zakres decyzji |
|---|---|---|
| 1 (najwyższy) | Użytkownik (Operator) | Decyzja nadrzędna wobec każdego elementu procesu; może w dowolnym momencie zmienić konfigurację, zatrzymać proces lub przejąć bezpośrednie sterowanie — żaden poziom niższy nie ogranicza tego uprawnienia |
| 2 | Always On Display (tryb operatora) | Interwencja w przebieg procesu z dowolnego miejsca, w imieniu i pod nadzorem użytkownika |
| 3 | Coordinator | Planowanie, podział pracy, sterowanie procesem, zarządzanie kolejką |
| 4 | Executor 3 / Validator | Ocena jakości wyników wykonawców; zdolność zwrotu etapu do Coordinatora (`retry`); we wcieleniu Arbitrator — rozstrzyganie konfliktów między wynikami Executora 1 i Executora 2 |
| 5 | Executor 1 / Executor 2 | Realizacja przydzielonych zadań w ramach ustalonego trybu współpracy; brak uprawnień do zmiany planu procesu |
| 6 (najniższy) | Subagent Network | Realizacja wąsko zdefiniowanych podzadań przydzielonych przez wykonawcę macierzystego; wynik podlega agregacji i ocenie na poziomach wyższych |

### 8.2. Rozstrzyganie konfliktów

Konflikt między wynikami Executora 1 i Executora 2 — możliwy w trybie pracy iteracyjnej lub naprzemiennej — rozstrzyga rola Executora 3 / Validatora we wcieleniu Arbitrator. Rozstrzygnięcie Arbitratora trafia do Coordinatora, który decyduje o dalszym przebiegu procesu: kontynuacja, powtórzenie etapu (`retry`) lub eskalacja do użytkownika.

### 8.3. Schemat hierarchii

```
Użytkownik (Operator)                          ── decyzja nadrzędna, bez ograniczeń
        │
Always On Display (tryb operatora)             ── interwencja z dowolnego miejsca (Mobile)
        │
    Coordinator                                ── planowanie i sterowanie procesem
        │
Executor 3 / Validator (w tym Arbitrator)       ── kontrola jakości, rozstrzyganie konfliktów
        │
   ┌────┴────┐
Executor 1  Executor 2                          ── realizacja pracy
   │            │
Subagent    Subagent
Network     Network                             ── realizacja podzadań (do 15 podagentów każdy)
```

### 8.4. Konfigurowalność hierarchii

Hierarchia wynika z domyślnego podziału odpowiedzialności ról, a nie ze sztywnej reguły platformy — użytkownik może w każdej chwili pominąć dowolny poziom pośredni i zainterweniować bezpośrednio, a skład ról uczestniczących w procesie (na przykład pominięcie Executora 3 / Validatora w prostszych zespołach) jest przedmiotem konfiguracji (rozdz. 14). Żaden poziom hierarchii nie jest bramą blokującą — jest domyślnym porządkiem, który ustępuje przed decyzją poziomu wyższego w dowolnej chwili.

---

## 9. Okna środowiska — makiety tekstowe

### 9.1. Anatomia wspólna okna środowiska

Każde okno środowiska — Chat Window, Execution Loop Window, cztery okna robocze ról oraz okna pomocnicze — zajmuje własną kolumnę obszaru roboczego i podlega tej samej anatomii co pozostałe okna operacyjne platformy. Kolumny sąsiadują ze sobą poziomo; regulacji podlega wyłącznie ich szerokość.

```
 ┌─ Kolumna okna ──────────────────────────────┐
 │ Nagłówek: nazwa okna · znaczniki kontekstowe │
 │           [Agent ▼] [Kanał ▼]           [⋮] │
 ├──────────────────────────────────────────────┤
 │ OBSZAR ZAWARTOŚCI                            │
 │   strumień komunikatów, plan, kolejka        │
 │   albo panel oceny — zależnie od okna;       │
 │   aktualizacja na żywo kanałem WebSocket     │
 │                                              │
 │                                              │
 ├──────────────────────────────────────────────┤
 │ Pole wprowadzania / sterowanie [Operacje ▼]  │
 └──────────────────────────────────────────────┘
   Stan procesu sesji utrzymywany po stronie
   serwera — rozłączenie klienta nie zamyka okna
```

### 9.2. Tabela okien środowiska

| Okno | Kanał / rola | Cel | Zawartość | Rozmieszczenie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| **Chat Window** | Użytkownik ↔ Wykonawca | Główne okno komunikacji, centralny punkt pracy i podstawowy mechanizm sterowania wszystkimi procesami środowiska | Polecenia w języku naturalnym, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań, wyjaśnianie wyniku i kontekstu | Lewa kolumna, stała, pełna wysokość obszaru roboczego; obecne w każdej przestrzeni roboczej środowiska | 1 | Widoczne bez interakcji |
| **Execution Loop Window** | Koordynator ↔ Wykonawca | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór, orkiestracja i kontrola realizacji procesów przy pracy wielowątkowej | Bieżące zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, komunikaty sterujące, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli, sterowanie przebiegiem | Kolumna sąsiadująca z Chat Window | 1 | Widoczne bez interakcji |
| Executor Chat (Executor 1) | Executor 1 | Realizacja pracy zleconej przez Koordynatora lub użytkownika | Historia poleceń i wyników; dostęp do pamięci projektu; panel Subagent Network | Kolumna obszaru roboczego ról | 1 | Widoczne bez interakcji |
| Executor Chat (Executor 2) | Executor 2 | Równoległa realizacja drugiego toru pracy | Jak Executor 1, w trybie współpracy ustalonym z Executorem 1 | Kolumna obszaru roboczego ról | 1 | Widoczne bez interakcji |
| Coordinator Chat | Coordinator | Planowanie i sterowanie procesem | Plan etapów, przypisania zadań, budowane prompty, stan kolejki | Kolumna obszaru roboczego ról | 1 | Widoczne bez interakcji |
| Results Analyzer | Executor 3 / Validator | Kontrola jakości, zgodności lub bezpieczeństwa pracy pozostałych ról | Wyniki pracy Executora 1 i Executora 2 przekazane do oceny; ustalenia kontrolne | Kolumna obszaru roboczego ról | 1 | Widoczne bez interakcji |
| Panel Subagent Network | Executor 1 lub Executor 2 | Podgląd i sterowanie podagentami uruchomionymi przez wykonawcę | Karty do 15 podagentów, paski postępu, akcje `split` i `merge` | Kolumna boczna, otwierana jako rozszerzenie boczne okna wykonawcy | 3 | Rozwinięcie `Subagent Network ▼` w oknie wykonawcy |
| Okno konfiguracji punktów izolacji | Wszystkie role | Ustalenie profilu izolacji kontekstu i izolacji technicznej dla wskazanej roli | Selektor zasięgu, macierz izolacji, profil i podgląd polityki efektywnej | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 | Znacznik profilu izolacji w nagłówku okna roli, paleta poleceń albo polecenie w Chat Window |

Zachowanie i powiązane operacje okien ról:

| Okno | Zachowanie | Powiązane operacje |
|---|---|---|
| Executor Chat (Executor 1) | Wykonuje zadania pobierane z kolejki; nie zarządza procesem | Wykonanie zadania, uruchomienie podagentów (do 15), zwrot wyniku |
| Executor Chat (Executor 2) | Tryb współpracy sterowany ustawieniem konfiguracyjnym (rozdz. 3.4) | Wykonanie zadania równoległego, wymiana wyników z Executorem 1 |
| Coordinator Chat | Nie tworzy produktu końcowego; steruje: start, stop, pauza, wznowienie, przekazanie, powtórzenie, walidacja | Planowanie etapów, podział pracy, budowa promptu, sterowanie kolejką |
| Results Analyzer | Wcielenie roli sterowane ustawieniem konfiguracyjnym (rozdz. 3.5) | Ocena wyniku, zgłoszenie niezgodności, rozstrzygnięcie konfliktu |
| Execution Loop Window | Prezentuje pętlę wykonawczą niezależnie od tego, które okno roli jest wiodące; sterowanie przebiegiem działa na całym procesie | `pause`, `resume`, przerwanie, korekta zlecenia, ponowienie zadania (`retry`), podgląd dekompozycji zlecenia |

### 9.3. Makieta główna — pełny układ okna środowiska MultitaskingAI

Makieta przedstawia interfejs w stanie spoczynku: widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3. Rozmieszczenie jest funkcjonalne, nie jest projektem graficznym.

```
 Makieta — Środowisko MultitaskingAI, stan spoczynku
 ══════════════════════════════════════════════════════════════════════════════
  Panel        │ Chat Window       │ Execution Loop    │ Obszar roboczy ról
  orkiestracji │ Użytkownik ↔      │ Koordynator ↔     │
  (boczna      │ Wykonawca         │ Wykonawca         │ Coordinator Chat ⋮
   nawigacja)  │                   │                   │  plan etapów
               │ [MultitaskingAI]  │ Zlecenie: budowa  │  ● uruchomiony
  ▸ Zespoły    │ [Zespół „Budowa   │ aplikacji         │
  ▸ Role       │  aplikacji"]      │ ▸ dekompozycja ▼  │ Executor Chat 1  ⋮
  ▸ Kolejki    │ [Executor 1] [⋮]  │ Zadania: 4        │  ● pracuje
  ▸ Orkiestr.  │                   │ ▸ Kolejka ▼       │  ▸ Subagent Network ▼
  ▸ Harmono-   │ strumień poleceń  │ ▸ Kontrola ▼      │
    gram i     │ i wyników         │                   │ Executor Chat 2  ⋮
    automatyki │                   │ pętla: cykl 3     │  ● pracuje
  ▸ Monitor    │                   │ ▸ Sterowanie ▼    │  ▸ Subagent Network ▼
    procesu    │                   │                   │
      [☰]      │ [📎] polecenie… ➤ │ [⏸] [▶] [↻] [⋮]  │ Results Analyzer ⋮
               │                   │                   │  ○ oczekuje na wynik
 ══════════════════════════════════════════════════════════════════════════════
   Znaczniki kontekstowe otwierają selektory (warstwa 2); ▼ i ⋮ rozwijają
   zestawy akcji (warstwa 3); ◐ Always On Display pozostaje dostępny nad
   kolumnami obszaru roboczego
```

### 9.4. Makieta — Chat Window

```
 Makieta — Chat Window (Użytkownik ↔ Wykonawca), stan spoczynku
 ┌──────────────────────────────────────────────┐
 │ Chat Window · [MultitaskingAI] [Zespół „Budo-│
 │ wa aplikacji"] [Executor 1] [Ultra]      [⋮] │
 ├──────────────────────────────────────────────┤
 │ Użytkownik: zbuduj API płatności i formularz │
 │ Wykonawca:  przyjęto; zlecenie przekazane do │
 │             pętli wykonawczej                │
 │ Wykonawca:  etap „Architektura" zakończony   │
 │             [ Pokaż wynik ] [ Wyjaśnij ]     │
 │                                              │
 ├──────────────────────────────────────────────┤
 │ [📎]  Wpisz polecenie…                    [➤]│
 └──────────────────────────────────────────────┘
   Lewa kolumna, stała, pełna wysokość obszaru
   roboczego; obecna w każdej przestrzeni
   roboczej środowiska
```

### 9.5. Makieta — Execution Loop Window

```
 Makieta — Execution Loop Window (Koordynator ↔ Wykonawca), stan spoczynku
 ┌──────────────────────────────────────────────┐
 │ Execution Loop · [Proces: Budowa aplikacji]  │
 │ pętla: cykl 3 · ● uruchomiona            [⋮] │
 ├──────────────────────────────────────────────┤
 │ Zlecenie: budowa aplikacji płatniczej        │
 │ ▸ Dekompozycja na zadania ▼                  │
 │   #128 API płatności      Executor 1 ● w toku│
 │   #129 Formularz UI       Executor 2 ● w toku│
 │   #130 Walidacja bezp.    Validator  ○ czeka │
 ├──────────────────────────────────────────────┤
 │ Koordynator → Executor 1: prompt etapu 2     │
 │ Executor 1 → Koordynator: wynik cząstkowy    │
 │ Koordynator → Validator: skierowanie (route) │
 ├──────────────────────────────────────────────┤
 │ ▸ Kontrola jakości ▼   ▸ Kolejka ▼           │
 │ [⏸ wstrzymaj] [▶ wznów] [↻ ponów] [⋮]        │
 └──────────────────────────────────────────────┘
   Kolumna sąsiadująca z Chat Window
```

### 9.6. Makieta — Executor Chat

```
 Makieta — Okno robocze: Executor Chat (Executor 1 lub Executor 2)
 ┌──────────────────────────────────────────────┐
 │ Executor 1 · [Agent Backend ▼] [Kanał: CLI ▼]│
 │ [Izolacja: dziedziczona]                 [⋮] │
 ├──────────────────────────────────────────────┤
 │ [pobrano z kolejki] Zadanie #128 — API       │
 │ Executor 1: analizuję specyfikację…          │
 │ Executor 1: dzielę zadanie na 3 zakresy      │
 ├──────────────────────────────────────────────┤
 │ ▸ Subagent Network (4/15) ▼                  │
 ├──────────────────────────────────────────────┤
 │ [📎]  Wpisz polecenie dla wykonawcy…      [➤]│
 └──────────────────────────────────────────────┘
   Panel Subagent Network otwiera się jako
   kolumna boczna — rozszerzenie boczne okna
```

### 9.7. Makieta — panel Subagent Network jako rozszerzenie boczne

```
 Makieta — Executor Chat + panel Subagent Network (warstwa 3, rozwinięty)
 ┌───────────────────────────┬──────────────────────────┐
 │ Executor Chat (Executor 1)│ Subagent Network  4/15 ⋮ │
 │                           │ Agent Backend  ████░ 80% │
 │ historia poleceń i wyników│ Agent API      █████ 100%│
 │                           │ Agent Database ██░░░ 30% │
 │                           │ Agent Security ███░░ 55% │
 │                           │                          │
 │ [📎] polecenie…        [➤]│ [ + podagent ] [ Scal ▼ ]│
 └───────────────────────────┴──────────────────────────┘
```

### 9.8. Makieta — Coordinator Chat

```
 Makieta — Okno robocze: Coordinator Chat
 ┌──────────────────────────────────────────────┐
 │ Coordinator · [Agent Project Lead ▼]         │
 │ [Kanał: API ▼]                           [⋮] │
 ├──────────────────────────────────────────────┤
 │ Plan etapów                                  │
 │  1. ✓ Architektura           Executor 1      │
 │  2. ● Backend                Executor 1      │
 │  3. ● Frontend               Executor 2      │
 │  4. ○ Walidacja bezpieczeństwa Executor 3    │
 │  5. ○ Wdrożenie                              │
 ├──────────────────────────────────────────────┤
 │ ▸ Kreator promptu ▼    Kolejka lokalna: 4    │
 ├──────────────────────────────────────────────┤
 │ [ Sterowanie ▼ ]  Start · Stop · Pauza ·     │
 │                   Wznów · Przekaż · Powtórz ·│
 │                   Waliduj                    │
 └──────────────────────────────────────────────┘
```

### 9.9. Makieta — Results Analyzer

```
 Makieta — Okno robocze: Results Analyzer (Executor 3 / Validator)
 ┌──────────────────────────────────────────────┐
 │ Executor 3 / Validator                       │
 │ [Wcielenie: Security Auditor ▼]          [⋮] │
 ├──────────────────────────────────────────────┤
 │ Wynik przekazany do oceny: Zadanie #128      │
 │ kryteria wcielenia · ustalenia kontrolne     │
 │ ▸ Porównanie wyników ▼  (przy rozbieżności)  │
 ├──────────────────────────────────────────────┤
 │ Uzasadnienie oceny…                          │
 ├──────────────────────────────────────────────┤
 │ [ ✓ Zatwierdź ]  [ ✕ Odrzuć ]  [ ⚖ ▼ ]       │
 └──────────────────────────────────────────────┘
   Panel porównania wyników rozwija się jako
   kolumna boczna z dwiema kolumnami wyników
```

### 9.10. Makieta — okno konfiguracji punktów izolacji, zasięg „Rola"

Pełny opis w rozdziale 13; poniżej forma okna, tożsama z jego zastosowaniem globalnym, z zasięgiem ustawionym na „Rola". Okno należy do warstwy 4 i otwiera się jako rozszerzenie boczne.

```
 Makieta — Okno konfiguracji punktów izolacji, zasięg: Rola (Executor 1)
 ┌────────────────┬──────────────────────────────┬─────────────────────┐
 │ SELEKTOR       │ MACIERZ IZOLACJI             │ PROFIL I PODGLĄD    │
 │ ZASIĘGU        │                              │                     │
 │  globalny      │ Izolacja kontekstu:          │ [ Zapisz profil ]   │
 │  środowisko    │  historia   [współdz.|odr.]  │ [ Wczytaj profil ]  │
 │  moduł         │  pamięć     [współdz.|odr.]  │ [ Przypisz do roli ]│
 │  para modułów  │  kontekst   [współdz.|odr.]  │                     │
 │  projekt       │                              │ Warstwa:            │
 │  karta sesji   │ Izolacja techniczna (8):     │  ( ) domyślna       │
 │▸ rola          │  katalog roboczy   [wł|wył]  │  ( ) sesji          │
 │  (Executor 1)  │  środowisko proc.  [wł|wył]  │                     │
 │                │  dane/konfig.modelu[wł|wył]  │ Polityka efektywna: │
 │                │  dostęp sieciowy   [wł|wył]  │  (podgląd reguł     │
 │                │  odczyt/zapis plik.[wł|wył]  │   dziedziczonych)   │
 │                │  konto i token     [wł|wył]  │                     │
 │                │  model procesu     [wł|wył]  │ [?] objaśnienia     │
 │                │  serwer wykonania  [wł|wył]  │     kontekstowe     │
 └────────────────┴──────────────────────────────┴─────────────────────┘
   stan wyjściowy wszystkich przełączników: WYŁĄCZONY (pełny dostęp)
```

---

## 10. Katalog elementów interfejsu

Poniższy katalog zbiera w jednym miejscu każdy pojedynczy element interfejsu wymieniony w rozdziałach 2–9, uzupełniony o elementy wspólne środowiska niewymienione osobno wcześniej. Konwencja kolumn jest jednolita dla całego katalogu.

### 10.1. Elementy globalne środowiska

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pasek kontekstu procesu | Listwa kontekstu nad kolumnami obszaru roboczego, złożona ze znaczników kontekstowych i kart sesji | Identyfikacja platformy i procesu, przełączanie kart sesji, dostęp do selektorów kontekstu | Listwa (`.dn-pasek`) ze znacznikami `[Danaco Console] [MultitaskingAI] [Zespół] [Rola] [Model]` | Domyślny · ze znacznikiem aktywnym | Kliknięcie znacznika otwiera odpowiedni selektor; listwa nie zmienia stanu przy nawigacji wewnętrznej | Nad kolumnami obszaru roboczego każdej karty sesji | 1 | Widoczny bez interakcji |
| Karta sesji | Zakładka pozioma reprezentująca jeden proces orkiestracji | Przełączanie między równolegle prowadzonymi zespołami/procesami | Zakładka (`.dn-zakladki`) | Domyślna · aktywna · z aktywnością w tle (plakietka) · zamykana | Kliknięcie przełącza widoczny proces; własny układ, historia i kontekst każdej karty | Pasek kontekstu procesu | 1 | Widoczna bez interakcji |
| Przycisk „+” nowej karty sesji | Ikona dodania karty | Otwarcie nowego procesu orkiestracji równolegle do istniejących | Mała ikonka (`plus`) | Domyślny · najechanie | Otwiera nową, pustą kartę sesji środowiska MultitaskingAI | Pasek kontekstu procesu, przy ostatniej karcie | 1 | Widoczny bez interakcji |
| Przycisk zamknięcia karty | Ikona „x” na karcie sesji | Zamknięcie procesu orkiestracji w danej karcie | Mała ikonka (`zamknij`) | Domyślny (widoczny na najechanie) | Zamyka kartę; stan procesu po stronie serwera zachowany zgodnie z trwałością sesji | Każda karta sesji | 2 | Ujawniany najechaniem na kartę |
| Przełącznik motywu jasny/ciemny | Ikona słońca/księżyca | Zmiana motywu wizualnego całej platformy | Mała ikonka (`slonce`/`ksiezyc`) | Jasny · ciemny | Przełącza zestaw tokenów kolorów bez przeładowania okna | Listwa ustawień, pasek kontekstu procesu | 3 | Menu kebab (⋮) paska kontekstu procesu |
| Pole wyszukiwania | Pole tekstowe z ikoną lupy | Szybkie odnalezienie roli, zadania lub zależności w bieżącym procesie | Pole (`.dn-input`) z ikoną `szukaj` | Domyślne (puste) · z treścią · brak wyników | Filtruje widoczne listy w otwartej sekcji panelu orkiestracji | Pasek kontekstu procesu, panele list | 2 | Ikona `szukaj`, skrót klawiszowy |

### 10.2. Elementy panelu orkiestracji

Pełny katalog elementów panelu — pozycje nawigacyjne sześciu sekcji, uchwyt przeciągania, przełącznik widoczności, przycisk zwinięcia — przedstawia rozdział 6.4.

### 10.3. Elementy sekcji Zespoły

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Lista zapisanych zespołów | Zestawienie presetów konfiguracji zespołu | Przegląd i wybór gotowej konfiguracji ról, powiązań i kolejek | Lista kart (`.dn-karta--pozycja`) | Domyślna · pusty stan · z zaznaczoną pozycją | Kliknięcie wczytuje podgląd konfiguracji zespołu | Sekcja Zespoły | 1 | Widoczna bez interakcji |
| Przycisk „Zapisz zespół” | Kontrolka zapisu bieżącej konfiguracji jako presetu | Utrwalenie bieżącego układu ról, trybów współpracy i kolejek do ponownego użycia | Przycisk (`.dn-btn--zloty`) | Domyślny · ładowanie (zapis w toku) · zapisany (potwierdzenie) | Otwiera pole nazwy presetu, następnie dodaje go do listy | Nagłówek sekcji Zespoły | 1 | Widoczny bez interakcji |
| Przycisk „Wczytaj zespół” | Kontrolka zastosowania zapisanego presetu | Ustawienie ról, trybów i kolejek zgodnie z wybranym presetem | Przycisk (`.dn-btn--zarys`) | Domyślny · ładowanie | Nadpisuje bieżącą konfigurację ról wybranym presetem od razu; modal potwierdzenia przed nadpisaniem jest ustawieniem konfiguracyjnym włączanym przez Operatora, domyślnie wyłączonym — po wykonaniu dostępne jest „Cofnij” | Wiersz listy zespołów | 2 | Akcja w wierszu listy zespołów |
| Przycisk „Duplikuj” | Ikona kopiowania presetu | Utworzenie nowej, edytowalnej kopii istniejącego zespołu | Mała ikonka | Domyślny · najechanie | Dodaje do listy kopię presetu z sufiksem „(kopia)” | Wiersz listy zespołów, menu kontekstowe | 3 | Menu kebab (⋮) wiersza |
| Przycisk „Usuń” | Ikona kosza | Usunięcie zapisanego presetu zespołu | Mała ikonka (`kosz`) | Domyślny · ładowanie · usunięty (z „Cofnij”) | Usuwa pozycję z listy od razu, z dostępnym „Cofnij” w toaście po wykonaniu; modal potwierdzenia przed usunięciem jest ustawieniem konfiguracyjnym włączanym przez Operatora, domyślnie wyłączonym | Wiersz listy zespołów, menu kontekstowe | 3 | Menu kebab (⋮) wiersza |

### 10.4. Elementy okien ról — wspólne i swoiste

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Karta roli | Kontener reprezentujący jedną z czterech ról w widoku sekcji Role | Podgląd i wejście do okna roboczego danej roli | Duża karta (`.dn-karta`) | Domyślna · aktywna (rola pracuje) · pusta (brak przypisania) · błąd (niedostępny kanał modelu) | Kliknięcie otwiera pełne okno robocze roli | Sekcja Role panelu orkiestracji, makieta główna | 1 | Widoczna bez interakcji |
| Awatar roli | Element identyfikacji wizualnej przypisanego wykonawcy | Rozróżnienie, czy rolę pełni agent, czy model bazowy | Awatar (`.dn-awatar`), wariant kwadratowy dla agenta, kołowy dla modelu bazowego | Domyślny · ze wskaźnikiem statusu (`-stan`) | Najechanie pokazuje nazwę i tożsamość przypisanego wykonawcy | Karta roli, nagłówek okna roboczego | 1 | Widoczny bez interakcji |
| Selektor przypisania wykonawcy | Pole wyboru agenta z modułu Agents lub modelu bazowego | Ustalenie, kto faktycznie pełni daną rolę | Pole wyboru (`.dn-select`) | Domyślny (brak przypisania) · wybrany agent · wybrany model bazowy | Zmiana natychmiast podmienia definicję działającą w roli | Karta roli, nagłówek okna roboczego | 2 | Znacznik kontekstowy `Agent ▼` |
| Selektor kanału modelu | Pole wyboru jednego z czterech kanałów | Ustalenie sposobu połączenia z modelem dla danej roli, nadpisujące kanał domyślny agenta | Pole wyboru (`.dn-select`) | Domyślny (dziedziczony z definicji agenta) · nadpisany | Zmiana kanału wymaga ponownego uruchomienia procesu roli | Nagłówek okna roboczego roli | 2 | Znacznik kontekstowy `Kanał ▼` |
| Plakietka statusu roli | Wskaźnik bieżącego stanu wykonania | Szybka orientacja, czy rola pracuje, oczekuje, czy zgłosiła błąd | Plakietka + kropka (`.dn-plakietka--stan`, `.dn-kropka`) | Bezczynna · pracuje · oczekuje na zależność · błąd | Kliknięcie otwiera log ostatniej operacji roli | Karta roli, makieta główna, pasek stanu | 1 | Widoczna bez interakcji |
| Wskaźnik dostępu do pamięci projektu | Mała ikona sygnalizująca aktywne korzystanie z pamięci | Informacja, czy rola czerpie z pamięci kontekstowej projektu | Mała ikonka | Aktywny (korzysta) · nieaktywny | Kliknięcie otwiera podgląd wpisów pamięci wykorzystywanych przez rolę | Nagłówek okna roboczego roli | 2 | Znacznik kontekstowy w nagłówku okna roli |
| Plakietka profilu izolacji roli | Etykieta wskazująca aktywny profil izolacji | Informacja, czy rola działa na profilu domyślnym, czy na profilu dedykowanym | Plakietka (`.dn-plakietka`) | Domyślny (dziedziczony) · profil dedykowany | Kliknięcie otwiera okno konfiguracji punktów izolacji z zasięgiem „Rola” (rozdz. 13) | Nagłówek okna roboczego roli | 2 | Znacznik kontekstowy w nagłówku okna roli |
| Obszar historii poleceń i wyników | Główny obszar treści okna roboczego wykonawcy | Prezentacja przebiegu pracy: polecenia, wyniki pośrednie, status | Duży obszar zawartości | Pusty (przed pierwszym zadaniem) · aktywny strumień (WebSocket) · przewijalny | Przewijanie ujawnia wcześniejsze etapy; nowe wpisy dopisywane na żywo | Executor Chat, Coordinator Chat, Results Analyzer | 1 | Widoczny bez interakcji |
| Pole wprowadzania poleceń | Pole tekstowe do wydania polecenia wykonawcy | Ręczne polecenie użytkownika wobec danej roli, poza kolejką | Pole tekstowe (`.dn-textarea`) | Domyślne (puste) · wypełnione · wysyłanie (ładowanie) | Wysłanie dodaje polecenie do historii i przekazuje je roli | Executor Chat, Coordinator Chat | 1 | Widoczne bez interakcji |
| Przycisk wyślij | Ikona wysłania treści pola poleceń | Przekazanie wpisanego polecenia | Mała ikonka (`wyslij`) | Domyślny · ładowanie | Przy wypełnionym polu wysyła treść, czyści pole i dopisuje wpis do historii; kliknięcie przy pustym polu pokazuje komunikat kontekstowy „Wpisz polecenie przed wysłaniem” zamiast blokować przycisk | Przy polu wprowadzania poleceń | 1 | Widoczny bez interakcji |
| Przycisk załącznika | Ikona spinacza | Dołączenie pliku jako kontekstu polecenia | Mała ikonka (`spinacz`) | Domyślny · załączono (plakietka z liczbą plików) | Otwiera okno wyboru pliku; podgląd załącznika w historii | Przy polu wprowadzania poleceń | 1 | Widoczny bez interakcji |
| Wskaźnik pobrania zadania (dequeue) | Etykieta w historii sygnalizująca automatyczne pobranie zadania z kolejki | Odróżnienie zadań zainicjowanych przez kolejkę od poleceń ręcznych | Mała etykieta / plakietka | Widoczna przy zadaniu pobranym z kolejki | Kliknięcie otwiera szczegóły zadania w sekcji Kolejki | Obszar historii Executor Chat | 1 | Widoczny bez interakcji |

### 10.5. Elementy Subagent Network

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Panel Subagent Network | Rozwijany panel boczny w oknie Executor Chat | Podgląd i sterowanie do 15 podagentami uruchomionymi przez wykonawcę | Panel boczny średniej wielkości | Zwinięty (0 podagentów) · rozwinięty · pełny (15/15) | Kliknięcie nagłówka rozwija lub zwija panel | Executor Chat (Executor 1 i Executor 2, niezależnie) | 3 | Rozwinięcie `Subagent Network ▼` jako kolumna boczna |
| Przycisk „Uruchom Subagent Network” | Kontrolka wywołania akcji `split` | Podział złożonego zadania na zakresy przydzielane podagentom | Przycisk (`.dn-btn--zloty`) | Domyślny · ładowanie | Otwiera kreator podziału zadania na zakresy i wcielenia podagentów; przy zadaniu ocenionym jako zbyt proste lub przy osiągniętym referencyjnym limicie 15 podagentów wyświetla dodatkowo komunikat kontekstowy sygnalizujący sytuację, a akcja pozostaje dostępna | Panel Subagent Network | 2 | Akcja w oknie wykonawcy |
| Karta podagenta | Pojedyncza pozycja reprezentująca jednego z do 15 podagentów | Podgląd wcielenia, postępu i statusu podagenta | Mała karta (`.dn-karta--pozycja`) | Oczekuje · pracuje · zakończony (sukces) · błąd | Kliknięcie otwiera szczegółowy log pracy podagenta | Panel Subagent Network | 3 | Widoczna po rozwinięciu panelu Subagent Network |
| Pasek postępu podagenta | Wskaźnik procentowego zaawansowania zakresu | Orientacja, jak blisko ukończenia jest dany podagent | Pasek postępu (mały) | 0–100% · nieokreślony (ładowanie) | Aktualizuje się na żywo kanałem WebSocket | Karta podagenta | 3 | Widoczny po rozwinięciu panelu Subagent Network |
| Przycisk „Scal wyniki” (merge) | Kontrolka wywołania akcji `merge` | Agregacja wyników wszystkich zakończonych podagentów w jeden rezultat | Przycisk (`.dn-btn--zloty`) | Domyślny · ładowanie | Scala wyniki i przekazuje rezultat do dalszego przepływu kolejki; kliknięcie przed zakończeniem wszystkich podagentów wyświetla komunikat kontekstowy o niekompletnym zestawie wyników, pozostawiając decyzję o scaleniu Operatorowi | Panel Subagent Network | 2 | Akcja w panelu Subagent Network |
| Przycisk „+ dodaj podagenta” | Ikona rozszerzenia sieci podagentów | Ręczne dodanie kolejnego podagenta poza automatycznym podziałem | Mała ikonka (`plus`) | Domyślny | Otwiera formularz wyboru wcielenia nowego podagenta; po osiągnięciu referencyjnego limitu 15 podagentów wyświetla dodatkowo komunikat kontekstowy sygnalizujący przekroczenie, bez blokowania dodania kolejnego | Panel Subagent Network | 3 | Rozwinięcie panelu Subagent Network |

### 10.6. Elementy okna Coordinator Chat — swoiste

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Panel planu etapów | Lista etapów procesu z przypisaną rolą i statusem | Wizualizacja planu zbudowanego przez Coordinatora | Lista/tabela średniej wielkości | Pusty (nowy proces) · z etapami w różnych stanach | Kliknięcie etapu otwiera jego szczegóły i powiązane zadania kolejki | Coordinator Chat | 1 | Widoczny bez interakcji |
| Kreator promptu (Prompt Builder wewnętrzny) | Panel formularza budowy promptu dla wskazanego wykonawcy | Dynamiczne tworzenie i optymalizacja poleceń przekazywanych rolom wykonawczym | Panel formularza | Pusty · w edycji · wysłany do kolejki | Wysłanie dodaje prompt jako zadanie (`enqueue`) do kolejki wskazanej roli | Coordinator Chat | 2 | Rozwinięcie `Kreator promptu ▼` |
| Widok stanu kolejki (miniaturowy) | Skrócone zestawienie liczby zadań w kolejkach procesu | Szybki podgląd obciążenia kolejki bez przechodzenia do sekcji Kolejki | Mała tabela | Aktualny · odświeżany na żywo | Kliknięcie przenosi do pełnego widoku sekcji Kolejki | Coordinator Chat | 1 | Widoczny bez interakcji |
| Zestaw przycisków sterowania procesem | Siedem przycisków: Start, Stop, Pauza, Wznów, Przekaż, Powtórz, Waliduj | Bezpośrednie sterowanie przebiegiem procesu z poziomu okna Coordinatora | Zestaw przycisków ikonowych (`.dn-btn-ikona`) | Każdy przycisk niezależnie: dostępny · ładowanie | Wywołuje odpowiadającą akcję lub mechanizm (rozdz. 3.3); przycisk nieadekwatny do bieżącego stanu procesu (np. „Wznów” przy procesie już uruchomionym) pozostaje klikalny i pokazuje komunikat kontekstowy zamiast się blokować | Coordinator Chat, element zbiorczy `Sterowanie ▼` | 3 | Rozwinięcie `Sterowanie ▼` |
| Selektor trybu współpracy | Pole wyboru jednego z czterech trybów (rozdz. 3.4) | Ustalenie relacji między Executorem 1 a Executorem 2 dla bieżącego procesu | Pole wyboru (`.dn-select`) | Domyślny (praca niezależna) · zmieniony | Zmiana trybu natychmiast wpływa na sposób kierowania zadań między wykonawcami | Coordinator Chat, karta roli Executor 2 | 2 | Znacznik kontekstowy `Tryb ▼` |

### 10.7. Elementy okna Results Analyzer — swoiste

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Panel porównania wyników | Dwukolumnowy widok zestawiający wynik Executora 1 i Executora 2 | Ułatwienie oceny lub rozstrzygnięcia rozbieżności między wykonawcami | Duży panel dwukolumnowy | Ukryty (jeden wynik) · widoczny (dwa wyniki do porównania) | Przewijanie i podświetlanie różnic (analogicznie do Diff/Grep Panel) | Results Analyzer | 3 | Rozwinięcie `Porównanie wyników ▼` jako kolumna boczna |
| Selektor wcielenia | Pole wyboru jednego z siedmiu wcieleń roli (rozdz. 3.5) | Ustalenie, w jakiej funkcji kontrolnej działa Executor 3 / Validator w bieżącym etapie | Pole wyboru (`.dn-select`) | Domyślny (Validator) · zmieniony | Zmiana wcielenia dostosowuje zestaw kryteriów oceny prezentowanych w panelu | Nagłówek Results Analyzer | 2 | Znacznik kontekstowy `Wcielenie ▼` |
| Przycisk „Zatwierdź” | Kontrolka pozytywnej oceny etapu | Uznanie wyniku za zgodny i przekazanie procesu dalej | Przycisk (`.dn-btn--zloty`), ikona `ptaszek` | Domyślny · ładowanie · zatwierdzony | Kieruje wynik do kolejnego etapu orkiestracji (`route`) | Results Analyzer | 1 | Widoczny bez interakcji |
| Przycisk „Odrzuć / zgłoś niezgodność” | Kontrolka negatywnej oceny etapu | Zwrot etapu do Coordinatora z uwagami, w celu powtórzenia | Przycisk (`.dn-btn--blad`) | Domyślny · ładowanie | Wywołuje `retry` z dołączonym uzasadnieniem | Results Analyzer | 1 | Widoczny bez interakcji |
| Przycisk „Rozstrzygnij” (Arbitrator) | Kontrolka dostępna we wcieleniu Arbitrator | Rozstrzygnięcie, które podejście spośród wyników Executora 1 i Executora 2 przyjąć | Przycisk (`.dn-btn--zarys`), ikona `waga` | Widoczny wyłącznie przy wykrytej rozbieżności · ładowanie | Wywołuje `branch`, kierując dalszy przepływ na wybraną ścieżkę | Results Analyzer, przy aktywnym panelu porównania | 3 | Rozwinięcie `⚖ ▼` przy wykrytej rozbieżności |
| Pole uzasadnienia oceny | Pole tekstowe towarzyszące decyzji | Zapisanie uwag przekazywanych do Coordinatora wraz z decyzją | Pole tekstowe (`.dn-textarea`) | Domyślne (puste) · wypełnione | Treść dołączana jest do zdarzenia `retry` lub do historii etapu | Results Analyzer | 1 | Widoczne bez interakcji |

### 10.8. Elementy silnika kolejek i sekcji Orkiestracja

Pełny katalog przedstawiają rozdziały 4.5 i 5.4.

### 10.9. Elementy sekcji Harmonogram i automatyki oraz Monitor procesu

Katalog elementów integracji z modułem Automations przedstawia rozdział 7.6. Poniżej elementy swoiste sekcji Monitor procesu.

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Tabela przebiegów | Historia kolejnych uruchomień procesu | Przegląd wszystkich dotychczasowych cykli pracy zespołu | Duża tabela (`.dn-tabela`) | Domyślna · pusty stan (pierwszy przebieg jeszcze się nie zakończył) | Kliknięcie wiersza otwiera szczegóły przebiegu (log, czas trwania, wynik) | Sekcja Monitor procesu | 1 | Widoczna bez interakcji |
| Plakietka statusu przebiegu | Znacznik sukcesu, błędu lub trwania przebiegu | Szybka orientacja w wyniku każdego cyklu | Plakietka + kropka (`.dn-plakietka--stan`) | Sukces · w toku · błąd · wstrzymany | Kliknięcie przy błędzie otwiera bezpośrednio przyczynę niepowodzenia | Wiersz tabeli przebiegów | 1 | Widoczna bez interakcji |
| Widok hierarchii decyzji | Wizualizacja sześciu poziomów hierarchii (rozdz. 8) z zaznaczeniem, kto obecnie podejmuje decyzję | Zrozumienie, na jakim poziomie znajduje się bieżący etap procesu | Schemat/diagram interaktywny | Domyślny · z podświetlonym bieżącym poziomem | Kliknięcie poziomu pokazuje uprawnienia i zakres decyzji tego poziomu | Sekcja Monitor procesu | 2 | Zakładka sekcji Monitor procesu |
| Przełącznik trybu AOD w monitorze | Kontrolka tożsama z opisaną w rozdziale 2.6 | Zmiana trybu Always On Display bezpośrednio z poziomu monitora | Przełącznik segmentowy | Obserwator · operator | Jak w rozdziale 2.6 | Nagłówek sekcji Monitor procesu | 2 | Znacznik `AOD ▼` w nagłówku sekcji |
| Przycisk „Szczegóły przebiegu” | Odnośnik/przycisk otwierający pełny log cyklu | Głęboki wgląd w każdy krok danego przebiegu — pełna, nieograniczona widoczność procesu | Przycisk tekstowy / link | Domyślny · najechanie | Otwiera pełny, chronologiczny log zdarzeń przebiegu, łącznie z każdą akcją silnika kolejek | Wiersz tabeli przebiegów | 3 | Menu kebab (⋮) wiersza tabeli przebiegów |

### 10.10. Elementy Chat Window i Execution Loop Window

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Okno Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca | Wydawanie poleceń w języku naturalnym, odbiór strumienia wyników, zatwierdzanie i przerywanie działań, sterowanie wszystkimi procesami środowiska | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Domyślny · strumień aktywny · oczekiwanie na zatwierdzenie · błąd kanału | Wysłanie polecenia dopisuje wpis do strumienia i uruchamia odpowiadającą operację procesu | Każda przestrzeń robocza środowiska, każda karta sesji | 1 | Widoczne bez interakcji |
| Pole polecenia Chat Window | Pole tekstowe z ikoną załącznika i ikoną wysłania | Wprowadzenie polecenia dla Wykonawcy | Pole tekstowe (`.dn-textarea`) pełnej szerokości kolumny | Domyślne · wypełnione · wysyłanie | Wysyła treść i czyści pole; kliknięcie przy pustym polu pokazuje komunikat kontekstowy | Chat Window | 1 | Widoczne bez interakcji |
| Akcje wyniku w strumieniu | Zestaw akcji przy wpisie wyniku: pokazanie wyniku, wyjaśnienie, zatwierdzenie, przerwanie | Praca na wyniku bez opuszczania okna komunikacji | Zestaw zwinięty (`Operacje ▼`) | Domyślny · rozwinięty · ładowanie | Rozwinięcie prezentuje pełną listę akcji dla danego wpisu | Strumień Chat Window | 3 | Rozwinięcie `Operacje ▼` przy wpisie wyniku |
| Wskaźnik stanu wykonania w Chat Window | Zestawienie liczby aktywnych wątków, stanu bieżącego etapu i sygnalizacji oczekiwania na zatwierdzenie | Orientacja w stanie procesu bez opuszczania okna komunikacji | Zestaw plakietek (`.dn-plakietka--stan`) | Bezczynny · wątki aktywne · oczekiwanie na zatwierdzenie · błąd | Kliknięcie plakietki otwiera Execution Loop Window z kontekstem wskazanego wątku | Nagłówek Chat Window | 1 | Widoczny bez interakcji |
| Wybór wykonawcy i modelu w Chat Window | Selektor przypisania agenta lub modelu bazowego do rozmowy | Ustalenie, kto odpowiada na polecenia wydawane w oknie komunikacji | Znacznik kontekstowy z listą rozwijaną | Domyślny · rozwinięty · wybrany | Po wyborze element zwija się samoczynnie, a znacznik przyjmuje nazwę wykonawcy | Pasek kontekstu Chat Window | 2 | Znacznik `Wykonawca ▼` |
| Wybór trybu pracy i poziomu wysiłku | Selektory parametrów przebiegu polecenia | Dostrojenie sposobu i głębokości realizacji polecenia | Znaczniki kontekstowe z listami rozwijanymi | Domyślny · rozwinięty · wybrany | Zmiana obowiązuje od kolejnego polecenia; element zwija się samoczynnie | Pasek kontekstu Chat Window | 2 | Znaczniki `Tryb ▼`, `Wysiłek ▼` |
| Operacje diagnostyczne rozmowy | Podgląd surowego kontekstu, śledzenie wywołań, tryb administracyjny | Diagnoza przebiegu rozmowy i kontekstu przekazanego Wykonawcy | Panel diagnostyczny | Domyślny · otwarty · ładowanie | Otwarcie prezentuje surowy kontekst i ślad wywołań bieżącej rozmowy | Chat Window | 4 | Polecenie języka naturalnego, skrót klawiszowy, wyszukiwarka funkcji |
| Wybór wątku i filtr ról | Selektory zawężenia widoku pętli do wskazanego wątku, roli lub kolejki | Utrzymanie czytelności przy wielu równoległych torach wykonania | Znaczniki kontekstowe z listami rozwijanymi | Domyślny (wszystkie wątki) · zawężony | Zawęża strumień komunikatów i listę zadań do wybranego zakresu | Nagłówek Execution Loop Window | 2 | Znaczniki `Wątek ▼`, `Rola ▼` |
| Zestaw operacji na zadaniu pętli | Akcje `route`, `branch`, `condition`, `delay`, `split`, `merge`, zmiana priorytetu, przypisanie do innej roli | Sterowanie pojedynczym zadaniem bez przechodzenia do sekcji Kolejki | Zestaw zwinięty (`Operacje ▼`) w wierszu zadania | Domyślny · rozwinięty · ładowanie | Wywołuje odpowiednią akcję `queue.action`; przy niepasującym stanie zadania pokazuje komunikat kontekstowy | Panel dekompozycji zlecenia, Execution Loop Window | 3 | Menu kebab (⋮) przy zadaniu |
| Edycja reguł orkiestracji i dziennik komunikatów pętli | Definicja zależności oraz dziennik surowych komunikatów sterujących pętli | Diagnostyka wykonania i zmiana reguł przepływu bez opuszczania okna pętli | Panel edycji reguł i log chronologiczny | Domyślny · w edycji · zapisany · ostrzeżenie walidacji reguły | Zapis reguły dodaje ją do listy zależności orkiestracji (rozdz. 5); dziennik prezentuje pełny, nieograniczony zapis komunikatów | Execution Loop Window | 4 | Polecenie języka naturalnego, skrót klawiszowy, tryb administracyjny |
| Okno Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór, orkiestracja i kontrola realizacji procesów przy pracy wielowątkowej | Kolumna sąsiadująca z Chat Window | Domyślny · pętla uruchomiona · pętla wstrzymana · oczekiwanie na kontrolę jakości · błąd przebiegu | Sterowanie przebiegiem działa na całym procesie, niezależnie od aktywnego okna roli | Każda karta sesji | 1 | Widoczne bez interakcji |
| Panel dekompozycji zlecenia | Lista zadań, na które Koordynator rozłożył bieżące zlecenie, z przypisaną rolą i stanem | Zrozumienie, z czego składa się bieżące zlecenie i kto je realizuje | Lista pozycji (`.dn-karta--pozycja`) | Domyślna · pusty stan · z zadaniem zaznaczonym | Kliknięcie zadania otwiera jego szczegóły w sekcji Kolejki | Execution Loop Window | 1 | Widoczny bez interakcji |
| Strumień komunikatów sterujących | Zapis wymiany między Koordynatorem a Wykonawcami: prompty etapów, wyniki cząstkowe, skierowania, ponowienia | Nadzór nad przebiegiem koordynacji zadań | Obszar zawartości kolumny | Domyślny · strumień aktywny · przewijalny | Nowe komunikaty dopisywane na żywo kanałem WebSocket | Execution Loop Window | 1 | Widoczny bez interakcji |
| Wskaźniki przebiegu pętli | Numer cyklu, liczba zadań w toku, stan kontroli jakości | Bieżąca orientacja w stanie pętli wykonawczej | Zestaw plakietek (`.dn-plakietka--stan`) | Uruchomiona · wstrzymana · oczekuje na decyzję · błąd | Kliknięcie plakietki otwiera szczegóły odpowiadającego elementu pętli | Execution Loop Window, panel stanu procesu | 1 | Widoczne bez interakcji |
| Sterowanie przebiegiem pętli | Zestaw operacji: wstrzymanie, wznowienie, przerwanie, korekta zlecenia, ponowienie zadania | Bezpośrednia kontrola realizacji procesu | Zestaw przycisków ikonowych (`.dn-btn-ikona`) oraz element zbiorczy | Domyślny · ładowanie · wykonano | Każda operacja pozostaje klikalna; przy niepasującym stanie pokazuje komunikat kontekstowy | Execution Loop Window | 1 | Widoczne bez interakcji |
| Wynik kontroli jakości w pętli | Blok prezentujący ocenę Executora 3 / Validatora i decyzję o ponowieniu | Domknięcie pętli decyzją o kontynuacji albo powtórzeniu | Blok treści z plakietką stanu | Brak oceny · ocena pozytywna · ocena negatywna | Kliknięcie otwiera Results Analyzer z kontekstem ocenianego zadania | Execution Loop Window | 1 | Widoczny bez interakcji |

---

---

## 11. Diagramy przepływu — orkiestracja i kolejkowanie

Rozdział zbiera i domyka w jednym miejscu wszystkie diagramy przepływu środowiska — pełną pętlę orkiestracji, cykl życia zadania w kolejce, eskalację przy niepowodzeniu oraz pracę ciągłą — jako punkt odniesienia niezależny od poszczególnych rozdziałów mechanizmu.

### 11.1. Pełny przepływ orkiestracji, z adnotacją okien i sekcji panelu

```
                    Użytkownik (Operator)                 ── decyzja nadrzędna, dowolny moment
                             │
                    Chat Window (Użytkownik ↔ Wykonawca)   ── lewa kolumna, stała
                    polecenie · zatwierdzenie · przerwanie         │
                             │                                     │
                    Always On Display                      ── Monitor procesu (rozdz. 8, sekcja panelu)
                    (obserwator / operator, rozdz. 2)              │
                             │                                     │ interwencja (Mobile, rozdz. 7.5)
                    Execution Loop Window                          │
                    (Koordynator ↔ Wykonawca, rozdz. 1.7, 9.5)     │
                    dekompozycja · przydział · nadzór ·            │
                    kontrola realizacji · sterowanie przebiegiem   │
                             │                                     │
                        Coordinator                                │
              (plan · podział pracy · budowa promptów              │
               · sterowanie · zarządzanie kolejką)                 │
                    okno: Coordinator Chat ────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
    Executor 1          Executor 2        Executor 3 / Validator
   (wykonanie)      (wykonanie równoległe)  (kontrola: Validator,
   okno: Executor         okno: Executor      Reviewer, Security Auditor,
   Chat                   Chat               Architect, QA Lead, Arbitrator …)
         │                   │               okno: Results Analyzer
         ▼                   ▼
  Subagent Network     Subagent Network
  (do 15 podagentów)   (do 15 podagentów)
  panel w Executor      panel w Executor
  Chat, sekcja Role     Chat, sekcja Role
         │                   │
         └────────►  Silnik kolejek  ◄────────┘
          sekcja Kolejki panelu orkiestracji (rozdz. 4)
          (enqueue · dequeue · delay · retry · pause · resume
           · split · merge · route · branch · condition)
                             │
                       Orkiestracja
        sekcja Orkiestracja panelu orkiestracji (rozdz. 5)
        (respektowanie zależności między rolami, kolejkami, zadaniami)
                             │
                  Integracja z Automations
          sekcja Harmonogram i automatyki (rozdz. 7)
          (harmonogram · cykliczność → praca ciągła 24/7/365)
                             │
                             └──────────► kolejny przebieg (enqueue) ──┐
                                                                        │
                                          ◄─────────────────────────────┘
```

Pętla wykonawcza całego przepływu prowadzona jest w Execution Loop Window: okno prezentuje dekompozycję zlecenia na zadania, przydział zadań rolom, stan kolejek, komunikaty sterujące, wyniki kontroli jakości oraz wskaźniki przebiegu, a także udostępnia sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Polecenia użytkownika wchodzą do przepływu przez Chat Window, obecne w każdej przestrzeni roboczej środowiska.

### 11.2. Przepływ silnika kolejek — cykl życia pojedynczego zadania

Pełny diagram i tabela akcji: rozdział 4.2–4.3. Poniżej wersja skrócona z zaznaczeniem punktów decyzyjnych.

```
        enqueue                dequeue           split/merge         route/branch/condition
  Coordinator ──► KOLEJKA ──► Wykonawca ──► Subagent Network ──► ORKIESTRACJA ──► kolejny krok
                     ▲                                                   │
                     │                                                   ▼
                  pause/resume                                  Executor 3 / Validator
                  (Coordinator, AOD)                                     │
                     │                                          zatwierdzenie / retry
                     └──────────────────────────────────────────────────┘
```

### 11.3. Przepływ eskalacji przy niepowodzeniu walidacji

```
Executor 1 ── route ──► Executor 3 / Validator ── ocena negatywna ──► Coordinator
                                                                          │
                                                                     retry (enqueue
                                                                     z nowym promptem)
                                                                          │
                                                                          ▼
                                                                    Executor 1
                                                                          │
                                                              (cykl powtarza się do
                                                               oceny pozytywnej lub
                                                               interwencji użytkownika
                                                               — poziom 1 hierarchii,
                                                               rozdz. 8.1 — w dowolnym
                                                               momencie)
```

### 11.4. Przepływ pracy ciągłej 24/7/365 z modułem Automations

```
 Scheduler (Automations) ── cyklicznie uruchamia ──► Coordinator ──► pełna pętla (11.1)
         ▲                                                                  │
         │                                                                  ▼
         │                                                    Execution Monitor (Automations)
         │                                                    + Monitor procesu (MultitaskingAI)
         │                                                                  │
         │                                          Always On Display (operator) ── wykrywa
         │                                          powtarzające się niepowodzenie
         │                                                                  │
         └──────────────── interwencja przez Mobile (pause/resume) ◄────────┘
```

### 11.5. Przepływ interwencji zdalnej (Mobile + Always On Display)

```
Proces działa w tle (praca ciągła)
        │
Always On Display wykrywa zdarzenie wymagające decyzji
        │
powiadomienie push ──► urządzenie mobilne Operatora
        │
Operator otwiera widok Mobile procesu MultitaskingAI
        │
   ┌────┴────┬─────────────┬──────────────┐
   ▼         ▼             ▼              ▼
zatwierdź  wstrzymaj    zmodyfikuj     przejmij bezpośrednie
 krok      (pause)      konfigurację    sterowanie (poziom 1
                         Coordinatora    hierarchii, rozdz. 8.1)
        │
proces wznawia się (resume) z zastosowaną decyzją
```

---

## 12. Stany

Rozdział zbiera stany domenowe (procesu, zadania, roli, Always On Display) oraz konwencję stanów interakcji elementów interfejsu, obowiązującą jednolicie w całym środowisku.

### 12.1. Stany przebiegu procesu (pętli orkiestracji)

| Stan | Znaczenie | Sygnalizacja wizualna | Możliwe przejścia |
|---|---|---|---|
| Nieuruchomiony | Zespół skonfigurowany, proces jeszcze nie wystartował | Plakietka neutralna „nieuruchomiony” | → uruchomiony (Start) |
| Uruchomiony | Proces aktywnie wykonuje etapy planu Coordinatora | Plakietka sukcesu, kropka pulsująca | → wstrzymany, → zakończony, → błąd |
| Wstrzymany | Proces zatrzymany akcją `pause`, stan zachowany | Plakietka ostrzeżenia „wstrzymany” | → uruchomiony (`resume`), → zatrzymany trwale (Stop) |
| Oczekujący na decyzję | Etap wymaga zatwierdzenia (Executor 3 / Validator, AOD operator lub użytkownik) | Plakietka info, ikona `oko` | → uruchomiony, → wstrzymany, → powtórzony (`retry`) |
| Zakończony sukcesem | Wszystkie etapy planu wykonane, wynik zaakceptowany | Plakietka sukcesu, ikona `ptaszek` | → nieuruchomiony (kolejny przebieg cyklu — rozdz. 7) |
| Zakończony błędem | Etap nieodwracalnie nieudany bez dalszego `retry` | Plakietka błędu | → nieuruchomiony (po interwencji użytkownika) |
| Zatrzymany trwale | Proces zakończony poleceniem Stop, bez automatycznego wznowienia | Plakietka neutralna „zatrzymany” | → nieuruchomiony (ponowna konfiguracja) |

### 12.2. Stany roli / okna roboczego

| Stan | Znaczenie | Sygnalizacja wizualna | Dotyczy |
|---|---|---|---|
| Bezczynna | Rola przypisana, brak aktywnego zadania | Kropka neutralna | Executor 1, Executor 2, Coordinator, Executor 3 / Validator |
| Pracuje | Rola aktywnie przetwarza zadanie | Kropka sukcesu, animacja | Wszystkie role |
| Oczekuje na zależność | Zadanie roli zablokowane regułą orkiestracji do czasu zakończenia etapu poprzedzającego | Kropka ostrzeżenia | Executor 1, Executor 2 |
| Oczekuje na ocenę | Wynik przekazany do Executora 3 / Validatora, oczekuje na werdykt | Kropka info | Executor 1, Executor 2 |
| Błąd wykonania | Zadanie zakończone niepowodzeniem technicznym | Kropka błędu | Wszystkie role |
| Brak przypisania | Karta roli utworzona, lecz bez wskazanego agenta lub modelu | Karta w stanie pustym (`.dn-pusty-stan`) | Wszystkie role |

### 12.3. Stany zadania w kolejce

| Stan | Znaczenie | Akcja wprowadzająca | Akcja wyprowadzająca |
|---|---|---|---|
| W kolejce | Zadanie oczekuje na pobranie | `enqueue` | `dequeue` |
| Odroczone | Zadanie czeka na upływ czasu lub spełnienie warunku | `delay` | automatyczne wznowienie do stanu „w kolejce” |
| W realizacji | Zadanie pobrane i wykonywane przez wykonawcę | `dequeue` | zakończenie pracy wykonawcy |
| Podzielone | Zadanie rozłożone na podzadania Subagent Network | `split` | `merge` |
| Scalone | Wyniki podzadań połączone w jeden rezultat | `merge` | przekazanie do orkiestracji |
| Skierowane | Zadanie lub wynik przekazane do innej roli lub kolejki | `route` | pojawienie się w kolejce lub oknie roli docelowej |
| Rozgałęzione | Przepływ podzielony na alternatywne ścieżki | `branch` | wybór ścieżki przez warunek lub decyzję Arbitratora |
| Wstrzymane | Kolejka, w której zadanie się znajduje, jest zatrzymana | `pause` (na poziomie kolejki) | `resume` |
| Do powtórzenia | Zadanie zwrócone po negatywnej ocenie lub błędzie | `retry` | powrót do stanu „w kolejce” (`enqueue`) |
| Zakończone | Zadanie w pełni obsłużone, wynik zaakceptowany | zatwierdzenie Executora 3 / Validatora lub Coordinatora | — |

### 12.4. Diagram przejść stanów zadania w kolejce

```
        enqueue                              delay
   ┌──────────────► [ W KOLEJCE ] ◄────────────────────┐
   │                      │                             │
   │                   dequeue                          │ (upływ czasu / warunek)
   │                      ▼                              │
   │              [ W REALIZACJI ] ─── split ──► [ PODZIELONE ] ─┐
   │                      │                                       │ merge
   │                pauza kolejki                                 ▼
   │                      ▼                              [ SCALONE ]
   │                [ WSTRZYMANE ] ── resume ──► z powrotem do „W REALIZACJI”
   │                      
   │              wynik gotowy
   │                      ▼
   │                 [ SKIEROWANE ] ── route ──► rola / kolejka docelowa
   │                      │
   │                 branch (wg reguły orkiestracji)
   │                      ▼
   │              [ ROZGAŁĘZIONE ] ── warunek (condition) ──► ścieżka A / ścieżka B
   │                      ▼
   │              ocena Executora 3 / Validatora
   │              ┌───────┴────────┐
   │         pozytywna          negatywna
   │              ▼                 ▼
   │        [ ZAKOŃCZONE ]    [ DO POWTÓRZENIA ]
   │                                 │
   └─────────────── retry ──────────┘
```

### 12.5. Stany Always On Display w kontekście procesu

| Tryb | Stan szczegółowy | Zachowanie |
|---|---|---|
| Obserwator | Bezczynny | Brak aktywnej sugestii, gotowość do podglądu |
| Obserwator | Doradza | Wyświetla dymek z sugestią lub ostrzeżeniem, bez możliwości interwencji |
| Operator | Bezczynny | Gotowość do interwencji, brak bieżącego zdarzenia |
| Operator | Doradza | Jak w trybie obserwatora, rozszerzone o przyciski akcji w dymku |
| Operator | Interweniuje | Aktywnie wykonuje zatwierdzenie, wstrzymanie lub modyfikację procesu, w tym przez Mobile |

Powyższe stany opisują funkcję wyłącznie w odniesieniu do procesu środowiska. Stany ogólnoplatformowe funkcji — spoczynek, sugestia oczekująca, rozmowa, nasłuch, mowa, wyciszenie, tryby obecności — podaje `funkcje-globalne/always-on-display.md`, rozdz. 9.1–9.2.

### 12.6. Konwencja stanów elementów interfejsu

Stany z poniższej tabeli obowiązują jednolicie dla wszystkich elementów interaktywnych katalogowanych w rozdziale 10, zgodnie z systemem wizualnym platformy.

| Stan | Reguła wizualna | Zastosowanie |
|---|---|---|
| Domyślny | Wygląd spoczynkowy elementu | Każdy element interaktywny |
| Najechanie (hover) | Subtelna zmiana tła lub uniesienie dla elementów akcentowanych złotem | Przyciski, karty, pozycje list |
| Aktywny / wciśnięty | Przesunięcie `translateY(1px)` lub zaznaczenie trwałe (zakładki, karty ról) | Przyciski, zakładki, pozycje nawigacyjne |
| Niegotowy (aktywny, z komunikatem) | Wygląd tożsamy ze stanem domyślnym, bez przyciemnienia i bez blokady kursora — element pozostaje w pełni klikalny | Akcje, których warunek nie jest jeszcze spełniony w bieżącym stanie procesu (np. `merge` przed zakończeniem wszystkich podagentów); kliknięcie wywołuje komunikat kontekstowy zamiast wykonania akcji lub wykonuje ją częściowo, z ostrzeżeniem |
| Ładowanie | Spinner (`.dn-spinner`) lub pasek postępu w miejscu treści | Akcje silnika kolejek w toku, zapis profilu, wczytanie zespołu |
| Ostrzeżenie | Obramowanie koloru ostrzegawczego, komunikat kontekstowy — sygnalizuje problem, nie blokuje zapisu ani wykonania | Pole z wartością nietypową, duplikat nazwy, reguła wewnętrznie sprzeczna |
| Błąd | Obramowanie lub tło koloru błędu, komunikat kontekstowy | Niepowodzenie techniczne poza kontrolą formularza — nieudane połączenie z Automations, brak dostępu do mikrofonu, błąd pobrania danych |
| Fokus | Pierścień złoty 2 px, wyzwalany przez nawigację klawiaturą | Wszystkie kontrolki interaktywne, obowiązkowo |

---

## 13. Konfigurowalność i punkty izolacji na poziomie roli

### 13.1. Poziom zasięgu „Rola” — pierwszeństwo najwyższe

Okno konfiguracji punktów izolacji udostępnia siedem poziomów zasięgu reguł izolacji, z których poziom „Rola (MultitaskingAI)” ma pierwszeństwo najwyższe — reguła ustalona na tym poziomie nadpisuje reguły odziedziczone ze wszystkich poziomów szerszych (globalny, środowisko, moduł, para modułów, projekt, karta sesji). Umożliwia to na przykład przypisanie Executorowi 1 pracującemu z zewnętrznymi źródłami odmiennej polityki dostępu sieciowego niż Executorowi 2 pracującemu wyłącznie na materiałach lokalnych — bez zmiany ustawień globalnych ani ustawień pozostałych ról.

```
 Pierwszeństwo poziomów zasięgu (od najniższego do najwyższego)

  globalny  →  środowisko  →  moduł  →  para modułów  →
  →  projekt  →  karta sesji  →  rola
  (warstwa bazowa)                        (pierwszeństwo najwyższe)

  brak ustawienia na poziomie  =  dziedziczenie z poziomu szerszego,
  nigdy blokada
```

### 13.2. Struktura profilu izolacji przypisanego roli

Profil izolacji przypisany roli grupuje wartości izolacji kontekstu (historia, pamięć, kontekst — współdzielone albo odrębne) oraz wartości ośmiu zakresów izolacji technicznej (katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania — każdy włączony albo wyłączony). Profil raz zbudowany może zostać przypisany wielokrotnie, do wielu ról lub sesji, bez ponownej konfiguracji.

| # | Zakres | Stan włączony | Stan wyłączony (domyślny) |
|---|---|---|---|
| 1 | Katalog roboczy sesji | Rola pracuje we własnym katalogu, niewidocznym dla innych ról | Rola korzysta ze wspólnego katalogu roboczego |
| 2 | Środowisko procesu | Rola otrzymuje własny zestaw zmiennych środowiskowych | Proces dziedziczy wspólne środowisko uruchomieniowe |
| 3 | Katalog danych i konfiguracji modelu | Rola ma własny katalog danych kanału modelu | Katalog współdzielony z innymi rolami tego samego kanału |
| 4 | Dostęp sieciowy | Rola otrzymuje odrębny, ograniczony dostęp sieciowy | Rola korzysta ze wspólnego dostępu sieciowego |
| 5 | Zakres odczytu i zapisu plików | Dostęp ograniczony do jawnie dozwolonych ścieżek | Pełny dostęp do zasobów plikowych w ramach uprawnień użytkownika |
| 6 | Konto i token per sesja | Rola korzysta z własnego, dedykowanego tokenu | Role korzystają ze wspólnych danych dostępowych platformy |
| 7 | Model procesu | Rola uruchamia własną, niezależną instancję procesu modelu | Rola korzysta ze wspólnej puli procesów modelu |
| 8 | Serwer wykonania | Rola przypisana do dedykowanego serwera wykonania | Rola wykonywana na współdzielonym serwerze platformy |

### 13.3. Atrybuty roli w modelu danych

Encja „Rola w MultitaskingAI” definiowana jest atrybutami: identyfikator, przypisanie do sesji lub orkiestracji, rodzaj roli, przypisany model (agent lub model bazowy), konfiguracja współpracy. Profil izolacji jest przypisywalny do sesji, roli lub projektu — w warstwie komunikacji przypisanie realizuje polecenie `isolation.profile.assign`.

### 13.4. Stan wyjściowy — wykonanie jest domyślne

Dopóki użytkownik nie utworzy profilu izolacji ani nie zmieni ustawienia dla żadnej roli, obowiązuje stan wyjściowy: odrębny kontekst (historia, pamięć, kontekst) per karta sesji, przy jednoczesnym braku aktywnej izolacji technicznej — żaden z ośmiu zakresów nie jest domyślnie włączony, więc procesy ról nie są technicznie ograniczane. Wszelka dalsza izolacja na poziomie roli — zarówno poluzowanie podziału kontekstu przez współdzielenie między rolami tego samego procesu, jak i włączenie zakresów izolacji technicznej — jest świadomą decyzją użytkownika, a nie ustawieniem narzuconym przez środowisko. Rozszerzenia podłączone do agenta (wtyczki, umiejętności, konektory, serwery MCP) działają w pełnym zakresie zarejestrowanym w Permissions Center, dopóki użytkownik świadomie go nie zawęzi.

### 13.5. Katalog elementów interfejsu okna izolacji

Pełny katalog przedstawia rozdział 4.4 dokumentu źródłowego Koncepcja platformy oraz makieta w rozdziale 9.10 niniejszego dokumentu; poniżej elementy swoiste zastosowaniu na poziomie „Rola”.

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pozycja „Rola” w selektorze zasięgu | Wpis listy poziomów zasięgu | Wybranie poziomu „Rola” jako kontekstu edytowanej macierzy izolacji | Pozycja listy (`.dn-karta--pozycja`) | Domyślna · zaznaczona | Po zaznaczeniu ujawnia dodatkowy selektor konkretnej roli (Executor 1/2/Coordinator/Validator) | Panel selektora zasięgu okna konfiguracji izolacji | 4 | Okno konfiguracji punktów izolacji |
| Selektor konkretnej roli | Pole wyboru, której z czterech ról dotyczy edytowany profil | Doprecyzowanie, do jakiej roli bieżącego procesu przypiąć profil | Pole wyboru (`.dn-select`) | Widoczny wyłącznie przy zasięgu „Rola” | Zmiana przełącza podgląd aktualnej polityki efektywnej danej roli | Panel selektora zasięgu, przy pozycji „Rola” | 4 | Okno konfiguracji punktów izolacji, przy pozycji „Rola” |
| Ikona objaśnienia kontekstowego | Znak zapytania przy każdym przełączniku macierzy | Wyjaśnienie działania ustawienia i jego wpływu na aplikację | Mała ikonka (`info`) | Domyślna · rozwinięta (tooltip widoczny) | Najechanie lub kliknięcie pokazuje tekst objaśnienia | Każdy wiersz macierzy izolacji | 2 | Najechanie lub kliknięcie ikony `[?]` |

---

## 14. Szablony konfiguracji zespołu

Poniższe szablony instancjonują szablon roli w postaci gotowych do zastosowania konfiguracji zespołu — pokazują złożenie już ustalonych elementów, nie wprowadzają nowych ról, trybów ani akcji.

### 14.1. Zespół „Budowa aplikacji”

Oparty na scenariuszu realizacji produktu cyfrowego w module Apps.

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie etapów, przypisanie pracy backendowi i frontendowi | — | Nie | Domyślny (dziedziczony) | Tak |
| Executor 1 | Realizacja backendu | Praca niezależna | Tak — Agent Backend, Agent API, Agent Database | Domyślny (dziedziczony) | Tak |
| Executor 2 | Realizacja frontendu (zasoby z modułu Design) | Praca niezależna | Tak — Agent UI | Domyślny (dziedziczony) | Tak |
| Executor 3 / Validator | Wcielenie Security Auditor, następnie QA Lead w kolejnych przebiegach | — | Nie | Domyślny (dziedziczony) | Tak |
| Always On Display | Operator procesu | — | — | — | — |

### 14.2. Zespół „Badanie i redakcja”

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie etapów badania i redakcji raportu | — | Nie | Domyślny | Nie (proces jednorazowy) |
| Executor 1 | Zbieranie i analiza źródeł | Przekazywanie wyników → Executor 2 | Nie | Domyślny | Nie |
| Executor 2 | Redakcja raportu końcowego na podstawie ustaleń Executora 1 | Przekazywanie wyników | Nie | Domyślny | Nie |
| Executor 3 / Validator | Wcielenie Reviewer — ocena spójności i kompletności raportu | — | Nie | Domyślny | Nie |
| Always On Display | Obserwator | — | — | — | — |

### 14.3. Zespół „Pętla ciągła 24/7”

Konfiguracja w pełni spięta z modułem Automations, przeznaczona do pracy bez stałego nadzoru użytkownika.

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie i sterowanie każdym cyklicznym przebiegiem | — | Nie | Odrębny profil roli (dostęp sieciowy wyłączony jako świadomy wybór) | Tak — Scheduler, Orchestrator |
| Executor 1 | Wykonanie zadań pobieranych z kolejki cyklicznej | Praca iteracyjna z Executorem 2 | Tak, w miarę potrzeby zadania | Odrębny profil roli | Tak — Queue Manager |
| Executor 2 | Poprawa i uzupełnienie wyników Executora 1 w kolejnym cyklu | Praca iteracyjna | Tak, w miarę potrzeby zadania | Odrębny profil roli | Tak — Queue Manager |
| Executor 3 / Validator | Wcielenie Validator — ocena jakości przed zamknięciem każdego cyklu, sterowana ustawieniem konfiguracyjnym | — | Nie | Odrębny profil roli | Tak — Execution Monitor |
| Always On Display | Operator, z interwencją przez Mobile | — | — | — | — |

### 14.4. Formularz pusty konfiguracji zespołu

| Pole | Wartość |
|---|---|
| Nazwa zespołu | |
| Cel procesu | |
| Coordinator — przypisany agent/model | |
| Executor 1 — zadanie | |
| Executor 1 — przypisany agent/model | |
| Executor 1 — Subagent Network (wcielenia podagentów) | |
| Executor 2 — zadanie | |
| Executor 2 — przypisany agent/model | |
| Executor 2 — Subagent Network (wcielenia podagentów) | |
| Tryb współpracy Executor 1 / Executor 2 | niezależna \| przekazywanie \| naprzemienna \| iteracyjna |
| Executor 3 / Validator — wcielenie | |
| Executor 3 / Validator — przypisany agent/model | |
| Zasięg i nazwa kolejki (kolejek) | |
| Zależności orkiestracji | |
| Profil izolacji — poziom „Rola” (per rola, jeśli odrębny od domyślnego) | |
| Powiązanie z Automations (tak/nie) | |
| Harmonogram (jeśli powiązanie = tak) | |
| Tryb Always On Display | obserwator \| operator |

---

## 15. Scenariusze operacyjne

### 15.1. Eskalacja przy niepowodzeniu walidacji

Executor 1 kończy etap pracy i przekazuje wynik do Executora 3 / Validatora (`route`). Validator ocenia wynik negatywnie. Ocena trafia do Coordinatora, który wydaje polecenie `retry` — zadanie wraca do kolejki (`enqueue`) z zaktualizowanym promptem zbudowanym przez Coordinatora na podstawie uwag Validatora. Cykl powtarza się do uzyskania oceny pozytywnej lub do interwencji użytkownika, który w każdej chwili może przejąć bezpośrednie sterowanie (poziom 1 hierarchii, rozdz. 8.1).

### 15.2. Rozstrzygnięcie konfliktu w trybie iteracyjnym

Executor 1 i Executor 2 pracują w trybie iteracyjnym nad tym samym artefaktem i po kilku cyklach ich wyniki są sprzeczne. Coordinator kieruje oba wyniki do Executora 3 / Validatora we wcieleniu Arbitrator. Arbitrator rozstrzyga, które podejście przyjąć, a decyzja trafia z powrotem do Coordinatora, który wydaje polecenie `branch`, kierując dalszy przepływ pracy na wybraną ścieżkę.

### 15.3. Praca ciągła z interwencją zdalną

Zespół skonfigurowany według szablonu „Pętla ciągła 24/7” (rozdz. 14.3) działa w trybie 24/7/365, uruchamiany cyklicznie przez Scheduler modułu Automations. Always On Display, działający w trybie operatora, wykrywa powtarzające się niepowodzenie walidacji w trzecim kolejnym cyklu i powiadamia użytkownika przez funkcję Mobile. Użytkownik, znajdując się poza stanowiskiem roboczym, wstrzymuje kolejkę (`pause`) z poziomu urządzenia mobilnego, analizuje przyczynę i wznawia proces (`resume`) po wprowadzeniu poprawki w instrukcjach Coordinatora.

### 15.4. Zespół w pełni oprzyrządowany — od agentów do pętli ciągłej

Scenariusz ilustruje maksymalne wykorzystanie dostępnego oprzyrządowania: od budowy agentów po pełną, autonomiczną pętlę.

| Krok | Działanie | Miejsce |
|---|---|---|
| 1 | Operator buduje cztery agenty w module Agents — po jednym pod każdą rolę — dobierając model bazowy, kanał, skille, konektory, serwery MCP i pamięć dla każdego | Agent Builder, Skills Manager, Connectors Manager |
| 2 | Operator ustala zakres uprawnień każdego agenta w Permissions Center: dostęp do rozszerzeń, modułów i zasobów, izolację techniczną, zdolność uruchamiania Subagent Network | Permissions Center |
| 3 | Operator otwiera sekcję Zespoły panelu orkiestracji i tworzy nowy zespół, przypisując cztery agenty do czterech ról | Panel orkiestracji ▸ Zespoły, Role |
| 4 | Operator ustala tryb współpracy Executor 1 / Executor 2 oraz definiuje kolejkę lokalną z regułami `branch` i `condition` | Panel orkiestracji ▸ Role, Kolejki |
| 5 | Operator definiuje zależności sekwencyjne i warunkowe w sekcji Orkiestracja, wiążąc etapy backendu, frontendu i walidacji | Panel orkiestracji ▸ Orkiestracja |
| 6 | Executor 1 i Executor 2, w toku pracy, uruchamiają własne instancje Subagent Network — łącznie do 30 jednostek wykonawczych działających jednocześnie | Executor Chat (×2) |
| 7 | Operator spina proces z modułem Automations — harmonogram cykliczny, Queue Manager, Orchestrator, Execution Monitor | Panel orkiestracji ▸ Harmonogram i automatyki |
| 8 | Always On Display przechodzi w tryb operatora; Mobile pozostaje kanałem interwencji na czas pracy ciągłej | Monitor procesu |
| 9 | Proces działa autonomicznie 24 godziny na dobę, 7 dni w tygodniu, 365 dni w roku, z pełną widocznością każdego kroku w Monitorze procesu i możliwością interwencji użytkownika w dowolnej chwili | Cała pętla (rozdz. 11.1) |

---

## 16. Zgodność z zasadami nadrzędnymi platformy

| Zasada nadrzędna | Realizacja w środowisku MultitaskingAI |
|---|---|
| 1. Pełna kompozycyjność | Role, tryby współpracy, akcje silnika kolejek i zależności orkiestracji są komponentami swobodnie składanymi w dowolną konfigurację zespołu — szablony w rozdziale 14 są przykładami takich kompozycji, nie zamkniętym katalogiem |
| 2. Pełna konfigurowalność (zasada centralna) | Przypisanie ról, tryb współpracy wykonawców, zasięg i akcje kolejek, zależności orkiestracji, integracja z Automations oraz punkty izolacji na poziomie roli konfiguruje się w całości z okna konfiguracji i z panelu orkiestracji; brak ustawienia oznacza wartość domyślną |
| 3. Jawność i konfigurowalność zależności | Powiązanie silnika kolejek MultitaskingAI z silnikiem kolejek Automations oraz zależności zdefiniowane w orkiestracji są jawnymi, odwracalnymi decyzjami użytkownika, nie regułami wbudowanymi na stałe |
| 4. Rozszerzenie orkiestracji (cel) | MultitaskingAI jest obecnym, najpełniejszym wcieleniem tego kierunku rozwoju — do dwóch równoległych wykonawców, każdy z możliwością uruchomienia do 15 podagentów, co daje potencjalnie do 30 jednostek wykonawczych działających jednocześnie w ramach jednego procesu |
| Klucze jawne | Pasek widoczności kontekstu (rozdz. 20), dziennik zdarzeń środowiska (rozdz. 21) i objaśnienia kontekstowe czynią jawnym, co która rola widzi i co dzieje się w procesie — środowisko nie zawiera ukrytej maszynerii |
| Stopniowe ujawnianie funkcjonalności | Warstwy widoczności (rozdz. 1.6) utrzymują wizualną prostotę interfejsu przy wielu równoległych wątkach: widoczna jest warstwa 1, pozostałe warstwy pozostają zwinięte i osiągalne jednym kliknięciem, skrótem albo poleceniem w Chat Window |
| Dwa kanały komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca) są elementami pierwszoplanowymi środowiska, obecnymi w każdej przestrzeni roboczej i w każdym układzie kolumn (rozdz. 1.7, 9.2) |
| Zero blokerów (dyrektywa nadrzędna) | Żaden mechanizm opisany w niniejszym dokumencie — hierarchia decyzji, ocena Executora 3 / Validatora, izolacja, uprawnienia — nie jest bramą blokującą. Wszystkie są ustawieniami konfiguracyjnymi ze stanem wyjściowym „wyłączone” lub „pełny dostęp”; domyślne zachowanie systemu to wykonanie polecenia. Funkcje powłoki — skrócone sterowanie (rozdz. 19.1), akcje kolejek (rozdz. 19.2), zamknięcie karty sesji (rozdz. 18) i współdzielenie kontekstu (rozdz. 20) — pozostają zawsze klikalne; przy niepasującym stanie wyświetlają komunikat kontekstowy, a akcje nieodwracalne mają „Cofnij" |

---

## 17. Powłoka środowiska — nawigacja i przestrzeń robocza okien

Powłoka środowiska jest ramą operacyjną orkiestracji wieloagentowej. Rozdziały 17–21 opisują jej funkcje: nawigację, przestrzeń roboczą okien ról, karty sesji, panel stanu procesu, skróty i paletę poleceń, współdzielenie kontekstu oraz funkcje wspólne. Każda funkcja działa bez blokujących walidacji — brak ustawienia oznacza wartość domyślną, nigdy przerwanie pracy.

### 17.1. Nawigacja środowiska

Nawigację środowiska tworzą: panel orkiestracji jako kolumna boczna (rozdz. 6), pasek kontekstu procesu ze znacznikami kontekstowymi i kartami procesów oraz Chat Window, z którego poleceniem języka naturalnego osiągalna jest każda sekcja i każde okno.

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Plakietki liczbowe sekcji panelu | Na pozycjach nawigacyjnych sześciu sekcji panelu orkiestracji wyświetlane są żywe liczniki: zadań w kolejce, aktywnych ról, naruszonych zależności, oczekujących decyzji | Orientacja o stanie procesu bez wchodzenia w sekcję | Pozycja nawigacyjna sekcji (rozdz. 6.4); kanał WebSocket | 1 | Widoczne bez interakcji |
| Zwijanie panelu do ikon | Panel orkiestracji zwęża się do samych ikon sekcji, powiększając kolumny okien ról; stan zapamiętywany per proces | Więcej szerokości dla kolumn okien ról przy szerokich procesach | Przycisk zwinięcia panelu (rozdz. 6.4) | 1 | Ikona `☰` w nagłówku panelu |
| Kolejność i widoczność sekcji | Przeciąganie pozycji oraz przełączniki widoczności sekcji; „Przywróć domyślne" pozostaje dostępne zawsze | Dopasowanie panelu do charakteru procesu bez utraty funkcji | Uchwyt przeciągania, przełącznik widoczności (rozdz. 6.4, 6.5) | 4 | Okno konfiguracji panelu orkiestracji |
| Szybki skok między procesami | Nakładka z miniaturami wszystkich otwartych kart sesji, wybór klawiaturą lub myszą | Nawigacja między wieloma równoległymi zespołami bez sięgania do paska kart | Karta sesji (rozdz. 10.1); rozdz. 19.2 | 3 | Skrót klawiszowy `Ctrl+Tab` |
| Ścieżka nawigacyjna | Ciąg znaczników kontekstowych: Środowisko ▸ Proces ▸ Sekcja ▸ Rola / okno; każdy segment klikalny | Świadomość położenia w gęstym środowisku orkiestracji | Pasek kontekstu procesu (rozdz. 10.1) | 1 | Widoczna bez interakcji |
| Przełącznik środowisk | Wyjście z MultitaskingAI do TalkIn, WorkSpace, CodeStudio lub strony głównej; procesy MultitaskingAI działają dalej w tle | Płynne przechodzenie między środowiskami bez utraty stanu procesów | interfejs-uzytkownika/strona-glowna-i-nawigacja.md (rozdz. 8.2–8.3); trwałość sesji | 2 | Znacznik kontekstowy `[MultitaskingAI]` w pasku kontekstu |
| Nawigacja klawiaturą po sekcjach | Klawisze 1–6 przełączają wprost do sekcji panelu; strzałki poruszają się po pozycjach list | Praca bez myszy przy sterowaniu procesem | Rozdz. 19.2 | 4 | Skrót klawiszowy |
| Kotwice do ról w Monitorze procesu | Kliknięcie statusu roli w sekcji Monitor procesu przenosi wprost do jej okna roboczego z podświetleniem | Szybkie dojście od diagnozy do miejsca interwencji | Widok hierarchii decyzji, plakietka statusu roli (rozdz. 10.4, 10.9) | 1 | Kliknięcie plakietki statusu |
| Menu kontekstowe pozycji nawigacyjnej | Zestaw akcji sekcji: otwarcie w oknie oderwanym, przypięcie, ukrycie, konfiguracja | Skrócenie częstych operacji na sekcjach | Pozycja nawigacyjna (rozdz. 6.4); rozdz. 17.2 | 3 | Menu kebab (⋮) pozycji, prawy przycisk myszy |

### 17.2. Przestrzeń robocza okien ról

Przestrzeń robocza środowiska rozmieszcza okna w kolumnach sąsiadujących poziomo. Chat Window zajmuje lewą kolumnę o pełnej wysokości, Execution Loop Window kolumnę sąsiadującą, okna ról kolumnę dominującą, a panele pomocnicze otwierają się jako kolejne rozszerzenia boczne. Wielokolumnowa przestrzeń robocza jest cechą wyróżniającą MultitaskingAI spośród środowisk jednookiennych.

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Menedżer układu kolumn | Przełącza rozmieszczenie okien ról między układami kolumnowymi: jedna kolumna (skupienie), dwie kolumny, trzy kolumny, cztery kolumny (cztery role), trzy kolumny wykonawcze z kolumną koordynacji | Dopasowanie szerokości kolumn do składu zespołu i etapu pracy | Makieta główna (rozdz. 9.3); karty ról (rozdz. 10.4) | 2 | Element zbiorczy `Układ ▼` w pasku kontekstu procesu |
| Zapisywane układy kolumn | Zapis bieżącego rozmieszczenia, szerokości i zwinięć kolumn jako nazwany preset; wczytanie jednym kliknięciem; presety wspólne z presetami zespołu | Powrót do sprawdzonego układu pracy dla danego typu procesu | Sekcja Zespoły (rozdz. 6.2, 14); architektura/model-konfiguracji.md | 3 | Rozwinięcie `Układ ▼ ▸ Presety` |
| Zmiana szerokości kolumn | Przeciąganie granicy między kolumnami zmienia ich szerokość; podwójne kliknięcie granicy wyrównuje szerokości | Ręczne doważenie uwagi między wykonawcami a koordynacją | Makieta główna (rozdz. 9.3) | 1 | Przeciągnięcie granicy kolumn |
| Tryb skupienia na roli | Wybrana kolumna roli rozszerza się na całą szerokość obszaru roboczego, pozostałe zwijają się do wąskich kolumn ze wskaźnikami statusu; Chat Window i Execution Loop Window pozostają widoczne | Głęboka praca nad jednym torem bez utraty podglądu pozostałych | Karta roli, plakietka statusu roli (rozdz. 10.4) | 2 | Ikona skupienia w nagłówku kolumny, skrót klawiszowy |
| Zwijanie kolumny roli | Kolumna roli zwija się do wąskiej listwy z nazwą roli i statusem; kliknięcie rozwija ją ponownie | Chwilowe uproszczenie widoku przy pracy nad wybranymi rolami | Nagłówek okna roli (rozdz. 9.1) | 1 | Kliknięcie nagłówka kolumny |
| Panel Subagent Network jako rozszerzenie boczne | Panel do 15 podagentów otwiera się jako osobna kolumna boczna obok okna wykonawcy, z listą kart podagentów | Czytelny podgląd wielu podagentów przy pełnej sieci (15/15) | Panel Subagent Network (rozdz. 10.5, makieta 9.7) | 3 | Rozwinięcie `Subagent Network ▼` |
| Okna oderwane | Okno roli lub sekcja wypinane są do osobnego okna systemowego — praca na wielu monitorach; okno oderwane zachowuje układ kolumnowy | Wykorzystanie drugiego ekranu do rozdzielenia wykonania od nadzoru | Klient Tauri (architektura/architektura.md); trwałość stanu sesji | 3 | Menu kebab (⋮) nagłówka okna |
| Dokowanie i przeciąganie kolumn | Przeciągnięcie nagłówka okna roli w inne miejsce układu kolumnowego, ze strefami zrzutu i podglądem pozycji docelowej | Swobodna aranżacja przestrzeni bez sztywnego porządku kolumn | Menedżer układu kolumn; makieta główna (rozdz. 9.3) | 2 | Przeciągnięcie nagłówka kolumny |
| Sprzężenie przewijania porównania | W Results Analyzer przewijanie dwóch kolumn wyników (Executor 1 i Executor 2) jest sprzężone | Rzetelne porównanie rozbieżnych wyników linia po linii | Panel porównania wyników (rozdz. 10.7) | 3 | Przełącznik w nagłówku panelu porównania |
| Zakładki wewnątrz okna roli | W jednej kolumnie roli działają zakładki: historia, plan, pamięć, log — bez mnożenia kolumn | Więcej treści roli w tej samej szerokości | Komponent `.dn-zakladki` (Załącznik A.2) | 1 | Kliknięcie zakładki |
| Migawka układu | Bieżące rozmieszczenie kolumn i statusów eksportowane jest jako obraz lub schemat tekstowy | Dokumentowanie konfiguracji zespołu i dzielenie się układem | Rozdz. 21; interfejs-uzytkownika/system-wizualny.md | 3 | Menu kebab (⋮) paska kontekstu procesu |

---

## 18. Karty sesji — procesy orkiestracji i trwałość stanu

Karta sesji reprezentuje jeden proces orkiestracji z własnym zespołem ról, kolejkami, zależnościami i historią przebiegów. Stan procesu utrzymywany jest po stronie serwera niezależnie od otwartego klienta.

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nazywanie i kolorowanie kart sesji | Karcie procesu nadawana jest nazwa zespołu, kolor i ikona; kolor przenosi się na panel stanu procesu | Rozróżnienie wielu równoległych procesów na pierwszy rzut oka | Karta sesji (rozdz. 10.1) | 3 | Menu kebab (⋮) karty |
| Przypinanie kart | Istotny proces przypinany jest na początku paska kart; karta przypięta nie zamyka się przypadkowo | Ochrona długotrwałych procesów ciągłych przed zamknięciem | Karta sesji (rozdz. 10.1) | 3 | Menu kebab (⋮) karty |
| Grupowanie kart sesji | Powiązane procesy łączone są w nazwaną grupę (na przykład „Projekt X — backend, frontend, QA") ze wspólnym zwijaniem | Porządek przy dziesiątkach równoległych zespołów | Pasek kontekstu procesu (rozdz. 10.1) | 3 | Menu kebab (⋮) karty, przeciągnięcie karty na kartę |
| Migawka procesu | Pełny stan procesu — role, kolejki, zależności, historia — zapisywany jest jako punkt przywracania; powrót następuje jednym kliknięciem | Bezpieczne eksperymentowanie z konfiguracją zespołu i powrót do znanego stanu | Trwałość stanu po stronie serwera (rozdz. 9.1); architektura/model-danych.md | 3 | Menu kebab (⋮) karty, paleta poleceń |
| Wznawianie procesu po rozłączeniu | Karta sesji odtwarza pełny stan po zamknięciu klienta lub utracie połączenia, ponieważ stan procesu żyje po stronie serwera | Praca ciągła 24/7/365 bez zależności od otwartego okna | Trwałość sesji (rozdz. 9.1); integracja z Automations (rozdz. 7) | 1 | Automatyczne przy otwarciu karty |
| Duplikowanie procesu | Karta sesji klonowana jest z konfiguracją zespołu, bez historii przebiegów — jako punkt startu wariantu | Szybki start kolejnego procesu z tego samego wzorca | Sekcja Zespoły (rozdz. 14) | 3 | Menu kebab (⋮) karty |
| Porównywanie dwóch przebiegów | Dwa przebiegi tego samego procesu zestawiane są w dwóch kolumnach: czas, koszt, wynik, log decyzji | Ocena, która konfiguracja zespołu działa lepiej | Tabela przebiegów, Monitor procesu (rozdz. 10.9) | 3 | Zaznaczenie dwóch wierszy tabeli przebiegów, akcja „Porównaj" |
| Wskaźnik aktywności w tle na karcie | Na karcie nieaktywnego procesu plakietka pulsuje przy pracy w tle: nowy wynik, oczekująca decyzja, błąd | Świadomość zdarzeń w procesach, które nie są na wierzchu | Karta sesji, stan „z aktywnością w tle" (rozdz. 10.1) | 1 | Widoczny bez interakcji |
| Zamknięcie z zachowaniem stanu | Zamknięcie karty nie kończy procesu po stronie serwera; toast „Cofnij" przywraca widok; trwałe zakończenie procesu jest osobną, świadomą akcją | Zero utraty pracy przy przypadkowym zamknięciu | Trwałość sesji; komponent `.dn-toast` (Załącznik A.2) | 1 | Ikona zamknięcia karty |
| Lista wszystkich procesów | Lista wszystkich otwartych i uśpionych procesów z filtrem po stanie: uruchomiony, wstrzymany, ciągły | Zarządzanie liczbą procesów przekraczającą pojemność paska kart | Stany przebiegu (rozdz. 12.1); rozdz. 17.1 | 3 | Element zbiorczy `Procesy ▼` w pasku kontekstu |

---

## 19. Panel stanu procesu, skróty klawiszowe i paleta poleceń

### 19.1. Panel stanu procesu

Panel stanu procesu jest wąską kolumną boczną po prawej stronie obszaru roboczego, zbierającą wskaźniki bieżącego procesu i skrócone sterowanie. Panel otwiera się jako rozszerzenie boczne i zwija do listwy znaczników.

```
 Makieta — Panel stanu procesu (kolumna boczna), stan spoczynku
 ┌──────────────────────────┐
 │ Stan procesu         [⋮] │
 ├──────────────────────────┤
 │ Przebieg   ● uruchomiony │
 │ Kolejka    4 zadania     │
 │ Automations ● spięte     │
 │ Najbliższy cykl   22:00  │
 │ Jednostki  4 role ·      │
 │            6/30 podagen. │
 │ Zależności ● zgodne      │
 │ AOD        obserwator ▼  │
 │ Połączenie ● na żywo     │
 ├──────────────────────────┤
 │ [ Sterowanie ▼ ]         │
 └──────────────────────────┘
```

| Segment | Zawartość | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Segment kolejek | Liczba zadań w kolejce, w tym wstrzymanych; kliknięcie otwiera sekcję Kolejki | Stały wgląd w obciążenie kolejki bez przełączania sekcji | Tabela kolejek (rozdz. 4.5); silnik kolejek (rozdz. 4) | 1 | Widoczny bez interakcji |
| Segment integracji Automations | Plakietka „spięte / niespięte / błąd" oraz czas najbliższego cyklu; kliknięcie otwiera sekcję Harmonogram i automatyki | Natychmiastowa wiedza, czy proces działa w trybie ciągłym | Plakietka statusu połączenia (rozdz. 7.6) | 1 | Widoczny bez interakcji |
| Segment stanu przebiegu | Plakietka stanu pętli: nieuruchomiony, uruchomiony, wstrzymany, oczekuje na decyzję, błąd — wraz z kropką stanu | Jeden punkt prawdy o stanie procesu | Stany przebiegu procesu (rozdz. 12.1) | 1 | Widoczny bez interakcji |
| Sterowanie skrócone | Zestaw Start, Pauza, Wznów, Stop dostępny bez przechodzenia do Coordinator Chat i Execution Loop Window | Sterowanie procesem z każdego miejsca środowiska | Zestaw przycisków sterowania procesem (rozdz. 10.6, 3.3) | 3 | Element zbiorczy `Sterowanie ▼` |
| Licznik kosztu i zużycia | Bieżący koszt i zużycie tokenów przebiegu oraz sumarycznie, z rozbiciem per rola po najechaniu | Kontrola budżetu procesu wieloagentowego (do 30 jednostek wykonawczych) | architektura/integracja-modeli.md; kanały modeli (rozdz. 3.1) | 2 | Najechanie na segment, rozwinięcie szczegółów |
| Licznik jednostek wykonawczych | Liczba aktywnych ról i podagentów, na przykład „4 role · 6/30 podagentów" | Świadomość skali równoległości i obciążenia | Subagent Network (rozdz. 3.2) | 1 | Widoczny bez interakcji |
| Wskaźnik zgodności zależności | Plakietka „zgodny / oczekujący / naruszony"; kliknięcie przy naruszeniu otwiera szczegóły konfliktu | Wczesne wykrycie złamania reguły orkiestracji | Plakietka statusu zgodności (rozdz. 5.4) | 1 | Widoczny bez interakcji |
| Wskaźnik trybu Always On Display | Tryb obserwatora lub operatora; kliknięcie przełącza tryb | Jawność, kto może interweniować w proces | Przełącznik trybu AOD (rozdz. 2.6, 10.9) | 2 | Znacznik `AOD ▼` |
| Wskaźnik przebiegu pętli wykonawczej | Numer bieżącego cyklu pętli, liczba zadań w dekompozycji zlecenia i stan kontroli jakości | Bieżąca orientacja w przebiegu pętli prowadzonej w Execution Loop Window | Execution Loop Window (rozdz. 1.7, 9.5) | 1 | Widoczny bez interakcji |
| Konfiguracja segmentów panelu | Włączenie, wyłączenie i kolejność segmentów panelu stanu procesu | Dopasowanie panelu do tego, co w danym procesie istotne | Rozdz. 22; architektura/model-konfiguracji.md | 4 | Okno konfiguracji |
| Wskaźnik połączenia na żywo | Stan kanału WebSocket: na żywo, ponawianie, rozłączony — bez blokowania pracy | Zaufanie do aktualności statusów ról i kolejek | architektura/architektura.md (WebSocket/JSON) | 1 | Widoczny bez interakcji |

### 19.2. Skróty klawiszowe i paleta poleceń

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Paleta poleceń środowiska | Nakładka z wyszukiwaniem wszystkich poleceń powłoki: przełączenie sekcji, uruchomienie i wstrzymanie procesu, dodanie kolejki, przypisanie agenta do roli, zmiana układu kolumn, wczytanie zespołu, otwarcie Execution Loop Window | Dostęp do każdej funkcji powłoki bez szukania po menu | Wszystkie sekcje panelu orkiestracji | 4 | Skrót `Ctrl/Cmd + K`, polecenie w Chat Window |
| Skróty sterowania procesem | Klawisze do operacji Start, Pauza, Wznów, Stop, Waliduj oraz `retry` ostatniego zadania | Szybkie sterowanie pętlą bez myszy | Akcje silnika kolejek (rozdz. 4.2); rozdz. 19.1 | 4 | Skrót klawiszowy |
| Skróty akcji kolejek | Przypisania klawiszowe do jedenastu akcji silnika kolejek na zaznaczonym zadaniu | Sprawna praca w gęstej sekcji Kolejki | Zestaw przycisków akcji (rozdz. 4.5) | 4 | Skrót klawiszowy |
| Skróty przełączania okien | Klawisze do skoku między Chat Window, Execution Loop Window, Executor 1, Executor 2, Coordinator i Validatorem oraz do trybu skupienia | Płynne przechodzenie między torami pracy | Karty ról (rozdz. 10.4); rozdz. 17.2 | 4 | Skrót klawiszowy |
| Skróty warstwy centralnej | Wywołanie dymka Always On Display, mikrofonu i przełącznika trybu AOD | Natychmiastowy dostęp do nadzoru i komunikacji głosowej | Elementy AOD (rozdz. 2.6) | 4 | Skrót klawiszowy |
| Edytor skrótów | Okno konfiguracji, w którym każdy skrót powłoki jest podglądalny i przypisywalny na nowo; „Przywróć domyślne" pozostaje dostępne | Dopasowanie skrótów do nawyków Operatora | Rozdz. 22; architektura/model-konfiguracji.md | 4 | Okno konfiguracji |
| Ściągawka skrótów | Nakładka z listą wszystkich aktywnych skrótów, pogrupowaną tematycznie | Przypomnienie skrótów bez opuszczania środowiska | Edytor skrótów | 4 | Skrót klawiszowy |

---

## 20. Współdzielenie kontekstu między rolami, procesami i modułami

Mostki kontekstu są jawne i sterowane ustawieniami konfiguracyjnymi, spójnymi z siedmioma poziomami zasięgu izolacji (rozdz. 13). Stan wyjściowy: kontekst odrębny per proces i per rola; każde współdzielenie jest świadomą decyzją użytkownika, nigdy domyślnym otwarciem.

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Współdzielenie kontekstu między rolami procesu | Historia, pamięć lub kontekst udostępniane są jawnie wybranym rolom tego samego procesu, na przykład Executorowi 1 i Executorowi 2 w trybie iteracyjnym | Wspólna baza wiedzy wykonawców bez ręcznego kopiowania | Okno konfiguracji izolacji, poziom „Rola" (rozdz. 13); tryby współpracy (rozdz. 3.4) | 4 | Okno konfiguracji punktów izolacji |
| Most kontekstu między procesami | Jednorazowe przekazanie wyniku lub pamięci z jednego procesu orkiestracji do drugiego, na przykład wyniku badania do procesu redakcji | Łańcuchowanie procesów bez eksportu ręcznego | Poziom zasięgu „karta sesji" (rozdz. 13.1); architektura/model-danych.md | 3 | Menu kebab (⋮) karty sesji, polecenie w Chat Window |
| Wciągnięcie kontekstu projektu z modułu Workspace | Pamięć kontekstowa i pliki projektu z modułu Workspace podpinane są jako wejście dla ról procesu | Role pracują na wiedzy projektu bez powielania materiałów | Moduł Workspace; pamięć poziomu „projekt" (rozdz. 3.1) | 3 | Rozwinięcie `Kontekst ▼` w nagłówku okna roli |
| Przekazanie wyniku do modułu docelowego | Rezultat procesu kierowany jest do modułu tematycznego (Apps, Studio, Library) jako artefakt, jednym poleceniem z okna roli lub z Chat Window | Domknięcie drogi „proces → produkt" bez wychodzenia ze środowiska | Powiązania modułów (architektura/koncepcja-platformy.md); akcja `route` (rozdz. 4.2) | 3 | Menu kebab (⋮) wyniku, polecenie w Chat Window |
| Pasek widoczności kontekstu | Przy każdej roli pokazywane są źródła kontekstu, z których rola korzysta: pamięć projektu, historia współdzielona, wynik innej roli — jawnie i klikalnie | Pełna przejrzystość tego, co widzi dana rola | Wskaźnik dostępu do pamięci (rozdz. 10.4); zasada kluczy jawnych | 1 | Widoczny bez interakcji |
| Odpięcie kontekstu jednym ruchem | Dowolne współdzielenie cofane jest natychmiast, z powrotem do stanu odrębnego; toast zawiera „Cofnij" | Odwracalność każdej decyzji o współdzieleniu | Komponent `.dn-toast` | 2 | Kliknięcie źródła w pasku widoczności kontekstu |
| Współdzielona tablica ustaleń procesu | Wspólny, lekki notatnik procesu widoczny dla wszystkich ról oraz dla Always On Display — ustalenia, blokady, decyzje | Jedno miejsce prawdy dla całego zespołu w toku pętli | Obszar historii (rozdz. 10.4); Monitor procesu (rozdz. 6.2); Execution Loop Window (rozdz. 9.5) | 3 | Rozwinięcie `Tablica ustaleń ▼` w Execution Loop Window |

---

## 21. Funkcje wspólne środowiska

| Funkcja | Działanie | Wartość dla użytkownika | Powiązania | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Wyszukiwanie globalne środowiska | Jedno pole odnajduje role, zadania, kolejki, zależności, przebiegi i procesy w bieżącym oraz w pozostałych procesach; wyniki pogrupowane | Szybkie dotarcie do dowolnego elementu w gęstym środowisku | Pole wyszukiwania (rozdz. 10.1); indeks po stronie rdzenia | 2 | Ikona `szukaj`, skrót klawiszowy |
| Centrum powiadomień | Mechanizm w postaci jednolitej dla całej platformy — `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6. Specyfika MultitaskingAI: zakres źródeł obejmuje wszystkie procesy orkiestracji równocześnie, a filtr źródła zawęża listę do wskazanego procesu, roli albo kolejki | Nic istotnego nie umyka przy wielu równoległych procesach | Plakietka powiadomień AOD (rozdz. 2.6); powiadomienia wypychane Mobile (rozdz. 7.6); encja `powiadomienie` — `architektura/model-danych.md`, rozdz. 18.4 | 3 | Plakietka powiadomień z licznikiem w pasku kontekstu |
| Dziennik zdarzeń środowiska | Pełny, chronologiczny, nieograniczony log każdej akcji powłoki i silnika kolejek per przebieg, z eksportem | Rozliczalność i diagnoza procesu wieloagentowego | Przycisk „Szczegóły przebiegu" (rozdz. 10.9); zasada pełnej widoczności | 4 | Paleta poleceń, sekcja Monitor procesu |
| Motyw i gęstość interfejsu | Motyw jasny lub ciemny oraz gęstość interfejsu: komfortowa albo kompaktowa — dla widoku wielokolumnowego | Czytelność przy wielu kolumnach naraz i oszczędność szerokości | Przełącznik motywu (rozdz. 10.1); interfejs-uzytkownika/system-wizualny.md | 3 | Menu kebab (⋮) paska kontekstu, okno konfiguracji |
| Tryb skupienia środowiska | Wycisza powiadomienia i zwija kolumny pomocnicze, pozostawiając Chat Window, Execution Loop Window, kolumny ról i panel stanu procesu | Głęboka praca nad procesem bez rozproszeń | Centrum powiadomień; tryb skupienia na roli (rozdz. 17.2) | 2 | Znacznik `Skupienie ▼`, skrót klawiszowy |
| Eksport i import konfiguracji środowiska | Cała konfiguracja powłoki — układy kolumn, skróty, segmenty panelu stanu, sekcje — zapisywana jest i wczytywana jako plik przenośny | Przeniesienie ustawień między stanowiskami i kopia zapasowa | architektura/model-konfiguracji.md; rozdz. 17.2, 19, 22 | 4 | Okno konfiguracji |
| Tryb prezentacji | Widok tylko do oglądania pełnej pętli — makieta główna wraz z Monitorem procesu — do pokazania procesu na żywo bez ryzyka przypadkowej akcji | Demonstracja i nadzór z ekranu współdzielonego | Makieta główna (rozdz. 9.3); Monitor procesu | 3 | Paleta poleceń, menu kebab (⋮) paska kontekstu |
| Widok Mobile procesu | Zwarty układ powłoki dla urządzenia mobilnego: stan procesu, skrócone sterowanie, zatwierdzanie zdalne; kolumny układają się jedna pod drugą wyłącznie w tym widoku, zachowując porządek Chat Window → Execution Loop Window → okno roli | Interwencja i nadzór spoza stanowiska roboczego w trybie ciągłym | Przycisk zatwierdzenia zdalnego (rozdz. 7.6); funkcja Mobile (rozdz. 2.4) | 1 | Otwarcie procesu na urządzeniu mobilnym |
| Objaśnienia kontekstowe | Przy każdym elemencie konfiguracji powłoki ikona `info` z wyjaśnieniem działania i wpływu ustawienia | Obniżenie progu wejścia do najbardziej złożonego środowiska platformy | Ikona objaśnienia (rozdz. 13.5); interfejs-uzytkownika/system-wizualny.md | 2 | Najechanie lub kliknięcie ikony `[?]` |
| Pusty stan z podpowiedzią startu | Nowy proces bez zespołu pokazuje kroki startowe — wczytanie zespołu, przypisanie ról, dodanie kolejki — zamiast pustej przestrzeni | Prowadzenie użytkownika od zera bez ograniczania swobody | Komponent `.dn-pusty-stan` (Załącznik A.2); sekcja Zespoły | 1 | Widoczny bez interakcji |

---

## 22. Punkty sterowania z okna konfiguracji

Wszystkie poniższe zakresy konfiguruje się z okna konfiguracji platformy; brak ustawienia oznacza wartość domyślną, nigdy blokadę. Poziom zasięgu odwołuje się do siedmiu poziomów hierarchii izolacji (globalny → środowisko → moduł → para modułów → projekt → karta sesji → rola), gdzie poziom węższy ma pierwszeństwo.

| Punkt sterowania | Czym steruje | Poziom zasięgu | Wartość domyślna |
|---|---|---|---|
| Kolejność i widoczność sekcji panelu | Układ sześciu sekcji panelu orkiestracji (rozdz. 17.1) | środowisko / karta sesji | sześć sekcji, kolejność z rozdz. 6.5 |
| Stan zwinięcia panelu nawigacji | Panel pełny albo ikonowy (rozdz. 17.1) | karta sesji | rozwinięty |
| Domyślny układ kolumn okien ról | Układ startowy nowego procesu (rozdz. 17.2) | środowisko / zespół | cztery kolumny albo układ z presetu zespołu |
| Presety układu kolumn | Nazwane rozmieszczenia kolumn (rozdz. 17.2) | zespół / projekt | brak — tworzone przez Operatora |
| Szerokość kolumny Chat Window | Szerokość lewej kolumny okna komunikacji Użytkownik ↔ Wykonawca | środowisko / karta sesji | szerokość standardowa kolumny komunikacji |
| Otwarcie i szerokość kolumny Execution Loop Window | Czy kolumna pętli wykonawczej jest rozwinięta przy otwarciu procesu i jaką ma szerokość | środowisko / karta sesji | rozwinięta, szerokość standardowa |
| Zakres treści Execution Loop Window | Które elementy pętli są prezentowane: dekompozycja zlecenia, kolejka i stan zadań, komunikaty sterujące, wyniki kontroli jakości, wskaźniki przebiegu | proces | wszystkie elementy widoczne |
| Segmenty i kolejność panelu stanu procesu | Które wskaźniki widoczne w panelu stanu (rozdz. 19.1) | środowisko | wszystkie segmenty widoczne |
| Skróty klawiszowe | Przypisania klawiszy powłoki (rozdz. 19.2) | globalny / środowisko | zestaw domyślny |
| Współdzielenie kontekstu | Zakres i kierunek mostków kontekstu (rozdz. 20) | rola / karta sesji | odrębny kontekst per proces i per rola |
| Tryb Always On Display | Obserwator albo operator (rozdz. 2.2) | proces / środowisko | obserwator |
| Powiadomienia i ich waga | Które klasy zdarzeń zgłasza centrum powiadomień (rozdz. 21; `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6) | środowisko / proces | wszystkie klasy, waga normalna |
| Motyw i gęstość | Motyw jasny lub ciemny, gęstość komfortowa lub kompaktowa (rozdz. 21) | globalny | motyw platformy, gęstość komfortowa |
| Potwierdzenia akcji nieodwracalnych | Modal przed usunięciem zespołu, kolejki albo procesu | globalny / środowisko | wyłączone (akcja od razu + „Cofnij") |
| Profil izolacji roli | Osiem zakresów izolacji technicznej oraz izolacja kontekstu (rozdz. 13.2) | rola (pierwszeństwo najwyższe) | wyłączone / pełny dostęp |
| Powiązanie z modułem Automations | Spięcie silnika kolejek ze Scheduler i Queue Manager (rozdz. 7) | proces | niespięte |
| Widoczność źródeł kontekstu roli | Czy pasek widoczności kontekstu jest pokazywany (rozdz. 20) | środowisko | widoczny |
| Warstwy widoczności elementów | Które elementy warstw 2–4 są ujawniane w interfejsie danej roli użytkownika (rozdz. 1.6) | rola / środowisko | warstwa 1 widoczna, warstwy 2–4 zwinięte |

**Szablon konfiguracji powłoki środowiska:**

```
srodowisko_multitaskingai:
  panel_orkiestracji:
    zwiniety: nie
    sekcje: [Zespoly, Role, Kolejki, Orkiestracja, "Harmonogram i automatyki", "Monitor procesu"]
  uklad_kolumn:
    domyslny: "cztery kolumny"
    chat_window: "lewa, stala, pelna wysokosc"
    execution_loop_window: "kolumna sasiadujaca, rozwinieta"
    presety: []                 # tworzone przez Operatora, wspolne z presetami zespolu
  panel_stanu_procesu:
    segmenty: [przebieg, kolejki, automations, petla, koszt, jednostki, zaleznosci, aod, polaczenie]
  skroty: domyslne              # personalizowalne, "Przywroc domyslne" zawsze dostepne
  warstwy_widocznosci:
    warstwa_1: widoczna
    warstwy_2_4: zwiniete       # wyzwalacze: znaczniki kontekstowe, ▼, ⋮, ☰
  wspoldzielenie_kontekstu:
    domyslne: odrebny           # per proces i per rola; mostki wlaczane swiadomie
    pierwszenstwo: rola         # profil roli wygrywa z ustawieniem karty sesji
  always_on_display:
    tryb: obserwator            # przelaczalny na operatora w dowolnej chwili
  potwierdzenia_nieodwracalne: wylaczone   # akcja od razu + "Cofnij"
  # brak dowolnego klucza = wartosc domyslna, nigdy blokada
```

---

## 23. Słownik pojęć

**Always On Display** — globalny agent towarzyszący, nieposiadający własnego środowiska ani modułu, obecny jednocześnie we wszystkich częściach platformy; w kontekście MultitaskingAI pełni funkcję warstwy centralnej jako obserwator lub operator procesu (rozdz. 2).

**Executor** — rola odpowiedzialna za faktyczne wykonanie pracy, bez odpowiedzialności za zarządzanie procesem (rozdz. 3.1, 3.4).

**Subagent Network** — mechanizm pozwalający wykonawcy uruchomić jednocześnie do 15 wyspecjalizowanych podagentów, których wyniki są następnie agregowane przez wykonawcę macierzystego (rozdz. 3.2).

**Coordinator** — rola odpowiedzialna za planowanie, podział pracy, budowę promptów, sterowanie procesem oraz zarządzanie kolejką (rozdz. 3.3).

**Validator** — jedno z przykładowych wcieleń roli Executor 3, odpowiedzialne za kontrolę jakości pracy pozostałych modeli (rozdz. 3.5).

**Silnik kolejek** — mechanizm pozwalający definiować kolejki na pięciu poziomach zasięgu i obsługujący jedenaście akcji (rozdz. 4).

**Orkiestracja** — mechanizm definiujący zależności pomiędzy modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami oraz zapewniający ich respektowanie w toku wykonania (rozdz. 5).

**Panel orkiestracji** — boczna nawigacja środowiska MultitaskingAI, złożona z sześciu sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu (rozdz. 6).

**Hierarchia decyzji** — porządek pierwszeństwa decyzyjnego między użytkownikiem, Always On Display, Coordinatorem, Executorem 3 / Validatorem, wykonawcami i Subagent Network, konfigurowalny i możliwy do pominięcia przez użytkownika w każdej chwili (rozdz. 8).

**Executor Chat** — okno robocze Executora 1 lub Executora 2 (rozdz. 9.6).

**Coordinator Chat** — okno robocze Coordinatora (rozdz. 9.8).

**Results Analyzer** — okno robocze Executora 3 / Validatora (rozdz. 9.9).

**Izolacja kontekstu** i **izolacja techniczna** — dwa rodzaje izolacji, przypisywalne na poziomie roli z pierwszeństwem najwyższym w hierarchii zasięgów (rozdz. 13).

Pozostałe pojęcia platformowe (Środowisko, Moduł, Komponent własny, Czat, Sesja, Projekt, Centrum dowodzenia, Agent, Mobile) zachowują brzmienie ustalone w słowniczku Koncepcji platformy i nie są w niniejszym dokumencie powtarzane.

---

## Załącznik A. Katalog ikon i komponentów systemu wizualnego wykorzystanych w środowisku

### A.1. Ikony

| Ikona | Zastosowanie w środowisku MultitaskingAI |
|---|---|
| `uzytkownik` | Awatar/profil ról; sekcja Role panelu orkiestracji |
| `gwiazdka` | Sekcja Zespoły (presety) |
| `filtr` / `strzalka-prawo` | Sekcja Kolejki (przepływ zadań) |
| `link-zewnetrzny` | Sekcja Orkiestracja (zależność) |
| `zegar` / `kalendarz` | Sekcja Harmonogram i automatyki |
| `oko` | Sekcja Monitor procesu; wcielenie Reviewer |
| `tarcza` | Wcielenie Security Auditor; Permissions Center |
| `waga` | Wcielenie Arbitrator |
| `ptaszek` / `ptaszek-kolo` | Wcielenie QA Lead; potwierdzenie, walidacja |
| `dokument` / `kod` | Wcielenia Product Owner / Architect |
| `odswiez` | Akcja `retry` |
| `uruchom` | Akcja `resume`, start procesu |
| `zatrzymaj` | Akcja `pause` |
| `plus` | Nowa kolejka, nowa karta sesji, dodanie podagenta lub zależności |
| `kosz` | Usunięcie zespołu, kolejki, zależności |
| `dzwonek` | Powiadomienia Always On Display |
| `koperta` | Powiadomienia (kontekst Assistant/Mobile) |
| `menu` | Zwinięcie/rozwinięcie panelu orkiestracji |
| `szukaj` | Pole wyszukiwania list i pozycji |
| `spinacz` | Załącznik w polu poleceń Executor Chat / Coordinator Chat |
| `wyslij` | Wysłanie polecenia |
| `info` / `ostrzezenie` / `blad` | Plakietki i alerty stanu procesu, roli, zadania |
| `slonce` / `ksiezyc` | Przełącznik motywu |
| `globus` | Dostęp sieciowy w macierzy izolacji |
| `klodka` | Konto i token per sesja w macierzy izolacji |
| `folder` | Katalog roboczy sesji w macierzy izolacji |

Pełna galeria wszystkich 47 ikon z podglądem SVG: `ikony/indeks.html`.

### A.2. Komponenty

| Komponent | Klasa bazowa | Zastosowanie w środowisku |
|---|---|---|
| Przycisk | `.dn-btn` (`--glowny`, `--zloty`, `--zarys`, `--duch`, `--blad`) | Sterowanie procesem, akcje kolejki, decyzje Results Analyzer |
| Przycisk ikonowy | `.dn-btn-ikona` | Jedenaście akcji silnika kolejek, sterowanie Coordinatora |
| Pole formularza | `.dn-pole`, `.dn-input`, `.dn-select`, `.dn-textarea` | Selektory przypisania, kanału, priorytetu, pola poleceń |
| Przełącznik | `.dn-suwak` | Macierz izolacji, widoczność sekcji panelu, powiązanie z Automations |
| Karta | `.dn-karta` (`--interaktywna`, `--akcent`, `--pozycja`) | Karty ról, karty podagentów, pozycje list |
| Plakietka | `.dn-plakietka` (`--stan`, `--zloto`, `--sukces`, `--ostrz`, `--blad`, `--info`) | Statusy ról, przebiegów, zgodności zależności |
| Kropka stanu | `.dn-kropka` | Status roli i przebiegu wewnątrz plakietki |
| Zakładki | `.dn-zakladki` | Karty sesji, zakładki wewnętrzne okien |
| Tabela | `.dn-tabela`, `.dn-tabela-owijka` | Tabela kolejek, tabela przebiegów (Monitor procesu) |
| Awatar | `.dn-awatar` (`--kwadrat` dla agenta, kołowy dla modelu bazowego, `-stan`) | Awatar ról, awatar Always On Display |
| Modal | `.dn-modal` | Potwierdzenie akcji nieodwracalnej (usunięcie zespołu, kolejki) — ustawienie konfiguracyjne włączane przez Operatora, domyślnie wyłączone, akcja wykonuje się od razu |
| Toast | `.dn-toast` | Powiadomienia interwencji (Mobile), potwierdzenia zapisu, opcja „Cofnij” po akcjach nieodwracalnych |
| Pusty stan | `.dn-pusty-stan` | Sekcja Kolejki, Zespoły lub Orkiestracja przed pierwszą konfiguracją |
| Spinner | `.dn-spinner` | Stan ładowania akcji silnika kolejek |

---

## Załącznik B. Macierz zgodności ze źródłami

| Rozdział niniejszego dokumentu | architektura/koncepcja-platformy.md | architektura/architektura.md | architektura/model-danych.md | architektura/kontrakty-komunikacji.md | specyfikacje/specyfikacja-okien-operacyjnych.md | specyfikacje/specyfikacja-agentow.md | interfejs-uzytkownika/system-wizualny.md | architektura/izolacja-i-zaleznosci.md |
|---|---|---|---|---|---|---|---|---|
| 1. Miejsce w platformie, warstwy widoczności i dwa kanały komunikacji | rozdz. 3.4, 5, 6, 7.4, 8, 9.8, 9.14, 11.1, 11.3–11.5, 12 | rozdz. 7, 8 | rozdz. 7.2, 7.4 | rozdz. 9.1, 9.2 | rozdz. 2.4 | rozdz. 2 | — | — |
| 2. Warstwa centralna | rozdz. 10.1, 10.2, 11.2 | — | — | — | — | — | rozdz. 8 | — |
| 3. Role | rozdz. 9.15, 11.3, Zał. B, D.2 | rozdz. 8 | rozdz. 8.2, 8.3 | rozdz. 5.8 | rozdz. 7.2 | rozdz. 6, 8 | — | — |
| 4. Silnik kolejek | rozdz. 9.3, 11.3, 11.4, 13.4 | — | rozdz. 10.3 | rozdz. 11.3 | — | — | — | — |
| 5. Orkiestracja | rozdz. 11.5, 13.5 | — | rozdz. 10.7 | rozdz. 11.4 | — | — | — | — |
| 6. Panel orkiestracji | rozdz. 6, 8, 11.8, 13 pkt 2 | — | — | — | — | — | rozdz. 13 | — |
| 7. Integracja z Automations | rozdz. 9.3, 10.1, 11.6 | — | — | — | rozdz. 6.3 | — | — | — |
| 8. Hierarchia decyzji | rozdz. 10.1, 10.2, 11.2, 11.3 | — | — | — | — | — | — | — |
| 9. Makiety okien ról | Zał. A | — | — | — | rozdz. 7 | rozdz. 6.2 | rozdz. 13 | — |
| 10. Katalog elementów | Zał. D | — | — | — | rozdz. 2.4, 3 | rozdz. 3, 4 | rozdz. 8, 13 | — |
| 11. Diagramy przepływu | rozdz. 11.9 | — | — | rozdz. 11.3, 11.4 | — | — | — | — |
| 12. Stany | — | — | — | — | — | — | rozdz. 6 | — |
| 13. Izolacja na poziomie roli | rozdz. 4.4–4.6 | rozdz. 9, 10, 12 | rozdz. 13.2, 15.5, 17.1 | rozdz. 15 | rozdz. 8 | rozdz. 4 | rozdz. 12 | rozdz. 3, 4, 6, 7 |
| 14. Szablony zespołu | rozdz. 12 pkt 1, Zał. D.2, E.1 | — | — | — | — | rozdz. 8 | — | rozdz. 7.2 |
| 15. Scenariusze | rozdz. 10.1, 10.2, 11.3, Zał. E.1 | — | — | — | — | — | — | rozdz. 13 |
| 16. Zgodność z zasadami | rozdz. 12 | — | rozdz. 1 | — | — | rozdz. 10 | rozdz. 14 | rozdz. 1 |
| 17. Powłoka — nawigacja i przestrzeń robocza | rozdz. 6, 8, 11.8 | rozdz. 7 | — | — | rozdz. 2, 7 | — | rozdz. 13 | — |
| 18. Karty sesji i trwałość stanu | rozdz. 9.8, 11.1 | — | rozdz. 7.2 | — | rozdz. 7 | — | — | rozdz. 7.3 |
| 19. Panel stanu, skróty, paleta poleceń | rozdz. 9.3, 11.4 | — | — | — | rozdz. 2.4 | — | rozdz. 8 | — |
| 20. Współdzielenie kontekstu | rozdz. 4.4–4.6 | rozdz. 12 | rozdz. 15.5 | rozdz. 15 | rozdz. 8 | rozdz. 4 | — | rozdz. 3, 4, 6, 7 |
| 21. Funkcje wspólne środowiska | rozdz. 8, 10.1 | — | — | — | rozdz. 3 | — | rozdz. 6, 8, 13 | — |
| 22. Punkty sterowania z okna konfiguracji | rozdz. 12, 13 | — | — | — | — | rozdz. 4 | — | rozdz. 7 |

---

*Koniec dokumentu. Środowisko MultitaskingAI — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
