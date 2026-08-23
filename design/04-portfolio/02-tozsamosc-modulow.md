# Danaco Console — Portfolio Design Identity — P2: Tożsamość piętnastu modułów

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Produkt (warstwa funkcjonalna)** | Danaco Pilot — Platforma AI Workspace OS (dokumentacja projektowa v1.0) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Dokument** | Plansza portfolio P2 — tożsamość piętnastu modułów platformy (opracowanie towarzyszące planszy `02-tozsamosc-modulow.html`) |
| **Odbiorcy** | Designer (tożsamość wizualna modułu, ikona, waga, miejsce w nawigacji) · Deweloper (kompletny inwentarz okien, macierz dostępności) · Odbiorca portfolio (co platforma faktycznie zawiera) |
| **Zakres** | Karta tożsamości każdego z 15 modułów · macierz moduł × środowisko · pełny inwentarz 77 wystąpień okien operacyjnych · grupowanie wg charakteru pracy · trzy przypadki szczególne klasyfikacji · mapa powiązań międzymodułowych |
| **Czego NIE zawiera** | Specyfikacji wnętrza okien (to zakres `dok/projekt-ui/moduly/*.md` rozdz. 3 oraz prototypów w `WYNIK/05-okna/moduly/`) · nowych ustaleń projektowych · żadnego modułu, okna ani nazwy bez pokrycia w dokumentacji |
| **Zasada nadrzędna** | Każda nazwa, liczba i przynależność na planszy pochodzi z dokumentacji merytorycznej albo z odczytu realnego pliku w repozytorium wynikowym. Zero elementów ilustracyjnych. |

---

## Spis treści

