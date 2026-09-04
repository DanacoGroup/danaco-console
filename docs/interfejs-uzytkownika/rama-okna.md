# Danaco Console — Rama okna aplikacji

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
| **Tytuł** | Rama okna aplikacji — belka tytułowa, szyna nawigacji, pasek edycji, pasek stanu i specyfikacja menu rozwijanych |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant · deweloper · redaktor treści interfejsu |
| **Przeznaczenie** | Pełny materiał źródłowy dla czterech pasów ramy obecnych w każdym oknie platformy: co zawierają, jak się zachowują, jakie niosą czynności i co dokładnie kryje każde menu rozwijane |
| **Zakres** | Belka tytułowa · szyna nawigacji środowisk i modułów · pasek edycji · pasek stanu · okno „Dostosuj wstążkę” · menu aplikacji z podmenu i pozostałe menu rozwijane · etykiety kontrolek · przypisywanie skrótów klawiszowych · zachowanie na progach szerokości · okna nakładkowe jako okna |
| **Poza zakresem** | Pas kart sesji, szyna sesji, panel orkiestracji, okna operacyjne — [Elementy okien](elementy-okien.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Elementy okien](elementy-okien.md) · [Przepływ okien](przeplyw-okien.md) · [Katalog komponentów](katalog-komponentow.md) · [System wizualny](system-wizualny.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) |
| **Prototypy odniesienia** | `design/05-okna/` — komplet 38 prototypów, rama występuje w każdym bez wyjątku · `design/05-okna/WZORZEC-STANOWISKA.html` — wzorzec kanoniczny |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `window`, `config`, `settings`, `session`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css` · `design/zasoby/rama.css` · `design/zasoby/prototyp.css` · `design/zasoby/okna-modalne.css` · `design/zasoby/rama.js` · `design/zasoby/okna-modalne.js` |
| **Zasada nadrzędna** | Danaco Console nie narzuca twardych blokad — każda kontrolka ramy odpowiada czynnością albo komunikatem mówiącym, co się stanie i gdzie tę czynność wykonać; żadna pozycja menu nie traci klikalności |

---

## Spis treści

1. [Czym jest rama okna](#1-czym-jest-rama-okna)
2. [Belka tytułowa](#2-belka-tytułowa)
   - [2.1 Skład](#21-skład)
   - [2.2 Tytuł widoku — reguła budowy](#22-tytuł-widoku--reguła-budowy)
   - [2.3 Sterowanie oknem systemowym](#23-sterowanie-oknem-systemowym)
3. [Szyna nawigacji](#3-szyna-nawigacji)
   - [3.1 Skład](#31-skład)
   - [3.2 Dwa rodzaje pozycji](#32-dwa-rodzaje-pozycji)
   - [3.3 Rozwijanie środowiska](#33-rozwijanie-środowiska)
   - [3.4 Etykieta pozycji](#34-etykieta-pozycji)
   - [3.5 Strefa pracy — nowa sesja, historia, nowy projekt](#35-strefa-pracy--nowa-sesja-historia-nowy-projekt)
   - [3.6 Szybki wybór](#36-szybki-wybór)
   - [3.7 Menu wywołane z szyny](#37-menu-wywołane-z-szyny)
4. [Pasek edycji](#4-pasek-edycji)
   - [4.1 Forma pasa — belka nakładana](#41-forma-pasa--belka-nakładana)
   - [4.2 Układ i rodziny kontrolek](#42-układ-i-rodziny-kontrolek)
   - [4.3 Forma kontrolki](#43-forma-kontrolki)
5. [Pasek stanu](#5-pasek-stanu)
6. [Nazewnictwo dwóch szyn](#6-nazewnictwo-dwóch-szyn)
7. [Okno „Dostosuj wstążkę”](#7-okno-dostosuj-wstążkę)
   - [7.1 Zakładka „Składniki wstążki”](#71-zakładka-składniki-wstążki)
   - [7.2 Zakładka „Karty sesji”](#72-zakładka-karty-sesji)
   - [7.3 Zakładka „Wygląd i zachowanie”](#73-zakładka-wygląd-i-zachowanie)
   - [7.4 Zakładka „Personalizacja widoku”](#74-zakładka-personalizacja-widoku)
   - [7.5 Stopka okna](#75-stopka-okna)
8. [Etykiety kontrolek i ruch](#8-etykiety-kontrolek-i-ruch)
9. [Menu aplikacji](#9-menu-aplikacji)
10. [Menu „Znajdź w widoku”](#10-menu-znajdź-w-widoku)
11. [Menu „Zrzut ekranu”](#11-menu-zrzut-ekranu)
12. [Menu „Schowek”](#12-menu-schowek)
13. [Menu „Ustawienia i konfiguracja”](#13-menu-ustawienia-i-konfiguracja)
14. [Menu „Profil użytkownika”](#14-menu-profil-użytkownika)
15. [Menu strefy pracy i rozwinięcie „Więcej”](#15-menu-strefy-pracy-i-rozwinięcie-więcej)
   - [15.1 Menu „Nowa sesja”](#151-menu-nowa-sesja)
   - [15.2 Menu „Nowy projekt”](#152-menu-nowy-projekt)
   - [15.3 Rozwinięcie „Więcej”](#153-rozwinięcie-więcej)
16. [Przypisywanie skrótów klawiszowych](#16-przypisywanie-skrótów-klawiszowych)
17. [Zachowanie na progach szerokości](#17-zachowanie-na-progach-szerokości)
18. [Warstwy widoczności i dostępność](#18-warstwy-widoczności-i-dostępność)
19. [Pole wyszukiwania i filtr zakresu](#19-pole-wyszukiwania-i-filtr-zakresu)
20. [Personalizacja pasów](#20-personalizacja-pasów)
21. [Okna nakładkowe jako okna](#21-okna-nakładkowe-jako-okna)
22. [Wykazy zamykające i kryteria odbioru](#22-wykazy-zamykające-i-kryteria-odbioru)
   - [22.1 Wykaz komend kontraktu](#221-wykaz-komend-kontraktu)
   - [22.2 Wykaz żetonów](#222-wykaz-żetonów)
   - [22.3 Wykaz komponentów](#223-wykaz-komponentów)
   - [22.4 Wykaz skrótów klawiszowych](#224-wykaz-skrótów-klawiszowych)
   - [22.5 Wykaz stanów kontrolek](#225-wykaz-stanów-kontrolek)
   - [22.6 Wykaz punktów łamania](#226-wykaz-punktów-łamania)
   - [22.7 Kryteria odbioru](#227-kryteria-odbioru)
23. [Załącznik. Pełny wykaz komend kontraktu ramy okna](#załącznik-pełny-wykaz-komend-kontraktu-ramy-okna)
   - [Obszar `window` — 7 komend](#obszar-window--7-komend)
   - [Obszar `session` — 19 komend](#obszar-session--19-komend)
   - [Obszar `home` — 1 komenda](#obszar-home--1-komenda)
   - [Obszar `config` — 1 komenda](#obszar-config--1-komenda)

---

## 1. Czym jest rama okna

Rama okna to **cztery pasy obecne w każdym oknie platformy**, niezależnie od
środowiska, modułu i widoku. Rama nie należy do żadnego środowiska ani modułu —
należy do aplikacji.

| Pas | Wymiar | Żeton | Tło | Co niesie |
|---|---:|---|---|---|
| Szyna nawigacji | szer. 48 px | `--dn-wym-belka` (żeton dzielony — brak osobnego żetonu szerokości szyny w `zetony.css`) | atrament w obu motywach | menu aplikacji, cztery środowiska, moduły środowiska rozwiniętego, profil |
| Belka tytułowa | wys. 48 px | `--dn-wym-belka` | atrament w obu motywach | znak i nazwa produktu, tytuł widoku, sterowanie oknem systemowym |
| Pasek edycji | wys. 48 px | `--dn-wym-pasek` | `--dn-panel` | czynności okna i treści, wyszukiwanie, karty otwartych sesji, widok, ustawienia |
| Pasek stanu | wys. 28 px | `--dn-wym-stan` | `--dn-panel` | operator, położenie, sesje czynne, motyw, obciążenie zasobów |

**Szyna biegnie od samej góry okna.** Zaczyna się przy górnej krawędzi, nie pod
belką tytułową, a belka tytułowa zaczyna się dopiero na prawo od niej. Szerokość
szyny równa się wysokości belki — obie miary niesie jeden i ten sam żeton
`--dn-wym-belka` = 48 px (`zetony.css` nie wyodrębnia osobnego żetonu
szerokości szyny), więc
pierwsza pozycja szyny — menu aplikacji — zajmuje kwadrat w narożniku okna.
Kwadrat ten jest stałym punktem odniesienia: kursor prowadzony w narożnik trafia
w menu aplikacji bez celowania.

**Szyna i belka łączą się bez kreski.** Oba pasy niosą tę samą barwę atramentu
i nie rozdziela ich ani obrys, ani kreska: rama okna czyta się jako jedna bryła
otaczająca treść, nie jako zestaw pasów. Od treści oddziela ramę różnica
odcienia — atrament wobec powierzchni — a nie linia.

**Podział ról między pasami.**

| Pas | Odpowiada na pytanie |
|---|---|
| Belka tytułowa | „W jakim programie jestem i jak steruję jego oknem?” |
| Szyna nawigacji | „Gdzie chcę pracować?” — wybór środowiska i modułu |
| Pasek edycji | „Co zrobić z tym, co widzę?” — czynności na oknie i treści |
| Pasek stanu | „W jakim stanie jest praca i maszyna?” |

**Dlaczego nawigacja stoi pionowo, a nie w pasku edycji.** Środowisk są cztery,
a modułów w środowisku od sześciu do dziewięciu — trzynaście punktów wejścia.
Pas poziomy mieści je wyłącznie kosztem ścieżki, wyszukiwania i czynności.
Szyna pionowa daje im miejsce bez kompromisu, a rozwinięcie modułów pod
środowiskiem pokazuje przynależność wprost: moduł stoi **pod** swoim
środowiskiem, nie obok niego.

**Gdzie rama występuje.** Bez wyjątku: okno startowe, okno rejestracji
i logowania, przygotowanie środowiska pracy, strona główna, przedsionki czterech
środowisk, przestrzenie robocze, okna platformowe (konfiguracja, ustawienia,
Mobile, Always On Display).

```
┌──────┬───────────────────────────────────────────────────────────────────┐
│  ☰   │ BELKA TYTUŁOWA · 48 px                                            │
│ 48×48├───────────────────────────────────────────────────────────────────┤
│+ ⏱ 🗀+│ ◈ Danaco Console   Raport końcowy — WorkSpace › Studio   ─ □ ✕    │
│ ──── │  strefa 1 · praca                                                 │
│ ▣ 🗨  ├───────────────────────────────────────────────────────────────────┤
│ ▣ 📁  │ ╭─ PASEK EDYCJI · 48 px — belka nakładana ──────────────────────╮ │
│  · ▪ │ │ ⇤ │ ← → │ ⌂ ⟳ │ ⛶ 📋 │ [ szukaj ] │ karty sesji │ ⋯ │ Dostosuj │ │
│  · ▪ │ ╰───────────────────────────────────────────────────────────────╯ │
│  · ▪ │                                                                   │
│ ▣ ⌨  │   ▣ środowisko — pozycja w ramce                                  │
│ ▣ ⚡  │   ▪ moduł — pozycja bez ramki, wcięta                             │
│ ──── │                                                                   │
│  ◉   │   ── treść widoku ──                                              │
├──────┴───────────────────────────────────────────────────────────────────┤
│ PASEK STANU · 28 px                                                      │
│ ◉ Operator │ WorkSpace · Studio │ Sesje: 7      ◐ ciemny │ CPU │ RAM     │
└──────────────────────────────────────────────────────────────────────────┘
```

Szyna nawigacji i pasek stanu biegną przez całą wysokość i całą szerokość okna;
belka tytułowa oraz pasek edycji mieszczą się w kolumnie na prawo od szyny.

---

## 2. Belka tytułowa

**Klasa bazowa:** `.dn-belka`

### 2.1. Skład

| Element | Klasa | Położenie | Co niesie |
|---|---|---|---|
| Znak i nazwa produktu | `.dn-belka-marka`, `.dn-belka-nazwa` | lewa krawędź | sygnet marki 18 × 18 px oraz napis „Danaco **Console**” krojem nagłówkowym |
| Tytuł bieżącego widoku | `.dn-belka-tytul` | środek | nazwa pracy w toku i miejsce, w którym się toczy |
| Sterowanie oknem | `.dn-belka-okno`, `.dn-belka-btn` | prawa krawędź | minimalizacja · maksymalizacja · zamknięcie |

### 2.2. Tytuł widoku — reguła budowy

Tytuł składa się z dwóch członów rozdzielonych półpauzą:

```
<nazwa pracy w toku> — <miejsce, w którym się toczy>
```

| Człon | Skąd pochodzi | Przykład |
|---|---|---|
| Nazwa pracy | Nazwa aktywnej karty sesji albo wskazanej pozycji szyny sesji | „Raport końcowy — redakcja” |
| Miejsce | Ścieżka położenia, dwa ostatnie człony | „WorkSpace › Studio” |

Gdy praca w toku nie jest wskazana — na przykład w przedsionku środowiska albo
na stronie głównej — tytuł zawiera samo miejsce. Tytuł jest przycinany
wielokropkiem, nigdy nie łamie się na dwa wiersze.

### 2.3. Sterowanie oknem systemowym

| Kontrolka | Wymiar | Wskazanie kursorem | Zachowanie po naciśnięciu |
|---|---|---|---|
| Minimalizuj okno | 46 × 48 px | tło `--dn-rama-hover` | Okno schodzi do paska zadań systemu; stan pracy pozostaje nietknięty |
| Maksymalizuj okno | 46 × 48 px | tło `--dn-rama-hover` | Okno wypełnia ekran; ponowne naciśnięcie przywraca rozmiar |
| Zamknij okno | 46 × 48 px | tło `--dn-blad-tekst`, napis w kontrze | Okno zostaje zamknięte; sesje pracujące w tle trwają dalej po stronie serwera |

**Przeciąganie okna.** Obszar belki poza kontrolkami jest uchwytem
przeciągania okna natywnego (`-webkit-app-region: drag`). Obszar kontrolek jest
z przeciągania wyłączony — inaczej naciśnięcie zamknięcia rozpoczynałoby
przeciąganie.

---

## 3. Szyna nawigacji

**Klasa bazowa:** `.dn-szyna-nawigacji`

Pionowy pas przy lewej krawędzi okna, biegnący **od samej góry okna** aż do
paska stanu. Niesie atrament ramy — odwrócony względem powierzchni okna w obu
motywach — dzięki czemu czyta się jako część ramy, nie jako panel treści.

Szerokość szyny równa się wysokości belki tytułowej, więc **głowa szyny jest
kwadratem w narożniku okna**, a belka tytułowa zaczyna się dokładnie na prawo
od niego. Między szyną a belką nie ma żadnej kreski — oba pasy mają tę samą
barwę i łączą się w jedną bryłę ramy.

### 3.1. Skład

| Odcinek | Zawartość | Klasa |
|---|---|---|
| Głowa | Menu aplikacji (hamburger) — kwadrat narożnika | `.dn-szyna-poz` z `data-menu` |
| **Strefa 1 — praca** | Nowa sesja · Historia sesji · Nowy projekt | `.dn-szyna-sekcja` |
| **Strefa 2 — środowiska** | Cztery środowiska; pod środowiskiem rozwiniętym — jego moduły | `.dn-szyna-nawigacji-lista` |
| **Strefa 3 — szybki wybór** | Komponenty własne i okna platformowe wraz z pozycją „Dodaj skrót” | `.dn-szyna-skroty` |
| Stopka | Ustawienia i konfiguracja · motyw jasny i ciemny · profil użytkownika | `.dn-szyna-stopka` |

Poniższy szkic pokazuje mapę stref całej ramy — cztery pasy z rozdz. 1 wraz
z trzema strefami szyny nawigacji i dziesięcioma rodzinami paska edycji
z rozdz. 4.2.

```
┌────────────────────────────────────────────────────────────────────────┐
│ BELKA TYTUŁOWA — marka · tytuł widoku · sterowanie oknem systemowym    │
├──────┬─────────────────────────────────────────────────────────────────┤
│ GŁOWA│ PASEK EDYCJI — belka nakładana                                  │
│ menu │ ┌─────────────────────────────────────────────────────────────┐ │
│ apki │ │ 1·2·3·4 │ [ szukaj ] │ ── pasmo kart sesji ── │ 7·8·9 │ 10  │ │
├──────┤ └─────────────────────────────────────────────────────────────┘ │
│STREFA│                                                                 │
│  1   │  Nowa sesja · Historia sesji · Nowy projekt                     │
│praca │                                                                 │
├──────┤                                                                 │
│STREFA│                    T R E Ś Ć   W I D O K U                      │
│  2   │  ▣ TalkIn  ▣ WorkSpace  ▣ CodeStudio  ▣ MultitaskingAI          │
│środo-│  (moduły rozwijają się pod środowiskiem rozwiniętym)            │
│wiska │                                                                 │
├──────┤                                                                 │
│STREFA│  Assistant · Agents · Automations · Workspace · Mobile ·        │
│  3   │  Always On Display · + Dodaj skrót                              │
│skróty│                                                                 │
├──────┼─────────────────────────────────────────────────────────────────┤
│STOPKA│                                                                 │
│⚙ ◐ ◉ │                                                                 │
├──────┴─────────────────────────────────────────────────────────────────┤
│ PASEK STANU — ◉ operator · położenie · sesje czynne│◐ motyw · CPU · RAM│
└────────────────────────────────────────────────────────────────────────┘
```

Legenda: `1`–`4` rodziny lewej strony paska edycji (układ, nawigacja historii,
widok, treść i przechwytywanie); `7`–`9` rodziny prawej strony (kontekst
i praca w tle, widok i powiadomienia, rozwinięcie „Więcej”); `10` przycisk
„Dostosuj”; `⚙` ustawienia i konfiguracja; `◐` motyw jasny/ciemny; `◉` profil
użytkownika.

**Trzy strefy, trzy pytania.** Strefa pracy odpowiada na pytanie „co chcę
zacząć albo do czego wrócić”, strefa środowisk — „gdzie chcę pracować”, strefa
szybkiego wyboru — „po co sięgam niezależnie od miejsca”. Strefy rozdziela
**delikatna kreska** `--dn-rama-obrys-mocny` szerokości 24 px: na tyle słaba, by
nie rysować pasa w pasie, i na tyle obecna, by trzy zbiory nie zlewały się
w jedną kolumnę ikon.

**Przewijanie.** Głowa i stopka stoją nieruchomo; trzy strefy leżą we wspólnym
pasie `.dn-szyna-tresc`, który przewija się, gdy okno jest za niskie na komplet
pozycji. Etykieta pozycji ma położenie ustalane względem okna (`position: fixed`,
współrzędne nadaje `rama.js` w chwili wskazania kursorem), więc przewijany kontener jej
nie przycina.

Stopka szyny skupia to, co dotyczy **aplikacji i konta**, a nie bieżącego
widoku: konfigurację, motyw i profil. Kolejność jest stała — ustawienia,
motyw, profil — a profil zamyka szynę przy dolnej krawędzi.

### 3.2. Dwa rodzaje pozycji

Hierarchia jest niesiona **kształtem, nie barwą**: środowisko ma ramkę, moduł
jej nie ma.

| Cecha | Środowisko | Moduł |
|---|---|---|
| Klasa | `.dn-szyna-poz--srodowisko` | `.dn-szyna-poz--modul` |
| Wymiar | 40 × 40 px | 34 × 34 px |
| Ramka w spoczynku | **brak** | **brak** |
| Ramka przy wskazaniu kursorem | 1 px, `--dn-rama-obrys-mocny` | brak; samo tło `--dn-rama-hover` |
| Ikona | 18 × 18 px (`--dn-wym-ikona`), kreska 1,9 px | 16 × 16 px (`--dn-wym-ikona-sm`), kreska 1,9 px |
| Rola | punkt wejścia do środowiska; rozwija listę modułów | punkt wejścia do modułu |

**Ramka jest znakiem stanu, nie rangi.** W spoczynku żadna pozycja szyny nie ma
ramki — pas pozostaje spokojny. Ramkę nosi wyłącznie **środowisko rozwinięte**
(rozdz. 3.3); wskazanie środowiska kursorem zapowiada ją obrysem słabszym. Rangę
pozycji niesie wymiar i wcięcie, nie obramowanie.

**Czytelność kreski.** Ikony obu rodzajów pozycji stoją w pełnej jasności ramy
(`--dn-rama-tekst`), nie w barwie drugorzędnej: pozycje leżą na atramencie,
na którym kreska osłabiona gaśnie.

**Jeden byt — jeden rysunek.** Godła środowisk w szynie są dokładnie tymi
samymi rysunkami, które niosą karty środowisk na stronie głównej, a ikony
modułów — tymi, które niosą kafle modułów w przedsionku. Szyna nie ma prawa do
własnego wariantu znaku: rozpoznanie ikony jest wynikiem powtórzenia, a dwa
warianty tego samego bytu to dwa byty w oczach Operatora.

### 3.3. Rozwijanie środowiska

Naciśnięcie środowiska rozwija jego moduły **w dół**, przesuwając środowiska
następne niżej. **Rozwinięte jest najwyżej jedno środowisko naraz**: wybór
kolejnego zwija poprzednie. Ponowne naciśnięcie środowiska rozwiniętego zwija
je bez otwierania innego.

| Stan | Wskazanie |
|---|---|
| Środowisko zwinięte | bez ramki, tło przezroczyste, ikona `--dn-rama-tekst` |
| Środowisko rozwinięte | **ramka `--dn-rama-tekst`** oraz tło `--dn-rama-hover` |
| Moduł bieżący | wypełnienie sygnałowe `--dn-sygnal-wypelnienie`, ikona w kontrze |

**Dlaczego ramka, a nie wypełnienie.** Wypełnienie sygnałowe jest zarezerwowane
dla wskazania **modułu bieżącego** — miejsca, w którym Operator faktycznie
pracuje. Środowisko rozwinięte jest stanem listy, nie miejscem pracy; dostaje
więc znak innego rodzaju: obramowanie. Dwa różne znaki dla dwóch różnych
rzeczy, oba czytelne bez odwołania do samej barwy.

### 3.4. Etykieta pozycji

Każda pozycja szyny — środowisko i moduł — niesie etykietę ujawnianą
wskazaniem kursorem i fokusem, wysuwaną **w prawo od szyny**. Etykieta jest
dwuczłonowa:

| Człon | Zawartość | Krój |
|---|---|---|
| Nazwa | Nazwa środowiska albo modułu | krój bazowy, waga półgruba |
| Opis | Jedno zdanie o przeznaczeniu | krój bazowy, stopień najmniejszy, barwa trzeciorzędna |

Etykieta stoi na powierzchni `--dn-powierzchnia` z obrysem i cieniem —
czytelna w obu motywach. Szerokość do 280 px; opis łamie się na wiersze.
Etykieta **ustępuje, gdy menu pozycji jest otwarte** — menu samo się nazywa.

### 3.5. Strefa pracy — nowa sesja, historia, nowy projekt

**Klasa:** `.dn-szyna-sekcja`, pozycje `.dn-szyna-poz--praca`

Pierwsza strefa szyny, nad środowiskami. Trzy pozycje w stałej kolejności od
góry: **Nowa sesja**, **Historia sesji**, **Nowy projekt**. Stoją w szynie, a nie
na wstążce, bo nie są czynnościami na bieżącej treści — zaczynają pracę albo do
niej wracają, niezależnie od tego, co jest teraz otwarte.

**Te same drogi prowadzą z przedsionka.** Listwa przedsionka środowiska niesie
przycisk „Nowy projekt”, a szyna sesji przedsionka — „Pełna historia sesji”.
Obie kontrolki otwierają te same okna co strefa pracy szyny, z tą różnicą, że
środowisko jest już znane z kontekstu i menu wyboru nie pojawia się.

#### Nowa sesja — wymuszony wybór środowiska

Sesja zawsze należy do środowiska; poza środowiskiem nie istnieje. Kontrolka
wywołana spoza środowiska **nie zgaduje** — otwiera menu z czterema
środowiskami i czeka na wskazanie.

| Sytuacja | Zachowanie |
|---|---|
| Operator jest w przestrzeni roboczej środowiska | Sesja powstaje w tym środowisku, bez pytania; wchodzi jako kolejna karta |
| Operator jest w Centrum dowodzenia, w oknie platformowym albo w historii sesji | Menu żąda wskazania środowiska |
| Środowisko wskazane | Otwiera się **przedsionek** tego środowiska |

**Modułu menu nie pyta.** Moduł wiodący wybiera się kaflem w przedsionku — i to
dopiero ten wybór otwiera przestrzeń roboczą. Pytanie o moduł w tym miejscu
dublowałoby wybór, który przedsionek prowadzi lepiej: kaflem, z opisem i liczbą
sesji przy każdym module.

#### Historia sesji — sekcja ustawień otwierana jako nakładka

Historia sesji **nie jest osobnym oknem pełnoekranowym**. Jest sekcją okna
ustawień i konfiguracji, otwieraną jako **nakładka nad bieżącym widokiem**:
Operator zagląda do rejestru bez opuszczania miejsca pracy i wraca jednym
zamknięciem. Nakładka otwiera się z sekcją „Historia sesji” już wybraną, a lewa
kolumna pozostaje pełną nawigacją ustawień — z rejestru przechodzi się wprost do
konta, uwierzytelniania czy izolacji kontekstu.

Trzy drogi prowadzą do tej samej nakładki: pozycja **Historia sesji** w strefie
pracy szyny, przycisk **Pełna historia sesji** w szynie sesji przedsionka oraz
sekcja **Historia sesji** w oknie ustawień i konfiguracji.

Poniższy przebieg pokazuje, że niezależnie od drogi wejścia nakładka jest
tym samym oknem, zasilanym tymi samymi dwiema komendami kontraktu (rozdz.
22.1).

```
Operator          Punkt wejścia (A/B/C)              Warstwa serwerowa
   │                     │                                   │
   ├─ A: „Historia       │                                   │
   │  sesji” w strefie   │                                   │
   │  pracy szyny        │                                   │
   ├────────────────────►│                                   │
   ├─ B: „Pełna historia │                                   │
   │  sesji” w szynie    │                                   │
   │  sesji przedsionka  │                                   │
   ├────────────────────►│                                   │
   ├─ C: sekcja „Historia│                                   │
   │  sesji” w oknie     │                                   │
   │  ustawień           │                                   │
   ├────────────────────►│                                   │
   │                     │  session.list (sesje czynne,      │
   │                     │  zakończone)                      │
   │                     ├──────────────────────────────────►│
   │                     │  session.archive.list (archiwalne)│
   │                     ├──────────────────────────────────►│
   │                     │  sessions[] · total · presence[]  │
   │                     │◄──────────────────────────────────┤
   │  nakładka nad       │                                   │
   │  bieżącym widokiem, │                                   │
   │  sekcja „Historia   │                                   │
   │  sesji” wybrana     │                                   │
   │◄────────────────────┤                                   │
