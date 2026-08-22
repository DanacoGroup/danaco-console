# Danaco Console — Środowisko WorkSpace

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
| **Tytuł** | Środowisko WorkSpace — powłoka środowiska, komplet okien, warstwy widoczności, nawigacja modułów, karty sesji, układy okien, współdzielenie kontekstu, punkty sterowania |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Materiał wykonawczy do projektowania i budowy interfejsu środowiska WorkSpace: co powstaje, gdzie leży, w jakiej formie i do czego służy — dla każdego okna środowiska, każdego elementu powłoki i każdego przepływu pracy |
| **Zakres** | powłoka środowiska WorkSpace: obszar roboczy jako kontener, okna komunikacji operacyjnej, boczna nawigacja modułów, karty sesji, pasek kontekstu środowiska, układy okien, magistrala kontekstu, funkcje wspólne, personalizacja z okna konfiguracji, wykazy normatywne, skróty klawiszowe, kryteria odbioru |
| **Poza zakresem** | wnętrze modułu Workspace — jego okna operacyjne, projekty, zadania, pamięć kontekstowa, biblioteka projektu i przypisania ekspertów — opisuje [Moduł WorkSpace](../moduly/workspace.md); niniejszy dokument opisuje środowisko i nie powtarza treści dokumentu modułu |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Moduł WorkSpace](../moduly/workspace.md) · [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Architektura techniczna](../architektura/architektura.md) |
| **Prototypy odniesienia** | `design/05-okna/srodowiska/workspace.html` · `design/05-okna/srodowiska/workspace-przedsionek.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `environment`, `workspace`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `design/zasoby/rama.css`, `design/zasoby/prototyp.css` · `design/05-okna/srodowiska/` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie; klucze konfiguracji jawne; każdy stan domyślny jest wyjściowy i zmienialny z okna konfiguracji; brak ustawienia oznacza wartość domyślną, nigdy blokadę |

Stos techniczny odniesienia powłoki środowiska: powłoka kliencka Tauri z silnikiem prezentacji, tokeny `design/zasoby/zetony/zetony.css`, komponenty `.dn-*` z `design/zasoby/css/komponenty.css`. Rdzeń serwera: Go — procesy sesji, trwałość układów, indeks wyszukiwania SQLite FTS5, dopasowanie rozmyte `sahilm/fuzzy`, liczenie tokenów `pkoukk/tiktoken-go`. Trwałość stanu powłoki: SQLite. Klucze konfiguracji powłoki środowiska (`shell.*` — `shell.layout.save`, `shell.tab.group`, `shell.context.handoff`) są kluczami warstwy konfiguracji, odrębnymi od komend kontraktu obszaru `workspace`, który niesie komendy modułu Workspace (rozdz. 9.1) — rozróżnienie jest przedmiotem rozdz. 7.

---

## Spis treści

1. [Przeznaczenie i kontekst środowiska](#1-przeznaczenie-i-kontekst-środowiska)
   - [1.1 Rola środowiska](#11-rola-środowiska)
   - [1.2 Charakter pracy w środowisku](#12-charakter-pracy-w-środowisku)
   - [1.3 Dwa kanały komunikacji operacyjnej środowiska](#13-dwa-kanały-komunikacji-operacyjnej-środowiska)
   - [1.4 Miejsce środowiska w architekturze platformy](#14-miejsce-środowiska-w-architekturze-platformy)
   - [1.5 Kanoniczne rozmieszczenie obszaru roboczego](#15-kanoniczne-rozmieszczenie-obszaru-roboczego)
2. [Moduły dostępne w środowisku](#2-moduły-dostępne-w-środowisku)
   - [2.1 Przedsionek środowiska — widok wejściowy](#21-przedsionek-środowiska--widok-wejściowy)
3. [Komplet okien środowiska](#3-komplet-okien-środowiska)
   - [3.1 Przegląd okien](#31-przegląd-okien)
   - [3.2 Warstwy widoczności w środowisku](#32-warstwy-widoczności-w-środowisku)
4. [Specyfikacja okien](#4-specyfikacja-okien)
   - [4.0 Konwencja opisu](#40-konwencja-opisu)
   - [4.1 Chat Window (kanał Użytkownik ↔ Wykonawca)](#41-chat-window-kanał-użytkownik--wykonawca)
   - [4.2 Execution Loop Window (kanał Koordynator ↔ Wykonawca)](#42-execution-loop-window-kanał-koordynator--wykonawca)
   - [4.3 Boczna nawigacja modułów](#43-boczna-nawigacja-modułów)
   - [4.4 Pasek kart sesji](#44-pasek-kart-sesji)
   - [4.5 Pasek kontekstu środowiska](#45-pasek-kontekstu-środowiska)
   - [4.6 Obszar roboczy jako kontener i menedżer układów okien](#46-obszar-roboczy-jako-kontener-i-menedżer-układów-okien)
   - [4.7 Wyszukiwarka środowiska i paleta poleceń](#47-wyszukiwarka-środowiska-i-paleta-poleceń)
   - [4.8 Menedżer sesji i grup kart](#48-menedżer-sesji-i-grup-kart)
   - [4.9 Magistrala kontekstu i przekazanie artefaktu](#49-magistrala-kontekstu-i-przekazanie-artefaktu)
   - [4.10 Centrum powiadomień i panel szybkiej konfiguracji sesji](#410-centrum-powiadomień-i-panel-szybkiej-konfiguracji-sesji)
5. [Przepływy pracy](#5-przepływy-pracy)
   - [5.1 Rozpoczęcie projektu w środowisku](#51-rozpoczęcie-projektu-w-środowisku)
   - [5.2 Zlecenie realizowane w pętli wykonawczej ponad modułami](#52-zlecenie-realizowane-w-pętli-wykonawczej-ponad-modułami)
   - [5.3 Przeniesienie materiału między modułami](#53-przeniesienie-materiału-między-modułami)
   - [5.4 Praca równoległa nad wieloma kartami jednego dnia](#54-praca-równoległa-nad-wieloma-kartami-jednego-dnia)
   - [5.5 Wznowienie pracy po przerwie](#55-wznowienie-pracy-po-przerwie)
   - [5.6 Nadzór nad procesem spoza stanowiska roboczego](#56-nadzór-nad-procesem-spoza-stanowiska-roboczego)
6. [Stany i powiązania](#6-stany-i-powiązania)
   - [6.1 Stany powłoki środowiska](#61-stany-powłoki-środowiska)
   - [6.2 Stany karty sesji](#62-stany-karty-sesji)
   - [6.3 Powiązania z modułami i funkcjami globalnymi](#63-powiązania-z-modułami-i-funkcjami-globalnymi)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
   - [7.1 Mapa punktów sterowania](#71-mapa-punktów-sterowania)
   - [7.2 Szablon konfiguracji środowiska WorkSpace](#72-szablon-konfiguracji-środowiska-workspace)
8. [Scenariusze użycia](#8-scenariusze-użycia)
   - [8.1 Kancelaria prowadząca kilka spraw równolegle](#81-kancelaria-prowadząca-kilka-spraw-równolegle)
   - [8.2 Zespół produktowy prowadzący projekt od materiału do produktu](#82-zespół-produktowy-prowadzący-projekt-od-materiału-do-produktu)
   - [8.3 Praca ciągła nad projektem cyklicznym](#83-praca-ciągła-nad-projektem-cyklicznym)
   - [8.4 Praca głęboko skupiona nad jednym materiałem](#84-praca-głęboko-skupiona-nad-jednym-materiałem)
   - [8.5 Prowadzenie jednego projektu na wielu urządzeniach](#85-prowadzenie-jednego-projektu-na-wielu-urządzeniach)
9. [Wykazy normatywne i kryteria odbioru](#9-wykazy-normatywne-i-kryteria-odbioru)
   - [9.1 Komendy kontraktu](#91-komendy-kontraktu)
   - [9.2 Żetony projektowe](#92-żetony-projektowe)
   - [9.3 Komponenty interfejsu](#93-komponenty-interfejsu)
   - [9.4 Etykiety interfejsu z prototypu](#94-etykiety-interfejsu-z-prototypu)
   - [9.5 Komunikaty](#95-komunikaty)
   - [9.6 Punkty łamania](#96-punkty-łamania)
   - [9.7 Stany kontrolek](#97-stany-kontrolek)
   - [9.8 Kryteria odbioru](#98-kryteria-odbioru)
10. [Załącznik A. Skróty klawiszowe środowiska](#załącznik-a-skróty-klawiszowe-środowiska)
11. [Załącznik B. Pełny wykaz komend kontraktu środowiska WorkSpace](#załącznik-b-pełny-wykaz-komend-kontraktu-środowiska-workspace)
   - [Obszar `workspace` — 59 komend](#obszar-workspace--59-komend)

---

## 1. Przeznaczenie i kontekst środowiska

### 1.1. Rola środowiska

Środowisko odpowiada na pytanie „w jakim trybie pracuję?”, moduł — na pytanie „jakie zadanie wykonuję?”. WorkSpace jest środowiskiem produktywności, organizacji i realizacji projektów: przekształca zamierzenia w zorganizowane działania i produkty. Istotne są w nim struktura projektu, powtarzalność procesów oraz koordynacja zadań w czasie.

Powłoka środowiska WorkSpace to warstwa nadrzędna wobec modułów — wszystko, co otacza okna operacyjne i pozostaje na ekranie niezależnie od tego, który moduł jest otwarty w bieżącej karcie: rama okna (szyna nawigacji, belka tytułowa, wstążka pozioma z pasmem kart sesji, pasek stanu), boczna nawigacja modułów, obszar roboczy rozumiany jako kontener oraz dwa okna komunikacji operacyjnej — Chat Window i Execution Loop Window. Powłoka pełni funkcję, którą w osobnym systemie operacyjnym pełnią menedżer okien, menedżer sesji i pulpit.

### 1.2. Charakter pracy w środowisku

Praca w WorkSpace jest wielomodułowa i długotrwała. Jeden projekt sięga równocześnie po Research (materiał źródłowy), Studio (redakcja), Design (makiety), Apps (produkt), Workspace (zarządzanie projektem) i Agents (delegacja zadań). Miejscem, w którym rozstrzyga się produktywność, jest powłoka: sposób utrzymywania wielu kart obok siebie, sposób przenoszenia materiału między modułami oraz odtwarzanie układu pracy z poprzedniego dnia.

```
 Środowisko WorkSpace — jeden projekt, wiele modułów w powłoce
 ─────────────────────────────────────────────────────────────
   Research  ─(materiał)─►  Studio  ─(treść)─►  Apps
      │                        │                 │
      └────────── Workspace (projekt nadrzędny) ─┘
                         │
                     Agents (delegacja zadań)

   Powłoka spina karty, przenosi kontekst i utrwala układ okien.
```

### 1.3. Dwa kanały komunikacji operacyjnej środowiska

Środowisko udostępnia w każdej karcie i w każdym module oba kanały komunikacji operacyjnej platformy, w tym samym miejscu układu.

| Kanał | Okno | Rola w środowisku |
|---|---|---|
| Użytkownik ↔ Wykonawca | **Chat Window** | Główne okno komunikacji. Centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami środowiska: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie i przerywanie działań oraz wyjaśnianie wyniku i kontekstu |
| Koordynator ↔ Wykonawca | **Execution Loop Window** | Pierwszoplanowy element środowiska. Prowadzi pętlę wykonawczą, koordynację zadań, nadzór nad realizacją, orkiestrację działań i kontrolę realizacji procesów w obrębie karty sesji, ponad podziałem na moduły |

Role: **Użytkownik** zleca i zatwierdza; **Koordynator** jest komponentem orkiestrującym platformy, dekomponuje zlecenie oraz przydziela i nadzoruje zadania; **Wykonawca** to AI, agent lub system wykonawczy realizujący zadania.

### 1.4. Miejsce środowiska w architekturze platformy

```
STRONA GŁÓWNA
    │  wybór środowiska
    ▼
WorkSpace
    │
    ├── Powłoka środowiska
    │       szyna nawigacji · belka tytułowa · wstążka (pasmo kart sesji) ·
    │       boczna nawigacja modułów · obszar roboczy (kontener)
    │
    ├── Kanały komunikacji operacyjnej
    │       Chat Window (Użytkownik ↔ Wykonawca)
    │       Execution Loop Window (Koordynator ↔ Wykonawca)
    │
    ├── Dziewięć modułów w bocznej nawigacji (rozdz. 2)
    │
    └── Karty sesji
            każda karta = odrębna przestrzeń robocza z własnym
            modułem, układem okien i kontekstem

PONAD CAŁĄ STRUKTURĄ
    Always On Display — pływający agent towarzyszący
    Mobile — nadzór i interwencja zdalna nad procesami środowiska
```

### 1.5. Kanoniczne rozmieszczenie obszaru roboczego

Środowisko stosuje wyłącznie układ pionowy — podział lewa–prawa. Wszystkie okna rozmieszczone są w kolumnach sąsiadujących poziomo; regulacji podlega wyłącznie szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Skrajna lewa | Boczna nawigacja modułów środowiska |
| Lewa, stała, pełna wysokość | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca, otwierana | **Execution Loop Window** — pętla wykonawcza Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu otwartego w bieżącej karcie |
| Kolejne kolumny boczne | Okna pomocnicze powłoki otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego |

```
 Makieta — środowisko WorkSpace w stanie spoczynku
 ══════════╤═════════════════════╤══════════════════════════╤══════════════
  Boczna   │ Chat Window         │ Obszar roboczy modułu    │ Panel
  nawigacja│ Użytkownik ↔        │ (moduł bieżącej karty)   │ pomocniczy
  modułów  │ Wykonawca           │                          │ (rozszerzenie
           │                     │ [Workspace] [Projekt A]  │  boczne)
  Studio   │ Historia rozmowy    │ [Fable 5] [Ultra]  ⋮     │
  Workspace│                     │                          │
  Browser  │ ┌─────────────────┐ │                          │
  Research │ │ pole poleceń    │ │                          │
  Library  │ └─────────────────┘ │                          │
  Roundtab.│         [ ⏎ ]       │                          │
  Design   │ ─────────────────── │                          │
  Apps     │ Execution Loop ▸    │                          │
  Agents   │ Koordynator ↔       │                          │
        ☰  │ Wykonawca           │                          │
 ══════════╧═════════════════════╧══════════════════════════╧══════════════
