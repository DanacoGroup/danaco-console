# Danaco Console — Środowisko WorkSpace: dokumentacja projektowa

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
| **Tytuł** | Środowisko WorkSpace — powłoka środowiska, komplet okien, warstwy widoczności, nawigacja modułów, karty sesji, układy okien, współdzielenie kontekstu, punkty sterowania |
| **Przeznaczenie** | Materiał wykonawczy do projektowania i budowy interfejsu środowiska WorkSpace: co powstaje, gdzie leży, w jakiej formie i do czego służy — dla każdego okna środowiska, każdego elementu powłoki i każdego przepływu pracy |
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Środowisko** | WorkSpace — środowisko produktywności, organizacji i realizacji projektów |
| **Zakres dokumentu** | Powłoka środowiska WorkSpace: obszar roboczy jako kontener, okna komunikacji operacyjnej, boczna nawigacja modułów, karty sesji, pasek kontekstu środowiska, układy okien, magistrala kontekstu, funkcje wspólne, personalizacja z okna konfiguracji |
| **Odesłanie** | Wnętrze modułu Workspace — jego okna operacyjne, projekty, zadania, pamięć kontekstowa, biblioteka projektu i Agent Manager — opisuje `moduly/workspace.md`. Niniejszy dokument opisuje środowisko i nie powtarza treści dokumentu modułu |
| **Stos techniczny odniesienia** | Powłoka kliencka: Tauri z silnikiem prezentacji, tokeny `tokens.css`, komponenty `.dn-*` z `components.css`. Kanał sterujący: WebSocket z kluczami jawnymi (`shell.layout.save`, `shell.tab.group`, `shell.context.handoff`). Rdzeń serwera: Go — procesy sesji, trwałość układów, indeks wyszukiwania SQLite FTS5, dopasowanie rozmyte `sahilm/fuzzy`, liczenie tokenów `pkoukk/tiktoken-go`. Trwałość stanu powłoki: SQLite |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie; klucze jawne; każdy stan domyślny jest wyjściowy i zmienialny z okna konfiguracji; brak ustawienia oznacza wartość domyślną, nigdy blokadę |

---

## Spis treści