```

Legenda: `A` pozycja strefy pracy szyny nawigacji (rozdz. 3.5); `B` przycisk
przedsionka środowiska; `C` sekcja okna ustawień i konfiguracji (rozdz. 13).
Trzy drogi różnią się wyłącznie miejscem wywołania — okno docelowe, jego stan
początkowy i komendy zasilające są identyczne.

Rejestr obejmuje **wszystkie** sesje konta, ze wszystkich czterech środowisk,
w trzech zbiorach.

| Zbiór | Co obejmuje | Co robi wejście w pozycję |
|---|---|---|
| **Sesje czynne** | Sesje pracujące, także te trwające w tle po zamknięciu aplikacji — proces sesji żyje po stronie serwera | Powrót do sesji w stanie, w jakim została pozostawiona |
| **Sesje zakończone** | Sesje domknięte poleceniem Operatora albo wyczerpaniem zadania; zapis rozmowy i wytworów pozostaje | Otwarcie nowej sesji z odtworzonym kontekstem poprzedniej |
| **Sesje archiwalne** | Sesje przeniesione do archiwum konta — nie liczą się do limitu sesji czynnych ani do zajętości pamięci projektu | Wymaga wcześniejszego przywrócenia; dostępny też sam zapis do pobrania |

**Jedna oś kolumn.** Nagłówek kolumn i wszystkie wiersze trzech zbiorów
korzystają z jednej siatki (`--hs-siatka`), więc środowisko, moduł, czas
ostatniej pracy i zapis stoją w jednej osi przez całą wysokość rejestru —
niezależnie od tego, że każdy zbiór ma inne działania.

**Zależności niewidoczne w oknie.** Zakończenie sesji **nie usuwa jej wytworów**:
dokumenty, wersje w Session Repository i wpisy pamięci pozostają w projekcie, do
którego sesja należała — rejestr sesji jest rejestrem pracy, nie repozytorium.
Archiwizacja przenosi zapis sesji do archiwum konta wraz z pełnym kontekstem
rozmowy i listą plików; przywrócenie odtwarza jedno i drugie. Filtr środowiska,
zakres dat i sortowanie zawężają wszystkie trzy zbiory naraz, nigdy jeden
osobno.

#### Nowy projekt — wybór środowiska, potem Workspace

Projekt także należy do środowiska, więc i tu menu żąda wskazania. Po wskazaniu
**nie otwiera się przedsionek**: otwiera się od razu okno **Workspace** tego
środowiska w trybie zakładania projektu
(`platformowe/nowy-projekt.html`). Przedsionek nie jest po drodze, bo projekt nie
jest pracą w module — jest ramą, w której moduły później pracują.

Okno zakładania projektu zbiera pięć grup ustaleń:

| Grupa | Co ustala | Czego nie widać wprost |
|---|---|---|
| Tożsamość projektu | nazwa, katalog na dysku, jednostka pracy, termin, opis | Katalog jest **jedynym** miejscem zapisu wytworów; poza nim sesja nie zapisze niczego bez osobnej zgody Operatora. Nazwa pojawia się w szynie sesji, w historii sesji i w pasku stanu |
| Instrukcje i model | instrukcja stała, model wiodący, profil izolacji kontekstu | Instrukcja obowiązuje **każdą** sesję projektu i każdego agenta; sesja może ją uzupełnić, nie może jej znieść. Profil izolacji rozstrzyga, po co sesja może sięgnąć: katalog projektu, biblioteka środowiska, sieć |
| Zespół agentów | agenci przypisani do projektu | Agenci pracują na tej samej instrukcji stałej; zestaw zmienia się później w module Agents, także w trakcie trwania projektu |
| Automatyzacje | czynności wykonywane bez polecenia | Każdą można wstrzymać w module Automations; wstrzymanie nie usuwa jej z projektu |
| Materiały wejściowe | pliki wniesione na starcie | Trafiają do Project Library projektu i są widoczne dla wszystkich jego sesji, nie tylko dla pierwszej |

Okno kończy się trzema wyjściami: **założenie projektu wraz z otwarciem
pierwszej sesji**, **założenie bez sesji** (projekt czeka w Workspace
środowiska i pojawia się w historii sesji dopiero po pierwszym otwarciu) oraz
**zapis jako szablon projektu** do ponownego użycia.

### 3.6. Szybki wybór

**Klasa:** `.dn-szyna-skroty`

Pod czterema środowiskami stoi grupa szybkiego wyboru: komponenty własne oraz
okna platformowe, do których Operator sięga niezależnie od tego, w którym
środowisku pracuje.

| Pozycja | Co otwiera |
|---|---|
| Assistant | Profil asystenta — dialog głosowy i tekstowy nad bieżącą treścią |
| Agents | Agenci autonomiczni wykonujący zadania według zdefiniowanych ról |
| Automations | Automatyki uruchamiane zdarzeniem, harmonogramem albo poleceniem |
| Workspace | Projekt roboczy — zasoby, zadania i pamięć kontekstowa |
| Mobile | Widok mobilny platformy |
| Always On Display | Awatar pływający ponad powłoką |
| **Dodaj skrót** (`+`) | Dodanie kolejnej pozycji szybkiego wyboru |

Grupa odcina się od środowisk **odstępem, nie kreską** — szyna nie niesie
żadnych linii. Pozycje mają wymiar modułu (34 × 34 px) i nie mają ramki: nic nie
rozwijają, otwierają wskazany komponent.

**Pozycja „Dodaj skrót”** jest jedyną w szynie o obrysie przerywanym. Obrys
przerywany mówi o miejscu, które dopiero czeka na wypełnienie — nie o gotowej
pozycji. Szybki wybór jest zbiorem Operatora: stoi w nim to, co ulubione albo
najczęściej używane, a nie zestaw narzucony przez platformę.

### 3.7. Menu wywołane z szyny

Szyna stoi przy krawędzi okna, więc jej menu — aplikacji i profilu —
otwierają się **w prawo od szyny**, nie w dół.

| Menu | Kotwica | Uzasadnienie |
|---|---|---|
| Menu aplikacji (głowa szyny) | górą do pozycji | Menu jest długie; kotwiczenie górą daje mu całą wysokość okna |
| Ustawienia i konfiguracja (stopka szyny) | dołem do pozycji | Pozycja stoi przy dolnej krawędzi; kotwiczenie dołem trzyma menu w oknie |
| Profil użytkownika (stopka szyny) | dołem do pozycji | Jak wyżej |

Wysokość menu jest ograniczona wysokością okna pomniejszoną o belkę tytułową
i pasek stanu; przy przekroczeniu menu przewija się wewnątrz siebie.

---

## 4. Pasek edycji

**Klasa bazowa:** `.dn-narzedzia`

### 4.1. Forma pasa — belka nakładana

Pasek edycji nie jest pasem rozpiętym od krawędzi do krawędzi. Jest **belką
nakładaną**: leży na własnej tacy z odstępem 8 px ze wszystkich stron, ma
zaokrąglone narożniki `--dn-r-lg` i **nie ma obrysu w żadnym miejscu**. Od tacy
oddziela go wyłącznie różnica odcienia.

| Warstwa | Żeton | Motyw jasny | Motyw ciemny |
|---|---|---|---|
| Belka (wstążka) | `--dn-wstazka` | `#FFFFFF` | `#2A2A2A` |
| Taca pod belką | `--dn-wstazka-taca` | `#E3E3E3` | `#0F0F0F` |