```

---

## 2. Moduły dostępne w środowisku

Boczna nawigacja środowiska WorkSpace udostępnia dziewięć modułów. Kolejność źródłowa i przynależność do grup podlegają konfiguracji z okna konfiguracji.

| Moduł | Rola w środowisku WorkSpace | Grupa domyślna |
|---|---|---|
| Studio | Redakcja i tworzenie treści projektu | Tworzenie |
| Workspace | Zarządzanie projektem: projekty, zadania, pamięć kontekstowa, biblioteka projektu, agenci projektu — pełna specyfikacja w opracowaniu [Moduł WorkSpace](../moduly/workspace.md) | Organizacja |
| Browser | Praca z materiałem sieciowym w obrębie projektu | Materiał |
| Research | Pozyskiwanie i porządkowanie materiału źródłowego, ustalenia badawcze | Materiał |
| Library | Trwałe repozytorium materiałów środowiska | Materiał |
| Roundtable | Wypracowanie stanowiska przez wiele modeli równolegle | Organizacja |
| Design | Makiety, projekt wizualny, materiały graficzne projektu | Tworzenie |
| Apps | Budowa produktu i aplikacji na podstawie treści i makiet | Tworzenie |
| Agents | Delegacja zadań projektu do agentów użytkownika | Organizacja |

Grupy porządkują listę według roli w pracy projektowej: **Materiał** (Research, Library, Browser), **Tworzenie** (Studio, Design, Apps), **Organizacja** (Workspace, Roundtable, Agents). Ukrycie modułu z listy jest preferencją widoku, nie ograniczeniem dostępu — moduł ukryty pozostaje osiągalny z palety poleceń środowiska i z wyszukiwarki środowiska.

Moduł Workspace jest modułem właściwym środowiska i punktem odniesienia dla projektu nadrzędnego: powłoka kieruje do niego artefakty odbierane z pozostałych modułów i wiąże z nim grupę kart projektu. Zakres funkcjonalny tego modułu — Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager — opisuje dokument modułu.

---

### 2.1. Przedsionek środowiska — widok wejściowy

**Zasada wejścia.** Kliknięcie karty środowiska na stronie głównej otwiera
**przedsionek środowiska**, nie przestrzeń roboczą. Operator nie wskazał
jeszcze modułu, więc żaden moduł nie zostaje otwarty; otwarcie pierwszego
w kolejności byłoby decyzją podjętą za Operatora. Przedsionek pełni wobec
modułów tę samą rolę, którą Centrum dowodzenia pełni wobec środowisk:
pokazuje zbiór i pozwala wybrać, nie wybiera sam.

**Trzy rejony przedsionka.**

| Rejon | Zawartość | Waga |
|---|---|---|
| Szyna sesji („Twoje sesje”) | Wykaz sesji Operatora w tym środowisku, zgrupowany projektami, z filtrem zakresu (czynne · wszystkie · zakończone) i wskaźnikiem pracy w tle. W stopce — wejście do pełnej historii sesji | Kolumna przy lewej krawędzi, waga średnia |
| Płótno — nagłówek środowiska | Godło, nazwa, motto i zdanie o przedmiocie pracy środowiska | Waga główna |
| Płótno — kafle modułów | Dziewięć kafli modułów, każdy z ikoną, nazwą, zdaniem o przeznaczeniu i miarą bieżącego obłożenia | Siatka kafli, waga pośrednia |
| Płótno — listwa działań środowiska | Trzy pozycje: **Nowy projekt**, **Konfiguracja środowiska**, **Ustawienia** — działania dotyczące środowiska jako całości | Listwa jednorzędowa, waga najniższa |

**Czego przedsionek nie zawiera.** Chat Window, Execution Loop Window,
bocznej nawigacji modułów, pasa kart sesji ani żadnego okna operacyjnego. Wszystkie te
elementy należą do przestrzeni roboczej i pojawiają się dopiero po wyborze
modułu. Moduł Workspace występuje w przedsionku **wyłącznie jako kafel**;
jego okna otwierają się po wejściu w moduł, nie wcześniej.

**Dwie drogi wyjścia z przedsionka.** Kliknięcie kafla otwiera przestrzeń
roboczą z tym modułem jako wiodącym i zakłada pierwszą kartę sesji. Kliknięcie
pozycji w szynie sesji wznawia sesję istniejącą — przestrzeń robocza otwiera
się od razu w module tej sesji, wraz z jej układem, historią i kontekstem.

**Powrót do przedsionka.** Przedsionek jest osiągalny z przestrzeni roboczej
przez kontrolkę trybów na pasku (ponowne wskazanie bieżącego środowiska) oraz
przez powrót ze strony głównej. Powrót nie zamyka kart sesji — trwają one
w tle i pozostają widoczne w szynie jako sesje czynne.

**Prototyp odniesienia:** `design/05-okna/srodowiska/workspace-przedsionek.html`.

**Przebieg wejścia do środowiska.** Wybór karty WorkSpace na stronie głównej wysyła komendę kontraktu `environment.enter` (obszar `environment`, rozdz. 9.1); rdzeń zwraca definicję środowiska, wykaz modułów i wykaz sesji karty, na podstawie których powłoka otwiera przedsionek.

```
Operator                Powłoka                       Rdzeń
   │                       │                             │
   │  wybór karty          │                             │
   │  „WorkSpace”          │                             │
   ├──────────────────────►│                             │
   │                       │  environment.enter          │
   │                       │  (environmentId,            │
   │                       │   clientId, sessionId?)     │
   │                       ├────────────────────────────►│
   │                       │                             │
   │                       │  environment · modules      │
   │                       │  sessions · focusedSessionId│
   │                       │◄────────────────────────────┤
   │  przedsionek          │                             │
   │  (szyna sesji,        │                             │
   │   siatka modułów)     │                             │
   │◄──────────────────────┤                             │
   │                       │                             │
```

Legenda: `environment.enter` — komenda kontraktu obszaru `environment`; pola żądania `environmentId`, `clientId` (wymagane), `sessionId` (opcjonalne — puste otwiera przedsionek bez karty wybranej); pola wyniku `environment`, `modules`, `sessions` (wymagane), `focusedSessionId` (opcjonalne). Niepuste `sessionId` pomija przedsionek i otwiera od razu przestrzeń roboczą wskazanej karty w module, do którego karta należy.

---

## 3. Komplet okien środowiska

### 3.1. Przegląd okien

| Okno | Rola w środowisku | Forma wiodąca | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca — centralny punkt pracy i podstawowy mechanizm sterowania procesami środowiska | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji |
| Execution Loop Window | Pętla wykonawcza Koordynator ↔ Wykonawca — koordynacja zadań, nadzór, orkiestracja i kontrola realizacji procesów | Kolumna sąsiadująca z Chat Window, otwierana | 1 | Nagłówek pętli widoczny stale; rozwinięcie kolumny kliknięciem `▸` |
| Boczna nawigacja modułów | Wybór modułu bieżącej karty, wskaźniki pracy w tle | Skrajna lewa kolumna, stała | 1 | Widoczna bez interakcji |
| Pasek kart sesji | Zarządzanie równoległymi przestrzeniami roboczymi | Pasek nad obszarem roboczym | 1 | Widoczny bez interakcji |
| Pasek kontekstu środowiska | Znaczniki kontekstu: środowisko, projekt, moduł, model, wykonawca, polityka kontekstu, procesy w tle | Wiersz znaczników w nagłówku obszaru roboczego | 1 | Widoczny bez interakcji; kliknięcie znacznika otwiera selektor |
| Obszar roboczy jako kontener | Kompozycja okien operacyjnych modułu w kolumnach | Prawa, dominująca kolumna | 1 | Widoczny bez interakcji |
| Panel szybkiej konfiguracji sesji | Model, profil asystenta, motyw, retencja historii, polityka współdzielenia, przypisany projekt | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie znacznika kontekstu w pasku kontekstu |
| Wyszukiwarka środowiska | Odnajdywanie kart, modułów, projektów, materiałów Library, zapisanych sesji i układów | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Pole wyszukiwania wstążki (z filtrem zakresu środowisk) |
| Podgląd polityki kontekstu | Prezentacja zasięgu kontekstu obowiązującego bieżącą kartę przed wysłaniem polecenia | Panel popover przy pasku kontekstu | 2 | Kliknięcie znacznika polityki kontekstu |
| Menedżer układów okien | Zapis, wczytanie i reset układu okien; podział kolumn, dokowanie, tryb skupienia | Panel popover nad obszarem roboczym | 3 | Menu kebab (⋮) nagłówka obszaru roboczego |
| Przełącznik i wyszukiwarka kart | Lista wszystkich kart z filtrem tekstowym i miniaturą, lista ostatnio zamkniętych | Panel popover przy pasku kart | 3 | Menu kebab (⋮) paska kart |
| Menedżer sesji i grup kart | Zapisane przestrzenie projektu, grupy kart, przypięcia | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu kart sesji |
| Magistrala kontekstu | Bufor roboczy na fragmenty, pliki i odnośniki krążące między kartami i modułami | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu hamburger (☰) bocznej nawigacji; upuszczenie materiału na krawędź obszaru |
| Panel przekazania artefaktu | Kierowanie wyniku okna do innego modułu, do Library albo do projektu w module Workspace | Panel popover przy elemencie wyniku | 3 | Menu kontekstowe wyniku |
| Centrum powiadomień | Rejestr zdarzeń środowiska; mechanizm w postaci jednolitej dla całej platformy — [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md), rozdz. 11.6 | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Plakietka powiadomień z licznikiem w pasku kontekstu |
| Paleta poleceń środowiska | Jeden punkt dostępu do wszystkich poleceń powłoki | Nakładka wyszukiwania nad obszarem roboczym | 4 | Skrót `Ctrl/Cmd + K` albo polecenie języka naturalnego w Chat Window |
| Mapa skrótów | Podgląd i przepięcie każdego skrótu powłoki, rozwiązywanie konfliktów | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 | Skrót `Ctrl/Cmd + /`, okno konfiguracji |
| Tryb czystego ekranu | Ukrycie nawigacji, kart i pasków do samego obszaru roboczego | Stan całego środowiska | 4 | Polecenie palety poleceń, skrót klawiszowy |

Osiemnaście okien i elementów powłoki łącznie. Chat Window i Execution Loop Window są oknami wspólnymi wszystkim modułom platformy — w środowisku WorkSpace występują w rekonfiguracji właściwej pracy projektowej opisanej w rozdz. 4.1 i 4.2.

### 3.2. Warstwy widoczności w środowisku

Zasada nadrzędna interfejsu: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Środowisko WorkSpace ujawnia możliwości powłoki stopniowo — zależnie od kontekstu karty, roli użytkownika i wykonywanej czynności. Złożoność środowiska istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, kart, układów, agentów i ustawień nie wpływa na postrzeganą prostotę interfejsu.

| Warstwa | Zawartość w środowisku WorkSpace | Sposób wywołania |
|---|---|---|
| 1 — Zawsze widoczna | Chat Window, nagłówek Execution Loop Window, boczna nawigacja modułów, pasek kart sesji, pasek kontekstu, obszar roboczy modułu, wskaźniki stanu wykonania. Ponad 80% powierzchni interfejsu | Widoczne bez interakcji |
| 2 — Widoczna na żądanie | Wybór modelu i wykonawcy karty, wybór modułu, wybór projektu, polityka kontekstu, retencja historii, motyw, wyszukiwarka środowiska | Znacznik kontekstu, ikona, przełącznik; po użyciu element zwija się samoczynnie |
| 3 — Rozwinięcia kontekstowe | Menedżer układów okien, przełącznik kart, menedżer sesji i grup, magistrala kontekstu, przekazanie artefaktu, centrum powiadomień, menu karty, menu modułu | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 — Funkcje eksperckie | Paleta poleceń środowiska, mapa skrótów, tryb czystego ekranu, pamięć układu per klasa urządzenia, przelanie karty do innego środowiska, diagnostyka kanału sterującego | Skrót klawiszowy, polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, okno konfiguracji |

Mechanizmy ukrywania funkcjonalności stosowane w środowisku:

- **Menu progresywne.** Wybór modelu, wykonawcy i modułu prezentowany jest jako jeden element zwinięty (`Wykonawca ▼`), a lista pozycji rozwija się po kliknięciu.
- **Panele wysuwane.** Menedżer sesji, magistrala kontekstu, centrum powiadomień i wyszukiwarka środowiska działają jako kolumny boczne; po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Operacje na karcie i na układzie okien występują jako jeden element zbiorczy (`Operacje ▼`, `⋮`), którego rozwinięcie zawiera pełną listę akcji.
- **Znaczniki kontekstowe.** Środowisko, projekt, moduł, model i wykonawca występują jako lekkie znaczniki paska kontekstu, na przykład `[Danaco Console] [WorkSpace] [Projekt Alfa] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja środowiska jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest w środowisku wykluczone — ukrycie zmniejsza chaos wizualny i nie utrudnia dostępu.

Makiety w niniejszym dokumencie przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `▸`, `⋮`, `☰`).

---

## 4. Specyfikacja okien

