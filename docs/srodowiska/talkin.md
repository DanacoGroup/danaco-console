# Danaco Console — Środowisko TalkIn: dokumentacja projektowa

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
| **Tytuł** | Środowisko TalkIn — dokumentacja projektowa: przeznaczenie, moduły, komplet okien, specyfikacja okien, przepływy pracy, stany i powiązania, punkty sterowania, scenariusze użycia |
| **Przeznaczenie** | Pełnozakresowy materiał źródłowy do projektowania i budowy interfejsu środowiska TalkIn: co istnieje, gdzie leży, w jakiej formie i do czego służy — dla każdego okna, każdego elementu i każdego przepływu pracy |
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Środowisko** | TalkIn — środowisko wiedzy, komunikacji i pracy z treścią (Koncepcja platformy, rozdz. 9.1) |
| **Autorytatywne źródło** | architektura/koncepcja-platformy.md |
| **Źródła pomocnicze** | interfejs-uzytkownika/strona-glowna-i-nawigacja.md · specyfikacje/specyfikacja-okien-operacyjnych.md · specyfikacje/specyfikacja-modulow.md · architektura/model-konfiguracji.md · interfejs-uzytkownika/system-wizualny.md · architektura/izolacja-i-zaleznosci.md · interfejs-uzytkownika/elementy-okien.md · interfejs-uzytkownika/przeplyw-okien.md |
| **Opracowanie** | Danaco Console — dokumentacja projektu UI |

**Zasada nadrzędna obowiązująca cały dokument.** Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód. Domyślne zachowanie systemu to wykonanie polecenia. Żaden element interfejsu nie traci klikalności; niekompletność sygnalizowana jest komunikatem kontekstowym, nie wyłączeniem kontrolki. Izolacja techniczna, izolacja kontekstu i wszelkie ograniczenia uprawnień są ustawieniami konfiguracyjnymi Operatora ze stanem wyjściowym „pełny dostęp”.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
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

## Wprowadzenie

TalkIn jest środowiskiem wiedzy, komunikacji i pracy z treścią. Jednostką organizacji pracy jest użytkownik prowadzący dialog z Wykonawcą nad treścią: jej tworzeniem, analizą, przekształcaniem i pozyskiwaniem ze źródeł. Środowisko udostępnia dziewięć modułów operujących na wspólnym przedmiocie — treści — oraz środowiskową warstwę spajającą, która pozwala jednemu zasobowi krążyć między modułami bez kopiowania.

Dwa okna komunikacji operacyjnej stanowią elementy pierwszoplanowe środowiska. **Chat Window** — kanał Użytkownik ↔ Wykonawca — jest centralnym punktem pracy i podstawowym mechanizmem sterowania wszystkimi procesami realizowanymi przez platformę; w TalkIn, środowisku z natury komunikacyjnym, zajmuje lewą kolumnę obszaru roboczego na pełną wysokość i jest obecne w każdym module w tym samym miejscu układu. **Execution Loop Window** — kanał Koordynator ↔ Wykonawca — prowadzi pętlę wykonawczą, koordynuje zadania, nadzoruje realizację, orkiestruje działania i kontroluje przebieg procesów komunikacyjnych; otwiera się jako kolumna sąsiadująca z Chat Window.

Dokument czyta się w czterech warstwach nałożonych na siebie:

| Warstwa dokumentu | Co dostarcza | Rozdziały |
|---|---|---|
| Kontekst | Miejsce środowiska w platformie, jego przedmiot pracy i zestaw modułów | 1–2 |
| Forma | Komplet okien i makiety kolumnowe — co gdzie leży w przestrzeni roboczej | 3–4 |
| Zachowanie w czasie | Przepływy pracy, stany, powiązania między oknami i modułami | 5–6 |
| Sterowanie | Punkty konfiguracji, scenariusze użycia, mapa skrótów | 7–8, załącznik A |

---

## 1. Przeznaczenie i kontekst środowiska

### 1.1. Przedmiot pracy środowiska

TalkIn obsługuje zadania, w których centralnym przedmiotem pracy jest treść. Dokument powstały w module Studio, wynik badania z modułu Research, plik z modułu Library, przekład z modułu Translate, wątek z modułu Roundtable i wypowiedź z modułu Assistant są wariantami tego samego zasobu. Powłoka środowiska dostarcza warstwę spajającą, dzięki której zasób przechodzi między modułami jako referencja, a nie jako kopia.

| Wymiar | Charakterystyka TalkIn |
|---|---|
| Przedmiot pracy | Treść — tworzenie, analiza, przekształcanie, pozyskiwanie ze źródeł |
| Jednostka organizacji pracy | Użytkownik wspierany przez Wykonawcę w ramach jednego zadania |
| Zawartość bocznej nawigacji | Lista dziewięciu modułów środowiska |
| Podstawowa jednostka pracy | Moduł → okno operacyjne |
| Kanał sterowania | Chat Window — polecenie w języku naturalnym |
| Kanał nadzoru | Execution Loop Window — pętla wykonawcza Koordynator ↔ Wykonawca |
| Warstwa spajająca | Magistrala kontekstu, kontekst wiedzy środowiska, referencje krzyżowe, oś czasu treści |

### 1.2. Miejsce środowiska w platformie

```
STRONA GŁÓWNA
    │  wybór środowiska w strefie 1
    ▼
TalkIn
    │
    ├── Boczna nawigacja modułów (rozdz. 2)
    │       Studio · Workspace · Browser · Research · Library ·
    │       Translate · Roundtable · Assistant · Agents
    │
    ├── Obszar roboczy karty sesji (rozdz. 3–4)
    │       Chat Window · Execution Loop Window · okna operacyjne modułów ·
    │       okna pomocnicze jako rozszerzenia boczne
    │
    └── Karty sesji
            każda karta = odrębna przestrzeń pracy nad treścią,
            z własnym modułem, układem kolumn i kontekstem

PONAD CAŁĄ STRUKTURĄ
    Always On Display — globalny agent towarzyszący
    Mobile — monitorowanie i kontynuacja pracy poza stanowiskiem
```

### 1.3. Granica powłoki środowiska i modułów

| Warstwa | Za co odpowiada |
|---|---|
| Powłoka środowiska | Nawigacja między modułami, układ i zarządzanie kolumnami okien, sesje i karty, kolumna stanu, skróty klawiszowe, współdzielenie kontekstu, funkcje wspólne, personalizacja środowiska |
| Moduł | Wnętrze okien operacyjnych (Studio Editor, Research Workspace, Library Explorer, Voice Console i pozostałe), narzędzia i przepływy pracy modułu |
| Funkcje globalne | Mobile, Always On Display — działają ponad środowiskami |

### 1.4. Układ obszaru roboczego

Obszar roboczy środowiska zbudowany jest wyłącznie z kolumn sąsiadujących poziomo. Regulacji podlega szerokość kolumn.

| Kolumna | Zawartość |
|---|---|
| Lewa krawędź | Boczna nawigacja modułów |
| Lewa, stała, pełna wysokość | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca |
| Kolumna sąsiadująca (otwierana) | **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca |
| Prawa, dominująca | Obszar roboczy modułu — okno operacyjne aktywnego modułu |
| Kolejne kolumny boczne | Okna pomocnicze otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego |
| Prawa krawędź | Kolumna stanu środowiska |

Makieta środowiska w stanie spoczynku — widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3:

```
 ═══════════════════════════════════════════════════════════════════════════════
  TalkIn ▼ │ Chat Window            │ Obszar roboczy modułu       │ Kolumna
  ─────────│ Użytkownik ↔ Wykonawca │ Studio Editor               │ stanu
  Studio   │                        │                             │
  Workspace│ [Fable 5] [API] ▼      │ nagłówek okna        ⋮      │ ● API
  Browser  │                        │ ─────────────────────────   │ ▮▮▮▯ 62%
  Research │  historia rozmowy      │                             │ ⟳ 2
  Library  │  polecenia i wyniki    │  treść okna operacyjnego    │ ◈ 5
  Translate│                        │                             │ ☰
  Roundtable                        │                             │
  Assistant│ ─────────────────────  │                             │
  Agents   │ pole polecenia   ⋮ ➤   │                             │
  ─────────│                        │                             │
  ☰        │ Execution Loop ▼       │                             │
 ═══════════════════════════════════════════════════════════════════════════════
```

Po otwarciu okna pętli wykonawczej i okna pomocniczego układ przyjmuje postać:

```
 ═══════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window     │ Execution Loop  │ Obszar roboczy │ Panel
  nawigacja │ Użytkownik ↔    │ Koordynator ↔   │ modułu         │ pomocniczy
  modułów   │ Wykonawca       │ Wykonawca       │                │ (rozszerzenie
            │                 │                 │                │  boczne)
            │ historia        │ zlecenie        │ okno operacyjne│
            │ rozmowy         │ i dekompozycja  │ modułu         │ Magistrala
            │                 │ ───────────     │                │ kontekstu
            │                 │ kolejka zadań   │                │
            │                 │ ───────────     │                │
            │ pole polecenia  │ sterowanie      │                │
            │            ⋮ ➤  │ przebiegiem ⋮   │                │
 ═══════════════════════════════════════════════════════════════════════════════
```

### 1.5. Warstwy widoczności w środowisku

Interfejs TalkIn ujawnia możliwości środowiska stopniowo. Jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Dziewięć modułów, komplet okien pomocniczych, warstwa spajająca i pełny zakres personalizacji nie wpływają na postrzeganą prostotę interfejsu — złożoność pozostaje w architekturze.

| Warstwa | Nazwa | Zawartość w środowisku TalkIn | Sposób wywołania |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, okno operacyjne aktywnego modułu, boczna nawigacja modułów, pasek kart sesji, kolumna stanu, znaczniki kontekstu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Execution Loop Window, wybór modelu i kanału, wybór modułu docelowego akcji „Wyślij do…”, wybór presetu układu kolumn, poziom wysiłku Wykonawcy, Magistrala kontekstu | Kliknięcie znacznika kontekstowego, ikona, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Menu środowiska, menu karty sesji, menu układu kolumn, akcje karty zasobu, ustawienia szybkie sesji, warianty eksportu | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Edytor mapy skrótów, okno punktów izolacji, indeks wiedzy środowiska, graf referencji krzyżowych, diagnostyka kanału połączenia, eksport migawki środowiska | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, paleta poleceń, tryb administracyjny, konfiguracja roli |

Mechanizmy ukrywania funkcjonalności zastosowane w środowisku:

- **Menu progresywne.** Wybór modułu docelowego, modelu, kanału i profilu personalizacji prezentowany jest jako jeden element zwinięty (`Moduł ▼`, `Model ▼`, `Profil ▼`).
- **Panele wysuwane.** Magistrala kontekstu, notatnik środowiska, centrum powiadomień i oś czasu treści są kolumnami bocznymi otwieranymi jako rozszerzenia boczne; po zamknięciu znikają całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Presety układu, zapis sesji, eksport, reset układu, ściągawka skrótów i personalizacja występują jako jeden element zbiorczy `Operacje ▼` w menu środowiska.
- **Znaczniki kontekstowe.** Środowisko, moduł, model, kanał i profil personalizacji występują jako lekkie znaczniki w pasku kontekstu Chat Window, na przykład `[TalkIn] [Studio] [Fable 5] [API] [Redaktor]`; kliknięcie znacznika otwiera odpowiedni selektor.

**Zasada jednego kliknięcia.** Każda ukryta funkcja środowiska jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window.

---

## 2. Moduły dostępne w środowisku

### 2.1. Zestaw modułów

Boczna nawigacja środowiska udostępnia dziewięć modułów. Każdy moduł otwiera własne okno operacyjne w kolumnie obszaru roboczego, zachowując Chat Window w lewej kolumnie.

| Moduł | Przedmiot pracy | Okno operacyjne | Grupa nawigacji |
|---|---|---|---|
| **Studio** | Tworzenie i redakcja treści długiej | Studio Editor | Treść |
| **Library** | Repozytorium dokumentów, plików i zbiorów wiedzy | Library Explorer | Treść |
| **Browser** | Praca ze stronami i źródłami sieciowymi | Browser View | Pozyskiwanie |
| **Research** | Badanie zagadnienia, zbieranie i porządkowanie źródeł | Research Workspace | Pozyskiwanie |
| **Translate** | Przekład i praca dwujęzyczna nad treścią | Translate Workbench | Komunikacja |
| **Roundtable** | Debata wielomodelowa i wypracowanie konsensusu | Roundtable Board | Komunikacja |
| **Assistant** | Dialog głosowy i tekstowy nad bieżącą treścią | Voice Console | Komunikacja |
| **Workspace** | Organizacja projektu, zasobów i porządku pracy | Workspace Board | Zaplecze |
| **Agents** | Definicje agentów wykorzystywanych jako Wykonawcy | Agents Registry | Zaplecze |

### 2.2. Boczna nawigacja modułów

Boczna nawigacja jest kolumną na lewej krawędzi środowiska, złożoną z pozycji `.dn-karta--pozycja` pogrupowanych w cztery sekcje zwijalne. Nagłówek kolumny zawiera nazwę środowiska działającą jako przełącznik środowisk.

| Element nawigacji | Co to jest | Do czego służy | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Nagłówek „TalkIn ▼” | Nazwa środowiska jako przycisk | Przełączenie środowiska bez powrotu na stronę główną; karty bieżącego środowiska trwają w tle | 1 | Widoczny; kliknięcie rozwija listę środowisk |
| Sekcja modułów | Zwijalna grupa nazwana: Treść · Pozyskiwanie · Komunikacja · Zaplecze | Skrócenie drogi wzroku do modułu | 1 | Widoczna; kliknięcie nagłówka zwija sekcję |
| Pozycja modułu | Wiersz z ikoną i nazwą modułu | Otwarcie modułu w obszarze roboczym bieżącej karty | 1 | Widoczna; kliknięcie otwiera moduł |
| Sekcja „Ostatnie” | Lista trzech do pięciu modułów odwiedzonych w bieżącej sesji połączenia | Powrót do przełączanego modułu bez szukania na liście | 1 | Widoczna nad grupami |
| Plakietka stanu przy module | Kropka lub liczba przy pozycji | Sygnalizacja aktywności modułu w tle | 1 | Widoczna przy module z aktywnym procesem |
| Podgląd modułu (peek) | Popover z nazwami okien modułu i ostatnim stanem | Decyzja o przełączeniu bez utraty bieżącego widoku | 2 | Najechanie na pozycję modułu |
| Pole filtrowania modułów | Pole `.dn-input` w nagłówku kolumny | Dojście do modułu z klawiatury | 2 | Skrót klawiszowy lub kliknięcie ikony `szukaj` |
| Menu kontekstowe pozycji | Lista akcji: „Otwórz w nowej karcie”, „Otwórz obok”, „Przypnij”, „Ukryj” | Otwarcie modułu bez utraty bieżącej pracy, personalizacja kolejności | 3 | Prawy klawisz myszy na pozycji lub kebab (⋮) pozycji |
| Uchwyt zmiany kolejności | Strefa przeciągania pozycji | Ustawienie własnej kolejności i przypięć modułów | 3 | Przeciągnięcie pozycji |
| Przełącznik trybu kompaktowego (☰) | Kontrolka zwężająca kolumnę do samych ikon | Więcej szerokości dla obszaru roboczego | 2 | Kliknięcie ikony ☰ u podstawy kolumny nawigacji |

### 2.3. Otwieranie modułu w kolumnie obszaru roboczego

Obszar roboczy karty mieści od jednej do trzech kolumn modułowych. Każda kolumna prowadzi niezależny moduł z własnym oknem operacyjnym; Chat Window pozostaje wspólne dla całej karty i zachowuje lewą kolumnę.

| Sposób otwarcia | Skutek |
|---|---|
| Kliknięcie pozycji modułu | Podmiana modułu w aktywnej kolumnie obszaru roboczego |
| „Otwórz obok” z menu kontekstowego | Dodanie nowej kolumny modułowej po prawej stronie aktywnej |
| Przeciągnięcie pozycji na kolumnę | Osadzenie modułu w tej kolumnie |
| Przeciągnięcie pozycji na krawędź obszaru | Utworzenie nowej kolumny modułowej w miejscu zrzutu |
| „Otwórz w nowej karcie” | Utworzenie karty sesji z tym modułem i domyślnym układem kolumn |

---

## 3. Komplet okien środowiska

Środowisko TalkIn udostępnia komplet okien podzielony na cztery grupy: okna komunikacji operacyjnej, okna operacyjne modułów, okna pomocnicze powłoki oraz nakładki. Kolumna „Warstwa” podaje warstwę widoczności, kolumna „Sposób wywołania” — mechanizm otwarcia.

### 3.1. Okna komunikacji operacyjnej