Para żetonów istnieje osobno, ponieważ żetony ogólne na to zadanie nie
wystarczają: w motywie jasnym `--dn-panel` i `--dn-tlo` dzieli sześć stopni
skali — za mało, by belka odcięła się bez obrysu. Odległość obu odcieni jest
dobrana tak, by różnica czytała się jednakowo w obu motywach.

> **Reguła.** Tam, gdzie wystarcza różnica odcienia, kreska jest nadmiarem.
> Pasek edycji nie ma ani obrysu, ani kreski dolnej.

### 4.2. Układ i rodziny kontrolek

Układ pasa jest **przeglądarkowy**: sterowanie po lewej, pole wyszukiwania
zaraz za nim, a cała reszta szerokości należy do kart otwartych sesji. Rodziny
ikon prawej strony zamykają pas.

Kontrolki są pogrupowane w **rodziny** o wspólnym przeznaczeniu. Rodziny
rozdziela kreska pionowa wysokości 20 px w barwie `--dn-obrys-mocny`; kreska
nigdy nie występuje między pojedynczymi ikonami wewnątrz jednej rodziny.

> **Reguła.** Odstęp bez kreski czyta się jako brak ikony. Rodziny rozdziela
> zawsze kreska, nigdy sam odstęp.

| Kolejność | Rodzina | Kontrolki |
|---:|---|---|
| 1 | Układ | zwinięcie szyny sesji |
| 2 | Nawigacja historii | wstecz · do przodu |
| 3 | Widok | Centrum dowodzenia · odświeżenie widoku |
| 4 | Treść i przechwytywanie | zrzut ekranu · schowek |
| 5 | — | pole wyszukiwania, szerokość 340 px, minimum 180 px |
| 6 | — | **pasmo kart otwartych sesji** (`.dn-narzedzia-karty`), zajmuje resztę szerokości |
| 7 | Kontekst i praca w tle | kolejka zadań · magistrala kontekstu · izolacja kontekstu |
| 8 | Widok i powiadomienia | powiadomienia · tryb skupienia · pełny ekran |
| 9 | Rozwinięcie | **„Więcej”** (`⋯`) — znajdź w widoku · dokumentacja platformy · skróty klawiszowe |
| 10 | — | **przycisk „Dostosuj”** (`.dn-narzedzia-dostosuj`), zawsze ostatni, przyklejony do prawej krawędzi wstążki |

**Czego na wstążce nie ma.** Nowa sesja, historia sesji i nowy projekt stoją
w **strefie pracy szyny nawigacji** (rozdz. 3.5). Nie są czynnościami na
bieżącej treści, więc nie należą do pasa czynności — należą do pasa, który
prowadzi pracę.

**Rozwinięcie „Więcej”.** Trzy pozycje o najrzadszym użyciu — wyszukiwanie
w widoku wraz z jego ustawieniami, dokumentacja platformy i skróty klawiszowe —
stoją zwinięte pod jedną kontrolką `⋯`. Zwinięcie nie odbiera im skrótów
klawiszowych: `Ctrl+F`, `F1` i `Ctrl+/` działają niezależnie od tego, czy
kontrolka jest widoczna. Skład rozwinięcia zmienia okno „Dostosuj”.

Rodziny 1–4 stoją po lewej stronie pasa, dokładnie tak, jak w oknie
przeglądarki: powrót i przejście dalej, odświeżenie i dom, narzędzia treści,
pole wyszukiwania. Za pasmem kart sesji stoją rodziny 7–10 — czynności, które
nie dotyczą bieżącej treści, tylko pracy Operatora jako całości — a zamyka pas
przycisk „Dostosuj”.

**Podział prawej i lewej strony.** Po lewej stoi to, co działa **na tym, co
widać**: historia widoku, wyszukiwanie w widoku, przechwycenie widoku. Po
prawej — to, co dotyczy **pracy poza widokiem**: nowa sesja, kolejka zadań,
zasięg izolacji, pomoc. Pasmo kart sesji rozdziela oba porządki.

**Pasmo kart otwartych sesji.** Miejsce za polem wyszukiwania należy do kart
sesji otwartych w oknie — tak jak karty w przeglądarce. Pasmo jest elastyczne:
przy braku kart pozostaje puste, przy nadmiarze kart przycina je własną
krawędzią.

**Ścieżki położenia pasek edycji nie niesie.** Miejsce, w którym Operator
pracuje, nazywa tytuł widoku w belce tytułowej oraz — wewnątrz przestrzeni
roboczej — nagłówek okna roboczego. Powtórzenie tej samej informacji w pasie
edycji zabierało tylko miejsce kartom sesji.

**Czego pasek edycji nie zawiera.** Menu aplikacji, ustawienia i konfiguracja,
przełącznik motywu oraz profil użytkownika stoją w szynie nawigacji — pierwsze
w jej głowie, trzy pozostałe w stopce. Wszystkie cztery dotyczą aplikacji jako
całości, nie bieżącego widoku, więc ich miejscem jest pas należący do aplikacji,
nie pas czynności na treści. Wyboru środowiska i modułu pasek nie prowadzi:
należy on w całości do szyny. Wskaźnik sesji czynnych stoi w pasku stanu oraz
przy każdej sesji w szynie sesji; w pasku edycji byłby powtórzeniem.

**Pole wyszukiwania stoi zaraz za rodzinami lewej strony**, w miejscu, które
w przeglądarce zajmuje pole adresu. Jest polem wprowadzania, nie kontrolką
ikonową, i nie należy do żadnej rodziny. Skrót `Alt + D` ustawia w nim kursor
z dowolnego miejsca aplikacji — kombinacja odpowiada nawykowi przeniesionemu
z pola adresu przeglądarki, w którego miejscu pole stoi. Skrót `Ctrl/Cmd + K`
należy wyłącznie do palety poleceń i nie jest z polem wyszukiwania współdzielony.
Szerokość pola wynosi 340 px i zwęża się do minimum 180 px
(`.dn-narzedzia-szukaj`, `design/zasoby/rama.css`).

### 4.3. Forma kontrolki

| Cecha | Wartość |
|---|---|
| Klasa | `.dn-nrz-btn` |
| Wymiar | 34 × 34 px (`--dn-wym-ikonowy`), ikona 18 × 18 px (`--dn-wym-ikona`); 40 × 40 px na wskazaniu dotykowym (`@media (pointer: coarse)`) i w gęstości przestronnej |
| Obrys | **brak w każdym stanie** — kontrolka nie udaje przycisku obramowanego |
| Spoczynek | tło przezroczyste, ikona `--dn-tekst-2` |
| Wskazanie kursorem | tło `--dn-hover`, ikona `--dn-tekst`, uniesienie o 1 px |
| Naciśnięcie | powrót do położenia wyjściowego |
| Fokus | pierścień `--dn-fokus` 2 px z odsunięciem |
| Menu otwarte | tło `--dn-powierzchnia-2`, ikona `--dn-tekst` |
| Dwustanowa | `aria-pressed`; stan włączony jak wyżej |

**Dobór ikon.** Ikony rodziny „treść i przechwytywanie” nie mogą mieć
kształtu pełnej ramki — prostokąt w rozmiarze kontrolki czyta się jako obrys
przycisku, nie jako znaczenie. Zrzut ekranu niesie sylwetkę aparatu, schowek —
sylwetkę podkładki z klipsem, zwinięcie szyny — kreskę z grotem.

---

## 5. Pasek stanu

**Klasa bazowa:** `.dn-stan`

Pas o wysokości 28 px przy dolnej krawędzi okna, na całą jego szerokość.
Niesie stan pracy i stan maszyny — wyłącznie odczyt, bez czynności.

