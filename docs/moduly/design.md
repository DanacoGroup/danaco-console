# Danaco Console — Moduł Design — dokument projektowy

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
| **Moduł** | Design (rozdz. 11.9 Koncepcji platformy, 4.9 Specyfikacji modułów) |
| **Środowiska dostępności** | WorkSpace, CodeStudio (moduł niedostępny w TalkIn) |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji; nie tworzy komponentu własnego |
| **Odbiorcy dokumentu** | Designer (co, gdzie, w jakiej formie, do czego), Operator (co personalizuje), Deweloper (co zbudować, na czym) |
| **Źródło faktów** | architektura/koncepcja-platformy.md (rozdz. 8, 9.9), specyfikacje/specyfikacja-modulow.md (4.9), specyfikacje/specyfikacja-okien-operacyjnych.md (6.9), interfejs-uzytkownika/system-wizualny.md, architektura/model-konfiguracji.md |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Komplet okien operacyjnych modułu — przegląd](#2-komplet-okien-operacyjnych-modułu--przegląd)
3. [Specyfikacja okien operacyjnych — pełny arsenał narzędzi](#3-specyfikacja-okien-operacyjnych--pełny-arsenał-narzędzi)
4. [Katalog funkcji i narzędzi](#4-katalog-funkcji-i-narzędzi)
5. [Katalog elementów interfejsu](#5-katalog-elementów-interfejsu)
6. [Przepływy pracy w module](#6-przepływy-pracy-w-module)
7. [Punkty sterowania z okna konfiguracji](#7-punkty-sterowania-z-okna-konfiguracji)
8. [Stany, dane i powiązania z innymi modułami](#8-stany-dane-i-powiązania-z-innymi-modułami)
9. [Scenariusze użycia](#9-scenariusze-użycia)
10. [Zgodność z rdzeniem platformy](#10-zgodność-z-rdzeniem-platformy)

---

## 1. Przeznaczenie i kontekst

### 1.1. Rola modułu

Design jest **kompletnym warsztatem wizualnym platformy** — pojedynczym miejscem, w którym powstaje, jest edytowany, porządkowany i wydawany każdy zasób graficzny wykorzystywany w pozostałych modułach. Moduł pokrywa osiem obszarów tematycznych:

| Obszar | Przedmiot pracy |
|---|---|
| A. Generowanie AI | Tworzenie grafiki z opisu, wariacje, przekształcenia neuronowe |
| B. Edycja rastrowa | Retusz, kadrowanie, warstwy, maski, korekcja barwna, filtry |
| C. Grafika wektorowa | Rysunek ścieżkowy, kształty, operacje logiczne, wektoryzacja |
| D. Projektowanie UI i makiety | Ramki, komponenty, auto-layout, prototypowanie, wireframe |
| E. Tokeny i system projektowy | Kolory, typografia, odstępy, motywy, eksport tokenów do kodu |
| F. Kolor | Palety, kontrast, harmonie, symulacja wad wzroku, ekstrakcja |
| G. Ikony i typografia | Zestawy ikon, favicony, sprite'y, dobór par krojów, podgląd webfontów |
| H. Eksport, podgląd i handoff | Formaty wyjściowe, kompresja, skalowanie @1/2/3x, inspekcja, kod |

Moduł obejmuje ponadto **warstwy przekrojowe**: bibliotekę zasobów (Assets Panel), wersjonowanie i współpracę (kursory, komentarze), animację lekką (Lottie, GIF, APNG) oraz import zasobów zewnętrznych.

### 1.2. Dla kogo

Design jest przestrzenią dla użytkowników **generujących i edytujących materiały graficzne** — od pojedynczych ilustracji i elementów brandingowych po kompletne makiety interfejsu. Adresatem są projektanci, zespoły marketingowe oraz zespoły produktowe budujące jednocześnie w module Apps (9.14) warstwę frontendową wymagającą zasobów wizualnych.

### 1.3. Po co

Moduł łączy w jednym miejscu **generowanie grafiki przez AI, jej edycję i zestawianie w spójną kompozycję** — z bezpośrednim wsparciem AI na każdym etapie procesu twórczego, zamiast przełączania się między osobnym narzędziem generującym a osobnym narzędziem porządkującym wynik.

### 1.4. Charakter pracy

| Cecha | Opis |
|---|---|
| Punkt wyjścia | Zamiar wizualny wyrażony poleceniem, a nie gotowy plik do edycji |
| Rytm pracy | Cykl: sformułowanie polecenia → wygenerowanie → ocena → zestawienie lub odrzucenie → iteracja |
| Wynik | Zasób wizualny gromadzony w Assets Panel, gotowy do dalszego wykorzystania w innych modułach |
| Dostępność środowiskowa | WorkSpace (praca projektowa i produktowa) oraz CodeStudio (budowa interfejsu aplikacji) — moduł nieobecny w TalkIn, gdzie centralnym przedmiotem pracy jest treść, nie grafika |

### 1.5. Granica tematyczna

Granica tematyczna jest jawna i utrzymuje czytelny podział między modułami platformy:

| Poza zakresem Design | Uzasadnienie | Moduł właściwy |
|---|---|---|
| Redagowanie treści dokumentów (tekst długi, PDF/DOCX, korekta, streszczenia) | Przedmiotem Design jest grafika, nie proza | **Studio** (4.1) |
| Budowa działającego frontendu produktu (kod komponentów, routing, stan) | Design wytwarza zasoby i makiety, nie kod aplikacji | **Apps — Frontend Workspace** (4.14) |
| Trwałe repozytorium plików całej organizacji | Assets Panel jest podręczną biblioteką modułu, nie centralnym archiwum | **Library** (4.6) |
| Montaż i renderowanie wideo (timeline, ścieżki, klatki kluczowe wideo) | Design obejmuje animację lekką (Lottie, GIF), nie postprodukcję filmu | poza platformą |
| Modelowanie i rendering 3D (sceny, siatki, materiały PBR) | Zakres modułu to 2D oraz lekki mockup 3D urządzeń | poza platformą |
| Przygotowalnia poligraficzna (CMYK profilowany, spady, pasery) | Moduł wydaje PDF i konwersję CMYK, nie pełny prepress | częściowo, Grupa H |

Design **wydaje** zasoby do Studio i Apps przez jawne powiązania konfigurowalne — nie przejmuje ich zadań.

### 1.6. Miejsce modułu w architekturze platformy

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
| Tokens & System Panel | Tokeny, motywy, kontrast, przewodnik stylu | Kolumna boczna, otwierana jako rozszerzenie boczne, drzewo tokenów i tabele | 3 | Menu kebab (⋮) obszaru roboczego, polecenie w Chat Window | Nie |

```
 Makieta całościowa — Moduł: Design          Dostępność: WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window            │ Design Board                │ Kolumny
  nawigacja │ Użytkownik ↔ Wykonawca │  kanwa robocza kompozycji   │ boczne
  modułów   │                        │                             │ (rozsze-
            │ ────────────────────   │  [Danaco Console][Design]   │  rzenia
            │ Execution Loop Window  │  [Ultra] [Tryb: kanwa ▼]  ⋮ │  boczne)
            │ Koordynator ↔ Wykonawca│                             │
            │  zlecenie · zadania    │                             │
 ═══════════════════════════════════════════════════════════════════════════
```

### 2.1. Warstwy widoczności w module

Moduł stosuje regułę stopniowego ujawniania funkcjonalności: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna.

| Warstwa | Zastosowanie w module Design |
|---|---|
| 1 | Chat Window, Execution Loop Window, kanwa Design Board, pasek kontekstu ze znacznikami projektu, środowiska, modelu i wykonawcy, wskaźnik stanu generowania |
| 2 | Assets Panel, Prompt Builder, Preview Window, wybór silnika generującego, tryb generowania, tryb narzędzia kanwy, poziom powiększenia — wywoływane znacznikiem kontekstowym lub przyciskiem, zwijane samoczynnie po użyciu |
| 3 | Tokens & System Panel, operacje logiczne na ścieżkach, wyrównanie i rozmieszczenie, warianty eksportu, tryb wsadowy, zestaw akcji zasobu — menu kebab (⋮), menu hamburger (☰), menu kontekstowe kanwy, panele popover |
| 4 | ControlNet i mapy sterujące, rozdzielenie generacji na warstwy, symulacja wad wzroku, inspektor glifów, eksport kodu widoku, konfiguracja endpointów obliczeniowych — polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny |

Każda ukryta funkcja modułu osiągalna jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wydanym w Chat Window. Makiety w rozdz. 3 przedstawiają stan spoczynku interfejsu: elementy warstwy 1 oraz zwinięte wyzwalacze warstw wyższych.

---

## 3. Specyfikacja okien operacyjnych — pełny arsenał narzędzi

### 3.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

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
| Ocena i iteracja konwersacyjna | Polecenia modyfikujące ostatnio wygenerowany zasób („jaśniejsza wersja", „usuń tło", „powiększ margines") | 2 | Pole polecenia przy przypiętej miniaturze |
| Skrócona ścieżka generowania | Uruchomienie generowania wprost z okna komunikacji, bez otwierania Prompt Buildera | 2 | Polecenie generujące w polu wprowadzania |
| Akcje na wiadomości | Kopiuj polecenie, zapisz jako szablon w Prompt Builderze, powtórz z modyfikacją | 3 | Menu kebab (⋮) wiadomości |
| Wybór wykonawcy i modelu | Zmiana modelu generującego i wykonawcy dla bieżącej sesji | 2 | Znacznik kontekstowy `[Ultra]`, `[Fable 5]` |
| Polecenia funkcji eksperckich | Uruchomienie operacji warstwy 4 (rozdzielenie na warstwy, mapy sterujące, eksport kodu widoku) | 4 | Polecenie języka naturalnego |

**Zachowanie i stany:** okno dostrojone do kontekstu bieżącej kompozycji na Design Board; historia domyślnie odrębna per karta sesji. Stan „generowanie w toku" — wskaźnik postępu przy miniaturze oczekiwanego wyniku. Przerwanie zlecenia wydane w Chat Window zatrzymuje pętlę wykonawczą prezentowaną w Execution Loop Window.

```
 Makieta — Chat Window w module Design (stan spoczynku)
 ═════════════════════════════════════
  Użytkownik ↔ Wykonawca            ⋮
  [Danaco Console] [Design] [Ultra]
 ─────────────────────────────────────
  Historia: „Ilustracja bohatera
  strony — wariant złoty"
  [miniatura ostatniego zasobu ▤]
 ─────────────────────────────────────
  ┌─────────────────────────────────┐
  │ Polecenie…                      │
  └─────────────────────────────────┘
   @  📎                  [ Wyślij ⏎ ]
 ═════════════════════════════════════
```

### 3.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Pole | Treść |
|---|---|
| Cel | Prowadzenie pętli wykonawczej modułu: dekompozycja zlecenia graficznego na zadania, koordynacja, nadzór nad realizacją i kontrola jakości wyniku |
| Zawartość | Bieżące zlecenie i jego dekompozycja na zadania modułu (generowanie, przekształcenie, kadrowanie, eksport, wydanie), kolejka i stan zadań, wymiana komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Bieżące zlecenie i dekompozycja | Prezentacja zlecenia graficznego rozłożonego na zadania modułu | 1 | Widoczne bez interakcji |
| Kolejka i stan zadań | Lista zadań pętli (prompt, generowanie, upscaling, usunięcie tła, wektoryzacja, eksport, wydanie) wraz ze stanem każdego | 1 | Widoczna bez interakcji |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator ↔ Wykonawca dla zadań tego modułu | 1 | Widoczna bez interakcji |
| Wskaźniki przebiegu pętli | Liczba iteracji, czas zadania, obciążenie kanału modelu, koszt zlecenia | 1 | Widoczne bez interakcji |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | 1 | Widoczne bez interakcji |
| Wyniki kontroli jakości | Ocena wyniku wobec kryteriów zlecenia (proporcje, paleta, kontrast, format) i decyzja o ponowieniu | 2 | Znacznik `Kontrola ▼` przy zadaniu |
| Szczegóły zadania | Parametry wykonania: model, seed, endpoint, wejściowe zasoby, czas i koszt | 3 | Menu kebab (⋮) zadania |
| Ponowienie z modyfikacją | Powtórzenie zadania ze zmienionym parametrem bez powrotu do Prompt Buildera | 3 | Menu kontekstowe zadania |
| Dziennik pętli | Pełny zapis przebiegu pętli wykonawczej wraz z komunikatami diagnostycznymi | 4 | Polecenie języka naturalnego, tryb administracyjny |

**Zachowanie i stany:** okno prezentuje pętlę wykonawczą zleceń modułu Design. Stan „pętla bezczynna" — okno zwinięte do znacznika stanu przy Chat Window. Stan „pętla w biegu" — kolumna rozwinięta z aktywną kolejką zadań. Stan „oczekiwanie na decyzję" — zadanie wstrzymane do rozstrzygnięcia w Preview Window, przy czym pozostałe zadania kolejki biegną dalej.

```
 Makieta — Execution Loop Window (stan spoczynku)
 ═════════════════════════════════════
  Koordynator ↔ Wykonawca           ⋮
 ─────────────────────────────────────
  Zlecenie: „Zestaw ilustracji
  brandingowych — kampania Q3"
 ─────────────────────────────────────
  ▸ 1. Kompozycja polecenia   gotowe
  ▸ 2. Generowanie 4 wariantów  w toku
  ▸ 3. Usunięcie tła          w kolejce
  ▸ 4. Eksport @1x/@2x        w kolejce
 ─────────────────────────────────────
  Iteracja 2 · kanał modelu 41%
  [ Wstrzymaj ] [ Przerwij ]  Kontrola ▼
 ═════════════════════════════════════
```

### 3.3. Design Board

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
| Inspektor właściwości | Wymiary, wypełnienie, obrys, efekty, więzy responsywne, auto-layout | 2 | Zaznaczenie elementu |
| Panel komponentów | Komponenty i warianty stanu, biblioteka UI systemu wizualnego | 2 | Znacznik `Komponenty ▼` |
| Siatka i linie pomocnicze | Wsparcie precyzyjnego rozmieszczenia elementów, z przyciąganiem | 2 | Znacznik `Siatka ▼` |
| Szablony układu | Schematy zestawień: tablica nastroju, siatka porównawcza, arkusz brandingowy, formaty społecznościowe | 3 | Menu hamburger (☰) kanwy |
| Narzędzia wyrównania i rozmieszczenia | Wyrównanie do siatki, rozłożenie równomierne, grupowanie i rozgrupowanie | 3 | Pasek kontekstowy zaznaczenia |
| Operacje logiczne na ścieżkach | Suma, różnica, przecięcie, wykluczenie | 3 | Menu kontekstowe zaznaczenia |
| Adnotacje i komentarze | Przypinane uwagi do elementów kompozycji, wątki, oznaczenia osób | 3 | Menu kebab (⋮), skrót klawiszowy |
| Wersjonowanie kompozycji | Historia stanów tablicy, nazwane wersje, powrót, porównanie | 3 | Menu kebab (⋮) tablicy |
| Eksport kompozycji i fragmentu | Zapis tablicy lub wskazanego obszaru jako plik graficzny lub PDF | 3 | Menu kebab (⋮), pasek kontekstowy zaznaczenia |
| Zaznaczenie wielokrotne | Operacje zbiorcze na wielu elementach (przesunięcie, usunięcie, eksport) | 2 | Zaznaczenie ramką lub z klawiszem modyfikującym |
| Kursor współpracy | Widoczność pozycji innych osób pracujących nad tą samą tablicą | 1 | Widoczny przy sesji współdzielonej |
| Biblioteka elementów pomocniczych | Kształty, linie prowadzące, ramki mockupów i elementy szkicowe | 3 | Menu hamburger (☰) kanwy |
| Tryb podglądu tokenów | Przełączenie kanwy między motywem jasnym i ciemnym systemu tokenów | 3 | Menu kebab (⋮) kanwy |
| Rozdzielenie generacji na warstwy | Rozkład wygenerowanego obrazu na edytowalne obiekty | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** odbiera zasoby wygenerowane i zgromadzone w Assets Panel przez przeciągnięcie na kanwę. Stan pusty — zachęta do wygenerowania pierwszego zasobu poleceniem w Chat Window. Stan „edycja współdzielona" — widoczne kursory innych uczestników.

```
 Makieta — Design Board (stan spoczynku)
 ═══════════════════════════════════════════════════════════
  Kompozycja: „Strona główna — wariant A"                  ⋮
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

### 3.4. Assets Panel

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
| Wyszukiwanie | Po nazwie, tagu lub treści polecenia, które wygenerowało zasób | 2 | Pole wyszukiwania |
| Filtrowanie wg typu | Ilustracja, element brandingowy, makieta interfejsu, ikona, zdjęcie | 2 | Znacznik `Filtr ▼` |
| Filtr po pochodzeniu | AI, edytowane, importowane, wektor, token | 2 | Znacznik `Filtr ▼` |
| Tagi i kolekcje | Nadawanie etykiet i grupowanie zasobów tematycznie | 2 | Plakietka tagu, pole tagów zasobu |
| Oznaczenie ulubionych | Szybki dostęp do najczęściej wykorzystywanych zasobów | 2 | Ikona gwiazdki na miniaturze |
| Panel metadanych i prowenancji | Polecenie źródłowe, model, seed, endpoint, rozdzielczość, data, łańcuch edycji | 3 | Menu kebab (⋮) zasobu |
| Historia wariantów | Kolejne warianty tego samego zasobu z porównaniem | 3 | Menu kebab (⋮) zasobu |
| Pobieranie i eksport zbiorczy | Zaznaczenie wielu zasobów i wydanie w wybranych formatach i skalach | 3 | Menu kebab (⋮) zaznaczenia |
| Import zasobów zewnętrznych | Wyszukanie i wstawienie zasobów z bibliotek zewnętrznych oraz treści próbnych | 3 | Menu hamburger (☰) panelu |
| Przeciągnięcie na Design Board | Przeniesienie zasobu wprost do kompozycji roboczej | 1 | Przeciągnięcie miniatury |
| Wysłanie do modułu docelowego | Przekazanie zasobu do Studio Editor lub Frontend Workspace | 2 | Przycisk `→ Moduł` |
| Usuwanie i archiwizacja | Porządkowanie zbioru zasobów — usunięcie następuje od razu, z dostępnym „Cofnij"; potwierdzenie przy usunięciu jest ustawieniem konfiguracyjnym | 3 | Menu kebab (⋮) zasobu |

**Zachowanie i stany:** zasoby trafiają tu bezpośrednio po zakończeniu zadania generowania w pętli wykonawczej. Stan „nowy" — plakietka przy zasobach dodanych od ostatniej wizyty w panelu.

```
 Makieta — Assets Panel (stan spoczynku)
 ═══════════════════════════════
  Zasoby (48)                  ☰
  [Szukaj]        [Filtr ▼]  ⋮
 ───────────────────────────────
  🖼 bohater_v3 ★
  🖼 logo_zloto
  🖼 ikona_produktu   [nowy]
 ───────────────────────────────
  [ → Design Board ]  [ → Moduł ]
 ═══════════════════════════════
```

### 3.5. Prompt Builder

| Pole | Treść |
|---|---|
| Cel | Precyzyjne formułowanie poleceń generujących grafikę |
| Zawartość | Struktura promptu: tryb generowania, styl, kompozycja, parametry |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

**Pełny arsenał narzędzi:**

| Narzędzie | Funkcja | Warstwa | Sposób wywołania |
|---|---|---|---|
| Pola strukturalne promptu | Temat, styl, kompozycja, oświetlenie, paleta barw, proporcje kadru, wykluczenia | 1 | Widoczne po otwarciu panelu |
| Przycisk „⚡ Generuj" | Uruchomienie generowania na podstawie skompletowanego promptu | 1 | Widoczny po otwarciu panelu |
| Tryb generowania | Text-to-image, image-to-image, inpainting, outpainting, ikony, makieta z opisu | 2 | Znacznik `Tryb ▼` |
| Wybór silnika generującego | Przełącznik modelu generującego grafikę i endpointu obliczeniowego | 2 | Znacznik kontekstowy modelu |
| Suwaki parametrów generowania | Kreatywność wobec promptu, ziarno losowości, liczba wariantów, siła wpływu referencji | 2 | Znacznik `Parametry ▼` |
| Obraz referencyjny | Wgranie obrazu jako punktu odniesienia generowania | 2 | Ikona załącznika referencji |
| Biblioteka stylów | Presety stylistyczne do zastosowania i dalszej modyfikacji | 3 | Lista rozwijana pola stylu |
| Paleta z Tokens & System Panel | Wstrzyknięcie zestawu tokenów kolorów jako ograniczenia barwnego generacji | 3 | Menu kebab (⋮) pola palety |
| Generowanie wsadowe | Utworzenie kilku wariantów jednym poleceniem, prezentowanych do wyboru | 2 | Pole `Warianty ▼` |
| Historia poleceń | Lista wcześniej użytych promptów z ponownym użyciem | 3 | Przycisk `Historia` |
| Zapis jako szablon | Utrwalenie skonstruowanego promptu do wielokrotnego użycia | 3 | Menu `Szablony ▼` |
| Wersjonowanie i porównanie promptu | Zestawienie dwóch wersji polecenia obok siebie wraz z ich wynikami | 3 | Menu kebab (⋮) panelu |
| Licznik długości polecenia | Orientacyjna liczba znaków i tokenów promptu | 2 | Widoczny przy polu tematu po rozwinięciu parametrów |
| Panel map sterujących | Wymuszenie pozy, głębi lub krawędzi w generacji na podstawie mapy sterującej | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** zlecenie generowania trafia do pętli wykonawczej, a jego przebieg widoczny jest w Execution Loop Window; wynik trafia do Assets Panel po zakończeniu zadania. Przycisk „⚡ Generuj" jest zawsze aktywny — przy pustym polu wyświetlany jest komunikat podpowiadający, nie blokada.

```
 Makieta — Prompt Builder (stan spoczynku)
 ═════════════════════════════════════
  Nowy prompt                        ⋮
  [Tryb ▼] [Szablony ▼] [Parametry ▼]
 ─────────────────────────────────────
  Temat:       [                    ]
  Styl:        [                  ▼ ]
  Paleta barw: [                  ▼ ]
  Proporcje:   [ 16:9 ▼ ] Warianty:[4▼]
  Wykluczenia: [                    ]
 ─────────────────────────────────────
                       [ ⚡ Generuj ]
 ═════════════════════════════════════
```

### 3.6. Preview Window

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
| Porównanie wariantów obok siebie | Widok siatki wariantów tego samego polecenia | 3 | Menu kebab (⋮) panelu |
| Osadzenie w ramce urządzenia | Wstawienie zasobu w ramkę telefonu, laptopa lub wizytówki | 3 | Menu kebab (⋮) panelu |
| Podgląd prototypu klikalnego | Nawigacja między ekranami makiety wg grafu połączeń | 3 | Menu kebab (⋮) panelu |
| Przełącznik motywu jasny/ciemny | Ocena zasobu w obu motywach systemu tokenów | 3 | Menu kebab (⋮) panelu |
| Diff wizualny | Nałożenie dwóch wersji z podświetleniem różnic pikselowych i strukturalnych | 3 | Menu kebab (⋮) panelu |
| Adnotacja zwrotna | Naniesienie uwagi wprost na podglądzie, przekładanej na kolejne polecenie modyfikujące | 3 | Menu kontekstowe podglądu |
| Podgląd dostępności | Nałożenie oceny kontrastu i symulacji wad wzroku na zasób | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |

**Zachowanie i stany:** aktualizuje się po każdej zmianie zaakceptowanej w Design Board oraz po zakończeniu zadania pętli wykonawczej. Stan „oczekuje decyzji" — widoczne trzy akcje decyzyjne, dostępne od razu; ich obecność nie wstrzymuje pracy w pozostałych oknach modułu ani biegu pozostałych zadań pętli.

```
 Makieta — Preview Window (stan spoczynku)
 ═══════════════════════════════════
  Podgląd: bohater_v3, wariant 2/4 ⋮
  [Format ▼] [Tło ▼] [Przed/po ▼]
 ───────────────────────────────────
  ┌───────────────────────────────┐
  │                               │
  │      [ podgląd grafiki ]      │
  │                               │
  └───────────────────────────────┘
 ───────────────────────────────────
  [ Akceptuj ] [ Regeneruj ] [ Odrzuć ]
 ═══════════════════════════════════
```

### 3.7. Tokens & System Panel

| Pole | Treść |
|---|---|
| Cel | Definiowanie i wydawanie systemu projektowego: tokenów, motywów, reguł dostępności i przewodnika stylu |
| Zawartość | Drzewo tokenów, motywy, tabela kontrastu, powiązania z komponentami |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |

| Sekcja | Zawartość | Warstwa | Sposób wywołania |
|---|---|---|---|
| Drzewo tokenów | Kolory, typografia, odstępy, promienie, cienie w strukturze W3C Design Tokens | 3 | Menu kebab (⋮) obszaru roboczego |
| Motywy | Zestawy nadpisań (jasny, ciemny, warianty marki) wraz z przełącznikiem podglądu | 3 | Zakładka `Motywy` panelu |
| Kontrast i dostępność | Tabela par tokenów z oceną WCAG oraz symulacja wad wzroku | 3 | Zakładka `Dostępność` panelu |
| Eksport i import | Wydanie tokenów do CSS, SCSS, Tailwind, JS, iOS, Android; wczytanie z JSON | 3 | Menu kebab (⋮) panelu |
| Powiązania | Wskazanie komponentów korzystających z danego tokenu | 3 | Kliknięcie tokenu w drzewie |
| Przewodnik stylu | Generowana strona-dokumentacja systemu projektowego | 4 | Polecenie języka naturalnego, wyszukiwarka funkcji |
| Inspektor glifów i par krojów | Przegląd znaków, ligatur, wariantów OpenType oraz dobór par krojów | 4 | Polecenie języka naturalnego, tryb administracyjny |

**Zachowanie i stany:** zmiana wartości tokenu propaguje się do komponentów, kanwy i podglądu. Stan „naruszenie kontrastu" — plakietka ostrzegawcza przy parze tokenów łamiącej wymagany próg; ostrzeżenie nie wstrzymuje pracy.

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

Każda pozycja opisana jest przez nazwę, działanie i zależności (biblioteki backendu w Go, formaty, modele, integracje). Backend platformy jest w Go; funkcje wymagające modelu neuronowego wykonywane są przez kanał modelu (API) lub przez własny endpoint obliczeniowy, zgodnie ze strategią modeli platformy. Zadania katalogu realizowane są w pętli wykonawczej prezentowanej w Execution Loop Window.

### 4.1. Grupa A — Generowanie grafiki przez AI

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| A1 | Text-to-image | Generuje obraz z opisu skomponowanego w Prompt Builderze | Modele obrazowe przez kanał API; własny endpoint GPU; format PNG/WEBP |
| A2 | Image-to-image | Przekształca obraz referencyjny wg promptu, z regulacją siły wpływu | Model obrazowy z parametrem `strength`; wejście: dowolny raster |
| A3 | Inpainting (domalowanie fragmentu) | Zamienia zaznaczony obszar obrazu nowym, wpasowanym w otoczenie | Model inpaint i maska alfa (PNG RGBA); `govips` do składania masek |
| A4 | Outpainting (rozszerzenie kadru) | Domalowuje treść poza pierwotną ramką obrazu | Model outpaint; `disintegration/imaging` do powiększania płótna |
| A5 | Wariacje | Tworzy N wariantów tego samego zasobu (seed i odchylenie) do wyboru | Parametr seed i liczba wariantów; generowanie wsadowe Prompt Buildera |
| A6 | Upscaling neuronowy | Powiększa obraz 2×/4×/8× bez utraty ostrości | Real-ESRGAN / SwinIR przez endpoint GPU; ONNX przez `yalue/onnxruntime_go` |
| A7 | Usuwanie tła | Odcina obiekt od tła, tworzy kanał alfa | U²-Net / RMBG-1.4 (ONNX) przez endpoint; wynik PNG RGBA |
| A8 | Transfer stylu | Przenosi styl obrazu-wzorca na treść obrazu-bazy | IP-Adapter / model transferu stylu; obraz referencyjny Prompt Buildera |
| A9 | Wektoryzacja (raster→SVG) | Zamienia bitmapę na czyste ścieżki wektorowe | `dennwc/gotrace` (potrace) lub wektoryzacja przez kanał API; wynik SVG |
| A10 | Generowanie ikon i logotypów | Tworzy spójny zestaw ikon i znaków w jednym stylu | Model obrazowy z wymuszeniem stylu przez preset; wyjście SVG/PNG |
| A11 | Rozdzielenie na warstwy | Rozkłada wygenerowany obraz na obiekty i warstwy edytowalne | Segment Anything przez endpoint; maski przenoszone na warstwy Design Board |
| A12 | Naprawa twarzy i detalu | Poprawia zniekształcone twarze oraz drobne detale generacji | GFPGAN / CodeFormer przez endpoint GPU |
| A13 | Sterowanie kompozycją (mapy sterujące) | Wymusza pozę, szkielet, głębię lub krawędzie w generacji wg szkicu | ControlNet (canny, depth, pose) przez endpoint; mapa sterująca jako obraz |

### 4.2. Grupa B — Edycja rastrowa

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| B1 | Kadrowanie i prostowanie | Przycina, obraca, prostuje horyzont, zmienia proporcje | `disintegration/imaging` (Crop, Rotate); reguły proporcji 1:1, 16:9, 4:5 |
| B2 | Skalowanie i resampling | Zmienia rozmiar z wyborem algorytmu (Lanczos, biliniowy) | `govips` / `bimg` (libvips, Lanczos3); zachowanie proporcji |
| B3 | Warstwy rastrowe | Niezależne warstwy z trybami mieszania i krycia | Kompozytor warstw (multiply, screen, overlay) — compositing `govips` |
| B4 | Maski warstw | Nieniszcząca maska alfa, malowana lub z zaznaczenia | Kanał alfa PNG; pędzel maski w kliencie (Canvas API), złożenie w Go |
| B5 | Korekcja barwna | Jasność, kontrast, ekspozycja, temperatura, krzywe, HSL | `govips` (LUT, krzywe); podgląd na żywo w kliencie |
| B6 | Retusz (klonowanie i leczenie) | Klonuje i wygładza fragmenty, usuwa niedoskonałości | Inpaint lokalny (A3) oraz pędzel klonujący |
| B7 | Filtry i efekty | Rozmycie, wyostrzenie, ziarno, winieta, cień, poświata | `govips` (gaussian blur, sharpen); efekty warstwy |
| B8 | Zaznaczanie obiektu | Automatyczne zaznaczenie obiektu jednym kliknięciem | Segmentacja (A11) lub różdżka progowa; maska zaznaczenia |
| B9 | Korekcja perspektywy | Prostuje zniekształcenia perspektywiczne i obiektywu | Transformacja homograficzna; `fogleman/gg` do przekształceń |
| B10 | Tryb wsadowy edycji | Stosuje ten sam zestaw operacji do wielu zasobów naraz | Pipeline operacji zapisany jako preset; kolejka modułu Automations |

### 4.3. Grupa C — Grafika i edycja wektorowa

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| C1 | Narzędzie pióra (ścieżki Béziera) | Rysuje i edytuje krzywe oraz węzły ścieżek | SVG `<path>`; edycja węzłów w kliencie; serializacja `ajstarks/svgo` |
| C2 | Kształty podstawowe i złożone | Prostokąty, elipsy, wielokąty, gwiazdy, zaokrąglenia | Prymitywy SVG; parametry promienia i liczby wierzchołków |
| C3 | Operacje logiczne (boolean) | Suma, różnica, przecięcie, wykluczenie ścieżek | Booleany `tdewolff/canvas`; wynik jako pojedyncza ścieżka |
| C4 | Obrys i wypełnienie | Grubość, zakończenia, gradient, deseń, reguła wypełnienia | Atrybuty SVG (stroke, fill, gradient defs); edytor gradientu |
| C5 | Tekst na ścieżce i obrys tekstu | Układa tekst wzdłuż krzywej, zamienia w kontury | `golang/freetype` / `go-text/typesetting`; SVG `textPath` |
| C6 | Optymalizacja SVG | Czyści i minimalizuje kod SVG bez utraty jakości | Reguły czyszczenia (metadane, scalanie ścieżek); wyjście SVG |
| C7 | Symbole i instancje | Definicja wielokrotnego elementu z propagacją zmian | SVG `<symbol>` / `<use>`; rejestr symboli w projekcie |
| C8 | Siatka i przyciąganie precyzyjne | Rozmieszczanie do siatki, prowadnic, punktów, pikseli | Silnik przyciągania w kliencie; jednostki px/rem |
| C9 | Eksport wektorowy | Wydaje czysty SVG, PDF wektorowy i EPS | `go-pdf/fpdf` (wektor→PDF); serializacja SVG |

### 4.4. Grupa D — Projektowanie UI i makiety

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| D1 | Ramki i obszary robocze | Wydzielone ekrany o zdefiniowanych rozmiarach urządzeń | Presety rozmiarów urządzeń; model sceny Design Board |
| D2 | Komponenty i warianty | Element wielokrotny z wariantami stanu (hover, active, disabled) | Rejestr komponentów oraz zestaw właściwości wariantu |
| D3 | Auto-layout | Automatyczne rozmieszczanie z odstępami i dopasowaniem do treści | Silnik flex/stack w kliencie; parametry gap, padding, wyrównanie |
| D4 | Więzy responsywne | Zachowanie elementów przy zmianie rozmiaru ramki | Reguły kotwiczenia i rozciągania; przeliczanie układu |
| D5 | Prototypowanie i przejścia | Łączy ekrany interakcjami i animacjami przejść | Graf połączeń ekranów; podgląd klikalny w Preview Window |
| D6 | Wireframe niskiej wierności | Szybkie makiety szkicowe z biblioteką prostych elementów | Biblioteka elementów szkicowych Design Board |
| D7 | Biblioteka UI | Gotowe komponenty (przyciski, pola, karty) do składania makiet | Komponenty zgodne z systemem wizualnym Danaco; import z tokenów (E) |
| D8 | Generowanie makiety z opisu | Tworzy szkielet ekranu z polecenia tekstowego | Model generujący układ → drzewo komponentów; mapowanie na D2/D7 |
| D9 | Import ze zrzutu ekranu | Odtwarza edytowalną makietę z obrazu istniejącego interfejsu | Rozpoznanie układu (segmentacja i OCR `otiai10/gosseract`) → komponenty |
| D10 | Siatki układu | Kolumny, rynny, moduły, siatka bazowa | Definicja siatki (kolumny, gap, margines); nakładka na ramkę |

### 4.5. Grupa E — Tokeny projektowe i system projektowy

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| E1 | Tokeny kolorów | Definicja nazwanych barw semantycznych (primary, surface, danger) | Format W3C Design Tokens (JSON, `$value`/`$type`); walidacja aliasów |
| E2 | Tokeny typografii | Skala rozmiarów, krój, interlinia, grubość, tracking | Tokeny typograficzne W3C; podgląd skali |
| E3 | Tokeny odstępów i promieni | Skala spacing, radius, rozmiary, cienie, z-index | Tokeny wymiarowe W3C; jednostki px/rem |
| E4 | Motyw jasny i ciemny | Warianty tokenów dla trybów i marek | Zestawy nadpisań tokenów; przełącznik trybu w podglądzie |
| E5 | Eksport tokenów do kodu | Wydaje tokeny jako CSS vars, SCSS, Tailwind, JS, iOS, Android | Silnik transformacji tokenów; wyjścia `.css`, `.scss`, `tailwind.config`, `.json`, `.xml` |
| E6 | Import tokenów | Wczytuje istniejące tokeny do modułu | Parser JSON W3C i `tailwind.config`; mapowanie na model tokenów |
| E7 | Powiązanie tokenów z komponentami | Komponenty UI (D2) czerpią wartości z tokenów, zmiana propaguje | Referencje tokenów we właściwościach komponentu; reaktywna aktualizacja |
| E8 | Dokumentacja systemu | Generuje przewodnik: kolory, typografia, komponenty | Render HTML przewodnika; wydanie do Library i Studio |

### 4.6. Grupa F — Kolor

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| F1 | Generator palet | Tworzy paletę z reguły harmonii (mono, komplementarna, triada) | `lucasb-eyer/go-colorful` (konwersje HSL/Lab, harmonie) |
| F2 | Ekstrakcja palety z obrazu | Wyciąga dominujące barwy z zasobu | Kwantyzacja (median cut, k-means); `go-colorful` |
| F3 | Kontroler kontrastu WCAG | Liczy współczynnik kontrastu i ocenę AA/AAA dla par barw | Wzór luminancji WCAG 2.1; `go-colorful`; podgląd wyniku |
| F4 | Symulacja wad wzroku | Podgląd projektu w protanopii, deuteranopii i tritanopii | Macierze symulacji daltonizmu na buforze obrazu |
| F5 | Edytor gradientów | Gradienty liniowe, promieniste i kątowe z wieloma stopniami | SVG/CSS gradient defs; interpolacja w przestrzeni Lab (`go-colorful`) |
| F6 | Konwersja przestrzeni barw | HEX ↔ RGB ↔ HSL ↔ Lab ↔ CMYK, próbki nazwane | `go-colorful`; tablice barw nazwanych |
| F7 | Sprawdzian dostępności kolorem | Wskazuje pary tokenów łamiące kontrast w całym systemie | Reguły WCAG na zestawie tokenów (E1); raport naruszeń |

### 4.7. Grupa G — Ikony i typografia

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| G1 | Biblioteka ikon | Przeszukiwalny katalog ikon otwartoźródłowych do wstawienia | Zestawy SVG (Lucide, Heroicons, Tabler); indeks nazw i tagów |
| G2 | Edytor ikony na siatce | Rysuje i poprawia ikonę na siatce 24×24 z wyrównaniem do pikseli | Edytor SVG (Grupa C), siatka ikony, kształty prowadzące |
| G3 | Generowanie zestawu ikon | Tworzy spójny komplet ikon w jednym stylu z listy pojęć | Model (A10) z wymuszeniem siatki i grubości; wyjście SVG |
| G4 | Font ikon i sprite | Pakuje ikony w font ikon lub sprite SVG z symbolami | Generator fontu przez moduł Terminal lub sprite `<symbol>`; `.woff2`, `.svg` |
| G5 | Generator faviconów | Wydaje komplet faviconów i ikon aplikacji ze źródła | `Kodeworks/golang-image-ico`; render 16/32/180/512 px; `.ico`, `.png`, manifest |
| G6 | Dobór par krojów | Proponuje pasujące zestawienia krojów nagłówka i tekstu | Reguły typograficzne i model; podgląd zestawień |
| G7 | Podgląd i osadzenie webfontów | Podgląd kroju w tekście próbnym, generacja `@font-face` | Metadane fontu (`golang/freetype`); podzbiór glifów |
| G8 | Inspektor glifów | Przegląd znaków, ligatur i wariantów OpenType kroju | Parser OpenType; render glifów |

### 4.8. Grupa H — Eksport, podgląd, handoff

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| H1 | Eksport wieloformatowy | Wydaje zasób jako PNG, JPG, WEBP, AVIF, SVG, PDF, ICO | `govips` / `bimg`, `chai2010/webp`, AVIF przez libvips; SVG i PDF natywnie |
| H2 | Kompresja i optymalizacja | Redukuje wagę pliku sterując jakością i paletą | `govips` (jakość, chroma subsampling), pngquant/oxipng przez Terminal |
| H3 | Skalowanie @1x/@2x/@3x | Generuje warianty gęstości pikseli jednym poleceniem | Mnożniki eksportu; resize `govips`; nazwy `@2x`, `@3x` |
| H4 | Cięcie na fragmenty | Eksportuje wskazane obszary makiety jako osobne pliki | Definicja ramek eksportu; wsadowe cięcie `imaging.Crop` |
| H5 | Handoff i inspekcja | Udostępnia wymiary, odstępy, kolory, typografię i CSS elementu | Odczyt modelu sceny → panel wartości i generacja CSS |
| H6 | Eksport kodu widoku | Zamienia makietę lub element w kod widoku (CSS, Tailwind, React) | Serializacja drzewa komponentów → JSX/HTML+CSS; mapowanie tokenów (E5); wydanie do Apps |
| H7 | Eksport tablicy jako PDF lub obraz | Zapisuje całą kompozycję Design Board do jednego pliku | `go-pdf/fpdf`; render kanwy → PDF/PNG |
| H8 | Osadzenie w ramce urządzenia | Wstawia zasób w realistyczną ramkę telefonu, laptopa lub wizytówki | Biblioteka ramek PNG/SVG; kompozycja perspektywiczna (`fogleman/gg`) |
| H9 | Eksport animacji lekkiej | Wydaje animację jako JSON Lottie, GIF lub APNG | Format Lottie JSON; `image/gif` (biblioteka standardowa); render klatek |
| H10 | Metadane i profil barwny | Osadza i oczyszcza EXIF, ICC oraz dane autorstwa przy eksporcie | `dsoprea/go-exif`; osadzanie profili ICC przez libvips |

### 4.9. Grupa I — Marketing, szablony, zasoby zewnętrzne

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| I1 | Szablony formatów społecznościowych | Gotowe rozmiary postów, relacji i okładek dla platform społecznościowych | Presety rozmiarów i szablony układu Design Board |
| I2 | Zestawy rozmiarów kampanii | Generuje ten sam projekt w wielu formatach reklamowych naraz | Reguła przeskalowania układu (auto-layout D3) na zestaw ramek |
| I3 | Biblioteka szablonów brandingowych | Arkusze brandingu, tablice nastroju, prezentacje, wizytówki | Szablony układu Design Board; katalog szablonów modułu |
| I4 | Import zasobów zewnętrznych | Wyszukuje i wstawia zasoby z bibliotek zdjęć i ilustracji | Integracje API bibliotek zewnętrznych (klucze jawne w oknie konfiguracji); import do Assets Panel |
| I5 | Placeholdery treści | Wstawia realistyczne treści próbne (nazwiska, teksty, awatary) | Generator danych próbnych i awatary AI (A1); zasilanie makiet |
| I6 | Znak wodny i branding wsadowy | Nakłada logo lub znak wodny na zestaw zasobów | Kompozycja warstwy (B3) w trybie wsadowym (B10) |

### 4.10. Grupa J — Współpraca, wersjonowanie, organizacja

| # | Nazwa | Co robi | Zależności |
|---|---|---|---|
| J1 | Warstwy i drzewo obiektów | Panel warstw z widocznością, blokadą, grupami i zagnieżdżeniem | Model sceny Design Board |
| J2 | Wersjonowanie kompozycji | Historia stanów tablicy, nazwane wersje, powrót, porównanie | Snapshoty sceny (Model danych — wersjonowanie); diff wizualny |
| J3 | Diff wizualny zasobów | Nakłada dwie wersje i podświetla różnice pikselowe oraz strukturalne | `corona10/goimagehash` (pHash), różnica pikseli; suwak przed/po |
| J4 | Komentarze i adnotacje | Przypięte uwagi do elementów, wątki, oznaczenia osób | Model adnotacji; powiązanie z Chat Window |
| J5 | Kursory i obecność | Widoczność kursorów współpracujących osób na tablicy | Kanał WebSocket platformy |
| J6 | Tagi, kolekcje, wyszukiwanie | Etykietowanie i przeszukiwanie zasobów po nazwie, tagu i prompcie | Assets Panel; indeks pełnotekstowy promptów |
| J7 | Metadane pochodzenia (prowenancja) | Zapisuje model, prompt, seed, datę i łańcuch edycji zasobu | Panel metadanych Assets Panel; zapis w Modelu danych; zgodność z blokiem prowenancji klienta |

Grupy A–J obejmują łącznie **73 funkcje i narzędzia** pokrywające obszar tematyczny modułu. Część funkcji opiera się wyłącznie na kodzie klienta (kanwa, edytor SVG), część na bibliotekach Go w backendzie (raster, eksport, kolor), a część korzysta z zasobu GPU i modelu neuronowego (generowanie, upscaling, segmentacja).

---

## 5. Katalog elementów interfejsu

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Gdzie występuje |
|---|---|---|---|---|---|---|---|---|
| Pole polecenia | `.dn-pole-polecenia` na całą szerokość kolumny | Wydawanie poleceń Wykonawcy w języku naturalnym | Element wiodący lewej kolumny | 1 | Widoczne bez interakcji | domyślny, aktywny, wysyłanie | Wysyła polecenie do Wykonawcy i uruchamia zadanie w pętli wykonawczej | Chat Window |
| Karta zadania pętli | Wiersz kolejki `.dn-karta--zadanie` | Reprezentacja pojedynczego zadania pętli wykonawczej | Wąska karta listowa | 1 | Widoczna bez interakcji | w kolejce, w toku, gotowe, ponowione, przerwane | Kliknięcie rozwija szczegóły; menu kebab udostępnia ponowienie | Execution Loop Window |
| Sterowanie przebiegiem pętli | `.dn-btn--zarys` w zestawie trzech | Wstrzymanie, wznowienie i przerwanie pętli wykonawczej | Trzy małe przyciski w rzędzie | 1 | Widoczne bez interakcji | domyślny, wstrzymany, przerwany | Zmienia stan pętli natychmiast, bez modala potwierdzającego | Execution Loop Window |
| Znacznik kontekstowy | `.dn-plakietka--kontekst` | Prezentacja i zmiana środowiska, projektu, modelu, wykonawcy | Bardzo mała pigułka | 1 | Widoczny bez interakcji; kliknięcie otwiera selektor | domyślny, otwarty selektor | Kliknięcie otwiera selektor warstwy 2, który zwija się po wyborze | Chat Window, Design Board |
| Miniatura zasobu | Kafel obrazu w siatce `.dn-karta--interaktywna` | Reprezentacja pojedynczego zasobu wizualnego | Mały–średni kafel kwadratowy | 1 | Widoczna po otwarciu panelu | domyślny, hover, nowy (plakietka), wybrany (obrys złoty) | Kliknięcie otwiera Preview Window; przeciągnięcie przenosi na Design Board | Assets Panel |
| Przycisk „⚡ Generuj" | `.dn-btn--zloty` | Najważniejsze CTA okna Prompt Builder | Średni przycisk wypełniony, cień złoty | 1 | Widoczny po otwarciu panelu | domyślny, generowanie (spinner) | Zawsze aktywny — zleca generowanie do pętli wykonawczej; przy pustym poleceniu komunikat podpowiadający zamiast blokady | Prompt Builder |
| Suwak parametru | `.dn-suwak` w wariancie zakresowym | Regulacja parametru generowania | Mały element liniowy | 2 | Znacznik `Parametry ▼` | domyślny, przeciągany | Zmienia wartość liczbową parametru promptu | Prompt Builder |
| Pole stylu i palety | `.dn-select` | Wybór presetu stylistycznego lub palety barw | Małe pole rozwijane | 2 | Kliknięcie pola | domyślny, otwarte (lista pozycji) | Wybór pozycji aktualizuje podgląd struktury promptu | Prompt Builder |
| Suwak porównawczy przed/po | Uchwyt przeciągalny na osi poziomej | Wizualne porównanie dwóch wersji zasobu | Element średni, na całą szerokość podglądu | 2 | Znacznik `Przed/po ▼` | domyślny, przeciągany | Przesuwa granicę widoczności między wersją „przed" i „po" | Preview Window |
| Przełącznik tła podglądu | Mały zestaw ikon (szachownica, biel, czerń) | Zmiana tła pod zasobem z przezroczystością | Bardzo mała grupa przycisków ikonowych | 2 | Znacznik `Tło ▼` | domyślny, wybrany | Zmienia tło podglądu natychmiastowo | Preview Window |
| Przyciski decyzyjne | `.dn-btn` w trzech wariantach (główny, zarys, duch) | Zamknięcie cyklu oceny wygenerowanego zasobu | Trzy małe–średnie przyciski w rzędzie | 1 | Widoczne po otwarciu panelu | domyślne, po decyzji | Akceptacja przenosi zasób do Assets Panel jako finalny; odrzucenie usuwa wariant i zachowuje go w historii | Preview Window |
| Panel warstw | Drzewo z ikonami widoczności i blokady | Zarządzanie kolejnością i widocznością elementów kompozycji | Kolumna boczna, wąska | 2 | Znacznik `Warstwy ▼` | domyślny, warstwa ukryta, zablokowana | Kliknięcie ikony przełącza widoczność lub blokadę warstwy | Design Board |
| Uchwyty wyrównania i rozmieszczenia | Mały pasek ikon przy zaznaczeniu | Wyrównanie i rozłożenie zaznaczonych elementów | Bardzo mały pasek kontekstowy | 3 | Pasek kontekstowy zaznaczenia wielokrotnego | widoczny tylko przy zaznaczeniu wielokrotnym | Kliknięcie stosuje wybrane wyrównanie do zaznaczonych elementów | Design Board |
| Znacznik adnotacji | Mała pinezka na kanwie | Przypięcie komentarza do elementu kompozycji | Bardzo mała ikonka | 3 | Menu kebab (⋮), skrót klawiszowy | domyślny, rozwinięty (dymek treści) | Kliknięcie otwiera i zamyka treść komentarza | Design Board |
| Tagi i kolekcje | `.dn-plakietka` w wariancie neutralnym | Kategoryzacja zasobów | Bardzo małe pigułki | 2 | Widoczne przy miniaturze zasobu | domyślny, aktywny filtr (akcent złoty) | Kliknięcie filtruje widok Assets Panel do wybranego tagu | Assets Panel |
| Skrót „Wyślij do modułu" | `.dn-btn--duch`, ikonowy | Przekazanie zasobu do Studio lub Apps | Mały przycisk bez obrysu | 2 | Przycisk `→ Moduł` | domyślny, hover | Otwiera wybór modułu docelowego i wykonuje transfer | Assets Panel, Preview Window |
| Wiersz tokenu | `.dn-wiersz--token` z próbką barwy | Prezentacja i edycja pojedynczego tokenu systemu | Wąski wiersz drzewa | 3 | Menu kebab (⋮) obszaru roboczego | domyślny, edytowany, naruszenie kontrastu | Zmiana wartości propaguje się do komponentów i kanwy | Tokens & System Panel |

---

## 6. Przepływy pracy w module

### 6.1. Od polecenia do zaakceptowanego zasobu

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

### 6.2. Zestawianie kompozycji z wielu zasobów

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

### 6.3. Budowa i wydanie systemu projektowego

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

### 6.4. Przekazanie zasobu do modułu docelowego

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

Operator personalizuje moduł z okna konfiguracji (architektura/model-konfiguracji.md), z zachowaniem warstwowości (globalna → środowisko → projekt → sesja) i zasady „brak ustawienia = wartość domyślna". Wszystkie klucze są jawne, żaden nie blokuje pracy.

| # | Punkt sterowania | Co personalizuje | Zakres konfiguracji (5.x) | Wartość domyślna |
|---|---|---|---|---|
| K1 | Domyślny silnik generujący | Model text-to-image używany bez ręcznego wyboru | 5.4 Zachowanie modeli | pierwszy skonfigurowany kanał |
| K2 | Endpointy GPU zadań neuronowych | Adresy własnych endpointów dla upscalingu, usuwania tła, segmentacji i map sterujących | 5.4, 5.7 Rozszerzenia | brak (funkcje AI przez kanał API) |
| K3 | Klucze API modeli i bibliotek zewnętrznych | Jawne klucze dostępu do modeli obrazowych i bibliotek zasobów (I4) | 5.4, 5.8 Integracje | brak (wprowadzane przez Operatora) |
| K4 | Domyślne formaty i jakość eksportu | Preset formatu (PNG, WEBP, AVIF), poziom kompresji, zestaw skal @1/2/3x | 5.1 Aplikacja | PNG, jakość 90, @1x |
| K5 | Domyślna paleta i system tokenów | Zestaw tokenów (E1–E4) wstrzykiwany do nowych projektów | 5.8 Integracje | system wizualny Danaco |
| K6 | Reguły dostępności | Próg kontrastu (AA/AAA) i ostrzeżenia w Tokens & System Panel | 5.1 Aplikacja | AA, ostrzeżenia włączone |
| K7 | Powiązanie z modułem Studio | Sposób przekazywania zasobów do Studio Editor (ręcznie lub synchronizacja) | 5.8 Integracje | wyłączone (ręczne wysłanie) |
| K8 | Powiązanie z modułem Apps | Sposób przekazywania kodu i zasobów do Frontend Workspace | 5.8 Integracje | wyłączone (ręczne wysłanie) |
| K9 | Powiązanie z Library | Archiwizacja zaakceptowanych zasobów w Library | 5.8 Integracje | wyłączone |
| K10 | Wsadowe zadania przez Automations | Korzystanie trybu wsadowego (B10, I2, I6) z kolejek modułu Automations | 5.3 Akcje, 5.8 | wyłączone |
| K11 | Potwierdzanie usuwania zasobu | Wymóg potwierdzenia usunięcia w Assets Panel (zawsze z „Cofnij") | 5.1 Aplikacja | bez potwierdzenia, „Cofnij" dostępne |
| K12 | Zapis prowenancji | Szczegółowość metadanych pochodzenia (J7) | 5.9 Historia, 5.10 Pamięć | zapis pełny |
| K13 | Domyślny tryb Preview | Motyw (jasny, ciemny), domyślna ramka urządzenia, tło przezroczystości | 5.1, 5.12 Karty sesji | motyw systemowy, bez ramki, szachownica |
| K14 | Widoczność wariantów i seed | Liczba wariantów domyślnych, przypinanie seed dla powtarzalności | 5.4 | 4 warianty, seed swobodny |
| K15 | Ścieżki eksportu docelowego | Katalog i kolekcja Library oraz miejsce w projekcie Apps dla wydań | 5.8, 5.13 Komponenty | brak (wybór przy wydaniu) |
| K16 | Zachowanie Execution Loop Window | Widoczność kolumny pętli w spoczynku, poziom szczegółowości komunikatów sterujących, próg automatycznego ponowienia zadania | 5.3 Akcje, 5.9 Historia | kolumna zwinięta w spoczynku, komunikaty skrócone, jedno ponowienie |
| K17 | Warstwy widoczności interfejsu modułu | Zestaw funkcji ujawnianych w warstwach 2–4 zależnie od roli użytkownika | 5.1 Aplikacja, 5.11 Role | warstwy 1–3 dla roli podstawowej, warstwa 4 dla roli rozszerzonej |

Każdy punkt sterowania opatrzony jest objaśnieniem kontekstowym `[?]` — zgodnie z 3.3 Modelu konfiguracji.

---

## 8. Stany, dane i powiązania z innymi modułami

### 8.1. Dane wykorzystywane przez AI w module

| Źródło danych | Okno pochodzenia | Uwaga |
|---|---|---|
| Polecenia Użytkownika | Chat Window | Treść polecenia, kontekst rozmowy, załączniki referencyjne |
| Zlecenia i zadania pętli | Execution Loop Window | Dekompozycja zlecenia, stan zadań, wyniki kontroli jakości |
| Polecenia generujące | Prompt Builder | Struktura promptu i historia poleceń |
| Zgromadzone zasoby | Assets Panel | Metadane, tagi, warianty, prowenancja |
| Kompozycja robocza | Design Board | Układ, warstwy, adnotacje, komponenty |
| System projektowy | Tokens & System Panel | Tokeny, motywy, reguły dostępności |

### 8.2. Stany zasobu wizualnego

| Stan | Znaczenie | Gdzie widoczny |
|---|---|---|
| W generowaniu | Zasób w trakcie tworzenia przez model | Execution Loop Window, Assets Panel (miniatura w budowie) |
| Nowy | Wygenerowany, jeszcze nieoceniony | Assets Panel (plakietka „nowy") |
| Zaakceptowany | Zatwierdzony jako wynik finalny | Preview Window, Assets Panel |
| Odrzucony | Wariant odrzucony, zachowany w historii | Preview Window (historia wariantów) |
| W kompozycji | Zasób umieszczony na Design Board | Design Board (panel warstw) |
| Przekazany | Wysłany do modułu docelowego | Assets Panel (znacznik miejsca docelowego) |

### 8.3. Powiązania konfigurowalne z innymi modułami

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
| Studio | Zasoby wygenerowane w Design wykorzystywane przy redagowaniu dokumentów | Konfiguracyjne | Assets Panel — „Wyślij do modułu" |
| Apps | Zasoby i kod widoku wykorzystywane przy budowie interfejsu produktu, w oknie Frontend Workspace | Konfiguracyjne | Assets Panel, Preview Window — „Wyślij do modułu" |
| Library | Archiwizacja zaakceptowanych zasobów i przewodnika stylu | Konfiguracyjne | Okno konfiguracji, klucz K9 |
| Automations | Kolejkowanie zadań wsadowych modułu | Konfiguracyjne | Okno konfiguracji, klucz K10 |

Powiązania nie są aktywne domyślnie. Przekazanie zasobu do innego modułu odbywa się jako jawna, pojedyncza decyzja użytkownika podejmowana z poziomu Assets Panel lub Preview Window albo — zgodnie z ustawieniem konfiguracyjnym — mechanizmem synchronizacji; ręczne potwierdzanie każdego przekazania jest ustawieniem konfiguracyjnym, nie wymogiem.

---

## 9. Scenariusze użycia

### 9.1. Budowa spójnego zestawu ilustracji brandingowych

1. Użytkownik otwiera moduł Design w środowisku WorkSpace.
2. W Chat Window zleca przygotowanie zestawu ilustracji; Koordynator rozkłada zlecenie na zadania widoczne w Execution Loop Window.
3. W Prompt Builderze Użytkownik definiuje styl bazowy (paleta granat i złoto, styl geometryczny) i zapisuje go jako szablon.
4. Wykonawca generuje kolejne ilustracje na podstawie szablonu; pętla wykonawcza prowadzi zadania generowania, usunięcia tła i eksportu, zachowując spójność stylistyczną zestawu.
5. Zaakceptowane warianty trafiają do Assets Panel, oznaczone wspólnym tagiem kolekcji „Branding — kampania Q3".
6. Na Design Board Użytkownik zestawia całość w jedną tablicę prezentacyjną i eksportuje ją jako PDF do przeglądu wewnętrznego.

### 9.2. Projektowanie interfejsu produktu równolegle z modułem Apps

1. Zespół pracujący nad produktem w module Apps (środowisko CodeStudio) potrzebuje zestawu elementów interfejsu.
2. W Tokens & System Panel projektant ustala tokeny kolorów, typografii i odstępów oraz sprawdza kontrast par barw.
3. W module Design projektant generuje ikony i elementy graficzne interfejsu, oceniając każdy wariant w Preview Window.
4. Zaakceptowane elementy oraz tokeny wydane jako zmienne CSS przekazywane są do okna Frontend Workspace modułu Apps przyciskiem „Wyślij do modułu".
5. Zespół frontendowy w module Apps korzysta z przekazanych zasobów bez opuszczania własnej przestrzeni roboczej.

### 9.3. Iteracyjna korekta ilustracji na podstawie uwag

1. Wygenerowana ilustracja trafia do oceny zespołu w Design Board.
2. Uwagi nanoszone są jako adnotacje przypięte do konkretnych fragmentów kompozycji.
3. Na podstawie zebranych adnotacji Użytkownik formułuje w Chat Window polecenie modyfikujące („ciemniejsze tło, usuń element w lewym rogu kadru").
4. Execution Loop Window prezentuje zadanie korekty, wynik kontroli jakości i decyzję o ponowieniu.
5. Nowy wariant trafia do Preview Window z aktywnym suwakiem porównawczym przed/po, ułatwiającym ocenę zmiany względem poprzedniej wersji.

---

## 10. Zgodność z rdzeniem platformy

| Zasada nadrzędna | Realizacja w module |
|---|---|
| Dwa kanały komunikacji | Chat Window (Użytkownik ↔ Wykonawca) jako centralny punkt pracy i podstawowy mechanizm sterowania; Execution Loop Window (Koordynator ↔ Wykonawca) jako okno pętli wykonawczej zadań modułu |
| Układ pionowy | Chat Window w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej, Design Board w prawej kolumnie dominującej, panele jako rozszerzenia boczne |
| Warstwy widoczności | Cztery warstwy ujawniania funkcji (rozdz. 2.1); w spoczynku widoczna warstwa 1 oraz zwinięte wyzwalacze |
| Zero blokad | „⚡ Generuj" zawsze aktywny; usuwanie z „Cofnij"; decyzje w Preview Window nieblokujące; brak modali-bramek |
| Klucze jawne | Klucze modeli i integracji (K3) wprowadzane wprost przez Operatora, widoczne w oknie konfiguracji |
| Pełna konfigurowalność | Siedemnaście punktów sterowania (rozdz. 7) z warstwowością i objaśnieniami `[?]` |
| Jawność zależności | Powiązania z Studio, Apps, Library i Automations (K7–K10) jako świadome ustawienia, nie ukryte reguły |
| Konfigurowalność modeli | Wybór silnika i endpointu (K1, K2) zgodny ze strategią modeli platformy |
| Audytowalność | Prowenancja zasobu (J7) — model, prompt, seed, łańcuch edycji — zgodna z blokiem prowenancji klienta |

---

*Koniec dokumentu. Moduł Design — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
