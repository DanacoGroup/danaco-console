# Danaco Console — Moduł Terminal

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
| **Tytuł** | Moduł Terminal |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant · Operator |
| **Przeznaczenie** | Ustala zakres funkcjonalny, komplet okien operacyjnych i zachowanie modułu Terminal — warstwy wykonawczej środowiska CodeStudio — jako źródło wykonawcze dla dewelopera i projektanta. |
| **Zakres** | komplet okien operacyjnych modułu Terminal (Chat Window, Execution Loop Window, Terminal Tabs, Output Console, Process Monitor, Session Manager, Task & Schedule, Script Library), katalog funkcji, komendy obszaru `terminal`, żetony i komponenty widoku, stany kontrolek, przepływy pracy i scenariusze użycia |
| **Poza zakresem** | edycja kodu i drzewo repozytorium — [Moduł Developer](developer.md); analiza i grupowanie błędów — [Moduł Diagnostics](diagnostics.md); budowanie wizualnych przebiegów automatyzacji — [Moduł Automations](automations.md); przeglądanie stron i DOM — [Moduł Browser](browser.md) (rozdz. 1.5) |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Moduł Developer](developer.md) · [Moduł Diagnostics](diagnostics.md) · [Moduł Automations](automations.md) · [Moduł Browser](browser.md) · [Model danych](../architektura/model-danych.md) · [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/terminal.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `terminal`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` · `design/05-okna/moduly/terminal.html` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad; klucze i parametry jawne; izolacja techniczna procesu jest jawną decyzją Operatora, nie wymogiem; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
   - [1.1 Definicja](#11-definicja)
   - [1.2 Dla kogo](#12-dla-kogo)
   - [1.3 Po co — wartość modułu](#13-po-co--wartość-modułu)
   - [1.4 Granica tematyczna — zakres modułu](#14-granica-tematyczna--zakres-modułu)
   - [1.5 Granica negatywna — obszary poza modułem](#15-granica-negatywna--obszary-poza-modułem)
   - [1.6 Miejsce w architekturze platformy](#16-miejsce-w-architekturze-platformy)
   - [1.7 Dostępność i forma udostępnienia](#17-dostępność-i-forma-udostępnienia)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - [2.1 Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
   - [2.2 Makieta zbiorcza modułu w stanie spoczynku](#22-makieta-zbiorcza-modułu-w-stanie-spoczynku)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
   - [3.1 Chat Window (główne okno komunikacji Użytkownik ↔ Wykonawca)](#31-chat-window-główne-okno-komunikacji-użytkownik--wykonawca)
   - [3.2 Execution Loop Window (okno pętli wykonawczej Koordynator ↔ Wykonawca)](#32-execution-loop-window-okno-pętli-wykonawczej-koordynator--wykonawca)
   - [3.3 Terminal Tabs](#33-terminal-tabs)
   - [3.4 Output Console](#34-output-console)
   - [3.5 Process Monitor](#35-process-monitor)
   - [3.6 Session Manager](#36-session-manager)
   - [3.7 Task & Schedule](#37-task--schedule)
   - [3.8 Script Library](#38-script-library)
4. [Przepływy pracy](#4-przepływy-pracy)
   - [4.1 Przepływ podstawowy — polecenie wydane ręcznie przez Operatora](#41-przepływ-podstawowy--polecenie-wydane-ręcznie-przez-operatora)
   - [4.2 Przepływ rozszerzony — polecenie zlecone z Chat Window](#42-przepływ-rozszerzony--polecenie-zlecone-z-chat-window)
   - [4.3 Przepływ pętli wykonawczej — zlecenie wielozadaniowe](#43-przepływ-pętli-wykonawczej--zlecenie-wielozadaniowe)
   - [4.4 Przepływ rozszerzony — awaria procesu i przekazanie do Diagnostics](#44-przepływ-rozszerzony--awaria-procesu-i-przekazanie-do-diagnostics)
   - [4.5 Przepływ zbiorczy — równoległa praca wielu kart i powłok](#45-przepływ-zbiorczy--równoległa-praca-wielu-kart-i-powłok)
   - [4.6 Przepływ zadania cyklicznego i wyzwalacza](#46-przepływ-zadania-cyklicznego-i-wyzwalacza)
5. [Stany, dane i powiązania](#5-stany-dane-i-powiązania)
   - [5.1 Model stanów procesu](#51-model-stanów-procesu)
   - [5.2 Model danych wykorzystywany przez moduł](#52-model-danych-wykorzystywany-przez-moduł)
   - [5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Terminal](#53-izolacja-i-konfigurowalność--punkty-właściwe-modułowi-terminal)
   - [5.4 Powiązania z innymi modułami](#54-powiązania-z-innymi-modułami)
6. [Scenariusze użycia](#6-scenariusze-użycia)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
   - [7.1 Powłoki, sesje i emulator terminala](#71-powłoki-sesje-i-emulator-terminala)
   - [7.2 Zdalne wykonanie i połączenia](#72-zdalne-wykonanie-i-połączenia)
   - [7.3 Skrypty, snippety i biblioteka](#73-skrypty-snippety-i-biblioteka)
   - [7.4 Zadania, harmonogram i automatyzacja poleceniowa](#74-zadania-harmonogram-i-automatyzacja-poleceniowa)
   - [7.5 Sekrety, środowisko i uwierzytelnienie CLI](#75-sekrety-środowisko-i-uwierzytelnienie-cli)
   - [7.6 Wynik, obserwowalność i rejestracja](#76-wynik-obserwowalność-i-rejestracja)
   - [7.7 Narzędzia wiersza poleceń wbudowane w moduł](#77-narzędzia-wiersza-poleceń-wbudowane-w-moduł)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)
10. [Punkty łamania i kryteria odbioru](#10-punkty-łamania-i-kryteria-odbioru)
   - [10.1 Punkty łamania](#101-punkty-łamania)
   - [10.2 Kryteria odbioru](#102-kryteria-odbioru)
11. [Załącznik — pełny wykaz komend kontraktu modułu Terminal](#załącznik--pełny-wykaz-komend-kontraktu-modułu-terminal)
   - [Obszar `terminal` — 27 komend](#obszar-terminal--27-komend)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja

Terminal jest modułem pracy z konsolami i środowiskami wykonawczymi — bezpośrednim dostępem do powłok systemowych i narzędzi wiersza poleceń z poziomu platformy, bez przełączania się do zewnętrznej aplikacji terminala. Moduł integruje w jednej przestrzeni czynności dotąd rozproszone między osobne okna systemowe: prowadzenie równoległych sesji powłok, gromadzenie i przeszukiwanie ich wyniku oraz obserwację i kontrolę uruchomionych procesów — w tym procesów zainicjowanych nie przez Operatora, lecz poleceniem wydanym przez Wykonawcę. Terminal jest warstwą wykonawczą środowiska CodeStudio: to w nim faktycznie uruchamia się to, co Developer i Diagnostics jedynie planują lub analizują.

Terminal jest kompletnym środowiskiem powłoki, sesji, skryptów, zadań i zdalnego wykonania. Jedna karta modułu zastępuje cały stos osobnych aplikacji trzymanych obok siebie przez inżyniera: emulator terminala, multiplekser sesji, klienta SSH i SFTP, menedżera połączeń zdalnych, harmonogram zadań, runner skryptów, menedżera sekretów wiersza poleceń, konsolę szeregową, klienta kontenerów oraz narzędzia rejestracji i odtwarzania sesji. Każda zdolność modułu jest dostępna bez warunku wstępnego (zero blokad), a jej widoczność, zakres i parametry Operator personalizuje z okna konfiguracji. Żadna funkcja nie wymusza izolacji ani uwierzytelnienia jako warunku startu — izolacja i sekrety są jawnymi, odwracalnymi decyzjami.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Terminal |
|---|---|
| Deweloperzy warstwy serwerowej i warstwy klienckiej | Uruchamianie poleceń budowania, instalacji zależności, serwerów deweloperskich |
| Inżynierowie DevOps i administratorzy systemów | Praca z wieloma powłokami jednocześnie, w tym sesjami na hostach zdalnych |
| Inżynierowie automatyzacji i skryptowania | Uruchamianie i obserwacja skryptów Node.js oraz Python, iteracyjne poprawki na żywo |
| Inżynierowie jakości i testerzy | Uruchamianie zestawów testów z poziomu wiersza poleceń, obserwacja wyniku w czasie rzeczywistym |
| Operatorzy nadzorujący pracę Wykonawcy | Podgląd i kontrola poleceń wykonywanych autonomicznie w toku pracy nienadzorowanej |

### 1.3. Po co — wartość modułu

| Problem klasycznego rozproszenia narzędzi | Rozwiązanie w Terminal |
|---|---|
| Terminal systemowy i okno komunikacji z Wykonawcą to dwie osobne aplikacje | Terminal Tabs i Chat Window w jednej karcie sesji, nad tym samym środowiskiem wykonawczym |
| Brak wglądu w to, co dokładnie uruchomił Wykonawca w tle | Process Monitor rozróżnia procesy zainicjowane przez Operatora i przez Wykonawcę |
| Brak wglądu w dekompozycję zlecenia na zadania powłoki | Execution Loop Window pokazuje kolejkę zadań, komunikaty sterujące i wynik kontroli jakości |
| Rozproszony wynik wielu równoległych powłok | Output Console gromadzi wynik wszystkich kart w jednym, przeszukiwalnym strumieniu |
| Konieczność ręcznego zabijania zawieszonych procesów przez menedżera systemowego | Process Monitor pozwala zakończyć proces bez opuszczania platformy |
| Utrata kontekstu powłoki po zamknięciu okna systemowego | Trwałość stanu procesu sesji po stronie serwera — rozłączenie klienta nie przerywa działania powłoki |
| Rozdzielenie hostów zdalnych, harmonogramu i biblioteki skryptów między osobne programy | Session Manager, Task & Schedule i Script Library w tej samej przestrzeni modułu |

### 1.4. Granica tematyczna — zakres modułu

| Obszar | Zakres w module Terminal |
|---|---|
| Powłoka interaktywna | PowerShell, CMD, Bash/Zsh/Fish (WSL i natywne), Git Bash; REPL: Node.js, Python, Deno, .NET (`dotnet script`), SQL |
| Sesje i multipleksacja | Karty, panele, układy zapisywane; trwałość sesji po stronie serwera; wznowienie po rozłączeniu; nazwane sesje trwałe |
| Zdalne wykonanie | SSH, SFTP, przekierowanie portów, agent SSH i ProxyJump, sesje na wielu hostach jednocześnie z rozgłaszaniem polecenia, `docker exec`, `kubectl exec` |
| Transfer plików | SFTP/SCP, `rsync`, dwupanelowy menedżer plików zdalny↔lokalny, edycja pliku zdalnego w miejscu |
| Skrypty | Edytor skryptów wieloliniowy z uruchomieniem, biblioteka skryptów, parametryzacja, snippety, generowanie z języka naturalnego |
| Zadania i automatyzacja | Runner zadań (`Makefile`/`npm scripts`/`Taskfile`), harmonogram (cron/at), wyzwalacze (watch plików, zdarzenia), potoki poleceń |
| Sekrety i środowisko | Menedżer zmiennych środowiskowych i plików `.env`, magazyn sekretów sesji, wstrzykiwanie tajnych wartości bez ujawnienia w historii |
| Obserwowalność | Rejestr procesów, wykresy zasobów, rejestracja i odtwarzanie sesji (asciinema), eksport transkryptów i dzienników |
| Powłoka portowa i sprzętowa | Konsola szeregowa (RS‑232/USB), połączenie Telnet (urządzenia sieciowe) |

### 1.5. Granica negatywna — obszary poza modułem

| Poza zakresem | Właściwe miejsce na platformie | Uzasadnienie granicy |
|---|---|---|
| Edycja kodu projektu, drzewo repozytorium, historia commitów w widoku graficznym | Moduł **Developer** (Code Editor, Project Tree, Git Panel) | Terminal jest warstwą wykonania, nie warsztatem edycji; wskaźnik gałęzi Git w Terminalu linkuje do Developera |
| Analiza i grupowanie błędów, odtwarzanie objawów w ustrukturyzowanej sesji diagnostycznej | Moduł **Diagnostics** (Errors Panel, Diagnostics Center) | Terminal wykrywa i przekazuje błąd; klasyfikacja i śledztwo należą do Diagnostics |
| Wizualne budowanie przebiegów automatyzacji między aplikacjami (nody, konektory) | Moduł **Automations** | Terminal automatyzuje polecenia powłoki, nie integracje między usługami zewnętrznymi |
| Przeglądanie stron i praca z DOM | Moduł **Browser** | Terminal wykonuje `curl`/`httpie`, lecz nie renderuje stron |
| Trwałe zarządzanie infrastrukturą chmury jako produkt (provisioning, pulpity IaC) | Poza platformą | Terminal uruchamia `terraform`/`kubectl` jako polecenia, nie zastępuje pulpitów IaC |

Granica jest linkująca, nie blokująca: z każdego z powyższych obszarów Terminal ma jawne, konfigurowalne przejście (rozdz. 5.4 i 8), a nie ścianę.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: CodeStudio
        ▼
ŚRODOWISKO: CODESTUDIO ── boczna nawigacja modułów
        │        (Workspace · Roundtable · Design · Terminal ·
        │         Developer · Diagnostics · Apps · Agents)
        ▼
MODUŁ: TERMINAL ────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Terminal Tabs ·
  Output Console · Process Monitor · Session Manager ·
  Task & Schedule · Script Library
        │
        ▼
KARTA SESJI — własny katalog roboczy, środowisko procesu i kontekst
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 5.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | CodeStudio (wyłącznie) |
| Środowiska TalkIn i WorkSpace | Moduł niedostępny — praca z powłoką nie należy do trybu wiedzy ani produktywności projektowej |
| Komponent własny | Terminal nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 8 (łącznie z Chat Window i Execution Loop Window) |
| Typ pracy | Sesyjna, wielopowłokowa — równoległe karty powłok w obrębie jednej karty sesji platformy |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Rola w module | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca (wspólne wszystkim modułom) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Centralny punkt pracy i podstawowy mechanizm sterowania procesami: polecenia w języku naturalnym, strumień odpowiedzi, zatwierdzanie i przerywanie działań, wyjaśnianie wyniku | 1 | Widoczne bez interakcji — okno stałe lewej kolumny |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość | Pętla wykonawcza modułu: dekompozycja zlecenia na polecenia powłoki, kolejka i stan zadań, kontrola jakości uruchomień, sterowanie przebiegiem | 2 | Znacznik „Execution Loop ▸” w stopce Chat Window; kolumna zwija się po zamknięciu przebiegu |
| 3 | Terminal Tabs | Okno edycyjne (interakcja wiersza poleceń) | Prawa kolumna, dominująca — pasek kart nad obszarem powłoki | Równoległe sesje powłok — PowerShell, CMD, Bash, Node.js, Python, SSH, kontener, pod, port szeregowy | 1 | Widoczne bez interakcji — wiodące okno obszaru roboczego |
| 4 | Output Console | Podgląd i porównanie / monitor | Prawa kolumna, dominująca, obok Terminal Tabs | Zagregowany wynik poleceń ze wszystkich kart, parser strukturalny, różnicowanie, odtwarzanie | 1 | Widoczne bez interakcji — druga część obszaru roboczego |
| 5 | Process Monitor | Monitor procesu | Wąskie okno pomocnicze po prawej stronie okna centralnego, otwierane na żądanie, na żywo | Obserwacja i kontrola uruchomionych procesów, drzewo procesów, metryki zasobów | 2 | Ikona szybkiego dostępu w prawym górnym rogu okna centralnego albo polecenie „pokaż procesy”; okno pomocnicze znika po zamknięciu |
| 6 | Session Manager | Nawigacja / katalog | Kolumna boczna, otwierana jako rozszerzenie boczne | Książka hostów i połączeń, sesje trwałe, profile powłok, klucze i tunele | 2 | Znacznik `[Host ▾]` w pasku kontekstu, menu ☰ modułu; kolumna zwija się po wybraniu hosta |
| 7 | Task & Schedule | Monitor / edytor zadań | Prawa kolumna, dominująca — zakładka obszaru roboczego | Runner zadań, harmonogram, wyzwalacze, potoki, historia przebiegów | 3 | Zakładka obszaru roboczego wywoływana z elementu zbiorczego `Operacje ▾` |
| 8 | Script Library | Katalog / edytor | Kolumna boczna, otwierana jako rozszerzenie boczne; dostępna z palety poleceń | Repozytorium skryptów i snippetów, edytor, parametryzacja, linter | 3 | Menu ☰ modułu, paleta poleceń `Ctrl/Cmd + K`, polecenie języka naturalnego „otwórz bibliotekę skryptów” |


### 2.1. Warstwy widoczności w module

Interfejs modułu Terminal ujawnia swoje możliwości stopniowo. Obowiązuje zasada nadrzędna: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Moduł zastępuje emulator terminala, multiplekser sesji, klienta SSH i SFTP, harmonogram zadań, runner skryptów i menedżera sekretów, lecz liczba tych zdolności nie wpływa na złożoność wizualną przestrzeni pracy — złożoność pozostaje w architekturze modułu i ujawnia się dopiero w chwili wystąpienia potrzeby. Każdy element interfejsu modułu należy do dokładnie jednej z czterech warstw widoczności.

**Warstwa 1 — zawsze widoczna.** Chat Window w lewej kolumnie, obszar powłoki aktywnej karty Terminal Tabs wraz z paskiem kart, strumień Output Console, pasek kontekstu karty (host, katalog roboczy, gałąź repozytorium), wskaźniki stanu karty i stanu procesu oraz — po otwarciu — kolejka zadań i strumień komunikatów sterujących Execution Loop Window. Elementy warstwy 1 zajmują ponad 80% powierzchni interfejsu modułu i są dostępne bez jakiejkolwiek interakcji.

**Warstwa 2 — widoczna na żądanie.** Znaczniki kontekstowe paska sesji: `[Danaco Console] [Ubuntu] [Bash ▾] [build-01 ▾] [Model ▾]`, selektor powłoki przycisku `[+ ▾]`, selektor widoku i formatu Output Console, filtr i grupowanie Process Monitor, pole Grep, formularz parametrów skryptu, przyciski „Uruchom teraz”, „Wstaw do terminala”, „Przerwij” i „Zakończ” oraz wyzwalacze otwarcia okna pomocniczego Process Monitor i kolumn Session Manager i Execution Loop Window. Wąskie okno po prawej stronie okna centralnego jest miejscem wszystkich narzędzi pomocniczych modułu i otwiera się wyłącznie na żądanie, z komponentu (ikony) w prawym górnym rogu okna centralnego. Element warstwy 2 przywoływany jest kliknięciem znacznika, ikony lub przełącznika i **zwija się samoczynnie po użyciu** — po wybraniu powłoki, hosta, modelu czy formatu selektor znika z przestrzeni roboczej.

**Warstwa 3 — rozwinięcia kontekstowe.** Menu ⋮ wiersza procesu, wiersza zadania i bloku polecenia, menu ☰ modułu (Script Library, import konfiguracji SSH, ustawienia widoku), panel popover ustawień sesji uruchamiany ikoną ⚙, panel szczegółów procesu z drzewem potomków, pasek transportu odtwarzania nagrania, edytor potoku, przełączniki prezentacji strumienia oraz elementy zbiorcze `Operacje ▾`, `Przebieg ▾` i `Snippety ▾`, które zastępują rozłożone paski dziesiątek przycisków. Panele i menu warstwy 3 znikają całkowicie po zamknięciu.

**Warstwa 4 — funkcje eksperckie.** Rozgłaszanie polecenia do wielu paneli i hostów, operacje na grupach hostów, menedżer kluczy SSH i polityka `known_hosts`, tunele dynamiczne SOCKS, maskowanie wartości wrażliwych według wzorców, konsola szeregowa i sesja Telnet, limity buforowania i progi powiadomień oraz narzędzia diagnostyczne niskiego poziomu rejestru procesów. Elementy warstwy 4 nie występują w interfejsie: przywołuje się je wyłącznie poleceniem języka naturalnego w Chat Window („rozgłoś to polecenie na grupę Produkcja”, „wygeneruj klucz SSH dla web-01”), skrótem klawiszowym, wyszukiwarką funkcji palety poleceń `Ctrl/Cmd + K` albo w trybie administracyjnym. Użytkownik podstawowy modułu ich nie widzi.

**Zasada jednego kliknięcia.** Każda funkcja modułu ukryta w warstwach 2–4 pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji Terminala głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny przestrzeni powłoki, nie utrudnia dostępu do zdolności modułu.

### 2.2. Makieta zbiorcza modułu w stanie spoczynku

```
 Makieta zbiorcza — Moduł Terminal              Dostępność: CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Terminal Tabs                       │ ☰
  nawigacja │ Użytkownik ↔         │ ● Bash ×│ [+ ▾]                     │
  modułów   │ Wykonawca            ├─────────────────────────────────────┤
  (poza     │                      │ [Danaco Console][Ubuntu][build-01 ▾]│
  zakresem  │ polecenia ·          │ [Bash ▾] ~/danaco/api · main        │
  tego      │ strumień odpowiedzi  ├─────────────────────────────────────┤
  dokumentu)│                      │ obszar powłoki aktywnej karty       │
            │                      │ (warstwa 1 — ponad 80% powierzchni) │
            │                      │                                     │
            │                      │ ─────────────────────────────────── │
            │                      │ Output Console — strumień na żywo   │
            │                      │                                     │
            │ ──────────────────── │                       [Operacje ▾]⋮ │
            │ Execution Loop ▸     │                                     │
            │ [Model ▾]  [Procesy ▾] [Host ▾]                            │
 ═══════════════════════════════════════════════════════════════════════════
```

Układ modułu jest wyłącznie pionowy: kolumny sąsiadują poziomo, regulacji podlega wyłącznie ich szerokość. Chat Window zajmuje lewą kolumnę o pełnej wysokości obszaru roboczego, Execution Loop Window otwiera się jako kolumna sąsiadująca, obszar roboczy modułu zajmuje kolumnę prawą i dominującą, a okna pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego.

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (główne okno komunikacji Użytkownik ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |
| Warstwa widoczności | Warstwa 1 — zawsze widoczne; okno stałe lewej kolumny, dostępne bez interakcji |

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą i stanowi centralny punkt pracy w module Terminal oraz podstawowy mechanizm sterowania wszystkimi procesami wykonywanymi w powłokach modułu. Użytkownik wydaje w nim polecenia w języku naturalnym, obserwuje strumień odpowiedzi i wyników, zatwierdza oraz przerywa działania, a także uzyskuje wyjaśnienie wyniku i kontekstu. Okno występuje w tym samym miejscu układu we wszystkich modułach i środowiskach platformy — w lewej kolumnie obszaru roboczego.

**Zawartość i pełny arsenał funkcji.**

- Historia komunikacji w kontekście bieżącej karty Terminal Tabs, jej hosta oraz katalogu roboczego.
- Pole wprowadzania poleceń, z rozpoznawaniem intencji: pytanie do Wykonawcy kontra polecenie do wykonania w powłoce.
- Strumień odpowiedzi Wykonawcy przekazywany na żywo kanałem WebSocket.
- Zatwierdzanie i przerywanie działań uruchomionych w powłokach modułu bezpośrednio z okna komunikacji.
- Generowanie polecenia powłoki na podstawie opisu w języku naturalnym, w postaci „znajdź procesy nasłuchujące na porcie 3000”.
- Wstawienie wygenerowanego polecenia do aktywnej karty Terminal Tabs bez automatycznego uruchomienia — do przeglądu przed wykonaniem.
- Uruchomienie polecenia bezpośrednio przez Wykonawcę w aktywnej karcie, z pełnym zapisem w Output Console i wpisem w Process Monitor oznaczonym inicjatorem „Wykonawca”.
- Wyjaśnienie zaznaczonego fragmentu wyniku z Output Console („wyjaśnij ten błąd”, „co oznacza ten kod wyjścia”).
- Generowanie skryptu wieloliniowego (PowerShell, Bash, Python) gotowego do zapisania w Script Library.
- Przekształcenie polecenia w skrypt z parametrami, zadanie harmonogramu z wyrażeniem cron oraz przeniesienie tunelu portowego do książki hostów — jednym poleceniem w języku naturalnym.
- Cytowanie poleceń i wyniku w treści komunikacji blokiem czcionki `--dn-ff-mono`.
- Załączanie plików — pliku logu albo skryptu — do wiadomości jako kontekstu zapytania.
- Przekazanie fragmentu z Output Console jako kontekstu jednym poleceniem („wyjaśnij zaznaczenie”).
- Skrót do modułu Diagnostics — przekazanie bieżącego problemu jako materiału źródłowego.
- Wskaźnik, której karty Terminal Tabs dotyczą polecenia wydawane z poziomu okna komunikacji (fokus kontekstowy).
- Otwarcie Execution Loop Window dla zlecenia wymagającego dekompozycji na wiele poleceń lub hostów.
- Szybkie menu kontekstowe konfiguracji bieżącej sesji: kanał modelu, zakres pamięci, izolacja techniczna procesu.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window — Użytkownik ↔ Wykonawca  │  Obszar roboczy modułu
  (lewa kolumna, pełna wysokość)        │  (prawa kolumna, dominująca)
 ───────────────────────────────────────┤
  [Danaco Console][Model ▾]  ⚙          │  Terminal Tabs
  kontekst: karta „Bash” · ~/api        │  Output Console
 ───────────────────────────────────────┤
  Historia komunikacji (przewijana)     │
   ┌─ Użytkownik ───────────────────┐   │
   │ „uruchom testy i pokaż wynik”  │   │
   └────────────────────────────────┘   │
   ┌─ Wykonawca (strumień na żywo) ─┐   │
   │ proponowane polecenie:         │   │
   │  npm test -- --watch=false     │   │
   │                             ⋮  │   │
   └────────────────────────────────┘   │
 ───────────────────────────────────────┤
  [ 📎 ]  Pole poleceń                  │
  ………………………………………  [ Wyślij ▶ ]         │
 ───────────────────────────────────────┤
  Execution Loop ▸                      │
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Wskaźnik kontekstu karty | Etykieta tekstowa w nagłówku | Pokazuje, której karty Terminal Tabs dotyczą bieżące polecenia | Mały tekst `--dn-fs-xs` z ikoną powłoki | domyślny · aktualizowany przy zmianie karty aktywnej | Kliknięcie otwiera listę kart do zmiany fokusu | Nagłówek Chat Window | 1 | Widoczny bez interakcji; kliknięcie rozwija listę kart (warstwa 2), która zwija się po wyborze |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie pytania lub polecenia do Wykonawcy | Duże pole w dolnej części lewej kolumny, rozciągliwe w pionie | domyślny · fokus · z załącznikiem · błąd (pusta wysyłka) | Enter wysyła, Shift+Enter nowa linia | Lewa kolumna, poniżej historii | 1 | Widoczne bez interakcji |
| Blok proponowanego polecenia | czcionka `--dn-ff-mono`, tło drugorzędne | Prezentacja wygenerowanego polecenia powłoki | Średni blok w treści odpowiedzi | domyślny · zaznaczony do edycji | Kliknięcie w treść pozwala edytować polecenie przed wstawieniem | Treść odpowiedzi Wykonawcy | 1 | Widoczny bez interakcji jako element strumienia odpowiedzi |
| Przycisk „Wstaw do terminala” | Przycisk `--zarys`, mały | Przeniesienie polecenia do aktywnej karty bez uruchomienia | Mały przycisk tekstowy w pasku akcji | domyślny · wskazanie kursorem | Wkleja polecenie w wierszu poleceń aktywnej karty, kursor gotowy do Enter | Pod blokiem proponowanego polecenia | 2 | Ujawnia się pod blokiem proponowanego polecenia i znika wraz z zamknięciem bloku |
| Przycisk „Uruchom teraz” | Przycisk `--sygnal`, mały | Wykonanie polecenia bezpośrednio przez Wykonawcę | Mały przycisk wyróżniony akcentem sygnałowym — sygnalizuje wagę akcji, nie blokuje jej wykonania | domyślny · ładowanie (polecenie w toku) | Uruchamia polecenie w aktywnej karcie; wynik trafia do Output Console, wpis w Process Monitor oznaczony jako „inicjator: Wykonawca” | Pod blokiem proponowanego polecenia | 2 | Ujawnia się pod blokiem proponowanego polecenia i znika wraz z zamknięciem bloku |
| Przycisk „Przerwij” | Przycisk `--blad`, mały | Przerwanie działania uruchomionego z okna komunikacji | Mały przycisk tekstowy w pasku stanu odpowiedzi | ukryty · widoczny (działanie w toku) | Wysyła sygnał zatrzymania do procesu i odnotowuje przerwanie w Process Monitor | Pasek stanu odpowiedzi Wykonawcy | 2 | Ujawnia się wyłącznie na czas trwania działania |
| Przycisk otwarcia pętli wykonawczej | Przycisk `--zarys`, mały | Otwarcie Execution Loop Window jako kolumny sąsiadującej | Mały przycisk tekstowy w stopce kolumny | domyślny · aktywny (kolumna otwarta) | Otwiera kolumnę pętli wykonawczej z bieżącym zleceniem | Dolna część lewej kolumny | 2 | Znacznik „Execution Loop ▸” w stopce kolumny |
| Przycisk ikony ustawień sesji | Ikona (ustawienia) | Szybka zmiana konfiguracji warstwy sesji | Mała ikona 36×36 px | domyślny · rozwinięte menu | Otwiera uproszczone menu kontekstowe: kanał modelu, pamięć, izolacja techniczna | Prawy róg nagłówka | 3 | Ikona ⚙ w nagłówku otwiera panel popover: kanał modelu, pamięć, izolacja techniczna |
| Przycisk załącznika | Ikona (spinacz) | Dołączenie pliku logu albo skryptu do wiadomości | Mała ikona | domyślny · aktywny (plik dołączony) | Otwiera okno wyboru pliku lub akceptuje przeciągnięcie | Lewa strona pola poleceń | 2 | Ikona 📎 przy polu poleceń |

---

### 3.2. Execution Loop Window (okno pętli wykonawczej Koordynator ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Pętla prowadzona w granicach karty sesji; widoczność zadań innych kart konfigurowalna (rozdz. 5.3) |
| Warstwa widoczności | Warstwa 2 — otwierane znacznikiem „Execution Loop ▸” w stopce Chat Window albo poleceniem języka naturalnego; kolumna znika po zamknięciu przebiegu |

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą i odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. W module Terminal pętla dotyczy zadań powłokowych: Koordynator rozkłada zlecenie Użytkownika na sekwencję poleceń, przydziela je kartom Terminal Tabs, hostom z książki połączeń, kontenerom i podom, po czym nadzoruje ich wykonanie, ocenia kody wyjścia i wzorce błędów rozpoznane w Output Console oraz decyduje o ponowieniu, korekcie parametrów lub zatrzymaniu przebiegu. Okno otwierane jest jako kolumna sąsiadująca z głównym oknem komunikacji.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie i jego dekompozycja na zadania powłokowe: polecenie, karta docelowa, host, katalog roboczy, zestaw zmiennych środowiskowych.
- Kolejka zadań ze stanem każdego zadania: oczekujące, przydzielone, w toku, zakończone sukcesem, zakończone błędem, ponawiane, wstrzymane.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą: przydział zadania, potwierdzenie podjęcia, raport postępu, zgłoszenie błędu, żądanie decyzji.
- Wyniki kontroli jakości uruchomienia: kod wyjścia, rozpoznane wzorce błędów, wynik analizy statycznej skryptu, porównanie z wynikiem poprzedniego przebiegu.
- Decyzje o ponowieniu zadania — ponowienie z tymi samymi parametrami, ponowienie po korekcie polecenia, przekazanie zadania do innej karty lub innego hosta.
- Wskaźniki przebiegu pętli: liczba zadań zakończonych i pozostałych, czas trwania przebiegu, liczba ponowień, wskaźnik współbieżności.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia w toku.
- Graf zależności zadań dla potoków poleceń — kolejność kroków, warunki sukcesu i porażki, gałęzie równoległe.
- Przejście z każdego zadania do karty źródłowej w Terminal Tabs, do bloku wyniku w Output Console i do wiersza procesu w Process Monitor.
- Powiązanie zadania pętli z pozycją harmonogramu lub wyzwalaczem z okna Task & Schedule.
- Eksport przebiegu pętli jako dziennika zawierającego zlecenie, zadania, komunikaty sterujące i wyniki kontroli.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window        │ Execution Loop Window       │ Obszar roboczy modułu
  Użytkownik ↔       │ Koordynator ↔ Wykonawca     │ (Terminal Tabs,
  Wykonawca          │ (kolumna sąsiadująca)       │  Output Console)
                     ├─────────────────────────────┤
                     │ Zlecenie: „zbuduj i wdróż   │
                     │ na dwóch hostach”        ⋮  │
                     │ przebieg 00:02:41 · ponow.1 │
                     ├─────────────────────────────┤
                     │ ✓ 1 npm ci        Bash      │
                     │ ✓ 2 npm run build Bash      │
                     │ ● 3 rsync → web-01 SSH      │
                     │ ○ 4 rsync → web-02 SSH      │
                     │ ○ 5 kontrola kodów wyjścia  │
                     ├─────────────────────────────┤
                     │ Koordynator → Wykonawca:    │
                     │  „zadanie 3 przydzielone”   │
                     │ Wykonawca → Koordynator:    │
                     │  „postęp 62%, bez błędów”   │
                     ├─────────────────────────────┤
                     │ [ Przebieg ▾ ]           ⋮  │
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok tekstowy z treścią zlecenia i wskaźnikami | Pokazuje przedmiot pętli i przebieg jej realizacji | Średni blok z etykietą i licznikami | w toku · wstrzymany · zakończony · przerwany | Kliknięcie rozwija pełną treść zlecenia i jego historię korekt | Góra kolumny pętli | 1 | Widoczny bez interakcji po otwarciu kolumny pętli |
| Wiersz zadania | Wiersz kolejki z ikoną stanu, poleceniem i kartą docelową | Reprezentuje jedno zadanie powłokowe pętli | Średni wiersz `.dn-tabela` z kropką stanu | oczekujące · przydzielone · w toku (kropka pulsująca) · sukces · błąd · ponawiane | Kliknięcie otwiera szczegóły zadania i przejścia do karty, wyniku i procesu | Kolejka zadań | 1 | Widoczny bez interakcji — kolejka zadań |
| Plakietka kontroli jakości | Mała plakietka `.dn-plakietka` z kodem wyjścia | Prezentuje wynik kontroli uruchomienia | Mała pigułka, kolor zależny od wyniku | sukces (zielona) · ostrzeżenie (żółta) · błąd (czerwona) | Kliknięcie pokazuje rozpoznany wzorzec błędu i proponowaną korektę | Prawa krawędź wiersza zadania | 1 | Widoczna bez interakcji przy wierszu zadania |
| Strumień komunikatów sterujących | Lista wiadomości Koordynator ↔ Wykonawca | Pokazuje przydziały, potwierdzenia, raporty postępu i żądania decyzji | Średni blok tekstowy, przewijany | domyślny · oczekiwanie na decyzję (wyróżnienie) | Żądanie decyzji ujawnia przyciski wyboru przekazywane Koordynatorowi | Środek kolumny pętli | 1 | Widoczny bez interakcji |
| Wskaźniki przebiegu | Zestaw liczników i paska postępu | Liczba zadań, czas przebiegu, ponowienia, współbieżność | Małe liczniki z paskiem postępu | aktualizowany na żywo · zatrzymany | — (informacyjny) | Nagłówek zlecenia | 1 | Widoczne bez interakcji |
| Przyciski sterowania przebiegiem | Grupa przycisków: wstrzymaj, wznów, przerwij | Sterowanie pętlą wykonawczą | Małe przyciski tekstowe, „Przerwij” w wariancie `--blad` | domyślny · ładowanie (sygnał wysłany) | Wstrzymanie zatrzymuje przydzielanie kolejnych zadań; przerwanie kończy zadania w toku i odnotowuje przerwanie w Process Monitor | Dolna część kolumny pętli | 2 | Element zbiorczy `Przebieg ▾` w stopce kolumny; zwija się po wydaniu polecenia |
| Przycisk „Koryguj zlecenie” | Przycisk `--zarys`, mały | Zmiana treści zlecenia bez zamykania pętli | Mały przycisk tekstowy | domyślny · edycja otwarta | Otwiera pole edycji zlecenia; Koordynator przelicza dekompozycję na zadania | Dolna część kolumny pętli | 3 | Menu ⋮ nagłówka zlecenia |
| Graf zależności zadań | Diagram kroków potoku | Pokazuje kolejność, warunki i gałęzie równoległe | Średni diagram, zwijany | zwinięty · rozwinięty | Kliknięcie kroku ustawia fokus na odpowiadającym zadaniu w kolejce | Pod kolejką zadań | 3 | Menu ⋮ nagłówka zlecenia albo polecenie „pokaż graf zadań”; panel zwija się po zamknięciu |

---

### 3.3. Terminal Tabs

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne (interaktywna sesja wiersza poleceń) |
| Waga wizualna | Prawa kolumna, dominująca — pasek kart nad obszarem powłoki |
| Izolacja domyślna | Każda karta działa jako odrębna sesja powłoki; katalog roboczy i środowisko procesu współdzielone w obrębie karty sesji, chyba że włączona izolacja techniczna (rozdz. 5.3) |
| Warstwa widoczności | Warstwa 1 — wiodące okno obszaru roboczego, widoczne bez interakcji |

**Zawartość i pełny arsenał funkcji.**

*Powłoki i sesje.*
- Otwarcie nowej karty z menu wyboru powłoki: PowerShell, CMD, Bash, Zsh, Fish, Git Bash, Node.js (REPL), Python (REPL), Deno, `dotnet script`, SQL, sesja SSH do hosta zdalnego, powłoka kontenera (`docker exec`), powłoka poda (`kubectl exec`), konsola szeregowa, sesja Telnet.
- Profile powłok — zapisane konfiguracje: typ powłoki, katalog startowy, zmienne środowiskowe, polecenie startowe, motyw, ikona, nazwa profilu.
- Multipleksacja paneli: podział karty na siatkę sąsiadujących paneli z zagnieżdżaniem, powiększeniem pojedynczego panelu i zapisywanymi układami.
- Pasek rozgłaszania — jedno wpisywane polecenie trafia równolegle do zaznaczonych paneli i hostów.
- Duplikowanie karty (nowa karta z tym samym katalogiem roboczym i powłoką).
- Trwałe sesje nazwane — sesja żyje po stronie serwera po rozłączeniu klienta, ponowne podpięcie przywraca bufor ekranu i stan.

*Zarządzanie kartami.*
- Zmiana nazwy karty, przypinanie karty, zmiana kolejności przez przeciąganie.
- Zamknięcie karty; ostrzeżenie wizualne (plakietka i okno nakładkowe potwierdzenia), gdy w karcie działa proces — sygnalizacja wagi decyzji, nie blokada zamknięcia.
- Zamknięcie pozostałych kart, przeniesienie karty do nowej kolumny obszaru roboczego.
- Wskaźnik stanu karty: aktywna, w tle, proces działa, proces zakończony, błąd.
- Wskaźnik typu połączenia karty: lokalne, SSH, kontener, pod, port szeregowy, Telnet.

*Praca w powłoce.*
- Emulacja terminala VT100/xterm: sekwencje ANSI, kolor 24‑bitowy, tryb aplikacyjny myszy.
- Integracja powłoki: znaczniki OSC 7 (katalog roboczy) i OSC 133 (granice poleceń, kod wyjścia) przekazywane automatycznie z powłoki.
- Bloki poleceń — każde polecenie i jego wynik jako osobny, składany blok z akcjami: kopiuj, ponów, udostępnij, oznacz.
- Autouzupełnianie inline ścieżek, poleceń i flag oraz podpowiedź z historii.
- Historia poleceń per karta (strzałki góra/dół), wyszukiwanie wstecz w historii.
- Zaznaczanie bloków tekstu, kopiowanie i wklejanie, przeciągnij‑i‑upuść pliku jako ścieżki.
- Wyszukiwanie w obrębie bieżącej karty (Grep lokalny okna powłoki) z wyrażeniami regularnymi, podświetleniem i licznikiem trafień.
- Breadcrumb katalogu roboczego z możliwością szybkiej zmiany.
- Nagłówek karty prezentujący host, katalog roboczy, gałąź repozytorium Git i nazwę profilu.
- Paleta poleceń (`Ctrl/Cmd + K`) — szybkie akcje, wyszukiwanie profili, snippetów, hostów i zadań.

*Personalizacja.*
- Zmiana rozmiaru czcionki, gęstości wierszy, ligatur i interlinii.
- Schemat kolorystyczny terminala, zsynchronizowany z motywem jasnym/ciemnym platformy lub ustawiony niezależnie.
- Zapis układu kart i paneli jako migawki karty sesji, przywracany przy wznowieniu sesji.

*Eksport i integracja.*
- Eksport transkryptu karty do pliku tekstowego.
- Przekazanie zawartości karty do Output Console w widoku pojedynczej karty.
- Ponowne uruchomienie bieżącego polecenia z modyfikacją parametrów.
- Zapis polecenia lub bloku poleceń do Script Library jako skryptu albo snippetu.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Terminal Tabs — prawa kolumna, dominująca     │ ☰
  Użytkownik ↔ │ ● PowerShell ×│ ○ Bash ×│  [+ ▾]              │
  Wykonawca    ├───────────────────────────────────────────────┤
               │ [Danaco Console][Ubuntu][build-01 ▾][Bash ▾]  │
  ──────────── │ ~/danaco/api · main                           │
  Execution    ├───────────────────────────────────────────────┤
  Loop ▸       │ PS ~/danaco/api> npm run build                │
               │ > kompilacja w toku…                          │
               │ PS ~/danaco/api>▌                             │
               │                                               │
               ├───────────────────────────────────────────────┤
               │                            [ Operacje ▾ ]  ⋮  │
  [Procesy ▾] [Host ▾]
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Karta powłoki | Zakładka `.dn-zakladki` z ikoną typu powłoki | Reprezentuje jedną, odrębną sesję powłoki | Mała zakładka z etykietą i kropką stanu | aktywna (podkreślenie sygnałowe) · w tle · proces działa (kropka pulsująca) · błąd (kropka czerwona) | Kliknięcie przełącza fokus; przeciągnięcie zmienia kolejność | Pasek kart nad obszarem powłoki | 1 | Widoczna bez interakcji — pasek kart |
| Przycisk „Nowa karta” | Przycisk ikonowy `.dn-btn-ikona` (plus) z rozwijanym menu | Otwarcie nowej sesji powłoki wybranego typu | Mała ikona z małą strzałką rozwinięcia | domyślny · rozwinięte menu wyboru powłoki | Otwiera nową kartę z wybraną powłoką i katalogiem startowym | Lewy kraniec paska kart | 2 | Przycisk `[+ ▾]` paska kart; lista powłok zwija się po wyborze |
| Ikona zamknięcia karty | Mały „×” na karcie | Zamknięcie sesji powłoki | Mała ikona (zamknij), widoczna przy wskazaniu kursorem na kartę | domyślny · potwierdzenie wymagane (proces aktywny) | Zamyka kartę natychmiast; gdy proces działa, otwiera `.dn-modal` z przyciskiem `--blad` sygnalizującym nieodwracalność, nie blokującym akcji | Prawa strona etykiety karty | 2 | Ujawnia się przy wskazaniu kursorem na kartę |
| Wskaźnik typu połączenia | Mała ikona przy etykiecie karty | Rozróżnia połączenie lokalne, SSH, kontener, pod, port szeregowy, Telnet | Bardzo mała ikona | domyślny · rozłączony (wyszarzony) | Kliknięcie otwiera wpis hosta w Session Manager | Lewa strona etykiety karty | 1 | Widoczny bez interakcji przy etykiecie karty |
| Breadcrumb ścieżki | Pasek stanu z segmentami katalogu | Pokazuje i pozwala zmienić katalog roboczy karty | Mały pasek tekstowy pod paskiem kart | domyślny · edytowalny (klik w segment) | Kliknięcie segmentu wykonuje `cd` do tego katalogu | Pasek stanu karty | 1 | Widoczny bez interakcji — pasek kontekstu karty |
| Wskaźnik gałęzi Git | Mała plakietka `.dn-plakietka` z ikoną gałęzi | Informuje o aktywnej gałęzi repozytorium w katalogu roboczym | Mała pigułka tekstowa | ukryty (katalog spoza repozytorium) · widoczny | Kliknięcie otwiera Git Panel modułu Developer dla tego repozytorium | Pasek stanu karty, obok breadcrumb | 1 | Widoczny bez interakcji jako znacznik kontekstowy w pasku karty |
| Obszar powłoki | Emulator terminala | Właściwa przestrzeń wpisywania i odczytu poleceń | Duży, dominujący blok tekstu mono, przewijany | gotowa do wpisu · polecenie w toku (kursor zajęty) · zablokowana do odczytu (podgląd historii) | Wpisywanie poleceń, zaznaczanie, kopiowanie | Środek prawej kolumny | 1 | Widoczny bez interakcji — ponad 80% powierzchni obszaru roboczego |
| Blok polecenia | Składany blok obejmujący polecenie i jego wynik | Wydziela pojedyncze uruchomienie z ciągłego strumienia | Średni blok z paskiem akcji | zwinięty · rozwinięty · oznaczony | Pasek akcji udostępnia: kopiuj, ponów, udostępnij, zapisz do Script Library | Obszar powłoki | 1 | Widoczny bez interakcji; pasek akcji bloku należy do warstwy 3 i otwiera się z menu ⋮ bloku |
| Przycisk podziału panelu | Para przycisków ikonowych | Podział karty na dodatkowy panel powłoki w siatce | Małe ikony w pasku narzędzi karty | domyślny · aktywny (podział zastosowany) | Dzieli obszar karty na niezależne panele powłoki | Pasek narzędzi karty | 3 | Element zbiorczy `Operacje ▾` w pasku kontekstu karty |
| Przełącznik rozgłaszania | Przycisk przełączający z licznikiem paneli | Kierowanie wpisu do wielu paneli i hostów jednocześnie | Mały przycisk tekstowy z plakietką liczby celów | wyłączony · włączony (obramowanie ostrzegawcze) | Wpis z klawiatury trafia do wszystkich zaznaczonych paneli | Pasek narzędzi karty | 4 | Polecenie języka naturalnego „rozgłaszaj polecenie na zaznaczone panele”, skrót klawiszowy albo wyszukiwarka funkcji `Ctrl/Cmd + K` |
| Selektor schematu kolorów | Rozwijana lista | Zmiana palety kolorystycznej terminala | Mały przycisk z etykietą | domyślny (zgodny z motywem platformy) · niestandardowy | Zmienia kolorystykę tekstu i tła obszaru powłoki | Pasek narzędzi karty | 3 | Menu ⋮ karty, sekcja prezentacji; zwija się po wyborze |
| Regulator wielkości czcionki | Para małych przycisków „A–”/„A+” | Zmniejszenie i zwiększenie rozmiaru czcionki terminala | Bardzo małe przyciski ikonowe | domyślny | Skaluje tekst we wszystkich kartach | Pasek narzędzi karty | 3 | Menu ⋮ karty, sekcja prezentacji |

---

### 3.4. Output Console

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor procesu / podgląd |
| Waga wizualna | Prawa kolumna, dominująca, obok Terminal Tabs |
| Izolacja domyślna | Wynik grupowany per karta sesji; widok zagregowany łączy wyłącznie karty tej samej sesji |
| Warstwa widoczności | Warstwa 1 — druga część obszaru roboczego, widoczna bez interakcji |

**Zawartość i pełny arsenał funkcji.**

- Widok zagregowany — strumień wyjścia ze wszystkich otwartych kart Terminal Tabs, ze znacznikiem karty źródłowej i hosta przy każdym bloku.
- Widok pojedynczej karty — filtr zawężający strumień do jednej sesji powłoki.
- Filtrowanie po poziomie: standardowe wyjście, standardowy błąd, ostrzeżenie, błąd (rozpoznane wzorcem).
- Wyszukiwanie pełnotekstowe (Grep) z obsługą wyrażeń regularnych, z podświetleniem i licznikiem trafień.
- Podświetlanie składni wyjścia (kolory ANSI konsoli zachowane i zinterpretowane).
- Parser strukturalny wyniku — rozpoznanie JSON, YAML i tabel oraz przełącznik prezentacji: surowe, tabela, drzewo, JSON.
- Panel różnicowania — porównanie wyjścia dwóch uruchomień obok siebie z kolorowym oznaczeniem zmian.
- Znaczniki czasu przy każdej linii (włączane i wyłączane).
- Zawijanie linii lub przewijanie poziome (przełącznik).
- Automatyczne przewijanie do najnowszego wpisu, z możliwością zatrzymania przy przeglądaniu historii wstecz.
- Limit buforowanych linii, konfigurowalny per karta sesji.
- Rejestracja sesji jako odtwarzalnego nagrania z osią czasu oraz odtwarzanie krok po kroku lub w przyspieszeniu, z paskiem transportu.
- Maskowanie wrażliwych wartości — zaciemnianie tokenów i kluczy w wyniku według wzorców, także w transkrypcie i nagraniu.
- Wykrywanie wzorców błędów i wyjątków — automatyczne oznaczenie linii plakietką błędu.
- Oznaczenie fragmentu jako kontekst dla Chat Window („wyjaśnij”, „napraw”).
- Przekazanie fragmentu jednym kliknięciem do modułu Diagnostics jako materiał źródłowy.
- Podział widoku — jednoczesny podgląd dwóch kart obok siebie.
- Eksport wybranego zakresu wyjścia do plików `.txt`, `.log`, `.md` oraz `.html` z zachowaniem kolorów.
- Kopiowanie zaznaczonego fragmentu.
- Wyczyszczenie widoku (dotyczy wyłącznie prezentacji — nie usuwa danych z rejestru procesu).
- Ponowne uruchomienie ostatniego polecenia bezpośrednio z widoku wyniku.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Output Console — prawa kolumna, dominująca    │ ☰
  Użytkownik ↔ │ [Zagregowany ▾] [Format ▾]             🔍  ⋮  │
  Wykonawca    ├───────────────────────────────────────────────┤
               │ [Bash]      12:04:01  $ npm run build         │
  ──────────── │ [Bash]      12:04:03  > kompilacja (2,1 s)    │
  Execution    │ [PowerShell]12:04:07  PS> Get-Process node    │
  Loop ▸       │ [Python]    12:04:12  Traceback … ⚠ błąd      │
               │        [ Wyjaśnij ][ Przekaż do Diagnostics ] │
               │                                               │
               ├───────────────────────────────────────────────┤
               │                            [ Operacje ▾ ]     │
  [Procesy ▾] [Host ▾]
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor widoku | Rozwijana lista | Przełącza między widokiem zagregowanym a widokiem jednej karty | Mały przycisk z etykietą | zagregowany (domyślny) · pojedyncza karta | Filtruje strumień do wybranego źródła | Nagłówek okna | 2 | Znacznik kontekstowy `[Zagregowany ▾]` w nagłówku; zwija się po wyborze |
| Selektor formatu wyniku | Rozwijana lista | Przełącza prezentację: surowe, tabela, drzewo, JSON | Mały przycisk z etykietą | surowe (domyślny) · struktura rozpoznana · struktura nierozpoznana (wyszarzony) | Renderuje rozpoznaną strukturę jako tabelę lub drzewo | Nagłówek okna | 2 | Znacznik kontekstowy `[Format ▾]` w nagłówku; zwija się po wyborze |
| Pole Grep | Pole wyszukiwania z przełącznikami | Wyszukiwanie wzorca w zgromadzonym wyjściu | Pole `.dn-pole` + ikony‑przełączniki (regex, wielkość liter) | puste · z wynikami · brak trafień | Enter uruchamia wyszukiwanie; podświetla i zlicza trafienia | Nagłówek okna | 2 | Ikona 🔍 w nagłówku albo skrót wyszukiwania; pole zwija się po wyczyszczeniu zapytania |
| Pasek transportu odtwarzania | Zestaw przycisków i suwak osi czasu | Odtwarzanie nagrania sesji krok po kroku lub w przyspieszeniu | Mały pasek z suwakiem | zatrzymany · odtwarzanie · przewijanie | Ustawia pozycję nagrania i tempo odtwarzania | Nagłówek okna | 3 | Menu ⋮ okna, pozycja odtwarzania nagrania; pasek znika po zamknięciu nagrania |
| Znacznik źródła bloku | Mała plakietka z nazwą karty i hosta | Wskazuje, z której karty i hosta pochodzi dana linia wyniku | Mała pigułka `.dn-plakietka` | domyślny | Kliknięcie przełącza widok na tę kartę | Przed każdym blokiem wyniku | 1 | Widoczny bez interakcji przed blokiem wyniku |
| Linia wyniku | Wiersz tekstu mono | Pojedynczy wpis wyjścia polecenia | Tekst `--dn-ff-mono`, tło zależne od poziomu | standardowe wyjście · ostrzeżenie (żółte obramowanie) · błąd (czerwone obramowanie, ikona `blad`) · wartość zamaskowana | Zaznaczenie fragmentu ujawnia pasek akcji kontekstowych | Ciało konsoli | 1 | Widoczna bez interakcji — zasadnicza treść okna |
| Pasek akcji linii błędu | Mały pasek pod linią błędu | Szybkie przejście do wyjaśnienia lub diagnozy | Dwa małe przyciski tekstowe | ukryty · widoczny (po wykryciu błędu) | „Wyjaśnij” otwiera Chat Window z kontekstem; „Przekaż do Diagnostics” tworzy wpis w Errors Panel | Pod linią rozpoznaną jako błąd | 2 | Ujawnia się po rozpoznaniu błędu albo po zaznaczeniu fragmentu; znika po wykonaniu akcji |
| Przycisk „Różnicuj” | Przycisk `--zarys`, mały | Porównanie wyniku dwóch uruchomień | Mały przycisk tekstowy | domyślny · aktywny (panel różnic otwarty) | Otwiera panel z zestawieniem dwóch przebiegów obok siebie | Pasek narzędzi okna | 3 | Element zbiorczy `Operacje ▾` w stopce okna |
| Przełączniki prezentacji | Trzy pola `.dn-check` | Zawijanie linii, znaczniki czasu, auto‑przewijanie | Małe pola wyboru z etykietą | zaznaczony · niezaznaczony | Zmienia natychmiast sposób prezentacji strumienia | Pasek narzędzi okna | 3 | Menu ⋮ okna, sekcja prezentacji strumienia |
| Przycisk „Eksportuj” | Przycisk `--zarys`, mały | Zapis widocznego lub całego zakresu wyniku do pliku | Mały przycisk tekstowy | domyślny · ładowanie | Generuje plik `.txt`, `.log`, `.md` albo `.html` i uruchamia pobranie | Pasek narzędzi okna, prawa strona | 3 | Element zbiorczy `Operacje ▾` w stopce okna |
| Przycisk „Wyczyść” | Przycisk `--duch`, mały | Czyszczenie widoku konsoli | Mała ikona (kosz) | domyślny | Czyści wyłącznie widok; dane procesu pozostają w rejestrze i historii sesji | Nagłówek okna | 3 | Menu ⋮ okna |

---

### 3.5. Process Monitor

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor procesu |
| Waga wizualna | Wąskie okno pomocnicze po prawej stronie okna centralnego, otwierane na żądanie, aktualizowane na żywo |
| Izolacja domyślna | Rejestr obejmuje procesy karty sesji bieżącej; widoczność procesów innych kart konfigurowalna |
| Warstwa widoczności | Warstwa 2 — otwierany na żądanie ikoną szybkiego dostępu w prawym górnym rogu okna centralnego aplikacji albo poleceniem „pokaż procesy”; okno pomocnicze znika po zamknięciu i nigdy nie jest stale widoczne |

**Reguła miejsca.** Wąskie okno po prawej stronie okna centralnego jest miejscem
wszystkich narzędzi pomocniczych platformy; Process Monitor jest jednym z nich.
Narzędzie pomocnicze otwiera się wyłącznie na żądanie, z komponentu (ikony)
w prawym górnym rogu okna centralnego, i nie zajmuje przestrzeni roboczej,
dopóki Operator go nie wywoła.

**Zawartość i pełny arsenał funkcji.**

- Tabela aktywnych i zakończonych procesów: identyfikator uruchomienia, polecenie, karta źródłowa, host, inicjator (Operator albo Wykonawca), czas uruchomienia, czas trwania, zużycie procesora, zużycie pamięci, operacje wejścia‑wyjścia, status.
- Rozróżnienie wizualne procesów zainicjowanych przez Wykonawcę (odrębna plakietka) — pełny wgląd w to, co zostało uruchomione bez bezpośredniego wpisu Operatora.
- Powiązanie wiersza procesu z zadaniem pętli wykonawczej, które go zleciło.
- Filtrowanie: aktywne, zakończone, zakończone błędem, uruchomione przez Wykonawcę, uruchomione przez Operatora.
- Wyszukiwanie po nazwie polecenia lub identyfikatorze procesu.
- Podgląd szczegółów procesu: pełne polecenie z parametrami, katalog roboczy, zmienne środowiskowe, drzewo procesów potomnych.
- Drzewo procesów — hierarchia proces nadrzędny/potomny, jak proces menedżera pakietów uruchamiający proces środowiska wykonawczego.
- Akcja: zakończenie procesu (pojedynczego albo grupy zaznaczonych), z sygnałem łagodnym lub wymuszonym.
- Akcja: wstrzymanie i wznowienie procesu, gdy powłoka i system to wspierają.
- Akcja: ponowne uruchomienie ostatniego zakończonego polecenia z tymi samymi parametrami.
- Mini‑wykres obciążenia procesora, pamięci i operacji wejścia‑wyjścia, na żywo, per proces oraz zbiorczo dla karty sesji.
- Przypięcie procesu jako obserwowanego — pozostaje na górze listy niezależnie od sortowania.
- Powiadomienie o zakończeniu długo działającego procesu, z progiem czasu ustawianym w oknie konfiguracji.
- Sortowanie kolumn: czas uruchomienia, czas trwania, zużycie zasobów, status.
- Grupowanie: po karcie źródłowej, po hoście, po inicjatorze, po statusie.
- Eksport migawki listy procesów.
- Odświeżanie na żywo kanałem WebSocket, ze wskaźnikiem stanu połączenia z rejestrem procesów rdzenia serwera.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Obszar roboczy modułu   │ Process Monitor (okno pomocnicze)
  Użytkownik ↔ │ Terminal Tabs           │ [Filtr ▾] [Grupuj ▾]        ● ⋮
  Wykonawca    │ Output Console          ├──────────────────────────────────
               │                         │ ● npm run dev    Bash   Operator
  ──────────── │                         │   CPU ▂▃▅▂ RAM 210 MB  00:03:12
  Execution    │                         │   [ Zakończ ]                 ⋮
  Loop ▸       │                         ├──────────────────────────────────
               │                         │ ● pytest -k api  Python  Wykonawca
               │                         │   CPU ▁▂▇▃ RAM 96 MB   00:00:41
               │                         │   zadanie pętli: 4  [ Zakończ ] ⋮
               │                         ├──────────────────────────────────
               │                         │ ○ git fetch --all  PS  Operator
               │                         │   kod wyjścia: 0              ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Wiersz procesu | Wiersz tabeli `.dn-tabela` z kropką stanu | Reprezentuje jeden proces w rejestrze | Średni wiersz z kolumnami danych i mini‑wykresem | uruchomiony (kropka zielona pulsująca) · wstrzymany (kropka żółta) · zakończony sukcesem (kropka szara) · zakończony błędem (kropka czerwona) | Kliknięcie rozwija szczegóły procesu i drzewo potomków | Ciało tabeli | 1 | Widoczny bez interakcji po otwarciu kolumny |
| Plakietka inicjatora | Mała plakietka `.dn-plakietka` z ikoną | Rozróżnia proces uruchomiony przez Operatora od uruchomionego przez Wykonawcę | Mała pigułka z etykietą „Operator” albo „Wykonawca” | domyślny | Kliknięcie plakietki „Wykonawca” otwiera w Chat Window wiadomość, która zainicjowała proces | Kolumna „Inicjator” | 1 | Widoczna bez interakcji w wierszu procesu |
| Odnośnik zadania pętli | Mała etykieta z numerem zadania | Wiąże proces z zadaniem pętli wykonawczej | Bardzo mała etykieta tekstowa | ukryty (proces spoza pętli) · widoczny | Kliknięcie ustawia fokus na zadaniu w Execution Loop Window | Wiersz procesu, pod poleceniem | 1 | Widoczny bez interakcji dla procesów zleconych przez pętlę |
| Mini‑wykres obciążenia | Mały wykres typu sparkline | Pokazuje trend zużycia procesora, pamięci i operacji wejścia‑wyjścia w czasie | Bardzo mały wykres liniowy, wysokość jednego wiersza | aktualizowany na żywo · płaski (proces zakończony) | Wskazanie kursorem pokazuje dokładne wartości chwilowe | Kolumna obciążenia w wierszu procesu | 1 | Widoczny bez interakcji w wierszu procesu |
| Przycisk „Zakończ” | Przycisk `--blad`, mały | Zatrzymanie działającego procesu | Mały przycisk tekstowy w kolorze sygnalizującym wagę akcji | domyślny · ładowanie (sygnał wysłany) | Otwiera `.dn-modal` z potwierdzeniem dla procesów z niezapisanym stanem; po potwierdzeniu kończy proces natychmiast — potwierdzenie sygnalizuje nieodwracalność, nie blokuje wykonania | Kolumna akcji, wiersz procesu aktywnego | 2 | Ujawnia się w wierszu procesu aktywnego |
| Przycisk „Uruchom ponownie” | Przycisk `--zarys`, mały | Ponowne wykonanie zakończonego polecenia | Mały przycisk tekstowy | domyślny | Wstawia i uruchamia to samo polecenie w karcie źródłowej | Kolumna akcji, wiersz procesu zakończonego | 3 | Menu ⋮ wiersza procesu zakończonego |
| Menu „⋮” wiersza | Ikona rozwijanego menu `wiecej` | Dodatkowe akcje: przypnij, eksportuj wpis, przejdź do karty źródłowej, wstrzymaj | Mała ikona 24 px | domyślny · rozwinięte | Otwiera menu kontekstowe wiersza | Prawa krawędź wiersza | 3 | Ikona ⋮ przy krawędzi wiersza; menu kontekstowe zwija się po wyborze akcji |
| Panel szczegółów procesu | Rozwijany blok pod wierszem | Pełne polecenie, katalog roboczy, zmienne środowiskowe, drzewo potomków | Średni panel, tekst mono dla poleceń | zwinięty · rozwinięty | Rozwija się po kliknięciu wiersza; kod poleceń czcionką `--dn-ff-mono` | Pod wierszem procesu | 3 | Rozwinięcie kontekstowe wiersza — kliknięcie wiersza; zwija się ponownym kliknięciem |
| Filtr i grupowanie | Dwie rozwijane listy w nagłówku | Zawężenie i porządkowanie listy procesów | Małe przyciski z etykietą | domyślny (wszystkie / karta) · zmieniony | Przelicza widoczną i pogrupowaną listę | Nagłówek kolumny | 2 | Znaczniki kontekstowe `[Filtr ▾] [Grupuj ▾]` w nagłówku kolumny; zwijają się po wyborze |
| Wskaźnik połączenia na żywo | Mała kropka statusu przy nagłówku | Sygnalizuje aktywność kanału WebSocket z rejestrem procesów | Bardzo mała kropka `.dn-kropka` | połączony (zielona) · rozłączony (szara, dane ostatnie znane) | — (informacyjny) | Nagłówek kolumny | 1 | Widoczny bez interakcji w nagłówku kolumny |

---

### 3.6. Session Manager

| Aspekt | Wartość |
|---|---|
| Typologia | Nawigacja / katalog |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne, zwijana |
| Izolacja domyślna | Książka hostów wspólna dla karty sesji; osobne poświadczenia per karta ustawiane zakresem izolacji „konto i token per sesja” |
| Warstwa widoczności | Warstwa 2 — otwierany znacznikiem `[Host ▾]` paska kontekstu albo z menu ☰ modułu; kolumna zwija się po wybraniu hosta |

**Zawartość i pełny arsenał funkcji.**

- Drzewo hostów z folderami i tagami: dane połączenia, port, użytkownik, klucz, jump‑host, notatka.
- Szybkie łączenie jednym kliknięciem: SSH, konsola szeregowa, Telnet, kontener, pod.
- Import i eksport pliku `~/.ssh/config`.
- Lista sesji trwałych z podglądem stanu oraz akcjami podpięcia i odpięcia (attach/detach).
- Menedżer profili powłok — tworzenie, edycja, kolejność, profil domyślny per środowisko i projekt.
- Menedżer kluczy SSH: generowanie, import, powiązanie klucza z hostem, obsługa agenta.
- Panel tuneli portowych: przekierowania lokalne, zdalne i dynamiczne (SOCKS) ze stanem połączenia i przepustowością.
- Grupy hostów dla operacji wielohostowych — jedna operacja na całej grupie, agregacja wyników per host.
- Polityka `known_hosts` i podgląd odcisków kluczy hostów.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Obszar roboczy modułu   │ Session Manager (kolumna boczna)
  Użytkownik ↔ │ Terminal Tabs           │ 🔍                          ☰
  Wykonawca    │ Output Console          ├──────────────────────────────────
               │                         │ ▾ Produkcja
  ──────────── │                         │    ● web-01   ssh
  Execution    │                         │    ● web-02   ssh
  Loop ▸       │                         │ ▸ Laboratorium
               │                         ├──────────────────────────────────
               │                         │ Sesje trwałe
               │                         │  ● build-nocny   [ Podepnij ]
               │                         ├──────────────────────────────────
               │                         │ [ Tunele ▾ ]                  ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Węzeł hosta | Wiersz drzewa z ikoną typu połączenia | Reprezentuje jeden wpis książki hostów | Mały wiersz z etykietą i kropką dostępności | rozłączony · połączony · niedostępny | Kliknięcie otwiera nową kartę powłoki z tym hostem | Drzewo hostów | 1 | Widoczny bez interakcji po otwarciu kolumny |
| Folder i tag | Węzeł grupujący z etykietą | Porządkuje hosty według środowiska i roli | Mały węzeł zwijany | zwinięty · rozwinięty | Rozwija listę hostów; menu kontekstowe uruchamia operację na grupie | Drzewo hostów | 1 | Widoczny bez interakcji; operacje na grupie w menu kontekstowym (warstwa 3) |
| Wiersz sesji trwałej | Wiersz listy z nazwą sesji i czasem życia | Pokazuje sesję żyjącą po stronie serwera | Mały wiersz z przyciskiem akcji | odpięta · podpięta · zakończona | „Podepnij” przywraca bufor i stan sesji w nowej karcie | Lista sesji trwałych | 1 | Widoczny bez interakcji — sekcja sesji trwałych |
| Wiersz tunelu | Wiersz z portem lokalnym, celem i stanem | Prezentuje przekierowanie portu | Mały wiersz z kropką stanu i licznikiem ruchu | nieaktywny · aktywny · błąd | Kliknięcie otwiera parametry tunelu i pozwala go zatrzymać lub wznowić | Panel tuneli | 2 | Znacznik `[Tunele ▾]` w stopce kolumny; sekcja zwija się po użyciu |
| Przycisk importu konfiguracji SSH | Przycisk `--zarys`, mały | Wczytanie wpisów z pliku `~/.ssh/config` | Mały przycisk tekstowy | domyślny · ładowanie | Dodaje wpisy do drzewa hostów z zachowaniem folderów | Stopka kolumny | 3 | Menu ☰ kolumny, pozycja importu i eksportu |
| Menedżer kluczy | Sekcja listy kluczy z akcjami | Generowanie, import i przypisanie kluczy do hostów | Średnia lista zwijana | zwinięty · rozwinięty | Otwiera formularz generowania albo importu klucza | Stopka kolumny | 4 | Polecenie języka naturalnego „wygeneruj klucz SSH dla hosta”, wyszukiwarka funkcji `Ctrl/Cmd + K` albo tryb administracyjny |

---

### 3.7. Task & Schedule

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor / edytor zadań |
| Waga wizualna | Prawa kolumna, dominująca — zakładka obszaru roboczego |
| Izolacja domyślna | Zadania wykryte w katalogu roboczym karty sesji; zasięg skanowania i strefa czasowa harmonogramu ustawiane w oknie konfiguracji |
| Warstwa widoczności | Warstwa 3 — zakładka obszaru roboczego wywoływana z elementu zbiorczego `Operacje ▾` |

**Zawartość i pełny arsenał funkcji.**

- Wykryte zadania projektu z manifestów `package.json`, `Makefile`, `Taskfile.yml` i `justfile`, każde z przyciskiem uruchomienia.
- Lista zadań cyklicznych z wyrażeniem cron i kalendarzem nadchodzących uruchomień.
- Wyzwalacze plikowe — uruchomienie polecenia przy zmianie plików i katalogów, z wzorcami glob i tłumieniem powtórzeń.
- Wyzwalacze zdarzeniowe — webhook, wpis w module Diagnostics, zakończenie innego zadania; realizowane jawnym powiązaniem z modułem Automations.
- Edytor potoku: kroki, warunki sukcesu i porażki, zmienne, gałęzie równoległe, limit współbieżności.
- Kolejka przebiegów i historia: status, czas, wynik, akcja ponowienia.
- Bramy wykonania równoległego — uruchomienie wielu zadań z limitem współbieżności i agregacją wyników.
- Powiadomienia o zakończeniu i o błędzie, z progiem czasu i wyborem kanału.
- Przekazanie zadania do pętli wykonawczej jako zlecenia obsługiwanego przez Koordynatora.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Task & Schedule — prawa kolumna, dominująca   │ ☰
  Użytkownik ↔ │ [Zadania projektu][Harmonogram]               │
  Wykonawca    ├───────────────────────────────────────────────┤
               │ package.json  build      [ Uruchom ]        ⋮ │
  ──────────── │ package.json  test       [ Uruchom ]        ⋮ │
  Execution    │ Makefile      deploy     [ Uruchom ]        ⋮ │
  Loop ▸       ├───────────────────────────────────────────────┤
               │ cron 0 3 * * *  kopia zapasowa  → jutro 03:00 │
               │ [ Wyzwalacze ▾ ]                              │
               ├───────────────────────────────────────────────┤
               │ Historia: ✓ build 12:04 (2,1 s) ✗ test 12:07  │
               │                             [ Operacje ▾ ]  ⋮ │
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Wiersz zadania projektu | Wiersz listy ze źródłem manifestu i nazwą zadania | Uruchomienie zadania wykrytego w projekcie | Mały wiersz z przyciskiem `--zarys` | domyślny · w toku · zakończone | Uruchamia zadanie w karcie Terminal Tabs, wynik trafia do Output Console | Zakładka „Zadania projektu” | 1 | Widoczny bez interakcji po otwarciu zakładki |
| Wiersz harmonogramu | Wiersz z wyrażeniem cron i czasem następnego uruchomienia | Zarządzanie zadaniem cyklicznym | Średni wiersz z kropką stanu | włączone · wyłączone · uruchamiane | Kliknięcie otwiera edytor wyrażenia cron i parametrów zadania | Zakładka „Harmonogram” | 1 | Widoczny bez interakcji — sekcja harmonogramu |
| Wiersz wyzwalacza | Wiersz ze wzorcem obserwacji lub zdarzeniem | Uruchamianie polecenia na zmianę plików albo zdarzenie | Średni wiersz z kropką stanu | aktywny · zatrzymany · błąd | Kliknięcie otwiera parametry: wzorzec, tłumienie powtórzeń, polecenie | Zakładka „Wyzwalacze” | 2 | Znacznik `[Wyzwalacze ▾]`; sekcja zwija się po wyborze pozycji |
| Edytor potoku | Diagram kroków z formularzem parametrów | Budowa sekwencji poleceń z warunkami i równoległością | Duży obszar edycyjny | edycja · walidacja · zapisany | Zapis tworzy potok uruchamialny i widoczny w historii przebiegów | Zakładka „Potoki” | 3 | Element zbiorczy `Operacje ▾`, pozycja edycji potoku; obszar znika po zamknięciu |
| Wiersz przebiegu | Wiersz historii ze statusem, czasem i wynikiem | Przegląd wykonanych przebiegów | Średni wiersz z plakietką statusu | sukces · błąd · przerwany | „Ponów” uruchamia przebieg z tymi samymi parametrami | Zakładka „Historia przebiegów” | 1 | Widoczny bez interakcji — historia przebiegów |
| Przycisk „Przekaż do pętli” | Przycisk `--zarys`, mały | Przekazanie zadania Koordynatorowi jako zlecenia pętli wykonawczej | Mały przycisk tekstowy | domyślny | Otwiera Execution Loop Window ze zleceniem zbudowanym z zadania | Wiersz przebiegu i wiersz zadania | 3 | Menu ⋮ wiersza zadania i wiersza przebiegu |

---

### 3.8. Script Library

| Aspekt | Wartość |
|---|---|
| Typologia | Katalog / edytor |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne; dostępna również z palety poleceń |
| Izolacja domyślna | Biblioteka wspólna dla projektu; zakres widoczności ustawiany warstwowo w oknie konfiguracji |
| Warstwa widoczności | Warstwa 3 — otwierana z menu ☰ modułu, palety poleceń `Ctrl/Cmd + K` albo poleceniem języka naturalnego; kolumna znika po zamknięciu |

**Zawartość i pełny arsenał funkcji.**

- Lista skryptów z tagami, opisem, wersją i czasem ostatniego uruchomienia.
- Edytor skryptów wieloliniowy dla PowerShell, Bash i Python z podświetlaniem składni oraz akcją uruchomienia.
- Parametryzacja skryptu — deklaracja parametrów renderuje formularz wypełniany przed uruchomieniem.
- Snippety poleceń przypisane do skrótu i aliasu, wstawiane z palety poleceń.
- Linter i formatowanie skryptu przed uruchomieniem, z listą uwag i akcją naprawy.
- Zapis polecenia lub bloku poleceń z Terminal Tabs jako nowego skryptu albo snippetu.
- Generowanie skryptu z opisu w języku naturalnym z poziomu Chat Window i zapis do biblioteki.
- Zarządzanie aliasami i funkcjami powłoki z interfejsu, synchronizowane między sesjami.
- Zestaw plików konfiguracyjnych powłoki wersjonowany i przenośny między hostami, zintegrowany z repozytorium przez moduł Developer.

**Makieta tekstowa.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window  │ Obszar roboczy modułu   │ Script Library (kolumna boczna)
  Użytkownik ↔ │ Terminal Tabs           │ 🔍                          ☰
  Wykonawca    │ Output Console          ├──────────────────────────────────
               │                         │ #deploy  wdrozenie.ps1   v3
  ──────────── │                         │ #db      migracja.sh     v7
  Execution    │                         │ #test    smoke.py        v2
  Loop ▸       │                         ├──────────────────────────────────
               │                         │ Edytor: migracja.sh
               │                         │  ✓ linter: bez uwag
               │                         │  [ Uruchom ]                   ⋮
               │                         ├──────────────────────────────────
               │                         │ [ Snippety ▾ ]
 ═══════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Wiersz skryptu | Wiersz listy z tagiem, nazwą i wersją | Reprezentuje skrypt w bibliotece | Mały wiersz z etykietą i plakietką wersji | domyślny · wybrany · uruchamiany | Kliknięcie otwiera skrypt w edytorze kolumny | Lista skryptów | 1 | Widoczny bez interakcji po otwarciu kolumny |
| Obszar edytora | Edytor kodu z podświetlaniem składni | Redagowanie treści skryptu | Średni obszar tekstu mono | edycja · zapisany · z uwagami analizy statycznej | Zapis tworzy nową wersję skryptu w bibliotece | Środek kolumny | 1 | Widoczny bez interakcji po wybraniu skryptu |
| Formularz parametrów | Zestaw pól wygenerowany z deklaracji parametrów | Wypełnienie wartości przed uruchomieniem | Mały formularz nad przyciskiem uruchomienia | pusty · wypełniony · błąd walidacji | Uruchomienie przekazuje wartości do skryptu w karcie powłoki | Pod obszarem edytora | 2 | Ujawnia się dla skryptu z zadeklarowanymi parametrami; znika po uruchomieniu |
| Wskaźnik analizy statycznej | Etykieta z liczbą uwag | Prezentuje wynik statycznej analizy skryptu | Mała etykieta z ikoną | bez uwag · uwagi · błędy składni | Kliknięcie rozwija listę uwag z przejściem do linii | Pod obszarem edytora | 1 | Widoczny bez interakcji pod obszarem edytora; lista uwag rozwija się jako panel warstwy 3 |
| Przycisk „Uruchom” | Przycisk `--sygnal`, mały | Wykonanie skryptu w karcie powłoki | Mały przycisk wyróżniony akcentem sygnałowym | domyślny · ładowanie | Uruchamia skrypt, wynik trafia do Output Console, wpis do Process Monitor | Stopka edytora | 1 | Widoczny bez interakcji w stopce edytora |
| Sekcja snippetów i aliasów | Lista zwijana | Zarządzanie snippetami, aliasami i funkcjami powłoki | Mała lista z akcjami | zwinięta · rozwinięta | Wstawia snippet do aktywnej karty albo otwiera edycję aliasu | Stopka kolumny | 3 | Element zbiorczy `Snippety ▾` w stopce kolumny; zwija się po wstawieniu snippetu |

---

## 4. Przepływy pracy

### 4.1. Przepływ podstawowy — polecenie wydane ręcznie przez Operatora

```
 Terminal Tabs              Output Console            Process Monitor
 ─────────────              ──────────────             ───────────────
 1. wybór/otwarcie
    karty powłoki  ───────►
 2. wpisanie
    polecenia      ──────────────► zapis wyniku
    i Enter                          na żywo
                                        │
                                        ▼
                                 wpis w rejestrze  ─────────────►  wiersz procesu
                                 aktywnych procesów                „uruchomiony”,
                                                                    inicjator: Operator
                                        │                                │
                                        ▼                                ▼
 3. obserwacja wyniku  ◄────────  strumień wyjścia          proces kończy się
    w toku                                                   naturalnie lub jest
                                                               zakończony ręcznie
                                        │                                │
                                        └───────────┬────────────────────┘
                                                    │
                                                    ▼
                                         status „zakończony” w obu oknach
```

### 4.2. Przepływ rozszerzony — polecenie zlecone z Chat Window

```
Chat Window (Użytkownik ↔ Wykonawca)
        │  Użytkownik: „uruchom testy i pokaż wynik”
        ▼
Wykonawca proponuje polecenie ── Użytkownik wybiera:
        │                        │
        │              [Wstaw do terminala]   [Uruchom teraz]
        │                        │                    │
        │                        ▼                    ▼
        │              polecenie czeka         Wykonawca uruchamia
        │              w wierszu poleceń       polecenie bezpośrednio
        │              do ręcznego Enter       w karcie Terminal Tabs
        │                                             │
        └─────────────────────────────────────────────┤
                                                      ▼
                                          Output Console — wynik oznaczony
                                          źródłem karty; Process Monitor —
                                          wpis z plakietką „inicjator: Wykonawca”
                                                      │
                                                      ▼
                                      Użytkownik obserwuje przebieg i, w razie
                                      potrzeby, przerywa działanie z Chat Window
                                      albo kończy proces z Process Monitor
```

### 4.3. Przepływ pętli wykonawczej — zlecenie wielozadaniowe

```
Chat Window                    Execution Loop Window            Terminal Tabs
(Użytkownik ↔ Wykonawca)       (Koordynator ↔ Wykonawca)        Output Console
──────────────────────         ──────────────────────────       ──────────────
 Użytkownik zleca:
 „zbuduj i wdróż na
  dwóch hostach”     ────►  Koordynator dekomponuje
                            zlecenie na zadania powłokowe
                                       │
                                       ▼
                            kolejka zadań: npm ci ·
                            npm run build · rsync web-01 ·
                            rsync web-02 · kontrola kodów
                                       │  przydział zadania
                                       ▼
                            Wykonawca podejmuje zadanie ────►  polecenie w karcie,
                                       ▲                        strumień wyniku
                                       │  raport postępu               │
                                       │◄──────────────────────────────┘
                                       ▼
                            kontrola jakości: kod wyjścia,
                            wzorzec błędu, porównanie z
                            poprzednim przebiegiem
                                       │
                        ┌──────────────┴──────────────┐
                        ▼                             ▼
                 wynik poprawny               wynik błędny —
                 → kolejne zadanie            ponowienie albo
                                              korekta polecenia
                        │                             │
                        └──────────────┬──────────────┘
                                       ▼
                       zakończenie przebiegu · dziennik pętli ·
                       podsumowanie przekazane do Chat Window
```

Użytkownik steruje przebiegiem z obu kolumn: z Chat Window przerywa i koryguje zlecenie w języku naturalnym, z Execution Loop Window wstrzymuje, wznawia i przerywa pętlę oraz decyduje o ponowieniu pojedynczego zadania.

### 4.4. Przepływ rozszerzony — awaria procesu i przekazanie do Diagnostics

```
Output Console                                    Diagnostics › Errors Panel
──────────────                                    ──────────────────────────
 wykrycie wzorca błędu w wyniku
 polecenia — nieobsłużony wyjątek
        │
        ▼
 pasek akcji linii błędu
 [ Wyjaśnij ]  [ Przekaż do Diagnostics ]
        │                    │
        ▼                    ▼
 Chat Window            nowy wpis błędu z pełnym
 (wyjaśnienie            kontekstem: polecenie, karta,
  kontekstowe)            katalog roboczy, znacznik czasu  ──────►
                                                                     agregacja w
                                                                     Diagnostics Center
                                                                     i dalsza analiza
```

### 4.5. Przepływ zbiorczy — równoległa praca wielu kart i powłok

```
Karta Terminal Tabs: PowerShell    Karta Terminal Tabs: Bash        Karta Terminal Tabs: Python
  build serwera Windows              wdrożenie na hoście Linux         uruchomienie zestawu testów
        │                                    │                                 │
        └────────────────┬───────────────────┴────────────────┬────────────────┘
                         ▼                                    ▼
                 Output Console — widok zagregowany     Process Monitor — wszystkie
                 wszystkich trzech kart jednocześnie,    procesy trzech kart w jednym
                 z filtrowaniem po źródle                rejestrze, grupowanie „po karcie”

           stan wyjściowy: katalog roboczy i środowisko procesu współdzielone w karcie sesji
           po włączeniu izolacji technicznej z okna konfiguracji punktów izolacji (rozdz. 5.3):
           odrębny katalog roboczy per karta
```

### 4.6. Przepływ zadania cyklicznego i wyzwalacza

```
Task & Schedule                       Terminal Tabs / Output Console
───────────────                       ──────────────────────────────
 wyrażenie cron albo wzorzec watch
        │  nadejście terminu lub zmiana pliku
        ▼
 uruchomienie zadania w karcie ───────►  polecenie wykonane, strumień
 wskazanej w konfiguracji                wyniku w Output Console
        │                                          │
        ▼                                          ▼
 wpis w historii przebiegów            wiersz procesu w Process Monitor
 (status, czas, kod wyjścia)           z inicjatorem i odnośnikiem zadania
        │
        ▼
 przekroczenie progu czasu lub błąd → powiadomienie kanałem ustawionym
 w oknie konfiguracji; przy zleceniu wielokrokowym przebieg prowadzi
 Koordynator w Execution Loop Window
```

---

## 5. Stany, dane i powiązania

### 5.1. Model stanów procesu

Moduł operuje na dwóch poziomach procesu, celowo rozróżnionych, aby uniknąć pomylenia procesu sesji platformy z pojedynczym poleceniem powłoki.

| Poziom | Czym jest | Gdzie widoczny |
|---|---|---|
| Proces sesji (`proces_sesji`) | Odrębny, izolowany proces serwera obsługujący całą kartę sesji modułu Terminal — własny katalog roboczy i środowisko | Trwa niezależnie od stanu połączenia klienta; niewidoczny wprost, stanowi „podłoże” dla kart powłok |
| Proces uruchomieniowy (wpis rejestru procesów) | Pojedyncze polecenie wykonane w jednej z kart Terminal Tabs, kluczowane identyfikatorem uruchomienia | Prezentowany wprost w Process Monitor, z wynikiem w Output Console |

```
                    ┌──────────────────┐
                    │  Karta utworzona │  (nowa sesja powłoki)
                    └────────┬─────────┘
                             │ polecenie — od Operatora, Wykonawcy
                             │ albo z zadania pętli wykonawczej
                             ▼
                    ┌──────────────────┐
        ┌──────────►│  Uruchomiony     │  wpis w rejestrze procesów,
        │           └────────┬─────────┘  strumień wyniku na żywo
        │                    │
        │        ┌───────────┼─────────────┐
        │        ▼           ▼             ▼
        │  ┌──────────┐┌──────────┐┌────────────────┐
        │  │Zakończony││Zakończony││   Zakończony   │
        │  │ sukcesem ││  błędem  ││ przez Operatora│
        │  └────┬─────┘└────┬─────┘└───────┬────────┘
        │       │           │              │  wykrycie wzorca błędu
        │       │           └──────────────┼───────────► Diagnostics
        │       │                          │             (Errors Panel)
        └───────┴──────────────────────────┘
          „Uruchom ponownie” z Process Monitor albo decyzja
          o ponowieniu podjęta w Execution Loop Window
```

Stan procesu sesji po stronie serwera jest trwały niezależnie od stanu połączenia klienta — rozłączenie nie przerywa działających powłok ani uruchomionych poleceń; po ponownym połączeniu Terminal Tabs przywraca karty w stanie, w jakim zostały pozostawione, Process Monitor odzyskuje pełny obraz rejestru procesów, a Execution Loop Window wznawia prezentację przebiegu pętli.

### 5.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Terminal |
|---|---|
| `karta_sesji`, `sesja` | Nośnik kontekstu bieżącej pracy w powłoce; jedna sesja modułu na kartę |
| `proces_sesji` | Proces serwera obsługujący całą kartę sesji Terminal — pole `katalog_roboczy` odpowiada katalogowi startowemu kart Terminal Tabs, pole `stan` (`uruchomiony`/`wstrzymany`/`zakonczony`) opisuje żywotność sesji niezależnie od pojedynczych poleceń |
| `wiadomosc`, `zalacznik_wiadomosci` | Historia Chat Window, w tym polecenia zaproponowane i uruchomione przez Wykonawcę |
| `zlecenie`, `zadanie_petli`, `komunikat_sterujacy` | Treść pętli wykonawczej prezentowanej w Execution Loop Window: dekompozycja zlecenia, kolejka i stan zadań, wymiana komunikatów Koordynator ↔ Wykonawca, wyniki kontroli jakości i decyzje o ponowieniu |
| `kanal_modelu` | Model bazowy przypisany karcie sesji, obsługujący komunikację i generowanie poleceń |
| `ustawienie` | Parametry modułu podlegające zasadzie „brak ustawienia = wartość domyślna” — limit buforowanych linii Output Console, próg powiadomienia o długo trwającym procesie, zasięg skanowania manifestów zadań |
| `host_polaczenia`, `sesja_trwala`, `tunel_portowy`, `klucz_ssh` | Wpisy Session Managera: książka hostów, sesje żyjące po stronie serwera, przekierowania portów i materiał kluczowy |
| `skrypt`, `snippet`, `parametr_skryptu` | Zawartość Script Library wraz z deklaracjami parametrów renderowanymi jako formularz |
| `zadanie_harmonogramu`, `wyzwalacz`, `przebieg_zadania` | Zawartość okna Task & Schedule: wyrażenia cron, wzorce obserwacji plików, zdarzenia oraz historia przebiegów |
| `sekret_sesji`, `zmienna_srodowiskowa` | Magazyn sekretów i warstwowe zmienne środowiskowe wstrzykiwane do poleceń bez zapisu w historii i transkrypcie |
| `profil_izolacji`, `regula_izolacji_technicznej` | Konfiguracja ośmiu zakresów izolacji technicznej procesu (rozdz. 5.3) — zakres najściślej właściwy modułowi Terminal ze wszystkich piętnastu modułów platformy |
| `powiazanie_komponentu` | Jawne powiązania Terminal ◄──► Developer, Terminal ◄──► Diagnostics, Terminal ◄──► Automations |

Rejestr procesów uruchomieniowych prezentowany w Process Monitor jest strukturą operacyjną rdzenia serwera (odrębną od trwałych encji SQLite) — aktualizowaną i przekazywaną na żywo kanałem WebSocket, zgodnie z zasadą synchronizacji na żywo właściwą oknom monitorującym procesy. Tym samym kanałem przekazywany jest stan kolejki zadań pętli wykonawczej.

### 5.3. Izolacja i konfigurowalność — punkty właściwe modułowi Terminal

Terminal jest modułem, w którym wszystkie osiem zakresów izolacji technicznej ma bezpośrednie, dosłowne zastosowanie — moduł jest samą warstwą wykonania procesów. Stan wyjściowy każdego zakresu jest wyłączony: karty powłok współdzielą katalog roboczy, środowisko procesu i dostęp sieciowy karty sesji, a proces powłoki nie jest ograniczany technicznie ponad to, co wynika z uprawnień Operatora na serwerze. Każde zawężenie jest jawną, odwracalną decyzją podjętą z okna konfiguracji punktów izolacji.

| Zakres izolacji technicznej | Stan wyjściowy | Zastosowanie w module Terminal | Typowy powód włączenia |
|---|---|---|---|
| Katalog roboczy sesji | Wyłączony (współdzielony) | Katalog startowy wszystkich kart powłok karty sesji | Odseparowanie plików roboczych karty diagnostycznej od karty produkcyjnej |
| Środowisko procesu | Wyłączony (współdzielony) | Zmienne środowiskowe dostępne poleceniom we wszystkich kartach | Praca z odrębnym zestawem zmiennych dla środowiska testowego |
| Katalog danych i konfiguracji modelu | Wyłączony (współdzielony) | Dane pomocnicze kanału modelu obsługującego Chat Window i pętlę wykonawczą | Odrębna konfiguracja modelu dla sesji o podwyższonej poufności |
| Dostęp sieciowy | Wyłączony (współdzielony) | Połączenia wychodzące poleceń uruchamianych w kartach, jak `curl` i `npm install` | Ograniczenie sesji wykonującej skrypty z materiału niezaufanego |
| Zakres odczytu i zapisu plików | Wyłączony (pełny dostęp) | Uprawnienia poleceń do plików poza katalogiem roboczym karty | Sesja pracująca wyłącznie na wskazanym podkatalogu repozytorium |
| Konto i token per sesja | Wyłączony (współdzielone dane dostępowe) | Dane logowania wykorzystywane przez polecenia CLI wymagające uwierzytelnienia oraz profile poświadczeń `aws`/`gcloud`/`az`/`kubectl` | Karta z osobnym tokenem, aby nie wyczerpywać wspólnego limitu |
| Model procesu | Wyłączony (współdzielony) | Instancja procesu wykonawczego obsługującego Chat Window i Execution Loop Window | Rola Executor w środowisku MultitaskingAI korzystająca z Terminala na dedykowanym procesie |
| Serwer wykonania | Wyłączony (współdzielony) | Fizyczny lub logiczny serwer wykonujący powłoki karty | Sesja SSH pracująca na dedykowanym hoście zdalnym |

Zgodnie z zasadą nadrzędną platformy żaden z powyższych zakresów nie jest wymuszony — moduł jest w pełni funkcjonalny bez jakiejkolwiek konfiguracji izolacji, a każdy zakres jest odwracalną, jawną decyzją Operatora podejmowaną z panelu macierzy izolacji ([Koncepcja platformy](../architektura/koncepcja-platformy.md#64-okno-konfiguracji-punktów-izolacji), rozdz. 6.4).

### 5.4. Powiązania z innymi modułami

```
                     ┌────────────────────────────────┐
   Developer  ◄──────┤            TERMINAL            ├───────►  Diagnostics
   (warstwa          │ warstwa wykonawcza poleceń dla │          (odtwarzanie
    wykonawcza       │ zapytań Build Output i poleceń │          i weryfikacja objawów
    poleceń Code     │wydawanych z modułu Diagnostics)│          błędu w rzeczywistym
    Editor)          └──────────────────┬─────────────┘          środowisku wykonawczym)
                                        │
                                        ▼
                                   Automations
                        (wyzwalacze zdarzeniowe zadań powłoki)
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Developer | Terminal dostarcza warstwy wykonawczej dla poleceń wydawanych z poziomu Developer (budowanie, uruchamianie, instalacja zależności); Script Library korzysta z repozytorium Git prowadzonego w module Developer | Konfiguracyjne | Build Output / Git Panel → Terminal Tabs; Script Library → Git Panel |
| Diagnostics | Terminal dostarcza warstwy wykonawczej do odtwarzania i weryfikacji objawów błędu w rzeczywistym środowisku wykonawczym | Konfiguracyjne | Diagnostics Center → Terminal Tabs; Output Console → Errors Panel |
| Automations | Wyzwalacze zdarzeniowe zadań powłoki (webhook, zdarzenie zewnętrzne, zakończenie innego zadania) realizowane jawnym powiązaniem, nie ukrytą regułą | Konfiguracyjne | Automations → Task & Schedule |

