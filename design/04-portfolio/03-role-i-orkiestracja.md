# Danaco Console — Portfolio Design Identity — P3: Role i orkiestracja (MultitaskingAI)

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
| **Dokument** | Plansza portfolio P3 — tożsamość systemu ról i orkiestracji środowiska MultitaskingAI (opracowanie towarzyszące planszy `03-role-i-orkiestracja.html`) |
| **Odbiorcy** | Designer (masa wizualna węzła, hierarchia, forma paneli i diagramów) · Deweloper (komplet nazw ról, akcji, sekcji, stanów i przejść) · Odbiorca portfolio (jak wygląda zespół modeli w działaniu) |
| **Zakres** | Pięć pozycji zespołu · asymetria koordynator–wykonawca jako decyzja projektowa · sześć sekcji panelu orkiestracji · silnik kolejek (5 zasięgów, 11 akcji, cykl życia zadania) · trzy rodzaje zależności i mapa zależności · sześć poziomów hierarchii decyzji · trzy szablony zespołu · praca ciągła 24/7/365 z modułem Automations · relacja Agents → MultitaskingAI |
| **Czego NIE zawiera** | Specyfikacji wnętrza okien roboczych (to zakres `dok/projekt-ui/srodowiska/multitaskingai.md` rozdz. 9–10 oraz prototypów w `WYNIK/05-okna/srodowiska/`) · nowych ról, akcji, sekcji ani stanów · żadnej nazwy, liczby ani mechanizmu bez pokrycia w dokumentacji |
| **Zasada nadrzędna** | Każda nazwa własna, każda liczba i każde przejście stanu na planszy pochodzi z dokumentacji merytorycznej albo z odczytu realnego pliku w repozytorium wynikowym. Zero elementów ilustracyjnych. |

---

## Spis treści

