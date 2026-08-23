# Moduł Browser — dokumentacja projektowa

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
| **Tytuł** | Moduł Browser — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować) |
| **Środowiska dostępności** | TalkIn, WorkSpace |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł nie tworzy komponentu własnego |
| **Stos techniczny odniesienia** | Rdzeń serwera: Go (orkiestracja, procesy sesji). Silnik prezentacji: Chromium / WebView2 / WKWebView zależnie od platformy |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 4, 9.4, 12, Załącznik A) · Specyfikacja modułów (rozdz. 4.4) · Specyfikacja okien operacyjnych (rozdz. 5, 6.4, 10) · System wizualny (rozdz. 6, 8, 10) · Model danych (rozdz. 8, 16, 17, Załącznik F.2) · Izolacja i zależności (rozdz. 1–6) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie; klucze jawne; domyślne zachowanie modułu = wykonanie |

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

Browser jest modułem współdzielonego przeglądania internetu — przestrzenią, w której użytkownik i Wykonawca patrzą na tę samą stronę w tym samym czasie i pracują nad tą samą treścią równolegle, zamiast wymieniać się linkami i wklejanymi fragmentami tekstu. Wspólny podgląd jest istotą modułu: to, co widzi Wykonawca, jest dokładnie tym, co widzi użytkownik, bez pośredniczącego opisu.

Moduł jest pełną, osadzoną przeglądarką klasy desktopowej ze współdzielonym podglądem Wykonawcy oraz warstwą automatyzacji web. Pokrywa cztery obszary pracy:

| Obszar | Zakres w module Browser |
|---|---|
| Przeglądanie i nawigacja | Karty, grupy kart, przestrzenie robocze, zakładki, historia, wyszukiwarki, profile, kontenery tożsamości |
| Przechwytywanie treści | Zrzuty (widoczny obszar / cała strona / region / element), archiwum stron (MHTML/WARC/single-file HTML), wycinki, czytnik, tłumaczenie, synteza mowy |
| Ekstrakcja i strukturyzacja | Wyodrębnianie tabel, list, danych kontaktowych, cen, schematów `JSON-LD`/`microdata`; zbieranie danych z wielu podstron; eksport do CSV/JSON/XLSX/Markdown |
| Automatyzacja web (RPA) | Nagrywanie i odtwarzanie makr, scenariusze warunkowe, wypełnianie formularzy, harmonogramy, śledzenie zmian stron, Wykonawca operujący na stronie |

### 1.2. Granica tematyczna modułu

| Poza granicą modułu | Gdzie to należy |
|---|---|
| Głęboka, wielowątkowa synteza badawcza z wielu źródeł | Moduł **Research** (Browser zasila go źródłami i notatkami) |
| Trwałe katalogowanie i tagowanie repozytorium materiałów | Moduł **Library** (Browser przekazuje artefakty) |
| Silnik tłumaczenia jako usługa | Moduł **Translate** (Browser wywołuje go, nie implementuje) |
| Redakcja i obróbka dokumentów wyodrębnionych ze stron | Moduł **Studio** (Browser dostarcza wsad) |
| Pełne środowisko programistyczne i repozytoria | Moduł **Developer** / **CodeStudio** (Browser udostępnia narzędzia inspekcyjne, nie środowisko programistyczne) |

Granica jest wykonawcza, nie licencyjna: powiązania są jawne i konfigurowalne.

### 1.3. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Browser |
|---|---|
| Analitycy i researcherzy | Wspólna analiza stron konkurencji, ofert, dokumentacji technicznej online |
| Zespoły zakupowe | Porównanie ofert wielu dostawców w czasie rzeczywistym z asystą Wykonawcy |
| Dziennikarze i weryfikatorzy faktów | Wspólna weryfikacja informacji na źródłowej stronie internetowej |
| Zespoły badawcze (Research) | Zbieranie źródeł internetowych zasilających dalszą pracę badawczą |
| Zespoły operacyjne | Powtarzalne przebiegi przeglądania, wypełniania formularzy i pobierania danych |
| Użytkownicy indywidualni | Szybkie zrozumienie długiej strony bez samodzielnego jej czytania w całości |

### 1.4. Po co — wartość modułu

| Problem pracy bez wspólnego podglądu | Rozwiązanie w Browser |
|---|---|
| Konieczność opisywania Wykonawcy treści strony słowami | Browser Window daje Wykonawcy ten sam widok co użytkownikowi |
| Utrata śladu, skąd pochodzi dana informacja | Sources Panel gromadzi adresy i tytuły odwiedzonych stron automatycznie |
| Istotne fragmenty giną w gąszczu przeglądanych stron | Notes Panel wiąże notatkę bezpośrednio ze źródłem, z którego pochodzi |
| Powrót do materiału wymaga ponownego przeszukania internetu | Zebrane źródła i notatki pozostają dostępne w module i przekazywalne dalej |
| Powtarzalne czynności na stronach wykonywane ręcznie | Automation Studio prowadzi makra i scenariusze pod nadzorem pętli wykonawczej |

### 1.5. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn albo WorkSpace
        ▼
ŚRODOWISKO (TalkIn / WorkSpace) ── boczna nawigacja modułów
        │
        ▼
MODUŁ: BROWSER ──────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Browser Window ·
  Sources Panel · Notes Panel · Automation Studio ·
  Capture & Monitor Panel
        │
        ▼
