# Danaco Console — Konstrukcja znaku

## Dokumentacja konstrukcyjna · rysunki techniczne i zasady stosowania

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment |
| **Warstwa** | Marka i identyfikacja · dokumentacja konstrukcyjna znaku |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski · geometria **zamknięta** |
| **Data** | 2026-08-14 |
| **Odbiorcy** | projektant wizualny · wdrożeniowiec front-end · dostawca materiałów firmowych · drukarnia · zakład grawerski i hafciarski |
| **Zakres** | moduł konstrukcyjny, rysunki techniczne sygnetu i lockupów, pole ochronne, wielkości minimalne, skalowanie, wyrównanie optyczne, znak w interfejsie, barwy i kontrasty, katalog naruszeń, reprodukcja |
| **Poza zakresem** | idea i narracja znaku (→ `ksiega-znaku.md`), typografia marki (→ `typografia-i-glos.md`), emblematy i ikonografia (→ `emblematy-i-ikony.md`), zastosowania nośnikowe (→ `zastosowania-marki.md`) |
| **Źródła wiążące** | kontrakt systemu projektowego · kierunek systemu projektowego · pliki `zasoby/marka/logo/*.svg` (odczytane co do współrzędnej) · `zasoby/zetony/zetony.css` · `zasoby/zetony/kontrasty.json` |
| **Rysunki** | `03-marka/konstrukcja/rys-01…14.svg` + rasteryzacje `konstrukcja/png/*.png` |
| **Prototyp interaktywny** | `03-marka/konstrukcja-znaku.html` |

---

## Spis treści

