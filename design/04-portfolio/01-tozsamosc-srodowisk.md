# Danaco Console — Plansza portfolio P1: Tożsamość czterech środowisk

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Warstwa funkcjonalna** | Danaco Pilot — Platforma AI Workspace OS (dokumentacja projektowa v1.0) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Rodzaj opracowania** | Plansza portfolio design identity (nie dokumentacja techniczna) |
| **Odbiorcy** | Właściciel · Designer · Deweloper · odbiorca portfolio |
| **Zakres** | Tożsamość wizualna czterech środowisk platformy: emblematy, skład nawigacji, powłoki, macierz modułów, ścieżka wejścia, uzasadnienie wspólnej palety |
| **Czego NIE zawiera** | Wnętrza okien operacyjnych modułów, wnętrza okien ról MultitaskingAI, okna Konfiguracji i Ustawień — te mają własne prototypy w `05-okna/` |
| **Plik towarzyszący** | `04-portfolio/01-tozsamosc-srodowisk.html` |

---

## Spis treści

1. [Czym jest ta plansza](#1-czym-jest-ta-plansza)
2. [Cztery karty tożsamości](#2-cztery-karty-tożsamości)
3. [Anatomia czterech emblematów](#3-anatomia-czterech-emblematów)
4. [Prawo kropki — cztery położenia](#4-prawo-kropki--cztery-położenia)
5. [Porównanie powłok — boczna nawigacja wobec panelu orkiestracji](#5-porównanie-powłok--boczna-nawigacja-wobec-panelu-orkiestracji)
6. [Macierz moduł × środowisko](#6-macierz-moduł--środowisko)
7. [Arytmetyka macierzy — co z niej wynika](#7-arytmetyka-macierzy--co-z-niej-wynika)
8. [Ścieżka wejścia do środowiska](#8-ścieżka-wejścia-do-środowiska)
9. [Jeden język wizualny dla czterech środowisk](#9-jeden-język-wizualny-dla-czterech-środowisk)
10. [Kafle wejściowe strefy 1 — forma realna](#10-kafle-wejściowe-strefy-1--forma-realna)
11. [Decyzje projektowe planszy](#11-decyzje-projektowe-planszy)
12. [Źródła](#12-źródła)
13. [Kontrola jakości](#13-kontrola-jakości)

---

## 1. Czym jest ta plansza

Plansza P1 pokazuje **tożsamość czterech środowisk Danaco Console w działaniu** — nie opisuje ich, lecz zestawia obok siebie tak, aby różnica była widoczna, zanim czytelnik przeczyta którykolwiek napis.

Platforma ma cztery środowiska: **TalkIn**, **WorkSpace**, **CodeStudio**, **MultitaskingAI**. Trzy pierwsze są środowiskami modułowymi — ich boczna nawigacja jest listą modułów. Czwarte nie ma modułów: zamiast bocznej nawigacji udostępnia **panel orkiestracji** o sześciu sekcjach. Ta jedna różnica strukturalna jest osią całej planszy.

Druga oś: **środowisk nie rozróżnia barwa.** Wszystkie cztery pracują na jednej monochromatycznej palecie i jednym błękicie sygnałowym. Rozróżnia je wyłącznie **emblemat** i **skład nawigacji**. Uzasadnienie tej decyzji — rozdział 9.

### 1.1. Trójstopniowa hierarchia platformy

```
ŚRODOWISKO      „w jakim trybie pracuję?"       TalkIn · WorkSpace · CodeStudio · MultitaskingAI
      ▼
MODUŁ           „jakie zadanie wykonuję?"       15 modułów (albo ROLA w MultitaskingAI)
      ▼
OKNO OPERACYJNE „jakim narzędziem realizuję?"   Chat Window (wspólne) + okna właściwe modułowi
```

Źródło: KANON rozdz. 6 · `projekt-ui/README.md` rozdz. 2.

### 1.2. Siedem zagadnień planszy

| # | Sekcja planszy | Co zestawia | Źródło |
|---|---|---|---|
| 1 | Karty tożsamości | emblemat, nazwa, istota, liczba i lista modułów, charakter pracy, różnica | KANON 6 · `elementy-okien.md` 2.3, 3.4 |
| 2 | Anatomia emblematu | metafora, konstrukcja na siatce 24, położenie kropki | `03-marka/emblematy-i-ikony.md` 2.3, 3 |
| 3 | Porównanie powłok | cztery szkielety obok siebie, wymiary z żetonów | `elementy-okien.md` 3.1, 4.2, 4.6 |
| 4 | Macierz moduł × środowisko | 15 modułów × 4 środowiska, sprzężone podświetlanie | KANON 6 · `README.md` 5 |
| 5 | Ścieżka wejścia | Strefa 1 → karta → powłoka | `elementy-okien.md` 2.6 |
| 6 | Jeden język wizualny | dlaczego jedna paleta i jeden sygnał | `KIERUNEK.md` 3.1, 3.4, 4 |
| 7 | Kafle strefy 1 | `.dn-karta-srodowiska` w formie realnej | `komponenty.css` · `elementy-okien.md` 2.3 |

---

## 2. Cztery karty tożsamości

### 2.1. Zestawienie zbiorcze

| | **TalkIn** | **WorkSpace** | **CodeStudio** | **MultitaskingAI** |
|---|---|---|---|---|
| **Emblemat** | `srodowisko-talkin.svg` | `srodowisko-workspace.svg` | `srodowisko-codestudio.svg` | `srodowisko-multitaskingai.svg` |
| **Motto** | Myśl. Analizuj. Rozumiej. | Planuj. Organizuj. Realizuj. | Projektuj. Buduj. Rozwijaj. | Deleguj. Koordynuj. Nadzoruj. |
| **Opis trybu pracy na karcie** | Wiedza, komunikacja i praca z treścią | Produktywność, organizacja i realizacja projektów | Programowanie | Orkiestracja autonomicznej pracy ciągłej |
| **Kolejność w strefie 1** | 1 | 2 | 3 | 4 |
| **Boczny panel** | boczna nawigacja modułów | boczna nawigacja modułów | boczna nawigacja modułów | **panel orkiestracji** |
| **Liczba pozycji** | **9** | **9** | **8** | **6 sekcji** |
| **Jednostka organizacji pracy** | pojedynczy użytkownik wspierany przez AI | pojedynczy użytkownik wspierany przez AI | pojedynczy użytkownik wspierany przez AI | zespół modeli i agentów realizujących wspólny proces |
| **Podstawowa jednostka pracy** | moduł → okno operacyjne | moduł → okno operacyjne | moduł → okno operacyjne | rola → akcja silnika kolejek → orkiestracja |
| **Co reprezentuje karta sesji** | jeden moduł otwarty w danej chwili | jw. | jw. | jeden proces orkiestracji |
| **Domyślny tryb pracy** | interakcja konwersacyjna, jednorazowa lub sesyjna | jw. | jw. | zdolność do pracy ciągłej po spięciu z Automations |
| **Rola Always On Display** | doradztwo kontekstowe | doradztwo kontekstowe | doradztwo kontekstowe | doradztwo kontekstowe **oraz** obserwator albo operator procesu |

Motta pochodzą z manifestu ikon (`zasoby/ikony/manifest.json`, pole `zastosowanie` dla czterech pozycji `srodowisko-*`). Opisy trybu pracy — z `elementy-okien.md` Tabela 4. Pozostałe wiersze — z `srodowiska/multitaskingai.md` rozdz. 1.1.

### 2.2. Skład bocznej nawigacji — kolejność źródłowa

Kolejność wyświetlania jest wiążąca; nie jest alfabetyczna ani przypadkowa.

| Środowisko | Poz. | Lista modułów w kolejności wyświetlania |
|---|:---:|---|
| **TalkIn** | 9 | Studio · Workspace · Browser · Research · Library · Translate · Roundtable · Assistant · Agents |
| **WorkSpace** | 9 | Studio · Workspace · Browser · Research · Library · Roundtable · Design · Apps · Agents |
| **CodeStudio** | 8 | Workspace · Roundtable · Design · Terminal · Developer · Diagnostics · Apps · Agents |
| **MultitaskingAI** | 6 | Zespoły · Role · Kolejki · Orkiestracja · Harmonogram i automatyki · Monitor procesu |

Źródło: `elementy-okien.md` Tabela 11 (środowiska modułowe) i Tabela 16 (panel orkiestracji).
Uwaga źródłowa: kolejność i widoczność sekcji panelu orkiestracji są **konfigurowalne**; powyższa lista to wartość domyślna (`elementy-okien.md` Makieta 17).

### 2.3. Czym każde środowisko różni się od pozostałych

| Środowisko | Różnica wynikająca wprost z macierzy |
|---|---|
| **TalkIn** | Jedyne środowisko z modułami **Translate** i **Assistant** — oba nie występują nigdzie indziej. Najdłuższa lista treściowa: Studio, Research, Library, Browser, Translate. |
| **WorkSpace** | Jedyne środowisko łączące komplet modułów treściowych (Studio, Research, Library, Browser) z wytwórczymi (**Design**, **Apps**). Most między TalkIn a CodeStudio — dzieli 7 modułów z TalkIn i 5 z CodeStudio. |
| **CodeStudio** | Jedyne środowisko z **Terminal**, **Developer**, **Diagnostics**. Jedyne, które **nie ma** Studio, Research, Library ani Browser. Najkrótsza lista (8). |
| **MultitaskingAI** | Jedyne środowisko **bez bocznej nawigacji modułów**. Nie występuje w macierzy jako kolumna udostępniająca moduły — operuje na rolach, kolejkach i orkiestracji. |

### 2.4. Wspólny rdzeń trzech środowisk modułowych

Trzy moduły występują we **wszystkich** trzech środowiskach modułowych:

| Moduł | Okno wiodące | Dlaczego wszędzie |
|---|---|---|
| **Roundtable** | Model Panels | współpraca wielu modeli nad jednym problemem jest potrzebna w każdym trybie pracy |
| **Workspace** | Project Dashboard | izolowana przestrzeń projektowa jest ramą dla każdego rodzaju pracy; jednocześnie komponent własny strefy 2 („Projekt") |
| **Agents** | Agent Builder | fabryka ekspertów zasila każde środowisko; w MultitaskingAI występuje jako rola/ekspert |

Czwarty moduł obecny we wszystkich trzech, ale **bez okna w bocznej nawigacji**: **Automations** — konfigurowany wyłącznie ze strony głównej (strefa 2), a gotowe automatyki trafiają do sesji jako komponenty własne.

---

## 3. Anatomia czterech emblematów

### 3.1. Gramatyka wspólna — pięć reguł

| # | Reguła | Wartość | Zakres |
|---|---|---|---|
| 1 | Siatka | 24 × 24 jednostki | emblematy + ikony modułów |
| 2 | Obrys | 1,75 j., `linecap: round`, `linejoin: round` | emblematy + ikony modułów |
| 3 | Wypełnienie | `none` na obrysach | emblematy + ikony modułów |
| 4 | Barwa | `currentColor` — jedna barwa na cały znak | emblematy + ikony modułów |
| 5 | **Kropka sygnału** | dokładnie **jedna** wypełniona kropka | **wyłącznie** emblematy i sygnet |

Reguła piąta jest tym, co czyni emblemat emblematem. Ikona `terminal` i emblemat `srodowisko-codestudio` dzielą siatkę, obrys i barwę — różni je wyłącznie obecność wypełnionej kropki.

### 3.2. Metafora i ciężar optyczny

| Emblemat | Metafora | Ciężar optyczny | Dominanta kierunkowa |
|---|---|---|---|
| TalkIn | dymek rozmowy z liniami tekstu | średni — jedna duża forma zamknięta | pozioma (linie tekstu) |
| WorkSpace | cztery moduły projektu | wysoki — cztery bryły równej masy | siatka 2 × 2, bez dominanty |
| CodeStudio | ramka terminala z zachętą `>_` | wysoki — pełna ramka | pozioma (ramka 18 × 16) |
| MultitaskingAI | rdzeń koordynatora + trzej wykonawcy | niski — punkty i cienkie łączniki | promienista, oś pionowa |

Zestaw jest zrównoważony celowo: dwie formy zamknięte (TalkIn, CodeStudio), jedna siatkowa (WorkSpace), jedna otwarta promienista (MultitaskingAI). Cztery emblematy stojące obok siebie w strefie 1 nie zlewają się w plamę.

### 3.3. TalkIn — dymek rozmowy

| Element | Ścieżka / atrybuty | Rola |
|---|---|---|
| Dymek | `M20.5 5.5a2.5 2.5 0 0 0-2.5-2.5H6A2.5 2.5 0 0 0 3.5 5.5v8A2.5 2.5 0 0 0 6 16h1v4l4.4-4H18a2.5 2.5 0 0 0 2.5-2.5Z` | forma zamknięta, promień naroża 2,5 |
| Linia tekstu 1 | `M7.5 7.5h9` (dł. 9 j.) | wypowiedź zakończona |
| Linia tekstu 2 | `M7.5 10.8h5` (dł. 5 j.) | wypowiedź w toku — dlatego krótsza |
| **Kropka sygnału** | `circle 16,2 · 10,8 · r 1,5` | „tu biegnie odpowiedź" — na osi drugiej linii, 3,7 j. za jej końcem |

**Dlaczego kropka jest tam, gdzie jest.** Kropka i druga linia mają identyczne `cy` = 10,8. Odstęp między końcem linii (12,5) a lewą krawędzią kropki (14,7) wynosi 2,2 j. — czyta się jako pauza przed dalszym ciągiem.

### 3.4. WorkSpace — cztery moduły projektu

| Element | Współrzędne | Rola |
|---|---|---|
| Moduł ↖ | `rect 3,5 · 3,5 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↗ | `rect 13,3 · 3,5 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↙ | `rect 3,5 · 13,3 · 7,2 × 7,2 · rx 2` | moduł projektu |
| Moduł ↘ | `rect 13,3 · 13,3 · 7,2 × 7,2 · rx 2` | moduł **aktywny** |
| **Kropka sygnału** | `circle 16,9 · 16,9 · r 1,6` | środek modułu ↘: 13,3 + 7,2/2 = 16,9 — dokładnie, w obu osiach |

**Dlaczego kropka jest tam, gdzie jest.** Jest jedynym elementem łamiącym symetrię czterokrotną. Emblemat bez kropki byłby obojętnym znakiem „siatka"; kropka nadaje mu kierunek czytania ku prawej dolnej ćwiartce — tam, gdzie w interfejsie kończy się przepływ pracy.

### 3.5. CodeStudio — ramka terminala

| Element | Współrzędne | Rola |
|---|---|---|
| Ramka | `rect 3 · 4 · 18 × 16 · rx 2,5` | okno terminala |
| Grot zachęty | `M7 9.2l3.1 2.8L7 14.8` | znak `>` — jedyny bezpośredni cytat z sygnetu w zestawie emblematów |
| Linia wyniku | `M12.6 15.4h4` (dł. 4 j.) | odpowiedź powłoki |
| **Kropka sygnału** | `circle 18,4 · 9,9 · r 1,5` | proces w tle; prawa krawędź kropki 19,9 j. — 0,225 j. od wewnętrznej krawędzi ramki |

**Dlaczego kropka jest tam, gdzie jest.** Prawy górny róg okna terminala to w interfejsie miejsce wskaźnika stanu procesu. Kropka jest tam maksymalnie dosunięta — 0,225 j. luzu do ramki. To najciaśniejsze miejsce w całym systemie godeł.

### 3.6. MultitaskingAI — rdzeń i trzej wykonawcy

| Element | Współrzędne | Rola |
|---|---|---|
| **Rdzeń — kropka sygnału** | `circle 12 · 12 · r 2,6` | Coordinator; środek geometryczny siatki co do jednostki |
| Węzeł górny | `circle 12 · 4,6 · r 1,9` | Executor 1 |
| Węzeł lewy dolny | `circle 4,9 · 16,4 · r 1,9` | Executor 2 |
| Węzeł prawy dolny | `circle 19,1 · 16,4 · r 1,9` | Executor 3 / Validator |
| Łącznik górny | `M12 9.4V6.5` (dł. 2,9 j.) | delegacja w górę |
| Łącznik lewy | `M9.8 13.4l-3.3 1.9` (dł. 3,81 j.) | delegacja w dół-lewo |
| Łącznik prawy | `M14.2 13.4l3.3 1.9` (dł. 3,81 j.) | delegacja w dół-prawo |

**Dlaczego kropka jest tam, gdzie jest.** To jedyny emblemat, w którym kropka nie jest satelitą, lecz środkiem. W MultitaskingAI praca nie „dzieje się gdzieś w rogu okna" — praca **jest** koordynacją. Kropka ma też największy promień w zestawie (2,6 wobec 1,5–1,6), bo nosi sygnał **i** jest węzłem konstrukcji.

### 3.7. Pomiary plików realnych

| Emblemat | Plik | Bajtów | Elementów obrysowych | Pole zadruku (x) | Pole zadruku (y) |
|---|---|---:|:---:|---|---|
| TalkIn | `srodowisko-talkin.svg` | 450 | 3 | 2,625 – 21,375 | 2,125 – 20,875 |
| WorkSpace | `srodowisko-workspace.svg` | 505 | 4 | 2,625 – 21,375 | 2,625 – 21,375 |
| CodeStudio | `srodowisko-codestudio.svg` | 388 | 3 | 2,125 – 21,875 | 3,125 – 20,875 |
| MultitaskingAI | `srodowisko-multitaskingai.svg` | 467 | 6 | 2,125 – 21,975 | 1,825 – 19,175 |

Wielkości plików odczytane poleceniem `wc -c` z katalogu `zasoby/marka/srodowiska/`. Pole zadruku liczone z uwzględnieniem połowy obrysu (0,875 j.) po każdej stronie.

---

## 4. Prawo kropki — cztery położenia

### 4.1. Pomiar

| Emblemat | cx | cy | r | Ø | Ćwiartka | Odległość od środka siatki |
|---|---:|---:|---:|---:|---|---:|
| TalkIn | 16,2 | 10,8 | 1,5 | 3,0 | prawa górna (nieznacznie) | 4,36 j. |
| WorkSpace | 16,9 | 16,9 | 1,6 | 3,2 | prawa dolna | 6,93 j. |
| CodeStudio | 18,4 | 9,9 | 1,5 | 3,0 | prawa górna | 6,73 j. |
| MultitaskingAI | 12,0 | 12,0 | 2,6 | 5,2 | **środek** | 0,00 j. |

### 4.2. Reguła i jej jedyny wyjątek

```
   cx kropek satelickich:  16,2 · 16,9 · 18,4
   ────────────────────────────────────────────
   pas pionowy x ∈ [16,2 ; 18,4] — szerokość 2,2 j. na siatce 24
                                 = 9,2 % szerokości znaku

   MultitaskingAI: cx = 12,0   ← WYJĄTEK ŚWIADOMY
```

Trzy z czterech kropek leżą w wąskim pasie po prawej stronie znaku. Czwarta stoi na osi. To reguła z jednym umotywowanym wyjątkiem, a wyjątek niesie różnicę znaczeniową: **trzy środowiska mają pracę obok siebie, MultitaskingAI ma pracę w sobie.**

### 4.3. Zbieżność z żetonem systemowym

| Render | Ø kropki satelickiej | Zgodność z żetonem |
|---:|---:|---|
| 16 px | 2,00 px | poniżej żetonu — wariant 16 px powiększa kropkę |
| 24 px | 3,00 px | — |
| 32 px | 4,00 px | — |
| **48 px** | **6,00 px** | **= `--dn-wym-kropka` (6 px) co do piksela** |
| 64 px | 8,00 px | — |

Kropka emblematu wyrenderowana w 48 px jest **dokładnie tą samą kropką**, którą interfejs stawia przy karcie sesji i przy nadawcy piszącym. Rodzina domyka się na jednej liczbie.

---

## 5. Porównanie powłok — boczna nawigacja wobec panelu orkiestracji

### 5.1. Cztery pasy powłoki modułowej

```
┌──────────────────────────────────────────────────────────────┐
│ PASEK GÓRNY                              48 px · atramentowy │
│ sygnet · dom · szukaj · konfiguracja · Mobile · AOD · motyw   │
├──────────────────────────────────────────────────────────────┤
│ PAS KART SESJI                                        36 px  │
│ [● Studio ✕] [ Research ✕] [+]                                │
├──────────────┬───────────────────────────────────────────────┤
│ BOCZNA       │ OBSZAR ROBOCZY                                │
│ NAWIGACJA    │ ┌───────────────────────────────────────────┐ │
│ MODUŁÓW      │ │ okna operacyjne wybranego modułu          │ │
│ 224 px       │ ├───────────────────────────────────────────┤ │
│ wiersz 32 px │ │ Chat Window — pas komunikacji · 320 px    │ │
│              │ └───────────────────────────────────────────┘ │
└──────────────┴───────────────────────────────────────────────┘
```

### 5.2. Powłoka MultitaskingAI

```
┌──────────────────────────────────────────────────────────────┐
│ PASEK GÓRNY                     48 px · IDENTYCZNY z pow. mod.│
├──────────────────────────────────────────────────────────────┤
│ PAS KART SESJI    36 px    karta = PROCES ORKIESTRACJI        │
├──────────────┬───────────────────────────────────────────────┤
│ PANEL        │ OBSZAR ROBOCZY                                │
│ ORKIESTRACJI │  Executor 1 │ Executor 2 │ Coordinator        │
│ 224 px       │  Executor 3 / Validator                       │
│ 6 sekcji     │ ┌───────────────────────────────────────────┐ │
│              │ │ Chat Window — per rola                    │ │
└──────────────┴─┴───────────────────────────────────────────┴─┘
```

### 5.3. Zestawienie strukturalne

| Cecha | Powłoka modułowa (TalkIn / WorkSpace / CodeStudio) | Powłoka MultitaskingAI |
|---|---|---|
| Pasek górny | `.dn-pasek`, 48 px, atramentowy w obu motywach | **identyczny — bez zmian** |
| Co organizuje boczny panel | zadania modułowe — lista modułów | zespół modeli i agentów realizujących wspólny proces |
| Zawartość bocznego panelu | lista modułów wg macierzy dostępności | panel orkiestracji — 6 sekcji stałych |
| Po czym się nawiguje | po modułach | po rolach, kolejkach i orkiestracji |
| Co reprezentuje karta sesji | jeden moduł otwarty w danej chwili | jeden proces orkiestracji |
| Tytuł karty sesji | nazwa modułu otwartego w karcie | nazwa procesu/zespołu orkiestracji |
| Zawartość obszaru roboczego | okna operacyjne modułu + Chat Window | zawartość wybranej sekcji panelu (np. 4 karty ról) |
| Kontrolka „+" | otwiera pustą kartę, oczekuje wyboru modułu | otwiera pusty proces orkiestracji, oczekuje konfiguracji zespołu |
| Poziom zasięgu izolacji karty | „karta sesji" — poziom szósty z siedmiu | karta sesji + zagnieżdżony poziom „Rola" — najwęższy z siedmiu |

Źródło: `elementy-okien.md` Tabela 14, Tabela 15, Schemat 6.

### 5.4. Wymiary z żetonów użyte w szkieletach

| Element | Żeton | Wartość |
|---|---|---|
| Pasek górny | `--dn-wym-pasek` | 48 px |
| Pas kart sesji | `--dn-wym-pas-kart` | 36 px |
| Boczna nawigacja / panel orkiestracji | `--dn-wym-boczna` | 224 px |
| Pas komunikacji (Chat Window) | `--dn-wym-pas-komunikacji` | 320 px |
| Kontrolka / pozycja nawigacji | `--dn-wym-kontrolka` | 32 px |
| Wiersz tabeli | `--dn-wym-wiersz` | 36 px |
| Kropka sygnału | `--dn-wym-kropka` | 6 px |
| Wstęga aktywności | `--dn-wym-wstega` | 2 px |

Na planszy szkielety czterech powłok są wyświetlane w skali 40 % (transformacja `scale`), a wszystkie wymiary wewnętrzne pozostają zapisane w żetonach — proporcje są prawdziwe, nie przerysowane.

---

## 6. Macierz moduł × środowisko

`●` = okno w bocznej nawigacji · `○` = wyłącznie komponent własny strefy 2 (bez okna modułowego)

| # | Moduł | Okno wiodące | TalkIn | WorkSpace | CodeStudio | MultitaskingAI |
|---:|---|---|:---:|:---:|:---:|:---:|
| 1 | Studio | Studio Editor | ● | ● | | |
| 2 | Research | Research Workspace | ● | ● | | |
| 3 | Library | Library Explorer | ● | ● | | |
| 4 | Translate | Translation Panels | ● | | | |
| 5 | Browser | Browser Window | ● | ● | | |
| 6 | Assistant | Voice Console | ● | | | |
| 7 | Roundtable | Model Panels | ● | ● | ● | |
| 8 | Workspace | Project Dashboard | ● | ● | ● | |
| 9 | Automations | Workflow Builder | ○ | ○ | ○ | integracja 24/7 |
| 10 | Design | Design Board | | ● | ● | |
| 11 | Apps | Product Builder | | ● | ● | |
| 12 | Terminal | Terminal Tabs | | | ● | |
| 13 | Developer | Code Editor | | | ● | |
| 14 | Diagnostics | Diagnostics Center | | | ● | |
| 15 | Agents | Agent Builder | ● | ● | ● | rola/ekspert |
| | **Liczba pozycji w bocznej nawigacji** | | **9** | **9** | **8** | panel orkiestracji |

Źródło: KANON rozdz. 6 · `projekt-ui/README.md` rozdz. 5.

### 6.1. Podwójna rola trzech modułów

| Moduł | Jako moduł | Jako komponent własny strefy 2 | Etykieta działania na kaflu |
|---|---|---|---|
| **Assistant** | okno w bocznej nawigacji TalkIn | Profil asystenta | „Ustaw profil asystenta" |
| **Workspace** | okno w bocznej nawigacji trzech środowisk | Projekt | „Załóż projekt" |
| **Automations** | **jedyny moduł bez okna w bocznej nawigacji** | Automatyka | „Zbuduj automatykę" |

Czwarty kafel strefy 2 — **Agents** — prowadzi do Agent Buildera; Agents ma również okno w bocznej nawigacji wszystkich trzech środowisk modułowych.

---

## 7. Arytmetyka macierzy — co z niej wynika

### 7.1. Rozkład modułów po liczbie środowisk

| Obecność | Liczba modułów | Które |
|---:|:---:|---|
| w 3 środowiskach modułowych | 3 | Roundtable · Workspace · Agents |
| w 2 środowiskach modułowych | 6 | Studio · Research · Library · Browser (TalkIn + WorkSpace) · Design · Apps (WorkSpace + CodeStudio) |
| w 1 środowisku modułowym | 5 | Translate · Assistant (TalkIn) · Terminal · Developer · Diagnostics (CodeStudio) |
| w 0 środowiskach (tylko strefa 2) | 1 | Automations |
| **Razem** | **15** | |

### 7.2. Części wspólne par środowisk

| Para | Wspólnych modułów | Które |
|---|:---:|---|
| TalkIn ∩ WorkSpace | **7** | Studio · Workspace · Browser · Research · Library · Roundtable · Agents |
| WorkSpace ∩ CodeStudio | **5** | Workspace · Roundtable · Design · Apps · Agents |
| TalkIn ∩ CodeStudio | **3** | Workspace · Roundtable · Agents |

**Wniosek projektowy.** WorkSpace jest mostem: dzieli 7 modułów z TalkIn i 5 z CodeStudio, podczas gdy TalkIn i CodeStudio mają wspólny wyłącznie rdzeń trzech modułów. Kolejność kart w strefie 1 (TalkIn → WorkSpace → CodeStudio → MultitaskingAI) odwzorowuje tę bliskość — sąsiadujące karty są sobie najbliższe składem.

### 7.3. Pokrycie macierzy

| Środowisko | Modułów z okna w bocznej nawigacji | Z 15 modułów | Udział |
|---|:---:|:---:|---:|
| TalkIn | 9 | 15 | 60 % |
| WorkSpace | 9 | 15 | 60 % |
| CodeStudio | 8 | 15 | 53 % |
| MultitaskingAI | 0 | 15 | — (panel orkiestracji, 6 sekcji) |

Wyliczenie własne z macierzy KANON rozdz. 6; udziały zaokrąglone do pełnego procenta. MultitaskingAI nie występuje w macierzy jako kolumna udostępniająca moduły — myślnik nie oznacza braku funkcji, lecz inną strukturę.

---

## 8. Ścieżka wejścia do środowiska

### 8.1. Cztery kroki nawigacji

```
STRONA GŁÓWNA — Centrum dowodzenia
      │  KROK 1 — klik karty środowiska (strefa 1)
      ▼
POWŁOKA ŚRODOWISKA   TalkIn · WorkSpace · CodeStudio · MultitaskingAI
      │  KROK 2 — boczna nawigacja modułów albo panel orkiestracji
      ▼
MODUŁ  (Studio, Research, Developer …)   albo   ROLA / SEKCJA panelu
      │  KROK 3 — otwarcie okien modułu w obszarze roboczym
      ▼
OBSZAR ROBOCZY = Chat Window + okna modułu
      │  KROK 4 — praca odbywa się w karcie sesji bieżącej
      ▼
[ Karta sesji A · Studio ]  [ Karta sesji B · Research ]  [ + nowa karta ]
```

Źródło: `elementy-okien.md` Diagram 2.

### 8.2. Trzy stany przejścia odwzorowane na planszy

| Stan | Co się dzieje | Odwzorowanie wizualne |
|---|---|---|
| **1 · Strefa 1** | cztery karty środowisk w spoczynku | `.dn-karta-srodowiska` bez wstęgi, bez cienia |
| **2 · Karta wybrana** | wskazanie karty | wstęga górna 2 px w barwie sygnału, uniesienie −2 px, cień |
| **3 · Powłoka** | otwarcie przestrzeni roboczej | pasek górny + pas kart + boczny panel właściwy środowisku |

Czas przejścia: `--dn-czas-3` (0,22 s), krzywa `--dn-ease`. Trzy kroki odtwarzają się sekwencyjnie; przy `prefers-reduced-motion` żetony wyłączają ruch, a stany pozostają rozróżnialne wstęgą i etykietą.

### 8.3. Co zostaje po wyjściu

Klik ikony „dom" w pasku górnym zamyka powłokę bieżącego środowiska i pokazuje stronę główną; **karty sesji bieżącego środowiska pozostają aktywne w tle** i odtwarzają pełny stan przy powrocie. Dlatego kafle strefy 1 na planszy pokazują wskaźnik pracy w tle (tętniąca kropka) — to nie ozdoba, lecz jedyna informacja o tym, że coś biegnie poza widokiem.

---

## 9. Jeden język wizualny dla czterech środowisk

### 9.1. Co rozróżnia, a co nie rozróżnia środowisk

| Środek wyrazu | Rozróżnia? | Uzasadnienie |
|---|:---:|---|
| **Emblemat** | **tak** | cztery różne metafory, cztery różne ciężary optyczne, cztery różne położenia kropki |
| **Skład nawigacji** | **tak** | 9 / 9 / 8 pozycji modułów albo 6 sekcji panelu orkiestracji |
| **Nazwa (Space Grotesk)** | **tak** | krój nagłówkowy niesie wejście i tożsamość |
| Barwa | **nie** | jedna paleta monochromatyczna w obu motywach |
| Sygnał | **nie** | jeden błękit sygnałowy `--dn-sygnal-500` dla wszystkich czterech |
| Typografia interfejsu | **nie** | IBM Plex Sans wszędzie; IBM Plex Mono dla danych wszędzie |
| Gęstość i wymiary | **nie** | te same żetony: 48 / 36 / 224 / 32 px |
| Pasek górny | **nie** | atramentowy w obu motywach, identyczny w czterech środowiskach |

### 9.2. Uzasadnienie z KIERUNEK.md

Kontrakt kierunku brzmi: platforma operacyjna — **kokpit dowodzenia** — w języku **monochromatycznej precyzji**: odcienie bieli w motywie jasnym, odcienie czerni w motywie ciemnym, **jeden chłodny sygnał**.

Z tego kontraktu wynikają trzy rozstrzygnięcia wiążące dla środowisk:

1. **Neutralne są czysto neutralne.** Równe składowe RGB, żadnego podbarwienia. Gdyby środowiska miały własne barwy, monochrom przestałby istnieć już przy pierwszym kafelku strefy 1.
2. **Jeden sygnał, jedno znaczenie.** Błękit sygnałowy oznacza fokus, stan aktywny, postęp i pracę w tle. Cztery barwy środowisk odebrałyby sygnałowi jednoznaczność: operator nie wiedziałby, czy błysk oznacza „tu biegnie praca", czy „to środowisko CodeStudio".
3. **Stan nigdy samym kolorem.** Zasada dostępności (WCAG 2.1 AA) wymaga, aby stan niosła ikona albo etykieta. Barwa środowiska byłaby czwartym kanałem informacji bez odpowiednika tekstowego.

### 9.3. Zablokowany odruch

Katalog anty-domyślnych KANON rozdz. 1 wymienia wprost: **„Dziewięć kolorów tła dla dziewięciu nadawców → trzy klasy semantyczne + ikona + etykieta + plakietka roli."** Ta sama zasada rozciąga się na środowiska: **cztery barwy dla czterech środowisk → jeden monochrom + cztery emblematy + cztery składy nawigacji.**

Na planszy zablokowany odruch jest zademonstrowany działającym przyciskiem „Nadaj środowisku własną barwę" — zgodnie z zasadą **zero blokad** (ADL-017) przycisk jest klikalny i odpowiada komunikatem zamiast być wyszarzony.

### 9.4. Typografia w tożsamości środowiska

| Rola | Krój | Gdzie w tożsamości środowiska |
|---|---|---|
| Nagłówki, tytuły środowisk, logotyp | **Space Grotesk** 500–700 | nazwa środowiska na karcie strefy 1, nagłówek bocznej nawigacji |
| Interfejs i treść | **IBM Plex Sans** 400–700 | opis trybu pracy, pozycje bocznej nawigacji |
| Dane, identyfikatory, terminal | **IBM Plex Mono** 400–600 | motto na karcie, liczby pozycji, współrzędne emblematu, metadane sesji |

Różnica krojów niesie znaczenie: **Space Grotesk = wejście do środowiska i tożsamość marki; Plex Sans = praca; Plex Mono = maszyna.** Motto środowiska stoi krojem mono celowo — jest etykietą maszynową, nie hasłem reklamowym.

Test diakrytyków obowiązujący cały pakiet: `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`.

---

## 10. Kafle wejściowe strefy 1 — forma realna

### 10.1. Anatomia karty środowiska

| Element | Klasa | Forma |
|---|---|---|
| Karta (całość) | `.dn-karta-srodowiska` | promień `--dn-r-xl` (14 px), odstęp wewnętrzny `--dn-od-6` (24 px), obrys 1 px |
| Wstęga aktywności | `.dn-karta-srodowiska::before` | 2 px (`--dn-wym-wstega`), barwa `--dn-kropka`, `opacity: 0` w spoczynku |
| Emblemat | `.dn-karta-srodowiska-godlo` | 40 × 40 px, `currentColor` |
| Nazwa | `.dn-karta-srodowiska-tytul` | Space Grotesk, `--dn-fs-3xl` (30 px), waga 600 |
| Motto | `.dn-karta-srodowiska-motto` | IBM Plex Mono, `--dn-fs-sm` (12 px), `--dn-tekst-3` |
| Opis trybu pracy | `.dn-karta-srodowiska-opis` | IBM Plex Sans, `--dn-fs-base` (13 px), `--dn-tekst-2`, szerokość maks. 44 znaki |

### 10.2. Trzy stany karty

| Stan | Wstęga | Obrys | Cień | Przesunięcie |
|---|---|---|---|---|
| Spoczynek | ukryta (`opacity: 0`) | `--dn-obrys` | brak | 0 |
| Najechanie / `aria-current` | widoczna, 2 px | `--dn-obrys-mocny` | `--dn-cien-2` | `translateY(-2px)` |
| Fokus klawiaturą | jw. | jw. | jw. | jw. + pierścień `--dn-fokus` 2 px z odsunięciem 2 px |

Wszystkie przejścia trwają `--dn-czas-2` (0,16 s) z krzywą `--dn-ease`.

### 10.3. Stan aktywnej sesji w tle

Dokumentacja opisuje wskaźnik stanu karty sesji jako sygnalizację aktywnej pracy w tle: „brak aktywności — niewidoczny; aktywny — kropka widoczna, kolor wg statusu". Na kafelku strefy 1 ten sam mechanizm sygnalizuje, że w środowisku pozostały sesje otwarte po wyjściu na stronę główną.

Wartości liczbowe pokazane na planszy są **przykładowe** i oznaczone jako przykładowe. Zbudowano je ze świata produktu: nazwy modułów istniejących w danym środowisku, nazwy sekcji panelu orkiestracji, terminologia kart sesji.

| Środowisko | Przykładowy stan sesji w tle | Wskaźnik |
|---|---|---|
| TalkIn | 3 karty sesji · 2 czynne (Studio, Research) | tętno kropki |
| WorkSpace | 1 karta sesji · 1 czynna (Project Dashboard) | tętno kropki |
| CodeStudio | 1 karta sesji · bez pracy w tle | kropka neutralna |
| MultitaskingAI | 1 proces orkiestracji · Monitor procesu czynny | tętno kropki |

---

## 11. Decyzje projektowe planszy

| # | Decyzja | Uzasadnienie | Podstawa |
|---:|---|---|---|
| D1 | Emblematy wklejone **inline z realnych plików**, nie odrysowane | plansza ma być surowo wierna; każdy `d=` odpowiada bajtowi w `zasoby/marka/srodowiska/*.svg` | zasada nadrzędna zespołu |
| D2 | Szkielety powłok jako **żywy DOM w skali 40 %**, nie obrazki | proporcje pozostają prawdziwe (wymiary z żetonów), a szkielet reaguje na najechanie i klik | KANON 11.4 |
| D3 | Macierz z **sprzężonym podświetlaniem** w obie strony | macierz jest tabelą relacji; statyczna tabela nie pokazuje, że moduł „widzi" środowiska, a środowisko „widzi" moduły | KANON 6 |
| D4 | Znaki `●` / `○` renderowane **jako SVG**, nie jako znak typograficzny | ten sam obrys 1,75 i ta sama siatka 24 co reszta ikonografii; zakaz emoji | KANON 1, 4 |
| D5 | Siatka 24 na emblemacie jako **przełączalna nakładka**, nie stały element | nakładka jest narzędziem analizy, nie częścią znaku; znak musi dać się obejrzeć czysto | `03-marka/emblematy-i-ikony.md` 2.4 |
| D6 | Powiększenie kropki przez **przycięcie viewBox**, nie przez skalowanie rastra | powiększenie pozostaje wektorowe i pokazuje realne współrzędne | — |
| D7 | „Nadaj środowisku własną barwę" jako **działający przycisk z komunikatem** | demonstracja zasady zero blokad wprost na planszy; wyszarzony przycisk byłby złamaniem ADL-017 | KANON 8 |
| D8 | Motta środowisk krojem **mono**, nie nagłówkowym | motto jest etykietą maszynową; nazwa środowiska pozostaje jedynym miejscem dla Space Grotesk na karcie | `KIERUNEK.md` 3.2 |
| D9 | Dane sesji w tle **oznaczone jako przykładowe** w widocznej etykiecie, nie tylko w przypisie | zakaz zmyślonych metryk; czytelnik portfolio musi wiedzieć, co jest pomiarem, a co ilustracją | KANON 10.3 |
| D10 | Stara warstwa wizualna (granat + złoto) **nie jest cytowana** mimo obecności w dokumentacji źródłowej | `KIERUNEK.md` uznaje ją za zastępczą i nieobowiązującą; wszystkie wzmianki o „wstędze złotej" z `elementy-okien.md` czytane są jako „wstęga sygnału" | `KIERUNEK.md` metryka |
| D11 | Pokrycie macierzy podane jako **wyliczenie własne** z jawnym oznaczeniem | liczby 60 / 60 / 53 % nie występują w dokumentacji — są arytmetyką z macierzy i muszą być tak opisane | KANON 10.5 |
| D12 | Panel orkiestracji pokazany **w tej samej pozycji** co boczna nawigacja | dokumentacja stwierdza wprost: „inny panel, ta sama pozycja" | `elementy-okien.md` Schemat 6 |

---

## 12. Źródła

| Zakres | Plik | Rozdziały |
|---|---|---|
| Jedyne źródło prawdy pakietu | `WYNIK/KANON.md` | 1, 2, 3, 4, 5, 6, 8, 9, 10, 11 |
| Indeks projektu UI, macierz moduł × środowisko | `dok/projekt-ui/README.md` | 2, 3.5, 4, 5 |
| Strefa 1, powłoki środowisk, panel orkiestracji | `dok/projekt-ui/przeplyw/elementy-okien.md` | 2.3, 2.4, 2.6, 3.1–3.5, 4.1–4.6 |
| Środowisko MultitaskingAI | `dok/projekt-ui/srodowiska/multitaskingai.md` | 1.1–1.5 |
| Kierunek projektowy | `design/opracowania/design/01-kierunek/KIERUNEK.md` | 1, 2, 3.1, 3.2, 3.4, 3.5, 3.6, 4 |
| Konstrukcja emblematów | `WYNIK/03-marka/emblematy-i-ikony.md` | 2.3, 2.4, 3.1–3.5, 4.1–4.3 |
| Pliki emblematów | `WYNIK/zasoby/marka/srodowiska/*.svg` | 4 pliki, odczyt bezpośredni |
| Manifest ikon (motta środowisk) | `WYNIK/zasoby/ikony/manifest.json` | pozycje `srodowisko-*` |
| Klasy komponentów | `WYNIK/zasoby/css/komponenty.css` | `.dn-karta-srodowiska`, `.dn-boczna`, `.dn-karty-sesji`, `.dn-pasek`, `.dn-kafel`, `.dn-listwa` |
| Warstwa prototypu | `WYNIK/zasoby/prototyp.css`, `WYNIK/zasoby/prototyp.js` | `.pt-*`, atrybuty `data-*` |
| Prototypy powłok | `WYNIK/05-okna/srodowiska/{talkin,workspace,codestudio,multitaskingai}.html` | — |
| Prototyp strony głównej | `WYNIK/05-okna/przeplyw/03-centrum-dowodzenia.html` | strefa 1, 2, 3 |

---

## 13. Kontrola jakości

- [x] Wszystkie nazwy własne zgodne z dokumentacją — zero parafraz nazw modułów, okien, sekcji
- [x] Zero elementów bez pokrycia w dokumentacji — każde środowisko, moduł, sekcja i okno pochodzi z KANON albo z plików źródłowych
- [x] Liczby 9 / 9 / 8 / 6 potwierdzone w dwóch niezależnych miejscach (KANON rozdz. 6 i `elementy-okien.md` Tabela 11/16)
- [x] Współrzędne emblematów odczytane z realnych plików SVG, nie z opisu
- [x] Wielkości plików odczytane poleceniem systemowym
- [x] Wyliczenia własne (pokrycie, części wspólne par) oznaczone jako wyliczenia własne
- [x] Zero wartości szesnastkowych w warstwie stylu planszy — wyłącznie `var(--dn-*)`
- [x] Zero `#000000`
- [x] Zero atrybutu wyłączającego kontrolkę (zasada zero blokad, ADL-017)
- [x] Zero emoji jako ikon — wyłącznie inline SVG 24 × 24, obrys 1,75, `currentColor`
- [x] Zero Lorem ipsum, zmyślonych nazwisk, zmyślonych metryk, fikcyjnych firm
- [x] Dane przykładowe oznaczone jako przykładowe i zbudowane ze świata produktu
- [x] Oba motywy działają; pasek górny atramentowy w obu
- [x] Stan nigdy samym kolorem — zawsze towarzyszy ikona albo etykieta
- [x] Polskie diakrytyki poprawne w całym pliku

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