### 4.0. Konwencja opisu

Każde okno opisano celem, umiejscowieniem, zawartością, pełnym arsenałem narzędzi z podaniem warstwy widoczności i sposobu wywołania, zachowaniem i stanami oraz makietą tekstową. Wszystkie makiety przedstawiają układ pionowy (podział lewa–prawa) w stanie spoczynku; regulacji podlega wyłącznie szerokość kolumn. Stany wspólne wszystkim oknom środowiska: trwałość stanu procesu sesji po rozłączeniu klienta, aktualizacja na żywo kanałem WebSocket dla okien monitorujących, objaśnienie kontekstowe `[?]` przy każdym elemencie konfiguracji, odrębny kontekst per karta sesji ze współdzieleniem sterowanym z okna konfiguracji punktów izolacji.

### 4.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Główne okno komunikacji między Użytkownikiem a Wykonawcą (AI, agent, system wykonawczy). Stanowi centralny punkt pracy w środowisku i podstawowy mechanizm sterowania wszystkimi procesami środowiska — poleceniami wydawanymi wobec kart, modułów, układów okien, materiałów i zadań projektu |
| Waga i miejsce | Lewa kolumna, stała, pełna wysokość obszaru roboczego; obecna w każdej karcie i w każdym module środowiska w tym samym miejscu układu |
| Zawartość | Historia rozmowy karty; pole wprowadzania poleceń w języku naturalnym; strumień odpowiedzi i wyników na żywo; sterowanie zatwierdzaniem i przerywaniem działań; wiersz znaczników kontekstu |
| Warstwa | 1 — zawsze widoczna |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kompozytor poleceń | Pole wieloliniowe z formatowaniem Markdown na żywo i skrótami klawiszowymi | 1 | Widoczne bez interakcji |
| Strumień odpowiedzi | Odbiór wyników na żywo kanałem WebSocket, z przerwaniem generowania w toku | 1 | Widoczny bez interakcji |
| Zatwierdzanie i przerywanie działań | Akceptacja wyniku, żądanie korekty, przerwanie działania Wykonawcy | 1 | Przyciski strumienia odpowiedzi |
| Zlecenie realizacji zadania | Przekazanie polecenia do Koordynatora — zlecenie pojawia się w Execution Loop Window jako pozycja dekomponowana na zadania | 1 | Przycisk kompozytora |
| Sterowanie powłoką poleceniem języka naturalnego | Otwarcie modułu w nowej karcie, wczytanie układu, zapis sesji, przełączenie motywu, wejście w tryb czystego ekranu — wydane słownie w polu poleceń | 1 | Treść polecenia |
| Wybór wykonawcy | Przełącznik modelu albo agenta jako adresata polecenia | 2 | Znacznik `Wykonawca ▼` |
| Wybór modelu i poziomu wysiłku | Model i intensywność pracy Wykonawcy dla bieżącej karty | 2 | Znaczniki paska kontekstu |
| Wzmianki kontekstowe (`@`) | Odwołanie do karty, modułu, projektu, materiału Library albo agenta wprost w treści polecenia | 2 | Znak `@` w polu poleceń |
| Załączniki | Dołączanie plików z magistrali kontekstu, z Library albo z urządzenia; podgląd miniatur przed wysłaniem | 2 | Ikona spinacza |
| Wskaźnik zasięgu kontekstu | Adnotacja o polityce kontekstu obowiązującej bieżącą kartę | 2 | Znacznik polityki kontekstu |
| Menedżer szablonów poleceń | Zapisane, wielokrotnego użytku polecenia i fragmenty wstawiane skrótem | 3 | Lista rozwijana `Szablony ▾` |
| Akcje wiadomości | Kopiuj, cytuj, edytuj i wyślij ponownie, rozgałęź wątek, odłóż na magistrali kontekstu, przekaż do modułu | 3 | Menu kontekstowe wiadomości |
| Blok kodu w odpowiedzi | Podświetlanie składni, kopiowanie, przekazanie do modułu docelowego | 3 | Menu kebab (⋮) bloku |
| Wyszukiwanie w historii | Przeszukiwanie treści rozmowy bieżącej karty | 3 | Menu kebab (⋮) nagłówka okna |
| Eksport rozmowy | Zapis wątku jako dokumentu Markdown albo PDF do Library lub do projektu | 3 | Menu kebab (⋮) nagłówka okna |
| Wyjaśnienie wyniku i kontekstu | Rozwinięcie informacji o materiałach, pamięci i instrukcjach użytych przez Wykonawcę | 3 | Odnośnik pod odpowiedzią |
| Diagnostyka kanału sterującego | Podgląd komunikatów WebSocket bieżącej karty | 4 | Paleta poleceń środowiska |

**Zachowanie i stany:** kontekst rozmowy jest odrębny per karta sesji; współdzielenie historii, pamięci albo kontekstu między wskazanymi kartami ustanawia się z panelu szybkiej konfiguracji sesji zgodnie z regułami okna konfiguracji punktów izolacji. Stan ładowania — animowany wskaźnik pracy Wykonawcy. Stan błędu — komunikat w treści strumienia z ponowieniem. Stan rozłączenia — okno prezentuje ostatni znany stan i wznawia strumień po odzyskaniu kanału.

```
 Makieta — Chat Window w środowisku WorkSpace (stan spoczynku)
 ══════════╤══════════════════════════════╤═══════════════════════════════
  Boczna   │ Chat Window                  │ Obszar roboczy modułu
  nawigacja│ Użytkownik ↔ Wykonawca       │
  modułów  │                              │ [WorkSpace] [Projekt Alfa]
           │ Historia rozmowy karty       │ [Studio] [Fable 5] [Ultra]  ⋮
  Studio   │ (przewijana, znaczniki czasu)│
  Workspace│                              │
  Browser  │ [@][📎][Szablony ▾]           │
  Research │ [Wykonawca ▼]                │
  Library  │ ┌──────────────────────────┐ │
  Roundtab.│ │ pole poleceń — Markdown  │ │
  Design   │ └──────────────────────────┘ │
  Apps     │      [ Zatwierdź ] [ ⏎ ]     │
  Agents   │ ──────────────────────────── │
        ☰  │ Execution Loop ▸             │
 ══════════╧══════════════════════════════╧═══════════════════════════════
```

### 4.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Okno pętli wykonawczej prezentujące komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów środowiska. W środowisku WorkSpace pętla obejmuje zadania sięgające wielu modułów jednej karty i jednej grupy kart projektu |
| Waga i miejsce | Kolumna sąsiadująca z Chat Window, otwierana; pełna wysokość obszaru roboczego |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania, kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli, sterowanie przebiegiem |
| Warstwa | 1 — nagłówek pętli i wskaźniki stanu widoczne stale; rozwinięcie kolumny jednym kliknięciem |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Karta bieżącego zlecenia | Treść zlecenia przyjętego z Chat Window, karta i projekt źródłowy, wskazany wykonawca, znacznik czasu przyjęcia | 1 | Widoczna po rozwinięciu kolumny |
| Dekompozycja zlecenia na zadania | Drzewo zadań i podzadań wyprowadzonych przez Koordynatora; każde zadanie z wykonawcą, modułem realizacji, terminem i kryterium ukończenia | 1 | Widoczna po rozwinięciu kolumny |
| Kolejka i stan zadań | Lista zadań w stanach: oczekujące, w realizacji, w kontroli, ukończone, ponawiane, przerwane | 1 | Widoczna po rozwinięciu kolumny |
| Wskaźniki przebiegu pętli | Liczba zadań w każdym stanie, czas trwania przebiegu, tempo realizacji, obłożenie przypisanych wykonawców | 1 | Nagłówek pętli — widoczny także przy kolumnie zwiniętej |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie pętli oraz korekta zlecenia bez utraty dotychczasowych wyników | 1 | Przyciski kolumny pętli |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator → Wykonawca i Wykonawca → Koordynator: przydział zadania, raport postępu, zgłoszenie przeszkody, zwrot wyniku | 2 | Zakładka strumienia w kolumnie pętli |
| Wyniki kontroli jakości | Ocena rezultatu zadania wobec kryterium ukończenia z uzasadnieniem i decyzją o przyjęciu albo ponowieniu | 2 | Rozwinięcie pozycji zadania |
| Orkiestracja międzymodułowa | Przydział zadań pętli do modułów środowiska: materiał z Research, redakcja w Studio, makieta w Design, budowa w Apps, ewidencja w module Workspace | 2 | Rozwinięcie pozycji zadania |
| Decyzje o ponowieniu | Rejestr ponowień zadania z licznikiem prób i zmienionymi parametrami zlecenia | 3 | Menu kebab (⋮) pozycji zadania |
| Odbiór artefaktów | Kierowanie wyników zadań do Library albo do projektu w module Workspace | 3 | Menu kebab (⋮) pozycji zadania |
| Powiązanie z kartami środowiska | Otwarcie karty modułu, w którym realizowane jest wskazane zadanie pętli | 3 | Menu kontekstowe zadania |
| Dziennik przebiegu | Pełny, chronologiczny zapis pętli do podglądu i eksportu jako dokument | 3 | Menu kebab (⋮) nagłówka pętli |
| Nadzór zdalny przez Mobile | Zatwierdzenie kroku, wstrzymanie i wznowienie pętli z urządzenia zdalnego | 4 | Funkcja globalna Mobile, paleta poleceń |
| Harmonogram przebiegu | Wpięcie pętli w automatykę cykliczną modułu Automations | 4 | Paleta poleceń, okno konfiguracji |

**Zachowanie i stany:** okno otwierane jest jako kolumna sąsiadująca z głównym oknem komunikacji i aktualizuje się na żywo kanałem WebSocket. Stan przebiegu utrzymuje się po rozłączeniu klienta — po powrocie okno prezentuje bieżący stan pętli. Stan pusty — brak aktywnego zlecenia, widoczna lista zakończonych przebiegów karty. Stan przeszkody — zadanie zgłoszone jako zablokowane wyróżnione wizualnie, z akcją korekty zlecenia. Stan pracy w tle — wskaźnik przebiegu widoczny na karcie sesji i przy pozycji modułu w bocznej nawigacji.

```
 Makieta — Execution Loop Window (kolumna rozwinięta)
 ══════════╤══════════════════════╤═══════════════════╤════════════════════
  Boczna   │ Chat Window          │ Execution Loop    │ Obszar roboczy
  nawigacja│ Użytkownik ↔         │ Koordynator ↔     │ modułu
  modułów  │ Wykonawca            │ Wykonawca         │
           │ (zwinięte ▸)         │                   │ Zadania i wyniki
  Studio   │                      │ Zlecenie:         │ aktualizowane
  Workspace│                      │ „Przygotuj raport │ na żywo
  Browser  │                      │  tygodniowy”      │
  Research │                      │ ├ 1 Materiał   ✔  │
  Library  │                      │ ├ 2 Analiza    ▶  │
  Roundtab.│                      │ ├ 3 Redakcja   ⏳  │
  Design   │                      │ └ 4 Kontrola   ⏳  │
  Apps     │                      │                   │
  Agents   │                      │ Ponowienia: 0     │
           │                      │ Czas: 04:12       │
        ☰  │                      │ [Wstrzymaj][⋮]    │
 ══════════╧══════════════════════╧═══════════════════╧════════════════════
```

### 4.3. Boczna nawigacja modułów

| Pole | Treść |
|---|---|
| Cel | Wybór modułu bieżącej karty i orientacja w pracy toczącej się w innych kartach środowiska |
| Waga i miejsce | Skrajna lewa kolumna, stała, pełna wysokość |
| Zawartość | Dziewięć pozycji modułów, sekcja przypiętych, grupy zwijane, wskaźniki aktywności i liczniki kart |
| Warstwa | 1 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista modułów | Dziewięć pozycji w kolejności konfigurowalnej | 1 | Widoczna bez interakcji |
| Wskaźnik aktywności modułu | Kropka `.dn-kropka` przy pozycji, gdy w innej karcie trwa w tym module praca w tle | 1 | Widoczny bez interakcji |
| Licznik kart modułu | Plakietka z liczbą kart, w których moduł jest otwarty | 1 | Widoczny bez interakcji |
| Sekcja przypiętych modułów | Moduły najczęściej używane w bieżącym projekcie, nad grupami | 1 | Widoczna, gdy zawiera pozycje |
| Grupy modułów | Zwijane sekcje: Materiał, Tworzenie, Organizacja | 2 | Kliknięcie nagłówka grupy |
| Otwarcie modułu w nowej karcie | `Ctrl/Cmd` z kliknięciem pozycji otwiera moduł w nowej karcie zamiast przeładować bieżącą | 2 | Modyfikator klawiszowy |
| Podgląd modułu | Mini-podgląd zestawu okien i opisu tematycznego modułu przed przeładowaniem karty | 2 | Wskazanie kursorem z przytrzymaniem |
| Zwinięcie panelu do ikon | Zwężenie nawigacji do samych godeł modułów; szerokość przechodzi do obszaru roboczego | 2 | Przełącznik u szczytu panelu |
| Menu modułu | Przypnij, otwórz w nowej karcie, ukryj z listy, przypisz skrót | 3 | Menu kontekstowe pozycji |
| Personalizacja widoczności pozycji | Ukrycie modułów nieużywanych w danym profilu pracy jako preferencja widoku | 3 | Menu hamburger (☰) panelu |

### 4.4. Pasek kart sesji

