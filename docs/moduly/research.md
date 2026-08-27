*Dokument specyfikuje interfejs modułu Research Danaco Console: okna, makiety, elementy, warstwy widoczności oraz scenariusze pracy.*

# Moduł Research — dokumentacja projektowa

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
| **Tytuł** | Moduł Research — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować): komplet okien operacyjnych, katalog funkcji i narzędzi modułu oraz punkty sterowania z okna konfiguracji |
| **Środowiska dostępności** | TalkIn, WorkSpace |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 4, 9.5, 12, Załącznik A) · Specyfikacja modułów 1.1 (rozdz. 4.5) · Specyfikacja okien operacyjnych 1.1 (rozdz. 5, 6.5, 10) · System wizualny (rozdz. 6, 8, 10) · Model danych 1.1 (rozdz. 8, 16, 17, Załącznik F.2) · Izolacja i zależności (rozdz. 1–6) · Architektura rdzenia (Go, SQLite, WebSocket/JSON, klient Tauri) · Model konfiguracji |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); zero blokad w interfejsie; klucze i integracje jawne; izolacja i uprawnienia są ustawieniami konfiguracyjnymi Operatora, nie wymogiem; domyślne zachowanie modułu = wykonanie |

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

Research jest modułem realizacji badań, analiz i opracowań — przestrzenią, w której zbieranie źródeł, odnotowywanie ustaleń cząstkowych i kompletowanie raportu końcowego następują w jednym, ciągłym procesie, zamiast rozpraszać się między notatnik, przeglądarkę i osobne narzędzie do pisania. Moduł nadaje strukturę pracy badawczej, która w pojedynczej rozmowie z modelem pozostałaby płaska i trudna do prześledzenia wstecz.

Research jest kompletnym stanowiskiem pracy badawczej — jednym miejscem, w którym analityk, konsultant, badacz akademicki, dziennikarz śledczy i zespół produktowy przechodzą całą drogę: od pytania badawczego, przez wyszukiwanie i pozyskanie źródeł, ich lekturę, ekstrakcję i weryfikację, po syntezę, cytowanie i eksport gotowego opracowania. Użytkownik nie sięga po żaden osobny program z obszaru badań: moduł obejmuje zakres menedżera bibliografii, wyszukiwarki literatury, czytnika PDF z adnotacjami, narzędzia do web-clippingu, agregatora RSS, programu do przeglądów systematycznych, transkryptora nagrań, narzędzia do analizy jakościowej i generatora cytowań.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Research |
|---|---|
| Analitycy rynku i strategii | Analizy konkurencyjne, benchmarki, opracowania sektorowe z wielu źródeł |
| Konsultanci i doradcy | Raporty dla klienta łączące dane pierwotne i wtórne w spójny dokument |
| Zespoły produktowe | Badania użytkownika, analiza rynku docelowego, walidacja hipotez |
| Dziennikarze i redaktorzy śledczy | Systematyczne gromadzenie i krzyżowa weryfikacja źródeł |
| Zespoły akademickie i eksperckie | Przeglądy literatury, opracowania z odniesieniami do źródeł |

### 1.3. Po co — wartość modułu

| Problem pracy badawczej rozproszonej | Rozwiązanie w Research |
|---|---|
| Źródła gubią się w zakładkach przeglądarki i notatkach | Sources Manager gromadzi je w jednym miejscu z metadanymi i oceną wiarygodności |
| Wnioski cząstkowe giną w długiej historii czatu | Findings Panel odnotowuje ustalenia powiązane bezpośrednio ze źródłem |
| Budowa raportu od zera z rozproszonych fragmentów | Report Builder komponuje dokument końcowy z materiału już zgromadzonego |
| Ręczne dostosowanie formatu na potrzeby odbiorcy | Export Panel eksportuje raport do formatu docelowego w każdej chwili — skompletowanie raportu to sugestia, nie warunek |
| Rozproszenie wyszukiwania między wyszukiwarki, bazy naukowe i czytniki kanałów | Discovery Panel scala wyszukiwanie webowe, naukowe, semantyczne i pełnotekstowe w jednym oknie |
| Lektura i adnotacje prowadzone w osobnym czytniku | Reading View udostępnia lekturę, podświetlenia, wypisy i ekstrakcję w obrębie modułu |

### 1.4. Zakres tematyczny modułu

| Obszar | Zakres w Research |
|---|---|
| Wyszukiwanie | Wyszukiwanie webowe i naukowe, przeszukiwanie własnych źródeł, wyszukiwanie semantyczne, monitorowanie tematów (alerty, RSS) |
| Pozyskiwanie źródeł | Web-clipping, import PDF/DOCX/EPUB, import z DOI/ISBN/arXiv/PubMed, zrzuty stron, transkrypcja audio/wideo |
| Zarządzanie źródłami | Katalog, metadane, ocena wiarygodności, deduplikacja, tagi, kolekcje, biblioteka referencji |
| Lektura i ekstrakcja | Czytnik z adnotacjami, podświetlenia, wypisy, OCR skanów, ekstrakcja tabel i danych liczbowych |
| Analiza | Kodowanie jakościowe, wykrywanie sprzeczności, mapa powiązań, analiza porównawcza źródeł, ekstrakcja twierdzeń |
| Synteza | Odnotowywanie ustaleń, grupowanie w wątki, budowa raportu, streszczenia, przeglądy literatury |
| Cytowania | Generowanie cytatów i bibliografii w stylach CSL, przypisy, zarządzanie referencjami, sprawdzanie kompletności cytowań |
| Eksport | PDF, DOCX, Markdown, HTML, PPTX, XLSX, LaTeX/BibTeX, przekazanie do Library/Studio/Roundtable |

### 1.5. Granice zewnętrzne modułu

| Poza zakresem | Gdzie realizowane | Uzasadnienie granicy |
|---|---|---|
| Trwałe, wielomodułowe przechowywanie plików i wersjonowanie | **Library** (Versioning Panel) | Research operuje sesyjnie; trwały magazyn to rola Library — Research zasila go i z niego czerpie |
| Zaawansowana redakcja językowa i WYSIWYG końcowego dokumentu | **Studio** (Studio Editor) | Research kompletuje treść merytoryczną; szlif redakcyjny to domena Studio |
| Wielomodelowa konfrontacja wniosków | **Roundtable** | Research przekazuje syntezę do oceny wielomodelowej |
| Wspólne przeglądanie stron przez użytkownika i AI na żywo | **Browser** | Research odbiera z Browser gotowe źródła, nie duplikuje przeglądarki współdzielonej |
| Głęboka analiza statystyczna dużych zbiorów danych, wizualizacje BI | **moduł analizy danych / Studio** | Research wykonuje proste zestawienia i wykresy poglądowe z danych z ustaleń, nie zastępuje narzędzia analitycznego |
| Cykliczne uruchamianie badań według harmonogramu | **Automations / Agents** | Research udostępnia operacje; ich cykliczność orkiestruje Automations |

Granica jest konfigurowalna, nie sztywna — powiązania Browser ◄──► Research, Research ───► Library, Research ◄──► Studio, Research ───► Roundtable są jawnymi kanałami ustanawianymi w oknie konfiguracji (rozdz. 5.3–5.4).

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn albo WorkSpace
        ▼
ŚRODOWISKO (TalkIn / WorkSpace) ── boczna nawigacja modułów
        │
        ▼
MODUŁ: RESEARCH ────────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Research Workspace ·
  Discovery Panel · Sources Manager · Reading View ·
  Findings Panel · Report Builder · Export Panel
        │
        ▼
KARTA SESJI — własny układ, historia i kontekst badania
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 5.3)
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | TalkIn, WorkSpace |
| Środowisko CodeStudio | Moduł niedostępny |
| Komponent własny | Research nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 9 (łącznie z Chat Window i Execution Loop Window) |
| Typ pracy | Sesyjna, projekt badawczy prowadzony od zbierania źródeł po eksport raportu |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja Użytkownik ↔ Wykonawca (wspólne wszystkim modułom) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji przy każdym wejściu do modułu | Centralny punkt pracy: polecenia w języku naturalnym, strumień odpowiedzi i wyników, zatwierdzanie i przerywanie działań |
| 2 | Execution Loop Window | Komunikacja Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana | 2 | Znacznik stanu pętli w pasku kontekstu (`⟳ Pętla ▼`) — jedno kliknięcie otwiera kolumnę; w stanie spoczynku widoczny jest wyłącznie znacznik jako wskaźnik stanu wykonania warstwy 1 | Pętla wykonawcza badania: dekompozycja zlecenia na zadania, kolejka i stan zadań, kontrola jakości, sterowanie przebiegiem |
| 3 | Research Workspace | Zarządzanie źródłami i ustaleniami / punkt wejścia | Prawa kolumna, dominująca | 1 | Widoczne bez interakcji jako aktywne okno wiodące modułu | Struktura badania, nawigacja do pozostałych okien |
| 4 | Discovery Panel | Wyszukiwanie i odkrywanie źródeł | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Skrót nawigacyjny `→ Discovery` w stopce Research Workspace, polecenie języka naturalnego w Chat Window; po zamknięciu panel znika z przestrzeni roboczej | Wyszukiwanie webowe, naukowe, semantyczne i pełnotekstowe, snowballing, monitory tematów |
| 5 | Sources Manager | Zarządzanie źródłami i ustaleniami | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Skrót nawigacyjny `→ Sources`, kliknięcie licznika źródeł w panelu postępu, odnośnik cytowania w odpowiedzi Wykonawcy | Gromadzenie i katalogowanie źródeł, biblioteka referencji |
| 6 | Reading View | Okno lektury i ekstrakcji | Kolumna boczna, otwierana jako rozszerzenie boczne, szeroka | 2 | Przycisk „Czytaj” przy pozycji źródła, kliknięcie kotwicy ustalenia | Lektura źródła, podświetlenia i adnotacje, wypisy, ekstrakcja tabel i twierdzeń |
| 7 | Findings Panel | Zarządzanie źródłami i ustaleniami | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Skrót nawigacyjny `→ Findings`, przycisk „Dodaj jako ustalenie” w dymku odpowiedzi, przycisk „→ ustalenie” w Reading View | Odnotowywanie ustaleń cząstkowych, kodowanie jakościowe |
| 8 | Report Builder | Okno edycyjne / kreator | Prawa kolumna, dominująca po aktywacji | 2 → 1 | Skrót nawigacyjny `→ Report Builder` albo przycisk „→ Report Builder” przy ustaleniu; po aktywacji przejmuje rolę aktywnego okna wiodącego i należy do warstwy 1 | Kompletowanie raportu końcowego |
| 9 | Export Panel | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne, wąska | 2 | Przycisk „→ Export Panel” w stopce Report Builder, skrót nawigacyjny `→ Export`, polecenie eksportu w Chat Window | Eksport raportu do formatu docelowego |

```
 Makieta zbiorcza — stan spoczynku            Dostępność: TalkIn, WorkSpace
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window                    │ Research Workspace
  nawigacja │ Użytkownik ↔ Wykonawca         │ (aktywne okno wiodące)
  modułów   │                                │
  (poza     │ [Research][Fable 5][⟳ Pętla ▼] │ Temat badania · etapy ·
  zakresem  │                                │ materiały badania
  tego      │ strumień rozmowy               │
  dokumentu)│                                │ postęp badania
            │                                │
            │ Pole poleceń …        [ ⋮ ]    │ [ Okna modułu ▼ ]
 ═══════════════════════════════════════════════════════════════════════════
  Warstwa 1: Chat Window · Research Workspace · pasek kontekstu · postęp
  Zwinięte wyzwalacze: [⟳ Pętla ▼] [ ⋮ ] [ Okna modułu ▼ ] · znaczniki kontekstu
```

### 2.1. Warstwy widoczności w module

Moduł Research stosuje regułę stopniowego ujawniania funkcjonalności: **jeżeli funkcja
nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Pełny arsenał
badawczy — dostawcy wyszukiwania, kodowanie jakościowe, protokół PRISMA, style CSL,
ekstrakcja tabel, OCR, monitory tematów — istnieje w architekturze modułu i pozostaje
poza polem widzenia do chwili wystąpienia potrzeby użycia. Liczba okien, narzędzi
i ustawień modułu nie wpływa na postrzeganą prostotę jego interfejsu.

| Warstwa | Zawartość w module Research | Sposób dostępu |
|---|---|---|
| 1 — zawsze widoczna | Chat Window (kanał Użytkownik ↔ Wykonawca), aktywne okno wiodące obszaru roboczego (Research Workspace, po aktywacji Report Builder), pasek kontekstu badania (moduł, model, wykonawca, stan pętli), nagłówek zakresu badania, pasek etapów, panel postępu badania, wskaźniki stanu wykonania zadań. Zajmuje ponad 80% powierzchni interfejsu modułu | Widoczne bez interakcji |
| 2 — widoczna na żądanie | Execution Loop Window, Discovery Panel, Sources Manager, Reading View, Findings Panel, Export Panel, selektor kanału modelu, selektor stylu cytatu, selektor szablonu raportu, selektor formatu eksportu, przełączniki widoku (lista, oś czasu, Evidence Map, Kanban, tabela dowodów; chronologia, wątki, kodowanie), przełącznik trybu wyszukiwania | Znacznik kontekstowy, ikona, przycisk lub przełącznik; element zbiorczy `Okna modułu ▼`, `Widok ▼`, `Styl cytatu ▼`; po użyciu element zwija się samoczynnie, a zamknięty panel znika całkowicie z przestrzeni roboczej |
| 3 — rozwinięcia kontekstowe | Zestawy akcji pozycji źródła (Czytaj, Cytuj, Powiąż z ustaleniem, Załączniki, Scal duplikat, Usuń z listy), zestawy akcji pozycji wyniku wyszukiwania (Dodaj do źródeł, Cytowania, Odrzuć, Podgląd), zestaw filtrów Discovery Panel, pasek operacji na źródle w Reading View (streszczenie, tabela, twierdzenia, zapytanie, OCR), operacje kontekstowe na zaznaczeniu w Report Builder, menedżer przypisów i bibliografii, panel wersji i różnic, ustawienia personalizacji eksportu, książka kodów | Menu kebab (`⋮`, `⋯`), menu hamburger (`☰`), menu kontekstowe zaznaczenia, panel popover, lista rozwijana, element zbiorczy `Operacje ▼` |
| 4 — funkcje eksperckie | Tryb przeglądu systematycznego z diagramem PRISMA, macierz kod × źródło i eksport macierzy kodowania, reguły kolekcji inteligentnych, progi deduplikacji i punktacji wiarygodności, import wsadowy adresów i ponowne przetworzenie kolekcji, edycja stylu CSL własnego, konfiguracja dostawców wyszukiwania i kluczy, monitorowany folder importu, diagnostyka pętli wykonawczej i śladu prowenancji | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, okno konfiguracji (rozdz. 8). Użytkownik podstawowy nie widzi tych elementów |

**Zasada jednego kliknięcia.** Każda ukryta funkcja modułu jest osiągalna jednym

kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego
w Chat Window. Ukrycie zmniejsza chaos wizualny i nie wydłuża drogi dostępu —
zagnieżdżanie funkcji głęboko w hierarchii menu jest w module wykluczone.

Makiety tekstowe okien w rozdziale 3 przedstawiają stan spoczynku interfejsu: widoczne
są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (`▼`, `⋮`, `☰`)
i znaczniki kontekstowe. Elementy warstw 2–4 opisane są w tabelach elementów każdego
okna wraz z podaniem warstwy i sposobu wywołania.