Współobecność Terminala, Developera i Diagnostics w jednym środowisku CodeStudio czyni powiązanie typowe w praktyce, lecz — zgodnie z zasadą jawności i konfigurowalności zależności — nie jest ono wbudowaną na stałe zależnością techniczną; ustanawia je Operator z okna konfiguracji. Zasada zero blokad obowiązuje w całym zakresie modułu: potwierdzenia sygnalizują wagę akcji (zakończenie procesu, zamknięcie karty z działającym procesem, przerwanie pętli wykonawczej), lecz nigdy nie blokują jej wykonania.

---

## 6. Scenariusze użycia

**Scenariusz 1 — uruchomienie budowania z obserwacją wyniku.**
Deweloper otwiera kartę PowerShell w Terminal Tabs, przechodzi do katalogu repozytorium przez breadcrumb i wpisuje polecenie budowania. Output Console gromadzi wynik na żywo; Process Monitor pokazuje proces jako aktywny z rosnącym czasem trwania. Po zakończeniu status zmienia się na „zakończony sukcesem”, a deweloper przechodzi do modułu Developer, aby zobaczyć wynik w Build Output.

**Scenariusz 2 — polecenie wydane przez Wykonawcę w toku pracy nienadzorowanej.**
W ramach roli Executor środowiska MultitaskingAI Wykonawca korzystający z modułu Terminal uruchamia zestaw testów po każdej zmianie kodu. Każde uruchomienie pojawia się w Process Monitor z plakietką „inicjator: Wykonawca”, a Użytkownik — poprzez funkcję globalną Mobile lub wracając do sesji — zachowuje pełny wgląd w to, co dokładnie zostało wykonane, oraz możliwość natychmiastowego przerwania działania z Chat Window.

