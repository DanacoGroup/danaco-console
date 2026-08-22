# Danaco Console — Moduł Design

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
| **Tytuł** | Moduł Design |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant (co, gdzie, w jakiej formie) · deweloper (co zbudować) |
| **Przeznaczenie** | Ustala interfejs modułu Design: komplet okien operacyjnych, katalog elementów, przepływy pracy, stany oraz punkty sterowania z okna konfiguracji |
| **Zakres** | okna operacyjne modułu dostępnego w WorkSpace i CodeStudio, katalog elementów interfejsu, przepływy pracy, komendy kontraktu obszaru `design` |
| **Poza zakresem** | system wizualny platformy jako taki — [System wizualny](../interfejs-uzytkownika/system-wizualny.md) |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/design.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `design`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie; klucze jawne; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
   - [1.1 Rola modułu](#11-rola-modułu)
   - [1.2 Dla kogo](#12-dla-kogo)
   - [1.3 Po co](#13-po-co)
   - [1.4 Charakter pracy](#14-charakter-pracy)
   - [1.5 Granica tematyczna](#15-granica-tematyczna)
   - [1.6 Miejsce modułu w architekturze platformy](#16-miejsce-modułu-w-architekturze-platformy)
2. [Komplet okien operacyjnych modułu — przegląd](#2-komplet-okien-operacyjnych-modułu--przegląd)
   - [2.1 Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
3. [Specyfikacja okien operacyjnych — pełny arsenał narzędzi](#3-specyfikacja-okien-operacyjnych--pełny-arsenał-narzędzi)
   - [3.1 Chat Window (kanał Użytkownik ↔ Wykonawca)](#31-chat-window-kanał-użytkownik--wykonawca)
   - [3.2 Execution Loop Window (kanał Koordynator ↔ Wykonawca)](#32-execution-loop-window-kanał-koordynator--wykonawca)
   - [3.3 Design Board](#33-design-board)
   - [3.4 Assets Panel](#34-assets-panel)
   - [3.5 Prompt Builder](#35-prompt-builder)
   - [3.6 Preview Window](#36-preview-window)
   - [3.7 Tokens & System Panel](#37-tokens--system-panel)
4. [Katalog funkcji i narzędzi](#4-katalog-funkcji-i-narzędzi)
   - [4.1 Grupa A — Generowanie grafiki przez AI](#41-grupa-a--generowanie-grafiki-przez-ai)
   - [4.2 Grupa B — Edycja rastrowa](#42-grupa-b--edycja-rastrowa)
   - [4.3 Grupa C — Grafika i edycja wektorowa](#43-grupa-c--grafika-i-edycja-wektorowa)
   - [4.4 Grupa D — Projektowanie UI i makiety](#44-grupa-d--projektowanie-ui-i-makiety)
   - [4.5 Grupa E — Tokeny projektowe i system projektowy](#45-grupa-e--tokeny-projektowe-i-system-projektowy)
   - [4.6 Grupa F — Kolor](#46-grupa-f--kolor)
   - [4.7 Grupa G — Ikony i typografia](#47-grupa-g--ikony-i-typografia)
   - [4.8 Grupa H — Eksport, podgląd, handoff](#48-grupa-h--eksport-podgląd-handoff)
   - [4.9 Grupa I — Marketing, szablony, zasoby zewnętrzne](#49-grupa-i--marketing-szablony-zasoby-zewnętrzne)
   - [4.10 Grupa J — Współpraca, wersjonowanie, organizacja](#410-grupa-j--współpraca-wersjonowanie-organizacja)
5. [Katalog elementów interfejsu](#5-katalog-elementów-interfejsu)
6. [Przepływy pracy w module](#6-przepływy-pracy-w-module)
   - [6.1 Od polecenia do zaakceptowanego zasobu](#61-od-polecenia-do-zaakceptowanego-zasobu)
   - [6.2 Zestawianie kompozycji z wielu zasobów](#62-zestawianie-kompozycji-z-wielu-zasobów)
   - [6.3 Budowa i wydanie systemu projektowego](#63-budowa-i-wydanie-systemu-projektowego)
   - [6.4 Przekazanie zasobu do modułu docelowego](#64-przekazanie-zasobu-do-modułu-docelowego)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
8. [Stany, dane i powiązania z innymi modułami](#8-stany-dane-i-powiązania-z-innymi-modułami)
   - [8.1 Dane wykorzystywane przez AI w module](#81-dane-wykorzystywane-przez-ai-w-module)
   - [8.2 Stany zasobu wizualnego](#82-stany-zasobu-wizualnego)
   - [8.3 Powiązania konfigurowalne z innymi modułami](#83-powiązania-konfigurowalne-z-innymi-modułami)
9. [Scenariusze użycia](#9-scenariusze-użycia)
   - [9.1 Budowa spójnego zestawu ilustracji brandingowych](#91-budowa-spójnego-zestawu-ilustracji-brandingowych)
   - [9.2 Projektowanie interfejsu produktu równolegle z modułem Apps](#92-projektowanie-interfejsu-produktu-równolegle-z-modułem-apps)
   - [9.3 Iteracyjna korekta ilustracji na podstawie uwag](#93-iteracyjna-korekta-ilustracji-na-podstawie-uwag)
10. [Zgodność z rdzeniem platformy](#10-zgodność-z-rdzeniem-platformy)
11. [Kryteria odbioru](#11-kryteria-odbioru)
12. [Załącznik — pełny wykaz komend kontraktu modułu Design](#załącznik--pełny-wykaz-komend-kontraktu-modułu-design)
   - [Obszar `design` — 106 komend](#obszar-design--106-komend)

---

## 1. Przeznaczenie i kontekst

### 1.1 Rola modułu

Design jest **kompletnym warsztatem wizualnym platformy** — pojedynczym miejscem, w którym powstaje, jest edytowany, porządkowany i wydawany każdy zasób graficzny wykorzystywany w pozostałych modułach. Moduł pokrywa dziesięć obszarów tematycznych, rozwiniętych w katalog funkcji rozdz. 4
jako grupy A–J:

| Obszar | Przedmiot pracy |
|---|---|
| A. Generowanie AI | Tworzenie grafiki z opisu, wariacje, przekształcenia neuronowe |
| B. Edycja rastrowa | Retusz, kadrowanie, warstwy, maski, korekcja barwna, filtry |
| C. Grafika wektorowa | Rysunek ścieżkowy, kształty, operacje logiczne, wektoryzacja |
| D. Projektowanie UI i makiety | Ramki, komponenty, auto-układ, prototypowanie, makiety szkicowe |
| E. Tokeny i system projektowy | Kolory, typografia, odstępy, motywy, eksport tokenów do kodu |
| F. Kolor | Palety, kontrast, harmonie, symulacja wad wzroku, ekstrakcja |
| G. Ikony i typografia | Zestawy ikon, favicony, sprite SVG, dobór par krojów, podgląd krojów sieciowych |
| H. Eksport, podgląd, handoff | Formaty wyjściowe, kompresja, skalowanie @1/2/3x, inspekcja, kod |
| I. Marketing, szablony, zasoby zewnętrzne | Szablony formatów społecznościowych, zestawy rozmiarów kampanii, biblioteka szablonów brandingowych, import zasobów zewnętrznych, treści próbne, branding wsadowy |
| J. Współpraca, wersjonowanie, organizacja | Warstwy i drzewo obiektów, wersjonowanie kompozycji, porównanie wizualne, komentarze, kursory i obecność, etykiety i wyszukiwanie, metadane pochodzenia |

Objaśnienie terminu: `handoff` — przekazanie zatwierdzonego zasobu wizualnego wraz z jego
specyfikacją (wymiary, żetony, warianty eksportu) osobie albo modułowi realizującemu
wdrożenie; termin zapisywany w tym opracowaniu w postaci nieodmiennej.

Moduł obejmuje ponadto **warstwy przekrojowe**: bibliotekę zasobów (Assets Panel), wersjonowanie i współpracę (kursory, komentarze), animację lekką (Lottie, GIF, APNG) oraz import zasobów zewnętrznych.

### 1.2 Dla kogo

Design jest przestrzenią dla użytkowników **generujących i edytujących materiały graficzne** — od pojedynczych ilustracji i elementów brandingowych po kompletne makiety interfejsu. Adresatem są projektanci, zespoły marketingowe oraz zespoły produktowe budujące jednocześnie w module Apps ([Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) rozdz. 4.14) warstwę kliencką wymagającą zasobów wizualnych.

### 1.3 Po co

Moduł łączy w jednym miejscu **generowanie grafiki przez AI, jej edycję i zestawianie w spójną kompozycję** — z bezpośrednim wsparciem AI na każdym etapie procesu twórczego, zamiast przełączania się między osobnym narzędziem generującym a osobnym narzędziem porządkującym wynik.

### 1.4 Charakter pracy

| Cecha | Opis |
|---|---|
| Punkt wyjścia | Zamiar wizualny wyrażony poleceniem, a nie gotowy plik do edycji |
| Rytm pracy | Cykl: sformułowanie polecenia → wygenerowanie → ocena → zestawienie lub odrzucenie → iteracja |
| Wynik | Zasób wizualny gromadzony w Assets Panel, gotowy do dalszego wykorzystania w innych modułach |
| Dostępność środowiskowa | WorkSpace (praca projektowa i produktowa) oraz CodeStudio (budowa interfejsu aplikacji) — moduł nieobecny w TalkIn, gdzie centralnym przedmiotem pracy jest treść, nie grafika |

### 1.5 Granica tematyczna

Granica tematyczna jest jawna i utrzymuje czytelny podział między modułami platformy:

| Poza zakresem Design | Uzasadnienie | Moduł właściwy |
|---|---|---|
| Redagowanie treści dokumentów (tekst długi, PDF/DOCX, korekta, streszczenia) | Przedmiotem Design jest grafika, nie proza | **Studio** ([Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) rozdz. 4.1) |
| Budowa działającej warstwy klienckiej produktu (kod komponentów, routing, stan) | Design wytwarza zasoby i makiety, nie kod aplikacji | **Apps — Frontend Workspace** ([Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) rozdz. 4.14) |
| Trwałe repozytorium plików całej organizacji | Assets Panel jest podręczną biblioteką modułu, nie centralnym archiwum | **Library** ([Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) rozdz. 4.6) |
| Montaż i renderowanie wideo (timeline, ścieżki, klatki kluczowe wideo) | Design obejmuje animację lekką (Lottie, GIF), nie postprodukcję filmu | poza platformą |
| Modelowanie i rendering 3D (sceny, siatki, materiały PBR) | Zakres modułu to 2D oraz lekka makieta 3D urządzeń | poza platformą |
| Przygotowalnia poligraficzna (CMYK profilowany, spady, pasery) | Moduł wydaje PDF i konwersję CMYK, nie pełny prepress | częściowo, Grupa H |

Design **wydaje** zasoby do Studio i Apps przez jawne powiązania konfigurowalne — nie przejmuje ich zadań.

### 1.6 Miejsce modułu w architekturze platformy

```
ŚRODOWISKO: WorkSpace lub CodeStudio
        │  boczna nawigacja modułów
        ▼
MODUŁ: Design
        │
        ▼
OKNA OPERACYJNE
Chat Window · Execution Loop Window · Design Board · Assets Panel ·
Prompt Builder · Preview Window · Tokens & System Panel
        │
        ▼
Zasób wizualny gotowy do przekazania:
        ├──► Studio (Studio Editor)      — ilustracja w dokumencie
        └──► Apps (Frontend Workspace)   — element interfejsu produktu
```

---

## 2. Komplet okien operacyjnych modułu — przegląd

| Okno | Rola w module | Forma wiodąca | Warstwa | Sposób wywołania | Punkt wejścia? |
|---|---|---|---|---|---|
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Nie — stale obecne |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; prowadzenie i nadzór zadań generowania, edycji i wydania zasobów | Kolumna sąsiadująca z Chat Window | 1 | Widoczne bez interakcji przy aktywnym zleceniu; poza zleceniem zwinięte do znacznika stanu pętli | Nie — stale obecne |
| Design Board | Kanwa: edycja, komponowanie, makiety | Prawa kolumna dominująca, kanwa nieskończona | 1 | Widoczne bez interakcji | **Tak** |
| Assets Panel | Biblioteka zasobów i pochodzenie | Kolumna boczna, otwierana jako rozszerzenie boczne, siatka miniatur | 2 | Znacznik `Zasoby ▼` w pasku kontekstu | Nie |
| Prompt Builder | Formułowanie poleceń generujących grafikę | Kolumna boczna, otwierana jako rozszerzenie boczne, formularz strukturalny | 2 | Przycisk `Prompt Builder`, polecenie w Chat Window | Nie |
| Preview Window | Ocena, podgląd kontekstowy, decyzja | Kolumna boczna, otwierana jako rozszerzenie boczne, z trybem powiększenia | 2 | Kliknięcie miniatury zasobu, znacznik `Podgląd ▼` | Nie |
| Tokens & System Panel | Tokeny, motywy, kontrast, przewodnik stylu | Kolumna boczna, otwierana jako rozszerzenie boczne, drzewo tokenów i tabele | 3 | Menu kebab `⋮` obszaru roboczego, polecenie w Chat Window | Nie |

```
 Makieta całościowa — Moduł: Design          Dostępność: WorkSpace, CodeStudio
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window            │ Design Board                │ Kolumny
  nawigacja │ Użytkownik ↔ Wykonawca │  kanwa robocza kompozycji   │ boczne
  modułów   │                        │                             │ (rozsze-
            │ ────────────────────   │  [Danaco Console][Design]   │  rzenia
            │ Execution Loop Window  │  [Ultra] [Tryb: kanwa ▼]  ⋮ │  boczne)
            │ Koordynator ↔ Wykonawca│                             │
            │  zlecenie · zadania    │                             │
 ════════════════════════════════════════════════════════════════════════════
```

### 2.1 Warstwy widoczności w module

Moduł stosuje regułę stopniowego ujawniania funkcjonalności: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna.

| Warstwa | Zastosowanie w module Design |
|---|---|
| 1 | Chat Window, Execution Loop Window, kanwa Design Board, pasek kontekstu ze znacznikami projektu, środowiska, modelu i wykonawcy, wskaźnik stanu generowania |
| 2 | Assets Panel, Prompt Builder, Preview Window, wybór silnika generującego, tryb generowania, tryb narzędzia kanwy, poziom powiększenia — wywoływane znacznikiem kontekstowym lub przyciskiem, zwijane samoczynnie po użyciu |
| 3 | Tokens & System Panel, operacje logiczne na ścieżkach, wyrównanie i rozmieszczenie, warianty eksportu, tryb wsadowy, zestaw akcji zasobu — menu kebab `⋮`, menu `☰`, menu kontekstowe kanwy, panele popover |
| 4 | ControlNet i mapy sterujące, rozdzielenie generacji na warstwy, symulacja wad wzroku, inspektor glifów, eksport kodu widoku, konfiguracja punktów końcowych obliczeń — polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny |

Objaśnienie terminu: `popover` — panel przywiązany do kontrolki, otwierany nad treścią
okna i zamykany kliknięciem poza jego obszarem, bez przesłonięcia całego widoku; termin
zapisywany w tym opracowaniu w postaci nieodmiennej.

Każda ukryta funkcja modułu osiągalna jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Makiety w rozdz. 3 przedstawiają stan spoczynku interfejsu: elementy warstwy 1 oraz zwinięte wyzwalacze warstw wyższych.

---

## 3. Specyfikacja okien operacyjnych — pełny arsenał narzędzi

### 3.1 Chat Window (kanał Użytkownik ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Główne okno komunikacji między Użytkownikiem a Wykonawcą; centralny punkt pracy w module i podstawowy mechanizm sterowania wszystkimi procesami twórczymi |
| Zawartość | Polecenia w języku naturalnym, strumień odpowiedzi i wyników, historia rozmowy o bieżącej kompozycji, zatwierdzanie i przerywanie działań |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Pole polecenia w języku naturalnym | Zlecanie generowania, edycji i wydania zasobu, zatwierdzanie i przerywanie działań Wykonawcy | 1 | Widoczne bez interakcji |
| Strumień odpowiedzi i wyników | Prezentacja wyników pracy Wykonawcy wraz z miniaturami zasobów | 1 | Widoczny bez interakcji |
| Wskaźnik zasobu w toku | Miniatura ostatnio wygenerowanego zasobu przypięta nad polem wprowadzania, jako punkt odniesienia kolejnego polecenia | 1 | Widoczny bez interakcji |
| Wzmianki (`@`) | Odwołanie do zasobu z Assets Panel lub elementu Design Board w treści polecenia | 2 | Znak `@` w polu polecenia |
| Załączniki referencyjne | Dołączenie obrazu jako odniesienia stylu z urządzenia lub z Assets Panel | 2 | Ikona załącznika |
| Ocena i iteracja konwersacyjna | Polecenia modyfikujące ostatnio wygenerowany zasób („jaśniejsza wersja”, „usuń tło”, „powiększ margines”) | 2 | Pole polecenia przy przypiętej miniaturze |
| Skrócona ścieżka generowania | Uruchomienie generowania wprost z okna komunikacji, bez otwierania Prompt Buildera | 2 | Polecenie generujące w polu wprowadzania |
| Akcje na wiadomości | Kopiuj polecenie, zapisz jako szablon w Prompt Builderze, powtórz z modyfikacją | 3 | Menu kebab `⋮` wiadomości |
| Wybór wykonawcy i modelu | Zmiana modelu generującego i wykonawcy dla bieżącej sesji | 2 | Znacznik kontekstowy `[Ultra]`, `[Fable 5]` |
| Polecenia funkcji eksperckich | Uruchomienie operacji warstwy 4 (rozdzielenie na warstwy, mapy sterujące, eksport kodu widoku) | 4 | Polecenie języka naturalnego |

**Zachowanie i stany:** okno dostrojone do kontekstu bieżącej kompozycji na Design Board; historia domyślnie odrębna per karta sesji. Stan „generowanie w toku” — wskaźnik postępu przy miniaturze oczekiwanego wyniku. Przerwanie zlecenia wydane w Chat Window zatrzymuje pętlę wykonawczą prezentowaną w Execution Loop Window.

```
 Makieta — Chat Window w module Design (stan spoczynku)
 ═════════════════════════════════════
  Użytkownik ↔ Wykonawca            ⋮
  [Danaco Console] [Design] [Ultra]
 ─────────────────────────────────────
  Historia: „Ilustracja bohatera
  strony — wariant złoty”
  [miniatura ostatniego zasobu ▤]
 ─────────────────────────────────────
  ┌─────────────────────────────────┐
  │ Polecenie…                      │
  └─────────────────────────────────┘
   @  📎                  [ Wyślij ⏎ ]
 ═════════════════════════════════════
```

### 3.2 Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Prowadzenie pętli wykonawczej modułu: dekompozycja zlecenia graficznego na zadania, koordynacja, nadzór nad realizacją i kontrola jakości wyniku |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania modułu (generowanie, przekształcenie, kadrowanie, eksport, wydanie), kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Bieżące zlecenie i dekompozycja | Prezentacja zlecenia graficznego rozłożonego na zadania modułu | 1 | Widoczne bez interakcji |
| Kolejka i stan zadań | Lista zadań pętli (prompt, generowanie, powiększenie neuronowe, usunięcie tła, wektoryzacja, eksport, wydanie) wraz ze stanem każdego | 1 | Widoczna bez interakcji |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator ↔ Wykonawca dla zadań tego modułu | 1 | Widoczna bez interakcji |
| Wskaźniki przebiegu pętli | Liczba iteracji, czas zadania, obciążenie kanału modelu, koszt zlecenia | 1 | Widoczne bez interakcji |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | 1 | Widoczne bez interakcji |
| Wyniki kontroli jakości | Ocena wyniku wobec kryteriów zlecenia (proporcje, paleta, kontrast, format) i decyzja o ponowieniu | 2 | Znacznik `Kontrola ▼` przy zadaniu |
| Szczegóły zadania | Parametry wykonania: model, seed, punkt końcowy, wejściowe zasoby, czas i koszt | 3 | Menu kebab `⋮` zadania |
| Ponowienie z modyfikacją | Powtórzenie zadania ze zmienionym parametrem bez powrotu do Prompt Buildera | 3 | Menu kontekstowe zadania |
| Dziennik pętli | Pełny zapis przebiegu pętli wykonawczej wraz z komunikatami diagnostycznymi | 4 | Polecenie języka naturalnego, tryb administracyjny |

**Zachowanie i stany:** okno prezentuje pętlę wykonawczą zleceń modułu Design. Stan „pętla bezczynna” — okno zwinięte do znacznika stanu przy Chat Window. Stan „pętla w biegu” — kolumna rozwinięta z aktywną kolejką zadań. Stan „oczekiwanie na decyzję” — zadanie wstrzymane do rozstrzygnięcia w Preview Window, przy czym pozostałe zadania kolejki biegną dalej.

```
 Makieta — Execution Loop Window (stan spoczynku)
 ═══════════════════════════════════════
  Koordynator ↔ Wykonawca              ⋮
 ───────────────────────────────────────
  Zlecenie: „Zestaw ilustracji
  brandingowych — kampania Q3”
 ───────────────────────────────────────
  ▸ 1. Kompozycja polecenia    gotowe
  ▸ 2. Generowanie 4 wariantów w toku
  ▸ 3. Usunięcie tła           w kolejce
  ▸ 4. Eksport @1x/@2x         w kolejce
 ───────────────────────────────────────
  Iteracja 2 · kanał modelu 41%
  [ Wstrzymaj ] [ Przerwij ]  Kontrola ▼
 ═══════════════════════════════════════
```

### 3.3 Design Board

| Pole | Treść |
|---|---|
| Cel | Kanwa modułu: edycja rastrowa i wektorowa, komponowanie kompozycji, budowa makiet interfejsu |
| Zawartość | Zgromadzone i zestawione zasoby graficzne, ścieżki, kształty, ramki UI, komponenty |
| Waga wizualna | Prawa kolumna dominująca obszaru roboczego |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Kanwa nieskończona | Nieograniczona przestrzeń robocza z przewijaniem i powiększeniem | 1 | Widoczna bez interakcji |
| Pasek kontekstu kanwy | Znaczniki projektu, środowiska, modelu, trybu narzędzia i poziomu powiększenia | 1 | Widoczny bez interakcji |
| Tryby narzędzia kanwy | Zaznaczanie, pióro, kształt, tekst, pędzel-maska, retusz, ramka UI | 2 | Znacznik `Tryb: … ▼` |
| Panel warstw z drzewem i grupami | Kolejność, widoczność, blokada, zagnieżdżenie elementów kompozycji | 2 | Znacznik `Warstwy ▼` |
| Inspektor właściwości | Wymiary, wypełnienie, obrys, efekty, więzy responsywne, auto-układ | 2 | Zaznaczenie elementu |
| Panel komponentów | Komponenty i warianty stanu, biblioteka UI systemu wizualnego | 2 | Znacznik `Komponenty ▼` |
| Siatka i linie pomocnicze | Wsparcie precyzyjnego rozmieszczenia elementów, z przyciąganiem | 2 | Znacznik `Siatka ▼` |
| Szablony układu | Schematy zestawień: tablica nastroju, siatka porównawcza, arkusz brandingowy, formaty społecznościowe | 3 | Menu `☰` kanwy |
| Narzędzia wyrównania i rozmieszczenia | Wyrównanie do siatki, rozłożenie równomierne, grupowanie i rozgrupowanie | 3 | Pasek kontekstowy zaznaczenia |
| Operacje logiczne na ścieżkach | Suma, różnica, przecięcie, wykluczenie | 3 | Menu kontekstowe zaznaczenia |
| Adnotacje i komentarze | Przypinane uwagi do elementów kompozycji, wątki, oznaczenia osób | 3 | Menu kebab `⋮`, skrót klawiszowy |
| Wersjonowanie kompozycji | Historia stanów tablicy, nazwane wersje, powrót, porównanie | 3 | Menu kebab `⋮` tablicy |
| Eksport kompozycji i fragmentu | Zapis tablicy lub wskazanego obszaru jako plik graficzny lub PDF | 3 | Menu kebab `⋮`, pasek kontekstowy zaznaczenia |
| Zaznaczenie wielokrotne | Operacje zbiorcze na wielu elementach (przesunięcie, usunięcie, eksport) | 2 | Zaznaczenie ramką lub z klawiszem modyfikującym |
| Kursor współpracy | Widoczność pozycji innych osób pracujących nad tą samą tablicą | 1 | Widoczny przy sesji współdzielonej |
| Biblioteka elementów pomocniczych | Kształty, linie prowadzące, ramki makiet i elementy szkicowe | 3 | Menu `☰` kanwy |
| Tryb podglądu tokenów | Przełączenie kanwy między motywem jasnym i ciemnym systemu tokenów | 3 | Menu kebab `⋮` kanwy |
| Rozdzielenie generacji na warstwy | Rozkład wygenerowanego obrazu na edytowalne obiekty | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** odbiera zasoby wygenerowane i zgromadzone w Assets Panel przez przeciągnięcie na kanwę. Stan pusty — zachęta do wygenerowania pierwszego zasobu poleceniem w Chat Window. Stan „edycja współdzielona” — widoczne kursory innych uczestników.

```
 Makieta — Design Board (stan spoczynku)
 ═══════════════════════════════════════════════════════════
  Kompozycja: „Strona główna — wariant A”                  ⋮
  [Design] [Tryb: zaznaczanie ▼] [Warstwy ▼] [Siatka ▼]  ☰
 ───────────────────────────────────────────────────────────
  ┌──────────────────────────────────────────────────┐
  │   ┌────────┐        ┌──────────────┐             │
  │   │ logo   │        │ ilustracja   │             │
  │   └────────┘        │ bohatera     │             │
  │                     └──────────────┘             │
  │        siatka pomocnicza (przyciąganie aktywne)  │
  └──────────────────────────────────────────────────┘
 ═══════════════════════════════════════════════════════════
```

### 3.4 Assets Panel

| Pole | Treść |
|---|---|
| Cel | Przechowywanie wygenerowanych i zgromadzonych zasobów wraz z ich pochodzeniem |
| Zawartość | Wygenerowane grafiki, ilustracje, elementy brandingowe, zasoby importowane |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Siatka miniatur zasobów | Prezentacja zbioru zasobów modułu | 1 | Widoczna po otwarciu panelu |
| Widok siatki i listy | Przełącznik prezentacji miniatur | 2 | Przełącznik `Siatka \| Lista` |
| Wyszukiwanie | Po nazwie, etykiecie lub treści polecenia, które wygenerowało zasób | 2 | Pole wyszukiwania |
| Filtrowanie wg typu | Ilustracja, element brandingowy, makieta interfejsu, ikona, zdjęcie | 2 | Znacznik `Filtr ▼` |
| Filtr po pochodzeniu | AI, edytowane, importowane, wektor, token | 2 | Znacznik `Filtr ▼` |
| Etykiety i kolekcje | Nadawanie etykiet i grupowanie zasobów tematycznie | 2 | Plakietka etykiety, pole etykiet zasobu |
| Oznaczenie ulubionych | Szybki dostęp do najczęściej wykorzystywanych zasobów | 2 | Ikona gwiazdki na miniaturze |
| Panel metadanych i prowenancji | Polecenie źródłowe, model, seed, punkt końcowy, rozdzielczość, data, łańcuch edycji | 3 | Menu kebab `⋮` zasobu |
| Historia wariantów | Kolejne warianty tego samego zasobu z porównaniem | 3 | Menu kebab `⋮` zasobu |
| Pobieranie i eksport zbiorczy | Zaznaczenie wielu zasobów i wydanie w wybranych formatach i skalach | 3 | Menu kebab `⋮` zaznaczenia |
| Import zasobów zewnętrznych | Wyszukanie i wstawienie zasobów z bibliotek zewnętrznych oraz treści próbnych | 3 | Menu `☰` panelu |
| Przeciągnięcie na Design Board | Przeniesienie zasobu wprost do kompozycji roboczej | 1 | Przeciągnięcie miniatury |
| Wysłanie do modułu docelowego | Przekazanie zasobu do Studio Editor lub Frontend Workspace | 2 | Przycisk `→ Moduł` |
| Usuwanie i archiwizacja | Porządkowanie zbioru zasobów — usunięcie następuje od razu, z dostępnym „Cofnij”; potwierdzenie przy usunięciu jest ustawieniem konfiguracyjnym | 3 | Menu kebab `⋮` zasobu |

**Zachowanie i stany:** zasoby trafiają tu bezpośrednio po zakończeniu zadania generowania w pętli wykonawczej. Stan „nowy” — plakietka przy zasobach dodanych od ostatniej wizyty w panelu.

```
 Makieta — Assets Panel (stan spoczynku)
 ════════════════════════════════
  Zasoby (48)                  ☰
  [Szukaj]        [Filtr ▼]    ⋮
 ────────────────────────────────
  🖼 bohater_v3 ★
  🖼 logo_zloto
  🖼 ikona_produktu   [nowy]
 ────────────────────────────────
  [ → Design Board ]  [ → Moduł ]
 ════════════════════════════════
```

### 3.5 Prompt Builder

| Pole | Treść |
|---|---|
| Cel | Precyzyjne formułowanie poleceń generujących grafikę |
| Zawartość | Struktura promptu: tryb generowania, styl, kompozycja, parametry |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Pola strukturalne promptu | Temat, styl, kompozycja, oświetlenie, paleta barw, proporcje kadru, wykluczenia | 1 | Widoczne po otwarciu panelu |
| Przycisk „⚡ Generuj” | Uruchomienie generowania na podstawie skompletowanego promptu | 1 | Widoczny po otwarciu panelu |
| Tryb generowania | Text-to-image, image-to-image, inpainting, outpainting, ikony, makieta z opisu | 2 | Znacznik `Tryb ▼` |
| Wybór silnika generującego | Przełącznik modelu generującego grafikę i punktu końcowego obliczeń | 2 | Znacznik kontekstowy modelu |
| Suwaki parametrów generowania | Kreatywność wobec promptu, ziarno losowości, liczba wariantów, siła wpływu referencji | 2 | Znacznik `Parametry ▼` |
| Obraz referencyjny | Wgranie obrazu jako punktu odniesienia generowania | 2 | Ikona załącznika referencji |
| Biblioteka stylów | Profile stylistyczne do zastosowania i dalszej modyfikacji | 3 | Lista rozwijana pola stylu |
| Paleta z Tokens & System Panel | Wstrzyknięcie zestawu tokenów kolorów jako ograniczenia barwnego generacji | 3 | Menu kebab `⋮` pola palety |
| Generowanie wsadowe | Utworzenie kilku wariantów jednym poleceniem, prezentowanych do wyboru | 2 | Pole `Warianty ▼` |
| Historia poleceń | Lista wcześniej użytych promptów z ponownym użyciem | 3 | Przycisk `Historia` |
| Zapis jako szablon | Utrwalenie skonstruowanego promptu do wielokrotnego użycia | 3 | Menu `Szablony ▼` |
| Wersjonowanie i porównanie promptu | Zestawienie dwóch wersji polecenia obok siebie wraz z ich wynikami | 3 | Menu kebab `⋮` panelu |
| Licznik długości polecenia | Orientacyjna liczba znaków i tokenów promptu | 2 | Widoczny przy polu tematu po rozwinięciu parametrów |
| Panel map sterujących | Wymuszenie pozy, głębi lub krawędzi w generacji na podstawie mapy sterującej | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** zlecenie generowania trafia do pętli wykonawczej, a jego przebieg widoczny jest w Execution Loop Window; wynik trafia do Assets Panel po zakończeniu zadania. Przycisk „⚡ Generuj” jest zawsze aktywny — przy pustym polu wyświetlany jest komunikat podpowiadający, nie blokada.

```
 Makieta — Prompt Builder (stan spoczynku)
 ══════════════════════════════════════
  Nowy prompt                         ⋮
  [Tryb ▼] [Szablony ▼] [Parametry ▼]
 ──────────────────────────────────────
  Temat:       [                    ]
  Styl:        [                  ▼ ]
  Paleta barw: [                  ▼ ]
  Proporcje:   [ 16:9 ▼ ] Warianty:[4▼]
  Wykluczenia: [                    ]
 ──────────────────────────────────────
                        [ ⚡ Generuj ]
 ══════════════════════════════════════
```

### 3.6 Preview Window

| Pole | Treść |
|---|---|
| Cel | Ocena wyniku, podgląd kontekstowy i decyzja o zasobie |
| Zawartość | Wyrenderowany podgląd grafiki, wariantów i prototypu |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne, z trybem powiększenia |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Podgląd zasobu | Prezentacja wygenerowanego lub edytowanego zasobu | 1 | Widoczny po otwarciu panelu |
| Akceptacja, regeneracja, odrzucenie | Trzy akcje decyzyjne wobec wyświetlanego wyniku | 1 | Widoczne po otwarciu panelu |
| Powiększenie pełnoekranowe | Podgląd zasobu w maksymalnym rozmiarze, z powiększeniem i przesuwaniem | 2 | Ikona powiększenia, skrót klawiszowy |
| Suwak porównawczy przed/po | Zestawienie oryginału i zmodyfikowanej wersji na jednym podglądzie | 2 | Znacznik `Przed/po ▼` |
| Przełącznik tła | Tło przezroczyste, białe, czarne lub barwa własna | 2 | Znacznik `Tło ▼` |
| Wybór formatu eksportu | PNG, SVG, JPG, WEBP, AVIF, PDF wraz z kompresją i skalowaniem | 2 | Znacznik `Format ▼` |
| Porównanie wariantów obok siebie | Widok siatki wariantów tego samego polecenia | 3 | Menu kebab `⋮` panelu |
| Osadzenie w ramce urządzenia | Wstawienie zasobu w ramkę telefonu, laptopa lub wizytówki | 3 | Menu kebab `⋮` panelu |
| Podgląd prototypu klikalnego | Nawigacja między ekranami makiety wg grafu połączeń | 3 | Menu kebab `⋮` panelu |
| Przełącznik motywu jasny/ciemny | Ocena zasobu w obu motywach systemu tokenów | 3 | Menu kebab `⋮` panelu |
| Porównanie wizualne | Nałożenie dwóch wersji z podświetleniem różnic pikselowych i strukturalnych | 3 | Menu kebab `⋮` panelu |
| Adnotacja zwrotna | Naniesienie uwagi wprost na podglądzie, przekładanej na kolejne polecenie modyfikujące | 3 | Menu kontekstowe podglądu |
| Podgląd dostępności | Nałożenie oceny kontrastu i symulacji wad wzroku na zasób | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** aktualizuje się po każdej zmianie zaakceptowanej w Design Board oraz po zakończeniu zadania pętli wykonawczej. Stan „oczekuje decyzji” — widoczne trzy akcje decyzyjne, dostępne od razu; ich obecność nie wstrzymuje pracy w pozostałych oknach modułu ani biegu pozostałych zadań pętli.

```
 Makieta — Preview Window (stan spoczynku)
 ══════════════════════════════════════
  Podgląd: bohater_v3, wariant 2/4    ⋮
  [Format ▼] [Tło ▼] [Przed/po ▼]
 ──────────────────────────────────────
  ┌──────────────────────────────────┐
  │                                  │
  │       [ podgląd grafiki ]        │
  │                                  │
  └──────────────────────────────────┘
 ──────────────────────────────────────
  [ Akceptuj ] [ Regeneruj ] [ Odrzuć ]
 ══════════════════════════════════════
```

### 3.7 Tokens & System Panel

| Pole | Treść |
|---|---|
| Cel | Definiowanie i wydawanie systemu projektowego: tokenów, motywów, reguł dostępności i przewodnika stylu |
| Zawartość | Drzewo tokenów, motywy, tabela kontrastu, powiązania z komponentami |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

| Sekcja | Zawartość | Warstwa | Sposób wywołania |
|---|---|---|---|
| Drzewo tokenów | Kolory, typografia, odstępy, promienie, cienie w strukturze W3C Design Tokens | 3 | Menu kebab `⋮` obszaru roboczego |
| Motywy | Zestawy nadpisań (jasny, ciemny, warianty marki) wraz z przełącznikiem podglądu | 3 | Zakładka `Motywy` panelu |
| Kontrast i dostępność | Tabela par tokenów z oceną WCAG oraz symulacja wad wzroku | 3 | Zakładka `Dostępność` panelu |
| Eksport i import | Wydanie tokenów do CSS, SCSS, Tailwind, JS, iOS, Android; wczytanie z JSON | 3 | Menu kebab `⋮` panelu |
| Powiązania | Wskazanie komponentów korzystających z danego tokenu | 3 | Kliknięcie tokenu w drzewie |
| Przewodnik stylu | Generowana strona-dokumentacja systemu projektowego | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |
| Inspektor glifów i par krojów | Przegląd znaków, ligatur, wariantów OpenType oraz dobór par krojów | 4 | Polecenie języka naturalnego, tryb administracyjny |

**Zachowanie i stany:** zmiana wartości tokenu propaguje się do komponentów, kanwy i podglądu. Stan „naruszenie kontrastu” — plakietka ostrzegawcza przy parze tokenów łamiącej wymagany próg; ostrzeżenie nie wstrzymuje pracy.

```
 Makieta — Tokens & System Panel (stan spoczynku)
 ═══════════════════════════════
  Tokeny i system              ⋮
  [Tokeny] [Motywy] [Dostępność]
 ───────────────────────────────
  ▸ kolor
      primary        #0B1F3A
      surface        #F5F5F2
      danger         #B3261E
  ▸ typografia
  ▸ odstępy
  ▸ promienie
 ═══════════════════════════════
```

---

## 4. Katalog funkcji i narzędzi

Każda pozycja opisana jest przez nazwę, działanie i zależności (biblioteki warstwy serwerowej w Go, formaty, modele, integracje). Warstwa serwerowa (rdzeń) platformy jest w Go; funkcje wymagające modelu neuronowego wykonywane są przez kanał modelu (API) lub przez własny punkt końcowy obliczeń, zgodnie ze strategią modeli platformy. Zadania katalogu realizowane są w pętli wykonawczej prezentowanej w Execution Loop Window.

### 4.1 Grupa A — Generowanie grafiki przez AI

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Text-to-image | Generuje obraz z opisu skomponowanego w Prompt Builderze | Modele obrazowe przez kanał API; własny punkt końcowy GPU; format PNG/WEBP |
| 2 | Image-to-image | Przekształca obraz referencyjny wg promptu, z regulacją siły wpływu | Model obrazowy z parametrem `strength`; wejście: dowolny raster |
| 3 | Inpainting (domalowanie fragmentu) | Zamienia zaznaczony obszar obrazu nowym, wpasowanym w otoczenie | Model inpaint i maska alfa (PNG RGBA); `govips` do składania masek |
| 4 | Outpainting (rozszerzenie kadru) | Domalowuje treść poza pierwotną ramką obrazu | Model outpaint; `disintegration/imaging` do powiększania płótna |
| 5 | Wariacje | Tworzy N wariantów tego samego zasobu (seed i odchylenie) do wyboru | Parametr seed i liczba wariantów; generowanie wsadowe Prompt Buildera |
| 6 | Upscaling neuronowy | Powiększa obraz 2×/4×/8× bez utraty ostrości | Real-ESRGAN / SwinIR przez punkt końcowy GPU; ONNX przez `yalue/onnxruntime_go` |
| 7 | Usuwanie tła | Odcina obiekt od tła, tworzy kanał alfa | U²-Net / RMBG-1.4 (ONNX) przez punkt końcowy; wynik PNG RGBA |
| 8 | Transfer stylu | Przenosi styl obrazu-wzorca na treść obrazu-bazy | IP-Adapter / model transferu stylu; obraz referencyjny Prompt Buildera |
| 9 | Wektoryzacja (raster→SVG) | Zamienia bitmapę na czyste ścieżki wektorowe | `dennwc/gotrace` (potrace) lub wektoryzacja przez kanał API; wynik SVG |
| 10 | Generowanie ikon i logotypów | Tworzy spójny zestaw ikon i znaków w jednym stylu | Model obrazowy z wymuszeniem stylu przez preset; wyjście SVG/PNG |
| 11 | Rozdzielenie na warstwy | Rozkłada wygenerowany obraz na obiekty i warstwy edytowalne | Segment Anything przez punkt końcowy; maski przenoszone na warstwy Design Board |
| 12 | Naprawa twarzy i detalu | Poprawia zniekształcone twarze oraz drobne detale generacji | GFPGAN / CodeFormer przez punkt końcowy GPU |
| 13 | Sterowanie kompozycją (mapy sterujące) | Wymusza pozę, szkielet, głębię lub krawędzie w generacji wg szkicu | ControlNet (canny, depth, pose) przez punkt końcowy; mapa sterująca jako obraz |

### 4.2 Grupa B — Edycja rastrowa

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Kadrowanie i prostowanie | Przycina, obraca, prostuje horyzont, zmienia proporcje | `disintegration/imaging` (Crop, Rotate); reguły proporcji 1:1, 16:9, 4:5 |
| 2 | Skalowanie i resampling | Zmienia rozmiar z wyborem algorytmu (Lanczos, biliniowy) | `govips` / `bimg` (libvips, Lanczos3); zachowanie proporcji |
| 3 | Warstwy rastrowe | Niezależne warstwy z trybami mieszania i krycia | Kompozytor warstw (multiply, screen, overlay) — compositing `govips` |
| 4 | Maski warstw | Nieniszcząca maska alfa, malowana lub z zaznaczenia | Kanał alfa PNG; pędzel maski w kliencie (Canvas API), złożenie w Go |
| 5 | Korekcja barwna | Jasność, kontrast, ekspozycja, temperatura, krzywe, HSL | `govips` (LUT, krzywe); podgląd na żywo w kliencie |
| 6 | Retusz (klonowanie i leczenie) | Klonuje i wygładza fragmenty, usuwa niedoskonałości | Inpaint lokalny („Inpainting (domalowanie fragmentu)”) oraz pędzel klonujący |
| 7 | Filtry i efekty | Rozmycie, wyostrzenie, ziarno, winieta, cień, poświata | `govips` (gaussian blur, sharpen); efekty warstwy |
| 8 | Zaznaczanie obiektu | Automatyczne zaznaczenie obiektu jednym kliknięciem | Segmentacja („Rozdzielenie na warstwy”) lub różdżka progowa; maska zaznaczenia |
| 9 | Korekcja perspektywy | Prostuje zniekształcenia perspektywiczne i obiektywu | Transformacja homograficzna; `fogleman/gg` do przekształceń |
| 10 | Tryb wsadowy edycji | Stosuje ten sam zestaw operacji do wielu zasobów naraz | Potok operacji zapisany jako profil; kolejka modułu Automations |

### 4.3 Grupa C — Grafika i edycja wektorowa

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Narzędzie pióra (ścieżki Béziera) | Rysuje i edytuje krzywe oraz węzły ścieżek | SVG `<path>`; edycja węzłów w kliencie; serializacja `ajstarks/svgo` |
| 2 | Kształty podstawowe i złożone | Prostokąty, elipsy, wielokąty, gwiazdy, zaokrąglenia | Prymitywy SVG; parametry promienia i liczby wierzchołków |
| 3 | Operacje logiczne (boolean) | Suma, różnica, przecięcie, wykluczenie ścieżek | Booleany `tdewolff/canvas`; wynik jako pojedyncza ścieżka |
| 4 | Obrys i wypełnienie | Grubość, zakończenia, gradient, deseń, reguła wypełnienia | Atrybuty SVG (stroke, fill, gradient defs); edytor gradientu |
| 5 | Tekst na ścieżce i obrys tekstu | Układa tekst wzdłuż krzywej, zamienia w kontury | `golang/freetype` / `go-text/typesetting`; SVG `textPath` |
| 6 | Optymalizacja SVG | Czyści i minimalizuje kod SVG bez utraty jakości | Reguły czyszczenia (metadane, scalanie ścieżek); wyjście SVG |
| 7 | Symbole i instancje | Definicja wielokrotnego elementu z propagacją zmian | SVG `<symbol>` / `<use>`; rejestr symboli w projekcie |
| 8 | Siatka i przyciąganie precyzyjne | Rozmieszczanie do siatki, prowadnic, punktów, pikseli | Silnik przyciągania w kliencie; jednostki px/rem |
| 9 | Eksport wektorowy | Wydaje czysty SVG, PDF wektorowy i EPS | `go-pdf/fpdf` (wektor→PDF); serializacja SVG |

### 4.4 Grupa D — Projektowanie UI i makiety

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Ramki i obszary robocze | Wydzielone ekrany o zdefiniowanych rozmiarach urządzeń | Profile rozmiarów urządzeń; model sceny Design Board |
| 2 | Komponenty i warianty | Element wielokrotny z wariantami stanu (wskazanie kursorem, aktywny, wyłączony) | Rejestr komponentów oraz zestaw właściwości wariantu |
| 3 | Auto-układ | Automatyczne rozmieszczanie z odstępami i dopasowaniem do treści | Silnik flex/stack w kliencie; parametry gap, padding, wyrównanie |
| 4 | Więzy responsywne | Zachowanie elementów przy zmianie rozmiaru ramki | Reguły kotwiczenia i rozciągania; przeliczanie układu |
| 5 | Prototypowanie i przejścia | Łączy ekrany interakcjami i animacjami przejść | Graf połączeń ekranów; podgląd klikalny w Preview Window |
| 6 | Makieta szkicowa niskiej wierności | Szybkie makiety szkicowe z biblioteką prostych elementów | Biblioteka elementów szkicowych Design Board |
| 7 | Biblioteka UI | Gotowe komponenty (przyciski, pola, karty) do składania makiet | Komponenty zgodne z systemem wizualnym Danaco; import z tokenów (E) |
| 8 | Generowanie makiety z opisu | Tworzy szkielet ekranu z polecenia tekstowego | Model generujący układ → drzewo komponentów; mapowanie na „Komponenty i warianty” i „Biblioteka UI” |
| 9 | Import ze zrzutu ekranu | Odtwarza edytowalną makietę z obrazu istniejącego interfejsu | Rozpoznanie układu (segmentacja i OCR `otiai10/gosseract`) → komponenty |
| 10 | Siatki układu | Kolumny, rynny, moduły, siatka bazowa | Definicja siatki (kolumny, gap, margines); nakładka na ramkę |

### 4.5 Grupa E — Tokeny projektowe i system projektowy

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Tokeny kolorów | Definicja nazwanych barw semantycznych (primary, surface, danger) | Format W3C Design Tokens (JSON, `$value`/`$type`); walidacja aliasów |
| 2 | Tokeny typografii | Skala rozmiarów, krój, interlinia, grubość, tracking | Tokeny typograficzne W3C; podgląd skali |
| 3 | Tokeny odstępów i promieni | Skala spacing, radius, rozmiary, cienie, z-index | Tokeny wymiarowe W3C; jednostki px/rem |
| 4 | Motyw jasny i ciemny | Warianty tokenów dla trybów i marek | Zestawy nadpisań tokenów; przełącznik trybu w podglądzie |
| 5 | Eksport tokenów do kodu | Wydaje tokeny jako CSS vars, SCSS, Tailwind, JS, iOS, Android | Silnik transformacji tokenów; wyjścia `.css`, `.scss`, `tailwind.config`, `.json`, `.xml` |
| 6 | Import tokenów | Wczytuje istniejące tokeny do modułu | Parser JSON W3C i `tailwind.config`; mapowanie na model tokenów |
| 7 | Powiązanie tokenów z komponentami | Komponenty UI („Komponenty i warianty”) czerpią wartości z tokenów, zmiana propaguje | Referencje tokenów we właściwościach komponentu; reaktywna aktualizacja |
| 8 | Dokumentacja systemu | Generuje przewodnik: kolory, typografia, komponenty | Render HTML przewodnika; wydanie do Library i Studio |

### 4.6 Grupa F — Kolor

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Generator palet | Tworzy paletę z reguły harmonii (mono, komplementarna, triada) | `lucasb-eyer/go-colorful` (konwersje HSL/Lab, harmonie) |
| 2 | Ekstrakcja palety z obrazu | Wyciąga dominujące barwy z zasobu | Kwantyzacja (median cut, k-means); `go-colorful` |
| 3 | Kontroler kontrastu WCAG | Liczy współczynnik kontrastu i ocenę AA/AAA dla par barw | Wzór luminancji WCAG 2.1; `go-colorful`; podgląd wyniku |
| 4 | Symulacja wad wzroku | Podgląd projektu w protanopii, deuteranopii i tritanopii | Macierze symulacji daltonizmu na buforze obrazu |
| 5 | Edytor gradientów | Gradienty liniowe, promieniste i kątowe z wieloma stopniami | SVG/CSS gradient defs; interpolacja w przestrzeni Lab (`go-colorful`) |
| 6 | Konwersja przestrzeni barw | HEX ↔ RGB ↔ HSL ↔ Lab ↔ CMYK, próbki nazwane | `go-colorful`; tablice barw nazwanych |
| 7 | Sprawdzian dostępności kolorem | Wskazuje pary tokenów łamiące kontrast w całym systemie | Reguły WCAG na zestawie tokenów („Tokeny kolorów”); raport naruszeń |

### 4.7 Grupa G — Ikony i typografia

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Biblioteka ikon | Przeszukiwalny katalog ikon otwartoźródłowych do wstawienia | Zestawy SVG (Lucide, Heroicons, Tabler); indeks nazw i etykiet |
| 2 | Edytor ikony na siatce | Rysuje i poprawia ikonę na siatce 24×24 z wyrównaniem do pikseli | Edytor SVG (Grupa C), siatka ikony, kształty prowadzące |
| 3 | Generowanie zestawu ikon | Tworzy spójny komplet ikon w jednym stylu z listy pojęć | Model („Generowanie ikon i logotypów”) z wymuszeniem siatki i grubości; wyjście SVG |
| 4 | Font ikon i sprite | Pakuje ikony w font ikon lub sprite SVG z symbolami | Generator fontu przez moduł Terminal lub sprite `<symbol>`; `.woff2`, `.svg` |
| 5 | Generator faviconów | Wydaje komplet faviconów i ikon aplikacji ze źródła | `Kodeworks/golang-image-ico`; render 16/32/180/512 px; `.ico`, `.png`, manifest |
| 6 | Dobór par krojów | Proponuje pasujące zestawienia krojów nagłówka i tekstu | Reguły typograficzne i model; podgląd zestawień |
| 7 | Podgląd i osadzenie krojów sieciowych | Podgląd kroju w tekście próbnym, generacja `@font-face` | Metadane fontu (`golang/freetype`); podzbiór glifów |
| 8 | Inspektor glifów | Przegląd znaków, ligatur i wariantów OpenType kroju | Parser OpenType; render glifów |

### 4.8 Grupa H — Eksport, podgląd, handoff

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Eksport wieloformatowy | Wydaje zasób jako PNG, JPG, WEBP, AVIF, SVG, PDF, ICO | `govips` / `bimg`, `chai2010/webp`, AVIF przez libvips; SVG i PDF natywnie |
| 2 | Kompresja i optymalizacja | Redukuje wagę pliku sterując jakością i paletą | `govips` (jakość, chroma subsampling), pngquant/oxipng przez Terminal |
| 3 | Skalowanie @1x/@2x/@3x | Generuje warianty gęstości pikseli jednym poleceniem | Mnożniki eksportu; resize `govips`; nazwy `@2x`, `@3x` |
| 4 | Cięcie na fragmenty | Eksportuje wskazane obszary makiety jako osobne pliki | Definicja ramek eksportu; wsadowe cięcie `imaging.Crop` |
| 5 | Handoff i inspekcja | Udostępnia wymiary, odstępy, kolory, typografię i CSS elementu | Odczyt modelu sceny → panel wartości i generacja CSS |
| 6 | Eksport kodu widoku | Zamienia makietę lub element w kod widoku (CSS, Tailwind, React) | Serializacja drzewa komponentów → JSX/HTML+CSS; mapowanie tokenów („Eksport tokenów do kodu”); wydanie do Apps |
| 7 | Eksport tablicy jako PDF lub obraz | Zapisuje całą kompozycję Design Board do jednego pliku | `go-pdf/fpdf`; render kanwy → PDF/PNG |
| 8 | Osadzenie w ramce urządzenia | Wstawia zasób w realistyczną ramkę telefonu, laptopa lub wizytówki | Biblioteka ramek PNG/SVG; kompozycja perspektywiczna (`fogleman/gg`) |
| 9 | Eksport animacji lekkiej | Wydaje animację jako JSON Lottie, GIF lub APNG | Format Lottie JSON; `image/gif` (biblioteka standardowa); render klatek |
| 10 | Metadane i profil barwny | Osadza i oczyszcza EXIF, ICC oraz dane autorstwa przy eksporcie | `dsoprea/go-exif`; osadzanie profili ICC przez libvips |

### 4.9 Grupa I — Marketing, szablony, zasoby zewnętrzne

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Szablony formatów społecznościowych | Gotowe rozmiary postów, relacji i okładek dla platform społecznościowych | Profile rozmiarów i szablony układu Design Board |
| 2 | Zestawy rozmiarów kampanii | Generuje ten sam projekt w wielu formatach reklamowych naraz | Reguła przeskalowania układu (auto-układ) na zestaw ramek |
| 3 | Biblioteka szablonów brandingowych | Arkusze brandingu, tablice nastroju, prezentacje, wizytówki | Szablony układu Design Board; katalog szablonów modułu |
| 4 | Import zasobów zewnętrznych | Wyszukuje i wstawia zasoby z bibliotek zdjęć i ilustracji | Integracje API bibliotek zewnętrznych (klucze jawne w oknie konfiguracji); import do Assets Panel |
| 5 | Treści próbne | Wstawia realistyczne treści próbne (nazwiska, teksty, awatary) | Generator danych próbnych i awatary AI („Text-to-image”); zasilanie makiet |
| 6 | Znak wodny i branding wsadowy | Nakłada logo lub znak wodny na zestaw zasobów | Kompozycja warstwy („Warstwy rastrowe”) w trybie wsadowym („Tryb wsadowy edycji”) |

### 4.10 Grupa J — Współpraca, wersjonowanie, organizacja

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| 1 | Warstwy i drzewo obiektów | Panel warstw z widocznością, blokadą, grupami i zagnieżdżeniem | Model sceny Design Board |
| 2 | Wersjonowanie kompozycji | Historia stanów tablicy, nazwane wersje, powrót, porównanie | Migawki sceny ([Model danych](../architektura/model-danych.md) — wersjonowanie); porównanie wizualne |
| 3 | Porównanie wizualne zasobów | Nakłada dwie wersje i podświetla różnice pikselowe oraz strukturalne | `corona10/goimagehash` (pHash), różnica pikseli; suwak przed/po |
| 4 | Komentarze i adnotacje | Przypięte uwagi do elementów, wątki, oznaczenia osób | Model adnotacji; powiązanie z Chat Window |
| 5 | Kursory i obecność | Widoczność kursorów współpracujących osób na tablicy | Kanał WebSocket platformy |
| 6 | Etykiety, kolekcje, wyszukiwanie | Etykietowanie i przeszukiwanie zasobów po nazwie, etykiecie i prompcie | Assets Panel; indeks pełnotekstowy promptów |
| 7 | Metadane pochodzenia (prowenancja) | Zapisuje model, prompt, seed, datę i łańcuch edycji zasobu | Panel metadanych Assets Panel; zapis w Modelu danych; zgodność z blokiem prowenancji klienta |

Grupy A–J obejmują łącznie **88 funkcji i narzędzi** pokrywających obszar tematyczny modułu. Część funkcji opiera się wyłącznie na kodzie klienta (kanwa, edytor SVG), część na bibliotekach Go warstwy serwerowej (raster, eksport, kolor), a część korzysta z zasobu GPU i modelu neuronowego (generowanie, powiększanie neuronowe, segmentacja).

---

## 5. Katalog elementów interfejsu

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole polecenia | `.dn-pole-kontrolka` na całą szerokość kolumny | Wydawanie poleceń Wykonawcy w języku naturalnym | Element wiodący lewej kolumny | 1 | Widoczne bez interakcji | domyślny, aktywny, wysyłanie | Wysyła polecenie do Wykonawcy i uruchamia zadanie w pętli wykonawczej | Chat Window |
| Karta zadania pętli | Wiersz kolejki `.dn-karta` | Reprezentacja pojedynczego zadania pętli wykonawczej | Wąska karta listowa | 1 | Widoczna bez interakcji | w kolejce, w toku, gotowe, ponowione, przerwane | Kliknięcie rozwija szczegóły; menu kebab udostępnia ponowienie | Execution Loop Window |
| Sterowanie przebiegiem pętli | `.dn-btn--zarys` w zestawie trzech | Wstrzymanie, wznowienie i przerwanie pętli wykonawczej | Trzy małe przyciski w rzędzie | 1 | Widoczne bez interakcji | domyślny, wstrzymany, przerwany | Zmienia stan pętli natychmiast, bez okna nakładkowego z potwierdzeniem | Execution Loop Window |
| Znacznik kontekstowy | `.dn-plakietka--sygnal` | Prezentacja i zmiana środowiska, projektu, modelu, wykonawcy | Bardzo mała pigułka | 1 | Widoczny bez interakcji; kliknięcie otwiera selektor | domyślny, otwarty selektor | Kliknięcie otwiera selektor warstwy 2, który zwija się po wyborze | Chat Window, Design Board |
| Miniatura zasobu | Kafel obrazu w siatce `.dn-karta--klikalna` | Reprezentacja pojedynczego zasobu wizualnego | Mały–średni kafel kwadratowy | 1 | Widoczna po otwarciu panelu | domyślny, wskazanie kursorem, nowy (plakietka), wybrany (obrys sygnałowy) | Kliknięcie otwiera Preview Window; przeciągnięcie przenosi na Design Board | Assets Panel |
| Przycisk „⚡ Generuj” | `.dn-btn--sygnal` | Najważniejsze CTA okna Prompt Builder | Średni przycisk wypełniony, cień `--dn-cien-sygnal` | 1 | Widoczny po otwarciu panelu | domyślny, generowanie (wskaźnik postępu) | Zawsze aktywny — zleca generowanie do pętli wykonawczej; przy pustym poleceniu komunikat podpowiadający zamiast blokady | Prompt Builder |
| Suwak parametru | `.dn-suwak` w wariancie zakresowym | Regulacja parametru generowania | Mały element liniowy | 2 | Znacznik `Parametry ▼` | domyślny, przeciągany | Zmienia wartość liczbową parametru promptu | Prompt Builder |
| Pole stylu i palety | `.dn-pole-kontrolka` | Wybór profilu stylistycznego lub palety barw | Małe pole rozwijane | 2 | Kliknięcie pola | domyślny, otwarte (lista pozycji) | Wybór pozycji aktualizuje podgląd struktury promptu | Prompt Builder |
| Suwak porównawczy przed/po | Uchwyt przeciągalny na osi poziomej | Wizualne porównanie dwóch wersji zasobu | Element średni, na całą szerokość podglądu | 2 | Znacznik `Przed/po ▼` | domyślny, przeciągany | Przesuwa granicę widoczności między wersją „przed” i „po” | Preview Window |
| Przełącznik tła podglądu | Mały zestaw ikon (szachownica, biel, czerń) | Zmiana tła pod zasobem z przezroczystością | Bardzo mała grupa przycisków ikonowych | 2 | Znacznik `Tło ▼` | domyślny, wybrany | Zmienia tło podglądu natychmiastowo | Preview Window |
| Przyciski decyzyjne | `.dn-btn` w trzech wariantach (główny, zarys, duch) | Zamknięcie cyklu oceny wygenerowanego zasobu | Trzy małe–średnie przyciski w rzędzie | 1 | Widoczne po otwarciu panelu | domyślne, po decyzji | Akceptacja przenosi zasób do Assets Panel jako finalny; odrzucenie usuwa wariant i zachowuje go w historii | Preview Window |
| Panel warstw | Drzewo z ikonami widoczności i blokady | Zarządzanie kolejnością i widocznością elementów kompozycji | Kolumna boczna, wąska | 2 | Znacznik `Warstwy ▼` | domyślny, warstwa ukryta, zablokowana | Kliknięcie ikony przełącza widoczność lub blokadę warstwy | Design Board |
| Uchwyty wyrównania i rozmieszczenia | Mały pasek ikon przy zaznaczeniu | Wyrównanie i rozłożenie zaznaczonych elementów | Bardzo mały pasek kontekstowy | 3 | Pasek kontekstowy zaznaczenia wielokrotnego | widoczny tylko przy zaznaczeniu wielokrotnym | Kliknięcie stosuje wybrane wyrównanie do zaznaczonych elementów | Design Board |
| Znacznik adnotacji | Mała pinezka na kanwie | Przypięcie komentarza do elementu kompozycji | Bardzo mała ikonka | 3 | Menu kebab `⋮`, skrót klawiszowy | domyślny, rozwinięty (dymek treści) | Kliknięcie otwiera i zamyka treść komentarza | Design Board |
| Etykiety i kolekcje | `.dn-plakietka` w wariancie neutralnym | Kategoryzacja zasobów | Bardzo małe pigułki | 2 | Widoczne przy miniaturze zasobu | domyślny, aktywny filtr (akcent sygnałowy) | Kliknięcie filtruje widok Assets Panel do wybranej etykiety | Assets Panel |
| Skrót „Wyślij do modułu” | `.dn-btn--duch`, ikonowy | Przekazanie zasobu do Studio lub Apps | Mały przycisk bez obrysu | 2 | Przycisk `→ Moduł` | domyślny, wskazanie kursorem | Otwiera wybór modułu docelowego i wykonuje transfer | Assets Panel, Preview Window |
| Wiersz tokenu | `.dn-listwa-pozycja` z próbką barwy | Prezentacja i edycja pojedynczego tokenu systemu | Wąski wiersz drzewa | 3 | Menu kebab `⋮` obszaru roboczego | domyślny, edytowany, naruszenie kontrastu | Zmiana wartości propaguje się do komponentów i kanwy | Tokens & System Panel |

---

## 6. Przepływy pracy w module

### 6.1 Od polecenia do zaakceptowanego zasobu

```
Chat Window — polecenie Użytkownika (kanał Użytkownik ↔ Wykonawca)
        │
        ▼
Execution Loop Window — Koordynator dekomponuje zlecenie na zadania
        │  (prompt → generowanie → przekształcenie → eksport)
        ▼
Prompt Builder — doprecyzowanie polecenia (temat, styl, paleta, parametry)
        │  [ ⚡ Generuj ]
        ▼
Wykonawca realizuje zadanie; przebieg widoczny w Execution Loop Window
        │
        ▼
Wynik trafia do Assets Panel (jeden lub kilka wariantów wsadowych)
        │
        ▼
Preview Window — ocena wyniku
        │
        ├──► [Akceptuj]    → zasób finalny w Assets Panel
        ├──► [Regeneruj]   → ponowienie zadania w pętli wykonawczej
        └──► [Odrzuć]      → wariant odrzucony, historia zachowana
```

### 6.2 Zestawianie kompozycji z wielu zasobów

```
Assets Panel (zasoby zgromadzone)
        │  przeciągnięcie na kanwę
        ▼
Design Board — zestawienie, wyrównanie, warstwy, siatka pomocnicza
        │
        ├──► adnotacje i komentarze do elementów
        ├──► wersjonowanie kompozycji (kolejne stany tablicy)
        │
        ▼
Eksport kompozycji całościowej lub przekazanie pojedynczych
elementów dalej (Studio, Apps)
```

### 6.3 Budowa i wydanie systemu projektowego

```
Tokens & System Panel — definicja tokenów (kolory, typografia, odstępy)
        │
        ├──► kontrola kontrastu WCAG i symulacja wad wzroku
        ├──► powiązanie tokenów z komponentami Design Board
        │
        ▼
Chat Window — polecenie wydania systemu
        │
        ▼
Execution Loop Window — zadania eksportu (CSS, SCSS, Tailwind, JS, iOS, Android)
        │
        ▼
Wydanie do modułu Apps oraz przewodnik stylu do Library
```

### 6.4 Przekazanie zasobu do modułu docelowego

```
Assets Panel / Preview Window
        │  [ → Studio ]                    [ → Apps · Frontend Workspace ]
        ▼                                              ▼
Studio Editor (moduł Studio)             Frontend Workspace (moduł Apps)
zasób wykorzystany przy redagowaniu      zasób wykorzystany jako element
dokumentu                                 interfejsu produktu
   (powiązanie konfiguracyjne — jawne, ustanawiane przez użytkownika)
```

---

## 7. Punkty sterowania z okna konfiguracji

Operator personalizuje moduł z okna konfiguracji ([Model konfiguracji](../architektura/model-konfiguracji.md)), z zachowaniem warstwowości (globalna → środowisko → projekt → sesja) i zasady „brak ustawienia = wartość domyślna”. Wszystkie klucze są jawne, żaden nie blokuje pracy.

| # | Punkt sterowania | Co personalizuje | Zakres konfiguracji (5.x) | Wartość domyślna |
|---|---|---|---|---|
| 1 | Domyślny silnik generujący | Model text-to-image używany bez ręcznego wyboru | 5.4 Zachowanie modeli | pierwszy skonfigurowany kanał |
| 2 | Punkty końcowe GPU zadań neuronowych | Adresy własnych punktów końcowych dla powiększania neuronowego, usuwania tła, segmentacji i map sterujących | 5.4, 5.7 Rozszerzenia | brak (funkcje AI przez kanał API) |
| 3 | Klucze API modeli i bibliotek zewnętrznych | Jawne klucze dostępu do modeli obrazowych i bibliotek zasobów („Import zasobów zewnętrznych”) | 5.4, 5.8 Integracje | brak (wprowadzane przez Operatora) |
| 4 | Domyślne formaty i jakość eksportu | Profil formatu (PNG, WEBP, AVIF), poziom kompresji, zestaw skal @1/2/3x | 5.1 Aplikacja | PNG, jakość 90, @1x |
| 5 | Domyślna paleta i system tokenów | Zestaw tokenów („Tokeny kolorów”–„Motyw jasny i ciemny”) wstrzykiwany do nowych projektów | 5.8 Integracje | system wizualny Danaco |
| 6 | Reguły dostępności | Próg kontrastu (AA/AAA) i ostrzeżenia w Tokens & System Panel | 5.1 Aplikacja | AA, ostrzeżenia włączone |
| 7 | Powiązanie z modułem Studio | Sposób przekazywania zasobów do Studio Editor (ręcznie lub synchronizacja) | 5.8 Integracje | wyłączone (ręczne wysłanie) |
| 8 | Powiązanie z modułem Apps | Sposób przekazywania kodu i zasobów do Frontend Workspace | 5.8 Integracje | wyłączone (ręczne wysłanie) |
| 9 | Powiązanie z Library | Archiwizacja zaakceptowanych zasobów w Library | 5.8 Integracje | wyłączone |
| 10 | Wsadowe zadania przez Automations | Korzystanie trybu wsadowego („Tryb wsadowy edycji”, „Zestawy rozmiarów kampanii”, „Znak wodny i branding wsadowy”) z kolejek modułu Automations | 5.3 Akcje, 5.8 | wyłączone |
| 11 | Potwierdzanie usuwania zasobu | Wymóg potwierdzenia usunięcia w Assets Panel (zawsze z „Cofnij”) | 5.1 Aplikacja | bez potwierdzenia, „Cofnij” dostępne |
| 12 | Zapis prowenancji | Szczegółowość metadanych pochodzenia („Metadane pochodzenia (prowenancja)”) | 5.9 Historia, 5.10 Pamięć | zapis pełny |
| 13 | Domyślny tryb Preview | Motyw (jasny, ciemny), domyślna ramka urządzenia, tło przezroczystości | 5.1, 5.12 Karty sesji | motyw systemowy, bez ramki, szachownica |
| 14 | Widoczność wariantów i seed | Liczba wariantów domyślnych, przypinanie seed dla powtarzalności | 5.4 | 4 warianty, seed swobodny |
| 15 | Ścieżki eksportu docelowego | Katalog i kolekcja Library oraz miejsce w projekcie Apps dla wydań | 5.8, 5.13 Komponenty | brak (wybór przy wydaniu) |
| 16 | Zachowanie Execution Loop Window | Widoczność kolumny pętli w spoczynku, poziom szczegółowości komunikatów sterujących, próg automatycznego ponowienia zadania | 5.3 Akcje, 5.9 Historia | kolumna zwinięta w spoczynku, komunikaty skrócone, jedno ponowienie |
| 17 | Warstwy widoczności interfejsu modułu | Zestaw funkcji ujawnianych w warstwach 2–4 zależnie od roli użytkownika | 5.1 Aplikacja, 5.16 Warstwy widoczności funkcji | warstwy 1–3 dla roli podstawowej, warstwa 4 dla roli rozszerzonej |

Każdy punkt sterowania opatrzony jest objaśnieniem kontekstowym `[?]` — zgodnie z 3.3 Modelu konfiguracji.

---

## 8. Stany, dane i powiązania z innymi modułami

### 8.1 Dane wykorzystywane przez AI w module

| Źródło danych | Okno pochodzenia | Uwaga |
|---|---|---|
| Polecenia Użytkownika | Chat Window | Treść polecenia, kontekst rozmowy, załączniki referencyjne |
| Zlecenia i zadania pętli | Execution Loop Window | Dekompozycja zlecenia, stan zadań, wyniki kontroli jakości |
| Polecenia generujące | Prompt Builder | Struktura promptu i historia poleceń |
| Zgromadzone zasoby | Assets Panel | Metadane, etykiety, warianty, prowenancja |
| Kompozycja robocza | Design Board | Układ, warstwy, adnotacje, komponenty |
| System projektowy | Tokens & System Panel | Tokeny, motywy, reguły dostępności |

### 8.2 Stany zasobu wizualnego

| Stan | Znaczenie | Gdzie widoczny |
|---|---|---|
| W generowaniu | Zasób w trakcie tworzenia przez model | Execution Loop Window, Assets Panel (miniatura w budowie) |
| Nowy | Wygenerowany, jeszcze nieoceniony | Assets Panel (plakietka „nowy”) |
| Zaakceptowany | Zatwierdzony jako wynik finalny | Preview Window, Assets Panel |
| Odrzucony | Wariant odrzucony, zachowany w historii | Preview Window (historia wariantów) |
| W kompozycji | Zasób umieszczony na Design Board | Design Board (panel warstw) |
| Przekazany | Wysłany do modułu docelowego | Assets Panel (znacznik miejsca docelowego) |

### 8.3 Powiązania konfigurowalne z innymi modułami

```
                       ┌──────────────────────────────┐
                       │            DESIGN            │
                       │      (moduł, ta karta)       │
                       └───────┬──────────────┬───────┘
                               │              │
                     ───────►  ▼              ▼  ◄───────
                    Studio                      Apps
              (Studio Editor —          (Frontend Workspace —
           zasoby przy redagowaniu       zasoby i kod widoku
                dokumentów)             przy budowie interfejsu)
```

| Moduł docelowy | Charakter powiązania | Typ | Miejsce ustanowienia |
|---|---|---|---|
| Studio | Zasoby wygenerowane w Design wykorzystywane przy redagowaniu dokumentów | Konfiguracyjne | Assets Panel — „Wyślij do modułu” |
| Apps | Zasoby i kod widoku wykorzystywane przy budowie interfejsu produktu, w oknie Frontend Workspace | Konfiguracyjne | Assets Panel, Preview Window — „Wyślij do modułu” |
| Library | Archiwizacja zaakceptowanych zasobów i przewodnika stylu | Konfiguracyjne | Okno konfiguracji, klucz „Powiązanie z Library” |
| Automations | Kolejkowanie zadań wsadowych modułu | Konfiguracyjne | Okno konfiguracji, klucz „Wsadowe zadania przez Automations” |

Powiązania nie są aktywne domyślnie. Przekazanie zasobu do innego modułu odbywa się jako jawna, pojedyncza decyzja użytkownika podejmowana z poziomu Assets Panel lub Preview Window albo — zgodnie z ustawieniem konfiguracyjnym — mechanizmem synchronizacji; ręczne potwierdzanie każdego przekazania jest ustawieniem konfiguracyjnym, nie wymogiem.

---

## 9. Scenariusze użycia

### 9.1 Budowa spójnego zestawu ilustracji brandingowych

1. Użytkownik otwiera moduł Design w środowisku WorkSpace.
2. W Chat Window zleca przygotowanie zestawu ilustracji; Koordynator rozkłada zlecenie na zadania widoczne w Execution Loop Window.
3. W Prompt Builderze Użytkownik definiuje styl bazowy (paleta granat i złoto, styl geometryczny) i zapisuje go jako szablon.
4. Wykonawca generuje kolejne ilustracje na podstawie szablonu; pętla wykonawcza prowadzi zadania generowania, usunięcia tła i eksportu, zachowując spójność stylistyczną zestawu.
5. Zaakceptowane warianty trafiają do Assets Panel, oznaczone wspólnym tagiem kolekcji „Branding — kampania Q3”.
6. Na Design Board Użytkownik zestawia całość w jedną tablicę prezentacyjną i eksportuje ją jako PDF do przeglądu wewnętrznego.

### 9.2 Projektowanie interfejsu produktu równolegle z modułem Apps

1. Zespół pracujący nad produktem w module Apps (środowisko CodeStudio) potrzebuje zestawu elementów interfejsu.
2. W Tokens & System Panel projektant ustala tokeny kolorów, typografii i odstępów oraz sprawdza kontrast par barw.
3. W module Design projektant generuje ikony i elementy graficzne interfejsu, oceniając każdy wariant w Preview Window.
4. Zaakceptowane elementy oraz tokeny wydane jako zmienne CSS przekazywane są do okna Frontend Workspace modułu Apps przyciskiem „Wyślij do modułu”.
5. Zespół warstwy klienckiej w module Apps korzysta z przekazanych zasobów bez opuszczania własnej przestrzeni roboczej.

### 9.3 Iteracyjna korekta ilustracji na podstawie uwag

1. Wygenerowana ilustracja trafia do oceny zespołu w Design Board.
2. Uwagi nanoszone są jako adnotacje przypięte do konkretnych fragmentów kompozycji.
3. Na podstawie zebranych adnotacji Użytkownik formułuje w Chat Window polecenie modyfikujące („ciemniejsze tło, usuń element w lewym rogu kadru”).
4. Execution Loop Window prezentuje zadanie korekty, wynik kontroli jakości i decyzję o ponowieniu.
5. Nowy wariant trafia do Preview Window z aktywnym suwakiem porównawczym przed/po, ułatwiającym ocenę zmiany względem poprzedniej wersji.

---

## 10. Zgodność z rdzeniem platformy

| Zasada nadrzędna | Realizacja w module |
|---|---|
| Dwa kanały komunikacji | Chat Window (Użytkownik ↔ Wykonawca) jako centralny punkt pracy i podstawowy mechanizm sterowania; Execution Loop Window (Koordynator ↔ Wykonawca) jako okno pętli wykonawczej zadań modułu |
| Układ pionowy | Chat Window w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej, Design Board w prawej kolumnie dominującej, panele jako rozszerzenia boczne |
| Warstwy widoczności | Cztery warstwy ujawniania funkcji (rozdz. 2.1); w spoczynku widoczna warstwa 1 oraz zwinięte wyzwalacze |
| Zero blokad | „⚡ Generuj” zawsze aktywny; usuwanie z „Cofnij”; decyzje w Preview Window nieblokujące; brak okien nakładkowych w roli bramek |
| Klucze jawne | Klucze modeli i integracji („Klucze API modeli i bibliotek zewnętrznych”) wprowadzane wprost przez Operatora, widoczne w oknie konfiguracji |
| Pełna konfigurowalność | Siedemnaście punktów sterowania (rozdz. 7) z warstwowością i objaśnieniami `[?]` |
| Jawność zależności | Powiązania z Studio, Apps, Library i Automations („Powiązanie z modułem Studio”–„Wsadowe zadania przez Automations”) jako świadome ustawienia, nie ukryte reguły |
| Konfigurowalność modeli | Wybór silnika i punktu końcowego („Domyślny silnik generujący”, „Punkty końcowe GPU zadań neuronowych”) zgodny ze strategią modeli platformy |
| Audytowalność | Prowenancja zasobu („Metadane pochodzenia (prowenancja)”) — model, prompt, seed, łańcuch edycji — zgodna z blokiem prowenancji klienta |

---

## 11. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Design Board otwiera się z bocznej nawigacji WorkSpace i CodeStudio, nie z TalkIn | Próba odnalezienia modułu Design w bocznej nawigacji trzech środowisk |
| Żaden przycisk akcji modułu nie występuje w stanie zablokowanym | Przegląd kontrolek CTA katalogu elementów rozdz. 5 pod kątem obecności `disabled`/`not-allowed` |
| Przycisk „⚡ Generuj” pozostaje klikalny przy pustym poleceniu i sygnalizuje komunikatem, nie blokadą | Kliknięcie „⚡ Generuj” z pustym polem polecenia w Prompt Builderze |
| Wszystkie 106 komend obszaru `design` mają pokrycie w interfejsie modułu albo jawne wskazanie miejsca wywołania | Zestawienie Załącznika z katalogiem elementów rozdz. 5 |
| Przekazanie zasobu do modułu docelowego jest jawną, pojedynczą decyzją, nie powiązaniem aktywnym domyślnie | Test: nowo wygenerowany zasób nie trafia automatycznie do Studio ani Apps bez wybrania „→ Moduł” |
| Suwak porównawczy przed/po w Preview Window działa na każdym wygenerowanym wariancie | Otwarcie Preview Window z co najmniej dwoma wariantami zasobu |
| Stany kontrolek z katalogu elementów rozdz. 5 są zaimplementowane w komplecie | Przegląd katalogu elementów pod kątem kompletu stanów: domyślny, wskazanie kursorem, ognisko, ładowanie, ostrzeżenie, błąd |

---

## Załącznik — pełny wykaz komend kontraktu modułu Design

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `design` — 106 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `design.annotation.list` | Zwraca adnotacje kompozycji wraz z wątkami | `boardId:string` (wym)<br>`openOnly:bool` (opc) | `annotations:DesignAnnotation[]` (wym)<br>`total:int` (wym) |
| `design.annotation.set` | Zakłada albo zmienia adnotację przypiętą do warstwy kompozycji. Pole `note` warstwy niesie **jedno** zdanie bez autora i bez wątku; opracowanie wymaga wątków i oznaczeń osób, a tego jedno pole nie unosi | `boardId:string` (wym)<br>`layerId:string` (opc)<br>`annotationId:string` (opc)<br>`parentId:string` (opc)<br>`text:string` (wym)<br>`resolved:bool` (opc) | `annotation:DesignAnnotation` (wym) |
| `design.asset.content.get` | Oddaje **treść** zasobu z magazynu rdzenia. Pole `uri` zasobu jest ścieżką w systemie plików rdzenia (magazyn oddaje `filepath`), więc przeglądarka nie wczyta spod niego niczego — także wtedy, gdy zasób powstał bez zarzutu. Komenda dotyczy **każdego** zasobu magazynu, nie tylko obrazu: jeden magazyn obsługuje rodziny `design.*`, `document.*`, `media.*` i `archive.*`. Rdzeń odmawia zasobu, którego nie zna, i zasobu, którego treści nie ma pod sumą kontrolną — nigdy nie oddaje bajtów zastępczych | `assetId:string` (wym)<br>`disposition:AssetContentDisposition` (opc)<br>`maxBytes:int` (opc) | `assetId:string` (wym)<br>`contentBase64:string` (opc)<br>`mediaType:string` (wym)<br>`sizeBytes:int` (wym)<br>`checksum:string` (wym) |
| `design.asset.decision.list` | Zwraca rozstrzygnięcia zasobów okna. Bez odczytu Preview Window po odświeżeniu strony nie wiedziałoby, który wariant został już przyjęty | `windowId:string` (wym)<br>`assetId:string` (opc) | `records:DesignAssetDecisionRecord[]` (wym)<br>`total:int` (wym) |
| `design.asset.decision.set` | Zapisuje rozstrzygnięcie Operatora wobec zasobu: przyjęty, odrzucony, przekazany albo bez rozstrzygnięcia. Preview Window ma trzy akcje decyzyjne, a magazyn nie miał gdzie zapisać ich wyniku — oznaczenie ulubionego znaczy co innego i służy do czego innego | `assetId:string` (wym)<br>`decision:DesignAssetDecision` (wym)<br>`note:string` (opc) | `record:DesignAssetDecisionRecord` (wym) |
| `design.asset.export` | Wydaje zasób Operatorowi jako plik we wskazanym formacie i skali. Różni się od `image.convert` celem: konwersja zakłada **nowy** zasób w magazynie i tam się kończy, eksport oddaje bajty gotowe do zapisania poza produktem | `assetId:string` (wym)<br>`format:string` (wym)<br>`scale:float` (opc)<br>`quality:int` (opc)<br>`stripMetadata:bool` (opc)<br>`colorProfile:string` (opc) | `contentBase64:string` (wym)<br>`fileName:string` (wym)<br>`mediaType:string` (wym)<br>`sizeBytes:int` (wym) |
| `design.asset.export.batch` | Wydaje wiele zasobów naraz, każdy w komplecie wskazanych skal. Pokrywa eksport zbiorczy Assets Panel oraz zestawy rozmiarów kampanii; odmowa jednego zasobu **nie** wstrzymuje pozostałych — wynik niesie bilans przyjętych i odrzuconych | `assetIds:string[]` (wym)<br>`format:string` (wym)<br>`scales:float[]` (opc)<br>`quality:int` (opc)<br>`archive:bool` (opc) | `files:DesignExportFile[]` (wym)<br>`exported:int` (wym)<br>`rejected:DesignExportRejection[]` (opc) |
| `design.asset.favorite.set` | Oznacza zasób jako ulubiony albo zdejmuje oznaczenie. Komenda `design.asset.list` umiała już filtrować po ulubionych, ale nikt nie miał jak ich ustawić | `assetId:string` (wym)<br>`favorite:bool` (wym) | `asset:DesignAsset` (wym) |
| `design.asset.generate` | Generuje zasób wizualny z promptu strukturalnego kanałem modelu obrazowego. Treść trafia do magazynu rdzenia pod sumą kontrolną, więc zasób nie zależy od żadnego pliku zewnętrznego. Brak kanału obrazowego albo brak poświadczenia to **odmowa nazywająca brak** — rdzeń nigdy nie zakłada zasobu bez bajtów obrazu | `windowId:string` (wym)<br>`prompt:DesignPrompt` (wym)<br>`referenceAssetId:string` (opc)<br>`channelId:string` (opc)<br>`kind:DesignAssetKind` (opc) | `assets:DesignAsset[]` (wym)<br>`processId:string` (opc) |
| `design.asset.list` | Zwraca zasoby Assets Panel | `windowId:string` (opc)<br>`kind:DesignAssetKind` (opc)<br>`tags:string[]` (opc)<br>`favoriteOnly:bool` (opc)<br>`limit:int` (opc) | `assets:DesignAsset[]` (wym)<br>`total:int` (opc) |
| `design.asset.remove` | Usuwa zasób z Assets Panel. Rodzaj zmiany `deleted` zdarzenia `design.asset.changed` nie miał dotąd żadnego nadawcy | `assetId:string` (wym) | `removed:bool` (wym) |
| `design.asset.tag.set` | Ustawia etykiety zasobu; `design.asset.list` zawężał po nich, choć nadać ich nie było czym | `assetId:string` (wym)<br>`tags:string[]` (wym) | `asset:DesignAsset` (wym) |
| `design.asset.upload` | Wnosi zasób wizualny wskazany przez Operatora do Assets Panel. Treść wchodzi do magazynu rdzenia pod sumą kontrolną — zasób przestaje zależeć od pliku, który Operator może nadpisać albo skasować. Komenda dotyczy zasobu **już istniejącego**: Operator ma plik albo wskazuje odsyłacz. Zasób mający dopiero powstać zamawia się komendą `design.asset.generate` | `windowId:string` (wym)<br>`name:string` (opc)<br>`kind:DesignAssetKind` (wym)<br>`contentBase64:string` (opc)<br>`sourcePath:string` (opc)<br>`format:string` (opc)<br>`width:int` (opc)<br>`height:int` (opc)<br>`tags:string[]` (opc) | `asset:DesignAsset` (wym) |
| `design.board.export` | Wyrysowuje całą kompozycję albo wskazany obszar do jednego pliku. Dziś kompozycja jeździ do rdzenia i z powrotem jako układ warstw i nie ma drogi wyjścia poza rdzeń | `boardId:string` (wym)<br>`format:string` (wym)<br>`region:DesignBoardRegion` (opc)<br>`scale:float` (opc) | `contentBase64:string` (wym)<br>`fileName:string` (wym)<br>`mediaType:string` (wym) |
| `design.board.list` | Zwraca kompozycje okna wraz z warstwami. Komenda `design.board.update` jechała w jedną stronę — po odświeżeniu okna plansza nie wracała | `windowId:string` (wym) | `boards:DesignBoard[]` (wym)<br>`total:int` (wym) |
| `design.board.update` | Zapisuje układ kompozycji Design Board | `windowId:string` (wym)<br>`boardId:string` (opc)<br>`name:string` (opc)<br>`layers:DesignBoardLayer[]` (opc) | `board:DesignBoard` (wym) |
| `design.board.version.list` | Zwraca wersje kompozycji — sam ciąg postaci układu, bez warstw, żeby wykaz nie ważył tyle co cały zapis | `boardId:string` (wym)<br>`limit:int` (opc) | `versions:DesignBoardVersion[]` (wym)<br>`total:int` (wym) |
| `design.board.version.restore` | Przywraca układ kompozycji z wersji. Przywrócenie **zakłada** nową wersję z układu sprzed przywrócenia, żeby cofnięcie się samo nie kasowało stanu, który Operator właśnie porzucił | `versionId:string` (wym) | `board:DesignBoard` (wym)<br>`supersededVersion:DesignBoardVersion` (opc) |
| `design.board.version.save` | Utrwala bieżący układ kompozycji jako nazwaną wersję. Komenda `design.board.update` zapisuje układ **bieżący** i zastępuje poprzedni, więc ciągu postaci tablicy nie ma dziś skąd wziąć | `boardId:string` (wym)<br>`name:string` (opc)<br>`note:string` (opc) | `version:DesignBoardVersion` (wym) |
| `design.campaign.set.build` | Wydaje komplet materiałów kampanii w każdym wskazanym rozmiarze naraz | `boardId:string` (wym)<br>`sizes:DesignCampaignSize[]` (wym)<br>`format:string` (opc)<br>`windowId:string` (opc) | `sizes:DesignCampaignSize[]` (wym)<br>`exported:int` (wym)<br>`failedSizes:string[]` (opc) |
| `design.chart.render` | Wyrysowuje wykres z danych jako zasób magazynu. Dane przychodzą seriami, a nie obrazem — wykres da się więc przerysować po zmianie liczb, zamiast rysować go od nowa | `kind:DesignChartKind` (wym)<br>`series:DesignChartSeries[]` (wym)<br>`categories:string[]` (opc)<br>`title:string` (opc)<br>`width:float` (opc)<br>`height:float` (opc)<br>`format:string` (opc)<br>`tokenSetId:string` (opc)<br>`legend:bool` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`mediaType:string` (wym) |
| `design.collection.assign` | Przypisuje zasoby do kolekcji albo je z niej zdejmuje. Zestaw jest **dokładką** albo **odjęciem**, nie zastąpieniem — inaczej niż przy etykietach, bo kolekcja bywa duża i przepisywanie jej w całości przy każdej zmianie jest drogą do zgubienia zawartości | `collectionId:string` (wym)<br>`assetIds:string[]` (wym)<br>`remove:bool` (opc) | `collection:DesignCollection` (wym)<br>`changed:int` (wym) |
| `design.collection.create` | Zakłada kolekcję zasobów. Etykiety zasobu już są (`design.asset.tag.set`), ale kolekcja jest bytem osobnym: ma nazwę, opis i porządek, a etykieta jest tylko słowem | `windowId:string` (wym)<br>`name:string` (wym)<br>`description:string` (opc) | `collection:DesignCollection` (wym) |
| `design.collection.list` | Zwraca kolekcje okna wraz z licznikiem zasobów. Bez tego kolekcja założona nie miałaby jak wrócić na ekran po odświeżeniu | `windowId:string` (wym)<br>`assetId:string` (opc) | `collections:DesignCollection[]` (wym)<br>`total:int` (wym) |
| `design.color.accessibility.audit` | Sprawdza całe zestawienie żetonów i wskazuje pary łamiące kontrast. Czynność **czyta** i niczego nie zmienia | `tokenSetId:string` (wym)<br>`minimumRatio:float` (opc) | `violations:DesignContrastResult[]` (wym)<br>`checked:int` (wym)<br>`passed:int` (wym) |
| `design.color.contrast.check` | Liczy współczynnik kontrastu pary barw i ocenę wedle WCAG | `foreground:string` (wym)<br>`background:string` (wym)<br>`fontSize:float` (opc)<br>`bold:bool` (opc) | `result:DesignContrastResult` (wym) |
| `design.color.convert` | Przelicza barwę między przestrzeniami i oddaje ją we wszystkich naraz | `value:string` (wym)<br>`space:DesignColorSpace` (opc) | `color:DesignColorValue` (wym) |
| `design.color.gradient.set` | Zakłada albo zmienia gradient wielostopniowy z interpolacją w przestrzeni percepcyjnej | `boardId:string` (wym)<br>`gradient:DesignGradient` (wym)<br>`pathId:string` (opc)<br>`layerId:string` (opc) | `gradient:DesignGradient` (wym)<br>`previewSvg:string` (opc) |
| `design.color.palette.extract` | Wyciąga barwy dominujące z zasobu obrazowego wraz z ich udziałem | `assetId:string` (wym)<br>`count:int` (opc) | `colors:DesignPaletteColor[]` (wym)<br>`pixelsSampled:int` (opc) |
| `design.color.palette.generate` | Tworzy paletę z reguły harmonii wobec barwy wiodącej | `baseColor:string` (wym)<br>`harmony:DesignColorHarmony` (wym)<br>`count:int` (opc)<br>`windowId:string` (opc) | `colors:DesignPaletteColor[]` (wym)<br>`baseColor:string` (wym) |
| `design.color.vision.simulate` | Oddaje zasób obrazowy w symulacji wady widzenia barw. Wynik jest nowym zasobem — źródło zostaje nietknięte | `assetId:string` (wym)<br>`vision:DesignColorVision` (wym)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym) |
| `design.component.instance.add` | Wstawia instancję komponentu do ramki albo na kanwę | `componentId:string` (wym)<br>`boardId:string` (wym)<br>`variant:string` (opc)<br>`frameId:string` (opc)<br>`x:float` (opc)<br>`y:float` (opc) | `layer:DesignBoardLayer` (wym)<br>`component:DesignComponent` (wym) |
| `design.component.list` | Zwraca komponenty okna wraz z wariantami i liczbą instancji | `windowId:string` (wym)<br>`componentId:string` (opc) | `components:DesignComponent[]` (wym)<br>`total:int` (wym) |
| `design.component.save` | Utrwala komponent wielokrotnego użycia wraz z wariantami stanu | `windowId:string` (wym)<br>`name:string` (wym)<br>`componentId:string` (opc)<br>`variants:DesignComponentVariant[]` (opc)<br>`tokenSetId:string` (opc) | `component:DesignComponent` (wym)<br>`propagatedTo:int` (opc) |
| `design.constraint.set` | Ustawia więzy responsywne warstw wobec ramki | `frameId:string` (wym)<br>`constraints:DesignConstraint[]` (wym) | `constraints:DesignConstraint[]` (wym)<br>`changed:int` (wym) |
| `design.diagram.render` | Wyrysowuje schemat z węzłów i połączeń jako zasób magazynu | `kind:DesignDiagramKind` (wym)<br>`nodes:DesignDiagramNode[]` (wym)<br>`edges:DesignDiagramEdge[]` (opc)<br>`title:string` (opc)<br>`direction:DesignLayoutDirection` (opc)<br>`format:string` (opc)<br>`tokenSetId:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`mediaType:string` (wym)<br>`unplacedNodeIds:string[]` (opc) |
| `design.favicon.build` | Wydaje komplet ikon aplikacji ze źródła wraz z manifestem | `assetId:string` (wym)<br>`sizes:int[]` (opc)<br>`includeManifest:bool` (opc)<br>`windowId:string` (opc) | `assetIds:string[]` (wym)<br>`manifest:string` (opc)<br>`sizes:int[]` (wym) |
| `design.font.glyphs.get` | Zwraca znaki kroju wraz z ich miarą i konturem | `fontFamily:string` (wym)<br>`from:int` (opc)<br>`to:int` (opc)<br>`limit:int` (opc) | `glyphs:DesignGlyph[]` (wym)<br>`total:int` (wym) |
| `design.font.pair.suggest` | Proponuje zestawienia krojów nagłówka i tekstu wraz z powodem | `windowId:string` (wym)<br>`mood:string` (opc)<br>`baseFont:string` (opc)<br>`count:int` (opc)<br>`channelId:string` (opc) | `pairs:DesignFontPair[]` (wym)<br>`source:string` (wym) |
| `design.font.preview` | Oddaje podgląd kroju w tekście próbnym wraz z gotową regułą osadzenia | `fontFamily:string` (wym)<br>`sampleText:string` (opc)<br>`sizes:float[]` (opc)<br>`subsetText:string` (opc) | `previewSvg:string` (wym)<br>`fontFace:string` (opc)<br>`available:bool` (wym)<br>`glyphCount:int` (opc) |
| `design.frame.list` | Zwraca ramki kompozycji wraz z ich układem i siatką | `boardId:string` (wym) | `frames:DesignFrame[]` (wym)<br>`total:int` (wym)<br>`devicePresets:string[]` (opc) |
| `design.frame.remove` | Usuwa ramkę kompozycji. Warstwy ramki zostają na kanwie — usunięcie ramki nie jest usunięciem pracy, która w niej leżała | `frameId:string` (wym) | `removed:bool` (wym)<br>`releasedLayerIds:string[]` (opc) |
| `design.frame.resize.apply` | Zmienia rozmiar ramki wraz z zastosowaniem więzów responsywnych. Osobno od `design.frame.set`, bo tam zmiana rozmiaru jest zapisem nastawy, a tu przeliczeniem układu wedle więzów | `frameId:string` (wym)<br>`width:float` (wym)<br>`height:float` (wym) | `frame:DesignFrame` (wym)<br>`layers:DesignBoardLayer[]` (wym) |
| `design.frame.set` | Zakłada ramkę albo zmienia jej rozmiar, układ i siatkę. Ramka jest ekranem makiety: to ona, a nie kanwa, wyznacza obszar wydania | `boardId:string` (wym)<br>`name:string` (wym)<br>`frameId:string` (opc)<br>`width:float` (opc)<br>`height:float` (opc)<br>`x:float` (opc)<br>`y:float` (opc)<br>`devicePreset:string` (opc)<br>`grid:DesignGrid` (opc)<br>`layout:DesignAutoLayout` (opc) | `frame:DesignFrame` (wym) |
| `design.grid.set` | Ustawia siatkę układu ramki albo całej kompozycji | `grid:DesignGrid` (wym)<br>`frameId:string` (opc)<br>`boardId:string` (opc) | `grid:DesignGrid` (wym) |
| `design.icon.generate` | Tworzy komplet ikon w jednym stylu z wykazu pojęć, z wymuszeniem siatki i grubości obrysu | `windowId:string` (wym)<br>`concepts:string[]` (wym)<br>`gridSize:int` (opc)<br>`strokeWidth:float` (opc)<br>`styleReferenceIconId:string` (opc)<br>`channelId:string` (opc) | `icons:DesignIcon[]` (wym)<br>`failedConcepts:string[]` (opc) |
| `design.icon.library.search` | Przeszukuje katalog ikon otwartoźródłowych wniesiony do rdzenia | `query:string` (opc)<br>`set:string` (opc)<br>`limit:int` (opc) | `icons:DesignIcon[]` (wym)<br>`total:int` (wym)<br>`sets:string[]` (opc) |
| `design.icon.set` | Zakłada albo zmienia ikonę własną na siatce z wyrównaniem do pikseli | `windowId:string` (wym)<br>`name:string` (wym)<br>`svg:string` (wym)<br>`iconId:string` (opc)<br>`gridSize:int` (opc)<br>`strokeWidth:float` (opc)<br>`tags:string[]` (opc) | `icon:DesignIcon` (wym)<br>`gridWarnings:string[]` (opc) |
| `design.icon.sprite.build` | Pakuje ikony w sprite SVG z symbolami albo w font ikon | `iconIds:string[]` (wym)<br>`kind:DesignIconSpriteKind` (wym)<br>`name:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`fileName:string` (wym)<br>`included:int` (wym) |
| `design.largeformat.tile` | Dzieli materiał wielkoformatowy na kafle o rozmiarze nośnika, z zakładką na sklejenie. Odpowiedź niesie układ kafli, więc Operator wie, który kafel gdzie idzie, zanim cokolwiek wydrukuje | `assetId:string` (wym)<br>`tileWidthMm:float` (wym)<br>`tileHeightMm:float` (wym)<br>`overlapMm:float` (opc)<br>`targetWidthMm:float` (opc)<br>`targetHeightMm:float` (opc)<br>`markers:bool` (opc)<br>`windowId:string` (opc) | `tiles:DesignTile[]` (wym)<br>`rows:int` (wym)<br>`columns:int` (wym)<br>`effectiveDpi:int` (opc) |
| `design.layout.auto` | Przelicza układ automatyczny ramki i przestawia położenia warstw. Rdzeń oddaje warstwy po przeliczeniu, więc klient nie liczy układu drugi raz i nie ma jak się rozjechać | `frameId:string` (wym)<br>`layout:DesignAutoLayout` (wym)<br>`layerIds:string[]` (opc) | `layers:DesignBoardLayer[]` (wym)<br>`contentWidth:float` (opc)<br>`contentHeight:float` (opc) |
| `design.mockup.generate` | Buduje szkielet ekranu z polecenia tekstowego: ramka, komponenty i układ. Wynik jest makietą do poprawienia, nie projektem końcowym, i tak się o nim mówi | `windowId:string` (wym)<br>`boardId:string` (wym)<br>`prompt:string` (wym)<br>`devicePreset:string` (opc)<br>`tokenSetId:string` (opc)<br>`channelId:string` (opc) | `frame:DesignFrame` (wym)<br>`layers:DesignBoardLayer[]` (wym)<br>`componentIds:string[]` (opc) |
| `design.mockup.import` | Odtwarza edytowalną makietę ze zrzutu ekranu: rozpoznaje układ i tekst, składa warstwy. Odpowiedź mówi, czego rdzeń nie rozpoznał — bilans zamiast ciszy | `windowId:string` (wym)<br>`boardId:string` (wym)<br>`assetId:string` (wym)<br>`recognizeText:bool` (opc)<br>`devicePreset:string` (opc) | `frame:DesignFrame` (wym)<br>`layers:DesignBoardLayer[]` (wym)<br>`unrecognizedRegions:int` (opc) |
| `design.photo.background.remove` | Odcina tło i zostawia kanał krycia. Gdy kanał modelu obrazowego stoi, liczy kanał; gdy nie stoi, liczy rachunek wkompilowany — odpowiedź mówi to polem `computedBy` | `assetId:string` (wym)<br>`tolerance:float` (opc)<br>`channelId:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`computedBy:DesignPhotoComputeRoute` (wym)<br>`hasAlpha:bool` (wym)<br>`transparentShare:float` (wym) |
| `design.photo.batch.apply` | Powtarza ten sam zestaw czynności na wielu zasobach | `assetIds:string[]` (wym)<br>`operations:DesignPhotoOperation[]` (opc)<br>`presetId:string` (opc)<br>`windowId:string` (opc) | `assets:DesignAsset[]` (wym)<br>`applied:int` (wym)<br>`failedAssetIds:string[]` (opc) |
| `design.photo.color.correct` | Koryguje barwę zdjęcia: jasność, kontrast, ekspozycja, temperatura barwowa, gamma, nasycenie, krzywe i HSL | `assetId:string` (wym)<br>`brightness:float` (opc)<br>`contrast:float` (opc)<br>`exposure:float` (opc)<br>`temperature:float` (opc)<br>`gamma:float` (opc)<br>`saturation:float` (opc)<br>`hueShiftDeg:float` (opc)<br>`lightness:float` (opc)<br>`curves:json` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`histogramShift:float` (wym) |
| `design.photo.crop` | Kadruje zdjęcie, prostuje horyzont i sprowadza do proporcji. Wynik jest wariantem źródła — oryginał zostaje nietknięty | `assetId:string` (wym)<br>`x:float` (opc)<br>`y:float` (opc)<br>`width:float` (opc)<br>`height:float` (opc)<br>`aspectRatio:string` (opc)<br>`straightenDeg:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`width:int` (wym)<br>`height:int` (wym) |
| `design.photo.enhance` | Poprawia jakość zdjęcia: auto-poziomy, auto-kontrast, odszumienie i wyostrzenie | `assetId:string` (wym)<br>`autoLevels:bool` (opc)<br>`autoContrast:bool` (opc)<br>`denoise:float` (opc)<br>`sharpen:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`appliedSteps:string[]` (wym) |
| `design.photo.expand` | Rozszerza kadr poza pierwotną ramkę. Gdy kanał modelu obrazowego stoi, liczy kanał; gdy nie stoi, liczy rachunek wkompilowany — odpowiedź mówi to polem `computedBy` | `assetId:string` (wym)<br>`left:int` (opc)<br>`right:int` (opc)<br>`top:int` (opc)<br>`bottom:int` (opc)<br>`prompt:string` (opc)<br>`channelId:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`computedBy:DesignPhotoComputeRoute` (wym)<br>`width:int` (wym)<br>`height:int` (wym) |
| `design.photo.filter.apply` | Nakłada filtr obrazu: rozmycie, wyostrzenie, ziarno, winieta, sepia, monochrom, poświata albo cień | `assetId:string` (wym)<br>`filter:DesignPhotoFilter` (wym)<br>`amount:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym) |
| `design.photo.history.get` | Oddaje łańcuch edycji zasobu: czynności wraz z nastawami i drogą rachunku | `assetId:string` (wym)<br>`limit:int` (opc) | `edits:DesignPhotoEdit[]` (wym)<br>`total:int` (wym) |
| `design.photo.inpaint` | Domalowuje obszar z maski. Gdy kanał modelu obrazowego stoi, liczy kanał; gdy nie stoi, liczy rachunek wkompilowany — odpowiedź mówi to polem `computedBy` | `assetId:string` (wym)<br>`maskAssetId:string` (opc)<br>`regions:DesignPhotoRegion[]` (opc)<br>`prompt:string` (opc)<br>`channelId:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`computedBy:DesignPhotoComputeRoute` (wym) |
| `design.photo.layer.composite` | Składa warstwy rastrowe w jeden obraz trybami mieszania i kryciem | `layers:DesignPhotoLayer[]` (wym)<br>`width:int` (opc)<br>`height:int` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`composited:int` (wym)<br>`skippedLayerAssetIds:string[]` (opc) |
| `design.photo.mask.set` | Zakłada maskę nieniszczącą na zasobie: źródło zostaje nietknięte, a maska zapisuje się jako czynność łańcucha edycji | `assetId:string` (wym)<br>`maskAssetId:string` (wym)<br>`invert:bool` (opc)<br>`featherPx:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`maskAssetId:string` (wym) |
| `design.photo.metadata.get` | Oddaje zmierzone właściwości zasobu obrazowego wraz z danymi EXIF | `assetId:string` (wym) | `metadata:DesignPhotoMetadata` (wym) |
| `design.photo.preset.list` | Zwraca nastawy warsztatu fotografii w oknie | `windowId:string` (wym) | `presets:DesignPhotoPreset[]` (wym)<br>`total:int` (wym) |
| `design.photo.preset.save` | Zapisuje zestaw czynności warsztatu fotografii pod nazwą | `windowId:string` (wym)<br>`name:string` (wym)<br>`operations:DesignPhotoOperation[]` (wym)<br>`presetId:string` (opc) | `preset:DesignPhotoPreset` (wym) |
| `design.photo.resample` | Przelicza rozdzielczość zdjęcia wskazanym filtrem | `assetId:string` (wym)<br>`width:int` (opc)<br>`height:int` (opc)<br>`filter:DesignPhotoResampleFilter` (opc)<br>`keepAspectRatio:bool` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`width:int` (wym)<br>`height:int` (wym) |
| `design.photo.retouch` | Retusz obszarów: klonowanie ze wskazanego źródła albo leczenie z otoczenia | `assetId:string` (wym)<br>`regions:DesignPhotoRegion[]` (wym)<br>`mode:DesignPhotoRetouchMode` (opc)<br>`source:DesignPhotoPoint` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`regionsApplied:int` (wym)<br>`regionsSkipped:string[]` (opc) |
| `design.photo.select.object` | Zaznacza obiekt i oddaje maskę jako zasób | `assetId:string` (wym)<br>`point:DesignPhotoPoint` (wym)<br>`tolerance:float` (opc)<br>`windowId:string` (opc) | `mask:DesignAsset` (wym)<br>`coverage:float` (wym)<br>`bounds:DesignPhotoRegion` (opc) |
| `design.photo.transform` | Obraca, odbija i koryguje perspektywę oraz zniekształcenia obiektywu | `assetId:string` (wym)<br>`rotateDeg:float` (opc)<br>`flipHorizontal:bool` (opc)<br>`flipVertical:bool` (opc)<br>`perspective:DesignPhotoPoint[]` (opc)<br>`lensDistortion:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`width:int` (wym)<br>`height:int` (wym) |
| `design.photo.upscale` | Powiększa zdjęcie z wyostrzeniem po powiększeniu. Gdy kanał modelu obrazowego stoi, liczy kanał; gdy nie stoi, liczy rachunek wkompilowany — odpowiedź mówi to polem `computedBy` | `assetId:string` (wym)<br>`factor:int` (wym)<br>`sharpenAfter:bool` (opc)<br>`channelId:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`computedBy:DesignPhotoComputeRoute` (wym)<br>`width:int` (wym)<br>`height:int` (wym) |
| `design.photo.vectorize` | Zamienia raster w rysunek wektorowy przez progowanie i obrysowanie konturów rachunkiem wkompilowanym | `assetId:string` (wym)<br>`colors:int` (opc)<br>`threshold:float` (opc)<br>`smoothing:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`pathCount:int` (wym)<br>`droppedRegions:int` (opc) |
| `design.presence.report` | Zgłasza obecność i położenie kursora Operatora na kompozycji. Rdzeń rozgłasza je pozostałym zdarzeniem `design.board.presence`. Zgłoszenie jest **ulotne** — nie zapisuje się w bazie, bo położenie kursora sprzed godziny nie jest wiedzą o niczym | `boardId:string` (wym)<br>`x:float` (opc)<br>`y:float` (opc)<br>`selectedLayerIds:string[]` (opc)<br>`leaving:bool` (opc) | `participants:DesignPresence[]` (wym) |
| `design.print.export` | Wydaje materiał gotowy do druku: przestrzeń barw, spady, znaczniki cięcia i pasowania, osadzony profil ICC. Wydanie idzie przez kontrolę przeddrukową i **odmawia**, gdy kontrola znajdzie wadę o wadze błędu — plik nie do druku wydany jako gotowy do druku jest gorszy niż odmowa | `format:string` (wym)<br>`boardId:string` (opc)<br>`assetId:string` (opc)<br>`frameId:string` (opc)<br>`profileId:string` (opc)<br>`profile:DesignPrintProfile` (opc)<br>`skipPreflight:bool` (opc)<br>`windowId:string` (opc)<br>`templateId:string` (opc)<br>`pageOrder:int[]` (opc)<br>`binding:DesignPrintBinding` (opc) | `asset:DesignAsset` (wym)<br>`fileName:string` (wym)<br>`mediaType:string` (wym)<br>`sizeBytes:int` (wym)<br>`preflightSkipped:bool` (wym)<br>`issues:DesignPreflightIssue[]` (opc)<br>`pageCount:int` (opc) |
| `design.print.paper.list` | Zwraca nośniki znane rdzeniowi wraz z ich wymiarami | `family:string` (opc) | `sizes:DesignPaperSize[]` (wym)<br>`total:int` (wym) |
| `design.print.preflight` | Sprawdza materiał przed drukiem: rozdzielczość obrazów, spady, przestrzeń barw, osadzenie krojów, cienkie linie i drobny tekst. Czynność **czyta** i niczego nie zmienia; zastrzeżenie o wadze błędu zatrzyma druk, więc Operator ma je zobaczyć **tutaj**, a nie w drukarni | `boardId:string` (opc)<br>`assetId:string` (opc)<br>`profileId:string` (opc)<br>`profile:DesignPrintProfile` (opc) | `issues:DesignPreflightIssue[]` (wym)<br>`errors:int` (wym)<br>`warnings:int` (wym)<br>`ready:bool` (wym) |
| `design.print.profile.list` | Zwraca profile wydania do druku wraz z nośnikami znanymi rdzeniowi | `windowId:string` (wym) | `profiles:DesignPrintProfile[]` (wym)<br>`paperSizes:DesignPaperSize[]` (wym)<br>`iccProfiles:string[]` (opc) |
| `design.print.profile.set` | Zakłada albo zmienia profil wydania do druku: spady, znaczniki, przestrzeń barw, rozdzielczość i profil ICC | `windowId:string` (wym)<br>`profile:DesignPrintProfile` (wym)<br>`profileId:string` (opc) | `profile:DesignPrintProfile` (wym) |
| `design.product.mockup.render` | Nakłada projekt na zdjęcie produktu wedle obszaru nałożenia i oddaje makietę produktową | `designAssetId:string` (wym)<br>`productAssetId:string` (wym)<br>`x:float` (wym)<br>`y:float` (wym)<br>`width:float` (wym)<br>`height:float` (wym)<br>`rotation:float` (opc)<br>`opacity:float` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym) |
| `design.prompt.history.list` | Zwraca prompty wydane w tym oknie wraz z zasobami, które z nich powstały. Domyka też prowenancję: dziś zasób niesie pole `promptId`, którego rdzeń **nie** wypełnia, bo nie ma przekładu klucza wiersza promptu na kod kontraktu — ta komenda ten przekład wnosi | `windowId:string` (wym)<br>`limit:int` (opc) | `prompts:DesignPromptRecord[]` (wym)<br>`total:int` (wym) |
| `design.prompt.template.list` | Zwraca szablony promptów okna | `windowId:string` (wym) | `templates:DesignPromptTemplate[]` (wym)<br>`total:int` (wym) |
| `design.prompt.template.save` | Utrwala prompt strukturalny jako szablon do wielokrotnego użycia. Dziś historia promptów i szablony żyją w oknie do zamknięcia karty przeglądarki i tyle o nich wiadomo | `windowId:string` (wym)<br>`name:string` (wym)<br>`prompt:DesignPrompt` (wym)<br>`templateId:string` (opc) | `template:DesignPromptTemplate` (wym) |
| `design.prototype.get` | Zwraca graf prototypu: ramki wraz z połączeniami między nimi | `boardId:string` (wym)<br>`startFrameId:string` (opc) | `frames:DesignFrame[]` (wym)<br>`links:DesignPrototypeLink[]` (wym)<br>`unreachableFrameIds:string[]` (opc) |
| `design.prototype.link.remove` | Usuwa połączenie prototypu | `linkId:string` (wym) | `removed:bool` (wym) |
| `design.prototype.link.set` | Łączy dwa ekrany prototypu przejściem albo zmienia istniejące połączenie | `boardId:string` (wym)<br>`fromFrameId:string` (wym)<br>`toFrameId:string` (wym)<br>`trigger:DesignPrototypeTrigger` (wym)<br>`transition:DesignPrototypeTransition` (wym)<br>`linkId:string` (opc)<br>`durationMs:int` (opc)<br>`layerId:string` (opc) | `link:DesignPrototypeLink` (wym) |
| `design.stock.import` | Wciąga zasób z katalogu zewnętrznego do magazynu okna wraz z zapisem licencji | `provider:string` (wym)<br>`externalId:string` (wym)<br>`windowId:string` (wym) | `asset:DesignAsset` (wym)<br>`license:string` (opc) |
| `design.stock.search` | Przeszukuje katalogi zasobów zewnętrznych wskazane w ustawieniach okna. Brak skonfigurowanego dostawcy jest **odpowiedzią nazywającą brak**, nie ciszą | `query:string` (wym)<br>`windowId:string` (wym)<br>`provider:string` (opc)<br>`limit:int` (opc) | `assets:DesignStockAsset[]` (wym)<br>`providersQueried:string[]` (wym)<br>`providersFailed:string[]` (opc) |
| `design.styleguide.publish` | Wydaje przewodnik systemu projektowego do modułu docelowego. Przeglądarka składa przewodnik sama i oddaje go plikiem, ale wydania go do Library albo Studio nie ma czym zlecić | `tokenSetId:string` (wym)<br>`targetModuleId:string` (wym)<br>`collectionId:string` (opc) | `assetId:string` (wym)<br>`published:bool` (wym) |
| `design.template.apply` | Zakłada kompozycję z szablonu wraz z podstawieniem treści | `templateId:string` (wym)<br>`windowId:string` (wym)<br>`boardId:string` (opc)<br>`replacements:json` (opc)<br>`pageNumber:int` (opc) | `board:DesignBoard` (wym)<br>`unmatchedNames:string[]` (opc)<br>`pages:DesignTemplatePage[]` (opc) |
| `design.template.list` | Zwraca szablony okna | `windowId:string` (wym)<br>`kind:DesignTemplateKind` (opc) | `templates:DesignTemplate[]` (wym)<br>`total:int` (wym) |
| `design.template.save` | Utrwala szablon materiału wraz z jego warstwami | `windowId:string` (wym)<br>`name:string` (wym)<br>`kind:DesignTemplateKind` (wym)<br>`width:float` (wym)<br>`height:float` (wym)<br>`templateId:string` (opc)<br>`layers:DesignBoardLayer[]` (opc)<br>`description:string` (opc)<br>`pages:DesignTemplatePage[]` (opc) | `template:DesignTemplate` (wym) |
| `design.tokenset.export` | Wydaje zestaw żetonów w postaci przyjmowanej przez kod. Klient składa dziś zmienne CSS, SCSS, konfigurację Tailwind i moduł JavaScript sam i nie potrzebuje do tego rdzenia; ta komenda jest potrzebna dla postaci, których przeglądarka złożyć nie może, oraz dla wydania idącego **do innego modułu** zamiast do pliku | `tokenSetId:string` (wym)<br>`target:DesignTokenTarget` (wym)<br>`targetModuleId:string` (opc) | `content:string` (wym)<br>`fileName:string` (wym)<br>`delivered:bool` (opc) |
| `design.tokenset.import` | Wczytuje zestaw żetonów z zapisu zewnętrznego i zakłada z niego zestaw modułu. Rdzeń **nie** nadpisuje motywu produktu — motyw jest własnością powłoki; import zakłada byt obok niego i oddaje różnicę wobec żetonów wskazanego motywu | `windowId:string` (wym)<br>`name:string` (wym)<br>`contentBase64:string` (wym)<br>`format:string` (opc) | `tokenSet:DesignTokenSet` (wym)<br>`unknownNames:string[]` (opc) |
| `design.tokenset.list` | Zwraca zestawy żetonów okna wraz z ich żetonami | `windowId:string` (wym)<br>`tokenSetId:string` (opc) | `tokenSets:DesignTokenSet[]` (wym)<br>`total:int` (wym) |
| `design.tokenset.save` | Utrwala zestaw żetonów systemu projektowego. Dziś Tokens & System Panel czyta żetony z motywu obowiązującego i nie ma ich gdzie odłożyć — kontrakt nie zna bytu zestawu żetonów | `windowId:string` (wym)<br>`name:string` (wym)<br>`tokens:DesignToken[]` (wym)<br>`tokenSetId:string` (opc)<br>`theme:string` (opc) | `tokenSet:DesignTokenSet` (wym) |
| `design.vector.boolean` | Wykonuje operację logiczną na ścieżkach i oddaje wynik jako jedną ścieżkę. Ścieżki źródłowe znikają albo zostają wedle wskazania — bo suma dwóch kształtów bywa krokiem pośrednim, a bywa wynikiem końcowym | `pathIds:string[]` (wym)<br>`operation:DesignBooleanOp` (wym)<br>`keepSources:bool` (opc) | `path:DesignVectorPath` (wym)<br>`removedPathIds:string[]` (opc) |
| `design.vector.export` | Wydaje ścieżki wektorowe jako czysty SVG, PDF wektorowy albo EPS. Wydanie zostaje wektorem — rasteryzacja byłaby tutaj stratą, której nie da się cofnąć | `boardId:string` (wym)<br>`target:DesignVectorExportTarget` (wym)<br>`pathIds:string[]` (opc)<br>`frameId:string` (opc) | `contentBase64:string` (wym)<br>`fileName:string` (wym)<br>`mediaType:string` (wym)<br>`sizeBytes:int` (wym) |
| `design.vector.optimize` | Czyści i skraca zapis SVG bez zmiany obrazu. Odpowiedź niesie ubytek zmierzony, żeby Operator widział, ile naprawdę ubyło, zamiast czytać obietnicę | `pathIds:string[]` (opc)<br>`boardId:string` (opc)<br>`precision:int` (opc) | `sizeBytes:int` (wym)<br>`savedBytes:int` (wym)<br>`pathsAffected:int` (wym) |
| `design.vector.path.list` | Zwraca ścieżki wektorowe kompozycji | `boardId:string` (wym)<br>`layerId:string` (opc) | `paths:DesignVectorPath[]` (wym)<br>`total:int` (wym) |
| `design.vector.path.remove` | Usuwa ścieżkę wektorową z kompozycji | `pathId:string` (wym) | `removed:bool` (wym) |
| `design.vector.path.set` | Zakłada albo zmienia ścieżkę wektorową kompozycji. Węzły przychodzą w całości — zmiana jednego węzła idzie tą samą drogą co narysowanie ścieżki, żeby nie było dwóch prawd o jej kształcie | `boardId:string` (wym)<br>`nodes:DesignVectorNode[]` (wym)<br>`pathId:string` (opc)<br>`layerId:string` (opc)<br>`closed:bool` (opc)<br>`fill:DesignFill` (opc)<br>`stroke:DesignStroke` (opc)<br>`name:string` (opc) | `path:DesignVectorPath` (wym) |
| `design.vector.shape.add` | Wstawia kształt podstawowy jako ścieżkę wektorową. Kształt powstaje od razu jako węzły, więc da się go dalej edytować piórem — nie jest osobnym bytem, który potem trzeba zamieniać | `boardId:string` (wym)<br>`kind:DesignShapeKind` (wym)<br>`x:float` (wym)<br>`y:float` (wym)<br>`width:float` (wym)<br>`height:float` (wym)<br>`cornerRadius:float` (opc)<br>`points:int` (opc)<br>`innerRadius:float` (opc)<br>`layerId:string` (opc)<br>`fill:DesignFill` (opc)<br>`stroke:DesignStroke` (opc) | `path:DesignVectorPath` (wym) |
| `design.vector.symbol.list` | Zwraca symbole kompozycji wraz z liczbą instancji | `boardId:string` (wym) | `symbols:DesignSymbol[]` (wym)<br>`total:int` (wym) |
| `design.vector.symbol.set` | Zakłada symbol wielokrotnego użycia albo zmienia jego definicję. Zmiana definicji propaguje do wszystkich instancji — po to symbol jest | `boardId:string` (wym)<br>`name:string` (wym)<br>`symbolId:string` (opc)<br>`pathIds:string[]` (opc)<br>`layerIds:string[]` (opc) | `symbol:DesignSymbol` (wym)<br>`propagatedTo:int` (wym) |
| `design.vector.text.path` | Układa tekst wzdłuż ścieżki albo zamienia go w kontury. Zamiana w kontury jest nieodwracalna dla wyniku, dlatego tekst źródłowy zostaje | `boardId:string` (wym)<br>`text:string` (wym)<br>`fontFamily:string` (wym)<br>`fontSize:float` (wym)<br>`pathId:string` (opc)<br>`outline:bool` (opc)<br>`x:float` (opc)<br>`y:float` (opc)<br>`fill:DesignFill` (opc)<br>`stroke:DesignStroke` (opc) | `path:DesignVectorPath` (wym)<br>`outlined:bool` (wym) |

**Zdarzenia obszaru `design` — 3:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `design.asset.changed` | Zmiana zasobu wizualnego; zasób trafia do Assets Panel po wygenerowaniu | — |
| `design.board.changed` | Zmiana kompozycji Design Board. Dziś kompozycja jeździ w obie strony komendami, ale drugie okno i drugie połączenie nie dowiadują się o zmianie niczym — zdarzenie obszaru jest dokładnie jedno i dotyczy zasobu | — |
| `design.board.presence` | Obecność i położenie kursorów na kompozycji. Zdarzenie **ulotne** — nie ma odpowiednika w bazie i nie jest odtwarzane po ponownym połączeniu | — |

Razem w wykazie: **106 komend** z 1 obszaru kontraktu.

---

*Koniec dokumentu. Moduł Design — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