| Okno | Kanał | Rola w środowisku | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| **Chat Window** | Użytkownik ↔ Wykonawca | Centralny punkt pracy i podstawowy mechanizm sterowania wszystkimi procesami środowiska; przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie i przerywanie działań oraz wyjaśnianie wyniku i kontekstu | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji w każdym module środowiska |
| **Execution Loop Window** | Koordynator ↔ Wykonawca | Prowadzenie pętli wykonawczej, koordynacja zadań, nadzór nad realizacją, orkiestracja działań i kontrola realizacji procesów komunikacyjnych | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 2 | Kliknięcie nagłówka „Execution Loop ▼” u podstawy Chat Window, skrót klawiszowy albo polecenie w Chat Window |

### 3.2. Okna operacyjne modułów

| Okno | Moduł | Zawartość | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Studio Editor | Studio | Edytor treści długiej z konspektem, wersjami i podglądem | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Studio w bocznej nawigacji |
| Library Explorer | Library | Drzewo zbiorów, lista dokumentów, podgląd pliku | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Library |
| Browser View | Browser | Widok strony, pasek adresu, historia odwiedzin karty | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Browser |
| Research Workspace | Research | Zapytanie badawcze, lista źródeł, notatki i synteza | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Research |
| Translate Workbench | Translate | Widok dwujęzyczny, słownik terminów, kontrola spójności | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Translate |
| Roundtable Board | Roundtable | Wypowiedzi modeli, panel moderatora, panel konsensusu | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Roundtable |
| Voice Console | Assistant | Nasłuch głosu, transkrypcja, sterowanie syntezą mowy | Prawa kolumna dominująca albo kolumna boczna | 1 | Kliknięcie pozycji Assistant |
| Workspace Board | Workspace | Struktura projektu, zasoby, zadania porządkowe | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Workspace |
| Agents Registry | Agents | Lista definicji agentów, przypisania, kanały modeli | Prawa kolumna dominująca | 1 | Kliknięcie pozycji Agents |

### 3.3. Okna pomocnicze powłoki

| Okno pomocnicze | Zawartość | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Magistrala kontekstu | Referencje do zasobów odłożonych z dowolnego modułu i dostępnych we wszystkich pozostałych | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie licznika `◈` w kolumnie stanu albo skrót klawiszowy |
| Oś czasu treści | Chronologiczny zapis tego, co powstało w bieżącej sesji i w którym module, z przejściem do źródła | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Pozycja „Oś czasu treści” w menu środowiska (☰) |
| Notatnik środowiska | Lekki obszar na notatki i wklejki wspólny dla kart środowiska, przenoszalny do Library i Studio | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Pozycja „Notatnik” w menu środowiska (☰) albo skrót klawiszowy |
| Centrum powiadomień | Rejestr zdarzeń modułów środowiska; mechanizm w postaci jednolitej dla całej platformy — `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6. Specyfika TalkIn: wyzwalacz umieszczony w kolumnie stanu środowiska, zakres źródeł obejmuje wszystkie moduły bieżącej karty | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Plakietka powiadomień z licznikiem w kolumnie stanu |
| Lista procesów w tle | Aktywne procesy kart środowiska z akcjami wstrzymania i podglądu | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie licznika `⟳` w kolumnie stanu |
| Graf referencji krzyżowych | Powiązania „ten wynik powstał z tych źródeł i tego dokumentu” utrzymywane przez powłokę między modułami | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 | Polecenie w Chat Window, paleta poleceń albo akcja karty zasobu |
| Kontekst wiedzy środowiska | Zbiory z Library aktywne jako tło wiedzy dla wszystkich modułów środowiska w bieżącej sesji | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 | Paleta poleceń albo znacznik kontekstu wiedzy w Chat Window |

### 3.4. Nakładki

| Nakładka | Zawartość | Warstwa | Sposób wywołania |
|---|---|---|---|
| Paleta poleceń | Pole wpisania i wykonania dowolnej akcji powłoki: przełącz moduł, otwórz kartę, uruchom preset układu, zmień model, otwórz okno pomocnicze | 4 | Skrót `Ctrl/Cmd + K` albo kliknięcie podpowiedzi skrótu w kolumnie stanu |
| Menedżer kart | Siatka miniatur wszystkich otwartych kart z wyszukiwaniem, zamykaniem zbiorczym i przywracaniem zamkniętych | 3 | Skrót klawiszowy albo pozycja w menu środowiska (☰) |
| Karta zasobu | Jednolity podgląd zasobu z akcjami „otwórz w…”, „dołącz”, „do magistrali” | 3 | Kliknięcie referencji zasobu w dowolnym module albo na magistrali kontekstu |
| Ściągawka skrótów | Lista aktywnych skrótów środowiska generowana z mapy skrótów | 3 | Klawisz `?` |
| Tryb odczytu | Widok treści aktywnego okna operacyjnego w postaci czytelnej, bez chromu interfejsu | 3 | Pozycja w menu środowiska (☰) albo skrót klawiszowy |
| Edytor mapy skrótów | Widok wszystkich skrótów powłoki ze zmianą przypisań i wykrywaniem konfliktów | 4 | Okno konfiguracji, paleta poleceń albo polecenie w Chat Window |
| Okno punktów izolacji | Siedem poziomów zasięgu izolacji dla karty i relacji między modułami | 4 | Kliknięcie znacznika zasięgu izolacji w kolumnie stanu albo okno konfiguracji |

### 3.5. Elementy stałe powłoki

| Element | Zawartość | Waga wizualna | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Pasek górny | Identyfikacja platformy, pole wyszukiwania ze znacznikiem zakresu, ustawienia szybkie sesji, menu środowiska (☰) | Pasek na szczycie okna środowiska | 1 | Widoczny bez interakcji |
| Pasek kart sesji | Karty sesji, przycisk „+”, selektor presetu układu, nawigacja wstecz/dalej | Pasek pod paskiem górnym | 1 | Widoczny bez interakcji |
| Boczna nawigacja modułów | Dziewięć modułów w czterech sekcjach, sekcja „Ostatnie”, przełącznik trybu kompaktowego | Kolumna na lewej krawędzi | 1 | Widoczna bez interakcji |
| Kolumna stanu środowiska | Agregacja stanu warstwy sesji i środowiska; wszystkie pozycje klikalne | Wąska kolumna na prawej krawędzi obszaru roboczego, pełna wysokość | 1 | Widoczna bez interakcji |

---

## 4. Specyfikacja okien

### 4.1. Anatomia wspólna okna

Każde okno środowiska podlega tej samej anatomii co pozostałe okna operacyjne platformy.

```
 ┌─ Nagłówek okna ─────────────────────────────
 │   nazwa okna · znaczniki kontekstu · ⋮
 ├─────────────────────────────────────────────
 │   OBSZAR ZAWARTOŚCI
 │     strumień, treść lub zestawienie —
 │     aktualizacja na żywo kanałem WebSocket
 ├─────────────────────────────────────────────
 │   PASEK AKCJI OKNA
 │     akcje zwinięte do elementu zbiorczego
 ├─────────────────────────────────────────────
 │   STAN PROCESU SESJI (po stronie serwera)
 │     trwały — rozłączenie klienta nie zamyka okna
 └─────────────────────────────────────────────
```

### 4.2. Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca

Chat Window jest elementem pierwszoplanowym środowiska TalkIn. Zajmuje lewą kolumnę obszaru roboczego na pełną wysokość, występuje w tym samym miejscu układu we wszystkich dziewięciu modułach i stanowi podstawowy mechanizm sterowania procesami platformy: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Polecenie wydane w Chat Window uruchamia operacje we wszystkich modułach środowiska oraz w warstwie spajającej.

Makieta w stanie spoczynku:

```
 ┌─ Chat Window · Użytkownik ↔ Wykonawca ──────────
 │  [TalkIn] [Studio] [Fable 5] [API] [Redaktor]  ⋮
 ├─────────────────────────────────────────────────
 │   Użytkownik
 │     polecenie w języku naturalnym
 │
 │   Wykonawca
 │     strumień odpowiedzi i wyników
 │     ◈ referencja zasobu · [Zatwierdź] [Przerwij]
 │
 ├─────────────────────────────────────────────────
 │   pole polecenia                        ⋮   ➤
 ├─────────────────────────────────────────────────
 │   Execution Loop ▼      Koordynator ↔ Wykonawca
 └─────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nagłówek okna | Belka z nazwą kanału i znacznikami kontekstu | Identyfikacja kanału i bieżącego kontekstu pracy | Domyślny | 1 | Widoczny |
