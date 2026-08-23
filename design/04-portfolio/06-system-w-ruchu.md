# Danaco Console — Plansza portfolio P6: System w ruchu

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
| **Odbiorcy** | Właściciel · Designer prowadzący makiety · Deweloper wdrażający · odbiorca portfolio |
| **Zakres** | Kontrakt ruchu `INTENSYWNOSC_RUCHU` 3/10 · cztery czasy i jedna krzywa · katalog siedmiu dozwolonych animacji · katalog dziesięciu zakazanych · `prefers-reduced-motion` · ruch jako nośnik informacji · wydajność · oś czasu złożonego przejścia |
| **Czego NIE zawiera** | definicji żetonów (→ `01-dokumentacja-md/04-tokens.md`), anatomii komponentów (→ `06-components.md`), sygnalizacji i stanów (→ plansza P4), makiet okien (→ `05-okna/`) |
| **Plik towarzyszący** | `04-portfolio/06-system-w-ruchu.html` |

---

## Spis treści

1. [Czym jest ta plansza](#1-czym-jest-ta-plansza)
2. [Kontrakt ruchu — 3/10 operacyjnie](#2-kontrakt-ruchu--310-operacyjnie)
3. [Cztery czasy i jeden ease](#3-cztery-czasy-i-jeden-ease)
4. [Katalog dozwolonych animacji](#4-katalog-dozwolonych-animacji)
5. [Katalog zakazanych animacji](#5-katalog-zakazanych-animacji)
6. [`prefers-reduced-motion`](#6-prefers-reduced-motion)
7. [Ruch jako nośnik informacji — pary porównawcze](#7-ruch-jako-nośnik-informacji--pary-porównawcze)
8. [Wydajność ruchu](#8-wydajność-ruchu)
9. [Oś czasu złożonego przejścia](#9-oś-czasu-złożonego-przejścia)
10. [Decyzje projektowe planszy](#10-decyzje-projektowe-planszy)
11. [Źródła](#11-źródła)
12. [Kontrola jakości](#12-kontrola-jakości)

---

## 1. Czym jest ta plansza

Danaco Console jest **kokpitem dowodzenia**, nie stroną marketingową. Kontrakt kierunku projektowego
ustawia pokrętło `INTENSYWNOSC_RUCHU` na **3/10** i to jedna liczba, z której wynika cały system ruchu:
mikroprzejścia **100–220 ms** plus **jeden ruch znaczący** — tętno kropki sygnału, czytane jako
„tu biegnie praca".

Plansza P6 zbiera ten system w jedno miejsce i pokazuje go **w działaniu**, a nie w opisie:
cztery czasy odpalane obok siebie na tym samym elemencie, krzywa `--dn-ease` narysowana jako
wykres z biegnącym po niej wskaźnikiem, siedem dozwolonych animacji uruchamianych przyciskiem,
dziesięć zakazanych z uzasadnieniem, symulator ograniczonego ruchu i rozbita na fazy oś czasu
przejścia z okna operacyjnego do modala.

**Nic na planszy nie jest wymyślone.** Wszystkie czasy, krzywe, nazwy klas i klatek kluczowych
pochodzą z `zasoby/zetony/zetony.css` (sekcja 5 i 14), `zasoby/css/komponenty.css`,
`zasoby/prototyp.css` oraz z rozdziałów 8–14 opracowania
`01-dokumentacja-md/08-handoff-motion-dostepnosc.md`. Wszystkie nazwy okien, modułów i ról
pochodzą z KANON rozdz. 6 i 7.

### 1.1. Jedno zdanie kontraktowe

> `INTENSYWNOSC_RUCHU` = **3/10** — mikroprzejścia 100–220 ms + **jeden ruch znaczący**
> (tętno kropki). Aplikacja robocza, nie strona marketingowa.
>
> — KIERUNEK.md rozdz. 2, powtórzone w KANON rozdz. 1 i w `08-handoff-motion-dostepnosc.md` rozdz. 8.1

---

## 2. Kontrakt ruchu — 3/10 operacyjnie

Liczba `3/10` sama w sobie nic nie znaczy. Dokumentacja rozkłada ją na **sześć mierzalnych wymiarów**
z testem zaliczenia — to jest właściwa treść kontraktu.

### 2.1. Sześć wymiarów pomiaru

| Wymiar | Wartość kontraktowa | Test zaliczenia |
|---|---|---|
| Górny czas przejścia | **220 ms** | żadna animacja interfejsu nie trwa dłużej; wyjątki: tętno 2,4 s (ciągły) i obrót spinnera 0,8 s (cykliczny) |
| Dolny czas przejścia | **100 ms** | poniżej tej wartości ruch jest nieczytelny — lepiej zmienić stan natychmiast |
| Liczba ruchów ciągłych na widok | **1** | tętno w miejscu, gdzie biegnie praca; dwa tętna w jednym widoku = błąd projektu |
| Ruch bez informacji | **0** | każda animacja odpowiada na pytanie „co się właśnie zmieniło w systemie?" |
| Dystans przesunięcia | **≤ 8 px** (`--dn-od-2`) | wejście warstwy `translateY(8px)` → 0 |
| Skala | **0,98 → 1** | wyłącznie wejście modala; brak „bounce", brak przeskoku ponad 1 |

### 2.2. Trzy pytania przed dodaniem animacji

```
1. Czy ruch niesie informację o stanie systemu?   NIE → NIE ANIMUJ
2. Czy mieści się w 100–220 ms?                   NIE → skróć albo zrezygnuj
3. Czy jest w katalogu dozwolonych (rozdz. 4)?    NIE → ZAKAZANA (rozdz. 5)
                                                  TAK → DOZWOLONA
```

### 2.3. Budżet ruchu na widok

| Pozycja | Limit |
|---|---|
| animacje ciągłe (`infinite`) | **1** — tętno; spinner jest ciągły wyłącznie na czas trwania operacji |
| jednoczesne przejścia wywołane jednym działaniem | **≤ 3** (wjazd modala = `opacity` + `transform` + nakładka) |
| elementy animowane przy przełączeniu motywu | `body` + powierzchnie dziedziczące — nigdy przejścia per komponent |

### 2.4. Jak to zmierzyć na gotowym oknie

| Krok | Co się robi | Wynik zaliczający |
|---|---|---|
| 1 | wyszukanie `animation:` i `transition:` w arkuszu okna | zero czasów wpisanych wprost — wyłącznie `var(--dn-czas-*)` |
| 2 | zliczenie `infinite` w widoku | **1** |
| 3 | wyszukanie krzywych | wyłącznie `var(--dn-ease)`; `linear` tylko przy `dn-obrot` |
| 4 | wyszukanie animowanych własności układu | brak `height`, `top`, `left`, `margin`; `width` wyłącznie w `.dn-postep-wartosc` |
| 5 | przełączenie systemowej preferencji ograniczonego ruchu | ruch znika, informacja zostaje (rozdz. 6) |

---

## 3. Cztery czasy i jeden ease

### 3.1. Krzywa

| Żeton | Wartość | Charakter |
|---|---|---|
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | szybki start, miękkie dojście do celu — ruch instrumentu, nie sprężyny |

**Jedna krzywa dla całego systemu.** Zakaz `ease-in-out`, `linear` (poza obrotem spinnera),
własnych `cubic-bezier` i wszelkich krzywych sprężystych z przeregulowaniem.

Odczyt krzywej: pierwszy punkt kontrolny `(0.2, 0)` leży nisko i blisko początku — ruch rusza
gwałtownie; drugi `(0, 1)` przyciąga tor do wartości docelowej długo przed końcem czasu — element
dojeżdża do celu i „siada". Operator widzi skutek działania niemal natychmiast, a mimo to nic
nie przeskakuje.

### 3.2. Cztery czasy

| Żeton | Wartość | Kiedy | Przykłady zastosowania |
|---|---|---|---|
| `--dn-czas-1` | **0,1 s** | mikroreakcja kontrolki — reakcja na kursor i naciśnięcie | `.dn-btn:hover`, `.dn-btn:active`, wiersz `.dn-tabela:hover`, `.dn-check::before` |
| `--dn-czas-2` | **0,16 s** | przejście barwy i przełączenie motywu | `body`, `.dn-zakladka`, `.dn-boczna-pozycja`, `.dn-przelacznik`, `.dn-karta-srodowiska::before` |
| `--dn-czas-3` | **0,22 s** | wejście i wyjście warstwy | `.dn-modal[open]`, `.dn-toast`, `.dn-tooltip-tresc`, `.dn-postep-wartosc` |
| `--dn-czas-tetno` | **2,4 s** | jedyny ruch ciągły | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after`, `.pt-tetno` |

### 3.3. Tabela decyzyjna „który czas"

| Co się zmienia | Czas | Własność CSS |
|---|---|---|
| tło kontrolki pod kursorem | `--dn-czas-1` | `background-color` |
| przesunięcie 1 px przy naciśnięciu | `--dn-czas-1` | `transform` |
| znacznik pola wyboru | `--dn-czas-1` | `transform: scale()` |
| barwa tekstu pozycji nawigacji | `--dn-czas-2` | `color`, `background-color` |
| przełączenie motywu całej strony | `--dn-czas-2` | `background-color`, `color` |
| gałka przełącznika | `--dn-czas-2` | `left`, `background-color` |
| obrys pola przy fokusie | `--dn-czas-2` | `border-color`, `box-shadow` |
| wjazd modala | `--dn-czas-3` | `opacity`, `transform` |
| wjazd powiadomienia | `--dn-czas-3` | `opacity`, `transform` |
| pojawienie się dymka | `--dn-czas-3` | `opacity`, `transform` |
| przyrost paska postępu | `--dn-czas-3` | `width` |
| tętno kropki pracy | `--dn-czas-tetno` | `box-shadow` |

### 3.4. Dwa czasy poza skalą

| Animacja | Czas | Dlaczego poza skalą |
|---|---|---|
| `dn-obrot` (spinner) | **0,8 s**, `linear` | ruch cykliczny bez punktu docelowego; `--dn-ease` wprowadziłby pulsowanie prędkości czytane jako zacinanie. Jedyne dopuszczone `linear` w systemie |
| `dn-tetno` | **2,4 s**, `--dn-ease` | ruch ciągły niosący stan „praca biegnie"; rytm wolniejszy od tętna spoczynkowego człowieka — nie pobudza, a pozostaje zauważalny peryferyjnie. **0,42 Hz**, siedmiokrotnie poniżej progu WCAG 2.3.1 (3 Hz) |

### 3.5. Demonstrator porównawczy (na planszy HTML)

Plansza uruchamia **ten sam element** czterema czasami równocześnie: cztery tory, cztery znaczniki,
jedno naciśnięcie. Wynik odczytywalny bez stopera:

| Tor | Czas | Co widać |
|---|---|---|
| 1 | 0,1 s | ruch praktycznie równoczesny z kliknięciem — czytany jako „reakcja", nie „animacja" |
| 2 | 0,16 s | wyraźne przejście, wciąż poniżej progu oczekiwania |
| 3 | 0,22 s | najdłuższy dopuszczalny — czytany jako „coś się pojawiło" |
| 4 | 2,4 s | nie jest przejściem; to rytm — pokazany dla skali, nigdy nie stosowany do przejść |

Obok torów: krzywa `--dn-ease` narysowana jako `<path>` w SVG (`M0,100 C20,100 100,0 100,0`
w układzie 100×100, oś Y odwrócona) z biegnącym po niej wskaźnikiem odpalanym `offset-path`.

---

## 4. Katalog dozwolonych animacji

**Katalog jest zamknięty.** Siedem pozycji. Animacja spoza katalogu wymaga decyzji Właściciela.

| # | Animacja | Nośnik | Czas | Własność | Informacja niesiona |
|---|---|---|---|---|---|
| **M1** | Mikroreakcja kontrolki | `.dn-btn`, `.dn-btn-ikona`, `.dn-karta--klikalna`, wiersz `.dn-tabela`, `.dn-zakladka` | `--dn-czas-1` | `background-color`, `transform: translateY(1px)` | „element jest interaktywny i przyjął moje działanie" |
| **M2** | Przejście barwy | `body`, `.dn-boczna-pozycja`, `.dn-przelacznik`, `.dn-pole-kontrolka:focus` | `--dn-czas-2` | `color`, `background-color`, `border-color`, `box-shadow` | „zmienił się stan wyboru / motyw / fokus" |
| **M3** | Wejście i wyjście warstwy | `.dn-modal[open]`, `.dn-toast`, `.dn-tooltip-tresc` | `--dn-czas-3` | `opacity` 0→1, `translateY(8px)`→0, modal dodatkowo `scale(0.98)`→1 | „pojawiła się nowa warstwa nad treścią" |
| **M4** | Tętno kropki | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after` | `--dn-czas-tetno` | `box-shadow` | **„tu biegnie praca"** — jedyny ruch ciągły |
| **M5** | Postęp kolejki | `.dn-postep-wartosc` | `--dn-czas-3` | `width` | „proces posunął się o zmierzoną wartość" |
| **M6** | Wskaźnik pracy wpisu | `.dn-btn[aria-busy='true']::after`, `.dn-spinner` | 0,8 s `linear` | `transform: rotate()` | „operacja trwa, czas nieznany" |
| **M7** | Przełączenie motywu | `body` + powierzchnie dziedziczące | `--dn-czas-2` | `background-color`, `color` | „zmieniłem motyw — to ta sama treść" |

### 4.1. Pozycje planszy a pozycje katalogu

Plansza HTML pokazuje **dwanaście demonstratorów**; każdy jest zastosowaniem jednej z siedmiu
pozycji katalogu — nie nową animacją:

| Demonstrator planszy | Pozycja katalogu | Czas |
|---|---|---|
| mikroreakcja kontrolki (hover / active) | M1 | `--dn-czas-1` |
| przejście barwy pozycji bocznej nawigacji | M2 | `--dn-czas-2` |
| przełączenie motywu | M7 (= M2 na `body`) | `--dn-czas-2` |
| wejście modala | M3 | `--dn-czas-3` |
| wejście i wyjście powiadomienia | M3 | `--dn-czas-3` |
| wejście panelu | M3 | `--dn-czas-3` |
| przejście widoku (`opacity` + `translateY(4px)`) | M3 (wariant wewnątrz okna) | `--dn-czas-3` |
| tętno kropki sygnału | M4 | `--dn-czas-tetno` |
| postęp kolejki | M5 | `--dn-czas-3` |
| wskaźnik pracy wpisu | M6 | 0,8 s `linear` |
| rozwinięcie sekcji (`<details>`) | M2 na treści o znanej wysokości | `--dn-czas-2` |
| zwijanie bocznej nawigacji | M2 na szerokości kolumny siatki powłoki | `--dn-czas-2` |

**Uwaga do dwóch ostatnich pozycji.** Dokumentacja rozdz. 11 (X7) zakazuje animowania własności
układu. Rozwinięcie sekcji realizuje więc natywny `<details>` z przejściem `opacity` + `transform`
na treści, a nie `height`. Zwijanie bocznej nawigacji zmienia **jedną deklarację
`grid-template-columns`** na kontenerze `.pt-cialo` (`.pt-cialo--zwinieta` w `prototyp.css`) —
przeliczenie dotyczy jednego kontenera powłoki, nie kaskady elementów. Odnotowane w rozdz. 10.

### 4.2. Wzorce implementacyjne — cytat z biblioteki

**M1 — mikroreakcja:**

```css
.dn-btn {
  transition:
    background-color var(--dn-czas-1) var(--dn-ease),
    border-color     var(--dn-czas-1) var(--dn-ease),
    color            var(--dn-czas-1) var(--dn-ease),
    transform        var(--dn-czas-1) var(--dn-ease);
}
.dn-btn:hover  { background: var(--dn-hover); }
.dn-btn:active { transform: translateY(1px); }
```

Reguła: własności wylicza się jawnie. `transition: all` animuje także własności układu i łamie X7.

**M4 — tętno:**

```css
.dn-kropka--tetno { animation: dn-tetno var(--dn-czas-tetno) var(--dn-ease) infinite; }

@keyframes dn-tetno {
  0%, 100% { box-shadow: 0 0 0 0  var(--dn-fokus-cien); }
  40%      { box-shadow: 0 0 0 5px transparent; }
}
```

Animowany jest `box-shadow`, nie `width`/`height` — pierścień rozchodzi się bez wpływu na układ sąsiadów.

**M3 — wejście modala:**

```css
.dn-modal[open] { animation: dn-wejscie var(--dn-czas-3) var(--dn-ease); }

@keyframes dn-wejscie {
  from { opacity: 0; transform: translateY(var(--dn-od-2)) scale(0.98); }
}
```

Klatka `to` pominięta świadomie — stan docelowy to stan spoczynkowy elementu.

**M5 — postęp:**

```css
.dn-postep-tor     { height: 4px; background: var(--dn-powierzchnia-2); overflow: hidden; }
.dn-postep-wartosc { height: 100%; background: var(--dn-sygnal-wypelnienie);
                     transition: width var(--dn-czas-3) var(--dn-ease); }
```

**M6 — wskaźnik pracy:**

```css
.dn-spinner { animation: dn-obrot 0.8s linear infinite; }
@keyframes dn-obrot { to { transform: rotate(360deg); } }
```

**M7 — motyw:**

```css
body {
  transition:
    background-color var(--dn-czas-2) var(--dn-ease),
    color            var(--dn-czas-2) var(--dn-ease);
}
```

### 4.3. M4 — tętno w pięciu miejscach kontraktu marki

| Miejsce | Forma | Ruch |
|---|---|---|
| godło (kropka zamykająca `»».`) | kropka sygnału | **statyczna** — znak się nie animuje |
| emblematy czterech środowisk | dokładnie jedna wypełniona kropka na emblemat | **statyczna** |
| karty sesji | pulsująca kropka = proces w tle | animowana |
| okno komunikacji | kropka przy nadawcy aktywnie piszącym | animowana (`.dn-wpis--pracuje`) |
| Always On Display | rdzeń awatara, pierścień wokół rdzenia | animowana |

**Reguła jednego tętna.** W jednym widoku pulsuje jedno miejsce. Gdy pracują trzy sesje —
pulsuje karta sesji aktywnej, pozostałe niosą kropkę statyczną plus liczbę w plakietce.

### 4.4. M5 — kiedy postęp, kiedy spinner

| Sytuacja | Komponent | Dlaczego |
|---|---|---|
| znany procent lub znana liczba kroków | `.dn-postep` + `.dn-postep-etykieta` (mono, `tabular-nums`) | pasek bez liczby to ozdoba; liczba bez paska to surowa dana |
| kolejka kroków o znanej liście | `.dn-kolejka` + `.dn-krok--pracuje` | postęp jest listą, nie paskiem |
| czas i zakres nieznane | `.dn-spinner` albo `.dn-btn[aria-busy='true']` | pasek udający postęp przy nieznanym zakresie to kłamstwo interfejsu |

---

## 5. Katalog zakazanych animacji

| # | Zakazane | Dlaczego szkodzi w kokpicie |
|---|---|---|
| **X1** | **Parallax** — warstwy przesuwające się z różną prędkością przy przewijaniu | ruch nie niesie żadnej informacji o stanie systemu; przy gęstości 8/10 rozbija odczyt danych; obciąża wątek kompozycji przy każdej klatce przewijania |
| **X2** | **Scrollytelling** — treść odsłaniana i animowana wraz z przewijaniem | wzorzec strony narracyjnej; okno robocze nie opowiada historii, pokazuje stan. Operator przewija, żeby czytać dane, nie żeby uruchamiać przedstawienie |
| **X3** | **Animacje dekoracyjne** — unoszące się kształty, animowane gradienty, cząstki, „oddychające" tła | anty-domyślne kierunku: zero dekoracji bez funkcji; gradient jest wyłącznie ilustracyjny i nigdy nie jest tłem sekcji |
| **X4** | **Bounce / przeregulowanie** — `cubic-bezier` przekraczający 1, sprężyny, `elastic` | system ma jedną krzywą `--dn-ease`; przeregulowanie sugeruje fizyczność obiektu, a kontrolki kokpitu nie są obiektami fizycznymi |
| **X5** | **Ruch bez informacji** — animacja wejścia strony, kaskadowe pojawianie się list, animowane liczniki | animowany licznik **fałszuje pomiar**: pokazuje wartości, których system nigdy nie zmierzył |
| **X6** | **Więcej niż jeden ruch ciągły na widok** | tętno traci znaczenie „tu biegnie praca", gdy pulsuje wszystko |
| **X7** | **Animacja własności układu** — `height`, `top`, `left`, `margin`, `padding` | wymusza przeliczenie geometrii w każdej klatce; jedyne dopuszczone `width` to `.dn-postep-wartosc` w torze o stałych wymiarach |
| **X8** | **Migotanie** powyżej 3 Hz | kryterium WCAG 2.3.1, ryzyko napadu światłoczułego. Tętno 2,4 s = 0,42 Hz |
| **X9** | **Rozmycie jako efekt** — `backdrop-filter` poza nakładką modala | anty-domyślne „glassmorfizm wszędzie" → powierzchnie kryjące; jedyne dopuszczone rozmycie to `blur(2px)` na `.dn-modal::backdrop` |
| **X10** | **Animowane przewijanie sterowane skryptem** — globalne `scroll-behavior: smooth` | odbiera Operatorowi kontrolę nad tempem czytania; przy `prefers-reduced-motion` system wymusza `scroll-behavior: auto` |

### 5.1. Zakazy dodatkowe wynikające z zakazów podstawowych

| Odruch | Który zakaz go obejmuje |
|---|---|
| „animacja wjazdu kart środowisk na Centrum dowodzenia" | X5 |
| „delikatny puls przycisku wysyłki, żeby zachęcić" | X3 + X6 |
| „licznik zadań w kolejce odliczający od 0 do 47" | X5 |
| „miękkie rozwinięcie panelu przez `height: auto`" | X7 |
| „przyciemnianie tła paska górnego przy przewijaniu" | X1 |
| „gradient przesuwający się w tle Always On Display" | X3 |

---

## 6. `prefers-reduced-motion`

### 6.1. Obsługa globalna, nie per komponent

Blok żyje w `zasoby/zetony/zetony.css`, sekcja 14 — wiersze 435–441:

```css
@media (prefers-reduced-motion: reduce) {
  :root {
    --dn-czas-1: 0.01ms;
    --dn-czas-2: 0.01ms;
    --dn-czas-3: 0.01ms;
    --dn-czas-tetno: 0.01ms;
  }
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

| Mechanizm | Co obejmuje | Dlaczego potrzebny |
|---|---|---|
| podmiana żetonów czasu | wszystko, co używa `--dn-czas-*` | źródło prawdy — komponent nie musi wiedzieć o preferencji |
| reguła `*` z `!important` | animacje o czasach spoza skali (0,8 s spinnera) i kod zewnętrzny | siatka bezpieczeństwa |

`animation-iteration-count: 1` zatrzymuje pętle po pierwszym przebiegu — animacja nie „ucina się"
w losowej klatce.

### 6.2. Tętno → pierścień statyczny

Skrócenie czasu wygasza ruch, ale usunęłoby też **informację**. Dlatego biblioteka podstawia
zamiennik statyczny (`prototyp.css`, `.pt-tetno`; analogicznie `.dn-kropka--tetno`):

```css
@media (prefers-reduced-motion: reduce) {
  .pt-tetno { animation: none; box-shadow: 0 0 0 2px var(--dn-fokus-cien); }
}
```

```
   RUCH DOZWOLONY                    RUCH OGRANICZONY
   ((( ● )))  pierścień              ( ● )  pierścień statyczny 2 px
   rozchodzi się, pętla 2,4 s        obecny zawsze, bez ruchu
   informacja: praca biegnie         informacja: praca biegnie
   nośnik: RUCH                      nośnik: KSZTAŁT
```

**Zasada ogólna: ograniczenie ruchu nie może usuwać informacji.** Jeżeli ruch był jedynym
nośnikiem, przy `reduce` informację przejmuje kształt, obrys albo etykieta.

### 6.3. Tabela zachowań przy `reduce`

| Animacja | Zachowanie | Nośnik informacji po zmianie |
|---|---|---|
| M1 mikroreakcja | zmiana natychmiastowa | barwa tła (`--dn-hover`) |
| M2 przejście barwy | zmiana natychmiastowa | barwa docelowa |
| M3 wejście warstwy | warstwa pojawia się bez wjazdu | obecność warstwy + cień + nakładka |
| **M4 tętno** | **pierścień statyczny 2 px** | **kształt (pierścień)** |
| M5 postęp | pasek skacze do wartości | liczba w etykiecie mono |
| M6 spinner | zatrzymany po pierwszym obrocie | `aria-busy="true"` + tekst „Trwa…" |
| M7 motyw | przełączenie natychmiastowe | nowe barwy |

### 6.4. Czego NIE robić

| Zakaz | Uzasadnienie |
|---|---|
| dublowanie zapytania `prefers-reduced-motion` w każdym komponencie | obsługa jest globalna; duplikat rozjeżdża się przy pierwszej zmianie |
| wykrywanie preferencji w JavaScript i wyłączanie klas | preferencja może zmienić się w trakcie sesji; CSS reaguje sam |
| całkowite usunięcie zamiennika przy `reduce` | usuwa informację o stanie systemu |

### 6.5. Symulator na planszy

Plansza ma przełącznik „symuluj ograniczony ruch". Ustawia on na korzeniu planszy atrybut
`data-ruch="ograniczony"`, który **lokalnie i wyłącznie na potrzeby demonstracji** podmienia te
same żetony czasu i podstawia pierścień statyczny. Symulator jest ilustracją mechanizmu, nie jego
implementacją — implementacja żyje w żetonach i działa niezależnie od tego przełącznika.

---

## 7. Ruch jako nośnik informacji — pary porównawcze

Trzy pary. W każdej: ten sam przepływ, raz z ruchem, raz bez. Bez ruchu przepływ dalej działa,
ale Operator traci konkretną informację o stanie systemu.

### 7.1. Para A — praca w tle sesji

| | Z ruchem | Bez ruchu |
|---|---|---|
| Nośnik | kropka tętniąca przy karcie sesji | kropka statyczna |
| Odczyt | „w tej sesji **teraz** coś się dzieje" | „ta sesja istnieje" |
| Co się gubi | rozróżnienie sesji pracującej od bezczynnej bez czytania etykiet | — |
| Rekompensata przy `reduce` | pierścień statyczny 2 px + plakietka liczby | wymagana |

Zastosowanie: pas kart sesji w powłokach TalkIn, WorkSpace, CodeStudio, MultitaskingAI.

### 7.2. Para B — przejście stanu kroku w kolejce

| | Z ruchem | Bez ruchu |
|---|---|---|
| Nośnik | `.dn-postep-wartosc` przesuwa się w `--dn-czas-3`; `.dn-krok--pracuje` → `--poprawny` | wartość skacze |
| Odczyt | „proces posunął się **o tyle**" — widać kierunek i wielkość skoku | „proces jest w innym miejscu niż przed chwilą" |
| Co się gubi | wielkość przyrostu, czyli tempo pracy kolejki | — |
| Rekompensata przy `reduce` | liczba w etykiecie mono (`tabular-nums`) | wymagana |

Zastosowanie: Queue Manager i Execution Monitor modułu Automations, silnik kolejek MultitaskingAI.

### 7.3. Para C — pojawienie się warstwy

| | Z ruchem | Bez ruchu |
|---|---|---|
| Nośnik | modal `opacity` 0→1 + `translateY(8px)` + `scale(0.98)` w `--dn-czas-3` | modal pojawia się natychmiast |
| Odczyt | „nad treścią **stanęła** nowa warstwa; poprzednia nadal jest pod spodem" | „ekran wygląda inaczej" |
| Co się gubi | relacja przestrzenna warstwy do tła — skąd przyszła i co przykryła | — |
| Rekompensata przy `reduce` | nakładka `--dn-nakladka` + `blur(2px)` + cień warstwy | wymagana |

Zastosowanie: modal „Pula kont Code CLI" w Oknie Konfiguracji, okno punktów izolacji, potwierdzenia
w Permissions Center modułu Agents.

### 7.4. Wniosek

Ruch nie jest w tym systemie warstwą estetyczną, którą można zdjąć bez konsekwencji. Jest
**trzecim nośnikiem stanu** obok barwy i kształtu — i jako jedyny z trzech może zostać wyłączony
decyzją systemową użytkownika. Dlatego każde jego użycie ma zdefiniowany zamiennik.

---

## 8. Wydajność ruchu

### 8.1. Etapy potoku renderowania

| Własność | Etap | Koszt | Status |
|---|---|---|---|
| `transform` | kompozycja | najniższy | **preferowana** |
| `opacity` | kompozycja | najniższy | **preferowana** |
| `box-shadow` | malowanie | średni | dopuszczona: tętno, fokus |
| `color`, `background-color` | malowanie | średni | dopuszczona: M1, M2 |
| `width`, `height`, `top`, `left`, `margin`, `padding` | układ | najwyższy | **zakazane** poza `.dn-postep-wartosc` |

### 8.2. Dlaczego `.dn-postep-wartosc` może animować `width`

Element leży w `.dn-postep-tor` z `overflow: hidden` i stałą wysokością 4 px. Przeliczenie układu
ogranicza się do jednego niezależnego elementu i nie propaguje na rodzeństwo. Zamiana na
`transform: scaleX()` rozciągałaby zaokrąglenia narożników — pogorszenie wizualne bez zysku
pomiarowego przy jednym pasku na widok.

### 8.3. `will-change` — zasada trzech warunków

Stosować **wyłącznie** gdy spełnione są równocześnie:

1. animacja jest długa lub cykliczna (tętno, spinner),
2. element ma znany, ograniczony rozmiar,
3. deklaracja jest **zdejmowana** po zakończeniu animacji.

| Zakaz | Skutek |
|---|---|
| `will-change` na wielu elementach naraz | każdy dostaje własną warstwę kompozycji — zużycie pamięci graficznej rośnie liniowo |
| `will-change` na stałe w arkuszu | warstwa nigdy nie znika; przeglądarka traci możliwość optymalizacji |
| `will-change: all` | przeciwieństwo optymalizacji |

**Stan biblioteki:** `komponenty.css` nie deklaruje `will-change` w żadnym miejscu. Animowane
elementy są małe (kropka 6 px, spinner 14 px, pierścień AOD 36 px) i przeglądarka promuje je do
warstwy samodzielnie. Wprowadzenie `will-change` wymaga pomiaru pokazującego zysk.

### 8.4. Unikanie przeliczeń układu (layout thrash)

| Wzorzec błędny | Wzorzec poprawny |
|---|---|
| w pętli: odczyt `offsetHeight` → zapis `style.height` → odczyt → zapis | wszystkie odczyty **przed** wszystkimi zapisami |
| animowanie `height` przy rozwijaniu sekcji | natywny `<details>` albo `opacity` + `transform` na treści o znanej wysokości |
| animowanie `top`/`left` przy przesuwaniu warstwy | `transform: translateY()` |
| pomiar szerokości paska postępu w skrypcie przy każdej klatce | zmiana deklaracji `width` raz; przejście wykonuje CSS |

---

## 9. Oś czasu złożonego przejścia

Przejście z okna operacyjnego do modala — na przykładzie **Okna Konfiguracji** i modala
**„Pula kont Code CLI"** (KANON rozdz. 7.3, poz. 8). Cała sekwencja mieści się w `--dn-czas-3`
= 220 ms i składa się z trzech przejść jednoczesnych (limit budżetu z rozdz. 2.3).

| Faza | Zakres czasu | Co się dzieje | Własność | Informacja |
|---|---|---|---|---|
| **F0** | 0 ms | naciśnięcie przycisku otwierającego | `transform: translateY(1px)` (M1, 100 ms) | „przyjąłem działanie" |
| **F1** | 0–220 ms | nakładka `::backdrop` wchodzi z `opacity` 0→1, `blur(2px)` | `opacity` | „treść pod spodem przestaje być aktywna" |
| **F2** | 0–220 ms | modal `opacity` 0→1 | `opacity` | „jest nowa warstwa" |
| **F3** | 0–220 ms | modal `translateY(8px)` → 0 oraz `scale(0.98)` → 1 | `transform` | „warstwa **stanęła nad** treścią, przyszła z dołu" |
| **F4** | 220 ms | fokus przechodzi na pierwszą kontrolkę modala | brak animacji | „tu teraz pracuję" |

### 9.1. Rozkład na klatki

```
   0 ms        55 ms       110 ms      165 ms      220 ms
   ├───────────┼───────────┼───────────┼───────────┤
   F1 nakładka  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▶ opacity 1
   F2 modal     ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▶ opacity 1
   F3 tor       8px ──── 3,1px ─ 1,1px ─ 0,2px ─── 0
                scale 0,98 ──────────────────── 1,00
   F4 fokus                                      ● pierwsza kontrolka
```

Wartości toru F3 wynikają z krzywej `cubic-bezier(0.2, 0, 0, 1)`: po połowie czasu element pokonał
już około 86% dystansu. To jest właśnie „ruch instrumentu" — energia na starcie, cisza na końcu.

### 9.2. Zamknięcie modala

| Faza | Czas | Co się dzieje |
|---|---|---|
| Z0 | 0 ms | `Esc`, przycisk zamknięcia albo kliknięcie w nakładkę |
| Z1 | 0–220 ms | nakładka i modal wygaszają `opacity` do 0 |
| Z2 | 220 ms | fokus wraca na kontrolkę, która modal otworzyła |

**Bez ruchu powrotnego w osi pionowej.** Wyjście warstwy nie odtwarza wejścia w odwrotną stronę —
element znika, a nie „chowa się z powrotem". Odwrócenie toru sugerowałoby fizyczne schowanie
obiektu w miejscu, którego w kokpicie nie ma.

### 9.3. Ta sama oś dla powiadomienia

Toast dzieli klatki kluczowe z modalem (`dn-wejscie`): jedno wejście warstwy = jedna definicja
ruchu. Różnica leży wyłącznie w warstwie z-index (powiadomienie 1000, modal 900) i w tym, że
powiadomienie nie przejmuje fokusu.

---

## 10. Decyzje projektowe planszy

| # | Decyzja | Uzasadnienie |
|---|---|---|
| **D1** | Demonstrator czterech czasów uruchamia **ten sam element** czterokrotnie zamiast czterech różnych | porównanie ma dotyczyć wyłącznie czasu; różny kształt zaburzyłby odczyt |
| **D2** | Krzywa `--dn-ease` narysowana jako SVG `<path>` z biegnącym wskaźnikiem | krzywa opisana liczbami jest nieczytelna; narysowana i przebiegnięta staje się argumentem |
| **D3** | Tor 2,4 s obecny w demonstratorze mimo że nie jest czasem przejścia | pokazuje **skalę** — bez niego 220 ms wygląda na dużo, z nim widać, że jest jedenastokrotnie krótsze |
| **D4** | Symulator ograniczonego ruchu działa na atrybucie `data-ruch` planszy, nie na klasach komponentów | dokumentacja zakazuje wykrywania preferencji w JS; symulator jest wyraźnie oznaczony jako demonstracja mechanizmu |
| **D5** | Katalog zakazanych ma tę samą wagę wizualną co katalog dozwolonych | w opracowaniu portfolio to, czego system nie robi, definiuje go tak samo jak to, co robi |
| **D6** | Żaden demonstrator zakazanej animacji nie jest odtwarzany | plansza nie może uruchamiać parallaxu ani bounce'u „dla przykładu" — złamałaby własny kontrakt. Zakazy pokazujemy tekstem i statycznym diagramem |
| **D7** | Rozwinięcie sekcji zrealizowane natywnym `<details>` | X7 zakazuje animowania `height`; natywny element daje semantykę i klawiaturę bez kodu |
| **D8** | Zwijanie bocznej nawigacji zmienia `grid-template-columns` jednego kontenera | to jedyna droga do zwinięcia bez animowania `width` każdej pozycji; przeliczenie dotyczy jednego kontenera powłoki |
| **D9** | Wszystkie fragmenty CSS na planszy mają przycisk kopiowania (`data-kopiuj`) | plansza jest źródłem gotowego kodu dla dewelopera; przepisywanie ręczne to droga do rozjazdu z biblioteką |
| **D10** | Oś czasu przejścia liczona na przykładzie realnego modala „Pula kont Code CLI" | KANON rozdz. 10 zakazuje przykładów spoza domeny produktu |

---

## 11. Źródła

| Zakres | Plik i rozdział |
|---|---|
| Pokrętło `INTENSYWNOSC_RUCHU` 3/10 | `KANON.md` rozdz. 1 (Trzy pokrętła); `KIERUNEK.md` rozdz. 2, wiersz 29 |
| Cztery czasy i krzywa | `KANON.md` rozdz. 2 (Ruch); `zasoby/zetony/zetony.css` sekcja 5, wiersze 126–130 |
| Wzorzec animacji prototypów | `KANON.md` rozdz. 11 (Standard techniczny, „Wzorzec animacji") |
| Kontrakt ruchu — sześć wymiarów | `01-dokumentacja-md/08-handoff-motion-dostepnosc.md` rozdz. 8.2 |
| Trzy pytania przed animacją | tamże, rozdz. 8.3 |
| Cztery czasy, tabela decyzyjna | tamże, rozdz. 9.1–9.4 |
| Katalog dozwolonych M1–M7 | tamże, rozdz. 10, 10.1, 10.2 |
| Katalog zakazanych X1–X10 | tamże, rozdz. 11 |
| Wzorce implementacyjne | tamże, rozdz. 12.1–12.7 (cytat z biblioteki) |
| `prefers-reduced-motion` | tamże, rozdz. 13.1–13.4; `zasoby/zetony/zetony.css` sekcja 14, wiersze 435–441 |
| Wydajność, `will-change`, layout thrash | tamże, rozdz. 14.1–14.5 |
| Tętno jako element sygnaturowy | `KANON.md` rozdz. 3 (Kropka sygnału); `KIERUNEK.md` rozdz. 3.4, wiersz 87 |
| Klasy komponentów i klatki kluczowe | `zasoby/css/komponenty.css` (`.dn-postep` 867–893, `.dn-krok` 896–944, `.dn-toast` 1000–1024) |
| `.pt-tetno`, `.pt-wejscie`, `.pt-cialo--zwinieta` | `zasoby/prototyp.css` wiersze 30, 288–309 |
| Interakcje sterowane atrybutami | `zasoby/prototyp.js` (przełączanie widoków, modale, `data-kopiuj`, `data-symuluj`, `data-postep-do`, spis treści) |
| Nazwy okien użytych w przykładach | `KANON.md` rozdz. 7.3 (Okno Konfiguracji, modal „Pula kont Code CLI"), 7.5 (Queue Manager, Execution Monitor, Permissions Center) |

---

## 12. Kontrola jakości

| Kryterium | Wynik |
|---|---|
| Wartości szesnastkowe w `<style>` planszy | **0** — wyłącznie `var(--dn-*)` |
| Wystąpienia `disabled` | **0** |
| Emoji jako ikony | **0** — wyłącznie inline SVG 24×24, obrys 1,75, `currentColor`, pliki z `zasoby/ikony/svg/` |
| Ruchy ciągłe na planszy | tętno demonstracyjne odpalane przyciskiem; poza demonstratorem widok bez ruchu ciągłego |
| Czasy użyte w arkuszu lokalnym | wyłącznie `var(--dn-czas-1/2/3/-tetno)` |
| Krzywe użyte w arkuszu lokalnym | wyłącznie `var(--dn-ease)`; `linear` wyłącznie przy demonstratorze M6 |
| Animowane własności układu | brak; `width` wyłącznie w demonstratorze M5 wewnątrz toru |
| `prefers-reduced-motion` | obsłużone globalnie żetonami; plansza nie dubluje zapytania poza jawnie oznaczonym symulatorem i zamiennikiem `.pt-tetno` |
| Oba motywy | działają; pasek górny atramentowy w obu |
| Dostępność | `lang="pl"`, `role`/`aria`, odnośnik pominięcia, fokus widoczny, spis sekcji z `aria-current` |
| Nazwy własne | wyłącznie z KANON rozdz. 6 i 7 |
| Powrót do indeksu | `../INDEKS.html` — dwa odnośniki (pasek górny, spis sekcji) |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