KARTA SESJI — wspólny podgląd i lista źródeł bieżącej sesji
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 6.3)
```

### 1.6. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | TalkIn, WorkSpace |
| Środowisko CodeStudio | Moduł niedostępny |
| Komponent własny | Browser nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 7 (łącznie z Chat Window i Execution Loop Window) |
| Charakter podglądu | Współdzielony — użytkownik i Wykonawca operują na tym samym, aktualnym stanie strony |

---

## 2. Katalog funkcji i narzędzi

Każda pozycja: nazwa, działanie oraz zależności (biblioteki, formaty, integracje). Nazwy bibliotek Go wskazują realne rozwiązania implementacyjne rdzenia.

### 2.1. Nawigacja, karty i przestrzenie robocze

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Karty i grupy kart** | Wiele stron równolegle w oknie; grupowanie kart w nazwane, kolorowane zestawy zwijane jednym kliknięciem | WebView2 / Chromium multi-tab; stan w `karta_sesji` |
| **Przestrzenie robocze** | Odrębne zestawy kart, zakładek i historii per temat pracy, przełączane bez utraty stanu | Trwały stan sesji (Go), `profil_izolacji` |
| **Drzewo kart / karty pionowe** | Boczna, hierarchiczna lista kart z zagnieżdżeniem otwarć | Komponent listy `.dn-*`; graf otwarć w stanie karty |
| **Zawieszanie kart** | Usypianie nieużywanych kart dla oszczędności pamięci, wznowienie kliknięciem | Zarządzanie procesami sesji rdzenia (Go), zdarzenia widoczności |
| **Sesje i przywracanie** | Zapis pełnego zestawu otwartych kart jako nazwanej sesji; przywrócenie po restarcie serwera | Trwałość stanu (rdzeń Go), format zrzutu sesji `JSON` |
| **Kontenery tożsamości** | Izolowane profile ciasteczek i magazynu per karta (dwa konta tego samego serwisu jednocześnie) | Chromium partitioned storage; `profil_izolacji`, `regula_izolacji_kontekstu` |
| **Multi-wyszukiwarka i słowa-klucze** | Pasek adresu z konfigurowalnymi silnikami (`!g`, `!yt`) i wyszukiwaniem wewnątrz platformy | Konfiguracja silników; integracja z modułem Research |
| **Gesty myszy i szybkie polecenia** | Nawigacja gestami oraz szybkie polecenia paska adresu wywoływane znakiem `/` | Warstwa skrótów klienta; paleta poleceń `.dn-paleta` (`Ctrl/Cmd + K`) |

### 2.2. Zakładki, historia i sesje

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Menedżer zakładek z tagami** | Zakładki z folderami, tagami, notatką i pełnotekstowym wyszukiwaniem | Indeks pełnotekstowy (`blevesearch/bleve`); eksport `Netscape bookmarks HTML`, `JSON` |
| **Historia z osią czasu i wyszukiwaniem treści** | Historia przeszukiwalna po treści odwiedzonych stron, nie tylko po tytule i adresie | Indeks treści (Bleve); migawki tekstu strony |
| **Kolejka czytania** | Odłożenie strony do przeczytania wraz z przypomnieniem | Kolejka w stanie sesji; integracja z Capture & Monitor Panel |
| **Import i eksport profilu** | Przeniesienie zakładek, referencji do sekretów i historii między przeglądarkami | Parsery `Netscape HTML`, `CSV`, `JSON`; format `sqlite` historii |

### 2.3. Przechwytywanie treści i archiwizacja

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Zrzut ekranu — cztery tryby** | Widoczny obszar, cała przewijana strona, wybrany region, pojedynczy element DOM | CDP `Page.captureScreenshot` (`chromedp/chromedp`); wyjście `PNG`/`JPEG`/`WebP` |
| **Adnotacja zrzutu** | Strzałki, wyróżnienia, rozmycie danych wrażliwych, numeracja kroków | Canvas klienta; zapis jako `artefakt` (Model danych, Zał. F.2) |
| **Archiwum strony (single-file)** | Zapis kompletnej strony z zasobami do jednego pliku wiernego oryginałowi | `MHTML`, `single-file HTML` (data-URI inlining), `WARC`; CDP snapshot |
| **Wycinek regionu** | Zaznaczenie fragmentu strony i zapis jako czysty Markdown lub HTML z zachowaniem źródła | Konwersja HTML→Markdown (`JohannesKaufmann/html-to-markdown`); trafia do Notes Panel |
| **Zapis strony do PDF** | Renderowanie bieżącej strony lub trybu czytnika do PDF | CDP `Page.printToPDF`; `pdfcpu/pdfcpu` do scalania wielu stron |
| **Nagrywanie ekranu i GIF** | Zapis sekwencji interakcji na stronie jako wideo lub GIF | CDP screencast frames; enkoder `GIF`/`WebM` (Go `image/gif`, `ffmpeg`) |
| **Migawka tekstowa dla Wykonawcy** | Zwięzły, oczyszczony zrzut treści strony podawany Wykonawcy jako kontekst | `go-shiori/go-readability`; `accessibility tree` z CDP |

### 2.4. Czytnik, dostępność i tłumaczenie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Tryb czytnika** | Oczyszczenie strony z reklam i nawigacji do czystego tekstu, sterowanie typografią, tłem, szerokością kolumny | `go-shiori/go-readability` |
| **Czytnik bez sieci** | Zapis oczyszczonych artykułów do czytania bez połączenia, biblioteka artykułów | Zapis oczyszczonego HTML jako `artefakt`; indeks Bleve |
| **Odczyt na głos** | Czytanie treści strony syntezą mowy z podświetlaniem czytanego fragmentu | Web Speech API klienta lub kanał modelu mowy; sterowanie tempem i głosem |
| **Tłumaczenie strony w miejscu** | Tłumaczenie widocznej treści z zachowaniem układu, na żądanie lub automatyczne per domena | Integracja z modułem **Translate**; wykrywanie języka (`pemistahl/lingua-go`) |
| **Tryby dostępności** | Nakładki: kontrast, krój dla dyslektyków, powiększenie, maska czytania | Wstrzykiwane arkusze CSS klienta; profil w konfiguracji |
| **Wymuszony tryb ciemny** | Ciemny motyw dla stron bez własnego wsparcia | Filtr CSS i inwersja per domena; ustawienie w konfiguracji |

### 2.5. Ekstrakcja i strukturyzacja danych

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Wyodrębnianie tabel** | Wykrycie tabel na stronie i eksport do CSV/XLSX/JSON/Markdown, z łączeniem tabel z wielu podstron | `PuerkitoBio/goquery`; `encoding/csv`, `qax-os/excelize` |
| **Ekstrakcja danych strukturalnych** | Odczyt `JSON-LD`, `microdata`, `Open Graph`, RDFa ze strony | Parsery `JSON-LD`; goquery do `<meta>` i `<script type=ld+json>` |
| **Zbieranie danych wskazaniem** | Wskazanie elementów na stronie i reguła zbierania z wielu podstron oraz paginacji | Silnik selektorów CSS/XPath; `gocolly/colly`; wynik `JSON`/`CSV` |
| **Zbieranie danych kontaktowych** | Wykrycie adresów e-mail, telefonów, adresów i profili społecznościowych na stronie lub domenie | Wyrażenia regularne z walidacją; eksport `vCard`/`CSV` |
| **Ekstrakcja mediów** | Lista i pobranie obrazów, wideo oraz plików ze strony z filtrem rozmiaru i typu | CDP wykaz zasobów; menedżer pobrań (rozdz. 2.9) |
| **Porównanie struktur danych** | Zestawienie wyodrębnionych danych z dwóch stron lub dwóch momentów w tabeli różnic | Diff strukturalny; `sergi/go-diff` |
| **Ekstrakcja prowadzona językiem naturalnym** | Polecenie opisane słowami („zbierz nazwy i ceny") zwracane jako JSON zgodny ze schematem | `kanal_modelu`; walidacja `JSON Schema` |

### 2.6. Automatyzacja web (RPA)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Nagrywarka makr** | Rejestracja kliknięć, wpisów, przewinięć i nawigacji jako powtarzalnego makra | CDP event recording; zapis kroków jako `JSON` scenariusza |
| **Odtwarzanie i harmonogram makr** | Uruchamianie makra ręcznie, cyklicznie (cron) lub po zdarzeniu | Harmonogram rdzenia (Go, cron); integracja z modułem **Automations** |
| **Scenariusze warunkowe** | Kroki z warunkami, pętlami po liście, obsługą błędów i ponowieniami | Silnik przepływu; opis scenariusza w `JSON`/`YAML` |
| **Autowypełnianie formularzy** | Zapisane profile danych wstawiane do formularzy; masowe wypełnianie z listy `CSV` | Mapowanie pól; wsad `CSV`; klucze jawne w konfiguracji |
| **Wykonawca operujący na stronie** | Wykonawca realizuje zadanie opisane słownie, działając na drzewie DOM | CDP sterowanie (`chromedp/chromedp`, `go-rod/rod`); pętla obserwacja→akcja z `kanal_modelu` |
| **Testy i asercje strony** | Sprawdzenie, że strona zawiera lub nie zawiera elementu — nadzór poprawności | `go-rod/rod`, `playwright-community/playwright-go`; raport asercji |
| **Zarządzanie kolejką zadań** | Wiele scenariuszy w kolejce z priorytetami, statusem i logiem wykonania | Kolejka rdzenia; log jako `artefakt`; przekazanie do **Automations** |

Pętlą wykonawczą tych zadań steruje Koordynator, a jej przebieg prezentuje Execution Loop Window (rozdz. 4.2).

### 2.7. Współpraca Użytkownik–Wykonawca (rdzeń modułu)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Wspólny podgląd na żywo** | Wykonawca widzi dokładnie to, co użytkownik, łącznie ze zmianami po interakcji | CDP DOM i accessibility snapshot strumieniowany kanałem WebSocket |
| **Wskazywanie i cytowanie fragmentu** | Zaznaczenie fragmentu jako przedmiotu polecenia; odpowiedź odsyła do miejsca na stronie | Mapowanie zaznaczenia na węzeł DOM; kotwice cytowań |
| **Pasek pływający operacji** | Nad zaznaczeniem: Wyjaśnij / Wyodrębnij / Notatka / Tłumacz / Zapytaj | Kontekstowy pasek `.dn-*`; kanał modelu |
| **Podsumowanie strony i pytania o treść** | „Streść", „wypunktuj wnioski", „odpowiedz na pytanie z tej strony" | go-readability jako wsad; `kanal_modelu`; strumień WebSocket |
| **Rozmowa obejmująca wiele kart** | Zapytanie obejmujące treść kilku otwartych kart jednocześnie | Agregacja migawek kart; widok podzielony |
| **Wskaźnik obecności Wykonawcy** | Sygnał, że Wykonawca aktywnie odbiera bieżący widok | Plakietka `.dn-plakietka--info`; stan strumienia |

### 2.8. Prywatność, bezpieczeństwo i higiena przeglądania

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Blokada reklam i mechanizmów śledzących** | Filtrowanie reklam, skryptów śledzących i banerów według list reguł | Silnik reguł `EasyList`/`ABP filter`; CDP `Network` interception; listy aktualizowalne |
| **Automatyczne odrzucanie ciasteczek** | Wybór najbardziej prywatnego ustawienia na banerach zgód | Reguły dopasowania banerów; polityka prywatności platformy |
| **Tryb efemeryczny** | Sesja bez zapisu historii, ciasteczek i magazynu; czyszczenie przy zamknięciu | Chromium ephemeral partition; `profil_izolacji` |
| **Izolacja sieciowa sesji** | Wydzielony, ograniczony dostęp sieciowy dla przeglądania materiałów wrażliwych | Proces sesji rdzenia (Go); panel macierzy izolacji |
| **Zarządzanie uprawnieniami stron** | Per domena: kamera, mikrofon, lokalizacja, powiadomienia, JavaScript, autoodtwarzanie | Chromium permission API; polityka per domena w konfiguracji |
| **Referencje sekretów (klucze jawne)** | Odwołania do sekretów logowania trzymanych w warstwie sekretów platformy | Warstwa sekretów platformy; jawne klucze konfiguracji |

> **Zasada bezpieczeństwa.** Moduł operuje na referencjach do sekretów (klucze jawne, konfigurowalne), a faktyczne wprowadzanie haseł, danych kart i tożsamości pozostaje po stronie użytkownika lub dedykowanego menedżera. Wykonawca nie realizuje samodzielnie płatności ani przelewów.

### 2.9. Pobrania, media i pliki

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Menedżer pobrań** | Kolejka, wznawianie, pobieranie wielowątkowe, kategorie, reguły folderów | HTTP range requests (Go `net/http`); kolejka rdzenia |
| **Podgląd i konwersja obrazów** | Podgląd, kadrowanie i konwersja formatów pobranych obrazów bez opuszczania modułu | Go `image`, `disintegration/imaging`; formaty `PNG/JPEG/WebP/AVIF` |
| **Odczyt schowka i szybkie otwarcie adresu** | Wykrycie adresu w schowku i wskazanie go do otwarcia; wklejanie obrazów jako źródeł | Schowek klienta; walidacja adresu |

### 2.10. Dokumenty online w oknie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Podgląd PDF** | Otwarcie PDF w oknie z wyszukiwaniem, zakładkami, adnotacją i wyodrębnianiem tekstu | PDF.js w webview lub serwerowo `gen2brain/go-fitz` (MuPDF), `pdfcpu/pdfcpu` |
| **Ekstrakcja tekstu z PDF** | Wydobycie tekstu i tabel z PDF do dalszej pracy Wykonawcy | `ledongthuc/pdf`, `go-fitz`; tabele wykrywane heurystykami układu |
| **Rozpoznawanie tekstu (OCR)** | Rozpoznanie tekstu ze skanowanych PDF i obrazów na stronie | `otiai10/gosseract` (Tesseract); języki `pl+eng` |
| **Podgląd DOCX/XLSX/PPTX** | Otwarcie dokumentów biurowych linkowanych na stronach bez zewnętrznego programu | `unidoc/unioffice` lub renderowanie do HTML |

### 2.11. Narzędzia inspekcyjne

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Inspektor DOM i stylów** | Podgląd drzewa elementów, stylów i atrybutów; wskazanie elementu na stronie | CDP domeny `DOM` i `CSS` |
| **Monitor sieci** | Rejestr żądań i odpowiedzi, nagłówków oraz czasów; eksport `HAR` | CDP `Network`; format `HAR` |
| **Konsola i log błędów** | Odczyt komunikatów `console.*` oraz błędów strony | CDP `Runtime` i `Log` |
| **Emulacja urządzeń** | Podgląd strony w rozdzielczościach mobilnych i tabletowych, orientacja, gęstość pikseli | CDP `Emulation.setDeviceMetricsOverride` |
| **Podgląd źródła i różnic** | Widok źródła HTML oraz porównanie dwóch wersji strony w czasie | `sergi/go-diff`; podświetlanie składni |

### 2.12. Monitorowanie i kanały

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Śledzenie zmian strony** | Cykliczne sprawdzanie wybranego fragmentu i alert przy zmianie (cena, dostępność, treść) | Harmonogram (Go cron); diff (`sergi/go-diff`); powiadomienia platformy |
| **Czytnik RSS/Atom** | Subskrypcja kanałów, agregacja, oznaczanie przeczytanych, wykrycie kanału na stronie | `mmcdole/gofeed`; formaty `RSS`/`Atom`/`JSON Feed` |
| **Alerty i powiadomienia** | Reguły powiadomień (e-mail, push, w aplikacji) dla monitorów i scenariuszy | Kanał powiadomień platformy; integracja z **Automations** |
| **Kolektor źródeł do Research** | Przekazywanie monitorowanych i odwiedzonych stron jako źródeł badawczych | `powiazanie_komponentu` Browser → Research |

### 2.13. Integracje międzymodułowe

| Integracja | Co daje | Kierunek | Zależności |
|---|---|---|---|
| **Browser → Research** | Zebrane źródła i notatki zasilają Sources Manager oraz Findings Panel | jednokierunkowe, konfigurowalne | `powiazanie_komponentu`; Sources Panel → Sources Manager |
| **Browser → Library** | Zrzuty, archiwa i wycinki trafiają do trwałego repozytorium | jednokierunkowe, konfigurowalne | Notes Panel → Library Explorer; `artefakt`, `etykieta_artefaktu` |
| **Browser ↔ Translate** | Tłumaczenie stron i wycinków przez silnik Translate | dwukierunkowe wywołania | Kanał Translate; wykrywanie języka |
| **Browser → Studio** | Wyodrębnione tabele i teksty jako wsad do redakcji dokumentów | jednokierunkowe | Format `Markdown`/`CSV`/`XLSX`; artefakt |
| **Browser → Automations** | Makra i scenariusze jako kroki większych automatyzacji z harmonogramem | jednokierunkowe | Scenariusz `JSON`/`YAML`; harmonogram rdzenia |
| **Browser → Agents** | Udostępnienie zadań przeglądania agentom platformy | jednokierunkowe | Pętla wykonawcza; CDP sterowanie |

### 2.14. Zależności krytyczne — zestawienie

| Obszar | Kluczowa zależność | Rozwiązania Go i formaty |
|---|---|---|
| Sterowanie i przechwytywanie | Chrome DevTools Protocol | `chromedp/chromedp`, `go-rod/rod`, `playwright-community/playwright-go` |
| Czytnik i oczyszczanie | Readability | `go-shiori/go-readability` |
| Parsowanie HTML i ekstrakcja | Selektory DOM | `PuerkitoBio/goquery`, `gocolly/colly` |
| PDF | Render i tekst | `gen2brain/go-fitz` (MuPDF), `pdfcpu/pdfcpu`, `ledongthuc/pdf` |
| Rozpoznawanie tekstu | Tesseract | `otiai10/gosseract` (języki `pl+eng`) |
| Dokumenty biurowe | DOCX/XLSX/PPTX | `unidoc/unioffice`, eksport XLSX `qax-os/excelize` |
| Różnice | Tekst i struktury | `sergi/go-diff` |
| Indeks i wyszukiwanie | Pełny tekst | `blevesearch/bleve` |
| RSS/Atom | Kanały | `mmcdole/gofeed` |
| Wykrywanie języka | Język strony | `pemistahl/lingua-go` |
| Konwersja treści | HTML→Markdown | `JohannesKaufmann/html-to-markdown` |
| Formaty archiwum | Wierne zapisy | `MHTML`, `WARC`, `single-file HTML`, `HAR` |

---

## 3. Komplet okien operacyjnych modułu

Moduł rozmieszcza okna wyłącznie w układzie pionowym (podział lewa–prawa). Kolumny sąsiadują poziomo, regulacji podlega wyłącznie ich szerokość.

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Centralny punkt pracy; polecenia w języku naturalnym, strumień odpowiedzi, zatwierdzanie i przerywanie działań |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana | 1 | Widoczne po podjęciu zlecenia; wywołanie znacznikiem `Pętla ▼` w pasku kontekstu | Pętla wykonawcza zadań przeglądania, ekstrakcji i automatyzacji web |
| 3 | Browser Window | Okno główne — osadzona przeglądarka | Prawa kolumna, dominująca | 1 | Widoczne bez interakcji | Współdzielony podgląd strony |
| 4 | Sources Panel | Okno list i źródeł | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik `Źródła ▼` w pasku kontekstu | Lista źródeł z sesji |
| 5 | Notes Panel | Zarządzanie źródłami i ustaleniami | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik `Notatki ▼` w pasku kontekstu | Notatki powiązane ze źródłem |
| 6 | Automation Studio | Okno operacji i scenariuszy | Kolumna boczna, otwierana jako rozszerzenie boczne w pełnej wysokości | 3 | Menu `Operacje ▼` w pasku kontekstu; polecenie języka naturalnego w Chat Window | Makra, scenariusze warunkowe, harmonogramy, profile autowypełniania |
| 7 | Capture & Monitor Panel | Okno przechwytywania i monitorowania | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu `Operacje ▼`; menu kebab (⋮) pozycji źródła | Zrzuty i archiwa, pobrania, monitory zmian, kanały RSS/Atom |

```
 Makieta zbiorcza — Moduł Browser         Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window        │ Browser Window          │ Kolumna
  nawigacja │ Użytkownik ↔       │ współdzielony podgląd   │ boczna
  modułów   │ Wykonawca          │ strony (użytkownik      │ (rozszerzenie
  (poza     │                    │ + Wykonawca, ta sama    │  boczne,
  zakresem  │ ───────────────    │ treść)                  │  warstwa 2–3)
  dokumentu)│ Execution Loop     │                         │
            │ Koordynator ↔      │                         │
            │ Wykonawca          │                         │
            │                    │                         │
  ─────────────────────────────────────────────────────────────────────
  [Browser] [TalkIn] [Karta sesji] [Model ▼]   Źródła ▼  Notatki ▼  ⋮  ☰
 ═══════════════════════════════════════════════════════════════════════