**Scenariusz 3 — zlecenie wielozadaniowe prowadzone przez pętlę wykonawczą.**
Użytkownik zleca w Chat Window zbudowanie aplikacji i wdrożenie jej na dwóch hostach. Koordynator rozkłada zlecenie na pięć zadań powłokowych i prezentuje kolejkę w Execution Loop Window. Wykonawca podejmuje zadania kolejno, raportując postęp; kontrola jakości po zadaniu trzecim wykrywa niezerowy kod wyjścia, Koordynator ponawia zadanie z poprawionym parametrem ścieżki, a Użytkownik obserwuje przebieg i wskaźnik ponowień, zachowując możliwość wstrzymania pętli.

**Scenariusz 4 — diagnoza błędu i przekazanie do Diagnostics.**
Skrypt Pythona kończy się nieobsłużonym wyjątkiem widocznym w Output Console. Operator klika „Przekaż do Diagnostics” przy oznaczonej linii błędu — nowy wpis pojawia się w Errors Panel modułu Diagnostics z pełnym kontekstem: treścią błędu, katalogiem roboczym i znacznikiem czasu, gotowy do dalszej analizy bez ręcznego kopiowania.

**Scenariusz 5 — praca równoległa na wielu hostach.**
Inżynier DevOps otwiera z Session Managera cztery karty: PowerShell lokalnie, Bash lokalnie oraz dwie sesje SSH do odrębnych hostów wdrożeniowych. Dla sesji SSH włącza w oknie konfiguracji punktów izolacji zakres „serwer wykonania” i „konto i token per sesja”, aby żadna z sesji nie współdzieliła danych dostępowych z pozostałymi. Pasek rozgłaszania kieruje to samo polecenie do obu hostów, a Output Console w widoku zagregowanym pozwala obserwować wdrożenie jednocześnie.