| Znacznik środowiska `[TalkIn]` | Lekka etykieta w pasku kontekstu | Potwierdzenie środowiska; wejście do przełącznika środowisk | Domyślny · aktywny | 1 | Widoczny; kliknięcie otwiera selektor |
| Znacznik modułu `[Studio]` | Lekka etykieta w pasku kontekstu | Wskazanie modułu, którego dotyczy rozmowa | Domyślny · zmieniony ręcznie | 1 | Widoczny; kliknięcie otwiera listę modułów |
| Znacznik modelu `[Fable 5]` | Lekka etykieta w pasku kontekstu | Wskazanie Wykonawcy odpowiadającego w tej karcie | Domyślny · nadpisany | 1 | Widoczny; kliknięcie otwiera selektor modelu |
| Znacznik kanału `[API]` | Lekka etykieta w pasku kontekstu | Wskazanie kanału połączenia z modelem (API, CLI, SSH, HTTP) | Domyślny · nadpisany · niedostępny | 1 | Widoczny; kliknięcie otwiera selektor kanału |
| Znacznik profilu `[Redaktor]` | Lekka etykieta w pasku kontekstu | Wskazanie aktywnego profilu personalizacji środowiska | Domyślny · własny | 1 | Widoczny; kliknięcie otwiera przełącznik profilu |
| Strumień rozmowy | Obszar zawartości z wypowiedziami Użytkownika i Wykonawcy | Prezentacja przebiegu pracy i wyników | Pusty · strumień aktywny · przewijalny | 1 | Widoczny |
| Referencja zasobu `◈` | Odnośnik do zasobu wytworzonego lub użytego w rozmowie | Przejście do zasobu, dołączenie go, odłożenie na magistralę kontekstu | Domyślna · odłożona do magistrali kontekstu | 1 | Widoczna w strumieniu; kliknięcie otwiera kartę zasobu |
| Przycisk „Zatwierdź” | Kontrolka akceptacji działania oczekującego na decyzję | Zatwierdzenie kroku realizowanego przez Wykonawcę | Domyślny · ładowanie · zatwierdzony | 1 | Widoczny przy kroku oczekującym na decyzję |
| Przycisk „Przerwij” | Kontrolka zatrzymania bieżącego działania | Przerwanie operacji Wykonawcy w toku | Domyślny · ładowanie | 1 | Widoczny podczas pracy Wykonawcy |
| Pole polecenia | Pole `.dn-textarea` | Wydanie polecenia w języku naturalnym | Puste · wypełnione · wysyłanie | 1 | Widoczne |
| Przycisk wyślij `➤` | Ikona wysłania treści pola | Przekazanie polecenia Wykonawcy | Domyślny · ładowanie | 1 | Widoczny; przy pustym polu wyświetla komunikat kontekstowy „Wpisz polecenie przed wysłaniem” |
| Kebab pola polecenia `⋮` | Zbiorcze menu akcji pola | Załącznik, wybór poziomu wysiłku, wybór trybu odpowiedzi, kompresja kontekstu, rozgałęzienie rozmowy | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ |
| Nagłówek „Execution Loop ▼” | Zwinięty wyzwalacz okna pętli wykonawczej | Otwarcie kolumny Execution Loop Window obok Chat Window | Zwinięty · rozwinięty | 2 | Kliknięcie nagłówka, skrót klawiszowy albo polecenie w Chat Window |
| Kebab okna `⋮` | Zbiorcze menu okna | Zakres kontekstu Chat Window, historia i retencja, eksport rozmowy, szerokość kolumny | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ w nagłówku |
| Wskaźnik zapełnienia kontekstu | Pasek z wartością procentową | Sygnał, kiedy skompresować lub rozgałęzić rozmowę | Poniżej progu · przy progu · powyżej progu | 2 | Widoczny w kolumnie stanu; kliknięcie otwiera ustawienia kontekstu okna |

### 4.3. Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca

Execution Loop Window jest elementem pierwszoplanowym środowiska. Prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań i kontrolę realizacji procesów komunikacyjnych środowiska: dekompozycję zlecenia użytkownika na zadania modułowe, kolejkowanie ich, kontrolę jakości wyników i decyzje o ponowieniu. Okno otwiera się jako kolumna sąsiadująca z Chat Window.

Makieta:

```
 ┌─ Execution Loop Window · Koordynator ↔ Wykonawca ──
 │  zlecenie: „Raport z materiałów zebranych w Research” ⋮
 ├────────────────────────────────────────────────────
 │  DEKOMPOZYCJA ZLECENIA
 │    1. Research  · zebranie źródeł        ✔ zakończone
 │    2. Library   · osadzenie w zbiorze    ✔ zakończone
 │    3. Studio    · redakcja raportu       ⟳ w toku
 │    4. Translate · wersja obcojęzyczna    ◌ oczekuje
 ├────────────────────────────────────────────────────
 │  WYMIANA KOMUNIKATÓW STERUJĄCYCH
 │    Koordynator → Wykonawca: zadanie 3, kontekst 1–2
 │    Wykonawca → Koordynator: wynik częściowy, 62%
 ├────────────────────────────────────────────────────
 │  KONTROLA JAKOŚCI
 │    zadanie 2 · wynik przyjęty
 │    zadanie 3 · ocena po zakończeniu
 ├────────────────────────────────────────────────────
 │  PRZEBIEG PĘTLI   ▮▮▮▮▮▯▯▯  3/4
 │  [Wstrzymaj] [Wznów] [Przerwij] [Korekta zlecenia] ⋮
 └────────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nagłówek zlecenia | Belka z treścią bieżącego zlecenia | Identyfikacja procesu prowadzonego przez pętlę | Domyślny · zlecenie skorygowane | 2 | Widoczny po otwarciu okna |
| Lista dekompozycji zlecenia | Uporządkowany wykaz zadań wywiedzionych ze zlecenia | Wgląd w podział pracy między moduły środowiska | Pusta · w toku · zakończona | 2 | Widoczna po otwarciu okna |
| Pozycja zadania | Wiersz z modułem docelowym, opisem i statusem | Śledzenie pojedynczego zadania pętli | Oczekuje · w toku · zakończone · błąd · ponowione | 2 | Widoczna; kliknięcie otwiera szczegóły zadania |
| Kolejka zadań | Zestawienie zadań oczekujących na wykonanie | Ustalenie kolejności realizacji | Pusta · zapełniona · wstrzymana | 2 | Widoczna po otwarciu okna |
| Strumień komunikatów sterujących | Obszar wymiany komunikatów Koordynator ↔ Wykonawca | Wgląd w treść koordynacji i nadzoru | Pusty · strumień aktywny | 2 | Widoczny po otwarciu okna |
| Panel kontroli jakości | Zestawienie ocen wyników zadań | Podstawa decyzji o przyjęciu wyniku albo ponowieniu zadania | Brak ocen · wynik przyjęty · wynik odrzucony | 2 | Widoczny po otwarciu okna |
| Decyzja o ponowieniu | Kontrolka ponownego uruchomienia zadania z korektą | Naprawa wyniku odrzuconego w kontroli jakości | Domyślna · ładowanie | 3 | Kebab pozycji zadania albo panel kontroli jakości |
| Wskaźnik przebiegu pętli | Pasek postępu z licznikiem zadań | Orientacja w zaawansowaniu całego zlecenia | 0–100% · nieokreślony | 2 | Widoczny po otwarciu okna |
| Przycisk „Wstrzymaj” | Kontrolka zatrzymania pętli | Wstrzymanie realizacji bez utraty stanu zadań | Domyślny · wstrzymane | 2 | Widoczny w pasku sterowania przebiegiem |
| Przycisk „Wznów” | Kontrolka podjęcia wstrzymanej pętli | Kontynuacja realizacji od zadania bieżącego | Domyślny · ładowanie | 2 | Widoczny w pasku sterowania przebiegiem |
| Przycisk „Przerwij” | Kontrolka zakończenia pętli | Zamknięcie realizacji zlecenia | Domyślny · ładowanie | 2 | Widoczny w pasku sterowania przebiegiem |
| Przycisk „Korekta zlecenia” | Kontrolka zmiany treści zlecenia w toku | Zmiana zakresu pracy bez zakładania nowego zlecenia | Domyślny · edycja | 2 | Widoczny w pasku sterowania przebiegiem |
| Kebab okna `⋮` | Zbiorcze menu okna | Historia przebiegów, eksport przebiegu, przypisanie Wykonawcy do zadań, szerokość kolumny | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ w nagłówku |

### 4.4. Obszar roboczy modułu

Prawa, dominująca kolumna obszaru roboczego mieści okno operacyjne aktywnego modułu. Anatomia okna jest wspólna dla dziewięciu modułów; wnętrze obszaru zawartości należy do modułu.

| Element okna operacyjnego | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Nagłówek okna operacyjnego | Belka z nazwą okna i modułu | Identyfikacja modułu prowadzonego w kolumnie | Domyślny · zmaksymalizowany | 1 | Widoczny |
| Obszar zawartości modułu | Główny obszar treści modułu | Praca nad treścią właściwa dla modułu | Pusty · aktywny · ładowanie · błąd | 1 | Widoczny |
| Uchwyt szerokości kolumny | Pionowa granica między kolumnami | Regulacja szerokości kolumny modułowej | Domyślny · przeciąganie | 1 | Przeciągnięcie granicy |
| Akcja „Wyślij do…” | Element zbiorczy przekazania treści do innego modułu | Łańcuch pracy nad treścią jako jawne przekazanie | Domyślna · rozwinięta | 2 | Zaznaczenie treści albo kebab okna |
| Akcja „Do magistrali” | Kontrolka odłożenia zasobu na magistralę kontekstu | Udostępnienie zasobu wszystkim modułom środowiska | Domyślna · odłożono | 2 | Zaznaczenie treści albo kebab okna |
| Maksymalizacja okna | Rozszerzenie kolumny modułowej na cały obszar roboczy | Powiększenie okna bez zmiany presetu układu | Domyślna · zmaksymalizowane | 2 | Dwuklik nagłówka okna albo skrót klawiszowy |
| Kebab okna `⋮` | Zbiorcze menu okna operacyjnego | Akcje modułu, eksport zawartości, tryb odczytu, reset szerokości kolumny | Domyślny · rozwinięty | 3 | Kliknięcie ikony ⋮ w nagłówku |
| Tryb skupienia | Ukrycie bocznej nawigacji, paska kart i kolumny stanu | Praca nad długą treścią bez rozproszeń | Wyłączony · włączony | 3 | Skrót klawiszowy albo pozycja w menu środowiska |

### 4.5. Kolumna stanu środowiska

Kolumna stanu jest wąską kolumną na prawej krawędzi obszaru roboczego, o pełnej wysokości. Agreguje stan warstwy sesji i środowiska. Każda pozycja jest klikalna i prowadzi do właściwego ustawienia albo okna pomocniczego.

```
 ┌─ Kolumna stanu ─┐
 │  ● API          │  kanał i model bieżącej karty
 │  ▮▮▮▯ 62%       │  zapełnienie kontekstu rozmowy
 │  ⟳ 2            │  procesy w tle środowiska
 │  ⇄ połączony    │  stan kanału połączenia
 │  ⛉ sesja        │  zasięg izolacji bieżącej karty
 │  ◈ 5            │  elementy magistrali kontekstu
 │  ✉ 3            │  powiadomienia środowiska
 │  ◑ Redaktor     │  profil personalizacji i motyw
 │  Ctrl/Cmd + K   │  skrót palety poleceń
 └─────────────────┘