```

### 3.1. Warstwy widoczności w module

Moduł stosuje regułę stopniowego ujawniania funkcjonalności: element interfejsu nieprzydatny w bieżącym zadaniu przeglądania pozostaje niewidoczny. W stanie spoczynku widoczne są wyłącznie Chat Window, Execution Loop Window, Browser Window oraz pasek kontekstu ze zwiniętymi wyzwalaczami.

| Warstwa | Zastosowanie w module Browser | Sposób wywołania |
|---|---|---|
| 1 | Chat Window, Execution Loop Window, obszar renderowania strony, pasek adresu, rząd kart, pasek kontekstu, wskaźnik obecności Wykonawcy | Widoczne bez interakcji |
| 2 | Sources Panel, Notes Panel, selektor kanału modelu, wybór przestrzeni roboczej, tryb czytnika, wybór trybu podglądu | Znacznik kontekstowy, ikona, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Automation Studio, Capture & Monitor Panel, zestaw operacji na stronie (zrzut, archiwum, ekstrakcja, tłumaczenie, podział widoku), akcje pozycji źródła i notatki | Menu `Operacje ▼`, menu kebab (⋮), menu hamburger (☰), pasek pływający zaznaczenia |
| 4 | Narzędzia inspekcyjne (DOM, sieć, konsola, emulacja urządzeń), macierz izolacji sesji, edytor scenariusza w widoku `JSON`/`YAML`, limity kroków Wykonawcy | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny. Użytkownik podstawowy nie widzi tych elementów |

Każda ukryta funkcja modułu jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window.

---

## 4. Specyfikacja okien operacyjnych

### 4.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — główne okno komunikacji Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Warstwa widoczności | 1 |
| Rola | Centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami modułu |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 6.3) |

**Zawartość i pełny arsenał funkcji.**

- Polecenia w języku naturalnym dotyczące treści aktualnie wyświetlanej strony: „streść tę stronę", „wyodrębnij tabelę cen", „porównaj tę ofertę z poprzednio odwiedzoną".
- Strumień odpowiedzi Wykonawcy na żywo kanałem WebSocket, z odwołaniami do konkretnych fragmentów widocznej strony.
- Zatwierdzanie i przerywanie działań Wykonawcy, w tym przebiegów automatyzacji uruchomionych w module.
- Polecenia nawigacyjne: „znajdź na tej stronie sekcję dotyczącą cennika", „przejdź do kolejnej strony wyników".
- Wyszukiwanie powiązanych informacji w internecie w kontekście przeglądanej strony, bez opuszczania modułu.
- Zapytanie obejmujące treść kilku otwartych kart jednocześnie.
- Wstawienie odpowiedzi jako notatki bezpośrednio do Notes Panel, powiązanej z bieżącym źródłem.
- Cytowanie fragmentu strony w odpowiedzi z odnośnikiem do pozycji w Sources Panel.
- Wywołanie funkcji warstwy 4 poleceniem języka naturalnego (narzędzia inspekcyjne, macierz izolacji, edytor scenariusza).
- Historia poleceń, regeneracja odpowiedzi, przekazanie wyniku do modułu Studio.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Chat Window ───────────────────┐
│ Użytkownik ↔ Wykonawca    [ ⋮ ] │
├─────────────────────────────────┤
│ ┌─ Użytkownik ────────────────┐ │
│ │ „wyodrębnij tabelę cen      │ │
│ │  z tej strony"              │ │
│ └─────────────────────────────┘ │
│ ┌─ Wykonawca (strumień) ──────┐ │
│ │ tabela wyodrębniona         │ │
│ │ z sekcji „Cennik" …         │ │
│ │ [ ⋮ ]                       │ │
│ └─────────────────────────────┘ │
├─────────────────────────────────┤
│ Pole poleceń ………………  [ Wyślij ▶]│
│ [Browser] [Karta sesji] [Model ▼]│
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Strumień rozmowy | Lista dymków Użytkownika i Wykonawcy | Prezentacja poleceń i odpowiedzi | Blok dominujący kolumny | 1 | Widoczny bez interakcji | ładowanie · strumień · zakończony | Przewijanie do najnowszej wypowiedzi |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia dotyczącego przeglądanej strony | Duże pole w stopce kolumny | 1 | Widoczne bez interakcji | domyślny · fokus | Enter wysyła polecenie |
| Pasek kontekstu | Rząd lekkich znaczników `[Browser] [Karta sesji] [Model ▼]` | Prezentacja i zmiana kontekstu pracy | Bardzo lekki pasek | 1 | Widoczny bez interakcji | domyślny · najechanie | Kliknięcie znacznika otwiera odpowiedni selektor |
| Selektor kanału modelu | Menu progresywne `Model ▼` | Wybór modelu analizującego współdzielony podgląd | Zwinięty element z etykietą | 2 | Kliknięcie znacznika w pasku kontekstu | zwinięty · rozwinięty · ładowanie | Wybór zmienia model kolejnych analiz i zwija menu |
| Menu akcji dymka | Menu kebab (⋮) przy odpowiedzi | Dodaj jako notatkę · Kopiuj · Regeneruj · Wyślij do Studio | Mała ikona w rogu dymka | 3 | Kliknięcie ⋮ | zwinięte · rozwinięte | Wybór akcji wykonuje operację i zwija menu |
| Odnośnik cytowania fragmentu strony | Mała plakietka w treści odpowiedzi | Wskazanie miejsca na stronie, z którego pochodzi fragment | Plakietka `.dn-plakietka--info` | 1 | Widoczna w treści odpowiedzi | domyślny · najechanie | Kliknięcie przewija Browser Window do wskazanego fragmentu |
| Sterowanie przebiegiem | Przycisk przerwania działania | Zatrzymanie realizowanego zadania Wykonawcy | Mały przycisk przy strumieniu | 1 | Widoczny w trakcie strumienia | ukryty · aktywny | Przerywa zadanie i zapisuje stan w Execution Loop Window |
| Historia poleceń | Lista wcześniejszych poleceń | Powrót do wcześniejszego polecenia | Panel wysuwany | 3 | Menu kebab (⋮) nagłówka | zwinięta · rozwinięta | Wybór wstawia polecenie do pola |

---

### 4.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — okno pętli wykonawczej Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Warstwa widoczności | 1 |
| Rola | Prowadzenie pętli wykonawczej zadań modułu, koordynacja, nadzór nad realizacją i kontrola przebiegu procesów przeglądania |
| Izolacja domyślna | Pętla prowadzona w obrębie karty sesji; log przebiegu zapisywany jako `artefakt` |

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie modułu i jego dekompozycja na zadania: otwarcie i nawigacja stron, migawka treści dla Wykonawcy, ekstrakcja tabel i danych strukturalnych, archiwizacja strony, zrzut ekranu, tłumaczenie, odtworzenie makra, przebieg scenariusza warunkowego, sprawdzenie monitora zmian.
- Kolejka i stan zadań przeglądania z priorytetami: oczekujące, w realizacji, zakończone, ponawiane, przerwane.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą: przydział zadania, potwierdzenie, raport wyniku, zgłoszenie błędu strony.
- Wyniki kontroli jakości: asercje obecności elementów na stronie, walidacja wyodrębnionych danych względem `JSON Schema`, kontrola kompletności archiwum, oraz decyzje o ponowieniu zadania.
- Wskaźniki przebiegu pętli: liczba kroków, czas realizacji, wykorzystanie limitów kroków Wykonawcy operującego na stronie.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia bez utraty dotychczasowych wyników.
- Powiązanie każdego zadania z pozycją w Sources Panel, Notes Panel lub Capture & Monitor Panel powstałą w wyniku jego realizacji.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Execution Loop Window ─────────┐
│ Koordynator ↔ Wykonawca   [ ⋮ ] │
├─────────────────────────────────┤
│ Zlecenie: „zbierz cenniki        │
│ trzech dostawców"                │
│ ───────────────────────────────  │
│ ▶ 1. otwórz stronę A     zakończ.│
│ ▶ 2. wyodrębnij tabelę   zakończ.│
│ ▶ 3. otwórz stronę B     w toku  │
│ ▶ 4. kontrola jakości    oczekuje│
│ ───────────────────────────────  │
│ Kroki 12/40 · ponowienia 1       │
├─────────────────────────────────┤
│ [ Wstrzymaj ] [ Przerwij ] [ ⋮ ] │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok tekstowy z treścią zlecenia | Prezentacja bieżącego zlecenia modułu | Blok w nagłówku kolumny | 1 | Widoczny bez interakcji | domyślny · skorygowany | Kliknięcie otwiera pole korekty zlecenia |
| Lista zadań pętli | Kolejka zadań przeglądania | Śledzenie dekompozycji i stanu realizacji | Lista dominująca kolumny | 1 | Widoczna bez interakcji | oczekujące · w realizacji · zakończone · ponawiane · przerwane | Kliknięcie zadania rozwija komunikaty sterujące i wynik |
| Wskaźnik przebiegu pętli | Licznik kroków, ponowień i czasu | Kontrola wykorzystania limitów przebiegu | Lekki pasek wskaźników | 1 | Widoczny bez interakcji | w normie · przy limicie · przekroczony | Przekroczenie limitu wstrzymuje pętlę i zgłasza to w Chat Window |
| Sterowanie przebiegiem | Przyciski Wstrzymaj / Wznów / Przerwij | Kontrola realizacji procesów modułu | Rząd małych przycisków w stopce kolumny | 1 | Widoczne bez interakcji | aktywne · nieaktywne | Zmienia stan pętli i odnotowuje decyzję w logu |
| Karta komunikatu sterującego | Wpis wymiany Koordynator ↔ Wykonawca | Prezentacja przydziału, potwierdzenia i raportu | Karta `.dn-karta--pozycja` | 2 | Kliknięcie zadania na liście | zwinięta · rozwinięta | Rozwinięcie pokazuje pełną treść komunikatu |
| Wynik kontroli jakości | Plakietka wyniku asercji i walidacji | Ocena poprawności wyniku zadania | Plakietka przy zadaniu | 2 | Kliknięcie zadania na liście | zaliczone · odrzucone · ponowione | Odrzucenie uruchamia ponowienie zadania |
| Menu operacji pętli | Menu kebab (⋮) | Korekta zlecenia, eksport logu, przekazanie do Automations | Mała ikona w nagłówku i stopce | 3 | Kliknięcie ⋮ | zwinięte · rozwinięte | Wybór akcji wykonuje operację i zwija menu |
| Limity kroków Wykonawcy | Zestaw parametrów przebiegu | Ograniczenie liczby kroków i zakresu działania na stronie | Panel wysuwany | 4 | Polecenie języka naturalnego w Chat Window lub tryb administracyjny | domyślny · zmieniony | Zmiana obowiązuje od kolejnego zadania pętli |

---

### 4.3. Browser Window

| Aspekt | Wartość |
|---|---|
| Typologia | Okno główne — osadzona przeglądarka współdzielona z Wykonawcą |
| Waga wizualna | Prawa kolumna, dominująca |
| Warstwa widoczności | 1 |
| Izolacja domyślna | Podgląd wspólny w obrębie karty sesji; historia przeglądania odrębna per karta |

**Zawartość i pełny arsenał funkcji.**

*Nawigacja.*
- Pasek adresu z historią, podpowiedziami, multi-wyszukiwarką (`!g`, `!yt`) i szybkimi poleceniami wywoływanymi znakiem `/`.
- Karty przeglądania z grupami kart oraz przełącznikiem przestrzeni roboczych, wewnątrz jednego okna, niezależnie od kart sesji platformy.
- Przyciski wstecz, dalej i odśwież; zakładki; historia przeglądania sesji.
- Wyszukiwanie w treści bieżącej strony.

*Widoczność dla Wykonawcy.*
- Wspólny podgląd — Wykonawca odbiera dokładnie renderowaną treść strony w czasie rzeczywistym, łącznie ze zmianami po interakcji użytkownika (przewinięcie, rozwinięcie sekcji, wypełnienie formularza).
- Wskaźnik obecności Wykonawcy na stronie.

*Praca z treścią strony.*
- Tryb czytnika — oczyszczenie strony z reklam i elementów pobocznych do czystego tekstu; wymuszony tryb ciemny i nakładki dostępności.
- Zaznaczenie fragmentu strony jako przedmiotu polecenia, z paskiem pływającym: Wyjaśnij / Wyodrębnij / Notatka / Tłumacz / Zapytaj.
- Podświetlanie i adnotowanie fragmentów strony bezpośrednio na podglądzie.
- Zrzut ekranu w czterech trybach oraz archiwum strony, zapisywane jako źródło wizualne.
- Wyodrębnianie danych strukturalnych: tabel, list, danych kontaktowych, cen — do formatu gotowego dalszego przetwarzania.
- Tłumaczenie strony na żądanie mechanizmem modułu Translate.
- Wbudowany podgląd plików PDF i dokumentów biurowych bez opuszczania okna.
- Widok podzielony na 2–4 kolumny do bezpośredniego porównania stron.
- Nagrywanie makra przeglądania przekazywane do Automation Studio.
- Narzędzia inspekcyjne: DOM i style, monitor sieci `HAR`, konsola, emulacja urządzeń, podgląd źródła i różnic.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Browser Window ──────────────────────────────┐
│ [‹][›][⟳] [ adres-strony.com/cennik ……… ] [ ⋮ ]│
│ [ Karta 1 · Cennik ][ Karta 2 · Konkurent ][+] │
├───────────────────────────────────────────────┤
│                                               │
│      wyrenderowana treść strony               │
│      internetowej (widoczna jednocześnie      │
│      użytkownikowi i Wykonawcy)               │
│                                               │
│                                               │
│                                    (◉ Wykonawca)│
├───────────────────────────────────────────────┤
│ [Ubuntu] [Karta sesji] [Czytnik ▼] Operacje ▼ ☰│
└───────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Pasek adresu | Pole `.dn-input` z historią, podpowiedziami i szybkimi poleceniami wywoływanymi znakiem `/` | Nawigacja do adresu lub wyszukiwanie | Średnie pole w nagłówku okna | 1 | Widoczny bez interakcji | pusty · z podpowiedziami · ładowanie strony (`.dn-spinner`) | Enter ładuje stronę i aktualizuje wspólny podgląd |
| Przyciski nawigacji | Ikony wstecz, dalej, odśwież | Standardowa nawigacja przeglądarki | Małe ikony 24 px | 1 | Widoczne bez interakcji | domyślny · wyłączony | Zmienia wyświetlaną stronę |
| Karty przeglądania | Rząd zakładek z grupami | Wiele stron równolegle w oknie | Rząd `.dn-zakladki` | 1 | Widoczny bez interakcji | aktywna · nieaktywna · ładowanie · zawieszona | Kliknięcie przełącza wyświetlaną stronę |
| Obszar renderowania strony | Wyrenderowana treść internetowa | Główna, współdzielona przestrzeń podglądu | Blok dominujący kolumny | 1 | Widoczny bez interakcji | ładowanie · wyrenderowany · błąd wczytania | Przewijanie, klikanie, zaznaczanie jak w przeglądarce |
| Wskaźnik obecności Wykonawcy | Mała plakietka `.dn-plakietka--info` | Sygnalizacja, że Wykonawca odbiera bieżący widok | Bardzo mała plakietka | 1 | Widoczna bez interakcji | nieaktywny · aktywny | Element informacyjny |
| Pasek kontekstu okna | Znaczniki `[Ubuntu] [Karta sesji] [Czytnik ▼]` | Prezentacja i zmiana kontekstu przeglądania | Bardzo lekki pasek | 1 | Widoczny bez interakcji | domyślny · najechanie | Kliknięcie znacznika otwiera selektor |
| Menu przestrzeni roboczych | Menu progresywne | Przełączenie zestawu kart, zakładek i historii | Zwinięty element w pasku kontekstu | 2 | Kliknięcie znacznika karty sesji | zwinięte · rozwinięte | Wybór przełącza przestrzeń i zwija menu |
| Przełącznik trybu czytnika | Menu progresywne `Czytnik ▼` | Oczyszczenie strony, typografia, tryb ciemny, nakładki dostępności | Zwinięty element w pasku kontekstu | 2 | Kliknięcie znacznika | zwinięty · rozwinięty · włączony | Przełącza renderowanie na uproszczony widok tekstowy |
| Ikona zakładki | Ikona gwiazdki | Dodanie strony do ulubionych sesji | Mała ikona | 2 | Menu kebab (⋮) paska adresu | nieaktywna · aktywna | Zapisuje adres na liście zakładek sesji |
| Pasek pływający zaznaczenia | Kontekstowy pasek narzędzi | Operacje na zaznaczonym fragmencie strony | Mały pasek nad zaznaczeniem | 3 | Zaznaczenie fragmentu strony | ukryty · widoczny | Wysyła polecenie do Chat Window albo tworzy notatkę |
| Menu `Operacje ▼` | Grupowanie logiczne akcji | Podział widoku, zrzut ekranu, archiwum strony, wyodrębnianie danych, tłumaczenie, nagranie makra | Zwinięty element zbiorczy w pasku kontekstu | 3 | Kliknięcie `Operacje ▼` | zwinięte · rozwinięte · ładowanie | Wybór wykonuje operację i zwija menu |
| Menu hamburger (☰) | Menu okna | Zakładki, historia, sesje, import i eksport profilu, pobrania | Mała ikona w pasku kontekstu | 3 | Kliknięcie ☰ | zwinięte · rozwinięte | Otwiera panel wysuwany właściwej sekcji |
| Narzędzia inspekcyjne | Panel DOM, sieci `HAR`, konsoli i emulacji urządzeń | Inspekcja strony i diagnostyka | Panel wysuwany | 4 | Polecenie języka naturalnego w Chat Window lub skrót klawiszowy | zamknięty · otwarty | Otwiera panel inspekcji nad obszarem renderowania |

---

### 4.4. Sources Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno list i źródeł |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa widoczności | 2 — wywołanie znacznikiem `Źródła ▼` w pasku kontekstu |
| Izolacja domyślna | Lista narasta w miarę przeglądania kolejnych stron w bieżącej karcie sesji |

**Zawartość i pełny arsenał funkcji.**

- Automatyczne dodawanie źródła przy odwiedzeniu nowej strony (adres, tytuł, migawka, znacznik czasu).
- Ręczne dodanie źródła adresem, bez odwiedzenia strony w bieżącej sesji.
- Grupowanie źródeł w zestawy tematyczne w obrębie sesji.
- Ocena istotności źródła (oznaczenie „kluczowe").
- Podgląd migawki strony bez ponownego jej ładowania.
- Wykrywanie ponownie odwiedzonej strony (scalanie zamiast duplikowania wpisu).
- Eksport listy źródeł jako bibliografii w formatach `BibTeX`, `RIS`, `CSV`, `Markdown` oraz przekazanie zbiorcze do modułu Research.
- Usunięcie źródła z listy sesji.
- Wyszukiwanie i filtrowanie listy według tytułu, domeny, daty i oznaczenia istotności.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Sources Panel ─────────────────┐
│ [ 🔍 szukaj ]             [ ⋮ ] │
├─────────────────────────────────┤
│ ★ Cennik — adres-strony.com     │
│   12:04                    [ ⋮ ]│
├─────────────────────────────────┤
│ ○ Oferta konkurenta —           │
│   konkurent.pl/oferta 11:58 [⋮ ]│
├─────────────────────────────────┤
│ ○ Artykuł branżowy —            │
│   branza.com/news/12  11:40 [⋮ ]│
├─────────────────────────────────┤
│ Operacje ▼                      │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Pole wyszukiwania | Pole `.dn-input` z ikoną | Odnalezienie źródła na liście | Małe pole w nagłówku panelu | 1 | Widoczne po otwarciu panelu | puste · z wynikami | Filtruje listę na żywo |
| Pozycja źródła | Wiersz karty `.dn-karta--pozycja` | Reprezentuje jedną odwiedzoną lub dodaną stronę | Średni wiersz z ikoną istotności, tytułem, domeną i czasem | 1 | Widoczna po otwarciu panelu | domyślny · kluczowe · najechanie | Kliknięcie otwiera stronę w Browser Window |
| Ikona istotności | Ikona gwiazdki | Oznaczenie źródła jako kluczowego dla sesji | Bardzo mała ikona | 1 | Widoczna przy pozycji | nieaktywna · aktywna | Kliknięcie przełącza status istotności |
| Filtry listy | Zestaw filtrów: domena, data, istotność | Zawężenie listy źródeł | Panel wysuwany | 2 | Znacznik filtru w nagłówku panelu | zwinięte · rozwinięte | Zastosowanie filtru zwija panel |
| Menu pozycji (⋮) | Menu kebab pozycji źródła | Podgląd migawki, grupowanie, usunięcie, eksport pojedynczy, dodanie monitora zmian | Mała ikona przy pozycji | 3 | Kliknięcie ⋮ | zwinięte · rozwinięte | Otwiera menu kontekstowe pozycji |
| Menu `Operacje ▼` panelu | Grupowanie logiczne akcji | Dodanie źródła ręcznie, eksport listy, przekazanie do Research | Zwinięty element zbiorczy w stopce panelu | 3 | Kliknięcie `Operacje ▼` | zwinięte · rozwinięte · ładowanie | Wykonuje operację i zwija menu |
| Reguły scalania duplikatów | Parametry wykrywania ponownych odwiedzin | Sterowanie scalaniem wpisów | Panel konfiguracji | 4 | Polecenie języka naturalnego w Chat Window lub tryb administracyjny | domyślny · zmieniony | Zmiana obowiązuje dla kolejnych wpisów |

---

### 4.5. Notes Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Zarządzanie źródłami i ustaleniami |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa widoczności | 2 — wywołanie znacznikiem `Notatki ▼` w pasku kontekstu |
| Izolacja domyślna | Notatki powiązane ze źródłem z Sources Panel bieżącej sesji |

**Zawartość i pełny arsenał funkcji.**

- Ręczne dodanie notatki dotyczącej bieżącej lub wcześniej odwiedzonej strony.
- Automatyczne dodanie notatki z odpowiedzi Wykonawcy (akcja „Dodaj jako notatkę" w Chat Window albo w pasku pływającym Browser Window).
- Powiązanie notatki z konkretnym fragmentem strony — cytat źródłowy zachowany razem z notatką.
- Klasyfikacja notatki: obserwacja, cytat, pytanie otwarte, wniosek.
- Edycja treści notatki, usuwanie, przypinanie notatek najważniejszych na początku listy.
- Grupowanie notatek według źródła lub według własnych wątków tematycznych.
- Wyszukiwanie pełnotekstowe w treści notatek.
- Eksport notatek jako osobnego dokumentu oraz przekazanie zbiorcze do modułu Research (Findings Panel) albo Library.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Notes Panel ───────────────────┐
│ Notatki sesji     [ Widok ▼ ][⋮]│
├─────────────────────────────────┤
│ 📌 wniosek · źródło: Cennik      │
│ „ceny wzrosły o 8% względem      │
│  ubiegłego kwartału"       [ ⋮ ]│
├─────────────────────────────────┤
│ 💬 cytat · Oferta konkurenta     │
│ „gwarancja rozszerzona w cenie   │
│  podstawowej"              [ ⋮ ]│
├─────────────────────────────────┤
│ ❓ pytanie otwarte               │
│ „czy dotyczy to rynku B2B?"[ ⋮ ]│
├─────────────────────────────────┤
│ Operacje ▼                      │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Pozycja notatki | Karta `.dn-karta--pozycja` | Reprezentuje pojedynczą notatkę | Średnia karta z ikoną klasyfikacji | 1 | Widoczna po otwarciu panelu | wniosek · cytat · pytanie otwarte · obserwacja · przypięta | Kliknięcie rozwija treść notatki |
| Ikona klasyfikacji | Mała ikona przy notatce | Wizualne rozróżnienie typu notatki | Bardzo mała ikona | 1 | Widoczna w nagłówku karty | domyślny | Kliknięcie otwiera wybór klasyfikacji |
| Odnośnik do źródła | Mały link tekstowy z ikoną `link-zewnetrzny` | Powiązanie notatki z konkretną stroną w Sources Panel | Mały tekst | 1 | Widoczny w treści karty | domyślny | Kliknięcie otwiera stronę w Browser Window |
| Przełącznik widoku | Menu progresywne `Widok ▼` | Grupowanie według źródła, wątku albo chronologii | Zwinięty element w nagłówku panelu | 2 | Kliknięcie `Widok ▼` | zwinięty · rozwinięty | Zmienia prezentację listy i zwija menu |
| Wyszukiwanie pełnotekstowe | Pole `.dn-input` | Odnalezienie notatki po treści | Panel wysuwany | 2 | Znacznik wyszukiwania w nagłówku panelu | puste · z wynikami | Filtruje listę na żywo |
| Menu pozycji (⋮) | Menu kebab notatki | Edycja, przypięcie, zmiana klasyfikacji, usunięcie, otwarcie źródła | Mała ikona przy karcie | 3 | Kliknięcie ⋮ | zwinięte · rozwinięte | Otwiera menu kontekstowe notatki |
| Menu `Operacje ▼` panelu | Grupowanie logiczne akcji | Nowa notatka, eksport, przekazanie do Research lub Library | Zwinięty element zbiorczy w stopce panelu | 3 | Kliknięcie `Operacje ▼` | zwinięte · rozwinięte · wysłano | Wykonuje operację i zwija menu |
| Schemat eksportu ustaleń | Definicja pól przekazywanych do Findings Panel | Sterowanie kształtem przekazywanych ustaleń | Panel konfiguracji | 4 | Polecenie języka naturalnego w Chat Window lub tryb administracyjny | domyślny · zmieniony | Zmiana obowiązuje dla kolejnych eksportów |

---

### 4.6. Automation Studio

| Aspekt | Wartość |
|---|---|
| Typologia | Okno operacji i scenariuszy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne, pełna wysokość obszaru roboczego |
| Warstwa widoczności | 3 — wywołanie z menu `Operacje ▼` albo poleceniem języka naturalnego w Chat Window |
| Izolacja domyślna | Scenariusze zapisane w obrębie karty sesji; przekazanie do Automations konfigurowalne |

**Zawartość i pełny arsenał funkcji.**

- Lista makr i scenariuszy z ostatnim statusem wykonania.
- Nagrywarka kroków: start, stop, krok pojedynczy.
- Edytor scenariusza: kroki, warunki, pętle, ponowienia, obsługa błędów — w formie kart oraz w widoku `JSON`/`YAML`.
- Harmonogram uruchomienia: jednorazowo, cyklicznie (cron), po zdarzeniu.
- Kolejka wykonania z logiem i artefaktami wyniku, prowadzona przez Koordynatora i prezentowana w Execution Loop Window.
- Katalog profili autowypełniania z kluczami jawnymi oraz wsadem masowym z pliku `CSV`.
- Limity kroków i zakresu działania Wykonawcy operującego na stronie.
- Przekazanie scenariusza do modułu Automations.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Automation Studio ─────────────┐
│ Makra i scenariusze       [ ⋮ ] │
├─────────────────────────────────┤
│ ▶ Zbiórka cenników    ok  [ ⋮ ] │
│   cron · co 24 h                │
├─────────────────────────────────┤
│ ▶ Formularz zgłoszeń  ok  [ ⋮ ] │
│   po zdarzeniu                  │
├─────────────────────────────────┤
│ ▶ Monitor dostępności ponow.[⋮ ]│
│   cron · co 1 h                 │
├─────────────────────────────────┤
│ [ ● Nagrywaj ]      Operacje ▼  │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Lista scenariuszy | Lista kart makr i scenariuszy | Przegląd i uruchamianie przebiegów | Lista dominująca kolumny | 1 | Widoczna po otwarciu okna | ok · ponawiany · błąd · wstrzymany | Kliknięcie otwiera szczegóły scenariusza |
| Przycisk nagrywania | Przycisk `--zloty` z ikoną | Rejestracja kroków przeglądania jako makra | Mały przycisk w stopce kolumny | 1 | Widoczny po otwarciu okna | gotowy · nagrywanie · zatrzymany | Rozpoczyna rejestrację kroków w Browser Window |
| Harmonogram scenariusza | Zestaw ustawień uruchomienia | Wybór trybu: jednorazowo, cron, po zdarzeniu | Panel wysuwany | 2 | Znacznik harmonogramu przy scenariuszu | zwinięty · rozwinięty | Zapisuje harmonogram i zwija panel |
| Edytor kroków | Widok kart kroków scenariusza | Budowa i korekta przebiegu | Panel wysuwany | 3 | Menu kebab (⋮) scenariusza | zwinięty · rozwinięty · zmieniony | Zmiany obowiązują od kolejnego uruchomienia |
| Menu `Operacje ▼` okna | Grupowanie logiczne akcji | Uruchomienie, duplikowanie, eksport, przekazanie do Automations | Zwinięty element zbiorczy w stopce kolumny | 3 | Kliknięcie `Operacje ▼` | zwinięte · rozwinięte | Wykonuje operację i zwija menu |
| Profile autowypełniania | Katalog profili z kluczami jawnymi | Wstawianie zapisanych danych do formularzy | Panel wysuwany | 3 | Menu `Operacje ▼` | zwinięty · rozwinięty | Wybór profilu przypisuje go do kroku scenariusza |
| Widok `JSON`/`YAML` scenariusza | Tekstowa reprezentacja przebiegu | Precyzyjna edycja warunków, pętli i ponowień | Panel pełnej wysokości kolumny | 4 | Polecenie języka naturalnego w Chat Window lub skrót klawiszowy | zamknięty · otwarty · walidacja | Zapis waliduje strukturę i aktualizuje karty kroków |
| Limity przebiegu | Parametry liczby kroków i zakresu działania | Ograniczenie działania Wykonawcy na stronie | Panel konfiguracji | 4 | Tryb administracyjny lub polecenie w Chat Window | domyślny · zmieniony | Obowiązuje od kolejnego przebiegu pętli |

---

### 4.7. Capture & Monitor Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno przechwytywania i monitorowania |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa widoczności | 3 — wywołanie z menu `Operacje ▼` albo z menu kebab (⋮) pozycji źródła |
| Izolacja domyślna | Zrzuty, archiwa i monitory prowadzone w obrębie karty sesji |

**Zawartość i pełny arsenał funkcji.**

- Galeria zrzutów i archiwów bieżącej sesji z adnotacją (strzałki, wyróżnienia, rozmycie danych wrażliwych, numeracja kroków).
- Menedżer pobrań: kolejka, wznawianie, reguły folderów, kategorie plików.
- Lista monitorów śledzenia zmian z historią różnic i alertami.
- Czytnik RSS/Atom z wykryciem kanału na przeglądanej stronie.
- Kolejka czytania z przypomnieniami.
- Podgląd i konwersja obrazów oraz ekstrakcja mediów ze strony.
- Eksport zebranego materiału do modułów Library i Research.

**Makieta tekstowa — stan spoczynku.**

```
┌─ Capture & Monitor ─────────────┐
│ Materiał sesji     [ Widok ▼ ][⋮]│
├─────────────────────────────────┤
│ ▣ Zrzut · Cennik  12:06   [ ⋮ ] │
├─────────────────────────────────┤
│ ▣ Archiwum · Oferta 11:59 [ ⋮ ] │
├─────────────────────────────────┤
│ ◔ Monitor · cena A   zmiana [⋮ ]│
├─────────────────────────────────┤
│ ⤓ Pobranie · raport.pdf   [ ⋮ ] │
├─────────────────────────────────┤
│ Operacje ▼                      │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| Lista materiału sesji | Wspólna lista zrzutów, archiwów, monitorów i pobrań | Przegląd całego przechwyconego materiału | Lista dominująca kolumny | 1 | Widoczna po otwarciu panelu | zapisany · w toku · zmiana wykryta · błąd | Kliknięcie otwiera podgląd pozycji |
| Podgląd pozycji | Widok zrzutu, archiwum albo różnicy monitora | Ocena przechwyconego materiału | Panel wysuwany | 1 | Kliknięcie pozycji | zamknięty · otwarty | Prezentuje materiał bez opuszczania panelu |
| Przełącznik widoku | Menu progresywne `Widok ▼` | Podział na zrzuty, archiwa, monitory, pobrania, kanały | Zwinięty element w nagłówku panelu | 2 | Kliknięcie `Widok ▼` | zwinięty · rozwinięty | Zawęża listę i zwija menu |
| Menu pozycji (⋮) | Menu kebab pozycji | Adnotacja, konwersja, ponowne pobranie, usunięcie, eksport | Mała ikona przy pozycji | 3 | Kliknięcie ⋮ | zwinięte · rozwinięte | Otwiera menu kontekstowe pozycji |
| Edytor adnotacji zrzutu | Warstwa rysowania na zrzucie | Strzałki, wyróżnienia, rozmycie, numeracja kroków | Panel wysuwany | 3 | Menu kebab (⋮) pozycji zrzutu | zamknięty · otwarty · zapisany | Zapis tworzy `artefakt` w bieżącej sesji |
| Menu `Operacje ▼` panelu | Grupowanie logiczne akcji | Nowy monitor, subskrypcja kanału, eksport do Library i Research | Zwinięty element zbiorczy w stopce panelu | 3 | Kliknięcie `Operacje ▼` | zwinięte · rozwinięte · wysłano | Wykonuje operację i zwija menu |
| Parametry monitora | Częstotliwość sprawdzeń, próg zmiany, kanał alertu | Sterowanie czułością śledzenia zmian | Panel konfiguracji | 4 | Polecenie języka naturalnego w Chat Window lub tryb administracyjny | domyślny · zmieniony | Obowiązuje od kolejnego sprawdzenia |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — wspólna analiza strony

