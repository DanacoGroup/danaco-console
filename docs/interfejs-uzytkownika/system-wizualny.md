# Danaco Console — System wizualny

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
| **Źródło** | Koncepcja platformy i architektura |

Dokument opracowuje system projektowy platformy Danaco Console na bazie przewodnika marki i towarzyszących mu plików źródłowych (`tokens.css`, `tokens.json`, `components.css`, katalog `ikony/`). Ustala tokeny kolorów w motywie jasnym i ciemnym, typografię, cienie, siatkę kolumnową, warstwy widoczności, ikonografię i bibliotekę komponentów, a następnie opisuje ich zastosowanie w oknach komunikacji, w oknach operacyjnych modułów oraz w kolumnach strony głównej. Zbiór tokenów i komponentów jest tożsamy z ustalonym w przewodniku marki; niniejszy dokument porządkuje ich zastosowanie w strukturze platformy opisanej w Koncepcji i w Architekturze.

Zastosowania systemu wizualnego w strukturze platformy, wraz z podstawą w Koncepcji platformy:

| Zastosowanie | Rozdział niniejszego dokumentu | Podstawa w Koncepcji platformy |
|---|---|---|
| Warstwy widoczności | 3a | rozdz. 4 (warstwy widoczności interfejsu) |
| Okna komunikacji — Chat Window i Execution Loop Window | 10.2 | rozdz. 3, Załącznik A |
| Okna operacyjne modułów | 10 | Załącznik A (macierz okien operacyjnych) |
| Strona główna — trzy strefy kolumnowe | 11 | rozdz. 7.4 (forma prezentacji stref) |
| Okno konfiguracji punktów izolacji | 12 | rozdz. 6.4–6.6 |
| Panel orkiestracji środowiska MultitaskingAI | 13 | rozdz. 13.8 |

---

## Spis treści

