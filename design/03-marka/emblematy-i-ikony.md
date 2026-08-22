# Emblematy środowisk, ikona aplikacji, favicon
## System godeł pochodnych Danaco Console

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment |
| **Warstwa** | marka · system godeł pochodnych |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Odbiorcy** | projektant wizualny · zespół wdrożeniowy front-end · osoba pakująca aplikację (Tauri, PWA) · dostawca materiałów firmowych |
| **Zakres** | hierarchia godeł, konstrukcja czterech emblematów środowisk, drabina rozmiarowa 16–64 px, rastry, ikona aplikacji, wariant maskowalny, favicon, ICO, manifest, snippet `<head>`, zasady doboru |
| **Poza zakresem** | konstrukcja sygnetu i logotypu (→ `03-marka/ksiega-znaku.md`), zestaw ikon interfejsu (→ `zasoby/ikony/manifest.json`), typografia i barwa systemu (→ kontrakt systemu projektowego) |
| **Źródła wiążące** | kontrakt systemu projektowego · kierunek systemu projektowego · realne pliki `zasoby/marka/**` · `zasoby/css/komponenty.css` · `zasoby/zetony/zetony.css` |
| **Wytworzone pliki** | `03-marka/emblematy/` — 102 pliki |

> **Zasada nadrzędna tego opracowania.** Nie powstał tu ani jeden nowy znak.
> Wszystko poniżej jest **pochodną zatwierdzonej geometrii**: sygnetu „Delegacja"
> (siatka 96 × 96) i czterech emblematów środowisk (siatka 24 × 24, obrys 1,75).
> Wolno było zbudować warianty rozmiarowe, konfiguracje i zastosowania.
> Nie wolno było ruszyć krzywych — i nie ruszono.

---

## Spis treści

