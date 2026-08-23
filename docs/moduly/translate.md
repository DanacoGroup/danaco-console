# Moduł Translate — dokumentacja projektowa

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
| **Tytuł** | Moduł Translate — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować) |
| **Środowiska dostępności** | TalkIn |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowiska TalkIn; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 4, 9.7, 12, Załącznik A) · Specyfikacja modułów 1.1 (rozdz. 4.7) · Specyfikacja okien operacyjnych 1.1 (rozdz. 5, 6.7, 10) · System wizualny (rozdz. 6, 8, 10) · Model danych 1.1 (rozdz. 8, 17, Załącznik F.4) · Model konfiguracji 1.0 (rozdz. 5) · Izolacja i zależności (rozdz. 1–6) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); zero blokad, klucze jawne; izolacja i uprawnienia pozostają w gestii Operatora; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Katalog funkcji i narzędzi](#2-katalog-funkcji-i-narzędzi)
3. [Komplet okien operacyjnych modułu](#3-komplet-okien-operacyjnych-modułu)
4. [Specyfikacja okien operacyjnych](#4-specyfikacja-okien-operacyjnych)
5. [Przepływy pracy](#5-przepływy-pracy)
6. [Stany, dane i powiązania](#6-stany-dane-i-powiązania)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
8. [Scenariusze użycia](#8-scenariusze-użycia)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja i rola modułu

Translate jest **stanowiskiem tłumacza i lokalizatora** platformy Danaco Console: jednym miejscem, w którym powstaje, jest weryfikowana, ujednolicana terminologicznie i wydawana treść wielojęzyczna — od pojedynczego zdania, przez dokument z zachowanym formatem, po komplet plików lokalizacyjnych produktu. Moduł łączy w jednej przestrzeni tłumaczenie maszynowe, pracę wspomaganą (CAT), zarządzanie terminologią, pamięć tłumaczeń, korektę językową, edycję napisów oraz konwersję formatów lokalizacyjnych.

Praca prowadzona jest w trybie wielozadaniowym: jeden tekst źródłowy tłumaczony jest jednocześnie na wiele języków docelowych, widocznych obok siebie w osobnych panelach, zamiast sekwencyjnie, fragment po fragmencie, w pojedynczym oknie komunikacji. Moduł nadaje pracy tłumaczeniowej strukturę wymaganą przy pracy zespołowej: spójną terminologię, równoległość języków i pełną widoczność tekstu źródłowego przez cały czas pracy.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Translate |
|---|---|
| Zespoły tłumaczeniowe i lokalizacyjne | Równoległe tłumaczenie tego samego materiału na wiele rynków językowych |
| Działy komunikacji międzynarodowej | Spójna terminologia korporacyjna we wszystkich wersjach językowych |
| Użytkownicy indywidualni wielojęzyczni | Szybkie przełączanie się między językami komunikacji bez utraty kontekstu |
| Zespoły redakcyjne pracujące ze Studio | Wersje wielojęzyczne dokumentu redagowanego równolegle w module Studio |
| Zespoły lokalizacji oprogramowania | Obieg zasobów lokalizacyjnych z ochroną kluczy, zmiennych i reguł językowych |

### 1.3. Po co — wartość modułu

| Problem tłumaczenia sekwencyjnego | Rozwiązanie w Translate |
|---|---|
| Tłumaczenie fragment po fragmencie w oknie komunikacji | Translation Panels tłumaczą cały tekst źródłowy na wiele języków jednocześnie |
| Niespójna terminologia między tłumaczeniami | Glossary & Termbase Manager wymusza jeden, zdefiniowany odpowiednik pojęcia we wszystkich językach |
| Utrata śledzenia, co dokładnie jest tekstem źródłowym | Source Panel pozostaje stale widoczny, niezależny od liczby aktywnych języków docelowych |
| Ręczne odtwarzanie kontekstu przy każdej zmianie tekstu | Zmiana w Source Panel uruchamia aktualizację wszystkich paneli tłumaczeń |
| Rozproszenie pracy między osobne narzędzia formatów, terminologii i kontroli jakości | Jedna przestrzeń łączy silniki tłumaczenia, pamięć tłumaczeń, terminologię, formaty i kontrolę jakości |

### 1.4. Granica wewnętrzna — obszary pokrywane przez moduł

| Obszar tematyczny | Zakres w module Translate |
|---|---|
| Tłumaczenie maszynowe i wspomagane (MT/CAT) | Tłumaczenie równoległe na wiele języków, post-edycja, tłumaczenie pivotowe, tłumaczenie zwrotne |
| Pamięć tłumaczeń (Translation Memory) | Gromadzenie, dopasowanie rozmyte, konkordancja, wyrównanie tekstów dwujęzycznych |
| Terminologia i glosariusze | Bazy terminologiczne, wymuszanie odpowiedników, wykrywanie niespójności, listy „nie tłumacz” |
| Korekta i redakcja językowa | Sprawdzanie gramatyki, stylu, rejestru, czytelności — dla każdego języka docelowego osobno |
| Tłumaczenie dokumentów z zachowaniem formatu | DOCX, PDF, PPTX, XLSX, ODT, Markdown, HTML z odtworzeniem układu, stylów i osadzeń |
| Lokalizacja oprogramowania | Formaty zasobów: XLIFF, PO/POT, JSON, YAML, RESX, Android XML, iOS `.strings`/`.stringsdict`, Java `.properties` |
| Napisy i treść mówiona | SRT, WebVTT, TTML, STL; synteza mowy (odsłuch), transkrypcja, dubbing kierowany |
| Kontrola jakości tłumaczenia (QA/LQA) | Liczby, daty, waluty, placeholdery, długość, spójność, pominięcia, znaki niedrukowalne |

### 1.5. Granica zewnętrzna — obszary poza modułem

| Obszar | Dlaczego poza Translate | Gdzie na platformie |
|---|---|---|
| Redakcja jednojęzyczna dokumentu (pisanie, przepisywanie, streszczanie) | To praca autorska, nie tłumaczeniowa | Moduł **Studio** (powiązanie pary modułów Studio ◄──► Translate) |
| Trwałe magazynowanie i katalogowanie plików wynikowych | Translate wydaje artefakty, nie jest repozytorium | Moduł **Library** |
| Badania wieloźródłowe i raporty z tłumaczeń obcojęzycznych źródeł | To praca badawcza | Moduł **Research** |
| Cykliczne, bezobsługowe tłumaczenie wsadowe wg harmonogramu | To orkiestracja procesu, nie edycja | Moduł **Automations** (Translate udostępnia operacje jako kroki) |
| Generowanie grafiki z lokalizowanym tekstem | To praca wizualna | Moduł **Design** |

Zasada rozgraniczenia: **Translate odpowiada za treść w wielu językach; sąsiednie moduły odpowiadają za jej pochodzenie, przechowywanie i dalszy obieg.** Powiązania są jawne i konfigurowalne (rozdz. 7), nigdy wbudowane na stałe.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn
        ▼
ŚRODOWISKO: TalkIn ── boczna nawigacja modułów
        │
        ▼
MODUŁ: TRANSLATE ────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Source Panel ·
  Translation Panels (N) · Glossary & Termbase Manager ·
  Translation Memory Panel · Format Studio · QA & Review Center
        │
        ▼
KARTA SESJI — tekst źródłowy i zestaw języków docelowych
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 6.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowisko, w którym moduł jest widoczny w bocznej nawigacji | Wyłącznie TalkIn |
| Środowiska WorkSpace, CodeStudio | Moduł niedostępny |
| Komponent własny | Translate nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 8 (łącznie z Chat Window i Execution Loop Window); liczba instancji Translation Panels zależna od liczby wybranych języków |
| Charakter pracy | Wielozadaniowa — wiele języków docelowych aktualizowanych równolegle |

---

## 2. Katalog funkcji i narzędzi

Każda pozycja: **nazwa · co robi · zależności (biblioteki / formaty / integracje)**. Wszystkie funkcje działają domyślnie (zero blokad); klucze dostępowe silników zewnętrznych są jawne i wprowadzane w oknie konfiguracji (zakres 5.4 „Zachowanie modeli”, 5.7 „Rozszerzenia”).

### 2.1. Silniki tłumaczenia i tryby pracy

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.1.1 | **Multi-Engine Translation** | Tłumaczenie jednego segmentu przez wiele silników równolegle (model platformy + zewnętrzne API) z porównaniem wariantów | Kanał modelu (Architektura rozdz. 9); adaptery API silników tłumaczenia; klucze z okna konfiguracji |
| 2.1.2 | **Parallel Panels (N języków)** | Ten sam tekst źródłowy tłumaczony jednocześnie na wiele języków w siatce paneli | Mechanizm bazowy modułu (rozdz. 4.4) |
| 2.1.3 | **Post-Editing Mode (MTPE)** | Tryb korekty tłumaczenia maszynowego z rejestrem różnic MT→wersja finalna i klasyfikacją zmian (light/full PE) | Silnik diff (`sergi/go-diff`); metryka edit-distance |
| 2.1.4 | **Pivot Translation** | Tłumaczenie przez język pośredni, gdy brak bezpośredniej pary (np. PL→EN→JA), z oznaczeniem trasy | Graf par językowych; polityka wyboru pivota w konfiguracji |
| 2.1.5 | **Back-Translation** | Tłumaczenie zwrotne wyniku na język źródłowy jako kontrola sensu, prezentowane obok oryginału | Silnik MT; widok trójkolumnowy (źródło · tłumaczenie · zwrotne) |
| 2.1.6 | **Register/Tone per Panel** | Niezależny rejestr każdego panelu: formalny, nieformalny, techniczny, marketingowy, prawniczy, dziecięcy | Model bazowy; presety tonu w konfiguracji |
| 2.1.7 | **Bulk/Batch Translate** | Tłumaczenie zbioru plików lub segmentów jednym zleceniem, z kolejką postępu | Kolejka zadań (Automations); tryb w tle |
| 2.1.8 | **Adaptive/Domain MT** | Dostrojenie tłumaczenia do dziedziny i do zatwierdzonej pamięci tłumaczeń (styl klienta) | TM (2.3); glosariusz (2.4); kontekst modelu |

### 2.2. Segmentacja, wyrównanie i struktura tekstu

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.2.1 | **Sentence/Segment Engine (SRX)** | Podział tekstu na jednostki tłumaczeniowe wg reguł segmentacji (zdania/akapity), numeracja wspólna dla wszystkich paneli | `neurosnap/sentences` (Punkt), reguły SRX; `blevesearch/segment` (granice ICU) |
| 2.2.2 | **Bitext Aligner** | Wyrównanie istniejącego tekstu źródłowego i jego tłumaczenia w pary segmentów → budowa TM z gotowych materiałów | Algorytm w stylu Hunalign; podobieństwo długości + kotwice; eksport TMX |
| 2.2.3 | **Concordance Search** | Wyszukanie fragmentu w całej pamięci tłumaczeń wraz z kontekstem i wcześniejszym przekładem | Indeks pełnotekstowy (`blevesearch/bleve`); TM |
| 2.2.4 | **Re-segmentation & Merge/Split** | Ręczne łączenie i dzielenie segmentów bez utraty powiązań z panelami | Model segmentów modułu; synchronizacja paneli |
| 2.2.5 | **Tag/Placeholder Protection** | Wykrycie i zablokowanie znaczników formatu, zmiennych i placeholderów (`{name}`, `%s`, `<b>`) tak, by nie zostały przetłumaczone | Parser wzorców; walidacja liczby i kolejności tagów źródło↔cel |

### 2.3. Pamięć tłumaczeń (Translation Memory)

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.3.1 | **TM Store & Fuzzy Match** | Zapamiętuje zatwierdzone pary segmentów i podpowiada je z dopasowaniem rozmytym (procent podobieństwa nad progiem) | Levenshtein (`agnivade/levenshtein`) + podobieństwo semantyczne (embeddingi modelu); próg z konfiguracji |
| 2.3.2 | **TMX Import/Export** | Wymiana pamięci z zewnętrznymi narzędziami w standardzie TMX 1.4b | `encoding/xml`; schemat TMX 1.4b |
| 2.3.3 | **Context (101%) Match** | Rozróżnia dopasowanie w identycznym kontekście (poprzedni/następny segment) od zwykłego 100% | Metadane kontekstu w rekordzie TM |
| 2.3.4 | **TM Maintenance** | Czyszczenie duplikatów, scalanie pamięci, masowa podmiana, filtrowanie wg pola | Operacje wsadowe na indeksie TM |
| 2.3.5 | **Leverage/Pre-translate** | Wstępne wypełnienie nowego materiału trafieniami z TM przed uruchomieniem MT | TM (2.3.1); reguły progu i blokady niskich trafień |

### 2.4. Terminologia i glosariusze

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.4.1 | **Termbase Manager** | Baza terminów z odpowiednikiem per język, kategorią dziedziny, definicją i przykładem użycia | Model danych terminu; kategorie dziedzin |
| 2.4.2 | **TBX/CSV Import/Export** | Wymiana bazy terminologicznej w standardzie TBX oraz w CSV | `encoding/xml` (TBX), `encoding/csv` |
| 2.4.3 | **Term Enforcement** | Automatyczne wstawianie zatwierdzonego odpowiednika przy generowaniu tłumaczeń | Glosariusz + wstrzyknięcie do promptu/API |
| 2.4.4 | **Inconsistency Detection** | Wskazuje, gdzie ten sam termin źródłowy przetłumaczono różnie w jednym języku | Analiza wystąpień w panelach |
| 2.4.5 | **Do-Not-Translate (DNT)** | Oznacza nazwy marek, jednostki, kody jako niepodlegające tłumaczeniu | Lista DNT; ochrona przy MT (2.2.5) |
| 2.4.6 | **Term Extraction** | Wskazuje kandydatów na terminy z tekstu źródłowego (statystyka + model) do zasilenia bazy | Analiza częstości + model; słowa stop per język |
| 2.4.7 | **Forbidden/Preferred Terms** | Listy terminów zabronionych i preferowanych zgodnie z przewodnikiem stylu klienta | Reguły w termbase; walidacja QA |

### 2.5. Korekta i redakcja językowa

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.5.1 | **Grammar & Spelling Check** | Sprawdzanie gramatyki, ortografii i interpunkcji dla każdego języka docelowego osobno | Silnik LanguageTool (HTTP API) jako rozszerzenie lub kontrola modelem; słowniki per język |
| 2.5.2 | **Style & Register Check** | Ocena zgodności rejestru z zamierzonym tonem panelu (formalność, długość zdań, strona bierna) | Reguły stylu + model; presety tonu (2.1.6) |
| 2.5.3 | **Readability Scoring** | Wskaźniki czytelności per język (Flesch, LIX, Gunning Fog, Jasnopis dla PL) | Metryki czytelności; parametry językowe |
| 2.5.4 | **False-Friends & Interference** | Wykrywa kalki, fałszywych przyjaciół i interferencję z języka źródłowego | Baza par ryzykownych; analiza modelem |
| 2.5.5 | **Punctuation & Typography Locale** | Poprawia cudzysłowy, spacje, myślniki i typografię zgodnie z konwencją języka docelowego | Reguły typograficzne CLDR per locale |
| 2.5.6 | **Consistency Checker** | Wykrywa niespójne tłumaczenie identycznych zdań i różnice w wersji tego samego terminu | Indeks segmentów; TM; termbase |

### 2.6. Kontrola jakości tłumaczenia (QA / LQA)

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.6.1 | **Numeric/Date/Currency Check** | Sprawdza zgodność liczb, dat, walut i jednostek między źródłem a tłumaczeniem, z uwzględnieniem konwencji locale | `golang.org/x/text/number`, `x/text/currency`, CLDR |
| 2.6.2 | **Placeholder/Tag Integrity** | Weryfikuje kompletność i kolejność placeholderów oraz znaczników | Parser tagów (2.2.5) |
| 2.6.3 | **Length/Fit Check** | Kontrola limitu znaków/pikseli (UI, napisy, meta) i sygnalizacja przekroczeń | Metryka długości; profile limitów (napisy: CPS, znaki/linię) |
| 2.6.4 | **Omission/Untranslated Detect** | Wykrywa segmenty pominięte, nieprzetłumaczone lub identyczne ze źródłem | Analiza statusów segmentów |
| 2.6.5 | **LQA Scorecard (MQM/DQF)** | Ocena jakości wg modelu błędów MQM/DQF z kategorią, wagą i wynikiem końcowym | Model kategorii MQM; formularz oceny; eksport raportu |
| 2.6.6 | **QA Profiles** | Zestawy reguł QA włączane per projekt/klient (co sprawdzać, jaka waga ostrzeżeń) | Konfiguracja reguł (rozdz. 7) |

### 2.7. Tłumaczenie dokumentów z zachowaniem formatu

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.7.1 | **DOCX Round-Trip** | Ekstrahuje tekst z DOCX, tłumaczy, odtwarza dokument ze stylami, tabelami, nagłówkami i osadzeniami | `unidoc/unioffice` (OOXML read/write) |
| 2.7.2 | **PDF Extract & Reflow** | Wydobywa tekst i układ z PDF (w tym warstwowych), tłumaczy, składa PDF z zachowanym rozkładem | `gen2brain/go-fitz` (MuPDF: tekst z pozycją), `pdfcpu/pdfcpu` (składanie i stemplowanie) |
| 2.7.3 | **OCR for Scanned Docs** | Rozpoznaje tekst ze skanów i obrazów przed tłumaczeniem, z detekcją języka i układu | `otiai10/gosseract` (Tesseract) |
| 2.7.4 | **PPTX/XLSX Translate** | Tłumaczy prezentacje i arkusze z zachowaniem slajdów, notatek, komórek i formuł tekstowych | `unidoc/unioffice` (PPTX), `qax-os/excelize` (XLSX) |
| 2.7.5 | **Markdown/HTML Translate** | Tłumaczy treść, chroniąc składnię MD/HTML, linki, bloki kodu i atrybuty | `yuin/goldmark` (MD AST), `golang.org/x/net/html` (HTML) |
| 2.7.6 | **ODT/ODF Translate** | Tłumaczy dokumenty OpenDocument z zachowaniem stylów | Parser XML ODF (`encoding/xml` po rozpakowaniu ZIP) |
| 2.7.7 | **Layout Diff / Visual Compare** | Porównuje wygląd dokumentu przed i po tłumaczeniu (przepełnienia, złamany układ) | Render `go-fitz`; różnica obrazów; `hexops/gotextdiff` |
| 2.7.8 | **Bilingual Export** | Wydaje plik dwujęzyczny (źródło + tłumaczenie) w XLIFF, DOCX-tabela lub TXT | `encoding/xml` (XLIFF 2.1), `unidoc/unioffice` (DOCX) |

### 2.8. Lokalizacja oprogramowania i treści cyfrowych

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.8.1 | **Localization File Formats** | Import/eksport zasobów lokalizacyjnych z ochroną kluczy i zmiennych | Parsery: JSON, YAML (`gopkg.in/yaml.v3`), `.properties`, Android `strings.xml`, iOS `.strings`/`.stringsdict`, RESX, gettext PO (`leonelquinteros/gotext`) |
| 2.8.2 | **XLIFF 1.2 / 2.1 Round-Trip** | Pełen obieg standardu wymiany lokalizacyjnej z metadanymi stanu, notatkami i TM | `encoding/xml`; schematy XLIFF 1.2 i 2.1 |
| 2.8.3 | **Pseudolocalization** | Generuje pseudotłumaczenie (wydłużenie, znaki diakrytyczne, ramki) do testu UI przed realnym tłumaczeniem | Transformacja znaków; reguły ekspansji |
| 2.8.4 | **Plural & Gender Rules** | Obsługa form liczby mnogiej i rodzaju wg reguł CLDR (ICU MessageFormat) | `golang.org/x/text/feature/plural`, CLDR plural rules |
| 2.8.5 | **Placeholder Style Mapping** | Rozpoznaje i mapuje style zmiennych między platformami (`%s` ↔ `{0}` ↔ `{{name}}`) | Reguły wzorców; walidacja |
| 2.8.6 | **Key Context & Screenshots** | Wiąże klucz lokalizacyjny z kontekstem i zrzutem ekranu miejsca użycia | Powiązanie z Library/Design (integracja) |

### 2.9. Napisy, treść mówiona i głos

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.9.1 | **Subtitle Formats** | Import/eksport i tłumaczenie napisów z zachowaniem taktowania | `asticode/go-astisub` (SRT, WebVTT, TTML, STL) |
| 2.9.2 | **Timing & CPS Control** | Kontrola liczby znaków na sekundę, długości linii i minimalnego czasu wyświetlania | Metryki CPS; reguły napisowe per język |
| 2.9.3 | **Text-to-Speech Readback** | Odsłuch tłumaczenia syntezą mowy jako kontrola brzmieniowa i podkład do dubbingu | API TTS; głos per język w konfiguracji |
| 2.9.4 | **Speech-to-Text Dictation** | Dyktowanie tłumaczenia głosem oraz transkrypcja nagrań źródłowych | API STT |
| 2.9.5 | **Dubbing/Voiceover Script** | Przygotowanie skryptu dubbingu z segmentacją kwestii, oznaczeniem mówców i długością pod lip-sync | Segmentacja (2.2); TTS (2.9.3); limity długości |

### 2.10. Praca zespołowa, wydanie i integracje

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.10.1 | **Vendor Handoff Package** | Pakuje materiał dla zewnętrznego tłumacza (XLIFF + TM + termbase + instrukcje) i przyjmuje zwrot | XLIFF (2.8.2), TMX (2.3.2), TBX (2.4.2), ZIP |
| 2.10.2 | **Review & Approval Workflow** | Przebieg: tłumaczenie → korekta → akceptacja, ze statusem i śladem autora zmiany | Statusy segmentów; historia (Model danych) |
| 2.10.3 | **Word/Char Count & Cost Estimate** | Statystyka objętości z analizą TM (nowe / fuzzy / powtórzenia) i wyceną wg stawek | Analiza względem TM; siatka stawek w konfiguracji |
| 2.10.4 | **Studio Bridge** | Odbiór zaznaczenia ze Studio jako tekstu źródłowego i zwrot wersji wielojęzycznej do Studio Editor | Powiązanie pary modułów Studio ◄──► Translate (rozdz. 7) |
| 2.10.5 | **Library Sync** | Zapis wydanych plików dwujęzycznych i TM/termbase jako artefaktów do modułu Library | Integracja Translate ◄──► Library |
| 2.10.6 | **Automations Steps** | Udostępnia operacje (tłumacz, QA, eksport) jako kroki procesu wsadowego uruchamianego przez Automations | Kolejka zadań; tryb w tle |
| 2.10.7 | **Translation API/Webhook** | Wystawia tłumaczenie i QA jako operacje wywoływane z zewnątrz (rozszerzenie/konektor) | Kontrakt rozszerzeń (Architektura rozdz. 14); MCP/konektor |

---

## 3. Komplet okien operacyjnych modułu

Układ modułu jest wyłącznie pionowy — okna rozmieszczone są w kolumnach sąsiadujących poziomo. Lewa kolumna, stała, pełna wysokość obszaru roboczego, należy do Chat Window; kolumna sąsiadująca otwiera Execution Loop Window; prawa, dominująca część przestrzeni należy do okien roboczych modułu; okna zarządcze otwierane są jako rozszerzenia boczne po prawej stronie obszaru roboczego.

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Główne okno komunikacji, centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji przy aktywnym zleceniu; wskaźnik przebiegu pętli zawsze widoczny | Pętla wykonawcza zadań tłumaczeniowych: dekompozycja zlecenia, kolejka zadań per język, kontrola jakości, ponowienia |
| 3 | Source Panel | Okno edycyjne | Lewa część obszaru roboczego, stały punkt odniesienia | 1 | Widoczne bez interakcji | Tekst źródłowy, segmentacja, ochrona tagów |
| 4 | Translation Panels | Okno edycyjne (instancja wielokrotna) | Obszar roboczy, siatka równoległych kolumn językowych | 1 | Widoczne bez interakcji | Tłumaczenia równoległe na wiele języków |
| 5 | Glossary & Termbase Manager | Okno zarządcy | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Znacznik terminologii, menu `⋮`, skrót `Ctrl/Cmd + G` | Spójność terminologii, bazy terminów, listy DNT |
| 6 | Translation Memory Panel | Okno zarządcy | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `⋮` obszaru roboczego, polecenie języka naturalnego | Pamięć tłumaczeń, konkordancja, wyrównanie, pre-translate |
| 7 | Format Studio | Okno robocze z podglądem | Kolumna sąsiadująca w obszarze roboczym, otwierana przy pracy z dokumentem | 2 | Znacznik formatu dokumentu w pasku kontekstu | Tłumaczenie z zachowaniem formatu, OCR, formaty lokalizacyjne |
| 8 | QA & Review Center | Okno zarządcy | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `⋮`, polecenie języka naturalnego, skrót `Ctrl/Cmd + Shift + Q` | Kontrola jakości, korekta językowa, przebieg akceptacji |

**Makieta zbiorcza — stan spoczynku interfejsu.**

```
 Makieta zbiorcza — Moduł Translate       Dostępność: TalkIn
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window        │ Source Panel  │ Translation Panels
  nawigacja │ Użytkownik ↔       │ [1] …         │ ┌── EN ──┬── DE ──┬── FR ──┐
  modułów   │ Wykonawca          │ [2] …         │ │[1] …   │[1] …   │[1] …   │
            │                    │ [3] ▌…▐       │ │[2] …   │[2] …   │[2] …   │
            │ ────────────────   │ [4] ⚠ …       │ │[3] ▌…▐ │[3] …   │[3] …   │
            │ Execution Loop     │               │ │[4] ⚠   │[4] ⚠   │[4] ⚠   │
            │ Koordynator ↔      │               │ └────────┴────────┴────────┘
            │ Wykonawca          │               │
            │ zadania: 3/12 ▶    │               │  [ + Język ▼ ]        [ ⋮ ]
 ═══════════════════════════════════════════════════════════════════════════
  [Danaco Console] [TalkIn] [Translate] [PL → EN·DE·FR] [Kanał modelu ▼]  [☰]
 ═══════════════════════════════════════════════════════════════════════════
```

W stanie spoczynku widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3: znaczniki kontekstowe paska kontekstu, `▼`, `⋮`, `☰`. Glossary & Termbase Manager, Translation Memory Panel, Format Studio i QA & Review Center pozostają zwinięte do chwili wywołania i po zamknięciu znikają całkowicie z przestrzeni roboczej.

### 3.1. Warstwy widoczności w module

| Warstwa | Co należy do warstwy w module Translate | Sposób dostępu |
|---|---|---|
| 1 — zawsze widoczna | Chat Window, Execution Loop Window ze wskaźnikiem przebiegu pętli, Source Panel z segmentami, siatka Translation Panels, statusy segmentów, pasek kontekstu ze znacznikami języków | Bez interakcji |
| 2 — widoczna na żądanie | Wybór kanału modelu i silnika tłumaczenia, dobór języków i wariantów, selektor tonu panelu, próg dopasowania pamięci tłumaczeń, wybór profilu QA, Format Studio | Znacznik kontekstowy, przycisk `+ Język ▼`, selektor tonu `▼`; element zwija się po użyciu |
| 3 — rozwinięcia kontekstowe | Glossary & Termbase Manager, Translation Memory Panel, QA & Review Center, operacje eksportu i importu, akcje segmentu, konkordancja, ujednolicenie terminu, pakiet dla wykonawcy zewnętrznego | Menu `⋮` przy segmencie i przy panelu, menu `☰` obszaru roboczego, panel popover podpowiedzi pamięci tłumaczeń |
| 4 — funkcje eksperckie | Reguły segmentacji SRX, mapowanie stylów placeholderów, konfiguracja pivota, konserwacja pamięci tłumaczeń, karta oceny LQA (MQM/DQF), pseudolokalizacja, wystawienie operacji jako API/webhook, macierz izolacji technicznej | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, okno konfiguracji |

**Zasada jednego kliknięcia.** Każdy z powyższych elementów osiągalny jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window.

---

## 4. Specyfikacja okien operacyjnych

### 4.1. Chat Window — komunikacja Użytkownik ↔ Wykonawca

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja (kanał Użytkownik ↔ Wykonawca) |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Warstwa | 1 |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 6.3) |

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą, centralnym punktem pracy w module i podstawowym mechanizmem sterowania wszystkimi procesami tłumaczeniowymi: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie i przerywanie działań oraz wyjaśnianie wyniku i kontekstu.

**Zawartość i pełny arsenał funkcji.**

- Polecenia w kontekście tłumaczenia: „dostosuj ton wersji niemieckiej do bardziej formalnego rejestru”, „sprawdź spójność terminologii we wszystkich panelach”, „wyjaśnij, dlaczego francuskie tłumaczenie różni się długością”.
- Strumień odpowiedzi Wykonawcy na żywo kanałem WebSocket.
- Odwołania do konkretnego języka lub panelu tłumaczenia przez wzmiankę w treści polecenia.
- Polecenia zbiorcze stosowane do wszystkich aktywnych paneli jednocześnie („skróć wszystkie wersje o 15%”).
- Zlecanie operacji całego modułu: uruchomienie tłumaczenia wsadowego, kontroli jakości, eksportu zbiorczego, pakietu dla wykonawcy zewnętrznego.
- Zatwierdzanie i przerywanie działań prowadzonych w pętli wykonawczej.
- Wyjaśnienia różnic kulturowych i idiomatycznych między wersjami językowymi.
- Wstawienie odpowiedzi bezpośrednio do wskazanego Translation Panel.
- Dostęp do funkcji warstwy 4 poleceniem języka naturalnego (reguły SRX, konserwacja pamięci tłumaczeń, karta oceny LQA).
- Historia poleceń, regeneracja odpowiedzi.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Chat Window ────────────────┐
│ Użytkownik ↔ Wykonawca  [ ⋮ ]│
├──────────────────────────────┤
│ ┌─ Użytkownik ─────────────┐ │
│ │ „sprawdź spójność        │ │
│ │  terminologii”           │ │
│ └──────────────────────────┘ │
│ ┌─ Wykonawca ──────────────┐ │
│ │ rozbieżność: „klient”    │ │
│ │ → DE: 2 warianty …       │ │
│ │ [ Glossary ] [ Popraw ]  │ │
│ └──────────────────────────┘ │
├──────────────────────────────┤
│ [Kanał modelu ▼]  [ ☰ ]      │
│ Pole poleceń ……… [ Wyślij ▶ ]│
└──────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Strumień komunikacji | Lista dymków Użytkownika i Wykonawcy | Prezentacja poleceń i wyników | Ciało kolumny | 1 | Widoczny bez interakcji | domyślny · strumieniowanie | Przewija się do najnowszej odpowiedzi |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia w języku naturalnym | Duże pole w stopce kolumny | 1 | Widoczne bez interakcji | domyślny · fokus | Enter wysyła polecenie |
| Selektor kanału modelu | Element zwinięty `Kanał modelu ▼` | Wybór modelu tłumaczącego bieżącej karty | Mały wyzwalacz w stopce kolumny | 2 | Kliknięcie znacznika; zwija się po wyborze | zwinięty · rozwinięty · ładowanie | Wybór zmienia model dla kolejnych tłumaczeń |
| Menu operacji | Menu `☰` | Zbiór operacji modułu wywoływanych z okna komunikacji | Jeden element zbiorczy | 3 | Kliknięcie `☰` | zwinięte · rozwinięte | Rozwija listę operacji (tłumaczenie wsadowe, QA, eksport) |
| Menu kontekstowe dymka | Menu `⋮` | Warianty operacji na odpowiedzi | Mały wyzwalacz | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Regeneracja, kopiowanie, wstawienie do wskazanego panelu |
| Przycisk „Glossary” | Mały przycisk `--zarys` | Przejście do rozbieżności terminologicznej | Przycisk w pasku akcji dymka | 3 | Widoczny w dymku wyniku | domyślny | Otwiera Glossary & Termbase Manager na wskazanym terminie |
| Przycisk „Popraw” | Mały przycisk `--zloty` | Zastosowanie spójnego terminu we wszystkich panelach | Mały przycisk CTA | 3 | Widoczny w dymku wyniku | domyślny · zastosowano | Ujednolica termin we wszystkich Translation Panels |

---

### 4.2. Execution Loop Window — komunikacja Koordynator ↔ Wykonawca

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja (kanał Koordynator ↔ Wykonawca) |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |
| Warstwa | 1 (wskaźnik przebiegu pętli oraz kolejka zadań); warstwa 3 dla wariantów sterowania przebiegiem |
| Izolacja domyślna | Pętla właściwa bieżącej karcie sesji |

Okno prowadzi pętlę wykonawczą zadań tłumaczeniowych modułu. Koordynator dekomponuje zlecenie Użytkownika na zadania właściwe Translate — segmentację tekstu źródłowego, tłumaczenie per język docelowy, wymuszenie terminologii, kontrolę jakości, korektę językową, odtworzenie formatu i eksport — przydziela je Wykonawcom, nadzoruje realizację i decyduje o ponowieniu zadania, którego wynik nie przeszedł kontroli.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie tłumaczeniowe i jego dekompozycja na zadania: segmentacja, tłumaczenie panelu EN, tłumaczenie panelu DE, wymuszenie terminologii, kontrola jakości panelu, odtworzenie formatu, eksport.
- Kolejka i stan zadań per język docelowy oraz per zakres segmentów, z licznikiem ukończonych zadań.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą: przydział zadania, raport wykonania, zgłoszenie niezgodności terminologicznej, żądanie ponowienia.
- Wyniki kontroli jakości zadania (liczby, placeholdery, długość, pominięcia) i decyzje o ponowieniu tłumaczenia segmentu.
- Wskaźniki przebiegu pętli: liczba zadań ukończonych, w toku i ponowionych, czas przebiegu, wykorzystanie silników tłumaczenia.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia bez utraty wykonanych zadań.
- Ślad realizacji zadania powiązany z segmentem — przejście z pozycji kolejki do odpowiadającego segmentu w Translation Panel.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Execution Loop ─────────────┐
│ Koordynator ↔ Wykonawca [ ⋮ ]│
├──────────────────────────────┤
│ Zlecenie: tłumaczenie PL →   │
│ EN · DE · FR, 18 segmentów   │
├──────────────────────────────┤
│ ✓ segmentacja SRX            │
│ ✓ terminologia — wymuszona   │
│ ▶ tłumaczenie EN   14/18     │
│ ▶ tłumaczenie DE   12/18     │
│ ▶ tłumaczenie FR    8/18     │
│ ⟳ QA panelu EN — ponowienie  │
│   segmentu [4]               │
├──────────────────────────────┤
│ Pętla: 3/12 zadań · 1 ponow. │
│ [ Wstrzymaj ] [ Przerwij ]   │
└──────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok tekstu ze zleceniem i jego zakresem | Identyfikacja bieżącego zlecenia tłumaczeniowego | Blok w nagłówku kolumny | 1 | Widoczny bez interakcji | domyślny · zlecenie skorygowane | Kliknięcie rozwija pełną treść zlecenia |
| Kolejka zadań | Lista pozycji zadań ze stanem | Wgląd w dekompozycję zlecenia i postęp per język | Ciało kolumny | 1 | Widoczna bez interakcji | oczekujące · w toku · ukończone · ponowione · przerwane | Kliknięcie pozycji przewija Translation Panel do odpowiadających segmentów |
| Wskaźnik przebiegu pętli | Licznik zadań i ponowień | Ocena stanu realizacji zlecenia | Mały pasek w stopce kolumny | 1 | Widoczny bez interakcji | w toku · ukończony · wstrzymany | Kliknięcie przewija do pierwszego zadania nieukończonego |
| Komunikat sterujący | Dymek wymiany Koordynator ↔ Wykonawca | Prezentacja przydziału, raportu i decyzji o ponowieniu | Mały dymek w kolejce | 1 | Widoczny przy zadaniu | przydział · raport · niezgodność · ponowienie | Rozwija szczegóły kontroli jakości zadania |
| Przycisk „Wstrzymaj” | Przycisk `--zarys` | Wstrzymanie pętli wykonawczej bez utraty wyników | Mały przycisk w stopce kolumny | 1 | Widoczny bez interakcji | domyślny · wstrzymana (etykieta „Wznów”) | Zatrzymuje przydział kolejnych zadań |
| Przycisk „Przerwij” | Przycisk `--zarys` | Zakończenie bieżącego zlecenia | Mały przycisk w stopce kolumny | 1 | Widoczny bez interakcji | domyślny | Kończy pętlę, zachowuje segmenty przetłumaczone |
| Menu sterowania przebiegiem | Menu `⋮` | Warianty sterowania: korekta zlecenia, zmiana kolejności zadań, wymuszenie ponowienia, zmiana silnika dla zadania | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji na pętli |
| Rejestr przebiegu pętli | Pełny zapis komunikatów sterujących i decyzji | Analiza przebiegu i przyczyn ponowień | Panel wysuwany | 4 | Polecenie języka naturalnego w Chat Window lub skrót klawiszowy | zwinięty · rozwinięty | Otwiera pełny rejestr przebiegu bieżącego zlecenia |

---

### 4.3. Source Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne |
| Waga wizualna | Lewa część obszaru roboczego, stały punkt odniesienia dla wszystkich paneli tłumaczeń |
| Warstwa | 1 |
| Izolacja domyślna | Tekst źródłowy odrębny per karta sesji |

**Zawartość i pełny arsenał funkcji.**

- Wprowadzenie tekstu źródłowego przez wpisanie, wklejenie, przeciągnięcie pliku lub import (TXT, DOCX, PDF, PPTX, XLSX, ODT, Markdown, HTML, plik lokalizacyjny), a także wynik OCR ze skanu.
- Automatyczne rozpoznanie języka źródłowego z ręczną korektą; wykrycie treści wielojęzycznej w jednym pliku.
- Segmentacja SRX na jednostki tłumaczeniowe (zdania/akapity) z numeracją segmentów wspólną dla wszystkich paneli; ręczne łączenie i dzielenie segmentów.
- Ochrona tagów i placeholderów — podświetlenie elementów nietłumaczalnych.
- Podświetlenie segmentu odpowiadającego fragmentowi aktualnie edytowanemu w dowolnym Translation Panel (synchronizacja wzajemna).
- Licznik słów, znaków i segmentów; analiza względem pamięci tłumaczeń (nowe / fuzzy / powtórzenia); szacowany czas i koszt.
- Historia zmian tekstu źródłowego — każda zmiana uruchamia aktualizację wszystkich paneli tłumaczeń, z wyraźnym oznaczeniem segmentów wymagających ponownego tłumaczenia.
- Oznaczenie segmentu zatwierdzonego we wszystkich językach: segment pozostaje edytowalny, a jego zmiana oznacza tłumaczenia jako wymagające rewizji.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Source Panel ──────────────────────┐
│ [ Polski ▼ ]   słów 342 · segm. 18  │
├─────────────────────────────────────┤
│ [1] Tekst źródłowy pierwszego       │
│     segmentu do tłumaczenia…        │
│ [2] Kolejny segment tekstu…         │
│ [3] ▌segment tłumaczony w DE▐       │
│ [4] Segment zmieniony ⚠             │
│  …                                  │
├─────────────────────────────────────┤
│ [ Wejście ▼ ]                [ ⋮ ]  │
└─────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Segment tekstu źródłowego | Ponumerowany blok tekstu | Jednostka tłumaczeniowa odpowiadająca fragmentom we wszystkich panelach | Blok z numerem, ciało panelu | 1 | Widoczny bez interakcji | domyślny · aktualnie tłumaczony (podświetlony) · zmieniony/wymaga rewizji (obramowanie ostrzegawcze) | Kliknięcie przewija wszystkie Translation Panels do odpowiadającego segmentu |
| Licznik słów i segmentów | Mały tekst pomocniczy | Orientacja w rozmiarze materiału | Bardzo mały tekst `.dn-tekst-3` | 1 | Widoczny bez interakcji | aktualizowany na żywo | Kliknięcie rozwija analizę względem pamięci tłumaczeń |
| Znacznik tagu i placeholdera | Podświetlenie fragmentu nietłumaczalnego | Ochrona znaczników formatu i zmiennych | Wyróżnienie w tekście segmentu | 1 | Widoczny bez interakcji | chroniony · naruszony | Wskazuje odpowiadający tag w panelach tłumaczeń |
| Selektor języka źródłowego | Element zwinięty `Polski ▼` | Wskazanie lub korekta języka tekstu źródłowego | Znacznik w nagłówku panelu | 2 | Kliknięcie znacznika; zwija się po wyborze | rozpoznany automatycznie · ręcznie ustawiony | Zmiana wpływa na segmentację i tłumaczenia |
| Wyzwalacz wejścia | Element zbiorczy `Wejście ▼` | Wczytanie tekstu źródłowego (wklejenie, import pliku, OCR skanu) | Jeden element zwinięty w stopce panelu | 2 | Kliknięcie `▼`; zwija się po wyborze | zwinięty · rozwinięty · ładowanie | Wypełnia panel treścią wyekstrahowaną ze wskazanego źródła |
| Menu operacji na tekście | Menu `⋮` | Segmentacja ponowna, łączenie i dzielenie segmentów, historia zmian źródła | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Przelicza numerację i granice segmentów we wszystkich panelach |
| Reguły segmentacji SRX | Zestaw reguł podziału tekstu | Dostrojenie granic segmentów do materiału i języka | Formularz reguł | 4 | Polecenie języka naturalnego w Chat Window lub okno konfiguracji | domyślne · własne | Przebudowuje segmentację całego materiału |

---

### 4.4. Translation Panels

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne (instancja wielokrotna) |
| Waga wizualna | Obszar roboczy, siatka równoległych kolumn — jedna na każdy wybrany język docelowy |
| Warstwa | 1 |
| Izolacja domyślna | Każdy panel odrębny w obrębie karty; wspólny tekst źródłowy z Source Panel |

**Zawartość i pełny arsenał funkcji.**

- Dodanie i usunięcie języka docelowego oraz wariantu (EN-US/EN-GB, PT-BR/PT-PT) — nowa kolumna pojawia się w siatce natychmiast po wyborze.
- Tłumaczenie wielosilnikowe generowane po każdej zmianie tekstu źródłowego, z porównaniem wariantów i ręczną korektą każdego segmentu.
- Widok segmentowy odpowiadający 1:1 segmentom Source Panel, z synchronicznym przewijaniem wszystkich kolumn.
- Podpowiedzi z pamięci tłumaczeń (dopasowanie rozmyte z procentem) oraz z terminologii (podświetlenie wymuszonych odpowiedników).
- Tryb post-edycji z rejestrem różnic między wynikiem maszynowym a wersją finalną.
- Kontrola jakości per panel: liczby, daty, symbole waluty, placeholdery, długość względem limitu, segmenty pominięte.
- Tłumaczenie zwrotne — podgląd tłumaczenia z powrotem na język źródłowy jako kontrola sensu.
- Odsłuch tłumaczenia syntezą mowy oraz dyktowanie korekty głosem.
- Wybór rejestru i tonu per panel niezależnie (formalny, nieformalny, techniczny, marketingowy, prawniczy, dziecięcy).
- Status segmentu: nieprzetłumaczony, przetłumaczony automatycznie, zweryfikowany przez człowieka, wymaga rewizji.
- Eksport panelu jako plik dwujęzyczny (XLIFF, DOCX, TXT) albo w formacie docelowym z zachowaniem układu.
- Widok porównawczy dwóch paneli obok siebie niezależnie od Source Panel.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Translation Panels ────────────────────────────────┐
│ [ + Język ▼ ]   Synchronizacja ✓            [ ⋮ ]   │
├──────────────┬──────────────┬───────────────────────┤
│ EN  ton ▼    │ DE  ton ▼    │ FR  ton ▼             │
│ 14/18 ✓      │ 12/18 ✓      │ 8/18 ✓                │
├──────────────┼──────────────┼───────────────────────┤
│[1] translated│[1] übersetzt │[1] texte traduit…     │
│[2] …         │[2] …         │[2] …                  │
│[3] ▌edycja▐  │[3] …         │[3] …                  │
│[4] ⚠ rewizja │[4] ⚠ rewizja │[4] ⚠ rewizja          │
└──────────────┴──────────────┴───────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Segment tłumaczenia | Edytowalny blok tekstu odpowiadający segmentowi źródłowemu | Prezentacja i korekta tłumaczenia | Blok z numerem, ciało kolumny | 1 | Widoczny bez interakcji | nieprzetłumaczony · automatyczny · zweryfikowany · wymaga rewizji (obramowanie ostrzegawcze) | Edycja ręczna oznacza segment jako zweryfikowany |
| Nagłówek kolumny językowej | Pasek z kodem języka i statusem | Identyfikacja panelu i postęp tłumaczenia | Pasek na szczycie każdej kolumny | 1 | Widoczny bez interakcji | domyślny | Kliknięcie kodu języka otwiera wybór wariantu (EN-US/EN-GB) |
| Wskaźnik statusu tłumaczenia | Licznik ułamkowy z ikoną `ptaszek` | Postęp weryfikacji segmentów panelu | Mały tekst z ikoną | 1 | Widoczny bez interakcji | w toku · ukończony | Przewija do pierwszego segmentu nieukończonego |
| Przełącznik synchronizacji przewijania | `.dn-suwak` | Zsynchronizowane przewijanie wszystkich kolumn razem z Source Panel | Mały przełącznik w nagłówku okna | 1 | Widoczny bez interakcji | włączony (domyślny) · wyłączony | Włącza i wyłącza wspólne przewijanie |
| Wyzwalacz „+ Język” | Element zwinięty `+ Język ▼` z listą języków | Dodanie kolumny tłumaczenia równoległego | Mały wyzwalacz CTA | 2 | Kliknięcie `▼`; zwija się po wyborze | zwinięty · rozwinięty | Otwiera nową kolumnę w siatce i uruchamia pierwsze tłumaczenie |
| Selektor tonu panelu | Element zwinięty `ton ▼` | Niezależny wybór rejestru językowego dla kolumny | Znacznik w nagłówku kolumny | 2 | Kliknięcie znacznika; zwija się po wyborze | zwinięty · rozwinięty | Generuje ponownie tłumaczenie segmentów w wybranym tonie |
| Selektor silnika tłumaczenia | Znacznik kontekstowy silnika | Wybór silnika lub zestawu silników dla kolumny | Znacznik w pasku kontekstu | 2 | Kliknięcie znacznika | zwinięty · rozwinięty | Przełącza silnik dla kolejnych tłumaczeń kolumny |
| Podpowiedź pamięci tłumaczeń | Panel popover nad segmentem | Podpowiedź z wcześniej przetłumaczonego, zbliżonego segmentu | Mały popover z procentem dopasowania | 3 | Fokus segmentu; dopasowanie powyżej progu | ukryty · widoczny | Kliknięcie wstawia podpowiedź do segmentu |
| Menu segmentu | Menu `⋮` przy segmencie | Warianty operacji: tłumaczenie zwrotne, porównanie wariantów silników, odsłuch TTS, konkordancja, historia post-edycji | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji segmentu |
| Menu operacji panelu | Menu `⋮` okna | Kontrola jakości, eksport panelu, widok porównawczy dwóch paneli, usunięcie kolumny | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Podświetla segmenty z wykrytymi niezgodnościami; generuje plik XLIFF/DOCX/TXT |
| Mapowanie stylów placeholderów | Reguły konwersji zmiennych | Przekład stylu zmiennych między platformami | Formularz reguł | 4 | Polecenie języka naturalnego lub okno konfiguracji | domyślne · własne | Konwertuje zmienne w całym materiale |

---

### 4.5. Glossary & Termbase Manager

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządcy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 3 |
| Izolacja domyślna | Baza terminologiczna wspólna dla wszystkich Translation Panels bieżącej karty; udostępnienie szersze konfigurowalne |

**Zawartość i pełny arsenał funkcji.**

- Termin źródłowy z odpowiednikiem per język, kategorią dziedziny, definicją, przykładem użycia i statusem (zatwierdzony, kandydat, zabroniony).
- Kategoryzacja terminów wg dziedziny: prawnicze, techniczne, medyczne, marketingowe, nazwy własne.
- Automatyczne stosowanie zdefiniowanego odpowiednika przy generowaniu tłumaczeń w Translation Panels.
- Wykrywanie niespójności — wskazanie miejsc, gdzie ten sam termin źródłowy przetłumaczono różnie w obrębie jednego języka.
- Ekstrakcja kandydatów na terminy z tekstu źródłowego.
- Listy preferowane i zabronione zgodne z przewodnikiem stylu klienta.
- Import i eksport bazy w formatach TBX i CSV.
- Historia zmian wpisu słownikowego oraz adnotacja uzasadniająca wybór odpowiednika.
- Wyszukiwanie w bazie, filtrowanie wg dziedziny lub języka.
- Oznaczenie terminu jako „nie tłumacz” (nazwy marek, jednostki miary o utrwalonym zapisie).
- Podgląd wszystkich wystąpień terminu we wszystkich aktywnych panelach jednym kliknięciem.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Glossary & Termbase ───────────────┐
│ [ 🔍 szukaj ]      [ Dziedzina ▼ ]  │
├─────────────────────────────────────┤
│ „klient”                      [ ⋮ ] │
│   EN customer · DE Kunde · FR client│
│   ⚠ niespójność w DE: „Klientel” 2× │
├─────────────────────────────────────┤
│ „Danaco Console” 🔒 nie tłumacz [ ⋮ ]│
│   EN/DE/FR: Danaco Console          │
├─────────────────────────────────────┤
│ [ + Termin ]                 [ ☰ ]  │
└─────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Pozycja terminu | Karta `.dn-karta--pozycja` | Reprezentuje jeden termin z odpowiednikami we wszystkich językach | Średnia karta z listą par język–odpowiednik | 1 (w otwartym oknie) | Widoczna po otwarciu okna | domyślny · niespójność wykryta (obramowanie ostrzegawcze) · „nie tłumacz” (ikona kłódki) | Kliknięcie przewija panele do wystąpień terminu |
| Pole wyszukiwania terminu | Pole `.dn-input` z ikoną `szukaj` | Odnalezienie wpisu w bazie | Małe pole w nagłówku kolumny | 1 (w otwartym oknie) | Widoczne po otwarciu okna | puste · z wynikami | Filtruje listę terminów na żywo |
| Ikona „nie tłumacz” | Mała ikona `klodka` | Oznaczenie terminu jako niepodlegającego tłumaczeniu | Bardzo mała ikona przy nazwie terminu | 1 (w otwartym oknie) | Widoczna przy pozycji | nieaktywna · aktywna | Kliknięcie przełącza status terminu |
| Filtr dziedziny | Element zwinięty `Dziedzina ▼` | Zawężenie bazy do jednej kategorii tematycznej | Znacznik w nagłówku kolumny | 2 | Kliknięcie znacznika; zwija się po wyborze | wszystkie (domyślny) · wybrana dziedzina | Filtruje widoczne wpisy |
| Przycisk „+ Termin” | Przycisk `--zarys` | Ręczne dodanie wpisu do bazy | Mały przycisk w stopce kolumny | 2 | Widoczny po otwarciu okna | domyślny | Otwiera formularz terminu źródłowego i odpowiedników per język |
| Menu akcji terminu | Menu `⋮` przy pozycji | Ujednolicenie, podgląd wystąpień, edycja, historia wpisu, oznaczenie zabronionego | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Podmienia niespójne wystąpienia w Translation Panels; podświetla wystąpienia |
| Menu bazy | Menu `☰` | Import i eksport TBX/CSV, ekstrakcja terminów, listy preferowane i zabronione | Jeden element zbiorczy | 3 | Kliknięcie `☰` | zwinięte · rozwinięte | Import wczytuje plik TBX/CSV; eksport generuje plik do pobrania |
| Reguły wymuszania terminologii | Zestaw reguł wstrzykiwania odpowiedników | Sterowanie sposobem wymuszania terminu w silnikach tłumaczenia | Formularz reguł | 4 | Polecenie języka naturalnego lub okno konfiguracji | domyślne · własne | Zmienia sposób wstrzyknięcia terminu do zapytania silnika |

---

### 4.6. Translation Memory Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządcy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 3 |
| Izolacja domyślna | Pamięć właściwa karcie sesji; zasięg karta / projekt / zespół konfigurowalny |

Pamięć tłumaczeń jako aktywne, przeszukiwalne źródło ponownego wykorzystania materiału.

**Zawartość i pełny arsenał funkcji.**

- Przegląd i edycja par segmentów; pola metadanych (projekt, klient, data, autor, kontekst).
- **Concordance Search** — wyszukiwanie fragmentu w całej pamięci z kontekstem i wcześniejszym przekładem.
- **Bitext Aligner** — wyrównanie gotowych tekstów dwujęzycznych w pary segmentów do budowy pamięci.
- Import i eksport TMX; konserwacja pamięci: duplikaty, scalanie, masowa podmiana, filtrowanie wg pola.
- **Pre-translate/Leverage** — wstępne wypełnienie nowego materiału trafieniami z pamięci przed uruchomieniem tłumaczenia maszynowego.
- Próg dopasowania rozmytego oraz polityka dopasowania kontekstowego (101%).

**Makieta tekstowa — stan spoczynku.**

```
┌─ Translation Memory ────────────────┐
│ [ 🔍 konkordancja ]   [ Zasięg ▼ ]  │
├─────────────────────────────────────┤
│ PL „umowa ramowa”              98%  │
│ EN „framework agreement”     [ ⋮ ]  │
│ projekt: Alfa · 2026-07-14          │
├─────────────────────────────────────┤
│ PL „warunki płatności”         86%  │
│ EN „payment terms”           [ ⋮ ]  │
├─────────────────────────────────────┤
│ Próg dopasowania: 75%        [ ☰ ]  │
└─────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Para segmentów | Karta z segmentem źródłowym i przekładem | Jednostka pamięci tłumaczeń | Średnia karta z metadanymi | 1 (w otwartym oknie) | Widoczna po otwarciu okna | domyślna · edytowana · duplikat | Kliknięcie wstawia przekład do aktywnego segmentu |
| Pole konkordancji | Pole `.dn-input` z ikoną `szukaj` | Wyszukanie fragmentu w całej pamięci | Małe pole w nagłówku kolumny | 1 (w otwartym oknie) | Widoczne po otwarciu okna | puste · z wynikami | Prezentuje trafienia z kontekstem |
| Wskaźnik dopasowania | Procent podobieństwa | Ocena przydatności trafienia | Mały tekst przy parze | 1 (w otwartym oknie) | Widoczny przy pozycji | poniżej progu · powyżej progu · kontekstowe (101%) | Kliknięcie prezentuje różnicę względem segmentu bieżącego |
| Selektor zasięgu pamięci | Element zwinięty `Zasięg ▼` | Wybór zasięgu: karta, projekt, zespół | Znacznik w nagłówku kolumny | 2 | Kliknięcie znacznika; zwija się po wyborze | karta · projekt · zespół | Przełącza zbiór trafień |
| Suwak progu dopasowania | Ustawienie progu procentowego | Sterowanie czułością podpowiedzi | Mały suwak w stopce kolumny | 2 | Widoczny po otwarciu okna | domyślny · własny | Zmienia zbiór podpowiedzi w Translation Panels |
| Menu pary segmentów | Menu `⋮` | Edycja pary, usunięcie, oznaczenie duplikatu, przejście do kontekstu | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji na parze |
| Menu pamięci | Menu `☰` | Import i eksport TMX, pre-translate, wyrównanie tekstów dwujęzycznych | Jeden element zbiorczy | 3 | Kliknięcie `☰` | zwinięte · rozwinięte | Uruchamia wskazaną operację na pamięci |
| Konserwacja pamięci | Operacje wsadowe na indeksie | Czyszczenie duplikatów, scalanie pamięci, masowa podmiana | Panel wysuwany | 4 | Polecenie języka naturalnego lub okno konfiguracji | zwinięty · rozwinięty | Wykonuje operację wsadową na wskazanym zbiorze par |

---

### 4.7. Format Studio

| Aspekt | Wartość |
|---|---|
| Typologia | Okno robocze z podglądem dokumentu |
| Waga wizualna | Kolumna sąsiadująca w obszarze roboczym, otwierana przy pracy z dokumentem |
| Warstwa | 2 |
| Izolacja domyślna | Dokument wejściowy i wynikowy właściwe karcie sesji |

Tłumaczenie dokumentów z zachowaniem formatu oraz obieg zasobów lokalizacyjnych.

**Zawartość i pełny arsenał funkcji.**

- Podgląd dokumentu (DOCX, PDF, PPTX, XLSX, ODT, Markdown, HTML) w kolumnie sąsiadującej z warstwą segmentów.
- OCR skanów przed tłumaczeniem; detekcja układu, tabel i kolumn.
- Odtworzenie dokumentu wynikowego ze stylami, osadzeniami i strukturą.
- **Layout Diff / Visual Compare** — wykrycie przepełnień i złamanego układu po tłumaczeniu.
- Obsługa formatów lokalizacyjnych (XLIFF, PO, JSON, YAML, RESX, Android XML, iOS `.strings`) z ochroną kluczy i zmiennych.
- Pseudolokalizacja do testu interfejsu; reguły liczby mnogiej i rodzaju (CLDR/ICU).
- Eksport docelowy w formacie natywnym oraz eksport dwujęzyczny.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Format Studio ─────────────────────┐
│ raport-2026.docx      [ Podgląd ▼ ] │
├─────────────────────────────────────┤
│ ┌ podgląd ──────┐ ┌ segmenty ─────┐ │
│ │ Nagłówek 1    │ │[1] Nagłówek…  │ │
│ │ ─────────     │ │[2] Akapit…    │ │
│ │ Akapit tekstu │ │[3] Tabela 1.1 │ │
│ │ [ tabela ]    │ │[4] ⚠ przepeł. │ │
│ └───────────────┘ └───────────────┘ │
├─────────────────────────────────────┤
│ [ Eksport ▼ ]                [ ⋮ ]  │
└─────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Podgląd dokumentu | Renderowany widok strony dokumentu | Kontrola układu i wierności formatu | Kolumna podglądu | 1 (w otwartym oknie) | Widoczny po otwarciu okna | źródłowy · wynikowy · porównanie | Kliknięcie fragmentu przewija warstwę segmentów |
| Warstwa segmentów dokumentu | Lista segmentów powiązanych z miejscami w dokumencie | Edycja treści z zachowaniem osadzenia | Kolumna segmentów | 1 (w otwartym oknie) | Widoczna po otwarciu okna | domyślny · przepełnienie (obramowanie ostrzegawcze) | Zaznacza odpowiadające miejsce w podglądzie |
| Znacznik formatu dokumentu | Znacznik kontekstowy z nazwą pliku | Identyfikacja materiału i otwarcie okna | Znacznik w pasku kontekstu | 2 | Kliknięcie znacznika | brak dokumentu · dokument wczytany | Otwiera Format Studio na wczytanym dokumencie |
| Selektor trybu podglądu | Element zwinięty `Podgląd ▼` | Wybór widoku: źródło, wynik, porównanie układów | Znacznik w nagłówku okna | 2 | Kliknięcie `▼`; zwija się po wyborze | zwinięty · rozwinięty | Przełącza widok kolumny podglądu |
| Wyzwalacz eksportu | Element zbiorczy `Eksport ▼` | Wydanie dokumentu w formacie natywnym lub dwujęzycznym | Mały wyzwalacz w stopce okna | 2 | Kliknięcie `▼` | zwinięty · rozwinięty · generowanie | Generuje plik wynikowy do pobrania |
| Menu operacji dokumentu | Menu `⋮` | OCR skanu, ponowna detekcja układu, porównanie układów, ochrona kluczy zasobów | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Uruchamia wskazaną operację na dokumencie |
| Pseudolokalizacja | Transformacja treści do testu interfejsu | Test wydłużenia i znaków diakrytycznych przed tłumaczeniem | Operacja wywoływana | 4 | Polecenie języka naturalnego lub skrót klawiszowy | domyślna · zastosowana | Generuje pseudotłumaczenie zasobu |
| Reguły liczby mnogiej i rodzaju | Zestaw reguł CLDR/ICU | Poprawność form gramatycznych w zasobach lokalizacyjnych | Formularz reguł | 4 | Okno konfiguracji lub polecenie języka naturalnego | domyślne · własne | Stosuje reguły do zasobu docelowego |

---

### 4.8. QA & Review Center

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządcy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 3 |
| Izolacja domyślna | Profile kontroli jakości właściwe karcie sesji; udostępnienie projektowe konfigurowalne |

Kontrola jakości, korekta językowa i przebieg akceptacji w jednym miejscu.

**Zawartość i pełny arsenał funkcji.**

- Uruchomienie profili kontroli jakości na wszystkich panelach: liczby, daty, waluty, placeholdery, długość, pominięcia, spójność.
- Korekta językowa per język: gramatyka, ortografia, styl, rejestr, czytelność, typografia locale, fałszywi przyjaciele.
- **LQA Scorecard (MQM/DQF)** — ocena jakości z kategorią błędu, wagą i wynikiem końcowym; eksport raportu.
- Lista ostrzeżeń per panel z przejściem do segmentu; poprawki masowe.
- Przebieg akceptacji: tłumaczenie → korekta → zatwierdzenie, ze śladem autora i statusem.
- Pakiet dla wykonawcy zewnętrznego (XLIFF + TM + termbase + instrukcje) i odbiór zwrotu.
- Eksport zbiorczy wszystkich paneli oraz statystyka objętości z wyceną.

**Makieta tekstowa — stan spoczynku.**

```
┌─ QA & Review Center ────────────────┐
│ [ Profil QA ▼ ]        [ Uruchom ▶ ]│
├─────────────────────────────────────┤
│ Panel FR — 3 ostrzeżenia      [ ⋮ ] │
│  ⚠ [7] brak wartości liczbowej      │
│  ⚠ [11] przekroczony limit długości │
│  ⚠ [14] placeholder {name} pominięty│
├─────────────────────────────────────┤
│ Panel DE — 1 ostrzeżenie      [ ⋮ ] │
│  ⚠ [4] niespójność terminu „klient” │
├─────────────────────────────────────┤
│ Akceptacja: korekta 2/3      [ ☰ ]  │
└─────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Lista ostrzeżeń per panel | Grupa pozycji ostrzeżeń z numerem segmentu | Przegląd wyników kontroli jakości | Ciało kolumny | 1 (w otwartym oknie) | Widoczna po otwarciu okna | brak ostrzeżeń · ostrzeżenia · błąd blokujący | Kliknięcie przewija Translation Panel do segmentu |
| Wskaźnik przebiegu akceptacji | Licznik etapów przebiegu | Stan zatwierdzania materiału | Mały pasek w stopce kolumny | 1 (w otwartym oknie) | Widoczny po otwarciu okna | tłumaczenie · korekta · zatwierdzone | Kliknięcie prezentuje ślad autora zmiany |
| Przycisk „Uruchom” | Przycisk `--zloty` | Uruchomienie kontroli jakości na wszystkich panelach | Mały przycisk CTA w nagłówku kolumny | 1 (w otwartym oknie) | Widoczny po otwarciu okna | domyślny · ładowanie · wynik | Zleca zadanie kontroli do pętli wykonawczej |
| Selektor profilu QA | Element zwinięty `Profil QA ▼` | Wybór zestawu reguł kontroli | Znacznik w nagłówku kolumny | 2 | Kliknięcie `▼`; zwija się po wyborze | domyślny · profil projektu · profil klienta | Zmienia zakres reguł kolejnego przebiegu |
| Menu ostrzeżeń panelu | Menu `⋮` | Poprawki masowe, pominięcie ostrzeżenia, przejście do korekty językowej panelu | Jeden element zbiorczy | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Stosuje poprawkę do wszystkich wskazanych segmentów |
| Menu wydania | Menu `☰` | Eksport zbiorczy, pakiet dla wykonawcy zewnętrznego, odbiór zwrotu, statystyka objętości i wycena | Jeden element zbiorczy | 3 | Kliknięcie `☰` | zwinięte · rozwinięte | Generuje pakiet lub raport do pobrania |
| Karta oceny LQA (MQM/DQF) | Formularz oceny jakości z kategoriami błędów | Ocena jakości materiału z wagą i wynikiem końcowym | Panel wysuwany | 4 | Polecenie języka naturalnego lub skrót klawiszowy | pusta · wypełniona · zamknięta | Generuje raport oceny do eksportu |
| Definicje profili QA | Zestawy reguł i wag ostrzeżeń | Dostrojenie kontroli do wymagań projektu lub klienta | Formularz reguł | 4 | Okno konfiguracji | domyślne · własne | Zapisuje profil dostępny w selektorze |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — tłumaczenie równoległe od tekstu źródłowego

```
 Chat Window        Execution Loop        Source Panel /        Glossary &
 Użytkownik ↔       Koordynator ↔         Translation Panels    Termbase
 Wykonawca          Wykonawca
 ───────────        ──────────────        ──────────────────    ───────────
 1. zlecenie
    „przetłumacz
     na EN, DE, FR”
        │
        ▼
                    2. dekompozycja na
                       zadania per język
                       i zakres segmentów
                            │
                            ▼
                                          3. segmentacja źródła,
                                             tłumaczenie równoległe
                                             wszystkich kolumn
                                                   │
                                                   ▼
                    4. kontrola jakości  ◄──────────┤
                       zadania, decyzja             │
                       o ponowieniu                 ▼
                            │             5. korekta ręczna ───►  6. spójność
                            │                segmentów                terminologii,
                            ▼                                          ujednolicenie
 7. zatwierdzenie                          8. eksport plików
    wyniku                                    dwujęzycznych
```

### 5.2. Przepływ — reakcja na zmianę tekstu źródłowego

```
Source Panel — edycja segmentu [4]
        │
        ▼
Oznaczenie segmentu [4] we wszystkich Translation Panels
jako „wymaga rewizji” (bez usunięcia dotychczasowego tłumaczenia)
        │
        ▼
Execution Loop Window — Koordynator dopisuje do kolejki
zadanie ponownego tłumaczenia segmentu [4] dla każdego języka
        │
        ├──► Panel EN — ponowne tłumaczenie, do weryfikacji
        ├──► Panel DE — ponowne tłumaczenie, do weryfikacji
        └──► Panel FR — ponowne tłumaczenie, do weryfikacji
                            │
                            ▼
             korekta ręczna → oznaczenie „zweryfikowany”
```

### 5.3. Przepływ — praca dwujęzyczna z modułem Studio

```
Studio › Tools Panel                          Translate › Source Panel
─────────────────────                         ─────────────────────────
 operacja „Tłumaczenie”      ─────────────►    tekst zaznaczenia jako
 na zaznaczonym fragmencie                     tekst źródłowy
 (powiązanie Studio ◄──► Translate,                     │
  konfigurowalne — rozdz. 6.3)                          ▼
                                              Execution Loop — zadania
                                              tłumaczenia per język
                                                         │
                                                         ▼
                                              Translation Panels — wynik
                                              w wybranych językach
                                                         │
                                                         ▼
                                              powrót do Studio Editor jako
                                              wersja wielojęzyczna dokumentu
```

### 5.4. Przepływ — kontrola jakości przed eksportem zbiorczym

```
Chat Window — polecenie „uruchom kontrolę jakości wszystkich paneli”
        │
        ▼
Execution Loop Window — zadania kontroli per panel językowy
        │
        ▼
QA & Review Center — sprawdzenie: liczby, daty, placeholdery,
                     długość względem limitu, segmenty pominięte
        │
        ▼
Lista ostrzeżeń per panel ──► korekta ręczna segmentów oznaczonych
        │
        ▼
Ponowienie zadania kontroli przez Koordynatora do wyniku bez ostrzeżeń
        │
        ▼
Eksport zbiorczy wszystkich paneli (XLIFF/DOCX/TXT)
```

### 5.5. Przepływ — tłumaczenie dokumentu z zachowaniem formatu

```
Source Panel — import dokumentu (DOCX/PDF/PPTX/XLSX/ODT)
        │           (skan → OCR przed segmentacją)
        ▼
Format Studio — detekcja układu, warstwa segmentów obok podglądu
        │
        ▼
Execution Loop Window — zadania: tłumaczenie, terminologia,
                        kontrola długości, odtworzenie formatu
        │
        ▼
Layout Diff — wykrycie przepełnień i złamanego układu
        │
        ▼
Eksport docelowy w formacie natywnym oraz eksport dwujęzyczny
```

---

## 6. Stany, dane i powiązania

### 6.1. Model stanów segmentu tłumaczenia

```
┌────────────────┐  automatyczne     ┌────────────────┐   korekta ręczna   ┌────────────────┐
│ Nieprzetłuma-   │  tłumaczenie      │ Przetłumaczony  │   i akceptacja     │ Zweryfikowany   │
│ czony           │──────────────────►│ automatycznie   │───────────────────►│                 │
└────────────────┘                    └────────────────┘                    └───────┬────────┘
                                                                                       │
                                                                     zmiana tekstu źródłowego
                                                                                       ▼
                                                                              ┌────────────────┐
                                                                              │ Wymaga rewizji  │
                                                                              └────────────────┘
```

### 6.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Translate |
|---|---|
| `sesja`, `karta_sesji` | Nośnik tekstu źródłowego i zestawu aktywnych języków docelowych |
| `wiadomosc` | Historia Chat Window oraz komunikaty sterujące Execution Loop Window |
| `zadanie`, `kolejka_zadan` | Dekompozycja zlecenia tłumaczeniowego, stan zadań per język i ponowienia |
| `artefakt`, `wersja_artefaktu` | Tekst źródłowy, dokumenty wynikowe, pliki dwujęzyczne, pamięć i bazy terminologiczne jako artefakty |
| `kanal_modelu` | Silniki wykorzystywane do generowania tłumaczeń |
| `ustawienie` | Parametry modułu: próg dopasowania pamięci tłumaczeń, profile QA, reguły segmentacji, siatka stawek |
| `powiazanie_komponentu` | Jawne powiązania Studio ◄──► Translate oraz Translate ◄──► Library (Załącznik F.4 Modelu danych) |
| `profil_izolacji`, `regula_izolacji_kontekstu` | Współdzielenie kontekstu na poziomie „para modułów” przy pracy dwujęzycznej |

### 6.3. Izolacja i konfigurowalność — punkty właściwe modułowi Translate

| Punkt izolacji | Stan wyjściowy (domyślny) | Co obejmuje ustawienie konfiguracyjne | Gdzie |
|---|---|---|---|
| Para modułów Studio–Translate | Operacje kontekstowe AI Studio nie są dostępne w Translate domyślnie | Współdzielenie kontekstu na poziomie „para modułów”, zgodnie ze scenariuszem F.4 Modelu danych | Okno konfiguracji, poziom zasięgu „para modułów” |
| Baza terminologiczna | Właściwa bieżącej karcie sesji | Udostępnienie tej samej bazy kilku kartom lub projektom zespołu tłumaczeniowego | Okno konfiguracji, poziom „karta sesji” lub „projekt” |
| Pamięć tłumaczeń | Właściwa bieżącej karcie sesji | Zasięg karta / projekt / zespół, próg dopasowania, polityka dopasowania kontekstowego | Okno konfiguracji, zakres „Pamięć” |
| Historia i pamięć Chat Window | Odrębna per karta sesji | Współdzielenie między kartami prowadzącymi tłumaczenie tego samego materiału | Okno konfiguracji, poziom „karta sesji” |
| Rejestr przebiegu pętli wykonawczej | Odrębny per karta sesji | Retencja rejestru i udostępnienie zespołowi projektowemu | Okno konfiguracji, zakres „Historia” |
| Izolacja techniczna procesu sesji (8 zakresów) | Żaden zakres domyślnie nie jest aktywny | Odrębne konto i token oraz brak dostępu sieciowego przy tłumaczeniu materiałów poufnych | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

### 6.4. Powiązania z innymi modułami

```
        Studio  ◄──►  TRANSLATE  ◄──►  Library
                          ▲
                          │
                     Automations
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Studio | Współdzielenie operacji kontekstowych AI przy pracy dwujęzycznej; Translate jako ogniwo pośrednie dokumentu wielojęzycznego | Konfiguracyjne | Tools Panel (Studio) ◄──► Source Panel / Translation Panels (Translate) |
| Library | Zapis wydanych plików dwujęzycznych, pamięci tłumaczeń i baz terminologicznych jako artefaktów | Konfiguracyjne | QA & Review Center / Format Studio (Translate) → Library |
| Automations | Udostępnienie operacji modułu (tłumaczenie, kontrola jakości, eksport) jako kroków procesu wsadowego | Konfiguracyjne | Execution Loop Window (Translate) ◄──► Automations |

Moduł pozostaje w pełni funkcjonalny również bez żadnego powiązania międzymodułowego — każde powiązanie jest jawne i ustanawiane z okna konfiguracji.

---

## 7. Punkty sterowania z okna konfiguracji

Operator personalizuje moduł z okna konfiguracji (Model konfiguracji, rozdz. 5) na właściwej warstwie zasięgu: globalna → środowisko → projekt → sesja/karta. Zero blokad — brak ustawienia oznacza wartość domyślną, nie zatrzymanie pracy.

| Zakres konfiguracji (Model konfiguracji) | Co personalizuje Operator w module Translate |
|---|---|
| 5.4 Zachowanie modeli | Silniki tłumaczenia i ich kolejność; **klucze API jawne**; domyślny silnik per para językowa; parametry MT (adaptacyjność, temperatura) |
| 5.7 Rozszerzenia | Włączenie konektorów: korekta językowa, OCR, silniki TTS/STT, dostawcy MT; źródło Danaco Plugin lub Personal |
| 5.8 Integracje | Powiązanie pary modułów **Studio ◄──► Translate**; **Translate ◄──► Library** (zapis artefaktów); udostępnienie operacji do **Automations**; webhooki i konektory zewnętrzne |
| 5.10 Pamięć + własne | Zasięg **pamięci tłumaczeń**: karta / projekt / zespół; próg dopasowania rozmytego; polityka dopasowania kontekstowego (101%); włączenie pre-translate |
| 5.12 Karty sesji | Współdzielenie **bazy terminologicznej** między kartami i projektami; migawka układu kolumn ośmiu okien; przypisanie projektu tłumaczeniowego do karty |
| 5.9 Historia | Współdzielenie historii Chat Window między kartami tłumaczącymi ten sam materiał; retencja rejestru przebiegu pętli wykonawczej |
| 5.11 Izolacja | Osiem zakresów izolacji technicznej dla materiałów poufnych (odrębne konto i token, brak dostępu sieciowego przy tłumaczeniu materiałów objętych NDA) |
| Pętla wykonawcza (5.1/5.6) | Liczba zadań realizowanych równolegle w Execution Loop Window; polityka ponowienia zadania po niepowodzeniu kontroli jakości; próg jakości wymagany do zamknięcia zadania; zakres autonomii Wykonawcy przy korekcie segmentu; zakres zatwierdzeń wymaganych od Użytkownika |
| Presety modułu (5.1/5.6) | Zestaw domyślnych języków docelowych; presety tonu i rejestru; profile QA (co sprawdzać, wagi ostrzeżeń); listy DNT i terminów zabronionych; reguły segmentacji SRX; reguły typografii locale; siatka stawek do wyceny; profile limitów długości (napisy CPS, interfejs) |

Zasada kluczy jawnych: dane dostępowe silników zewnętrznych (tłumaczenie maszynowe, OCR, TTS, STT) są widoczne i edytowalne w zakresie 5.4 „Dane dostępowe kanału”, przechowywane poza bazą danych zgodnie z Modelem danych (rozdz. 8.1) — nigdy nie są ukryte ani zaszyte na stałe.

---

## 8. Scenariusze użycia

**Scenariusz 1 — lokalizacja materiału marketingowego na trzy rynki.**
Zespół lokalizacyjny wkleja tekst kampanii do Source Panel i zleca w Chat Window tłumaczenie na EN, DE i FR. Execution Loop Window prezentuje dekompozycję zlecenia na zadania per język, a Koordynator nadzoruje ich realizację. Zespół ustawia ton „marketingowy” dla każdej kolumny niezależnie, a Glossary & Termbase Manager wymusza spójne tłumaczenie nazwy produktu we wszystkich trzech wersjach. Po korekcie ręcznej zespół eksportuje trzy pliki DOCX gotowe do dalszej dystrybucji.

**Scenariusz 2 — reakcja na zmianę tekstu źródłowego w trakcie pracy.**
Klient prosi o zmianę jednego akapitu. Redaktor edytuje segment w Source Panel — wszystkie aktywne Translation Panels oznaczają odpowiadający segment jako „wymaga rewizji”, a Koordynator dopisuje do kolejki zadanie ponownego tłumaczenia tego segmentu dla każdego języka. Tłumacze poprawiają wynik niezależnie w swoich kolumnach.

**Scenariusz 3 — praca dwujęzyczna zainicjowana w Studio.**
Redaktor dokumentu w module Studio zaznacza akapit i wybiera operację „Tłumaczenie” z Tools Panel. Dzięki skonfigurowanemu powiązaniu na poziomie pary modułów tekst trafia do Source Panel modułu Translate, a wynik po zaakceptowaniu wraca do Studio Editor jako wersja dwujęzyczna dokumentu.

**Scenariusz 4 — wykrycie niespójności terminologicznej.**
Wykonawca sygnalizuje w Chat Window, że termin „klient” występuje w wersji niemieckiej w dwóch różnych formach. Użytkownik otwiera Glossary & Termbase Manager, przegląda oba wystąpienia i jednym kliknięciem ujednolica odpowiednik we wszystkich segmentach kolumny DE.

**Scenariusz 5 — kontrola jakości przed wysyłką do klienta.**
Przed eksportem zespół poleceniem w Chat Window uruchamia kontrolę jakości wszystkich paneli. Execution Loop Window prowadzi zadania kontroli per panel, a QA & Review Center wskazuje segment w kolumnie FR, w którym pominięto wartość liczbową obecną w tekście źródłowym. Po korekcie Koordynator ponawia zadanie kontroli, a zespół eksportuje zbiorczo wszystkie panele w formacie XLIFF.

**Scenariusz 6 — tłumaczenie umowy z zachowaniem formatu.**
Prawnik importuje umowę w formacie DOCX do Source Panel. Format Studio prezentuje podgląd dokumentu obok warstwy segmentów, a pętla wykonawcza prowadzi zadania tłumaczenia, wymuszenia terminologii prawniczej i kontroli długości. Layout Diff wskazuje przepełnienie w tabeli, po korekcie moduł wydaje dokument wynikowy ze stylami i tabelami oraz plik dwujęzyczny.

**Scenariusz 7 — lokalizacja zasobów interfejsu produktu.**
Zespół produktowy wczytuje pliki JSON i Android XML do Format Studio. Klucze i zmienne pozostają chronione, reguły liczby mnogiej CLDR stosowane są per język, a pseudolokalizacja pozwala sprawdzić wydłużenie tekstu w interfejsie przed właściwym tłumaczeniem. Kontrola długości sygnalizuje przekroczenia limitów, po czym zasoby wracają w formatach natywnych.

**Scenariusz 8 — przekazanie materiału wykonawcy zewnętrznemu.**
Kierownik projektu otwiera QA & Review Center i generuje pakiet zawierający XLIFF, pamięć tłumaczeń, bazę terminologiczną i instrukcje. Po zwrocie pakietu moduł wczytuje tłumaczenia, uruchamia profil kontroli jakości i prowadzi przebieg akceptacji: tłumaczenie → korekta → zatwierdzenie, ze śladem autora każdej zmiany.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + Shift + L` | Dodanie nowego języka docelowego | Translation Panels |
| `Alt + ↑` / `Alt + ↓` | Przejście do poprzedniego/następnego segmentu wymagającego rewizji | Translation Panels |
| `Ctrl/Cmd + G` | Otwarcie Glossary & Termbase Manager | dowolne okno modułu |
| `Ctrl/Cmd + M` | Otwarcie Translation Memory Panel | dowolne okno modułu |
| `Ctrl/Cmd + Shift + Q` | Otwarcie QA & Review Center | dowolne okno modułu |
| `Ctrl/Cmd + Shift + E` | Wstrzymanie i wznowienie pętli wykonawczej | Execution Loop Window |
| `Ctrl/Cmd + Shift + R` | Otwarcie rejestru przebiegu pętli wykonawczej | Execution Loop Window |
| `Ctrl/Cmd + K` | Wyszukiwarka funkcji modułu | dowolne okno modułu |
| Ikona `globus` | Wybór/zmiana języka | Source Panel, Translation Panels |
| Ikona `ostrzezenie` | Segment wymaga rewizji / niespójność terminologii | Translation Panels, QA & Review Center |
| Ikona `ptaszek` | Segment zweryfikowany · zadanie ukończone | Translation Panels, Execution Loop Window |
| Ikona `petla` | Zadanie ponowione w pętli wykonawczej | Execution Loop Window |
| Ikona `klodka` | Termin oznaczony „nie tłumacz” | Glossary & Termbase Manager |
| Ikona `pobierz` | Eksport panelu, pamięci, bazy terminów lub dokumentu | Translation Panels, Translation Memory Panel, Format Studio, QA & Review Center |
| Ikona `szukaj` | Wyszukiwanie terminu i konkordancja | Glossary & Termbase Manager, Translation Memory Panel |
| Znacznik `⋮` | Rozwinięcie kontekstowe zestawu akcji | wszystkie okna modułu |
| Znacznik `☰` | Rozwinięcie zbiorczego menu operacji | Chat Window, Glossary & Termbase Manager, Translation Memory Panel, QA & Review Center |
| Znacznik `▼` | Rozwinięcie selektora warstwy 2 | wszystkie okna modułu |

*Koniec dokumentu. Moduł Translate — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