1. [Zasady systemu wizualnego marki](#1-zasady-systemu-wizualnego-marki)
2. [Tokeny kolorów](#2-tokeny-kolorów)
3. [Typografia](#3-typografia)
3a. [Warstwy widoczności w systemie wizualnym](#3a-warstwy-widoczności-w-systemie-wizualnym)
4. [Siatka kolumnowa, odstępy, promienie](#4-siatka-kolumnowa-odstępy-promienie)
5. [Cienie i ruch](#5-cienie-i-ruch)
6. [Stany interakcji i fokus](#6-stany-interakcji-i-fokus)
7. [Ikonografia](#7-ikonografia)
8. [Biblioteka komponentów](#8-biblioteka-komponentów)
9. [Motyw jasny i ciemny — mechanizm przełączania](#9-motyw-jasny-i-ciemny--mechanizm-przełączania)
10. [Zastosowanie w oknach operacyjnych modułów](#10-zastosowanie-w-oknach-operacyjnych-modułów)
11. [Zastosowanie na stronie głównej — Centrum dowodzenia](#11-zastosowanie-na-stronie-głównej--centrum-dowodzenia)
12. [Zastosowanie w oknie konfiguracji punktów izolacji](#12-zastosowanie-w-oknie-konfiguracji-punktów-izolacji)
13. [Zastosowanie w panelu orkiestracji środowiska MultitaskingAI](#13-zastosowanie-w-panelu-orkiestracji-środowiska-multitaskingai)
14. [Zgodność z zasadami nadrzędnymi platformy](#14-zgodność-z-zasadami-nadrzędnymi-platformy)
- [Załącznik A. Pełna tabela tokenów kolorów semantycznych](#załącznik-a-pełna-tabela-tokenów-kolorów-semantycznych)
- [Załącznik B. Tabela klas komponentów](#załącznik-b-tabela-klas-komponentów)
- [Załącznik C. Katalog ikon i typowe zastosowanie](#załącznik-c-katalog-ikon-i-typowe-zastosowanie)

---

## 1. Zasady systemu wizualnego marki

### 1.1. Charakter marki

System wizualny Danaco opiera się na parze kolorów **granat + złoto**, w stylu określonym w przewodniku marki jako „kancelaryjno-nowoczesny" — powaga kancelarii prawnej połączona z klarownością współczesnego oprogramowania. Rolę obu kolorów marki porządkuje poniższa tabela.

| Kolor | Rola w systemie | Zastosowanie | Ograniczenie |
|---|---|---|---|
| Granat | Kolor dominujący | Paski nawigacji, powierzchnie trybu ciemnego, główne akcenty marki | — |
| Złoto | Akcent sygnaturowy | Wywołania do działania (CTA), stany aktywne i najechania, podkreślenia, godło | Stosowane oszczędnie; nigdy nie wypełnia dużych powierzchni interfejsu |

Ta sama warstwa wizualna obsługuje wszystkie produkty rodziny Danaco (Hub, Mail, Lynx, aplikację mobilną) oraz platformę Danaco Console opisaną w Koncepcji. Danaco Console przejmuje tokeny i komponenty bez modyfikacji, dopasowując wyłącznie ich zastosowanie do własnej struktury środowisk, modułów i okien operacyjnych (rozdziały 10–13 niniejszego dokumentu).

### 1.2. Źródła systemu

System wizualny jest opisany w niniejszym dokumencie wraz z towarzyszącymi mu plikami źródłowymi. Poniższa tabela porządkuje rolę poszczególnych plików; niniejszy dokument omawia ich zawartość merytorycznie i wskazuje zastosowanie w platformie Danaco Console, nie zastępując ich jako źródła wartości.

| Plik | Rola |
|---|---|
| `tokens.css` | Zmienne CSS (`--dn-*`) dla motywu jasnego i ciemnego. Źródło prawdy dla warstwy webowej. |
| `tokens.json` | Te same wartości w formacie maszynowym — dla warstwy JavaScript oraz aplikacji mobilnej (rozdz. 4 Architektury: interfejs TypeScript, silniki webview per platforma). |
| `components.css` | Biblioteka komponentów gotowych do użycia, prefiks klas `.dn-`. Wymaga wczytanego `tokens.css`. |
| `manifest.json` | Snapshot systemu wizualnego — identyfikator wersji, suma kontrolna plików, zasada jawnej synchronizacji. |
| `ikony/` | Zestaw 47 ikon SVG w jednolitej siatce, wraz z galerią `indeks.html`. |

### 1.3. Zgodność z zasadami nadrzędnymi platformy

System wizualny podlega tym samym zasadom nadrzędnym, co reszta platformy (Koncepcja platformy, rozdz. 14). Motyw jasny i ciemny są dwiema równoprawnymi, w pełni obsłużonymi wersjami tego samego systemu tokenów, nie wariantem podstawowym i pochodnym — wybór między nimi jest ustawieniem konfigurowalnym z okna konfiguracji (Architektura, rozdz. 13), a nie sztywną regułą narzuconą przez platformę. Zgodnie z zasadą braku twardych blokad wbudowanych na stałe (dokument nadrzędny „Zero blokerów"), żaden token ani komponent nie koduje zachowania niedostępnego do zmiany z poziomu konfiguracji tam, gdzie konfiguracja tej warstwy jest przewidziana — kolor motywu, typografia i gęstość odstępów pozostają spójne z pozostałymi ustawieniami platformy sterowanymi z okna konfiguracji. Pełne zestawienie zgodności z czterema zasadami nadrzędnymi zawiera rozdział 14.

---

## 2. Tokeny kolorów

### 2.1. Kolory marki

Granat i złoto tworzą pełne skale prymitywów, z których budowane są tokeny semantyczne używane w komponentach. Komponenty platformy odwołują się wyłącznie do tokenów semantycznych (rozdz. 2.2–2.4) — wartości prymitywne poniżej stanowią ich podstawę.

| Skala | Token | Wartości (od najciemniejszej) |
|---|---|---|
| Granat | `--dn-granat-900…100` + `--dn-granat-marka` | `#0B1120 · #10192E · #16223C · #1E2E4E · #26395C · #3A4E78 · #5C6F98 · #9BA9C6 · #D7DEEB` · marka `#1A2A4A` |
| Złoto | `--dn-zloto-700…100` + `--dn-zloto-ghost` | `#8A6D2E · #A6812F · #C9A24B (500, bazowe) · #D9B968 (400, na ciemnym) · #E6C67A (300) · #F5ECD4 (100) · #FBF7EC (ghost)` |
| Szarość | `--dn-szary-0…900` | `#FFFFFF → #12161F`, dziewięć kroków między bielą a niemal-czernią |

Warstwy koloru układają się w łańcuch od wartości bazowych, przez tokeny semantyczne, po komponenty:

```
PRYMITYWY (skale bazowe)          TOKENY SEMANTYCZNE               KOMPONENTY (.dn-*)
────────────────────────          ──────────────────────           ──────────────────
Granat  900…100 + marka   ──┐
Złoto   700…100 + ghost   ──┼──►  --dn-tlo · --dn-powierzchnia  ──►  .dn-btn · .dn-karta
Szarość 0…900             ──┘      --dn-tekst · --dn-akcent · …        .dn-tabela · .dn-pole …
                                   (wartość dla motywu                (odwołują się WYŁĄCZNIE
   wartości bazowe,                 jasnego i ciemnego                 do tokenów semantycznych,
   nie używane wprost                w tokens.css / tokens.json)        nigdy do prymitywów)
```

**Zasada użycia:** złoto jako tekst na jasnym tle stosuje się wyłącznie przez token `--dn-akcent-txt` (`#8A6D2E`) — bazowe złoto `#C9A24B` ma na bieli kontrast niewystarczający dla drobnego tekstu (zob. rozdz. 2.5). Kolory nie są nigdy kodowane na sztywno w komponentach — wyłącznie przez tokeny.

### 2.2. Kolory semantyczne — motyw jasny

| Token | Wartość | Zastosowanie |
|---|---|---|
| `--dn-tlo` | `#F7F8FA` | Tło głównego obszaru roboczego |
| `--dn-powierzchnia` | `#FFFFFF` | Powierzchnia kart, pól, paneli |
| `--dn-powierzchnia-2` | `#EEF1F6` | Powierzchnia drugorzędna (nagłówki tabel, stopki kart, tło hover) |
| `--dn-panel` | `#FFFFFF` | Tło modali i paneli bocznych |
| `--dn-hover` | `#EEF1F6` | Tło stanu najechania |
| `--dn-nakladka` | `rgba(11,17,32,0.42)` | Nakładka pod modalem |
| `--dn-obrys` / `-subtelny` / `-mocny` | `#E2E6EE` / `#EEF1F6` / `#CBD2DF` | Trzy poziomy intensywności obrysu |
| `--dn-tekst` | `#12161F` | Tekst podstawowy |
| `--dn-tekst-2` | `#4C566B` | Tekst drugorzędny |
| `--dn-tekst-3` | `#6E7A93` | Tekst trzeciorzędny — wyłącznie duży lub metadane (rozdz. 2.5) |
| `--dn-tekst-inv` | `#FFFFFF` | Tekst na powierzchniach ciemnych (granat, złoto) |
| `--dn-marka` / `-hover` | `#1A2A4A` / `#1E2E4E` | Kolor marki i jego stan najechania |
| `--dn-akcent` / `-txt` / `-tlo` / `-obrys` | `#C9A24B` / `#8A6D2E` / `#FBF7EC` / `#E6C67A` | Złoto jako wypełnienie, tekst, tło plakietki, obrys |
| `--dn-focus` / `-cien` | `#C9A24B` / `rgba(201,162,75,0.35)` | Pierścień fokusu klawiaturowego |

### 2.3. Kolory semantyczne — motyw ciemny

| Token | Wartość | Zastosowanie |
|---|---|---|
| `--dn-tlo` | `#0B1120` | Tło głównego obszaru roboczego |
| `--dn-powierzchnia` | `#16223C` | Powierzchnia kart, pól, paneli |
| `--dn-powierzchnia-2` | `#10192E` | Powierzchnia drugorzędna |
| `--dn-panel` | `#16223C` | Tło modali i paneli bocznych |
| `--dn-hover` | `rgba(255,255,255,0.06)` | Tło stanu najechania |
| `--dn-nakladka` | `rgba(3,6,14,0.62)` | Nakładka pod modalem |
| `--dn-obrys` / `-subtelny` / `-mocny` | `#26395C` / `#1E2E4E` / `#3A4E78` | Trzy poziomy intensywności obrysu |
| `--dn-tekst` | `#E8EDF7` | Tekst podstawowy |
| `--dn-tekst-2` | `#9AA6BF` | Tekst drugorzędny |
| `--dn-tekst-3` | `#6B7E9C` | Tekst trzeciorzędny |
| `--dn-tekst-inv` | `#0B1120` | Tekst na powierzchniach jasnych (złoto) |
| `--dn-marka` / `-hover` | `#26395C` / `#3A4E78` | Kolor marki w trybie ciemnym (jaśniejszy niż w jasnym, dla odróżnienia od tła) |
| `--dn-akcent` / `-txt` / `-tlo` / `-obrys` | `#D9B968` / `#E6C67A` / `rgba(201,162,75,0.12)` / `#A6812F` | Złoto podniesione o jeden krok jasności dla kontrastu na ciemnym tle |
| `--dn-focus` / `-cien` | `#D9B968` / `rgba(217,185,104,0.40)` | Pierścień fokusu klawiaturowego |

### 2.4. Kolory statusów

Każdy status ma trzy warianty tokenu: `-fg` (tekst), `-bg` (tło plakietki/alertu), `-br` (obrys) — dodatkowo wariant `-dk` używany na powierzchniach ciemnych niezależnie od aktywnego motywu (np. pasek marki, toast).

| Status | Jasny — tekst | Jasny — tło | Jasny — obrys | Ciemny — tekst |
|---|---|---|---|---|
| Sukces | `#1F7A50` | `#E7F4EE` | `#B9E0CE` | `#4CAF7D` |
| Ostrzeżenie | `#8A6212` | `#FBF1DA` | `#ECD8A6` | `#E0A63C` |
| Błąd | `#B23A34` | `#FBECEB` | `#F0C9C6` | `#E06A66` |
| Informacja | `#2E5C94` | `#E9F1FA` | `#C6DBF0` | `#5B9BD8` |

Statusy nigdy nie polegają wyłącznie na kolorze: łączone są z ikoną i/lub etykietą tekstową (plakietka z kropką i tekstem, alert z ikoną — rozdz. 8.4, 8.12).

### 2.5. Kontrast i dostępność

Wszystkie pary tekst/tło systemu spełniają standard WCAG 2.1 AA (próg ≥ 4,5:1 dla tekstu normalnego, ≥ 3:1 dla tekstu dużego i elementów graficznych), zweryfikowany metodą kontrastu względnej luminancji sRGB. Poniżej zestawienie skrócone dla par najczęściej używanych w komponentach.

| Para | Motyw | Kontrast | Wynik |
|---|---|---|---|
| Tekst podstawowy / tło | Jasny | 16,66 | AA |
| Tekst podstawowy / tło | Ciemny | 16,04 | AA |
| Tekst-2 / tło | Jasny | 6,93 | AA |
| Tekst-2 / tło | Ciemny | 7,69 | AA |
| Tekst-3 / tło | Jasny | 4,06 | AA — tylko tekst duży lub drugorzędny (zob. uwaga niżej) |
| Złoto-tekst / powierzchnia | Jasny | 4,88 | AA |
| Złoto / powierzchnia | Ciemny | 8,35 | AA |
| Biel / przycisk granatowy | Jasny | 14,25 | AA |
| Granat / przycisk złoty | Jasny | 6,63 | AA |

**Uwaga o `tekst-3`:** na jasnym tle kontrast wynosi 4,06:1, poniżej progu AA dla drobnego tekstu. Token jest przeznaczony wyłącznie dla tekstu dużego (≥ 18,66 px przy grubości 700, albo ≥ 24 px) lub dla treści drugorzędnej niebędącej nośnikiem informacji krytycznej — metadanych, znaczników czasu, etykiet pomocniczych. Do treści czytelnej stosuje się `--dn-tekst-2` lub `--dn-tekst`. Ta sama zasada obowiązuje we wszystkich zastosowaniach opisanych w rozdziałach 10–13.

Dodatkowe reguły dostępności: fokus klawiaturowy jest zawsze widoczny (pierścień złoty, rozdz. 6); `color-scheme` przełącza się per motyw dla poprawnego wyglądu natywnych kontrolek przeglądarki/webview.

---

## 3. Typografia

### 3.1. Kroje i role

| Rola | Krój | Token | Zastosowanie |
|---|---|---|---|
| Nagłówki | Cormorant Garamond (serif) | `--dn-ff-naglowek` | Tytuły kart, modali, godło, tytuły kart środowisk strony głównej (rozdz. 11.1) |
| Interfejs / treść | Inter, awaryjnie Segoe UI, system-ui | `--dn-ff-bazowa` | Cały tekst interfejsu, treść czatu, etykiety, formularze |
| Dane techniczne / kod | JetBrains Mono | `--dn-ff-mono` | Blok kodu, dane techniczne, identyfikatory w oknach programistycznych (rozdz. 10.3) |

Serif jest zarezerwowany wyłącznie dla nagłówków, tytułów kart i modali oraz godła marki — nigdy dla treści ciągłej. To rozróżnienie ma bezpośrednie znaczenie dla strefy 1 strony głównej (rozdz. 11.1), gdzie tytuł każdej z czterech kart środowisk jest złożony krojem nagłówkowym zgodnie z ustaleniem rozdz. 7.4 Koncepcji.

### 3.2. Skala rozmiarów

| Stopień | Rozmiar (px) | Typowe zastosowanie |
|---|---|---|
| `xs` | 11 | Etykiety wersalikowe, nagłówki kolumn tabeli, etykiety sekcji |
| `sm` | 12 | Tekst pomocniczy, metadane, opis trybu pracy karty środowiska |
| `base` | 14 | Tekst podstawowy interfejsu i treści czatu |
| `md` | 15 | Wyróżniony tekst treści |
| `lg` | 17 | Mniejszy stopień skali nagłówkowej |
| `xl` | 21 | Skala nagłówkowa |
| `2xl` | 26 | Tytuły kart środowisk strony głównej (wariant mniejszy) |
| `3xl` | 34 | Tytuły kart środowisk strony głównej (wariant większy) |
| `display` | 44 | Największy stopień ekspozycyjny |

### 3.3. Grubości, interlinie i tracking

| Grubość | Wartość |
|---|---|
| `normal` | 400 |
| `medium` | 500 |
| `semibold` | 600 |
| `bold` | 700 |
| `heavy` | 800 |

| Parametr | Wartości |
|---|---|
| Interlinie | `ciasny 1,25` · `bazowy 1,5` · `luzny 1,65` |
| Tracking | nagłówek `0,01em` · wersaliki `0,12em` |

### 3.4. Zasady stosowania

- Tekst podstawowy interfejsu: 14 px (`base`), interlinia 1,5 (`bazowy`).
- Etykiety wersalikowe (np. nagłówki kolumn tabeli, etykiety sekcji): `font-size: 11px`, `letter-spacing: 0,12em`, `text-transform: uppercase`, kolor `--dn-tekst-3`.
- Nagłówki kart i modali: krój nagłówkowy, grubość `semibold`, interlinia `ciasny`.
- Tytuły kart środowisk (strefa 1 strony głównej): krój nagłówkowy w rozmiarze `2xl` lub `3xl`, zgodnie z wagą wizualną strefy głównej (rozdz. 11.1).

---

## 3a. Warstwy widoczności w systemie wizualnym

### 3a.1. Zasada nadrzędna

Interfejs Danaco Console ujawnia funkcjonalność stopniowo: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w warstwie wizualnej do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu ani na liczbę elementów rysowanych w stanie spoczynku.

Każdy element interfejsu opisywany w niniejszym dokumencie należy do dokładnie jednej z czterech warstw widoczności. Przy opisie okna, panelu, paska narzędzi i pojedynczej funkcji podaje się jej warstwę.

### 3a.2. Cztery warstwy

| Warstwa | Nazwa | Zawartość | Sposób dostępu | Tokeny wiodące |
|---|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, Execution Loop Window, aktywne okno wiodące obszaru roboczego, pasek kontekstu, boczna nawigacja modułów, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu. | Widoczna bez interakcji | `--dn-tlo`, `--dn-powierzchnia`, `--dn-tekst`, `--dn-cien-1` |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie | `--dn-powierzchnia-2`, `--dn-tekst-2`, `--dn-akcent-tlo`, `--dn-cien-2` |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana | `--dn-panel`, `--dn-obrys`, `--dn-cien-3` |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów | `--dn-panel`, `--dn-nakladka`, `--dn-cien-lg`, `--dn-ff-mono` |

### 3a.3. Hierarchia wizualna wynikająca z warstw

Warstwa determinuje wagę wizualną elementu — powierzchnię, kontrast, stopień pisma i głębokość cienia.

| Warstwa | Udział powierzchni | Kontrast tekstu | Stopień pisma | Głębokość |
|---|---|---|---|---|
| 1 | Ponad 80% powierzchni interfejsu | `--dn-tekst` (16,66 / 16,04) | `base` 14 px, nagłówki `lg`–`xl` | Płasko, `--dn-cien-1` |
| 2 | Wyzwalacze zajmują pojedyncze wiersze i znaczniki paska kontekstu | `--dn-tekst-2` | `sm` 12 px, etykiety `xs` 11 px wersalikami | `--dn-cien-2` po rozwinięciu |
| 3 | Powierzchnia rozwinięcia nie przekracza powierzchni okna wywołującego | `--dn-tekst-2`, pozycja aktywna `--dn-akcent-txt` | `base` 14 px | `--dn-cien-3`, nad warstwą 1 |
| 4 | Brak stałej powierzchni w stanie spoczynku | `--dn-tekst` na `--dn-panel` | `base` 14 px, dane techniczne `--dn-ff-mono` | `--dn-cien-lg` + `--dn-nakladka` |

Złoto (`--dn-akcent`) występuje wyłącznie na elemencie aktywnym oraz na jednym CTA widoku, niezależnie od warstwy — reguła oszczędnego złota (rozdz. 1.1) obowiązuje we wszystkich czterech warstwach.

### 3a.4. Wzorce wizualne mechanizmów ukrywania

| Mechanizm | Warstwa | Wzorzec wizualny w stanie spoczynku | Stan rozwinięty | Tokeny |
|---|---|---|---|---|
| Element zbiorczy `Operacje ▼` | 3 | Jeden przycisk `.dn-btn--zarys` z etykietą i znakiem `▼` (ikona `strzalka-dol`) w miejscu zestawu przycisków | Lista akcji w panelu popover, pozycje `.dn-karta--pozycja` | `--dn-obrys`, `--dn-tekst-2`; rozwinięcie `--dn-panel` + `--dn-cien-3` |
| Menu progresywne `Agent ▼` | 2 | Jeden element zwinięty z nazwą bieżącego wyboru i znakiem `▼` | Lista wyborów jednorodnych; pozycja bieżąca oznaczona `--dn-akcent-tlo` / `--dn-akcent-txt` | `--dn-powierzchnia-2`, `--dn-akcent-tlo` |
| Menu kebab `⋮` | 3 | `.dn-btn-ikona` 36×36 px z ikoną `wiecej`, bez obrysu, w skrajnej prawej kolumnie wiersza lub nagłówka karty | Menu kontekstowe wyrównane do krawędzi wyzwalacza | `--dn-tekst-3`; najechanie `--dn-hover` |
| Menu hamburger `☰` | 3 | `.dn-btn-ikona` z ikoną `menu` w nagłówku kolumny nawigacji | Rozwinięcie kolumny nawigacji do pełnej listy pozycji | `--dn-powierzchnia-2`, `--dn-obrys-subtelny` |
| Panel popover | 3 | Niewidoczny; wyzwalany elementem warstwy 2 lub 3 | Powierzchnia `--dn-panel`, promień `lg` (14 px), obrys `--dn-obrys`, cień `--dn-cien-3`; szerokość nie przekracza szerokości kolumny wywołującej | `--dn-panel`, `--dn-cien-3` |
| Panel wysuwany | 3 | Niewidoczny; wyzwalany ikoną w pasku kontekstu | Kolumna boczna wysuwana od prawej krawędzi obszaru roboczego, promień `0`, obrys lewy `--dn-obrys`, cień `--dn-cien-lg` | `--dn-panel`, `--dn-cien-lg` |
| Znacznik kontekstowy (tag) | 2 | `.dn-plakietka` w promieniu `pill`, stopień `xs` 11 px, tekst `--dn-tekst-2`, tło `--dn-powierzchnia-2` | Kliknięcie otwiera selektor w panelu popover | `--dn-powierzchnia-2`, `--dn-obrys-subtelny` |
| Pasek kontekstu | 1 | Wiersz znaczników w nagłówku kolumny Chat Window: `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Znaczniki pozostają widoczne; rozwija się wyłącznie wybrany selektor | `--dn-powierzchnia-2`, `--dn-tekst-2` |
| Wyszukiwarka funkcji | 4 | Niewidoczna; wyzwalana skrótem klawiszowym | Pole `.dn-input` z ikoną `szukaj` w panelu nakładkowym wyśrodkowanym nad obszarem roboczym, lista wyników `.dn-karta--pozycja` | `--dn-panel`, `--dn-nakladka`, `--dn-cien-lg` |

Stany wyzwalaczy warstw 2–3 są wspólne i zgodne z rozdz. 6: spoczynek — tekst `--dn-tekst-2` bez tła; najechanie — tło `--dn-hover`; fokus klawiaturowy — pierścień złoty 2 px z odsunięciem 2 px; rozwinięcie aktywne — tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt`, znak `▼` obrócony o 180°; wyłączenie — `opacity: 0.5`.

### 3a.5. Animacja pojawienia i zniknięcia

| Zdarzenie | Czas | Krzywa | Przekształcenie |
|---|---|---|---|
| Pojawienie panelu popover i menu kontekstowego | `--dn-czas-2` 0,18 s | `--dn-ease` | `opacity 0 → 1`, `translateY(-4px) → 0` |
| Zniknięcie panelu popover i menu kontekstowego | `--dn-czas-1` 0,12 s | `--dn-ease` | `opacity 1 → 0`, `translateY(0) → (-4px)` |
| Wysunięcie panelu bocznego | `--dn-czas-3` 0,24 s | `--dn-ease` | zmiana szerokości kolumny od 0 do szerokości docelowej |
| Schowanie panelu bocznego | `--dn-czas-2` 0,18 s | `--dn-ease` | zmiana szerokości kolumny do 0; po zakończeniu panel znika z przestrzeni roboczej całkowicie |
| Obrót znaku `▼` wyzwalacza | `--dn-czas-1` 0,12 s | `--dn-ease` | `rotate(0) → rotate(180deg)` |
| Pojawienie wyszukiwarki funkcji | `--dn-czas-2` 0,18 s | `--dn-ease` | `opacity 0 → 1`, `scale(0.98) → 1`; nakładka `--dn-nakladka` narasta równolegle |

Wszystkie powyższe przejścia respektują `prefers-reduced-motion: reduce` i skracają się wtedy do 0,01 ms (rozdz. 5.2). Zniknięcie elementu warstw 2–3 zwalnia zajmowaną powierzchnię w całości — po zamknięciu panel nie pozostawia rezerwy miejsca ani śladu obrysu w kolumnie.

### 3a.6. Zasada jednego kliknięcia

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione: rozwinięcie warstwy 3 prezentuje pełną listę akcji, bez podmenu drugiego poziomu. Ukrycie zmniejsza chaos wizualny i nie wydłuża drogi dostępu.

### 3a.7. Makiety w stanie spoczynku

Makiety ASCII niniejszego dokumentu rysuje się w stanie spoczynku interfejsu: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki paska kontekstu, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisuje się w tabelach elementów okna z podaniem warstwy i sposobu wywołania.

```
 ═══════════════════════════════════════════════════════════════════════════
  ☰ Boczna │ Chat Window              │ Obszar roboczy modułu     │ Pasek
    nawi-  │ Użytkownik ↔ Wykonawca   │                       ⋮   │ kontekstu
    gacja  │ [Danaco Console][Ubuntu] │                           │ (znaczniki
    modu-  │ [Worktree][Fable 5]      │   okno wiodące            │  warstwy 2)
    łów    │ Agent ▼      Operacje ▼  │   (warstwa 1)             │
           │ ──────────────────────── │                           │
           │ Execution Loop Window    │                           │
           │ Koordynator ↔ Wykonawca  │                           │
 ═══════════════════════════════════════════════════════════════════════════
   warstwa 1: kolumny robocze          warstwa 2: znaczniki i ▼
   warstwa 3: ⋮ ☰ i panele wywoływane  warstwa 4: bez reprezentacji w spoczynku
```

---

## 4. Siatka kolumnowa, odstępy, promienie

### 4.1. Odstępy

Skala odstępów opiera się na jednostce 4 px; wszystkie odstępy w interfejsie są wielokrotnością tej jednostki — nie wprowadza się wartości spoza skali.

| Token | `--dn-od-1` | `-2` | `-3` | `-4` | `-5` | `-6` | `-8` | `-10` | `-12` | `-16` |
|---|---|---|---|---|---|---|---|---|---|---|
| Wartość (px) | 4 | 8 | 12 | 16 | 20 | 24 | 32 | 40 | 48 | 64 |

### 4.2. Promienie zaokrągleń

| Token | `xs` | `sm` | `md` | `lg` | `xl` | `pill` |
|---|---|---|---|---|---|---|
| Wartość (px) | 4 | 8 | 10 | 14 | 20 | 999 |
| Zastosowanie | drobne akcenty, blok kodu | przyciski, pola formularzy | karty pojedyncze, modale (mniejsze) | karty, modale | duże powierzchnie, karty środowisk strony głównej | plakietki, przełączniki |

### 4.3. Siatka kolumnowa interfejsu

Obowiązuje wyłącznie układ pionowy — podział lewa–prawa. Wszystkie okna robocze, okna komunikacji, przestrzenie modułowe i okna pomocnicze rozmieszczone są w kolumnach sąsiadujących poziomo, na pełną wysokość obszaru roboczego. Kanoniczne rozmieszczenie kolumn:

| Kolumna | Zawartość | Waga wizualna |
|---|---|---|
| Boczna nawigacja modułów | Lista modułów środowiska | Kolumna stała, pełna wysokość obszaru roboczego |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca | Kolumna sąsiadująca, otwierana, pełna wysokość obszaru roboczego |
| Obszar roboczy modułu | Okna edycyjne, podglądu, monitory | Kolumna dominująca, największy udział powierzchni |
| Okna pomocnicze i panele | Panele kontekstowe, panele wysuwane | Kolumna boczna, otwierana jako rozszerzenie boczne, po prawej stronie obszaru roboczego |

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │                          │  boczne)
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
```

Wewnątrz każdej kolumny panele i listy stosują odstępy `--dn-od-3` do `--dn-od-5` między elementami, zgodnie z gęstością przyjętą w `components.css` dla kart, tabel i pól formularzy. Kolumny rozdzielone są obrysem `--dn-obrys-subtelny` o grubości 1 px, bez cienia — rozdzielenie kolumn jest wyłącznie liniowe.

### 4.4. Waga wizualna i responsywność kolumn

Regulacji podlega wyłącznie szerokość kolumn; wysokość każdej kolumny jest pełną wysokością obszaru roboczego. Wagę wizualną kolumny wyznacza jej szerokość, kontrast powierzchni i gęstość treści.

| Kolumna | Szerokość bazowa | Zakres regulacji | Powierzchnia |
|---|---|---|---|
| Boczna nawigacja modułów | 240 px | 200–320 px | `--dn-powierzchnia-2` |
| Chat Window | 30% szerokości obszaru roboczego | 24–40% | `--dn-powierzchnia` |
| Execution Loop Window | 24% szerokości obszaru roboczego | 20–34% | `--dn-powierzchnia-2` |
| Obszar roboczy modułu | Pozostała szerokość | Wypełnia resztę, minimum 40% | `--dn-tlo` |
| Panel pomocniczy | 320 px | 280–420 px | `--dn-panel` |

Responsywność realizowana jest wyłącznie przez zwężanie i zwijanie kolumn — kolejność kolumn od lewej do prawej pozostaje niezmienna, a żadna kolumna nie przechodzi pod inną:

| Dostępna szerokość | Zachowanie kolumn |
|---|---|
| Powyżej 1600 px | Wszystkie kolumny w szerokościach bazowych |
| 1280–1600 px | Panel pomocniczy zwinięty do wyzwalacza ikonowego w kolumnie skrajnej prawej; pozostałe kolumny w szerokościach bazowych |
| 1024–1280 px | Boczna nawigacja modułów zwinięta do kolumny ikonowej 56 px, rozwijana wyzwalaczem `☰` |
| 768–1024 px | Execution Loop Window zwinięte do wyzwalacza w nagłówku Chat Window; Chat Window i obszar roboczy zajmują po połowie szerokości |
| Poniżej 768 px | Widoczna jedna kolumna naraz, wybierana przełącznikiem kolumn; kolejność przełączania odpowiada kolejności kolumn od lewej do prawej |

Zwinięcie i rozwinięcie kolumny przebiega jako zmiana szerokości w czasie `--dn-czas-3` (rozdz. 3a.5).

---

## 5. Cienie i ruch

### 5.1. Cienie

| Token | Jasny | Ciemny |
|---|---|---|
| `--dn-cien-1` | `0 1px 2px rgba(18,22,31,0.05)` | `0 1px 2px rgba(0,0,0,0.30)` |
| `--dn-cien-2` | `0 2px 8px rgba(18,22,31,0.06), 0 1px 2px rgba(18,22,31,0.04)` | `0 3px 10px rgba(0,0,0,0.38)` |
| `--dn-cien-3` | `0 8px 24px rgba(18,22,31,0.10)` | `0 10px 30px rgba(0,0,0,0.48)` |
| `--dn-cien-lg` | `0 20px 48px rgba(18,22,31,0.16)` | `0 24px 60px rgba(0,0,0,0.60)` |
| `--dn-cien-zloto` | `0 4px 18px rgba(201,162,75,0.22)` | `0 4px 20px rgba(217,185,104,0.18)` |

W trybie ciemnym cienie są głębsze — przełączenie następuje automatycznie wraz z motywem (rozdz. 9), bez interwencji na poziomie komponentu. Cień złoty jest zarezerwowany dla akcentów CTA i stanu najechania kart (m.in. karty środowisk strony głównej, rozdz. 11.1).

### 5.2. Ruch

| Parametr | Wartość |
|---|---|
| Krzywa | `--dn-ease`: `cubic-bezier(0.4, 0, 0.2, 1)` |
| Czasy | `--dn-czas-1` 0,12 s · `--dn-czas-2` 0,18 s · `--dn-czas-3` 0,24 s |

System respektuje `prefers-reduced-motion: reduce` — wszystkie animacje i przejścia skracają się wtedy do 0,01 ms (reguła zdefiniowana globalnie w `components.css`), bez konieczności osobnej konfiguracji per komponent.

---

## 6. Stany interakcji i fokus

| Stan | Reguła |
|---|---|
| Hover | subtelna zmiana tła (`--dn-hover`) albo uniesienie `translateY(-1px)`/`(-2px)` dla elementów akcentowanych złotem |
| Active | wciśnięcie `translateY(1px)` |
| Focus | pierścień złoty 2 px + odsunięcie 2 px, wyzwalany wyłącznie przez `:focus-visible` (nawigacja klawiaturą), obowiązkowy dla wszystkich kontrolek interaktywnych |
| Disabled | `opacity: 0.5`, kursor `not-allowed`, bez efektów najechania i bez pierścienia fokusu |
| Zaznaczenie / stan aktywny | akcent złoty: podkreślenie zakładki, obrys karty, tło `--dn-akcent-tlo` |

Pierścień fokusu jest zdefiniowany globalnie i nie podlega usunięciu bez zapewnienia równoważnego wskaźnika — dotyczy to również okien operacyjnych modułów i okna konfiguracji punktów izolacji (rozdz. 12), gdzie liczba kontrolek interaktywnych (przełączniki macierzy izolacji) jest szczególnie wysoka.

---

## 7. Ikonografia

### 7.1. Zasady rysunku

| Parametr rysunku | Wartość |
|---|---|
| Siatka | 24×24 |
| Grubość obrysu | 1,7 px |
| Wypełnienie | `fill="none"` |
| Obrys | `stroke="currentColor"` — kolor dziedziczony z koloru tekstu kontekstu |
| Zakończenia i łączenia linii | zaokrąglone |
| Rozmiary renderowania | 16, 18, 24 px; warianty `.dn-ikona--sm` (14 px), `.dn-ikona--lg` (20 px) — rozdz. 8.12 |

Dziedziczenie koloru z kontekstu (`currentColor`) pozwala jednej ikonie automatycznie zmieniać barwę wraz z motywem lub stanem — na przykład na `--dn-akcent-txt` dla stanu aktywnego.

### 7.2. Zestaw i galeria

Zestaw liczy 47 ikon, udokumentowanych w galerii z wyszukiwarką i przełącznikiem motywu pod adresem `ikony/indeks.html`. Pełny katalog nazw wraz z typowym zastosowaniem w oknach operacyjnych platformy przedstawia Załącznik C.

### 7.3. Godło marki

Godło `ikony/logo-danaco.svg` jest jedyną ikoną pełnokolorową w zestawie: tarcza granatowa z monogramem „D" w złocie. Stosowane w kolumnie nawigacji (rozdz. 8.7) oraz — w większym rozmiarze — jako element godła na kartach środowisk strony głównej (rozdz. 11.1), zgodnie z rozdz. 7.4 Koncepcji.

---

## 8. Biblioteka komponentów

Biblioteka `components.css` (prefiks `.dn-`) dostarcza gotowe komponenty interfejsu zgodne z tokenami rozdziałów 2–7. Poniższe podrozdziały porządkują komponenty według rodzaju wraz z ich wariantami i zastosowaniem w platformie; skrótową tabelę klas zawiera Załącznik B.

### 8.1. Przyciski

Klasa bazowa `.dn-btn`, rozmiary `--sm`/`--lg`. Osobna klasa `.dn-btn-ikona` dla przycisków czysto ikonowych (36×36 px).

| Wariant | Wygląd | Zastosowanie (zasada doboru) |
|---|---|---|
| `--glowny` | Wypełnienie granatowe, tekst odwrócony | Akcja podstawowa okna (np. „Zapisz") |
| `--zloty` | Gradient złoty, cień `--dn-cien-zloto` — sygnaturowy akcent CTA | Jeden, najważniejszy CTA na widoku (np. „Uruchom proces" w Execution Monitor) |
| `--zarys` | Obrys, wypełnienie tłem akcentu przy najechaniu | Akcje drugorzędne |
| `--duch` | Bez obrysu, przezroczysty | Akcje drugorzędne |
| `--blad` | Sygnalizacja wagi akcji | Akcje nieodwracalne lub niebezpieczne |

Zgodnie z zasadą pełnej konfigurowalności i braku twardych blokad wariant `--blad` sygnalizuje wagę akcji wizualnie, nie blokuje jej wykonania.

### 8.2. Pola formularzy

| Klasa | Rola |
|---|---|
| `.dn-pole` | Kontener z etykietą `.dn-pole-etykieta`; wariant błędu `--blad` |
| `.dn-input`, `.dn-select`, `.dn-textarea` | Pola wprowadzania — spójny obrys, fokus przez obrys `--dn-focus` + poświatę `--dn-focus-cien` |
| `.dn-suwak` | Przełącznik (toggle) — tor w gradiencie złotym w stanie zaznaczonym |
| `.dn-check` | Checkbox / radio z `accent-color` złotym |
| `.dn-pomoc` | Pomoc kontekstowa |

Komponent `.dn-suwak` jest podstawowym elementem wizualnym macierzy izolacji w oknie konfiguracji punktów izolacji (rozdz. 12.2) — każdy z ośmiu zakresów izolacji technicznej i trzech pozycji izolacji kontekstu jest przełącznikiem tego typu.

### 8.3. Karty

| Wariant | Charakterystyka |
|---|---|
| `.dn-karta` (bazowy) | Płaska wewnątrz panelu (bez własnego obrysu/cienia — panel jest kontenerem) |
| `--interaktywna` | Wskazanie stanu najechania tłem, nie uniesieniem |
| `--akcent` | Wstęga akcentowa 3 px w kolorze złota na krawędzi karty |
| `--pozycja` | Wariant listy pozycji — wiersze rozdzielone cienką kreską dolną zamiast obramowania każdej pozycji osobno |

Struktura karty:

```
.dn-karta [.dn-karta--interaktywna] [.dn-karta--akcent]
├── .dn-karta-naglowek
│      └── .dn-karta-tytul      (krój nagłówkowy --dn-ff-naglowek)
├── .dn-karta-tresc
└── .dn-karta-stopka            (akcje: .dn-btn …)
```

Karty środowisk strony głównej (rozdz. 11.1) są głównym zastosowaniem wariantu `--akcent` i `--interaktywna` łącznie: wstęga i cień złoty pojawiają się na stanie aktywnym/najechania, zgodnie z zasadą oszczędnego złota (rozdz. 1.1, rozdz. 7.4 Koncepcji).

### 8.4. Plakietki i statusy

| Klasa | Warianty | Charakterystyka i zastosowanie |
|---|---|---|
| `.dn-plakietka` | `--zloto`, `--sukces`, `--ostrz`, `--blad`, `--info`, `--wersaliki` | Pigułka z obrysem; `--wersaliki` = etykieta kapitalikowa |
| `.dn-kropka` | `--sukces`, `--ostrz`, `--blad` | Kropka stanu z poświatą koloru tła statusu |
| `.dn-plakietka--stan` | — | Łączy pigułkę z kropką (np. status uruchomienia procesu w Execution Monitor lub Monitorze procesu środowiska MultitaskingAI, rozdz. 13) |

### 8.5. Zakładki

| Wariant | Charakterystyka | Zastosowanie |
|---|---|---|
| `.dn-zakladki` | Wariant podkreślony — dolna kreska złota na zakładce aktywnej | Przełączanie widoków w jednym oknie operacyjnym (np. panele modeli w Roundtable) |
| `.dn-zakladki--pigulki` | Wariant segmentowany — tło subtelne na zakładce aktywnej, bez własnego obramowania grupy | Przełączniki gęstsze, jak wybór warstwy konfiguracji w oknie punktów izolacji (rozdz. 12.3) |

### 8.6. Tabela

| Cecha | Realizacja |
|---|---|
| Nagłówek | Wersalikowy, na tle `--dn-powierzchnia-2` |
| Separacja wierszy | Kreska subtelna |
| Wariant `--paski` | Naprzemienne tło wierszy |
| Stan najechania | Podświetlenie wiersza |
| Owijka | `.dn-tabela-owijka` |

`.dn-tabela` jest podstawowym komponentem dla okien typu monitor i manager (rozdz. 10.5, 10.7) — Execution Monitor, Process Monitor, Queue Manager, Permissions Center.

### 8.7. Kolumna nawigacji

| Element | Klasa | Charakterystyka | Warstwa |
|---|---|---|---|
| Kolumna | `.dn-pasek` w orientacji pionowej | Tło marki (granat), pełna wysokość obszaru roboczego | 1 |
| Blok logotypu | `.dn-pasek-marka` | Krój nagłówkowy, u szczytu kolumny | 1 |
| Pole wyszukiwania | `.dn-pasek-szukaj` | Pole `.dn-input` z ikoną `szukaj` w obrębie kolumny | 1 |
| Wyzwalacz zwinięcia | `.dn-btn-ikona` z ikoną `menu` (`☰`) | Zwija kolumnę do szerokości ikonowej 56 px (rozdz. 4.4) | 3 |
| Przyciski ikonowe | `.dn-btn-ikona` | Wariant na tle ciemnym, ułożone pionowo | 1 |

Kolumna nawigacji stanowi wzorzec bocznej nawigacji modułów okna środowiska (rozdz. 8 Koncepcji), wspólny dla wszystkich czterech środowisk. Zajmuje skrajną lewą kolumnę układu (rozdz. 4.3) i nie dzieli obszaru roboczego w żadnym innym kierunku.

### 8.8. Awatar

| Wariant | Charakterystyka |
|---|---|
| `.dn-awatar` (bazowy) | Koło z gradientem granatowym |
| `--zloto` | Gradient złoty |
| `--sm` / `--lg` | Rozmiary |
| `--kwadrat` | Identyfikacja nie-osobowa (np. agenta) |
| `.dn-awatar-stan` | Wskaźnik statusu online |

Zastosowanie: identyfikacja modeli i ról w oknach wielomodelowych (Model Panels w Roundtable, role w panelu orkiestracji MultitaskingAI — rozdz. 13) oraz identyfikacja agentów utworzonych w module Agents.

### 8.9. Modal

| Element | Rola |
|---|---|
| `.dn-nakladka` | Nakładka pod modalem |
| `.dn-modal` | Okno modalne |
| Nagłówek modala | Tytuł w kroju nagłówkowym |
| Treść modala | Przewijana |
| Stopka modala | Akcje wyrównane do prawej |

Zastosowanie: okno konfiguracji uruchamiane z okna operacyjnego (Architektura, rozdz. 13), potwierdzenia akcji nieodwracalnych, kreatory tworzenia komponentu własnego w strefie 2 strony głównej (rozdz. 11.2).

### 8.10. Toast

| Cecha | Realizacja |
|---|---|
| Tło | Granatowe ciemne |
| Kreska lewa | W kolorze statusu (`--sukces` / `--blad` / `--ostrz`) |

Zastosowanie: potwierdzenia zapisu profilu izolacji, zakończenia procesu automatyki, komunikatów funkcji globalnej Mobile (rozdz. 12.1 Koncepcji).

### 8.11. Tooltip

`.dn-tooltip` — dymek na tle granatowym ciemnym, wyzwalany najechaniem lub fokusem. Jest nośnikiem objaśnień kontekstowych `[?]` przy każdym elemencie konfiguracji, zgodnie z zasadą przyjętą w Architekturze (rozdz. 13): „każdy element konfiguracji zawiera objaśnienie kontekstowe opisujące jego działanie i wpływ na aplikację". Ten wzorzec jest bezpośrednio zastosowany w macierzy izolacji okna konfiguracji punktów izolacji (rozdz. 12.2).

### 8.12. Dodatki

| Klasa | Rola | Zastosowanie / warianty |
|---|---|---|
| `.dn-spinner` | Wskaźnik ładowania, obrót ciągły | Respektuje `prefers-reduced-motion` |
| `.dn-separator` | Pozioma kreska rozdzielająca sekcje | — |
| `.dn-kod` | Blok treści technicznej (polecenie, ścieżka, ładunek narzędzia) krojem mono na powierzchni drugorzędnej | Terminal, Developer, Diagnostics; podgląd polityki efektywnej okna punktów izolacji (rozdz. 12.3) |
| `.dn-alert` | Komunikat blokowy z kreską lewą semantyczną | Warianty `--info` / `--sukces` / `--ostrz` / `--blad`; lekki, bez własnej powierzchni ani obrysu |
| `.dn-pusty-stan` | Komunikat pustego widoku (ikona + tytuł + opis) | Panele list bez zawartości: Sources Manager, Findings Panel, Queue Manager przed pierwszą konfiguracją |
| `.dn-ikona` | Nośnik ikon SVG maskowanych kolorem `currentColor` (rozdz. 7.1) | Warianty `--sm` / `--lg` |

---

## 9. Motyw jasny i ciemny — mechanizm przełączania

Motyw jest sterowany atrybutem `data-theme="light|dark"` na elemencie głównym dokumentu. Mechanizm przełączania przebiega następująco:

```
Element główny dokumentu:   data-theme="light"   albo   data-theme="dark"
        │
        ▼
tokens.css / tokens.json  —  każdy token semantyczny (rozdz. 2.2–2.4)
        │                     ma zdefiniowaną wartość dla OBU motywów
        ▼
Komponenty (.dn-*)  —  odwołują się wyłącznie do tokenów semantycznych,
        │              nigdy do wartości prymitywnych
        ▼
Przełączenie motywu = zmiana wartości zmiennych, nie reguł komponentów
        │              przejście animowane (--dn-czas-2: kolor tła i tekstu)
        ▼              respektuje prefers-reduced-motion
Spójny wygląd w obu motywach bez zmiany reguł komponentów
```

Struktura definicji tokenów per motyw (fragment; pełny zestaw wartości — rozdz. 2.2 i 2.3):

```
:root[data-theme="light"] {
  --dn-tlo:          #F7F8FA;
  --dn-powierzchnia: #FFFFFF;
  --dn-tekst:        #12161F;
  --dn-akcent:       #C9A24B;
  /* … pełny zestaw tokenów semantycznych, rozdz. 2.2 */
}
:root[data-theme="dark"] {
  --dn-tlo:          #0B1120;
  --dn-powierzchnia: #16223C;
  --dn-tekst:        #E8EDF7;
  --dn-akcent:       #D9B968;
  /* … pełny zestaw tokenów semantycznych, rozdz. 2.3 */
}
```

Wybór motywu jest ustawieniem użytkownika dostępnym z okna konfiguracji (Architektura, rozdz. 13 — warstwa domyślna i warstwa sesji), zgodnie z zasadą pełnej konfigurowalności (Koncepcja, rozdz. 14, zasada 2):

| Źródło wyboru motywu | Zachowanie |
|---|---|
| Brak jawnego ustawienia | Wartość domyślna z preferencji systemowej (`prefers-color-scheme`) |
| Ustawienie jawne | Z warstwy domyślnej lub warstwy sesji, analogicznie do mechanizmu opisanego dla punktów izolacji w rozdz. 6.6 Koncepcji |

---

## 10. Zastosowanie w oknach operacyjnych modułów

### 10.1. Typologia okien operacyjnych

Załącznik A Koncepcji wymienia okna operacyjne wszystkich piętnastu modułów. Okna te, mimo różnorodności nazw, powtarzają ograniczoną liczbę wzorców wizualnych — poniższa typologia porządkuje je względem komponentów opisanych w rozdziale 8, tak aby jeden zestaw tokenów i komponentów obsłużył całą macierz okien operacyjnych (Koncepcja, Załącznik A) bez potrzeby definiowania odrębnego stylu dla każdego z nich.

| Typologia | Podrozdział | Przykładowe okna (rozdz. 11 Koncepcji) |
|---|---|---|
| Chat Window | 10.2 | wspólne dla wszystkich piętnastu modułów |
| Execution Loop Window | 10.2 | wspólne dla wszystkich piętnastu modułów |
| Okna edycyjne | 10.3 | Studio Editor, Code Editor, Design Board, Source Panel / Translation Panels |
| Okna list i źródeł | 10.4 | Sources Panel, Sources Manager, Assets Panel, Library Explorer, Project Tree |
| Okna monitorów i dashboardów | 10.5 | Execution Monitor, Process Monitor, Diagnostics Center, Build Output, Actions Monitor |
| Okna kreatorów (builderów) | 10.6 | Workflow Builder, Agent Builder, Product Builder, Architecture Designer, Prompt Builder |
| Okna zarządców (managerów) | 10.7 | Queue Manager, Orchestrator, Connectors Manager, Permissions Center, Agent Manager, Glossary Manager |

### 10.2. Okna komunikacji — Chat Window i Execution Loop Window

Platforma prowadzi dwa kanały komunikacji operacyjnej, prezentowane w dwóch oknach warstwy 1 (rozdz. 3a.2), rozmieszczonych w sąsiadujących kolumnach po lewej stronie obszaru roboczego. Oba okna są elementami pierwszoplanowymi architektury wizualnej i występują w każdym module i w każdym środowisku w tym samym miejscu układu.

```
 ═══════════════════════════════════════════════════════════════════
  Chat Window                    │ Execution Loop Window
  Użytkownik ↔ Wykonawca         │ Koordynator ↔ Wykonawca
  ─────────────────────────────  │ ────────────────────────────────
  [Danaco Console][Ubuntu]       │ Zlecenie: „…"          ⋮
  [Worktree][Fable 5][Ultra]     │ Zadania:  ● ● ○ ○
  ─────────────────────────────  │ ────────────────────────────────
  ▸ wiadomość Użytkownika        │ ▸ komunikat sterujący
  ▸ odpowiedź Wykonawcy          │ ▸ wynik kontroli jakości
  ─────────────────────────────  │ ────────────────────────────────
  [ pole wprowadzania ]  Agent ▼ │ [wstrzymaj][wznów][przerwij]
 ═══════════════════════════════════════════════════════════════════
   kolumna lewa, stała,            kolumna sąsiadująca, otwierana,
   pełna wysokość                  pełna wysokość
```

**Chat Window — kanał Użytkownik ↔ Wykonawca.** Główne okno komunikacji i centralny punkt pracy użytkownika: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu.

| Aspekt | Realizacja |
|---|---|
| Warstwa widoczności | 1 — zawsze widoczna |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego; szerokość bazowa 30% obszaru roboczego (rozdz. 4.4) |
| Powierzchnia | `--dn-powierzchnia`, obrys prawy `--dn-obrys-subtelny`, bez cienia własnego |
| Komponenty bazowe | `.dn-karta--pozycja` (wiersz wiadomości), `.dn-awatar` (identyfikacja uczestnika: Użytkownik, Wykonawca), `.dn-kod` (fragmenty kodu cytowane w rozmowie), `.dn-input` / `.dn-textarea` (pole wprowadzania), `.dn-plakietka` (znaczniki paska kontekstu) |
| Typografia | Treść wiadomości `--dn-ff-bazowa`, stopień `base` 14 px, interlinia `bazowy` 1,5, grubość `normal`; autor wiadomości grubość `semibold`, stopień `sm` 12 px; znacznik czasu `--dn-tekst-3`, stopień `xs` 11 px; dane techniczne `--dn-ff-mono` |
| Hierarchia wizualna | Strumień rozmowy jest treścią o najwyższym kontraście w kolumnie (`--dn-tekst`); pasek kontekstu i pole wprowadzania mają kontrast obniżony (`--dn-tekst-2` na `--dn-powierzchnia-2`) |
| Akcent | `--dn-akcent` wyłącznie na przycisku wysłania (`.dn-btn--zloty`, ikona `wyslij`) jako jedynym CTA kolumny |
| Elementy warstw 2–3 | Pasek kontekstu ze znacznikami `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` (warstwa 2), menu progresywne `Agent ▼` i `Operacje ▼` (warstwy 2–3), menu kebab `⋮` wiersza wiadomości (warstwa 3) |
| Zakres | Okno wspólne dla wszystkich piętnastu modułów (Koncepcja, rozdz. 2.4, Załącznik A) |
| Zasada | Ten sam zestaw tokenów we wszystkich piętnastu wystąpieniach; zmienia się kontekst i dostępne narzędzia, nie wygląd (rozdz. 2.4 Koncepcji) |

**Execution Loop Window — kanał Koordynator ↔ Wykonawca.** Okno pętli wykonawczej prezentujące bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korektę zlecenia.

| Aspekt | Realizacja |
|---|---|
| Warstwa widoczności | 1 — zawsze widoczna |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego; szerokość bazowa 24% obszaru roboczego (rozdz. 4.4) |
| Powierzchnia | `--dn-powierzchnia-2`, obrysy boczne `--dn-obrys-subtelny` — powierzchnia drugorzędna odróżnia kanał sterujący od kanału użytkownika bez obniżania jego rangi |
| Komponenty bazowe | `.dn-tabela` (kolejka i stan zadań), `.dn-plakietka--stan` + `.dn-kropka` (stan zadania i wynik kontroli jakości), `.dn-karta--pozycja` (dekompozycja zlecenia), `.dn-kod` (komunikaty sterujące i ładunki narzędzi), `.dn-btn-ikona` (`uruchom`, `zatrzymaj`, `odswiez`), `.dn-spinner` (przebieg pętli) |
| Typografia | Nagłówek zlecenia `--dn-ff-bazowa`, grubość `semibold`, stopień `md` 15 px; nazwy zadań stopień `base` 14 px; komunikaty sterujące i identyfikatory `--dn-ff-mono`, stopień `sm` 12 px; etykiety kolumn `xs` 11 px wersalikami, tracking `0,12em`, kolor `--dn-tekst-3` |
| Hierarchia wizualna | Bieżące zlecenie ma najwyższy kontrast kolumny; kolejka zadań kontrast pośredni (`--dn-tekst-2`); zapis komunikatów sterujących kontrast najniższy, krojem mono. Zadanie aktywne wyróżnione tłem `--dn-akcent-tlo` i tekstem `--dn-akcent-txt` |
| Akcent | `--dn-akcent` wyłącznie na zadaniu aktywnym i na wskaźniku przebiegu pętli; sterowanie przebiegiem korzysta z `.dn-btn--zarys`, przerwanie z `.dn-btn--blad` |
| Elementy warstw 2–3 | Korekta zlecenia i warianty ponowienia w elemencie zbiorczym `Operacje ▼` (warstwa 3), menu kebab `⋮` wiersza zadania (warstwa 3), wybór poziomu szczegółowości zapisu jako znacznik kontekstowy (warstwa 2) |
| Zakres | Okno wspólne dla wszystkich piętnastu modułów oraz dla Monitora procesu środowiska MultitaskingAI (rozdz. 13) |

Relacja hierarchiczna obu okien: Chat Window jest oknem wiodącym pary — ma większą szerokość bazową, jaśniejszą powierzchnię i wyższy kontrast treści. Execution Loop Window jest oknem równorzędnym rangą architektoniczną, a odróżnionym wizualnie powierzchnią drugorzędną i gęstszą, techniczną typografią. Oba okna zachowują pełną wysokość obszaru roboczego; regulacji podlega wyłącznie ich szerokość.

### 10.3. Okna edycyjne

| Aspekt | Realizacja |
|---|---|
| Komponenty bazowe | Pasek narzędzi; `--dn-ff-bazowa` (treść dokumentowa), `--dn-ff-mono` + `.dn-kod` (kod), `.dn-karta` w układzie siatki (podgląd zasobów) |
| Okna | Studio Editor, Code Editor, Design Board; cytowania kodu poza edytorem w Diff/Grep Panel i Build Output |
| Uwaga | Treść tekstowa i dokumentowa (Studio Editor) krojem bazowym; kod (Code Editor) krojem mono; Design Board opiera ramę na wspólnych tokenach odstępów i obrysów |

### 10.4. Okna list i źródeł

| Aspekt | Realizacja |
|---|---|
| Komponenty bazowe | `.dn-karta--pozycja`: `.dn-ikona` + `.dn-karta-poz-nazwa` + `.dn-karta-poz-meta` + `.dn-karta-poz-akcje`; `.dn-input`; `.dn-pusty-stan` |
| Okna | Sources Panel, Sources Manager, Assets Panel, Library Explorer, Project Tree, Project Library |
| Uwaga | Metadane tokenem `--dn-tekst-3` (rozdz. 2.5); filtrowanie i wyszukiwanie ikonami `filtr` i `szukaj` w polu `.dn-input`; pusty stan list przez `.dn-pusty-stan` |

### 10.5. Okna monitorów i dashboardów

| Aspekt | Realizacja |
|---|---|
| Komponenty bazowe | `.dn-tabela` (lista przebiegów/zdarzeń/procesów) + `.dn-plakietka--stan` + `.dn-kropka`; `.dn-btn-ikona`; `.dn-kod` |
| Okna | Execution Monitor, Process Monitor, Diagnostics Center, Build Output, Actions Monitor, Monitor procesu środowiska MultitaskingAI (rozdz. 13); logi i komunikaty: Logs Viewer, Errors Panel |
| Uwaga | Status każdego wiersza sukces/ostrzeżenie/błąd (rozdz. 2.4); akcje sterujące procesem `uruchom`/`zatrzymaj`/`odswiez` w `.dn-btn-ikona`; logi i komunikaty techniczne w `.dn-kod` |

### 10.6. Okna kreatorów (builderów)

| Aspekt | Realizacja |
|---|---|
| Komponenty bazowe | `.dn-zakladki` lub kroki formularza (`.dn-pole` w układzie pionowym); `.dn-karta` boczna (podsumowanie bieżącej konfiguracji); `.dn-btn--zloty` |
| Okna | Workflow Builder, Agent Builder, Product Builder, Architecture Designer, Prompt Builder |
| Uwaga | Akcja finalizująca kreator (utworzenie automatyki, agenta, architektury) korzysta z `.dn-btn--zloty` jako jedynego CTA widoku, zgodnie z zasadą oszczędnego złota (rozdz. 8.1) |

### 10.7. Okna zarządców (managerów)

| Aspekt | Realizacja |
|---|---|
| Komponenty bazowe | `.dn-tabela` z kolumną akcji (`.dn-btn-ikona`) lub z przełącznikiem `.dn-suwak` w wierszu, gdy zarządzana właściwość jest binarna |
| Okna | Queue Manager, Orchestrator, Connectors Manager, Permissions Center, Agent Manager, Glossary Manager |
| Uwaga | Przełącznik w wierszu dla właściwości binarnej (uprawnienie w Permissions Center, powiązanie w Connectors Manager); wzorzec przełączników w tabeli rozwinięty w oknie punktów izolacji (rozdz. 12.2), gdzie skala przełączników jest największa w całej platformie |

---

## 11. Zastosowanie na stronie głównej — Centrum dowodzenia

Rozdział 7.4 Koncepcji określa formę prezentacji trzech stref strony głównej jako trzy zróżnicowane formy dobrane do wagi i funkcji każdej strefy. Strefy rozmieszczone są w trzech kolumnach sąsiadujących poziomo, na pełną wysokość obszaru roboczego; waga strefy wyrażona jest szerokością kolumny, krojem tytułu i intensywnością akcentu, nie położeniem względem krawędzi ekranu.

```
CENTRUM DOWODZENIA — strona główna (układ kolumnowy, stan spoczynku)
 ═══════════════════════════════════════════════════════════════════════════
  STREFA 1 · środowiska   │ STREFA 2 · komponenty  │ STREFA 3 · ustawienia
  (waga główna)           │ (waga pośrednia)       │ (waga najniższa)
  ──────────────────────  │ ─────────────────────  │ ─────────────────────
  ┌────────────────────┐  │ ┌───────────────────┐  │  ustawienia
  │ TalkIn             │  │ │ Automations       │  │  Mobile
  └────────────────────┘  │ └───────────────────┘  │  Always On Display
  ┌────────────────────┐  │ ┌───────────────────┐  │
  │ WorkSpace          │  │ │ Agents            │  │  Operacje ▼
  └────────────────────┘  │ └───────────────────┘  │
  ┌────────────────────┐  │ ┌───────────────────┐  │
  │ CodeStudio         │  │ │ Workspace         │  │
  └────────────────────┘  │ └───────────────────┘  │
  ┌────────────────────┐  │ ┌───────────────────┐  │
  │ MultitaskingAI     │  │ │ Assistant         │  │
  └────────────────────┘  │ └───────────────────┘  │
                          │                        │
  .dn-karta               │ .dn-karta              │ .dn-karta--pozycja
  --interaktywna --akcent │ (bez --akcent)         │ (kolumna narzędziowa)
 ═══════════════════════════════════════════════════════════════════════════
```

| Strefa | Kolumna | Komponent bazowy | Warianty i tokeny | Krój tytułu / etykiety | Cień i stan |
|---|---|---|---|---|---|
| Strefa 1 — karty środowisk | Kolumna lewa, dominująca, pełna wysokość | `.dn-karta` | `--interaktywna` + `--akcent` (wstęga 3 px, `--dn-akcent`) aktywny wyłącznie na stanie aktywnym i najechania | Nagłówkowy (`--dn-ff-naglowek`, rozmiar `2xl`/`3xl`) | `--dn-cien-zloto`, uniesienie `translateY(-2px)` |
| Strefa 2 — kafle komponentów | Kolumna środkowa, węższa, pełna wysokość | `.dn-karta` | Bez `--akcent` i bez cienia złotego na spoczynku — wizualnie lżejsze | Bazowy, grubość `semibold` (nie nagłówkowy) | Wejście otwiera `.dn-modal` lub okno konfiguracji komponentu |
| Strefa 3 — kolumna ustawień | Kolumna prawa, najwęższa, pełna wysokość | `.dn-karta--pozycja` | Tło `--dn-powierzchnia-2`; elementy `.dn-btn-ikona` z etykietą | Bez kroju nagłówkowego | Akcent złoty wyłącznie w stanie fokusu/najechania pojedynczego elementu |

Szerokości kolumn stref: strefa 1 — 44% obszaru roboczego, strefa 2 — 34%, strefa 3 — 22%; regulacji podlega wyłącznie szerokość, kolejność kolumn od lewej do prawej pozostaje niezmienna. Zachowanie przy zmniejszaniu dostępnej szerokości opisuje rozdz. 4.4.

### 11.1. Strefa 1 — karty środowisk

Cztery karty środowisk — TalkIn, WorkSpace, CodeStudio, MultitaskingAI — są zbudowane na komponencie `.dn-karta` w wariancie `--interaktywna`, z wariantem `--akcent` aktywowanym wyłącznie na stanie aktywnym i najechania (rozdz. 7.4 Koncepcji: „akcent złoty rezerwowany jest dla stanu aktywnego i najechania"). Tytuł karty jest złożony krojem nagłówkowym (rozmiar `2xl`/`3xl`, rozdz. 3.1), opis trybu pracy — krojem bazowym (rozmiar `sm`/`base`). Godło środowiska korzysta z maski ikony `.dn-ikona--lg` lub, dla marki platformy jako całości, z `logo-danaco.svg`. Cztery karty ułożone są jedna pod drugą w obrębie kolumny strefy 1, z równą wagą wizualną, stanowiąc środek ciężkości strony — zgodnie z metaforą centrum dowodzenia przyjętą w rozdz. 7.4 Koncepcji.

### 11.2. Strefa 2 — kafle komponentów własnych

Cztery kafle — Automations, Agents, Workspace, Assistant — są zbudowane na tym samym komponencie `.dn-karta`, lecz bez wariantu `--akcent` i bez cienia złotego na spoczynku, przez co są wizualnie lżejsze od kart strefy 1 (rozdz. 7.4 Koncepcji). Etykieta kafla jest zorientowana na działanie twórcze (np. „Utwórz automatykę", „Nowy agent") i korzysta z kroju bazowego o grubości `semibold`, nie z kroju nagłówkowego — odróżnienie krojów między strefą 1 a strefą 2 oddziela w odbiorze „wejście do środowiska" od „zbudowania komponentu". Wejście w kafel otwiera modal (`.dn-modal`, rozdz. 8.9) lub okno konfiguracji komponentu własnego (Koncepcja, rozdz. 7.2).

### 11.3. Strefa 3 — kolumna ustawień

Kolumna ustawień (okno konfiguracji, Mobile, Always On Display) jest zbudowana na wzorcu listy pionowej `.dn-karta--pozycja` w wersji lekkiej: tło `--dn-powierzchnia-2` zamiast tła marki, elementy jako `.dn-btn-ikona` z etykietą tekstową obok ikony, jedna pozycja na wiersz.

| Element kolumny | Ikona (Załącznik C) | Warstwa widoczności |
|---|---|---|
| Okno konfiguracji | `ustawienia` | 1 |
| Always On Display | `dzwonek` / awatar | 1 |
| Mobile | ikona urządzenia mobilnego | 1 |
| Ustawienia zaawansowane | `wiecej` w elemencie zbiorczym `Operacje ▼` | 3 |

Kolumna nie stosuje kroju nagłówkowego ani akcentu złotego poza stanem fokusu/najechania pojedynczego elementu — jej waga wizualna jest najniższa spośród trzech stref, sygnalizując warstwę narzędziową, stale dostępną (rozdz. 7.4 Koncepcji).

---

## 12. Zastosowanie w oknie konfiguracji punktów izolacji

Rozdział 6.4 Koncepcji określa układ okna konfiguracji punktów izolacji jako trzypanelowy: selektor zasięgu, macierz izolacji, profil i podgląd. Poniższy schemat mapuje ten układ na komponenty niniejszego systemu wizualnego, a dalsze podrozdziały opisują każdy panel.

```
OKNO KONFIGURACJI PUNKTÓW IZOLACJI  (układ trzypanelowy — rozdz. 6.4 Koncepcji)
┌────────────────────┬──────────────────────────────┬────────────────────┐
│ SELEKTOR ZASIĘGU    │ MACIERZ IZOLACJI             │ PROFIL I PODGLĄD    │
│ (lewy · 12.1)       │ (środkowy · 12.2)           │ (prawy · 12.3)      │
│                     │                              │                     │
│ .dn-karta--pozycja  │ Izolacja kontekstu (3):      │ .dn-btn--zarys:     │
│  • Globalny         │   historia  [.dn-suwak]      │  Zapisz / Wczytaj / │
│  • Środowisko       │   pamięć    [.dn-suwak]      │  Przypisz profil    │
│  • Moduł            │   kontekst  [.dn-suwak]      │ (.dn-btn--zloty dla │
│  • Para modułów     │                              │  „Zapisz profil")   │
│  • Projekt          │ Izolacja techniczna (8):     │                     │
│  • Karta sesji      │   [.dn-suwak] × 8 + ikony    │ Warstwa: .dn-check  │
│  • Rola             │   + [?] .dn-tooltip          │  ( ) domyślna       │
│                     │                              │  ( ) sesji          │
│ zaznaczenie:        │ układ: .dn-tabela lub        │                     │
│  --dn-akcent-tlo    │ .dn-karta--pozycja           │ Polityka efektywna: │
│  --dn-akcent-txt    │ (etykieta | przełącznik)     │  .dn-kod            │
└────────────────────┴──────────────────────────────┴────────────────────┘
```

### 12.1. Panel selektora zasięgu (lewy)

| Aspekt | Realizacja |
|---|---|
| Komponent | `.dn-karta--pozycja` bez akcji w wierszu |
| Zawartość | Siedem poziomów zasięgu: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola (rozdz. 6.5 Koncepcji) |
| Zaznaczenie pozycji | Tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt` — analogicznie do stanu aktywnego zakładki (`.dn-zakladka[aria-selected="true"]`, rozdz. 8.5) |
| Kolejność listy | Odpowiada porządkowi pierwszeństwa zasięgu (rosnąco ku dołowi, zgodnie z tabelą rozdz. 6.5 Koncepcji), co wizualnie odwzorowuje regułę pierwszeństwa zasięgu najbardziej szczegółowego |

### 12.2. Panel macierzy izolacji (środkowy)

Macierz izolacji dzieli się na dwie grupy przełączników `.dn-suwak` (rozdz. 8.2): trzy pozycje izolacji kontekstu i osiem pozycji izolacji technicznej.

| Izolacja kontekstu | Stan przełącznika `.dn-suwak` |
|---|---|
| historia | współdzielone / odrębne |
| pamięć | współdzielone / odrębne |
| kontekst | współdzielone / odrębne |

Osiem zakresów izolacji technicznej (Architektura, rozdz. 12) korzysta dodatkowo z ikon kontekstowych ułatwiających skanowanie wzrokowe listy; zestaw ikon dobrano z katalogu Załącznika C, bez wprowadzania nowych ikon poza istniejącym zestawem 47.

| Zakres izolacji technicznej | Ikona kontekstowa (Załącznik C) | Stan wyjściowy |
|---|---|---|
| katalog roboczy sesji | `folder` | wyłączony |
| środowisko procesu | `ustawienia` | wyłączony |
| katalog danych i konfiguracji modelu | `plik` | wyłączony |
| dostęp sieciowy | `globus` | wyłączony |
| zakres odczytu i zapisu plików | `dokument` | wyłączony |
| konto i token per sesja | `klodka` | wyłączony |
| model procesu | `kod` | wyłączony |
| serwer wykonania | `uruchom` | wyłączony |

Każdy wiersz macierzy stosuje układ tabelaryczny (`.dn-tabela` bez widocznego nagłówka kolumnowego albo `.dn-karta--pozycja` z przełącznikiem jako akcją) z etykietą pozycji po lewej, przełącznikiem po prawej i ikoną objaśnienia kontekstowego `[?]` w `.dn-tooltip` (rozdz. 8.11), zgodnie z zasadą Architektury (rozdz. 13): każdy element konfiguracji niesie objaśnienie swojego działania i wpływu na aplikację.

Zgodnie z zasadą braku twardych blokad (Koncepcja, rozdz. 6.4, rozdz. 14) żaden przełącznik macierzy nie jest zablokowany do edycji ani nie wymusza stanu włączonego — stan wyjściowy wszystkich ośmiu przełączników izolacji technicznej jest „wyłączony" (rozdz. 6.6 Koncepcji), co odwzorowuje się wizualnie pozycją toru przełącznika po lewej stronie (`.dn-suwak` bez `input:checked`).

### 12.3. Panel profilu i podglądu (prawy)

| Element | Komponent | Uwaga |
|---|---|---|
| Zapis, wczytanie i przypisanie profilu izolacji | `.dn-btn--zarys` | Akcje drugorzędne wobec samej konfiguracji macierzy |
| Akcja „Zapisz profil" | `.dn-btn--zloty` | Jako jedyne CTA panelu może przyjąć wariant złoty |
| Wybór warstwy konfiguracji (domyślna / sesji) | `.dn-check` (układ radio) | rozdz. 6.6 Koncepcji, rozdz. 14 Architektury |
| Podgląd polityki efektywnej | `.dn-kod` | Wynikowy zestaw reguł po uwzględnieniu dziedziczenia między poziomami zasięgu |

Podgląd polityki efektywnej jest prezentowany w `.dn-kod` (rozdz. 8.12) z formatowaniem tekstowym analogicznym do schematów stosowanych w Koncepcji i w Architekturze. Przykładowy format podglądu:

```
POLITYKA EFEKTYWNA — poziom: projekt „Alfa"        (podgląd w .dn-kod)
  izolacja kontekstu:
    historia   : odrębna       (poziom: projekt)
    pamięć     : odrębna       (poziom: projekt)
    kontekst   : odrębny       (poziom: projekt)
  izolacja techniczna:
    zakres odczytu i zapisu plików : współdzielony (poziom: projekt)
    katalog roboczy sesji          : wyłączony     (dziedziczone: globalny)
    dostęp sieciowy                : wyłączony     (dziedziczone: globalny)
    … pozostałe zakresy            : wyłączone     (dziedziczone: globalny)
  warstwa: domyślna
```

Przykładowy przepływ konfiguracji profilu izolacji z użyciem komponentów systemu (odpowiada scenariuszowi rozdz. E.3 Koncepcji):

1. Selektor zasięgu (`.dn-karta--pozycja`) — wybór poziomu „projekt"; wybrana pozycja podświetlona tłem `--dn-akcent-tlo`.
2. Macierz izolacji — ustawienie historii, pamięci i kontekstu na „odrębne" (`.dn-suwak`), z pozostawieniem zakresu odczytu i zapisu plików współdzielonym.
3. Podgląd polityki efektywnej (`.dn-kod`) — weryfikacja reguł wynikowych po dziedziczeniu.
4. Akcja „Zapisz profil" (`.dn-btn--zloty`) — potwierdzenie zapisu komunikatem `.dn-toast--sukces`.

---

## 13. Zastosowanie w panelu orkiestracji środowiska MultitaskingAI

Rozdział 13.8 Koncepcji określa zawartość bocznej nawigacji środowiska MultitaskingAI jako panel orkiestracji złożony z sześciu sekcji: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu. Panel zajmuje miejsce analogiczne do bocznej nawigacji modułów w pozostałych trzech środowiskach — wizualnie jest więc zbudowany na tym samym wzorcu listy pionowej co nawigacja modułowa, z sekcjami jako pozycjami `.dn-karta--pozycja` lub `.dn-zakladki` w układzie pionowym.

```
PANEL ORKIESTRACJI — boczna nawigacja MultitaskingAI (rozdz. 13.8 Koncepcji)
┌──────────────────────────────────────────────┐
│  ● Zespoły         (gwiazdka / plik)          │  .dn-karta--pozycja
│  ● Role            (uzytkownik)               │  cztery karty ról + .dn-awatar
│  ● Kolejki         (filtr / strzalka-prawo)   │  .dn-tabela + .dn-btn-ikona
│  ● Orkiestracja    (link-zewnetrzny)          │  .dn-karta--pozycja + .dn-plakietka
│  ● Harmonogram…    (zegar / kalendarz)        │  wzorzec Scheduler (10.6 / 10.7)
│  ● Monitor procesu (oko)                      │  .dn-tabela + .dn-plakietka--stan
└──────────────────────────────────────────────┘
   kolejność i widoczność sekcji — konfigurowalne (zestaw domyślny)
```

| Sekcja | Ikona (Załącznik C) | Wzorzec wizualny |
|---|---|---|
| Zespoły | `gwiazdka` lub `plik` (konfiguracja zapisana) | Lista `.dn-karta--pozycja` z akcjami zapisu/wczytania/duplikowania |
| Role | `uzytkownik` | Cztery karty ról (`.dn-karta`) z `.dn-awatar` per rola i przypisanym agentem |
| Kolejki | `filtr` lub `strzalka-prawo` (przepływ) | `.dn-tabela` z akcjami enqueue/dequeue/retry jako `.dn-btn-ikona` |
| Orkiestracja | `link-zewnetrzny` (zależność) | Lista zależności `.dn-karta--pozycja`, status zgodności `.dn-plakietka` |
| Harmonogram i automatyki | `zegar`/`kalendarz` | Lista wpiętych automatyk, wzorzec Scheduler (rozdz. 10.6/10.7) |
| Monitor procesu | `oko` | Wzorzec monitora (rozdz. 10.5): `.dn-tabela` + `.dn-plakietka--stan` |

Sekcja Role prezentuje cztery okna robocze (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) jako karty z `.dn-awatar` — wariant `--kwadrat` dla agenta przypisanego z modułu Agents, wariant kołowy dla modelu bazowego bez przypisanego agenta. Mechanizm Subagent Network, uruchamiany przez wykonawcę, jest prezentowany jako rozwinięcie karty roli (do 15 pozycji `.dn-karta--pozycja` zagnieżdżonych), a nie jako odrębna sekcja panelu, zgodnie z jego statusem mechanizmu, a nie roli zespołu (Koncepcja, Załącznik B).

Zgodnie z zasadą pełnej konfigurowalności (Koncepcja, rozdz. 6, rozdz. 14) kolejność i widoczność sekcji panelu orkiestracji są konfigurowalne z okna konfiguracji — wizualnie odwzorowane możliwością przeciągnięcia lub ukrycia pozycji listy sekcji, analogicznie do mechanizmu profili izolacji (rozdz. 12.3); przedstawiony w tabeli powyżej zestaw i dobór ikon jest zestawem domyślnym.

---

## 14. Zgodność z zasadami nadrzędnymi platformy

System wizualny nie wprowadza żadnego ograniczenia sprzecznego z czterema zasadami nadrzędnymi funkcjonalności platformy (Koncepcja, rozdz. 14):

| Zasada nadrzędna (Koncepcja, rozdz. 14) | Realizacja w systemie wizualnym |
|---|---|
| Pełna kompozycyjność | Jeden zestaw tokenów i komponentów obsługuje wszystkie środowiska, moduły i okna operacyjne bez wariantów wykluczających się wzajemnie; komponent użyty w jednym oknie (np. `.dn-tabela`) jest dostępny i spójny w każdym innym oknie, które go potrzebuje. |
| Pełna konfigurowalność | Motyw (jasny/ciemny) jest ustawieniem z okna konfiguracji, nie regułą sztywną; kolejność i widoczność sekcji list (panel orkiestracji, panele boczne) są konfigurowalne, a system wizualny nie koduje wyjątków od tej zasady. |
| Jawność i konfigurowalność zależności | Objaśnienia kontekstowe `[?]` (rozdz. 8.11) czynią każdy token i przełącznik konfiguracji jawnie opisanym w miejscu użycia, zgodnie z zasadą Architektury (rozdz. 13). |
| Rozszerzenie orkiestracji | Biblioteka komponentów (tabele, plakietki stanu, karty ról, awatary) skaluje się wraz ze wzrostem liczby ról, agentów i procesów równoległych w środowisku MultitaskingAI, bez potrzeby definiowania nowych wzorców wizualnych przy każdym rozszerzeniu skali. |

Ponadto system wizualny jest zgodny z dyrektywą nadrzędną „Zero blokerów" obowiązującą w projekcie: żaden token, komponent ani wzorzec zastosowania opisany w niniejszym dokumencie nie tworzy blokady niedostępnej do zmiany z poziomu konfiguracji tam, gdzie dana warstwa jest konfigurowalna (motyw, kolejność sekcji, stan przełączników izolacji) — brak jawnego ustawienia przyjmuje wartość domyślną, nigdy stan zablokowany.

---

## Załącznik A. Pełna tabela tokenów kolorów semantycznych

| Token | Jasny | Ciemny |
|---|---|---|
| `tlo` | `#F7F8FA` | `#0B1120` |
| `powierzchnia` | `#FFFFFF` | `#16223C` |
| `powierzchnia-2` | `#EEF1F6` | `#10192E` |
| `panel` | `#FFFFFF` | `#16223C` |
| `hover` | `#EEF1F6` | `rgba(255,255,255,0.06)` |
| `nakladka` | `rgba(11,17,32,0.42)` | `rgba(3,6,14,0.62)` |
| `obrys` | `#E2E6EE` | `#26395C` |
| `obrys-subtelny` | `#EEF1F6` | `#1E2E4E` |
| `obrys-mocny` | `#CBD2DF` | `#3A4E78` |
| `tekst` | `#12161F` | `#E8EDF7` |
| `tekst-2` | `#4C566B` | `#9AA6BF` |
| `tekst-3` | `#6E7A93` | `#6B7E9C` |
| `tekst-inv` | `#FFFFFF` | `#0B1120` |
| `marka` | `#1A2A4A` | `#26395C` |
| `marka-hover` | `#1E2E4E` | `#3A4E78` |
| `akcent` | `#C9A24B` | `#D9B968` |
| `akcent-txt` | `#8A6D2E` | `#E6C67A` |
| `akcent-tlo` | `#FBF7EC` | `rgba(201,162,75,0.12)` |
| `akcent-obrys` | `#E6C67A` | `#A6812F` |
| `focus` | `#C9A24B` | `#D9B968` |
| `focus-cien` | `rgba(201,162,75,0.35)` | `rgba(217,185,104,0.40)` |
| `sukces-fg / bg / br` | `#1F7A50 / #E7F4EE / #B9E0CE` | `#4CAF7D / rgba(76,175,125,0.14) / rgba(76,175,125,0.34)` |
| `ostrzezenie-fg / bg / br` | `#8A6212 / #FBF1DA / #ECD8A6` | `#E0A63C / rgba(224,166,60,0.14) / rgba(224,166,60,0.34)` |
| `blad-fg / bg / br` | `#B23A34 / #FBECEB / #F0C9C6` | `#E06A66 / rgba(224,106,102,0.14) / rgba(224,106,102,0.34)` |
| `info-fg / bg / br` | `#2E5C94 / #E9F1FA / #C6DBF0` | `#5B9BD8 / rgba(91,155,216,0.14) / rgba(91,155,216,0.34)` |

Źródło wartości: `tokens.json`, sekcja `semantyczne`.

---

## Załącznik B. Tabela klas komponentów

| Komponent | Klasa bazowa | Warianty |
|---|---|---|
| Przycisk | `.dn-btn` | `--glowny`, `--zloty`, `--zarys`, `--duch`, `--blad`, `--sm`, `--lg` |
| Przycisk ikonowy | `.dn-btn-ikona` | — |
| Pole formularza | `.dn-pole`, `.dn-input`, `.dn-select`, `.dn-textarea` | `--blad` |
| Przełącznik | `.dn-suwak` | — |
| Checkbox / radio | `.dn-check` | — |
| Karta | `.dn-karta` | `--interaktywna`, `--akcent`, `--pozycja` |
| Plakietka | `.dn-plakietka` | `--zloto`, `--sukces`, `--ostrz`, `--blad`, `--info`, `--wersaliki`, `--stan` |
| Kropka stanu | `.dn-kropka` | `--sukces`, `--ostrz`, `--blad` |
| Zakładki | `.dn-zakladki` | `--pigulki` |
| Tabela | `.dn-tabela`, `.dn-tabela-owijka` | `--paski` |
| Kolumna nawigacji | `.dn-pasek` | — |
| Awatar | `.dn-awatar` | `--zloto`, `--sm`, `--lg`, `--kwadrat`, `-stan` |
| Modal | `.dn-nakladka`, `.dn-modal` | — |
| Toast | `.dn-toast` | `--sukces`, `--blad`, `--ostrz` |
| Tooltip | `.dn-tooltip` | — |
| Spinner | `.dn-spinner` | — |
| Separator | `.dn-separator` | — |
| Blok kodu | `.dn-kod` | — |
| Alert | `.dn-alert` | `--info`, `--sukces`, `--ostrz`, `--blad` |
| Pusty stan | `.dn-pusty-stan` | — |
| Ikona | `.dn-ikona` | `--sm`, `--lg` |

Źródło: `components.css`.

---

## Załącznik C. Katalog ikon i typowe zastosowanie

| Ikona | Typowe zastosowanie w platformie |
|---|---|
| `dom` | Strona główna / Centrum dowodzenia |
| `ustawienia` | Okno konfiguracji, kolumna ustawień (strefa 3) |
| `menu` | Wyzwalacz `☰` — rozwinięcie kolumny nawigacji / panelu orkiestracji |
| `wiecej` | Menu kebab `⋮` — menu kontekstowe wiersza / karty |
| `szukaj` | Pole wyszukiwania (kolumna nawigacji, panele list, wyszukiwarka funkcji) |
| `filtr` | Filtrowanie list (Sources Manager, Findings Panel, Kolejki) |
| `plus` | Dodanie pozycji (nowa karta sesji, nowy komponent własny) |
| `zamknij` | Zamknięcie karty sesji, modala, toastu |
| `kosz` | Usunięcie pozycji |
| `archiwum` | Session Repository, historia wersji |
| `pobierz` | Eksport (Export Panel, Deployment Panel) |
| `wyslij` | Wysłanie wiadomości w Chat Window |
| `spinacz` | Załącznik w Chat Window |
| `odpowiedz` | Odpowiedź / cytowanie w rozmowie |
| `dokument` | Studio, Library, dokumenty; zakres odczytu i zapisu plików (izolacja techniczna) |
| `plik` | Project Library, Library Explorer; katalog danych i konfiguracji modelu (izolacja techniczna) |
| `folder` | Katalog roboczy sesji (izolacja techniczna), Project Tree |
| `obraz` | Design, Assets Panel |
| `kod` | Developer, Code Editor, model procesu (izolacja techniczna) |
| `karta-przegladarki` | Browser Window |
| `link-zewnetrzny` | Źródło zewnętrzne, zależność w Orkiestracji |
| `globus` | Dostęp sieciowy (izolacja techniczna) |
| `koperta` | Assistant, powiadomienia e-mail |
| `dzwonek` | Powiadomienia, Always On Display |
| `uzytkownik` | Awatar / profil / role MultitaskingAI |
| `klodka` | Konto i token per sesja (izolacja techniczna) |
| `tarcza` | Security Auditor, Permissions Center |
| `waga` | Arbitrator (rola Executor 3 / Validator) |
| `oko` | Podgląd, Monitor procesu |
| `gwiazdka` | Oznaczenie ważności, Zespoły (presety) |
| `zegar` | Harmonogram, Scheduler |
| `kalendarz` | Harmonogram i automatyki |
| `odswiez` | Ponowienie (retry) w kolejce lub procesie |
| `uruchom` | Start procesu, serwer wykonania (izolacja techniczna) |
| `zatrzymaj` | Zatrzymanie procesu |
| `ptaszek` / `ptaszek-kolo` | Potwierdzenie, walidacja (Validator) |
| `info` | Alert / plakietka informacyjna |
| `ostrzezenie` | Alert / plakietka ostrzeżenia |
| `blad` | Alert / plakietka błędu |
| `strzalka-lewo` / `strzalka-prawo` | Nawigacja kart, Diff Panel, przepływ Orkiestracji |
| `strzalka-dol` | Znak `▼` wyzwalacza menu progresywnego i elementu zbiorczego |
| `slonce` / `ksiezyc` | Przełącznik motywu jasny/ciemny |

Pełna galeria z podglądem SVG: `ikony/indeks.html`.

---

*Koniec dokumentu. Danaco Console — System wizualny, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