| Pole | Treść |
|---|---|
| Cel | Zarządzanie równoległymi przestrzeniami roboczymi jednego środowiska; karta sesji reprezentuje jedną samodzielną przestrzeń z własnym modułem, układem okien i kontekstem |
| Waga i miejsce | Pasek nad obszarem roboczym, na całej szerokości kolumn roboczych |
| Zawartość | Karty otwarte, grupy kart, karty przypięte, wskaźniki pracy w tle, wyzwalacz nowej karty |
| Warstwa | 1 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Karty otwarte | Przełączanie między przestrzeniami roboczymi | 1 | Kliknięcie karty |
| Wskaźnik pracy w tle karty | Kropka aktywności `●` sygnalizująca trwający przebieg pętli wykonawczej albo proces modułu | 1 | Widoczny bez interakcji |
| Karty przypięte | Węższa forma bez tytułu, na stałym początku paska — stały punkt odniesienia projektu | 1 | Widoczne, gdy istnieją |
| Grupy kart | Nazwane, barwione grupy odpowiadające projektowi albo wątkowi pracy; zwijanie grupy do jednej etykiety | 2 | Kliknięcie etykiety grupy |
| Podgląd pracy w tle | Dymek z opisem trwającego procesu, postępem i przejściem do Execution Loop Window albo do funkcji Mobile | 2 | Wskazanie kursorem wskaźnika karty |
| Zmiana kolejności kart | Przeciąganie karty w obrębie paska i między grupami | 2 | Przeciągnięcie |
| Menu karty | Zamknij inne, zamknij po prawej, przypnij, dodaj do grupy, zmień nazwę, duplikuj | 3 | Menu kontekstowe karty |
| Ręczna nazwa karty | Własna etykieta zamiast nazwy modułu, na przykład „Studio — rozdział 3” | 3 | Menu kontekstowe karty |
| Duplikowanie karty | Kopia bieżącej karty z tym samym modułem i układem, z odrębnym kontekstem | 3 | Menu kontekstowe karty |
| Lista i wyszukiwarka kart | Lista wszystkich kart z filtrem tekstowym i miniaturą, skok do karty po nazwie | 3 | Menu kebab (⋮) paska kart |
| Ostatnio zamknięte karty | Lista niedawno zamkniętych kart z przywróceniem stanu, dopóki proces sesji trwa w tle | 3 | Menu kebab (⋮) paska kart |
| Przelanie karty do innego środowiska | Przeniesienie karty modułu wspólnego do TalkIn albo CodeStudio, gdy moduł tam występuje | 4 | Paleta poleceń, menu kontekstowe karty |

```
 Makieta — pasek kart sesji i grupy kart (stan spoczynku)
 ══════════╤═══════════════════════════════════════════════════════════════
  Boczna   │ [◈ Alfa][Studio ●][Research][Design]  [+]                  ⋮
  nawigacja├───────────────────────────────────────────────────────────────
  modułów  │ Chat Window          │ Obszar roboczy modułu
           │ Użytkownik ↔         │
  Studio   │ Wykonawca            │ [WorkSpace] [Projekt Alfa] [Studio]
  Workspace│                      │
  Browser  │ ──────────────────── │
  Research │ Execution Loop ▸     │
        ☰  │ Koordynator ↔ Wyk.   │
 ══════════╧═══════════════════════════════════════════════════════════════
```

### 4.5. Pasek kontekstu środowiska

| Pole | Treść |
|---|---|
| Cel | Jawna prezentacja kontekstu bieżącej karty i szybkie przejście do selektorów: środowisko, projekt, moduł, model, wykonawca, polityka kontekstu, stan kanału sterującego, procesy w tle, zużycie kontekstu |
| Waga i miejsce | Wiersz lekkich znaczników w nagłówku obszaru roboczego, ponad zawartością modułu |
| Zawartość | Zestaw znaczników kontekstowych i wskaźników stanu, konfigurowalny co do zestawu i kolejności |
| Warstwa | 1 — wiersz znaczników; selektory otwierane ze znaczników należą do warstwy 2 |

| Znacznik | Treść i działanie | Warstwa | Sposób wywołania |
|---|---|---|---|
| Środowisko i projekt | `[WorkSpace] [Projekt Alfa]` — kliknięcie otwiera selektor projektu przypisanego karcie | 2 | Kliknięcie znacznika |
| Moduł bieżącej karty | Nazwa modułu; kliknięcie otwiera przełącznik modułu | 2 | Kliknięcie znacznika |
| Model i poziom wysiłku | `[Fable 5] [Ultra]` — zmiana modelu i intensywności pracy bez otwierania okna konfiguracji | 2 | Kliknięcie znacznika |
| Wykonawca | Model albo agent przypisany karcie | 2 | Kliknięcie znacznika |
| Polityka kontekstu | Zasięg kontekstu karty: odrębny albo współdzielony; kliknięcie otwiera podgląd „co widzi ta rozmowa” | 2 | Kliknięcie znacznika |
| Stan kanału i sesji | Stan kanału WebSocket i procesu sesji: połączony, wznawianie, praca w tle | 1 | Widoczny bez interakcji |
| Procesy w tle | Liczba aktywnych procesów środowiska; kliknięcie otwiera Execution Loop Window albo funkcję Mobile | 1 | Widoczny bez interakcji |
| Zużycie kontekstu | Bieżące zużycie kontekstu i tokenów karty, liczone `pkoukk/tiktoken-go` | 2 | Kliknięcie wskaźnika |
| Komunikaty środowiska | Zakotwiczenie komunikatów `.dn-toast` o zdarzeniach powłoki: zapisano układ, zakończono proces | 1 | Pojawiają się przy zdarzeniu |
| Szybkie przełączniki | Motyw jasny i ciemny, zwinięcie bocznej nawigacji, wywołanie Always On Display, wejście w funkcję Mobile | 3 | Menu kebab (⋮) paska kontekstu |

### 4.6. Obszar roboczy jako kontener i menedżer układów okien

| Pole | Treść |
|---|---|
| Cel | Kompozycja okien operacyjnych modułu otwartego w bieżącej karcie oraz utrwalanie tej kompozycji |
| Waga i miejsce | Prawa, dominująca kolumna obszaru roboczego; kolumny wewnętrzne rozmieszczone poziomo |
| Zawartość | Okna operacyjne modułu, uchwyty szerokości kolumn, nagłówek z paskiem kontekstu i menu operacji |
| Warstwa | 1 — kontener i okna modułu; 3 — menedżer układów |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kolumny okien modułu | Rozmieszczenie okien operacyjnych w sąsiadujących kolumnach; regulacji podlega wyłącznie szerokość | 1 | Przeciągnięcie uchwytu szerokości |
| Podział kontenera na kolumny | Ustawienie dwóch okien operacyjnych obok siebie zamiast na zakładkach | 3 | Menedżer układów okien, skrót |
| Dokowanie okien | Przeciągnięcie okna operacyjnego do krawędzi kontenera; strefy zrzutu z podpowiedzią wyrównania | 3 | Przeciągnięcie nagłówka okna |
| Tryb skupienia okna | Powiększenie jednego okna operacyjnego na cały kontener z powrotem do układu | 2 | Ikona okna, skrót |
| Siatka przyciągania | Przyciąganie krawędzi kolumn do wspólnej siatki i do siebie nawzajem | 2 | Działa przy zmianie szerokości |
| Zapisane układy okien | Zapis bieżącego rozmieszczenia pod nazwą i wczytanie go jednym poleceniem; klucze `shell.layout.save`, `shell.layout.apply` | 3 | Menedżer układów okien |
| Reset układu do domyślnego modułu | Przywrócenie układu zdefiniowanego przez moduł | 3 | Menedżer układów okien |
| Pamięć układu per karta i per moduł | Utrwalenie rozmieszczenia osobno dla każdej karty i domyślnie dla każdego modułu | 2 | Działa samoczynnie, sterowana z okna konfiguracji |
| Pamięć układu per klasa urządzenia | Odrębny układ, zwinięcia i motyw dla komputera, tabletu i telefonu | 4 | Okno konfiguracji |
| Tryb czystego ekranu | Ukrycie bocznej nawigacji, kart i pasków do samego obszaru roboczego | 4 | Paleta poleceń, skrót |

### 4.7. Wyszukiwarka środowiska i paleta poleceń

| Pole | Treść |
|---|---|
| Cel | Odnajdywanie wszystkiego, co istnieje w środowisku, oraz jeden punkt dostępu do wszystkich poleceń powłoki |
| Waga i miejsce | Wyszukiwarka — kolumna boczna otwierana jako rozszerzenie boczne; paleta poleceń — nakładka wyszukiwania nad obszarem roboczym |
| Zawartość | Wyniki wyszukiwania z filtrami zakresu; lista poleceń powłoki z dopasowaniem rozmytym |
| Warstwa | Wyszukiwarka — 2; paleta poleceń — 4 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Zakres wyszukiwania | Karty, moduły, projekty, materiały Library, zapisane sesje i układy; indeks po stronie Go, SQLite FTS5 | 2 | Pole wyszukiwania wstążki (z filtrem zakresu środowisk) |
| Filtry zakresu | Zawężenie wyników do wybranego typu obiektu i wybranego projektu | 2 | Znaczniki filtrów w panelu wyników |
| Paleta poleceń środowiska | Przełącz do modułu, otwórz kartę, wczytaj układ, wczytaj sesję, przełącz motyw, otwórz konfigurację; dopasowanie rozmyte `sahilm/fuzzy` | 4 | `Ctrl/Cmd + K` |
| Wyszukiwanie funkcji | Odnalezienie funkcji powłoki po nazwie, wraz z przypisanym jej skrótem | 4 | Paleta poleceń |
| Skróty do zapisanych sesji i grup kart | Przypisanie skrótu do zapisanej sesji albo grupy kart | 4 | Mapa skrótów |
| Mapa skrótów | Podgląd i przepięcie każdego skrótu powłoki, rozwiązywanie konfliktów | 4 | `Ctrl/Cmd + /`, okno konfiguracji |

### 4.8. Menedżer sesji i grup kart

| Pole | Treść |
|---|---|
| Cel | Utrwalenie i odtworzenie całego kontekstu pracy projektowej: zestawu kart, ich modułów, układów okien i grup |
| Waga i miejsce | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zawartość | Lista zapisanych przestrzeni projektu, lista grup kart, lista przypięć |
| Warstwa | 3 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Zapis przestrzeni projektu | Zapis całego zestawu otwartych kart pod nazwą; klucz `shell.session.save` | 3 | Przycisk panelu |
| Odtworzenie przestrzeni projektu | Wczytanie zapisanego zestawu kart jako kompletu; klucz `shell.session.restore` | 3 | Kliknięcie pozycji listy |
| Grupy kart | Tworzenie, nazywanie i barwienie grup; klucz `shell.tab.group` | 3 | Przycisk panelu |
| Szybki start projektu | Otwarcie kompletu kart według szablonu startu: karta modułu Workspace z projektem oraz przygotowane karty Research i Studio w jednej grupie | 3 | Przycisk panelu, paleta poleceń |
| Przypisanie skrótu | Powiązanie zapisanej sesji albo grupy kart ze skrótem klawiszowym | 4 | Mapa skrótów |

### 4.9. Magistrala kontekstu i przekazanie artefaktu

Mechanizm magistrali kontekstu — przeznaczenie, sposób przekazania artefaktu, zasięg, warstwę widoczności i sposób wywołania oraz relacje do modułu Library, punktów izolacji i encji artefaktu — opisuje [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) (rozdz. 6a). Specyfika środowiska WorkSpace: zasięg domyślny magistrali obejmuje grupę kart projektu, a moduł Workspace jest celem odbioru artefaktów kierowanych do projektu.

| Pole | Treść |
|---|---|
| Cel | Przenoszenie materiału i kontekstu między modułami i kartami bez ręcznego kopiowania oraz kierowanie wyników do trwałego miejsca |
| Waga i miejsce | Magistrala kontekstu — kolumna boczna otwierana jako rozszerzenie boczne; przekazanie artefaktu — panel popover przy elemencie wyniku |
| Zawartość | Fragmenty, pliki i odnośniki odłożone w sesji; lista celów przekazania |
| Warstwa | 3 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Odłożenie artefaktu do magistrali kontekstu | Przeciągnięcie fragmentu, pliku albo odnośnika na krawędź obszaru roboczego | 3 | Przeciągnięcie, menu kontekstowe wyniku |
| Pobranie artefaktu z magistrali kontekstu | Wstawienie odłożonego artefaktu do pola poleceń albo do okna operacyjnego modułu | 3 | Przeciągnięcie z panelu magistrali kontekstu |
| Przekazanie artefaktu między modułami | Wysłanie wyniku z jednego okna do innego modułu w nowej albo istniejącej karcie: Research → Studio, Studio → Apps, Design → Apps; klucz `shell.context.handoff` | 3 | Menu kontekstowe wyniku |
| Odbiór artefaktów do Library i projektu | Skierowanie wyniku dowolnego modułu do Library albo do projektu w module Workspace | 3 | Menu kontekstowe wyniku |
| Magistrala kontekstu o zasięgu projektu | Przypięte materiały, ustalenia i pamięć dostępne każdej karcie grupy projektu na poziomie zasięgu „projekt” | 2 | Znacznik polityki kontekstu, panel szybkiej konfiguracji sesji |
| Współdzielenie historii i pamięci między kartami | Ustanowienie współdzielenia historii, pamięci albo kontekstu między dwiema wskazanymi kartami | 2 | Panel szybkiej konfiguracji sesji |
| Podgląd polityki efektywnej | Prezentacja zasięgu obowiązującego bieżącą kartę oraz treści widocznej dla rozmowy przed wysłaniem polecenia | 2 | Kliknięcie znacznika polityki kontekstu |
| Odesłanie do okna punktów izolacji | Przejście do pełnej definicji profili i macierzy zasięgu — powłoka reguły stosuje i prezentuje, definiuje je okno punktów izolacji | 3 | Odnośnik panelu |

### 4.10. Centrum powiadomień i panel szybkiej konfiguracji sesji

