*Dokument specyfikuje interfejs modułu Apps Danaco Console: okna, makiety, elementy, warstwy widoczności i stany.*

# Danaco Console — Moduł Apps — dokument projektowy

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
| **Moduł** | Apps (rozdz. 11.14 Koncepcji platformy) |
| **Środowiska dostępności** | WorkSpace, CodeStudio |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| **Odbiorcy dokumentu** | Designer (co, gdzie, w jakiej formie, do czego), Deweloper (co zbudować) |
| **Stos techniczny odniesienia** | Rdzeń serwera: Go (rejestr rozszerzeń, procesy sesji, sterowanie MCP i konektorami). Kanał sterujący: WebSocket (`extension.list` / `extension.install` / `extension.toggle`, `config.get` / `config.set`) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w UI; klucze jawne; kontrola przez stan wyjściowy i zakres uprawnień, nie przez blokadę; domyślne zachowanie = wykonanie |
| **Źródło faktów** | architektura/koncepcja-platformy.md (rozdz. 4.4–4.6, 8, 9.14, 11.3, 11.9), specyfikacje/specyfikacja-modulow.md (4.14), specyfikacje/specyfikacja-okien-operacyjnych.md (6.14), specyfikacje/specyfikacja-agentow.md, architektura/rozszerzenia.md, architektura/kontrakty-komunikacji.md (rozdz. 1, 5.9), architektura/model-danych.md (rozdz. 9), architektura/architektura.md (rozdz. 9, 12), interfejs-uzytkownika/system-wizualny.md, architektura/model-konfiguracji.md |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Komplet okien operacyjnych modułu — przegląd](#2-komplet-okien-operacyjnych-modułu--przegląd)
3. [Specyfikacja okien operacyjnych — pełny arsenał narzędzi](#3-specyfikacja-okien-operacyjnych--pełny-arsenał-narzędzi)
4. [Katalog elementów interfejsu](#4-katalog-elementów-interfejsu)
5. [Przepływy pracy w module](#5-przepływy-pracy-w-module)
6. [Stany, dane i powiązania z innymi modułami](#6-stany-dane-i-powiązania-z-innymi-modułami)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
9. [Scenariusze użycia](#9-scenariusze-użycia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Dla kogo

Apps jest przestrzenią dla zespołów i użytkowników realizujących **pełne projekty aplikacyjne** — od architektury rozwiązania, przez warstwę frontendową i backendową, po wdrożenie i udostępnienie gotowego oprogramowania. Adresatem są zespoły produktowe budujące aplikacje webowe, mobilne i desktopowe jako spójny produkt, a nie pojedyncze, rozproszone fragmenty kodu, oraz Operatorzy zarządzający katalogiem aplikacji, rozszerzeń i integracji zewnętrznych platformy.

### 1.2. Po co

Moduł spina w **jeden, ciągły proces** wszystkie etapy powstawania produktu cyfrowego, integrując funkcjonalności modułów Developer i Terminal w ramach jednej, ustrukturyzowanej przestrzeni, korzystając — decyzją użytkownika — z zasobów wizualnych modułu Design oraz ze wsparcia orkiestracji wielomodelowej środowiska MultitaskingAI przy pracy nad złożonymi projektami.

### 1.3. Rola modułu — dwie strony jednego cyklu życia oprogramowania

Moduł Apps obejmuje pełen cykl życia oprogramowania na platformie. Jego pierwsza strona to **budowa produktu** — od Architecture Designer, przez Frontend Workspace i Backend Workspace, po Deployment Panel. Druga, komplementarna strona to **dystrybucja i konsumpcja aplikacji oraz rozszerzeń** — katalog, instalacja, uprawnienia i integracje zewnętrzne. Moduł jest **operacyjnym frontem warstwy rozszerzeń platformy** (`architektura/architektura.md`, rozdz. 12; `architektura/rozszerzenia.md`) i zastępuje dziesiątki osobnych programów w obszarze katalogu aplikacji i rozszerzeń, instalacji, uprawnień oraz integracji zewnętrznych (serwery MCP, konektory).

```
 ═══════════════════════════════════════════════════════════════════════════
  STRONA BUDOWY                    │ STRONA DYSTRYBUCJI I KONSUMPCJI
  ───────────────────────────────  │ ───────────────────────────────────────
  Architecture Designer            │ Publikacja i pakowanie
  Frontend / Backend Workspace     │ Katalog rozszerzeń → Instalacja
  Deployment Panel                 │ Uprawnienia i zaufanie
                                   │ Integracje (MCP, konektory)
                                   │ Wykorzystanie operacyjne
                                   │ (moduł Agents, moduły, środowiska)
 ═══════════════════════════════════════════════════════════════════════════
```

Strona dystrybucji i konsumpcji pokrywa sześć obszarów:

| Obszar | Zakres w module Apps |
|---|---|
| Katalog rozszerzeń | Przeglądanie, wyszukiwanie, filtrowanie, kuracja i zestawianie aplikacji, wtyczek, umiejętności, konektorów i serwerów MCP z obu źródeł (Danaco Plugin, Personal) |
| Instalacja i cykl życia | Instalacja Personal (przesłanie pliku), włączanie i wyłączanie, wersjonowanie, aktualizacja, cofnięcie do wcześniejszej wersji, dziennik zmian |
| Uprawnienia, zaufanie i sandbox | Przegląd i nadawanie zakresu uprawnień, izolacja techniczna, weryfikacja podpisów i pochodzenia, audyt uprawnień |
| Serwery MCP | Rejestr serwerów MCP, transporty (stdio/SSE/HTTP), odkrywanie narzędzi, zasobów i promptów, inspektor i próbne wywołania |
| Konektory i integracje zewnętrzne | Katalog konektorów, uwierzytelnianie (OAuth2, klucze, tokeny), testy połączeń, webhooki, mapowanie danych |
| Publikacja i pakowanie | Zapakowanie produktu zbudowanego w module w dystrybuowalną aplikację lub rozszerzenie; prywatny rejestr organizacji |

### 1.4. Granica tematyczna modułu

| Poza granicą | Gdzie to należy |
|---|---|
| Trwały, kanoniczny rejestr rozszerzeń i jego rozgłaszanie | Rejestr żyje w rdzeniu serwera i sekcji „rozszerzenia" **okna konfiguracji** (`architektura/rozszerzenia.md`, rozdz. 7); Apps jest jego operacyjnym frontem, nie drugim źródłem prawdy |
| Podłączenie rozszerzenia jako wykonawcy zadań agenta | Moduł **Agents** — Skills Manager, Connectors Manager, Permissions Center (`architektura/rozszerzenia.md`, rozdz. 6); Apps przygotowuje i testuje, Agents konsumuje |
| Definicja profili izolacji technicznej (macierz, poziomy zasięgu) | Okno **punktów izolacji** (`architektura/model-konfiguracji.md`; `architektura/koncepcja-platformy.md`, rozdz. 4.4–4.6); Apps prezentuje i stosuje zakresy, nie definiuje macierzy |
| Właściwe środowisko programistyczne rozszerzenia | Moduł **Developer** / środowisko **CodeStudio**; Apps udostępnia pakowanie i manifest, nie IDE |
| Orkiestracja wieloagentowa uruchomień rozszerzeń | Środowisko **MultitaskingAI**; Apps dostarcza integracje, orkiestruje je MultitaskingAI |

Granica jest wykonawcza, nie licencyjna: powiązania są jawne i konfigurowalne. Zgodnie z rozdziałem 8 `architektura/rozszerzenia.md`, moduł **nie wprowadza twardej blokady instalacji** — kontrola odbywa się przez stan wyjściowy (Personal = wyłączone) oraz przez zakres uprawnień nadawany przy podłączeniu.

### 1.5. Charakter pracy

| Cecha | Opis |
|---|---|
| Jednostka pracy | Produkt — od architektury po wdrożenie i publikację, jeden spójny przebieg |
| Rytm pracy | Architektura → praca równoległa nad frontendem i backendem → wdrożenie → pakowanie i dystrybucja |
| Skala | Pojedynczy moduł integrujący komponenty programistyczne (kod, kontrola wersji, konsole) oraz warstwę rozszerzeń i integracji w jeden proces |
| Zastosowanie orkiestracji | Realizacja złożonych projektów w module Apps jest naturalnym zastosowaniem środowiska MultitaskingAI — w szczególności podziału na Executor 1 i Executor 2 przy równoległej pracy nad backendem i frontendem |

### 1.6. Miejsce modułu w architekturze platformy

```
ŚRODOWISKO: WorkSpace lub CodeStudio
        │  boczna nawigacja modułów
        ▼
MODUŁ: Apps
        │
        ▼
OKNA OPERACYJNE
Chat Window · Execution Loop Window · Product Builder ·
Architecture Designer · Frontend Workspace · Backend Workspace ·
Deployment Panel · App Catalog · Installed Apps Manager ·
Permissions & Trust Center · Integrations Hub ·
MCP & Connector Console · Publisher Panel
        │
        ├──► Developer, Terminal   — integracja funkcjonalności w jednym procesie
        ├──► Design                — zasoby wizualne warstwy frontendowej
        ├──► Agents                — konsumpcja przygotowanych rozszerzeń
        ├──► Automations           — automatyka zdarzeniowa i health-checki
        ├──► Diagnostics           — logi i błędy integracji
        └──► MultitaskingAI        — Executor 1 / Executor 2, praca równoległa
                                      nad złożonym projektem
```

---

## 2. Komplet okien operacyjnych modułu — przegląd

| Okno | Rola w module | Forma wiodąca | Warstwa | Sposób wywołania | Punkt wejścia? |
|---|---|---|---|---|---|
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Nie — stale obecne |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; prowadzenie i nadzór zadań modułu | Kolumna sąsiadująca z Chat Window | 1 | Widoczne przy aktywnym zleceniu; poza nim otwierane wyzwalaczem `Pętla wykonawcza ▼` | Nie — pierwszoplanowe |
| Product Builder | Nadrzędny widok procesu budowy produktu | Obszar roboczy, prawa kolumna dominująca | 1 | Widok domyślny obszaru roboczego | **Tak** |
| Architecture Designer | Projektowanie architektury rozwiązania | Kanwa robocza w obszarze roboczym | 1 | Kafel nawigacyjny Product Buildera | Nie |
| Frontend Workspace | Praca nad warstwą frontendową | Obszar roboczy, kolumny równoległe | 1 | Kafel nawigacyjny Product Buildera | Nie |
| Backend Workspace | Praca nad warstwą backendową | Obszar roboczy, kolumny równoległe | 1 | Kafel nawigacyjny Product Buildera | Nie |
| Deployment Panel | Zarządzanie udostępnieniem produktu | Obszar roboczy | 1 | Kafel nawigacyjny Product Buildera | Nie |
| App Catalog | Katalog aplikacji, wtyczek, umiejętności, konektorów i serwerów MCP | Obszar roboczy, siatka kart | 1 | Kafel nawigacyjny; polecenie w Chat Window | Nie |
| Installed Apps Manager | Rejestr, stan włączenia, wersje i cykl życia rozszerzeń | Obszar roboczy, lista | 1 | Kafel nawigacyjny; wyzwalacz `Zainstalowane ▼` w App Catalog | Nie |
| Permissions & Trust Center | Uprawnienia, sandbox, podpisy i pochodzenie | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik uprawnień na karcie rozszerzenia; menu `⋮` pozycji katalogu | Nie |
| Integrations Hub | Konektory i serwery MCP w jednym widoku | Obszar roboczy | 1 | Kafel nawigacyjny | Nie |
| MCP & Connector Console | Inspektor narzędzi, próbne wywołania, log JSON-RPC | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `⋮` pozycji Integrations Hub; przycisk „Inspektor" | Nie |
| Publisher Panel | Pakowanie, podpis, publikacja do prywatnego rejestru | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `Operacje ▼` w Deployment Panel i Installed Apps Manager | Nie |

```
 Makieta całościowa — Moduł: Apps            Dostępność: WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window          │ Obszar roboczy modułu   │ Panel
  nawigacja  │ Użytkownik ↔         │ Product Builder         │ pomocniczy
  modułów    │ Wykonawca            │                         │ (rozszerzenie
             │                      │ [Architektura ▼]        │  boczne)
  Apps       │ ────────────────     │ [Frontend] [Backend]    │
             │ Execution Loop       │ [Wdrożenie] [Katalog]   │  ⋮
             │ Koordynator ↔        │                         │
             │ Wykonawca            │                         │
             │                      │                         │
  [Danaco Console] [CodeStudio] [Portal klienta] [Fable 5]   ☰  ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

### 2.1. Warstwy widoczności w module

Moduł Apps stosuje regułę stopniowego ujawniania funkcjonalności: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Złożoność katalogu rozszerzeń, integracji, uprawnień i narzędzi budowy produktu istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia.

| Warstwa | Zawartość w module Apps | Sposób dostępu |
|---|---|---|
| 1 | Chat Window, Execution Loop Window, aktywne okno wiodące obszaru roboczego (Product Builder, Architecture Designer, Frontend/Backend Workspace, Deployment Panel, App Catalog, Installed Apps Manager, Integrations Hub), pasek kontekstu produktu, wskaźniki stanu budowania i wdrożenia | Widoczne bez interakcji |
| 2 | Wybór wykonawcy (Executor 1 / Executor 2), wybór środowiska wdrożenia, wybór transportu MCP, wybór stosu technologicznego, przełącznik podglądu responsywnego, fasety katalogu, Permissions & Trust Center | Znacznik kontekstowy, ikona, przełącznik; element zwija się samoczynnie po użyciu |
| 3 | Zestawy akcji na rozszerzeniu i wdrożeniu (`Operacje ▼`), warianty strategii wdrożenia, MCP & Connector Console, Publisher Panel, porównanie pozycji katalogu, import i eksport konfiguracji | Menu kebab `⋮`, menu hamburger `☰`, panel popover, lista rozwijana |
| 4 | Skaner manifestu, audyt uprawnień i użycia, log JSON-RPC, log diagnostyczny rozszerzenia, edytor macierzy zakresów, podpisywanie pakietu kluczem wydawcy, tryb administracyjny rejestru | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window. Znaczniki kontekstowe paska modułu — `[Danaco Console] [CodeStudio] [Portal klienta] [Fable 5]` — otwierają po kliknięciu odpowiedni selektor.

---

## 3. Specyfikacja okien operacyjnych — pełny arsenał narzędzi

### 3.1. Chat Window — kanał Użytkownik ↔ Wykonawca

| Pole | Treść |
|---|---|
| Cel | Główne okno komunikacji między Użytkownikiem a Wykonawcą; centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami modułu — budową produktu, katalogiem rozszerzeń, instalacją, uprawnieniami i integracjami |
| Zawartość | Historia rozmowy właściwa bieżącemu produktowi i bieżącemu kontekstowi modułu; strumień odpowiedzi i wyników na żywo; zatwierdzanie i przerywanie działań; wyjaśnienie wyniku i kontekstu |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Pole polecenia w języku naturalnym | Przyjmowanie poleceń wobec dowolnego procesu modułu i prezentacja strumienia odpowiedzi | 1 | Widoczne stale |
| Sterowanie przebiegiem | Zatwierdzenie, przerwanie i korekta działania Wykonawcy | 1 | Przyciski strumienia odpowiedzi |
| Wskaźnik kontekstu | Adnotacja w nagłówku wskazująca okno i etap, którego dotyczy rozmowa | 1 | Widoczny stale |
| Wzmianki (`@`) | Odwołanie do komponentu architektury, pliku frontendu lub backendu, środowiska wdrożeniowego, rozszerzenia, serwera MCP lub konektora | 2 | Znak `@` w polu polecenia |
| Wybór wykonawcy | Przełącznik między pracą własną, Wykonawcą ogólnym oraz Executorem 1 i Executorem 2 środowiska MultitaskingAI | 2 | Znacznik `Wykonawca ▼` w pasku kontekstu |
| Generowanie kodu z opisu | Polecenie przekładane na fragment kodu wprost we właściwym oknie obszaru roboczego | 2 | Polecenie w polu wejściowym |
| Załączniki | Dołączenie specyfikacji, zrzutu ekranu błędu, manifestu lub pliku referencyjnego | 2 | Ikona załącznika |
| Akcje na wiadomości | Kopiowanie fragmentu kodu, zastosowanie wprost w edytorze, rozgałęzienie wątku | 3 | Menu `⋮` wiadomości |
| Asystent doboru rozszerzeń | Rekomendacja i zestawienie rozszerzeń pod opisany cel wraz z wyjaśnieniem różnic i wymaganych uprawnień | 3 | Menu `Operacje ▼` |
| Wyjaśnienie uprawnień i ryzyka | Streszczenie zakresu dostępu, jaki uzyska rozszerzenie, przed jego włączeniem | 3 | Menu `Operacje ▼`; znacznik uprawnień karty rozszerzenia |
| Konfiguracja konektora z opisu | Polecenie języka naturalnego przekładane na wstępną konfigurację konektora przedstawianą do zatwierdzenia | 4 | Polecenie języka naturalnego |
| Diagnoza błędu integracji | Odczyt logu JSON-RPC lub błędu konektora, wskazanie przyczyny i kroku naprawy | 4 | Polecenie języka naturalnego |

**Zachowanie i stany:** okno rekonfiguruje się do kontekstu okna aktualnie otwartego w obszarze roboczym. Stan „oczekuje na wynik równoległy" — gdy polecenie dotyczy jednocześnie frontendu i backendu, widoczne są dwa niezależne wskaźniki postępu.

```
 Makieta — Chat Window w module Apps (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window — Użytkownik ↔ Wykonawca              ⋮
  Produkt „Portal klienta" · kontekst: Backend Workspace
  ───────────────────────────────────────────────
  ▸ Polecenie: „Dodaj endpoint GET /zamowienia z paginacją"
  ▸ Wykonawca: punkt końcowy dodany, testy przeszły
  ───────────────────────────────────────────────
  ┌───────────────────────────────────────────┐
  │ Polecenie…                              @ │
  └───────────────────────────────────────────┘
  [Wykonawca ▼]  [Fable 5 ▼]              [ Wyślij ⏎ ]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.2. Execution Loop Window — kanał Koordynator ↔ Wykonawca

| Pole | Treść |
|---|---|
| Cel | Prowadzenie pętli wykonawczej na zadaniach modułu Apps: koordynacja, nadzór nad realizacją, orkiestracja działań i kontrola przebiegu procesów budowy produktu oraz obsługi rozszerzeń |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania modułu, kolejka i stan zadań, wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli, sterowanie przebiegiem |
| Waga wizualna | Kolumna sąsiadująca z Chat Window |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Widok zlecenia i dekompozycji | Prezentacja zlecenia modułu — budowy komponentu architektury, implementacji punktu końcowego, instalacji zestawu rozszerzeń, uruchomienia wdrożenia — rozłożonego na zadania | 1 | Widoczny przy aktywnym zleceniu |
| Kolejka zadań modułu | Lista zadań frontendu, backendu, wdrożenia, instalacji i testu integracji wraz ze stanem: oczekujące, w toku, zakończone, ponawiane | 1 | Widoczna przy aktywnym zleceniu |
| Strumień komunikatów sterujących | Wymiana poleceń i raportów między Koordynatorem a Wykonawcą dla zadań modułu | 1 | Widoczny przy aktywnym zleceniu |
| Wskaźniki przebiegu pętli | Liczba iteracji, czas trwania, udział zadań zakończonych powodzeniem | 1 | Widoczne przy aktywnym zleceniu |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia | 1 | Przyciski nagłówka okna |
| Wyniki kontroli jakości | Raport walidatora dla zadań modułu: wynik budowania, testy punktów końcowych, walidacja architektury, walidacja manifestu, test połączenia konektora | 2 | Znacznik `Kontrola ▼` |
| Przydział wykonawców | Podgląd, które zadania modułu realizuje Executor 1, a które Executor 2 | 2 | Znacznik `Wykonawcy ▼` |
| Decyzje o ponowieniu | Historia ponowień zadania wraz z przyczyną i zmienionym parametrem | 3 | Menu `⋮` zadania |
| Korekta zlecenia w toku | Zmiana zakresu zlecenia bez przerywania pętli | 3 | Menu `Operacje ▼` |
| Log przebiegu pętli | Pełny zapis komunikatów sterujących do diagnostyki | 4 | Polecenie języka naturalnego w Chat Window; skrót klawiszowy |

**Zachowanie i stany:** okno otwiera się jako kolumna sąsiadująca z Chat Window w chwili przyjęcia zlecenia przez Koordynatora. Stany pętli: planowanie, wykonanie, kontrola, ponowienie, zakończenie. Stan „ponowienie" oznaczany jest przy zadaniu, którego kontrola jakości zwróciła zastrzeżenie.

```
 Makieta — Execution Loop Window (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Execution Loop — Koordynator ↔ Wykonawca          ⋮
  Zlecenie: „Portal klienta — moduł zamówień"
  ───────────────────────────────────────────────
  Zadania
  ✔  Kontrakt API /zamowienia
  ▶  Backend: implementacja punktu końcowego
  ○  Frontend: widok listy zamówień
  ○  Wdrożenie: staging
  ───────────────────────────────────────────────
  Iteracja 3 · 12 min          [Kontrola ▼] [Wykonawcy ▼]
  [ Wstrzymaj ] [ Wznów ] [ Przerwij ] [ Koryguj ]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.3. Product Builder

| Pole | Treść |
|---|---|
| Cel | Nadrzędny widok procesu budowy produktu |
| Zawartość | Etapy budowy produktu od architektury po wdrożenie i publikację |
| Waga wizualna | Prawa kolumna dominująca — obszar roboczy modułu |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Śledzenie etapów | Tracker etapów: architektura → frontend i backend (równolegle) → wdrożenie, z oznaczeniem etapu bieżącego | 1 | Widoczny stale |
| Nawigacja do pozostałych okien | Kafle prowadzące do Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel, App Catalog, Integrations Hub z licznikiem otwartych spraw | 1 | Widoczne stale |
| Panel kondycji projektu | Zbiorczy wskaźnik stanu produktu: liczba otwartych zadań, błędów budowania, ostatnie wdrożenie | 1 | Widoczny stale |
| Metadane produktu | Nazwa, opis, platformy docelowe (webowa, mobilna, desktopowa), repozytorium źródłowe | 2 | Znacznik kontekstowy produktu |
| Przegląd przypisania wykonawców | Widok, które etapy realizuje użytkownik, a które Executor 1 i Executor 2 środowiska MultitaskingAI | 2 | Znacznik `Wykonawcy ▼` |
| Lista kamieni milowych | Wykaz głównych celów projektu z terminami i statusem realizacji | 2 | Znacznik `Kamienie milowe ▼` |
| Widżety integracji | Stan powiązań z modułami Design, Developer, Terminal, Agents i Automations | 3 | Menu `☰` nagłówka |
| Centrum dokumentacji | Odnośniki do dokumentacji technicznej produktu, generowanej i aktualizowanej w toku pracy | 3 | Menu `☰` nagłówka |
| Dziennik wydań | Lista kolejnych wersji produktu wraz z notatkami wydania | 3 | Menu `Operacje ▼` |
| Widok całościowej osi czasu | Chronologia zdarzeń projektu — zmiany architektury, zatwierdzenia zmian, wdrożenia — w jednym zestawieniu | 3 | Menu `Operacje ▼` |

**Zachowanie i stany:** punkt wejścia integrujący pozostałe okna modułu, otwierany po wybraniu produktu. Stan „projekt nowy" — zachęta do rozpoczęcia od Architecture Designer.

```
 Makieta — Product Builder (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window        │ Product Builder                          ☰  ⋮
  Użytkownik ↔       │ Produkt: „Portal klienta"
  Wykonawca          │ ─────────────────────────────────────────────
                     │ Architektura ─► Frontend/Backend ─► Wdrożenie
  ─────────────────  │      ✔              ▶ w toku            ○
  Execution Loop     │ ─────────────────────────────────────────────
  Koordynator ↔      │ ┌────────────┬────────────┬────────────┐
  Wykonawca          │ │Architecture│  Frontend  │  Backend   │
                     │ │  gotowa    │  3 zadania │  5 zadań   │
                     │ ├────────────┼────────────┼────────────┤
                     │ │ Deployment │ App Catalog│Integrations│
                     │ │  v1.2 live │  24 pozycje│  6 aktywnych│
                     │ └────────────┴────────────┴────────────┘
                     │ Kondycja: 8 zadań · 0 błędów budowania
                     │ [Wykonawcy ▼] [Kamienie milowe ▼]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.4. Architecture Designer

| Pole | Treść |
|---|---|
| Cel | Projektowanie architektury rozwiązania |
| Zawartość | Struktura komponentów rozwiązania i ich zależności |
| Waga wizualna | Prawa kolumna dominująca — kanwa robocza |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kanwa diagramu komponentów | Przeciągane węzły: frontend, backend, baza danych, API, integracja zewnętrzna | 1 | Widoczna stale |
| Linie zależności | Rysowanie powiązań i kierunku przepływu danych między komponentami | 1 | Widoczne stale |
| Selektor stosu technologicznego | Przypisanie technologii do każdego komponentu | 2 | Kliknięcie węzła; znacznik `Stos ▼` |
| Panel kontraktu API | Definiowanie punktów końcowych, metod i struktur danych wymienianych między frontendem a backendem | 2 | Kliknięcie etykiety „API" na linii zależności |
| Walidacja architektury | Wykrycie komponentów bez połączeń i brakujących zależności; ostrzeżenie sygnalizowane przed etapem implementacji, nieblokujące dalszej pracy | 2 | Przycisk `Waliduj` |
| Szablony architektury | Wzorce: monolit, mikroserwisy, architektura bezserwerowa — jako punkt wyjścia do dalszej edycji | 3 | Menu `Operacje ▼` |
| Wizualizacja przepływu danych | Podgląd kierunku i rodzaju danych przepływających między komponentami | 3 | Menu `⋮` kanwy |
| Wersjonowanie architektury | Historia kolejnych wersji diagramu i porównanie zmian | 3 | Menu `Operacje ▼` |
| Adnotacje | Notatki projektowe przypięte do komponentów lub połączeń | 3 | Menu kontekstowe węzła |
| Eksport diagramu | Zapis jako obraz lub dokument do dokumentacji technicznej | 3 | Menu `Operacje ▼` |

**Zachowanie i stany:** stanowi punkt wyjścia procesu — poprzedza pracę w Frontend Workspace i Backend Workspace. Stan „architektura zatwierdzona" odsłania zadania pochodne w Product Builderze, bez ograniczania dalszej edycji.

```
 Makieta — Architecture Designer (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window        │ Architecture Designer — „Portal klienta"   ⋮
  Użytkownik ↔       │ [Waliduj]                        [Operacje ▼]
  Wykonawca          │ ─────────────────────────────────────────────
                     │  ┌──────────┐      ┌──────────┐   ┌────────┐
  ─────────────────  │  │ Frontend │─API─►│ Backend  │──►│ Baza   │
  Execution Loop     │  │ (React)  │      │ (Node.js)│   │ danych │
  Koordynator ↔      │  └──────────┘      └────┬─────┘   └────────┘
  Wykonawca          │                         │
                     │                  ┌──────▼──────┐
                     │                  │ Integracja  │
                     │                  │ zewnętrzna  │
                     │                  └─────────────┘
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.5. Frontend Workspace

| Pole | Treść |
|---|---|
| Cel | Praca nad warstwą frontendową produktu |
| Zawartość | Kod i zasoby interfejsu użytkownika |
| Waga wizualna | Prawa kolumna dominująca, podzielona na kolumny robocze |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Edytor kodu interfejsu | Edycja z podświetlaniem składni, uzupełnianiem i wsparciem Wykonawcy, analogiczny do Code Editor modułu Developer | 1 | Widoczny stale |
| Drzewo komponentów | Nawigacja po strukturze komponentów interfejsu | 1 | Widoczne stale |
| Podgląd na żywo | Natychmiastowe odświeżanie wyniku pracy równolegle z edycją | 1 | Widoczny stale |
| Wskaźnik budowania i lintowania | Status kompilacji frontendu i wyników statycznej analizy kodu | 1 | Widoczny stale |
| Podgląd responsywny | Symulacja szerokości ekranu: desktop, tablet, telefon | 2 | Grupa ikon przełącznika podglądu |
| Selektor zasobów z Design | Wybór i wstawienie zasobu wizualnego wprost z Assets Panel modułu Design | 2 | Znacznik `Zasoby Design ▼` |
| Edytor motywu i stylów | Zmiana zmiennych stylistycznych (kolory, typografia, odstępy) spójnie z systemem wizualnym produktu | 3 | Menu `Operacje ▼` |
| Mapa routingu | Przegląd i edycja tras oraz widoków aplikacji | 3 | Menu `☰` |
| Inspektor stanu aplikacji | Podgląd bieżącego stanu danych interfejsu w toku działania podglądu na żywo | 3 | Menu `⋮` podglądu |
| Integracja z Git Panel | Współdzielone z modułem Developer zarządzanie zmianami i historią zatwierdzeń warstwy frontendowej | 3 | Menu `Operacje ▼` |
| Integracja z Terminal | Uruchamianie poleceń budowania i narzędzi wiersza poleceń bez opuszczania okna | 3 | Menu `Operacje ▼`; skrót klawiszowy |

**Zachowanie i stany:** korzysta z zasobów wizualnych z Assets Panel modułu Design; praca prowadzona równolegle z Backend Workspace. Stan „konflikt scalania" sygnalizowany w Git Panel.

```
 Makieta — Frontend Workspace (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ Drzewo      │ Edytor kodu       │ Podgląd na żywo
  Użytkownik ↔    │ komponentów │ interfejsu        │ ┌───────────────┐
  Wykonawca       │             │                   │ │   [podgląd]   │
                  │ ▸ App       │                   │ └───────────────┘
  ──────────────  │ ▸ Layout    │                   │ [▭ ▯ ▫]
  Execution Loop  │ ▸ Zamówienia│                   │
  Koordynator ↔   │             │                   │
  Wykonawca       │ Build: ✔ 0 błędów · 2 ostrzeżenia
                  │ [Zasoby Design ▼]   [Operacje ▼]  ☰  ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.6. Backend Workspace

| Pole | Treść |
|---|---|
| Cel | Praca nad warstwą backendową produktu |
| Zawartość | Kod i konfiguracja usług serwerowych |
| Waga wizualna | Prawa kolumna dominująca, podzielona na kolumny robocze |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Edytor kodu usług | Edycja logiki serwerowej z podświetlaniem składni i wsparciem Wykonawcy | 1 | Widoczny stale |
| Eksplorator punktów końcowych API | Lista tras REST i GraphQL wraz z metodami, parametrami i statusem | 1 | Widoczny stale |
| Wskaźnik budowania | Status kompilacji lub uruchomienia usług backendowych | 1 | Widoczny stale |
| Konstruktor zapytań testowych | Wysyłanie zapytań testowych do punktu końcowego z podglądem odpowiedzi, nagłówków i czasu wykonania | 2 | Kliknięcie wiersza punktu końcowego |
| Panel zmiennych środowiskowych | Konfiguracja parametrów usług per środowisko | 2 | Znacznik `Środowisko ▼` |
| Podgląd schematu bazy danych | Wizualizacja tabel i relacji wraz z edycją struktury | 3 | Menu `☰` |
| Podgląd logów na żywo | Strumień logów usług backendowych w czasie rzeczywistym | 3 | Menu `⋮` |
| Mapa zależności usług | Wizualizacja powiązań między usługami backendu, spójna z Architecture Designer | 3 | Menu `Operacje ▼` |
| Powiązanie z kolejkami | Przejście do zadań w tle powiązanych z modułem Automations | 3 | Menu `Operacje ▼` |
| Integracja z Terminal | Uruchamianie poleceń serwerowych i migracji bez opuszczania okna | 3 | Menu `Operacje ▼`; skrót klawiszowy |

**Zachowanie i stany:** praca prowadzona równolegle z Frontend Workspace. Stan „usługa zatrzymana" sygnalizowany na mapie zależności usług.

```
 Makieta — Backend Workspace (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ Punkty końcowe │ Edytor kodu usług │ Wynik zapytania
  Użytkownik ↔    │ GET /zamowienia│                   │ 200 OK · 84 ms
  Wykonawca       │ POST /klienci  │                   │
                  │ GET /faktury   │                   │
  ──────────────  │                │                   │
  Execution Loop  │                │                   │
  Koordynator ↔   │                │                   │
  Wykonawca       │ Build: ✔ usługi uruchomione
                  │ [Środowisko ▼]      [Operacje ▼]  ☰  ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.7. Deployment Panel

| Pole | Treść |
|---|---|
| Cel | Zarządzanie udostępnieniem produktu |
| Zawartość | Ustawienia wdrożenia, status publikacji |
| Waga wizualna | Prawa kolumna dominująca — obszar roboczy |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Panel stanu produktu | Dostępność usługi, czas nieprzerwanego działania, wynik ostatniego sprawdzenia kondycji | 1 | Widoczny stale |
| Wdrożenie jednym kliknięciem | Uruchomienie procesu wdrożenia bieżącej wersji produktu do wybranego środowiska | 1 | Przycisk `🚀 Wdróż` |
| Historia wdrożeń | Chronologiczna lista wdrożeń z wersją, czasem trwania i wynikiem | 1 | Widoczna stale |
| Selektor środowiska | Środowisko deweloperskie, testowe i produkcyjne — dowolna liczba środowisk konfigurowanych przez użytkownika | 2 | Znacznik `Środowisko ▼` |
| Strategia wdrożenia | Ustawienie konfiguracyjne sposobu wprowadzenia nowej wersji: natychmiastowe, etapowe lub równoległe | 2 | Znacznik `Strategia ▼` |
| Cofnięcie wdrożenia | Powrót do poprzednio wdrożonej wersji, wykonywany od razu; ustawienie `deployment.rollback.confirm` włącza potwierdzenie, a po wykonaniu dostępne jest cofnięcie samej operacji | 2 | Przycisk `Cofnij do tej wersji` w wierszu historii |
| Podgląd logów wdrożenia na żywo | Strumień zdarzeń procesu wdrożeniowego w czasie rzeczywistym | 3 | Menu `⋮` wiersza wdrożenia |
| Edytor notatek wydania | Opis zmian towarzyszący każdemu wdrożeniu | 3 | Menu `Operacje ▼` |
| Lista artefaktów budowania | Pliki wynikowe procesu budowania, gotowe do wdrożenia i do pakowania | 3 | Menu `Operacje ▼` |
| Przekazanie do Publisher Panel | Skierowanie artefaktu wdrożenia do pakowania i publikacji jako rozszerzenie | 3 | Menu `Operacje ▼` |
| Ustawienia domeny i DNS | Konfiguracja adresu docelowego produktu | 3 | Menu `☰` |
| Zarządzanie zmiennymi i danymi dostępowymi środowiska | Konfiguracja wartości właściwych każdemu środowisku wdrożeniowemu | 3 | Menu `☰` |
| Ustawienia skalowania | Liczba instancji usługi i reguły automatycznego skalowania | 4 | Tryb administracyjny; polecenie języka naturalnego |
| Powiązanie z Automations | Wpięcie wdrożenia jako kroku automatyki wyzwalanego zdarzeniem | 4 | Konfiguracja powiązania komponentu |

**Zachowanie i stany:** aktywny po zakończeniu prac w Frontend Workspace i Backend Workspace. Stany wdrożenia: przygotowywane, w toku, zakończone sukcesem, zakończone błędem, wycofane.

```
 Makieta — Deployment Panel (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ Deployment Panel — „Portal klienta"        ☰  ⋮
  Użytkownik ↔    │ [Środowisko ▼] [Strategia ▼]      [ 🚀 Wdróż ]
  Wykonawca       │ ─────────────────────────────────────────────
                  │ Status: ● aktywne · v1.2.0 · dostępność 99,98%
  ──────────────  │ ─────────────────────────────────────────────
  Execution Loop  │ Historia wdrożeń
  Koordynator ↔   │ v1.2.0  2026-08-05 14:20  ✔  [Cofnij do tej wersji] ⋮
  Wykonawca       │ v1.1.3  2026-07-29 09:05  ✔                         ⋮
                  │ v1.1.2  2026-07-22 11:40  ✕                         ⋮
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.8. App Catalog

| Pole | Treść |
|---|---|
| Cel | Katalog aplikacji, wtyczek, umiejętności, konektorów i serwerów MCP z obu źródeł rejestru |
| Zawartość | Siatka kart pozycji katalogu z opisem, wersją, źródłem i stanem |
| Waga wizualna | Prawa kolumna dominująca — obszar roboczy |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Siatka kart katalogu | Jeden widok wszystkich rozszerzeń rejestru z obu źródeł (Danaco Plugin, Personal) | 1 | Widoczna stale |
| Wyszukiwarka pełnotekstowa | Szukanie po nazwie, opisie, kategorii, udostępnianych narzędziach i znacznikach, z podpowiedziami | 1 | Pole wyszukiwania nagłówka |
| Przycisk „Zainstaluj / Włącz" | Uruchomienie instalacji lub zmiana stanu włączenia pozycji | 1 | Przycisk karty |
| Filtry i fasety | Zawężanie po rodzaju, źródle, stanie, transporcie MCP, dostawcy konektora i poziomie uprawnień | 2 | Znacznik `Filtry ▼` |
| Przełącznik widoku | Zmiana prezentacji między siatką kart a listą | 2 | Ikona przełącznika widoku |
| Karta szczegółów rozszerzenia | Pełna metryka: opis, wersja, dziennik zmian, wykaz udostępnianych narzędzi i zasobów, wymagane uprawnienia, zależności, pochodzenie i podpis | 2 | Kliknięcie karty pozycji |
| Kolekcje kuratorskie | Nazwane, oznaczone kolorem zestawy rozszerzeń instalowane i aktywowane grupowo | 3 | Menu `☰` nagłówka |
| Zestawienie porównawcze | Zestawienie 2–4 rozszerzeń tego samego rodzaju w tabeli cech | 3 | Menu `⋮` karty → „Dodaj do zestawienia" |
| Rekomendacje Wykonawcy | Wskazanie rozszerzeń pod opisane zadanie wraz z uzasadnieniem wyboru | 3 | Menu `Operacje ▼`; polecenie w Chat Window |
| Prywatny rejestr organizacji | Wewnętrzny katalog rozszerzeń zbudowanych i opublikowanych w organizacji, prezentowany obok Danaco Plugin | 3 | Znacznik źródła w pasku kontekstu |

**Zachowanie i stany:** punkt wejścia strony dystrybucji i konsumpcji. Stan pozycji: dostępna, zainstalowana wyłączona, zainstalowana włączona, dostępna aktualizacja.

```
 Makieta — App Catalog (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ App Catalog                    [ Szukaj… ]  ☰  ⋮
  Użytkownik ↔    │ [Filtry ▼]  [Zainstalowane ▼]      [▦ ▤]
  Wykonawca       │ ─────────────────────────────────────────────
                  │ ┌────────────┬────────────┬────────────┐
  ──────────────  │ │ Konektor   │ Serwer MCP │ Umiejętność│
  Execution Loop  │ │ CRM        │ Repozytoria│ Faktury    │
  Koordynator ↔   │ │ ● włączony │ ○ wyłączony│ ● włączona │
  Wykonawca       │ │ [Otwórz] ⋮ │ [Włącz]  ⋮ │ [Otwórz] ⋮ │
                  │ └────────────┴────────────┴────────────┘
                  │ [Danaco Plugin] [Personal] [Rejestr organizacji]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.9. Installed Apps Manager

| Pole | Treść |
|---|---|
| Cel | Rejestr zainstalowanych rozszerzeń, ich stan włączenia, wersje i cykl życia |
| Zawartość | Lista rozszerzeń z rodzajem, źródłem, wersją i stanem |
| Waga wizualna | Prawa kolumna dominująca — obszar roboczy |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista zainstalowanych rozszerzeń | Rodzaj, źródło, wersja i stan każdej pozycji rejestru | 1 | Widoczna stale |
| Przełącznik stanu włączenia | Natychmiastowe, odwracalne włączenie i wyłączenie bez restartu serwera; wyłączenie nie usuwa pozycji z rejestru | 1 | Przełącznik wiersza |
| Wskaźnik dostępnej aktualizacji | Oznaczenie pozycji z nowszą wersją w rejestrze wraz z przyciskiem aktualizacji | 1 | Widoczny w wierszu |
| Instalacja Personal | Wskazanie pliku lub pakietu, przesłanie kanałem WebSocket, zapis w katalogu użytkownika i rejestracja encji, potwierdzenie; nowo zainstalowane rozszerzenie pozostaje wyłączone | 2 | Przycisk `Zainstaluj plik` |
| Menedżer wersji | Śledzenie wersji, przypinanie wersji i oznaczanie zgodności semantycznej | 2 | Znacznik `Wersja ▼` wiersza |
| Cofnięcie do wcześniejszej wersji | Powrót do poprzednio zainstalowanej wersji rozszerzenia Personal, wykonywany od razu; ustawienie `extension.rollback.confirm` włącza potwierdzenie, a po wykonaniu dostępne jest cofnięcie samej operacji | 3 | Menu `⋮` wiersza |
| Instalacja z manifestu zestawu | Zainstalowanie i skonfigurowanie całego zestawu rozszerzeń z jednego pliku definicji, odtwarzające środowisko | 3 | Menu `Operacje ▼` |
| Import i eksport konfiguracji | Wyeksportowanie stanu i konfiguracji zestawu rozszerzeń oraz odtworzenie na innej instancji, bez sekretów w jawnym eksporcie | 3 | Menu `Operacje ▼` |
| Dziennik cyklu życia | Chronologia zdarzeń instalacji, aktualizacji, włączeń, wyłączeń i cofnięć wersji per rozszerzenie | 3 | Menu `⋮` wiersza |
| Tryb administracyjny rejestru | Operacje zbiorcze na rejestrze i podgląd zapisu kanonicznego | 4 | Tryb administracyjny; polecenie języka naturalnego |

**Zachowanie i stany:** każda zmiana stanu włączenia jest rozgłaszana na wszystkie urządzenia sesji (rozdz. 1).

```
 Makieta — Installed Apps Manager (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ Installed Apps Manager                     ☰  ⋮
  Użytkownik ↔    │ [Zainstaluj plik]  [Operacje ▼]
  Wykonawca       │ ─────────────────────────────────────────────
                  │ Nazwa            Rodzaj     Źródło     Stan
  ──────────────  │ Konektor CRM     konektor   Danaco   [●] v2.1 ⋮
  Execution Loop  │ Serwer MCP Repo  serwer MCP Personal [○] v0.9 ⋮
  Koordynator ↔   │ Umiejętność PDF  umiejętność Danaco  [●] v1.4 ⋮
  Wykonawca       │ Wtyczka Raporty  wtyczka    Personal [●] v3.0 ⋮ ↑
                  │ [Wersja ▼]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.10. Permissions & Trust Center

| Pole | Treść |
|---|---|
| Cel | Uprawnienia, izolacja wykonania, podpisy i pochodzenie rozszerzeń |
| Zawartość | Zakres uprawnień, status sandboxu, wynik weryfikacji podpisu, audyt użycia |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Podgląd wymaganych uprawnień | Prezentacja deklarowanego zakresu dostępu (sieć, odczyt i zapis plików) z objaśnieniem `[?]` każdego uprawnienia | 2 | Znacznik uprawnień karty rozszerzenia |
| Weryfikacja podpisu i pochodzenia | Sprawdzenie podpisu cyfrowego i sumy kontrolnej pakietu, oznaczenie źródła: Danaco Plugin, zweryfikowany wydawca, niezweryfikowany Personal | 2 | Widoczna po otwarciu panelu |
| Status sandboxu wykonania | Informacja o izolowanym środowisku uruchomienia kodu rozszerzenia i jego ograniczeniach | 2 | Widoczny po otwarciu panelu |
| Przegląd i nadanie zakresu uprawnień | Ustalenie, do jakich zasobów rozszerzenie podłączone do agenta ma dostęp, spójnie z Permissions Center modułu Agents | 3 | Menu `Operacje ▼` |
| Podgląd macierzy izolacji | Prezentacja, jak rozszerzenie mieści się w profilu izolacji roli agenta na poziomie zasięgu „Rola" | 3 | Menu `⋮` panelu |
| Zakres współdzielenia referencji sekretu | Ustalenie, które rozszerzenia i role mają dostęp do danej referencji sekretu | 3 | Menu `Operacje ▼` |
| Skaner manifestu | Ostrzeżenie o szerokich uprawnieniach, nieznanym wydawcy lub braku podpisu — sygnał, nie brama | 4 | Polecenie języka naturalnego; tryb administracyjny |
| Audyt uprawnień i użycia | Zestawienie, które rozszerzenia mają jakie uprawnienia, kiedy i przez którego agenta były użyte, wraz z wykryciem uprawnień nadmiarowych | 4 | Tryb administracyjny; wyszukiwarka funkcji |

Żadne ostrzeżenie ani brak podpisu nie blokuje instalacji ani włączenia. Kontrola pozostaje po stronie Operatora — przez świadome włączenie i przez zakres uprawnień. Dane dostępowe konektorów i serwerów MCP są **kluczami jawnymi** w warstwie sekretów; moduł operuje na referencjach do sekretów, a faktyczne wprowadzenie hasła lub tokenu pozostaje po stronie Operatora. Rozszerzenia nie wykonują samodzielnie płatności ani przelewów.

```
 Makieta — Permissions & Trust Center (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window   │ App Catalog            │ Permissions & Trust Center   ⋮
  Użytkownik ↔  │                        │ Rozszerzenie: Konektor CRM
  Wykonawca     │ ┌──────────┐           │ ────────────────────────────
                │ │ Konektor │           │ Uprawnienia
  ────────────  │ │ CRM   🔒 │           │ ▸ sieć: domena crm.* [?]
  Execution     │ │ [Otwórz] │           │ ▸ pliki: odczyt /dane  [?]
  Loop          │ └──────────┘           │ ────────────────────────────
  Koordynator ↔ │                        │ Podpis: ✔ zweryfikowany
  Wykonawca     │                        │ Sandbox: ● aktywny
                │                        │ [Operacje ▼]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.11. Integrations Hub

| Pole | Treść |
|---|---|
| Cel | Konektory i serwery MCP w jednym widoku operacyjnym |
| Zawartość | Rejestr serwerów MCP, katalog konektorów, stan zdrowia integracji |
| Waga wizualna | Prawa kolumna dominująca — obszar roboczy |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Wspólna lista integracji | Konektory i serwery MCP z ich rodzajem, stanem i dostawcą | 1 | Widoczna stale |
| Panel zdrowia integracji | Zbiorczy status wszystkich włączonych konektorów i serwerów MCP: dostępność, błędy, opóźnienia | 1 | Widoczny stale |
| Rejestr serwerów MCP | Dodanie serwera MCP (adres, transport, dane dostępowe) jako rozszerzenia rodzaju „serwer MCP" | 2 | Przycisk `Dodaj integrację` |
| Wybór transportu MCP | Podłączenie przez stdio, SSE oraz Streamable HTTP wraz z testem transportu | 2 | Znacznik `Transport ▼` |
| Kreator uwierzytelniania | Prowadzone ustanowienie połączenia: OAuth2, klucz API, token, Basic — z kluczami jawnymi w warstwie sekretów | 2 | Przycisk `Połącz` wiersza konektora |
| Test połączenia | Wywołanie próbne konektora z raportem statusu, kodu i czasu odpowiedzi | 2 | Przycisk `Testuj` wiersza |
| Odkrywanie narzędzi, zasobów i promptów | Pobranie listy `tools`, `resources` i `prompts` udostępnianych przez serwer wraz ze schematami wejścia | 2 | Znacznik `Narzędzia ▼` wiersza serwera |
| Import definicji z OpenAPI i GraphQL | Zbudowanie konektora z opisu API wraz z automatycznym wykryciem operacji | 3 | Menu `Operacje ▼` |
| Odbiornik i weryfikacja webhooków | Rejestracja adresu webhooka, weryfikacja podpisu HMAC zdarzeń przychodzących, podgląd ładunku | 3 | Menu `Operacje ▼` |
| Katalog webhooków wychodzących | Konfiguracja zdarzeń platformy wypychanych do systemów zewnętrznych | 3 | Menu `Operacje ▼` |
| Mapowanie i transformacja danych | Odwzorowanie pól między systemem zewnętrznym a encjami platformy wraz z transformacjami | 3 | Menu `⋮` wiersza |
| Podgląd limitów szybkości | Prezentacja limitów usługi, liczników zużycia i kolejkowania wywołań konektora | 3 | Menu `⋮` wiersza |
| Rejestr referencji sekretów i rotacja | Odwołania do danych dostępowych, przypomnienia o wygaśnięciu tokenów i planowa rotacja kluczy | 3 | Menu `☰` nagłówka |
| Monitor zdrowia serwera MCP | Sprawdzenie dostępności, czasu odpowiedzi, liczby udostępnianych narzędzi i błędów handshake | 3 | Menu `⋮` wiersza serwera |
| Metryki użycia i kosztu | Liczba wywołań, powodzenia i błędy, czas odpowiedzi oraz szacunkowy koszt płatnych integracji w oknie czasu | 4 | Wyszukiwarka funkcji; polecenie języka naturalnego |
| Alerty o awarii integracji | Reguła powiadomienia, gdy konektor lub serwer MCP przestaje odpowiadać albo przekracza próg błędów | 4 | Tryb administracyjny; konfiguracja powiązania z Automations |
| Katalog gotowych serwerów MCP | Kuratorowana lista znanych serwerów MCP z opisem i dodaniem jednym kliknięciem | 3 | Menu `☰` nagłówka |

```
 Makieta — Integrations Hub (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window     │ Integrations Hub          [Dodaj integrację]  ☰  ⋮
  Użytkownik ↔    │ Zdrowie: ● 5 sprawnych · ▲ 1 opóźnienie
  Wykonawca       │ ─────────────────────────────────────────────
                  │ Nazwa           Rodzaj      Stan      Akcje
  ──────────────  │ Konektor CRM    konektor    ● 120 ms  [Testuj] ⋮
  Execution Loop  │ Serwer MCP Repo serwer MCP  ● 40 ms   [Narzędzia ▼] ⋮
  Koordynator ↔   │ Konektor Poczta konektor    ▲ 900 ms  [Testuj] ⋮
  Wykonawca       │ Serwer MCP Pliki serwer MCP ○ wyłączony        ⋮
                  │ [Transport ▼]
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.12. MCP & Connector Console

| Pole | Treść |
|---|---|
| Cel | Inspektor narzędzi, próbne wywołania i diagnostyka protokołu |
| Zawartość | Definicje narzędzi serwera MCP, formularz argumentów, odpowiedzi, log JSON-RPC |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista odkrytych narzędzi, zasobów i promptów | Prezentacja definicji udostępnianych przez serwer wraz ze schematami | 3 | Otwarcie konsoli z menu `⋮` Integrations Hub |
| Formularz argumentów | Pola wejściowe generowane ze schematu `JSON Schema` narzędzia | 3 | Wybór narzędzia z listy |
| Próbne wywołanie | Wykonanie `tools/call` lub testu konektora z podglądem odpowiedzi surowej i sformatowanej | 3 | Przycisk `Wywołaj` |
| Piaskownica testu umiejętności i wtyczki | Uruchomienie umiejętności lub wtyczki na przykładowym wejściu bez podłączania do agenta produkcyjnego | 3 | Zakładka `Piaskownica` |
| Wybór podzbioru narzędzi dla agenta | Wskazanie narzędzi serwera przekazywanych do Connectors Manager modułu Agents | 3 | Przycisk `Przekaż do Agents` |
| Podgląd logu JSON-RPC | Strumień komunikatów żądanie, odpowiedź i notyfikacja między platformą a serwerem MCP | 4 | Polecenie języka naturalnego; skrót klawiszowy |
| Log diagnostyczny rozszerzenia | Strumień zdarzeń i błędów wybranego rozszerzenia, spójny z modułem Diagnostics | 4 | Polecenie języka naturalnego; wyszukiwarka funkcji |

```
 Makieta — MCP & Connector Console (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window   │ Integrations Hub    │ MCP & Connector Console       ⋮
  Użytkownik ↔  │                     │ Serwer: Repozytoria
  Wykonawca     │ Serwer MCP Repo     │ ────────────────────────────
                │ ● 40 ms             │ Narzędzia
  ────────────  │ [Narzędzia ▼]  ⋮    │ ▸ repo.search
  Execution     │                     │ ▸ repo.read_file
  Loop          │                     │ ▸ repo.commit
  Koordynator ↔ │                     │ ────────────────────────────
  Wykonawca     │                     │ Argumenty: [ query ]
                │                     │ [ Wywołaj ]  Odpowiedź: —
 ═══════════════════════════════════════════════════════════════════════════
```

### 3.13. Publisher Panel

| Pole | Treść |
|---|---|
| Cel | Pakowanie, podpisywanie i publikacja rozszerzenia do prywatnego rejestru organizacji |
| Zawartość | Kreator pakietu, edytor manifestu, walidator kontraktu, dziennik wydań |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kreator pakietu rozszerzenia | Zapakowanie produktu zbudowanego w module lub artefaktu wdrożenia w dystrybuowalne rozszerzenie z manifestem | 3 | Menu `Operacje ▼` w Deployment Panel |
| Edytor manifestu | Uzupełnienie tożsamości (identyfikator, nazwa, wersja), deklaracji narzędzi, wymaganych uprawnień i zależności, z objaśnieniami `[?]` | 3 | Zakładka `Manifest` |
| Walidator zgodności z kontraktem | Sprawdzenie, czy pakiet spełnia jednolity kontrakt rozszerzenia (tożsamość, stan, zakres); ostrzeżenia nieblokujące | 3 | Przycisk `Waliduj` |
| Dziennik wydań | Notatki kolejnych wersji rozszerzenia powiązane z pakietami | 3 | Zakładka `Wydania` |
| Publikacja do prywatnego rejestru | Umieszczenie rozszerzenia w wewnętrznym rejestrze organizacji, dostępnym w App Catalog obok Danaco Plugin | 3 | Przycisk `Opublikuj` |
| Podpisywanie pakietu | Nadanie podpisu i sumy kontrolnej wydawcy przed publikacją | 4 | Tryb administracyjny; klucz wydawcy z warstwy sekretów |

```
 Makieta — Publisher Panel (stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  Chat Window   │ Deployment Panel    │ Publisher Panel               ⋮
  Użytkownik ↔  │                     │ Pakiet: portal-klienta 1.2.0
  Wykonawca     │ v1.2.0 ✔            │ ────────────────────────────
                │ [Operacje ▼]        │ [Pakiet] [Manifest] [Wydania]
  ────────────  │                     │ Identyfikator: portal-klienta
  Execution     │                     │ Wersja: 1.2.0
  Loop          │                     │ Uprawnienia: sieć, pliki [?]
  Koordynator ↔ │                     │ ────────────────────────────
  Wykonawca     │                     │ [ Waliduj ]   [ Opublikuj ]
 ═══════════════════════════════════════════════════════════════════════════
```

---

## 4. Katalog elementów interfejsu

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole polecenia | Obszar wprowadzania polecenia w języku naturalnym | Sterowanie wszystkimi procesami modułu | Element duży, lewa kolumna | 1 | Widoczne bez interakcji | domyślny, aktywny, oczekiwanie na wynik | Wysyła polecenie do Wykonawcy, strumień odpowiedzi wypełnia okno | Chat Window |
| Kolejka zadań pętli | Lista zadań zlecenia z oznaczeniem stanu | Nadzór nad przebiegiem pętli wykonawczej | Element średni, kolumna sąsiadująca | 1 | Widoczna przy aktywnym zleceniu | oczekujące (○), w toku (▶), zakończone (✔), ponawiane (↻) | Kliknięcie zadania odsłania komunikaty sterujące jego dotyczące | Execution Loop Window |
| Sterowanie przebiegiem pętli | Grupa przycisków wstrzymania, wznowienia, przerwania i korekty | Kontrola realizacji zlecenia | Mała grupa przycisków | 1 | Widoczna przy aktywnym zleceniu | domyślny, wstrzymane, przerwane | Zmienia stan pętli natychmiast i odnotowuje zmianę w strumieniu komunikatów | Execution Loop Window |
| Tracker etapów | Pasek z węzłami etapów | Orientacja w postępie budowy produktu | Element średni, pełna szerokość nagłówka obszaru roboczego | 1 | Widoczny bez interakcji | ukończony (✔), w toku (▶), oczekujący (○) | Kliknięcie etapu przenosi do właściwego okna | Product Builder |
| Kafel nawigacyjny okna | `.dn-karta--interaktywna` z licznikiem | Przejście do jednego z pozostałych okien modułu | Średni panel w siatce | 1 | Widoczny bez interakcji | domyślny, hover, z licznikiem otwartych spraw | Otwiera docelowe okno w obszarze roboczym | Product Builder |
| Znacznik kontekstowy | Lekka plakietka paska kontekstu (`[Portal klienta]`, `[Fable 5]`) | Prezentacja i zmiana kontekstu pracy | Bardzo mały element | 1 | Widoczny bez interakcji | domyślny, aktywny | Kliknięcie otwiera odpowiedni selektor i zwija go po wyborze | Wszystkie okna modułu |
| Węzeł komponentu architektury | Blok na kanwie | Reprezentacja komponentu rozwiązania | Średni blok, kanwa swobodna | 1 | Widoczny bez interakcji | domyślny, zaznaczony, ostrzeżenie walidacji | Kliknięcie odsłania panel właściwości i selektor stosu technologicznego | Architecture Designer |
| Etykieta „API" na linii zależności | Mała podpisana strzałka | Oznaczenie kontraktu komunikacji między komponentami | Bardzo mały element na linii | 2 | Kliknięcie linii zależności | domyślny, najechanie | Otwiera panel kontraktu API | Architecture Designer |
| Przełącznik podglądu responsywnego | Grupa trzech ikon | Symulacja szerokości ekranu w podglądzie na żywo | Bardzo mała grupa przycisków ikonowych | 2 | Ikona nagłówka podglądu | domyślny, wybrany | Zmienia natychmiast szerokość ramki podglądu i zwija się | Frontend Workspace |
| Wskaźnik budowania | `.dn-plakietka--stan` | Sygnalizacja wyniku ostatniej kompilacji | Mała pigułka z kropką | 1 | Widoczny bez interakcji | sukces, ostrzeżenie, błąd, w toku | Kliknięcie rozwija listę błędów i ostrzeżeń | Frontend Workspace, Backend Workspace |
| Wiersz punktu końcowego API | Wiersz `.dn-tabela` z plakietką metody | Lista dostępnych tras API | Wiersz tabeli, waga średnia | 1 | Widoczny bez interakcji | domyślny, wybrany | Kliknięcie ładuje szczegóły do konstruktora zapytań | Backend Workspace |
| Konstruktor zapytania testowego | Formularz metody, adresu i nagłówków | Ręczne wywołanie punktu końcowego API | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie wiersza punktu końcowego | domyślny, wysyłanie, odpowiedź otrzymana | Wysyła zapytanie i prezentuje odpowiedź | Backend Workspace |
| Selektor środowiska wdrożenia | `.dn-select` w pasku kontekstu | Wybór środowiska docelowego wdrożenia | Mały element, zwijany | 2 | Znacznik `Środowisko ▼` | domyślny, przełączony | Zmienia kontekst pozostałych elementów panelu i zwija się | Deployment Panel |
| Przycisk „🚀 Wdróż" | `.dn-btn--zloty` | Najważniejsza akcja okna — uruchomienie wdrożenia | Średni przycisk wypełniony, cień złoty | 1 | Widoczny bez interakcji | domyślny, w toku, ostrzeżenie | Uruchamia wdrożenie i odsłania podgląd logów na żywo; przy nieudanym budowaniu wyświetla ostrzeżenie, pozostając aktywnym | Deployment Panel |
| Przycisk „Cofnij do tej wersji" | `.dn-btn--zarys` | Powrót do wcześniejszej wersji — akcja o podwyższonej wadze | Mały przycisk z obrysem | 2 | Widoczny w wierszu historii | domyślny, w toku | Zawsze aktywny; uruchamia powrót do wersji od razu, a przy włączonym ustawieniu `deployment.rollback.confirm` poprzedza go potwierdzeniem; po wykonaniu dostępne „Cofnij tę zmianę" | Deployment Panel |
| Modal potwierdzenia powrotu do wersji | `.dn-modal` | Potwierdzenie sterowane ustawieniem konfiguracyjnym; nie jest warunkiem wykonania akcji | Duży panel wyśrodkowany | 3 | Wywoływany ustawieniem `deployment.rollback.confirm` | otwarty, zamknięty | Przy włączonym ustawieniu wyświetla się przed operacją; potwierdzenie ją uruchamia, anulowanie zamyka bez skutku | Deployment Panel, Installed Apps Manager |
| Wskaźnik dostępności | Mały pierścień lub pasek procentowy | Informacja o dostępności usługi | Bardzo mały element panelu stanu | 1 | Widoczny bez interakcji | w normie, obniżony | Najechanie pokazuje wartość liczbową i okres pomiaru | Deployment Panel |
| Karta pozycji katalogu | `.dn-karta` z nazwą, rodzajem, źródłem i stanem | Prezentacja rozszerzenia w katalogu | Średni panel w siatce | 1 | Widoczna bez interakcji | dostępna, zainstalowana wyłączona, zainstalowana włączona, dostępna aktualizacja | Kliknięcie odsłania kartę szczegółów rozszerzenia | App Catalog |
| Znacznik uprawnień rozszerzenia | Mała plakietka z symbolem zakresu | Sygnalizacja zakresu dostępu żądanego przez rozszerzenie | Bardzo mały element karty | 2 | Kliknięcie znacznika | wąski zakres, szeroki zakres, brak podpisu | Otwiera Permissions & Trust Center jako rozszerzenie boczne | App Catalog, Installed Apps Manager |
| Przełącznik stanu włączenia | `.dn-przelacznik` w wierszu rejestru | Włączenie i wyłączenie rozszerzenia | Mały przełącznik | 1 | Widoczny bez interakcji | włączone, wyłączone, w toku zmiany | Wysyła `extension.toggle`, zmiana rozgłaszana na wszystkie urządzenia | Installed Apps Manager |
| Plakietka zdrowia integracji | `.dn-plakietka--stan` z czasem odpowiedzi | Status konektora lub serwera MCP | Mała pigułka z kropką | 1 | Widoczna bez interakcji | sprawny, opóźnienie, błąd, wyłączony | Kliknięcie odsłania szczegóły ostatnich sprawdzeń kondycji | Integrations Hub |
| Formularz argumentów narzędzia MCP | Pola generowane ze schematu `JSON Schema` | Próbne wywołanie narzędzia serwera MCP | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `⋮` wiersza integracji | domyślny, walidacja argumentów, wywołanie w toku, odpowiedź otrzymana | Wykonuje `tools/call` i prezentuje odpowiedź surową oraz sformatowaną | MCP & Connector Console |
| Menu `Operacje ▼` | Element zbiorczy grupujący akcje okna | Dostęp do pełnej listy akcji bez zajmowania przestrzeni roboczej | Mały element zwinięty | 3 | Kliknięcie elementu zbiorczego | zwinięte, rozwinięte | Rozwija listę akcji; wybór akcji zwija menu | Wszystkie okna modułu |
| Przełącznik wykonawcy | `.dn-select` w pasku kontekstu Chat Window | Wybór, czy polecenie kierowane jest do użytkownika, Wykonawcy ogólnego czy roli MultitaskingAI | Małe pole zwijane | 2 | Znacznik `Wykonawca ▼` | domyślny, wybrany | Zmienia adresata kolejnego polecenia i zwija się | Chat Window, Execution Loop Window, Product Builder |

---

## 5. Przepływy pracy w module

### 5.1. Od architektury do wdrożenia

```
Product Builder — otwarcie lub utworzenie produktu
        │
        ▼
Chat Window — zlecenie Użytkownika
        │
        ▼
Execution Loop Window — Koordynator dekomponuje zlecenie
na zadania modułu i przydziela je Wykonawcom
        │
        ▼
Architecture Designer — projekt komponentów i zależności
        │
        ▼
        ┌─────────────────────┬─────────────────────┐
        ▼                                            ▼
  Frontend Workspace                          Backend Workspace
  (praca równoległa)                          (praca równoległa)
        │                                            │
        └─────────────────────┬─────────────────────┘
                                ▼
                    Deployment Panel
              wybór środowiska → [ 🚀 Wdróż ]
                                │
                                ▼
                 Status publikacji · logi na żywo
                 · panel stanu produktu
```

### 5.2. Pętla wykonawcza na zadaniach modułu

```
Chat Window — Użytkownik zleca zadanie
(„zbuduj widok listy zamówień i punkt końcowy /zamowienia")
        │
        ▼
Execution Loop Window — Koordynator ↔ Wykonawca
  ├─ dekompozycja zlecenia na zadania modułu
  ├─ kolejka: kontrakt API → backend → frontend → wdrożenie
  ├─ komunikaty sterujące i raporty realizacji
  ├─ kontrola jakości: budowanie, testy punktu końcowego,
  │  walidacja architektury
  └─ decyzja: zakończenie zadania albo ponowienie
        │
        ▼
Obszar roboczy — wynik widoczny w Backend Workspace
i Frontend Workspace, stan etapu w Product Builderze
        │
        ▼
Chat Window — Wykonawca przedstawia wynik,
Użytkownik zatwierdza albo koryguje zlecenie
```

### 5.3. Praca równoległa Executor 1 / Executor 2 (MultitaskingAI)

```
Praca nad frontendem i backendem prowadzona przez
jednego Wykonawcę
                              │
          decyzja użytkownika: powiązanie z MultitaskingAI
          (rozdz. 6.3 niniejszego dokumentu)
                              │
                              ▼
        Koordynator (MultitaskingAI) planuje podział pracy
                              │
              ┌───────────────┴───────────────┐
              ▼                                ▼
        Executor 1                       Executor 2
   Frontend Workspace                Backend Workspace
   (realizacja UI)                   (realizacja API i usług)
              │                                │
              └───────────────┬───────────────┘
                                ▼
              Execution Loop Window — scalenie wyników
              i kontrola jakości prowadzona równolegle;
              ostrzeżenia nie wstrzymują przekazania
              do Deployment Panel
```

### 5.4. Instalacja i podłączenie rozszerzenia

```
App Catalog — wyszukanie pozycji, karta szczegółów
        │
        ▼
Installed Apps Manager — instalacja Personal
(wskaż plik → prześlij → rejestracja encji → potwierdzenie)
        │  nowo zainstalowane rozszerzenie pozostaje wyłączone
        ▼
Permissions & Trust Center — podgląd uprawnień,
podpisu i pochodzenia; nadanie zakresu dostępu
        │
        ▼
Installed Apps Manager — włączenie rozszerzenia
(`extension.toggle`, rozgłoszenie na wszystkie urządzenia)
        │
        ▼
Moduł Agents — podłączenie rozszerzenia jako wykonawcy
zadań agenta w Skills / Connectors Manager
```

### 5.5. Podłączenie serwera MCP i próbne wywołanie

```
Integrations Hub — dodanie serwera MCP
(adres, transport stdio/SSE/HTTP, referencja sekretu)
        │
        ▼
Handshake i odkrywanie: tools · resources · prompts
        │
        ▼
MCP & Connector Console — inspektor definicji narzędzia,
formularz argumentów ze schematu, próbne `tools/call`
        │
        ▼
Wybór podzbioru narzędzi → przekazanie do modułu Agents
        │
        ▼
Integrations Hub — monitor zdrowia serwera,
metryki użycia i kosztu
```

### 5.6. Integracja zasobów wizualnych z modułem Design

```
Moduł Design — Assets Panel (zasoby zaakceptowane)
        │  „Wyślij do modułu" → Apps · Frontend Workspace
        ▼
Frontend Workspace — selektor zasobów z Design
        │
        ▼
Zasób osadzony w komponencie interfejsu, widoczny
natychmiast w podglądzie na żywo
```

### 5.7. Cykl wydania z automatycznym wdrożeniem i publikacją

```
Backend/Frontend Workspace — zmiany zatwierdzone (Git Panel)
        │
        ▼
Automations — automatyka wyzwalana zdarzeniem „zatwierdzenie zmian"
        │
        ▼
Deployment Panel — wdrożenie uruchomione automatycznie
        │
        ▼
Execution Monitor (moduł Automations) — status przebiegu wdrożenia
        │
        ▼
Deployment Panel — panel stanu produktu aktualizuje się po
zakończeniu wdrożenia (sukces, błąd, powrót do wersji)
        │
        ▼
Publisher Panel — pakowanie artefaktu, manifest, walidacja
kontraktu, podpis i publikacja do prywatnego rejestru
        │
        ▼
App Catalog — rozszerzenie dostępne obok Danaco Plugin
```

---

## 6. Stany, dane i powiązania z innymi modułami

### 6.1. Dane wykorzystywane przez Wykonawcę w module

| Źródło danych | Okno pochodzenia | Uwaga |
|---|---|---|
| Zlecenie i stan pętli wykonawczej | Execution Loop Window | Dekompozycja, kolejka zadań, wyniki kontroli jakości |
| Architektura rozwiązania | Architecture Designer | Komponenty, zależności, kontrakty API |
| Stan pracy nad frontendem | Frontend Workspace | Kod, komponenty, zasoby wizualne |
| Stan pracy nad backendem | Backend Workspace | Kod usług, schemat danych, konfiguracja |
| Status wdrożenia | Deployment Panel | Historia, logi, stan produkcyjny |
| Rejestr rozszerzeń | App Catalog, Installed Apps Manager | Sześcioatrybutowy model encji Rozszerzenie: Identyfikator, Nazwa, Źródło, Stan włączenia, Konfiguracja, Lokalizacja (`architektura/model-danych.md`, rozdz. 9) |
| Zakresy uprawnień i wynik weryfikacji podpisu | Permissions & Trust Center | Manifest uprawnień, podpis, pochodzenie |
| Definicje narzędzi i schematy MCP | Integrations Hub, MCP & Connector Console | `tools`, `resources`, `prompts` wraz ze schematami argumentów |

### 6.2. Stany produktu i rozszerzenia w module

| Stan | Znaczenie | Gdzie widoczny |
|---|---|---|
| W architekturze | Trwa projektowanie struktury rozwiązania | Product Builder (tracker etapów) |
| W budowie | Trwa równoległa praca nad frontendem i backendem | Product Builder, Frontend/Backend Workspace |
| Gotowy do wdrożenia | Budowanie zakończone, status widoczny w panelu | Deployment Panel |
| Wdrożony | Produkt aktywny w wybranym środowisku | Deployment Panel (panel stanu) |
| Wycofany | Ostatnie wdrożenie cofnięte do wcześniejszej wersji | Deployment Panel (historia wdrożeń) |
| Opublikowany | Produkt zapakowany i umieszczony w prywatnym rejestrze organizacji | Publisher Panel, App Catalog |
| Zainstalowane wyłączone | Rozszerzenie zarejestrowane, stan wyjściowy Personal | Installed Apps Manager |
| Zainstalowane włączone | Rozszerzenie aktywne i dostępne modułom oraz agentom | Installed Apps Manager, App Catalog |
| Integracja sprawna / z ostrzeżeniem / niedostępna | Wynik sprawdzenia kondycji konektora lub serwera MCP | Integrations Hub (panel zdrowia) |

### 6.3. Powiązania konfigurowalne z innymi modułami i środowiskami

```
                    Developer, Terminal
                    (integracja funkcjonalności
                     w jednym procesie budowy)
                              │
                              ▼
                    ┌───────────────────┐
   Design ────────► │       APPS        │ ────────► Agents
  (zasoby wizualne  │ (moduł, ta karta) │  (umiejętności, wtyczki,
   frontendu)       └────────┬──────────┘   konektory, serwery MCP)
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
      MultitaskingAI    Automations     Diagnostics
   (Executor 1 / 2 —   (automatyka     (logi i błędy
    praca równoległa)   zdarzeniowa)     integracji)
```

| Moduł lub środowisko docelowe | Charakter powiązania | Typ | Miejsce ustanowienia |
|---|---|---|---|
| Developer, Terminal | Apps integruje funkcjonalności obu modułów w ramach jednego, ustrukturyzowanego procesu budowy produktu | Konfiguracyjne | Frontend Workspace, Backend Workspace (edytor, Git Panel, konsola) |
| Design | Zasoby wizualne dla warstwy frontendowej pochodzą z modułu Design | Konfiguracyjne | Frontend Workspace — selektor zasobów |
| Agents | Przygotowane i przetestowane umiejętności, wtyczki, konektory i serwery MCP są podłączane w Skills Manager, Connectors Manager i Permissions Center | Konfiguracyjne, jednokierunkowe | MCP & Connector Console — przekazanie podzbioru narzędzi |
| Okno konfiguracji | Apps jest operacyjnym frontem sekcji „rozszerzenia" rejestru; zmiany rozgłaszane na wszystkie urządzenia (`extension.*`, `config.*`) | Konfiguracyjne, dwukierunkowe | Installed Apps Manager; sekcja „rozszerzenia" okna konfiguracji |
| MultitaskingAI | Realizacja złożonych projektów jest naturalnym zastosowaniem środowiska MultitaskingAI, w szczególności podziału na Executor 1 i Executor 2 przy równoległej pracy nad backendem i frontendem; rozszerzenia wchodzą w skład profili ról z izolacją na poziomie „Rola" | Konfiguracyjne | Okno konfiguracji; Chat Window — przełącznik wykonawcy |
| Automations | Wdrożenie, instalacja, aktualizacja i sprawdzenia kondycji są wyzwalane automatyką zdarzeniową; webhooki działają jako wyzwalacze | Konfiguracyjne | Deployment Panel, Integrations Hub — powiązanie z Automations |
| Diagnostics | Logi i błędy integracji zasilają diagnostykę | Konfiguracyjne, jednokierunkowe | MCP & Connector Console — log diagnostyczny rozszerzenia |

Powiązania nie są aktywne domyślnie: praca w module Apps działa w pełni samodzielnie bez połączenia z Developer, Terminal, Design, Agents, Automations, Diagnostics czy MultitaskingAI — każde z tych powiązań jest świadomą decyzją podejmowaną z poziomu okna konfiguracji lub bezpośrednio z okna operacyjnego, w którym ma zastosowanie.

---

## 7. Katalog funkcji i narzędzi

Każda pozycja katalogu podaje nazwę funkcji, jej działanie oraz zależności techniczne — biblioteki, formaty i protokoły. Wszystkie funkcje respektują jednolity kontrakt rozszerzenia (`architektura/rozszerzenia.md`, rozdz. 3) i sześcioatrybutowy model encji Rozszerzenie (`architektura/model-danych.md`, rozdz. 9).

### 7.1. Katalog rozszerzeń

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Katalog rozszerzeń** | Jeden widok wszystkich rozszerzeń rejestru — wtyczki, umiejętności, konektory, serwery MCP — z obu źródeł, z kartą pozycji (opis, wersja, źródło, stan) | `extension.list` (WebSocket); render karty `.dn-karta`; encja Rozszerzenie (`architektura/model-danych.md`, rozdz. 9) |
| **Wyszukiwarka pełnotekstowa katalogu** | Szukanie po nazwie, opisie, kategorii, udostępnianych narzędziach i znacznikach, z podpowiedziami | Indeks pełnotekstowy Go (`blevesearch/bleve`); metadane manifestu |
| **Filtry i fasety** | Zawężanie po rodzaju, źródle, stanie włączenia, transporcie MCP, dostawcy konektora i poziomie uprawnień | Fasety nad indeksem `bleve`; schemat kategorii |
| **Kuratorskie kolekcje i zestawy** | Nazwane, oznaczone kolorem zestawy rozszerzeń instalowane i aktywowane grupowo | Definicja kolekcji `JSON`; wsad grupowy do `extension.toggle` |
| **Karta szczegółów rozszerzenia** | Pełna metryka: opis, wersja, dziennik zmian, wykaz udostępnianych narzędzi i zasobów, wymagane uprawnienia, zależności, pochodzenie i podpis | Manifest rozszerzenia (`JSON`/`YAML`); podpis (rozdz. 7.3) |
| **Zestawienie porównawcze rozszerzeń** | Zestawienie 2–4 rozszerzeń tego samego rodzaju w tabeli cech | Diff metadanych; `sergi/go-diff` dla różnic opisu |
| **Rekomendacje Wykonawcy** | Wskazanie rozszerzeń pod opisane zadanie wraz z uzasadnieniem wyboru | `kanal_modelu` (`architektura/integracja-modeli.md`); kontekst katalogu; walidacja dopasowania do uprawnień |
| **Prywatny rejestr organizacji** | Wewnętrzny katalog rozszerzeń zbudowanych i opublikowanych w organizacji, prezentowany obok Danaco Plugin | Rejestr serwera (Go); artefakty pakietów (rozdz. 7.7) |

### 7.2. Instalacja i cykl życia

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Instalacja Personal** | Wskazanie pliku lub pakietu → przesłanie kanałem WebSocket → zapis w katalogu użytkownika i rejestracja encji → potwierdzenie; nowo zainstalowane rozszerzenie pozostaje wyłączone | `extension.install` (`architektura/rozszerzenia.md`, rozdz. 5.1); przesył WebSocket; katalog użytkownika na serwerze |
| **Przełącznik stanu włączenia** | Natychmiastowe, odwracalne włączenie i wyłączenie bez restartu serwera; wyłączenie nie usuwa pozycji z rejestru | `extension.toggle`; rozgłoszenie zmiany do wszystkich urządzeń (`architektura/kontrakty-komunikacji.md`, rozdz. 1) |
| **Menedżer wersji** | Śledzenie wersji zainstalowanych rozszerzeń, oznaczanie dostępnych aktualizacji, przypinanie wersji | Wersjonowanie semantyczne (`Masterminds/semver`); porównanie z rejestrem |
| **Aktualizacja** | Danaco Plugin — wraz z pakietem serwera; Personal — ponowne przesłanie pod tym samym identyfikatorem, z zachowaniem konfiguracji i stanu włączenia | `extension.install` (ponowne); zachowanie atrybutów Konfiguracja i Stan włączenia |
| **Cofnięcie do wcześniejszej wersji** | Powrót do poprzednio zainstalowanej wersji rozszerzenia Personal, wykonywany od razu; ustawienie `extension.rollback.confirm` włącza potwierdzenie, a po wykonaniu dostępne jest cofnięcie samej operacji | Przechowywanie poprzednich artefaktów; wzorzec spójny z Deployment Panel |
| **Instalacja masowa z manifestu zestawu** | Zainstalowanie i skonfigurowanie całego zestawu rozszerzeń z jednego pliku definicji, odtwarzające środowisko | Manifest zestawu `YAML`/`JSON`; kolejka `extension.install` |
| **Import i eksport konfiguracji rozszerzeń** | Wyeksportowanie stanu i konfiguracji zestawu rozszerzeń oraz odtworzenie na innym serwerze lub instancji | `config.get` / `config.set`; format eksportu `JSON`, bez sekretów w jawnym eksporcie |
| **Dziennik cyklu życia** | Chronologia zdarzeń instalacji, aktualizacji, włączeń, wyłączeń i cofnięć wersji per rozszerzenie | Log jako `artefakt` (`architektura/model-danych.md`, Zał. F); znaczniki czasu zdarzeń |

### 7.3. Uprawnienia, zaufanie i sandbox

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Podgląd wymaganych uprawnień** | Przed włączeniem prezentuje deklarowany zakres dostępu (sieć, odczyt i zapis plików) czytelnie, z objaśnieniem `[?]` każdego uprawnienia | Manifest uprawnień rozszerzenia; zakresy izolacji (`architektura/architektura.md`, rozdz. 9) |
| **Przegląd i nadanie zakresu uprawnień** | Ustala, do jakich zasobów rozszerzenie podłączone do agenta ma dostęp, spójnie z Permissions Center modułu Agents | Permissions Center (`architektura/rozszerzenia.md`, rozdz. 6, 8); zakresy sieć i pliki |
| **Sandbox wykonania rozszerzenia** | Uruchomienie kodu rozszerzenia w izolowanym środowisku o ograniczonym dostępie, spójnie z profilem izolacji roli | Izolacja procesów Go; runtime WASM (`tetratelabs/wazero`); `hashicorp/go-plugin` (proces potomny gRPC) |
| **Weryfikacja podpisu i pochodzenia** | Sprawdzenie podpisu cyfrowego i sumy kontrolnej pakietu, oznaczenie źródła: Danaco Plugin, zweryfikowany wydawca, niezweryfikowany Personal | Podpis Ed25519 (`crypto/ed25519`); sumy `SHA-256` |
| **Podgląd macierzy izolacji rozszerzenia** | Prezentacja, jak rozszerzenie mieści się w profilu izolacji roli agenta na poziomie zasięgu „Rola", bez odrębnej konfiguracji per rozszerzenie | Okno punktów izolacji (`architektura/model-konfiguracji.md`); poziom zasięgu „Rola" (`architektura/koncepcja-platformy.md`, rozdz. 4.5) |
| **Audyt uprawnień i użycia** | Zestawienie, które rozszerzenia mają jakie uprawnienia, kiedy i przez którego agenta były użyte; wykrycie uprawnień nadmiarowych | Log użycia (`artefakt`); korelacja z agentami (`specyfikacje/specyfikacja-agentow.md`) |
| **Skaner manifestu** | Ostrzeżenie o szerokich uprawnieniach, nieznanym wydawcy lub braku podpisu — sygnał, nie brama | Reguły heurystyczne nad manifestem; plakietka `.dn-plakietka--stan` |

### 7.4. Serwery MCP (Model Context Protocol)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rejestr serwerów MCP** | Dodanie serwera MCP (adres, transport, dane dostępowe) jako rozszerzenia rodzaju „serwer MCP"; lista wraz z jego stanem | Konfiguracja MCP jako atrybut Konfiguracja encji (`architektura/model-danych.md`, rozdz. 9); Connectors Manager (`architektura/rozszerzenia.md`, rozdz. 6) |
| **Obsługa transportów** | Podłączenie przez stdio (proces lokalny), SSE oraz Streamable HTTP; wybór i test transportu | Klient MCP w Go (`modelcontextprotocol/go-sdk`, `mark3labs/mcp-go`); JSON-RPC 2.0 |
| **Odkrywanie narzędzi, zasobów i promptów** | Pobranie listy `tools`, `resources` i `prompts` udostępnianych przez serwer wraz ze schematami wejścia | Handshake MCP (`initialize`, `tools/list`, `resources/list`, `prompts/list`); `JSON Schema` |
| **Inspektor MCP** | Interaktywny podgląd definicji narzędzia (nazwa, opis, schemat argumentów) i próbne wywołanie z podglądem odpowiedzi surowej i sformatowanej | `tools/call` JSON-RPC; walidacja argumentów `santhosh-tekuri/jsonschema`; render odpowiedzi |
| **Katalog gotowych serwerów MCP** | Kuratorowana lista znanych serwerów MCP z opisem i dodaniem jednym kliknięciem | Lista `JSON`; metadane serwera |
| **Monitor zdrowia serwera MCP** | Sprawdzenie dostępności, czasu odpowiedzi, liczby udostępnianych narzędzi i błędów handshake | Cykliczne sprawdzenie kondycji (Go cron); status `.dn-plakietka--stan` |
| **Mapowanie narzędzi MCP na agenta** | Wybór, które narzędzia serwera będą dostępne agentowi, przekazywany do Connectors Manager | Wybór podzbioru narzędzi; przekazanie do modułu Agents |
| **Podgląd logu JSON-RPC** | Strumień komunikatów żądanie, odpowiedź i notyfikacja między platformą a serwerem MCP, do diagnostyki | Log kanału MCP; format `JSON-RPC 2.0`; korelacja żądanie↔odpowiedź |

### 7.5. Konektory i integracje zewnętrzne

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Katalog konektorów** | Lista konektorów do zewnętrznych systemów i usług (repozytoria, systemy biznesowe, usługi sieciowe) z obu źródeł | Rejestr rozszerzeń; rodzaj „konektor" (`architektura/rozszerzenia.md`, rozdz. 1) |
| **Kreator uwierzytelniania** | Prowadzone ustanowienie połączenia: OAuth2 (przekierowanie zgody), klucz API, token, Basic — z kluczami jawnymi w warstwie sekretów | `golang.org/x/oauth2`; warstwa sekretów platformy; wprowadzanie danych logowania po stronie Operatora |
| **Test połączenia** | Wywołanie próbne konektora z raportem statusu, kodu i czasu odpowiedzi | Klient HTTP (`net/http`); parser odpowiedzi; raport `.dn-tabela` |
| **Import definicji z OpenAPI i GraphQL** | Zbudowanie konektora z opisu API (OpenAPI 3, schemat GraphQL) wraz z automatycznym wykryciem operacji | `getkin/kin-openapi`; introspekcja GraphQL; mapowanie operacji |
| **Odbiornik i weryfikacja webhooków** | Rejestracja adresu webhooka, weryfikacja podpisu HMAC zdarzeń przychodzących, podgląd ładunku | HMAC (`crypto/hmac`); endpoint serwera; log zdarzeń |
| **Mapowanie i transformacja danych** | Odwzorowanie pól między systemem zewnętrznym a encjami platformy wraz z transformacjami | Reguły mapowania `JSON`; wyrażenia transformacji; walidacja `JSON Schema` |
| **Podgląd limitów szybkości** | Prezentacja limitów usługi, liczników zużycia i kolejkowania wywołań konektora | Liczniki rdzenia; nagłówki `Retry-After` i `X-RateLimit-*` |
| **Katalog webhooków wychodzących** | Konfiguracja zdarzeń platformy wypychanych do systemów zewnętrznych, jawna i konfigurowalna | Kanał zdarzeń; podpis wychodzący HMAC; integracja z **Automations** |

### 7.6. Wtyczki i umiejętności

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Menedżer wtyczek** | Zarządzanie wtyczkami dostarczającymi nowe funkcje, narzędzia i akcje modułom oraz agentom; podgląd udostępnianych narzędzi | Rodzaj „wtyczka" (`architektura/rozszerzenia.md`, rozdz. 1); Skills Manager (Agents) |
| **Menedżer umiejętności** | Katalog umiejętności — zdefiniowanych sposobów wykonania zadań — z zakresem zastosowania edytowalnym w konfiguracji | Rodzaj „umiejętność"; atrybut Konfiguracja = zakres zastosowania (`architektura/model-danych.md`, rozdz. 9) |
| **Podgląd narzędzi udostępnianych przez wtyczkę** | Rozwinięcie, jakie konkretne narzędzia i akcje wtyczka wnosi i z jakimi parametrami | Manifest wtyczki; schematy narzędzi `JSON Schema` |
| **Piaskownica testu umiejętności** | Uruchomienie umiejętności lub wtyczki na przykładowym wejściu bez podłączania do agenta produkcyjnego | Sandbox (rozdz. 7.3); `kanal_modelu` do próbnego przebiegu |
| **Powiązanie z rolą MultitaskingAI** | Wskazanie, które umiejętności i wtyczki wchodzą w skład profilu roli (Executor, Validator) | Profil roli (`srodowiska/multitaskingai.md`, rozdz. 13); poziom zasięgu „Rola" |

### 7.7. Publikacja, pakowanie i wydawanie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Kreator pakietu rozszerzenia** | Zapakowanie produktu zbudowanego w module lub artefaktu w dystrybuowalne rozszerzenie z manifestem | Archiwum `zip`/`tar.gz`; manifest `JSON`/`YAML`; walidacja pól kontraktu (`architektura/rozszerzenia.md`, rozdz. 3) |
| **Edytor manifestu** | Uzupełnienie tożsamości (identyfikator, nazwa, wersja), deklaracji narzędzi, wymaganych uprawnień i zależności | Schemat manifestu; objaśnienia `[?]`; semver (`Masterminds/semver`) |
| **Podpisywanie pakietu** | Nadanie podpisu i sumy kontrolnej wydawcy przed publikacją do prywatnego rejestru | Ed25519 (`crypto/ed25519`); `SHA-256`; klucz wydawcy w warstwie sekretów |
| **Publikacja do prywatnego rejestru** | Umieszczenie rozszerzenia w wewnętrznym rejestrze organizacji, dostępnym w katalogu obok Danaco Plugin | Rejestr serwera (Go); artefakt pakietu; wersjonowanie |
| **Dziennik wydań** | Notatki kolejnych wersji rozszerzenia powiązane z pakietami | Format `Markdown`; powiązanie z semver |
| **Walidator zgodności z kontraktem** | Sprawdzenie, czy pakiet spełnia jednolity kontrakt (tożsamość, stan, zakres) przed publikacją; ostrzeżenia nieblokujące | Reguły kontraktu (`architektura/rozszerzenia.md`, rozdz. 3); raport walidacji |

### 7.8. Obserwowalność, koszt i zdrowie integracji

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Panel zdrowia integracji** | Zbiorczy status wszystkich włączonych konektorów i serwerów MCP: dostępność, błędy, opóźnienia | Sprawdzenia kondycji (Go cron); agregacja statusów; `.dn-plakietka--stan` |
| **Metryki użycia rozszerzeń** | Liczba wywołań, powodzenia i błędy, czas odpowiedzi per rozszerzenie i narzędzie, w oknie czasu | Liczniki rdzenia; szereg czasowy; wykresy `.dn-*` |
| **Śledzenie kosztu wywołań** | Szacunkowy koszt płatnych konektorów i serwerów oraz wywołań modelu związanych z rozszerzeniem | Liczniki użycia; cennik konfigurowalny; powiązanie z `kanal_modelu` |
| **Alerty o awarii integracji** | Reguła powiadomienia, gdy konektor lub serwer MCP przestaje odpowiadać albo przekracza próg błędów | Progi alertów; kanał powiadomień platformy; integracja z **Automations** |
| **Log diagnostyczny per rozszerzenie** | Strumień zdarzeń i błędów wybranego rozszerzenia do diagnostyki, spójny z modułem Diagnostics | Log jako `artefakt`; powiązanie z modułem **Diagnostics** |

### 7.9. Sekrety i dane dostępowe (klucze jawne)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rejestr referencji sekretów** | Przechowywanie odwołań (kluczy jawnych) do danych dostępowych konektorów i serwerów MCP w warstwie sekretów platformy | Warstwa sekretów platformy; klucze jawne; atrybut Konfiguracja encji |
| **Rotacja i ważność danych dostępowych** | Przypomnienie o wygaśnięciu tokenów OAuth i planowa rotacja kluczy | Znaczniki ważności; harmonogram (Go cron); `golang.org/x/oauth2` (odświeżanie tokenu) |
| **Zakres współdzielenia sekretu** | Ustalenie, które rozszerzenia i role mają dostęp do danej referencji sekretu | Poziomy zasięgu (`architektura/koncepcja-platformy.md`, rozdz. 4.5); Permissions Center |

Moduł operuje na referencjach do sekretów; faktyczne wprowadzenie hasła, tokenu czy danych karty pozostaje po stronie Operatora lub dedykowanego menedżera. Rozszerzenia nie wykonują samodzielnie płatności ani przelewów.

### 7.10. Współpraca człowiek–Wykonawca w module

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Asystent doboru rozszerzeń** | Rekomendacja i zestawienie rozszerzeń pod cel, wyjaśnienie różnic i wymaganych uprawnień | `kanal_modelu`; kontekst katalogu i uprawnień |
| **Konfiguracja konektora z opisu** | Polecenie w języku naturalnym przekładane na wstępną konfigurację konektora, przedstawianą do zatwierdzenia | `kanal_modelu`; mapowanie na pola kreatora (rozdz. 7.5); potwierdzenie Operatora |
| **Wyjaśnienie uprawnień i ryzyka** | Streszczenie zakresu dostępu rozszerzenia i konsekwencji jego włączenia | Manifest uprawnień; `kanal_modelu` |
| **Diagnoza błędu integracji** | Odczyt logu JSON-RPC lub błędu konektora, wskazanie przyczyny i kroku naprawy | Log diagnostyczny (rozdz. 7.8); `kanal_modelu`; powiązanie z **Diagnostics** |

### 7.11. Zależności techniczne — zestawienie

| Obszar | Kluczowa zależność | Biblioteki Go, formaty i protokoły |
|---|---|---|
| Rejestr i cykl życia | Kontrakt komunikacji | `extension.list` / `extension.install` / `extension.toggle`, `config.get` / `config.set` (WebSocket) |
| Serwery MCP | Model Context Protocol | `modelcontextprotocol/go-sdk`, `mark3labs/mcp-go`; JSON-RPC 2.0; transporty stdio, SSE, Streamable HTTP |
| Walidacja narzędzi MCP | Schematy argumentów | `JSON Schema` (`santhosh-tekuri/jsonschema`) |
| Konektory — uwierzytelnianie | OAuth2 i tokeny | `golang.org/x/oauth2`; klucze jawne w warstwie sekretów |
| Konektory — definicje API | OpenAPI i GraphQL | `getkin/kin-openapi`; introspekcja GraphQL |
| Webhooki | Weryfikacja podpisu | `crypto/hmac`; endpoint serwera |
| Sandbox rozszerzeń | Izolacja wykonania | `tetratelabs/wazero` (WASM), `hashicorp/go-plugin` (gRPC), izolacja procesów |
| Podpis i pochodzenie | Integralność pakietu | `crypto/ed25519`, `SHA-256` |
| Pakowanie i wersje | Dystrybucja | `zip`, `tar.gz`, manifest `JSON`/`YAML`, semver (`Masterminds/semver`) |
| Katalog i wyszukiwanie | Indeks pełnotekstowy | `blevesearch/bleve` |
| Różnice metadanych | Porównania pozycji katalogu | `sergi/go-diff` |
| Harmonogram i zdrowie | Cykliczne sprawdzenia | Go cron; integracja z **Automations** |

---

## 8. Punkty sterowania z okna konfiguracji

Zgodnie z Modelem konfiguracji (warstwy: aplikacja, proces, akcja, sesja) i zasadą „brak ustawienia = wartość domyślna = wykonanie", Operator personalizuje moduł Apps bez blokad. Rejestr rozszerzeń pozostaje kanonicznie w sekcji „rozszerzenia" okna konfiguracji (`architektura/rozszerzenia.md`, rozdz. 7); poniższe punkty rozciągają jego sterowanie na moduł.

| Zakres | Co Operator personalizuje | Warstwa / miejsce |
|---|---|---|
| **Kanały komunikacji** | Szerokość kolumny Chat Window, warunek otwarcia Execution Loop Window przy przyjęciu zlecenia, poziom szczegółowości strumienia komunikatów sterujących | Aplikacja / sesja |
| **Pętla wykonawcza modułu** | Zakres zadań dekomponowanych automatycznie, próg ponowienia zadania, zestaw kontroli jakości uruchamianych po zadaniu (budowanie, testy punktów końcowych, walidacja architektury, walidacja manifestu) | Procesy / akcje |
| **Warstwy widoczności** | Przypisanie funkcji do warstw 2–4 dla roli użytkownika, zestaw znaczników kontekstowych paska modułu, widoczność funkcji eksperckich | Aplikacja / konfiguracja roli |
| **Źródła katalogu** | Widoczność Danaco Plugin i Personal, włączenie prywatnego rejestru organizacji, kuratorowane listy serwerów MCP | Aplikacja |
| **Polityka instalacji Personal** | Stan wyjściowy nowych rozszerzeń (wyłączone), potwierdzenie przy instalacji | Aplikacja (`architektura/rozszerzenia.md`, rozdz. 5.2, 8) |
| **Aktualizacje i wersje** | Automatyczne sprawdzanie aktualizacji, przypinanie wersji, ustawienie `extension.rollback.confirm` | Procesy / akcje |
| **Wdrożenia** | Domyślne środowisko, domyślna strategia wprowadzenia nowej wersji, ustawienie `deployment.rollback.confirm`, reguły skalowania | Procesy / akcje |
| **Zaufanie i podpisy** | Wymagany poziom podpisu do oznaczenia „zweryfikowane", reguły skanera manifestu (progi ostrzeżeń, nieblokujące) | Aplikacja / izolacja |
| **Domyślne uprawnienia** | Domyślny zakres nadawany rozszerzeniu przy podłączeniu (sieć, pliki), spójny z profilem izolacji roli | Okno punktów izolacji (macierz, poziom „Rola") |
| **Transporty MCP** | Domyślny transport (stdio, SSE, HTTP), limity handshake, częstotliwość sprawdzeń kondycji | Procesy |
| **Konektory** | Domyślny sposób uwierzytelniania, polityka rotacji tokenów, limity szybkości, obsługa webhooków | Procesy / integracje |
| **Sekrety (klucze jawne)** | Zakres współdzielenia referencji sekretów, przypomnienia o wygaśnięciu | Sekrety / izolacja |
| **Obserwowalność i koszt** | Progi alertów awarii, cennik płatnych integracji, kanał powiadomień (push, e-mail, w aplikacji) | Procesy / integracje |
| **Powiązania modułów** | Włączenie i kierunek powiązań Apps → Agents, MultitaskingAI, Automations, Diagnostics oraz Design, Developer, Terminal | Komponenty (`powiazanie_komponentu`) |
| **Publikacja** | Adres i widoczność prywatnego rejestru, klucz wydawcy, reguły walidacji kontraktu | Aplikacja |

---

## 9. Scenariusze użycia

### 9.1. Budowa produktu od podstaw przez pojedynczy zespół

1. Zespół otwiera moduł Apps w środowisku WorkSpace i tworzy nowy produkt w Product Builderze; polecenie kieruje z Chat Window, a Koordynator otwiera Execution Loop Window z dekompozycją zlecenia.
2. W Architecture Designer zespół projektuje strukturę: frontend, backend, baza danych, jedna integracja zewnętrzna — z jawnie zdefiniowanym kontraktem API między frontendem a backendem.
3. Po zatwierdzeniu architektury praca biegnie równolegle: jedna osoba w Frontend Workspace nad interfejsem, korzystając z zasobów przekazanych z modułu Design, druga w Backend Workspace nad punktami końcowymi API; kolejka zadań i wyniki kontroli jakości widoczne są w Execution Loop Window.
4. Po zakończeniu obu torów pracy zespół otwiera Deployment Panel, wybiera środowisko testowe, uruchamia wdrożenie i weryfikuje wynik w panelu stanu produktu.
5. Po pozytywnej weryfikacji powtarza wdrożenie do środowiska produkcyjnego ze strategią etapową.

### 9.2. Złożony projekt orkiestrowany przez MultitaskingAI

1. Projekt przekracza możliwości efektywnej pracy jednoosobowej — zespół wiąże moduł Apps ze środowiskiem MultitaskingAI.
2. Koordynator planuje podział: Executor 1 obejmuje Frontend Workspace, Executor 2 obejmuje Backend Workspace, pracując według wspólnej architektury zatwierdzonej w Architecture Designer.
3. Execution Loop Window prezentuje kolejkę zadań obu torów, komunikaty sterujące i wskaźniki przebiegu pętli; Validator, wcielony jako Security Auditor, przegląda wyniki równolegle z dalszą pracą — kontrola nie wstrzymuje przekazania do Deployment Panel, a zastrzeżenia sygnalizowane są jako ostrzeżenia.
4. Cały przebieg nadzorowany jest przez Always On Display, z interwencją użytkownika w dowolnym momencie z poziomu Chat Window, Product Buildera lub funkcji globalnej Mobile.

### 9.3. Naprawa błędu produkcyjnego i powrót do wcześniejszej wersji

1. Panel stanu produktu w Deployment Panel sygnalizuje spadek dostępności po najnowszym wdrożeniu.
2. Zespół otwiera podgląd logów wdrożenia oraz — pomocniczo — moduł Diagnostics, aby zlokalizować przyczynę.
3. Do czasu przygotowania poprawki zespół korzysta z przycisku „Cofnij do tej wersji" i od razu wraca do ostatniej stabilnej wersji produktu; potwierdzenie sterowane ustawieniem `deployment.rollback.confirm` nie warunkuje wykonania.
4. Poprawka przygotowywana jest równolegle w Backend Workspace pod nadzorem Execution Loop Window, a po zweryfikowaniu w środowisku testowym zostaje wdrożona ponownie do produkcji.

### 9.4. Podłączenie serwera MCP i przekazanie narzędzi agentowi

1. Operator otwiera Integrations Hub i dodaje serwer MCP, wskazując adres, transport stdio oraz referencję sekretu z warstwy kluczy jawnych.
2. Po handshake platforma odkrywa listę `tools`, `resources` i `prompts`; Operator otwiera MCP & Connector Console i wykonuje próbne wywołanie wybranego narzędzia, sprawdzając odpowiedź surową i sformatowaną.
3. W Permissions & Trust Center Operator przegląda wymagane uprawnienia, wynik weryfikacji podpisu i status sandboxu, po czym nadaje zakres dostępu.
4. Operator wybiera podzbiór narzędzi i przekazuje go do Connectors Manager modułu Agents; Integrations Hub monitoruje odtąd zdrowie serwera, metryki użycia i koszt wywołań.

### 9.5. Publikacja zbudowanego produktu jako rozszerzenia

1. Po wdrożeniu produkcyjnym zespół otwiera z Deployment Panel menu `Operacje ▼` i kieruje artefakt budowania do Publisher Panel.
2. W edytorze manifestu uzupełnia identyfikator, wersję semantyczną, deklarację udostępnianych narzędzi, wymagane uprawnienia i zależności.
3. Walidator zgodności z kontraktem zwraca raport; ostrzeżenia nie wstrzymują publikacji, a zespół rozstrzyga je świadomie.
4. Pakiet zostaje podpisany kluczem wydawcy i opublikowany do prywatnego rejestru organizacji, po czym pojawia się w App Catalog obok pozycji Danaco Plugin i jest instalowany przez pozostałe zespoły.

---

*Koniec dokumentu. Moduł Apps — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*

*© 2026 Danaco Holding Group Sp. z o.o. — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