```

| Pozycja kolumny stanu | Co pokazuje | Zachowanie po kliknięciu | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Kanał i model bieżącej karty | Aktywny kanał (API, CLI, SSH, HTTP) i nazwę modelu przypisanego karcie | Otwiera ustawienia szybkie sesji ze zmianą modelu i kanału | 1 | Widoczna |
| Zapełnienie kontekstu | Pasek i wartość procentową zapełnienia okna kontekstu bieżącej rozmowy | Otwiera ustawienia kontekstu Chat Window: kompresja, rozgałęzienie | 1 | Widoczna |
| Procesy w tle środowiska | Liczbę aktywnych procesów w kartach środowiska | Otwiera listę procesów w tle jako rozszerzenie boczne | 1 | Widoczna |
| Stan połączenia | Status kanału WebSocket z serwerem | Otwiera szczegóły połączenia z ponowieniem | 1 | Widoczna |
| Zasięg izolacji bieżącej karty | Poziom zasięgu obowiązujący kartę: globalny, środowisko, projekt, sesja | Otwiera okno punktów izolacji | 1 | Widoczna |
| Licznik magistrali kontekstu | Liczbę zasobów odłożonych na środowiskowej magistrali kontekstu | Otwiera magistralę kontekstu jako rozszerzenie boczne | 1 | Widoczna |
| Plakietka powiadomień | Liczbę zdarzeń nowych z modułów środowiska | Otwiera centrum powiadomień jako rozszerzenie boczne (`interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 11.6) | 1 | Widoczna |
| Profil personalizacji i motyw | Nazwę aktywnego profilu personalizacji oraz motyw jasny lub ciemny | Otwiera przełącznik profilu i motywu | 1 | Widoczna |
| Skrót palety poleceń | Podpowiedź skrótu otwierającego paletę poleceń | Otwiera paletę poleceń | 1 | Widoczna |

Zestaw i kolejność pozycji kolumny stanu podlegają ustawieniu konfiguracyjnemu w warstwie środowiska i sesji (rozdz. 7).

### 4.6. Pasek kart sesji

| Element paska kart | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Karta sesji | Zakładka reprezentująca jedną przestrzeń pracy nad treścią | Przełączanie między równolegle prowadzonymi przestrzeniami | Domyślna · aktywna · przypięta · z aktywnością w tle · zamykana | 1 | Widoczna; kliknięcie przełącza kartę |
| Kropka stanu karty | Wskaźnik pracy w tle z typem procesu | Wiedza, która karta pracuje, gdy widoczna jest inna | Bezczynna · badanie · przekład · nasłuch głosu · model w Roundtable · błąd | 1 | Widoczna przy karcie z aktywnym procesem |
| Przycisk „+” | Ikona dodania karty | Otwarcie nowej przestrzeni pracy | Domyślny · najechanie | 1 | Widoczny przy ostatniej karcie |
| Nawigacja wstecz/dalej | Para przycisków historii karty | Cofnięcie i powtórzenie nawigacji po modułach i stanach w obrębie karty | Domyślna · nieaktywna (pusty stos) | 1 | Widoczna |
| Selektor presetu układu | Element zwinięty `Układ ▼` | Uruchomienie zapisanego układu kolumn jednym kliknięciem | Domyślny · rozwinięty · preset aktywny | 2 | Kliknięcie elementu `Układ ▼` |
| Etykieta grupy kart | Kolorowana etykieta łącząca karty w nazwaną grupę | Porządek przy wielu kartach jednego tematu | Domyślna · zwinięta grupa | 2 | Widoczna przy pierwszej karcie grupy |
| Menu karty | Lista akcji karty: zmiana nazwy, przypnij, duplikuj, dodaj do grupy, przekaż na Mobile, zamknij | Zarządzanie pojedynczą kartą | Domyślne · rozwinięte | 3 | Prawy klawisz myszy na karcie albo kebab (⋮) karty |
| Przywracanie zamkniętej karty | Stos zamkniętych kart bieżącej sesji połączenia | Odzyskanie zamkniętej pracy | Dostępne · pusty stos | 3 | Skrót klawiszowy albo pozycja w menu środowiska |

### 4.7. Pasek górny

| Element paska górnego | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Identyfikacja platformy | Znak i nazwa Danaco Console | Powrót na stronę główną | Domyślna | 1 | Widoczna |
| Pole wyszukiwania | Pole `.dn-pasek-szukaj` ze znacznikiem zakresu | Znalezienie treści środowiska z jednego pola: karty, moduły, zbiory Library, wyniki Research, historia rozmów | Puste · z treścią · brak wyników | 1 | Widoczne |
| Znacznik zakresu wyszukiwania | Segment przy polu: „TalkIn”, „platforma”, „bieżąca karta” | Pewność, gdzie prowadzone jest wyszukiwanie | Środowisko · platforma · karta | 1 | Widoczny; kliknięcie przełącza zakres |
| Ustawienia szybkie sesji | Rozwijany panel przy ikonie `ustawienia` | Zmiana modelu i kanału, motywu, retencji historii karty, zasięgu izolacji karty, kontekstu wiedzy, profilu personalizacji | Domyślne · rozwinięte | 2 | Kliknięcie ikony `ustawienia` |
| Przełącznik gęstości interfejsu | Kontrolka zmiany zagęszczenia powłoki: komfortowa albo zwarta | Dopasowanie liczby widocznej treści do ekranu | Komfortowa · zwarta | 3 | Pozycja w ustawieniach szybkich sesji |
| Menu środowiska (☰) | Zbiorcze menu `Operacje ▼`: presety układu, zapis sesji roboczej, eksport migawki, reset układu, ściągawka skrótów, personalizacja, okna pomocnicze | Jeden dostęp do rzadszych funkcji powłoki | Domyślne · rozwinięte | 3 | Kliknięcie ikony ☰ |