1. [Czym jest tożsamość modułu](#1-czym-jest-tożsamość-modułu)
2. [Piętnaście modułów — tabela zbiorcza](#2-piętnaście-modułów--tabela-zbiorcza)
3. [Karty tożsamości modułów](#3-karty-tożsamości-modułów)
4. [Grupowanie wg charakteru pracy](#4-grupowanie-wg-charakteru-pracy)
5. [Inwentarz okien w liczbach](#5-inwentarz-okien-w-liczbach)
6. [Macierz moduł × środowisko](#6-macierz-moduł--środowisko)
7. [Trzy przypadki szczególne](#7-trzy-przypadki-szczególne)
8. [Mapa powiązań między modułami](#8-mapa-powiązań-między-modułami)
9. [Decyzje projektowe planszy](#9-decyzje-projektowe-planszy)
10. [Źródła](#10-źródła)

---

## 1. Czym jest tożsamość modułu

Moduł jest środkowym stopniem trójstopniowej hierarchii platformy: **Środowisko** („w jakim trybie pracuję?”) → **Moduł** („jakie zadanie wykonuję?”) → **Okno operacyjne** („jakim narzędziem realizuję?”) — KANON rozdz. 6; `projekt-ui/README.md` rozdz. 2.

Na tożsamość pojedynczego modułu w warstwie wizualnej składa się siedem rozstrzygalnych cech. Wszystkie siedem jest udokumentowanych — żadnej nie wymyślono na potrzeby planszy:

| # | Cecha tożsamości | Gdzie rozstrzygnięta |
|---|---|---|
| 1 | **Ikona modułu** — jeden znak z zestawu, siatka 24×24, obrys 1,75, `currentColor` | `zasoby/ikony/manifest.json`, pole `zastosowanie` (przypisanie ikona → moduł jest tam jawne) |
| 2 | **Nazwa własna** — po angielsku, bez tłumaczenia i parafrazy | KANON rozdz. 10 pkt 1–2 |
| 3 | **Istota** — jedno zdanie odpowiadające na pytanie „jakie zadanie wykonuję?” | `projekt-ui/README.md` rozdz. 3.5, kolumna „Istota” |
| 4 | **Środowiska dostępności** — w których powłokach moduł ma pozycję w bocznej nawigacji | KANON rozdz. 6 (macierz) · `README.md` rozdz. 5 · nagłówek pliku modułu |
| 5 | **Komplet okien operacyjnych** z jawnie wskazanym **oknem wiodącym** (punkt wejścia) | KANON rozdz. 7.5 · `README.md` rozdz. 4 · rozdz. 2 pliku modułu |
| 6 | **Podwójna rola** — czy moduł jest jednocześnie komponentem własnym strefy 2 | KANON rozdz. 6 · `README.md` rozdz. 3.5, kolumna „Komp.” · L9 |
| 7 | **Powiązania z innymi modułami** — jawne, konfigurowalne, nigdy domyślnie aktywne | rozdz. 5.4 albo 6.3 pliku modułu |

Ósma cecha jest wspólna wszystkim piętnastu i dlatego nie różnicuje: **Chat Window**. Pas komunikacji jest oknem wspólnym każdego modułu, rekonfigurowanym w jego kontekście (KANON rozdz. 7.4; `przeplyw/elementy-okien.md` rozdz. 3.6).

---

## 2. Piętnaście modułów — tabela zbiorcza

Okno wiodące (punkt wejścia) zapisane **pogrubieniem**. Kolumna „Liczba” obejmuje Chat Window. Kolumna „Komp.” — czy moduł jest jednocześnie komponentem własnym strefy 2 strony głównej.

| Moduł | Ikona | Okna operacyjne | Liczba | Środowiska | Komp. |
|---|---|---|---:|---|---|
| Studio | `dokument` | **Studio Editor**, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window | 6 | TalkIn, WorkSpace | — |
| Research | `badanie` | **Research Workspace**, Sources Manager, Findings Panel, Report Builder, Export Panel | 6 | TalkIn, WorkSpace | — |
| Library | `biblioteka` | **Library Explorer**, Tags & Collections, File Preview, Versioning Panel | 5 | TalkIn, WorkSpace | — |
| Translate | `tlumacz` | Source Panel, **Translation Panels**, Glossary Manager | 4 | TalkIn | — |
| Browser | `karta-okna` | **Browser Window**, Sources Panel, Notes Panel | 4 | TalkIn, WorkSpace | — |
| Assistant | `mikrofon` | **Voice Console**, Actions Monitor, Activity Feed | 4 | TalkIn | **Tak** — Profil asystenta |
| Roundtable | `debata` | **Model Panels**, Debate Panel, Moderator Panel, Consensus Panel | 5 | TalkIn, WorkSpace, CodeStudio | — |
| Workspace | `warstwy` | **Project Dashboard**, Instructions Panel, Context Memory, Project Library, Agent Manager | 6 | TalkIn, WorkSpace, CodeStudio | **Tak** — Projekt |
| Automations | `automatyzacja` | **Workflow Builder**, Scheduler, Queue Manager, Orchestrator, Execution Monitor | 6 | brak okna w bocznej nawigacji | **Tak** — Automatyka |
| Design | `paleta` | **Design Board**, Assets Panel, Prompt Builder, Preview Window | 5 | WorkSpace, CodeStudio | — |
| Apps | `aplikacje` | **Product Builder**, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel | 6 | WorkSpace, CodeStudio | — |
| Terminal | `terminal` | **Terminal Tabs**, Output Console, Process Monitor | 4 | CodeStudio | — |
| Developer | `kod` | **Code Editor**, Project Tree, Git Panel, Build Output | 5 | CodeStudio | — |
| Diagnostics | `diagnostyka` | **Diagnostics Center**, Logs Viewer, Errors Panel, Recommendations Panel | 5 | CodeStudio | — |
| Agents | `agenci` | **Agent Builder**, Model Configuration, Skills Manager, Connectors Manager, Permissions Center | 6 | TalkIn, WorkSpace, CodeStudio; rola/ekspert w MultitaskingAI | **Tak** — Agent (ekspert) |

**Suma:** 62 okna właściwe modułom + 15 wystąpień Chat Window = **77 wystąpień okien operacyjnych**; okien odrębnych: 62 + 1 = **63**.

---

## 3. Karty tożsamości modułów

Każda karta podaje: ikonę z zestawu, istotę wg `README.md` rozdz. 3.5, środowiska dostępności, komplet okien z oknem wiodącym, podwójną rolę oraz udokumentowane powiązania. Odnośnik prowadzi do klikalnego prototypu okien modułu w `WYNIK/05-okna/moduly/`.

### 3.1. Studio

- **Ikona:** `dokument` (lucide:file-text — „Dokument, opracowanie — moduł Studio”)
- **Istota:** zaawansowana praca z tekstem: edytor + operacje AI na fragmencie + historia wersji w jednym oknie.
- **Definicja źródłowa:** „Studio jest modułem zaawansowanej pracy z tekstem, dokumentami i treścią… Moduł nie jest edytorem tekstu z dopiętym czatem — jest jedną przestrzenią, w której edytor, operacje kontekstowe AI i historia wersji współdzielą ten sam dokument w czasie rzeczywistym.” (`moduly/studio.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace
- **Okna (6):** Chat Window · **Studio Editor** (wiodące) · Tools Panel · Diff/Grep Panel · Session Repository · Preview Window
- **Komponent własny strefy 2:** nie
- **Powiązania:** Translate ↔ Studio · Studio → Library · Studio ↔ Research · Design → Studio (`studio.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/studio.html`

### 3.2. Research

- **Ikona:** `badanie` (lucide:telescope — „Moduł Research — analizy i raporty”)
- **Istota:** badania i opracowania: źródła → ustalenia cząstkowe → raport końcowy w ciągłym procesie.
- **Definicja źródłowa:** „Research jest modułem realizacji badań, analiz i opracowań… Moduł nadaje strukturę pracy badawczej, która w pojedynczej rozmowie z modelem pozostałaby płaska i trudna do prześledzenia wstecz.” (`moduly/research.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace
- **Okna (6):** Chat Window · **Research Workspace** (wiodące) · Sources Manager · Findings Panel · Report Builder · Export Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Browser ↔ Research · Library ↔ Research · Research ↔ Studio · Research ↔ Roundtable (`research.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/research.html`

### 3.3. Library

- **Ikona:** `biblioteka` (lucide:library-big — „Moduł Library — repozytorium wiedzy”)
- **Istota:** centralne repozytorium wiedzy i plików — trwała pamięć zewnętrzna wspólna wielu modułom.
- **Definicja źródłowa:** „Library jest centralnym repozytorium wiedzy i plików platformy… Pełni funkcję trwałej pamięci zewnętrznej, wspólnej dla wielu środowisk i modułów.” (`moduly/library.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace
- **Okna (5):** Chat Window · **Library Explorer** (wiodące) · Tags & Collections · File Preview · Versioning Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Studio → Library · Research → Library · Browser → Library · Library ↔ Project Library (moduł Workspace) (`library.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/library.html`

### 3.4. Translate

- **Ikona:** `tlumacz` (lucide:languages — „Moduł Translate”)
- **Istota:** tłumaczenie wielojęzyczne równoległe — jeden tekst źródłowy, wiele paneli docelowych naraz.
- **Definicja źródłowa:** „Translate jest modułem wielojęzycznych tłumaczeń w trybie wielozadaniowym — przestrzenią, w której jeden tekst źródłowy tłumaczony jest jednocześnie na wiele języków docelowych, widocznych obok siebie w osobnych panelach.” (`moduly/translate.md` rozdz. 1.1)
- **Środowiska:** TalkIn (jedyny moduł wyłącznie jednośrodowiskowy razem z Assistant)
- **Okna (4):** Chat Window · Source Panel · **Translation Panels** (wiodące) · Glossary Manager
- **Komponent własny strefy 2:** nie
- **Powiązania:** Studio ↔ Translate — „Translate jest modułem o najwęższym zbiorze powiązań konfigurowalnych na platformie” (`translate.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/translate.html`

### 3.5. Browser

- **Ikona:** `karta-okna` (lucide:app-window — „Okno, karta przeglądarki, moduł Browser”)
- **Istota:** współdzielone przeglądanie — użytkownik i AI patrzą na tę samą stronę w tym samym czasie.
- **Definicja źródłowa:** „Wspólny podgląd jest istotą modułu: to, co widzi AI, jest dokładnie tym, co widzi użytkownik, bez pośredniczącego opisu.” (`moduly/browser.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace
- **Okna (4):** Chat Window · **Browser Window** (wiodące) · Sources Panel · Notes Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Browser ↔ Research · Browser → Library (`browser.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/browser.html`

### 3.6. Assistant

- **Ikona:** `mikrofon` (lucide:mic — „Moduł Assistant — interfejs głosowy”)
- **Istota:** komunikacja głosowa z AI; jednocześnie komponent własny (Profil asystenta) i moduł operacyjny.
- **Definicja źródłowa:** „W odróżnieniu od pozostałych czternastu modułów platformy, Assistant jest jednocześnie komponentem własnym — Profilem asystenta, konfigurowanym w strefie 2 strony głównej — i modułem operacyjnym osadzonym w środowisku TalkIn.” (`moduly/assistant.md` rozdz. 1.1)
- **Środowiska:** TalkIn (nagłówek pliku używa formuły „Środowisko wykorzystania operacyjnego”)
- **Okna (4):** Chat Window (uzupełniający głos) · **Voice Console** (wiodące) · Actions Monitor · Activity Feed
- **Komponent własny strefy 2:** **tak — Profil asystenta**
- **Powiązania:** Assistant ═ Always On Display (pokrewieństwo funkcjonalne, nie łączenie kontekstu) · dowolny moduł jako krok realizacji polecenia — zdolność operacyjna profilu, nie stałe powiązanie modułowe (`assistant.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/assistant.html`

### 3.7. Roundtable

- **Ikona:** `debata` (lucide:messages-square — „Moduł Roundtable — debata modeli”)
- **Istota:** współpraca wielu modeli nad jednym problemem: równoległe odpowiedzi → krytyka → konsensus.
- **Definicja źródłowa:** „…wartość nie wynika z pojedynczej, izolowanej odpowiedzi jednego modelu, lecz z konfrontacji różnych perspektyw, wzajemnej krytyki i wypracowania wspólnego stanowiska.” (`moduly/roundtable.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace, CodeStudio
- **Okna (5):** Chat Window · **Model Panels** (wiodące) · Debate Panel · Moderator Panel · Consensus Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Roundtable ═ MultitaskingAI (pokrewieństwo funkcjonalne — „ta sama idea wielomodelowości”) · Research → Roundtable · Roundtable → Studio · Roundtable → Research (`roundtable.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/roundtable.html`

### 3.8. Workspace

- **Ikona:** `warstwy` (lucide:layers — „Moduł Workspace — przestrzeń projektowa”)
- **Istota:** izolowana przestrzeń projektowa: rozdzielenie kontekstu, instrukcji i pamięci między projektami.
- **Definicja źródłowa:** „Moduł adresuje sytuację, w której jedna, wspólna historia rozmowy z AI i jedna, wspólna pamięć przestają wystarczać, bo ustalenia jednego przedsięwzięcia zaczynają przenikać do innego.” (`moduly/workspace.md` rozdz. 1.1)
- **Środowiska:** TalkIn, WorkSpace, CodeStudio
- **Okna (6):** Chat Window · **Project Dashboard** (wiodące) · Instructions Panel · Context Memory · Project Library · Agent Manager
- **Komponent własny strefy 2:** **tak — Projekt**
- **Powiązania:** Agents → Workspace (Agent Manager) · Library ↔ Project Library · Workspace → Automations · przypisanie profilu izolacji na poziomie „Projekt” (`workspace.md` rozdz. 6.3)
- **Prototyp:** `05-okna/moduly/workspace.html`

### 3.9. Automations

- **Ikona:** `automatyzacja` (lucide:workflow — „Moduł Automations — procesy i kolejki”)
- **Istota:** automatyzacja powtarzalnych procesów bez nadzoru; jedyny moduł bez okna w bocznej nawigacji.
- **Definicja źródłowa:** „Moduł oddziela pracę konwersacyjną, jednorazową (typową dla pozostałych czternastu modułów) od pracy proceduralnej, uruchamianej według harmonogramu lub zdarzenia.” (`moduly/automations.md` rozdz. 1.2)
- **Środowiska:** brak okna modułowego w jakimkolwiek środowisku — wejście wyłącznie ze strony głównej, strefa 2; opcjonalna integracja 24/7 z MultitaskingAI
- **Okna (6):** Chat Window · **Workflow Builder** (wiodące) · Scheduler · Queue Manager · Orchestrator · Execution Monitor
- **Komponent własny strefy 2:** **tak — Automatyka**
- **Powiązania:** Automations → MultitaskingAI (harmonogramy i kolejki pętli pracy ciągłej) · Automations → dowolny moduł platformy (krok „akcja w module” w Workflow Builder) (`automations.md` rozdz. 6.3)
- **Prototyp:** `05-okna/moduly/automations.html`

### 3.10. Design

- **Ikona:** `paleta` (lucide:palette — „Moduł Design — materiały wizualne”)
- **Istota:** generowanie i edycja grafiki przez AI, zestawianie w kompozycję.
- **Definicja źródłowa:** „Design jest przestrzenią dla użytkowników generujących i edytujących materiały graficzne — od pojedynczych ilustracji i elementów brandingowych po kompletne makiety interfejsu.” (`moduly/design.md` rozdz. 1.1)
- **Środowiska:** WorkSpace, CodeStudio (moduł niedostępny w TalkIn — zapis jawny w nagłówku pliku)
- **Okna (5):** Chat Window · **Design Board** (wiodące) · Assets Panel · Prompt Builder · Preview Window
- **Komponent własny strefy 2:** nie
- **Powiązania:** Design → Studio (zasoby przy redagowaniu dokumentów) · Design → Apps (zasoby dla Frontend Workspace) (`design.md` rozdz. 6.3)
- **Prototyp:** `05-okna/moduly/design.html`

### 3.11. Apps

- **Ikona:** `aplikacje` (lucide:layout-grid — „Moduł Apps — budowa produktów”)
- **Istota:** pełne projekty aplikacyjne: architektura → frontend/backend → wdrożenie; integruje Developer i Terminal.
- **Definicja źródłowa:** „Moduł spina w jeden, ciągły proces wszystkie etapy powstawania produktu cyfrowego, integrując funkcjonalności modułów Developer i Terminal w ramach jednej, ustrukturyzowanej przestrzeni.” (`moduly/apps.md` rozdz. 1.2)
- **Środowiska:** WorkSpace, CodeStudio
- **Okna (6):** Chat Window · **Product Builder** (wiodące) · Architecture Designer · Frontend Workspace · Backend Workspace · Deployment Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Developer → Apps · Terminal → Apps · Design → Apps · Apps → MultitaskingAI (Executor 1/2 przy pracy równoległej) · Apps → Automations (Deployment Panel) (`apps.md` rozdz. 6.3)
- **Prototyp:** `05-okna/moduly/apps.html`

### 3.12. Terminal

- **Ikona:** `terminal` (lucide:square-terminal — „Moduł Terminal — powłoki wykonawcze”)
- **Istota:** konsole i środowiska wykonawcze — warstwa wykonawcza CodeStudio (powłoki, wynik, procesy).
- **Definicja źródłowa:** „Terminal jest warstwą wykonawczą środowiska CodeStudio: to w nim faktycznie uruchamia się to, co Developer i Diagnostics jedynie planują lub analizują.” (`moduly/terminal.md` rozdz. 1.1)
- **Środowiska:** CodeStudio
- **Okna (4):** Chat Window · **Terminal Tabs** (wiodące) · Output Console · Process Monitor
- **Komponent własny strefy 2:** nie
- **Powiązania:** Developer → Terminal · Diagnostics ↔ Terminal (`terminal.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/terminal.html`

### 3.13. Developer

- **Ikona:** `kod` (lucide:code-xml — „Fragment kodu, moduł Developer”)
- **Istota:** tworzenie i rozwój kodu: edytor + drzewo projektu + Git + wynik budowania ze wsparciem AI.
- **Definicja źródłowa:** „…jedną przestrzenią roboczą, w której struktura projektu, treść kodu, historia zmian i wynik budowania współdzielą ten sam kontekst przekazywany modelowi.” (`moduly/developer.md` rozdz. 1.1)
- **Środowiska:** CodeStudio
- **Okna (5):** Chat Window · **Code Editor** (wiodące) · Project Tree · Git Panel · Build Output
- **Komponent własny strefy 2:** nie
- **Powiązania:** Developer → Terminal · Developer ↔ Diagnostics · Developer → Apps (`developer.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/developer.html`

### 3.14. Diagnostics

- **Ikona:** `diagnostyka` (lucide:stethoscope — „Moduł Diagnostics — analiza problemów”)
- **Istota:** analiza i usuwanie problemów: logi → agregacja stanu → rekomendacje naprawcze, zawsze do przeglądu Operatora.
- **Definicja źródłowa:** „…punktem wyjścia nie jest tworzenie nowej funkcjonalności, lecz zrozumienie przyczyny istniejącego błędu lub spadku wydajności na podstawie logów, komunikatów błędów i obserwowalnych objawów.” (`moduly/diagnostics.md` rozdz. 1.1)
- **Środowiska:** CodeStudio
- **Okna (5):** Chat Window · **Diagnostics Center** (wiodące) · Logs Viewer · Errors Panel · Recommendations Panel
- **Komponent własny strefy 2:** nie
- **Powiązania:** Diagnostics ↔ Developer (Recommendations Panel → Code Editor) · Diagnostics ↔ Terminal (`diagnostics.md` rozdz. 5.4)
- **Prototyp:** `05-okna/moduly/diagnostics.html`

### 3.15. Agents

- **Ikona:** `agenci` (lucide:users — „Zespół, moduł Agents / role”)
- **Istota:** fabryka ekspertów — model bazowy plus tożsamość, umiejętności, konektory i zakres możliwości daje nazwaną, skonfigurowaną jednostkę AI.
- **Definicja źródłowa:** „Agents jest fabryką, która zamienia surowy model bazowy w nazwaną, skonfigurowaną jednostkę AI.” (`okna/agents.md`, nota terminologiczna; rozdz. 1.1 — łańcuch wartości)
- **Środowiska:** TalkIn, WorkSpace, CodeStudio (wykonawca) oraz MultitaskingAI (rola: Executor 1/2, Coordinator, Executor 3/Validator)
- **Okna (6):** Chat Window · **Agent Builder** (wiodące, okno nadrzędne) · Model Configuration · Skills Manager · Connectors Manager · Permissions Center
- **Komponent własny strefy 2:** **tak — Agent (ekspert)**; kafel „Agents” w strefie 2, drugi punkt wejścia to boczna nawigacja trzech środowisk (`agents.md` rozdz. 1.2)
- **Powiązania:** Agents → Workspace (Agent Manager) · Agents → wszystkie środowiska i moduły jako wykonawca · Agents → MultitaskingAI (role zespołu)
- **Prototyp:** `05-okna/moduly/agents.html`

---

## 4. Grupowanie wg charakteru pracy

Dokumentacja **nie definiuje kategorii modułów** — nie ma w niej rozdziału „rodziny modułów”. Grupowanie poniżej jest warstwą redakcyjną planszy i zostało wyprowadzone z dwóch źródeł jednocześnie, aby nie było kategorią wymyśloną:

1. **Opisu produktu powtarzanego w nagłówku każdego dokumentu źródłowego:** „…system operacyjny dla sztucznej inteligencji, integrujący **komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania**”.
2. **Udokumentowanych powiązań** modułów (rozdz. 5.4 / 6.3) — moduły trafiają do jednej grupy wtedy, gdy dokumentacja opisuje między nimi bezpośrednie powiązanie.

| Grupa | Moduły | Kotwica w opisie produktu | Uzasadnienie z powiązań |
|---|---|---|---|
| **Praca z treścią** | Studio, Translate, Design | „tworzenie treści”, „projektowanie” | `studio.md` 5.4 wiąże wszystkie trzy: Translate ↔ Studio ← Design |
| **Badanie i wiedza** | Research, Browser, Library | „zarządzanie wiedzą” | `browser.md` 5.4: Browser ↔ Research → Library; `library.md` 5.4 pokazuje Library jako cel wszystkich trzech |
| **Współpraca modeli** | Roundtable, Assistant | „komunikacja” | jedyne dwa moduły, których dokumentacja opisuje powiązanie typu **pokrewieństwo funkcjonalne** z warstwą platformy (Roundtable ═ MultitaskingAI, Assistant ═ Always On Display), a nie powiązanie konfiguracyjne z innym modułem |
| **Wytwarzanie** | Workspace, Developer, Apps | „rozwój oprogramowania” | `developer.md` 5.4: Developer → Apps; `apps.md` 1.2: Apps integruje Developer i Terminal; `workspace.md` 6.3: projekt jako jednostka organizacji pracy nad wytworem, z Agent Manager i Project Library |
| **Warstwa wykonawcza** | Terminal, Diagnostics, Automations | „automatyzacje procesów” | `terminal.md` 1.1: „Terminal jest warstwą wykonawczą środowiska CodeStudio”; Diagnostics ↔ Terminal; Automations = praca proceduralna uruchamiana wg harmonogramu lub zdarzenia |
| **Fabryka ekspertów** | Agents | — (warstwa ponad modułami) | `agents.md` rozdz. 1: „Agents jako fabryka ekspertów — pozycja w platformie”; jedyny komponent własny o zastosowaniu nieograniczonym do jednego kontekstu pracy (rozdz. 1.2) |

**Kontrola sumy:** 3 + 3 + 2 + 3 + 3 + 1 = **15**.

### 4.1. Co grupowanie ujawnia

| Obserwacja | Podstawa |
|---|---|
| Grupa „badanie i wiedza” jest jedyną domkniętą — wszystkie jej powiązania biegną wewnątrz grupy plus do Studio | `research.md`, `browser.md`, `library.md` rozdz. 5.4 |
| Grupa „warstwa wykonawcza” jest jedyną, w której żaden moduł nie wytwarza artefaktu treści — przedmiotem pracy jest bieg procesu | `terminal.md`, `diagnostics.md`, `automations.md` rozdz. 1 |
| Grupa „fabryka ekspertów” jest jednoelementowa i zasila wszystkie pozostałe grupy | `agents.md` rozdz. 1.2–1.3 |
| Trzy z sześciu grup mają moduł o podwójnej roli (Assistant, Workspace, Automations) — czwarty komponent własny, Agent, tworzy grupę samodzielnie | KANON rozdz. 6; `agents.md` rozdz. 1.2 |

---

## 5. Inwentarz okien w liczbach

| Miara | Wartość | Skąd |
|---|---:|---|
| Moduły platformy | 15 | KANON rozdz. 6 |
| Pliki dokumentacji modułowej w `moduly/` | 14 | `README.md` rozdz. 3.5 (Agents w `okna/`) |
| Okna właściwe modułom (bez Chat Window) | 62 | suma kolumn KANON rozdz. 7.5 |
| Wystąpienia Chat Window | 15 | jedno na moduł, KANON rozdz. 7.4 |
| Wystąpienia okien operacyjnych łącznie | 77 | 62 + 15 |
| Okna odrębne (Chat Window liczone raz) | 63 | 62 + 1 |
| Okna wiodące (punkty wejścia) | 15 | jedno na moduł, KANON rozdz. 7.5 |
| Moduły o podwójnej roli wg `README.md` 3.5 | 3 | Assistant, Workspace, Automations |
| Kafle komponentów własnych w strefie 2 | 4 | Automations, Agents, Workspace, Assistant — KANON rozdz. 6 |

### 5.1. Rozkład liczby okien

| Liczba okien (z Chat Window) | Moduły | Liczba modułów |
|---:|---|---:|
| 6 | Studio, Research, Workspace, Automations, Apps, Agents | 6 |
| 5 | Library, Roundtable, Design, Developer, Diagnostics | 5 |
| 4 | Translate, Browser, Assistant, Terminal | 4 |

Rozkład jest trójstopniowy i pokrywa się z wagą wizualną modułu: moduły sześciookienne prowadzą proces wieloetapowy (badanie, projekt, produkt, ekspert), czterookienne obsługują jedną czynność w kilku widokach.

---

## 6. Macierz moduł × środowisko

`●` = okno w bocznej nawigacji · `○` = wyłącznie komponent własny strefy 2 (bez okna modułowego).

| Moduł | TalkIn | WorkSpace | CodeStudio | MultitaskingAI |
|---|:-:|:-:|:-:|:-:|
| Studio | ● | ● | | |
| Research | ● | ● | | |
| Library | ● | ● | | |
| Translate | ● | | | |
| Browser | ● | ● | | |
| Assistant | ● | | | |
| Roundtable | ● | ● | ● | |
| Workspace | ● | ● | ● | |
| Automations | ○ | ○ | ○ | integracja 24/7 |
| Design | | ● | ● | |
| Apps | | ● | ● | |
| Terminal | | | ● | |
| Developer | | | ● | |
| Diagnostics | | | ● | |
| Agents | ● | ● | ● | rola/ekspert |
| **Pozycji w bocznej nawigacji** | **9** | **9** | **8** | panel orkiestracji (6 sekcji) |

### 6.1. Skład bocznej nawigacji

| Środowisko | Motto (manifest ikon) | Moduły w bocznej nawigacji |
|---|---|---|
| TalkIn | Myśl. Analizuj. Rozumiej. | Studio, Research, Library, Translate, Browser, Assistant, Roundtable, Workspace, Agents |
| WorkSpace | Planuj. Organizuj. Realizuj. | Studio, Research, Library, Browser, Roundtable, Workspace, Design, Apps, Agents |
| CodeStudio | Projektuj. Buduj. Rozwijaj. | Roundtable, Workspace, Design, Apps, Terminal, Developer, Diagnostics, Agents |
| MultitaskingAI | Deleguj. Koordynuj. Nadzoruj. | brak modułów — panel orkiestracji (6 sekcji), role zamiast modułów |

**Uwaga dokumentacyjna (L5):** dokładny skład bocznej nawigacji WorkSpace i CodeStudio nie jest jawnie wyliczony w żadnym pliku zestawu — powyższe składy odczytano z macierzy KANON rozdz. 6 / `README.md` rozdz. 5 i one same są w README oznaczone jako rekonstrukcja do potwierdzenia.

### 6.2. Zasięg modułów

| Zasięg | Moduły | Liczba |
|---|---|---:|
| Trzy środowiska modułowe | Roundtable, Workspace, Agents | 3 |
| Dwa środowiska | Studio, Research, Library, Browser, Design, Apps | 6 |
| Jedno środowisko | Translate, Assistant (TalkIn), Terminal, Developer, Diagnostics (CodeStudio) | 5 |
| Zero okien modułowych | Automations | 1 |

---

## 7. Trzy przypadki szczególne

### 7.1. Automations — jedyny moduł bez okna w bocznej nawigacji

**Fakt:** Automations, jako jedyny z piętnastu modułów, nie ma pozycji w bocznej nawigacji żadnego środowiska. Cała praca projektowa nad automatyką odbywa się z poziomu strony głównej, ze strefy 2 (kafel „Automations” → „Zbuduj automatykę”).

**Cytat źródłowy:** „Brak okna w bocznej nawigacji jest architektoniczną decyzją, nie brakiem — odzwierciedla naturę automatyki jako wytworu konfigurowanego raz, a wykorzystywanego wielokrotnie, w wielu miejscach platformy jednocześnie.” (`automations.md` rozdz. 1.4)

**Konsekwencja projektowa:** w każdej powłoce środowiska pozycja „Automations” **nie może się pojawić** w bocznej nawigacji — ani jako pozycja wyszarzona, ani jako pozycja „niedostępna”. Zgodnie z zasadą zero blokad (KANON rozdz. 8) obowiązuje reguła „brak metody = mniej segmentów, nie zablokowany segment”: moduł jest po prostu nieobecny w tej liście, a obecny w strefie 2.

**Cykl pracy:** budowa (strona główna) → zapis jako komponent własny → wpięcie do sesji modułu docelowego → wykonanie samodzielne wg harmonogramu lub zdarzenia.

### 7.2. Agents — 15. moduł, którego dokument leży w katalogu `okna/`

**Fakt:** Agents jest 15. modułem platformy (Koncepcja platformy rozdz. 9.15), ale jego dokument projektowy leży w `projekt-ui/okna/agents.md`, a nie w `projekt-ui/moduly/`. Katalog `moduly/` zawiera 14 plików.

**Zapis niespójności:** pozycja **L1** w `README.md` rozdz. 7 — „Niespójność klasyfikacji”, priorytet **Ś**. Proponowane rozstrzygnięcie: przenieść dokument do `moduly/agents.md` albo pozostawić jawną notę w README i wyrównać odsyłacze w `elementy-okien.md` rozdz. 0.2.

**Konsekwencje widoczne w danych:**

| Skutek | Opis |
|---|---|
| Tabela `README.md` rozdz. 3.5 obejmuje 14 modułów | Agents pojawia się dopiero w rozdz. 3.3 wśród okien platformowych |
| Kolumna „Komp.” nie odnotowuje Agents | Choć Agent (ekspert) jest jednym z czterech komponentów własnych — `agents.md` rozdz. 1.2 |
| Liczba „15 modułów” wymaga każdorazowego dopisku | „…wraz z 15. modułem Agents, którego dokument leży w `okna/`” — `moduly/README.md` |

**Rozstrzygnięcie na planszy:** Agents prezentowany jest jako pełnoprawny 15. moduł, na równi z pozostałymi, z jawną adnotacją o miejscu dokumentu.

### 7.3. Apps — moduł integrujący Developer i Terminal w jeden proces

**Fakt:** Apps nie duplikuje modułów Developer i Terminal, lecz integruje ich funkcjonalności w jednym, ustrukturyzowanym procesie budowy produktu: architektura → praca równoległa nad frontendem i backendem → wdrożenie.

**Cytat źródłowy:** „Moduł spina w jeden, ciągły proces wszystkie etapy powstawania produktu cyfrowego, integrując funkcjonalności modułów Developer i Terminal w ramach jednej, ustrukturyzowanej przestrzeni…” (`apps.md` rozdz. 1.2)

**Gdzie integracja jest widoczna w interfejsie:** Frontend Workspace i Backend Workspace zawierają edytor, Git Panel i konsolę — czyli komponenty właściwe Developer i Terminal (`apps.md` rozdz. 6.3, kolumna „Miejsce ustanowienia”).

**Zastrzeżenie z dokumentacji:** „praca nad produktem w module Apps działa w pełni samodzielnie bez połączenia z Developer, Terminal, Design czy MultitaskingAI — każde z tych powiązań jest świadomą decyzją”. Integracja jest więc opisem zakresu funkcji, nie twardą zależnością techniczną.

### 7.4. Rozbieżność liczby okien modułu Apps — odnotowanie

Przy zestawianiu inwentarza ujawniła się rozbieżność wewnątrz zestawu źródeł:

| Źródło | Okna własne Apps | Liczba łącznie |
|---|---|---:|
| KANON rozdz. 7.5 | Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel | 6 |
| `moduly/apps.md` rozdz. 2 | Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel | 6 |
| `projekt-ui/README.md` rozdz. 4 | Product Builder, Frontend Workspace, Backend Workspace, Deployment Panel — **brak Architecture Designer** | 5 |

**Rozstrzygnięcie planszy:** przyjęto zapis KANON rozdz. 7.5, zgodny z rozdz. 2 pliku modułu (sześć okien). Tabela zbiorcza README rozdz. 4 pomija Architecture Designer — rozbieżność odnotowano tutaj, nie usuwając jej samodzielnie ze źródła.

---

## 8. Mapa powiązań między modułami

Wykaz obejmuje **wyłącznie powiązania opisane w dokumentacji** (rozdz. 5.4 albo 6.3 pliku modułu). Kierunek `→` oznacza „zasila / przekazuje do”, `↔` — powiązanie obustronne, `═` — pokrewieństwo funkcjonalne z warstwą platformy (nie z modułem).

| # | Powiązanie | Charakter | Okno źródłowe → okno docelowe | Źródło |
|---|---|---|---|---|
| 1 | Studio ↔ Translate | współdzielenie operacji kontekstowych AI przy pracy dwujęzycznej | Tools Panel ↔ Source Panel / Translation Panels | `studio.md` 5.4, `translate.md` 5.4 |
| 2 | Studio → Library | repozytorium źródłowe i cel dla artefaktów | Preview Window / Studio Editor → Library Explorer | `studio.md` 5.4, `library.md` 5.4 |
| 3 | Studio ↔ Research | redakcja wyników badania i dalsze badanie | Report Builder ↔ Studio Editor | `studio.md` 5.4, `research.md` 5.4 |
| 4 | Design → Studio | zasoby wizualne przy redagowaniu dokumentów | Assets Panel → Studio Editor | `studio.md` 5.4, `design.md` 6.3 |
| 5 | Browser ↔ Research | wspólne źródła i zasilanie badania | Sources Panel → Sources Manager | `browser.md` 5.4, `research.md` 5.4 |
| 6 | Browser → Library | trwałe przechowanie odnotowanych materiałów | Notes Panel → Library Explorer | `browser.md` 5.4 |
| 7 | Research ↔ Library | trwałe przechowanie raportów; materiał źródłowy badania | Report Builder → Library Explorer; Library Explorer → Sources Manager | `research.md` 5.4, `library.md` 5.4 |
| 8 | Research ↔ Roundtable | wielomodelowa prezentacja wyników; stanowisko jako ustalenie | Findings Panel → Model Panels; Consensus Panel → Findings Panel | `research.md` 5.4, `roundtable.md` 5.4 |
| 9 | Roundtable → Studio | dalsza redakcja uzgodnionego stanowiska | Consensus Panel → Studio Editor | `roundtable.md` 5.4 |
| 10 | Workspace ↔ Library | Project Library jako odpowiednik ograniczony do projektu | Library Explorer ↔ Project Library | `library.md` 5.4, `workspace.md` 6.3 |
| 11 | Workspace → Automations | powiązanie projektu z procesem automatycznym | Project Dashboard → Automations | `workspace.md` 6.3 |
| 12 | Agents → Workspace | agenci udostępniani jako wykonawcy zadań projektu | Agent Builder → Agent Manager | `workspace.md` 6.3 |
| 13 | Design → Apps | zasoby wizualne dla warstwy frontendowej | Assets Panel / Preview Window → Frontend Workspace | `design.md` 6.3, `apps.md` 6.3 |
| 14 | Developer → Apps | Developer jako komponent budowy produktu | Code Editor / Git Panel → Frontend/Backend Workspace | `developer.md` 5.4, `apps.md` 6.3 |
| 15 | Terminal → Apps | konsola jako komponent budowy produktu | Terminal Tabs → Frontend/Backend Workspace | `apps.md` 6.3 |
| 16 | Apps → Automations | wdrożenie wyzwalane automatyką zdarzeniową | Deployment Panel → Automations | `apps.md` 6.3 |
| 17 | Developer → Terminal | warstwa wykonawcza poleceń budowania i uruchamiania | Build Output / Git Panel → Terminal Tabs | `developer.md` 5.4, `terminal.md` 5.4 |
| 18 | Developer ↔ Diagnostics | analiza błędów; wdrażanie poprawek | Build Output → Errors Panel; Recommendations Panel → Code Editor | `developer.md` 5.4, `diagnostics.md` 5.4 |
| 19 | Diagnostics ↔ Terminal | odtwarzanie i weryfikacja objawów błędu | Diagnostics Center → Terminal Tabs; Output Console → Errors Panel | `diagnostics.md` 5.4, `terminal.md` 5.4 |
| 20 | Agents → wszystkie moduły | ekspert jako wykonawca w każdym module i środowisku | Agent Builder → wybór wykonawcy w sesji modułu | `agents.md` 1.2–1.3 |
| 21 | Automations → wszystkie moduły | krok „akcja w module” w procesie automatycznym | Workflow Builder → dowolny moduł | `automations.md` 6.3 |

### 8.1. Powiązania ponadmodułowe (poza mapą modułów)

| Powiązanie | Charakter | Źródło |
|---|---|---|
| Assistant ═ Always On Display | pokrewieństwo funkcjonalne — wspólny charakter komunikacji głosowej; bez łączenia kontekstu, historii i pamięci | `assistant.md` 5.4 |
| Roundtable ═ MultitaskingAI | ta sama idea wielomodelowości: debata i konsensus vs. podział ról i orkiestracja | `roundtable.md` 5.4 |
| Apps → MultitaskingAI | Executor 1 / Executor 2 przy pracy równoległej nad frontendem i backendem | `apps.md` 6.3 |
| Automations → MultitaskingAI | harmonogramy i kolejki pętli pracy ciągłej 24/7/365 | `automations.md` 6.3 |
| Agents → MultitaskingAI | ekspert przypisany do trwałej roli w zespole modeli | `agents.md` 1.3 |

### 8.2. Odczyt mapy

| Obserwacja | Podstawa |
|---|---|
| **Library jest węzłem o największej liczbie wejść** — zasilana przez Studio, Research, Browser i Workspace | `library.md` 5.4 |
| **Translate ma najwęższy zbiór powiązań** — jedno, ze Studio; działa w pełni samodzielnie | `translate.md` 5.4 (zapis jawny) |
| **Assistant nie ma stałego powiązania międzymodułowego** — ma zdolność operacyjną obejmowania dowolnego modułu jako kroku zlecenia | `assistant.md` 5.4 |
| **Trójkąt CodeStudio** — Developer, Terminal, Diagnostics wiążą się wzajemnie w każdą stronę, a wszystkie trzy zasilają Apps | `developer.md`, `terminal.md`, `diagnostics.md` 5.4 |
| **Żadne powiązanie nie jest domyślnie aktywne** — każde ustanawia Operator z okna konfiguracji albo z okna operacyjnego | zapisy zamykające rozdz. 5.4 / 6.3 wszystkich plików modułowych |

---

## 9. Decyzje projektowe planszy

| # | Decyzja | Uzasadnienie |
|---|---|---|
| 1 | **Ikona modułu odczytana z `manifest.json`, nie dobrana intuicyjnie** | Manifest zawiera jawne przypisania w polu `zastosowanie` („Moduł Library — repozytorium wiedzy”, „Moduł Translate”, „Moduł Assistant — interfejs głosowy”…). Wszystkie 15 ikon ma tam potwierdzenie; żadna nie została dorysowana. |
| 2 | **Karta modułu ma stałą siatkę siedmiu pól** | Tożsamość ma się porównywać, nie opowiadać. Stała siatka pozwala czytać kolumnowo (wszystkie okna wiodące, wszystkie środowiska) — zgodnie z `GESTOSC_WIZUALNA` 8/10. |
| 3 | **Okno wiodące wyróżnione plakietką, nie kolorem** | KANON rozdz. 9: stan nigdy samym kolorem. Plakietka `.dn-plakietka--sygnal` niesie etykietę „wiodące”. |
| 4 | **Filtr środowiska jako grupa przełączników, nie lista rozwijana** | Cztery środowiska plus „wszystkie” mieszczą się w jednym rzędzie; przełączniki pokazują cały zakres wyboru bez otwierania. Wybór nie blokuje niczego — zmienia tylko widoczność kart. |
| 5 | **Zero `disabled` na całej planszy** | KANON rozdz. 8 (ADL-017). Filtr, który nie ma trafień, pokazuje stan pusty z komunikatem, nie wyłącza przycisku. |
| 6 | **Automations pokazany bez ikony „brak”** | Moduł nie jest ułomny — jego nieobecność w bocznej nawigacji jest decyzją architektoniczną. Karta niesie plakietkę „strefa 2” zamiast pustego pola środowisk. |
| 7 | **Agents pokazany jako pełnoprawny 15. moduł, z adnotacją L1** | Klasyfikacja wg Koncepcji platformy rozdz. 9.15, a nie wg położenia pliku. Adnotacja o katalogu `okna/` towarzyszy karcie, żeby nie ukrywać niespójności. |
| 8 | **Grupowanie oznaczone jawnie jako warstwa redakcyjna** | Dokumentacja nie definiuje kategorii modułów. Każde przypisanie ma kotwicę w opisie produktu i w udokumentowanym powiązaniu — bez tego byłoby kategorią wymyśloną, czego zakazuje zasada nadrzędna zespołu. |
| 9 | **Mapa powiązań rysuje wyłącznie krawędzie z tabel 5.4 / 6.3** | Powiązania „Agents → wszystkie” i „Automations → wszystkie” rysowane są odrębną klasą i pojawiają się dopiero po wskazaniu węzła — inaczej 28 dodatkowych krawędzi zasłoniłoby czytelny rdzeń mapy. |
| 10 | **Diagram jest SVG w `currentColor`, bez rastrów** | Ten sam wymóg co dla ikon (KANON rozdz. 4): jedna geometria, dwa motywy, skalowanie bez utraty. |
| 11 | **Liczby na planszy policzone, nie oszacowane** | 62 okna własne, 77 wystąpień, 9/9/8 pozycji w nawigacji — każda liczba wyprowadzona z tabeli KANON rozdz. 7.5 i macierzy rozdz. 6, a istnienie 15 plików prototypów potwierdzone odczytem katalogu `05-okna/moduly/`. |
| 12 | **Rozbieżność liczby okien Apps odnotowana, nie ukryta** | Zasada „nie usuwaj niespójności źródła po cichu” — plansza podaje rozstrzygnięcie i wskazuje miejsce rozbieżności (rozdz. 7.4). |

---

## 10. Źródła

| Zakres | Plik | Rozdziały |
|---|---|---|
| Kanon projektu | `WYNIK/KANON.md` | 4 (ikony), 6 (architektura, macierz), 7.4–7.5 (inwentarz okien), 8 (zero blokad), 9 (dostępność), 10 (redakcja) |
| Indeks projektu UI | `dok/projekt-ui/README.md` | 2 (model pojęciowy), 3.3, 3.5 (istota modułów), 4 (inwentarz okien), 5 (macierz), 6 (mapa zależności), 7 (luki L1, L5, L9) |
| Moduły | `dok/projekt-ui/moduly/*.md` | nagłówek (środowiska, forma udostępnienia), 1 (przeznaczenie), 2 (komplet okien), 5.4 albo 6.3 (powiązania) |
| Moduł Agents | `dok/projekt-ui/okna/agents.md` | nota terminologiczna, 1.1 (łańcuch wartości), 1.2 (komponenty własne), 1.3 (gdzie ekspert działa), 2 (mapa okien) |
| Katalog modułów | `dok/projekt-ui/moduly/README.md` | zasada „jeden plik = jeden moduł”, struktura 7 rozdziałów |
| Elementy powłoki | `dok/projekt-ui/przeplyw/elementy-okien.md` | 3.4 (boczna nawigacja), 3.5 (obszar roboczy), 4.4 (panel orkiestracji), Zał. A (katalog ikon) |
| Ikony | `WYNIK/zasoby/ikony/manifest.json`, `WYNIK/zasoby/ikony/svg/*.svg` | pola `nazwa`, `zrodlo`, `zastosowanie` |
| Emblematy i motta środowisk | `WYNIK/zasoby/marka/srodowiska/*.svg`, `manifest.json` | pozycje `srodowisko-*` |
| Prototypy okien | `WYNIK/05-okna/moduly/*.html` | 15 plików — istnienie potwierdzone odczytem katalogu |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