1. [Pozycja opracowania w systemie marki](#1-pozycja-opracowania-w-systemie-marki)
2. [System godeł — hierarchia i wspólna gramatyka](#2-system-godeł--hierarchia-i-wspólna-gramatyka)
3. [Cztery emblematy środowisk — konstrukcja](#3-cztery-emblematy-środowisk--konstrukcja)
4. [Prawo kropki — analiza czterech położeń](#4-prawo-kropki--analiza-czterech-położeń)
5. [Barwa emblematu — rozstrzygnięcie](#5-barwa-emblematu--rozstrzygnięcie)
6. [Drabina rozmiarowa 16 → 64 px](#6-drabina-rozmiarowa-16--64-px)
7. [Rastry PNG emblematów](#7-rastry-png-emblematów)
8. [Ikona aplikacji](#8-ikona-aplikacji)
9. [Favicon](#9-favicon)
10. [Snippet `<head>` i manifest](#10-snippet-head-i-manifest)
11. [Zasady doboru — sygnet, emblemat, ikona aplikacji](#11-zasady-doboru--sygnet-emblemat-ikona-aplikacji)
12. [Podgląd w kontekście](#12-podgląd-w-kontekście)
13. [Inwentarz wytworzonych plików](#13-inwentarz-wytworzonych-plików)
14. [Kontrola jakości](#14-kontrola-jakości)
15. [Decyzje projektowe](#15-decyzje-projektowe)

---

## 1. Pozycja opracowania w systemie marki

### 1.1 Trzy dokumenty, trzy zakresy

| Dokument | Odpowiada na pytanie | Zakres |
|---|---|---|
| `03-marka/ksiega-znaku.md` | „Jak wygląda znak marki i czego z nim nie wolno?" | sygnet, logotyp, lockupy, pole ochronne, katalog naruszeń |
| **`03-marka/emblematy-i-ikony.md`** (ten plik) | „Jak znak rozgałęzia się w rodzinę godeł pochodnych?" | emblematy środowisk, ikona aplikacji, favicon, drabina rozmiarowa, pliki wynikowe |
| `zasoby/ikony/manifest.json` | „Jakiej ikony użyć do danego pojęcia?" | 82 ikony interfejsu z biblioteki Lucide |

Rozgraniczenie jest ostre: **księga znaku mówi o marce, to opracowanie o rodzinie
godeł, manifest ikon o słowniku interfejsu.** Trzy poziomy nie mieszają się —
emblemat środowiska nigdy nie stanie na wizytówce, sygnet nigdy nie stanie
w bocznej nawigacji jako pozycja modułu.

### 1.2 Co to opracowanie wnosi ponad stan zastany

Stan zastany (`zasoby/marka/**`) zawierał: cztery emblematy w jednym rozmiarze,
ikonę aplikacji w pięciu rastrach, favicon w trzech rastrach i jednym ICO.
Opracowanie dokłada:

- **drabinę rozmiarową** emblematów (16 / 20 / 24 / 32 / 48 / 64 px) z regułą
  upraszczania poniżej 24 px — dotąd nieistniejącą,
- **komplet rastrów** emblematów (24 / 48 / 96 / 192 px × jasny / ciemny),
- **komplet ikony aplikacji** obejmujący rozmiary Tauri i PWA oraz ICO
  wielorozmiarowy z osobną grafiką na każdy rozmiar,
- **naprawę trzech wad plików zastanych** — udokumentowanych pomiarowo
  w rozdz. 9.6,
- **generator** (`emblematy/generator-godel.py`) — cały komplet odtwarza się
  jednym poleceniem z tej samej geometrii źródłowej.

### 1.3 Uwaga o dokumencie poprzednie wydanie przewodnika marki

Plik poprzednie wydanie przewodnika marki (v1.0,
para **granat + złoto**, tarcza kancelaryjna z monogramem „D") opisuje
**poprzednią** warstwę wizualną. Zastąpił go kierunek systemu projektowego i kontrakt systemu projektowego:
monochrom + jeden błękit sygnałowy, sygnet „Delegacja". W tym opracowaniu
przewodnik v1.0 pełni wyłącznie rolę historyczną — **żadna jego wartość barwna,
typograficzna ani konstrukcyjna nie obowiązuje.**

---

## 2. System godeł — hierarchia i wspólna gramatyka

### 2.1 Trzy piętra godła

```
                       ┌──────────────────────────────────────┐
   PIĘTRO 1            │        SYGNET „DELEGACJA"            │   siatka 96 × 96
   ZNAK NADRZĘDNY      │        »». — dwa groty + kropka      │   pełne krzywe
                       │        barwy własne marki            │   ≥ 24 px
                       └───────────────────┬──────────────────┘
                                           │  „czyj to produkt"
                    ┌──────────────────────┼──────────────────────┐
                    ▼                      ▼                      ▼
   PIĘTRO 2   ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐
   EMBLEMATY  │  TalkIn   │  │ WorkSpace │  │CodeStudio │  │Multitask… │  siatka 24 × 24
   ŚRODOWISK  │   dymek   │  │  4 moduły │  │  ramka >_ │  │  4 węzły  │  obrys 1,75
              └───────────┘  └───────────┘  └───────────┘  └───────────┘  currentColor
                    │              │              │              │
                    └──────────────┴──────┬───────┴──────────────┘
                                          │  „gdzie jestem"
                                          ▼
   PIĘTRO 3   ┌──────────────────────────────────────────────────────────┐
   IKONY      │  82 ikony Lucide: studio · research · terminal · agent…  │  siatka 24 × 24
   MODUŁÓW    │  bez kropki sygnału — kropka należy do godeł             │  obrys 1,75
              └──────────────────────────────────────────────────────────┘  currentColor
                                          │  „czym pracuję"
                                          ▼
              ┌──────────────────────────────────────────────────────────┐
   ODNOGA     │  IKONA APLIKACJI · FAVICON — sygnet w kaflu kryjącym     │  1024 / 96
   SYSTEMOWA  │  grunt --dn-rama · promień naroża 21,875 %               │  barwy własne
              └──────────────────────────────────────────────────────────┘
```

Hierarchia jest **jednokierunkowa**: piętro niższe cytuje wyższe, nigdy odwrotnie.
Emblemat cytuje sygnet (kropka), ikona modułu nie cytuje niczego (jest słownikiem),
ikona aplikacji jest oprawą sygnetu, nie osobnym znakiem.

### 2.2 Co łączy — pięć elementów wspólnej gramatyki

| # | Element wspólny | Wartość | Zakres obowiązywania |
|---|---|---|---|
| 1 | **Siatka** | 24 × 24 j. | emblematy + ikony modułów (sygnet: 96 × 96 = czterokrotność) |
| 2 | **Obrys** | 1,75 j., `linecap: round`, `linejoin: round` | emblematy + ikony modułów |
| 3 | **Wypełnienie** | `none` na obrysach | emblematy + ikony modułów |
| 4 | **Barwa** | `currentColor` — jedna barwa na cały znak | emblematy + ikony modułów |
| 5 | **Kropka sygnału** | dokładnie **jedna** wypełniona kropka | **wyłącznie** emblematy i sygnet — ikony modułów jej nie mają |

Piąty punkt jest tym, co czyni emblemat emblematem. Ikona `terminal` i emblemat
`srodowisko-codestudio` leżą na tej samej siatce, mają ten sam obrys i tę samą
barwę — różni je wyłącznie obecność wypełnionej kropki. **Kropka to znak
przynależności do rodziny godeł.**

### 2.3 Co rozróżnia — metafora i ciężar

| Emblemat | Metafora | Ciężar optyczny | Dominanta kierunkowa |
|---|---|---|---|
| TalkIn | dymek rozmowy z liniami tekstu | średni — jedna duża forma zamknięta | pozioma (linie tekstu) |
| WorkSpace | cztery moduły projektu | wysoki — cztery bryły równej masy | siatka 2 × 2, bez dominanty |
| CodeStudio | ramka terminala z zachętą `>_` | wysoki — pełna ramka | pozioma (ramka 18 × 16) |
| MultitaskingAI | rdzeń koordynatora + trzej wykonawcy | niski — punkty i cienkie łączniki | promienista, oś pionowa |

Zestaw jest zrównoważony celowo: dwie formy zamknięte (TalkIn, CodeStudio),
jedna siatkowa (WorkSpace), jedna otwarta promienista (MultitaskingAI). Cztery
emblematy stojące obok siebie na stronie głównej nie zlewają się w plamę.

### 2.4 Czego w gramatyce nie ma i być nie może

- **gradientu** — gradienty są wyłącznie ilustracyjne (kontrakt systemu projektowego), nigdy w godle,
- **drugiej barwy** — jeden emblemat = jedna barwa, także kropka,
- **wypełnienia obszarowego** — poza kropką wszystko jest obrysem,
- **cienia, poświaty, obrysu wtórnego**,
- **emoji jako godła** — bezwzględny zakaz z kontraktu systemu projektowego.

---

## 3. Cztery emblematy środowisk — konstrukcja

Wszystkie pomiary poniżej odczytano z realnych plików
`zasoby/marka/srodowiska/srodowisko-*.svg`. Jednostka = 1/24 boku siatki.

### 3.1 Tabela zbiorcza

| Emblemat | Plik źródłowy | Elementów obrysowych | Kropka (cx, cy, r) | Pole zadruku (x) | Pole zadruku (y) |
|---|---|---|---|---|---|
| TalkIn | `srodowisko-talkin.svg` | 3 | 16,2 · 10,8 · **1,5** | 2,625 – 21,375 | 2,125 – 20,875 |
| WorkSpace | `srodowisko-workspace.svg` | 4 | 16,9 · 16,9 · **1,6** | 2,625 – 21,375 | 2,625 – 21,375 |
| CodeStudio | `srodowisko-codestudio.svg` | 3 | 18,4 · 9,9 · **1,5** | 2,125 – 21,875 | 3,125 – 20,875 |
| MultitaskingAI | `srodowisko-multitaskingai.svg` | 6 | 12,0 · 12,0 · **2,6** | 2,125 – 21,975 | 1,825 – 19,175 |

Pole zadruku liczone z uwzględnieniem połowy obrysu (0,875 j.) po każdej stronie.

### 3.2 TalkIn — dymek rozmowy

**Metafora.** Środowisko rozmowy: dymek z dwiema liniami tekstu, przy końcu
drugiej linii biegnie odpowiedź. Ogon dymka skierowany w dół-lewo — mowa wychodzi
od Operatora.

**Konstrukcja na siatce 24:**

```
   0    3.5              12               20.5   24
 0 ┌───────────────────────────────────────────┐
   │                                           │
 3 │   ╭───────────────────────────────────╮   │  łuk narożny r = 2,5
   │   │                                   │   │
7.5│   │   ────────────────────────        │   │  linia 1: M7.5 7.5 h9
   │   │   x 7,5 → 16,5                    │   │
10.8   │   ──────────────      ●           │   │  linia 2: M7.5 10.8 h5
   │   │   x 7,5 → 12,5     kropka         │   │  KROPKA 16,2 · 10,8 · r 1,5
   │   │                    16,2 · 10,8    │   │
16 │   ╰──┬──────────────────────────────╯    │
   │      │╲                                  │  ogon: v4 → l4.4 -4
20 │      └─╲                                 │  wierzchołek (7 · 20)
   │                                          │
24 └──────────────────────────────────────────┘
```

| Element | Ścieżka / atrybuty | Rola |
|---|---|---|
| Dymek | `M20.5 5.5a2.5 2.5 0 0 0-2.5-2.5H6A2.5 2.5 0 0 0 3.5 5.5v8A2.5 2.5 0 0 0 6 16h1v4l4.4-4H18a2.5 2.5 0 0 0 2.5-2.5Z` | forma zamknięta, promień naroża 2,5 |
| Linia tekstu 1 | `M7.5 7.5h9` (dł. 9 j.) | wypowiedź zakończona |
| Linia tekstu 2 | `M7.5 10.8h5` (dł. 5 j.) | wypowiedź w toku — dlatego krótsza |
| **Kropka sygnału** | `circle 16,2 · 10,8 · r 1,5`, `fill: currentColor` | „tu biegnie odpowiedź" — leży **na osi drugiej linii**, 3,7 j. za jej końcem |

**Dlaczego kropka jest tam, gdzie jest.** Kropka i druga linia mają identyczne
`cy` = 10,8. Odstęp między końcem linii (12,5) a lewą krawędzią kropki (14,7)
wynosi 2,2 j. — czyta się jako pauza przed dalszym ciągiem, nie jako przypadek.

**Zastosowanie:** karta środowiska TalkIn na stronie głównej (40 px), pozycja
środowiska w bocznej nawigacji (20 px), tytuł powłoki w pasku górnym (20 px),
karta sesji (16 px), lista wyboru środowiska w Centrum poleceń (16 px).

### 3.3 WorkSpace — cztery moduły projektu

**Metafora.** Środowisko pracy projektowej: cztery moduły równej masy; w jednym
z nich — prawym dolnym — biegnie praca.

**Konstrukcja na siatce 24:**

```
    0   3.5        10.7  13.3       20.5  24
  0 ┌─────────────────────────────────────────┐
    │                                         │
3.5 │   ┌──────────┐    ┌──────────┐          │  cztery prostokąty
    │   │          │    │          │          │  7,2 × 7,2 · rx 2
    │   │          │    │          │          │
10.7│   └──────────┘    └──────────┘          │  szczelina 2,6 j.
    │        ▲              ▲                 │
13.3│   ┌──────────┐    ┌──────────┐          │
    │   │          │    │    ●     │          │  KROPKA 16,9 · 16,9 · r 1,6
    │   │          │    │  16,9    │          │  = geometryczny środek modułu
20.5│   └──────────┘    └──────────┘          │
    │                                         │
 24 └─────────────────────────────────────────┘
```

| Element | Współrzędne | Rola |
|---|---|---|
| Moduł ↖ | `rect 3,5 · 3,5 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↗ | `rect 13,3 · 3,5 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↙ | `rect 3,5 · 13,3 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↘ | `rect 13,3 · 13,3 · 7,2 × 7,2 · rx 2` | moduł **aktywny** |
| **Kropka sygnału** | `circle 16,9 · 16,9 · r 1,6` | środek modułu ↘: 13,3 + 7,2/2 = 16,9 — **dokładnie**, w obu osiach |

**Dlaczego kropka jest tam, gdzie jest.** Jest jedynym elementem łamiącym
symetrię czterokrotną. Emblemat bez kropki byłby obojętnym znakiem „siatka";
kropka nadaje mu kierunek czytania (ku prawej dolnej ćwiartce — tam, gdzie
w interfejsie kończy się przepływ pracy).

**Zastosowanie:** karta środowiska WorkSpace, boczna nawigacja WorkSpace,
kafel „Workspace (Projekt)" strefy 2 Centrum dowodzenia, nagłówek okna
`Project Dashboard`.

### 3.4 CodeStudio — ramka terminala

**Metafora.** Środowisko wytwarzania: ramka okna terminala, zachęta `>` i linia
wyniku pod nią; w prawym górnym rogu — kropka procesu.

**Konstrukcja na siatce 24:**

```
    0   3            10.1      16.6    21  24
  0 ┌────────────────────────────────────────┐
    │                                        │
  4 │  ┌───────────────────────────────────┐ │  rect 18 × 16 · rx 2,5
    │  │                          ●        │ │  KROPKA 18,4 · 9,9 · r 1,5
9.2 │  │   ╲                    18,4·9,9   │ │
    │  │    ╲                              │ │  grot: M7 9.2 l3.1 2.8 L7 14.8
 12 │  │     ▶ 10,1 · 12                   │ │  wierzchołek na osi y = 12
    │  │    ╱                              │ │
14.8│  │   ╱                               │ │
    │  │          ──────────               │ │  linia wyniku: M12.6 15.4 h4
15.4│  │          12,6 → 16,6              │ │
 20 │  └───────────────────────────────────┘ │
 24 └────────────────────────────────────────┘
```

| Element | Współrzędne | Rola |
|---|---|---|
| Ramka | `rect 3 · 4 · 18 × 16 · rx 2,5` | okno terminala |
| Grot zachęty | `M7 9.2l3.1 2.8L7 14.8` | znak `>` — **jedyny bezpośredni cytat z sygnetu** w całym zestawie emblematów |
| Linia wyniku | `M12.6 15.4h4` (dł. 4 j.) | odpowiedź powłoki |
| **Kropka sygnału** | `circle 18,4 · 9,9 · r 1,5` | proces w tle; prawa krawędź kropki 19,9 j. — 0,225 j. od wewnętrznej krawędzi ramki |

**Dlaczego kropka jest tam, gdzie jest.** Prawy górny róg okna terminala to
w interfejsie miejsce wskaźnika stanu procesu. Kropka jest tam **maksymalnie
dosunięta** — 0,225 j. luzu do ramki. To najciaśniejsze miejsce w całym systemie
godeł i przy skalowaniu poniżej 24 px wymaga korekty (rozdz. 6.2).

**Zastosowanie:** karta środowiska CodeStudio, boczna nawigacja CodeStudio,
nagłówki okien `Terminal Tabs`, `Code Editor`, `Build Output`.

### 3.5 MultitaskingAI — rdzeń i trzej wykonawcy

**Metafora.** Środowisko orkiestracji: wypełniony rdzeń w środku (Coordinator)
i trzy węzły obrysowe (Executor 1, Executor 2, Executor 3 / Validator) połączone
z rdzeniem. Emblemat powtarza opowieść sygnetu: zlecenie idzie od koordynatora
do wykonawców.

**Konstrukcja na siatce 24:**

```
    0        4.9        12        19.1     24
  0 ┌───────────────────────────────────────┐
    │                                       │
4.6 │                  ◯ 12 · 4,6 · r 1,9   │  węzeł górny (Executor 1)
    │                  │                    │
6.5 │                  │  łącznik M12 9.4 V6.5
    │                  │                    │
 12 │        ╱────────●────────╲            │  RDZEŃ 12 · 12 · r 2,6
    │       ╱      12 · 12      ╲           │  (jedyna kropka wypełniona)
13.4│      ╱  M9.8 13.4       M14.2 13.4 ╲  │
    │     ╱   l-3.3 1.9        l3.3 1.9   ╲ │
16.4│    ◯                              ◯   │  węzły dolne (Executor 2,
    │  4,9·16,4                    19,1·16,4│   Executor 3 / Validator)
 24 └───────────────────────────────────────┘
```

| Element | Współrzędne | Rola |
|---|---|---|
| **Rdzeń — kropka sygnału** | `circle 12 · 12 · r 2,6`, `fill: currentColor` | Coordinator; **środek geometryczny siatki co do jednostki** |
| Węzeł górny | `circle 12 · 4,6 · r 1,9` | Executor 1 |
| Węzeł lewy dolny | `circle 4,9 · 16,4 · r 1,9` | Executor 2 |
| Węzeł prawy dolny | `circle 19,1 · 16,4 · r 1,9` | Executor 3 / Validator |
| Łącznik górny | `M12 9.4V6.5` (dł. 2,9 j.) | delegacja w górę; **wchodzi 0,875 j. w obręb węzła** |
| Łącznik lewy | `M9.8 13.4l-3.3 1.9` (dł. 3,81 j.) | delegacja w dół-lewo; kończy się **na krawędzi** węzła |
| Łącznik prawy | `M14.2 13.4l3.3 1.9` (dł. 3,81 j.) | delegacja w dół-prawo; symetryczny do lewego |

**Dlaczego kropka jest tam, gdzie jest.** To jedyny emblemat, w którym kropka
**nie jest satelitą, lecz środkiem**. Uzasadnienie merytoryczne: w MultitaskingAI
praca nie „dzieje się gdzieś w rogu okna" — praca **jest** koordynacją. Kropka
ma też największy promień w zestawie (2,6 wobec 1,5–1,6), bo pełni podwójną
funkcję: nosi sygnał i jest węzłem konstrukcji.

**Trzy pomiary warte odnotowania:**

- węzły dolne są od osi pionowej odsunięte symetrycznie: 12 − 4,9 = 7,1;
  19,1 − 12 = 7,1 — **symetria dokładna**,
- pole zadruku w pionie: 1,825 – 19,175 j., środek 10,5 j. — **1,5 j. powyżej
  środka siatki**; to kompensacja optyczna trójkąta o podstawie na dole,
- węzeł górny nie leży w narożniku pola — emblemat jest jedynym w zestawie,
  który **nie wypełnia siatki po brzegi** (pole zadruku 19,85 × 17,35 j.).

**Zastosowanie:** karta środowiska MultitaskingAI, **panel orkiestracji**
(zamiast bocznej nawigacji — MultitaskingAI jej nie ma), nagłówki okien
`Coordinator Chat`, `Results Analyzer`, plakietka roli w oknie komunikacji.

---

## 4. Prawo kropki — analiza czterech położeń

### 4.1 Pomiar

| Emblemat | cx | cy | r | Ø | Ćwiartka | Odległość od środka siatki |
|---|---:|---:|---:|---:|---|---:|
| TalkIn | 16,2 | 10,8 | 1,5 | 3,0 | prawa górna (nieznacznie) | 4,36 j. |
| WorkSpace | 16,9 | 16,9 | 1,6 | 3,2 | prawa dolna | 6,93 j. |
| CodeStudio | 18,4 | 9,9 | 1,5 | 3,0 | prawa górna | 6,73 j. |
| MultitaskingAI | 12,0 | 12,0 | 2,6 | 5,2 | **środek** | 0,00 j. |

### 4.2 Reguła i jej jedyny wyjątek

```
   cx wszystkich kropek satelickich:  16,2 · 16,9 · 18,4
   ─────────────────────────────────────────────────────
   pas pionowy x ∈ [16,2 ; 18,4]  ─ szerokość 2,2 j. na siatce 24
                                    = 9,2 % szerokości znaku

   MultitaskingAI: cx = 12,0  ← WYJĄTEK ŚWIADOMY
```

Trzy z czterech kropek leżą w wąskim pasie po prawej stronie znaku. Czwarta stoi
na osi. To nie jest niekonsekwencja — to **reguła z jednym umotywowanym
wyjątkiem**, a wyjątek niesie różnicę znaczeniową: trzy środowiska mają pracę
*obok* siebie, MultitaskingAI ma pracę *w sobie*.

### 4.3 Średnica kropki a żeton systemowy

Kropka satelicka ma Ø ≈ 3 j. na siatce 24. Po wyrenderowaniu:

| Render | Ø kropki | Zgodność z żetonem |
|---:|---:|---|
| 16 px | 2,00 px | poniżej żetonu — dlatego wariant 16 px powiększa kropkę (rozdz. 6.2) |
| 24 px | 3,00 px | — |
| 32 px | 4,00 px | — |
| **48 px** | **6,00 px** | **= `--dn-wym-kropka` (6 px) co do piksela** |
| 64 px | 8,00 px | — |

To zbieżność strukturalna, nie przypadek: kropka emblematu wyrenderowana
w 48 px jest **dokładnie tą samą kropką**, którą interfejs stawia przy karcie
sesji i przy nadawcy piszącym. Rodzina domyka się na jednej liczbie.

### 4.4 Diagram nakładkowy — wspólna gramatyka

```
   siatka 24 × 24 · nakładka czterech emblematów

   ┌───────────────────────────────────────────────┐
   │                                               │
   │              ┌ MultitaskingAI ┐               │
   │              │   węzeł górny  │               │
   │              └───────┬────────┘               │
   │  ╔════════════════════════════════════════╗   │  ← pole zadruku wspólne
   │  ║                              ░░░       ║   │     x 2,125 – 21,975
   │  ║                              ░C░ 18,4  ║   │     y 1,825 – 21,375
   │  ║              ●             ░░T░ 16,2   ║   │
   │  ║           rdzeń M            ░░░       ║   │  ← PAS KROPEK
   │  ║           12 · 12            ░W░ 16,9  ║   │     x ∈ [16,2 ; 18,4]
   │  ║                              ░░░       ║   │
   │  ╚════════════════════════════════════════╝   │
   │                                               │
   └───────────────────────────────────────────────┘
        T = TalkIn   W = WorkSpace   C = CodeStudio   M = MultitaskingAI
```

---

## 5. Barwa emblematu — rozstrzygnięcie

### 5.1 Napięcie do rozstrzygnięcia

kierunek systemu projektowego mówi: „w **emblematach czterech środowisk**: każdy
zawiera dokładnie jedną wypełnioną kropkę" — w akapicie o **kropce sygnału**,
która jest z definicji błękitna. Realne pliki `srodowisko-*.svg` mają jednak
`fill="currentColor"` także na kropce. `komponenty.css` rozstrzyga to
jednoznacznie: `.dn-karta-srodowiska-godlo { color: var(--dn-tekst); }`,
a pozycja aktywna bocznej nawigacji dostaje sygnał **kreską 2 px obok ikony**
(`::before { background: var(--dn-kropka); }`), nie przebarwieniem emblematu.

### 5.2 Rozstrzygnięcie [NIENEGOCJOWALNE]

> **Emblemat środowiska jest zawsze jednobarwny. Kropka nigdy nie odrywa się
> barwą od obrysu.**

Kropka w emblemacie jest **cytatem formy**, nie cytatem barwy. Cytatem barwy
jest sygnet (kropka `--dn-sygnal-500`) i elementy interfejsu (`.dn-kropka`,
wstęga karty, kreska aktywności).

**Uzasadnienie ilościowe.** Boczna nawigacja WorkSpace ma 9 pozycji, pas kart
sesji do 8 kart, strona główna 4 karty środowisk. Gdyby każdy emblemat niósł
własną błękitną kropkę, na jednym ekranie zapaliłoby się do **21 punktów
sygnału naraz** — barwa przestałaby cokolwiek znaczyć, a to jest wprost
zablokowany odruch z kontraktu systemu projektowego („Dziewięć kolorów tła dla dziewięciu
nadawców").

### 5.3 Cztery dopuszczone barwy emblematu

| Wariant | Żeton | Wartość | Kiedy |
|---|---|---|---|
| na tle jasnym | `--dn-tekst` (motyw jasny) | `#181818` | karta środowiska, nawigacja, nagłówek |
| na tle ciemnym | `--dn-tekst` (motyw ciemny) | `#ECECEC` | jw. w motywie ciemnym |
| środowisko aktywne, motyw jasny | `--dn-sygnal` (motyw jasny) | `#2457C9` | wyłącznie gdy wdrożenie przebarwia **cały** emblemat |
| środowisko aktywne, motyw ciemny | `--dn-sygnal` (motyw ciemny) | `#8FB2F5` | jw. |

W interfejsie stosuje się `currentColor`. Cztery warianty z wypaloną barwą
(**16 plików** w `emblematy/svg/warianty/` — 4 emblematy × 4 barwy) służą
osadzeniom, które `currentColor` nie obsługują: rastrom, dokumentom biurowym,
prezentacjom, poczcie.

### 5.4 Kontrast — pomiar

| Para | Kontrast | Próg (grafika nietekstowa, WCAG 1.4.11) | Wynik |
|---|---:|---:|---|
| `#181818` na `#F4F4F4` (tło jasne) | 16,14 : 1 | 3,0 : 1 | ✓ |
| `#181818` na `#FFFFFF` (powierzchnia) | 17,76 : 1 | 3,0 : 1 | ✓ |
| `#ECECEC` na `#0F0F0F` (tło ciemne) | 16,74 : 1 | 3,0 : 1 | ✓ |
| `#ECECEC` na `#181818` (powierzchnia) | 15,03 : 1 | 3,0 : 1 | ✓ |
| `#2457C9` na `#FFFFFF` | 6,42 : 1 | 3,0 : 1 | ✓ |
| `#8FB2F5` na `#181818` | 8,04 : 1 | 3,0 : 1 | ✓ |

Wartości zgodne z `zasoby/zetony/kontrasty.json`.

---

## 6. Drabina rozmiarowa 16 → 64 px

### 6.1 Tabela kompensacji

| Render | Plik | Siatka | Obrys (j.) | Obrys efektywny | Ø kropki | Detal |
|---:|---|---:|---:|---:|---:|---|
| **16 px** | `…-16.svg` | 24 | **1,90** | 1,27 px | 2,4 – 4,0 px | **uproszczony** |
| **20 px** | `…-20.svg` | 24 | **1,90** | 1,58 px | 2,50 px | pełny |
| **24 px** | `…-24.svg` | 24 | 1,75 | **1,75 px** | 3,00 px | pełny — **kanoniczny** |
| **32 px** | `…-32.svg` | 24 | 1,75 | 2,33 px | 4,00 px | pełny |
| **48 px** | `…-48.svg` | 24 | 1,75 | 3,50 px | 6,00 px | pełny |
| **64 px** | `…-64.svg` | 24 | 1,75 | 4,67 px | 8,00 px | pełny |

**Reguła obrysu.** Powyżej i włącznie z 24 px obrys pozostaje kanoniczny (1,75 j.).
Poniżej 24 px podnosi się do 1,90 j. — kompensacja optyczna, bo 1,75 j. przy
16 px daje 1,17 px, czyli linię cieńszą od piksela, którą rasteryzator rozmywa
na dwa rzędy o połowicznym kryciu. 1,90 j. daje 1,27 px — nadal poniżej piksela,
ale z wyraźnie mocniejszym rdzeniem. **Krzywe pozostają nietknięte.**

**Rozmiar 40 px** (karta środowiska, `.dn-karta-srodowiska-godlo`) obsługuje
plik `…-48.svg` lub dowolny plik ≥ 24 px — SVG jest niezależny od rozdzielczości,
a atrybuty `width`/`height` w plikach są wyłącznie **domyślnym** rozmiarem
renderowania; CSS je nadpisuje.

### 6.2 Reguły upraszczania przy 16 px

Zasada nadrzędna: **usuń najdrobniejszy element, zachowaj kropkę, nie ruszaj
metafory.**

#### TalkIn 16 px

| Zabieg | Z | Na |
|---|---|---|
| usunięto linię tekstu 1 (dłuższą) | `M7.5 7.5h9` | — |
| linię tekstu 2 podniesiono do optycznego środka dymka | `M7.5 10.8h5` | `M7.5 9.4h4.6` |
| kropkę podniesiono wraz z linią i powiększono | `16,2 · 10,8 · r 1,5` | `16,2 · 9,4 · r 1,9` |

Pozostaje para „linia + kropka" — nośnik znaczenia. Odstęp linia–kropka wzrósł
z 2,2 do 2,2 j. przy krótszej linii (koniec 12,1; lewa krawędź kropki 14,3),
co przy 16 px daje 1,47 px przerwy — powyżej progu sklejania.

#### WorkSpace 16 px

| Zabieg | Z | Na |
|---|---|---|
| moduły zwężono (szersza szczelina) | `7,2 × 7,2` @ 3,5 / 13,3 | `6,9 × 6,9` @ 3,55 / 13,55 |
| promień naroża zmniejszono | `rx 2` | `rx 1,4` |
| **usunięto obrys modułu aktywnego** | `rect 13,3 · 13,3` | — |
| kropkę powiększono do rangi modułu | `r 1,6` | `r 3,0` |

**Uzasadnienie usunięcia czwartego obrysu.** Przy 16 px wnętrze modułu ma
3,3 px, a kropka 2,5 px — obrys i kropka sklejają się w nieczytelną plamę
(zweryfikowano rasteryzacją). Emblemat czyta się wtedy jako „trzy kwadraty
i klaks". Po usunięciu obrysu czyta się jako **trzy moduły i jeden aktywny** —
znaczenie nienaruszone, forma czytelna. Liczba wypełnionych kropek pozostaje
dokładnie jedna.

#### CodeStudio 16 px

| Zabieg | Z | Na |
|---|---|---|
| **usunięto linię wyniku** | `M12.6 15.4h4` | — |
| ramkę zwężono o 0,2 j. z każdej strony | `rect 3 · 4 · 18 × 16 · rx 2,5` | `rect 3,2 · 4,2 · 17,6 × 15,6 · rx 2,2` |
| grot przesunięto o 0,2 j. w prawo (środek pionowy ramki) | `M7 9.2l3.1 2.8L7 14.8` | `M7.2 9.4l3.1 2.8-3.1 2.8` |
| kropkę odsunięto od ramki i powiększono | `18,4 · 9,9 · r 1,5` | `17,4 · 10,0 · r 1,85` |

Luz kropka–ramka wzrósł z 0,225 do 0,60 j. (0,40 px przy 16 px). Bez tej korekty
kropka zlewa się z prawą krawędzią ramki.

#### MultitaskingAI 16 px

| Zabieg | Z | Na |
|---|---|---|
| **nie usunięto niczego** | — | — |
| rdzeń powiększono | `r 2,6` | `r 2,9` |
| węzły powiększono | `r 1,9` | `r 2,0` |
| łączniki przeliczono na nowe promienie | `M12 9.4V6.5` · `M9.8 13.4l-3.3 1.9` · `M14.2 13.4l3.3 1.9` | `M12 9.1V6.6` · `M9.55 13.55l-2.95 1.8` · `M14.45 13.55l2.95 1.8` |

**Dlaczego nic nie ubyło.** Usunięcie któregokolwiek węzła zmieniłoby liczbę
wykonawców z trzech na dwóch — a to jest fakt produktowy (Executor 1, Executor 2,
Executor 3 / Validator), nie ozdoba. Usunięcie łączników zamieniłoby orkiestrację
w cztery niezależne punkty. Uproszczenie polegało więc wyłącznie na
**pogrubieniu i przeliczeniu**, nie na ujęciu.

### 6.3 Kontrola czytelności — metoda

Każdy wariant 16 i 20 px zrasteryzowano i obejrzano w powiększeniu 14 ×
(interpolacja `NEAREST`, bez wygładzania). Kryteria przyjęcia:

1. kropka jest widoczna jako **osobna forma** (nie zlana z sąsiadem),
2. sylwetka jest **rozróżnialna od pozostałych trzech** przy jednoczesnym oglądzie,
3. żadna szczelina nie zamyka się całkowicie,
4. znak nie „ciemnieje" — stosunek zadruku do pola pozostaje poniżej 45 %.

---

## 7. Rastry PNG emblematów

### 7.1 Zakres

32 pliki: 4 emblematy × 4 rozmiary (24 / 48 / 96 / 192 px) × 2 warianty barwne.

| Wariant | Barwa | Przeznaczenie |
|---|---|---|
| `…-jasny.png` | `#181818` | osadzenia na tle jasnym: dokumenty, prezentacje, e-mail |
| `…-ciemny.png` | `#ECECEC` | osadzenia na tle ciemnym |

Wszystkie z **kanałem alfa** (tło przezroczyste) i wszystkie z geometrii
kanonicznej 24 j. — rozmiar najmniejszy w komplecie to 24 px, więc uproszczenie
16 px nie ma tu zastosowania.

### 7.2 Kiedy raster, kiedy SVG

| Sytuacja | Format |
|---|---|
| interfejs aplikacji (wpięcie inline, `currentColor`) | **SVG** |
| dokumentacja HTML, prototypy | **SVG inline** |
| dokument biurowy, prezentacja, e-mail, komunikator | **PNG** |
| skrót aplikacji w manifeście (`shortcuts[].icons`) | **PNG 192** |
| druk | **SVG** (raster wyłącznie ≥ 300 dpi w docelowym rozmiarze) |

---

## 8. Ikona aplikacji

### 8.1 Analiza plików zastanych

Katalog `zasoby/marka/ikona-aplikacji/` zawierał:

| Plik | Rozmiar | Ocena |
|---|---|---|
| `ikona-aplikacji.svg` | 1024 × 1024, `rx 224` | **podstawa — zachowana bez zmian** |
| `ikona-maskowalna.svg` | 1024 × 1024, `rx 0`, mniejsza skala | **podstawa — zachowana bez zmian** |
| `ikona-1024/512/192/180.png` | rastry | komplet niepełny — brak 32, 64, 128, 256 |
| `ikona-maskowalna-512.png` | raster | ✓ |
| brak | — | **brak ICO** — Tauri i Windows go wymagają |
| brak | — | **brak rozmiaru 128@2x** — Tauri go wymaga |

### 8.2 Konstrukcja — pomiar

```
   1024 × 1024
   ┌─────────────────────────────────────────────────────────┐  rx 224
   │ ╭─────────────────────────────────────────────────────╮ │  = 21,875 %
   │ │                                                     │ │
   │ │                                                     │ │
   │ │            ▶▶  ●                                    │ │
   │ │           sygnet 96 × 96 · scale 6,6133             │ │
   │ │           translate(194,6 · 194,6)                  │ │
   │ │                                                     │ │
   │ ╰─────────────────────────────────────────────────────╯ │
   └─────────────────────────────────────────────────────────┘
      │←274,0→│                                    │←237,5→│
```

| Wielkość | Wartość | Uwaga |
|---|---:|---|
| Płótno | 1024 × 1024 px | — |
| Promień naroża | 224 px = **21,875 %** | ta sama proporcja co kafel faviconu (21/96) |
| Grunt | `#131313` = `--dn-rama` = `szary-925` | jedyna powierzchnia marki niezmienna w obu motywach |
| Skala sygnetu | 6,6133 | 96 j. → 634,9 px |
| Odsunięcie | `translate(194,6 · 194,6)` | — |
| Pole zadruku (x) | 274,0 – 786,5 px, szer. 512,5 px = **50,1 %** | — |
| Pole zadruku (y) | 366,5 – 657,5 px, wys. 291,0 px | **wyśrodkowane co do dziesiątej px** |
| Margines lewy / prawy | 274,0 / 237,5 px | różnica 36,5 px |
| Oś masy grotów (x) | 472,4 px | 39,6 px w lewo od osi płótna |

**Dlaczego marginesy poziome nie są równe.** Sygnet ma masę skupioną w dwóch
grotach i lekką, oderwaną kropkę po prawej. Wyśrodkowanie pola zadruku
(274,0 = 237,5) zepchnęłoby oś masy grotów o **76 px** w lewo od osi płótna
i ikona wyglądałaby na przesuniętą. Zastana kompozycja jest kompromisem:
oś masy grotów leży 39,6 px w lewo, oś pola zadruku 18,25 px w prawo od osi
płótna. **Kompozycji nie zmieniono** — jest zatwierdzona i spójna z rastrami
już wydanymi.

### 8.3 Pole bezpieczne

| Reguła | Wartość | Wynik |
|---|---:|---|
| minimalny margines od krawędzi płótna | 237,5 px = **23,2 %** | ✓ (zalecenie ≥ 20 %) |
| minimalny margines od łuku naroża | > 240 px | ✓ |
| stosunek zadruku do płótna | ok. 14,2 % powierzchni | ✓ (poniżej progu „ciężkiej" ikony 25 %) |

> **REGUŁA.** W polu bezpiecznym ikony aplikacji nie stawia się niczego:
> ani plakietki wersji, ani napisu „beta", ani liczby powiadomień.
> Odznaki nakłada **system operacyjny**, na swoich zasadach.

### 8.4 Wariant maskowalny — Android adaptive icon

Android przycina ikonę dowolną maską producenta (koło, squircle, „teardrop",
kwadrat zaokrąglony). Bezpieczna jest wyłącznie **centralna strefa o średnicy
80 % płótna**.

| Wielkość | Wartość | Wynik |
|---|---:|---|
| Płótno | 1024 × 1024 px, `rx 0` (grunt do krawędzi) | — |
| Skala sygnetu | 5,5467 (mniejsza niż w ikonie zwykłej) | — |
| Odsunięcie | `translate(245,8 · 245,8)` | — |
| Pole zadruku | 312,4 – 742,2 × 390,0 – 634,1 px, szer. **42,0 %** | — |
| Promień strefy bezpiecznej | 409,6 px (80 % ⇒ r = 0,4 × 1024) | — |
| **Najdalszy punkt zadruku od środka** | **248,4 px** | ✓ |
| Zapas | 161,2 px = **39,3 %** promienia strefy | ✓ z dużym marginesem |

Najdalszym punktem zadruku jest zewnętrzna krawędź kropki sygnału
(środek po transformacji 706,2 · 598,0; promień 36,1 px). Nawet najbardziej
agresywna maska kołowa nie utnie ani fragmentu kropki.

```
        maska producenta (dowolna)
        ┌───────────────────────────────┐
        │        ...............        │
        │     ..                 ..     │  ← strefa bezpieczna Ø 80 %
        │   ..    ▶▶  ●            ..   │     r = 409,6 px
        │  .      zadruk r ≤ 248,4   .  │
        │   ..                     ..   │     zapas 161,2 px
        │     ..                 ..     │
        │        ...............        │
        └───────────────────────────────┘
              grunt do krawędzi (rx 0)
```

### 8.5 Rozmiary Tauri

| Plik wymagany przez Tauri | Odpowiednik w komplecie | Rozmiar |
|---|---|---:|
| `32x32.png` | `ikona-32.png` | 32 |
| `128x128.png` | `ikona-128.png` | 128 |
| `128x128@2x.png` | `ikona-128@2x.png` | 256 |
| `icon.png` | `ikona-1024.png` | 1024 |
| `icon.ico` | `ikona-aplikacji.ico` | 16–256, 7 wpisów |
| `icon.icns` (macOS) | zbuduj z `ikona-1024.png` narzędziem `iconutil` | — |

**Uproszczenie dla rozmiarów ≤ 32 px.** Rastry 32 px budowane są z wariantu
`ikona-uproszczona.svg` — jeden grot i kropka zamiast dwóch grotów. Przy 32 px
kafel ma ok. 8 px marginesu, a znak ok. 16 px szerokości; podwójny grot rozpada
się wtedy na cztery szare kreski. Reguła progu 24 px z kontraktu systemu projektowego obowiązuje
**dla wielkości renderowanego znaku, nie dla wielkości kafla**.

### 8.6 Rozmiary PWA

| Rozmiar | Plik | Rola |
|---:|---|---|
| 192 | `icon-192.png` | ikona instalacji, ekran domowy Android |
| 512 | `icon-512.png` | ekran powitalny, sklepy aplikacji |
| 512 maskowalna | `ikona-maskowalna-512.png` | `purpose: "maskable"` |
| dowolny | `favicon.svg` | `sizes: "any"` — wpis SVG w manifeście |
| 180 | `apple-touch-icon.png` | iOS / iPadOS (poza manifestem) |

### 8.7 ICO aplikacji — grafika dobrana do rozmiaru

`ikona-aplikacji.ico` zawiera **siedem osobnych obrazów**, nie jeden przeskalowany:

| Wpis | Źródło | Powód |
|---:|---|---|
| 16 px | `ikona-uproszczona.svg` | jeden grot |
| 24 px | `ikona-uproszczona.svg` | jeden grot |
| 32 px | `ikona-uproszczona.svg` | jeden grot |
| 48 px | `ikona-aplikacji.svg` | pełny sygnet |
| 64 px | `ikona-aplikacji.svg` | pełny sygnet |
| 128 px | `ikona-aplikacji.svg` | pełny sygnet |
| 256 px | `ikona-aplikacji.svg` | pełny sygnet |

Każdy wpis to skompresowany PNG w kontenerze ICO (obsługiwane przez Windows
Vista+ i wszystkie przeglądarki). Plik: 14 308 B.

---

## 9. Favicon

### 9.1 Trzy nośniki, trzy role

| Nośnik | Plik | Kiedy działa | Tło |
|---|---|---|---|
| **podstawowy** | `favicon.svg` | wszystkie współczesne przeglądarki | przezroczyste; barwy przełącza `prefers-color-scheme` |
| **zapasowy** | `favicon.ico` | starsze przeglądarki, żądanie `/favicon.ico` bez `<link>` | grunt kryjący `#131313` |
| **systemowy** | `apple-touch-icon.png` | iOS / iPadOS „Dodaj do ekranu początkowego" | grunt kryjący, **bez alfy, bez zaokrąglenia** |

### 9.2 SVG adaptacyjny

```
   viewBox 0 0 96 96 · tło przezroczyste
   znak: sygnet UPROSZCZONY (jeden grot + kropka)
     grot   M18 18 H35 L64 48 L35 78 H18 L45 48 Z
     kropka cx 79 · cy 69 · r 9

   .znak   { fill: #181818 }   .kropka { fill: #3B6FE0 }
   @media (prefers-color-scheme: dark) {
     .znak { fill: #ECECEC }   .kropka { fill: #5C8CEC }
   }
```

Favicon renderuje się najczęściej w 16–20 px, czyli **poniżej progu 24 px** —
dlatego wariant uproszczony jest tu wariantem podstawowym, nie awaryjnym.

### 9.3 Kafel rastrowy — dlaczego grunt kryjący

Raster nie umie przełączyć barwy z motywem systemu. Pomiar dla znaku
atramentowego (`#181818`) na pasku kart w motywie ciemnym (typowo `#35363A`)
daje kontrast **1,6 : 1** — znak znika.

Rozstrzygnięcie: rastry faviconu stoją na **kryjącym gruncie `#131313`**
(`--dn-rama`) w kwadracie o promieniu naroża 21 j. na 96 — proporcja **21/96 =
0,21875**, identyczna z `224/1024` ikony aplikacji. Kafel faviconu i ikona
aplikacji to ta sama bryła w dwóch skalach.

| Para | Kontrast | Wynik |
|---|---:|---|
| `#ECECEC` na `#131313` (znak w kaflu) | 16,32 : 1 | ✓ |
| `#131313` na `#F4F4F4` (kafel na jasnym pasku kart) | 15,42 : 1 | ✓ |
| `#131313` na `#35363A` (kafel na ciemnym pasku kart) | 1,84 : 1 | kafel ledwie odcina się od tła — **odcina go znak w środku**, nie krawędź |

### 9.4 ICO wielorozmiarowy

`favicon.ico` — trzy osobne obrazy:

| Wpis | Źródło | Znak |
|---:|---|---|
| 16 px | `favicon-kafel-uproszczony.svg` | jeden grot + kropka, skala 0,72 |
| 32 px | `favicon-kafel.svg` | pełny sygnet, skala 0,62 |
| 48 px | `favicon-kafel.svg` | pełny sygnet, skala 0,62 |

Plik: 2 539 B. Zbudowany bezpośrednim zapisem kontenera ICO (nagłówek
`ICONDIR` + 3 × `ICONDIRENTRY` + trzy bloki PNG), a nie przeskalowaniem
jednego źródła — dzięki temu wpis 16 px ma własną, uproszczoną grafikę.

**Pomiary kompozycji kafla:**

| Wariant | Skala | Pole zadruku (x) | Pole zadruku (y) | Marginesy |
|---|---:|---|---|---|
| pełny | 0,62 | 23,975 – 72,025 | 34,360 – 61,640 | L 23,98 · P 23,97 · G 34,36 · D 34,36 |
| uproszczony | 0,72 | 22,800 – 73,200 | 26,400 – 69,600 | L 22,80 · P 22,80 · G 26,40 · D 26,40 |

Wariant uproszczony jest **wyśrodkowany dokładnie w obu osiach** — sylwetka
„jeden grot + kropka" nie ma oderwanej masy po prawej w takim stopniu jak
sygnet pełny, więc kompensacja optyczna nie jest potrzebna.

### 9.5 apple-touch-icon

| Cecha | Wartość | Uzasadnienie |
|---|---|---|
| Rozmiar | 180 × 180 px | wymóg iOS dla ekranów @3x |
| Promień naroża | **0 (kwadrat)** | maskę nakłada system; własne zaokrąglenie daje podwójny łuk |
| Kanał alfa | **usunięty** (`RGB`) | iOS nie obsługuje przezroczystości w tej ikonie — kompozytuje ją na czarno |
| Grunt | `#131313` | jak ikona aplikacji |
| Znak | sygnet pełny, kompozycja jak w ikonie aplikacji | — |

### 9.6 Wady plików zastanych — wykryte i naprawione

Trzy defekty wykryte pomiarowo w plikach `zasoby/marka/favicon/`:

| # | Plik zastany | Defekt | Pomiar | Naprawa |
|---|---|---|---|---|
| **1** | `favicon-16.png` | **kropka sygnału nie istnieje** — przy 16 px promień 9 j. na 96 daje 1,5 px i rasteryzator gasi ją do zera | histogram: **0 pikseli błękitu**; barwy w pliku wyłącznie `#181818` + alfa | kafel uproszczony ze skalą 0,72 → **2 piksele błękitu** przy 16 px |
| **2** | `favicon-16/32/48.png` | tło przezroczyste + znak atramentowy — znak niewidoczny na ciemnym pasku kart | kontrast `#181818` / `#35363A` = **1,6 : 1** | grunt kryjący `#131313`, znak `#ECECEC` — 16,32 : 1 |
| **3** | `apple-touch-icon.png` | identyczny z `ikona-180.png`, czyli z wypalonym `rx 224`; dodatkowo **białe naroża** po spłaszczeniu alfy | piksel (0,0) = `#FFFFFF` | kwadrat `rx 0`, tryb `RGB`, piksel (0,0) = `#131313` |

> Naprawy dotyczą **wyłącznie plików wytworzonych w `03-marka/emblematy/`**.
> Pliki w `zasoby/marka/` pozostawiono nietknięte — ich podmiana jest decyzją
> wdrożeniową, nie projektową.

---

## 10. Snippet `<head>` i manifest

### 10.1 Snippet — gotowy do skopiowania

Plik: `emblematy/favicon/naglowek-snippet.html`

```html
<!-- 1. Nośnik podstawowy: SVG adaptacyjny -->
<link rel="icon" href="/favicon.svg" type="image/svg+xml">

<!-- 2. Zapas rastrowy: ICO 16/32/48 na kryjącym gruncie -->
<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 48x48">

<!-- 3. iOS / iPadOS — kwadrat bez zaokrąglenia i bez alfy -->
<link rel="apple-touch-icon" href="/apple-touch-icon.png">

<!-- 4. Manifest aplikacji internetowej -->
<link rel="manifest" href="/site.webmanifest">

<!-- 5. Barwa paska systemowego -->
<meta name="theme-color" content="#F4F4F4" media="(prefers-color-scheme: light)">
<meta name="theme-color" content="#131313" media="(prefers-color-scheme: dark)">

<!-- 6. Kafel Windows -->
<meta name="msapplication-TileColor" content="#131313">
<meta name="msapplication-TileImage" content="/icon-192.png">

<!-- 7. Nazwa aplikacji w trybie samodzielnym -->
<meta name="application-name" content="Danaco Console">
<meta name="apple-mobile-web-app-title" content="Danaco Console">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
```

**Kolejność jest znacząca.** Przeglądarka obsługująca SVG bierze wpis pierwszy
i ignoruje ICO. Starsza pomija SVG i pobiera ICO. Odwrócenie kolejności sprawia,
że część przeglądarek utknie na rastrze.

**Dlaczego dwa `theme-color`.** W motywie jasnym pasek systemowy przyjmuje
`#F4F4F4` (papier roboczy `--dn-tlo`), w ciemnym `#131313` (rama kokpitu
`--dn-rama`). Rama jest w systemie zawsze atramentowa, więc w motywie ciemnym
pasek systemowy przedłuża pasek górny aplikacji bez szwu.

### 10.2 site.webmanifest

Plik: `emblematy/favicon/site.webmanifest`

| Pole | Wartość | Źródło |
|---|---|---|
| `name` | `Danaco Console — AI Operating Environment` | kontrakt systemu projektowego |
| `short_name` | `Danaco Console` | — |
| `lang` | `pl` | kontrakt systemu projektowego |
| `background_color` | `#0F0F0F` | `--dn-tlo` (motyw ciemny) = `szary-950` |
| `theme_color` | `#131313` | `--dn-rama` = `szary-925` |
| `display` | `standalone` | — |
| `display_override` | `["window-controls-overlay", "standalone"]` | pasek górny aplikacji przejmuje pas tytułowy okna |
| `icons` | 192, 512, 512 maskowalna, SVG `any` | rozdz. 8.6 |
| `shortcuts` | cztery środowiska: TalkIn, WorkSpace, CodeStudio, MultitaskingAI | kontrakt systemu projektowego |

**Skróty do środowisk** używają rastrów emblematów 192 px — to jedyne miejsce
w systemie, w którym emblemat środowiska pełni rolę ikony na poziomie systemu
operacyjnego. Uzasadnienie: skrót prowadzi **do środowiska**, nie do aplikacji;
gdyby wszystkie cztery skróty miały ikonę aplikacji, byłyby nierozróżnialne.

---

## 11. Zasady doboru — sygnet, emblemat, ikona aplikacji

### 11.1 Drzewo decyzyjne

```
   Co oznaczasz tym znakiem?
   │
   ├─ PRODUKT / FIRMĘ (nadawcę)                    ──►  SYGNET  (lub logotyp)
   │   pasek górny · dokument · wizytówka · stopka
   │   ekran powitalny · list · faktura
   │   │
   │   └─ mniej niż 24 px?                         ──►  SYGNET UPROSZCZONY
   │
   ├─ ŚRODOWISKO (gdzie użytkownik jest)           ──►  EMBLEMAT ŚRODOWISKA
   │   karta środowiska · boczna nawigacja
   │   tytuł powłoki · karta sesji · skrót systemowy
   │
   ├─ MODUŁ / NARZĘDZIE (czym pracuje)             ──►  IKONA Z zasoby/ikony/
   │   pozycja bocznej nawigacji · przycisk
   │   zakładka · pozycja menu
   │
   └─ APLIKACJĘ W SYSTEMIE (plik wykonywalny)      ──►  IKONA APLIKACJI
       pulpit · pasek zadań · ekran domowy
       lista aplikacji · karta przeglądarki (favicon)
```

### 11.2 Tabela miejsc

| Miejsce | Znak | Rozmiar | Barwa |
|---|---|---:|---|
| Pasek górny aplikacji | sygnet | 24–28 px | `--dn-rama-tekst` + kropka `--dn-sygnal-500` |
| Ekran startowy (okno ładowania) | logo pionowy | ≥ 96 px | barwy własne |
| Karta środowiska (strona główna) | **emblemat** | 40 px | `--dn-tekst` |
| Boczna nawigacja — nagłówek środowiska | **emblemat** | 20 px | `--dn-tekst-2` |
| Boczna nawigacja — pozycja modułu | ikona modułu | 20 px | `--dn-tekst-2` / `--dn-tekst` |
| Pas kart sesji | **emblemat** | 16 px | `--dn-tekst-2` |
| Panel orkiestracji (MultitaskingAI) | **emblemat** MultitaskingAI | 20 px | `--dn-tekst` |
| Modal, toast, tooltip | ikona z zestawu | 16 px | wg klasy semantycznej |
| Karta przeglądarki | **favicon** | 16–20 px | kafel `#131313` |
| Pasek zadań / dok | **ikona aplikacji** | 32–64 px | kafel `#131313` |
| Ekran domowy telefonu | **ikona aplikacji** (maskowalna) | 192 px | kafel `#131313` |
| Lista aplikacji systemu | **ikona aplikacji** | 48–128 px | kafel `#131313` |
| Skrót do środowiska (manifest) | **emblemat** | 192 px | `#ECECEC` |
| Papier firmowy, faktura, wizytówka | logotyp / sygnet | wg makiety | barwy własne |

### 11.3 Sześć zakazów [NIENEGOCJOWALNE]

1. **Emblemat środowiska nie jest znakiem marki** — nie stanie na papierze
   firmowym, wizytówce, fakturze ani w stopce dokumentu.
2. **Sygnet nie jest ikoną interfejsu** — nie stanie w bocznej nawigacji jako
   pozycja modułu ani w przycisku.
3. **Ikona aplikacji nie jest godłem wewnątrz aplikacji** — kafel z gruntem
   `#131313` nie pojawia się w treści okna; wewnątrz stoi sam sygnet.
4. **Nie miesza się pięter** — cztery karty środowisk mają cztery emblematy,
   nie cztery sygnety i nie cztery ikony aplikacji.
5. **Nie dorabia się piątego emblematu** — środowiska są cztery
   (kontrakt systemu projektowego). Komponenty własne strefy 2 mają **kafle z ikonami**,
   nie emblematy.
6. **Nie przebarwia się samej kropki emblematu** (rozdz. 5.2).

---

## 12. Podgląd w kontekście

### 12.1 Karta przeglądarki

```
   ┌──────────────────────────────────────────────────────────────┐
   │ ▣ Danaco Console — Centrum dowo…  ✕ │  + │                   │  ← favicon 16 px
   ├──────────────────────────────────────────────────────────────┤
   │ ←  →  ⟳   console.danaco-group.pl/                           │
   └──────────────────────────────────────────────────────────────┘

   ▣ = kafel #131313 · jeden grot #ECECEC · kropka #5C8CEC (2 px błękitu)
```

Kontrola: przy 16 px w kaflu widoczne są **trzy formy** — grunt, grot, kropka.
W pliku zastanym widoczne były dwie (kropka nie istniała — defekt pierwszy).

### 12.2 Pasek zadań

```
   ┌──────────────────────────────────────────────────────────────┐
   │  ▣    ▣    ▣    ▣    ▣                                      │  ← ikona 32 px
   │  ‾‾‾‾                                                        │     (wariant
   └──────────────────────────────────────────────────────────────┘      uproszczony)
      ↑
      Danaco Console — aktywna, wskaźnik pod kaflem
```

Przy 32 px kafel ma 224/1024 × 32 ≈ 7 px promienia naroża i ok. 16 px znaku —
stąd wariant z jednym grotem.

### 12.3 Ekran domowy telefonu

```
   ┌───────────────────────────┐
   │  9:41              ▮▮▮ ⌁  │
   │                           │
   │   ▢     ▢     ▢     ▢     │
   │                           │
   │   ▢    ╔═══╗  ▢     ▢     │  ← ikona maskowalna 192 px
   │        ║▶▶ ●║               (maska producenta: koło / squircle /
   │        ╚═══╝                 kwadrat zaokrąglony — zadruk
   │       Danaco                 w strefie r ≤ 248,4 z 409,6)
   │       Console
   │                           │
   └───────────────────────────┘
```

### 12.4 Lista aplikacji systemu

```
   ┌──────────────────────────────────────────────────────────────┐
   │  ▣  Danaco Console                          2,4.0   Otwórz   │  ← 48 px
   │     AI Operating Environment · Danaco Holding Group          │
   ├──────────────────────────────────────────────────────────────┤
   │  ▣  Danaco Console (Mobile)                 2,4.0   Otwórz   │
   └──────────────────────────────────────────────────────────────┘
```

Przy 48 px kafel niesie **pełny sygnet** — dwa groty i kropkę.

### 12.5 Cztery karty środowisk — strona główna

```
   ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
   │ ▭●           │ │ ▫▫           │ │ ▭>_          │ │  ◯           │  ← emblemat
   │              │ │ ▫◉           │ │              │ │ ◯●◯          │     40 px
   │ TalkIn       │ │ WorkSpace    │ │ CodeStudio   │ │ Multitask…AI │
   │ 9 modułów    │ │ 9 modułów    │ │ 8 modułów    │ │ panel orkie… │
   └──────────────┘ └──────────────┘ └──────────────┘ └──────────────┘
     wstęga 2 px --dn-kropka pojawia się na :hover i [aria-current]
```

Emblematy są tu **atramentowe** (`--dn-tekst`). Sygnał niesie wstęga górna karty,
nie emblemat — zgodnie z rozstrzygnięciem rozdz. 5.2.

---

## 13. Inwentarz wytworzonych plików

Katalog: `/home/claude/danaco/WYNIK/03-marka/emblematy/` — **102 pliki**.

### 13.1 `emblematy/svg/` — 24 pliki

```
srodowisko-talkin-{16,20,24,32,48,64}.svg
srodowisko-workspace-{16,20,24,32,48,64}.svg
srodowisko-codestudio-{16,20,24,32,48,64}.svg
srodowisko-multitaskingai-{16,20,24,32,48,64}.svg
```

Wariant interfejsowy: `fill="none"`, `stroke="currentColor"`, kropka
`fill="currentColor"`. Warianty 16 px uproszczone wg rozdz. 6.2.

### 13.2 `emblematy/svg/warianty/` — 16 plików

```
srodowisko-{talkin,workspace,codestudio,multitaskingai}-jasny.svg          #181818
srodowisko-{talkin,workspace,codestudio,multitaskingai}-ciemny.svg         #ECECEC
srodowisko-{talkin,workspace,codestudio,multitaskingai}-sygnal-jasny.svg   #2457C9
srodowisko-{talkin,workspace,codestudio,multitaskingai}-sygnal-ciemny.svg  #8FB2F5
```

### 13.3 `emblematy/png/` — 32 pliki

```
srodowisko-<środowisko>-{24,48,96,192}-{jasny,ciemny}.png
```

### 13.4 `emblematy/ikona-aplikacji/` — 16 plików

| Plik | Bajtów | Rola |
|---|---:|---|
| `ikona-aplikacji.svg` | 417 | źródło — kafel `rx 224` |
| `ikona-maskowalna.svg` | 415 | źródło — Android adaptive, `rx 0` |
| `ikona-kwadratowa.svg` | 415 | źródło — apple-touch-icon, `rx 0`, pełna skala |
| `ikona-uproszczona.svg` | 335 | źródło dla rastrów ≤ 32 px |
| `ikona-32.png` | 826 | Tauri `32x32.png` |
| `ikona-64.png` | 1 848 | Windows, Linux |
| `ikona-128.png` | 3 313 | Tauri `128x128.png` |
| `ikona-128@2x.png` | 5 549 | Tauri `128x128@2x.png` |
| `ikona-180.png` | 3 995 | iOS (kafel zaokrąglony) |
| `ikona-192.png` | 4 458 | PWA |
| `ikona-256.png` | 5 817 | Windows |
| `ikona-512.png` | 11 295 | PWA, sklepy |
| `ikona-1024.png` | 35 710 | źródło `icon.png`, macOS `.icns` |
| `ikona-maskowalna-192.png` | 2 504 | Android adaptive |
| `ikona-maskowalna-512.png` | 5 494 | `purpose: "maskable"` |
| `ikona-aplikacji.ico` | 14 308 | 7 wpisów: 16, 24, 32, 48, 64, 128, 256 |

### 13.5 `emblematy/favicon/` — 13 plików

| Plik | Bajtów | Rola |
|---|---:|---|
| `favicon.svg` | 418 | nośnik podstawowy, adaptacyjny |
| `favicon-kafel.svg` | 427 | kafel z pełnym sygnetem (źródło rastrów 32/48) |
| `favicon-kafel-uproszczony.svg` | 352 | kafel z jednym grotem (źródło rastru 16) |
| `favicon-16.png` | 408 | — |
| `favicon-32.png` | 776 | — |
| `favicon-48.png` | 1 301 | — |
| `favicon.ico` | 2 539 | 3 wpisy: 16, 32, 48 |
| `apple-touch-icon.png` | 2 190 | 180 px, `rx 0`, RGB bez alfy |
| `icon-192.png` | 4 458 | PWA |
| `icon-512.png` | 11 295 | PWA |
| `ikona-maskowalna-512.png` | 5 494 | PWA maskowalna |
| `site.webmanifest` | 1 628 | manifest |
| `naglowek-snippet.html` | 2 062 | snippet `<head>` |

### 13.6 `emblematy/generator-godel.py` — 1 plik

Odtwarza cały komplet: `python3 generator-godel.py`. Zależności: `cairosvg`,
`Pillow`. Zawiera kontrolę końcową — liczy piksele błękitu w rastrach faviconu
i przerywa, jeżeli kropka sygnału zniknęła (regresja pierwszego defektu).

---

## 14. Kontrola jakości

| # | Sprawdzenie | Wynik |
|---|---|---|
| 1 | krzywe sygnetu i emblematów nietknięte | ✓ — warianty ≥ 24 px różnią się wyłącznie atrybutami `width`/`height` |
| 2 | dokładnie jedna wypełniona kropka w każdym emblemacie, w każdym rozmiarze | ✓ — 24 pliki sprawdzone |
| 3 | zero `#000000` | ✓ — najciemniejsza wartość `#0F0F0F` |
| 4 | zero gradientów w godłach | ✓ |
| 5 | emblematy dziedziczą `currentColor` | ✓ — 24 pliki bazowe |
| 6 | kropka nie odrywa się barwą od obrysu | ✓ — 16 wariantów barwnych, każdy jednobarwny |
| 7 | kropka sygnału widoczna w faviconie 16 px | ✓ — 2 px błękitu (przed: 0) |
| 8 | zadruk ikony maskowalnej w strefie 80 % | ✓ — 248,4 z 409,6 px |
| 9 | `apple-touch-icon` bez alfy i bez zaokrąglenia | ✓ — RGB, piksel (0,0) = `#131313` |
| 10 | ICO wielorozmiarowe czytelne przez PIL i przeglądarki | ✓ — 3 i 7 wpisów odczytanych |
| 11 | manifest poprawny składniowo (JSON) | ✓ |
| 12 | nazwy własne środowisk zgodne z dokumentacją | ✓ — TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| 13 | nazwy plików kebab-case | ✓ |
| 14 | polskie diakrytyki poprawne | ✓ |

---

## 15. Decyzje projektowe

### Emblemat pozostaje jednobarwny — kropka nie dostaje błękitu

**Rozstrzygnięcie.** Kropka w emblemacie środowiska jest cytatem **formy**
sygnetu, nie cytatem **barwy**. Cały emblemat ma jedną barwę, także kropka.

**Uzasadnienie.** Realne pliki `srodowisko-*.svg` mają `fill="currentColor"`
na kropce, a `komponenty.css` ustawia `.dn-karta-srodowiska-godlo { color:
var(--dn-tekst) }` i sygnalizuje aktywność **kreską 2 px obok ikony**, nie
przebarwieniem. Do 21 emblematów na jednym ekranie oznaczałoby do 21 punktów
sygnału — barwa przestałaby cokolwiek znaczyć.

**Konsekwencja.** Rastry PNG emblematów są jednobarwne. Nie istnieje plik
z niebieską kropką w emblemacie.

### Obrys poniżej 24 px rośnie do 1,90 j.

**Rozstrzygnięcie.** Warianty 16 i 20 px mają obrys 1,90 j. zamiast kanonicznych
1,75 j.

**Uzasadnienie.** 1,75 j. przy 16 px daje 1,17 px — linię cieńszą od piksela,
którą rasteryzator rozmywa na dwa rzędy o połowicznym kryciu. Znak szarzeje
i traci kontur. 1,90 j. daje 1,27 px z wyraźniejszym rdzeniem. Wartość kanoniczna
obowiązuje od 24 px w górę — czyli wszędzie tam, gdzie kontrakt systemu projektowego ją
przewiduje.

**Czego nie zrobiono.** Nie zmieniono żadnej współrzędnej w wariantach 20 px —
wyłącznie grubość obrysu.

### WorkSpace 16 px traci obrys czwartego modułu

**Rozstrzygnięcie.** W wariancie 16 px emblemat WorkSpace ma trzy obrysowe moduły
i jedną kropkę o promieniu 3,0 j. w miejscu czwartego.

**Uzasadnienie.** Zweryfikowano rasteryzacją: wnętrze modułu przy 16 px ma
3,3 px, kropka 2,5 px — obrys i kropka sklejają się w plamę. Po usunięciu obrysu
emblemat czyta się jako „trzy moduły i jeden aktywny" — znaczenie nienaruszone,
liczba wypełnionych kropek nadal dokładnie jedna.

**Wariant odrzucony.** Wypełnienie całego czwartego modułu — łamie zasadę
„dokładnie jedna wypełniona **kropka**" (kwadrat nie jest kropką).

### MultitaskingAI nie traci przy 16 px żadnego elementu

**Rozstrzygnięcie.** Wariant 16 px zachowuje rdzeń, trzy węzły i trzy łączniki;
uproszczenie polega wyłącznie na pogrubieniu i przeliczeniu łączników.

**Uzasadnienie.** Liczba wykonawców (Executor 1, Executor 2, Executor 3 /
Validator) jest faktem produktowym z kontraktu systemu projektowego. Usunięcie węzła zmienia
treść znaku. Usunięcie łączników zamienia orkiestrację w cztery niezależne punkty.

### Kompozycji ikony aplikacji nie zmieniono mimo nierównych marginesów

**Rozstrzygnięcie.** Zachowano zastane `translate(194,6 · 194,6) scale(6,6133)`,
choć margines lewy (274,0 px) różni się od prawego (237,5 px) o 36,5 px.

**Uzasadnienie.** Wyśrodkowanie pola zadruku zepchnęłoby oś masy dwóch grotów
o 76 px w lewo od osi płótna — sygnet ma masę skupioną po lewej i lekką,
oderwaną kropkę po prawej. Zastana kompozycja jest kompromisem między osią pola
zadruku a osią masy. Zmiana rozjechałaby się też z rastrami już wydanymi
w `zasoby/marka/`.

**Odnotowano jako obserwację, nie jako defekt.**

### Rastry ikony aplikacji ≤ 32 px używają wariantu uproszczonego

**Rozstrzygnięcie.** `ikona-32.png` i wpisy ICO 16 / 24 / 32 budowane są
z `ikona-uproszczona.svg` (jeden grot).

**Uzasadnienie.** Próg 24 px z kontraktu systemu projektowego dotyczy **renderowanego znaku**,
nie kafla. Przy kaflu 32 px znak ma ok. 16 px szerokości — podwójny grot rozpada
się na cztery szare kreski. To ta sama reguła, nie odstępstwo od niej.

### Rastry faviconu stoją na kryjącym gruncie

**Rozstrzygnięcie.** `favicon-16/32/48.png` i `favicon.ico` to kafle `#131313`
ze znakiem `#ECECEC` i kropką `#5C8CEC`. `favicon.svg` pozostaje przezroczysty
i adaptacyjny.

**Uzasadnienie.** Raster nie przełączy barwy z motywem systemu. Znak atramentowy
na ciemnym pasku kart daje 1,6 : 1 — znika. Grunt kryjący daje 16,32 : 1
niezależnie od paska. Promień naroża kafla 21/96 = 0,21875 jest identyczny
z 224/1024 ikony aplikacji — to ta sama bryła w dwóch skalach.

**Koszt świadomy.** Favicon ma widoczny kwadrat. Przyjęty, bo alternatywą jest
znak niewidoczny dla połowy użytkowników.

### ICO budowane bezpośrednim zapisem kontenera

**Rozstrzygnięcie.** Oba pliki ICO powstają przez zapis nagłówka `ICONDIR`
i tablicy `ICONDIRENTRY` z osobnymi blokami PNG, a nie przez
`Image.save(..., sizes=[...])`.

**Uzasadnienie.** Zapis przez `sizes` przeskalowuje **jedno** źródło, więc wpis
16 px byłby zmniejszoną wersją pełnego sygnetu — nieczytelną. Zapis bezpośredni
pozwala dać każdemu rozmiarowi własną grafikę, co jest oryginalnym powodem
istnienia formatu ICO.

### apple-touch-icon jest kwadratem bez alfy

**Rozstrzygnięcie.** 180 × 180 px, `rx 0`, tryb `RGB`.

**Uzasadnienie.** iOS nakłada własną maskę — własne zaokrąglenie daje podwójny
łuk. iOS nie obsługuje przezroczystości w tej ikonie i kompozytuje ją na czarno,
co przy zastanym pliku dawało białe naroża — defekt trzeci.

### Manifest niesie skróty do czterech środowisk

**Rozstrzygnięcie.** `shortcuts[]` zawiera cztery pozycje z rastrami emblematów
192 px.

**Uzasadnienie.** Skrót prowadzi do środowiska, nie do aplikacji — cztery skróty
z tą samą ikoną aplikacji byłyby nierozróżnialne. To jedyne miejsce, w którym
emblemat pełni rolę ikony na poziomie systemu operacyjnego; zakaz z rozdz. 11.3
pkt 1 dotyczy materiałów firmowych, nie interfejsu systemowego.

### Plików w `zasoby/marka/` nie podmieniono

**Rozstrzygnięcie.** Trzy wykryte defekty naprawiono wyłącznie
w plikach wytworzonych w `03-marka/emblematy/`.

**Uzasadnienie.** Podmiana zasobów wspólnych zmieniłaby wygląd wszystkich
prototypów i opracowań innych zespołów w trakcie pracy. Zalecenie podmiany
odnotowano; decyzja należy do wdrożenia.

### Komplet odtwarzalny generatorem

**Rozstrzygnięcie.** Wszystkie 101 plików graficznych i konfiguracyjnych powstaje
z jednego skryptu, w którym geometria źródłowa jest wpisana raz, jako stała.

**Uzasadnienie.** Ręczne utrzymanie 32 rastrów, 40 plików SVG i 2 kontenerów ICO
gwarantuje rozjazd przy pierwszej korekcie. Generator zawiera też kontrolę
regresji — liczy piksele błękitu w rastrach faviconu i sygnalizuje zanik kropki.

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
