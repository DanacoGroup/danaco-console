*Dokument specyfikuje interfejs modułu Agents Danaco Console: okna, makiety, elementy, warstwy widoczności, stany i punkty sterowania.*

# Moduł Agents — dokumentacja projektowa

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
| **Tytuł** | Moduł Agents — pełnozakresowa dokumentacja projektowa |
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować): komplet okien operacyjnych, katalog elementów każdego okna, diagram cyklu tworzenia eksperta, stany oraz punkty sterowania z okna konfiguracji |
| **Środowiska dostępności** | TalkIn, WorkSpace, CodeStudio; moduł osiągalny także ze strefy 2 strony głównej |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł tworzy komponent własny rodzaju Agent (ekspert) |
| **Data opracowania** | 2026-08-06 |
| **Dokument nadrzędny** | `specyfikacje/specyfikacja-agentow.md` — model agentów, role i pętla wykonawcza; niniejszy dokument opisuje wyłącznie interfejs modułu Agents |
| **Źródła** | Koncepcja platformy (rozdz. 1–2, 6–7, 10–11.15, 13–14, Załącznik D.1) · Specyfikacja modułów (rozdz. 4.15) · Specyfikacja okien operacyjnych (rozdz. 6.15, 7, 8) · Strona główna i nawigacja (rozdz. 3.3, sekcja Role) · System wizualny (rozdz. 1–14, Załączniki A–C) · Model konfiguracji (rozdz. 3–7) · `specyfikacje/specyfikacja-agentow.md` (całość) · `srodowiska/multitaskingai.md` · Izolacja i zależności · Rozszerzenia (całość) · Bezpieczeństwo i uwierzytelnianie (rozdz. 2) |
| **Zakres** | Sześć okien operacyjnych modułu Agents — Chat Window, Execution Loop Window, Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center — przy czym Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca) są kolumnami wspólnymi wszystkim pozostałym widokom modułu |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); zero blokad w interfejsie; klucze i integracje jawne; izolacja i zakres możliwości są ustawieniami konfiguracyjnymi Operatora, nie wymogiem; domyślne zachowanie modułu = wykonanie |

Moduł Agents udostępnia sześć okien operacyjnych (Koncepcja platformy, rozdz. 11.15; Specyfikacja okien operacyjnych, rozdz. 6.15): **Chat Window**, **Agent Builder**, **Model Configuration**, **Skills Manager**, **Connectors Manager**, **Permissions Center**. Niniejszy dokument projektuje wszystkie sześć na poziomie gotowym do wykonania — makieta, katalog elementów, stany — z Permissions Center rozumianym wyłącznie jako **konfiguracja zakresu możliwości eksperta**, nigdy jako mechanizm blokujący (rozdz. 9 niniejszego dokumentu; podstawa: Specyfikacja agentów, rozdz. 6).

**Konwencja terminologiczna.** Specyfikacja platformy nazywa wytwór modułu Agents „agentem” (Koncepcja platformy, rozdz. 2.3 i 11.15; encja Agent, Model danych rozdz. 8.2). Niniejszy dokument, pisany z perspektywy interfejsu, posługuje się zamiennie określeniem **ekspert** — Agents jest fabryką, która zamienia surowy model bazowy w nazwaną, skonfigurowaną jednostkę AI. Oba terminy oznaczają dokładnie ten sam komponent własny; nazwy okien, ról i mechanizmów pozostają dosłownie takie, jak w dokumentach źródłowych.

**Agent jako Wykonawca.** W architekturze platformy agent (ekspert) występuje w roli **Wykonawcy** — AI, agenta lub systemu wykonawczego realizującego zadania. Moduł Agents jest miejscem, w którym Wykonawca powstaje i jest konfigurowany. Wykonawca komunikuje się z Użytkownikiem przez **Chat Window** (kanał Użytkownik ↔ Wykonawca) oraz z **Koordynatorem** — komponentem orkiestrującym platformy — przez **Execution Loop Window** (kanał Koordynator ↔ Wykonawca). Relacje te opisuje rozdz. 1.4, a udział Wykonawcy w pętli wykonawczej — rozdz. 10.4.

**Rozgraniczenie wobec dokumentu systemowego.** Model agentów — definicja agenta jako komponentu własnego, siedem komponentów definicji, role agenta w środowiskach, cykl życia agenta, pełny przebieg pętli wykonawczej, model danych i kontrakty komunikacji agenta ma odrębny dokument systemowy. Niniejszy dokument opisuje wyłącznie **interfejs modułu Agents**: okna, ich makiety, elementy, warstwy widoczności, stany oraz punkty sterowania z okna konfiguracji. Ustalenia modelowe nie są tu powtarzane, lecz przywoływane odesłaniem do właściwego rozdziału dokumentu systemowego.

**Warstwy widoczności.** Każdy element interfejsu opisany w niniejszym dokumencie należy do dokładnie jednej z czterech warstw widoczności platformy. Katalogi elementów rozdz. 2, 5–9 podają warstwę i sposób wywołania każdego elementu; makiety rysowane są w stanie spoczynku interfejsu — widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3. Zasadę i zestawienie warstw dla modułu Agents zawiera rozdz. 2.5.

---

## Spis treści