Mechanizm centrum powiadomień — taksonomię klas zdarzeń, postać wizualną, umiejscowienie w układzie pionowym, warstwę widoczności, sposób wywołania, stany oraz działania dostępne z poziomu powiadomienia — opisuje karta komponentu w opracowaniu [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) (rozdz. 11.6). Poniżej podano wyłącznie to, co właściwe środowisku WorkSpace.

| Pole | Treść |
|---|---|
| Cel | Zbiorcza prezentacja zdarzeń środowiska oraz zmiana ustawień warstwy sesji bez opuszczania pracy |
| Waga i miejsce | Obie funkcje jako kolumny boczne, otwierane jako rozszerzenia boczne |
| Specyfika środowiska | Zakres źródeł zdarzeń obejmuje grupę kart projektu: przebiegi pętli wykonawczej kart grupy, automatyki wpięte do projektu, wzmianki oraz terminy zadań projektu prowadzonego w module Workspace; filtr źródła domyślnie zawęża listę do projektu bieżącej grupy kart |
| Warstwa | Centrum powiadomień — 3; panel szybkiej konfiguracji sesji — 2 |

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Filtr projektu | Zawężenie listy zdarzeń do projektu bieżącej grupy kart | 3 | Znacznik projektu w nagłówku kolumny |
| Ustawienia warstwy sesji | Model, profil asystenta, motyw, retencja historii, polityka współdzielenia, przypisany projekt | 2 | Kliknięcie znacznika kontekstu |
| Odnośnik do okna konfiguracji | Przejście do pełnego zakresu ustawień środowiska | 2 | Odnośnik panelu |
| Powrót na stronę główną z zachowaniem tła | Wyjście ze środowiska; karty trwają w tle i odtwarzają pełny stan przy powrocie | 1 | Kontrolka „Centrum dowodzenia” na wstążce |
| Punkt dostępu Always On Display | Przywołanie pływającego agenta towarzyszącego ponad bieżącą przestrzenią, w każdej karcie i module | 2 | Szybki wybór szyny nawigacji, skrót |
| Punkt dostępu Mobile | Wejście w funkcję globalną Mobile — nadzór nad procesami środowiska działającymi w tle; nie tworzy nowej przestrzeni roboczej | 2 | Szybki wybór szyny nawigacji |

---

## 5. Przepływy pracy

### 5.1. Rozpoczęcie projektu w środowisku

1. Użytkownik wybiera środowisko WorkSpace na stronie głównej. Powłoka otwiera kartę z modułem Workspace i obydwoma kanałami komunikacji operacyjnej.
2. Użytkownik wydaje w Chat Window polecenie założenia projektu w języku naturalnym. Wykonawca zakłada projekt w module Workspace, powłoka wiąże go ze znacznikiem projektu w pasku kontekstu.
3. Szybki start projektu otwiera komplet kart według szablonu: karta modułu Workspace oraz karty Research i Studio w jednej, nazwanej grupie.
4. Użytkownik zapisuje układ okien pierwszej karty pod nazwą i zapisuje całą przestrzeń projektu w menedżerze sesji.

### 5.2. Zlecenie realizowane w pętli wykonawczej ponad modułami

1. Użytkownik formułuje zlecenie w Chat Window i przekazuje je do Koordynatora.
2. Execution Loop Window przyjmuje zlecenie, przedstawia jego dekompozycję na zadania i przypisuje każdemu zadaniu wykonawcę oraz moduł realizacji.
3. Koordynator prowadzi pętlę: przydziela zadania, odbiera raporty postępu, kieruje wyniki do kontroli jakości i podejmuje decyzje o ponowieniu.
4. Powłoka sygnalizuje przebieg wskaźnikiem pracy w tle na karcie sesji i przy pozycji modułu w bocznej nawigacji.
5. Użytkownik nadzoruje przebieg w Execution Loop Window: wstrzymuje, wznawia, przerywa albo koryguje zlecenie bez utraty dotychczasowych wyników.
6. Wyniki zadań trafiają do Library albo do projektu w module Workspace; dziennik przebiegu zapisuje się jako dokument projektu.

### 5.3. Przeniesienie materiału między modułami

1. W karcie modułu Research użytkownik zaznacza ustalenia i odkłada je na magistrali kontekstu.
2. Przekazanie artefaktu kieruje ustalenia do modułu Studio — powłoka otwiera nową kartę Studio w tej samej grupie projektu i wstawia materiał jako szkic.
3. Znacznik polityki kontekstu prezentuje zasięg obowiązujący nową kartę; podgląd polityki efektywnej pokazuje treść widoczną dla rozmowy.
4. Gotowa treść przekazywana jest z modułu Studio do modułu Apps, a wynik końcowy odbierany do Library.

### 5.4. Praca równoległa nad wieloma kartami jednego dnia

1. Użytkownik utrzymuje kilkanaście kart zorganizowanych w grupy odpowiadające projektom.
2. Karta modułu Workspace z projektem nadrzędnym pozostaje przypięta na początku paska.
3. Przełącznik kart z filtrem tekstowym prowadzi do karty po nazwie; skróty klawiszowe przenoszą między kartami i modułami.
4. Wskaźniki pracy w tle wskazują karty, w których trwa przebieg pętli wykonawczej albo proces modułu.
5. Zamknięcie karty przez pomyłkę cofa lista ostatnio zamkniętych kart, dopóki proces sesji trwa w tle.

### 5.5. Wznowienie pracy po przerwie

1. Użytkownik wraca na stronę główną i wybiera środowisko WorkSpace. Karty trwają w tle i odtwarzają pełny stan.
2. Menedżer sesji wczytuje zapisaną przestrzeń projektu jako komplet kart z ich modułami, układami i grupami.
3. Powłoka stosuje układ zapamiętany dla bieżącej klasy urządzenia.
4. Execution Loop Window prezentuje bieżący stan pętli wykonawczych, które trwały pod nieobecność użytkownika.

### 5.6. Nadzór nad procesem spoza stanowiska roboczego

1. Przebieg pętli wykonawczej trwa w tle po rozłączeniu klienta.
2. Funkcja Mobile udostępnia stan pętli, kolejkę zadań i wskaźniki przebiegu.
3. Użytkownik zatwierdza krok, wstrzymuje albo wznawia przebieg z urządzenia zdalnego.
4. Always On Display sygnalizuje zdarzenia wymagające decyzji i prowadzi do właściwej karty po powrocie do stanowiska.

---

## 6. Stany i powiązania

### 6.1. Stany powłoki środowiska

| Stan | Wyzwolenie | Prezentacja |
|---|---|---|
| Gotowość | Kanał sterujący połączony, karta aktywna | Pełny zestaw znaczników paska kontekstu, wskaźnik kanału w stanie połączonym |
| Praca w tle | Trwa przebieg pętli wykonawczej albo proces modułu | Kropka aktywności na karcie i przy pozycji modułu, licznik procesów w pasku kontekstu |
| Wznawianie | Klient odzyskuje kanał po rozłączeniu | Wskaźnik kanału w stanie wznawiania, okna monitorujące odtwarzają bieżący stan |
| Odtwarzanie sesji | Wczytanie zapisanej przestrzeni projektu albo powrót ze strony głównej | Karty odtwarzane kolejno z układami okien i grupami |
| Skupienie | Tryb skupienia okna albo tryb czystego ekranu | Ukrycie elementów powłoki poza obszarem roboczym, wyjście jednym skrótem |
| Zdarzenie środowiska | Zakończenie procesu, zapis układu, wpięta automatyka | Komunikat `.dn-toast` zakotwiczony przy pasku kontekstu, pozycja w centrum powiadomień |

### 6.2. Stany karty sesji

| Stan | Znaczenie |
|---|---|
| Aktywna | Karta widoczna w obszarze roboczym, kanały komunikacji operacyjnej podłączone do jej kontekstu |
| W tle | Karta niewidoczna, proces sesji trwa, przebiegi pętli postępują |
| Przypięta | Karta na stałym początku paska, w węższej formie bez tytułu |
| W grupie | Karta należy do nazwanej grupy projektu i korzysta z szyny kontekstu tej grupy, gdy zasięg jest ustawiony na projekt |
| Zapisana | Karta należy do zapisanej przestrzeni projektu i odtwarza się razem z nią |
| Zamknięta niedawno | Karta dostępna do przywrócenia, dopóki proces sesji trwa w tle |

### 6.3. Powiązania z modułami i funkcjami globalnymi

| Powiązanie | Zakres | Sterowanie |
|---|---|---|
| Moduł Workspace | Projekt nadrzędny grupy kart, cel odbioru artefaktów, źródło terminów i zadań prezentowanych w centrum powiadomień. Zakres funkcjonalny modułu — [Moduł WorkSpace](../moduly/workspace.md) | `shell.context.sink`, znacznik projektu |
| Moduł Library | Trwałe repozytorium środowiska, cel odbioru artefaktów i zakres wyszukiwania środowiska | `shell.context.sink`, `shell.search.scope` |
| Moduł Agents | Źródło wykonawców przypisywanych zadaniom pętli wykonawczej | Selektor wykonawcy w pasku kontekstu |
| Moduł Automations | Harmonogram przebiegów pętli wykonawczej i praca cykliczna | Panel harmonogramu Execution Loop Window |
| Okno punktów izolacji | Definicja profili i macierzy zasięgu kontekstu; powłoka reguły stosuje i prezentuje | `shell.context.scope` |
| Always On Display | Pływający agent towarzyszący ponad bieżącą przestrzenią, obecny w każdej karcie i module | Szybki wybór szyny nawigacji |
| Mobile | Nadzór i interwencja zdalna nad procesami środowiska działającymi w tle | Szybki wybór szyny nawigacji |
| Środowiska TalkIn i CodeStudio | Przelanie karty modułu wspólnego między środowiskami | Paleta poleceń |

---

## 7. Punkty sterowania z okna konfiguracji

Wszystkie funkcje powłoki środowiska są sterowalne z okna konfiguracji, w zakresie „Aplikacja i procesy” oraz w warstwie sesji. Każdy element ma objaśnienie kontekstowe `[?]`. Klucze są jawne, brak ustawienia oznacza wartość domyślną.

### 7.1. Mapa punktów sterowania

| Obszar | Punkt sterowania | Klucz | Wartość domyślna |
|---|---|---|---|
| Nawigacja | Grupy, kolejność i widoczność modułów | `shell.nav.groups`, `shell.nav.order`, `shell.nav.hidden` | Dziewięć modułów, kolejność źródłowa, wszystkie widoczne |
| Nawigacja | Przypięte moduły | `shell.nav.pinned` | Brak przypięć |
| Nawigacja | Zwinięcie panelu do ikon | `shell.nav.rail` | Panel rozwinięty |
| Nawigacja | Zachowanie kliknięcia modułu: przeładowanie karty albo nowa karta | `shell.nav.open_mode` | Przeładowanie bieżącej karty |
| Układ okien | Zapisane układy i układ domyślny per moduł | `shell.layout.presets`, `shell.layout.default` | Układ zdefiniowany przez moduł |
| Układ okien | Siatka przyciągania | `shell.layout.snap` | Włączona |
| Układ okien | Zapamiętane położenie paska kolumny Chat Window — ustawiane przesunięciem paska, nie suwakiem okna konfiguracji | `shell.chat.width` | Kolumna lewa, domyślne położenie startowe |
| Układ okien | Stan kolumny Execution Loop Window | `shell.loop.column` | Zwinięta do nagłówka pętli |
| Układ okien | Pamięć układu per klasa urządzenia | `shell.layout.per_device` | Włączona |
| Karty sesji | Zasięg kontekstu nowej karty | `shell.tab.isolation` | Odrębna historia, pamięć i kontekst |
| Karty sesji | Zapisane sesje i grupy kart | `shell.session.saved`, `shell.tab.groups` | Brak zapisanych |
| Karty sesji | Przechowywanie ostatnio zamkniętych kart | `shell.tab.recent_keep` | Pamiętane w czasie trwania sesji połączenia |
| Pasek kontekstu | Zestaw i kolejność znaczników | `shell.context_bar.items` | Środowisko, projekt, moduł, model, wykonawca, polityka kontekstu, procesy, zużycie |
| Skróty | Mapa skrótów powłoki | `shell.keymap` | Zestaw domyślny, w pełni przepinany |
| Współdzielenie | Zasięg szyny kontekstu projektu | `shell.context.scope` | Odrębny, bez współdzielenia |
| Współdzielenie | Cele odbioru artefaktów | `shell.context.sink` | Library oraz projekt w module Workspace |
| Wyszukiwanie | Zakres wyszukiwania środowiska | `shell.search.scope` | Karty, moduły, projekty, Library |
| Funkcje wspólne | Motyw środowiska | `shell.theme` | Zgodny z ustawieniem systemowym |
| Funkcje wspólne | Powiadomienia środowiska | `shell.notify` | Włączone, bez dźwięku |
| Funkcje wspólne | Szablon szybkiego startu projektu | `shell.project.template` | Workspace, Research, Studio |

### 7.2. Szablon konfiguracji środowiska WorkSpace