---

## 3. Specyfikacja okien operacyjnych

### 3.1. Chat Window (okno wspólne)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI / agent / system wykonawczy). Stanowi centralny punkt pracy użytkownika w module Research i podstawowy mechanizm sterowania wszystkimi procesami badawczymi: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu badania. Okno zajmuje to samo miejsce układu we wszystkich modułach i środowiskach — lewą kolumnę obszaru roboczego.

**Zawartość i pełny arsenał funkcji.**

- Historia rozmowy z Wykonawcą w kontekście bieżącego badania.
- Polecenia analityczne: „streść ten zbiór źródeł”, „wskaż sprzeczności między źródłem A i B”, „zaproponuj strukturę raportu”, „znajdź sprzeczności”.
- Strumień odpowiedzi i wyników na żywo kanałem WebSocket.
- Odwołania kontekstowe `@źródło`, `@ustalenie`, `@kolekcja` przez wzmiankę w treści polecenia.
- Zatwierdzanie i przerywanie działań realizowanych przez Wykonawcę.
- Generowanie pytań badawczych pogłębiających temat na podstawie dotychczasowych ustaleń.
- Wstawienie odpowiedzi bezpośrednio jako nowego wpisu w Findings Panel jednym kliknięciem.
- Cytowanie źródła w odpowiedzi z automatycznym odnośnikiem do pozycji w Sources Manager.
- Podsumowanie stanu badania na żądanie („gdzie jesteśmy”).
- Historia poleceń, regeneracja odpowiedzi, rozgałęzianie wątku.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window — Użytkownik ↔ Wykonawca   │  Research Workspace
  [Research][Fable 5][⟳ Pętla ▼]         │  (aktywne okno wiodące)
  ─────────────────────────────────────  │
  ┌─ Użytkownik ───────────────────────┐ │  Temat badania
  │ „wskaż sprzeczności między         │ │  etapy · postęp
  │  źródłem 3 i 7”                    │ │
  └────────────────────────────────────┘ │
  ┌─ Wykonawca (strumień na żywo) ─────┐ │
  │ analiza z odnośnikami [3] [7] …    │ │
  │                            [ ⋮ ]   │ │
  └────────────────────────────────────┘ │
  ─────────────────────────────────────  │
  Pole poleceń …            [ ⋮ ][ ▶ ]   │
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Strumień rozmowy | Historia wymiany Użytkownik ↔ Wykonawca | Prezentacja poleceń, odpowiedzi i wyników badania | Duży obszar, dominujący w oknie | 1 | Widoczny bez interakcji | domyślny · strumieniowanie odpowiedzi | Przewijanie historii, zaznaczenie fragmentu odpowiedzi | Ciało Chat Window |
| Pasek kontekstu badania | Rząd lekkich znaczników `[Research] [Fable 5] [⟳ Pętla]` | Wskazanie modułu, kanału modelu, wykonawcy i stanu pętli wykonawczej | Mały rząd plakietek | 1 | Widoczny bez interakcji; kliknięcie znacznika otwiera odpowiedni selektor | domyślny · pętla w toku · pętla wstrzymana | Kliknięcie znacznika modelu otwiera selektor kanału, kliknięcie znacznika pętli otwiera Execution Loop Window | Nagłówek Chat Window |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzanie polecenia analitycznego | Duże pole w stopce okna | 1 | Widoczne bez interakcji | domyślny · fokus · z odwołaniem do źródła | Enter wysyła, `@` wywołuje odwołanie do źródła/ustalenia/kolekcji | Stopka Chat Window |
| Selektor kanału modelu | Rozwijana lista wywoływana ze znacznika kontekstu | Wybór i zmiana modelu bazowego bieżącej karty | Mały przycisk z etykietą | 2 | Kliknięcie znacznika `[Fable 5]` w pasku kontekstu; po wyborze lista zwija się samoczynnie | domyślny · rozwinięty · ładowanie | Wybór zmienia model dla kolejnych analiz | Pasek kontekstu Chat Window |
| Odwołanie `@źródło` | Podpowiedź kontekstowa | Wskazanie konkretnego źródła jako przedmiotu pytania | Mała rozwijana lista przy polu poleceń | 2 | Znak `@` w polu poleceń | ukryty · rozwinięty (po `@`) | Wstawia odnośnik do pozycji Sources Manager w treści polecenia | Pole poleceń |
| Odnośnik cytowania | Mała plakietka numeryczna w treści odpowiedzi | Wskazanie źródła, z którego pochodzi fragment odpowiedzi | Mała plakietka `.dn-plakietka--info` | 1 | Widoczny w treści odpowiedzi bez interakcji | domyślny · najechanie (podgląd źródła) | Kliknięcie przewija Sources Manager do wskazanej pozycji | W treści odpowiedzi Wykonawcy |
| Menu operacji odpowiedzi `⋮` | Menu kebab zbierające akcje dymka | Grupowanie akcji: Dodaj jako ustalenie, Kopiuj, Regeneruj, Rozgałęź wątek | Mała ikona w rogu dymka | 3 | Kliknięcie ikony `⋮` w dymku odpowiedzi | zwinięte · rozwinięte | Rozwija komplet akcji dotyczących odpowiedzi | Dymek odpowiedzi Wykonawcy |
| Przycisk „Dodaj jako ustalenie” | Pozycja menu operacji odpowiedzi | Przeniesienie odpowiedzi Wykonawcy do Findings Panel | Pozycja menu, mała | 3 | Menu `⋮` dymka odpowiedzi | domyślny · najechanie · dodano (potwierdzenie) | Tworzy nowy wpis w Findings Panel z treścią odpowiedzi | Menu operacji odpowiedzi |
| Przycisk „Przerwij” | Mały przycisk `--zarys` | Przerwanie działania realizowanego przez Wykonawcę | Mały przycisk przy strumieniu odpowiedzi | 1 | Ujawnia się samoczynnie wyłącznie na czas trwania działania (wskaźnik stanu wykonania) | ukryty (brak działania) · aktywny (działanie w toku) | Zatrzymuje bieżące działanie i odnotowuje przerwanie w Execution Loop Window | Przy dymku odpowiedzi w trakcie działania |
| Menu operacji poleceń `⋮` | Menu kebab przy polu poleceń | Grupowanie operacji wejściowych: załącznik, operacje predefiniowane, historia poleceń | Mała ikona przy polu poleceń | 3 | Kliknięcie ikony `⋮` albo znak `/` w polu poleceń | zwinięte · rozwinięte | Rozwija listę operacji i szablonów poleceń | Stopka Chat Window |
| Polecenie „Podsumuj stan badania” | Pozycja menu operacji poleceń | Generuje syntezę dotychczasowych źródeł i ustaleń | Pozycja menu, mała | 3 | Menu operacji poleceń `⋮` albo polecenie języka naturalnego | domyślny | Wysyła predefiniowane polecenie podsumowujące | Menu operacji poleceń |
| Tryby zaawansowane badania | Polecenia eksperckie modułu (tryb przeglądu systematycznego, diagnostyka pętli, konfiguracja dostawców wyszukiwania) | Uruchomienie operacji spoza codziennego przebiegu badania | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w polu poleceń, wyszukiwarka funkcji, skrót klawiszowy, okno konfiguracji | niedostępny dla użytkownika podstawowego · dostępny wg roli | Uruchamia operację i odnotowuje ją w Execution Loop Window | Chat Window, poza interfejsem graficznym |

---