**Scenariusz 6 — zawieszony proces wymagający interwencji.**
Proces uruchomiony przez wcześniejsze polecenie przestaje odpowiadać. Operator odnajduje go w Process Monitor po długim czasie trwania, rozwija szczegóły procesu, aby potwierdzić polecenie i katalog roboczy, po czym wybiera „Zakończ”. Okno nakładkowe potwierdzenia sygnalizuje wagę decyzji ze względu na niezapisany stan procesu; po potwierdzeniu proces kończy się natychmiast, a karta powłoki wraca do stanu gotowości.

**Scenariusz 7 — zadanie cykliczne i skrypt z parametrami.**
Administrator zapisuje w Script Library skrypt kopii zapasowej z zadeklarowanymi parametrami środowiska i wersji, uruchamia go raz z formularza dla weryfikacji, po czym w oknie Task & Schedule tworzy zadanie cykliczne z wyrażeniem cron. Sekrety dostępu wstrzykiwane są z magazynu sekretów sesji i maskowane w Output Console oraz w nagraniu sesji. Historia przebiegów gromadzi status każdego uruchomienia, a przekroczenie progu czasu wywołuje powiadomienie.

---

## 7. Katalog funkcji i narzędzi

Każda pozycja obejmuje nazwę funkcji, opis działania oraz zależności techniczne (biblioteki, formaty, integracje). Zależności warstwy serwerowej realizuje rdzeń serwera w Go; warstwa terminala w kliencie oparta jest na `xterm.js`.

