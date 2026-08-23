# Danaco Console — Moduł Automations — dokument projektowy

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
| **Moduł** | Automations (rozdz. 11.3 Koncepcji platformy; specyfikacje/specyfikacja-modulow.md 4.3) |
| **Środowiska dostępności** | Brak okna modułowego w jakimkolwiek środowisku — wyłącznie strona główna, strefa 2 |
| **Forma udostępnienia** | Komponent własny (Automatyka); gotowe automatyki trafiają do sesji modułów docelowych jako komponenty własne |
| **Odbiorcy dokumentu** | Designer (co, gdzie, w jakiej formie, do czego), Deweloper (co zbudować) |
| **Zasada nadrzędna** | Zero blokad; klucze i zależności jawne; wszystko sterowane z okna konfiguracji; automatyka jako komponent własny |
| **Źródło faktów** | architektura/koncepcja-platformy.md (rozdz. 5.2, 8, 9.3, 11.3–11.6), specyfikacje/specyfikacja-modulow.md (4.3), specyfikacje/specyfikacja-okien-operacyjnych.md (6.3, 7.3), interfejs-uzytkownika/strona-glowna-i-nawigacja.md (3.3, 5.1), architektura/model-konfiguracji.md (5.2, 5.3, 5.4, 5.8), interfejs-uzytkownika/system-wizualny.md |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Komplet okien operacyjnych modułu — przegląd](#2-komplet-okien-operacyjnych-modułu--przegląd)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
4. [Katalog funkcji i narzędzi](#4-katalog-funkcji-i-narzędzi)
5. [Katalog elementów interfejsu](#5-katalog-elementów-interfejsu)
6. [Przepływy pracy w module](#6-przepływy-pracy-w-module)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
8. [Stany, dane i powiązania z innymi modułami](#8-stany-dane-i-powiązania-z-innymi-modułami)
9. [Scenariusze użycia](#9-scenariusze-użycia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Rola modułu

**Automations jest silnikiem procesów platformy** — miejscem, w którym powtarzalna, proceduralna praca (wyzwalana zdarzeniem lub harmonogramem, nie rozmową) jest projektowana jako przepływ, spinana w potok wielu modułów i modeli, kolejkowana, orkiestrowana według zależności i nadzorowana przez cały cykl życia przebiegu. Moduł pokrywa całą warstwę spinającą narzędzia platformy: integrację aplikacji, orkiestrację przepływów, silnik procesów, planowanie zadań, kolejkowanie, dostarczanie i weryfikację wywołań przychodzących, robotyzację czynności powtarzalnych oraz lekkie przetwarzanie danych w toku przepływu.

### 1.2. Dla kogo

Automations jest przestrzenią dla użytkowników przekształcających **powtarzalne czynności** — cykliczne raporty, regularne przetwarzanie danych, zaplanowane wywołania modeli — w procesy działające **bez stałego nadzoru**. Adresatem są zarówno pojedynczy użytkownicy automatyzujący własną pracę, jak i zespoły budujące operacyjne zaplecze procesowe dla całych projektów prowadzonych w innych modułach platformy.

### 1.3. Po co

Moduł oddziela pracę **konwersacyjną, jednorazową** (typową dla pozostałych czternastu modułów) od pracy **proceduralnej, uruchamianej według harmonogramu lub zdarzenia**. Automatyka raz zbudowana działa samodzielnie — bez konieczności, by użytkownik za każdym razem inicjował ją rozmową — i zostaje wpięta jako gotowy komponent własny w sesję dowolnego modułu docelowego.

### 1.4. Zakres tematyczny modułu

| Obszar | Zakres |
|---|---|
| Wyzwalanie | Harmonogramy czasowe (cron, kalendarz, interwał), wyzwalacze zdarzeniowe (webhook, zmiana pliku, wynik modelu, zdarzenie modułu lub środowiska), wyzwalacze ręczne, wyzwalacze łańcuchowe (automatyka wyzwala automatykę) |
| Projektowanie przepływu | Kanwa wizualna węzłów, przepływ warunkowy, pętle, rozgałęzienia, podprzepływy, mapowanie danych między krokami, obsługa błędów jako część przepływu |
| Spinanie modułów i modeli | Krok „akcja w module” (dowolny z 15 modułów), krok „wywołanie modelu”, potoki wielomodelowe (kilka modeli w jednym przebiegu), pętle iteracyjne model↔model |
| Kolejkowanie | Pełny silnik kolejek (jedenaście akcji), priorytety, zasięgi kolejki, współbieżność, ograniczanie przepustowości, kolejka zadań martwych, ponawianie z wycofaniem |
| Orkiestracja | Graf zależności (DAG), grupy równoległe i sekwencyjne, ścieżka krytyczna, bramki dołączenia (fan-in/fan-out), maszyna stanów przebiegu |
| Nadzór i niezawodność | Monitor przebiegów na żywo, logi, metryki, alarmy, ponowienia, idempotencja, punkty wznowienia, transakcje kompensujące (saga) |
| Utrzymanie | Wersjonowanie definicji, testy i symulacje, import/eksport, biblioteka szablonów, sekrety i poświadczenia jawnie zarządzane, audyt |

### 1.5. Granice modułu

| Poza zakresem | Gdzie realizowane | Dlaczego |
|---|---|---|
| Praca konwersacyjna, jednorazowa | Pozostałe 14 modułów | Automations obsługuje procedurę, nie dialog |
| Trwałe przechowywanie plików i wiedzy | Library / Project Library | Automatyka *przenosi* dane, nie jest ich repozytorium |
| Definiowanie tożsamości i uprawnień agenta | Agents (Agent Builder, Permissions Center) | Automations *używa* agenta jako wykonawcy kroku, nie tworzy go |
| Warstwa centralna, role wykonawcze, orkiestracja ról | Środowisko MultitaskingAI (rozdz. 13 Koncepcji) | Automations *dostarcza* zaplecze harmonogramów i kolejek; nie zastępuje warstwy ról |
| Edycja kodu, konsole, budowanie | Developer / Terminal | Automatyka *wyzwala* te operacje jako kroki, nie jest środowiskiem programistycznym |
| Wizualizacja biznesowa danych wynikowych | Research / Studio / Workspace Dashboard | Automatyka *generuje i dostarcza* dane, prezentację wykonuje moduł docelowy |

Granica jest prosta: **Automations jest czasownikiem, nie rzeczownikiem** — spina, wyzwala, kolejkuje i nadzoruje, a przedmiotem pracy pozostają zasoby żyjące w innych modułach.

### 1.6. Charakter pracy — cecha wyróżniająca

| Cecha | Opis |
|---|---|
| Brak okna modułowego | Automations, jako jedyny z piętnastu modułów, nie ma pozycji w bocznej nawigacji żadnego środowiska — cała praca projektowa nad automatyką odbywa się z poziomu strony głównej |
| Jednostka wytwórcza | Automatyka — nazwany komponent własny, zapisywany w pamięci aplikacji po zbudowaniu |
| Cykl pracy | Budowa (strona główna) → zapis jako komponent własny → wpięcie do sesji modułu docelowego → wykonanie samodzielne wg harmonogramu lub zdarzenia |
| Nadzór | Monitorowanie przebiegów odbywa się z tego samego miejsca budowy (Execution Loop Window i Execution Monitor) oraz — dla procesów długotrwałych — z poziomu funkcji globalnej Mobile |

### 1.7. Miejsce modułu w architekturze platformy

```
STRONA GŁÓWNA — Centrum dowodzenia
├─ Strefa 1 · środowiska            (Automations NIE występuje tutaj)
├─ Strefa 2 · komponenty własne
│    kafel „Automations” → „Zbuduj automatykę”
│         │
│         ▼
│    OKNA OPERACYJNE MODUŁU AUTOMATIONS
│    Chat Window · Execution Loop Window · Workflow Builder ·
│    Scheduler · Queue Manager · Orchestrator · Execution Monitor
│         │
│         ▼
│    AUTOMATYKA zapisana jako komponent własny
│         │
│         ▼
└─ Strefa 3 · ustawienia    ▼
                    wpięta w sesji modułu docelowego
                    (dowolny z piętnastu modułów platformy)
                            │
                            ▼
            powiązanie ze środowiskiem MultitaskingAI
            (harmonogramy i kolejki pętli pracy ciągłej
             — rozdz. 8.3 niżej)
```

Brak okna w bocznej nawigacji jest architektoniczną decyzją, nie brakiem — odzwierciedla naturę automatyki jako wytworu konfigurowanego raz, a wykorzystywanego wielokrotnie, w wielu miejscach platformy jednocześnie.

---

## 2. Komplet okien operacyjnych modułu — przegląd

| Okno | Rola w module | Forma wiodąca i waga wizualna | Warstwa | Sposób wywołania | Punkt wejścia? |
|---|---|---|---|---|---|
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Nie — stale obecne |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; orkiestracja i nadzór uruchomień przepływów pracy | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Widoczne przy aktywnym zleceniu; poza tym otwierane znacznikiem `Pętla wykonawcza` w pasku kontekstu | Nie — pierwszoplanowe |
| Workflow Builder | Projektowanie procesu automatycznego | Prawa kolumna dominująca, kanwa robocza | 1 | Widoczne bez interakcji | **Tak** |
| Scheduler | Ustalanie harmonogramu i cykliczności | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Harmonogram ▼` | Nie |
| Queue Manager | Zarządzanie kolejką zadań procesu | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Kolejka ▼` | Nie |
| Orchestrator | Zależności między etapami procesu | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Zależności ▼` | Nie |
| Execution Monitor | Obserwacja statusu uruchomień i błędów | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Przebiegi ▼` | Nie |

```
 Makieta całościowa — Moduł: Automations · stan spoczynku
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Workflow Builder          │ Rozszerzenie
  nawigacja │ Użytkownik ↔         │ — kanwa procesu           │ boczne
  modułów   │ Wykonawca            │                           │ (zwinięte)
            │                      │  ┌────────┐  ┌────────┐   │
            │ ──────────────────   │  │ Start  │─►│ Pobierz│   │      ⋮
            │ Execution Loop       │  └────────┘  └────────┘   │
            │ Koordynator ↔        │                           │
            │ Wykonawca            │                           │
            │                      │                           │
  ═════════════════════════════════════════════════════════════════════════
  [Automations] [Raport tygodniowy sprzedaży] [Harmonogram ▼] [Kolejka ▼]
  [Zależności ▼] [Przebiegi ▼] [Fable 5] [Ultra]                        ☰
 ═══════════════════════════════════════════════════════════════════════════
```

### 2.1. Warstwy widoczności w module

Interfejs modułu ujawnia funkcje stopniowo: jeżeli funkcja nie jest potrzebna do realizacji aktualnej czynności, pozostaje niewidoczna. W stanie spoczynku widoczne są wyłącznie elementy warstwy 1 — Chat Window, Execution Loop Window, kanwa Workflow Buildera, pasek kontekstu ze znacznikami i wskaźniki stanu wykonania.

| Warstwa | Zakres w module Automations | Sposób wywołania |
|---|---|---|
| 1 | Chat Window, Execution Loop Window, kanwa Workflow Buildera, pasek kontekstu automatyki, wskaźnik stanu bieżącego przebiegu | Widoczne bez interakcji |
| 2 | Scheduler, Queue Manager, Orchestrator, Execution Monitor, wybór modelu kroku, wybór wykonawcy kroku, zasięg kolejki, poziom wysiłku modelu | Znacznik kontekstowy, ikona, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Zestawy akcji na węźle, akcje na zadaniu kolejki, ustawienia szybkie harmonogramu, warianty ponowienia przebiegu | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Skarbiec poświadczeń, dziennik audytu, izolacja procesu automatyki, uruchomienie wsteczne, odtworzenie przebiegu z ładunku, edycja definicji w formacie strukturalnym | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny |

Każda ukryta funkcja pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego.

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (Użytkownik ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Rola | Główne okno komunikacji między Użytkownikiem a Wykonawcą; centralny punkt pracy w module i podstawowy mechanizm sterowania wszystkimi procesami automatyki |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Zawartość | Historia rozmowy właściwa bieżącej automatyce; strumień odpowiedzi na żywo; zatwierdzanie i przerywanie działań |
| Warstwa | 1 |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Generowanie kroku z opisu | Polecenie w języku naturalnym przekształcane przez Wykonawcę w gotowy krok Workflow Buildera; stan „sugestia oczekująca”, zastosowanie jawne | 1 | Wpisanie polecenia |
| Generowanie całego przepływu z opisu | Opis celu procesu przekształcany w kompletny szkic przepływu do dalszej edycji | 1 | Wpisanie polecenia |
| Wyjaśnienie procesu | Zapytanie o działanie istniejącego przepływu — Wykonawca opisuje przebieg krok po kroku na podstawie definicji | 1 | Wpisanie polecenia |
| Debugowanie konwersacyjne | Przekazanie logu błędu z Execution Monitor i uzyskanie sugerowanej poprawki kroku | 1 | Wpisanie polecenia lub akcja „Omów w Chat Window” z logu |
| Zatwierdzanie i przerywanie działań | Akceptacja lub zatrzymanie operacji zleconej Wykonawcy | 1 | Przycisk w strumieniu odpowiedzi |
| Wzmianki (`@`) | Odwołanie do konkretnego kroku, kolejki lub uruchomienia wprost w treści polecenia | 2 | Znak `@` w polu polecenia |
| Szablony poleceń | Zapisane polecenia typowe dla budowy automatyk (np. „dodaj krok warunkowy po kroku 3”) | 2 | Element zwinięty `Szablon ▼` |
| Akcje na wiadomości | Kopiuj, zastosuj sugestię wprost do Workflow Buildera, odrzuć | 3 | Menu kebab (⋮) przy wiadomości |
| Wskaźnik kontekstu | Adnotacja, której automatyki dotyczy bieżąca rozmowa | 1 | Znacznik w pasku kontekstu |

**Zachowanie i stany:** okno rekonfigurowane do kontekstu budowanej lub przeglądanej automatyki; historia domyślnie odrębna per automatyka, współdzielenie sterowane z okna konfiguracji. Stan „sugestia oczekująca” — sugerowana zmiana kroku wymaga jawnego zastosowania przez Użytkownika.

```
 Makieta — Chat Window (stan spoczynku)
 ═══════════════════════════════════════════
  Chat Window — Użytkownik ↔ Wykonawca     ⋮
  [Raport tygodniowy sprzedaży]
 ───────────────────────────────────────────
  Historia rozmowy
  ▸ Użytkownik: „Dodaj krok wysyłający
    raport mailem po kroku 4”
  ▸ Wykonawca: sugestia kroku
    [ Zastosuj ] [ Odrzuć ]
 ───────────────────────────────────────────
  ┌───────────────────────────────────────┐
  │ Polecenie…                            │
  └───────────────────────────────────────┘
  [ Szablon ▼ ]              [ Wyślij ⏎ ]
 ═══════════════════════════════════════════
```

### 3.2. Execution Loop Window (Koordynator ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Rola | Okno pętli wykonawczej; mechanizm orkiestracji i nadzoru uruchomień przepływów pracy. Koordynator dekomponuje zlecenie na zadania kroków, przydziela je Wykonawcom, nadzoruje realizację i decyduje o ponowieniach |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli, sterowanie przebiegiem |
| Warstwa | 1 |

W module Automations pętla wykonawcza Koordynator ↔ Wykonawca jest mechanizmem centralnym: każde uruchomienie przepływu pracy — wyzwolone harmonogramem, zdarzeniem lub ręcznie — przebiega jako pętla, w której Koordynator kolejno wydaje zadania kroków, odbiera wyniki, ocenia je względem warunków przepływu i steruje dalszym przebiegiem. Okno prezentuje ten przebieg w czasie rzeczywistym i udostępnia sterowanie nim.

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Widok zlecenia i dekompozycji | Bieżące zlecenie uruchomienia przepływu rozłożone na zadania kroków | 1 | Widoczne bez interakcji |
| Kolejka i stan zadań pętli | Lista zadań pętli ze stanem: zaplanowane, w toku, zakończone, błędne | 1 | Widoczne bez interakcji |
| Strumień komunikatów sterujących | Wymiana poleceń i potwierdzeń między Koordynatorem a Wykonawcą | 1 | Widoczne bez interakcji |
| Wskaźniki przebiegu pętli | Numer iteracji, czas trwania, liczba ponowień, postęp względem grafu zależności | 1 | Widoczne bez interakcji |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia w toku | 1 | Przyciski w nagłówku okna |
| Wyniki kontroli jakości | Ocena wyniku kroku względem warunku przepływu i decyzja o ponowieniu lub przejściu dalej | 2 | Rozwinięcie wiersza zadania |
| Wybór wykonawcy zadania | Przypisanie modelu lub agenta do zadania pętli | 2 | Element zwinięty `Wykonawca ▼` |
| Punkty wznowienia pętli | Wznowienie od zapisanego punktu bez powtarzania ukończonych kroków | 3 | Menu kebab (⋮) przy zleceniu |
| Podgląd i odtworzenie ładunku | Wgląd w dane wejściowe i wyjściowe zadania oraz odtworzenie przebiegu z tym ładunkiem | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny |

**Zachowanie i stany:** okno aktualizuje się na żywo kanałem WebSocket. Stany pętli: bezczynna, w toku, wstrzymana, ponawiająca zadanie, zakończona sukcesem, zakończona błędem, przerwana.

```
 Makieta — Execution Loop Window (stan spoczynku)
 ═══════════════════════════════════════════════
  Execution Loop — Koordynator ↔ Wykonawca     ⋮
  Zlecenie: uruchomienie „Raport tygodniowy”
  [ ⏸ Wstrzymaj ] [ ⏹ Przerwij ]
 ───────────────────────────────────────────────
  Zadania pętli
  1 Pobierz dane          ✔ zakończone
  2 Waliduj wynik         ● w toku
  3 Wygeneruj dokument    ○ zaplanowane
  4 Wyślij raport         ○ zaplanowane
 ───────────────────────────────────────────────
  Komunikaty sterujące
  ▸ Koordynator → Wykonawca: zadanie 2
  ▸ Wykonawca → Koordynator: wynik cząstkowy
 ───────────────────────────────────────────────
  Iteracja 1 · 00:42 · ponowienia 0
  [ Wykonawca ▼ ]
 ═══════════════════════════════════════════════
```

### 3.3. Workflow Builder (punkt wejścia)

| Pole | Treść |
|---|---|
| Rola | Projektowanie procesu automatycznego na kanwie węzłowej |
| Waga wizualna | Prawa kolumna dominująca, kanwa robocza |
| Zawartość | Kroki procesu, warunki, kolejność wykonania, zmienne przepływu |
| Warstwa | 1 |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kanwa wizualna | Przestrzeń robocza z przeciąganymi węzłami kroków, przewijalna i skalowalna (zoom, pan), minimapa, autoukład | 1 | Widoczna bez interakcji |
| Paleta typów kroku | Zestaw węzłów: wywołanie modelu, akcja w module, HTTP/API, warunek, pętla, opóźnienie, podprzepływ, transformacja danych, punkt kontrolny wymagający potwierdzenia | 2 | Element zwinięty `Krok ▼` w pasku kontekstu |
| Łączniki między krokami | Rysowanie i edycja połączeń określających kolejność wykonania | 1 | Przeciągnięcie od portu węzła |
| Panel zmiennych | Definiowanie zmiennych procesu przekazywanych między krokami | 2 | Znacznik `Zmienne ▼` |
| Mapowanie danych | Wizualne łączenie wyjść jednego kroku z wejściami kolejnego, wyrażenia i szablony na polach | 2 | Kliknięcie łącznika |
| Test i symulacja | Uruchomienie procesu w trybie testowym, bez wpięcia produkcyjnego, z podglądem wyniku każdego kroku i danymi próbnymi na wejściach | 2 | Przycisk `Test` w nagłówku |
| Walidacja procesu | Wykrycie kroków bez połączenia, cykli i brakujących parametrów — ostrzeżenie sygnalizowane na kanwie, zapis pozostaje możliwy | 1 | Automatyczna, sygnalizacja na węźle |
| Wersjonowanie i porównanie | Historia kolejnych wersji procesu, porównanie strukturalne dwóch wersji, przywrócenie wcześniejszej | 3 | Menu `Wersje ⋮` |
| Duplikowanie i szablonowanie | Kopia procesu jako punkt wyjścia; zapis dowolnego przepływu jako szablonu z parametrami | 3 | Menu kebab (⋮) nagłówka |
| Notatki przy kroku | Dokumentacja opisowa dołączana do pojedynczego kroku, widoczna na kanwie jako rozwijalna adnotacja | 3 | Menu kontekstowe węzła |
| Cofnij / ponów | Pełna historia edycji kanwy w ramach bieżącej sesji pracy | 2 | Skrót klawiszowy, ikona |
| Tagowanie i wyszukiwanie | Kategoryzacja automatyk oraz wyszukiwanie po nazwie, tagu, module docelowym lub wyzwalaczu | 3 | Menu hamburger (☰) |
| Import i eksport definicji | Zapis i wczytanie przepływu w formacie strukturalnym (JSON/YAML) | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** zmiany zapisywane jako definicja automatyki — komponentu własnego. Stan „niezapisane zmiany”, stan „ostrzeżenia walidacji” (żółte oznaczenie problematycznego kroku — sygnalizacja, zapis pozostaje możliwy), stan „test w toku” (animacja przebiegu na kanwie), rozdzielenie wersji roboczej od opublikowanej.

```
 Makieta — Workflow Builder (stan spoczynku)
 ═══════════════════════════════════════════════════════════
  Automatyka: „Raport tygodniowy sprzedaży”   [Test][Zapisz] ⋮
 ───────────────────────────────────────────────────────────
  ┌─────────┐    ┌─────────┐    ┌───────────┐    ┌─────────┐
  │ Start   │───►│ Pobierz │───►│ Warunek:  │───►│ Wyślij  │
  │ (event) │    │ dane    │    │ dane > 0? │ Tak│ raport  │
  └─────────┘    └─────────┘    └─────┬─────┘    └─────────┘
                                      │ Nie
                                      ▼
                                ┌───────────┐
                                │ Powiadom  │
                                │ o braku   │
                                └───────────┘
 ───────────────────────────────────────────────────────────
  [ Krok ▼ ] [ Zmienne ▼ ] [ Wersje ⋮ ]                    ☰
 ═══════════════════════════════════════════════════════════
```

### 3.4. Scheduler

| Pole | Treść |
|---|---|
| Rola | Ustalanie momentu i zdarzenia startu procesu |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zawartość | Lista wyzwalaczy automatyki, kalendarz uruchomień, podgląd najbliższych terminów |
| Warstwa | 2 — wywołanie znacznikiem `Harmonogram ▼` |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kreator wyrażenia cyklicznego | Budowa reguły cron w formie przyjaznej, z podglądem wyniku w składni technicznej i listą najbliższych uruchomień | 2 | Widoczny po otwarciu okna |
| Wzorce cykliczności | Gotowe wzorce: co minutę, godzinnie, codziennie, co tydzień, co miesiąc, kwartalnie, niestandardowo | 2 | Lista rozwijana |
| Harmonogram interwałowy | Uruchomienie co N sekund, minut lub godzin niezależnie od zegara ściennego | 2 | Zakładka trybu reguły |
| Widok kalendarza uruchomień | Miesięczny i tygodniowy podgląd zaplanowanych oraz historycznych przebiegów | 2 | Przełącznik widoku |
| Strefy czasowe i zmiana czasu | Reguła harmonogramu w wybranej strefie, poprawna obsługa przejścia czasu letniego i zimowego | 2 | Element zwinięty `Strefa ▼` |
| Okna wykonania | Ograniczenie przedziałów czasu, w których uruchomienie następuje (np. wyłącznie godziny robocze) | 3 | Menu kebab (⋮) reguły |
| Wyzwalacz webhook | Unikatowy adres wejściowy automatyki; przyjęcie ładunku HTTP z weryfikacją podpisu (HMAC) i deduplikacją | 2 | Pozycja listy `+ Dodaj wyzwalacz` |
| Wyzwalacz zmiany pliku | Uruchomienie na dodanie, zmianę lub usunięcie pliku w Library albo Project Library | 2 | Pozycja listy `+ Dodaj wyzwalacz` |
| Wyzwalacz na wyniku modelu | Uruchomienie, gdy wynik wywołania modelu spełnia warunek (np. klasyfikacja „pilne”) | 2 | Pozycja listy `+ Dodaj wyzwalacz` |
| Wyzwalacz zdarzenia modułu | Reakcja na zdarzenie dowolnego modułu (zatwierdzenie zmian w Developer, wdrożenie w Apps, nowy plik w Library) | 2 | Pozycja listy `+ Dodaj wyzwalacz` |
| Wyzwalacze łańcuchowe | Zakończenie jednej automatyki — sukcesem lub błędem — wyzwala kolejną | 2 | Pozycja listy `+ Dodaj wyzwalacz` |
| Wiele wyzwalaczy na proces | Powiązanie jednego przepływu z kilkoma niezależnymi regułami czasu i zdarzeń jednocześnie | 2 | Lista wyzwalaczy |
| Wstrzymanie i wznowienie harmonogramu | Przełącznik zatrzymujący cykliczne uruchamianie bez usuwania definicji | 2 | Suwak w nagłówku |
| Uruchomienie wsteczne | Wykonanie przebiegów dla przeszłych, pominiętych terminów w zadanym zakresie dat | 4 | Polecenie języka naturalnego, tryb administracyjny |
| Historia wyzwoleń | Zapis rzeczywistych momentów wyzwolenia z przyczyną: harmonogram, zdarzenie, uruchomienie ręczne | 3 | Menu kebab (⋮) |
| Nadzór obecności uruchomień | Alarm, gdy oczekiwane uruchomienie nie nastąpiło w oknie tolerancji | 3 | Panel popover reguły alarmowej |

**Zachowanie i stany:** harmonogram obowiązuje po powiązaniu z przepływem. Stan „aktywny” i „wstrzymany” sygnalizowany przełącznikiem i plakietką. Stan „najbliższe uruchomienie za…” widoczny jako licznik.

```
 Makieta — Scheduler (rozszerzenie boczne, stan spoczynku)
 ═════════════════════════════════════════════
  Harmonogram — „Raport tygodniowy sprzedaży”
  [●aktywny]                                 ⋮
 ─────────────────────────────────────────────
  Wzorzec: ( ) codziennie (•) co tydzień
           ( ) co miesiąc ( ) własny
  Dzień: [ poniedziałek ▼ ]
  Godzina: [ 07:00 ▼ ]  Strefa: [ CET ▼ ]
  Wyrażenie cron:  0 7 * * 1
 ─────────────────────────────────────────────
  Najbliższe uruchomienia:
  2026-08-10 · 2026-08-17 · 2026-08-24
 ─────────────────────────────────────────────
  Wyzwalacze:  [ + Dodaj wyzwalacz ]
 ═════════════════════════════════════════════
```

### 3.5. Queue Manager

| Pole | Treść |
|---|---|
| Rola | Zarządzanie kolejką zadań automatyki na żywo |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zawartość | Lista zadań oczekujących i przetwarzanych, wskaźnik głębokości kolejki |
| Warstwa | 2 — wywołanie znacznikiem `Kolejka ▼` |

**Pełny zestaw akcji silnika kolejek:**

| Akcja | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| `enqueue` | Dodanie nowego zadania do kolejki | 3 | Menu kebab (⋮) wiersza zadania |
| `dequeue` | Zdjęcie zadania z kolejki przed wykonaniem | 3 | Menu kebab (⋮) wiersza zadania |
| `delay` | Opóźnienie wykonania zadania o wskazany czas | 3 | Menu kebab (⋮) wiersza zadania |
| `retry` | Ponowna próba wykonania zadania zakończonego błędem | 3 | Menu kebab (⋮) wiersza zadania |
| `pause` | Wstrzymanie przetwarzania kolejki | 2 | Przycisk w nagłówku okna |
| `resume` | Wznowienie przetwarzania wstrzymanej kolejki | 2 | Przycisk w nagłówku okna |
| `split` | Podział zadania na mniejsze podzadania | 3 | Menu kebab (⋮) wiersza zadania |
| `merge` | Scalenie kilku zadań w jedno | 3 | Menu kebab (⋮) zaznaczenia |
| `route` | Skierowanie zadania do innej kolejki lub innego wykonawcy | 3 | Menu kebab (⋮) wiersza zadania |
| `branch` | Rozgałęzienie przetwarzania zadania na kilka torów równoległych | 3 | Menu kebab (⋮) wiersza zadania |
| `condition` | Warunkowe przetworzenie zadania w zależności od wartości parametru | 3 | Menu kebab (⋮) wiersza zadania |

**Dodatkowe narzędzia interfejsu:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista zadań z kolumnami statusu | Oczekujące, w toku, zakończone, błędne, z przeciąganiem między priorytetami | 1 | Widoczna po otwarciu okna |
| Zmiana priorytetu | Przeciągnięcie zadania wyżej lub niżej w kolejce | 2 | Przeciągnięcie wiersza |
| Zasięg kolejki | Globalna, lokalna, dla modelu, dla agenta, dla projektu | 2 | Element zwinięty `Zasięg ▼` |
| Współbieżność i ograniczanie przepustowości | Limit równoległych zadań na kolejkę i zasięg oraz ograniczenie tempa przetwarzania | 3 | Panel popover ustawień kolejki |
| Ponawianie z wycofaniem | Automatyczne ponowienie z wykładniczym wycofaniem i rozproszeniem, z limitem prób | 3 | Panel popover ustawień kolejki |
| Kolejka zadań martwych | Zadania trwale nieudane trafiają do osobnej kolejki do przeglądu ręcznego lub ponowienia zbiorczego | 3 | Zakładka `Zadania martwe` |
| Idempotencja zadań | Klucz idempotencji zapobiegający podwójnemu wykonaniu przy ponowieniu lub duplikacie wywołania przychodzącego | 4 | Tryb administracyjny, wyszukiwarka funkcji |
| Wizualizacja głębokości kolejki | Wykres liczby zadań oczekujących w czasie, obciążenie per zasięg | 1 | Widoczna po otwarciu okna |
| Szuflada szczegółów zadania | Ładunek zadania, znaczniki czasu, dotychczasowe próby, powiązane uruchomienie, logi | 3 | Kliknięcie wiersza zadania |
| Filtrowanie i wyszukiwanie | Po statusie, dacie, przepływie źródłowym, zasięgu | 2 | Ikona filtra |
| Akcje zbiorcze | Zastosowanie dowolnej z jedenastu akcji do wielu zaznaczonych zadań naraz | 3 | Menu kebab (⋮) zaznaczenia |
| Ręczne wstrzyknięcie zadania | „Uruchom teraz” — dodanie zadania do kolejki poza harmonogramem | 2 | Przycisk w nagłówku okna |
| Harmonogram opóźnień | Zadania zaplanowane na wskazany termin i cykliczne wewnątrz kolejki | 3 | Panel popover ustawień kolejki |

**Zachowanie i stany:** kolejka przetwarzana zgodnie z regułami silnika kolejek; aktualizacja na żywo kanałem WebSocket. Stan „kolejka wstrzymana” sygnalizowany plakietką ostrzegawczą na nagłówku okna.

```
 Makieta — Queue Manager (rozszerzenie boczne, stan spoczynku)
 ═══════════════════════════════════════════════
  Kolejka — „Raport tygodniowy sprzedaży”
  [ Zasięg ▼ ]  [ ⏸ pauza ]                    ⋮
 ───────────────────────────────────────────────
  # 1  Pobierz dane        [w toku]           ⋮
  # 2  Waliduj wynik       [oczekuje]         ⋮
  # 3  Wyślij raport       [błąd ✕]           ⋮
 ───────────────────────────────────────────────
  Głębokość kolejki (7 dni):  ▂▃▅▂▇▃▂
 ═══════════════════════════════════════════════
```

### 3.6. Orchestrator

| Pole | Treść |
|---|---|
| Rola | Definiowanie zależności między krokami procesu jako grafu |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zawartość | Zależności i kolejność wykonania kroków przepływu, widok zorientowany na relacje |
| Warstwa | 2 — wywołanie znacznikiem `Zależności ▼` |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Graf zależności (DAG) | Wizualny widok kroków jako węzłów i zależności jako strzałek, oddzielny od kanwy edycyjnej Workflow Buildera | 1 | Widoczny po otwarciu okna |
| Dodanie i usunięcie zależności | Rysowanie i kasowanie połączeń między krokami | 2 | Przeciągnięcie od portu węzła |
| Grupy równoległe i sekwencyjne | Oznaczenie zbioru kroków jako wykonywanych równolegle albo w ścisłej kolejności | 2 | Zaznaczenie węzłów, element zwinięty `Grupa ▼` |
| Bramki dołączenia | Rozejście przebiegu na wiele torów i scalenie wyników według reguły: wszystkie, dowolny, licznik | 3 | Menu kontekstowe węzła rozgałęzienia |
| Reguły warunkowe gałęzi | Warunek decydujący, którą gałąź zależności respektować w toku wykonania | 3 | Panel popover krawędzi |
| Ścieżka krytyczna | Podświetlenie najdłuższego łańcucha zależności determinującego czas całego przebiegu | 2 | Przełącznik `Ścieżka krytyczna` |
| Walidacja grafu | Wykrycie cykli i zależności sprzecznych — ostrzeżenie przed zapisem, zapis pozostaje możliwy | 1 | Automatyczna, sygnalizacja na węźle |
| Maszyna stanów przebiegu | Formalny model stanów kroku i całego przebiegu z przejściami: zaplanowany → w toku → sukces, błąd, przerwany | 3 | Zakładka `Stany` |
| Transakcje kompensujące | Kroki wycofujące skutki przy błędzie w połowie procesu | 4 | Polecenie języka naturalnego, tryb administracyjny |
| Powiązanie z orkiestracją MultitaskingAI | Wskaźnik i spięcie: zależności procesu respektowane przez warstwę orkiestracji środowiska MultitaskingAI | 2 | Znacznik kontekstowy `MultitaskingAI` |
| Powiększenie i przewijanie grafu | Nawigacja po złożonych procesach wielokrokowych | 1 | Kółko myszy, przeciągnięcie |
| Eksport mapy zależności | Zapis grafu jako obrazu lub dokumentu do dokumentacji procesu | 3 | Menu kebab (⋮) nagłówka |

**Zachowanie i stany:** zależności respektowane w toku wykonania procesu — zadanie kolejnego kroku nie rozpoczyna się przed zakończeniem kroku, od którego zależy. Stan „konflikt zależności” sygnalizowany czerwonym obrysem węzła — ostrzeżenie, zapis pozostaje możliwy.

```
 Makieta — Orchestrator (rozszerzenie boczne, stan spoczynku)
 ═══════════════════════════════════════════════
  Zależności — „Raport tygodniowy sprzedaży”   ⋮
 ───────────────────────────────────────────────
        ┌──────────┐
        │ Pobierz  │
        │ dane     │
        └────┬─────┘
             │ zależność: zakończone
        ┌────▼─────┐        ┌──────────┐
        │ Waliduj  │───────►│ Wyślij   │ ← ścieżka
        │ wynik    │        │ raport   │   krytyczna
        └──────────┘        └──────────┘
 ───────────────────────────────────────────────
  [ Grupa ▼ ]  [ MultitaskingAI ]
 ═══════════════════════════════════════════════
```

### 3.7. Execution Monitor

| Pole | Treść |
|---|---|
| Rola | Obserwacja kolejnych uruchomień automatyki i niezawodność przebiegów |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zawartość | Status każdego przebiegu, błędy wymagające interwencji, metryki |
| Warstwa | 2 — wywołanie znacznikiem `Przebiegi ▼` |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Tabela i oś czasu przebiegów | Historia uruchomień z datą, czasem trwania, statusem końcowym i wyzwalaczem | 1 | Widoczna po otwarciu okna |
| Status na żywo | Aktualizacja stanu trwającego przebiegu kanałem WebSocket, krok po kroku, z podświetleniem bieżącego kroku | 1 | Widoczny po otwarciu okna |
| Przeglądarka logów przebiegu | Pełny zapis zdarzeń pojedynczego uruchomienia, z filtrowaniem po poziomie (informacja, ostrzeżenie, błąd), wyszukiwaniem i eksportem | 3 | Akcja `log` w wierszu przebiegu |
| Szuflada szczegółów błędu | Treść błędu, ślad stosu, krok, w którym wystąpił, sugerowana przyczyna | 3 | Kliknięcie plakietki błędu |
| Ponowne uruchomienie | Powtórzenie przebiegu zakończonego błędem, całości lub od nieudanego kroku, z zachowaniem kontekstu | 2 | Przycisk `Ponów` w wierszu przebiegu |
| Przerwanie trwającego przebiegu | Zatrzymanie procesu w toku — natychmiastowe, z mechanizmem „Cofnij” dostępnym przez krótki czas po akcji | 1 | Przycisk `Przerwij` w nagłówku |
| Panel metryk | Wskaźnik powodzenia w czasie, średni czas trwania, liczba przebiegów w okresie, opóźnienie p95, koszt wyrażony liczbą tokenów i wywołań | 2 | Zakładka `Metryki` |
| Reguły alarmowania | Warunki i kanały powiadomień: błąd, przekroczenie czasu, brak uruchomienia, spadek wskaźnika powodzenia — kanał funkcji globalnej Mobile i pozostałe kanały | 3 | Panel popover `Reguły alarmowania` |
| Budżety czasu przebiegu | Maksymalny czas przebiegu i kroku oraz alarm przy jego przekroczeniu | 3 | Panel popover reguły alarmowej |
| Punkty wznowienia | Zapis stanu przebiegu pozwalający wznowić od miejsca przerwania bez powtarzania ukończonych kroków | 3 | Menu kebab (⋮) wiersza przebiegu |
| Filtrowanie i wyszukiwanie | Po zakresie dat, statusie, nazwie przepływu | 2 | Ikona filtra |
| Eksport raportu przebiegów | Zapis zestawienia jako dokument do audytu i sprawozdawczości | 3 | Menu kebab (⋮) nagłówka |
| Drążenie do poziomu kroku | Rozwinięcie pojedynczego przebiegu do statusu każdego kroku osobno, ze skrótem do Orchestratora i logów kroku | 2 | Akcja `kroki` w wierszu przebiegu |
| Podgląd i odtworzenie ładunku | Wgląd w dane wejściowe i wyjściowe każdego kroku oraz odtworzenie przebiegu z tym ładunkiem | 4 | Polecenie języka naturalnego, tryb administracyjny |

**Zachowanie i stany:** okno aktualizuje się na żywo w miarę przebiegu wykonania. Stany przebiegu: zaplanowany, w toku, zakończony sukcesem, zakończony błędem, przerwany ręcznie.

```
 Makieta — Execution Monitor (rozszerzenie boczne, stan spoczynku)
 ═══════════════════════════════════════════════
  Przebiegi — „Raport tygodniowy sprzedaży”    ⋮
 ───────────────────────────────────────────────
  2026-08-03 07:00  ✔ sukces  2m 14s  [log][kroki]
  2026-07-27 07:00  ✕ błąd    0m 48s  [log][kroki]
  2026-07-20 07:00  ✔ sukces  2m 09s  [log][kroki]
 ───────────────────────────────────────────────
  Wskaźnik powodzenia (8 tyg.):  ▇▇▇▁▇▇▇▇  87%
 ═══════════════════════════════════════════════
```

---

## 4. Katalog funkcji i narzędzi

Katalog obejmuje dziewięć rodzin funkcjonalnych. Każda pozycja opisana jest nazwą, działaniem i zależnościami technicznymi (biblioteki backendu w języku Go, formaty, integracje). Backend platformy realizowany jest w Go, spójnie z warstwą serwera i kanałem WebSocket kontraktów komunikacji.

### 4.1. Rodzina Projektowanie przepływu (Workflow Builder)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Kanwa węzłowa | Wizualny edytor przeciąganych węzłów kroków z łącznikami, powiększaniem i przewijaniem, minimapą i autoukładem | Warstwa klienta: React Flow lub własna kanwa SVG; serwer: model grafu w JSON |
| Biblioteka typów kroku | Zestaw węzłów: wywołanie modelu, akcja w module, HTTP/API, warunek, pętla, opóźnienie, podprzepływ, transformacja danych, punkt kontrolny | Rejestr typów kroków |
| Mapowanie danych | Wizualne łączenie wyjść jednego kroku z wejściami kolejnego, wyrażenia i szablony na polach | `text/template`, `tidwall/gjson`, `PaesslerAG/jsonpath` |
| Silnik wyrażeń warunkowych | Ewaluacja warunków rozgałęzień i filtrów w bezpiecznym języku wyrażeń, bez wykonywania dowolnego kodu | `google/cel-go` |
| Podprzepływy reużywalne | Wydzielenie fragmentu procesu jako nazwanego, wielokrotnego bloku wpinanego w inne automatyki | Odwołanie do definicji automatyki jako komponentu własnego |
| Pętle i iteracje | Węzeł iterujący po kolekcji, pętla warunkowa, pętla po stronach interfejsu programistycznego (paginacja) z limitem bezpieczeństwa | Iterator w silniku wykonania; `time` na limity |
| Transformacje danych w przepływie | Węzły: JSON↔YAML↔CSV, parsowanie i serializacja, filtr, sortowanie, agregacja, wybór pól, scalanie kolekcji | `encoding/json`, `gopkg.in/yaml.v3`, `encoding/csv` |
| Test i symulacja przebiegu | Uruchomienie w trybie testowym bez wpięcia produkcyjnego, podgląd wyniku każdego kroku, dane próbne na wejściach | Silnik wykonania w trybie próbnym; izolacja efektów ubocznych |
| Walidacja definicji | Wykrycie kroków bez połączenia, cykli i brakujących parametrów — ostrzeżenie na kanwie, zapis pozostaje możliwy | `heimdalr/dag`, własny walidator schematu |
| Import i eksport definicji | Zapis i wczytanie przepływu w JSON/YAML wraz ze zgodnością formatu wymiany | `encoding/json`, `gopkg.in/yaml.v3`, `santhosh-tekuri/jsonschema` |
| Wersjonowanie i porównanie | Historia wersji definicji, porównanie strukturalne grafu i pól, przywrócenie wcześniejszej wersji | `wI2L/jsondiff`; magazyn wersji w bazie |
| Duplikowanie i szablonowanie | Kopia procesu jako punkt wyjścia; zapis dowolnego przepływu jako szablonu z parametrami | Parametryzacja definicji (zmienne szablonu) |
| Notatki i dokumentacja kroku | Adnotacje opisowe przy węźle, opis całej automatyki, generowany diagram dokumentacyjny | Render do Mermaid/DOT (`goccy/go-graphviz`) |
| Tagowanie i wyszukiwanie | Kategoryzacja automatyk, wyszukiwanie po nazwie, tagu, module docelowym i wyzwalaczu | Indeks pełnotekstowy (Postgres `tsvector`) |
| Cofnij i ponów edycję | Pełna historia edycji kanwy w sesji pracy | Stos komend po stronie klienta |

### 4.2. Rodzina Wyzwalacze i harmonogramy (Scheduler)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Kreator wyrażenia cron | Budowa reguły cron w formie przyjaznej, z podglądem składni technicznej i najbliższych uruchomień | `robfig/cron/v3`, `adhocore/gronx` |
| Wzorce cykliczności | Zestawy gotowe: co minutę, godzinnie, codziennie, co tydzień, co miesiąc, kwartalnie, niestandardowo | Warstwa wzorców nad kreatorem wyrażenia cron |
| Harmonogram interwałowy | Uruchomienie co N sekund, minut lub godzin niezależnie od zegara ściennego | `time.Ticker`, trwały stan ostatniego uruchomienia |
| Widok kalendarza uruchomień | Miesięczny i tygodniowy podgląd zaplanowanych oraz historycznych przebiegów | Komponent kalendarza klienta; projekcja z reguł po stronie serwera |
| Strefy czasowe i zmiana czasu | Reguła harmonogramu w wybranej strefie, poprawna obsługa przejścia czasu letniego i zimowego | IANA tz (`time.LoadLocation`), baza `tzdata` |
| Okna wykonania | Ograniczenie przedziałów, w których uruchomienie następuje (np. wyłącznie godziny robocze) | Reguły filtrujące na warstwie planera |
| Wyzwalacz webhook | Unikatowy adres wejściowy automatyki; przyjęcie ładunku HTTP z weryfikacją podpisu (HMAC) i deduplikacją | `net/http`, `crypto/hmac`, magazyn kluczy idempotencji |
| Wyzwalacz zmiany pliku | Uruchomienie na dodanie, zmianę lub usunięcie pliku w Library albo Project Library | `fsnotify/fsnotify`; zdarzenia magazynu Library |
| Wyzwalacz na wyniku modelu | Uruchomienie, gdy wynik wywołania modelu spełnia warunek (np. klasyfikacja „pilne”) | Silnik wyrażeń warunkowych + hak zdarzenia modelu |
| Wyzwalacz zdarzenia modułu | Reakcja na zdarzenie dowolnego modułu: zatwierdzenie zmian w Developer, wdrożenie w Apps, nowy plik w Library | Szyna zdarzeń platformy (WebSocket, kontrakty komunikacji) |
| Wyzwalacze łańcuchowe | Zakończenie jednej automatyki — sukcesem lub błędem — wyzwala kolejną | Zdarzenia cyklu życia przebiegu |
| Wiele wyzwalaczy na proces | Jeden przepływ powiązany z wieloma niezależnymi regułami czasu i zdarzeń jednocześnie | Lista wyzwalaczy w definicji automatyki |
| Wstrzymanie i wznowienie | Przełącznik zatrzymujący cykliczność bez usuwania definicji; stan aktywny lub wstrzymany | Flaga stanu; planer pomija wstrzymane |
| Uruchomienie wsteczne | Wykonanie przebiegów dla przeszłych, pominiętych terminów w zadanym zakresie dat | Generator terminów z reguły + kolejkowanie |
| Historia wyzwoleń | Zapis rzeczywistych momentów wyzwolenia z przyczyną: harmonogram, zdarzenie, uruchomienie ręczne | Tabela zdarzeń wyzwoleń |
| Nadzór obecności uruchomień | Alarm, gdy oczekiwane uruchomienie nie nastąpiło w oknie tolerancji | Nadzorca na warstwie planera + reguły alarmowania |

### 4.3. Rodzina Kolejkowanie (Queue Manager)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Silnik kolejek — jedenaście akcji | Pełny zestaw: `enqueue`, `dequeue`, `delay`, `retry`, `pause`, `resume`, `split`, `merge`, `route`, `branch`, `condition` | `riverqueue/river` (Postgres — jedna baza, spójność transakcyjna z resztą danych platformy) |
| Zasięgi kolejki | Kolejka globalna, lokalna, dla modelu, dla agenta, dla projektu — spójne z MODEL-KONFIGURACJI 5.3 | Klucz zasięgu w definicji kolejki |
| Priorytety i zmiana kolejności | Przeciąganie zadań między priorytetami; kolejki priorytetowe | Kolejki wagowe silnika |
| Współbieżność i ograniczanie przepustowości | Limit równoległych zadań na kolejkę i zasięg; ograniczanie tempa przetwarzania i szczytów obciążenia | Semafory i mechanizm żetonów (`golang.org/x/time/rate`) |
| Ponawianie z wycofaniem | Ponowienie z wykładniczym wycofaniem i rozproszeniem, z limitem prób | `cenkalti/backoff/v4` |
| Kolejka zadań martwych | Zadania trwale nieudane trafiają do osobnej kolejki do przeglądu ręcznego lub ponowienia zbiorczego | Osobna kolejka + akcja `route` |
| Idempotencja zadań | Klucz idempotencji zapobiegający podwójnemu wykonaniu przy ponowieniu lub duplikacie wywołania przychodzącego | Magazyn kluczy idempotencji z czasem życia |
| Wizualizacja głębokości | Wykres liczby zadań oczekujących w czasie, obciążenie per zasięg | Metryki kolejki + komponent wykresu |
| Szuflada szczegółów zadania | Ładunek, znaczniki czasu, dotychczasowe próby, powiązany przebieg, logi | Widok szczegółu; dane z magazynu zadań |
| Akcje zbiorcze | Zastosowanie dowolnej z jedenastu akcji do wielu zaznaczonych zadań naraz | Operacje wsadowe silnika kolejek |
| Ręczne wstrzyknięcie zadania | Dodanie zadania do kolejki poza harmonogramem | `enqueue` z priorytetem |
| Harmonogram opóźnień | Zadania zaplanowane na wskazany termin i cykliczne wewnątrz kolejki | Kolumna `scheduled_at` w magazynie zadań |
| Filtrowanie i wyszukiwanie zadań | Po statusie, dacie, przepływie źródłowym, zasięgu | Indeks zadań |

### 4.4. Rodzina Orkiestracja i zależności (Orchestrator)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Graf zależności (DAG) | Widok kroków jako węzłów i zależności jako strzałek, zorientowany na relacje, nie na treść | `heimdalr/dag`, `gonum.org/v1/gonum/graph` |
| Grupy równoległe i sekwencyjne | Oznaczenie zbioru kroków jako równoległych albo ściśle sekwencyjnych | Model grupowania w definicji |
| Bramki dołączenia | Rozejście na wiele torów i scalenie wyników według reguły: wszystkie, dowolny, licznik | Reguły synchronizacji w silniku wykonania |
| Reguły warunkowe gałęzi | Warunek decydujący, którą zależność respektować w toku wykonania | Silnik wyrażeń warunkowych |
| Ścieżka krytyczna | Podświetlenie najdłuższego łańcucha zależności determinującego czas przebiegu | Algorytm najdłuższej ścieżki na grafie (`gonum`) |
| Walidacja grafu | Wykrycie cykli i zależności sprzecznych — ostrzeżenie przed zapisem, zapis pozostaje możliwy | `heimdalr/dag` (detekcja cyklu) |
| Maszyna stanów przebiegu | Formalny model stanów kroku i całego przebiegu z przejściami: zaplanowany → w toku → sukces, błąd, przerwany | `qmuntal/stateless` |
| Transakcje kompensujące | Kroki wycofujące skutki przy błędzie w połowie procesu (wycofanie rozproszone) | Wzorzec saga w silniku wykonania; kroki kompensujące |
| Powiązanie z orkiestracją MultitaskingAI | Wskaźnik i spięcie: zależności procesu respektowane przez warstwę ról MultitaskingAI po skonfigurowaniu | Integracja z panelem orkiestracji (Koncepcja 11.6) |
| Eksport mapy zależności | Zapis grafu jako obraz lub dokument (PNG, SVG, DOT, Mermaid) do dokumentacji | `goccy/go-graphviz`, render Mermaid |

### 4.5. Rodzina Spinanie modułów i modeli w potoki

| Nazwa | Co robi | Zależności |
|---|---|---|
| Krok „akcja w module” | Wywołanie operacji dowolnego z 15 modułów jako kroku (np. „wygeneruj raport w Research”, „wdroż w Apps”, „zapisz w Library”) | Kontrakty komunikacji modułów; jawne powiązanie (MODEL-KONFIGURACJI 5.8) |
| Krok „wywołanie modelu” | Wywołanie modelu z promptem, parametrami i kanałem (API, CLI, SSH, HTTP), wynik jako dane kroku | Warstwa kanałów modeli (Architektura 8) |
| Potok wielomodelowy | Kilka różnych modeli w jednym przebiegu — jeden klasyfikuje, drugi redaguje, trzeci weryfikuje | Rejestr modeli; mapowanie danych |
| Pętla iteracyjna model↔model | Naprzemienna praca dwóch modeli aż do warunku stopu (generacja↔krytyka), spójna z relacjami wykonawców (Koncepcja 11.3) | Warunek stopu z silnika wyrażeń + limit iteracji |
| Krok agenta | Wpięcie agenta — komponentu własnego z modułu Agents — jako Wykonawcy kroku, z jego umiejętnościami i konektorami | Odwołanie do agenta; uprawnienia z Permissions Center |
| Krok HTTP / REST / GraphQL | Wywołanie zewnętrznego interfejsu programistycznego z metodą, nagłówkami, ciałem, uwierzytelnieniem i parsowaniem odpowiedzi | `go-resty/resty`; uwierzytelnienie Bearer i OAuth2 (`golang.org/x/oauth2`) |
| Krok rozszerzenia (MCP) | Wywołanie narzędzia z rejestru rozszerzeń (Danaco Plugin, rozszerzenie osobiste, serwer MCP) jako kroku | Kontrakt rozszerzeń (Architektura 12); klient MCP |
| Krok skryptu w piaskownicy | Wykonanie krótkiego skryptu transformującego dane w izolowanym środowisku, bez dostępu do systemu poza zakresem ustawionym przez Operatora | Piaskownica `goja` lub izolowany proces; izolacja z MODEL-KONFIGURACJI 5.2 |
| Krok kolejki międzymodułowej | Przekazanie zadania do kolejki innego modułu lub środowiska i oczekiwanie na wynik | Silnik kolejek + zasięgi kolejki |
| Emisja zdarzeń wyjściowych | Emisja zdarzenia platformy na zakończenie kroku lub przebiegu, wyzwalająca inne automatyki i moduły | Szyna zdarzeń; wyzwalacze zdarzenia modułu i łańcuchowe |

### 4.6. Rodzina Nadzór, niezawodność i alarmowanie (Execution Monitor)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Oś czasu i tabela przebiegów | Historia uruchomień z datą, czasem trwania, statusem i wyzwalaczem | Magazyn przebiegów; komponent osi czasu |
| Status na żywo krok po kroku | Aktualizacja trwającego przebiegu kanałem WebSocket, podświetlenie bieżącego kroku | Kanał WebSocket (kontrakty komunikacji) |
| Przeglądarka logów przebiegu | Pełny zapis zdarzeń przebiegu z filtrem po poziomie, wyszukiwaniem i eksportem | Strukturalne logi (`log/slog`); magazyn logów |
| Szuflada szczegółów błędu | Treść, ślad stosu, krok wystąpienia, sugerowana przyczyna wyznaczona przez Wykonawcę | Wychwyt błędu + podpowiedź modelu |
| Ponowienie przebiegu | Powtórzenie całości lub od nieudanego kroku, z zachowaniem kontekstu | Punkty wznowienia |
| Przerwanie trwającego przebiegu | Natychmiastowe zatrzymanie z mechanizmem „Cofnij” przez krótki czas po akcji | Sygnał anulowania (`context.Context`) |
| Panel metryk | Wskaźnik powodzenia w czasie, średni czas trwania, liczba przebiegów, opóźnienie p95, koszt wyrażony liczbą tokenów i wywołań | `prometheus/client_golang`; agregacja metryk |
| Reguły alarmowania | Warunki i kanały powiadomień: błąd, przekroczenie czasu, brak uruchomienia, spadek wskaźnika — kanał Mobile i pozostałe | Silnik wyrażeń warunkowych; integracja z funkcją globalną Mobile |
| Punkty wznowienia | Zapis stanu przebiegu pozwalający wznowić od miejsca przerwania po awarii, bez powtarzania ukończonych kroków | Trwały zapis stanu kroku; idempotencja zadań |
| Drążenie do poziomu kroku | Rozwinięcie przebiegu do statusu każdego kroku, skrót do Orchestratora i logów kroku | Powiązanie krok ↔ log ↔ stan |
| Eksport raportu przebiegów | Zestawienie jako dokument (CSV, PDF, Markdown) do audytu i sprawozdawczości | `encoding/csv`; render PDF (`johnfercher/maroto`) |
| Podgląd i odtworzenie ładunku | Wgląd w dane wejściowe i wyjściowe każdego kroku oraz odtworzenie przebiegu z tym ładunkiem | Magazyn ładunków z retencją i redakcją sekretów |
| Budżety czasu i poziom obsługi | Maksymalny czas przebiegu i kroku oraz alarm przy przekroczeniu | Timery + reguły alarmowania |

### 4.7. Rodzina Asystent budowy automatyki (Chat Window)

| Nazwa | Co robi | Zależności |
|---|---|---|
| Generowanie kroku z opisu | Polecenie w języku naturalnym przekształcane w gotowy krok Workflow Buildera; stan „sugestia oczekująca”, zastosowanie jawne | Wywołanie modelu; mapowanie na typy kroku |
| Generowanie całego przepływu z opisu | Opis celu procesu przekształcany w kompletny szkic przepływu do dalszej edycji | Model + walidacja definicji |
| Wyjaśnienie procesu | Wykonawca opisuje działanie istniejącego przepływu krok po kroku na podstawie definicji | Odczyt definicji + model |
| Debugowanie konwersacyjne | Przekazanie logu błędu z Execution Monitor i uzyskanie sugerowanej poprawki kroku | Kontekst logu przebiegu + model |
| Wzmianki (`@`) | Odwołanie do kroku, kolejki, uruchomienia wprost w treści polecenia | Indeks encji automatyki |
| Szablony poleceń | Zapisane polecenia typowe dla budowy automatyk | Magazyn szablonów per użytkownik |
| Akcje na wiadomości | Kopiuj, zastosuj sugestię wprost do Workflow Buildera, odrzuć | Zastosowanie zmiany różnicowej do definicji |
| Sterowanie pętlą wykonawczą z poziomu rozmowy | Zlecenie uruchomienia, wstrzymania, wznowienia i korekty przebiegu poleceniem języka naturalnego | Kontrakt sterowania Execution Loop Window |

### 4.8. Rodzina Sekrety, poświadczenia i bezpieczeństwo operacyjne

| Nazwa | Co robi | Zależności |
|---|---|---|
| Skarbiec poświadczeń | Jawnie zarządzane klucze i tokeny kroków (interfejsy programistyczne, HTTP, konektory), przechowywane poza bazą danych, przywoływane referencją z kroków | Magazyn sekretów poza bazą (Model danych 8.1); szyfrowanie (`crypto/aes`, KMS) |
| Referencje kluczy w krokach | Krok odwołuje się do sekretu przez nazwę lub referencję, nigdy nie zawiera wartości w definicji | Rozwiązywanie referencji w czasie wykonania |
| Redakcja sekretów w logach | Maskowanie wartości wrażliwych w logach i podglądzie ładunku | Reguły redakcji na warstwie logowania |
| OAuth2 i odświeżanie tokenów | Obsługa przepływów OAuth2 dla konektorów, odświeżanie wygasłych tokenów | `golang.org/x/oauth2`; magazyn tokenów odświeżających |
| Dziennik audytu | Zapis, kto i kiedy zmienił definicję, uruchomił przebieg, przerwał go lub zmienił harmonogram | Tabela audytu; niezmienny log |
| Izolacja procesu automatyki | Uruchomienie kroków w izolowanym procesie, katalogu i sieci według zakresów izolacji technicznej — sterowane z okna konfiguracji | MODEL-KONFIGURACJI 5.2 (poziomy zasięgu izolacji) |

### 4.9. Rodzina Biblioteka i dystrybucja

| Nazwa | Co robi | Zależności |
|---|---|---|
| Biblioteka szablonów przepływów | Gotowe wzorce procesów: raport cykliczny, przetwarzanie wsadowe, przenoszenie danych, potok wielomodelowy — do zastosowania i edycji | Repozytorium szablonów; parametryzacja definicji |
| Katalog konektorów | Rejestr dostępnych akcji w modułach, wywołań HTTP, rozszerzeń MCP i agentów jako klocków kroku | Rejestr typów kroku + kontrakt rozszerzeń |
| Współdzielenie automatyk w organizacji | Udostępnianie gotowych automatyk jako komponentów własnych w obrębie organizacji | Pamięć aplikacji (zasób użytkownika) |
| Parametryzacja przy wpięciu | Przy wpinaniu automatyki do modułu docelowego uzupełnienie parametrów: odbiorca, zakres dat, zasób | Zmienne szablonu; formularz wpięcia |
| Wersje publikowane i robocze | Rozdzielenie wersji roboczej, edytowanej, od opublikowanej, wykonywanej produkcyjnie | Wersjonowanie definicji + flaga publikacji |

Rodziny Sekrety, poświadczenia i bezpieczeństwo operacyjne oraz Biblioteka i dystrybucja nie mają odrębnych okien — przenikają wszystkie siedem okien modułu: skarbiec poświadczeń jest dostępny z każdego kroku wymagającego uwierzytelnienia, biblioteka szablonów otwiera się przy tworzeniu automatyki, a dziennik audytu i redakcja sekretów działają w tle Execution Monitora i Execution Loop Window.

---

## 5. Katalog elementów interfejsu

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa i sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|
| Węzeł kroku | Blok na kanwie, kształt zależny od typu kroku | Reprezentacja pojedynczego kroku procesu | Średni blok, kanwa swobodna | 1 — widoczny bez interakcji | domyślny, zaznaczony (obrys złoty), błąd walidacji (obrys czerwony), w trakcie testu (animacja pulsu) | Kliknięcie otwiera panel właściwości kroku; przeciągnięcie zmienia położenie | Workflow Builder |
| Łącznik między krokami | Linia ze strzałką | Określenie kolejności i zależności wykonania | Cienka linia, waga niska | 1 — widoczny bez interakcji | domyślny, najechanie (pogrubienie), nieprawidłowy (przerywany, czerwony) | Kliknięcie zaznacza do edycji lub usunięcia | Workflow Builder, Orchestrator |
| Znacznik kontekstowy okna | Lekka pigułka w pasku kontekstu | Wywołanie okna pomocniczego jako rozszerzenia bocznego | Bardzo mała pigułka z `▼` | 2 — kliknięcie znacznika | zwinięty (domyślny), aktywny (okno otwarte) | Otwiera odpowiednie okno jako kolumnę boczną; ponowne kliknięcie zamyka | Pasek kontekstu obszaru roboczego |
| Przycisk „Test” | `.dn-btn--zarys`, rozmiar `sm` | Uruchomienie symulacji procesu bez wpięcia produkcyjnego | Mały przycisk z obrysem | 2 — przycisk nagłówka | domyślny, w toku (spinner), zakończony (plakietka wyniku) | Uruchamia przebieg testowy widoczny na kanwie | Workflow Builder |
| Przycisk „Zapisz” | `.dn-btn--glowny` | Zapis definicji automatyki jako komponentu własnego | Mały–średni przycisk wypełniony | 1 — widoczny bez interakcji | domyślny (zawsze aktywny — nie blokowany walidacją), ostrzeżenie walidacji (plakietka obok przycisku), zapisano (potwierdzenie toast) | Zapisuje wersję i zamyka stan „niezapisane zmiany” niezależnie od ostrzeżeń walidacji | Workflow Builder |
| Wiersz zadania pętli wykonawczej | Wiersz listy zadań ze wskaźnikiem stanu | Prezentacja zadania zleconego Wykonawcy przez Koordynatora | Wiersz listy, waga średnia | 1 — widoczny bez interakcji | zaplanowane, w toku (puls), zakończone, błędne, ponawiane | Kliknięcie rozwija wynik kontroli jakości i komunikaty sterujące zadania | Execution Loop Window |
| Sterowanie przebiegiem pętli | Zestaw przycisków `.dn-btn--zarys` i `.dn-btn--blad` | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia w toku | Mała grupa przycisków w nagłówku | 1 — widoczne bez interakcji | bezczynne (nieaktywne), przebieg w toku (aktywne) | Zmienia stan pętli natychmiast; przerwanie udostępnia „Cofnij” przez krótki czas | Execution Loop Window |
| Suwak `[wł|wył]` harmonogramu | `.dn-suwak` | Aktywacja i wstrzymanie cykliczności bez usuwania definicji | Mały przełącznik | 2 — widoczny po otwarciu Schedulera | aktywny (tor złoty), wstrzymany (tor neutralny) | Natychmiastowa zmiana stanu harmonogramu | Scheduler |
| Kreator cron | Zestaw pól `.dn-select` + podgląd tekstowy | Budowa reguły cykliczności bez znajomości składni | Panel średni, kilka pól w rzędzie | 2 — widoczny po otwarciu Schedulera | domyślny, nieprawidłowa kombinacja (obrys błędu na polu) | Aktualizuje na żywo podgląd wyrażenia i najbliższych uruchomień | Scheduler |
| Wiersz zadania kolejki | `.dn-tabela` wiersz z menu akcji | Prezentacja pojedynczego zadania i dostęp do jedenastu akcji silnika kolejek | Wiersz tabeli, waga średnia | 1 w oknie kolejki; akcje w warstwie 3 przez menu kebab (⋮) | oczekuje, w toku, błąd (plakietka czerwona), zakończone | Menu kebab (⋮) rozwija dostępne akcje (`enqueue`, `dequeue`, `delay`, `retry`, …) | Queue Manager |
| Wskaźnik głębokości kolejki | Mały wykres słupkowy w linii | Orientacyjny obraz obciążenia kolejki w czasie | Mały, w nagłówku okna | 1 w oknie kolejki | niski, średni, wysoki poziom (kolor) | Najechanie pokazuje wartość liczbową dnia | Queue Manager |
| Węzeł grafu zależności | Blok w widoku Orchestratora | Reprezentacja kroku w kontekście relacji, nie treści | Średni blok, kanwa grafu | 1 w oknie zależności | domyślny, na ścieżce krytycznej (podświetlenie złote), konflikt (czerwony) | Kliknięcie pokazuje listę zależności wchodzących i wychodzących | Orchestrator |
| Wiersz historii przebiegu | `.dn-tabela--paski` wiersz | Prezentacja jednego uruchomienia automatyki | Wiersz tabeli | 1 w oknie przebiegów | sukces (zielony), błąd (czerwony), w toku (pulsujący), przerwany (szary) | Kliknięcie rozwija log i status krokowy | Execution Monitor |
| Plakietka statusu przebiegu | `.dn-plakietka--stan` | Szybka sygnalizacja wyniku uruchomienia | Bardzo mała pigułka z kropką | 1 — widoczna bez interakcji | sukces, ostrzeżenie, błąd, informacja | Element informacyjny | Execution Monitor, Execution Loop Window, Queue Manager, Project Dashboard |
| Przycisk „Ponów” | `.dn-btn--zarys`, ikonowy | Powtórzenie nieudanego przebiegu lub kroku | Mały przycisk | 2 — widoczny w wierszu przebiegu | domyślny, w toku ponowienia (spinner) | Uruchamia nowy przebieg od punktu błędu lub od początku | Execution Monitor |
| Przycisk „Przerwij” | `.dn-btn--blad` | Zatrzymanie trwającego przebiegu — akcja o podwyższonej wadze | Mały przycisk wariantu błędu | 1 — widoczny w nagłówku | zawsze aktywny (domyślny), w toku zatrzymywania (spinner) | Bez aktywnego przebiegu: komunikat „brak przebiegu do przerwania”. Z przebiegiem w toku: zatrzymuje proces natychmiast i udostępnia „Cofnij” przez krótki czas | Execution Monitor, Execution Loop Window |
| Dialog potwierdzenia przerwania | `.dn-modal` wywoływany po włączeniu ustawienia `potwierdzenie_przerwania` | Dodatkowe zabezpieczenie przed przypadkowym przerwaniem — nie jest domyślną bramką | Duży panel wyśrodkowany | 3 — wywoływany akcją przerwania przy włączonym ustawieniu | ustawienie wyłączone (domyślnie), ustawienie włączone; przy włączonym: otwarty lub zamknięty | Przy wyłączonym ustawieniu przerwanie następuje od razu, „Cofnij” dostępne po akcji. Przy włączonym dialog poprzedza przerwanie; potwierdzenie zatrzymuje proces, anulowanie zamyka dialog bez skutku | Execution Monitor; ustawienie w oknie konfiguracji |
| Panel metryk (wykres) | Wykres liniowy lub słupkowy | Wizualizacja wskaźnika powodzenia w czasie | Średni panel w kolumnie okna przebiegów | 2 — zakładka `Metryki` | z danymi, brak danych (stan pusty) | Najechanie punktu pokazuje wartość i datę | Execution Monitor |
| Kafel „Zbuduj automatykę” | `.dn-karta`, wariant kafla strefy 2 | Wejście do modułu Automations ze strony głównej | Średni kafel w siatce czterech | 1 — widoczny bez interakcji | domyślny, najechanie (akcent złoty) | Otwiera Workflow Builder w trybie tworzenia nowej automatyki | Strefa 2 strony głównej |

---

## 6. Przepływy pracy w module

### 6.1. Budowa i uruchomienie automatyki od podstaw

```
Strona główna, strefa 2 → kafel „Automations” → „Zbuduj automatykę”
        │
        ▼
Chat Window — polecenie języka naturalnego opisujące cel procesu
        │
        ▼
Workflow Builder — projektowanie kroków, warunków, kolejności
        │
        ├──► Orchestrator  — ustalenie zależności między krokami
        ├──► Scheduler     — ustalenie cykliczności lub zdarzenia wyzwalającego
        │
        ▼
Zapis jako AUTOMATYKA (komponent własny)
        │
        ▼
Wpięcie automatyki do sesji modułu docelowego (np. Developer, Research)
        │
        ▼
Wyzwolenie → Execution Loop Window: Koordynator dekomponuje zlecenie
             na zadania kroków i przydziela je Wykonawcom
        │
        ▼
Queue Manager przetwarza zadania · Execution Monitor rejestruje przebieg
```

### 6.2. Pętla wykonawcza uruchomienia przepływu pracy

```
Wyzwalacz (harmonogram · zdarzenie · uruchomienie ręczne · Chat Window)
        │
        ▼
Execution Loop Window — Koordynator przyjmuje zlecenie uruchomienia
        │
        ▼
Dekompozycja przepływu na zadania kroków wg grafu zależności
        │
        ▼
┌───────────────────────────────────────────────────────┐
│ Koordynator → Wykonawca: zadanie kroku                │
│ Wykonawca → Koordynator: wynik kroku                  │
│ Koordynator: kontrola jakości wyniku                  │
│    ├─ warunek spełniony → następne zadanie            │
│    └─ warunek niespełniony → ponowienie z wycofaniem  │
└───────────────────────────────────────────────────────┘
        │  (iteracja aż do wyczerpania grafu zadań)
        ▼
Zamknięcie przebiegu · emisja zdarzenia wyjściowego
        │
        ▼
Execution Monitor — zapis wyniku, metryk i logów przebiegu
```

Użytkownik steruje pętlą z Chat Window poleceniem języka naturalnego oraz bezpośrednio z Execution Loop Window przyciskami wstrzymania, wznowienia, przerwania i korekty zlecenia.

### 6.3. Reakcja na błąd wykonania

```
Execution Loop Window: zadanie kroku zakończone błędem
        │
        ▼
Koordynator: decyzja o ponowieniu albo przekazaniu błędu dalej
        │
        ▼
Execution Monitor: przebieg zakończony błędem (✕)
        │
        ▼
Szuflada szczegółów błędu — treść, krok, sugerowana przyczyna
        │
        ├──► [Ponów]  ─────────────────────► nowy przebieg od punktu błędu
        │
        ├──► Queue Manager → akcja `retry` lub `route` na zadaniu źródłowym
        │
        └──► Chat Window → przekazanie logu, prośba o sugestię poprawki
                    │
                    ▼
             sugestia zastosowana w Workflow Builderze → zapis nowej wersji
```

### 6.4. Powiązanie ze środowiskiem MultitaskingAI

```
Stan wyjściowy: Automations działa niezależnie od MultitaskingAI
                (połączenie nie jest domyślne)
        │
        │  ustawienie w oknie konfiguracji:
        │  powiązanie silników kolejek
        ▼
Silnik kolejek Automations  ◄────────────►  Silnik kolejek MultitaskingAI
Scheduler, Execution Loop Window,  ───────►  panel orkiestracji, sekcja
Execution Monitor                            „Harmonogram i automatyki”
        │
        ▼
Automatyka stanowi operacyjne zaplecze harmonogramów i kolejek
pętli pracy ciągłej 24/7/365 (rozdz. 13.6 Koncepcji platformy)
```

### 6.5. Wpięcie gotowej automatyki do dowolnego modułu

```
AUTOMATYKA (komponent własny, zapisana)
        │
        ▼
Wybór modułu docelowego podczas pracy w sesji, np.:
        ├──► Developer   — uruchomienie testów po zatwierdzeniu zmian
        ├──► Research    — cykliczne odświeżenie źródeł
        ├──► Workspace   — cykliczny raport projektu (rozdz. 7.3 dokumentu Workspace)
        └──► Apps        — wdrożenie po zatwierdzeniu (Deployment Panel)
```

---

## 7. Punkty sterowania z okna konfiguracji

Zachowania modułu są personalizowane z okna konfiguracji, zgodnie z warstwowością (globalna → środowisko → projekt → sesja) i zasadą „brak ustawienia = wartość domyślna”. Poniższe punkty rozwijają zakres **5.3 Akcje**, **5.4 Zachowanie modeli** i **5.8 Integracje** modelu konfiguracji. Żaden nie jest twardą blokadą — każdy jest ustawieniem konfiguracyjnym.

| Punkt sterowania | Co Operator personalizuje | Zakres (MODEL-KONFIGURACJI) | Wartość domyślna |
|---|---|---|---|
| Akcje kolejki dostępne | Które z jedenastu akcji silnika kolejek są udostępnione | 5.3 Akcje | pełny zestaw jedenastu |
| Zasięg kolejki | Globalna, lokalna, dla modelu, dla agenta, dla projektu | 5.3 Akcje | lokalna |
| Obsługa błędów kolejki | `retry`, `route` lub `pause` przy błędzie zadania | 5.3 Akcje | `retry` |
| Polityka ponawiania | Liczba prób, typ wycofania (stały, wykładniczy), rozproszenie, limit czasu | 5.3 Akcje | 3 próby, wykładniczy z rozproszeniem |
| Współbieżność i przepustowość | Limit równoległych zadań i tempa przetwarzania per zasięg | 5.3 Akcje | bez ograniczenia |
| Kolejka zadań martwych | Włączenie kolejki zadań martwych i polityka przenoszenia | 5.3 Akcje | włączona, przenoszenie po wyczerpaniu prób |
| Tryb pracy w tle zadania | Wykonanie w tle bez blokowania karty sesji | 5.3 Akcje | tryb w tle aktywny |
| Widoczność Execution Loop Window | Stałe wyświetlanie okna pętli wykonawczej albo otwieranie znacznikiem `Pętla wykonawcza` | 5.1 Aplikacja / ustawienia okna modułu | stałe wyświetlanie przy aktywnym zleceniu |
| Poziom szczegółowości pętli wykonawczej | Zakres komunikatów sterujących Koordynator ↔ Wykonawca prezentowanych w oknie pętli | 5.1 Aplikacja / 5.9 Historia | komunikaty zadań i decyzje o ponowieniu |
| Kontrola jakości w pętli wykonawczej | Warunek oceny wyniku kroku i decyzja o ponowieniu, przejściu dalej lub zatrzymaniu | 5.3 Akcje / 5.4 Zachowanie modeli | ocena warunkiem przepływu, ponowienie przy niespełnieniu |
| Limit iteracji pętli | Maksymalna liczba iteracji pętli wykonawczej dla jednego zlecenia | 5.3 Akcje | 10 |
| Powiązanie z silnikiem kolejek MultitaskingAI | Spięcie silnika kolejek Automations z silnikiem MultitaskingAI | 5.3 Akcje / 5.8 Integracje | wyłączone |
| Kanał modelu w krokach | API, CLI, SSH lub HTTP dla kroków wywołania modelu | 5.4 Zachowanie modeli | API |
| Model bazowy kroku | Wybór modelu per krok i per automatyka; różne modele w potoku | 5.4 Zachowanie modeli | zależny od kanału |
| Relacje Wykonawców w pętli model↔model | Niezależna, przekazywanie, naprzemienna, iteracyjna | 5.4 Zachowanie modeli | niezależna |
| Powiązanie automatyki z modułem docelowym | Do których modułów kierowane są kroki „akcja w module” | 5.8 Integracje | brak powiązania (jawnie dodawane) |
| Integracja MultitaskingAI ↔ Automations | Ujawnienie harmonogramów i kolejek w panelu orkiestracji | 5.8 Integracje | wyłączona |
| Strefa czasowa harmonogramów | Domyślna strefa dla reguł Schedulera | 5.1 Aplikacja / warstwa środowiska | CET |
| Okna wykonania | Dozwolone przedziały czasu uruchomień | 5.3 Akcje | bez ograniczenia |
| Reguły alarmowania | Warunki i kanały powiadomień: błąd, przekroczenie budżetu czasu, brak uruchomienia | 5.1 Aplikacja (Powiadomienia) / 5.8 | przy błędzie → Mobile |
| Skarbiec poświadczeń | Zarządzanie kluczami kroków przechowywanymi poza bazą danych | 5.7 Rozszerzenia / Model danych 8.1 | pusty, klucze jawnie dodawane |
| Izolacja procesu automatyki | Zakresy izolacji technicznej kroków: katalog, sieć, konto i token, model procesu, serwer wykonania | 5.2 Procesy / 5.11 Izolacja | wszystkie wyłączone |
| Retencja przebiegów i logów | Czas przechowywania historii przebiegów, logów i ładunków | 5.9 Historia | bez limitu |
| Redakcja sekretów w logach | Maskowanie wartości wrażliwych | 5.7 / 5.9 | włączona |
| Potwierdzenie przed przerwaniem | Dialog poprzedzający zatrzymanie przebiegu | ustawienia okna modułu | wyłączone (przerwanie natychmiastowe + „Cofnij”) |
| Profil automatyki | Zapis konfiguracji modułu jako nazwanego profilu do wielokrotnego użycia | 7. Profile konfiguracji | brak |

---

## 8. Stany, dane i powiązania z innymi modułami

### 8.1. Dane wykorzystywane przez Wykonawcę w module

| Źródło danych | Okno pochodzenia | Uwaga |
|---|---|---|
| Definicja procesu | Workflow Builder | Kroki, warunki, zmienne |
| Reguły harmonogramu | Scheduler | Cykliczność i zdarzenia wyzwalające |
| Stan i pozycje kolejki | Queue Manager | Zadania oczekujące i przetwarzane |
| Zależności procesu | Orchestrator | Kolejność respektowana w toku wykonania |
| Zlecenie i stan pętli wykonawczej | Execution Loop Window | Dekompozycja zlecenia, komunikaty sterujące, decyzje o ponowieniu |
| Status i logi wykonania | Execution Monitor | Historia przebiegów i błędów |

### 8.2. Stany automatyki jako komponentu własnego

| Stan | Znaczenie | Gdzie widoczny |
|---|---|---|
| Robocza | W trakcie budowy, niezapisana lub niewpięta do żadnej sesji | Workflow Builder |
| Aktywna | Zapisana, harmonogram włączony, gotowa do wykonania | Scheduler (suwak wł.), Execution Monitor |
| Wstrzymana | Harmonogram wyłączony, definicja zachowana | Scheduler (suwak wył.) |
| W toku wykonania | Przebieg aktualnie trwa | Execution Loop Window (pętla w toku), Execution Monitor (status na żywo) |
| Błędna | Ostatni przebieg zakończony błędem wymagającym interwencji | Execution Monitor (plakietka błędu) |

### 8.3. Powiązania konfigurowalne z innymi modułami i środowiskami

```
                    AUTOMATIONS  (moduł, ta karta)
                            │
        ┌───────────────────┼────────────────────┐
        ▼                                        ▼
   MultitaskingAI                         Dowolny moduł platformy
   harmonogramy i kolejki                 (Studio, Developer, Research,
   pętli pracy ciągłej                     Workspace, Apps, …)
   (rozdz. 13.6 Koncepcji)                 zadania objęte procesem
```

| Moduł lub środowisko docelowe | Charakter powiązania | Typ | Miejsce ustanowienia |
|---|---|---|---|
| MultitaskingAI | Automations stanowi operacyjne zaplecze harmonogramów i kolejek pętli pracy ciągłej; silnik kolejek Automations powiązany z silnikiem kolejek MultitaskingAI | Konfiguracyjne | Okno konfiguracji oraz ustawienia okna modułu Automations |
| Dowolny moduł platformy | Procesy zdefiniowane w Automations obejmują zadania z dowolnego innego modułu | Konfiguracyjne | Workflow Builder — krok „akcja w module” |
| Agents | Agent wpięty jako Wykonawca kroku, z własnymi umiejętnościami i konektorami | Konfiguracyjne | Workflow Builder — krok agenta; uprawnienia z Permissions Center |

Żadne z powyższych powiązań nie jest domyślnie aktywne — automatyka nowo zbudowana działa w pełnej odrębności od środowiska MultitaskingAI, dopóki połączenie nie zostanie świadomie ustanowione.

---

## 9. Scenariusze użycia

### 9.1. Cykliczny raport sprzedażowy bez stałego nadzoru

1. Użytkownik otwiera moduł Automations ze strefy 2 strony głównej i wybiera „Zbuduj automatykę”.
2. W Chat Window opisuje cel procesu; Wykonawca przedstawia szkic przepływu jako sugestię oczekującą, Użytkownik ją stosuje.
3. W Workflow Builderze doprecyzowuje proces: pobranie danych, walidacja, generowanie dokumentu, wysyłka pocztą.
4. W Orchestratorze ustala, że wysyłka nie następuje przed pozytywną walidacją danych.
5. W Schedulerze ustawia cykliczność „co poniedziałek, 7:00”.
6. Automatyka trafia jako komponent własny do sesji modułu Workspace, gdzie zasila projekt raportowy.
7. Co tydzień proces uruchamia się samodzielnie: Execution Loop Window prowadzi pętlę Koordynator ↔ Wykonawca krok po kroku, Execution Monitor rejestruje przebieg, a błąd — na przykład brak danych źródłowych — trafia do Mobile jako powiadomienie zgodnie z regułą alarmowania.

### 9.2. Testy uruchamiane po każdej zmianie w repozytorium

1. W module Developer (CodeStudio) powstaje potrzeba uruchamiania testów po każdym zatwierdzeniu zmian.
2. W Automations powstaje przepływ z wyzwalaczem zdarzeniowym powiązanym ze zdarzeniem repozytorium.
3. Queue Manager przyjmuje zadanie „uruchom testy” natychmiast po wyzwoleniu zdarzenia, korzystając z akcji `enqueue`.
4. Koordynator w Execution Loop Window przydziela zadanie Wykonawcy, odbiera wynik i ocenia go względem warunku przepływu.
5. Wynik trafia do Execution Monitor; przy błędzie akcja `route` kieruje zadanie do kolejki przeglądu ręcznego zamiast automatycznego ponowienia.

### 9.3. Pełna pętla pracy ciągłej ze środowiskiem MultitaskingAI

1. Zespół projektowy prowadzi proces budowy produktu w module Apps częściowo autonomicznie poza godzinami pracy.
2. W oknie konfiguracji Użytkownik świadomie łączy silnik kolejek Automations z silnikiem kolejek środowiska MultitaskingAI.
3. Harmonogram i kolejki widoczne dotąd wyłącznie w Automations pojawiają się dodatkowo w sekcji „Harmonogram i automatyki” panelu orkiestracji MultitaskingAI.
4. Koordynator środowiska MultitaskingAI planuje kolejne odcinki pracy, korzystając z tego samego zaplecza harmonogramowego, a pętle wykonawcze poszczególnych automatyk pozostają widoczne w Execution Loop Window; Always On Display nadzoruje całość z perspektywy Użytkownika.

### 9.4. Potok wielomodelowy z kontrolą jakości w pętli wykonawczej

1. Użytkownik buduje automatykę porządkującą materiały źródłowe: pierwszy model klasyfikuje dokumenty, drugi redaguje streszczenia, trzeci weryfikuje zgodność streszczenia ze źródłem.
2. W Workflow Builderze każdemu krokowi przypisany jest inny model, a mapowanie danych łączy wyjścia kolejnych kroków.
3. Koordynator w Execution Loop Window prowadzi pętlę: zleca klasyfikację, odbiera wynik, ocenia go warunkiem przepływu, po czym zleca redakcję i weryfikację.
4. Wynik weryfikacji niespełniający warunku powoduje ponowienie kroku redakcji z korektą zlecenia; licznik iteracji i ponowień widoczny jest w oknie pętli.
5. Zamknięty przebieg emituje zdarzenie wyjściowe, które wyzwala kolejną automatykę publikującą materiał w module docelowym.

---

*Koniec dokumentu. Moduł Automations — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