| Pozycja | Zawartość | Strona |
|---|---|---|
| Operator | Ikona sylwetki, słowo „Operator” i adres konta | lewa |
| Położenie | Środowisko i moduł bieżący | lewa |
| Sesje czynne | Liczba sesji pracujących w tle | lewa |
| Widok | Motyw bieżący — jasny albo ciemny — wraz z ikoną | prawa |
| CPU | Udział procesora wraz z torem wypełnienia | prawa |
| RAM | Zajętość pamięci wraz z torem wypełnienia | prawa |

**Miary zasobów.** Udział procesora i zajętość pamięci odczytuje **powłoka
natywna** aplikacji — warstwa systemowa ma dostęp do liczników procesu
i maszyny. W prototypie wartości są przykładowe: przeglądarka nie udostępnia
tych liczników. Tor wypełnienia przechodzi w barwę ostrzeżenia po
przekroczeniu progu wysokiego obciążenia.

**Bez ruchu.** Pasek stanu nie animuje wartości. Zmiana liczby albo szerokości
toru następuje skokowo — animowany licznik fałszowałby pomiar, pokazując
wartości, których system nigdy nie zmierzył.

**Zwężanie.** Poniżej 960 px pozycje drugorzędne — położenie, CPU i RAM —
ustępują; operator, sesje czynne i motyw pozostają.

---

## 6. Nazewnictwo dwóch szyn

Platforma ma dwa pionowe pasy przy lewej krawędzi i nie wolno ich mylić.
Określenie „szyna boczna” jest nazwą dwuznaczną i nie występuje w dokumentacji
ani w interfejsie.

| Nazwa wiążąca | Czym jest | Gdzie stoi | Kto ją zwija |
|---|---|---|---|
| **Szyna nawigacji** | Pas środowisk i modułów, część ramy okna | Skrajnie po lewej, od górnej krawędzi okna na całą jego wysokość | Menu aplikacji → Widok → Szyna nawigacji (`Ctrl+Shift+B`) |
| **Szyna sesji** | Kolumna sesji Operatora, część przestrzeni roboczej albo przedsionka | Na prawo od szyny nawigacji, wewnątrz treści widoku | Kontrolka „Zwiń szynę sesji” w pasku edycji oraz menu aplikacji → Widok → Szyna sesji (`Ctrl+B`) |

Kolumna listy modułów wewnątrz przestrzeni roboczej nosi nazwę **boczna
nawigacja modułów** i jest odrębna od obu szyn.

---

## 7. Okno „Dostosuj wstążkę”

**Wywołanie:** przycisk `.dn-narzedzia-dostosuj` zamykający pasek edycji od
prawej. **Klasa okna:** `.dn-modal.dn-modal--szeroki`, identyfikator
`dostosuj-wstazke`.

Wstążka nie jest układem narzuconym raz na zawsze. Operator sam rozstrzyga,
z czego się składa, w jakiej stoi kolejności i jak wyglądają na niej karty
sesji. Zasięg ustawienia — jeden skład dla wszystkich środowisk albo osobny
skład dla każdego z nich — rozstrzyga sam Operator w zakładce „Wygląd
i zachowanie”; niezależnie od wyboru ustawienie należy do konta i wędruje z nim
na każde urządzenie.

**Dlaczego przycisk niesie napis, a nie samą ikonę.** Pozostałe pozycje pasa są
czynnościami na treści; ta jedna otwiera konfigurację samego pasa. Napis
oddziela ją od rodzin ikon skuteczniej niż jakikolwiek dobór kształtu.

### 7.1. Zakładka „Składniki wstążki”

Trzy kolumny — **Poza paskami**, **Wstążka pozioma — kolejność od lewej** oraz
**Szyna pionowa — kolejność od góry** — a pod każdą własne przyciski
przenoszenia. Przycisk nazywa kierunek wprost, więc kolumny nie potrzebują
wspólnego środka.

| Kolumna | Kontrolki pod nią |
|---|---|
| Poza paskami | `Na wstążkę →` · `Na szynę →` |
| Wstążka pozioma | `← Odłóż` · `Na szynę →` · `↑` · `↓` |
| Szyna pionowa | `← Odłóż` · `Na wstążkę →` · `↑` · `↓` |

**Ten sam składnik nie stoi w dwóch miejscach naraz.** Przeniesienie na drugi
pas zdejmuje go z pierwszego — inaczej ta sama czynność miałaby dwa punkty
wywołania o dwóch różnych wyglądach, a Operator dwa miejsca do sprawdzenia.

**Skład stały.** Cztery karty środowisk pracy oraz menu aplikacji, ustawienia
i profil użytkownika należą do stałego składu szyny: nie podlegają przenoszeniu
ani usunięciu. Bez nich rama przestałaby prowadzić nawigację, a dostosowanie ma
zmieniać wygodę, nie odbierać drogi.

Pozycja wskazana niesie wypełnienie sygnałowe **oraz kreskę przy krawędzi** —
wskazanie nie opiera się na samej barwie. Każda z tych kontrolek odpowiada
także wtedy, gdy nie wskazano żadnej pozycji: komunikatem mówiącym, czego
brakuje, nie brakiem reakcji.

Rodzina pozycji jest wypisana przy nazwie, ponieważ to ona rozstrzyga o kresce
rozdzielającej: pozycje jednej rodziny stoją obok siebie, kreska pada między
rodzinami.

### 7.2. Zakładka „Karty sesji”

| Grupa | Rozstrzygnięcie |
|---|---|
| Zawartość karty | nazwa sesji i moduł · sama nazwa sesji · sama ikona modułu |
| Szerokość karty | stała 200 px · elastyczna |
| Nadmiar kart | przewijanie pasma · zwężanie kart · lista rozwijana |
| Oznaczenia karty | wskaźnik pracy w tle · przycisk zamknięcia · numer karty |

### 7.3. Zakładka „Wygląd i zachowanie”

| Grupa | Rozstrzygnięcie |
|---|---|
| Postać kontrolek | same ikony · ikony z etykietami · napisy tylko w wybranych rodzinach |
| Forma wstążki | wstążka nakładana · wstążka na całą szerokość |
| Gęstość | zwarta 48 px · swobodna 56 px |
| Zachowanie | ukrycie w trybie skupienia · pole wyszukiwania na wstążce · jeden skład we wszystkich środowiskach |

### 7.4. Zakładka „Personalizacja widoku”

Czwarta zakładka okna „Dostosuj wstążkę”, zweryfikowana w `design/05-okna/WZORZEC-STANOWISKA.html`
(`data-przelacz="dost-person"`) — pominięta w poprzedniej redakcji niniejszego dokumentu, obecna w
znaczniku wzorca kanonicznego jako czwarty przycisk zakładki (`role="tab"`), obok „Składniki”,
„Karty sesji” i „Wygląd i zachowanie” (rozdz. 7.1–7.3). Wprowadzenie zakładki (`.dn-dost-lid`),
cytat dosłowny:

> „Personalizacja dotyczy wyglądu obu pasów — wstążki poziomej i szyny pionowej. Barwy pobierane są
> z palety systemu projektowego, więc każdy wybór pozostaje czytelny w motywie jasnym i ciemnym.”

Sześć grup ustaleń (`fieldset.dn-dost-grupa`), w kolejności występowania w znaczniku:

| Grupa | Rodzaj kontrolki | Opcje (pierwsza — domyślnie wybrana) |
|---|---|---|
| Kolorystyka pasów | cztery pola radiowe z próbką barwy | Atrament (`--dn-rama`) · Grafit (`--dn-szary-800`) · Powierzchnia (`--dn-panel`) · Sygnał przygaszony (`--dn-sygnal-800`) |
| Kolor podświetlenia ikon | cztery pola radiowe z próbką barwy | Błękit sygnałowy (`--dn-sygnal-500`) · Zieleń (`--dn-zielen-400`) · Bursztyn (`--dn-bursztyn-400`) · Neutralny (`--dn-szary-400`) |
| Rodzaj separacji pozycji | trzy pola radiowe z opisem | „Kreska między rodzinami” (opis: „Pionowa kreska rozdziela rodziny; wewnątrz rodziny odstępu nie ma.”) · „Sam odstęp” („Rodziny rozdziela większy odstęp bez kreski.”) · „Tło rodziny” („Każda rodzina stoi na własnym, delikatnie wyodrębnionym tle.”) |
| Akcja po kliknięciu pozycji | trzy pola radiowe z opisem | „Wykonaj od razu” („Kliknięcie uruchamia czynność bez pytania.”) · „Pokaż podpowiedź, potem wykonaj” („Pierwsze kliknięcie odsłania opis, drugie wykonuje.”) · „Rozwiń menu wariantów” („Kliknięcie otwiera menu wariantów czynności zamiast wykonywać domyślny.”) |
| Odstępy i pozycjonowanie | dwa pola wyboru, oba domyślnie zaznaczone | „Pozwól wstawiać puste odstępy” („W oknie składników pojawia się pozycja „Pusty odstęp” do rozsuwania ikon wedle uznania.”) · „Wyrównaj rodziny do krawędzi” („Pierwsza rodzina przy lewej, ostatnia przy prawej krawędzi pasa.”) |
| Grupowanie w jedną ikonę | trzy pola wyboru, wszystkie domyślnie zaznaczone | „Zezwalaj na ikony zbiorcze” („Kilka pozycji można spiąć pod jedną ikoną z rozwijanym menu „⋯”.”) · „Zwijaj rzadko używane pozycje szyny” („Moduły używane rzadko chowają się pod ikoną zbiorczą szybkiego wyboru.”) · „Pozwól usuwać pozycje szybkiego wyboru” („Moduł zdjęty z szybkiego wyboru pozostaje dostępny w przedsionku środowiska.”) |

Nota zamykająca zakładkę (`.dn-dost-nota`), cytat dosłowny:

> „Pozycje zbiorcze buduje się w zakładce „Składniki”: zaznacz kilka pozycji jednego pasa i użyj
> „Spnij w ikonę zbiorczą”. Ikona zbiorcza otwiera po kliknięciu małe menu ze spiętymi pozycjami.”

Ta nota wiąże rozdz. 7.4 z rozdz. 7.1 wprost: mechanizm ikony zbiorczej (grupa „Grupowanie w jedną
ikonę”) jest ustawieniem włącz/wyłącz tej zakładki, lecz jego budowa — wybór, które konkretnie
pozycje spina jedna ikona — odbywa się w zakładce „Składniki”, nie tutaj. Cztery pierwsze grupy
(kolorystyka, podświetlenie, separacja, akcja) są zestawami wyboru jednokrotnego (pola radiowe,
klasa `.dn-radio`); dwie ostatnie (odstępy, grupowanie) są zestawami pól wyboru niezależnych
(`.dn-check`), każde przełączalne osobno.

### 7.5. Stopka okna

| Kontrolka | Czynność |
|---|---|
| `Przywróć układ domyślny` | Wraca skład wstążki i wszystkie rozstrzygnięcia do wartości fabrycznych |
| `Anuluj` | Zamyka okno bez zapisania zmian |
| `Zastosuj` | Zapisuje układ i zamyka okno; potwierdzenie niesie liczbę składników pozostałych na wstążce |

Stopka jest wspólna wszystkim czterem zakładkom (rozdz. 7.1–7.4) — „Zastosuj” zapisuje jednocześnie
rozstrzygnięcia poczynione w każdej z nich, nie tylko w zakładce aktualnie widocznej.

---

## 8. Etykiety kontrolek i ruch

**Etykieta.** Każda kontrolka pasa niesie etykietę zapisaną w atrybucie
`data-etykietka`, ujawnianą **wskazaniem kursorem i fokusem klawiatury**. Etykieta
pojawia się pod kontrolką, na powierzchni `--dn-powierzchnia` z obrysem
i cieniem — czytelnie w obu motywach. Etykieta ustępuje, gdy menu kontrolki
jest otwarte: otwarte menu przejmuje uwagę i samo nazywa czynności.

**Ruch.** Uniesienie o 1 px w czasie `--dn-czas-1` należy do wzorca
mikroreakcji kontrolki z katalogu dozwolonych animacji. Etykieta pojawia się
przejściem krycia i przesunięcia w czasie `--dn-czas-2`.

**Ograniczony ruch.** Przy `prefers-reduced-motion` uniesienie ustaje,
przejście etykiety ustaje; zmiana tła i sama etykieta pozostają. Informacja nie
ginie — zmienia się wyłącznie sposób jej podania.

---

## 9. Menu aplikacji

Kontrolka: hamburger, głowa szyny nawigacji. Najszersze menu platformy. **Menu
jest dwupoziomowe**: górna lista jest wąska i niesie same tytuły sekcji, a
szczegóły rozwijają się jako **podmenu w prawo** dopiero po wskazaniu kursorem albo
fokusie. Dzięki temu menu nie zasłania całego ekranu od góry do dołu — otwiera
się tyle, ile Operator wskaże.