### 7.1. Powłoki, sesje i emulator terminala

| Funkcja | Co robi | Zależności |
|---|---|---|
| Wielopowłokowy emulator | Pełna emulacja terminala (VT100/xterm, sekwencje ANSI, 24‑bitowy kolor, tryb aplikacyjny myszy) w karcie | Klient: `xterm.js` + `xterm-addon-fit`, `xterm-addon-web-links`, `xterm-addon-webgl` (render), `xterm-addon-unicode11`. Serwer: `creack/pty` (Unix) / ConPTY przez `UiPath/pty` lub `photostorm/pty` (Windows) |
| Menedżer profili powłok | Nazwane profile: typ powłoki, katalog startowy, zmienne, polecenie startowe, motyw, ikona | Format profilu: JSON w `ustawienie`; integracja z warstwą konfiguracji sesja/projekt |
| Multiplekser paneli i układów | Podział karty na siatkę sąsiadujących paneli z zagnieżdżaniem, zapisywane układy, powiększenie panelu | Logika układu w kliencie; stan paneli w `migawka` karty sesji |
| Trwałe sesje nazwane (detach/attach) | Sesja żyje po stronie serwera po rozłączeniu klienta; ponowne podpięcie przywraca bufor i stan | `proces_sesji` (rdzeń serwera); bufor ekranu w pamięci rdzenia; kanał WebSocket do wznowienia |
| Rozgłaszanie wpisu do wielu paneli | Jedno wpisywane polecenie trafia równolegle do zaznaczonych paneli i hostów | Multipleksacja strumieni wejścia w rdzeniu |
| Paleta poleceń | `Ctrl/Cmd + K`: szybkie akcje, wyszukiwanie profili, snippetów, hostów, zadań | Indeks w kliencie; źródła: profile, biblioteka skryptów, książka hostów |
| Bloki poleceń | Każde polecenie i jego wynik jako osobny, składany blok z akcjami (kopiuj, ponów, udostępnij, oznacz) | Parser granic poleceń (OSC 133 shell integration); render bloków w `xterm.js` |
| Integracja powłoki (shell hooks) | Znaczniki OSC 7 (katalog roboczy) i OSC 133 (granice poleceń, kod wyjścia) przekazywane automatycznie z powłoki | Skrypty integracyjne dla PowerShell/Bash/Zsh; sekwencje OSC |
| Wyszukiwanie w buforze | Grep w obrębie karty z wyrażeniami regularnymi, podświetleniem i licznikiem | `xterm-addon-search` |
| Autouzupełnianie inline i podpowiedzi | Uzupełnianie ścieżek, poleceń i flag; podpowiedź z historii | Provider uzupełnień w rdzeniu; słowniki flag CLI; historia per karta |

