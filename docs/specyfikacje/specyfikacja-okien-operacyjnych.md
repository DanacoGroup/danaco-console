# Danaco Console — Specyfikacja okien operacyjnych

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
| **Źródło** | Koncepcja platformy i architektura |

Dokument stanowi katalog wszystkich typów okien operacyjnych występujących w piętnastu modułach platformy oraz w środowisku MultitaskingAI. Dla każdego okna określono cel, zawartość, zachowanie, warstwę widoczności, sposób wywołania oraz powiązane operacje. Dokument porządkuje okna wymienione w Koncepcji platformy (rozdział 11, rozdział 13, Załącznik A) i w Architekturze (rozdział 5), stosując się do zasady nadrzędnej pełnej konfigurowalności i braku twardych blokad.

---

## Spis treści

1. [Wprowadzenie](#1-wprowadzenie)
2. [Pojęcie okna operacyjnego](#2-pojęcie-okna-operacyjnego)
3. [Zasady wspólne dla okien operacyjnych](#3-zasady-wspólne-dla-okien-operacyjnych)
3a. [Warstwy widoczności — stopniowe ujawnianie funkcjonalności](#3a-warstwy-widoczności--stopniowe-ujawnianie-funkcjonalności)
4. [Klasyfikacja okien według charakteru pracy](#4-klasyfikacja-okien-według-charakteru-pracy)
5. [Okna komunikacji operacyjnej — Chat Window i Execution Loop Window](#5-okna-komunikacji-operacyjnej--chat-window-i-execution-loop-window)
6. [Katalog okien operacyjnych per moduł](#6-katalog-okien-operacyjnych-per-moduł)
   - 6.0. [Konwencja makiet i typologia kolumn](#60-konwencja-makiet-i-typologia-kolumn)
   - 6.1. [Studio](#61-studio)
   - 6.2. [Workspace](#62-workspace)
   - 6.3. [Automations](#63-automations)
   - 6.4. [Browser](#64-browser)
   - 6.5. [Research](#65-research)
   - 6.6. [Library](#66-library)
   - 6.7. [Translate](#67-translate)
   - 6.8. [Roundtable](#68-roundtable)
   - 6.9. [Design](#69-design)
   - 6.10. [Assistant](#610-assistant)
   - 6.11. [Terminal](#611-terminal)
   - 6.12. [Developer](#612-developer)
   - 6.13. [Diagnostics](#613-diagnostics)
   - 6.14. [Apps](#614-apps)
   - 6.15. [Agents](#615-agents)
7. [Okna robocze środowiska MultitaskingAI](#7-okna-robocze-środowiska-multitaskingai)
8. [Okno konfiguracji punktów izolacji](#8-okno-konfiguracji-punktów-izolacji)
9. [Macierz zbiorcza okien operacyjnych per moduł](#9-macierz-zbiorcza-okien-operacyjnych-per-moduł)
10. [Słowniczek i zasady nazewnictwa okien](#10-słowniczek-i-zasady-nazewnictwa-okien)

---

## 1. Wprowadzenie

### 1.1. Cel dokumentu

Dokument odpowiada na pytanie, jaką konkretną przestrzeń roboczą użytkownik otrzymuje po wejściu w każdy z piętnastu modułów platformy oraz w cztery role środowiska MultitaskingAI. O ile Koncepcja platformy (rozdział 11) opisuje moduły przez ich cel, funkcjonalności i powiązania, o tyle niniejsza specyfikacja schodzi o jeden poziom niżej — do pojedynczego okna operacyjnego — i dla każdego z nich ustala cel, zawartość, zachowanie, warstwę widoczności, sposób wywołania oraz powiązane operacje w sposób jednolity dla całej platformy.

| Poziom opisu | Dokument | Jednostka opisu |
|---|---|---|
| Cel, funkcjonalności, powiązania modułu | Koncepcja platformy, rozdział 11 | moduł |
| Cel, zawartość, zachowanie, warstwa widoczności, sposób wywołania, powiązane operacje | Niniejsza specyfikacja, rozdziały 5–8 | pojedyncze okno operacyjne |

### 1.2. Zakres

| Status | Element zakresu |
|---|---|
| Obejmuje | Katalog okien operacyjnych piętnastu modułów platformy (rozdział 6), zgodnie z Załącznikiem A Koncepcji platformy |
| Obejmuje | Chat Window oraz Execution Loop Window jako okna komunikacji operacyjnej wspólne wszystkim modułom (rozdział 5) |
| Obejmuje | Warstwy widoczności okien i elementów interfejsu wraz ze sposobem ich wywołania (rozdział 3a) |
| Obejmuje | Okna robocze czterech ról środowiska MultitaskingAI (rozdział 7) |
| Obejmuje | Okno konfiguracji punktów izolacji jako okno o zasięgu globalnym, powiązane z pracą we wszystkich modułach (rozdział 8) |
| Nie obejmuje | Szczegóły projektu wizualnego okien: wymiary, tokeny kolorów i typografii oraz warstwa graficzna — przedmiot przewodnika marki |
| Nie obejmuje | Model danych okien — przedmiot odrębnego dokumentu (Model danych) |
| Nie obejmuje | Kontrakty komunikacji klient–serwer — przedmiot odrębnego dokumentu (Kontrakty komunikacji) |

Makiety tekstowe zamieszczone w rozdziałach 5–8 przedstawiają **schematyczne, funkcjonalne rozmieszczenie okien w kolumnach** — ich sąsiedztwo i porządek od lewej do prawej — a nie warstwę graficzną. Makiety rysowane są w stanie spoczynku interfejsu, zgodnie z regułą warstw widoczności (rozdział 3a): widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

### 1.3. Źródła

| Dokument | Wykorzystane części |
|---|---|
| Koncepcja platformy | Rozdziały 1–14 oraz Załączniki A–E; w szczególności rozdział 9 (Moduły), rozdział 11 (MultitaskingAI — rozwinięcie środowiska), Załącznik A (Macierz okien operacyjnych per moduł) |
| Architektura | Rozdział 5 (Interfejs użytkownika), rozdział 9 (Izolacja) |

Wszystkie nazwy własne — środowisk, modułów, okien, ról — przejęto bez zmian z Koncepcji platformy.

---

## 2. Pojęcie okna operacyjnego

### 2.1. Definicja

Okno operacyjne jest pojedynczym elementem interfejsu wchodzącym w skład zestawu okien danego modułu — stanowi właściwą przestrzeń wykonywania zadań przypisaną do konkretnej funkcjonalności modułu (Koncepcja platformy, rozdział 2.2 i Słowniczek). Każdy moduł otwiera własny, dedykowany zestaw okien operacyjnych, zastępujący zestaw okien poprzedniego modułu w tej samej karcie sesji (rozdział 2.2 i rozdział 8 Koncepcji platformy).

### 2.2. Uwaga terminologiczna: dwa znaczenia „okna operacyjnego”

Nazwa „okno operacyjne” funkcjonuje w dokumentacji platformy w dwóch odrębnych znaczeniach, celowo zbieżnych jedynie brzmieniowo — analogicznie do rozróżnienia między warstwą modułów (Workspace Layer) a modułem Workspace, wyjaśnionego w rozdziale 5.2 Koncepcji platformy.

| Znaczenie | Źródło | Definicja | Liczność |
|---|---|---|---|
| Okno operacyjne aplikacji | Architektura, rozdział 13 | Jedno z dwóch okien aplikacji klienckiej, obok okna konfiguracji: główne środowisko pracy, w którym toczą się rozmowy, sesje i panele funkcjonalne, ze zmianą konfiguracji bieżącej sesji przez uproszczone menu kontekstowe | Jedno okno systemowe aplikacji |
| Okno operacyjne modułu | Koncepcja platformy, rozdział 2.2, Słowniczek, Załącznik A | Pojedynczy element zestawu okien danego modułu, na przykład Studio Editor, Workflow Builder lub Code Editor | Wiele okien, po kilka na moduł |

Niniejszy dokument posługuje się wyłącznie drugim znaczeniem — okno operacyjne jako element zestawu okien modułu. Wszystkie okna operacyjne katalogowane w rozdziałach 5–8 renderują się wewnątrz okna operacyjnego aplikacji w pierwszym znaczeniu (Architektura, rozdział 13), które stanowi ich wspólną powłokę techniczną.

### 2.3. Miejsce w hierarchii platformy

Okno operacyjne jest najniższym poziomem trójstopniowej hierarchii organizacji pracy przyjętej przez platformę (Koncepcja platformy, Wprowadzenie i rozdział 3). Poniższy schemat przedstawia tę hierarchię wraz z rolą obu okien komunikacji operacyjnej na poziomie okna operacyjnego.

```
   ŚRODOWISKO                MODUŁ                   OKNO OPERACYJNE
 (ogólny tryb pracy)  →  (konkretne zadanie   →   (narzędzia realizacji
                          biznesowe)                zadania)
        │                     │                          │
        └─────────────────────┴──────────────┬───────────┘
                                              │
                        Chat Window  (Użytkownik ↔ Wykonawca)
                        Execution Loop Window  (Koordynator ↔ Wykonawca)
              oba okna obecne w każdym module; Chat Window prowadzi
              komunikację użytkownika z Wykonawcą, Execution Loop Window
              prowadzi pętlę wykonawczą, niezależnie od tego, które
              z pozostałych okien modułu jest aktualnie używane
```

Środowisko wyznacza ogólny tryb pracy, moduł doprecyzowuje go do konkretnego zadania biznesowego, a zestaw okien operacyjnych danego modułu dostarcza narzędzi do jego realizacji.

### 2.4. Anatomia okna operacyjnego

Poniższy schemat i tabela porządkują elementy wspólne każdego okna operacyjnego — wynikające bezpośrednio z zasad wspólnych z rozdziału 3 oraz z reguły warstw widoczności z rozdziału 3a. Stanowią wzorzec anatomii, do którego odnoszą się szczegółowe katalogi okien w rozdziałach 5–8.

```
 ┌─ Nagłówek okna ─────────────────────────────────────────────
 │   nazwa okna · przynależność do bieżącego modułu · warstwa
 ├─────────────────────────────────────────────────────────────
 │   OBSZAR ZAWARTOŚCI
 │     treść właściwa oknu (edytor, lista, monitor, podgląd …)
 │     aktualizacja na żywo kanałem WebSocket — dotyczy okien
 │     monitorujących procesy oraz okien strumienia odpowiedzi
 ├─────────────────────────────────────────────────────────────
 │   WYZWALACZE WARSTW 2–3 (stan spoczynku)
 │     znaczniki kontekstowe · `▼` · `⋮` · `☰`
 │     rozwinięcie jednym kliknięciem, zwinięcie po użyciu
 ├─────────────────────────────────────────────────────────────
 │   STAN PROCESU SESJI (po stronie serwera)
 │     trwały — rozłączenie klienta nie zamyka okna;
 │     izolacja kontekstu domyślnie odrębna per karta sesji,
 │     zakres współdzielenia konfigurowalny (rozdział 8)
 └─────────────────────────────────────────────────────────────
```

| Element anatomii | Rola | Podstawa (zasada wspólna) |
|---|---|---|
| Nagłówek okna | Identyfikacja okna, jego przynależności do modułu i warstwy widoczności | Rozdział 2.1, rozdział 3a |
| Obszar zawartości | Prezentacja treści właściwej oknu; aktualizacja na żywo tam, gdzie okno monitoruje proces lub odbiera strumień odpowiedzi | Zasada 6 (synchronizacja na żywo) |
| Wyzwalacze warstw 2–3 | Udostępnienie ustawień i akcji ukrytych w stanie spoczynku, osiągalnych jednym kliknięciem | Rozdział 3a |
| Elementy konfiguracji z objaśnieniem `[?]` | Udostępnienie ustawień z opisem działania i wpływu na aplikację | Zasada 4 (objaśnienia kontekstowe) |
| Stan procesu sesji | Trwałość okna po rozłączeniu klienta; przywrócenie po ponownym połączeniu | Zasada 5 (trwałość stanu procesu sesji) |
| Domyślna izolacja kontekstu | Odrębna historia, pamięć i kontekst per karta sesji; współdzielenie konfigurowalne | Zasada 2 (konfigurowalna izolacja kontekstu) |

Konfigurowalną strukturę zestawu okien modułu — do której odnoszą się katalogi rozdziału 6 — porządkuje poniższy szablon redakcyjny. Wszystkie pola podlegają zasadzie „brak ustawienia oznacza wartość domyślną”.

```
zestaw_okien_modułu:
  moduł:               <nazwa modułu>              # np. Studio
  środowiska_dostępne: [ <środowiska> ]            # wg macierzy dostępności (Koncepcja, rozdz. 10)
  okna_komunikacji:
    - Chat Window            # lewa kolumna stała  (Użytkownik ↔ Wykonawca)
    - Execution Loop Window  # kolumna sąsiadująca (Koordynator ↔ Wykonawca)
  okna_operacyjne:
    - nazwa:                  <nazwa okna>         # np. Studio Editor
      cel:                    <do czego służy okno>
      zawartość:              <co okno prezentuje>
      zachowanie:             <kiedy i jak okno się aktualizuje>
      warstwa:                <1 | 2 | 3 | 4>           # rozdział 3a
      sposób_wywołania:       <widoczne bez interakcji | znacznik | przycisk |
                               menu ⋮ | menu ☰ | polecenie języka naturalnego>
      kolumna:                <lewa stała | sąsiadująca | dominująca | boczna>
      izolacja_domyślna:      odrębna_per_karta_sesji   # konfigurowalna (rozdział 8)
      powiązane_operacje:     [ <operacje> ]
      powiązania_międzymodułowe: [ <okna innych modułów> ]  # jawne, konfigurowalne (Koncepcja, rozdz. 6)
```

---

## 3. Zasady wspólne dla okien operacyjnych

Poniższe zasady obowiązują jednolicie dla wszystkich okien operacyjnych katalogowanych w niniejszym dokumencie, niezależnie od modułu, w którym się znajdują.

| Zasada | Treść | Podstawa |
|---|---|---|
| Układ pionowy lewa–prawa | Wszystkie okna rozmieszczone są w kolumnach sąsiadujących poziomo: lewa kolumna stała — Chat Window; kolumna sąsiadująca — Execution Loop Window; kolumna dominująca — obszar roboczy modułu; kolejne kolumny boczne — okna pomocnicze otwierane jako rozszerzenia boczne. Regulacji podlega wyłącznie szerokość kolumn | Rozdział 6.0 |
| Przeładowanie przy zmianie modułu | Zestaw okien operacyjnych jest właściwością modułu, nie karty sesji. Przełączenie modułu w obrębie jednej karty sesji zamyka okna poprzedniego modułu i otwiera zestaw okien nowego modułu | Koncepcja platformy, rozdział 2.2 i rozdział 8 |
| Konfigurowalna izolacja kontekstu | Historia, pamięć i kontekst każdego okna są domyślnie odrębne dla każdej nowej karty sesji; zakres współdzielenia między kartami, modułami, środowiskami i projektami jest sterowany z okna konfiguracji punktów izolacji. Żadne okno nie narzuca izolacji jako reguły sztywnej | Rozdział 8; Koncepcja platformy, rozdział 6 |
| Brak twardych blokad | Żadne okno operacyjne nie ogranicza na stałe dostępu do funkcji platformy. Brak ustawienia w oknie konfiguracji oznacza wartość domyślną, a nie zablokowaną funkcję | Koncepcja platformy, rozdział 14, zasada 2 |
| Objaśnienia kontekstowe | Elementy konfiguracji dostępne z okna operacyjnego lub jego okna kontekstowego opatrzone są objaśnieniem oznaczonym znakiem `[?]`, opisującym działanie ustawienia i jego wpływ na aplikację; obejmuje to w szczególności macierz izolacji z rozdziału 8 | Architektura, rozdział 13 |
| Trwałość stanu procesu sesji | Każde okno działa w ramach procesu sesji po stronie serwera. Rozłączenie klienta nie kończy sesji ani nie zamyka stanu okna — stan jest trwały, a ponowne połączenie przywraca okno w stanie, w jakim zostało pozostawione | Architektura, rozdział 7 |
| Synchronizacja na żywo | Zawartość okna aktualizuje się w czasie rzeczywistym kanałem WebSocket — dotyczy to w szczególności okien monitorujących procesy (Execution Monitor, Process Monitor, Build Output, Actions Monitor) oraz okien odbierających strumień odpowiedzi modelu (Chat Window, Execution Loop Window i pochodne) | Architektura, rozdział 11 i 16 |
| Wspólność okien komunikacji | Chat Window i Execution Loop Window występują w każdym z piętnastu modułów — stanowią tę samą warstwę komunikacji operacyjnej, rekonfigurowaną do kontekstu bieżącego modułu | Rozdział 5 |
| Przypisanie warstwy widoczności | Każde okno, panel i pojedyncza funkcja należy do dokładnie jednej z czterech warstw widoczności i ma określony sposób wywołania | Rozdział 3a |

---

## 3a. Warstwy widoczności — stopniowe ujawnianie funkcjonalności

Zasadą nadrzędną interfejsu Danaco Console jest reguła: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

Każdy element interfejsu opisywany w niniejszej specyfikacji należy do dokładnie jednej z czterech warstw widoczności. Przy opisie okna, panelu, paska narzędzi i pojedynczej funkcji podaje się jej warstwę oraz sposób wywołania — w kolumnach „Warstwa” i „Sposób wywołania” tabel rozdziałów 5–9.

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, Execution Loop Window, aktywne okno wiodące obszaru roboczego, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy, panele pomocnicze i monitory otwierane jako rozszerzenia boczne | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w oknie komunikacji, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Mechanizmy ukrywania funkcjonalności stosowane w oknach operacyjnych platformy:

- **Menu progresywne.** Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`), a lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Funkcjonalność umieszczana jest w panelach bocznych, panelach wysuwanych, oknach popover i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji.
- **Znaczniki kontekstowe.** Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

Makiety w rozdziałach 5–8 rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 — lewa kolumna Chat Window, kolumna sąsiadująca Execution Loop Window, kolumna dominująca z oknem wiodącym modułu — oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są w tabelach okien z podaniem warstwy i sposobu wywołania.

---

## 4. Klasyfikacja okien według charakteru pracy

Poniższa klasyfikacja porządkuje okna operacyjne katalogowane w rozdziałach 5–7 według charakteru wykonywanej w nich pracy. Klasyfikacja ma charakter porządkujący i pomocniczy — nie jest elementem struktury platformy, a jedynie ułatwia orientację w katalogu.

| Kategoria | Charakter pracy | Warstwa | Przykładowe okna |
|---|---|---|---|
| Komunikacja użytkownika | Rozmowa z Wykonawcą, wydawanie poleceń, zatwierdzanie i przerywanie działań | 1 | Chat Window, Voice Console, Executor Chat |
| Pętla wykonawcza | Koordynacja zadań, nadzór, orkiestracja, kontrola realizacji | 1 | Execution Loop Window, Coordinator Chat |
| Edycja treści | Tworzenie i przekształcanie materiału | 1 | Studio Editor, Code Editor, Design Board, Source Panel, Translation Panels |
| Podgląd i porównanie | Prezentacja wyniku, zestawienie wersji | 2 | Preview Window, Diff/Grep Panel, Build Output, File Preview |
| Repozytorium i biblioteka | Trwałe przechowywanie materiałów | 2–3 | Session Repository, Library Explorer, Project Library, Versioning Panel |
| Konstruktor / builder | Budowa struktury lub procesu | 1–2 | Workflow Builder, Agent Builder, Product Builder, Prompt Builder, Architecture Designer |
| Kolejka i harmonogram | Planowanie i sterowanie wykonaniem w czasie | 2–3 | Scheduler, Queue Manager, Orchestrator |
| Monitor procesu | Obserwacja przebiegu wykonania na żywo | 2 | Execution Monitor, Process Monitor, Actions Monitor, Activity Feed, Diagnostics Center, Results Analyzer |
| Panel narzędziowy | Zestaw operacji lub ustawień pomocniczych | 3 | Tools Panel, Glossary Manager, Tags & Collections, Permissions Center |
| Zarządzanie źródłami i ustaleniami | Gromadzenie i porządkowanie materiału badawczego | 2 | Sources Panel, Sources Manager, Findings Panel, Notes Panel |
| Konfiguracja | Ustalanie parametrów pracy modułu, agenta lub sesji | 2–4 | Instructions Panel, Model Configuration, Connectors Manager, okno konfiguracji punktów izolacji |

---

## 5. Okna komunikacji operacyjnej — Chat Window i Execution Loop Window

Platforma prowadzi dwa kanały komunikacji operacyjnej. Oba są elementami pierwszoplanowymi architektury interfejsu, należą do warstwy 1 i występują w każdym z piętnastu modułów oraz w środowisku MultitaskingAI. Zajmują dwie pierwsze kolumny układu: lewą kolumnę stałą i kolumnę z nią sąsiadującą.

### 5.1. Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy). Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę.

| Aspekt anatomii | Charakterystyka |
|---|---|
| Kanał | Użytkownik ↔ Wykonawca |
| Warstwa | 1 — zawsze widoczna |
| Sposób wywołania | Widoczne bez interakcji; lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Cel | Centralny punkt pracy użytkownika i podstawowy mechanizm sterowania procesami platformy; przyjmowanie poleceń w języku naturalnym, prezentacja wyników, zatwierdzanie i przerywanie działań, wyjaśnianie wyniku i kontekstu |
| Zawartość | Historia rozmowy; pole wprowadzania poleceń; strumień odpowiedzi modelu przekazywany na żywo (Architektura, rozdział 11); pasek znaczników kontekstowych; dostęp do operacji właściwych bieżącemu modułowi |
| Zachowanie | Przy każdej zmianie środowiska lub modułu rekonfigurowane są: wygląd, możliwości, historia, dostępne narzędzia i kontekst (Koncepcja platformy, rozdział 2.4). Domyślnie każda nowa karta sesji otrzymuje odrębną historię i pamięć czatu; zakres współdzielenia między kartami, modułami lub środowiskami sterowany jest z okna konfiguracji punktów izolacji (rozdział 8) |
| Powiązane operacje | Wydawanie poleceń Wykonawcy właściwych bieżącemu modułowi; zatwierdzanie i przerywanie działań; odbiór wyników generowanych przez inne okna modułu; wywołanie funkcji warstwy 4 poleceniem języka naturalnego |

### 5.2. Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window jest oknem pętli wykonawczej prezentującym komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów.

| Aspekt anatomii | Charakterystyka |
|---|---|
| Kanał | Koordynator ↔ Wykonawca |
| Warstwa | 1 — zawsze widoczna |
| Sposób wywołania | Widoczne bez interakcji; kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |
| Cel | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór nad realizacją, orkiestracja działań, kontrola realizacji procesów |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania; kolejka i stan zadań; wymiana komunikatów sterujących; wyniki kontroli jakości i decyzje o ponowieniu; wskaźniki przebiegu pętli |
| Zachowanie | Aktualizuje się na żywo kanałem WebSocket w miarę przebiegu pętli; rekonfiguruje zakres zadań przy zmianie modułu; stan pętli jest trwały po stronie procesu sesji |
| Powiązane operacje | Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia; ponowienie zadania; przegląd dekompozycji zlecenia; przegląd wyników kontroli jakości |

### 5.3. Role uczestników kanałów

| Rola | Zakres odpowiedzialności |
|---|---|
| Użytkownik | Zleca i zatwierdza |
| Koordynator | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela i nadzoruje zadania |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania |

### 5.4. Rozmieszczenie okien komunikacji

Miejsce obu okien komunikacji w układzie kolumnowym dowolnego modułu przedstawia makieta w stanie spoczynku.

```
 Makieta — okna komunikacji w układzie kolumnowym dowolnego modułu
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window          │ Execution Loop       │ Obszar roboczy
  nawigacja  │ Użytkownik ↔         │ Window               │ modułu
  modułów    │ Wykonawca            │ Koordynator ↔        │ (okno wiodące
             │                      │ Wykonawca            │  — katalog
             │  historia rozmowy    │                      │  rozdziału 6)
             │  pole poleceń        │  zlecenie i jego     │
             │  strumień odpowiedzi │  dekompozycja        │
             │  na żywo (WebSocket) │  kolejka i stan zadań│
             │                      │  komunikaty sterujące│
             │  [Danaco Console]    │  kontrola jakości    │
             │  [TalkIn] [Fable 5]▼ │  wskaźniki pętli   ⋮ │            ☰
 ═══════════════════════════════════════════════════════════════════════════
```

Przy zmianie modułu żadne z okien komunikacji nie znika — zmienia się kontekst poleceń, zakres zadań pętli i dostępne operacje.

| Moduł | Kontekst Chat Window po rekonfiguracji | Kontekst Execution Loop Window po rekonfiguracji |
|---|---|---|
| Studio | Operacje kontekstowe wobec dokumentu (na przykład wynik operacji z Tools Panel) | Dekompozycja zlecenia redakcyjnego na zadania, kontrola jakości wyniku |
| Automations | Polecenia wobec kolejki zadań automatyki | Nadzór nad przebiegami workflow, ponowienia i korekty zlecenia |
| Developer | Polecenia wobec repozytorium powiązanego z sesją | Koordynacja zadań budowania, testów i poprawek |

---

## 6. Katalog okien operacyjnych per moduł

Poniższe podrozdziały katalogują okna operacyjne piętnastu modułów platformy w kolejności przyjętej w rozdziale 11 Koncepcji platformy. Środowiska, w których dany moduł jest dostępny, wskazuje macierz dostępności modułów (Koncepcja platformy, rozdział 10) i przywołano je w nagłówku każdego modułu jako „Dostępność”. Chat Window i Execution Loop Window, opisane w rozdziale 5, występują w każdym module i zajmują dwie pierwsze kolumny każdej makiety; tabele obejmują pozostałe okna właściwe danemu modułowi.

### 6.0. Konwencja makiet i typologia kolumn

Każdy moduł opisano trzema elementami: **tabelą anatomii okien** (cel, zawartość, zachowanie, warstwa, sposób wywołania, powiązane operacje — po jednym wierszu na okno), **makietą układu kolumn** (rozmieszczenie okien od lewej do prawej) oraz **tabelą powiązań międzymodułowych**.

Obowiązuje wyłącznie układ pionowy — podział lewa–prawa. Typologia wag wizualnych okien wyrażona jest w kolumnach.

| Kolumna | Waga wizualna | Zawartość | Warstwa |
|---|---|---|---|
| Boczna nawigacja | Kolumna skrajnie lewa, stała, wąska | Lista modułów; dla środowiska MultitaskingAI zastępuje ją panel orkiestracji (rozdział 7) | 1 |
| Lewa kolumna komunikacji | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Chat Window — kanał Użytkownik ↔ Wykonawca | 1 |
| Kolumna sąsiadująca | Kolumna sąsiadująca z lewą kolumną komunikacji, pełna wysokość | Execution Loop Window — kanał Koordynator ↔ Wykonawca | 1 |
| Kolumna dominująca | Kolumna o największej szerokości | Obszar roboczy modułu: okno wiodące, edytory, podglądy, monitory wywołane do pracy | 1 |
| Kolejne kolumny boczne | Kolumna boczna, otwierana jako rozszerzenie boczne | Okna pomocnicze i panele, po prawej stronie obszaru roboczego; po zamknięciu znikają z przestrzeni roboczej | 2–3 |

Regulacji podlega wyłącznie **szerokość** kolumn. Makiety rysowane są w stanie spoczynku: kolumny boczne warstw 2–3 nie są rozwinięte, a ich wyzwalacze zaznaczono jako `▼`, `⋮`, `☰` po prawej stronie kolumny dominującej. Nazwa okna podana jest przed znakiem `│`, skrót jego funkcji po nim.

### 6.1. Studio

Moduł zaawansowanej pracy z tekstem, dokumentami i treścią (Koncepcja platformy, rozdział 11.1). **Dostępność:** TalkIn, WorkSpace.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Studio Editor | 1 | Widoczne bez interakcji — kolumna dominująca | Edycja treści dokumentu | Treść dokumentu w formatach PDF, DOCX, TXT, Markdown | Wczytuje dokument z Library lub z urządzenia; umożliwia zaznaczenie fragmentu jako przedmiotu operacji | Wczytanie dokumentu, edycja, zaznaczenie fragmentu, zapis |
| Tools Panel | 3 | Menu zbiorcze `Operacje ▼` — rozszerzenie boczne | Dostęp do operacji kontekstowych właściwych dokumentowi | Lista operacji: korekta, zmiana stylu, tłumaczenie, streszczenie, rozwinięcie treści | Operacje stosowane do zaznaczonego fragmentu lub całego dokumentu, zlecane z Chat Window | Wywołanie operacji na fragmencie z Studio Editor |
| Diff/Grep Panel | 2 | Przycisk porównania — rozszerzenie boczne | Porównanie wersji dokumentu i przeszukiwanie treści | Widok różnic przed i po zmianie; wyniki wyszukiwania wzorca w treści | Otwiera się po wykonaniu operacji lub na żądanie użytkownika; po zamknięciu znika z przestrzeni roboczej | Akceptacja lub odrzucenie zmiany, przeszukiwanie dokumentu |
| Preview Window | 2 | Przełącznik podglądu — rozszerzenie boczne | Podgląd wynikowego dokumentu | Wyrenderowana treść w docelowym formacie wyjściowym | Aktualizuje się po zmianach zaakceptowanych w Studio Editor | Eksport dokumentu, przekazanie do modułu Library |
| Session Repository | 3 | Menu kebab (⋮) — rozszerzenie boczne | Historia wersji pracy nad dokumentem w bieżącej sesji | Kolejne wersje dokumentu ze znacznikami czasu | Narasta w toku sesji; nie usuwa wersji wcześniejszych bez decyzji użytkownika | Powrót do wcześniejszej wersji, przegląd historii zmian |

```
 Makieta — Moduł: Studio          Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Studio Editor      │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  treść dokumentu   │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  PDF · DOCX · TXT  │ ☰
             │                   │                    │  Markdown          │
             │ historia rozmowy  │ zlecenie i zadania │  zaznaczenie       │
             │ pole poleceń      │ kolejka i stan     │  fragmentu jako    │
             │ strumień na żywo  │ kontrola jakości   │  przedmiot operacji│
 ═══════════════════════════════════════════════════════════════════════════
   kolumna     lewa kolumna        kolumna              kolumna dominująca   kolumny
   stała       stała               sąsiadująca                               boczne
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Dalsze przetwarzanie wyników pracy | Sources Manager, Report Builder (moduł Research) |
| Przechowywanie wyników | Library Explorer (moduł Library) |

### 6.2. Workspace

Moduł izolowanego środowiska projektowego (Koncepcja platformy, rozdział 11.2). **Dostępność:** TalkIn, WorkSpace, CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Project Dashboard | 1 | Widoczne bez interakcji — kolumna dominująca | Ogólny widok stanu projektu | Zestawienie zadań, statusu prac, powiązanych zasobów projektu | Otwiera się jako punkt wejścia po wybraniu projektu | Przegląd stanu projektu, nawigacja do pozostałych okien modułu |
| Project Library | 2 | Znacznik projektu — rozszerzenie boczne | Repozytorium plików ograniczone do projektu | Dokumenty i artefakty właściwe jednemu projektowi | Odpowiednik modułu Library ograniczony do zakresu projektu | Dodanie pliku, przegląd zasobów projektu |
| Agent Manager | 2 | Element zwinięty `Agent ▼` — rozszerzenie boczne | Przypisywanie agentów do projektu | Lista agentów skonfigurowanych w module Agents, dostępnych do przypisania | Agenci przypisani domyślnie działają w zakresie projektu | Przypisanie agenta jako Wykonawcy zadania w projekcie |
| Instructions Panel | 3 | Menu kebab (⋮) — rozszerzenie boczne | Definiowanie instrukcji systemowych właściwych projektowi | Treść instrukcji obowiązujących Wykonawcę w ramach danego projektu | Instrukcje obowiązują domyślnie w zakresie projektu; współdzielenie z innym projektem konfigurowalne | Edycja instrukcji, przypisanie do projektu |
| Context Memory | 3 | Menu kebab (⋮) — rozszerzenie boczne | Pamięć kontekstowa projektu | Zapisane ustalenia, fakty i decyzje właściwe projektowi | Domyślnie odrębna per projekt; zakres współdzielenia sterowany z okna konfiguracji punktów izolacji (rozdział 8) | Zapis ustalenia, przegląd pamięci, edycja wpisu |

```
 Makieta — Moduł: Workspace       Dostępność: TalkIn, WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Project Dashboard  │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  zadania projektu  │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  status prac       │ ☰
             │                   │                    │  zasoby projektu   │
             │ historia rozmowy  │ dekompozycja       │                    │
             │ pole poleceń      │ zlecenia projektu  │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Agent Manager korzysta z utworzonych agentów | Agent Builder (moduł Agents) |
| Project Library jako odpowiednik ograniczony do projektu | Library Explorer (moduł Library) |

### 6.3. Automations

Moduł budowania i wykonywania procesów automatycznych; konfigurowany ze strony głównej, bez własnego okna w bocznej nawigacji środowisk (Koncepcja platformy, rozdział 11.3). **Dostępność:** konfiguracja ze strony głównej (strefa komponentów własnych); gotowe automatyki trafiają do sesji jako komponenty własne.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Workflow Builder | 1 | Widoczne bez interakcji — kolumna dominująca | Projektowanie procesu automatycznego | Kroki procesu, warunki, kolejność wykonania | Zmiany zapisywane jako definicja automatyki — komponentu własnego | Dodanie kroku, ustalenie warunku, zapis workflow |
| Execution Monitor | 2 | Wskaźnik przebiegu — rozszerzenie boczne | Obserwacja kolejnych uruchomień automatyki | Status każdego przebiegu, błędy wymagające interwencji | Aktualizuje się na żywo w miarę przebiegu wykonania | Przegląd statusu, interwencja przy błędzie, ponowne uruchomienie |
| Queue Manager | 2 | Znacznik kolejki — rozszerzenie boczne | Zarządzanie kolejką zadań automatyki | Lista zadań oczekujących i przetwarzanych | Kolejka przetwarzana zgodnie z regułami silnika kolejek | Dodanie zadania do kolejki, zmiana priorytetu |
| Scheduler | 2 | Przycisk harmonogramu — rozszerzenie boczne | Ustalanie cykliczności uruchomień | Harmonogram: częstotliwość, godziny, zdarzenia wyzwalające | Harmonogram obowiązuje po powiązaniu z workflow | Ustalenie cykliczności, edycja harmonogramu |
| Orchestrator | 3 | Menu kebab (⋮) — rozszerzenie boczne | Definiowanie zależności między krokami procesu | Zależności i kolejność wykonania kroków workflow | Zależności respektowane w toku wykonania procesu | Ustalenie zależności, zmiana kolejności kroków |

```
 Makieta — Moduł: Automations     Konfiguracja ze strony głównej
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Workflow Builder   │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  kroki procesu     │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  warunki           │ ☰
             │                   │                    │  kolejność         │
             │ polecenia wobec   │ nadzór przebiegów  │  wykonania         │
             │ kolejki zadań     │ ponowienia, korekty│                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzyśrodowiskowe (jawne, konfigurowalne; rozdział 13.6 Koncepcji) | Zastosowanie |
|---|---|
| Queue Manager, Orchestrator, Execution Monitor | Operacyjne zaplecze sekcji Kolejki, Orkiestracja oraz Harmonogram i automatyki panelu orkiestracji środowiska MultitaskingAI (rozdział 7) |

### 6.4. Browser

Moduł współdzielonego przeglądania internetu (Koncepcja platformy, rozdział 11.4). **Dostępność:** TalkIn, WorkSpace.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Browser Window | 1 | Widoczne bez interakcji — kolumna dominująca | Współdzielony podgląd strony internetowej | Wyrenderowana strona widoczna jednocześnie użytkownikowi i Wykonawcy | Podgląd wspólny — Wykonawca odpowiada na pytania dotyczące wyświetlanej treści | Nawigacja do strony, przewijanie, zaznaczenie fragmentu |
| Sources Panel | 2 | Znacznik źródeł — rozszerzenie boczne | Lista źródeł wykorzystanych w toku pracy | Adresy i tytuły stron odwiedzonych w sesji | Narasta w miarę przeglądania kolejnych stron | Dodanie źródła, przegląd listy, usunięcie pozycji |
| Notes Panel | 2 | Przycisk notatki — rozszerzenie boczne | Odnotowywanie istotnych fragmentów treści | Notatki użytkownika i Wykonawcy dotyczące przeglądanej strony | Notatka powiązana z konkretnym źródłem z Sources Panel | Dodanie notatki, edycja, powiązanie ze źródłem |

```
 Makieta — Moduł: Browser         Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Browser Window     │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  współdzielony     │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  podgląd strony    │ ☰
             │                   │                    │  (użytkownik +     │
             │ pytania o treść   │ zadania zbierania  │   Wykonawca)       │
             │ wyświetlanej      │ i weryfikacji      │                    │
             │ strony            │ źródeł             │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Sources Panel zasila gromadzenie źródeł badania | Sources Manager (moduł Research) |
| Notes Panel zasila trwałe repozytorium | Library Explorer (moduł Library) |

### 6.5. Research

Moduł realizacji badań, analiz i opracowań (Koncepcja platformy, rozdział 11.5). **Dostępność:** TalkIn, WorkSpace.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Research Workspace | 1 | Widoczne bez interakcji — kolumna dominująca | Główna przestrzeń prowadzenia badania | Struktura badania: zakres, etapy przebiegu badania, powiązane materiały | Punkt wejścia integrujący pozostałe okna modułu | Definiowanie zakresu badania, nawigacja do pozostałych okien |
| Sources Manager | 2 | Znacznik źródeł — rozszerzenie boczne | Gromadzenie źródeł badania | Lista źródeł wraz z metadanymi (typ, pochodzenie, data pozyskania) | Przyjmuje źródła przekazane z modułu Browser | Dodanie źródła, katalogowanie, ocena wiarygodności |
| Findings Panel | 2 | Przycisk ustaleń — rozszerzenie boczne | Odnotowywanie ustaleń cząstkowych | Wnioski i obserwacje powstające w miarę postępu analizy | Narasta w toku badania, powiązane ze źródłami z Sources Manager | Zapis ustalenia, powiązanie ze źródłem, edycja |
| Report Builder | 2 | Przycisk raportu — rozszerzenie boczne | Kompletowanie raportu końcowego | Struktura raportu złożona z ustaleń i wniosków | Buduje dokument końcowy na podstawie zawartości Findings Panel | Kompozycja raportu, redakcja treści |
| Export Panel | 3 | Menu kebab (⋮) — rozszerzenie boczne | Eksport raportu do formatu docelowego | Ustawienia formatu eksportu i miejsca docelowego | Aktywny po skompletowaniu raportu w Report Builder | Eksport dokumentu, wybór formatu wyjściowego |

```
 Makieta — Moduł: Research        Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Research Workspace │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  zakres badania    │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  etapy przebiegu   │ ☰
             │                   │                    │  powiązane         │
             │ polecenia badawcze│ podział badania    │  materiały         │
             │ i wyjaśnienia     │ na zadania, kontro-│                    │
             │ wyniku            │ la kompletności    │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Sources Manager korzysta z zebranych źródeł | Sources Panel (moduł Browser) |
| Sources Manager korzysta z przechowywanych dokumentów | Library Explorer (moduł Library) |
| Raport z Report Builder dalej redagowany | Studio Editor (moduł Studio) |
| Raport prezentowany wielomodelowo | Model Panels (moduł Roundtable) |

### 6.6. Library

Moduł centralnego repozytorium wiedzy i plików (Koncepcja platformy, rozdział 11.6). **Dostępność:** TalkIn, WorkSpace.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Library Explorer | 1 | Widoczne bez interakcji — kolumna dominująca | Przegląd i nawigacja po zgromadzonych materiałach | Struktura plików i zasobów wiedzy | Odbiera artefakty generowane w innych modułach | Nawigacja, wyszukiwanie, otwarcie zasobu |
| File Preview | 2 | Wybór pliku — rozszerzenie boczne | Podgląd zawartości pliku | Wyrenderowana zawartość pliku bez opuszczania modułu | Otwiera się po wybraniu pliku w Library Explorer | Podgląd pliku, przejście do edycji we właściwym module |
| Tags & Collections | 3 | Menu hamburger (☰) — rozszerzenie boczne | Porządkowanie materiałów | Etykiety i kolekcje przypisane zasobom | Etykiety i kolekcje definiowane swobodnie przez użytkownika | Nadanie etykiety, utworzenie kolekcji, przypisanie zasobu |
| Versioning Panel | 3 | Menu kebab (⋮) — rozszerzenie boczne | Śledzenie wersji dokumentu | Kolejne wersje tego samego dokumentu, w tym wygenerowane automatycznie | Wersje narastają wraz z każdą zmianą dokumentu w dowolnym module | Przegląd wersji, przywrócenie wcześniejszej wersji |

```
 Makieta — Moduł: Library         Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Library Explorer   │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  struktura plików  │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  i zasobów wiedzy  │ ☰
             │                   │                    │  odbiór artefaktów │
             │ wyszukiwanie      │ zadania porządko-  │  z innych modułów  │
             │ i polecenia       │ wania i indeksacji │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Library Explorer jako centralne repozytorium | Studio Editor (Studio), Research Workspace (Research), Browser Window (Browser) |
| Odpowiednik ograniczony do projektu | Project Library (moduł Workspace) |

### 6.7. Translate

Moduł wielojęzycznych tłumaczeń w trybie wielozadaniowym (Koncepcja platformy, rozdział 11.7). **Dostępność:** TalkIn.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Source Panel | 1 | Widoczne bez interakcji — kolumna dominująca, część lewa | Umieszczenie tekstu źródłowego | Tekst przeznaczony do tłumaczenia na wiele języków jednocześnie | Zmiana tekstu źródłowego uruchamia aktualizację wszystkich Translation Panels | Wprowadzenie lub wklejenie tekstu źródłowego |
| Translation Panels | 1 | Widoczne bez interakcji — kolumna dominująca, kolumny języków | Równoległe tłumaczenia tekstu źródłowego | Osobna kolumna tłumaczenia dla każdego wybranego języka docelowego | Kolumny aktualizują się jednocześnie po zmianie tekstu źródłowego | Wybór języków docelowych, przegląd i korekta tłumaczenia |
| Glossary Manager | 3 | Menu kebab (⋮) — rozszerzenie boczne | Zapewnienie spójności terminologii | Preferowane odpowiedniki kluczowych pojęć w poszczególnych językach | Stosowany automatycznie przy generowaniu tłumaczeń w Translation Panels | Definiowanie odpowiednika terminu, edycja glosariusza |

```
 Makieta — Moduł: Translate       Dostępność: TalkIn
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window    │ Execution Loop  │ Source  │ Transl. │ Transl.│ ▼
  nawigacja  │ Użytkownik ↔   │ Koordynator ↔   │ Panel   │ Panel   │ Panel  │ ⋮
  modułów    │ Wykonawca      │ Wykonawca       │ tekst   │ język 1 │ język 2│ ☰
             │                │                 │ źródłowy│         │        │
             │ polecenia      │ zadania tłumacz.│         │ aktualizacja     │
             │ redakcyjne     │ kontrola spój-  │         │ jednoczesna po   │
             │                │ ności terminów  │         │ zmianie źródła   │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne; rozdział 6.5 i Załącznik E.4 Koncepcji) | Okno / moduł docelowy |
|---|---|
| Współdzielenie operacji kontekstowych przy pracy dwujęzycznej | Tools Panel (moduł Studio) |

### 6.8. Roundtable

Moduł współpracy wielu modeli nad wspólnym problemem (Koncepcja platformy, rozdział 11.8). **Dostępność:** TalkIn, WorkSpace, CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Model Panels | 1 | Widoczne bez interakcji — kolumna dominująca, kolumny modeli | Równoległe odpowiedzi wielu modeli na to samo zagadnienie | Osobna kolumna odpowiedzi dla każdego uczestniczącego modelu | Kolumny aktualizują się jednocześnie po zadaniu pytania | Wybór uczestniczących modeli, przegląd odpowiedzi |
| Debate Panel | 1 | Widoczne bez interakcji — kolumna dominująca, część dalsza | Rejestrowanie wymiany argumentów między modelami | Zapis kolejnych wypowiedzi modeli odnoszących się do siebie nawzajem | Narasta w toku debaty inicjowanej przez użytkownika lub Moderator Panel | Uruchomienie kolejnej tury debaty, przegląd argumentów |
| Moderator Panel | 2 | Przycisk moderacji — rozszerzenie boczne | Ukierunkowanie przebiegu dyskusji | Narzędzia sterowania: temat, kolejność głosu, zakończenie tury | Aktywny w toku całej sesji Roundtable | Ukierunkowanie dyskusji, zamknięcie tury, wybór kolejnego zagadnienia |
| Consensus Panel | 2 | Znacznik stanowiska — rozszerzenie boczne | Gromadzenie finalnego, uzgodnionego stanowiska | Wypracowane wnioski końcowe wynikające z debaty | Aktualizuje się po zakończeniu tury lub całej debaty | Zapis stanowiska końcowego, eksport wniosków |

```
 Makieta — Moduł: Roundtable      Dostępność: TalkIn, WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window    │ Execution Loop  │ Model │ Model │ Model │ De-│ ▼
  nawigacja  │ Użytkownik ↔   │ Koordynator ↔   │ Panel │ Panel │ Panel │ ba-│ ⋮
  modułów    │ Wykonawca      │ Wykonawca       │  (1)  │  (2)  │  (3)  │ te │ ☰
             │                │                 │       │       │       │ Pa-│
             │ zadanie pytania│ kolejność głosu │ odpowiedzi równo-     │ nel│
             │ i wyjaśnienia  │ tury, kontrola  │ ległe, aktualizacja   │    │
             │                │ jakości wniosku │ jednoczesna           │    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzyśrodowiskowe | Charakter |
|---|---|
| Rozwinięcie idei wielomodelowości w formie zorientowanej na role | Okna robocze środowiska MultitaskingAI (rozdział 7) |

### 6.9. Design

Moduł projektowania i tworzenia zasobów wizualnych (Koncepcja platformy, rozdział 11.9). **Dostępność:** WorkSpace, CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Design Board | 1 | Widoczne bez interakcji — kolumna dominująca | Zestawianie i praca koncepcyjna nad kompozycją wizualną | Zgromadzone i zestawione zasoby graficzne | Odbiera zasoby wygenerowane i zgromadzone w Assets Panel | Zestawienie kompozycji, edycja układu |
| Assets Panel | 2 | Znacznik zasobów — rozszerzenie boczne | Przechowywanie wygenerowanych i zgromadzonych zasobów | Wygenerowane grafiki, ilustracje, elementy brandingowe | Zasoby trafiają tu bezpośrednio po wygenerowaniu | Zapis zasobu, przegląd, usunięcie |
| Prompt Builder | 2 | Przycisk generowania — rozszerzenie boczne | Precyzyjne formułowanie poleceń generujących grafikę | Struktura promptu: styl, kompozycja, parametry generowania | Wynik generowania trafia do Assets Panel | Budowa promptu, uruchomienie generowania |
| Preview Window | 2 | Przełącznik podglądu — rozszerzenie boczne | Podgląd wygenerowanego lub edytowanego zasobu | Wyrenderowany podgląd grafiki | Aktualizuje się po każdej zmianie zaakceptowanej w Design Board | Podgląd zasobu, akceptacja wyniku |

```
 Makieta — Moduł: Design          Dostępność: WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Design Board       │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  zestawianie       │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  kompozycji        │ ☰
             │                   │                    │  praca koncepcyjna │
             │ polecenia         │ zadania generowa-  │  nad układem       │
             │ generujące        │ nia, kontrola      │                    │
             │ grafikę           │ zgodności ze stylem│                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Assets Panel udostępnia zasoby przy redagowaniu dokumentów | Studio Editor (moduł Studio) |
| Assets Panel udostępnia zasoby przy budowie interfejsu | Frontend Workspace (moduł Apps) |

### 6.10. Assistant

Moduł naturalnej komunikacji głosowej; konfigurowany jako komponent własny na stronie głównej, wykorzystywany operacyjnie w środowisku TalkIn (Koncepcja platformy, rozdział 11.10). **Dostępność:** TalkIn (komponent własny konfigurowany na stronie głównej).

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Voice Console | 1 | Widoczne bez interakcji — kolumna dominująca | Wydawanie poleceń głosowych | Interfejs nagrywania i rozpoznawania mowy, synteza odpowiedzi głosowej | Aktywny w trakcie interakcji głosowej z użytkownikiem | Wydanie polecenia głosowego, odsłuch odpowiedzi |
| Actions Monitor | 2 | Wskaźnik realizacji — rozszerzenie boczne | Śledzenie bieżącego statusu realizacji poleceń | Status realizacji zlecenia wieloetapowego | Aktualizuje się na żywo w toku realizacji zlecenia | Przegląd statusu, wstrzymanie lub anulowanie działania |
| Activity Feed | 3 | Menu kebab (⋮) — rozszerzenie boczne | Chronologiczny zapis wykonanych działań | Lista zrealizowanych zleceń wraz z wynikiem | Pozwala odtworzyć przebieg wieloetapowego zlecenia zrealizowanego głosowo | Przegląd historii działań, powrót do wyniku wcześniejszego zlecenia |

```
 Makieta — Moduł: Assistant       Dostępność: TalkIn (komponent własny)
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Voice Console      │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  nagrywanie mowy   │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  rozpoznawanie     │ ☰
             │                   │                    │  synteza głosu     │
             │ transkrypcja      │ dekompozycja zlece-│                    │
             │ poleceń głosowych │ nia wieloetapowego │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie | Charakter |
|---|---|
| Wspólny charakter komunikacji głosowej z funkcją globalną Always On Display | Assistant działa jako moduł osadzony w środowisku TalkIn, a nie jako warstwa obecna ponad całą platformą |

### 6.11. Terminal

Moduł pracy z konsolami i środowiskami wykonawczymi (Koncepcja platformy, rozdział 11.11). **Dostępność:** CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Terminal Tabs | 1 | Widoczne bez interakcji — kolumna dominująca, karty | Prowadzenie równoległych sesji w różnych powłokach | Karty odpowiadające sesjom PowerShell, CMD, Bash i innym | Każda karta działa jako odrębna sesja powłoki | Otwarcie nowej karty, przełączanie między powłokami |
| Output Console | 1 | Widoczne bez interakcji — kolumna dominująca | Gromadzenie wyniku poleceń | Wynik wykonania poleceń ze wszystkich otwartych kart | Aktualizuje się na żywo w miarę wykonywania poleceń | Przegląd wyniku, przewijanie historii wyjścia |
| Process Monitor | 2 | Wskaźnik procesów — rozszerzenie boczne | Obserwacja i kontrola uruchomionych procesów | Lista aktywnych procesów, w tym zainicjowanych przez Wykonawcę | Aktualizuje się na żywo zgodnie z rejestrem procesów rdzenia serwera | Podgląd procesu, zakończenie procesu |

```
 Makieta — Moduł: Terminal        Dostępność: CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Terminal Tabs      │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  PowerShell · CMD  │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  Bash · …          │ ☰
             │                   │                    │ ─────────────────  │
             │ polecenia powłoki │ kolejka poleceń    │ Output Console     │
             │ i wyjaśnienia     │ kontrola wyniku    │  wynik poleceń ze  │
             │ wyniku            │ wykonania          │  wszystkich kart   │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Terminal jako warstwa wykonawcza poleceń | Code Editor (moduł Developer), Diagnostics Center (moduł Diagnostics) |

### 6.12. Developer

Moduł tworzenia i rozwoju kodu (Koncepcja platformy, rozdział 11.12). **Dostępność:** CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Code Editor | 1 | Widoczne bez interakcji — kolumna dominująca | Edycja kodu źródłowego | Treść plików źródłowych repozytorium | Wsparcie Wykonawcy udzielane w kontekście danego repozytorium przez Chat Window | Edycja kodu, generowanie, refaktoryzacja |
| Project Tree | 1 | Widoczne bez interakcji — kolumna dominująca, część lewa | Nawigacja po strukturze projektu | Struktura katalogów i plików repozytorium | Odzwierciedla aktualny stan repozytorium | Nawigacja po plikach, otwarcie pliku w Code Editor |
| Git Panel | 2 | Znacznik repozytorium — rozszerzenie boczne | Zarządzanie zmianami w repozytorium | Status zmian, historia commitów, gałęzie | Odzwierciedla stan repozytorium powiązanego z sesją | Zatwierdzenie zmian, przegląd historii, zmiana gałęzi |
| Build Output | 2 | Wskaźnik budowania — rozszerzenie boczne | Podgląd wyniku kompilacji lub budowania | Log procesu budowania, wynik testów | Aktualizuje się na żywo w toku budowania | Uruchomienie budowania, przegląd wyniku, przegląd błędów |

```
 Makieta — Moduł: Developer       Dostępność: CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window    │ Execution Loop  │ Project │ Code Editor    │ ▼
  nawigacja  │ Użytkownik ↔   │ Koordynator ↔   │ Tree    │  treść plików  │ ⋮
  modułów    │ Wykonawca      │ Wykonawca       │ struktu-│  źródłowych    │ ☰
             │                │                 │ ra kata-│  wsparcie      │
             │ polecenia wobec│ zadania budowa- │ logów   │  Wykonawcy     │
             │ repozytorium   │ nia, testów     │ i plików│                │
             │ sesji          │ i poprawek      │         │                │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Warstwa wykonawcza poleceń | Terminal Tabs (moduł Terminal) |
| Analiza wykrytych błędów | Diagnostics Center (moduł Diagnostics) |
| Komponent budowy produktów | Frontend Workspace, Backend Workspace (moduł Apps) |

### 6.13. Diagnostics

Moduł analizy oraz usuwania problemów technicznych (Koncepcja platformy, rozdział 11.13). **Dostępność:** CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Diagnostics Center | 1 | Widoczne bez interakcji — kolumna dominująca | Agregacja spójnego obrazu stanu systemu | Zestawienie danych z Logs Viewer i Errors Panel | Punkt wejścia integrujący pozostałe okna modułu | Przegląd zagregowanego stanu, uruchomienie analizy |
| Logs Viewer | 2 | Znacznik logów — rozszerzenie boczne | Przegląd logów systemowych | Zapis zdarzeń systemu w porządku chronologicznym | Aktualizuje się na żywo w miarę powstawania nowych wpisów | Przeszukiwanie logów, filtrowanie po zdarzeniu |
| Errors Panel | 2 | Znacznik błędów — rozszerzenie boczne | Przegląd zarejestrowanych błędów | Lista błędów wraz z kontekstem ich wystąpienia | Zasila Diagnostics Center materiałem źródłowym do analizy | Przegląd błędu, oznaczenie jako rozwiązany |
| Recommendations Panel | 2 | Przycisk rekomendacji — rozszerzenie boczne | Prezentacja sugerowanych kroków naprawczych | Rekomendacje wypracowane na podstawie zebranych danych | Aktualizuje się po każdej analizie przeprowadzonej w Diagnostics Center | Przegląd rekomendacji, zastosowanie poprawki |

```
 Makieta — Moduł: Diagnostics     Dostępność: CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Diagnostics Center │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  zagregowany obraz │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  stanu systemu     │ ☰
             │                   │                    │                    │
             │ pytania o przyczy-│ zadania analizy    │                    │
             │ nę i przebieg     │ i weryfikacji      │                    │
             │ błędu             │ poprawek           │                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe (jawne, konfigurowalne) | Okno / moduł docelowy |
|---|---|
| Wdrażanie poprawek | Code Editor (moduł Developer) |
| Odtwarzanie i weryfikacja objawów błędu | Terminal Tabs (moduł Terminal) |

### 6.14. Apps

Moduł budowy kompletnych produktów cyfrowych (Koncepcja platformy, rozdział 11.14). **Dostępność:** WorkSpace, CodeStudio.

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Product Builder | 1 | Widoczne bez interakcji — kolumna dominująca | Nadrzędny widok procesu budowy produktu | Kolejne czynności budowy produktu od architektury po wdrożenie | Punkt wejścia integrujący pozostałe okna modułu | Nawigacja między czynnościami budowy produktu |
| Architecture Designer | 2 | Przycisk architektury — rozszerzenie boczne | Projektowanie architektury rozwiązania | Struktura komponentów rozwiązania i ich zależności | Punkt wyjścia procesu — poprzedza pracę w Frontend Workspace i Backend Workspace | Definiowanie komponentów, ustalenie zależności |
| Frontend Workspace | 2 | Znacznik warstwy interfejsu — rozszerzenie boczne | Praca nad warstwą frontendową produktu | Kod i zasoby interfejsu użytkownika | Korzysta z zasobów wizualnych z Assets Panel modułu Design | Edycja interfejsu, podgląd wyniku |
| Backend Workspace | 2 | Znacznik warstwy serwerowej — rozszerzenie boczne | Praca nad warstwą backendową produktu | Kod i konfiguracja usług serwerowych | Prowadzona równolegle z Frontend Workspace | Edycja logiki serwerowej, konfiguracja usług |
| Deployment Panel | 3 | Menu kebab (⋮) — rozszerzenie boczne | Zarządzanie udostępnieniem produktu | Ustawienia wdrożenia, status publikacji | Aktywny po zakończeniu prac w Frontend Workspace i Backend Workspace | Uruchomienie wdrożenia, przegląd statusu publikacji |

```
 Makieta — Moduł: Apps            Dostępność: WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Product Builder    │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  przebieg budowy   │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  produktu:         │ ☰
             │                   │                    │  architektura →    │
             │ polecenia wobec   │ podział pracy na   │  interfejs →       │
             │ produktu          │ tory frontend      │  usługi →          │
             │                   │ i backend          │  wdrożenie         │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie międzymodułowe / środowiskowe (jawne, konfigurowalne) | Okno / moduł / środowisko docelowe |
|---|---|
| Integracja funkcjonalności w procesie budowy produktu | Code Editor (Developer), Terminal Tabs (Terminal) |
| Podział na Executor 1 i Executor 2 przy pracy nad backendem i frontendem | Okna robocze środowiska MultitaskingAI (rozdział 7) |

### 6.15. Agents

Moduł tworzenia i zarządzania własnymi agentami; konfigurowany jako komponent własny na stronie głównej, dostępny także jako okno modułowe w środowiskach wskazanych w macierzy dostępności modułów (Koncepcja platformy, rozdział 11.15). **Dostępność:** TalkIn, WorkSpace, CodeStudio (komponent własny konfigurowany na stronie głównej, dodatkowo okno modułowe).

| Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|
| Agent Builder | 1 | Widoczne bez interakcji — kolumna dominująca | Tworzenie tożsamości agenta | Nazwa, opis, instrukcje systemowe agenta | Punkt wejścia integrujący pozostałe okna modułu | Utworzenie agenta, edycja tożsamości |
| Model Configuration | 2 | Element zwinięty `Model ▼` — rozszerzenie boczne | Wybór modelu bazowego agenta | Dostępne modele i parametry ich działania | Wybór modelu obowiązuje dla danego agenta | Wybór modelu bazowego, ustalenie parametrów |
| Skills Manager | 2 | Znacznik umiejętności — rozszerzenie boczne | Dobór umiejętności agenta | Lista dostępnych skilli możliwych do przypisania agentowi | Przypisane skille dostępne agentowi w pracy operacyjnej | Dodanie skilla, usunięcie skilla |
| Connectors Manager | 3 | Menu hamburger (☰) — rozszerzenie boczne | Podłączenie integracji zewnętrznych | Lista dostępnych konektorów i pluginów | Podłączone konektory dostępne agentowi w toku pracy | Podłączenie konektora, konfiguracja integracji |
| Permissions Center | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny albo konfiguracja roli | Ustalenie zakresu uprawnień agenta | Lista uprawnień możliwych do przyznania lub odebrania | Zakres uprawnień obowiązuje przed udostępnieniem agenta do pracy operacyjnej; użytkownik podstawowy nie widzi tego okna | Przyznanie uprawnienia, odebranie uprawnienia |

```
 Makieta — Moduł: Agents          Dostępność: TalkIn, WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window       │ Execution Loop     │ Agent Builder      │ ▼
  nawigacja  │ Użytkownik ↔      │ Koordynator ↔      │  tożsamość agenta: │ ⋮
  modułów    │ Wykonawca         │ Wykonawca          │  nazwa · opis      │ ☰
             │                   │                    │  instrukcje        │
             │ polecenia budowy  │ zadania weryfikacji│  systemowe         │
             │ agenta            │ konfiguracji agenta│                    │
 ═══════════════════════════════════════════════════════════════════════════
```

| Powiązanie (jawne, konfigurowalne) | Okno / środowisko docelowe |
|---|---|
| Przypisanie agentów do projektów | Agent Manager (moduł Workspace) |
| Pełnienie ról wykonawczych | Okna robocze środowiska MultitaskingAI (rozdział 7) |

---

## 7. Okna robocze środowiska MultitaskingAI

### 7.1. Charakter okien środowiska MultitaskingAI

Środowisko MultitaskingAI nie organizuje pracy modułowej pojedynczego użytkownika, lecz zespół modeli i agentów realizujących wspólny proces (Koncepcja platformy, rozdział 13.1). Jego boczna nawigacja nie zawiera listy modułów, lecz panel orkiestracji — sześć sekcji odpowiadających własnym warstwom sterowania środowiska (rozdział 13.8 Koncepcji platformy).

| Element | Charakterystyka |
|---|---|
| Jednostka organizacji | Zespół modeli i agentów realizujących wspólny proces, a nie pojedynczy użytkownik z modułami |
| Boczna nawigacja | Panel orkiestracji (sześć sekcji), a nie lista modułów |
| Okna komunikacji operacyjnej | Chat Window w lewej kolumnie stałej; Execution Loop Window w kolumnie sąsiadującej, którego rozwinięciem właściwym temu środowisku jest Coordinator Chat |
| Sekcje z własnymi oknami roboczymi | Cztery sekcje operują oknami przypisanymi bezpośrednio rolom zespołu |
| Sekcje korzystające z okien modułu Automations | Pozostałe dwie sekcje, po skonfigurowaniu powiązań, korzystają z okien skatalogowanych w rozdziale 6.3 |

### 7.2. Okna czterech ról

Sekcja Role panelu orkiestracji udostępnia cztery okna robocze, po jednym na każdą rolę zespołu (Koncepcja platformy, rozdział 13.3 i 13.8). Nazwy okien odpowiadają charakterowi odpowiedzialności każdej roli: role wykonawcze (Executor 1, Executor 2) korzystają z okna typu Executor Chat, rola koordynująca korzysta z okna typu Coordinator Chat — instancji Execution Loop Window właściwej temu środowisku — a rola kontrolna (Executor 3 / Validator) korzysta z okna typu Results Analyzer, którego zawartość odpowiada jej funkcji kontroli jakości pracy pozostałych modeli.

| Rola | Okno | Warstwa | Sposób wywołania | Cel | Zawartość | Zachowanie | Powiązane operacje |
|---|---|---|---|---|---|---|---|
| Coordinator (koordynator) | Coordinator Chat — instancja Execution Loop Window | 1 | Widoczne bez interakcji — kolumna sąsiadująca z Chat Window | Prowadzenie pętli wykonawczej: planowanie, koordynacja, nadzór, orkiestracja, kontrola realizacji | Plan etapów przebiegu, przypisania zadań, budowane prompty dla Wykonawców, stan kolejki, wskaźniki przebiegu pętli | Nie tworzy produktu końcowego; steruje: start, stop, pauza, wznowienie, przekazanie, powtórzenie, walidacja | Planowanie przebiegu, podział pracy, budowa promptu, sterowanie kolejką, korekta zlecenia |
| Executor 1 (główny wykonawca) | Executor Chat | 1 | Widoczne bez interakcji — kolumna dominująca, część lewa | Realizacja pracy zleconej przez Koordynatora lub użytkownika | Historia poleceń i wyników pracy; dostęp do pamięci projektu; panel Subagent Network | Wykonuje zadania pobierane z kolejki; nie zarządza procesem | Wykonanie zadania, uruchomienie podagentów (do 15), zwrot wyniku |
| Executor 2 (równoległy wykonawca) | Executor Chat | 1 | Widoczne bez interakcji — kolumna dominująca, część prawa | Równoległa realizacja drugiego toru pracy | Jak Executor 1, w trybie współpracy ustalonym z Executorem 1 | Tryb współpracy konfigurowalny: praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna | Wykonanie zadania równoległego, wymiana wyników z Executorem 1 |
| Executor 3 / Validator (czwarty model) | Results Analyzer | 2 | Wskaźnik kontroli jakości — rozszerzenie boczne | Kontrola jakości, zgodności lub bezpieczeństwa pracy pozostałych ról | Wyniki pracy Executora 1 i Executora 2 przekazane do oceny; ustalenia kontrolne | Wcielenie roli konfigurowalne: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator | Ocena wyniku, zgłoszenie niezgodności, rozstrzygnięcie konfliktu między modelami |

Executor 1 i Executor 2 korzystają z okna tego samego typu (Executor Chat), otwartego w dwóch niezależnych instancjach — po jednej na każdą rolę wykonawczą, zgodnie z zasadą, że MultitaskingAI udostępnia cztery okna robocze odpowiadające czterem rolom (Koncepcja platformy, rozdział 13.8). Rozmieszczenie okien czterech ról wraz z panelem orkiestracji przedstawia makieta w stanie spoczynku.

```
 Makieta — Środowisko: MultitaskingAI (obszar ról)
 ═══════════════════════════════════════════════════════════════════════════
  Panel      │ Chat Window    │ Coordinator Chat │ Executor Chat │ Executor │ ▼
  orkiestracji│ Użytkownik ↔  │ (Execution Loop  │ (Executor 1)  │ Chat     │ ⋮
  (boczna    │ Wykonawca      │  Window)         │               │ (Exec. 2)│ ☰
  nawigacja):│                │ Koordynator ↔    │ wykonanie     │ drugi,   │
   Zespoły   │ zlecenia       │ Wykonawca        │ zadań         │ równoległ│
   Role      │ użytkownika    │                  │ z kolejki     │ y tor    │
   Kolejki   │ zatwierdzenia  │ plan przebiegu   │               │ pracy    │
   Orkiestra-│ przerwania     │ podział pracy    │ Subagent      │ Subagent │
    cja      │ wyjaśnienia    │ budowa promptów  │ Network       │ Network  │
   Harmono-  │ wyniku         │ stan kolejki     │ (do 15 pod-   │ (do 15   │
    gram i   │                │ kontrola jakości │  agentów)     │  pod-    │
    automat. │                │ sterowanie pętlą │               │  agentów)│
   Monitor   │                │                  │               │          │
    procesu  │                │                  │               │          │
 ═══════════════════════════════════════════════════════════════════════════
   kolumna     lewa kolumna     kolumna            kolumna dominująca   kolumny
   stała       stała            sąsiadująca                             boczne
```

Results Analyzer należy do warstwy 2 i otwiera się jako kolumna boczna po prawej stronie obszaru roboczego; w stanie spoczynku reprezentuje go wskaźnik kontroli jakości.

### 7.3. Pozostałe sekcje panelu orkiestracji

Poniższe sekcje panelu orkiestracji nie posiadają odrębnych okien roboczych poza oknami ról opisanymi w rozdziale 7.2 — realizują swoje funkcje przez zapisane konfiguracje oraz, po skonfigurowaniu integracji (Koncepcja platformy, rozdział 13.6), przez okna modułu Automations skatalogowane w rozdziale 6.3 niniejszego dokumentu.

| Sekcja | Warstwa | Sposób wywołania | Zawartość | Wykorzystywane okna |
|---|---|---|---|---|
| Zespoły | 2 | Pozycja panelu orkiestracji | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek; zapis, wczytanie, duplikowanie | Konfiguracje ról z okien opisanych w rozdziale 7.2 |
| Kolejki | 2 | Pozycja panelu orkiestracji | Definicje kolejek globalnych, lokalnych, dla modeli, agentów i projektów oraz ich akcje (silnik kolejek, rozdział 13.4 Koncepcji platformy) | Queue Manager (moduł Automations, rozdział 6.3), po skonfigurowaniu powiązania |
| Orkiestracja | 3 | Pozycja panelu orkiestracji, menu kebab (⋮) | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami (rozdział 13.5 Koncepcji platformy) | Orchestrator (moduł Automations, rozdział 6.3), po skonfigurowaniu powiązania |
| Harmonogram i automatyki | 2 | Pozycja panelu orkiestracji | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations (rozdział 13.6 Koncepcji platformy) | Scheduler, Execution Monitor (moduł Automations, rozdział 6.3) |
| Monitor procesu | 2 | Pozycja panelu orkiestracji, wskaźnik stanu wykonania | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display (rozdział 13.2 Koncepcji platformy) | Results Analyzer oraz Execution Monitor, uzupełnione widokiem Always On Display |

Zgodnie z zasadą pełnej konfigurowalności kolejność i widoczność sekcji panelu orkiestracji podlegają konfiguracji; zestaw przedstawiony powyżej jest zestawem domyślnym (Koncepcja platformy, rozdział 13.8).

---

## 8. Okno konfiguracji punktów izolacji

### 8.1. Miejsce okna w katalogu

Okno konfiguracji punktów izolacji nie jest oknem modułowym przypisanym do jednego z piętnastu modułów — jest oknem o zasięgu globalnym, otwieranym z okna konfiguracji (Architektura, rozdział 13) i powiązanym z pracą we wszystkich modułach i oknach skatalogowanych w rozdziałach 5–7, ponieważ to ono ustala, w jakim zakresie ich historia, pamięć, kontekst oraz proces techniczny sesji są odrębne lub współdzielone (Koncepcja platformy, rozdział 6.4–6.6).

| Cecha | Wartość |
|---|---|
| Kategoria okna | Okno globalne (nieprzypisane do jednego modułu) |
| Warstwa | 4 — funkcje eksperckie |
| Sposób wywołania | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, skrót klawiszowy albo okno konfiguracji aplikacji (Architektura, rozdział 13) |
| Rozmieszczenie | Kolumna boczna otwierana jako rozszerzenie boczne po prawej stronie obszaru roboczego |
| Zasięg oddziaływania | Wszystkie moduły i okna z rozdziałów 5–7 |
| Podstawa | Koncepcja platformy, rozdział 6.4–6.6 |

### 8.2. Cel

Jedno, spójne miejsce konfiguracji dwóch rodzajów izolacji.

| Rodzaj izolacji | Zakres |
|---|---|
| Izolacja kontekstu | Historia, pamięć i kontekst współdzielony między kartami sesji, modułami, środowiskami i projektami |
| Izolacja techniczna procesu sesji | Osiem zakresów zdefiniowanych w rozdziale 9 dokumentu Architektury: katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania |

### 8.3. Zawartość

Okno dzieli się na trzy kolumny sąsiadujące poziomo.

| Kolumna | Warstwa | Zawartość | Zachowanie |
|---|---|---|---|
| Selektor zasięgu (lewa) | 4 | Siedem poziomów zasięgu: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola | Wybór poziomu determinuje, do którego zakresu odnoszą się ustawienia w kolumnie macierzy izolacji |
| Macierz izolacji (środkowa, dominująca) | 4 | Przełączniki izolacji kontekstu (historia, pamięć, kontekst — każdy: współdzielone / odrębne) oraz izolacji technicznej (osiem zakresów — każdy: włączony / wyłączony) | Każdy przełącznik opatrzony objaśnieniem kontekstowym `[?]` |
| Profil i podgląd (prawa) | 4 | Zapis, wczytanie i przypisanie profilu izolacji; wybór warstwy konfiguracji (domyślna albo sesji); podgląd polityki efektywnej | Podgląd polityki efektywnej uwzględnia dziedziczenie między poziomami zasięgu |

Układ trzech kolumn przedstawia makieta.

```
 Makieta — Okno konfiguracji punktów izolacji (warstwa 4)
 ┌────────────────┬──────────────────────────────┬─────────────────────┐
 │ SELEKTOR       │ MACIERZ IZOLACJI             │ PROFIL I PODGLĄD    │
 │ ZASIĘGU        │                              │                     │
 │                │ Izolacja kontekstu:          │ [ Zapisz profil ]   │
 │ • globalny     │   historia   [współdz.|odr.] │ [ Wczytaj profil ]  │
 │ • środowisko   │   pamięć     [współdz.|odr.] │ [ Przypisz do… ]    │
 │ • moduł        │   kontekst   [współdz.|odr.] │                     │
 │ • para modułów │                              │ Warstwa:            │
 │ • projekt      │ Izolacja techniczna (8):     │  ( ) domyślna       │
 │ • karta sesji  │   katalog roboczy   [wł|wył] │  ( ) sesji          │
 │ • rola         │   środowisko proc.  [wł|wył] │                     │
 │                │   dane/konfig.modelu[wł|wył] │ Polityka efektywna: │
 │                │   dostęp sieciowy   [wł|wył] │  (podgląd reguł     │
 │                │   odczyt/zapis plik.[wł|wył] │   dziedziczonych)   │
 │                │   konto i token     [wł|wył] │                     │
 │                │   model procesu     [wł|wył] │ [?] objaśnienia     │
 │                │   serwer wykonania  [wł|wył] │     kontekstowe     │
 └────────────────┴──────────────────────────────┴─────────────────────┘
   kolumna lewa      kolumna dominująca             kolumna prawa
```

Zestaw ustawień izolacji zapisuje się jako nazwany profil izolacji. Poniższy szablon porządkuje jego pola; wszystkie podlegają zasadzie „brak ustawienia oznacza wartość domyślną” (dziedziczenie z poziomu szerszego, brak aktywnej izolacji technicznej).

```
profil_izolacji:
  nazwa:          <nazwa własna profilu>
  poziom_zasięgu: <globalny | środowisko | moduł | para modułów |
                   projekt | karta sesji | rola>
  izolacja_kontekstu:
    historia:  <współdzielona | odrębna>
    pamięć:    <współdzielona | odrębna>
    kontekst:  <współdzielony  | odrębny>
  izolacja_techniczna:                       # osiem zakresów; każdy: włączony | wyłączony
    katalog_roboczy_sesji:              wyłączony
    środowisko_procesu:                 wyłączony
    katalog_danych_i_konfiguracji_modelu: wyłączony
    dostęp_sieciowy:                    wyłączony
    odczyt_i_zapis_plików:              wyłączony
    konto_i_token_per_sesja:            wyłączony
    model_procesu:                      wyłączony
    serwer_wykonania:                   wyłączony
  warstwa:      <domyślna | sesji>
  przypisanie:  <sesja | rola>
```

### 8.4. Zachowanie

Reguły z różnych poziomów zasięgu rozstrzyga zasada pierwszeństwa zasięgu najbardziej szczegółowego: ustawienie z poziomu węższego ma pierwszeństwo przed ustawieniem z poziomu szerszego, a brak ustawienia na danym poziomie oznacza dziedziczenie wartości z poziomu bezpośrednio szerszego, aż do poziomu globalnego (Koncepcja platformy, rozdział 6.5).

```
 Pierwszeństwo poziomów zasięgu (od najniższego do najwyższego)

  globalny  →  środowisko  →  moduł  →  para modułów  →
  →  projekt  →  karta sesji  →  rola
  (warstwa bazowa)                        (pierwszeństwo najwyższe)

  brak ustawienia na poziomie  =  dziedziczenie z poziomu szerszego
```

Konfiguracja izolacji ma strukturę warstwową.

| Warstwa konfiguracji | Zakres obowiązywania | Relacja do pozostałych warstw |
|---|---|---|
| Warstwa domyślna | Obowiązuje przy każdej nowej sesji | Warstwa bazowa (Koncepcja platformy, rozdział 6.6) |
| Warstwa sesji | Zmiany dokonane dla sesji bieżącej | Nakłada się na warstwę domyślną bez jej modyfikacji (Architektura, rozdział 13) |

Zgodnie z zasadą braku twardych blokad okno nigdy nie wymusza izolacji — udostępnia ją jako ustawienie konfiguracyjne. Stanem wyjściowym jest odrębna historia, pamięć i kontekst dla każdej nowej karty sesji, przy jednoczesnym braku aktywnej izolacji technicznej — żaden z ośmiu zakresów nie jest domyślnie włączony (Koncepcja platformy, rozdział 6.6).

### 8.5. Powiązane operacje

| Operacja | Opis |
|---|---|
| Wybór poziomu zasięgu | Ustawienie poziomu w selektorze zasięgu |
| Ustawienie przełącznika izolacji | Zmiana izolacji kontekstu lub izolacji technicznej w macierzy |
| Zapis, wczytanie i przypisanie profilu izolacji | Przypisanie profilu do sesji lub roli środowiska MultitaskingAI (rozdział 7.2) |
| Przełączenie warstwy konfiguracji | Zmiana między warstwą domyślną a warstwą sesji |
| Przegląd polityki efektywnej | Weryfikacja wynikowego zestawu reguł przed zatwierdzeniem zmiany |

---

## 9. Macierz zbiorcza okien operacyjnych per moduł

Poniższa macierz zbiera w jednym miejscu okna operacyjne wszystkich piętnastu modułów oraz okna robocze środowiska MultitaskingAI, zgodnie z Załącznikiem A Koncepcji platformy, rozszerzonym o rozdział 11 tego samego dokumentu. Chat Window i Execution Loop Window należą do warstwy 1 i występują w każdym module — kolumna „Liczba okien” obejmuje oba.

| Moduł / środowisko | Okna warstwy 1 (widoczne bez interakcji) | Okna warstw 2–4 | Sposób wywołania okien warstw 2–4 | Liczba okien |
|---|---|---|---|---|
| Studio (6.1) | Chat Window, Execution Loop Window, Studio Editor | Diff/Grep Panel (2), Preview Window (2), Tools Panel (3), Session Repository (3) | Przycisk, przełącznik, menu `Operacje ▼`, menu kebab (⋮) | 7 |
| Workspace (6.2) | Chat Window, Execution Loop Window, Project Dashboard | Project Library (2), Agent Manager (2), Instructions Panel (3), Context Memory (3) | Znacznik, element zwinięty `Agent ▼`, menu kebab (⋮) | 7 |
| Automations (6.3) | Chat Window, Execution Loop Window, Workflow Builder | Execution Monitor (2), Queue Manager (2), Scheduler (2), Orchestrator (3) | Wskaźnik przebiegu, znacznik kolejki, przycisk, menu kebab (⋮) | 7 |
| Browser (6.4) | Chat Window, Execution Loop Window, Browser Window | Sources Panel (2), Notes Panel (2) | Znacznik źródeł, przycisk notatki | 5 |
| Research (6.5) | Chat Window, Execution Loop Window, Research Workspace | Sources Manager (2), Findings Panel (2), Report Builder (2), Export Panel (3) | Znacznik, przycisk, menu kebab (⋮) | 7 |
| Library (6.6) | Chat Window, Execution Loop Window, Library Explorer | File Preview (2), Tags & Collections (3), Versioning Panel (3) | Wybór pliku, menu hamburger (☰), menu kebab (⋮) | 6 |
| Translate (6.7) | Chat Window, Execution Loop Window, Source Panel, Translation Panels | Glossary Manager (3) | Menu kebab (⋮) | 5 |
| Roundtable (6.8) | Chat Window, Execution Loop Window, Model Panels, Debate Panel | Moderator Panel (2), Consensus Panel (2) | Przycisk moderacji, znacznik stanowiska | 6 |
| Design (6.9) | Chat Window, Execution Loop Window, Design Board | Assets Panel (2), Prompt Builder (2), Preview Window (2) | Znacznik zasobów, przycisk generowania, przełącznik podglądu | 6 |
| Assistant (6.10) | Chat Window, Execution Loop Window, Voice Console | Actions Monitor (2), Activity Feed (3) | Wskaźnik realizacji, menu kebab (⋮) | 5 |
| Terminal (6.11) | Chat Window, Execution Loop Window, Terminal Tabs, Output Console | Process Monitor (2) | Wskaźnik procesów | 5 |
| Developer (6.12) | Chat Window, Execution Loop Window, Project Tree, Code Editor | Git Panel (2), Build Output (2) | Znacznik repozytorium, wskaźnik budowania | 6 |
| Diagnostics (6.13) | Chat Window, Execution Loop Window, Diagnostics Center | Logs Viewer (2), Errors Panel (2), Recommendations Panel (2) | Znacznik logów, znacznik błędów, przycisk rekomendacji | 6 |
| Apps (6.14) | Chat Window, Execution Loop Window, Product Builder | Architecture Designer (2), Frontend Workspace (2), Backend Workspace (2), Deployment Panel (3) | Przycisk, znacznik warstwy, menu kebab (⋮) | 7 |
| Agents (6.15) | Chat Window, Execution Loop Window, Agent Builder | Model Configuration (2), Skills Manager (2), Connectors Manager (3), Permissions Center (4) | Element zwinięty `Model ▼`, znacznik, menu hamburger (☰), polecenie języka naturalnego | 7 |
| MultitaskingAI — role (7.2) | Chat Window, Coordinator Chat (instancja Execution Loop Window), Executor Chat (Executor 1), Executor Chat (Executor 2) | Results Analyzer (2) | Wskaźnik kontroli jakości | 5 |
| Okno globalne (8) | — | Okno konfiguracji punktów izolacji (4) | Polecenie języka naturalnego, wyszukiwarka funkcji, skrót klawiszowy, okno konfiguracji aplikacji | 1 |

Poniższa tabela odnotowuje okna występujące w więcej niż jednym miejscu katalogu oraz okna o nazwach zbliżonych brzmieniowo, lecz odrębnych — zgodnie z uwagą przyjętą w Załączniku A Koncepcji platformy.

| Obserwacja | Szczegół |
|---|---|
| Okna wspólne wszystkim modułom | Chat Window i Execution Loop Window (występują w każdym z piętnastu modułów, warstwa 1, rozdział 5) |
| Okno wspólne dwóm modułom | Preview Window (Studio oraz Design — w obu funkcja podglądu wyniku pracy, warstwa 2) |
| Trzy różne okna o nazwach zbliżonych brzmieniowo | Process Monitor (Terminal), Execution Monitor (Automations), Actions Monitor (Assistant) |
| Dwa różne okna repozytorium | Project Library (moduł Workspace) oraz Library Explorer (moduł Library) |
| Okno o dwóch nazwach w zależności od środowiska | Execution Loop Window w modułach; Coordinator Chat jako jego instancja w środowisku MultitaskingAI (rozdział 7.2) |

---

## 10. Słowniczek i zasady nazewnictwa okien

| Termin | Definicja | Odsyłacz |
|---|---|---|
| Okno operacyjne | Pojedynczy element zestawu okien danego modułu, stanowiący właściwą przestrzeń wykonywania zadań przypisaną do konkretnej funkcjonalności modułu; nie należy mylić z oknem operacyjnym aplikacji (Architektura, rozdział 13) | Rozdział 2.2 |
| Zestaw okien modułu | Wszystkie okna operacyjne właściwe danemu modułowi, otwierane łącznie przy jego wybraniu i zamykane łącznie przy przełączeniu na inny moduł | Koncepcja platformy, rozdział 2.2 |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy użytkownika i podstawowy mechanizm sterowania procesami platformy; lewa kolumna, stała, pełna wysokość obszaru roboczego, warstwa 1 | Rozdział 5.1 |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca: koordynacja zadań, nadzór, orkiestracja, kontrola realizacji; kolumna sąsiadująca z Chat Window, warstwa 1 | Rozdział 5.2 |
| Użytkownik | Rola zlecająca i zatwierdzająca | Rozdział 5.3 |
| Koordynator | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela i nadzoruje zadania | Rozdział 5.3 |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania | Rozdział 5.3 |
| Układ pionowy (podział lewa–prawa) | Jedyny obowiązujący sposób rozmieszczenia okien: kolumny sąsiadujące poziomo, o regulowanej szerokości | Rozdział 6.0 |
| Okno pomocnicze | Okno otwierane jako rozszerzenie boczne w kolumnie bocznej, po prawej stronie obszaru roboczego; po zamknięciu znika z przestrzeni roboczej | Rozdział 6.0 |
| Warstwa widoczności | Jedna z czterech warstw (1–4), do której należy każde okno, panel i pojedyncza funkcja; wyznacza, czy element jest widoczny bez interakcji, na żądanie, w rozwinięciu kontekstowym, czy jako funkcja ekspercka | Rozdział 3a |
| Sposób wywołania | Konkretny mechanizm ujawnienia elementu warstw 2–4: znacznik kontekstowy, przycisk, przełącznik, menu kebab (⋮), menu hamburger (☰), polecenie języka naturalnego, skrót klawiszowy | Rozdział 3a |
| Okno wspólne | Okno operacyjne występujące w więcej niż jednym module w niezmienionej roli; oknami wspólnymi wszystkim piętnastu modułom są Chat Window i Execution Loop Window, a Preview Window jest oknem wspólnym dwóm modułom (Studio, Design) | Rozdział 5 |
| Okno robocze roli | Okno operacyjne środowiska MultitaskingAI przypisane jednej z czterech ról zespołu (Executor 1, Executor 2, Coordinator, Executor 3 / Validator), a nie modułowi | Rozdział 7.2 |
| Okno globalne | Okno operacyjne dostępne niezależnie od aktywnego modułu, otwierane z okna konfiguracji aplikacji lub poleceniem języka naturalnego; jedynym oknem tej kategorii jest okno konfiguracji punktów izolacji | Rozdział 8 |
| Powiązane operacje | Konkretne czynności wykonywane przez Użytkownika, Koordynatora lub Wykonawcę w danym oknie operacyjnym lub między tym oknem a innym oknem tego samego bądź innego modułu | Kolumna „Powiązane operacje” tabel rozdziałów 5–8 |
| Zachowanie okna | Sposób, w jaki okno operacyjne reaguje na zdarzenia: moment otwarcia, warunki aktualizacji zawartości, zakres domyślnej izolacji kontekstu oraz relacja do innych okien tego samego modułu | Kolumna „Zachowanie” tabel rozdziałów 5–8 |

---

*Koniec dokumentu. Danaco Console — Specyfikacja okien operacyjnych, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