### 4.8. Magistrala kontekstu

Mechanizm magistrali kontekstu — przeznaczenie, sposób przekazania artefaktu, zasięg, warstwę widoczności i sposób wywołania oraz relacje do modułu Library, punktów izolacji i encji artefaktu — opisuje `interfejs-uzytkownika/przeplyw-okien.md` (rozdz. 6a). Specyfika środowiska TalkIn: wyzwalaczem magistrali jest licznik `◈` w kolumnie stanu, a zasięg domyślny obejmuje wszystkie moduły środowiska, zgodnie z konwersacyjnym charakterem pracy — zasób odłożony w jednym module jest dostępny w pozostałych bez przełączania karty.

| Element magistrali kontekstu | Co to jest | Do czego służy | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Lista zasobów | Zestawienie odłożonych referencji z modułem pochodzenia | Przegląd wspólnego kontekstu międzymodułowego | Pusta · zapełniona · filtrowana | 2 | Widoczna po otwarciu magistrali kontekstu |
| Pozycja zasobu | Wiersz z typem zasobu, nazwą i modułem pochodzenia | Praca z pojedynczym zasobem | Domyślna · zaznaczona · użyta w bieżącej karcie | 2 | Widoczna; kliknięcie otwiera kartę zasobu |
| Akcja „Otwórz w…” | Element zbiorczy wyboru modułu docelowego | Otwarcie zasobu w wybranym module środowiska | Domyślna · rozwinięta | 2 | Kliknięcie elementu `Otwórz w ▼` |
| Akcja „Dołącz do rozmowy” | Kontrolka dołączenia zasobu jako kontekstu polecenia | Przekazanie zasobu Wykonawcy w Chat Window | Domyślna · dołączono | 2 | Kliknięcie akcji na pozycji zasobu |
| Przeciągnięcie zasobu | Przeciągnięcie pozycji do kolumny modułu | Przeniesienie treści gestem, bez schowka systemowego | Domyślne · przeciąganie · strefa zrzutu podświetlona | 2 | Przeciągnięcie pozycji |
| Retencja magistrali kontekstu | Ustawienie okresu przechowywania referencji | Kontrola trwałości wspólnego kontekstu | Sesja · projekt · środowisko | 4 | Okno konfiguracji |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — polecenie w Chat Window

```
 Użytkownik ── polecenie w języku naturalnym ──► Chat Window
                                                     │
                                                     ▼
                                                Koordynator
                                     dekompozycja zlecenia na zadania
                                                     │
                        ┌────────────────────────────┼───────────────────────┐
                        ▼                            ▼                       ▼
                   zadanie 1                    zadanie 2               zadanie n
                   moduł A                      moduł B                 moduł C
                        │                            │                       │
                        └────────────► Wykonawca ◄───┴───────────────────────┘
                                          │
                             wynik zadania │ komunikat sterujący
                                          ▼
                              Execution Loop Window
                              kontrola jakości wyniku
                                          │
                            ┌─────────────┴─────────────┐
                     wynik przyjęty              wynik odrzucony
                            │                           │
                            ▼                           ▼
                  aktualizacja okna            ponowienie zadania
                  operacyjnego modułu          z korektą zakresu
                            │
                            ▼
                       Chat Window
             prezentacja wyniku i referencji zasobu
                            │
                            ▼
                  Użytkownik zatwierdza
```

### 5.2. Przepływ treści między modułami

```
  Research ── wynik badania ──► Magistrala kontekstu ──► Studio
      │                              ▲                 │
      │                              │                 │
      └──── „Wyślij do…” ────────────┘                 │
                                                       ▼
  Library ◄──── osadzenie dokumentu ──────────── Studio Editor
      │                                                │
      │                                                ▼
      └──── kontekst wiedzy środowiska ──►  Translate ──► przekład
                                                       │
                                                       ▼
                                            Roundtable ──► ocena wielomodelowa
                                                       │
                                                       ▼
                                          referencje krzyżowe
                            „ten raport powstał z tych źródeł i tego dokumentu”
```

Trzy mechanizmy przenoszenia treści działają równolegle:

| Mechanizm | Kierunek | Efekt |
|---|---|---|
| Magistrala kontekstu | Dowolny moduł → wszystkie moduły | Referencja dostępna we wszystkich modułach środowiska |
| „Wyślij do…” | Moduł źródłowy → wskazany moduł docelowy | Jednokrotne przekazanie treści jako wejścia modułu docelowego |
| Przeciągnięcie treści | Kolumna źródłowa → kolumna docelowa | Osadzenie treści w miejscu zrzutu gestem |

### 5.3. Przepływ pracy w karcie sesji

```
  Nowa karta ──► wybór modułu ──► układ kolumn ──► praca nad treścią
       │              │                │                   │
       │              │                │                   ▼
       │              │                │            zapis artefaktu
       │              │                │                   │
       │              │                ▼                   ▼
       │              │        preset układu       oś czasu treści
       │              │        zapisany                    │
       │              ▼                                    ▼
       │       zapamiętany układ per moduł        referencje krzyżowe
       ▼
  zapis sesji roboczej ──► wznowienie sesji w innym dniu
```

### 5.4. Przepływ nawigacji między modułami

| Krok | Działanie użytkownika | Zachowanie środowiska |
|---|---|---|
| 1 | Kliknięcie pozycji modułu w bocznej nawigacji | Aktywna kolumna obszaru roboczego przyjmuje okno operacyjne modułu; Chat Window pozostaje bez zmian w lewej kolumnie |
| 2 | Powrót do modułu odwiedzonego wcześniej | Środowisko odtwarza zapamiętany układ kolumn tego modułu w tej karcie |
| 3 | „Otwórz obok” z menu kontekstowego pozycji | Powstaje dodatkowa kolumna modułowa po prawej stronie aktywnej |
| 4 | „Otwórz w nowej karcie” | Powstaje karta sesji z tym modułem i presetem układu domyślnym dla karty |
| 5 | Przełączenie środowiska nagłówkiem „TalkIn ▼” | Karty środowiska trwają w tle wraz ze stanem procesów sesji po stronie serwera |

### 5.5. Przepływ nadzoru nad pętlą wykonawczą

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Chat Window | Użytkownik wydaje zlecenie w języku naturalnym |
| 2 | Execution Loop Window | Koordynator prezentuje dekompozycję zlecenia i kolejkę zadań |
| 3 | Execution Loop Window | Wykonawca podejmuje zadanie z kolejki; strumień komunikatów sterujących pokazuje przebieg |
| 4 | Okno operacyjne modułu | Wynik zadania osadza się w module docelowym |
| 5 | Execution Loop Window | Panel kontroli jakości ocenia wynik; wynik odrzucony wraca do kolejki jako zadanie ponowione |
| 6 | Chat Window | Wykonawca przedstawia wynik końcowy z referencjami zasobów; użytkownik zatwierdza albo poleca korektę |
| 7 | Execution Loop Window | Użytkownik steruje przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia |

---

## 6. Stany i powiązania

### 6.1. Stany karty sesji

| Stan | Wyzwalacz | Prezentacja | Wyjście ze stanu |
|---|---|---|---|
| Bezczynna | Brak procesów w tle karty | Karta bez kropki stanu | Uruchomienie procesu |
| Pracująca | Aktywny proces modułu karty | Kropka stanu z typem procesu, licznik `⟳` w kolumnie stanu | Zakończenie procesu |
| Oczekująca na decyzję | Krok procesu oczekuje na zatwierdzenie | Kropka stanu w kolorze uwagi, powiadomienie w centrum powiadomień | Zatwierdzenie albo przerwanie w Chat Window |
| Błąd | Niedostępny kanał modelu albo błąd zadania | Kropka stanu w kolorze błędu, wpis w strumieniu Chat Window | Ponowienie, zmiana kanału albo korekta zlecenia |
| Przypięta | Przypięcie z menu karty | Karta zwężona do ikony na początku paska | Odpięcie z menu karty |
| Zamknięta | Zamknięcie karty | Karta trafia na stos zamkniętych kart | Przywrócenie zamkniętej karty |

### 6.2. Stany połączenia