1. [Agents jako fabryka ekspertów — pozycja w platformie](#1-agents-jako-fabryka-ekspertów--pozycja-w-platformie)
   - [1.4. Użytkownik, Koordynator, Wykonawca — relacje i dwa kanały komunikacji](#14-użytkownik-koordynator-wykonawca--relacje-i-dwa-kanały-komunikacji)
2. [Mapa okien modułu Agents](#2-mapa-okien-modułu-agents)
   - [2.5. Warstwy widoczności w oknie](#25-warstwy-widoczności-w-oknie)
3. [Diagram cyklu tworzenia eksperta](#3-diagram-cyklu-tworzenia-eksperta)
4. [Stany](#4-stany)
5. [Agent Builder — okno nadrzędne](#5-agent-builder--okno-nadrzędne)
6. [Model Configuration](#6-model-configuration)
7. [Skills Manager](#7-skills-manager)
8. [Connectors Manager](#8-connectors-manager)
9. [Permissions Center](#9-permissions-center)
10. [Od eksperta do Wykonawcy — zastosowanie operacyjne](#10-od-eksperta-do-wykonawcy--zastosowanie-operacyjne)
    - [10.4. Udział Wykonawców w pętli wykonawczej](#104-udział-wykonawców-w-pętli-wykonawczej)
11. [Katalog komponentów i ikon systemu wizualnego w module Agents](#11-katalog-komponentów-i-ikon-systemu-wizualnego-w-module-agents)
12. [Słowniczek pojęć](#12-słowniczek-pojęć)
13. [Punkty sterowania z okna konfiguracji](#13-punkty-sterowania-z-okna-konfiguracji)
- [Załącznik A. Szablon definicji eksperta jako formularz UI](#załącznik-a-szablon-definicji-eksperta-jako-formularz-ui)

---

## 1. Agents jako fabryka ekspertów — pozycja w platformie

### 1.1. Łańcuch wartości

```
Diagram — łańcuch wartości modułu Agents

  MODEL BAZOWY                  surowa zdolność obliczeniowa, bez tożsamości
  (inteligencja)                 wybór i kanał połączenia → Model Configuration
        │
        │   + tożsamość, instrukcje systemowe, pamięć  ─────  Agent Builder
        │   + umiejętności                              ─────  Skills Manager
        │   + pluginy, konektory, serwery MCP            ─────  Connectors Manager
        │   + zakres możliwości                          ─────  Permissions Center
        ▼
     EKSPERT (AGENT)             nazwana, skonfigurowana jednostka AI —
     = WYKONAWCA                  komponent własny, zasób Operatora
                                  w pamięci aplikacji
        │
        ▼   wykorzystanie operacyjne (rozdz. 10 niniejszego dokumentu)
   ┌────────────┬───────────────────┬────────────────────┬──────────────────┐
   ▼            ▼                   ▼                     ▼
 TalkIn       WorkSpace           CodeStudio            MultitaskingAI
 wykonawca    wykonawca /         wykonawca             rola: Executor 1/2 ·
 doraźny      przypisanie do      doraźny               Coordinator ·
              projektu (Agent                            Executor 3/Validator
              Manager, moduł                              (sekcja Role)
              Workspace)
                                                                 │
                                                                 ▼
                                                          ZESPÓŁ WYKONAWCÓW
                                                   Koordynator prowadzi pętlę
                                                wykonawczą, silnik kolejek, integracja
                                                z Automations → praca ciągła 24/7/365
```

Podstawa: Koncepcja platformy, rozdz. 1 (pięć elementów łańcucha organizacji pracy), rozdz. 9.15, rozdz. 11; Specyfikacja agentów, rozdz. 1 i 7–8.

### 1.2. Agent jako komponent własny — pozycja wśród czterech rodzajów

Pozycję agenta wśród czterech rodzajów komponentów własnych oraz jego status komponentu platformowego ponad środowiskami ustala `specyfikacje/specyfikacja-agentow.md`, rozdz. 1.1–1.2. Dla interfejsu wynika stąd jedna konsekwencja projektowa: agent jest jedynym komponentem własnym, którego zastosowanie operacyjne nie jest ograniczone do jednego kontekstu pracy, więc moduł Agents musi być osiągalny z dwóch punktów wejścia prowadzących do tej samej definicji.

| Punkt wejścia | Kontekst |
|---|---|
| Strona główna, strefa 2 — kafel „Agents” („Skonfiguruj agenta”) | Tworzenie i zarządzanie ekspertami jako takimi |
| Boczna nawigacja modułów środowiska (TalkIn, WorkSpace, CodeStudio) | Tworzenie lub edycja eksperta w toku bieżącej pracy nad zadaniem |

### 1.3. Gdzie ekspert działa — punkty styku interfejsu

Zakres i charakter wykorzystania agenta w każdym środowisku ustala `specyfikacje/specyfikacja-agentow.md`, rozdz. 7–9. Poniższa tabela wskazuje wyłącznie miejsce interfejsu, z którego przypisanie następuje.

| Środowisko | Miejsce interfejsu, z którego następuje przypisanie | Rozdział niniejszego dokumentu |
|---|---|---|
| TalkIn | Selektor komponentu własnego w sesji modułu (`Agent ▼` w pasku kontekstu) | 10.1 |
| WorkSpace | Selektor `Agent ▼` w sesji modułu oraz okno Agent Manager modułu Workspace | 10.1, 10.2 |
| CodeStudio | Selektor komponentu własnego w sesji modułu (`Agent ▼` w pasku kontekstu) | 10.1 |
| MultitaskingAI | Sekcja Role panelu orkiestracji (`srodowiska/multitaskingai.md`, rozdz. 6.2) | 10.3 |

### 1.4. Użytkownik, Koordynator, Wykonawca — relacje i dwa kanały komunikacji

Moduł Agents zarządza agentami, czyli **Wykonawcami** platformy. Wykonawca uczestniczy w dwóch kanałach komunikacji operacyjnej, obecnych w każdym oknie modułu jako elementy pierwszoplanowe architektury.

| Rola | Czym jest | Co robi |
|---|---|---|
| **Użytkownik** | Operator platformy | Zleca zadania, zatwierdza i przerywa działania, wyjaśnia kontekst |
| **Koordynator** | Komponent orkiestrujący platformy | Dekomponuje zlecenie na zadania, przydziela je Wykonawcom, nadzoruje realizację, prowadzi kontrolę jakości i decyduje o ponowieniach |
| **Wykonawca** | AI, agent lub system wykonawczy — ekspert utworzony w module Agents | Realizuje przydzielone zadania, raportuje stan i wynik |

**Kanał pierwszy — Chat Window (Użytkownik ↔ Wykonawca).** Główne okno komunikacji między Użytkownikiem a Wykonawcą. Stanowi centralny punkt pracy Użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. W module Agents obsługuje budowę, testowanie i wybór Wykonawcy. Zajmuje lewą kolumnę obszaru roboczego, stałą, o pełnej wysokości — w tym samym miejscu układu we wszystkich pięciu oknach modułu.

**Kanał drugi — Execution Loop Window (Koordynator ↔ Wykonawca).** Okno pętli wykonawczej prezentujące komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Zawiera bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie, korektę zlecenia. Otwierane jest jako kolumna sąsiadująca z Chat Window.

Definicja Wykonawcy zapisana w module Agents jest tym, co Koordynator otrzymuje przy przydzieleniu zadania: tożsamość, instrukcje systemowe, model bazowy i kanał, umiejętności, rozszerzenia, pamięć oraz zakres możliwości. Szczegóły udziału Wykonawcy w pętli wykonawczej opisuje rozdz. 10.4.

---

## 2. Mapa okien modułu Agents

### 2.1. Komplet okien operacyjnych modułu i ich relacje

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca (wspólne wszystkim modułom) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji przy każdym wejściu do modułu | Centralny punkt pracy: polecenia w języku naturalnym, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań, budowa i testowanie Wykonawcy w rozmowie |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window | 1 | Widoczne bez interakcji; zwinięcie do znacznika kolumny przy wąskim obszarze roboczym, rozwinięcie jednym kliknięciem | Pętla wykonawcza zlecenia: dekompozycja na zadania, kolejka i stan zadań, komunikaty sterujące, kontrola jakości, sterowanie przebiegiem |
| 3 | Agent Builder | Okno edycyjne / kreator — okno nadrzędne modułu | Prawa kolumna, dominująca | 1 | Widoczne bez interakcji jako aktywne okno wiodące modułu | Biblioteka ekspertów i edytor definicji: tożsamość, instrukcje systemowe, pamięć, historia wersji; scala wynik czterech okien konfiguracyjnych w jeden zapisywalny komponent własny |
| 4 | Model Configuration | Panel konfiguracyjny | Zakładka edytora w kolumnie dominującej | 2 | Zakładka `Model` paska zakładek edytora, polecenie języka naturalnego w Chat Window | Wybór modelu bazowego i kanału połączenia (API, CLI, SSH, HTTP) |
| 5 | Skills Manager | Rejestr / panel narzędziowy | Zakładka edytora w kolumnie dominującej | 2 | Zakładka `Skille` paska zakładek edytora, polecenie języka naturalnego w Chat Window | Podłączanie umiejętności do definicji eksperta |
| 6 | Connectors Manager | Rejestr / panel narzędziowy | Zakładka edytora w kolumnie dominującej | 2 | Zakładka `Konektory` paska zakładek edytora, polecenie języka naturalnego w Chat Window | Podłączanie pluginów, konektorów i serwerów MCP do definicji eksperta |
| 7 | Permissions Center | Panel konfiguracyjny zakresu możliwości | Zakładka edytora w kolumnie dominującej | 4 | Zakładka `Uprawnienia` paska zakładek edytora, tryb administracyjny, konfiguracja roli, polecenie języka naturalnego w Chat Window | Konfiguracja zakresu możliwości eksperta w czterech grupach (rozdz. 9); nigdy mechanizm blokujący |
| — | Panel Historia wersji | Panel pomocniczy | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Przycisk `Historia wersji` w nagłówku edytora, menu kebab (⋮) karty eksperta | Przegląd i przywracanie wcześniejszych wersji definicji eksperta |

Chat Window i Execution Loop Window są kolumnami wspólnymi: występują w tym samym miejscu układu niezależnie od tego, która zakładka edytora jest otwarta (rozdz. 2.4). Cztery okna konfiguracyjne — Model Configuration, Skills Manager, Connectors Manager, Permissions Center — działają jako zakładki edytora okna nadrzędnego Agent Builder i zapisują swój wynik w jednej definicji eksperta.

```
Schemat — okna modułu Agents i ich wzajemne relacje

                       CHAT WINDOW  +  EXECUTION LOOP WINDOW
              lewa kolumna: Użytkownik ↔ Wykonawca · kolumna sąsiadująca:
              Koordynator ↔ Wykonawca — obecne w każdym z pięciu poniższych
                    widoków jako stałe kolumny lewej strony układu
                                      │
                               AGENT BUILDER                         ◄── okno nadrzędne
              widok „Biblioteka ekspertów” · widok „Edytor” (tożsamość,
              instrukcje systemowe, pamięć, historia wersji) — scala wynik
              czterech okien konfiguracyjnych w jeden zapisywalny komponent własny
        ┌────────────────┬───────────────────┬────────────────────┬───────────────┐
        ▼                ▼                    ▼                     ▼
   MODEL                SKILLS              CONNECTORS            PERMISSIONS
   CONFIGURATION         MANAGER             MANAGER                CENTER
   model bazowy          umiejętności        pluginy ·              zakres możliwości
   + kanał połączenia     (skille)            konektory ·            eksperta
   (API/CLI/SSH/HTTP)                          serwery MCP           (4 grupy — rozdz. 9)
        │                     │                     │                      │
        │                     └──────────┬──────────┘                      │
        │              rejestr rozszerzeń (okno konfiguracji platformy —    │
        │              poza modułem Agents; podłączenie do TEGO eksperta    │
        │              rozstrzyga się tutaj, rozdz. 7–8)                    │
        │                                                                   │
        └───────────────────────────────┬───────────────────────────────────┘
                                         ▼
                        ZAPIS W AGENT BUILDER → NOWA WERSJA EKSPERTA
```

Podstawa: Specyfikacja agentów, rozdz. 4.1; Rozszerzenia, rozdz. 6.

### 2.2. Legenda ikon okien (System wizualny, 47-ikonowy zestaw — Załącznik C)

| Okno | Ikona | Uzasadnienie doboru |
|---|---|---|
| Agent Builder | `uzytkownik` | Tworzenie i identyfikacja tożsamości jednostki AI (ikona już przypisana do „Awatar / profil / role”) |
| Model Configuration | `uruchom` | Uruchomienie połączenia z modelem przez wybrany kanał |
| Skills Manager | `gwiazdka` | Umiejętności jako wyróżnione, przypisane kwalifikacje eksperta |
| Connectors Manager | `link-zewnetrzny` | Integracja zewnętrzna — ikona już przypisana wprost do tego znaczenia w Załączniku C |
| Permissions Center | `tarcza` | Załącznik C Systemu wizualnego przypisuje tę ikonę dosłownie: „Security Auditor, Permissions Center” |

Żadna z pięciu ikon nie wykracza poza istniejący zestaw 47 ikon (System wizualny, rozdz. 7.2 i Załącznik C) — zgodnie z zasadą przyjętą dla macierzy izolacji (rozdz. 12.2 Systemu wizualnego).

### 2.3. Makieta powłoki modułu

```
Makieta — powłoka modułu Agents w środowisku (TalkIn / WorkSpace / CodeStudio), stan spoczynku

 ══════════════════════════════════════════════════════════════════════════════
  .dn-pasek   Danaco Console · [nazwa środowiska]   [ szukaj ]   ⚙ 👤
 ══════════════════════════════════════════════════════════════════════════════
  BOCZNA      │ CHAT WINDOW          │ OBSZAR ROBOCZY MODUŁU AGENTS │ PANEL
  NAWIGACJA   │ Użytkownik ↔         │                              │ POMOCNICZY
  MODUŁÓW     │ Wykonawca            │ 🧑 Agent Builder             │ (rozszerzenie
              │                      │ [ Biblioteka ▼ ]  [＋ Nowy ] │  boczne,
  Studio      │ ┌──────────────────┐ │ ─────────────────────────── │  zwinięty)
  Workspace   │ │ strumień         │ │ Zakładki edytora:            │
  Browser     │ │ rozmowy          │ │ 🧑Tożsamość │🚀Model │⭐Skille│    ⋮
  Research    │ └──────────────────┘ │ 🔗Konektory │🛡Uprawnienia   │
  ► Agents    │ [ Napisz polecenie ] │ ─────────────────────────── │
  …           │                      │ OBSZAR ZAWARTOŚCI ZAKŁADKI   │
              │ ──────────────────── │ (formularz / rejestr /       │
              │ EXECUTION LOOP       │  macierz — wg zakładki)      │
              │ WINDOW               │                              │
              │ Koordynator ↔        │ ┌──────────────────────────┐ │
              │ Wykonawca            │ │ PODSUMOWANIE DEFINICJI   │ │
              │ zlecenie · kolejka   │ │ (.dn-karta boczna)       │ │
              │ zadań · kontrola     │ └──────────────────────────┘ │
              │ jakości · sterowanie │                              │
 ══════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu: [Danaco Console] [Ubuntu] [Agent ▼] [Model ▼] [Ultra]
 ══════════════════════════════════════════════════════════════════════════════
```

Układ jest wyłącznie pionowy — kolumny sąsiadują poziomo, regulacji podlega wyłącznie ich szerokość. Chat Window zajmuje lewą kolumnę, stałą, o pełnej wysokości obszaru roboczego; Execution Loop Window otwiera się jako kolumna sąsiadująca; obszar roboczy modułu Agents jest kolumną dominującą; panele pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego.

Anatomia zgodna z wzorcem wspólnym każdego okna operacyjnego platformy (Specyfikacja okien operacyjnych, rozdz. 2.4): nagłówek, obszar zawartości, elementy konfiguracji z objaśnieniem `[?]`, trwały stan procesu sesji.

### 2.4. Chat Window i Execution Loop Window — dwa kanały komunikacji operacyjnej

Chat Window i Execution Loop Window nie są odrębnymi, samodzielnie nawigowalnymi oknami modułu Agents — są stałymi kolumnami układu, obecnymi po lewej stronie każdego z pozostałych pięciu okien (Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center), zgodnie ze wzorcem wspólnym platformy (Specyfikacja okien operacyjnych, rozdz. 5). Chat Window prowadzi kanał Użytkownik ↔ Wykonawca, Execution Loop Window — kanał Koordynator ↔ Wykonawca (rozdz. 1.4). Niniejsza sekcja domyka ich projekt na poziomie gotowym do wykonania w kontekście modułu Agents, uzupełniając odwołania rozproszone w katalogach elementów rozdz. 5–9.

#### 2.4.1. Makieta

```
Makieta — Chat Window i Execution Loop Window jako kolumny lewej strony układu
(stan spoczynku; identyczne w każdym z pięciu okien modułu Agents)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW            │ ⟳ EXECUTION LOOP WINDOW      │ OBSZAR ROBOCZY
  Użytkownik ↔ Wykonawca    │ Koordynator ↔ Wykonawca      │ MODUŁU AGENTS
 ───────────────────────────┼──────────────────────────────┤
  Użytkownik: „dodaj temu   │ Zlecenie: rozbudowa definicji│ (zawartość
  ekspertowi umiejętność    │ Wykonawcy „Agent Redaktor”   │  wg rozdz.
  streszczeń”               │ ──────────────────────────── │  5–9)
                            │ Zadania:                     │
  Wykonawca: znaleziono     │  ▣ 1 wyszukanie umiejętności │
  „Streszczenia wielo-      │  ⟳ 2 przypisanie do definicji│
  języczne” w rejestrze —   │  ○ 3 kontrola jakości        │
  przypisać teraz?          │  ○ 4 zapis nowej wersji      │
  [ Przypisz ] [ Otwórz ⋮ ] │ ──────────────────────────── │
                            │ Kontrola jakości: oczekuje   │
                            │ Ponowienia: 0                │
  ──────────────────────────│ Przebieg: ▮▮▮▯▯ 2/4          │
  [ Napisz polecenie…     ] │ [ Sterowanie ▼ ]             │
  [ Agent ▼ ] [ Wyślij ]    │                              │
 ══════════════════════════════════════════════════════════════════════════════
```

Obie kolumny sąsiadują poziomo z obszarem roboczym modułu. Regulacji podlega wyłącznie ich szerokość. Element `[ Agent ▼ ]` jest zwiniętym menu progresywnym — lista Wykonawców rozwija się dopiero po kliknięciu. Element `[ Sterowanie ▼ ]` grupuje akcje przebiegu pętli (wstrzymanie, wznowienie, przerwanie, korekta zlecenia).

#### 2.4.2. Katalog elementów — Chat Window

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Nagłówek kolumny | Etykieta tekstowa z ikoną czatu i podpisem „Użytkownik ↔ Wykonawca” | Identyfikuje kanał komunikacji | 1 | Widoczny bez interakcji | Mały, stała wysokość | Domyślny | — (statyczny) | Górna krawędź lewej kolumny, w każdym z pięciu okien |
| Strumień wiadomości | Lista wypowiedzi Użytkownika i odpowiedzi Wykonawcy | Kontekst konwersacyjny wsparcia budowy i testowania Wykonawcy | 1 | Widoczny bez interakcji | Środkowa część kolumny, przewijalna | Pusty (podpowiedź startowa) · strumień na żywo (WebSocket) · błąd połączenia | Odpowiedzi zawierają skróty akcji osadzone w treści prowadzące wprost do właściwej zakładki | Środek lewej kolumny |
| Akcje kontekstowe w odpowiedzi | `.dn-btn--zarys` małe, osadzone w treści odpowiedzi Wykonawcy | Skrót do wykonania sugerowanej czynności bez opuszczania rozmowy | 1 | Widoczne bez interakcji | Małe przyciski inline | Domyślny · hover | Wykonuje akcję i potwierdza wynik w strumieniu | Wewnątrz odpowiedzi Wykonawcy |
| Zestaw dalszych akcji odpowiedzi | Menu kebab (⋮) przy akcjach kontekstowych | Grupuje warianty operacji wynikające z odpowiedzi | 3 | Menu kebab (⋮) | Zwinięty wyzwalacz | Zwinięty (domyślny) · rozwinięty | Rozwija pełną listę wariantów, po wyborze zwija się samoczynnie | Przy akcjach kontekstowych |
| Pole wpisu | `.dn-input` pełnej szerokości kolumny | Wpisanie polecenia w języku naturalnym | 1 | Widoczne bez interakcji | Duże pole tekstowe, dół kolumny | Domyślny · fokus · wypełnione | Enter lub „Wyślij” wysyła polecenie | Dolna część lewej kolumny |
| Selektor Wykonawcy `Agent ▼` | Menu progresywne — jeden element zwinięty | Wybór Wykonawcy obsługującego bieżącą rozmowę | 2 | Kliknięcie znacznika `Agent ▼`; po wyborze element zwija się samoczynnie | Mały znacznik kontekstowy przy polu wpisu | Zwinięty (domyślny) · rozwinięty (lista Wykonawców) · wybrany (nazwa w znaczniku) | Wybór ustawia Wykonawcę rozmowy i zwija listę | Przy polu wpisu |
| Przycisk „Wyślij” | `.dn-btn--zloty` mały | Wysyła treść pola wpisu | 1 | Widoczny bez interakcji | Mały, przy polu wpisu | Zawsze w pełni klikalny — pole puste przy kliknięciu sygnalizuje komunikatem „wpisz polecenie”, kontrolka nie jest blokowana | Wysyła treść i czyści pole; strumień aktualizuje się na żywo | Obok pola wpisu |
| Tryb administracyjny rozmowy | Zestaw poleceń diagnostycznych i administracyjnych kanału | Operacje eksperckie na kanale komunikacji | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji lub tryb administracyjny | Bez reprezentacji graficznej w stanie spoczynku | Ukryty dla użytkownika podstawowego | Wywołanie otwiera właściwe narzędzie w kolumnie roboczej | Chat Window |

#### 2.4.3. Katalog elementów — Execution Loop Window

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Nagłówek kolumny | Etykieta z podpisem „Koordynator ↔ Wykonawca” | Identyfikuje kanał pętli wykonawczej | 1 | Widoczny bez interakcji | Mały, stała wysokość | Domyślny | — | Górna krawędź kolumny sąsiadującej |
| Blok „Zlecenie i dekompozycja” | Nazwa bieżącego zlecenia oraz jego rozkład na zadania dokonany przez Koordynatora | Pokazuje, co i jak zostało rozłożone na zadania | 1 | Widoczny bez interakcji | Średni blok, góra kolumny | Brak zlecenia · zlecenie w toku · zlecenie zakończone | Kliknięcie zadania odsłania jego szczegóły w kolumnie | Górna część kolumny |
| Kolejka i stan zadań | Lista zadań ze znacznikami stanu (oczekuje, w toku, zakończone, ponawiane) | Prezentuje kolejkę i postęp realizacji | 1 | Widoczna bez interakcji | Lista, środek kolumny | Pusta · w toku · zakończona | Kliknięcie pozycji rozwija komunikaty sterujące tego zadania | Środek kolumny |
| Strumień komunikatów sterujących | Wymiana komunikatów Koordynator ↔ Wykonawca | Prezentuje przydzielenia, raporty stanu i decyzje sterujące | 1 | Widoczny bez interakcji | Przewijalny blok | Na żywo · bezczynny · błąd połączenia | — | Środek kolumny |
| Wynik kontroli jakości | Blok z rezultatem sprawdzenia zadania przez Koordynatora | Informuje o przyjęciu wyniku lub skierowaniu do ponowienia | 1 | Widoczny bez interakcji | Mały blok z plakietką | Oczekuje · przyjęty · skierowany do ponowienia | Kliknięcie odsłania uzasadnienie decyzji | Kolumna pętli wykonawczej |
| Licznik ponowień | Tekst „Ponowienia: N” | Pokazuje liczbę powtórzonych prób realizacji zadania | 1 | Widoczny bez interakcji | Bardzo mała waga | Domyślny | Najechanie pokazuje listę ponowień w `.dn-tooltip` | Kolumna pętli wykonawczej |
| Wskaźnik przebiegu pętli | Pasek postępu z licznikiem zadań | Pokazuje zaawansowanie zlecenia | 1 | Widoczny bez interakcji | Mały pasek | Domyślny · wstrzymany · przerwany | — | Dolna część kolumny |
| Grupa „Sterowanie ▼” | Element zbiorczy grupujący akcje: wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Sterowanie przebiegiem pętli wykonawczej | 3 | Kliknięcie elementu `Sterowanie ▼`; rozwinięcie zawiera pełną listę akcji | Jeden element zbiorczy zamiast zestawu przycisków | Zwinięty (domyślny) · rozwinięty | Wybór akcji wykonuje ją natychmiast i zwija listę | Dół kolumny pętli wykonawczej |
| Parametry pętli wykonawczej | Ustawienia liczby ponowień, progu kontroli jakości i trybu przydziału zadań | Dostrojenie zachowania pętli dla bieżącego zlecenia | 2 | Znacznik kontekstowy w kolumnie; po użyciu zwija się samoczynnie | Panel popover | Zwinięty · rozwinięty | Zmiana obowiązuje od najbliższego zadania | Kolumna pętli wykonawczej |
| Diagnostyka niskiego poziomu pętli | Podgląd surowych komunikatów sterujących i śladu wykonania | Rozpoznanie zachowania pętli na poziomie technicznym | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji lub tryb administracyjny | Bez reprezentacji graficznej w stanie spoczynku | Ukryta dla użytkownika podstawowego | Otwiera podgląd jako rozszerzenie boczne | Kolumna pętli wykonawczej |

#### 2.4.4. Stany

| Stan | Wyzwalacz | Co widać |
|---|---|---|
| Kolumna zwinięta do znacznika | Wąski obszar roboczy | Znacznik kolumny z ikoną kanału; rozwinięcie jednym kliknięciem |
| Kolumna rozwinięta | Kliknięcie znacznika kolumny lub pola wpisu | Pełny strumień i pole wpisu widoczne |
| Strumień na żywo | Aktywne połączenie WebSocket | Treść pojawia się przyrostowo; `.dn-spinner` w trakcie oczekiwania |
| Błąd połączenia | Utrata WebSocket | `.dn-alert--blad` w kolumnie; pole wpisu pozostaje aktywne — treść wysłana w tym stanie jest kolejkowana do ponowienia po przywróceniu połączenia |
| Pętla wstrzymana | „Wstrzymaj” z grupy `Sterowanie ▼` | Kolejka zadań zamrożona, wskaźnik przebiegu w stanie wstrzymanym, Chat Window pozostaje w pełni czynny |

### 2.5. Warstwy widoczności w oknie

Interfejs modułu Agents ujawnia funkcjonalność stopniowo: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Liczba Wykonawców, umiejętności, rozszerzeń i ustawień zakresu nie wpływa na postrzeganą prostotę okna — złożoność istnieje w architekturze i pozostaje niewidoczna do chwili wystąpienia potrzeby użycia.

| Warstwa | Nazwa | Co należy do niej w module Agents | Sposób wywołania |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, Execution Loop Window, nagłówek okna, pasek zakładek edytora, obszar zawartości bieżącej zakładki, panel podsumowania definicji, wskaźniki stanu i przebiegu, stopka z „Zapisz eksperta” | Widoczne bez interakcji |
| 2 | Widoczna na żądanie | Selektor Wykonawcy `Agent ▼`, wybór modelu bazowego, wybór kanału modelu, filtry Biblioteki i rejestrów, parametry pętli wykonawczej | Znacznik kontekstowy, przycisk, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Menu akcji karty Wykonawcy (`⋮`), grupa `Operacje ▼` rejestrów, grupa `Sterowanie ▼` pętli, panele szczegółów umiejętności i rozszerzeń, panel Historia wersji, selektor grup zakresu | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Definiowanie ról Wykonawców, konfiguracja uprawnień i zakresu możliwości (Permissions Center), limity zasobów (liczba podagentów, macierz izolacji technicznej), tryb administracyjny, diagnostyka niskiego poziomu pętli wykonawczej, podgląd polityki efektywnej | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Mechanizmy ukrywania funkcjonalności zastosowane w module Agents:

- **Menu progresywne.** Lista Wykonawców prezentowana jest jako jeden element zwinięty `Agent ▼`; pozycje rozwijają się po kliknięciu.
- **Panele wysuwane.** Historia wersji, szczegóły umiejętności i konfiguracja instancji rozszerzenia otwierają się jako rozszerzenia boczne i po zamknięciu znikają całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Zestawy akcji rejestrów i pętli wykonawczej występują jako jeden element zbiorczy (`Operacje ▼`, `Sterowanie ▼`).
- **Znaczniki kontekstowe.** Środowisko, model i Wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Agent ▼] [Model ▼] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

Każda ukryta funkcja modułu Agents jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego.

---

## 3. Diagram cyklu tworzenia eksperta

```
Diagram — cykl tworzenia i życia eksperta

  START            STREFA 2 STRONY GŁÓWNEJ, kafel „Agents”
                    albo BOCZNA NAWIGACJA modułu Agents w środowisku
                                    │
                                    ▼
                        AGENT BUILDER → [ ＋ Nowy ekspert ]
                                    │
                                    ▼
                   ╔═══════════════════════════════════════╗
                   ║   KONFIGURACJA SIEDMIU KOMPONENTÓW     ║   kolejność dowolna —
                   ╠═══════════════════════════════════════╣   pełna kompozycyjność,
                   ║ 1. Tożsamość             Agent Builder ║   żaden krok nie jest
                   ║ 2. Instrukcje systemowe  Agent Builder ║   wymuszony sekwencyjnie
                   ║ 3. Model bazowy + kanał  Model Config. ║   (Koncepcja platformy,
                   ║ 4. Umiejętności          Skills Mgr.   ║    rozdz. 12, zasada 1)
                   ║ 5. Pluginy/konektory/MCP Connectors    ║
                   ║ 6. Pamięć                Agent Builder ║
                   ║ 7. Zakres możliwości     Permissions   ║
                   ╚═══════════════════════╤═══════════════╝
                                            ▼
                                  [ Zapisz eksperta ]
                                            │
                       walidacja formy — informacyjna, nigdy blokująca:
                       nazwa pusta → uzupełniana nazwą roboczą „Nowy ekspert”;
                       nazwa zduplikowana w zasięgu widoczności → ostrzeżenie,
                       zapis przebiega dalej bez przerwania
                                            │
                            ┌───────────────┴────────────────┐
                            ▼                                  ▼
                  ZAPIS → NOWA WERSJA (w. 1)                BŁĄD ZAPISU
                  stan: ZAPISANY, AKTYWNY                   (np. utrata połączenia
                  zdarzenie `agent.changed`                   z serwerem)
                  (rozgłoszone na wszystkie urządzenia)       → .dn-toast--blad
                                                               → treść formularza
                                                                 zachowana, nic nie
                                                                 ginie — ponów
                                    │
                                    ▼
                     EKSPERT DOSTĘPNY OPERACYJNIE
                     (zasób Operatora w pamięci aplikacji,
                      widoczny w Bibliotece ekspertów)
                                    │
              ┌─────────────────────┼──────────────────────┬────────────────┐
              ▼                     ▼                        ▼                ▼
        wykonawca doraźny    Agent Manager             sekcja Role      pozostaje
        w TalkIn / WorkSpace  moduł Workspace —          panel           w bibliotece
        / CodeStudio          przypisanie do projektu     orkiestracji    bez użycia
        (rozdz. 10.1)         (rozdz. 10.2)                MultitaskingAI
                                                            (rozdz. 10.3)
              └─────────────────────┴──────────┬───────────┘
                                                ▼
                                       EKSPERT W UŻYCIU
                                                │
                              ┌─────────────────┴──────────────────┐
                              ▼                                     ▼
                    dalsza edycja definicji                 brak dalszych zmian
                    (dowolna z pięciu zakładek                → obowiązuje
                     Agent Builder)                              wersja aktywna
                              │
                              ▼
                    [ Zapisz zmiany ] → NOWA WERSJA (w. 2, 3, …)
                    poprzednie wersje zachowane w panelu „Historia wersji”;
                    zmiana obowiązuje od najbliższego wywołania — NIE przerywa
                    procesów, którym definicja została już dostarczona
                    (Specyfikacja okien operacyjnych, rozdz. 3 — trwałość stanu
                     procesu sesji)
                              │
                              ▼
                   decyzja Operatora
                   [ Przywróć wcześniejszą wersję ]  albo  [ Archiwizuj eksperta ]
                              │                                    │
                              ▼                                    ▼
                   wskazana wersja staje się              stan: ZARCHIWIZOWANY —
                   wersją AKTYWNĄ; poprzednia               niewybieralny jako NOWY
                   pozostaje w historii                      wykonawca; zachowany
                                                               w bibliotece, w historii
                                                               wersji i w przypisaniach
                                                               już istniejących
```

Ani przywrócenie wersji, ani archiwizacja nie usuwają danych — zgodnie z zasadą braku twardych blokad (Koncepcja platformy, rozdz. 6 i 14, zasada 2) obie operacje są odwracalnymi decyzjami Operatora, analogicznymi do `extension.toggle` zastosowanego wobec rozszerzeń (Rozszerzenia, rozdz. 5.2).

---

## 4. Stany

### 4.1. Stany eksperta jako komponentu własnego

| Stan | Znaczenie | Jak powstaje | Widoczny w |
|---|---|---|---|
| Szkic | Definicja rozpoczęta w Agent Builder, jeszcze nigdy niezapisana | Kliknięcie „＋ Nowy ekspert” | Nagłówek edytora — plakietka `.dn-plakietka--wersaliki` „SZKIC” |
| Niezapisane zmiany | Zmiana dokonana na zapisanej już definicji, jeszcze niezapisana jako nowa wersja | Edycja dowolnego pola w dowolnej zakładce | Nagłówek edytora — kropka `.dn-kropka--ostrz` przy nazwie |
| Zapisany / Aktywny | Wersja faktycznie wykorzystywana przy wybraniu eksperta jako wykonawcy | Kliknięcie „Zapisz eksperta” (walidacja formy jest informacyjna, nigdy nie blokuje zapisu — rozdz. 4.2); domyślnie najnowsza wersja | Biblioteka ekspertów i panel Historia wersji — plakietka `.dn-plakietka--sukces` „Aktywna” |
| Przypisany | Wykorzystywany w co najmniej jednym projekcie (Agent Manager) lub jednej roli (sekcja Role) | Jawna decyzja Operatora poza modułem Agents | Karta eksperta w bibliotece — plakietka „N przypisań” |
| Zarchiwizowany | Wycofany z aktywnego użycia jako nowy wykonawca | „Archiwizuj eksperta” | Biblioteka ekspertów, filtr „Zarchiwizowane” — plakietka tekstem `--dn-tekst-3` |

### 4.2. Stany interfejsu — reguła wspólna (System wizualny, rozdz. 6) i jej zastosowanie w module Agents

| Stan | Reguła wizualna | Przykład w module Agents |
|---|---|---|
| Domyślny | Wygląd spoczynkowy | Karta eksperta w bibliotece przed najechaniem |
| Hover / najechanie | `--dn-hover` albo uniesienie `translateY` dla elementów złotych | Karta eksperta unosi się `-1px`; wiersz rejestru w Skills/Connectors Manager podświetla tło |
| Fokus (klawiatura) | Pierścień złoty 2 px, wyłącznie `:focus-visible` | Każde pole formularza, każdy przełącznik macierzy uprawnień |
| Aktywny / zaznaczony | Akcent złoty: obrys, podkreślenie, tło `--dn-akcent-tlo` | Zakładka bieżąca edytora; pozycja wybrana w selektorze grup Permissions Center |
| Nieaktywny kontekstowo (pole bez zastosowania — nie: brak zgody) | `opacity: 0.5`, bez kursora `not-allowed` — pole pominięte w tab-order, nie zablokowane | Pole tokenu CLI, dopóki kanał „CLI” nie jest wybrany — pole nie ma zastosowania w bieżącym kontekście, nie „czeka na zgodę” |
| Ładowanie | `.dn-spinner` | Lista modeli w Model Configuration; rejestr rozszerzeń przy pierwszym otwarciu Skills/Connectors Manager |
| Ostrzeżenie | `.dn-kropka--ostrz` / obramowanie bursztynowe pola — sygnalizuje problem, nigdy nie blokuje akcji | Nazwa zduplikowana w zasięgu widoczności; pole nazwy puste w chwili zapisu (uzupełniane nazwą roboczą) |
| Błąd | `.dn-alert--blad` / `.dn-pole--blad` — wyłącznie awarie techniczne, nigdy wynik walidacji formy | Utrata połączenia przy zapisie; nieprawidłowy adres serwera MCP |

**Uwaga o kontrolkach akcji (przyciski, CTA).** W całym module Agents żaden przycisk akcji — „Zapisz eksperta”, „Testuj połączenie”, „Resetuj do pełnego dostępu” i pozostałe — nie występuje w stanie zablokowanym (`disabled`/`not-allowed`) jako sposób wymuszenia kompletności czy kolejności. Każdy pozostaje w pełni klikalny w każdej chwili; niegotowość lub brak danych sygnalizowana jest dopiero po kliknięciu — komunikatem, ostrzeżeniem lub dymkiem — albo opisowo obok kontrolki, nigdy przez odebranie interakcji. Stan „nieaktywny kontekstowo” powyżej dotyczy wyłącznie pól danych bez zastosowania w bieżącym kontekście (np. pole właściwe innemu kanałowi) — nigdy przycisków akcji i nigdy oczekiwania na zgodę, autoryzację czy zatwierdzenie. Żaden element modułu Agents nie blokuje akcji Operatora jako mechanizm kontroli dostępu (rozdz. 9.1).

---

## 5. Agent Builder — okno nadrzędne

### 5.1. Cel, zawartość ogólna i pozycja w przepływie

| Aspekt | Realizacja |
|---|---|
| Cel | Tworzenie i edycja definicji eksperta; scala wynik czterech okien konfiguracyjnych w jeden zapisywalny komponent własny (Specyfikacja agentów, rozdz. 4.1) |
| Zawartość natywna | Tożsamość (nazwa, opis przeznaczenia, moduły zastosowania, zasięg widoczności), instrukcje systemowe, konfiguracja pamięci, historia wersji |
| Zawartość osiągalna z tego okna | Model Configuration, Skills Manager, Connectors Manager, Permissions Center — jako zakładki tego samego edytora |
| Typologia komponentowa | Okno kreatora (buildera) — System wizualny, rozdz. 10.6: `.dn-zakladki` lub kroki formularza, `.dn-karta` boczna z podsumowaniem, `.dn-btn--zloty` jako jedyne CTA finalizujące |
| Dwa widoki | **Biblioteka ekspertów** (punkt wejścia — przegląd istniejących) i **Edytor** (tworzenie/edycja jednego eksperta) |

### 5.2. Widok „Biblioteka ekspertów”

#### 5.2.1. Makieta

```
Makieta — Agent Builder, widok: Biblioteka ekspertów (stan spoczynku)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW      │ ⟳ EXECUTION LOOP  │ 🧑 AGENT BUILDER — Biblioteka
  Użytkownik ↔        │ Koordynator ↔     │                      [＋ Nowy ]
  Wykonawca           │ Wykonawca         │ ─────────────────────────────────
 ─────────────────────┼───────────────────┤ [ 🔍 szukaj… ]  [ Filtry ▼ ]
  strumień rozmowy    │ zlecenie i        │ ─────────────────────────────────
                      │ dekompozycja      │ ┌──────────┐ ┌──────────┐
                      │ kolejka zadań     │ │▣ Agent   │ │▣ Agent   │
                      │ kontrola jakości  │ │  Redaktor│ │  Backend │
                      │ ponowienia: 0     │ │ w.3 ·    │ │ w.7 ·    │
                      │ przebieg ▮▮▯▯     │ │ globalny │ │ globalny │
 ─────────────────────┤                   │ │[Aktywny] │ │[Aktywny] │
  [ Napisz polecenie ]│ [ Sterowanie ▼ ]  │ │ ⋮        │ │ ⋮        │
  [ Agent ▼ ][Wyślij ]│                   │ └──────────┘ └──────────┘
                      │                   │ … kolejne karty w siatce …
 ══════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu: [Danaco Console] [Ubuntu] [Agent ▼] [Model ▼] [Ultra]
 ══════════════════════════════════════════════════════════════════════════════
```

#### 5.2.2. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Nagłówek widoku | Etykieta tekstowa z ikoną `uzytkownik` | Identyfikuje okno i jego przynależność do modułu Agents | 1 | Widoczny bez interakcji | Napis krojem nagłówkowym, duża waga, stały u góry okna | Domyślny | — (statyczny) | Górna krawędź obu widoków Agent Buildera |
| Przycisk „＋ Nowy ekspert” | Przycisk pierwszoplanowy `.dn-btn--zloty` z ikoną `plus` | Rozpoczyna tworzenie nowej definicji eksperta | 1 | Widoczny bez interakcji | Duży przycisk, jedyny CTA złoty widoku — zasada oszczędnego złota (System wizualny, rozdz. 8.1) | Domyślny · hover (uniesienie 2 px) · aktywny (wciśnięcie) · fokus (pierścień złoty) | Otwiera widok Edytor z pustym formularzem w stanie „Szkic” | Prawy górny róg widoku Biblioteka |
| Pole wyszukiwania | Pole tekstowe `.dn-input` z ikoną `szukaj` | Filtruje listę po nazwie lub opisie przeznaczenia | 2 | Kliknięcie ikony `szukaj` w nagłówku widoku; po użyciu pole zwija się samoczynnie | Pole średnie, pełna szerokość paska filtrów przy braku innych filtrów | Domyślny · fokus · wypełnione (z przyciskiem `zamknij` czyszczącym) | Filtruje siatkę kart na żywo, bez przeładowania widoku | Pasek filtrów widoku Biblioteka |
| Filtr „Środowisko” | Rozwijana lista `.dn-select` (wielokrotny wybór) | Zawęża do ekspertów wykorzystywanych w TalkIn / WorkSpace / CodeStudio / MultitaskingAI | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka w pasku filtrów | Domyślny · rozwinięty · zaznaczone pozycje (plakietka liczby) | Zawęża siatkę kart | Pasek filtrów |
| Filtr „Widoczność” | `.dn-select`: Globalny / Projektowy | Zawęża wg zasięgu widoczności (Specyfikacja agentów, rozdz. 5.2) | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka | Jak wyżej | Zawęża siatkę kart | Pasek filtrów |
| Filtr „Stan” | `.dn-select`: Aktywny / Szkic / Zarchiwizowany | Zawęża wg stanu eksperta (rozdz. 4.1) | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka | Jak wyżej | Zawęża siatkę kart | Pasek filtrów |
| Karta eksperta | `.dn-karta` w układzie siatki | Reprezentuje jedną zapisaną definicję eksperta; punkt wejścia do edycji | 1 | Widoczna bez interakcji | Karta średnia, siatka 3–4 kolumn responsywnie | Domyślny · hover (`--interaktywna`, podświetlenie tła) · zaznaczona (przy akcji zbiorczej, jeśli skonfigurowana) | Kliknięcie otwiera widok Edytor z załadowaną definicją (wersja aktywna) | Główny obszar widoku Biblioteka |
| Awatar eksperta na karcie | `.dn-awatar--kwadrat` | Identyfikacja niepersonalna jednostki AI (System wizualny, rozdz. 8.8: wariant kwadratowy zarezerwowany dla agentów) | 1 | Widoczny bez interakcji | Mała ikonka, lewy górny róg karty | Domyślny | — | Karta eksperta |
| Nazwa eksperta (na karcie) | Tekst krojem bazowym, `semibold` | Główny identyfikator karty | 1 | Widoczna bez interakcji | Średnia waga tekstowa, jedna linia z przycięciem | Domyślny | — | Karta eksperta |
| Opis przeznaczenia (na karcie) | Tekst `--dn-tekst-2`, maks. dwie linie | Skrót zadania eksperta | 1 | Widoczny bez interakcji | Mała waga tekstowa | Domyślny | — | Karta eksperta |
| Wiersz „model · kanał” | Tekst pomocniczy `--dn-tekst-3` | Skrót konfiguracji z Model Configuration | 1 | Widoczny bez interakcji | Bardzo mała waga (metadane) | Domyślny | — | Karta eksperta |
| Plakietka wersji | `.dn-plakietka` mała, np. „w.3” | Numer bieżącej wersji aktywnej | 1 | Widoczna bez interakcji | Bardzo mała, pigułka | Domyślny | Najechanie pokazuje tooltip z datą ostatniego zapisu | Karta eksperta |
| Plakietka widoczności | `.dn-plakietka` „globalny” / „projektowy” | Zasięg dostępności eksperta | 1 | Widoczna bez interakcji | Bardzo mała, pigułka | Domyślny | — | Karta eksperta |
| Plakietka stanu | `.dn-plakietka--stan` z `.dn-kropka` (sukces/ostrzeżenie/tekst-3) | Aktywny / Szkic / Zarchiwizowany | 1 | Widoczna bez interakcji | Mała, pigułka z kropką | sukces (aktywny) · ostrzeżenie (szkic) · neutralny (zarchiwizowany) | — | Karta eksperta |
| Licznik przypisań | Tekst „N przypisań” / „1 rola” | Sygnalizuje wykorzystanie w Agent Manager lub sekcji Role | 1 | Widoczny bez interakcji | Bardzo mała waga | Domyślny · zero (tekst „bez przypisań”, `--dn-tekst-3`) | Najechanie pokazuje tooltip z listą projektów/ról | Karta eksperta |
| Menu akcji karty | Trzy przyciski tekstowe/`.dn-btn-ikona`: Edytuj · Duplikuj · Więcej (`wiecej`) | Szybkie akcje bez otwierania pełnego edytora | 3 | Menu kebab (⋮) na karcie | Mały pasek na stopce karty | Domyślny · hover | „Edytuj” → widok Edytor; „Duplikuj” → nowy szkic wypełniony tą definicją; „Więcej” rozwija: Archiwizuj (`archiwum`), Usuń (`kosz`, wariant `.dn-btn--blad` — sygnalizuje wagę, nie blokuje) | Stopka karty eksperta |
| Pusty stan | `.dn-pusty-stan` — ikona `uzytkownik`, tytuł, opis, CTA | Wyświetlany, gdy filtr lub brak jakichkolwiek ekspertów nie zwraca wyników | 1 | Widoczny bez interakcji | Duży, wyśrodkowany blok zastępujący siatkę | Aktywny tylko przy pustym wyniku | CTA „Utwórz pierwszego eksperta” otwiera widok Edytor | Środek widoku Biblioteka |
| Chat Window i Execution Loop Window | Dwie stałe kolumny komunikacji operacyjnej (rozdz. 2.4) | Kanał Użytkownik ↔ Wykonawca oraz kanał Koordynator ↔ Wykonawca — wsparcie konwersacyjne przy przeglądaniu i wyborze Wykonawcy oraz podgląd pętli wykonawczej | 1 | Widoczne bez interakcji | Lewa kolumna, stała, pełna wysokość obszaru roboczego; kolumna sąsiadująca dla pętli wykonawczej | Domyślny · strumień na żywo (WebSocket) | Polecenia w kontekście modułu Agents (np. „znajdź Wykonawcę do…”) | Lewa strona układu w każdym widoku modułu Agents |

### 5.3. Widok „Edytor” — pasek zakładek i zakładka Tożsamość

#### 5.3.1. Makieta

```
Makieta — Agent Builder, widok: Edytor, zakładka Tożsamość (stan spoczynku)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW   │ ⟳ EXECUTION LOOP │ ← Biblioteka  [ Nazwa ] [SZKIC] ⏺ zapisano
  Użytkownik ↔     │ Koordynator ↔    │ ────────────────────────────┬───────────
  Wykonawca        │ Wykonawca        │ 🧑Tożsamość │🚀Model │⭐Skille│ PODSUMO-
 ──────────────────┼──────────────────┤ 🔗Konektory │🛡Uprawnienia   │ WANIE
  strumień rozmowy │ zlecenie i       │ Nazwa       [ .dn-input ]   │ DEFINICJI
                   │ dekompozycja     │ Opis        [ .dn-textarea ]│ Model: —
                   │ kolejka zadań    │ Moduły  ☐TalkIn ☐WorkSpace  │ Kanał: —
                   │ kontrola jakości │         ☐CodeStudio         │ Skille: 0
                   │ ponowienia: 0    │ ⓘ rola w MultitaskingAI     │ Rozszerz.: 0
                   │ przebieg ▮▮▯▯    │   przypisywana osobno       │ Uprawnienia:
                   │                  │ Widoczność ( )Globalny      │  pełny dostęp
                   │                  │            ( )Projektowy    │ Pamięć: sesja
                   │                  │ Instrukcje systemowe        │ ──────────
                   │                  │ ┌─────────────────────────┐ │ HISTORIA
                   │                  │ │ .dn-textarea duże       │ │ WERSJI ▼
                   │                  │ └─────────────────────────┘ │  w.3 Aktywna
 ──────────────────┤                  │ Pamięć [ Poziomy ▼ ]        │  w.2, w.1
  [ Napisz polec. ]│ [ Sterowanie ▼ ] │ ────────────────────────────┴───────────
  [ Agent ▼ ][Wyśl]│                  │            [ Anuluj ]  [ Zapisz eksperta ]
 ══════════════════════════════════════════════════════════════════════════════
```

#### 5.3.2. Katalog elementów — nagłówek i pasek zakładek

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Odnośnik „← Biblioteka” | Link tekstowy z ikoną `strzalka-lewo` | Powrót do widoku Biblioteka | 1 | Widoczny bez interakcji | Mały, w lewym rogu nagłówka | Domyślny · hover | Powrót; przy niezapisanych zmianach — nieblokujący toast „zmiany zachowane jako szkic” (nie modal zatrzymujący) | Nagłówek edytora |
| Pole nazwy (inline, w nagłówku) | `.dn-input` bez obramowania w stanie spoczynku, edytowalne po kliknięciu | Szybka zmiana nazwy bez wchodzenia w zakładkę Tożsamość | 1 | Widoczny bez interakcji | Duża waga — krój nagłówkowy przy odczycie, `.dn-input` przy edycji | Domyślny (odczyt) · edycja (fokus) · ostrzeżenie (nazwa pusta lub zduplikowana — obramowanie bursztynowe, nieblokujące) | Zapis nazwy przy utracie fokusu; nazwa pusta → uzupełniana nazwą roboczą „Nowy ekspert”; duplikat oznacza pole z komunikatem ostrzegawczym — zapis i pozostałe pola pozostają w pełni dostępne | Nagłówek edytora, zsynchronizowane z polem „Nazwa eksperta” w zakładce Tożsamość |
| Plakietka stanu | `.dn-plakietka--wersaliki` | SZKIC / AKTYWNY / ZARCHIWIZOWANY (rozdz. 4.1) | 1 | Widoczna bez interakcji | Mała pigułka obok nazwy | Wg stanu encji | — | Nagłówek edytora |
| Wskaźnik zapisu | Mały tekst z kropką `.dn-kropka` | Informuje, czy bieżący stan formularza jest zapisany | 1 | Widoczny bez interakcji | `--dn-tekst-3`, bardzo mała waga | „zapisano o HH:MM” (sukces) · „zapisywanie…” (`.dn-spinner`) · „niezapisane zmiany” (ostrzeżenie) · „błąd zapisu” (błąd, z linkiem „ponów”) | Aktualizuje się na żywo przy każdej zmianie pola i po każdej próbie zapisu | Nagłówek edytora, obok plakietki stanu |
| Pasek zakładek | `.dn-zakladki` z ikonami | Nawigacja między pięcioma komponentami definicji | 1 | Widoczny bez interakcji | Pełna szerokość okna, pas średniej wysokości pod nagłówkiem | Zakładka bieżąca (podkreślenie złote) · pozostałe (domyślny) · fokus klawiaturowy | Przełącza obszar zawartości bez przeładowania okna; stan pozostałych zakładek zachowany | Pod nagłówkiem edytora, widoczny we wszystkich pięciu zakładkach |

#### 5.3.3. Katalog elementów — zakładka Tożsamość

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole „Nazwa eksperta” | `.dn-input` jednowierszowe | Nazwa własna odróżniająca eksperta jako nazwaną jednostkę (Specyfikacja agentów, rozdz. 5.2) | 1 | Widoczny bez interakcji | Duże, u góry treści zakładki | Domyślny · fokus · ostrzeżenie (pusta/zduplikowana — sygnalizowane, nigdy nieblokujące zapisu) | Synchronizuje się z polem nazwy w nagłówku | Zakładka Tożsamość |
| Pole „Opis przeznaczenia” | `.dn-textarea` mała (2–3 wiersze) | Krótki opis zadania eksperta | 1 | Widoczny bez interakcji | Średnie, pod nazwą | Domyślny · fokus · puste (dopuszczalne — brak wymogu) | Aktualizuje opis widoczny na karcie w Bibliotece | Zakładka Tożsamość |
| Grupa „Moduły zastosowania” | Zestaw `.dn-check` (checkbox) — TalkIn, WorkSpace, CodeStudio + informacja pomocnicza o roli w MultitaskingAI | Środowiska i moduły wykorzystania operacyjnego (Specyfikacja agentów, rozdz. 5.2) | 1 | Widoczna bez interakcji | Mała grupa kontrolek w karcie | Domyślny (żadne niezaznaczone = dostępny wszędzie, brak ograniczenia) · zaznaczone | Zaznaczenie zawęża, gdzie ekspert pojawia się jako wykonawca doraźny w wyborze komponentu własnego | Zakładka Tożsamość |
| Adnotacja „rola w MultitaskingAI przypisywana osobno” | Tekst pomocniczy z ikoną `info` | Wyjaśnia, że rola w środowisku MultitaskingAI ustala się w sekcji Role, nie tutaj (rozdz. 10.3) | 1 | Widoczna bez interakcji | Mała, `--dn-tekst-3` | Domyślny | Kliknięcie prowadzi (kontekstowo) do wyjaśnienia w Chat Window | Zakładka Tożsamość, przy grupie modułów |
| Grupa „Zasięg widoczności” | `.dn-check` w układzie radio: Globalny / Projektowy | Zakres dostępności eksperta (Specyfikacja agentów, Załącznik D.1) | 1 | Widoczna bez interakcji | Mała grupa kontrolek | Domyślny (Globalny) · Projektowy wybrany | Projektowy odsłania pole wyboru projektu macierzystego | Zakładka Tożsamość |
| Pole „Instrukcje systemowe” | `.dn-textarea` duża, wieloliniowa | Tekstowy opis sposobu działania eksperta, przekazywany modelowi przy każdym wywołaniu (Specyfikacja agentów, rozdz. 5.3) | 1 | Widoczny bez interakcji | Największy pojedynczy element treści zakładki — duży panel | Domyślny · fokus · puste (dopuszczalne) · licznik znaków przy zbliżaniu się do limitu technicznego | Treść zapisywana wprost do definicji; podgląd w panelu podsumowania skraca do pierwszej linii | Zakładka Tożsamość, środek |
| Grupa „Konfiguracja pamięci” | Zestaw `.dn-suwak` — Globalna, Projekt, Sesja, Środowisko, Wyłączona | Poziomy pamięci, z których ekspert korzysta domyślnie (Specyfikacja agentów, rozdz. 5.6) | 3 | Element zbiorczy `Poziomy ▼` — rozwinięcie zawiera pełną listę poziomów pamięci | Mała grupa przełączników w karcie | Domyślny (Sesja włączona) · dowolna kombinacja poziomów · „Wyłączona” wyklucza pozostałe (wzajemnie wykluczający się przełącznik) | Zmiana natychmiast odzwierciedlona w panelu podsumowania | Zakładka Tożsamość, dół |
| Panel „Podsumowanie definicji” | `.dn-karta` boczna, stała podczas przewijania treści zakładki | Podgląd wszystkich siedmiu komponentów na żywo, niezależnie od aktywnej zakładki | 1 | Widoczny bez interakcji | Średni panel, prawa kolumna edytora | Aktualizuje się przy każdej zmianie w dowolnej zakładce | Kliknięcie wiersza podsumowania (np. „Model: —”) przenosi do właściwej zakładki | Prawa kolumna wszystkich pięciu zakładek edytora |
| Przycisk „Anuluj” | `.dn-btn--duch` | Odrzuca niezapisane zmiany bieżącej sesji edycji | 1 | Widoczny bez interakcji | Mały, obok „Zapisz eksperta” | Domyślny · hover | Przywraca ostatnią zapisaną wersję do formularza; nie usuwa historii | Stopka edytora |
| Przycisk „Zapisz eksperta” | `.dn-btn--zloty` | Zapisuje bieżący stan formularza jako nową wersję | 1 | Widoczny bez interakcji | Duży, jedyny CTA złoty widoku Edytor | Domyślny — zawsze w pełni klikalny, również z nazwą pustą · ładowanie (`.dn-spinner` w przycisku) · błąd (toast, wyłącznie awaria techniczna) | Sukces: nazwa pusta uzupełniana automatycznie nazwą roboczą „Nowy ekspert” z ostrzeżeniem; nowa wersja, plakietka „zapisano”, karta w Bibliotece aktualizuje się na żywo (`agent.changed`) | Stopka edytora, widoczna na wszystkich pięciu zakładkach |

### 5.4. Panel „Historia wersji”

#### 5.4.1. Makieta

```
Makieta — Agent Builder, panel Historia wersji (rozszerzenie boczne, stan rozwinięty)

 ══════════════════════════════════════════════════════════════════════════════
  💬 Chat  │ ⟳ Execution │ OBSZAR ROBOCZY — Edytor │ HISTORIA WERSJI      [ ✕ ]
  Window   │ Loop Window │ (zakładka bieżąca)      │ Agent Redaktor
 ──────────┼─────────────┼─────────────────────────┼───────────────────────────
  strumień │ kolejka     │ formularz zakładki      │ ● w.3 Aktywna 12:04
  rozmowy  │ zadań       │                         │   instrukcje, +1 umiejętność
           │ kontrola    │                         │   [ Podgląd ]
           │ jakości     │                         │ ○ w.2  wczoraj, 18:41
           │ ponowienia  │                         │   +konektor „repozytorium”
           │ przebieg    │                         │   [ Podgląd ] [ Przywróć ]
 ──────────┤             │                         │ ○ w.1  3 dni temu, 09:12
  [ pole ] │ [ Sterow. ▼]│                         │   [ Podgląd ] [ Przywróć ]
  [Agent ▼]│             │                         │ [ Duplikuj jako nowy… ]
 ══════════════════════════════════════════════════════════════════════════════
```

#### 5.4.2. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Otwierany panel „Historia wersji” | Panel boczny/rozwijany (wariant `.dn-modal` lub panel dokujący) z ikoną `archiwum` | Przegląd wszystkich zapisanych wersji eksperta | 3 | Element `Historia wersji ▼` w panelu podsumowania — panel otwiera się jako rozszerzenie boczne | Średni panel, wywoływany z linku „Zobacz wszystkie” w podsumowaniu (rozdz. 5.3.3) | Zwinięty (domyślnie, widoczne 3 ostatnie w podsumowaniu) · rozwinięty (pełna lista) | Kliknięcie „Zobacz wszystkie” rozwija pełną listę | Prawa kolumna edytora / panel wywoływany |
| Wiersz wersji | `.dn-karta--pozycja` | Reprezentuje jeden zapisany stan definicji | 3 | Widoczny po otwarciu panelu Historia wersji | Mały wiersz listy | Aktywna (`.dn-kropka--sukces` wypełniona) · archiwalna (`.dn-kropka` pusta) | — | Panel Historia wersji |
| Skrót zmian | Tekst `--dn-tekst-2`, jedna–dwie linie | Streszcza różnicę względem poprzedniej wersji (pole, umiejętności, rozszerzenia) | 3 | Widoczny po otwarciu panelu Historia wersji | Mała waga | Domyślny | Najechanie pokazuje pełną listę zmienionych pól w `.dn-tooltip` | Wiersz wersji |
| Przycisk „Podgląd” | `.dn-btn--duch`, mały | Otwiera daną wersję w trybie tylko do odczytu, bez zmiany wersji aktywnej | 3 | Widoczny po otwarciu panelu Historia wersji | Mały | Domyślny | Otwiera formularz edytora w trybie podglądu (pola nieedytowalne, baner „podgląd wersji w.N”) | Każdy wiersz wersji |
| Przycisk „Przywróć tę wersję” | `.dn-btn--zarys`, mały, ikona `odswiez` | Ustawia wskazaną wersję jako aktywną | 3 | Widoczny po otwarciu panelu Historia wersji | Mały | Domyślny · niedostępny dla wersji już aktywnej (naturalnie ukryty, nie „zablokowany”) | Tworzy nową wersję o treści identycznej ze wskazaną (nie usuwa wersji pośrednich); wersja aktywna zmienia się natychmiast | Wiersz wersji innej niż aktywna |
| Przycisk „Duplikuj eksperta jako nowy…” | `.dn-btn--zarys`, ikona brak (etykieta tekstowa) | Tworzy odrębnego, nowego eksperta zaczynającego się od bieżącej definicji | 3 | Widoczny po otwarciu panelu Historia wersji | Mały, stopka panelu | Domyślny | Otwiera widok Edytor w stanie „Szkic” z polami wypełnionymi kopią definicji; nazwa jest wstępnie zmieniana (np. sufiks „kopia”), a przy nazwie już istniejącej pojawia się ostrzeżenie — nigdy blokada zapisu, spójnie z zasadą walidacji ostrzegawczej | Stopka panelu Historia wersji |

### 5.5. Stany szczegółowe okna Agent Builder

| Stan okna | Wyzwalacz | Co widać |
|---|---|---|
| Pusta biblioteka | Brak jakiegokolwiek zapisanego eksperta | `.dn-pusty-stan` z CTA „Utwórz pierwszego eksperta” |
| Przeglądanie | Wejście do modułu, co najmniej jeden ekspert istnieje | Siatka kart, filtry aktywne |
| Tworzenie nowego | Kliknięcie „＋ Nowy ekspert” | Edytor, wszystkie pola puste, zakładka Tożsamość domyślnie aktywna, plakietka „SZKIC” |
| Edycja | Wybór karty z Biblioteki | Edytor wypełniony danymi wersji aktywnej |
| Podgląd wersji archiwalnej | „Podgląd” w Historii wersji | Edytor w trybie tylko do odczytu, baner informacyjny u góry, brak przycisku „Zapisz” |
| Zapisywanie | Kliknięcie „Zapisz eksperta” | Przycisk ze spinnerem, pola tymczasowo tylko do odczytu |
| Zapisano | Potwierdzenie serwera | Toast sukcesu, wskaźnik zapisu aktualizuje znacznik czasu, nowa pozycja w Historii wersji |
| Błąd zapisu | Utrata połączenia lub odrzucenie przez serwer | Toast błędu z treścią przyczyny, formularz pozostaje edytowalny z wprowadzonymi danymi |

---

## 6. Model Configuration

### 6.1. Cel i zawartość

Model Configuration ustala model bazowy eksperta oraz kanał połączenia platformy z tym modelem (Specyfikacja agentów, rozdz. 5.1). Wartość zapisana tutaj jest domyślna i nadpisywalna per rola w chwili przypisania eksperta w sekcji Role środowiska MultitaskingAI (rozdz. 8.3 Specyfikacji agentów; rozdz. 10.3 niniejszego dokumentu).

### 6.2. Makieta

```
Makieta — zakładka Model Configuration (stan spoczynku)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW   │ ⟳ EXECUTION LOOP │ 🧑Tożsamość │🚀Model Configuration │⭐Skille
  Użytkownik ↔     │ Koordynator ↔    │ 🔗Konektory │🛡Uprawnienia  ┬──────────────
  Wykonawca        │ Wykonawca        │ Model bazowy [ Model ▼ ]    │ PODSUMOWANIE
 ──────────────────┼──────────────────┤ ┌───────┐┌───────┐┌───────┐ │ DEFINICJI
  strumień rozmowy │ kolejka zadań    │ │Model A││Model B││Model C│ │ (jak w
                   │ kontrola jakości │ │○      ││●      ││○      │ │  rozdz.
                   │ ponowienia: 0    │ └───────┘└───────┘└───────┘ │  5.3.3)
                   │ przebieg ▮▮▯▯    │ Kanał ⦿API ○CLI ○SSH ○HTTP  │
                   │                  │ ⓘ API — klucz dostępu,      │
                   │                  │   żądanie HTTP              │
                   │                  │ Dane dostępowe              │
                   │                  │ [ sk-ant-…  🙈 ukryj·zmień ]│
                   │                  │ [ Testuj połączenie ]       │
 ──────────────────┤                  │ ⓘ nadpisywalne per rola     │
  [ Napisz polec. ]│ [ Sterowanie ▼ ] │ ────────────────────────────┴──────────────
  [ Agent ▼ ][Wyśl]│                  │            [ Anuluj ]  [ Zapisz eksperta ]
 ══════════════════════════════════════════════════════════════════════════════
```

### 6.3. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Siatka kart „Model bazowy” | Karty wyboru `.dn-karta--interaktywna` w układzie poziomym | Wybór modelu AI stanowiącego podstawę eksperta | 2 | Kliknięcie znacznika `Model ▼`; po wyborze lista zwija się samoczynnie | Duży panel, główny element zakładki | Domyślny · wybrany (`--akcent`, obrys złoty) · wygaszony jako wskazówka, gdy model wymaga innego kanału niż obecnie wybrany — bez ukrywania, zawsze klikalny | Wybór ustawia model bazowy definicji; jeśli karta wymagała innego kanału, wybór przełącza automatycznie segment „Kanał modelu” na właściwy | Górna część zakładki Model Configuration |
| Segment „Kanał modelu” | `.dn-zakladki--pigulki`: API, CLI, SSH, HTTP | Sposób połączenia platformy z modelem (Architektura, rozdz. 9) | 2 | Segment kanału w zakładce; po wyborze pozostaje zwinięty do wybranej wartości | Średni pasek segmentowy | Domyślny (API zaznaczone) · wybrany segment (tło akcentu) | Zmiana kanału odsłania właściwe mu pola poniżej (dane dostępowe, host, adres) | Środek zakładki |
| Objaśnienie kontekstowe kanału | `.dn-alert--info` / `.dn-tooltip` `[?]` | Opisuje zastosowanie typowe wybranego kanału (tabela Specyfikacji agentów, rozdz. 5.1) | 3 | Objaśnienie `[?]` — dymek kontekstowy | Mały blok tekstu pod segmentem | Zmienia treść wraz z wybranym kanałem | — | Pod segmentem „Kanał modelu” |
| Pole „Dane dostępowe kanału” | `.dn-input` jawny + `.dn-btn-ikona` „ukryj/maskuj” + link „zmień” | Odwołanie do danych dostępowych kanału (klucz API, token CLI, dane SSH), przechowywanych poza bazą danych (Model danych, rozdz. 8.1) | 2 | Widoczne po wyborze kanału wymagającego danych dostępowych | Średnie pole | Domyślny (wartość jawna, w pełni widoczna dla Operatora — zgodnie z zasadą kluczy jawnych) · zamaskowany (ukrycie jednym kliknięciem, do ponownego odsłonięcia w każdej chwili) · puste (wymagane przed użyciem operacyjnym, nie przed zapisem) | „ukryj/maskuj” przełącza natychmiast — zawsze dostępne, nigdy zablokowane; „zmień” otwiera sekcję „integracje” okna konfiguracji platformy (poza modułem Agents) — rozdzielenie decyzji analogiczne do Rozszerzeń, rozdz. 7; domyślnym trybem jest wartość jawna; domyślne maskowanie jest ustawieniem konfiguracyjnym sterowanym z okna konfiguracji platformy | Zależnie od wybranego kanału |
| Przycisk „Testuj połączenie” | `.dn-btn--zarys`, ikona `odswiez` | Weryfikuje osiągalność kanału przed zapisem | 2 | Przycisk w zakładce Model Configuration | Mały przycisk | Domyślny — zawsze w pełni klikalny; brak wybranych danych dostępowych → komunikat „uzupełnij dane dostępowe” przy kliknięciu, zamiast blokady · ładowanie · sukces (`.dn-plakietka--sukces` „połączono”) · błąd (`.dn-plakietka--blad` z treścią) | Wynik widoczny lokalnie w zakładce; nie warunkuje możliwości zapisania eksperta | Obok pola danych dostępowych |
| Adnotacja „nadpisywalne per rola” | Tekst pomocniczy z ikoną `info` | Przypomina, że wartość jest domyślna, nadpisywalna przy przypisaniu do roli MultitaskingAI (Specyfikacja agentów, rozdz. 8.3) | 1 | Widoczna bez interakcji | Mała, `--dn-tekst-3` | Domyślny | — | Dół zakładki |

### 6.4. Stany

| Stan | Wyzwalacz | Co widać |
|---|---|---|
| Ładowanie listy modeli | Pierwsze otwarcie zakładki w danej sesji | `.dn-spinner` w miejscu siatki kart |
| Kanał zmieniony | Wybór innego segmentu | Pola specyficzne dla poprzedniego kanału znikają, pojawiają się właściwe nowemu (bez utraty już wprowadzonych wspólnych danych) |
| Test połączenia w toku | Kliknięcie „Testuj połączenie” | Przycisk ze spinnerem, reszta zakładki pozostaje edytowalna |
| Błąd testu | Nieosiągalny kanał | `.dn-alert--blad` z treścią przyczyny, możliwość zapisania eksperta mimo błędu testu (test jest narzędziem pomocniczym, nie bramką) |

---

## 7. Skills Manager

### 7.1. Cel i zawartość

Skills Manager dobiera umiejętności (skille) rozszerzające zakres działania eksperta o wyspecjalizowane procedury lub wiedzę proceduralną (Specyfikacja agentów, rozdz. 5.4; Rozszerzenia, rozdz. 1 i 6). Okno operuje na wspólnym rejestrze rozszerzeń platformy (Rozszerzenia, rozdz. 7) — rejestracja i włączenie umiejętności w rejestrze to decyzja podejmowana w oknie konfiguracji platformy; **przypisanie** zarejestrowanej i włączonej umiejętności do TEGO eksperta jest odrębną decyzją podejmowaną tutaj (Rozszerzenia, rozdz. 7, tabela dwóch decyzji).

### 7.2. Makieta

```
Makieta — zakładka Skills Manager (stan spoczynku)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW   │ ⟳ EXECUTION LOOP │ 🧑Tożsamość │🚀Model │⭐Skills Manager
  Użytkownik ↔     │ Koordynator ↔    │ 🔗Konektory │🛡Uprawnienia ┬───────────────
  Wykonawca        │ Wykonawca        │ [ 🔍 szukaj ] [ Filtry ▼ ] │ PODSUMOWANIE
 ──────────────────┼──────────────────┤ ────────────────────────── │ DEFINICJI
  strumień rozmowy │ kolejka zadań    │ ⭐ Korekta i redakcja      │ Umiejętności
                   │ kontrola jakości │    Danaco Plugin  [ ▣ ]  ⋮ │ przypisane: 2
                   │ ponowienia: 0    │ ⭐ Streszczenia wielojęz.  │
                   │ przebieg ▮▮▯▯    │    Danaco Plugin  [ ▣ ]  ⋮ │
                   │                  │ ⭐ Analiza sentymentu      │
                   │                  │    Personal       [ ☐ ]  ⋮ │
                   │                  │    ⓘ wyłączona w rejestrze │
                   │                  │ [ Operacje ▼ ]             │
 ──────────────────┤                  │ ───────────────────────────┴───────────────
  [ Napisz polec. ]│ [ Sterowanie ▼ ] │            [ Anuluj ]  [ Zapisz eksperta ]
  [ Agent ▼ ][Wyśl]│                  │
 ══════════════════════════════════════════════════════════════════════════════
```

### 7.3. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole wyszukiwania umiejętności | `.dn-input` z ikoną `szukaj` | Filtruje rejestr po nazwie | 2 | Kliknięcie ikony `szukaj` w pasku zakładki | Średnie pole | Domyślny · fokus · wypełnione | Filtruje listę na żywo | Górny pasek zakładki |
| Filtr „Źródło” | `.dn-select`: Danaco Plugin / Personal | Zawęża wg pochodzenia (Rozszerzenia, rozdz. 2) | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka | Domyślny · rozwinięty | Filtruje listę | Górny pasek |
| Filtr „Stan” | `.dn-select`: Przypisane temu ekspertowi / Nieprzypisane / Wyłączone w rejestrze | Zawęża wg statusu przypisania | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka | Domyślny · rozwinięty | Filtruje listę | Górny pasek |
| Wiersz umiejętności | `.dn-karta--pozycja` z ikoną `gwiazdka` | Reprezentuje jedną zarejestrowaną umiejętność | 1 | Widoczny bez interakcji | Średni wiersz listy (`.dn-tabela` bez widocznego nagłówka lub karta pozycji) | Włączona w rejestrze / Wyłączona w rejestrze (`.dn-plakietka`) | Kliknięcie nazwy rozwija panel szczegółów (opis, parametry konfiguracyjne właściwe tej umiejętności, objaśnienie `[?]`) | Lista główna zakładki |
| Przełącznik „Przypisana temu ekspertowi” | `.dn-suwak` w wierszu | Decyduje, czy TEN ekspert korzysta z danej umiejętności | 1 | Widoczny bez interakcji | Mały przełącznik, prawa krawędź wiersza | Domyślny (wyłączony) · włączony — zawsze w pełni klikalny, również gdy umiejętność jest wyłączona w rejestrze globalnym | Włączenie zapisuje przypisanie do definicji eksperta; jeśli umiejętność jest wciąż wyłączona w rejestrze globalnym, zapisuje się jako „przypisana, oczekuje aktywacji w rejestrze” (adnotacja `[?]` wskazuje, gdzie włączyć) i zaczyna działać automatycznie po włączeniu jej tam; wyłączenie usuwa przypisanie bez usuwania umiejętności z rejestru platformy | Każdy wiersz umiejętności |
| Adnotacja „włącz w oknie konfiguracji” | Link tekstowy z ikoną `info` | Prowadzi do rejestru globalnego, gdy umiejętność jest wyłączona tam, a Operator chce ją przypisać | 3 | Widoczna w wierszu umiejętności wyłączonej w rejestrze | Mała, `--dn-tekst-3` | Widoczna wyłącznie przy umiejętnościach wyłączonych w rejestrze | Otwiera sekcję „rozszerzenia” okna konfiguracji platformy | Wiersz umiejętności wyłączonej w rejestrze |
| Przycisk „Zainstaluj umiejętność Personal…” | `.dn-btn--zarys` z ikoną `plus` | Skrót do instalacji nowego rozszerzenia typu Personal (Rozszerzenia, rozdz. 5.1) | 3 | Grupa `Operacje ▼` w stopce listy | Mały przycisk, stopka listy | Domyślny · hover | Otwiera przepływ instalacji (wskazanie pliku → przesłanie → rejestracja → potwierdzenie) w oknie konfiguracji platformy | Stopka zakładki Skills Manager |
| Panel szczegółów umiejętności | Rozwijany fragment wiersza | Prezentuje opis i parametry konfiguracyjne (`config.get`/`config.set`) właściwe tej instancji rozszerzenia (Rozszerzenia, rozdz. 4) | 3 | Kliknięcie nazwy umiejętności — panel rozwija się w miejscu | Średni, rozwijany w miejscu | Zwinięty (domyślnie) · rozwinięty | Edycja parametru zapisuje się natychmiast w atrybucie Konfiguracja rozszerzenia | Pod wierszem umiejętności po kliknięciu |

### 7.4. Stany

| Stan | Wyzwalacz | Co widać |
|---|---|---|
| Ładowanie rejestru | Pierwsze otwarcie zakładki | `.dn-spinner` w miejscu listy |
| Pusty rejestr | Brak jakiejkolwiek zarejestrowanej umiejętności na serwerze | `.dn-pusty-stan` z CTA „Zainstaluj pierwszą umiejętność” |
| Brak wyników filtra | Filtr/wyszukiwanie nie zwraca wierszy | `.dn-pusty-stan` z opisem aktywnego filtra |
| Błąd instalacji Personal | Przesłanie pliku nieudane | `.dn-toast--blad` z treścią przyczyny; formularz instalacji pozostaje otwarty do ponowienia |

---

## 8. Connectors Manager

### 8.1. Cel i zawartość

Connectors Manager podłącza do definicji eksperta pluginy, konektory i serwery MCP — rozszerzenia otwierające ekspertowi dostęp poza platformę (Specyfikacja agentów, rozdz. 5.5; Rozszerzenia, rozdz. 1 i 6). Podobnie jak w Skills Manager, okno rozstrzyga wyłącznie **czy rozszerzenie jest podłączone do definicji eksperta** — zakres, w jakim ekspert może je faktycznie wykorzystać (odczyt, zapis, wywołanie akcji), rozstrzyga Permissions Center (rozdz. 9).

```
Schemat — rodzaje rozszerzeń i ich okna w module Agents

  ROZSZERZENIE  (jednolity kontrakt integracji — Rozszerzenia, rozdz. 3)
  │
  ├── Umiejętność (skill) ───────────────────────────► Skills Manager (rozdz. 7)
  │     zdefiniowany sposób wykonania zadania
  │
  ├── Wtyczka ─────────────┐
  ├── Konektor              ├──────────────────────────► Connectors Manager (ten rozdział)
  └── Serwer MCP ──────────┘
        podłączenie do zewnętrznego systemu, usługi
        lub serwera zgodnego z Model Context Protocol

  Podział przyjęty w niniejszym projekcie UI podąża za oknami wskazanymi
  w Specyfikacji agentów (rozdz. 5.4–5.5): Skills Manager = umiejętności;
  Connectors Manager = wtyczki, konektory i serwery MCP łącznie — zgodnie
  z formułą zadania „skille / pluginy / konektory / MCP”. Rozszerzenia
  (rozdz. 6) grupują wtyczkę koncepcyjnie ze skillami na poziomie kontraktu
  integracji; oba ujęcia korzystają z tego samego, jednolitego kontraktu
  (Rozszerzenia, rozdz. 3) i tego samego rejestru — różni się wyłącznie
  okno, w którym Operator dokonuje przypisania do eksperta.
```

### 8.2. Makieta

```
Makieta — zakładka Connectors Manager (stan spoczynku)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW   │ ⟳ EXECUTION LOOP │ 🧑Tożsamość │🚀Model │⭐Skille
  Użytkownik ↔     │ Koordynator ↔    │ 🔗Connectors Manager │🛡Uprawn. ┬─────────
  Wykonawca        │ Wykonawca        │ [ 🔍 szukaj ] [ Rodzaj ▼ ]     │ PODSUMO-
 ──────────────────┼──────────────────┤ ────────────────────────────── │ WANIE
  strumień rozmowy │ kolejka zadań    │ 🔗 Repozytorium kodu           │ Rozszerz.
                   │ kontrola jakości │    [Konektor]      [ ▣ ]     ⋮ │ podłączone
                   │ ponowienia: 0    │ 🔗 Serwer MCP „LexAI”          │ : 2
                   │ przebieg ▮▮▯▯    │    [Serwer MCP]    [ ▣ ]     ⋮ │
                   │                  │ 🔗 Słownik terminologiczny     │
                   │                  │    [Wtyczka]       [ ☐ ]     ⋮ │
                   │                  │ ⓘ zakres wykorzystania —       │
                   │                  │   zakładka Permissions Center →│
                   │                  │ [ Operacje ▼ ]                 │
 ──────────────────┤                  │ ───────────────────────────────┴─────────
  [ Napisz polec. ]│ [ Sterowanie ▼ ] │            [ Anuluj ]  [ Zapisz eksperta ]
  [ Agent ▼ ][Wyśl]│                  │
 ══════════════════════════════════════════════════════════════════════════════
```

### 8.3. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole wyszukiwania | `.dn-input` z ikoną `szukaj` | Filtruje rejestr po nazwie | 2 | Kliknięcie ikony `szukaj` w nagłówku widoku; po użyciu pole zwija się samoczynnie | Średnie pole | Domyślny · fokus · wypełnione | Filtruje listę na żywo | Górny pasek |
| Filtr „rodzaj” | `.dn-zakladki--pigulki`: Wszystkie / Wtyczki / Konektory / Serwery MCP | Zawęża wg rodzaju rozszerzenia (rozdz. 8.1 schemat) | 3 | Grupa `Rodzaj ▼` — lista rozwijana | Segmentowy pasek | Domyślny (Wszystkie) · wybrany segment | Filtruje listę | Górny pasek |
| Filtr „Źródło” | `.dn-select`: Danaco Plugin / Personal | Zawęża wg pochodzenia | 3 | Grupa `Filtry ▼` — lista rozwijana | Mała kontrolka | Domyślny · rozwinięty | Filtruje listę | Górny pasek |
| Wiersz rozszerzenia | `.dn-karta--pozycja` z ikoną `link-zewnetrzny` | Reprezentuje jeden zarejestrowany plugin, konektor lub serwer MCP | 1 | Widoczny bez interakcji | Średni wiersz listy | Podłączony / Niepodłączony (adnotacja w wierszu) | Kliknięcie nazwy rozwija panel konfiguracji instancji | Lista główna |
| Plakietka rodzaju | `.dn-plakietka--wersaliki`: „Plugin” / „Konektor” / „Serwer MCP” | Rozróżnia rodzaj rozszerzenia w wierszu | 1 | Widoczna bez interakcji | Bardzo mała pigułka | Domyślny | — | Wiersz rozszerzenia |
| Przełącznik „Podłączony temu ekspertowi” | `.dn-suwak` w wierszu | Decyduje, czy TEN ekspert ma dostęp do danego rozszerzenia | 1 | Widoczny bez interakcji | Mały przełącznik | Domyślny (wyłączony) · włączony — zawsze w pełni klikalny, również gdy rozszerzenie jest wyłączone w rejestrze globalnym | Włączenie podłącza rozszerzenie do definicji eksperta (jeśli wciąż wyłączone w rejestrze globalnym — zapisane jako „podłączone, oczekuje aktywacji w rejestrze”, zaczyna działać automatycznie po włączeniu go tam); zakres wykorzystania pozostaje pełny dopóki nieskonfigurowany inaczej w Permissions Center (rozdz. 9.6) | Każdy wiersz |
| Panel konfiguracji instancji (dla konektora/serwera MCP) | Rozwijany fragment wiersza z polami `.dn-input` (adres serwera, dane dostępowe — domyślnie jawne, z maskowaniem jednym kliknięciem, jak w Model Configuration, rozdz. 6.3) | Parametry połączenia właściwe temu egzemplarzowi rozszerzenia (Rozszerzenia, rozdz. 4) | 3 | Kliknięcie nazwy rozszerzenia — panel rozwija się w miejscu | Średni, rozwijany | Zwinięty · rozwinięty · ostrzeżenie (adres wygląda na nieprawidłowy — sygnał, nie blokada zapisu) | Edycja zapisuje się w atrybucie Konfiguracja | Pod wierszem po kliknięciu |
| Baner „Zakres wykorzystania → Permissions Center” | `.dn-alert--info` | Jednoznacznie oddziela decyzję „czy podłączone” (tu) od decyzji „w jakim zakresie” (Permissions Center) | 1 | Widoczny bez interakcji | Pełna szerokość, pod listą | Stały | Kliknięcie „→” przełącza na zakładkę Permissions Center | Dół listy, nad stopką |
| Przycisk „Podłącz nowy plugin / konektor / serwer MCP Personal…” | `.dn-btn--zarys` z ikoną `plus` | Skrót do instalacji nowego rozszerzenia typu Personal | 3 | Grupa `Operacje ▼` w stopce listy | Mały przycisk | Domyślny · hover | Otwiera przepływ instalacji w oknie konfiguracji platformy | Stopka zakładki |

### 8.4. Stany

Analogiczne do Skills Manager (rozdz. 7.4): ładowanie rejestru, pusty rejestr, brak wyników filtra, błąd instalacji Personal, dodatkowo — błąd walidacji adresu serwera MCP (`.dn-pole--blad` w panelu konfiguracji instancji).

---

## 9. Permissions Center

### 9.1. Cel — konfiguracja możliwości, nie kontrola dostępu

**Permissions Center nie jest bramką ani mechanizmem kontroli dostępu do aplikacji.** Rozstrzyga wyłącznie **zakres operacyjny eksperta jako wykonawcy** — jak szeroko ekspert, którego Operator już uruchomił, może działać w jego imieniu (Specyfikacja agentów, rozdz. 6.1). Rozróżnienie to jest ustalone wprost w dokumentach źródłowych i musi pozostać czytelne w interfejsie:

| Zagadnienie | Kto/co odpowiada | Co rozstrzyga | Podstawa |
|---|---|---|---|
| Uwierzytelnianie | Logowanie Operatora do aplikacji | **Jedyny** mechanizm kontroli dostępu do platformy jako takiej | Architektura, rozdz. 15; Bezpieczeństwo i uwierzytelnianie, rozdz. 2 |
| **Permissions Center** | **Ten ekspert, ta definicja** | **Konfiguracja zakresu możliwości** jednostki AI działającej w imieniu już uwierzytelnionego Operatora | Specyfikacja agentów, rozdz. 6.1; Bezpieczeństwo i uwierzytelnianie, rozdz. 2 |
| Izolacja techniczna „konto i token per sesja” | Proces pojedynczej sesji | Domyślnie wyłączone zawężenie techniczne — jedna z ośmiu pozycji macierzy „Izolacja techniczna”, czyli jednej z czterech grup zakresu Permissions Center (rozdz. 9.2) | Architektura, rozdz. 12; Bezpieczeństwo i uwierzytelnianie, rozdz. 2 |

```
Schemat — trzy odrębne mechanizmy, jedno źródło nieporozumień do uniknięcia w UI

  UWIERZYTELNIANIE              PERMISSIONS CENTER            IZOLACJA TECHNICZNA
  (poza modułem Agents)          (ten dokument, rozdz. 9)       (jedna z 4 grup
                                                                  Permissions Center)
  „Czy TO URZĄDZENIE          „Jak SZEROKO może działać        „Czy TEN PROCES sesji
   ma dostęp do aplikacji      TEN EKSPERT w imieniu            ma OTRZYMAĆ WŁASNE,
   jako właściciel konta?”     już zalogowanego                 ODRĘBNE konto/token,
                                Operatora?”                      zamiast korzystać ze
                                                                  wspólnych?”
        │                              │                                │
        ▼                              ▼                                ▼
  logowanie · hasło ·           4 grupy zakresu:               ustawienie wewnątrz
  Windows Hello ·               rozszerzenia · moduły          grupy „izolacja
                                                               techniczna”
  e-mail uwierzytelniający      i zasoby · izolacja             Permissions Center —
  (rozdz. 4–8 Bezpieczeństwa)   techniczna · MultitaskingAI     domyślnie wyłączona
                                (rozdz. 9.2 niniejszego)
```

**Ta ramka jest obowiązkowym elementem interfejsu** — wyświetlana jako baner informacyjny na górze okna Permissions Center (rozdz. 9.3), aby Operator nigdy nie mylił zawężenia zakresu eksperta z blokadą dostępu do własnej aplikacji.

### 9.2. Cztery grupy zakresu

| Grupa | Przedmiot konfiguracji | Powiązane okno lub mechanizm |
|---|---|---|
| Dostęp do rozszerzeń | Które z podłączonych w Connectors Manager pluginów, konektorów i serwerów MCP ekspert może rzeczywiście wykorzystać oraz w jakim zakresie (odczyt, zapis, wywołanie akcji) | Connectors Manager (rozdz. 8) |
| Dostęp do modułów i zasobów | W których modułach platformy ekspert może działać jako wykonawca oraz do jakich zasobów modułowych (np. zapis plików w Library, wykonywanie poleceń w Terminal) ma dostęp | Macierz dostępności modułów (Koncepcja platformy, rozdz. 10) jako rama działania |
| Zakresy izolacji technicznej | Domyślne ustawienie ośmiu zakresów izolacji technicznej właściwe temu ekspertowi | Okno konfiguracji punktów izolacji (rozdz. 9.5 niniejszego dokumentu) |
| Zakres działania w MultitaskingAI | Możliwość uruchamiania Subagent Network (do 15 podagentów) oraz udział w kolejkach i orkiestracji, gdy ekspert pełni rolę wykonawczą | Panel orkiestracji, sekcja Role (rozdz. 10.3 niniejszego dokumentu) |

### 9.3. Makieta

```
Makieta — zakładka Permissions Center (stan spoczynku; funkcje warstwy 4)

 ══════════════════════════════════════════════════════════════════════════════
  💬 CHAT WINDOW   │ ⟳ EXECUTION LOOP │ 🧑Tożsamość │🚀Model │⭐Skille │🔗Konektory
  Użytkownik ↔     │ Koordynator ↔    │ 🛡Permissions Center
  Wykonawca        │ Wykonawca        │ ─────────────────────────────────────────
 ──────────────────┼──────────────────┤ ⚑ Stan wyjściowy: PEŁNY DOSTĘP OPERACYJNY.
  strumień rozmowy │ kolejka zadań    │   Centrum uprawnień nie kontroluje dostępu
                   │ kontrola jakości │   do aplikacji — konfiguruje zakres
                   │ ponowienia: 0    │   działania TEGO WYKONAWCY.
                   │ przebieg ▮▮▯▯    │              [ Resetuj do pełnego dostępu ]
                   │                  │ ──────────────┬──────────────────────────
                   │                  │ GRUPY ZAKRESU │ ZAWARTOŚĆ GRUPY
                   │                  │ ► Rozszerzenia│ Repozytorium kodu
                   │                  │   Moduły      │   ▣odczyt ▣zapis ▣akcja
                   │                  │   i zasoby    │ Serwer MCP „LexAI”
                   │                  │   Izolacja    │   ▣odczyt ▣zapis ▣akcja
                   │                  │   techniczna  │
                   │                  │   Multitasking│ [ Polityka efektywna ▼ ]
 ──────────────────┤                  │ ──────────────┴──────────────────────────
  [ Napisz polec. ]│ [ Sterowanie ▼ ] │            [ Anuluj ]  [ Zapisz eksperta ]
  [ Agent ▼ ][Wyśl]│                  │
 ══════════════════════════════════════════════════════════════════════════════
```

Grupa „Izolacja techniczna” otwiera macierz ośmiu przełączników identyczną z oknem konfiguracji punktów izolacji (System wizualny, rozdz. 12.2) — patrz rozdz. 9.5 niniejszego dokumentu.

### 9.4. Katalog elementów

| Element | Co to jest | Do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Baner „Stan wyjściowy: pełny dostęp” | `.dn-alert--info`, pełna szerokość, ikona `tarcza` | Ustala ramę odbioru całego okna — zakres, nie blokada (rozdz. 9.1) | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Duży, nieprzewijalny pasek w nagłówku zakładki | Stały, niezależny od grupy wybranej po lewej | — (statyczny, tylko przycisk „Resetuj” jest interaktywny) | Góra zakładki Permissions Center |
| Przycisk „Resetuj do pełnego dostępu” | `.dn-btn--zarys` z ikoną `odswiez` | Jednym kliknięciem przywraca stan wyjściowy we wszystkich czterech grupach | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Mały przycisk w prawym rogu banera | Domyślny — zawsze w pełni klikalny, również gdy żadne zawężenie nie zostało wprowadzone | Przywraca natychmiast pełny dostęp we wszystkich grupach, bez modala potwierdzenia; `.dn-toast--sukces` „przywrócono pełny dostęp” z akcją „Cofnij” (przywraca poprzednią konfigurację zawężeń) — przy braku uprzednich zawężeń toast „już pełny dostęp, brak zmian” | Baner górny |
| Selektor grup | Lista `.dn-karta--pozycja` bez akcji w wierszu, cztery pozycje | Wybór jednej z czterech grup zakresu do edycji | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Mała lista, lewa kolumna | Pozycja wybrana (tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt`) · pozostałe (domyślny) | Kliknięcie przełącza zawartość prawej kolumny | Lewa kolumna okna |
| Tabela „Dostęp do rozszerzeń” | `.dn-tabela` — wiersz na każde rozszerzenie podłączone w Connectors Manager | Ustala zakres wykorzystania (odczyt / zapis / wywołanie akcji) per rozszerzenie | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Średnia tabela | Domyślny (wszystkie trzy zaznaczone = pełny dostęp) · częściowo zawężony (część zaznaczeń odznaczona) | Odznaczenie zawęża zakres dla tego jednego rozszerzenia; nie wpływa na pozostałe | Prawa kolumna, grupa „Dostęp do rozszerzeń” |
| Checklista „Dostęp do modułów i zasobów” | Lista `.dn-check` — piętnaście modułów platformy + adnotacje o zasobach przykładowych (zapis plików w Library, polecenia w Terminal) | Ustala, w których modułach ekspert może działać jako wykonawca | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Średnia lista, grupowana wg środowisk | Domyślny (wszystkie zaznaczone) · częściowo odznaczone | Odznaczenie modułu wyklucza eksperta z wyboru jako wykonawcy w tym module | Prawa kolumna, grupa „Dostęp do modułów i zasobów” |
| Macierz „Izolacja techniczna” | Osiem `.dn-suwak` z ikonami (`folder`, `ustawienia`, `plik`, `globus`, `dokument`, `klodka`, `kod`, `uruchom`) i objaśnieniem `[?]` każdy | Domyślne ustawienie ośmiu zakresów izolacji technicznej właściwe temu ekspertowi (identyczna z macierzą okna konfiguracji punktów izolacji) | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Średnia macierz, identyczna forma jak System wizualny, rozdz. 12.2 | Każdy przełącznik: wyłączony (domyślnie, tor po lewej) · włączony | Włączenie pojedynczego zakresu zawęża wyłącznie ten jeden aspekt procesu tego eksperta | Prawa kolumna, grupa „Izolacja techniczna” |
| Przełącznik „Subagent Network” | `.dn-suwak` | Możliwość uruchamiania Subagent Network, gdy ekspert pełni rolę wykonawczą w MultitaskingAI | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Mały przełącznik | Domyślny (włączony) · wyłączony | Wyłączenie usuwa możliwość uruchamiania podagentów przez tego eksperta w dowolnej roli | Prawa kolumna, grupa „MultitaskingAI” |
| Pole „Liczba dozwolonych podagentów” | Suwak liczbowy / `.dn-input` numeryczne, zakres 1–15 | Ustala górną liczbę jednoczesnych podagentów (górny limit 15 jest parametrem technicznym platformy, nie decyzją produktową — Model konfiguracji, rozdz. 2) | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Mała kontrolka, aktywna tylko gdy „Subagent Network” włączony | Domyślny (15 — maksimum techniczne) · wartość zmniejszona | Zmiana zapisuje się jako część definicji eksperta | Prawa kolumna, grupa „MultitaskingAI”, pod przełącznikiem |
| Panel „Polityka efektywna dla tego eksperta” | `.dn-kod` | Podgląd wynikowego zestawu ustawień po uwzględnieniu wszystkich czterech grup, w formie identycznej z podglądem polityki efektywnej okna konfiguracji punktów izolacji (System wizualny, rozdz. 12.3) | 4 | Tryb administracyjny lub konfiguracja roli; polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji. Użytkownik podstawowy nie widzi tego elementu | Duży blok tekstowy techniczny, pod selektorem grup | Aktualizuje się na żywo przy każdej zmianie w dowolnej grupie | — (podgląd, bez interakcji poza kopiowaniem treści) | Dół okna, nad stopką |

### 9.5. Relacja do okna konfiguracji punktów izolacji

Permissions Center i okno konfiguracji punktów izolacji (dostępne z okna konfiguracji platformy, poza modułem Agents) są **dwoma punktami wejścia do tej samej macierzy izolacji technicznej** (Specyfikacja agentów, rozdz. 6.3). Ustawienia zapisane w Permissions Center pełnią funkcję wartości wyjściowej właściwej temu ekspertowi — obowiązują wszędzie, gdzie ekspert działa, dopóki nie zostaną nadpisane przez regułę zapisaną na poziomie zasięgu bardziej szczegółowym.

```
Schemat — Permissions Center wobec siedmiu poziomów zasięgu izolacji

DEFINICJA EKSPERTA                        OKNO KONFIGURACJI PUNKTÓW IZOLACJI
Permissions Center                        (poza modułem Agents)
(wartość wyjściowa
 właściwa temu ekspertowi)
        │                                          │
        │            zasięg: Globalny ─────────────┤  (najniższe pierwszeństwo)
        │            zasięg: Środowisko ────────────┤
        │            zasięg: Moduł ──────────────────┤
        │            zasięg: Para modułów ─────────────┤
        │            zasięg: Projekt ─────────────────────┤
        │            zasięg: Karta sesji ───────────────────┤
        └───────────► zasięg: Rola (MultitaskingAI) ────────────┤  (najwyższe)
                                                          ▼
                                              POLITYKA EFEKTYWNA
                                     (ten sam panel `.dn-kod`, wynik po
                                      uwzględnieniu dziedziczenia)
```

Konsekwencja projektowa: gdy ekspert pełni rolę w MultitaskingAI i tej roli przypisano własny profil izolacji (poziom zasięgu „Rola” — najwyższe pierwszeństwo, Koncepcja platformy rozdz. 6.5), profil roli nadpisuje ustawienia z Permissions Center tego eksperta **wyłącznie w kontekście tej roli** — poza nią ekspert zachowuje ustawienia zapisane we własnym Permissions Center (Specyfikacja agentów, rozdz. 8.5).

### 9.6. Stan wyjściowy — pełny dostęp operacyjny

| Grupa zakresu | Stan wyjściowy | Przykład świadomego zawężenia przez Operatora |
|---|---|---|
| Dostęp do rozszerzeń | Pełny dostęp (odczyt + zapis + wywołanie akcji) do każdego podłączonego rozszerzenia | Konektor tylko do odczytu dla eksperta w roli kontrolnej |
| Dostęp do modułów i zasobów | Pełny dostęp do wszystkich modułów, w których ekspert został wykorzystany jako wykonawca | Zawężenie do modułów Developer i Terminal dla eksperta wyspecjalizowanego w kodzie |
| Izolacja techniczna (8 zakresów) | Żaden zakres nie jest włączony — pełna operacyjna swoboda procesu | Włączenie dostępu sieciowego i zakresu odczytu/zapisu plików dla eksperta pracującego autonomicznie 24/7 |
| MultitaskingAI | Subagent Network dostępny, do 15 podagentów | Wyłączenie Subagent Network dla eksperta w roli Executor 3 / Validator |

Stan wyjściowy odpowiada zasadzie „brak ustawienia = wartość domyślna” (Koncepcja platformy, rozdz. 6.5) oraz zasadzie braku twardych blokad wbudowanych na stałe (rozdz. 4 i 12). Permissions Center nigdy nie wymusza zawężenia — udostępnia je jako możliwość (Specyfikacja agentów, rozdz. 6.4).

### 9.7. Stany szczegółowe okna

| Stan | Wyzwalacz | Co widać |
|---|---|---|
| Pełny dostęp (wyjściowy) | Nowy ekspert, żadna grupa nieskonfigurowana | Wszystkie przełączniki w pozycji domyślnej; przycisk „Resetuj” aktywny — kliknięcie potwierdza toastem „już pełny dostęp, brak zmian” |
| Częściowo zawężony | Co najmniej jedno ustawienie zmienione w dowolnej grupie | Plakietka przy nazwie grupy w selektorze: „zawężone”; panel Polityka efektywna odzwierciedla zmianę |
| Resetowanie w toku | Kliknięcie „Resetuj do pełnego dostępu” | Krótka animacja przywrócenia wszystkich przełączników; `.dn-toast--sukces` |
| Rozszerzenie usunięte z Connectors Manager po skonfigurowaniu zakresu | Operator odłącza rozszerzenie w Connectors Manager, dla którego wcześniej zawężono zakres | Wiersz znika z tabeli „Dostęp do rozszerzeń”; zapisane wcześniej zawężenie zachowane w tle na wypadek ponownego podłączenia tego samego rozszerzenia |

---

## 10. Od eksperta do Wykonawcy — zastosowanie operacyjne

**Zakres tego rozdziału.** Poniżej wyłącznie punkty styku modułu Agents z oknami innych modułów — Agent Manager (moduł Workspace) i sekcja Role (panel orkiestracji środowiska MultitaskingAI). Pełny projekt UI tych okien jest przedmiotem odrębnych dokumentów projektu UI właściwych modułowi Workspace i środowisku MultitaskingAI; niniejszy rozdział zamyka ciągłość cyklu życia eksperta rozpoczętego w rozdz. 3 oraz opisuje udział Wykonawców w pętli wykonawczej (rozdz. 10.4).

```
Diagram — trzy gałęzie wykorzystania operacyjnego eksperta (nie wykluczają się wzajemnie)

                         EKSPERT ZAPISANY W AGENT BUILDERZE
                                        │
        ┌───────────────────────────────┼────────────────────────────────┐
        ▼                                ▼                                 ▼
  10.1 Wykonawca doraźny         10.2 Agent Manager                10.3 Sekcja Role
  wybór komponentu               (moduł Workspace)                  (panel orkiestracji
  własnego w sesji                                                   MultitaskingAI)
  TalkIn / WorkSpace /            przypisanie do PROJEKTU             przypisanie do ROLI:
  CodeStudio                      — działa domyślnie w                Executor 1, Executor 2,
                                    zakresie tego projektu             Coordinator, Executor 3
                                                                        / Validator
```

### 10.1. Wykorzystanie w TalkIn, WorkSpace i CodeStudio

Ekspert nie otwiera własnego okna operacyjnego w tych trzech środowiskach — jest wybierany jako wykonawca w kontekście modułu, w którym Operator aktualnie pracuje, analogicznie do sposobu wpinania automatyki w sesji modułu Developer (Specyfikacja agentów, rozdz. 7.1).

### 10.2. Przypisanie do projektu — Agent Manager

| Aspekt | Ustalenie |
|---|---|
| Okno | Agent Manager, moduł Workspace |
| Mechanizm | Operator wskazuje w Agent Manager, którzy z zapisanych ekspertów mają zastosowanie w danym projekcie |
| Zakres domyślny | Ekspert przypisany do projektu działa domyślnie w zakresie tego projektu (instrukcje systemowe i pamięć kontekstowa projektu łączą się z instrukcjami i pamięcią eksperta) |
| Relacja do Permissions Center | Zakres możliwości ustalony w Permissions Center obowiązuje niezależnie od przypisania do projektu — projekt dodaje kontekst (instrukcje, pamięć), nie zmienia zakresu uprawnień |

### 10.3. Przypisanie do roli — sekcja Role panelu orkiestracji

| Aspekt | Ustalenie |
|---|---|
| Okno | Sekcja Role, panel orkiestracji środowiska MultitaskingAI |
| Mechanizm | Operator wskazuje dla danej roli (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) jednego z zapisanych ekspertów albo model bazowy bez pośrednictwa eksperta |
| Co niesie przypisanie | Cała definicja eksperta: tożsamość, instrukcje systemowe, umiejętności, rozszerzenia, pamięć, zakres możliwości z Permissions Center |
| Nadpisanie | Definicja jest punktem wyjścia, nadpisywalnym ustawieniami zapisanymi na poziomie zasięgu „Rola” w oknie konfiguracji punktów izolacji (rozdz. 9.5 niniejszego dokumentu) |
| Kanał modelu per rola | Kanał zapisany w Model Configuration eksperta może zostać dla danej roli nadpisany (rozdz. 6.1) |

### 10.4. Udział Wykonawców w pętli wykonawczej — odwzorowanie w interfejsie

Pełny przebieg pętli wykonawczej, role jej uczestników, kolejka i stan zadań, komunikaty sterujące, kontrola jakości oraz decyzja o ponowieniu opisane są w `specyfikacje/specyfikacja-agentow.md`, rozdz. 2.4–2.9. Koordynator nie jest agentem konfigurowanym w tym module — jest komponentem orkiestrującym platformy; agenci konfigurowani w module Agents występują wyłącznie jako Wykonawcy. Niniejszy podrozdział wskazuje, w którym elemencie interfejsu modułu każdy etap pętli jest widoczny.

| Etap pętli (`specyfikacje/specyfikacja-agentow.md`, rozdz. 2.4) | Element interfejsu modułu Agents | Warstwa |
|---|---|---|
| Dekompozycja zlecenia | Blok „Zlecenie i dekompozycja”, Execution Loop Window (rozdz. 2.4.3) | 1 |
| Przydział zadania | Kolejka i stan zadań, Execution Loop Window | 1 |
| Realizacja | Strumień komunikatów sterujących, Execution Loop Window | 1 |
| Kontrola jakości | Blok „Wynik kontroli jakości”, Execution Loop Window | 1 |
| Ponowienie | Licznik ponowień, Execution Loop Window | 1 |
| Sterowanie przebiegiem | Grupa `Sterowanie ▼`, Execution Loop Window | 3 |
| Dostrojenie parametrów pętli | Panel „Parametry pętli wykonawczej”, Execution Loop Window | 2 |
| Diagnostyka przebiegu | Podgląd surowych komunikatów sterujących i śladu wykonania | 4 |

Definicja Wykonawcy zbudowana w tym module jest tym, co Koordynator otrzymuje przy przydzieleniu zadania. Poniższa tabela wskazuje, które okno modułu ustala każdy element tej definicji; znaczenie poszczególnych komponentów definicji w samej pętli opisuje `specyfikacje/specyfikacja-agentow.md`, rozdz. 5 i 2.4.

| Komponent definicji | Okno modułu ustalające jego wartość |
|---|---|
| Tożsamość i instrukcje systemowe | Agent Builder → zakładka Tożsamość (rozdz. 5.3) |
| Model bazowy i kanał | Model Configuration (rozdz. 6); kanał nadpisywalny per rola (rozdz. 10.3) |
| Umiejętności | Skills Manager (rozdz. 7) |
| Rozszerzenia — pluginy, konektory, serwery MCP | Connectors Manager (rozdz. 8) |
| Pamięć | Agent Builder → zakładka Tożsamość (rozdz. 5.3) |
| Zakres możliwości, w tym Subagent Network | Permissions Center (rozdz. 9) |

W środowisku MultitaskingAI role Executor 1, Executor 2 i Executor 3 / Validator są rolami Wykonawców, a rola Coordinator jest miejscem, w którym Koordynator platformy działa na definicji wskazanej w sekcji Role (rozdz. 10.3; `srodowiska/multitaskingai.md`, rozdz. 3).

---

## 11. Katalog komponentów i ikon systemu wizualnego w module Agents

### 11.1. Komponenty `.dn-*` użyte w module Agents

| Komponent | Warianty użyte | Gdzie w module Agents |
|---|---|---|
| `.dn-btn` | `--zloty`, `--zarys`, `--duch`, `--blad`, `--sm` | CTA finalizujące (Zapisz), akcje drugorzędne (Anuluj, Testuj połączenie, Resetuj), akcja niebezpieczna (Usuń eksperta) |
| `.dn-btn-ikona` | — | Menu akcji karty eksperta, przyciski w wierszach rejestrów |
| `.dn-input` / `.dn-textarea` / `.dn-select` | `--blad` | Wszystkie pola formularza pięciu zakładek edytora |
| `.dn-suwak` | — | Przypisanie umiejętności/rozszerzeń, macierz izolacji technicznej, Subagent Network |
| `.dn-check` | układ radio i multi-select | Moduły zastosowania, zasięg widoczności, poziomy pamięci, checklista modułów w Permissions Center |
| `.dn-karta` | `--interaktywna`, `--akcent`, `--pozycja` | Karty ekspertów w Bibliotece, karty modeli w Model Configuration, wiersze rejestrów |
| `.dn-plakietka` | `--zloto`, `--sukces`, `--ostrz`, `--blad`, `--wersaliki`, `--stan` | Stan eksperta, wersja, widoczność, rodzaj rozszerzenia |
| `.dn-zakladki` | `--pigulki` | Pasek zakładek edytora, segment „Kanał modelu”, filtr „rodzaj” w Connectors Manager |
| `.dn-tabela` | — | Tabela „Dostęp do rozszerzeń” w Permissions Center |
| `.dn-awatar` | `--kwadrat` | Identyfikacja eksperta na kartach i w odwołaniach z sekcji Role |
| `.dn-modal` | — | Panel Historia wersji (wariant dokujący) |
| `.dn-toast` | `--sukces`, `--blad` | Potwierdzenie zapisu, błąd zapisu, błąd instalacji rozszerzenia |
| `.dn-tooltip` | — | Wszystkie objaśnienia kontekstowe `[?]` |
| `.dn-alert` | `--info`, `--blad` | Baner Permissions Center, komunikaty błędów kanału/adresu |
| `.dn-pusty-stan` | — | Pusta Biblioteka, pusty rejestr Skills/Connectors Manager |
| `.dn-kod` | — | Panel „Polityka efektywna dla tego eksperta” |
| `.dn-spinner` | — | Ładowanie list, zapisywanie, test połączenia |

### 11.2. Ikony (z zestawu 47 — System wizualny, Załącznik C)

| Ikona | Zastosowanie w module Agents |
|---|---|
| `uzytkownik` | Ikona okna Agent Builder; awatar (wariant `--kwadrat`) eksperta |
| `uruchom` | Ikona okna Model Configuration; „serwer wykonania” w macierzy izolacji |
| `gwiazdka` | Ikona okna Skills Manager; wiersze umiejętności |
| `link-zewnetrzny` | Ikona okna Connectors Manager; wiersze pluginów/konektorów/serwerów MCP |
| `tarcza` | Ikona okna Permissions Center (dosłowne dopasowanie w Załączniku C) |
| `plus` | „＋ Nowy ekspert”, „Zainstaluj/Podłącz nowe rozszerzenie” |
| `szukaj` / `filtr` | Wyszukiwanie i filtrowanie w Bibliotece i rejestrach |
| `archiwum` | Historia wersji; archiwizacja eksperta |
| `kosz` | Usunięcie eksperta (wariant `.dn-btn--blad`) |
| `odswiez` | Testuj połączenie; Przywróć wersję; Resetuj do pełnego dostępu |
| `info` | Adnotacje pomocnicze i objaśnienia kontekstowe |
| `folder`, `ustawienia`, `plik`, `globus`, `dokument`, `klodka`, `kod` (dodatkowo) | Osiem pozycji macierzy izolacji technicznej w Permissions Center — identycznie jak w oknie konfiguracji punktów izolacji |
| `wiecej` | Menu „więcej akcji” na karcie eksperta |
| `strzalka-lewo` | Powrót z edytora do Biblioteki |

---

## 12. Słowniczek pojęć

| Pojęcie | Definicja w niniejszym dokumencie |
|---|---|
| Ekspert | Określenie interfejsowe komponentu własnego rodzaju Agent (Koncepcja platformy, rozdz. 2.3, 11.15) — nazwanej, skonfigurowanej jednostki AI utworzonej w module Agents |
| Fabryka ekspertów | Metafora projektowa modułu Agents — miejsce zamieniające surowy model bazowy w nazwanego, skonfigurowanego eksperta |
| Biblioteka ekspertów | Widok wejściowy okna Agent Builder — przegląd wszystkich zapisanych ekspertów Operatora |
| Wersja eksperta | Zapisany stan definicji eksperta powstały przy każdym „Zapisz eksperta”; poprzednie wersje zachowane w panelu Historia wersji |
| Wersja aktywna | Wersja faktycznie wykorzystywana, gdy ekspert zostaje wybrany jako wykonawca |
| Podłączenie rozszerzenia | Decyzja podejmowana w Skills Manager lub Connectors Manager: czy zarejestrowane na platformie rozszerzenie jest częścią definicji tego eksperta |
| Zakres możliwości | Decyzja podejmowana w Permissions Center: w jakim stopniu podłączony ekspert może faktycznie wykorzystać swoje rozszerzenia, moduły i mechanizmy platformy |
| Polityka efektywna | Wynikowy zestaw reguł zakresu możliwości po uwzględnieniu wszystkich czterech grup Permissions Center oraz — dla eksperta w roli MultitaskingAI — dziedziczenia z poziomów zasięgu izolacji |
| Stan wyjściowy | Wartość obowiązująca w danym polu lub przełączniku, dopóki Operator świadomie jej nie zmieni; nigdy stan zablokowany |
| Użytkownik | Operator platformy — zleca zadania i zatwierdza działania; komunikuje się z Wykonawcą przez Chat Window |
| Koordynator | Komponent orkiestrujący platformy — dekomponuje zlecenie, przydziela i nadzoruje zadania, prowadzi kontrolę jakości i decyduje o ponowieniach; komunikuje się z Wykonawcą przez Execution Loop Window |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania — ekspert utworzony w module Agents |
| Chat Window | Główne okno komunikacji, kanał Użytkownik ↔ Wykonawca; lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Execution Loop Window | Okno pętli wykonawczej, kanał Koordynator ↔ Wykonawca; kolumna sąsiadująca z Chat Window |
| Warstwa widoczności | Jedna z czterech warstw ujawniania funkcjonalności interfejsu (rozdz. 2.5), przypisana każdemu elementowi okna |

---

## 13. Punkty sterowania z okna konfiguracji

Wszystkie poniższe punkty Operator personalizuje z okna konfiguracji platformy; wartości domyślne oznaczają działanie bez blokad, klucze i integracje pozostają jawne. Brak ustawienia oznacza wartość domyślną, nigdy brak dostępności funkcji.

| Grupa | Punkt sterowania | Zakres personalizacji | Domyślnie |
|---|---|---|---|
| **Okna modułu** | Zakładka otwierana przy wejściu do edytora | Tożsamość, Model, Skille, Konektory, Uprawnienia | Tożsamość |
| | Widok wejściowy okna Agent Builder | Biblioteka ekspertów albo ostatnio edytowany ekspert | Biblioteka ekspertów |
| | Otwarcie i szerokość kolumny Execution Loop Window | Kolumna rozwinięta przy wejściu do modułu albo zwinięta do znacznika; szerokość kolumny | Rozwinięta, szerokość standardowa |
| | Szerokość kolumny Chat Window | Szerokość lewej kolumny okna komunikacji Użytkownik ↔ Wykonawca | Szerokość standardowa kolumny komunikacji |
| **Definicja eksperta** | Domyślny model bazowy nowego eksperta | Wskazanie modelu podstawianego przy tworzeniu nowej definicji | Brak — wybór przy pierwszym zapisie |
| | Domyślny kanał modelu | API, CLI, SSH, HTTP | API |
| | Domyślne poziomy pamięci | Globalna, projekt, sesja, środowisko, wyłączona | Sesja |
| | Domyślna widoczność eksperta | Globalny albo projektowy | Globalny |
| | Automatyczne wersjonowanie | Nowa wersja przy każdym zapisie albo wersja na żądanie | Nowa wersja przy każdym zapisie |
| | Retencja historii wersji | Liczba przechowywanych wersji definicji | Bez limitu |
| **Rozszerzenia** | Rejestr umiejętności i rozszerzeń | Które rejestry są przeszukiwane w Skills Manager i Connectors Manager | Wszystkie zarejestrowane |
| | Automatyczne podłączanie rozszerzeń | Czy nowo zarejestrowane rozszerzenie jest podłączane do nowych ekspertów | Wyłączone (podłączenie jawne) |
| | Test połączenia konektora | Wykonanie testu przy podłączeniu rozszerzenia | Włączony |
| **Zakres możliwości** | Stan wyjściowy Permissions Center | Pełny dostęp albo zestaw zawężeń przyjęty jako punkt wyjścia nowych ekspertów | Pełny dostęp (rozdz. 9.6) |
| | Osiem zakresów izolacji technicznej | Wartość wyjściowa właściwa nowemu ekspertowi (rozdz. 9.4) | Żaden zakres nieaktywny |
| | Subagent Network | Dostępność mechanizmu i górna liczba jednoczesnych podagentów (1–15) | Dostępny, 15 |
| **Pętla wykonawcza** | Widoczność Execution Loop Window | Otwarcie okna przy starcie zlecenia albo na żądanie | Otwarcie przy starcie zlecenia |
| | Próg ponowienia zadania | Liczba ponowień zadania po negatywnym wyniku kontroli jakości | 2 |
| | Zakres treści Execution Loop Window | Które elementy pętli są prezentowane: dekompozycja zlecenia, kolejka i stan zadań, komunikaty sterujące, wynik kontroli jakości, wskaźniki przebiegu | Wszystkie elementy widoczne |
| **Warstwy widoczności** | Ujawnianie elementów warstw 2–4 | Które elementy warstw 2–4 są ujawniane w interfejsie danej roli użytkownika (rozdz. 2.5) | Warstwa 1 widoczna, warstwy 2–4 zwinięte |
| **Powiązania modułów** | Agents ───► Workspace (Agent Manager) | Przypisywanie ekspertów do projektów | Włączone |
| | Agents ───► MultitaskingAI (sekcja Role) | Przypisywanie ekspertów do ról zespołu | Włączone |
| | Agents ◄──► Automations | Udostępnianie ekspertów automatykom cyklicznym | Wyłączone (włączenie jawne) |
| **Potwierdzenia i cofanie** | Potwierdzenia akcji nieodwracalnych | Modal przed usunięciem eksperta albo wersji definicji | Wyłączone (akcja od razu + „Cofnij”) |

---

## Załącznik A. Szablon definicji eksperta jako formularz UI

Poniższy szablon przekłada szablon definicji agenta (Specyfikacja agentów, rozdz. 10) na konkretne kontrolki zaprojektowane w niniejszym dokumencie.

```
EKSPERT                                                          okno / zakładka
  nazwa:                 [ .dn-input, nagłówek edytora ]         Agent Builder
  przeznaczenie:         [ .dn-textarea małe ]                   Agent Builder
                                                                  → zakładka Tożsamość

  model_bazowy:          [ siatka kart .dn-karta--interaktywna ] Model Configuration
  kanał_modelu:          [ .dn-zakladki--pigulki: API|CLI|SSH|HTTP ]
                                                                  Model Configuration

  instrukcje_systemowe:  [ .dn-textarea duże ]                   Agent Builder
                                                                  → zakładka Tożsamość

  umiejętności:          [ .dn-tabela / .dn-karta--pozycja +      Skills Manager
                            .dn-suwak per wiersz ]
  rozszerzenia:          [ .dn-tabela / .dn-karta--pozycja +      Connectors Manager
                            .dn-suwak per wiersz; plakietka
                            rodzaju: Plugin|Konektor|Serwer MCP ]

  pamięć:                [ .dn-suwak ×5: globalna|projekt|sesja|  Agent Builder
                            środowisko|wyłączona ]                → zakładka Tożsamość

  zakres_możliwości:                                              Permissions Center
    rozszerzenia:        [ .dn-tabela: ▣odczyt ▣zapis ▣akcja      (stan wyjściowy:
                            per podłączone rozszerzenie ]           pełny dostęp)
    moduły_i_zasoby:     [ .dn-check ×15 modułów ]
    izolacja_techniczna: [ .dn-suwak ×8 — identyczne z oknem
                            konfiguracji punktów izolacji ]
    multitaskingai:
      subagent_network:  [ .dn-suwak + pole liczbowe 1–15 ]

  moduły_zastosowania:   [ .dn-check ×3: TalkIn|WorkSpace|        Agent Builder
                            CodeStudio ]                          → zakładka Tożsamość
  widoczność:            [ .dn-check radio: globalny|projektowy ]

  wersja:                [ generowana automatycznie przy         Agent Builder
                            każdym zapisie; panel Historia         → panel Historia
                            wersji ]                                wersji
  powiązania:            [ tylko do odczytu — liczniki na         Agent Builder
                            karcie; edycja w Agent Manager i       (Biblioteka),
                            sekcji Role, poza modułem Agents ]     poza zakresem
                                                                  tego dokumentu
```

---

*Koniec dokumentu. Moduł Agents — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*

*© 2026 Danaco Holding Group Sp. z o.o. — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
