# Danaco Console — Specyfikacja modułów

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
| **Tytuł** | Specyfikacja modułów |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Ustala operacyjną specyfikację piętnastu modułów platformy — cel, okna operacyjne z warstwami widoczności, funkcjonalności, przebieg pracy, dane wykorzystywane przez AI oraz powiązania konfigurowalne z innymi modułami — jako rozwinięcie wykonawcze rozdziału 11 dokumentu Koncepcja platformy |
| **Zakres** | piętnaście modułów platformy (rozdz. 4), macierz dostępności modułów w środowiskach (rozdz. 3), warstwy widoczności funkcji modułu (rozdz. 3a), zbiorcze zestawienie powiązań międzymodułowych (rozdz. 5) |
| **Poza zakresem** | środowisko MultitaskingAI jako całość — warstwa centralna, role, silnik kolejek, panel orkiestracji — [Koncepcja platformy](../architektura/koncepcja-platformy.md) rozdz. 13; makiety okien i katalogi elementów interfejsu — [Specyfikacja okien operacyjnych](specyfikacja-okien-operacyjnych.md); model agentów — [Specyfikacja agentów](specyfikacja-agentow.md) |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Architektura techniczna](../architektura/architektura.md) · [Specyfikacja okien operacyjnych](specyfikacja-okien-operacyjnych.md) · [Specyfikacja agentów](specyfikacja-agentow.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary piętnastu modułów, rozdz. 4) · [Koncepcja platformy](../architektura/koncepcja-platformy.md) rozdz. 9, 11, 14 · `design/05-okna/moduly/` |
| **Zasada nadrzędna** | Żadne powiązanie międzymodułowe nie jest domyślnie aktywne — każda zależność jest jawną, konfigurowalną decyzją użytkownika podejmowaną w oknie konfiguracji |

Dokument stanowi operacyjną specyfikację piętnastu modułów platformy Danaco Console. Rozwija w formie wykonawczej ustalenia rozdziału 11 dokumentu Koncepcja platformy — dla każdego modułu określa cel, okna operacyjne wraz z warstwami widoczności, funkcjonalności, przebieg pracy, dane wykorzystywane przez AI oraz powiązania konfigurowalne z innymi modułami. Każdy moduł udostępnia dwa okna komunikacji operacyjnej jako elementy pierwszoplanowe: Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca). Układ interfejsu jest pionowy — okna rozmieszczone są w kolumnach sąsiadujących poziomo.

---

## Spis treści