### 7.2. Zdalne wykonanie i połączenia

| Funkcja | Co robi | Zależności |
|---|---|---|
| Klient SSH | Sesja SSH jako karta powłoki; klucze, hasło, agent, `keyboard-interactive` | `golang.org/x/crypto/ssh`; klucze OpenSSH/PEM; agent `SSH_AUTH_SOCK` |
| Książka hostów | Katalog hostów z folderami, tagami, danymi połączenia, jump‑hostami; szybkie łączenie | Wpisy w konfiguracji + import/eksport `~/.ssh/config` |
| ProxyJump / bastion | Łańcuch przeskoków przez host pośredni | `x/crypto/ssh` z zagnieżdżonym `Dial` |
| Przekierowanie portów | Lokalne, zdalne i dynamiczne (SOCKS) tunele z panelem stanu | `x/crypto/ssh` kanały; SOCKS5 w Go |
| SFTP / transfer plików | Panel plików zdalnych, wysyłka i pobieranie, wznawianie, edycja w miejscu | `pkg/sftp` (Go); format SFTP v3 |
| Rsync / synchronizacja | Przyrostowa synchronizacja katalogów lokalny↔zdalny z podglądem różnic | Wywołanie `rsync` lub natywna implementacja delta w Go |
| Sesje na wielu hostach | Jedna operacja na grupie hostów, agregacja wyników per host | Pula połączeń `x/crypto/ssh`; agregacja w Output Console |
| `docker exec` / attach | Karta powłoki wewnątrz kontenera; lista kontenerów | Docker Engine API (`/containers/{id}/exec`) przez socket |
| `kubectl exec` / logs | Powłoka w podzie, strumień logów, wybór kontekstu i przestrzeni nazw | `client-go` (Kubernetes); kontekst z `kubeconfig` |
| Konsola szeregowa | Połączenie RS‑232/USB‑serial (prędkość, parzystość, bity) jako karta | `go.bug.st/serial`; wybór portu COM/tty |
| Telnet | Sesja Telnet do urządzeń sieciowych i starszych systemów | Natywna implementacja Telnet (negocjacja opcji) w Go |
| Odporność łącza | Utrzymanie sesji przy zmianie adresu IP i zaniku łącza, echo lokalne | UDP + przewidywanie lokalne oparte na trwałej sesji rdzenia |