```
 Chat Window        Execution Loop      Browser Window      Sources / Notes
 ───────────        ──────────────      ──────────────      ───────────────
 1. polecenie ─────► 2. dekompozycja
    użytkownika         zlecenia
                        │
                        ├──► zadanie: otwórz stronę ─────► render + migawka
                        │
                        ├──► zadanie: migawka treści ────► odbiór przez
                        │                                   Wykonawcę
                        │
 4. odpowiedź ◄─────── 3. kontrola jakości              ──► źródło dodane
    z cytowaniem          i raport wyniku                    automatycznie
                        │
                        └──────────────────────────────► notatka powiązana
                                                          ze źródłem
```

### 5.2. Przepływ — porównanie ofert

```
 Browser Window — Karta 1 (oferta A)  │  Browser Window — Karta 2 (oferta B)
 ────────────────────────────────────  │  ────────────────────────────────────
              └────────► widok podzielony ◄────────┘
                                │
 Chat Window ──► „porównaj cenniki obu ofert"
                                │
 Execution Loop ──► zadania: migawka A · migawka B · zestawienie · kontrola
                                │
 Notes Panel ──► wniosek z porównania, powiązany z oboma źródłami
```

### 5.3. Przepływ — przebieg automatyzacji web

```
 Chat Window ──► zlecenie: „zbierz cenniki trzech dostawców"
                                │
 Execution Loop Window ──► dekompozycja na zadania i kolejka
                                │
   ┌────────────┬───────────────┼───────────────┬────────────┐
   ▼            ▼               ▼               ▼            ▼
 nawigacja   ekstrakcja      walidacja       ponowienie    zapis
 (Browser)   tabel           JSON Schema     zadania       artefaktu
   │            │               │               │            │
   └────────────┴───────────────┴───────────────┴────────────┘
                                │
 Automation Studio ──► scenariusz zapisany · harmonogram · log przebiegu
                                │
 Capture & Monitor Panel ──► artefakty wyniku i monitory zmian
```