Siedem sekcji górnych, każda z podmenu: **Plik**, **Edycja**, **Widok**,
**Przejdź**, **Narzędzia**, **Dla programistów**, **Pomoc**. Mechanika podmenu
stoi wyłącznie na `:hover` i `:focus-within` (klasy `.sta-podmenu`,
`.sta-podmenu-tresc`) — bez skryptu.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Plik | Nowa sesja | czynność | `Ctrl+N` |
| Plik | Nowe okno aplikacji | czynność | `Ctrl+Shift+N` |
| Plik | Otwórz projekt… | czynność | `Ctrl+O` |
| Plik | Zapisz sesję | czynność | `Ctrl+S` |
| Plik | Eksportuj sesję… | czynność | `Ctrl+Shift+E` |
| Plik | Importuj materiał… | czynność | — |
| Edycja | Cofnij | czynność | `Ctrl+Z` |
| Edycja | Ponów | czynność | `Ctrl+Y` |
| Edycja | Kopiuj | czynność | `Ctrl+C` |
| Edycja | Wklej | czynność | `Ctrl+V` |
| Edycja | Znajdź w widoku | czynność | `Ctrl+F` |
| Edycja | Znajdź i zamień | czynność | `Ctrl+H` |
| Widok | Szyna nawigacji | przełącznik | `Ctrl+Shift+B` |
| Widok | Szyna sesji | przełącznik | `Ctrl+B` |
| Widok | Pas kart sesji | przełącznik | — |
| Widok | Pasek stanu | przełącznik | — |
| Widok | Tryb skupienia | przełącznik | `Ctrl+Shift+F` |
| Widok | Pełny ekran | przełącznik | `F11` |
| Widok | Powiększ | czynność | `Ctrl+=` |
| Widok | Pomniejsz | czynność | `Ctrl+-` |
| Widok | Odśwież widok | czynność | `F5` |
| Przejdź | Centrum dowodzenia | czynność | `Ctrl+Home` |
| Przejdź | Przedsionek środowiska | czynność | `Ctrl+E` |
| Przejdź | Następna karta sesji | czynność | `Ctrl+Tab` |
| Przejdź | Poprzednia karta sesji | czynność | `Ctrl+Shift+Tab` |
| Przejdź | Paleta poleceń | czynność | `Ctrl+K` |
| Przejdź | Historia sesji | czynność (nakładka ustawień) | — |
| Narzędzia | Menedżer sesji | czynność | — |
| Narzędzia | Magistrala kontekstu | czynność | — |
| Narzędzia | Kolejka zadań | czynność | — |
| Narzędzia | Izolacja kontekstu | czynność | — |
| Narzędzia | Punkty izolacji | czynność | — |
| Narzędzia | Funkcje doświadczalne | czynność | — |
| **Dla programistów** | Terminal aplikacji | czynność | — |
| Dla programistów | Konsola i wykaz błędów | czynność | — |
| Dla programistów | Podgląd sieci | czynność | — |
| Dla programistów | Dziennik zdarzeń | czynność | — |
| Dla programistów | Ponowne wczytanie widoku | czynność | `F5` |
| Dla programistów | Restart aplikacji | czynność | — |
| Dla programistów | Naprawa aplikacji | czynność | — |
| Dla programistów | Sprawdź dostępność aktualizacji | czynność | — |
| Dla programistów | Otwórz katalog dzienników | czynność | — |
| Dla programistów | Zgłoś diagnostykę | czynność | — |
| Pomoc | Instrukcja użytkowania | otwiera okno | — |
| Pomoc | Instrukcja instalacji | otwiera okno | — |
| Pomoc | Dokumentacja platformy | czynność | `F1` |
| Pomoc | Skróty klawiszowe | czynność | `Ctrl+/` |
| Pomoc | Makiety i opracowania okien | czynność | — |
| Pomoc | Sprawdź dostępność aktualizacji | czynność | — |
| Pomoc | Zgłoś obserwację | czynność | — |
| Pomoc | Wersja i licencja | otwiera okno | — |
| Pomoc | O programie | czynność | — |

**Sekcja „Dla programistów”** zbiera narzędzia utrzymania i diagnostyki
aplikacji, które w produkcie natywnym odpowiadają czynnościom powłoki: terminal
aplikacji, konsola i wykaz błędów, podgląd sieci, dziennik zdarzeń, restart oraz
naprawa aplikacji (weryfikacja i odbudowa lokalnych danych) i sprawdzenie
dostępności aktualizacji.

**Sekcja „Pomoc”** prowadzi do trzech okien: **Instrukcji użytkowania**
(`platformowe/instrukcja-uzytkowania.html`), **Instrukcji instalacji**
(`platformowe/instalator.html`) oraz **Wersji i licencji**. Instrukcja
użytkowania jest oknem nakładkowym ze spisem rozdziałów i makietami układu ramy
oraz przestrzeni roboczej.

**Reguła doboru.** Do menu aplikacji trafia czynność spełniająca dwa warunki:
dotyczy aplikacji lub sesji jako całości oraz nie ma własnego miejsca w oknie
roboczym. Czynność należąca do jednego modułu należy do jego okna, nie tutaj.

---

## 10. Menu „Znajdź w widoku”

Kontrolka: lupa, rodzina „treść i przechwytywanie”. Wyszukiwanie w obrębie bieżącego widoku — odrębne od pola wyszukiwania po prawej stronie pasa, które przeszukuje sesję, projekt i środowisko. Menu otwiera pole wpisu wraz z ustawieniami zakresu i sposobu dopasowania.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Zakres wyszukiwania | Okno robocze wiodące | wybór jednokrotny | — |
| Zakres wyszukiwania | Wszystkie okna bieżącej karty | wybór jednokrotny | — |
| Zakres wyszukiwania | Historia rozmowy | wybór jednokrotny | — |
| Sposób dopasowania | Rozróżniaj wielkość liter | przełącznik | — |
| Sposób dopasowania | Całe słowa | przełącznik | — |
| Sposób dopasowania | Wyrażenie regularne | przełącznik | — |
| Sposób dopasowania | Znajdź i zamień | czynność | `Ctrl+H` |
| Sposób dopasowania | Ostatnie wyszukiwania | czynność | — |
| Sposób dopasowania | *Skrót:* Znajdź w widoku | pole skrótu | `Ctrl + F` |

**Rozgraniczenie.** Lupa po lewej szuka **w tym, co widać**. Pole po prawej szuka **w tym, co jest** — w sesjach, projektach i materiałach środowiska.

---

## 11. Menu „Zrzut ekranu”

Kontrolka: aparat, rodzina „treść i przechwytywanie”. Menu dzieli się na trzy sekcje: rodzaj zrzutu, miejsce docelowe oraz ustawienia wraz z przypisaniem skrótów.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Wykonaj zrzut | Całe okno aplikacji | czynność | `PrtSc` |
| Wykonaj zrzut | Zaznaczony obszar | czynność | `Ctrl+PrtSc` |
| Wykonaj zrzut | Wskazane okno robocze | czynność | `Alt+PrtSc` |
| Wykonaj zrzut | Widok przewijany w całości | czynność | — |
| Wykonaj zrzut | Zrzut z opóźnieniem 5 s | czynność | — |
| Po wykonaniu | Do schowka | wybór jednokrotny | — |
| Po wykonaniu | Do pliku w katalogu zrzutów | wybór jednokrotny | — |
| Po wykonaniu | Do modułu Library | wybór jednokrotny | — |
| Po wykonaniu | Do Chat Window jako załącznik | wybór jednokrotny | — |
| Ustawienia | Otwórz w edytorze adnotacji po wykonaniu | przełącznik | — |
| Ustawienia | Ukryj dane wrażliwe automatycznie | przełącznik | — |
| Ustawienia | Katalog docelowy i wzorzec nazwy… | czynność | — |
| Ustawienia | *Skrót:* Skrót zrzutu obszaru | pole skrótu | `Ctrl + PrtSc` |
| Ustawienia | *Skrót:* Skrót zrzutu okna | pole skrótu | `Alt + PrtSc` |

**Miejsce docelowe jest wyborem jednokrotnym** — zrzut trafia w jedno miejsce, nie w kilka naraz. Wybór obowiązuje do zmiany, także po ponownym uruchomieniu aplikacji.

**Ukrycie danych wrażliwych** zamazuje pola haseł, tokeny i wartości oznaczone w oknie jako poufne, zanim obraz opuści aplikację.

---

## 12. Menu „Schowek”

Kontrolka: podkładka z klipsem, rodzina „treść i przechwytywanie”. Menu daje dostęp do historii schowka, jego ustawień i skrótów.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Schowek | Otwórz schowek | czynność | `Ctrl+Shift+V` |
| Schowek | Wklej ostatnią pozycję | czynność | `Ctrl+V` |
| Schowek | Wklej jako zwykły tekst | czynność | `Ctrl+Shift+T` |
| Schowek | Historia schowka | czynność | — |
| Schowek | Wyczyść schowek | czynność | — |
| Ostatnie pozycje | „Sprawozdanie końcowe za trzeci kwartał” — fragment | czynność | — |
| Ostatnie pozycje | raport-koncowy.md — ścieżka pliku | czynność | — |
| Ostatnie pozycje | Zrzut obszaru — 12:04 | czynność | — |
| Ustawienia | Liczba przechowywanych pozycji: 25 | czynność | — |
| Ustawienia | Pomijaj treść z pól haseł | przełącznik | — |
| Ustawienia | Współdziel schowek między urządzeniami | przełącznik | — |
| Ustawienia | Czyść schowek przy zamknięciu aplikacji | przełącznik | — |
| Ustawienia | *Skrót:* Otwarcie schowka | pole skrótu | `Ctrl + Shift + V` |
| Ustawienia | *Skrót:* Wklejenie z listy | pole skrótu | `Ctrl + Alt + V` |

**Wklejanie z listy.** Przytrzymanie kombinacji wklejania otwiera listę przechowywanych pozycji; zwolnienie klawiszy wstawia pozycję wskazaną. Wzorzec jest zgodny z zachowaniem znanym z systemu operacyjnego.

**Pomijanie treści z pól haseł** jest domyślnie włączone: treść skopiowana z pola maskowanego nie trafia do historii.

---

## 13. Menu „Ustawienia i konfiguracja”

Kontrolka: koło zębate w **stopce szyny nawigacji**, nad przełącznikiem motywu. Menu rozdziela trzy zakresy konfiguracji platformy — aplikacja, środowisko, sesja — oraz prowadzi do punktów izolacji, uprawnień, rozszerzeń i utrzymania.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Aplikacja | Ustawienia aplikacji | czynność | `Ctrl+,` |
| Aplikacja | Język i format zapisu | czynność | — |
| Aplikacja | Motyw i gęstość widoku | czynność | — |
| Aplikacja | Skróty klawiszowe | czynność | — |
| Aplikacja | Powiadomienia | czynność | — |
| Aplikacja | Urządzenia i synchronizacja | czynność | — |
| Środowisko i sesja | Konfiguracja środowiska bieżącego | czynność | — |
| Środowisko i sesja | Konfiguracja karty sesji | czynność | — |
| Środowisko i sesja | Modele, kanały i poziom wysiłku | czynność | — |
| Środowisko i sesja | Magistrala kontekstu i pamięć | czynność | — |
| Środowisko i sesja | Punkty izolacji | czynność | — |
| Środowisko i sesja | Uprawnienia i zaufanie | czynność | — |
| Rozszerzenia | Wtyczki i konektory | czynność | — |
| Rozszerzenia | Serwery narzędziowe | czynność | — |
| Rozszerzenia | Funkcje doświadczalne | czynność | — |
| Utrzymanie | Kopie zapasowe i przywracanie | czynność | — |
| Utrzymanie | Pamięć podręczna i dane lokalne | czynność | — |
| Utrzymanie | Diagnostyka aplikacji | czynność | — |

**Trzy zakresy.** Ustawienie zmienione w zakresie aplikacji obowiązuje wszędzie; w zakresie środowiska — w jednym środowisku; w zakresie sesji — w jednej karcie. Zakres węższy nadpisuje szerszy.

---

## 14. Menu „Profil użytkownika”

Kontrolka: sylwetka zamykająca **stopkę szyny nawigacji**. Menu otwiera dane konta, urządzenia, status obecności i wyjście z aplikacji. Nagłówek menu nosi adres zalogowanego Operatora oraz plan i nazwę urządzenia.

| Sekcja | Pozycja | Rodzaj | Skrót |
|---|---|---|---|
| Zalogowany operator | Dane konta | czynność | — |
| Zalogowany operator | Plan i rozliczenia | czynność | — |
| Zalogowany operator | Urządzenia i sesje logowania | czynność | — |
| Zalogowany operator | Bezpieczeństwo i uwierzytelnianie dwuskładnikowe | czynność | — |
| Zalogowany operator | Powiadomienia konta | czynność | — |
| Zalogowany operator | Status: dostępny | wybór jednokrotny | — |
| Zalogowany operator | Status: zajęty | wybór jednokrotny | — |
| Zalogowany operator | Status: niewidoczny | wybór jednokrotny | — |
| Zalogowany operator | Pomoc i wsparcie | czynność | — |
| Zalogowany operator | Wyloguj | czynność | — |

**Status obecności** jest wyborem jednokrotnym i widzą go pozostali uczestnicy pracy zespołowej. Status „niewidoczny” nie wstrzymuje pracy zadań w tle.

---

## 15. Menu strefy pracy i rozwinięcie „Więcej”

### 15.1. Menu „Nowa sesja”

Kontrolka: znak `+`, pierwsza pozycja strefy pracy szyny.

| Sekcja | Pozycja | Rodzaj |
|---|---|---|
| Nowa sesja — w którym środowisku | TalkIn · WorkSpace · CodeStudio · MultitaskingAI | czynność, jedna z czterech |
| — | Powtórz ostatnią sesję (`Ctrl+Shift+N`) | czynność |

Menu zamyka zdanie wyjaśniające, że sesja otwiera przedsionek środowiska,
a moduł wiodący wybiera się dopiero kaflem w przedsionku.

### 15.2. Menu „Nowy projekt”

Kontrolka: katalog ze znakiem `+`, trzecia pozycja strefy pracy.

| Sekcja | Pozycja | Rodzaj |
|---|---|---|
| Nowy projekt — w którym środowisku | TalkIn · WorkSpace · CodeStudio · MultitaskingAI | czynność, jedna z czterech |
| — | Otwórz projekt istniejący (`Ctrl+O`) | czynność |
| — | Utwórz z szablonu projektu | czynność |

Menu zamyka zdanie o tym, że po wskazaniu środowiska otwiera się od razu jego
Workspace w trybie zakładania projektu, z pominięciem przedsionka.

### 15.3. Rozwinięcie „Więcej”

Kontrolka: `⋯`, przedostatnia pozycja wstążki, tuż przed przyciskiem „Dostosuj”.

| Sekcja | Pozycja | Skrót |
|---|---|---|
| Znajdź w widoku | pole wyszukiwania, zakres, rozróżnianie wielkości liter, znajdź i zamień | `Ctrl+F`, `Ctrl+H` |
| Pomoc | Dokumentacja platformy | `F1` |
| Pomoc | Skróty klawiszowe | `Ctrl+/` |
| Pomoc | Zgłoś obserwację | — |