### 3.2. Execution Loop Window (okno pętli wykonawczej)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana |
| Izolacja domyślna | Pętla wykonawcza prowadzona w kontekście bieżącej karty sesji badania |

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą i odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów badawczych. W module Research pętla obejmuje zadania właściwe pracy badawczej: rozpisanie zlecenia badawczego na zapytania wyszukiwawcze, pozyskanie i wczytanie źródeł, ekstrakcję treści i metadanych, ocenę wiarygodności, wydobycie twierdzeń, wykrycie sprzeczności, syntezę ustaleń, złożenie sekcji raportu i wygenerowanie eksportu.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie badawcze i jego dekompozycja na zadania (wyszukaj → pozyskaj → wczytaj → wydobądź → zweryfikuj → zsyntetyzuj → złóż raport).
- Kolejka zadań i stan każdego zadania: oczekuje, w realizacji, zakończone, ponowione, przerwane.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, z treścią zlecenia cząstkowego i raportem zwrotnym.
- Wyniki kontroli jakości: kompletność metadanych źródła, zakotwiczenie ustalenia w cytacie, pokrycie pytań badawczych, kompletność cytowań — wraz z decyzją o ponowieniu zadania.
- Wskaźniki przebiegu pętli: liczba iteracji, czas realizacji, liczba ponowień, stopień ukończenia zlecenia.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia.
- Podgląd zadań zbiorczych: import wsadowy adresów, ponowne przetworzenie kolekcji źródeł, odświeżenie monitorów tematów.
- Ślad prowenancji każdej operacji przekazywany do rejestru ustaleń.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window        │ Execution Loop Window        │ Research
  Użytkownik ↔       │ Koordynator ↔ Wykonawca      │ Workspace
  Wykonawca          │ ──────────────────────────   │ (aktywne okno
                     │ Zlecenie: „analiza           │  wiodące)
  [⟳ Pętla ▼]        │  konkurencyjna segmentu X”   │
                     │ Iteracja 3 · ponowienia: 1   │ Temat badania
                     │ ──────────────────────────   │ etapy · postęp
                     │ ✓ wyszukaj źródła      (14)  │
                     │ ✓ wczytaj i wydobądź   (14)  │
                     │ ● zweryfikuj wiarygodność    │
                     │ ○ zsyntetyzuj ustalenia      │
                     │ ○ złóż sekcje raportu        │
                     │ ──────────────────────────   │
                     │ ⚠ Kontrola jakości: 2 uwagi  │
                     │ ──────────────────────────   │
                     │ [ Sterowanie pętlą ▼ ]       │
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Znacznik stanu pętli | Lekki znacznik `⟳ Pętla ▼` w pasku kontekstu | Sygnalizacja przebiegu pętli wykonawczej i wejście do okna | Bardzo mała plakietka | 1 | Widoczny bez interakcji jako wskaźnik stanu wykonania | pętla bezczynna · w toku · wstrzymana · uchybienia kontroli jakości | Kliknięcie otwiera Execution Loop Window jako kolumnę sąsiadującą | Pasek kontekstu Chat Window |
| Nagłówek zlecenia | Blok tekstowy z treścią zlecenia badawczego | Wskazanie, co realizuje bieżąca pętla | Średni nagłówek z licznikiem iteracji | 1 (w otwartym oknie) | Widoczny po otwarciu okna | domyślny · po korekcie zlecenia | Kliknięcie rozwija pełną treść zlecenia i jego historię | Nagłówek okna |
| Lista zadań | Lista pozycji z oznaczeniem stanu | Prezentacja dekompozycji zlecenia i postępu | Średnia lista `.dn-karta--pozycja` | 1 (w otwartym oknie) | Widoczna po otwarciu okna | oczekuje (○) · w realizacji (●) · zakończone (✓) · ponowione · przerwane | Kliknięcie zadania rozwija komunikaty sterujące i wynik zadania | Ciało okna |
| Komunikat sterujący | Dymek komunikatu Koordynator → Wykonawca lub odwrotnie | Podgląd treści zlecenia cząstkowego i raportu zwrotnego | Mały dymek z oznaczeniem roli | 3 | Kliknięcie pozycji zadania na liście | zwinięty · rozwinięty | Rozwija pełną treść komunikatu z sygnaturą czasu | Wewnątrz rozwiniętego zadania |
| Blok kontroli jakości | Blok `.dn-alert--ostrz` z listą uchybień | Prezentacja wyniku kontroli i decyzji o ponowieniu | Średni blok z kreską lewą ostrzegawczą | 1 (w otwartym oknie; ukryty przy braku uchybień) | Ujawnia się samoczynnie wyłącznie po wykryciu uchybienia | brak uchybień (ukryty) · uchybienia wykryte | Kliknięcie pozycji przenosi do źródła lub ustalenia wymagającego uzupełnienia | Ciało okna, pod listą zadań |
| Wskaźniki przebiegu pętli | Rząd liczników | Liczba iteracji, ponowień, czas realizacji, stopień ukończenia | Mała karta `.dn-karta` z licznikami | 1 (w otwartym oknie) | Widoczne po otwarciu okna | aktualizowane na żywo | Kliknięcie licznika ponowień filtruje listę zadań do ponowionych | Nagłówek okna |
| Element zbiorczy „Sterowanie pętlą ▼” | Grupowanie logiczne akcji sterujących | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia w jednym elemencie | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Sterowanie pętlą ▼`; po wyborze akcji element zwija się samoczynnie | pętla w toku · wstrzymana · przerwana | „Koryguj” otwiera pole edycji zlecenia i uruchamia ponowną dekompozycję | Stopka okna |
| Diagnostyka pętli i śladu prowenancji | Podgląd surowych komunikatów sterujących, sygnatur czasu i pełnego śladu operacji | Analiza przebiegu pętli na poziomie technicznym | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Otwiera pełny zapis przebiegu z możliwością eksportu | Execution Loop Window, tryb administracyjny |

---

### 3.3. Research Workspace

| Aspekt | Wartość |
|---|---|
| Typologia | Punkt wejścia integrujący pozostałe okna (zarządzanie źródłami i ustaleniami) |
| Waga wizualna | Prawa kolumna, dominująca |
| Izolacja domyślna | Struktura badania odrębna per karta sesji |

**Zawartość i pełny arsenał funkcji.**

- Definiowanie zakresu badania: temat, pytania badawcze, granice tematyczne, odbiorca raportu, protokół badania.
- Podział badania na etapy (rozpoznanie → zbieranie → analiza → synteza → raport) z checklistą i własną strukturą etapów.
- Widok osi czasu badania — chronologia dodawania źródeł i ustaleń.
- Widok Evidence Map — graf źródeł, ustaleń i twierdzeń z widocznymi zależnościami, poparciem i sprzecznościami.
- Widok tablicy Kanban etapów badania dla pracy zespołowej nad jednym tematem.
- Widok tabeli dowodów (evidence table) zestawiającej ustalenia ze źródłami i oceną.
- Panel postępu: liczba źródeł, liczba ustaleń, pokrycie pytań badawczych źródłami, procent skompletowania raportu, źródła przeczytane i nieprzeczytane.
- Diagram PRISMA dla trybu przeglądu systematycznego: zidentyfikowane → przesiane → włączone.
- Szybka nawigacja do Discovery Panel, Sources Manager, Reading View, Findings Panel, Report Builder i Export Panel z jednego miejsca.
- Notatka robocza ogólna, niezwiązana z konkretnym źródłem (hipotezy, pytania otwarte).
- Wykrywanie luk badawczych — karta sugestii wskazująca obszary tematu niepokryte dotychczasowymi źródłami.
- Wskaźnik świeżości badania (data ostatniej aktualizacji źródeł, sugestia odświeżenia po czasie).

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window     │ Research Workspace
  Użytkownik ↔    │ Temat: „Analiza konkurencyjna — segment X”     [ ⋮ ]
  Wykonawca       │ Etapy: [Rozpoznanie ✓][Zbieranie ●][Analiza ○][Raport ○]
                  │ ──────────────────────────────────────────────────────
  [⟳ Pętla ▼]     │ [ Widok: Lista ▼ ]
                  │
                  │ ┌─ postęp badania ─────────────────────────────────┐
                  │ │ Źródła: 14   Ustalenia: 22                       │
                  │ │ Pokrycie pytań: 71%   Raport: 40%                │
                  │ └──────────────────────────────────────────────────┘
                  │
                  │ Sugestia: „brak źródeł dot. rynku azjatyckiego”
                  │                                       [ Zbadaj → ]
                  │ ──────────────────────────────────────────────────────
                  │ [ Okna modułu ▼ ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zakresu badania | Pole tekstowe z tytułem | Definicja tematu, pytań badawczych i granic badania | Duży nagłówek | 1 | Widoczny bez interakcji | domyślny · edycja | Kliknięcie otwiera formularz zakresu (temat, pytania, odbiorca, protokół) | Nagłówek okna |
| Pasek etapów | Rząd kroków z oznaczeniem statusu i checklistą | Wizualizacja postępu badania | Średni pasek z plakietkami `.dn-plakietka--stan` | 1 | Widoczny bez interakcji | ukończony (✓) · bieżący (●) · nierozpoczęty (○) | Kliknięcie etapu przewija widok do materiału właściwego etapowi | Pod nagłówkiem |
| Panel postępu | Karta z licznikami | Liczbowy przegląd stanu badania i pokrycia pytań | Mała/średnia karta `.dn-karta` | 1 | Widoczny bez interakcji | domyślny, aktualizowany na żywo | Kliknięcie licznika przenosi do odpowiedniego okna | Ciało okna |
| Karta sugestii | Karta z tekstem i przyciskiem akcji | Wskazanie luki badawczej lub kolejnego kroku | Średnia karta `.dn-karta--akcent` | 1 (ujawniana samoczynnie) | Pojawia się wyłącznie po wykryciu luki badawczej | ukryta (brak sugestii) · widoczna | Kliknięcie „Zbadaj” otwiera Chat Window z wypełnionym poleceniem | Ciało okna |
| Menu zakresu badania `⋮` | Menu kebab nagłówka | Grupowanie akcji zakresu: edycja tematu, pytania badawcze, odbiorca, protokół badania, notatka robocza, wskaźnik świeżości | Mała ikona w nagłówku | 3 | Kliknięcie ikony `⋮` w nagłówku okna | zwinięte · rozwinięte | Rozwija komplet akcji dotyczących zakresu badania | Nagłówek okna |
| Element zbiorczy „Widok ▼” | Menu progresywne reprezentacji materiału | Zmiana reprezentacji materiału badania: lista, oś czasu, Evidence Map, Kanban, tabela dowodów | Mały element zwinięty | 2 | Kliknięcie elementu `Widok ▼`; po wyborze lista zwija się samoczynnie | lista (domyślny) · oś czasu · Evidence Map · Kanban · tabela dowodów | Przełącza sposób prezentacji źródeł i ustaleń | Górna część ciała okna |
| Element zbiorczy „Okna modułu ▼” | Grupowanie logiczne skrótów nawigacyjnych | Przejście do Discovery Panel, Sources Manager, Reading View, Findings Panel, Report Builder i Export Panel | Średni przycisk zbiorczy | 2 | Kliknięcie elementu `Okna modułu ▼`; po wyborze okna element zwija się samoczynnie | domyślny · rozwinięty | Otwiera wskazane okno jako rozszerzenie boczne albo przełącza obszar roboczy | Stopka okna |
| Wskaźnik pokrycia pytań badawczych | Pasek postępu z listą pytań | Wskazanie, które pytania badawcze mają poparcie w źródłach | Średni pasek z rozwijaną listą | 3 | Kliknięcie licznika pokrycia w panelu postępu | pełne pokrycie · częściowe · brak źródeł dla pytania | Kliknięcie pytania filtruje Sources Manager do powiązanych źródeł | Panel postępu |
| Evidence Map (widok grafu) | Interaktywny graf węzłów | Wizualizacja relacji poparcia, sprzeczności i cytowania między źródłami, ustaleniami i twierdzeniami | Duży obszar canvas | 2 | Wybór pozycji „Evidence Map” w elemencie `Widok ▼` | pusty stan · wyrenderowany · powiększony (po kliknięciu węzła) | Kliknięcie węzła otwiera szczegóły źródła lub ustalenia | Ciało okna, w widoku „Evidence Map” |
| Diagram PRISMA | Schemat przepływu z licznikami | Dokumentacja przesiewu w przeglądzie systematycznym | Średni blok SVG | 4 | Włączenie trybu przeglądu systematycznego poleceniem języka naturalnego w Chat Window albo ustawieniem w oknie konfiguracji; niewidoczny dla użytkownika podstawowego | pusty · wypełniony liczbami etapów przesiewu | Kliknięcie etapu pokazuje listę pozycji włączonych i wyłączonych | Ciało okna, w trybie przeglądu systematycznego |

---

### 3.4. Discovery Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Wyszukiwanie i odkrywanie źródeł |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Zapytania i wyniki w kontekście bieżącej karty sesji; dostawcy wyszukiwania i klucze konfigurowane w oknie konfiguracji |

**Zawartość i pełny arsenał funkcji.**

- Zunifikowane pole zapytania z przełącznikiem trybu: web · naukowe · własne semantyczne · pełnotekstowe (FTS).
- Filtry wyników: rok, dziedzina, dostęp otwarty, typ pozycji, dostawca.
- Query Assistant — przekształcenie pytania badawczego w zestaw zapytań wyszukiwawczych z operatorami i filtrami.
- Lista wyników z metadanymi (tytuł, autorzy, rok, źródło, fragment) i akcją „→ Dodaj do źródeł”.
- Rozstrzyganie identyfikatorów DOI / ISBN / PMID do pełnych metadanych pozycji.
- Panel snowballingu — lista prac cytowanych i cytujących wybraną pozycję, z przejściem do Evidence Map.
- Zakładka monitorów tematów i kanałów RSS/Atom ze skrzynką „nowe źródła”.
- Import wsadowy listy adresów przekazywany do kolejki zadań pętli wykonawczej.
- Oznaczenie wyniku jako odrzuconego z uzasadnieniem — pozycja zasila liczniki diagramu PRISMA.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window     │ Research Workspace       │ Discovery Panel
  Użytkownik ↔    │ (aktywne okno wiodące)   │ ─────────────────────────
  Wykonawca       │                          │ [ Tryb: Web ▼ ]
                  │ Temat badania            │ Zapytanie: „udział rynkowy
  [⟳ Pętla ▼]     │ etapy · postęp           │  segment X 2024”     [ ▶ ]
                  │                          │ [ Filtry ▼ ]
                  │                          │ ─────────────────────────
                  │                          │ ▸ Raport rynkowy 2024
                  │                          │   Crossref · DOI 10.xxxx
                  │                          │   [ → Dodaj ]     [ ⋮ ]
                  │                          │ ▸ Artykuł przeglądowy
                  │                          │   OpenAlex · 2023
                  │                          │   [ → Dodaj ]     [ ⋮ ]
                  │                          │ ─────────────────────────
                  │                          │ Monitory / RSS: 6 nowych
                  │                          │ [ Operacje panelu ▼ ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Pole zapytania | Pole `.dn-input` z przyciskiem | Wprowadzenie zapytania wyszukiwawczego | Średnie pole tekstowe | 1 (w otwartym panelu) | Widoczne po otwarciu panelu | puste · z zapytaniem · wyszukiwanie w toku | Enter uruchamia wyszukiwanie u wybranych dostawców | Nagłówek panelu |
| Lista wyników | Lista pozycji z metadanymi | Prezentacja trafień wyszukiwania | Duży obszar, dominujący w panelu | 1 (w otwartym panelu) | Widoczna po uruchomieniu wyszukiwania | pusty stan · z wynikami · ładowanie | Przewijanie listy, rozwinięcie pozycji | Ciało panelu |
| Element zbiorczy „Tryb ▼” | Menu progresywne trybu wyszukiwania | Wybór między wyszukiwaniem webowym, naukowym, semantycznym i pełnotekstowym | Mały element zwinięty | 2 | Kliknięcie elementu `Tryb ▼`; po wyborze element zwija się samoczynnie | web (domyślny) · naukowe · semantyczne · FTS | Zmienia dostawcę zapytania i zestaw dostępnych filtrów | Nagłówek panelu |
| Element zbiorczy „Filtry ▼” | Grupowanie logiczne kontrolek zawężających | Zawężenie wyników wg roku, dziedziny, dostępu otwartego, typu, dostawcy | Mały element zwinięty | 3 | Kliknięcie elementu `Filtry ▼` — rozwija panel popover z kompletem kontrolek | zwinięty (bez filtrów) · zwinięty z licznikiem aktywnych filtrów · rozwinięty | Zawęża listę wyników bez ponownego wpisywania zapytania | Pod polem zapytania |
| Pozycja wyniku | Wiersz karty `.dn-karta--pozycja` | Reprezentuje pojedynczy wynik z metadanymi i dostawcą | Średni wiersz | 1 (w otwartym panelu) | Widoczna po uruchomieniu wyszukiwania | nowy · dodany do źródeł · odrzucony · duplikat wykryty | Rozwija akcje pozycji | Lista główna panelu |
| Przycisk „→ Dodaj do źródeł” | Mały przycisk `--zloty` | Przeniesienie pozycji do Sources Manager z metadanymi | Mały przycisk CTA | 1 (przy pozycji wyniku) | Widoczny przy każdej pozycji jako akcja podstawowa | domyślny · dodano | Tworzy pozycję w Sources Manager i uruchamia pozyskanie pełnego tekstu | Pozycja wyniku |
| Menu pozycji wyniku `⋮` | Menu kebab pozycji | Grupowanie akcji dodatkowych: Cytowania (snowballing), Odrzuć z uzasadnieniem, Podgląd, Rozstrzygnij identyfikator | Mała ikona przy pozycji | 3 | Kliknięcie ikony `⋮` przy pozycji wyniku | zwinięte · rozwinięte | Rozwija komplet akcji pozycji | Prawa krawędź pozycji wyniku |
| Menu „Cytowania” | Rozwijane menu snowballingu | Podgląd prac cytowanych i cytujących | Małe menu z dwiema listami | 3 | Pozycja „Cytowania” w menu `⋮` pozycji wyniku | zwinięte · rozwinięte · ładowanie relacji | Wybór pozycji dodaje ją do wyników albo otwiera Evidence Map | Menu pozycji wyniku |
| Karta Query Assistant | Karta `.dn-karta--akcent` z sugerowanymi zapytaniami | Przekształcenie pytania badawczego w zapytania z operatorami | Średnia karta | 2 | Kliknięcie ikony asystenta przy polu zapytania albo polecenie w Chat Window; po zastosowaniu karta zwija się samoczynnie | ukryta · widoczna · zastosowana | Kliknięcie wstawia zapytanie do pola i uruchamia wyszukiwanie | Pod polem zapytania |
| Zakładka monitorów i RSS | Lista kanałów z licznikiem nowych pozycji | Przegląd nowych publikacji z monitorowanych tematów i kanałów | Średnia lista `.dn-karta--pozycja` | 2 | Kliknięcie licznika „nowe źródła” w stopce panelu | brak nowych · nowe pozycje · odświeżanie | Kliknięcie pozycji przenosi ją do listy wyników z akcją dodania | Stopka panelu |
| Element zbiorczy „Operacje panelu ▼” | Grupowanie logiczne operacji zbiorczych | Import wsadowy adresów, konfiguracja monitorów tematów, odświeżenie kanałów, zarządzanie dostawcami wyszukiwania | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Operacje panelu ▼`; po wyborze akcji element zwija się samoczynnie | domyślny · rozwinięty | Rozwija listę operacji zbiorczych panelu | Stopka panelu |
| Import wsadowy adresów | Pole wklejenia listy adresów do przetworzenia w tle | Zbiorcze pozyskanie wielu źródeł jednym zleceniem | Panel wywoływany | 4 | Pozycja w elemencie `Operacje panelu ▼`, polecenie języka naturalnego w Chat Window albo skrót klawiszowy | domyślny · przetwarzanie (postęp w Execution Loop Window) | Tworzy zadanie zbiorcze w pętli wykonawczej, status per pozycja | Operacje panelu |
| Konfiguracja dostawców wyszukiwania i kluczy | Ustawienia dostawców zapytań i kluczy dostępu | Zmiana zestawu dostawców webowych i naukowych oraz kluczy jawnych | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), polecenie języka naturalnego w Chat Window, tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Zmienia zestaw dostawców obsługujących pole zapytania | Poza interfejsem panelu |

---

### 3.5. Sources Manager

| Aspekt | Wartość |
|---|---|
| Typologia | Zarządzanie źródłami i ustaleniami |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Lista źródeł odrębna per karta sesji; przyjmuje źródła z Browser i Discovery Panel |

**Zawartość i pełny arsenał funkcji.**

