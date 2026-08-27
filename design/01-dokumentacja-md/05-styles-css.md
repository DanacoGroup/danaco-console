# Danaco Console — STYLES + CSS: style i architektura arkuszy

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — architektura arkuszy stylów |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 · fundament systemu projektowego |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | deweloper wdrażający warstwę wizualną w `budowa/client/src/`; projektant utrzymujący bibliotekę `.dn-*`; recenzent arkusza przed scaleniem; autor nowego okna operacyjnego |
| **Zakres** | warstwa fundamentu (`fundament.css`), biblioteka komponentów (`komponenty.css`), arkusze widoku (`rama.css`, `stanowisko.css`, `przedsionek.css`), kolejność importu, konwencja nazewnicza, kontrakt specyficzności, wzorce stanu / fokusu / przejścia, podział na pliki per komponent, obsługa motywu i preferencji środowiska, antywzorce, lista kontrolna przeglądu |
| **Czego NIE zawiera** | wartości żetonów (są w opracowaniu A4 „Żetony” i w `zetony.css` / `zetony.json`), pomiarów kontrastu (`kontrasty.json`), katalogu ikon (`manifest.json`), makiet okien (katalog `05-okna/`), zasad marki i znaku (opracowania 03-marka / 09-brand-system), warstwy JavaScript poza `wspolne.js` |

---

## Spis treści

