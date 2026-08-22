# Danaco Console — Środowisko CodeStudio

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
| **Tytuł** | Środowisko CodeStudio |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Pełnozakresowy materiał źródłowy do projektowania i budowy interfejsu środowiska programistycznego CodeStudio: co istnieje, gdzie leży w układzie, w jakiej formie i do czego służy — dla każdego okna, panelu i elementu |
| **Zakres** | przeznaczenie i kontekst środowiska, osiem modułów, przedsionek wejściowy, komplet okien z warstwami widoczności, specyfikacja okien, przepływy pracy, stany i powiązania, punkty sterowania z okna konfiguracji, scenariusze użycia, wykazy normatywne, skróty klawiszowe, kryteria odbioru |
| **Poza zakresem** | wnętrze poszczególnych modułów — [Developer](../moduly/developer.md) · [Terminal](../moduly/terminal.md) · [Diagnostics](../moduly/diagnostics.md) · [Apps](../moduly/apps.md) · [Roundtable](../moduly/roundtable.md) · [Design](../moduly/design.md) · [Agents](../moduly/agents.md) · [Workspace](../moduly/workspace.md); mechanizm magistrali kontekstu — [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) rozdz. 6a |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [Strona główna i nawigacja](../interfejs-uzytkownika/strona-glowna-i-nawigacja.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Specyfikacja agentów](../specyfikacje/specyfikacja-agentow.md) · [Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) · [Rozszerzenia](../architektura/rozszerzenia.md) |
| **Prototypy odniesienia** | `design/05-okna/srodowiska/codestudio.html` · `design/05-okna/srodowiska/codestudio-przedsionek.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `environment`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` · `design/05-okna/srodowiska/` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność; zero blokad w interfejsie (brak nieaktywnych kontrolek, okien nakładkowych-bramek, walidacji blokującej); klucze konfiguracji jawne; izolacja i uprawnienia są ustawieniami konfiguracyjnymi Operatora, nie wymogiem platformy; domyślne zachowanie środowiska to wykonanie polecenia |

Dokument opisuje środowisko CodeStudio w postaci, w jakiej ma powstać w platformie Danaco Console: przeznaczenie środowiska, moduły przez nie udostępniane, komplet okien z przypisaniem do warstw widoczności, specyfikację każdego okna wraz z makietami tekstowymi, przepływy pracy, katalog stanów i powiązań, punkty sterowania dostępne z okna konfiguracji oraz scenariusze użycia. Dwa kanały komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) i Execution Loop Window (Koordynator ↔ Wykonawca) — stanowią elementy pierwszoplanowe środowiska; wszystkie pozostałe okna są wobec nich narzędziami roboczymi.

---

## Spis treści