```
# Środowisko WorkSpace — konfiguracja powłoki
# Brak ustawienia oznacza wartość domyślną, nigdy blokadę
srodowisko_workspace:
  nawigacja:
    grupy:            domyslne        # Materiał · Tworzenie · Organizacja
    kolejnosc:        zrodlowa        # Studio·Workspace·Browser·Research·Library·
                                      # Roundtable·Design·Apps·Agents
    przypiete:        []
    rail:             wylaczony
    tryb_otwarcia:    przeladuj       # wartość alternatywna: nowa_karta
  uklad_okien:
    presety:          []
    domyslny:         modul
    siatka_snap:      wlaczona
    chat_window:      kolumna_lewa
    execution_loop:   kolumna_zwinieta
    uklad_per_urzadzenie: wlaczony
  karty_sesji:
    izolacja_karty:   odrebna         # historia, pamięć i kontekst odrębne
    grupy_kart:       []
    zapisane_sesje:   []
  pasek_kontekstu:
    znaczniki:        [srodowisko, projekt, modul, model, wykonawca,
                       polityka_kontekstu, procesy, zuzycie]
  skroty:
    mapa:             domyslna        # w pełni przepinalna
  wspoldzielenie_kontekstu:
    zasieg:           odrebny         # poziom „projekt” ustawiany świadomie
    odbior_artefaktow: [library, workspace]
    odsylacz_izolacja: okno_punktow_izolacji
  funkcje_wspolne:
    wyszukiwanie_zakres: [karty, moduly, projekty, library]
    motyw:            systemowy
    powiadomienia:    wlaczone
    szablon_startu:   [workspace, research, studio]
  zrodlo_konfiguracji: "okno konfiguracji"
  przywrocenie_domyslnych: dostepne
```

Współdzielenie kontekstu jest świadomą decyzją użytkownika na jawnie wskazanym poziomie zasięgu; powłoka nie łączy kontekstów samoczynnie. Każdy stan domyślny jest wyjściowy i przywracalny; żadne ustawienie nie zamyka dostępu do modułu ani nie wymusza układu okien.

---

## 8. Scenariusze użycia

### 8.1. Kancelaria prowadząca kilka spraw równolegle

Każda sprawa ma własną grupę kart: karta modułu Workspace ze sprawą, karta Research z materiałem źródłowym, karta Studio z redakcją pisma. Karta sprawy wiodącej pozostaje przypięta. Zlecenie „przygotuj projekt pisma na podstawie ustaleń” przechodzi z Chat Window do Koordynatora; Execution Loop Window prowadzi zadania: zebranie materiału w module Research, redakcja w module Studio, kontrola jakości wobec kryterium ukończenia. Gotowe pismo trafia do projektu w module Workspace, dziennik przebiegu — do akt sprawy.

### 8.2. Zespół produktowy prowadzący projekt od materiału do produktu

Grupa kart projektu obejmuje moduły Research, Design, Studio, Apps i Workspace. Magistrala kontekstu o zasięgu projektu udostępnia każdej karcie przypięte materiały i ustalenia projektu. Makieta z modułu Design przekazywana jest do modułu Apps przekazaniem artefaktu, treść z modułu Studio — tą samą drogą. Execution Loop Window prowadzi budowę produktu jako jedną pętlę obejmującą zadania w trzech modułach, a wskaźniki przebiegu widoczne są na kartach i przy pozycjach modułów.

### 8.3. Praca ciągła nad projektem cyklicznym

Przebieg pętli wykonawczej wpięty jest w automatykę modułu Automations i realizuje cykliczne zadanie projektowe. Użytkownik nadzoruje przebieg przez funkcję Mobile: odbiera zdarzenia z centrum powiadomień, zatwierdza kroki oczekujące na decyzję i wstrzymuje przebieg w razie zgłoszenia przeszkody. Po powrocie do stanowiska zapisana przestrzeń projektu odtwarza pełny układ pracy, a Execution Loop Window prezentuje historię przebiegów.

### 8.4. Praca głęboko skupiona nad jednym materiałem

Użytkownik wczytuje zapisany układ „Redakcja”: obszar roboczy podzielony na kolumnę materiału źródłowego i kolumnę redakcyjną, Chat Window w kolumnie lewej, Execution Loop Window zwinięta do nagłówka pętli. Tryb czystego ekranu ukrywa nawigację i karty. Wyjście z trybu oraz powrót do pełnego środowiska następuje jednym skrótem, bez zmiany kontekstu karty.

### 8.5. Prowadzenie jednego projektu na wielu urządzeniach

Ten sam projekt otwierany jest na komputerze i na tablecie. Pamięć układu per klasa urządzenia utrzymuje odrębne rozmieszczenie kolumn dla każdego ekranu przy zachowaniu tej samej przestrzeni projektu, tych samych kart i tego samego kontekstu. Procesy trwają w tle niezależnie od urządzenia, z którego prowadzony jest nadzór.

---

## 9. Wykazy normatywne i kryteria odbioru

### 9.1. Komendy kontraktu

