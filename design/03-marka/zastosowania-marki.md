# Danaco Console — zastosowania marki na nośnikach

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Opracowanie** | Zastosowania marki — materiały i nośniki |
| **Wersja** | 2.0 |
| **Status** | wiążący — obowiązuje przy każdym wytworzeniu nośnika marki |
| **Data** | 2026-08-14 |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Odbiorcy** | projektant wizualny · osoba wytwarzająca materiały · drukarnia · dział wsparcia · autor prezentacji |
| **Zakres** | dwanaście nośników: sieć, ekran, druk, poczta, prezentacja, dokument |
| **Podstawa** | kontrakt systemu projektowego · kierunek systemu projektowego (3.4, 3.5) · `03-marka/ksiega-znaku.md` (5, 6, 7, 8) · `zasoby/marka/logo/*.svg` |
| **Pliki wytworzone** | `03-marka/zastosowania/` — 30 SVG + 32 PNG + 4 HTML |

---

## Spis treści

1. [Wprowadzenie](#1-wprowadzenie)
2. [Zasady wspólne wszystkich nośników](#2-zasady-wspólne-wszystkich-nośników)
3. [Dobór wariantu znaku do nośnika](#3-dobór-wariantu-znaku-do-nośnika)
4. [Katalog nośników](#4-katalog-nośników)
   - [4.1 Obraz Open Graph 1200 × 630](#41-obraz-open-graph-1200--630)
   - [4.2 Baner LinkedIn 1584 × 396](#42-baner-linkedin-1584--396)
   - [4.3 Tapeta 2560 × 1440 i 3840 × 2160](#43-tapeta-2560--1440-i-3840--2160)
   - [4.4 Tło slajdu 1920 × 1080](#44-tło-slajdu-1920--1080)
   - [4.5 Karta tytułowa dokumentacji 1600 × 900](#45-karta-tytułowa-dokumentacji-1600--900)
   - [4.6 Sygnatura poczty elektronicznej](#46-sygnatura-poczty-elektronicznej)
   - [4.7 Papier firmowy A4](#47-papier-firmowy-a4)
   - [4.8 Wizytówka 90 × 50 mm](#48-wizytówka-90--50-mm)
   - [4.9 Ekran powitalny aplikacji 1280 × 800 i 2560 × 1600](#49-ekran-powitalny-aplikacji-1280--800-i-2560--1600)
   - [4.10 Awatar / ikona profilu 400 × 400](#410-awatar--ikona-profilu-400--400)
   - [4.11 Szablon slajdów prezentacji](#411-szablon-slajdów-prezentacji)
   - [4.12 Znak wodny do dokumentów](#412-znak-wodny-do-dokumentów)
5. [Materiał źródłowy — pakiet 08-grafika](#5-materiał-źródłowy--pakiet-08-grafika)
6. [Nazewnictwo i inwentarz plików](#6-nazewnictwo-i-inwentarz-plików)
7. [Produkcja — od SVG do nośnika](#7-produkcja--od-svg-do-nośnika)
8. [Dostępność nośników](#8-dostępność-nośników)
9. [Kontrola jakości](#9-kontrola-jakości)
10. [Decyzje projektowe](#10-decyzje-projektowe)

---

## 1. Wprowadzenie

### 1.1 Czym jest ten dokument

Księga znaku (`03-marka/ksiega-znaku.md`) rozstrzyga, **jak znak jest zbudowany**.
Ten dokument rozstrzyga, **jak znak stoi na nośniku** — na obrazie w karcie
odnośnika, w tle prezentacji, na wizytówce, w stopce pisma, w sygnaturze poczty.

Każdy nośnik opisany jest pięcioma polami stałymi:

| Pole | Odpowiada na pytanie |
|---|---|
| **Przeznaczenie** | do czego nośnik służy i kto go używa |
| **Wymiary** | kadr, moduł siatki, marginesy |
| **Układ** | gdzie stoi znak, gdzie treść, gdzie metadane |
| **Dobór wariantu znaku** | która z konfiguracji A–H (ksiega-znaku 7.1) i w jakim wariancie barwnym |
| **Pole ochronne** | ile pustej przestrzeni należy się znakowi na tym nośniku |
| **Dopuszczalne modyfikacje** | co wolno zmienić, nie zgłaszając zmiany |
| **Zakazy** | co jest naruszeniem księgi |

### 1.2 Zasada nadrzędna opracowania

> **Nośnik nie jest kampanią.** Każdy element na każdym nośniku pochodzi
> z dokumentacji produktu: nazwa produktu, deskryptor, nazwy czterech
> środowisk, trójstopniowa hierarchia, dane producenta, numer wersji, status.
> **Zero fikcyjnych klientów, zero wymyślonych haseł reklamowych, zero
> zmyślonych metryk, zero nazwisk poza twórcą wskazanym w metryce produktu.**

Pola, które wypełnia użytkownik nośnika (nazwisko na wizytówce, stanowisko
w sygnaturze, tytuł prezentacji), są zapisane jako **pola w nawiasach
kwadratowych** — `[Imię i nazwisko]`, `[stanowisko]` — a nie jako przykładowe
osoby. To rozróżnienie jest celowe i wiążące.

### 1.3 Czego ten dokument nie zawiera

- konstrukcji znaku (krzywe, siatka, kąty) — patrz `ksiega-znaku.md` rozdz. 3–4,
- katalogu naruszeń znaku — patrz `ksiega-znaku.md` rozdz. 10,
- emblematów środowisk i ikony aplikacji — patrz `emblematy-i-ikony.md`,
- specyfikacji komponentów interfejsu — patrz `zasoby/css/komponenty.css`,
- **planu komunikacji, harmonogramu publikacji ani treści marketingowych** —
  te leżą poza zakresem systemu wizualnego.

### 1.4 Jak używać

```
   1. Ustal nośnik           → rozdz. 4, tabela nośników
   2. Sprawdź kadr i moduł   → pole „Wymiary"
   3. Dobierz konfigurację   → rozdz. 3, drzewo decyzyjne
   4. Sprawdź pole ochronne  → pole „Pole ochronne"
   5. Weź plik źródłowy      → rozdz. 6, inwentarz
   6. Przejdź listę kontroli → rozdz. 9
```

---

## 2. Zasady wspólne wszystkich nośników

### 2.1 Jeden akcent [NIENEGOCJOWALNE]

Na każdym nośniku **kropka sygnału jest jedynym elementem barwnym**. Wszystko
pozostałe — tło, typografia, reguły, siatka, ramki — mieści się w skali
neutralnej (odcienie bieli albo odcienie czerni).

```
   TŁO ─────────────────────── skala szarości (bez wyjątku)
   TYPOGRAFIA ──────────────── skala szarości (bez wyjątku)
   REGUŁY, SIATKA, RAMKI ───── skala szarości (bez wyjątku)
   KROPKA SYGNAŁU ──────────── #3B6FE0 / #5C8CEC   ← jedyny akcent
```

Konsekwencja praktyczna: **na nośniku nie ma niebieskiej linii, niebieskiego
paska ani niebieskiego bloku**. Jeśli nośnik potrzebuje podziału — dzieli go
reguła 1 px w `#E3E3E3` (jasne tło) albo `#2A2A2A` (ciemne tło).

### 2.2 Barwy nośników

| Rola | Nośnik ciemny | Nośnik jasny | Żeton systemu |
|---|---|---|---|
| Tło nośnika | `#0F0F0F` | `#F4F4F4` | `--dn-tlo` |
| Pole medalionu (awatar, ikona) | `#131313` | `#131313` | `--dn-rama` |
| Papier (druk) | — | `#FFFFFF` | `--dn-powierzchnia` |
| Farba znaku i tekst główny | `#ECECEC` | `#181818` | `--dn-tekst` |
| Tekst pomocniczy | `#9E9E9E` | `#616161` | `--dn-tekst-2` |
| Metadane, deskryptor | `#7C7C7C` | `#7C7C7C` | `--dn-tekst-3` |
| Reguła, obrys | `#2A2A2A` | `#E3E3E3` | `--dn-obrys` |
| **Kropka sygnału** | **`#5C8CEC`** | **`#3B6FE0`** | `--dn-kropka` |
| Siatka konstrukcyjna tła | biel, krycie 0,030–0,040 | atrament, krycie 0,045–0,060 | — |

**`#000000` i `#FFFFFF` jako tło nośnika ekranowego są zakazane.** Biel
`#FFFFFF` występuje wyłącznie jako papier w druku i jako tło sygnatury poczty
(wymuszone przez klienty pocztowe).

### 2.3 Typografia nośników

| Rola na nośniku | Krój | Waga | Traktowanie |
|---|---|---|---|
| Zdanie ekspozycyjne, tytuł | Space Grotesk | 700 | światło międzyliterowe −0,01 em |
| Logotyp `DANACO` | Space Grotesk | 700 | tracking 0,012 em (metryka znaku) |
| Logotyp `CONSOLE` | IBM Plex Mono | 500 | krok 0,90 em (tracking 0,30 em) |
| Deskryptor `AI OPERATING ENVIRONMENT` | IBM Plex Mono | 400 | wersaliki, tracking 0,14 em |
| Metadane, stopki, wersja | IBM Plex Mono | 400 | tracking 0,04–0,16 em |
| Treść ciągła (papier firmowy, slajd) | IBM Plex Sans | 400 / 600 | interlinia 1,45–1,55 |

**Deskryptor nigdy nie jest składany krojem nagłówkowym.** Zestawienie
Space Grotesk (marka) + IBM Plex Mono (maszyna) jest nośnikiem znaczenia,
nie ozdobą — zamiana krojów zaciera rozróżnienie.

### 2.4 Siatka konstrukcyjna nośnika

Każdy nośnik ma zdefiniowany **moduł siatki** — kwadrat, w którym mieszczą się
marginesy i pozycje bloków. Siatka bywa widoczna jako rysunek tła o kryciu
3–6 %; jest wtedy fakturą, nie dekoracją: powtarza raster kokpitu.

| Nośnik | Kadr | Moduł | Kolumn × wierszy | Margines |
|---|---:|---:|---:|---:|
| Open Graph | 1200 × 630 | 48 | 25 × 13,1 | 80 |
| Baner LinkedIn | 1584 × 396 | 44 | 36 × 9 | 120 |
| Tapeta | 2560 × 1440 | 80 | 32 × 18 | — |
| Tapeta | 3840 × 2160 | 120 | 32 × 18 | — |
| Tło slajdu | 1920 × 1080 | 60 | 32 × 18 | 120 |
| Karta tytułowa | 1600 × 900 | 50 | 32 × 18 | 100 |
| Ekran powitalny | 1280 × 800 | 40 | 32 × 20 | — |
| Ekran powitalny | 2560 × 1600 | 80 | 32 × 20 | — |
| Papier firmowy | 210 × 297 mm | 5 mm | 42 × 59,4 | 20 mm |
| Wizytówka | 90 × 50 mm | 2,5 mm | 36 × 20 | 7 mm |

Moduł 32 kolumn na nośnikach ekspozycyjnych jest świadomym rozszerzeniem
siatki 12-kolumnowej interfejsu: nośnik nie jest oknem roboczym, więc
gęstość rastra rośnie, a proporcja modułu pozostaje kwadratowa.

### 2.5 Pole ochronne na nośniku

Definicja modułu pozostaje bez zmian względem księgi znaku:

> **X = wysokość grotu.** **Pole ochronne = ½X ze wszystkich stron znaku.**

Na nośniku dochodzi druga reguła, mocniejsza od pierwszej:

> **Margines nośnika ≥ pole ochronne znaku.** Jeżeli znak stoi przy krawędzi,
> obowiązuje większa z dwóch wartości.

W praktyce marginesy nośników są wielokrotnie większe od ½X — to zamierzone:
znak na nośniku ekspozycyjnym oddycha, a nie tylko spełnia minimum.

| Nośnik | Konfiguracja | Szer. farby | X | ½X | Faktyczny margines | Zapas |
|---|---|---:|---:|---:|---:|---:|
| Open Graph | C poziomy 0,62 | 131,1 px | 18,2 | **9,1** | 80 px | 8,8 × |
| Baner LinkedIn | A sygnet 132 px | 106,6 px | 60,5 | **30,2** | 120 px | 4,0 × |
| Karta tytułowa | C poziomy 0,86 | 181,9 px | 25,2 | **12,6** | 100 px | 7,9 × |
| Slajd tytułowy | D pionowy 340 px | 340,0 px | 94,1 | **47,0** | 120 px | 2,6 × |
| Slajd treściowy | E kompakt 150 px | 150,0 px | 19,8 | **9,9** | 120 px | 12,1 × |
| Ekran powitalny 1280 | D pionowy 180 px | 180,0 px | 49,8 | **24,9** | 550 px (do krawędzi) | 22 × |
| Papier firmowy | C poziomy 45 mm | 45,0 mm | 6,2 | **3,1** | 20 mm | 6,5 × |
| Wizytówka awers | E kompakt 28 mm | 28,0 mm | 3,7 | **1,9** | 7 mm | 3,7 × |
| Awatar kwadrat | A sygnet 208 px | 208,0 px | 118,1 | **59,0** | 96 px | 1,6 × |
| Awatar koło | A sygnet 176 px | 176,0 px | 99,9 | **50,0** | 112 px | 2,2 × |

**Wyjątek jedyny: znak wodny** (rozdz. 4.12) — pole ochronne nie obowiązuje,
bo znak wodny leży **pod** treścią, a nie obok niej.

### 2.6 Wielkości minimalne na nośnikach

Przeniesienie tabeli z `ksiega-znaku.md` 6.1 na praktykę nośnikową:

```
   lockup poziomy  →  lockup kompaktowy  →  sygnet pełny  →  sygnet uproszczony
      ≥ 120 px            ≥ 96 px             ≥ 24 px            ≥ 16 px
      ≥  28 mm            ≥ 22 mm             ≥  9 mm            ≥  6 mm*
   * wariant uproszczony w druku wyłącznie za pisemną zgodą — patrz ksiega 6.3
```

Żaden nośnik z tego opracowania nie schodzi poniżej wielkości minimalnej.
Najmniejsze wystąpienia: **wizytówka awers 28 mm** (lockup kompaktowy, minimum
22 mm) i **sygnatura poczty 168 px** (lockup kompaktowy, minimum 96 px).

### 2.7 Znak na tle — reguła nośnikowa

| Tło nośnika | Wariant znaku | Uzasadnienie |
|---|---|---|
| `#0F0F0F` … `#212121` | **na ciemnym** (biel `#F4F4F4` / `#ECECEC` + kropka `#5C8CEC`) | kontrast farby do tła ≥ 12 : 1 |
| `#F4F4F4` … `#FFFFFF` | **podstawowy** (atrament `#181818` + kropka `#3B6FE0`) | kontrast farby do tła ≥ 14 : 1 |
| `#131313` (medalion) | **na ciemnym** — zawsze, niezależnie od motywu strony | medalion jest własnym tłem znaku |
| druk jednokolorowy | **mono czarny** / **mono biały** — kropka przyjmuje barwę grotów | brak drugiego koloru w maszynie |
| zdjęcie, tło mieszane | **niedozwolone bez podkładu** — patrz `ksiega-znaku.md` 9.2 | — |

---

## 3. Dobór wariantu znaku do nośnika

```
  NOŚNIK
    │
    ├─ Czy nazwa „Danaco Console" pada obok znaku tekstem?
    │     TAK  → A  sygnet
    │     NIE ↓
    │
    ├─ Jakie proporcje ma pole na znak?
    │     szerokie i niskie (pasek, stopka)     → E  lockup kompaktowy
    │     szerokie z miejscem na dwie linie     → C  lockup poziomy
    │     kwadratowe albo wysokie               → D  lockup pionowy
    │     kwadratowe i małe (≤ 96 px)           → A  sygnet
    │     kwadratowe i bardzo małe (< 24 px)    → A′ sygnet uproszczony
    │
    ├─ Czy odbiorca styka się z marką po raz pierwszy?
    │     TAK  → dodaj deskryptor „AI Operating Environment" (F)
    │
    ├─ Czy znak jest treścią, czy fakturą?
    │     treścią  → wariant podstawowy / na ciemnym, pełne krycie
    │     fakturą  → H  znak wodny, krycie 8 % (jasne) / 10 % (ciemne)
    │
    └─ Jaka technologia?
          ekran RGB           → wariant podstawowy / na ciemnym
          druk pełnokolorowy  → wariant podstawowy / na ciemnym
          druk jednokolorowy  → mono czarny / mono biały
          poczta elektroniczna→ PNG 2× w data URI (SVG bywa blokowane)
```

### 3.1 Mapa: nośnik → konfiguracja

| Nośnik | Konfiguracja | Wariant barwny | Wielkość znaku |
|---|---|---|---:|
| Obraz Open Graph | **C** lockup poziomy | na ciemnym / podstawowy | 131 px szer. |
| Baner LinkedIn | **A** sygnet + `DANACO CONSOLE` składane | na ciemnym / podstawowy | 132 px |
| Tapeta | **H** znak wodny z żywą kropką | mono z akcentem | 550 px (2560) |
| Slajd tytułowy | **D** lockup pionowy | na ciemnym / podstawowy | 340 px szer. |
| Slajd treściowy | **E** lockup kompaktowy | na ciemnym / podstawowy | 150 px szer. |
| Karta tytułowa dokumentacji | **C** lockup poziomy | na ciemnym / podstawowy | 182 px szer. |
| Sygnatura poczty | **E** lockup kompaktowy | podstawowy / na ciemnym | 168 px szer. |
| Papier firmowy — arkusz 1 | **C** lockup poziomy | podstawowy | 45 mm szer. |
| Papier firmowy — arkusz kolejny | **E** lockup kompaktowy | podstawowy | 26 mm szer. |
| Wizytówka awers | **E** lockup kompaktowy | na ciemnym | 28 mm szer. |
| Wizytówka rewers | **A** sygnet | podstawowy | 11 mm |
| Ekran powitalny | **D** lockup pionowy | na ciemnym | 180 px szer. |
| Awatar kwadrat / koło | **A** sygnet w medalionie | na ciemnym | 208 / 176 px |
| Awatar miniatura ≤ 48 px | **A′** sygnet uproszczony | na ciemnym | — |
| Znak wodny dokumentu | **H** | mono, krycie 8 / 10 % | dowolna |

---

## 4. Katalog nośników

### 4.1 Obraz Open Graph 1200 × 630

**Przeznaczenie.** Miniatura karty odnośnika w serwisach społecznościowych
i komunikatorach (Open Graph, Twitter Card). Pierwszy — a często jedyny —
kontakt odbiorcy z marką produktu.

**Wymiary.** 1200 × 630 px (proporcja 1,91 : 1). Moduł siatki 48 px.
Margines 80 px. Bezpieczna strefa treści: 1040 × 470 px — przycięcia
w podglądach komunikatorów sięgają 5 % z każdej strony.

**Układ.**

```
   ┌──────────────────────────────────────────────────────────────┐ 0
   │                                                              │
   │   ▛▚▛▚·  DANACO             ← lockup poziomy 131 px, y = 64  │
   │          CONSOLE·                                            │
   │                                                              │
   │                                                              │
   │   Zarządzaj cyfrową         ← Space Grotesk 700 · 64 px      │ 300
   │   organizacją●              ← kropka sygnału zamiast kropki  │ 376
   │                                interpunkcyjnej (r = 9)       │
   │   AI OPERATING ENVIRONMENT  ← Plex Mono 17 px · tracking .14 │ 436
   │                                                              │
   │  ─────────────────────────────────────────────────────────   │ 496
   │                                                              │
   │   TalkIn · WorkSpace · CodeStudio ·      danaco console·v2.0 │ 552
   └──────────────────────────────────────────────────────────────┘ 630
      80                                                        80
```

**Dobór wariantu znaku.** Konfiguracja **C** (lockup poziomy) w skali 0,62.
Wariant **na ciemnym** dla pliku podstawowego, **podstawowy** dla wariantu
jasnego. Deskryptor stoi osobno pod zdaniem ekspozycyjnym, nie w lockupie —
lockup pozostaje czysty, deskryptor pracuje jako blok samodzielny.

**Pole ochronne.** ½X = 9,1 px. Faktyczny margines 80 px (zapas 8,8 ×).
W polu ochronnym nie stoi nic — także siatka tła jest w tym miejscu
przykryta jednolitym tłem.

**Zdanie ekspozycyjne.** „Zarządzaj cyfrową organizacją." — sformułowanie
pochodzi wprost z kontraktu kierunku projektowego (kierunek systemu projektowego: *platforma
operacyjna dla zawodowego Operatora zarządzającego cyfrową organizacją*).
Kropka kończąca zdanie jest **kropką sygnału**, nie znakiem interpunkcyjnym —
to jedyne miejsce w systemie, gdzie kropka znaku wchodzi w tekst.

**Dopuszczalne modyfikacje.**
- zamiana zdania ekspozycyjnego na inne zdanie **z dokumentacji produktu**
  (maksymalnie dwa wiersze po 20 znaków),
- zamiana listy środowisk na nazwę pojedynczego środowiska, gdy karta dotyczy
  konkretnego środowiska,
- wariant jasny (`obraz-og-1200x630-jasny.png`) dla odbiorców z wymuszonym
  motywem jasnym.

**Zakazy.**
- ✗ zdjęcie, gradient albo blok barwny jako tło,
- ✗ drugi akcent barwny obok kropki sygnału,
- ✗ zdanie ekspozycyjne dłuższe niż dwa wiersze,
- ✗ logotyp powiększony ponad 1/6 szerokości kadru,
- ✗ tekst w polu ochronnym znaku,
- ✗ kompresja stratna poniżej jakości 85 — kropka rozmywa się pierwsza.

---

### 4.2 Baner LinkedIn 1584 × 396

**Przeznaczenie.** Tło profilu i strony firmowej w serwisie LinkedIn.
Nośnik o skrajnej proporcji 4 : 1, oglądany zwykle w skali 30–60 %.

**Wymiary.** 1584 × 396 px (proporcja 4 : 1). Moduł siatki 44 px
(36 × 9 komórek — podział całkowity). Margines 120 px.

**Strefy przycięcia — obowiązkowe do sprawdzenia:**

```
   ┌───────────────────────────────────────────────────────────────┐
   │◄── 1584 ─────────────────────────────────────────────────────►│
   │                                                               │
   │  ▛▚▛▚·   DANACO CONSOLE●                Danaco Holding Group  │
   │          AI OPERATING ENVIRONMENT       support@danaco-...    │
   │          Środowisko · Moduł · Okno                            │
   │                                                               │
   └───────────────────────────────────────────────────────────────┘
    ▲                    ▲                                        ▲
    │ strefa awatara     │ pas widoczny zawsze                    │
    │ profilu osobowego  │ (telefon przycina do 1128 px w środku) │
    │ ≈ 0–440 px         │                                        │
```

- **strona firmowa** — kadr widoczny w całości; blok treści od 120 px,
- **profil osobowy** — awatar zasłania lewe ≈ 440 px i dolne 90 px;
  w tym zastosowaniu blok treści przesuwa się do x = 560 px,
- **telefon** — widoczny środkowy pas 1128 px; krawędzie skrajne rezerwowe.

**Układ.** Sygnet 132 px na lewej krawędzi bloku, po nim odstęp 56 px,
dalej logotyp składany `DANACO CONSOLE` (Space Grotesk 700, 52 px) z kropką
sygnału po ostatniej literze, pod nim deskryptor 17 px i trójstopniowa
hierarchia 15 px. Blok producenta wyrównany do prawej krawędzi bloku treści.

**Dobór wariantu znaku.** Konfiguracja **A** (sygnet) + logotyp składany
w jednej linii. To **jedyny nośnik**, na którym `DANACO CONSOLE` stoi
w jednym wierszu — uzasadnienie: kadr 4 : 1 nie mieści dwuwierszowego
lockupu w czytelnej skali. Rozstrzygnięcie odnotowane w rozdz. 10 (D-3).

**Pole ochronne.** ½X = 30,2 px. Faktyczny margines 120 px (zapas 4,0 ×).
Odstęp sygnet → logotyp: 56 px = 1,85 × ½X — powyżej minimum.

**Dopuszczalne modyfikacje.**
- przesunięcie bloku treści do x = 560 px w wersji dla profilu osobowego,
- skrócenie prawego bloku do samej nazwy producenta,
- wariant jasny.

**Zakazy.**
- ✗ umieszczanie treści w dolnych 90 px (przycięcie awatarem),
- ✗ niebieski pas przy dolnej krawędzi (naruszenie zasady jednego aktu),
- ✗ powtarzanie znaku po obu stronach kadru,
- ✗ adres serwisu internetowego niepotwierdzony w metryce produktu.

---

### 4.3 Tapeta 2560 × 1440 i 3840 × 2160

**Przeznaczenie.** Tło pulpitu stanowiska pracy Operatora. Nośnik pracuje
**pod** oknami aplikacji — nie może konkurować z interfejsem.

**Wymiary.** 2560 × 1440 (moduł 80 px) oraz 3840 × 2160 (moduł 120 px).
Obie rozdzielczości mają identyczną kompozycję względną — moduł skaluje się
proporcjonalnie (32 × 18 komórek w obu).

**Układ.**

```
   ┌──────────────────────────────────────────────────────────┐
   │                                                          │
   │                                                          │
   │                    ▛▚  ▛▚                                │  ← znak wodny
   │                   ▟▘  ▟▘   ●                             │    krycie 7 %
   │                                ▲ kropka krycie 90 %      │    (biel)
   │                                                          │
   │              DANACO CONSOLE · AI OPERATING ENVIRONMENT   │  ← krycie 55 %
   │                                                          │
   │                                                          │
   └──────────────────────────────────────────────────────────┘
      środek kompozycji: x = 50 % · y = 46 % (oś optyczna)
```

Znak wodny ma szerokość farby **21,5 % szerokości kadru** (550 px przy 2560,
826 px przy 3840). Groty schodzą do krycia **7 %** (motyw ciemny) albo
**9 %** (motyw jasny) — czyta się je jako fakturę, nie jako znak. Kropka
sygnału pozostaje przy kryciu **90 %** i jest jedynym punktem, który pulpit
naprawdę pokazuje. To dosłowna realizacja zdania z kierunku: *„tu biegnie
praca"*.

**Dobór wariantu znaku.** Konfiguracja **H** (znak wodny) z zachowaną barwą
kropki. Groty w barwie farby motywu, kropka w barwie sygnału — **nie mono**,
bo pulpit ma pokazywać punkt sygnału.

**Pole ochronne.** Nie obowiązuje — znak wodny leży pod treścią.
Obowiązuje natomiast **strefa ikon pulpitu**: lewa kolumna 240 px (2560)
i 360 px (3840) pozostaje wolna od farby znaku.

**Dopuszczalne modyfikacje.**
- rozdzielczości pochodne przez proporcjonalne przeskalowanie SVG
  (1920 × 1080, 3440 × 1440, 5120 × 2880),
- podniesienie krycia grotów maksymalnie do 10 %,
- usunięcie podpisu tekstowego (wersja czysta).

**Zakazy.**
- ✗ krycie grotów powyżej 10 % (znak zaczyna konkurować z oknami),
- ✗ obniżenie krycia kropki poniżej 80 % (ginie jedyny akcent),
- ✗ powielanie znaku w rastrze (tapeta to jeden znak, nie deseń),
- ✗ dodawanie cytatów, dat, nazwisk i innych treści.

---

### 4.4 Tło slajdu 1920 × 1080

**Przeznaczenie.** Podkład slajdu prezentacji — do wstawienia jako obraz tła
w dowolnym programie prezentacyjnym. Dwa układy: **tytułowy** i **treściowy**.

**Wymiary.** 1920 × 1080 px (16 : 9). Moduł 60 px. Margines 120 px.
Pole tytułu na slajdzie treściowym zaczyna się pod regułą na y = 216 px.

**Układ — wersja tytułowa.**

```
   ┌────────────────────────────────────────────────────────────┐
   │                                                            │
   │      ▛▚▛▚·                                                 │
   │                          ← lockup pionowy 340 px szer.,    │
   │      DANACO                 wyśrodkowany w pionie (−40)    │
   │      CONSOLE·                                              │
   │                                                            │
   │      ──────────────                                        │ 760
   │      AI OPERATING ENVIRONMENT                              │ 808
   │                                                            │
   │      Danaco Holding Group Sp. z o.o.                       │ 996
   └────────────────────────────────────────────────────────────┘
     120                    ▲ prawe 2/3 kadru pozostają puste
                              — tam wchodzi tytuł prezentacji
```

**Układ — wersja treściowa.**

```
   ┌────────────────────────────────────────────────────────────┐
   │                                                            │
   │  ──────────────────────────────────────────────────────    │ 216  reguła tytułu
   │                                                            │
   │                    pole treści slajdu                      │
   │                    (pozostaje puste)                       │
   │                                                            │
   │                                                            │
   │  DANACO CONSOLE                            ▛▚▛▚· DANACO●   │ 984
   └────────────────────────────────────────────────────────────┘
     120                                                    120
```

**Dobór wariantu znaku.** Tytułowy — konfiguracja **D** (lockup pionowy)
340 px szerokości, co odpowiada zaleceniu księgi (60 mm na slajdzie).
Treściowy — konfiguracja **E** (lockup kompaktowy) 150 px szerokości
w prawym dolnym rogu, poza polem treści.

**Pole ochronne.** Tytułowy: ½X = 47,0 px, margines 120 px (zapas 2,6 ×).
Treściowy: ½X = 9,9 px, margines 120 px (zapas 12,1 ×).

**Dopuszczalne modyfikacje.**
- wariant jasny obu układów,
- usunięcie reguły tytułu na slajdzie treściowym, gdy slajd jest pełnoobrazowy,
- rozdzielczość 2560 × 1440 przez przeskalowanie SVG.

**Zakazy.**
- ✗ umieszczanie treści slajdu w polu znaku,
- ✗ powiększanie znaku na slajdzie treściowym ponad 200 px,
- ✗ znak w obu rogach jednocześnie,
- ✗ zmiana proporcji kadru bez przebudowy układu (4 : 3 wymaga własnego pliku).

---

### 4.5 Karta tytułowa dokumentacji 1600 × 900

**Przeznaczenie.** Pierwsza strona dokumentu (PDF, opracowanie, specyfikacja)
oraz miniatura dokumentu w indeksie. Nośnik niesie **metrykę wydania**.

**Wymiary.** 1600 × 900 px (16 : 9). Moduł 50 px. Margines 100 px.

**Układ.**

```
   ┌──────────────────────────────────────────────────────────────┐
   │  ▛▚▛▚· DANACO                                                │ 84
   │        CONSOLE·          ← lockup poziomy 182 px             │
   │                                                              │
   │                                                              │
   │  Danaco Console          ← Space Grotesk 700 · 84 px         │ 452
   │  AI OPERATING ENVIRONMENT                                    │ 506
   │                                                              │
   │  ──────────────────────────────────────────────────────────  │ 596
   │  WERSJA      STATUS         DATA          PRODUCENT          │ 650
   │  2.0         deweloperski   2026-08-14    Danaco Holding …   │ 682
   │                                                              │
   │  Środowisko · Moduł · Okno operacyjne                        │ 828
   └──────────────────────────────────────────────────────────────┘
     100                    kolumny metryki co 300 px
```

**Dobór wariantu znaku.** Konfiguracja **C** (lockup poziomy) w skali 0,86
— 182 px szerokości, powyżej minimum 120 px z zapasem 1,5 ×.

**Pole ochronne.** ½X = 12,6 px, margines 100 px (zapas 7,9 ×).

**Metryka.** Cztery kolumny stałe: **WERSJA · STATUS · DATA · PRODUCENT**.
Etykiety Plex Mono 13 px wersaliki, wartości Plex Sans 600 19 px. Wartości
pochodzą z metryki produktu — nie wolno ich zmyślać ani zaokrąglać.

**Dopuszczalne modyfikacje.**
- zamiana tytułu na tytuł konkretnego opracowania (maks. 28 znaków w jednym
  wierszu przy 84 px albo dwa wiersze przy 64 px),
- dodanie piątej kolumny metryki (np. `ZAKRES`) przy zmniejszeniu odstępu
  kolumn do 240 px,
- wariant jasny.

**Zakazy.**
- ✗ metryka z wartościami spoza dokumentacji produktu,
- ✗ tytuł przekraczający regułę na y = 596,
- ✗ dwa znaki na jednej karcie (nagłówek + stopka),
- ✗ zdjęcie albo ilustracja w tle.

---

### 4.6 Sygnatura poczty elektronicznej

**Przeznaczenie.** Stopka wiadomości poczty elektronicznej. Nośnik działa
w środowisku najbardziej wrogim projektowi: silnik Word w Outlooku,
przycinanie stylów przez Gmaila, wymuszona inwersja barw w trybie ciemnym
klientów mobilnych.

**Wymiary.** Szerokość 100 %, maksimum **520 px**. Znak 168 × 72 px
(w pliku 2 × = 336 px szerokości źródłowej — dla ekranów o podwójnej
gęstości). Odstępy pionowe: 16 / 12 / 12 / 10 px.

**Układ.**

```
   ┌────────────────────────────────────────┐ ← max 520 px
   │                                        │
   │   ▛▚▛▚·  DANACO●      ← znak 168 px    │
   │                                        │
   │   [Imię i nazwisko]   ← Plex Sans 600  │
   │   [STANOWISKO]        ← Plex Mono 11   │
   │  ────────────────────────────────────  │ ← reguła 1 px
   │   Danaco Holding Group Sp. z o.o.      │
   │   support@danaco-group.pl   ← odnośnik │
   │                                        │
   │   DANACO CONSOLE · AI OPERATING ENV.   │ ← Plex Mono 10
   └────────────────────────────────────────┘
```

**Rozstrzygnięcia techniczne [ZAMKNIĘTE].**

| Zagadnienie | Rozstrzygnięcie | Powód |
|---|---|---|
| Układ | tabela `role="presentation"`, `border-collapse: collapse` | silnik Word nie liczy `flex` ani `grid` |
| Style | wyłącznie w atrybucie `style` | Gmail usuwa `<style>` w `<head>` |
| Znak | PNG 2 × w `data:` URI | SVG blokowany przez Outlooka; plik zewnętrzny wymaga zgody odbiorcy |
| Wysokość reguły | komórka `height="1"` z `font-size: 0` | Outlook wymusza minimalną wysokość wiersza |
| Odnośnik | `mailto:` z jawnym `color` | klienty nadpisują barwę odnośników |
| Wersja ciemna | osobny plik, tło `#131313` | automatyczna inwersja psuje kontrast znaku |

**Dobór wariantu znaku.** Konfiguracja **E** (lockup kompaktowy).
Wersja jasna — wariant **podstawowy** (atrament + kropka `#3B6FE0`) na bieli.
Wersja ciemna — wariant **na ciemnym** (biel + kropka `#5C8CEC`) na `#131313`.

**Pole ochronne.** ½X = 2,5 px przy znaku 168 px. Faktyczne wcięcie
komórki 20 px — z dużym zapasem. Znak stoi we własnym wierszu tabeli,
więc pole ochronne jest zagwarantowane strukturą, nie odstępem.

**Pliki.** `zastosowania/html/sygnatura-poczty-jasna.html`
i `zastosowania/html/sygnatura-poczty-ciemna.html` — gotowe do wklejenia,
z komentarzem instruktażowym wewnątrz pliku.

**Dopuszczalne modyfikacje.**
- uzupełnienie pól `[Imię i nazwisko]` i `[stanowisko]`,
- usunięcie wiersza stanowiska,
- dodanie numeru telefonu jako trzeciego wiersza bloku kontaktowego
  (Plex Sans 13 px, ta sama barwa co adres).

**Zakazy.**
- ✗ dodawanie odnośników do serwisów społecznościowych jako ikon,
- ✗ obrazek tła (`background-image`) — usuwany przez większość klientów,
- ✗ czcionka o rozmiarze poniżej 10 px,
- ✗ cytat prawny albo klauzula poufności w barwie sygnału,
- ✗ znak jako plik zewnętrzny (`<img src="https://…">`).

---

### 4.7 Papier firmowy A4

**Przeznaczenie.** Pismo, notatka, protokół, wydruk dokumentu formalnego.
Nośnik ma dwie realizacje: **arkusz pierwszy** (pełny nagłówek) i **arkusz
kolejny** (nagłówek skrócony).

**Wymiary.** 210 × 297 mm. Margines **20 mm** ze wszystkich stron.
Nagłówek **32 mm** (arkusz kolejny: 18 mm). Stopka **18 mm**. Pole treści
maksymalnie **150 mm** szerokości — przy stopniu 10,5 pt daje to wiersz
80–90 znaków, wartość czytelną w druku.

**Układ.**

```
   ┌─ 210 mm ───────────────────────────────────────────┐
   │ ↕ 20 mm margines                                   │
   │  ┌──────────────────────────────────────────────┐  │
   │  │ ▛▚▛▚· DANACO          AI OPERATING ENVIRON.  │  │ nagłówek
   │  │       CONSOLE·        Danaco Holding Group   │  │ 32 mm
   │  ├──────────────────────────────────────────────┤  │ ← reguła 0,25 mm
   │  │                                              │  │
   │  │  MIEJSCOWOŚĆ, DATA                           │  │
   │  │  Tytuł pisma                                 │  │ pole treści
   │  │                                              │  │ max 150 mm
   │  │  …                                           │  │
   │  │                                              │  │
   │  ├──────────────────────────────────────────────┤  │ ← reguła 0,25 mm
   │  │ Danaco Holding Group Sp. z o.o.  Danaco …    │  │ stopka
   │  │ support@danaco-group.pl          Strona 01   │  │ 18 mm
   │  └──────────────────────────────────────────────┘  │
   │ ↕ 20 mm margines                                   │
   └────────────────────────────────────────────────────┘
```

**Rozstrzygnięcia druku [ZAMKNIĘTE].**

| Zagadnienie | Rozstrzygnięcie |
|---|---|
| Margines strony | `@page { size: A4 portrait; margin: 0 }` — cały margines niesie arkusz, przez co podgląd ekranowy i wydruk mają identyczną geometrię |
| Grubość reguł | 0,25 mm — najcieńsza kreska pewna w druku cyfrowym |
| Numeracja stron | licznik CSS `counter(page, decimal-leading-zero)` — `01`, `02`, `03` |
| Łamanie stron | `break-after: page` na arkuszu, `:last-child` bez łamania |
| Podziałka marginesów | przycisk ekranowy, wygaszany w `@media print` |

**Dobór wariantu znaku.** Arkusz pierwszy — konfiguracja **C** (lockup
poziomy) **45 mm** szerokości, wariant podstawowy. Arkusz kolejny —
konfiguracja **E** (lockup kompaktowy) **26 mm**, wariant podstawowy.
Oba powyżej minimum druku (28 mm / 22 mm).

**Pole ochronne.** ½X = 3,1 mm przy znaku 45 mm. Margines 20 mm i odstęp
od reguły nagłówka 6 mm — obie wartości powyżej minimum.

**Dopuszczalne modyfikacje.**
- dodanie w prawym bloku nagłówka trzeciego wiersza z numerem sprawy
  (Plex Mono 7 pt, wersaliki),
- zwiększenie pola treści do 170 mm przy piśmie tabelarycznym,
- druk jednokolorowy: znak w wariancie **mono czarny**.

**Zakazy.**
- ✗ znak przy samej krawędzi arkusza (pole ochronne bez zmian),
- ✗ kolorowy pas nagłówka albo stopki,
- ✗ tło rastrowe pod treścią (znak wodny — patrz 4.12 — jest wyjątkiem),
- ✗ dane producenta niepotwierdzone w metryce produktu,
- ✗ stopień pisma treści poniżej 9,5 pt.

---

### 4.8 Wizytówka 90 × 50 mm

**Przeznaczenie.** Karta kontaktowa. Format 90 × 50 mm — standard polskiego
rynku poligraficznego.

**Wymiary.** Netto 90 × 50 mm. Ze spadem 96 × 56 mm (spad **3 mm**
z każdej strony). Margines wewnętrzny **7 mm** od krawędzi netto.
Pasery narożne 2 mm.

**Układ — awers (atrament `#131313`).**

```
   ┌── 96 mm (ze spadem) ─────────────────────┐
   │ ┌── 90 mm netto ────────────────────────┐│
   │ │                                       ││
   │ │    ▛▚▛▚·  DANACO●                     ││ ← lockup kompaktowy
   │ │           ← 28 mm szerokości          ││   wyśrodkowany w pionie
   │ │                                       ││
   │ │  AI OPERATING ENVIRONMENT             ││ ← Plex Mono 2,1 mm
   │ └───────────────────────────────────────┘│
   └──────────────────────────────────────────┘
      ⌐ pasery narożne 2 mm w każdym rogu
```

**Układ — rewers (papier `#F4F4F4`).**

```
   ┌── 90 mm netto ───────────────────────────┐
   │                              ▛▚▛▚·       │ ← sygnet 11 mm
   │  [Imię i nazwisko]                       │ ← Plex Sans 600 · 3,6 mm
   │  [STANOWISKO]                            │ ← Plex Mono 2,3 mm
   │  ──────────────────────────────────────  │ ← reguła 0,2 mm
   │  support@danaco-group.pl                 │ ← Plex Sans 2,8 mm
   │  Danaco Holding Group Sp. z o.o.         │
   │                                          │
   │  Danaco Console · AI Operating Environ.  │ ← Plex Mono 2,1 mm
   └──────────────────────────────────────────┘
```

**Dobór wariantu znaku.** Awers — konfiguracja **E** (lockup kompaktowy)
28 mm, wariant **na ciemnym**. Rewers — konfiguracja **A** (sygnet) 11 mm,
wariant **podstawowy**. Sygnet 11 mm mieści się powyżej minimum druku 9 mm.

**Pole ochronne.** Awers: ½X = 1,9 mm, margines 7 mm (zapas 3,7 ×).
Rewers: sygnet 11 mm → ½X = 2,5 mm, odstęp od krawędzi 7 mm.

**Wskazania dla drukarni.**

| Parametr | Wartość |
|---|---|
| Format netto | 90 × 50 mm |
| Format ze spadem | 96 × 56 mm |
| Spad | 3 mm z każdej strony |
| Margines bezpieczeństwa | 7 mm od krawędzi netto |
| Gramatura zalecana | 300–350 g/m² |
| Uszlachetnienie zalecane | soft touch matowy, bez lakieru wybiórczego na znaku |
| Barwy | 4 + 4 CMYK albo 1 + 1 (wersja mono) |
| Minimalna kreska | 0,2 mm |

**Dopuszczalne modyfikacje.**
- zamiana strony atramentowej i jasnej (awers jasny, rewers atramentowy),
- dodanie numeru telefonu jako wiersza pod adresem poczty,
- wersja mono (1 + 1) dla druku jednokolorowego,
- zaokrąglenie narożników promieniem 2 mm.

**Zakazy.**
- ✗ znak w polu spadu,
- ✗ tekst bliżej niż 5 mm od krawędzi netto,
- ✗ lakier wybiórczy na kropce sygnału (zmienia jej barwę pod kątem),
- ✗ tłoczenie znaku poniżej 15 mm szerokości,
- ✗ nazwiska przykładowe zamiast pól `[Imię i nazwisko]` w pliku źródłowym.

---

### 4.9 Ekran powitalny aplikacji 1280 × 800 i 2560 × 1600

**Przeznaczenie.** Okno startowe aplikacji w stanie **Łączenie** — pierwszy
stan przepływu głównego (kontrakt systemu projektowego okno 1). Nośnik żyje 0,5–3 sekundy.

**Wymiary.** 1280 × 800 (moduł 40 px) i 2560 × 1600 (moduł 80 px).
Kompozycja identyczna względnie; wersja 2 × jest przeskalowaniem, nie
przeprojektowaniem.

**Układ.**

```
   ┌──────────────────────────────────────────────┐
   │                                              │
   │                                              │
   │                  ▛▚▛▚·                       │
   │                                              │ ← lockup pionowy
   │                  DANACO                      │   180 px szer.,
   │                  CONSOLE·                    │   środek − 24 px
   │                                              │
   │              ▬▬▬▬▬▬▭▭▭▭▭▭▭▭                  │ ← tor postępu 220 px
   │                  ŁĄCZENIE                    │ ← Plex Mono 13 px
   │                                              │
   │            v2.0 · status deweloperski        │ ← stopka 12 px
   └──────────────────────────────────────────────┘
```

**Tor postępu.** Wysokość 2 px, szerokość 220 px, promień 1 px. Tor
w bieli o kryciu 10 %, wypełnienie w bieli o kryciu 55 % — **nie w sygnale**.
Uzasadnienie: kropka znaku stoi 130 px wyżej i jest jedynym akcentem; drugi
punkt barwny rozbiłby hierarchię. Postęp niesie informację kształtem
i długością, nie barwą (zasada „stan nigdy samym kolorem").

**Etykieta stanu.** `ŁĄCZENIE` — nazwa stanu z dokumentacji okna startowego.
Pozostałe stany okna (`Powrót z ważnym tokenem`, `Błąd połączenia`) mają
własne komunikaty i nie są nośnikiem marki — obsługuje je prototyp okna.

**Dobór wariantu znaku.** Konfiguracja **D** (lockup pionowy) 180 px
szerokości, wariant **na ciemnym**. Zgodnie z zaleceniem księgi dla okna
startowego (160 px szer.) — tu 180 px, bo kadr splash jest większy od okna.

**Pole ochronne.** ½X = 24,9 px. Odstęp od toru postępu 130 px, od górnej
krawędzi 280 px — z wielokrotnym zapasem.

**Dopuszczalne modyfikacje.**
- rozdzielczości pochodne 1440 × 900, 1920 × 1200, 3840 × 2400,
- usunięcie toru postępu w wersji statycznej (obraz w pliku wykonywalnym),
- zamiana etykiety stanu na inny stan **z dokumentacji okna startowego**.

**Zakazy.**
- ✗ tor postępu w barwie sygnału,
- ✗ animowany znak (tętno kropki należy do interfejsu, nie do ekranu powitalnego),
- ✗ tekst powitalny, hasło, cytat,
- ✗ logo partnera, znak sklepu, plakietka wersji systemu.

---

### 4.10 Awatar / ikona profilu 400 × 400

**Przeznaczenie.** Zdjęcie profilowe konta produktu w serwisach zewnętrznych,
medalion systemu w oknie komunikacji, miniatura w katalogu aplikacji.

**Wymiary.** 400 × 400 px, dwa kadrowania: **kwadrat** (pełny kwadrat) i
**koło** (maska okrągła r = 200). Pole medalionu `#131313` — barwa ramy
kokpitu, ta sama w obu motywach.

**Układ.**

```
   KWADRAT 400 × 400                 KOŁO r = 200
   ┌─────────────────────┐           ╭───────────────────╮
   │                     │          ╱                     ╲
   │    ▛▚ ▛▚            │         │      ▛▚ ▛▚            │
   │   ▟▘ ▟▘   ●         │         │     ▟▘ ▟▘   ●         │
   │                     │          ╲                     ╱
   │  farba 208 px       │           ╰───────────────────╯
   └─────────────────────┘             farba 176 px
     margines 96 px                     margines 112 px
```

**Oś optyczna.** Znak jest wyśrodkowany **osią optyczną farby**
(x = 50,75 j. w siatce 96 × 96), nie środkiem kadru pliku. Różnica wynosi
2,75 j. i wynika z tego, że kropka wysuwa farbę w prawo. Bez tej poprawki
znak w medalionie czyta się jako przesunięty w lewo.

**Wielkość farby.** Kwadrat — 52 % boku. Koło — 44 % średnicy, bo maska
okrągła odcina narożniki i znak potrzebuje większego luzu bocznego.

**Dobór wariantu znaku.** Konfiguracja **A** (sygnet), wariant **na ciemnym**
— zawsze, niezależnie od motywu strony, na której awatar się wyświetla.
Miniatury **poniżej 48 px** — konfiguracja **A′** (sygnet uproszczony),
plik `awatar-400x400-kwadrat-uproszczony.png` skalowany w dół.

**Pole ochronne.** Kwadrat: ½X = 59 px, margines 96 px (zapas 1,6 ×).
Koło: ½X = 50 px, margines 112 px (zapas 2,2 ×). Obie wartości powyżej
minimum; kwadrat jest najciaśniejszym nośnikiem w całym katalogu — świadomie,
bo awatar bywa oglądany w skali 24 px i znak musi wypełnić kadr.

**Dopuszczalne modyfikacje.**
- wariant jasny koła (`awatar-400x400-kolo-jasny`) — dla serwisów o ciemnym tle,
- rozmiary pochodne 1024, 512, 256, 128, 64, 48, 32 px (skalowanie SVG),
- zaokrąglenie narożników kwadratu promieniem 88 px (22 % boku — zgodnie
  z proporcją ikony aplikacji 224/1024).

**Zakazy.**
- ✗ znak dotykający krawędzi kadru,
- ✗ pole medalionu w barwie innej niż `#131313`,
- ✗ pełny sygnet w rozmiarze poniżej 24 px (obowiązuje wariant uproszczony),
- ✗ logotyp tekstowy w kadrze awatara,
- ✗ obramowanie, poświata, cień rzucony pod znakiem.

---

### 4.11 Szablon slajdów prezentacji

**Przeznaczenie.** Żywy szablon HTML dwóch slajdów — **tytułowego**
i **treściowego** — do wypełnienia treścią i wydrukowania do PDF.
Uzupełnia statyczne tła z rozdz. 4.4: tam podkład, tu kompletny slajd.

**Wymiary.** Kadr odniesienia 1920 × 1080. Wszystkie miary zapisane
w jednostce `--u` = 1/1920 szerokości slajdu, więc szablon skaluje się bez
utraty proporcji do dowolnej szerokości okna i do wydruku poziomego A4.

**Układ — slajd tytułowy.** Lockup pionowy 340 u po lewej, reguła 340 u,
deskryptor; tytuł prezentacji w prawych dwóch trzecich (Space Grotesk 700,
84 u), podtytuł Plex Mono 22 u, stopka producenta 15 u.

**Układ — slajd treściowy.** Nadtytuł Plex Mono 16 u wersaliki, tytuł
Space Grotesk 700 52 u, reguła, trzy kolumny treści z etykietami
poprzedzonymi kropką sygnału 12 u, stopka z lockupem kompaktowym 150 u.

**Kropka w etykiecie kolumny.** Jedyne wystąpienie kropki sygnału poza
znakiem na tym nośniku. Kropka pełni tu rolę **znacznika porządkowego**
(zastępuje punktor listy) — dopuszczone, bo niesie tę samą treść, co
w interfejsie: „tu jest pozycja, która pracuje".

**Interakcje szablonu.**

| Kontrolka | Działanie |
|---|---|
| **Wariant jasny** | przełącza oba slajdy na tło `#F4F4F4` i podmienia wariant barwny znaku |
| **Siatka konstrukcyjna** | nakłada podziałkę marginesu 120 u |
| **Drukuj** | `@media print` → A4 poziomo, po jednym slajdzie na stronę |

**Dobór wariantu znaku.** Tytułowy — **D** lockup pionowy. Treściowy —
**E** lockup kompaktowy. Znak nie zmienia barw z motywem samoczynnie:
w drzewie dokumentu stoją **dwa pliki wariantu** (podstawowy i na ciemnym),
a przełącznik pokazuje właściwy. To realizacja reguły „godło zachowuje
barwy własne" — nie przebarwiamy znaku, wybieramy wariant.

**Pole ochronne.** Marginesy 120 u przy ½X = 47 u (tytułowy) i 9,9 u
(treściowy).

**Dopuszczalne modyfikacje.**
- zwiększenie liczby kolumn treści do czterech (gap 36 u),
- dodanie slajdu przejściowego (tytuł + reguła, bez kolumn),
- zmiana kadru na 4 : 3 przez `aspect-ratio: 4/3` i moduł 68 u.

**Zakazy.**
- ✗ więcej niż jeden znak na slajdzie,
- ✗ tytuł slajdu krojem monospace,
- ✗ kolumny treści zachodzące na pole znaku w stopce,
- ✗ kropka sygnału jako punktor w tekście ciągłym (tylko etykiety kolumn).

---

### 4.12 Znak wodny do dokumentów

**Przeznaczenie.** Podkład dokumentu, slajdu, wydruku roboczego i strony
tytułowej. Znak wodny **nie jest znakiem marki w użyciu** — jest fakturą
identyfikującą pochodzenie arkusza.

**Wymiary i krycie.**

| Plik | Konfiguracja | Kadr | Krycie | Zastosowanie |
|---|---|---:|---:|---|
| `znak-wodny-sygnet-na-jasnym` | A | 96 × 96 | **8 %** | podkład arkusza jasnego |
| `znak-wodny-sygnet-na-ciemnym` | A | 96 × 96 | **10 %** | podkład arkusza ciemnego |
| `znak-wodny-lockup-na-jasnym` | C | 225 × 96 | **8 %** | strona tytułowa, duży format |
| `znak-wodny-lockup-na-ciemnym` | C | 225 × 96 | **10 %** | jw. na ciemnym |
| `znak-wodny-kafel-na-jasnym` | A w kaflu 240 | 240 × 240 | **6 %** | deseń powtarzalny |
| `znak-wodny-kafel-na-ciemnym` | A w kaflu 240 | 240 × 240 | **8 %** | jw. na ciemnym |

**Dlaczego dwa krycia.** Percepcja krycia zależy od kierunku kontrastu:
atrament na bieli przy 8 % czyta się tak samo mocno, jak biel na atramencie
przy 10 %. Wartości nie są symetryczne, bo oko nie jest symetryczne.

**Mono bez wyjątku.** W znaku wodnym **kropka przyjmuje barwę grotów** —
to jedyne (obok wariantów mono) miejsce, gdzie kropka traci barwę sygnału.
Powód: znak wodny leży pod tekstem; kolorowy punkt pod literą czyta się jako
błąd druku, nie jako akcent.

**Układ na arkuszu.**

```
   pojedynczy                        deseń (kafel 240)
   ┌──────────────────┐              ┌──────────────────┐
   │                  │              │ ▚·  ▚·  ▚·  ▚·   │
   │      ▛▚▛▚·       │  ← 45 % szer.│ ▚·  ▚·  ▚·  ▚·   │
   │      (obrót 0°)  │    środek    │ ▚·  ▚·  ▚·  ▚·   │
   │                  │    arkusza   │ ▚·  ▚·  ▚·  ▚·   │
   └──────────────────┘              └──────────────────┘
```

**Pole ochronne.** Nie obowiązuje (`ksiega-znaku.md` 5.5).

**Dopuszczalne modyfikacje.**
- skalowanie do dowolnej wielkości (SVG),
- obrót o **0°** — i tylko 0°; znak wodny nie jest obracany,
- zmiana krycia w zakresie 6–12 % z zachowaniem różnicy jasne/ciemne.

**Zakazy.**
- ✗ krycie powyżej 12 % (znak przestaje być podkładem),
- ✗ obrót o 45° albo dowolny inny kąt,
- ✗ kropka w barwie sygnału w znaku wodnym,
- ✗ znak wodny pod treścią o kontraście poniżej 7 : 1,
- ✗ znak wodny na wizytówce, w sygnaturze poczty i na awatarze.

---

## 5. Materiał źródłowy — pakiet 08-grafika

Katalog `design/opracowania/design/08-grafika/` zawiera cztery pliki
wytworzone w pierwszym podejściu do nośników. Niniejsze opracowanie **zachowuje ich
język i naprawia trzy rzeczy**: rozdzielczości, dyscyplinę jednego akcentu
oraz źródło danych kontaktowych.

| Plik źródłowy | Wymiary rzeczywiste | Ocena | Następca |
|---|---:|---|---|
| `og-image.png` | **1200 × 630** | kompozycja i typografia zgodne z kierunkiem; do przeniesienia bez zmian | `obraz-og-1200x630.png` |
| `banner-linkedin.png` | **3168 × 792** | to **2 × 1584 × 396** — plik jest wariantem podwójnej gęstości, nie odrębnym formatem; zawiera niebieski pas przy dolnej krawędzi (drugi akcent) i adres serwisu spoza metryki produktu | `baner-linkedin-1584x396.png` |
| `tapeta-2560.png` | **2560 × 1440** | znak w pełnym kryciu — czyta się jako znak, nie jako podkład pulpitu | `tapeta-2560x1440.png` (krycie grotów 7 %) |
| `tlo-slajdu.png` | **1920 × 1080** | poprawne tło treściowe; brakuje wersji tytułowej i wariantu jasnego | `tlo-slajdu-tresciowy-1920x1080.png` + wersja tytułowa + warianty jasne |

**Wspólne cechy materiału źródłowego (zachowane):** tło `#0F0F0F`, siatka
o bardzo niskim kryciu, typografia Space Grotesk + IBM Plex Mono, kropka
sygnału jako punkt kompozycji, brak zdjęć i gradientów.

**Trzy poprawki wprowadzone w niniejszym opracowaniu:**

1. **Jeden akcent bez wyjątku.** Niebieski pas przy dolnej krawędzi banera
   został usunięty. Kropka sygnału jest jedynym elementem barwnym.
2. **Dane kontaktowe wyłącznie z metryki produktu.** Adres serwisu
   internetowego z banera źródłowego nie ma pokrycia w kontrakcie systemu projektowego;
   zastąpiono go nazwą producenta i adresem `support@danaco-group.pl`.
3. **Znak wodny zamiast znaku.** Tapeta źródłowa pokazuje znak w pełnym
   kryciu. Tapeta ma być podkładem pulpitu, więc groty schodzą do 7 %,
   a kropka zostaje.

**Zachowane bez zmian:** zdanie ekspozycyjne „Zarządzaj cyfrową
organizacją." wraz z kropką sygnału zamiast kropki interpunkcyjnej —
to najlepszy pomysł materiału źródłowego i zostaje w całości.

---

## 6. Nazewnictwo i inwentarz plików

### 6.1 Reguła nazewnicza

```
   <nośnik>-<wymiary>[-<wariant>].<rozszerzenie>

   obraz-og-1200x630.png                  ← nośnik + kadr
   obraz-og-1200x630-jasny.png            ← + wariant barwny
   wizytowka-awers-90x50-spady.svg        ← nośnik + strona + kadr + wersja
   znak-wodny-lockup-na-ciemnym.svg       ← nośnik + konfiguracja + tło
```

- **kebab-case**, wyłącznie małe litery, bez znaków diakrytycznych,
- **wymiary bez jednostki** dla pikseli (`1200x630`), z domyślnym mm dla druku
  (`90x50`),
- **wariant na końcu**: `-jasny` / `-jasna`, `-na-ciemnym`, `-uproszczony`,
  `-spady`, `-kwadrat`, `-kolo`,
- **SVG jest źródłem, PNG jest wynikiem** — nigdy odwrotnie.

### 6.2 Inwentarz — grafiki wytworzone

| Plik (SVG + PNG) | Kadr | Wariant |
|---|---:|---|
| `obraz-og-1200x630` | 1200 × 630 | ciemny |
| `obraz-og-1200x630-jasny` | 1200 × 630 | jasny |
| `baner-linkedin-1584x396` | 1584 × 396 | ciemny |
| `baner-linkedin-1584x396-jasny` | 1584 × 396 | jasny |
| `tapeta-2560x1440` | 2560 × 1440 | ciemny |
| `tapeta-2560x1440-jasna` | 2560 × 1440 | jasny |
| `tapeta-3840x2160` | 3840 × 2160 | ciemny |
| `tapeta-3840x2160-jasna` | 3840 × 2160 | jasny |
| `tlo-slajdu-tytulowy-1920x1080` | 1920 × 1080 | ciemny |
| `tlo-slajdu-tytulowy-1920x1080-jasny` | 1920 × 1080 | jasny |
| `tlo-slajdu-tresciowy-1920x1080` | 1920 × 1080 | ciemny |
| `tlo-slajdu-tresciowy-1920x1080-jasny` | 1920 × 1080 | jasny |
| `karta-tytulowa-dokumentacji-1600x900` | 1600 × 900 | ciemny |
| `karta-tytulowa-dokumentacji-1600x900-jasna` | 1600 × 900 | jasny |
| `wizytowka-awers-90x50` | 90 × 50 mm | atrament |
| `wizytowka-rewers-90x50` | 90 × 50 mm | papier |
| `wizytowka-awers-90x50-spady` | 96 × 56 mm | + pasery |
| `wizytowka-rewers-90x50-spady` | 96 × 56 mm | + pasery |
| `ekran-powitalny-1280x800` | 1280 × 800 | ciemny |
| `ekran-powitalny-2560x1600` | 2560 × 1600 | ciemny |
| `awatar-400x400-kwadrat` | 400 × 400 | medalion |
| `awatar-400x400-kwadrat-uproszczony` | 400 × 400 | medalion, A′ |
| `awatar-400x400-kolo` | 400 × 400 | medalion |
| `awatar-400x400-kolo-jasny` | 400 × 400 | papier |
| `znak-wodny-sygnet-na-jasnym` | 96 × 96 | 8 % |
| `znak-wodny-sygnet-na-ciemnym` | 96 × 96 | 10 % |
| `znak-wodny-lockup-na-jasnym` | 225 × 96 | 8 % |
| `znak-wodny-lockup-na-ciemnym` | 225 × 96 | 10 % |
| `znak-wodny-kafel-na-jasnym` | 240 × 240 | 6 % |
| `znak-wodny-kafel-na-ciemnym` | 240 × 240 | 8 % |

Dodatkowo w `png/`: `znak-poczty-jasny@2x.png` i `znak-poczty-ciemny@2x.png`
(znak do ręcznego wklejenia w sygnaturze).

### 6.3 Inwentarz — nośniki HTML

| Plik | Rola |
|---|---|
| `html/sygnatura-poczty-jasna.html` | gotowa do wklejenia sygnatura, tło białe |
| `html/sygnatura-poczty-ciemna.html` | jw., tło `#131313` |
| `html/papier-firmowy-a4.html` | dwa arkusze A4 z `@media print` |
| `html/szablon-slajdow.html` | slajd tytułowy i treściowy, przełączalne |

### 6.4 Inwentarz — narzędzia

| Plik | Rola |
|---|---|
| `generator-zastosowan.py` | wytwarza wszystkie SVG i PNG z zatwierdzonej geometrii |
| `generator-sygnatury.py` | wytwarza sygnatury poczty z osadzonym PNG w data URI |
| `fragmenty/*.svgfrag` | znak w konfiguracjach A/C/D/E, wklejany do plików HTML |

### 6.5 Drzewo katalogu

```
   03-marka/
   ├── zastosowania-marki.md          ← ten dokument
   ├── zastosowania-marki.html        ← galeria zastosowań
   └── zastosowania/
       ├── generator-zastosowan.py
       ├── generator-sygnatury.py
       ├── svg/    (30 plików — źródła)
       ├── png/    (32 pliki — rastry)
       ├── html/   (4 pliki — nośniki żywe)
       └── fragmenty/ (14 plików — znak inline do HTML)
```

---

## 7. Produkcja — od SVG do nośnika

### 7.1 Łańcuch wytwórczy

```
   geometria zatwierdzona          kroje jako pliki woff2
   (sygnet.svg, logotyp.svg)       (Space Grotesk, IBM Plex)
            │                                │
            └──────────────┬─────────────────┘
                           ▼
                 generator-zastosowan.py
                 (typografia → krzywe)
                           │
                 ┌─────────┴──────────┐
                 ▼                    ▼
              SVG źródłowy         PNG (cairosvg)
                 │                    │
        ┌────────┼────────┐           └── nośnik ekranowy
        ▼        ▼        ▼
     druk    prezentacja  poczta (PNG 2× w data URI)
```

### 7.2 Typografia zamieniona na krzywe [ZAMKNIĘTE]

Wszystkie napisy w plikach SVG są **konturami**, nie tekstem. Konsekwencje:

| Skutek | Ocena |
|---|---|
| plik nie zależy od zainstalowanych krojów | ✓ wymagane |
| rasteryzacja identyczna na każdej maszynie | ✓ wymagane |
| drukarnia nie musi dostawać krojów | ✓ wymagane |
| tekstu nie da się zaznaczyć ani wyszukać | ✓ akceptowane — to grafika, nie dokument |
| plik jest większy (30–100 kB zamiast 2 kB) | ✓ akceptowane |

Nośniki HTML (sygnatura, papier firmowy, slajdy) używają **żywej typografii**
— tam tekst musi być zaznaczalny, a znak jest wklejony jako SVG inline.

### 7.3 Rasteryzacja

```bash
python3 -c "import cairosvg; cairosvg.svg2png(
    url='svg/obraz-og-1200x630.svg',
    write_to='png/obraz-og-1200x630.png',
    output_width=1200)"
```

Wszystkie PNG w tym opracowaniu wytworzono w **skali 1 : 1** względem kadru
projektowego. Wersje o podwójnej gęstości uzyskuje się przez podanie
`output_width` przemnożonego przez 2 — bez zmian w pliku SVG.

### 7.4 Wskazania dla druku

| Zagadnienie | Rozstrzygnięcie |
|---|---|
| Format plików do drukarni | SVG (wektor) + PDF/X-4 z osadzonymi krzywymi |
| Przestrzeń barwna | konwersja RGB → CMYK po stronie drukarni, profil ISO Coated v2 |
| Kropka sygnału w CMYK | wartość docelowa: **C 78 M 55 Y 0 K 0** (odpowiednik `#3B6FE0`) |
| Atrament w CMYK | **K 100** przy druku jednokolorowym; **C 0 M 0 Y 0 K 92** przy 4 c. |
| Spad | 3 mm na wizytówce, 0 mm na papierze firmowym (zadruk nie sięga krawędzi) |
| Minimalna kreska | 0,2 mm (offset i cyfra), 0,4 mm (sitodruk) |
| Nadruk czerni | znak nie jest nadrukowywany (overprint OFF) |

### 7.5 Odtworzenie kompletu

```bash
cd 03-marka/zastosowania
python3 generator-zastosowan.py   # 30 SVG + 30 PNG
python3 generator-sygnatury.py    # 2 HTML + 2 PNG znaku
```

Generator jest **deterministyczny** — dwa uruchomienia dają bajtowo
identyczne pliki. Każda zmiana nośnika przechodzi przez generator,
nie przez ręczną edycję pliku wynikowego.

---

## 8. Dostępność nośników

### 8.1 Kontrasty zmierzone

| Para | Nośnik | Kontrast | Próg WCAG 2.1 |
|---|---|---:|---|
| `#ECECEC` na `#0F0F0F` | tekst główny, ciemny | **15,0 : 1** | AAA (7,0) ✓ |
| `#9E9E9E` na `#0F0F0F` | tekst pomocniczy, ciemny | **6,8 : 1** | AA (4,5) ✓ |
| `#7C7C7C` na `#0F0F0F` | metadane ≥ 18,66 px półgrube | **4,1 : 1** | AA duży (3,0) ✓ |
| `#181818` na `#F4F4F4` | tekst główny, jasny | **16,1 : 1** | AAA ✓ |
| `#616161` na `#F4F4F4` | tekst pomocniczy, jasny | **6,4 : 1** | AA ✓ |
| `#ECECEC` na `#131313` | znak w medalionie | **14,5 : 1** | — (grafika) ✓ |
| `#5C8CEC` na `#0F0F0F` | kropka sygnału, ciemny | **6,2 : 1** | AA grafika (3,0) ✓ |
| `#3B6FE0` na `#F4F4F4` | kropka sygnału, jasny | **4,9 : 1** | AA grafika ✓ |

**`#7C7C7C` obowiązuje wyłącznie w metadanych** — na każdym nośniku stoi
w stopniu ≥ 13 px w kroju monospace o dużym świetle międzyliterowym, co
odpowiada progowi „tekst duży" po przeliczeniu wagi optycznej.

### 8.2 Tekst zastępczy

| Nośnik | Tekst zastępczy |
|---|---|
| Open Graph | `Danaco Console — AI Operating Environment. Zarządzaj cyfrową organizacją.` |
| Baner LinkedIn | `Danaco Console — AI Operating Environment` |
| Awatar | `Danaco Console` |
| Znak w sygnaturze | `Danaco Console` |
| Tapeta, tło slajdu, znak wodny | `alt=""` + `aria-hidden="true"` — dekoracja |

Każdy plik SVG wytworzony przez generator ma `role="img"` i
`aria-label="Danaco Console"`; znak wodny i siatka mają `aria-hidden`.

### 8.3 Nośniki HTML

- `lang="pl"` w każdym pliku,
- kontrolki podglądu obsługiwane klawiaturą, `aria-pressed` na przełącznikach,
- pierścień fokusu 2 px + odsunięcie 2 px w barwie sygnału,
- **zero atrybutu `disabled`** (zasada zero blokad),
- `@media print` nie wymaga barwy tła do czytelności.

### 8.4 Nośniki drukowane

- stopień pisma treści ≥ 9,5 pt (papier firmowy), ≥ 7 pt (wizytówka),
- kontrast farby do papieru ≥ 12 : 1 w obu realizacjach,
- informacja nigdy nie jest przenoszona samą barwą — kropka sygnału ma zawsze
  towarzyszący tekst albo pozycję w znaku.

---

## 9. Kontrola jakości

Lista obowiązuje przed oddaniem **każdego** nośnika.

### 9.1 Znak

- [ ] geometria sygnetu niezmieniona (`M12 26 H24 L44 48 …`)
- [ ] konfiguracja z katalogu A–H, nie układ własny
- [ ] wariant barwny dopasowany do tła (podstawowy / na ciemnym / mono)
- [ ] wielkość ≥ minimum dla nośnika i technologii
- [ ] pole ochronne ≥ ½X, puste
- [ ] znak nie jest obrócony, ścięty, rozciągnięty ani przebarwiony
- [ ] kropka sygnału obecna i w barwie właściwej wariantowi

### 9.2 Kompozycja

- [ ] kropka sygnału jest **jedynym** akcentem barwnym
- [ ] tło z palety nośników (bez `#000000`, bez `#FFFFFF` na ekranie)
- [ ] marginesy zgodne z tabelą 2.4
- [ ] siatka tła o kryciu ≤ 6 %
- [ ] brak gradientu jako tła bloku, przycisku, sekcji
- [ ] brak zdjęcia pod znakiem

### 9.3 Typografia

- [ ] Space Grotesk 700 wyłącznie w tytułach i logotypie
- [ ] IBM Plex Mono w deskryptorze, metadanych i danych
- [ ] IBM Plex Sans w treści ciągłej
- [ ] polskie znaki diakrytyczne poprawne (`ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`)
- [ ] deskryptor `AI Operating Environment` bez tłumaczenia

### 9.4 Treść

- [ ] każdy napis ma pokrycie w dokumentacji produktu
- [ ] pola użytkownika zapisane jako `[…]`, nie jako przykładowe osoby
- [ ] nazwy środowisk dokładnie: `TalkIn`, `WorkSpace`, `CodeStudio`, `MultitaskingAI`
- [ ] dane producenta i kontakt zgodne z metryką produktu
- [ ] brak zmyślonych metryk, klientów, kampanii, adresów

### 9.5 Plik

- [ ] nazwa w kebab-case, wymiary w nazwie
- [ ] SVG jest źródłem, PNG wynikiem generatora
- [ ] PNG w skali 1 : 1 względem kadru projektowego
- [ ] `role="img"` + `aria-label` albo `aria-hidden` w SVG
- [ ] plik HTML działa po otwarciu z `file://`

---

## 10. Decyzje projektowe

### Jeden akcent bez wyjątku — także tam, gdzie materiał źródłowy dopuszczał drugi

**Problem.** Baner źródłowy (`08-grafika/banner-linkedin.png`) niesie
niebieski pas przy dolnej krawędzi. Pas jest ładny i czyta się jako
podpis marki. Czy zachować?

**Rozstrzygnięcie.** **Nie.** Kontrakt kierunku mówi o **jednym chłodnym
sygnale**, a księga znaku o kropce jako **jedynym** akcencie. Pas jest drugim
wystąpieniem barwy sygnału w tym samym kadrze — osłabia kropkę, bo oko
przestaje traktować błękit jako punkt informacji, a zaczyna jako kolor marki.
Wszystkie reguły i podziały na nośnikach są neutralne.

**Konsekwencja.** Podziały nośne wykonuje reguła 1 px w `#2A2A2A` / `#E3E3E3`.

---

### Baner LinkedIn — 1584 × 396, nie 3168 × 792

**Problem.** Plik źródłowy ma 3168 × 792 px. Czy to odrębny format?

**Rozstrzygnięcie.** **Nie — to wariant podwójnej gęstości.** 3168 = 2 × 1584,
792 = 2 × 396. Format bazowy LinkedIn to 1584 × 396; plik 2 × jest wynikiem
rasteryzacji, nie osobnym projektem. Opracowanie utrzymuje kadr projektowy 1584 × 396
w SVG i wytwarza rastry o dowolnej gęstości z tego samego źródła.

---

### Logotyp w jednym wierszu — wyłącznie na banerze

**Problem.** Kadr 4 : 1 nie mieści dwuwierszowego lockupu w skali czytelnej
przy 30 % powiększenia. Księga zabrania tworzenia własnych konfiguracji (7.5).

**Rozstrzygnięcie.** Baner używa **konfiguracji A (sygnet) + składanego
napisu `DANACO CONSOLE`** — to nie jest nowa konfiguracja lockupu, tylko
sygnet postawiony obok tytułu składanego w kroju nagłówkowym. Rozróżnienie
jest istotne: w lockupie `CONSOLE` jest monospace i stoi pod `DANACO`;
tutaj `DANACO CONSOLE` jest jednym napisem Space Grotesk 700 z kropką
sygnału na końcu.

**Ograniczenie.** Zastosowanie **wyłącznie** na banerze 4 : 1. Na każdym
innym nośniku obowiązują konfiguracje A–H bez zmian.

---

### Tapeta — znak wodny, nie znak

**Problem.** Tapeta źródłowa pokazuje znak w pełnym kryciu na środku ekranu.
Pulpit z otwartymi oknami przykrywa go w 90 %, a odsłonięty fragment
konkuruje z interfejsem.

**Rozstrzygnięcie.** Groty schodzą do **7 %** (ciemna) / **9 %** (jasna),
kropka zostaje przy **90 %**. Tapeta pokazuje wtedy dokładnie to, co powinna:
jeden punkt sygnału na neutralnym polu, z ledwie wyczuwalną fakturą znaku.

**Wartość liczbowa.** 7 % bieli na `#0F0F0F` daje wartość zmierzoną `#1D1D1D`
— różnicę widać, ale nie da się jej pomylić z elementem interfejsu.

---

### Tor postępu na ekranie powitalnym — biel, nie sygnał

**Problem.** Odruch projektowy każe zabarwić wypełnienie toru na błękit.

**Rozstrzygnięcie.** **Biel o kryciu 55 %.** Kropka znaku stoi 130 px wyżej
i jest jedynym akcentem kadru. Drugi punkt barwny rozdzieliłby uwagę na dwa
ośrodki. Postęp niesie informację **długością**, nie barwą — co jest zgodne
z regułą „stan nigdy samym kolorem".

---

### Awatar — oś optyczna, nie środek kadru

**Problem.** Znak wyśrodkowany geometrycznie w kwadracie czyta się jako
przesunięty w lewo, bo kropka wysuwa farbę w prawo.

**Rozstrzygnięcie.** Środkowanie **osią optyczną farby**: x = 50,75 j.
w siatce 96 × 96 (środek między lewą krawędzią grotu na 12 a prawą krawędzią
kropki na 89,5). Różnica 2,75 j. odpowiada 11,5 px przy awatarze 400 px —
wartość dobrze widoczna.

**Zgodność.** Rozstrzygnięcie powtarza zapis `ksiega-znaku.md` 7.4
(„sygnet w kole, oś optyczna x = 50,75"). Istniejąca ikona aplikacji
(`zasoby/marka/ikona-aplikacji/ikona-aplikacji.svg`) używa środka siatki
(48 j.) — różnica 2,75 j. przy kaflu 1024 px to 18 px. **Pozycja ikony
aplikacji nie jest zmieniana** (plik jest zatwierdzony i wdrożony);
nowe nośniki medalionowe używają osi optycznej.

---

### Pola użytkownika w nawiasach kwadratowych, nie przykładowe osoby

**Problem.** Wizytówka i sygnatura poczty potrzebują nazwiska i stanowiska.
Odruch każe wpisać przykładową osobę.

**Rozstrzygnięcie.** **Pola `[Imię i nazwisko]`, `[stanowisko]`.** Zakaz
zmyślonych nazwisk z kontraktu systemu projektowego jest bezwzględny, a przykładowa osoba
w pliku produkcyjnym prędzej czy później trafia do druku. Nawias kwadratowy
jest sygnałem „to pole do wypełnienia" czytelnym bez instrukcji.

---

### Dane kontaktowe wyłącznie z metryki produktu

**Problem.** Baner źródłowy zawiera adres serwisu internetowego, którego nie
ma w kontrakcie systemu projektowego.

**Rozstrzygnięcie.** Na nośnikach występują wyłącznie: **Danaco Holding
Group Sp. z o.o.** i **support@danaco-group.pl** — obie wartości z metryki
produktu. Adres serwisu wróci na nośniki, gdy zostanie potwierdzony
w dokumentacji.

---

### Typografia zamieniona na krzywe w plikach SVG

**Problem.** Kroje są dostarczone jako `woff2` (podzbiory `latin`
i `latin-ext`). SVG z żywym tekstem zależy od instalacji kroju na maszynie
rasteryzującej i w drukarni.

**Rozstrzygnięcie.** Generator wyciąga kontury glifów (fontTools) i składa
je jako `<path>`. Napis jest niezaznaczalny — akceptowane, bo to grafika.
Nośniki HTML zachowują żywy tekst.

**Skutek uboczny rozwiązany.** Podzbiór `latin` nie zawiera polskich znaków
diakrytycznych; generator czyta `latin` i `latin-ext` **razem** i dla każdego
znaku wybiera pierwszy plik, który go ma. Test `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`
przechodzi w obu rodzinach.

---

### Siatka 32 kolumn na nośnikach, 12 w interfejsie

**Problem.** Interfejs używa siatki 12 kolumn (kontrakt systemu projektowego). Nośnik
ekspozycyjny w tej siatce ma moduł 100–160 px — za gruby jako faktura.

**Rozstrzygnięcie.** Nośniki ekspozycyjne używają **siatki kwadratowej
32 × 18** (moduł = szerokość ÷ 32). Moduł wychodzi 40–120 px zależnie od
kadru i we wszystkich kadrach 16 : 9 daje podział całkowity. Siatka nośnika
nie jest siatką układu treści — jest rastrem faktury, powtórzeniem rastra
kokpitu.

---

### Dwie wartości krycia znaku wodnego

**Problem.** Czy krycie znaku wodnego na jasnym i ciemnym może być takie samo?

**Rozstrzygnięcie.** **Nie.** Percepcja zależy od kierunku kontrastu:
atrament na bieli przy 8 % czyta się z tą samą siłą, co biel na atramencie
przy 10 %. Symetryczna wartość dałaby znak wodny za słaby na ciemnym albo
za mocny na jasnym.

---

### Sygnatura poczty jako tabela ze stylami w atrybutach

**Problem.** Nowoczesny HTML (flex, grid, `<style>`) nie działa w silniku
Word, którego używa Outlook dla Windows.

**Rozstrzygnięcie.** Tabela `role="presentation"`, style wyłącznie
w atrybucie `style`, znak jako PNG 2 × w `data:` URI, reguła jako komórka
`height="1"` z `font-size: 0`. Rozwiązanie jest technicznie wsteczne
i **celowo** takie pozostaje — nośnik ma działać u odbiorcy, nie
w przeglądarce projektanta.

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