1. [Wprowadzenie](#wprowadzenie)
2. [Przeznaczenie i kontekst środowiska](#1-przeznaczenie-i-kontekst-środowiska)
   - [1.1 Definicja środowiska](#11-definicja-środowiska)
   - [1.2 Zakres środowiska i rozgraniczenie od modułów](#12-zakres-środowiska-i-rozgraniczenie-od-modułów)
   - [1.3 Położenie w hierarchii platformy](#13-położenie-w-hierarchii-platformy)
   - [1.4 Kanoniczne rozmieszczenie obszaru roboczego](#14-kanoniczne-rozmieszczenie-obszaru-roboczego)
3. [Moduły dostępne w środowisku](#2-moduły-dostępne-w-środowisku)
   - [2.1 Osiem modułów CodeStudio](#21-osiem-modułów-codestudio)
   - [2.2 Boczna nawigacja modułów](#22-boczna-nawigacja-modułów)
   - [2.3 Magistrala kontekstu środowiska](#23-magistrala-kontekstu-środowiska)
   - [2.4 Przedsionek środowiska — widok wejściowy](#24-przedsionek-środowiska--widok-wejściowy)
4. [Komplet okien środowiska](#3-komplet-okien-środowiska)
   - [3.1 Warstwy widoczności w środowisku](#31-warstwy-widoczności-w-środowisku)
   - [3.2 Tabela kompletu okien](#32-tabela-kompletu-okien)
   - [3.3 Listwa stanu środowiska](#33-listwa-stanu-środowiska)
5. [Specyfikacja okien](#4-specyfikacja-okien)
   - [4.1 Makieta główna — pełny układ obszaru roboczego](#41-makieta-główna--pełny-układ-obszaru-roboczego)
   - [4.2 Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca](#42-chat-window--główne-okno-komunikacji-użytkownik--wykonawca)
   - [4.3 Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca](#43-execution-loop-window--okno-pętli-wykonawczej-koordynator--wykonawca)
   - [4.4 Okna robocze modułu Developer](#44-okna-robocze-modułu-developer)
   - [4.5 Karty sesji](#45-karty-sesji)
   - [4.6 Okna pomocnicze powłoki](#46-okna-pomocnicze-powłoki)
   - [4.7 Always On Display i Mobile w środowisku](#47-always-on-display-i-mobile-w-środowisku)
6. [Przepływy pracy](#5-przepływy-pracy)
   - [5.1 Przepływ podstawowy — zlecenie programistyczne przez pętlę wykonawczą](#51-przepływ-podstawowy--zlecenie-programistyczne-przez-pętlę-wykonawczą)
   - [5.2 Przepływ artefaktu magistralą kontekstu](#52-przepływ-artefaktu-magistralą-kontekstu)
   - [5.3 Przepływ pracy z kartą sesji](#53-przepływ-pracy-z-kartą-sesji)
   - [5.4 Przepływ narady inżynierskiej](#54-przepływ-narady-inżynierskiej)
   - [5.5 Przepływ pracy ciągłej z modułem Automations](#55-przepływ-pracy-ciągłej-z-modułem-automations)
7. [Stany i powiązania](#6-stany-i-powiązania)
   - [6.1 Stany przebiegu pętli wykonawczej](#61-stany-przebiegu-pętli-wykonawczej)
   - [6.2 Stany zadania w kolejce pętli](#62-stany-zadania-w-kolejce-pętli)
   - [6.3 Stany karty sesji](#63-stany-karty-sesji)
   - [6.4 Stany repozytorium i kontroli wersji](#64-stany-repozytorium-i-kontroli-wersji)
   - [6.5 Powiązania między oknami](#65-powiązania-między-oknami)
8. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
9. [Scenariusze użycia](#8-scenariusze-użycia)
   - [8.1 Naprawa usterki zgłoszonej w repozytorium](#81-naprawa-usterki-zgłoszonej-w-repozytorium)
   - [8.2 Diagnoza niepowodzenia budowania](#82-diagnoza-niepowodzenia-budowania)
   - [8.3 Decyzja projektowa wsparta naradą modeli](#83-decyzja-projektowa-wsparta-naradą-modeli)
   - [8.4 Wydanie produktu](#84-wydanie-produktu)
   - [8.5 Praca nad dwoma zadaniami równolegle](#85-praca-nad-dwoma-zadaniami-równolegle)
   - [8.6 Nocny przebieg zadań środowiska](#86-nocny-przebieg-zadań-środowiska)
10. [Wykazy normatywne i kryteria odbioru](#9-wykazy-normatywne-i-kryteria-odbioru)
   - [9.1 Komendy kontraktu](#91-komendy-kontraktu)
   - [9.2 Żetony projektowe](#92-żetony-projektowe)
   - [9.3 Komponenty interfejsu](#93-komponenty-interfejsu)
   - [9.4 Etykiety interfejsu z prototypu](#94-etykiety-interfejsu-z-prototypu)
   - [9.5 Komunikaty](#95-komunikaty)
   - [9.6 Punkty łamania](#96-punkty-łamania)
   - [9.7 Kryteria odbioru](#97-kryteria-odbioru)
11. [Załącznik A. Skróty klawiszowe środowiska](#załącznik-a-skróty-klawiszowe-środowiska)
   - [A.1. Sterowanie kanałami komunikacji](#a1-sterowanie-kanałami-komunikacji)
   - [A.2. Nawigacja](#a2-nawigacja)
   - [A.3. Karty sesji i układ kolumn](#a3-karty-sesji-i-układ-kolumn)
   - [A.4. Praca inżynierska](#a4-praca-inżynierska)
   - [A.5. Funkcje eksperckie](#a5-funkcje-eksperckie)
12. [Załącznik B. Pełny wykaz komend kontraktu środowiska CodeStudio](#załącznik-b-pełny-wykaz-komend-kontraktu-środowiska-codestudio)
   - [Obszar `developer` — 50 komend](#obszar-developer--50-komend)
   - [Obszar `terminal` — 27 komend](#obszar-terminal--27-komend)
   - [Obszar `diagnostics` — 9 komend](#obszar-diagnostics--9-komend)

---

## Wprowadzenie

CodeStudio jest środowiskiem programistycznym i inżynierskim platformy Danaco Console. Użytkownik prowadzi w nim jednocześnie edycję kodu, wykonanie poleceń, diagnostykę, budowę produktu i naradę modeli — w jednej karcie sesji, bez przełączania się między aplikacjami i bez utraty kontekstu. Środowisko dostarcza to, co w innych warsztatach inżynierskich rozdzielone jest między menedżer okien systemu operacyjnego, menedżer sesji, listwę stanu edytora, paletę poleceń i ręczne przenoszenie artefaktów między narzędziami — jako spójną warstwę środowiska nadrzędną wobec modułów.

Osią pracy w CodeStudio są dwa kanały komunikacji operacyjnej. Chat Window jest głównym oknem komunikacji Użytkownika z Wykonawcą, centralnym punktem pracy i podstawowym mechanizmem sterowania wszystkimi procesami środowiska; zajmuje lewą kolumnę obszaru roboczego, stałą, o pełnej wysokości. Execution Loop Window jest oknem pętli wykonawczej między Koordynatorem a Wykonawcą; przy zadaniach programistycznych prowadzi dekompozycję zlecenia, kolejkę zadań, nadzór nad realizacją, kontrolę jakości, ponowienia i sterowanie przebiegiem. Otwiera się jako kolumna sąsiadująca z Chat Window.

Dokument czyta się w czterech warstwach nałożonych na siebie:

| Warstwa dokumentu | Co dostarcza | Rozdziały |
|---|---|---|
| Mechanizm | Przeznaczenie środowiska, moduły, kanały komunikacji, magistrala kontekstu | 1–2 |
| Forma | Komplet okien z warstwami widoczności, makiety tekstowe, katalogi elementów | 3–4 |
| Zachowanie w czasie | Przepływy pracy, stany, powiązania między oknami i modułami | 5–6 |
| Sterowanie | Klucze konfiguracji, scenariusze użycia, skróty klawiszowe | 7–8, Załącznik A |

---

## 1. Przeznaczenie i kontekst środowiska

### 1.1. Definicja środowiska

CodeStudio organizuje pracę inżynierską wokół zadania, nie wokół pojedynczego narzędzia. Jednostką organizacji pracy jest karta sesji powiązana z repozytorium i gałęzią; w obrębie karty osiem modułów środowiska pracuje na wspólnym kontekście — tym samym stanie repozytorium, tym samym zestawie zmiennych, tej samej kolejce błędów i tym samym koszyku kontekstu podawanym Wykonawcy.

| Cecha | Charakterystyka w CodeStudio |
|---|---|
| Środek ciężkości pracy | Warsztat wielookienny — kod, wykonanie i wynik widoczne równocześnie w sąsiadujących kolumnach |
| Kanał sterowania procesami | Chat Window — polecenia w języku naturalnym, zatwierdzanie i przerywanie działań |
| Kanał prowadzenia pracy | Execution Loop Window — dekompozycja zlecenia programistycznego, kolejka zadań, kontrola jakości, ponowienia |
| Układ obszaru roboczego | Wyłącznie pionowy, podział lewa–prawa; regulacji podlega szerokość kolumn |
| Powiązania międzymodułowe | Częste i wspierane magistralą kontekstu: błąd → Diagnostics, plik → Developer, polecenie → Terminal |
| Karty sesji | Odpowiadają zadaniom inżynierskim, powiązane z repozytorium i gałęzią |
| Stan procesów | Trwały po stronie serwera — rozłączenie klienta nie przerywa budowania, testów ani pętli wykonawczej |

### 1.2. Zakres środowiska i rozgraniczenie od modułów

Powłoka środowiska dostarcza modułom ramy: miejsce, w którym się rysują, sposób przełączania, wymianę kontekstu i wspólne kanały komunikacji. Nie powiela funkcji modułów.

| Obszar | Właściwe miejsce |
|---|---|
| Wnętrze edytora kodu, drzewo projektu, klient Git, debugger | Moduł Developer |
| Konsole, sesje terminala, procesy wykonawcze | Moduł Terminal |
| Analiza logów, agregacja stanu, ustalenia naprawcze | Moduł Diagnostics |
| Złożenie kompletnego produktu — warstwa kliencka (powłoka), warstwa serwerowa (rdzeń), wdrożenie | Moduł Apps |
| Narada wielu modeli, konsensus | Moduł Roundtable |
| Definicja i cykl życia agenta lub eksperta | Moduł Agents |
| Nakładka na model, Panel prowenancji | Okno Konfiguracji |
| Konto, uwierzytelnianie, parowanie urządzeń, konta modeli | Okno Ustawień |

### 1.3. Położenie w hierarchii platformy

```
STRONA GŁÓWNA
    │  wybór środowiska w strefie 1
    ▼
CodeStudio
    │
    ├── Boczna nawigacja modułów (osiem pozycji — rozdz. 2)
    │
    ├── Obszar roboczy karty sesji — układ kolumnowy lewa–prawa
    │       Chat Window (lewa kolumna, stała)
    │       Execution Loop Window (kolumna sąsiadująca, otwierana)
    │       Okna robocze modułu (kolumna dominująca)
    │       Panele pomocnicze (kolejne kolumny boczne)
    │
    └── Karty sesji
            każda karta = jedno zadanie inżynierskie, z własnym
            repozytorium, gałęzią, układem kolumn i kontekstem

PONAD CAŁĄ STRUKTURĄ
    Always On Display — agent towarzyszący, świadomy kontekstu inżynierskiego
    Mobile — podgląd procesów i interwencja zdalna
    Automations — harmonogram i wykonanie cykliczne zadań środowiska
```

### 1.4. Kanoniczne rozmieszczenie obszaru roboczego

| Kolumna | Zawartość | Regulacja |
|---|---|---|
| Boczna nawigacja modułów | Osiem pozycji modułów, grupy, przypięcia, wskaźniki aktywności | Szerokość; tryb wąski (ikony) i szeroki (ikona z nazwą) |
| Lewa, stała, pełna wysokość obszaru roboczego | **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca | Wyłącznie szerokość |
| Kolumna sąsiadująca (otwierana) | **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca | Wyłącznie szerokość |
| Kolumna dominująca | Obszar roboczy modułu — okna edycyjne, wykonawcze, podglądu i monitory | Szerokość; podział na kolumny zagnieżdżone |
| Kolejne kolumny boczne | Okna pomocnicze otwierane jako rozszerzenia boczne, po prawej stronie obszaru roboczego | Wyłącznie szerokość |

---

## 2. Moduły dostępne w środowisku

### 2.1. Osiem modułów CodeStudio

| Moduł | Rola w środowisku | Okna wiodące modułu | Wkład do magistrali kontekstu |
|---|---|---|---|
| **Workspace** | Materiały projektu, notatki, dokumentacja zadania | Menedżer dokumentów, Podgląd dokumentu | Dokumenty i fragmenty jako kontekst dla Wykonawcy |
| **Roundtable** | Narada wielu modeli nad zagadnieniem inżynierskim | Okno narady, Moderator Panel, Consensus Panel | Ustalenia narady jako kontekst decyzji projektowej |
| **Design** | Warstwa wizualna produktu, komponenty, tokeny | Podgląd projektu, Katalog komponentów | Komponenty i tokeny jako kontekst implementacji |
| **Terminal** | Wykonanie poleceń, procesy, sesje powłoki | Konsola, Menedżer procesów | Wyniki poleceń, kody wyjścia, strumienie wyjścia |
| **Developer** | Edycja kodu, drzewo projektu, kontrola wersji, debugowanie | Project Tree, Code Editor, Git Panel, Debugger | Pliki, zaznaczenia, gałąź, zmiany do zatwierdzenia |
| **Diagnostics** | Analiza logów, agregacja stanu, ustalenia naprawcze | Panel Problemy, Analizator logów | Błędy i ostrzeżenia do wspólnego bufora środowiska |
| **Apps** | Złożenie produktu, budowanie, wdrożenie | Build Output, Panel wdrożenia | Wynik budowania, artefakty, status wdrożenia |
| **Agents** | Definicja agentów i ekspertów obsługujących zadania środowiska | Katalog agentów, Edytor agenta | Definicje wykonawców przypisywanych do zadań pętli |

Mapa modułów w siatce przedsionka (`.pd-siatka`, `design/05-okna/srodowiska/codestudio-przedsionek.html`), w kolejności występowania w źródle:

```
 Mapa modułów środowiska CodeStudio — siatka kafli przedsionka
 ┌───────────────┬───────────────┬───────────────┬───────────────┐
 │  Developer    │  Terminal     │  Diagnostics  │  Apps         │
 │  drzewo, kod, │  polecenia,   │  logi, stan,  │  budowa,      │
 │  Git, debug   │  procesy      │  naprawy      │  wydanie      │
 ├───────────────┼───────────────┼───────────────┼───────────────┤
 │  Design       │  Workspace    │  Roundtable   │  Agents       │
 │  makiety,     │  materiały,   │  narada       │  definicje    │
 │  komponenty   │  notatki      │  wielomodelowa│  wykonawców   │
 └───────────────┴───────────────┴───────────────┴───────────────┘
```

### 2.2. Boczna nawigacja modułów

| Funkcja nawigacji | Działanie | Warstwa | Sposób wywołania |
|---|---|---|---|
| Lista pozycji modułów | Osiem pozycji `.dn-szyna-poz.dn-szyna-poz--modul` ze stanem wybranym | 1 | Widoczna bez interakcji |
| Wskaźnik aktywności modułu | Plakietka z liczbą aktywnych procesów (Terminal), błędów (Diagnostics), zmian (Developer) | 1 | Widoczna bez interakcji |
| Grupy modułów | Sekcje „Kod” (Developer, Terminal, Diagnostics), „Produkt” (Apps, Design), „Współpraca” (Roundtable, Workspace, Agents), zwijane nagłówkami | 1 | Widoczna bez interakcji; nagłówek grupy zwija sekcję |
| Sekcja „Ostatnie” | Moduły uporządkowane wg czasu ostatniego otwarcia w karcie | 2 | Przełącznik widoczności w nagłówku panelu |
| Przypinanie i kolejność | Przypięcie modułu na początku listy, przeciąganie pozycji | 3 | Menu kontekstowe pozycji |
| Otwarcie w nowej karcie lub kolumnie obok | Otwarcie modułu jako nowa karta sesji albo jako kolumna sąsiadująca | 3 | Menu kontekstowe pozycji |
| Podgląd modułu | Mini-podgląd zestawu okien i stanu modułu bez przeładowania | 2 | Wskazanie kursorem pozycji |
| Szybkie przełączanie modułu | Lista „Idź do modułu…” z filtrem tekstowym | 2 | Skrót klawiszowy, paleta poleceń |
| Ukrycie pozycji | Pozycja znika z listy, moduł pozostaje osiągalny z palety poleceń i przełącznika „Idź do modułu” | 3 | Menu kontekstowe pozycji |
| Tryb wąski i szeroki | Panel z samymi ikonami albo z ikoną i nazwą; zwinięcie do krawędzi | 2 | Przycisk zwinięcia w nagłówku panelu |
| Przełącznik środowiska | Przejście do innego środowiska platformy bez powrotu na stronę główną | 2 | Nagłówek panelu — odznaka „CodeStudio ▼” |

### 2.3. Magistrala kontekstu środowiska

Mechanizm magistrali kontekstu — przeznaczenie, sposób przekazania artefaktu, zasięg, warstwę widoczności i sposób wywołania oraz relacje do modułu Library, punktów izolacji i encji artefaktu — opisuje [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) (rozdz. 6a). Poniżej podano wyłącznie specyfikę środowiska CodeStudio: rodzaje artefaktów inżynierskich krążących magistralą oraz mechanizmy karty sesji, które z niej korzystają.

| Mechanizm | Działanie | Warstwa | Sposób wywołania |
|---|---|---|---|
| Koszyk kontekstu | Wspólny zestaw kontekstu karty: pliki, zaznaczenia, błędy, wyniki, fragmenty — podawany Wykonawcy w Chat Window i pętli wykonawczej | 2 | Znacznik kontekstu „Koszyk (N)” nad polem poleceń Chat Window |
| Przekaż do modułu | Skierowanie artefaktu do wskazanego modułu: błąd z Build Output do Diagnostics, plik z Project Tree do narady Roundtable | 3 | Menu kontekstowe artefaktu |
| Wspólne zaznaczenie kodu | Zaznaczenie w Code Editor dostępne jako kontekst dla Chat Window, Execution Loop Window, Roundtable i Diagnostics | 1 | Dostępne bez interakcji po zaznaczeniu |
| Wspólny bufor błędów | Błędy budowania, testów i analizy statycznej spływają do jednej kolejki środowiska, z której korzystają moduły i Wykonawca | 1 | Widoczny jako licznik w listwie stanu |
| Wspólne repozytorium i gałąź | Repozytorium i gałąź powiązane z kartą są wspólne dla wszystkich modułów karty | 1 | Znacznik kontekstu `[repo] [gałąź]` |
| Wspólne zmienne środowiskowe | Zestaw zmiennych karty dostępny dla Terminal, Developer i Apps; klucze jawne, wartości maskowane | 3 | Menu szybkiej konfiguracji sesji |
| Pamięć kontekstu sesji | Trwały kontekst karty przekazywany modelowi, spójny z ustawieniem izolacji | 4 | Okno konfiguracji, polecenie języka naturalnego |
| Podgląd złożonego kontekstu | Panel prowenancji pokazujący, co zostało podane modelowi w bieżącym wywołaniu | 4 | Polecenie w Chat Window, okno Konfiguracji |
| Reguły przepływu artefaktów | Jawne reguły kierowania artefaktów — nowy błąd budowania trafia do koszyka Diagnostics | 4 | Okno konfiguracji, klucz `codestudio.kontekst.reguly` |

---

### 2.4. Przedsionek środowiska — widok wejściowy

**Zasada wejścia.** Kliknięcie karty środowiska na stronie głównej otwiera
**przedsionek środowiska**, nie przestrzeń roboczą. Operator nie wskazał
jeszcze modułu, więc żaden moduł nie zostaje otwarty; otwarcie pierwszego
w kolejności byłoby decyzją podjętą za Operatora. Przedsionek pełni wobec
modułów tę samą rolę, którą Centrum dowodzenia pełni wobec środowisk:
pokazuje zbiór i pozwala wybrać, nie wybiera sam.

**Trzy rejony przedsionka.**

| Rejon | Zawartość | Waga |
|---|---|---|
| Szyna sesji („Twoje sesje”) | Wykaz sesji Operatora w tym środowisku, zgrupowany repozytoriami, z filtrem zakresu (czynne · wszystkie · zakończone) i wskaźnikiem pracy w tle. W stopce — wejście do pełnej historii sesji | Kolumna przy lewej krawędzi, waga średnia |
| Płótno — nagłówek środowiska | Godło, nazwa, motto i zdanie o przedmiocie pracy środowiska | Waga główna |
| Płótno — kafle modułów | Osiem kafli modułów, każdy z ikoną, nazwą, zdaniem o przeznaczeniu i miarą bieżącego obłożenia | Siatka kafli, waga pośrednia |
| Płótno — listwa działań środowiska | Trzy pozycje: **Nowe repozytorium**, **Konfiguracja środowiska**, **Ustawienia** — działania dotyczące środowiska jako całości | Listwa jednorzędowa, waga najniższa |

Nagłówek środowiska niesie brzmienie dosłowne z prototypu: nadtytuł „Środowisko pracy”, tytuł „CodeStudio”, motto „Projektuj. Buduj. Rozwijaj.”, opis „Wytwarzanie oprogramowania: repozytorium, edycja kodu, powłoka systemowa, diagnostyka i wydanie produktu. Osiem modułów wokół jednego repozytorium.”

```
 Szkic przedsionka środowiska CodeStudio
 ┌──────────────────┬────────────────────────────────────────────────────────┐
 │ SZYNA SESJI      │  ⬡  Środowisko pracy                                   │
 │ „Twoje sesje”    │     CODESTUDIO                                         │
 │ ┌─ 3 repo ·7 ses │     Projektuj. Buduj. Rozwijaj.                        │
 │ [Czynne][Wszyst] │  ────────────────────────────────────────────────────  │
 │ [Zakończone]     │  Moduły środowiska                                     │
 │                  │  ┌──────────┬──────────┬──────────┬──────────┐         │
 │ danaco-console   │  │Developer │Terminal  │Diagnost. │   Apps   │         │
 │ ● Naprawa        │  ├──────────┼──────────┼──────────┼──────────┤         │
 │   wycieku pamięci│  │ Design   │Workspace │Roundtable│ Agents   │         │
 │   Developer·prac.│  └──────────┴──────────┴──────────┴──────────┘         │
 │ ○ Refaktoryzacja │                                                        │
 │   warstwy        │  ────────────────────────────────────────────────────  │
 │   kontraktów     │  [ Nowe repozytorium ] [ Konfiguracja ] [ Ustawienia]  │
 │ ──────────────   │                                                        │
 │ [Pełna historia  │                                                        │
 │  sesji]          │                                                        │
 └──────────────────┴────────────────────────────────────────────────────────┘
   pd-szyna              pd-plotno (pd-naglowek → pd-siatka → pd-listwa)
```

**Czego przedsionek nie zawiera.** Chat Window, Execution Loop Window,
bocznej nawigacji modułów, pasa kart sesji ani żadnego okna operacyjnego. Wszystkie te
elementy należą do przestrzeni roboczej i pojawiają się dopiero po wyborze
modułu. Moduł Workspace występuje w przedsionku **wyłącznie jako kafel**;
jego okna otwierają się po wejściu w moduł, nie wcześniej.

**Dwie drogi wyjścia z przedsionka.** Kliknięcie kafla otwiera przestrzeń
roboczą z tym modułem jako wiodącym i zakłada pierwszą kartę sesji. Kliknięcie
pozycji w szynie sesji wznawia sesję istniejącą — przestrzeń robocza otwiera
się od razu w module tej sesji, wraz z jej układem, historią i kontekstem.

**Powrót do przedsionka.** Przedsionek jest osiągalny z przestrzeni roboczej
przez kontrolkę trybów na pasku (ponowne wskazanie bieżącego środowiska) oraz
przez powrót ze strony głównej. Powrót nie zamyka kart sesji — trwają one
w tle i pozostają widoczne w szynie jako sesje czynne.

**Prototyp odniesienia:** `design/05-okna/srodowiska/codestudio-przedsionek.html`.

Przebieg wejścia do środowiska opiera się na komendzie kontraktu `environment.enter`
(obszar `environment`, `budowa/shared/contract.json`); pole `sessionId` puste otwiera
kartę pustą — przedsionek — zamiast wznawiać sesję istniejącą:

```
 Operator                Powłoka                        Rdzeń
    │                       │                              │
    │  kliknięcie karty     │                              │
    │  środowiska           │                              │
    ├──────────────────────►│                              │
    │                       │  environment.enter           │
    │                       │  { environmentId, clientId,  │
    │                       │    sessionId: "" }           │
    │                       ├─────────────────────────────►│
    │                       │                              │
    │                       │  environment · modules: []   │
    │                       │  sessions · focusedSessionId:│
    │                       │  "" (puste)                  │
    │                       │◄─────────────────────────────┤
    │  przedsionek          │                              │
    │  (szyna sesji +       │                              │
    │   kafle modułów)      │                              │
    │◄──────────────────────┤                              │
    │                       │                              │
    │  kliknięcie kafla     │                              │
    │  modułu               │                              │
    ├──────────────────────►│                              │
    │                       │  environment.enter           │
    │                       │  { …, sessionId: "<nowa>" }  │
    │                       ├─────────────────────────────►│
    │                       │  modules: [8 pozycji]        │
    │                       │  focusedSessionId: "<nowa>"  │
    │                       │◄─────────────────────────────┤
    │  przestrzeń robocza   │                              │
    │  modułu wybranego     │                              │
    │◄──────────────────────┤                              │
```

Legenda: `environment.enter` — komenda kontraktu obszaru `environment`; pola żądania
`environmentId`, `clientId`, `sessionId` (opcjonalne — puste otwiera kartę pustą, czyli
przedsionek); pola wyniku `environment`, `modules` (puste dla przedsionka, ośmioelementowe
po wyborze modułu), `sessions`, `focusedSessionId` (puste w przedsionku, ustawione po
wyborze modułu albo wznowieniu sesji z szyny).

---

## 3. Komplet okien środowiska

### 3.1. Warstwy widoczności w środowisku

Interfejs CodeStudio ujawnia możliwości środowiska stopniowo: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Liczba modułów, okien, paneli i ustawień nie wpływa na postrzeganą prostotę interfejsu — złożoność środowiska istnieje w architekturze i pozostaje niewidoczna do chwili wystąpienia potrzeby użycia danej funkcji.

| Warstwa | Zawartość w CodeStudio | Sposób dostępu |
|---|---|---|
| 1 — zawsze widoczna | Chat Window, boczna nawigacja modułów, aktywne okno wiodące modułu, znaczniki kontekstu, listwa stanu środowiska, wskaźniki przebiegu pętli wykonawczej | Widoczne bez interakcji; ponad 80% powierzchni interfejsu |
| 2 — widoczna na żądanie | Execution Loop Window, wybór modelu i wykonawcy, wybór repozytorium i gałęzi, przełącznik trybu pracy, poziom wysiłku, podgląd modułu, przełącznik szerokości kolumn | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 — rozwinięcia kontekstowe | Zestawy operacji na artefaktach, ustawienia szybkie sesji, warianty uruchomienia zadania, operacje na kartach sesji, operacje na pozycjach nawigacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 — funkcje eksperckie | Panel prowenancji, ustawienia izolacji kontekstu, tryb klawiatury modalnej, reguły przepływu artefaktów, mapa skrótów, narzędzia diagnostyczne niskiego poziomu, tryb administracyjny | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, okno konfiguracji, konfiguracja roli |

Mechanizmy ukrywania stosowane w środowisku:

- **Menu progresywne.** Wybór wykonawcy, modelu i profilu uruchomienia prezentowany jest jako jeden zwinięty element (`Wykonawca ▼`), rozwijany kliknięciem.
- **Panele wysuwane.** Okna pomocnicze otwierają się jako kolumny boczne po prawej stronie obszaru roboczego i po zamknięciu znikają całkowicie z przestrzeni roboczej.
- **Grupowanie logiczne akcji.** Budowanie, testy, lintowanie i wdrożenie występują jako jeden element zbiorczy (`Zadania ▼`) z pełną listą w rozwinięciu.
- **Znaczniki kontekstowe.** Środowisko, repozytorium, gałąź, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [danaco-console] [feature/auth] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor.

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu nie występuje.

### 3.2. Tabela kompletu okien

| Okno | Kanał / moduł | Rola w środowisku | Położenie w układzie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| **Chat Window** | Kanał Użytkownik ↔ Wykonawca | Główne okno komunikacji, centralny punkt pracy i podstawowy mechanizm sterowania procesami środowiska | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji |
| **Execution Loop Window** | Kanał Koordynator ↔ Wykonawca | Pętla wykonawcza zadania programistycznego: dekompozycja, kolejka, nadzór, kontrola jakości, ponowienia, sterowanie przebiegiem | Kolumna sąsiadująca z Chat Window | 2 | Przycisk „Pętla wykonawcza” w listwie stanu; skrót klawiszowy; automatycznie po przyjęciu zlecenia wieloetapowego |
| **Project Tree** | Developer | Struktura repozytorium, wybór plików do edycji i do koszyka kontekstu | Kolumna dominująca, region lewy | 1 | Widoczne po wejściu w moduł Developer |
| **Code Editor** | Developer | Edycja plików źródłowych, zaznaczenia podawane jako kontekst | Kolumna dominująca, region główny | 1 | Otwarcie pliku z Project Tree, palety poleceń albo polecenia w Chat Window |
| **Git Panel** | Developer | Gałęzie, zmiany do zatwierdzenia, historia, operacje kontroli wersji | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik gałęzi w listwie stanu; skrót klawiszowy |
| **Debugger** | Developer | Punkty wstrzymania, stos wywołań, zmienne, krokowanie | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Uruchomienie profilu debugowania z menu `Zadania ▼` |
| **Konsola** | Terminal | Wykonanie poleceń powłoki, strumienie wyjścia procesów | Kolumna dominująca, region wyników | 1 | Wejście w moduł Terminal; skrót klawiszowy; polecenie w Chat Window |
| **Menedżer procesów** | Terminal | Wykaz żywych procesów karty i całego środowiska | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Segment „Aktywne procesy” listwy stanu |
| **Build Output** | Apps | Przebieg i wynik budowania, artefakty, błędy budowania | Kolumna dominująca, region wyników | 1 | Uruchomienie zadania budowania; segment statusu budowania |
| **Panel wdrożenia** | Apps | Stan wdrożenia produktu, cele wdrożenia, dziennik operacji | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Menu `Zadania ▼` → wdrożenie |
| **Panel Problemy** | Diagnostics | Wspólny bufor błędów i ostrzeżeń środowiska, skok do źródła | Kolumna boczna, otwierana jako rozszerzenie boczne | 1 | Segment diagnostyki w listwie stanu |
| **Analizator logów** | Diagnostics | Agregacja logów, korelacja zdarzeń, ustalenia naprawcze | Kolumna dominująca, region wyników | 2 | Wejście w moduł Diagnostics; przekazanie błędu magistralą kontekstu |
| **Okno narady** | Roundtable | Równoległe odpowiedzi kilku modeli na zagadnienie inżynierskie | Kolumna dominująca, region główny | 2 | Wejście w moduł Roundtable; przekazanie artefaktu „Wyślij do Roundtable” |
| **Menedżer dokumentów** | Workspace | Materiały i notatki zadania, dokumentacja repozytorium | Kolumna dominująca, region główny | 2 | Wejście w moduł Workspace |
| **Podgląd projektu** | Design | Widok warstwy wizualnej produktu, komponenty, tokeny | Kolumna dominująca, region główny | 2 | Wejście w moduł Design |
| **Katalog agentów** | Agents | Wykaz agentów i ekspertów przypisywanych do zadań pętli wykonawczej | Kolumna dominująca, region główny | 2 | Wejście w moduł Agents; selektor `Wykonawca ▼` |
| **Koszyk kontekstu** | Powłoka środowiska | Zestaw artefaktów podawanych Wykonawcy jako kontekst | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik „Koszyk (N)” nad polem poleceń Chat Window |
| **Paleta poleceń** | Powłoka środowiska | Wywołanie dowolnego polecenia powłoki i modułów, nawigacja bez myszy | Nakładka nad obszarem roboczym | 4 | Skrót `Ctrl/Cmd + K` ([Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md), rozdz. 15.8) |
| **Dziennik zdarzeń** | Powłoka środowiska | Budowania, testy, zatwierdzenia, uruchomienia automatyk, powiadomienia modułów | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Segment powiadomień listwy stanu |
| **Menedżer układów kolumn** | Powłoka środowiska | Zapisane układy kolumn wywoływane jednym poleceniem | Panel popover nad listwą stanu | 3 | Menu `Układ ▼` w listwie stanu; paleta poleceń |
| **Szybka konfiguracja sesji** | Powłoka środowiska | Model i kanał, układ kolumn, izolacja karty, powiązane repozytorium, zestaw zmiennych | Panel popover przy belce tytułowej | 3 | Znacznik kontekstu w belce tytułowej |
| **Panel prowenancji** | Okno Konfiguracji | Podgląd złożonego kontekstu przekazanego modelowi | Okno platformowe | 4 | Polecenie w Chat Window; okno Konfiguracji |
| **Mapa skrótów** | Okno Konfiguracji | Rejestr skrótów środowiska z edycją i wykrywaniem konfliktów | Okno platformowe | 4 | Skrót klawiszowy; wyszukiwarka funkcji |
| **Always On Display** | Warstwa globalna | Agent towarzyszący świadomy kontekstu inżynierskiego — komentuje przebieg budowania, testów i pętli wykonawczej | Element pływający nad obszarem roboczym | 2 | Kliknięcie awatara; polecenie głosowe |
| **Mobile** | Warstwa globalna | Podgląd procesów środowiska i interwencja zdalna | Kanał zewnętrzny | 2 | Sparowane urządzenie; segment procesów listwy stanu |

### 3.3. Listwa stanu środowiska

Listwa stanu jest ciągiem segmentów osadzonym w pasku stanu okna. Agreguje sygnały modułów i sesji; każdy segment jest klikalny i otwiera właściwe okno lub menu. Listwa informuje i skraca drogę — nigdy nie blokuje.

| Segment | Co pokazuje | Kliknięcie | Warstwa | Źródło |
|---|---|---|---|---|
| Gałąź Git | Bieżąca gałąź i stan względem zdalnej | Menu gałęzi → Git Panel | 1 | Developer |
| Repozytorium karty | Repozytorium powiązane z kartą sesji | Selektor repozytorium | 1 | Powłoka |
| Diagnostyka | Liczba błędów i ostrzeżeń (✖ N ⚠ M) | Panel Problemy | 1 | Diagnostics |
| Status budowania | Wynik ostatniego budowania i czas trwania | Build Output | 1 | Developer, Apps |
| Status testów | Zdane, niezdane i pominięte z ostatniego przebiegu | Panel testów | 1 | Developer |
| Pętla wykonawcza | Stan pętli i liczba zadań w kolejce | Execution Loop Window | 1 | Koordynator |
| Aktywne procesy | Liczba żywych procesów Terminala, automatyk i narad | Menedżer procesów | 1 | Terminal, Automations |
| Model i kanał sesji | Aktywny wykonawca i kanał modelu bieżącej karty | Szybka konfiguracja sesji | 2 | Model konfiguracji |
| Zmiany do zatwierdzenia | Liczba plików do zatwierdzenia | Git Panel | 2 | Developer |
| Pozycja kursora | Wiersz i kolumna, zaznaczenie, wcięcia, kodowanie, znak końca linii | Menu wcięć i kodowania | 2 | Code Editor |
| Wskaźnik izolacji | Czy karta działa w izolacji kontekstu i zasobów; stan wyjściowy: wyłączona | Okno konfiguracji → izolacja | 3 | Izolacja i zależności |
| Powiadomienia | Licznik nieprzeczytanych zdarzeń środowiska | Dziennik zdarzeń | 3 | Powłoka |
| Menu `Układ ▼` | Zapisane układy kolumn karty | Menedżer układów kolumn | 3 | Powłoka |
| Menu `Zadania ▼` | Budowanie, testy, lintowanie, uruchomienie, wdrożenie, debugowanie | Lista zadań środowiska | 3 | Powłoka, Terminal |

---

## 4. Specyfikacja okien

### 4.1. Makieta główna — pełny układ obszaru roboczego

Makieta przedstawia interfejs w stanie spoczynku: widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

```
 Makieta — Środowisko CodeStudio, pełen układ karty sesji
 ═════════════════════════════════════════════════════════════════════════════════════════════
  Belka + wstążka │ Danaco Console  [CodeStudio ▼] [danaco-console] [feature/auth] [Fable 5 ▼]
                  │ Karty sesji: [ Fix #421 auth ] [ Wydanie 2.0 ] [ + ]
                  │ Listwa stanu: ⑂ feature/auth · ✖2 ⚠5 · build ✔ 41 s · testy 128/130
                  │               · pętla ● 3 zadania · procesy 2 · [Układ ▼] [Zadania ▼] ⋮
 ─────────────────────────────────────────────────────────────────────────────────────────────
  Boczna     │ Chat Window          │ Execution Loop      │ Obszar roboczy       │ Panel
  nawigacja  │ Użytkownik ↔         │ Koordynator ↔       │ modułu Developer     │ pomoc-
  modułów    │ Wykonawca            │ Wykonawca           │                      │ niczy
             │                      │                     │  Project Tree │ Code │ (roz-
  ▸ Workspace│ ▸ historia poleceń   │ Zlecenie:           │   src/        │ Edi- │ szerze-
  ▸ Roundtab.│   i wyników          │  „Napraw logowanie  │    auth/      │ tor  │ nie
  ▸ Design   │ ▸ strumień odpowiedzi│   przez token”      │     token.ts  │      │ boczne)
  ▸ Terminal │ ▸ zatwierdzanie i    │                     │     login.ts  │ ───  │
  ▸ Developer│   przerywanie działań│ Kolejka zadań:      │ ───────────── │ Kon- │ Git
    ●2       │                      │  ✓ analiza modułu   │ Build Output  │ sola │ Panel
  ▸ Diagnost.│ ───────────────────  │  ● poprawka token.ts│  ✔ 41 s       │      │
    ✖2       │ [Koszyk (4) ▼]       │  ⏸ testy jednostkowe│               │      │ ⋮
  ▸ Apps     │ [ 📎 ] Wpisz          │  ⏸ przegląd zmian   │               │      │
  ▸ Agents   │ polecenie…      [ ➤ ]│ [⏸] [▶] [■] [↻]     │               │      │
 ═════════════════════════════════════════════════════════════════════════════════════════════
   Lewa kolumna stała,  Kolumna sąsiadująca,   Kolumna dominująca,     Kolumny boczne,
   pełna wysokość       otwierana              regulowana szerokość    rozszerzenia
```

### 4.2. Chat Window — główne okno komunikacji Użytkownik ↔ Wykonawca

Chat Window jest centralnym punktem pracy w CodeStudio i podstawowym mechanizmem sterowania wszystkimi procesami środowiska. Przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Zajmuje lewą kolumnę obszaru roboczego, stałą, o pełnej wysokości; regulacji podlega wyłącznie szerokość. Okno występuje w tym samym miejscu układu we wszystkich ośmiu modułach środowiska.

```
 Makieta — Chat Window (stan spoczynku)
 ─────────────────────────────────────────────────────────────────
  Nagłówek okna   │ Wykonawca: Fable 5 ▼ · kanał: API · Ultra ▼
 ─────────────────────────────────────────────────────────────────
  Strumień        │ Użytkownik: Napraw odświeżanie tokenu
  komunikacji     │ Wykonawca: analizuję auth/token.ts…
                  │ Wykonawca: przygotowałem zmianę w 2 plikach
                  │            [ Pokaż różnice ] [ Zatwierdź ✓ ]
                  │ Użytkownik: uruchom testy jednostkowe
                  │ Wykonawca: 128/130 zdanych, 2 niezdane
                  │            [ Wyślij do Diagnostics ] ⋮
 ─────────────────────────────────────────────────────────────────
  Pasek kontekstu │ [danaco-console] [feature/auth] [Koszyk (4) ▼]
 ─────────────────────────────────────────────────────────────────
  Pole poleceń    │ [ 📎 ] Wpisz polecenie…            [ ⏹ ] [ ➤ ]
 ─────────────────────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Strumień komunikacji | Przewijana historia poleceń Użytkownika i wyników Wykonawcy | Podgląd przebiegu pracy i wyników w jednym miejscu | Obszar przewijany złożony z wpisów `.dn-wpis` (warianty `--czlowiek` · `--inteligencja` · `--system` · `--pracuje`); kontener regionu **[DO DECYZJI OPERATORA]** — arkusze źródłowe nie definiują odrębnej klasy kontenera strumienia | Domyślny · odpowiedź w toku · przerwany · błąd kanału | 1 | Widoczny bez interakcji |
| Pole poleceń | Pole tekstowe przyjmujące polecenia w języku naturalnym | Zlecanie pracy Wykonawcy i sterowanie procesami środowiska | Pole (`.dn-pole`, kontrolka `.dn-pole-kontrolka`) z przyciskiem wysłania | Puste · z treścią · wysyłanie · zablokowane kanałem | 1 | Widoczne bez interakcji |
| Przycisk przerwania | Kontrolka zatrzymania bieżącego działania Wykonawcy | Natychmiastowe przerwanie generowania lub wykonania | Ikona (`stop`) w polu poleceń | Nieaktywny wizualnie przy braku działania · aktywny w trakcie | 1 | Widoczny w trakcie działania |
| Kontrolka zatwierdzenia wyniku | Przycisk akceptacji zmiany przygotowanej przez Wykonawcę | Wprowadzenie zmiany do repozytorium karty | Przycisk (`.dn-btn.dn-btn--sygnal`) w wiadomości | Domyślny · zatwierdzony · wycofany („Cofnij”) | 1 | Pojawia się przy wiadomości z wynikiem |
| Podgląd różnic | Rozwijany blok z różnicami plików objętych zmianą | Ocena zakresu zmiany przed zatwierdzeniem | Blok rozwijany w strumieniu | Zwinięty · rozwinięty | 2 | Przycisk „Pokaż różnice” |
| Pasek kontekstu | Ciąg znaczników repozytorium, gałęzi i koszyka kontekstu | Podgląd i zmiana kontekstu przekazywanego Wykonawcy | Znaczniki (`.dn-etykietka`) | Domyślny · aktywny selektor | 1 | Widoczny bez interakcji; kliknięcie otwiera selektor |
| Selektor wykonawcy | Zwinięty wybór agenta lub modelu bazowego | Ustalenie, kto realizuje polecenia w tej karcie | Menu progresywne `Wykonawca ▼` | Domyślny · rozwinięty · zmieniony | 2 | Kliknięcie w nagłówku okna |
| Selektor poziomu wysiłku | Zwinięty wybór intensywności pracy modelu | Dopasowanie nakładu do złożoności zadania | Menu progresywne `Ultra ▼` | Domyślny · rozwinięty | 2 | Kliknięcie w nagłówku okna |
| Załącznik kontekstu | Kontrolka dodania pliku lub fragmentu do koszyka | Uzupełnienie kontekstu bieżącego polecenia | Ikona (`spinacz`) | Domyślny · z liczbą pozycji | 2 | Kliknięcie w polu poleceń |
| Menu operacji wiadomości | Zestaw operacji na pojedynczym wyniku: przekaż do modułu, skopiuj, wyjaśnij, ponów | Kierowanie wyniku do właściwego narzędzia | Menu kebab (⋮) | Zwinięte · rozwinięte | 3 | Kliknięcie ⋮ przy wiadomości |
| Uchwyt szerokości kolumny | Pionowy uchwyt na prawej krawędzi okna | Regulacja szerokości lewej kolumny | Uchwyt przeciągania krawędzi kolumny — **[DO DECYZJI OPERATORA]**, brak klasy uchwytu podziału w arkuszach źródłowych | Domyślny · przeciąganie | 2 | Przeciągnięcie krawędzi |

### 4.3. Execution Loop Window — okno pętli wykonawczej Koordynator ↔ Wykonawca

Execution Loop Window prowadzi pętlę wykonawczą zadania programistycznego. Koordynator przyjmuje zlecenie sformułowane przez Użytkownika w Chat Window, dokonuje jego dekompozycji na zadania, buduje kolejkę, przydziela zadania Wykonawcy, nadzoruje realizację, przeprowadza kontrolę jakości wyniku, zarządza ponowieniami i udostępnia sterowanie przebiegiem. Okno otwiera się jako kolumna sąsiadująca z Chat Window; regulacji podlega wyłącznie szerokość.

```
 Makieta — Execution Loop Window (stan spoczynku)
 ─────────────────────────────────────────────────────────────────────────
  Nagłówek okna   │ Pętla wykonawcza · Koordynator ↔ Wykonawca
                  │ przebieg #14 · ● uruchomiona
 ─────────────────────────────────────────────────────────────────────────
  Zlecenie        │ „Napraw odświeżanie tokenu i pokryj testami”
  i dekompozycja  │ źródło: Chat Window · repo: danaco-console
                  │ [ Skoryguj zlecenie ]
 ─────────────────────────────────────────────────────────────────────────
  Kolejka zadań   │ ✓ 1. Analiza modułu auth              12 s
                  │ ● 2. Poprawka auth/token.ts        w toku
                  │ ⏸ 3. Testy jednostkowe            oczekuje
                  │ ⏸ 4. Przegląd zmian               oczekuje
 ─────────────────────────────────────────────────────────────────────────
  Wymiana         │ Koordynator → Wykonawca: zadanie 2, zakres…
  komunikatów     │ Wykonawca → Koordynator: zmiana w 2 plikach
  sterujących     │ Koordynator: kontrola jakości — 1 uwaga
                  │ Koordynator → Wykonawca: ponowienie zadania 2
 ─────────────────────────────────────────────────────────────────────────
  Wskaźniki       │ zadania 1/4 · ponowienia 1 · czas 03:12
  przebiegu       │ kontrola jakości: 1 uwaga otwarta
 ─────────────────────────────────────────────────────────────────────────
  Sterowanie      │ [ ⏸ Wstrzymaj ] [ ▶ Wznów ] [ ■ Przerwij ] [ ↻ Ponów ]
 ─────────────────────────────────────────────────────────────────────────
```

| Element okna | Co to jest | Do czego służy | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Nagłówek pętli | Identyfikacja przebiegu i jego stan | Rozpoznanie, który przebieg jest oglądany | Nagłówek okna z plakietką stanu | Nieuruchomiona · uruchomiona · wstrzymana · zakończona · błąd | 1 | Widoczny po otwarciu okna |
| Blok zlecenia | Treść zlecenia przyjętego od Użytkownika wraz ze źródłem i kontekstem | Jednoznaczne ustalenie, co realizuje pętla | Blok zawartości (`.dn-karta`) | Domyślny · skorygowany | 1 | Widoczny po otwarciu okna |
| Dekompozycja zlecenia | Rozbicie zlecenia na zadania wykonawcze przez Koordynatora | Podział pracy programistycznej na kroki nadzorowane osobno | Lista pozycji z numeracją | Wyliczona · zmieniona po korekcie zlecenia | 1 | Powstaje po przyjęciu zlecenia |
| Kolejka zadań | Uporządkowany wykaz zadań ze stanem i czasem realizacji | Podgląd postępu i kolejności pracy | Lista kroków (`.dn-kolejka` z pozycjami `.dn-krok`, warianty `--pracuje` · `--poprawny` · `--bledy` · `--wstrzymany`) | Oczekujące · w toku · zakończone · niezdane · pominięte | 1 | Widoczna po otwarciu okna |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator ↔ Wykonawca | Wgląd w przydziały, wyniki cząstkowe i decyzje sterujące | Obszar strumienia | Domyślny · nowy komunikat · zwinięty do podsumowań | 1 | Widoczna po otwarciu okna |
| Wynik kontroli jakości | Ustalenia kontroli wyniku zadania wraz z uwagami | Rozstrzygnięcie, czy zadanie przechodzi dalej, czy wraca do ponowienia | Blok ustaleń z licznikiem uwag | Brak uwag · uwagi otwarte · uwagi rozstrzygnięte | 1 | Powstaje po zakończeniu zadania |
| Decyzja o ponowieniu | Zapis ponowienia zadania wraz z uzasadnieniem | Śledzenie, ile razy i dlaczego zadanie było powtarzane | Pozycja strumienia z licznikiem | Domyślna · przekroczony limit ponowień | 1 | Powstaje po negatywnej kontroli jakości |
| Wskaźniki przebiegu pętli | Liczniki zadań, ponowień i czasu przebiegu | Ocena tempa i kosztu realizacji zlecenia | Pasek liczników | Domyślny · przebieg wstrzymany · przebieg zakończony | 1 | Widoczne po otwarciu okna |
| Sterowanie przebiegiem | Kontrolki wstrzymania, wznowienia, przerwania i ponowienia | Bezpośrednie sterowanie pętlą przez Użytkownika | Grupa przycisków (`.dn-btn`) | Domyślne · w trakcie operacji · niedostępna operacja bez zmiany stanu | 1 | Widoczne po otwarciu okna |
| Korekta zlecenia | Kontrolka zmiany treści zlecenia w trakcie przebiegu | Doprecyzowanie zakresu bez przerywania pętli | Przycisk (`.dn-btn--zarys`) | Domyślny · edycja · zastosowana | 2 | Przycisk „Skoryguj zlecenie” |
| Przydział wykonawcy zadania | Wybór agenta lub modelu realizującego wskazane zadanie kolejki | Dopasowanie wykonawcy do charakteru zadania | Menu progresywne przy pozycji kolejki | Domyślny · rozwinięty · zmieniony | 2 | Kliknięcie pozycji kolejki |
| Operacje na zadaniu | Zestaw operacji pozycji kolejki: przesuń, podziel, scal, pomiń, ponów, przekaż | Ręczna korekta planu pracy | Menu kebab (⋮) przy pozycji | Zwinięte · rozwinięte | 3 | Kliknięcie ⋮ przy pozycji kolejki |
| Reguły ponowień | Ustawienie liczby ponowień i warunku eskalacji | Ustalenie granicy automatycznego powtarzania zadania | Panel ustawień pętli | Domyślne · zmienione | 4 | Okno konfiguracji, polecenie w Chat Window |
| Uchwyt szerokości kolumny | Pionowy uchwyt na prawej krawędzi okna | Regulacja szerokości kolumny pętli | Uchwyt przeciągania krawędzi kolumny — **[DO DECYZJI OPERATORA]**, brak klasy uchwytu podziału w arkuszach źródłowych | Domyślny · przeciąganie | 2 | Przeciągnięcie krawędzi |

### 4.4. Okna robocze modułu Developer

```
 Makieta — Kolumna dominująca w module Developer (stan spoczynku)
 ══════════════════════════════════════════════════════════════════════════════
  Chat Window   │ Project Tree      │ Code Editor                  │ Git
  Użytkownik ↔  │  danaco-console/  │ [token.ts] [login.ts] [+]    │ Panel
  Wykonawca     │   src/            │ ─────────────────────────    │ (roz-
                │    auth/          │  1  export async function    │ szerze-
                │     token.ts  ●   │  2    refreshToken(t: string)│ nie
                │     login.ts      │  3  {                        │ boczne)
                │   tests/          │ …                            │
                │ ───────────────── │ ──────────────────────────   │ ⑂ feature
                │ Build Output ✔    │ Konsola  $ npm test          │   /auth
                │  41 s · 0 błędów  │  128 passed, 2 failed        │ 4 zmiany
 ══════════════════════════════════════════════════════════════════════════════
```

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Project Tree | Drzewo plików repozytorium karty | Nawigacja po strukturze projektu, wskazanie plików do edycji i do koszyka | Region kolumny dominującej, lista drzewiasta | Domyślny · plik zmieniony · plik otwarty · filtr aktywny | 1 | Widoczny w module Developer |
| Code Editor | Edytor plików źródłowych z zakładkami plików | Edycja kodu i tworzenie zaznaczeń podawanych jako kontekst | Region główny kolumny dominującej | Domyślny · zmiana niezapisana · tylko odczyt · konflikt scalania | 1 | Otwarcie pliku |
| Zakładki plików | Poziomy pasek otwartych plików wewnątrz Code Editora | Przełączanie między plikami bez utraty pozycji kursora | Zakładki (`.dn-zakladki`) | Domyślna · aktywna · zmieniona · przypięta | 1 | Widoczne po otwarciu pliku |
| Build Output | Region wyników z przebiegiem budowania | Podgląd przebiegu i błędów budowania | Region wyników kolumny dominującej | W toku · zakończone powodzeniem · zakończone błędem | 1 | Uruchomienie zadania budowania |
| Konsola | Region wykonania poleceń powłoki | Uruchamianie poleceń i podgląd strumieni wyjścia | Region wyników kolumny dominującej | Bezczynna · proces w toku · proces zakończony | 1 | Wejście w Terminal; polecenie w Chat Window |
| Git Panel | Panel kontroli wersji: gałęzie, zmiany, historia | Operacje kontroli wersji na repozytorium karty | Kolumna boczna, rozszerzenie boczne | Czysto · zmiany oczekujące · rozjazd z gałęzią zdalną · konflikt | 2 | Znacznik gałęzi w listwie stanu |
| Debugger | Panel punktów wstrzymania, stosu i zmiennych | Krokowe badanie wykonania programu | Kolumna boczna, rozszerzenie boczne | Nieaktywny · wstrzymany na punkcie · krokowanie · zakończony | 2 | Uruchomienie profilu debugowania |
| Zakładki regionu wyników | Przełączanie zawartości regionu wyników: budowanie, testy, problemy, konsola | Utrzymanie jednej kolumny wyników zamiast wielu równoległych | Zakładki regionu | Domyślna · aktywna · z plakietką liczby | 2 | Kliknięcie zakładki regionu |
| Uchwyty szerokości regionów | Pionowe uchwyty między regionami kolumny dominującej | Regulacja proporcji regionów | Uchwyt przeciągania krawędzi kolumny — **[DO DECYZJI OPERATORA]**, brak klasy uchwytu podziału w arkuszach źródłowych | Domyślny · przeciąganie | 2 | Przeciągnięcie krawędzi |
| Menu operacji pliku | Zestaw operacji na pozycji drzewa: otwórz, dodaj do koszyka, wyślij do modułu, zmień nazwę, usuń | Kierowanie pliku do właściwego narzędzia i operacje plikowe | Menu kontekstowe | Zwinięte · rozwinięte | 3 | Menu kontekstowe pozycji drzewa |

### 4.5. Karty sesji

Karta sesji odpowiada jednemu zadaniu inżynierskiemu. Pamięta repozytorium, gałąź, układ kolumn, koszyk kontekstu, historię Chat Window i stan pętli wykonawczej. Procesy karty trwają po stronie serwera niezależnie od tego, czy karta jest widoczna.

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|
| Karta sesji | Zakładka reprezentująca jedno zadanie inżynierskie | Przełączanie między równolegle prowadzonymi zadaniami | Zakładka (`.dn-zakladki`) | Domyślna · aktywna · z żywym procesem · przypięta · zamykana | 1 | Widoczna w pasmie kart sesji na wstążce |
| Nazwa własna karty | Tytuł zadania w miejsce nazwy modułu | Rozpoznanie zadania na pasku wielu kart | Tekst zakładki | Automatyczna · własna | 2 | Podwójne kliknięcie tytułu |
| Wskaźnik gałęzi na karcie | Znacznik gałęzi powiązanej z kartą | Rozróżnienie kart pracujących na różnych gałęziach | Znacznik (`.dn-etykietka`) | Domyślny · rozjazd z gałęzią zdalną | 1 | Widoczny na karcie |
| Wskaźnik pracy w tle | Plakietka aktywności procesów karty | Rozpoznanie kart z trwającym budowaniem, pętlą lub naradą | Plakietka (`.dn-kropka`) | Brak aktywności · praca w toku · wynik gotowy · błąd | 1 | Widoczny na karcie |
| Przycisk nowej karty | Kontrolka „+” otwarcia nowego zadania | Rozpoczęcie pracy nad kolejnym zadaniem | Ikona (`plus`) | Domyślny · wskazanie kursorem | 1 | Widoczny przy ostatniej karcie |
| Grupy kart | Kolorowe grupowanie kart wg projektu lub tematu | Porządkowanie wielu równoległych zadań | Grupa zakładek | Rozwinięta · zwinięta | 3 | Menu kontekstowe karty |
| Operacje na karcie | Zestaw: przypnij, duplikuj, rozdziel, zamknij, zostaw proces w tle | Zarządzanie cyklem życia karty i jej procesów | Menu kontekstowe | Zwinięte · rozwinięte | 3 | Menu kontekstowe karty |
| Wznawianie zamkniętych kart | Wykaz ostatnio zamkniętych kart z pełnym stanem | Powrót do zadania zamkniętego przez pomyłkę | Lista w menu przepełnienia | Domyślna · pusta | 3 | Menu ☰ paska kart; skrót klawiszowy |
| Zestaw kart jako obszar pracy | Zapisany zestaw otwartych kart wywoływany jednym poleceniem | Odtworzenie kompletu zadań projektu | Pozycja katalogu obszarów pracy | Zapisany · wczytany | 3 | Paleta poleceń; menu ☰ paska kart |
| Podgląd karty w tle | Miniatura stanu karty ze wskaźnikiem pracy | Ocena postępu bez przełączania karty | Panel popover | Zwinięty · rozwinięty | 2 | Wskazanie kursorem karty |

### 4.6. Okna pomocnicze powłoki

| Okno | Zawartość | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Koszyk kontekstu | Pliki, zaznaczenia, błędy, wyniki i fragmenty bieżącej karty; usuwanie i czyszczenie pozycji | Kolumna boczna, rozszerzenie boczne | Pusty · z pozycjami · przekazany do wywołania | 2 | Znacznik „Koszyk (N)” w Chat Window |
| Panel Problemy | Wspólny bufor błędów i ostrzeżeń z filtrami i skokiem do źródła | Kolumna boczna, rozszerzenie boczne | Bez problemów · z problemami · filtr aktywny | 1 | Segment diagnostyki listwy stanu |
| Menedżer procesów | Wykaz żywych procesów karty i środowiska z operacjami zatrzymania i podglądu | Kolumna boczna, rozszerzenie boczne | Brak procesów · procesy aktywne · proces zakończony błędem | 3 | Segment „Aktywne procesy” |
| Dziennik zdarzeń | Przewijany dziennik budowań, testów, zatwierdzeń, automatyk i powiadomień z filtrami | Kolumna boczna, rozszerzenie boczne | Domyślny · filtr aktywny · nowe wpisy | 3 | Segment powiadomień |
| Paleta poleceń | Pole z filtrem tekstowym po wszystkich poleceniach powłoki i modułów; tryby `>` polecenie, `@` symbol, `#` plik, `:` wiersz, `/` moduł | Nakładka nad obszarem roboczym | Puste · z wynikami · brak wyników | 4 | Skrót `Ctrl/Cmd + K` |
| Menedżer układów kolumn | Zapisane układy kolumn karty: „Kodowanie”, „Debugowanie”, „Wydanie”, „Narada”; zapis, wczytanie, przywrócenie układu bazowego | Panel popover | Domyślny · układ wczytany · układ zapisany | 3 | Menu `Układ ▼` |
| Szybka konfiguracja sesji | Model i kanał, wykonawca, izolacja karty, powiązane repozytorium, zestaw zmiennych, gęstość interfejsu, motyw | Panel popover | Domyślny · zmiana zastosowana | 3 | Znacznik kontekstu w belce tytułowej |
| Ściągawka skrótów | Wykaz skrótów aktywnych w bieżącym kontekście | Nakładka nad obszarem roboczym | Zwinięta · rozwinięta | 4 | Skrót klawiszowy |
| Panel prowenancji | Złożony kontekst przekazany modelowi w bieżącym wywołaniu | Okno platformowe | Domyślny · wywołanie wybrane | 4 | Polecenie w Chat Window; okno Konfiguracji |
| Mapa skrótów | Rejestr skrótów z edycją, wyszukiwaniem konfliktów i profilami | Okno platformowe | Domyślny · konflikt wykryty · profil zmieniony | 4 | Skrót klawiszowy; wyszukiwarka funkcji |

### 4.7. Always On Display i Mobile w środowisku

| Element | Rola w CodeStudio | Forma i waga | Stany | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|
| Awatar Always On Display | Agent towarzyszący świadomy kontekstu inżynierskiego: komentuje wynik budowania, przebieg testów i stan pętli wykonawczej | Element pływający nad obszarem roboczym (`.dn-awatar` w wariancie powiększonym) | Spoczynek · aktywny · powiadomienie · ukryty | 2 | Kliknięcie awatara; polecenie głosowe |
| Dymek kontekstowy | Podsumowanie stanu karty z ustaleniami i kontrolkami decyzji | Panel popover przy awatarze | Zwinięty · rozwinięty · oczekuje na decyzję | 2 | Kliknięcie awatara |
| Kanał Mobile | Podgląd procesów środowiska i interwencja zdalna: zatwierdzenie kroku, wstrzymanie pętli, wznowienie | Kanał zewnętrzny sparowanego urządzenia | Rozłączony · połączony · oczekuje na decyzję | 2 | Sparowane urządzenie; segment procesów listwy stanu |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — zlecenie programistyczne przez pętlę wykonawczą

```
 Użytkownik
     │  polecenie w języku naturalnym
     ▼
 CHAT WINDOW  (lewa kolumna, stała)
     │  przekazanie zlecenia wraz z koszykiem kontekstu
     ▼
 KOORDYNATOR
     │  dekompozycja zlecenia na zadania
     ▼
 EXECUTION LOOP WINDOW  (kolumna sąsiadująca)
     │  kolejka zadań · przydział wykonawcy
     ▼
 WYKONAWCA
     │  realizacja zadania w oknach modułu (Code Editor, Konsola, Build Output)
     ▼
 KOORDYNATOR — kontrola jakości wyniku
     │
     ├── wynik przyjęty ──► kolejne zadanie kolejki
     │                          │
     │                          └── kolejka pusta ──► wynik zbiorczy do Chat Window
     │                                                      │
     │                                                      ▼
     │                                              Użytkownik zatwierdza
     │
     └── uwagi kontroli ──► ponowienie zadania z korektą zakresu
                                │
                                └── przekroczony limit ponowień ──►
                                        wstrzymanie pętli, decyzja Użytkownika
                                        w Chat Window
```

### 5.2. Przepływ artefaktu magistralą kontekstu

Mechanizm magistrali opisuje [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md), rozdz. 6a.

```
 Build Output: błąd kompilacji
     │  operacja „Wyślij do Diagnostics” (menu kebab wyniku)
     ▼
 Wspólny bufor błędów środowiska
     │
     ├──► Panel Problemy — skok do źródła w Code Editorze
     ├──► Koszyk kontekstu — pozycja dostępna Wykonawcy
     └──► Execution Loop Window — nowe zadanie „Napraw błąd kompilacji”
                │
                ▼
          Wykonawca realizuje zadanie na pliku wskazanym przez błąd
```

### 5.3. Przepływ pracy z kartą sesji

```
 Nowa karta  ──►  wybór repozytorium i gałęzi
     │
     ▼
 Karta powiązana z repozytorium
     │  moduły karty pracują na wspólnym stanie
     ├── Developer: edycja
     ├── Terminal: wykonanie
     ├── Diagnostics: analiza
     └── Apps: budowanie i wdrożenie
     │
     ▼
 Zamknięcie karty
     │
     ├── proces żywy ──► proces pozostaje w tle, karta wznawialna z pełnym stanem
     └── brak procesów ──► karta trafia do wykazu ostatnio zamkniętych
```

### 5.4. Przepływ narady inżynierskiej

```
 Zagadnienie projektowe w Chat Window
     │  operacja „Wyślij do Roundtable”
     ▼
 Okno narady (kolumna dominująca) — kilka modeli odpowiada równolegle
     │
     ▼
 Consensus Panel — ustalenie wspólne
     │
     ▼
 Ustalenie trafia do koszyka kontekstu karty
     │
     ▼
 Chat Window: zlecenie implementacji ustalenia
     │
     ▼
 Execution Loop Window: dekompozycja i realizacja
```

### 5.5. Przepływ pracy ciągłej z modułem Automations

```
 Harmonogram Automations
     │  wyzwolenie zadania środowiska (budowanie nocne, przegląd zależności)
     ▼
 Koordynator — pętla wykonawcza w karcie powiązanej z repozytorium
     │
     ▼
 Wykonawca realizuje kolejkę zadań bez obecności Użytkownika
     │
     ├── wynik przyjęty ──► wpis w dzienniku zdarzeń, powiadomienie
     └── uwagi kontroli ──► Always On Display zgłasza decyzję
                                │
                                ▼
                          Mobile — zatwierdzenie, wstrzymanie
                          lub korekta zlecenia z dowolnego miejsca
```

---

## 6. Stany i powiązania

### 6.1. Stany przebiegu pętli wykonawczej

| Stan | Znaczenie | Sygnalizacja wizualna | Możliwe przejścia |
|---|---|---|---|
| Nieuruchomiona | Zlecenie przyjęte, dekompozycja gotowa, wykonanie niewystartowane | Plakietka neutralna „nieuruchomiona” | → uruchomiona |
| Uruchomiona | Wykonawca realizuje zadania kolejki pod nadzorem Koordynatora | Plakietka sukcesu, kropka pulsująca | → wstrzymana, → zakończona powodzeniem, → zakończona błędem |
| Wstrzymana | Przebieg zatrzymany, stan zadań zachowany | Plakietka ostrzeżenia „wstrzymana” | → uruchomiona, → przerwana |
| Oczekująca na decyzję | Kontrola jakości zgłosiła uwagę wymagającą rozstrzygnięcia przez Użytkownika | Plakietka informacyjna, ikona `oko` | → uruchomiona, → wstrzymana, → ponowienie zadania |
| Zakończona powodzeniem | Wszystkie zadania kolejki wykonane, wynik przyjęty kontrolą jakości | Plakietka sukcesu, ikona `ptaszek` | → nieuruchomiona (kolejne zlecenie) |
| Zakończona błędem | Zadanie nieudane po wyczerpaniu ponowień | Plakietka błędu | → nieuruchomiona po korekcie zlecenia |
| Przerwana | Przebieg zatrzymany poleceniem Użytkownika, bez wznowienia | Plakietka neutralna „przerwana” | → nieuruchomiona |

### 6.2. Stany zadania w kolejce pętli

| Stan | Znaczenie | Wejście w stan | Wyjście ze stanu |
|---|---|---|---|
| Oczekujące | Zadanie w kolejce, poprzedzające zadania niezakończone | Dekompozycja zlecenia | Przydział wykonawcy |
| W toku | Wykonawca realizuje zadanie | Przydział wykonawcy | Zwrot wyniku, przerwanie |
| W kontroli jakości | Wynik zwrócony, Koordynator ocenia zgodność ze zleceniem | Zwrot wyniku | Przyjęcie wyniku, zgłoszenie uwag |
| Ponawiane | Zadanie wraca do realizacji z korektą zakresu | Zgłoszenie uwag kontroli | Zwrot poprawionego wyniku |
| Zakończone | Wynik przyjęty | Przyjęcie wyniku kontrolą jakości | — |
| Niezdane | Wyczerpany limit ponowień | Przekroczenie limitu ponowień | Korekta zlecenia przez Użytkownika |
| Pominięte | Zadanie wyłączone z przebiegu poleceniem Użytkownika | Operacja „Pomiń” | Przywrócenie do kolejki |

### 6.3. Stany karty sesji

| Stan | Znaczenie | Sygnalizacja wizualna |
|---|---|---|
| Bezczynna | Karta otwarta, brak procesów i pętli w toku | Zakładka bez plakietki |
| Praca w toku | Trwa budowanie, test, proces terminala, narada albo pętla wykonawcza | Plakietka aktywności, kropka pulsująca |
| Wynik gotowy | Proces zakończony, wynik nieodczytany | Plakietka informacyjna |
| Błąd | Proces karty zakończony niepowodzeniem | Plakietka błędu |
| Proces w tle | Karta zamknięta, proces po stronie serwera trwa | Pozycja w wykazie procesów, wznawialna |

### 6.4. Stany repozytorium i kontroli wersji

| Stan | Znaczenie | Sygnalizacja w listwie stanu |
|---|---|---|
| Czysto | Brak zmian nieskomitowanych | `⑂ gałąź` bez licznika |
| Zmiany oczekujące | Pliki zmienione, niezatwierdzone | `⑂ gałąź · N zmian` |
| Rozjazd z gałęzią zdalną | Zatwierdzenia lokalne lub zdalne do zsynchronizowania | `⑂ gałąź · ↑N ↓M` |
| Konflikt scalania | Scalanie wymaga rozstrzygnięcia konfliktu | Segment w stanie błędu, Git Panel w stanie konfliktu |

### 6.5. Powiązania między oknami

| Zdarzenie źródłowe | Okno źródłowe | Skutek | Okno docelowe |
|---|---|---|---|
| Polecenie wieloetapowe | Chat Window | Otwarcie pętli wykonawczej i dekompozycja zlecenia | Execution Loop Window |
| Wynik zbiorczy pętli | Execution Loop Window | Podsumowanie i wezwanie do zatwierdzenia | Chat Window |
| Zaznaczenie kodu | Code Editor | Zaznaczenie dostępne jako kontekst wywołania | Chat Window, Execution Loop Window, Roundtable, Diagnostics |
| Błąd budowania | Build Output | Wpis do wspólnego bufora błędów | Panel Problemy, Koszyk kontekstu |
| Kliknięcie problemu | Panel Problemy | Otwarcie pliku i ustawienie kursora w miejscu błędu | Code Editor |
| Zmiana gałęzi | Git Panel | Rekonfiguracja kontekstu wszystkich modułów karty | Wszystkie okna karty |
| Zatwierdzenie zmiany | Chat Window | Zapis zmiany i aktualizacja stanu kontroli wersji | Git Panel, listwa stanu |
| Uruchomienie zadania | Menu `Zadania ▼` | Wykonanie polecenia i strumień wyjścia | Konsola, Build Output |
| Ustalenie narady | Consensus Panel | Wpis ustalenia do koszyka kontekstu | Koszyk kontekstu, Chat Window |
| Przydział wykonawcy | Execution Loop Window | Zmiana agenta realizującego zadanie kolejki | Katalog agentów |
| Decyzja zdalna | Mobile | Zatwierdzenie, wstrzymanie albo wznowienie pętli | Execution Loop Window |

---

## 7. Punkty sterowania z okna konfiguracji

Zestawienie jawnych kluczy sterujących środowiskiem CodeStudio. Każdy klucz ma zasięg i stan wyjściowy; żaden nie tworzy blokady — klucze sterują wyłącznie zachowaniem i wyglądem środowiska. Klucze zasięgu sesji dostępne są także z szybkiej konfiguracji sesji.

| Klucz | Czym steruje | Zasięg | Stan wyjściowy |
|---|---|---|---|
| `codestudio.nawigacja.grupy` | Skład i kolejność grup modułów w bocznej nawigacji | Środowisko | Grupy „Kod / Produkt / Współpraca” |
| `codestudio.nawigacja.przypięte` · `.kolejność` · `.widoczność` | Przypięcie, kolejność i ukrywanie pozycji modułów | Środowisko | Osiem modułów, kolejność źródłowa |
| `codestudio.nawigacja.tryb` | Panel wąski, szeroki albo zwinięty do krawędzi | Środowisko | Szeroki |
| `codestudio.chat.szerokosc` | Szerokość lewej kolumny Chat Window | Sesja / Moduł | Wartość bazowa układu |
| `codestudio.chat.kontekstAutomatyczny` | Automatyczne podawanie zaznaczenia kodu jako kontekstu wywołania | Sesja | Włączone |
| `codestudio.petla.otwieranie` | Warunek otwarcia Execution Loop Window: przy zleceniu wieloetapowym albo wyłącznie na żądanie | Sesja | Przy zleceniu wieloetapowym |
| `codestudio.petla.szerokosc` | Szerokość kolumny pętli wykonawczej | Sesja / Moduł | Wartość bazowa układu |
| `codestudio.petla.ponowienia` | Liczba ponowień zadania i warunek eskalacji do Użytkownika | Środowisko / Sesja | Trzy ponowienia, eskalacja do Chat Window |
| `codestudio.petla.kontrolaJakosci` | Zakres kontroli jakości wyniku zadania: zgodność ze zleceniem, testy, lintowanie, przegląd zmian | Środowisko | Zgodność ze zleceniem i testy |
| `codestudio.petla.wykonawcaDomyslny` | Agent lub model przypisywany zadaniom kolejki bez wskazania własnego | Środowisko / Sesja | Wykonawca sesji |
| `codestudio.uklad.kolumny` | Skład i kolejność kolumn obszaru roboczego | Sesja / Moduł | Chat Window · pętla · moduł · panele boczne |
| `codestudio.uklad.zapisane` | Katalog zapisanych układów kolumn i ich wywołanie | Środowisko | Zestaw bazowy: Kodowanie, Debugowanie, Wydanie, Narada |
| `codestudio.uklad.domyslny.<modul>` | Układ kolumn wczytywany przy wejściu w moduł | Moduł | Układ bazowy modułu |
| `codestudio.uklad.gestosc` | Gęstość interfejsu: kompaktowa albo komfortowa | Sesja | Komfortowa |
| `codestudio.karty.nazwa` | Tytuł karty: nazwa modułu albo nazwa zadania | Sesja | Nazwa zadania |
| `codestudio.karty.grupy` · `.przypięte` · `.przepelnienie` | Grupowanie, przypinanie i zachowanie paska przy wielu kartach | Środowisko | Bez grup, przewijanie |
| `codestudio.karty.repoPowiązanie` | Powiązanie karty z repozytorium i gałęzią | Sesja | Wg wybranego repozytorium |
| `codestudio.karty.procesTla` | Los procesu przy zamknięciu karty: pozostawienie w tle albo zatrzymanie | Sesja | Pozostawienie w tle |
| `codestudio.pasekStanu.segmenty` · `.widocznosc` | Skład, kolejność i widoczność segmentów listwy stanu | Środowisko | Zestaw inżynierski (rozdz. 3.3) |
| `codestudio.skroty.*` · `.profil` · `.klawiaturaModalna` | Mapa skrótów, profil przypisań, tryb klawiatury modalnej, akordy | Środowisko | Profil „Domyślny”, klawiatura modalna wyłączona |
| `codestudio.kontekst.koszyk` · `.przekazywanie` · `.reguly` | Koszyk kontekstu, kierowanie artefaktów do modułów, reguły przepływu | Sesja / Środowisko | Koszyk włączony, reguły wyłączone |
| `codestudio.kontekst.izolacja` | Izolacja kontekstu i zasobów karty | Sesja | **Wyłączona — pełny dostęp** |
| `codestudio.kontekst.zmienne` | Wspólne zmienne środowiskowe karty; klucze jawne, wartości maskowane | Sesja / Środowisko | Puste |
| `codestudio.wspolne.wyszukiwanie` · `.zadania` · `.dziennik` · `.powiadomienia` | Zakresy wyszukiwania, katalog zadań, retencja dziennika, kanały powiadomień | Środowisko | Włączone, zakresy bazowe |
| `codestudio.wyglad.motyw` | Motyw jasny albo ciemny | Sesja | Wg ustawień aplikacji |
| `codestudio.globalne.aod` · `.mobile` · `.skupienie` | Always On Display, kanał Mobile, tryb skupienia | Środowisko / Sesja | AOD i Mobile dostępne, skupienie wyłączone |

Powiązanie z oknami platformowymi: pełny podgląd złożonego kontekstu udostępnia Panel prowenancji okna Konfiguracji; konto, uwierzytelnianie, parowanie urządzeń i konta modeli prowadzi okno Ustawień.

---

## 8. Scenariusze użycia

### 8.1. Naprawa usterki zgłoszonej w repozytorium

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Użytkownik otwiera nową kartę sesji i wiąże ją z repozytorium oraz gałęzią zadania | Karta sesji, selektor repozytorium |
| 2 | Użytkownik opisuje usterkę i wskazuje pliki podejrzane; pliki trafiają do koszyka kontekstu | Chat Window, Project Tree, Koszyk kontekstu |
| 3 | Koordynator dokonuje dekompozycji zlecenia na analizę, poprawkę, testy i przegląd zmian | Execution Loop Window |
| 4 | Wykonawca realizuje zadania kolejki; zmiany widoczne na bieżąco w edytorze, testy w konsoli | Code Editor, Konsola |
| 5 | Kontrola jakości zgłasza uwagę do poprawki; zadanie wraca do ponowienia z korektą zakresu | Execution Loop Window |
| 6 | Wynik zbiorczy trafia do Chat Window; Użytkownik przegląda różnice i zatwierdza zmianę | Chat Window, Git Panel |

### 8.2. Diagnoza niepowodzenia budowania

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Budowanie kończy się błędem; segment statusu budowania przechodzi w stan błędu | Build Output, listwa stanu |
| 2 | Użytkownik kieruje błąd do Diagnostics operacją menu kontekstowego wyniku | Build Output, Panel Problemy |
| 3 | Analizator logów koreluje błąd ze zdarzeniami przebiegu i wskazuje miejsce w kodzie | Analizator logów, Code Editor |
| 4 | Użytkownik zleca poprawkę w Chat Window z błędem podanym jako kontekst | Chat Window, Koszyk kontekstu |
| 5 | Pętla wykonawcza realizuje poprawkę i uruchamia ponowne budowanie | Execution Loop Window, Build Output |

### 8.3. Decyzja projektowa wsparta naradą modeli

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Użytkownik formułuje zagadnienie architektoniczne i kieruje je do narady | Chat Window, Okno narady |
| 2 | Modele odpowiadają równolegle; moderator porządkuje wypowiedzi | Okno narady, Moderator Panel |
| 3 | Ustalenie wspólne trafia do koszyka kontekstu karty | Consensus Panel, Koszyk kontekstu |
| 4 | Użytkownik zleca implementację ustalenia | Chat Window |
| 5 | Koordynator dokonuje dekompozycji i prowadzi realizację | Execution Loop Window |

### 8.4. Wydanie produktu

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Użytkownik wczytuje zapisany układ kolumn „Wydanie” | Menedżer układów kolumn |
| 2 | Zadanie budowania produkcyjnego uruchamiane jest z menu zadań | Menu `Zadania ▼`, Build Output |
| 3 | Wynik testów i przegląd zmian potwierdzają gotowość | Panel testów, Git Panel |
| 4 | Panel wdrożenia prowadzi operację wydania; dziennik zdarzeń rejestruje przebieg | Panel wdrożenia, Dziennik zdarzeń |
| 5 | Użytkownik obserwuje przebieg wydania z urządzenia zdalnego | Mobile, Always On Display |

### 8.5. Praca nad dwoma zadaniami równolegle

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Dwie karty sesji wiążą się z dwiema gałęziami tego samego repozytorium | Karty sesji, selektor repozytorium |
| 2 | Pierwsza karta prowadzi pętlę wykonawczą nad zadaniem długim | Execution Loop Window |
| 3 | Użytkownik przełącza się na drugą kartę; procesy pierwszej trwają po stronie serwera | Karty sesji, listwa stanu |
| 4 | Wskaźnik pracy w tle na pierwszej karcie sygnalizuje gotowy wynik | Karta sesji |
| 5 | Powrót do pierwszej karty odtwarza układ kolumn, historię i stan pętli | Wszystkie okna karty |

### 8.6. Nocny przebieg zadań środowiska

| Krok | Działanie | Okna |
|---|---|---|
| 1 | Harmonogram Automations wyzwala zadanie przeglądu zależności w karcie projektu | Automations, Execution Loop Window |
| 2 | Koordynator prowadzi kolejkę zadań bez obecności Użytkownika | Execution Loop Window |
| 3 | Kontrola jakości zgłasza uwagę wymagającą decyzji | Execution Loop Window, Always On Display |
| 4 | Użytkownik rozstrzyga uwagę z urządzenia zdalnego | Mobile |
| 5 | Przebieg kończy się wpisem w dzienniku zdarzeń i powiadomieniem | Dziennik zdarzeń |

---

## 9. Wykazy normatywne i kryteria odbioru

### 9.1. Komendy kontraktu

Środowisko CodeStudio korzysta z komend trzech obszarów kontraktu: `developer` (50 komend), `terminal` (27 komend), `diagnostics` (9 komend) — razem 86 komend. Wejście do środowiska i wybór modułu korzystają dodatkowo z dwóch komend obszaru `environment` (`environment.list`, `environment.enter` — rozdz. 2.4). Pełny wykaz nazw komend, pól żądania i pól wyniku dla wszystkich 86 komend znajduje się w [Załączniku B — pełnym wykazie komend kontraktu środowiska CodeStudio](#załącznik-b-pełny-wykaz-komend-kontraktu-środowiska-codestudio), na końcu niniejszego dokumentu.

| Obszar | Liczba komend | Reprezentatywne komendy | Zastosowanie w środowisku |
|---|---|---|---|
| `developer` | 50 | `developer.file.open`, `developer.git.action`, `developer.build.run`, `developer.refactor.apply` | Code Editor, Project Tree, Git Panel, Build Output (rozdz. 4.4) |
| `terminal` | 27 | `terminal.session.open`, `terminal.command.exec`, `terminal.process.list`, `terminal.process.kill` | Konsola, Menedżer procesów (rozdz. 3.2) |
| `diagnostics` | 9 | `diagnostics.error.list`, `diagnostics.error.update`, `diagnostics.recommendation.update` | Konsola i wykaz błędów, Diagnostics (rozdz. 2.1) |

### 9.2. Żetony projektowe

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Plakietka „42/42 testy przechodzą” | `--dn-sukces-tekst`, `--dn-sukces-tlo`, `--dn-sukces-obrys` | Kolor stanu powodzenia budowania i testów w listwie stanu (rozdz. 3.3) |
| Plakietka „ostrzeżenie: DeprecationWarning” | `--dn-ostrzezenie-tekst`, `--dn-ostrzezenie-tlo`, `--dn-ostrzezenie-obrys` | Kolor ostrzeżenia w wyniku pytest i w Diagnostics |
| Plakietka błędu budowania | `--dn-blad-tekst`, `--dn-blad-tlo`, `--dn-blad-obrys` | Kolor niepowodzenia budowania, zliczenie błędów w listwie stanu |
| Kropka sygnału procesu aktywnego | `--dn-kropka`, `--dn-czas-tetno` | Wskaźnik „pracuje” przy karcie sesji i przy zadaniu w toku |
| Fokus kontrolek | `--dn-fokus` | Obrys fokusu klawiaturowego pól, przycisków, pozycji list |
| Przycisk główny akcji (`.dn-btn--sygnal`) | `--dn-sygnal`, `--dn-sygnal-wypelnienie`, `--dn-sygnal-wypelnienie-hover` | „Zatwierdź zmiany”, „Wypchnij (push)”, akcje pierwszoplanowe |
| Odstępy paneli i formularzy | `--dn-od-2`, `--dn-od-4`, `--dn-od-6` | Odstęp wewnętrzny kart, wierszy listy zmian, pól formularza commitu |
| Promień zaokrąglenia kart i pól | `--dn-r-sm`, `--dn-r-md`, `--dn-r-xl` | Przyciski i pola (`--dn-r-sm`), karty zadań (`--dn-r-md`), karta środowiska w przedsionku (`--dn-r-xl`) |
| Typografia bloku kodu | `--dn-ff-mono`, `--dn-fs-sm` | Code Editor, blok diff, wynik `pytest` w terminalu |
| Typografia interfejsu | `--dn-ff-bazowa`, `--dn-fs-base`, `--dn-lh-bazowy` | Tekst kontrolek, etykiet i list w gęstości zwartej |
| Czas przejścia warstw | `--dn-czas-3` | Otwarcie/zamknięcie panelu bocznego, palety poleceń, dymka funkcji |

### 9.3. Komponenty interfejsu

| Klasa `.dn-*` | Rola | Modyfikatory użyte w środowisku | Stany |
|---|---|---|---|
| `.dn-btn` | Przycisk | `--sygnal` (akcja główna: zatwierdź, wypchnij), `--zarys` (akcja pomocnicza: cofnij, anuluj) | domyślny · wskazanie kursorem · wciśnięty · nieaktywny · ładowanie |
| `.dn-pole` | Wrapper pola formularza | — | domyślny · błąd |
| `.dn-pole-kontrolka` | Kontrolka pola (tekst, wybór, wieloliniowa) | — | puste · wypełnione · fokus · błąd |
| `.dn-karta` | Karta | — | domyślna · klikalna · wybrana |
| `.dn-kolejka` | Wiersz zadania w kolejce/planie | — | oczekuje · w toku · zakończone · błąd |
| `.dn-krok` | Krok planu dekompozycji zadania | — | oczekuje · w toku · zaakceptowany · odrzucony |
| `.dn-kropka` | Wskaźnik semantyczny stanu | — | aktywny (tętno) · nieaktywny |
| `.dn-etykietka` | Etykieta filtru/plakietka drobna | — | domyślna · aktywna |
| `.dn-awatar` | Awatar Operatora/Always On Display | — | domyślny |
| `.dn-szyna-poz` | Pozycja szyny nawigacji modułów | `--modul` | domyślna · aktywna · wskazana |
| `.dn-wpis` | Wpis w strumieniu Chat Window / Execution Loop Window | — | domyślny · strumieniowanie |
| `.dn-zakladki` | Pas kart sesji | — | aktywna · w tle · z pracą w toku |

### 9.4. Etykiety interfejsu z prototypu

Wykaz dosłownych brzmień z `design/05-okna/srodowiska/codestudio.html` i `codestudio-przedsionek.html`, nieopisanych już w rozdziałach 1–8.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Panele modułu Developer | „Code Editor”, „Project Tree”, „Git Panel”, „Build Output” | Rozdz. 4.4 |
| Panele uniwersalne | „Plan”, „Zadania w tle”, „Kolejka”, „Subagenci”, „Artefakty” | Pasek boczny okien roboczych |
| Akcje karty sesji | „Zmień nazwę”, „Rozgałęź”, „Archiwizuj” | Menu kebab karty sesji |
| Tryb uprawnień | „Ręczny”, „Auto — domyślny — zero blokad” | Panel ustawień Execution Loop Window (rozdz. 7) |
| Akcja zatwierdzenia zmian | „Zatwierdź zmiany”, „Z pominięciem zgód”, „Przyjmij zmiany” | Git Panel |
| Pole opisu commitu | „Opis commitu”, „Wygeneruj przez AI” | Git Panel |
| Akcje repozytorium | „Zatwierdź (commit)”, „Wypchnij (push)” | Git Panel |
| Nagłówek środowiska w przedsionku | „Środowisko pracy”, „CODESTUDIO”, „Projektuj. Buduj. Rozwijaj.” | Rozdz. 2.4 |
| Listwa działań przedsionka | „Nowe repozytorium”, „Konfiguracja środowiska”, „Ustawienia” | Rozdz. 2.4 |
| Szyna sesji przedsionka | „Twoje sesje”, filtry „Czynne” · „Wszystkie” · „Zakończone” | Rozdz. 2.4 |

### 9.5. Komunikaty

| Sytuacja | Dosłowna treść | Rodzaj |
|---|---|---|
| Budowanie zakończone powodzeniem | „42 testy przechodzą, pokrycie 91%” | Potwierdzenie |
| Wynik uruchomienia pytest w terminalu | „42 passed, 1 warning in 3.14s — pokrycie 91%” | Potwierdzenie z ostrzeżeniem |
| Ostrzeżenie wykryte w budowaniu | „ostrzeżenie: DeprecationWarning w silnik.py:88” | Ostrzeżenie |
| Status budowania w listwie stanu | „build ✔ 41 s” albo „build ✖ (N błędów)” | Stan |
| Wynik testów w listwie stanu | „testy 128/130” | Stan |
| Rozpoczęcie kolekcji testów | „# kolekcjonowanie przypadków…” | Stan pośredni |
| Brak przypisanego repozytorium | „Repozytorium: Brak” | Stan pusty |
| Zmiany oczekujące na zatwierdzenie | „Zmiany do zatwierdzenia (2)”, „+34 −8” | Stan |
| Kliknięcie przycisku bez uprawnienia w trybie obserwatora | Komunikat kontekstowy ze wskazaniem przełączenia w tryb operatora — przycisk pozostaje klikalny (zasada zero blokad) | Komunikat kontekstowy |

### 9.6. Punkty łamania

Wartości progów z `design/zasoby/zetony/zetony.css`, rozdział „Siatka i punkty łamania”.

| Żeton | Próg szerokości | Zachowanie układu w środowisku CodeStudio |
|---|---|---|
| `--dn-bp-w1` | 640px | Widok mobilny — środowisko udostępniane wyłącznie poprzez funkcję globalną Mobile ([Mobile](../funkcje-globalne/mobile.md)), nie jako stanowisko stacjonarne w pełnym układzie |
| `--dn-bp-w2` | 960px | Boczna nawigacja modułów (szyna) zwija się do samych ikon; szyna sesji przedsionka zwija się analogicznie |
| `--dn-bp-w3` | 1280px | Pełny kokpit — Chat Window, Execution Loop Window i okna robocze modułu widoczne jednocześnie w układzie docelowym |
| `--dn-bp-w4` | 1600px | Szerokie biurko — para Chat Window / Execution Loop Window (Koordynator ↔ Wykonawca) rozwija się do pełnej szerokości roboczej bez nakładania się na panele modułu |

Poniżej `--dn-bp-w2` okna robocze modułu Developer i panele uniwersalne (Plan, Kolejka, Subagenci, Artefakty) przechodzą z układu kolumnowego na zakładkowy — wybór między nimi odbywa się `.dn-zakladki`, nie jednoczesnym wyświetlaniem.

### 9.7. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Wejście do środowiska otwiera przedsionek, nie przestrzeń roboczą | Kliknięcie karty CodeStudio na stronie głównej; obserwacja, że żaden moduł nie jest otwarty automatycznie |
| Wszystkie 86 komend obszarów `developer`, `terminal`, `diagnostics` mają pokrycie w interfejsie | Zestawienie wykazu rozdz. 9.1 z panelami Code Editor, Project Tree, Git Panel, Build Output, Konsola, Menedżer procesów, Panel Problemy i Analizator logów |
| Żaden przycisk nie jest trwale nieaktywny | Przegląd wszystkich kontrolek wykazanych w rozdz. 9.3 pod kątem stanu `disabled`; zero wystąpień poza stanem „ładowanie” |
| Chat Window i Execution Loop Window są obecne w każdym module środowiska | Otwarcie kolejno wszystkich ośmiu modułów (rozdz. 2.1); potwierdzenie stałej obecności obu kanałów |
| Powrót z przestrzeni roboczej do przedsionka nie zamyka kart sesji | Powrót do przedsionka, sprawdzenie szyny sesji — karty czynne pozostają widoczne jako sesje w toku |
| Układ odpowiada progom łamania rozdz. 9.6 na czterech szerokościach referencyjnych | Zmiana szerokości okna kolejno przez 640px, 960px, 1280px, 1600px; obserwacja zachowań z tabeli rozdz. 9.6 |
| Kolory stanu budowania i testów pochodzą wyłącznie z żetonów rozdz. 9.2 | Przegląd arkusza stylów zastosowanego do listwy stanu i plakietek — zero wartości szesnastkowych wpisanych na sztywno |

---

## Załącznik A. Skróty klawiszowe środowiska

Skróty należą do warstwy 4 — stanowią ekspercką drogę dostępu do funkcji osiągalnych także kliknięciem albo poleceniem w Chat Window. Pełny rejestr z edycją przypisań i wykrywaniem konfliktów udostępnia Mapa skrótów okna Konfiguracji.

### A.1. Sterowanie kanałami komunikacji

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + L` | Ustawienie kursora w polu poleceń Chat Window |
| `Ctrl/Cmd + Shift + L` | Otwarcie i zamknięcie Execution Loop Window |
| `Ctrl/Cmd + Enter` | Zatwierdzenie wyniku przedstawionego przez Wykonawcę |
| `Esc` | Przerwanie bieżącego działania Wykonawcy |
| `Ctrl/Cmd + Shift + Space` | Wstrzymanie i wznowienie pętli wykonawczej |
| `Ctrl/Cmd + Shift + R` | Ponowienie bieżącego zadania kolejki |
| `Ctrl/Cmd + Shift + K` | Otwarcie koszyka kontekstu |

### A.2. Nawigacja

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + K` | Paleta poleceń |
| `Ctrl/Cmd + P` | Przełącznik uniwersalny „Idź do…” |
| `Ctrl/Cmd + 1`–`8` | Przejście do modułu wg kolejności bocznej nawigacji |
| `Ctrl/Cmd + Tab` | Przełączenie na poprzednio używany moduł |
| `Ctrl/Cmd + B` | Zwinięcie i rozwinięcie bocznej nawigacji modułów |
| `Ctrl/Cmd + Shift + F` | Wyszukiwanie w repozytorium karty |

### A.3. Karty sesji i układ kolumn

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + T` | Nowa karta sesji |
| `Ctrl/Cmd + W` | Zamknięcie karty sesji |
| `Ctrl/Cmd + Shift + T` | Przywrócenie ostatnio zamkniętej karty |
| `Ctrl/Cmd + Alt + →` / `←` | Przejście do sąsiedniej karty |
| `Ctrl/Cmd + \` | Podział kolumny dominującej na regiony |
| `Ctrl/Cmd + Shift + \` | Powiększenie bieżącego okna do pełni kolumny dominującej |
| `Ctrl/Cmd + Alt + U` | Menedżer układów kolumn |
| `Ctrl/Cmd + Alt + 0` | Przywrócenie układu bazowego karty |

### A.4. Praca inżynierska

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + Shift + B` | Uruchomienie budowania |
| `Ctrl/Cmd + Shift + U` | Uruchomienie testów |
| `Ctrl/Cmd + Shift + M` | Panel Problemy |
| `Ctrl/Cmd + Shift + G` | Git Panel |
| `Ctrl/Cmd + ~` | Konsola modułu Terminal |
| `F5` | Uruchomienie profilu debugowania |
| `F9` | Ustawienie punktu wstrzymania w Code Editorze |

### A.5. Funkcje eksperckie

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + Alt + K` | Mapa skrótów |
| `Ctrl/Cmd + Alt + ?` | Ściągawka skrótów bieżącego kontekstu |
| `Ctrl/Cmd + Alt + P` | Panel prowenancji |
| `Ctrl/Cmd + Alt + I` | Ustawienia izolacji kontekstu karty |
| `Ctrl/Cmd + Alt + D` | Dziennik zdarzeń środowiska |
| `Ctrl/Cmd + Alt + A` | Always On Display |


---

## Załącznik B. Pełny wykaz komend kontraktu środowiska CodeStudio

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `developer` — 50 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `developer.file.open` | Wczytuje plik repozytorium do Code Editor | `windowId:string` (wym)<br>`path:string` (wym) | `file:DeveloperFile` (wym) |
| `developer.file.save` | Zapisuje treść pliku repozytorium | `windowId:string` (wym)<br>`path:string` (wym)<br>`content:string` (wym)<br>`createVersion:bool` (opc) | `file:DeveloperFile` (wym) |
| `developer.tree.get` | Zwraca drzewo projektu odzwierciedlające stan repozytorium | `windowId:string` (wym)<br>`path:string` (opc)<br>`depth:int` (opc)<br>`includeHidden:bool` (opc) | `root:string` (wym)<br>`nodes:DeveloperTreeNode[]` (wym) |
| `developer.git.action` | Wykonuje czynność na repozytorium powiązanym z sesją | `windowId:string` (wym)<br>`action:GitActionKind` (wym)<br>`paths:string[]` (opc)<br>`message:string` (opc)<br>`branch:string` (opc)<br>`remote:string` (opc)<br>`force:bool` (opc) | `result:GitActionResult` (wym) |
| `developer.build.run` | Uruchamia albo przerywa budowanie; log idzie na żywo do Build Output | `windowId:string` (wym)<br>`task:string` (wym)<br>`arguments:string[]` (opc)<br>`stop:bool` (opc) | `build:DeveloperBuild` (wym) |
| `developer.file.version.list` | Zwraca wersje pliku założone przy zapisie z createVersion | `windowId:string` (wym)<br>`path:string` (wym)<br>`limit:int` (opc) | `versions:DeveloperFileVersion[]` (wym)<br>`total:int` (opc) |
| `developer.file.version.restore` | Przywraca treść pliku z zapisanej wersji | `windowId:string` (wym)<br>`versionId:string` (wym)<br>`writeToDisk:bool` (opc) | `file:DeveloperFile` (wym) |
| `developer.file.create` | Zakłada plik albo katalog w katalogu roboczym okna | `windowId:string` (wym)<br>`path:string` (wym)<br>`kind:TreeNodeKind` (wym)<br>`content:string` (opc) | `node:DeveloperTreeNode` (wym) |
| `developer.file.rename` | Zmienia nazwę pliku albo katalogu bez zmiany jego miejsca | `windowId:string` (wym)<br>`path:string` (wym)<br>`newName:string` (wym) | `node:DeveloperTreeNode` (wym) |
| `developer.file.delete` | Usuwa wskazane węzły drzewa projektu | `windowId:string` (wym)<br>`paths:string[]` (wym)<br>`recursive:bool` (opc) | `deletedPaths:string[]` (wym)<br>`undoToken:string` (opc) |
| `developer.file.move` | Przenosi węzły drzewa do wskazanego katalogu | `windowId:string` (wym)<br>`paths:string[]` (wym)<br>`targetPath:string` (wym) | `nodes:DeveloperTreeNode[]` (wym) |
| `developer.symbol.navigate` | Zwraca miejsca symbolu wskazane przez serwer języka | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`column:int` (wym)<br>`kind:SymbolNavigationKind` (wym) | `symbols:DeveloperSymbol[]` (wym)<br>`serverAvailable:bool` (wym) |
| `developer.format.run` | Formatuje dokument albo zaznaczenie wedle reguł stylu repozytorium | `windowId:string` (wym)<br>`path:string` (wym)<br>`content:string` (opc)<br>`startLine:int` (opc)<br>`endLine:int` (opc)<br>`writeToDisk:bool` (opc) | `file:DeveloperFile` (wym)<br>`changed:bool` (wym)<br>`formatter:string` (opc) |
| `developer.lint.get` | Zwraca zgłoszenia analizy statycznej dla wskazanych plików | `windowId:string` (wym)<br>`paths:string[]` (opc)<br>`limit:int` (opc) | `diagnostics:DeveloperDiagnostic[]` (wym)<br>`linterAvailable:bool` (wym)<br>`truncated:bool` (opc) |
| `developer.refactor.apply` | Wykonuje refaktoryzację semantyczną obejmującą wiele plików | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`column:int` (wym)<br>`kind:RefactorKind` (wym)<br>`newName:string` (opc)<br>`preview:bool` (opc) | `edits:DeveloperTextEdit[]` (wym)<br>`applied:bool` (wym)<br>`changedPaths:string[]` (opc) |
| `developer.grep.search` | Szuka wzorca w całym repozytorium katalogu roboczego okna | `windowId:string` (wym)<br>`pattern:string` (wym)<br>`regex:bool` (opc)<br>`caseSensitive:bool` (opc)<br>`include:string[]` (opc)<br>`exclude:string[]` (opc)<br>`limit:int` (opc) | `matches:DeveloperGrepMatch[]` (wym)<br>`total:int` (opc)<br>`truncated:bool` (wym) |
| `developer.grep.replace` | Zamienia trafienia wzorca w repozytorium; warstwa funkcji eksperckich | `windowId:string` (wym)<br>`pattern:string` (wym)<br>`replacement:string` (wym)<br>`regex:bool` (opc)<br>`paths:string[]` (opc)<br>`preview:bool` (opc) | `edits:DeveloperTextEdit[]` (wym)<br>`changedPaths:string[]` (wym)<br>`applied:bool` (wym) |
| `developer.git.status` | Oddaje stan repozytorium katalogu roboczego okna bez wykonywania czynności | `windowId:string` (wym) | `status:DeveloperGitStatus` (wym) |
| `developer.git.diff` | Zwraca różnice między wskazanymi stanami repozytorium | `windowId:string` (wym)<br>`paths:string[]` (opc)<br>`staged:bool` (opc)<br>`fromRef:string` (opc)<br>`toRef:string` (opc)<br>`contextLines:int` (opc) | `hunks:GitDiffHunk[]` (wym)<br>`binaryPaths:string[]` (opc) |
| `developer.git.log` | Zwraca historię zatwierdzeń repozytorium | `windowId:string` (wym)<br>`branch:string` (opc)<br>`path:string` (opc)<br>`author:string` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`limit:int` (opc) | `commits:GitCommit[]` (wym)<br>`total:int` (opc) |
| `developer.git.branch.list` | Zwraca gałęzie repozytorium wraz z rozbieżnością wobec gałęzi zdalnej | `windowId:string` (wym)<br>`includeRemote:bool` (opc) | `branches:GitBranch[]` (wym)<br>`current:string` (opc) |
| `developer.git.conflict.get` | Oddaje trzy wersje pliku skonfliktowanego do widoku trójstronnego | `windowId:string` (wym)<br>`path:string` (wym) | `conflict:GitConflict` (wym) |
| `developer.git.conflict.resolve` | Oznacza konflikt jako rozstrzygnięty wskazaną wersją albo treścią własną | `windowId:string` (wym)<br>`path:string` (wym)<br>`resolution:ConflictResolutionKind` (wym)<br>`content:string` (opc) | `resolved:bool` (wym)<br>`remainingPaths:string[]` (wym) |
| `developer.build.list` | Zwraca przebiegi budowania z dziennika okna | `windowId:string` (wym)<br>`status:BuildStatus` (opc)<br>`limit:int` (opc) | `builds:DeveloperBuild[]` (wym)<br>`total:int` (opc) |
| `developer.build.log.get` | Rozwija odnośnik logu przebiegu budowania | `buildId:string` (wym)<br>`tail:int` (opc)<br>`fromLine:int` (opc) | `lines:string[]` (wym)<br>`truncated:bool` (wym)<br>`truncatedFrom:int` (opc) |
| `developer.test.result.get` | Zwraca wynik zestawu testów przebiegu budowania | `buildId:string` (wym)<br>`status:TestStatus` (opc) | `results:DeveloperTestResult[]` (wym)<br>`passed:int` (wym)<br>`failed:int` (wym)<br>`skipped:int` (wym)<br>`durationMs:int64` (opc) |
| `developer.coverage.get` | Zwraca pokrycie kodu testami dla przebiegu budowania | `buildId:string` (wym)<br>`path:string` (opc) | `files:DeveloperCoverage[]` (wym)<br>`percent:int` (wym)<br>`threshold:int` (opc) |
| `developer.debug.session.start` | Rozpoczyna sesję debugowania pod adapterem właściwym językowi | `windowId:string` (wym)<br>`configurationId:string` (opc)<br>`program:string` (opc)<br>`arguments:string[]` (opc)<br>`adapter:string` (opc)<br>`stopOnEntry:bool` (opc) | `session:DebugSession` (wym) |
| `developer.debug.session.control` | Steruje przebiegiem sesji debugowania | `sessionId:string` (wym)<br>`step:DebugStepKind` (wym)<br>`threadId:string` (opc) | `session:DebugSession` (wym) |
| `developer.breakpoint.set` | Zakłada albo zdejmuje punkt przerwania na marginesie Code Editora | `windowId:string` (wym)<br>`path:string` (wym)<br>`line:int` (wym)<br>`kind:BreakpointKind` (opc)<br>`condition:string` (opc)<br>`hitCondition:string` (opc)<br>`logMessage:string` (opc)<br>`remove:bool` (opc) | `breakpoints:Breakpoint[]` (wym) |
| `developer.debug.scope.get` | Zwraca stos wywołań, zakresy i zmienne zatrzymanego procesu | `sessionId:string` (wym)<br>`frameId:string` (opc)<br>`variablesRef:string` (opc) | `frames:DebugFrame[]` (wym)<br>`scopes:DebugScope[]` (wym)<br>`variables:DebugVariable[]` (wym) |
| `developer.debug.evaluate` | Wykonuje wyrażenie w kontekście zatrzymanej ramki | `sessionId:string` (wym)<br>`frameId:string` (wym)<br>`expression:string` (wym)<br>`assignTo:string` (opc) | `value:string` (wym)<br>`type:string` (opc)<br>`variablesRef:string` (opc) |
| `developer.api.request` | Wykonuje zapytanie HTTP w sieci okna modułu | `windowId:string` (wym)<br>`method:string` (wym)<br>`url:string` (wym)<br>`headers:json` (opc)<br>`body:string` (opc)<br>`bodyKind:string` (opc)<br>`environmentId:string` (opc)<br>`timeoutMs:int` (opc) | `response:ApiResponse` (wym) |
| `developer.api.collection.save` | Zapisuje kolekcję zapytań wraz ze środowiskami | `windowId:string` (wym)<br>`collectionId:string` (opc)<br>`name:string` (wym)<br>`requests:json` (wym)<br>`environments:json` (opc) | `collection:ApiCollection` (wym) |
| `developer.api.collection.list` | Zwraca kolekcje zapytań okna | `windowId:string` (wym)<br>`collectionId:string` (opc) | `collections:ApiCollection[]` (wym) |
| `developer.api.openapi.import` | Buduje kolekcję zapytań z kontraktu OpenAPI | `windowId:string` (wym)<br>`path:string` (opc)<br>`url:string` (opc)<br>`collectionName:string` (opc) | `collection:ApiCollection` (wym)<br>`requestCount:int` (wym) |
| `developer.data.connection.set` | Zakłada albo zmienia definicję połączenia bazodanowego | `windowId:string` (wym)<br>`connectionId:string` (opc)<br>`name:string` (wym)<br>`engine:DataEngine` (wym)<br>`host:string` (opc)<br>`port:int` (opc)<br>`database:string` (wym)<br>`user:string` (opc)<br>`credentialRef:string` (opc)<br>`readOnly:bool` (opc) | `connection:DataConnection` (wym) |
| `developer.data.connection.list` | Zwraca połączenia bazodanowe okna | `windowId:string` (wym) | `connections:DataConnection[]` (wym) |
| `developer.data.schema.get` | Zwraca drzewo schematu połączenia w liście płaskiej | `connectionId:string` (wym)<br>`path:string` (opc)<br>`depth:int` (opc) | `nodes:DataSchemaNode[]` (wym) |
| `developer.data.query.run` | Wykonuje zapytanie na połączeniu albo oddaje jego plan | `connectionId:string` (wym)<br>`sql:string` (wym)<br>`limit:int` (opc)<br>`explain:bool` (opc)<br>`transaction:bool` (opc) | `result:DataQueryResult` (wym) |
| `developer.data.migration.run` | Uruchamia migracje schematu z katalogu roboczego okna | `connectionId:string` (wym)<br>`windowId:string` (wym)<br>`direction:string` (wym)<br>`target:string` (opc)<br>`dryRun:bool` (opc) | `applied:string[]` (wym)<br>`pending:string[]` (wym)<br>`output:string` (opc) |
| `developer.container.list` | Zwraca kontenery i obrazy silnika kontenerów | `windowId:string` (wym)<br>`all:bool` (opc)<br>`includeImages:bool` (opc) | `containers:ContainerInfo[]` (wym)<br>`images:ImageInfo[]` (opc)<br>`engineAvailable:bool` (wym) |
| `developer.container.action` | Wykonuje czynność cyklu życia kontenera | `containerId:string` (wym)<br>`action:ContainerActionKind` (wym)<br>`tail:int` (opc) | `container:ContainerInfo` (wym)<br>`output:string` (opc) |
| `developer.image.build` | Buduje obraz kontenera z pliku Dockerfile repozytorium | `windowId:string` (wym)<br>`dockerfile:string` (wym)<br>`tag:string` (wym)<br>`buildArgs:json` (opc)<br>`push:bool` (opc)<br>`registry:string` (opc) | `imageId:string` (wym)<br>`logRef:string` (opc) |
| `developer.compose.up` | Podnosi albo zatrzymuje stos usług opisany plikiem compose | `windowId:string` (wym)<br>`file:string` (wym)<br>`services:string[]` (opc)<br>`down:bool` (opc) | `services:ContainerInfo[]` (wym)<br>`output:string` (opc) |
| `developer.dependency.list` | Zwraca drzewo zależności z manifestu repozytorium | `windowId:string` (wym)<br>`manifest:string` (opc)<br>`outdatedOnly:bool` (opc)<br>`depth:int` (opc) | `dependencies:DependencyNode[]` (wym)<br>`manifest:string` (wym) |
| `developer.scan.run` | Uruchamia skan zależności, sekretów, kodu albo licencji | `windowId:string` (wym)<br>`kinds:ScanKind[]` (wym)<br>`paths:string[]` (opc) | `scan:ScanRun` (wym) |
| `developer.scan.result.list` | Zwraca zgłoszenia skanu bezpieczeństwa | `scanId:string` (opc)<br>`windowId:string` (opc)<br>`kind:ScanKind` (opc)<br>`severity:ProblemSeverity` (opc)<br>`limit:int` (opc) | `findings:ScanFinding[]` (wym)<br>`total:int` (opc) |
| `developer.contextual.op` | Wykonuje operację kontekstową modelu nad treścią repozytorium | `windowId:string` (wym)<br>`operation:ContextualOpKind` (wym)<br>`path:string` (opc)<br>`selection:string` (opc)<br>`contextPaths:string[]` (opc)<br>`instruction:string` (opc)<br>`targetLanguage:string` (opc) | `result:string` (wym)<br>`edits:DeveloperTextEdit[]` (opc)<br>`messageId:string` (opc) |
| `developer.toolchain.check` | Sprawdza obecność programów zewnętrznych w ścieżce wykonywalnej rdzenia | `programs:string[]` (wym) | `programs:ToolchainProgram[]` (wym) |

**Zdarzenia obszaru `developer` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `developer.build.changed` | Przyrost budowania; log narasta w toku przebiegu | — |
| `developer.debug.changed` | Zmiana stanu sesji debugowania: zatrzymanie, wznowienie, zakończenie | — |

### Obszar `terminal` — 27 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `terminal.session.open` | Otwiera kartę terminala jako odrębną sesję powłoki | `windowId:string` (wym)<br>`shell:TerminalShell` (wym)<br>`workingDir:string` (opc)<br>`title:string` (opc)<br>`environment:EnvironmentVariable[]` (opc)<br>`remoteTarget:string` (opc)<br>`remotePort:int` (opc)<br>`hostId:string` (opc)<br>`containerRef:TerminalContainerRef` (opc)<br>`serialDevice:string` (opc)<br>`serialBaudRate:int` (opc) | `session:TerminalSession` (wym) |
| `terminal.command.exec` | Uruchamia polecenie w karcie terminala; wynik idzie strumieniem fragmentów | `sessionId:string` (wym)<br>`command:string` (wym)<br>`initiator:ProcessInitiator` (opc)<br>`timeoutMs:int` (opc) | `process:TerminalProcess` (wym) |
| `terminal.process.list` | Zwraca procesy rejestru rdzenia, w tym uruchomione przez model | `windowId:string` (opc)<br>`sessionId:string` (opc)<br>`status:TerminalProcessStatus` (opc)<br>`initiator:ProcessInitiator` (opc) | `processes:TerminalProcess[]` (wym) |
| `terminal.process.kill` | Kończy proces sygnałem łagodnym albo wymuszonym | `processId:string` (wym)<br>`force:bool` (opc) | `process:TerminalProcess` (wym) |
| `terminal.output.stream` | Zapisuje okno na zbiorcze wyjście wszystkich otwartych kart terminala i zwraca ogon historii | `sessionId:string` (opc)<br>`windowId:string` (opc)<br>`terminalSessionIds:string[]` (opc)<br>`tail:int` (opc) | `lines:TerminalOutputLine[]` (wym)<br>`subscribed:bool` (wym) |
| `terminal.output.read` | Oddaje wyjście jednego procesu terminala — stdout, stderr, kod wyjścia i stan — jako jedną odpowiedź, bez zapisywania się na strumień. Komenda `terminal.command.exec` kończy się w chwili startu procesu (kompilacja trwa dłużej niż każde sensowne oczekiwanie na odpowiedź, a rozłączenie klienta nie ma prawa jej przerwać), więc wyniku nieść nie może i nigdy nie będzie mogła. Tędy model czyta, co polecenie wypisało. Źródłem jest ten sam dziennik zbiorczego wyjścia, na którym stoi `terminal.output.stream` — drugiego strumienia nie ma. Historia żyje jeden bieg rdzenia: po ponownym uruchomieniu wyjście jest puste, a odpowiedź mówi to wprost | `processId:string` (wym)<br>`tail:int` (opc)<br>`waitMs:int` (opc) | `stdout:string` (wym)<br>`stderr:string` (wym)<br>`exitCode:int` (opc)<br>`status:TerminalProcessStatus` (wym)<br>`truncated:bool` (wym)<br>`truncatedBytes:int` (opc) |
| `terminal.session.close` | Zamyka kartę powłoki. Dziś zamknięcie karty żyje wyłącznie w widoku klienta: powłoka i jej procesy biegną dalej, a rdzeń o zamknięciu nie wie | `sessionId:string` (wym)<br>`force:bool` (opc) | `session:TerminalSession` (wym)<br>`stoppedProcessIds:string[]` (opc) |
| `terminal.session.list` | Zwraca karty powłoki znane rdzeniowi. Rdzeń odtwarza karty przy starcie, ale klient po ponownym połączeniu nie ma jak ich zobaczyć i zaczyna wykaz od pustego | `windowId:string` (opc)<br>`status:TerminalSessionStatus` (opc)<br>`includeExited:bool` (opc) | `sessions:TerminalSession[]` (wym)<br>`total:int` (wym) |
| `terminal.file.read` | Oddaje treść pliku z katalogu roboczego karty. Bez tej komendy klient czyta manifest projektu poleceniem powłoki, więc wykrycie zadań zależy od programu wypisującego plik i od składni każdej z sześciu powłok | `sessionId:string` (wym)<br>`path:string` (wym)<br>`tail:int` (opc)<br>`maxBytes:int` (opc) | `content:string` (wym)<br>`path:string` (wym)<br>`truncated:bool` (wym)<br>`truncatedBytes:int64` (opc)<br>`sizeBytes:int64` (opc) |
| `terminal.process.suspend` | Wstrzymuje albo wznawia proces rejestru rdzenia. Dziś rdzeń umie proces wyłącznie zakończyć, więc długie zadanie da się tylko przerwać | `processId:string` (wym)<br>`resume:bool` (opc) | `process:TerminalProcess` (wym)<br>`supported:bool` (wym) |
| `terminal.host.save` | Zapisuje wpis książki hostów albo podmienia istniejący. Bez tej komendy książka żyje jedno posiedzenie przeglądarki i ginie przy odświeżeniu strony | `host:TerminalHost` (wym) | `host:TerminalHost` (wym)<br>`created:bool` (wym) |
| `terminal.host.list` | Zwraca książkę hostów Operatora | `group:string` (opc)<br>`query:string` (opc) | `hosts:TerminalHost[]` (wym)<br>`total:int` (wym) |
| `terminal.host.remove` | Usuwa wpis książki hostów. Karty już otwarte do tego hosta biegną dalej | `hostId:string` (wym) | `removed:bool` (wym) |
| `terminal.script.save` | Zapisuje skrypt albo wycinek kodu biblioteki jako kolejną wersję. Bez tej komendy biblioteka żyje jedno posiedzenie, a jedyną drogą jej zachowania jest wywóz do pliku | `script:TerminalScript` (wym) | `script:TerminalScript` (wym)<br>`created:bool` (wym) |
| `terminal.script.list` | Zwraca pozycje biblioteki skryptów | `kind:TerminalScriptKind` (opc)<br>`tag:string` (opc)<br>`query:string` (opc) | `scripts:TerminalScript[]` (wym)<br>`total:int` (wym) |
| `terminal.script.remove` | Usuwa pozycję biblioteki wraz ze wszystkimi jej wersjami | `scriptId:string` (wym) | `removed:bool` (wym) |
| `terminal.script.lint` | Poddaje treść skryptu analizie statycznej i formatowaniu. Analizę prowadzą programy spoza instalacji Danaco Console, więc odpowiedź mówi wprost, czy narzędzie było dostępne — pusty wykaz uwag przy braku narzędzia znaczyłby fałszywie treść bez zastrzeżeń | `content:string` (wym)<br>`shell:TerminalShell` (wym)<br>`format:bool` (opc) | `findings:TerminalLintFinding[]` (wym)<br>`formatted:string` (opc)<br>`analyzerAvailable:bool` (wym)<br>`analyzer:string` (wym) |
| `terminal.tunnel.open` | Zakłada przekierowanie portu. Dziś tunel da się założyć wyłącznie poleceniem wydanym w karcie, a wtedy jego stan i przepustowość są niewidoczne | `windowId:string` (wym)<br>`kind:TerminalTunnelKind` (wym)<br>`hostId:string` (opc)<br>`remoteTarget:string` (opc)<br>`localPort:int` (opc)<br>`remoteHost:string` (opc)<br>`remotePort:int` (opc) | `tunnel:TerminalTunnel` (wym) |
| `terminal.tunnel.list` | Zwraca przekierowania portów wraz z ich stanem | `windowId:string` (opc)<br>`status:TerminalTunnelStatus` (opc) | `tunnels:TerminalTunnel[]` (wym)<br>`total:int` (wym) |
| `terminal.tunnel.close` | Zamyka przekierowanie portu | `tunnelId:string` (wym) | `tunnel:TerminalTunnel` (wym) |
| `terminal.key.generate` | Wytwarza parę kluczy SSH na maszynie rdzenia. Hasło klucza wchodzi odwołaniem do sejfu, nigdy treścią — zgodnie z zasadą zapisaną w kontrakcie przy zmiennej środowiska | `name:string` (wym)<br>`keyType:TerminalKeyType` (wym)<br>`comment:string` (opc)<br>`passphraseRef:string` (opc) | `key:TerminalSshKey` (wym) |
| `terminal.key.import` | Wciąga do wykazu klucz leżący już na maszynie rdzenia. Klucz wskazuje się ścieżką, a nie treścią: materiał kluczowy nie ma powodu przechodzić przez łącze | `name:string` (wym)<br>`path:string` (wym) | `key:TerminalSshKey` (wym) |
| `terminal.key.list` | Zwraca wykaz kluczy SSH znanych rdzeniowi wraz z ich odciskami | — | `keys:TerminalSshKey[]` (wym)<br>`total:int` (wym) |
| `terminal.key.remove` | Zdejmuje klucz z wykazu. Wpisy książki hostów wskazujące ten klucz tracą wskazanie i wracają do klucza domyślnego konfiguracji maszyny | `keyId:string` (wym)<br>`deleteFiles:bool` (opc) | `removed:bool` (wym)<br>`detachedHostIds:string[]` (opc) |
| `terminal.watch.start` | Zakłada obserwację plików uruchamiającą polecenie przy ich zmianie. Kontrakt daje dziś wyzwalacz plikowy automatyce, a nie karcie powłoki | `sessionId:string` (wym)<br>`pattern:string` (wym)<br>`command:string` (wym)<br>`debounceMs:int` (opc)<br>`recursive:bool` (opc) | `watch:TerminalWatch` (wym) |
| `terminal.watch.stop` | Zatrzymuje obserwację plików. Polecenie już uruchomione biegnie dalej | `watchId:string` (wym) | `watch:TerminalWatch` (wym) |
| `terminal.watch.list` | Zwraca obserwacje plików wraz z licznikiem wyzwoleń | `windowId:string` (opc)<br>`status:TerminalWatchStatus` (opc) | `watches:TerminalWatch[]` (wym)<br>`total:int` (wym) |

**Zdarzenia obszaru `terminal` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `terminal.process.changed` | Zmiana procesu rejestru rdzenia serwera | — |

### Obszar `diagnostics` — 9 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `diagnostics.log.query` | Przeszukuje dziennik z filtrowaniem po poziomie, źródle i zakresie czasu | `pattern:string` (opc)<br>`regex:bool` (opc)<br>`level:LogLevel` (opc)<br>`source:string` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`deduplicate:bool` (opc)<br>`limit:int` (opc)<br>`processId:string` (opc)<br>`traceId:string` (opc) | `entries:LogEntry[]` (wym)<br>`total:int` (opc)<br>`truncated:bool` (opc) |
| `diagnostics.error.list` | Zwraca błędy zgrupowane po odcisku wraz z kontekstem wystąpienia | `status:DiagnosticErrorStatus` (opc)<br>`priority:DiagnosticPriority` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`limit:int` (opc)<br>`source:string` (opc)<br>`errorCode:string` (opc) | `errors:DiagnosticError[]` (wym)<br>`total:int` (opc) |
| `diagnostics.analyze.run` | Uruchamia analizę zagregowanego stanu systemu | `windowId:string` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`errorIds:string[]` (opc)<br>`compareAnalysisId:string` (opc) | `analysis:DiagnosticAnalysis` (wym) |
| `diagnostics.recommendation.list` | Zwraca rekomendacje poprawek powstałe z analizy | `analysisId:string` (opc)<br>`status:RecommendationStatus` (opc)<br>`priority:DiagnosticPriority` (opc)<br>`limit:int` (opc) | `recommendations:DiagnosticRecommendation[]` (wym) |
| `diagnostics.error.update` | Zapisuje stan, priorytet i notatkę błędu diagnostycznego | `errorId:string` (wym)<br>`status:DiagnosticErrorStatus` (opc)<br>`priority:DiagnosticPriority` (opc)<br>`note:string` (opc) | `error:DiagnosticError` (wym) |
| `diagnostics.recommendation.update` | Zapisuje stan rekomendacji naprawczej wraz z komentarzem Operatora | `recommendationId:string` (wym)<br>`status:RecommendationStatus` (wym)<br>`comment:string` (opc) | `recommendation:DiagnosticRecommendation` (wym) |
| `diagnostics.problem.report` | Wnosi błąd albo ostrzeżenie do wspólnego bufora problemów środowiska; piszą tu budowanie, testy i kontrola stylu, a czytają `diagnostics.error.list` oraz `diagnostics.log.query` | `source:string` (wym)<br>`level:LogLevel` (wym)<br>`message:string` (wym)<br>`filePath:string` (opc)<br>`line:int` (opc)<br>`column:int` (opc)<br>`errorCode:string` (opc) | `error:DiagnosticError` (wym) |
| `diagnostics.metrics.query` | Zwraca szeregi czasowe miar platformy wraz z percentylami p50, p95 i p99 oraz znacznikami wdrożeń naniesionymi na tę samą oś | `metrics:MetricKind[]` (opc)<br>`fromTime:int64` (opc)<br>`toTime:int64` (opc)<br>`buckets:int` (opc) | `series:MetricSeries[]` (wym)<br>`deployments:DeploymentMarker[]` (opc)<br>`sampledSince:int64` (opc) |
| `diagnostics.profile.capture` | Zdejmuje profil procesu platformy i oddaje go drzewem ramek do wykresu płomieniowego | `kind:ProfileKind` (wym)<br>`depth:int` (opc) | `root:ProfileFrame` (wym)<br>`kind:ProfileKind` (wym)<br>`unit:string` (wym)<br>`capturedAt:int64` (wym) |

**Zdarzenia obszaru `diagnostics` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `diagnostics.analysis.changed` | Zmiana analizy diagnostycznej; odświeża wykaz rekomendacji zwracany komendą `diagnostics.recommendation.list` | — |


Razem w wykazie: **86 komend** z 3 obszarów kontraktu.

---

*Koniec dokumentu. Środowisko CodeStudio — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