Środowisko WorkSpace korzysta z komend dwóch obszarów kontraktu: `environment` (2 komendy — wejście do środowiska, rozdz. 2.1) oraz `workspace` (59 komend — moduł Workspace, [Moduł WorkSpace](../moduly/workspace.md)). Pełny wykaz nazw komend, pól żądania i pól wyniku dla wszystkich 59 komend obszaru `workspace` znajduje się w [Załączniku B. Pełnym wykazie komend kontraktu środowiska WorkSpace](#załącznik-b-pełny-wykaz-komend-kontraktu-środowiska-workspace), na końcu niniejszego dokumentu. Powłoka środowiska — przedmiot niniejszego dokumentu — nie ma własnego obszaru kontraktu; jej ustawienia są kluczami konfiguracji `shell.*` (rozdz. 7), odrębnymi kategorialnie od komend kontraktu.

| Obszar | Liczba komend | Reprezentatywne komendy | Zastosowanie w środowisku |
|---|---|---|---|
| `environment` | 2 | `environment.list`, `environment.enter` | Wejście do środowiska z przedsionka (rozdz. 2.1) |
| `workspace` | 59 | `workspace.project.list`, `workspace.task.create`, `workspace.board.get`, `workspace.context.get` | Moduł Workspace — Project Dashboard, Task Board, Context Memory, Project Library ([Moduł WorkSpace](../moduly/workspace.md)) |

### 9.2. Żetony projektowe

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Kropka sygnału procesu aktywnego | `--dn-kropka`, `--dn-czas-tetno` | Wskaźnik pracy w tle przy karcie sesji i przy pozycji modułu z procesem w toku |
| Fokus kontrolek | `--dn-fokus` | Obrys fokusu klawiaturowego pól, przycisków, pozycji list |
| Przycisk główny akcji | `--dn-sygnal`, `--dn-sygnal-wypelnienie`, `--dn-sygnal-wypelnienie-hover` | „Nowy projekt”, akcje pierwszoplanowe przedsionka i paska kontekstu |
| Stan powodzenia zapisu układu | `--dn-sukces-tekst`, `--dn-sukces-tlo`, `--dn-sukces-obrys` | Potwierdzenie `.dn-toast` po `shell.layout.save` |
| Stan ostrzeżenia terminu zadania | `--dn-ostrzezenie-tekst`, `--dn-ostrzezenie-tlo`, `--dn-ostrzezenie-obrys` | Plakietka terminu zbliżającego się w pasku kontekstu środowiska |
| Odstępy paneli i list | `--dn-od-2`, `--dn-od-4`, `--dn-od-6` | Odstęp wewnętrzny kart projektów, wierszy szyny sesji, pól panelu konfiguracji |
| Promień zaokrąglenia | `--dn-r-sm`, `--dn-r-md`, `--dn-r-xl` | Pola i przyciski (`--dn-r-sm`), karty projektów (`--dn-r-md`), karta środowiska w przedsionku (`--dn-r-xl`) |
| Typografia interfejsu | `--dn-ff-bazowa`, `--dn-fs-base`, `--dn-lh-bazowy` | Tekst kontrolek, etykiet i list w gęstości zwartej |
| Czas przejścia warstw | `--dn-czas-3` | Otwarcie/zamknięcie panelu bocznego, palety poleceń, menedżera układów |

### 9.3. Komponenty interfejsu

| Klasa `.dn-*` | Rola | Modyfikatory użyte w środowisku | Stany |
|---|---|---|---|
| `.dn-btn` | Przycisk | `--sygnal` (akcja główna), `--zarys` (akcja pomocnicza), `--duch` (akcja trzeciorzędna) | domyślny · wskazanie kursorem · wciśnięty · ładowanie |
| `.dn-pole` | Wrapper pola formularza | — | domyślny · błąd |
| `.dn-pole-kontrolka` | Kontrolka pola (tekst, wybór, wieloliniowa) | — | puste · wypełnione · fokus · błąd |
| `.dn-karta` | Karta | `--klikalna`, `--wybrana` | domyślna · wskazanie kursorem · wybrana |
| `.dn-plakietka` | Plakietka statusu | `--sukces`, `--ostrzezenie`, `--blad`, `--informacja` | domyślna |
| `.dn-kropka` | Wskaźnik semantyczny stanu | — | aktywny (tętno) · nieaktywny |
| `.dn-zakladki` | Pas kart sesji | — | aktywna · w tle · z pracą w toku |
| `.dn-toast` | Dymek powiadomienia | — | widoczny · wygaszany |
| `.dn-pusty-stan` | Stan pusty listy albo panelu | — | domyślny |
| `.dn-spinner` | Wskaźnik ładowania | — | aktywny |

### 9.4. Etykiety interfejsu z prototypu

Wykaz dosłownych brzmień z `design/05-okna/srodowiska/workspace.html` i `workspace-przedsionek.html`, nieopisanych już w rozdziałach 1–8.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Nagłówek środowiska w przedsionku | „WorkSpace”, „Projekty, procesy i produkty” | Rozdz. 2.1 |
| Siatka modułów przedsionka | „Projekty, procesy i produkty — dziewięć modułów.” | Rozdz. 2.1 |
| Listwa działań przedsionka | „Nowy projekt”, „Otwórz projekt istniejący”, „Utwórz z szablonu projektu” | Rozdz. 2.1 |
| Wybór środowiska dla nowego projektu | „Nowy projekt — w którym środowisku”, „Założenie projektu — wskazanie środowiska otwiera od razu jego Workspace.” | Rozdz. 2.1 |
| Szyna sesji przedsionka | „Twoje sesje”, licznik „5 projektów · 7 sesji” | Rozdz. 2.1 |
| Opis kafla Workspace | „Projekty, zadania, pamięć kontekstowa i biblioteka projektu.” | Rozdz. 2 |
| Opis kafla Studio | „Redakcja i tworzenie treści projektu” | Rozdz. 2 |
| Opis kafla Design | „Makiety, projekt wizualny i materiały graficzne” | Rozdz. 2 |
| Opis kafla Browser | „Praca z materiałem sieciowym w obrębie projektu” | Rozdz. 2 |
| Opis kafla Agents | „Delegacja zadań projektu do agentów użytkownika” | Rozdz. 2 |

### 9.5. Komunikaty

| Sytuacja | Dosłowna treść albo wzorzec | Rodzaj |
|---|---|---|
| Układ okien zapisany | Dymek powiadomienia potwierdzający zapis pod nazwą układu (`shell.layout.save`) | Potwierdzenie |
| Przestrzeń projektu odtworzona | Dymek powiadomienia potwierdzający wczytanie kompletu kart (`shell.session.restore`) | Potwierdzenie |
| Artefakt przekazany między modułami | Dymek powiadomienia ze wskazaniem modułu docelowego (`shell.context.handoff`) | Potwierdzenie |
| Zdarzenie powłoki (zakończony proces, wpięta automatyka) | Dymek powiadomienia zakotwiczony przy pasku kontekstu, pozycja w centrum powiadomień | Powiadomienie |
| Odrzucenie polecenia w trybie obserwatora | Komunikat kontekstowy ze wskazaniem przełączenia w tryb operatora — kontrolka pozostaje klikalna | Komunikat kontekstowy |
| Brak zapisanych układów albo przestrzeni projektu | Stan pusty panelu z podpowiedzią zapisu bieżącego rozmieszczenia | Stan pusty |
| Przekroczenie limitu WIP na tablicy zadań modułu | Ostrzeżenie niesione polem `wipExceeded` wyniku `workspace.task.move` — zapis odbywa się mimo ostrzeżenia | Ostrzeżenie |

### 9.6. Punkty łamania

Wartości progów z `design/zasoby/zetony/zetony.css`, rozdział „Siatka i punkty łamania”.

| Żeton | Próg szerokości | Zachowanie układu w środowisku WorkSpace |
|---|---|---|
| `--dn-bp-w1` | 640px | Widok mobilny — środowisko udostępniane przez funkcję globalną Mobile ([Mobile](../funkcje-globalne/mobile.md)) |
| `--dn-bp-w2` | 960px | Boczna nawigacja modułów zwija się do samych ikon; szyna sesji przedsionka zwija się analogicznie |
| `--dn-bp-w3` | 1280px | Pełny układ — Chat Window, obszar roboczy modułu i kolumna stanu środowiska widoczne jednocześnie |
| `--dn-bp-w4` | 1600px | Szerokie biurko — Chat Window i Execution Loop Window rozwijają się do pełnej szerokości roboczej bez zwijania obszaru roboczego modułu |

Poniżej `--dn-bp-w2` menedżer sesji i grup kart (rozdz. 4.8) oraz centrum powiadomień (rozdz. 4.10) przechodzą z prezentacji jako rozszerzenie boczne na prezentację w nakładce.

### 9.7. Stany kontrolek

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| Przycisk „Nowy projekt” (`.dn-btn--sygnal`) | Wypełnienie `--dn-sygnal-wypelnienie` | Wypełnienie `--dn-sygnal-wypelnienie-hover` | Przyciemnienie o jeden krok | Obrys `--dn-fokus` | nie występuje (zero blokad) | Spinner (`.dn-spinner`) w miejscu ikony | — | — |
| Karta projektu (`.dn-karta--klikalna`) | Tło `--dn-tlo` | Tło `--dn-powierzchnia`, uniesienie cienia | Przyciśnięcie | Obrys `--dn-fokus` | nie występuje | Kropka `--dn-kropka` z tętnem | Karta „Załóż pierwszy projekt” | Plakietka `--dn-blad-tlo` przy projekcie z niepowodzeniem synchronizacji |
| Pole filtrowania listy | Obrys `--dn-obrys` | Obrys `--dn-sygnal-obrys` | — | Obrys `--dn-fokus` | nie występuje | — | Tekst podpowiedzi „Szukaj projektu…” | — |
| Pozycja bocznej nawigacji modułu | Tekst `--dn-atrament` | Tło `--dn-sygnal-tlo` | Przyciśnięcie | Obrys `--dn-fokus` | nie występuje | — | — | — |
| Karta zadania na tablicy (moduł Workspace) | Tło `--dn-tlo` | Uniesienie cienia | Przeciąganie (`workspace.task.move`) | Obrys `--dn-fokus` | nie występuje | — | Kolumna bez zadań — `.dn-pusty-stan` | Plakietka terminu przekroczonego `--dn-blad-tekst` |

**Nota terminologiczna.** Żeton `--dn-sygnal-wypelnienie-hover` jest identyfikatorem
technicznym zapisanym w `design/zasoby/zetony/zetony.css` i pozostaje w tej postaci
niezmieniony. Człon `hover` w jego nazwie oznacza stan opisany w niniejszym dokumencie
i w całym zbiorze jako **wskazanie kursorem** — kolumna „Wskazanie kursorem” tabeli
powyżej oraz kolumna żetonów w rozdziale 9.2 mówią o tym samym stanie kontrolki.

### 9.8. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Wejście do środowiska otwiera przedsionek, nie przestrzeń roboczą | Kliknięcie karty WorkSpace na stronie głównej; obserwacja, że żaden moduł nie jest otwarty automatycznie |
| Wszystkie dziewięć modułów jest osiągalnych z bocznej nawigacji | Przegląd bocznej nawigacji względem wykazu rozdz. 2 |
| Chat Window jest obecne w każdym module środowiska w tym samym miejscu układu | Otwarcie kolejno wszystkich dziewięciu modułów; potwierdzenie stałej pozycji lewej kolumny |
| Wszystkie 59 komend obszaru `workspace` mają pokrycie w oknach modułu Workspace | Zestawienie wykazu Załącznika B z oknami Project Dashboard, Task Board, Context Memory, Project Library |
| Żaden przycisk nie jest trwale nieaktywny | Przegląd kontrolek wykazanych w rozdz. 9.3 i 9.7 pod kątem stanu `disabled`; zero wystąpień poza stanem „ładowanie” |
| Zapisany układ okien odtwarza się identycznie po ponownym uruchomieniu | Zapis układu, zamknięcie i ponowne otwarcie karty; porównanie rozmieszczenia kolumn |
| Pamięć układu per klasa urządzenia utrzymuje odrębne rozmieszczenia | Otwarcie tego samego projektu na dwóch szerokościach referencyjnych z tabeli rozdz. 9.6; potwierdzenie niezależnego zapisu |
| Układ odpowiada progom łamania rozdz. 9.6 na czterech szerokościach referencyjnych | Zmiana szerokości okna kolejno przez 640px, 960px, 1280px, 1600px; obserwacja zachowań z tabeli rozdz. 9.6 |

---

## Załącznik A. Skróty klawiszowe środowiska

Wszystkie skróty są przepinalne z mapy skrótów i z okna konfiguracji; konflikty rozwiązuje mapa skrótów.

| Obszar | Skrót | Działanie |
|---|---|---|
| Powłoka | `Ctrl/Cmd + K` | Paleta poleceń środowiska |
| Powłoka | `Ctrl/Cmd + /` | Mapa skrótów |
| Powłoka | `Ctrl/Cmd + F` | Wyszukiwarka środowiska |
| Powłoka | `Ctrl/Cmd + ,` | Okno konfiguracji |
| Moduły | `Ctrl/Cmd + 1` … `Ctrl/Cmd + 9` | Przejście do modułu w kolejności bocznej nawigacji |
| Moduły | `Ctrl/Cmd + B` | Zwinięcie bocznej nawigacji do ikon |
| Karty | `Ctrl/Cmd + T` | Nowa karta sesji |
| Karty | `Ctrl/Cmd + W` | Zamknięcie karty |
| Karty | `Ctrl/Cmd + Shift + T` | Przywrócenie ostatnio zamkniętej karty |
| Karty | `Ctrl/Cmd + Tab` | Następna karta |
| Karty | `Ctrl/Cmd + Shift + Tab` | Poprzednia karta |
| Karty | `Ctrl/Cmd + Shift + E` | Przełącznik i wyszukiwarka kart |
| Komunikacja | `Ctrl/Cmd + Enter` | Wysłanie polecenia z Chat Window |
| Komunikacja | `Esc` | Przerwanie działania Wykonawcy |
| Pętla wykonawcza | `Ctrl/Cmd + L` | Rozwinięcie i zwinięcie kolumny Execution Loop Window |
| Pętla wykonawcza | `Ctrl/Cmd + .` | Wstrzymanie i wznowienie przebiegu pętli |
| Układ okien | `Ctrl/Cmd + \` | Podział kontenera na kolumny |
| Układ okien | `Ctrl/Cmd + Shift + M` | Tryb skupienia okna |
| Układ okien | `Ctrl/Cmd + Shift + R` | Reset układu do domyślnego modułu |
| Układ okien | `Ctrl/Cmd + Shift + S` | Zapis bieżącego układu okien |
| Sesje | `Ctrl/Cmd + Shift + O` | Menedżer sesji i grup kart |
| Środowisko | `Ctrl/Cmd + Shift + F` | Tryb czystego ekranu |
| Środowisko | `Ctrl/Cmd + Shift + D` | Przełącznik motywu jasny i ciemny |
| Środowisko | `Ctrl/Cmd + Shift + A` | Always On Display |
| Środowisko | `Ctrl/Cmd + Shift + N` | Centrum powiadomień |
| Dostępność | `Tab` | Kolejność fokusu: szyna nawigacji → belka tytułowa → wstążka (pasmo kart sesji) → boczna nawigacja modułów → Chat Window → Execution Loop Window → obszar roboczy → kolumny boczne |


---

## Załącznik B. Pełny wykaz komend kontraktu środowiska WorkSpace

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `workspace` — 59 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `workspace.enter` | Przeładowuje przestrzeń roboczą karty sesji na wskazany moduł | `sessionId:string` (wym)<br>`moduleId:string` (wym)<br>`windowId:string` (opc) | `session:Session` (wym)<br>`module:Module` (wym)<br>`window:Window` (wym)<br>`operationalWindowCodes:string[]` (wym) |
| `workspace.dashboard.get` | Zwraca zestawienie stanu projektu dla Project Dashboard | `projectId:string` (wym) | `dashboard:WorkspaceDashboard` (wym) |
| `workspace.instructions.set` | Zapisuje instrukcję systemową projektu | `projectId:string` (wym)<br>`content:string` (wym)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `instructions:WorkspaceInstructions` (wym) |
| `workspace.context.set` | Zapisuje albo zmienia wpis pamięci projektu | `projectId:string` (wym)<br>`content:string` (wym)<br>`entryId:string` (opc)<br>`pinned:bool` (opc)<br>`origin:MemoryEntryOrigin` (opc)<br>`scope:ConfigScope` (opc)<br>`tags:string[]` (opc) | `entry:WorkspaceMemoryEntry` (wym) |
| `workspace.context.get` | Zwraca wpisy pamięci projektu wraz z ich zasięgiem | `projectId:string` (wym)<br>`includeShared:bool` (opc)<br>`limit:int` (opc)<br>`tag:string` (opc)<br>`query:string` (opc) | `entries:WorkspaceMemoryEntry[]` (wym) |
| `workspace.library.list` | Zwraca pliki biblioteki projektu | `projectId:string` (wym)<br>`query:string` (opc)<br>`limit:int` (opc)<br>`tag:string` (opc) | `files:LibraryFile[]` (wym)<br>`total:int` (opc) |
| `workspace.agent.assign` | Przypisuje eksperta jako wykonawcę zadań w projekcie | `projectId:string` (wym)<br>`agentId:string` (wym)<br>`role:string` (opc)<br>`defaultExecutor:bool` (opc)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `assignment:WorkspaceAgentAssignment` (wym) |
| `workspace.agent.unassign` | Odłącza eksperta od projektu. Nie usuwa eksperta z biblioteki modułu Agents — znosi wyłącznie jego przypisanie do tego projektu | `projectId:string` (wym)<br>`agentId:string` (wym) | `unassigned:bool` (wym) |
| `workspace.project.status.set` | Ustawia stan projektu. Droga do stanu `paused`, którego opracowanie wymaga, a którego żaden dzisiejszy uchwyt nie potrafi zapisać | `projectId:string` (wym)<br>`status:WorkspaceProjectStatus` (wym) | `project:WorkspaceProject` (wym) |
| `workspace.instructions.version.list` | Zwraca wykaz wersji instrukcji systemowych projektu dla panelu „Wersje” | `projectId:string` (wym)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `versions:WorkspaceInstructionsVersion[]` (wym)<br>`total:int` (wym) |
| `workspace.instructions.version.restore` | Przywraca wskazaną wersję instrukcji jako obowiązującą. Przywrócenie zakłada wersję nową o treści wersji wskazanej — historia nie jest przepisywana | `projectId:string` (wym)<br>`versionId:string` (wym) | `instructions:WorkspaceInstructions` (wym)<br>`version:WorkspaceInstructionsVersion` (wym) |
| `workspace.task.create` | Zakłada zadanie projektu | `projectId:string` (wym)<br>`title:string` (wym)<br>`description:string` (opc)<br>`status:WorkspaceTaskStatus` (opc)<br>`priority:WorkspaceTaskPriority` (opc)<br>`assigneeKind:WorkspaceAssigneeKind` (opc)<br>`assigneeId:string` (opc)<br>`parentTaskId:string` (opc)<br>`boardColumnId:string` (opc)<br>`startAt:int64` (opc)<br>`dueAt:int64` (opc)<br>`estimateMinutes:int` (opc)<br>`milestone:bool` (opc)<br>`recurrenceRule:string` (opc)<br>`labels:string[]` (opc) | `task:WorkspaceTask` (wym) |
| `workspace.task.update` | Zmienia zadanie projektu. Pola pominięte zostają bez zmiany — wywołanie nie jest podmianą całego zadania | `taskId:string` (wym)<br>`title:string` (opc)<br>`description:string` (opc)<br>`status:WorkspaceTaskStatus` (opc)<br>`priority:WorkspaceTaskPriority` (opc)<br>`assigneeKind:WorkspaceAssigneeKind` (opc)<br>`assigneeId:string` (opc)<br>`startAt:int64` (opc)<br>`dueAt:int64` (opc)<br>`estimateMinutes:int` (opc)<br>`spentMinutes:int` (opc)<br>`progressPercent:int` (opc)<br>`milestone:bool` (opc)<br>`recurrenceRule:string` (opc)<br>`labels:string[]` (opc)<br>`checklist:WorkspaceTaskChecklistItem[]` (opc) | `task:WorkspaceTask` (wym)<br>`rescheduledTaskIds:string[]` (opc) |
| `workspace.task.delete` | Usuwa zadanie projektu wraz z jego listą kontrolną i zależnościami | `taskId:string` (wym) | `deleted:bool` (wym)<br>`deletedSubtaskIds:string[]` (opc) |
| `workspace.task.list` | Zwraca zadania projektu dla widoku listy węzła planowania | `projectId:string` (wym)<br>`status:WorkspaceTaskStatus` (opc)<br>`assigneeId:string` (opc)<br>`label:string` (opc)<br>`query:string` (opc)<br>`dueBefore:int64` (opc)<br>`includeDone:bool` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `tasks:WorkspaceTask[]` (wym)<br>`total:int` (wym) |
| `workspace.task.move` | Przenosi kartę zadania na tablicy — między kolumnami albo w obrębie kolumny. Osobno od `workspace.task.update`, bo przeciągnięcie karty zmienia stan i porządek naraz, a jest czynnością wykonywaną jednym ruchem myszy | `taskId:string` (wym)<br>`boardColumnId:string` (opc)<br>`status:WorkspaceTaskStatus` (opc)<br>`beforeTaskId:string` (opc)<br>`afterTaskId:string` (opc) | `task:WorkspaceTask` (wym)<br>`wipExceeded:bool` (opc) |
| `workspace.task.dependency.set` | Zakłada zależność między dwoma zadaniami. Rdzeń odmawia założenia zależności domykającej cykl, bo cyklu nie da się ułożyć w czasie | `projectId:string` (wym)<br>`predecessorTaskId:string` (wym)<br>`successorTaskId:string` (wym)<br>`kind:WorkspaceDependencyKind` (opc)<br>`lagMinutes:int` (opc) | `dependency:WorkspaceTaskDependency` (wym)<br>`rescheduledTaskIds:string[]` (opc) |
| `workspace.task.dependency.remove` | Znosi zależność między zadaniami | `dependencyId:string` (wym) | `removed:bool` (wym) |
| `workspace.board.get` | Zwraca tablicę zadań projektu — kolumny wraz z kartami w kolejności | `projectId:string` (wym)<br>`includeDone:bool` (opc) | `board:WorkspaceBoard` (wym) |
| `workspace.schedule.get` | Zwraca harmonogram projektu dla widoku osi czasu — słupki zadań, zależności i ścieżkę krytyczną. Zadanie bez obu granic czasu do harmonogramu nie wchodzi | `projectId:string` (wym)<br>`from:int64` (opc)<br>`to:int64` (opc) | `bars:WorkspaceScheduleBar[]` (wym)<br>`dependencies:WorkspaceTaskDependency[]` (wym)<br>`unscheduledTaskIds:string[]` (opc) |
| `workspace.calendar.get` | Zwraca pozycje kalendarza projektu we wskazanej siatce | `projectId:string` (wym)<br>`span:WorkspaceCalendarSpan` (wym)<br>`anchorAt:int64` (wym) | `entries:WorkspaceCalendarEntry[]` (wym)<br>`from:int64` (wym)<br>`to:int64` (wym) |
| `workspace.calendar.import` | Wciąga kalendarz iCal do projektu jako zadania z terminem | `projectId:string` (wym)<br>`content:string` (wym)<br>`fileName:string` (opc)<br>`asTasks:bool` (opc) | `imported:int` (wym)<br>`skipped:int` (wym)<br>`skippedReasons:string[]` (opc)<br>`taskIds:string[]` (opc) |
| `workspace.calendar.export` | Zapisuje kalendarz projektu w zapisie iCal | `projectId:string` (wym)<br>`from:int64` (opc)<br>`to:int64` (opc) | `content:string` (wym)<br>`fileName:string` (wym)<br>`exported:int` (wym) |
| `workspace.note.save` | Zapisuje notatkę projektu. Puste `noteId` zakłada notatkę nową; podane zmienia istniejącą. Rdzeń przy zapisie przelicza odnośniki treści | `projectId:string` (wym)<br>`noteId:string` (opc)<br>`title:string` (wym)<br>`content:string` (wym)<br>`parentNoteId:string` (opc)<br>`tags:string[]` (opc) | `note:WorkspaceNote` (wym)<br>`linkedNames:string[]` (opc)<br>`missingNames:string[]` (opc) |
| `workspace.note.get` | Zwraca notatkę projektu wraz z jej treścią | `noteId:string` (wym) | `note:WorkspaceNote` (wym) |
| `workspace.note.list` | Zwraca notatki projektu | `projectId:string` (wym)<br>`query:string` (opc)<br>`tag:string` (opc)<br>`parentNoteId:string` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `notes:WorkspaceNote[]` (wym)<br>`total:int` (wym) |
| `workspace.note.delete` | Usuwa notatkę projektu | `noteId:string` (wym)<br>`withChildren:bool` (opc) | `deleted:bool` (wym)<br>`deletedNoteIds:string[]` (opc)<br>`orphanedBacklinkCount:int` (opc) |
| `workspace.note.tree.get` | Zwraca drzewo stron projektu do nawigacji zakładki „Notatki i wiki” | `projectId:string` (wym) | `nodes:WorkspaceNoteNode[]` (wym) |
| `workspace.note.backlink.list` | Zwraca odnośniki wsteczne strony — panel „Co odsyła tutaj” | `noteId:string` (wym)<br>`limit:int` (opc)<br>`offset:int` (opc) | `backlinks:WorkspaceBacklink[]` (wym)<br>`total:int` (wym) |
| `workspace.knowledge.graph.get` | Zwraca graf wiedzy projektu. Nazwa rodziny `workspace.knowledge` nie miesza się z rodziną `knowledge.*`: tamta prowadzi wskaźnik znaczenia, ta rysuje sieć odnośników między bytami projektu | `projectId:string` (wym)<br>`rootId:string` (opc)<br>`rootKind:WorkspaceEntityKind` (opc)<br>`depth:int` (opc)<br>`kinds:WorkspaceEntityKind[]` (opc)<br>`maxNodes:int` (opc) | `graph:WorkspaceKnowledgeGraph` (wym) |
| `workspace.canvas.get` | Zwraca tablicę wizualną projektu | `projectId:string` (wym)<br>`canvasId:string` (opc) | `canvas:WorkspaceCanvas` (opc)<br>`canvasIds:string[]` (opc) |
| `workspace.canvas.save` | Zapisuje tablicę wizualną projektu | `projectId:string` (wym)<br>`canvasId:string` (opc)<br>`name:string` (opc)<br>`scene:json` (wym) | `canvas:WorkspaceCanvas` (wym) |
| `workspace.library.text.extract` | Wydobywa tekst z pliku projektu i dopisuje go do wskaźnika wyszukiwania. Obejmuje rozpoznanie tekstu z obrazów i skanów | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`method:WorkspaceExtractionMethod` (opc)<br>`languages:string[]` (opc)<br>`force:bool` (opc) | `extraction:WorkspaceTextExtraction` (wym) |
| `workspace.library.duplicate.list` | Wykrywa pliki projektu o treści identycznej. Sam wykaz niczego nie scala | `projectId:string` (wym)<br>`minSizeBytes:int64` (opc) | `groups:WorkspaceDuplicateGroup[]` (wym)<br>`total:int` (wym)<br>`reclaimableBytes:int64` (opc) |
| `workspace.library.duplicate.merge` | Scala duplikaty w jeden plik. Etykiety, kolekcje i wersje plików scalanych przechodzą na plik zachowany, żeby scalenie nie gubiło dorobku | `projectId:string` (wym)<br>`keepFileId:string` (wym)<br>`mergedFileIds:string[]` (wym) | `file:LibraryFile` (wym)<br>`merged:int` (wym)<br>`reclaimedBytes:int64` (opc) |
| `workspace.search.project` | Przeszukuje projekt po słowach — zadania, notatki, pliki, wpisy pamięci i instrukcje jednym pytaniem. Nie zastępuje `library.file.search`: tamta komenda przeszukuje bibliotekę centralną, ta jeden projekt i więcej niż pliki. Wyszukiwanie po znaczeniu prowadzi rodzina `knowledge.*` | `projectId:string` (wym)<br>`query:string` (wym)<br>`kinds:WorkspaceEntityKind[]` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `hits:WorkspaceSearchHit[]` (wym)<br>`total:int` (wym) |
| `workspace.activity.list` | Zwraca oś czasu aktywności projektu — chronologiczny zapis zdarzeń | `projectId:string` (wym)<br>`kinds:WorkspaceActivityKind[]` (opc)<br>`entityKind:WorkspaceEntityKind` (opc)<br>`entityId:string` (opc)<br>`from:int64` (opc)<br>`to:int64` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `entries:WorkspaceActivityEntry[]` (wym)<br>`total:int` (wym) |
| `workspace.comment.add` | Dokłada komentarz przy zadaniu, notatce, pliku albo wpisie pamięci | `projectId:string` (wym)<br>`targetKind:WorkspaceEntityKind` (wym)<br>`targetId:string` (wym)<br>`content:string` (wym)<br>`parentCommentId:string` (opc) | `comment:WorkspaceComment` (wym)<br>`mentionedIds:string[]` (opc) |
| `workspace.comment.list` | Zwraca komentarze projektu albo wątek jednego bytu | `projectId:string` (wym)<br>`targetKind:WorkspaceEntityKind` (opc)<br>`targetId:string` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `comments:WorkspaceComment[]` (wym)<br>`total:int` (wym) |
| `workspace.comment.delete` | Usuwa komentarz | `commentId:string` (wym)<br>`withReplies:bool` (opc) | `deleted:bool` (wym)<br>`deletedCommentIds:string[]` (opc) |
| `workspace.project.list` | Zwraca projekty przestrzeni roboczej — wykaz do wyboru projektu w oknie Project Dashboard i w bocznej nawigacji | `query:string` (opc)<br>`includeArchived:bool` (opc)<br>`includeDeleted:bool` (opc)<br>`limit:int` (opc) | `projects:WorkspaceProject[]` (wym)<br>`total:int` (wym) |
| `workspace.project.create` | Zakłada projekt przestrzeni roboczej. Projekt jest operacyjny od razu, przed uzupełnieniem instrukcji, pamięci i biblioteki | `name:string` (wym)<br>`description:string` (opc)<br>`ownerNote:string` (opc)<br>`owner:string` (opc)<br>`projectId:string` (opc) | `project:WorkspaceProject` (wym) |
| `workspace.project.update` | Zmienia nazwę, opis, właściciela albo notatkę właściciela projektu; pola pominięte zostają bez zmiany | `projectId:string` (wym)<br>`name:string` (opc)<br>`description:string` (opc)<br>`ownerNote:string` (opc)<br>`owner:string` (opc) | `project:WorkspaceProject` (wym) |
| `workspace.project.duplicate` | Klonuje projekt: instrukcje, wpisy pamięci, przypisania ekspertów i zadania trafiają do projektu nowego pod nową nazwą | `projectId:string` (wym)<br>`name:string` (opc) | `project:WorkspaceProject` (wym)<br>`copiedTaskCount:int` (wym)<br>`copiedMemoryEntryCount:int` (wym) |
| `workspace.project.delete` | Usuwa projekt miękko: znika z wykazu i z pulpitu, a dane zostają do chwili cofnięcia komendą `workspace.project.restore` | `projectId:string` (wym) | `project:WorkspaceProject` (wym)<br>`deleted:bool` (wym) |
| `workspace.project.restore` | Cofa miękkie usunięcie projektu — projekt wraca do wykazu ze wszystkimi danymi | `projectId:string` (wym) | `project:WorkspaceProject` (wym)<br>`restored:bool` (wym) |
| `workspace.project.export` | Składa jednoplikowe podsumowanie stanu projektu w Markdown — zadania, pliki, ustalenia, ekspertów i oś czasu — i odkłada je jako plik biblioteki projektu | `projectId:string` (wym) | `fileName:string` (wym)<br>`content:string` (wym)<br>`file:LibraryFile` (opc) |
| `workspace.instructions.get` | Zwraca instrukcje systemowe obowiązujące w projekcie po rozstrzygnięciu warstw oraz treść własną wskazanej warstwy zapisu | `projectId:string` (wym)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `instructions:WorkspaceInstructions` (wym)<br>`layerContent:string` (wym)<br>`layerPresent:bool` (wym) |
| `workspace.agent.assignment.list` | Zwraca przypisania ekspertów do projektu wraz z rolą, wskazaniem domyślnego wykonawcy, zakresem dostępu do okien i zasięgiem przypisania | `projectId:string` (wym) | `assignments:WorkspaceAgentAssignment[]` (wym) |
| `workspace.agent.access.set` | Ustala, do których okien projektu przypisany ekspert ma dostęp, oraz zasięg przypisania; klucze okien są jawne | `projectId:string` (wym)<br>`agentId:string` (wym)<br>`windowAccess:string[]` (wym)<br>`scope:ConfigScope` (opc) | `assignment:WorkspaceAgentAssignment` (wym) |
| `workspace.library.upload` | Wgrywa plik do biblioteki projektu — do katalogu roboczego projektu, skąd czyta go `workspace.library.list` | `projectId:string` (wym)<br>`name:string` (wym)<br>`contentBase64:string` (wym)<br>`path:string` (opc) | `file:LibraryFile` (wym) |
| `workspace.library.delete` | Usuwa pliki z biblioteki projektu | `projectId:string` (wym)<br>`fileIds:string[]` (wym) | `deletedCount:int` (wym)<br>`deletedFileIds:string[]` (wym) |
| `workspace.library.preview` | Zwraca podgląd pliku biblioteki projektu bez opuszczania modułu: tekst dla dokumentów, obraz dla grafik | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`maxChars:int` (opc) | `file:LibraryFile` (wym)<br>`kind:LibraryPreviewKind` (wym)<br>`text:string` (opc)<br>`contentBase64:string` (opc)<br>`truncated:bool` (opc) |
| `workspace.library.move` | Przenosi albo zmienia nazwę pliku biblioteki projektu w obrębie katalogu projektu | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`path:string` (wym) | `file:LibraryFile` (wym) |
| `workspace.library.tag.set` | Nadaje etykiety plikowi biblioteki projektu; zestaw zastępuje poprzedni | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`tags:string[]` (wym) | `file:LibraryFile` (wym) |
| `workspace.context.search` | Odnajduje wpisy pamięci projektu znaczeniowo bliskie zapytaniu. Różni się od `workspace.context.get`: tamta komenda zawęża wykaz frazą występującą w treści, ta liczy podobieństwo znaczeniowe wpisu do zapytania i oddaje wpisy uszeregowane według niego, także gdy żadne słowo zapytania nie stoi we wpisie dosłownie | `projectId:string` (wym)<br>`query:string` (wym)<br>`limit:int` (opc)<br>`minScore:int` (opc)<br>`includeShared:bool` (opc) | `matches:WorkspaceMemoryMatch[]` (wym)<br>`searchedCount:int` (wym) |
| `workspace.library.version.list` | Zwraca historię wersji pliku biblioteki projektu, od najnowszej. Wersja powstaje przy każdym zapisie pliku w projekcie i niesie własną treść, więc historia jest odtwarzalna bez systemu kontroli wersji | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`limit:int` (opc) | `versions:WorkspaceFileVersion[]` (wym) |
| `workspace.library.version.compare` | Porównuje dwie wersje dokumentu projektu wiersz po wierszu. Brak drugiej wersji porównuje wskazaną wersję z treścią bieżącą pliku | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`versionId:string` (wym)<br>`otherVersionId:string` (opc) | `lines:WorkspaceDiffLine[]` (wym)<br>`addedCount:int` (wym)<br>`removedCount:int` (wym)<br>`binary:bool` (opc) |
| `workspace.library.version.restore` | Przywraca treść wskazanej wersji jako bieżącą treść pliku projektu. Treść sprzed przywrócenia zostaje odłożona jako kolejna wersja, więc przywrócenie samo jest odwracalne | `projectId:string` (wym)<br>`fileId:string` (wym)<br>`versionId:string` (wym) | `file:LibraryFile` (wym)<br>`version:WorkspaceFileVersion` (wym) |

**Zdarzenia obszaru `workspace` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `workspace.project.changed` | Zmiana projektu przestrzeni roboczej | — |


Razem w wykazie: **59 komend** z jednego obszaru kontraktu.

---

*Koniec dokumentu. Środowisko WorkSpace — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