---

## 16. Przypisywanie skrótów klawiszowych

Trzy menu — zrzut ekranu, schowek i znajdowanie w widoku — zawierają pola
przypisania skrótu. Pole przyjmuje kombinację naciśniętą przez Operatora.

| Etap | Zachowanie |
|---|---|
| Ustawienie kursora w polu | Zapis dotychczasowy ustępuje wezwaniu „naciśnij kombinację…” |
| Naciśnięcie kombinacji | Klawisze modyfikujące same nie kończą wprowadzania; kombinacja zapisuje się dopiero z klawiszem znakowym |
| Kombinacja wolna | Zapis przyjęty; komunikat potwierdza przypisanie |
| Kombinacja zajęta | Zapis **nie zostaje przyjęty**; komunikat wskazuje czynność, która ją zajmuje |
| Klawisz `Escape` | Zapis poprzedni wraca; pole traci kursor |

**Zero blokad.** Pole pozostaje aktywne również wtedy, gdy kombinacja jest
zajęta — Operator dostaje komunikat i może wprowadzić inną, nie napotyka
kontrolki wyłączonej.

**Konflikty ze skrótami systemowymi** rozstrzyga system operacyjny; platforma
odnotowuje kombinację, lecz nie przechwytuje jej wcześniej niż system.

---

## 17. Zachowanie na progach szerokości

| Próg | Belka tytułowa | Szyna nawigacji | Pasek edycji | Pasek stanu |
|---|---|---|---|---|
| powyżej 1400 px | pełny skład | 48 px, pozycje 40 px | pełny skład; pole wyszukiwania 340 px | pełny skład |
| 1400 px | bez zmian | bez zmian | pole wyszukiwania zwęża się w kierunku minimum 180 px | bez zmian |
| 1200 px | bez zmian | bez zmian | rodzina „widok i powiadomienia” ustępuje; napis przycisku „Dostosuj” zwija się do ikony | bez zmian |
| 960 px | tytuł widoku ustępuje; znak i sterowanie oknem pozostają | pozycje 36 px | wszystkie rodziny prawej strony ustępują; ich czynności pozostają w menu aplikacji | położenie, CPU i RAM ustępują |

**Kolejność ustępowania.** Pierwsze ustępują rodziny prawej strony, ostatnie —
pasmo kart sesji. Karty są pracą Operatora; ikony funkcji mają swoje miejsce
także w menu aplikacji.

**Próg wysokości okna.** Poniżej 900 px szyna nawigacji przewija się, a etykiety
pozycji ustępują: przewijanie wymaga przycinania, a przycięty kontener nie
wypuszcza etykiety poza swoją krawędź.

**Reguła ustępowania.** Element może zniknąć z pasa, nie może zniknąć
z systemu: każda czynność usunięta z pasa przy węższym oknie pozostaje
wywoływalna z menu aplikacji albo skrótem klawiszowym.

---

## 18. Warstwy widoczności i dostępność

| Element | Warstwa | Uzasadnienie |
|---|---|---|
| Belka tytułowa wraz ze sterowaniem oknem | 1 | Widoczna bez interakcji w każdym oknie |
| Szyna nawigacji — środowiska | 1 | Widoczne bez interakcji w każdym oknie |
| Pasek edycji wraz z rodzinami kontrolek | 1 | Widoczny bez interakcji w każdym oknie |
| Pasek stanu | 1 | Widoczny bez interakcji w każdym oknie |
| Szyna nawigacji — moduły środowiska | 2 | Ujawniane rozwinięciem środowiska |
| Etykieta kontrolki i etykieta pozycji szyny | 2 | Ujawniane wskazaniem kursorem albo fokusem |
| Menu rozwijane | 2 | Ujawniane naciśnięciem kontrolki |
| Pola przypisania skrótu | 3 | Wewnątrz menu, po jego rozwinięciu |

**Dostępność.**

| Wymaganie | Rozwiązanie |
|---|---|
| Etykieta dla czytnika ekranu | Każda kontrolka ma `aria-label` powtarzający treść etykiety wizualnej |
| Ikona dekoracyjna | Każdy znak graficzny wewnątrz kontrolki ma `aria-hidden="true"` |
| Stan menu | `aria-expanded` na kontrolce; `role="menu"` na treści; `role="menuitem"`, `menuitemcheckbox` albo `menuitemradio` na pozycjach |
| Nawigacja klawiaturą | Cały pas osiągalny klawiszem `Tab`; menu zamyka się klawiszem `Escape` |
| Fokus widoczny | Pierścień 2 px w barwie `--dn-fokus` z odsunięciem 2 px — nigdy usuwany |
| Kontrast etykiety | Powierzchnia treści z obrysem, nie tło ramy — czytelność zachowana w obu motywach |
| Ograniczony ruch | `prefers-reduced-motion` wygasza uniesienie kontrolki i przejście etykiety |

---

## 19. Pole wyszukiwania i filtr zakresu

Pole wyszukiwania na wstążce jest szersze niż pojedyncza kontrolka — mieści do
540 px, bo przeszukuje **sesję, projekt
i środowisko** naraz, a nie tylko bieżący widok (tym zajmuje się osobne menu
„Znajdź w widoku”, rozdz. 10).

Przy prawej krawędzi pola stoi **filtr zakresu** (`.dn-szukaj-filtr`) z napisem
„Wszystkie”. Rozwija listę środowisk do zaznaczenia: Operator zawęża
wyszukiwanie do wybranych środowisk, nie opuszczając pola. Domyślnie zaznaczone
są wszystkie cztery; odznaczenie któregoś wyłącza je z wyników. Skrót
`Alt + D` ustawia kursor w polu z dowolnego miejsca aplikacji. Szerokość pola
wynosi 340 px i zwęża się do minimum 180 px (`.dn-narzedzia-szukaj`,
`design/zasoby/rama.css`).

---

## 20. Personalizacja pasów

Okno „Dostosuj” (rozdz. 7) ma czwartą zakładkę — **Personalizacja widoku** —
obejmującą wygląd **obu pasów naraz**: wstążki poziomej i szyny pionowej.

| Grupa | Rozstrzygnięcie |
|---|---|
| Kolorystyka pasów | atrament · grafit · powierzchnia · sygnał przygaszony |
| Kolor podświetlenia ikon | błękit sygnałowy · zieleń · bursztyn · neutralny |
| Rodzaj separacji pozycji | kreska między rodzinami · sam odstęp · tło rodziny |
| Akcja po kliknięciu pozycji | wykonaj od razu · podpowiedź, potem wykonaj · rozwiń menu wariantów |
| Odstępy i pozycjonowanie | puste odstępy do rozsuwania ikon · wyrównanie rodzin do krawędzi |
| Grupowanie w jedną ikonę | ikony zbiorcze z menu „⋯” · zwijanie rzadko używanych pozycji szyny · usuwanie pozycji szybkiego wyboru |

**Barwy pochodzą z palety systemu**, nie z dowolnego wyboru koloru: każda opcja
pozostaje czytelna w motywie jasnym i ciemnym. **Ikony zbiorcze** pozwalają
spiąć kilka pozycji jednego pasa pod jedną ikoną, która po kliknięciu otwiera
małe menu ze spiętymi pozycjami — na przykład trzy rzadko używane moduły
szybkiego wyboru chowają się pod jedną ikoną rozwijaną. Moduł zdjęty
z szybkiego wyboru pozostaje dostępny w przedsionku środowiska.

**Dostosowanie szyny** wywołuje przycisk `Dostosuj pasek` w stopce szyny
(`data-dostosuj-szyne`), pod ikoną ustawień i konfiguracji — otwiera to samo
okno „Dostosuj”, w którym szyna pionowa jest jedną z trzech kolumn przenoszenia
składników (rozdz. 7.1).

---

## 21. Okna nakładkowe jako okna

Wszystkie okna nakładkowe platformy (`<dialog class="dn-modal">` — Dostosuj,
Ustawienia i konfiguracja, Historia sesji, Instrukcja użytkowania i inne) niosą
**ramę okna**: przeciąganie, zmianę rozmiaru oraz minimalizację
i maksymalizację. Zachowanie dokłada `design/zasoby/okna-modalne.js` (przyciski
`.dn-btn` w `.dn-modal-stopka`, wiersz `.dn-modal-naglowek`), aktywne wyłącznie
na dialogach oznaczonych `data-okno="1"`.

| Czynność | Uchwyt | Uwaga |
|---|---|---|
| Przeciąganie | nagłówek `.dn-modal-naglowek` | zawsze co najmniej 80 px nagłówka pozostaje w oknie |
| Zmiana rozmiaru | osiem uchwytów przy krawędziach i narożnikach | minimum 420 × 260 px |
| Minimalizacja | przycisk `--min` albo kliknięcie nagłówka | okno dokuje przy dolnej krawędzi od prawej |
| Maksymalizacja | przycisk `--max` albo dwuklik nagłówka | wypełnia okno z marginesem 16 px |

Klawiaturą: gdy nagłówek ma fokus, strzałki przesuwają okno, a `Shift`+strzałki
zmieniają rozmiar (krok 16 px); `Escape` zamyka. Stan (położenie, rozmiar,
zminimalizowanie) trzymają atrybuty `data-okno-*`. Przy `prefers-reduced-motion`
przejścia ustają.

---

## 22. Wykazy zamykające i kryteria odbioru

Rozdział zamyka dokument czterema wykazami obowiązkowymi wskazanymi w
[Standardzie redakcyjnym i językowym](../STANDARD-REDAKCYJNY-I-JEZYKOWY.md)
rozdz. 4.2 — komendy kontraktu, żetony, komponenty i skróty klawiszowe — oraz
dwoma wykazami specyficznymi dla katalogu `interfejs-uzytkownika/` — stany
kontrolek i punkty łamania — a kończy go zbiór kryteriów odbioru.

### 22.1. Wykaz komend kontraktu

Rama okna sama nie przenosi danych — komendy poniżej obsługują okna
nakładkowe otwierane z ramy: menu „Ustawienia i konfiguracja” (rozdz. 13)
oraz nakładkę „Historia sesji” wywoływaną ze strefy pracy szyny (rozdz. 3.5).
Źródło: `budowa/shared/contract.json`, obszary `session` i `config`.

| Komenda | Obszar | Pola żądania | Pola wyniku | Zdarzenia | Kody błędów |
|---|---|---|---|---|---|
| `config.window.open` | `config` | `area?: SessionConfigArea` · `scope?: ConfigScope` · `sessionId?: string` · `windowId?: string` | `opened: bool` · `area?: SessionConfigArea` · `scope?: ConfigScope` | — | `validation_failed` |
| `session.list` | `session` | `status?: SessionStatus` · `limit?: int` · `offset?: int` · `includePresence?: bool` | `sessions: Session[]` · `total: int` · `presence?: SessionPresence[]` | `session.changed` | `validation_failed` |
| `session.archive.list` | `session` | `offset?: int` · `limit?: int` | `sessions: Session[]` · `total: int` | `session.changed` | `validation_failed` |
| `session.resume` | `session` | `sessionId: string` | `session: Session` · `windows: Window[]` | `session.changed` · `session.focus.changed` | `not_found` · `conflict` |
| `session.close` | `session` | `sessionId: string` | `session: Session` | `session.changed` | `not_found` |
| `session.archive` | `session` | `sessionIds: string[]` | `archivedIds: string[]` | `session.changed` | `not_found` |
| `session.restore` | `session` | `sessionIds: string[]` | `restoredIds: string[]` | `session.changed` | `not_found` |
| `session.rename` | `session` | `sessionId: string` · `title: string` | `session: Session` | `session.changed` | `validation_failed` · `not_found` |
| `session.copy` | `session` | `sessionId: string` · `title?: string` | `session: Session` · `copiedMessages: int` | `session.changed` | `not_found` |
| `session.delete` | `session` | `sessionIds: string[]` · `confirm: bool` | `deletedIds: string[]` · `deletedCount: int` · `missingIds?: string[]` | `session.changed` | `validation_failed` · `not_found` |

**Uwaga o zdarzeniach.** `contract.json` nie wiąże jawnie zdarzenia
z komendą, która je wywołuje — kolumna „Zdarzenia” wskazuje zdarzenie
`session.changed` (obszar `session`, zmiana sesji) jako jedyne zdarzenie
zwrotne obszaru, którego nasłuchuje szyna sesji i nakładka „Historia sesji”
przy każdej z ośmiu komend mutujących listę. `session.resume` dokłada
`session.focus.changed`, ponieważ wznowienie przenosi ognisko na wznowioną
kartę.

### 22.2. Wykaz żetonów

Każdy żeton poniżej występuje w treści dokumentu i jest zweryfikowany wobec
`design/zasoby/zetony/zetony.css`. Dwa żetony omówione w rozdz. 1 i rozdz.
4.3 — `--dn-wym-szyna` i `--dn-tło` — nie istnieją w źródle; poprawiono je na
`--dn-wym-belka` (żeton dzielony) i `--dn-tlo` (zapis bez polskiego znaku,
zgodny z konwencją identyfikatorów kodu).