| Stan | Prezentacja w kolumnie stanu | Zachowanie środowiska |
|---|---|---|
| Połączony | `⇄ połączony` | Strumienie okien aktualizują się na żywo kanałem WebSocket |
| Wznawianie | `⇄ wznawianie` | Interfejs pozostaje sterowalny; polecenia kolejkowane do czasu wznowienia |
| Rozłączony | `⇄ offline` | Stan procesów sesji trwa po stronie serwera; okna odtwarzają się po wznowieniu połączenia |

### 6.3. Stany pętli wykonawczej

| Stan pętli | Prezentacja | Przejścia |
|---|---|---|
| Bezczynna | Pusta lista dekompozycji | Wydanie zlecenia w Chat Window → W toku |
| W toku | Wskaźnik przebiegu, aktywny strumień komunikatów | Wstrzymaj → Wstrzymana; Przerwij → Przerwana; ostatnie zadanie zakończone → Zakończona |
| Wstrzymana | Wskaźnik przebiegu zatrzymany, kolejka wstrzymana | Wznów → W toku; Przerwij → Przerwana |
| Ponowienie | Zadanie odrzucone w kontroli jakości wraca do kolejki | Podjęcie zadania → W toku |
| Przerwana | Pętla zamknięta, stan zadań zachowany | Nowe zlecenie → W toku |
| Zakończona | Wskaźnik przebiegu pełny, wynik w Chat Window | Nowe zlecenie → W toku |

### 6.4. Powiązania między oknami

| Okno źródłowe | Okno docelowe | Charakter powiązania |
|---|---|---|
| Chat Window | Execution Loop Window | Zlecenie wydane w Chat Window uruchamia pętlę wykonawczą i zapełnia listę dekompozycji |
| Execution Loop Window | Okno operacyjne modułu | Zadanie pętli osadza wynik w module docelowym |
| Execution Loop Window | Chat Window | Wynik zadania i decyzja kontroli jakości trafiają do strumienia rozmowy jako wpis z referencją zasobu |
| Okno operacyjne modułu | Magistrala kontekstu | Akcja „Do magistrali” dodaje referencję zasobu dostępną wszystkim modułom |
| Magistrala kontekstu | Chat Window | Akcja „Dołącz do rozmowy” przekazuje zasób jako kontekst polecenia |
| Magistrala kontekstu | Okno operacyjne modułu | Akcja „Otwórz w…” i przeciągnięcie osadzają zasób w module docelowym |
| Okno operacyjne modułu | Okno operacyjne modułu | Akcja „Wyślij do…” przekazuje treść jako wejście modułu docelowego |
| Wszystkie okna | Oś czasu treści | Każdy wytworzony artefakt zapisuje wpis z modułem pochodzenia i czasem |
| Wszystkie okna | Graf referencji krzyżowych | Powstanie artefaktu z innych zasobów utrwala powiązanie proweniencji |
| Wszystkie okna | Kolumna stanu | Stan kanału, kontekstu, procesów i izolacji agreguje się w kolumnie stanu |
| Wszystkie okna | Centrum powiadomień | Zdarzenie wymagające uwagi tworzy powiadomienie środowiska |

### 6.5. Powiązania kontekstu między modułami

| Poziom zasięgu | Co jest współdzielone | Stan wyjściowy |
|---|---|---|
| Sesja | Kontekst Chat Window bieżącej karty | Odrębny dla każdej karty |
| Para modułów | Historia i pamięć Chat Window między dwoma wskazanymi modułami karty | Odrębna; współdzielenie włącza ustawienie konfiguracyjne |
| Projekt | Magistrala kontekstu, oś czasu treści, referencje krzyżowe | Wspólne dla kart projektu |
| Środowisko | Kontekst wiedzy środowiska — zbiory z Library aktywne jako tło wiedzy wszystkich modułów | Wspólny dla środowiska |
| Globalny | Profil personalizacji, motyw, mapa skrótów | Wspólne dla platformy |

### 6.6. Stany elementów wspólnych

| Element | Stany |
|---|---|
| Pozycja modułu w bocznej nawigacji | Domyślna · aktywna · z plakietką aktywności · przypięta · ukryta |
| Kolumna modułowa | Domyślna · zmaksymalizowana · w trybie skupienia · w trakcie zmiany szerokości |
| Chat Window | Bezczynne · Wykonawca pracuje · oczekiwanie na zatwierdzenie · kontekst przy progu zapełnienia · błąd kanału |
| Execution Loop Window | Zwinięte · rozwinięte, pętla bezczynna · rozwinięte, pętla w toku · rozwinięte, pętla wstrzymana |
| Magistrala kontekstu | Zamknięta · otwarta pusta · otwarta zapełniona · otwarta filtrowana |
| Pozycja kolumny stanu | Wartość aktualna · wartość przy progu · wartość w stanie błędu |

---

## 7. Punkty sterowania z okna konfiguracji

Model konfiguracji platformy obejmuje warstwy: **globalny → środowisko → projekt → sesja**; izolacja ma osobne okno o siedmiu poziomach zasięgu. Brak jawnego ustawienia oznacza wartość domyślną i zachowanie kanoniczne. Wszystkie poniższe zakresy są sterowane przez Operatora bez blokad w interfejsie.

### 7.1. Zakresy sterowania

| Zakres | Co Operator ustawia | Warstwa |
|---|---|---|
| Boczna nawigacja | Grupowanie modułów w sekcje, kolejność pozycji, przypięcia, ukrycia, tryb kompaktowy, obecność sekcji „Ostatnie” | Środowisko |
| Układ obszaru roboczego | Kanoniczny układ kolumn nowej karty, liczba kolumn modułowych, szerokość Chat Window, szerokość kolumny Execution Loop Window | Środowisko / projekt |
| Presety układu | Definicje i nazwy presetów kolumn, preset domyślny nowej karty | Środowisko / projekt |
| Chat Window | Zasięg kontekstu (odrębny albo para modułów), retencja historii, próg ostrzeżenia o zapełnieniu kontekstu, zachowanie kompresji, domyślny poziom wysiłku Wykonawcy | Sesja / projekt |
| Execution Loop Window | Domyślna widoczność okna, zakres prezentowanego strumienia komunikatów sterujących, próg kontroli jakości uruchamiający ponowienie, limit ponowień zadania | Środowisko / projekt |
| Karty i sesje | Zachowanie przycisku „+”, nazewnictwo automatyczne albo ręczne, przywracanie kart przy powrocie, próg ostrzeżenia przy zamykaniu wielu kart | Środowisko / sesja |
| Kolumna stanu | Zestaw widocznych pozycji i ich kolejność, obecność całej kolumny | Środowisko / sesja |
| Skróty klawiszowe | Mapa skrótów powłoki, rozwiązywanie konfliktów, zestaw domyślny albo własny | Globalny / środowisko |
| Współdzielenie kontekstu | Zbiory Library aktywne jako kontekst wiedzy środowiska, retencja magistrali kontekstu, jawność referencji krzyżowych, zakres osi czasu treści | Projekt / sesja / para modułów |
| Funkcje wspólne | Obecność notatnika środowiska, centrum powiadomień, kanałów Mobile i Always On Display, domyślny zakres wyszukiwania | Środowisko |
| Motyw i gęstość | Motyw jasny albo ciemny, gęstość interfejsu, profil personalizacji jako nazwany zestaw ustawień | Globalny / środowisko / sesja |
| Izolacja powłoki | Poziom zasięgu karty i relacji między modułami — co współdzielone, co odrębne; stan wyjściowy „pełny dostęp” | Okno punktów izolacji (siedem poziomów) |

### 7.2. Profile personalizacji środowiska

Zestaw ustawień z tabeli 7.1 zapisuje się jako nazwany profil personalizacji środowiska — na przykład „Redaktor”, „Badacz”, „Tłumacz”. Profil przełącza się jednym kliknięciem znacznika `◑` w kolumnie stanu. Zawartość profilu jest jawna i podlega podglądowi przed przełączeniem.

| Profil | Preset układu kolumn | Moduły przypięte | Kontekst wiedzy |
|---|---|---|---|
| Redaktor | Chat Window · Studio Editor · Library Explorer | Studio, Library | Zbiory redakcyjne projektu |
| Badacz | Chat Window · Research Workspace · Browser View · Library Explorer | Research, Browser, Library | Zbiory źródłowe projektu |
| Tłumacz | Chat Window · Studio Editor · Translate Workbench | Studio, Translate | Zbiory terminologiczne projektu |
| Panel ekspertów | Chat Window · Execution Loop Window · Roundtable Board | Roundtable, Agents | Zbiory tematyczne projektu |

### 7.3. Ustawienia szybkie sesji