1. [Przeznaczenie i kontekst środowiska](#1-przeznaczenie-i-kontekst-środowiska)
2. [Moduły dostępne w środowisku](#2-moduły-dostępne-w-środowisku)
3. [Komplet okien środowiska](#3-komplet-okien-środowiska)
4. [Specyfikacja okien](#4-specyfikacja-okien)
5. [Przepływy pracy](#5-przepływy-pracy)
6. [Stany i powiązania](#6-stany-i-powiązania)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
8. [Scenariusze użycia](#8-scenariusze-użycia)
- [Załącznik A. Skróty klawiszowe środowiska](#załącznik-a-skróty-klawiszowe-środowiska)

---

## 1. Przeznaczenie i kontekst środowiska

### 1.1. Rola środowiska

Środowisko odpowiada na pytanie „w jakim trybie pracuję?", moduł — na pytanie „jakie zadanie wykonuję?". WorkSpace jest środowiskiem produktywności, organizacji i realizacji projektów: przekształca zamierzenia w zorganizowane działania i produkty. Istotne są w nim struktura projektu, powtarzalność procesów oraz koordynacja zadań w czasie.

Powłoka środowiska WorkSpace to warstwa nadrzędna wobec modułów — wszystko, co otacza okna operacyjne i pozostaje na ekranie niezależnie od tego, który moduł jest otwarty w bieżącej karcie: pasek górny z paskiem kontekstu, pasek kart sesji, boczna nawigacja modułów, obszar roboczy rozumiany jako kontener oraz dwa okna komunikacji operacyjnej — Chat Window i Execution Loop Window. Powłoka pełni funkcję, którą w osobnym systemie operacyjnym pełnią menedżer okien, menedżer sesji i pulpit.

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
    │       pasek górny · pasek kontekstu · karty sesji ·
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
 ══════════╤═════════════════════╤══════════════════════════╤═════════════
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
 ══════════╧═════════════════════╧══════════════════════════╧═════════════
```

---

## 2. Moduły dostępne w środowisku

Boczna nawigacja środowiska WorkSpace udostępnia dziewięć modułów. Kolejność źródłowa i przynależność do grup podlegają konfiguracji z okna konfiguracji.

| Moduł | Rola w środowisku WorkSpace | Grupa domyślna |
|---|---|---|
| Studio | Redakcja i tworzenie treści projektu | Tworzenie |
| Workspace | Zarządzanie projektem: projekty, zadania, pamięć kontekstowa, biblioteka projektu, agenci projektu — pełna specyfikacja w `moduly/workspace.md` | Organizacja |
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
| Wyszukiwarka środowiska | Odnajdywanie kart, modułów, projektów, materiałów Library, zapisanych sesji i układów | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Pole wyszukiwania paska górnego |
| Podgląd polityki kontekstu | Prezentacja zasięgu kontekstu obowiązującego bieżącą kartę przed wysłaniem polecenia | Panel popover przy pasku kontekstu | 2 | Kliknięcie znacznika polityki kontekstu |
| Menedżer układów okien | Zapis, wczytanie i reset układu okien; podział kolumn, dokowanie, tryb skupienia | Panel popover nad obszarem roboczym | 3 | Menu kebab (⋮) nagłówka obszaru roboczego |
| Przełącznik i wyszukiwarka kart | Lista wszystkich kart z filtrem tekstowym i miniaturą, lista ostatnio zamkniętych | Panel popover przy pasku kart | 3 | Menu kebab (⋮) paska kart |
| Menedżer sesji i grup kart | Zapisane przestrzenie projektu, grupy kart, przypięcia | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu kart sesji |
| Magistrala kontekstu | Bufor roboczy na fragmenty, pliki i odnośniki krążące między kartami i modułami | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu hamburger (☰) bocznej nawigacji; upuszczenie materiału na krawędź obszaru |
| Panel przekazania artefaktu | Kierowanie wyniku okna do innego modułu, do Library albo do projektu w module Workspace | Panel popover przy elemencie wyniku | 3 | Menu kontekstowe wyniku |
| Centrum powiadomień | Rejestr zdarzeń środowiska; mechanizm w postaci jednolitej dla całej platformy — `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6 | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Plakietka powiadomień z licznikiem w pasku kontekstu |
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
  Browser  │ [@][📎][Szablony ▾]          │
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
  Browser  │                      │  tygodniowy"      │
  Research │                      │ ├ 1 Materiał   ✔  │
  Library  │                      │ ├ 2 Analiza    ▶  │
  Roundtab.│                      │ ├ 3 Redakcja   ⏳ │
  Design   │                      │ └ 4 Kontrola   ⏳ │
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
| Podgląd modułu | Mini-podgląd zestawu okien i opisu tematycznego modułu przed przeładowaniem karty | 2 | Najechanie z przytrzymaniem |
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
| Podgląd pracy w tle | Dymek z opisem trwającego procesu, postępem i przejściem do Execution Loop Window albo do funkcji Mobile | 2 | Najechanie na wskaźnik karty |
| Zmiana kolejności kart | Przeciąganie karty w obrębie paska i między grupami | 2 | Przeciągnięcie |
| Menu karty | Zamknij inne, zamknij po prawej, przypnij, dodaj do grupy, zmień nazwę, duplikuj | 3 | Menu kontekstowe karty |
| Ręczna nazwa karty | Własna etykieta zamiast nazwy modułu, na przykład „Studio — rozdział 3" | 3 | Menu kontekstowe karty |
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
| Polityka kontekstu | Zasięg kontekstu karty: odrębny albo współdzielony; kliknięcie otwiera podgląd „co widzi ta rozmowa" | 2 | Kliknięcie znacznika |
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
| Zakres wyszukiwania | Karty, moduły, projekty, materiały Library, zapisane sesje i układy; indeks po stronie Go, SQLite FTS5 | 2 | Pole wyszukiwania paska górnego |
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

Mechanizm magistrali kontekstu — przeznaczenie, sposób przekazania artefaktu, zasięg, warstwę widoczności i sposób wywołania oraz relacje do modułu Library, punktów izolacji i encji artefaktu — opisuje `interfejs-uzytkownika/przeplyw-okien.md` (rozdz. 6a). Specyfika środowiska WorkSpace: zasięg domyślny magistrali obejmuje grupę kart projektu, a moduł Workspace jest celem odbioru artefaktów kierowanych do projektu.

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
| Magistrala kontekstu o zasięgu projektu | Przypięte materiały, ustalenia i pamięć dostępne każdej karcie grupy projektu na poziomie zasięgu „projekt" | 2 | Znacznik polityki kontekstu, panel szybkiej konfiguracji sesji |
| Współdzielenie historii i pamięci między kartami | Ustanowienie współdzielenia historii, pamięci albo kontekstu między dwiema wskazanymi kartami | 2 | Panel szybkiej konfiguracji sesji |
| Podgląd polityki efektywnej | Prezentacja zasięgu obowiązującego bieżącą kartę oraz treści widocznej dla rozmowy przed wysłaniem polecenia | 2 | Kliknięcie znacznika polityki kontekstu |
| Odesłanie do okna punktów izolacji | Przejście do pełnej definicji profili i macierzy zasięgu — powłoka reguły stosuje i prezentuje, definiuje je okno punktów izolacji | 3 | Odnośnik panelu |

### 4.10. Centrum powiadomień i panel szybkiej konfiguracji sesji

Mechanizm centrum powiadomień — taksonomię klas zdarzeń, postać wizualną, umiejscowienie w układzie pionowym, warstwę widoczności, sposób wywołania, stany oraz działania dostępne z poziomu powiadomienia — opisuje karta komponentu w `interfejs-uzytkownika/katalog-komponentow.md` (rozdz. 11.6). Poniżej podano wyłącznie to, co właściwe środowisku WorkSpace.

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
| Powrót na stronę główną z zachowaniem tła | Wyjście ze środowiska; karty trwają w tle i odtwarzają pełny stan przy powrocie | 1 | Ikona domu paska górnego |
| Punkt dostępu Always On Display | Przywołanie pływającego agenta towarzyszącego ponad bieżącą przestrzenią, w każdej karcie i module | 2 | Ikona paska górnego, skrót |
| Punkt dostępu Mobile | Wejście w funkcję globalną Mobile — nadzór nad procesami środowiska działającymi w tle; nie tworzy nowej przestrzeni roboczej | 2 | Ikona paska górnego |

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
| Moduł Workspace | Projekt nadrzędny grupy kart, cel odbioru artefaktów, źródło terminów i zadań prezentowanych w centrum powiadomień. Zakres funkcjonalny modułu — `moduly/workspace.md` | `shell.context.sink`, znacznik projektu |
| Moduł Library | Trwałe repozytorium środowiska, cel odbioru artefaktów i zakres wyszukiwania środowiska | `shell.context.sink`, `shell.search.scope` |
| Moduł Agents | Źródło wykonawców przypisywanych zadaniom pętli wykonawczej | Selektor wykonawcy w pasku kontekstu |
| Moduł Automations | Harmonogram przebiegów pętli wykonawczej i praca cykliczna | Panel harmonogramu Execution Loop Window |
| Okno punktów izolacji | Definicja profili i macierzy zasięgu kontekstu; powłoka reguły stosuje i prezentuje | `shell.context.scope` |
| Always On Display | Pływający agent towarzyszący ponad bieżącą przestrzenią, obecny w każdej karcie i module | Punkt dostępu paska górnego |
| Mobile | Nadzór i interwencja zdalna nad procesami środowiska działającymi w tle | Punkt dostępu paska górnego |
| Środowiska TalkIn i CodeStudio | Przelanie karty modułu wspólnego między środowiskami | Paleta poleceń |

---

## 7. Punkty sterowania z okna konfiguracji

Wszystkie funkcje powłoki środowiska są sterowalne z okna konfiguracji, w zakresie „Aplikacja i procesy" oraz w warstwie sesji. Każdy element ma objaśnienie kontekstowe `[?]`. Klucze są jawne, brak ustawienia oznacza wartość domyślną.

### 7.1. Mapa punktów sterowania

| Obszar | Punkt sterowania | Klucz | Wartość domyślna |
|---|---|---|---|
| Nawigacja | Grupy, kolejność i widoczność modułów | `shell.nav.groups`, `shell.nav.order`, `shell.nav.hidden` | Dziewięć modułów, kolejność źródłowa, wszystkie widoczne |
| Nawigacja | Przypięte moduły | `shell.nav.pinned` | Brak przypięć |
| Nawigacja | Zwinięcie panelu do ikon | `shell.nav.rail` | Panel rozwinięty |
| Nawigacja | Zachowanie kliknięcia modułu: przeładowanie karty albo nowa karta | `shell.nav.open_mode` | Przeładowanie bieżącej karty |
| Układ okien | Zapisane układy i układ domyślny per moduł | `shell.layout.presets`, `shell.layout.default` | Układ zdefiniowany przez moduł |
| Układ okien | Siatka przyciągania | `shell.layout.snap` | Włączona |
| Układ okien | Szerokość kolumny Chat Window | `shell.chat.width` | Kolumna lewa, szerokość podstawowa |
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
    zasieg:           odrebny         # poziom „projekt" ustawiany świadomie
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

Każda sprawa ma własną grupę kart: karta modułu Workspace ze sprawą, karta Research z materiałem źródłowym, karta Studio z redakcją pisma. Karta sprawy wiodącej pozostaje przypięta. Zlecenie „przygotuj projekt pisma na podstawie ustaleń" przechodzi z Chat Window do Koordynatora; Execution Loop Window prowadzi zadania: zebranie materiału w module Research, redakcja w module Studio, kontrola jakości wobec kryterium ukończenia. Gotowe pismo trafia do projektu w module Workspace, dziennik przebiegu — do akt sprawy.

### 8.2. Zespół produktowy prowadzący projekt od materiału do produktu

Grupa kart projektu obejmuje moduły Research, Design, Studio, Apps i Workspace. Magistrala kontekstu o zasięgu projektu udostępnia każdej karcie przypięte materiały i ustalenia projektu. Makieta z modułu Design przekazywana jest do modułu Apps przekazaniem artefaktu, treść z modułu Studio — tą samą drogą. Execution Loop Window prowadzi budowę produktu jako jedną pętlę obejmującą zadania w trzech modułach, a wskaźniki przebiegu widoczne są na kartach i przy pozycjach modułów.

### 8.3. Praca ciągła nad projektem cyklicznym

Przebieg pętli wykonawczej wpięty jest w automatykę modułu Automations i realizuje cykliczne zadanie projektowe. Użytkownik nadzoruje przebieg przez funkcję Mobile: odbiera zdarzenia z centrum powiadomień, zatwierdza kroki oczekujące na decyzję i wstrzymuje przebieg w razie zgłoszenia przeszkody. Po powrocie do stanowiska zapisana przestrzeń projektu odtwarza pełny układ pracy, a Execution Loop Window prezentuje historię przebiegów.

### 8.4. Praca głęboko skupiona nad jednym materiałem

Użytkownik wczytuje zapisany układ „Redakcja": obszar roboczy podzielony na kolumnę materiału źródłowego i kolumnę redakcyjną, Chat Window w kolumnie lewej, Execution Loop Window zwinięta do nagłówka pętli. Tryb czystego ekranu ukrywa nawigację i karty. Wyjście z trybu oraz powrót do pełnego środowiska następuje jednym skrótem, bez zmiany kontekstu karty.

### 8.5. Prowadzenie jednego projektu na wielu urządzeniach

Ten sam projekt otwierany jest na komputerze i na tablecie. Pamięć układu per klasa urządzenia utrzymuje odrębne rozmieszczenie kolumn dla każdego ekranu przy zachowaniu tej samej przestrzeni projektu, tych samych kart i tego samego kontekstu. Procesy trwają w tle niezależnie od urządzenia, z którego prowadzony jest nadzór.

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
| Dostępność | `Tab` | Kolejność fokusu: pasek górny → karty sesji → boczna nawigacja → Chat Window → Execution Loop Window → obszar roboczy → kolumny boczne |

---

*Koniec dokumentu. Środowisko WorkSpace — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