| Żeton | Wartość źródłowa | Zastosowanie w ramie |
|---|---|---|
| `--dn-wym-belka` | `48px` | wysokość belki tytułowej i paska górnego; szerokość szyny nawigacji (żeton dzielony) |
| `--dn-wym-pasek` | `48px` | wysokość paska edycji |
| `--dn-wym-stan` | `28px` | wysokość paska stanu |
| `--dn-wym-ikonowy` | `34px` (`40px` na wskazaniu dotykowym i w gęstości przestronnej) | wymiar kontrolki `.dn-nrz-btn` |
| `--dn-wym-ikona` | `18px` | ikona wewnątrz `.dn-nrz-btn` i `.dn-szyna-poz` (środowisko) |
| `--dn-wym-ikona-sm` | `16px` | ikona wewnątrz `.dn-szyna-poz--modul` |
| `--dn-r-lg` | `10px` | zaokrąglenie tacy paska edycji |
| `--dn-czas-1` | `0,1 s` | mikroreakcje — wskazanie kursorem, naciśnięcie |
| `--dn-czas-2` | `0,16 s` | przejścia barw, ujawnienie etykiety |
| `--dn-panel` | `#FAFAFA` (motyw jasny) | tło paska edycji i paska stanu |
| `--dn-tlo` | `#F4F4F4` (motyw jasny) | odniesienie przy uzasadnieniu żetonów dzielonych wstążki |
| `--dn-powierzchnia` | `#FFFFFF` (motyw jasny) | tło etykiety pozycji i dymka podpowiedzi |
| `--dn-powierzchnia-2` | `#ECECEC` (motyw jasny) | tło kontrolki paska edycji przy otwartym menu |
| `--dn-tekst` | `#181818` (motyw jasny) | barwa ikony kontrolki przy wskazaniu kursorem i przy menu otwartym |
| `--dn-tekst-2` | `#616161` (motyw jasny) | barwa ikony kontrolki w spoczynku |
| `--dn-rama-tekst` | `#ECECEC` (motyw ciemny bazowy ramy) | barwa ikon pozycji szyny nawigacji |
| `--dn-rama-hover` | `rgba(255,255,255,0.08)` | tło pozycji szyny i przycisków belki przy wskazaniu kursorem |
| `--dn-rama-obrys-mocny` | `rgba(255,255,255,0.24)` | ramka środowiska rozwiniętego i wskazanego kursorem |
| `--dn-sygnal-wypelnienie` | pochodna `--dn-sygnal-600` | wypełnienie modułu bieżącego w szynie |
| `--dn-obrys-mocny` | `#C0C0C0` (motyw jasny) | kreska pionowa rozdzielająca rodziny kontrolek paska edycji |
| `--dn-hover` | `rgba(10,10,10,0.05)` | tło kontrolki `.dn-nrz-btn` przy wskazaniu kursorem |
| `--dn-fokus` | pochodna `--dn-sygnal-500` | pierścień fokusu wszystkich kontrolek ramy |
| `--dn-blad-tekst` | pochodna `--dn-czerwien-700` | tło przycisku „Zamknij okno” przy wskazaniu kursorem |
| `--dn-wstazka` | `#FFFFFF` (motyw jasny) | belka nakładana paska edycji |
| `--dn-wstazka-taca` | `#E3E3E3` (motyw jasny) | taca pod belką nakładaną |

### 22.3. Wykaz komponentów

Wszystkie klasy poniżej są zweryfikowane wobec `design/zasoby/rama.css`,
`design/zasoby/css/komponenty.css`, `design/zasoby/prototyp.css` i
`design/zasoby/okna-modalne.css` — kompletny zbiór klas `.dn-*` i `.sta-*`
użytych w tym dokumencie, bez braków.

| Klasa | Rola | Modyfikatory | Stany |
|---|---|---|---|
| `.dn-belka` | pas belki tytułowej | — | — |
| `.dn-belka-marka` | sygnet i nazwa produktu | — | — |
| `.dn-belka-nazwa` | napis „Danaco Console” | — | — |
| `.dn-belka-tytul` | tytuł widoku bieżącego | — | ukryty poniżej 960 px |
| `.dn-belka-okno` | kontener sterowania oknem | — | — |
| `.dn-belka-btn` | przycisk sterowania oknem systemowym | `--zamknij` | spoczynek · wskazanie kursorem · fokus |
| `.dn-szyna-nawigacji` | pas szyny nawigacji | — | — |
| `.dn-szyna-nawigacji-lista` | lista środowisk strefy 2 | — | — |
| `.dn-szyna-sekcja` | grupa strefy pracy | — | — |
| `.dn-szyna-skroty` | grupa szybkiego wyboru, strefa 3 | — | — |
| `.dn-szyna-stopka` | stopka szyny — ustawienia, motyw, profil | — | — |
| `.dn-szyna-tresc` | kontener przewijany trzech stref | — | przewijanie poniżej 900 px wysokości |
| `.dn-szyna-poz` | pozycja szyny — klasa bazowa | `--srodowisko` · `--modul` · `--praca` · `--dodaj` | spoczynek · wskazanie kursorem · fokus · rozwinięte · bieżące |
| `.dn-szyna-poz--srodowisko` | pozycja środowiska | — | zwinięte · rozwinięte (`aria-expanded`) · bieżące (`data-biezace`) |
| `.dn-szyna-poz--modul` | pozycja modułu | — | spoczynek · bieżący (`aria-current`) |
| `.dn-szyna-poz--praca` | pozycja strefy pracy (Nowa sesja, Historia, Nowy projekt) | — | spoczynek · wskazanie kursorem |
| `.dn-narzedzia` | pas paska edycji | — | — |
| `.dn-narzedzia-karty` | pasmo kart otwartych sesji | — | puste · przewijane · przycięte |
| `.dn-narzedzia-dostosuj` | przycisk „Dostosuj” zamykający pasek | — | spoczynek · wskazanie kursorem |
| `.dn-nrz-btn` | kontrolka ikonowa paska edycji | — | spoczynek · wskazanie kursorem · naciśnięcie · fokus · menu otwarte (`aria-expanded`) · dwustanowa (`aria-pressed`) |
| `.dn-stan` | pas paska stanu | — | — |
| `.dn-szukaj-filtr` | filtr zakresu pola wyszukiwania | — | spoczynek · rozwinięty |
| `.dn-modal` | okno nakładkowe niosące ramę | `--szeroki` | normalny · zminimalizowane · zmaksymalizowane (`data-okno-stan`) |
| `.dn-btn` w `.dn-modal-stopka` | przycisk sterowania oknem nakładkowym | `--min` · `--max` | spoczynek · wskazanie kursorem |
| `.dn-modal-naglowek` | uchwyt przeciągania okna nakładkowego | — | fokus (klawiaturowe przesuwanie) |
| bez odrębnej klasy w arkuszu | uchwyt zmiany rozmiaru przy krawędziach i narożnikach okna nakładkowego (osiem sztuk) | `--n` `--s` `--e` `--w` `--ne` `--nw` `--se` `--sw` | spoczynek · przeciąganie |
| `.sta-podmenu` | kontener podmenu menu aplikacji | — | zwinięte · rozwinięte (`:hover`, `:focus-within`) |
| `.sta-podmenu-tresc` | treść rozwiniętego podmenu | — | ukryta · widoczna |

### 22.4. Wykaz skrótów klawiszowych

Wykaz konsoliduje wszystkie kombinacje przywołane w rozdziałach 9–15 —
menu aplikacji, menu „Znajdź w widoku”, menu „Zrzut ekranu”, menu „Schowek”
oraz menu „Ustawienia i konfiguracja”. Kolizje odnotowuje kolumna „Kolizja”;
puste pole znaczy brak kolizji zarejestrowanej w tym dokumencie.

| Kombinacja | Działanie | Zasięg | Kolizja |
|---|---|---|---|
| `Ctrl+N` | Nowa sesja | globalny | — |
| `Ctrl+Shift+N` | Nowe okno aplikacji · Powtórz ostatnią sesję (menu „Nowa sesja”) | globalny | dwie czynności różnych menu dzielą kombinację |
| `Ctrl+O` | Otwórz projekt… · Otwórz projekt istniejący (menu „Nowy projekt”) | globalny | dwie czynności różnych menu dzielą kombinację |
| `Ctrl+S` | Zapisz sesję | globalny | — |
| `Ctrl+Shift+E` | Eksportuj sesję… | globalny | — |
| `Ctrl+Z` | Cofnij | widok | — |
| `Ctrl+Y` | Ponów | widok | — |
| `Ctrl+C` | Kopiuj | widok | — |
| `Ctrl+V` | Wklej · Wklej ostatnią pozycję schowka | widok | dwie czynności różnych menu dzielą kombinację |
| `Ctrl+F` | Znajdź w widoku | widok | — |
| `Ctrl+H` | Znajdź i zamień | widok | — |
| `Ctrl+Shift+B` | Szyna nawigacji — przełącznik widoczności | globalny | — |
| `Ctrl+B` | Szyna sesji — przełącznik widoczności | globalny | — |
| `Ctrl+Shift+F` | Tryb skupienia | globalny | — |
| `F11` | Pełny ekran | globalny | koliduje z pełnym ekranem systemu operacyjnego — rozstrzyga system |
| `Ctrl+=` | Powiększ | widok | — |
| `Ctrl+-` | Pomniejsz | widok | — |
| `F5` | Odśwież widok · Ponowne wczytanie widoku (sekcja „Dla programistów”) | widok | dwie pozycje menu aplikacji dzielą kombinację |
| `Ctrl+Home` | Centrum dowodzenia | globalny | — |
| `Ctrl+E` | Przedsionek środowiska | globalny | — |
| `Ctrl+Tab` | Następna karta sesji | widok | — |
| `Ctrl+Shift+Tab` | Poprzednia karta sesji | widok | — |
| `Ctrl+K` | Paleta poleceń | globalny | — |
| `Alt+D` | Ustawienie kursora w polu wyszukiwania paska edycji | globalny | — |
| `F1` | Dokumentacja platformy | globalny | — |
| `Ctrl+/` | Skróty klawiszowe | globalny | — |
| `Ctrl+,` | Ustawienia aplikacji | globalny | — |
| `PrtSc` | Zrzut całego okna aplikacji | globalny | koliduje ze skrótem systemowym zrzutu ekranu — rozstrzyga system |
| `Ctrl+PrtSc` | Zrzut zaznaczonego obszaru | globalny | — |
| `Alt+PrtSc` | Zrzut wskazanego okna roboczego | globalny | — |
| `Ctrl+Shift+V` | Otwórz schowek | globalny | — |
| `Ctrl+Shift+T` | Wklej jako zwykły tekst | widok | — |
| `Ctrl+Alt+V` | Wklejenie z listy schowka (przypisywalne) | globalny | — |

**Pola przypisania skrótu.** Trzy kombinacje — skrót zrzutu obszaru
(domyślnie `Ctrl+PrtSc`), skrót zrzutu okna (domyślnie `Alt+PrtSc`) i
wklejenie z listy (domyślnie `Ctrl+Alt+V`) — Operator zmienia w polach
przypisania opisanych w rozdz. 16; wartości w tabeli są domyślne fabrycznie.

### 22.5. Wykaz stanów kontrolek

Osiem stanów za standardem rozdz. 4.2. Kontrolki nawigacyjne ramy nie noszą
danych, więc stany „ładowanie”, „pusty” i „błąd” ich nie dotyczą z wyjątkiem
filtra zakresu i pasma kart sesji, które prezentują zbiory zmienne w czasie —
w pozostałych wierszach komórka niesie adnotację „nie dotyczy” wraz z
przyczyną, zgodnie z rozdz. 4.2 standardu.

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| `.dn-nrz-btn` | tło przezroczyste, ikona `--dn-tekst-2` | tło `--dn-hover`, ikona `--dn-tekst`, uniesienie 1 px | powrót do położenia wyjściowego | pierścień `--dn-fokus` 2 px z odsunięciem | `[disabled]`, kursor domyślny | nie dotyczy — czynność natychmiastowa, bez odpowiedzi sieci | nie dotyczy | nie dotyczy — błąd czynności zgłasza dymek powiadomienia, nie stan kontrolki |
| `.dn-belka-btn` | tło przezroczyste, ikona `--dn-rama-tekst-2` | tło `--dn-rama-hover`, ikona `--dn-rama-tekst` (`--zamknij`: tło `--dn-blad-tekst`) | brak odrębnego stanu — działanie wykonuje się na puszczeniu | pierścień `--dn-fokus`, odsunięcie ujemne o `--dn-wym-fokus` | nie dotyczy — trzy kontrolki systemowe zawsze dostępne | nie dotyczy | nie dotyczy | nie dotyczy |
| `.dn-szyna-poz--srodowisko` | bez ramki, tło przezroczyste | ramka `--dn-rama-obrys-mocny` | brak odrębnego stanu wizualnego | pierścień `--dn-fokus` | nie dotyczy — cztery środowiska zawsze dostępne | nie dotyczy | nie dotyczy | nie dotyczy |
| `.dn-szyna-poz--modul` | bez ramki, tło przezroczyste | tło `--dn-rama-hover` | brak odrębnego stanu wizualnego | pierścień `--dn-fokus` | nie dotyczy | nie dotyczy | nie dotyczy | nie dotyczy |
| `.dn-szukaj-filtr` | napis „Wszystkie”, tło przezroczyste | tło `--dn-rama-hover` | rozwija listę środowisk | pierścień `--dn-fokus` | nie dotyczy | nie dotyczy | wszystkie środowiska odznaczone — komunikat „Wybierz co najmniej jedno środowisko” | nie dotyczy |
| `.dn-narzedzia-karty` (pasmo kart sesji) | karty w kolejności otwarcia | karta wskazana kursorem unosi wskaźnik zamknięcia | przełączenie ogniska na kartę | pierścień `--dn-fokus` na karcie aktywnej klawiaturą | nie dotyczy | nie dotyczy — otwarcie karty jest lokalne, okno komunikacji wczytuje się wewnątrz karty, nie na karcie samej | brak kart otwartych — pasmo puste, bez komunikatu | nie dotyczy |
| `.dn-btn` w `.dn-modal-stopka` | tło przezroczyste | tło `--dn-rama-hover` | wykonanie czynności (minimalizacja/maksymalizacja) | pierścień `--dn-fokus` | nie dotyczy | nie dotyczy | nie dotyczy | nie dotyczy |
| pole przypisania skrótu (rozdz. 16) | zapis dotychczasowej kombinacji | nie dotyczy — pole nie reaguje na wskazanie kursorem | tekst „naciśnij kombinację…” po ustawieniu kursora | pierścień `--dn-fokus` | nie dotyczy | nie dotyczy | brak kombinacji przypisanej — pole puste z tekstem podpowiedzi | kombinacja zajęta — komunikat wskazujący czynność, która ją zajmuje |

### 22.6. Wykaz punktów łamania