Ustawienia szybkie sesji w pasku górnym obejmują podzbiór zakresów z tabeli 7.1, sterowany bez otwierania pełnego okna konfiguracji: model i kanał, motyw, retencja historii karty, zasięg izolacji karty, aktywny kontekst wiedzy, profil personalizacji, gęstość interfejsu. Panel zawiera odnośnik do pełnego okna konfiguracji.

---

## 8. Scenariusze użycia

### 8.1. Redakcja raportu ze zgromadzonych źródeł

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Pasek kart | Użytkownik otwiera kartę i wybiera profil personalizacji „Badacz”; układ kolumn przyjmuje Chat Window, Research Workspace, Browser View, Library Explorer |
| 2 | Chat Window | Polecenie: „Zbierz źródła o rynku energii rozproszonej z ostatnich dwóch lat i osadź je w zbiorze projektu” |
| 3 | Execution Loop Window | Koordynator dzieli zlecenie na zadania: zapytanie badawcze w Research, weryfikacja stron w Browser, osadzenie w Library |
| 4 | Research Workspace | Wynik badania pojawia się jako lista źródeł z syntezą; użytkownik odkłada wybrane pozycje na magistralę kontekstu |
| 5 | Boczna nawigacja | Przełączenie na profil „Redaktor”; kolumny przyjmują Studio Editor i Library Explorer, Chat Window zachowuje kontekst rozmowy |
| 6 | Magistrala kontekstu | Przeciągnięcie źródeł do Studio Editor osadza je jako materiał wyjściowy |
| 7 | Chat Window | Polecenie redakcji raportu; Wykonawca prowadzi pracę w Studio Editor, prezentując wyniki w strumieniu rozmowy |
| 8 | Execution Loop Window | Kontrola jakości ocenia kompletność cytowań; zadanie z brakami wraca do kolejki jako ponowione |
| 9 | Chat Window | Użytkownik zatwierdza wynik końcowy; raport zapisuje się do Library z referencjami krzyżowymi do źródeł |

### 8.2. Przekład dokumentu z kontrolą terminologii

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Kolumna stanu | Przełączenie profilu na „Tłumacz”; kolumny przyjmują Studio Editor i Translate Workbench |
| 2 | Chat Window | Polecenie: „Przełóż ten dokument, trzymając się zbioru terminologicznego projektu” |
| 3 | Execution Loop Window | Dekompozycja: pobranie zbioru terminów z Library, przekład sekcjami, kontrola spójności terminologicznej |
| 4 | Translate Workbench | Widok dwujęzyczny wypełnia się sekcjami w miarę postępu pętli |
| 5 | Execution Loop Window | Kontrola jakości wykrywa niespójność terminu; zadanie sekcji wraca do kolejki z korektą zakresu |
| 6 | Chat Window | Wykonawca przedstawia listę rozstrzygnięć terminologicznych do zatwierdzenia |
| 7 | Oś czasu treści | Wpis rejestruje powstanie przekładu z odesłaniem do dokumentu źródłowego i zbioru terminów |

### 8.3. Debata wielomodelowa nad wariantem treści

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Boczna nawigacja | Otwarcie modułu Roundtable obok Studio poleceniem „Otwórz obok” |
| 2 | Studio Editor | Zaznaczenie fragmentu i akcja „Wyślij do…” z modułem docelowym Roundtable |
| 3 | Roundtable Board | Modele przedstawiają stanowiska wobec fragmentu; panel konsensusu porządkuje wnioski |
| 4 | Execution Loop Window | Koordynator nadzoruje kolejność wypowiedzi i zamyka rundę po osiągnięciu warunku zakończenia |
| 5 | Chat Window | Wykonawca streszcza wynik debaty i proponuje wariant redakcyjny |
| 6 | Studio Editor | Zatwierdzony wariant zastępuje fragment; graf referencji krzyżowych utrwala pochodzenie zmiany |

### 8.4. Wznowienie złożonej pracy po przerwie

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Menu środowiska (☰) | Wczytanie zapisanej sesji roboczej |
| 2 | Pasek kart | Karty odtwarzają się wraz z modułami, układem kolumn i kontekstem każdej z nich |
| 3 | Kolumna stanu | Licznik procesów w tle pokazuje zadania kontynuowane po stronie serwera |
| 4 | Execution Loop Window | Pętla wstrzymana przed przerwą podejmuje pracę po przycisku „Wznów” |
| 5 | Oś czasu treści | Przegląd chronologii pracy przywraca orientację w tym, co powstało przed przerwą |

### 8.5. Nadzór nad pracą poza stanowiskiem

| Krok | Miejsce | Działanie |
|---|---|---|
| 1 | Menu karty | Akcja „Przekaż na Mobile” oznacza kartę do monitorowania z urządzenia mobilnego |
| 2 | Centrum powiadomień | Zdarzenia karty tworzą powiadomienia przekazywane kanałem Mobile |
| 3 | Mobile | Użytkownik zatwierdza krok oczekujący na decyzję albo wstrzymuje pętlę wykonawczą |
| 4 | Always On Display | Agent towarzyszący przedstawia podsumowanie stanu procesu i sugeruje kolejny krok |
| 5 | Execution Loop Window | Po powrocie do stanowiska pełny przebieg pętli jest dostępny w historii przebiegów |

---

## Załącznik A. Skróty klawiszowe środowiska

Mapa skrótów jest w całości konfigurowalna z edytora mapy skrótów (rozdz. 7.1). Konflikt przypisań sygnalizowany jest ostrzeżeniem, nie blokadą.

### A.1. Sterowanie powłoką

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + K` | Paleta poleceń środowiska |
| `?` | Ściągawka skrótów |
| `Ctrl/Cmd + Shift + U` | Przełącznik profilu personalizacji |
| `Ctrl/Cmd + ,` | Okno konfiguracji |
| `Ctrl/Cmd + Shift + M` | Menedżer kart |

### A.2. Moduły i nawigacja

| Skrót | Działanie |
|---|---|
| `Alt + 1` … `Alt + 9` | Przejście do modułu według kolejności bocznej nawigacji |
| `Alt + \`` | Powrót do ostatnio używanego modułu |
| `Ctrl/Cmd + B` | Tryb kompaktowy bocznej nawigacji |
| `Ctrl/Cmd + E` | Pole filtrowania modułów |
| `Alt + ←` / `Alt + →` | Nawigacja wstecz i dalej w obrębie karty |

### A.3. Karty sesji

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + T` | Nowa karta sesji |
| `Ctrl/Cmd + W` | Zamknięcie karty |
| `Ctrl/Cmd + Shift + T` | Przywrócenie zamkniętej karty |
| `Ctrl + Tab` | Następna karta |
| `Ctrl + Shift + Tab` | Poprzednia karta |
| `Ctrl/Cmd + D` | Duplikowanie karty |

### A.4. Układ kolumn

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + \` | Dodanie kolumny modułowej po prawej stronie aktywnej |
| `Ctrl/Cmd + Shift + \` | Zamknięcie aktywnej kolumny modułowej |
| `Ctrl/Cmd + Shift + F` | Tryb skupienia |
| `Ctrl/Cmd + Shift + Enter` | Maksymalizacja aktywnego okna operacyjnego |
| `Ctrl/Cmd + Alt + R` | Przywrócenie kanonicznego układu kolumn |
| `Ctrl/Cmd + Alt + ←` / `→` | Zmiana szerokości aktywnej kolumny |

### A.5. Okna komunikacji operacyjnej

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + L` | Fokus na pole polecenia Chat Window |
| `Ctrl/Cmd + Shift + L` | Powrót fokusu do okna operacyjnego modułu |
| `Ctrl/Cmd + J` | Execution Loop Window |
| `Ctrl/Cmd + Enter` | Zatwierdzenie kroku oczekującego na decyzję |
| `Ctrl/Cmd + .` | Przerwanie bieżącego działania Wykonawcy |
| `Ctrl/Cmd + Shift + K` | Kompresja kontekstu rozmowy |

### A.6. Warstwa spajająca treść

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + Shift + S` | Magistrala kontekstu |
| `Ctrl/Cmd + Shift + N` | Notatnik środowiska |
| `Ctrl/Cmd + Shift + O` | Oś czasu treści |
| `Ctrl/Cmd + Shift + D` | Odłożenie zaznaczonej treści na magistralę kontekstu |
| `Ctrl/Cmd + Shift + G` | Akcja „Wyślij do…” dla zaznaczonej treści |
| `Ctrl/Cmd + F` | Wyszukiwanie w zakresie środowiska |

---

*Koniec dokumentu. Środowisko TalkIn — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
