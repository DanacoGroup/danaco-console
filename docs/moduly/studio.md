# Danaco Console — Moduł Studio

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
| **Tytuł** | Moduł Studio |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant · Operator |
| **Przeznaczenie** | Ustala zakres funkcjonalny, komplet okien operacyjnych i zachowanie modułu Studio — zaawansowanego stanowiska pracy z tekstem, dokumentami i treścią — jako źródło wykonawcze dla dewelopera i projektanta. |
| **Zakres** | komplet okien operacyjnych modułu Studio (Chat Window, Execution Loop Window, Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Ingest/OCR Panel, Preview Window), katalog funkcji, komendy obszaru `studio`, żetony i komponenty widoku, stany kontrolek, przebiegi pracy i scenariusze użycia |
| **Poza zakresem** | praca nad kodem źródłowym i repozytoria Git — Developer/Terminal; Studio niedostępny w środowisku CodeStudio (rozdz. 1.5) |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Moduł Translate](translate.md) · [Moduł Library](library.md) · [Moduł Research](research.md) · [Model danych](../architektura/model-danych.md) · [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) · [Mobile](../funkcje-globalne/mobile.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/studio.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `studio`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` · `design/05-okna/moduly/studio.html` |
| **Zasada nadrzędna** | Zero blokad, klucze jawne, pełna kompozycyjność i pełna konfigurowalność z okna konfiguracji; domyślne zachowanie modułu = wykonanie. Każda funkcja działa bez wymuszonej konfiguracji, a każde ograniczenie jest jawnym, odwracalnym ustawieniem Operatora |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
   - [1.1 Definicja](#11-definicja)
   - [1.2 Dla kogo](#12-dla-kogo)
   - [1.3 Po co — wartość modułu](#13-po-co--wartość-modułu)
   - [1.4 Granica tematyczna — zakres modułu](#14-granica-tematyczna--zakres-modułu)
   - [1.5 Granica tematyczna — poza zakresem modułu](#15-granica-tematyczna--poza-zakresem-modułu)
   - [1.6 Miejsce w architekturze platformy](#16-miejsce-w-architekturze-platformy)
   - [1.7 Dostępność i forma udostępnienia](#17-dostępność-i-forma-udostępnienia)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - [2.1 Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
   - [3.1 Chat Window (okno wspólne)](#31-chat-window-okno-wspólne)
   - [3.2 Execution Loop Window](#32-execution-loop-window)
   - [3.3 Studio Editor](#33-studio-editor)
   - [3.4 Tools Panel](#34-tools-panel)
   - [3.5 Diff/Grep Panel](#35-diffgrep-panel)
   - [3.6 Session Repository](#36-session-repository)
   - [3.7 Preview Window](#37-preview-window)
   - [3.8 Ingest/OCR Panel](#38-ingestocr-panel)
4. [Przepływy pracy](#4-przepływy-pracy)
   - [4.1 Przepływ podstawowy — redakcja fragmentu z korektą AI](#41-przepływ-podstawowy--redakcja-fragmentu-z-korektą-ai)
   - [4.2 Przepływ cyfryzacji — skan do dokumentu edytowalnego](#42-przepływ-cyfryzacji--skan-do-dokumentu-edytowalnego)
   - [4.3 Przepływ rozszerzony — praca dwujęzyczna z modułem Translate](#43-przepływ-rozszerzony--praca-dwujęzyczna-z-modułem-translate)
   - [4.4 Przepływ rozszerzony — redakcja raportu przekazanego z Research](#44-przepływ-rozszerzony--redakcja-raportu-przekazanego-z-research)
   - [4.5 Przepływ zbiorczy — wiele dokumentów równolegle](#45-przepływ-zbiorczy--wiele-dokumentów-równolegle)
5. [Stany, dane i powiązania](#5-stany-dane-i-powiązania)
   - [5.1 Model stanów sesji modułu](#51-model-stanów-sesji-modułu)
   - [5.2 Model danych wykorzystywany przez moduł](#52-model-danych-wykorzystywany-przez-moduł)
   - [5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Studio](#53-izolacja-i-konfigurowalność--punkty-właściwe-modułowi-studio)
   - [5.4 Powiązania z innymi modułami](#54-powiązania-z-innymi-modułami)
6. [Scenariusze użycia](#6-scenariusze-użycia)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
   - [7.1 Rodzina A — Wczytywanie, formaty i cyfryzacja](#71-rodzina-a--wczytywanie-formaty-i-cyfryzacja)
   - [7.2 Rodzina B — Edycja i tworzenie treści](#72-rodzina-b--edycja-i-tworzenie-treści)
   - [7.3 Rodzina C — Operacje kontekstowe AI na treści](#73-rodzina-c--operacje-kontekstowe-ai-na-treści)
   - [7.4 Rodzina D — Konwersje i eksport](#74-rodzina-d--konwersje-i-eksport)
   - [7.5 Rodzina E — Diff, porównania i wyszukiwanie](#75-rodzina-e--diff-porównania-i-wyszukiwanie)
   - [7.6 Rodzina F — Wersjonowanie i historia](#76-rodzina-f--wersjonowanie-i-historia)
   - [7.7 Rodzina G — Operacje na PDF](#77-rodzina-g--operacje-na-pdf)
   - [7.8 Rodzina H — Bezpieczeństwo i poufność dokumentu](#78-rodzina-h--bezpieczeństwo-i-poufność-dokumentu)
   - [7.9 Rodzina I — Arkusze i dane tabelaryczne](#79-rodzina-i--arkusze-i-dane-tabelaryczne)
   - [7.10 Rodzina J — Kolaboracja i przekazywanie](#710-rodzina-j--kolaboracja-i-przekazywanie)
   - [7.11 Rodzina K — Produktywność i sterowanie](#711-rodzina-k--produktywność-i-sterowanie)
   - [7.12 Rodzina L — Załączniki multimedialne i specjalne](#712-rodzina-l--załączniki-multimedialne-i-specjalne)
   - [7.13 Podstawa technologiczna modułu](#713-podstawa-technologiczna-modułu)
   - [7.14 Przebiegi wybranych funkcji i przypadki brzegowe](#714-przebiegi-wybranych-funkcji-i-przypadki-brzegowe)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)
10. [Punkty łamania i kryteria odbioru](#10-punkty-łamania-i-kryteria-odbioru)
   - [10.1 Punkty łamania](#101-punkty-łamania)
   - [10.2 Kryteria odbioru](#102-kryteria-odbioru)
11. [Załącznik — pełny wykaz komend kontraktu modułu Studio](#11-załącznik--pełny-wykaz-komend-kontraktu-modułu-studio)
   - [11.1 Obszar `studio` — 179 komend](#111-obszar-studio--179-komend)

---

## 1. Przeznaczenie i kontekst

### 1.1 Definicja

Studio jest modułem zaawansowanej pracy z tekstem, dokumentami i treścią. Integruje w jednym oknie roboczym trzy czynności dotąd rozproszone między osobne narzędzia: edycję treści, wsparcie AI działające bezpośrednio na zaznaczonym fragmencie oraz porównywanie wersji dokumentu w czasie. Moduł nie jest edytorem tekstu z dopiętym czatem — jest jedną przestrzenią, w której edytor, operacje kontekstowe AI i historia wersji współdzielą ten sam dokument w czasie rzeczywistym.

Studio stanowi kompletne stanowisko pracy z dokumentami: obejmuje wczytanie dowolnego obsługiwanego formatu, cyfryzację skanów przez OCR, konwersje między formatami, redakcję wspieraną AI, porównania różnicowe, wersjonowanie, eksport oraz zabezpieczenie i podpis dokumentu. Jeden moduł pokrywa cały obszar pracy dokumentowej, bez sięgania po osobne aplikacje desktopowe i usługi sieciowe.

### 1.2 Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Studio |
|---|---|
| Redaktorzy i copywriterzy | Przepisywanie, zmiana stylu, dopracowanie tonu dużych partii tekstu |
| Analitycy i konsultanci | Przekształcanie notatek roboczych w dokumenty finalne, streszczenia długich materiałów |
| Prawnicy i specjaliści dokumentacji | Praca na wersjach dokumentu z pełną historią zmian i porównaniem redakcji |
| Zespoły wielojęzyczne | Praca dwujęzyczna przy współdzieleniu operacji kontekstowych AI z modułem Translate |
| Twórcy raportów końcowych | Redakcja materiału przekazanego z modułu Research |

### 1.3 Po co — wartość modułu

| Problem klasycznego rozproszenia narzędzi | Rozwiązanie w Studio |
|---|---|
| Edytor tekstu i okno czatu AI to dwie osobne aplikacje | Studio Editor i Chat Window w jednej karcie sesji, nad tym samym dokumentem |
| Brak podglądu, co dokładnie zmieniło AI | Diff/Grep Panel pokazuje różnicę przed/po dla każdej operacji |
| Utrata wcześniejszych redakcji przy nadpisaniu pliku | Session Repository gromadzi każdą wersję ze znacznikiem czasu |
| Ręczne przełączanie się między formatem roboczym a wynikowym | Preview Window renderuje na żywo format docelowy |
| Skan lub zdjęcie dokumentu wymaga osobnego programu OCR | Ingest/OCR Panel prowadzi cyfryzację wewnątrz sesji modułu |

### 1.4 Granica tematyczna — zakres modułu

| Obszar | Zakres pokrycia w Studio |
|---|---|
| Formaty tekstowe i dokumentowe | PDF, DOCX/DOC/ODT/RTF, TXT, Markdown/MDX, HTML, LaTeX, EPUB, arkusze (XLSX/ODS/CSV), prezentacje (odczyt), e-mail (EML/MSG odczyt), obrazy dokumentów (PNG/JPG/TIFF) jako wejście OCR |
| Wczytywanie i cyfryzacja | Import lokalny, z Library, ze schowka, z URL, ze skanera i kamery; OCR skanów i zdjęć, rozpoznanie układu, prostowanie skosu i czyszczenie obrazu |
| Redakcja i tworzenie | Edytor WYSIWYG i tekstowy, style typografii, tabele, obrazy, przypisy, spisy treści, szablony dokumentów, składanie z bloków |
| Operacje kontekstowe AI | Korekta, styl, ton, streszczenia, ekstrakcja, przepisanie strukturalne, tłumaczenie, wzbogacanie, weryfikacja, pytania i odpowiedzi nad dokumentem (RAG dokumentowy) |
| Porównania i wyszukiwanie | Diff tekstowy i wizualny, Grep/regex w dokumencie i w historii, wyszukiwanie semantyczne, porównanie treści z materiałem źródłowym |
| Wersjonowanie | Historia sesyjna, etykiety, gałęzie, przywracanie, eksport historii, znaczniki autora (`uzytkownik`/`model`) |
| Konwersje i eksport | Każdy-do-każdego w granicach obsługiwanych formatów, druk, podgląd i konfiguracja wydruku, znak wodny, redakcja poufności, kompresja, scalanie i dzielenie PDF |
| Bezpieczeństwo dokumentu | Szyfrowanie PDF, podpis cyfrowy, usuwanie metadanych, redakcja poufności (trwałe zamazanie), kontrola dostępu na poziomie artefaktu — każde jako jawne, odwracalne ustawienie |

### 1.5 Granica tematyczna — poza zakresem modułu

| Poza zakresem | Uzasadnienie | Właściwy moduł |
|---|---|---|
| Praca nad kodem źródłowym, repozytoria Git | Studio niedostępny w CodeStudio; edycja kodu to inna dyscyplina | Developer, Terminal |
| Generowanie i edycja grafiki koncepcyjnej, ilustracji, brandingu | Studio korzysta z gotowych zasobów wizualnych, ale ich nie tworzy | Design |
| Trwałe repozytorium wiedzy całej organizacji, katalogowanie, kolekcje | Studio wersjonuje w obrębie sesji; magazyn trwały należy do innego modułu | Library |
| Wieloźródłowe badania, benchmarki, budowa raportu z surowych źródeł | Studio redaguje gotowy materiał, nie prowadzi badania | Research |
| Zaawansowane, wielozadaniowe tłumaczenie z zarządzaniem terminologią i pamięcią tłumaczeń | Studio tłumaczy fragmenty kontekstowo; pełna lokalizacja należy do innego modułu | Translate |
| Arkusze jako narzędzie analityczne (pełne modele, tabele przestawne, wykresy analityczne) | Studio odczytuje i edytuje arkusze jako dokumenty, nie jest kalkulatorem analitycznym | Workspace/Research + analiza danych |

Granica jest miękka i konfigurowalna: powiązania `Studio ◄──► Translate`, `Studio ───► Library`, `Studio ◄──► Research`, `Design ───► Studio` pozwalają wywołać sąsiedni moduł bez opuszczania sesji, gdy Operator je ustanowi.

### 1.6 Miejsce w architekturze platformy

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

### 1.7 Dostępność i forma udostępnienia

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

### 2.1 Warstwy widoczności w module

Moduł Studio realizuje regułę stopniowego ujawniania funkcjonalności: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Pełny arsenał modułu — dwanaście rodzin funkcji, warsztat PDF, cyfryzacja OCR, wersjonowanie z gałęziami i zabezpieczenie dokumentu — istnieje w architekturze modułu i pozostaje niewidoczny w interfejsie do chwili wystąpienia potrzeby użycia. Liczba obsługiwanych formatów, operacji kontekstowych i paneli nie wpływa na postrzeganą prostotę okna roboczego.

**Warstwa 1 — zawsze widoczna.** Chat Window w lewej kolumnie wraz z polem poleceń i strumieniem odpowiedzi, Studio Editor jako okno wiodące z obszarem treści dokumentu, pasek statusu dokumentu (status zapisu, liczba słów, numer wersji), pasek kontekstu ze znacznikami modelu, wykonawcy i formatu docelowego oraz wskaźniki stanu wykonania (wskaźnik pracy Wykonawcy, wskaźnik przebiegu pętli). Warstwa ta zajmuje ponad 80% powierzchni modułu.

**Warstwa 2 — widoczna na żądanie.** Selektor kanału i modelu bieżącej karty, przełącznik trybu pracy edytora (WYSIWYG / Markdown / źródło / nakładka OCR), przełącznik zakresu operacji (zaznaczenie / cały dokument / wszystkie karty), selektor formatu podglądu i profilu eksportu, Execution Loop Window oraz Preview Window. Elementy te wywołuje pojedyncze kliknięcie ikony, przełącznika lub znacznika kontekstowego w pasku kontekstu (`[Fable 5] [Ultra] [Format: PDF] [⇄ pętla]`); po użyciu element zwija się samoczynnie.

**Warstwa 3 — rozwinięcia kontekstowe.** Tools Panel z pełnym katalogiem operacji kontekstowych AI, Diff/Grep Panel, Session Repository i Ingest/OCR Panel, pasek pływający zaznaczenia, menu „Warsztat PDF”, menu „Zabezpiecz”, menu wersji `⋯`, panel Znajdź/Zamień, filtry różnic i filtry historii. Wywołanie następuje przez elementy zbiorcze `Operacje ▼` i `Narzędzia ▼`, menu kebab (⋮), menu hamburger (☰), menu kontekstowe zaznaczenia, panel popover lub listę rozwijaną. Po zamknięciu panel znika całkowicie z przestrzeni roboczej.

**Warstwa 4 — funkcje eksperckie.** Numeracja Bates i formularze AcroForm, szyfrowanie i uprawnienia PDF, podpis cyfrowy i jego weryfikacja, redakcja poufności i usuwanie metadanych, scalanie gałęzi dokumentu z rozwiązywaniem konfliktów per fragment, korekta rozpoznania OCR na warstwie tekstowej, wybór silnika OCR i progu pewności rozpoznania, operacje wsadowe i łańcuchy operacji na wielu kartach, definiowanie własnych promptów jako nazwanych narzędzi, eksport paczki redakcyjnej oraz diagnostyka przebiegu pętli wykonawczej. Dostęp prowadzi przez polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarkę funkcji, tryb administracyjny albo konfigurację roli. Użytkownik podstawowy nie widzi tych elementów.

**Zasada jednego kliknięcia.** Każda funkcja modułu ukryta w warstwach 2–4 jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Zagnieżdżanie funkcji Studio głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny okna roboczego, nie utrudnia dostępu do funkcji.

---

## 3. Specyfikacja okien operacyjnych

### 3.1 Chat Window (okno wspólne)

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
- Odwołania do zaznaczonego fragmentu dokumentu jako przedmiotu polecenia („zastosuj do zaznaczenia”).
- Wywoływanie operacji z Tools Panel bezpośrednio z poziomu wiadomości (skrót poleceniowy `/`).
- Pytania i odpowiedzi nad dokumentem z cytowaniem miejsca w treści (RAG dokumentowy, funkcja „Pytania i odpowiedzi nad dokumentem (RAG)”).
- Dyktowanie polecenia głosem i odczyt odpowiedzi na głos (funkcja „Dyktowanie i odczyt na głos”).
- Cytowanie fragmentów Diff/Grep w treści odpowiedzi (blok czcionką `--dn-ff-mono`).
- Załączanie plików do wiadomości (obsługiwane formaty jak w Studio Editor).
- Historia poleceń (strzałka górna/dolna do poprzednich promptów).
- Regeneracja odpowiedzi, edycja własnego polecenia i ponowne wysłanie, rozgałęzianie wątku odpowiedzi.
- Kopiowanie odpowiedzi, wstawienie odpowiedzi bezpośrednio do Studio Editor w miejscu kursora.
- Zatwierdzanie i przerywanie działań zleconych Wykonawcy.
- Wskaźnik pracy Wykonawcy i wskaźnik długości kontekstu bieżącej rozmowy.

**Makieta tekstowa.**

```
┌─ Chat Window ──────────────────────┐
│ Użytkownik ↔ Wykonawca           ⋮ │
├────────────────────────────────────┤
│  Historia rozmowy (przewijana)     │
│   ┌─ Użytkownik ────────────────┐  │
│   │ „popraw styl zaznaczonego   │  │
│   │  akapitu”                   │  │
│   └─────────────────────────────┘  │
│   ┌─ Wykonawca (strumień) ──────┐  │
│   │ tekst odpowiedzi budowany   │  │
│   │ token po tokenie …          │  │
│   └─────────────────────────────┘  │
├────────────────────────────────────┤
│ [ 📎 ] [ / ]                       │
│ Pole poleceń ………… [ Wyślij ▶ ]     │
│ [Fable 5] [Ultra] [⇄ pętla]        │
└────────────────────────────────────┘
   lewa kolumna, stała, pełna
   wysokość obszaru roboczego
   stan spoczynku: warstwa 1 oraz
   znaczniki kontekstowe i menu ⋮
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor kanału modelu | Rozwijana lista w nagłówku | Wybór/zmiana modelu bazowego bieżącej karty | Mały przycisk z etykietą + strzałka | spoczynek · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty · ładowanie przełączania modelu (`.dn-spinner`, `[aria-busy='true']`) | Otwiera listę kanałów modelu; wybór zmienia model dla kolejnych wiadomości | Nagłówek Chat Window | 2 | Kliknięcie znacznika modelu w pasku kontekstu; lista zwija się samoczynnie po wyborze |
| Pole poleceń | Wieloliniowe pole tekstowe `.dn-prompt` | Wprowadzanie polecenia do Wykonawcy | Duże pole na końcu kolumny, rozciągliwe w pionie | spoczynek · wskazanie kursorem (`:hover`) · ognisko (`.dn-prompt:focus`, żeton `--dn-fokus`) · z załącznikiem · błąd pustej wysyłki (`.dn-pole-blad`, `--dn-blad-obrys`) | Enter wysyła, Shift+Enter nowa linia | Zamknięcie kolumny Chat Window | 1 | Widoczne bez interakcji |
| Przycisk „Wyślij” | Przycisk ikonowy `.dn-btn-ikona` (wyslij) | Wysłanie polecenia | Mała ikona 36×36 px, akcent sygnałowy gdy pole niepuste | spoczynek · wskazanie kursorem (`.dn-btn-ikona:hover`) · wciśnięcie (`.dn-btn-ikona:active`) · ognisko (`.dn-btn-ikona:focus`, żeton `--dn-fokus`) · nieaktywny (pole puste) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Wysyła polecenie, czyści pole, uruchamia strumień odpowiedzi; przy pustym polu kliknięcie sygnalizuje brak treści komunikatem, bez wysyłki | Obok pola poleceń | 1 | Widoczny bez interakcji |
| Przycisk załącznika | Przycisk ikonowy `.dn-btn-ikona` (spinacz) | Dołączenie pliku do wiadomości | Mała ikona | spoczynek · wskazanie kursorem (`.dn-btn-ikona:hover`) · ognisko (`.dn-btn-ikona:focus`, żeton `--dn-fokus`) · aktywny — plik dołączony (`[aria-pressed='true']`) | Otwiera okno wyboru pliku lub akceptuje przeciągnięcie | Lewa strona pola poleceń | 1 | Widoczny bez interakcji |
| Dymek wiadomości użytkownika | Karta wiadomości wyrównana do prawej (`.dn-wpis--czlowiek`) | Prezentacja wysłanego polecenia | Średni blok tekstu z awatarem `.dn-awatar` | spoczynek · wskazanie kursorem (`:hover` ujawnia ikonę edycji) | Edytowalny — ikona ołówek przy wskazaniu kursorem | Historia rozmowy | 1 | Widoczny bez interakcji |
| Dymek odpowiedzi Wykonawcy | Karta wiadomości wyrównana do lewej (`.dn-wpis--inteligencja`) | Prezentacja odpowiedzi Wykonawcy | Średni blok tekstu, treść może zawierać fragment czcionką `--dn-ff-mono` | strumieniowanie — kursor migający (`.dn-spinner`, `[aria-busy='true']`) · kompletny · wskazanie kursorem (`:hover` ujawnia pasek akcji) · błąd generowania (`--dn-blad-tekst`) | Po wskazaniu kursorem ujawnia pasek akcji (Wstaw / Kopiuj / Regeneruj / Przerwij) | Historia rozmowy | 1 | Widoczny bez interakcji |
| Przycisk „Wstaw do dokumentu” | Przycisk `.dn-btn--zarys`, mały | Przeniesienie treści odpowiedzi do Studio Editor | Mały przycisk tekstowy w pasku akcji dymka | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Wstawia treść w miejscu kursora edytora lub w miejscu zaznaczenia | Pasek akcji dymka odpowiedzi | 3 | Pasek akcji dymka ujawniany po najechaniu na odpowiedź |
| Przycisk „Przerwij” | Przycisk `.dn-btn--duch`, mały | Zatrzymanie trwającego działania Wykonawcy | Mały przycisk tekstowy | ukryty (brak działania) · aktywny · wskazanie kursorem (`.dn-btn--duch:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Przerywa strumień i zadanie, pozostawiając dokument w stanie sprzed operacji | Pasek akcji dymka odpowiedzi | 3 | Pasek akcji dymka ujawniany po najechaniu na odpowiedź w trakcie strumienia |
| Wskaźnik strumienia | Migający kursor / `.dn-spinner` | Sygnalizacja generowania odpowiedzi na żywo | Mały, subtelny | ukryty · aktywny (`.dn-spinner`, `[aria-busy='true']`) | Znika po zakończeniu strumienia | Koniec treści generowanej odpowiedzi | 1 | Widoczny bez interakcji w trakcie generowania |
| Skrót „/ operacje” | Menu podpowiedzi poleceń | Szybkie wywołanie operacji Tools Panel z poziomu czatu | Rozwijana lista przy polu poleceń | ukryty · rozwinięty (po wpisaniu „/”) | Wstawia szablon polecenia operacji do pola | Przy polu poleceń, wyzwalane znakiem „/” | 3 | Wpisanie znaku „/” w polu poleceń |

---

### 3.2 Execution Loop Window

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — kanał Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Pętla wykonawcza odrębna per karta sesji; przebieg powiązany z dokumentem bieżącej karty |

Execution Loop Window prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą realizującym zadania. Okno odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów modułu Studio. Zlecenie dokumentowe wydane w Chat Window (na przykład „przygotuj wersję finalną raportu”: cyfryzacja skanu, korekta, streszczenie zarządcze, eksport PDF) Koordynator rozkłada tu na zadania jednostkowe, przydziela je Wykonawcy, kontroluje jakość wyniku każdego zadania i decyduje o ponowieniu.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie dokumentowe i jego dekompozycja na zadania (wczytanie i OCR, operacja kontekstowa AI, wygenerowanie różnicy, zapis wersji, konwersja, eksport).
- Kolejka zadań i stan każdego zadania (oczekuje, w realizacji, kontrola jakości, ponowienie, zakończone, przerwane).
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, z pełnym zapisem przebiegu.
- Wyniki kontroli jakości zadania — zgodność wyniku operacji AI z zakresem zaznaczenia, poprawność warstwy tekstowej po OCR, kompletność renderu eksportu — oraz decyzje o ponowieniu zadania.
- Wskaźniki przebiegu pętli: liczba zadań w kolejce, liczba ponowień, czas realizacji, zużycie kontekstu modelu.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia w toku.
- Powiązanie zadania z artefaktem i wersją, do których się odnosi — kliknięcie zadania otwiera odpowiadającą mu pozycję w Session Repository lub blok różnicy w Diff/Grep Panel.
- Przebieg łańcuchów operacji i operacji wsadowych na wielu dokumentach (funkcje „Makra i łańcuchy operacji” i „Operacje wsadowe na wielu dokumentach”) jako zadania jednej pętli.

**Makieta tekstowa.**

```
┌─ Execution Loop Window ─────────⋮─┐
│ Zlecenie: „wersja finalna raportu”│
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
| Nagłówek zlecenia | Blok tytułowy z treścią zlecenia | Wskazanie zlecenia prowadzonego przez pętlę | Średni blok tekstu, wyróżnienie sygnałowe | aktywne · wstrzymane · zakończone · przerwane | Kliknięcie rozwija pełną treść zlecenia i jego kontekst | Góra kolumny | 2 | Widoczny po otwarciu Execution Loop Window znacznikiem [⇄ pętla] |
| Pozycja zadania | Wiersz kolejki `.dn-karta--klikalna` | Reprezentuje jedno zadanie dekompozycji | Średni wiersz z ikoną stanu i etykietą | oczekuje · w realizacji (`.dn-spinner`) · kontrola jakości · ponowienie · zakończone · przerwane | Kliknięcie otwiera powiązaną wersję w Session Repository lub różnicę w Diff/Grep Panel | Kolejka zadań | 2 | Widoczna po otwarciu Execution Loop Window |
| Strumień komunikatów sterujących | Lista wiadomości Koordynator ↔ Wykonawca | Podgląd wymiany sterującej | Blok tekstu czcionką `--dn-ff-mono`, przewijany | strumieniowanie · kompletny | Przewijanie; kopiowanie pojedynczego komunikatu | Środek kolumny | 3 | Rozwinięcie sekcji komunikatów w Execution Loop Window |
| Znacznik kontroli jakości | Plakietka `.dn-plakietka` przy zadaniu | Wynik kontroli jakości zadania | Mała pigułka z ikoną | zaliczone · odrzucone (ponowienie) | Kliknięcie ujawnia kryterium i uzasadnienie decyzji | Prawa krawędź pozycji zadania | 2 | Widoczny przy pozycji zadania; kliknięcie ujawnia kryterium i uzasadnienie (warstwa 3) |
| Wskaźniki przebiegu pętli | Wiersz liczbowy | Liczba zadań, ponowień, czas, zużycie kontekstu | Mały pasek tekstowy `--dn-fs-xs` | aktualizowany na żywo | — (informacyjny) | Zamknięcie kolumny, nad sterowaniem | 1 | Widoczne bez interakcji jako znacznik stanu wykonania w pasku kontekstu Chat Window |
| Przyciski sterowania przebiegiem | Zestaw trzech przycisków | Wstrzymanie, wznowienie, przerwanie pętli | Małe przyciski `--zarys`, w rzędzie | domyślny · aktywny · nieosiągalny stan ukryty | Wstrzymanie zatrzymuje kolejkę po bieżącym zadaniu; wznowienie kontynuuje; przerwanie kończy pętlę i pozostawia wersje już zapisane | Zamknięcie kolumny | 2 | Widoczne po otwarciu Execution Loop Window |
| Przycisk „Koryguj zlecenie” | Przycisk `.dn-btn--sygnal`, pełna szerokość kolumny | Zmiana treści zlecenia w toku pętli | Średni przycisk CTA | domyślny · ładowanie | Otwiera pole edycji zlecenia; zatwierdzenie powoduje ponowną dekompozycję zadań niezrealizowanych | Zamknięcie kolumny | 3 | Menu kebab (⋮) Execution Loop Window |

---

### 3.3 Studio Editor

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
- Pasek narzędzi formatowania: pogrubienie, kursywa, podkreślenie, przekreślenie, nagłówki poziomów od pierwszego do szóstego, cytat blokowy, lista punktowana i numerowana, lista zadań, tabela, blok kodu, link, obraz, separator poziomy.
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
- Komentarze przypięte do fragmentu, wątki komentarzy, oznaczenie „rozwiązane”.
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
│                                                           │
│                                                           │
│                                                           │
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
| Pasek narzędzi formatowania | Rząd przycisków ikonowych `.dn-btn-ikona` | Formatowanie zaznaczonego tekstu | Pasek poziomy, wysokość jednego wiersza, stały na początku kolumny | spoczynek · wskazanie kursorem (`.dn-btn-ikona:hover`) · wciśnięcie (`.dn-btn-ikona:active`) · ognisko (`.dn-btn-ikona:focus`, żeton `--dn-fokus`) · aktywny — styl zastosowany do zaznaczenia (`[aria-pressed='true']`) · nieaktywny (brak otwartego dokumentu) | Kliknięcie zmienia formatowanie zaznaczenia lub pozycji kursora; bez otwartego dokumentu kliknięcie sygnalizuje brak treści do sformatowania | Otwarcie kolumny Studio Editor | 3 | Element zbiorczy `Operacje ▼` w pasku Studio Editor |
| Przełącznik trybu pracy | Rozwijana lista trybów | Wybór WYSIWYG / Markdown / źródło / nakładka OCR | Mały przycisk z etykietą | WYSIWYG (domyślny) · Markdown · źródło · nakładka OCR | Przerysowuje obszar treści w wybranym trybie, zachowując pozycję kursora | Pasek narzędzi | 2 | Kliknięcie znacznika trybu w pasku kontekstu; zwija się po wyborze |
| Obszar treści | Edytowalny obszar dokumentu | Główna przestrzeń pisania i odczytu | Duży, dominujący blok, przewijany w pionie | pusty stan (`.dn-pusty-stan`, brak dokumentu) · edycja · tylko odczyt (podgląd) | Wpisywanie tekstu, zaznaczanie, wklejanie | Środek kolumny | 1 | Widoczny bez interakcji |
| Pasek pływający zaznaczenia | Mały kontekstowy pasek narzędzi | Szybki dostęp do operacji AI na zaznaczeniu | Mały, unoszący się nad zaznaczonym tekstem | ukryty · widoczny (po zaznaczeniu ≥ 1 znaku) | Kliknięcie operacji otwiera Tools Panel z wypełnionym kontekstem lub wykonuje operację wprost | Nad zaznaczonym fragmentem | 3 | Zaznaczenie fragmentu treści (menu kontekstowe zaznaczenia) |
| Znajdź/Zamień | Ikona lupy + panel rozwijany | Wyszukiwanie i zamiana wzorca w treści | Mała ikona w pasku narzędzi, rozwija pole `.dn-pole` | zwinięty · rozwinięty · ognisko pola (`.dn-pole-kontrolka:focus`, żeton `--dn-fokus`) · trafienia podświetlone · brak trafień (`.dn-pole-blad`, `--dn-blad-obrys`) | Enter przechodzi do kolejnego trafienia; „Zamień wszystko” prezentuje podgląd zmian przed zatwierdzeniem | Prawy róg paska narzędzi | 3 | Pozycja `Znajdź/Zamień` w menu `Narzędzia ▼` lub skrót klawiszowy |
| Pasek statusu dokumentu | Wiersz informacyjny | Liczba słów, status zapisu, numer wersji | Mały, niski pasek, tekst pomocniczy | zapisano · zapisywanie (`.dn-spinner`) · niezapisane zmiany | Kliknięcie numeru wersji otwiera Session Repository na tej pozycji | Zamknięcie kolumny Studio Editor | 1 | Widoczny bez interakcji |
| Selektor stylu akapitu | Rozwijana lista | Zastosowanie predefiniowanego stylu typografii | Mały przycisk z etykietą | spoczynek · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty | Zmienia styl akapitu, w którym znajduje się kursor | Pasek narzędzi | 3 | Element zbiorczy `Operacje ▼`, grupa stylów typografii |
| Minimapa nagłówków | Wąska kolumna zwijana | Nawigacja po nagłówkach dokumentu | Wąska kolumna, chowana i wysuwana | zwinięta · rozwinięta | Kliknięcie pozycji przewija dokument do nagłówka | Lewa krawędź obszaru treści | 2 | Kliknięcie uchwytu przy lewej krawędzi obszaru treści; zwija się po wyborze nagłówka |
| Znacznik śledzenia zmian | Oznaczenie wstawienia lub usunięcia w treści | Prezentacja zmiany z autorem | Podkreślenie lub przekreślenie w kolorze autora | wstawienie · usunięcie · zaakceptowana · odrzucona | Kliknięcie ujawnia autora i przyciski akceptacji lub odrzucenia | Obszar treści | 2 | Widoczny po włączeniu śledzenia zmian znacznikiem kontekstowym trybu redakcji |
| Komentarz redakcyjny | Karta komentarza przypięta do fragmentu | Dyskusja redakcyjna nad fragmentem | Mała karta w kolumnie marginesu | otwarty · rozwiązany | Kliknięcie rozwija wątek i pole odpowiedzi | Prawa krawędź obszaru treści | 3 | Pozycja `Komentarz` w menu kontekstowym zaznaczenia; karta znika po zamknięciu wątku |

---

### 3.4 Tools Panel

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
| Przełącznik zakresu | Zestaw przycisków radiowych | Wybór, czy operacja dotyczy zaznaczenia, całego dokumentu czy wszystkich kart | Mały zestaw `.dn-check` typu radio | spoczynek — zaznaczenie, gdy istnieje · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · wybrany: cały dokument / wszystkie karty (`[aria-pressed='true']`) | Zmienia zakres wszystkich kolejnych operacji | Góra panelu | 2 | Widoczny po otwarciu Tools Panel; zwija się do znacznika zakresu po wyborze |
| Grupa operacji | Nagłówek kategorii z listą pozycji | Porządkuje operacje tematycznie | Nagłówek `.dn-karta-tytul` + lista `.dn-karta--klikalna` | zwinięta · rozwinięta | Kliknięcie nagłówka zwija/rozwija grupę | Ciało panelu | 3 | Rozwinięcie kategorii w Tools Panel |
| Pozycja operacji | Pojedynczy wiersz listy `.dn-karta--klikalna` | Wybór konkretnej operacji AI | Mały wiersz z ikoną i etykietą | spoczynek · wskazanie kursorem (`.dn-karta--klikalna:hover`) · wciśnięcie (`.dn-karta--klikalna:active`) · wybrany (`.dn-karta--wybrana`) | Kliknięcie zaznacza operację i (dla operacji z parametrem) rozwija podmenu | Wewnątrz grupy operacji | 3 | Rozwinięcie grupy operacji |
| Podmenu z listą rozwijaną | Np. „Zmień styl ▾” | Wybór wariantu operacji (np. konkretny styl docelowy) | Mała rozwijana lista przy pozycji | zwinięte · rozwinięte | Wybór wariantu zamyka listę i ustawia parametr operacji | Przy pozycjach parametryzowanych | 3 | Kliknięcie pozycji parametryzowanej — lista rozwijana |
| Pozycja łańcucha operacji | Wiersz listy `.dn-karta--klikalna` z ikoną sekwencji | Uruchomienie sekwencji operacji jednym poleceniem | Mały wiersz z licznikiem kroków | spoczynek · wskazanie kursorem (`.dn-karta--klikalna:hover`) · w realizacji (`.dn-spinner`, `[aria-busy='true']`) | Przekazuje sekwencję do Execution Loop Window jako zadania jednej pętli | Grupa „Łańcuchy” | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji lub grupa „Łańcuchy” dostępna po włączeniu operacji zaawansowanych w konfiguracji roli |
| Przycisk „Uruchom operację” | Główny CTA panelu | Wysłanie wybranej operacji do wykonania | Duży przycisk `.dn-btn--sygnal`, pełna szerokość kolumny | spoczynek · wskazanie kursorem (`.dn-btn--sygnal:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · nieaktywny (brak wybranej operacji) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Uruchamia operację na wybranej pozycji; bez wybranej operacji kliknięcie sygnalizuje to komunikatem zamiast uruchomienia. Wynik trafia do Chat Window i Diff/Grep Panel, a przebieg do Execution Loop Window | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Tools Panel |
| Wskaźnik zakresu operacji | Mała etykieta nad listą | Przypomnienie, ile znaków i słów obejmie operacja | Tekst pomocniczy `--dn-fs-xs` | domyślny · aktualizowany na żywo | Aktualizuje się przy każdej zmianie zaznaczenia w edytorze | Nad grupami operacji | 3 | Widoczny po otwarciu Tools Panel |

---

### 3.5 Diff/Grep Panel

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
│                                    │
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
| Selektory wersji | Dwie rozwijane listy | Wybór porównywanej pary wersji | Małe przyciski z etykietą numeru wersji | spoczynek — ostatnie dwie wersje · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty | Zmiana przelicza widok różnicowy | Nagłówek panelu | 3 | Widoczne po otwarciu Diff/Grep Panel z menu `Narzędzia ▼` |
| Przełącznik trybu | Zakładki `.dn-zakladki--pigulki` | Wybór widoku: dwukolumnowy / scalony / wizualny | Mały segmentowany przełącznik | scalony (domyślny) · dwukolumnowy · wizualny · wskazanie kursorem (`.dn-zakladka:hover`) · ognisko (`.dn-zakladka:focus`, żeton `--dn-fokus`) · wybrany (`.dn-zakladka[aria-selected='true']`) | Przełącza układ prezentacji różnic | Nagłówek panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |
| Pole Grep | Pole wyszukiwania z przełącznikami | Wyszukiwanie wzorca w treści i historii | Pole `.dn-pole` + ikony-przełączniki (regex, wielkość liter, tryb semantyczny) | puste · ognisko (`.dn-pole-kontrolka:focus`, żeton `--dn-fokus`) · z wynikami · brak trafień (`.dn-pole-blad`, `--dn-blad-obrys`) | Enter uruchamia wyszukiwanie, podświetla trafienia w treści | Nagłówek panelu | 3 | Widoczne po otwarciu Diff/Grep Panel lub skrót klawiszowy wyszukiwania |
| Blok różnicy | Fragment tekstu z oznaczeniem koloru | Prezentacja pojedynczej zmiany | Średni blok, kolor tła zależny od typu zmiany | dodanie — tło `--dn-sukces-tlo`, obrys `--dn-sukces-obrys` · usunięcie — tło `--dn-blad-tlo`, przekreślenie · zmiana formatu — obrys `--dn-ostrzezenie-obrys` · wskazanie kursorem (`:hover`) · wybrany do akceptacji/odrzucenia (`[aria-selected='true']`) | Kliknięcie zaznacza blok do akceptacji lub odrzucenia | Ciało panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |
| Przycisk „Akceptuj zmianę” | Mały przycisk `.dn-btn--zarys` | Zatwierdzenie pojedynczej różnicy do wersji roboczej | Mały przycisk przy bloku różnicy | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Wprowadza zmianę do Studio Editor | Przy każdym bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Przycisk „Odrzuć zmianę” | Mały przycisk `.dn-btn--duch` | Pominięcie różnicy | Mały przycisk przy bloku różnicy | spoczynek · wskazanie kursorem (`.dn-btn--duch:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Zachowuje treść sprzed zmiany w dokumencie roboczym | Przy każdym bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Ikona adnotacji | Ikona dymka (odpowiedz) | Dodanie komentarza do zmiany | Mała ikona | brak adnotacji · z adnotacją (wypełniona) | Otwiera małe pole tekstowe komentarza | Przy bloku różnicy | 3 | Kliknięcie bloku różnicy |
| Przycisk „Eksportuj raport zmian” | Przycisk `.dn-btn--zarys`, pełna szerokość kolumny | Zapis widoku różnicowego jako dokumentu | Średni przycisk | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Generuje dokument raportu redakcji i udostępnia go do pobrania lub do Library | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Diff/Grep Panel albo polecenie języka naturalnego w Chat Window |
| Pasek statystyki | Wiersz podsumowania | Liczbowe podsumowanie zmian | Mały pasek tekstowy | domyślny, aktualizowany na żywo | — (informacyjny) | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Diff/Grep Panel |

---

### 3.6 Session Repository

| Aspekt | Wartość |
|---|---|
| Typologia | Repozytorium i biblioteka |
| Waga wizualna | Kolumna boczna wąska, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Historia wersji odrębna per karta sesji; narasta w toku sesji, nic nie jest usuwane bez decyzji użytkownika |

**Zawartość i pełny arsenał funkcji.**

- Chronologiczna lista wersji dokumentu ze znacznikiem czasu i autorem zmiany (`uzytkownik` lub `model`).
- Etykietowanie wersji własną nazwą (na przykład „wersja do akceptacji klienta”) i oznaczanie wersji kluczowych.
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
│   „korekta stylu”      │
├────────────────────────┤
│ ○ Wersja 6 · 11:47   ⋯ │
│   użytkownik           │
├────────────────────────┤
│ ○ Wersja 5 · 11:20   ⋯ │
│   🏷 „do akceptacji”    │
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
| Filtr historii | Rozwijana lista | Zawężenie listy wersji wg autora lub etykiety | Mały przycisk z etykietą | spoczynek — wszystkie · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty · aktywny filtr (`[aria-pressed='true']`) | Przelicza widoczną listę wersji | Góra panelu | 3 | Widoczny po otwarciu Session Repository |
| Pozycja wersji | Wiersz karty `.dn-karta--klikalna` | Reprezentuje jedną zapisaną wersję dokumentu | Średni wiersz z kropką stanu, znacznikiem czasu, autorem | bieżąca — kropka `.dn-kropka--sygnal`, wyróżnienie `.dn-karta--wybrana` · archiwalna — kropka `.dn-kropka--neutralna` · wskazanie kursorem (`.dn-karta--klikalna:hover`) · wciśnięcie (`.dn-karta--klikalna:active`) | Rozwija zestaw akcji (Podgląd / Przywróć / Porównaj / Rozgałęź) | Lista główna panelu | 3 | Widoczna po otwarciu Session Repository znacznikiem wersji w pasku statusu |
| Etykieta wersji | Mała plakietka `.dn-plakietka` | Własna nazwa nadana wersji przez użytkownika | Mała pigułka z ikoną (gwiazdka) | brak etykiety · z etykietą | Kliknięcie pozwala edytować treść etykiety | Pod znacznikiem czasu w pozycji wersji | 3 | Menu `⋯` pozycji wersji |
| Menu „⋯” | Ikona rozwijanego menu (wiecej) `.dn-btn-ikona` | Dodatkowe akcje na wersji (etykieta, eksport, odwołanie do wersji, usunięcie z listy lokalnie) | Mała ikona 24 px | spoczynek · wskazanie kursorem (`.dn-btn-ikona:hover`) · ognisko (`.dn-btn-ikona:focus`, żeton `--dn-fokus`) · rozwinięte (`[aria-pressed='true']`) | Otwiera menu kontekstowe pozycji | Prawa krawędź pozycji wersji | 3 | Kliknięcie ikony menu przy pozycji wersji |
| Przycisk „Przywróć” | Przycisk `.dn-btn--zarys`, mały | Ustawienie wybranej wersji jako bieżącej | Mały przycisk tekstowy | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Tworzy nową wersję na końcu listy z treścią wersji przywróconej | Zestaw akcji pozycji wersji | 3 | Rozwinięcie zestawu akcji pozycji wersji |
| Przycisk „Porównaj” | Przycisk `.dn-btn--duch`, mały | Otwarcie Diff/Grep Panel dla tej wersji względem bieżącej | Mały przycisk tekstowy | spoczynek · wskazanie kursorem (`.dn-btn--duch:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Otwiera Diff/Grep Panel z wypełnioną parą wersji | Zestaw akcji pozycji wersji | 3 | Rozwinięcie zestawu akcji pozycji wersji |
| Przycisk „Rozgałęź” | Przycisk `.dn-btn--duch`, mały | Utworzenie wariantu dokumentu od tej wersji | Mały przycisk tekstowy | spoczynek · wskazanie kursorem (`.dn-btn--duch:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Tworzy gałąź dokumentu i otwiera ją jako bieżącą treść edytora | Zestaw akcji pozycji wersji | 4 | Menu `⋯` pozycji wersji po włączeniu pracy na gałęziach w konfiguracji roli albo polecenie języka naturalnego w Chat Window |
| Przycisk „Eksportuj historię” | Przycisk `.dn-btn--zarys`, pełna szerokość kolumny | Pobranie całej historii wersji jako archiwum | Średni przycisk | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Generuje archiwum z manifestem i udostępnia plik do pobrania | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Session Repository albo polecenie języka naturalnego w Chat Window |
| Przycisk „Paczka redakcyjna” | Przycisk `.dn-btn--zarys`, pełna szerokość kolumny | Spakowanie dokumentu, historii, raportu zmian i adnotacji | Średni przycisk | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Generuje archiwum przekazania i udostępnia je do pobrania lub do Library | Zamknięcie kolumny panelu | 4 | Menu kebab (⋮) Session Repository albo polecenie języka naturalnego w Chat Window |

---

### 3.7 Preview Window

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
- Profile eksportu — nazwane zestawy ustawień (na przykład „PDF do druku”, „DOCX dla klienta”).
- Podgląd i konfiguracja wydruku: marginesy, orientacja, nagłówek, stopka, numeracja, skala.
- Znak wodny „wersja robocza” sterowany ustawieniem konfiguracyjnym podglądu, bez wpływu na treść dokumentu.
- Warsztat PDF dostępny z menu eksportu: scalanie i dzielenie, porządkowanie stron, kompresja i optymalizacja, znak wodny i stemplowanie, formularze AcroForm, numeracja Bates, ekstrakcja obrazów i załączników, drzewo zakładek.
- Zabezpieczenie dokumentu dostępne z menu eksportu: szyfrowanie i uprawnienia PDF, podpis cyfrowy i jego weryfikacja, redakcja poufności, usuwanie metadanych, wykrywanie danych wrażliwych.

**Makieta tekstowa.**

```
┌─ Preview Window ─────────────────⋮─┐
│ [ Format: PDF ▼ ]  [ Widok ▼ ]     │
├────────────────────────────────────┤
│                                    │
│    ┌───────────────────────────┐   │
│    │ wyrenderowana strona      │   │
│    │ dokumentu w formacie      │   │
│    │ docelowym                 │   │
│    └───────────────────────────┘   │
│                                    │
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
| Selektor formatu | Rozwijana lista | Wybór formatu renderowania podglądu | Mały przycisk z etykietą | spoczynek — PDF · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty | Przerenderowuje dokument w wybranym formacie | Nagłówek okna | 2 | Kliknięcie znacznika formatu docelowego w pasku kontekstu; zwija się po wyborze |
| Przełącznik widoku | Para przycisków `.dn-zakladki--pigulki` | Strona pojedyncza / przewijanie ciągłe | Mały segmentowany przełącznik | strona (domyślny) · ciągły · wskazanie kursorem (`.dn-zakladka:hover`) · ognisko (`.dn-zakladka:focus`, żeton `--dn-fokus`) · wybrany (`.dn-zakladka[aria-selected='true']`) | Zmienia sposób prezentacji wielostronicowego dokumentu | Nagłówek okna | 2 | Widoczny po otwarciu Preview Window |
| Suwak powiększenia | Rozwijana lista procentowa | Zmiana skali podglądu | Mały przycisk z etykietą procentową | 100% (domyślny) · niestandardowy · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) | Skaluje renderowaną stronę | Nagłówek okna | 2 | Widoczny po otwarciu Preview Window |
| Obszar renderowania | Wyrenderowana treść dokumentu | Wizualna prezentacja wyniku | Duży, centralny blok | ładowanie (`.dn-spinner`, `[aria-busy='true']`) · wyrenderowany · pusty stan — brak dokumentu (`.dn-pusty-stan`) | Przewijanie, powiększanie; brak edycji bezpośredniej | Środek okna | 2 | Widoczny po otwarciu Preview Window znacznikiem formatu docelowego |
| Przycisk „Widok kolumnowy” | Przycisk `.dn-btn--zarys` | Zestawienie Studio Editor i Preview w sąsiadujących kolumnach | Mały przycisk tekstowy | wyłączony · włączony — podświetlony (`.dn-btn--wybrany`, `[aria-pressed='true']`) · wskazanie kursorem (`.dn-btn--zarys:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) | Dzieli obszar roboczy modułu na dwie sąsiadujące kolumny z synchronicznym przewijaniem | Zamknięcie kolumny okna | 2 | Widoczny po otwarciu Preview Window |
| Przycisk „Eksportuj” | Przycisk `.dn-btn--sygnal` z rozwijanym menu formatów i profili | Pobranie dokumentu w formacie docelowym | Średni przycisk CTA | spoczynek · wskazanie kursorem (`.dn-btn--sygnal:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · rozwinięte menu · ładowanie eksportu (`.dn-spinner`, `[aria-busy='true']`) | Generuje plik według wybranego profilu i uruchamia pobranie | Zamknięcie kolumny okna | 3 | Pozycja `Eksport` w menu `Operacje ▼` lub w Preview Window |
| Menu „Warsztat PDF” | Rozwijane menu operacji PDF | Scalanie, dzielenie, porządkowanie stron, kompresja, stemplowanie, formularze, Bates, zakładki | Średni przycisk z rozwinięciem | spoczynek · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięte · ładowanie operacji (`.dn-spinner`, `[aria-busy='true']`) | Wykonuje operację na bieżącym dokumencie i zapisuje wynik jako nową wersję | Zamknięcie kolumny okna | 4 | Wyszukiwarka funkcji, skrót klawiszowy albo polecenie języka naturalnego w Chat Window; menu widoczne po włączeniu operacji zaawansowanych w konfiguracji roli |
| Menu „Zabezpiecz” | Rozwijane menu operacji ochronnych | Szyfrowanie, podpis, redakcja poufności, usuwanie metadanych | Średni przycisk z rozwinięciem | spoczynek · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięte · ładowanie operacji (`.dn-spinner`, `[aria-busy='true']`) | Nakłada wybrane zabezpieczenie i zapisuje wynik jako nową wersję | Zamknięcie kolumny okna | 4 | Wyszukiwarka funkcji albo polecenie języka naturalnego w Chat Window; menu widoczne po włączeniu operacji zaawansowanych w konfiguracji roli |
| Przycisk „Wyślij do Library” | Przycisk `.dn-btn--zarys` | Przekazanie bieżącej wersji do modułu Library | Mały przycisk tekstowy | spoczynek · wskazanie kursorem (`.dn-btn--zarys:hover`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · wysłano — potwierdzenie `.dn-toast--sukces` | Tworzy artefakt w Library, gdy powiązanie skonfigurowane | Zamknięcie kolumny okna | 3 | Pozycja `Wyślij do Library` w menu `Operacje ▼` lub w Preview Window |

---

### 3.8 Ingest/OCR Panel

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
| Selektor źródła | Zestaw przycisków źródeł | Wybór drogi wczytania materiału | Mały segmentowany przełącznik | pliki (domyślny) · folder · URL · skaner · e-mail · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · wybrany (`[aria-pressed='true']`) | Otwiera właściwy dialog wczytania lub pole adresu | Góra panelu | 3 | Widoczny po otwarciu Ingest/OCR Panel pozycją `Wczytaj i cyfryzuj` w menu `Operacje ▼` |
| Pozycja kolejki | Wiersz `.dn-karta--klikalna` z nazwą pliku | Reprezentuje jeden materiał wejściowy | Średni wiersz z ikoną stanu i postępem | oczekuje · przetwarzanie (`.dn-spinner`, `[aria-busy='true']` z procentem) · gotowy (`.dn-kropka--sukces`) · błąd odczytu (`.dn-kropka--blad`, `--dn-blad-tekst`) · wskazanie kursorem (`.dn-karta--klikalna:hover`) | Kliknięcie otwiera podgląd warstwy tekstowej tej pozycji | Kolejka wczytywania | 3 | Widoczna po otwarciu Ingest/OCR Panel |
| Selektor silnika OCR | Rozwijana lista | Wybór silnika rozpoznawania | Mały przycisk z etykietą | lokalny (domyślny) · usługa zewnętrzna · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty | Zmienia silnik dla kolejnych pozycji kolejki | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu albo polecenie języka naturalnego w Chat Window |
| Selektor języków | Rozwijana lista wielokrotnego wyboru | Wybór języków rozpoznawania | Mały przycisk z listą znaczników `.dn-plakietka` | pol + eng (domyślny) · zestaw własny · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) · rozwinięty | Zmienia zestaw danych językowych dla rozpoznawania | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu albo polecenie języka naturalnego w Chat Window |
| Przełączniki pre-OCR | Zestaw pól wyboru `.dn-check` | Sterowanie czyszczeniem obrazu przed rozpoznaniem | Małe pola wyboru w rzędach | włączone (`[aria-pressed='true']`) · wyłączone · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) | Zmiana przelicza przetwarzanie pozycji oczekujących | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu |
| Przełączniki rozpoznania układu | Zestaw pól wyboru `.dn-check` | Sterowanie wykrywaniem kolumn, tabel, przypisów | Małe pola wyboru w rzędach | włączone (`[aria-pressed='true']`) · wyłączone · wskazanie kursorem (`:hover`) · ognisko (`:focus`, żeton `--dn-fokus`) | Zmienia zakres odtwarzanej struktury dokumentu | Ustawienia panelu | 4 | Rozwinięcie ustawień zaawansowanych panelu |
| Podgląd warstwy tekstowej | Obraz skanu z nałożonym tekstem | Kontrola jakości rozpoznania | Średni blok z przełącznikiem przezroczystości | warstwa ukryta · warstwa widoczna · korekta | Kliknięcie słowa otwiera pole korekty rozpoznanego tekstu | Środek panelu | 3 | Kliknięcie pozycji kolejki; korekta rozpoznania stanowi warstwę 4 |
| Przycisk „Przekaż do edytora” | Główny CTA panelu | Przyjęcie wyniku cyfryzacji jako dokumentu roboczego | Duży przycisk `.dn-btn--sygnal`, pełna szerokość kolumny | spoczynek · wskazanie kursorem (`.dn-btn--sygnal:hover`) · wciśnięcie (`.dn-btn:active`) · ognisko (`.dn-btn:focus`, żeton `--dn-fokus`) · nieaktywny (kolejka bez pozycji gotowej) · ładowanie (`.dn-spinner`, `[aria-busy='true']`) | Otwiera dokument w Studio Editor i zapisuje pierwszą wersję w Session Repository | Zamknięcie kolumny panelu | 3 | Widoczny po otwarciu Ingest/OCR Panel |

---

## 4. Przepływy pracy

### 4.1 Przepływ podstawowy — redakcja fragmentu z korektą AI

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

### 4.2 Przepływ cyfryzacji — skan do dokumentu edytowalnego

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

### 4.3 Przepływ rozszerzony — praca dwujęzyczna z modułem Translate

```
Studio Editor (dokument źródłowy)
        │  zaznaczenie fragmentu do wersji dwujęzycznej
        ▼
Tools Panel ── operacja „Tłumaczenie” ──► współdzielenie operacji kontekstowych AI
        │                                  (powiązanie Studio ◄──► Translate,
        │                                   konfigurowalne — rozdz. 5.3)
        ▼
Translate — Translation Panels (wynik równoległy w językach docelowych)
        │
        ▼
Powrót do Studio Editor jako wersja wielojęzyczna dokumentu
        │
        ▼
Session Repository — zapis wersji z etykietą „wielojęzyczna”
```

### 4.4 Przepływ rozszerzony — redakcja raportu przekazanego z Research

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

### 4.5 Przepływ zbiorczy — wiele dokumentów równolegle

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

### 5.1 Model stanów sesji modułu

```
                    ┌────────────────┐
                    │  Pusta sesja   │  (brak dokumentu)
                    └───────┬────────┘
                            │ wczytanie / cyfryzacja / nowy dokument
                            ▼
                    ┌────────────────┐
        ┌──────────►│  Edycja        │◄───────────────┐
        │           └───────┬────────┘                │
        │                   │ zlecenie operacji AI    │ odrzucenie zmiany
        │                   ▼                         │
        │           ┌────────────────┐                │
        │           │  Pętla         │  (Execution    │
        │           │  wykonawcza    │   Loop Window) │
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

### 5.2 Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Studio |
|---|---|
| `sesja`, `karta_sesji` | Nośnik kontekstu bieżącej pracy nad dokumentem |
| `wiadomosc`, `zalacznik_wiadomosci` | Historia Chat Window i pliki załączone do poleceń |
| `artefakt` | Dokument roboczy jako plik z odwołaniem w bazie |
| `wersja_artefaktu` | Każda pozycja w Session Repository; pole `autor` rozróżnia zmianę `uzytkownik`/`model`; nośnik etykiet, gałęzi i adnotacji |
| `kanal_modelu` | Model bazowy przypisany karcie sesji, wybierany w Chat Window; źródło embeddingów wyszukiwania semantycznego oraz przetwarzania mowy |
| `ustawienie` | Parametry modułu podlegające zasadzie „brak ustawienia = wartość domyślna” (interwał autozapisu, prompty operacji, profile eksportu, mapa skrótów) |
| `profil_izolacji`, `regula_izolacji_kontekstu` | Zakres współdzielenia historii, pamięci i kontekstu między kartami i modułami oraz zakres widoczności artefaktu |
| `powiazanie_komponentu` | Jawne powiązania Studio ◄──► Translate, Studio ───► Library, Studio ◄──► Research, Design ───► Studio |

### 5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Studio

| Punkt izolacji | Stan wyjściowy (domyślny) | Co można skonfigurować | Gdzie |
|---|---|---|---|
| Historia Chat Window | Odrębna per karta sesji | Współdzielenie historii między kartami pracującymi nad powiązanymi dokumentami | Okno konfiguracji punktów izolacji, poziom „karta sesji” lub „moduł” |
| Przebieg Execution Loop Window | Odrębny per karta sesji | Współdzielenie kolejki zadań między kartami przy operacjach wsadowych | Okno konfiguracji punktów izolacji, poziom „moduł” lub „projekt” |
| Pamięć kontekstowa | Odrębna per karta sesji | Współdzielenie pamięci projektu, gdy Studio używany w ramach projektu Workspace | Poziom zasięgu „projekt” |
| Para modułów Studio–Translate | Operacje kontekstowe AI nie są współdzielone domyślnie | Współdzielenie mechanizmu operacji kontekstowych AI przy pracy dwujęzycznej | Poziom zasięgu „para modułów” |
| Repozytorium źródłowe (Library) | Brak automatycznego powiązania | Ustanowienie Library jako źródła wejściowego i celu eksportu artefaktów | Okno konfiguracji, powiązanie komponentu Studio ───► Library |
| Widoczność artefaktu | Zakres sesji | Rozszerzenie widoczności dokumentu do projektu lub zakresu współdzielonego | Okno konfiguracji punktów izolacji |
| Izolacja techniczna procesu sesji (8 zakresów) | Żaden zakres domyślnie nie jest aktywny — proces sesji nie jest ograniczany | Włączenie wybranych zakresów (np. odrębny katalog roboczy) dla sesji redagującej dokumenty poufne | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

Zgodnie z zasadą nadrzędną platformy żaden z powyższych punktów nie jest wymuszony — moduł działa w pełni funkcjonalnie bez jakiejkolwiek konfiguracji izolacji, a każdy punkt jest odwracalną, jawną decyzją Operatora.

### 5.4 Powiązania z innymi modułami

```
       Translate ◄──►  STUDIO  ───►  Library
     (operacje AI              (repozytorium źródłowe
      przy pracy                i cel dla artefaktów)
      dwujęzycznej,
      konfigurowalne)                  ▲
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
Użytkownik wczytuje dokument DOCX do Studio Editor, zaznacza akapit budzący wątpliwości stylistyczne, z paska pływającego wybiera „Korekta”. Koordynator rozkłada zlecenie na zadania widoczne w Execution Loop Window, odpowiedź Wykonawcy trafia do Chat Window, a Diff/Grep Panel prezentuje różnicę. Użytkownik akceptuje zmianę fragmentarycznie — tylko dwie z trzech zaproponowanych poprawek — pozostawiając trzecią bez zmian. Session Repository zapisuje nową wersję.

**Scenariusz 2 — cyfryzacja umowy ze skanu.**
Kancelaria wczytuje osiemnastostronicowy skan umowy do Ingest/OCR Panel. Panel prostuje skos, odszumia obrazy i rozpoznaje tekst w językach polskim i angielskim. Execution Loop Window prowadzi kolejkę stron jako zadania pętli i ponawia rozpoznanie dwóch stron o niskiej pewności. Użytkownik koryguje kilka słów w podglądzie warstwy tekstowej, przekazuje wynik do Studio Editor, a z Preview Window eksportuje przeszukiwalny PDF z numeracją Bates.

**Scenariusz 3 — praca dwujęzyczna.**
Zespół lokalizacyjny ustanawia w oknie konfiguracji współdzielenie operacji kontekstowych AI między Studio a Translate na poziomie pary modułów. Redakcja dokumentu źródłowego w Studio Editor odzwierciedla się w tłumaczeniu zaznaczenia bezpośrednio z Tools Panel, bez opuszczania modułu Studio.

**Scenariusz 4 — redakcja raportu z Research.**
Analityk kończy pracę badawczą w module Research; z Report Builder przekazuje skompletowany raport do Studio. W Studio Editor prowadzi ostateczną redakcję językową, generuje streszczenie zarządcze operacją z Tools Panel, porównuje wynik z materiałem wejściowym w Diff/Grep Panel, a następnie eksportuje finalny PDF z Preview Window.

**Scenariusz 5 — praca równoległa nad wieloma dokumentami.**
Redaktor otwiera trzy karty sesji modułu Studio, każdą z osobnym dokumentem klienta. Domyślna izolacja kontekstu zapewnia, że polecenia wydane w jednej karcie nie wpływają na kontekst pozostałych. Dla dwóch kart dotyczących tego samego projektu redaktor włącza współdzielenie pamięci na poziomie „projekt”, a operację korekty uruchamia wsadowo — Execution Loop Window prowadzi ją jako jedną pętlę obejmującą oba dokumenty.

**Scenariusz 6 — powrót do wcześniejszej redakcji.**
Po serii zmian klient prosi o powrót do wersji sprzed dwóch dni. Użytkownik odnajduje ją w Session Repository po etykiecie „do akceptacji klienta”, otwiera podgląd, a następnie przywraca ją jako wersję bieżącą — operacja tworzy nową pozycję na końcu historii, zachowując wszystkie wersje pośrednie.

**Scenariusz 7 — przekazanie paczki redakcyjnej.**
Po zamknięciu redakcji użytkownik usuwa metadane dokumentu, nakłada podpis cyfrowy z Preview Window i generuje z Session Repository paczkę redakcyjną: dokument finalny, pełną historię wersji, raport zmian i adnotacje w jednym archiwum przekazania.

---

## 7. Katalog funkcji i narzędzi

Katalog obejmuje dwanaście rodzin funkcji. Każda pozycja opisuje nazwę, działanie oraz zależności techniczne (biblioteki, formaty, integracje). Nazwy bibliotek podano jako konkretne komponenty w Go — języku zaplecza platformy — wraz z rezerwą; narzędzia zewnętrzne (LibreOffice, Pandoc, Tesseract) wywoływane są jako procesy pomocnicze przez moduł Terminal lub bezpośrednio przez warstwę serwera.

### 7.1 Rodzina A — Wczytywanie, formaty i cyfryzacja

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Uniwersalny importer formatów | Rozpoznaje typ pliku po sygnaturze i wczytuje PDF/DOCX/ODT/RTF/TXT/MD/HTML/EPUB/XLSX/CSV do wspólnego modelu dokumentu | `h2non/filetype` (detekcja MIME); parsery per format; model dokumentu wewnętrzny |
| 2 | Ekstraktor tekstu z PDF | Wyciąga tekst, strukturę akapitów i pozycje z PDF tekstowego, zachowując kolejność czytania | `ledongthuc/pdf` lub `gen2brain/go-fitz` (MuPDF, wyższa jakość układu); format PDF |
| 3 | Silnik OCR skanów i zdjęć | Rozpoznaje tekst z obrazów i PDF-skanów wielojęzycznie, zwraca warstwę tekstową z pozycjami słów | `otiai10/gosseract/v2` → Tesseract (`libtesseract`, `tessdata` — pol, eng, deu…); dla PDF-skanu render przez go-fitz; usługa chmurowa (Google Vision / Azure) sterowana jawnym kluczem w oknie konfiguracji |
| 4 | Wstępne czyszczenie obrazu | Prostuje skos, usuwa szum, binaryzuje, przycina marginesy, zwiększa kontrast przed rozpoznaniem | `disintegration/imaging` lub `anthonynsimon/bild`; formaty PNG/JPG/TIFF |
| 5 | Rozpoznanie układu | Wykrywa kolumny, tabele, nagłówki, stopki, przypisy — odtwarza strukturę logiczną, nie tylko surowy tekst | go-fitz (bloki tekstu z bbox); reguły heurystyczne i model AI |
| 6 | Import ze schowka z rozpoznaniem struktury | Wkleja HTML/RTF/tekst ze schowka i odtwarza nagłówki, listy, tabele | Parser HTML `net/html`; konwersja HTML → model dokumentu |
| 7 | Import z URL i zrzut strony do dokumentu | Pobiera stronę WWW, oczyszcza z nawigacji i reklam, zapisuje jako czytelny dokument | `go-shiori/go-readability`; powiązanie z modułem Browser |
| 8 | Import e-maila jako dokumentu | Wczytuje EML/MSG, wyodrębnia treść, załączniki i nagłówki do dokumentu roboczego | `net/mail`, `emersion/go-message`, `mnako/letters`; formaty EML/MSG |
| 9 | Kolejka wsadowego wczytywania | Wczytuje wiele plików naraz (folder, ZIP) z jednakowym potokiem rozpoznawania i normalizacji | `archive/zip`; kolejka zadań (powiązanie z Automations) |

### 7.2 Rodzina B — Edycja i tworzenie treści

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Edytor WYSIWYG dokumentu | Edycja bogata: style znaku i akapitu, nagłówki poziomów od pierwszego do szóstego, listy, cytaty, tabele, obrazy, linki, separatory | Edytor po stronie klienta (ProseMirror/TipTap w warstwie web); model dokumentu ↔ serwer |
| 2 | Edytor Markdown z podglądem | Tryb źródłowy Markdown/MDX z podglądem na żywo, tabelami GFM, przypisami, blokami kodu | `yuin/goldmark` (CommonMark + GFM, tabele, footnotes, task-list) do renderu |
| 3 | Szablony dokumentów | Predefiniowane układy (pismo, umowa, raport, notatka, oferta) z polami do wypełnienia | `text/template`/`html/template`; biblioteka szablonów w Library |
| 4 | Style typografii i motywy dokumentu | Zestawy stylów (czcionki, interlinia, marginesy) stosowane jednym kliknięciem, spójne z systemem wizualnym | Definicje stylów w modelu dokumentu; tokeny typografii platformy |
| 5 | Tabele zaawansowane | Wstawianie i edycja tabel: scalanie komórek, nagłówki, sortowanie, suma kolumny | Model tabeli w dokumencie; render do DOCX/PDF/MD |
| 6 | Przypisy, bibliografia, cytowania | Zarządza przypisami dolnymi i końcowymi oraz cytowaniami w stylach APA, MLA, Chicago | Parser stylów cytowania (CSL); baza źródeł (powiązanie z Research/Library) |
| 7 | Spis treści i numeracja automatyczna | Generuje spis treści z nagłówków, numeruje strony, rozdziały, ilustracje | Analiza drzewa nagłówków; render paginacji |
| 8 | Kotwice, zakładki i nawigacja | Zakładki w dokumencie i minimapa nagłówków do szybkiego skoku | Indeks nagłówków; minimapa w oknie edytora |
| 9 | Komentarze i adnotacje redakcyjne | Komentarze przypięte do fragmentu, wątki, oznaczenie „rozwiązane” | Model adnotacji w `wersja_artefaktu`; render do PDF (warstwa adnotacji) |
| 10 | Śledzenie zmian | Rejestruje wstawienia i usunięcia z autorem, akceptacja i odrzucenie pojedynczo | Warstwa zmian nad modelem dokumentu; integracja z rodziną E |
| 11 | Liczniki i statystyki dokumentu | Słowa, znaki, zdania, akapity, strony, czas czytania, terminy unikalne | Analizator tekstu (segmentacja zdań i słów) |

### 7.3 Rodzina C — Operacje kontekstowe AI na treści

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Korekta ortograficzno-gramatyczna | Poprawia błędy pisowni, gramatyki i interpunkcji na zaznaczeniu lub całości | `kanal_modelu`; warstwa lokalna `hunspell`/`client9/misspell` jako szybkie rozpoznanie wstępne |
| 2 | Zmiana stylu i rejestru | Przełącza styl (formalny/nieformalny/techniczny/marketingowy) i rejestr językowy | `kanal_modelu`; szablony promptów operacji (konfigurowalne) |
| 3 | Zmiana tonu | Nadaje ton (neutralny/perswazyjny/empatyczny/stanowczy) bez zmiany treści | `kanal_modelu` |
| 4 | Skracanie i rozwijanie | Kondensuje lub rozszerza fragment do zadanej długości lub proporcji | `kanal_modelu`; wskaźnik długości („Liczniki i statystyki dokumentu”) |
| 5 | Streszczenia wielopoziomowe | Streszczenie jednozdaniowe, akapitowe, wypunktowane, streszczenie zarządcze na wstępie | `kanal_modelu` |
| 6 | Ekstrakcja strukturalna | Wyciąga tezy, działania, decyzje, terminy, encje (osoby, kwoty, daty) | `kanal_modelu`; reguły rozpoznawania encji |
| 7 | Przepisanie strukturalne | Akapit → lista, lista → tabela, proza → punkty, generacja tytułów i podtytułów | `kanal_modelu`; operacje na modelu dokumentu |
| 8 | Parafraza | Przepisuje frazę innym sformułowaniem, z wariantami do wyboru | `kanal_modelu` |
| 9 | Uproszczenie języka | Upraszcza żargon, skraca zdania, podnosi czytelność do zadanego poziomu | `kanal_modelu` + metryki czytelności („Kontrola czytelności i jakości”) |
| 10 | Tłumaczenie kontekstowe | Tłumaczy zaznaczenie lub całość na wybrany język, wersja dwujęzyczna równoległa | `kanal_modelu`; powiązanie `Studio ◄──► Translate` (konfigurowalne) |
| 11 | Wzbogacanie treści | Dodaje kontekst, przykłady, kontrargumenty, pytania kontrolne, definicje | `kanal_modelu` |
| 12 | Pytania i odpowiedzi nad dokumentem (RAG) | Odpowiada na pytania w oparciu o treść dokumentu z cytowaniem miejsca | `kanal_modelu`; indeks wektorowy (`blevesearch/bleve` + embeddingi) lub kontekst pełnotekstowy |
| 13 | Kontrola czytelności i jakości | Wskaźniki (Flesch, Gunning fog, długość zdań), powtórzenia, strona bierna, spójność terminologii | Metryki liczone lokalnie; `kanal_modelu` dla sugestii |
| 14 | Weryfikacja i wskazanie źródeł | Oznacza fragmenty wymagające źródła i potencjalne nieścisłości, porównuje z materiałem z Library | `kanal_modelu`; powiązanie z Library/Research |
| 15 | Operacje własne | Zapis własnych promptów jako nazwanych narzędzi wielokrotnego użytku | Szablony promptów w `ustawienie`; okno konfiguracji |

### 7.4 Rodzina D — Konwersje i eksport

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Konwerter uniwersalny dokumentów | Konwertuje między MD/DOCX/HTML/LaTeX/EPUB/RTF/ODT/PDF/TXT w dowolnym kierunku | `pandoc` wywoływany jako proces zewnętrzny (najszersza matryca formatów) |
| 2 | Renderowanie do PDF | Składa dokument do PDF z typografią, paginacją, nagłówkiem i stopką, spisem treści | `go-pdf/fpdf` lub `signintech/gopdf`; ścieżka HTML→PDF przez `chromedp` dla wiernego układu |
| 3 | Eksport do DOCX | Zapisuje dokument jako edytowalny DOCX ze stylami i tabelami | `unidoc/unioffice` (klucz komercyjny jawny) lub `fumiama/go-docx`; rezerwa: Pandoc |
| 4 | Eksport do arkusza | Zapisuje tabele dokumentu jako arkusz XLSX lub CSV | `qax-os/excelize` (XLSX), `encoding/csv` |
| 5 | Eksport do EPUB | Składa dokument wielorozdziałowy w e-book | `bmaupin/go-epub` + goldmark do treści |
| 6 | Eksport do HTML | Renderuje dokument jako samodzielny HTML z osadzonymi stylami i obrazami | goldmark/`html/template`; inline CSS |
| 7 | Podgląd i konfiguracja wydruku | Marginesy, orientacja, nagłówek i stopka, numeracja, skala, podgląd stron | Warstwa layoutu druku; render przez „Renderowanie do PDF” |
| 8 | Konwersja obraz → dokument | Zamienia skan lub zdjęcie w przeszukiwalny PDF albo edytowalny DOCX/TXT | „Silnik OCR skanów i zdjęć” (rozpoznanie tekstu) + „Renderowanie do PDF”/„Eksport do DOCX”; warstwa tekstowa nad obrazem |
| 9 | Eksport z profilami | Nazwane profile eksportu (np. „PDF do druku”, „DOCX dla klienta”) z zapamiętanymi ustawieniami | Profile w `ustawienie` (konfigurowalne) |

### 7.5 Rodzina E — Diff, porównania i wyszukiwanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Diff tekstowy słowo po słowie | Różnica dwóch wersji z oznaczeniem dodań, usunięć i zmian, tryb scalony i dwukolumnowy | `sergi/go-diff` (diff-match-patch) lub `hexops/gotextdiff` (Myers) |
| 2 | Diff dowolnej pary wersji | Porównuje nie tylko sąsiednie, lecz dowolne dwie pozycje z historii | „Diff tekstowy słowo po słowie” + `wersja_artefaktu` (Session Repository) |
| 3 | Diff wizualny PDF i układu | Porównuje wyrenderowany wygląd — nakładka i podświetlenie zmian graficznych | Render obu wersji przez go-fitz → porównanie pikselowe (`imaging`) |
| 4 | Grep i wyszukiwanie regex | Wyszukuje wzorzec (w tym regex) w dokumencie i w całej historii wersji | `regexp` stdlib; podświetlanie trafień |
| 5 | Znajdź i zamień wsadowo | Zamiana wzorca w dokumencie, podgląd przed zatwierdzeniem, cofnięcie | `regexp`; operacja na modelu dokumentu z migawką |
| 6 | Wyszukiwanie semantyczne | Znajduje fragmenty znaczeniowo bliskie zapytaniu, nie tylko dosłowne | `blevesearch/bleve` + embeddingi z `kanal_modelu` |
| 7 | Porównanie ze źródłem | Zestawia dokument roboczy z materiałem wejściowym (Research/Library) i wskazuje rozbieżności | „Diff tekstowy słowo po słowie” + powiązanie międzymodułowe |
| 8 | Filtrowanie i statystyka różnic | Filtruje różnice (dodania, usunięcia, format), liczy zmienione słowa i znaki | Agregacja wyniku „Diff tekstowy słowo po słowie” |
| 9 | Adnotacje i eksport raportu zmian | Komentarz przy różnicy; eksport całości różnic jako osobny dokument redakcji | Model adnotacji + „Renderowanie do PDF” (PDF) |

### 7.6 Rodzina F — Wersjonowanie i historia

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Historia wersji sesyjna | Chronologiczna lista wersji ze znacznikiem czasu i autorem (`uzytkownik`/`model`) | `wersja_artefaktu`; nic nie znika bez decyzji użytkownika |
| 2 | Etykiety i kamienie milowe | Nazwa własna wersji („do akceptacji klienta”), oznaczenie wersji kluczowych | Pole etykiety w `wersja_artefaktu` |
| 3 | Przywracanie bez utraty | Ustawia starszą wersję jako bieżącą, tworząc nową pozycję; historia pozostaje nienaruszona | „Historia wersji sesyjna”; polityka append-only |
| 4 | Gałęzie dokumentu | Tworzy równoległy wariant od wybranej wersji (np. dwie redakcje) | Model gałęzi nad `artefakt`; scalanie przez „Diff tekstowy słowo po słowie” |
| 5 | Scalanie wariantów | Łączy dwie gałęzie z rozwiązywaniem konfliktów per fragment | „Diff tekstowy słowo po słowie”/„Diff dowolnej pary wersji” + interfejs rozwiązywania konfliktów |
| 6 | Filtrowanie i przegląd historii | Filtr: zmiany użytkownika, zmiany AI, wersje etykietowane; podgląd bez opuszczania modułu | Zapytania nad `wersja_artefaktu` |
| 7 | Eksport i archiwizacja historii | Pobiera pojedynczą wersję lub całą historię jako archiwum ZIP z manifestem | `archive/zip`; manifest JSON |
| 8 | Autozapis i migawki | Automatyczny zapis w konfigurowalnym interwale; migawka przy każdej zaakceptowanej zmianie AI | Interwał w `ustawienie` (konfigurowalny) |

### 7.7 Rodzina G — Operacje na PDF

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Scalanie i dzielenie PDF | Łączy wiele PDF w jeden, dzieli po stronach, zakresach i zakładkach | `pdfcpu/pdfcpu` (merge, split, trim) |
| 2 | Porządkowanie stron | Zmiana kolejności, obrót, usuwanie, wstawianie, ekstrakcja stron | `pdfcpu` (rotate, remove, insert, extract) |
| 3 | Kompresja i optymalizacja | Zmniejsza rozmiar PDF przez kompresję obrazów i czyszczenie struktury | `pdfcpu` (optimize); rekompresja obrazów przez `imaging` |
| 4 | Znak wodny i stemplowanie | Nakłada znak wodny lub pieczęć tekstową albo graficzną | `pdfcpu` (stamp/watermark) |
| 5 | Formularze PDF | Wypełnianie i tworzenie pól formularza, odczyt i zapis danych AcroForm | `pdfcpu` (form fill/export) lub `unidoc/unipdf` |
| 6 | Numeracja Bates i nagłówki | Dodaje numerację prawną Bates oraz nagłówki i stopki na wszystkich stronach | `pdfcpu` (stamp z licznikiem) |
| 7 | Ekstrakcja obrazów i załączników | Wyciąga obrazy, czcionki i załączniki osadzone w PDF | `pdfcpu` (extract images/attachments) |
| 8 | Zakładki i nawigacja PDF | Tworzy i edytuje drzewo zakładek dokumentu | `pdfcpu` (bookmarks) / `unipdf` |

### 7.8 Rodzina H — Bezpieczeństwo i poufność dokumentu

Każda funkcja rodziny H jest jawnym, odwracalnym ustawieniem Operatora; żadna nie warunkuje pracy nad dokumentem.

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Szyfrowanie i hasło PDF | Nakłada i zdejmuje szyfrowanie oraz uprawnienia dokumentu | `pdfcpu` (encrypt/decrypt/permissions) |
| 2 | Podpis cyfrowy i pieczęć | Podpisuje PDF certyfikatem, weryfikuje podpis | `digitorus/pdfsign`; certyfikat użytkownika |
| 3 | Redakcja poufności | Trwale zamazuje wskazane fragmenty wraz z warstwą tekstową pod spodem | Nadpisanie warstwy tekstu + zaczernienie obrazu (pdfcpu + imaging) |
| 4 | Usuwanie metadanych | Czyści autora, historię, dane ukryte i komentarze przed wysyłką | `pdfcpu` (metadata); parsery DOCX i obrazów |
| 5 | Wykrywanie danych wrażliwych | Oznacza w treści numery identyfikacyjne, numery kart, adresy e-mail i dane osobowe do redakcji | Reguły regex + rozpoznawanie encji przez `kanal_modelu` |
| 6 | Kontrola dostępu do artefaktu | Ustawia zakres widoczności dokumentu (sesja, projekt, zakres współdzielony) | `profil_izolacji`, `regula_izolacji_kontekstu`; okno konfiguracji izolacji |

### 7.9 Rodzina I — Arkusze i dane tabelaryczne

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Odczyt i edycja arkusza | Otwiera XLSX/ODS/CSV, edytuje komórki, arkusze i proste formuły | `qax-os/excelize` (XLSX), `encoding/csv` |
| 2 | Tabela dokumentu ↔ arkusz | Konwertuje tabelę z dokumentu do arkusza i odwrotnie | Model tabeli („Tabele zaawansowane”) ↔ excelize |
| 3 | Ekstrakcja tabel z PDF | Wykrywa i wyciąga tabele z PDF do arkusza lub CSV | go-fitz (pozycje) + heurystyka siatki; rezerwa: `tabula-java` jako proces |
| 4 | Korespondencja seryjna | Generuje serię dokumentów z szablonu i wierszy arkusza | „Szablony dokumentów” (szablony) + „Odczyt i edycja arkusza”; pętla generacji (Automations) |
| 5 | Wizualizacja tabeli | Wstawia wykres z danych tabeli do dokumentu | excelize (chart) lub render SVG po stronie klienta |

### 7.10 Rodzina J — Kolaboracja i przekazywanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Przekazanie do i z Library | Zapisuje wersję jako artefakt w Library, wczytuje dokument źródłowy z Library | Powiązanie `Studio ───► Library`; `powiazanie_komponentu` |
| 2 | Odbiór raportu z Research | Wczytuje skompletowany raport z Report Builder do redakcji finalnej | Powiązanie `Research ◄──► Studio` |
| 3 | Praca dwujęzyczna z Translate | Współdzieli operacje kontekstowe przy wersji dwujęzycznej | Powiązanie `Studio ◄──► Translate` |
| 4 | Wstawianie zasobów z Design | Osadza grafiki i elementy wizualne z Assets Panel | Powiązanie `Design ───► Studio` |
| 5 | Odwołanie do wersji | Generuje odwołanie do konkretnej wersji dokumentu w obrębie platformy | `wersja_artefaktu` + warstwa udostępnień (konfigurowalna) |
| 6 | Eksport paczki redakcyjnej | Pakuje dokument, historię, raport zmian i adnotacje do jednego archiwum przekazania | „Eksport i archiwizacja historii” + „Adnotacje i eksport raportu zmian”; `archive/zip` |

### 7.11 Rodzina K — Produktywność i sterowanie

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Paleta poleceń | Wywołanie dowolnej operacji modułu z klawiatury (`Ctrl/Cmd + K`) | Rejestr operacji modułu; skrót globalny |
| 2 | Operacje z czatu przez „/” | Wywołanie operacji Tools Panel z poziomu Chat Window skrótem `/` | Menu operacji + Chat Window |
| 3 | Makra i łańcuchy operacji | Sekwencja operacji AI wykonywana jednym poleceniem (np. korekta → streszczenie → eksport) | Definicje łańcuchów w `ustawienie`; Execution Loop Window prowadzi przebieg; powiązanie z Automations |
| 4 | Operacje wsadowe na wielu dokumentach | Ta sama operacja na wielu kartach lub plikach naraz | Kolejka + „Kolejka wsadowego wczytywania”; Execution Loop Window; Automations |
| 5 | Skróty klawiszowe konfigurowalne | Pełna mapa skrótów, edytowalna przez Operatora | Mapa skrótów w `ustawienie` |
| 6 | Tryb skupienia | Ukrywa kolumny boczne, centruje wiersz aktywny | Ustawienie widoku edytora |
| 7 | Dyktowanie i odczyt na głos | Wprowadzanie tekstu głosem i odczyt dokumentu | Powiązanie z Assistant (Voice Console); przetwarzanie mowy przez `kanal_modelu` |

### 7.12 Rodzina L — Załączniki multimedialne i specjalne

| # | Nazwa funkcji | Co robi | Zależności |
|---|---|---|---|
| 1 | Transkrypcja audio i wideo do dokumentu | Zamienia nagranie w tekst z podziałem na mówców i znacznikami czasu | Przetwarzanie mowy przez `kanal_modelu` (Whisper); formaty audio i wideo (MP3, WAV, MP4) |
| 2 | Osadzanie i podpisy obrazów | Wstawia obrazy z podpisem i tekstem alternatywnym generowanym przez AI | `imaging`; `kanal_modelu` dla opisu |
| 3 | Kod QR i kody kreskowe | Generuje i wstawia kod QR lub kreskowy (np. odwołanie do wersji) | `boombuler/barcode` |
| 4 | Wzory matematyczne | Wstawia i renderuje wzory z notacji LaTeX | Render KaTeX/MathJax po stronie klienta; eksport przez Pandoc |
| 5 | Diagramy z tekstu | Renderuje diagramy z opisu tekstowego osadzone w dokumencie | Mermaid po stronie klienta; eksport do SVG/PNG |

### 7.13 Podstawa technologiczna modułu

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

### 7.14 Przebiegi wybranych funkcji i przypadki brzegowe

Trzy funkcje katalogu wymagają rozpisania przebiegu i przypadków brzegowych ponad wiersz
tabeli, bo prowadzą je komendy o rozgałęzionym wyniku: „Rozpoznanie układu” (rodzina A),
„Scalanie wariantów” (rodzina F) i „Redakcja poufności” (rodzina H). Każdy przebieg
odwołuje się wyłącznie do komend obszaru `studio` z Załącznika rozdz. 11.

**Rozpoznanie układu (rodzina A, funkcja 5).** Odtworzenie struktury logicznej materiału
wejściowego — kolumn, tabel, nagłówków, stopek i przypisów — jako warstwy nad surowym
tekstem rozpoznania.

```
Ingest/OCR Panel        Wykonawca (OCR)              Studio Editor
────────────────        ───────────────              ─────────────
1. pozycja w kolejce
   studio.ingest.queue.add ──►
2. rozpoznanie z nastawami
   studio.ingest.recognize ──► item · layout · words
                               (bloki z bbox)
                                     │
3. ‹ podgląd warstwy ›  ◄───────────┘
   tekstowej na skanie
4. korekta rozpoznania
   studio.ingest.correction.set
        │
        ▼
5. przyjęcie do edytora
   studio.ingest.item.accept ──────────────────────► dokument z odtworzoną strukturą
```

Przebieg pracy:

1. Materiał wejściowy trafia do kolejki panelu komendą `studio.ingest.queue.add`.
2. `studio.ingest.recognize` przyjmuje `itemId` oraz `settings` (silnik, zestaw języków,
   próg pewności, czyszczenie obrazu, przełączniki rozpoznania układu) i zwraca `item`,
   `layout` (bloki z ramkami ograniczającymi) oraz `words` (słowa z pozycjami).
3. Panel prezentuje warstwę tekstową nałożoną na skan; Operator kontroluje jakość.
4. Rozbieżność rozpoznania Operator poprawia komendą `studio.ingest.correction.set`.
5. `studio.ingest.item.accept` przyjmuje wynik jako dokument roboczy Studio Editora
   i zakłada pierwszą wersję w Session Repository.

Przypadki brzegowe:

- **Brak nastaw rozpoznania.** Puste `settings` nie wstrzymują rozpoznania — pole
  przyjmuje wartości domyślne modułu (zasada zero blokad), rozpoznanie rusza z językami
  `pol + eng` i silnikiem lokalnym.
- **Niska pewność rozpoznania układu.** Blok poniżej progu pewności trafia do warstwy
  tekstowej oznaczony jako wymagający korekty; nie jest odrzucany, lecz przekazany
  Operatorowi do poprawy komendą `studio.ingest.correction.set` przed przyjęciem.
- **Materiał bez warstwy graficznej (dokument tekstowy).** Pozycja pomija czyszczenie
  obrazu i rozpoznanie tekstu, `layout` odtwarzany jest wprost ze struktury pliku.
- **Przyjęcie wielu pozycji naraz.** `studio.ingest.item.accept` przyjmuje `itemIds`
  jako zestaw; pozycje bez zakończonego rozpoznania nie wchodzą do przyjęcia i pozostają
  w kolejce w stanie oczekiwania.

**Scalanie wariantów (rodzina F, funkcja 5).** Połączenie dwóch gałęzi dokumentu
z rozstrzyganiem konfliktów fragment po fragmencie.

```
Session Repository       Diff/Grep Panel            Studio Editor
──────────────────       ───────────────            ─────────────
1. gałąź od wersji
   studio.branch.create ──►
2. scalenie gałęzi
   studio.branch.merge ──► merged · conflicts · document
                                │
        ┌───────────────────────┴───────────────────────┐
   brak konfliktów                               konflikt per fragment
        │                                              │
        ▼                                              ▼
3a. merged = true                        3b. conflicts wraca w wyniku
    dokument scalony ─────► edytor           Operator rozstrzyga fragment
                                             studio.diff.hunk.apply ──► edytor
```

Przebieg pracy:

1. Wariant powstaje komendą `studio.branch.create` (`documentId`, `fromVersionId`,
   `name`) — punktem startowym jest wskazana wersja.
2. `studio.branch.merge` przyjmuje `sourceBranchId`, `targetBranchId` oraz `resolutions`
   i zwraca `merged`, `conflicts` i `document`.
3. Fragment bez konfliktu wchodzi do dokumentu scalonego wprost; fragment sporny wraca
   w polu `conflicts` i Operator przenosi wybrane brzmienie komendą
   `studio.diff.hunk.apply` (`sourceVersionId`, `hunkIndex`, `rangeStart`, `rangeEnd`).

Przypadki brzegowe:

- **Konflikt nierozstrzygnięty.** Konflikt bez wpisu w `resolutions` nie jest
  rozstrzygany domyślem — wraca w polu `conflicts` z `merged = false`, a rozstrzygnięcie
  należy do Operatora. Zasada zero blokad nie znosi tu decyzji: brak decyzji zostawia
  fragment sporny, nie nadpisuje żadnej z wersji.
- **Gałąź już scalona.** Ponowne wywołanie `studio.branch.merge` na scalonej parze
  zwraca `merged = true` bez zmian, `conflicts` pusty — czynność jest idempotentna.
- **Scalenie odłożonych brzmień.** Brzmienie fragmentu odłożone przy poprzednim
  rozstrzygnięciu nie przepada — pozostaje dostępne w historii wersji obu gałęzi
  (Session Repository) i można je przenieść komendą `studio.diff.hunk.apply`.

**Redakcja poufności (rodzina H, funkcja 3).** Trwałe zamazanie wskazanych fragmentów
wraz z warstwą tekstową pod spodem. Czynność jest nieodwracalna dla wyniku i stanowi
zapisany wyjątek od zasady zero blokad — wynik powstaje z potwierdzeniem, źródło
zostaje nietknięte.

```
Studio Editor / Preview     Wykonawca                    wynik
───────────────────────     ─────────                    ─────
1. wykrycie danych wrażliwych
   studio.security.sensitive.detect ──► findings (CZYTA, nic nie zmienia)
        │
2. wskazanie regionów + potwierdzenie
   studio.security.redact ──────────► asset (kopia) · redacted
        │                             źródło nietknięte
        ▼
3. usunięcie metadanych przed wysyłką
   studio.security.metadata.strip ──► asset · removed
```

Przebieg pracy:

1. `studio.security.sensitive.detect` (`documentId`, `versionId`, `categories`) zwraca
   `findings` — czynność czyta i niczego nie zmienia.
2. Operator potwierdza zakres i uruchamia `studio.security.redact` (`assetId`,
   `regions`); wynik to nowy `asset` z policzoną liczbą `redacted`, źródło pozostaje
   nietknięte.
3. Przed przekazaniem dokumentu `studio.security.metadata.strip` czyści autora,
   historię i dane ukryte (`removed`).

Przypadki brzegowe:

- **Brak potwierdzenia.** Redakcja poufności nie rusza bez potwierdzenia — trwałe
  zamazanie jest zapisanym wyjątkiem od zero blokad (czynność nieodwracalna dla wyniku).
- **Region poza treścią.** Region wskazany poza zakresem dokumentu nie zamazuje niczego
  i nie zatrzymuje pozostałych regionów; `redacted` liczy wyłącznie regiony skuteczne.
- **Dane wrażliwe pominięte przy wykrywaniu.** `studio.security.sensitive.detect`
  wskazuje kategorie do redakcji, lecz nie warunkuje jej — Operator może wskazać region
  ręcznie bez uprzedniego wykrycia. Brak trafień w `findings` nie blokuje `redact`.
- **Zachowanie wybranych metadanych.** Pole `keepFields` w `studio.security.metadata.strip`
  pozostawia wskazane pola nietknięte; puste `keepFields` czyści komplet.

---

## 8. Punkty sterowania z okna konfiguracji

Wszystkie pozycje działają wg zasady „brak ustawienia = wartość domyślna”; żadna nie blokuje pracy. Warstwowość ustawień: globalna → środowisko → projekt → sesja.

| Grupa ustawień | Co Operator personalizuje | Wartość domyślna | Warstwa |
|---|---|---|---|
| Silnik rozpoznawania tekstu | Wybór silnika (lokalny Tesseract / usługa chmurowa), języki (`tessdata`), próg pewności, czyszczenie obrazu przed rozpoznaniem | Tesseract lokalny, języki pol+eng, czyszczenie automatyczne | globalna/sesja |
| Kanały modeli operacji AI | Który `kanal_modelu` obsługuje którą rodzinę operacji (korekta tańszym modelem, redakcja mocniejszym) | model domyślny karty | projekt/sesja |
| Szablony operacji (prompty) | Treść promptów operacji „Korekta ortograficzno-gramatyczna”–„Weryfikacja i wskazanie źródeł”, dodawanie operacji własnych „Operacje własne” | prompty fabryczne | globalna/projekt |
| Łańcuchy i makra | Definicje łańcuchów operacji („Makra i łańcuchy operacji”), operacje wsadowe („Operacje wsadowe na wielu dokumentach”) | brak definicji — tworzy je użytkownik | projekt |
| Formaty i konwersje | Domyślny format eksportu, ścieżka renderu PDF (fpdf/gopdf albo chromedp/Pandoc), profile eksportu („Eksport z profilami”) | PDF przez render docelowy | globalna/sesja |
| Autozapis i wersje | Interwał autozapisu („Autozapis i migawki”), migawka po operacji AI, polityka retencji historii | autozapis 30 s, migawka włączona | sesja |
| Diff | Tryb domyślny (scalony/dwukolumnowy), granularność (znak/słowo/zdanie), kolory oznaczeń | scalony, słowo | sesja |
| Redakcja i poufność | Reguły wykrywania danych wrażliwych („Wykrywanie danych wrażliwych”), automatyczne usuwanie metadanych przy eksporcie („Usuwanie metadanych”) | wykrywanie włączone, automatyczne czyszczenie wyłączone | projekt |
| Skróty klawiszowe | Pełna mapa skrótów („Skróty klawiszowe konfigurowalne”), paleta poleceń („Paleta poleceń”) | mapa domyślna | globalna |
| Widok edytora | Domyślny tryb pracy (WYSIWYG/Markdown), tryb skupienia („Tryb skupienia”), motyw typografii | WYSIWYG | sesja |
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

---

## 10. Punkty łamania i kryteria odbioru

### 10.1 Punkty łamania

Moduł Studio dzieli obszar roboczy na kolumny sąsiadujące poziomo (rozdz. 2). Poniższe progi, wspólne całej platformie, rozstrzygają zachowanie układu tych kolumn wraz ze zmianą szerokości okna aplikacji.

| Żeton | Próg szerokości | Zachowanie układu w module Studio |
|---|---|---|
| `--dn-bp-w1` | 640 px | Telefon poziomo — widok mobilny; okna operacyjne modułu prezentowane pojedynczo, pełny ekran, nawigacja powrotna zastępuje układ kolumnowy |
| `--dn-bp-w2` | 960 px | Tablet — boczna nawigacja modułów zwija się do samych ikon; kolumny modułu Studio zachowują układ, lecz kolumny boczne (rozszerzenia warstwy 2–3) otwierają się jako nakładka zamiast stałej kolumny |
| `--dn-bp-w3` | 1280 px | Biurko — pełny kokpit; wszystkie okna operacyjne modułu Studio wymienione w rozdz. 2 dostępne jednocześnie w układzie kolumnowym opisanym w tym rozdziale |
| `--dn-bp-w4` | 1600 px | Szerokie biurko — para Chat Window i Execution Loop Window prezentowana jednocześnie obok okna wiodącego modułu, bez wzajemnego przesłaniania |

Poniżej progu `--dn-bp-w1` układ kolumnowy modułu Studio nie jest dostępny — zachowanie na urządzeniach mobilnych ustala [Mobile](../funkcje-globalne/mobile.md).

### 10.2 Kryteria odbioru

Warunki sprawdzalne, których łączne spełnienie oznacza gotowość modułu Studio do odbioru.

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Wszystkie 8 okien operacyjnych z rozdz. 2 otwiera się dokładnie sposobem wywołania opisanym w tabeli okien i w katalogu elementów interfejsu rozdz. 3 | Przegląd manualny wg tabeli rozdz. 2 — każde okno wywołane, sprawdzone wejście i zamknięcie |
| Każda z 179 komend obszaru `studio` z Załącznika ma pokrycie w co najmniej jednym przebiegu pracy albo scenariuszu użycia (rozdz. 4, 6) | Zestawienie nazw komend z treścią rozdz. 4 i 6 |
| Cykl wczytanie → redakcja wspierana AI → porównanie różnicowe (Diff/Grep Panel) → wersjonowanie (Session Repository) → eksport (Preview Window) działa bez przerwy dla materiału każdego obsługiwanego formatu (rozdz. 1.4) | Przebieg manualny jednego pełnego cyklu redakcyjnego, wg scenariuszy rozdz. 6 |
| Żaden opisany element interfejsu nie odwołuje się do klasy `.dn-*` ani żetonu `--dn-*` nieobecnego w arkuszach `design/zasoby/` | `grep` nazwy klas i żetonów użytych w dokumencie względem arkuszy źródłowych |
| Ingest/OCR Panel otwiera się samoczynnie, gdy wejściem sesji jest skan, zdjęcie lub zestaw plików, bez odrębnego polecenia Operatora (rozdz. 2) | Przegląd sposobu wywołania Ingest/OCR Panel w tabeli okien rozdz. 2 |


---

## 11. Załącznik — pełny wykaz komend kontraktu modułu Studio

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### 11.1 Obszar `studio` — 179 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `studio.agents.claim` | Zajmuje fragment dokumentu dla wykonawcy, żeby dwóch agentów nie pisało po tym samym akapicie. Fragment zajęty przez kogo innego wraca odmowa nazywająca wykonawcę i zakres | `documentId:string` (wym)<br>`rangeStart:int` (wym)<br>`rangeEnd:int` (wym)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc)<br>`taskId:string` (opc)<br>`ttlSeconds:int` (opc) | `slot:StudioAgentSlot` (opc)<br>`claimed:bool` (wym)<br>`heldBy:StudioAgentSlot` (opc)<br>`refusalReason:string` (opc) |
| `studio.agents.conflicts.list` | Oddaje bilans spięć wykonawców o ten sam fragment — czyja zmiana weszła, czyja została odłożona i dlaczego. Odłożone brzmienie nie przepada | `documentId:string` (wym)<br>`resolved:bool` (opc)<br>`agentId:string` (opc)<br>`limit:int` (opc) | `conflicts:StudioAgentConflict[]` (wym)<br>`openCount:int` (wym) |
| `studio.agents.release` | Zwalnia fragment zajęty przez wykonawcę | `documentId:string` (wym)<br>`slotId:string` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc)<br>`all:bool` (opc) | `released:int` (wym)<br>`slots:StudioAgentSlot[]` (wym) |
| `studio.agents.settings.get` | Oddaje nastawy pętli wykonawczej i pracy wielu agentów wraz z zasięgiem, z którego przyszły. Oba narzędzia są domyślnie wyłączone | `documentId:string` (opc)<br>`windowId:string` (opc)<br>`sessionId:string` (opc) | `settings:StudioAgentSettings` (wym) |
| `studio.agents.settings.set` | Włącza albo wyłącza pętlę wykonawczą i pracę wielu agentów nad dokumentem — jawnym, odwracalnym ustawieniem Operatora. Wartości idą zasięgami rodziny config, nie osobnym magazynem | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc)<br>`executionLoopEnabled:bool` (opc)<br>`multiAgentEnabled:bool` (opc)<br>`maxConcurrentAgents:int` (opc)<br>`conflictPolicy:StudioAgentConflictPolicy` (opc)<br>`loopMaxIterations:int` (opc)<br>`loopNoProgressThreshold:int` (opc)<br>`requireFragmentClaim:bool` (opc) | `settings:StudioAgentSettings` (wym) |
| `studio.agents.slots.list` | Oddaje wykaz wykonawców pracujących nad dokumentem wraz z fragmentami, które zajęli — Operator widzi, kto nad czym pracuje, i model widzi, czego nie tykać | `documentId:string` (wym)<br>`state:StudioAgentSlotState` (opc)<br>`includeExpired:bool` (opc) | `slots:StudioAgentSlot[]` (wym)<br>`activeAgents:int` (wym) |
| `studio.annotation.add` | Zakłada adnotację przy fragmencie różnicy. Dziś czynność ta idzie drogą generyczną window.action i wraca odmowa not_found, bo katalog akcji nie ma jej wiersza | `documentId:string` (wym)<br>`hunkIndex:int` (wym)<br>`baseVersionId:string` (opc)<br>`targetVersionId:string` (opc)<br>`proposalId:string` (opc)<br>`body:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `annotation:StudioAnnotation` (wym) |
| `studio.annotation.list` | Zwraca adnotacje założone przy różnicach dokumentu | `documentId:string` (wym)<br>`baseVersionId:string` (opc)<br>`targetVersionId:string` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `annotations:StudioAnnotation[]` (wym) |
| `studio.apparatus.insert` | Zakłada element aparatu dokumentu — spis treści, spis ilustracji albo tabel, przypis dolny i końcowy, podpis, zakładkę, odwołanie wzajemne, odsyłacz, powołanie, bibliografię albo hasło indeksu | `documentId:string` (wym)<br>`kind:StudioApparatusKind` (wym)<br>`offset:int` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`text:string` (opc)<br>`label:string` (opc)<br>`targetId:string` (opc)<br>`targetUrl:string` (opc)<br>`levelsFrom:int` (opc)<br>`levelsTo:int` (opc)<br>`citationKey:string` (opc)<br>`sourceTitle:string` (opc)<br>`sourceAuthor:string` (opc)<br>`sourceYear:string` (opc)<br>`sourceUrl:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`item:StudioApparatusItem` (wym) |
| `studio.apparatus.list` | Oddaje aparat dokumentu — spisy, przypisy, podpisy, zakładki, odwołania, powołania i hasła indeksu | `documentId:string` (wym)<br>`kind:StudioApparatusKind` (opc)<br>`staleOnly:bool` (opc) | `items:StudioApparatusItem[]` (wym) |
| `studio.apparatus.refresh` | Odświeża aparat dokumentu — spis treści zgadza się z nagłówkami, przypisy są przenumerowane, spisy ilustracji i tabel policzone od nowa | `documentId:string` (wym)<br>`kind:StudioApparatusKind` (opc)<br>`itemId:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`items:StudioApparatusItem[]` (wym) |
| `studio.apparatus.remove` | Usuwa element aparatu dokumentu; przypisy pozostałe przenumerowują się | `documentId:string` (wym)<br>`itemId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `removed:bool` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`actionId:string` (opc) |
| `studio.asset.embed` | Osadza zasób wizualny z magazynu rdzenia w treści dokumentu we wskazanym miejscu. Realizuje powiązanie Design do Studio | `documentId:string` (wym)<br>`assetId:string` (wym)<br>`position:int` (opc)<br>`caption:string` (opc)<br>`altText:string` (opc) | `document:StudioDocument` (wym) |
| `studio.autosave.get` | Oddaje nastawy autozapisu wraz z czasem ostatniego udanego zapisu i uczciwym wskazaniem, czy ostatni zapis samoczynny się nie udał | `documentId:string` (opc)<br>`windowId:string` (opc) | `settings:StudioAutosaveSettings` (wym) |
| `studio.autosave.run` | Wykonuje zapis samoczynny wraz z postacią i odkłada wersję w osobnym szeregu autozapisu, żeby nie zaśmiecać historii wersji Operatora. Nieudany zapis wraca nazwany, nie przemilczany | `documentId:string` (wym)<br>`content:string` (opc)<br>`form:json` (opc)<br>`trigger:StudioBackupReason` (opc) | `saved:bool` (wym)<br>`failureReason:string` (opc)<br>`version:StudioVersion` (opc)<br>`backup:StudioDocumentBackup` (opc)<br>`savedAt:int64` (opc) |
| `studio.autosave.set` | Ustawia autozapis — odstęp, zapis przy zdarzeniach okna oraz zasadę wygasania kopii zapasowych jako jawne, odwracalne ustawienie Operatora | `documentId:string` (opc)<br>`windowId:string` (opc)<br>`enabled:bool` (wym)<br>`intervalSeconds:int` (opc)<br>`onBlur:bool` (opc)<br>`onClose:bool` (opc)<br>`onSwitch:bool` (opc)<br>`backupRetentionCount:int` (opc)<br>`backupRetentionHours:int` (opc) | `settings:StudioAutosaveSettings` (wym) |
| `studio.backup.create` | Zakłada kopię zapasową dokumentu niezależnie od historii wersji, żeby przetrwała awarię procesu i awarię zapisu | `documentId:string` (wym)<br>`reason:StudioBackupReason` (opc)<br>`content:string` (opc)<br>`form:json` (opc) | `backup:StudioDocumentBackup` (wym) |
| `studio.backup.list` | Oddaje wykaz kopii zapasowych wraz z czasem i rozmiarem oraz wskazaniem kopii niosącej zmiany niezapisane — po niej Studio samo zgłasza przywrócenie po nagłym zamknięciu | `documentId:string` (opc)<br>`windowId:string` (opc)<br>`unsavedOnly:bool` (opc) | `backups:StudioDocumentBackup[]` (wym)<br>`unsavedCount:int` (wym) |
| `studio.backup.restore` | Przywraca kopię zapasową — na miejsce albo do nowego dokumentu, żeby przywracanie samo nie kasowało tego, co jest | `backupId:string` (wym)<br>`asNewDocument:bool` (opc)<br>`windowId:string` (opc)<br>`title:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym) |
| `studio.batch.run` | Uruchamia tę samą operację kontekstową na wielu dokumentach naraz. Każdy dokument dostaje własne zadanie pętli; odmowa jednego nie wstrzymuje pozostałych | `windowId:string` (wym)<br>`documentIds:string[]` (wym)<br>`actionId:string` (wym)<br>`params:json` (opc) | `runId:string` (wym)<br>`accepted:int` (wym)<br>`rejected:StudioBatchRejection[]` (opc) |
| `studio.branch.create` | Zakłada gałąź dokumentu od wskazanej wersji jako punktu startowego. Służy prowadzeniu dwóch redakcji obok siebie | `documentId:string` (wym)<br>`fromVersionId:string` (wym)<br>`name:string` (wym) | `branch:StudioBranch` (wym)<br>`document:StudioDocument` (wym) |
| `studio.branch.list` | Zwraca gałęzie dokumentu wraz z ich punktami startowymi | `documentId:string` (wym) | `branches:StudioBranch[]` (wym) |
| `studio.branch.merge` | Scala gałąź dokumentu z gałęzią docelową, rozstrzygając konflikty per fragment. Konflikt nierozstrzygnięty wraca w wyniku zamiast być rozstrzygnięty domysłem | `sourceBranchId:string` (wym)<br>`targetBranchId:string` (wym)<br>`resolutions:StudioMergeResolution[]` (opc) | `merged:bool` (wym)<br>`conflicts:StudioMergeConflict[]` (opc)<br>`document:StudioDocument` (opc) |
| `studio.chain.list` | Zwraca łańcuchy operacji dostępne Operatorowi | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `chains:StudioChain[]` (wym) |
| `studio.chain.run` | Uruchamia łańcuch operacji na dokumencie. Przebieg prowadzi pętla wykonawcza okna, a każdy krok odkłada własną propozycję zmiany | `windowId:string` (wym)<br>`documentId:string` (wym)<br>`chainId:string` (wym)<br>`scope:StudioOperationScope` (wym)<br>`selectionStart:int` (opc)<br>`selectionEnd:int` (opc) | `runId:string` (wym)<br>`steps:int` (wym) |
| `studio.chain.save` | Zapisuje łańcuch operacji — sekwencję wykonywaną jednym poleceniem: korekta, streszczenie, wydanie | `chainId:string` (opc)<br>`name:string` (wym)<br>`steps:StudioChainStep[]` (wym)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `chain:StudioChain` (wym) |
| `studio.clipboard.copy` | Odkłada fragment dokumentu do schowka platformy wraz z jego postacią — dostępne modelowi tak samo jak Operatorowi | `documentId:string` (wym)<br>`rangeStart:int` (wym)<br>`rangeEnd:int` (wym)<br>`cut:bool` (opc)<br>`formatOnly:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `clipboardEntryId:string` (wym)<br>`text:string` (wym)<br>`form:StudioDocumentForm` (opc)<br>`actionId:string` (opc) |
| `studio.clipboard.paste` | Wkleja wpis schowka w miejsce kursora — z zachowaniem postaci albo jako czysty tekst, według jawnego wyboru Operatora | `documentId:string` (wym)<br>`offset:int` (wym)<br>`clipboardEntryId:string` (opc)<br>`text:string` (opc)<br>`mode:StudioPasteMode` (opc)<br>`replaceRangeStart:int` (opc)<br>`replaceRangeEnd:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.comment.add` | Zakłada komentarz redakcyjny przypięty do fragmentu dokumentu albo odpowiedź w wątku komentarza | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`selectionStart:int` (opc)<br>`selectionEnd:int` (opc)<br>`parentCommentId:string` (opc)<br>`body:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `comment:StudioComment` (wym) |
| `studio.comment.list` | Zwraca komentarze redakcyjne dokumentu wraz z wątkami odpowiedzi | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`includeResolved:bool` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `comments:StudioComment[]` (wym) |
| `studio.comment.resolve` | Oznacza wątek komentarza jako rozwiązany albo cofa to oznaczenie. Wątek nie znika: rozwiązanie jest stanem, nie usunięciem | `commentId:string` (wym)<br>`resolved:bool` (wym) | `comment:StudioComment` (wym) |
| `studio.contextual.op` | Wykonuje operację kontekstową Tools Panel na zaznaczeniu albo całym dokumencie | `windowId:string` (wym)<br>`documentId:string` (wym)<br>`actionId:string` (wym)<br>`scope:StudioOperationScope` (wym)<br>`selectionStart:int` (opc)<br>`selectionEnd:int` (opc)<br>`params:json` (opc) | `messageId:string` (wym)<br>`resultText:string` (opc)<br>`proposalId:string` (opc) |
| `studio.diff.compare` | Porównuje wersje dokumentu i wyszukuje wzorzec w treści | `documentId:string` (wym)<br>`baseVersionId:string` (opc)<br>`targetVersionId:string` (opc)<br>`proposalId:string` (opc)<br>`pattern:string` (opc)<br>`regex:bool` (opc) | `hunks:StudioDiffHunk[]` (opc)<br>`matches:StudioTextMatch[]` (opc) |
| `studio.diff.form.compare` | Porównuje postać dwóch dowolnych wersji dokumentu — zmiana kroju czy wcięcia jest widoczna jako zmiana, a nie milczy, bo litery zostały te same | `documentId:string` (wym)<br>`baseVersionId:string` (opc)<br>`targetVersionId:string` (opc)<br>`area:string` (opc) | `entries:StudioFormDiffEntry[]` (wym)<br>`added:int` (wym)<br>`removed:int` (wym)<br>`changed:int` (wym) |
| `studio.diff.hunk.apply` | Przenosi pojedynczy fragment ze wskazanej wersji do stanu bieżącego — praktyczny sens widoku różnicy, nie samo patrzenie | `documentId:string` (wym)<br>`sourceVersionId:string` (wym)<br>`hunkIndex:int` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`includeForm:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`document:StudioDocument` (wym) |
| `studio.diff.report.export` | Wydaje widok różnicowy jako osobny dokument redakcji: fragmenty, statystykę i adnotacje. Wynik jest zasobem magazynu rdzenia | `documentId:string` (wym)<br>`baseVersionId:string` (wym)<br>`targetVersionId:string` (opc)<br>`format:string` (wym)<br>`includeAnnotations:bool` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`sizeBytes:int` (wym) |
| `studio.diff.source` | Zestawia dokument roboczy z materiałem wejściowym z Library albo Research i wskazuje rozbieżności. Służy kontroli, czy redakcja nie odeszła od źródła | `documentId:string` (wym)<br>`sourceLibraryFileId:string` (opc)<br>`sourceAssetId:string` (opc)<br>`versionId:string` (opc) | `hunks:StudioDiffHunk[]` (wym)<br>`sourceResolved:bool` (wym) |
| `studio.diff.visual` | Porównuje wygląd dwóch wersji dokumentu: renderuje obie strony i wskazuje obszary różniące się graficznie. Służy tam, gdzie różnica tekstowa nie widzi zmiany układu | `documentId:string` (wym)<br>`baseVersionId:string` (wym)<br>`targetVersionId:string` (wym)<br>`pageFrom:int` (opc)<br>`pageTo:int` (opc)<br>`windowId:string` (opc) | `regions:StudioVisualDiffRegion[]` (wym)<br>`overlayAssetIds:string[]` (opc) |
| `studio.document.copy` | Zakłada kopię dokumentu wraz z całą postacią oraz — według jawnego wyboru Operatora — z historią wersji albo bez niej. Kopia jest osobnym dokumentem, nie drugim odwołaniem do tego samego | `documentId:string` (wym)<br>`title:string` (opc)<br>`windowId:string` (opc)<br>`includeVersions:bool` (opc)<br>`includeMarkup:bool` (opc)<br>`includeLocks:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`versionsCopied:int` (wym) |
| `studio.document.create` | Zakłada dokument pusty — nową stronę gotową do pisania — z arkuszem stylów i nastawami strony domyślnymi albo przejętymi z szablonu | `windowId:string` (wym)<br>`title:string` (opc)<br>`templateId:string` (opc)<br>`paperName:string` (opc)<br>`orientation:StudioPageOrientation` (opc)<br>`format:StudioDocumentFormat` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym) |
| `studio.document.export.batch` | Wydaje wiele dokumentów naraz do wskazanego formatu; odmowa jednego nie wstrzymuje pozostałych, a każdy wynik niesie swój bilans | `documentIds:string[]` (wym)<br>`format:StudioExportFormat` (wym)<br>`directory:string` (opc)<br>`profileId:string` (opc) | `results:StudioExportResult[]` (wym)<br>`succeeded:int` (wym)<br>`failed:int` (wym)<br>`failures:StudioSkippedItem[]` (opc) |
| `studio.document.export.format` | Wydaje dokument do formatu docelowego — txt, md, docx, odt, pdf, html albo rtf — wkompilowanym rachunkiem, bez programu zewnętrznego. Format uboższy niż dokument oddaje wykaz cech pominiętych | `documentId:string` (wym)<br>`format:StudioExportFormat` (wym)<br>`path:string` (opc)<br>`profileId:string` (opc)<br>`includeComments:bool` (opc)<br>`includeTrackedChanges:bool` (opc)<br>`templateId:string` (opc) | `result:StudioExportResult` (wym) |
| `studio.document.form.get` | Oddaje pełną postać dokumentu — arkusz stylów, nastawy strony, sekcje, bloki, tabele, obiekty, aparat i pola | `documentId:string` (wym)<br>`includeBlocks:bool` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `form:StudioDocumentForm` (wym) |
| `studio.document.form.save` | Zapisuje postać dokumentu wraz z treścią — arkusz stylów, sekcje, tabele, obiekty i aparat przestają ginąć przy zapisie | `documentId:string` (wym)<br>`form:json` (wym)<br>`content:string` (opc)<br>`title:string` (opc)<br>`createVersion:bool` (opc)<br>`series:StudioVersionSeries` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`version:StudioVersion` (opc)<br>`balance:StudioActionBalance` (wym) |
| `studio.document.format.set` | Przestawia format dokumentu Studia. Dziś format nadaje rdzeń przy wczytaniu i żadna komenda go nie zmienia: studio.document.save przyjmuje documentId, content, title i createVersion, ale nie format | `documentId:string` (wym)<br>`format:StudioDocumentFormat` (wym)<br>`convertContent:bool` (opc) | `document:StudioDocument` (wym) |
| `studio.document.image.import` | Wnosi obraz wprost w miejsce kursora — z pliku, z magazynu zasobów rdzenia, z modułu Design albo z bazy zdjęciowej | `documentId:string` (wym)<br>`offset:int` (wym)<br>`source:StudioObjectSource` (wym)<br>`path:string` (opc)<br>`assetId:string` (opc)<br>`designNodeId:string` (opc)<br>`photoBankId:string` (opc)<br>`libraryFileId:string` (opc)<br>`bytesBase64:string` (opc)<br>`widthMm:float` (opc)<br>`heightMm:float` (opc)<br>`altText:string` (opc)<br>`caption:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`object:StudioDocumentObject` (wym)<br>`provenance:StudioProvenance` (opc) |
| `studio.document.import.file` | Wnosi plik Operatora wprost do edytora z zachowaną postacią — docx, dotx, odt, ott, tekst czysty, markdown, RTF i HTML; zapis znaków rozpoznawany, także strony kodowe inne niż UTF-8 | `windowId:string` (wym)<br>`documentId:string` (opc)<br>`format:StudioImportFormat` (opc)<br>`path:string` (opc)<br>`libraryFileId:string` (opc)<br>`assetId:string` (opc)<br>`bytesBase64:string` (opc)<br>`encoding:string` (opc)<br>`insertAtOffset:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioImportBalance` (wym)<br>`provenance:StudioProvenance` (opc) |
| `studio.document.import.pdf` | Zamienia PDF na dokument edytowalny wprost w edytorze — tekst, akapity, tabele i obrazy odzyskane na tyle, na ile PDF je niesie. Odzyskanie jest odtworzeniem, nie odczytem, więc odpowiedź niesie bilans; PDF ze samych skanów kieruje na rozpoznanie tekstu | `windowId:string` (wym)<br>`documentId:string` (opc)<br>`path:string` (opc)<br>`assetId:string` (opc)<br>`libraryFileId:string` (opc)<br>`bytesBase64:string` (opc)<br>`pages:string` (opc)<br>`recoverTables:bool` (opc)<br>`recoverImages:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioImportBalance` (wym)<br>`ingestItemId:string` (opc) |
| `studio.document.open` | Wczytuje dokument do Studio Editor z repozytorium albo z urządzenia | `windowId:string` (wym)<br>`documentId:string` (opc)<br>`libraryFileId:string` (opc)<br>`path:string` (opc) | `document:StudioDocument` (wym) |
| `studio.document.save` | Zapisuje treść dokumentu i zakłada wersję w repozytorium sesji | `documentId:string` (wym)<br>`content:string` (wym)<br>`title:string` (opc)<br>`createVersion:bool` (opc)<br>`form:json` (opc)<br>`author:StudioAuthor` (opc)<br>`series:StudioVersionSeries` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`version:StudioVersion` (opc) |
| `studio.document.save.as` | Zapisuje dokument pod nową nazwą albo do wskazanego pliku, wraz z całą postacią | `documentId:string` (wym)<br>`title:string` (opc)<br>`path:string` (opc)<br>`format:StudioExportFormat` (opc)<br>`createVersion:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`export:StudioExportResult` (opc) |
| `studio.export.profile.list` | Zwraca profile wydania dostępne Operatorowi | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `profiles:StudioExportProfile[]` (wym) |
| `studio.export.profile.save` | Zapisuje nazwany profil wydania: format, marginesy, orientację, nagłówek, stopkę, numerację i skalę | `profileId:string` (opc)<br>`name:string` (wym)<br>`format:string` (wym)<br>`pageSetup:StudioPageSetup` (opc)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `profile:StudioExportProfile` (wym) |
| `studio.field.insert` | Wstawia pole dokumentu — numer strony, liczbę stron, datę, właściwość dokumentu albo pole obliczane | `documentId:string` (wym)<br>`offset:int` (wym)<br>`kind:StudioFieldKind` (wym)<br>`format:string` (opc)<br>`expression:string` (opc)<br>`propertyName:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`field:StudioDocumentField` (wym) |
| `studio.field.list` | Oddaje pola dokumentu wraz z wartościami i wskazaniem, które wymagają odświeżenia | `documentId:string` (wym)<br>`kind:StudioFieldKind` (opc) | `fields:StudioDocumentField[]` (wym) |
| `studio.field.refresh` | Odświeża pola dokumentu — wartości policzone od nowa | `documentId:string` (wym)<br>`fieldId:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`fields:StudioDocumentField[]` (wym) |
| `studio.format.case.set` | Przestawia wielkość liter wskazanego fragmentu | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`transform:StudioCaseTransform` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.format.character.get` | Oddaje styl znaku obowiązujący na wskazanym fragmencie wraz z tym, co we fragmencie niejednolite | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `character:StudioCharacterFormat` (wym)<br>`mixedFields:string[]` (opc)<br>`runs:StudioDocumentRun[]` (opc) |
| `studio.format.character.set` | Ustawia styl znaku na wskazanym fragmencie — krój, stopień, grubość, odmianę, podkreślenie, przekreślenie, indeksy, barwę, wyróżnienie, odstęp międzyliterowy, kapitaliki | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`fontFamily:string` (opc)<br>`fontSizePt:float` (opc)<br>`fontSizeStepPt:float` (opc)<br>`bold:bool` (opc)<br>`italic:bool` (opc)<br>`underline:StudioUnderlineStyle` (opc)<br>`strikethrough:bool` (opc)<br>`superscript:bool` (opc)<br>`subscript:bool` (opc)<br>`color:string` (opc)<br>`highlightColor:string` (opc)<br>`letterSpacingPt:float` (opc)<br>`smallCaps:bool` (opc)<br>`allCaps:bool` (opc)<br>`effect:StudioTextEffect` (opc)<br>`language:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.format.clear` | Czyści formatowanie wskazanego fragmentu — postać wraca do stylu nazwanego albo do postaci domyślnej | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`character:bool` (opc)<br>`paragraph:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.format.painter.apply` | Malarz formatów — nanosi pobraną postać na wskazany fragment | `documentId:string` (wym)<br>`clipId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`includeParagraph:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.format.painter.copy` | Malarz formatów — pobiera postać wskazanego fragmentu do naniesienia w innym miejscu; kopiuje postać, nie treść | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`includeParagraph:bool` (opc) | `clipId:string` (wym)<br>`character:StudioCharacterFormat` (opc)<br>`paragraph:StudioParagraphFormat` (opc) |
| `studio.format.paragraph.get` | Oddaje styl akapitu obowiązujący na wskazanym fragmencie | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `paragraph:StudioParagraphFormat` (wym)<br>`mixedFields:string[]` (opc) |
| `studio.format.paragraph.set` | Ustawia styl akapitu na wskazanym fragmencie — wyrównanie, wcięcia, odstępy, interlinię, tabulatory, obramowanie, cieniowanie, kontrolę wdów i sierot, razem z następnym | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`align:StudioTextAlign` (opc)<br>`firstLineIndentMm:float` (opc)<br>`indentLeftMm:float` (opc)<br>`indentRightMm:float` (opc)<br>`indentStepMm:float` (opc)<br>`spaceBeforePt:float` (opc)<br>`spaceAfterPt:float` (opc)<br>`lineSpacingRule:StudioLineSpacingRule` (opc)<br>`lineSpacingValue:float` (opc)<br>`tabStops:json` (opc)<br>`border:json` (opc)<br>`shadingColor:string` (opc)<br>`widowControl:bool` (opc)<br>`keepWithNext:bool` (opc)<br>`keepLines:bool` (opc)<br>`outlineLevel:int` (opc)<br>`rightToLeft:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.format.replace` | Znajdź i zamień wraz z postacią — szuka treści, postaci albo obojga i zamienia treść, postać albo oboje; fragmenty pod blokadą wracają bilansem | `documentId:string` (wym)<br>`findText:string` (opc)<br>`replaceText:string` (opc)<br>`findFormat:json` (opc)<br>`replaceFormat:json` (opc)<br>`findStyleName:string` (opc)<br>`replaceStyleName:string` (opc)<br>`matchCase:bool` (opc)<br>`wholeWord:bool` (opc)<br>`regex:bool` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`replaceAll:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`matches:int` (wym)<br>`replaced:int` (wym) |
| `studio.format.similar.select` | Zaznacz według podobnego formatowania — oddaje wszystkie fragmenty dokumentu o postaci zgodnej ze wskazanym wzorem | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`styleName:string` (opc)<br>`matchCharacter:bool` (opc)<br>`matchParagraph:bool` (opc) | `matches:StudioDocumentRun[]` (wym)<br>`count:int` (wym) |
| `studio.ingest.correction.set` | Poprawia rozpoznane słowo na warstwie tekstowej pozycji kolejki, przed przyjęciem wyniku do edytora | `itemId:string` (wym)<br>`wordIndex:int` (wym)<br>`text:string` (wym) | `item:StudioIngestItem` (wym) |
| `studio.ingest.device.list` | Zwraca urządzenia wejściowe widziane przez rdzeń: skanery i kamery. Bez tego wykazu wybór urządzenia nie ma z czego powstać | — | `devices:StudioInputDevice[]` (wym) |
| `studio.ingest.device.scan` | Pobiera obraz ze skanera albo kamery podłączonej do maszyny rdzenia i dokłada go jako pozycję kolejki | `windowId:string` (wym)<br>`deviceId:string` (opc)<br>`resolutionDpi:int` (opc)<br>`colorMode:string` (opc)<br>`pages:int` (opc) | `items:StudioIngestItem[]` (wym) |
| `studio.ingest.item.accept` | Przyjmuje wynik pozycji kolejki jako dokument roboczy Studio Editora i zakłada pierwszą wersję w repozytorium sesji | `windowId:string` (wym)<br>`itemIds:string[]` (wym)<br>`title:string` (opc)<br>`format:StudioDocumentFormat` (opc) | `document:StudioDocument` (wym)<br>`version:StudioVersion` (wym) |
| `studio.ingest.queue.add` | Dokłada materiał do kolejki wczytywania po stronie rdzenia: pojedynczy plik, katalog albo archiwum. Dziś kolejka jest wyłącznie kliencka i ginie z odświeżeniem okna | `windowId:string` (wym)<br>`sourcePaths:string[]` (opc)<br>`assetIds:string[]` (opc)<br>`archivePath:string` (opc)<br>`settings:StudioRecognitionSettings` (opc) | `items:StudioIngestItem[]` (wym) |
| `studio.ingest.queue.list` | Zwraca kolejkę wczytywania okna wraz ze stanem i wynikiem każdej pozycji | `windowId:string` (wym)<br>`pendingOnly:bool` (opc) | `items:StudioIngestItem[]` (wym) |
| `studio.ingest.recognize` | Rozpoznaje tekst pozycji kolejki z pełnym sterowaniem: silnik, zestaw języków, próg pewności, czyszczenie obrazu i rozpoznanie układu. Komenda document.text.extract robi to samo wąskim wejściem — jeden język, bez silnika, bez progu i bez układu | `itemId:string` (wym)<br>`settings:StudioRecognitionSettings` (opc) | `item:StudioIngestItem` (wym)<br>`layout:StudioLayoutBlock[]` (opc)<br>`words:StudioRecognizedWord[]` (opc) |
| `studio.ingest.url` | Pobiera stronę sieciową, oczyszcza ją z nawigacji i reklam i dokłada jako pozycję kolejki wczytywania. Migawkę strony oddaje browser.snapshot.get, ale wymaga okna modułu Browser i nie prowadzi do dokumentu Studia | `windowId:string` (wym)<br>`url:string` (wym)<br>`readability:bool` (opc)<br>`includeImages:bool` (opc) | `item:StudioIngestItem` (wym) |
| `studio.insert.from.library` | Wnosi plik, wzór, załącznik albo obraz z modułu Library wprost do dokumentu w miejsce kursora — nie do kolejki i nie do zasobów — wraz z zapisem pochodzenia | `documentId:string` (wym)<br>`offset:int` (wym)<br>`libraryFileId:string` (wym)<br>`asObject:bool` (opc)<br>`asAttachment:bool` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`provenance:StudioProvenance` (wym) |
| `studio.insert.from.web` | Wnosi fragment albo obraz ze strony sieci wprost do dokumentu w miejsce kursora, wraz z zapisem pochodzenia — adresem i czasem sięgnięcia | `documentId:string` (wym)<br>`offset:int` (wym)<br>`url:string` (wym)<br>`selector:string` (opc)<br>`text:string` (opc)<br>`imageUrl:string` (opc)<br>`asObject:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`provenance:StudioProvenance` (wym) |
| `studio.journal.list` | Oddaje odwracalny dziennik czynności dokumentu — każdą czynność wraz z autorem, opisem i zależnościami, po których widać, czego cofnąć nie da się samodzielnie | `documentId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`kind:StudioActionKind` (opc)<br>`state:StudioActionState` (opc)<br>`limit:int` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `actions:StudioDocumentAction[]` (wym)<br>`total:int` (wym) |
| `studio.journal.redo` | Ponawia cofniętą czynność dokumentu — obejmuje treść i postać | `documentId:string` (wym)<br>`actionIds:string[]` (wym) | `redone:string[]` (wym)<br>`document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym) |
| `studio.journal.revert` | Cofa wskazaną czynność dokumentu pojedynczo, także nie po kolei — obejmuje treść i postać. Czynność będąca podstawą późniejszej odmawia i nazywa zależność, zamiast zostawić dokument w stanie niespójnym | `documentId:string` (wym)<br>`actionIds:string[]` (wym)<br>`createVersion:bool` (opc) | `reverted:string[]` (wym)<br>`document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym) |
| `studio.list.apply` | Zakłada wypunktowanie, numerację albo listę wielopoziomową na wskazanym fragmencie | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`kind:StudioListKind` (wym)<br>`listId:string` (opc)<br>`level:int` (opc)<br>`numberFormat:StudioListNumberFormat` (opc)<br>`startAt:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`list:StudioListDefinition` (wym) |
| `studio.list.bullet.set` | Ustawia znak wypunktowania poziomu — znak gotowy, dowolny symbol, ikonę albo obraz własny — wraz z wcięciem, odstępem i wyrównaniem znaku | `documentId:string` (wym)<br>`listId:string` (wym)<br>`level:int` (wym)<br>`bulletSource:StudioBulletSource` (opc)<br>`bulletCharacter:string` (opc)<br>`bulletAssetId:string` (opc)<br>`indentMm:float` (opc)<br>`hangingMm:float` (opc)<br>`align:StudioTextAlign` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`list:StudioListDefinition` (wym) |
| `studio.list.level.indent` | Zwiększa albo zmniejsza poziom listy wskazanego fragmentu — wcięcie i odstęp poziomu idą za nim | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`step:int` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.list.numbering.set` | Ustawia format numeracji poziomu listy — arabską, rzymską, literową albo prawniczą wielopoziomową — wraz ze wzorem numeru i punktem startu | `documentId:string` (wym)<br>`listId:string` (wym)<br>`level:int` (wym)<br>`numberFormat:StudioListNumberFormat` (opc)<br>`pattern:string` (opc)<br>`startAt:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`list:StudioListDefinition` (wym) |
| `studio.list.restart` | Wznawia numerację listy od wskazanego miejsca | `documentId:string` (wym)<br>`listId:string` (wym)<br>`offset:int` (wym)<br>`startAt:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.lock.add` | Zakłada blokadę na zaznaczonym fragmencie — fragment, którego model nie tknie ani w treści, ani w postaci. Blokada obowiązuje w rdzeniu, przed dotknięciem treści | `documentId:string` (wym)<br>`rangeStart:int` (wym)<br>`rangeEnd:int` (wym)<br>`name:string` (wym)<br>`reason:string` (opc)<br>`scope:StudioLockScope` (opc) | `lock:StudioFragmentLock` (wym)<br>`locks:StudioFragmentLock[]` (wym) |
| `studio.lock.list` | Oddaje wykaz blokad dokumentu wraz z nazwą, powodem, zakresem i zasięgiem | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `locks:StudioFragmentLock[]` (wym) |
| `studio.lock.remove` | Zdejmuje blokadę fragmentu. Blokadę zdejmuje wyłącznie Operator — czynność modelu wraca odmowa nazywająca powód | `documentId:string` (wym)<br>`lockId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `removed:bool` (wym)<br>`locks:StudioFragmentLock[]` (wym) |
| `studio.markup.add` | Znakuje fragment — wyróżnienie barwą, znacznik własny Operatora albo propozycja zmiany na marginesie z brzmieniem, bez wchodzenia w treść. Znakowanie modelu jest podpisane jako model | `documentId:string` (wym)<br>`kind:StudioMarkupKind` (wym)<br>`rangeStart:int` (wym)<br>`rangeEnd:int` (wym)<br>`color:string` (opc)<br>`markType:string` (opc)<br>`body:string` (opc)<br>`suggestedText:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `markup:StudioMarkup` (wym)<br>`form:StudioDocumentForm` (opc)<br>`actionId:string` (opc) |
| `studio.markup.decide` | Rozstrzyga propozycję modelu z marginesu — przyjęcie wnosi proponowane brzmienie do treści jako zmianę śledzoną, odrzucenie zostawia treść nietkniętą | `documentId:string` (wym)<br>`markupIds:string[]` (wym)<br>`accept:bool` (wym)<br>`editedText:string` (opc)<br>`createVersion:bool` (opc) | `decided:int` (wym)<br>`document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym) |
| `studio.markup.list` | Oddaje wykaz znakowań dokumentu jako spis do przejścia — komentarze, propozycje, wyróżnienia i znaczniki — z filtrem według rodzaju, autora i stanu | `documentId:string` (wym)<br>`kind:StudioMarkupKind` (opc)<br>`author:StudioAuthor` (opc)<br>`state:StudioMarkupState` (opc)<br>`markType:string` (opc)<br>`includeComments:bool` (opc)<br>`includeTrackedChanges:bool` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `markups:StudioMarkup[]` (wym)<br>`comments:StudioComment[]` (opc)<br>`trackedChanges:StudioTrackedChange[]` (opc)<br>`total:int` (wym)<br>`openCount:int` (wym) |
| `studio.markup.remove` | Zdejmuje znakowanie fragmentu — wyróżnienie, znacznik albo propozycję | `documentId:string` (wym)<br>`markupId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `removed:bool` (wym)<br>`form:StudioDocumentForm` (opc) |
| `studio.markup.type.delete` | Usuwa rodzaj znacznika własnego. Rodzaju fabrycznego nie usuwa — odpowiada odmowa nazywająca powód | `name:string` (wym) | `deleted:bool` (wym) |
| `studio.markup.type.list` | Oddaje wykaz rodzajów znaczników własnych wraz z barwą i liczbą użyć | `documentId:string` (opc) | `markupTypes:StudioMarkupType[]` (wym) |
| `studio.markup.type.save` | Zakłada albo zmienia rodzaj znacznika własnego Operatora — do sprawdzenia, wymaga źródła, gotowe — wraz z nazwą i barwą | `name:string` (wym)<br>`label:string` (opc)<br>`color:string` (opc) | `markupType:StudioMarkupType` (wym) |
| `studio.model.changes.list` | Oddaje wszystkie zmiany wniesione przez model w całym dokumencie — treści i postaci — wraz z licznikiem do przełącznika podświetlenia | `documentId:string` (wym)<br>`includeDecided:bool` (opc)<br>`onlyFormChanges:bool` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `summary:StudioModelChangeSummary` (wym) |
| `studio.model.changes.navigate` | Wskazuje następną albo poprzednią zmianę modelu wraz z miejscem, do którego okno ma przewinąć | `documentId:string` (wym)<br>`fromOffset:int` (opc)<br>`direction:string` (wym)<br>`onlyFormChanges:bool` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `change:StudioTrackedChange` (opc)<br>`index:int` (opc)<br>`total:int` (wym) |
| `studio.model.changes.revert` | Cofa zmiany modelu — wszystkie albo wskazane — z zachowaniem zmian Operatora naniesionych w tym czasie. Nie jest to przywrócenie wersji sprzed, bo to skasowałoby pracę Operatora | `documentId:string` (wym)<br>`changeIds:string[]` (opc)<br>`all:bool` (opc)<br>`createBackup:bool` (opc)<br>`createVersion:bool` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `revertedCount:int` (wym)<br>`keptOperatorChanges:int` (wym)<br>`document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`backupId:string` (opc) |
| `studio.object.format.set` | Ustawia postać obiektu — rozmiar, przycięcie, opływanie tekstem, położenie i zakotwiczenie, warstwę, tekst zastępczy, obramowanie, wypełnienie, obrót | `documentId:string` (wym)<br>`objectId:string` (wym)<br>`widthMm:float` (opc)<br>`heightMm:float` (opc)<br>`keepAspect:bool` (opc)<br>`crop:json` (opc)<br>`wrap:StudioTextWrap` (opc)<br>`anchor:StudioAnchorKind` (opc)<br>`anchorOffset:int` (opc)<br>`positionXMm:float` (opc)<br>`positionYMm:float` (opc)<br>`zOrder:int` (opc)<br>`altText:string` (opc)<br>`border:json` (opc)<br>`fillColor:string` (opc)<br>`strokeColor:string` (opc)<br>`strokeWidthPt:float` (opc)<br>`shadow:bool` (opc)<br>`rotationDeg:float` (opc)<br>`innerText:string` (opc)<br>`caption:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`object:StudioDocumentObject` (wym) |
| `studio.object.insert` | Wstawia obiekt w miejsce kursora — obraz z pliku, z magazynu zasobów rdzenia, z modułu Design albo z bazy zdjęciowej; kształt, ikonę, pole tekstowe albo logo | `documentId:string` (wym)<br>`offset:int` (wym)<br>`kind:StudioObjectKind` (wym)<br>`source:StudioObjectSource` (opc)<br>`assetId:string` (opc)<br>`designNodeId:string` (opc)<br>`libraryFileId:string` (opc)<br>`sourceUrl:string` (opc)<br>`path:string` (opc)<br>`bytesBase64:string` (opc)<br>`shapeKind:StudioShapeKind` (opc)<br>`iconName:string` (opc)<br>`widthMm:float` (opc)<br>`heightMm:float` (opc)<br>`altText:string` (opc)<br>`innerText:string` (opc)<br>`caption:string` (opc)<br>`wrap:StudioTextWrap` (opc)<br>`anchor:StudioAnchorKind` (opc)<br>`fillColor:string` (opc)<br>`strokeColor:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`object:StudioDocumentObject` (wym) |
| `studio.object.list` | Oddaje obiekty osadzone w dokumencie | `documentId:string` (wym)<br>`kind:StudioObjectKind` (opc) | `objects:StudioDocumentObject[]` (wym) |
| `studio.object.remove` | Usuwa obiekt osadzony w dokumencie | `documentId:string` (wym)<br>`objectId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `removed:bool` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`actionId:string` (opc) |
| `studio.operation.delete` | Usuwa operację własną Operatora. Operacji fabrycznej nie usuwa — na nią odpowiada odmowa nazywająca ten fakt | `operationId:string` (wym) | `deleted:bool` (wym) |
| `studio.operation.list` | Zwraca operacje kontekstowe Tools Panel: fabryczne wraz z ich promptami oraz zapisane przez Operatora | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc)<br>`category:string` (opc) | `operations:StudioOperation[]` (wym) |
| `studio.operation.save` | Zapisuje własny prompt Operatora jako nazwane narzędzie wielokrotnego użytku, dostępne w Tools Panel na równi z operacjami fabrycznymi | `operationId:string` (opc)<br>`name:string` (wym)<br>`category:string` (wym)<br>`prompt:string` (wym)<br>`scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `operation:StudioOperation` (wym) |
| `studio.package.export` | Pakuje paczkę redakcyjną przekazania: dokument finalny, historię wersji, raport zmian i adnotacje w jednym archiwum | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`includeHistory:bool` (opc)<br>`includeDiffReport:bool` (opc)<br>`includeAnnotations:bool` (opc)<br>`documentFormat:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`entries:int` (wym)<br>`sizeBytes:int` (wym) |
| `studio.page.break.insert` | Wstawia podział strony, kolumny, sekcji albo wiersza we wskazanym miejscu treści | `documentId:string` (wym)<br>`offset:int` (wym)<br>`kind:StudioBreakKind` (wym)<br>`sectionStart:StudioSectionStart` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.page.envelope.set` | Ustawia nadruk koperty — adres adresata i nadawcy wraz z ich położeniem na kopercie | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`recipient:string` (opc)<br>`sender:string` (opc)<br>`includeSender:bool` (opc)<br>`recipientXMm:float` (opc)<br>`recipientYMm:float` (opc)<br>`senderXMm:float` (opc)<br>`senderYMm:float` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`envelope:StudioEnvelopeSetup` (wym) |
| `studio.page.headerfooter.get` | Oddaje nagłówki i stopki dokumentu albo wskazanej sekcji według zasięgu | `documentId:string` (wym)<br>`sectionId:string` (opc) | `headersFooters:StudioHeaderFooter[]` (wym) |
| `studio.page.headerfooter.set` | Ustawia nagłówek i stopkę — osobno dla sekcji, dla pierwszej strony i dla stron parzystych | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`scope:StudioHeaderScope` (wym)<br>`headerText:string` (opc)<br>`footerText:string` (opc)<br>`headerDistanceMm:float` (opc)<br>`footerDistanceMm:float` (opc)<br>`linkedToPrevious:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`headersFooters:StudioHeaderFooter[]` (wym) |
| `studio.page.numbering.set` | Ustawia numerację stron — format, punkt startu, wznowienie w sekcji i położenie numeru | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`enabled:bool` (wym)<br>`format:StudioPageNumberFormat` (opc)<br>`startAt:int` (opc)<br>`restartInSection:bool` (opc)<br>`showTotal:bool` (opc)<br>`position:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`numbering:StudioPageNumbering` (wym) |
| `studio.page.paper.list` | Oddaje wykaz formatów nośnika do wyboru — szereg A i B, Letter, Legal, Tabloid oraz koperty DL, C4, C5, C6 | `kind:StudioPaperKind` (opc) | `papers:StudioPaperFormat[]` (wym) |
| `studio.page.setup.get` | Oddaje nastawy strony dokumentu albo wskazanej sekcji | `documentId:string` (wym)<br>`sectionId:string` (opc) | `pageSetup:StudioPageSetup` (wym)<br>`sections:StudioSection[]` (opc) |
| `studio.page.setup.set` | Ustawia nastawy strony dokumentu albo pojedynczej sekcji — nośnik, orientację, marginesy, margines na oprawę, kolumny. Zmiana formatu przelicza układ i oddaje bilans tego, co się nie zmieściło | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`paperName:string` (opc)<br>`widthMm:float` (opc)<br>`heightMm:float` (opc)<br>`orientation:StudioPageOrientation` (opc)<br>`marginTopMm:float` (opc)<br>`marginBottomMm:float` (opc)<br>`marginLeftMm:float` (opc)<br>`marginRightMm:float` (opc)<br>`gutterMm:float` (opc)<br>`mirrorMargins:bool` (opc)<br>`marginPreset:string` (opc)<br>`columns:int` (opc)<br>`columnGapMm:float` (opc)<br>`columnRule:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`pageSetup:StudioPageSetup` (wym) |
| `studio.page.watermark.set` | Ustawia znak wodny dokumentu albo sekcji — napis albo obraz, wraz z kryciem i obrotem | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`kind:StudioWatermarkKind` (wym)<br>`text:string` (opc)<br>`assetId:string` (opc)<br>`opacity:float` (opc)<br>`angleDeg:float` (opc)<br>`color:string` (opc)<br>`fontSizePt:float` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`watermark:StudioWatermark` (wym) |
| `studio.pdf.bates` | Nakłada numerację prawną Bates oraz nagłówki i stopki na wszystkich stronach dokumentu | `assetId:string` (wym)<br>`prefix:string` (opc)<br>`startNumber:int` (opc)<br>`digits:int` (opc)<br>`header:string` (opc)<br>`footer:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`lastNumber:int` (wym) |
| `studio.pdf.bookmarks.set` | Tworzy albo przestawia drzewo zakładek dokumentu PDF | `assetId:string` (wym)<br>`bookmarks:StudioPdfBookmark[]` (wym)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym) |
| `studio.pdf.extract` | Wyciąga z dokumentu PDF osadzone obrazy i załączniki jako osobne zasoby magazynu | `assetId:string` (wym)<br>`images:bool` (opc)<br>`attachments:bool` (opc)<br>`windowId:string` (opc) | `assetIds:string[]` (wym)<br>`extracted:int` (wym) |
| `studio.pdf.form.fill` | Odczytuje i wypełnia pola formularza AcroForm dokumentu PDF | `assetId:string` (wym)<br>`values:json` (opc)<br>`flatten:bool` (opc)<br>`windowId:string` (opc) | `fields:StudioPdfFormField[]` (wym)<br>`asset:DesignAsset` (opc) |
| `studio.pdf.merge` | Łączy wiele dokumentów PDF w jeden, w kolejności wskazanej przez Operatora | `assetIds:string[]` (wym)<br>`name:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`pages:int` (wym) |
| `studio.pdf.optimize` | Zmniejsza rozmiar dokumentu PDF przez rekompresję obrazów i oczyszczenie struktury | `assetId:string` (wym)<br>`imageQuality:int` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`sizeBytes:int` (wym)<br>`savedBytes:int` (wym) |
| `studio.pdf.pages.reorder` | Porządkuje strony dokumentu PDF: zmiana kolejności, obrót, usunięcie, wstawienie i wyodrębnienie | `assetId:string` (wym)<br>`operations:StudioPdfPageOperation[]` (wym)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`pages:int` (wym) |
| `studio.pdf.split` | Dzieli dokument PDF po stronach, zakresach albo zakładkach. Każda część jest osobnym zasobem | `assetId:string` (wym)<br>`ranges:string[]` (opc)<br>`windowId:string` (opc) | `assetIds:string[]` (wym)<br>`parts:int` (wym) |
| `studio.pdf.stamp` | Nakłada na strony dokumentu PDF znak wodny albo pieczęć tekstową lub graficzną | `assetId:string` (wym)<br>`text:string` (opc)<br>`stampAssetId:string` (opc)<br>`pages:string` (opc)<br>`opacity:int` (opc)<br>`rotation:int` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym) |
| `studio.plan.create` | Rozkłada zlecenie dokumentowe Operatora na zadania — podstawę pętli wykonawczej. Sam rozkład niczego nie uruchamia | `documentId:string` (wym)<br>`order:string` (wym)<br>`tasks:json` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc) | `plan:StudioTaskPlan` (wym) |
| `studio.plan.get` | Oddaje rozkład zlecenia dokumentowego wraz ze stanem każdego zadania i tym, komu je powierzono | `planId:string` (opc)<br>`documentId:string` (opc)<br>`state:StudioTaskState` (opc) | `plan:StudioTaskPlan` (wym)<br>`plans:StudioTaskPlan[]` (opc) |
| `studio.plan.run` | Puszcza pętlę wykonawczą na rozkładzie. Gdy nastawa pętli jest wyłączona, wraca odmowa nazywająca brak nastawy, a nie cisza | `planId:string` (wym)<br>`maxIterations:int` (opc)<br>`agentIds:string[]` (opc)<br>`dryRun:bool` (opc) | `plan:StudioTaskPlan` (wym)<br>`started:bool` (wym)<br>`refusalReason:string` (opc)<br>`balance:StudioActionBalance` (opc) |
| `studio.plan.stop` | Zatrzymuje pętlę wykonawczą rozkładu; zadania w biegu wracają do stanu czekania, a to, co zrobiono, zostaje | `planId:string` (wym)<br>`reason:string` (opc) | `plan:StudioTaskPlan` (wym)<br>`stoppedTasks:int` (wym) |
| `studio.plan.task.update` | Przestawia stan zadania rozkładu, powierza je wykonawcy albo odkłada jego wynik | `taskId:string` (wym)<br>`state:StudioTaskState` (opc)<br>`result:string` (opc)<br>`failureReason:string` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc)<br>`instruction:string` (opc) | `task:StudioDocumentTask` (wym)<br>`plan:StudioTaskPlan` (wym) |
| `studio.preview.render` | Renderuje dokument w formacie docelowym do podglądu: typografia, paginacja, nagłówek i stopka. Zwraca strony jako zasoby, więc podgląd pokazuje układ, a nie sam tekst | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`format:string` (wym)<br>`pageFrom:int` (opc)<br>`pageTo:int` (opc)<br>`profileId:string` (opc)<br>`watermark:string` (opc)<br>`windowId:string` (opc) | `pageAssetIds:string[]` (wym)<br>`pages:int` (wym) |
| `studio.proposal.decide` | Przyjmuje albo odrzuca propozycję zmiany po stronie rdzenia, w całości albo fragmentami. Dziś decyzja zapada wyłącznie w kliencie i rdzeń o niej nie wie, więc propozycja zostaje w nim nierozstrzygnięta | `documentId:string` (wym)<br>`proposalId:string` (wym)<br>`hunkIndexes:int[]` (opc)<br>`accept:bool` (wym)<br>`createVersion:bool` (opc) | `document:StudioDocument` (wym)<br>`version:StudioVersion` (opc) |
| `studio.provenance.list` | Oddaje pochodzenie fragmentów dokumentu — skąd każdy wniesiony fragment pochodzi; podstawa pod podobieństwa w panelu redaktora i pod bibliografię | `documentId:string` (wym)<br>`kind:StudioProvenanceKind` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `entries:StudioProvenance[]` (wym) |
| `studio.repository.export` | Wydaje pojedynczą wersję albo całą historię dokumentu jako archiwum z manifestem | `documentId:string` (wym)<br>`versionIds:string[]` (opc)<br>`format:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`entries:int` (wym)<br>`sizeBytes:int` (wym) |
| `studio.repository.list` | Zwraca historię wersji dokumentu w repozytorium sesji | `documentId:string` (wym)<br>`limit:int` (opc) | `versions:StudioVersion[]` (wym) |
| `studio.repository.restore` | Przywraca wcześniejszą wersję dokumentu bez usuwania wersji nowszych | `documentId:string` (wym)<br>`versionId:string` (wym) | `document:StudioDocument` (wym) |
| `studio.ruler.tabstop.set` | Zakłada, przestawia albo zdejmuje tabulator na linijce — z rodzajem i znakiem wiodącym | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`positionMm:float` (wym)<br>`kind:StudioTabKind` (opc)<br>`leader:StudioTabLeader` (opc)<br>`remove:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`tabStops:StudioTabStop[]` (wym) |
| `studio.search.semantic` | Wyszukuje w dokumencie fragmenty znaczeniowo bliskie zapytaniu, nie tylko dopasowania dosłowne. Uzupełnia wyszukiwanie wzorca w studio.diff.compare | `documentId:string` (wym)<br>`query:string` (wym)<br>`versionId:string` (opc)<br>`limit:int` (opc)<br>`minScore:float` (opc) | `matches:StudioSemanticMatch[]` (wym)<br>`mode:string` (wym) |
| `studio.section.delete` | Usuwa sekcję — jej treść przechodzi do sekcji poprzedniej wraz z jej nastawami | `documentId:string` (wym)<br>`sectionId:string` (wym)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `deleted:bool` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym) |
| `studio.section.list` | Oddaje sekcje dokumentu wraz z ich nastawami strony, nagłówkami i numeracją | `documentId:string` (wym) | `sections:StudioSection[]` (wym) |
| `studio.section.save` | Zakłada sekcję o własnych nastawach albo zmienia zastaną — pismo z załącznikiem w orientacji poziomej zostaje jednym dokumentem | `documentId:string` (wym)<br>`sectionId:string` (opc)<br>`title:string` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`start:StudioSectionStart` (opc)<br>`pageSetup:json` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`section:StudioSection` (wym) |
| `studio.security.encrypt` | Nakłada albo zdejmuje szyfrowanie dokumentu PDF wraz z uprawnieniami. Czynność jest jawnym, odwracalnym ustawieniem Operatora i niczego nie warunkuje | `assetId:string` (wym)<br>`userPassword:string` (opc)<br>`ownerPassword:string` (opc)<br>`permissions:string[]` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`encrypted:bool` (wym) |
| `studio.security.metadata.strip` | Czyści autora, historię, dane ukryte i komentarze dokumentu przed wysyłką | `assetId:string` (wym)<br>`keepFields:string[]` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`removed:string[]` (wym) |
| `studio.security.redact` | Trwale zamazuje wskazane fragmenty wraz z warstwą tekstową pod spodem. Czynność jest nieodwracalna dla wyniku, dlatego źródło zostaje nietknięte | `assetId:string` (wym)<br>`regions:StudioRedactionRegion[]` (wym)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`redacted:int` (wym) |
| `studio.security.sensitive.detect` | Oznacza w treści numery identyfikacyjne, numery kart, adresy poczty i dane osobowe wymagające redakcji. Czynność czyta i niczego nie zmienia | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`categories:string[]` (opc) | `findings:StudioSensitiveFinding[]` (wym) |
| `studio.security.sign` | Podpisuje dokument PDF certyfikatem Operatora | `assetId:string` (wym)<br>`certificateId:string` (wym)<br>`reason:string` (opc)<br>`location:string` (opc)<br>`windowId:string` (opc) | `asset:DesignAsset` (wym)<br>`signedAt:int64` (wym) |
| `studio.security.sign.verify` | Sprawdza podpisy dokumentu PDF i zwraca wynik weryfikacji każdego z nich | `assetId:string` (wym) | `signatures:StudioSignature[]` (wym)<br>`allValid:bool` (wym) |
| `studio.style.apply` | Stosuje styl nazwany do wskazanego fragmentu | `documentId:string` (wym)<br>`name:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.style.delete` | Usuwa styl nazwany własny. Stylu fabrycznego nie usuwa — odpowiada odmowa nazywająca powód | `documentId:string` (wym)<br>`name:string` (wym)<br>`replaceWith:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `deleted:bool` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym) |
| `studio.style.list` | Oddaje arkusz stylów nazwanych dokumentu wraz z dziedziczeniem i liczbą miejsc użycia | `documentId:string` (wym)<br>`kind:StudioStyleKind` (opc)<br>`includeBuiltin:bool` (opc) | `styles:StudioNamedStyle[]` (wym) |
| `studio.style.save` | Zakłada styl nazwany albo zmienia zastany. Zmiana stylu przestawia wszystkie miejsca dokumentu, które go używają | `documentId:string` (wym)<br>`name:string` (wym)<br>`displayName:string` (opc)<br>`kind:StudioStyleKind` (opc)<br>`basedOn:string` (opc)<br>`nextStyle:string` (opc)<br>`character:json` (opc)<br>`paragraph:json` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `style:StudioNamedStyle` (wym)<br>`form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`actionId:string` (opc)<br>`change:StudioTrackedChange` (opc) |
| `studio.symbol.autoreplace.list` | Oddaje wykaz zasad autozamiany skrótów na znaki | `includeBuiltin:bool` (opc) | `rules:StudioAutoReplaceRule[]` (wym) |
| `studio.symbol.autoreplace.set` | Ustawia zasadę autozamiany skrótu na znak — nastawa jawna i odwracalna, z własnym wykazem Operatora | `shortcut:string` (wym)<br>`replacement:string` (opc)<br>`enabled:bool` (opc) | `rule:StudioAutoReplaceRule` (opc)<br>`removed:bool` (opc) |
| `studio.symbol.insert` | Wstawia znak specjalny w miejsce kursora — znak matematyczny, walutę, literę grecką, strzałkę, znak prawniczy, znak diakrytyczny albo interpunkcyjny niedostępny z klawiatury | `documentId:string` (wym)<br>`offset:int` (wym)<br>`character:string` (opc)<br>`code:string` (opc)<br>`name:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`symbol:StudioSymbol` (wym) |
| `studio.symbol.list` | Oddaje tablicę znaków do wybrania wraz z wyszukaniem po nazwie i po kodzie oraz znakami ostatnio użytymi | `query:string` (opc)<br>`category:string` (opc)<br>`recentOnly:bool` (opc)<br>`limit:int` (opc) | `symbols:StudioSymbol[]` (wym)<br>`categories:string[]` (opc) |
| `studio.table.convert` | Zamienia tekst na tabelę albo tabelę na tekst | `documentId:string` (wym)<br>`direction:StudioTableConvert` (wym)<br>`tableId:string` (opc)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc)<br>`separator:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`table:StudioDocumentTable` (opc) |
| `studio.table.format.set` | Ustawia postać tabeli i komórek — szerokości kolumn, obramowanie, cieniowanie, wyrównanie w komórce, styl tabeli, powtarzanie wiersza nagłówkowego | `documentId:string` (wym)<br>`tableId:string` (wym)<br>`row:int` (opc)<br>`column:int` (opc)<br>`columnWidthsMm:float[]` (opc)<br>`widthMm:float` (opc)<br>`styleName:string` (opc)<br>`border:json` (opc)<br>`shadingColor:string` (opc)<br>`verticalAlign:StudioVerticalAlign` (opc)<br>`align:StudioTextAlign` (opc)<br>`headerRows:int` (opc)<br>`repeatHeader:bool` (opc)<br>`caption:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`table:StudioDocumentTable` (wym) |
| `studio.table.insert` | Zakłada tabelę o wskazanym rozmiarze — przez podanie liczby wierszy i kolumn albo z szablonu gotowego | `documentId:string` (wym)<br>`offset:int` (wym)<br>`rows:int` (wym)<br>`columns:int` (wym)<br>`styleName:string` (opc)<br>`headerRows:int` (opc)<br>`repeatHeader:bool` (opc)<br>`widthMm:float` (opc)<br>`cells:json` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`table:StudioDocumentTable` (wym) |
| `studio.table.list` | Oddaje tabele dokumentu wraz z komórkami, szerokościami kolumn i stylem | `documentId:string` (wym)<br>`tableId:string` (opc) | `tables:StudioDocumentTable[]` (wym) |
| `studio.table.sort` | Sortuje zawartość tabeli po wskazanej kolumnie, z pominięciem wiersza nagłówkowego | `documentId:string` (wym)<br>`tableId:string` (wym)<br>`column:int` (wym)<br>`descending:bool` (opc)<br>`numeric:bool` (opc)<br>`skipHeader:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`table:StudioDocumentTable` (wym) |
| `studio.table.structure.edit` | Zmienia budowę tabeli — wstawia i usuwa wiersz oraz kolumnę, scala i dzieli komórki. Szerokości kolumn zostają policzone, nie zerowe | `documentId:string` (wym)<br>`tableId:string` (wym)<br>`operation:StudioTableStructureOp` (wym)<br>`row:int` (opc)<br>`column:int` (opc)<br>`count:int` (opc)<br>`rowSpan:int` (opc)<br>`columnSpan:int` (opc)<br>`before:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc)<br>`table:StudioDocumentTable` (wym) |
| `studio.template.apply` | Zakłada dokument z szablonu, wypełniając jego pola wskazanymi wartościami | `windowId:string` (wym)<br>`templateId:string` (wym)<br>`values:json` (opc)<br>`title:string` (opc) | `document:StudioDocument` (wym) |
| `studio.template.delete` | Usuwa szablon własny. Szablonu fabrycznego nie usuwa — odpowiada odmowa nazywająca powód | `templateId:string` (wym) | `deleted:bool` (wym) |
| `studio.template.export` | Oddaje szablon do pliku — dotx albo ott — wraz z arkuszem stylów, nastawami strony i polami | `templateId:string` (wym)<br>`format:StudioExportFormat` (opc)<br>`path:string` (opc) | `result:StudioExportResult` (wym) |
| `studio.template.field.list` | Oddaje pola szablonu do wypełnienia wraz z rodzajem, wartością domyślną i wymagalnością | `templateId:string` (wym) | `fields:StudioTemplateFieldSpec[]` (wym) |
| `studio.template.field.set` | Wskazuje pole szablonu do wypełnienia — nazwę, opis, wartość domyślną, rodzaj i to, czy jest wymagane | `templateId:string` (wym)<br>`name:string` (wym)<br>`label:string` (opc)<br>`description:string` (opc)<br>`kind:StudioTemplateFieldKind` (opc)<br>`required:bool` (opc)<br>`defaultValue:string` (opc)<br>`choices:string[]` (opc)<br>`anchorOffset:int` (opc)<br>`remove:bool` (opc) | `template:StudioTemplateDetail` (wym)<br>`fields:StudioTemplateFieldSpec[]` (wym) |
| `studio.template.fill` | Wypełnia pola szablonu wartościami i oddaje dokument gotowy — czynność dostępna także modelowi: zrób z tego wzór pisma i wypełnij dla tej sprawy | `templateId:string` (wym)<br>`documentId:string` (opc)<br>`windowId:string` (opc)<br>`values:json` (wym)<br>`title:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`missingRequired:string[]` (opc)<br>`balance:StudioActionBalance` (wym) |
| `studio.template.import` | Wnosi szablon z pliku Operatora — dotx albo ott — wraz z arkuszem stylów, nastawami strony, nagłówkiem, stopką i polami do wypełnienia | `name:string` (opc)<br>`format:StudioImportFormat` (opc)<br>`path:string` (opc)<br>`libraryFileId:string` (opc)<br>`bytesBase64:string` (opc)<br>`category:string` (opc) | `template:StudioTemplateDetail` (wym)<br>`balance:StudioImportBalance` (wym) |
| `studio.template.list` | Zwraca szablony dokumentów dostępne Operatorowi: pismo, umowa, raport, notatka, oferta oraz szablony własne | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `templates:StudioTemplate[]` (wym) |
| `studio.template.save` | Zakłada szablon pisma z bieżącego dokumentu wraz z arkuszem stylów, nastawami strony, nagłówkiem, stopką, logo, tabelami i blokadami wzorcowymi | `documentId:string` (opc)<br>`templateId:string` (opc)<br>`name:string` (wym)<br>`description:string` (opc)<br>`category:string` (opc)<br>`fields:json` (opc)<br>`includeLocks:bool` (opc)<br>`thumbnailAssetId:string` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `template:StudioTemplateDetail` (wym) |
| `studio.text.edit` | Zmienia treść wskazanego fragmentu bez przepisywania całości; blokada fragmentu jest sprawdzana w rdzeniu przed dotknięciem treści | `documentId:string` (wym)<br>`rangeStart:int` (wym)<br>`rangeEnd:int` (wym)<br>`text:string` (wym)<br>`keepFormat:bool` (opc)<br>`author:StudioAuthor` (opc)<br>`agentId:string` (opc)<br>`agentName:string` (opc)<br>`subagentId:string` (opc) | `form:StudioDocumentForm` (wym)<br>`balance:StudioActionBalance` (wym)<br>`change:StudioTrackedChange` (opc)<br>`actionId:string` (opc) |
| `studio.text.get` | Oddaje treść wskazanego fragmentu wraz z jego postacią znaku i akapitu | `documentId:string` (wym)<br>`rangeStart:int` (opc)<br>`rangeEnd:int` (opc) | `text:string` (wym)<br>`runs:StudioDocumentRun[]` (opc)<br>`paragraph:StudioParagraphFormat` (opc)<br>`character:StudioCharacterFormat` (opc)<br>`locks:StudioFragmentLock[]` (opc) |
| `studio.tracking.decide` | Przyjmuje albo odrzuca zmiany zarejestrowane przez śledzenie, pojedynczo albo grupą. Decyzja zapisuje się w treści dokumentu i zakłada wersję | `documentId:string` (wym)<br>`changeIds:string[]` (wym)<br>`accept:bool` (wym)<br>`createVersion:bool` (opc) | `document:StudioDocument` (wym)<br>`decided:int` (wym) |
| `studio.tracking.list` | Zwraca zmiany zarejestrowane przez śledzenie wraz z autorem każdej z nich | `documentId:string` (wym)<br>`versionId:string` (opc)<br>`pendingOnly:bool` (opc)<br>`agentId:string` (opc)<br>`subagentId:string` (opc) | `changes:StudioTrackedChange[]` (wym) |
| `studio.tracking.set` | Włącza albo wyłącza śledzenie zmian w dokumencie. Przy włączonym śledzeniu każdy zapis odkłada wstawienia i usunięcia jako zmiany do decyzji, zamiast nadpisywać treść | `documentId:string` (wym)<br>`enabled:bool` (wym) | `enabled:bool` (wym) |
| `studio.version.label.set` | Nadaje wersji nazwę własną albo ją zdejmuje i oznacza wersję jako kluczową | `versionId:string` (wym)<br>`label:string` (opc)<br>`milestone:bool` (opc) | `version:StudioVersion` (wym) |
| `studio.version.reference.create` | Tworzy odwołanie do konkretnej wersji dokumentu w obrębie platformy, o zasięgu widoczności wskazanym przez Operatora | `versionId:string` (wym)<br>`scope:ConfigScope` (opc)<br>`expiresAt:int64` (opc) | `reference:string` (wym)<br>`scope:ConfigScope` (wym) |
| `studio.version.restore.initial` | Wraca dokumentem do wersji założycielskiej jednym poleceniem, bez szukania jej w wykazie; wersje nowsze zostają, więc samo cofnięcie jest odwracalne | `documentId:string` (wym)<br>`createVersion:bool` (opc) | `document:StudioDocument` (wym)<br>`form:StudioDocumentForm` (wym)<br>`version:StudioVersion` (wym) |
| `studio.version.series.list` | Oddaje wersje dokumentu w rozbiciu na szereg Operatora i szereg autozapisu, żeby zapisy samoczynne były odróżnialne w wykazie | `documentId:string` (wym)<br>`series:StudioVersionSeries` (opc)<br>`milestonesOnly:bool` (opc)<br>`limit:int` (opc) | `versions:StudioVersion[]` (wym)<br>`operatorCount:int` (wym)<br>`autosaveCount:int` (wym)<br>`initialVersionId:string` (opc) |
| `studio.view.get` | Oddaje nastawy widoku okna pracy z dokumentem pamiętane przy dokumencie — tryb powierzchni, skalę, linijki, układ stron i przewijanie | `documentId:string` (opc)<br>`windowId:string` (opc) | `settings:StudioViewSettings` (wym) |
| `studio.view.set` | Ustawia nastawy widoku okna pracy z dokumentem. Gdzie da się zrobić dwojako i obie drogi mają sens, wybór należy do Operatora i jest jawnym, odwracalnym ustawieniem | `documentId:string` (opc)<br>`windowId:string` (opc)<br>`surfaceMode:StudioSurfaceMode` (opc)<br>`splitOrientation:StudioSplitOrientation` (opc)<br>`splitRatio:float` (opc)<br>`viewMode:StudioViewMode` (opc)<br>`zoomPercent:int` (opc)<br>`zoomPreset:StudioZoomPreset` (opc)<br>`rulersVisible:bool` (opc)<br>`rulerUnit:StudioRulerUnit` (opc)<br>`marginGuides:bool` (opc)<br>`formattingMarks:bool` (opc)<br>`pagesPerRow:int` (opc)<br>`spreadView:bool` (opc)<br>`scrollMode:StudioScrollMode` (opc)<br>`modelChangesHighlighted:bool` (opc)<br>`toolboxVisible:bool` (opc) | `settings:StudioViewSettings` (wym) |

**Zdarzenia obszaru `studio` — 3:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `studio.chain.progressed` | Postęp łańcucha operacji albo wsadu; zasila pętlę wykonawczą okna | `runId:string`, `stepIndex:int`, `stepCount:int`, `state:string`, `proposalId:string` |
| `studio.document.changed` | Zmiana dokumentu Studio; odświeża Preview Window i Session Repository | `change:ChangeKind`, `document:StudioDocument` |
| `studio.ingest.changed` | Zmiana stanu pozycji kolejki wczytywania; odświeża Ingest/OCR Panel bez odpytywania w pętli | `windowId:string`, `item:StudioIngestItem` |


Razem w wykazie: **179 komend** z 1 obszaru kontraktu.

---

*Koniec dokumentu. Moduł Studio — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
