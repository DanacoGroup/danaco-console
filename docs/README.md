# Danaco Console — dokumentacja systemowa platformy

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

Zbiór stanowi kompletną dokumentację systemową platformy Danaco Console w wersji 2.0 — docelową specyfikację produktu, architektury i funkcjonalności, przekazywaną zespołowi projektowemu i deweloperskiemu jako źródło prawdy. Opisuje koncepcję platformy, warstwy techniczne, model danych, model konfiguracji, kontrakty komunikacji, bezpieczeństwo, warstwę modeli, warstwę rozszerzeń, system wizualny oraz pełny projekt interfejsu użytkownika: środowiska, moduły, okna platformowe i funkcje globalne. Deweloper buduje na jej podstawie produkt, projektant prowadzi na jej podstawie projekt interfejsu.

Katalog główny zbioru zawiera wyłącznie niniejszy plik `README.md` — jedyny indeks zbioru. Wszystkie opracowania leżą w sześciu katalogach tematycznych: `architektura/`, `specyfikacje/`, `interfejs-uzytkownika/`, `moduly/`, `srodowiska/` oraz `funkcje-globalne/`. Podkatalogi nie zawierają własnych indeksów.

---

## Spis treści

1. [Indeks opracowań](#1-indeks-opracowań)
2. [Przewodnik czytania](#2-przewodnik-czytania)
3. [Zasady projektowe platformy](#3-zasady-projektowe-platformy)
4. [Terminologia wiążąca](#4-terminologia-wiążąca)
5. [Standard dokumentu i organizacja zbioru](#5-standard-dokumentu-i-organizacja-zbioru)

---

## 1. Indeks opracowań

Indeks obejmuje wszystkie opracowania zbioru w sześciu grupach odpowiadających katalogom tematycznym. Odesłania zapisane są jako ścieżki względne wobec katalogu głównego zbioru.

### 1.1. `architektura/` — poziom systemowy platformy

Grupa definiuje produkt jako system: koncepcję platformy, warstwy techniczne, modele danych i konfiguracji, kontrakty komunikacji, bezpieczeństwo, warstwę modeli, izolację oraz warstwę rozszerzeń. Grupa ustala rdzeń, na którym opierają się wszystkie pozostałe opracowania zbioru. Czyta ją w całości deweloper jako źródło prawdy technicznej; projektant czyta z niej koncepcję platformy oraz model konfiguracji.

| Dokument | Przeznaczenie |
|---|---|
| `architektura/koncepcja-platformy.md` | Dokument nadrzędny zbioru: charakter platformy, pojęcia podstawowe, komunikacja operacyjna, warstwy widoczności, warstwy architektury, konfigurowalność zależności, strona główna, nawigacja i sesje, cztery środowiska, moduły, funkcje globalne oraz zasady nadrzędne funkcjonalności. |
| `architektura/architektura.md` | Źródło prawdy architektonicznej: model wdrożenia, stos technologiczny, warstwy systemu, warstwa komunikacji operacyjnej, model sesji i procesów, realizacja zadań obliczeniowych, warstwa danych, komunikacja klient–serwer, izolacja, konfiguracja, rozszerzenia i synchronizacja. |
| `architektura/model-danych.md` | Pełny model danych platformy: konto i urządzenia, środowiska i moduły, komponenty własne, projekty, karty sesji i sesje, pamięć wielopoziomowa, zadania, kolejki, harmonogram i orkiestracja, kanały modeli, agenci, rozszerzenia oraz profile i punkty izolacji. |
| `architektura/model-konfiguracji.md` | Struktura okna konfiguracji i wszystkie zakresy ustawień — aplikacja, procesy, akcje, zachowanie i tożsamość modeli, prompty systemowe, rozszerzenia, integracje, historia, pamięć, izolacja, karty sesji, komponenty własne, oba kanały komunikacji operacyjnej i warstwy widoczności — wraz z warstwowością globalna → środowisko → projekt → sesja. |
| `architektura/kontrakty-komunikacji.md` | Kontrakty komunikacji między klientem a rdzeniem platformy: pojedynczy kanał WebSocket, schematy ładunków oraz diagramy sekwencji i przepływu dla strony głównej, środowisk, kart sesji, modułów, obu kanałów komunikacji operacyjnej, silnika kolejek, orkiestracji, agentów, pamięci, konfiguracji i rozszerzeń. |
| `architektura/integracja-modeli.md` | Warstwa dostawcy modelu: cztery kanały integracji modeli sztucznej inteligencji (API, CLI, SSH, HTTP), mechanizm adapterów, przypisanie modelu do sesji lub roli, tożsamość i persona modelu oraz miejsce konfiguracji kanału w oknie konfiguracji. |
| `architektura/izolacja-i-zaleznosci.md` | Konfigurowalność zależności i punkty izolacji: zakresy izolacji technicznej procesu sesji i izolacji kontekstu, łączenie zakresów w politykę, profile izolacji, przypisanie punktów do sesji, roli i projektu, współdzielenie historii i pamięci oraz interfejs konfiguracji punktów izolacji. |
| `architektura/bezpieczenstwo-i-uwierzytelnianie.md` | Model bezpieczeństwa: rejestracja konta właściciela, logowanie, metody dodatkowe uwierzytelniania, odzyskiwanie konta, token dostępu i nawiązywanie połączenia WebSocket, przechowywanie danych dostępowych poza bazą danych oraz autoryzacja w obu kanałach komunikacji operacyjnej. |
| `architektura/rozszerzenia.md` | Warstwa rozszerzeń: wtyczki, umiejętności, konektory i serwery MCP, jednolity kontrakt rozszerzenia, dwa źródła rozszerzeń (Danaco Plugin oraz Personal), cykl życia rozszerzenia oraz udział rozszerzeń w warstwie komunikacji operacyjnej. |

### 1.2. `specyfikacje/` — zakres funkcjonalny do implementacji

Grupa przekłada rdzeń architektoniczny na zakres funkcjonalny platformy: co robią moduły, jakie typy okien operacyjnych występują w systemie i jak działa warstwa agentowa. Grupa jest pomostem między architekturą a dokumentacją projektową. Czyta ją deweloper wyznaczający zakres implementacji oraz projektant ustalający repertuar okien i zachowań.

| Dokument | Przeznaczenie |
|---|---|
| `specyfikacje/specyfikacja-modulow.md` | Operacyjna specyfikacja modułów platformy: cel modułu, okna operacyjne wraz z warstwami widoczności, funkcjonalności, przepływy pracy, dane wykorzystywane przez AI oraz powiązania konfigurowalne między modułami. |
| `specyfikacje/specyfikacja-okien-operacyjnych.md` | Katalog wszystkich typów okien operacyjnych występujących w modułach i środowiskach: cel, zawartość, zachowanie, warstwa widoczności, sposób wywołania i powiązane operacje każdego okna. |
| `specyfikacje/specyfikacja-agentow.md` | Warstwa agentowa platformy: tworzenie agenta, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć i uprawnienia, agent jako komponent własny ponad środowiskami, pełny przebieg pętli wykonawczej oraz Centrum uprawnień. |

### 1.3. `interfejs-uzytkownika/` — warstwa kliencka ponad środowiskami

Grupa zawiera opracowania obowiązujące ponad środowiskami i modułami: punkt wejścia użytkownika i nawigację, przepływ okien, elementy okien, katalog komponentów, system projektowy oraz okna platformowe konfiguracji i ustawień. Grupa jest pierwszym materiałem projektanta interfejsu i wiążącym odniesieniem dla dewelopera warstwy klienckiej.

| Dokument | Przeznaczenie |
|---|---|
| `interfejs-uzytkownika/strona-glowna-i-nawigacja.md` | Strona główna platformy jako centrum dowodzenia i nawigacja: układ kolumnowy, trzy strefy wyboru, przepływ strona główna → środowisko → moduł → okno, mechanika kart sesji, przełączanie przestrzeni roboczej, boczna nawigacja modułów i panel orkiestracji. |
| `interfejs-uzytkownika/przeplyw-okien.md` | Przepływ okien od uruchomienia aplikacji, przez uwierzytelnienie i stronę główną, po okno robocze modułu. |
| `interfejs-uzytkownika/elementy-okien.md` | Specyfikacja elementów okien przepływu głównego wraz z warstwą widoczności i sposobem wywołania każdego elementu. |
| `interfejs-uzytkownika/katalog-komponentow.md` | Katalog komponentów interfejsu: budowa, warianty, stany, zachowanie i zastosowanie każdego komponentu. |
| `interfejs-uzytkownika/system-wizualny.md` | System projektowy platformy: tokeny kolorów w motywie jasnym i ciemnym, typografia, cienie, siatka kolumnowa, ikonografia i biblioteka komponentów wraz z ich zastosowaniem w oknach komunikacji, oknach operacyjnych modułów i kolumnach strony głównej. |
| `interfejs-uzytkownika/konfiguracja.md` | Okno Konfiguracji: struktura zakresów, pola sterujące wywołaniem modelu, konfiguracja obu kanałów komunikacji operacyjnej i warstw widoczności oraz Panel prowenancji wywołania. |
| `interfejs-uzytkownika/ustawienia.md` | Okno Ustawień poziomu aplikacji: ustawienia konta, urządzenia i zachowania aplikacji. |

### 1.4. `moduly/` — dokumentacja projektowa modułów

Grupa opisuje każdy moduł platformy jako gotowy do budowy zakres wykonawczy: komplet okien operacyjnych, katalog elementów każdego okna, funkcje, narzędzia, przepływy pracy, stany i punkty sterowania z okna konfiguracji. Każde opracowanie adresowane jest jednocześnie do projektanta (co, gdzie i w jakiej formie) oraz do dewelopera (co zbudować).

| Dokument | Przeznaczenie |
|---|---|
| `moduly/assistant.md` | Moduł Assistant — rozmowa z AI jako podstawowa forma pracy: okna operacyjne, funkcje, narzędzia i punkty sterowania modułu. |
| `moduly/agents.md` | Moduł Agents — fabryka ekspertów: tworzenie, konfigurowanie i przypisywanie agentów wraz z cyklem tworzenia eksperta i punktami sterowania. |
| `moduly/workspace.md` | Moduł Workspace — projekty, zadania i praca operacyjna: okna projektu, listy zadań, przepływy realizacji i powiązania z orkiestracją. |
| `moduly/studio.md` | Moduł Studio — tworzenie i redagowanie treści: katalog okien, funkcji, narzędzi i zależności modułu. |
| `moduly/design.md` | Moduł Design — projektowanie materiałów wizualnych: okna edycyjne i podglądu, narzędzia projektowe oraz przepływy pracy z materiałem graficznym. |
| `moduly/developer.md` | Moduł Developer — praca z kodem i repozytorium: pełny katalog funkcji, narzędzi, okien operacyjnych i punktów sterowania modułu. |
| `moduly/terminal.md` | Moduł Terminal — praca w powłoce systemowej: okna sesji powłoki, sterowanie wykonaniem poleceń i zakres uprawnień. |
| `moduly/browser.md` | Moduł Browser — praca z zasobami sieciowymi: okna przeglądania, pobieranie i przetwarzanie treści oraz powiązania z modułami wiedzy. |
| `moduly/research.md` | Moduł Research — badania, analizy i opracowania źródłowe: komplet okien operacyjnych, katalog funkcji i narzędzi oraz punkty sterowania modułu. |
| `moduly/library.md` | Moduł Library — zarządzanie wiedzą i zbiorem dokumentów: komplet okien operacyjnych, katalog funkcji i narzędzi wraz z zależnościami technicznymi. |
| `moduly/translate.md` | Moduł Translate — tłumaczenia i praca wielojęzyczna: okna tłumaczenia, kontrola jakości przekładu i zarządzanie terminologią. |
| `moduly/automations.md` | Moduł Automations — automatyzacje i przepływy procesów: budowa przepływu, wyzwalacze, harmonogram oraz integracja z silnikiem kolejek. |
| `moduly/apps.md` | Moduł Apps — budowa aplikacji biznesowych: okna projektowania aplikacji, model danych aplikacji, publikacja i uruchamianie. |
| `moduly/roundtable.md` | Moduł Roundtable — współpraca wielu modeli nad jednym zagadnieniem: okna dyskusji wykonawców, tryby uzgadniania i prezentacja wyniku zbiorczego. |
| `moduly/diagnostics.md` | Moduł Diagnostics — nadzór nad stanem platformy i diagnostyka: monitory, dzienniki, wskaźniki wykonania i narzędzia niskiego poziomu. |

### 1.5. `srodowiska/` — dokumentacja projektowa środowisk

Grupa opisuje cztery środowiska platformy: powłokę środowiska, dostępne w nim moduły, komplet okien, przepływy pracy i rozróżnienie projektowe między środowiskami. Grupa jest podstawą pracy projektanta nad przestrzenią roboczą i punktem odniesienia dla dewelopera przy budowie powłok środowiskowych.

| Dokument | Przeznaczenie |
|---|---|
| `srodowiska/talkin.md` | Środowisko TalkIn — wiedza, komunikacja i praca z treścią: okna środowiska, elementy powłoki i przepływy pracy. |
| `srodowiska/workspace.md` | Środowisko WorkSpace — produktywność, organizacja i realizacja projektów: okna środowiska, elementy powłoki i przepływy pracy. |
| `srodowiska/codestudio.md` | Środowisko CodeStudio — rozwój oprogramowania: okna środowiska, praca z repozytorium i przepływy wytwórcze. |
| `srodowiska/multitaskingai.md` | Środowisko MultitaskingAI — warstwa centralna: role zespołu, tryby współpracy wykonawców, silnik kolejek, orkiestracja i zależności, hierarchia decyzji, praca ciągła, makiety okien i katalog elementów interfejsu. |

### 1.6. `funkcje-globalne/` — funkcje obecne ponad środowiskami

Grupa opisuje funkcje działające we wszystkich środowiskach, modułach, projektach i sesjach, niezależnie od aktywnej przestrzeni roboczej. Czyta ją projektant przy projektowaniu zachowań globalnych interfejsu oraz deweloper przy budowie warstwy funkcji globalnych i protokołów urządzeń.

| Dokument | Przeznaczenie |
|---|---|
| `funkcje-globalne/always-on-display.md` | Funkcja globalna Always On Display — postać wizualna, reguły wyzwalania proaktywnych sugestii, katalog sugestii, tor głosowy, zachowanie per środowisko i per moduł, warstwy widoczności, stany i punkty sterowania. |
| `funkcje-globalne/mobile.md` | Funkcja globalna Mobile — zakres funkcjonalny na urządzeniu przenośnym, układ kolumnowy i punkty przełamania, komplet ekranów, oba kanały komunikacji operacyjnej, praca bez połączenia, powiadomienia wypychane oraz protokół parowania urządzenia. |

---

## 2. Przewodnik czytania

### 2.1. Ścieżka dewelopera

1. `architektura/koncepcja-platformy.md` — produkt jako całość: środowiska, moduły, sesje, zasady nadrzędne.
2. `architektura/architektura.md` — warstwy systemu, stos technologiczny, model sesji i procesów.
3. `architektura/model-danych.md` — encje, relacje i trwałość stanu.
4. `architektura/kontrakty-komunikacji.md` — protokół klient–rdzeń, schematy ładunków, sekwencje.
5. `architektura/model-konfiguracji.md` — zakresy ustawień i warstwowość konfiguracji.
6. `architektura/bezpieczenstwo-i-uwierzytelnianie.md`, `architektura/integracja-modeli.md`, `architektura/izolacja-i-zaleznosci.md`, `architektura/rozszerzenia.md` — warstwy systemowe budowane na ustalonym rdzeniu.
7. `specyfikacje/specyfikacja-modulow.md`, `specyfikacje/specyfikacja-okien-operacyjnych.md`, `specyfikacje/specyfikacja-agentow.md` — zakres funkcjonalny do implementacji.
8. `interfejs-uzytkownika/przeplyw-okien.md`, `interfejs-uzytkownika/elementy-okien.md`, `interfejs-uzytkownika/katalog-komponentow.md`, `interfejs-uzytkownika/strona-glowna-i-nawigacja.md` — warstwa kliencka: przepływ, elementy, komponenty, punkt wejścia.
9. `srodowiska/talkin.md`, `srodowiska/workspace.md`, `srodowiska/codestudio.md`, `srodowiska/multitaskingai.md`, następnie katalog `moduly/` — budowa powłok środowiskowych i modułów, moduł po module.
10. `interfejs-uzytkownika/konfiguracja.md`, `interfejs-uzytkownika/ustawienia.md`, `funkcje-globalne/always-on-display.md`, `funkcje-globalne/mobile.md` — okna platformowe, warstwa funkcji globalnych i urządzenie przenośne.

### 2.2. Ścieżka projektanta

1. `architektura/koncepcja-platformy.md` — struktura produktu i zasady nadrzędne.
2. Rozdział 3 niniejszego dokumentu — zasady projektowe wiążące dla całego interfejsu.
3. `interfejs-uzytkownika/system-wizualny.md` — tokeny, typografia, siatka kolumnowa, ikonografia, biblioteka komponentów.
4. `interfejs-uzytkownika/strona-glowna-i-nawigacja.md` — punkt wejścia użytkownika i nawigacja.
5. `interfejs-uzytkownika/przeplyw-okien.md` oraz `interfejs-uzytkownika/elementy-okien.md` — przepływ okien i elementy przepływu głównego.
6. `interfejs-uzytkownika/katalog-komponentow.md` — komponenty, warianty i stany.
7. `specyfikacje/specyfikacja-okien-operacyjnych.md`, `specyfikacje/specyfikacja-modulow.md`, `specyfikacje/specyfikacja-agentow.md` — typy okien operacyjnych, zakres modułów i warstwa agentowa.
8. `srodowiska/talkin.md`, `srodowiska/workspace.md`, `srodowiska/codestudio.md`, `srodowiska/multitaskingai.md` — projekt środowisk.
9. Katalog `moduly/` — projekt okien modułowych, moduł po module.
10. `interfejs-uzytkownika/konfiguracja.md`, `interfejs-uzytkownika/ustawienia.md` oraz katalog `funkcje-globalne/` — okna platformowe i funkcje globalne.
11. `architektura/model-konfiguracji.md` — punkty sterowania interfejsem po stronie konfiguracji.

---

## 3. Zasady projektowe platformy

Rozdział zbiera ustalenia wiążące dla całego zbioru. Każde opracowanie stosuje je bez wyjątku, a każdy element interfejsu opisywany w dokumentacji jest z nimi zgodny.

### 3.1. Dwa kanały komunikacji operacyjnej

Komunikacja operacyjna platformy prowadzona jest dwoma kanałami. Oba są elementami pierwszoplanowymi architektury i występują w każdym środowisku, module i oknie roboczym.

| Kanał | Okno | Rola |
|---|---|---|
| Użytkownik ↔ Wykonawca | **Chat Window** | Centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami platformy |
| Koordynator ↔ Wykonawca | **Execution Loop Window** | Pętla wykonawcza: koordynacja zadań, nadzór nad realizacją, orkiestracja działań i kontrola realizacji procesów |

**Chat Window** przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Każdy moduł i każde środowisko udostępnia to okno w tym samym miejscu układu — w lewej kolumnie obszaru roboczego.

**Execution Loop Window** zawiera bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem — wstrzymanie, wznowienie, przerwanie i korektę zlecenia. Okno otwierane jest jako kolumna sąsiadująca z Chat Window.

Role rozdzielone są jednoznacznie: **Użytkownik** zleca i zatwierdza; **Koordynator** jest komponentem orkiestrującym platformy — dekomponuje zlecenie, przydziela i nadzoruje zadania; **Wykonawca** jest AI, agentem lub systemem wykonawczym realizującym zadania.

### 3.2. Układ interfejsu — wyłącznie pionowy

Obowiązuje wyłącznie układ pionowy, oparty na podziale lewa–prawa. Wszystkie okna robocze, okna komunikacji, przestrzenie modułowe i okna pomocnicze rozmieszczone są w kolumnach sąsiadujących poziomo. Okna pomocnicze otwierane są jako rozszerzenia boczne. Regulacji podlega wyłącznie szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Lewa, stała, pełna wysokość | Chat Window — kanał Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca, otwierana | Execution Loop Window — kanał Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu: okna edycyjne, okna podglądu, monitory |
| Kolejne kolumny boczne | Okna pomocnicze i panele, otwierane jako rozszerzenia boczne po prawej stronie obszaru roboczego |

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

### 3.3. Cztery warstwy widoczności

Zasadą nadrzędną interfejsu jest reguła: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

Każdy element interfejsu należy do dokładnie jednej z czterech warstw widoczności; przy opisie okna, panelu, paska narzędzi i pojedynczej funkcji podaje się jej warstwę.

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu. | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wykonawcy, środowiska i trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w oknie komunikacji, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Ukrywanie funkcjonalności realizują cztery mechanizmy:

- **Menu progresywne.** Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`), a lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Funkcjonalność umieszczana jest w panelach bocznych, panelach wysuwanych, oknach popover i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji.
- **Znaczniki kontekstowe.** Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

Makiety w dokumentacji rysowane są w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3. Elementy warstw 2–4 opisywane są w tabelach elementów okna z podaniem warstwy i sposobu wywołania.

---

## 4. Terminologia wiążąca

Jedno pojęcie ma dokładnie jedną nazwę we wszystkich dokumentach zbioru. Synonimy traktuje się jak usterkę.

| Pojęcie | Brzmienie wiążące |
|---|---|
| Platforma | Danaco Console |
| Główne okno komunikacji | Chat Window (kanał Użytkownik ↔ Wykonawca) |
| Okno pętli wykonawczej | Execution Loop Window (kanał Koordynator ↔ Wykonawca) |
| Realizator zadań | Wykonawca (AI / agent / system wykonawczy) |
| Komponent orkiestrujący | Koordynator |
| Rozmieszczenie okien | Układ pionowy (podział lewa–prawa) |
| Okno dodatkowe | Okno pomocnicze otwierane jako rozszerzenie boczne |

---

## 5. Standard dokumentu i organizacja zbioru

**Nagłówek redakcyjny.** Każdy plik `.md` zbioru otwiera jednolity nagłówek redakcyjny w postaci tabeli o polach: Produkt, Rodzaj, Opis, Producent, Twórca, Wersja, Status, Data. Wersja każdego dokumentu wynosi **v2.0**, status — **Deweloperski**, data — **2026-08-06**.

**Stopka.** Każdy dokument zamyka dwuwierszowa stopka produktowa, oddzielona linią poziomą: wiersz pierwszy podaje kursywą nazwę produktu, rodzaj platformy oraz wersję `v2.0`, wiersz drugi — notę praw autorskich Danaco Holding Group Sp. z o.o. z odesłaniem do licencji oraz adresem kontaktowym. Wzorem obowiązującym jest stopka zamykająca niniejszy dokument.

**Konwencja nazw plików.** Nazwy plików zapisywane są małymi literami, wyrazy rozdziela się myślnikiem, bez polskich znaków diakrytycznych. Konwencja obowiązuje jednakowo we wszystkich katalogach tematycznych. Jedynym wyjątkiem jest `README.md`.

**Katalog główny.** W katalogu głównym zbioru nie umieszcza się opracowań — zawiera on wyłącznie plik `README.md`. Każde opracowanie merytoryczne leży w jednym z sześciu katalogów tematycznych.

**Przydział opracowania do katalogu według zakresu.** O przynależności dokumentu rozstrzyga jego zakres przedmiotowy:

| Katalog | Zakres opracowań |
|---|---|
| `architektura/` | Systemowy poziom produktu: koncepcja platformy, warstwy techniczne i stos, model danych, model konfiguracji, kontrakty komunikacji, bezpieczeństwo i uwierzytelnianie, integracja modeli, izolacja i zależności, warstwa rozszerzeń. |
| `specyfikacje/` | Przekrojowy zakres funkcjonalny ponad pojedynczym modułem i pojedynczym oknem: specyfikacja modułów, katalog typów okien operacyjnych, warstwa agentowa. |
| `interfejs-uzytkownika/` | Warstwa kliencka obowiązująca ponad środowiskami i modułami: strona główna i nawigacja, przepływ okien, elementy okien, katalog komponentów, system wizualny oraz okna platformowe konfiguracji i ustawień. |
| `moduly/` | Dokumentacja projektowa pojedynczego modułu platformy — jeden plik na moduł, nazwany identyfikatorem modułu. |
| `srodowiska/` | Dokumentacja projektowa pojedynczego środowiska platformy — jeden plik na środowisko, nazwany identyfikatorem środowiska. |
| `funkcje-globalne/` | Funkcje działające we wszystkich środowiskach, modułach, projektach i sesjach, niezależnie od aktywnej przestrzeni roboczej. |

**Odsyłacze.** Odesłania do innych opracowań zapisywane są w treści dokumentów jako ścieżki względne wobec katalogu głównego zbioru, na przykład `moduly/studio.md`, `architektura/model-danych.md`.

**Edycja w miejscu.** Dokument aktualizuje się w istniejącym pliku. Kopie robocze, kopie zapasowe, szkice i wersje poprzednie zapisywane obok pod inną nazwą są niedopuszczalne — zbiór zawiera wyłącznie materiały merytoryczne w wersji obowiązującej.

**Kompletność indeksu.** Niniejszy dokument jest jedynym indeksem zbioru. Każdy plik `.md` zbioru występuje w indeksie rozdziału 1, a każda pozycja indeksu wskazuje istniejący plik.

**Język.** Każdy dokument napisany jest profesjonalnym językiem formalnym, oznajmującym, w konwencji analityczno-architektonicznej, pełnymi nazwami i określeniami. Treść merytoryczna ujmowana jest w tabelach zestawczych, schematach tekstowych, diagramach przepływu, szablonach konfiguracji i scenariuszach użycia. Każdy opis elementu interfejsu podaje jego warstwę widoczności i sposób wywołania.

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