1. [Warstwa fundamentu — rola, zakres, granice](#1-warstwa-fundamentu--rola-zakres-granice)
2. [Kolejność importu arkuszy i jej konsekwencje](#2-kolejność-importu-arkuszy-i-jej-konsekwencje)
3. [Konwencja nazewnicza klas `.dn-*`](#3-konwencja-nazewnicza-klas-dn-)
4. [Specyficzność jako kontrakt — zakaz zagnieżdżania](#4-specyficzność-jako-kontrakt--zakaz-zagnieżdżania)
5. [Pełny wykaz reguł fundamentu z wartościami](#5-pełny-wykaz-reguł-fundamentu-z-wartościami)
6. [Wzorzec stanu](#6-wzorzec-stanu)
7. [Wzorzec fokusu globalnego](#7-wzorzec-fokusu-globalnego)
8. [Wzorzec przejścia i animacji](#8-wzorzec-przejścia-i-animacji)
9. [Podział `komponenty.css` na pliki per komponent](#9-podział-komponentycss-na-pliki-per-komponent)
10. [Wzorce użytkowe (utility) — co istnieje, czego się nie tworzy](#10-wzorce-użytkowe-utility--co-istnieje-czego-się-nie-tworzy)
11. [Obsługa motywu w CSS](#11-obsługa-motywu-w-css)
12. [Preferencje środowiska: ruch, dotyk, wymuszone barwy](#12-preferencje-środowiska-ruch-dotyk-wymuszone-barwy)
13. [Antywzorce CSS w tym systemie](#13-antywzorce-css-w-tym-systemie)
14. [Lista kontrolna przeglądu arkusza przed scaleniem](#14-lista-kontrolna-przeglądu-arkusza-przed-scaleniem)
15. [Decyzje projektowe](#15-decyzje-projektowe)

---

## 1. Warstwa fundamentu — rola, zakres, granice

### 1.1. Deklaracja z nagłówka pliku

Nagłówek `zasoby/css/fundament.css` brzmi dosłownie:

> Warstwa zerowa nad żetonami: wyzerowanie, podłoże dokumentu, typografia bazowa, fokus, zaznaczenie, paski przewijania, drobne wzorce tekstowe.
> Wymaga: `03-zetony/fonty.css` oraz `03-zetony/zetony.css`.

Zdanie „wymaga” jest wiążące: **fundament nie definiuje ani jednej wartości własnej.** Każda deklaracja sięga po żeton. Plik nie zawiera ani jednej wartości szesnastkowej — sprawdzone: zero trafień.

### 1.2. Siedem odpowiedzialności fundamentu

| # | Odpowiedzialność | Linie | Co konkretnie robi |
|---|---|---|---|
| 1 | **Wyzerowanie i model pudełka** | 9–13 | `box-sizing: border-box` na wszystkim wraz z pseudoelementami; `margin: 0` na wszystkim; `html { height: 100% }` |
| 2 | **Podłoże dokumentu** | 15–28 | `body` przyjmuje `--dn-tlo`, `--dn-tekst`, krój bazowy, stopień bazowy, interlinię bazową; `min-height: 100dvh`; przejście barw przy przełączeniu motywu |
| 3 | **Dziedziczenie krojów w kontrolkach** | 30–34 | media blokowe (`img, svg, video, canvas`); `input, button, textarea, select { font: inherit; color: inherit }`; łamanie długich wyrazów w `p, h1–h4` |
| 4 | **Typografia nagłówków** | 36–42 | `h1, h2, h3` **oraz** klasa `.dn-naglowek` — Space Grotesk, waga półgruba, światło międzyliterowe nagłówkowe, interlinia ciasna |
| 5 | **Fokus globalny** | 44–49 | jedna reguła `:focus-visible` dla całego produktu (rozdz. 7) |
| 6 | **Zaznaczenie i paski przewijania** | 51–64 | `::selection` w dwóch wariantach motywu; paski przewijania w standardzie `scrollbar-*` i `::-webkit-scrollbar` |
| 7 | **Drobne wzorce tekstowe** | 66–156 | odnośniki `a` oraz jedenaście klas pomocniczych: `.dn-etykieta-wersalikowa`, `.dn-etykieta-mono`, `.dn-dane`, `.dn-liczba`, `.dn-kod`, `.dn-kod--wiersz`, `.dn-kbd`, `.dn-separator`, `.dn-separator--pionowy`, `.dn-sr-only`, `.dn-naglowek` |

### 1.3. Granica: co należy do fundamentu, a co do komponentów

```
   ŻETONY (zetony.css)        ── wartości. Nie zna elementów HTML.
        │
   FUNDAMENT (fundament.css)  ── elementy HTML. Nie zna komponentów.
        │                        Jedyna warstwa, która wolno jej
        │                        stylować goły selektor elementowy.
        │
   KOMPONENTY (komponenty.css)── klasy .dn-*. Nie zna elementów HTML
        │                        inaczej niż przez kotwicę klasową.
        │
   OKNO (arkusz lokalny)      ── układ i kompozycja. Nie zmienia
                                 wyglądu komponentu, ustawia go w siatce.
```

**Test przynależności reguły do fundamentu** — reguła należy do fundamentu wtedy i tylko wtedy, gdy:
1. dotyczy elementu HTML, który występuje w każdym oknie produktu (`body`, `a`, `input`, pasek przewijania), **oraz**
2. jej brak spowodowałby, że okno wygląda niespójnie **jeszcze zanim** wstawi się do niego jakikolwiek komponent, **oraz**
3. nie da się jej wyrazić żetonem (żeton przenosi wartość, nie zachowanie).

Reguła, która nie przechodzi wszystkich trzech testów, należy do komponentu albo do arkusza lokalnego okna.

### 1.4. Dlaczego fundament jest tak mały (156 linii)

Kierunek projektowy określa `GESTOSC_WIZUALNA` na 8/10 i mówi „zero dekoracji bez funkcji”. Rozrost fundamentu ma dwa znane skutki uboczne, oba sprzeczne z kierunkiem:

- **stylowanie „w ciemno”** — reguła na gołym `button` albo `ul` dosięga elementów, których autor reguły nigdy nie widział; w kokpicie o 15 modułach to gwarantowana regresja,
- **niewidoczna specyficzność** — deklaracja elementowa przegrywa z każdą klasą, więc autor komponentu i tak ją nadpisuje; efektem jest martwy kod, który nadal trzeba czytać.

Fundament stylizuje **tylko to, czego komponent nie może przejąć**: podłoże strony, dziedziczenie krojów, fokus, pasek przewijania, zaznaczenie.

---

## 2. Kolejność importu arkuszy i jej konsekwencje

### 2.1. Kolejność obowiązująca

```html
<link rel="stylesheet" href="../zasoby/zetony/fonty.css">      <!-- 1 -->
<link rel="stylesheet" href="../zasoby/zetony/zetony.css">     <!-- 2 -->
<link rel="stylesheet" href="../zasoby/css/fundament.css">     <!-- 3 -->
<link rel="stylesheet" href="../zasoby/css/komponenty.css">    <!-- 4 -->
<style>/* 5 — style lokalne okna, prefiks lokalny */</style>
```

Kolejność jest identyczna w kontrakcie systemu projektowego (standard techniczny prototypów) i w opracowaniu o przekazaniu, ruchu i dostępności („import po żetonach, przed komponentami”). Nie jest kwestią gustu — każdy krok ma twarde uzasadnienie.

### 2.2. Uzasadnienie krok po kroku

| Krok | Arkusz | Co wnosi | Dlaczego przed następnym |
|---|---|---|---|
| 1 | `fonty.css` | reguły `@font-face` (Space Grotesk, IBM Plex Sans, IBM Plex Mono, `unicode-range`, licencje OFL) | Rodzina kroju musi być **zarejestrowana** zanim jakakolwiek reguła jej użyje. Deklaracja `@font-face` po pierwszym użyciu działa, ale wywołuje przeliczenie układu po dociągnięciu pliku — na kokpicie o gęstości zwartej widać to jako drgnięcie wierszy. |
| 2 | `zetony.css` | wszystkie `--dn-*`: prymitywy, żetony semantyczne obu motywów, typografia, przestrzeń, ruch, wymiary, warstwy | Własność niestandardowa musi istnieć **w chwili obliczania** reguły, która ją czyta. `var(--dn-tlo)` odczytane przed zdefiniowaniem `--dn-tlo` daje wartość nieprawidłową i element traci tło. |
| 3 | `fundament.css` | podłoże, dziedziczenie, fokus, drobne wzorce | Musi być **przed** komponentami, bo jego reguły elementowe mają być nadpisywalne klasą. Zamiana miejscami odwróciłaby zależność — patrz 2.3. |
| 4 | `komponenty.css` | biblioteka `.dn-*` | Wygrywa z fundamentem przy równej specyficzności, bo jest później. To jest zamierzone. |
| 5 | `rama.css` | belka tytułowa i pasek narzędzi okna | Rama jest wspólna wszystkim oknom platformy i musi móc nadpisać ustawienia powłoki (kierunek układu, rozdział wysokości), dlatego stoi po bibliotece komponentów. |
| 6 | `stanowisko.css` · `przedsionek.css` | układ wielookienny przestrzeni roboczej · widok wejściowy środowiska | Arkusze widoku: jeden opisuje pracę w module, drugi wejście w środowisko. Rozłączne — okno wczytuje ten, który go dotyczy. |
| 7 | styl lokalny okna | kompozycja, siatka, wyjątki jednego okna | Ostatni, bo wyjątek zawsze musi mieć możliwość zwycięstwa bez podnoszenia specyficzności. |

### 2.3. Dowód konieczności kolejności 3 → 4

Fundament zawiera regułę:

```css
:focus-visible {
  outline: var(--dn-wym-fokus) solid var(--dn-fokus);
  outline-offset: var(--dn-wym-fokus-odsuniecie);
  border-radius: var(--dn-r-xs);
}
```

Selektor `:focus-visible` ma specyficzność **(0,1,0)** — dokładnie tyle samo, co selektor klasowy `.dn-btn`, który deklaruje `border-radius: var(--dn-r-sm)`.

| Kolejność importu | Kto wygrywa `border-radius` przy fokusie klawiaturowym | Skutek widoczny |
|---|---|---|
| fundament → komponenty (**obowiązująca**) | `.dn-btn` — 6 px | przycisk zachowuje kształt; pierścień fokusu obrysowuje go zgodnie z promieniem własnym |
| komponenty → fundament (**błędna**) | `:focus-visible` — 3 px | przycisk **zmienia kształt w chwili sfokusowania**; narożniki „skaczą” z 6 px na 3 px |

Deklaracja `border-radius` w regule fokusu nie jest przypadkiem: obsługuje elementy, które **nie mają promienia własnego** (odnośnik w tekście, sfokusowany kontener listy, pozycja spisu treści). Dla nich pierścień byłby ostrym prostokątem, obcym dla systemu o promieniach 3–14 px. Reguła działa jak wartość domyślna, którą komponent nadpisuje samym faktem posiadania własnego promienia — **pod warunkiem, że jest importowany później.**

### 2.4. Co się dzieje przy braku któregoś arkusza

| Brakuje | Objaw |
|---|---|
| `fonty.css` | kroje z listy zapasowej (`Segoe UI`, `system-ui`, `Consolas`) — układ nie pęka, ale znikają wersaliki Space Grotesk i `tabular-nums` w danych |
| `zetony.css` | strona bez tła i bez barw tekstu; wszystkie `var(--dn-*)` nieprawidłowe; **awaria całkowita** |
| `fundament.css` | podwójne marginesy przeglądarki, brak `border-box`, brak fokusu na elementach spoza biblioteki, paski przewijania systemowe, brak `.dn-sr-only` (etykiety dla czytników stają się widoczne) |
| `komponenty.css` | goły dokument tekstowy — czytelny, bez kontrolek; **stan degradacji akceptowalny**, treść pozostaje dostępna |

Ostatni wiersz jest miarą jakości podziału: usunięcie warstwy komponentów pozostawia dokument czytelny. Tak wygląda poprawnie warstwowy system.

### 2.5. Miejsce `wspolne.js`

Skrypt `zasoby/wspolne.js` ładuje się **na końcu `<body>`**, nie w `<head>`. Ustawia `data-theme` na `<html>` (z zapisu w `localStorage`, w zapasie z `prefers-color-scheme`) i obsługuje `[data-przelacz-motyw]` oraz zapas dla poleceń natywnych `commandfor`/`command`. Skrypt **nie dodaje żadnych stylów** — zmienia wyłącznie atrybut, po którym kaskaduje `zetony.css`. To jedyny dozwolony sposób wpływania JavaScriptu na wygląd w tym systemie.

---

## 3. Konwencja nazewnicza klas `.dn-*`

### 3.1. Trzy formy nazwy

| Forma | Wzorzec | Znaczenie | Przykład z arkusza |
|---|---|---|---|
| **Blok** | `.dn-<komponent>` | samodzielny komponent; niesie układ, wymiary, obrys, tło | `.dn-btn`, `.dn-karta`, `.dn-wpis`, `.dn-tabela` |
| **Modyfikator** | `.dn-<komponent>--<wariant>` | odmiana bloku; **nigdy nie występuje sam** — wymaga bloku w atrybucie `class` | `.dn-btn--atrament`, `.dn-krok--wstrzymany` |
| **Element** | `.dn-<komponent>-<część>` | część składowa, która nie istnieje poza blokiem | `.dn-karta-tytul`, `.dn-wpis-medalion` |

Zapis w HTML:

```html
<button class="dn-btn dn-btn--atrament" type="button">Uruchom proces</button>
                ↑ blok    ↑ modyfikator
<div class="dn-karta">
  <div class="dn-karta-naglowek"><h3 class="dn-karta-tytul">Execution Monitor</h3></div>
  <div class="dn-karta-cialo">…</div>
</div>
        ↑ element — bez powtarzania klasy bloku
```

### 3.2. Dlaczego modyfikator nie działa samodzielnie

Reguły modyfikatorów są **niekompletne z założenia** — deklarują wyłącznie różnicę:

```css
.dn-btn--atrament {              /* nie ma tu wysokości, odstępów, kroju, promienia */
  background: var(--dn-atrament);
  border-color: var(--dn-atrament);
  color: var(--dn-atrament-tekst);
}
```

Konsekwencja praktyczna: `class="dn-btn--atrament"` bez `dn-btn` renderuje czarny prostokąt bez wysokości kontrolki i bez promienia. Recenzja arkusza sprawdza to punktem 14.2.

### 3.3. Trzy poziomy nazwy — przypadek karty środowiska

`.dn-karta-srodowiska-godlo` jest nazwą trójczłonową: blok `.dn-karta-srodowiska` + element `-godlo`. Reguła: **element dokleja się do pełnej nazwy bloku**, nie do jego skróconej formy.

```
.dn-karta-srodowiska          blok
.dn-karta-srodowiska-godlo    element  (godło środowiska)
.dn-karta-srodowiska-tytul    element  (TalkIn / WorkSpace / CodeStudio / MultitaskingAI)
.dn-karta-srodowiska-motto    element  (mono, --dn-tekst-3)
.dn-karta-srodowiska-opis     element  (max-width: 44ch)
```

### 3.4. Kolizja lekturalna: element czy osobny blok?

Konwencja ma jedną wadę wrodzoną, którą trzeba znać: `.dn-karta-tytul` i `.dn-karta-sesji` wyglądają identycznie, a znaczą co innego.

| Klasa | Faktycznie jest | Rozstrzygnięcie |
|---|---|---|
| `.dn-karta-tytul` | **elementem** `.dn-karta` | wymaga rodzica `.dn-karta` |
| `.dn-karta-cialo` | **elementem** `.dn-karta` | wymaga rodzica `.dn-karta` |
| `.dn-karta-naglowek` | **elementem** `.dn-karta` | wymaga rodzica `.dn-karta` |
| `.dn-karta-sesji` | **osobnym blokiem** | pas kart sesji powłoki środowiska; nie ma nic wspólnego z `.dn-karta` |
| `.dn-karta-srodowiska` | **osobnym blokiem** | karta środowiska — Centrum dowodzenia, strefa 1 |
| `.dn-karty-sesji` | **osobnym blokiem** (kontener) | liczba mnoga oznacza kontener pozycji |
| `.dn-btn-ikona` | **osobnym blokiem** | ma własną szerokość, wysokość, obrys i stany; nie dziedziczy z `.dn-btn` |

**Reguła rozstrzygająca:** jeżeli klasa deklaruje własne `display`, wymiary i tło — to blok, choćby jej nazwa wyglądała na element. Jeżeli deklaruje wyłącznie typografię i odstępy wewnątrz rodzica — to element.

**Reguła dla nowych nazw:** przed dodaniem `.dn-<x>-<y>` sprawdź, czy `<y>` może istnieć bez `<x>`. Jeśli tak — nadaj nazwę własną (`.dn-kafel`, a nie `.dn-karta-kafel`). Arkusz stosuje tę zasadę konsekwentnie: kafel komponentu własnego, karta środowiska, karta sesji, boczna nawigacja, listwa ustawień i Always On Display są **osobnymi blokami**, mimo pokrewieństwa wizualnego z kartą.

### 3.5. Wykaz odstępstw faktycznie obecnych w arkuszu

Recenzja musi je znać, żeby nie zgłaszać ich powtórnie jako błędów:

| Klasa | Odstępstwo | Powód zachowania |
|---|---|---|
| `.dn-radio` | powinna brzmieć `.dn-check--radio` (rozszerza `.dn-check`, wymaga jej w `class`) | nazwa krótka, zgodna z atrybutem `type="radio"`; zapisana w kontrakcie systemu projektowego jako wiążąca |
| `.dn-dane` | pełni **dwie role**: klasa pomocnicza fundamentu (krój mono + `tabular-nums`) i modyfikator komórki tabeli (`.dn-tabela td.dn-dane`) | rola jest ta sama — „to są dane maszynowe”; zdublowanie jest zamierzone |
| `.dn-liczba` | synonim `.dn-dane` w tej samej regule fundamentu (linia 92) | dwie nazwy dla jednej reguły; `.dn-liczba` czytelniejsza w tabelach liczbowych |
| `.dn-suwak` | w katalogu komponentów (v1.0) nazwa oznaczała przełącznik, w `komponenty.css` (v2.0) oznacza suwak zakresu | rozjazd między generacjami; wiążąca jest v2.0 — przełącznik to `.dn-przelacznik` |
| `.dn-naglowek` | klasa użytkowa w regule razem z `h1, h2, h3` | pozwala nadać brzmienie nagłówka elementowi, który semantycznie nagłówkiem nie jest |

### 3.6. Identyfikatory a język

Zasada kontrakt systemu projektowego: **komentarze i opisy po polsku, identyfikatory kodu i klasy po polsku wg istniejącej konwencji `.dn-*`, nazwy własne okien i modułów po angielsku.** W praktyce:

```css
/* Pas kart sesji powłoki środowiska — mechanika zakładek + wskaźnik pracy w tle */   ← komentarz PL
.dn-karty-sesji { … }                                                   ← klasa PL
```

```html
<span class="dn-plakietka dn-plakietka--rola">Coordinator</span>
                                              ↑ nazwa własna roli — angielska, bez tłumaczenia
```

Nazw własnych **nie tłumaczy się nigdy**: `Studio Editor`, `Workflow Builder`, `Chat Window`, `Queue Manager`, `Execution Monitor`, `Coordinator`, `Executor 1`, `Validator`, `Subagent Network`.

---

## 4. Specyficzność jako kontrakt — zakaz zagnieżdżania

### 4.1. Zmierzony sufit specyficzności

Pomiar na pełnym `komponenty.css` (1207 linii):

| Miara | Wynik | Znaczenie |
|---|---|---|
| Selektory identyfikatorowe (`#`) | **0** | żaden komponent nie da się przypiąć do jednego wystąpienia |
| Deklaracje `!important` | **0** | w `komponenty.css` i `fundament.css` nie ma ani jednej |
| Najwyższa specyficzność | **(0,3,0)** — `.dn-tooltip:hover .dn-tooltip-tresc` | trzy jednostki klasowe, zero identyfikatorów |
| Najczęstsza specyficzność stanu | **(0,2,0)** | `.dn-btn:hover`, `.dn-btn[aria-pressed='true']`, `.dn-krok--poprawny .dn-krok-znak` |
| Specyficzność bazy komponentu | **(0,1,0)** | `.dn-btn`, `.dn-karta`, `.dn-wpis` |

`!important` występuje w całym pakiecie **wyłącznie** w bloku `prefers-reduced-motion` w `zetony.css` (linie 442–447) — i to jest jego jedyne dopuszczone zastosowanie: awaryjne wygaszenie ruchu, które musi wygrać ze wszystkim, łącznie ze stylem lokalnym okna.

### 4.2. Drabina specyficzności — kontrakt

```
(0,1,0)  .dn-btn                              baza komponentu
         .dn-btn--sygnal                      wariant  ── ta sama waga, wygrywa kolejnością
            ▼
(0,2,0)  .dn-btn:hover                        stan wskaźnika
         .dn-btn--atrament:hover              stan wariantu ── musi być po wariancie
         .dn-btn[aria-pressed='true']         stan trwały
         .dn-btn.dn-btn--wybrany              stan trwały zapisany klasą
         .dn-krok--poprawny .dn-krok-znak     stan rodzica → dziecko
            ▼
(0,2,1)  .dn-btn[aria-busy='true']::after     stan + pseudoelement
            ▼
(0,2,2)  .dn-tabela tbody tr[aria-selected]   struktura natywna tabeli
            ▼
(0,3,0)  .dn-tooltip:hover .dn-tooltip-tresc  SUFIT — nie przekraczać
```

### 4.3. Dlaczego `.dn-btn.dn-btn--wybrany` ma podwojoną klasę

Zapis z linii 47–52:

```css
.dn-btn[aria-pressed='true'],
.dn-btn.dn-btn--wybrany { … }
```

Gdyby napisać samo `.dn-btn--wybrany` (0,1,0), reguła miałaby tę samą wagę co `.dn-btn--sygnal` i `.dn-btn--atrament`. O wyniku decydowałaby kolejność w pliku — czyli przypadek edycji. Podwojenie klasy podnosi wagę do (0,2,0) i **gwarantuje**, że stan „wybrany trwale” przykrywa barwę wariantu, niezależnie od tego, w którym miejscu pliku ktoś dopisze nowy wariant. To jest wzorzec do naśladowania przy każdym stanie kolidującym z wariantem.

### 4.4. Zakaz zagnieżdżania

**Zakazane** — reguła opisująca komponent przez kontekst rodzica:

```css
/* ZAKAZANE — komponent zmienia się od miejsca, w którym stoi */
.dn-boczna .dn-btn { border-radius: 0; }
.dn-modal .dn-pole-kontrolka { background: var(--dn-powierzchnia-2); }
.dn-karta > div > span { color: var(--dn-tekst-3); }
```

Trzy powody zakazu:

1. **Zerwanie kontraktu wizualnego.** Przycisk w bocznej nawigacji przestaje wyglądać jak przycisk w modalu. Makieta przestaje być wiarygodna, bo ten sam znacznik daje dwa różne wyniki.
2. **Nieograniczony przyrost specyficzności.** Każde zagnieżdżenie wymusza kolejne zagnieżdżenie u tego, kto chce to nadpisać. Po trzech iteracjach jedynym wyjściem jest `!important` — czyli koniec kontraktu.
3. **Niemożność podziału na pliki.** Rozdz. 9 dzieli bibliotekę na 14 arkuszy per komponent. Reguła `.dn-boczna .dn-btn` nie ma domu: nie należy ani do `boczna.css`, ani do `przycisk.css`. Podział staje się nierozstrzygalny.

**Dozwolone i faktycznie stosowane** — zagnieżdżenie w granicach **jednego** komponentu, czyli blok → własny element:

```css
.dn-wpis--czlowiek .dn-wpis-medalion { … }   /* wariant bloku → jego własny element */
.dn-krok--poprawny .dn-krok-znak { … }       /* jw. */
.dn-tooltip:hover .dn-tooltip-tresc { … }    /* stan bloku → jego własny element */
```

Kryterium: **oba człony selektora należą do tego samego komponentu i trafią do tego samego pliku.** Jeżeli tak — zagnieżdżenie jest opisem wewnętrznej mechaniki, a nie zależnością od kontekstu.

### 4.5. Selektory elementowe w komponentach

Komponent nie stylizuje gołego elementu. Pełny wykaz selektorów elementowych faktycznie obecnych w `komponenty.css` — **22 wystąpienia, wszystkie zakotwiczone klasą**:

| Postać | Wystąpienia | Uzasadnienie dopuszczenia |
|---|---|---|
| `textarea.dn-pole-kontrolka`, `select.dn-pole-kontrolka` | 2 | **kwalifikacja typu** — jedna klasa obsługuje trzy różne znaczniki natywne, a `textarea` potrzebuje `resize` i podwójnej wysokości, `select` — miejsca na strzałkę |
| `.dn-szukaj > svg`, `.dn-plakietka > svg`, `.dn-boczna-pozycja > svg`, `.dn-wpis-medalion > svg`, `.dn-krok-znak > svg`, `.dn-toast > svg`, `.dn-pusty-stan > svg`, `.dn-toast--{sukces,ostrzezenie,blad,informacja} > svg` | 11 | **ikona nie może nieść klasy** — wklejana inline z `zasoby/ikony/svg/`, bez modyfikacji pliku źródłowego; wymiar i barwa muszą przyjść z komponentu |
| `.dn-tabela th`, `.dn-tabela td`, `.dn-tabela tbody tr`, `tr:hover`, `tr[aria-selected]`, `td.dn-dane` | 6 | **struktura natywna tabeli** — `<tr>` i `<td>` generuje często silnik danych; oklasowanie każdej komórki jest niewykonalne |
| `.dn-pasek-logotyp small`, `.dn-aod-tresc strong` | 2 | **semantyka wewnątrz tekstu** — `<small>` niesie podpis `CONSOLE`, `<strong>` wyróżnia nazwę procesu w Always On Display |
| `.dn-szukaj > .dn-pole-kontrolka` | 1 | kompozycja dwóch komponentów w jednym opakowaniu (pole + ikona) |

**Zero gołych selektorów elementowych** (`div`, `p`, `button`, `ul` bez kotwicy klasowej) — sprawdzone. Reguła do zapamiętania: **selektor elementowy jest dozwolony wyłącznie jako kwalifikator albo jako dziecko bezpośrednie klasy `.dn-*`, i wyłącznie wtedy, gdy element nie może dostać klasy.**

---

## 5. Pełny wykaz reguł fundamentu z wartościami

Poniżej **komplet** `fundament.css` — wszystkie reguły, wszystkie deklaracje, wszystkie żetony. Kolumna „linie” odsyła do pliku źródłowego.

### 5.1. Wyzerowanie i podłoże

| Linie | Selektor | Deklaracje | Żeton / wartość | Cel |
|---|---|---|---|---|
| 9 | `*, *::before, *::after` | `box-sizing` | `border-box` | wymiar deklarowany = wymiar zajmowany; warunek działania siatki 4 px |
| 11 | `*` | `margin` | `0` | usunięcie marginesów przeglądarki; odstęp wyłącznie z `--dn-od-*` |
| 13 | `html` | `height` | `100%` | punkt odniesienia dla wysokości procentowych powłoki |
| 15–28 | `body` | `min-height` | `100dvh` | pełna wysokość widoku również na telefonie z paskiem systemowym |
| | | `background` | `var(--dn-tlo)` | `#F4F4F4` / `#0F0F0F` — nigdy czysta biel ani czerń |
| | | `color` | `var(--dn-tekst)` | `#181818` / `#ECECEC` |
| | | `font-family` | `var(--dn-ff-bazowa)` | IBM Plex Sans |
| | | `font-size` | `var(--dn-fs-base)` | 13 px (gęstość zwarta) |
| | | `line-height` | `var(--dn-lh-bazowy)` | 1.45 |
| | | `font-synthesis` | `none` | zakaz syntezy pogrubienia i kursywy — brakująca waga ma być widoczna, nie udawana |
| | | `text-rendering` | `optimizeLegibility` | podstawienia i kerning w małych stopniach |
| | | `-webkit-font-smoothing` | `antialiased` | wyrównanie grubości kroju na ciemnym tle |
| | | `transition` | `background-color`, `color` · `var(--dn-czas-2)` `var(--dn-ease)` | przełączenie motywu jako przejście 160 ms, nie przeskok |

### 5.2. Elementy bazowe

| Linie | Selektor | Deklaracje | Cel |
|---|---|---|---|
| 30 | `img, svg, video, canvas` | `display: block`, `max-width: 100%` | usunięcie szczeliny pod elementem liniowym; media nigdy nie rozpychają kontenera |
| 32 | `input, button, textarea, select` | `font: inherit`, `color: inherit` | **kluczowa reguła fundamentu** — kontrolki natywne przejmują krój i barwę dokumentu; bez niej każdy komponent formularza musiałby powtarzać deklarację kroju |
| 34 | `p, h1, h2, h3, h4` | `overflow-wrap: break-word` | długi identyfikator sesji albo ścieżka repozytorium nie rozrywa panelu |

### 5.3. Typografia nagłówków

| Linie | Selektor | Deklaracje | Żeton |
|---|---|---|---|
| 37–42 | `h1, h2, h3, .dn-naglowek` | `font-family` | `var(--dn-ff-naglowek)` — Space Grotesk |
| | | `font-weight` | `var(--dn-fw-polgruba)` — 600 |
| | | `letter-spacing` | `var(--dn-ls-naglowek)` — −0.01em |
| | | `line-height` | `var(--dn-lh-ciasny)` — 1.25 |

Komentarz w pliku brzmi: *„Space Grotesk niesie znaczenie: wejście i tożsamość”*. `h4` celowo **nie jest objęte** regułą — w kokpicie czwarty poziom pełni funkcję etykiety sekcji, a nie nagłówka, i przyjmuje wzorzec `.dn-etykieta-wersalikowa`.

### 5.4. Fokus, zaznaczenie, paski przewijania

| Linie | Selektor | Deklaracje | Żeton |
|---|---|---|---|
| 45–49 | `:focus-visible` | `outline` | `var(--dn-wym-fokus)` `solid` `var(--dn-fokus)` — 2 px, sygnał |
| | | `outline-offset` | `var(--dn-wym-fokus-odsuniecie)` — 2 px |
| | | `border-radius` | `var(--dn-r-xs)` — 3 px (wartość domyślna dla elementów bez promienia własnego) |
| 52 | `::selection` | `background` / `color` | `var(--dn-sygnal-tlo)` / `var(--dn-tekst)` |
| 53 | `:root[data-theme='dark'] ::selection` | `background` | `var(--dn-sygnal-obrys)` — w motywie ciemnym wersja półprzezroczysta |
| 56 | `*` | `scrollbar-width` / `scrollbar-color` | `thin` / `var(--dn-obrys-mocny)` `transparent` |
| 57 | `::-webkit-scrollbar` | `width` / `height` | `10px` / `10px` |
| 58–63 | `::-webkit-scrollbar-thumb` | `background` | `var(--dn-obrys-mocny)` |
| | | `border` | `3px solid transparent` — przezroczysta ramka zwęża uchwyt do 4 px |
| | | `border-radius` | `var(--dn-r-pill)` |
| | | `background-clip` | `content-box` — warunek działania sztuczki z ramką |
| 64 | `::-webkit-scrollbar-track` | `background` | `transparent` — tor niewidoczny; komentarz: *„instrument, nie ozdoba”* |

### 5.5. Odnośniki

| Linie | Selektor | Deklaracje |
|---|---|---|
| 67 | `a` | `color: var(--dn-sygnal)`, `text-decoration: none` |
| 68 | `a:hover` | `color: var(--dn-sygnal-mocny)`, `text-decoration: underline`, `text-underline-offset: 2px` |

Odnośnik jest jedynym miejscem, w którym sygnał pełni rolę barwy tekstu na powierzchni. `--dn-sygnal` przyjmuje `sygnal-600` na jasnym (5,5:1) i `sygnal-300` na ciemnym (8,0:1) — obie pary zmierzone w `kontrasty.json`.

### 5.6. Wzorce tekstowe — jedenaście klas

| Linie | Klasa | Deklaracje kluczowe | Zastosowanie w produkcie |
|---|---|---|---|
| 71–79 | `.dn-etykieta-wersalikowa` | `--dn-ff-bazowa`, `--dn-fs-xs` 11 px, `--dn-fw-polgruba`, `--dn-ls-wersaliki` 0.08em, `text-transform: uppercase`, `--dn-tekst-3` | nagłówki sekcji paneli, etykiety stref Centrum dowodzenia, podpisy grup w Oknie Konfiguracji |
| 82–89 | `.dn-etykieta-mono` | `--dn-ff-mono`, `--dn-fs-xs`, `--dn-fw-srednia`, `--dn-ls-mono-wersaliki` 0.14em, wersaliki, `--dn-tekst-3` | nagłówki kolumn, identyfikatory stref, znaczniki poziomów zasięgu |
| 92–95 | `.dn-dane`, `.dn-liczba` | `--dn-ff-mono`, `font-variant-numeric: tabular-nums` | każda liczba w kokpicie: liczniki kolejki, czasy, rozmiary, identyfikatory zadań |
| 98–110 | `.dn-kod` | `display: block`, padding `--dn-od-3`/`--dn-od-4`, obrys `--dn-obrys-subtelny`, `--dn-r-sm`, tło `--dn-powierzchnia-2`, `--dn-tekst-2`, mono `--dn-fs-sm`, `white-space: pre-wrap`, `overflow-x: auto` | Output Console, Logs Viewer, podgląd polityki efektywnej, ładunek narzędzia |
| 111–116 | `.dn-kod--wiersz` | `display: inline`, padding `1px var(--dn-od-1)`, `--dn-r-xs`, `white-space: nowrap` | nazwa żetonu albo klasy w zdaniu |
| 119–132 | `.dn-kbd` | `min-width: 20px`, obrys `--dn-obrys-mocny` z dolną krawędzią 2 px, `--dn-r-xs`, tło `--dn-powierzchnia`, mono `--dn-fs-xs` | skróty klawiaturowe w opisach i dymkach |
| 135–140 | `.dn-separator` | `height: 1px`, margines `--dn-od-4` pion, `background: var(--dn-obrys-subtelny)`, `border: none` | podział sekcji wewnątrz panelu bez nagłówka |
| 141–145 | `.dn-separator--pionowy` | `width: 1px`, `align-self: stretch`, margines `--dn-od-2` poziom | rozdzielenie grup na pasku górnym |
| 148–156 | `.dn-sr-only` | `position: absolute`, `1px × 1px`, `clip-path: inset(50%)`, `overflow: hidden`, `white-space: nowrap` | odnośnik „Przejdź do treści”, etykiety kontrolek ikonowych, opisy dla czytnika ekranu |
| 37 | `.dn-naglowek` | wspólna reguła z `h1, h2, h3` | brzmienie nagłówka na elemencie, który nagłówkiem semantycznie nie jest |

**Uwaga dostępnościowa do `.dn-sr-only`:** klasa używa `clip-path: inset(50%)` zamiast przestarzałego `clip: rect(…)` i **nie stosuje** `display: none` ani `visibility: hidden` — obie te deklaracje usunęłyby treść z drzewa dostępności, czyli zrobiłyby dokładnie odwrotnie do zamierzenia.

---

## 6. Wzorzec stanu

### 6.1. Zmierzony rozkład stanów w bibliotece

| Nośnik stanu | Wystąpień w `komponenty.css` | Rola |
|---|---|---|
| `:hover` | 19 | wskaźnik nad elementem |
| `:focus` | 7 | fokus pól (zastąpiony obrysem + poświatą) |
| `:checked` | 4 | stan kontrolek natywnych |
| `:active` | 3 | naciśnięcie |
| `:focus-visible` | 2 | fokus klawiaturowy komponentu (powtórzenie reguły globalnej) |
| `:focus-within` | 2 | fokus wewnątrz kontenera (`.dn-prompt`, `.dn-tooltip`) |
| `[aria-selected]` | 5 | wybór w zbiorze (zakładka, karta sesji, wiersz tabeli) |
| `[aria-current]` | 4 | bieżące położenie (boczna nawigacja, karta środowiska) |
| `[aria-pressed]` | 3 | przełącznik dwustanowy (przycisk, przycisk ikonowy) |
| `[aria-busy]` | 2 | praca w toku |
| `[aria-invalid]` | 2 | wartość nieprawidłowa |
| `[readonly]` | 1 | pole tylko do odczytu |
| modyfikator klasowy | 7 | stan trwały niereprezentowany przez ARIA (`--pracuje`, `--poprawny`, `--bledy`, `--wstrzymany`, `--wybrany`, `--wybrana`) |
| `[data-stan]` | **0** | **nie występuje** — patrz 6.5 |

### 6.2. Zasada wyboru nośnika — drabina decyzyjna

```
Czy stan ma odpowiednik w pseudoklasie natywnej (:hover/:active/:checked)?
   TAK ─► użyj pseudoklasy. Koniec.
   NIE ▼
Czy stan ma odpowiednik w atrybucie ARIA, który i tak trzeba ustawić
dla czytnika ekranu (aria-selected / aria-current / aria-pressed /
aria-busy / aria-invalid)?
   TAK ─► stylizuj po tym atrybucie. Jeden zapis obsługuje wygląd
          i dostępność; nie da się ich rozjechać. Koniec.
   NIE ▼
Czy stan jest trwały i wynika z danych (wynik kroku kolejki, klasa nadawcy)?
   TAK ─► modyfikator klasowy .dn-<komponent>--<stan>. Koniec.
   NIE ▼
Ostateczność: [data-stan="…"] — wyłącznie gdy żaden z powyższych
nie oddaje znaczenia (patrz 6.5).
```

Zysk z reguły „ARIA przed klasą” jest konkretny: element, który **wygląda** na wybrany, **jest** wybrany dla czytnika ekranu. Nie ma stanu widocznego, którego nie widzi technologia wspomagająca — bo to ten sam atrybut.

### 6.3. Wzorce zapisu — z arkusza

**Najechanie — zmienia tło, nie rozmiar:**

```css
.dn-btn:hover        { background: var(--dn-hover); }
.dn-btn-ikona:hover  { background: var(--dn-hover); color: var(--dn-tekst); }
.dn-zakladka:hover   { color: var(--dn-tekst); background: var(--dn-hover); }
```

`--dn-hover` to `rgba(10,10,10,.05)` na jasnym i `rgba(255,255,255,.06)` na ciemnym — nakładka półprzezroczysta, nie druga barwa. Dzięki temu działa na każdej powierzchni (`--dn-powierzchnia`, `--dn-panel`, `--dn-powierzchnia-2`) bez wariantów.

**Naciśnięcie — jeden piksel w dół:**

```css
.dn-btn:active       { transform: translateY(1px); }
.dn-btn-ikona:active { transform: translateY(1px); }
.dn-karta--klikalna:active { transform: translateY(0); }   /* powrót z uniesienia */
```

Karta klikalna unosi się przy najechaniu o −1 px, więc jej naciśnięcie to powrót do zera. Ruch jest zawsze o **1 px** i zawsze na osi pionowej.

**Stan trwały — ARIA:**

```css
.dn-zakladka[aria-selected='true'] { color: var(--dn-tekst); }
.dn-zakladka[aria-selected='true']::after { … height: 2px; background: var(--dn-kropka); }

.dn-boczna-pozycja[aria-current='page'] { background: var(--dn-sygnal-tlo); color: var(--dn-tekst); }
.dn-boczna-pozycja[aria-current='page']::before { … width: var(--dn-wym-wstega); height: 16px; }

.dn-tabela tbody tr[aria-selected='true'] { background: var(--dn-sygnal-tlo); }
```

Zwróć uwagę: **stan nigdy nie jest samą barwą.** Zakładka wybrana dostaje podkreślenie 2 px, pozycja nawigacji — kreskę sygnału po lewej, karta sesji — wstęgę górną. Barwa jest wzmocnieniem kształtu, nie jedynym nośnikiem. To realizacja zasady z kontraktu systemu projektowego.

**Praca w toku — Zasada zero blokad:**

```css
.dn-btn[aria-busy='true'] { cursor: progress; }
.dn-btn[aria-busy='true']::after { … animation: dn-obrot 0.8s linear infinite; }
```

Przycisk pracujący **pozostaje klikalny**. Zmienia się kursor i dochodzi wskaźnik, nie znika sprawczość. Nagłówek arkusza mówi to wprost: *„Zasada zero blokad: żaden wariant nie odbiera klikalności”*.

**Błąd:**

```css
.dn-btn[aria-invalid='true']            { border-color: var(--dn-blad-obrys); color: var(--dn-blad-tekst); }
.dn-pole-kontrolka[aria-invalid='true'] { border-color: var(--dn-blad-tekst); }
```

Pole z błędem pozostaje edytowalne. Obok stoi `.dn-pole-blad` z ikoną — znowu: nie sam kolor.

### 6.4. Stan wariantu — kolejność w pliku

Reguła stanu wariantu musi stać **po** regule wariantu, bo obie mają wagę (0,2,0) i (0,1,0):

```css
.dn-btn--atrament        { background: var(--dn-atrament); … }        /* (0,1,0) */
.dn-btn--atrament:hover  { background: var(--dn-atrament-hover); … }  /* (0,2,0) */
```

Gdyby wariant nie miał własnej reguły `:hover`, przejąłby `.dn-btn:hover { background: var(--dn-hover) }` — czyli półprzezroczystą nakładkę zamiast ciemniejszego atramentu. Dlatego **każdy wariant zmieniający tło musi zadeklarować własne najechanie.** Sprawdza to punkt 14.5 listy kontrolnej. W arkuszu robią to: `--atrament`, `--sygnal`, `--duch`, `--niebezpieczny`.

### 6.5. `[data-stan]` — dlaczego nie występuje

W obu arkuszach nie ma ani jednego selektora `[data-stan]`. Atrybuty `data-*` są w tym systemie zarezerwowane dla **trybów** i **zaczepów skryptu**, nie dla stanów komponentu:

| Atrybut | Miejsce | Rola |
|---|---|---|
| `data-theme` | `<html>` | tryb motywu — przełącza cały zestaw żetonów semantycznych |
| `data-gestosc` | `<html>` | tryb gęstości (`przestronna`) — przelicza wymiary |
| `data-przelacz-motyw` | przycisk na pasku | zaczep dla `wspolne.js`; **bez reguły CSS** |
| `data-motyw-ikona` | opakowanie ikony | zaczep dla `wspolne.js`; **bez reguły CSS** |

**Rozstrzygnięcie (decyzja D-A5-04):** `[data-stan]` jest dopuszczalny wyłącznie dla stanu, który jednocześnie (a) nie ma odpowiednika w ARIA, (b) nie jest trwały (więc modyfikator klasowy byłby mylący), (c) zmienia się wyłącznie po stronie skryptu. W obecnym zakresie okien taki przypadek nie występuje — cztery wyjścia kroku kolejki (`--pracuje`, `--poprawny`, `--bledy`, `--wstrzymany`) są trwałe i wynikają z danych, więc poprawnie zapisano je modyfikatorem. Zapis pozostaje w konwencji jako furtka, nie jako wzorzec.

---

## 7. Wzorzec fokusu globalnego

### 7.1. Reguła

```css
:focus-visible {
  outline: var(--dn-wym-fokus) solid var(--dn-fokus);   /* 2 px, błękit sygnałowy */
  outline-offset: var(--dn-wym-fokus-odsuniecie);       /* 2 px odsunięcia */
  border-radius: var(--dn-r-xs);                        /* 3 px — wartość domyślna */
}
```

Wartości pochodzą z kontraktu systemu projektowego („fokus 2 px + odsunięcie 2 px”) i rozdz. 9 („Pierścień fokusu: 2 px + odsunięcie 2 px, barwa `--dn-fokus`”). opracowanie o przekazaniu, ruchu i dostępności wymienia pierścień fokusu wśród rzeczy, **których nie wolno zmienić po cichu**.

### 7.2. Sześć powodów, dla których fokus jest globalny, a nie per komponent

| # | Powód | Konsekwencja praktyczna |
|---|---|---|
| 1 | **Zasięg poza biblioteką** | Fokus dostaje każdy element fokusowalny: odnośnik w treści, `<summary>` panelu „O tym opracowaniu”, wiersz listy z `tabindex`, kontrolka natywna bez klasy `.dn-*`. Reguła per komponent obsłużyłaby wyłącznie bibliotekę. |
| 2 | **Nowy komponent dziedziczy za darmo** | Klasa dopisana jutro do `komponenty.css` ma poprawny fokus, zanim ktokolwiek o nim pomyśli. Nie da się „zapomnieć fokusu” w nowym komponencie. |
| 3 | **Jedna wartość, jedno miejsce zmiany** | Zmiana grubości pierścienia to jedna deklaracja. Rozproszona wersja wymagałaby edycji kilkudziesięciu reguł i gwarantowała rozjazd. |
| 4 | **Symetria z `prefers-reduced-motion`** | kontrakt systemu projektowego stawia obie sprawy obok siebie: ruch ograniczony **globalnie w żetonach**, fokus **globalnie w fundamencie** — nie per komponent. Ten sam wzorzec architektoniczny. |
| 5 | **`:focus-visible`, nie `:focus`** | Pierścień pojawia się przy nawigacji klawiaturą, nie po kliknięciu myszą. Bez tego rozróżnienia projektanci zaczynają usuwać obrys, żeby „nie brzydził” — i kokpit traci dostępność. |
| 6 | **Audyt w jednym miejscu** | Recenzja dostępności sprawdza jedną regułę zamiast przeglądać cały arkusz. |

### 7.3. Dwa powtórzenia reguły w komponentach — czym są

`komponenty.css` powtarza reguły fokusu dokładnie dwa razy:

```css
.dn-btn:focus-visible       { outline: var(--dn-wym-fokus) solid var(--dn-fokus);
                              outline-offset: var(--dn-wym-fokus-odsuniecie); }
.dn-btn-ikona:focus-visible { outline: var(--dn-wym-fokus) solid var(--dn-fokus);
                              outline-offset: var(--dn-wym-fokus-odsuniecie); }
```

Powtórzenia są **funkcjonalnie zbędne** — reguła globalna daje ten sam wynik przy tych samych żetonach. Zachowano je świadomie, bo rozdz. 9 dzieli bibliotekę na osobne pliki: `przycisk.css` ma pozostać **przenośny** — czytelny i poprawny również wtedy, gdy trafi do przeglądu w oderwaniu od `fundament.css`. Odnotowane jako decyzja D-A5-05. Powielenia **nie wolno rozszerzać** na kolejne komponenty; jeśli ktoś je usunie, wygląd nie zmieni się o piksel.

### 7.4. Fokus pól — dlaczego `outline: none` jest tu dopuszczalne

Trzy miejsca gaszą pierścień:

```css
.dn-pole-kontrolka:focus { outline: none; border-color: var(--dn-fokus); box-shadow: var(--dn-cien-sygnal); }
.dn-pasek-szukaj:focus   { outline: none; border-color: var(--dn-fokus); box-shadow: var(--dn-cien-sygnal); }
.dn-prompt-obszar:focus  { outline: none; }
```

Warunki dopuszczalności — spełnione we wszystkich trzech:

1. **Zastąpienie, nie usunięcie.** Pierwsze dwa zamieniają obrys na parę „obrys kontrolki w barwie fokusu + poświata `--dn-cien-sygnal` 3 px”. To wskazanie **silniejsze** wizualnie niż pierścień, bo obejmuje cały kształt pola.
2. **Kontener przejmuje wskazanie.** `.dn-prompt-obszar` gasi obrys, ale rodzic reaguje: `.dn-prompt:focus-within { border-color: var(--dn-fokus); box-shadow: var(--dn-cien-sygnal); }`. Podświetla się całe pole promptu wraz z grotem i przybornikiem — to trafniejsze wskazanie niż obrys wokół samego obszaru tekstu.
3. **Selektor to `:focus`, nie `:focus-visible`.** Pole musi pokazywać aktywność również po kliknięciu myszą — użytkownik musi wiedzieć, gdzie trafi wpisywany tekst.

**Zakaz bezwzględny:** `outline: none` bez zastąpienia wskazaniem o równej albo większej sile. Recenzja traktuje samotne `outline: none` jako błąd blokujący scalenie (punkt 14.9).

### 7.5. Barwa pierścienia a rama kokpitu

`--dn-fokus` przyjmuje `sygnal-500` na jasnym i `sygnal-400` na ciemnym. Pasek górny jest atramentowy **w obu motywach**, więc w motywie jasnym pierścień pada na tło `szary-925` — a to ciemniejsza para niż zakładana. Dlatego wariant `.dn-btn-ikona--na-ramie` w stanie wciśniętym używa `--dn-sygnal-300` (jaśniejszego), a nie `--dn-sygnal`. Sam pierścień fokusu pozostaje wspólny: `sygnal-500` na `szary-925` daje kontrast wystarczający dla elementu nietekstowego.

---

## 8. Wzorzec przejścia i animacji

### 8.1. Zasada nadrzędna

Kierunek projektowy ustawia `INTENSYWNOSC_RUCHU` na **3/10**: mikroprzejścia 100–220 ms plus **jeden ruch znaczący** — tętno pracy w tle. Ruch zawsze niesie informację o stanie systemu; ruch dekoracyjny jest zakazany.

### 8.2. Cztery czasy i ich zakresy

| Żeton | Wartość | Zakres | Wystąpień w `komponenty.css` |
|---|---|---|---|
| `--dn-czas-1` | 0.10 s | mikroreakcje: naciśnięcie, tło wiersza tabeli, zaznaczenie pola wyboru | 6 |
| `--dn-czas-2` | 0.16 s | **czas domyślny** — barwy, obrysy, cienie, przełączenie motywu | 36 |
| `--dn-czas-3` | 0.22 s | wejście warstw: modal, powiadomienie, pasek postępu | 3 |
| `--dn-czas-tetno` | 2.40 s | tętno kropki sygnału — **jedyny ruch ciągły** | 3 |

Krzywa jest jedna dla całego systemu: `--dn-ease: cubic-bezier(0.2, 0, 0, 1)` — szybki start, długie wyhamowanie. Dobór czasu: **czas-2, chyba że masz powód.** Powodem dla czas-1 jest reakcja, która ma sprawiać wrażenie natychmiastowej; dla czas-3 — pojawienie się nowej warstwy, które ma być zauważone.

### 8.3. Zasada wyliczania właściwości

**Nigdy `transition: all`.** Arkusz wylicza właściwości co do jednej — 19 deklaracji `transition`, zero z `all`.

```css
.dn-btn {
  transition:
    background-color var(--dn-czas-2) var(--dn-ease),
    border-color     var(--dn-czas-2) var(--dn-ease),
    color            var(--dn-czas-2) var(--dn-ease),
    box-shadow       var(--dn-czas-2) var(--dn-ease),
    transform        var(--dn-czas-1) var(--dn-ease);
}
```

Trzy powody:

| # | Powód | Skutek `transition: all` |
|---|---|---|
| 1 | **Kontrola czasu per właściwość** | `all` wymusza jeden czas; przycisk straciłby szybkie 100 ms na `transform` przy 160 ms na barwach — naciśnięcie zrobiłoby się ospałe |
| 2 | **Ochrona przed animowaniem układu** | `all` animuje także `width`, `height`, `padding`, `top`; każda z nich wywołuje przeliczenie układu w każdej klatce. W tabeli o kilkuset wierszach to widoczne zacinanie |
| 3 | **Czytelność zamiaru** | wyliczona lista jest dokumentacją: widać, co się rusza i dlaczego |

### 8.4. Właściwości dopuszczone do przejść

| Kategoria | Właściwości | Koszt | Uwagi |
|---|---|---|---|
| **Kompozytowe** | `transform`, `opacity` | najniższy — bez przeliczania układu i malowania | preferowane; `translateY(1px)` przy naciśnięciu, `translateY(-1px)`/`(-2px)` przy uniesieniu karty |
| **Malowanie** | `background-color`, `border-color`, `color`, `box-shadow` | średni — bez przeliczania układu | trzon systemu; wszystkie stany barwne |
| **Widoczność** | `visibility` | niski | wyłącznie razem z `opacity` (dymek), żeby po zaniku nie przechwytywał wskaźnika |
| **Warunkowo** | `width`, `left` | wysoki | **dwa dopuszczone wyjątki**, patrz 8.5 |

### 8.5. Dwa uzasadnione wyjątki

```css
.dn-postep-wartosc { transition: width var(--dn-czas-3) var(--dn-ease); }
.dn-przelacznik::before { transition: left var(--dn-czas-2) var(--dn-ease), background-color …; }
```

- **`width` paska postępu** — jest treścią komunikatu („ile zrobiono”), a nie efektem. Wariant z `transform: scaleX()` zniekształciłby zaokrąglenie końcówki toru.
- **`left` suwaka przełącznika** — element ma 14 px i jest jedynym dzieckiem toru; koszt przeliczenia jest pomijalny, a `left` daje precyzyjne trafienie w `calc(100% - var(--dn-wym-przelacznik-wys) + 4px)` przy zmiennej szerokości toru w trybie dotykowym.

### 8.6. Trzy animacje w systemie

| Nazwa | Definicja | Użycie | Czas |
|---|---|---|---|
| `dn-tetno` | pierścień `box-shadow` 0 → 5 px, zanikający | `.dn-kropka--tetno`, `.dn-wpis--pracuje .dn-wpis-nadawca::after`, `.dn-aod-rdzen::after` | `--dn-czas-tetno` (2,4 s), `infinite` |
| `dn-wejscie` | `opacity` 0 → 1 z `translateY(8px) scale(0.98)` | `.dn-modal[open]`, `.dn-toast` | `--dn-czas-3` (0,22 s), raz |
| `dn-obrot` | `rotate(360deg)` | `.dn-spinner`, `.dn-btn[aria-busy='true']::after` | **`0.8s`** — literał, patrz niżej |

**Odnotowane odstępstwo:** `dn-obrot` jest jedyną animacją z czasem wpisanym wprost (`0.8s`, dwa wystąpienia: linie 65 i 1101), mimo deklaracji „zero wartości zaszytych” w nagłówku arkusza. Nie ma żetonu `--dn-czas-obrot`. Konsekwencja praktyczna jest zerowa — blok `prefers-reduced-motion` w `zetony.css` gasi wszystkie animacje przez `animation-duration: 0.01ms !important` niezależnie od tego, czy czas pochodzi z żetonu. Pozostawiono bez zmian; zgłoszone jako obserwacja przeglądu (rozdz. 14.13).

### 8.7. Tętno — jedyny ruch ciągły

```css
.dn-kropka--tetno { animation: dn-tetno var(--dn-czas-tetno) var(--dn-ease) infinite; }
@keyframes dn-tetno {
  0%, 100% { box-shadow: 0 0 0 0 var(--dn-fokus-cien); }
  40%      { box-shadow: 0 0 0 5px transparent; }
}
@media (prefers-reduced-motion: reduce) {
  .dn-kropka--tetno { box-shadow: 0 0 0 2px var(--dn-fokus-cien); }
}
```

Trzy cechy warte zapamiętania:

1. **Animowany jest `box-shadow`, nie `transform: scale()`** — kropka ma 6 px; skalowanie tak małego elementu daje rozmyte krawędzie. Rozchodzący się pierścień czyta się jako sygnał, a nie jako drganie.
2. **Szczyt na 40 %, nie na 50 %** — pierścień rozchodzi się szybciej niż wraca. To rytm pulsu, nie oddechu.
3. **Zastąpienie, nie wygaszenie.** Przy `prefers-reduced-motion` kropka dostaje **statyczny pierścień 2 px** — pozostaje odróżnialna od kropek stanu (`--sukces`, `--ostrzezenie`, `--blad`, `--neutralna`). Gdyby animację po prostu wyłączono, informacja „tu biegnie praca” zniknęłaby dla użytkownika, który ograniczył ruch. Wzorzec do naśladowania: **ruch niosący informację zastępuje się formą statyczną niosącą tę samą informację.**

---

## 9. Podział `komponenty.css` na pliki per komponent

### 9.1. Podstawa

Regulamin budowy przywołany w opracowaniu o przekazaniu, ruchu i dostępności: **jedna odpowiedzialność = jeden plik, arkusz ≤ 300 linii, komponenty sięgają wyłącznie po żetony semantyczne.** `komponenty.css` ma 1207 linii — czterokrotność progu. opracowanie o przekazaniu, ruchu i dostępności wskazuje cel `komponenty/` i wylicza czternaście plików.

### 9.2. Zmierzone granice sekcji w pliku źródłowym

Podział nie jest szacunkiem — sekcje są w pliku wydzielone komentarzami blokowymi. Sumy zweryfikowane: 1189 linii treści + 18 linii pustych między sekcjami = 1207. ✔

| Sekcja w `komponenty.css` | Linie | Objętość |
|---|---|---|
| Nagłówek arkusza | 1–9 | 9 |
| PRZYCISK | 11–164 | 154 |
| POLE FORMULARZA | 166–233 | 68 |
| WYBÓR: pole wyboru · przełącznik · radio | 235–333 | 99 |
| KROPKA SYGNAŁU | 335–360 | 26 |
| PLAKIETKA | 362–393 | 32 |
| KARTA (+ karta środowiska, kafel) | 395–523 | 129 |
| TABELA | 525–559 | 35 |
| ZAKŁADKI · KARTY SESJI | 561–637 | 77 |
| RAMA KOKPITU (pasek + listwa) | 639–718 | 80 |
| BOCZNA NAWIGACJA | 720–771 | 52 |
| WPIS OKNA KOMUNIKACJI | 773–862 | 90 |
| MONITOR WYKONANIA · KOLEJKA | 864–945 | 82 |
| MODAL · NAKŁADKA | 947–987 | 41 |
| POWIADOMIENIE | 989–1024 | 36 |
| DYMEK OBJAŚNIENIA | 1026–1058 | 33 |
| AWATAR · SPINNER · PUSTY STAN | 1060–1121 | 62 |
| POLE PROMPTU · PRZYBORNIK | 1123–1169 | 47 |
| ALWAYS ON DISPLAY | 1171–1207 | 37 |

### 9.3. Tabela podziału docelowego

Szacunek objętości = linie sekcji + 8 linii nagłówka komentarzowego pliku (tytuł, zależność od żetonów, zakres, data).

| # | Plik docelowy | Klasy `.dn-*` | Sekcje źródłowe | Linie | Szacunek |
|---|---|---|---|---|---|
| 1 | `przycisk.css` | `.dn-btn` · `--wybrany --atrament --sygnal --zarys --duch --niebezpieczny --sm --lg` · `.dn-btn-ikona` · `--na-ramie` | PRZYCISK | 154 | **~162** |
| 2 | `pole.css` | `.dn-pole` · `-etykieta -kontrolka -opis -blad` · `.dn-szukaj` | POLE FORMULARZA | 68 | **~76** |
| 3 | `wybor.css` | `.dn-wybor` · `.dn-check` · `.dn-radio` · `.dn-przelacznik` · `.dn-suwak` | WYBÓR | 99 | **~107** |
| 4 | `plakietka.css` | `.dn-plakietka` · `--sukces --ostrzezenie --blad --informacja --sygnal --rola` · `.dn-kropka` · `--tetno --sukces --ostrzezenie --blad --neutralna` · `@keyframes dn-tetno` | PLAKIETKA + KROPKA | 58 | **~66** |
| 5 | `karta.css` | `.dn-karta` · `--klikalna --wybrana` · `-naglowek -tytul -cialo` · `.dn-karta-srodowiska` + 4 elementy · `.dn-kafel` + 3 elementy | KARTA | 129 | **~137** |
| 6 | `tabela.css` | `.dn-tabela` (+ `th`, `td`, `tbody tr`, `td.dn-dane`) | TABELA | 35 | **~43** |
| 7 | `zakladki.css` | `.dn-zakladki` · `.dn-zakladka` · `.dn-karty-sesji` · `.dn-karta-sesji` | ZAKŁADKI · KARTY SESJI | 77 | **~85** |
| 8 | `pasek.css` | `.dn-pasek` · `-godlo -logotyp -szukaj -prawa` · `.dn-listwa` · `-pozycja` | RAMA KOKPITU | 80 | **~88** |
| 9 | `boczna.css` | `.dn-boczna` · `-naglowek -pozycja` | BOCZNA NAWIGACJA | 52 | **~60** |
| 10 | `wpis.css` | `.dn-wpis` · `--czlowiek --inteligencja --system --pracuje` · `-medalion -tozsamosc -nadawca -godzina -tresc` | WPIS | 90 | **~98** |
| 11 | `postep.css` | `.dn-postep` · `-tor -wartosc -etykieta` · `.dn-kolejka` · `.dn-krok` · `--pracuje --poprawny --bledy --wstrzymany` · `-znak -meta` | MONITOR · KOLEJKA | 82 | **~90** |
| 12 | `nakladka.css` | `.dn-modal` · `::backdrop` · `-naglowek -tytul -cialo -stopka` · `@keyframes dn-wejscie` | MODAL | 41 | **~49** |
| 13 | `powiadomienie.css` | `.dn-toasty` · `.dn-toast` · `--sukces --ostrzezenie --blad --informacja` · `-tytul -tresc` | POWIADOMIENIE | 36 | **~44** |
| 14 | `drobne.css` | `.dn-tooltip` · `-tresc` · `.dn-awatar` · `--sm --lg --kwadrat --inteligencja` · `-stan` · `.dn-spinner` · `@keyframes dn-obrot` · `.dn-pusty-stan` · `-tytul -opis` · `.dn-prompt` · `-grot -obszar` · `.dn-przybornik` · `.dn-aod` · `-rdzen -tresc` | DYMEK + AWATAR/SPINNER/PUSTY + PROMPT + AOD | 179 | **~187** |
| | **Razem** | 119 tokenów klasowych | 18 sekcji | 1180 | **~1292** |

**Wszystkie czternaście plików mieści się pod progiem 300 linii.** Najbliżej progu: `drobne.css` (~187) i `przycisk.css` (~162) — oba z zapasem ponad 100 linii.

### 9.4. Gdzie trafiają definicje `@keyframes`

Klatki kluczowe mieszkają w pliku komponentu, który ich używa jako pierwszy — nie w osobnym pliku „animacje”:

| Klatki | Plik | Uzasadnienie |
|---|---|---|
| `dn-tetno` | `plakietka.css` | kropka sygnału jest źródłem wzorca; `wpis.css` i `drobne.css` (Always On Display) korzystają z tej samej nazwy |
| `dn-wejscie` | `nakladka.css` | modal jest pierwszym użytkownikiem; `powiadomienie.css` korzysta z tej samej nazwy |
| `dn-obrot` | `drobne.css` | spinner jest źródłem; `przycisk.css` korzysta przy `[aria-busy]` |

**Konsekwencja wiążąca dla arkusza spinającego:** nazwy `@keyframes` są globalne w dokumencie, więc plik korzystający z cudzych klatek **musi być importowany po** pliku, który je definiuje. Kolejność w arkuszu spinającym `komponenty.css`:

```
plakietka.css   ─► definiuje dn-tetno
nakladka.css    ─► definiuje dn-wejscie
drobne.css      ─► definiuje dn-obrot
   ▼ dopiero potem pliki korzystające
wpis.css        ─► używa dn-tetno
powiadomienie   ─► używa dn-wejscie
przycisk.css    ─► używa dn-obrot
```

Zależność jest odwrotna do kolejności alfabetycznej i do kolejności w tabeli 9.3 — **arkusz spinający musi ją zapisać jawnie, z komentarzem.** To najczęstsza pułapka podziału: po rozbiciu pliku spinner przestaje się obracać, bo `@keyframes dn-obrot` ładuje się po `przycisk.css`. Odnotowane jako decyzja D-A5-06.

### 9.5. Trzy zasady podziału

1. **Plik ma jedną nazwę komponentu i wszystkie jego warianty.** `przycisk.css` zawiera osiem modyfikatorów `.dn-btn` i dwa `.dn-btn-ikona` — bo wariant bez bazy jest niekompletny (rozdz. 3.2).
2. **Komponenty pokrewne wizualnie i zawsze współwystępujące mogą dzielić plik.** `plakietka.css` mieści plakietkę i kropkę, bo kropka jest bezetykietową odmianą znacznika stanu; `zakladki.css` — zakładki i karty sesji, bo dzielą mechanikę `[aria-selected]`.
3. **`drobne.css` to nie śmietnik.** Mieści cztery komponenty, których osobne pliki miałyby po 33–62 linie: dymek, tożsamość (awatar + spinner + pusty stan), prompt, Always On Display. Kryterium: komponent bez własnych wariantów i bez własnej mechaniki stanu. Gdy któryś urośnie o warianty — wychodzi do własnego pliku.

---

## 10. Wzorce użytkowe (utility) — co istnieje, czego się nie tworzy

### 10.1. Jedenaście klas użytkowych, które istnieją

Wszystkie mieszkają w `fundament.css` (rozdz. 5.6). Dzielą się na trzy rodziny:

| Rodzina | Klasy | Wspólna cecha |
|---|---|---|
| **Brzmienie tekstu** | `.dn-naglowek`, `.dn-etykieta-wersalikowa`, `.dn-etykieta-mono`, `.dn-dane`, `.dn-liczba` | nadają rolę typograficzną dowolnemu elementowi; zero wpływu na układ |
| **Treść techniczna** | `.dn-kod`, `.dn-kod--wiersz`, `.dn-kbd` | prezentują treść maszynową; mają własne tło i obrys, ale nie własny układ |
| **Struktura minimalna** | `.dn-separator`, `.dn-separator--pionowy`, `.dn-sr-only` | jedno zadanie strukturalne, którego nie da się przypisać do żadnego komponentu |

**Test przyjęcia nowej klasy użytkowej — musi spełniać wszystkie cztery warunki:**

1. odpowiada **roli semantycznej**, nie wartości („to są dane maszynowe”, a nie „to jest 12 px”),
2. występuje w **co najmniej trzech różnych oknach** produktu,
3. **nie da się jej przypisać** do konkretnego komponentu bez sztuczności,
4. **nie parametryzuje się** — nie ma i nie będzie mieć rodzeństwa `-1`, `-2`, `-3`.

### 10.2. Czego się nie tworzy — i dlaczego

| Rodzina odrzucona | Przykłady zakazane | Uzasadnienie |
|---|---|---|
| **Odstępy** | `.mt-2`, `.p-4`, `.gap-3`, `.dn-od-4` jako klasa | Odstęp jest **cechą komponentu**, nie dokumentu. Kokpit ma jeden rytm — 12–16 px wewnątrz paneli — zapisany w `--dn-odstep-panel` i `--dn-odstep-sekcji`. Klasa odstępowa pozwala go ominąć w znaczniku, gdzie recenzja stylu jej nie widzi. Dodatkowo `data-gestosc="przestronna"` przelicza wymiary przez żetony — klasa `.mt-2` na to nie zareaguje i gęstość przestronna rozjedzie się w tych właśnie miejscach. |
| **Barwy** | `.text-blad`, `.bg-panel`, `.obrys-mocny` | Barwa niesie znaczenie stanu, a stan **nigdy nie jest samym kolorem** (kontrakt systemu projektowego). Klasa barwna zachęca dokładnie do naruszenia tej zasady: pomaluj na czerwono zamiast dodać ikonę i etykietę. Właściwym narzędziem jest wariant komponentu (`.dn-plakietka--blad` niesie tło, obrys, barwę **i** miejsce na ikonę 12 px). |
| **Układ** | `.flex`, `.grid`, `.items-center`, `.justify-between` | Układ wewnętrzny należy do komponentu, układ zewnętrzny — do arkusza lokalnego okna. Klasy układowe rozmywają tę granicę: znacznik zaczyna opisywać wygląd, a arkusz przestaje być kompletnym opisem systemu. Przy 12 modułach i kilkudziesięciu oknach oznacza to, że makieta nie jest już źródłem prawdy. |
| **Typografia parametryczna** | `.fs-sm`, `.fw-600`, `.lh-tight` | Duplikat żetonu z gorszym interfejsem. `.dn-etykieta-wersalikowa` mówi, **czym** jest tekst; `.fs-xs` mówi tylko, **jak duży** jest — i pozwala zbudować etykietę o właściwym rozmiarze, ale bez wersalików i światła. |
| **Widoczność i wymiar** | `.hidden`, `.w-full`, `.truncate` | `.hidden` bywa protezą blokady (Zasada zero blokad: element niegotowy nie znika, tylko komunikuje). Szerokość i przycięcie to sprawa kontekstu, w którym komponent stoi. |

### 10.3. Bilans

Klas użytkowych jest **jedenaście na 119 tokenów klasowych** — poniżej dziesięciu procent, wszystkie w warstwie fundamentu, żadna w bibliotece komponentów. Taki bilans jest miarą zdrowia systemu: rosnąca liczba klas użytkowych to sygnał, że biblioteka nie pokrywa faktycznych potrzeb i autorzy okien obchodzą ją bokiem.

---

## 11. Obsługa motywu w CSS

### 11.1. Zasada jednego zapisu

**Komponent nie wie, jaki motyw jest włączony.** Zapisuje się go **raz**; różnicę niesie żeton.

Dowód: `komponenty.css` zawiera **zero** selektorów `[data-theme]` — sprawdzone na 1207 liniach.

```css
/* POPRAWNIE — jeden zapis, dwa motywy */
.dn-karta {
  border: 1px solid var(--dn-obrys);
  background: var(--dn-powierzchnia);
  box-shadow: var(--dn-cien-1);
}

/* ZAKAZANE — duplikacja reguły */
.dn-karta { background: #FFFFFF; }
[data-theme='dark'] .dn-karta { background: #181818; }
```

Wersja zakazana ma cztery wady: podwaja objętość arkusza, podnosi specyficzność jednego motywu ponad drugi (`[data-theme='dark'] .dn-karta` = (0,2,0) kontra (0,1,0) — motyw ciemny staje się „mocniejszy”), wpisuje wartości szesnastkowe wprost i rozjeżdża się przy pierwszej zmianie palety, bo nikt nie pamięta o drugiej regule.

### 11.2. Trzy kategorie żetonów według zachowania przy przełączeniu

| Kategoria | Zachowanie | Przykłady | Kiedy sięgać |
|---|---|---|---|
| **Przełączalne (semantyczne)** | zmieniają wartość razem z `data-theme` | `--dn-tlo`, `--dn-powierzchnia`, `--dn-panel`, `--dn-tekst`, `--dn-obrys`, `--dn-sygnal`, `--dn-cien-*`, wszystkie żetony stanów | **domyślnie zawsze** — to 95 % przypadków |
| **Stałe ramy** | identyczne w obu motywach | `--dn-rama`, `--dn-rama-tekst`, `--dn-rama-tekst-2`, `--dn-rama-hover`, `--dn-rama-obrys` | wyłącznie w obrębie paska górnego i w dymku (`.dn-tooltip-tresc` używa `--dn-rama`) |
| **Prymitywy stałe** | wartość absolutna, świadomie niezmienna | `--dn-szary-0`, `--dn-szary-100` | **tylko** gdy powierzchnia pod spodem sama nie przełącza się z motywem |

### 11.3. Pięć dopuszczonych użyć prymitywu w komponentach

Zasada kontrakt systemu projektowego mówi: *„komponent sięga wyłącznie po żetony SEMANTYCZNE, nigdy po prymitywy”*. Arkusz łamie ją **pięć razy** — i każde złamanie jest uzasadnione tą samą przyczyną: **tło jest stałe w obu motywach, więc tekst na nim też musi być stały.**

| Linia | Reguła | Prymityw | Powierzchnia pod spodem |
|---|---|---|---|
| 90 | `.dn-btn--sygnal { color: … }` | `--dn-szary-0` | `--dn-sygnal-wypelnienie` — wypełnienie sygnałowe, ciemne w obu motywach |
| 307 | `.dn-przelacznik:checked::before { background: … }` | `--dn-szary-0` | tor w barwie sygnału po włączeniu |
| 1073 | `.dn-awatar { color: … }` | `--dn-szary-100` | `--dn-grad-atrament` — gradient atramentowy, ciemny w obu motywach |
| 1083 | `.dn-awatar--inteligencja { color: … }` | `--dn-szary-0` | `--dn-grad-sygnal` |
| 1196 | `.dn-aod-rdzen { color: … }` | `--dn-szary-0` | `--dn-grad-sygnal` |

Gdyby użyć `--dn-tekst-inv`, w motywie ciemnym dałby `szary-950` — czyli **ciemny tekst na ciemnoniebieskim wypełnieniu**. Test dopuszczalności prymitywu: *czy powierzchnia, na której leży ten tekst, zmienia się z motywem?* Jeśli **nie** — prymityw jest poprawny i wymaga komentarza w kodzie.

### 11.4. Jedyny wyjątek `[data-theme]` w fundamencie

```css
::selection { background: var(--dn-sygnal-tlo); color: var(--dn-tekst); }
:root[data-theme='dark'] ::selection { background: var(--dn-sygnal-obrys); }
```

Przyczyna wyjątku jest strukturalna, nie estetyczna: **nie istnieje żeton roli „zaznaczenie”.** W motywie jasnym poprawny jest kryjący `--dn-sygnal-tlo` (`sygnal-100`), w ciemnym ta sama rola wymaga wartości półprzezroczystej — a taką w tym miejscu palety niesie `--dn-sygnal-obrys` (`rgba(92,140,236,.40)`). Zaznaczenie musi prześwitywać, bo pada na tekst o nieznanej z góry barwie.

**Rozstrzygnięcie (D-A5-07):** wyjątek zostaje do czasu wprowadzenia żetonu roli `--dn-zaznaczenie-tlo` do `zetony.css`. Jest to jedyne miejsce w całym pakiecie CSS, gdzie motyw rozstrzyga się selektorem, i **nie wolno się na nie powoływać** jako na precedens dla komponentów.

### 11.5. Trzy mechanizmy wyboru motywu

```
1. data-theme="light|dark" na <html>   ── wybór jawny Operatora (wspolne.js + localStorage)
2. :root:not([data-theme]) + @media (prefers-color-scheme)  ── zapas systemowy
3. color-scheme: light|dark            ── informacja dla przeglądarki: paski
                                          przewijania, pola natywne, menu select
```

Bloki `@media (prefers-color-scheme)` powielają wartości motywów **dosłownie**. Komentarz w `zetony.css` (linie 324–326) mówi: *„powielenie jest świadome (mechanizm kaskady, nie drugie źródło prawdy)”*. Konsekwencja dla utrzymania: **zmiana wartości żetonu semantycznego wymaga edycji w dwóch miejscach** — w bloku `[data-theme]` i w bloku `@media`. Sprawdza to punkt 14.11.

### 11.6. Lista kontrolna „styl działa w obu motywach”

- [ ] reguła zapisana raz, bez selektora `[data-theme]`
- [ ] wszystkie barwy z żetonów semantycznych (nie z prymitywów, nie z wartości szesnastkowych)
- [ ] tło nakładane na wskaźnik to `--dn-hover` / `--dn-wcisniecie` (półprzezroczyste), nie druga barwa kryjąca
- [ ] cień z `--dn-cien-*` (dwa komplety tonowane per motyw), nie `rgba` wpisane wprost
- [ ] element na pasku górnym używa `--dn-rama-*`, nie żetonów przełączalnych
- [ ] jeśli użyto prymitywu — powierzchnia pod spodem jest stała w obu motywach i jest o tym komentarz
- [ ] sprawdzone wzrokowo w obu motywach, nie tylko w domyślnym

---

## 12. Preferencje środowiska: ruch, dotyk, wymuszone barwy

### 12.1. `prefers-reduced-motion` — obsługa globalna

Dwa poziomy, oba poza komponentem:

**Poziom 1 — żetony** (`zetony.css`, sekcja 14, linie 435–448):

```css
@media (prefers-reduced-motion: reduce) {
  :root { --dn-czas-1: 0.01ms; --dn-czas-2: 0.01ms; --dn-czas-3: 0.01ms; --dn-czas-tetno: 0.01ms; }
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

Skrócenie żetonów gasi wszystko, co czas czerpie z systemu. Blok `*` z `!important` domyka szczelinę: literały (jak `0.8s` w `dn-obrot`) i style lokalne okna. `animation-iteration-count: 1` zatrzymuje pętle po jednym przebiegu — spinner nie kręci się w nieskończoność.

**Poziom 2 — zastąpienie informacji** (`komponenty.css`, linie 353–355): jedyna reguła komponentowa reagująca na preferencję. Kropka tętna dostaje **statyczny pierścień 2 px**, żeby nie zgubić znaczenia „tu biegnie praca” (rozdz. 8.7).

**Zasada:** komponent obsługuje `prefers-reduced-motion` **tylko wtedy**, gdy jego ruch niesie informację, której trzeba dostarczyć inaczej. Wygaszanie ruchu nie jest sprawą komponentu — robi to warstwa żetonów.

### 12.2. `pointer: coarse` — cele dotykowe rosną żetonem

`zetony.css`, sekcja 12 (linie 406–415) — sześć żetonów wymiarowych:

| Żeton | Wskaźnik precyzyjny | Dotyk | Co rośnie |
|---|---|---|---|
| `--dn-wym-kontrolka` | 32 px | **40 px** | przyciski, pola, pozycje wyboru, zakładki, pozycje nawigacji |
| `--dn-wym-ikonowy` | 32 px | **40 px** | przyciski ikonowe (pasek górny, zamknięcie karty sesji) |
| `--dn-wym-wiersz` | 36 px | **44 px** | wiersze tabel, kroki kolejki |
| `--dn-wym-check` | 16 px | **20 px** | pola wyboru i opcje jednokrotne |
| `--dn-wym-przelacznik-szer` | 36 px | **44 px** | przełącznik — tor |
| `--dn-wym-przelacznik-wys` | 20 px | **24 px** | przełącznik — wysokość (suwak liczony `calc(wys − 6px)`) |

**Zero reguł komponentowych z `pointer: coarse`** — sprawdzone. Kierunek projektowy formułuje to jako *„żeton, nie wyjątek”*.

Dwie konsekwencje dla autora komponentu:

1. **Nie wpisuj wysokości kontrolki wprost.** `min-height: 32px` nie urośnie na tablecie; `min-height: var(--dn-wym-kontrolka)` urośnie.
2. **Wymiary pochodne licz przez `calc()` na żetonie.** Arkusz robi to konsekwentnie: `.dn-btn--sm` to `calc(var(--dn-wym-kontrolka) - var(--dn-od-1))`, suwak przełącznika to `calc(var(--dn-wym-przelacznik-wys) - 6px)`. Wartość pochodna skaluje się razem z bazą.

### 12.3. `forced-colors` — luka i rozstrzygnięcie

**Stan faktyczny: `zetony.css`, `fundament.css` i `komponenty.css` nie zawierają ani jednej reguły `@media (forced-colors: active)`.** Tryb wysokiego kontrastu systemu operacyjnego nie jest w pakiecie obsłużony.

Zachowanie bez obsługi (przeglądarki Chromium i Firefox na Windows z włączonym trybem kontrastu):

| Właściwość | Zachowanie | Skutek w kokpicie |
|---|---|---|
| `color`, `background-color`, `border-color` | zastąpione barwami systemu | monochrom + sygnał znikają; układ i czytelność pozostają |
| `box-shadow` | **usunięty** | znika poświata fokusu pól (`--dn-cien-sygnal`) i pierścień tętna kropki |
| `outline` | zachowany, przeliczony na barwę systemu | **pierścień fokusu działa** — reguła globalna z fundamentu przechodzi |
| `background-image` (gradienty) | usunięty | awatary i rdzeń Always On Display tracą wypełnienie |
| obrys 1 px | zachowany | struktura kart, paneli i tabel pozostaje czytelna |

**Ocena ryzyka:** poważna jest jedna pozycja — utrata `box-shadow` gasi **pierścień tętna kropki sygnału** i **poświatę fokusu pola**. W obu przypadkach system ma zapasowe wskazanie (kropka pozostaje wypełniona i towarzyszy jej etykieta; pole zmienia barwę obrysu na `--dn-fokus`, a obrys w trybie wymuszonym jest zachowywany), więc żadna informacja nie znika całkowicie.

**Rozstrzygnięcie (D-A5-08)** — reguła spójna z systemem, do dopisania w `fundament.css` przy najbliższym otwarciu pliku:

```css
/* --- Wymuszone barwy systemu — zachowanie wskazań niesionych cieniem ------ */
@media (forced-colors: active) {
  /* fokus pola: cień zastępuje obrys elementu w barwie podświetlenia systemu */
  .dn-pole-kontrolka:focus,
  .dn-pasek-szukaj:focus,
  .dn-prompt:focus-within { outline: var(--dn-wym-fokus) solid Highlight; outline-offset: 0; }
  /* stan wybrany: kolor tła znika — wskazanie przejmuje obrys */
  .dn-tabela tbody tr[aria-selected='true'],
  .dn-boczna-pozycja[aria-current='page'] { outline: 1px solid Highlight; }
  /* kropka sygnału: pierścień tętna znika — kropka dostaje obrys systemowy */
  .dn-kropka { outline: 1px solid CanvasText; }
}
```

Zapis używa **słów kluczowych systemowych** (`Highlight`, `CanvasText`) — jedynych wartości sensownych w tym trybie, bo żetony i tak są nadpisywane przez system. To jedyne miejsce w pakiecie, w którym wolno napisać barwę spoza `var(--dn-*)`, i wyłącznie w postaci słowa kluczowego, nigdy wartości szesnastkowej. Zakres zamierzony: **niecałe 15 linii w `fundament.css`**, zero zmian w komponentach.

### 12.4. Zapytania, których w tych arkuszach nie ma

| Zapytanie | Obecność | Uwaga |
|---|---|---|
| `@media (min-width: …)` | **0 wystąpień** | żaden punkt łamania nie znajduje się w `fundament.css` ani `komponenty.css` — komponent nie zna szerokości okna. Układ responsywny należy do arkusza lokalnego okna |
| `@container` | 0 | zapytania kontenerowe nie są w pakiecie używane |
| `@supports` | 0 | wszystkie użyte właściwości (`dvh`, `clip-path`, `:focus-visible`, `::backdrop`, `backdrop-filter`, `translate`) mają wsparcie w silnikach docelowych |
| `prefers-contrast` | 0 | patrz 12.3 — obsługiwane pośrednio przez `forced-colors` |

**Pułapka do zapamiętania:** żetony `--dn-bp-w1…w4` (640 / 960 / 1280 / 1600 px) **nie działają wewnątrz `@media`** — własności niestandardowe nie są dozwolone w warunku zapytania medialnego. W arkuszu lokalnym okna punkt łamania trzeba wpisać literałem, a literał **musi** być równy wartości żetonu, z komentarzem wskazującym żeton:

```css
/* Punkt łamania w2 = var(--dn-bp-w2) = 960 px — boczna nawigacja zwija się do ikon */
@media (min-width: 960px) { … }
```

---

## 13. Antywzorce CSS w tym systemie

Lista zakazów. Kolumna „stan” podaje wynik pomiaru na obu arkuszach źródłowych.

| # | Antywzorzec | Stan | Uzasadnienie zakazu | Zamiast tego |
|---|---|---|---|---|
| 1 | **Wartość szesnastkowa wprost** | 1 wystąpienie (linia 163) | Wartość poza żetonem nie zmieni się przy przełączeniu motywu ani przy korekcie palety; omija pomiar kontrastu z `kontrasty.json` | `var(--dn-*)` |
| 2 | **`#000000`** | 0 | Czysta czerń zabija głębię i wywołuje halację na OLED; skala kończy się na `#0A0A0A` (kierunek systemu projektowego) | `--dn-szary-950` (tło) / `--dn-szary-900` (tekst) |
| 3 | **Prymityw w komponencie** | 5 uzasadnionych + 1 nieuzasadnione | Prymityw nie przełącza się z motywem — komponent traci zdolność reakcji | żeton semantyczny; prymityw wyłącznie na powierzchni stałej (rozdz. 11.3) |
| 4 | **`disabled` / `aria-disabled` jako brama** | 0 | Zasada zero blokad: domyślne zachowanie systemu to wykonanie polecenia; wyszarzony przycisk nie mówi, czego brakuje | przycisk klikalny + komunikat po naciśnięciu albo opis obok |
| 5 | **`opacity: .5` jako „wyłączenie”** | 0 | Obejście zakazu z punktu 4; dodatkowo psuje kontrast wszystkich warstw pod spodem | `--dn-tekst-2` / `--dn-tekst-3` jeśli chodzi o hierarchię; plakietka jeśli chodzi o stan |
| 6 | **`!important`** | 0 w komponentach; 4 w `zetony.css` (blok ograniczonego ruchu) | Kończy kaskadę; jedyny sposób nadpisania to kolejny `!important` | popraw specyficzność albo kolejność w pliku; wzorzec podwojenia klasy (rozdz. 4.3) |
| 7 | **Selektor identyfikatorowy** | 0 | Specyficzność (1,0,0) — nie do pokonania klasą; wiąże styl z jednym wystąpieniem | klasa `.dn-*` |
| 8 | **Zagnieżdżenie przez kontekst rodzica** | 0 | Komponent zmienia się od miejsca, w którym stoi; podział na pliki staje się nierozstrzygalny (rozdz. 4.4) | modyfikator komponentu albo reguła układu w arkuszu okna |
| 9 | **Goły selektor elementowy w komponencie** | 0 | Dosięga elementów, których autor reguły nie widział | kwalifikator klasy (`textarea.dn-pole-kontrolka`) albo dziecko bezpośrednie (`> svg`) |
| 10 | **`transition: all`** | 0 | Animuje właściwości układu w każdej klatce; wymusza jeden czas dla wszystkiego | wyliczona lista właściwości (rozdz. 8.3) |
| 11 | **Animacja ciągła poza tętnem** | 1 (spinner — uzasadniona) | `INTENSYWNOSC_RUCHU` 3/10: jeden ruch znaczący na widok | tętno kropki; spinner wyłącznie na czas trwania operacji |
| 12 | **Gradient jako tło przycisku, karty, sekcji** | 0 | kontrakt systemu projektowego: gradient wyłącznie ilustracyjny — awatary, rdzeń AOD, grafiki brandowe | powierzchnia kryjąca `--dn-powierzchnia` / `--dn-panel` |
| 13 | **`backdrop-filter` poza nakładką modala** | 1 (`.dn-modal::backdrop`) | Anty-domyślne: „glassmorfizm wszędzie” — powierzchnie mają być kryjące | tło kryjące + cień `--dn-cien-*` |
| 14 | **`z-index` liczbowy wprost** | 0 | Wojna warstw; jedenaście poziomów jest zdefiniowanych (`--dn-z-*`) | `var(--dn-z-modal)`, `var(--dn-z-tooltip)`, `var(--dn-z-aod)`… |
| 15 | **Wartość pikselowa spoza skali 4 px** | 27 literałów (rozdz. 13.1) | Jednostka 4 px jest [NIENEGOCJOWALNA]; wartość spoza skali psuje rytm i nie reaguje na gęstość | `--dn-od-*`, `--dn-wym-*`, `calc()` na żetonie |
| 16 | **`[data-theme]` w komponencie** | 0 | Podwaja regułę i podnosi specyficzność jednego motywu ponad drugi (rozdz. 11.1) | żeton semantyczny |
| 17 | **Stan wyrażony samą barwą** | 0 | kontrakt systemu projektowego: stan zawsze niesie ikonę albo etykietę | barwa + kształt (wstęga, kreska, podkreślenie) + ikona albo etykieta |
| 18 | **`outline: none` bez zastąpienia** | 0 (3 zastąpienia poprawne) | Usuwa jedyne wskazanie fokusu dla nawigacji klawiaturą | obrys w barwie `--dn-fokus` + poświata `--dn-cien-sygnal` (rozdz. 7.4) |
| 19 | **Emoji jako ikona** | 0 | Anty-domyślne kierunek systemu projektowego; emoji renderuje się inaczej w każdym systemie i nie przyjmuje `currentColor` | SVG z `zasoby/ikony/svg/` (82 pliki), obrys 1,75 |
| 20 | **Jednostka `em`/`rem` w odstępach** | 0 | Odstęp zależałby od stopnia pisma rodzica — rytm 4 px przestaje być przewidywalny | `--dn-od-*` (px); `ch` dopuszczone dla szerokości bloku tekstu (`44ch`, `40ch`, `36ch`) |

### 13.1. Odstępstwa faktycznie obecne — pełny wykaz do przeglądu

Nagłówek `komponenty.css` deklaruje „zero wartości zaszytych”. Pomiar wskazuje **27 literałów pikselowych** i **jeden literał szesnastkowy**. Zestawienie z oceną:

| Kategoria | Wystąpienia | Ocena | Uzasadnienie |
|---|---|---|---|
| **Grubość obrysu 1 px / 2 px** | ~20 | **dopuszczalne** | Krawędź włosowa nie należy do skali odstępów; skalowanie jej z gęstością pogorszyłoby rysunek |
| **Geometria wewnętrzna znaku** | 10 px i 8 px (ptaszek, kropka radio), 12 px (ikona plakietki i kroku), 16 px (kciuk suwaka), 20 px (znak kroku), 28 px (ikona pustego stanu), 40 px (godło karty środowiska) | **dopuszczalne warunkowo** | Wymiary rysunku wewnątrz kontrolki; **brak żetonów** dla tych wielkości. Zalecenie: przy najbliższej edycji podpiąć pod `--dn-wym-ikona-sm` (14), `--dn-wym-ikona` (16), `--dn-wym-ikona-lg` (20), `--dn-wym-ikona-xl` (24) tam, gdzie wartości się pokrywają |
| **Wysokość toru** | 4 px (`.dn-postep-tor`, `.dn-suwak`) | **dopuszczalne** | Zgodne ze skalą (`--dn-od-1`); zapis literałem dla czytelności rysunku |
| **Ograniczenia szerokości** | 260 px (dymek), 420 px (wyszukiwanie na pasku), 860 px (wpis), 720 px (wysokość modala), 160 px (maks. wysokość promptu) | **do przeglądu** | Wielkości kompozycyjne bez żetonu; wpływają na czytelność wiersza tekstu. Zalecenie: rozważyć żetony `--dn-wym-*` |
| **Wysokość pola na pasku** | 30 px (`.dn-pasek-szukaj`) | **do poprawy** | Jedyna kontrolka niższa od `--dn-wym-kontrolka`, wpisana literałem; **nie urośnie** przy `pointer: coarse`. Zalecany zapis: `calc(var(--dn-wym-kontrolka) - var(--dn-od-1))` |
| **Czas animacji** | `0.8s` × 2 (`dn-obrot`) | **do przeglądu** | Brak żetonu `--dn-czas-obrot`; skutek praktyczny zerowy (rozdz. 8.6) |
| **Wartość szesnastkowa** | `var(--dn-sygnal-300, #8FB2F5)` (linia 163) | **do poprawy — priorytet** | Podwójne naruszenie: wartość szesnastkowa **i** prymityw w komponencie. Zapas jest zbędny — `--dn-sygnal-300` jest zdefiniowany bezwarunkowo w `:root`, więc nigdy nie zadziała. Zalecany zapis: `var(--dn-sygnal-300)` z komentarzem, że jest to prymityw celowy dla stanu wciśniętego na ramie (rozdz. 7.5) |
| **`rgba` wprost** | 1 (`.dn-pasek-szukaj { background: rgba(255,255,255,.06) }`) | **do poprawy** | Wartość jest identyczna z `--dn-rama-hover`; zalecany zapis: `var(--dn-rama-hover)` |

---

## 14. Lista kontrolna przeglądu arkusza przed scaleniem

Kolejność jest celowa: punkty blokujące (1–9) przed poprawkami jakościowymi (10–16). Punkt oznaczony **[B]** blokuje scalenie.

### 14.A. Wartości i żetony

| # | Sprawdzenie | Jak sprawdzić |
|---|---|---|
| 1 **[B]** | Zero wartości szesnastkowych | wyszukaj `#` w pliku — dopuszczalne wyłącznie w komentarzu |
| 2 **[B]** | Zero `#000000` i `rgb(0,0,0)` | wyszukaj `000000`, `rgb(0` — `rgba(0,0,0,…)` dopuszczalne wyłącznie w żetonach cieni motywu ciemnego |
| 3 **[B]** | Barwy wyłącznie z żetonów semantycznych | wyszukaj `--dn-szary-`, `--dn-sygnal-[0-9]`, `--dn-zielen`, `--dn-bursztyn`, `--dn-czerwien`; każde trafienie wymaga komentarza wg testu z rozdz. 11.3 |
| 4 | Odstępy z `--dn-od-*`, wymiary z `--dn-wym-*` | wyszukaj `px` — dopuszczalne: grubość obrysu, geometria rysunku wewnątrz znaku, wartości w `calc()` na żetonie |
| 5 | `z-index` wyłącznie z `--dn-z-*` | wyszukaj `z-index` |

### 14.B. Selektory i specyficzność

| # | Sprawdzenie | Jak sprawdzić |
|---|---|---|
| 6 **[B]** | Zero selektorów `#id` | wyszukaj `#` na początku selektora |
| 7 **[B]** | Zero `!important` | wyszukaj `!important` — jedyne dopuszczone miejsce to blok `prefers-reduced-motion` w `zetony.css` |
| 8 | Specyficzność ≤ (0,3,0) | policz jednostki klasowe w najdłuższym selektorze |
| 9 | Zero zagnieżdżeń przez kontekst rodzica | każdy selektor złożony ma **oba** człony z tego samego komponentu (rozdz. 4.4) |
| 10 | Selektor elementowy wyłącznie zakotwiczony | wyszukaj selektory zaczynające się od litery — dopuszczalne wyłącznie jako kwalifikator klasy |
| 11 | Modyfikator nie działa sam | reguła `--wariant` deklaruje wyłącznie różnicę, nie układ ani wymiary |

### 14.C. Stany i dostępność

| # | Sprawdzenie | Jak sprawdzić |
|---|---|---|
| 12 **[B]** | Zero `disabled`, zero `aria-disabled` jako bramy | wyszukaj `disabled` |
| 13 **[B]** | Zero `outline: none` bez zastąpienia | każde wystąpienie ma obok obrys w barwie `--dn-fokus` + poświatę albo reakcję `:focus-within` rodzica |
| 14 **[B]** | Stan nigdy samą barwą | każdy stan barwny ma towarzysza: ikonę, etykietę, wstęgę, kreskę albo podkreślenie |
| 15 | Komplet stanów komponentu interaktywnego | spoczynek · najechanie · naciśnięcie · fokus · wybrany · ładowanie · błąd (nagłówek `komponenty.css`) |
| 16 | Wariant zmieniający tło ma własne `:hover` | inaczej odziedziczy półprzezroczyste `--dn-hover` (rozdz. 6.4) |
| 17 | Stan po ARIA, gdy ARIA go opisuje | `[aria-selected]`, `[aria-current]`, `[aria-pressed]`, `[aria-busy]`, `[aria-invalid]` przed modyfikatorem klasowym |
| 18 | Cel dotykowy przez żeton | wysokości z `--dn-wym-kontrolka` / `--dn-wym-ikonowy` / `--dn-wym-wiersz`, nie literałem |

### 14.D. Motyw i ruch

| # | Sprawdzenie | Jak sprawdzić |
|---|---|---|
| 19 **[B]** | Zero `[data-theme]` w regule komponentu | wyszukaj `data-theme` |
| 20 **[B]** | Wygląd sprawdzony w obu motywach | przełącz `data-theme` i obejrzyj każdy stan, nie tylko spoczynek |
| 21 | Element na pasku górnym używa `--dn-rama-*` | pasek jest atramentowy w obu motywach |
| 22 | Zero `transition: all` | wyszukaj `transition` |
| 23 | Czasy z `--dn-czas-*` | wyszukaj `s ` i `ms` w wartościach `transition` / `animation` |
| 24 | Ruch niosący informację ma zamiennik statyczny | reguła w bloku `prefers-reduced-motion` (rozdz. 8.7) |
| 25 | Zmiana żetonu naniesiona w obu blokach | `[data-theme='…']` **oraz** `@media (prefers-color-scheme: …)` (rozdz. 11.5) |

### 14.E. Struktura pliku

| # | Sprawdzenie | Jak sprawdzić |
|---|---|---|
| 26 | Plik ≤ 300 linii | `wc -l` |
| 27 | Nagłówek komentarzowy obecny | tytuł, zależność od żetonów, zakres klas |
| 28 | `@keyframes` w pliku pierwszego użytkownika | i arkusz spinający importuje go **przed** plikami korzystającymi (rozdz. 9.4) |
| 29 | Komentarze po polsku, z polskimi znakami | przejrzyj wzrokowo — `ą ć ę ł ń ó ś ź ż` |
| 30 | Nazwy własne bez tłumaczenia | `Studio Editor`, `Workflow Builder`, `Chat Window`, `Coordinator`, `Executor 1`, `Validator` |
| 31 | Zero emoji | wyszukaj wzrokowo w komentarzach i w `content:` |
| 32 | Zero treści zastępczych | brak Lorem ipsum, brak zmyślonych nazwisk, brak zmyślonych metryk — również w komentarzach |

### 14.F. Odbiór wizualny

| # | Sprawdzenie |
|---|---|
| 33 | Komponent obejrzany w galerii `06-okna/komponenty.html` (brama odbioru z opracowania o przekazaniu, ruchu i dostępności) |
| 34 | Fokus widoczny na każdej kontrolce przy nawigacji Tab |
| 35 | Wygląd sprawdzony przy `pointer: coarse` (emulacja dotyku) |
| 36 | Wygląd sprawdzony przy `data-gestosc="przestronna"` |
| 37 | Wygląd sprawdzony przy `prefers-reduced-motion: reduce` |
| 38 | Plik otwiera się i działa z `file://`, bez serwera i bez budowania |

---

## 15. Decyzje projektowe

Rozstrzygnięcia podjęte ponad literalny zapis dokumentacji źródłowej, wraz z podstawą i uzasadnieniem. Każde jest zgodne z kontraktem systemu projektowego i katalogiem komponentów.

| # | Rozstrzygnięcie | Podstawa w źródłach | Uzasadnienie |
|---|---|---|---|
| **1** | **Test przynależności reguły do fundamentu** (trzy warunki, rozdz. 1.3) | Nagłówek `fundament.css` wylicza zakres, ale nie podaje kryterium rozstrzygania przypadków spornych | Bez kryterium fundament rośnie przy każdej wątpliwości; test przenosi decyzję z gustu na sprawdzalny warunek |
| **2** | **Reguła rozstrzygająca „element czy blok”** (rozdz. 3.4): klasa deklarująca własne `display`, wymiary i tło jest blokiem, choćby nazwa wyglądała na element | kontrakt systemu projektowego wylicza klasy, ale nie podaje reguły ich czytania; `komponenty.css` stosuje ją konsekwentnie (`.dn-karta-sesji`, `.dn-btn-ikona`, `.dn-kafel`), nie nazywając | Konwencja ma wadę wrodzoną — `.dn-karta-tytul` i `.dn-karta-sesji` są lekturalnie nierozróżnialne. Reguła wywiedziona z faktycznego zapisu arkusza, nie wymyślona |
| **3** | **Drabina wyboru nośnika stanu** (rozdz. 6.2): pseudoklasa → ARIA → modyfikator klasowy → `[data-stan]` | Nagłówek `komponenty.css` wymienia siedem stanów; rozkład nośników zmierzony w pliku (ARIA 17 wystąpień, modyfikator 7) | Kolejność wynika z pomiaru, nie z preferencji. „ARIA przed klasą” daje zysk dostępnościowy: nie da się rozjechać wyglądu ze stanem dla czytnika ekranu |
| **4** | **`[data-stan]` pozostaje furtką, nie wzorcem** — `data-*` zarezerwowane dla trybów `<html>` i zaczepów skryptu | Zero wystąpień `[data-stan]` w obu arkuszach; `data-theme` i `data-gestosc` na `:root`; `data-przelacz-motyw` bez reguły CSS w `wspolne.js` | Zadanie wymieniało `[data-stan]` wśród wzorców stanu; stan faktyczny jest inny. Odnotowano rzeczywistość i podano warunki dopuszczalności zamiast opisywać nieistniejący wzorzec |
| **5** | **Zachowanie dwóch powtórzeń reguły fokusu** w `.dn-btn` i `.dn-btn-ikona` mimo nadmiarowości | opracowanie o przekazaniu, ruchu i dostępności: podział na pliki per komponent; rozdz. 9 tego opracowania | `przycisk.css` ma być czytelny i poprawny w oderwaniu od fundamentu. Zakaz rozszerzania powtórzenia na kolejne komponenty |
| **6** | **Kolejność importu w arkuszu spinającym wyznaczają zależności `@keyframes`**, nie kolejność alfabetyczna | Trzy definicje klatek w `komponenty.css`, sześć miejsc użycia; opracowanie o przekazaniu, ruchu i dostępności nie porusza tematu | Najczęstsza pułapka podziału pliku: po rozbiciu spinner przestaje się obracać. Zależność musi być zapisana jawnie, z komentarzem |
| **7** | **Wyjątek `[data-theme]` dla `::selection` zostaje** do czasu wprowadzenia żetonu roli `--dn-zaznaczenie-tlo`; nie stanowi precedensu dla komponentów | `fundament.css` linie 52–53; brak żetonu roli „zaznaczenie” w `zetony.css` | Przyczyna wyjątku jest strukturalna (brak żetonu), nie estetyczna. Nazwanie luki chroni przed powoływaniem się na wyjątek w komponentach |
| **8** | **Reguła `forced-colors` do dopisania w `fundament.css`** (rozdz. 12.3), ~15 linii, słowa kluczowe systemowe, zero zmian w komponentach | Zero wystąpień `forced-colors` w pakiecie; kontrakt systemu projektowego: WCAG 2.1 AA jako warunek wejściowy | Jedyna luka o realnym skutku: usunięcie `box-shadow` gasi pierścień tętna i poświatę fokusu pola. Umiejscowienie w fundamencie zgodne z zasadą „preferencje środowiska poza komponentem” |
| **9** | **Cztery warunki przyjęcia klasy użytkowej** (rozdz. 10.1) i katalog pięciu odrzuconych rodzin (rozdz. 10.2) | Jedenaście istniejących klas w `fundament.css`, zero w `komponenty.css`; kontrakt systemu projektowego | Dokumentacja źródłowa nie zajmuje stanowiska wobec klas użytkowych. Bez zapisanej granicy każde okno dopisze własne — i rytm 4 px przestanie obowiązywać |
| **10** | **Podział na 14 plików z wyliczoną objętością** i przypisaniem klas (tabela 9.3) | opracowanie o przekazaniu, ruchu i dostępności wymienia nazwy plików bez objętości i bez przypisania klas; granice sekcji zmierzone w `komponenty.css` | Wykonawca potrzebuje sprawdzalnego przypisania. Suma zweryfikowana: 1189 + 18 = 1207 linii ✔; wszystkie pliki poniżej progu 300 linii |
| **11** | **Wykaz odstępstw od „zera wartości zaszytych”** z oceną trzystopniową: dopuszczalne / do przeglądu / do poprawy (rozdz. 13.1) | Nagłówek `komponenty.css` deklaruje „zero wartości zaszytych”; pomiar wskazuje 27 literałów pikselowych i 1 szesnastkowy | Deklaracja i stan faktyczny się rozchodzą. Milczenie oznaczałoby, że recenzja zgłasza te same pozycje przy każdym przeglądzie. Trzy pozycje oznaczono do poprawy: `var(--dn-sygnal-300, #8FB2F5)`, `rgba(255,255,255,.06)` w polu paska, `min-height: 30px` w polu paska |
| **12** | **Punkt łamania w arkuszu lokalnym zapisuje się literałem z komentarzem wskazującym żeton** | `--dn-bp-w1…w4` istnieją w `zetony.css`; zero zapytań szerokościowych w `fundament.css` i `komponenty.css` | Ograniczenie techniczne CSS: własności niestandardowe nie działają w warunku `@media`. Bez zapisanej reguły powstałby kod pozornie poprawny, który nigdy nie zadziała |
| **13** | **Uzasadnienie kolejności importu dowodem na `border-radius`** (rozdz. 2.3) zamiast odwołania do zwyczaju | `fundament.css` linia 48 (`border-radius: var(--dn-r-xs)` w regule fokusu); `komponenty.css` linia 22 (`border-radius: var(--dn-r-sm)` w `.dn-btn`); obie o specyficzności (0,1,0) | Kolejność importu bywa traktowana jako konwencja. Dowód czyni ją sprawdzalną: zamiana miejscami daje widoczny defekt — przycisk zmienia kształt przy fokusie klawiaturowym |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*

## Uzasadnienia reguł okien

Powody, dla których reguły poszczególnych okien są takie, a nie inne, prowadzi
[11 — uzasadnienia okien](11-uzasadnienia-okien.md). Arkusze okien niosą sam kod:
limit gęstości komentarzy wynosi 250 znaków na 1000 wierszy.
