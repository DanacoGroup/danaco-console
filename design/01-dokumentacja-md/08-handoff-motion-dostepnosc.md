# Danaco Console — Handoff, ruch i dostępność

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — kontrakt przekazania do wdrożenia, kontrakt ruchu, kontrakt dostępności |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | **Deweloper wdrażający** warstwę wizualną do `budowa/client` · **Projektant** prowadzący makiety i prototypy · **Odbierający jakość** (brama weryfikacyjna fali) |
| **Zakres** | (I) mapowanie pakietu design na strukturę kodu, zmiany kontraktowe, cztery fale wdrożenia; (II) katalog ruchu dozwolonego i zakazanego, wzorce implementacyjne, ograniczony ruch, wydajność; (III) WCAG 2.1 AA, komplet 33 pomiarów kontrastu, fokus, klawiatura, semantyka, czytniki ekranu, listy kontrolne |
| **Czego NIE zawiera** | definicji żetonów (→ `04-tokens`), anatomii komponentów (→ `06-components`), makiet okien (→ `05-okna/`), księgi znaku (→ `09-brand-system`), reguł typografii i siatki (→ `01-design-system`) |

---

## Spis treści

**CZĘŚĆ I — HANDOFF (przekazanie do wdrożenia)**

1. [Odbiorca, zakres, regulamin budowy](#1-odbiorca-zakres-regulamin-budowy)
2. [Pełna tabela mapowania plików pakietu](#2-pełna-tabela-mapowania-plików-pakietu)
3. [Zmiany kontraktowe w kodzie](#3-zmiany-kontraktowe-w-kodzie)
4. [Mapowanie ról kontraktu WebSocket na klasy `.dn-wpis--*`](#4-mapowanie-ról-kontraktu-websocket-na-klasy-dn-wpis--)
5. [Cztery fale wdrożenia i brama weryfikacyjna](#5-cztery-fale-wdrożenia-i-brama-weryfikacyjna)
6. [Czego nie wolno zmienić po cichu](#6-czego-nie-wolno-zmienić-po-cichu)
7. [Lista kontrolna odbioru wizualnego fali](#7-lista-kontrolna-odbioru-wizualnego-fali)

**CZĘŚĆ II — MOTION (ruch i animacja)**

8. [Kontrakt ruchu: `INTENSYWNOSC_RUCHU` 3/10](#8-kontrakt-ruchu-intensywnosc_ruchu-310)
9. [Cztery czasy i jeden ease](#9-cztery-czasy-i-jeden-ease)
10. [Katalog dozwolonych animacji](#10-katalog-dozwolonych-animacji)
11. [Katalog zakazanych animacji](#11-katalog-zakazanych-animacji)
12. [Wzorce implementacyjne](#12-wzorce-implementacyjne)
13. [`prefers-reduced-motion` — co dokładnie się dzieje](#13-prefers-reduced-motion--co-dokładnie-się-dzieje)
14. [Ruch a wydajność](#14-ruch-a-wydajność)

**CZĘŚĆ III — DOSTĘPNOŚĆ**

15. [WCAG 2.1 AA jako warunek wejściowy](#15-wcag-21-aa-jako-warunek-wejściowy)
16. [Kontrast: metoda, komplet 33 pomiarów, zasada tekst-3](#16-kontrast-metoda-komplet-33-pomiarów-zasada-tekst-3)
17. [Fokus](#17-fokus)
18. [Klawiatura](#18-klawiatura)
19. [Stan nigdy samym kolorem](#19-stan-nigdy-samym-kolorem)
20. [Semantyka](#20-semantyka)
21. [Cele dotykowe i `pointer: coarse`](#21-cele-dotykowe-i-pointer-coarse)
22. [Czytniki ekranu](#22-czytniki-ekranu)
23. [Lista kontrolna dostępności przed oddaniem okna](#23-lista-kontrolna-dostępności-przed-oddaniem-okna)
24. [Decyzje projektowe](#24-decyzje-projektowe)

---

# CZĘŚĆ I — HANDOFF

## 1. Odbiorca, zakres, regulamin budowy

### 1.1. Kto odbiera i co odbiera

| Rola | Co odbiera | Czym się rozlicza |
|---|---|---|
| **Deweloper wdrażający** | pakiet plików `zasoby/` (żetony, fundament, komponenty, ikony, fonty, marka) | działającą aplikacją na nowej palecie, bez `#000000`, bez `disabled`, w obu motywach |
| **Projektant prowadzący** | kontrakt ruchu i kontrakt dostępności | makietami i prototypami, które przechodzą listę kontrolną rozdz. 7 i 23 |
| **Odbierający jakość** | bramy weryfikacyjne czterech fal | porównaniem z makietami w obu motywach + przebiegiem `kontrasty.json` |

**Źródło kontraktu:** `design/opracowania/design/HANDOFF.md` (v2.0, 2026-08-11), rozszerzone o realny stan pakietu `WYNIK/zasoby/` policzony plik po pliku.

### 1.2. Regulamin budowy — cztery reguły nienaruszalne

| # | Reguła | Konsekwencja praktyczna |
|---|---|---|
| **R1** | **Jedna odpowiedzialność = jeden plik** | żeton barwy nie mieszka w pliku komponentu; komponent nie definiuje żetonu |
| **R2** | **Arkusz ≤ 300 linii** | `komponenty.css` (1207 linii) rozpada się na 14 plików per komponent; `zetony.css` (448 linii) na 11 plików per warstwa |
| **R3** | **Komponenty sięgają wyłącznie po żetony semantyczne** | `.dn-btn` używa `--dn-atrament`, nigdy `--dn-szary-900`; prymityw jest zakazany poza plikami motywu |
| **R4** | **Wartości przenosi się dosłownie** | podział pliku to operacja mechaniczna — zero „przy okazji poprawiłem odcień" |

### 1.3. Kierunek zależności arkuszy

```
   fonty.css            ← @font-face, licencje OFL, latin + latin-ext
       │
       ▼
   motyw/ (żetony)      ← prymitywy → semantyczne (jasny | ciemny) → niezależne
       │                  ŻADEN plik motywu nie zna klasy .dn-*
       ▼
   fundament.css        ← reset, podłoże, fokus globalny, wzorce tekstowe
       │                  zna wyłącznie var(--dn-*)
       ▼
   komponenty/*.css     ← biblioteka .dn-*
       │                  zna wyłącznie żetony SEMANTYCZNE
       ▼
   okna / powłoka       ← składanie z klas .dn-*; zero nowych wartości barw
```

**Zasada odwrotna jest błędem wdrożenia:** jeżeli plik komponentu definiuje `--dn-*`, warstwa została przekroczona.

---

## 2. Pełna tabela mapowania plików pakietu

### 2.1. Żetony → `budowa/client/src/motyw/`

Podział `zasoby/zetony/zetony.css` (448 linii) na jedenaście arkuszy. Kolumna „Sekcja źródła" wskazuje numerowane sekcje komentarzy w pliku wejściowym.

| Plik docelowy | Sekcja źródła | Zawartość | Szac. linii |
|---|---|---|---|
| `motyw/prymitywy.css` | 1 | skala szarości (18 stopni), rodzina sygnału (8), stany (12) | ~40 |
| `motyw/rama.css` | 2 | `--dn-rama`, `-tekst`, `-tekst-2`, `-hover`, `-obrys` | ~10 |
| `motyw/typografia.css` | 3 | trzy rodziny, 9 stopni, 4 wagi, 3 interlinie, 3 odstępy liter | ~30 |
| `motyw/przestrzen.css` | 4 (+ cienie z 10/11) | 11 odstępów, 6 promieni, komplety cieni per motyw | ~30 |
| `motyw/ruch.css` | 5 + 14 | `--dn-ease`, cztery czasy, blok `prefers-reduced-motion` | ~25 |
| `motyw/wymiary.css` | 6 + 12 + 13 | wymiary gęstości zwartej, blok `pointer: coarse`, blok `data-gestosc` | ~50 |
| `motyw/siatka.css` | 7 | punkty łamania w1–w4, `--dn-tresc-max`, kolumny, przerwa | ~10 |
| `motyw/warstwy.css` | 8 | 11 poziomów `--dn-z-*` | ~15 |
| `motyw/gradienty.css` | 9 | `--dn-grad-atrament`, `--dn-grad-sygnal` | ~8 |
| `motyw/semantyczne-jasny.css` | 10 + zapas `prefers-color-scheme: light` | role motywu jasnego | ~90 |
| `motyw/semantyczne-ciemny.css` | 11 + zapas `prefers-color-scheme: dark` | role motywu ciemnego | ~90 |
| `motyw/motyw.css` | — | **arkusz spinający** — wyłącznie `@import` w kolejności powyżej | ~15 |

**Uwaga do sekcji 7 (siatka).** `HANDOFF.md` nie wymienia osobnego pliku siatki. Sekcja jest zbyt mała, by stanowić plik, i zbyt obca, by leżeć w `przestrzen.css`. Rozstrzygnięcie: osobny `siatka.css` — odnotowane w rozdz. 24.

### 2.2. Pozostałe pliki żetonów i fontów

| Plik pakietu | Cel | Sposób | Liczba plików |
|---|---|---|---|
| `zasoby/zetony/zetony.json` | `motyw/zetony.json` | podmiana w całości | 1 |
| `zasoby/zetony/kontrasty.json` | `motyw/kontrasty.json` | **nowy plik** — wejście bramy weryfikacyjnej | 1 |
| `zasoby/zetony/fonty.css` | `motyw/fonty.css` | `@font-face` z `unicode-range`, `font-display: swap` | 1 |
| `zasoby/zetony/fonty/*.woff2` | `motyw/fonty/` | podzbiory `latin` + `latin-ext` | **20** |
| `zasoby/zetony/fonty/LICENCJA-*.txt` | `motyw/fonty/` | OFL 1.1 — IBM Plex, Space Grotesk | 2 |

Rozkład plików kroju: **IBM Plex Sans** 4 wagi × 2 podzbiory = 8 · **IBM Plex Mono** 3 wagi × 2 = 6 · **Space Grotesk** 3 wagi × 2 = 6.

### 2.3. CSS bazowy i komponenty

| Plik pakietu | Cel | Uwagi |
|---|---|---|
| `zasoby/css/fundament.css` (156 linii) | `motyw/fundament.css` | import **po** żetonach, **przed** komponentami; zawiera fokus globalny `:focus-visible` — **usunąć duplikat z dotychczasowego `motyw.css`** |
| `zasoby/css/komponenty.css` (1207 linii) | `komponenty/` — podział na 14 arkuszy | selektory `.dn-*` **bez zmian**; nazwy plików zachowane tam, gdzie się pokrywają z istniejącymi |

Podział `komponenty.css` zgodnie z progiem 300 linii (R2):

| # | Plik docelowy | Klasy przenoszone | Zakres linii źródła |
|---|---|---|---|
| 1 | `komponenty/przycisk.css` | `.dn-btn` (+ 8 modyfikatorów), `.dn-btn-ikona` (+ `--na-ramie`) | 14–166 |
| 2 | `komponenty/pole.css` | `.dn-pole` (`-etykieta -kontrolka -opis -blad`), `.dn-szukaj` | 169–234 |
| 3 | `komponenty/wybor.css` | `.dn-wybor`, `.dn-check`, `.dn-radio`, `.dn-przelacznik`, `.dn-suwak` | 238–333 |
| 4 | `komponenty/plakietka.css` | `.dn-kropka` (+ 5 wariantów), `.dn-plakietka` (+ 6 wariantów) | 335–393 |
| 5 | `komponenty/karta.css` | `.dn-karta`, `.dn-karta-srodowiska`, `.dn-kafel` | 395–524 |
| 6 | `komponenty/tabela.css` | `.dn-tabela` | 526–560 |
| 7 | `komponenty/zakladki.css` | `.dn-zakladki`, `.dn-zakladka`, `.dn-karty-sesji`, `.dn-karta-sesji` | 562–639 |
| 8 | `komponenty/pasek.css` | `.dn-pasek` (+ 5 elementów), `.dn-listwa` | 641–721 |
| 9 | `komponenty/boczna.css` | `.dn-boczna` (+ `-naglowek -pozycja`) | 723–780 |
| 10 | `komponenty/wpis.css` | `.dn-wpis` (+ 4 modyfikatory, 5 elementów) | 782–862 |
| 11 | `komponenty/postep.css` | `.dn-postep`, `.dn-kolejka`, `.dn-krok` (+ 4 modyfikatory) | 864–945 |
| 12 | `komponenty/nakladka.css` | `.dn-modal` (+ 4 elementy), `@keyframes dn-wejscie` | 947–987 |
| 13 | `komponenty/powiadomienie.css` | `.dn-toasty`, `.dn-toast` (+ 4 warianty) | 989–1026 |
| 14 | `komponenty/drobne.css` | `.dn-tooltip`, `.dn-awatar`, `.dn-spinner`, `.dn-pusty-stan`, `.dn-prompt`, `.dn-przybornik`, `.dn-aod` | 1028–1208 |

**Współdzielone klatki kluczowe.** `dn-tetno` żyje w `plakietka.css` (definicja przy `.dn-kropka`), a używają jej także `wpis.css` (`.dn-wpis--pracuje`) i `drobne.css` (`.dn-aod-rdzen`). `dn-wejscie` żyje w `nakladka.css`, używa jej także `powiadomienie.css`. `dn-obrot` żyje w `drobne.css` (`.dn-spinner`), używa jej także `przycisk.css` (`.dn-btn[aria-busy='true']::after`). **Kolejność importu w arkuszu spinającym musi zapewnić obecność definicji** — najprościej: `plakietka.css` → `wpis.css`, `nakladka.css` → `powiadomienie.css`, `drobne.css` przed lub po (klatki kluczowe nie podlegają kaskadzie kolejności, ale porządek czytelności ma znaczenie dla utrzymania).

### 2.4. Ikony, znak i zasoby marki

| Plik pakietu | Cel | Uwagi | Liczba |
|---|---|---|---|
| `zasoby/ikony/svg/*.svg` | `ikony/svg/` | **podmiana całego katalogu** (poprzedni zestaw = zastępczy); nazwy polskie, obrys 1,75, siatka 24×24 | **82** |
| `zasoby/ikony/manifest.json` | `ikony/manifest.json` | nowy plik; pola: `nazwa`, `zrodlo`, `zastosowanie` | 1 |
| `zasoby/marka/srodowiska/*.svg` | `ikony/svg/` | emblematy 4 środowisk — używane przez Centrum dowodzenia i pasek ramy | **4** |
| `zasoby/marka/logo/sygnet*.svg` | `ikony/svg/logo-danaco.svg` **oraz** `zasoby/marka/` | godło zachowuje barwy własne w obu motywach; wariant uproszczony poniżej 24 px | 6 |
| `zasoby/marka/logo/logotyp*.svg`, `logo-poziomy*.svg`, `logo-pionowy*.svg` | `zasoby/marka/logo/` | komplet wersji: na jasnym, na ciemnym, mono biały, mono czarny | 8 |
| `zasoby/marka/logo/png/*` | `zasoby/marka/logo/png/` | rastry @2x i 512 px do zastosowań poza wektorem | 6 |
| `zasoby/marka/favicon/*` | `client/public/` (Vite) | snippet `<head>` gotowy w `naglowek-snippet.html`; `site.webmanifest` | 10 |
| `zasoby/marka/ikona-aplikacji/*` | `desktop/` (Tauri: 32, 128, 128@2x z `ikona-1024.png`) + `client/public/` (PWA) | wariant **maskowalny** obowiązkowy dla manifestu | 7 |

### 2.5. Czego się NIE przenosi

| Plik pakietu | Dlaczego zostaje |
|---|---|
| `zasoby/wspolne.js` | warstwa makiet (przełącznik motywu + zapas Invoker Commands) — aplikacja ma własny mechanizm motywu |
| `zasoby/prototyp.css`, `zasoby/prototyp.js` | rusztowanie prototypów dokumentacyjnych, nie biblioteka produktu |
| `zasoby/WZORZEC-OKNA.html` | szablon makiety, materiał projektanta |
| `zasoby/marka/logo/alternatywy/*` | odrzucone koncepcje znaku — archiwum decyzji, nie zasób produkcyjny |

---

## 3. Zmiany kontraktowe w kodzie

Pięć zmian, których nie da się wykonać wyłącznie podmianą arkuszy.

### 3.1. `ikony/zrodla-ikon.ts` + `ikony.ts`

| Element | Stan docelowy |
|---|---|
| Wykaz nazw | zaktualizowany do `manifest.json` — **82 pozycje**, nazewnictwo polskie opisowe |
| Typowanie | literówka w nazwie ikony **ma dalej zatrzymywać kompilację** — wykaz jest typem, nie łańcuchem znaków |
| Kolor | `currentColor` — komponent nadaje barwę, ikona jej nie zna |
| Renderowanie | 14 / 16 / 20 / 24 px (żetony `--dn-wym-ikona-*`) |

> **Uwaga rozbieżności.** `HANDOFF.md` v2.0 podaje „82 pozycje". Zliczenie katalogu `zasoby/ikony/svg/` daje **82 pliki** — zgodność potwierdzona pomiarem, nie deklaracją.

### 3.2. `okno-komunikacji/okno.css`

| Krok | Działanie |
|---|---|
| 1 | **usunąć lokalne zmienne `--dc-*`** — okno komunikacji przestaje mieć własną paletę |
| 2 | przepiąć wpisy rozmowy na `.dn-wpis` z modyfikatorami `--czlowiek` / `--inteligencja` / `--system` oraz `--pracuje` |
| 3 | przepiąć pole wpisywania na `.dn-prompt` (`-grot`, `-obszar`) + `.dn-przybornik` |
| 4 | przepiąć wskaźnik postępu na `.dn-postep` (`-etykieta`, `-tor`, `-wartosc`) |
| 5 | zastosować mapowanie ról kontraktu (rozdz. 4) |

### 3.3. Motyw

| Element | Stan docelowy |
|---|---|
| Mechanizm | **bez zmian**: atrybut `data-theme` na `<html>` + zapas `prefers-color-scheme` |
| Wartości | nowe — z `motyw/semantyczne-jasny.css` i `semantyczne-ciemny.css` |
| `color-scheme` | ustawiane per motyw (już obecne w żetonach: `light` / `dark`) |
| Pasek górny | **zawsze atramentowy** (`--dn-rama`) — jedyny element nieprzełączający się z motywem |
| Równoprawność | oba motywy definiowane osobno, **nie wywodzone** jeden z drugiego |

### 3.4. Gęstość

| Element | Stan docelowy |
|---|---|
| Domyślna | **zwarta** — kontrolka 32 px, wiersz 36 px, pasek 48 px |
| Przełącznik | przyszłościowo `data-gestosc="przestronna"` na `<html>` — **żetony gotowe**, sekcja 13 `zetony.css` |
| Dotyk | `pointer: coarse` podnosi wymiary automatycznie (sekcja 12) — bez reguł per komponent |

### 3.5. Zero blokad (ADL-017)

| Zakaz | Zamiast tego |
|---|---|
| atrybut `disabled` na kontrolce | pełna klikalność + **komunikat po naciśnięciu** albo **opis obok** |
| krycie 0,5 + kursor „niedozwolone" | plakietka wyjaśniająca (`.dn-plakietka--ostrzezenie`) lub dymek (`.dn-tooltip`) |
| brama „najpierw uzupełnij pola" | przycisk klikalny; kliknięcie ujawnia braki komunikatem `.dn-pole-blad` |
| odliczanie blokujące „Wyślij ponownie" | odliczanie **wyłącznie informacyjne** — kliknięcie działa zawsze |
| wyszarzona metoda uwierzytelniania | metoda wyłączona w konfiguracji **nie jest renderowana** — „brak metody = mniej segmentów, nie zablokowany segment" |

**Jedyny dopuszczalny wyjątek:** `readonly` na polu formularza w widoku faktycznie tylko do odczytu (podgląd prowenancji, archiwum sesji). To nie jest brama — element fizycznie nie ma trybu edycji. Klasa `.dn-pole-kontrolka[readonly]` jest w bibliotece przewidziana.

---

## 4. Mapowanie ról kontraktu WebSocket na klasy `.dn-wpis--*`

### 4.1. Tabela mapowania

| Rola w kontrakcie | Klasa `.dn-wpis--*` | Medalion | Kreska lewa | Tło | Plakietka roli |
|---|---|---|---|---|---|
| `user` | `--czlowiek` | ikona `uzytkownik`, obrys atramentowy | `--dn-atrament` | powierzchnia | — |
| `assistant` | `--inteligencja` | ikona `agent`, tło `--dn-sygnal-tlo` | `--dn-sygnal-wypelnienie` | powierzchnia | nazwa modelu (przykładowo: `CLAUDE-OPUS`) |
| `system` | `--system` | ikona `info` | `--dn-obrys-mocny` | **przezroczyste** (wycofane) | — |
| `tool` | `--system` | ikona `terminal` | `--dn-obrys-mocny` | **przezroczyste** | **wymagana** — `.dn-plakietka--rola` = nazwa narzędzia |

### 4.2. Dlaczego `tool` dzieli klasę z `system`

Anty-domyślne kierunku zakazuje „dziewięciu kolorów tła dla dziewięciu nadawców". Kontrakt ma cztery role, ale nadawców w oknie komunikacji bywa więcej (wiele modeli, wiele narzędzi, wiele ról MultitaskingAI). Rozstrzygnięcie systemu: **trzy klasy semantyczne + ikona + etykieta + plakietka roli.** Rozróżnienie `tool` od `system` niesie **plakietka**, nie barwa — dzięki temu liczba nadawców może rosnąć bez rozrostu palety.

```
   .dn-wpis--czlowiek        .dn-wpis--inteligencja      .dn-wpis--system
   ┌─┬────────────────┐      ┌─┬────────────────┐        ┌─┬────────────────┐
   │▌│ ◉ Operator     │      │▌│ ◈ Asystent  ●  │        │▌│ ⓘ System       │
   │▌│ 14:22          │      │▌│ CLAUDE-OPUS    │        │▌│ [WEB-SEARCH]   │  ← plakietka roli
   │▌│ treść wpisu    │      │▌│ treść wpisu    │        │▌│ treść wpisu    │     (rola tool)
   └─┴────────────────┘      └─┴────────────────┘        └─┴────────────────┘
    ▲ atrament                ▲ sygnał   ▲ kropka         ▲ obrys mocny
                                          tętna (--pracuje)
```

### 4.3. Modyfikator `--pracuje`

Nakłada się na dowolną z trzech klas. Realizacja: `.dn-wpis--pracuje .dn-wpis-nadawca::after` — kropka 6 px w barwie `--dn-kropka` z animacją `dn-tetno`. **Nie zmienia klasy semantycznej wpisu** — praca to stan, nie rola.

### 4.4. Reguła nadawcy nierozpoznanego

Rola spoza kontraktu → `.dn-wpis--system` + plakietka roli z surową nazwą roli. **Nigdy nie renderować wpisu bez klasy semantycznej** — wpis bez kreski lewej czyta się jak awaria stylów.

---

## 5. Cztery fale wdrożenia i brama weryfikacyjna

### 5.1. Oś wdrożenia

```
  FALA 1              FALA 2               FALA 3              FALA 4
  ŻETONY              KOMPONENTY           IKONY I MARKA       OKNA
  ──────►             ──────►              ──────►             ──────►
  motyw/              komponenty/          ikony/              okno-komunikacji
  11 arkuszy          14 arkuszy           82 SVG              + powłoka
  + fonty             + fundament          + favicon           wg makiet
                                           + Tauri
  ┌──────────┐        ┌──────────┐         ┌──────────┐        ┌──────────┐
  │ BRAMA 1  │        │ BRAMA 2  │         │ BRAMA 3  │        │ BRAMA 4  │
  └──────────┘        └──────────┘         └──────────┘        └──────────┘
  aplikacja rusza     galeria              zero ikon           makiety okien
  na nowej palecie;   komponenty.html      zastępczych;        w obu motywach;
  komponenty stare    = wzorzec odbioru    kompilacja          role kontraktu
  jeszcze działają                         wykazu nazw         na klasach .dn-wpis
```

### 5.2. Zakres i brama każdej fali

| Fala | Zakres | Warunek wejścia | **Brama weryfikacyjna** |
|---|---|---|---|
| **1 · Żetony** | `motyw/` — 11 arkuszy + `zetony.json` + `kontrasty.json` + fonty (20 woff2 + 2 licencje) | brak | (a) aplikacja uruchamia się w obu motywach; (b) **przebieg `kontrasty.json`** — 33 pomiary, pomiar a nie deklaracja; (c) pasek górny atramentowy w obu motywach; (d) `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` renderuje się poprawnie we wszystkich trzech krojach |
| **2 · Komponenty** | `komponenty/` — 14 arkuszy + `fundament.css` | fala 1 zamknięta | (a) porównanie z galerią `06-okna/komponenty.html` w obu motywach; (b) zero `disabled` w drzewie DOM; (c) fokus widoczny na każdej kontrolce; (d) zero wartości szesnastkowych poza plikami motywu |
| **3 · Ikony i marka** | `ikony/svg/` (82 + 4 emblematy), `manifest.json`, favicon, ikony Tauri, PWA | fala 2 zamknięta | (a) kompilacja przechodzi po podmianie wykazu nazw; (b) zero ikon z poprzedniego (zastępczego) zestawu; (c) godło w wariancie uproszczonym poniżej 24 px; (d) ikona maskowalna obecna w manifeście |
| **4 · Okna** | `okno-komunikacji/` + powłoka wg makiet Centrum dowodzenia, powłoki środowiska, relacji koordynator–wykonawca | fala 3 zamknięta | (a) porównanie z makietami w obu motywach; (b) mapowanie ról kontraktu (rozdz. 4) sprawdzone na czterech rolach; (c) tętno kropki obecne dokładnie tam, gdzie biegnie praca; (d) lista kontrolna rozdz. 7 zaliczona |

### 5.3. Reguła bramy

> **Brama weryfikacyjna każdej fali: porównanie z makietami w obu motywach oraz przebieg `kontrasty.json` (pomiar, nie deklaracja).**

Trzy konsekwencje operacyjne:

1. **Oba motywy zawsze.** Zaliczenie w jednym motywie nie jest zaliczeniem. Motywy są równoprawne, więc bramę przechodzi się dwukrotnie.
2. **Pomiar, nie deklaracja.** „Kontrast wygląda dobrze" nie jest wynikiem. Wynikiem jest liczba z algorytmu WCAG (rozdz. 16.1) i werdykt wobec progu.
3. **Fala nie zamyka się częściowo.** Fala 2 nie startuje przed zamknięciem fali 1 — komponenty odczytują żetony, których jeszcze nie ma.

---

## 6. Czego nie wolno zmienić po cichu

| # | Zamrożone | Co się psuje przy cichej zmianie | Kto może zmienić |
|---|---|---|---|
| **Z1** | **Wartości żetonów** (`zetony.css` / `zetony.json` jako źródło prawdy) | rozjazd między pakietem a kodem; 33 pomiary kontrastu przestają obowiązywać | Właściciel + aktualizacja `kontrasty.json` |
| **Z2** | **Geometria znaku i emblematów** (krzywe, nie fonty) | znak przestaje być znakiem; wersje mono i uproszczona się rozjeżdżają | Właściciel |
| **Z3** | **Zasada jednego akcentu** — sygnał ≤ 5% ekranu, **nigdy tło sekcji** | kokpit przestaje wskazywać, gdzie biegnie praca; sygnał traci znaczenie | nikt — kontrakt kierunku |
| **Z4** | **Pierścień fokusu** (2 px + odsunięcie 2 px, `--dn-fokus`) i **zachowanie `prefers-reduced-motion`** | wypadnięcie z WCAG 2.1 AA (2.4.7); ruch wraca użytkownikom, którzy go wyłączyli | nikt — warunek wejściowy |
| **Z5** | **Gęstość zwarta jako domyślna** | rytm kokpitu; wszystkie makiety przestają pasować | Właściciel (przełącznik `data-gestosc` istnieje) |

**Dodatkowo zamrożone na mocy KANONU:**

| # | Zamrożone | Uzasadnienie |
|---|---|---|
| **Z6** | **Nazwy klas `.dn-*`** | wolno rozszerzać o nowe modyfikatory; **nie wolno zmieniać istniejących nazw** — są kontraktem między projektem a kodem |
| **Z7** | **Zakaz `#000000`** | skala kończy się na `#0A0A0A`; czysta czerń zabija głębię i haluje na OLED |
| **Z8** | **Zero blokad (ADL-017)** | decyzja rejestru; `disabled` nie wraca jako „drobne usprawnienie" |
| **Z9** | **Nazwy własne okien i modułów** (`Studio Editor`, `Workflow Builder`, `Chat Window`, `Coordinator`, `Executor 1`…) | zakaz tłumaczenia i parafrazowania |
| **Z10** | **Pasek górny atramentowy w obu motywach** | jedyny element nieprzełączający się z motywem — rama kokpitu |

---

## 7. Lista kontrolna odbioru wizualnego fali

Do przejścia **w obu motywach**, na każdej fali. Pozycja niezaliczona = fala niezamknięta.

### 7.1. Barwa i żetony

- [ ] Zero wartości szesnastkowych poza plikami `motyw/prymitywy.css`, `semantyczne-*.css`
- [ ] Zero `#000000` w całym drzewie
- [ ] Komponenty sięgają wyłącznie po żetony semantyczne (zero `--dn-szary-*` w `komponenty/`)
- [ ] Sygnał zajmuje ≤ 5% powierzchni ekranu i nie jest tłem żadnej sekcji
- [ ] Pasek górny atramentowy w motywie jasnym i ciemnym
- [ ] Gradient nie występuje jako tło przycisku, karty ani sekcji

### 7.2. Typografia

- [ ] Nagłówki krojem Space Grotesk, interfejs Plex Sans, dane Plex Mono
- [ ] Liczby w `tabular-nums`
- [ ] Diakrytyki `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` poprawne we wszystkich trzech krojach i wagach
- [ ] Stopień bazowy 13 px; brak stopni spoza skali

### 7.3. Przestrzeń i wymiary

- [ ] Wszystkie odstępy są krotnością 4 px
- [ ] Kontrolka 32 px, wiersz 36 px, pasek 48 px (gęstość zwarta)
- [ ] Promienie ze skali `--dn-r-*`, bez wartości własnych

### 7.4. Stany i zero blokad

- [ ] Zero atrybutów `disabled` w drzewie DOM
- [ ] Każdy stan niesie ikonę albo etykietę — nigdy sam kolor
- [ ] Niegotowość sygnalizowana opisem obok albo komunikatem po naciśnięciu
- [ ] Przycisk „Pomiń" w oknie logowania obecny i klikalny

### 7.5. Ruch

- [ ] Zero animacji dekoracyjnych, parallaxu i scrollytellingu
- [ ] Jeden ruch znaczący na widok (tętno) — dokładnie tam, gdzie biegnie praca
- [ ] Wszystkie przejścia w czasach `--dn-czas-*` i `--dn-ease`
- [ ] `prefers-reduced-motion` sprawdzone: tętno → pierścień statyczny

### 7.6. Dostępność

- [ ] Fokus widoczny na każdej kontrolce interaktywnej
- [ ] Pełne okno przechodzone klawiaturą, bez pułapki fokusu poza modalem
- [ ] Ikony dekoracyjne `aria-hidden="true"`, informacyjne z etykietą
- [ ] `lang="pl"` na dokumencie
- [ ] Przebieg 33 pomiarów kontrastu bez pozycji „NIE"

### 7.7. Ikony i marka

- [ ] Wyłącznie ikony z `ikony/svg/` — zero emoji, zero ikon spoza zestawu
- [ ] Obrys 1,75, `currentColor`, `fill: none`
- [ ] Godło w wariancie uproszczonym poniżej 24 px

---

# CZĘŚĆ II — MOTION

## 8. Kontrakt ruchu: `INTENSYWNOSC_RUCHU` 3/10

### 8.1. Zdanie kontraktowe

> `INTENSYWNOSC_RUCHU` = **3/10** — mikroprzejścia 100–220 ms + **jeden ruch znaczący** (tętno kropki). Aplikacja robocza, nie strona marketingowa.

### 8.2. Co 3/10 znaczy operacyjnie

| Wymiar | Wartość kontraktowa | Test zaliczenia |
|---|---|---|
| **Górny czas przejścia** | **220 ms** | żadna animacja interfejsu nie trwa dłużej; wyjątek: tętno 2,4 s (ruch ciągły) i obrót spinnera 0,8 s (ruch cykliczny) |
| **Dolny czas przejścia** | **100 ms** | poniżej tej wartości ruch jest nieczytelny — lepiej zmienić stan natychmiast niż animować niewidocznie |
| **Liczba ruchów ciągłych na widok** | **1** | tętno kropki w miejscu, gdzie biegnie praca. Dwa tętna w jednym widoku = błąd projektu |
| **Ruch bez informacji** | **0** | każda animacja odpowiada na pytanie „co się właśnie zmieniło w systemie?" |
| **Dystans przesunięcia** | **≤ 8 px** (`--dn-od-2`) | wejście warstwy: `translateY(8px)` → 0. Większy dystans czyta się jak strona marketingowa |
| **Skala** | **0,98 → 1** | wyłącznie wejście modala; brak „bounce", brak przeskoku ponad 1 |

### 8.3. Trzy pytania przed dodaniem animacji

```
   ┌───────────────────────────────────────────────┐
   │  1. Czy ruch niesie informację o stanie       │  NIE ─► NIE ANIMUJ
   │     systemu?                                  │
   └────────────────────┬──────────────────────────┘
                        │ TAK
   ┌────────────────────▼──────────────────────────┐
   │  2. Czy mieści się w 100–220 ms?              │  NIE ─► skróć albo
   └────────────────────┬──────────────────────────┘         zrezygnuj
                        │ TAK
   ┌────────────────────▼──────────────────────────┐
   │  3. Czy jest w katalogu dozwolonych           │  NIE ─► ZAKAZANA
   │     (rozdz. 10)?                              │         (rozdz. 11)
   └────────────────────┬──────────────────────────┘
                        │ TAK
                    ► DOZWOLONA
```

---

## 9. Cztery czasy i jeden ease

### 9.1. Krzywa

| Żeton | Wartość | Charakter |
|---|---|---|
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | szybki start, miękkie dojście do celu — ruch instrumentu, nie sprężyny |

**Jedna krzywa dla całego systemu.** Zakaz `ease-in-out`, `linear` (poza obrotem spinnera), `cubic-bezier` własnych i wszelkich krzywych „sprężystych" z przeregulowaniem.

### 9.2. Cztery czasy

| Żeton | Wartość | Kiedy | Przykłady zastosowania |
|---|---|---|---|
| `--dn-czas-1` | **0,1 s** | **mikroreakcja kontrolki** — reakcja na kursor i naciśnięcie | `.dn-btn:hover`, `.dn-btn:active`, `.dn-tabela tbody tr:hover`, `.dn-check::before` |
| `--dn-czas-2` | **0,16 s** | **przejście barwy i przełączenie motywu** | `body` (tło + tekst), `.dn-zakladka`, `.dn-boczna-pozycja`, `.dn-przelacznik`, `.dn-karta-srodowiska::before` |
| `--dn-czas-3` | **0,22 s** | **wejście i wyjście warstwy** | `.dn-modal[open]`, `.dn-toast`, `.dn-tooltip-tresc`, `.dn-postep-wartosc` |
| `--dn-czas-tetno` | **2,4 s** | **jedyny ruch ciągły** | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after` |

### 9.3. Tabela decyzyjna „który czas"

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

### 9.4. Dwa czasy poza skalą — uzasadnienie

| Animacja | Czas | Dlaczego poza skalą |
|---|---|---|
| `dn-obrot` (spinner) | **0,8 s**, `linear` | ruch cykliczny bez punktu docelowego; krzywa `--dn-ease` wprowadziłaby pulsowanie prędkości, czytane jako zacinanie. Jedyne dopuszczone `linear` w systemie |
| `dn-tetno` | **2,4 s**, `--dn-ease` | ruch ciągły niosący stan „praca biegnie"; 2,4 s to rytm wolniejszy od tętna spoczynkowego człowieka — nie wywołuje pobudzenia, a pozostaje zauważalny peryferyjnie |

---

## 10. Katalog dozwolonych animacji

Siedem pozycji. **Katalog jest zamknięty** — animacja spoza niego wymaga decyzji Właściciela.

| # | Animacja | Nośnik | Czas | Własność | Informacja niesiona |
|---|---|---|---|---|---|
| **M1** | **Mikroreakcja kontrolki** | `.dn-btn`, `.dn-btn-ikona`, `.dn-karta--klikalna`, wiersz `.dn-tabela`, `.dn-zakladka` | `--dn-czas-1` | `background-color`, `transform: translateY(1px)` | „element jest interaktywny i przyjął moje działanie" |
| **M2** | **Przejście barwy** | `body`, `.dn-boczna-pozycja`, `.dn-przelacznik`, `.dn-pole-kontrolka:focus` | `--dn-czas-2` | `color`, `background-color`, `border-color`, `box-shadow` | „zmienił się stan wyboru / motyw / fokus" |
| **M3** | **Wejście i wyjście warstwy** | `.dn-modal[open]`, `.dn-toast`, `.dn-tooltip-tresc` | `--dn-czas-3` | `opacity` 0→1, `transform: translateY(8px)`→0, modal dodatkowo `scale(0.98)`→1 | „pojawiła się nowa warstwa nad treścią" |
| **M4** | **Tętno kropki** | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after` | `--dn-czas-tetno` | `box-shadow` (pierścień rozchodzący się) | **„tu biegnie praca"** — jedyny ruch ciągły |
| **M5** | **Postęp kolejki** | `.dn-postep-wartosc` | `--dn-czas-3` | `width` | „proces posunął się o zmierzoną wartość" |
| **M6** | **Wskaźnik pracy wpisu** | `.dn-btn[aria-busy='true']::after`, `.dn-spinner` | 0,8 s `linear` | `transform: rotate()` | „operacja trwa, czas nieznany" |
| **M7** | **Przełączenie motywu** | `body` + wszystkie powierzchnie dziedziczące | `--dn-czas-2` | `background-color`, `color` | „zmieniłem motyw — to ta sama treść" |

### 10.1. M4 — tętno w szczegółach

Tętno jest **elementem sygnaturowym produktu**, nie ozdobą. Występuje dokładnie w pięciu miejscach kontraktu marki:

| Miejsce | Forma | Uwaga |
|---|---|---|
| godło | kropka zamykająca podwójny grot `»».` | **statyczna** — znak się nie animuje |
| emblematy czterech środowisk | dokładnie jedna wypełniona kropka na emblemat | **statyczna** |
| karty sesji | pulsująca kropka = proces w tle | animowana |
| okno komunikacji | kropka przy nadawcy aktywnie piszącym | animowana (`.dn-wpis--pracuje`) |
| Always On Display | rdzeń awatara (pierścień wokół rdzenia) | animowana |

**Reguła jednego tętna.** W jednym widoku pulsuje **jedno** miejsce. Jeżeli pracują trzy sesje jednocześnie — pulsuje karta sesji aktywnej, pozostałe niosą kropkę statyczną plus liczbę w plakietce.

### 10.2. M5 — kiedy postęp, kiedy spinner

| Sytuacja | Komponent | Dlaczego |
|---|---|---|
| znany procent lub znana liczba kroków | `.dn-postep` + `.dn-postep-etykieta` (mono, `tabular-nums`) | pasek bez liczby to ozdoba; liczba bez paska to surowa dana — razem tworzą pomiar |
| kolejka kroków o znanej liście | `.dn-kolejka` + `.dn-krok--pracuje` | postęp jest listą, nie paskiem |
| czas i zakres nieznane | `.dn-spinner` albo `.dn-btn[aria-busy='true']` | pasek udający postęp przy nieznanym zakresie to kłamstwo interfejsu |

---

## 11. Katalog zakazanych animacji

| # | Zakazane | Dlaczego — uzasadnienie z kontraktu |
|---|---|---|
| **X1** | **Parallax** (warstwy przesuwające się z różną prędkością przy przewijaniu) | ruch nie niesie żadnej informacji o stanie systemu; w kokpicie z gęstością 8/10 rozbija odczyt danych; obciąża wątek kompozycji przy każdej klatce przewijania |
| **X2** | **Scrollytelling** (treść odsłaniana i animowana wraz z przewijaniem) | wzorzec strony narracyjnej; okno robocze nie opowiada historii — pokazuje stan. Operator przewija, żeby **czytać dane**, nie żeby uruchamiać przedstawienie |
| **X3** | **Animacje dekoracyjne** (unoszące się kształty, animowane gradienty, cząstki, „oddychające" tła) | anty-domyślne kierunku: „zero dekoracji bez funkcji"; gradient jest wyłącznie ilustracyjny i **nigdy nie jest tłem sekcji** |
| **X4** | **Bounce / przeregulowanie** (`cubic-bezier` przekraczający 1, sprężyny, `elastic`) | system ma **jedną** krzywą `--dn-ease`; przeregulowanie sugeruje fizyczność obiektu, a kontrolki kokpitu nie są obiektami fizycznymi |
| **X5** | **Ruch bez informacji** (animacja przy wejściu na stronę, kaskadowe pojawianie się list, animowane liczniki odliczające do wartości) | każda animacja odpowiada na pytanie „co się zmieniło?"; animowany licznik **fałszuje pomiar** — pokazuje wartości, których system nigdy nie zmierzył |
| **X6** | **Więcej niż jeden ruch ciągły na widok** | tętno traci znaczenie „tu biegnie praca", gdy pulsuje wszystko |
| **X7** | **Animacja wysokości i szerokości układu** (`height`, `top`, `left`, `margin` na elementach układu) | wymusza przeliczenie układu w każdej klatce (rozdz. 14); jedyne dopuszczone `width` to `.dn-postep-wartosc` wewnątrz toru o stałych wymiarach |
| **X8** | **Migotanie** (częstotliwość powyżej 3 Hz) | kryterium WCAG 2.3.1; ryzyko napadu światłoczułego. Tętno 2,4 s = **0,42 Hz** — siedmiokrotnie poniżej progu |
| **X9** | **Rozmycie jako efekt** (`backdrop-filter` poza nakładką modala) | anty-domyślne: „glassmorfizm wszędzie" → powierzchnie kryjące; jedyne dopuszczone rozmycie to `blur(2px)` na `.dn-modal::backdrop` |
| **X10** | **Animowane przewijanie sterowane skryptem** (`scroll-behavior: smooth` wymuszone globalnie) | odbiera Operatorowi kontrolę nad tempem czytania; przy `prefers-reduced-motion` system wymusza `scroll-behavior: auto` |

---

## 12. Wzorce implementacyjne

Kod poniżej jest **cytatem z biblioteki** `zasoby/css/komponenty.css` i `zasoby/zetony/zetony.css` — nie propozycją.

### 12.1. Mikroreakcja kontrolki (M1)

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

**Reguła:** wylicz własności jawnie. `transition: all` animuje także własności układu i łamie X7.

### 12.2. Tętno (M4) — klatki kluczowe

```css
.dn-kropka--tetno { animation: dn-tetno var(--dn-czas-tetno) var(--dn-ease) infinite; }

@keyframes dn-tetno {
  0%, 100% { box-shadow: 0 0 0 0  var(--dn-fokus-cien); }
  40%      { box-shadow: 0 0 0 5px transparent; }
}
```

Animowany jest **`box-shadow`**, nie `width`/`height` — pierścień rozchodzi się bez wpływu na układ sąsiadów.

### 12.3. Wejście modala (M3)

```css
.dn-modal[open] { animation: dn-wejscie var(--dn-czas-3) var(--dn-ease); }

@keyframes dn-wejscie {
  from { opacity: 0; transform: translateY(var(--dn-od-2)) scale(0.98); }
}
```

Klatka `to` jest pominięta świadomie — stan docelowy to stan spoczynkowy elementu. Dystans 8 px = `--dn-od-2`.

Nakładka:

```css
.dn-modal::backdrop { background: var(--dn-nakladka); backdrop-filter: blur(2px); }
```

### 12.4. Powiadomienie (M3)

```css
.dn-toast { animation: dn-wejscie var(--dn-czas-3) var(--dn-ease); }
```

Toast **dzieli klatki kluczowe z modalem**. Jedno wejście warstwy = jedna definicja ruchu.

### 12.5. Postęp (M5)

```css
.dn-postep-tor      { height: 4px; background: var(--dn-powierzchnia-2); overflow: hidden; }
.dn-postep-wartosc  {
  height: 100%;
  background: var(--dn-sygnal-wypelnienie);
  transition: width var(--dn-czas-3) var(--dn-ease);
}
```

Wartość zmienia deklarację `width` (np. z `34%` na `61%`); przejście robi CSS. Etykieta obok podaje liczbę krojem mono.

### 12.6. Wskaźnik pracy (M6)

```css
.dn-spinner { animation: dn-obrot 0.8s linear infinite; }
@keyframes dn-obrot { to { transform: rotate(360deg); } }
```

### 12.7. Przełączenie motywu (M7)

```css
body {
  transition:
    background-color var(--dn-czas-2) var(--dn-ease),
    color            var(--dn-czas-2) var(--dn-ease);
}
```

Motyw przełącza się zmianą atrybutu `data-theme` na `<html>`; kaskada robi resztę. **Zakaz animowania każdego elementu z osobna** — to setki jednoczesnych przejść.

---

## 13. `prefers-reduced-motion` — co dokładnie się dzieje

### 13.1. Obsługa globalna, nie per komponent

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

**Dwa mechanizmy działają razem:**

| Mechanizm | Co obejmuje | Dlaczego potrzebny |
|---|---|---|
| podmiana żetonów czasu | wszystko, co używa `--dn-czas-*` | źródło prawdy — komponent nie musi wiedzieć o preferencji |
| reguła `*` z `!important` | animacje i przejścia z czasami spoza skali (np. `0.8s` spinnera) oraz kod zewnętrzny | siatka bezpieczeństwa |

`animation-iteration-count: 1` zatrzymuje pętle po pierwszym przebiegu — animacja nie „ucina się" w losowej klatce.

### 13.2. Tętno → pierścień statyczny

Skrócenie czasu do 0,01 ms wygasza ruch, ale usuwa też **informację**: kropka pracująca stałaby się nieodróżnialna od zwykłej. Dlatego biblioteka podstawia **zamiennik statyczny**:

```css
@media (prefers-reduced-motion: reduce) {
  .dn-kropka--tetno { box-shadow: 0 0 0 2px var(--dn-fokus-cien); }
}
```

```
   RUCH DOZWOLONY                    RUCH OGRANICZONY
   ┌────────────────┐                ┌────────────────┐
   │   ((( ● )))    │   2,4 s        │     ( ● )      │   pierścień
   │  pierścień     │   pętla        │   statyczny    │   2 px, bez ruchu
   │  rozchodzi się │                │   zawsze       │
   └────────────────┘                └────────────────┘
   informacja: praca biegnie         informacja: praca biegnie
   nośnik: ruch                      nośnik: kształt
```

**Zasada ogólna:** ograniczenie ruchu **nie może usuwać informacji**. Jeżeli ruch był jedynym jej nośnikiem, przy `reduce` informację przejmuje kształt, obrys albo etykieta.

### 13.3. Tabela zachowań przy `reduce`

| Animacja | Zachowanie przy `reduce` | Nośnik informacji po zmianie |
|---|---|---|
| M1 mikroreakcja | zmiana natychmiastowa | barwa tła (`--dn-hover`) |
| M2 przejście barwy | zmiana natychmiastowa | barwa docelowa |
| M3 wejście warstwy | warstwa pojawia się bez wjazdu | obecność warstwy + cień + nakładka |
| **M4 tętno** | **pierścień statyczny 2 px** | **kształt (pierścień)** |
| M5 postęp | pasek skacze do wartości | liczba w etykiecie mono |
| M6 spinner | zatrzymany po pierwszym obrocie | `aria-busy="true"` + tekst „Trwa…" |
| M7 motyw | przełączenie natychmiastowe | nowe barwy |

### 13.4. Czego NIE robić

| Zakaz | Uzasadnienie |
|---|---|
| dublowanie zapytania `prefers-reduced-motion` w każdym komponencie | obsługa jest globalna; duplikat rozjeżdża się przy pierwszej zmianie |
| wykrywanie preferencji w JavaScript i wyłączanie klas | preferencja może zmienić się w trakcie sesji; CSS reaguje sam |
| całkowite usunięcie zamiennika przy `reduce` | usuwa informację o stanie systemu (rozdz. 13.2) |

---

## 14. Ruch a wydajność

### 14.1. Dwie własności bezpieczne

| Własność | Etap potoku renderowania | Koszt |
|---|---|---|
| `transform` | wyłącznie **kompozycja** | najniższy — obsługiwana poza wątkiem głównym |
| `opacity` | wyłącznie **kompozycja** | najniższy |
| `box-shadow` | **malowanie** | średni — dopuszczony dla tętna (jeden mały element) i fokusu |
| `color`, `background-color` | **malowanie** | średni — dopuszczony dla mikroreakcji i przejść barwy |
| `width`, `height`, `top`, `left`, `margin`, `padding` | **układ** (przeliczenie geometrii) | najwyższy — **zakazane** poza `.dn-postep-wartosc` |

### 14.2. Dlaczego `.dn-postep-wartosc` może animować `width`

Element leży wewnątrz `.dn-postep-tor` z `overflow: hidden` i stałą wysokością 4 px. Przeliczenie układu ogranicza się do jednego niezależnego elementu i nie propaguje na rodzeństwo. Zamiana na `transform: scaleX()` rozciągałaby zaokrąglenia narożników — pogorszenie wizualne bez zysku pomiarowego przy jednym pasku na widok.

### 14.3. `will-change` — zasada trzech warunków

Stosować **wyłącznie** gdy spełnione są trzy warunki równocześnie:

1. animacja jest długa lub cykliczna (tętno, spinner),
2. element ma znany, ograniczony rozmiar,
3. deklaracja jest **zdejmowana** po zakończeniu animacji.

| Zakaz | Skutek |
|---|---|
| `will-change` na wielu elementach jednocześnie | każdy dostaje własną warstwę kompozycji — zużycie pamięci graficznej rośnie liniowo |
| `will-change` na stałe w arkuszu | warstwa nigdy nie znika; przeglądarka traci możliwość optymalizacji |
| `will-change: all` | przeciwieństwo optymalizacji |

**Stan biblioteki:** `komponenty.css` **nie deklaruje `will-change`** w żadnym miejscu. Animowane elementy są małe (kropka 6 px, spinner 14 px, pierścień AOD 36 px) i przeglądarka promuje je do warstwy samodzielnie. Wprowadzenie `will-change` wymaga pomiaru pokazującego zysk.

### 14.4. Unikanie przeliczeń układu (layout thrash)

| Wzorzec błędny | Wzorzec poprawny |
|---|---|
| w pętli: odczyt `offsetHeight` → zapis `style.height` → odczyt → zapis | wszystkie odczyty **przed** wszystkimi zapisami |
| animowanie `height` przy rozwijaniu sekcji | `<details>` natywny albo `opacity` + `transform` na treści o znanej wysokości |
| animowanie `top`/`left` przy przesuwaniu warstwy | `transform: translateY()` |
| pomiar szerokości paska postępu w skrypcie przy każdej klatce | zmiana deklaracji `width` raz; przejście wykonuje CSS |

### 14.5. Budżet ruchu na widok

| Pozycja | Limit |
|---|---|
| animacje ciągłe (`infinite`) | **1** — tętno; spinner jest ciągły wyłącznie na czas trwania operacji |
| jednoczesne przejścia wywołane jednym działaniem | ≤ 3 (np. wjazd modala = `opacity` + `transform` + nakładka) |
| elementy animowane przy przełączeniu motywu | `body` + powierzchnie dziedziczące — **nigdy** przejścia deklarowane per komponent |

---

# CZĘŚĆ III — DOSTĘPNOŚĆ

## 15. WCAG 2.1 AA jako warunek wejściowy

> **WCAG 2.1 AA to warunek wejściowy, nie cel.** Okno, które go nie spełnia, nie jest gotowe — nie jest „gotowe z zastrzeżeniem".

### 15.1. Kryteria, których system dotyka

| Kryterium | Poziom | Gdzie system je realizuje |
|---|---|---|
| **1.1.1** Treść nietekstowa | A | ikony dekoracyjne `aria-hidden="true"`; ikony informacyjne z `aria-label` albo tekstem obok |
| **1.3.1** Informacje i relacje | A | `<table>` z `<th scope>`, `<fieldset>`/`<legend>`, `role="tablist"`, hierarchia nagłówków |
| **1.3.2** Zrozumiała kolejność | A | kolejność DOM = kolejność odczytu; brak przestawiania układem |
| **1.4.1** Użycie koloru | A | **stan nigdy samym kolorem** (rozdz. 19) |
| **1.4.3** Kontrast (minimalny) | AA | 33 pomiary (rozdz. 16), reguła `--dn-tekst-3` |
| **1.4.4** Zmiana rozmiaru tekstu | AA | stopnie w px na skali; układ oparty o `grid`/`flex`, bez wysokości stałych na kontenerach treści |
| **1.4.10** Zawijanie treści (reflow) | AA | punkty łamania w1–w4; boczna nawigacja zwija się do ikon przy 960 px |
| **1.4.11** Kontrast elementów nietekstowych | AA | obrys kontrolki, pierścień fokusu, kropki stanu — pomiary 13, 14, 29 |
| **1.4.12** Odstępy w tekście | AA | interlinia 1,45 bazowa / 1,6 luźna; brak `height` stałych na blokach tekstu |
| **1.4.13** Treść przy najechaniu lub fokusie | AA | `.dn-tooltip` reaguje na `:hover` **i** `:focus-within`; znika po utracie obu |
| **2.1.1** Klawiatura | A | pełna nawigacja klawiaturą we wszystkich prototypach |
| **2.1.2** Brak pułapki klawiatury | A | pułapka wyłącznie w `<dialog>` (celowa, z wyjściem `Esc`) |
| **2.2.1** Dostosowanie czasu | A | odliczanie „Wyślij ponownie" **informacyjne**, nie blokuje |
| **2.2.2** Pauza, zatrzymanie, ukrycie | A | tętno = wskaźnik stanu; `prefers-reduced-motion` zastępuje je pierścieniem statycznym (rozdz. 24, D-08) |
| **2.3.1** Trzy błyski lub poniżej progu | A | tętno 0,42 Hz; brak migotania w systemie |
| **2.4.1** Możliwość pominięcia bloków | A | odnośnik „Przejdź do treści" (`.dn-sr-only`) jako pierwszy element `<body>` |
| **2.4.3** Kolejność fokusu | A | rozdz. 17.2 |
| **2.4.6** Nagłówki i etykiety | AA | każde pole ma `.dn-pole-etykieta` powiązaną atrybutem `for` |
| **2.4.7** Widoczny fokus | AA | pierścień 2 px + odsunięcie 2 px, `--dn-fokus` (rozdz. 17) |
| **2.5.1** Gesty wskazujące | A | brak gestów wielopunktowych i ścieżkowych jako jedynej drogi |
| **2.5.3** Etykieta w nazwie | A | `aria-label` przycisku ikonowego zawiera widoczną nazwę czynności |
| **3.1.1** Język strony | A | `lang="pl"` na `<html>` |
| **3.2.1** Po otrzymaniu fokusu | A | fokus nie wywołuje zmiany kontekstu |
| **3.2.2** Podczas wprowadzania danych | A | wpisanie wartości nie przenosi widoku ani nie zatwierdza formularza |
| **3.3.1** Identyfikacja błędu | A | `.dn-pole-blad` + `aria-invalid` + ikona |
| **3.3.2** Etykiety lub instrukcje | A | `.dn-pole-opis` przy polach wymagających wyjaśnienia |
| **3.3.3** Sugestia korekty błędu | AA | komunikat wskazuje brak, nie tylko fakt błędu |
| **4.1.2** Nazwa, rola, wartość | AA | `aria-pressed`, `aria-selected`, `aria-current`, `aria-busy`, `aria-expanded` |
| **4.1.3** Komunikaty o stanie | AA | `aria-live` dla powiadomień i zmian stanu (rozdz. 22) |

### 15.2. Kryteria poziomu AAA realizowane dobrowolnie

| Kryterium | Poziom | Realizacja |
|---|---|---|
| **1.4.6** Kontrast wzmocniony | AAA | 11 z 33 pomiarów osiąga AAA (≥ 7:1) — cała para tekst/tło w obu motywach |
| **2.5.5** Rozmiar celu | AAA | `pointer: coarse` podnosi kontrolkę do 40 px, wiersz do 44 px |
| **2.3.3** Animacja z interakcji | AAA | `prefers-reduced-motion` obsłużone globalnie |

---

## 16. Kontrast: metoda, komplet 33 pomiarów, zasada tekst-3

### 16.1. Metoda pomiaru

Kontrast liczony **algorytmem WCAG 2.1** — iloraz luminancji względnej sRGB:

```
   1. Kanał znormalizowany:      c = C / 255                    (C ∈ ⟨0,255⟩)

   2. Linearyzacja kanału:       c ≤ 0,03928  →  c / 12,92
                                 c >  0,03928  →  ((c + 0,055) / 1,055) ^ 2,4

   3. Luminancja względna:       L = 0,2126·R + 0,7152·G + 0,0722·B

   4. Kontrast pary:             K = (L_jaśniejsza + 0,05) / (L_ciemniejsza + 0,05)

   5. Barwa z kanałem alfa:      przed pomiarem złożyć z tłem
                                 wynik = barwa·α + tło·(1 − α)
```

### 16.2. Progi

| Rodzaj treści | Próg | Uzasadnienie |
|---|---|---|
| tekst normalny | **≥ 4,5 : 1** | WCAG 2.1, kryterium 1.4.3 |
| tekst duży — ≥ 18,66 px przy wadze 700 albo ≥ 24 px | **≥ 3 : 1** | WCAG 2.1, kryterium 1.4.3 |
| elementy graficzne i kontrolki (obrys, ikona niosąca znaczenie, pierścień fokusu) | **≥ 3 : 1** | WCAG 2.1, kryterium 1.4.11 |
| **próg systemowy** — obrys dekoracyjny rozdzielający powierzchnie o tej samej roli | **≥ 1,6 : 1** | **próg własny systemu**, poza WCAG — obrys nie niesie informacji, oddziela powierzchnie |

### 16.3. Komplet 33 pomiarów

Źródło: `zasoby/zetony/kontrasty.json`. Wartości cytowane, nie przeliczane ponownie.

| # | Para | Motyw | Barwa treści | Barwa podłoża | Kontrast | Próg | Werdykt |
|---|---|---|---|---|---|---|---|
| 01 | tekst / tło | jasny | `--dn-szary-900` #181818 | `--dn-szary-50` #F4F4F4 | **16,14** | 4,5 | AAA |
| 02 | tekst / powierzchnia | jasny | `--dn-szary-900` #181818 | `--dn-szary-0` #FFFFFF | **17,76** | 4,5 | AAA |
| 03 | tekst-2 / tło | jasny | `--dn-szary-600` #616161 | `--dn-szary-50` #F4F4F4 | **5,63** | 4,5 | AA |
| 04 | tekst-2 / powierzchnia | jasny | `--dn-szary-600` #616161 | `--dn-szary-0` #FFFFFF | **6,19** | 4,5 | AA |
| 05 | tekst-2 / powierzchnia-2 | jasny | `--dn-szary-600` #616161 | `--dn-szary-100` #ECECEC | **5,24** | 4,5 | AA |
| 06 | tekst-3 / tło | jasny | `--dn-szary-500` #7C7C7C | `--dn-szary-50` #F4F4F4 | **3,80** | 3,0 | AA (tylko metadane) |
| 07 | tekst-3 / powierzchnia | jasny | `--dn-szary-500` #7C7C7C | `--dn-szary-0` #FFFFFF | **4,17** | 3,0 | AA (tylko metadane) |
| 08 | biel / przycisk atrament | jasny | `--dn-szary-0` #FFFFFF | `--dn-szary-900` #181818 | **17,76** | 4,5 | AAA |
| 09 | sygnał-tekst / tło | jasny | `--dn-sygnal-600` #2457C9 | `--dn-szary-50` #F4F4F4 | **5,82** | 4,5 | AA |
| 10 | sygnał-tekst / powierzchnia | jasny | `--dn-sygnal-600` #2457C9 | `--dn-szary-0` #FFFFFF | **6,40** | 4,5 | AA |
| 11 | sygnał-tekst / sygnał-tło | jasny | `--dn-sygnal-600` #2457C9 | `--dn-sygnal-100` #EDF3FE | **5,74** | 4,5 | AA |
| 12 | biel / sygnał wypełnienie | jasny | `--dn-szary-0` #FFFFFF | `--dn-sygnal-600` #2457C9 | **6,40** | 4,5 | AA |
| 13 | obrys kontrolki / tło | jasny | `--dn-szary-300` #C0C0C0 | `--dn-szary-50` #F4F4F4 | **1,65** | 1,6 | próg systemowy |
| 14 | fokus sygnał / tło | jasny | `--dn-sygnal-500` #3B6FE0 | `--dn-szary-50` #F4F4F4 | **4,21** | 3,0 | AA |
| 15 | sukces tekst / tło stanu | jasny | `--dn-zielen-700` #1E7A46 | `--dn-zielen-100` #E7F5EE | **4,76** | 4,5 | AA |
| 16 | ostrzeżenie tekst / tło stanu | jasny | `--dn-bursztyn-700` #8A5B0C | `--dn-bursztyn-100` #FBF2DE | **5,26** | 4,5 | AA |
| 17 | błąd tekst / tło stanu | jasny | `--dn-czerwien-700` #C0362F | `--dn-czerwien-100` #FCEDEB | **4,84** | 4,5 | AA |
| 18 | sukces tekst / powierzchnia | jasny | `--dn-zielen-700` #1E7A46 | `--dn-szary-0` #FFFFFF | **5,35** | 4,5 | AA |
| 19 | ostrzeżenie tekst / powierzchnia | jasny | `--dn-bursztyn-700` #8A5B0C | `--dn-szary-0` #FFFFFF | **5,86** | 4,5 | AA |
| 20 | błąd tekst / powierzchnia | jasny | `--dn-czerwien-700` #C0362F | `--dn-szary-0` #FFFFFF | **5,51** | 4,5 | AA |
| 21 | tekst / tło | ciemny | `--dn-szary-100` #ECECEC | `--dn-szary-950` #0F0F0F | **16,23** | 4,5 | AAA |
| 22 | tekst / powierzchnia | ciemny | `--dn-szary-100` #ECECEC | `--dn-szary-900` #181818 | **15,03** | 4,5 | AAA |
| 23 | tekst-2 / tło | ciemny | `--dn-szary-400` #9E9E9E | `--dn-szary-950` #0F0F0F | **7,15** | 4,5 | AAA |
| 24 | tekst-2 / powierzchnia | ciemny | `--dn-szary-400` #9E9E9E | `--dn-szary-900` #181818 | **6,63** | 4,5 | AA |
| 25 | tekst-3 / powierzchnia | ciemny | `--dn-szary-500` #7C7C7C | `--dn-szary-900` #181818 | **4,25** | 3,0 | AA (tylko metadane) |
| 26 | atrament / przycisk biel | ciemny | `--dn-szary-950` #0F0F0F | `--dn-szary-100` #ECECEC | **16,23** | 4,5 | AAA |
| 27 | sygnał-tekst / tło | ciemny | `--dn-sygnal-300` #8FB2F5 | `--dn-szary-950` #0F0F0F | **9,00** | 4,5 | AAA |
| 28 | sygnał-tekst / powierzchnia | ciemny | `--dn-sygnal-300` #8FB2F5 | `--dn-szary-900` #181818 | **8,33** | 4,5 | AAA |
| 29 | fokus sygnał / tło | ciemny | `--dn-sygnal-400` #5C8CEC | `--dn-szary-950` #0F0F0F | **5,86** | 3,0 | AA |
| 30 | sukces tekst / powierzchnia | ciemny | `--dn-zielen-400` #4FBD85 | `--dn-szary-900` #181818 | **7,57** | 4,5 | AAA |
| 31 | ostrzeżenie tekst / powierzchnia | ciemny | `--dn-bursztyn-400` #E0A63C | `--dn-szary-900` #181818 | **8,19** | 4,5 | AAA |
| 32 | błąd tekst / powierzchnia | ciemny | `--dn-czerwien-400` #EE7168 | `--dn-szary-900` #181818 | **6,09** | 4,5 | AA |
| 33 | biel / sygnał wypełnienie | ciemny | `--dn-szary-0` #FFFFFF | `--dn-sygnal-500` #3B6FE0 | **4,63** | 3,0 | AA |

### 16.4. Zestawienie wyniku

| Miara | Wartość |
|---|---|
| Liczba pomiarów | **33** |
| Pary motywu jasnego | 20 |
| Pary motywu ciemnego | 13 |
| Werdykt **AAA** (≥ 7,0) | **11** |
| Werdykt **AA** | **21** |
| Próg systemowy (obrys) | **1** |
| Pozycji „NIE PRZECHODZI" | **0** |
| Najniższy pomiar | 1,65 (obrys kontrolki — próg systemowy 1,6) |
| Najwyższy pomiar | 17,76 (tekst / powierzchnia, jasny) |
| Najniższy pomiar tekstowy | 3,80 (`--dn-tekst-3` na tle jasnym — patrz 16.5) |

### 16.5. Zasada „tekst-3 tylko metadane" [NIENEGOCJOWALNE]

`--dn-tekst-3` (`#7C7C7C` w obu motywach) osiąga **3,80 : 1** na tle jasnym i **4,25 : 1** na powierzchni ciemnej — **poniżej progu 4,5 dla tekstu normalnego**. Wolno go użyć wyłącznie do:

| Dozwolone | Przykład z produktu |
|---|---|
| **tekst duży** — ≥ 18,66 px przy wadze 700 albo ≥ 24 px | podtytuł pustego stanu |
| **metadane** — znaczniki czasu, liczniki, identyfikatory pomocnicze | godzina wpisu (`.dn-wpis-godzina`), opis kafla (`.dn-kafel-opis`) |
| **etykiety wersalikowe** — 11 px, wersaliki, tracking | `.dn-etykieta-wersalikowa`, `.dn-etykieta-mono`, nagłówki kolumn `.dn-tabela th` |
| **tekst zastępczy pola** (`::placeholder`) | podpowiedź w `.dn-prompt-obszar` |
| **stan pustego widoku** | `.dn-pusty-stan` |

| Zakazane | Zamiast tego |
|---|---|
| treść czytana (akapit, opis funkcji) | `--dn-tekst-2` |
| etykieta pola formularza | `--dn-tekst` |
| wartość danej w tabeli | `--dn-tekst-2` (klasa `.dn-dane`) |
| komunikat błędu | `--dn-blad-tekst` |

**Reguła weryfikacji:** jeżeli treść w `--dn-tekst-3` niesie informację, bez której Operator nie wykona zadania — żeton jest użyty błędnie.

### 16.6. Pary bez pomiaru w komplecie

`kontrasty.json` mierzy **pary komponentowe**, nie iloczyn kartezjański palety. Pary poniżej wymagają pomiaru **przed** wprowadzeniem nowego zestawienia:

| Para wymagająca pomiaru przy nowym zestawieniu | Uwaga |
|---|---|
| tekst na `--dn-panel` w motywie jasnym (`#FAFAFA`) | pomiar 02 (na `#FFFFFF`) jest niższym oszacowaniem — panel jest ciemniejszy o jeden stopień |
| stany na `--dn-panel` | tło modala i toastu to `--dn-panel`, nie `--dn-powierzchnia` |
| tekst na `--dn-rama` (pasek górny) | rama jest stała w obu motywach; `--dn-rama-tekst` = `szary-100` na `szary-925` |
| barwy półprzezroczyste motywu ciemnego (`--dn-sygnal-tlo`, tła stanów) | **złożyć z podłożem przed pomiarem** (16.1, krok 5) |

---

## 17. Fokus

### 17.1. Pierścień

| Cecha | Wartość | Żeton |
|---|---|---|
| Grubość | **2 px** | `--dn-wym-fokus` |
| Odsunięcie | **2 px** | `--dn-wym-fokus-odsuniecie` |
| Barwa | sygnał — `sygnal-500` (jasny) / `sygnal-400` (ciemny) | `--dn-fokus` |
| Poświata | `rgba(59,111,224,.30)` / `rgba(92,140,236,.35)` | `--dn-fokus-cien` |
| Wyzwalanie | **wyłącznie nawigacja klawiaturą** | `:focus-visible` |
| Zakres | **każda kontrolka interaktywna, bez wyjątku** | fundament globalny |

```css
:focus-visible {
  outline: var(--dn-wym-fokus) solid var(--dn-fokus);
  outline-offset: var(--dn-wym-fokus-odsuniecie);
  border-radius: var(--dn-r-xs);
}
```

**Pomiar:** fokus na tle jasnym **4,21 : 1**, na tle ciemnym **5,86 : 1** — oba powyżej progu 3,0 dla elementów graficznych (1.4.11).

**Zakaz usunięcia.** `outline: none` bez równoważnego wskaźnika jest błędem wdrożenia. Dotyczy też widoków gęstych: macierzy przełączników okna izolacji, przybornika promptu, pasa kart sesji.

**Odsunięcie 2 px, nie 0.** Przy gęstości zwartej (kontrolki 32 px, wiersze 36 px) pierścień przylegający zlewa się z obrysem kontrolki. Odsunięcie tworzy przerwę, która czyta się jako oddzielna warstwa.

### 17.2. Kolejność tabulacji

| Zasada | Realizacja |
|---|---|
| **Kolejność DOM = kolejność fokusu** | zakaz `tabindex` dodatniego; wyłącznie `0` (włącz do kolejności) i `-1` (fokus programowy) |
| **Pominięcie bloków** | odnośnik „Przejdź do treści" (`.dn-sr-only`) jako pierwszy element `<body>`; staje się widoczny przy fokusie |
| **Kolejność w powłoce** | pasek górny → pas kart sesji → boczna nawigacja → obszar roboczy → pas komunikacji |
| **Grupy jednym przystankiem** | grupa opcji jednokrotnych (`radiogroup`), zestaw zakładek (`tablist`), lista kart sesji — `Tab` wchodzi w grupę, strzałki przemieszczają wewnątrz |

```
   KOLEJNOŚĆ TABULACJI — POWŁOKA ŚRODOWISKA (przykładowa)

   [Przejdź do treści]                     ← 0. pierwszy przystanek, .dn-sr-only
        │
   ┌────▼──────────────────────────────────────────────────────────┐
   │ 1 godło · 2 wyszukiwarka · 3 powiadomienia · 4 motyw · 5 konto │  pasek górny
   ├───────────────────────────────────────────────────────────────┤
   │ 6 karty sesji (Tab wchodzi → ←/→ przemieszcza) · 7 «+»         │  pas kart
   ├──────────────┬────────────────────────────────────────────────┤
   │ 8 boczna     │ 10 obszar roboczy okna operacyjnego            │
   │   nawigacja  │                                                │
   │   (↑/↓)      ├────────────────────────────────────────────────┤
   │ 9 stopka     │ 11 przybornik promptu · 12 pole · 13 wyślij     │  pas komunikacji
   └──────────────┴────────────────────────────────────────────────┘
```

### 17.3. Pułapka fokusu w modalu

**Jedyne dopuszczone miejsce pułapki.** Realizacja natywnym `<dialog>` + `showModal()`:

| Zachowanie | Realizacja |
|---|---|
| fokus wchodzi do modala przy otwarciu | `showModal()` ustawia fokus na pierwszym elemencie fokusowalnym |
| `Tab` krąży wewnątrz modala | natywna właściwość `<dialog>` otwartego modalnie |
| `Esc` zamyka | natywne zdarzenie `cancel` |
| kliknięcie w nakładkę zamyka | obsługa własna — porównanie celu zdarzenia z elementem `<dialog>` |
| fokus wraca do wyzwalacza po zamknięciu | natywne dla `<dialog>`; przy zamknięciu programowym — przywrócić jawnie |
| treść pod modalem niedostępna | `::backdrop` + `inert` na tle (natywne dla `showModal()`) |

**To nie jest blokada w rozumieniu ADL-017.** Nakładka ogranicza interakcję wyłącznie na czas otwarcia i zawsze udostępnia drogę zamknięcia (`Esc`, kontrolka zamknięcia, kliknięcie w nakładkę, akcja stopki).

### 17.4. Fokus a stany komponentów

| Komponent | Cecha szczególna fokusu |
|---|---|
| `.dn-btn` | pierścień + zachowany `transform` stanu `:active` |
| `.dn-btn-ikona--na-ramie` | pierścień w barwie sygnału na atramentowej ramie — kontrast 5,86 (pomiar 29) |
| `.dn-pole-kontrolka` | `:focus` zmienia obrys na `--dn-fokus` + poświata `--dn-cien-sygnal`; **`:focus`, nie `:focus-visible`** — pole tekstowe pokazuje aktywność także po kliknięciu |
| `.dn-prompt` | `:focus-within` na kontenerze — pierścień obejmuje grot, obszar i przybornik jako jedną kontrolkę |
| `.dn-tooltip` | `:focus-within` pokazuje treść dymka — dymek dostępny z klawiatury (1.4.13) |
| `.dn-karta-sesji` | pierścień + wstęga aktywności 2 px nie znikają jednocześnie |

---

## 18. Klawiatura

### 18.1. Zachowania per komponent

| Komponent | `Tab` | Strzałki | `Enter` | `Spacja` | `Esc` |
|---|---|---|---|---|---|
| `.dn-btn`, `.dn-btn-ikona` | przystanek | — | uruchamia | uruchamia | — |
| `.dn-pole-kontrolka` (jednowierszowe) | przystanek | kursor w tekście | zatwierdza formularz | znak spacji | czyści pole wyszukiwania |
| `.dn-prompt-obszar` | przystanek | kursor w tekście | nowa linia; **`Ctrl/Cmd + Enter` wysyła** | znak spacji | — |
| `.dn-check` | przystanek | — | — | przełącza | — |
| `.dn-radio` (grupa) | jeden przystanek na grupę | **↑ ↓ ← →** przemieszczają i wybierają | — | wybiera | — |
| `.dn-przelacznik` | przystanek | — | przełącza | przełącza | — |
| `.dn-suwak` | przystanek | ← → o krok; ↑ ↓ o krok | — | — | — |
| `.dn-wybor` (lista rozwijana) | przystanek | ↑ ↓ zmienia wartość | otwiera / zatwierdza | otwiera | zamyka bez zmiany |
| `.dn-zakladki` / `.dn-zakladka` | jeden przystanek na zestaw | ← → przemieszcza między zakładkami | aktywuje | aktywuje | — |
| `.dn-karty-sesji` | jeden przystanek na pas | ← → przemieszcza między kartami | otwiera kartę | otwiera kartę | — |
| `.dn-boczna-pozycja` | przystanek na pozycję | ↑ ↓ przemieszcza w obrębie nawigacji | otwiera moduł | otwiera moduł | — |
| `.dn-tabela` (wiersze wybieralne) | jeden przystanek na tabelę | ↑ ↓ przemieszcza wiersz; ← → kolumnę | otwiera pozycję | zaznacza wiersz | czyści zaznaczenie |
| `.dn-modal` | pułapka wewnątrz | — | akcja główna stopki | — | **zamyka** |
| `.dn-tooltip` | fokus na wyzwalaczu pokazuje dymek | — | — | — | ukrywa dymek |
| `.dn-toast` | przystanek na kontrolce zamknięcia | — | zamyka | zamyka | zamyka najnowszy |
| `.dn-kolejka` / `.dn-krok` | jeden przystanek na kolejkę | ↑ ↓ przemieszcza krok | otwiera szczegóły kroku | — | — |

### 18.2. Skróty udokumentowane w modułach

Cytat z dokumentacji modułów. **Zakres obowiązywania podany w kolumnie „Okno".**

| Skrót | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + Enter` | Wysłanie polecenia do AI | **Chat Window** (wspólne wszystkim modułom) |
| `Ctrl/Cmd + F` | Wyszukiwanie / Grep w bieżącym oknie | Logs Viewer, Errors Panel, Library Explorer, Studio Editor, Terminal Tabs |
| `Ctrl/Cmd + Shift + F` | Grep — znajdź w plikach / w Output Console | Code Editor, Output Console |
| `Ctrl/Cmd + S` | Zapis pliku (uzupełnia autozapis) | Code Editor, Studio Editor |
| `Ctrl/Cmd + P` | Szybkie otwarcie pliku | Code Editor, Project Tree |
| `Ctrl/Cmd + Shift + P` | Paleta poleceń edytora | Code Editor |
| `Ctrl/Cmd + B` · `I` · `U` | Pogrubienie, kursywa, podkreślenie | Studio Editor |
| `Ctrl/Cmd + D` | Zaznacz kolejne wystąpienie (wielokursorowość) | Code Editor |
| `Ctrl/Cmd + /` | Komentarz / odkomentowanie linii | Code Editor |
| `Ctrl/Cmd + T` | Nowa karta powłoki / nowa karta przeglądania | Terminal Tabs · Browser Window |
| `Ctrl/Cmd + W` | Zamknięcie bieżącej karty | Terminal Tabs |
| `Ctrl/Cmd + Tab` | Przełączenie na kolejną kartę | Terminal Tabs |
| `Ctrl/Cmd + Shift + \` | Podział karty pionowo | Terminal Tabs |
| `Ctrl/Cmd + L` | Fokus na pasku adresu | Browser Window |
| `Ctrl/Cmd + D` | Dodanie strony do zakładek | Browser Window |
| `Ctrl/Cmd + Shift + S` | Zrzut ekranu | Browser Window |
| `Ctrl/Cmd + Shift + A` | Uruchomienie pełnej analizy | Diagnostics Center |
| `Ctrl/Cmd + Shift + M` | Dodanie nowego uczestnika / modelu | Model Panels |
| `Tab` | Przejście między panelami uczestników | Model Panels |
| `Ctrl/Cmd + Enter` | Uruchomienie kolejnej tury | Debate Panel, Moderator Panel |
| `Ctrl/Cmd + K` | Szybkie dodanie źródła | Sources Manager |
| `Ctrl/Cmd + N` | Nowe ustalenie | Findings Panel |
| `Ctrl/Cmd + E` | Otwarcie Export Panel | Report Builder |
| `Ctrl/Cmd + Shift + L` | Dodanie nowego języka docelowego | Translation Panels |
| `Ctrl/Cmd + G` | Otwarcie Glossary Manager | dowolne okno modułu Translate |
| `Ctrl/Cmd + Klik` | Zaznaczenie wielokrotne | Library Explorer |

### 18.3. Reguły projektowania skrótów

| Reguła | Uzasadnienie |
|---|---|
| **Zakaz skrótów jednoklawiszowych** bez modyfikatora | WCAG 2.1.4 — użytkownik sterowania głosem uruchamiałby je przypadkiem |
| Skrót nie może być **jedyną** drogą do funkcji | każda funkcja osiągalna także przez kontrolkę widoczną |
| Kolizja z `Ctrl/Cmd + D` (Code Editor vs Browser Window) rozstrzygana **kontekstem okna** | skrót działa w oknie mającym fokus; okna nie współdzielą przestrzeni skrótów |
| Skrót ujawniany w interfejsie | dymek kontrolki zawiera zapis klawiszy w elemencie `.dn-kbd` |

---

## 19. Stan nigdy samym kolorem

> **Żaden stan — sukces, ostrzeżenie, błąd, informacja, aktywność, wybór, praca w tle — nie jest sygnalizowany wyłącznie barwą.**

### 19.1. Trzy nośniki

```
   BARWA          +      KSZTAŁT / IKONA      +      SŁOWO
   ────────              ───────────────             ─────
   szybka                jednoznaczna                pełna
   rozpoznawalność       przy daltonizmie            przy czytniku ekranu
   peryferyjna           i w druku mono              i przy niskim kontraście

   Każdy stan niesie CO NAJMNIEJ DWA nośniki. Barwa nigdy sama.
```

### 19.2. Realizacja per wzorzec

| Wzorzec | Barwa | Kształt / ikona | Słowo |
|---|---|---|---|
| Plakietka stanu (`.dn-plakietka--sukces` …) | tło + obrys + tekst w rodzinie stanu | ikona 12 px w plakietce | tekst plakietki |
| Powiadomienie (`.dn-toast--blad`) | kreska lewa 2 px | ikona `blad` | `.dn-toast-tytul` + `.dn-toast-tresc` |
| Błąd pola | obrys `--dn-blad-tekst` | ikona przy komunikacie | `.dn-pole-blad` — treść wskazująca brak |
| Kropka stanu (`.dn-kropka--sukces` …) | wypełnienie kropki | kropka **nigdy nie stoi sama** | etykieta obok kropki |
| Wybór trwały (`.dn-karta--wybrana`) | tło `--dn-sygnal-tlo` | obrys `--dn-sygnal-obrys` (zmiana kształtu) | `aria-selected` / `aria-current` |
| Pozycja bieżąca nawigacji | tło + barwa tekstu | **wstęga 2 px** przy krawędzi | `aria-current="page"` |
| Zakładka aktywna | barwa tekstu | **podkreślenie 2 px** (`::after`) | `aria-selected="true"` |
| Krok kolejki | kreska lewa w barwie wyniku | znak w kółku (`.dn-krok-znak` + ikona) | etykieta kroku + `.dn-krok-meta` |
| Praca w tle | kropka sygnałowa | tętno (ruch) albo pierścień (przy `reduce`) | plakietka „pracuje" lub `aria-busy` |
| Nadawca w oknie komunikacji | kreska lewa 2 px | medalion z ikoną roli | nazwa nadawcy + plakietka roli |

### 19.3. Test weryfikacyjny

| Test | Sposób przeprowadzenia | Warunek zaliczenia |
|---|---|---|
| **Test monochromatyczny** | zrzut ekranu przekonwertowany do skali szarości | wszystkie stany nadal rozróżnialne |
| **Test opisu** | odczytanie widoku bez patrzenia — sam tekst dostępny | każdy stan ma nazwę słowną |
| **Test kropki** | wyszukanie każdego `.dn-kropka` w drzewie | każda ma sąsiadujący tekst albo `aria-label` |

---

## 20. Semantyka

### 20.1. Landmarki powłoki

| Obszar | Element / rola | Etykieta |
|---|---|---|
| pasek górny | `<header class="dn-pasek">` → `role="banner"` | — |
| boczna nawigacja modułów | `<nav class="dn-boczna">` | `aria-label="Moduły środowiska"` |
| pas kart sesji | `<div class="dn-karty-sesji" role="tablist">` | `aria-label="Karty sesji"` |
| obszar roboczy | `<main id="tresc">` | `aria-label` = nazwa okna operacyjnego |
| pas komunikacji | `<section>` z `aria-label="Chat Window"` | nazwa własna okna |
| powiadomienia | `<div class="dn-toasty" role="status" aria-live="polite">` | — |
| panel orkiestracji (MultitaskingAI) | `<nav aria-label="Panel orkiestracji">` | 6 sekcji |

**Jeden `<main>` na dokument.** Modal nie jest `<main>` — jest `<dialog>`.

### 20.2. Nagłówki

| Poziom | Zastosowanie |
|---|---|
| `h1` | tytuł okna operacyjnego — **dokładnie jeden na widok** |
| `h2` | sekcja okna (panel, strefa) |
| `h3` | podsekcja w panelu |
| `h4` | etykieta grupy — w systemie realizowana klasą `.dn-etykieta-wersalikowa` |

**Zakaz przeskoku poziomu** (`h2` → `h4`). Nagłówek dobiera się rangą w strukturze, nie stopniem pisma — stopień nadaje żeton `--dn-fs-*`.

### 20.3. Atrybuty `aria-*` per komponent

| Komponent | Atrybuty | Uwaga |
|---|---|---|
| `.dn-btn` przełączający | `aria-pressed="true|false"` | CSS reaguje selektorem `[aria-pressed='true']` — **stan jest w ARIA, nie w klasie** |
| `.dn-btn` w trakcie pracy | `aria-busy="true"` | uruchamia spinner `::after`; przycisk **pozostaje klikalny** |
| `.dn-btn` z błędem | `aria-invalid="true"` | obrys w barwie błędu |
| `.dn-zakladka` | `role="tab"`, `aria-selected`, `aria-controls` | rodzic `role="tablist"`, panel `role="tabpanel"` |
| `.dn-karta-sesji` | `role="tab"`, `aria-selected` | pas kart = `tablist` |
| `.dn-boczna-pozycja` | `aria-current="page"` | pozycja bieżąca |
| `.dn-karta-srodowiska` | `aria-current="true"` | środowisko bieżące |
| `.dn-tabela` wiersz | `aria-selected="true"` | wiersz wybrany |
| `.dn-pole-kontrolka` | `aria-invalid`, `aria-describedby` → `.dn-pole-opis` / `.dn-pole-blad` | etykieta przez `<label for>` |
| `.dn-przelacznik` | `<input type="checkbox" role="switch">` | stan natywny |
| `.dn-modal` | `<dialog>` + `aria-labelledby` → `.dn-modal-tytul` | `showModal()` |
| `.dn-tooltip-tresc` | `role="tooltip"`, wyzwalacz `aria-describedby` | — |
| `.dn-postep` | `role="progressbar"`, `aria-valuenow/min/max`, `aria-valuetext` | `aria-valuetext` z jednostką („12 z 34 zadań") |
| `.dn-spinner` | `aria-hidden="true"` + tekst statusu obok | spinner sam nie mówi nic |
| `.dn-kropka` | `aria-hidden="true"` | znaczenie niesie etykieta obok |
| `.dn-awatar` (inicjały) | `aria-hidden="true"` + nazwa obok | inicjały nie są treścią |

### 20.4. Ikony — dekoracyjne kontra informacyjne

| Rodzaj | Kryterium | Realizacja |
|---|---|---|
| **dekoracyjna** | obok jest tekst niosący to samo znaczenie | `aria-hidden="true"`, brak `<title>` |
| **informacyjna** | ikona jest **jedynym** nośnikiem znaczenia (przycisk ikonowy) | `aria-label` na **przycisku**, ikona nadal `aria-hidden="true"` |
| **ilustracyjna** | emblemat środowiska, godło | `role="img"` + `aria-label` |

```
   ┌─ przycisk ikonowy — etykieta na PRZYCISKU, nie na ikonie ─────────────┐
   │ <button class="dn-btn-ikona" aria-label="Zamknij kartę sesji">        │
   │   <svg aria-hidden="true">…</svg>                                     │
   │ </button>                                                             │
   └───────────────────────────────────────────────────────────────────────┘
```

### 20.5. Język dokumentu

| Element | Wartość |
|---|---|
| `<html lang="pl">` | obowiązkowe w każdym pliku |
| wstawka obcojęzyczna (nazwa własna okna, fragment kodu) | `lang="en"` na elemencie — dotyczy `Studio Editor`, `Workflow Builder`, `Chat Window`, `Coordinator`, `Executor 1` |
| dane techniczne (ścieżki, identyfikatory) | bez zmiany języka — czytnik odczyta znak po znaku dzięki krojowi mono i `aria-label` opisowemu |

---

## 21. Cele dotykowe i `pointer: coarse`

### 21.1. Mechanizm

Cele rosną **żetonem, nie wyjątkiem** — jedno zapytanie medialne w `zetony.css`, zero reguł per komponent:

```css
@media (pointer: coarse) {
  :root {
    --dn-wym-kontrolka: 40px;
    --dn-wym-ikonowy: 40px;
    --dn-wym-wiersz: 44px;
    --dn-wym-check: 20px;
    --dn-wym-przelacznik-szer: 44px;
    --dn-wym-przelacznik-wys: 24px;
  }
}
```

### 21.2. Wymiary porównawczo

| Element | Wskaźnik precyzyjny | `pointer: coarse` | Przyrost |
|---|---|---|---|
| kontrolka (przycisk, pole) | 32 px | **40 px** | +25% |
| przycisk ikonowy | 32 px | **40 px** | +25% |
| wiersz tabeli / listy | 36 px | **44 px** | +22% |
| pole wyboru | 16 px | **20 px** | +25% |
| przełącznik | 36 × 20 px | **44 × 24 px** | +22% |
| kropka sygnału | 6 px | 6 px | bez zmiany — nie jest celem |
| ikona | 14/16/20/24 px | bez zmiany | ikona nie jest celem; celem jest przycisk |

### 21.3. Reguły uzupełniające

| Reguła | Uzasadnienie |
|---|---|
| **Zapytanie o wskaźnik, nie o szerokość ekranu** | tablet z rysikiem ma wskaźnik precyzyjny mimo dużego ekranu; laptop dotykowy ma oba |
| **Odstęp między celami ≥ 8 px** (`--dn-od-2`) | sąsiadujące cele 44 px bez odstępu i tak dają błędne trafienia |
| **Cel dotykowy większy niż grafika** | ikona 16 px w przycisku 40 px — powierzchnia trafienia jest przyciskiem |
| **Zakaz gestów wielopunktowych jako jedynej drogi** | WCAG 2.5.1 — każda czynność osiągalna pojedynczym dotknięciem |
| **Zakaz działania na `pointerdown`** dla czynności nieodwracalnych | WCAG 2.5.2 — działanie na `pointerup` pozwala cofnąć przez odsunięcie palca |

---

## 22. Czytniki ekranu

### 22.1. Co ogłaszać przy zmianie stanu

| Zdarzenie | Rejon | Grzeczność | Treść komunikatu (przykładowa, z domeny produktu) |
|---|---|---|---|
| zapis zakończony | `.dn-toasty` | `polite` | „Profil agenta zapisany" |
| błąd operacji w tle | `.dn-toasty` | `assertive` | „Zadanie w kolejce wstrzymane — zdarzenie nieoczekiwane" |
| krok kolejki zmienił wynik | rejon kolejki | `polite` | „Krok 3 z 7: poprawny" |
| postęp przekroczył próg | `role="progressbar"` | — | odczyt z `aria-valuetext` przy zmianie fokusu, **bez** ogłaszania każdego procentu |
| nadawca zaczął pisać | rejon wpisów | `polite` | „Asystent pisze" |
| nowa wiadomość w Chat Window | rejon wpisów | `polite` | nadawca + treść |
| przełączenie motywu | — | — | **nie ogłaszać** — zmiana wizualna bez znaczenia semantycznego |
| otwarcie modala | — | — | natywne dla `<dialog>` — czytnik ogłasza tytuł z `aria-labelledby` |
| błąd pola formularza | pole | — | powiązanie `aria-describedby` — odczyt przy fokusie pola |

### 22.2. Reguły rejonów żywych

| Reguła | Uzasadnienie |
|---|---|
| **`polite` domyślnie** | `assertive` przerywa odczyt w toku; zarezerwowane dla błędów wymagających reakcji |
| **Rejon istnieje w drzewie przed wstawieniem treści** | rejon utworzony razem z treścią bywa pomijany |
| **Zakaz ogłaszania ruchu** | tętno, spinner i wjazdy warstw są nieme; ogłasza się **wynik**, nie animację |
| **Zakaz zalewu komunikatów** | postęp co 1% = 100 komunikatów; ogłaszać progi (start, zakończenie, błąd) |
| **`aria-busy="true"` na kontenerze podczas przebudowy listy** | czytnik nie odczytuje stanów pośrednich |

### 22.3. Teksty alternatywne — reguła doboru

| Element | Tekst |
|---|---|
| przycisk ikonowy | **czynność**, nie nazwa ikony: „Zamknij kartę sesji", nie „ikona X" |
| kropka pracy | brak (`aria-hidden`) — pracę ogłasza `aria-busy` albo plakietka |
| awatar z inicjałami | brak (`aria-hidden`) — nazwa stoi obok |
| emblemat środowiska w karcie | nazwa środowiska: „TalkIn" |
| godło w pasku | „Danaco Console" |
| wykres / pasek postępu | `aria-valuetext` z jednostką: „12 z 34 zadań" |

---

## 23. Lista kontrolna dostępności przed oddaniem okna

Do przejścia **w obu motywach**, przed zgłoszeniem okna do odbioru.

### 23.1. Kontrast

- [ ] Każda para tekst/tło zmierzona algorytmem WCAG (rozdz. 16.1), nie oceniona wzrokiem
- [ ] Tekst normalny ≥ 4,5 : 1; tekst duży i elementy graficzne ≥ 3 : 1
- [ ] `--dn-tekst-3` wyłącznie na metadanych, etykietach wersalikowych i tekście dużym
- [ ] Barwy półprzezroczyste złożone z podłożem przed pomiarem
- [ ] Pierścień fokusu zmierzony na każdym podłożu, na którym występuje (w tym na ramie atramentowej)

### 23.2. Fokus i klawiatura

- [ ] Pierścień 2 px + odsunięcie 2 px widoczny na **każdej** kontrolce
- [ ] Zero `outline: none` bez równoważnego wskaźnika
- [ ] Całe okno przechodzone `Tab`-em; kolejność zgodna z układem wizualnym
- [ ] Zero `tabindex` dodatniego
- [ ] Grupy (zakładki, karty sesji, opcje jednokrotne) — jeden przystanek, strzałki wewnątrz
- [ ] Modal: pułapka fokusu, `Esc` zamyka, fokus wraca do wyzwalacza
- [ ] Odnośnik „Przejdź do treści" pierwszym elementem `<body>`
- [ ] Zero skrótów jednoklawiszowych bez modyfikatora

### 23.3. Stan i barwa

- [ ] Każdy stan niesie co najmniej dwa nośniki (barwa + kształt/ikona lub słowo)
- [ ] Test monochromatyczny zaliczony — wszystkie stany rozróżnialne w skali szarości
- [ ] Każda `.dn-kropka` ma sąsiadującą etykietę
- [ ] Zero `disabled`; niegotowość opisana obok albo komunikatem po naciśnięciu

### 23.4. Semantyka

- [ ] `lang="pl"` na dokumencie; wstawki obcojęzyczne oznaczone
- [ ] Dokładnie jeden `h1`; brak przeskoków poziomów nagłówków
- [ ] Landmarki obecne: `banner`, `nav` (z etykietą), `main`, rejon powiadomień
- [ ] Każde pole ma etykietę powiązaną `for`/`id`
- [ ] `aria-pressed` / `aria-selected` / `aria-current` / `aria-busy` odzwierciedlają stan faktyczny
- [ ] Ikony dekoracyjne `aria-hidden="true"`; ikony informacyjne z etykietą na kontrolce

### 23.5. Ruch

- [ ] Zero animacji spoza katalogu dozwolonych (rozdz. 10)
- [ ] Dokładnie jeden ruch ciągły na widok
- [ ] `prefers-reduced-motion` sprawdzone w przeglądarce z włączoną preferencją
- [ ] Tętno przy `reduce` zastąpione pierścieniem statycznym — informacja zachowana
- [ ] Zero animacji własności układu (poza `.dn-postep-wartosc`)

### 23.6. Dotyk i wskaźnik

- [ ] Sprawdzone przy `pointer: coarse` — kontrolki 40 px, wiersze 44 px
- [ ] Odstęp między sąsiadującymi celami ≥ 8 px
- [ ] Zero gestów wielopunktowych jako jedynej drogi do funkcji

### 23.7. Czytniki ekranu

- [ ] Powiadomienia w rejonie `aria-live="polite"`; błędy `assertive`
- [ ] Rejony żywe obecne w drzewie przed wstawieniem treści
- [ ] Zero ogłaszania ruchu i zmian czysto wizualnych
- [ ] `aria-valuetext` z jednostką na wskaźnikach postępu

---

## 24. Decyzje projektowe

Rozstrzygnięcia podjęte ponad dokumentację źródłową, wraz z uzasadnieniem.

| # | Rozstrzygnięcie | Sytuacja w źródłach | Decyzja i uzasadnienie |
|---|---|---|---|
| **D-01** | **`DOSTEPNOSC.md` v1.0 obowiązuje strukturą, nie wartościami** | Plik pakietu systemu wizualnego (v1.0, 2026-08-11) opisuje paletę złoto-granatową, pierścień fokusu `#C9A24B`, `--dn-tekst-3` = `#6E7A93` i dziewięć pomiarów. `KIERUNEK.md` v2.0 zastępuje całą warstwę wizualną. | Z v1.0 przenosimy **reguły**: progi WCAG, zasadę „tekst-3 tylko metadane", zasadę „stan nigdy samym kolorem", geometrię fokusu 2+2 px, globalną obsługę `prefers-reduced-motion`, zero blokad ADL-017. **Wartości barwne v1.0 są nieobowiązujące** — zastąpione żetonami v2.0. Dziewięć pomiarów v1.0 zastąpione **33 pomiarami** z `kontrasty.json`. |
| **D-02** | **Próg systemowy 1,6 : 1 dla obrysu dekoracyjnego** | `kontrasty.json` zawiera pomiar 13 (obrys kontrolki / tło = 1,65) z polem `prog: 1.6` — wartość spoza skali WCAG. | Wprowadzono **czwarty próg systemowy** obok 4,5 / 3,0. Dotyczy wyłącznie obrysów oddzielających powierzchnie o tej samej roli, nienoszących informacji. Obrys niosący informację (błąd pola, wybór) podlega progowi 3,0. Werdykt w tabeli nie brzmi „AA", lecz „próg systemowy" — żeby nie sugerować zgodności z kryterium, którego pomiar nie dotyczy. |
| **D-03** | **`kontrasty.json` przenosi się do `motyw/`** | `HANDOFF.md` §1 nie wymienia `kontrasty.json` w mapowaniu plików, choć §3 czyni z niego wejście bramy weryfikacyjnej. | Plik dodany do mapowania jako `motyw/kontrasty.json`. Brama wymagająca „przebiegu `kontrasty.json`" nie może działać na pliku, którego w repozytorium docelowym nie ma. |
| **D-04** | **`motyw/siatka.css` jako dwunasty arkusz** | `HANDOFF.md` wymienia dziesięć arkuszy motywu; sekcja 7 `zetony.css` (punkty łamania, siatka) nie ma przypisania. | Utworzono `siatka.css`. Alternatywa — dołączenie do `przestrzen.css` — mieszałaby dwie odpowiedzialności (rytm wewnętrzny kontra rusztowanie strony), łamiąc regułę R1. |
| **D-05** | **Kolejność importu wymuszona przez współdzielone klatki kluczowe** | Podział `komponenty.css` na 14 arkuszy rozdziela definicje `@keyframes` od części zastosowań (`dn-tetno`, `dn-wejscie`, `dn-obrot`). | Zapisano właściciela każdej klatki i listę zastosowań (rozdz. 2.3). Alternatywa — piętnasty arkusz `ruch-komponentow.css` — została odrzucona: klatka kluczowa animacji tętna należy do kropki sygnału, a rozdzielenie utrudniałoby czytanie komponentu. |
| **D-06** | **Rozszerzona lista czego się nie przenosi** | `HANDOFF.md` §1 nie wymienia `wspolne.js`, `prototyp.css`, `prototyp.js`, `WZORZEC-OKNA.html` ani katalogu `logo/alternatywy/`. | Dodano rozdz. 2.5. Bez jawnego wykluczenia rusztowanie makiet trafiłoby do produktu jako „pliki z pakietu". |
| **D-07** | **Katalog ruchu jako katalog zamknięty (7 pozycji + 10 zakazów)** | `KIERUNEK.md` podaje wartość pokrętła i wzorzec ogólny; nie wylicza animacji. | Wyprowadzono siedem pozycji dozwolonych **z realnego kodu** `komponenty.css` (każda ma nośnik i selektor) oraz dziesięć zakazów z anty-domyślnych i z fizyki renderowania. Katalog jest zamknięty: animacja spoza niego wymaga decyzji Właściciela — inaczej pokrętło 3/10 nie ma egzekwowalnej treści. |
| **D-08** | **Tętno wobec kryterium WCAG 2.2.2 (Pauza, zatrzymanie, ukrycie)** | Kryterium wymaga mechanizmu zatrzymania dla treści ruchomej trwającej ponad 5 s. Tętno jest ciągłe. | Rozstrzygnięcie: tętno jest **wskaźnikiem stanu**, a nie ruchomą treścią prezentacyjną — informuje, że praca biegnie, i mieści się w wyjątku „ruch jest istotny dla działania". Niezależnie od tego system dostarcza mechanizm wyłączenia — `prefers-reduced-motion` z zamiennikiem statycznym (rozdz. 13.2), przy zachowaniu informacji. Częstotliwość 0,42 Hz jest siedmiokrotnie poniżej progu 2.3.1. |
| **D-09** | **`--dn-czas-3` opisany jako 0,22 s, nie 0,24 s** | Katalog komponentów v1.0 opisuje wjazd modala jako „0,24 s (`--dn-czas-3`)"; `zetony.css` v2.0 definiuje `--dn-czas-3: 0.22s`. | Obowiązuje **wartość żetonu (0,22 s)** — żeton jest źródłem prawdy (Z1). Katalog v1.0 opisywał poprzednią generację. |
| **D-10** | **Dwa czasy poza skalą uznane i uzasadnione** | `zetony.css` deklaruje cztery czasy, a `komponenty.css` używa dodatkowo `0.8s linear` (spinner). | Uznano jako **jedyny dopuszczony wyjątek** z uzasadnieniem (rozdz. 9.4): ruch cykliczny bez punktu docelowego wymaga krzywej liniowej. Nie tworzy się dla niego żetonu, żeby nie sugerować, że skala czasów ma pięć pozycji. |
| **D-11** | **Skróty klawiszowe zebrane z modułów, `Ctrl/Cmd + K` jako propozycja** | Skróty rozproszone po plikach modułów. `Ctrl/Cmd + K` w znaczeniu „paleta poleceń" występuje **wyłącznie** w `propozycje-rozbudowy/`; w module Research ten sam skrót oznacza „Szybkie dodanie źródła". | Do tabeli 18.2 weszły wyłącznie skróty z dokumentacji modułów (kontrakt). Paleta poleceń środowiska pozostaje **propozycją rozbudowy** i nie jest wpisana jako obowiązująca — mimo że warstwa `--dn-z-centrum-polecen` (1300) istnieje w żetonach. |
| **D-12** | **Kolejność tabulacji powłoki ustalona wprost** | Dokumentacja opisuje anatomię powłoki, nie kolejność fokusu. | Ustalono kolejność zgodną z układem wizualnym: pasek → karty sesji → boczna nawigacja → obszar roboczy → pas komunikacji (rozdz. 17.2). Zgodne z kryterium 2.4.3 (kolejność zachowująca sens) i z hierarchią środowisko → moduł → okno. |
| **D-13** | **Rola `tool` mapowana na `--system` z obowiązkową plakietką** | `HANDOFF.md` §2 podaje mapowanie `tool → system` „z plakietką roli = nazwa narzędzia", nie rozstrzygając, czy plakietka jest obowiązkowa. | Plakietka roli **obowiązkowa** dla `tool`. Bez niej wpis narzędzia jest nieodróżnialny od komunikatu systemowego, co łamie zasadę „stan nigdy samym kolorem" przeniesioną na tożsamość nadawcy. |
| **D-14** | **Trzy testy weryfikacyjne zasady „stan nigdy samym kolorem"** | Zasada jest w źródłach nienegocjowalna, ale nie ma procedury sprawdzenia. | Dodano test monochromatyczny, test opisu i test kropki (rozdz. 19.3). Zasada bez procedury sprawdzenia jest deklaracją, nie kontraktem. |
| **D-15** | **`will-change` nieużywany, z zapisaną zasadą trzech warunków** | Biblioteka nie deklaruje `will-change`; źródła nie poruszają tematu. | Utrzymano stan zerowy i zapisano warunki wprowadzenia (rozdz. 14.3). Animowane elementy są małe; przedwczesna promocja do warstwy kompozycji kosztuje pamięć graficzną bez zysku. |
| **D-16** | **Pary bez pomiaru wskazane jawnie** | `kontrasty.json` mierzy pary komponentowe, nie iloczyn palety; brak pomiarów na `--dn-panel` i dla barw półprzezroczystych motywu ciemnego. | Dodano rozdz. 16.6 z listą par wymagających pomiaru przed nowym zestawieniem. Milczenie o lukach czytałoby się jako komplet. |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