### 5.4. Przepływ — zasilenie modułu Research

```
 Sources Panel (lista źródeł sesji)  │  Notes Panel (notatki sesji)
 ──────────────────────────────────  │  ───────────────────────────
              └───────────────┬──────────────────┘
                              ▼
                 powiązanie Browser ───► Research
                 (konfigurowalne, okno konfiguracji)
                              │
                              ▼
   Research › Sources Manager   │   Research › Findings Panel
   (źródła scalone z listą      │   (notatki jako ustalenia
    badania)                    │    cząstkowe)
```

### 5.5. Przepływ — trwałe przechowanie materiału

```
 Browser Window ──► zrzut ekranu albo archiwum strony
              │
              ▼
 Capture & Monitor Panel ──► artefakt wizualny w galerii sesji
              │
              ▼
 Sources Panel ──► pozycja źródła z odwołaniem do artefaktu
              │
              ▼
 powiązanie Browser ───► Library (konfigurowalne)
              │
              ▼
 Library › Library Explorer ──► trwałe przechowanie i katalogowanie
                                 w Tags & Collections
```

---

## 6. Stany, dane i powiązania

### 6.1. Model stanów sesji przeglądania

```
┌────────────────┐  otwarcie strony  ┌────────────────┐  zaznaczenie/polecenie ┌────────────────┐
│ Pusta karta    │──────────────────►│ Przeglądanie   │──────────────────────►│ Analiza wspólna│
│ (start.strona) │                   │ (wspólny widok)│                       │ z Wykonawcą    │
└────────────────┘                   └───────┬────────┘                       └───────┬────────┘
                                             │ kolejna strona                         │ notatka/źródło
                                             ▼                                        ▼
                                    ┌────────────────┐                       ┌────────────────┐
                                    │ Kolejne źródło │                       │ Materiał gotowy│
                                    │ dodane         │                       │ do przekazania │
                                    │ automatycznie  │                       │                │
                                    └───────┬────────┘                       └────────────────┘
                                            │ zlecenie powtarzalne
                                            ▼
                                    ┌────────────────┐
                                    │ Pętla          │
                                    │ wykonawcza     │
                                    │ (Koordynator)  │
                                    └────────────────┘
```

