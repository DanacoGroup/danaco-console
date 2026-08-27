# Moduł Studio — dokumentacja projektowa

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
| **Tytuł** | Moduł Studio — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować): katalog okien, funkcji, narzędzi i zależności modułu |
| **Środowiska dostępności** | TalkIn, WorkSpace |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 3, 4, 9.1, 12, Załącznik A) · Specyfikacja modułów 1.1 (rozdz. 4.1) · Specyfikacja okien operacyjnych 1.1 (rozdz. 2.4, 3, 5, 6.1, 10) · System wizualny (rozdz. 6, 8, 10) · Model danych 1.1 (rozdz. 8, 16, 17) · Model konfiguracji (rozdz. 4, 5, 6) · Izolacja i zależności (rozdz. 1–6) |
| **Zasada nadrzędna** | Zero blokad, klucze jawne, pełna kompozycyjność i pełna konfigurowalność z okna konfiguracji; domyślne zachowanie modułu = wykonanie. Każda funkcja działa bez wymuszonej konfiguracji, a każde ograniczenie jest jawnym, odwracalnym ustawieniem Operatora |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - 2.1. [Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
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

Studio jest modułem zaawansowanej pracy z tekstem, dokumentami i treścią. Integruje w jednym oknie roboczym trzy czynności dotąd rozproszone między osobne narzędzia: edycję treści, wsparcie AI działające bezpośrednio na zaznaczonym fragmencie oraz porównywanie wersji dokumentu w czasie. Moduł nie jest edytorem tekstu z dopiętym czatem — jest jedną przestrzenią, w której edytor, operacje kontekstowe AI i historia wersji współdzielą ten sam dokument w czasie rzeczywistym.

Studio stanowi kompletne stanowisko pracy z dokumentami: obejmuje wczytanie dowolnego obsługiwanego formatu, cyfryzację skanów przez OCR, konwersje między formatami, redakcję wspieraną AI, porównania różnicowe, wersjonowanie, eksport oraz zabezpieczenie i podpis dokumentu. Jeden moduł pokrywa cały obszar pracy dokumentowej, bez sięgania po osobne aplikacje desktopowe i usługi sieciowe.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Studio |
|---|---|
| Redaktorzy i copywriterzy | Przepisywanie, zmiana stylu, dopracowanie tonu dużych partii tekstu |
| Analitycy i konsultanci | Przekształcanie notatek roboczych w dokumenty finalne, streszczenia długich materiałów |
| Prawnicy i specjaliści dokumentacji | Praca na wersjach dokumentu z pełną historią zmian i porównaniem redakcji |
| Zespoły wielojęzyczne | Praca dwujęzyczna przy współdzieleniu operacji kontekstowych AI z modułem Translate |
| Twórcy raportów końcowych | Redakcja materiału przekazanego z modułu Research |

### 1.3. Po co — wartość modułu

| Problem klasycznego rozproszenia narzędzi | Rozwiązanie w Studio |
|---|---|
| Edytor tekstu i okno czatu AI to dwie osobne aplikacje | Studio Editor i Chat Window w jednej karcie sesji, nad tym samym dokumentem |
| Brak podglądu, co dokładnie zmieniło AI | Diff/Grep Panel pokazuje różnicę przed/po dla każdej operacji |
| Utrata wcześniejszych redakcji przy nadpisaniu pliku | Session Repository gromadzi każdą wersję ze znacznikiem czasu |
| Ręczne przełączanie się między formatem roboczym a wynikowym | Preview Window renderuje na żywo format docelowy |
| Skan lub zdjęcie dokumentu wymaga osobnego programu OCR | Ingest/OCR Panel prowadzi cyfryzację wewnątrz sesji modułu |

### 1.4. Granica tematyczna — zakres modułu

| Obszar | Zakres pokrycia w Studio |
|---|---|
| Formaty tekstowe i dokumentowe | PDF, DOCX/DOC/ODT/RTF, TXT, Markdown/MDX, HTML, LaTeX, EPUB, arkusze (XLSX/ODS/CSV), prezentacje (odczyt), e-mail (EML/MSG odczyt), obrazy dokumentów (PNG/JPG/TIFF) jako wejście OCR |
| Wczytywanie i cyfryzacja | Import lokalny, z Library, ze schowka, z URL, ze skanera i kamery; OCR skanów i zdjęć, rozpoznanie układu, prostowanie skosu i czyszczenie obrazu |
| Redakcja i tworzenie | Edytor WYSIWYG i tekstowy, style typografii, tabele, obrazy, przypisy, spisy treści, szablony dokumentów, składanie z bloków |
| Operacje kontekstowe AI | Korekta, styl, ton, streszczenia, ekstrakcja, przepisanie strukturalne, tłumaczenie, wzbogacanie, weryfikacja, pytania i odpowiedzi nad dokumentem (RAG dokumentowy) |
| Porównania i wyszukiwanie | Diff tekstowy i wizualny, Grep/regex w dokumencie i w historii, wyszukiwanie semantyczne, porównanie treści z materiałem źródłowym |
| Wersjonowanie | Historia sesyjna, etykiety, gałęzie, przywracanie, eksport historii, znaczniki autora (`uzytkownik`/`model`) |
| Konwersje i eksport | Każdy-do-każdego w granicach obsługiwanych formatów, druk, podgląd wydruku, znak wodny, redakcja poufności, kompresja, scalanie i dzielenie PDF |
| Bezpieczeństwo dokumentu | Szyfrowanie PDF, podpis cyfrowy, usuwanie metadanych, redakcja poufności (trwałe zamazanie), kontrola dostępu na poziomie artefaktu — każde jako jawne, odwracalne ustawienie |

### 1.5. Granica tematyczna — poza zakresem modułu

| Poza zakresem | Uzasadnienie | Właściwy moduł |
|---|---|---|
| Praca nad kodem źródłowym, repozytoria Git | Studio niedostępny w CodeStudio; edycja kodu to inna dyscyplina | Developer, Terminal |
| Generowanie i edycja grafiki koncepcyjnej, ilustracji, brandingu | Studio korzysta z gotowych zasobów wizualnych, ale ich nie tworzy | Design |
| Trwałe repozytorium wiedzy całej organizacji, katalogowanie, kolekcje | Studio wersjonuje w obrębie sesji; magazyn trwały należy do innego modułu | Library |
| Wieloźródłowe badania, benchmarki, budowa raportu z surowych źródeł | Studio redaguje gotowy materiał, nie prowadzi badania | Research |
| Zaawansowane, wielozadaniowe tłumaczenie z zarządzaniem terminologią i pamięcią tłumaczeń | Studio tłumaczy fragmenty kontekstowo; pełna lokalizacja należy do innego modułu | Translate |
| Arkusze jako narzędzie analityczne (pełne modele, tabele przestawne, wykresy analityczne) | Studio odczytuje i edytuje arkusze jako dokumenty, nie jest kalkulatorem analitycznym | Workspace/Research + analiza danych |

Granica jest miękka i konfigurowalna: powiązania `Studio ◄──► Translate`, `Studio ───► Library`, `Studio ◄──► Research`, `Design ───► Studio` pozwalają wywołać sąsiedni moduł bez opuszczania sesji, gdy Operator je ustanowi.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn albo WorkSpace
        ▼
ŚRODOWISKO (TalkIn / WorkSpace) ── boczna nawigacja modułów
        │
        ▼
MODUŁ: STUDIO ─────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Studio Editor ·
  Tools Panel · Diff/Grep Panel · Session Repository ·
  Preview Window · Ingest/OCR Panel
        │
        ▼
KARTA SESJI — własny układ, historia i kontekst
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 5.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | TalkIn, WorkSpace |
| Środowisko CodeStudio | Moduł niedostępny — praca nad kodem należy do Developer/Terminal |
| Komponent własny | Studio nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 8 (łącznie z Chat Window i Execution Loop Window) |
| Typ pracy | Sesyjna, jednodokumentowa lub wielodokumentowa w równoległych kartach |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Rola w module | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja (wspólne wszystkim modułom) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu | 1 | Widoczne bez interakcji od chwili otwarcia modułu |
| 2 | Execution Loop Window | Komunikacja (kanał Koordynator ↔ Wykonawca) | Kolumna sąsiadująca z Chat Window, otwierana | Pętla wykonawcza: dekompozycja zlecenia dokumentowego, kolejka zadań, kontrola jakości, sterowanie przebiegiem | 2 | Znacznik przebiegu pętli w pasku kontekstu Chat Window; otwiera się samoczynnie w chwili uruchomienia zlecenia wielozadaniowego |
| 3 | Studio Editor | Okno edycyjne | Prawa kolumna, dominująca (największa powierzchnia modułu) | Treść dokumentu, punkt wejścia pracy | 1 | Widoczne bez interakcji od chwili otwarcia modułu |
| 4 | Tools Panel | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne | Operacje kontekstowe AI na dokumencie | 3 | Element zbiorczy `Operacje ▼` w pasku Studio Editor, pasek pływający zaznaczenia, skrót poleceniowy `/` w Chat Window |
| 5 | Diff/Grep Panel | Podgląd i porównanie | Kolumna boczna, otwierana jako rozszerzenie boczne | Porównanie wersji, wyszukiwanie w treści | 3 | Pozycja `Porównaj` w menu `Narzędzia ▼`, kliknięcie wyniku operacji AI, kliknięcie pozycji wersji w Session Repository |
| 6 | Session Repository | Repozytorium i biblioteka | Kolumna boczna wąska, otwierana jako rozszerzenie boczne | Historia wersji bieżącej sesji | 3 | Kliknięcie znacznika wersji w pasku statusu Studio Editor, pozycja `Historia wersji` w menu `Narzędzia ▼` |
| 7 | Preview Window | Podgląd i porównanie | Kolumna boczna lub zakładka obszaru roboczego | Podgląd wyniku w formacie docelowym | 2 | Znacznik formatu docelowego w pasku kontekstu Studio Editor; przełącznik podglądu zwija się po zamknięciu okna |
| 8 | Ingest/OCR Panel | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne | Wczytywanie wsadowe, cyfryzacja skanów, rozpoznanie układu | 3 | Pozycja `Wczytaj i cyfryzuj` w menu `Operacje ▼`; panel otwiera się samoczynnie, gdy wejściem sesji jest skan, zdjęcie lub zestaw plików |

```
 Makieta zbiorcza — Moduł Studio          Dostępność: TalkIn, WorkSpace
 Stan spoczynku interfejsu — widoczna warstwa 1 i zwinięte wyzwalacze
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window               │ Studio Editor              │
  nawigacja │ Użytkownik ↔ Wykonawca    │ [Operacje ▼][Narzędzia ▼] ⋮│
  modułów   │                           │ ────────────────────────── │
  (poza     │ historia rozmowy ·        │ treść dokumentu            │
  zakresem  │ strumień odpowiedzi       │ PDF·DOCX·TXT·MD            │
  tego      │                           │                            │
  dokumentu)│                           │                            │
            │ ───────────────────────   │                            │
            │ [ 📎 ][ / ] pole poleceń  │ zapisano · słów · wersja 7 │
            │ [Fable 5][Ultra][⇄ pętla] │ [Format: PDF]              │
 ═══════════════════════════════════════════════════════════════════════════
   kolumna    lewa kolumna, stała         prawa kolumna, dominująca
   nawigacji  pełna wysokość              (warstwa 1)
              (warstwa 1)

   Warstwy 2–4 pozostają zwinięte: znacznik [⇄ pętla] otwiera Execution Loop
   Window, znaczniki kontekstu otwierają selektory modelu i formatu,
   elementy zbiorcze [Operacje ▼] i [Narzędzia ▼] oraz menu kebab ⋮ otwierają
   Tools Panel, Diff/Grep Panel, Session Repository i Ingest/OCR Panel jako
   kolumny boczne.
```

### 2.1. Warstwy widoczności w module

Moduł Studio realizuje regułę stopniowego ujawniania funkcjonalności: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Pełny arsenał modułu — dwanaście rodzin funkcji, warsztat PDF, cyfryzacja OCR, wersjonowanie z gałęziami i zabezpieczenie dokumentu — istnieje w architekturze modułu i pozostaje niewidoczny w interfejsie do chwili wystąpienia potrzeby użycia. Liczba obsługiwanych formatów, operacji kontekstowych i paneli nie wpływa na postrzeganą prostotę okna roboczego.

**Warstwa 1 — zawsze widoczna.** Chat Window w lewej kolumnie wraz z polem poleceń i strumieniem odpowiedzi, Studio Editor jako okno wiodące z obszarem treści dokumentu, pasek statusu dokumentu (status zapisu, liczba słów, numer wersji), pasek kontekstu ze znacznikami modelu, wykonawcy i formatu docelowego oraz wskaźniki stanu wykonania (wskaźnik pracy Wykonawcy, wskaźnik przebiegu pętli). Warstwa ta zajmuje ponad 80% powierzchni modułu.

**Warstwa 2 — widoczna na żądanie.** Selektor kanału i modelu bieżącej karty, przełącznik trybu pracy edytora (WYSIWYG / Markdown / źródło / nakładka OCR), przełącznik zakresu operacji (zaznaczenie / cały dokument / wszystkie karty), selektor formatu podglądu i profilu eksportu, Execution Loop Window oraz Preview Window. Elementy te wywołuje pojedyncze kliknięcie ikony, przełącznika lub znacznika kontekstowego w pasku kontekstu (`[Fable 5] [Ultra] [Format: PDF] [⇄ pętla]`); po użyciu element zwija się samoczynnie.

**Warstwa 3 — rozwinięcia kontekstowe.** Tools Panel z pełnym katalogiem operacji kontekstowych AI, Diff/Grep Panel, Session Repository i Ingest/OCR Panel, pasek pływający zaznaczenia, menu „Warsztat PDF", menu „Zabezpiecz", menu wersji `⋯`, panel Znajdź/Zamień, filtry różnic i filtry historii. Wywołanie następuje przez elementy zbiorcze `Operacje ▼` i `Narzędzia ▼`, menu kebab (⋮), menu hamburger (☰), menu kontekstowe zaznaczenia, panel popover lub listę rozwijaną. Po zamknięciu panel znika całkowicie z przestrzeni roboczej.

**Warstwa 4 — funkcje eksperckie.** Numeracja Bates i formularze AcroForm, szyfrowanie i uprawnienia PDF, podpis cyfrowy i jego weryfikacja, redakcja poufności i usuwanie metadanych, scalanie gałęzi dokumentu z rozwiązywaniem konfliktów per fragment, korekta rozpoznania OCR na warstwie tekstowej, wybór silnika OCR i progu pewności rozpoznania, operacje wsadowe i łańcuchy operacji na wielu kartach, definiowanie własnych promptów jako nazwanych narzędzi, eksport paczki redakcyjnej oraz diagnostyka przebiegu pętli wykonawczej. Dostęp prowadzi przez polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarkę funkcji, tryb administracyjny albo konfigurację roli. Użytkownik podstawowy nie widzi tych elementów.

**Zasada jednego kliknięcia.** Każda funkcja modułu ukryta w warstwach 2–4 jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji Studio głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny okna roboczego, nie utrudnia dostępu do funkcji.

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (okno wspólne)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — kanał Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agent lub system wykonawczy). Stanowi centralny punkt pracy użytkownika w module Studio i podstawowy mechanizm sterowania wszystkimi procesami dokumentowymi: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Okno zajmuje to samo miejsce układu we wszystkich modułach i środowiskach platformy — lewą kolumnę obszaru roboczego.