Rama korzysta z dwóch zestawów progów: rodziny żetonów globalnych
`--dn-bp-*` (siatka platformy, `zetony.css`) oraz **własnych progów ramy** —
1400 px, 1200 px i 960 px w szerokości oraz 900 px w wysokości. Z rodziną
żetonów globalnych pokrywa się wyłącznie próg 960 px; pozostałe są progami
własnymi ramy i taki stan jest wiążący.

| Żeton / próg | Wartość | Pochodzenie | Zachowanie ramy |
|---|---|---|---|
| `--dn-bp-w1` | `640px` | żeton globalny | poniżej progu — widok mobilny; poza zakresem tego dokumentu, patrz [Ustawienia](ustawienia.md) |
| `--dn-bp-w2` | `960px` | żeton globalny | zgodny z progiem ramy 960 px (rozdz. 17): tytuł widoku ustępuje, pozycje szyny zwężają się do 36 px, rodziny prawej strony paska edycji ustępują, pasek stanu traci pozycje drugorzędne |
| `--dn-bp-w3` | `1280px` | żeton globalny | brak progu ramy dedykowanego temu żetonowi |
| `--dn-bp-w4` | `1600px` | żeton globalny | brak progu ramy dedykowanego temu żetonowi |
| próg ramy 1400 px | `1400px` | próg własny ramy, rozdz. 17 | pole wyszukiwania zwęża się z 340 px w kierunku minimum 180 px |
| próg ramy 1200 px | `1200px` | próg własny ramy, rozdz. 17 | rodzina „widok i powiadomienia” ustępuje; napis „Dostosuj” zwija się do ikony |
| próg wysokości 900 px | `900px` | próg własny ramy, rozdz. 17 | szyna nawigacji zaczyna przewijać się; etykiety pozycji ustępują |

**Wyjątek zapisany wprost.** Progi 1400 px, 1200 px i 960 px oraz próg
wysokości 900 px są **własnymi progami ramy okna** i stanowią świadomy wyjątek
od rodziny żetonów `--dn-bp-*`, nie usterkę do ujednolicenia. Rama jest jedynym
bytem platformy, którego skład zmienia się wraz z ubywaniem miejsca na pasy
poziome, a nie wraz z przełamaniem siatki treści — dlatego jej progi wypadają
w innych miejscach niż progi siatki. Rodzina `--dn-bp-w1…w4` (640 px, 960 px,
1280 px, 1600 px) pozostaje rodziną progów platformy w niezmienionej postaci;
progów ramy nie przelicza się na żetony ani nie dopisuje do tej rodziny.

### 22.7. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Szerokość szyny nawigacji równa się wysokości belki tytułowej w każdym motywie | pomiar w narzędziu deweloperskim na `--dn-wym-belka` obu elementów |
| Rozwinięte jest najwyżej jedno środowisko naraz | rozwinięcie drugiego środowiska przy pierwszym już rozwiniętym — pierwsze zwija się automatycznie |
| Każda kontrolka usunięta z paska edycji przy zwężeniu okna pozostaje wywoływalna z menu aplikacji albo skrótem klawiszowym | zwężenie okna poniżej 960 px, sprawdzenie każdej ukrytej rodziny w menu aplikacji |
| Trzy okna nakładkowe (Dostosuj, Ustawienia i konfiguracja, Historia sesji) przeciągają się, zmieniają rozmiar, minimalizują i maksymalizują | test manualny trzech okien wraz z klawiaturą (`Tab`, strzałki, `Shift`+strzałki, `Escape`) |
| Pole przypisania skrótu odmawia zapisania kombinacji zajętej i wskazuje czynność, która ją zajmuje | próba przypisania `Ctrl+F` w polu skrótu zrzutu ekranu |
| Etykieta kontrolki i etykieta pozycji szyny ustępują, gdy menu tej kontrolki jest otwarte | otwarcie menu „Zrzut ekranu” przy kursorze wskazującym kontrolkę |
| Przy `prefers-reduced-motion` uniesienie kontrolki i przejście etykiety ustają, a informacja pozostaje dostępna | test z włączonym ograniczeniem ruchu systemu operacyjnego |
| Każdy żeton i każda klasa przywołane w dokumencie istnieją w `zetony.css`, `rama.css`, `komponenty.css`, `prototyp.css` albo `okna-modalne.css` | zestawienie automatyczne nazw przywołanych w dokumencie z plikami źródłowymi (wykonane przy redakcji niniejszego wydania — wynik rozdz. 22.2–22.3) |
| Trzy drogi do nakładki „Historia sesji” (strefa pracy szyny, przycisk przedsionka, sekcja ustawień) otwierają dokładnie tę samą nakładkę z tym samym stanem początkowym | otwarcie z każdej z trzech dróg, porównanie zaznaczonej sekcji i zawartości |
| Status obecności w menu „Profil użytkownika” jest widoczny pozostałym uczestnikom pracy zespołowej | test dwuklientowy — zmiana statusu na jednym kliencie, odczyt na drugim |

---

## Załącznik. Pełny wykaz komend kontraktu ramy okna

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `window` — 7 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `window.create` | Zakłada okno komunikacji w sesji | `sessionId:string` (wym)<br>`moduleId:string` (wym)<br>`modelChannelId:string` (wym)<br>`workingDirs:string[]` (wym)<br>`executionEnv:ExecutionEnv` (wym)<br>`permissionMode:PermissionMode` (wym)<br>`windowRole:WindowRole` (wym)<br>`coordinatorWindowId:string` (opc)<br>`title:string` (opc) | `window:Window` (wym) |
| `window.list` | Zwraca okna komunikacji | `sessionId:string` (opc)<br>`status:WindowStatus` (opc) | `windows:Window[]` (wym) |
| `window.state.get` | Zwraca stan okna komunikacji: parametry wykonania, stan procesu i konfigurację efektywną | `windowId:string` (wym)<br>`includeConfig:bool` (opc)<br>`includeAccess:bool` (opc) | `window:Window` (wym)<br>`processStatus:ProgressStatus` (wym)<br>`messageCount:int` (wym)<br>`lastMessageId:string` (opc)<br>`streaming:bool` (wym)<br>`config:ConfigEntry[]` (opc)<br>`accessGrants:AccessGrant[]` (opc)<br>`loop:LoopState` (opc) |
| `window.update` | Zmienia parametry okna komunikacji | `windowId:string` (wym)<br>`moduleId:string` (opc)<br>`modelChannelId:string` (opc)<br>`agentId:string` (opc)<br>`workingDirs:string[]` (opc)<br>`executionEnv:ExecutionEnv` (opc)<br>`permissionMode:PermissionMode` (opc)<br>`windowRole:WindowRole` (opc)<br>`coordinatorWindowId:string` (opc)<br>`title:string` (opc) | `window:Window` (wym) |
| `window.close` | Zamyka okno komunikacji | `windowId:string` (wym) | `window:Window` (wym) |
| `window.action` | Wykonuje akcję panelu akcji okna operacyjnego; wspólną każdemu oknu | `windowId:string` (wym)<br>`actionId:string` (wym)<br>`parameters:json` (opc) | `windowId:string` (wym)<br>`actionId:string` (wym)<br>`result:json` (opc) |
| `window.handoff` | Przekazuje zlecenie z okna koordynatora do okna wykonawcy; zakłada pozycję kolejki | `sessionId:string` (wym)<br>`fromWindowId:string` (wym)<br>`toWindowId:string` (wym)<br>`instruction:string` (wym)<br>`bundle:ContextBundle` (opc) | `window:Window` (wym)<br>`queueItemId:string` (wym) |

**Zdarzenia obszaru `window` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `window.changed` | Zmiana okna komunikacji | — |
| `window.state.changed` | Zmiana stanu okna operacyjnego; wspólna każdemu oknu | — |

### Obszar `session` — 19 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `session.create` | Zakłada sesję | `title:string` (opc)<br>`projectId:string` (opc)<br>`metadata:json` (opc) | `session:Session` (wym) |
| `session.list` | Zwraca listę sesji | `status:SessionStatus` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc)<br>`includePresence:bool` (opc) | `sessions:Session[]` (wym)<br>`total:int` (wym)<br>`presence:SessionPresence[]` (opc) |
| `session.focus` | Przenosi ognisko na wskazaną kartę sesji i opcjonalnie na okno w jej wnętrzu. Asystent przestawia ognisko Operatorowi polem `targetClientId` — inaczej jego posunięcia nie byłyby widoczne na ekranie | `sessionId:string` (wym)<br>`clientId:string` (wym)<br>`windowId:string` (opc)<br>`targetClientId:string` (opc) | `sessionId:string` (wym)<br>`windowId:string` (opc)<br>`previousSessionId:string` (opc)<br>`focusedAt:int64` (wym) |
| `session.bind` | Wiąże bieżące połączenie z sesją trwającą na rdzeniu; odtwarza jej okna | `sessionId:string` (wym)<br>`clientId:string` (wym)<br>`windowIds:string[]` (opc) | `session:Session` (wym)<br>`windows:Window[]` (wym)<br>`bound:bool` (wym)<br>`resumed:bool` (wym) |
| `session.open` | Otwiera sesję wraz z jej oknami | `sessionId:string` (wym) | `session:Session` (wym)<br>`windows:Window[]` (wym) |
| `session.close` | Zamyka sesję | `sessionId:string` (wym) | `session:Session` (wym) |
| `session.delete` | Usuwa trwale wskazane sesje wraz z całym ich zapisem; jedyna droga utraty danych sesji | `sessionIds:string[]` (wym)<br>`confirm:bool` (wym) | `deletedIds:string[]` (wym)<br>`deletedCount:int` (wym)<br>`missingIds:string[]` (opc) |
| `session.rename` | Zmienia nazwę sesji w historii | `sessionId:string` (wym)<br>`title:string` (wym) | `session:Session` (wym) |
| `session.copy` | Kopiuje sesję wraz z jej zapisem; kopia jest osobnym bytem historii | `sessionId:string` (wym)<br>`title:string` (opc) | `session:Session` (wym)<br>`copiedMessages:int` (wym) |
| `session.project.set` | Przenosi sesję do projektu; projekt wskazany albo zakładany nowy | `sessionIds:string[]` (wym)<br>`projectId:string` (opc)<br>`projectName:string` (opc) | `projectId:string` (wym)<br>`movedIds:string[]` (wym) |
| `session.project.clear` | Wyjmuje sesję z projektu; sesja zostaje w historii bez projektu | `sessionIds:string[]` (wym) | `clearedIds:string[]` (wym) |
| `session.archive` | Przenosi sesje do archiwum; zapis zostaje w całości, sesja znika z historii bieżącej | `sessionIds:string[]` (wym) | `archivedIds:string[]` (wym) |
| `session.restore` | Przywraca sesje z archiwum do historii bieżącej | `sessionIds:string[]` (wym) | `restoredIds:string[]` (wym) |
| `session.archive.list` | Zwraca sesje archiwum; wgląd z okna ustawień | `offset:int` (opc)<br>`limit:int` (opc) | `sessions:Session[]` (wym)<br>`total:int` (wym) |
| `session.resume` | Wznawia zamkniętą sesję wraz z jej oknami | `sessionId:string` (wym) | `session:Session` (wym)<br>`windows:Window[]` (wym) |
| `session.stop` | Zatrzymuje tury biegnące we wszystkich oknach sesji | `sessionId:string` (wym) | `sessionId:string` (wym)<br>`stoppedWindowIds:string[]` (wym) |
| `session.tool.attach` | Dokłada narzędzie albo umiejętność do sesji na czas jej trwania. Definicji eksperta nie zmienia: dołożenie żyje w stanie sesji, przeżywa rozłączenie klienta i kończy się wraz z sesją albo z usunięciem rozmowy | `sessionId:string` (wym)<br>`toolName:string` (wym)<br>`source:SessionToolSource` (opc) | `tool:SessionTool` (wym)<br>`tools:SessionTool[]` (wym)<br>`alreadyAttached:bool` (wym) |
| `session.tool.detach` | Zdejmuje dołożenie z sesji. Zestaw wraca do podstawy z definicji eksperta; sama definicja pozostaje niezmieniona | `sessionId:string` (wym)<br>`toolName:string` (wym) | `detached:bool` (wym)<br>`tools:SessionTool[]` (wym) |
| `session.tool.list` | Zwraca narzędzia dołożone do sesji. Zestaw narzędzi tury to definicja eksperta wraz z tymi dołożeniami — bez tego odczytu druga połowa zestawu byłaby niewidoczna | `sessionId:string` (wym) | `tools:SessionTool[]` (wym)<br>`total:int` (wym) |

**Zdarzenia obszaru `session` — 4:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `session.changed` | Zmiana sesji | — |
| `session.focus.changed` | Zmiana ogniska karty sesji; nośnik synchronizacji ogniska między urządzeniami konta | — |
| `session.tool.attached` | Narzędzie dołożone do sesji. Skoro model dostaje narzędzie, którego nie miał, Operator musi to zobaczyć — każde posunięcie widać na ekranie | — |
| `session.tool.detached` | Dołożenie zdjęte z sesji. Widoczność obowiązuje w obie strony: zestaw, który urósł na oczach Operatora, nie może skurczyć się po cichu — bez tego zdarzenia sąsiednie okno pokazywałoby narzędzie, którego model już nie ma | — |

### Obszar `home` — 1 komenda

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `home.enter` | Wejście na stronę główną; zwraca karty środowisk i sesje czynne konta | `clientId:string` (wym) | `environments:Environment[]` (wym)<br>`sessions:Session[]` (wym)<br>`focusedSessionId:string` (opc)<br>`presence:SessionPresence[]` (opc) |

### Obszar `config` — 1 komenda

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `config.window.open` | Otwiera okno konfiguracji na wskazanym zakresie i zwraca jego stan wyjściowy; komendę przywołuje menu „Ustawienia i konfiguracja” (rozdz. 13) | `area:SessionConfigArea` (opc)<br>`scope:ConfigScope` (opc)<br>`sessionId:string` (opc)<br>`windowId:string` (opc) | `opened:bool` (wym)<br>`area:SessionConfigArea` (opc)<br>`scope:ConfigScope` (opc) |

Razem w wykazie: **28 komend** z czterech obszarów kontraktu.

---

*Koniec dokumentu. Rama okna aplikacji — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
