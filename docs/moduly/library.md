# Moduł Library — dokumentacja projektowa

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
| **Tytuł** | Moduł Library — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować): komplet okien operacyjnych, katalog funkcji i narzędzi wraz z zależnościami technicznymi oraz punkty sterowania z okna konfiguracji |
| **Środowiska dostępności** | TalkIn, WorkSpace (odpowiednik projektowy: Project Library w module Workspace) |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 4, 9.6, 12, Załącznik A) · Specyfikacja modułów (rozdz. 4.6) · Specyfikacja okien operacyjnych (rozdz. 5, 6.6, 10) · System wizualny (rozdz. 6, 8, 10) · Model danych (rozdz. 16, 17, Załącznik F.2) · Izolacja i zależności (rozdz. 1–6) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); zero blokad w interfejsie; klucze i powiązania jawne; każda zdolność sterowana z okna konfiguracji; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - 2.1. [Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
   - 2.2. [Makieta zbiorcza modułu](#22-makieta-zbiorcza-modułu)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
4. [Przepływy pracy](#4-przepływy-pracy)
5. [Stany, dane i powiązania](#5-stany-dane-i-powiązania)
6. [Scenariusze użycia](#6-scenariusze-użycia)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja

Library jest centralnym repozytorium wiedzy i plików platformy — jednym, uporządkowanym miejscem przechowywania materiałów wykorzystywanych i wytwarzanych w toku pracy z AI. Pełni funkcję trwałej pamięci zewnętrznej, wspólnej dla wielu środowisk i modułów: dokument napisany w Studio, źródło zebrane w Browser, raport z Research czy zasób graficzny z Design trafiają do tego samego, jednego miejsca, zamiast pozostawać rozproszone w poszczególnych modułach, w których powstały.

Library jest zarazem **centralnym systemem zarządzania zasobami cyfrowymi** platformy — jednym miejscem, w którym plik jest przyjmowany, opisywany, klasyfikowany, przeszukiwany, wersjonowany, archiwizowany i udostępniany, niezależnie od modułu, w którym powstał. Moduł łączy w jednym środowisku pracy zdolności, które poza platformą rozdzielone są między klasy osobnych programów:

- **Zarządzanie zasobami cyfrowymi (DAM)** — biblioteka zasobów z metadanymi, miniaturami i galerią mediów.
- **Zarządzanie dokumentami (DMS)** — repozytorium dokumentów z wersjonowaniem, historią i regułami napływu.
- **Menedżer bibliografii i źródeł** — rekordy cytowań, import formatów bibliograficznych, generowanie bibliografii.
- **Wyszukiwanie pełnotekstowe i semantyczne** po treści plików, wraz z warstwą OCR.
- **Katalog mediów** — obraz, audio i wideo z podglądem, miniaturami i transkrypcją.
- **Archiwizacja długoterminowa** — pakowanie archiwalne, metadane utrwalenia, kontrola integralności.
- **Higiena plików** — deduplikacja, normalizacja nazw, wykrywanie niemal-duplikatów.
- **Uniwersalny podgląd i porównanie plików** — renderowanie szerokiego zestawu formatów i widok różnic.

### 1.2. Granica tematyczna modułu

| Poza granicą | Właściwy moduł / miejsce |
|---|---|
| Edycja treści dokumentu (redakcja tekstu, formatowanie) | Studio (Studio Editor) — Library przechowuje i podgląda, edycja otwiera się w module źródłowym |
| Repozytorium kodu, gałęzie, commity | Developer / Git Panel (środowisko CodeStudio) — Library nie prowadzi repozytorium kodu |
| Tworzenie grafiki i obróbka pikselowa | Design (Design Board) — Library kataloguje gotowy zasób graficzny, nie edytuje go |
| Orkiestracja procesów i harmonogramy | Automations — Library dostarcza reguły napływu, nie prowadzi kolejek zadań |
| Zewnętrzne pozyskiwanie źródeł z sieci | Browser / Research — Library przyjmuje ich wynik do trwałego przechowania |
| Trwała pamięć kontekstowa modeli | Zakres „Pamięć” okna konfiguracji — Library przechowuje pliki, nie wektory kontekstu rozmów |

Granica jest konfigurowalna, nie sztywna: powiązania Studio → Library, Research → Library, Browser → Library oraz Library ◄──► Project Library ustanawia Operator w oknie konfiguracji (zakres Integracje). Powiązania są jawne, a moduł działa bez dodatkowych warunków wstępnych.

### 1.3. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Library |
|---|---|
| Zespoły pracujące wielomodułowo | Jedno miejsce odnalezienia materiału niezależnie od tego, w którym module powstał |
| Administratorzy wiedzy organizacyjnej | Katalogowanie, tagowanie i porządkowanie zasobów firmowych |
| Redaktorzy i badacze | Śledzenie wersji dokumentu generowanych automatycznie przez AI w różnych modułach |
| Zespoły projektowe (Workspace) | Odpowiednik ograniczony do projektu — Project Library — dla materiałów jednego przedsięwzięcia |
| Operatorzy porządkujący archiwa | Przegląd, kolekcje, wersjonowanie i odnajdywanie materiałów historycznych |

### 1.4. Po co — wartość modułu

| Problem rozproszenia materiałów | Rozwiązanie w Library |
|---|---|
| Artefakty AI powstają w wielu modułach i giną w ich lokalnych historiach | Library Explorer odbiera je automatycznie jako centralny punkt zbiorczy |
| Brak spójnego systemu porządkowania | Tags & Collections nadaje strukturę niezależną od miejsca powstania pliku |
| Konieczność otwierania osobnej aplikacji do podglądu pliku | File Preview renderuje zawartość bez opuszczania modułu |
| Utrata śladu, która wersja dokumentu jest aktualna | Versioning Panel śledzi każdą wersję, również te wygenerowane automatycznie przez AI |
| Metadane, archiwizacja i higiena rozproszone między narzędziami | Metadata & Archive Panel skupia opis, utrwalenie i audyt zasobu w jednym oknie |

### 1.5. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn albo WorkSpace
        ▼
ŚRODOWISKO (TalkIn / WorkSpace) ── boczna nawigacja modułów
        │
        ▼
MODUŁ: LIBRARY ──────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Library Explorer ·
  Tags & Collections · File Preview · Versioning Panel ·
  Metadata & Archive Panel
        │
        ▼
KARTA SESJI — przegląd repozytorium w kontekście bieżącej pracy
  (repozytorium samo w sobie jest trwałe i wykracza poza jedną kartę)
```

### 1.6. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | TalkIn, WorkSpace |
| Środowisko CodeStudio | Moduł niedostępny — repozytorium kodu obsługuje Developer/Git Panel |
| Komponent własny | Library nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Odpowiednik ograniczony do projektu | Project Library (moduł Workspace) — ten sam wzorzec, zakres ograniczony do jednego projektu |
| Liczba okien operacyjnych | 7 (łącznie z Chat Window i Execution Loop Window) |
| Charakter przechowywania | Trwały — zasoby nie znikają wraz z zamknięciem karty sesji, w której powstały |
| Układ interfejsu | Układ pionowy (podział lewa–prawa); regulacji podlega wyłącznie szerokość kolumn |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji — okno obecne w każdym module i środowisku | Centralny punkt pracy i podstawowy mechanizm sterowania procesami repozytorium |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość | 1 | Widoczne bez interakcji w czasie trwania zlecenia; poza zleceniem reprezentowane wskaźnikiem stanu wykonania w Chat Window | Pętla wykonawcza zleceń repozytoryjnych: dekompozycja, kolejka zadań, kontrola jakości |
| 3 | Library Explorer | Okno list i źródeł | Prawa kolumna, dominująca | 1 | Widoczne bez interakcji — okno wiodące modułu | Przegląd, nawigacja, wyszukiwanie i odbiór artefaktów |
| 4 | Tags & Collections | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy kolekcji w pasku kontekstu Library Explorer albo polecenie w Chat Window; po użyciu panel zwija się | Katalogowanie materiałów, kolekcje, etykiety, taksonomia |
| 5 | File Preview | Podgląd i porównanie | Kolumna robocza wywoływana, szeroka | 2 | Podwójne kliknięcie karty pliku, klawisz `Spacja` na zaznaczeniu albo odnośnik cytowania źródła w Chat Window | Uniwersalny podgląd zawartości pliku i porównanie plików |
| 6 | Versioning Panel | Repozytorium i biblioteka | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Kliknięcie wskaźnika liczby wersji na karcie pliku albo pozycji „Wersje” w menu kebab (⋮) pliku | Śledzenie wersji dokumentu i różnic między wersjami |
| 7 | Metadata & Archive Panel | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu kebab (⋮) karty pliku, pozycja „Metadane i archiwum”; zakładki Archiwum i Higiena należą do warstwy 4 | Metadane, archiwizacja długoterminowa, higiena i audyt repozytorium |

### 2.1. Warstwy widoczności w module

Moduł Library stosuje regułę stopniowego ujawniania funkcjonalności: **jeżeli funkcja nie
jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Repozytorium obejmuje
kilkanaście grup zdolności — import, katalogowanie, wyszukiwanie, podgląd, wersjonowanie,
higienę, archiwizację, bibliografię, udostępnianie i audyt — a mimo to w stanie spoczynku
interfejs modułu prezentuje wyłącznie Chat Window, Library Explorer i pasek kontekstu
repozytorium. Złożoność modułu istnieje w architekturze i pozostaje niewidoczna do chwili
wystąpienia potrzeby użycia danej funkcji; liczba narzędzi, paneli i ustawień nie wpływa
na postrzeganą prostotę interfejsu.

| Warstwa | Zawartość w module Library | Sposób wywołania |
|---|---|---|
| 1 — zawsze widoczna | Chat Window, Library Explorer w widoku domyślnym (siatka miniatur), okruszkowa nawigacja repozytorium, pole wyszukiwania, pasek kontekstu ze znacznikami repozytorium, kolekcji i widoku, wskaźnik stanu wykonania zleceń oraz Execution Loop Window w czasie trwania zlecenia. Warstwa zajmuje ponad 80% powierzchni obszaru roboczego modułu | Widoczne bez interakcji |
| 2 — widoczna na żądanie | Przełącznik widoku Library Explorer (lista, galeria, oś czasu, mapa), rząd filtrów fasetowych, tryb wyszukiwania semantycznego i hybrydowego, panel Tags & Collections, File Preview, Versioning Panel, filtr historii wersji, wybór kanału modelu w Chat Window | Kliknięcie znacznika kontekstowego, ikony, przycisku lub przełącznika, podwójne kliknięcie karty pliku, klawisz `Spacja`; po zakończeniu czynności element zwija się samoczynnie |
| 3 — rozwinięcia kontekstowe | Zestaw operacji na pliku i na zaznaczeniu zbiorczym (tagowanie, przypisanie do kolekcji, eksport, archiwizacja, normalizacja nazw, usunięcie odwracalne), operatory zapytań i zapisane wyszukiwania, Metadata & Archive Panel z zakładką Metadane, akcje pozycji wersji (Podgląd, Przywróć, Porównaj), karty sugestii porządkujących i sugestii tagowania AI, akcje dymka odpowiedzi Wykonawcy | Menu kebab (⋮) karty pliku i pozycji wersji, menu hamburger (☰) repozytorium, menu kontekstowe zaznaczenia, element zbiorczy `Operacje ▼` paska zaznaczenia, panel popover filtra, lista rozwijana |
| 4 — funkcje eksperckie | Zakładki Archiwum i Higiena panelu Metadata & Archive Panel (PDF/A, BagIt, PREMIS/METS, polityki retencji, migawka repozytorium, weryfikacja integralności), dziennik audytu, definicje pól niestandardowych i schematu metadanych, reguły kolekcji inteligentnych i reguły napływu, wsadowa klasyfikacja i operacje wsadowe, eksport SKOS/RDF tezaurusa, eksport paczki migracyjnej, API repozytorium i webhooki, macierz izolacji dostępu do plików | Polecenie języka naturalnego w Chat Window, pasek `/operacje`, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny albo konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

**Zasada jednego kliknięcia.** Każda ukryta funkcja modułu Library jest osiągalna jednym
kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego
w Chat Window. Tags & Collections, File Preview, Versioning Panel oraz Metadata &
Archive Panel otwierają się z poziomu karty pliku bez przechodzenia przez ekrany
pośrednie; zagnieżdżanie funkcji głęboko w hierarchii menu jest w module zabronione.
Ukrycie zmniejsza chaos wizualny repozytorium, nie utrudnia dostępu do jego zdolności.

Mechanizmy ukrywania stosowane w module: menu progresywne (`Widok ▼`, `Filtry ▼`),
grupowanie logiczne akcji w jeden element zbiorczy (`Operacje ▼` zamiast zestawu
przycisków paska zaznaczenia), panele wysuwane znikające całkowicie z przestrzeni
roboczej po zamknięciu (Tags & Collections, Versioning Panel, Metadata & Archive Panel)
oraz znaczniki kontekstowe repozytorium, na przykład
`[Library] [Klienci/Klient X] [Siatka] [Kanał X]`; kliknięcie znacznika otwiera
odpowiedni selektor.

### 2.2. Makieta zbiorcza modułu

```
 Makieta zbiorcza — Moduł Library              Dostępność: TalkIn, WorkSpace
 Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ [Library][Klienci/Klient X]      ☰
  nawigacja │ Użytkownik ↔         │ [ 🔍 szukaj ]  [Widok ▼][Filtry ▼]
  modułów   │ Wykonawca            │
  (poza     │                      │ 🏠 / Klienci / Klient X /
  zakresem  │ strumień odpowiedzi  │
  dokumentu)│                      │  [📄] raport-Q2.pdf         ⋮
            │                      │  [📄] notatki.md            ⋮
            │                      │  [🖼] grafika-1.png         ⋮
            │                      │  [📄] umowa.docx            ⋮
            │ ──────────────────   │
            │ [📎] pole poleceń    │ Zaznaczono: 0
            │ [Kanał X ▼] [Wyślij▶]│ [ Operacje ▼ ]
 ═══════════════════════════════════════════════════════════════════════════
 Execution Loop Window (Koordynator ↔ Wykonawca) otwiera się jako kolumna
 sąsiadująca z Chat Window w czasie trwania zlecenia; poza zleceniem obecny
 jest wyłącznie wskaźnik stanu wykonania. Tags & Collections, File Preview,
 Versioning Panel i Metadata & Archive Panel są rozszerzeniami bocznymi
 warstw 2–3 — w stanie spoczynku nie zajmują przestrzeni roboczej.
```

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (okno wspólne)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — kanał Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI / agent / system wykonawczy). Stanowi centralny punkt pracy użytkownika w module Library i podstawowy mechanizm sterowania wszystkimi procesami repozytorium: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Okno zajmuje to samo miejsce układu we wszystkich modułach i środowiskach — lewą kolumnę obszaru roboczego.

**Zawartość i pełny arsenał funkcji.**

- Polecenia w kontekście zgromadzonych materiałów: „znajdź wszystkie dokumenty dotyczące klienta X”, „podsumuj zawartość tej kolekcji”, „porównaj te dwa pliki”.
- Zapytania naturalne tłumaczone na filtr fasetowy i zapytanie do indeksu repozytorium.
- Odpowiedzi z cytowaniem źródła — wskazanie konkretnego pliku i fragmentu, na którym oparta jest odpowiedź.
- Strumień odpowiedzi modelu na żywo kanałem WebSocket.
- Odwołania do konkretnego pliku lub kolekcji z Library Explorer wzmianką w treści polecenia.
- Generowanie opisu i metadanych pliku na podstawie jego treści (automatyczne streszczenie dodawane jako notatka pliku).
- Sugestie porządkujące: „te pliki wyglądają na duplikaty”, „ta kolekcja nie ma jeszcze etykiety tematycznej”.
- Wstawienie wyniku polecenia jako nowej notatki przypisanej do pliku.
- Tryb operacji wsadowych wywoływany z paska `/operacje` — klasyfikacja, tagowanie i archiwizacja całych zbiorów.
- Zatwierdzenie i przerwanie działania Wykonawcy, historia poleceń, regeneracja odpowiedzi.

**Makieta tekstowa.**

```
┌─ Chat Window ──────────────────┐   stan spoczynku
│ [Kanał X ▼]              [⋮]   │
├────────────────────────────────┤
│ ┌─ Użytkownik ───────────────┐ │
│ │ „znajdź dokumenty klienta X │ │
│ │  z Q2”                      │ │
│ └────────────────────────────┘ │
│ ┌─ Wykonawca (strumień) ─────┐ │
│ │ lista 6 plików z odnośnikami│ │
│ │ źródło: raport-Q2.pdf, s. 3 │ │
│ │                         [⋮] │ │
│ └────────────────────────────┘ │
├────────────────────────────────┤
│ [📎]  Pole poleceń             │
│                     [ Wyślij ▶]│
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia dotyczącego repozytorium | Duże pole w stopce kolumny | 1 | Widoczne bez interakcji | domyślny · fokus · z odwołaniem do pliku | Enter wysyła; `@` wywołuje odwołanie do pliku lub kolekcji | Stopka Chat Window |
| Strumień odpowiedzi | Blok treści odpowiedzi Wykonawcy | Prezentacja wyniku i przebiegu pracy | Duży obszar przewijany | 1 | Widoczne bez interakcji | oczekiwanie · strumień na żywo · zakończony | Przewijanie do najnowszej odpowiedzi | Ciało Chat Window |
| Odnośnik cytowania źródła | Odsyłacz do pliku i fragmentu | Wskazanie podstawy odpowiedzi Wykonawcy | Bardzo mały odsyłacz `.dn-tekst-3` | 1 | Widoczne bez interakcji w treści odpowiedzi | domyślny · najechanie (podgląd fragmentu) | Kliknięcie otwiera File Preview na wskazanym fragmencie | W treści odpowiedzi |
| Selektor kanału modelu | Zwinięty element `Kanał X ▼` w nagłówku | Wybór i zmiana modelu bazowego bieżącej karty | Mały przycisk z etykietą | 2 | Kliknięcie znacznika kontekstowego kanału; po wyborze lista zwija się samoczynnie | zwinięty (domyślny) · rozwinięty · ładowanie | Wybór zmienia model dla kolejnych poleceń | Nagłówek Chat Window |
| Przycisk przerwania | Mały przycisk `--zarys` | Zatrzymanie trwającego działania Wykonawcy | Mały przycisk przy strumieniu | 2 | Pojawia się wyłącznie na czas pracy Wykonawcy | ukryty · widoczny w trakcie pracy | Przerywa bieżące działanie, zachowuje wynik częściowy | Przy strumieniu odpowiedzi |
| Menu akcji dymka odpowiedzi | Menu kebab (⋮) przy odpowiedzi | Zbiorcze udostępnienie akcji: „Pokaż w Explorerze”, „Kopiuj”, „Regeneruj” | Zwinięta ikona kebab, pasek dymka | 3 | Kliknięcie ikony kebab (⋮) rozwija listę akcji; po wyborze menu zamyka się | zwinięty (domyślny) · rozwinięty | „Pokaż w Explorerze” filtruje Library Explorer do wskazanego zestawu plików | Pasek akcji dymka odpowiedzi |
| Menu historii poleceń | Pozycja menu kebab nagłówka | Dostęp do historii poleceń i wznowienia wątku | Pozycja listy rozwijanej | 3 | Menu kebab (⋮) nagłówka Chat Window | zwinięty · rozwinięty | Wybór pozycji przywraca wskazane polecenie do pola poleceń | Nagłówek Chat Window |
| Karta sugestii porządkującej | Wyróżniony blok w treści odpowiedzi | Wskazanie porządkowe (duplikat, brak etykiety) | Mały blok `.dn-alert--info` | 3 | Pojawia się w odpowiedzi Wykonawcy jako rozwinięcie kontekstowe wyniku | ukryty · widoczny | Kliknięcie „Zastosuj” wykonuje sugerowaną akcję porządkową | W treści odpowiedzi Wykonawcy |
| Tryb operacji wsadowych | Pasek `/operacje` | Klasyfikacja, tagowanie i archiwizacja całych zbiorów | Pasek poleceń rozszerzony | 4 | Polecenie języka naturalnego albo wpisanie `/operacje` w polu poleceń; niewidoczny dla użytkownika podstawowego | ukryty · aktywny | Zleca operację wsadową pętli wykonawczej | Pole poleceń Chat Window |

---

### 3.2. Execution Loop Window

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — kanał Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Przebieg pętli powiązany z kartą sesji; log zleceń zapisywany trwale przy repozytorium |

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów modułu Library. Zlecenie przyjęte w Chat Window Koordynator dekomponuje na zadania właściwe repozytorium — indeksacja napływających plików, ekstrakcja treści i metadanych, OCR skanów, generowanie miniatur, klasyfikacja i tagowanie, wykrywanie duplikatów, przeliczanie kolekcji inteligentnych, konwersja archiwalna, weryfikacja integralności — przydziela je Wykonawcy i nadzoruje ich wykonanie.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie repozytoryjne i jego dekompozycja na zadania (np. „skataloguj 240 plików importu” → ekstrakcja typu, sumy kontrolne, ekstrakcja tekstu, indeksacja, auto-metadane, sugestie etykiet).
- Kolejka zadań i ich stan: oczekujące, w toku, zakończone, ponowione, odrzucone.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, z ujawnieniem treści polecenia i odpowiedzi.
- Wyniki kontroli jakości: poprawność ekstrakcji treści, kompletność metadanych, zgodność sumy kontrolnej, trafność klasyfikacji — wraz z decyzją o ponowieniu zadania.
- Wskaźniki przebiegu pętli: liczba zadań, postęp, czas trwania, liczba ponowień, zadania zablokowane zależnością (np. brak warstwy OCR przed indeksacją).
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie i korekta zlecenia bez utraty zadań już wykonanych.
- Log przebiegu zasilający dziennik audytu repozytorium.

**Makieta tekstowa.**

```
┌─ Execution Loop Window ────────┐   stan spoczynku
│ Zlecenie: „skataloguj import    │
│ 240 plików”     [⏸]        [⋮] │
├────────────────────────────────┤
│ KOORDYNATOR → WYKONAWCA         │
│  ▸ zad. 1 wykrycie typu    ✓    │
│  ▸ zad. 2 sumy kontrolne   ✓    │
│  ▸ zad. 3 ekstrakcja treści ◐   │
│  ▸ zad. 4 OCR skanów      ⏳    │
│  ▸ zad. 5 indeksacja      ⏳    │
│  ▸ zad. 6 auto-metadane   ⏳    │
├────────────────────────────────┤
│ Postęp 2/6 · ponowień 1     [▼]│
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Pasek z treścią bieżącego zlecenia | Identyfikacja realizowanego zlecenia repozytoryjnego | Mały pasek, góra kolumny | 1 | Widoczne bez interakcji w czasie trwania zlecenia | domyślny · zlecenie skorygowane | — | Góra Execution Loop Window |
| Pozycja zadania | Wiersz z nazwą zadania i znacznikiem stanu | Reprezentacja pojedynczego zadania pętli | Średni wiersz `.dn-karta--pozycja` | 1 | Widoczne bez interakcji | oczekujące · w toku · zakończone · ponowione · zablokowane zależnością | Kliknięcie rozwija komunikaty sterujące i wynik zadania (warstwa 3) | Kolejka zadań |
| Wskaźnik przebiegu | Pasek postępu z licznikami | Stan zaawansowania pętli | Mały pasek `.dn-postep` z tekstem | 1 | Widoczne bez interakcji | domyślny · zakończony | — (informacyjny) | Stopka okna |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie i przerwanie jako jeden element sterujący | Kontrola realizacji pętli wykonawczej | Ikona przełączna `⏸`/`▶` z przerwaniem w menu kebab | 2 | Kliknięcie ikony albo skrót `Ctrl/Cmd + .` | aktywny · wstrzymany · przerwany | Wstrzymanie zatrzymuje kolejkę bez utraty stanu zadań | Nagłówek okna |
| Blok kontroli jakości | Wyróżniony blok z wynikiem sprawdzenia | Prezentacja oceny wyniku i decyzji o ponowieniu | Mały blok `.dn-alert--info` | 3 | Rozwinięcie wskaźnika przebiegu (`▼`) albo kliknięcie pozycji zadania | zwinięty (domyślny) · wynik pozytywny · wynik z ponowieniem | Kliknięcie otwiera zadanie, którego dotyczy ocena | Rozwinięcie pod kolejką zadań |
| Komunikaty sterujące Koordynator ↔ Wykonawca | Pełna treść poleceń i odpowiedzi zadania | Wgląd w wymianę sterującą pętli | Lista `.dn-tekst-3` w rozwinięciu zadania | 3 | Kliknięcie pozycji zadania w kolejce | zwinięty (domyślny) · rozwinięty | Rozwinięcie zwija się po ponownym kliknięciu | Rozwinięcie pozycji zadania |
| Korekta zlecenia | Pole zmiany zakresu zlecenia w trakcie realizacji | Modyfikacja zlecenia bez utraty zadań wykonanych | Pozycja menu kebab i pole tekstowe | 3 | Menu kebab (⋮) nagłówka okna, pozycja „Korekta zlecenia” | zwinięty · dostępny przy zleceniu wstrzymanym | Otwiera pole korekty; Koordynator przelicza dekompozycję zadań | Nagłówek okna |
| Log przebiegu pętli | Pełny zapis przebiegu zasilający dziennik audytu | Diagnostyka niskiego poziomu realizacji zleceń | Lista zdarzeń append-only | 4 | Polecenie języka naturalnego w Chat Window albo tryb administracyjny; niewidoczny dla użytkownika podstawowego | ukryty · otwarty | Otwiera log z filtrem zakresu i czasu | Dziennik audytu repozytorium |

---

### 3.3. Library Explorer

| Aspekt | Wartość |
|---|---|
| Typologia | Okno list i źródeł |
| Waga wizualna | Prawa kolumna, dominująca |
| Izolacja domyślna | Repozytorium trwałe, wspólne dla środowisk TalkIn i WorkSpace; udział poszczególnych modułów konfigurowalny |

**Zawartość i pełny arsenał funkcji.**

*Nawigacja i widoki.*
- Struktura folderów i przestrzeni oraz widok płaski z filtrowaniem.
- Widok siatki miniatur, widok listy szczegółowej, widok galerii (dla materiałów graficznych), widok osi czasu (chronologia dodania) oraz widok mapy oparty na współrzędnych geolokalizacji z metadanych EXIF.
- Okruszkowa nawigacja (breadcrumb) i szybki powrót do korzenia repozytorium.

*Odbiór artefaktów.*
- Automatyczny odbiór artefaktów generowanych przez AI w innych modułach (Studio, Research, Browser, Design) zgodnie z ustanowionymi powiązaniami.
- Ręczne przeciągnięcie i upuszczenie plików z urządzenia lokalnego.
- Import zbiorczy (wiele plików lub archiwum ZIP rozpakowywane automatycznie), import przez adres URL oraz import z folderu obserwowanego.

*Wyszukiwanie i filtrowanie.*
- Wyszukiwanie pełnotekstowe wewnątrz dokumentów (w tym OCR dla skanów i obrazów z tekstem).
- Filtry fasetowe łączone koniunkcyjnie: typ pliku, moduł pochodzenia, data, autor (użytkownik/model), etykieta, kolekcja, rozmiar.
- Wyszukiwanie semantyczne i hybrydowe — odnalezienie plików o zbliżonej treści, nie tylko dopasowaniu słów kluczowych.
- Operatory zapytań (AND/OR/NOT, zakresy dat, `typ:pdf`) oraz zapisane wyszukiwania (widoki inteligentne aktualizujące się automatycznie).
- Wyszukiwanie po podobieństwie obrazu.

*Operacje na plikach.*
- Otwarcie pliku we właściwym module źródłowym (np. dokument tekstowy otwiera się w Studio Editor).
- Zbiorcze operacje: przeniesienie, tagowanie, dodanie do kolekcji, eksport, archiwizacja, deduplikacja, normalizacja nazw, usunięcie (odwracalne przez kosz).
- Wykrywanie duplikatów i niemal-duplikatów (podobna treść, różne nazwy).
- Udostępnienie pliku odnośnikiem (gdy funkcja udostępniania skonfigurowana).
- Podgląd metadanych: rozmiar, typ, moduł i sesja pochodzenia, historia dostępu.

**Makieta tekstowa.**

```
┌─ Library Explorer ─────────────────────────┐  stan spoczynku
│ [Library] [Klienci/Klient X] [Siatka]    ☰  │
│ [ 🔍 szukaj ]        [ Widok ▼ ][ Filtry ▼ ]│
├─────────────────────────────────────────────┤
│ 🏠 / Klienci / Klient X /                    │
│                                              │
│  [📄] raport-Q2.pdf ⋮  [📄] notatki.md    ⋮  │
│  Studio · 3 wersje     Research · 1 wersja   │
│                                              │
│  [🖼] grafika-1.png ⋮  [📄] umowa.docx    ⋮  │
│  Design · 1 wersja     Browser · 2 wersje    │
├─────────────────────────────────────────────┤
│ Zaznaczono: 0                [ Operacje ▼ ]  │
└─────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pole wyszukiwania | Pole `.dn-input` z ikoną (szukaj) | Wyszukiwanie pełnotekstowe w treści repozytorium | Średnie pole, góra kolumny | 1 | Widoczne bez interakcji; `Ctrl/Cmd + F` ustawia w nim fokus | puste · z wynikami | Enter uruchamia wyszukiwanie na żywo | Górna krawędź okna |
| Pasek kontekstu repozytorium | Zestaw znaczników `[Library] [kolekcja] [widok]` | Wskazanie zasięgu i formy bieżącego przeglądu | Rząd lekkich znaczników kontekstowych | 1 | Widoczne bez interakcji | domyślny · znacznik aktywny | Kliknięcie znacznika otwiera odpowiedni selektor (warstwa 2) | Górna krawędź okna |
| Okruszkowa nawigacja | Ścieżka folderów | Orientacja w strukturze i szybki powrót | Mały pasek tekstowy z separatorami | 1 | Widoczne bez interakcji | domyślny | Kliknięcie segmentu przenosi do tego poziomu struktury | Nad listą plików |
| Karta pliku (widok siatki) | Miniatura + nazwa + metadane | Reprezentacja pojedynczego artefaktu | Średnia karta `.dn-karta` w układzie siatki | 1 | Widoczne bez interakcji | domyślny · zaznaczony (obrys złoty) · najechanie · ładowanie miniatury | Pojedyncze kliknięcie zaznacza, podwójne otwiera File Preview lub moduł źródłowy | Ciało okna, widok siatki |
| Plakietka modułu źródłowego | Mała etykieta pod nazwą pliku | Wskazanie modułu, w którym artefakt powstał | Bardzo mała plakietka `.dn-plakietka--wersaliki` | 1 | Widoczne bez interakcji na karcie pliku | domyślny | Kliknięcie filtruje widok do tego modułu źródłowego | Karta pliku |
| Wskaźnik liczby wersji | Mały tekst pomocniczy | Informacja o liczbie zapisanych wersji pliku | Bardzo mały tekst `.dn-tekst-3` | 1 | Widoczne bez interakcji na karcie pliku | jedna wersja · wiele wersji | Kliknięcie otwiera Versioning Panel dla tego pliku | Karta pliku |
| Pusty stan repozytorium | Komunikat z ikoną | Informacja o braku materiałów w bieżącym widoku | Średni blok `.dn-pusty-stan` wyśrodkowany | 1 | Widoczny wyłącznie gdy lista jest pusta | widoczny tylko gdy lista pusta | Zawiera skrót „Dodaj plik” | Środek okna, zamiast listy |
| Przełącznik widoku | Element zwinięty `Widok ▼` | Wybór formy prezentacji zawartości | Mały przycisk z listą rozwijaną | 2 | Kliknięcie `Widok ▼` albo znacznika kontekstowego widoku; po wyborze lista zwija się | zwinięty (siatka domyślnie) · rozwinięty (lista · galeria · oś czasu · mapa) | Przelicza prezentację tej samej listy plików | Górna krawędź okna |
| Rząd filtrów fasetowych | Element zwinięty `Filtry ▼` z panelem popover | Zawężenie widoku wg typu, modułu, daty, etykiety, kolekcji, rozmiaru | Mały przycisk; po rozwinięciu panel popover z fasetami | 2 | Kliknięcie `Filtry ▼`; panel zamyka się po zastosowaniu filtra | zwinięty (brak filtra) · rozwinięty · filtr aktywny (znacznik przy przycisku) | Każdy wybór filtruje listę na żywo, filtry łączą się koniunkcyjnie | Pod polem wyszukiwania |
| Przełącznik trybu wyszukiwania | Ikona trybu przy polu wyszukiwania | Przejście z wyszukiwania pełnotekstowego na semantyczne lub hybrydowe | Bardzo mała ikona w polu | 2 | Kliknięcie ikony trybu w polu wyszukiwania | pełnotekstowy (domyślny) · semantyczny · hybrydowy | Powtarza bieżące zapytanie w wybranym trybie | Pole wyszukiwania |
| Menu operacji na pliku | Menu kebab (⋮) karty pliku | Otwarcie pełnego zestawu akcji pojedynczego pliku | Zwinięta ikona kebab na karcie | 3 | Kliknięcie ikony kebab (⋮) albo menu kontekstowe karty | zwinięty (domyślny) · rozwinięty | Pozycje: Podgląd, Wersje, Metadane i archiwum, Tagowanie, Eksport, Archiwizacja, Usunięcie odwracalne | Karta pliku |
| Element zbiorczy „Operacje ▼” | Jeden element zamiast zestawu przycisków zbiorczych | Wykonanie operacji na wielu zaznaczonych plikach jednocześnie | Element zbiorczy przy dolnej krawędzi kolumny | 3 | Kliknięcie `Operacje ▼` po zaznaczeniu ≥1 pliku | zwinięty (domyślny) · rozwinięty (Tag · Do kolekcji · Eksport · Archiwizuj · Deduplikacja · Normalizacja nazw) | Wybrana akcja stosuje się do wszystkich zaznaczonych pozycji | Stopka okna |
| Menu repozytorium | Menu hamburger (☰) | Dostęp do zapisanych wyszukiwań, operatorów zapytań i importu | Zwinięta ikona hamburger | 3 | Kliknięcie ikony ☰ w pasku kontekstu | zwinięty (domyślny) · rozwinięty | Wybór pozycji stosuje zapisany widok albo otwiera formularz importu | Górna krawędź okna |
| Reguły napływu i import z folderu obserwowanego | Definicje automatycznego wciągania plików | Sterowanie źródłami zasilającymi repozytorium | Formularz reguły | 4 | Polecenie języka naturalnego w Chat Window albo okno konfiguracji (tryb administracyjny); niewidoczne dla użytkownika podstawowego | ukryty · otwarty | Zapis reguły uruchamia automatyczny napływ zgodnie z jej warunkiem | Konfiguracja modułu |

---

### 3.4. Tags & Collections

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Struktura tagów i kolekcji wspólna dla repozytorium; kolekcje projektowe ograniczone do Project Library |

**Zawartość i pełny arsenał funkcji.**

- Tworzenie kolekcji swobodnych (foldery tematyczne niezależne od struktury plików) oraz kolekcji inteligentnych (reguła automatycznego przypisania, np. „wszystkie pliki z modułu Research otagowane «rynek X»”).
- Zarządzanie słownikiem etykiet: dodawanie, łączenie duplikujących się etykiet, zmiana koloru etykiety, usuwanie nieużywanej etykiety.
- Taksonomia hierarchiczna i tezaurus — etykiety powiązane relacjami nadrzędny/podrzędny oraz synonimami, w modelu SKOS-lite, z eksportem SKOS/RDF.
- Przypisywanie wielu etykiet do jednego pliku i wielu plików do jednej kolekcji (relacja wiele-do-wielu).
- Hierarchia kolekcji (kolekcje nadrzędne i podkolekcje).
- Sugestie tagowania generowane przez AI na podstawie treści pliku, akceptowane pojedynczo lub zbiorczo.
- Statystyka użycia etykiety (liczba przypisanych plików) i podgląd zawartości kolekcji bez opuszczania panelu.
- Eksport struktury kolekcji jako mapy (do dokumentacji porządku repozytorium).

**Makieta tekstowa.**

```
┌─ Tags & Collections ───────────┐   stan spoczynku
│ [ Nowe ▼ ]                 [⋮] │
├────────────────────────────────┤
│ KOLEKCJE                        │
│  📁 Klienci                  ▼ │
│  📁 Materiały wewn. (28 plików) │
│  🤖 Research · rynek X   (auto) │
├────────────────────────────────┤
│ ETYKIETY / TEZAURUS          ▼ │
│  🏷 rynek-X (14)                │
│  🏷 poufne (3)                  │
│  🏷 do-przeglądu (7)            │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Drzewo kolekcji | Struktura zagnieżdżonych pozycji | Nawigacja po hierarchii kolekcji | Lista z wcięciami, ikonami folderów | 2 | Widoczne po otwarciu panelu; podkolekcje rozwija wyzwalacz `▼` przy pozycji nadrzędnej | zwinięte (domyślnie) · rozwinięte · zaznaczone | Kliknięcie filtruje Library Explorer do zawartości kolekcji | Część panelu z kolekcjami |
| Ikona kolekcji inteligentnej | Mała ikona (odswiez/automatyczna) | Odróżnienie kolekcji regułowej od swobodnej | Bardzo mała ikona przy nazwie | 2 | Widoczna przy nazwie kolekcji po otwarciu panelu | domyślny | Najechanie pokazuje regułę przypisania w dymku | Przy nazwie kolekcji inteligentnej |
| Licznik plików kolekcji | Mała liczba w nawiasie | Rozmiar kolekcji | Bardzo mały tekst pomocniczy | 2 | Widoczny przy nazwie kolekcji po otwarciu panelu | aktualizowany na żywo | — (informacyjny) | Obok nazwy kolekcji |
| Chmura etykiet | Zestaw plakietek `.dn-plakietka` | Przegląd i wybór etykiet repozytorium | Rząd małych pigułek z licznikiem | 2 | Widoczna po otwarciu panelu; pełny słownik rozwija wyzwalacz `▼` nagłówka sekcji | zwinięta do etykiet najczęściej używanych · rozwinięta · zaznaczona (filtr aktywny) | Kliknięcie filtruje Library Explorer do plików z tą etykietą | Część panelu z etykietami |
| Element zbiorczy „Nowe ▼” | Jeden element zamiast pary przycisków tworzenia | Utworzenie kolekcji swobodnej, kolekcji inteligentnej albo etykiety | Mały element zbiorczy `--zarys` | 3 | Kliknięcie `Nowe ▼` rozwija listę pozycji; po wyborze lista zwija się | zwinięty (domyślny) · rozwinięty | Otwiera formularz właściwy wybranej pozycji | Góra panelu |
| Menu zarządzania słownikiem | Menu kebab (⋮) panelu | Łączenie duplikujących się etykiet, zmiana koloru, usunięcie nieużywanej etykiety | Zwinięta ikona kebab | 3 | Kliknięcie ikony kebab (⋮) w nagłówku panelu | zwinięty (domyślny) · rozwinięty | Wybór pozycji otwiera właściwą operację słownikową | Góra panelu |
| Relacja tezaurusa | Wskazanie etykiety nadrzędnej lub synonimu | Utrzymanie spójnej taksonomii | Bardzo mały tekst przy etykiecie | 3 | Kliknięcie plakietki etykiety rozwija jej relacje | zwinięta (domyślna) · rozwinięta relacja | Kliknięcie przechodzi do etykiety powiązanej | Przy plakietce etykiety |
| Karta sugestii tagowania AI | Mały blok z proponowaną etykietą | Przyspieszenie katalogowania nowych plików | Mały blok `.dn-alert--info` z dwoma przyciskami | 3 | Pojawia się w panelu wyłącznie w chwili wystąpienia sugestii Wykonawcy | ukryty (brak sugestii) · widoczny | „Zastosuj” przypisuje etykietę, „Odrzuć” usuwa sugestię | Stopka panelu |
| Reguła kolekcji inteligentnej | Definicja warunku automatycznego przypisania plików | Sterowanie zawartością kolekcji regułowej | Formularz reguły | 4 | Polecenie języka naturalnego w Chat Window albo tryb administracyjny; niewidoczna dla użytkownika podstawowego | ukryta · otwarta | Zapis reguły uruchamia przeliczenie kolekcji | Konfiguracja kolekcji |
| Eksport tezaurusa i mapy kolekcji | Wywóz taksonomii w SKOS/RDF oraz mapy porządku repozytorium | Dokumentacja i przeniesienie struktury katalogowej | Operacja eksportu | 4 | Polecenie języka naturalnego w Chat Window albo wyszukiwarka funkcji; niewidoczny dla użytkownika podstawowego | ukryty · przetwarzanie · wynik gotowy | Generuje plik SKOS/RDF albo mapę Markdown/JSON do pobrania | Panel Tags & Collections |

---

### 3.5. File Preview

| Aspekt | Wartość |
|---|---|
| Typologia | Podgląd i porównanie |
| Waga wizualna | Kolumna robocza wywoływana, szeroka |
| Izolacja domyślna | Renderuje wybrany plik bez trwałego stanu własnego |

**Zawartość i pełny arsenał funkcji.**

- Renderowanie zawartości bez opuszczania modułu dla szerokiego zestawu formatów: dokumenty tekstowe (PDF, DOCX, ODT, RTF, Markdown, TXT), obrazy (w tym RAW i HEIC), arkusze kalkulacyjne (XLSX, CSV), kod źródłowy z podświetleniem składni, audio i wideo (odtwarzacz wbudowany), archiwa (podgląd zawartości bez rozpakowania).
- Podgląd wielostronicowy z miniaturami stron i szybkim przeskokiem.
- Powiększenie, dopasowanie do szerokości i strony, obrót (dla obrazów i skanów).
- Zaznaczenie fragmentu tekstu podglądu i skopiowanie go bezpośrednio, bez otwierania pliku w edytorze.
- Przejście „Otwórz w module źródłowym” — dokument tekstowy otwiera się w Studio Editor, kod w Developer, grafika w Design Board.
- Panel metadanych obok podglądu: rozmiar, typ MIME, data utworzenia i modyfikacji, moduł i sesja pochodzenia, lista etykiet i kolekcji.
- Odtwarzacz audio i wideo z osią czasu, oznaczeniami rozdziałów i transkrypcją pozyskaną z modułu Assistant.
- Porównanie dwóch plików obok siebie (widok podzielony) z różnicami treści dla dokumentów tekstowych, niezależnie od Versioning Panel.
- Degradacja bez błędu: gdy dla formatu brak renderera, okno prezentuje metadane i akcję pobrania zamiast komunikatu o błędzie.

**Makieta tekstowa.**

```
┌─ File Preview ──────────────────────────────┐  stan spoczynku
│ raport-Q2.pdf              [⋮]  [ ✕ Zamknij ]│
├───────────────────────────┬─────────────────┤
│ ┌───────────────────────┐ │ METADANE        │
│ │ wyrenderowana strona   │ │ Typ: PDF·2,4 MB │
│ │ dokumentu              │ │ Moduł: Studio   │
│ │                        │ │ Utw.: 2026-08-04│
│ └───────────────────────┘ │ Zm.: 2026-08-06 │
│ [‹ str. 3/12 ›]  [ 100% ▼ ]│ 🏷rynek-X 🏷poufne│
│                            │ Klienci/Klient X│
├───────────────────────────┴─────────────────┤
│ [ Otwórz w Studio ]            [ Operacje ▼ ]│
└──────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek podglądu | Pasek z nazwą pliku | Identyfikacja podglądanego pliku | Mały pasek, góra kolumny | 2 | Widoczny po otwarciu File Preview | domyślny | — | Góra File Preview |
| Przycisk „✕ Zamknij” | Ikona zamknięcia | Zamknięcie podglądu, powrót do Library Explorer | Mała ikona (zamknij) | 2 | Widoczna po otwarciu File Preview; klawisz `Esc` działa równoważnie | domyślny · najechanie | Zamyka kolumnę, zachowuje pozycję w Library Explorer | Prawy róg nagłówka |
| Obszar renderowania | Wyrenderowana zawartość pliku | Główna prezentacja treści | Duży blok, lewa część kolumny | 2 | Widoczny po otwarciu File Preview | ładowanie (`.dn-spinner`) · wyrenderowany · format nieobsługiwany (metadane i pobranie) | Przewijanie, powiększanie, zaznaczanie tekstu | Lewa część kolumny podglądu |
| Nawigator stron | Para strzałek + numer strony | Przechodzenie między stronami dokumentu wielostronicowego | Mały pasek pod obszarem renderowania | 2 | Widoczny wyłącznie dla dokumentu wielostronicowego | ukryty (plik jednostronicowy) · widoczny | Kliknięcie strzałki zmienia stronę | Pod obszarem renderowania |
| Panel metadanych | Kolumna z listą właściwości | Informacje o pochodzeniu i klasyfikacji pliku | Wąska kolumna prawa | 2 | Widoczny po otwarciu File Preview | domyślny | Kliknięcie etykiety lub kolekcji filtruje Library Explorer | Prawa część kolumny podglądu |
| Przycisk „Otwórz w [moduł źródłowy]” | Przycisk `--zloty` | Przejście do pełnej edycji we właściwym module | Średni przycisk CTA | 2 | Widoczny po otwarciu File Preview jako jedyna akcja pierwszoplanowa | domyślny | Otwiera plik w Studio Editor / Design Board / Code Editor, zależnie od typu | Stopka kolumny |
| Selektor powiększenia | Element zwinięty `100% ▼` | Powiększenie, dopasowanie do szerokości i do strony, obrót | Mały element zbiorczy z listą | 3 | Kliknięcie `100% ▼`; po wyborze lista zwija się | zwinięty (domyślny) · rozwinięty | Zmienia skalę lub orientację renderowanej zawartości | Pod obszarem renderowania |
| Element zbiorczy „Operacje ▼” | Jeden element zamiast zestawu przycisków stopki | Porównanie z innym plikiem, pobranie, kopiowanie fragmentu, udostępnienie odnośnikiem | Element zbiorczy w stopce kolumny | 3 | Kliknięcie `Operacje ▼` albo menu kebab (⋮) nagłówka | zwinięty (domyślny) · rozwinięty | „Porównaj z…” otwiera Library Explorer w trybie wyboru i widok dwóch plików obok siebie z podświetleniem różnic | Stopka kolumny |
| Odtwarzacz audio i wideo | Oś czasu, rozdziały i transkrypcja nagrania | Odsłuch i przegląd materiału multimedialnego | Odtwarzacz w obszarze renderowania | 3 | Pojawia się w miejsce obszaru renderowania dla pliku audio lub wideo | ukryty · odtwarzanie · wstrzymany | Kliknięcie znacznika transkrypcji przenosi do punktu nagrania | Obszar renderowania |
| Transkodowanie formatu przy podglądzie | Konwersja HEIC/MOV do formatu przeglądarki w locie | Renderowanie formatów nieobsługiwanych natywnie | Operacja procesowa bez reprezentacji stałej | 4 | Uruchamiane poleceniem języka naturalnego albo automatycznie przez pętlę wykonawczą; sterowanie w trybie administracyjnym | ukryte · przetwarzanie · wynik gotowy | Podgląd otwiera się na skonwertowanej reprezentacji pliku | File Preview |

---

### 3.6. Versioning Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Repozytorium i biblioteka |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Wersje narastają wraz z każdą zmianą dokumentu w dowolnym module |

**Zawartość i pełny arsenał funkcji.**

- Chronologiczna lista wersji dokumentu ze znacznikiem czasu, autorem (`uzytkownik` / `model`) i modułem, w którym zmiana powstała.
- Podgląd dowolnej wersji przez File Preview bez tworzenia nowej wersji bieżącej.
- Przywrócenie wcześniejszej wersji jako bieżącej — operacja nieusuwająca, dopisuje nową pozycję na końcu historii.
- Porównanie dwóch dowolnych wersji, również dokumentów binarnych po wyekstrahowanym tekście (współdzielony mechanizm z Diff/Grep Panel modułu Studio, gdy plik pochodzi z tego modułu).
- Etykietowanie wersji kamieni milowych („wersja podpisana”, „wersja wysłana do klienta”).
- Filtrowanie historii: tylko zmiany AI, tylko zmiany użytkownika, tylko wersje oznaczone etykietą.
- Eksport pełnej historii wersji jako archiwum lub jako pojedynczy raport zmian.
- Powiadomienie o nowej wersji wygenerowanej automatycznie w innym module, gdy plik jest aktualnie otwarty w podglądzie.
- Ustawienie okresu retencji historii wersji sterowane z okna konfiguracji.

**Makieta tekstowa.**

```
┌─ Versioning Panel ─────────────┐   stan spoczynku
│ raport-Q2.pdf              [⋮] │
│ [ Filtr ▼ ]                     │
├────────────────────────────────┤
│ ● Wersja 3 · dziś 12:04      ⋮ │
│   model · Studio                │
│   🏷 „wysłana do klienta”       │
├────────────────────────────────┤
│ ○ Wersja 2 · wczoraj 16:20   ⋮ │
│   użytkownik · Research         │
├────────────────────────────────┤
│ ○ Wersja 1 · 2026-08-04      ⋮ │
│   model · Research              │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pozycja wersji | Wiersz karty `.dn-karta--pozycja` | Reprezentuje jedną zapisaną wersję pliku | Średni wiersz z kropką stanu, znacznikiem czasu, autorem i modułem pochodzenia | 2 | Widoczna po otwarciu panelu | bieżąca (wypełniona kropka, wyróżnienie złote) · archiwalna · najechanie | Kliknięcie otwiera wersję w File Preview | Lista główna panelu |
| Plakietka modułu pochodzenia | Mała etykieta | Wskazanie, w którym module powstała dana wersja | Bardzo mała plakietka `.dn-plakietka--wersaliki` | 2 | Widoczna przy pozycji wersji po otwarciu panelu | domyślny | — (informacyjny) | Przy pozycji wersji |
| Powiadomienie o nowej wersji | Toast lub subtelny pasek | Informacja, że plik zmienił się w innym module podczas jego przeglądania | Mały pasek `.dn-alert--info` | 2 | Pojawia się samoczynnie w chwili powstania nowej wersji, kanałem zdarzeń | ukryty · widoczny | Kliknięcie odświeża panel do najnowszej wersji | Góra panelu, pojawia się na żywo |
| Filtr historii | Element zwinięty `Filtr ▼` | Zawężenie listy wersji wg autora lub etykiety | Mały przycisk z listą rozwijaną | 3 | Kliknięcie `Filtr ▼`; po wyborze lista zwija się | zwinięty (wszystkie) · rozwinięty · filtr aktywny | Przelicza widoczną listę wersji | Góra panelu |
| Menu akcji wersji | Menu kebab (⋮) pozycji wersji | Zbiorcze udostępnienie akcji: Podgląd, Przywróć, Porównaj, Etykieta kamienia milowego | Zwinięta ikona kebab przy pozycji | 3 | Kliknięcie ikony kebab (⋮) przy pozycji wersji | zwinięty (domyślny) · rozwinięty | Przywrócenie dopisuje nową pozycję na końcu historii, nie usuwa wersji | Pozycja wersji |
| Menu panelu wersji | Menu kebab (⋮) nagłówka | Eksport historii wersji jako archiwum albo raportu zmian | Zwinięta ikona kebab | 3 | Kliknięcie ikony kebab (⋮) w nagłówku panelu | zwinięty (domyślny) · rozwinięty · ładowanie | Generuje archiwum lub raport zmian do pobrania | Góra panelu |
| Retencja historii wersji | Okres przechowywania wersji dokumentu | Sterowanie długością zachowywanej historii | Ustawienie `retencja_historii` | 4 | Polecenie języka naturalnego w Chat Window albo okno konfiguracji (tryb administracyjny); niewidoczna dla użytkownika podstawowego | ukryta · ustalona | Zapis okresu zmienia zakres przechowywanej historii wersji | Konfiguracja modułu |

---

### 3.7. Metadata & Archive Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Operuje na metadanych i stanie archiwalnym artefaktu; dziennik audytu zapisywany trwale, przyrostowo |

Metadata & Archive Panel skupia trzy obszary pracy z zasobem, które w innym układzie rozpraszają się między pozostałe okna: opis, utrwalenie i higienę.

**Zawartość i pełny arsenał funkcji.**

*Metadane.*
- Edycja pól schematu Dublin Core oraz pól niestandardowych zdefiniowanych przez Operatora per typ i kolekcja.
- Podgląd metadanych technicznych: EXIF, IPTC, XMP dla obrazów, właściwości dokumentów, ID3 dla audio.
- Auto-metadane AI: opis, słowa kluczowe, streszczenie, klasyfikacja typu treści, wykryty język.
- Rozpoznanie encji (nazwiska, firmy, daty, kwoty) jako pól wyszukiwalnych.
- Rekordy bibliograficzne: autor, tytuł, rok, DOI, ISBN, adnotacje i wypisy do źródła.

*Archiwizacja.*
- Konwersja do formatu archiwalnego PDF/A wraz z walidacją wyniku.
- Pakiet archiwalny BagIt z manifestem sum kontrolnych.
- Metadane utrwalenia w profilu PREMIS/METS.
- Polityki retencji i utylizacji zasobów, bez twardego usunięcia pozbawionego potwierdzenia.
- Migawka repozytorium i eksport paczki migracyjnej kolekcji.

*Higiena i audyt.*
- Raport duplikatów i niemal-duplikatów wraz z operacją scalenia.
- Raport plików osieroconych — bez etykiety, kolekcji lub powiązania.
- Weryfikacja integralności (fixity check) na podstawie sum kontrolnych.
- Normalizacja nazw plików według reguł nazewnictwa.
- Dziennik audytu: kto, co i kiedy — dostęp, zmiana, archiwizacja, przywrócenie, eksport.
- Pulpit stanu repozytorium: liczba plików, rozmiar, rozkład typów, moduły pochodzenia, statystyka użycia etykiet, raport retencji, historia dostępu do pliku.

**Makieta tekstowa.**

```
┌─ Metadata & Archive Panel ─────┐   stan spoczynku
│ raport-Q2.pdf              [⋮] │
│ [ Obszar: Metadane ▼ ]          │
├────────────────────────────────┤
│ DUBLIN CORE                     │
│  Tytuł:   Raport Q2 · Klient X  │
│  Twórca:  model · Studio        │
│  Temat:   🏷 rynek-X            │
│  Data:    2026-08-04            │
│  Prawa:   wewnętrzne            │
├────────────────────────────────┤
│ TECHNICZNE                   ▼ │
│  SHA-256: 3f9a…c21              │
└────────────────────────────────┘
 Zakładki Archiwum i Higiena należą do warstwy 4 — otwiera je
 polecenie w Chat Window albo tryb administracyjny.
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Formularz Dublin Core | Zestaw pól tekstowych | Opis zasobu w standardowym schemacie metadanych | Średni formularz `.dn-input` | 3 | Otwarcie panelu z menu kebab (⋮) karty pliku, pozycja „Metadane i archiwum” | domyślny · edytowany · uzupełniony przez AI (wyróżnienie) | Zapis aktualizuje encję `artefakt` i indeks | Zakładka Metadane |
| Selektor obszaru | Element zwinięty `Obszar: Metadane ▼` zamiast zestawu zakładek | Przełączanie między metadanymi, archiwum i higieną | Mały element zbiorczy z listą | 3 | Kliknięcie `Obszar ▼`; pozycje Archiwum i Higiena dostępne przy roli z uprawnieniem eksperckim | zwinięty (Metadane) · rozwinięty | Zmienia zawartość panelu bez utraty zaznaczenia pliku | Góra panelu |
| Blok metadanych technicznych | Lista właściwości odczytanych z pliku | Prezentacja EXIF/IPTC/XMP/ID3 i sumy kontrolnej | Wąska lista `.dn-tekst-3` | 3 | Rozwinięcie sekcji wyzwalaczem `▼` w nagłówku „TECHNICZNE” | zwinięty (domyślny) · rozwinięty · brak danych technicznych | — (informacyjny) | Zakładka Metadane |
| Pole niestandardowe | Pole metadanych zdefiniowane przez Operatora | Rozszerzenie opisu zasobu poza schemat podstawowy | Pozycja menu kebab i pole formularza | 4 | Menu kebab (⋮) panelu w trybie administracyjnym albo polecenie języka naturalnego; niewidoczne dla użytkownika podstawowego | ukryte · dodane | Otwiera wybór pola ze schematu konfiguracji | Zakładka Metadane |
| Operacje formatów archiwalnych | Element zbiorczy zamiast zestawu przycisków — PDF/A, BagIt, PREMIS/METS | Utrwalenie zasobu w formacie archiwalnym | Element zbiorczy `Utrwalenie ▼` | 4 | Zakładka Archiwum otwierana poleceniem w Chat Window albo w trybie administracyjnym | ukryty · przetwarzanie · wynik gotowy | Zleca zadanie pętli wykonawczej i pokazuje jego przebieg w Execution Loop Window | Zakładka Archiwum |
| Wskaźnik retencji | Tekst z okresem przechowywania | Informacja o polityce retencji zasobu | Mały tekst pomocniczy | 4 | Widoczny wyłącznie w zakładce Archiwum | bez limitu · okres ustalony · zbliża się koniec okresu | Kliknięcie otwiera regułę retencji w oknie konfiguracji | Zakładka Archiwum |
| Karty raportów higieny | Zestaw liczników z odsyłaczem | Duplikaty, pliki osierocone, wynik weryfikacji integralności | Małe karty `.dn-karta` | 4 | Zakładka Higiena otwierana poleceniem w Chat Window, wyszukiwarką funkcji albo w trybie administracyjnym | zgodne · wykryte odchylenia | Kliknięcie filtruje Library Explorer do wskazanego zbioru | Zakładka Higiena |
| Dziennik audytu | Przyrostowy log zdarzeń repozytorium | Ustalenie kto, co i kiedy zrobił z zasobem | Lista zdarzeń append-only z filtrem | 4 | Polecenie języka naturalnego w Chat Window albo tryb administracyjny; niewidoczny dla użytkownika podstawowego | ukryty · otwarty · ładowanie | Rozwija listę zdarzeń z filtrem zakresu i czasu | Zakładka Higiena |
| Migawka repozytorium i paczka migracyjna | Zapis stanu repozytorium i eksport kolekcji jako samodzielnego archiwum | Utrwalenie i przeniesienie całości zasobów | Operacja wsadowa | 4 | Polecenie języka naturalnego w Chat Window albo pasek `/operacje`; niewidoczna dla użytkownika podstawowego | ukryta · przetwarzanie · wynik gotowy | Zleca zadanie pętli wykonawczej; wynik trafia do pobrania | Zakładka Archiwum |

---

## 4. Przepływy pracy

### 4.1. Przepływ podstawowy — odbiór, katalogowanie i podgląd

```
 Moduł źródłowy            Library Explorer         Tags & Collections      File Preview
 (Studio/Research/          ─────────────────         ──────────────────      ────────────
  Browser/Design)
 1. artefakt AI      ────────►
    generowany                2. odbiór automatyczny
                                  w repozytorium
                                          │
                              3. nadanie etykiet   ─────────►
                                 i przypisanie
                                 do kolekcji
                                          │
                                                              4. podgląd bez
                                                                 opuszczania modułu ─►
```

### 4.2. Przepływ pętli wykonawczej — zlecenie katalogowania

```
Chat Window (Użytkownik ↔ Wykonawca)
   „skataloguj import 240 plików”
        │
        ▼
Execution Loop Window (Koordynator ↔ Wykonawca)
   dekompozycja: typ → suma kontrolna → ekstrakcja treści →
   OCR → indeksacja → auto-metadane → sugestie etykiet
        │  kontrola jakości po każdym zadaniu, ponowienie przy braku
        │  warstwy tekstu lub niekompletnych metadanych
        ▼
Library Explorer ──► pliki widoczne, przeszukiwalne, z miniaturami
        │
        ▼
Metadata & Archive Panel ──► uzupełnienie opisu, raport higieny
```

### 4.3. Przepływ rozszerzony — porządkowanie zbiorcze z sugestiami AI

```
Chat Window ──► polecenie: „znajdź duplikaty w kolekcji Klienci”
        │
        ▼
Library Explorer ──► lista kandydatów na duplikaty, zaznaczenie zbiorcze
        │
        ▼
Tags & Collections ──► sugestie tagowania AI dla plików bez etykiety
        │
        ▼
Pasek zbiorczych operacji ──► archiwizacja duplikatów, ujednolicenie etykiet
```

### 4.4. Przepływ rozszerzony — śledzenie wersji między modułami

```
Research › Report Builder      Studio › Studio Editor         Library › Versioning Panel
──────────────────────────     ──────────────────────         ──────────────────────────
 wersja 1 raportu    ────────►  redakcja językowa      ──────►  wersja 2 odnotowana
 (artefakt utworzony)                                             automatycznie,
                                                                   moduł pochodzenia: Studio
                                        │
                                        ▼
                              dalsza redakcja      ──────────►  wersja 3, etykieta
                                                                 „wysłana do klienta”
```

### 4.5. Przepływ rozszerzony — utrwalenie archiwalne

```
Library Explorer ──► wybór kolekcji do utrwalenia
        │
        ▼
Metadata & Archive Panel ──► konwersja PDF/A · pakiet BagIt · PREMIS/METS
        │
        ▼
Execution Loop Window ──► kolejka konwersji, weryfikacja integralności,
                          ponowienie zadań o wyniku niezgodnym
        │
        ▼
Eksport paczki migracyjnej ──► archiwum z manifestem i metadanymi JSON
```

### 4.6. Przepływ rozszerzony — odpowiednik projektowy (Project Library)

```
              LIBRARY (repozytorium centralne, TalkIn/WorkSpace)
                              ▲
                              │  materiały ogólnodostępne, poza projektami
                              │
          ┌───────────────────┼───────────────────┐
          │                                        │
 Workspace › Projekt A                    Workspace › Projekt B
 Project Library                          Project Library
 (zakres ograniczony                      (zakres ograniczony
  do jednego projektu)                     do jednego projektu)

  Project Library jest odpowiednikiem Library ograniczonym do projektu —
  nie jest odrębnym mechanizmem, lecz tym samym wzorcem w węższym zasięgu.
```

---

## 5. Stany, dane i powiązania

### 5.1. Model stanów repozytorium

```
┌───────────────┐  napływ pliku    ┌───────────────┐   nadanie etykiety   ┌───────────────┐
│ Plik poza      │  (moduł/import) │ Nieuporząd-    │   lub kolekcji       │ Skatalogowany  │
│ repozytorium   │────────────────►│ kowany         │─────────────────────►│                │
└───────────────┘                  └───────────────┘                       └───────┬───────┘
                                                                                     │
                                                                        nowa wersja z modułu
                                                                                     ▼
                                                                            ┌───────────────┐
                                                                            │ Wersjonowany   │
                                                                            │ (Versioning    │
                                                                            │  Panel)        │
                                                                            └───────┬───────┘
                                                                                     │ utrwalenie
                                                                                     ▼
                                                                            ┌───────────────┐
                                                                            │ Utrwalony      │
                                                                            │ (PDF/A, BagIt) │
                                                                            └───────┬───────┘
                                                                                     │ archiwizacja
                                                                                     ▼
                                                                            ┌───────────────┐
                                                                            │ Zarchiwizowany │
                                                                            │ (odwracalne)   │
                                                                            └───────────────┘
```

### 5.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Library |
|---|---|
| `artefakt` | Każdy plik widoczny w Library Explorer; pola `sesja_id`/`projekt_id` wskazują pochodzenie; pola rozszerzone niosą schemat Dublin Core i pola niestandardowe |
| `wersja_artefaktu` | Każda pozycja w Versioning Panel; pole `autor` rozróżnia `uzytkownik`/`model`, pole `moduł_pochodzenia` wskazuje miejsce zmiany |
| `kolekcja`, `kolekcja_artefakt` | Struktura kolekcji swobodnych i inteligentnych w Tags & Collections; reguła kolekcji inteligentnej zapisana przy encji `kolekcja` |
| `etykieta_artefaktu` | Etykiety przypisane plikom wraz z licznikiem użycia i relacjami tezaurusa |
| `sesja`, `karta_sesji` | Kontekst przeglądu repozytorium; samo repozytorium wykracza poza pojedynczą kartę |
| `powiazanie_komponentu` | Jawne powiązania Studio/Research/Browser ───► Library oraz Library ◄──► Project Library |

### 5.3. Izolacja i konfigurowalność — punkty właściwe modułowi Library

| Punkt izolacji | Stan wyjściowy (domyślny) | Co podlega konfiguracji | Gdzie |
|---|---|---|---|
| Odbiór artefaktów z innych modułów | Zgodnie z powiązaniami ustanowionymi w oknie konfiguracji dla każdego modułu źródłowego z osobna | Włączenie i wyłączenie automatycznego odbioru z konkretnego modułu | Okno konfiguracji, powiązanie komponentu |
| Zakres Library a Project Library | Repozytorium centralne i repozytoria projektowe pozostają odrębne | Współdzielenie wybranej kolekcji między repozytorium centralnym a projektem | Poziom zasięgu „projekt” lub „para modułów” |
| Retencja historii wersji | Przechowywanie bez limitu (brak ustawienia = wartość domyślna) | Ustalenie okresu retencji na dowolnym poziomie zasięgu | Okno konfiguracji, ustawienie `retencja_historii` |
| Izolacja techniczna dostępu do plików | Żaden z ośmiu zakresów nie jest domyślnie aktywny — pełny dostęp w ramach uprawnień Operatora | Ograniczenie zakresu odczytu i zapisu plików dla wybranej sesji lub roli | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

### 5.4. Powiązania z innymi modułami

```
   Studio ────►┐
                │
   Research ───►├───►  LIBRARY  ◄──►  Project Library (moduł Workspace)
                │        (repozytorium centralne)
   Browser ────►┘
   (Design ───► Studio i Apps — pośrednio zasila Library przez te moduły)
```

| Moduł źródłowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Studio | Przekazanie wygenerowanych artefaktów i wersji dokumentów | Konfiguracyjne | Preview Window / Studio Editor → Library Explorer |
| Research | Trwałe przechowanie raportów i zebranych źródeł | Konfiguracyjne | Report Builder / Sources Manager → Library Explorer |
| Browser | Trwałe przechowanie odnotowanych materiałów i archiwów WWW | Konfiguracyjne | Notes Panel / Sources Panel → Library Explorer |
| Assistant | Transkrypcja nagrań zasilająca warstwę tekstu w indeksie | Konfiguracyjne | Assistant → File Preview / indeks repozytorium |
| Automations | Reguły napływu, wyzwalacze zdarzeń, harmonogram konserwacji | Konfiguracyjne | Library Explorer ◄──► Automations |
| Workspace | Project Library jako odpowiednik ograniczony do zakresu projektu | Konfiguracyjne | Library Explorer ◄──► Project Library |

---

## 6. Scenariusze użycia

**Scenariusz 1 — centralne repozytorium wielomodułowe.**
Zespół pracuje równolegle w Studio, Research i Browser nad tym samym klientem. Każdy wygenerowany artefakt — dokument, raport, notatka źródłowa — trafia automatycznie do Library Explorer, gdzie jednym poleceniem w Chat Window użytkownik odnajduje wszystkie materiały dotyczące tego klienta niezależnie od modułu pochodzenia.

**Scenariusz 2 — porządkowanie z pomocą AI.**
Administrator wiedzy poleca w Chat Window odnalezienie duplikatów w kolekcji „Klienci”. Library Explorer prezentuje kandydatów, Tags & Collections sugeruje brakujące etykiety dla plików niekatalogowanych, a pasek zbiorczych operacji pozwala jednym kliknięciem zarchiwizować duplikaty.

**Scenariusz 3 — podgląd bez opuszczania modułu.**
Użytkownik przegląda listę plików klienta i otwiera File Preview dla pliku PDF, sprawdza metadane pochodzenia (moduł Studio, trzy wersje), a następnie otwiera dokument w Studio Editor do dalszej redakcji — bez pobierania pliku na dysk lokalny.

**Scenariusz 4 — śledzenie wersji generowanych automatycznie.**
Raport powstały w module Research trafia do Library jako wersja pierwsza. Po redakcji w Studio powstaje wersja druga, odnotowana automatycznie w Versioning Panel z adnotacją modułu pochodzenia „Studio”. Użytkownik oznacza finalną wersję etykietą „wysłana do klienta”.

**Scenariusz 5 — masowy import pod nadzorem pętli wykonawczej.**
Do repozytorium trafia archiwum z kilkuset skanami. Użytkownik zleca katalogowanie w Chat Window, a Execution Loop Window prezentuje dekompozycję zlecenia na zadania, kolejkę ich realizacji i wynik kontroli jakości: pliki bez warstwy tekstu wracają do OCR i zostają ponownie zindeksowane, po czym całość jest przeszukiwalna pełnotekstowo.

**Scenariusz 6 — utrwalenie archiwalne kolekcji.**
Kolekcja dokumentów zamkniętego przedsięwzięcia zostaje przekonwertowana do PDF/A i spakowana w BagIt z manifestem sum kontrolnych oraz metadanymi PREMIS/METS. Weryfikacja integralności potwierdza zgodność, a paczka migracyjna trafia do archiwum zewnętrznego.

**Scenariusz 7 — praca ze źródłami bibliograficznymi.**
Badacz importuje bazę źródeł w formacie BibTeX, uzupełnia rekordy metadanymi pobranymi po DOI, wiąże źródła z raportem powstałym w module Research i generuje bibliografię w wybranym stylu cytowania.

**Scenariusz 8 — repozytorium projektowe ograniczone do zespołu.**
Kierownik projektu w module Workspace korzysta z Project Library do przechowywania materiałów wyłącznie tego przedsięwzięcia, utrzymując rozdzielenie od repozytorium centralnego Library — zgodnie z domyślnym brakiem współdzielenia kontekstu między projektami.

---

## 7. Katalog funkcji i narzędzi

Legenda pozycji: **Nazwa** — co robi — **zależności** (biblioteki Go, formaty, integracje). Biblioteki wskazane poniżej to realne składniki ekosystemu Go albo silniki wywoływane jako proces zewnętrzny, dobrane pod backend platformy.

### 7.1. Napływ, import i akwizycja zasobów

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Automatyczny odbiór artefaktów** | Przyjmuje pliki generowane przez AI w Studio/Research/Browser/Design zgodnie z powiązaniem komponentu | `powiazanie_komponentu`; kolejka zdarzeń wewnętrznych; encja `artefakt` |
| **Przeciągnij-i-upuść z urządzenia** | Ręczne dodanie plików z dysku lokalnego pojedynczo lub hurtem | HTML5 File API (front); upload chunkowany (backend) |
| **Import archiwum ZIP i rozpakowanie** | Wrzucone archiwum rozpakowuje i kataloguje zawartość zachowując strukturę folderów | `archive/zip`, `mholt/archiver/v4`, `nwaples/rardecode` (RAR), `bodgit/sevenzip` (7z), `ulikunitz/xz` |
| **Import z folderu obserwowanego** | Wskazany katalog jest monitorowany, nowe pliki wciągane automatycznie | `fsnotify/fsnotify`; reguła napływu |
| **Import przez URL i z Browser** | Pobranie zasobu spod adresu jako trwałego pliku wraz z metadanymi źródła | `net/http`; integracja Browser → Library |
| **Import poczty i załączników** | Wciągnięcie plików z wiadomości (EML/MSG) z rozbiciem na załączniki | `emersion/go-message`, `DusanKasan/parsemail` |
| **Import archiwum WWW (WARC)** | Przyjęcie zarchiwizowanej strony jako pojedynczego zasobu z zachowaniem treści | `slyrz/warc`; format WARC; integracja Browser |
| **Skanowanie do repozytorium (ingest OCR)** | Skan lub zdjęcie dokumentu → rozpoznanie tekstu → indeksowanie treści | `otiai10/gosseract` (Tesseract); PDF/A jako wynik |
| **Deklaracja typu przy imporcie** | Wykrycie realnego typu pliku niezależnie od rozszerzenia | `gabriel-vasile/mimetype`, `h2non/filetype` |

### 7.2. Katalogowanie, metadane i klasyfikacja

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Ekstrakcja metadanych technicznych** | Odczyt EXIF/IPTC/XMP z obrazów, właściwości dokumentów, ID3 z audio | `barasher/go-exiftool`, `dsoprea/go-exif`, `rwcarlsen/goexif`; formaty EXIF/IPTC/XMP/ID3 |
| **Ekstrakcja tekstu z dokumentów** | Wydobycie treści z PDF/DOCX/XLSX/PPTX/ODT do indeksu i podglądu | PDF: `ledongthuc/pdf`, `pdfcpu/pdfcpu`; DOCX/XLSX/PPTX: `unidoc/unioffice`; XLSX: `qax-os/excelize`; serwer Apache Tika jako silnik zamienny |
| **Schemat metadanych Dublin Core** | Standaryzowany zestaw pól opisowych: tytuł, twórca, temat, data, prawa | Profil pól Dublin Core; encja `artefakt` (pola rozszerzone) |
| **Pola niestandardowe** | Operator definiuje własne pola metadanych per kolekcja i typ | Definicja schematu w konfiguracji; walidacja typów pól bez blokad |
| **Automatyczne streszczenie i opis AI** | Model generuje opis, słowa kluczowe i streszczenie na podstawie treści pliku | Kanał modelu; treść z ekstraktora; zapis jako notatka pliku |
| **Wykrywanie języka treści** | Automatyczne oznaczenie języka dokumentu | `pemistahl/lingua-go` |
| **Klasyfikacja typu treści AI** | Przypisanie kategorii (umowa, faktura, raport, notatka) modelem | Kanał modelu; taksonomia typów; sugestia z akceptacją |
| **Rozpoznanie encji (NER)** | Wyodrębnienie nazwisk, firm, dat i kwot jako metadanych wyszukiwalnych | Kanał modelu lub biblioteka NER; pola indeksu |
| **Odcisk treści (fixity/checksum)** | Suma kontrolna pliku dla integralności i deduplikacji | `crypto/sha256`, `crypto/blake2b` |

### 7.3. Kolekcje, etykiety i taksonomie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Kolekcje swobodne** | Foldery tematyczne niezależne od fizycznej struktury plików | `kolekcja`, `kolekcja_artefakt` (relacja wiele-do-wielu) |
| **Kolekcje inteligentne (regułowe)** | Automatyczne przypisanie plików spełniających regułę: moduł, etykieta, data | Silnik reguł; zapis reguły w `kolekcja`; przeliczanie na żywo |
| **Hierarchia kolekcji** | Kolekcje nadrzędne i podkolekcje o dowolnej głębokości | Struktura zagnieżdżona; pole rodzica |
| **Słownik etykiet (kontrolowany)** | Zarządzanie zbiorem etykiet: dodawanie, łączenie duplikatów, kolor, usuwanie nieużywanej | `etykieta_artefaktu`; licznik użycia |
| **Taksonomia hierarchiczna i tezaurus** | Etykiety powiązane relacjami nadrzędny/podrzędny oraz synonimami | Model SKOS-lite; eksport SKOS/RDF |
| **Sugestie tagowania AI** | Model proponuje etykiety na podstawie treści, akceptacja pojedyncza lub zbiorcza | Kanał modelu; karta sugestii z akcjami Zastosuj/Odrzuć |
| **Tagowanie zbiorcze** | Nadanie wielu etykiet wielu plikom jednocześnie | Pasek zbiorczych operacji |
| **Mapa struktury kolekcji (eksport)** | Wygenerowanie mapy porządku repozytorium jako dokumentacji | Eksport do Markdown/JSON; reprezentacja grafowa |

### 7.4. Wyszukiwanie: pełnotekstowe, semantyczne, fasetowe

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Indeks pełnotekstowy** | Wyszukiwanie w treści wszystkich dokumentów, nie tylko w nazwach | `blevesearch/bleve` (natywny Go); Meilisearch/Typesense/OpenSearch jako backend zamienny |
| **Wyszukiwanie fasetowe** | Zawężanie po typie, module, dacie, autorze, etykiecie, kolekcji, rozmiarze — koniunkcyjnie | Faset w indeksie; rząd filtrów interfejsu |
| **Wyszukiwanie semantyczne (wektorowe)** | Odnajduje pliki o zbliżonym znaczeniu, nie tylko dopasowaniu słów | Embeddingi z modelu; `pgvector` (PostgreSQL), Qdrant/Milvus albo `asg017/sqlite-vec` |
| **Wyszukiwanie hybrydowe** | Łączy trafność słów kluczowych z podobieństwem semantycznym | Bleve + wektory; scalanie wyników metodą RRF |
| **Wyszukiwanie w OCR** | Przeszukiwanie tekstu rozpoznanego ze skanów i obrazów | Warstwa tekstu z Tesseract w indeksie |
| **Zapisane wyszukiwania (widoki inteligentne)** | Zapytanie zapisane jako auto-aktualizujący się widok | Zapis zapytania; przeliczanie na żywo |
| **Wyszukiwanie po podobieństwie obrazu** | Odnajduje wizualnie podobne grafiki, nie po nazwie | `corona10/goimagehash` (pHash/dHash) lub embeddingi wizualne |
| **Zapytania naturalne przez Chat Window** | „znajdź wszystkie umowy klienta X z Q2” tłumaczone na filtr i wyszukiwanie | Kanał modelu → zapytanie do indeksu; integracja Chat Window |
| **Operatory zapytań** | Składnia zaawansowana: AND/OR/NOT, zakresy dat, `typ:pdf` | Parser zapytań Bleve/Lucene-like |

### 7.5. Podgląd, renderowanie i porównanie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Renderowanie PDF w oknie** | Podgląd stron PDF z miniaturami i przeskokiem, bez otwierania edytora | `gen2brain/go-fitz` (MuPDF) do rasteryzacji stron; PDF i PDF/A |
| **Podgląd DOCX/ODT/RTF** | Renderowanie dokumentów biurowych z zachowaniem układu | `unidoc/unioffice`; LibreOffice headless (`soffice --convert-to`) jako ścieżka pełnego układu |
| **Podgląd arkuszy** | Renderowanie XLSX/CSV jako tabeli z zakładkami arkuszy | `qax-os/excelize`; `encoding/csv` |
| **Podgląd kodu z podświetleniem** | Kod źródłowy z kolorowaniem składni | `alecthomas/chroma` |
| **Podgląd Markdown/HTML** | Renderowanie do sformatowanej treści | `yuin/goldmark`; sanityzacja `microcosm-cc/bluemonday` |
| **Odtwarzacz audio i wideo** | Wbudowane odtwarzanie z osią czasu, rozdziałami i transkrypcją | HTML5 `<audio>`/`<video>`; transkodowanie `ffmpeg` (przez `u2takey/ffmpeg-go`) |
| **Podgląd obrazów** | Powiększenie, dopasowanie, obrót; obsługa RAW i HEIC | `disintegration/imaging`, `govips` (libvips) dla HEIC/TIFF/WebP |
| **Podgląd wnętrza archiwum** | Lista zawartości ZIP/RAR/7z bez rozpakowywania | `archive/zip`, `bodgit/sevenzip`, `nwaples/rardecode` |
| **Kopiowanie fragmentu z podglądu** | Zaznaczenie i skopiowanie tekstu bez otwierania pliku | Warstwa tekstowa PDF i dokumentu |
| **Porównanie dwóch plików obok siebie** | Widok podzielony, różnice treści dla dokumentów tekstowych | `sergi/go-diff` (diff-match-patch), `pmezard/go-difflib` |
| **Format nieobsługiwany — degradacja** | Gdy brak renderera, okno pokazuje metadane i pobranie zamiast błędu | Reguła zastępcza interfejsu |

### 7.6. Wersjonowanie i śledzenie zmian

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Historia wersji dokumentu** | Chronologiczna lista wersji ze znacznikiem czasu, autorem (`uzytkownik`/`model`) i modułem | `wersja_artefaktu`; pola `autor`, `moduł_pochodzenia` |
| **Podgląd dowolnej wersji** | Otwarcie starszej wersji bez czynienia jej bieżącą | File Preview na wskazanej wersji |
| **Przywrócenie wersji (nieusuwające)** | Wskazana wersja staje się bieżącą, dopisując nową pozycję | Operacja append-only |
| **Porównanie dwóch wersji** | Różnice między dowolnymi wersjami, także dokumentów binarnych | `sergi/go-diff`; dla DOCX porównanie po wyekstrahowanym tekście; współdzielenie z Diff/Grep Panel Studio |
| **Etykiety kamieni milowych** | Oznaczenie wersji „podpisana”, „wysłana do klienta” | Pole etykiety wersji |
| **Filtrowanie historii** | Tylko zmiany AI, tylko użytkownika, tylko oznaczone | Filtr po `autor` i etykiecie |
| **Powiadomienie o nowej wersji na żywo** | Toast, gdy plik zmienia się w innym module podczas podglądu | Kanał WebSocket; zdarzenie wersji |
| **Eksport historii wersji** | Pełna historia jako archiwum lub raport zmian | Pakiet ZIP + raport Markdown/PDF |
| **Retencja wersji** | Ustalony okres przechowywania historii wersji | Ustawienie `retencja_historii` w oknie konfiguracji |

### 7.7. Deduplikacja, higiena i normalizacja

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Wykrywanie duplikatów dokładnych** | Identyczne pliki rozpoznane po sumie kontrolnej | `crypto/sha256`; indeks skrótów |
| **Wykrywanie niemal-duplikatów** | Podobna treść, różne nazwy — tekst i obraz | Tekst: shingling/MinHash; obraz: `corona10/goimagehash` |
| **Scalanie duplikatów** | Połączenie duplikatów w jeden zasób z zachowaniem etykiet i wersji | Operacja scalająca metadane |
| **Normalizacja nazw plików** | Reguły ujednolicania nazw: schemat, transliteracja, znaki | Silnik reguł nazewnictwa; `golang.org/x/text` |
| **Wykrywanie plików osieroconych** | Pliki bez etykiety, kolekcji lub powiązania | Reguła audytu |
| **Weryfikacja integralności (fixity check)** | Okresowe sprawdzenie sum kontrolnych i wykrycie uszkodzenia | Harmonogram (integracja Automations); log audytu |
| **Kompresja i optymalizacja** | Zmniejszenie obrazów i PDF bez utraty jakości katalogowej | `govips`, Ghostscript (`gs`) dla PDF |

### 7.8. Archiwizacja i długoterminowe przechowywanie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Archiwizacja odwracalna (kosz)** | Przeniesienie do archiwum z możliwością przywrócenia | Stan `zarchiwizowany` encji; operacja odwracalna |
| **Pakiet archiwalny BagIt** | Spakowanie zasobu z manifestem sum kontrolnych do przechowania lub przekazania | Generator BagIt (format Library of Congress) |
| **Konwersja do PDF/A** | Normalizacja dokumentów do formatu archiwalnego | Ghostscript (`gs` z profilem PDF/A), `pdfcpu` (walidacja) |
| **Metadane utrwalenia (PREMIS/METS)** | Zapis metadanych o pochodzeniu, zdarzeniach i integralności | Profil PREMIS/METS (XML); generator |
| **Polityki retencji i utylizacji** | Reguły „przechowuj N lat, potem oznacz do przeglądu” | Reguły retencji; harmonogram; brak twardego usunięcia bez potwierdzenia |
| **Eksport paczki migracyjnej** | Cała kolekcja jako samodzielne archiwum z indeksem, do pracy bez sieci | ZIP + manifest + eksport metadanych JSON |
| **Migawka repozytorium** | Zapis stanu całego repozytorium w punkcie czasu | Content-addressable store; MinIO (S3) |
| **Content-addressable storage** | Przechowywanie po skrócie treści — deduplikacja u podstaw | MinIO/S3; adresacja po SHA-256 |

### 7.9. Multimedia — obraz, audio, wideo

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Generowanie miniatur** | Miniatury dla obrazów, stron PDF i klatek wideo | `disintegration/imaging`, `go-fitz` (PDF), `ffmpeg` (klatka wideo) |
| **Podgląd galerii i lightbox** | Widok galerii dużych miniatur dla materiałów graficznych | Widok galerii Library Explorer |
| **Transkrypcja audio i wideo** | Zamiana mowy na tekst wyszukiwalny | Integracja modułu Assistant; warstwa transkryptu w indeksie |
| **Znaczniki i rozdziały na osi czasu** | Oznaczenia punktów w nagraniu z opisem | Metadane czasowe; odtwarzacz HTML5 |
| **Ekstrakcja klatek kluczowych** | Reprezentatywne klatki wideo jako podgląd | `ffmpeg` (detekcja scen) |
| **Odczyt danych geolokalizacji** | Widok mapy zdjęć na podstawie EXIF GPS | `dsoprea/go-exif`; współrzędne GPS |
| **Paleta barw zasobu graficznego** | Dominujące kolory jako metadana wyszukiwalna | `EdlinOrg/prominentcolor` |
| **Transkodowanie formatu przy podglądzie** | Konwersja HEIC/MOV do formatu przeglądarki w locie | `govips` (HEIC), `ffmpeg` (MOV→MP4) |

### 7.10. Bibliografia i menedżer źródeł

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rekordy bibliograficzne** | Zasób opatrzony metadanymi cytowania: autor, tytuł, rok, DOI, ISBN | Profil pól bibliograficznych; encja `artefakt` |
| **Import BibTeX/RIS/CSL-JSON** | Wciągnięcie bazy źródeł z formatów bibliograficznych | Parser BibTeX (`nickng/bibtex`), RIS, CSL-JSON |
| **Pobranie metadanych po DOI/ISBN** | Uzupełnienie rekordu z rejestrów CrossRef i OpenLibrary | `net/http`; integracja API sterowana konfiguracją, klucz jawny |
| **Generowanie cytowań i bibliografii** | Sformatowana bibliografia w wybranym stylu (APA, MLA, Chicago) | Przetwarzanie CSL (Citation Style Language); style CSL |
| **Powiązanie źródło ↔ dokument** | Odnośnik między cytowanym źródłem a raportem z Research | `powiazanie_komponentu`; integracja Research |
| **Adnotacje i wypisy** | Zaznaczenia i notatki do źródła przechowywane przy pliku | Warstwa adnotacji; notatka pliku |

### 7.11. Udostępnianie, eksport i publikacja

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Udostępnienie odnośnikiem** | Wygenerowanie odnośnika do pliku lub kolekcji | Konfiguracja udostępniania; token dostępu jawny |
| **Eksport zbiorczy** | Wiele plików jako jedno archiwum z zachowaniem struktury | `archive/zip` |
| **Eksport metadanych** | Wywóz metadanych repozytorium do CSV/JSON/XML | `encoding/csv`, `encoding/json`, `encoding/xml` |
| **Publikacja kolekcji jako katalogu** | Statyczny katalog przeglądowy kolekcji w HTML, do pracy bez sieci | Generator statyczny; szablon HTML |
| **Paczka do przekazania klientowi** | Wybrane pliki, indeks i podgląd w samodzielnym pakiecie | ZIP + indeks HTML |
| **Eksport do formatu archiwalnego** | Wywóz w PDF/A i BagIt do archiwum zewnętrznego | Zdolności rozdz. 7.8 |

### 7.12. Operacje AI na repozytorium

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Zapytania do repozytorium (RAG)** | „Podsumuj tę kolekcję”, „porównaj te dwa pliki” na treści zasobów | Kanał modelu; indeks wektorowy (rozdz. 7.4); Chat Window |
| **Auto-porządkowanie** | Sugestie: duplikaty, brak etykiety, błędna kolekcja | Kanał modelu; reguły audytu; karta sugestii |
| **Auto-metadane** | Opis, słowa kluczowe, streszczenie i typ dla nowych plików | Kanał modelu; ekstraktor treści |
| **Wsadowa klasyfikacja** | Sklasyfikowanie całej kolekcji nieuporządkowanych plików | Kanał modelu; pasek zbiorczych operacji; pętla wykonawcza |
| **Odpowiedź z cytowaniem źródła** | Odpowiedź Wykonawcy wskazuje konkretny plik i fragment | Indeks; odnośniki do artefaktów |
| **Wykrywanie treści wrażliwej** | Oznaczenie plików z danymi osobowymi lub poufnymi do przeglądu | Kanał modelu i NER; etykieta „poufne”; sugestia bez blokad |

### 7.13. Automatyzacje, reguły i integracje

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Reguły napływu** | „Pliki z modułu Research z etykietą X → kolekcja Y” | Silnik reguł; integracja Automations |
| **Wyzwalacze zdarzeń** | Akcja przy dodaniu, zmianie lub archiwizacji pliku | Zdarzenia repozytorium; integracja Automations |
| **Powiązania modułów** | Jawne spięcie Studio/Research/Browser → Library | `powiazanie_komponentu`; zakres Integracje |
| **Project Library** | Odpowiednik Library ograniczony do jednego projektu Workspace | Ten sam wzorzec; poziom zasięgu „projekt” |
| **Webhook i API repozytorium** | Zewnętrzny dostęp do zasobów przez API z jawnym kluczem | Endpoint REST; token z okna konfiguracji |
| **Harmonogram konserwacji** | Cyklicznie: reindeksacja, weryfikacja integralności, przeliczenie kolekcji inteligentnych | Integracja Automations |

### 7.14. Analityka, raporty i audyt repozytorium

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Pulpit stanu repozytorium** | Liczba plików, rozmiar, rozkład typów, moduły pochodzenia | Agregaty na `artefakt` |
| **Raport wykorzystania etykiet** | Statystyka użycia etykiet i kolekcji, etykiety nieużywane | Liczniki `etykieta_artefaktu` |
| **Dziennik audytu** | Kto, co i kiedy: dostęp, zmiana, archiwizacja, przywrócenie | Log zdarzeń append-only |
| **Raport duplikatów i higieny** | Zestawienie duplikatów, plików osieroconych i uszkodzonych | Zdolności rozdz. 7.7 |
| **Raport retencji** | Pliki zbliżające się do końca okresu przechowywania | Reguły retencji |
| **Historia dostępu do pliku** | Kiedy i skąd plik był otwierany | Log dostępu na `artefakt` |

### 7.15. Zależności techniczne o charakterze wiążącym

- **MuPDF (`go-fitz`)** obsługuje renderowanie PDF i miniatur stron. Ścieżką zamienną jest proces zewnętrzny (Ghostscript albo Poppler `pdftoppm`).
- **Tesseract** i **ffmpeg** są zależnościami procesowymi, nie czysto-Go: binaria dostarczane są w obrazie środowiska, a `gosseract` oraz `ffmpeg-go` stanowią ich warstwę wywołań.
- **Bleve** pokrywa wyszukiwanie bez usługi zewnętrznej. Przełączenie na Meilisearch, Typesense albo OpenSearch jest ustawieniem konfiguracyjnym i nie zmienia interfejsu.
- **Magazyn wektorów**: `sqlite-vec` obsługuje wdrożenie jednordzeniowe, `pgvector` wdrożenie oparte o PostgreSQL, Qdrant i Milvus wdrożenia dużej skali.
- **Renderowanie DOCX z pełnym układem** przebiega przez konwersję LibreOffice headless do PDF i dalej przez `go-fitz`; `unioffice` dostarcza treść i uproszczony układ w czystym Go.
- Wybór silnika wyszukiwania, magazynu wektorów, formatu archiwalnego oraz kluczy integracji jest jawny i sterowany z okna konfiguracji; domyślne zachowanie modułu = wykonanie.

---

## 8. Punkty sterowania z okna konfiguracji

Zgodnie z zasadą pełnej konfigurowalności każda zdolność modułu jest sterowalna, klucze są jawne, brak ustawienia oznacza wartość domyślną, a interfejs nie stawia blokad. Odwołania wskazują zakresy Modelu konfiguracji (rozdz. 5) oraz poziomy zasięgu (rozdz. 6).

| Punkt sterowania | Co Operator personalizuje | Zakres konfiguracji | Wartość domyślna |
|---|---|---|---|
| **Odbiór z modułów źródłowych** | Włączenie i wyłączenie automatycznego napływu ze Studio/Research/Browser/Design osobno | 5.8 Integracje (`powiazanie_komponentu`) | Zgodnie z ustanowionym powiązaniem |
| **Reguły napływu** | Definicje „moduł + etykieta → kolekcja”, foldery obserwowane | 5.7 Rozszerzenia / Automations | Brak reguł |
| **Silnik wyszukiwania** | Wybór backendu indeksu (Bleve natywny, Meilisearch, OpenSearch) i zakresu indeksacji | 5.2 Procesy | Bleve natywny |
| **Wyszukiwanie semantyczne** | Włączenie embeddingów, wybór magazynu wektorów (pgvector, Qdrant, sqlite-vec), model embeddingów | 5.4 Zachowanie modeli / 5.8 Integracje | Wyłączone |
| **Klucze integracji zewnętrznych** | Klucze API do CrossRef, OpenLibrary i OCR — jawnie w konfiguracji | 5.8 Integracje | Brak (pole jawne) |
| **Automatyczne metadane AI** | Czy nowe pliki otrzymują auto-opis, tagi i klasyfikację; który kanał modelu | 5.4 Zachowanie modeli | Sugestie, akceptacja ręczna |
| **Schemat metadanych** | Pola Dublin Core oraz pola niestandardowe per typ i kolekcja | 5.13 Komponenty własne | Dublin Core podstawowy |
| **Retencja historii wersji** | Okres przechowywania wersji na poziomie global/środowisko/projekt/sesja | 5.9 Historia (`retencja_historii`) | Bez limitu |
| **Polityki retencji zasobów** | Reguły archiwizacji i przeglądu; twarde usunięcie wymaga potwierdzenia | 5.2 Procesy | Brak polityki (przechowywanie trwałe) |
| **Format archiwalny** | Docelowy format normalizacji (PDF/A-2b, BagIt, PREMIS włączone lub wyłączone) | 5.2 Procesy | Bez normalizacji |
| **Udostępnianie odnośnikiem** | Włączenie funkcji, czas życia i zakres tokenów dostępu | 5.3 Akcje / 5.11 Izolacja | Wyłączone |
| **API repozytorium i webhooki** | Włączenie endpointu REST, token jawny, zakres uprawnień | 5.8 Integracje | Wyłączone |
| **Zakres Library ↔ Project Library** | Współdzielenie wybranej kolekcji między repozytorium centralnym a projektem | 6 Poziomy zasięgu (projekt / para modułów) | Odrębne |
| **Izolacja techniczna dostępu do plików** | Ograniczenie odczytu i zapisu dla sesji lub roli (osiem zakresów) | 5.11 Izolacja / rozdz. 6 macierz | Żaden zakres nieaktywny (pełny dostęp) |
| **Deduplikacja i higiena** | Tryb automatyczny albo na żądanie; próg podobieństwa niemal-duplikatów | 5.2 Procesy | Na żądanie |
| **Harmonogram konserwacji** | Częstotliwość reindeksacji, weryfikacji integralności, przeliczania kolekcji inteligentnych | 5.8 Integracje (Automations) | Wyłączony |
| **Miniatury i transkodowanie** | Rozmiary miniatur, formaty transkodowania podglądu, jakość | 5.2 Procesy | Standardowe |
| **Dziennik audytu** | Zakres logowania (dostęp, zmiany, eksport) i retencja logu | 5.9 Historia | Podstawowy (zmiany) |
| **Pętla wykonawcza modułu** | Liczba równoległych zadań, próg ponowienia po kontroli jakości, zakres logu przebiegu | 5.2 Procesy / 5.9 Historia | Ponowienie jednokrotne, log podstawowy |

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + F` | Wyszukiwanie w repozytorium | Library Explorer |
| `Ctrl/Cmd + Klik` | Zaznaczenie wielokrotne | Library Explorer |
| `Spacja` | Szybki podgląd zaznaczonego pliku | Library Explorer → File Preview |
| `Ctrl/Cmd + Enter` | Wysłanie polecenia | Chat Window |
| `Ctrl/Cmd + .` | Wstrzymanie i wznowienie pętli wykonawczej | Execution Loop Window |
| Ikona `folder` | Kolekcja swobodna | Tags & Collections |
| Ikona `odswiez` | Kolekcja inteligentna (regułowa) | Tags & Collections |
| Ikona `oko` | Podgląd pliku lub wersji | File Preview, Versioning Panel |
| Ikona `archiwum` | Archiwizacja i historia | Library Explorer, Versioning Panel, Metadata & Archive Panel |
| Ikona `pobierz` | Eksport pliku, historii lub paczki archiwalnej | File Preview, Versioning Panel, Metadata & Archive Panel |
| Ikona `gwiazdka` | Etykieta kamienia milowego | Versioning Panel |
| Ikona `tarcza` | Weryfikacja integralności (fixity) | Metadata & Archive Panel |
| Ikona `kosz` | Usunięcie (odwracalne) | Library Explorer |

*Koniec dokumentu. Moduł Library — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