### 6.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Browser |
|---|---|
| `sesja`, `karta_sesji` | Nośnik kontekstu przeglądania i listy źródeł |
| `wiadomosc` | Historia Chat Window oraz komunikaty sterujące Execution Loop Window |
| `artefakt` | Zrzuty ekranu, archiwa stron, wyodrębnione dane strukturalne i logi przebiegu pętli (Załącznik F.2 Modelu danych) |
| `etykieta_artefaktu` | Klasyfikacja notatek i grupowanie źródeł |
| `kanal_modelu` | Model bazowy analizujący współdzielony podgląd oraz prowadzący ekstrakcję |
| `powiazanie_komponentu` | Jawne powiązania Browser ───► Research, Library, Studio, Translate, Automations, Agents |
| `profil_izolacji`, `regula_izolacji_kontekstu` | Zakres współdzielenia historii przeglądania między kartami, kontenery tożsamości, tryb efemeryczny |

### 6.3. Izolacja i konfigurowalność — punkty właściwe modułowi Browser

| Punkt izolacji | Stan wyjściowy (domyślny) | Co konfiguruje Operator | Gdzie |
|---|---|---|---|
| Historia przeglądania i lista źródeł | Odrębna per karta sesji | Współdzielenie listy źródeł między kartami tego samego badania | Okno konfiguracji, poziom „karta sesji" |
| Zasilanie modułu Research | Brak automatycznego przekazywania | Ustanowienie stałego kanału Browser ───► Research | Okno konfiguracji, powiązanie komponentu |
| Trwałe przechowanie w Library | Brak automatycznego zapisu | Ustanowienie automatycznego przekazywania zrzutów i notatek do Library | Okno konfiguracji, powiązanie komponentu |
| Dostęp sieciowy procesu sesji | Współdzielony dostęp sieciowy serwera | Wydzielenie odrębnego, ograniczonego dostępu sieciowego dla sesji przeglądania materiałów wrażliwych | Okno konfiguracji punktów izolacji, panel macierzy izolacji |
| Kontenery tożsamości | Wspólny magazyn ciasteczek karty sesji | Wydzielenie izolowanych profili ciasteczek i magazynu per karta | Okno konfiguracji, `profil_izolacji` |
| Zakres pętli wykonawczej | Pętla ograniczona do zadań bieżącej karty sesji | Rozszerzenie zakresu pętli na zestaw kart badania | Okno konfiguracji, poziom procesu |