**Zawartość i pełny arsenał funkcji.**

- Historia rozmowy z Wykonawcą w kontekście bieżącego dokumentu Studio Editor.
- Pole wprowadzania poleceń z obsługą Markdown, wklejania plików i przeciągnij-upuść.
- Strumień odpowiedzi przekazywany na żywo kanałem WebSocket (odpowiedź buduje się token po tokenie); strumień jest trwały po stronie serwera.
- Wybór kanału modelu bieżącej karty (API / CLI / SSH / HTTP) oraz szybka zmiana modelu bazowego z poziomu menu kontekstowego okna.
- Odwołania do zaznaczonego fragmentu dokumentu jako przedmiotu polecenia („zastosuj do zaznaczenia").
- Wywoływanie operacji z Tools Panel bezpośrednio z poziomu wiadomości (skrót poleceniowy `/`).
- Pytania i odpowiedzi nad dokumentem z cytowaniem miejsca w treści (RAG dokumentowy, funkcja C12).
- Dyktowanie polecenia głosem i odczyt odpowiedzi na głos (funkcja K7).
- Cytowanie fragmentów Diff/Grep w treści odpowiedzi (blok `.dn-kod`).
- Załączanie plików do wiadomości (obsługiwane formaty jak w Studio Editor).
- Historia poleceń (strzałka górna/dolna do poprzednich promptów).
- Regeneracja odpowiedzi, edycja własnego polecenia i ponowne wysłanie, rozgałęzianie wątku odpowiedzi.
- Kopiowanie odpowiedzi, wstawienie odpowiedzi bezpośrednio do Studio Editor w miejscu kursora.
- Zatwierdzanie i przerywanie działań zleconych Wykonawcy.
- Wskaźnik pracy Wykonawcy i wskaźnik długości kontekstu bieżącej rozmowy.

**Makieta tekstowa.**

```
┌─ Chat Window ─────────────────────┐
│ Użytkownik ↔ Wykonawca          ⋮ │
├───────────────────────────────────┤
│  Historia rozmowy (przewijana)     │
│   ┌─ Użytkownik ────────────────┐  │
│   │ „popraw styl zaznaczonego   │  │
│   │  akapitu"                   │  │
│   └─────────────────────────────┘  │
│   ┌─ Wykonawca (strumień) ──────┐  │
│   │ tekst odpowiedzi budowany   │  │
│   │ token po tokenie …          │  │
│   └─────────────────────────────┘  │
├───────────────────────────────────┤
│ [ 📎 ] [ / ]                       │
│ Pole poleceń ………… [ Wyślij ▶ ]     │
│ [Fable 5] [Ultra] [⇄ pętla]        │
└───────────────────────────────────┘
   lewa kolumna, stała, pełna
   wysokość obszaru roboczego
   stan spoczynku: warstwa 1 oraz
   znaczniki kontekstowe i menu ⋮
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor kanału modelu | Rozwijana lista w nagłówku | Wybór/zmiana modelu bazowego bieżącej karty | Mały przycisk z etykietą + strzałka | domyślny · rozwinięty · ładowanie (przełączanie modelu) | Otwiera listę kanałów modelu; wybór zmienia model dla kolejnych wiadomości | Nagłówek Chat Window | 2 | Kliknięcie znacznika modelu w pasku kontekstu; lista zwija się samoczynnie po wyborze |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia do Wykonawcy | Duże pole na końcu kolumny, rozciągliwe w pionie | domyślny · fokus · z załącznikiem · błąd (pusta wysyłka) | Enter wysyła, Shift+Enter nowa linia | Zamknięcie kolumny Chat Window | 1 | Widoczne bez interakcji |
| Przycisk „Wyślij" | Przycisk ikonowy `.dn-btn-ikona` (wyslij) | Wysłanie polecenia | Mała ikona 36×36 px, akcent złoty gdy pole niepuste | domyślny · ładowanie | Wysyła polecenie, czyści pole, uruchamia strumień odpowiedzi; przy pustym polu kliknięcie sygnalizuje brak treści komunikatem, bez wysyłki | Obok pola poleceń | 1 | Widoczny bez interakcji |
| Przycisk załącznika | Ikona (spinacz) | Dołączenie pliku do wiadomości | Mała ikona | domyślny · aktywny (plik dołączony) | Otwiera okno wyboru pliku lub akceptuje przeciągnięcie | Lewa strona pola poleceń | 1 | Widoczny bez interakcji |
| Dymek wiadomości użytkownika | Karta wiadomości wyrównana do prawej | Prezentacja wysłanego polecenia | Średni blok tekstu z awatarem `.dn-awatar` | domyślny | Edytowalny (ikona ołówek przy najechaniu) | Historia rozmowy | 1 | Widoczny bez interakcji |
| Dymek odpowiedzi Wykonawcy | Karta wiadomości wyrównana do lewej | Prezentacja odpowiedzi Wykonawcy | Średni blok tekstu, treść może zawierać `.dn-kod` | strumieniowanie (kursor migający) · kompletny · błąd generowania | Po najechaniu ujawnia pasek akcji (Wstaw / Kopiuj / Regeneruj / Przerwij) | Historia rozmowy | 1 | Widoczny bez interakcji |
| Przycisk „Wstaw do dokumentu" | Przycisk `--zarys`, mały | Przeniesienie treści odpowiedzi do Studio Editor | Mały przycisk tekstowy w pasku akcji dymka | domyślny · najechanie | Wstawia treść w miejscu kursora edytora lub w miejscu zaznaczenia | Pasek akcji dymka odpowiedzi | 3 | Pasek akcji dymka ujawniany po najechaniu na odpowiedź |
| Przycisk „Przerwij" | Przycisk `--duch`, mały | Zatrzymanie trwającego działania Wykonawcy | Mały przycisk tekstowy | ukryty (brak działania) · aktywny | Przerywa strumień i zadanie, pozostawiając dokument w stanie sprzed operacji | Pasek akcji dymka odpowiedzi | 3 | Pasek akcji dymka ujawniany po najechaniu na odpowiedź w trakcie strumienia |
| Wskaźnik strumienia | Migający kursor / `.dn-spinner` | Sygnalizacja generowania odpowiedzi na żywo | Mały, subtelny | ukryty · aktywny | Znika po zakończeniu strumienia | Koniec treści generowanej odpowiedzi | 1 | Widoczny bez interakcji w trakcie generowania |
| Skrót „/ operacje" | Menu podpowiedzi poleceń | Szybkie wywołanie operacji Tools Panel z poziomu czatu | Rozwijana lista przy polu poleceń | ukryty · rozwinięty (po wpisaniu „/") | Wstawia szablon polecenia operacji do pola | Przy polu poleceń, wyzwalane znakiem „/" | 3 | Wpisanie znaku „/" w polu poleceń |

---

### 3.2. Execution Loop Window

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — kanał Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Pętla wykonawcza odrębna per karta sesji; przebieg powiązany z dokumentem bieżącej karty |

Execution Loop Window prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą realizującym zadania. Okno odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów modułu Studio. Zlecenie dokumentowe wydane w Chat Window (na przykład „przygotuj wersję finalną raportu": cyfryzacja skanu, korekta, streszczenie zarządcze, eksport PDF) Koordynator rozkłada tu na zadania jednostkowe, przydziela je Wykonawcy, kontroluje jakość wyniku każdego zadania i decyduje o ponowieniu.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie dokumentowe i jego dekompozycja na zadania (wczytanie i OCR, operacja kontekstowa AI, wygenerowanie różnicy, zapis wersji, konwersja, eksport).
- Kolejka zadań i stan każdego zadania (oczekuje, w realizacji, kontrola jakości, ponowienie, zakończone, przerwane).
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, z pełnym zapisem przebiegu.
- Wyniki kontroli jakości zadania — zgodność wyniku operacji AI z zakresem zaznaczenia, poprawność warstwy tekstowej po OCR, kompletność renderu eksportu — oraz decyzje o ponowieniu zadania.
- Wskaźniki przebiegu pętli: liczba zadań w kolejce, liczba ponowień, czas realizacji, zużycie kontekstu modelu.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia w toku.
- Powiązanie zadania z artefaktem i wersją, do których się odnosi — kliknięcie zadania otwiera odpowiadającą mu pozycję w Session Repository lub blok różnicy w Diff/Grep Panel.
- Przebieg łańcuchów operacji i operacji wsadowych na wielu dokumentach (funkcje K3, K4) jako zadania jednej pętli.

**Makieta tekstowa.**

```
┌─ Execution Loop Window ─────────⋮─┐
│ Zlecenie: „wersja finalna raportu"│
│ Koordynator ↔ Wykonawca           │
├───────────────────────────────────┤
│ Kolejka zadań                     │
│  ✓ 1. OCR skanu (18 stron)        │
│  ✓ 2. Korekta całości             │
│  ⟳ 3. Streszczenie zarządcze      │
│  ○ 4. Render PDF + eksport        │
├───────────────────────────────────┤
│ Komunikaty sterujące            ▼ │
├───────────────────────────────────┤
│ Przebieg: 4 zadania · 1 ponowienie│
│ [ Wstrzymaj ][ Wznów ][ Przerwij ]│
└───────────────────────────────────┘
   kolumna sąsiadująca z Chat Window
   okno warstwy 2 — otwierane znacznikiem
   [⇄ pętla]; strumień komunikatów
   sterujących i korekta zlecenia
   pozostają zwinięte (warstwa 3)
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok tytułowy z treścią zlecenia | Wskazanie zlecenia prowadzonego przez pętlę | Średni blok tekstu, wyróżnienie złote | aktywne · wstrzymane · zakończone · przerwane | Kliknięcie rozwija pełną treść zlecenia i jego kontekst | Góra kolumny | 2 | Widoczny po otwarciu Execution Loop Window znacznikiem [⇄ pętla] |
| Pozycja zadania | Wiersz kolejki `.dn-karta--pozycja` | Reprezentuje jedno zadanie dekompozycji | Średni wiersz z ikoną stanu i etykietą | oczekuje · w realizacji (`.dn-spinner`) · kontrola jakości · ponowienie · zakończone · przerwane | Kliknięcie otwiera powiązaną wersję w Session Repository lub różnicę w Diff/Grep Panel | Kolejka zadań | 2 | Widoczna po otwarciu Execution Loop Window |
| Strumień komunikatów sterujących | Lista wiadomości Koordynator ↔ Wykonawca | Podgląd wymiany sterującej | Blok tekstu monospace `.dn-kod`, przewijany | strumieniowanie · kompletny | Przewijanie; kopiowanie pojedynczego komunikatu | Środek kolumny | 3 | Rozwinięcie sekcji komunikatów w Execution Loop Window |
| Znacznik kontroli jakości | Plakietka `.dn-plakietka` przy zadaniu | Wynik kontroli jakości zadania | Mała pigułka z ikoną | zaliczone · odrzucone (ponowienie) | Kliknięcie ujawnia kryterium i uzasadnienie decyzji | Prawa krawędź pozycji zadania | 2 | Widoczny przy pozycji zadania; kliknięcie ujawnia kryterium i uzasadnienie (warstwa 3) |
| Wskaźniki przebiegu pętli | Wiersz liczbowy | Liczba zadań, ponowień, czas, zużycie kontekstu | Mały pasek tekstowy `.dn-tekst-3` | aktualizowany na żywo | — (informacyjny) | Zamknięcie kolumny, nad sterowaniem | 1 | Widoczne bez interakcji jako znacznik stanu wykonania w pasku kontekstu Chat Window |
| Przyciski sterowania przebiegiem | Zestaw trzech przycisków | Wstrzymanie, wznowienie, przerwanie pętli | Małe przyciski `--zarys`, w rzędzie | domyślny · aktywny · nieosiągalny stan ukryty | Wstrzymanie zatrzymuje kolejkę po bieżącym zadaniu; wznowienie kontynuuje; przerwanie kończy pętlę i pozostawia wersje już zapisane | Zamknięcie kolumny | 2 | Widoczne po otwarciu Execution Loop Window |
| Przycisk „Koryguj zlecenie" | Przycisk `.dn-btn--zloty`, pełna szerokość kolumny | Zmiana treści zlecenia w toku pętli | Średni przycisk CTA | domyślny · ładowanie | Otwiera pole edycji zlecenia; zatwierdzenie powoduje ponowną dekompozycję zadań niezrealizowanych | Zamknięcie kolumny | 3 | Menu kebab (⋮) Execution Loop Window |

---

### 3.3. Studio Editor

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne |
| Waga wizualna | Prawa kolumna, dominująca (największa powierzchnia modułu) |
| Izolacja domyślna | Treść dokumentu odrębna per karta sesji; udostępnienie przez Library konfigurowalne |

**Zawartość i pełny arsenał funkcji.**

*Formaty i wczytywanie.*
- Wczytanie dokumentu z Library, z urządzenia lokalnego (przeciągnij-upuść lub wybór pliku), z URL, ze schowka albo utworzenie dokumentu pustego.
- Obsługa formatów: PDF, DOCX/DOC/ODT/RTF, TXT, Markdown/MDX, HTML, LaTeX, EPUB, XLSX/ODS/CSV, EML/MSG; renderowanie WYSIWYG dla Markdown i DOCX, widok tekstowy dla TXT.
- Import ze schowka z automatycznym rozpoznaniem struktury (nagłówki, listy, tabele).
- Przełącznik trybu pracy: WYSIWYG / Markdown / źródło / nakładka OCR dla skanów.

*Formatowanie.*
- Pasek narzędzi formatowania: pogrubienie, kursywa, podkreślenie, przekreślenie, nagłówki H1–H6, cytat blokowy, lista punktowana i numerowana, lista zadań, tabela, blok kodu, link, obraz, separator poziomy.
- Style akapitu i motywy typografii dokumentu, spójne z systemem wizualnym platformy.
- Tabele zaawansowane: scalanie komórek, nagłówki, sortowanie, suma kolumny.
- Przypisy dolne i końcowe, bibliografia i cytowania w stylach APA, MLA, Chicago.
- Numeracja stron, rozdziałów i ilustracji oraz spis treści generowany automatycznie dla dokumentów wielostronicowych.
- Wzory matematyczne w notacji LaTeX oraz diagramy renderowane z opisu tekstowego.

*Zaznaczenie i operacje kontekstowe.*
- Zaznaczenie dowolnego fragmentu jako przedmiotu operacji AI — pasek pływający pojawia się nad zaznaczeniem z najczęstszymi operacjami (korekta, przepisz, streść, zmień styl).
- Zaznaczenie całego dokumentu jednym poleceniem.
- Kotwice, zakładki i minimapa nagłówków do szybkiej nawigacji.

*Redakcja zespołowa.*
- Komentarze przypięte do fragmentu, wątki komentarzy, oznaczenie „rozwiązane".
- Śledzenie zmian z autorem wstawienia i usunięcia, akceptacja i odrzucenie pojedynczo.

*Grep i wyszukiwanie lokalne.*
- Wyszukaj i zamień (z obsługą wyrażeń regularnych) w treści bieżącego dokumentu, z podglądem przed zatwierdzeniem i cofnięciem.
- Podświetlenie wszystkich wystąpień wzorca, nawigacja między trafieniami.

*Współpraca z pozostałymi oknami.*
- Wywołanie Tools Panel bezpośrednio z paska pływającego zaznaczenia.
- Automatyczne odłożenie migawki do Session Repository przy każdej zaakceptowanej zmianie.
- Podgląd na żywo w Preview Window przy przełączeniu widoku kolumnowego.
- Przekazanie skanu lub obrazu do Ingest/OCR Panel i przyjęcie warstwy tekstowej z powrotem.

*Praca skupiona i liczniki.*
- Tryb skupienia z ukryciem kolumn bocznych i centrowaniem wiersza aktywnego.
- Liczniki słów, znaków, zdań, akapitów, stron, czasu czytania i terminów unikalnych.

*Zapis i eksport.*
- Zapis ręczny i autozapis (interwał konfigurowalny z okna konfiguracji).
- Eksport do PDF, DOCX, Markdown, TXT, HTML, EPUB, XLSX, druk.
- Wysłanie bieżącej wersji do Library jednym poleceniem.

**Makieta tekstowa.**

```
┌─ Studio Editor ───────────────────────────────────────────┐
│ [ Operacje ▼ ] [ Narzędzia ▼ ]                          ⋮ │
├───────────────────────────────────────────────────────────┤
│  Tytuł dokumentu                                          │
│  ¶ Treść dokumentu — edytowalna, WYSIWYG lub tekstowa…    │
│                                                            │
│                                                            │
│                                                            │
├───────────────────────────────────────────────────────────┤
│ Status: zapisano 12:04 · słów: 1 284 · wersja 7           │
│ [Tryb: WYSIWYG] [Format: PDF]                             │
└───────────────────────────────────────────────────────────┘
   prawa kolumna, dominująca
   stan spoczynku: obszar treści i pasek statusu (warstwa 1),
   znaczniki kontekstowe trybu i formatu (warstwa 2),
   elementy zbiorcze i menu ⋮ (warstwa 3); pasek pływający
   zaznaczenia oraz Znajdź/Zamień ujawniają się na żądanie
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pasek narzędzi formatowania | Rząd przycisków ikonowych | Formatowanie zaznaczonego tekstu | Pasek poziomy, wysokość jednego wiersza, stały na początku kolumny | domyślny · aktywny (styl zastosowany do zaznaczenia podświetla ikonę) | Kliknięcie zmienia formatowanie zaznaczenia lub pozycji kursora; bez otwartego dokumentu kliknięcie sygnalizuje brak treści do sformatowania | Otwarcie kolumny Studio Editor | 3 | Element zbiorczy `Operacje ▼` w pasku Studio Editor |
| Przełącznik trybu pracy | Rozwijana lista trybów | Wybór WYSIWYG / Markdown / źródło / nakładka OCR | Mały przycisk z etykietą | WYSIWYG (domyślny) · Markdown · źródło · nakładka OCR | Przerysowuje obszar treści w wybranym trybie, zachowując pozycję kursora | Pasek narzędzi | 2 | Kliknięcie znacznika trybu w pasku kontekstu; zwija się po wyborze |
| Obszar treści | Edytowalny obszar dokumentu | Główna przestrzeń pisania i odczytu | Duży, dominujący blok, przewijany w pionie | pusty stan (`.dn-pusty-stan`, brak dokumentu) · edycja · tylko odczyt (podgląd) | Wpisywanie tekstu, zaznaczanie, wklejanie | Środek kolumny | 1 | Widoczny bez interakcji |
| Pasek pływający zaznaczenia | Mały kontekstowy pasek narzędzi | Szybki dostęp do operacji AI na zaznaczeniu | Mały, unoszący się nad zaznaczonym tekstem | ukryty · widoczny (po zaznaczeniu ≥ 1 znaku) | Kliknięcie operacji otwiera Tools Panel z wypełnionym kontekstem lub wykonuje operację wprost | Nad zaznaczonym fragmentem | 3 | Zaznaczenie fragmentu treści (menu kontekstowe zaznaczenia) |
| Znajdź/Zamień | Ikona lupy + panel rozwijany | Wyszukiwanie i zamiana wzorca w treści | Mała ikona w pasku narzędzi, rozwija pole `.dn-input` | zwinięty · rozwinięty · trafienia podświetlone | Enter przechodzi do kolejnego trafienia; „Zamień wszystko" prezentuje podgląd zmian przed zatwierdzeniem | Prawy róg paska narzędzi | 3 | Pozycja `Znajdź/Zamień` w menu `Narzędzia ▼` lub skrót klawiszowy |
| Pasek statusu dokumentu | Wiersz informacyjny | Liczba słów, status zapisu, numer wersji | Mały, niski pasek, tekst pomocniczy | zapisano · zapisywanie (`.dn-spinner`) · niezapisane zmiany | Kliknięcie numeru wersji otwiera Session Repository na tej pozycji | Zamknięcie kolumny Studio Editor | 1 | Widoczny bez interakcji |
| Selektor stylu akapitu | Rozwijana lista | Zastosowanie predefiniowanego stylu typografii | Mały przycisk z etykietą | domyślny · rozwinięty | Zmienia styl akapitu, w którym znajduje się kursor | Pasek narzędzi | 3 | Element zbiorczy `Operacje ▼`, grupa stylów typografii |
| Minimapa nagłówków | Wąska kolumna zwijana | Nawigacja po nagłówkach dokumentu | Wąska kolumna, chowana i wysuwana | zwinięta · rozwinięta | Kliknięcie pozycji przewija dokument do nagłówka | Lewa krawędź obszaru treści | 2 | Kliknięcie uchwytu przy lewej krawędzi obszaru treści; zwija się po wyborze nagłówka |
| Znacznik śledzenia zmian | Oznaczenie wstawienia lub usunięcia w treści | Prezentacja zmiany z autorem | Podkreślenie lub przekreślenie w kolorze autora | wstawienie · usunięcie · zaakceptowana · odrzucona | Kliknięcie ujawnia autora i przyciski akceptacji lub odrzucenia | Obszar treści | 2 | Widoczny po włączeniu śledzenia zmian znacznikiem kontekstowym trybu redakcji |
| Komentarz redakcyjny | Karta komentarza przypięta do fragmentu | Dyskusja redakcyjna nad fragmentem | Mała karta w kolumnie marginesu | otwarty · rozwiązany | Kliknięcie rozwija wątek i pole odpowiedzi | Prawa krawędź obszaru treści | 3 | Pozycja `Komentarz` w menu kontekstowym zaznaczenia; karta znika po zamknięciu wątku |

---

### 3.4. Tools Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Zestaw operacji wspólny; wynik operacji zapisuje się w kontekście bieżącej sesji |

**Zawartość i pełny arsenał funkcji — operacje kontekstowe AI.**

| Kategoria operacji | Operacje |
|---|---|
| Korekta i jakość tekstu | Korekta ortograficzno-gramatyczna, korekta interpunkcji, sprawdzenie spójności terminologii, wykrycie powtórzeń, wykrycie strony biernej, kontrola czytelności (Flesch, Gunning fog, długość zdań) |
| Przekształcenia stylu | Zmiana stylu (formalny / nieformalny / techniczny / marketingowy), zmiana tonu (neutralny / perswazyjny / empatyczny / stanowczy), skrócenie, rozwinięcie treści, uproszczenie języka, podniesienie rejestru językowego, parafraza z wariantami do wyboru |
| Streszczenia i ekstrakcja | Streszczenie jednozdaniowe, streszczenie akapitowe, wypunktowanie kluczowych tez, wyodrębnienie listy działań, wyodrębnienie decyzji, terminów i encji (osoby, kwoty, daty), wygenerowanie tytułu i podtytułów, wygenerowanie meta-opisu |
| Tłumaczenia | Tłumaczenie zaznaczenia lub całości na wybrany język, wersja dwujęzyczna równoległa (z użyciem Translate, gdy powiązanie skonfigurowane) |
| Przepisywanie strukturalne | Przekształcenie akapitu na listę, przekształcenie listy na tabelę, przekształcenie prozy na punkty, wygenerowanie spisu treści, wygenerowanie streszczenia zarządczego na początku dokumentu |
| Wzbogacanie treści | Rozwinięcie akapitu o kontekst, wygenerowanie przykładów, wygenerowanie definicji, wygenerowanie pytań kontrolnych, wygenerowanie kontrargumentów |
| Weryfikacja | Wskazanie fragmentów wymagających źródła, oznaczenie potencjalnych nieścisłości, porównanie z materiałem źródłowym z Library, pytania i odpowiedzi nad dokumentem z cytowaniem miejsca |
| Operacje własne | Zapis własnych promptów jako nazwanych narzędzi wielokrotnego użytku, dostępnych na równi z operacjami fabrycznymi |

Panel udostępnia także łańcuchy operacji — sekwencje wykonywane jednym poleceniem (na przykład korekta → streszczenie → eksport) — oraz operacje wsadowe obejmujące wiele kart i plików naraz. Przebieg łańcucha i wsadu prowadzi Execution Loop Window.

**Makieta tekstowa.**

```
┌─ Tools Panel ────────────────────⋮─┐
│ [ Zakres: zaznaczenie ▼ ]          │
├────────────────────────────────────┤
│  Korekta i jakość                ▼ │
│  Styl i ton                      ▼ │
│  Streszczenia                    ▼ │
│  Tłumaczenia                     ▼ │
│  Struktura                       ▼ │
│  Weryfikacja                     ▼ │
│  Operacje własne                 ▼ │
├────────────────────────────────────┤
│        [ Uruchom operację ]        │
└────────────────────────────────────┘
   kolumna boczna, rozszerzenie boczne
   panel warstwy 3 — wywoływany elementem
   zbiorczym [ Operacje ▼ ]; grupy operacji
   pozostają zwinięte, łańcuchy i operacje
   wsadowe należą do warstwy 4
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Przełącznik zakresu | Zestaw przycisków radiowych | Wybór, czy operacja dotyczy zaznaczenia, całego dokumentu czy wszystkich kart | Mały zestaw `.dn-check` typu radio | domyślny (zaznaczenie, gdy istnieje) · cały dokument · wszystkie karty | Zmienia zakres wszystkich kolejnych operacji | Góra panelu | 2 | Widoczny po otwarciu Tools Panel; zwija się do znacznika zakresu po wyborze |
| Grupa operacji | Nagłówek kategorii z listą pozycji | Porządkuje operacje tematycznie | Nagłówek `.dn-karta-tytul` + lista `.dn-karta--pozycja` | zwinięta · rozwinięta | Kliknięcie nagłówka zwija/rozwija grupę | Ciało panelu | 3 | Rozwinięcie kategorii w Tools Panel |
| Pozycja operacji | Pojedynczy wiersz listy | Wybór konkretnej operacji AI | Mały wiersz z ikoną i etykietą | domyślny · najechanie · wybrany | Kliknięcie zaznacza operację i (dla operacji z parametrem) rozwija podmenu | Wewnątrz grupy operacji | 3 | Rozwinięcie grupy operacji |
| Podmenu z listą rozwijaną | Np. „Zmień styl ▾" | Wybór wariantu operacji (np. konkretny styl docelowy) | Mała rozwijana lista przy pozycji | zwinięte · rozwinięte | Wybór wariantu zamyka listę i ustawia parametr operacji | Przy pozycjach parametryzowanych | 3 | Kliknięcie pozycji parametryzowanej — lista rozwijana |
| Pozycja łańcucha operacji | Wiersz listy z ikoną sekwencji | Uruchomienie sekwencji operacji jednym poleceniem | Mały wiersz z licznikiem kroków | domyślny · w realizacji | Przekazuje sekwencję do Execution Loop Window jako zadania jednej pętli | Grupa „Łańcuchy" | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji lub grupa „Łańcuchy" dostępna po włączeniu operacji zaawansowanych w konfiguracji roli |
| Przycisk „Uruchom operację" | Główny CTA panelu | Wysłanie wybranej operacji do wykonania | Duży przycisk `.dn-btn--zloty`, pełna szerokość kolumny | domyślny · ładowanie | Uruchamia operację na wybranej pozycji; bez wybranej operacji kliknięcie sygnalizuje to komunikatem zamiast uruchomienia. Wynik trafia do Chat Window i Diff/Grep Panel, a przebieg do Execution Loop Window | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Tools Panel |
| Wskaźnik zakresu operacji | Mała etykieta nad listą | Przypomnienie, ile znaków i słów obejmie operacja | Tekst pomocniczy `.dn-tekst-3` | domyślny · aktualizowany na żywo | Aktualizuje się przy każdej zmianie zaznaczenia w edytorze | Nad grupami operacji | 3 | Widoczny po otwarciu Tools Panel |

---

### 3.5. Diff/Grep Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Podgląd i porównanie |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne po operacji AI lub na żądanie |
| Izolacja domyślna | Operuje na wersjach z Session Repository bieżącej sesji |

**Zawartość i pełny arsenał funkcji.**

- Widok różnicowy dwukolumnowy lub scalony, z kolorowym oznaczeniem dodania, usunięcia i zmiany, z granularnością znaku, słowa lub zdania.
- Porównanie dowolnych dwóch wersji z Session Repository, nie tylko sąsiadujących.
- Diff wizualny wyrenderowanego układu PDF — nakładka i podświetlenie zmian graficznych.
- Grep — wyszukiwanie wzorca (w tym wyrażeń regularnych) w treści bieżącej lub w całej historii wersji dokumentu.
- Wyszukiwanie semantyczne — fragmenty znaczeniowo bliskie zapytaniu, nie tylko dosłowne trafienia.
- Porównanie ze źródłem: zestawienie dokumentu roboczego z materiałem wejściowym z Research lub Library ze wskazaniem rozbieżności.
- Filtrowanie różnic: tylko dodania, tylko usunięcia, tylko zmiany formatowania.
- Akceptacja lub odrzucenie zmiany fragmentarycznie (per akapit i zdanie) albo całościowo; decyzja zasila warstwę śledzenia zmian edytora.
- Adnotacje przy zmianach — komentarz użytkownika przypisany do konkretnej różnicy.
- Eksport widoku różnicowego jako osobny dokument (raport redakcji).
- Statystyka zmiany: liczba dodanych, usuniętych i zmienionych słów oraz znaków.
- Nawigacja klawiaturowa między kolejnymi różnicami.

**Makieta tekstowa.**

```
┌─ Diff/Grep Panel ────────────────⋮─┐
│ Porównaj: [ W6 ▼ ] → [ W7 ▼ ]      │
│ Tryb: [ scalony ▼ ]  Grep: [ 🔍 ]  │
├────────────────────────────────────┤
│  − usunięty fragment tekstu        │
│  + dodany fragment tekstu          │
│  ▪ fragment zmieniony (format)     │
│                                     │
├────────────────────────────────────┤
│ +48 słów · −12 słów · 3 zmiany     │
└────────────────────────────────────┘
   kolumna boczna, rozszerzenie boczne
   panel warstwy 3 — wywoływany z menu
   [ Narzędzia ▼ ]; akceptacja, odrzucenie
   i adnotacja ujawniają się po kliknięciu
   bloku różnicy, eksport raportu zmian
   należy do warstwy 4 (menu ⋮)
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektory wersji | Dwie rozwijane listy | Wybór porównywanej pary wersji | Małe przyciski z etykietą numeru wersji | domyślny (ostatnie dwie wersje) · rozwinięty | Zmiana przelicza widok różnicowy | Nagłówek panelu | 3 | Widoczne po otwarciu Diff/Grep Panel z menu `Narzędzia ▼` |
| Przełącznik trybu | Zakładki `.dn-zakladki--pigulki` | Wybór widoku: dwukolumnowy / scalony / wizualny | Mały segmentowany przełącznik | scalony (domyślny) · dwukolumnowy · wizualny | Przełącza układ prezentacji różnic | Nagłówek panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |
| Pole Grep | Pole wyszukiwania z przełącznikami | Wyszukiwanie wzorca w treści i historii | Pole `.dn-input` + ikony-przełączniki (regex, wielkość liter, tryb semantyczny) | puste · z wynikami · brak trafień | Enter uruchamia wyszukiwanie, podświetla trafienia w treści | Nagłówek panelu | 3 | Widoczne po otwarciu Diff/Grep Panel lub skrót klawiszowy wyszukiwania |
| Blok różnicy | Fragment tekstu z oznaczeniem koloru | Prezentacja pojedynczej zmiany | Średni blok, kolor tła zależny od typu zmiany | dodanie (zielone tło) · usunięcie (czerwone, przekreślenie) · zmiana formatu (żółte obramowanie) | Kliknięcie zaznacza blok do akceptacji lub odrzucenia | Ciało panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |
| Przycisk „Akceptuj zmianę" | Mały przycisk `--zarys` | Zatwierdzenie pojedynczej różnicy do wersji roboczej | Mały przycisk przy bloku różnicy | domyślny · najechanie | Wprowadza zmianę do Studio Editor | Przy każdym bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Przycisk „Odrzuć zmianę" | Mały przycisk `--duch` | Pominięcie różnicy | Mały przycisk przy bloku różnicy | domyślny · najechanie | Zachowuje treść sprzed zmiany w dokumencie roboczym | Przy każdym bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Ikona adnotacji | Ikona dymka (odpowiedz) | Dodanie komentarza do zmiany | Mała ikona | brak adnotacji · z adnotacją (wypełniona) | Otwiera małe pole tekstowe komentarza | Przy bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Przycisk „Eksportuj raport zmian" | Przycisk `--zarys`, pełna szerokość kolumny | Zapis widoku różnicowego jako dokumentu | Średni przycisk | domyślny · ładowanie | Generuje dokument raportu redakcji i udostępnia go do pobrania lub do Library | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Diff/Grep Panel albo polecenie języka naturalnego w Chat Window |
| Pasek statystyki | Wiersz podsumowania | Liczbowe podsumowanie zmian | Mały pasek tekstowy | domyślny, aktualizowany na żywo | — (informacyjny) | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |

---

### 3.6. Session Repository

| Aspekt | Wartość |
|---|---|
| Typologia | Repozytorium i biblioteka |
| Waga wizualna | Kolumna boczna wąska, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Historia wersji odrębna per karta sesji; narasta w toku sesji, nic nie jest usuwane bez decyzji użytkownika |

**Zawartość i pełny arsenał funkcji.**

- Chronologiczna lista wersji dokumentu ze znacznikiem czasu i autorem zmiany (`uzytkownik` lub `model`).
- Etykietowanie wersji własną nazwą (na przykład „wersja do akceptacji klienta") i oznaczanie wersji kluczowych.
- Podgląd dowolnej wersji bez opuszczania modułu oraz podgląd wariantów obok siebie.
- Przywrócenie wcześniejszej wersji jako bieżącej według polityki append-only — przywrócenie tworzy nową wersję na końcu listy, nic nie znika.
- Porównanie dowolnych dwóch wersji jednym kliknięciem — otwiera Diff/Grep Panel z wypełnioną parą.
- Gałęzie dokumentu: utworzenie wariantu z wybranej wersji jako punktu startowego oraz scalanie wariantów z rozwiązywaniem konfliktów per fragment.
- Eksport pojedynczej wersji lub całej historii jako archiwum ZIP z manifestem.
- Eksport paczki redakcyjnej: dokument, historia, raport zmian i adnotacje w jednym archiwum przekazania.
- Filtrowanie listy wersji: tylko zmiany użytkownika, tylko zmiany AI, wersje oznaczone etykietą.
- Wskaźnik liczby zmian od ostatniej etykiety.
- Wyszukiwanie w treści konkretnej wersji (przekazywane do Diff/Grep Panel jako Grep).
- Odwołanie do konkretnej wersji dokumentu w obrębie platformy, generowane jednym poleceniem.

**Makieta tekstowa.**

```
┌─ Session Repository ─⋮─┐
│ Filtr: [ wszystkie ▼ ] │
├────────────────────────┤
│ ● Wersja 7 · 12:04   ⋯ │
│   model                │
│   „korekta stylu"      │
├────────────────────────┤
│ ○ Wersja 6 · 11:47   ⋯ │
│   użytkownik           │
├────────────────────────┤
│ ○ Wersja 5 · 11:20   ⋯ │
│   🏷 „do akceptacji"    │
└────────────────────────┘
   wąska kolumna boczna
   panel warstwy 3 — wywoływany
   kliknięciem znacznika wersji
   w pasku statusu edytora; akcje
   pozycji rozwija menu ⋯,
   rozgałęzianie, eksport historii
   i paczka redakcyjna należą do
   warstwy 4
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Filtr historii | Rozwijana lista | Zawężenie listy wersji wg autora lub etykiety | Mały przycisk z etykietą | domyślny (wszystkie) · aktywny filtr | Przelicza widoczną listę wersji | Góra panelu | 3 | Widoczny po otwarciu Session Repository |
| Pozycja wersji | Wiersz karty `.dn-karta--pozycja` | Reprezentuje jedną zapisaną wersję dokumentu | Średni wiersz z kropką stanu, znacznikiem czasu, autorem | bieżąca (kropka wypełniona, wyróżnienie złote) · archiwalna (kropka pusta) · najechanie | Rozwija zestaw akcji (Podgląd / Przywróć / Porównaj / Rozgałęź) | Lista główna panelu | 3 | Widoczna po otwarciu Session Repository znacznikiem wersji w pasku statusu |
| Etykieta wersji | Mała plakietka `.dn-plakietka` | Własna nazwa nadana wersji przez użytkownika | Mała pigułka z ikoną (gwiazdka) | brak etykiety · z etykietą | Kliknięcie pozwala edytować treść etykiety | Pod znacznikiem czasu w pozycji wersji | 3 | Menu `⋯` pozycji wersji |
| Menu „⋯" | Ikona rozwijanego menu (wiecej) | Dodatkowe akcje na wersji (etykieta, eksport, odwołanie do wersji, usunięcie z listy lokalnie) | Mała ikona 24 px | domyślny · rozwinięte | Otwiera menu kontekstowe pozycji | Prawa krawędź pozycji wersji | 3 | Kliknięcie ikony menu przy pozycji wersji |
| Przycisk „Przywróć" | Przycisk `--zarys`, mały | Ustawienie wybranej wersji jako bieżącej | Mały przycisk tekstowy | domyślny · najechanie | Tworzy nową wersję na końcu listy z treścią wersji przywróconej | Zestaw akcji pozycji wersji | 3 | Rozwinięcie zestawu akcji pozycji wersji |
| Przycisk „Porównaj" | Przycisk `--duch`, mały | Otwarcie Diff/Grep Panel dla tej wersji względem bieżącej | Mały przycisk tekstowy | domyślny | Otwiera Diff/Grep Panel z wypełnioną parą wersji | Zestaw akcji pozycji wersji | 3 | Rozwinięcie zestawu akcji pozycji wersji |
| Przycisk „Rozgałęź" | Przycisk `--duch`, mały | Utworzenie wariantu dokumentu od tej wersji | Mały przycisk tekstowy | domyślny | Tworzy gałąź dokumentu i otwiera ją jako bieżącą treść edytora | Zestaw akcji pozycji wersji | 4 | Menu `⋯` pozycji wersji po włączeniu pracy na gałęziach w konfiguracji roli albo polecenie języka naturalnego w Chat Window |
| Przycisk „Eksportuj historię" | Przycisk `--zarys`, pełna szerokość kolumny | Pobranie całej historii wersji jako archiwum | Średni przycisk | domyślny · ładowanie | Generuje archiwum z manifestem i udostępnia plik do pobrania | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Session Repository albo polecenie języka naturalnego w Chat Window |
| Przycisk „Paczka redakcyjna" | Przycisk `--zarys`, pełna szerokość kolumny | Spakowanie dokumentu, historii, raportu zmian i adnotacji | Średni przycisk | domyślny · ładowanie | Generuje archiwum przekazania i udostępnia je do pobrania lub do Library | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Session Repository albo polecenie języka naturalnego w Chat Window |

---

### 3.7. Preview Window

| Aspekt | Wartość |
|---|---|
| Typologia | Podgląd i porównanie |
| Waga wizualna | Kolumna boczna lub zakładka obszaru roboczego obok Studio Editor |
| Izolacja domyślna | Renderuje bieżący stan dokumentu; brak własnego stanu trwałego |

**Zawartość i pełny arsenał funkcji.**

- Renderowanie dokumentu w formacie docelowym (typografia, style, numeracja stron) w czasie zbliżonym do rzeczywistego, po każdej zaakceptowanej zmianie w Studio Editor.
- Przełącznik formatu podglądu: PDF, HTML, Markdown wyrenderowany, DOCX, EPUB.
- Widok jednostronicowy i widok ciągły (przewijany) dla dokumentów wielostronicowych.
- Powiększenie i dopasowanie do szerokości lub strony.
- Widok kolumnowy: Studio Editor i Preview Window w sąsiadujących kolumnach, z synchronicznym przewijaniem.
- Eksport bezpośrednio z podglądu (PDF, DOCX, XLSX, EPUB, HTML, druk) oraz przekazanie do Library.
- Profile eksportu — nazwane zestawy ustawień (na przykład „PDF do druku", „DOCX dla klienta").
- Podgląd wydruku: marginesy, orientacja, nagłówek, stopka, numeracja, skala.
- Znak wodny „wersja robocza" sterowany ustawieniem konfiguracyjnym podglądu, bez wpływu na treść dokumentu.
- Warsztat PDF dostępny z menu eksportu: scalanie i dzielenie, porządkowanie stron, kompresja i optymalizacja, znak wodny i stemplowanie, formularze AcroForm, numeracja Bates, ekstrakcja obrazów i załączników, drzewo zakładek.
- Zabezpieczenie dokumentu dostępne z menu eksportu: szyfrowanie i uprawnienia PDF, podpis cyfrowy i jego weryfikacja, redakcja poufności, usuwanie metadanych, wykrywanie danych wrażliwych.

**Makieta tekstowa.**

```
┌─ Preview Window ─────────────────⋮─┐
│ [ Format: PDF ▼ ]  [ Widok ▼ ]     │
├────────────────────────────────────┤
│                                     │
│    ┌───────────────────────────┐   │
│    │ wyrenderowana strona      │   │
│    │ dokumentu w formacie      │   │
│    │ docelowym                 │   │
│    └───────────────────────────┘   │
│                                     │
├────────────────────────────────────┤
│ [ Widok kolumnowy ]  [ Operacje ▼ ]│
└────────────────────────────────────┘
   kolumna boczna lub zakładka
   okno warstwy 2 — otwierane znacznikiem
   formatu docelowego; eksport i wysłanie
   do Library prowadzi element zbiorczy
   [ Operacje ▼ ] (warstwa 3), warsztat PDF
   i zabezpieczenie dokumentu należą do
   warstwy 4
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor formatu | Rozwijana lista | Wybór formatu renderowania podglądu | Mały przycisk z etykietą | domyślny (PDF) · rozwinięty | Przerenderowuje dokument w wybranym formacie | Nagłówek okna | 2 | Kliknięcie znacznika formatu docelowego w pasku kontekstu; zwija się po wyborze |
| Przełącznik widoku | Para przycisków `.dn-zakladki--pigulki` | Strona pojedyncza / przewijanie ciągłe | Mały segmentowany przełącznik | strona (domyślny) · ciągły | Zmienia sposób prezentacji wielostronicowego dokumentu | Nagłówek okna | 2 | Widoczny po otwarciu Preview Window |
| Suwak powiększenia | Rozwijana lista procentowa | Zmiana skali podglądu | Mały przycisk z etykietą procentową | 100% (domyślny) · niestandardowy | Skaluje renderowaną stronę | Nagłówek okna | 2 | Widoczny po otwarciu Preview Window |
| Obszar renderowania | Wyrenderowana treść dokumentu | Wizualna prezentacja wyniku | Duży, centralny blok | ładowanie (`.dn-spinner`) · wyrenderowany · pusty stan (brak dokumentu) | Przewijanie, powiększanie; brak edycji bezpośredniej | Środek okna | 2 | Widoczny po otwarciu Preview Window znacznikiem formatu docelowego |
| Przycisk „Widok kolumnowy" | Przycisk `--zarys` | Zestawienie Studio Editor i Preview w sąsiadujących kolumnach | Mały przycisk tekstowy | wyłączony · włączony (podświetlony) | Dzieli obszar roboczy modułu na dwie sąsiadujące kolumny z synchronicznym przewijaniem | Zamknięcie kolumny okna | 2 | Widoczny po otwarciu Preview Window |
| Przycisk „Eksportuj" | Przycisk `--zloty` z rozwijanym menu formatów i profili | Pobranie dokumentu w formacie docelowym | Średni przycisk CTA | domyślny · rozwinięte menu · ładowanie eksportu | Generuje plik według wybranego profilu i uruchamia pobranie | Zamknięcie kolumny okna | 3 | Pozycja `Eksport` w menu `Operacje ▼` lub w Preview Window |
| Menu „Warsztat PDF" | Rozwijane menu operacji PDF | Scalanie, dzielenie, porządkowanie stron, kompresja, stemplowanie, formularze, Bates, zakładki | Średni przycisk z rozwinięciem | domyślny · rozwinięte · ładowanie operacji | Wykonuje operację na bieżącym dokumencie i zapisuje wynik jako nową wersję | Zamknięcie kolumny okna | 4 | Wyszukiwarka funkcji, skrót klawiszowy albo polecenie języka naturalnego w Chat Window; menu widoczne po włączeniu operacji zaawansowanych w konfiguracji roli |
| Menu „Zabezpiecz" | Rozwijane menu operacji ochronnych | Szyfrowanie, podpis, redakcja poufności, usuwanie metadanych | Średni przycisk z rozwinięciem | domyślny · rozwinięte · ładowanie operacji | Nakłada wybrane zabezpieczenie i zapisuje wynik jako nową wersję | Zamknięcie kolumny okna | 4 | Wyszukiwarka funkcji albo polecenie języka naturalnego w Chat Window; menu widoczne po włączeniu operacji zaawansowanych w konfiguracji roli |
| Przycisk „Wyślij do Library" | Przycisk `--zarys` | Przekazanie bieżącej wersji do modułu Library | Mały przycisk tekstowy | domyślny · wysłano (potwierdzenie `.dn-toast`) | Tworzy artefakt w Library, gdy powiązanie skonfigurowane | Zamknięcie kolumny okna | 3 | Pozycja `Wyślij do Library` w menu `Operacje ▼` lub w Preview Window |

---

### 3.8. Ingest/OCR Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Kolejka wczytywania odrębna per karta sesji; wynik trafia do dokumentu bieżącej karty |

Panel prowadzi cyfryzację materiału wejściowego. Otwiera się, gdy wejściem są skany, zdjęcia lub zestawy plików; przy dokumentach tekstowych pozostaje zwinięty i praca przebiega bez niego.

**Zawartość i pełny arsenał funkcji.**

- Kolejka wczytywania wsadowego: folder, archiwum ZIP, wiele plików naraz, jednakowy potok normalizacji.
- Wybór silnika OCR i języków rozpoznawania oraz progu pewności rozpoznania.
- Wstępne czyszczenie obrazu przed OCR: prostowanie skosu, usuwanie szumu, binaryzacja, przycięcie marginesów, podniesienie kontrastu.
- Rozpoznanie układu: kolumny, tabele, nagłówki, stopki, przypisy — odtworzenie struktury logicznej dokumentu.
- Podgląd warstwy tekstowej nałożonej na skan, z możliwością korekty rozpoznania przed przyjęciem do edytora.
- Import z URL z oczyszczeniem strony z nawigacji i reklam oraz import wiadomości e-mail z załącznikami.
- Przekazanie wyniku do Studio Editor jako dokumentu roboczego oraz zapis pierwszej wersji w Session Repository.

**Makieta tekstowa.**

```
┌─ Ingest/OCR Panel ───────────────⋮─┐
│ [ Źródło: Pliki ▼ ]                │
├────────────────────────────────────┤
│ Kolejka wczytywania                │
│  ✓ umowa-skan-01.tif   OCR gotowy  │
│  ⟳ umowa-skan-02.tif   OCR 62%     │
│  ○ umowa-skan-03.tif   oczekuje    │
├────────────────────────────────────┤
│ Ustawienia rozpoznawania         ▼ │
├────────────────────────────────────┤
│      [ Przekaż do edytora ]        │
└────────────────────────────────────┘
   kolumna boczna, rozszerzenie boczne
   panel warstwy 3 — wywoływany pozycją
   `Wczytaj i cyfryzuj` w menu [ Operacje ▼ ];
   silnik OCR, języki, czyszczenie obrazu,
   rozpoznanie układu i korekta rozpoznania
   pozostają zwinięte jako warstwa 4
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor źródła | Zestaw przycisków źródeł | Wybór drogi wczytania materiału | Mały segmentowany przełącznik | pliki (domyślny) · folder · URL · skaner · e-mail | Otwiera właściwy dialog wczytania lub pole adresu | Góra panelu | 3 | Widoczny po otwarciu Ingest/OCR Panel pozycją `Wczytaj i cyfryzuj` w menu `Operacje ▼` |
| Pozycja kolejki | Wiersz `.dn-karta--pozycja` z nazwą pliku | Reprezentuje jeden materiał wejściowy | Średni wiersz z ikoną stanu i postępem | oczekuje · przetwarzanie (`.dn-spinner` z procentem) · gotowy · błąd odczytu | Kliknięcie otwiera podgląd warstwy tekstowej tej pozycji | Kolejka wczytywania | 3 | Widoczna po otwarciu Ingest/OCR Panel |
| Selektor silnika OCR | Rozwijana lista | Wybór silnika rozpoznawania | Mały przycisk z etykietą | lokalny (domyślny) · usługa zewnętrzna | Zmienia silnik dla kolejnych pozycji kolejki | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu albo polecenie języka naturalnego w Chat Window |
| Selektor języków | Rozwijana lista wielokrotnego wyboru | Wybór języków rozpoznawania | Mały przycisk z listą znaczników | pol + eng (domyślny) · zestaw własny | Zmienia zestaw danych językowych dla rozpoznawania | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu albo polecenie języka naturalnego w Chat Window |
| Przełączniki pre-OCR | Zestaw pól wyboru `.dn-check` | Sterowanie czyszczeniem obrazu przed rozpoznaniem | Małe pola wyboru w rzędach | włączone · wyłączone | Zmiana przelicza przetwarzanie pozycji oczekujących | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu |
| Przełączniki rozpoznania układu | Zestaw pól wyboru `.dn-check` | Sterowanie wykrywaniem kolumn, tabel, przypisów | Małe pola wyboru w rzędach | włączone · wyłączone | Zmienia zakres odtwarzanej struktury dokumentu | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu |
| Podgląd warstwy tekstowej | Obraz skanu z nałożonym tekstem | Kontrola jakości rozpoznania | Średni blok z przełącznikiem przezroczystości | warstwa ukryta · warstwa widoczna · korekta | Kliknięcie słowa otwiera pole korekty rozpoznanego tekstu | Środek panelu | 3 | Kliknięcie pozycji kolejki; korekta rozpoznania stanowi warstwę 4 |
| Przycisk „Przekaż do edytora" | Główny CTA panelu | Przyjęcie wyniku cyfryzacji jako dokumentu roboczego | Duży przycisk `.dn-btn--zloty`, pełna szerokość kolumny | domyślny · ładowanie | Otwiera dokument w Studio Editor i zapisuje pierwszą wersję w Session Repository | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Ingest/OCR Panel |

---

## 4. Przepływy pracy

### 4.1. Przepływ podstawowy — redakcja fragmentu z korektą AI

```
 Studio Editor        Chat / Execution Loop      Diff/Grep Panel      Session Repository
 ─────────────        ─────────────────────      ───────────────      ──────────────────
 1. wczytanie
    dokumentu  ───────►
 2. zaznaczenie
    fragmentu  ─────────────► polecenie w Chat Window
                              (korekta / styl)
                                    │
                                    ▼
                              Koordynator dzieli
                              zlecenie na zadania
                              (Execution Loop)
                                    │
                                    ▼
                              Wykonawca realizuje,
                              strumień w Chat Window
                                    │
                                    ▼
 3. ‹ podgląd zmiany ›  ◄─────────────────────────  wygenerowana różnica
                                                             │
 4. akceptacja lub odrzucenie  ──────────────────────────────┘
        │
        ▼
 5. zapis kolejnej wersji  ────────────────────────────────────────────►  nowa pozycja
                                                                             w historii
```

### 4.2. Przepływ cyfryzacji — skan do dokumentu edytowalnego

```
Ingest/OCR Panel
   kolejka wczytywania → czyszczenie obrazu → rozpoznanie tekstu
   → rozpoznanie układu → podgląd warstwy tekstowej
        │
        ▼
Execution Loop Window — każda pozycja kolejki jako zadanie pętli,
   kontrola jakości warstwy tekstowej, ponowienie przy niskiej pewności
        │
        ▼
Studio Editor — dokument roboczy z odtworzoną strukturą
        │
        ▼
Session Repository — pierwsza wersja dokumentu
        │
        ▼
Preview Window — eksport przeszukiwalnego PDF lub edytowalnego DOCX
```

### 4.3. Przepływ rozszerzony — praca dwujęzyczna z modułem Translate

```
Studio Editor (dokument źródłowy)
        │  zaznaczenie fragmentu do wersji dwujęzycznej
        ▼
Tools Panel ── operacja „Tłumaczenie" ──► współdzielenie operacji kontekstowych AI
        │                                  (powiązanie Studio ◄──► Translate,
        │                                   konfigurowalne — rozdz. 5.3)
        ▼
Translate — Translation Panels (wynik równoległy w językach docelowych)
        │
        ▼
Powrót do Studio Editor jako wersja wielojęzyczna dokumentu
        │
        ▼
Session Repository — zapis wersji z etykietą „wielojęzyczna"
```

### 4.4. Przepływ rozszerzony — redakcja raportu przekazanego z Research

```
Research › Report Builder                Studio › Studio Editor
──────────────────────────               ──────────────────────
 raport końcowy ustaleń      ──────────►   wczytanie jako dokument
                                            źródłowy
                                                  │
                                                  ▼
                                     operacje kontekstowe AI: redakcja,
                                     dopracowanie stylu, streszczenie
                                     zarządcze na wstępie dokumentu
                                                  │
                                                  ▼
                                     Diff/Grep Panel — porównanie
                                     z materiałem wejściowym z Research
                                                  │
                                                  ▼
                                     Preview Window — podgląd finalny
                                                  │
                                                  ▼
                                     eksport lub wysłanie do Library
```

### 4.5. Przepływ zbiorczy — wiele dokumentów równolegle

```
Karta sesji A · Studio          Karta sesji B · Studio          Karta sesji C · Studio
  dokument 1                      dokument 2                      dokument 3
  (kontekst odrębny)               (kontekst odrębny)              (kontekst odrębny)

           domyślnie: pełna izolacja kontekstu między kartami (rozdz. 5.3)
           ustawienie konfiguracyjne: współdzielona pamięć projektu,
                        ustanawiana z okna konfiguracji punktów izolacji

           operacja wsadowa obejmująca wszystkie karty prowadzona jest
           jako jedna pętla w Execution Loop Window
```

---

## 5. Stany, dane i powiązania

### 5.1. Model stanów sesji modułu

```
                    ┌────────────────┐
                    │  Pusta sesja   │  (brak dokumentu)
                    └───────┬────────┘
                            │ wczytanie / cyfryzacja / nowy dokument
                            ▼
                    ┌────────────────┐
        ┌──────────►│  Edycja        │◄─────────────┐
        │           └───────┬────────┘               │
        │                   │ zlecenie operacji AI    │ odrzucenie zmiany
        │                   ▼                         │
        │           ┌────────────────┐                │
        │           │  Pętla         │  (Execution     │
        │           │  wykonawcza    │   Loop Window)  │
        │           └───────┬────────┘                │
        │                   │ wynik zaliczony         │
        │                   ▼                         │
        │           ┌────────────────┐                │
        │           │  Podgląd zmiany│────────────────┘
        │           │  (Diff/Grep)   │
        │           └───────┬────────┘
        │                   │ akceptacja
        │                   ▼
        │           ┌────────────────┐
        └───────────│  Nowa wersja   │──► Session Repository
                    └────────────────┘
```

Stan procesu sesji po stronie serwera jest trwały niezależnie od stanu połączenia klienta — rozłączenie nie przerywa trwającej operacji ani nie zamyka okna; po ponownym połączeniu Studio Editor przywraca dokument w stanie, w jakim został pozostawiony, a Execution Loop Window odtwarza przebieg pętli.

### 5.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Studio |
|---|---|
| `sesja`, `karta_sesji` | Nośnik kontekstu bieżącej pracy nad dokumentem |
| `wiadomosc`, `zalacznik_wiadomosci` | Historia Chat Window i pliki załączone do poleceń |
| `artefakt` | Dokument roboczy jako plik z odwołaniem w bazie |
| `wersja_artefaktu` | Każda pozycja w Session Repository; pole `autor` rozróżnia zmianę `uzytkownik`/`model`; nośnik etykiet, gałęzi i adnotacji |
| `kanal_modelu` | Model bazowy przypisany karcie sesji, wybierany w Chat Window; źródło embeddingów wyszukiwania semantycznego oraz przetwarzania mowy |
| `ustawienie` | Parametry modułu podlegające zasadzie „brak ustawienia = wartość domyślna" (interwał autozapisu, prompty operacji, profile eksportu, mapa skrótów) |
| `profil_izolacji`, `regula_izolacji_kontekstu` | Zakres współdzielenia historii, pamięci i kontekstu między kartami i modułami oraz zakres widoczności artefaktu |
| `powiazanie_komponentu` | Jawne powiązania Studio ◄──► Translate, Studio ───► Library, Studio ◄──► Research, Design ───► Studio |

### 5.3. Izolacja i konfigurowalność — punkty właściwe modułowi Studio

| Punkt izolacji | Stan wyjściowy (domyślny) | Co można skonfigurować | Gdzie |
|---|---|---|---|
| Historia Chat Window | Odrębna per karta sesji | Współdzielenie historii między kartami pracującymi nad powiązanymi dokumentami | Okno konfiguracji punktów izolacji, poziom „karta sesji" lub „moduł" |
| Przebieg Execution Loop Window | Odrębny per karta sesji | Współdzielenie kolejki zadań między kartami przy operacjach wsadowych | Okno konfiguracji punktów izolacji, poziom „moduł" lub „projekt" |
| Pamięć kontekstowa | Odrębna per karta sesji | Współdzielenie pamięci projektu, gdy Studio używany w ramach projektu Workspace | Poziom zasięgu „projekt" |
| Para modułów Studio–Translate | Operacje kontekstowe AI nie są współdzielone domyślnie | Współdzielenie mechanizmu operacji kontekstowych AI przy pracy dwujęzycznej | Poziom zasięgu „para modułów" |
| Repozytorium źródłowe (Library) | Brak automatycznego powiązania | Ustanowienie Library jako źródła wejściowego i celu eksportu artefaktów | Okno konfiguracji, powiązanie komponentu Studio ───► Library |
| Widoczność artefaktu | Zakres sesji | Rozszerzenie widoczności dokumentu do projektu lub zakresu współdzielonego | Okno konfiguracji punktów izolacji |
| Izolacja techniczna procesu sesji (8 zakresów) | Żaden zakres domyślnie nie jest aktywny — proces sesji nie jest ograniczany | Włączenie wybranych zakresów (np. odrębny katalog roboczy) dla sesji redagującej dokumenty poufne | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

Zgodnie z zasadą nadrzędną platformy żaden z powyższych punktów nie jest wymuszony — moduł działa w pełni funkcjonalnie bez jakiejkolwiek konfiguracji izolacji, a każdy punkt jest odwracalną, jawną decyzją Operatora.

### 5.4. Powiązania z innymi modułami

```
       Translate ◄──►  STUDIO  ───►  Library
     (operacje AI              (repozytorium źródłowe
      przy pracy                i cel dla artefaktów)
      dwujęzycznej,
      konfigurowalne)                 ▲
                                       │
                            Research ◄─┘
                    (dalsze przetwarzanie wyników
                     redakcji przy budowie raportów;
                     Studio ◄──► Research w obie strony)

       Design  ───►  Studio   (zasoby wizualne przy redagowaniu dokumentów)
       Automations ◄──► Studio (kolejki wsadowe i łańcuchy operacji)
       Assistant  ◄──► Studio (dyktowanie i odczyt na głos)
       Browser    ───► Studio (import stron WWW jako dokumentów)
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Translate | Współdzielenie operacji kontekstowych AI przy pracy dwujęzycznej | Konfiguracyjne | Tools Panel → Translation Panels |
| Library | Repozytorium źródłowe dokumentów wejściowych i cel dla artefaktów wygenerowanych | Konfiguracyjne | Preview Window / Studio Editor → Library Explorer |
| Research | Dalsze przetwarzanie wyników redakcji przy budowie raportów końcowych | Konfiguracyjne | Report Builder → Studio Editor |
| Design | Wykorzystanie zasobów wizualnych przy redagowaniu dokumentów | Konfiguracyjne | Assets Panel → Studio Editor |
| Automations | Kolejki wczytywania wsadowego, łańcuchy operacji, korespondencja seryjna | Konfiguracyjne | Ingest/OCR Panel / Tools Panel → kolejki Automations |
| Assistant | Dyktowanie polecenia i odczyt dokumentu na głos | Konfiguracyjne | Chat Window → Voice Console |
| Browser | Import strony WWW jako dokumentu roboczego | Konfiguracyjne | Ingest/OCR Panel → Studio Editor |

---

## 6. Scenariusze użycia

**Scenariusz 1 — korekta i akceptacja zmiany.**
Użytkownik wczytuje dokument DOCX do Studio Editor, zaznacza akapit budzący wątpliwości stylistyczne, z paska pływającego wybiera „Korekta". Koordynator rozkłada zlecenie na zadania widoczne w Execution Loop Window, odpowiedź Wykonawcy trafia do Chat Window, a Diff/Grep Panel prezentuje różnicę. Użytkownik akceptuje zmianę fragmentarycznie — tylko dwie z trzech zaproponowanych poprawek — pozostawiając trzecią bez zmian. Session Repository zapisuje nową wersję.

**Scenariusz 2 — cyfryzacja umowy ze skanu.**
Kancelaria wczytuje osiemnastostronicowy skan umowy do Ingest/OCR Panel. Panel prostuje skos, odszumia obrazy i rozpoznaje tekst w językach polskim i angielskim. Execution Loop Window prowadzi kolejkę stron jako zadania pętli i ponawia rozpoznanie dwóch stron o niskiej pewności. Użytkownik koryguje kilka słów w podglądzie warstwy tekstowej, przekazuje wynik do Studio Editor, a z Preview Window eksportuje przeszukiwalny PDF z numeracją Bates.

**Scenariusz 3 — praca dwujęzyczna.**
Zespół lokalizacyjny ustanawia w oknie konfiguracji współdzielenie operacji kontekstowych AI między Studio a Translate na poziomie pary modułów. Redakcja dokumentu źródłowego w Studio Editor odzwierciedla się w tłumaczeniu zaznaczenia bezpośrednio z Tools Panel, bez opuszczania modułu Studio.

**Scenariusz 4 — redakcja raportu z Research.**
Analityk kończy pracę badawczą w module Research; z Report Builder przekazuje skompletowany raport do Studio. W Studio Editor prowadzi ostateczną redakcję językową, generuje streszczenie zarządcze operacją z Tools Panel, porównuje wynik z materiałem wejściowym w Diff/Grep Panel, a następnie eksportuje finalny PDF z Preview Window.

**Scenariusz 5 — praca równoległa nad wieloma dokumentami.**
Redaktor otwiera trzy karty sesji modułu Studio, każdą z osobnym dokumentem klienta. Domyślna izolacja kontekstu zapewnia, że polecenia wydane w jednej karcie nie wpływają na kontekst pozostałych. Dla dwóch kart dotyczących tego samego projektu redaktor włącza współdzielenie pamięci na poziomie „projekt", a operację korekty uruchamia wsadowo — Execution Loop Window prowadzi ją jako jedną pętlę obejmującą oba dokumenty.

**Scenariusz 6 — powrót do wcześniejszej redakcji.**
Po serii zmian klient prosi o powrót do wersji sprzed dwóch dni. Użytkownik odnajduje ją w Session Repository po etykiecie „do akceptacji klienta", otwiera podgląd, a następnie przywraca ją jako wersję bieżącą — operacja tworzy nową pozycję na końcu historii, zachowując wszystkie wersje pośrednie.

**Scenariusz 7 — przekazanie paczki redakcyjnej.**
Po zamknięciu redakcji użytkownik usuwa metadane dokumentu, nakłada podpis cyfrowy z Preview Window i generuje z Session Repository paczkę redakcyjną: dokument finalny, pełną historię wersji, raport zmian i adnotacje w jednym archiwum przekazania.

---

## 7. Katalog funkcji i narzędzi

Katalog obejmuje dwanaście rodzin funkcji. Każda pozycja opisuje nazwę, działanie oraz zależności techniczne (biblioteki, formaty, integracje). Nazwy bibliotek podano jako konkretne komponenty w Go — języku zaplecza platformy — wraz z rezerwą; narzędzia zewnętrzne (LibreOffice, Pandoc, Tesseract) wywoływane są jako procesy pomocnicze przez moduł Terminal lub bezpośrednio przez warstwę serwera.

### 7.1. Rodzina A — Wczytywanie, formaty i cyfryzacja

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| A1 | Uniwersalny importer formatów | Rozpoznaje typ pliku po sygnaturze i wczytuje PDF/DOCX/ODT/RTF/TXT/MD/HTML/EPUB/XLSX/CSV do wspólnego modelu dokumentu | `h2non/filetype` (detekcja MIME); parsery per format; model dokumentu wewnętrzny |
| A2 | Ekstraktor tekstu z PDF | Wyciąga tekst, strukturę akapitów i pozycje z PDF tekstowego, zachowując kolejność czytania | `ledongthuc/pdf` lub `gen2brain/go-fitz` (MuPDF, wyższa jakość układu); format PDF |
| A3 | Silnik OCR skanów i zdjęć | Rozpoznaje tekst z obrazów i PDF-skanów wielojęzycznie, zwraca warstwę tekstową z pozycjami słów | `otiai10/gosseract/v2` → Tesseract (`libtesseract`, `tessdata` — pol, eng, deu…); dla PDF-skanu render przez go-fitz; usługa chmurowa (Google Vision / Azure) sterowana jawnym kluczem w oknie konfiguracji |
| A4 | Wstępne czyszczenie obrazu | Prostuje skos, usuwa szum, binaryzuje, przycina marginesy, zwiększa kontrast przed rozpoznaniem | `disintegration/imaging` lub `anthonynsimon/bild`; formaty PNG/JPG/TIFF |
| A5 | Rozpoznanie układu | Wykrywa kolumny, tabele, nagłówki, stopki, przypisy — odtwarza strukturę logiczną, nie tylko surowy tekst | go-fitz (bloki tekstu z bbox); reguły heurystyczne i model AI |
| A6 | Import ze schowka z rozpoznaniem struktury | Wkleja HTML/RTF/tekst ze schowka i odtwarza nagłówki, listy, tabele | Parser HTML `net/html`; konwersja HTML → model dokumentu |
| A7 | Import z URL i zrzut strony do dokumentu | Pobiera stronę WWW, oczyszcza z nawigacji i reklam, zapisuje jako czytelny dokument | `go-shiori/go-readability`; powiązanie z modułem Browser |
| A8 | Import e-maila jako dokumentu | Wczytuje EML/MSG, wyodrębnia treść, załączniki i nagłówki do dokumentu roboczego | `net/mail`, `emersion/go-message`, `mnako/letters`; formaty EML/MSG |
| A9 | Kolejka wsadowego wczytywania | Wczytuje wiele plików naraz (folder, ZIP) z jednakowym potokiem rozpoznawania i normalizacji | `archive/zip`; kolejka zadań (powiązanie z Automations) |

### 7.2. Rodzina B — Edycja i tworzenie treści

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| B1 | Edytor WYSIWYG dokumentu | Edycja bogata: style znaku i akapitu, nagłówki H1–H6, listy, cytaty, tabele, obrazy, linki, separatory | Edytor po stronie klienta (ProseMirror/TipTap w warstwie web); model dokumentu ↔ serwer |
| B2 | Edytor Markdown z podglądem | Tryb źródłowy Markdown/MDX z podglądem na żywo, tabelami GFM, przypisami, blokami kodu | `yuin/goldmark` (CommonMark + GFM, tabele, footnotes, task-list) do renderu |
| B3 | Szablony dokumentów | Predefiniowane układy (pismo, umowa, raport, notatka, oferta) z polami do wypełnienia | `text/template`/`html/template`; biblioteka szablonów w Library |
| B4 | Style typografii i motywy dokumentu | Zestawy stylów (czcionki, interlinia, marginesy) stosowane jednym kliknięciem, spójne z systemem wizualnym | Definicje stylów w modelu dokumentu; tokeny typografii platformy |
| B5 | Tabele zaawansowane | Wstawianie i edycja tabel: scalanie komórek, nagłówki, sortowanie, suma kolumny | Model tabeli w dokumencie; render do DOCX/PDF/MD |
| B6 | Przypisy, bibliografia, cytowania | Zarządza przypisami dolnymi i końcowymi oraz cytowaniami w stylach APA, MLA, Chicago | Parser stylów cytowania (CSL); baza źródeł (powiązanie z Research/Library) |
| B7 | Spis treści i numeracja automatyczna | Generuje spis treści z nagłówków, numeruje strony, rozdziały, ilustracje | Analiza drzewa nagłówków; render paginacji |
| B8 | Kotwice, zakładki i nawigacja | Zakładki w dokumencie i minimapa nagłówków do szybkiego skoku | Indeks nagłówków; minimapa w oknie edytora |
| B9 | Komentarze i adnotacje redakcyjne | Komentarze przypięte do fragmentu, wątki, oznaczenie „rozwiązane" | Model adnotacji w `wersja_artefaktu`; render do PDF (warstwa adnotacji) |
| B10 | Śledzenie zmian | Rejestruje wstawienia i usunięcia z autorem, akceptacja i odrzucenie pojedynczo | Warstwa zmian nad modelem dokumentu; integracja z rodziną E |
| B11 | Liczniki i statystyki dokumentu | Słowa, znaki, zdania, akapity, strony, czas czytania, terminy unikalne | Analizator tekstu (segmentacja zdań i słów) |

### 7.3. Rodzina C — Operacje kontekstowe AI na treści

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| C1 | Korekta ortograficzno-gramatyczna | Poprawia błędy pisowni, gramatyki i interpunkcji na zaznaczeniu lub całości | `kanal_modelu`; warstwa lokalna `hunspell`/`client9/misspell` jako szybkie rozpoznanie wstępne |
| C2 | Zmiana stylu i rejestru | Przełącza styl (formalny/nieformalny/techniczny/marketingowy) i rejestr językowy | `kanal_modelu`; szablony promptów operacji (konfigurowalne) |
| C3 | Zmiana tonu | Nadaje ton (neutralny/perswazyjny/empatyczny/stanowczy) bez zmiany treści | `kanal_modelu` |
| C4 | Skracanie i rozwijanie | Kondensuje lub rozszerza fragment do zadanej długości lub proporcji | `kanal_modelu`; wskaźnik długości (B11) |
| C5 | Streszczenia wielopoziomowe | Streszczenie jednozdaniowe, akapitowe, wypunktowane, streszczenie zarządcze na wstępie | `kanal_modelu` |
| C6 | Ekstrakcja strukturalna | Wyciąga tezy, działania, decyzje, terminy, encje (osoby, kwoty, daty) | `kanal_modelu`; reguły rozpoznawania encji |
| C7 | Przepisanie strukturalne | Akapit → lista, lista → tabela, proza → punkty, generacja tytułów i podtytułów | `kanal_modelu`; operacje na modelu dokumentu |
| C8 | Parafraza | Przepisuje frazę innym sformułowaniem, z wariantami do wyboru | `kanal_modelu` |
| C9 | Uproszczenie języka | Upraszcza żargon, skraca zdania, podnosi czytelność do zadanego poziomu | `kanal_modelu` + metryki czytelności (C13) |
| C10 | Tłumaczenie kontekstowe | Tłumaczy zaznaczenie lub całość na wybrany język, wersja dwujęzyczna równoległa | `kanal_modelu`; powiązanie `Studio ◄──► Translate` (konfigurowalne) |
| C11 | Wzbogacanie treści | Dodaje kontekst, przykłady, kontrargumenty, pytania kontrolne, definicje | `kanal_modelu` |
| C12 | Pytania i odpowiedzi nad dokumentem (RAG) | Odpowiada na pytania w oparciu o treść dokumentu z cytowaniem miejsca | `kanal_modelu`; indeks wektorowy (`blevesearch/bleve` + embeddingi) lub kontekst pełnotekstowy |
| C13 | Kontrola czytelności i jakości | Wskaźniki (Flesch, Gunning fog, długość zdań), powtórzenia, strona bierna, spójność terminologii | Metryki liczone lokalnie; `kanal_modelu` dla sugestii |
| C14 | Weryfikacja i wskazanie źródeł | Oznacza fragmenty wymagające źródła i potencjalne nieścisłości, porównuje z materiałem z Library | `kanal_modelu`; powiązanie z Library/Research |
| C15 | Operacje własne | Zapis własnych promptów jako nazwanych narzędzi wielokrotnego użytku | Szablony promptów w `ustawienie`; okno konfiguracji |

### 7.4. Rodzina D — Konwersje i eksport

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| D1 | Konwerter uniwersalny dokumentów | Konwertuje między MD/DOCX/HTML/LaTeX/EPUB/RTF/ODT/PDF/TXT w dowolnym kierunku | `pandoc` wywoływany jako proces zewnętrzny (najszersza matryca formatów) |
| D2 | Renderowanie do PDF | Składa dokument do PDF z typografią, paginacją, nagłówkiem i stopką, spisem treści | `go-pdf/fpdf` lub `signintech/gopdf`; ścieżka HTML→PDF przez `chromedp` dla wiernego układu |
| D3 | Eksport do DOCX | Zapisuje dokument jako edytowalny DOCX ze stylami i tabelami | `unidoc/unioffice` (klucz komercyjny jawny) lub `fumiama/go-docx`; rezerwa: Pandoc |
| D4 | Eksport do arkusza | Zapisuje tabele dokumentu jako arkusz XLSX lub CSV | `qax-os/excelize` (XLSX), `encoding/csv` |
| D5 | Eksport do EPUB | Składa dokument wielorozdziałowy w e-book | `bmaupin/go-epub` + goldmark do treści |
| D6 | Eksport do HTML | Renderuje dokument jako samodzielny HTML z osadzonymi stylami i obrazami | goldmark/`html/template`; inline CSS |
| D7 | Podgląd i konfiguracja wydruku | Marginesy, orientacja, nagłówek i stopka, numeracja, skala, podgląd stron | Warstwa layoutu druku; render przez D2 |
| D8 | Konwersja obraz → dokument | Zamienia skan lub zdjęcie w przeszukiwalny PDF albo edytowalny DOCX/TXT | A3 (rozpoznanie tekstu) + D2/D3; warstwa tekstowa nad obrazem |
| D9 | Eksport z profilami | Nazwane profile eksportu (np. „PDF do druku", „DOCX dla klienta") z zapamiętanymi ustawieniami | Profile w `ustawienie` (konfigurowalne) |

### 7.5. Rodzina E — Diff, porównania i wyszukiwanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| E1 | Diff tekstowy słowo po słowie | Różnica dwóch wersji z oznaczeniem dodań, usunięć i zmian, tryb scalony i dwukolumnowy | `sergi/go-diff` (diff-match-patch) lub `hexops/gotextdiff` (Myers) |
| E2 | Diff dowolnej pary wersji | Porównuje nie tylko sąsiednie, lecz dowolne dwie pozycje z historii | E1 + `wersja_artefaktu` (Session Repository) |
| E3 | Diff wizualny PDF i układu | Porównuje wyrenderowany wygląd — nakładka i podświetlenie zmian graficznych | Render obu wersji przez go-fitz → porównanie pikselowe (`imaging`) |
| E4 | Grep i wyszukiwanie regex | Wyszukuje wzorzec (w tym regex) w dokumencie i w całej historii wersji | `regexp` stdlib; podświetlanie trafień |
| E5 | Znajdź i zamień wsadowo | Zamiana wzorca w dokumencie, podgląd przed zatwierdzeniem, cofnięcie | `regexp`; operacja na modelu dokumentu z migawką |
| E6 | Wyszukiwanie semantyczne | Znajduje fragmenty znaczeniowo bliskie zapytaniu, nie tylko dosłowne | `blevesearch/bleve` + embeddingi z `kanal_modelu` |
| E7 | Porównanie ze źródłem | Zestawia dokument roboczy z materiałem wejściowym (Research/Library) i wskazuje rozbieżności | E1 + powiązanie międzymodułowe |
| E8 | Filtrowanie i statystyka różnic | Filtruje różnice (dodania, usunięcia, format), liczy zmienione słowa i znaki | Agregacja wyniku E1 |
| E9 | Adnotacje i eksport raportu zmian | Komentarz przy różnicy; eksport całości różnic jako osobny dokument redakcji | Model adnotacji + D2 (PDF) |

### 7.6. Rodzina F — Wersjonowanie i historia

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| F1 | Historia wersji sesyjna | Chronologiczna lista wersji ze znacznikiem czasu i autorem (`uzytkownik`/`model`) | `wersja_artefaktu`; nic nie znika bez decyzji użytkownika |
| F2 | Etykiety i kamienie milowe | Nazwa własna wersji („do akceptacji klienta"), oznaczenie wersji kluczowych | Pole etykiety w `wersja_artefaktu` |
| F3 | Przywracanie bez utraty | Ustawia starszą wersję jako bieżącą, tworząc nową pozycję; historia pozostaje nienaruszona | F1; polityka append-only |
| F4 | Gałęzie dokumentu | Tworzy równoległy wariant od wybranej wersji (np. dwie redakcje) | Model gałęzi nad `artefakt`; scalanie przez E1 |
| F5 | Scalanie wariantów | Łączy dwie gałęzie z rozwiązywaniem konfliktów per fragment | E1/E2 + interfejs rozwiązywania konfliktów |
| F6 | Filtrowanie i przegląd historii | Filtr: zmiany użytkownika, zmiany AI, wersje etykietowane; podgląd bez opuszczania modułu | Zapytania nad `wersja_artefaktu` |
| F7 | Eksport i archiwizacja historii | Pobiera pojedynczą wersję lub całą historię jako archiwum ZIP z manifestem | `archive/zip`; manifest JSON |
| F8 | Autozapis i migawki | Automatyczny zapis w konfigurowalnym interwale; migawka przy każdej zaakceptowanej zmianie AI | Interwał w `ustawienie` (konfigurowalny) |

### 7.7. Rodzina G — Operacje na PDF

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| G1 | Scalanie i dzielenie PDF | Łączy wiele PDF w jeden, dzieli po stronach, zakresach i zakładkach | `pdfcpu/pdfcpu` (merge, split, trim) |
| G2 | Porządkowanie stron | Zmiana kolejności, obrót, usuwanie, wstawianie, ekstrakcja stron | `pdfcpu` (rotate, remove, insert, extract) |
| G3 | Kompresja i optymalizacja | Zmniejsza rozmiar PDF przez kompresję obrazów i czyszczenie struktury | `pdfcpu` (optimize); rekompresja obrazów przez `imaging` |
| G4 | Znak wodny i stemplowanie | Nakłada znak wodny lub pieczęć tekstową albo graficzną | `pdfcpu` (stamp/watermark) |
| G5 | Formularze PDF | Wypełnianie i tworzenie pól formularza, odczyt i zapis danych AcroForm | `pdfcpu` (form fill/export) lub `unidoc/unipdf` |
| G6 | Numeracja Bates i nagłówki | Dodaje numerację prawną Bates oraz nagłówki i stopki na wszystkich stronach | `pdfcpu` (stamp z licznikiem) |
| G7 | Ekstrakcja obrazów i załączników | Wyciąga obrazy, czcionki i załączniki osadzone w PDF | `pdfcpu` (extract images/attachments) |
| G8 | Zakładki i nawigacja PDF | Tworzy i edytuje drzewo zakładek dokumentu | `pdfcpu` (bookmarks) / `unipdf` |

### 7.8. Rodzina H — Bezpieczeństwo i poufność dokumentu

Każda funkcja rodziny H jest jawnym, odwracalnym ustawieniem Operatora; żadna nie warunkuje pracy nad dokumentem.

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| H1 | Szyfrowanie i hasło PDF | Nakłada i zdejmuje szyfrowanie oraz uprawnienia dokumentu | `pdfcpu` (encrypt/decrypt/permissions) |
| H2 | Podpis cyfrowy i pieczęć | Podpisuje PDF certyfikatem, weryfikuje podpis | `digitorus/pdfsign`; certyfikat użytkownika |
| H3 | Redakcja poufności | Trwale zamazuje wskazane fragmenty wraz z warstwą tekstową pod spodem | Nadpisanie warstwy tekstu + zaczernienie obrazu (pdfcpu + imaging) |
| H4 | Usuwanie metadanych | Czyści autora, historię, dane ukryte i komentarze przed wysyłką | `pdfcpu` (metadata); parsery DOCX i obrazów |
| H5 | Wykrywanie danych wrażliwych | Oznacza w treści numery identyfikacyjne, numery kart, adresy e-mail i dane osobowe do redakcji | Reguły regex + rozpoznawanie encji przez `kanal_modelu` |
| H6 | Kontrola dostępu do artefaktu | Ustawia zakres widoczności dokumentu (sesja, projekt, zakres współdzielony) | `profil_izolacji`, `regula_izolacji_kontekstu`; okno konfiguracji izolacji |

### 7.9. Rodzina I — Arkusze i dane tabelaryczne

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| I1 | Odczyt i edycja arkusza | Otwiera XLSX/ODS/CSV, edytuje komórki, arkusze i proste formuły | `qax-os/excelize` (XLSX), `encoding/csv` |
| I2 | Tabela dokumentu ↔ arkusz | Konwertuje tabelę z dokumentu do arkusza i odwrotnie | Model tabeli (B5) ↔ excelize |
| I3 | Ekstrakcja tabel z PDF | Wykrywa i wyciąga tabele z PDF do arkusza lub CSV | go-fitz (pozycje) + heurystyka siatki; rezerwa: `tabula-java` jako proces |
| I4 | Korespondencja seryjna | Generuje serię dokumentów z szablonu i wierszy arkusza | B3 (szablony) + I1; pętla generacji (Automations) |
| I5 | Wizualizacja tabeli | Wstawia wykres z danych tabeli do dokumentu | excelize (chart) lub render SVG po stronie klienta |

### 7.10. Rodzina J — Kolaboracja i przekazywanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| J1 | Przekazanie do i z Library | Zapisuje wersję jako artefakt w Library, wczytuje dokument źródłowy z Library | Powiązanie `Studio ───► Library`; `powiazanie_komponentu` |
| J2 | Odbiór raportu z Research | Wczytuje skompletowany raport z Report Builder do redakcji finalnej | Powiązanie `Research ◄──► Studio` |
| J3 | Praca dwujęzyczna z Translate | Współdzieli operacje kontekstowe przy wersji dwujęzycznej | Powiązanie `Studio ◄──► Translate` |
| J4 | Wstawianie zasobów z Design | Osadza grafiki i elementy wizualne z Assets Panel | Powiązanie `Design ───► Studio` |
| J5 | Odwołanie do wersji | Generuje odwołanie do konkretnej wersji dokumentu w obrębie platformy | `wersja_artefaktu` + warstwa udostępnień (konfigurowalna) |
| J6 | Eksport paczki redakcyjnej | Pakuje dokument, historię, raport zmian i adnotacje do jednego archiwum przekazania | F7 + E9; `archive/zip` |

### 7.11. Rodzina K — Produktywność i sterowanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| K1 | Paleta poleceń | Wywołanie dowolnej operacji modułu z klawiatury (`Ctrl/Cmd + K`) | Rejestr operacji modułu; skrót globalny |
| K2 | Operacje z czatu przez „/" | Wywołanie operacji Tools Panel z poziomu Chat Window skrótem `/` | Menu operacji + Chat Window |
| K3 | Makra i łańcuchy operacji | Sekwencja operacji AI wykonywana jednym poleceniem (np. korekta → streszczenie → eksport) | Definicje łańcuchów w `ustawienie`; Execution Loop Window prowadzi przebieg; powiązanie z Automations |
| K4 | Operacje wsadowe na wielu dokumentach | Ta sama operacja na wielu kartach lub plikach naraz | Kolejka + A9; Execution Loop Window; Automations |
| K5 | Skróty klawiszowe konfigurowalne | Pełna mapa skrótów, edytowalna przez Operatora | Mapa skrótów w `ustawienie` |
| K6 | Tryb skupienia | Ukrywa kolumny boczne, centruje wiersz aktywny | Ustawienie widoku edytora |
| K7 | Dyktowanie i odczyt na głos | Wprowadzanie tekstu głosem i odczyt dokumentu | Powiązanie z Assistant (Voice Console); przetwarzanie mowy przez `kanal_modelu` |

### 7.12. Rodzina L — Załączniki multimedialne i specjalne

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| L1 | Transkrypcja audio i wideo do dokumentu | Zamienia nagranie w tekst z podziałem na mówców i znacznikami czasu | Przetwarzanie mowy przez `kanal_modelu` (Whisper); formaty audio i wideo (MP3, WAV, MP4) |
| L2 | Osadzanie i podpisy obrazów | Wstawia obrazy z podpisem i tekstem alternatywnym generowanym przez AI | `imaging`; `kanal_modelu` dla opisu |
| L3 | Kod QR i kody kreskowe | Generuje i wstawia kod QR lub kreskowy (np. odwołanie do wersji) | `boombuler/barcode` |
| L4 | Wzory matematyczne | Wstawia i renderuje wzory z notacji LaTeX | Render KaTeX/MathJax po stronie klienta; eksport przez Pandoc |
| L5 | Diagramy z tekstu | Renderuje diagramy z opisu tekstowego osadzone w dokumencie | Mermaid po stronie klienta; eksport do SVG/PNG |

### 7.13. Podstawa technologiczna modułu

| Obszar | Komponent podstawowy | Uzasadnienie | Rezerwa |
|---|---|---|---|
| PDF (operacje) | `pdfcpu/pdfcpu` | Czyste Go, najszerszy zakres operacji, licencja Apache-2.0, bez CGo | `unidoc/unipdf` przy głębszej manipulacji treści |
| PDF (render i ekstrakcja wierna) | `gen2brain/go-fitz` (MuPDF) | Najwyższa jakość renderu i tekstu; podstawa rozpoznawania skanów i różnicy wizualnej | `ledongthuc/pdf` (lekki, czyste Go, tylko tekst) |
| Rozpoznawanie tekstu | `otiai10/gosseract` + Tesseract | Standard offline, wielojęzyczny, bez kosztu za stronę; klucz usługi chmurowej jawny w oknie konfiguracji | Google Vision / Azure Read przy piśmie odręcznym |
| DOCX / XLSX / PPTX | `qax-os/excelize` (XLSX) + `unidoc/unioffice` (DOCX/PPTX) | excelize darmowy i kompletny; unioffice najpełniejszy dla formatów Office | Pandoc lub LibreOffice jako konwerter procesowy |
| Diff | `sergi/go-diff` | Dojrzały port diff-match-patch, semantyczne czyszczenie różnic | `hexops/gotextdiff` (Myers, unified) |
| Markdown | `yuin/goldmark` | Zgodny z CommonMark, rozszerzalny (GFM, footnotes, tables) | `gomarkdown/markdown` |
| Konwersja uniwersalna | `pandoc` jako proces | Najszerszy zakres formatów | LibreOffice headless (wierność Office) |
| Indeks semantyczny | `blevesearch/bleve` + embeddingi z `kanal_modelu` | Wyszukiwanie znaczeniowe i pytania nad dokumentem | Kontekst pełnotekstowy modelu |

---

## 8. Punkty sterowania z okna konfiguracji

Wszystkie pozycje działają wg zasady „brak ustawienia = wartość domyślna"; żadna nie blokuje pracy. Warstwowość ustawień: globalna → środowisko → projekt → sesja.

| Grupa ustawień | Co Operator personalizuje | Wartość domyślna | Warstwa |
|---|---|---|---|
| Silnik rozpoznawania tekstu | Wybór silnika (lokalny Tesseract / usługa chmurowa), języki (`tessdata`), próg pewności, czyszczenie obrazu przed rozpoznaniem | Tesseract lokalny, języki pol+eng, czyszczenie automatyczne | globalna/sesja |
| Kanały modeli operacji AI | Który `kanal_modelu` obsługuje którą rodzinę operacji (korekta tańszym modelem, redakcja mocniejszym) | model domyślny karty | projekt/sesja |
| Szablony operacji (prompty) | Treść promptów operacji C1–C14, dodawanie operacji własnych C15 | prompty fabryczne | globalna/projekt |
| Łańcuchy i makra | Definicje łańcuchów operacji (K3), operacje wsadowe (K4) | brak definicji — tworzy je użytkownik | projekt |
| Formaty i konwersje | Domyślny format eksportu, ścieżka renderu PDF (fpdf/gopdf albo chromedp/Pandoc), profile eksportu (D9) | PDF przez render docelowy | globalna/sesja |
| Autozapis i wersje | Interwał autozapisu (F8), migawka po operacji AI, polityka retencji historii | autozapis 30 s, migawka włączona | sesja |
| Diff | Tryb domyślny (scalony/dwukolumnowy), granularność (znak/słowo/zdanie), kolory oznaczeń | scalony, słowo | sesja |
| Redakcja i poufność | Reguły wykrywania danych wrażliwych (H5), automatyczne usuwanie metadanych przy eksporcie (H4) | wykrywanie włączone, automatyczne czyszczenie wyłączone | projekt |
| Skróty klawiszowe | Pełna mapa skrótów (K5), paleta poleceń (K1) | mapa domyślna | globalna |
| Widok edytora | Domyślny tryb pracy (WYSIWYG/Markdown), tryb skupienia (K6), motyw typografii | WYSIWYG | sesja |
| Pętla wykonawcza | Widoczność Execution Loop Window przy starcie karty, próg kontroli jakości, limit ponowień zadania | okno otwierane przy zleceniu wielozadaniowym, limit ponowień 3 | projekt/sesja |
| Powiązania międzymodułowe | Ustanowienie `Studio ◄──► Translate`, `Studio ───► Library`, `Studio ◄──► Research`, `Design ───► Studio` | brak powiązań | projekt |
| Izolacja sesji | Współdzielenie historii Chat Window i pamięci między kartami, 8 zakresów izolacji technicznej procesu sesji | pełna izolacja per karta, żaden zakres techniczny nieaktywny | okno izolacji |
| Klucze integracji | Jawne klucze do usług zewnętrznych (rozpoznawanie tekstu w chmurze, model przez API, licencja `unioffice`) — widoczne i edytowalne | puste | globalna |

Zasada kluczy jawnych: klucze do usług zewnętrznych są wprost widoczne i edytowalne w oknie konfiguracji — nigdy ukryte ani zaszyte w kodzie.

Moduł pozostaje zgodny z zasadą nadrzędną platformy: **zero blokad** — żadna funkcja nie wymusza konfiguracji ani nie blokuje pracy, nie występują nieaktywne przyciski akcji ani modalne bramki; **klucze jawne** — klucze usług i licencji widoczne w oknie konfiguracji; **pełna konfigurowalność** — każde zachowanie sterowalne z okna konfiguracji, warstwowo. Domyślnym zachowaniem każdej funkcji jest wykonanie; izolacja, poufność i powiązania są odwracalnymi, jawnymi decyzjami Operatora, nie wymogami.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + B`, `I`, `U` | Pogrubienie, kursywa, podkreślenie | Studio Editor |
| `Ctrl/Cmd + F` | Otwarcie Znajdź/Zamień | Studio Editor |
| `Ctrl/Cmd + S` | Zapis ręczny (uzupełnia autozapis) | Studio Editor |
| `Ctrl/Cmd + K` | Paleta poleceń modułu | Wszystkie okna modułu |
| `Ctrl/Cmd + Enter` | Wysłanie polecenia do Wykonawcy | Chat Window |
| `↑` / `↓` w pustym polu poleceń | Nawigacja po historii poleceń | Chat Window |
| `/` | Wywołanie menu operacji Tools Panel z poziomu czatu | Chat Window |
| `Esc` | Wstrzymanie pętli wykonawczej | Execution Loop Window |
| Ikona `olowek` | Edycja | Studio Editor, Session Repository |
| Ikona `strzalka-lewo` / `strzalka-prawo` | Nawigacja między różnicami | Diff/Grep Panel |
| Ikona `zegar` | Znacznik czasu wersji | Session Repository |
| Ikona `gwiazdka` | Etykieta wersji | Session Repository |
| Ikona `pobierz` | Eksport | Preview Window, Session Repository |
| Ikona `oko` | Podgląd | Session Repository, Diff/Grep Panel, Ingest/OCR Panel |
| Ikona `szukaj` | Wyszukiwanie / Grep | Diff/Grep Panel, Studio Editor |
| Ikona `petla` | Stan zadania pętli wykonawczej | Execution Loop Window |
| Ikona `skan` | Pozycja kolejki cyfryzacji | Ingest/OCR Panel |

*Koniec dokumentu. Moduł Studio — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