1. [Czym jest tożsamość systemu ról](#1-czym-jest-tożsamość-systemu-ról)
2. [Pięć pozycji zespołu — zestawienie zbiorcze](#2-pięć-pozycji-zespołu--zestawienie-zbiorcze)
3. [Karty tożsamości pięciu pozycji](#3-karty-tożsamości-pięciu-pozycji)
4. [Asymetria koordynator–wykonawca jako decyzja projektowa](#4-asymetria-koordynatorwykonawca-jako-decyzja-projektowa)
5. [Sześć sekcji panelu orkiestracji](#5-sześć-sekcji-panelu-orkiestracji)
6. [Silnik kolejek — zasięgi, jedenaście akcji, cykl życia](#6-silnik-kolejek--zasięgi-jedenaście-akcji-cykl-życia)
7. [Stany — proces, rola, zadanie](#7-stany--proces-rola-zadanie)
8. [Orkiestracja — trzy rodzaje zależności i mapa](#8-orkiestracja--trzy-rodzaje-zależności-i-mapa)
9. [Hierarchia decyzji — sześć poziomów](#9-hierarchia-decyzji--sześć-poziomów)
10. [Trzy szablony konfiguracji zespołu](#10-trzy-szablony-konfiguracji-zespołu)
11. [Praca ciągła 24/7/365 — integracja z Automations](#11-praca-ciągła-247365--integracja-z-automations)
12. [Agents → MultitaskingAI: rola jest ekspertem z fabryki](#12-agents--multitaskingai-rola-jest-ekspertem-z-fabryki)
13. [Decyzje projektowe planszy](#13-decyzje-projektowe-planszy)
14. [Źródła](#14-źródła)

---

## 1. Czym jest tożsamość systemu ról

MultitaskingAI jest jedynym środowiskiem platformy, w którym jednostką organizacji pracy nie jest pojedynczy użytkownik wspierany przez AI, lecz **zespół modeli i agentów realizujących wspólny proces** (`multitaskingai.md`, Wprowadzenie oraz rozdz. 1.1). Konsekwencja tego rozstrzygnięcia jest widoczna w każdym elemencie interfejsu:

| Cecha | TalkIn / WorkSpace / CodeStudio | MultitaskingAI |
|---|---|---|
| Jednostka organizacji pracy | Pojedynczy użytkownik wspierany przez AI, w ramach jednego zadania | Zespół modeli i agentów realizujących wspólny proces |
| Zawartość bocznej nawigacji | Lista modułów dostępnych w środowisku | Panel orkiestracji — sześć sekcji sterowania zespołem |
| Podstawowa jednostka pracy | Moduł → okno operacyjne | Rola → akcja silnika kolejek → orkiestracja |
| Obecność w macierzy dostępności modułów | Tak — jako kolumna macierzy | Nie — środowisko nie udostępnia modułów w bocznej nawigacji |
| Domyślny tryb pracy | Interakcja konwersacyjna, jednorazowa lub sesyjna | Zdolność do pracy ciągłej po skonfigurowaniu integracji z Automations |
| Rola Always On Display | Doradztwo kontekstowe | Doradztwo kontekstowe oraz funkcja obserwatora lub operatora procesu |

Źródło tabeli: `multitaskingai.md` rozdz. 1.1.

**Fundamentalne rozstrzygnięcie systemu ról** (rozdz. 3, akapit wprowadzający): podział na rolę wykonującą (Executor) i rolę zarządzającą (Coordinator) jest celowy — „obie odpowiedzialności są celowo rozdzielone, aby żaden pojedynczy model nie musiał jednocześnie planować i realizować”. To zdanie jest źródłem całej warstwy wizualnej planszy: **rozdzielenie planowania od wykonania musi być widoczne, zanim czytelnik przeczyta jakąkolwiek etykietę.**

Środowisko odróżnia się także od modułu Roundtable (rozdz. 1.2):

| Cecha | Moduł Roundtable | Środowisko MultitaskingAI |
|---|---|---|
| Forma wielomodelowości | Debata i wypracowanie konsensusu | Zorganizowany podział pracy między wyspecjalizowane role |
| Sposób pracy modeli | Kilka modeli odpowiada równolegle na to samo zagadnienie | Role realizują rozdzielone zakresy wspólnego procesu |
| Mechanizmy porządkujące | Moderator Panel, Consensus Panel | Podział ról, silnik kolejek, orkiestracja zależności |
| Punkt ciężkości | Debata | Zorganizowany podział pracy |

---

## 2. Pięć pozycji zespołu — zestawienie zbiorcze

Środowisko udostępnia **cztery okna robocze** odpowiadające czterem rolom oraz **jeden mechanizm** — Subagent Network — uruchamiany przez rolę wykonawczą. Razem: pięć pozycji zespołu.

### 2.1. Cel, wejście, wyjście (rozdz. 3.0)

| Rola | Cel roli | Wejście | Wyjście |
|---|---|---|---|
| **Executor 1 — główny wykonawca** | Faktyczna realizacja pracy: tworzenie dokumentów, analiza danych, programowanie, projektowanie, budowa aplikacji, przetwarzanie materiałów; nie zarządza procesem | Zadanie z kolejki przypisanej roli, polecenie użytkownika, prompt zbudowany przez Coordinatora | Rezultat pracy, status wykonania do silnika kolejek, podzadania do Subagent Network przy złożonych zleceniach |
| **Subagent Network** | Mechanizm uruchamiany przez wykonawcę, nie odrębna rola | Zakres zadania wydzielony przez wykonawcę (`split`) | Wynik cząstkowy agregowany przez wykonawcę (`merge`) |
| **Coordinator — koordynator** | Nie tworzy końcowego produktu; projektuje i kontroluje sposób jego wytwarzania | Cel procesu ustalony przez użytkownika, wyniki cząstkowe wykonawców, raporty Executora 3 / Validatora | Plan pracy, prompty dla wykonawców, polecenia sterujące kolejką i orkiestracją |
| **Executor 2 — równoległy wykonawca** | Drugi niezależny model wykonawczy, procesy równoległe (backend/frontend, badania/raport, kod/dokumentacja) | Jak Executor 1; dodatkowo wyniki pośrednie Executora 1, zależnie od trybu współpracy | Jak Executor 1; wynik przekazywany zależnie od trybu współpracy |
| **Executor 3 / Validator — czwarty model** | Funkcja kontrolna lub doradcza domykająca zespół; odpowiada za jakość, zgodność lub rozstrzyganie rozbieżności | Rezultat pracy wykonawców, wynik pracy Subagent Network po agregacji | Ocena, raport, decyzja o przekazaniu dalej lub zwrocie do Coordinatora (`retry`) |

### 2.2. Zestawienie zbiorcze ról (rozdz. 3.6)

| Rola | Tworzy produkt końcowy | Zarządza procesem | Uruchamia Subagent Network | Przypisanie wykonawcy | Sekcja panelu orkiestracji | Okno robocze |
|---|---|---|---|---|---|---|
| Executor 1 | Tak | Nie | Tak (do 15) | Agent lub model bazowy | Role | Executor Chat |
| Subagent Network | Nie — wynik agreguje wykonawca | Nie | — | Podagenci uruchamiani przez Executora | Role (rozwinięcie karty roli) | — (zagnieżdżone w Executor Chat) |
| Coordinator | Nie | Tak | Nie | Agent lub model bazowy | Role | Coordinator Chat |
| Executor 2 | Tak | Nie | Tak (do 15) | Agent lub model bazowy | Role | Executor Chat |
| Executor 3 / Validator | Zależnie od wcielenia (zwykle nie) | Nie | Nie | Agent lub model bazowy | Role | Results Analyzer |

Warstwą centralną nadrzędną wobec wszystkich pięciu pozycji jest **Always On Display** (rozdz. 2, rozdz. 3.6 przypis).

---

## 3. Karty tożsamości pięciu pozycji

### 3.1. Executor 1 — główny wykonawca

| Grupa narzędzi | Zawartość |
|---|---|
| Przypisanie wykonawcy | Agent skonfigurowany w module Agents (siedem komponentów definicji: model bazowy, tożsamość, instrukcje systemowe, skille, pluginy i konektory, pamięć, uprawnienia) albo model bazowy podłączony bezpośrednio jednym z czterech kanałów |
| Kanał modelu | API · CLI · SSH · HTTP — wybór dokonywany per rola, nadpisujący kanał domyślny zapisany w definicji agenta |
| Rozszerzenia dostępne wykonawcy | Wtyczki i umiejętności (Skills Manager); konektory i serwery MCP (Connectors Manager) |
| Pamięć | Poziomy: globalna, projekt, sesja, środowisko |
| Akcje silnika kolejek dostępne roli | `dequeue` · `split` · `merge` |
| Uruchamianie Subagent Network | Do 15 równoczesnych podagentów |
| Profil izolacji | Poziom zasięgu „Rola” — najwyższe pierwszeństwo spośród siedmiu poziomów; konfigurowalny, domyślnie brak aktywnej izolacji technicznej |
| Okno robocze | **Executor Chat** |

Charakterystyka (rozdz. 3.1): rola podstawowa — punkt wyjścia dla każdej konfiguracji zespołu.

### 3.2. Subagent Network

| Aspekt | Zawartość |
|---|---|
| Charakter | Mechanizm, nie odrębna rola zespołu — uruchamiany przez wykonawcę (Executor 1 lub Executor 2) |
| Zakres | Do 15 wyspecjalizowanych podagentów jednocześnie |
| Przykładowe wcielenia | Agent UI, Agent Backend, Agent API, Agent Security, Agent Database, Agent Testing — zestaw przykładowy, nie zamknięty |
| Oprzyrządowanie podagenta | Pełna definicja komponentu własnego rodzaju „agent” albo wąsko sprofilowany model bazowy dedykowany jednemu zakresowi zadania |
| Mechanizm działania | `split` → przydział zakresów → `merge` (agregacja przez wykonawcę macierzystego) |
| Punkty izolacji | Brak odrębnego poziomu zasięgu — dziedziczy profil izolacji roli macierzystej |

Executor 2 dysponuje **własną, niezależną instancją** Subagent Network (rozdz. 3.4) — stąd do **30 jednostek wykonawczych** działających jednocześnie w jednym procesie (rozdz. 16, wiersz 4).

### 3.3. Coordinator — koordynator

**Pięć obszarów odpowiedzialności** (rozdz. 3.3):

| Obszar | Zawartość |
|---|---|
| Planowanie | Etapy projektu, zależności, harmonogramy, logika przepływu pracy |
| Podział pracy | Przypisywanie ról, zadań, odpowiedzialności, modeli wykonawczych |
| Budowa promptów | Dynamiczne tworzenie i optymalizacja promptów dla wykonawców |
| Sterowanie procesem | Start, stop, pauza, wznowienie, przekazanie, powtórzenie, walidacja |
| Zarządzanie kolejką | Budowa zaawansowanych struktur kolejkujących |

**Sterowanie procesem → akcje silnika kolejek:**

| Operacja sterująca | Odpowiadająca akcja lub mechanizm |
|---|---|
| Start | Uruchomienie przepływu zadań |
| Stop | Trwałe zakończenie przepływu |
| Pauza | `pause` |
| Wznowienie | `resume` |
| Przekazanie | `route` — skierowanie zadania do innej roli |
| Powtórzenie | `retry` |
| Walidacja | Skierowanie wyniku do Executora 3 / Validatora przed uznaniem etapu za zakończony |

**Oprzyrządowanie:** widok planu (panel etapów projektu z zależnościami i harmonogramem), kreator promptów, widok stanu kolejki, akcje `enqueue`, `delay`, `retry`, `pause`, `resume`, `route`, `branch`, `condition` — pełny zestaw poza `dequeue`, `split`, `merge`. Okno robocze: **Coordinator Chat**.

### 3.4. Executor 2 — równoległy wykonawca

Oprzyrządowanie tożsame z Executorem 1. Rolę odróżniają **cztery tryby współpracy** z Executorem 1 (rozdz. 3.4):

| Tryb | Opis | Kierunek przepływu danych | Przykład zastosowania |
|---|---|---|---|
| Praca niezależna | Odrębne zadania bez wzajemnej zależności; synchronizacja przez wspólną kolejkę lub wspólny cel procesu | Brak przepływu bezpośredniego — tory zbiegają się w Coordinatorze lub Executorze 3 / Validatorze | Backend i frontend budowane równolegle w module Apps |
| Przekazywanie wyników | Wynik jednego wykonawcy jest wejściem drugiego; przepływ jednokierunkowy | Executor 1 → Executor 2 (lub odwrotnie), przez `route` | Badania (Executor 1) przekazane do redakcji raportu (Executor 2) |
| Praca naprzemienna | Wykonawcy przejmują zadanie kolejno, każdy realizując kolejny etap tego samego procesu | Executor 1 → Executor 2 → Executor 1 → …, sterowane przez Coordinatora | Kod i dokumentacja techniczna aktualizowane naprzemiennie |
| Praca iteracyjna | Wykonawcy poprawiają wzajemnie swoje rezultaty w powtarzających się cyklach | Dwukierunkowa pętla Executor 1 ⇄ Executor 2, zamykana decyzją Coordinatora lub oceną Executora 3 / Validatora | Iteracyjne dopracowywanie kodu i testów do przejścia walidacji |

### 3.5. Executor 3 / Validator — czwarty model

**Siedem przykładowych wcieleń** (rozdz. 3.5) — zestawienie przykładowe, nie wyczerpujące:

| Wcielenie | Funkcja | Ikona odniesienia |
|---|---|---|
| Validator | Kontrola jakości pracy pozostałych modeli przed uznaniem etapu za zakończony | `ptaszek-kolo` |
| Reviewer | Recenzja wyników pod kątem poprawności i kompletności | `oko` |
| Security Auditor | Ocena bezpieczeństwa rozwiązania wytworzonego przez wykonawców | `tarcza` |
| Architect | Ocena zgodności rozwiązania z założeniami architektonicznymi | `kod` |
| Product Owner | Ocena zgodności rezultatu z wymaganiami biznesowymi | `dokument` |
| QA Lead | Przygotowanie przypadków testowych i raportów jakości | `ptaszek` |
| Arbitrator | Rozstrzyganie konfliktów powstałych między wynikami Executora 1 i Executora 2 | `waga` |

**Uwaga wykonawcza dla planszy:** ikona `waga` przypisana wcieleniu Arbitrator **nie występuje** w komplecie 82 plików `WYNIK/zasoby/ikony/svg/` (potwierdzone odczytem katalogu). Plansza nie rysuje jej zastępczo — wcielenie Arbitrator opatrzone jest ikoną `walidator` (rodzina roli) i etykietą tekstową, zgodnie z zasadą „stan nigdy samym kolorem, zawsze ikona lub etykieta” (KANON rozdz. 9). Odnotowane w decyzjach projektowych (rozdz. 13).

**Oprzyrządowanie:** wejście oceny (rezultat pracy Executora 1 i Executora 2, wynik Subagent Network po agregacji), akcje `retry` oraz udział w `condition`/`branch` jako źródło warunku, panel porównania wyników dwóch wykonawców (wcielenie Arbitrator). Okno robocze: **Results Analyzer**.

---

## 4. Asymetria koordynator–wykonawca jako decyzja projektowa

### 4.1. Podstawa kontraktowa

`KIERUNEK.md` rozdz. 2 (pokrętło `WARIANCJA_PROJEKTOWA` = 4/10): *„Kokpit wymaga przewidywalności. **Asymetria tylko tam, gdzie niesie hierarchię** (strefy E1, relacja koordynator–wykonawca). Zero ozdobnego chaosu.”*

`KIERUNEK.md` rozdz. 4, katalog anty-domyślnych: odruch **„trzy równe karty funkcji”** ma być zastąpiony przez **„strefy o malejącej masie (E1), asymetria koordynator–wykonawca (E4)”**. To samo powtarza KANON rozdz. 1.

Wniosek dla planszy: **układ zespołu nie może być siatką pięciu równych kart.** Węzły muszą mieć różną masę wizualną, a różnica masy musi być czytelna jako informacja, nie jako ozdoba.

### 4.2. Reguła masy przyjęta na planszy

Plansza wiąże masę wizualną węzła **wprost z poziomem hierarchii decyzji** (rozdz. 8.1 dokumentu źródłowego). Dokumentacja nie rozstrzyga formy graficznej — rozstrzygnięcie należy do warstwy wizualnej i zostaje tu odnotowane (KANON rozdz. 10 pkt 5).

> **masa wizualna = 7 − poziom hierarchii decyzji**

| Poziom (rozdz. 8.1) | Podmiot | Masa | Odwzorowanie wizualne na planszy |
|:--:|---|:--:|---|
| 1 | Użytkownik (Operator) | 6 | Pas nadrzędny nad całym układem, pełna szerokość, obrys mocny |
| 2 | Always On Display (tryb operatora) | 5 | Pas warstwy centralnej, pełna szerokość, kropka sygnału z tętnem |
| 3 | Coordinator | 4 | Węzeł najszerszy spośród ról, dwie kolumny siatki, nagłówek `--dn-fs-lg` |
| 4 | Executor 3 / Validator | 3 | Węzeł pojedynczej kolumny, nagłówek `--dn-fs-md`, obrys standardowy |
| 5 | Executor 1 · Executor 2 | 2 | Para węzłów **równych sobie** — symetria wewnątrz pary, asymetria wobec Coordinatora |
| 6 | Subagent Network | 1 | Węzły najmniejsze, zagnieżdżone wewnątrz węzła wykonawcy macierzystego |

**Dlaczego para wykonawców pozostaje symetryczna.** Executor 2 ma oprzyrządowanie „tożsame z Executorem 1” (rozdz. 3.4) i ten sam poziom hierarchii (rozdz. 8.1, poziom 5). Zróżnicowanie ich masy byłoby asymetrią bez pokrycia w hierarchii — czyli dokładnie tym „ozdobnym chaosem”, którego zakazuje kontrakt kierunku. Symetria pary jest tu równie znaczącą decyzją jak asymetria wobec koordynatora.

### 4.3. Przełącznik „pokaż masę wizualną”

Plansza udostępnia przełącznik nakładający na diagram oznaczenia wagi: numer poziomu hierarchii, wartość masy oraz opis odwzorowania. Bez przełącznika diagram czyta się jako układ zespołu; z przełącznikiem — jako dowód, że układ jest wyprowadzony z hierarchii, a nie z estetyki. Ruch przy przełączeniu ogranicza się do mikroprzejścia (`--dn-czas-2`), zgodnie z `INTENSYWNOSC_RUCHU` 3/10.

---

## 5. Sześć sekcji panelu orkiestracji

Panel orkiestracji jest **odpowiednikiem bocznej nawigacji modułów** w pozostałych środowiskach: nawiguje po tym, czym w tym środowisku faktycznie się steruje — po rolach, kolejkach i orkiestracji — a nie po modułach (rozdz. 6.1).

| # | Sekcja | Zawartość | Typowe działania | Powiązane mechanizmy i okna | Rozdział źródła |
|:--:|---|---|---|---|:--:|
| 1 | **Zespoły** | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek | Zapis nowego zespołu, wczytanie zapisanego zespołu, duplikowanie, usunięcie | Szablony konfiguracji zespołu | 14 |
| 2 | **Role** | Cztery okna robocze: Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów; uruchomienie i podgląd Subagent Network | Przypisanie agenta lub modelu bazowego, wybór kanału modelu, ustalenie trybu współpracy, aktywacja Subagent Network | Agent Builder, Permissions Center | 3, 9 |
| 3 | **Kolejki** | Definicje kolejek pięciu zasięgów oraz ich akcje | Utworzenie kolejki, przypisanie zasięgu, konfiguracja reguł warunkowych, podgląd zawartości | Silnik kolejek, Queue Manager | 4 |
| 4 | **Orkiestracja** | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami | Definicja zależności sekwencyjnych, warunkowych, równoległych | Orchestrator | 5 |
| 5 | **Harmonogram i automatyki** | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations | Ustalenie reguły czasowej i cykliczności, wpięcie lub odpięcie automatyki | Scheduler, Execution Monitor | 7 |
| 6 | **Monitor procesu** | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display | Podgląd stanu ról i kolejek, historia przebiegów, przełączenie AOD między trybami | Always On Display, Mobile | 2, 8, 12 |

Źródło: `multitaskingai.md` rozdz. 6.2 — nazwy sekcji przejęte dosłownie, bez skrótów.

**Konfigurowalność (rozdz. 6.5).** Kolejność i widoczność sekcji podlegają konfiguracji; zestaw sześciu sekcji jest **zestawem domyślnym**. Panel nigdy nie wymusza ukrycia żadnej sekcji ani nie blokuje przywrócenia zestawu domyślnego — „brak ustawienia oznacza wartość domyślną, nigdy brak dostępności”.

**Ikony sekcji** (Załącznik A.1 dokumentu źródłowego, zweryfikowane wobec `WYNIK/zasoby/ikony/svg/`): `gwiazdka` (Zespoły), `uzytkownik` (Role), `filtr` (Kolejki), `link-zewnetrzny` (Orkiestracja), `zegar` (Harmonogram i automatyki), `oko` (Monitor procesu).

---

## 6. Silnik kolejek — zasięgi, jedenaście akcji, cykl życia

W warstwie komunikacji silnik jest realizowany poleceniem `queue.action`, przyjmującym jako parametr jedną z jedenastu akcji (rozdz. 4).

### 6.1. Pięć zasięgów kolejek (rozdz. 4.1)

| Zasięg | Opis | Przykład zastosowania |
|---|---|---|
| Globalna | Obejmuje wszystkie procesy MultitaskingAI użytkownika | Kolejka wspólna dla wszystkich aktywnych zespołów |
| Lokalna | Obejmuje jeden proces orkiestracji (jedną kartę sesji) | Kolejka etapów pojedynczego projektu realizowanego w module Apps |
| Dla modelu | Przypisana konkretnemu modelowi bazowemu | Kolejka zadań przetwarzanych przez jeden, wskazany model |
| Dla agenta | Przypisana konkretnemu agentowi z modułu Agents | Kolejka zadań przypisanych agentowi pełniącemu rolę Executora 1 |
| Dla projektu | Przypisana projektowi prowadzonemu w module Workspace | Kolejka zadań realizowanych w ramach jednego, izolowanego projektu |

### 6.2. Jedenaście akcji (rozdz. 4.2) — komplet

| # | Akcja | Grupa funkcjonalna | Działanie | Typowy inicjator | Reprezentacja w interfejsie |
|:--:|---|---|---|---|---|
| 1 | `enqueue` | Obieg zadań | Dodanie nowego zadania do wskazanej kolejki | Coordinator, Automations (po integracji) | Przycisk ikonowy w wierszu tabeli kolejki, ikona `plus` |
| 2 | `dequeue` | Obieg zadań | Pobranie kolejnego zadania z kolejki do wykonania | Executor 1, Executor 2 | Automatyczne wywołanie przy wolnym slocie wykonawcy; widoczne jako etykieta w Executor Chat |
| 3 | `delay` | Sterowanie czasem | Odroczenie wykonania zadania o zadany czas lub do spełnienia warunku | Coordinator | Pole czasu w wierszu zadania, ikona `zegar` |
| 4 | `retry` | Ponawianie | Ponowienie zadania po niepowodzeniu wykonania lub po negatywnej ocenie walidacji | Coordinator, Executor 3 / Validator | Przycisk ikonowy, ikona `odswiez` |
| 5 | `pause` | Wstrzymywanie | Wstrzymanie przetwarzania wskazanej kolejki | Coordinator, Always On Display (operator) | Przycisk ikonowy, ikona `zatrzymaj` |
| 6 | `resume` | Wstrzymywanie | Wznowienie wstrzymanej kolejki | Coordinator, Always On Display (operator) | Przycisk ikonowy, ikona `uruchom` |
| 7 | `split` | Podział i scalanie | Podział zadania na mniejsze podzadania | Executor 1 lub Executor 2, przy uruchamianiu Subagent Network | Przycisk „Uruchom Subagent Network” w Executor Chat |
| 8 | `merge` | Podział i scalanie | Scalenie wyników wielu zadań lub podzadań w jeden rezultat | Executor 1 lub Executor 2, po agregacji wyników Subagent Network | Przycisk „Scal wyniki” w panelu Subagent Network |
| 9 | `route` | Kierowanie warunkowe | Skierowanie zadania do wskazanej roli lub kolejki na podstawie reguły | Coordinator, orkiestracja | Selektor roli docelowej w wierszu zadania, ikona `strzalka-prawo` |
| 10 | `branch` | Kierowanie warunkowe | Rozgałęzienie przepływu pracy na alternatywne ścieżki | Coordinator | Kreator reguły rozgałęzienia w sekcji Orkiestracja |
| 11 | `condition` | Kierowanie warunkowe | Warunkowe wykonanie kolejnego kroku procesu w zależności od wyniku poprzedniego | Coordinator, orkiestracja | Kreator warunku w sekcji Orkiestracja |

**Rozkład akcji na role:**

| Rola | Akcje dostępne | Źródło |
|---|---|---|
| Executor 1 / Executor 2 | `dequeue`, `split`, `merge` | rozdz. 3.1, 3.4 |
| Coordinator | `enqueue`, `delay`, `retry`, `pause`, `resume`, `route`, `branch`, `condition` | rozdz. 3.3 |
| Executor 3 / Validator | `retry`; udział w `condition`/`branch` jako źródło warunku | rozdz. 3.5 |
| Always On Display (operator) | `pause`, `resume` | rozdz. 2.3 |

Podział jest rozłączny w części wykonawczej: wykonawcy nie dysponują żadną akcją sterującą przepływem, koordynator nie dysponuje żadną akcją wykonawczą. To ta sama zasada, co asymetria z rozdz. 4 niniejszego opracowania — wyrażona w zestawie akcji.

### 6.3. Cykl życia zadania w kolejce (rozdz. 4.3)

```
Coordinator ── enqueue ──► KOLEJKA (globalna / lokalna / modelu / agenta / projektu)
                                │
                    (pause / resume — wg decyzji Coordinatora lub AOD)
                                │
                          dequeue ──► Executor 1 / Executor 2
                                            │
                              split (opcjonalnie) ──► Subagent Network
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
                    (ocena jakości — opcjonalna,
                     konfigurowalna, pomijalna)
                              │
                    retry (przy niepowodzeniu) ──► z powrotem do enqueue
```

### 6.4. Symulacja przebiegu na planszy

Plansza odwzorowuje cykl życia jako **działającą symulację**: przycisk przesuwa przykładowe zadania przez kolejne stany, zmieniając klasę `.dn-krok--*` i wartość `.dn-postep`. Dane przykładowe pochodzą z makiet dokumentu (rozdz. 4.6 i 9.4) i są jako przykładowe oznaczone:

| Zadanie | Rola | Ścieżka symulacji | Źródło danych |
|---|---|---|---|
| Zadanie #128 — API płatności | Executor 1 | `enqueue` → `dequeue` → `split` → `merge` → `route` → ocena pozytywna → Zakończone | Makieta rozdz. 9.4 |
| Zadanie #129 — komponent UI formularza | Executor 2 | `enqueue` → `dequeue` → `route` → ocena negatywna → Do powtórzenia → `retry` → W kolejce | Makieta rozdz. 9.6 (przedmiot porównania), scenariusz rozdz. 15.1 |

Kolejki przykładowe (makieta rozdz. 4.6): **Kolejka główna** (zasięg lokalna, 4 zadania, priorytet 1) · **Kolejka QA** (zasięg dla agenta, 2 zadania, priorytet 2) · **Kolejka Backend** (zasięg dla projektu, 6 zadań, priorytet 1).

---

## 7. Stany — proces, rola, zadanie

### 7.1. Stany przebiegu procesu (rozdz. 12.1)

| Stan | Znaczenie | Sygnalizacja wizualna | Możliwe przejścia |
|---|---|---|---|
| Nieuruchomiony | Zespół skonfigurowany, proces jeszcze nie wystartował | Plakietka neutralna „nieuruchomiony” | → uruchomiony (Start) |
| Uruchomiony | Proces aktywnie wykonuje etapy planu Coordinatora | Plakietka sukcesu, kropka pulsująca | → wstrzymany, → zakończony, → błąd |
| Wstrzymany | Proces zatrzymany akcją `pause`, stan zachowany | Plakietka ostrzeżenia „wstrzymany” | → uruchomiony (`resume`), → zatrzymany trwale (Stop) |
| Oczekujący na decyzję | Etap wymaga zatwierdzenia (Executor 3 / Validator, AOD operator lub użytkownik) | Plakietka info, ikona `oko` | → uruchomiony, → wstrzymany, → powtórzony (`retry`) |
| Zakończony sukcesem | Wszystkie etapy planu wykonane, wynik zaakceptowany | Plakietka sukcesu, ikona `ptaszek` | → nieuruchomiony (kolejny przebieg cyklu) |
| Zakończony błędem | Etap nieodwracalnie nieudany bez dalszego `retry` | Plakietka błędu | → nieuruchomiony (po interwencji użytkownika) |
| Zatrzymany trwale | Proces zakończony poleceniem Stop, bez automatycznego wznowienia | Plakietka neutralna „zatrzymany” | → nieuruchomiony (ponowna konfiguracja) |

### 7.2. Stany roli / okna roboczego (rozdz. 12.2)

| Stan | Znaczenie | Sygnalizacja wizualna | Dotyczy |
|---|---|---|---|
| Bezczynna | Rola przypisana, brak aktywnego zadania | Kropka neutralna | Executor 1, Executor 2, Coordinator, Executor 3 / Validator |
| Pracuje | Rola aktywnie przetwarza zadanie | Kropka sukcesu, animacja | Wszystkie role |
| Oczekuje na zależność | Zadanie roli zablokowane regułą orkiestracji do czasu zakończenia etapu poprzedzającego | Kropka ostrzeżenia | Executor 1, Executor 2 |
| Oczekuje na ocenę | Wynik przekazany do Executora 3 / Validatora, oczekuje na werdykt | Kropka info | Executor 1, Executor 2 |
| Błąd wykonania | Zadanie zakończone niepowodzeniem technicznym | Kropka błędu | Wszystkie role |
| Brak przypisania | Karta roli utworzona, lecz bez wskazanego agenta lub modelu | Karta w stanie pustym (`.dn-pusty-stan`) | Wszystkie role |

### 7.3. Stany zadania w kolejce (rozdz. 12.3)

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

**Odwzorowanie stanów zadania na klasy komponentów** (decyzja warstwy wizualnej, KANON rozdz. 5):

| Stan zadania | Klasa kroku | Uzasadnienie |
|---|---|---|
| W kolejce · Odroczone | `.dn-krok` (bez modyfikatora) | Stan spoczynkowy — brak akcentu |
| W realizacji · Podzielone · Scalone · Skierowane · Rozgałęzione | `.dn-krok--pracuje` | Praca w toku — akcent sygnałowy |
| Zakończone | `.dn-krok--poprawny` | Wynik zaakceptowany — akcent sukcesu |
| Do powtórzenia | `.dn-krok--bledy` | Zwrot po negatywnej ocenie — akcent ostrzeżenia |
| Wstrzymane | `.dn-krok--wstrzymany` | Kolejka zatrzymana — akcent zatrzymania |

### 7.4. Stan „niegotowy” — realizacja zasady zero blokad (rozdz. 12.6)

Dokumentacja definiuje osobny stan interfejsu: **„Niegotowy (aktywny, z komunikatem)”** — „wygląd tożsamy ze stanem domyślnym, bez przyciemnienia i bez blokady kursora — element pozostaje w pełni klikalny”. Zastosowanie: akcje, których warunek nie jest jeszcze spełniony (np. `merge` przed zakończeniem wszystkich podagentów). Kliknięcie wywołuje komunikat kontekstowy zamiast wykonania akcji. Plansza odwzorowuje ten stan atrybutem `data-komunikat` (toast), nigdy atrybutem `disabled` (KANON rozdz. 8).

---

## 8. Orkiestracja — trzy rodzaje zależności i mapa

### 8.1. Funkcja orkiestracji (rozdz. 5.1)

Orkiestracja pozwala definiować zależności pomiędzy modelami, agentami, zadaniami, kolejkami, automatyzacjami oraz projektami. W warstwie komunikacji definicja zależności jest realizowana poleceniem `orchestration.define`. Warstwa orkiestracji odpowiada za to, że zadanie Executora 2 nie rozpocznie się przed zakończeniem etapu przypisanego Executorowi 1, jeśli taka zależność została ustalona.

### 8.2. Trzy rodzaje zależności (rozdz. 5.2)

| Rodzaj zależności | Opis | Mechanizm realizujący |
|---|---|---|
| **Sekwencyjna** | Etap kolejnej roli rozpoczyna się dopiero po zakończeniu etapu roli poprzedzającej | `route` w połączeniu z regułą zależności zdefiniowaną przez Coordinatora |
| **Warunkowa** | Dalszy przebieg procesu zależy od wyniku poprzedniego etapu (np. oceny Executora 3 / Validatora) | `condition`, `branch` |
| **Równoległa niezależna** | Dwa lub więcej tory pracy przebiegają bez wzajemnego oczekiwania | Brak zależności zdefiniowanej w orkiestracji — odpowiednik trybu pracy niezależnej (rozdz. 3.4) |

### 8.3. Odpowiedzialność (rozdz. 5.3)

Zależności ustala Coordinator w ramach obszarów „Zarządzanie kolejką” i „Planowanie”; ich respektowanie w toku wykonania zapewnia warstwa orkiestracji działająca w powiązaniu z silnikiem kolejek. **Żadna zależność nie jest wbudowana na stałe** — każda jest możliwością świadomie ustanowioną, w każdej chwili odwracalną.

### 8.4. Mapa zależności (rozdz. 5.5) — siedem krawędzi

Plansza odwzorowuje diagram z rozdz. 5.5 jako interaktywny graf SVG. Krawędzie i ich rodzaje:

| # | Od | Do | Rodzaj zależności | Mechanizm |
|:--:|---|---|---|---|
| 1 | Coordinator | Executor 1 (backend) | Sekwencyjna | `route` + reguła |
| 2 | Coordinator | Executor 2 (frontend) | Sekwencyjna | `route` + reguła |
| 3 | Executor 1 | Executor 2 | Równoległa niezależna | brak reguły w orkiestracji |
| 4 | Executor 1 | Executor 3 / Validator | Sekwencyjna (zbieg torów) | `route` |
| 5 | Executor 2 | Executor 3 / Validator | Sekwencyjna (zbieg torów) | `route` |
| 6 | Executor 3 / Validator | Kolejny etap procesu | Warunkowa — ocena pozytywna | `condition` |
| 7 | Executor 3 / Validator | Coordinator | Warunkowa — ocena negatywna | `retry` |

**Zachowanie interaktywne:** najechanie lub sfokusowanie węzła podświetla cały łańcuch zależności, w którym węzeł uczestniczy — krawędzie przychodzące i wychodzące wraz z węzłami sąsiednimi. Odwzorowuje to element „Widok grafu zależności” z rozdz. 5.4: *„Kliknięcie węzła podświetla powiązane zależności i otwiera panel szczegółów”*.

**Plakietka statusu zgodności** (rozdz. 5.4) — trzy stany: Zgodny (sukces) · Oczekujący (info) · Naruszony (błąd). Kliknięcie przy stanie „naruszony” otwiera szczegóły konfliktu.

---

## 9. Hierarchia decyzji — sześć poziomów

### 9.1. Poziomy (rozdz. 8.1)

| Poziom | Podmiot | Zakres decyzji |
|:--:|---|---|
| **1 (najwyższy)** | Użytkownik (Operator) | Decyzja nadrzędna wobec każdego elementu procesu; może w dowolnym momencie zmienić konfigurację, zatrzymać proces lub przejąć bezpośrednie sterowanie — żaden poziom niższy nie ogranicza tego uprawnienia |
| **2** | Always On Display (tryb operatora) | Interwencja w przebieg procesu z dowolnego miejsca, w imieniu i pod nadzorem użytkownika |
| **3** | Coordinator | Planowanie, podział pracy, sterowanie procesem, zarządzanie kolejką |
| **4** | Executor 3 / Validator | Ocena jakości wyników wykonawców; zdolność zwrotu etapu do Coordinatora (`retry`); we wcieleniu Arbitrator — rozstrzyganie konfliktów między wynikami Executora 1 i Executora 2 |
| **5** | Executor 1 / Executor 2 | Realizacja przydzielonych zadań w ramach ustalonego trybu współpracy; brak uprawnień do zmiany planu procesu |
| **6 (najniższy)** | Subagent Network | Realizacja wąsko zdefiniowanych podzadań przydzielonych przez wykonawcę macierzystego; wynik podlega agregacji i ocenie na poziomach wyższych |

### 9.2. Rozstrzyganie konfliktów (rozdz. 8.2)

Konflikt między wynikami Executora 1 i Executora 2 — możliwy w trybie pracy iteracyjnej lub naprzemiennej — rozstrzyga rola Executora 3 / Validatora we wcieleniu **Arbitrator**. Rozstrzygnięcie Arbitratora trafia do Coordinatora, który decyduje o dalszym przebiegu procesu: **kontynuacja**, **powtórzenie etapu (`retry`)** lub **eskalacja do użytkownika**.

Ścieżka rozstrzygnięcia w trzech krokach:

```
Executor 1 ⇄ Executor 2      konflikt wyników (tryb iteracyjny lub naprzemienny)
        │
        ▼  Coordinator kieruje oba wyniki do oceny
Executor 3 / Validator — wcielenie Arbitrator
        │
        ▼  rozstrzygnięcie
Coordinator ── branch ──►  kontynuacja  ·  retry  ·  eskalacja do użytkownika (poziom 1)
```

### 9.3. Konfigurowalność hierarchii (rozdz. 8.4)

Hierarchia wynika z **domyślnego podziału odpowiedzialności ról**, a nie ze sztywnej reguły platformy. Użytkownik może w każdej chwili pominąć dowolny poziom pośredni i zainterweniować bezpośrednio; skład ról uczestniczących w procesie (np. pominięcie Executora 3 / Validatora w prostszych zespołach) jest przedmiotem konfiguracji. **Żaden poziom hierarchii nie jest bramą blokującą** — jest domyślnym porządkiem, który ustępuje przed decyzją poziomu wyższego w dowolnej chwili.

Legenda przyjęta na planszy dla diagramu hierarchii:

| Znacznik | Znaczenie |
|---|---|
| Pas pełnej szerokości | Poziom nadrzędny wobec całego procesu (1, 2) |
| Węzeł z obrysem mocnym | Poziom sterujący procesem (3) |
| Węzeł z obrysem standardowym | Poziom oceniający (4) |
| Para węzłów równych | Poziom wykonawczy (5) |
| Węzły zagnieżdżone | Poziom podzadań (6) |
| Znacznik „pomijalny” | Poziom, który konfiguracja zespołu może wyłączyć (4 — rozdz. 8.4) |

---

## 10. Trzy szablony konfiguracji zespołu

Szablony **nie wprowadzają nowych ról, trybów ani akcji** — pokazują złożenie już ustalonych elementów (rozdz. 14, zdanie wprowadzające).

### 10.1. Zespół „Budowa aplikacji” (rozdz. 14.1)

Oparty na scenariuszu realizacji produktu cyfrowego w module Apps.

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie etapów, przypisanie pracy backendowi i frontendowi | — | Nie | Domyślny (dziedziczony) | Tak |
| Executor 1 | Realizacja backendu | Praca niezależna | Tak — Agent Backend, Agent API, Agent Database | Domyślny (dziedziczony) | Tak |
| Executor 2 | Realizacja frontendu (zasoby z modułu Design) | Praca niezależna | Tak — Agent UI | Domyślny (dziedziczony) | Tak |
| Executor 3 / Validator | Wcielenie Security Auditor, następnie QA Lead w kolejnych przebiegach | — | Nie | Domyślny (dziedziczony) | Tak |
| Always On Display | Operator procesu | — | — | — | — |

### 10.2. Zespół „Badanie i redakcja” (rozdz. 14.2)

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie etapów badania i redakcji raportu | — | Nie | Domyślny | Nie (proces jednorazowy) |
| Executor 1 | Zbieranie i analiza źródeł | Przekazywanie wyników → Executor 2 | Nie | Domyślny | Nie |
| Executor 2 | Redakcja raportu końcowego na podstawie ustaleń Executora 1 | Przekazywanie wyników | Nie | Domyślny | Nie |
| Executor 3 / Validator | Wcielenie Reviewer — ocena spójności i kompletności raportu | — | Nie | Domyślny | Nie |
| Always On Display | Obserwator | — | — | — | — |

### 10.3. Zespół „Pętla ciągła 24/7” (rozdz. 14.3)

Konfiguracja w pełni spięta z modułem Automations, przeznaczona do pracy bez stałego nadzoru użytkownika.

| Rola | Zadanie | Tryb współpracy | Subagent Network | Profil izolacji | Powiązanie z Automations |
|---|---|---|---|---|---|
| Coordinator | Planowanie i sterowanie każdym cyklicznym przebiegiem | — | Nie | Odrębny profil roli (dostęp sieciowy wyłączony jako świadomy wybór) | Tak — Scheduler, Orchestrator |
| Executor 1 | Wykonanie zadań pobieranych z kolejki cyklicznej | Praca iteracyjna z Executorem 2 | Tak, w miarę potrzeby zadania | Odrębny profil roli | Tak — Queue Manager |
| Executor 2 | Poprawa i uzupełnienie wyników Executora 1 w kolejnym cyklu | Praca iteracyjna | Tak, w miarę potrzeby zadania | Odrębny profil roli | Tak — Queue Manager |
| Executor 3 / Validator | Wcielenie Validator — ocena jakości przed zamknięciem każdego cyklu, konfigurowalna i pomijalna | — | Nie | Odrębny profil roli | Tak — Execution Monitor |
| Always On Display | Operator, z interwencją przez Mobile | — | — | — | — |

### 10.4. Co odróżnia szablony — odczyt porównawczy

| Wymiar | Budowa aplikacji | Badanie i redakcja | Pętla ciągła 24/7 |
|---|---|---|---|
| Tryb współpracy wykonawców | Praca niezależna | Przekazywanie wyników | Praca iteracyjna |
| Subagent Network | Tak (E1: 3 wcielenia, E2: 1 wcielenie) | Nie | Tak, w miarę potrzeby zadania |
| Wcielenie Executora 3 | Security Auditor → QA Lead | Reviewer | Validator |
| Profil izolacji | Domyślny (dziedziczony) | Domyślny | Odrębny profil roli |
| Automations | Tak | Nie (proces jednorazowy) | Tak — Scheduler, Orchestrator, Queue Manager, Execution Monitor |
| Tryb Always On Display | Operator | Obserwator | Operator, z interwencją przez Mobile |

Trzy szablony przechodzą przez **wszystkie cztery tryby współpracy poza pracą naprzemienną**, przez **trzy z siedmiu wcieleń** Executora 3 oraz przez **oba tryby Always On Display** — to najkrótsza możliwa demonstracja zakresu konfiguracyjnego środowiska.

---

## 11. Praca ciągła 24/7/365 — integracja z Automations

### 11.1. Warunek uruchomienia (rozdz. 7.1)

Integracja **nie jest domyślna**: platforma jej nie wymusza, lecz udostępnia jako możliwość skonfigurowania w dowolnym momencie — decyzją użytkownika podjętą w oknie konfiguracji oraz w ustawieniach okna modułu Automations.

### 11.2. Mapowanie mechanizmów (rozdz. 7.2)

| Mechanizm MultitaskingAI | Okno modułu Automations | Funkcja po spięciu |
|---|---|---|
| Silnik kolejek (rozdz. 4) | **Queue Manager** | Trwałe zarządzanie kolejkami zadań poza pojedynczą sesją |
| Harmonogram pracy ciągłej (sekcja Harmonogram i automatyki) | **Scheduler** | Cykliczne uruchamianie procesu według reguły czasowej |
| Orkiestracja (rozdz. 5) | **Orchestrator** | Trwałe egzekwowanie zależności między kolejnymi przebiegami procesu |
| Monitor procesu (rozdz. 8, 12) | **Execution Monitor** | Status każdego przebiegu oraz sygnalizacja błędów wymagających interwencji |

### 11.3. Schemat spięcia (rozdz. 7.3)

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
```

### 11.4. Podział odpowiedzialności w pętli ciągłej (rozdz. 7.4)

| Uczestnik pętli ciągłej | Odpowiedzialność |
|---|---|
| Coordinator | Planowanie procesu |
| Moduł Automations | Zarządzanie harmonogramami i kolejkami |
| Executor 1, Executor 2 (wraz z ich Subagent Network) | Realizacja zadań |
| Executor 3 / Validator | Kontrola jakości |
| Always On Display | Nadzór całości z perspektywy użytkownika — w trybie obserwatora lub operatora |

**Odczyt projektowy:** spięcie z Automations **nie zmienia podziału ról** — przenosi wyłącznie zarządzanie harmonogramami i kolejkami z sesji do modułu trwałego. Na planszy skutkiem przełączenia trybu ciągłego jest więc pojawienie się czwartego uczestnika (moduł Automations) w tym samym układzie, nie przebudowa układu.

### 11.5. Rola funkcji Mobile (rozdz. 7.5)

Praca w trybie 24/7/365 z natury przebiega **poza stałą obecnością użytkownika przy stanowisku roboczym**. Funkcja Mobile udostępnia kanał zatwierdzania, wstrzymywania i modyfikowania uruchomionych procesów z dowolnego miejsca — mechanizm komplementarny wobec Always On Display działającego w trybie operatora.

Rozdział ról w interwencji zdalnej (rozdz. 2.4):

| Element | Rola w interwencji |
|---|---|
| Always On Display (tryb operatora) | Obserwuje przebieg procesu i **inicjuje** interwencję (zatwierdzenie kroku, wstrzymanie procesu) |
| Mobile | Udostępnia **kanał**, przez który interwencja jest faktycznie wykonywana z dowolnego urządzenia |

### 11.6. Elementy interfejsu integracji (rozdz. 7.6)

| Element | Forma | Stany |
|---|---|---|
| Przełącznik „Powiąż z Automations” | Duży przełącznik z etykietą | Domyślny (wyłączony) · włączony · ładowanie (w trakcie spinania) |
| Plakietka statusu połączenia | Mała plakietka stanu | Niespięte (neutralny) · spięte (sukces) · błąd połączenia |
| Pole harmonogramu | Pole formularza + selektor cykliczności | Domyślny (brak harmonogramu) · skonfigurowany · ostrzeżenie (reguła sprzeczna) |
| Lista wpiętych automatyk | Lista pozycji | Domyślna · pusty stan · z pozycją aktywną |
| Przycisk zatwierdzenia zdalnego | Przycisk pełnej szerokości (kontekst mobilny) | Domyślny · ładowanie · potwierdzony |
| Powiadomienie push interwencji | Toast / powiadomienie systemowe | Wysłane · odczytane · z akcją szybkiej odpowiedzi |

**Reguła zero blokad w tym miejscu:** reguła harmonogramu wewnętrznie sprzeczna „zostaje zapisana z ostrzeżeniem zamiast blokady zapisu” (rozdz. 7.6).

### 11.7. Dwa tryby Always On Display względem procesu (rozdz. 2.2)

| Tryb | Zakres działania | Typowe zastosowanie |
|---|---|---|
| **Obserwator** | Podgląd przebiegu pętli, statusów przebiegów i hierarchii decyzji (sekcja Monitor procesu); zgłaszanie sugestii i ostrzeżeń jako proaktywne doradztwo, bez ingerencji w przebieg procesu | Nadzór nad procesem o niskim ryzyku, decyzje pozostają w gestii Coordinatora i użytkownika |
| **Operator** | Wszystkie uprawnienia obserwatora, rozszerzone o zdolność interwencji: zatwierdzanie, wstrzymywanie kroków procesu, interwencja z dowolnego miejsca | Nadzór nad procesem działającym w trybie ciągłym, gdzie brak stałej obecności użytkownika wymaga zdolności do samodzielnej interwencji |

Wybór trybu „nie jest bramką bezpieczeństwa, lecz przełącznikiem zakresu działania” (rozdz. 2.2) i pozostaje odwracalny w dowolnym momencie.

---

## 12. Agents → MultitaskingAI: rola jest ekspertem z fabryki

### 12.1. Łańcuch wartości (agents.md rozdz. 1.1)

```
MODEL BAZOWY                  surowa zdolność obliczeniowa, bez tożsamości
(inteligencja)                 wybór i kanał połączenia → Model Configuration
      │
      │   + tożsamość, instrukcje systemowe, pamięć  ─────  Agent Builder
      │   + umiejętności                              ─────  Skills Manager
      │   + pluginy, konektory, serwery MCP            ─────  Connectors Manager
      │   + zakres możliwości                          ─────  Permissions Center
      ▼
   EKSPERT (AGENT)             nazwana, skonfigurowana jednostka AI
      │
      ▼   wykorzystanie operacyjne
 TalkIn · WorkSpace · CodeStudio · MultitaskingAI (rola: Executor 1/2 · Coordinator · Executor 3/Validator)
                                                          │
                                                          ▼
                                                   ZESPÓŁ EKSPERTÓW
                                          orkiestracja, silnik kolejek,
                                       opcjonalnie integracja z Automations
                                            → pętla pracy ciągłej 24/7/365
```

### 12.2. Trzy gałęzie wykorzystania operacyjnego eksperta (agents.md rozdz. 10)

| Gałąź | Kontekst | Charakter przypisania |
|---|---|---|
| 10.1 — Wykonawca doraźny | TalkIn / WorkSpace / CodeStudio — wybór komponentu własnego w sesji | Doraźny, per zadanie |
| 10.2 — Agent Manager (moduł Workspace) | Przypisanie do PROJEKTU | Trwały, w zakresie projektu |
| 10.3 — Sekcja Role (panel orkiestracji MultitaskingAI) | Przypisanie do ROLI: Executor 1, Executor 2, Coordinator, Executor 3 / Validator | **Trwały, per rola** |

### 12.3. Co niesie przypisanie do roli (agents.md rozdz. 10.3)

| Aspekt | Ustalenie |
|---|---|
| Okno | Sekcja Role, panel orkiestracji środowiska MultitaskingAI |
| Mechanizm | Operator wskazuje dla danej roli jednego z zapisanych ekspertów **albo model bazowy bez pośrednictwa eksperta** |
| Co niesie przypisanie | Cała definicja eksperta: tożsamość, instrukcje systemowe, umiejętności, rozszerzenia, pamięć, zakres możliwości z Permissions Center |
| Nadpisanie | Definicja jest **punktem wyjścia**, nadpisywalnym ustawieniami zapisanymi na poziomie zasięgu „Rola” w oknie konfiguracji punktów izolacji |
| Kanał modelu per rola | Kanał zapisany w Model Configuration eksperta może zostać dla danej roli **nadpisany** |

**Odczyt projektowy dla planszy.** Rola w MultitaskingAI nie jest nową jednostką — jest **stanowiskiem**, na które Operator obsadza wcześniej zbudowanego eksperta. Dlatego karta roli na planszy pokazuje dwie warstwy: *stanowisko* (nazwa roli, cel, wejście, wyjście, dostępne akcje kolejki — niezmienne) oraz *obsada* (przypisany agent lub model bazowy, kanał modelu, profil izolacji — konfigurowalne, nadpisywalne). Ta dwuwarstwowość jest wprost odwzorowaniem relacji z `agents.md` rozdz. 10.3.

### 12.4. Pierwszeństwo zasięgu „Rola” (rozdz. 13.1)

```
 Pierwszeństwo poziomów zasięgu (od najniższego do najwyższego)

  globalny  →  środowisko  →  moduł  →  para modułów  →
  →  projekt  →  karta sesji  →  rola
  (warstwa bazowa)                        (pierwszeństwo najwyższe)

  brak ustawienia na poziomie  =  dziedziczenie z poziomu szerszego, nigdy blokada
```

Poziom „Rola (MultitaskingAI)” ma pierwszeństwo najwyższe spośród siedmiu poziomów zasięgu — reguła ustalona na tym poziomie nadpisuje reguły odziedziczone ze wszystkich poziomów szerszych.

---

## 13. Decyzje projektowe planszy

| # | Decyzja | Uzasadnienie | Podstawa |
|:--:|---|---|---|
| 1 | Masa wizualna węzła = 7 − poziom hierarchii decyzji | Asymetria musi nieść hierarchię, nie estetykę; wiązanie masy z udokumentowanym poziomem czyni różnicę rozmiarów sprawdzalną | KIERUNEK.md rozdz. 2 i 4; `multitaskingai.md` rozdz. 8.1 |
| 2 | Para Executor 1 / Executor 2 pozostaje symetryczna | Ten sam poziom hierarchii i „oprzyrządowanie tożsame”; zróżnicowanie byłoby asymetrią bez pokrycia | `multitaskingai.md` rozdz. 3.4, 8.1 |
| 3 | Przełącznik „pokaż masę wizualną” zamiast stałych adnotacji | Diagram ma się czytać najpierw jako układ zespołu, dopiero na żądanie jako dowód reguły; jeden ruch informacyjny, mikroprzejście 160 ms | KANON rozdz. 1 (`INTENSYWNOSC_RUCHU` 3/10) |
| 4 | Jedenaście akcji jako jedna tabela sortowalna z filtrem, nie jako jedenaście kafli | Akcje są zbiorem o jednorodnej strukturze (akcja · grupa · działanie · inicjator · reprezentacja); tabela zachowuje gęstość zwartą 8/10 | KANON rozdz. 1, rozdz. 5; `multitaskingai.md` rozdz. 4.2 |
| 5 | Symulacja cyklu życia sterowana jednym przyciskiem, z krokami `.dn-krok--*` i paskiem `.dn-postep` | Ruch niesie informację o stanie systemu — dokładnie ten przypadek opisuje wzorzec animacji KANON rozdz. 11 | KANON rozdz. 11; `multitaskingai.md` rozdz. 4.3, 12.3 |
| 6 | Mapa zależności jako SVG z podświetlaniem łańcucha przy najechaniu i fokusie | Element „Widok grafu zależności” przewiduje podświetlanie powiązanych zależności; obsługa fokusu wynika z warunku dostępności | `multitaskingai.md` rozdz. 5.4; KANON rozdz. 9 |
| 7 | Wcielenie Arbitrator bez ikony `waga` | Ikona `waga` nie istnieje w komplecie 82 plików `zasoby/ikony/svg/` (potwierdzone odczytem); KANON zakazuje ikon spoza zestawu, więc użyto ikony rodziny roli (`walidator`) z etykietą tekstową | KANON rozdz. 4; odczyt katalogu ikon |
| 8 | Stan „niegotowy” odwzorowany atrybutem `data-komunikat`, nigdy `disabled` | Dokumentacja definiuje ten stan jako „w pełni klikalny”, a KANON zakazuje `disabled` bezwzględnie | `multitaskingai.md` rozdz. 12.6; KANON rozdz. 8 |
| 9 | Przełącznik trybu ciągłego dodaje uczestnika (moduł Automations), nie przebudowuje układu | Rozdz. 7.4 utrzymuje ten sam podział ról po spięciu — przenosi wyłącznie zarządzanie harmonogramami i kolejkami | `multitaskingai.md` rozdz. 7.4 |
| 10 | Karta roli w dwóch warstwach: stanowisko (niezmienne) i obsada (konfigurowalna) | Rola jest stanowiskiem obsadzanym ekspertem z modułu Agents; definicja eksperta jest „punktem wyjścia, nadpisywalnym” | `agents.md` rozdz. 10.3 |
| 11 | Szablony zespołu rozwijane w miejscu, bez modala | Modal jest w tym środowisku zarezerwowany dla opcjonalnego potwierdzenia akcji nieodwracalnej i domyślnie wyłączony | `multitaskingai.md` Załącznik A.2, wiersz „Modal” |
| 12 | Dane przykładowe wyłącznie z makiet dokumentu i oznaczone jako przykładowe | Zakaz zmyślonych metryk i nazw; każde zadanie, kolejka i podagent ma numer rozdziału źródła | KANON rozdz. 10 pkt 3–4 |

---

## 14. Źródła

| Zakres | Plik | Rozdziały |
|---|---|---|
| Kontrakt kierunku, żetony, komponenty, zasady redakcyjne | `WYNIK/KANON.md` | 1, 2, 4, 5, 6, 7.6, 8, 9, 10, 11 |
| Role zespołu — cel, wejście, wyjście, oprzyrządowanie | `dok/projekt-ui/srodowiska/multitaskingai.md` | 3.0–3.6 |
| Silnik kolejek — zasięgi, akcje, cykl życia, elementy sekcji Kolejki | `dok/projekt-ui/srodowiska/multitaskingai.md` | 4.1–4.6 |
| Orkiestracja i zależności, mapa zależności | `dok/projekt-ui/srodowiska/multitaskingai.md` | 5.1–5.5 |
| Panel orkiestracji — sześć sekcji, elementy, konfigurowalność | `dok/projekt-ui/srodowiska/multitaskingai.md` | 6.1–6.5 |
| Integracja z Automations — praca ciągła 24/7/365 | `dok/projekt-ui/srodowiska/multitaskingai.md` | 7.1–7.6 |
| Hierarchia decyzji i rozstrzyganie konfliktów | `dok/projekt-ui/srodowiska/multitaskingai.md` | 8.1–8.4 |
| Okna robocze ról — makiety tekstowe (dane przykładowe) | `dok/projekt-ui/srodowiska/multitaskingai.md` | 9.1–9.7 |
| Katalog elementów interfejsu | `dok/projekt-ui/srodowiska/multitaskingai.md` | 10.1–10.9 |
| Stany procesu, roli, zadania; konwencja stanów interfejsu | `dok/projekt-ui/srodowiska/multitaskingai.md` | 12.1–12.6 |
| Izolacja na poziomie roli — pierwszeństwo zasięgu | `dok/projekt-ui/srodowiska/multitaskingai.md` | 13.1–13.5 |
| Szablony konfiguracji zespołu | `dok/projekt-ui/srodowiska/multitaskingai.md` | 14.1–14.4 |
| Scenariusze operacyjne | `dok/projekt-ui/srodowiska/multitaskingai.md` | 15.1–15.4 |
| Ikony i komponenty środowiska | `dok/projekt-ui/srodowiska/multitaskingai.md` | Załącznik A.1, A.2 |
| Warstwa centralna Always On Display, tryby, Mobile | `dok/projekt-ui/srodowiska/multitaskingai.md` | 2.1–2.6 |
| Agents jako fabryka ekspertów, łańcuch wartości | `dok/projekt-ui/okna/agents.md` | 1.1–1.3 |
| Od eksperta do wykonawcy — przypisanie do roli | `dok/projekt-ui/okna/agents.md` | 10, 10.1–10.3 |
| Asymetria koordynator–wykonawca, anty-domyślne | `design/opracowania/design/01-kierunek/KIERUNEK.md` | 2, 4, 5 |
| Klasy komponentów `.dn-*` | `WYNIK/zasoby/css/komponenty.css` | odczyt pliku |
| Klasy warstwy prototypu `.pt-*` i atrybuty `data-*` | `WYNIK/zasoby/prototyp.css`, `WYNIK/zasoby/prototyp.js` | odczyt plików |
| Ikony (82 pliki) i ich zastosowania | `WYNIK/zasoby/ikony/svg/*.svg`, `WYNIK/zasoby/ikony/manifest.json` | odczyt katalogu i manifestu |
| Emblemat środowiska | `WYNIK/zasoby/marka/srodowiska/srodowisko-multitaskingai.svg` | odczyt pliku |
| Prototypy okien środowiska | `WYNIK/05-okna/srodowiska/multitaskingai.html`, `mtai-okna-rol.html`, `mtai-izolacja-i-zespoly.html` | odczyt katalogu |
| Prototyp modułu Automations i Agents | `WYNIK/05-okna/moduly/automations.html`, `WYNIK/05-okna/moduly/agents.html` | odczyt katalogu |
| Prototypy platformowe (AOD, Mobile) | `WYNIK/05-okna/platformowe/always-on-display.html`, `mobile.html` | odczyt katalogu |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