### 6.4. Powiązania z innymi modułami

```
        BROWSER   ◄──►  Research    (wspólne źródła i zasilanie badania)
           │      ◄──►  Translate   (tłumaczenie stron i wycinków)
           ├──────►  Library     (trwałe przechowanie odnotowanych materiałów)
           ├──────►  Studio      (wyodrębnione dane jako wsad redakcyjny)
           ├──────►  Automations (makra i scenariusze jako kroki automatyzacji)
           └──────►  Agents      (zadania przeglądania dla agentów platformy)
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Research | Zasilanie pracy badawczej zebranymi źródłami | Konfiguracyjne | Sources Panel → Sources Manager |
| Library | Trwałe przechowanie odnotowanych materiałów | Konfiguracyjne | Capture & Monitor Panel → Library Explorer |
| Translate | Tłumaczenie stron i wycinków | Konfiguracyjne | Browser Window → silnik Translate |
| Studio | Wyodrębnione tabele i teksty jako wsad redakcyjny | Konfiguracyjne | Browser Window → Studio Editor |
| Automations | Makra i scenariusze jako kroki większych automatyzacji | Konfiguracyjne | Automation Studio → Automations |
| Agents | Udostępnienie zadań przeglądania agentom platformy | Konfiguracyjne | Execution Loop Window → Agents |

---

## 7. Punkty sterowania z okna konfiguracji

Zgodnie z Modelem konfiguracji (warstwy: aplikacja / proces / akcja / sesja) Operator personalizuje moduł Browser bez blokad — brak ustawienia oznacza wartość domyślną, a wartością domyślną jest wykonanie.

| Zakres | Co Operator personalizuje | Warstwa i miejsce |
|---|---|---|
| **Silniki wyszukiwania** | Domyślna wyszukiwarka, skróty w pasku adresu, wyszukiwanie w platformie względem sieci | Aplikacja / sesja |
| **Karty i przestrzenie** | Zawieszanie kart (próg czasu i pamięci), drzewo kart, domyślna przestrzeń robocza | Aplikacja |
| **Kanał modelu podglądu** | Model analizujący wspólny podgląd; odrębny model dla ekstrakcji i dla rozmowy | Zachowanie modeli (`kanal_modelu`) |
| **Chat Window** | Szerokość lewej kolumny, zakres historii per karta sesji, zestaw akcji dymka | Aplikacja / sesja |
| **Execution Loop Window** | Szerokość kolumny pętli, próg automatycznego otwarcia po podjęciu zlecenia, zakres prezentowanych komunikatów sterujących, polityka ponowień zadań, limity kroków Wykonawcy, kanał zgłaszania przerwań do Chat Window | Procesy / sesja |
| **Tryb czytnika** | Typografia, szerokość kolumny tekstu, tło, włączanie per domena; synteza mowy (głos, tempo) | Aplikacja / akcja |
| **Tłumaczenie** | Tłumaczenie per język i domena; wybór silnika Translate | Integracje |
| **Prywatność i blokowanie** | Listy filtrów reklam i mechanizmów śledzących, poziom blokowania, polityka ciasteczek, per domena JavaScript i multimedia | Aplikacja / izolacja |
| **Izolacja sesji** | Kontenery tożsamości, tryb efemeryczny, wydzielony dostęp sieciowy dla materiałów wrażliwych | Okno punktów izolacji (macierz, 7 poziomów) |
| **Uprawnienia stron** | Domyślne i per domena: kamera, mikrofon, lokalizacja, powiadomienia, autoodtwarzanie | Aplikacja / akcja |
| **Ekstrakcja danych** | Formaty eksportu (CSV/JSON/XLSX/Markdown), schematy pól, model prowadzący ekstrakcję | Akcja |
| **Automatyzacja** | Zakres działania Wykonawcy na stronie, limity kroków, harmonogramy, profile autowypełniania (klucze jawne) | Procesy / akcje |
| **Powiązania modułów** | Włączenie i kierunek Browser → Research / Library / Studio / Translate / Automations / Agents | Komponenty (`powiazanie_komponentu`) |
| **Monitorowanie i alerty** | Częstotliwość sprawdzeń, kanał powiadomień (push, e-mail, w aplikacji), progi zmian | Procesy / integracje |
| **Menedżer pobrań** | Foldery docelowe według typu, liczba wątków, otwieranie po pobraniu, skanowanie | Aplikacja |
| **Warstwy widoczności** | Przypisanie ról do warstw 3 i 4, widoczność narzędzi inspekcyjnych, zestaw akcji w menu `Operacje ▼` | Aplikacja / konfiguracja roli |
| **Skróty i gesty** | Przypisanie skrótów, gesty myszy, paleta poleceń (`Ctrl/Cmd + K`) | Aplikacja |

---

## 8. Scenariusze użycia

**Scenariusz 1 — wspólna weryfikacja informacji.**
Dziennikarz otwiera artykuł budzący wątpliwości w Browser Window. Wykonawca, odbierając tę samą treść, wskazuje w Chat Window niespójność między nagłówkiem a treścią artykułu. Fragment trafia jako cytat do Notes Panel, a strona zapisuje się w Sources Panel z pełnymi metadanymi.

**Scenariusz 2 — porównanie ofert w widoku podzielonym.**
Zespół zakupowy otwiera dwie oferty konkurencyjne w dwóch kartach Browser Window, włącza widok podzielony i poleca w Chat Window porównanie cenników. Koordynator dekomponuje polecenie na zadania migawek i zestawienia, a Execution Loop Window pokazuje ich przebieg. Wniosek zapisuje się w Notes Panel jako pozycja powiązana z obydwoma źródłami.

**Scenariusz 3 — zasilenie badania rynkowego.**
Analityk prowadzi sesję przeglądania stron konkurencji. Dzięki skonfigurowanemu powiązaniu z modułem Research zebrane źródła i notatki trafiają do Sources Manager i Findings Panel, gdzie kontynuowana jest systematyczna analiza.

**Scenariusz 4 — wyodrębnienie danych strukturalnych.**
Użytkownik trafia na stronę z tabelą cennikową i wydaje polecenie wyodrębnienia danych. Wykonawca generuje ustrukturyzowaną tabelę, Koordynator weryfikuje ją względem schematu `JSON Schema`, a wynik trafia do Studio lub do dalszej analizy w Research.

**Scenariusz 5 — trwałe archiwum źródeł wizualnych.**
Podczas researchu wizualnego użytkownik wykonuje zrzuty ekranu kilku stron z przykładami interfejsów konkurencji. Materiał gromadzi się w Capture & Monitor Panel, a powiązanie Browser ───► Library przenosi zrzuty do centralnego repozytorium, gdzie zespół projektowy taguje je i wykorzystuje w module Design.

**Scenariusz 6 — powtarzalna zbiórka cenników.**
Zespół operacyjny nagrywa w Automation Studio makro obejścia trzech stron dostawców i wyodrębnienia tabel cen. Koordynator prowadzi pętlę wykonawczą według harmonogramu cron, Execution Loop Window prezentuje kolejkę zadań i wyniki kontroli jakości, a wykryte różnice cen generują alert w Capture & Monitor Panel.

**Scenariusz 7 — nadzór nad przebiegiem na stronie.**
Użytkownik zleca w Chat Window wypełnienie formularza zgłoszeniowego z danych profilu. Execution Loop Window prezentuje kolejne kroki Wykonawcy na stronie wraz z wykorzystaniem limitu kroków. Użytkownik wstrzymuje przebieg przed krokiem wysyłki, koryguje zlecenie i wznawia pętlę bez utraty dotychczasowych wyników.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + T` | Nowa karta przeglądania | Browser Window |
| `Ctrl/Cmd + L` | Fokus na pasku adresu | Browser Window |
| `Ctrl/Cmd + D` | Dodanie strony do zakładek | Browser Window |
| `Ctrl/Cmd + Shift + S` | Zrzut ekranu | Browser Window |
| `Ctrl/Cmd + Shift + I` | Narzędzia inspekcyjne (warstwa 4) | Browser Window |
| `Ctrl/Cmd + Shift + E` | Otwarcie i zamknięcie kolumny pętli wykonawczej | Execution Loop Window |
| `Ctrl/Cmd + K` | Paleta poleceń i wyszukiwarka funkcji | Chat Window, Browser Window |
| `Ctrl/Cmd + .` | Przerwanie realizowanego zadania | Chat Window, Execution Loop Window |
| Ikona `karta-przegladarki` | Karta przeglądania | Browser Window |
| Ikona `globus` | Źródło internetowe | Sources Panel |
| Ikona `gwiazdka` | Zakładka i źródło kluczowe | Browser Window, Sources Panel |
| Ikona `obraz` | Zrzut ekranu jako artefakt wizualny | Capture & Monitor Panel, Sources Panel |
| Ikona `odpowiedz` | Notatka i cytat | Notes Panel |
| Ikona `link-zewnetrzny` | Odnośnik do źródła | Notes Panel |
| Ikona `petla` | Zadanie pętli wykonawczej | Execution Loop Window |
| Ikona `zegar` | Harmonogram i monitor zmian | Automation Studio, Capture & Monitor Panel |
| Znak `⋮` | Menu kebab — rozwinięcia kontekstowe (warstwa 3) | Wszystkie okna |
| Znak `☰` | Menu hamburger — sekcje okna (warstwa 3) | Browser Window |
| Znak `▼` | Menu progresywne i znacznik kontekstowy (warstwa 2) | Wszystkie okna |

*Koniec dokumentu. Moduł Browser — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