### 7.3. Skrypty, snippety i biblioteka

| Funkcja | Co robi | Zależności |
|---|---|---|
| Edytor skryptów | Wieloliniowy edytor (PowerShell/Bash/Python) z podświetlaniem i uruchomieniem | Monaco Editor w kliencie; wykonanie przez `proces_sesji` |
| Biblioteka skryptów | Repozytorium skryptów z tagami, opisem, wersją, czasem ostatniego uruchomienia | Encja skryptu w konfiguracji i projekcie; format: plik + metadane YAML front‑matter |
| Snippety poleceń | Krótkie fragmenty przypisane do skrótu i aliasu, wstawiane z palety | Indeks snippetów; format YAML; podstawianie zmiennych |
| Parametryzacja skryptów | Skrypt z deklarowanymi parametrami renderuje formularz przed uruchomieniem | Parser deklaracji parametrów; format: front‑matter `args:` |
| Generowanie skryptu z języka naturalnego | Chat Window buduje skrypt wieloliniowy gotowy do zapisu w bibliotece | Kanał modelu; integracja z biblioteką skryptów |
| Linter i formatowanie skryptów | Statyczna analiza i formatowanie przed uruchomieniem | Integracja `shellcheck`, `shfmt`, PSScriptAnalyzer, `ruff`/`black` jako narzędzia CLI |
| Aliasy i funkcje powłoki | Zarządzanie aliasami i funkcjami z interfejsu, synchronizowane między sesjami | Zapis do profilu powłoki; warstwa konfiguracji |
| Pliki konfiguracyjne powłoki i synchronizacja profilu | Wersjonowany zestaw plików konfiguracyjnych powłoki, przenośny między hostami | Repozytorium plików konfiguracyjnych powłoki; integracja z Git (moduł Developer) |

### 7.4. Zadania, harmonogram i automatyzacja poleceniowa

| Funkcja | Co robi | Zależności |
|---|---|---|
| Runner zadań | Wykrywa i uruchamia zadania z `package.json`, `Makefile`, `Taskfile.yml`, `justfile` | Parsery: `Makefile`, JSON (`package.json`), YAML (`Taskfile`), `justfile`; wykonanie w rdzeniu |
| Harmonogram zadań | Cykliczne uruchomienia po wyrażeniu cron lub o wskazanym czasie; kalendarz nadchodzących uruchomień | `robfig/cron` (Go); format crontab; encja zadania harmonogramu |
| Wyzwalacze plikowe | Uruchomienie polecenia przy zmianie plików i katalogów (tłumienie powtórzeń, wzorce) | `fsnotify` (Go); wzorce glob |
| Potoki i sekwencje poleceń | Łańcuch kroków z warunkami (sukces/porażka), zmiennymi i równoległością | DSL potoku (YAML); silnik wykonania w rdzeniu; graf zależności kroków |
| Kolejka zadań i przebiegi | Historia przebiegów zadania: status, czas, wynik, ponowienie; kolejkowanie | Encja przebiegu; rejestr w SQLite; strumień na żywo |
| Powiadomienia o zakończeniu | Alert (globalny i Mobile) po zakończeniu długiego zadania oraz przy błędzie | Integracja z funkcją globalną Mobile i powiadomieniami platformy; próg czasu w konfiguracji |
| Wyzwalacze zdarzeniowe | Uruchomienie polecenia na zdarzenie: webhook, wpis w Diagnostics, koniec innego zadania | Jawne, konfigurowalne powiązanie z modułem Automations |
| Bramy wykonania równoległego | Uruchomienie N zadań równolegle z limitem współbieżności i agregacją | Pula workerów w rdzeniu; limit z konfiguracji |

### 7.5. Sekrety, środowisko i uwierzytelnienie CLI

| Funkcja | Co robi | Zależności |
|---|---|---|
| Menedżer zmiennych środowiskowych | Edytor zmiennych per sesja, projekt i host, warstwowe nadpisania, podgląd efektywnego zestawu zmiennych | Warstwowość konfiguracji (globalna→środowisko→projekt→sesja); format `.env` |
| Ładowarka plików `.env` | Automatyczne wczytanie `.env` i `.env.local` po wejściu do katalogu | Parser `.env`; hak katalogu (OSC 7) |
| Magazyn sekretów sesji | Sekrety wstrzykiwane do polecenia bez zapisu w historii i transkrypcie; maskowanie w wyniku | Szyfrowany magazyn (klucz jawnie zarządzany zgodnie z zasadą platformy); integracja `regula_izolacji_technicznej` „konto i token per sesja” |
| Maskowanie wrażliwych wartości | Automatyczne zaciemnianie tokenów i kluczy w wyniku według wzorców | Reguły wzorców (wyrażenia regularne); filtr strumienia Output Console |
| Menedżer kluczy SSH | Generowanie, import, agent, przypisanie klucza do hosta | `x/crypto/ssh`; formaty OpenSSH/PEM/PPK (import) |
| Profile poświadczeń CLI | Przełączanie profili `aws`/`gcloud`/`az`/`kubectl` per karta | Wstrzykiwanie zmiennych i kontekstu; izolacja „konto i token per sesja” |

### 7.6. Wynik, obserwowalność i rejestracja

| Funkcja | Co robi | Zależności |
|---|---|---|
| Rejestracja sesji | Nagranie sesji jako odtwarzalny plik z osią czasu | Format `asciicast` v2 (JSON‑lines); rejestrator w rdzeniu |
| Odtwarzanie | Przewijanie wykonania krok po kroku lub w przyspieszeniu w Output Console | Player oparty na `asciicast`; render w kliencie |
| Eksport transkryptu i dziennika | Zapis karty lub zakresu do `.txt`, `.log`, `.md`, `.html` (z kolorami) | Konwersja ANSI→HTML; szablony eksportu |
| Parser strukturalny wyniku | Rozpoznanie JSON i tabel w wyniku oraz render jako tabela lub drzewo | Silnik zgodny z `jq` (`itchyny/gojq`); wykrywanie formatu |
| Wykresy zasobów procesu | Sparkline CPU/RAM/IO per proces i zbiorczo, na żywo | Metryki z rdzenia (`shirou/gopsutil`); strumień WebSocket |
| Rejestr i drzewo procesów | Pełny rejestr uruchomień z hierarchią potomków, filtrami, akcjami (zakończenie, wstrzymanie, ponowienie) | `gopsutil`; sygnały POSIX / Job Objects (Windows) |
| Wykrywanie wzorców błędów | Oznaczenie linii błędu, kodu wyjścia i wyjątku plakietką oraz akcją | Reguły wzorców; przekazanie do Diagnostics (Errors Panel) |
| Różnicowanie wyników | Porównanie wyjścia dwóch uruchomień obok siebie | Silnik diff w Go (`sergi/go-diff`); render kolorowy |
| Dziennik pętli wykonawczej | Zapis zlecenia, zadań, komunikatów sterujących i wyników kontroli jakości | Encje `zlecenie`, `zadanie_petli`, `komunikat_sterujacy`; eksport do `.md` i `.json` |

### 7.7. Narzędzia wiersza poleceń wbudowane w moduł

| Funkcja | Co robi | Zależności |
|---|---|---|
| Klient HTTP w karcie | Wykonanie i podgląd żądań HTTP z formularzem i historią | Klient HTTP Go (`net/http`); kolekcje w formacie HAR i `.http` |
| Podgląd i edycja plików strukturalnych | Podgląd i edycja JSON/YAML/TOML/CSV z walidacją | `gojq`, `goccy/go-yaml`, `BurntSushi/toml`; parser CSV |
| Menedżer pakietów — pulpit | Jednolity widok `npm`/`pip`/`cargo`/`go`/`winget`/`brew`: instalacja, aktualizacja, audyt | Adaptery per menedżer; parsowanie manifestów (`package.json`, `requirements.txt`, `Cargo.toml`, `go.mod`) |
| Konwertery i kodowanie | base64/hex, kodowanie URL, dekodowanie JWT, UUID, skróty (md5/sha) | Biblioteki standardowe Go (`crypto`, `encoding`) |
| Kalkulator poleceń i zmienne | Zmienne robocze i wyrażenia dostępne w powłoce, w tym wynik poprzedniego polecenia | Podstawienia w warstwie rdzenia; „ostatni wynik” jako zmienna |
| Menedżer plików dwupanelowy | Lokalny↔lokalny i lokalny↔zdalny: kopiowanie, przenoszenie, usuwanie, podgląd | Warstwa systemu plików w rdzeniu + SFTP dla zasobów zdalnych |

---

## 8. Punkty sterowania z okna konfiguracji

Wszystkie punkty są jawne, warstwowe (globalny → środowisko → projekt → sesja) i opatrzone objaśnieniem kontekstowym; stan wyjściowy nie blokuje pracy.

