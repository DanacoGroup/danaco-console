*Dokument specyfikuje interfejs modułu Diagnostics Danaco Console: okna, makiety, elementy, warstwy widoczności i stany.*

# Moduł Diagnostics — dokumentacja projektowa

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
| **Tytuł** | Moduł Diagnostics — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować) |
| **Środowiska dostępności** | CodeStudio (wyłącznie) |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowiska CodeStudio; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 3, 4, 7.3, 9.13, 10.2, 12, Załącznik A) · Specyfikacja modułów (rozdz. 4.11–4.13) · Specyfikacja okien operacyjnych (rozdz. 2.4, 3, 5, 6.13, 10) · System wizualny (rozdz. 8, 10) · Model danych (rozdz. 8, 15, 17) · Model konfiguracji (rozdz. 3–6) · Integracja modeli · Izolacja i zależności (rozdz. 1–6, Załącznik A) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie; klucze konfiguracji jawne; izolacja, uprawnienia i redakcja danych wrażliwych są ustawieniami konfiguracyjnymi Operatora, nie wymogiem; domyślnie poprawki wypracowane przez Wykonawcę otwierają się do przeglądu Operatora, przy czym jest to zachowanie konfigurowalne; domyślne zachowanie modułu = wykonanie (podgląd, agregacja, analiza) |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Katalog funkcji i narzędzi](#2-katalog-funkcji-i-narzędzi)
3. [Komplet okien operacyjnych modułu](#3-komplet-okien-operacyjnych-modułu)
4. [Specyfikacja okien operacyjnych](#4-specyfikacja-okien-operacyjnych)
5. [Przepływy pracy](#5-przepływy-pracy)
6. [Punkty sterowania z okna konfiguracji](#6-punkty-sterowania-z-okna-konfiguracji)
7. [Stany, dane i powiązania](#7-stany-dane-i-powiązania)
8. [Scenariusze użycia](#8-scenariusze-użycia)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja i rola modułu

Diagnostics jest modułem analizy oraz usuwania problemów technicznych — przestrzenią, w której punktem wyjścia nie jest tworzenie nowej funkcjonalności, lecz zrozumienie przyczyny istniejącego błędu lub spadku wydajności na podstawie logów, komunikatów błędów, metryk, śladów wykonania i obserwowalnych objawów. Moduł integruje w jednym miejscu czynności rozproszone między osobne narzędzia: gromadzenie materiału źródłowego, agregację tego materiału w spójny obraz stanu systemu oraz wypracowanie przez Wykonawcę konkretnych, uzasadnionych rekomendacji naprawczych — zamiast pozostawiania użytkownikowi ręcznego przeszukiwania dzienników zdarzeń i samodzielnego wnioskowania o przyczynie.

Diagnostics stanowi kompletną warstwę obserwowalności (observability) platformy AI Workspace OS: jedną przestrzeń pełniącą funkcję agregatora logów, rejestru błędów, warstwy metryk i monitorowania wydajności aplikacji, śledzenia rozproszonego, warstwy prowenancji wywołań modeli, pulpitu zużycia kont i kosztów, monitora dostępności, silnika alertów oraz profilera — wszystko nad **jednym, spójnym strumieniem telemetrii** platformy, powiązanym wprost z encjami modelu danych (`kanal_modelu`, `proces_sesji`, `karta_sesji`, `zadanie`).

Operator nadzorujący pracę Danaco Console prowadzi pełny cykl diagnostyki bez uruchamiania zewnętrznego systemu monitorowania: obserwacja → wykrycie → korelacja → prowenancja → przyczyna → rekomendacja → kontrola kosztu → potwierdzenie skuteczności. Szczególny nacisk — właściwy platformie zbudowanej wokół modeli AI — spoczywa na dwóch obszarach: **prowenancji każdego wywołania modelu** (co, przez jaki kanał, z jakim promptem, kosztem i wynikiem) oraz **zużyciu kont i budżetów** (ile tokenów, żądań i środków zużyła każda sesja, kanał modelu, środowisko i konto dostawcy).

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Diagnostics |
|---|---|
| Deweloperzy prowadzący utrzymanie systemu | Szybkie zlokalizowanie przyczyny zgłoszonego błędu na podstawie logów produkcyjnych |
| Inżynierowie niezawodności i DevOps | Analiza spadku wydajności, korelacja zdarzeń z wielu komponentów w czasie |
| Inżynierowie jakości | Weryfikacja, czy zgłoszony defekt jest znanym, powtarzającym się wzorcem błędu |
| Liderzy techniczni | Przegląd stanu zdrowia systemu przed wydaniem, ocena ryzyka technicznego |
| Operatorzy nadzorujący pracę autonomiczną | Analiza przyczyn niepowodzeń procesów wykonywanych bez nadzoru w środowisku MultitaskingAI |
| Osoby odpowiedzialne za koszt pracy modeli | Kontrola zużycia tokenów i środków per kanał modelu, konto, sesja i środowisko |

### 1.3. Po co — wartość modułu

| Problem klasycznego rozproszenia narzędzi | Rozwiązanie w Diagnostics |
|---|---|
| Logi, rejestr błędów i wiedza o rozwiązaniu rozproszone w osobnych systemach | Logs Viewer, Errors Panel i Recommendations Panel w jednej przestrzeni, nad tym samym materiałem |
| Ręczne przeszukiwanie tysięcy linii logu w poszukiwaniu przyczyny | Diagnostics Center agreguje materiał w spójny, zwięzły obraz stanu systemu |
| Powtarzające się błędy analizowane od nowa za każdym razem | Grupowanie błędów po odcisku (fingerprint) z licznikiem wystąpień i historią częstotliwości |
| Wywołanie modelu jako czarna skrzynka bez wglądu w koszt i przebieg | Provenance Explorer rejestruje każde wywołanie kanału modelu wraz z drzewem śladu, tokenami i kosztem |
| Zużycie środków znane dopiero z rozliczenia dostawcy | Usage & Cost rozlicza tokeny, żądania i koszt per kanał, konto, sesja, projekt i środowisko |
| Rekomendacja naprawcza bez uzasadnienia i bez odniesienia do źródła | Każda rekomendacja powiązana wprost z logiem, błędem, metryką lub śladem, na podstawie którego powstała |
| Wdrożenie poprawki wymaga ręcznego przepisania jej do edytora | Zastosowanie rekomendacji otwiera Code Editor modułu Developer z gotową zmianą do przeglądu |

### 1.4. Granica tematyczna modułu

| Obszar | Zakres |
|---|---|
| Logi | Agregacja strumienia logów na żywo i archiwalnych, wyszukiwanie pełnotekstowe i regex, filtry, deduplikacja, korelacja po identyfikatorze śledzenia, parsowanie strukturalne (JSON/logfmt), archiwum i retencja |
| Błędy | Rejestr błędów, grupowanie po odcisku (fingerprint), stos wywołań ze źródłem, wykrywanie regresji, przypisanie do wdrożenia i commitu, cykl życia statusu |
| Prowenancja wywołań AI | Pełny zapis każdego wywołania `kanal_modelu`: prompt, odpowiedź, tokeny, koszt, opóźnienie, użyte narzędzia, trafienia cache, łańcuch wywołań agenta i podagentów, drzewo śledzenia (trace) |
| Zużycie kont i kosztów | Rozliczenie tokenów, żądań i wydatków per kanał modelu, dostawca, konto i token, sesja, projekt oraz środowisko; limity, budżety i prognozy zużycia |
| Metryki i wydajność | Wskaźniki czasu odpowiedzi, przepustowości, zużycia procesora i pamięci procesów sesji, opóźnień komponentów; szeregi czasowe, percentyle, wykresy trendu, profilowanie |
| Kontrola stanu | Wskaźniki kondycji per komponent, usługa i model, sondy zdrowia, monitoring dostępności, monitoring syntetyczny, historia incydentów |
| Alerty | Reguły progowe i anomalii, kanały powiadomień, wyciszenia, eskalacje, powiązanie z Always On Display |
| Rekomendacje i analiza | Agregacja materiału w spójny obraz, wypracowanie i priorytetyzacja rekomendacji naprawczych, ocena skuteczności, powiązanie ze źródłem |

### 1.5. Rozgraniczenie z innymi modułami

| Poza zakresem | Właściwy moduł / powód |
|---|---|
| Edycja i wdrożenie poprawki (Code Editor, Git) | **Developer** — Diagnostics przekazuje mu rekomendację jako gotową zmianę do przeglądu (Recommendations Panel → Code Editor) |
| Odtworzenie objawu w rzeczywistej powłoce, wykonanie polecenia diagnostycznego | **Terminal** — Diagnostics zleca polecenie weryfikujące, wynik wraca do Logs Viewer i Output Console |
| Orkiestracja ról wykonawczych, silnik kolejek pętli pracy ciągłej | **MultitaskingAI** — Diagnostics odbiera zdarzenia niepowodzeń procesów jako materiał źródłowy (Errors Panel) |
| Definiowanie harmonogramów i procesów automatycznych | **Automations** — Diagnostics obserwuje ich przebiegi, nie tworzy ich |
| Budżetowanie i rozliczenia księgowe organizacji | **poza platformą** — Diagnostics mierzy i prognozuje zużycie techniczne, nie prowadzi księgowości ani rozliczeń finansowych |
| Konfiguracja kanałów modeli i poświadczeń dostawców | **Okno konfiguracji / Integracja modeli** — Diagnostics czyta telemetrię kanałów, nie zarządza ich definicją |

Rozgraniczenie jest miękkie i konfigurowalne — metryki zużycia kont zasilają alerty budżetowe, a prowenancja wywołania otwiera powiązany błąd w Errors Panel albo zadanie poprawki w module Developer. Powiązania te są jawne w oknie konfiguracji (rozdz. 6), nigdy wbudowane na stałe. Diagnostyka nie jest bramką: żaden wskaźnik „czerwony" nie blokuje pracy platformy — sygnalizuje i rekomenduje, decyzję podejmuje Operator.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: CodeStudio
        ▼
ŚRODOWISKO: CODESTUDIO ── boczna nawigacja modułów
        │        (Workspace · Roundtable · Design · Terminal ·
        │         Developer · Diagnostics · Apps · Agents)
        ▼
MODUŁ: DIAGNOSTICS ──────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Diagnostics Center ·
  Logs Viewer · Errors Panel · Recommendations Panel ·
  Observability Tools
        │
        ▼
KARTA SESJI — materiał diagnostyczny właściwy sesji, własny kontekst
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 7.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | CodeStudio (wyłącznie) |
| Środowiska TalkIn i WorkSpace | Moduł niedostępny — diagnostyka techniczna należy wyłącznie do trybu programistycznego |
| Komponent własny | Diagnostics nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 7 (łącznie z Chat Window i Execution Loop Window) |
| Typ pracy | Sesyjna, dochodzeniowa — od materiału źródłowego, przez agregację, do rekomendacji |

---

## 2. Katalog funkcji i narzędzi

Każda pozycja: **nazwa** — co robi · zależności (biblioteki, formaty, integracje). Biblioteki backendu podano w wariancie Go (rdzeń platformy); telemetria opiera się na standardzie **OpenTelemetry** (OTel) jako wspólnym kontrakcie sygnałów: logów, metryk i śladów. Jeden kontrakt sygnałów obsługuje logi, błędy, metryki, ślady, koszt i alerty bez powielania instrumentacji.

### 2.1. Agregacja i przeszukiwanie logów (Logs Viewer)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Strumień logów na żywo (WebSocket)** | Chronologiczny strumień z kolorowaniem wg poziomu, auto-przewijanie z pauzą | kanał WebSocket `proces_sesji`; bufor pierścieniowy w kliencie |
| **Wyszukiwanie pełnotekstowe i regex (Grep)** | Wzorzec w treści logów z podświetleniem i licznikiem trafień, RE2 | indeks pełnotekstowy `blevesearch/bleve`; RE2 (`regexp`); przeszukiwanie archiwum plikowego |
| **Parsowanie strukturalne** | Rozbicie wpisu JSON/logfmt na pola, filtrowanie i sortowanie po polu (np. `status=500`) | `encoding/json`, parser logfmt; schemat pola z próbki |
| **Filtry: poziom / źródło / zakres czasu** | Debugowanie…krytyczny, proces/moduł/karta sesji/komponent, od-do z szybkimi skrótami | zapytanie po `proces_sesji` i `karta_sesji`; predykaty pól |
| **Deduplikacja z licznikiem powtórzeń** | Grupowanie identycznych wpisów, licznik ×N w oknie czasu | odcisk treści (hash znormalizowany); okno agregacji |
| **Korelacja po `trace_id` / `session_id`** | Zebranie wszystkich wpisów jednego wywołania lub sesji w jeden wątek | pole korelacyjne OTel `trace_id`; powiązanie z Provenance Explorer (2.4) |
| **Podział widoku dwóch źródeł lub okresów** | Porównanie logów dwóch komponentów albo stanu przed wdrożeniem i po nim w sąsiadujących kolumnach | dwa niezależne strumienie; wspólna oś czasu |
| **Zakładki i przypięcia wpisów** | Oznaczanie interesujących linii do szybkiego powrotu | przechowanie w `ustawienie` zasięgu sesja |
| **Archiwum i retencja** | Dostęp do logów spoza bufora bieżącej sesji, polityka retencji per zasięg | magazyn logów (pliki, magazyn obiektowy); polityka w `ustawienie` |
| **Eksport zakresu** | Zapis widocznego zakresu do pliku `.log`, `.jsonl` lub CSV | serializatory; formaty `.log`/JSON Lines/CSV |
| **Przekazanie fragmentu** | Zaznaczenie → Chat Window (kontekst) albo nowy wpis w Errors Panel | `wiadomosc`, `powiazanie_komponentu` |

### 2.2. Rejestr i grupowanie błędów (Errors Panel)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Tabela błędów z odciskiem (fingerprint)** | Grupowanie podobnych błędów w jeden wpis z licznikiem wystąpień | normalizacja stosu i hash odcisku; agregacja po kluczu |
| **Pełny stos wywołań ze źródłem** | Rozwinięty stos wywołań, mapowanie ramki do pliku i linii, kontekst wystąpienia | parser stosu per język; mapy źródeł dla JS |
| **Cykl życia statusu** | Nowy → w analizie → rozwiązany → ignorowany, przygaszanie rozwiązanych | model błędu; przejścia bez modali-bramek |
| **Priorytet niezależny od poziomu logu** | Ręczne nadanie: krytyczny, wysoki, średni, niski | plakietka priorytetu; sortowanie |
| **Wykres częstotliwości w czasie** | Trend liczby wystąpień — narasta, stały, wygasa | szereg czasowy zdarzeń; render sparkline i słupków po stronie klienta |
| **Wykrywanie regresji** | Sygnalizacja powrotu błędu oznaczonego jako rozwiązany | porównanie statusu z bieżącym strumieniem; znacznik regresji |
| **Powiązanie z wdrożeniem i commitem** | Skojarzenie błędu ze zbieżną czasowo zmianą (Git Panel modułu Developer) | `powiazanie_komponentu` Diagnostics↔Developer; oś czasu wdrożeń |
| **Notatki własne Operatora** | Adnotacje inline widoczne przy kolejnych analizach | pole edytowalne, autozapis; zasięg sesja lub projekt |
| **Alert „nowy typ błędu"** | Powiadomienie o pierwszym wystąpieniu nieznanego odcisku | reguła w silniku alertów (2.8) |
| **Eksport listy, w tym odfiltrowanej** | Zapis błędów do pliku raportu | serializator CSV/JSON/Markdown |

### 2.3. Zagregowany obraz stanu i pulpit (Diagnostics Center)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Panel stanu zagregowanego** | Liczba aktywnych błędów, ostrzeżeń, status ogólny kondycji | agregat z Errors Panel oraz kontroli stanu (2.7) |
| **Siatka wskaźników kondycji** | Plakietki stanu per komponent, usługa i model: sukces, ostrzeżenie, błąd | sondy zdrowia (2.7); kliknięcie → filtr Logs Viewer i Errors Panel |
| **Oś czasu zdarzeń** | Chronologiczne zestawienie zdarzeń z logów, błędów i wdrożeń | scalony strumień zdarzeń; markery wdrożeń |
| **Uruchomienie pełnej analizy** | Agregacja bieżącego materiału w spójny wniosek i wygenerowanie rekomendacji | `zadanie` + `kanal_modelu`; wynik → Recommendations Panel |
| **Porównanie z migawką** | Zestawienie stanu bieżącego ze stanem sprzed wdrożenia lub incydentu | zapis migawki stanu; różnica wskaźników |
| **Widok wg komponentu, modułu lub środowiska** | Przełączenie perspektywy agregacji | wymiary (labels) telemetrii OTel |
| **Powiadomienie o zdarzeniu krytycznym** | Sygnał widoczny niezależnie od przeglądanej sekcji | kanał zdarzeń krytycznych; próg w `ustawienie` |
| **Eksport raportu diagnostycznego** | Zwięzłe podsumowanie stanu do pliku | szablon raportu; Markdown/HTML/PDF |
| **Pulpity konfigurowalne (widgety)** | Układanie własnych kafli (metryka, log, błąd, koszt) w siatce | definicja pulpitu w `ustawienie`; siatka kafli po stronie klienta |

### 2.4. Prowenancja wywołań AI (Provenance Explorer)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rejestr wywołań kanału modelu** | Lista każdego wywołania: kanał, model, sesja, czas, tokeny wejścia i wyjścia, koszt, status, opóźnienie | przechwyt na warstwie `kanal_modelu`; magazyn śladów |
| **Drzewo śladu wywołania** | Hierarchia span: prompt → narzędzia → podagenci → odpowiedź, z czasami cząstkowymi | OTel Trace API (`go.opentelemetry.io/otel`); format OTLP; span dla wywołań AI |
| **Podgląd promptu i odpowiedzi** | Pełna treść wejścia i wyjścia z podświetleniem, w tym wiadomości systemowe i kontekst | render `wiadomosc`; redakcja danych wrażliwych sterowana ustawieniem (rozdz. 6) |
| **Rozbicie kosztu i tokenów wywołania** | Tokeny promptu, uzupełnienia i cache, koszt wg cennika kanału, udział narzędzi | liczenie tokenów (`tiktoken-go`, licznik dostawcy); cennik z konfiguracji kanału |
| **Trafienia cache promptu** | Oznaczenie odczytów i zapisów cache oraz oszczędności tokenów | metadane cache z odpowiedzi dostawcy; agregacja |
| **Łańcuch agenta i podagentów** | Powiązanie wywołań w ramach jednego zadania agentowego lub roli MultitaskingAI | korelacja `zadanie` i roli; zgodność ze **SPECYFIKACJA-AGENTOW** |
| **Powtórzenie wywołania (replay)** | Ponowne wykonanie tego samego wywołania z podglądem różnicy odpowiedzi | ponowne zlecenie przez `kanal_modelu`; różnica odpowiedzi |
| **Ocena i adnotacja jakości** | Ręczne oznaczenie trafności odpowiedzi, znacznik do dalszej analizy | pole oceny; zasięg sesja lub projekt |
| **Filtr i wyszukiwanie po wywołaniach** | Po modelu, sesji, koszcie, opóźnieniu, statusie, treści promptu | indeks śladów; predykaty pól |
| **Wykrycie nieudanych i wolnych wywołań** | Lista przekroczeń limitu czasu, błędów dostawcy, wywołań powyżej progu opóźnienia | progi w `ustawienie`; powiązanie z Errors Panel |
| **Eksport śladów (OTLP/JSON)** | Zapis wybranych śladów do pliku lub zewnętrznego zbieracza | format OTLP/JSON; eksporter OTel |

### 2.5. Zużycie kont, kosztów i budżetów (Usage & Cost)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Zużycie per kanał modelu i dostawca** | Tokeny, liczba żądań i koszt w podziale na kanały i dostawców, w czasie | agregacja śladów (2.4) po wymiarze `kanal_modelu`; szereg czasowy |
| **Zużycie per konto i token** | Rozliczenie na konkretne dane dostępowe (`profil_izolacji` — token per sesja) | wymiar konta z `profil_izolacji`; agregacja kosztu |
| **Zużycie per sesja, projekt i środowisko** | Który kontekst pracy zużył ile tokenów i środków | wymiary `karta_sesji`, projekt, środowisko; grupowanie |
| **Limity i budżety** | Miękkie progi zużycia z ostrzeżeniem przy zbliżaniu się i przekroczeniu, bez blokady | progi w `ustawienie`; reguła → silnik alertów (2.8) |
| **Prognoza zużycia** | Ekstrapolacja bieżącego tempa do końca okresu rozliczeniowego | regresja liniowa i EWMA (`montanaflynn/stats`) |
| **Rozbicie kosztu na operacje** | Udział promptu, uzupełnienia, cache i narzędzi w koszcie okresu | agregacja składników z 2.4 |
| **Porównanie modeli pod kątem kosztu i efektu** | Zestawienie kosztu i opóźnienia alternatywnych kanałów dla podobnych zadań | agregat śladów; zgodność ze strategią modeli platformy |
| **Wykrycie anomalii kosztowej** | Sygnał nagłego skoku zużycia względem historii | detektor anomalii (odchylenie od EWMA); alert |
| **Raport rozliczeniowy** | Eksport zużycia i kosztu do CSV/JSON/PDF za wybrany okres | serializatory; szablon raportu |

> **Uwaga o granicy.** Usage & Cost mierzy i prognozuje **zużycie techniczne** (tokeny, żądania, szacowany koszt wg cennika kanału). Nie jest systemem księgowym i nie realizuje operacji finansowych — nie inicjuje płatności, nie zmienia planów u dostawców, nie przenosi środków. Progi budżetowe **ostrzegają**, nigdy nie wstrzymują pracy platformy siłą (zero blokad).

### 2.6. Metryki, wydajność i profilowanie (Metrics & Performance)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Metryki wydajności na żywo** | Czas odpowiedzi, przepustowość, procesor i pamięć procesów sesji, opóźnienia komponentów | OTel Metrics, `prometheus/client_golang`; źródło `proces_sesji` |
| **Szeregi czasowe i percentyle** | Wykresy trendu, p50/p90/p95/p99, agregacja po oknie | magazyn szeregów czasowych (`prometheus/tsdb`, VictoriaMetrics); histogramy |
| **Mini-wykresy z rozwinięciem** | Kompaktowy trend każdego wskaźnika, kliknięcie → pełny wykres | render po stronie klienta; dane z magazynu szeregów czasowych |
| **Definiowalne zapytania metryk** | Własne agregaty i formuły nad metrykami (rate, sum, avg by label) | silnik zapytań zgodny z PromQL; parser wyrażeń |
| **Profilowanie procesora i pamięci (pprof)** | Zrzut i podgląd profilu procesu jako wykres płomieniowy | `net/http/pprof`, `google/pprof`; render wykresu płomieniowego po stronie klienta |
| **Profilowanie ciągłe** | Ciągłe próbkowanie profilu z historią i porównaniem okresów, sterowane ustawieniem `diagnostics.metryki.profilowanie` | agent próbkujący; magazyn profili |
| **Korelacja metryka ↔ zdarzenie** | Nałożenie markerów wdrożeń i incydentów na wykres metryki | wspólna oś czasu ze zdarzeniami (2.3) |
| **Detekcja anomalii wskaźnika** | Sygnał odchylenia metryki od wzorca sezonowego i historii | EWMA, odchylenie standardowe (`montanaflynn/stats`); próg → alert |
| **Eksport metryk (OTLP/Prometheus)** | Udostępnienie metryk zewnętrznemu zbieraczowi | eksporter OTel i Prometheus; format OTLP |

### 2.7. Kontrola stanu i dostępność (Health & Uptime)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Sondy zdrowia** | Okresowe sprawdzenie kondycji komponentu, usługi lub modelu (ping, punkt końcowy, zapytanie) | harmonogram sond (`zadanie`, Scheduler); sonda HTTP/TCP (`net/http`) |
| **Monitor dostępności** | Historia dostępności, procent dostępności, czas ostatniej niedostępności per usługa | seria stanów sond; agregat SLA/SLO |
| **Monitoring syntetyczny** | Odtwarzalny scenariusz wywołania (np. testowy prompt do modelu) jako sonda | scenariusz jako `zadanie`; wykonanie przez warstwę Terminal lub HTTP |
| **Budżet błędów (SLO)** | Cel dostępności i zużycie budżetu błędów w okresie | definicja SLO w `ustawienie`; wyliczenie z serii sond |
| **Historia i oś incydentów** | Rejestr okresów degradacji z czasem trwania i zasięgiem | log stanów; powiązanie z Errors Panel |
| **Mapa zależności komponentów** | Graf zależności usług i propagacji stanu (błąd usługi A → ryzyko B) | graf `powiazanie_komponentu`; render po stronie klienta |
| **Kontrola stanu samej platformy** | Kondycja procesów sesji, kolejek, kanałów WebSocket, silnika modeli | telemetria `proces_sesji`; metryki wewnętrzne |

### 2.8. Alerty i powiadomienia (Alerts)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Reguły progowe** | Alert, gdy metryka, koszt lub liczba błędów przekroczy próg w oknie czasu | ewaluator reguł nad magazynem szeregów czasowych i strumieniem; próg w `ustawienie` |
| **Reguły anomalii** | Alert na odchylenie od wzorca bez sztywnego progu | detektor (EWMA, odchylenie) współdzielony z 2.5 i 2.6 |
| **Kanały powiadomień** | Dostarczenie alertu: powiadomienie w aplikacji, e-mail, komunikator, webhook | rozszerzenia i integracje (Model konfiguracji, rozdz. 5.7); format webhook |
| **Always On Display** | Alert krytyczny wypychany do globalnej funkcji nadzoru nad pracą nienadzorowaną | powiązanie z Always On Display (Koncepcja, rozdz. 12.2) |
| **Wyciszenia i okna serwisowe** | Czasowe wstrzymanie alertów, grupowanie duplikatów | reguły wyciszeń w `ustawienie` |
| **Eskalacja i potwierdzanie** | Ścieżka eskalacji przy braku reakcji, potwierdzenie obsłużenia | stan alertu; harmonogram eskalacji |
| **Historia alertów** | Rejestr wyzwoleń, czasu reakcji i rozwiązania | log alertów; powiązanie z incydentami (2.7) |

### 2.9. Rekomendacje i operacje kontekstowe (Recommendations Panel)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rekomendacje naprawcze z priorytetem** | Lista sugestii Wykonawcy: tytuł, opis, źródło, pewność, priorytet | `kanal_modelu`; kontekst z Logs Viewer, Errors Panel, Metrics & Performance, Provenance Explorer |
| **Podgląd poprawki** | Fragment kodu, polecenie powłoki lub zmiana ustawienia w bloku `.dn-kod` | render różnicy; `hexops/gotextdiff` |
| **Zastosuj poprawkę → Developer** | Otwarcie Code Editor z gotową zmianą do przeglądu; tryb bezpośredni sterowany ustawieniem | `powiazanie_komponentu` Diagnostics→Developer |
| **Analiza przyczyny źródłowej** | Wykonawca wiąże objawy z wielu okien w jedną hipotezę przyczyny | agregat wielosygnałowy: log, błąd, metryka, ślad |
| **Korelacja wieloźródłowa** | Automatyczne zszycie logu, błędu, metryki i śladu tego samego zdarzenia | wspólny `trace_id` i oś czasu |
| **Inne rozwiązanie** | Wygenerowanie odmiennego rozwiązania tego samego problemu | `kanal_modelu`; nowa karta rekomendacji |
| **Ocena skuteczności po wdrożeniu** | Obserwacja, czy błąd źródłowy powrócił w oknie obserwacji po zastosowaniu | powiązanie z Errors Panel; okno obserwacji |
| **Podsumowanie stanu na żądanie** | Zwięzły opis kondycji systemu w języku naturalnym | `kanal_modelu` + agregat Diagnostics Center |
| **Wyjaśnienie wpisu, błędu lub śladu** | Kontekst zaznaczonego materiału → wyjaśnienie i sugestia kroku | `kanal_modelu`; materiał z aktywnego okna |

### 2.10. Raporty, migawki i wymiana danych

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Migawka stanu** | Zamrożenie pełnego obrazu diagnostycznego do porównania i archiwum | serializacja stanu; katalog roboczy sesji |
| **Raport diagnostyczny i powypadkowy** | Złożenie logów, błędów, metryk, kosztu i rekomendacji w jeden dokument | szablon; Markdown/HTML/PDF |
| **Eksport uniwersalny** | Logi (`.jsonl`), metryki (Prometheus/OTLP), ślady (OTLP), błędy i koszt (CSV/JSON) | eksportery OTel; serializatory formatów |
| **Import telemetrii zewnętrznej** | Wczytanie logów i śladów z zewnętrznego systemu jako materiał analizy | parsery OTLP/JSON; zasięg izolacji sieciowej |
| **Współdzielenie raportu z zespołem** | Przekazanie raportu jako artefaktu (np. do Library) | `powiazanie_komponentu`; artefakt pliku |

### 2.11. Zestawienie zależności bibliotecznych (backend Go)

| Obszar | Rozwiązanie wiodące | Rozwiązanie równoległe | Uwaga |
|---|---|---|---|
| Kontrakt telemetrii | **OpenTelemetry Go** (`go.opentelemetry.io/otel`) | natywne kolektory per sygnał | jeden standard dla logów, metryk i śladów; eksport OTLP |
| Indeks logów pełnotekstowy | `blevesearch/bleve` (czysty Go) | przeszukiwanie plików archiwum | bleve dla bufora i indeksu; przeszukiwanie plikowe dla surowych logów |
| Regex wyszukiwania | RE2 (`regexp`) | — | bezpieczny, bez katastrofalnego nawracania |
| Magazyn metryk | `prometheus/tsdb` (osadzony) | VictoriaMetrics | wariant osadzony bez zależności zewnętrznej; VictoriaMetrics przy dużym wolumenie |
| Klient metryk | `prometheus/client_golang` | OTel Metrics SDK | zależnie od modelu eksportu |
| Ślady | OTel Trace z magazynem śladów | — | OTLP jako format transportu i eksportu |
| Liczenie tokenów | `tiktoken-go` | licznik zwracany przez dostawcę | dla modeli spoza rodziny — licznik dostawcy z odpowiedzi |
| Profilowanie | `net/http/pprof` + `google/pprof` | agent profilowania ciągłego | pprof natywny dla Go |
| Statystyka i anomalie | `montanaflynn/stats` | własne EWMA | percentyle, odchylenia, detekcja skoków |
| Różnica poprawki i stanu | `hexops/gotextdiff` | `sergi/go-diff` | render różnicy rekomendacji i porównań migawek |
| Sondy zdrowia | `net/http` (sonda HTTP/TCP) | biblioteki dostępności | proste sondy; scenariusze syntetyczne przez Terminal |
| Kanały alertów | rozszerzenia i integracje (Model konfiguracji, rozdz. 5.7) | webhook generyczny | e-mail i komunikatory jako integracje, nie rdzeń |

---

## 3. Komplet okien operacyjnych modułu

Układ modułu jest **wyłącznie pionowy — podział lewa–prawa**. Chat Window zajmuje lewą kolumnę o pełnej wysokości obszaru roboczego, Execution Loop Window otwiera się jako kolumna sąsiadująca, obszar roboczy modułu zajmuje kolumnę dominującą po prawej, a okna pomocnicze otwierają się jako kolejne rozszerzenia boczne po prawej stronie obszaru roboczego. Regulacji podlega wyłącznie szerokość kolumn.

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja (kanał Użytkownik ↔ Wykonawca) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy i podstawowy mechanizm sterowania procesami diagnostycznymi |
| 2 | Execution Loop Window | Komunikacja (kanał Koordynator ↔ Wykonawca) | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji po otwarciu kolumny pętli; znacznik `Pętla ▼` w pasku kontekstu | Prowadzenie pętli wykonawczej nad zadaniami diagnostycznymi: dekompozycja zlecenia, kolejka zadań, kontrola jakości, sterowanie przebiegiem |
| 3 | Diagnostics Center | Monitor procesu / pulpit | Kolumna dominująca (obszar roboczy), punkt wejścia | 1 | Widoczne bez interakcji | Agregacja spójnego obrazu stanu systemu |
| 4 | Logs Viewer | Monitor procesu (log na żywo) | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Logi`, skrót klawiszowy, polecenie w Chat Window | Przegląd logów systemowych |
| 5 | Errors Panel | Okno zarządcy (managera) | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `Błędy`, przycisk „Analizuj", polecenie w Chat Window | Rejestr zarejestrowanych błędów |
| 6 | Recommendations Panel | Podgląd i porównanie | Kolumna boczna wynikowa, otwierana jako rozszerzenie boczne | 2 | Otwiera się z wynikiem analizy; znacznik `Rekomendacje` | Sugerowane kroki naprawcze |
| 7 | Observability Tools | Kontener zakładek obserwowalności | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny | Provenance Explorer, Usage & Cost, Metrics & Performance, Health & Uptime, Alerts |

```
 Makieta zbiorcza — Moduł Diagnostics (stan spoczynku)   Dostępność: CodeStudio
 ═════════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window          │ Execution Loop     │ Diagnostics Center
  nawigacja  │ Użytkownik ↔         │ Koordynator ↔      │ zagregowany obraz
  modułów    │ Wykonawca            │ Wykonawca          │ stanu systemu
             │                      │                    │
  [CodeStudio]│ historia rozmowy    │ zlecenie i jego    │ stan ogólny · błędy
  [Ubuntu]   │ i strumień           │ dekompozycja       │ wskaźniki kondycji
  [Fable 5]  │ odpowiedzi           │ kolejka zadań      │ oś czasu zdarzeń
  [Ultra]    │                      │ przebieg pętli     │ wskaźniki wydajności
             │                      │                    │
             │ [Agent ▼] [⋮]        │ [Sterowanie ▼] [⋮] │ [Zakres ▼] [Operacje ▼]
             │ pole poleceń    [▶]  │                    │ [☰]
 ═════════════════════════════════════════════════════════════════════════════
   Logs Viewer, Errors Panel, Recommendations Panel i Observability Tools
   otwierają się jako kolejne kolumny boczne po prawej stronie obszaru roboczego.
```

### 3.1. Warstwy widoczności w module

Moduł stosuje regułę stopniowego ujawniania funkcjonalności: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. W stanie spoczynku interfejs prezentuje wyłącznie elementy warstwy 1 — Chat Window, Execution Loop Window, Diagnostics Center, pasek kontekstu oraz zwinięte wyzwalacze `▼`, `⋮` i `☰`.

| Warstwa | Zawartość w module Diagnostics | Sposób dostępu |
|---|---|---|
| 1 | Chat Window, Execution Loop Window, Diagnostics Center, pasek kontekstu, wskaźniki stanu wykonania i plakietka stanu ogólnego | Widoczne bez interakcji |
| 2 | Wybór kanału modelu i wykonawcy, zakres czasu, filtry poziomu i źródła, otwarcie Logs Viewer, Errors Panel i Recommendations Panel | Znacznik kontekstowy, ikona, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Zestawy akcji nad wpisem i błędem, ustawienia szybkie sesji, warianty eksportu, operacje na rekomendacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | **Znaczna część funkcji diagnostycznych — funkcje eksperckie:** Provenance Explorer, Usage & Cost, Metrics & Performance, Health & Uptime, Alerts, profilowanie pprof i ciągłe, definiowalne zapytania metryk, powtórzenie wywołania modelu, reguły anomalii, monitoring syntetyczny, budżet błędów SLO, import telemetrii zewnętrznej, redakcja danych wrażliwych w śladach | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

Każda ukryta funkcja pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Widoczność zakładek kontenera Observability Tools steruje Operator kluczami `diagnostics.narzedzia.*` (rozdz. 6.6); brak włączenia zakładki oznacza jej ukrycie, nie zablokowanie.

---

## 4. Specyfikacja okien operacyjnych

### 4.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja |
| Rola | Główne okno komunikacji Użytkownik ↔ Wykonawca; centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami diagnostycznymi modułu |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Warstwa | 1 |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 7.3) |

**Zawartość i pełny arsenał funkcji.**

- Przyjmowanie poleceń w języku naturalnym i prezentacja strumienia odpowiedzi oraz wyników; zatwierdzanie i przerywanie działań Wykonawcy; wyjaśnianie wyniku i kontekstu.
- Historia rozmowy z Wykonawcą w kontekście bieżącego materiału diagnostycznego — aktywnego filtra Logs Viewer, zaznaczonych pozycji Errors Panel i wybranego śladu w Provenance Explorer.
- Polecenia analityczne: „przeanalizuj ten log", „co powoduje ten błąd", „porównaj to zdarzenie z poprzednim tygodniem", „zaproponuj poprawkę", „pokaż koszt tej sesji".
- Odwołanie do zaznaczonego wpisu w Logs Viewer, Errors Panel lub Provenance Explorer jako przedmiotu zapytania, bez ręcznego kopiowania treści.
- Uruchomienie funkcji eksperckich warstwy 4 poleceniem języka naturalnego — otwarcie Provenance Explorer, Usage & Cost, Metrics & Performance, Health & Uptime i Alerts.
- Generowanie podsumowania stanu systemu na żądanie, niezależnie od pełnej analizy uruchamianej z Diagnostics Center.
- Zlecenie pętli wykonawczej — polecenie złożone przekazywane Koordynatorowi, którego przebieg prezentuje Execution Loop Window.
- Eskalacja do modułu Developer — przygotowanie poprawki i przekazanie jej do Code Editor, gdy powiązanie skonfigurowane.
- Uruchomienie polecenia diagnostycznego w module Terminal bezpośrednio z poziomu okna komunikacji, gdy powiązanie skonfigurowane.
- Cytowanie logów, komunikatów błędów i fragmentów śladów blokiem `.dn-kod`, krojem mono.
- Załączanie plików do wiadomości — zrzutu ekranu błędu, pliku zrzutu awaryjnego, eksportu logu zewnętrznego.
- Historia poleceń (nawigacja strzałkami), ponowne zadanie pytania z modyfikacją zakresu czasu lub filtra.
- Wskaźnik zakresu analizowanego materiału (liczba logów, błędów i śladów uwzględnionych w bieżącym kontekście).
- Pasek kontekstu ze znacznikami środowiska, repozytorium, kanału modelu i wykonawcy; kliknięcie znacznika otwiera odpowiedni selektor.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Chat Window ───────────────────┐
│ [Danaco Console] [Ubuntu]       │
│ [Fable 5] [Ultra]          [⋮]  │
├─────────────────────────────────┤
│ Historia rozmowy (przewijana)   │
│ ┌─ Użytkownik ────────────────┐ │
│ │ „dlaczego usługa płatności  │ │
│ │  zwraca błąd 500?"          │ │
│ └─────────────────────────────┘ │
│ ┌─ Wykonawca (strumień) ──────┐ │
│ │ Na podstawie logu z         │ │
│ │ 12:04:07 przyczyną jest…    │ │
│ │ ```                         │ │
│ │ TimeoutError: usługa        │ │
│ │ płatności                   │ │
│ │ ```                         │ │
│ │ [ Otwórz w Centrum ]        │ │
│ │ [ Poproś o poprawkę ]       │ │
│ └─────────────────────────────┘ │
├─────────────────────────────────┤
│ [📎] [Agent ▼]  pole poleceń    │
│ ……………………………………………     [ ▶ ] │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pasek kontekstu | Zestaw lekkich znaczników | Prezentuje środowisko, repozytorium, kanał modelu i wykonawcę | Mała linia znaczników `--dn-tekst-3` | 1 | Widoczny bez interakcji | domyślny · znacznik aktywny | Kliknięcie znacznika otwiera odpowiedni selektor (warstwa 2) | Nagłówek Chat Window |
| Wskaźnik zakresu materiału | Etykieta tekstowa w nagłówku | Pokazuje liczbę logów, błędów i śladów uwzględnionych w kontekście rozmowy | Mały tekst `--dn-tekst-3` z licznikiem | 1 | Widoczny bez interakcji | domyślny · aktualizowany przy zmianie filtra | Kliknięcie otwiera podgląd listy uwzględnionych pozycji | Nagłówek Chat Window |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie pytania lub polecenia diagnostycznego | Duże pole zamykające kolumnę, rozciągliwe w pionie | 1 | Widoczne bez interakcji | domyślny · fokus · z załącznikiem · błąd (pusta wysyłka) | Enter wysyła, Shift+Enter nowa linia | Koniec kolumny Chat Window |
| Selektor wykonawcy i kanału modelu | Element zwinięty `Agent ▼` | Wybór wykonawcy, kanału modelu i poziomu wysiłku | Mały element zwinięty | 2 | Kliknięcie znacznika `Agent ▼`; po wyborze element zwija się samoczynnie | zwinięty · rozwinięty | Ustawia wykonawcę dla bieżącej karty sesji | Pasek poleceń |
| Blok cytowanego logu, błędu lub śladu | `.dn-kod`, krój mono | Prezentacja przywołanego fragmentu logu, stosu wywołań lub span śladu | Średni blok w treści wiadomości | 1 | Widoczny w treści odpowiedzi | domyślny | Kliknięcie „Otwórz w Centrum" przenosi do materiału źródłowego | Treść wiadomości |
| Przycisk „Poproś o poprawkę" | Przycisk `--zloty`, mały | Skierowanie zapytania o konkretną rekomendację naprawczą | Mały przycisk CTA, akcent złoty | 1 | Widoczny pod odpowiedzią diagnostyczną | domyślny · ładowanie | Generuje wpis w Recommendations Panel powiązany z tym wątkiem rozmowy | Pod odpowiedzią Wykonawcy |
| Menu kebab sesji | Ikona `⋮` | Zestaw akcji i ustawień szybkich sesji | Mała ikona 36×36 px | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Otwiera menu: kanał modelu, zakres pamięci, izolacja techniczna, eksport rozmowy | Prawy róg nagłówka |
| Przycisk załącznika | Ikona (spinacz) | Dołączenie pliku (np. zrzutu awaryjnego) do wiadomości | Mała ikona | 2 | Kliknięcie ikony lub przeciągnięcie pliku | domyślny · aktywny (plik dołączony) | Otwiera okno wyboru pliku lub przyjmuje przeciągnięcie | Lewa strona pola poleceń |
| Wywołanie funkcji eksperckiej | Polecenie języka naturalnego | Otwarcie okna warstwy 4 (Provenance Explorer, Usage & Cost, Metrics & Performance, Health & Uptime, Alerts) | Brak elementu wizualnego | 4 | Polecenie w polu poleceń, skrót klawiszowy, wyszukiwarka funkcji | — | Otwiera wskazane okno jako kolumnę boczną | Pole poleceń |

---

### 4.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja pętli wykonawczej |
| Rola | Prowadzenie pętli wykonawczej nad zadaniami diagnostycznymi modułu: koordynacja zadań, nadzór nad realizacją, orkiestracja działań i kontrola realizacji procesów |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |
| Warstwa | 1 |
| Izolacja domyślna | Pętla właściwa karcie sesji diagnostycznej; zakres zadań zgodny z zakresem materiału sesji |

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie diagnostyczne i jego dekompozycja na zadania — pobranie zakresu logów, zbudowanie odcisków błędów, korelacja śladów po `trace_id`, wyliczenie metryk okresu, wypracowanie rekomendacji.
- Kolejka i stan zadań diagnostycznych: oczekujące, w toku, zakończone, ponowione, przerwane.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą — przydział zadania, raport z wykonania, żądanie uzupełnienia materiału.
- Wyniki kontroli jakości i decyzje o ponowieniu — Koordynator ponawia zadanie analizy, gdy materiał źródłowy okazał się niekompletny albo pewność rekomendacji jest niska.
- Wskaźniki przebiegu pętli: liczba zadań, czas trwania przebiegu, liczba ponowień, koszt wywołań modelu w ramach pętli.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia bez utraty dotychczasowego materiału.
- Powiązanie każdego zadania pętli z oknem źródłowym — Logs Viewer, Errors Panel, Provenance Explorer, Metrics & Performance — i przejście do tego materiału.
- Rejestr przebiegu pętli zasilający Errors Panel wpisem przy niepowodzeniu zadania oraz Provenance Explorer śladami wywołań modelu wykonanych w pętli.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Execution Loop Window ─────────┐
│ Koordynator ↔ Wykonawca    [⋮]  │
├─────────────────────────────────┤
│ Zlecenie: „ustal przyczynę      │
│ błędów usługi płatności"        │
├─────────────────────────────────┤
│ Zadania                         │
│  ● pobranie logów zakresu   ✓   │
│  ● odcisk błędów            ✓   │
│  ● korelacja śladów        ⟳    │
│  ○ wyliczenie metryk okresu     │
│  ○ wypracowanie rekomendacji    │
├─────────────────────────────────┤
│ Komunikaty sterujące            │
│  Koordynator → Wykonawca:       │
│   „uzupełnij materiał o ślady   │
│    kanału modelu"               │
├─────────────────────────────────┤
│ Przebieg: 3/5 · ponowienia 1    │
│ [ Sterowanie ▼ ]                │
└─────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok tekstowy | Prezentuje bieżące zlecenie diagnostyczne przekazane Koordynatorowi | Średni blok z tytułem zlecenia | 1 | Widoczny bez interakcji | domyślny · zlecenie skorygowane | Kliknięcie otwiera pełną treść zlecenia | Nagłówek kolumny |
| Lista zadań pętli | Lista pozycji ze znacznikiem stanu | Dekompozycja zlecenia na zadania diagnostyczne i ich stan | Średnia lista z ikonami stanu | 1 | Widoczna bez interakcji | oczekujące · w toku · zakończone · ponowione · przerwane | Kliknięcie zadania otwiera okno źródłowe z materiałem tego zadania | Środkowa część kolumny |
| Strumień komunikatów sterujących | Lista wiadomości Koordynator ↔ Wykonawca | Prezentacja wymiany sterującej w toku pętli | Średni blok przewijany | 1 | Widoczny bez interakcji | domyślny · nowy komunikat | Kliknięcie komunikatu rozwija pełną treść i powiązane zadanie | Dalsza część kolumny |
| Wskaźniki przebiegu pętli | Linia liczników | Liczba zadań, czas przebiegu, ponowienia, koszt wywołań modelu | Mała linia `--dn-tekst-3` | 1 | Widoczne bez interakcji | przebieg w toku · przebieg zakończony | Kliknięcie licznika kosztu otwiera Usage & Cost dla tej pętli | Stopka kolumny |
| Element zbiorczy „Sterowanie" | Element zwinięty `Sterowanie ▼` | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Mały element zwinięty | 2 | Kliknięcie `Sterowanie ▼`; po wyborze element zwija się samoczynnie | zwinięty · rozwinięty | Zmienia stan przebiegu pętli bez utraty materiału | Stopka kolumny |
| Wynik kontroli jakości | Plakietka przy zadaniu | Rezultat kontroli jakości zadania i decyzja o ponowieniu | Mała pigułka `.dn-plakietka` | 3 | Rozwinięcie pozycji zadania | przyjęty · ponowiony · odrzucony | Rozwija uzasadnienie decyzji Koordynatora | Pozycja zadania |
| Menu kebab pętli | Ikona `⋮` | Zestaw akcji nad przebiegiem: eksport przebiegu, powtórzenie pętli, przypięcie materiału | Mała ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Otwiera menu akcji przebiegu | Prawy róg nagłówka |
| Parametry pętli wykonawczej | Zestaw ustawień przebiegu | Limit ponowień, próg pewności rekomendacji, tryb autonomii Wykonawcy | Brak elementu wizualnego w stanie spoczynku | 4 | Polecenie języka naturalnego, tryb administracyjny, okno konfiguracji (rozdz. 6.7) | — | Ustawia parametry kolejnych przebiegów pętli | Konfiguracja modułu |

---

### 4.3. Diagnostics Center

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor procesu / pulpit |
| Waga wizualna | Kolumna dominująca obszaru roboczego, punkt wejścia integrujący pozostałe okna modułu |
| Warstwa | 1 |
| Izolacja domyślna | Agreguje materiał widoczny w ramach karty sesji; zakres źródeł konfigurowalny |

**Zawartość i pełny arsenał funkcji.**

- Panel zagregowanego stanu systemu: liczba aktywnych błędów, liczba ostrzeżeń, ogólny status kondycji.
- Wskaźniki kondycji per komponent, usługa i model — plakietki stanu: sukces, ostrzeżenie, błąd.
- Oś czasu zdarzeń — chronologiczne zestawienie zdarzeń pochodzących z Logs Viewer, Errors Panel oraz markerów wdrożeń w jednym widoku.
- Uruchomienie pełnej analizy — Wykonawca agreguje bieżący materiał w spójny obraz i generuje wniosek wraz z rekomendacjami; zlecenie złożone prowadzi Koordynator w Execution Loop Window.
- Wskaźniki wydajności: czas odpowiedzi, zużycie procesora i pamięci procesów sesji, opóźnienia komponentów.
- Mini-wykresy trendu dla każdego wskaźnika wydajności i kosztu, rozwijane do pełnego widoku.
- Filtrowanie widoku według zakresu czasu: ostatnia godzina, dzień, tydzień, zakres własny.
- Widok według komponentu, modułu platformy lub środowiska.
- Porównanie stanu bieżącego z migawką wcześniejszą, na przykład sprzed wdrożenia zmiany.
- Powiadomienie o nowym zdarzeniu krytycznym, widoczne niezależnie od aktualnie przeglądanej sekcji.
- Pulpity konfigurowalne — własne kafle metryki, logu, błędu i kosztu układane w siatce.
- Odnośniki szybkiego przejścia do Logs Viewer, Errors Panel, Recommendations Panel i zakładek Observability Tools z zachowanym kontekstem filtra.
- Eksport raportu diagnostycznego — zwięzłe podsumowanie stanu systemu do pliku.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Diagnostics Center ─────────────────────────────────┐
│ [ Zakres ▼ ]              [ Operacje ▼ ]        [☰]  │
├──────────────────────────────────────────────────────┤
│ Stan ogólny: ⚠ wymaga uwagi   Błędy: 3  Ostrzeżenia: 7│
│ ┌──────────────┐┌──────────────┐┌──────────────┐     │
│ │ usługa API   ││ baza danych  ││ kolejka zadań│     │
│ │ ● błąd       ││ ● sukces     ││ ● ostrzeżenie│     │
│ └──────────────┘└──────────────┘└──────────────┘     │
├──────────────────────────────────────────────────────┤
│ Oś czasu: 11:58 ● wdrożenie  12:04 ● błąd 500  12:07⚠│
├──────────────────────────────────────────────────────┤
│ Czas odpowiedzi ▂▃▇▆   CPU ▁▂▃▅   Pamięć ▂▂▃▃        │
│ Koszt okresu ▁▂▂▅                                    │
└──────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Plakietka stanu ogólnego | `.dn-plakietka--stan` | Jednozdaniowe podsumowanie kondycji systemu | Średnia pigułka z ikoną i etykietą | 1 | Widoczna bez interakcji | dobry stan · wymaga uwagi · krytyczny | Informacyjna, prowadzi wzrok do dalszej analizy | Nagłówek okna |
| Karta wskaźnika kondycji | Mała karta `.dn-karta` z kropką stanu | Reprezentuje stan pojedynczego komponentu, usługi lub modelu | Mała karta w siatce | 1 | Widoczna bez interakcji | sukces (zielona) · ostrzeżenie (żółta) · błąd (czerwona) | Kliknięcie filtruje Logs Viewer i Errors Panel do tego komponentu | Siatka wskaźników kondycji |
| Oś czasu zdarzeń | Lista znaczników chronologicznych | Prezentacja kolejności zdarzeń w wybranym zakresie czasu | Średni pas przewijany | 1 | Widoczna bez interakcji | pusta · z pozycjami | Kliknięcie znacznika otwiera szczegóły zdarzenia w Logs Viewer lub Errors Panel | Środkowa część okna |
| Mini-wykres wskaźnika | Mały wykres sparkline | Trend czasu odpowiedzi, zużycia procesora, pamięci i kosztu | Bardzo mały wykres liniowy z etykietą | 1 | Widoczny bez interakcji | aktualizowany na żywo | Kliknięcie rozwija pełny wykres w Metrics & Performance | Dalsza część kolumny |
| Selektor zakresu czasu | Element zwinięty `Zakres ▼` | Zawężenie agregowanego materiału do wskazanego okresu | Mały element zwinięty | 2 | Kliknięcie `Zakres ▼`; po wyborze zwija się samoczynnie | domyślny (ostatnia godzina) · rozwinięty · zakres własny | Przelicza wszystkie wskaźniki i oś czasu do wybranego okresu | Nagłówek okna |
| Element zbiorczy „Operacje" | Element zwinięty `Operacje ▼` | Uruchomienie pełnej analizy, porównanie z migawką, eksport raportu | Mały element zbiorczy z akcentem złotym dla analizy | 2 | Kliknięcie `Operacje ▼` | zwinięty · rozwinięty · ładowanie (analiza w toku) | Wynik analizy trafia do Recommendations Panel, podsumowanie do Chat Window, przebieg do Execution Loop Window | Nagłówek okna |
| Menu hamburger widoku | Ikona `☰` | Widok wg komponentu, modułu i środowiska; układ kafli pulpitu | Mała ikona | 3 | Kliknięcie `☰` | zwinięte · rozwinięte | Przełącza perspektywę agregacji i układ kafli | Nagłówek okna |
| Powiadomienie o zdarzeniu krytycznym | Pasek sygnalizacyjny | Sygnał zdarzenia krytycznego niezależny od przeglądanej sekcji | Wąski pas sygnalizacyjny w nagłówku | 1 | Pojawia się po przekroczeniu progu `diagnostics.center.prog_krytyczny` | ukryty · widoczny | Kliknięcie otwiera zdarzenie źródłowe | Nagłówek okna |
| Definicja pulpitu konfigurowalnego | Zestaw kafli własnych | Budowa własnego pulpitu z metryk, logów, błędów i kosztu | Brak elementu wizualnego w stanie spoczynku | 4 | Polecenie języka naturalnego, tryb administracyjny, konfiguracja modułu | — | Zapisuje definicję pulpitu w `ustawienie` | Konfiguracja modułu |

---

### 4.4. Logs Viewer

| Aspekt | Wartość |
|---|---|
| Typologia | Monitor procesu (log na żywo) |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne; materiał źródłowy dla Diagnostics Center |
| Warstwa | 2 |
| Izolacja domyślna | Strumień logów właściwy zakresowi wskazanemu w konfiguracji (proces sesji, moduł, środowisko) |

**Zawartość i pełny arsenał funkcji.**

- Strumień logów na żywo kanałem WebSocket, prezentowany chronologicznie.
- Filtrowanie po poziomie: debugowanie, informacja, ostrzeżenie, błąd, krytyczny.
- Filtrowanie po źródle: proces, moduł, karta sesji, komponent.
- Wyszukiwanie pełnotekstowe (Grep) z obsługą wyrażeń regularnych RE2, licznikiem i podświetleniem trafień.
- Parsowanie strukturalne wpisów JSON i logfmt, filtrowanie oraz sortowanie po polu.
- Korelacja po `trace_id` i `session_id` — zebranie wpisów jednego wywołania lub sesji w jeden wątek, przejście do Provenance Explorer.
- Zakres czasu od-do, z szybkimi skrótami: ostatnie piętnaście minut, godzina, dzień.
- Podświetlanie i kolorowanie wpisów według poziomu ważności.
- Znaczniki czasu przy każdej linii, ze zmianą strefy czasowej prezentacji.
- Zawijanie linii albo przewijanie poziome (przełącznik).
- Przypinanie interesujących wpisów jako zakładek do szybkiego powrotu.
- Grupowanie powtarzających się wpisów (deduplikacja) z licznikiem powtórzeń.
- Automatyczne przewijanie do najnowszego wpisu, wstrzymywane podczas przeglądania historii.
- Podział widoku — porównanie logów z dwóch źródeł lub okresów w sąsiadujących kolumnach.
- Limit buforowania wpisów, dostęp do archiwum starszych logów i polityka retencji per zasięg.
- Eksport wybranego zakresu logów do pliku `.log`, `.jsonl` lub CSV.
- Przekazanie zaznaczonego fragmentu do Chat Window albo utworzenie z niego nowego wpisu w Errors Panel.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Logs Viewer ────────────────────────────────────────┐
│ [ Poziom ▼ ] [ Źródło ▼ ] [ Zakres ▼ ]  🔍     [⋮]   │
├──────────────────────────────────────────────────────┤
│ 12:03:58  INFO   api   żądanie POST /platnosci       │
│ 12:04:02  WARN   api   opóźniona odpowiedź usługi    │
│ 12:04:07  ERROR  api   TimeoutError: usługa płatności│
│           ×12 (powtórzenie w ciągu 2 minut)          │
│ 12:04:09  INFO   api   ponowienie żądania            │
├──────────────────────────────────────────────────────┤
│ trace_id: 7f3a… · bufor 5000 wpisów                  │
└──────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Linia logu | Wiersz tekstu mono ze znacznikiem czasu | Pojedynczy wpis dziennika zdarzeń | Tekst `--dn-ff-mono`, kolor zależny od poziomu | 1 | Widoczna bez interakcji | debugowanie · informacja · ostrzeżenie (żółte) · błąd (czerwone) · krytyczny (czerwone, wytłuszczone) | Kliknięcie rozwija pełny kontekst wpisu | Ciało konsoli |
| Licznik powtórzeń | Mała plakietka pod zdeduplikowaną linią | Sygnalizuje, ile razy podobny wpis wystąpił w oknie czasu | Mała pigułka z liczbą | 1 | Widoczny przy ≥ 2 wystąpieniach | ukryty · widoczny | Kliknięcie rozwija listę wszystkich wystąpień ze znacznikami czasu | Pod linią zdeduplikowaną |
| Filtr poziomu | Element zwinięty `Poziom ▼` | Zawężenie strumienia do wybranych poziomów ważności | Mały element zwinięty | 2 | Kliknięcie `Poziom ▼`; po wyborze zwija się samoczynnie | domyślny (wszystkie) · zawężony | Przelicza widoczny strumień logów | Nagłówek okna |
| Filtr źródła | Element zwinięty `Źródło ▼` | Zawężenie strumienia do komponentu, modułu lub karty sesji | Mały element zwinięty | 2 | Kliknięcie `Źródło ▼` | domyślny (wszystkie) · zawężony | Przelicza widoczny strumień logów | Nagłówek okna |
| Pole Grep | Ikona wyszukiwania rozwijana do pola | Wyszukiwanie wzorca w treści logów | Ikona `🔍` rozwijana do pola `.dn-input` z przełącznikami regex i wielkości liter | 2 | Kliknięcie ikony lub skrót `Ctrl/Cmd + F` | zwinięte · puste · z wynikami · brak trafień | Enter uruchamia wyszukiwanie, podświetla i zlicza trafienia | Nagłówek okna |
| Znacznik `trace_id` | Etykieta korelacyjna | Powiązanie wpisu z wywołaniem modelu | Mała etykieta mono | 2 | Kliknięcie znacznika w stopce lub w rozwiniętym wpisie | domyślny · aktywny | Otwiera Provenance Explorer z tym śladem | Stopka okna, rozwinięty wpis |
| Menu kebab wpisu i widoku | Ikona `⋮` | Zestaw akcji: przypięcie wpisu, podział widoku, przełącznik auto-przewijania, zawijanie linii, eksport zakresu, przekazanie fragmentu | Mała ikona | 3 | Kliknięcie `⋮` lub menu kontekstowe na wpisie | zwinięte · rozwinięte | Wykonuje wybraną akcję nad wpisem albo widokiem | Nagłówek okna, wiersz wpisu |
| Podział widoku | Dwie niezależne kolumny przewijania | Porównanie logów dwóch źródeł lub okresów | Podział kolumny okna na dwie kolumny | 3 | Pozycja „Podziel" w menu `⋮` | wyłączony · włączony | Dzieli okno na dwie sąsiadujące kolumny przewijania | Ciało okna |
| Parsowanie strukturalne i filtr po polu | Zestaw predykatów pól | Rozbicie wpisu JSON/logfmt i filtrowanie po wartości pola | Brak elementu wizualnego w stanie spoczynku | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji, konfiguracja `diagnostics.logi.parsowanie` | — | Zawęża strumień do wpisów spełniających predykat | Konfiguracja modułu, Chat Window |
| Archiwum i retencja | Dostęp do logów spoza bufora | Odczyt materiału historycznego i polityka retencji | Brak elementu wizualnego w stanie spoczynku | 4 | Polecenie języka naturalnego, tryb administracyjny | — | Odczytuje archiwum na żądanie, zgodnie z `diagnostics.logi.retencja` | Konfiguracja modułu, Chat Window |

---

### 4.5. Errors Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządcy (managera) |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne; materiał źródłowy dla Diagnostics Center |
| Warstwa | 2 |
| Izolacja domyślna | Rejestr błędów właściwy zakresowi wskazanemu w konfiguracji |

**Zawartość i pełny arsenał funkcji.**

- Tabela zarejestrowanych błędów: znacznik czasu, typ błędu, komponent źródłowy, komunikat, status.
- Grupowanie podobnych błędów po odcisku (fingerprint) z licznikiem wystąpień, zamiast powielania identycznych wpisów.
- Szczegóły błędu: pełny stos wywołań ze wskazaniem pliku i linii, kontekst wystąpienia, powiązane linie z Logs Viewer, powiązany ślad wywołania modelu.
- Filtrowanie po statusie: nowy, w analizie, rozwiązany, ignorowany.
- Filtrowanie po komponencie źródłowym i poziomie ważności: krytyczny, wysoki, średni, niski.
- Wyszukiwanie po treści komunikatu błędu.
- Zmiana statusu błędu — oznaczenie jako rozwiązany, ignorowany lub w analizie, bez modali-bramek.
- Przypisanie priorytetu błędowi, niezależnie od jego pierwotnego poziomu ważności w logu.
- Wykres częstotliwości wystąpień błędu w czasie — czy problem narasta, jest stały, czy wygasa.
- Wykrywanie regresji — sygnalizacja powrotu błędu oznaczonego wcześniej jako rozwiązany.
- Alert przy pierwszym wystąpieniu nieznanego odcisku błędu.
- Notatki własne przy błędzie — adnotacje Operatora widoczne przy kolejnych analizach.
- Powiązanie błędu z konkretnym commitem lub wdrożeniem, gdy zbieżne czasowo (widoczne po stronie Git Panel modułu Developer).
- Przekazanie błędu do Diagnostics Center w celu agregacji, do Recommendations Panel po analizę, do Execution Loop Window jako zlecenie pętli albo do modułu Developer jako zadanie poprawki.
- Eksport listy błędów, w tym wyłącznie zakresu odfiltrowanego.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Errors Panel ───────────────────────────────────────┐
│ [ Status ▼ ] [ Komponent ▼ ] [ Priorytet ▼ ]    [⋮]  │
├──────────────────────────────────────────────────────┤
│ ✖ TimeoutError: usługa płatności   api    12:04:07   │
│   ×12 wystąpień · krytyczny · ↻ regresja        [⋮]  │
├──────────────────────────────────────────────────────┤
│ ⚠ Nieudana walidacja formularza  frontend  11:47 [⋮] │
├──────────────────────────────────────────────────────┤
│ ● Nowy odcisk błędu wykryty o 12:09                  │
└──────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Wiersz błędu (zgrupowany) | Wiersz `.dn-tabela` z ikoną i licznikiem | Reprezentuje jeden odcisk błędu i liczbę jego wystąpień | Średni wiersz z plakietką priorytetu | 1 | Widoczny bez interakcji | nowy · w analizie · rozwiązany (przygaszony) · ignorowany (przygaszony) · regresja | Kliknięcie rozwija szczegóły: stos wywołań, kontekst, notatki | Ciało tabeli |
| Plakietka priorytetu | Mała plakietka koloru | Sygnalizuje wagę błędu | Mała pigułka `.dn-plakietka` | 1 | Widoczna bez interakcji | krytyczny (czerwona) · wysoki (pomarańczowa) · średni (żółta) · niski (neutralna) | Kliknięcie otwiera listę wyboru priorytetu | Wiersz błędu |
| Znacznik regresji | Mała ikona `↻` przy wierszu | Sygnalizuje powrót błędu oznaczonego jako rozwiązany | Bardzo mała ikona | 1 | Pojawia się po wykryciu regresji | ukryty · widoczny | Kliknięcie otwiera historię statusów tego odcisku | Wiersz błędu |
| Filtry statusu, komponentu i priorytetu | Trzy elementy zwinięte `▼` | Zawężenie listy błędów wg wybranych kryteriów | Małe elementy zwinięte | 2 | Kliknięcie znacznika; po wyborze zwija się samoczynnie | domyślny · zawężony | Przelicza widoczną listę błędów | Nagłówek panelu |
| Panel szczegółów błędu | Rozwijany blok pod wierszem | Pełny stos wywołań, kontekst, powiązane logi i ślady, notatki, wykres częstotliwości | Średni panel, tekst mono dla stosu wywołań | 3 | Kliknięcie wiersza błędu | zwinięty · rozwinięty | Rozwija szczegóły; `Esc` zwija panel | Pod wierszem błędu |
| Menu kebab błędu | Ikona `⋮` przy wierszu | Zestaw akcji: „Analizuj", „Rozwiąż", „Ignoruj", zlecenie pętli, przekazanie do Developer, eksport | Mała ikona | 3 | Kliknięcie `⋮` lub menu kontekstowe wiersza | zwinięte · rozwinięte · ładowanie (analiza) | „Analizuj" otwiera Diagnostics Center z tym błędem jako materiałem wejściowym; zlecenie pętli otwiera Execution Loop Window | Wiersz błędu |
| Pole notatki | Małe pole tekstowe z ikoną dymka | Adnotacja własna Operatora przy błędzie | Małe pole edytowalne inline | 3 | Rozwinięcie panelu szczegółów | puste · z treścią | Zapisuje się automatycznie, widoczne przy kolejnych analizach tego błędu | Panel szczegółów błędu |
| Wykres częstotliwości | Mały wykres słupkowy lub liniowy | Trend liczby wystąpień błędu w czasie | Mały wykres w panelu szczegółów | 3 | Rozwinięcie panelu szczegółów | rosnący · malejący · stabilny | Najechanie pokazuje liczbę wystąpień w danym przedziale | Panel szczegółów błędu |
| Reguły normalizacji odcisku | Zestaw reguł grupowania | Sterowanie sposobem budowy odcisku błędu per język | Brak elementu wizualnego w stanie spoczynku | 4 | Tryb administracyjny, konfiguracja `diagnostics.bledy.reguly_odcisku` | — | Przebudowuje grupowanie wpisów rejestru | Konfiguracja modułu |
| Powiązanie z wdrożeniem i commitem | Odnośnik do Git Panel | Skojarzenie błędu ze zbieżną czasowo zmianą | Brak elementu wizualnego w stanie spoczynku | 4 | Polecenie języka naturalnego, konfiguracja `diagnostics.bledy.powiazanie_wdrozenia` | — | Otwiera powiązaną zmianę w Git Panel modułu Developer | Panel szczegółów błędu |

---

### 4.6. Recommendations Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Podgląd i porównanie |
| Waga wizualna | Kolumna boczna wynikowa, otwierana jako rozszerzenie boczne, aktualizowana po każdej analizie |
| Warstwa | 2 |
| Izolacja domyślna | Rekomendacje właściwe karcie sesji; historia zastosowanych rekomendacji narasta w toku sesji |

**Zawartość i pełny arsenał funkcji.**

- Lista rekomendacji wypracowanych przez Wykonawcę: tytuł, opis, powiązany błąd, log, metryka lub ślad źródłowy, poziom pewności.
- Priorytetyzacja rekomendacji: krytyczna, zalecana, uzupełniająca.
- Podgląd proponowanej poprawki — fragment kodu, polecenie powłoki lub zmiana konfiguracji — w bloku `.dn-kod`.
- Analiza przyczyny źródłowej — powiązanie objawów z wielu okien w jedną hipotezę przyczyny.
- Korelacja wieloźródłowa — zszycie logu, błędu, metryki i śladu tego samego zdarzenia po `trace_id` i osi czasu.
- Akcja „zastosuj poprawkę" — domyślnie przekazuje zmianę do modułu Developer i otwiera Code Editor z gotową zmianą do przeglądu; tryb ten jest konfigurowalny — Operator włącza bezpośrednie zastosowanie poprawki z pominięciem etapu przeglądu, zgodnie z zasadą pełnej konfigurowalności.
- Akcja „odrzuć rekomendację", z komentarzem doprecyzowującym wykorzystywanym do poprawy trafności kolejnych sugestii.
- Akcja „inne rozwiązanie" — wygenerowanie odmiennego rozwiązania dla tego samego problemu.
- Historia zastosowanych rekomendacji, z odnotowaniem, czy powiązany błąd ustąpił po wdrożeniu.
- Powiązanie każdej rekomendacji z oknem źródłowym — Logs Viewer, Errors Panel, Diagnostics Center, Provenance Explorer lub Metrics & Performance.
- Ocena skuteczności — po zastosowaniu poprawki panel oznacza, czy błąd źródłowy pojawił się ponownie w kolejnym oknie obserwacji.
- Eksport listy rekomendacji jako raportu do przekazania zespołowi.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Recommendations Panel ──────────────────────────────┐
│ [ Priorytet ▼ ]                                 [⋮]  │
├──────────────────────────────────────────────────────┤
│ 🔴 krytyczna · pewność: wysoka                       │
│ Dodaj limit czasu i ponowienie próby dla klienta     │
│ usługi płatności                                     │
│ źródło: Errors Panel — TimeoutError (×12)            │
│ ┌ podgląd poprawki ─────────────────────────────┐    │
│ │ client.timeout = 5000;                        │    │
│ │ client.retry = { attempts: 3 };               │    │
│ └───────────────────────────────────────────────┘    │
│ [ Zastosuj poprawkę ]                           [⋮]  │
├──────────────────────────────────────────────────────┤
│ 🟡 zalecana · pewność: średnia                       │
│ Zwiększ pulę połączeń do bazy danych            [⋮]  │
└──────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Karta rekomendacji | `.dn-karta` z plakietką priorytetu | Reprezentuje jedną sugestię naprawczą Wykonawcy | Średnia karta z tytułem, opisem i odnośnikiem do źródła | 1 | Widoczna bez interakcji | krytyczna (czerwona plakietka) · zalecana (żółta) · uzupełniająca (neutralna) · zastosowana (przygaszona z ikoną potwierdzenia) | Kliknięcie odnośnika źródła otwiera powiązany błąd, log lub ślad | Lista główna panelu |
| Wskaźnik pewności | Mała etykieta tekstowa | Informuje o poziomie pewności co do trafności rekomendacji | Mały tekst `--dn-tekst-3` z ikoną | 1 | Widoczny bez interakcji | wysoka · średnia · niska | Informacyjny, wpływa na kolejność domyślnego sortowania | Nagłówek karty rekomendacji |
| Blok podglądu poprawki | `.dn-kod`, krój mono | Prezentacja proponowanej zmiany przed zastosowaniem | Średni blok wewnątrz karty rekomendacji | 1 | Widoczny bez interakcji | domyślny | Kliknięcie rozwija pełną treść, gdy poprawka jest dłuższa niż podgląd | Ciało karty rekomendacji |
| Przycisk „Zastosuj poprawkę" | Przycisk `--zloty`, mały | Przekazanie poprawki do modułu Developer w celu przeglądu i wdrożenia | Mały przycisk CTA, akcent złoty — sygnalizuje wagę akcji, nie blokuje jej wykonania | 1 | Widoczny na karcie rekomendacji | domyślny · ładowanie | Otwiera Code Editor modułu Developer z gotową zmianą przygotowaną do przeglądu i zatwierdzenia | Stopka karty rekomendacji |
| Znacznik skuteczności | Mała ikona przy zastosowanej rekomendacji | Informuje, czy błąd źródłowy ustąpił po wdrożeniu poprawki | Bardzo mała ikona (ptaszek, ostrzeżenie) | 1 | Pojawia się po zakończeniu okna obserwacji | oczekuje obserwacji · potwierdzona skuteczność · błąd powrócił | Kliknięcie otwiera powiązany błąd w Errors Panel dla porównania | Karta zastosowanej rekomendacji |
| Filtr priorytetu | Element zwinięty `Priorytet ▼` | Zawężenie listy rekomendacji wg wagi | Mały element zwinięty | 2 | Kliknięcie `Priorytet ▼`; po wyborze zwija się samoczynnie | domyślny (wszystkie) · zawężony | Przelicza widoczną listę rekomendacji | Nagłówek panelu |
| Menu kebab rekomendacji | Ikona `⋮` | Zestaw akcji: „Odrzuć" z komentarzem, „Inne rozwiązanie", podgląd analizy przyczyny źródłowej, eksport raportu | Mała ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte · ładowanie | „Odrzuć" oznacza rekomendację jako odrzuconą i zasila kontekst kolejnych analiz; „Inne rozwiązanie" dodaje nową kartę dla tego samego źródła | Karta rekomendacji |
| Analiza przyczyny źródłowej | Blok wnioskowania wieloźródłowego | Hipoteza przyczyny zszyta z logu, błędu, metryki i śladu | Średni blok rozwijany | 3 | Pozycja w menu `⋮` karty | zwinięty · rozwinięty | Prezentuje łańcuch dowodowy z odnośnikami do materiału źródłowego | Karta rekomendacji |
| Tryb bezpośredniego zastosowania poprawki | Ustawienie konfiguracyjne przepływu | Pominięcie etapu przeglądu w Code Editor | Brak elementu wizualnego w stanie spoczynku | 4 | Tryb administracyjny, okno konfiguracji | — | Zmienia domyślny przepływ akcji „Zastosuj poprawkę" | Konfiguracja modułu |

---

### 4.7. Observability Tools (kontener zakładek, funkcje eksperckie)

| Aspekt | Wartość |
|---|---|
| Typologia | Kontener zakładek obserwowalności |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne po prawej stronie obszaru roboczego |
| Warstwa | 4 |
| Izolacja domyślna | Telemetria właściwa zakresowi karty sesji; zakres źródeł i retencja konfigurowalne |

Kontener zbiera pięć zakładek funkcji eksperckich. Użytkownik podstawowy nie widzi tych elementów; wywołuje je polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji albo tryb administracyjny. Widocznością zakładek steruje Operator kluczami `diagnostics.narzedzia.*`; brak włączenia zakładki oznacza jej ukrycie, nie zablokowanie.

| Zakładka | Zawartość | Warstwa | Sposób wywołania |
|---|---|---|---|
| **Provenance Explorer** | Rejestr wywołań kanału modelu, drzewo śladu, podgląd promptu i odpowiedzi, rozbicie kosztu i tokenów, trafienia cache, łańcuch agenta i podagentów, powtórzenie wywołania, ocena jakości, filtr po wywołaniach, wykrycie nieudanych i wolnych wywołań, eksport śladów (2.4) | 4 | Polecenie w Chat Window, znacznik `trace_id` w Logs Viewer, skrót klawiszowy, wyszukiwarka funkcji |
| **Usage & Cost** | Zużycie per kanał modelu, dostawca, konto, token, sesja, projekt i środowisko; limity i budżety, prognoza, rozbicie kosztu na operacje, porównanie modeli, wykrycie anomalii kosztowej, raport rozliczeniowy (2.5) | 4 | Polecenie w Chat Window, licznik kosztu w Execution Loop Window, wyszukiwarka funkcji |
| **Metrics & Performance** | Metryki na żywo, szeregi czasowe i percentyle, mini-wykresy z rozwinięciem, definiowalne zapytania metryk, profilowanie pprof i ciągłe, korelacja metryka ↔ zdarzenie, detekcja anomalii, eksport metryk (2.6) | 4 | Polecenie w Chat Window, mini-wykres w Diagnostics Center, tryb administracyjny |
| **Health & Uptime** | Sondy zdrowia, monitor dostępności, monitoring syntetyczny, budżet błędów SLO, historia i oś incydentów, mapa zależności komponentów, kontrola stanu samej platformy (2.7) | 4 | Polecenie w Chat Window, karta wskaźnika kondycji w Diagnostics Center, tryb administracyjny |
| **Alerts** | Reguły progowe i anomalii, kanały powiadomień, powiązanie z Always On Display, wyciszenia i okna serwisowe, eskalacja i potwierdzanie, historia alertów (2.8) | 4 | Polecenie w Chat Window, powiadomienie o zdarzeniu krytycznym, tryb administracyjny |

**Makieta tekstowa (stan spoczynku — kontener zamknięty).**

```
┌─ Diagnostics Center ────────┬─ Observability Tools ──┐
│ stan ogólny · wskaźniki     │ [Prowenancja] [Koszt]  │
│ kondycji · oś czasu         │ [Metryki] [Stan]       │
│ wskaźniki wydajności        │ [Alerty]          [⋮]  │
│                             ├────────────────────────┤
│ [ Zakres ▼ ] [ Operacje ▼ ] │ Wybierz zakładkę albo  │
│ [☰]                         │ wydaj polecenie w      │
│                             │ Chat Window            │
└─────────────────────────────┴────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pasek zakładek kontenera | Lista zakładek obserwowalności | Przełączanie między pięcioma zakładkami eksperckimi | Wąski pas zakładek w nagłówku kolumny | 4 | Otwarcie kontenera poleceniem, skrótem lub z wyszukiwarki funkcji | zakładka aktywna · zakładka ukryta konfiguracją | Przełącza zawartość kolumny bocznej | Nagłówek kolumny |
| Wiersz wywołania modelu | Wiersz `.dn-tabela` z kosztem i opóźnieniem | Reprezentuje jedno wywołanie `kanal_modelu` | Średni wiersz z licznikami tokenów i kosztu | 4 | Zakładka Provenance Explorer | zakończone · nieudane · wolne (powyżej progu) | Kliknięcie rozwija drzewo śladu z podglądem promptu i odpowiedzi | Provenance Explorer |
| Drzewo śladu | Hierarchia span | Przebieg wywołania: prompt → narzędzia → podagenci → odpowiedź | Średni blok hierarchiczny z czasami cząstkowymi | 4 | Rozwinięcie wiersza wywołania | zwinięte · rozwinięte | Kliknięcie span otwiera jego treść i metadane | Provenance Explorer |
| Kafel budżetu | Karta `.dn-karta` z progiem i prognozą | Zużycie względem miękkiego budżetu okresu | Średnia karta z paskiem postępu | 4 | Zakładka Usage & Cost | poniżej progu · próg ostrzeżenia · przekroczony | Sygnalizuje i rekomenduje, nie wstrzymuje pracy platformy | Usage & Cost |
| Wykres szeregu czasowego | Wykres liniowy z percentylami | Trend metryki w wybranym oknie z p50/p95/p99 | Średni wykres z osią czasu | 4 | Zakładka Metrics & Performance | dane dostępne · brak danych w oknie | Najechanie prezentuje wartości; markery wdrożeń nakładane na oś | Metrics & Performance |
| Wykres płomieniowy profilu | Render profilu pprof | Podgląd zużycia procesora i pamięci procesu | Duży wykres płomieniowy | 4 | Polecenie języka naturalnego, tryb administracyjny | zrzut na żądanie · profilowanie ciągłe | Kliknięcie ramki zawęża profil do wybranego poddrzewa | Metrics & Performance |
| Karta sondy zdrowia | Karta z historią stanów | Definicja sondy i historia dostępności komponentu | Średnia karta z paskiem stanów | 4 | Zakładka Health & Uptime | dostępny · degradacja · niedostępny | Kliknięcie otwiera oś incydentów i powiązane błędy | Health & Uptime |
| Reguła alertu | Wiersz reguły z progiem i kanałem | Definicja warunku wyzwolenia alertu | Średni wiersz z warunkiem i kanałem powiadomienia | 4 | Zakładka Alerts, tryb administracyjny | aktywna · wyciszona · w eskalacji | Zmiana reguły obowiązuje od kolejnej ewaluacji; alert krytyczny trafia do Always On Display | Alerts |
| Menu kebab kontenera | Ikona `⋮` | Zestaw akcji: eksport telemetrii, import telemetrii zewnętrznej, migawka stanu, redakcja danych wrażliwych | Mała ikona | 4 | Kliknięcie `⋮` w nagłówku kontenera | zwinięte · rozwinięte | Wykonuje wybraną operację ekspercką nad telemetrią | Nagłówek kolumny |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — od objawu do rekomendacji

```
 Logs Viewer / Errors Panel      Diagnostics Center        Recommendations Panel
 ───────────────────────         ──────────────────         ──────────────────────
 1. gromadzenie materiału
    źródłowego          ───────►
 2. uruchomienie
    pełnej analizy       ──────────────► agregacja logów,
                                          błędów, metryk i
                                          śladów w spójny
                                          obraz stanu systemu
                                                │
                                                ▼
 3. wniosek diagnostyczny  ◄────────────────────┘
    i wygenerowane
    rekomendacje                                         ─────────►  lista rekomendacji
                                                                       z priorytetem
                                                                       i podglądem poprawki
```

### 5.2. Przepływ pętli wykonawczej — Koordynator ↔ Wykonawca nad zadaniem diagnostycznym

```
Chat Window                     Execution Loop Window            Okna materiałowe
──────────────                  ─────────────────────            ─────────────────
 Użytkownik zleca:
 „ustal przyczynę błędów
  usługi płatności"
        │
        ▼
 przekazanie zlecenia  ────────► Koordynator dekomponuje
                                  zlecenie na zadania:
                                  · pobranie logów zakresu ─────► Logs Viewer
                                  · odcisk błędów          ─────► Errors Panel
                                  · korelacja śladów       ─────► Provenance Explorer
                                  · wyliczenie metryk      ─────► Metrics & Performance
                                          │
                                          ▼
                                  Wykonawca realizuje zadania,
                                  raportuje wynik każdego kroku
                                          │
                                          ▼
                                  kontrola jakości wyniku
                                  · materiał kompletny → dalej
                                  · materiał niepełny  → ponowienie
                                          │
                                          ▼
                                  wypracowanie rekomendacji ────► Recommendations Panel
        ◄─────────────────────────  podsumowanie przebiegu
        ▼
 Użytkownik zatwierdza,
 koryguje zlecenie
 albo przerywa przebieg
```

### 5.3. Przepływ zastosowania poprawki w module Developer

```
Recommendations Panel                        Developer › Code Editor / Git Panel
──────────────────────                       ────────────────────────────────────
 wybór rekomendacji
 [ Zastosuj poprawkę ]
        │
        ▼
 przekazanie proponowanej zmiany  ──────────►  otwarcie pliku z gotową zmianą
                                                do przeglądu
                                                       │
                                                       ▼
                                                przegląd i korekta
                                                       │
                                                       ▼
                                                zatwierdzenie zmiany (Git Panel)
        │
        ◄──────────────────────────────────────────────┘
        ▼
 oznaczenie rekomendacji jako zastosowanej,
 rozpoczęcie obserwacji skuteczności
```

### 5.4. Przepływ odtworzenia objawu w module Terminal

```
Diagnostics Center                          Terminal › Terminal Tabs
───────────────────                         ─────────────────────────
 hipoteza: usługa nie odpowiada
 na porcie wskazanym w konfiguracji
        │
        │  polecenie diagnostyczne wydane z Chat Window
        ▼
 „sprawdź, czy usługa nasłuchuje" ─────────►  uruchomienie polecenia
                                               weryfikującego w rzeczywistym
                                               środowisku wykonawczym
                                                       │
                                                       ▼
                                               wynik trafia z powrotem
        ◄──────────────────────────────────────  do Output Console i,
        ▼                                         na żądanie, do Logs Viewer
 potwierdzenie lub odrzucenie
 hipotezy diagnostycznej
```

### 5.5. Przepływ prowenancji i kosztu wywołania modelu

```
Logs Viewer                Provenance Explorer            Usage & Cost
────────────               ────────────────────           ─────────────
 wpis ze znacznikiem
 trace_id
        │
        ▼
 przejście po korelacji ──► drzewo śladu wywołania:
                             prompt → narzędzia →
                             podagenci → odpowiedź
                                     │
                                     │ tokeny, koszt, opóźnienie
                                     ▼
                             wykrycie wywołania wolnego
                             lub nieudanego  ─────────────► agregacja kosztu
                                     │                        per kanał, konto,
                                     ▼                        sesja i środowisko
                             wpis w Errors Panel                     │
                                                                     ▼
                                                            próg budżetu →
                                                            alert ostrzegawczy
```

### 5.6. Przepływ nadzoru nad procesem nienadzorowanym w MultitaskingAI

```
Środowisko MultitaskingAI                    Diagnostics
──────────────────────────                   ───────────
 rola Executor 3 / Validator we wcieleniu
 Security Auditor napotyka niepowodzenie
 procesu pracy ciągłej
        │
        ▼
 zdarzenie trafia do rejestru zdarzeń  ─────────────►  Errors Panel — nowy wpis
 procesu (Monitor procesu panelu                        z pełnym kontekstem procesu
 orkiestracji)                                                  │
                                                                  ▼
                                                         Execution Loop Window —
                                                         zlecenie analizy przyczyny
                                                                  │
                                                                  ▼
                                                         Diagnostics Center — analiza
                                                         w ramach nadzoru Always On
                                                         Display nad przebiegiem pętli
```

---

## 6. Punkty sterowania z okna konfiguracji

Wszystkie punkty są jawne, opatrzone objaśnieniem `[?]`, przechowywane jako `ustawienie` z kluczem (Model konfiguracji, rozdz. 5.3). Wartości podlegają zasadzie „brak ustawienia = wartość domyślna". Zero blokad — konfiguracja steruje **widocznością, zakresem i progami sygnalizacji**, nigdy nie wyłącza działania platformy siłą i nie zamienia progu w bramkę.

### 6.1. Logi

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.logi.poziom_domyslny` | Domyślny poziom filtra strumienia | wszystkie |
| `diagnostics.logi.zrodlo_domyslne` | Domyślne źródło (proces, moduł, karta sesji) | bieżąca sesja |
| `diagnostics.logi.bufor` | Rozmiar bufora wpisów na żywo | 5000 wpisów |
| `diagnostics.logi.retencja` | Polityka retencji archiwum per zasięg | 30 dni |
| `diagnostics.logi.parsowanie` | Parsowanie strukturalne JSON/logfmt | auto-wykrycie |
| `diagnostics.logi.dedup_okno` | Okno deduplikacji powtórzeń | 2 minuty |

### 6.2. Błędy

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.bledy.reguly_odcisku` | Reguły normalizacji stosu do odcisku | wbudowane per język |
| `diagnostics.bledy.status_domyslny` | Domyślny filtr statusu | nowy |
| `diagnostics.bledy.alert_nowy_typ` | Alert przy pierwszym wystąpieniu odcisku | włączony |
| `diagnostics.bledy.powiazanie_wdrozenia` | Kojarzenie błędu ze zbieżnym commitem (Developer) | włączone |

### 6.3. Prowenancja wywołań AI

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.prowenancja.wlaczona` | Zbieranie śladów wywołań kanałów modeli | włączone |
| `diagnostics.prowenancja.zapis_tresci` | Zapisywanie pełnej treści promptu i odpowiedzi | włączony |
| `diagnostics.prowenancja.redakcja` | Redakcja danych wrażliwych w treści śladu | wyłączona |
| `diagnostics.prowenancja.prog_wolnego` | Próg opóźnienia oznaczający wywołanie jako wolne | 10 s |
| `diagnostics.prowenancja.retencja` | Retencja śladów | 30 dni |
| `diagnostics.prowenancja.eksport` | Cel eksportu OTLP (zewnętrzny zbieracz) | brak |

### 6.4. Zużycie i budżety

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.koszt.cennik` | Cennik per kanał modelu (do wyliczeń) | z definicji kanału |
| `diagnostics.koszt.okres_rozliczeniowy` | Okres agregacji zużycia | miesięczny |
| `diagnostics.koszt.budzet` | Miękkie budżety per zasięg (sesja, projekt, środowisko, konto) | brak |
| `diagnostics.koszt.prog_ostrzezenia` | Procent budżetu wyzwalający ostrzeżenie | 80% |
| `diagnostics.koszt.anomalia` | Detekcja skoków kosztu | włączona |

### 6.5. Metryki i wydajność

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.metryki.retencja_tsdb` | Retencja szeregów czasowych | 15 dni |
| `diagnostics.metryki.percentyle` | Zestaw prezentowanych percentyli | p50/p95/p99 |
| `diagnostics.metryki.profilowanie` | Profilowanie pprof i profilowanie ciągłe | na żądanie |
| `diagnostics.metryki.zapytania_wlasne` | Zdefiniowane zapytania i formuły | brak |

### 6.6. Kontrola stanu, alerty i widoczność zakładek Observability Tools

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.narzedzia.provenance` | Widoczność zakładki Provenance Explorer | widoczna |
| `diagnostics.narzedzia.usage_cost` | Widoczność zakładki Usage & Cost | widoczna |
| `diagnostics.narzedzia.metrics` | Widoczność zakładki Metrics & Performance | ukryta |
| `diagnostics.narzedzia.health` | Widoczność zakładki Health & Uptime | ukryta |
| `diagnostics.narzedzia.alerts` | Widoczność zakładki Alerts | ukryta |
| `diagnostics.health.sondy` | Zdefiniowane sondy zdrowia i częstotliwość | brak |
| `diagnostics.alerty.kanaly` | Kanały powiadomień (aplikacja, e-mail, komunikator, webhook) | aplikacja |
| `diagnostics.alerty.always_on` | Wypychanie alertów krytycznych do Always On Display | włączone |
| `diagnostics.center.prog_krytyczny` | Próg powiadomienia o zdarzeniu krytycznym | z wartości domyślnej |

### 6.7. Pętla wykonawcza (Execution Loop Window)

| Klucz | Co Operator personalizuje | Domyślna |
|---|---|---|
| `diagnostics.petla.widocznosc` | Otwarcie kolumny pętli wykonawczej przy starcie sesji modułu | włączone |
| `diagnostics.petla.limit_ponowien` | Maksymalna liczba ponowień zadania diagnostycznego w przebiegu | 3 |
| `diagnostics.petla.prog_pewnosci` | Minimalny poziom pewności rekomendacji przyjmowany bez ponowienia | średni |
| `diagnostics.petla.autonomia` | Zakres samodzielności Wykonawcy w pętli: raportowanie kroków albo realizacja ciągła do wyniku | raportowanie kroków |
| `diagnostics.petla.zrodla_zadan` | Okna materiałowe uwzględniane w dekompozycji zlecenia | Logs Viewer, Errors Panel, Provenance Explorer |
| `diagnostics.petla.rejestr_przebiegu` | Zapis przebiegu pętli do rejestru sesji i eksport przebiegu | włączony |

### 6.8. Izolacja i zależności (powiązania jawne)

Zgodnie z rozdz. 7.3 oraz oknem konfiguracji punktów izolacji sterowaniu podlegają: zakres odczytu plików dziennika spoza katalogu roboczego sesji, dostęp sieciowy do zewnętrznych systemów monitorowania, katalog roboczy sesji dla raportów i migawek, model procesu analizy oraz konto i token per sesja (dane dostępowe do zewnętrznych źródeł logów i monitoringu). Powiązania `Diagnostics→Developer` (Recommendations Panel → Code Editor), `Diagnostics→Terminal` (odtworzenie objawu) oraz `MultitaskingAI→Diagnostics` (zdarzenia niepowodzeń procesów) konfigurowane są jako `powiazanie_komponentu`. Redakcja danych wrażliwych w prowenancji (`diagnostics.prowenancja.redakcja`) jest ustawieniem konfiguracyjnym Operatora, nie wymogiem.

---

## 7. Stany, dane i powiązania

### 7.1. Model stanów rekomendacji

```
                    ┌────────────────┐
                    │  Zebrany       │  (log / błąd / ślad / metryka)
                    │  materiał      │
                    └───────┬────────┘
                            │ uruchomienie analizy (Diagnostics Center
                            │ albo zlecenie pętli w Execution Loop Window)
                            ▼
                    ┌────────────────┐
                    │  Analiza w toku│  (agregacja, wnioskowanie Wykonawcy)
                    └───────┬────────┘
                            │ wniosek gotowy
                            ▼
                    ┌────────────────┐
        ┌──────────►│  Rekomendacja  │
        │           │  zaproponowana │
        │           └───────┬────────┘
        │                   │
        │      ┌────────────┼────────────┐
        │      ▼            ▼            ▼
        │ ┌─────────┐ ┌───────────┐ ┌───────────────┐
        │ │Odrzucona│ │   Inne    │ │  Zastosowana  │
        │ │         │ │rozwiązanie│ │ (przekazana   │
        │ └─────────┘ └─────┬─────┘ │  do Developer)│
        │                   │       └───────┬───────┘
        └───────────────────┘               │ obserwacja skuteczności
                                            ▼
                                  ┌───────────────────┐
                                  │ Skuteczność       │
                                  │ potwierdzona /    │
                                  │ błąd powrócił     │
                                  └───────────────────┘
```

Stan procesu sesji po stronie serwera jest trwały niezależnie od stanu połączenia klienta — analiza uruchomiona w Diagnostics Center oraz przebieg pętli prowadzony w Execution Loop Window kończą się niezależnie od rozłączenia klienta; po ponownym połączeniu Recommendations Panel prezentuje gotowy wynik, a Execution Loop Window — pełny rejestr przebiegu.

### 7.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Diagnostics |
|---|---|
| `karta_sesji`, `sesja` | Nośnik kontekstu bieżącego dochodzenia diagnostycznego |
| `proces_sesji` | Proces serwera obsługujący sesję modułu; źródło danych o zużyciu zasobów prezentowanych w Diagnostics Center i Metrics & Performance |
| `wiadomosc`, `zalacznik_wiadomosci` | Historia Chat Window, w tym cytowane logi, stosy wywołań i załączone zrzuty awaryjne; treść promptu i odpowiedzi w Provenance Explorer |
| `zadanie` | Uruchomienie pełnej analizy, zadanie pętli wykonawczej i sonda zdrowia jako jednostki pracy, z powiązaniem z silnikiem kolejek |
| `kanal_modelu` | Model bazowy przypisany karcie sesji, generujący analizy i rekomendacje; wymiar prowenancji, zużycia i kosztu |
| `ustawienie` | Parametry modułu podlegające zasadzie „brak ustawienia = wartość domyślna" — zakres czasu Diagnostics Center, progi alertów, budżety, retencja, parametry pętli wykonawczej |
| `profil_izolacji`, `regula_izolacji_technicznej` | Konfiguracja izolacji technicznej procesu diagnostycznego (rozdz. 7.3), w szczególności zakresu odczytu logów spoza katalogu roboczego sesji oraz wymiaru konta i tokenu w rozliczeniu zużycia |
| `powiazanie_komponentu` | Jawne powiązania Diagnostics ◄──► Developer, Diagnostics ◄──► Terminal, MultitaskingAI ──► Diagnostics; graf zależności komponentów w Health & Uptime |

Strumień logów, rejestr błędów, ślady wywołań i metryki prezentowane w oknach modułu pochodzą z bieżących zdarzeń systemowych przekazywanych na żywo kanałem WebSocket; materiał historyczny wykraczający poza bufor bieżącej sesji odczytywany jest na żądanie z archiwum dziennika zdarzeń, magazynu śladów i magazynu szeregów czasowych, zgodnie z zasadą synchronizacji na żywo właściwą oknom monitorującym procesy. Wspólnym kontraktem sygnałów jest OpenTelemetry, a polem korelacyjnym spinającym logi, błędy, ślady i metryki — `trace_id`.

### 7.3. Izolacja i konfigurowalność — punkty właściwe modułowi Diagnostics

| Zakres izolacji | Stan wyjściowy | Zastosowanie w module Diagnostics | Typowy powód włączenia |
|---|---|---|---|
| Zakres odczytu i zapisu plików | Wyłączony (pełny dostęp) | Dostęp Logs Viewer do plików dziennika zdarzeń poza katalogiem roboczym sesji | Ograniczenie sesji diagnostycznej do logów jednego, wskazanego komponentu |
| Dostęp sieciowy | Wyłączony (współdzielony) | Pobieranie zdarzeń z zewnętrznych systemów monitorowania i import telemetrii zewnętrznej | Odseparowanie sesji analizującej materiał pochodzący z zewnętrznego, niezaufanego źródła |
| Katalog roboczy sesji | Wyłączony (współdzielony) | Miejsce przechowywania eksportowanych raportów diagnostycznych i migawek stanu | Odrębny katalog dla sesji prowadzącej dochodzenie objęte poufnością |
| Model procesu | Wyłączony (współdzielony) | Instancja procesu wykonawczego modelu obsługującego analizę oraz pętlę wykonawczą | Odrębny proces modelu dla sesji analizującej duży wolumen logów |
| Konto i token per sesja | Wyłączony (współdzielone dane dostępowe) | Dane dostępowe do zewnętrznych systemów logowania i monitorowania; wymiar rozliczenia zużycia w Usage & Cost | Odrębny token dostępu przy analizie logów środowiska produkcyjnego klienta |

Zgodnie z zasadą nadrzędną platformy żaden z powyższych punktów nie jest wymuszony — moduł jest w pełni funkcjonalny bez jakiejkolwiek konfiguracji izolacji. Domyślnym przepływem zastosowania rekomendacji jest przekazanie zmiany do Code Editor modułu Developer do przeglądu, po czym zatwierdzenie następuje w Git Panel. Przepływ ten jest w pełni konfigurowalny — Operator włącza bezpośrednie zastosowanie poprawki z pominięciem etapu przeglądu — zgodnie z zasadą pełnej konfigurowalności; nie jest to bramka wymuszona przez moduł. Próg budżetu i próg kondycji sygnalizują oraz rekomendują zamiast blokować, a niedostępność integracji ukrywa zakładkę zamiast prezentować zablokowany element interfejsu.

### 7.4. Powiązania z innymi modułami

```
                     ┌────────────────────────────────┐
    Developer ◄──────┤          DIAGNOSTICS           ├──────►  Terminal
   (wdrażanie        │  agregacja logów, błędów,      │        (odtwarzanie i
    poprawek         │  metryk i śladów; prowenancja  │         weryfikacja objawów
    wynikających     │  wywołań, kontrola kosztu,     │         błędu w rzeczywistym
    z analizy)       │  rekomendacje naprawcze)       │         środowisku wykonawczym)
                     └────────────────────────────────┘
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Developer | Diagnostics współpracuje z modułem Developer przy wdrażaniu poprawek wynikających z analizy oraz przy kojarzeniu błędów ze zbieżnymi zmianami | Konfiguracyjne | Recommendations Panel → Code Editor; Errors Panel → Git Panel |
| Terminal | Diagnostics współpracuje z Terminalem przy odtwarzaniu i weryfikacji objawów błędu oraz przy wykonywaniu scenariuszy monitoringu syntetycznego | Konfiguracyjne | Chat Window / Diagnostics Center → Terminal Tabs; Health & Uptime → Terminal Tabs |
| MultitaskingAI | Diagnostics odbiera zdarzenia niepowodzeń procesów pracy ciągłej jako materiał źródłowy | Konfiguracyjne | Monitor procesu panelu orkiestracji → Errors Panel |

Współobecność Diagnostics, Developer i Terminala w jednym środowisku CodeStudio czyni powiązanie typowe w praktyce, lecz — zgodnie z zasadą jawności i konfigurowalności zależności — nie jest ono wbudowaną na stałe zależnością techniczną; ustanawia je Operator z okna konfiguracji.

---

## 8. Scenariusze użycia

**Scenariusz 1 — od zgłoszenia błędu do zagregowanego obrazu przyczyny.**

Zespół otrzymuje zgłoszenie o błędach usługi płatności. Operator otwiera Errors Panel, odnajduje zgrupowany wpis z dwunastoma wystąpieniami, wybiera „Analizuj" z menu wiersza — Diagnostics Center agreguje ten błąd z powiązanymi logami z ostatniej godziny i przedstawia zwięzły wniosek: usługa zewnętrzna przekracza limit czasu odpowiedzi.

**Scenariusz 2 — zlecenie złożone prowadzone w pętli wykonawczej.**

Operator formułuje w Chat Window zlecenie „ustal przyczynę błędów usługi płatności". Koordynator dekomponuje je w Execution Loop Window na zadania: pobranie logów zakresu, zbudowanie odcisków błędów, korelację śladów wywołań modelu i wyliczenie metryk okresu. Wykonawca realizuje zadania i raportuje wynik każdego kroku; kontrola jakości odrzuca pierwszy wynik korelacji jako niekompletny i Koordynator ponawia zadanie z szerszym zakresem czasu. Po zakończeniu przebiegu Recommendations Panel prezentuje rekomendacje, a Chat Window — podsumowanie przebiegu.

**Scenariusz 3 — zastosowanie poprawki z pełnym przeglądem.**

Na podstawie analizy Recommendations Panel przedstawia rekomendację dodania limitu czasu i ponowienia próby połączenia. Operator wybiera „Zastosuj poprawkę" — Code Editor modułu Developer otwiera się z gotową zmianą przygotowaną do przeglądu. Po weryfikacji poprawności Operator zatwierdza zmianę w Git Panel; rekomendacja zostaje oznaczona jako zastosowana i wchodzi w okres obserwacji skuteczności.

**Scenariusz 4 — weryfikacja hipotezy w rzeczywistym środowisku.**

Przed wdrożeniem poprawki Operator potwierdza hipotezę o przeciążonej puli połączeń do bazy danych. Z poziomu Chat Window zleca sprawdzenie liczby aktywnych połączeń — polecenie wykonuje się w module Terminal, a wynik wraca do Diagnostics, potwierdzając lub obalając hipotezę przed podjęciem dalszych kroków.

**Scenariusz 5 — prowenancja i koszt wywołania modelu.**

Wpis logu ze znacznikiem `trace_id` prowadzi Operatora do Provenance Explorer, gdzie drzewo śladu ujawnia, że jedno z wywołań podagenta przekracza próg opóźnienia i zużywa nieproporcjonalnie dużo tokenów. Usage & Cost pokazuje udział tego kanału modelu w koszcie okresu i sygnalizuje zbliżenie do miękkiego progu budżetu sesji. Operator zmienia kanał modelu dla tego typu zadań, a porównanie kosztu i opóźnienia potwierdza efekt zmiany.

**Scenariusz 6 — przegląd stanu systemu przed wydaniem.**

Przed planowanym wydaniem lider techniczny otwiera Diagnostics Center, ustawia zakres czasu na ostatni tydzień i przegląda wskaźniki kondycji wszystkich komponentów oraz oś czasu zdarzeń. Health & Uptime potwierdza zużycie budżetu błędów poniżej limitu, brak nowych błędów krytycznych i stabilny trend czasu odpowiedzi; wynik eksportowany jest jako raport dla zespołu.

**Scenariusz 7 — nadzór nad procesem autonomicznym środowiska MultitaskingAI.**

Agent pełniący rolę Executora w nienadzorowanej pętli pracy ciągłej napotyka powtarzające się niepowodzenie kroku procesu. Zdarzenie trafia do Errors Panel modułu Diagnostics z pełnym kontekstem przebiegu, a reguła alertu wypycha sygnał do Always On Display. Operator zleca analizę przyczyny w Execution Loop Window, przegląda wniosek w Diagnostics Center i wprowadza poprawkę w module Developer, zanim proces wznowi kolejne uruchomienia według harmonogramu.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + F` | Wyszukiwanie / Grep w bieżącym oknie | Logs Viewer, Errors Panel |
| `Ctrl/Cmd + Shift + A` | Uruchomienie pełnej analizy | Diagnostics Center |
| `Ctrl/Cmd + Shift + L` | Otwarcie kolumny pętli wykonawczej | Execution Loop Window |
| `Ctrl/Cmd + K` | Wyszukiwarka funkcji — dostęp do funkcji eksperckich warstwy 4 | Wszystkie okna modułu |
| `Ctrl/Cmd + Enter` | Wysłanie polecenia do Wykonawcy | Chat Window |
| `/` | Wywołanie menu operacji z poziomu okna komunikacji | Chat Window |
| `Esc` | Zwinięcie rozwiniętego panelu szczegółów | Errors Panel, Recommendations Panel, Provenance Explorer |
| Ikona `oko` | Podgląd, Monitor procesu | Diagnostics Center |
| Ikona `blad` | Alert / plakietka błędu | Errors Panel, Logs Viewer |
| Ikona `ostrzezenie` | Alert / plakietka ostrzeżenia | Diagnostics Center, Logs Viewer, Alerts |
| Ikona `info` | Alert / plakietka informacyjna | Diagnostics Center |
| Ikona `zegar` | Zakres czasu, oś czasu zdarzeń | Diagnostics Center, Logs Viewer, Metrics & Performance |
| Ikona `odswiez` | Ponowienie analizy, polecenia albo zadania pętli | Diagnostics Center, Recommendations Panel, Execution Loop Window |
| Ikona `petla` | Przebieg pętli wykonawczej, ponowienie zadania | Execution Loop Window |
| Ikona `ptaszek` / `ptaszek-kolo` | Potwierdzenie skuteczności poprawki, przyjęcie wyniku zadania | Recommendations Panel, Execution Loop Window |
| Ikona `tarcza` | Kontekst bezpieczeństwa (np. wcielenie Security Auditor) | Errors Panel, Recommendations Panel |
| Ikona `kod` | Blok treści technicznej (log, stos wywołań, poprawka, prompt) | Logs Viewer, Errors Panel, Recommendations Panel, Provenance Explorer |
| Ikona `moneta` | Koszt wywołania, zużycie budżetu | Usage & Cost, Execution Loop Window |
| Ikona `wykres` | Szereg czasowy, percentyle, profil | Metrics & Performance |
| Ikona `puls` | Sonda zdrowia, dostępność usługi | Health & Uptime, Diagnostics Center |
| Ikona `dzwonek` | Reguła alertu, powiadomienie | Alerts |
| Ikona `pobierz` | Eksport raportu, zakresu logów, metryk lub śladów | Diagnostics Center, Logs Viewer, Recommendations Panel, Observability Tools |

*Koniec dokumentu. Moduł Diagnostics — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*

*© 2026 Danaco Holding Group Sp. z o.o. — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