- Dodawanie źródła ręcznie (adres URL, plik, cytat, identyfikator DOI/ISBN/PMID), z importu zbiorczego, z monitorowanego folderu oraz z modułu Browser i Discovery Panel.
- Biblioteka referencji z pełnymi metadanymi w schemacie CSL-JSON: autor, rok, wydawca, DOI, typ pozycji, pochodzenie, data pozyskania.
- Ocena wiarygodności źródła (skala punktowa lub etykieta: pierwotne / wtórne / nieporównane) — ustalana ręcznie lub sugerowana według kryteriów recenzji, cytowalności i domeny.
- Katalogowanie: tagi tematyczne, kolekcje swobodne i kolekcje inteligentne oparte na regułach, przypisanie do etapu badania i do konkretnego pytania badawczego.
- Deduplikacja: wykrywanie duplikatów i niemal-duplikatów (ten sam DOI/URL, zbliżona treść) z sugestią scalenia — dodanie źródła pozostaje możliwe.
- Import zbiorczy listy źródeł (BibTeX, RIS, CSL-JSON, EndNote XML, CSV) oraz eksport bibliografii w tych formatach i w postaci sformatowanej listy w wybranym stylu CSL.
- Menedżer załączników: pełny tekst PDF, migawka strony, notatki źródła, kontrola kompletności.
- Podgląd migawki treści źródła bez opuszczania modułu oraz przejście do Reading View.
- Wyszukiwanie i filtrowanie: po typie, wiarygodności, tagu, dacie, etapie badania, stanie lektury.
- Automatyczne generowanie cytatu w wybranym stylu CSL na potrzeby Report Builder.
- Citation Health Check: kontrola kompletności metadanych do cytowania, wykrywanie braków i uzupełnienie z Crossref.
- Flaga wycofania lub korekty pracy na podstawie danych Retraction Watch.
- Oznaczenie źródła jako „do weryfikacji” lub „zweryfikowane” oraz oznaczenie stanu lektury.
- Powiązanie źródła z jednym lub wieloma wpisami Findings Panel.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window     │ Research Workspace       │ Sources Manager
  Użytkownik ↔    │ (aktywne okno wiodące)   │ ─────────────────────────
  Wykonawca       │                          │ [ + Dodaj źródło ]
                  │ Temat badania            │ [ 🔍 filtr/szukaj ]
  [⟳ Pętla ▼]     │ etapy · postęp           │ ─────────────────────────
                  │                          │ ● [3] raport-rynkowy.pdf
                  │                          │   wiarygodność: wysoka
                  │                          │   🏷rynek · przeczytane
                  │                          │   ustalenia: 4
                  │                          │   [ Czytaj ]      [ ⋯ ]
                  │                          │ ─────────────────────────
                  │                          │ ○ [7] wywiad-branzowy.mp4
                  │                          │   wiarygodność: średnia
                  │                          │   [ Czytaj ]      [ ⋯ ]
                  │                          │ ─────────────────────────
                  │                          │ ⚠ [9] artykul-x.html
                  │                          │   do weryfikacji · brak roku
                  │                          │ ─────────────────────────
                  │                          │ [ Operacje panelu ▼ ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Lista źródeł | Lista pozycji z metadanymi | Katalog źródeł zgromadzonych w badaniu | Duży obszar, dominujący w panelu | 1 (w otwartym panelu) | Widoczna po otwarciu panelu | pusty stan · z pozycjami | Przewijanie listy, rozwinięcie pozycji | Ciało panelu |
| Przycisk „+ Dodaj źródło” | Przycisk `--zloty`, mały | Dodanie nowego źródła ręcznie lub z identyfikatora | Mały przycisk CTA | 1 (w otwartym panelu) | Widoczny po otwarciu panelu jako akcja podstawowa | domyślny · otwarty formularz | Otwiera formularz z polami URL/plik/cytat/DOI oraz metadanymi CSL | Nagłówek panelu |
| Pole filtr/szukaj | Pole `.dn-input` z ikoną (filtr) | Zawężenie listy źródeł | Małe pole tekstowe | 1 (w otwartym panelu) | Widoczne po otwarciu panelu | puste · z wynikami · brak trafień | Filtruje listę na żywo wg wpisanej frazy lub aktywnych filtrów | Nagłówek panelu |
| Pozycja źródła | Wiersz karty `.dn-karta--pozycja` | Reprezentuje pojedyncze źródło z metadanymi | Średni wiersz z kropką wiarygodności, tytułem, metadanymi | 1 (w otwartym panelu) | Widoczna po otwarciu panelu | zweryfikowane (kropka zielona) · do weryfikacji (kropka ostrzegawcza) · wycofane (flaga) · najechanie | Rozwija akcje pozycji | Lista główna panelu |
| Numer odnośnika | Mała plakietka liczbowa | Numer cytowania używany w Chat Window i Report Builder | Mała plakietka `.dn-plakietka` | 1 | Widoczny przy pozycji bez interakcji | domyślny | Kliknięcie kopiuje odnośnik do schowka | Przy każdej pozycji źródła |
| Etykiety tematyczne | Zestaw małych plakietek | Tagi i kolekcje przypisane źródłu | Bardzo małe pigułki `.dn-plakietka--wersaliki` | 1 | Widoczne przy pozycji bez interakcji | brak tagów · z tagami · kolekcja inteligentna | Kliknięcie tagu filtruje listę do tego tagu | Pod tytułem pozycji źródła |
| Wskaźnik stanu lektury | Mała plakietka tekstowa | Informacja, czy źródło przeczytano, przejrzano czy pominięto | Bardzo mała plakietka | 1 | Widoczny przy pozycji bez interakcji | nieprzeczytane · przejrzane · przeczytane · pominięte | Kliknięcie zmienia stan lektury i aktualizuje panel postępu | Przy pozycji źródła |
| Przycisk „Czytaj” | Mały przycisk `--zarys` | Otwarcie źródła w Reading View | Mały przycisk tekstowy | 1 (przy pozycji źródła) | Widoczny przy każdej pozycji jako akcja podstawowa | domyślny · ładowanie | Otwiera Reading View jako rozszerzenie boczne z wybranym źródłem | Pozycja źródła |
| Alert Citation Health Check | Blok `.dn-alert--ostrz` | Sygnalizacja braków metadanych, niecytowanych źródeł i cytowań bez wpisu | Średni blok z kreską lewą ostrzegawczą | 1 (ujawniany samoczynnie) | Pojawia się wyłącznie przy pozycji z niekompletnymi metadanymi | ukryty (metadane kompletne) · widoczny | „Uzupełnij metadane” pobiera brakujące pola z Crossref | Przy pozycji z brakami |
| Menu pozycji źródła `⋯` | Ikona rozwijanego menu | Grupowanie akcji: Cytuj, Powiąż z ustaleniem, Załączniki, Scalenie duplikatu, Flaga wycofania, Usunięcie z listy | Mała ikona (wiecej) | 3 | Kliknięcie ikony `⋯` przy pozycji źródła | zwinięte · rozwinięte | Otwiera menu kontekstowe pozycji | Prawa krawędź pozycji źródła |
| Przycisk „Cytuj” | Pozycja menu pozycji źródła | Generowanie gotowego cytatu w wybranym stylu | Pozycja menu, mała | 3 | Menu `⋯` pozycji źródła albo skrót klawiszowy | domyślny · skopiowano | Kopiuje sformatowany cytat do schowka | Menu pozycji źródła |
| Selektor stylu cytatu | Rozwijana lista stylów CSL | Wybór formatu generowanego cytatu | Mały element zwinięty `Styl cytatu ▼` | 2 | Kliknięcie znacznika stylu w nagłówku panelu albo pozycja w elemencie `Operacje panelu ▼`; po wyborze lista zwija się samoczynnie | domyślny (APA) · rozwinięty · styl własny | Zmienia format cytatów i bibliografii końcowej | Nagłówek panelu |
| Element zbiorczy „Operacje panelu ▼” | Grupowanie logiczne operacji panelu | Eksport bibliografii, import zbiorczy referencji, deduplikacja, styl cytatu, kolekcje | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Operacje panelu ▼`; po wyborze akcji element zwija się samoczynnie | domyślny · rozwinięty | Rozwija listę operacji panelu | Stopka panelu |
| Eksport bibliografii | Operacja pobrania listy źródeł w formacie BibTeX/RIS/CSL-JSON lub jako sformatowanej listy | Przekazanie bibliografii poza moduł | Panel wywoływany | 3 | Pozycja w elemencie `Operacje panelu ▼` | domyślny · ładowanie | Generuje plik bibliografii do pobrania | Operacje panelu |
| Reguły kolekcji inteligentnych i progi deduplikacji | Definicje regułowe kolekcji, progi podobieństwa i kryteria punktacji wiarygodności | Sterowanie zachowaniem katalogowania i wykrywania duplikatów | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), polecenie języka naturalnego w Chat Window, tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Zmienia sposób przypisywania źródeł do kolekcji i wykrywania duplikatów | Poza interfejsem panelu |
| Monitorowany folder importu | Wskazanie katalogu, z którego nowe dokumenty wchodzą jako źródła | Automatyczne zasilanie badania dokumentami z dysku | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), polecenie języka naturalnego w Chat Window | niedostępny dla użytkownika podstawowego · dostępny wg roli | Uruchamia obserwację katalogu i ekstrakcję metadanych | Poza interfejsem panelu |

---

### 3.6. Reading View

| Aspekt | Wartość |
|---|---|
| Typologia | Okno lektury i ekstrakcji |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne, szeroka |
| Izolacja domyślna | Adnotacje i wypisy zapisywane przy źródle bieżącej karty sesji |

**Zawartość i pełny arsenał funkcji.**

- Renderowanie źródła w trybie lektury: PDF, HTML, DOCX, EPUB, Markdown, transkrypt ze znacznikami czasu.
- Miniatury stron, nawigacja po strukturze dokumentu, wyszukiwanie w treści z podświetleniem trafień.
- Podświetlenia w kolorach, notatki na marginesie i zakładki, kotwiczone do pozycji (strona i przesunięcie dla PDF, selektor CSS dla HTML, znacznik czasu dla transkryptu).
- Zamiana podświetlenia w ustalenie jednym działaniem, z zachowaniem odnośnika do fragmentu.
- Panel wypisów zbierający wszystkie podświetlenia i notatki z jednego lub wielu źródeł.
- Streszczenie pojedynczego źródła: abstrakt roboczy, tezy, metodologia, wnioski.
- Ekstrakcja tabel i danych liczbowych do postaci ustrukturyzowanej gotowej do wykresu lub arkusza.
- Wydobycie kluczowych twierdzeń, danych liczbowych i podmiotów (osoby, organizacje, daty) z sugestią zapisu jako ustalenia.
- Cross-Document Chat: rozmowa oparta na treści wybranych źródeł, z odpowiedziami zakotwiczonymi w cytatach i numerach stron.
- OCR skanów i obrazów czyniący dokument przeszukiwalnym i cytowalnym.
- Oznaczenie stanu lektury źródła, przenoszone do panelu postępu w Research Workspace.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window   │ Research        │ Reading View — raport-rynkowy.pdf
  Użytkownik ↔  │ Workspace       │ ──────────────────────────────────
  Wykonawca     │ (aktywne okno   │ [Miniatury] │ str. 12 / 48
                │  wiodące)       │  ▤ 10       │
  [⟳ Pętla ▼]   │                 │  ▤ 11       │ „Segment X rośnie
                │ Temat badania   │  ▤ 12 ◄     │  o 12% rocznie od
                │ etapy · postęp  │  ▤ 13       │  2023 roku” ▮▮▮▮
                │                 │             │  ← podświetlenie
                │                 │             │
                │                 │ ──────────────────────────────────
                │                 │ [ → ustalenie ]  [ Operacje ▼ ]
                │                 │ ──────────────────────────────────
                │                 │ [ Wypisy (7) ▼ ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Obszar treści źródła | Renderowany dokument | Lektura treści źródła | Duży obszar, dominujący w oknie | 1 (w otwartym oknie) | Widoczny po otwarciu okna | ładowanie · wyrenderowany · skan bez warstwy tekstowej | Zaznaczenie fragmentu wywołuje pasek akcji podświetlenia | Ciało okna |
| Pas miniatur stron | Wąska lista miniatur | Nawigacja po stronach lub sekcjach dokumentu | Wąska kolumna wewnątrz okna | 1 (w otwartym oknie) | Widoczny po otwarciu okna | domyślny · strona bieżąca wyróżniona | Kliknięcie miniatury przenosi widok do strony | Lewa krawędź okna |
| Podświetlenie | Zaznaczony fragment w kolorze | Trwałe oznaczenie istotnego fragmentu | Kolorowe tło fragmentu | 1 | Widoczne w treści bez interakcji | cztery kolory znaczników · z notatką · bez notatki | Kliknięcie otwiera notatkę na marginesie i akcję „→ ustalenie” | W obszarze treści źródła |
| Przycisk „→ ustalenie” | Mały przycisk `--zloty` | Zamiana podświetlenia w ustalenie z odnośnikiem do fragmentu | Mały przycisk CTA | 1 (akcja podstawowa okna) | Widoczny po zaznaczeniu fragmentu | domyślny · dodano | Tworzy wpis w Findings Panel z cytatem i kotwicą | Pasek akcji podświetlenia |
| Pasek akcji podświetlenia | Menu kontekstowe zaznaczenia | Akcje na zaznaczonym fragmencie: kolor znacznika, notatka, kopiowanie cytatu | Mały pasek przy zaznaczeniu | 3 | Zaznaczenie fragmentu treści; po użyciu pasek znika | ukryty (brak zaznaczenia) · widoczny | Realizuje akcję na zaznaczeniu i zapisuje kotwicę pozycji | Przy zaznaczeniu w treści źródła |
| Notatka na marginesie | Mała karta przy krawędzi treści | Komentarz roboczy do podświetlonego fragmentu | Mała karta `.dn-karta` | 3 | Kliknięcie podświetlenia albo pozycja w pasku akcji podświetlenia | zwinięta · rozwinięta · edycja | Zapis treści notatki wraz z kotwicą pozycji | Prawa krawędź obszaru treści |
| Element zbiorczy „Operacje ▼” | Grupowanie logiczne operacji na źródle | Streszczenie źródła, ekstrakcja tabeli, wydobycie twierdzeń, zapytanie o źródło (Cross-Document Chat), OCR, stan lektury | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Operacje ▼`; po wyborze akcji element zwija się samoczynnie | domyślny · ładowanie operacji | Uruchamia operację jako zadanie widoczne w Execution Loop Window | Pod obszarem treści |
| Element zbiorczy „Wypisy ▼” | Zwinięta lista podświetleń i notatek z licznikiem | Zbiorczy przegląd materiału wydobytego ze źródła lub kolekcji | Średni element zwinięty | 2 | Kliknięcie elementu `Wypisy (n) ▼`; po zamknięciu lista zwija się do licznika | pusty stan · z wypisami · grupowanie po źródle lub tagu | Kliknięcie wypisu przenosi widok do kotwicy w treści | Stopka okna |
| Wynik ekstrakcji tabeli | Podgląd tabeli ustrukturyzowanej | Prezentacja danych wydobytych ze strony dokumentu | Średnia tabela z akcją zapisu | 3 | Pozycja „Wyodrębnij tabelę” w elemencie `Operacje ▼` | wykryta · zweryfikowana · zapisana jako dane | Zapis tworzy artefakt danych dostępny dla wykresów i eksportu XLSX | Nakładka nad obszarem treści |
| Wskaźnik OCR | Plakietka stanu przetwarzania | Informacja o rozpoznawaniu tekstu w skanie | Mała plakietka `.dn-plakietka--stan` | 1 (ujawniana samoczynnie) | Pojawia się wyłącznie przy skanie bez warstwy tekstowej i na czas przetwarzania | niepotrzebny (ukryty) · w toku · zakończony | Po zakończeniu treść skanu wchodzi do wyszukiwania pełnotekstowego | Nagłówek okna |
| Cross-Document Chat na wyborze źródeł | Rozmowa oparta na treści wielu wskazanych źródeł z odpowiedziami zakotwiczonymi w cytatach | Analiza porównawcza korpusu źródeł badania | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, skrót klawiszowy | niedostępny dla użytkownika podstawowego · dostępny wg roli | Zwraca odpowiedź z odnośnikami do stron i fragmentów | Chat Window i Reading View |
| Parametry OCR i przetwarzania wstępnego | Ustawienia pakietów językowych i obróbki obrazu przed rozpoznaniem tekstu | Dostosowanie rozpoznawania tekstu do materiału źródłowego | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Zmienia zachowanie operacji OCR | Poza interfejsem okna |

---

### 3.7. Findings Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Zarządzanie źródłami i ustaleniami |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Ustalenia narastają w toku sesji, powiązane ze źródłami z Sources Manager |

**Zawartość i pełny arsenał funkcji.**

- Zapis ustalenia cząstkowego — wniosku, obserwacji, cytatu — z powiązaniem z jednym lub wieloma źródłami i kotwicą do fragmentu.
- Klasyfikacja ustalenia: fakt, hipoteza, opinia eksperta, dana liczbowa, cytat bezpośredni.
- Ocena wagi ustalenia (kluczowe / poboczne) wpływająca na kolejność w Report Builder.
- Kodowanie jakościowe: kody tematyczne przypisywane fragmentom, książka kodów, macierz kod × źródło, zliczanie wystąpień.
- Wykrywanie sprzeczności między ustaleniami pochodzącymi z różnych źródeł, z heurystyką rozbieżności liczbowych.
- Weryfikacja twierdzenia względem dostępnych źródeł i wyszukiwania webowego, ze wskazaniem poparcia, braku poparcia lub sprzeczności.
- Grupowanie ustaleń w wątki tematyczne — swobodnie definiowane przez użytkownika oraz automatyczne klastrowanie.
- Scalanie powtórzonych ustaleń z różnych źródeł w jedno z wieloma odnośnikami.
- Edycja treści ustalenia, dodawanie komentarza roboczego, oznaczanie jako „wymaga potwierdzenia”.
- Widok chronologiczny, widok grupowany tematycznie i widok kodowania jakościowego.
- Przeciągnięcie ustalenia bezpośrednio do struktury Report Builder.
- Wyszukiwanie pełnotekstowe w treści ustaleń.
- Pełny ślad prowenancji ustalenia: źródło, fragment, autor zmiany (Użytkownik albo Wykonawca), znacznik czasu.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window     │ Research Workspace      │ Findings Panel
  Użytkownik ↔    │ (aktywne okno wiodące)  │ ──────────────────────────
  Wykonawca       │                         │ [ + Nowe ustalenie ]
                  │ Temat badania           │ [ Widok: Chronologia ▼ ]
  [⟳ Pętla ▼]     │ etapy · postęp          │ ──────────────────────────
                  │                         │ ★ Kluczowe · fakt  [3][7]
                  │                         │  „Segment X rośnie o 12%
                  │                         │   rocznie od 2023 roku”
                  │                         │  kody: #wzrost #rynek
                  │                         │  [→ Report Builder] [ ⋮ ]
                  │                         │ ──────────────────────────
                  │                         │ ⚠ Sprzeczność  [7] vs [9]
                  │                         │  „rozbieżne dane udziału”
                  │                         │                     [ ⋮ ]
                  │                         │ ──────────────────────────
                  │                         │ [ Operacje panelu ▼ ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Lista ustaleń | Lista kart ustaleń | Rejestr ustaleń cząstkowych badania | Duży obszar, dominujący w panelu | 1 (w otwartym panelu) | Widoczna po otwarciu panelu | pusty stan · z ustaleniami | Przewijanie listy, rozwinięcie karty | Ciało panelu |
| Przycisk „+ Nowe ustalenie” | Przycisk `--zloty`, mały | Ręczne dodanie ustalenia | Mały przycisk CTA | 1 (w otwartym panelu) | Widoczny po otwarciu panelu jako akcja podstawowa | domyślny | Otwiera pole edycji nowego wpisu | Nagłówek panelu |
| Pozycja ustalenia | Karta `.dn-karta--pozycja` | Reprezentuje pojedyncze ustalenie | Średnia karta z oznaczeniem wagi i typu | 1 (w otwartym panelu) | Widoczna po otwarciu panelu | kluczowe (gwiazdka) · poboczne · sprzeczność (plakietka ostrzegawcza) · wymaga potwierdzenia · scalone | Kliknięcie rozwija treść, kody, ślad prowenancji i akcje | Lista główna panelu |
| Plakietka typu ustalenia | Mała etykieta | Klasyfikacja: fakt / hipoteza / opinia / dana liczbowa / cytat | Mała pigułka `.dn-plakietka` | 1 | Widoczna przy karcie bez interakcji | domyślny | Kliknięcie otwiera listę klasyfikacji i pozwala ją zmienić | Nagłówek karty ustalenia |
| Odnośniki do źródeł | Zestaw małych plakietek numerycznych | Wskazanie źródeł, na których oparte jest ustalenie | Bardzo małe plakietki `.dn-plakietka--info` | 1 | Widoczne przy karcie bez interakcji | jedno źródło · wiele źródeł | Kliknięcie przewija do pozycji w Sources Manager lub do kotwicy w Reading View | Nagłówek karty ustalenia |
| Kody tematyczne | Zestaw pigułek z kodami | Kodowanie jakościowe fragmentu | Bardzo małe pigułki `.dn-plakietka--wersaliki` | 1 | Widoczne przy karcie bez interakcji | bez kodów · z kodami · kod sugerowany | Kliknięcie kodu filtruje listę; menu otwiera książkę kodów | Ciało karty ustalenia |
| Alert sprzeczności | Blok `.dn-alert--ostrz` | Sygnalizacja rozbieżności między źródłami | Średni blok z kreską lewą ostrzegawczą | 1 (ujawniany samoczynnie) | Pojawia się wyłącznie po wykryciu sprzeczności | ukryty (brak sprzeczności) · widoczny · rozstrzygnięta | „Porównaj źródła” otwiera zestawienie różnicowe dwóch źródeł | Nad kartą ustalenia, gdy wykryto sprzeczność |
| Przycisk „→ Report Builder” | Mały przycisk `--zarys` | Przeniesienie ustalenia do struktury raportu | Mały przycisk tekstowy | 1 (przy karcie ustalenia) | Widoczny przy karcie jako akcja podstawowa | domyślny · dodano | Wstawia ustalenie jako sekcję lub akapit w Report Builder z przypisem | Karta ustalenia |
| Element zbiorczy „Widok ▼” | Menu progresywne reprezentacji ustaleń | Chronologia / wątki tematyczne / kodowanie jakościowe | Mały element zwinięty | 2 | Kliknięcie elementu `Widok ▼`; po wyborze element zwija się samoczynnie | chronologia (domyślny) · wątki · kodowanie | Zmienia sposób prezentacji listy ustaleń | Nagłówek panelu |
| Menu karty ustalenia `⋮` | Menu kebab karty | Grupowanie akcji: Edytuj, Fact-Check, Zmień wagę, Oznacz „wymaga potwierdzenia”, Oznacz sprzeczność rozstrzygniętą, Powiąż ze źródłem, Prowenancja | Mała ikona przy karcie | 3 | Kliknięcie ikony `⋮` przy karcie ustalenia | zwinięte · rozwinięte | Rozwija komplet akcji karty | Prawa krawędź karty ustalenia |
| Przycisk „Fact-Check” | Pozycja menu karty ustalenia | Weryfikacja twierdzenia względem korpusu źródeł i wyszukiwania | Pozycja menu, mała | 3 | Menu `⋮` karty ustalenia albo polecenie w Chat Window | domyślny · weryfikacja w toku · wynik (poparcie / brak / sprzeczność) | Zapisuje wynik jako adnotację ustalenia z odnośnikami | Menu karty ustalenia |
| Blok prowenancji | Lista wpisów śladu | Pochodzenie i historia zmian ustalenia | Mała lista z sygnaturami czasu | 3 | Pozycja „Prowenancja” w menu `⋮` karty ustalenia | zwinięty · rozwinięty | Rozwija pełny ślad: źródło, fragment, autor zmiany, czas | Ciało karty ustalenia |
| Element zbiorczy „Operacje panelu ▼” | Grupowanie logiczne operacji panelu | Grupowanie automatyczne w wątki, scalanie duplikatów, książka kodów, wyszukiwanie pełnotekstowe w ustaleniach | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Operacje panelu ▼`; po wyborze akcji element zwija się samoczynnie | domyślny · rozwinięty · ładowanie operacji | Rozwija listę operacji panelu i uruchamia wybraną | Stopka panelu |
| Macierz kod × źródło | Tabela zliczeń wystąpień | Analiza jakościowa: rozkład kodów w źródłach | Średnia tabela przewijana | 4 | Widok „Kodowanie” uruchamiany poleceniem języka naturalnego w Chat Window, skrótem klawiszowym albo z książki kodów; niewidoczny dla użytkownika podstawowego | pusta · wypełniona | Kliknięcie komórki pokazuje fragmenty objęte kodem; eksport do XLSX | Ciało panelu, w widoku „Kodowanie” |
| Parametry wykrywania sprzeczności i klastrowania | Progi rozbieżności liczbowych i kryteria grupowania tematycznego | Sterowanie czułością detektora sprzeczności i grupowania ustaleń | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Zmienia zachowanie wykrywania sprzeczności i grupowania | Poza interfejsem panelu |

---

### 3.8. Report Builder

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne / kreator |
| Waga wizualna | Prawa kolumna, dominująca po aktywacji |
| Izolacja domyślna | Buduje dokument na podstawie zawartości Findings Panel bieżącej sesji |

**Zawartość i pełny arsenał funkcji.**

- Struktura raportu w postaci konspektu (sekcje, podsekcje) — budowana ręcznie lub sugerowana na podstawie ustaleń.
- Wstawianie ustaleń z Findings Panel do wybranej sekcji przeciągnięciem lub poleceniem, z automatycznym przypisem CSL.
- Szablony raportu: streszczenie zarządcze, analiza konkurencyjna, analiza SWOT, raport benchmarkowy, nota badawcza, przegląd literatury.
- Składanie przeglądu literatury: tabela dowodów, synteza narracyjna, luki badawcze, protokół PRISMA z diagramem przepływu.
- Automatyczne generowanie streszczenia zarządczego na podstawie kluczowych ustaleń.
- Wstawianie tabel, macierzy porównawczej, osi czasu, wykresów poglądowych z danych liczbowych z ustaleń oraz cytatów źródłowych z automatycznym przypisem.
- Zarządzanie przypisami dolnymi i końcowymi: numeracja, przypisy skrócone (ibid., op. cit.), renumeracja przy edycji.
- Bibliografia końcowa generowana z Sources Manager w wybranym stylu CSL, z wyborem zakresu: wszystkie źródła albo tylko cytowane.
- Wstawianie odnośnika w tekście `[n]` albo (Autor, rok) z jednoczesnym wpisem do bibliografii.
- Operacje kontekstowe analogiczne do modułu Studio: korekta, zmiana stylu, streszczenie, rozwinięcie — na zaznaczonym fragmencie raportu.
- Widok podzielony: konspekt w lewej kolumnie okna, treść sekcji w prawej.
- Licznik pokrycia ustaleń (ile z zebranych ustaleń trafiło do raportu, ile pozostaje poza dokumentem).
- Wersjonowanie raportu — kolejne kompletacje jako odrębne wersje z porównaniem różnic.
- Tryb recenzji: komentarze do fragmentów raportu przed finalizacją, z wątkami dyskusji.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window   │ Report Builder (aktywne okno wiodące)
  Użytkownik ↔  │ [ Szablon: Analiza konkurencyjna ▼ ]   Pokrycie: 78%
  Wykonawca     │ Wersja: 3                                       [ ⋮ ]
                │ ────────────────────┬─────────────────────────────────
  [⟳ Pętla ▼]   │ KONSPEKT            │ TREŚĆ SEKCJI
                │ 1. Streszczenie ✓   │ ## 2. Kontekst
                │ 2. Kontekst    ●    │ Segment X rośnie
                │ 3. Konkurenci  ○    │ o 12% rocznie [3] …
                │ 4. Wnioski     ○    │
                │                     │ [przeciągnij ustalenie tutaj]
                │ [ + Dodaj sekcję ]  │
                │ ────────────────────┴─────────────────────────────────
                │ [ Operacje raportu ▼ ]              [ → Export Panel ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Panel konspektu | Lista sekcji z drzewem podsekcji | Nawigacja i struktura dokumentu | Wąska kolumna wewnątrz okna | 1 | Widoczny bez interakcji po aktywacji okna wiodącego | ukończona sekcja (✓) · bieżąca (●) · pusta (○) | Kliknięcie sekcji przenosi widok treści do niej | Lewa kolumna okna |
| Obszar treści sekcji | Edytowalny blok tekstu | Redakcja treści wybranej sekcji | Duży blok, prawa kolumna okna | 1 | Widoczny bez interakcji po aktywacji okna wiodącego | pusty (podpowiedź przeciągnięcia ustalenia) · z treścią | Edycja tekstu, przeciąganie ustaleń, zaznaczanie do operacji kontekstowych | Prawa kolumna okna |
| Wskaźnik pokrycia ustaleń | Mała plakietka procentowa | Informacja, ile ustaleń trafiło do raportu | Mała plakietka `.dn-plakietka--info` | 1 | Widoczny bez interakcji jako wskaźnik stanu | aktualizowana na żywo | Kliknięcie pokazuje listę ustaleń jeszcze niewykorzystanych | Nagłówek okna |
| Przycisk „+ Dodaj sekcję” | Mały przycisk `--zarys` | Rozszerzenie struktury o nową sekcję | Mały przycisk tekstowy | 1 | Widoczny w stopce konspektu jako akcja podstawowa | domyślny | Dodaje pozycję do konspektu i umożliwia nadanie nazwy | Stopka panelu konspektu |
| Strefa upuszczenia ustalenia | Wskazany obszar w treści sekcji | Cel przeciągnięcia karty z Findings Panel | Subtelnie obramowany blok, widoczny przy przeciąganiu | 1 (ujawniana samoczynnie) | Pojawia się wyłącznie w trakcie przeciągania karty ustalenia | spoczynek (niewidoczna) · podświetlona (podczas przeciągania) | Upuszczenie wstawia treść ustalenia z automatycznym przypisem CSL | Wewnątrz obszaru treści sekcji |
| Element zbiorczy „Szablon ▼” | Menu progresywne wzorców struktury | Wybór wzorca struktury raportu | Mały element zwinięty | 2 | Kliknięcie elementu `Szablon ▼`; po wyborze lista zwija się samoczynnie | domyślny (pusty konspekt) · wybrany szablon | Wypełnia konspekt predefiniowaną strukturą sekcji | Nagłówek okna |
| Menu raportu `⋮` | Menu kebab nagłówka | Grupowanie ustawień dokumentu: tryb recenzji, styl CSL, zakres bibliografii, ustawienia numeracji przypisów | Mała ikona w nagłówku | 3 | Kliknięcie ikony `⋮` w nagłówku okna | zwinięte · rozwinięte | Rozwija komplet ustawień dokumentu | Nagłówek okna |
| Tryb recenzji | Przełącznik z wątkami komentarzy | Komentowanie fragmentów raportu przed finalizacją | Mały przełącznik + karty komentarzy | 3 | Pozycja w menu `⋮` nagłówka | wyłączony · włączony · wątek rozwiązany | Zaznaczenie fragmentu tworzy wątek komentarza zakotwiczony w treści | Nagłówek okna |
| Pasek operacji kontekstowych | Menu kontekstowe zaznaczenia | Korekta, streszczenie, zmiana stylu, rozwinięcie zaznaczonego fragmentu | Mały pasek przy obszarze treści | 3 | Zaznaczenie fragmentu treści; po użyciu pasek znika | ukryty (brak zaznaczenia) · widoczny | Realizuje operację na zaznaczeniu, analogicznie do Tools Panel modułu Studio | Przy obszarze treści sekcji |
| Element zbiorczy „Operacje raportu ▼” | Grupowanie logiczne akcji dokumentu | Generowanie streszczenia zarządczego, bibliografia, menedżer przypisów, wersje i różnice, wstawki (macierz, oś czasu, wykres, tabela dowodów) | Średni przycisk zbiorczy | 3 | Kliknięcie elementu `Operacje raportu ▼`; po wyborze akcji element zwija się samoczynnie | domyślny · rozwinięty · ładowanie operacji | Rozwija listę operacji i uruchamia wybraną | Stopka okna |
| Menedżer przypisów i bibliografii | Panel z listą przypisów i pozycji literatury | Numeracja, przypisy skrócone, bibliografia końcowa | Średni panel wywoływany | 3 | Pozycja w elemencie `Operacje raportu ▼`; po zamknięciu panel znika z przestrzeni roboczej | domyślny · styl zmieniony · renumeracja | Zmiana stylu CSL przelicza wszystkie cytowania i bibliografię | Stopka okna |
| Panel wersji i różnic | Lista wersji raportu z porównaniem | Przegląd kolejnych kompletacji dokumentu | Średni panel wywoływany | 3 | Pozycja w elemencie `Operacje raportu ▼` albo kliknięcie znacznika wersji w nagłówku | jedna wersja · wiele wersji · widok różnic | Wybór dwóch wersji pokazuje zestawienie różnicowe treści | Stopka okna |
| Wstawka macierzy, osi czasu, wykresu | Blok osadzony w treści sekcji | Prezentacja danych z ustaleń w postaci zestawienia lub grafiki | Średni blok z paskiem edycji | 3 | Pozycja w elemencie `Operacje raportu ▼` albo polecenie w Chat Window | pusty · z danymi · przeliczany | Kliknięcie otwiera edycję kryteriów zestawienia i wyboru danych | Obszar treści sekcji |
| Przycisk „Generuj streszczenie zarządcze” | Pozycja w elemencie zbiorczym operacji | Złożenie streszczenia z kluczowych ustaleń | Pozycja menu, średnia | 3 | Element `Operacje raportu ▼` albo polecenie języka naturalnego w Chat Window | domyślny · ładowanie | Wstawia lub aktualizuje sekcję streszczenia na początku dokumentu | Operacje raportu |
| Przycisk „→ Export Panel” | Przycisk `--zarys` | Przejście do eksportu raportu | Średni przycisk | 1 | Widoczny w stopce okna jako akcja podstawowa | domyślny | Otwiera Export Panel jako rozszerzenie boczne z aktualnym dokumentem | Stopka okna |
| Edycja stylu CSL własnego i szablonu raportu | Definicja własnego stylu cytowania i własnego wzorca struktury dokumentu | Dostosowanie raportu do wymagań wydawcy lub odbiorcy instytucjonalnego | Bez reprezentacji graficznej w stanie spoczynku | 4 | Okno konfiguracji (rozdz. 8), polecenie języka naturalnego w Chat Window, tryb administracyjny | niedostępny dla użytkownika podstawowego · dostępny wg roli | Udostępnia nowy styl i szablon w elementach `Szablon ▼` i menu `⋮` | Poza interfejsem okna |

---

### 3.9. Export Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne, wąska; dostępna zawsze (skompletowanie raportu to sugestia, nie warunek odblokowania) |
| Izolacja domyślna | Operuje na dokumencie bieżącej sesji z Report Builder |

**Zawartość i pełny arsenał funkcji.**

- Wybór formatu eksportu: PDF, DOCX, Markdown, HTML, PPTX (struktura slajdów generowana z konspektu), XLSX (tabele danych, macierz porównawcza, tabela dowodów, macierz kodowania), LaTeX z plikiem BibTeX.
- Wybór miejsca docelowego: pobranie lokalne, wysłanie do Library, przekazanie do modułu Studio do dalszej redakcji, przekazanie do Roundtable do konfrontacji wniosków.
- Personalizacja eksportu: strona tytułowa, spis treści, bibliografia, stopka z datą i wersją, diagram PRISMA, tabela dowodów.
- Podgląd dokumentu przed eksportem w formacie docelowym (mechanizm współdzielony z Preview Window modułu Studio i File Preview modułu Library).
- Zapis szablonu eksportu (kombinacja formatu i ustawień) do ponownego użycia.
- Historia eksportów bieżącej sesji — lista wygenerowanych plików z ponownym pobraniem.
- Udostępnienie linku do raportu, gdy funkcja udostępniania jest skonfigurowana.

**Makieta tekstowa.**

```
 ═════════════════════════════════════════════════════════════════════════
  Chat Window   │ Report Builder          │ Export Panel
  Użytkownik ↔  │ (aktywne okno wiodące)  │ ────────────────────────────
  Wykonawca     │                         │ [ Format: PDF ▼ ]
                │ konspekt · treść sekcji │ [ Cel: Pobierz ▼ ]
  [⟳ Pętla ▼]   │                         │ [ Zawartość dokumentu ▼ ]
                │                         │ ────────────────────────────
                │                         │ [ Historia eksportów (2) ▼ ]
                │                         │ ────────────────────────────
                │                         │ [ Eksportuj teraz ]
 ═════════════════════════════════════════════════════════════════════════
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Przycisk „Eksportuj teraz” | Przycisk `--zloty`, pełna szerokość | Wygenerowanie i dostarczenie pliku | Duży przycisk CTA | 1 (w otwartym panelu) | Widoczny po otwarciu panelu jako akcja podstawowa | domyślny · ładowanie · zakończono (toast potwierdzenia) | Generuje plik zgodnie z ustawieniami i realizuje wybrany cel | Stopka panelu |
| Element zbiorczy „Format ▼” | Menu progresywne formatów wyjściowych | Wybór formatu pliku: PDF, DOCX, Markdown, HTML, PPTX, XLSX, LaTeX z BibTeX | Mały element zwinięty | 2 | Kliknięcie elementu `Format ▼`; po wyborze lista zwija się samoczynnie | domyślny (PDF) · rozwinięty | Zmienia zestaw ustawień personalizacji poniżej | Nagłówek panelu |
| Element zbiorczy „Cel ▼” | Menu progresywne miejsc docelowych | Wybór miejsca docelowego: Pobierz, Library, Studio, Roundtable | Mały element zwinięty | 2 | Kliknięcie elementu `Cel ▼`; po wyborze element zwija się samoczynnie | Pobierz (domyślny) · Library · Studio · Roundtable | Determinuje działanie przycisku „Eksportuj teraz” | Nagłówek panelu |
| Element zbiorczy „Zawartość dokumentu ▼” | Grupowanie logiczne pól wyboru zawartości | Włączenie i wyłączenie strony tytułowej, spisu treści, bibliografii, stopki z datą i wersją, diagramu PRISMA, tabeli dowodów | Średni element zwinięty | 3 | Kliknięcie elementu `Zawartość dokumentu ▼` — rozwija panel popover z kompletem pól wyboru | zwinięty z podsumowaniem wyboru · rozwinięty | Zmienia skład generowanego pliku | Ciało panelu |
| Podgląd przed eksportem | Wyświetlenie dokumentu w formacie docelowym przed generowaniem | Kontrola wyniku przed wygenerowaniem pliku | Pozycja menu wywołująca podgląd | 3 | Pozycja w elemencie `Zawartość dokumentu ▼` albo skrót klawiszowy | domyślny · ładowanie podglądu | Otwiera podgląd analogiczny do Preview Window; po zamknięciu znika z przestrzeni roboczej | Ciało panelu |
| Element zbiorczy „Historia eksportów ▼” | Zwinięta lista pozycji ze znacznikiem czasu i licznikiem | Wygenerowane pliki bieżącej sesji | Mały element zwinięty | 2 | Kliknięcie elementu `Historia eksportów (n) ▼`; po zamknięciu lista zwija się do licznika | pusty stan · z pozycjami | Kliknięcie „Pobierz ponownie” udostępnia wcześniej wygenerowany plik | Stopka panelu |
| Selektor szablonu eksportu | Rozwijana lista zapisanych kombinacji formatu i ustawień | Ponowne użycie zapisanych ustawień eksportu | Mały element zwinięty | 3 | Pozycja w elemencie `Zawartość dokumentu ▼` albo menu kontekstowe nagłówka panelu | brak szablonów · szablon wybrany | Wypełnia format i zawartość zapisanymi wartościami | Ciało panelu |
| Udostępnienie linku do raportu | Wygenerowanie odnośnika do raportu poza modułem | Przekazanie raportu odbiorcy bez pobierania pliku | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, okno konfiguracji (rozdz. 8), tryb administracyjny; dostępne, gdy funkcja udostępniania jest skonfigurowana | niedostępny dla użytkownika podstawowego · dostępny wg roli | Generuje odnośnik do raportu i odnotowuje go w historii eksportów | Poza interfejsem panelu |

---

## 4. Przepływy pracy

### 4.1. Przepływ podstawowy — od pytania badawczego do raportu

```
 Discovery Panel     Sources Manager      Reading View        Findings Panel
 ───────────────     ───────────────      ────────────        ──────────────
 1. zapytanie
    i wyniki  ─────► 2. dodanie źródła
                        z metadanymi
                     3. ocena wiarygodności,
                        katalogowanie
                                    ─────► 4. lektura,
                                              podświetlenia,
                                              wypisy
                                                       ─────► 5. odnotowanie
                                                                 ustalenia,
                                                                 kodowanie
                                                                     │
                                                                     ▼
                                          Report Builder ◄───────────┘
                                          6. przeciągnięcie ustalenia
                                             do sekcji konspektu,
                                             przypis CSL
                                                     │
                                                     ▼
                                          Export Panel
                                          7. wybór formatu i eksport
```

Przebieg pętli wykonawczej towarzyszącej temu przepływowi widoczny jest w Execution Loop Window: Koordynator dekomponuje zlecenie badawcze na zadania wyszukiwania, pozyskania, ekstrakcji, weryfikacji i syntezy, a Wykonawca raportuje wynik każdego zadania wraz z kontrolą jakości.

### 4.2. Przepływ rozszerzony — badanie zasilane z modułu Browser

```
Browser › Browser Window            Research › Sources Manager
────────────────────────           ───────────────────────────
 wspólny podgląd strony  ─────────►  źródło dodane automatycznie
 (Użytkownik + Wykonawca)            wraz z metadanymi (adres,
                                      data pozyskania)
                                             │
Browser › Sources Panel  ─────────►  scalenie z listą źródeł badania
 (lista źródeł sesji)                 (powiązanie Browser ───► Research,
                                       konfigurowalne)
                                             │
                                             ▼
                                   Reading View — lektura i ekstrakcja
                                             │
                                             ▼
                                   Findings Panel — dalsza analiza
```

### 4.3. Przepływ rozszerzony — wykrycie i rozstrzygnięcie sprzeczności

```
Sources Manager ──► dwa źródła o rozbieżnych danych
        │
        ▼
Findings Panel ──► automatyczne wykrycie sprzeczności (alert)
        │
        ▼
Chat Window ──► polecenie: „porównaj metodologię źródeł 7 i 9”
        │
        ▼
Execution Loop Window ──► zadanie weryfikacji: pobranie fragmentów
        │                  metodologicznych, porównanie, kontrola jakości
        ▼
Findings Panel ──► ustalenie rozstrzygające, oznaczenie sprzeczności
                     jako „rozstrzygnięta” z uzasadnieniem i prowenancją
        │
        ▼
Report Builder ──► wstawienie rozstrzygniętego ustalenia z przypisem
                     do obu źródeł
```

### 4.4. Przepływ rozszerzony — przegląd literatury według protokołu PRISMA

```
Discovery Panel ──► zapytania w bazach naukowych, zidentyfikowane pozycje
        │
        ▼
Sources Manager ──► deduplikacja, przesianie tytułów i abstraktów,
        │             odrzucenia z uzasadnieniem
        ▼
Reading View ──► lektura pełnych tekstów, ekstrakcja danych do tabeli dowodów
        │
        ▼
Findings Panel ──► kodowanie jakościowe, synteza ustaleń
        │
        ▼
Report Builder ──► szablon „przegląd literatury”: tabela dowodów,
        │             synteza narracyjna, diagram PRISMA
        ▼
Export Panel ──► PDF, DOCX lub LaTeX z bibliografią w wybranym stylu CSL
```

### 4.5. Przepływ rozszerzony — dalsza redakcja i prezentacja wielomodelowa

```
Research › Report Builder                    Studio › Studio Editor
──────────────────────────                   ──────────────────────
 raport skompletowany     ─────────────────►  dalsza redakcja językowa
                                               (powiązanie konfigurowalne)

Research › Report Builder                    Roundtable › Model Panels
──────────────────────────                   ─────────────────────────
 wnioski robocze przed     ─────────────────►  wielomodelowa konfrontacja
 ostatecznym raportem                          wniosków przed finalizacją
                                               (decyzja Użytkownika)
```

---

## 5. Stany, dane i powiązania

### 5.1. Model stanów badania

```
 ┌───────────────┐   dodanie      ┌───────────────┐   odnotowanie    ┌───────────────┐
 │ Zakres         │   pierwszego   │ Zbieranie      │   pierwszego     │ Analiza        │
 │ zdefiniowany   │───────────────►│ źródeł         │─────────────────►│ i synteza      │
 └───────────────┘   źródła        └───────────────┘   ustalenia      └───────┬───────┘
                                                                                │
                                                                    skompletowanie
                                                                    konspektu
                                                                                ▼
                                                                       ┌───────────────┐
                                                                       │ Redakcja       │
                                                                       │ raportu        │
                                                                       └───────┬───────┘
                                                                                │ eksport
                                                                                ▼
                                                                       ┌───────────────┐
                                                                       │ Raport         │
                                                                       │ wyeksportowany │
                                                                       └───────────────┘
```

Etapy nie są sekwencją wymuszoną — Research Workspace pozwala wracać do wcześniejszego etapu (na przykład dodać nowe źródło po rozpoczęciu redakcji raportu) bez ograniczeń, zgodnie z zasadą pełnej kompozycyjności.

### 5.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Research |
|---|---|
| `sesja`, `karta_sesji` | Nośnik kontekstu badania |
| `wiadomosc` | Historia Chat Window i komunikaty sterujące Execution Loop Window |
| `artefakt`, `wersja_artefaktu` | Źródła jako pliki (F.2), migawki, adnotacje, dane wydobyte z tabel oraz raport i jego kolejne wersje w Report Builder |
| `kolekcja`, `kolekcja_artefakt`, `etykieta_artefaktu` | Grupowanie źródeł tematycznie w Sources Manager, kolekcje inteligentne, kody jakościowe |
| `kanal_modelu` | Kanał wykonawczy używany do analizy, syntezy, embeddingów, transkrypcji i rozpoznawania tabel |
| `powiazanie_komponentu` | Jawne powiązania Browser ───► Research, Research ───► Library, Research ◄──► Studio, Research ───► Roundtable, Research ───► Automations/Agents |
| `profil_izolacji`, `regula_izolacji_kontekstu` | Zakres współdzielenia historii badania między kartami i projektami |

### 5.3. Izolacja i konfigurowalność — punkty właściwe modułowi Research

| Punkt izolacji | Stan wyjściowy (domyślny) | Co można skonfigurować | Gdzie |
|---|---|---|---|
| Historia i pamięć badania | Odrębna per karta sesji | Współdzielenie między kartami prowadzącymi powiązane badania | Okno konfiguracji, poziom „karta sesji” lub „moduł” |
| Źródła z Browser | Brak automatycznego zasilania | Ustanowienie stałego kanału Browser ───► Research | Okno konfiguracji, powiązanie komponentu |
| Dokumenty z Library | Brak automatycznego dostępu | Ustanowienie Library jako źródła materiałów wejściowych | Okno konfiguracji, powiązanie komponentu |
| Zapytania wychodzące do sieci | Wszystkie integracje dozwolone | Ograniczenie zbioru integracji z prawem wyjścia do sieci dla badań wrażliwych | Okno konfiguracji, grupa integracji |
| Izolacja techniczna procesu sesji (8 zakresów) | Żaden zakres domyślnie nie jest aktywny | Włączenie odrębnego dostępu sieciowego dla badania wrażliwego | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

### 5.4. Powiązania z innymi modułami

```
   Browser  ◄──►  RESEARCH  ───►  Library
 (wspólne źródła         (wykorzystanie dokumentów
  i zasilanie badania)     przechowywanych centralnie)
                                    │
                                    ▼
                    Studio  ◄──►  RESEARCH   (redakcja / dalsze badanie)
                                    │
                                    ├──────────────► Roundtable
                                    │        (wielomodelowa prezentacja
                                    │         wyników przed wnioskami)
                                    ▼
                          Automations / Agents
                    (cykliczne monitorowanie tematów
                     i cotygodniowe przeglądy źródeł)
```

| Moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Browser | Zasilanie pracy badawczej zebranymi źródłami | Konfiguracyjne | Sources Panel → Sources Manager |
| Library | Wykorzystanie dokumentów przechowywanych centralnie jako materiału źródłowego i trwały magazyn raportów | Konfiguracyjne | Library Explorer → Sources Manager · Export Panel → Library |
| Studio | Dalsza redakcja wyników pracy badawczej | Konfiguracyjne | Report Builder → Studio Editor |
| Roundtable | Wielomodelowa prezentacja wyników przed wnioskami końcowymi | Konfiguracyjne | Report Builder / Findings Panel → Model Panels |
| Automations / Agents | Cykliczne uruchamianie operacji badawczych: monitorowanie tematów, odświeżanie źródeł, przeglądy okresowe | Konfiguracyjne | Discovery Panel / Execution Loop Window → Automations |
| Assistant | Współdzielenie mechanizmu transkrypcji nagrań | Konfiguracyjne | Sources Manager → mechanizm transkrypcji |

---

## 6. Scenariusze użycia

**Scenariusz 1 — analiza konkurencyjna od podstaw.**

Analityk otwiera moduł Research, definiuje zakres w Research Workspace („analiza konkurencyjna segmentu X”), zleca w Chat Window rozpoznanie tematu. Koordynator rozpisuje zlecenie w Execution Loop Window na zadania wyszukiwania i pozyskania źródeł, Discovery Panel zwraca wyniki webowe i naukowe, a wybrane pozycje trafiają do Sources Manager z pełnymi metadanymi. Analityk czyta kluczowe pozycje w Reading View, zamienia podświetlenia w ustalenia, a po zgromadzeniu 14 źródeł i 22 ustaleń wybiera szablon „Analiza konkurencyjna” w Report Builder, przeciąga kluczowe ustalenia do sekcji i generuje streszczenie zarządcze.

**Scenariusz 2 — badanie zasilane wspólnym przeglądaniem.**

Zespół prowadzi wspólną analizę stron konkurencji w module Browser. Zebrane tam źródła, dzięki wcześniej skonfigurowanemu powiązaniu, trafiają do Sources Manager modułu Research, gdzie są katalogowane, deduplikowane i oceniane pod kątem wiarygodności, a następnie czytane i adnotowane w Reading View.

**Scenariusz 3 — wykrycie i rozstrzygnięcie sprzeczności danych.**

Findings Panel sygnalizuje sprzeczność między dwoma źródłami podającymi różne wartości udziału rynkowego. Użytkownik poleca w Chat Window porównanie metodologii obu źródeł; Execution Loop Window pokazuje przebieg zadania weryfikacji i wynik kontroli jakości. Użytkownik zapisuje rozstrzygające ustalenie i oznacza sprzeczność jako rozstrzygniętą — z pełnym śladem prowenancji.

**Scenariusz 4 — przegląd literatury według protokołu.**

Badaczka prowadzi przegląd systematyczny: w Discovery Panel wykonuje zapytania w Crossref, OpenAlex i PubMed, w Sources Manager deduplikuje i przesiewa pozycje z uzasadnieniem odrzuceń, w Reading View wydobywa dane do tabeli dowodów, w Findings Panel koduje fragmenty, a w Report Builder składa przegląd z diagramem PRISMA i bibliografią w stylu czasopisma docelowego.

**Scenariusz 5 — eksport raportu w wielu formatach dla różnych odbiorców.**

Po skompletowaniu raportu w Report Builder użytkownik generuje w Export Panel wersję PDF ze stroną tytułową i bibliografią dla zarządu, wersję DOCX bez strony tytułowej do dalszej redakcji w module Studio oraz arkusz XLSX z macierzą porównawczą dla zespołu analitycznego.

**Scenariusz 6 — konfrontacja wniosków w Roundtable przed finalizacją.**

Przed napisaniem ostatecznych wniosków analityk przekazuje wstępną syntezę ustaleń z Findings Panel do modułu Roundtable, gdzie kilka modeli ocenia mocne i słabe strony wnioskowania z różnych perspektyw, zanim wersja końcowa trafi z powrotem do Report Builder.

**Scenariusz 7 — badanie prowadzone cyklicznie.**

Zespół konfiguruje monitory tematów i kanały RSS w Discovery Panel oraz powiązanie Research ───► Automations. Nowe publikacje trafiają do skrzynki „nowe źródła”, a cotygodniowy przegląd jest składany przez pętlę wykonawczą i publikowany jako kolejna wersja raportu w Library.

---

## 7. Katalog funkcji i narzędzi

Każda pozycja: **nazwa** · co robi · zależności techniczne (biblioteki Go, formaty, integracje). Rdzeń wykonawczy w Go, magazyn SQLite, kanał WebSocket/JSON, klient Tauri — nazwy bibliotek dobrane pod ten stos.

### 7.1. Wyszukiwanie i odkrywanie źródeł

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Web Search Bridge** | Wyszukiwanie webowe z poziomu modułu; wyniki wpadają wprost na kandydatów do Sources Manager z metadanymi (tytuł, URL, fragment) | Integracje: konfigurowalne API wyszukiwania (Brave Search API, Bing Web Search, SerpAPI, DuckDuckGo), klucze jawne w oknie konfiguracji; klient HTTP `net/http` |
| **Scholarly Search** | Wyszukiwanie w bazach naukowych po zapytaniu i filtrach (rok, dziedzina, dostęp otwarty) | Integracje: Crossref REST, OpenAlex, Semantic Scholar Graph API, arXiv API, PubMed E-utilities, CORE; parser JSON `encoding/json` |
| **DOI / ISBN / PMID Resolver** | Po wklejeniu identyfikatora pobiera pełne metadane pozycji i tworzy referencję | Crossref (`https://api.crossref.org/works/{doi}`), OpenLibrary (ISBN), PubMed; mapowanie na CSL-JSON |
| **Semantic Search (własne źródła)** | Wyszukiwanie po znaczeniu w całym korpusie zebranych źródeł i ustaleń, nie po słowach kluczowych | Embeddingi z `kanal_modelu`; magazyn wektorów: rozszerzenie SQLite `sqlite-vec` / `sqlite-vss` albo in-process HNSW (`hnswlib` binding); indeks aktualizowany przy dodaniu źródła |
| **Full-Text Search (FTS)** | Wyszukiwanie pełnotekstowe w treści źródeł, adnotacji i ustaleń, z podświetleniem trafień | SQLite **FTS5** (tokenizer unicode61 + porter); snippet i highlight z FTS5 |
| **Topic Monitor / Alerts** | Zapisany temat monitorowany cyklicznie; nowe pasujące publikacje trafiają do skrzynki „nowe źródła” | RSS/Atom przez `mmcdole/gofeed`; zapytania Crossref/OpenAlex z filtrem daty; wyzwalanie przez Automations (konfigurowalne) |
| **RSS / Atom Reader** | Agregacja kanałów źródeł, blogów i feedów czasopism; pozycja jednym kliknięciem staje się źródłem | `mmcdole/gofeed` (RSS 1.0/2.0, Atom, JSON Feed) |
| **Snowballing (cytowania wstecz i wprzód)** | Z pozycji rozwija listę prac cytowanych i cytujących, budując graf literatury | Crossref/OpenAlex relacje `referenced-works` i `cited-by`; graf w SQLite; render w Evidence Map |
| **Query Assistant** | Przekształca pytanie badawcze w zestaw skutecznych zapytań wyszukiwawczych i operatorów (boolean, filtry) | `kanal_modelu`; szablony operatorów per dostawca wyszukiwania |

### 7.2. Pozyskiwanie i wczytywanie źródeł

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Web Clipper** | Zapisuje stronę jako czyste źródło: wyodrębniona treść, autor, data, tytuł, obraz wiodący; usuwa nawigację i reklamy | Ekstrakcja treści: `go-shiori/go-readability` (port Readability); pobranie: `net/http`; render JS przez `chromedp` (Chrome headless) dla stron dynamicznych, sterowany ustawieniem konfiguracyjnym |
| **Page Archiver (migawka)** | Zapisuje pełną, niezmienną migawkę strony (HTML z zasobami albo PDF) na wypadek zniknięcia źródła | `chromedp` (print-to-PDF, pełnostronicowy zrzut); pakowanie WARC/ZIP; integracja Wayback Machine API dla wersji archiwalnej |
| **Import PDF** | Wczytuje dokument PDF, ekstrahuje tekst, strukturę stron, metadane i obrazy; przygotowuje do czytnika i wyszukiwania pełnotekstowego | Ekstrakcja tekstu: `ledongthuc/pdf` albo `go-fitz` (binding MuPDF, wierniejszy układ); render stron do podglądu: `go-fitz`; metadane: `pdfcpu` |
| **Import DOCX / ODT / EPUB / TXT / MD** | Wczytuje dokumenty tekstowe różnych formatów i normalizuje je do wewnętrznej reprezentacji treści | DOCX: `nguyenthenguyen/docx` / `unidoc/unioffice`; EPUB: `taylorskalyo/goreader` albo własny czytnik ZIP+XHTML; MD: `yuin/goldmark` |
| **OCR skanów i obrazów** | Rozpoznaje tekst w zeskanowanych dokumentach PDF i obrazach, czyniąc je przeszukiwalnymi i cytowalnymi | `otiai10/gosseract` (binding **Tesseract**, pakiety językowe pol+eng); przetwarzanie wstępne obrazu `disintegration/imaging` |
| **Import bibliografii zbiorczej** | Wczytuje listę referencji z pliku i tworzy pozycje w bibliotece referencji | BibTeX: `nickng/bibtex`; RIS: własny parser; CSL-JSON: `encoding/json`; EndNote XML |
| **Audio/Video → transkrypt** | Transkrybuje nagranie wywiadu lub wykładu; transkrypt staje się cytowalnym źródłem ze znacznikami czasu | Transkrypcja: `whisper.cpp` (binding CGo) albo konfigurowalne API (kanał modelu obsługujący audio); diaryzacja mówców sterowana ustawieniem konfiguracyjnym; współdzielenie mechanizmu z modułem Assistant |
| **Clipboard / Quote Capture** | Wklejony fragment tekstu staje się cytatem-źródłem z zachowaniem odnośnika do miejsca pochodzenia | SQLite (`artefakt` typu cytat); powiązanie z URL/DOI, gdy obecne |
| **Batch URL Import** | Wklejenie listy adresów naraz; każdy przetwarzany przez Web Clipper w tle | Kolejka zadań w rdzeniu Go (pula goroutines), status per pozycja przez WebSocket, przebieg widoczny w Execution Loop Window |

### 7.3. Zarządzanie źródłami i referencjami

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Reference Library** | Trwała biblioteka referencji z pełnymi metadanymi (autor, rok, wydawca, DOI, typ pozycji) — rdzeń menedżera bibliografii | Model danych: `artefakt` + `etykieta_artefaktu` + `kolekcja`; schemat pól zgodny z typami pozycji CSL-JSON |
| **Credibility Scoring** | Ocena wiarygodności (pierwotne/wtórne, punktacja) ręczna lub sugerowana według kryteriów recenzji, cytowalności i domeny | `kanal_modelu` do sugestii; dane cytowalności z OpenAlex/Crossref; reguły punktacji konfigurowalne |
| **Deduplikacja** | Wykrywa duplikaty i niemal-duplikaty (ten sam DOI/URL, zbliżona treść) i proponuje scalenie | Hash URL/DOI; podobieństwo treści przez embeddingi (próg konfigurowalny); SimHash `mfonda/simhash` jako szybki pre-filtr |
| **Tagi i kolekcje** | Katalogowanie wielowymiarowe: tagi tematyczne, kolekcje swobodne i inteligentne (regułowe), przypisanie do etapu i pytania badawczego | `kolekcja`/`kolekcja_artefakt` (współdzielone z Library); reguły kolekcji inteligentnych ewaluowane w Go |
| **Bibliography Export** | Eksportuje pełną bibliografię lub wybór do BibTeX, RIS, CSL-JSON i do sformatowanej listy w stylu CSL | Serializatory: `nickng/bibtex`, CSL-JSON; formatowanie przez procesor CSL (rozdz. 7.6) |
| **Attachment Manager** | Wiąże ze źródłem załączniki: PDF pełnotekstowy, migawkę, notatki; kontroluje kompletność | `artefakt` powiązany relacją; przekazanie do Library przez `powiazanie_komponentu` |
| **Watched Folder Import** | Wskazany folder lokalny lub katalog Library jest monitorowany — nowe dokumenty PDF wchodzą jako źródła z automatyczną ekstrakcją metadanych | `fsnotify/fsnotify`; ekstrakcja metadanych z PDF (DOI z tekstu → Crossref) |

### 7.4. Lektura, adnotacje i ekstrakcja

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Reading View / czytnik** | Renderuje źródło (PDF, HTML, DOCX, EPUB) w trybie lektury z podglądem stron i miniaturami | Render: `go-fitz` (PDF → obraz strony), `goldmark` (MD → HTML), sanitizacja HTML `bluemonday`; podgląd współdzielony z Library File Preview |
| **Highlights & Annotations** | Podświetlanie fragmentów w kolorach, notatki na marginesie, zakładki; każde podświetlenie staje się ustaleniem | Kotwice pozycji (strona + przesunięcie, selektor CSS dla HTML, znacznik czasu dla transkryptu); zapis w SQLite; eksport adnotacji do Findings Panel |
| **Excerpt Extractor** | Zbiera wszystkie podświetlenia i wypisy z jednego lub wielu źródeł w jedną listę wypisów gotową do syntezy | Zapytanie po kotwicach; grupowanie po źródle i tagu |
| **Table & Figure Extractor** | Wykrywa i wyodrębnia tabele oraz dane liczbowe z dokumentów PDF do postaci ustrukturyzowanej (do wykresu lub arkusza XLSX) | Detekcja tabel: `go-fitz` (pozycje słów) z heurystyką siatki; złożone tabele przez model wizyjny z `kanal_modelu`; wynik zapisywany jako `artefakt` typu dane |
| **Claim / Entity Extraction** | Wydobywa ze źródła kluczowe twierdzenia, dane liczbowe i podmioty (osoby, organizacje, daty) i proponuje je jako ustalenia | `kanal_modelu`; wynik mapowany na wpisy Findings Panel z odnośnikiem do fragmentu |
| **Per-Source Summary** | Streszczenie pojedynczego źródła (abstrakt roboczy, tezy, metodologia, wnioski) generowane na żądanie | `kanal_modelu`; zapis jako notatka źródła |
| **Cross-Document Chat (RAG)** | Rozmowa oparta na treści wybranych źródeł, z odpowiedziami zakotwiczonymi w cytatach i numerach stron | Retrieval z indeksu wektorowego (rozdz. 7.1) i FTS; grounding promptu; cytowania wstawiane jako odnośniki `[n]` |
| **Reading Progress** | Śledzi, które źródła przeczytano, przejrzano i pominięto oraz pokrycie pytań badawczych źródłami | Pola stanu w `artefakt`; wskaźnik w Research Workspace |

### 7.5. Analiza i synteza

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Findings Ledger** | Rejestr ustaleń cząstkowych z klasyfikacją (fakt/hipoteza/opinia/dana/cytat), wagą i powiązaniem do źródeł | `artefakt` i wpisy w SQLite; relacja wiele-do-wielu ustalenie ↔ źródło |
| **Qualitative Coding** | Kodowanie fragmentów kodami tematycznymi, budowa książki kodów, zliczanie wystąpień — analiza jakościowa | Kody jako etykiety; macierz kod × źródło; eksport do XLSX; sugestie kodów przez `kanal_modelu` |
| **Contradiction Detector** | Wykrywa sprzeczności między ustaleniami i źródłami oraz sygnalizuje je do rozstrzygnięcia | `kanal_modelu` (porównanie par twierdzeń); heurystyka rozbieżności liczbowych; alert w Findings Panel |
| **Evidence Map / graf** | Interaktywny graf: źródła, ustalenia, twierdzenia i ich relacje (poparcie, sprzeczność, cytowanie) | Graf w SQLite; render w kliencie (układ siła-ukierunkowany); dane z relacji ustaleń i snowballingu |
| **Comparative Matrix** | Zestawienie wielu źródeł lub wariantów wg wspólnych kryteriów w tabeli porównawczej (konkurenci × cechy) | Model danych ustaleń → pivot; generowanie kryteriów przez `kanal_modelu`; eksport XLSX |
| **Thematic Grouping** | Grupowanie ustaleń w wątki tematyczne przez klastrowanie | Embeddingi ustaleń i klastrowanie (k-means / HDBSCAN w Go albo przez model); nazwy wątków z `kanal_modelu` |
| **Research Gap Finder** | Wskazuje obszary tematu niepokryte źródłami oraz pytania bez odpowiedzi | `kanal_modelu` porównujący zakres pytań z pokryciem źródeł; wynik jako karta sugestii w Research Workspace |
| **Simple Data Charts** | Poglądowe wykresy (słupki, linie, udziały) z danych liczbowych zebranych w ustaleniach | Render SVG w kliencie; dane z `artefakt` typu dane; granica: głęboka analiza należy do modułu analizy danych |
| **Timeline Builder** | Chronologiczna oś wydarzeń i publikacji zbudowana z dat w ustaleniach i źródłach | Render w kliencie; dane z pól dat; eksport do raportu |

### 7.6. Cytowania i bibliografia

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Citation Generator (CSL)** | Generuje cytat w tekście i pozycję bibliograficzną w wybranym stylu (APA, MLA, Chicago, IEEE, Vancouver, styl prawniczy, style czasopism) | **Citation Style Language**: repozytorium stylów CSL (około 2500 plików `.csl`) i pliki locale; procesor CSL — most do `citeproc` (port `citeproc-go`) albo `citeproc-js` uruchamiany w izolowanym środowisku JS wywoływanym z rdzenia; dane wejściowe CSL-JSON |
| **In-text Citation Insert** | Wstawia odnośnik `[n]` albo (Autor, rok) w Report Builder z automatycznym wpisem do bibliografii końcowej | Mapowanie `artefakt` → numer/klucz; synchronizacja z Bibliography Builder |
| **Footnote / Endnote Manager** | Zarządza przypisami dolnymi i końcowymi, numeracją oraz przypisami skróconymi (ibid., op. cit.) | Logika CSL dla przypisów skróconych; renumeracja przy edycji |
| **Bibliography Builder** | Składa listę literatury końcowej z wykorzystanych źródeł w wybranym stylu, sortowaną i odfiltrowaną | Procesor CSL; źródła z Reference Library; wybór zakresu: wszystkie albo tylko cytowane |
| **Citation Health Check** | Sprawdza kompletność metadanych do cytowania, wykrywa braki (brak roku, autora, DOI), źródła niecytowane i cytowania bez wpisu | Walidacja pól CSL-JSON; raport braków; uzupełnienie z Crossref |
| **Retraction / Integrity Check** | Sygnalizuje, czy cytowana praca została wycofana lub skorygowana | Integracja Retraction Watch (Crossref Labs API); flaga przy pozycji |
| **Style Editor** | Podgląd i personalizacja stylu cytowania (własny wariant CSL) | Edycja plików `.csl` (XML); walidacja schematu CSL |

### 7.7. Opracowanie i eksport

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Report Builder (konspekt → dokument)** | Buduje raport z ustaleń: konspekt, sekcje, wstawianie ustaleń, szablony (analiza konkurencyjna, SWOT, benchmark, nota badawcza, przegląd literatury, streszczenie zarządcze) | Edytor treści w kliencie; szablony w SQLite; operacje kontekstowe z `kanal_modelu` |
| **Literature Review Assembler** | Składa przegląd literatury: tabela dowodów, synteza narracyjna, luki badawcze, zgodnie z protokołem PRISMA dla przeglądów systematycznych | Szablon PRISMA; diagram przepływu (zidentyfikowane → przesiane → włączone); render SVG; dane z etapów przeglądu |
| **Executive Summary Generator** | Składa streszczenie zarządcze z kluczowych ustaleń | `kanal_modelu`; wejście: ustalenia oznaczone jako kluczowe |
| **Export: PDF** | Generuje dokument PDF ze stroną tytułową, spisem treści, bibliografią i stopką z wersją | Silnik: HTML→PDF przez `chromedp` (Chrome print) dla wiernego składu albo `maroux/go-wkhtmltopdf` / `unidoc/unipdf` dla generacji natywnej; szablony HTML/CSS druku |
| **Export: DOCX** | Generuje edytowalny dokument DOCX do dalszej pracy w module Studio | `unidoc/unioffice` (tworzenie DOCX ze stylami, nagłówkami, tabelami i przypisami) |
| **Export: Markdown / HTML** | Czysty Markdown albo samodzielny plik HTML | `goldmark` (render MD → HTML); serializacja Markdown z wewnętrznej reprezentacji |
| **Export: PPTX** | Slajdy generowane ze struktury konspektu (tytuł sekcji → slajd, kluczowe ustalenia → punkty) | `unidoc/unioffice` (PPTX) albo generacja OOXML; współpraca z modułem Design w zakresie motywu |
| **Export: XLSX** | Arkusz z tabelami danych, macierzą porównawczą, tabelą dowodów i macierzą kodowania | `qax-os/excelize` |
| **Export: LaTeX + BibTeX** | Źródło LaTeX z osadzonym `\cite` i plikiem `.bib` dla środowisk akademickich | Szablon LaTeX; `nickng/bibtex`; mapowanie CSL-JSON → BibTeX |
| **Preview przed eksportem** | Wierny podgląd dokumentu w formacie docelowym przed generacją | Mechanizm współdzielony z Library File Preview i Studio Preview |
| **Cel eksportu** | Pobranie lokalne, wysłanie do Library, przekazanie do Studio (redakcja) albo do Roundtable (konfrontacja) | `powiazanie_komponentu`; kanały konfigurowalne |

### 7.8. Współpraca, orkiestracja i higiena badania

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Research Protocol / Stages** | Definiowalne etapy badania (rozpoznanie → zbieranie → analiza → synteza → raport) z checklistą i statusem | Model etapów w SQLite; Kanban i oś czasu w Research Workspace |
| **Provenance & Audit Trail** | Pełny ślad pochodzenia każdego ustalenia: z jakiego źródła, z jakiego fragmentu, przez kogo (Użytkownik albo Wykonawca) i kiedy | `wersja_artefaktu.autor`; log operacji pętli wykonawczej; spójne z blokiem prowenancji platformy (`prowenancja*.ts`) |
| **Fact-Check Assistant** | Weryfikuje pojedyncze twierdzenie względem dostępnych źródeł i wskazuje poparcie, brak poparcia albo sprzeczność | RAG po korpusie i wyszukiwanie webowe; `kanal_modelu`; wynik jako adnotacja ustalenia |
| **Duplicate Findings Merge** | Scala powtórzone ustalenia z różnych źródeł w jedno z wieloma odnośnikami | Podobieństwo embeddingów; potwierdzenie użytkownika |
| **Session Snapshots / wersje raportu** | Kolejne kompletacje raportu jako wersje, z porównaniem różnic | `wersja_artefaktu`; diff tekstu: `sergi/go-diff` (Diff-Match-Patch); mechanizm współdzielony z Library Versioning |
| **Review Mode / komentarze** | Komentarze do fragmentów raportu przed finalizacją, tryb recenzji | Kotwice fragmentów; wątki komentarzy w SQLite |
| **Handoff do Automations/Agents** | Udostępnia operacje (monitoruj temat, odśwież źródła, wygeneruj cotygodniowy przegląd) do orkiestracji cyklicznej | `powiazanie_komponentu` → Automations/Agents; operacje jako wyzwalane akcje |

### 7.9. Zależności techniczne o charakterze przekrojowym

| Zależność | Charakter | Rozstrzygnięcie w module |
|---|---|---|
| Procesor CSL | Brak w pełni natywnej biblioteki Go | Moduł korzysta z `citeproc-js` osadzonego w izolowanym środowisku JS wywoływanym z rdzenia Go; interfejs wywołania jest wspólny z portem `citeproc-go`, dane wejściowe i wyjściowe pozostają w CSL-JSON |
| Render PDF wysokiej jakości (`go-fitz`) | Zależność natywna MuPDF przez CGo | Konfiguracja w buildzie wieloplatformowym Tauri/Go obejmuje biblioteki MuPDF dla wszystkich obsługiwanych platform |
| Magazyn wektorowy | Wybór między rozszerzeniem SQLite a osobnym serwerem | Moduł używa rozszerzenia SQLite `sqlite-vec` dla spójności z jednym plikiem danych platformy, bez osobnego serwera wektorowego |
| Headless Chrome (`chromedp`) | Jedna zależność obsługująca wiele funkcji | Ten sam proces obsługuje Web Clipper, Page Archiver i eksport HTML→PDF |
| Tesseract i `whisper.cpp` | Zależności natywne przez CGo | Pakiety językowe OCR (pol, eng i kolejne) oraz model transkrypcji dostarczane z buildem, sterowane ustawieniami konfiguracyjnymi |
| SQLite FTS5 | Wbudowany mechanizm magazynu | Indeks pełnotekstowy budowany przy dodaniu i aktualizacji źródła, adnotacji i ustalenia |

---

## 8. Punkty sterowania z okna konfiguracji

Wszystkie poniższe punkty Operator personalizuje z okna konfiguracji; wartości domyślne oznaczają działanie bez blokad, klucze i integracje pozostają jawne.

| Grupa | Punkt sterowania | Zakres personalizacji | Domyślnie |
|---|---|---|---|
| **Integracje wyszukiwania** | Dostawca wyszukiwania webowego | Wybór i klucz API (Brave, Bing, SerpAPI, DuckDuckGo), wielu dostawców naraz | Brak dostawcy → tryb wklejania wyników |
| | Bazy naukowe | Włączenie Crossref, OpenAlex, Semantic Scholar, arXiv, PubMed, CORE; klucze gdy wymagane; adres kontaktowy do puli grzecznościowej Crossref | Otwarte API włączone bez klucza |
| | Archiwizacja i Wayback | Włączenie archiwum migawek i integracji Wayback | Migawka lokalna włączona, Wayback wyłączony |
| **Modele i kanały wykonawcze** | `kanal_modelu` per operacja | Osobny kanał do syntezy, ekstrakcji, embeddingów, transkrypcji i rozpoznawania tabel | Kanał domyślny sesji |
| | Automatyzm sugestii | Które sugestie generują się samoczynnie (luki, tagowanie, sprzeczności, kody), a które na żądanie | Wszystkie na żądanie |
| | Grounding cytowań | Wymóg zakotwiczenia odpowiedzi RAG w cytatach (rygor cytowania) | Włączony |
| **Pętla wykonawcza** | Widoczność Execution Loop Window | Otwarcie okna przy starcie zlecenia albo na żądanie | Otwarcie przy starcie zlecenia |
| | Próg ponowienia zadania | Liczba ponowień zadania po negatywnym wyniku kontroli jakości | 2 |
| | Zakres kontroli jakości | Które kontrole obowiązują: kompletność metadanych, zakotwiczenie ustalenia, pokrycie pytań, kompletność cytowań | Wszystkie włączone |
| **Cytowania** | Domyślny styl CSL | Wybór z repozytorium stylów; własne warianty `.csl` | APA |
| | Wykrywanie wycofań | Włączenie Retraction Watch | Włączone |
| **Źródła** | Watched Folders | Ścieżki folderów lokalnych i katalogów Library do automatycznego importu | Brak |
| | Kanały RSS i monitory tematów | Lista feedów, częstotliwość odświeżania | Brak |
| | Próg deduplikacji | Czułość wykrywania niemal-duplikatów (podobieństwo) | Średnia |
| | OCR | Języki Tesseract (pol, eng i kolejne), automatyczne OCR dla skanów | pol+eng, automatyczne włączone |
| **Powiązania modułów** | Browser ◄──► Research | Stały kanał zasilania źródłami z Browser | Wyłączone (włączenie jawne) |
| | Research ───► Library | Przekazywanie raportów i źródeł do magazynu trwałego | Wyłączone |
| | Research ◄──► Studio | Przekazanie raportu do redakcji | Wyłączone |
| | Research ───► Roundtable | Przekazanie syntezy do konfrontacji | Wyłączone |
| | Research ───► Automations/Agents | Udostępnienie operacji cyklicznych | Wyłączone |
| **Eksport** | Domyślny format i zawartość | Format, strona tytułowa, spis treści, bibliografia, stopka | PDF, wszystko włączone |
| | Silnik PDF | `chromedp` (wierność składu) albo generacja natywna (lekkość) | `chromedp` |
| | Szablony raportu i eksportu | Zestaw dostępnych szablonów, szablony własne | Pełny zestaw |
| **Izolacja i prywatność** | Historia badania | Odrębna per karta, współdzielona między kartami albo na poziomie modułu | Per karta sesji |
| | 8 zakresów izolacji technicznej | Odrębny dostęp sieciowy i zasobowy dla badania wrażliwego | Żaden aktywny (pełny dostęp) |
| | Retencja migawek i wersji | Okres przechowywania | Bez limitu |
| | Zapytania wychodzące do sieci | Które integracje mają prawo wychodzić do sieci (tryb offline dla źródeł wrażliwych) | Wszystkie dozwolone |

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + Shift + K` | Szybkie dodanie źródła | Sources Manager |
| `Ctrl/Cmd + N` | Nowe ustalenie | Findings Panel |
| `Ctrl/Cmd + F` | Wyszukiwanie pełnotekstowe w korpusie badania | Discovery Panel |
| `Ctrl/Cmd + H` | Podświetlenie zaznaczonego fragmentu | Reading View |
| `Ctrl/Cmd + L` | Otwarcie i zamknięcie Execution Loop Window | Wszystkie okna modułu |
| `@` | Odwołanie do źródła, ustalenia lub kolekcji w poleceniu | Chat Window |
| `Ctrl/Cmd + E` | Otwarcie Export Panel | Report Builder |
| `Ctrl/Cmd + K` | Paleta poleceń | Wszystkie okna modułu |
| Ikona `dokument` | Źródło typu plik lub dokument | Sources Manager |
| Ikona `globus` | Źródło typu strona internetowa | Sources Manager |
| Ikona `waga` | Ocena wiarygodności źródła | Sources Manager |
| Ikona `zakreslacz` | Podświetlenie fragmentu | Reading View |
| Ikona `petla` | Stan pętli wykonawczej | Execution Loop Window |
| Ikona `ostrzezenie` | Sprzeczność wykryta | Findings Panel |
| Ikona `gwiazdka` | Ustalenie kluczowe | Findings Panel |
| Ikona `cytat` | Generowanie cytatu w stylu CSL | Sources Manager · Report Builder |
| Ikona `pobierz` | Eksport raportu | Export Panel |
| Ikona `link-zewnetrzny` | Otwarcie źródła w Browser | Sources Manager |
| Ikona `archiwum` | Historia eksportów | Export Panel |

*Koniec dokumentu. Moduł Research — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*

*© 2026 Danaco Holding Group Sp. z o.o. — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