| Grupa | Punkt konfiguracji | Co personalizuje Operator |
|---|---|---|
| Emulator | Motyw i schemat kolorów | Paleta ANSI, tło, kursor; zgodność z motywem platformy lub ustawienie niezależne |
| Emulator | Czcionka i gęstość | Krój mono, rozmiar, ligatury, interlinia, renderer (canvas/webgl) |
| Emulator | Integracja powłoki | Ustawienie OSC 7/133, bloki poleceń, autouzupełnianie inline |
| Sesje | Powłoka domyślna i profile | Domyślny profil per środowisko i projekt; edycja i kolejność profili |
| Sesje | Trwałość sesji | Przeżywanie rozłączenia przez sesje; limit czasu bezczynności; automatyczne podpięcie |
| Sesje | Bufor i historia | Limit buforowanych linii; zakres i współdzielenie historii poleceń |
| Okna komunikacji | Chat Window | Kanał modelu, zakres pamięci, szerokość lewej kolumny, zakres kontekstu przekazywanego z aktywnej karty |
| Okna komunikacji | Execution Loop Window | Otwieranie kolumny pętli, próg automatycznego ponowienia zadania, limit współbieżności zadań pętli, zakres dziennika pętli |
| Zdalne | Książka hostów | Wpisy, foldery, tagi, jump‑hosty; import `~/.ssh/config` |
| Zdalne | Klucze i agent | Ścieżki kluczy, użycie agenta, polityka `known_hosts` |
| Zdalne | Tunele domyślne | Predefiniowane przekierowania portów uruchamiane z sesją |
| Zadania | Źródła zadań | Zakres skanowanych manifestów (`package.json`/`Makefile`/`Taskfile`/`just`) |
| Zadania | Harmonogram i wyzwalacze | Praca schedulera, strefa czasowa, wzorce obserwacji plików, limit współbieżności |
| Zadania | Powiadomienia | Próg czasu długiego zadania; kanały (Mobile, globalne); powiadomienie o błędzie |
| Sekrety | Menedżer zmiennych i `.env` | Automatyczne ładowanie `.env`; kolejność warstw nadpisań; widoczność efektywnego zestawu zmiennych |
| Sekrety | Maskowanie | Wzorce zaciemniania w wyniku; wykluczenie z transkryptu i rejestracji |
| Sekrety | Konto i token per sesja | Przypisanie osobnych poświadczeń do karty (zakres izolacji technicznej) |
| Obserwowalność | Rejestracja sesji | Automatyczne nagrywanie wybranych profili; format i retencja `asciicast` |
| Obserwowalność | Metryki procesów | Częstotliwość próbkowania CPU/RAM/IO; progi alertów |
| Izolacja | Osiem zakresów izolacji technicznej | Katalog roboczy, środowisko procesu, dane modelu, sieć, odczyt i zapis plików, konto i token, model procesu, serwer wykonania (per karta i sesja) |
| Powiązania | Jawne linki modułowe | Terminal↔Developer, Terminal↔Diagnostics, Terminal↔Automations (ustanawiane, nie wbudowane) |
| Bezpieczeństwo działania | Potwierdzenia wagi | Zakres akcji (zakończenie procesu z niezapisanym stanem, zamknięcie karty z procesem, przerwanie pętli) pokazujących okno nakładkowe sygnalizujące wagę — sygnał, nie blokada |

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + T` | Nowa karta powłoki | Terminal Tabs |
| `Ctrl/Cmd + W` | Zamknięcie bieżącej karty | Terminal Tabs |
| `Ctrl/Cmd + Tab` | Przełączenie na kolejną kartę | Terminal Tabs |
| `Ctrl/Cmd + Shift + \` | Podział karty na sąsiadujący panel | Terminal Tabs |
| `Ctrl/Cmd + K` | Paleta poleceń: profile, snippety, hosty, zadania | Terminal Tabs, Script Library |
| `↑` / `↓` w wierszu poleceń | Nawigacja po historii poleceń karty | Terminal Tabs |
| `Ctrl/Cmd + F` | Wyszukiwanie w bieżącej karcie | Terminal Tabs |
| `Ctrl/Cmd + Shift + F` | Otwarcie Grep w Output Console | Output Console |
| `Ctrl/Cmd + Enter` | Wysłanie polecenia do Wykonawcy | Chat Window |
| `Ctrl/Cmd + Shift + L` | Otwarcie kolumny pętli wykonawczej | Chat Window, Execution Loop Window |
| `Ctrl/Cmd + .` | Wstrzymanie i wznowienie przebiegu pętli | Execution Loop Window |
| `/` | Wywołanie menu operacji z poziomu okna komunikacji | Chat Window |
| Ikona `uruchom` | Start procesu / uruchomienie polecenia | Process Monitor, Chat Window, Script Library |
| Ikona `zatrzymaj` | Zatrzymanie procesu, przerwanie przebiegu pętli | Process Monitor, Execution Loop Window |
| Ikona `odswiez` | Ponowne uruchomienie polecenia lub zadania | Process Monitor, Output Console, Execution Loop Window |
| Ikona `folder` | Katalog roboczy sesji (izolacja techniczna) | Terminal Tabs |
| Ikona `globus` | Dostęp sieciowy (izolacja techniczna) | okno konfiguracji punktów izolacji |
| Ikona `klodka` | Konto i token per sesja, magazyn sekretów | okno konfiguracji punktów izolacji, Session Manager |
| Ikona `blad` | Oznaczenie linii wyniku jako błędu | Output Console |
| Ikona `oko` | Podgląd szczegółów procesu | Process Monitor |
| Ikona `zegar` | Zadanie cykliczne i wyzwalacz | Task & Schedule |
| Ikona `serwer` | Host i połączenie zdalne | Session Manager, Terminal Tabs |
| Ikona `kod` | Blok treści technicznej (polecenie, ścieżka, skrypt) | Chat Window, Output Console, Process Monitor, Script Library |

---

## 10. Punkty łamania i kryteria odbioru

### 10.1. Punkty łamania

Moduł Terminal dzieli obszar roboczy na kolumny sąsiadujące poziomo (rozdz. 2). Poniższe progi, wspólne całej platformie, rozstrzygają zachowanie układu tych kolumn wraz ze zmianą szerokości okna aplikacji.

| Żeton | Próg szerokości | Zachowanie układu w module Terminal |
|---|---|---|
| `--dn-bp-w1` | 640 px | Telefon poziomo — widok mobilny; okna operacyjne modułu prezentowane pojedynczo, pełny ekran, nawigacja powrotna zastępuje układ kolumnowy |
| `--dn-bp-w2` | 960 px | Tablet — boczna nawigacja modułów zwija się do samych ikon; kolumny modułu Terminal zachowują układ, lecz kolumny boczne (rozszerzenia warstwy 2–3) otwierają się jako nakładka zamiast stałej kolumny |
| `--dn-bp-w3` | 1280 px | Biurko — pełny kokpit; wszystkie okna operacyjne modułu Terminal wymienione w rozdz. 2 dostępne jednocześnie w układzie kolumnowym opisanym w tym rozdziale |
| `--dn-bp-w4` | 1600 px | Szerokie biurko — para Chat Window i Execution Loop Window prezentowana jednocześnie obok okna wiodącego modułu, bez wzajemnego przesłaniania |

Poniżej progu `--dn-bp-w1` układ kolumnowy modułu Terminal nie jest dostępny — zachowanie na urządzeniach mobilnych ustala [Mobile](../funkcje-globalne/mobile.md).

### 10.2. Kryteria odbioru

Warunki sprawdzalne, których łączne spełnienie oznacza gotowość modułu Terminal do odbioru.

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Wszystkie 8 okien operacyjnych z rozdz. 2 otwiera się dokładnie sposobem wywołania opisanym w tabeli okien i w katalogu elementów interfejsu rozdz. 3 | Przegląd manualny wg tabeli rozdz. 2 — każde okno wywołane, sprawdzone wejście i zamknięcie |
| Każda z 27 komend obszaru `terminal` z Załącznika ma pokrycie w co najmniej jednym przepływie pracy albo scenariuszu użycia (rozdz. 4, 6) | Zestawienie nazw komend z treścią rozdz. 4 i 6 |
| Process Monitor rozróżnia w każdym wpisie procesy zainicjowane przez Operatora od procesów zainicjowanych poleceniem Wykonawcy (rozdz. 1.3) | Przegląd kolumny „inicjator” w opisie Process Monitor, rozdz. 3 |
| Żaden opisany element interfejsu nie odwołuje się do klasy `.dn-*` ani żetonu `--dn-*` nieobecnego w arkuszach `design/zasoby/` | `grep` nazwy klas i żetonów użytych w dokumencie względem arkuszy źródłowych |
| Rozłączenie klienta nie przerywa działania powłoki — sesja pozostaje trwała po stronie serwerowej (rozdz. 1.3) | Przegląd opisu trwałości sesji w rozdz. 1.3 i stanów Terminal Tabs w rozdz. 3 |


---

## Załącznik — pełny wykaz komend kontraktu modułu Terminal

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `terminal` — 27 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `terminal.command.exec` | Uruchamia polecenie w karcie terminala; wynik idzie strumieniem fragmentów | `sessionId:string` (wym)<br>`command:string` (wym)<br>`initiator:ProcessInitiator` (opc)<br>`timeoutMs:int` (opc) | `process:TerminalProcess` (wym) |
| `terminal.file.read` | Oddaje treść pliku z katalogu roboczego karty. Bez tej komendy klient czyta manifest projektu poleceniem powłoki, więc wykrycie zadań zależy od programu wypisującego plik i od składni każdej z sześciu powłok | `sessionId:string` (wym)<br>`path:string` (wym)<br>`tail:int` (opc)<br>`maxBytes:int` (opc) | `content:string` (wym)<br>`path:string` (wym)<br>`truncated:bool` (wym)<br>`truncatedBytes:int64` (opc)<br>`sizeBytes:int64` (opc) |
| `terminal.host.list` | Zwraca książkę hostów Operatora | `group:string` (opc)<br>`query:string` (opc) | `hosts:TerminalHost[]` (wym)<br>`total:int` (wym) |
| `terminal.host.remove` | Usuwa wpis książki hostów. Karty już otwarte do tego hosta biegną dalej | `hostId:string` (wym) | `removed:bool` (wym) |
| `terminal.host.save` | Zapisuje wpis książki hostów albo podmienia istniejący. Bez tej komendy książka żyje jedno posiedzenie przeglądarki i ginie przy odświeżeniu strony | `host:TerminalHost` (wym) | `host:TerminalHost` (wym)<br>`created:bool` (wym) |
| `terminal.key.generate` | Wytwarza parę kluczy SSH na maszynie rdzenia. Hasło klucza wchodzi odwołaniem do sejfu, nigdy treścią — zgodnie z zasadą zapisaną w kontrakcie przy zmiennej środowiska | `name:string` (wym)<br>`keyType:TerminalKeyType` (wym)<br>`comment:string` (opc)<br>`passphraseRef:string` (opc) | `key:TerminalSshKey` (wym) |
| `terminal.key.import` | Wciąga do wykazu klucz leżący już na maszynie rdzenia. Klucz wskazuje się ścieżką, a nie treścią: materiał kluczowy nie ma powodu przechodzić przez łącze | `name:string` (wym)<br>`path:string` (wym) | `key:TerminalSshKey` (wym) |
| `terminal.key.list` | Zwraca wykaz kluczy SSH znanych rdzeniowi wraz z ich odciskami | — | `keys:TerminalSshKey[]` (wym)<br>`total:int` (wym) |
| `terminal.key.remove` | Zdejmuje klucz z wykazu. Wpisy książki hostów wskazujące ten klucz tracą wskazanie i wracają do klucza domyślnego konfiguracji maszyny | `keyId:string` (wym)<br>`deleteFiles:bool` (opc) | `removed:bool` (wym)<br>`detachedHostIds:string[]` (opc) |
| `terminal.output.read` | Oddaje wyjście jednego procesu terminala — stdout, stderr, kod wyjścia i stan — jako jedną odpowiedź, bez zapisywania się na strumień. Komenda `terminal.command.exec` kończy się w chwili startu procesu (kompilacja trwa dłużej niż każde sensowne oczekiwanie na odpowiedź, a rozłączenie klienta nie ma prawa jej przerwać), więc wyniku nieść nie może i nigdy nie będzie mogła. Stąd model czyta, co polecenie wypisało. Źródłem jest ten sam dziennik zbiorczego wyjścia, na którym stoi `terminal.output.stream` — drugiej pompy nie ma. Historia żyje jeden bieg rdzenia: po restarcie wyjście jest puste i odpowiedź mówi to wprost | `processId:string` (wym)<br>`tail:int` (opc)<br>`waitMs:int` (opc) | `stdout:string` (wym)<br>`stderr:string` (wym)<br>`exitCode:int` (opc)<br>`status:TerminalProcessStatus` (wym)<br>`truncated:bool` (wym)<br>`truncatedBytes:int` (opc) |
| `terminal.output.stream` | Zapisuje okno na zbiorcze wyjście wszystkich otwartych kart terminala i zwraca ogon historii | `sessionId:string` (opc)<br>`windowId:string` (opc)<br>`terminalSessionIds:string[]` (opc)<br>`tail:int` (opc) | `lines:TerminalOutputLine[]` (wym)<br>`subscribed:bool` (wym) |
| `terminal.process.kill` | Kończy proces sygnałem łagodnym albo wymuszonym | `processId:string` (wym)<br>`force:bool` (opc) | `process:TerminalProcess` (wym) |
| `terminal.process.list` | Zwraca procesy rejestru rdzenia, w tym uruchomione przez model | `windowId:string` (opc)<br>`sessionId:string` (opc)<br>`status:TerminalProcessStatus` (opc)<br>`initiator:ProcessInitiator` (opc) | `processes:TerminalProcess[]` (wym) |
| `terminal.process.suspend` | Wstrzymuje albo wznawia proces rejestru rdzenia. Dziś rdzeń umie proces wyłącznie zakończyć, więc długie zadanie da się tylko przerwać | `processId:string` (wym)<br>`resume:bool` (opc) | `process:TerminalProcess` (wym)<br>`supported:bool` (wym) |
| `terminal.script.lint` | Poddaje treść skryptu analizie statycznej i formatowaniu. Analizę prowadzą programy spoza instalacji Danaco Console, więc odpowiedź mówi wprost, czy narzędzie było dostępne — pusty wykaz uwag przy braku narzędzia znaczyłby fałszywie treść bez zastrzeżeń | `content:string` (wym)<br>`shell:TerminalShell` (wym)<br>`format:bool` (opc) | `findings:TerminalLintFinding[]` (wym)<br>`formatted:string` (opc)<br>`analyzerAvailable:bool` (wym)<br>`analyzer:string` (wym) |
| `terminal.script.list` | Zwraca pozycje biblioteki skryptów | `kind:TerminalScriptKind` (opc)<br>`tag:string` (opc)<br>`query:string` (opc) | `scripts:TerminalScript[]` (wym)<br>`total:int` (wym) |
| `terminal.script.remove` | Usuwa pozycję biblioteki wraz ze wszystkimi jej wersjami | `scriptId:string` (wym) | `removed:bool` (wym) |
| `terminal.script.save` | Zapisuje skrypt albo snippet biblioteki jako kolejną wersję. Bez tej komendy biblioteka żyje jedno posiedzenie, a jedyną drogą jej zachowania jest wywóz do pliku | `script:TerminalScript` (wym) | `script:TerminalScript` (wym)<br>`created:bool` (wym) |
| `terminal.session.close` | Zamyka kartę powłoki. Dziś zamknięcie karty żyje wyłącznie w widoku klienta: powłoka i jej procesy biegną dalej, a rdzeń o zamknięciu nie wie | `sessionId:string` (wym)<br>`force:bool` (opc) | `session:TerminalSession` (wym)<br>`stoppedProcessIds:string[]` (opc) |
| `terminal.session.list` | Zwraca karty powłoki znane rdzeniowi. Rdzeń odtwarza karty przy starcie, ale klient po ponownym połączeniu nie ma jak ich zobaczyć i zaczyna wykaz od pustego | `windowId:string` (opc)<br>`status:TerminalSessionStatus` (opc)<br>`includeExited:bool` (opc) | `sessions:TerminalSession[]` (wym)<br>`total:int` (wym) |
| `terminal.session.open` | Otwiera kartę terminala jako odrębną sesję powłoki | `windowId:string` (wym)<br>`shell:TerminalShell` (wym)<br>`workingDir:string` (opc)<br>`title:string` (opc)<br>`environment:EnvironmentVariable[]` (opc)<br>`remoteTarget:string` (opc)<br>`remotePort:int` (opc)<br>`hostId:string` (opc)<br>`containerRef:TerminalContainerRef` (opc)<br>`serialDevice:string` (opc)<br>`serialBaudRate:int` (opc) | `session:TerminalSession` (wym) |
| `terminal.tunnel.close` | Zamyka przekierowanie portu | `tunnelId:string` (wym) | `tunnel:TerminalTunnel` (wym) |
| `terminal.tunnel.list` | Zwraca przekierowania portów wraz z ich stanem | `windowId:string` (opc)<br>`status:TerminalTunnelStatus` (opc) | `tunnels:TerminalTunnel[]` (wym)<br>`total:int` (wym) |
| `terminal.tunnel.open` | Zakłada przekierowanie portu. Dziś tunel da się założyć wyłącznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowość są niewidoczne | `windowId:string` (wym)<br>`kind:TerminalTunnelKind` (wym)<br>`hostId:string` (opc)<br>`remoteTarget:string` (opc)<br>`localPort:int` (opc)<br>`remoteHost:string` (opc)<br>`remotePort:int` (opc) | `tunnel:TerminalTunnel` (wym) |
| `terminal.watch.list` | Zwraca obserwacje plików wraz z licznikiem wyzwoleń | `windowId:string` (opc)<br>`status:TerminalWatchStatus` (opc) | `watches:TerminalWatch[]` (wym)<br>`total:int` (wym) |
| `terminal.watch.start` | Zakłada obserwację plików uruchamiającą polecenie przy ich zmianie. Kontrakt daje dziś wyzwalacz plikowy automatyce, a nie karcie powłoki | `sessionId:string` (wym)<br>`pattern:string` (wym)<br>`command:string` (wym)<br>`debounceMs:int` (opc)<br>`recursive:bool` (opc) | `watch:TerminalWatch` (wym) |
| `terminal.watch.stop` | Zatrzymuje obserwację plików. Polecenie już uruchomione biegnie dalej | `watchId:string` (wym) | `watch:TerminalWatch` (wym) |

**Zdarzenia obszaru `terminal` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `terminal.process.changed` | Zmiana procesu rejestru rdzenia serwera | `change:ChangeKind`, `process:TerminalProcess` |


Razem w wykazie: **27 komend** z jednego obszaru kontraktu.

---

*Koniec dokumentu. Moduł Terminal — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