1. [Status i sposób czytania dokumentu](#1-status-i-sposób-czytania-dokumentu)
2. [Moduł konstrukcyjny X](#2-moduł-konstrukcyjny-x)
3. [Rysunek konstrukcyjny sygnetu](#3-rysunek-konstrukcyjny-sygnetu)
4. [Kąt rozwarcia i rzeczywista grubość ramienia](#4-kąt-rozwarcia-i-rzeczywista-grubość-ramienia)
5. [Kropka sygnału i linia bazowa](#5-kropka-sygnału-i-linia-bazowa)
6. [Wariant uproszczony](#6-wariant-uproszczony)
7. [Konstrukcja lockupu poziomego](#7-konstrukcja-lockupu-poziomego)
8. [Konstrukcja lockupu pionowego](#8-konstrukcja-lockupu-pionowego)
9. [Logotyp samodzielny](#9-logotyp-samodzielny)
10. [Pole ochronne](#10-pole-ochronne)
11. [Wielkości minimalne](#11-wielkości-minimalne)
12. [Skalowanie](#12-skalowanie)
13. [Wyrównanie optyczne](#13-wyrównanie-optyczne)
14. [Znak na siatce interfejsu](#14-znak-na-siatce-interfejsu)
15. [Barwy znaku i dozwolone pary tło/znak](#15-barwy-znaku-i-dozwolone-pary-tłoznak)
16. [Katalog błędnych użyć](#16-katalog-błędnych-użyć)
17. [Reprodukcja](#17-reprodukcja)
18. [Lista kontrolna poprawnego użycia znaku](#18-lista-kontrolna-poprawnego-użycia-znaku)
19. [Wykaz rysunków i plików](#19-wykaz-rysunków-i-plików)
20. [Decyzje projektowe](#20-decyzje-projektowe)

---

## 1. Status i sposób czytania dokumentu

### 1.1 Czym jest ten dokument

To jest **warstwa konstrukcyjna** dokumentacji marki: rysunek techniczny znaku,
a nie opowieść o znaku. Wszystkie liczby w tym dokumencie zostały **odczytane
z zatwierdzonych plików SVG**, a nie przyjęte. Tam, gdzie liczba jest
wyprowadzona (kąt, grubość prostopadła, środek masy), podano wzór wyprowadzenia,
żeby dało się go sprawdzić.

Znak jest **zamknięty**. Ten dokument opisuje, jak go budować, mierzyć,
umieszczać i reprodukować — nie zawiera propozycji jego zmiany.

### 1.2 Poziomy zobowiązania

W tekście występują trzy znaczniki. Znaczą dokładnie to, co poniżej.

| Znacznik | Znaczenie | Odstępstwo |
|---|---|---|
| **[ZAMKNIĘTE]** | wartość wynika z pliku źródłowego albo z geometrii; zmiana zmienia znak | niemożliwe |
| **[REGUŁA]** | reguła stosowania obowiązująca we wszystkich materiałach | wyłącznie za pisemną zgodą Właściciela znaku |
| **[ZALECENIE]** | rozstrzygnięcie warsztatowe; wolno odstąpić, jeżeli sytuacja tego wymaga | decyzja projektanta, warto odnotować |

### 1.3 Jednostki

| Zapis | Znaczenie |
|---|---|
| **j.** | jednostka siatki konstrukcyjnej — 1/96 kadru sygnetu; wielkość względna, nie fizyczna |
| **X** | moduł konstrukcyjny = wysokość grotu; w kadrze sygnetu **X = 44 j.** |
| **px** | piksel renderowania na ekranie (CSS px, gęstość 1×) |
| **mm** | milimetr na nośniku drukowanym |

Wszystkie wymiary w rozdziałach 2–9 podano w **j.**, bo znak jest niezależny od
wielkości. Wymiary fizyczne pojawiają się dopiero w rozdziale 11 (wielkości
minimalne) i 17 (reprodukcja).

### 1.4 Jak korzystać z rysunków

Każdy rysunek istnieje w dwóch postaciach:

- **SVG** — `03-marka/konstrukcja/rys-NN-*.svg` — do powiększania i pomiaru,
- **PNG** — `03-marka/konstrukcja/png/rys-NN-*.png` — do wklejania w dokumenty,
  które nie przyjmują SVG (1600 px szerokości).

W prototypie `konstrukcja-znaku.html` te same rysunki są **wklejone inline**
i mają przełączalne warstwy: `siatka · osie · wymiary · pole ochronne · sam znak`.

---

## 2. Moduł konstrukcyjny X

### 2.1 Definicja [ZAMKNIĘTE]

> **X = wysokość grotu = 44 j.**
> Mierzone od linii górnej `y = 26` do linii bazowej `y = 70` w kadrze `96 × 96`.

Moduł jest przywiązany do **farby**, nie do kadru pliku. To rozstrzygnięcie ma
konsekwencję praktyczną: wysokość grotu da się zmierzyć na wydruku, na
tabliczce, na zdjęciu szyldu i na zrzucie ekranu — bez dostępu do pliku
źródłowego i bez wiedzy o tym, jaki był kadr. Reguła oparta na kadrze byłaby
niemierzalna wszędzie tam, gdzie znak został wycięty z kadru.

```
                 kadr 96 × 96 j.
   ┌───────────────────────────────────────────┐
   │                                           │
   │   y=26  ─────────────────────────         │  ← linia górna
   │         ▛▚▖    ▛▚▖                        │
   │          ▚▚▖    ▚▚▖                       │
   │   y=48  ──▚▚────▚▚────────────  oś        │   X = 44 j.
   │          ▞▞     ▞▞        ●               │
   │         ▟▞     ▟▞         ▲               │
   │   y=70  ──────────────────┴──────         │  ← LINIA BAZOWA
   │                                           │
   └───────────────────────────────────────────┘
     x=12                x=72  x=76,5   x=89,5
```

### 2.2 Wyprowadzenie wszystkich wymiarów z X [ZAMKNIĘTE]

Rysunek: **`konstrukcja/rys-02-wymiarowanie.svg`**

| # | Wielkość | Jednostki | W module X | Skąd wynika |
|---|---|---:|---:|---|
| 1 | **X — wysokość grotu** | **44,000** | **1,000 X** | `y = 70 − 26` |
| 2 | połowa wysokości grotu | 22,000 | 0,500 X | oś pozioma `y = 48` |
| 3 | **szerokość grotu** | **32,000** | **0,727 X** | `x = 44 − 12` (od A do ostrza C) |
| 4 | **grubość ramienia mierzona poziomo** | **12,000** | **0,273 X** | `x = 24 − 12` (A → B) |
| 5 | **grubość ramienia prostopadle do osi** | **8,879** | **0,202 X** | `12 · 22 / 29,7321` |
| 6 | wybieg poziomy ramienia | 20,000 | 0,455 X | `x = 44 − 24` (B → C) |
| 7 | długość osi ramienia | 29,732 | 0,676 X | `√(20² + 22²)` |
| 8 | **rozstaw grotów (translacja)** | **28,000** | **0,636 X** | grot 2 = grot 1 + 28 j. w osi x |
| 9 | **prześwit między grotami** | **16,000** | **0,364 X** | `40 − 24`; stały na całej wysokości |
| 10 | **średnica kropki** | **13,000** | **0,295 X** | `2 · 6,5` |
| 11 | promień kropki | 6,500 | 0,148 X | `r` z pliku |
| 12 | **prześwit grot 2 → kropka** | **4,500** | **0,102 X** | `76,5 − 72` |
| 13 | zejście środka kropki poniżej osi | 15,500 | 0,352 X | `63,5 − 48` |
| 14 | **szerokość pola farby** | **77,500** | **1,761 X** | `89,5 − 12` |
| 15 | wysokość pola farby | 44,000 | 1,000 X | `70 − 26` |
| 16 | margines kadru z lewej | 12,000 | 0,273 X | `x = 12 − 0` |
| 17 | margines kadru z prawej | 6,500 | 0,148 X | `96 − 89,5` |
| 18 | margines kadru góra = dół | 26,000 | 0,591 X | `26 − 0` oraz `96 − 70` |

**Proporcja pola farby: 1,761 : 1.** Kwadratem jest kadr, nie znak. Ta różnica
jest źródłem większości błędów kadrowania (naruszenia 16.13 i 16.14).

### 2.3 Zależności wewnętrzne, które muszą być zachowane [ZAMKNIĘTE]

Poniższe równości nie są przypadkiem — sprawdzają poprawność każdej rekonstrukcji
znaku:

| Zależność | Zapis | Kontrola |
|---|---|---|
| prześwit = ½ szerokości grotu | `16 = 32 / 2` | 0,364 X = ½ · 0,727 X |
| rozstaw = szerokość grotu − ½ grubości ramienia... | `28 = 32 − 4` | grot 2 wchodzi 4 j. za ostrze grotu 1 w kadrze poziomym |
| rozstaw − grubość ramienia = wybieg | `28 − 12 = 16` | równe prześwitowi |
| średnica kropki ≈ grubość ramienia + 1 j. | `13 ≈ 12 + 1` | kropka jest optycznie odrobinę cięższa od kreski |
| kropka i groty mają wspólną linię bazową | `63,5 + 6,5 = 70` | [ZAMKNIĘTE] warunek czytania „»»." |
| kadr pionowo symetryczny | `26 = 96 − 70` | znak jest wyśrodkowany pionowo w kadrze |

### 2.4 Czego z X nie wolno wyprowadzać

X opisuje **znak**, nie interfejs. Przestrzeń interfejsu opiera się na własnej
jednostce **4 px** (`--dn-od-1`). Nie wolno:

- wyrażać odstępów interfejsu w X,
- wyrażać wielkości kontrolek w X,
- zaokrąglać wielkości znaku do wielokrotności X kosztem siatki 4 px.

Znak **wchodzi na siatkę interfejsu**, a nie odwrotnie: w pasku 48 px godło ma
28 px (7 × 4 px), nie „0,64 X".

---

## 3. Rysunek konstrukcyjny sygnetu

Rysunek: **`konstrukcja/rys-01-siatka-bazowa.svg`**

### 3.1 Siatka bazowa

| Parametr | Wartość |
|---|---|
| kadr | **96 × 96 j.** |
| podziałka pomocnicza | **4 j.** (24 pola w rzędzie) |
| podziałka główna | **12 j.** (8 pól w rzędzie) |
| oś pozioma | `y = 48` |
| oś pionowa | `x = 48` |
| linia górna farby | `y = 26` |
| **linia bazowa** | **`y = 70`** |

Kadr 96 j. jest **czterokrotnością siatki ikon interfejsu** (24 × 24) i
**dwudziestoczterokrotnością jednostki przestrzeni** (4 px). Znak, ikonografia
i układ interfejsu leżą na tej samej siatce — to nie jest zbieg okoliczności,
tylko warunek spójności.

### 3.2 Ścieżki źródłowe [ZAMKNIĘTE]

Zapis dokładnie taki, jak w `zasoby/marka/logo/sygnet.svg`:

```
grot 1   M12 26 H24 L44 48.0 L24 70 H12 L32 48.0 Z
grot 2   M40 26 H52 L72 48.0 L52 70 H40 L60 48.0 Z
kropka   cx="83" cy="63.5" r="6.5"
```

Ścieżki są **zamknięte** (`Z`), złożone wyłącznie z odcinków prostych,
bez krzywych i bez zaokrągleń narożników. Kropka jest **okręgiem**, nie ścieżką
— dzięki temu pozostaje idealnie okrągła po każdej transformacji i po
konwersji na krzywe w programach DTP.

### 3.3 Punkty konstrukcyjne grotu 1

| Punkt | Współrzędne | Rola |
|---|---:|---|
| **A** | 12 · 26 | narożnik górny lewy (początek ścieżki) |
| **B** | 24 · 26 | narożnik górny prawy — początek ramienia górnego |
| **C** | **44 · 48** | **ostrze zewnętrzne** — najdalszy punkt grotu |
| **D** | 24 · 70 | narożnik dolny prawy — koniec ramienia dolnego |
| **E** | 12 · 70 | narożnik dolny lewy |
| **F** | **32 · 48** | **wierzchołek wewnętrzny** — wcięcie |

### 3.4 Punkty konstrukcyjne grotu 2

Grot 2 to **czysta translacja** grotu 1 o **+28 j.** w osi `x`. Zero skalowania,
zero obrotu, zero korekt.

| Punkt | Współrzędne |
|---|---:|
| A′ | 40 · 26 |
| B′ | 52 · 26 |
| **C′** | **72 · 48** |
| D′ | 52 · 70 |
| E′ | 40 · 70 |
| F′ | 60 · 48 |

**Dlaczego translacja, a nie skalowanie.** Dwa groty tej samej wielkości czytają
się jako **powtórzenie** — ten sam gest wykonany dwa razy, czyli przekazanie
dalej. Grot mniejszy czytałby się jako perspektywa (jeden obiekt oddalający się),
grot większy — jako wzmocnienie (crescendo). Znak opowiada delegację, więc oba
groty są równe. [ZAMKNIĘTE]

### 3.5 Punkt kropki

| Parametr | Wartość |
|---|---:|
| środek **S** | 83 · 63,5 |
| promień | 6,5 |
| krawędź górna | `y = 57,0` |
| **krawędź dolna** | **`y = 70,0`** |
| krawędź lewa | `x = 76,5` |
| krawędź prawa | `x = 89,5` |

### 3.6 Kontrola rekonstrukcji

Jeżeli odtwarzasz znak od zera (np. w programie CAD albo przy przygotowaniu
matrycy), sprawdź cztery rzeczy w tej kolejności:

1. **Linia bazowa.** Dolne krawędzie obu grotów **i** dolna krawędź kropki leżą
   na `y = 70`. Jeżeli kropka wisi wyżej albo niżej — znak jest zepsuty.
2. **Prześwit.** Odległość między grotami wynosi 16 j. **na każdej wysokości**,
   nie tylko w połowie. Jeżeli prześwit się zwęża — grot 2 został obrócony
   albo przeskalowany.
3. **Ostrze.** `C = 44 · 48` i `C′ = 72 · 48` leżą **dokładnie na osi poziomej**.
4. **Kropka.** Odległość od ostrza `C′` (72) do lewej krawędzi kropki (76,5)
   wynosi 4,5 j. Jeżeli kropka „przykleiła się" do grotu — została przesunięta.

---

## 4. Kąt rozwarcia i rzeczywista grubość ramienia

Rysunek: **`konstrukcja/rys-03-katy-ramienia.svg`**

### 4.1 Kąty [ZAMKNIĘTE]

| Kąt | Wartość | Wyprowadzenie |
|---|---:|---|
| **kąt rozwarcia grotu** | **95,45°** | między wektorami `C→B` i `C→D`; `cos θ = −84 / 884` |
| nachylenie ramienia do poziomu | 47,73° | `arctan(22 / 20)` |
| nachylenie ramienia do pionu | 42,27° | `arctan(20 / 22)` |
| kąt między ramieniem a osią poziomą znaku | 47,73° | ten sam kąt, mierzony od osi `y = 48` |

**Kąt 95,45° jest rozwarty o 5,45° względem kąta prostego.** To jest cała
osobowość znaku:

| Kąt | Jak by się czytał |
|---|---|
| poniżej 90° (ostry) | **strzałka nawigacji** — „idź dalej", „następny slajd" |
| dokładnie 90° | **narożnik ramki** — element techniczny, nie znak |
| **95,45°** | **grot promptu w ruchu** — gest przekazania, nie polecenie nawigacji |
| powyżej 110° | **daszek / akcent** — traci kierunek, staje się ozdobą |

Zmiana tego kąta zmienia **znaczenie** znaku, a nie jego wygląd. [ZAMKNIĘTE]

### 4.2 Grubość ramienia — dwie różne liczby

To najczęstsze źródło pomyłki przy rekonstrukcji i przy zamówieniu grawerowania.

| Pomiar | Wartość | Kiedy używać |
|---|---:|---|
| grubość mierzona **poziomo** | **12,000 j.** | do rysowania ścieżki (to jest odległość A–B) |
| grubość mierzona **prostopadle do osi ramienia** | **8,879 j.** | do oceny, czy kreska przetrwa druk, grawer i renderowanie |

Wyprowadzenie:

```
wektor ramienia            (Δx, Δy) = (20, 22)
długość osi ramienia       L = √(20² + 22²) = √884 = 29,7321
współczynnik rzutowania    22 / 29,7321 = 0,73993
grubość prostopadła        12 × 0,73993 = 8,8793 j. = 0,2018 X
```

**Konsekwencja praktyczna [REGUŁA].** Wszystkie progi technologiczne (minimalna
kreska w druku, minimalna kreska w renderowaniu, minimalna szerokość rowka
grawerskiego) liczy się od **8,879 j.**, czyli od **0,2018 X** — nigdy od 12 j.
Liczenie od 12 j. zawyża rzeczywistą grubość o **35 %** i prowadzi do
zamówienia znaku, który się zaleje.

### 4.3 Tablica przeliczeniowa kreski

| Wysokość znaku | Kreska (0,2018 X) | Kropka (0,2955 X) | Prześwit (0,3636 X) |
|---:|---:|---:|---:|
| 16 px | 1,48 px | 2,17 px | 2,67 px |
| 24 px | 2,22 px | 3,25 px | 4,00 px |
| 32 px | 2,96 px | 4,33 px | 5,33 px |
| 48 px | 4,44 px | 6,50 px | 8,00 px |
| 96 px | 8,88 px | 13,00 px | 16,00 px |
| 9 mm | 0,83 mm | 1,22 mm | 1,50 mm |
| 15 mm | 1,39 mm | 2,03 mm | 2,50 mm |
| 25 mm | 2,31 mm | 3,38 mm | 4,17 mm |
| 40 mm | 3,70 mm | 5,41 mm | 6,67 mm |

Wysokość znaku = wysokość **grotu** (X), nie wysokość kadru. Jeżeli mierzysz kadr
96 j., przelicz: `X = 0,4583 × wysokość kadru`.

---

## 5. Kropka sygnału i linia bazowa

Rysunek: **`konstrukcja/rys-04-kropka-linia-bazowa.svg`**

### 5.1 Dlaczego kropka leży na linii bazowej, a nie na osi [ZAMKNIĘTE]

Znak czyta się `»».` — jak zapis, nie jak schemat. W zapisie **kropka stoi na
linii pisma**, a nie w połowie wysokości znaków. Dlatego dolna krawędź kropki
(`y = 70`) jest **wspólna** z dolnymi krawędziami obu grotów.

Konsekwencja: środek kropki (`y = 63,5`) leży **15,5 j. poniżej osi grotów**
(`y = 48`), czyli o **0,352 X**.

| Położenie kropki | Jak się czyta |
|---|---|
| **na linii bazowej (`y = 70`)** | **zdanie: „przekazane, przekazane, ruszyło"** — dwa człony i kropka |
| wyśrodkowana na osi (`y = 48`) | schemat blokowy: trzy elementy równorzędne w rzędzie |
| powyżej osi | wykrzyknik albo indeks górny — akcent, nie zamknięcie |
| oderwana od linii bazowej w dół | odsyłacz, przypis — kropka przestaje należeć do zdania |

### 5.2 Odległości kropki

| Pomiar | Wartość | W module X |
|---|---:|---:|
| prześwit od ostrza `C′` (72) do krawędzi kropki (76,5) | 4,500 | 0,102 X |
| odległość środków `C′ → S` | 19,007 | 0,432 X |
| przesunięcie środka: Δx | 11,000 | 0,250 X |
| przesunięcie środka: Δy | 15,500 | 0,352 X |
| kąt odcinka `C′ → S` do poziomu | 54,66° | `arctan(15,5 / 11)` |

Kąt 54,66° jest **większy** od nachylenia ramienia (47,73°). Kropka nie leży na
przedłużeniu ramienia — leży niżej. Gdyby leżała na przedłużeniu, czytałaby się
jako **trzeci grot w perspektywie**; leżąc niżej, czyta się jako **osobny znak
interpunkcyjny**.

### 5.3 Udział kropki w polu farby

| Element | Pole | Udział |
|---|---:|---:|
| grot 1 | 528,00 j² | 44,4 % |
| grot 2 | 528,00 j² | 44,4 % |
| **kropka sygnału** | **132,73 j²** | **11,2 %** |
| **razem** | **1 188,73 j²** | **100 %** |

Kropka zajmuje **jedenaście procent** pola farby, a niesie **jedyną barwę**
w całym systemie. To celowa asymetria: barwa jest zarezerwowana dla znaczenia
(„tu biegnie praca"), więc dostaje najmniejszą powierzchnię i największą uwagę.

### 5.4 Kropka a tętno

W interfejsie kropka sygnału pulsuje (`--dn-czas-tetno` = 2,4 s) — to jedyny
ruch ciągły w systemie. **W znaku marki kropka nie pulsuje nigdy** [REGUŁA]:

| Kontekst | Kropka |
|---|---|
| godło w pasku górnym | **statyczna** |
| lockup w oknie startowym | **statyczna** |
| favicon, ikona aplikacji | **statyczna** |
| materiały firmowe, druk | **statyczna** (oczywiste) |
| kropka stanu `.dn-kropka--tetno` w interfejsie | pulsuje — ale to **komponent**, nie znak |

Uzasadnienie: pulsujące godło znaczyłoby „marka pracuje", co jest bez sensu.
Pulsuje **wskaźnik procesu**, nie tożsamość.

---

## 6. Wariant uproszczony

Rysunek: **`konstrukcja/rys-05-wariant-uproszczony.svg`**

### 6.1 Ścieżki źródłowe [ZAMKNIĘTE]

Z pliku `zasoby/marka/logo/sygnet-uproszczony.svg`:

```
grot     M18 18 H35 L64 48.0 L35 78 H18 L45 48.0 Z
kropka   cx="79" cy="69" r="9"
```

Kadr pozostaje **96 × 96 j.** — ten sam co w wariancie pełnym. To jest istotne:
oba warianty są wymienne w tym samym kontenerze bez przeliczania układu.

### 6.2 Wariant uproszczony NIE jest wycinkiem wariantu pełnego [ZAMKNIĘTE]

Najczęstsze nieporozumienie. Poniżej porównanie realnych wartości:

| Parametr | Sygnet pełny | Sygnet uproszczony | Różnica |
|---|---:|---:|---:|
| **moduł X** (wysokość grotu) | 44,000 j. | **60,000 j.** | +36,4 % |
| szerokość grotu | 32,000 j. | 46,000 j. | +43,8 % |
| grubość ramienia poziomo | 12,000 j. | 17,000 j. | +41,7 % |
| **grubość ramienia prostopadle** | 8,879 j. | **12,223 j.** | +37,7 % |
| grubość prostopadła **w module X** | **0,2018 X** | **0,2037 X** | +0,9 % |
| średnica kropki | 13,000 j. | 18,000 j. | +38,5 % |
| średnica kropki **w module X** | **0,2955 X** | **0,3000 X** | +1,5 % |
| prześwit grot → kropka | 4,500 j. | 6,000 j. | +33,3 % |
| prześwit **w module X** | **0,1023 X** | **0,1000 X** | −2,2 % |
| **kąt rozwarcia** | **95,45°** | **91,94°** | −3,51° |
| pole farby | 77,5 × 44 j. | 70 × 60 j. | — |
| proporcja pola farby | 1,761 : 1 | 1,167 : 1 | — |

**Wnioski konstrukcyjne:**

1. **Proporcje wewnętrzne są zachowane** — grubość, kropka i prześwit
   w module X różnią się o mniej niż 2,5 %. To ten sam alfabet form.
2. **Kąt jest inny o 3,51°** — grot uproszczony jest odrobinę **ostrzejszy**.
   Przy jednym grocie kąt 95,45° czytałby się zbyt miękko; 91,94° przywraca
   kierunek, którego w wariancie pełnym dostarczało powtórzenie.
3. **Znak jest większy w kadrze** — pole farby zajmuje 70 × 60 j. zamiast
   77,5 × 44 j. Przy renderowaniu 16 px daje to kreskę 2,04 px zamiast 1,48 px:
   **wariant uproszczony odzyskuje próg czytelności, którego pełny nie ma.**

### 6.3 Próg przełączenia [REGUŁA]

> **Poniżej 24 px obowiązuje wariant uproszczony.**
> Poniżej 16 px znak nie występuje w ogóle.

Próg nie jest umowny. Wynika z dwóch progów renderowania: **kreska ≥ 2 px**,
**kropka ≥ 3 px średnicy** — poniżej tych wartości wygładzanie krawędzi zamienia
element w szarą plamę.

| Wielkość | Ramię pełnego | Kropka pełnego | Ramię uproszczonego | Kropka uproszczonego |
|---:|---:|---:|---:|---:|
| 12 px | 1,11 px ✗ | 1,62 px ✗ | 1,53 px ✗ | 2,25 px ✗ |
| **16 px** | 1,48 px ✗ | 2,17 px ✗ | **2,04 px ✓** | **3,00 px ✓** |
| 20 px | 1,85 px ✗ | 2,71 px ✗ | 2,55 px ✓ | 3,75 px ✓ |
| **24 px** | **2,22 px ✓** | **3,25 px ✓** | 3,06 px ✓ | 4,50 px ✓ |
| 32 px | 2,96 px ✓ | 4,33 px ✓ | 4,07 px ✓ | 6,00 px ✓ |
| 48 px | 4,44 px ✓ | 6,50 px ✓ | 6,11 px ✓ | 9,00 px ✓ |

**24 px jest ostatnią wielkością**, przy której ramię pełnego sygnetu utrzymuje
2 px. **16 px jest pierwszą wielkością**, przy której wariant uproszczony
odzyskuje oba progi (2,04 px i dokładnie 3,00 px). Stąd dwie granice.

### 6.4 Gdzie wariant uproszczony obowiązuje zawsze [REGUŁA]

| Zastosowanie | Powód |
|---|---|
| favicon 16 / 32 / 48 px | wielkość poniżej progu |
| ikona w pasku zadań systemu | wielkość niepewna, kontrolowana przez system |
| plakietka i medalion poniżej 24 px | wielkość poniżej progu |
| stempel, pieczątka | druk stemplowy zalewa prześwity poniżej 2 mm |
| haft | ścieg nie odwzorowuje prześwitu 16 j. przy rozmiarach użytkowych |
| grawer poniżej 15 mm | rowek zlewa prześwit |

---

## 7. Konstrukcja lockupu poziomego

Rysunek: **`konstrukcja/rys-06-lockup-poziomy.svg`**
Plik: `zasoby/marka/logo/logo-poziomy.svg` · kadr **225 × 96 j.**

### 7.1 Składniki i ich położenie [ZAMKNIĘTE]

| Składnik | Transformacja w pliku | Obwiednia farby (x) | Obwiednia farby (y) |
|---|---|---:|---:|
| **sygnet** | `translate(0,16) scale(0.6667)` | 8,00 … 59,75 | 33,25 … 62,75 |
| **DANACO** | `translate(76,51)`, glify `scale(0.034)` | 77,50 … 209,88 | 26,63 … 51,00 |
| **CONSOLE** | `translate(76,67)`, glify `scale(0.0125)` | 76,63 … 150,13 | 58,13 … 67,00 |
| **kropka po CONSOLE** | `cx=158 cy=63 r=3.4` | 154,50 … 161,50 | 59,60 … 66,40 |
| **obwiednia całości** | — | **8,00 … 209,88** | **26,63 … 67,25** |

### 7.2 Moduł X w lockupie poziomym

Sygnet jest przeskalowany do **0,6667**, więc:

```
X_lockup = 44 × 0,6667 = 29,333 j.
```

Wszystkie reguły odstępów i pola ochronnego liczy się od **tego** X, nie od 44.

### 7.3 Relacja sygnet ↔ logotyp [ZAMKNIĘTE]

| Pomiar | Wartość | Wyrażony inaczej |
|---|---:|---|
| **odstęp: prawa krawędź kadru sygnetu (64) → początek bloku tekstu (76)** | **12,000 j.** | **0,409 X** · **½ wersalika DANACO** |
| odstęp od obwiedni farby sygnetu (59,75) do farby DANACO (77,50) | 17,750 j. | 0,605 X |
| wysokość wersalika DANACO | 23,800 j. | 0,811 X |
| wysokość wersalika CONSOLE | 8,725 j. | 0,297 X |
| **interlinia: linia bazowa DANACO (51) → CONSOLE (67)** | **16,000 j.** | **0,672 wersalika DANACO** |
| szerokość wyrazu DANACO | 132,375 j. | 4,513 X |
| szerokość wyrazu CONSOLE (bez kropki) | 73,500 j. | 2,506 X |

**Reguła robocza [REGUŁA]:** odstęp między sygnetem a logotypem =
**½ wysokości wersalika `DANACO`**. Ta postać reguły jest odporna na
przeskalowanie sygnetu, bo odwołuje się do elementu, który zawsze jest obecny
w lockupie.

```
   ┌──── kadr sygnetu 0…64 ────┐  12  ┌──── blok tekstu 76…210 ────────┐
   │                           │ ←──→ │                                │
   │      ▛▚ ▛▚   ●            │      │   D A N A C O                  │  ← linia bazowa y=51
   │                           │      │   C O N S O L E ●              │  ← linia bazowa y=67
   └───────────────────────────┘      └────────────────────────────────┘
              ↑ oś grotów y = 48  ≈  optyczna oś bloku tekstu (47,1)
```

### 7.4 Wyrównanie pionowe [ZAMKNIĘTE]

| Wielkość | Wartość |
|---|---:|
| oś grotów sygnetu | `y = 48,00` |
| górna krawędź wersalika DANACO | `y = 27,20` |
| linia bazowa CONSOLE | `y = 67,00` |
| **optyczna oś bloku dwuwierszowego** | `(27,20 + 67,00) / 2 = 47,10` |
| **odchyłka sygnetu od osi bloku** | **+0,90 j.** (sygnet niżej o 0,03 X) |

Odchyłka 0,9 j. jest **celowa**. Blok tekstu ma masę skupioną w górnym wierszu
(DANACO 23,8 j. wysokości wobec CONSOLE 8,7 j.), więc jego środek optyczny leży
wyżej niż środek geometryczny. Sygnet przesunięty o 0,9 j. w dół **równoważy**
tę asymetrię.

**Nie wolno wyrównywać sygnetu do linii bazowej DANACO ani do środka
geometrycznego obwiedni tekstu** — obie te operacje rozjeżdżają lockup
(naruszenia 16.05 i 16.06).

### 7.5 Kropka po CONSOLE

| Parametr | Wartość | Reguła |
|---|---:|---|
| średnica | 6,800 j. | 0,779 wersalika CONSOLE |
| środek `cy` | 63,00 | połowa wysokości wersalika CONSOLE (62,64) |
| odstęp od litery `E` (150,13) do krawędzi kropki (154,50) | 4,375 j. | **½ wersalika CONSOLE** |

Kropka po `CONSOLE` powtarza kropkę sygnetu w skali typograficznej. Różnica
położenia jest zamierzona: **w sygnecie kropka stoi na linii bazowej,
w logotypie — na osi wersalika**, bo w wierszu wersalikowym oś optyczna jest
środkiem wersalika, nie linią pisma.

### 7.6 Warianty lockupu poziomego

| Plik | Groty | Kropka sygnetu | Logotyp | Kropka po CONSOLE |
|---|---|---|---|---|
| `logo-poziomy.svg` | `#181818` | `#3B6FE0` | `#181818` | `#3B6FE0` |
| `logo-poziomy-na-ciemnym.svg` | `#F4F4F4` | `#5C8CEC` | `#F4F4F4` | `#3B6FE0` |
| `logo-poziomy-mono-czarny.svg` | `#181818` | `#181818` | `#181818` | `#3B6FE0` |
| `logo-poziomy-mono-bialy.svg` | `#F4F4F4` | `#F4F4F4` | `#F4F4F4` | `#3B6FE0` |

**Uwaga wdrożeniowa.** W plikach mono kropka po `CONSOLE` pozostaje w barwie
sygnałowej. Przy zamówieniu reprodukcji **jednokolorowej** (jeden kolor farby,
grawer, tłoczenie) tę kropkę trzeba **ręcznie przestawić na kolor znaku** —
patrz rozdział 17.5. To jedyne miejsce w komplecie plików, które wymaga
ingerencji przed przekazaniem do warsztatu.

---

## 8. Konstrukcja lockupu pionowego

Rysunek: **`konstrukcja/rys-07-lockup-pionowy.svg`**
Plik: `zasoby/marka/logo/logo-pionowy.svg` · kadr **159 × 172 j.**

### 8.1 Składniki i ich położenie [ZAMKNIĘTE]

| Składnik | Transformacja | Obwiednia farby (x) | Obwiednia farby (y) |
|---|---|---:|---:|
| **sygnet** | `translate(31.7,0)` — **skala 1:1** | 43,63 … 121,25 | 26,00 … 70,00 |
| **DANACO** | `translate(12.0,136)` | 13,50 … 145,88 | 111,63 … 136,00 |
| **CONSOLE** | `translate(42.2,154)` | 42,88 … 116,25 | 145,13 … 154,00 |
| **kropka po CONSOLE** | `cx=124.2 cy=150 r=3.4` | 120,80 … 127,60 | 146,60 … 153,40 |
| **obwiednia całości** | — | **13,50 … 145,88** | **26,00 … 154,25** |

W lockupie pionowym sygnet występuje w **skali 1:1**, więc **X = 44 j.**

### 8.2 Trzy różne osie centrowania [ZAMKNIĘTE]

To najbardziej niejawna decyzja w całym systemie znaku — i najczęściej psuta
przy odtwarzaniu układu.

| Element | Sposób centrowania | Środek | Odchyłka od osi kadru (79,5) |
|---|---|---:|---:|
| **sygnet** | przez **kadr 96 × 96** (`31,7 … 127,7`) | 79,70 | +0,20 |
| **DANACO** | przez **obwiednię farby** | 79,69 | +0,19 |
| **CONSOLE** | przez **obwiednię farby BEZ kropki** | 79,56 | +0,06 |
| kropka po CONSOLE | **zwisa poza blok** — nie bierze udziału w centrowaniu | — | — |

**Wniosek [REGUŁA]:**

1. **Sygnet centruje się kadrem, nie farbą.** Kadr 96 × 96 zawiera już korektę
   optyczną (rozdz. 13). Centrowanie sygnetu obwiednią farby przesunęłoby go
   o **2,75 j.** w lewo i przechyliło cały układ.
2. **CONSOLE centruje się bez kropki.** Kropka jest elementem **zwisającym**,
   jak wiszący znak interpunkcyjny w składzie typograficznym. Gdyby brała
   udział w centrowaniu, wyraz `CONSOLE` przesunąłby się o **5,7 j.** w lewo
   względem `DANACO` — i wiersz drugi wyraźnie „uciekłby".

### 8.3 Odstępy pionowe

| Pomiar | Wartość | W module X |
|---|---:|---:|
| **linia bazowa sygnetu (70) → górna krawędź wersalika DANACO (112,2)** | **42,200 j.** | **0,959 X ≈ 1X** |
| linia bazowa DANACO (136) → linia bazowa CONSOLE (154) | 18,000 j. | 0,409 X |
| margines kadru: góra (0 → 26) | 26,000 j. | 0,591 X |
| margines kadru: dół (154,25 → 172) | 17,750 j. | 0,403 X |

**Reguła robocza [REGUŁA]:** odstęp sygnet → logotyp w układzie pionowym
wynosi **1X** (dokładnie 0,96 X — zaokrąglenie do 1X jest bezpieczne i łatwe
do skontrolowania na wydruku).

**Interlinia jest większa niż w lockupie poziomym** (18 j. wobec 16 j.). Powód:
w układzie pionowym oba wiersze są centrowane, więc oko potrzebuje więcej
prześwitu, żeby je rozdzielić — w układzie poziomym rozdziela je już wyrównanie
do lewej i obecność sygnetu.

### 8.4 Szerokość odniesienia

Szerokość wyrazu `DANACO` (**132,375 j.**) wyznacza szerokość całego układu
pionowego. Wszystkie pozostałe elementy są węższe:

```
   sygnet — kadr        96,000 j.   (72,5 % szerokości DANACO)
   sygnet — farba       77,625 j.   (58,6 %)
   DANACO               132,375 j.  (100 %)   ← szerokość odniesienia
   CONSOLE              73,375 j.   (55,4 %)
   CONSOLE + kropka     84,750 j.   (64,0 %)
```

Przy skalowaniu lockupu pionowego do zadanej szerokości **przelicza się od
`DANACO`**, nie od kadru: `skala = szerokość_docelowa / 159` (kadr) albo
`skala = szerokość_docelowa / 132,375` (farba) — konsekwentnie jedna metoda
w całym materiale.

---

## 9. Logotyp samodzielny

Plik: `zasoby/marka/logo/logotyp.svg` · kadr **149 × 56 j.**

### 9.1 Składniki [ZAMKNIĘTE]

| Składnik | Obwiednia (x) | Obwiednia (y) | Linia bazowa |
|---|---:|---:|---:|
| **DANACO** | 1,50 … 133,88 | 9,63 … 34,50 | `y = 34` |
| **CONSOLE** | 0,63 … 74,13 | 41,13 … 50,25 | `y = 50` |
| **kropka** | 78,50 … 85,50 | 42,50 … 49,50 | środek `y = 46` |

### 9.2 Moduł zastępczy [REGUŁA]

W logotypie samodzielnym **nie ma grotu**, więc nie ma X. Modułem zastępczym
jest **wysokość wersalika `DANACO` = 23,80 j.**

| Pomiar | Wartość | W module zastępczym |
|---|---:|---:|
| wysokość wersalika DANACO | 23,800 | 1,000 |
| wysokość wersalika CONSOLE | 8,725 | 0,367 |
| interlinia (34 → 50) | 16,000 | 0,672 |
| odstęp CONSOLE → kropka | 4,375 | 0,184 |
| **pole ochronne** | **11,900** | **0,500** |

Interlinia 16 j. jest identyczna z lockupem poziomym — logotyp samodzielny to
dokładnie ten sam blok tekstowy, wyjęty z lockupu.

### 9.3 Kiedy wolno użyć logotypu bez sygnetu [REGUŁA]

| Wolno | Nie wolno |
|---|---|
| stopka dokumentu, w którym sygnet już wystąpił na stronie tytułowej | pierwsze wystąpienie marki w materiale |
| pasek boczny prezentacji, gdy sygnet stoi na slajdzie tytułowym | ikona, favicon, awatar |
| nagłówek listy dokumentów w Library Explorer | ekran startowy i ekran logowania |
| stopka wiadomości e-mail | oznaczenie własności na przedmiocie |

Zasada nadrzędna: **logotyp samodzielny jest podpisem, nie znakiem.** Podpis
wymaga, żeby znak wystąpił wcześniej.

---

## 10. Pole ochronne

Rysunki: **`rys-08-pole-ochronne-sygnet.svg`**, **`rys-09-pole-ochronne-lockup-poziomy.svg`**,
**`rys-10-pole-ochronne-lockup-pionowy.svg`**

### 10.1 Dwa poziomy pola [REGUŁA]

> **Pole minimalne = ½X** — obowiązuje **zawsze i bezwzględnie**.
> **Pole ekspozycyjne = 1X** — obowiązuje przy **samodzielnej ekspozycji znaku**.

| Poziom | Wartość | Kiedy obowiązuje | Status |
|---|---:|---|---|
| **pole minimalne** | **½X** | zawsze — każde wystąpienie znaku bez wyjątku | **[REGUŁA]** |
| **pole ekspozycyjne** | **1X** | znak jako element samodzielny: strona tytułowa, okładka, plansza, znak na tle obcym (zdjęcie, blok barwny), oznaczenie na przedmiocie | **[ZALECENIE]** |

Rozróżnienie jest praktyczne, nie kosmetyczne. W kokpicie (pasek 48 px, karta,
medalion) znak sąsiaduje z elementami interfejsu i pole ½X jest wszystkim, co
zmieści się bez rozerwania układu. Na okładce, na tabliczce i na koszulce znak
stoi sam i pole ½X wygląda na przycięcie — 1X przywraca mu oddech.

Zgodność z `ksiega-znaku.md` rozdz. 5: **wartość wiążąca pozostaje ½X.**
Pole 1X jest **zaleceniem ekspozycyjnym**, nie zaostrzeniem reguły.

### 10.2 Pole liczy się od kadru, nie od farby [REGUŁA]

To rozstrzygnięcie ma trzy powody:

1. **Mierzalność.** Kadr jest krawędzią pliku i krawędzią kontenera w kodzie.
   Obwiednia farby wymaga pomiaru punktu skrajnego kropki, którego przy małych
   wielkościach nie da się wskazać.
2. **Korekta optyczna.** Kadr 96 × 96 zawiera już przesunięcie równoważące masę
   znaku (rozdz. 13). Liczenie pola od farby zniosłoby tę korektę.
3. **Powtarzalność.** Reguła oparta na kadrze daje ten sam wynik w SVG, w CSS,
   w programie DTP i na rysunku warsztatowym.

### 10.3 Tabela pól ochronnych [REGUŁA]

| Konfiguracja | Kadr | X w kadrze | Pole ½X | Kadr z polem ½X | Pole 1X | Kadr z polem 1X |
|---|---:|---:|---:|---:|---:|---:|
| **Sygnet pełny** | 96 × 96 | 44,000 | **22,00** | 140 × 140 | 44,00 | 184 × 184 |
| **Sygnet uproszczony** | 96 × 96 | 60,000 | **30,00** | 156 × 156 | 60,00 | 216 × 216 |
| **Lockup poziomy** | 225 × 96 | 29,333 | **14,67** | 254,3 × 125,3 | 29,33 | 283,7 × 154,7 |
| **Lockup pionowy** | 159 × 172 | 44,000 | **22,00** | 203 × 216 | 44,00 | 247 × 260 |
| **Logotyp samodzielny** | 149 × 56 | brak grotu → wersalik 23,80 | **11,90** | 172,8 × 79,8 | 23,80 | 196,6 × 103,6 |

### 10.4 Pole ochronne w pikselach — tabela robocza

Wyliczone dla pola **½X**, od **kadru**:

| Konfiguracja | Wysokość znaku | Pole ochronne | Konfiguracja | Wysokość znaku | Pole ochronne |
|---|---:|---:|---|---:|---:|
| Sygnet | 24 px | 5,5 px → **6 px** | Lockup poziomy | 32 px | 4,9 px → **5 px** |
| Sygnet | 28 px | 6,4 px → **7 px** | Lockup poziomy | 48 px | 7,3 px → **8 px** |
| Sygnet | 32 px | 7,3 px → **8 px** | Lockup poziomy | 64 px | 9,8 px → **10 px** |
| Sygnet | 48 px | 11,0 px → **11 px** | Lockup pionowy | 96 px | 12,3 px → **13 px** |
| Sygnet | 96 px | 22,0 px → **22 px** | Lockup pionowy | 160 px | 20,5 px → **21 px** |
| Sygnet | 128 px | 29,3 px → **30 px** | Logotyp | 56 px | 11,9 px → **12 px** |

Wzory (wysokość podana jako wysokość **kadru** w px):

```
sygnet pełny        pole = wysokość × 44 / 96 / 2 = wysokość × 0,2292
sygnet uproszczony  pole = wysokość × 60 / 96 / 2 = wysokość × 0,3125
lockup poziomy      pole = wysokość × 29,333 / 96 / 2 = wysokość × 0,1528
lockup pionowy      pole = wysokość × 44 / 172 / 2 = wysokość × 0,1279
logotyp samodzielny pole = wysokość × 23,8 / 56 / 2 = wysokość × 0,2125
```

Wynik **zaokrągla się w górę do pełnego piksela**, a w interfejsie —
**do najbliższej wielokrotności 4 px w górę**, żeby znak nie łamał siatki
przestrzeni.

### 10.5 Co NIE MOŻE wejść w pole ochronne [REGUŁA]

| Kategoria | Przykłady |
|---|---|
| **tekst** | podpisy, daty, numery wersji, adresy, hasła, nazwy wydarzeń |
| **inne znaki graficzne** | logotypy partnerów, plakietki certyfikatów, ikony sklepów z aplikacjami, kody QR |
| **linie i ramki** | separatory, obrysy, linijki, podkreślenia, ramki dekoracyjne |
| **krawędzie nośnika** | brzeg kartki, brzeg wizytówki, brzeg ekranu, spad drukarski, zgięcie |
| **krawędzie bloku** | brzeg zdjęcia, brzeg bloku barwnego, brzeg karty, brzeg modala |
| **elementy interfejsu** | przyciski, pola, ikony, plakietki, awatary, paski postępu |
| **zmiana tła** | granica gradientu, granica dwóch barw tła, granica dwóch powierzchni |

### 10.6 Co MOŻE wejść w pole ochronne

| Dopuszczone | Warunek |
|---|---|
| jednolita powierzchnia tła, na której znak stoi | musi być **jednolita** na całym polu |
| bardzo delikatna faktura papieru albo materiału | nie może tworzyć rozpoznawalnego wzoru |
| ten sam kolor tła w innej powierzchni (np. karta na tle) | granica powierzchni musi leżeć **poza** polem |

### 10.7 Wyjątki [ZALECENIE]

| Sytuacja | Rozstrzygnięcie |
|---|---|
| **pasek górny aplikacji (48 px)** | pole pionowe redukuje się do **½ wartości** — godło 28 px daje 10 px góra i dół; w poziomie obowiązuje pełne ½X (12 px od krawędzi ekranu) |
| **favicon i ikona aplikacji** | pole zawiera się w kadrze pliku — ikona ma własne marginesy zapisane w plikach `ikona-aplikacji/*` |
| **znak wodny** | pole nie obowiązuje: znak wodny leży **pod** treścią, nie obok niej |
| **stopka dokumentu przy krawędzi** | pole obowiązuje bez zmian; jeżeli się nie mieści — **zmniejsz znak, nie pole** |
| **medalion nadawcy w oknie komunikacji** | pole liczy się od krawędzi medalionu, nie od krawędzi wpisu |

---

## 11. Wielkości minimalne

Rysunek: **`konstrukcja/rys-11-drabina-wielkosci.svg`**

### 11.1 Tabela wielkości minimalnych [REGUŁA]

Podane wartości to **szerokość** znaku, z wyjątkiem sygnetu, gdzie podano
**wysokość kadru** (kadr jest kwadratowy).

| Konfiguracja | Min. ekran | Min. druk offset/cyfra | Min. sitodruk, flexo | Min. grawer, haft, tłoczenie | Uwagi |
|---|---:|---:|---:|---:|---|
| **Lockup poziomy** | **120 px** | **28 mm** | 36 mm | 40 mm | poniżej — użyj sygnetu, nie zmniejszaj |
| **Lockup pionowy** | **120 px** | **28 mm** | 36 mm | 40 mm | wysokość rośnie proporcjonalnie (× 1,082) |
| **Logotyp samodzielny** | **90 px** | **20 mm** | 26 mm | 30 mm | wyłącznie jako podpis (rozdz. 9.3) |
| **Sygnet pełny** | **24 px** | **9 mm** | 12 mm | 15 mm | **próg przełączenia wariantu** |
| **Sygnet uproszczony** | **16 px** | **6 mm** | 8 mm | 10 mm | dolna granica całego systemu |

**Poniżej wielkości minimalnej nie zmniejsza się znaku — zmienia się
konfigurację na prostszą:**

```
   lockup poziomy  →  lockup pionowy  →  logotyp  →  sygnet pełny  →  sygnet uproszczony
      ≥ 120 px          ≥ 120 px         ≥ 90 px       ≥ 24 px             ≥ 16 px
      ≥ 28 mm           ≥ 28 mm          ≥ 20 mm       ≥ 9 mm              ≥ 6 mm
```

### 11.2 Uzasadnienie progu 120 px dla lockupu

Przy szerokości lockupu 120 px sygnet ma wysokość:

```
X_px = 120 × 29,333 / 225 = 15,64 px
kreska = 15,64 × 0,2018 = 3,16 px  ✓
```

kreska jest bezpieczna, ale **wersalik CONSOLE** ma wtedy:

```
CONSOLE_cap = 120 × 8,725 / 225 = 4,65 px
```

4,65 px to **dolna granica czytelności wersalików rozstrzelonych w Plex Mono**.
Poniżej tej wartości `CONSOLE` zamienia się w szarą kreskę, a lockup przestaje
być lockupem — zostaje logotyp z szumem. Stąd 120 px, a nie mniej.

### 11.3 Uzasadnienie progu 9 mm dla druku

| Wysokość sygnetu | Kreska (0,2018 X) | Kropka (0,2955 X) | Prześwit (0,3636 X) | Ocena |
|---:|---:|---:|---:|---|
| 6 mm | 0,55 mm | 0,81 mm | 1,00 mm | ✗ kreska poniżej progu |
| 8 mm | 0,74 mm | 1,08 mm | 1,33 mm | ✗ kropka poniżej progu |
| **9 mm** | **0,83 mm** | **1,22 mm** | **1,50 mm** | **✓ oba progi spełnione** |
| 12 mm | 1,11 mm | 1,62 mm | 2,00 mm | ✓ |
| 15 mm | 1,39 mm | 2,03 mm | 2,50 mm | ✓ margines na grawer |

Progi technologiczne przyjęte za standardem druku offsetowego:
**kreska ≥ 0,8 mm** (poniżej farba zaczyna zamykać prześwity przy rozlewie),
**element okrągły ≥ 1,2 mm** (poniżej okrąg traci okrągłość).

Dla grawerowania, tłoczenia i haftu progi rosną do **1,5 mm / 2,0 mm** — stąd
minimum **15 mm** dla tych technik.

### 11.4 Wielkości zalecane [ZALECENIE]

| Zastosowanie | Konfiguracja | Wielkość |
|---|---|---:|
| favicon karty przeglądarki | sygnet uproszczony | 16 · 32 · 48 px |
| ikona aplikacji na pulpicie | sygnet pełny w kaflu | 512 · 1024 px |
| **pasek górny aplikacji (48 px)** | **sygnet pełny + logotyp** | **godło 28 px** |
| medalion nadawcy w oknie komunikacji | sygnet pełny | 28 px |
| karta środowiska — emblemat | emblemat środowiska | 24 px (ekspozycja 40 px) |
| **okno startowe (ładowanie)** | **lockup pionowy** | **160 px szerokości** |
| okno rejestracji i logowania | lockup pionowy | 140 px szerokości |
| Always On Display — rdzeń | sygnet pełny | 96 px |
| papier firmowy A4 | lockup poziomy | 45 mm szerokości |
| wizytówka 90 × 50 mm | lockup poziomy | 28 mm szerokości |
| prezentacja — slajd tytułowy | lockup pionowy | 60 mm szerokości |
| stopka dokumentu PDF | logotyp samodzielny | 22 mm szerokości |
| koszulka, torba (haft) | sygnet uproszczony | 40 mm |
| tabliczka grawerowana | sygnet pełny | 20 mm |

---

## 12. Skalowanie

### 12.1 Zasada podstawowa [REGUŁA]

> **Znak skaluje się wyłącznie proporcjonalnie.**
> Jedyna dopuszczalna transformacja to `scale(k, k)` — ten sam współczynnik
> w obu osiach.

Zakazane bezwzględnie:

| Operacja | Zapis | Dlaczego zakazana |
|---|---|---|
| skalowanie nieproporcjonalne | `scale(1.2, 1)` | zmienia kąt rozwarcia grotu — a kąt niesie znaczenie (rozdz. 4.1) |
| pochylenie | `skewX()`, `skewY()` | to samo; dodatkowo zamienia kropkę w elipsę |
| obrót | `rotate()` | znak ma kierunek czytania (od lewej do prawej); obrót go niszczy |
| odbicie | `scale(-1, 1)` | grot zaczyna wskazywać wstecz — znaczenie odwrotne |
| skalowanie samego logotypu w lockupie | — | zrywa relację ½ wersalika (rozdz. 7.3) |
| skalowanie samej kropki | — | zrywa proporcję 0,295 X |

### 12.2 Jak zmienia się kąt przy skalowaniu nieproporcjonalnym

Ilustracja skali problemu — kąt rozwarcia grotu przy rozciągnięciu w poziomie:

| Rozciągnięcie poziome | Efektywny `Δx` | Kąt rozwarcia | Ocena |
|---:|---:|---:|---|
| 0,80 | 16,0 | 107,9° | daszek — znak traci kierunek |
| 0,90 | 18,0 | 101,4° | wyraźnie miększy |
| **1,00** | **20,0** | **95,45°** | **poprawny** |
| 1,10 | 22,0 | 90,0° | narożnik ramki |
| 1,20 | 24,0 | 85,0° | strzałka nawigacji |
| 1,50 | 30,0 | 72,5° | ostrze — obcy znak |

Rozciągnięcie o **10 %** zamienia grot promptu w kąt prosty. To nie jest
subtelna zmiana proporcji, tylko zmiana pojęcia.

### 12.3 Skalowanie sygnetu — obrys nie występuje

Sygnet jest zbudowany z **wypełnień**, nie z obrysów. Nie ma tu problemu
skalowania kreski: `fill` skaluje się razem z kształtem, więc `scale(k)` daje
zawsze poprawny znak. Jedynym progiem jest wielkość minimalna (rozdz. 11).

### 12.4 Skalowanie emblematów środowisk — obrys 1,75 [REGUŁA]

Rysunek: **`konstrukcja/rys-13-emblemat-siatka.svg`**

Emblematy środowisk (`srodowisko-talkin`, `-workspace`, `-codestudio`,
`-multitaskingai`) są zbudowane inaczej: **siatka 24 × 24, `fill="none"`,
`stroke="currentColor"`, `stroke-width="1.75"`, `stroke-linecap/linejoin="round"`.**

Obrys **skaluje się razem z kształtem**, więc rzeczywista grubość kreski
w pikselach zmienia się liniowo z wielkością renderowania:

| Renderowanie | `stroke-width` w siatce | Rzeczywista kreska | Postępowanie |
|---:|---:|---:|---|
| 14 px | 1,75 | **1,02 px** | pozostaw 1,75 — przeglądarka wygładza; poniżej 14 px emblematu nie stosujemy |
| **16 px** | 1,75 | **1,17 px** | pozostaw 1,75 — wielkość zalecana dla list i pozycji nawigacji |
| **20 px** | 1,75 | **1,46 px** | pozostaw 1,75 — wielkość zalecana dla nagłówków paneli |
| **24 px** | 1,75 | **1,75 px** | renderowanie 1:1 — kreska idealnie na pikselu |
| 32 px | 1,75 | 2,33 px | pozostaw 1,75 |
| **40 px** | 1,75 | **2,92 px** | ekspozycja na karcie środowiska — pozostaw 1,75 |
| 64 px | 1,75 | 4,67 px | **obniż do 1,5** — inaczej emblemat zgrubieje optycznie |
| 96 px i więcej | 1,5 → 1,25 | 6,00 → 5,00 px | **obniż do 1,25** przy ekspozycji plakatowej |

**Zasada [REGUŁA]:** do **48 px włącznie** obrys pozostaje **1,75** bez wyjątku
— to jest ta sama liczba, która trzyma spójność całego zestawu ikon interfejsu.
Powyżej 48 px, wyłącznie w ekspozycji (karta tytułowa, plansza, materiał
drukowany), wolno obniżyć obrys do 1,5 lub 1,25, żeby zachować **wrażenie**
tej samej wagi kreski.

**Nie wolno używać `vector-effect="non-scaling-stroke"`.** Ten atrybut trzyma
grubość kreski w pikselach ekranu niezależnie od skali — czyli robi dokładnie
to, czego nie chcemy: emblemat 96 px z kreską 1,75 px wygląda jak rysunek
techniczny, nie jak znak.

### 12.5 Skalowanie w kodzie [REGUŁA]

| Sposób | Zapis | Ocena |
|---|---|---|
| **SVG inline + `width`/`height`** | `<svg width="28" height="28" viewBox="0 0 96 96">` | **zalecany** — `currentColor` działa, brak żądania sieciowego |
| SVG inline + CSS | `.godlo svg { width: 28px; height: auto; }` | zalecany |
| `<img src="sygnet.svg" width="28">` | — | dopuszczalny, ale `currentColor` nie działa |
| `background-image` | — | dopuszczalny dla tła dekoracyjnego; znak traci dostępność (brak `role="img"`) |
| `transform: scale()` na kontenerze | — | **niedopuszczalny** — rozmywa krawędzie i psuje wyrównanie do siatki |
| CSS `zoom` | — | **niedopuszczalny** — ten sam problem |

### 12.6 Skalowanie a siatka przestrzeni [REGUŁA]

W interfejsie znak przyjmuje **wyłącznie wielokrotności 4 px**:

```
   16 · 20 · 24 · 28 · 32 · 40 · 48 · 64 · 96 · 128 · 160 px
```

Wartości pośrednie (np. 30 px, 45 px) są zakazane — łamią siatkę przestrzeni
i powodują, że znak renderuje się na półpikselu.

---

## 13. Wyrównanie optyczne

Rysunek: **`konstrukcja/rys-04-kropka-linia-bazowa.svg`**

### 13.1 Trzy różne środki tego samego znaku [ZAMKNIĘTE]

| Środek | Wartość `x` | Jak liczony |
|---|---:|---|
| **środek kadru** | **48,000** | `96 / 2` |
| środek prostokąta farby | 50,750 | `(12 + 89,5) / 2` |
| **środek masy farby** | **46,578** | średnia ważona polem: groty (2 × 528 j²) + kropka (132,73 j²) |

Pionowo sytuacja jest prosta: środek kadru `y = 48` = środek prostokąta farby
`y = 48`; środek masy wypada na `y = 49,73` (kropka ciągnie masę w dół o 1,73 j.).

### 13.2 Dlaczego kadr wypada dokładnie tam, gdzie wypada

```
   masa 46,578          kadr 48,000          prostokąt farby 50,750
        │                    │                        │
   ─────●────────────────────●────────────────────────●─────►  x
        └──── 1,422 ─────────┴────────── 2,750 ───────┘
                  34 %                     66 %
```

Środek kadru **nie jest** średnią tych dwóch punktów. Leży **bliżej środka
masy** (34 % drogi), bo:

- **wyrównanie do prostokąta farby** przesunęłoby znak w prawo — groty
  (88,8 % masy) wylądowałyby po lewej stronie kontenera, a znak czytałby się
  jako „zsunięty w lewo z kropką doklejoną z prawej",
- **wyrównanie do środka masy** przesunęłoby znak w lewo — kropka podeszłaby
  do prawej krawędzi na 4,08 j. i wyglądałaby na przycinaną.

Kadr 96 × 96 jest **rozstrzygnięciem** tego konfliktu, zapisanym raz i na
zawsze w pliku. Dlatego reguła brzmi:

> **Znak wyrównuje się w kontenerze ZAWSZE przez kadr, nigdy przez obwiednię
> farby.** [REGUŁA]

### 13.3 Praktyka wyrównania w kontenerze

| Kontener | Postępowanie |
|---|---|
| kwadratowy (medalion, awatar, kafel) | `viewBox="0 0 96 96"`, `width = height = bok kontenera` — kadr wypełnia kontener, gotowe |
| prostokątny poziomy (pasek, przycisk) | wyśrodkuj **kadr** w pionie; w poziomie odsuń kadr od krawędzi o pole ochronne |
| prostokątny pionowy (karta) | wyśrodkuj **kadr** w poziomie |
| przy tekście w jednym wierszu | **oś grotów** (`y = 48` w kadrze, czyli połowa wysokości kadru) na wysokości **środka wersalika** tekstu |
| przy tekście dwuwierszowym | patrz rozdz. 7.4 — oś grotów na optycznej osi bloku |

### 13.4 Wyrównanie sygnetu do tekstu — wzór

Dla sygnetu o wysokości kadru `H` obok tekstu o wysokości wersalika `C`
z linią bazową na `B`:

```
oś grotów sygnetu  =  B − C / 2
górna krawędź kadru sygnetu  =  B − C / 2 − H / 2
```

Przykład (pasek górny, godło 28 px, `DANACO` 13 px Plex Sans, wersalik ≈ 9,3 px,
linia bazowa 30 px od góry paska):

```
oś grotów  = 30 − 9,3 / 2 = 25,35 px
górna krawędź kadru = 25,35 − 14 = 11,35 px  →  zaokrąglone do siatki: 10 px
```

Zaokrąglenie w dół do 10 px daje symetryczne marginesy 10/10 w pasku 48 px
i mieści się w siatce 4 px. To jest wartość zapisana w rysunku 12.

---

## 14. Znak na siatce interfejsu

Rysunek: **`konstrukcja/rys-12-znak-w-interfejsie.svg`**

### 14.1 Pasek górny — 48 px [REGUŁA]

Pasek górny (`--dn-wym-pasek` = 48 px) jest **zawsze atramentowy**
(`--dn-rama` = `szary-925` #131313) w obu motywach. Godło stoi na nim
w **wariancie na ciemnym tle**.

| Parametr | Wartość | Wyprowadzenie |
|---|---:|---|
| wysokość paska | **48 px** | `--dn-wym-pasek` |
| **wysokość godła (kadr)** | **28 px** | 7 × jednostka 4 px |
| margines górny i dolny | **10 px** | `(48 − 28) / 2` |
| odstęp od lewej krawędzi | **12 px** | 3 × jednostka 4 px = pole ochronne ½X (6,4 px) z zapasem |
| odstęp godło → logotyp | **12 px** | ½ wersalika `DANACO` w stopniu 13 px |
| barwa grotów | `--dn-rama-tekst` (#ECECEC) | — |
| barwa kropki | `--dn-sygnal-400` (#5C8CEC) | kontrast 5,68 : 1 na ramie |

**Wyjątek pola ochronnego [ZALECENIE].** Pole pionowe wynosi 10 px zamiast
wymaganych 6,4 px — jest **większe** niż minimum, więc nie ma naruszenia.
Problem pojawia się dopiero przy godle 32 px w pasku 48 px: pole spadłoby
do 8 px, co nadal spełnia ½X (7,3 px). Godło 36 px w pasku 48 px **narusza**
pole (6 px < 8,25 px) i jest zakazane.

### 14.2 Karta środowiska

| Element | Wielkość | Uwaga |
|---|---:|---|
| emblemat środowiska w liście | 24 px | siatka 24 × 24, obrys 1,75 |
| **emblemat w ekspozycji na karcie** | **40 px** | obrys pozostaje 1,75 |
| tytuł środowiska | `--dn-fs-3xl` (30 px) Space Grotesk | — |
| odstęp emblemat → tytuł | 24 px (`--dn-od-6`) | — |
| kropka sygnału w emblemacie | **dokładnie jedna** | wypełniona, `--dn-kropka` |

Na karcie środowiska **nie umieszcza się godła marki** — kartę oznacza emblemat
środowiska. Godło stoi wyżej, w pasku. Podwójne oznaczenie (godło + emblemat
w tym samym kadrze) jest naruszeniem 16.16.

### 14.3 Okno startowe i okno logowania

| Okno | Konfiguracja | Szerokość | Położenie |
|---|---|---:|---|
| **Okno startowe (ładowania)** | lockup pionowy na ciemnym | **160 px** | wyśrodkowany w obu osiach, przesunięty o −8 % wysokości względem środka geometrycznego |
| **Okno rejestracji i logowania** | lockup pionowy | **140 px** | nad blokiem formularza, odstęp 48 px (`--dn-od-12`) |

Przesunięcie w oknie startowym: znak wyśrodkowany geometrycznie w pełnym oknie
wygląda na osadzony za nisko, bo dolna część kadru jest pusta. Przesunięcie
o 8 % wysokości okna w górę przywraca wrażenie środka. [ZALECENIE]

### 14.4 Always On Display

| Element | Konfiguracja | Wielkość |
|---|---|---:|
| rdzeń awatara | **sygnet pełny** | 96 px |
| kropka sygnału w rdzeniu | pulsuje (`--dn-czas-tetno`) | — |

To **jedyne** miejsce w systemie, gdzie kropka **sygnetu** się porusza — bo
w AOD sygnet nie występuje jako godło marki, tylko jako **rdzeń wskaźnika
procesu**. Rozróżnienie jest formalne: w AOD znak jest komponentem `.dn-aod-rdzen`,
nie elementem tożsamości.

### 14.5 Medalion nadawcy w oknie komunikacji

| Rola nadawcy | Medalion | Zawartość |
|---|---:|---|
| `--inteligencja` | 28 px | **sygnet pełny** na tle `--dn-powierzchnia-2` |
| `--czlowiek` | 28 px | inicjały Operatora, krój mono |
| `--system` | 28 px | ikona z zestawu (`ustawienia`, `terminal`…) |

Sygnet w medalionie ma **28 px kadru**, czyli X = 12,83 px i kreskę 2,59 px —
powyżej progu. Poniżej 24 px medalion przechodzi na wariant uproszczony.

### 14.6 Tabela zbiorcza — wielkość znaku w oknach platformy

| Okno / element | Konfiguracja | Wielkość |
|---|---|---:|
| pasek górny (wszystkie powłoki) | sygnet pełny + logotyp | godło **28 px** |
| medalion nadawcy (Chat Window) | sygnet pełny | **28 px** |
| karta środowiska (Centrum dowodzenia) | emblemat środowiska | **24 / 40 px** |
| kafel komponentu własnego (strefa 2) | ikona z zestawu | **24 px** |
| Okno startowe | lockup pionowy | **160 px** |
| Okno rejestracji i logowania | lockup pionowy | **140 px** |
| Always On Display | sygnet pełny (rdzeń) | **96 px** |
| Okno Konfiguracji — nagłówek | logotyp samodzielny | **90 px** |
| favicon | sygnet uproszczony | **16 / 32 / 48 px** |
| ikona aplikacji | sygnet pełny w kaflu | **512 / 1024 px** |
| ikona maskowalna (Android) | sygnet pełny w kaflu bez zaokrągleń | **512 px** |

---

## 15. Barwy znaku i dozwolone pary tło/znak

### 15.1 Wartości dosłowne dla każdego wariantu [ZAMKNIĘTE]

Odczytane z plików. Kolumna „żeton" wskazuje odpowiednik w `zetony.css` —
**w kodzie interfejsu używa się żetonu, nie wartości szesnastkowej.**

| Plik | Groty | Żeton | Kropka | Żeton |
|---|---|---|---|---|
| `sygnet.svg` | **#181818** | `--dn-szary-900` | **#3B6FE0** | `--dn-sygnal-500` |
| `sygnet-na-ciemnym.svg` | **#F4F4F4** | `--dn-szary-50` | **#5C8CEC** | `--dn-sygnal-400` |
| `sygnet-mono-czarny.svg` | #181818 | `--dn-szary-900` | #181818 | `--dn-szary-900` |
| `sygnet-mono-bialy.svg` | #F4F4F4 | `--dn-szary-50` | #F4F4F4 | `--dn-szary-50` |
| `sygnet-uproszczony.svg` | #181818 | `--dn-szary-900` | #3B6FE0 | `--dn-sygnal-500` |
| `sygnet-uproszczony-na-ciemnym.svg` | #F4F4F4 | `--dn-szary-50` | #5C8CEC | `--dn-sygnal-400` |
| `logo-poziomy.svg` | #181818 | `--dn-szary-900` | #3B6FE0 | `--dn-sygnal-500` |
| `logo-poziomy-na-ciemnym.svg` | #F4F4F4 | `--dn-szary-50` | #5C8CEC | `--dn-sygnal-400` |
| `logo-pionowy.svg` | #181818 | `--dn-szary-900` | #3B6FE0 | `--dn-sygnal-500` |
| `logo-pionowy-na-ciemnym.svg` | #F4F4F4 | `--dn-szary-50` | #5C8CEC | `--dn-sygnal-400` |
| `favicon.svg` | #181818 / **#F4F4F4** (media query) | — | #3B6FE0 / **#5C8CEC** | — |
| `ikona-aplikacji.svg` | #F4F4F4 na kaflu **#131313** | `--dn-rama` | #5C8CEC | `--dn-sygnal-400` |

**Uwaga [ZAMKNIĘTE]:** atrament znaku na ciemnym tle to **`#F4F4F4`
(`szary-50`)**, a nie `#ECECEC` (`szary-100`, żeton tekstu w motywie ciemnym).
Znak jest o jeden stopień jaśniejszy od tekstu — celowo: godło ma stać przed
treścią, nie zlewać się z nią. Wyjątkiem jest **pasek górny**, gdzie godło
przyjmuje `--dn-rama-tekst` (#ECECEC), żeby zrównać się z logotypem w ramie.

### 15.2 Zasada nadrzędna barwy znaku [REGUŁA]

> **Godło zachowuje barwy własne w obu motywach.** Nie przełącza się z motywem
> poza wyborem wariantu „na jasnym / na ciemnym", który zależy od **tła, na
> którym stoi**, a nie od ustawienia motywu.

Konsekwencja praktyczna: w pasku górnym — który jest atramentowy **w obu
motywach** — godło zawsze jest w wariancie na ciemnym tle, także wtedy, gdy
Operator pracuje w motywie jasnym.

### 15.3 Tabela dozwolonych par tło/znak z pomiarem kontrastu

Pomiar według **WCAG 2.1 · współczynnik luminancji względnej**. Dla znaku
graficznego obowiązuje próg **3 : 1** (kryterium 1.4.11 — kontrast elementów
nietekstowych). Kolumna „ocena": ✓ spełnia, ✗ nie spełnia.

#### 15.3.1 Warianty na jasnym tle — groty `#181818`

| Tło | Żeton | Kontrast grotów | Ocena | Kontrast kropki `#3B6FE0` | Ocena |
|---|---|---:|:-:|---:|:-:|
| `#F4F4F4` | `--dn-tlo` (jasny) | **16,14 : 1** | ✓ | **4,21 : 1** | ✓ |
| `#FFFFFF` | `--dn-powierzchnia` (jasny) | **17,76 : 1** | ✓ | **4,63 : 1** | ✓ |
| `#FAFAFA` | `--dn-panel` (jasny) | **17,01 : 1** | ✓ | 4,44 : 1 | ✓ |
| `#ECECEC` | `--dn-powierzchnia-2` (jasny) | **15,03 : 1** | ✓ | **3,92 : 1** | ✓ |
| `#E3E3E3` | `--dn-obrys` (jasny) | 13,74 : 1 | ✓ | 3,58 : 1 | ✓ |
| `#D7D7D7` | `--dn-szary-200` | **12,34 : 1** | ✓ | **3,22 : 1** | ✓ |
| `#C0C0C0` | `--dn-szary-300` | 9,76 : 1 | ✓ | **2,55 : 1** | **✗** |
| `#9E9E9E` | `--dn-szary-400` | 6,63 : 1 | ✓ | 1,73 : 1 | **✗** |
| `#EDF3FE` | `--dn-sygnal-100` | **15,94 : 1** | ✓ | 4,15 : 1 | ✓ |

**Granica [REGUŁA]: `#D7D7D7` (`szary-200`) jest najciemniejszym jasnym tłem
dopuszczalnym dla wariantu z kolorową kropką.** Na tłach ciemniejszych kropka
spada poniżej 3 : 1 — wtedy obowiązuje **wariant mono** albo **wariant na
ciemnym tle**.

#### 15.3.2 Warianty na ciemnym tle — groty `#F4F4F4`

| Tło | Żeton | Kontrast grotów | Ocena | Kontrast kropki `#5C8CEC` | Ocena |
|---|---|---:|:-:|---:|:-:|
| `#0A0A0A` | `--dn-szary-1000` | **18,00 : 1** | ✓ | 6,53 : 1 | ✓ |
| `#0F0F0F` | `--dn-tlo` (ciemny) | **17,43 : 1** | ✓ | **5,86 : 1** | ✓ |
| `#131313` | `--dn-rama` | **16,89 : 1** | ✓ | **5,68 : 1** | ✓ |
| `#181818` | `--dn-powierzchnia` (ciemny) | **16,14 : 1** | ✓ | **5,43 : 1** | ✓ |
| `#212121` | `--dn-panel` (ciemny) | **14,64 : 1** | ✓ | 4,92 : 1 | ✓ |
| `#2A2A2A` | `--dn-szary-800` | **13,05 : 1** | ✓ | **4,39 : 1** | ✓ |
| `#3A3A3A` | `--dn-szary-750` | 10,34 : 1 | ✓ | 3,48 : 1 | ✓ |
| `#4A4A4A` | `--dn-szary-700` | 8,06 : 1 | ✓ | **2,71 : 1** | **✗** |
| `#616161` | `--dn-szary-600` | 5,63 : 1 | ✓ | 1,89 : 1 | **✗** |
| `#7C7C7C` | `--dn-szary-500` | **3,80 : 1** | ✓ | 1,28 : 1 | **✗** |
| `#173A8C` | `--dn-sygnal-800` | 9,45 : 1 | ✓ | — (mono) | — |

**Granica [REGUŁA]: `#3A3A3A` (`szary-750`) jest najjaśniejszym ciemnym tłem
dopuszczalnym dla wariantu z kolorową kropką.** Powyżej — wariant mono biały.

#### 15.3.3 Znak na tle sygnałowym

| Tło | Wariant znaku | Kontrast | Ocena | Uwaga |
|---|---|---:|:-:|---|
| `#3B6FE0` (`sygnal-500`) | mono biały `#F4F4F4` | 3,49 : 1 | ✓ | dopuszczalne wyłącznie w ekspozycji ≥ 48 px |
| `#3B6FE0` | mono czarny `#181818` | **3,84 : 1** | ✓ | preferowany na sygnale |
| `#3B6FE0` | pełny (kropka `#3B6FE0`) | — | **✗** | **kropka znika w tle — zakazane** |
| `#2457C9` (`sygnal-600`) | mono biały `#FFFFFF` | 6,40 : 1 | ✓ | — |
| `#173A8C` (`sygnal-800`) | na ciemnym `#F4F4F4` | 9,45 : 1 | ✓ | kropka `#5C8CEC` → 1,80 : 1 ✗ → **użyj mono białego** |
| `#EDF3FE` (`sygnal-100`) | pełny na jasnym | 15,94 / 4,15 | ✓ | — |

**Reguła bezwzględna [REGUŁA]: znaku w wariancie pełnym nie umieszcza się na
żadnym tle z rodziny sygnału.** Kropka sygnału musi być odróżnialna — na tle
sygnałowym nie jest.

### 15.4 Znak na zdjęciu i na tle niejednolitym [REGUŁA]

| Sytuacja | Rozwiązanie |
|---|---|
| zdjęcie jasne, jednolite w polu ochronnym | wariant mono czarny |
| zdjęcie ciemne, jednolite w polu ochronnym | wariant mono biały |
| zdjęcie o zmiennej jasności | **podkład kryjący** pod znak: prostokąt `--dn-szary-925` z zaokrągleniem `--dn-r-lg`, marginesy = pole ochronne 1X |
| gradient | jak wyżej — podkład kryjący; **gradient nigdy nie jest tłem znaku** |
| wideo | zawsze podkład kryjący |

Zakazane: półprzezroczysty podkład, rozmycie tła pod znakiem (glassmorfizm),
cień rzucany przez znak, obrys („outline") wokół znaku.

---

## 16. Katalog błędnych użyć

Dwadzieścia przypadków. Każdy opisany dwoma zdaniami: **co zrobiono źle**
i **dlaczego to szkodzi**. Wszystkie mają żywe odpowiedniki w prototypie
`konstrukcja-znaku.html` (rozdz. „Katalog naruszeń").

### 16.01 Rozciągnięcie w poziomie

- **Co zrobiono źle:** znak przeskalowano `scale(1.2, 1)`, żeby wypełnił
  szerszy kontener.
- **Dlaczego to szkodzi:** kąt rozwarcia grotu spada z 95,45° do 85,2° —
  grot promptu zamienia się w strzałkę nawigacji. Znak przestaje znaczyć
  „przekazanie", zaczyna znaczyć „dalej".

### 16.02 Ściśnięcie w pionie

- **Co zrobiono źle:** znak przeskalowano `scale(1, 0.8)`, żeby zmieścić się
  w niskim pasku.
- **Dlaczego to szkodzi:** kropka przestaje być okręgiem i staje się elipsą —
  a kropka jest **jedynym** elementem sygnaturowym całego systemu. Zdeformowana
  kropka podważa spójność wszystkich pozostałych wystąpień kropki w interfejsie.

### 16.03 Obrót

- **Co zrobiono źle:** znak obrócono o 45°, żeby wpisać go w narożnik.
- **Dlaczego to szkodzi:** znak ma kierunek czytania (od koordynatora do
  wykonawcy, od lewej do prawej). Obrót usuwa kierunek i zamienia znak
  w ornament.

### 16.04 Odbicie lustrzane

- **Co zrobiono źle:** znak odbito poziomo (`scale(-1, 1)`), żeby „patrzył
  w stronę treści".
- **Dlaczego to szkodzi:** groty wskazują wstecz, a kropka ląduje **przed**
  nimi. Znak czyta się `.««` — czyli „praca cofnięta przed zleceniem".
  Znaczenie jest odwrotne do zamierzonego.

### 16.05 Wyrównanie sygnetu do linii bazowej DANACO

- **Co zrobiono źle:** w lockupie poziomym sygnet postawiono dolną krawędzią
  na linii bazowej `DANACO` (`y = 51`).
- **Dlaczego to szkodzi:** sygnet podskakuje o 11,7 j. w górę i traci związek
  z drugim wierszem. Lockup rozpada się na dwa niezależne obiekty stojące obok
  siebie.

### 16.06 Wyrównanie sygnetu do środka geometrycznego bloku tekstu

- **Co zrobiono źle:** oś grotów ustawiono na `(26,63 + 67,25) / 2 = 46,94`
  zamiast na 48.
- **Dlaczego to szkodzi:** różnica jest mała (1,06 j.), ale znosi celową
  korektę optyczną: blok tekstu ma masę skupioną w górnym wierszu, więc sygnet
  wyrównany do środka geometrycznego wygląda na zawieszony za wysoko.

### 16.07 Wyrównanie sygnetu obwiednią farby zamiast kadrem

- **Co zrobiono źle:** w lockupie pionowym sygnet wyśrodkowano po obwiedni
  farby (`43,63 … 121,25`) zamiast po kadrze (`31,7 … 127,7`).
- **Dlaczego to szkodzi:** sygnet przesuwa się o 2,75 j. w lewo względem
  `DANACO`. Cały układ pionowy czyta się jako przechylony, mimo że każdy
  element z osobna jest wyśrodkowany.

### 16.08 Zmiana odstępu sygnet ↔ logotyp

- **Co zrobiono źle:** odstęp powiększono do 24 j., „żeby oddychało".
- **Dlaczego to szkodzi:** przy odstępie równym wysokości wersalika sygnet
  i logotyp przestają być jednym znakiem — czytają się jako godło **i** nazwa,
  czyli dwa oznaczenia. Reguła ½ wersalika utrzymuje je w jednej frazie.

### 16.09 Kropka wyśrodkowana na osi grotów

- **Co zrobiono źle:** kropkę przesunięto z `y = 63,5` na `y = 48`, „bo tak
  jest równiej".
- **Dlaczego to szkodzi:** znak przestaje być zdaniem `»».` i staje się
  schematem blokowym z trzema równorzędnymi elementami. Ginie hierarchia
  „dwa człony + zamknięcie".

### 16.10 Kropka w innej barwie

- **Co zrobiono źle:** kropkę pokolorowano na zielono, żeby oznaczyć „stan
  poprawny".
- **Dlaczego to szkodzi:** kropka sygnału jest **jedyną barwą marki** i znaczy
  dokładnie jedno: „tu biegnie praca". Barwy stanów (zieleń, bursztyn,
  czerwień) należą do komponentów interfejsu, nie do znaku. Kolorowa kropka
  w godle zamienia tożsamość we wskaźnik stanu.

### 16.11 Kropka usunięta

- **Co zrobiono źle:** z sygnetu usunięto kropkę, zostawiając same groty,
  „bo w mono i tak jest czarna".
- **Dlaczego to szkodzi:** bez kropki znak czyta się `»»` — otwarty cudzysłów
  albo znak przewijania. Znika cała treść: praca, która ruszyła. W wariancie
  mono kropka **pozostaje**, tylko przyjmuje barwę znaku.

### 16.12 Trzeci grot

- **Co zrobiono źle:** dodano trzeci grot, „żeby znak był bogatszy".
- **Dlaczego to szkodzi:** dwa groty to relacja **koordynator → wykonawca**,
  zapisana w architekturze platformy. Trzy groty to ciąg bez końca — i kolizja
  ze zwyczajowym znaczeniem `»»»` jako „przewiń do końca".

### 16.13 Kadrowanie po obwiedni farby

- **Co zrobiono źle:** kadr pliku przycięto do `12 … 89,5 × 26 … 70`,
  „żeby nie było pustego miejsca".
- **Dlaczego to szkodzi:** znika korekta optyczna zapisana w marginesach kadru,
  a proporcja zmienia się z 1 : 1 na 1,761 : 1. Wszystkie kontenery kwadratowe
  (medalion, favicon, kafel) zaczynają deformować znak.

### 16.14 Znak wpisany w koło

- **Co zrobiono źle:** sygnet wpisano w okrągły awatar bez przeliczenia
  wielkości.
- **Dlaczego to szkodzi:** kwadratowy kadr wpisany w koło traci narożniki —
  a w prawym dolnym narożniku leży kropka. Kropka zostaje przycięta albo
  wypada poza okrąg.

### 16.15 Ramka wokół znaku

- **Co zrobiono źle:** znak obwiedziono cienką ramką, żeby „oddzielić go od
  tła".
- **Dlaczego to szkodzi:** ramka wchodzi w pole ochronne (½X) i konkuruje
  z grotami o tę samą uwagę — bo grot jest zbudowany z prostych kresek.
  Jeżeli znak nie odcina się od tła, rozwiązaniem jest podkład kryjący,
  nie ramka.

### 16.16 Godło i emblemat środowiska w jednym kadrze

- **Co zrobiono źle:** na karcie środowiska umieszczono jednocześnie godło
  marki i emblemat `TalkIn`.
- **Dlaczego to szkodzi:** dwa oznaczenia w jednym kadrze znoszą hierarchię
  Środowisko → Moduł → Okno. Karta ma mówić „to jest TalkIn", a nie
  „to jest Danaco, a w środku TalkIn" — bo pierwsze zdanie już padło w pasku.

### 16.17 Znak na tle sygnałowym w wariancie pełnym

- **Co zrobiono źle:** godło z kolorową kropką postawiono na przycisku
  wypełnionym `--dn-sygnal-500`.
- **Dlaczego to szkodzi:** kropka jest w tej samej barwie co tło — kontrast
  1,00 : 1. Kropka znika, znak czyta się `»»` i traci znaczenie (patrz 16.11).

### 16.18 Znak na tle o zbyt niskim kontraście

- **Co zrobiono źle:** wariant pełny na tle `--dn-szary-300` (#C0C0C0).
- **Dlaczego to szkodzi:** groty mają 9,76 : 1 (dobrze), ale kropka **2,55 : 1**
  — poniżej progu 3 : 1 wymaganego dla elementów nietekstowych (WCAG 1.4.11).
  Znak jest niedostępny dla osób ze słabym widzeniem kontrastu.

### 16.19 Cień pod znakiem

- **Co zrobiono źle:** pod godło podłożono `--dn-cien-2`, „żeby odskoczyło
  od tła".
- **Dlaczego to szkodzi:** cień jest w systemie zarezerwowany dla **warstw**
  (modal, dymek, toast) i niesie informację „to jest wyżej". Znak nie jest
  warstwą. Dodatkowo cień rozmywa krawędzie grotów przy małych wielkościach.

### 16.20 Wariant pełny poniżej 24 px

- **Co zrobiono źle:** favicon 16 px wygenerowano z pliku `sygnet.svg`
  zamiast `sygnet-uproszczony.svg`.
- **Dlaczego to szkodzi:** ramię ma 1,48 px, kropka 2,17 px — oba poniżej
  progu renderowania. Znak zamienia się w szarą plamę bez rozpoznawalnego
  kształtu, a wygładzanie krawędzi zamyka prześwit 2,67 px.

### 16.21 Tętno kropki w godle

- **Co zrobiono źle:** kropce w godle w pasku górnym dodano animację tętna.
- **Dlaczego to szkodzi:** tętno znaczy „tu biegnie praca". Pulsujące godło
  twierdzi, że pracuje **marka**, co jest bez sensu i odbiera tętnu wartość
  informacyjną w miejscach, gdzie ono coś znaczy (karta sesji, wpis rozmowy).

### 16.22 Rekonstrukcja znaku z kroju pisma

- **Co zrobiono źle:** znak złożono ze znaków `»` i `.` w Space Grotesk,
  bo „tak samo wygląda".
- **Dlaczego to szkodzi:** cudzysłów francuski w każdym kroju ma inny kąt,
  inną grubość, inny prześwit i inne położenie względem linii bazowej. Znak
  przestaje być powtarzalny — a powtarzalność jest jedyną funkcją znaku.

---

## 17. Reprodukcja

### 17.1 Druk pełnokolorowy

| Parametr | Wartość |
|---|---|
| przestrzeń źródłowa | **sRGB** — `#181818` atrament, `#3B6FE0` kropka |
| konwersja do CMYK | **wyłącznie z profilu ICC nośnika** dostarczonego przez drukarnię (renderowanie perceptual dla zdjęć, **relative colorimetric dla znaku**) |
| nadruk (overprint) | **wyłączony** dla obu elementów |
| zalewkowanie (trapping) | **niepotrzebne** — elementy się nie stykają (prześwit 16 j.) |
| minimalna liniatura | 150 lpi dla wielkości ≥ 9 mm |

**Nie podajemy rozpisu CMYK w tym dokumencie [REGUŁA].** Rozpis zależy od
profilu papieru (powlekany / niepowlekany / offsetowy), a rozpis „uniwersalny"
daje na niepowlekanym papierze błękit przesunięty w fiolet — czyli dokładnie
to, czego kierunek projektowy zakazuje. Rozpis ustala się **na próbie
kontrolnej**.

### 17.2 Pantone

> **Numer Pantone: do ustalenia z drukarnią na wydruku kontrolnym.**

Nie podajemy numeru, bo:

1. dopasowanie `#3B6FE0` do palety Pantone zależy od tego, czy nośnik jest
   powlekany (C), niepowlekany (U) czy matowy (M) — to trzy różne numery,
2. wybór między „najbliższy numerycznie" a „najbliższy wizualnie" wymaga
   oceny na próbce, nie w pliku,
3. podanie numeru bez próby jest zobowiązaniem, którego nie da się dotrzymać.

**Procedura [REGUŁA]:**

1. przekaż drukarni plik SVG i wartość `#3B6FE0` w sRGB,
2. poproś o **trzy wydruki próbne** na docelowym nośniku (kandydaci z wachlarza),
3. wybierz na świetle D50,
4. **zapisz wybrany numer w tym dokumencie** — do tego czasu pozostaje
   „do ustalenia z drukarnią".

Atrament `#181818` w druku jednokolorowym: **czerń procesowa 100 %**, nie
czerń kompozytowa (bogata czerń rozjeżdża się przy małych elementach).

### 17.3 Grawer laserowy i mechaniczny

| Parametr | Wartość |
|---|---|
| **wariant** | **konturowy** — groty i kropka jako obrysy, nie wypełnienia |
| minimalna wysokość znaku | **15 mm** |
| minimalna szerokość rowka | **0,3 mm** |
| kropka | **okrąg konturowy**, nie plama — wypełnienie zlewa się z tłem materiału |
| prześwit między grotami przy 15 mm | 2,50 mm — bezpieczny |
| kierunek grawerowania | zgodny z osią ramienia (redukuje schodkowanie krawędzi) |

Plik roboczy: `sygnet.svg` z ręcznie zamienionym `fill` na
`fill="none" stroke="…" stroke-width="2"` — patrz rysunek
`rys-14-warianty-reprodukcyjne.svg`, kolumna „KONTUR / GRAWER".

### 17.4 Tłoczenie (suchy tłok, hot stamping)

| Parametr | Wartość |
|---|---|
| minimalna wysokość znaku | **15 mm** |
| minimalna szerokość elementu | **0,8 mm** przy wysokości 15 mm — spełnione (kreska 1,39 mm) |
| kropka | pełna, wypukła — jedyny element, który w tłoczeniu wolno pozostawić jako plamę |
| głębokość | 0,3–0,5 mm (do ustalenia z wykonawcą matrycy) |
| skos krawędzi | 45° — zgodny z nachyleniem ramienia (47,73°), więc wizualnie neutralny |

### 17.5 Znak jednokolorowy — procedura [REGUŁA]

Przy każdej reprodukcji jednokolorową farbą (sitodruk jednokolorowy, stempel,
grawer, tłoczenie, odlew) obowiązuje:

1. użyj pliku `sygnet-mono-czarny.svg` albo `sygnet-mono-bialy.svg`,
2. **dla lockupów** — użyj `logo-poziomy-mono-*.svg`, ale **sprawdź kropkę
   po `CONSOLE`**: w plikach mono pozostaje ona w barwie `#3B6FE0` i wymaga
   ręcznego przestawienia na kolor znaku,
3. **nie usuwaj kropki** — zmienia się jej barwa, nie jej istnienie
   (naruszenie 16.11),
4. sprawdź prześwit: przy wysokości znaku `H` prześwit wynosi `0,3636 × H` —
   musi być większy od dwukrotności rozlewu farby.

### 17.6 Haft

| Parametr | Wartość |
|---|---|
| **wariant** | **wyłącznie sygnet uproszczony** |
| minimalna wysokość | **40 mm** (kadr) → grot 25 mm |
| ścieg grotu | satynowy, kierunek prostopadły do osi ramienia |
| ścieg kropki | wypełniający, spiralny od środka |
| podkład | stabilizator wycinany — grot ma długie proste krawędzie, które bez podkładu falują |
| minimalna szerokość elementu | 1,2 mm → przy 40 mm kadru kreska ma 5,09 mm ✓ |

Wariant pełny **nie nadaje się do haftu w żadnym rozmiarze użytkowym**:
prześwit 16 j. przy kadrze 40 mm daje 6,67 mm, ale dwa równoległe groty
w ściegu satynowym ściągają materiał i prześwit się zamyka.

### 17.7 Odlew, frezowanie, wycinanie laserowe

| Technika | Minimalna wysokość | Uwaga |
|---|---:|---|
| odlew (mosiądz, aluminium) | **25 mm** | kropka jako element osobny — nie jest połączona z grotami |
| frezowanie CNC | **20 mm** | promień freza ≤ 1 mm; narożniki znaku są ostre i wymagają dobiegu |
| wycinanie laserowe w blasze | **30 mm** | **kropka wypadnie** — wymaga mostka albo osobnego elementu |
| wycinanie w folii (ploter) | **25 mm** | kropka jako osobny element aplikacji |

**Ostrzeżenie konstrukcyjne [REGUŁA]:** we wszystkich technikach wycinania
kropka jest **elementem niepołączonym** z resztą znaku. Musi być planowana
jako osobny detal — nigdy nie wolno łączyć jej z grotem mostkiem, bo mostek
staje się trzecim elementem znaku.

### 17.8 Tablica minimalnych wielkości według techniki

| Technika | Sygnet pełny | Sygnet uproszczony | Lockup |
|---|---:|---:|---:|
| druk offsetowy / cyfrowy | 9 mm | 6 mm | 28 mm |
| sitodruk, flexo, opakowanie | 12 mm | 8 mm | 36 mm |
| stempel | — | 10 mm | — |
| grawer laserowy | 15 mm | 10 mm | 40 mm |
| tłoczenie | 15 mm | 10 mm | 40 mm |
| haft | — | 40 mm | — |
| odlew | 25 mm | 18 mm | — |
| wycinanie laserowe | 30 mm | 20 mm | — |

---

## 18. Lista kontrolna poprawnego użycia znaku

Do przejścia **przed** oddaniem każdego materiału ze znakiem.

### 18.1 Geometria

- [ ] znak pochodzi z pliku źródłowego, nie został odrysowany ani złożony z kroju pisma
- [ ] skalowanie proporcjonalne — ten sam współczynnik w obu osiach
- [ ] brak obrotu, pochylenia, odbicia
- [ ] kąt rozwarcia grotu wynosi 95,45° (kontrola: ostrza `C` i `C′` leżą na osi poziomej)
- [ ] prześwit między grotami jest stały na całej wysokości
- [ ] kropka stoi na linii bazowej grotów
- [ ] kropka jest okręgiem, nie elipsą

### 18.2 Konfiguracja i wielkość

- [ ] wybrana konfiguracja jest właściwa dla nośnika (rozdz. 11.4)
- [ ] wielkość nie jest mniejsza od minimalnej dla tej konfiguracji i tej techniki
- [ ] poniżej 24 px użyto wariantu uproszczonego
- [ ] w interfejsie wielkość jest wielokrotnością 4 px
- [ ] lockup nie został rozbity na sygnet i logotyp umieszczone osobno

### 18.3 Pole ochronne

- [ ] pole ≥ ½X ze wszystkich stron, liczone **od kadru**
- [ ] pole ekspozycyjne 1X przy samodzielnym wystąpieniu znaku
- [ ] w polu nie ma tekstu, linii, ramek, innych znaków ani krawędzi nośnika
- [ ] tło w polu jest jednolite

### 18.4 Barwa i kontrast

- [ ] wariant znaku dobrany do **tła**, nie do ustawienia motywu
- [ ] kontrast grotów wobec tła ≥ 3 : 1
- [ ] **kontrast kropki wobec tła ≥ 3 : 1**
- [ ] znak nie stoi na żadnym tle z rodziny sygnału w wariancie pełnym
- [ ] w wariancie mono kropka jest obecna i ma barwę znaku
- [ ] w kodzie użyto żetonów `--dn-*`, nie wartości szesnastkowych

### 18.5 Kontekst

- [ ] znak nie ma cienia, obrysu ani podkładu półprzezroczystego
- [ ] kropka w godle nie jest animowana
- [ ] w jednym kadrze stoi jedno oznaczenie (godło **albo** emblemat środowiska)
- [ ] przy tle niejednolitym zastosowano podkład kryjący
- [ ] `role="img"` i `aria-label="Danaco Console"` przy wstawieniu inline

### 18.6 Pliki

- [ ] przekazano SVG, nie rasteryzację, wszędzie tam, gdzie SVG jest przyjmowany
- [ ] rasteryzacje mają co najmniej dwukrotność wielkości docelowej
- [ ] pliki mono sprawdzone pod kątem kropki po `CONSOLE` (rozdz. 17.5)
- [ ] favicon i ikona aplikacji pochodzą z katalogów `favicon/` i `ikona-aplikacji/`,
      nie z ręcznego przeskalowania sygnetu

---

## 19. Wykaz rysunków i plików

### 19.1 Rysunki konstrukcyjne

| Rysunek | Plik | Zawartość |
|---|---|---|
| 01 | `konstrukcja/rys-01-siatka-bazowa.svg` | siatka 4/12 j., osie, linia bazowa, punkty A–F i A′–F′, kropka S |
| 02 | `konstrukcja/rys-02-wymiarowanie.svg` | linie wymiarowe, wszystkie wymiary w j. i w module X |
| 03 | `konstrukcja/rys-03-katy-ramienia.svg` | kąt 95,45°, grubość pozioma 12 j., grubość prostopadła 8,879 j. |
| 04 | `konstrukcja/rys-04-kropka-linia-bazowa.svg` | kropka, linia bazowa, trzy środki kadru, udział w polu farby |
| 05 | `konstrukcja/rys-05-wariant-uproszczony.svg` | porównanie wariantów, tabela parametrów |
| 06 | `konstrukcja/rys-06-lockup-poziomy.svg` | konstrukcja lockupu poziomego, linie bazowe, odstępy |
| 07 | `konstrukcja/rys-07-lockup-pionowy.svg` | konstrukcja lockupu pionowego, oś symetrii, centrowanie |
| 08 | `konstrukcja/rys-08-pole-ochronne-sygnet.svg` | pole ½X i 1X dla sygnetu |
| 09 | `konstrukcja/rys-09-pole-ochronne-lockup-poziomy.svg` | pole ½X i 1X dla lockupu poziomego |
| 10 | `konstrukcja/rys-10-pole-ochronne-lockup-pionowy.svg` | pole ½X i 1X dla lockupu pionowego |
| 11 | `konstrukcja/rys-11-drabina-wielkosci.svg` | drabina 12–128 px, próg 24 px, progi renderowania |
| 12 | `konstrukcja/rys-12-znak-w-interfejsie.svg` | pasek 48 px, karta środowiska, okno startowe, tabela wielkości |
| 13 | `konstrukcja/rys-13-emblemat-siatka.svg` | cztery emblematy na siatce 24 × 24, obrys 1,75 a skala |
| 14 | `konstrukcja/rys-14-warianty-reprodukcyjne.svg` | pełny, mono, na ciemnym, kontur do grawerowania |

Rasteryzacje: `konstrukcja/png/rys-NN-*.png` (1600 px szerokości, tło białe).

### 19.2 Pliki znaku — kompletny wykaz

| Kategoria | Pliki |
|---|---|
| **sygnet** | `sygnet.svg` · `sygnet-na-ciemnym.svg` · `sygnet-mono-czarny.svg` · `sygnet-mono-bialy.svg` |
| **sygnet uproszczony** | `sygnet-uproszczony.svg` · `sygnet-uproszczony-na-ciemnym.svg` |
| **lockup poziomy** | `logo-poziomy.svg` · `logo-poziomy-na-ciemnym.svg` · `logo-poziomy-mono-czarny.svg` · `logo-poziomy-mono-bialy.svg` |
| **lockup pionowy** | `logo-pionowy.svg` · `logo-pionowy-na-ciemnym.svg` |
| **logotyp** | `logotyp.svg` · `logotyp-na-ciemnym.svg` |
| **rasteryzacje** | `png/sygnet-512.png` · `png/sygnet-na-ciemnym-512.png` · `png/logo-poziomy@2x.png` · `png/logo-poziomy-na-ciemnym@2x.png` · `png/logo-pionowy@2x.png` · `png/logotyp@2x.png` |
| **favicon** | `favicon.svg` · `favicon.ico` · `favicon-16/32/48.png` · `apple-touch-icon.png` · `icon-192.png` · `icon-512.png` · `site.webmanifest` |
| **ikona aplikacji** | `ikona-aplikacji.svg` · `ikona-maskowalna.svg` · `ikona-180/192/512/1024.png` · `ikona-maskowalna-512.png` |
| **emblematy środowisk** | `srodowisko-talkin.svg` · `srodowisko-workspace.svg` · `srodowisko-codestudio.svg` · `srodowisko-multitaskingai.svg` |

Wszystkie w `WYNIK/zasoby/marka/`.

### 19.3 Szybka tablica liczb

Do wydruku i powieszenia obok stanowiska:

```
   MODUŁ            X = 44 j. = wysokość grotu
   PROPORCJE        kreska 0,2018 X · kropka 0,2955 X · prześwit 0,3636 X
   KĄT              95,45°  [ZAMKNIĘTE]
   LINIA BAZOWA     y = 70 — wspólna dla grotów i kropki
   POLE OCHRONNE    minimum ½X · ekspozycja 1X · liczone od KADRU
   PRÓG WARIANTU    24 px  ·  dolna granica systemu 16 px
   MIN. DRUK        sygnet 9 mm · lockup 28 mm · grawer 15 mm · haft 40 mm
   KONTRAST         groty ≥ 3:1 · KROPKA ≥ 3:1 (to jest ten trudny)
   BARWY            #181818 / #3B6FE0   ·   #F4F4F4 / #5C8CEC
```

---

## 20. Decyzje projektowe

### Moduł X przywiązany do farby, nie do kadru

**Kontekst.** Modułem mogła być wysokość kadru (96 j.), połowa kadru (48 j.)
albo wysokość grotu (44 j.).

**Rozstrzygnięcie.** Modułem jest **X = wysokość grotu = 44 j.**

**Uzasadnienie.** Kadr jest bytem pliku i znika w chwili, gdy znak zostaje
wycięty, wydrukowany, wygrawerowany albo sfotografowany. Wysokość grotu da się
zmierzyć zawsze i wszędzie — na ekranie linijką ekranową, na wydruku suwmiarką,
na szyldzie taśmą. Reguła oparta na kadrze byłaby niemierzalna w połowie
przypadków, w których jest potrzebna.

**Konsekwencja.** Wszystkie proporcje w tym dokumencie podano dwukrotnie:
w jednostkach siatki (do rysowania) i w module X (do kontroli w terenie).

### Dwa poziomy pola ochronnego: ½X i 1X

**Kontekst.** `ksiega-znaku.md` rozdz. 5 ustala pole ochronne **½X** jako
[REGUŁA]. Brief dokumentacji konstrukcyjnej wskazuje **1X** jako wartość
docelową dla ekspozycji.

**Rozstrzygnięcie.** Wprowadzono **dwa poziomy**:
**½X — minimum obowiązujące zawsze** (wartość wiążąca, zgodna z Księgą Znaku),
**1X — pole ekspozycyjne, zalecane** przy samodzielnym wystąpieniu znaku.

**Uzasadnienie.** Jedna wartość nie obsługuje obu skrajności. W kokpicie
(pasek 48 px, medalion 28 px, kafel) pole 1X jest fizycznie niewykonalne —
godło 28 px wymagałoby 12,8 px marginesu z każdej strony, czyli 53,7 px
w pasku o wysokości 48 px. Na okładce, tabliczce i koszulce pole ½X wygląda
na przycięcie. Podział na „minimum" i „ekspozycja" rozwiązuje konflikt bez
osłabiania reguły: **wartość wiążąca pozostaje ½X**, 1X jest zaleceniem.

**Konsekwencja.** Kalkulator pola w prototypie HTML podaje obie wartości
jednocześnie i oznacza, która jest obowiązkowa.

### Grubość ramienia liczona prostopadle, nie poziomo

**Kontekst.** Ścieżka znaku daje grubość poziomą 12 j. To liczba, którą widać
w pliku i którą naturalnie się cytuje.

**Rozstrzygnięcie.** Wszystkie progi technologiczne liczy się od **grubości
prostopadłej 8,879 j. = 0,2018 X**.

**Uzasadnienie.** Farba, wiązka lasera i piksel „widzą" kreskę prostopadle do
jej osi, nie poziomo. Liczenie od 12 j. zawyża rzeczywistą grubość o 35 %
i prowadzi do zamówień, w których kreska teoretycznie mieści się w progu, a
w praktyce zalewa się przy pierwszym wydruku.

**Konsekwencja.** Tablica przeliczeniowa w rozdz. 4.3 i wszystkie progi
w rozdz. 11 i 17 są policzone od 0,2018 X.

### Wyrównanie w kontenerze przez kadr, nie przez farbę

**Kontekst.** Kadr 96 × 96 ma niesymetryczne marginesy poziome (12 j. z lewej,
6,5 j. z prawej). To wygląda na błąd i kusi, żeby „naprawić" kadrowanie.

**Rozstrzygnięcie.** **Kadr jest narzędziem wyrównania optycznego i pozostaje
niezmieniony.** Znak wyrównuje się w kontenerze zawsze przez kadr.

**Uzasadnienie.** Środek masy znaku leży na `x = 46,578`, środek prostokąta
farby na `x = 50,750`, środek kadru na `x = 48,000` — czyli w 34 % drogi od
masy do prostokąta. Wyrównanie do prostokąta zsuwa groty (88,8 % masy) w lewo;
wyrównanie do masy dosuwa kropkę do prawej krawędzi na 4,08 j. Kadr jest
zapisanym kompromisem.

**Konsekwencja.** Reguła „kadrowanie po obwiedni farby" jest naruszeniem
(16.13), a nie optymalizacją.

### Wariant uproszczony ma inny kąt niż wariant pełny

**Kontekst.** Naturalne oczekiwanie: wariant uproszczony to jeden grot
wycięty z wariantu pełnego i powiększony.

**Rozstrzygnięcie.** Wariant uproszczony jest **narysowany osobno**:
kąt 91,94° zamiast 95,45°, kreska 0,2037 X zamiast 0,2018 X, kropka 0,300 X
zamiast 0,2955 X. Wartości odczytane z pliku — nie zmieniamy ich.

**Uzasadnienie.** Przy jednym grocie kąt 95,45° traci kierunek, bo nie ma
powtórzenia, które ten kierunek niosło. Zaostrzenie o 3,51° przywraca
czytelność gestu. Pozostałe różnice (poniżej 2,5 % w module X) są korektami
optycznymi na małe wielkości.

**Konsekwencja.** Nie wolno generować wariantu uproszczonego przez wycięcie
grotu z `sygnet.svg`. Obowiązuje plik `sygnet-uproszczony.svg`.

### Kropka po `CONSOLE` centrowana na osi wersalika, nie na linii bazowej

**Kontekst.** W sygnecie kropka stoi na linii bazowej. W logotypie leży na
połowie wysokości wersalika `CONSOLE`. To wygląda na niekonsekwencję.

**Rozstrzygnięcie.** Obie pozycje są poprawne i zamierzone.

**Uzasadnienie.** W sygnecie kropka zamyka **zdanie zapisane grotami** — a
w zdaniu kropka stoi na linii pisma. W logotypie kropka stoi w **wierszu
wersalikowym**, którego oś optyczna to połowa wysokości wersalika, nie linia
bazowa. Kropka na linii bazowej w wierszu wersalikowym czytałaby się jako
zawieszona pod tekstem.

**Konsekwencja.** Przy odtwarzaniu logotypu kropkę centruje się na
`cy = baseline − cap/2`, nie stawia na baseline.

### Wyraz `CONSOLE` centrowany bez kropki

**Kontekst.** W lockupie pionowym wszystkie elementy są centrowane. Kropka
po `CONSOLE` zwiększa szerokość wiersza o 11,4 j.

**Rozstrzygnięcie.** `CONSOLE` centruje się **bez kropki**; kropka zwisa
poza blok.

**Uzasadnienie.** Kropka jest znakiem interpunkcyjnym, a w składzie
typograficznym interpunkcja zwisa poza blok wyrównania (hanging punctuation).
Włączenie kropki do centrowania przesunęłoby wyraz o 5,7 j. w lewo względem
`DANACO` — i drugi wiersz wyraźnie „uciekłby", mimo że matematycznie byłby
wyśrodkowany.

**Konsekwencja.** Wartość odczytana z pliku (`CONSOLE` bez kropki: środek
79,56 wobec osi kadru 79,5) potwierdza rozstrzygnięcie.

### Obrys emblematów nie skaluje się liniowo powyżej 48 px

**Kontekst.** Emblematy środowisk mają `stroke-width="1.75"` w siatce 24 × 24.
Przy renderowaniu 96 px daje to kreskę 7 px, która wygląda znacznie grubiej
niż 1,75 px przy 24 px.

**Rozstrzygnięcie.** Do **48 px włącznie** obrys pozostaje **1,75** bez wyjątku.
Powyżej 48 px, wyłącznie w ekspozycji, wolno obniżyć do **1,5** (64–96 px)
i **1,25** (powyżej 96 px).

**Uzasadnienie.** Do 48 px liniowe skalowanie obrysu jest wizualnie neutralne
i utrzymuje spójność z zestawem ikon interfejsu. Powyżej tej wielkości oko
zaczyna czytać kreskę jako „grubą" — bo w dużych formatach oczekuje kreski
cieńszej względem pola. Obniżenie zachowuje **wrażenie** tej samej wagi.

**Konsekwencja.** Zakaz `vector-effect="non-scaling-stroke"` — ten atrybut
rozwiązuje odwrotny problem i psuje wygląd w ekspozycji.

### Nie podajemy numeru Pantone

**Kontekst.** Księgi znaku zwykle podają numer Pantone jako wartość wiążącą.

**Rozstrzygnięcie.** Numer Pantone: **do ustalenia z drukarnią na wydruku
kontrolnym.** Podajemy wyłącznie procedurę ustalenia (rozdz. 17.2).

**Uzasadnienie.** Dopasowanie `#3B6FE0` różni się dla papieru powlekanego (C),
niepowlekanego (U) i matowego (M) — to trzy różne numery. Wybór wymaga oceny
wzrokowej na próbce w świetle D50. Podanie jednego numeru „na sucho" jest
zobowiązaniem, którego nie da się dotrzymać na wszystkich nośnikach, i naraża
markę na to, że najczęściej reprodukowana barwa będzie za każdym razem inna.

**Konsekwencja.** Do czasu przeprowadzenia próby kontrolnej materiały spotowe
zamawia się z podaniem sRGB i profilu ICC, a nie numeru katalogowego.

### Kropka jako element niepołączony we wszystkich technikach wycinania

**Kontekst.** W wycinaniu laserowym, plotterowym i w odlewie kropka jest
osobnym, niepodpartym elementem — wypada z materiału.

**Rozstrzygnięcie.** Kropka planowana jest jako **osobny detal**. Zakaz
łączenia jej mostkiem z grotem.

**Uzasadnienie.** Mostek staje się **trzecim elementem znaku** — widocznym,
niezamierzonym i niosącym fałszywe znaczenie („praca połączona ze zleceniem",
zamiast „praca, która ruszyła"). Prześwit 4,5 j. jest częścią znaku, tak samo
jak prześwit 16 j. między grotami.

**Konsekwencja.** W technikach, w których osobny detal jest niewykonalny
(np. wycinanie w cienkiej blasze bez montażu), znak **nie jest reprodukowany
w wariancie pełnym** — obowiązuje logotyp samodzielny albo rezygnacja ze znaku
na tym nośniku.

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