1. [Zakres i zasady redakcji specyfikacji](#1-zakres-i-zasady-redakcji-specyfikacji)
2. [Moduł a komponent własny](#2-moduł-a-komponent-własny)
3. [Macierz dostępności modułów](#3-macierz-dostępności-modułów)
4. [3a. Warstwy widoczności funkcji modułu](#3a-warstwy-widoczności-funkcji-modułu)
5. [3b. Macierz uprawnień — punkty izolacji i współdzielenie kontekstu](#3b-macierz-uprawnień--punkty-izolacji-i-współdzielenie-kontekstu)
6. [Specyfikacja modułów](#4-specyfikacja-modułów)
   - [4.1 Studio](#41-studio)
   - [4.2 Workspace](#42-workspace)
   - [4.3 Automations](#43-automations)
   - [4.4 Browser](#44-browser)
   - [4.5 Research](#45-research)
   - [4.6 Library](#46-library)
   - [4.7 Translate](#47-translate)
   - [4.8 Roundtable](#48-roundtable)
   - [4.9 Design](#49-design)
   - [4.10 Assistant](#410-assistant)
   - [4.11 Terminal](#411-terminal)
   - [4.12 Developer](#412-developer)
   - [4.13 Diagnostics](#413-diagnostics)
   - [4.14 Apps](#414-apps)
   - [4.15 Agents](#415-agents)
7. [Zbiorcze zestawienie powiązań międzymodułowych](#5-zbiorcze-zestawienie-powiązań-międzymodułowych)
8. [Zgodność z zasadami nadrzędnymi platformy](#6-zgodność-z-zasadami-nadrzędnymi-platformy)
9. [Kryteria odbioru](#7-kryteria-odbioru)

---

## 1. Zakres i zasady redakcji specyfikacji

Dokument doprecyzowuje operacyjnie piętnaście modułów platformy opisanych w rozdziale 11 Koncepcji platformy. Nie obejmuje środowiska MultitaskingAI jako całości — jest ono jednym z czterech środowisk platformy (rozdz. 9.4 Koncepcji), a nie modułem. Poniższa tabela rozgranicza zakres.

| Element | Status w specyfikacji | Miejsce pełnego opisu |
|---|---|---|
| Piętnaście modułów: Studio, Workspace, Automations, Browser, Research, Library, Translate, Roundtable, Design, Assistant, Terminal, Developer, Diagnostics, Apps, Agents | W zakresie — pełny opis operacyjny | Rozdział 4 tego dokumentu |
| Środowisko MultitaskingAI jako całość: warstwa centralna, role, silnik kolejek, orkiestracja, integracja z modułem Automations, panel orkiestracji zastępujący boczną nawigację modułów | Poza zakresem specyfikacji modułowej | Rozdział 13 Koncepcji platformy |
| Zastosowanie modułów w środowisku MultitaskingAI: Agents jako dostawca wykonawców do ról, Automations jako dostawca harmonogramów i kolejek | Odnotowane w podpunkcie „Powiązania konfigurowalne” właściwego modułu | Rozdział 4 tego dokumentu, z odesłaniem do rozdziału 13 Koncepcji |

Poniższy schemat przedstawia granicę zakresu wraz z punktami styku między modułami a środowiskiem MultitaskingAI.

```
   W ZAKRESIE SPECYFIKACJI                    POZA ZAKRESEM (Koncepcja, rozdz. 13)
   ─────────────────────────                  ──────────────────────────────────
   Piętnaście modułów platformy:               Środowisko MultitaskingAI (całość):
     Studio · Workspace · Automations            warstwa centralna
     Browser · Research · Library                role wykonawcze
     Translate · Roundtable · Design             silnik kolejek · orkiestracja
     Assistant · Terminal · Developer            integracja z modułem Automations
     Diagnostics · Apps · Agents                 panel orkiestracji (boczna nawigacja)
              │                                             ▲
              │   punkty styku odnotowane w                 │
              └──►  „Powiązania konfigurowalne”:  ──────────┘
                      Agents      ──► wykonawcy przypisywani do ról
                      Automations ──► harmonogramy i kolejki pętli pracy ciągłej
```

Każdy moduł opisano w jednolitej strukturze sześciu elementów, stosowanej jednakowo do wszystkich piętnastu modułów, niezależnie od tego, czy moduł jest dostępny jako okno w bocznej nawigacji środowisk, jako komponent własny konfigurowany na stronie głównej, czy w obu formach jednocześnie (rozróżnienie wyjaśnia rozdział 2).

| Element opisu modułu | Co zawiera |
|---|---|
| Cel | Przedmiot pracy modułu wraz z przeznaczeniem — dla kogo i po co moduł istnieje |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca) — elementy pierwszoplanowe każdego modułu |
| Okna operacyjne | Zestaw okien roboczych modułu, funkcja każdego z nich oraz warstwa widoczności |
| Funkcjonalności | Zakres operacji dostępnych w module |
| Przebieg pracy modułu | Typowy, domyślny przebieg pracy krok po kroku |
| Dane wykorzystywane przez AI | Źródła danych zasilające model w obrębie modułu |
| Powiązania konfigurowalne | Jawne powiązania z innymi modułami i funkcjami, ustanawiane przez użytkownika w oknie konfiguracji |

Specyfikacja pozostaje w pełnej zgodności z zasadami nadrzędnymi ustalonymi w rozdziale 14 Koncepcji platformy: pełną kompozycyjnością, pełną konfigurowalnością, jawnością i konfigurowalnością zależności oraz kierunkiem rozszerzania orkiestracji. W szczególności podpunkty „Powiązania konfigurowalne” nie stanowią zależności wbudowanych na stałe — każde powiązanie jest ustawieniem konfiguracyjnym ustanawianym z poziomu okna konfiguracji, o stanie domyślnym „nieaktywne” (rozdz. 6 i rozdz. 14, zasada 3 Koncepcji). Zgodność ze wszystkimi czterema zasadami omawia zbiorczo rozdział 6.

---

## 2. Moduł a komponent własny

Rozróżnienie między modułem a komponentem własnym, wprowadzone w rozdziałach 2.2 i 2.3 Koncepcji platformy, jest warunkiem poprawnej interpretacji specyfikacji zawartej w rozdziale 4. Poniższa tabela zestawia oba pojęcia.

| Cecha | Moduł | Komponent własny |
|---|---|---|
| Charakter | Element platformy dostępny dla wszystkich użytkowników | Nazwany wytwór konkretnego użytkownika |
| Powstanie | Wyspecjalizowany obszar roboczy dostarczony przez platformę | Tworzony i zapisywany w oknie konfiguracji na stronie głównej (strefa 2) lub w oknie modułu |
| Wyposażenie | Własny zestaw okien operacyjnych, dedykowane narzędzia, przebieg pracy oraz typ danych wykorzystywanych przez AI | Konkretna, zindywidualizowana instancja: automatyka, agent, projekt lub profil asystenta |
| Liczba na platformie | Piętnaście modułów wymienionych w rozdziale 4 | Dowolna liczba, zależnie od użytkownika |
| Dostępność | Zgodnie z macierzą dostępności modułów (rozdz. 3) | Wybieralny w trakcie sesji w module, w którym ma zastosowanie |

Cztery moduły platformy — Automations, Agents, Workspace i Assistant — są jednocześnie punktem wejścia do tworzenia komponentów własnych na stronie głównej; pozostałe jedenaście modułów nie tworzy komponentów własnych i jest dostępnych wyłącznie jako okno modułowe w bocznej nawigacji środowisk.

| Moduł | Ma okno modułowe w bocznej nawigacji środowisk | Jest punktem wejścia do komponentu własnego (strefa 2 strony głównej) | Rodzaj tworzonego komponentu |
|---|---|---|---|
| Studio | Tak | Nie | — |
| Workspace | Tak | Tak | Projekt |
| Automations | Nie — wyłącznie strefa 2 | Tak | Automatyka |
| Browser | Tak | Nie | — |
| Research | Tak | Nie | — |
| Library | Tak | Nie | — |
| Translate | Tak | Nie | — |
| Roundtable | Tak | Nie | — |
| Design | Tak | Nie | — |
| Assistant | Ma pozycję na lewym pasku nawigacji, lecz poza grupą środowisk — jest modułem aplikacji, nie modułem środowiska; konfiguracja w strefie 2 | Tak | Profil asystenta |
| Terminal | Tak | Nie | — |
| Developer | Tak | Nie | — |
| Diagnostics | Tak | Nie | — |
| Apps | Tak | Nie | — |
| Agents | Tak | Tak | Agent |

Utworzony komponent własny trafia do pamięci aplikacji jako zasób użytkownika i pozostaje wybieralny później, w module, w którym ma zastosowanie — niezależnie od tego, w którym module lub oknie konfiguracji powstał. Poniższy schemat przedstawia ten cykl.

```
   MIEJSCE UTWORZENIA                              WYKORZYSTANIE OPERACYJNE
   ─────────────────────                          ──────────────────────────
   Strona główna, strefa 2                          Sesja modułu docelowego, na przykład:
   (Automations · Agents ·        ─────►             automatyka  ──► sesja modułu Developer
    Workspace · Assistant)           │               agent       ──► rola w środowisku
        lub                          │                              MultitaskingAI albo
   okno modułu                       ▼                              projekt w module Workspace
                            Pamięć aplikacji                        (przez Agent Manager)
                            (zasób użytkownika,
                             wybieralny później)
```

**Maszyna stanów komponentu własnego.** Automatyka, agent, projekt i profil asystenta przechodzą przez ten sam ciąg stanów, niezależnie od modułu, w którym powstają — okno konfiguracji nadaje im wyłącznie odmienny zestaw pól.

```
                    ┌──────────────────┐
        edycja      │      ROBOCZY     │  utworzenie w oknie konfiguracji
     ┌──────────────┤  (w budowie)     │  (strona główna strefa 2 albo okno modułu)
     │              └────────┬─────────┘
     │                       │ zapis
     │                       ▼
     │              ┌──────────────────┐
     │              │     ZAPISANY     │  zasób użytkownika w pamięci aplikacji,
     │              │   (nieaktywny)   │  wybieralny w module, do którego należy
     │              └────────┬─────────┘
     │                       │ wybór w sesji / przypisanie do roli lub projektu
     │                       ▼
     │              ┌──────────────────┐
     └─────────────►│      AKTYWNY     │  wykorzystywany operacyjnie: automatyka
                    │   (w użyciu)     │  uruchomiona, agent jako Wykonawca,
                    └────────┬─────────┘  projekt otwarty, profil w rozmowie
                             │ zakończenie użycia / usunięcie
                             ▼
                    ┌──────────────────┐
                    │     USUNIĘTY     │  komponent trwale usuwany z pamięci
                    │                  │  aplikacji na decyzję użytkownika
                    └──────────────────┘
```

Przejście `AKTYWNY → ROBOCZY` (edycja w toku użycia) jest dopuszczalne wyłącznie dla automatyki i agenta — projekt oraz profil asystenta wracają do stanu ROBOCZY wyłącznie przez bezpośrednie otwarcie okna konfiguracji z poziomu stanu ZAPISANY. Stan USUNIĘTY jest stanem końcowym bez przejścia powrotnego; ponowne użycie wymaga utworzenia nowego komponentu własnego.

---

## 3. Macierz dostępności modułów

Moduł nie należy do jednego środowiska na wyłączność — jest dostępny w środowiskach, dla których stanowi naturalne narzędzie pracy. Poniższa macierz, zgodna z rozdziałem 10 Koncepcji platformy, wskazuje dostępność każdego z piętnastu modułów w trzech środowiskach udostępniających boczną nawigację modułów: TalkIn, WorkSpace i CodeStudio. Środowisko MultitaskingAI nie ma bocznej nawigacji modułów — udostępnia w jej miejsce panel orkiestracji opisany w rozdziale 13.8 Koncepcji — dlatego nie występuje jako kolumna macierzy.

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|---|---|---|
| Studio | TAK | TAK | NIE |
| Workspace | TAK | TAK | TAK |
| Automations | — | — | — |
| Browser | TAK | TAK | NIE |
| Research | TAK | TAK | NIE |
| Library | TAK | TAK | NIE |
| Translate | TAK | NIE | NIE |
| Roundtable | TAK | TAK | TAK |
| Design | NIE | TAK | TAK |
| Assistant | poza grupą | poza grupą | poza grupą |
| Terminal | NIE | NIE | TAK |
| Developer | NIE | NIE | TAK |
| Diagnostics | NIE | NIE | TAK |
| Apps | NIE | TAK | TAK |
| Agents | TAK | TAK | TAK |

Ujęcie macierzy z perspektywy środowiska — które moduły widnieją w bocznej nawigacji każdego z nich — przedstawia poniższy schemat.

```
   TalkIn ─────────────────► Studio · Workspace · Browser · Research · Library ·
                             Translate · Roundtable · Agents
   Poza grupą środowisk ───► Assistant (moduł aplikacji — własna pozycja
                             na lewym pasku nawigacji)
   WorkSpace ──────────────► Studio · Workspace · Browser · Research · Library ·
                             Roundtable · Design · Apps · Agents
   CodeStudio ─────────────► Workspace · Roundtable · Design · Terminal ·
                             Developer · Diagnostics · Apps · Agents
   MultitaskingAI ─────────► bez bocznej nawigacji modułów — panel orkiestracji
                             (rozdz. 13.8 Koncepcji)
   Strona główna, strefa 2 ─► Automations · Agents · Workspace · Assistant
                             (konfiguracja komponentów własnych)
```

Niezależnie od środowiska każdy moduł udostępnia oba okna komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca) — jako elementy pierwszoplanowe warstwy 1, w lewych kolumnach obszaru roboczego.

Cztery wiersze macierzy wymagają objaśnienia — zbiera je poniższa tabela.

| Moduł lub środowisko | Status w macierzy | Wyjaśnienie |
|---|---|---|
| Automations | „—” we wszystkich trzech kolumnach | Nie ma okna modułowego w żadnym środowisku; działa wyłącznie ze strony głównej (strefa 2); gotowe automatyki trafiają do sesji modułów docelowych jako komponenty własne |
| Assistant | „poza grupą” we wszystkich trzech kolumnach | Ma własną pozycję na lewym pasku nawigacji, lecz pozycja ta stoi poza grupą środowisk — Assistant jest modułem aplikacji, nie modułem środowiska, więc macierz środowiskowa nie rozstrzyga o jego dostępności. Konfigurowany jako komponent własny w strefie 2 strony głównej |
| Agents oraz Workspace | „TAK” w macierzy i dodatkowo konfiguracja na stronie głównej | Konfigurowane w strefie 2 strony głównej i jednocześnie dostępne jako okno modułowe we wskazanych środowiskach |
| MultitaskingAI | Brak jako kolumna | Nie ma bocznej nawigacji modułów; udostępnia panel orkiestracji (rozdz. 13.8 Koncepcji), dlatego nie występuje jako kolumna macierzy |

---

## 3a. Warstwy widoczności funkcji modułu

Interfejs każdego modułu podlega zasadzie nadrzędnej: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Moduł ujawnia swoje możliwości stopniowo — zależnie od kontekstu pracy, roli użytkownika i wykonywanej czynności — zachowując pełną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność modułu istnieje w jego architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba okien, paneli, narzędzi, ustawień i funkcji administracyjnych modułu nie wpływa na postrzeganą prostotę jego interfejsu.

Każde okno, każdy panel i każda funkcja modułu należy do dokładnie jednej z czterech warstw widoczności. Warstwę podaje kolumna „Warstwa” w zestawieniach okien operacyjnych w rozdziale 4.

| Warstwa | Nazwa | Zawartość w module | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, Execution Loop Window, aktywne okno wiodące modułu, kontekst pracy, boczna nawigacja modułów, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy modułu | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji modułu, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel nakładkowy, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje modułu, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

**Mechanizmy ukrywania funkcjonalności modułu.**

| Mechanizm | Działanie |
|---|---|
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`), a lista pozycji rozwija się po kliknięciu |
| Panele wysuwane | Funkcjonalność umieszczana jest w panelach bocznych, panelach wysuwanych, oknach nakładkowych i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor |

**Zasada jednego kliknięcia.** Każda ukryta funkcja modułu jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

**Rozmieszczenie okien modułu.** Obowiązuje układ pionowy — podział lewa–prawa. Wszystkie okna robocze, okna komunikacji operacyjnej i okna pomocnicze modułu rozmieszczone są w kolumnach sąsiadujących poziomo; regulacji podlega wyłącznie szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Lewa, stała, pełna wysokość obszaru roboczego | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca (otwierana) | **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu — okna wiodące, edycyjne, podglądu i monitory |
| Kolejne kolumny boczne | Okna pomocnicze i panele modułu — otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego |

Poniższa makieta przedstawia kanoniczny układ okien modułu w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki kontekstowe, `▼`, `⋮`, `☰`).

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │ (okno wiodące)           │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │ [Danaco Console] [Fable  │  boczne)
            │                      │  5] [Ultra]   Agent ▼  ⋮ │
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ════════════════════════════════════════════════════════════════════════════
```

**Chat Window — kanał Użytkownik ↔ Wykonawca.** Główne okno komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy). Stanowi centralny punkt pracy użytkownika w module i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez moduł: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Każdy moduł udostępnia to okno w tym samym miejscu układu — w lewej kolumnie obszaru roboczego. Warstwa 1.

**Execution Loop Window — kanał Koordynator ↔ Wykonawca.** Okno pętli wykonawczej prezentujące komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów modułu. Zawiera bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Okno otwierane jest jako kolumna sąsiadująca z Chat Window. Warstwa 1.

Role: **Użytkownik** — zleca i zatwierdza; **Koordynator** — komponent orkiestrujący platformy, dekomponuje zlecenie, przydziela i nadzoruje zadania; **Wykonawca** — AI, agent lub system wykonawczy realizujący zadania.

---

## 3b. Macierz uprawnień — punkty izolacji i współdzielenie kontekstu

Dane wykorzystywane przez AI w każdym module — historia rozmów, pamięć, kontekst dokumentu, materiały z innych modułów — podlegają punktom izolacji konfigurowalnym w oknie opisanym w rozdziale 6.4–6.6 Koncepcji platformy. Poniższa macierz zestawia cztery poziomy zasięgu z czterema rodzajami danych współdzielonych oraz stan wyjściowy dostępu, obowiązujący przy braku jawnego ustawienia użytkownika.

| Poziom zasięgu | Historia rozmów | Pamięć | Kontekst dokumentu / sesji | Materiały z innych modułów |
|---|---|---|---|---|
| Globalny | Dziedziczony przez wszystkie moduły jako warstwa bazowa | Dostępna wszystkim modułom, o ile nie nadpisana niżej | Nie dotyczy — zasięg globalny nie niesie kontekstu pojedynczej sesji | Widoczność zależna od jawnego powiązania konfigurowalnego (rozdz. 5) |
| Środowisko | Dziedziczony przez moduły danego środowiska (TalkIn, WorkSpace, CodeStudio) | Dziedziczona przez moduły danego środowiska | Współdzielony między modułami tego samego środowiska, jeśli ustawiono jawnie | Ograniczony do modułów dostępnych w tym środowisku (macierz rozdz. 3) |
| Projekt | Odrębna dla każdego projektu w module Workspace | Odrębna dla każdego projektu, nadpisuje poziom środowiska | Współdzielony przez wszystkie moduły przypisane do projektu (Agent Manager, Instructions Panel) | Widoczność ograniczona do modułów powiązanych z danym projektem |
| Sesja | Odrębna dla każdej nowej karty sesji — stan wyjściowy | Odrębna dla sesji, nadpisuje poziom projektu i środowiska | Właściwa wyłącznie bieżącej karcie sesji | Brak dostępu bez jawnego powiązania konfigurowalnego dla tej sesji |

**Reguła pierwszeństwa.** Ustawienie na poziomie bardziej szczegółowym (sesja) nadpisuje ustawienie na poziomie szerszym (projekt, środowisko, globalny); brak ustawienia na poziomie szczegółowym oznacza dziedziczenie z poziomu bezpośrednio szerszego, aż do poziomu globalnego stanowiącego warstwę bazową (rozdz. 6.5 Koncepcji platformy). Żadna komórka macierzy nie jest domyślnie współdzielona między modułami bez jawnego powiązania konfigurowalnego odnotowanego w rozdziale 5 — stan wyjściowy izolacji między modułami jest rozłączny, a każde odstępstwo od niego jest świadomą decyzją użytkownika.

---

## 4. Specyfikacja modułów

Poniżej pełna specyfikacja operacyjna piętnastu modułów, w kolejności przyjętej w rozdziale 11 Koncepcji platformy. Każdy moduł otwiera karta modułu zestawiająca cel, okna komunikacji operacyjnej, okna robocze, funkcjonalności i powiązania, po której następują tabele szczegółowe: okien operacyjnych wraz z warstwami widoczności (rozdz. 3a), przebiegu pracy, danych wykorzystywanych przez AI oraz powiązań konfigurowalnych. Chat Window oraz Execution Loop Window występują w każdym module jako elementy pierwszoplanowe. Środowiska, w których dany moduł jest dostępny, wskazuje macierz w rozdziale 3. Przegląd wszystkich piętnastu modułów przedstawia poniższa tabela.

Kolumna „Obszar kontraktu” wskazuje nazwę obszaru komend w `budowa/shared/contract.json` (klucz `obszary`, pole `nazwa`) realizującego dany moduł; wszystkie piętnaście wartości potwierdzono odczytem pola `nazwa` z listy `obszary` — każda nazwa jest obecna wśród 68 obszarów kontraktu.

| Moduł | Cel operacyjny (w skrócie) | Forma udostępnienia | Obszar kontraktu | Okna komunikacji operacyjnej |
|---|---|---|---|---|
| Studio | Zaawansowana praca z tekstem, dokumentami i treścią | Okno modułowe | `studio` | Chat Window · Execution Loop Window |
| Workspace | Izolowane środowisko projektowe | Komponent własny (Projekt) oraz okno modułowe | `workspace` | Chat Window · Execution Loop Window |
| Automations | Budowanie i wykonywanie procesów automatycznych | Komponent własny (Automatyka), wyłącznie strefa 2 | `automation` | Chat Window · Execution Loop Window |
| Browser | Współdzielone przeglądanie internetu | Okno modułowe | `browser` | Chat Window · Execution Loop Window |
| Research | Realizacja badań, analiz i opracowań | Okno modułowe | `research` | Chat Window · Execution Loop Window |
| Library | Centralne repozytorium wiedzy i plików | Okno modułowe | `library` | Chat Window · Execution Loop Window |
| Translate | Wielojęzyczne tłumaczenia w trybie wielozadaniowym | Okno modułowe | `translate` | Chat Window · Execution Loop Window |
| Roundtable | Współpraca wielu modeli AI nad wspólnym problemem | Okno modułowe | `roundtable` | Chat Window · Execution Loop Window |
| Design | Projektowanie i tworzenie zasobów wizualnych | Okno modułowe | `design` | Chat Window · Execution Loop Window |
| Assistant | Naturalna komunikacja głosowa z AI | Komponent własny (Profil asystenta); moduł aplikacji z własną pozycją na lewym pasku nawigacji, poza grupą środowisk | `assistant` | Chat Window · Execution Loop Window |
| Terminal | Praca z konsolami i środowiskami wykonawczymi | Okno modułowe | `terminal` | Chat Window · Execution Loop Window |
| Developer | Tworzenie i rozwój kodu | Okno modułowe | `developer` | Chat Window · Execution Loop Window |
| Diagnostics | Analiza oraz usuwanie problemów technicznych | Okno modułowe | `diagnostics` | Chat Window · Execution Loop Window |
| Apps | Budowa kompletnych produktów cyfrowych | Okno modułowe | `apps` | Chat Window · Execution Loop Window |
| Agents | Tworzenie i zarządzanie własnymi agentami AI | Komponent własny (Agent) oraz okno modułowe | `agent` | Chat Window · Execution Loop Window |

### 4.1. Studio

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Zaawansowana praca z tekstem, dokumentami oraz treścią |
| Przeznaczenie | Użytkownicy tworzący, redagujący i przekształcający materiały pisane — od pojedynczych notatek po obszerne dokumenty wielostronicowe; potrzeba jednego miejsca trwale integrującego edycję treści, wsparcie AI i porównywanie wersji dokumentu, zamiast rozpraszania między osobny edytor a osobne okno rozmowy |
| Środowiska dostępności | TalkIn, WorkSpace |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window |
| Funkcjonalności | Edycja dokumentów; obsługa formatów PDF, DOCX, TXT, Markdown; tłumaczenia; korekta; analiza treści; przepisywanie; zmiana stylu; streszczenia; rozwijanie treści; Diff; Grep; operacje kontekstowe AI |
| Powiązania konfigurowalne | Translate, Library, Research |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście edytowanego dokumentu; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Studio Editor | Edycja treści dokumentu | 1 |
| Tools Panel | Dostęp do operacji kontekstowych AI (korekta, przepisywanie, streszczenie, zmiana stylu) | 2 |
| Diff/Grep Panel | Porównywanie wersji dokumentu i wyszukiwanie w treści | 2 |
| Session Repository | Historia wersji pracy nad dokumentem | 3 |
| Preview Window | Podgląd wyniku pracy | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Wczytanie dokumentu do edytora | Studio Editor |
| 2 | Zlecenie AI operacji kontekstowej (korekta, zmiana stylu, streszczenie) na wybranym fragmencie | Chat Window lub Tools Panel |
| 3 | Porównanie wersji przed i po zmianie | Diff/Grep Panel |
| 4 | Akceptacja lub odrzucenie wyniku | Decyzja użytkownika (Studio Editor) |
| 5 | Zapis kolejnej wersji z możliwością powrotu bez utraty historii zmian | Session Repository |

**Dane wykorzystywane przez AI.**
- Treść dokumentu edytowanego w Studio Editor
- Wynik operacji Diff/Grep
- Historia wersji dokumentu z Session Repository
- Materiały wejściowe i artefakty wyjściowe wymieniane z Library, jeśli powiązanie ustanowiono
- Zakres pamięci sesji lub projektu, zgodnie z konfiguracją izolacji
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Translate | Operacje kontekstowe AI są mechanizmem natywnym Studio; ich wspólne wykorzystanie z Translate przy pracy dwujęzycznej jest możliwością do skonfigurowania przez użytkownika, a nie mechanizmem współdzielonym domyślnie | Konfiguracyjne |
| Library | Repozytorium źródłowe dla dokumentów wejściowych i miejsce docelowe dla wygenerowanych artefaktów, zgodnie z powiązaniem ustanowionym przez użytkownika | Konfiguracyjne |
| Research | Wyniki pracy w Studio mogą, decyzją użytkownika, być dalej przetwarzane w module Research przy budowie raportów końcowych | Konfiguracyjne |

### 4.2. Workspace

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Izolowane środowisko projektowe |
| Przeznaczenie | Użytkownicy prowadzący równolegle wiele odrębnych projektów, którym moduł pozwala skonfigurować rozdzielenie kontekstu, instrukcji i pamięci między nimi, tak aby ustalenia jednego projektu nie przenikały do innego wtedy, gdy izolacja jest pożądana |
| Środowiska dostępności | TalkIn, WorkSpace, CodeStudio |
| Forma udostępnienia | Komponent własny (Projekt) konfigurowany w strefie 2 strony głównej oraz okno modułowe w bocznej nawigacji |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager |
| Funkcjonalności | Separacja projektów; własne instrukcje; własna pamięć; własna biblioteka; własna konfiguracja modeli; zarządzanie agentami |
| Powiązania konfigurowalne | Agents, Library, Automations |

**Uwaga terminologiczna.** Moduł Workspace jest jednym z elementów warstwy modułów (Workspace Layer, rozdz. 5.2 Koncepcji), a nie jej odpowiednikiem — nazwa modułu odnosi się do konkretnej, izolowanej przestrzeni projektowej, nie do całej warstwy architektonicznej.

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście projektu; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Project Dashboard | Przegląd stanu i postępu projektu | 1 |
| Instructions Panel | Definiowanie instrukcji systemowych projektu | 2 |
| Context Memory | Pamięć kontekstowa dedykowana projektowi | 2 |
| Project Library | Biblioteka materiałów projektu | 3 |
| Agent Manager | Przypisywanie agentów do projektu | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Założenie nowego projektu jako komponentu własnego | Moduł Workspace (strefa 2) |
| 2 | Zdefiniowanie odrębnego zestawu instrukcji systemowych | Instructions Panel |
| 3 | Zbudowanie dedykowanej pamięci kontekstowej | Context Memory |
| 4 | Przypisanie do projektu konkretnych agentów | Agent Manager |
| 5 | Ustalenie, czy i w jakim stopniu konfiguracja projektu wpływa na inny, równolegle prowadzony projekt | Okno konfiguracji |

**Dane wykorzystywane przez AI.**
- Instrukcje systemowe projektu z Instructions Panel
- Pamięć kontekstowa projektu z Context Memory
- Materiały z Project Library
- Lista i konfiguracja przypisanych agentów z Agent Manager
- Stan i postęp projektu z Project Dashboard
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Agents | Korzystanie z agentów skonfigurowanych w module Agents, udostępnianych — decyzją użytkownika — jako wykonawcy zadań w ramach projektu | Konfiguracyjne |
| Library | Project Library pełni funkcję analogiczną do modułu Library, w zakresie ograniczonym do jednego projektu | Konfiguracyjne |
| Automations | Projekty prowadzone w Workspace mogą zostać powiązane, decyzją użytkownika, z procesami automatycznymi tworzonymi w module Automations | Konfiguracyjne |

### 4.3. Automations

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Budowanie i wykonywanie procesów automatycznych |
| Przeznaczenie | Użytkownicy przekształcający powtarzalne czynności — cykliczne raporty, regularne przetwarzanie danych, zaplanowane wywołania modeli — w procesy działające bez stałego nadzoru; operacyjne zaplecze zadań uruchamianych według harmonogramu lub zdarzenia, a nie w trybie konwersacyjnym |
| Środowiska dostępności | Brak okna modułowego w jakimkolwiek środowisku — wyłącznie strona główna (strefa 2) |
| Forma udostępnienia | Komponent własny (Automatyka); gotowe automatyki trafiają do sesji docelowych modułów jako komponenty własne |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Workflow Builder, Scheduler, Queue Manager, Orchestrator, Execution Monitor |
| Funkcjonalności | Harmonogramy; zadania cykliczne; przebieg pracy; kolejki; automatyczne wywołania modeli; pętle wielomodelowe; orkiestracja |
| Powiązania konfigurowalne | MultitaskingAI, dowolny moduł platformy |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI przy budowie procesu; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Workflow Builder | Projektowanie procesu automatycznego | 1 |
| Scheduler | Ustalanie harmonogramu i cykliczności | 2 |
| Queue Manager | Zarządzanie kolejką zadań procesu | 2 |
| Orchestrator | Zależności między etapami procesu | 3 |
| Execution Monitor | Obserwacja statusu kolejnych uruchomień i błędów | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Zaprojektowanie procesu automatycznego | Workflow Builder |
| 2 | Ustalenie cykliczności procesu | Scheduler |
| 3 | Konfiguracja kolejki zadań i zależności między etapami | Queue Manager, Orchestrator |
| 4 | Obserwacja kolejnych uruchomień, statusu każdego przebiegu oraz błędów wymagających interwencji | Execution Monitor |
| 5 | Wpięcie gotowej automatyki jako komponentu własnego do sesji modułu docelowego | Moduł docelowy |

**Dane wykorzystywane przez AI.**
- Definicja procesu z Workflow Builder
- Reguły harmonogramu ze Scheduler
- Stan i pozycje kolejki z Queue Manager
- Zależności procesu z Orchestrator
- Status i logi wykonania z Execution Monitor
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| MultitaskingAI | Po skonfigurowaniu integracji (rozdz. 13.6 Koncepcji) Automations może stanowić operacyjne zaplecze harmonogramów i kolejek dla środowiska MultitaskingAI; silnik kolejek modułu Automations może zostać powiązany z silnikiem kolejek środowiska MultitaskingAI. Połączenie nie jest domyślne i wymaga decyzji użytkownika w oknie konfiguracji oraz w ustawieniach okna modułu Automations | Konfiguracyjne |
| Dowolny moduł platformy | Procesy zdefiniowane w Automations mogą obejmować zadania z dowolnego innego modułu, zgodnie z zakresem powiązań ustanowionym przez użytkownika | Konfiguracyjne |

### 4.4. Browser

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Współdzielone przeglądanie internetu |
| Przeznaczenie | Sytuacje, w których użytkownik i AI muszą pracować nad tą samą treścią internetową w tym samym czasie — analizując stronę, porównując oferty lub weryfikując informację wspólnie, zamiast przekazywać sobie nawzajem linki i fragmenty tekstu |
| Środowiska dostępności | TalkIn, WorkSpace |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Browser Window, Sources Panel, Notes Panel |
| Funkcjonalności | Analiza stron; wyszukiwanie informacji; wspólny podgląd użytkownika i AI; analiza dokumentów online |
| Powiązania konfigurowalne | Research, Library |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście przeglądanej treści; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Browser Window | Wspólny podgląd strony przez użytkownika i AI | 1 |
| Sources Panel | Lista źródeł wykorzystanych w toku pracy | 2 |
| Notes Panel | Notatki z istotnych fragmentów | 2 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Otwarcie strony | Browser Window |
| 2 | Odpowiedzi AI na pytania o treść i wyszukiwanie powiązanych informacji przy wspólnym podglądzie | Browser Window, Chat Window |
| 3 | Odnotowanie istotnych fragmentów | Notes Panel |
| 4 | Gromadzenie wykorzystanych źródeł | Sources Panel |

**Dane wykorzystywane przez AI.**
- Treść przeglądanej strony z Browser Window
- Lista źródeł z Sources Panel
- Notatki z Notes Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Research | Źródła zebrane w module Browser mogą, decyzją użytkownika, zasilać dalszą pracę badawczą w module Research | Konfiguracyjne |
| Library | Odnotowane materiały mogą trafiać do Library w celu trwałego przechowania, zgodnie z powiązaniem ustanowionym przez użytkownika | Konfiguracyjne |

### 4.5. Research

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Realizacja badań, analiz i opracowań |
| Przeznaczenie | Użytkownicy prowadzący pracę wymagającą systematycznego zbierania, porządkowania i syntezowania informacji z wielu źródeł — od analiz rynkowych po opracowania konkurencyjne — w sposób bardziej ustrukturyzowany niż pojedyncza rozmowa z modelem |
| Środowiska dostępności | TalkIn, WorkSpace |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Research Workspace, Sources Manager, Findings Panel, Report Builder, Export Panel |
| Funkcjonalności | Raporty; analizy rynku; benchmarki; analiza konkurencji; wieloźródłowe badania |
| Powiązania konfigurowalne | Browser, Library, Studio, Roundtable |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku badania; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Research Workspace | Główna przestrzeń prowadzenia badania | 1 |
| Sources Manager | Gromadzenie i zarządzanie źródłami | 2 |
| Findings Panel | Odnotowywanie ustaleń cząstkowych | 2 |
| Report Builder | Kompletowanie raportu końcowego | 3 |
| Export Panel | Eksport raportu do formatu docelowego | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Gromadzenie źródeł | Sources Manager |
| 2 | Odnotowywanie ustaleń cząstkowych w miarę postępu analizy | Findings Panel |
| 3 | Kompletowanie raportu końcowego | Report Builder |
| 4 | Eksport raportu do formatu wymaganego przez odbiorcę | Export Panel |

**Dane wykorzystywane przez AI.**
- Zebrane źródła z Sources Manager
- Ustalenia cząstkowe z Findings Panel
- Treść raportu z Report Builder
- Format eksportu z Export Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Browser | Korzystanie ze źródeł zebranych w module Browser, zgodnie z powiązaniem skonfigurowanym przez użytkownika | Konfiguracyjne |
| Library | Korzystanie z dokumentów przechowywanych w Library, zgodnie z powiązaniem skonfigurowanym przez użytkownika | Konfiguracyjne |
| Studio | Wyniki pracy badawczej mogą być dalej redagowane w module Studio, decyzją użytkownika | Konfiguracyjne |
| Roundtable | Wyniki mogą być prezentowane wielomodelowo w module Roundtable przed opracowaniem ostatecznych wniosków, decyzją użytkownika | Konfiguracyjne |

### 4.6. Library

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Centralne repozytorium wiedzy i plików |
| Przeznaczenie | Użytkownicy potrzebujący jednego, uporządkowanego miejsca przechowywania materiałów wykorzystywanych i wytwarzanych w toku pracy z platformą; moduł pełni funkcję trwałej pamięci zewnętrznej wspólnej dla wielu środowisk i modułów |
| Środowiska dostępności | TalkIn, WorkSpace |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Library Explorer, Tags & Collections, File Preview, Versioning Panel |
| Funkcjonalności | Przechowywanie plików; katalogowanie wiedzy; odbiór artefaktów generowanych przez AI; zarządzanie dokumentami |
| Powiązania konfigurowalne | Studio, Research, Browser, Workspace |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście zgromadzonych materiałów; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Library Explorer | Przegląd i zarządzanie zgromadzonymi materiałami | 1 |
| Tags & Collections | Katalogowanie materiałów | 2 |
| File Preview | Podgląd zawartości pliku bez opuszczania modułu | 2 |
| Versioning Panel | Śledzenie kolejnych wersji dokumentu | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Napływ materiałów jako plików wejściowych lub artefaktów generowanych przez AI w innych modułach | Library Explorer |
| 2 | Porządkowanie materiałów | Tags & Collections |
| 3 | Przegląd zawartości bez opuszczania modułu | File Preview |
| 4 | Śledzenie kolejnych wersji tego samego dokumentu, w tym wygenerowanych automatycznie przez AI | Versioning Panel |

**Dane wykorzystywane przez AI.**
- Metadane i treść przechowywanych plików z Library Explorer
- Struktura tagów i kolekcji
- Historia wersji dokumentu z Versioning Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Studio, Research, Browser | Rola centralnego repozytorium dla materiałów wykorzystywanych w tych modułach, zgodnie z powiązaniami ustanowionymi przez użytkownika | Konfiguracyjne |
| Workspace | Project Library w module Workspace stanowi odpowiednik Library ograniczony do zakresu pojedynczego projektu | Konfiguracyjne |

### 4.7. Translate

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Wielojęzyczne tłumaczenia w trybie wielozadaniowym |
| Przeznaczenie | Użytkownicy pracujący równolegle z treścią w wielu językach — zespoły tłumaczeniowe, lokalizacyjne lub użytkownicy indywidualni obsługujący wielojęzyczną komunikację — którym nie wystarcza pojedyncze, sekwencyjne tłumaczenie fragmentów tekstu w oknie czatu |
| Środowiska dostępności | TalkIn |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Source Panel, Translation Panels, Glossary Manager |
| Funkcjonalności | Równoległe tłumaczenia; obsługa wielu języków; zarządzanie terminologią; lokalizacja treści |
| Powiązania konfigurowalne | Studio |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście tłumaczenia; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Source Panel | Tekst źródłowy | 1 |
| Translation Panels | Równoległe tłumaczenia na wiele języków | 2 |
| Glossary Manager | Zarządzanie terminologią i spójnością tłumaczeń | 2 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Umieszczenie tekstu źródłowego | Source Panel |
| 2 | Jednoczesne tłumaczenie na wiele języków widocznych w osobnych panelach | Translation Panels |
| 3 | Zapewnienie spójności terminologii przez definiowanie preferowanych odpowiedników kluczowych pojęć | Glossary Manager |

**Dane wykorzystywane przez AI.**
- Tekst źródłowy z Source Panel
- Treść tłumaczeń równoległych z Translation Panels
- Terminologia z Glossary Manager
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Studio | Translate może zostać skonfigurowany do korzystania z operacji kontekstowych AI właściwych modułowi Studio przy pracy dwujęzycznej i może pełnić rolę etapu pośredniego, gdy dokument wymaga wersji wielojęzycznej. Mechanizm operacji kontekstowych AI pozostaje natywnie mechanizmem Studio; jego wspólne wykorzystanie ustanawia użytkownik w oknie konfiguracji | Konfiguracyjne |

### 4.8. Roundtable

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Współpraca wielu modeli AI nad wspólnym problemem |
| Przeznaczenie | Zadania, w których pojedyncza odpowiedź jednego modelu jest niewystarczająca — gdy wartość wynika z konfrontacji różnych perspektyw, wzajemnej krytyki i wypracowania wspólnego stanowiska, a nie z jednej, izolowanej odpowiedzi |
| Środowiska dostępności | TalkIn, WorkSpace, CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Model Panels, Debate Panel, Moderator Panel, Consensus Panel |
| Funkcjonalności | Współpraca modeli; debaty; wymiana argumentów; porównywanie odpowiedzi; wypracowywanie konsensusu |
| Powiązania konfigurowalne | MultitaskingAI (pokrewieństwo funkcjonalne) |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku dyskusji; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Model Panels | Równoległe odpowiedzi kilku modeli | 1 |
| Debate Panel | Rejestr wymiany argumentów między modelami | 2 |
| Moderator Panel | Ukierunkowanie dyskusji przez użytkownika | 2 |
| Consensus Panel | Finalne, uzgodnione stanowisko | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Równoległe odpowiedzi kilku modeli na to samo zagadnienie | Model Panels |
| 2 | Rejestr wymiany argumentów między modelami | Debate Panel |
| 3 | Ukierunkowanie dyskusji przez użytkownika | Moderator Panel |
| 4 | Zgromadzenie finalnego, uzgodnionego stanowiska | Consensus Panel |

**Dane wykorzystywane przez AI.**
- Odpowiedzi poszczególnych modeli z Model Panels
- Przebieg debaty z Debate Panel
- Interwencje moderujące z Moderator Panel
- Wypracowane stanowisko z Consensus Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| MultitaskingAI | Roundtable stanowi funkcjonalne rozwinięcie idei wielomodelowości realizowanej w pełniejszej, zorientowanej na role formie przez środowisko MultitaskingAI (rozdz. 9.4 i 13 Koncepcji); Roundtable koncentruje się na debacie i konsensusie, a nie na podziale ról wykonawczych i orkiestracji procesu. Nie jest to powiązanie łączące kontekst, lecz pokrewieństwo funkcjonalne między dwoma odrębnymi mechanizmami platformy | Pokrewieństwo funkcjonalne |

### 4.9. Design

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Projektowanie i tworzenie zasobów wizualnych |
| Przeznaczenie | Użytkownicy potrzebujący generować i edytować materiały graficzne — od ilustracji i elementów brandingowych po makiety interfejsu — z bezpośrednim wsparciem AI na każdym etapie procesu twórczego |
| Środowiska dostępności | WorkSpace, CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Design Board, Assets Panel, Prompt Builder, Preview Window |
| Funkcjonalności | Generowanie grafiki; edycja obrazów; ilustracje; branding; UI/UX; materiały marketingowe |
| Powiązania konfigurowalne | Studio, Apps |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku pracy twórczej; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Design Board | Zestawianie i dalsza praca koncepcyjna nad kompozycją | 1 |
| Assets Panel | Wygenerowane i zgromadzone zasoby wizualne | 2 |
| Prompt Builder | Precyzyjne formułowanie poleceń generujących grafikę | 2 |
| Preview Window | Podgląd wyniku pracy | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Sformułowanie polecenia generującego grafikę | Prompt Builder |
| 2 | Napływ wygenerowanych i zgromadzonych zasobów | Assets Panel |
| 3 | Zestawianie zasobów i praca koncepcyjna nad spójną kompozycją wizualną | Design Board |

**Dane wykorzystywane przez AI.**
- Polecenia z Prompt Builder
- Zgromadzone zasoby z Assets Panel
- Kompozycja robocza z Design Board
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Studio | Zasoby wygenerowane w module Design mogą, w ramach powiązań ustanowionych przez użytkownika, być wykorzystywane w module Studio przy redagowaniu dokumentów | Konfiguracyjne |
| Apps | Zasoby mogą być wykorzystywane w module Apps przy budowie interfejsu produktu, zgodnie z powiązaniem ustanowionym przez użytkownika | Konfiguracyjne |

### 4.10. Assistant

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Naturalna komunikacja głosowa z AI |
| Przeznaczenie | Scenariusze, w których interakcja tekstowa jest mniej wygodna niż mowa — praca w ruchu, sterowanie zadaniami bez użycia klawiatury lub preferencja użytkownika co do formy komunikacji |
| Środowiska dostępności | TalkIn (wykorzystanie operacyjne) |
| Forma udostępnienia | Komponent własny (Profil asystenta) konfigurowany w strefie 2 strony głównej, wykorzystywany operacyjnie w środowisku TalkIn |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Voice Console, Actions Monitor, Activity Feed |
| Funkcjonalności | Komunikacja głosowa; wykonywanie poleceń; sterowanie zadaniami; obsługa aplikacji; realizacja działań wieloetapowych |
| Powiązania konfigurowalne | Always On Display (pokrewieństwo funkcjonalne) |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja tekstowa uzupełniająca; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Voice Console | Wydawanie poleceń głosowych | 1 |
| Actions Monitor | Bieżący status realizacji poleceń | 2 |
| Activity Feed | Chronologiczny zapis wykonanych działań | 2 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Konfiguracja profilu asystenta jako komponentu własnego | Strona główna (strefa 2) |
| 2 | Wydawanie poleceń głosowych w środowisku TalkIn | Voice Console |
| 3 | Bieżący status realizacji poleceń | Actions Monitor |
| 4 | Chronologiczny zapis wykonanych działań, pozwalający odtworzyć przebieg wieloetapowego zlecenia zrealizowanego głosowo | Activity Feed |

**Dane wykorzystywane przez AI.**
- Polecenia głosowe z Voice Console
- Status realizacji działań z Actions Monitor
- Chronologiczny zapis działań z Activity Feed
- Konfiguracja profilu asystenta
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Always On Display | Assistant współdzieli charakter komunikacji głosowej z funkcją globalną Always On Display (rozdz. 12.2 Koncepcji), różniąc się od niej tym, że działa jako pełnoprawny moduł osadzony w konkretnym środowisku (TalkIn), a nie jako warstwa obecna ponad całą platformą. Nie jest to powiązanie łączące kontekst, lecz pokrewieństwo funkcjonalne między modułem a funkcją globalną | Pokrewieństwo funkcjonalne |

### 4.11. Terminal

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Praca z konsolami i środowiskami wykonawczymi |
| Przeznaczenie | Użytkownicy techniczni potrzebujący bezpośredniego dostępu do powłok systemowych i narzędzi wiersza polecenia z poziomu platformy, bez konieczności przełączania się do zewnętrznej aplikacji terminala |
| Środowiska dostępności | CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Terminal Tabs, Output Console, Process Monitor |
| Funkcjonalności | PowerShell; CMD; Bash; Node.js; Python; narzędzia CLI |
| Powiązania konfigurowalne | Developer, Diagnostics |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście pracy w powłoce; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Terminal Tabs | Równoległe sesje w różnych powłokach | 1 |
| Output Console | Wynik wykonania poleceń | 2 |
| Process Monitor | Obserwacja i kontrola uruchomionych procesów | 2 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Otwarcie jednej lub kilku sesji powłok | Terminal Tabs |
| 2 | Gromadzenie wyniku poleceń | Output Console |
| 3 | Obserwacja i kontrola uruchomionych procesów, w tym zainicjowanych poleceniem wydanym przez AI | Process Monitor |

**Dane wykorzystywane przez AI.**
- Polecenia i wynik pracy powłok z Terminal Tabs i Output Console
- Stan uruchomionych procesów z Process Monitor
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Developer, Diagnostics | Terminal dostarcza warstwy wykonawczej dla poleceń wydawanych z poziomu modułów Developer i Diagnostics, z którymi jest najściślej powiązany funkcjonalnie w ramach środowiska CodeStudio. Współobecność tych modułów w jednym środowisku czyni powiązanie typowym w praktyce, lecz — zgodnie z zasadą jawności i konfigurowalności zależności (rozdz. 14, zasada 3 Koncepcji) — nie jest ono wbudowaną na stałe zależnością techniczną | Konfiguracyjne |

### 4.12. Developer

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Tworzenie i rozwój kodu |
| Przeznaczenie | Praca programistyczna wymagająca pełnego zestawu narzędzi inżynierskich — edytora kodu, kontroli wersji i podglądu wyniku budowania — zintegrowanych ze wsparciem AI działającym bezpośrednio w kontekście danego repozytorium |
| Środowiska dostępności | CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Code Editor, Project Tree, Git Panel, Build Output |
| Funkcjonalności | Generowanie kodu; refaktoryzacja; analiza architektury; dokumentacja; testowanie |
| Powiązania konfigurowalne | Terminal, Diagnostics, Apps |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w kontekście repozytorium; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Code Editor | Edycja kodu | 1 |
| Project Tree | Nawigacja po strukturze projektu | 2 |
| Git Panel | Zarządzanie zmianami i kontrola wersji | 2 |
| Build Output | Wynik kompilacji lub budowania | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Nawigacja po strukturze projektu | Project Tree |
| 2 | Edycja kodu przy wsparciu AI | Code Editor, Chat Window |
| 3 | Zarządzanie zmianami i kontrola wersji | Git Panel |
| 4 | Obserwacja wyniku kompilacji lub budowania | Build Output |

**Dane wykorzystywane przez AI.**
- Struktura projektu z Project Tree
- Treść kodu źródłowego z Code Editor
- Historia zmian z Git Panel
- Wynik budowania z Build Output
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Terminal | Developer korzysta z Terminala jako warstwy wykonawczej | Konfiguracyjne |
| Diagnostics | Developer korzysta z modułu Diagnostics przy analizie błędów wykrytych w toku pracy | Konfiguracyjne |
| Apps | Developer stanowi jeden z komponentów wykorzystywanych przy budowie kompletnych produktów w module Apps | Konfiguracyjne |

### 4.13. Diagnostics

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Analiza oraz usuwanie problemów technicznych |
| Przeznaczenie | Sytuacje, w których punktem wyjścia pracy jest nie tworzenie nowej funkcjonalności, lecz zrozumienie przyczyny istniejącego błędu lub spadku wydajności — na podstawie logów, komunikatów błędów i obserwowalnych objawów |
| Środowiska dostępności | CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Diagnostics Center, Logs Viewer, Errors Panel, Recommendations Panel |
| Funkcjonalności | Debugowanie; analiza logów; analiza błędów; diagnostyka wydajności |
| Powiązania konfigurowalne | Developer, Terminal |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku diagnozy; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Diagnostics Center | Agregacja spójnego obrazu stanu systemu | 1 |
| Logs Viewer | Przegląd logów | 2 |
| Errors Panel | Rejestr błędów | 2 |
| Recommendations Panel | Sugerowane kroki naprawcze | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Dostarczenie materiału źródłowego do analizy | Logs Viewer, Errors Panel |
| 2 | Agregacja materiału w spójny obraz stanu systemu | Diagnostics Center |
| 3 | Przedstawienie sugerowanych kroków naprawczych wypracowanych przez AI na podstawie zebranych danych | Recommendations Panel |

**Dane wykorzystywane przez AI.**
- Logi z Logs Viewer
- Rejestr błędów z Errors Panel
- Zagregowany stan systemu z Diagnostics Center
- Sugerowane kroki naprawcze z Recommendations Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Developer | Diagnostics współpracuje z modułem Developer przy wdrażaniu poprawek wynikających z analizy | Konfiguracyjne |
| Terminal | Diagnostics współpracuje z Terminalem przy odtwarzaniu i weryfikacji objawów błędu w rzeczywistym środowisku wykonawczym | Konfiguracyjne |

### 4.14. Apps

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Budowa kompletnych produktów cyfrowych |
| Przeznaczenie | Realizacja pełnych projektów aplikacyjnych — od architektury, przez warstwę frontend i backend, po wdrożenie — w jednym, spójnym procesie obejmującym wszystkie etapy powstawania produktu, a nie pojedyncze fragmenty kodu |
| Środowiska dostępności | WorkSpace, CodeStudio |
| Forma udostępnienia | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel |
| Funkcjonalności | Aplikacje webowe, mobilne, desktopowe; API; backend; frontend; repozytoria; dokumentacja techniczna |
| Powiązania konfigurowalne | Developer, Terminal, Design, MultitaskingAI |

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku budowy produktu; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Product Builder | Prowadzenie całościowego procesu budowy | 1 |
| Architecture Designer | Projektowanie architektury rozwiązania | 2 |
| Frontend Workspace | Praca nad warstwą frontendową | 2 |
| Backend Workspace | Praca nad warstwą backendową | 3 |
| Deployment Panel | Zarządzanie końcowym etapem udostępnienia produktu | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Zaprojektowanie architektury rozwiązania | Architecture Designer |
| 2 | Równoległa praca nad warstwą frontendową i backendową | Frontend Workspace, Backend Workspace |
| 3 | Zarządzanie końcowym etapem udostępnienia gotowego produktu | Deployment Panel |

Całościowy proces budowy prowadzi Product Builder, spinający powyższe etapy w jeden przebieg.

**Dane wykorzystywane przez AI.**
- Architektura rozwiązania z Architecture Designer
- Stan pracy nad frontendem z Frontend Workspace
- Stan pracy nad backendem z Backend Workspace
- Status wdrożenia z Deployment Panel
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Developer, Terminal | Apps integruje funkcjonalności modułów Developer i Terminal w ramach jednego, ustrukturyzowanego procesu budowy produktu | Konfiguracyjne |
| Design | Zasoby wizualne dla warstwy frontendowej mogą pochodzić z modułu Design, zgodnie z powiązaniem ustanowionym przez użytkownika | Konfiguracyjne |
| MultitaskingAI | Realizacja rozbudowanych projektów w module Apps jest naturalnym zastosowaniem środowiska MultitaskingAI (rozdz. 9.4 i 13 Koncepcji), w szczególności podziału na Executor 1 i Executor 2 przy równoległej pracy nad backendem i frontendem | Konfiguracyjne |

### 4.15. Agents

| Wymiar modułu | Charakterystyka |
|---|---|
| Cel | Tworzenie i zarządzanie własnymi agentami AI |
| Przeznaczenie | Użytkownicy chcący skonfigurować trwałe, wyspecjalizowane jednostki AI — z określoną tożsamością, zestawem umiejętności i uprawnień — wykorzystywane następnie jako wykonawcy zadań w innych częściach platformy, zamiast każdorazowego definiowania kontekstu od nowa |
| Środowiska dostępności | TalkIn, WorkSpace, CodeStudio |
| Forma udostępnienia | Komponent własny (Agent) konfigurowany w strefie 2 strony głównej oraz okno modułowe w bocznej nawigacji |
| Okna komunikacji operacyjnej | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) |
| Okna operacyjne | Chat Window, Execution Loop Window, Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center |
| Funkcjonalności | Wybór modelu bazowego; definiowanie tożsamości; konfiguracja instrukcji systemowych; dodawanie skilli, pluginów, konektorów; zarządzanie pamięcią; konfiguracja uprawnień |
| Powiązania konfigurowalne | Workspace, MultitaskingAI |

**Charakterystyka.** Agenci są komponentem platformowym działającym ponad wszystkimi środowiskami. Po utworzeniu mogą być wykorzystywani jako wykonawcy zadań w TalkIn, WorkSpace i CodeStudio, a po przypisaniu do roli — także w środowisku MultitaskingAI (rozdz. 9.4 i 13 Koncepcji).

**Okna operacyjne.**

Okna komunikacji operacyjnej modułu — Chat Window oraz Execution Loop Window — są elementami pierwszoplanowymi i zajmują lewe kolumny obszaru roboczego; pozostałe okna rozmieszczone są w kolumnach po ich prawej stronie. Kolumna „Warstwa” wskazuje warstwę widoczności zgodnie z rozdziałem 3a.

| Okno | Funkcja | Warstwa |
|---|---|---|
| Chat Window | Komunikacja z AI w toku konfiguracji agenta; centralny punkt pracy w module i podstawowy mechanizm sterowania procesami modułu (Użytkownik ↔ Wykonawca) | 1 |
| Execution Loop Window | Pętla wykonawcza modułu (Koordynator ↔ Wykonawca): dekompozycja zlecenia na zadania, kolejka i stan zadań, koordynacja, nadzór, orkiestracja i kontrola realizacji, sterowanie przebiegiem | 1 |
| Agent Builder | Tworzenie agenta | 1 |
| Model Configuration | Wybór modelu bazowego | 2 |
| Skills Manager | Dobór umiejętności agenta | 2 |
| Connectors Manager | Podłączenie integracji zewnętrznych | 3 |
| Permissions Center | Ustalenie zakresu uprawnień | 3 |

**Przebieg pracy modułu.**

| Krok | Działanie | Okno lub mechanizm |
|---|---|---|
| 1 | Utworzenie agenta | Agent Builder |
| 2 | Wybór modelu bazowego | Model Configuration |
| 3 | Dobór umiejętności | Skills Manager |
| 4 | Podłączenie integracji zewnętrznych | Connectors Manager |
| 5 | Ustalenie zakresu uprawnień przed udostępnieniem agenta do wykorzystania operacyjnego | Permissions Center |

**Dane wykorzystywane przez AI.**
- Konfiguracja modelu bazowego z Model Configuration
- Zestaw umiejętności ze Skills Manager
- Podłączone integracje z Connectors Manager
- Zakres uprawnień z Permissions Center
- Tożsamość i instrukcje systemowe agenta z Agent Builder
- Przebieg pętli wykonawczej, stan zadań i wyniki kontroli jakości z Execution Loop Window
- Polecenia, zatwierdzenia i przerwania wydane w Chat Window

**Powiązania konfigurowalne.**

| Moduł lub funkcja docelowa | Charakter powiązania | Typ |
|---|---|---|
| Workspace | Agenci skonfigurowani w tym module mogą być przypisywani do projektów w module Workspace przez Agent Manager | Konfiguracyjne |
| MultitaskingAI | Agenci mogą pełnić role wykonawcze (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) w środowisku MultitaskingAI, zgodnie z panelem orkiestracji opisanym w rozdziale 13.8 Koncepcji | Konfiguracyjne |

---

## 5. Zbiorcze zestawienie powiązań międzymodułowych

Rozdział zbiera w jednym miejscu wszystkie powiązania wymienione przy poszczególnych modułach w rozdziale 4. Zgodnie z zasadą jawności i konfigurowalności zależności (rozdz. 14, zasada 3 Koncepcji platformy) żadne z poniższych powiązań nie jest wbudowaną na stałe zależnością — każde jest możliwością do ustanowienia przez użytkownika z poziomu okna konfiguracji lub, tam gdzie dotyczy to modułu Automations, dodatkowo w ustawieniach okna tego modułu.

**Mapa powiązań międzymodułowych (diagram).** Poniższy schemat grupuje piętnaście modułów w trzy klastry funkcjonalne, wskazuje kierunek każdego powiązania konfigurowalnego oraz oddziela je od pokrewieństw funkcjonalnych.

```
LEGENDA
  ───►  powiązanie konfigurowalne (moduł źródłowy → moduł docelowy)
  ◄──►  powiązanie wzajemne (oba moduły wskazują na siebie)
  ═══   pokrewieństwo funkcjonalne (nie łączy kontekstu, historii ani pamięci)

ZBIEŻNOŚĆ KLASTRÓW
  Klaster treści      ─┐
  Klaster inżynierii  ─┼──►  środowisko MultitaskingAI (orkiestracja ról i kolejek, rozdz. 13)
  Klaster orkiestracji─┘
  łączniki między klastrami:  Design ──► Studio (treść) oraz Design ──► Apps (produkt)


[1] KLASTER TREŚCI, WIEDZY I KOMUNIKACJI                  środowiska: TalkIn, WorkSpace
    (węzły łączące: Studio i Research; Library jako centralne repozytorium)

    Studio      ◄──►  Translate      współdzielone operacje kontekstowe AI (praca dwujęzyczna)
    Studio      ◄──►  Research       redakcja wyników badawczych oraz dalsze przetwarzanie
    Studio      ───►  Library        repozytorium źródłowe i miejsce docelowe artefaktów
    Design      ───►  Studio         zasoby wizualne przy redagowaniu dokumentów
    Browser     ◄──►  Research       wspólne źródła oraz zasilanie pracy badawczej
    Browser     ───►  Library        trwałe przechowanie odnotowanych materiałów
    Research    ───►  Library        wykorzystanie dokumentów przechowywanych centralnie
    Research    ───►  Roundtable     wielomodelowa prezentacja wyników przed wnioskami
    Library     ───►  Studio · Research · Browser    rola centralnego repozytorium


[2] KLASTER INŻYNIERII I KODU                                    środowisko: CodeStudio
    (węzeł: Developer; Terminal jako warstwa wykonawcza)

    Terminal    ───►  Developer · Diagnostics    warstwa wykonawcza poleceń
    Developer   ◄──►  Terminal       wykorzystanie warstwy wykonawczej
    Developer   ◄──►  Diagnostics    analiza błędów oraz wdrażanie poprawek
    Diagnostics ───►  Terminal       odtwarzanie i weryfikacja objawów błędu
    Developer   ───►  Apps           komponent budowy kompletnych produktów
    Design      ───►  Apps           zasoby wizualne warstwy frontendu
    Apps        ───►  Developer · Terminal    integracja w jednym procesie budowy


[3] KLASTER ORKIESTRACJI                          zbieżność: MultitaskingAI
    (zbieżność ku środowisku MultitaskingAI — rozdz. 13 Koncepcji)

    Agents      ───►  Workspace         przypisanie agentów do projektów (Agent Manager)
    Workspace   ───►  Automations       powiązanie projektów z procesami automatycznymi
    Automations ───►  dowolny moduł     objęcie zadań dowolnego modułu procesem
    Agents      ───►  MultitaskingAI    role Executor 1, Executor 2, Coordinator, Executor 3 / Validator
    Automations ───►  MultitaskingAI    harmonogramy i kolejki pętli pracy ciągłej (rozdz. 13.6)
    Apps        ───►  MultitaskingAI    podział Executor 1 / Executor 2 (backend i frontend)


[4] POKREWIEŃSTWO FUNKCJONALNE (poza łączeniem kontekstu, historii i pamięci)

    Roundtable  ═══  MultitaskingAI    debata i konsensus wobec orkiestracji ról wykonawczych
    Assistant   ═══  Always On Display  moduł osadzony w TalkIn wobec nadrzędnej funkcji globalnej
```

**Tabela zbiorcza powiązań.** Poniższa tabela porządkuje wszystkie powiązania według modułu źródłowego.

| Moduł źródłowy | Moduł lub funkcja docelowa | Charakter powiązania |
|---|---|---|
| Studio | Translate | Współdzielenie operacji kontekstowych AI przy pracy dwujęzycznej |
| Studio | Library | Repozytorium źródłowe dokumentów wejściowych i miejsce docelowe artefaktów |
| Studio | Research | Dalsze przetwarzanie wyników redakcji przy budowie raportów |
| Workspace | Agents | Udostępnienie agentów jako wykonawców zadań projektu |
| Workspace | Library | Project Library jako odpowiednik Library ograniczony do projektu |
| Workspace | Automations | Powiązanie projektów z procesami automatycznymi |
| Automations | MultitaskingAI | Zaplecze harmonogramów i kolejek dla pętli pracy ciągłej (rozdz. 13.6 Koncepcji) |
| Automations | dowolny moduł platformy | Objęcie zadań z dowolnego modułu procesem automatycznym |
| Browser | Research | Zasilanie pracy badawczej zebranymi źródłami |
| Browser | Library | Trwałe przechowanie odnotowanych materiałów |
| Research | Browser | Wykorzystanie źródeł zebranych we wspólnym przeglądaniu |
| Research | Library | Wykorzystanie dokumentów przechowywanych centralnie |
| Research | Studio | Dalsza redakcja wyników pracy badawczej |
| Research | Roundtable | Wielomodelowa prezentacja wyników przed wnioskami końcowymi |
| Library | Studio, Research, Browser | Rola centralnego repozytorium materiałów |
| Translate | Studio | Wykorzystanie operacji kontekstowych AI i etap pośredni przy dokumentach wielojęzycznych |
| Design | Studio | Wykorzystanie zasobów wizualnych przy redagowaniu dokumentów |
| Design | Apps | Wykorzystanie zasobów wizualnych przy budowie interfejsu produktu |
| Terminal | Developer, Diagnostics | Warstwa wykonawcza poleceń w ramach środowiska CodeStudio |
| Developer | Terminal | Wykorzystanie warstwy wykonawczej |
| Developer | Diagnostics | Analiza błędów wykrytych w toku pracy |
| Developer | Apps | Komponent wykorzystywany przy budowie kompletnych produktów |
| Diagnostics | Developer | Wdrażanie poprawek wynikających z analizy |
| Diagnostics | Terminal | Odtwarzanie i weryfikacja objawów błędu |
| Apps | Developer, Terminal | Integracja funkcjonalności w jednym procesie budowy produktu |
| Apps | Design | Pozyskanie zasobów wizualnych warstwy frontendowej |
| Apps | MultitaskingAI | Realizacja rozbudowanych projektów przez podział Executor 1 / Executor 2 |
| Agents | Workspace | Przypisanie agentów do projektów przez Agent Manager |
| Agents | MultitaskingAI | Pełnienie ról wykonawczych (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) |

**Szablon powiązania konfigurowalnego.** Poniższy blok porządkuje redakcyjnie pola każdego powiązania wymienionego powyżej; wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna” (rozdz. 6 i rozdz. 14 Koncepcji).

```
POWIĄZANIE KONFIGUROWALNE (szablon redakcyjny)
  Moduł źródłowy:            < np. Studio >
  Moduł lub funkcja docelowa: < np. Translate >
  Na czym polega:            < np. współdzielenie operacji kontekstowych AI >
  Typ:                       konfiguracyjne  /  pokrewieństwo funkcjonalne
  Miejsce ustanowienia:      okno konfiguracji (rozdz. 6 Koncepcji)
                             [dla modułu Automations dodatkowo: ustawienia okna modułu]
  Stan domyślny:             nieaktywne — brak współdzielenia kontekstu, historii i pamięci
  Odwracalność:              tak — powiązanie jawne, w każdej chwili wyłączalne
```

Dwa moduły nie występują w tabeli jako źródła powiązań konfigurowalnych — Roundtable oraz Assistant. Ich relacje mają charakter pokrewieństwa funkcjonalnego, nie konfiguracyjnego, co zestawia poniższa tabela.

| Moduł | Element pokrewny | Charakter relacji |
|---|---|---|
| Roundtable | Środowisko MultitaskingAI (rozdz. 4.8) | Pokrewieństwo funkcjonalne — nie polega na łączeniu kontekstu, historii lub pamięci w rozumieniu rozdziału 6 Koncepcji platformy |
| Assistant | Funkcja globalna Always On Display (rozdz. 4.10) | Pokrewieństwo funkcjonalne — nie polega na łączeniu kontekstu, historii lub pamięci w rozumieniu rozdziału 6 Koncepcji platformy |

---

## 6. Zgodność z zasadami nadrzędnymi platformy

Specyfikacja modułów zawarta w rozdziale 4 pozostaje w pełnej zgodności z czterema zasadami nadrzędnymi funkcjonalności platformy, ustalonymi w rozdziale 14 Koncepcji platformy. Poniższa tabela wskazuje, jak każda z zasad jest realizowana.

| Zasada nadrzędna (rozdz. 14 Koncepcji) | Sposób realizacji w specyfikacji modułów |
|---|---|
| Pełna kompozycyjność | Przebieg pracy modułu opisany przy każdym module przedstawia typowy, domyślny przebieg pracy, a nie jedyną dopuszczalną sekwencję działań; użytkownik może komponować własne układy z dostępnych komponentów (prompt, zadanie, akcja, orkiestracja, subagenci) w każdym module, zgodnie z rozdziałem 10.2 dokumentu Model danych |
| Pełna konfigurowalność | Żadna funkcjonalność opisana przy module nie jest zaszyta na stałe w sposób niedostępny dla użytkownika; zakres danych wykorzystywanych przez AI — historia, pamięć, materiały z innych modułów — podlega punktom izolacji konfigurowalnym w oknie opisanym w rozdziale 6.4–6.6 Koncepcji; brak ustawienia oznacza dziedziczenie z poziomu szerszego, aż do poziomu globalnego stanowiącego warstwę bazową |
| Jawność i konfigurowalność zależności | Wszystkie powiązania wymienione w rozdziale 4 i zebrane zbiorczo w rozdziale 5 są jawnymi, nazwanymi możliwościami do ustanowienia przez użytkownika; żadne z nich nie jest domyślnie aktywne jako współdzielenie kontekstu, historii lub pamięci; każda zależność jest decyzją podejmowaną świadomie w oknie konfiguracji |
| Kierunek rozszerzania orkiestracji | Moduły Agents, Automations, Workspace i Apps współdziałają ze środowiskiem MultitaskingAI: agent pełni rolę wykonawcy, automatyka stanowi zaplecze harmonogramowe pętli pracy ciągłej, a złożone projekty budowane w module Apps korzystają z orkiestracji wielomodelowej opisanej w rozdziale 13 Koncepcji platformy. Pętlę wykonawczą każdego z tych procesów prezentuje Execution Loop Window |

Specyfikacja nie wprowadza żadnej blokady niemożliwej do wyłączenia z okna konfiguracji, a brak ustawienia w każdym z opisanych modułów oznacza wartość domyślną, a nie ograniczenie dostępnej funkcjonalności.

---

## 7. Kryteria odbioru

Opracowanie jest przyjmowane do budowy, gdy poniższe warunki są spełnione łącznie.

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Piętnaście modułów opisanych w rozdziale 4, po jednym podrozdziale na moduł | Podrozdziały 4.1–4.15 obecne, w kolejności rozdziału 11 Koncepcji platformy |
| Każdy moduł ma przypisany obszar kontraktu istniejący w `contract.json` | odczyt pola `nazwa` z listy `obszary` w `budowa/shared/contract.json` obejmuje wszystkie piętnaście wartości z kolumny „Obszar kontraktu” tabeli otwierającej rozdział 4 |
| Każdy moduł ma tabelę okien operacyjnych z warstwą widoczności | Tabela „Okno / Funkcja / Warstwa” obecna w każdym podrozdziale 4.1–4.15 |
| Każdy moduł ma opisany przebieg pracy | Nagłówek „**Przebieg pracy modułu.**” obecny w każdym podrozdziale 4.1–4.15; zero wystąpień zapisu „Workflow modułu” |
| Macierz dostępności modułów obejmuje wszystkie piętnaście modułów i trzy środowiska modułowe | Rozdział 3 zawiera piętnaście wierszy i kolumny TalkIn, WorkSpace, CodeStudio |
| Zbiorcze zestawienie powiązań międzymodułowych jest spójne z powiązaniami w rozdziale 4 | Każde powiązanie wymienione w podpunkcie „Powiązania konfigurowalne” modułu ma odpowiednik w rozdziale 5 |
| Formy wizualne zajmują co najmniej 30% wierszy dokumentu | Kontrola ręczna zgodnie z rozdziałem 3.1 [Standardu redakcyjnego i językowego](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| Objętość dokumentu nie mniejsza niż 85 000 znaków | liczba znaków pliku źródłowego (rozdz. 9 [Standardu redakcyjnego i językowego](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md)), licznik uruchomiony na `docs/specyfikacje/specyfikacja-modulow.md` |

---

*Koniec dokumentu. Specyfikacja modułów — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
