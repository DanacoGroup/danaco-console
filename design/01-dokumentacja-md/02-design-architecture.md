# Danaco Console — Design Architecture (architektura projektowa)

| Pozycja | Treść |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — architektura systemu projektowego |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 (opracowanie A2, zespół A — Fundament systemu projektowego) |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 (pakiet źródłowy: 2026-08-11) |
| **Odbiorcy** | Deweloper wdrażający warstwę wizualną do `budowa/client`, projektant systemu, osoba prowadząca odbiór wizualny, autor kolejnych opracowań pakietu |
| **Zakres** | Model warstwowy żetonów, architektura plików pakietu i katalogu docelowego, regulamin budowy arkuszy, mechanizm motywu i gęstości, skala warstw `z-index`, siatka i punkty łamania, kompozycja powłoki i okna operacyjnego, kontrakt komponent ↔ żeton, architektura ikon, architektura ruchu, bramy weryfikacyjne fal wdrożeniowych |
| **Czego NIE zawiera** | Wartości żetonów jako katalogu barw (opracowanie A1 — system żetonów), katalogu komponentów z pełnymi stanami (opracowanie zespołu B), księgi znaku i portfolio marki (zespół C), makiet okien operacyjnych (zespół D), specyfikacji funkcjonalnej modułów i logiki serwera |

---

## Spis treści

1. [Model warstwowy systemu wizualnego](#1-model-warstwowy-systemu-wizualnego)
2. [Architektura plików](#2-architektura-plików)
3. [Regulamin budowy](#3-regulamin-budowy)
4. [Mechanizm motywu](#4-mechanizm-motywu)
5. [Mechanizm gęstości](#5-mechanizm-gęstości)
6. [Architektura warstw (z-index)](#6-architektura-warstw-z-index)
7. [Architektura siatki i punktów łamania](#7-architektura-siatki-i-punktów-łamania)
8. [Architektura kompozycji okna — powłoka](#8-architektura-kompozycji-okna--powłoka)
9. [Model kompozycji okna operacyjnego](#9-model-kompozycji-okna-operacyjnego)
10. [Kontrakt komponent ↔ żeton](#10-kontrakt-komponent--żeton)
11. [Architektura ikon](#11-architektura-ikon)
12. [Architektura ruchu](#12-architektura-ruchu)
13. [Bramy weryfikacyjne fal wdrożeniowych](#13-bramy-weryfikacyjne-fal-wdrożeniowych)
14. [Decyzje projektowe](#14-decyzje-projektowe)

---

## 1. Model warstwowy systemu wizualnego

### 1.1. Cztery warstwy — deklaracja

Nagłówek `zasoby/zetony/zetony.css` deklaruje trzy warstwy żetonów. Czwartą warstwę
— komponentową — tworzy `zasoby/css/komponenty.css`. Razem dają model
czterowarstwowy, w którym **kierunek zależności jest jednokierunkowy i nieodwracalny**.

```
   ┌───────────────────────────────────────────────────────────────────────┐
   │  W4 · WARSTWA KOMPONENTOWA          komponenty.css · 119 klas .dn-*    │
   │      .dn-btn · .dn-pole · .dn-karta · .dn-wpis · .dn-modal …           │
   │      wolno czytać: W2 i W3        ·  ZAKAZ czytania: W1               │
   └───────────────────────────────▲───────────────────────────────────────┘
                                   │  var(--dn-tlo), var(--dn-od-3), var(--dn-z-modal)
   ┌───────────────────────────────┴────────────┬──────────────────────────┐
   │  W2 · SEMANTYKA PER MOTYW                  │  W3 · NIEZALEŻNE OD       │
   │      :root[data-theme='light']             │       MOTYWU              │
   │      :root[data-theme='dark']              │  typografia · przestrzeń  │
   │      @media (prefers-color-scheme)         │  ruch · wymiary · siatka  │
   │      role: tło, powierzchnia, tekst,       │  warstwy · gradienty      │
   │      obrys, sygnał, stany, cienie          │  rama kokpitu            │
   │      wolno czytać: W1                      │  wolno czytać: W1        │
   └───────────────────────────────▲────────────┴──────────────────────────┘
                                   │  var(--dn-szary-900), var(--dn-sygnal-500)
   ┌───────────────────────────────┴───────────────────────────────────────┐
   │  W1 · PRYMITYWY                     surowe skale — 18 szarości,        │
   │      --dn-szary-* · --dn-sygnal-*    8 stopni sygnału, 12 stanów       │
   │      --dn-zielen/bursztyn/czerwien-* ZAKAZ użycia poza W2/W3          │
   └───────────────────────────────────────────────────────────────────────┘
```

### 1.2. Charakterystyka warstw

| Warstwa | Nazwa | Miejsce w `zetony.css` | Selektor | Co zawiera | Kto ją czyta |
|---|---|---|---|---|---|
| **W1** | Prymitywy | sekcja 1 (l. 18–58) | `:root` | 18 kroków szarości, 8 stopni sygnału, 12 wartości stanów | wyłącznie W2 i W3 |
| **W2** | Semantyka per motyw | sekcje 10–11 + blok zapasowy (l. 198–401) | `:root[data-theme='light' \| 'dark']`, `@media (prefers-color-scheme)` | 26 ról semantycznych × 2 motywy, 4 rodziny stanów × 2 motywy, 5 cieni × 2 motywy | W4 |
| **W3** | Niezależne od motywu | sekcje 2–9 oraz 12–14 (l. 60–195, 403–448) | `:root`, `@media (pointer: coarse)`, `:root[data-gestosc]`, `@media (prefers-reduced-motion)` | rama kokpitu, typografia, przestrzeń, ruch, wymiary, siatka i łamanie, warstwy, gradienty | W4 |
| **W4** | Komponenty | `komponenty.css` (1207 linii) | klasy `.dn-*` | 119 tokenów klasowych w 18 sekcjach tematycznych | okna i prototypy |

### 1.3. Reguła kierunku — trzy zdania wiążące

1. **Prymitywu nie wolno użyć wprost w komponencie.** `komponenty.css` sięga po
   `var(--dn-tlo)`, nigdy po `var(--dn-szary-950)`. Jedyne udokumentowane odstępstwo:
   `.dn-awatar--inteligencja` używa `var(--dn-szary-0)` jako barwy tekstu na gradiencie
   sygnałowym — bo gradient jest identyczny w obu motywach i semantyka „tekst na
   gradiencie" nie ma odrębnego żetonu.
2. **Semantyka nie wywołuje semantyki drugiego motywu.** Motyw jasny i ciemny są
   definiowane osobno, nie wywodzone przez inwersję ani filtr (KIERUNEK.md 3.1:
   „Oba motywy równoprawne — definiowane osobno, nie wywodzone").
3. **Warstwa niezależna od motywu nie zawiera barwy** — z jednym wyjątkiem:
   **rama kokpitu** (sekcja 2), która jest barwą stałą w obu motywach i dlatego
   należy do W3, nie do W2.

### 1.4. Dlaczego rama kokpitu jest w warstwie niezależnej od motywu

`--dn-rama`, `--dn-rama-tekst`, `--dn-rama-tekst-2`, `--dn-rama-hover`,
`--dn-rama-obrys` są zdefiniowane raz, w bloku `:root`, poza obydwoma motywami.
KIERUNEK.md 3.3: *„Pasek górny jest zawsze atramentowy (`szary-925`) — w obu
motywach. […] To jedyny element, który nie przełącza się z motywem (poza godłem)."*

Konsekwencja architektoniczna: **każdy komponent osadzony na ramie musi mieć
własny modyfikator ramowy**, bo nie może polegać na żetonach motywu.
W `komponenty.css` realizuje to `.dn-btn-ikona--na-ramie` (l. 159–165) oraz
`.dn-pasek-szukaj` (l. 675–691), które czytają wyłącznie `--dn-rama-*`.

| Element na ramie | Klasa | Żetony, których używa |
|---|---|---|
| Godło i logotyp | `.dn-pasek-godlo`, `.dn-pasek-logotyp` | `--dn-rama-tekst`, `--dn-rama-tekst-2` |
| Wyszukiwarka paska | `.dn-pasek-szukaj` | `--dn-rama-hover`, `--dn-rama-obrys`, `--dn-rama-tekst`, `--dn-fokus` |
| Przyciski ikonowe paska | `.dn-btn-ikona--na-ramie` | `--dn-rama-tekst-2`, `--dn-rama-hover`, `--dn-rama-tekst` |
| Plakietka środowiska na pasku | `.dn-plakietka--rola` (nadpisanie ramowe) | `--dn-rama-hover`, `--dn-rama-tekst` |

### 1.5. Warstwa komponentowa — 18 sekcji tematycznych

| # | Sekcja `komponenty.css` | Linie | Liczba linii | Klasy bazowe |
|---|---|---:|---:|---|
| 1 | Przycisk | 11–164 | 154 | `.dn-btn`, `.dn-btn-ikona` |
| 2 | Pole formularza | 166–233 | 68 | `.dn-pole`, `.dn-szukaj` |
| 3 | Wybór: pole wyboru · przełącznik · radio | 235–333 | 99 | `.dn-wybor`, `.dn-check`, `.dn-radio`, `.dn-przelacznik`, `.dn-suwak` |
| 4 | Kropka sygnału | 335–360 | 26 | `.dn-kropka` |
| 5 | Plakietka — stany i role | 362–393 | 32 | `.dn-plakietka` |
| 6 | Karta | 395–523 | 129 | `.dn-karta`, `.dn-karta-srodowiska`, `.dn-kafel` |
| 7 | Tabela | 525–559 | 35 | `.dn-tabela` |
| 8 | Zakładki · karty sesji | 561–637 | 77 | `.dn-zakladki`, `.dn-zakladka`, `.dn-karty-sesji`, `.dn-karta-sesji` |
| 9 | Rama kokpitu — pasek + listwa | 639–718 | 80 | `.dn-pasek`, `.dn-listwa` |
| 10 | Boczna nawigacja modułów | 720–771 | 52 | `.dn-boczna` |
| 11 | Wpis okna komunikacji | 773–862 | 90 | `.dn-wpis` |
| 12 | Monitor wykonania · kolejka kroków | 864–945 | 82 | `.dn-postep`, `.dn-kolejka`, `.dn-krok` |
| 13 | Modal · nakładka | 947–987 | 41 | `.dn-modal` |
| 14 | Powiadomienie | 989–1024 | 36 | `.dn-toasty`, `.dn-toast` |
| 15 | Dymek objaśnienia | 1026–1058 | 33 | `.dn-tooltip` |
| 16 | Awatar · spinner · pusty stan | 1060–1121 | 62 | `.dn-awatar`, `.dn-spinner`, `.dn-pusty-stan` |
| 17 | Pole wpisywania promptu | 1123–1169 | 47 | `.dn-prompt`, `.dn-przybornik` |
| 18 | Always On Display | 1171–1207 | 37 | `.dn-aod` |

---

## 2. Architektura plików

### 2.1. Mapa pakietu design — źródło

```
design/
├── 01-kierunek/
│   └── KIERUNEK.md            kontrakt kierunku, trzy pokrętła, anty-domyślne
├── 02-marka/
│   ├── logo/                  sygnet · logotyp · lockupy · mono · kontrowe (+png/)
│   ├── favicon/               favicon.svg/ico/png · apple-touch · manifest · snippet
│   ├── ikona-aplikacji/       kafel SVG + PNG 180/192/512/1024 + wariant maskowalny
│   └── srodowiska/            4 emblematy: TalkIn · WorkSpace · CodeStudio · MultitaskingAI
├── 03-zetony/
│   ├── zetony.css             448 linii · 14 sekcji · źródło prawdy wartości
│   ├── zetony.json            odczyt maszynowy (te same wartości)
│   ├── kontrasty.json         33 zmierzone pary WCAG + paleta odniesienia
│   ├── fonty.css              @font-face z unicode-range (latin + latin-ext)
│   └── fonty/                 pliki .woff2 + licencje OFL 1.1
├── 04-css/
│   ├── fundament.css          157 linii · warstwa zerowa nad żetonami
│   └── komponenty.css         1207 linii · 119 klas .dn-* · 18 sekcji
├── 05-ikony/
│   ├── svg/                   82 pliki SVG (siatka 24, obrys 1,75)
│   ├── manifest.json          82 wpisy: nazwa · źródło · zastosowanie
│   └── ikony.html             galeria przeglądowa
├── 06-okna/                   E1 · E2+E3 · E4 · komponenty.html · wspolne.js
├── 07-ksiega/                 ksiega-marki.html + zrzuty/
├── 08-grafika/                og-image · banner-linkedin · tapeta-2560 · tlo-slajdu
├── README.md                  przewodnik po pakiecie + rejestr zamkniętych luk
└── HANDOFF.md                 mapowanie wdrożeniowe na budowa/client
```

### 2.2. Zależności między katalogami pakietu

```
  01-kierunek  ────────────────────────────────────────────────┐
      │  (kontrakt: monochrom, pokrętła, anty-domyślne)         │
      ▼                                                        ▼
  03-zetony ──► 04-css/fundament.css ──► 04-css/komponenty.css ──► 06-okna
      │                                        ▲                    ▲
      │                                        │                    │
      └──► 02-marka (barwy własne znaku)  05-ikony ─────────────────┘
                    │                                              │
                    └──► 08-grafika ◄──────────────────────────────┘
                                          07-ksiega (opisuje całość)
```

**Odczyt diagramu:** `01-kierunek` nie jest importowany przez nic — jest
kontraktem tekstowym. `03-zetony` jest korzeniem technicznym. `02-marka`
i `05-ikony` są warstwami zasobowymi: ikony czytają barwę przez `currentColor`
(czyli pośrednio z żetonów), znak marki ma barwy własne wpisane w SVG.

### 2.3. Struktura docelowa w `budowa/client/src` (wg HANDOFF.md §1)

```
budowa/client/
├── public/
│   ├── favicon.svg · favicon.ico · favicon-16/32/48.png
│   ├── apple-touch-icon.png · icon-192.png · icon-512.png
│   └── site.webmanifest
└── src/
    ├── motyw/
    │   ├── motyw.css                arkusz spinający — importuje całość
    │   ├── prymitywy.css            zetony.css §1 (szarość + sygnał)
    │   ├── stany.css                zetony.css §1 (stany) + §10/§11 (stany semantyczne)
    │   ├── rama.css                 zetony.css §2
    │   ├── typografia.css           zetony.css §3
    │   ├── przestrzen.css           zetony.css §4 + cienie z §10/§11
    │   ├── ruch.css                 zetony.css §5 + §14
    │   ├── wymiary.css              zetony.css §6 + §12 + §13
    │   ├── siatka.css               zetony.css §7        ← plik dodany (rozdz. 14, D-01)
    │   ├── warstwy.css              zetony.css §8
    │   ├── gradienty.css            zetony.css §9
    │   ├── semantyczne-jasny.css    zetony.css §10 + blok zapasowy „light"
    │   ├── semantyczne-ciemny.css   zetony.css §11 + blok zapasowy „dark"
    │   ├── fundament.css            04-css/fundament.css bez zmian
    │   ├── zetony.json              podmiana w całości
    │   └── fonty.css + fonty/       @font-face + .woff2 + licencje
    ├── komponenty/
    │   ├── przycisk.css · pole.css · wybor.css · plakietka.css · karta.css
    │   ├── tabela.css · zakladki.css · pasek.css · boczna.css · wpis.css
    │   ├── postep.css · nakladka.css · powiadomienie.css · drobne.css
    │   └── komponenty.css           arkusz spinający — importuje powyższe
    ├── ikony/
    │   ├── svg/                     82 pliki — PODMIANA CAŁEGO KATALOGU
    │   ├── manifest.json            nowy plik (82 pozycje)
    │   ├── zrodla-ikon.ts           wykaz nazw — literówka zatrzymuje kompilację
    │   └── ikony.ts                 komponent renderujący
    └── zasoby/
        └── marka/                   sygnet · logotyp · lockupy · emblematy środowisk
```

### 2.4. Mapowanie plik po pliku

| Pakiet design | Cel w `budowa/client` | Charakter operacji |
|---|---|---|
| `03-zetony/zetony.css` | `src/motyw/*.css` (13 plików) | podział wg regulaminu, wartości przenoszone **dosłownie** |
| `03-zetony/zetony.json` | `src/motyw/zetony.json` | podmiana w całości |
| `03-zetony/fonty.css` + `fonty/*.woff2` | `src/motyw/fonty.css` + `src/motyw/fonty/` | przeniesienie 1:1 z licencjami OFL |
| `04-css/fundament.css` | `src/motyw/fundament.css` | import po żetonach, przed komponentami; **usunąć duplikat fokusu** z dotychczasowego `motyw.css` |
| `04-css/komponenty.css` | `src/komponenty/*.css` (14 plików) | podział per komponent; selektory `.dn-*` bez zmian |
| `05-ikony/svg/*.svg` | `src/ikony/svg/` | **podmiana całego katalogu** — poprzedni zestaw był zastępczy |
| `05-ikony/manifest.json` | `src/ikony/manifest.json` | nowy plik |
| `02-marka/logo/sygnet*.svg`, `logotyp*.svg` | `src/ikony/svg/` (sygnet jako `logo-danaco.svg`) **oraz** `src/zasoby/marka/` | duplikat świadomy: kompatybilność importów + katalog marki |
| `02-marka/favicon/*` | `client/public/` | snippet `<head>` gotowy w `naglowek-snippet.html` |
| `02-marka/ikona-aplikacji/*` | `desktop/` (Tauri: 32, 128, 128@2x z `ikona-1024.png`) + `client/public/` (PWA) | wariant maskowalny dla manifestu |
| `02-marka/srodowiska/*.svg` | `src/ikony/svg/` | emblematy używane przez Centrum dowodzenia i pasek ramy |

### 2.5. Kolejność importu — kontrakt kaskady

```
  index.css
    │
    ├─ 1. motyw/fonty.css          (@font-face — musi być pierwszy)
    ├─ 2. motyw/motyw.css          (żetony: W1 → W3 → W2)
    │      ├─ prymitywy.css
    │      ├─ rama.css
    │      ├─ typografia.css
    │      ├─ przestrzen.css
    │      ├─ ruch.css
    │      ├─ wymiary.css
    │      ├─ siatka.css
    │      ├─ warstwy.css
    │      ├─ gradienty.css
    │      ├─ semantyczne-jasny.css
    │      ├─ semantyczne-ciemny.css
    │      └─ stany.css
    ├─ 3. motyw/fundament.css      (wyzerowanie, podłoże, fokus globalny)
    └─ 4. komponenty/komponenty.css (14 arkuszy per komponent)
```

**Reguła:** żaden arkusz z poziomu 4 nie może być zaimportowany przed poziomem 2.
Odwrócenie kolejności nie wywoła błędu — wywoła **ciche fałszowanie wartości**
(zmienne niezdefiniowane w chwili odczytu dają wartość pustą, a nie wyjątek).
To jedyny błąd tej architektury, którego kompilator nie zgłasza.

---

## 3. Regulamin budowy

### 3.1. Trzy zasady wiążące (HANDOFF.md, nagłówek)

| # | Zasada | Egzekwowanie |
|---|---|---|
| R1 | **Jedna odpowiedzialność = jeden plik** | podział per komponent; arkusz nie łączy dwóch niepowiązanych komponentów |
| R2 | **Arkusz ≤ 300 linii** | próg twardy; przekroczenie = sygnał do dalszego podziału |
| R3 | **Komponent sięga wyłącznie po żetony semantyczne** | przegląd kodu: wystąpienie `--dn-szary-*` / `--dn-sygnal-[0-9]` w `komponenty/` jest błędem |

### 3.2. Podział `zetony.css` (448 linii → 13 arkuszy + arkusz spinający)

| Arkusz docelowy | Sekcje źródłowe | Linie źródła | Szacowana objętość | Zawartość |
|---|---|---|---:|---|
| `prymitywy.css` | §1 (szarość, sygnał) | 22–50 | ~45 | 18 kroków szarości, 8 stopni sygnału |
| `stany.css` | §1 (stany) + §10/§11 (stany semantyczne) + odpowiadające bloki zapasowe | 52–58, 240–251, 303–314, 346–355, 383–394 | ~70 | 12 prymitywów stanów, 4 rodziny × 3 role × 2 motywy |
| `rama.css` | §2 | 60–67 | ~20 | 5 żetonów ramy kokpitu (stałe w obu motywach) |
| `typografia.css` | §3 | 69–98 | ~45 | 3 kroje, 9 stopni, 4 wagi, 3 interlinie, 3 tracking |
| `przestrzen.css` | §4 + cienie z §10/§11 i bloków zapasowych | 100–121, 253–259, 316–321, 356–360, 395–399 | ~60 | 11 odstępów, 6 promieni, 5 cieni × 2 motywy |
| `ruch.css` | §5 + §14 | 123–130, 432–448 | ~40 | 1 ease, 4 czasy, globalna obsługa `prefers-reduced-motion` |
| `wymiary.css` | §6 + §12 + §13 | 132–161, 403–415, 417–430 | ~70 | 28 wymiarów, 6 nadpisań dotykowych, 9 nadpisań gęstości przestronnej |
| `siatka.css` | §7 | 163–172 | ~25 | 4 punkty łamania, `--dn-tresc-max`, 12 kolumn, przerwa |
| `warstwy.css` | §8 | 174–187 | ~28 | 11 poziomów `--dn-z-*` |
| `gradienty.css` | §9 | 189–195 | ~20 | 2 gradienty + zapisany zakres użycia |
| `semantyczne-jasny.css` | §10 (bez stanów i cieni) + blok zapasowy „light" | 201–238, 328–345 | ~90 | 26 ról semantycznych motywu jasnego |
| `semantyczne-ciemny.css` | §11 (bez stanów i cieni) + blok zapasowy „dark" | 265–301, 364–382 | ~90 | 26 ról semantycznych motywu ciemnego |
| `motyw.css` | — (nowy) | — | ~25 | wyłącznie `@import` w ustalonej kolejności + nagłówek |

**Suma docelowa:** ~628 linii wobec 448 linii źródła. Przyrost ≈ 40% to nagłówki
komentarzowe i domknięcia selektorów, które w pliku scalonym występowały raz.
Przyrost jest kosztem świadomym i zapisanym: R1 przed zwięzłością.

### 3.3. Podział `komponenty.css` (1207 linii → 14 arkuszy + arkusz spinający)

| Arkusz docelowy | Sekcje źródłowe | Linie źródła | Linii | Klasy przenoszone |
|---|---|---:|---:|---|
| `przycisk.css` | Przycisk | 11–164 | 154 | `.dn-btn` + 9 modyfikatorów, `.dn-btn-ikona` + `--na-ramie` |
| `pole.css` | Pole formularza | 166–233 | 68 | `.dn-pole`, `-etykieta`, `-kontrolka`, `-opis`, `-blad`, `.dn-szukaj` |
| `wybor.css` | Wybór | 235–333 | 99 | `.dn-wybor`, `.dn-check`, `.dn-radio`, `.dn-przelacznik`, `.dn-suwak` |
| `plakietka.css` | Kropka + Plakietka | 335–393 | 58 | `.dn-kropka` + 5 modyfikatorów, `.dn-plakietka` + 6 modyfikatorów |
| `karta.css` | Karta | 395–523 | 129 | `.dn-karta` + 2 modyfikatory + 3 podklasy, `.dn-karta-srodowiska` + 4 podklasy, `.dn-kafel` + 3 podklasy |
| `tabela.css` | Tabela | 525–559 | 35 | `.dn-tabela`, `.dn-dane` (komórka) |
| `zakladki.css` | Zakładki · karty sesji | 561–637 | 77 | `.dn-zakladki`, `.dn-zakladka`, `.dn-karty-sesji`, `.dn-karta-sesji` |
| `pasek.css` | Rama kokpitu | 639–718 | 80 | `.dn-pasek` + 4 podklasy, `.dn-listwa` + `-pozycja` |
| `boczna.css` | Boczna nawigacja | 720–771 | 52 | `.dn-boczna`, `-naglowek`, `-pozycja` |
| `wpis.css` | Wpis okna komunikacji | 773–862 | 90 | `.dn-wpis` + 4 modyfikatory + 5 podklas |
| `postep.css` | Monitor wykonania · kolejka | 864–945 | 82 | `.dn-postep` + 3 podklasy, `.dn-kolejka`, `.dn-krok` + 4 modyfikatory + 2 podklasy |
| `nakladka.css` | Modal | 947–987 | 41 | `.dn-modal` + 4 podklasy, `::backdrop`, `@keyframes dn-wejscie` |
| `powiadomienie.css` | Toast | 989–1024 | 36 | `.dn-toasty`, `.dn-toast` + 4 modyfikatory + 2 podklasy |
| `drobne.css` | Tooltip + Awatar/Spinner/Pusty stan + Prompt + AOD | 1026–1207 | 179 | `.dn-tooltip`, `.dn-awatar`, `.dn-spinner`, `.dn-pusty-stan`, `.dn-prompt`, `.dn-przybornik`, `.dn-aod` |
| `komponenty.css` | — (nowy) | — | ~20 | wyłącznie `@import` |

**Kontrola progu R2:** największy arkusz to `drobne.css` (179 linii) — mieści się.
Drugi w kolejności `przycisk.css` (154) i `karta.css` (129) również. Żaden arkusz
nie zbliża się do progu 300 linii, więc podział z HANDOFF.md jest wystarczający
bez dalszego rozdrabniania.

### 3.4. Reguła dalszego podziału — kiedy dzielić ponownie

```
   arkusz rośnie
        │
        ├── < 240 linii  ─────────────────► zostaw jak jest
        │
        ├── 240–300 linii ────────────────► ostrzeżenie: zaplanuj podział
        │                                    przy najbliższej zmianie
        │
        └── > 300 linii  ─────────────────► podziel NATYCHMIAST:
                                             1) wydziel modyfikatory do
                                                <komponent>-warianty.css
                                             2) jeśli nie wystarczy — wydziel
                                                stany do <komponent>-stany.css
```

Kolejność cięcia jest ustalona, bo utrzymuje kaskadę: baza → warianty → stany.
Odwrócenie tej kolejności w importach cofa specyficzność i psuje `:hover`.

### 3.5. Zakazy w warstwie komponentowej

| Zakaz | Uzasadnienie | Gdzie egzekwowany |
|---|---|---|
| Wartość szesnastkowa wprost | żetony są jedynym źródłem prawdy | przegląd kodu; wyjątek: brak |
| Prymitywy `--dn-szary-*` / `--dn-sygnal-[0-9]{3}` | łamie kontrakt warstw (rozdz. 1.3) | przegląd kodu; wyjątek udokumentowany: `.dn-awatar--inteligencja` |
| `!important` | kaskada jest jednokierunkowa i wystarczająca | przegląd kodu |
| Atrybut `disabled` | zasada zero blokad (ADL-017) | przegląd kodu + lista sprawdzeń |
| Zagnieżdżenie selektorów głębsze niż 2 poziomy | koszt utrzymania i nieprzewidywalna specyficzność | przegląd kodu |
| Media query barwna wewnątrz komponentu | motyw obsługiwany wyłącznie w warstwie żetonów | przegląd kodu |
| `@media (prefers-reduced-motion)` wewnątrz komponentu | obsługa jest globalna (zetony.css §14) | przegląd kodu |

---

## 4. Mechanizm motywu

### 4.1. Trzy elementy mechanizmu

```
   ┌──────────────────────────────────────────────────────────────────────┐
   │  1. data-theme            atrybut na <html> — wybór jawny Operatora  │
   │       ├─ "light"  ──► :root[data-theme='light']  { … 26 ról … }      │
   │       └─ "dark"   ──► :root[data-theme='dark']   { … 26 ról … }      │
   ├──────────────────────────────────────────────────────────────────────┤
   │  2. prefers-color-scheme  gdy atrybutu NIE MA — decyduje system      │
   │       @media (prefers-color-scheme: light) :root:not([data-theme])   │
   │       @media (prefers-color-scheme: dark)  :root:not([data-theme])   │
   ├──────────────────────────────────────────────────────────────────────┤
   │  3. color-scheme          deklaracja dla kontrolek natywnych         │
   │       ustawiana W KAŻDYM z czterech bloków                           │
   └──────────────────────────────────────────────────────────────────────┘
```

### 4.2. Tabela stanów mechanizmu

| Stan `<html>` | Preferencja systemu | Selektor rozstrzygający | Wynik | `color-scheme` |
|---|---|---|---|---|
| `data-theme="light"` | dowolna | `:root[data-theme='light']` | motyw jasny | `light` |
| `data-theme="dark"` | dowolna | `:root[data-theme='dark']` | motyw ciemny | `dark` |
| brak atrybutu | jasna | `@media (prefers-color-scheme: light) :root:not([data-theme])` | motyw jasny | `light` |
| brak atrybutu | ciemna | `@media (prefers-color-scheme: dark) :root:not([data-theme])` | motyw ciemny | `dark` |
| brak atrybutu | brak deklaracji systemu | żaden blok W2 nie wchodzi | **W1 i W3 działają, W2 pusta** | nieustawiony |

Ostatni wiersz jest ważny: bez atrybutu i bez preferencji systemowej role
semantyczne pozostają nierozwinięte. Dlatego `wspolne.js` **zawsze** ustawia
`data-theme` przy starcie — czyta zapis z `localStorage`, a przy jego braku
odczytuje `matchMedia('(prefers-color-scheme: dark)')` i zapisuje wynik jawnie.
Mechanizm zapasowy w `@media` obsługuje wyłącznie dokumenty statyczne, otwierane
bez skryptu.

### 4.3. Dlaczego powielenie wartości w `@media` jest świadome

`zetony.css` l. 324–326 zapisuje to wprost:

> *„Brak jawnego wyboru — rozstrzyga preferencja systemu. Wartości muszą być
> identyczne z blokami motywów; powielenie jest świadome (mechanizm kaskady,
> nie drugie źródło prawdy — patrz HANDOFF.md)."*

Rozwinięcie argumentacji:

| Rozwiązanie | Dlaczego odrzucone / przyjęte |
|---|---|
| `:root[data-theme='light'], @media(...) :root:not([data-theme])` — wspólna lista selektorów | **Odrzucone.** Selektora spoza media query nie można połączyć przecinkiem z selektorem wewnątrz media query — to nie jest legalna składnia CSS. |
| Warstwa pośrednia `--dn-motyw-tlo` wskazywana raz | **Odrzucone.** Wprowadza piąty poziom pośrednictwa dla całej palety; koszt czytelności większy niż koszt powielenia. |
| Preprocesor generujący oba bloki z jednego źródła | **Odrzucone.** Pakiet ma działać z `file://` bez budowania (KANON §11.1). Preprocesor przenosi źródło prawdy poza `zetony.css`. |
| `@container style()` / `light-dark()` | **Odrzucone na dziś.** Cel produkcyjny to Tauri WebView2; oparcie mechanizmu motywu na funkcji o niepewnym wsparciu byłoby ryzykiem bez zysku. |
| **Powielenie jawne, oznaczone komentarzem** | **Przyjęte.** Jedno źródło prawdy pozostaje w sekcjach 10–11; blok `@media` jest ich mechanicznym odbiciem. |

### 4.4. Kontrakt utrzymaniowy powielenia

> **Zmiana wartości w sekcji 10 lub 11 musi być w tym samym zatwierdzeniu
> odzwierciedlona w odpowiadającym bloku `@media (prefers-color-scheme)`.**

Po podziale na pliki (rozdz. 3.2) blok zapasowy trafia do tego samego arkusza,
co blok główny danego motywu — `semantyczne-jasny.css` zawiera oba bloki jasne,
`semantyczne-ciemny.css` oba ciemne. To celowe: rozbieżność jest wtedy widoczna
w jednym otwartym pliku, a nie rozrzucona po dwóch.

### 4.5. Co się przełącza, a co nie

| Element | Przełącza się z motywem | Uzasadnienie |
|---|---|---|
| Tło strony, powierzchnie, panele | tak | rdzeń motywu |
| Tekst (3 stopnie) i obrysy (3 stopnie) | tak | rdzeń motywu |
| Sygnał (tekst, wypełnienie, tło, obrys, kropka) | tak | inne stopnie rodziny dla różnych kontrastów |
| Cienie (5 stopni) | tak | dwa komplety tonowane per motyw |
| Stany (sukces, ostrzeżenie, błąd, informacja) | tak | pełny komplet obu motywów (zamknięcie luki L-P-19) |
| **Pasek górny (rama kokpitu)** | **nie** | KIERUNEK.md 3.3 — stanowisko dowodzenia ma stały dom |
| **Godło i emblematy środowisk** | **nie** | geometria i barwy własne znaku (HANDOFF §4) |
| Gradienty `--dn-grad-*` | nie | zdefiniowane raz na prymitywach, zakres ilustracyjny |
| Typografia, przestrzeń, wymiary, ruch, warstwy | nie | warstwa W3 z definicji |

### 4.6. Animacja przełączenia

`fundament.css` nakłada na `body` przejście `background-color` i `color` czasem
`--dn-czas-2` (0,16 s) z krzywą `--dn-ease`. Przy `prefers-reduced-motion`
czas skraca się globalnie do 0,01 ms (zetony.css §14) — bez osobnej reguły.
Żaden komponent nie definiuje własnego przejścia motywu.

---

## 5. Mechanizm gęstości

### 5.1. Dwie osie, nie jedna

Gęstość w Danaco Console ma **dwie niezależne osie**, które mogą działać razem:

```
   OŚ 1 — URZĄDZENIE (automatyczna)          OŚ 2 — PREFERENCJA (jawna)
   @media (pointer: coarse)                  :root[data-gestosc='przestronna']
        │                                          │
        │ wykrycie wskaźnika grubego               │ decyzja Operatora
        │ (dotyk, rysik)                           │ (ustawienia · wygląd)
        ▼                                          ▼
   6 nadpisań wymiarów                        9 nadpisań: wymiary + typografia
   (bez zmiany typografii)                    (stopień bazowy 13 → 14 px)
```

### 5.2. Tabela nadpisań — obie osie

| Żeton | Zwarta (domyślna) | `pointer: coarse` | `data-gestosc="przestronna"` |
|---|---:|---:|---:|
| `--dn-fs-base` | 13 px | 13 px | **14 px** |
| `--dn-lh-bazowy` | 1.45 | 1.45 | **1.5** |
| `--dn-wym-kontrolka` | 32 px | **40 px** | **40 px** |
| `--dn-wym-ikonowy` | 32 px | **40 px** | **40 px** |
| `--dn-wym-wiersz` | 36 px | **44 px** | **44 px** |
| `--dn-wym-check` | 16 px | **20 px** | 16 px |
| `--dn-wym-przelacznik-szer` | 36 px | **44 px** | 36 px |
| `--dn-wym-przelacznik-wys` | 20 px | **24 px** | 20 px |
| `--dn-wym-pasek` | 48 px | 48 px | **56 px** |
| `--dn-wym-pas-kart` | 36 px | 36 px | **40 px** |
| `--dn-odstep-panel` | 12 px | 12 px | **20 px** |
| `--dn-odstep-sekcji` | 24 px | 24 px | **32 px** |

### 5.3. Kolejność kaskady obu osi

W `zetony.css` blok `@media (pointer: coarse)` (§12, l. 403–415) stoi **przed**
blokiem `:root[data-gestosc='przestronna']` (§13, l. 417–430). Oba mają
selektor o specyficzności `:root` (0,1,0), więc rozstrzyga kolejność zapisu:
**gęstość przestronna wygrywa z dotykiem** tam, gdzie oba nadpisują ten sam żeton.

Praktycznie oznacza to: na tablecie z ustawioną gęstością przestronną kontrolka
ma 40 px (obie osie zgadzają się), a pole wyboru wraca do 16 px, bo gęstość
przestronna nie nadpisuje `--dn-wym-check`.

> **Rozstrzygnięcie utrzymaniowe:** przy podziale na pliki (rozdz. 3.2) obie
> sekcje trafiają do `wymiary.css` — w tej samej kolejności co w źródle.
> Rozdzielenie ich na dwa pliki wprowadziłoby zależność od kolejności importów
> tam, gdzie dziś zależność jest zamknięta w jednym arkuszu.

### 5.4. Zasada „żeton, nie wyjątek"

KIERUNEK.md 3.7 oraz README pakietu (L-P-14) rozstrzygają: cele dotykowe rosną
**żetonem**, nie regułą per komponent. Konsekwencja: `komponenty.css` nie zawiera
ani jednego `@media (pointer: coarse)`. Komponent, który czyta
`var(--dn-wym-kontrolka)`, dostosowuje się bez własnego kodu.

**Sprawdzenie negatywne:** komponent, który wpisałby 32 px wprost, przestałby
reagować na obie osie. To jeden z powodów, dla których zakaz wartości wprost
(rozdz. 3.5) jest zakazem twardym, a nie zaleceniem.

### 5.5. Wartość domyślna i jej status

Gęstość **zwarta** jest decyzją Właściciela (README, L-P-13; KIERUNEK.md
`GESTOSC_WIZUALNA = 8/10`). HANDOFF §4 zapisuje: *„Gęstość — domyślnie zwarta;
przełącznik przyszłościowo przez `data-gestosc="przestronna"` na `<html>`
(żetony gotowe)."* HANDOFF §4 wymienia ją także wśród rzeczy, których nie wolno
zmienić po cichu.

Wariant przestronny jest **przygotowany, nie wdrożony** — żetony istnieją,
przełącznik w Oknie Ustawień (zakres „wygląd") jest miejscem przewidzianym
na jego udostępnienie.

---

## 6. Architektura warstw (z-index)

### 6.1. Jedenaście poziomów — pełna tabela

| # | Żeton | Wartość | Warstwa | Co się na niej znajduje | Dlaczego tu, a nie wyżej/niżej |
|---:|---|---:|---|---|---|
| 1 | `--dn-z-podloga` | 0 | Podłoga | obszar roboczy, okna operacyjne, treść dokumentu | poziom odniesienia; wszystko inne mierzy się względem niego |
| 2 | `--dn-z-przybornik` | 10 | Przybornik | przyklejone nagłówki tabel (`.dn-tabela th`), przybornik promptu | musi zakryć przewijaną treść, ale zostaje w obrębie swojego okna |
| 3 | `--dn-z-pasek` | 100 | Pasek górny | `.dn-pasek` — rama kokpitu | rama zamyka kompozycję; przechodzi nad treścią wszystkich okien |
| 4 | `--dn-z-boczna` | 200 | Boczna nawigacja | `.dn-boczna`, panel orkiestracji | nad paskiem w osi pionowej cieni i menu rozwijanych; nawigacja ma pierwszeństwo przed treścią |
| 5 | `--dn-z-pas-komunikacji` | 300 | Pas komunikacji | Chat Window w trybie pływającym / nakładającym | wspólny dla wszystkich modułów — nie może zniknąć pod panelem modułu |
| 6 | `--dn-z-nakladka` | 800 | Nakładka | przyciemnienie pod modalem (`::backdrop`) | przerwa 300 → 800 rezerwuje pasmo na przyszłe warstwy powłoki bez przenumerowania |
| 7 | `--dn-z-modal` | 900 | Modal | `.dn-modal` | dokładnie jeden krok nad własną nakładką |
| 8 | `--dn-z-powiadomienie` | 1000 | Powiadomienie | `.dn-toasty`, `.dn-toast` | komunikat o wyniku operacji musi być widoczny także przy otwartym modalu |
| 9 | `--dn-z-tooltip` | 1100 | Dymek objaśnienia | `.dn-tooltip-tresc` | objaśnia elementy modala i toasta — musi być nad nimi |
| 10 | `--dn-z-aod` | 1200 | Always On Display | `.dn-aod` | funkcja globalna, obecna ponad bieżącą przestrzenią roboczą (dokumentacja E12) |
| 11 | `--dn-z-centrum-polecen` | 1300 | Centrum poleceń | Command Center (`Ctrl K`) | **zawsze najwyżej** — wywoływane z każdego stanu, także znad AOD |

### 6.2. Stos w perspektywie

```
                                            ┌─────────────────────────────┐
                                       1300 │  CENTRUM POLECEŃ            │
                                            └─────────────────────────────┘
                                       1200 │  ALWAYS ON DISPLAY          │
                                       1100 │  DYMEK OBJAŚNIENIA          │
                                       1000 │  POWIADOMIENIE (toast)      │
                                        900 │  MODAL                      │
                                        800 │  NAKŁADKA (::backdrop)      │
              ── pasmo rezerwowe 301–799 ─────────────────────────────────
                                        300 │  PAS KOMUNIKACJI            │
                                        200 │  BOCZNA NAWIGACJA           │
                                        100 │  PASEK GÓRNY (rama)         │
              ── pasmo rezerwowe 11–99 ──────────────────────────────────
                                         10 │  PRZYBORNIK / nagłówki tabel│
                                          0 │  PODŁOGA — obszar roboczy   │
                                            └─────────────────────────────┘
```

### 6.3. Trzy zasady skali

| Zasada | Treść | Skutek praktyczny |
|---|---|---|
| **Z1 — skok dziesiętny** | poziomy rosną krokami 10 / 100, nigdy o 1 | zawsze jest miejsce na warstwę pośrednią bez przenumerowania reszty |
| **Z2 — dwa pasma rezerwowe** | 11–99 (wnętrze okna) oraz 301–799 (powłoka) | nowa warstwa powłoki nie zmusza do ruszania warstw nakładkowych |
| **Z3 — brak `z-index` wprost** | komponent zapisuje `z-index: var(--dn-z-*)`, nigdy liczbę | skala jest jednym miejscem prawdy; audyt to jeden `grep` |

### 6.4. Faktyczne użycie w `komponenty.css`

| Linia | Klasa | Żeton | Kontekst |
|---:|---|---|---|
| 536 | `.dn-tabela th` | `--dn-z-przybornik` | `position: sticky` — nagłówek nad przewijanymi wierszami |
| 995 | `.dn-toasty` | `--dn-z-powiadomienie` | `position: fixed`, prawy dolny róg |
| 1034 | `.dn-tooltip-tresc` | `--dn-z-tooltip` | dymek nad każdą warstwą poza AOD i Centrum poleceń |
| 1178 | `.dn-aod` | `--dn-z-aod` | `position: fixed`, pływający awatar |

Modal i nakładka nie potrzebują deklaracji `z-index`: natywny `<dialog>`
otwarty przez `showModal()` trafia do warstwy najwyższej przeglądarki, a jego
`::backdrop` jest z definicji tuż pod nim. Żetony `--dn-z-nakladka`
i `--dn-z-modal` pozostają w skali dla **zapasowych implementacji modala
w prototypach** oraz dla warstw naśladujących modal bez elementu `<dialog>`.

### 6.5. Dlaczego boczna nawigacja jest nad paskiem

Pasek górny (100) i boczna nawigacja (200) nie nachodzą na siebie geometrycznie
— pasek zajmuje pas 48 px u góry, boczna zaczyna się pod nim. Kolejność ma
znaczenie dla **elementów wysuwanych**: menu kontekstowe pozycji bocznej
nawigacji, rozwinięcie „Subagent Network" w powłoce MultitaskingAI oraz cień
rzucany przez boczną w widoku zwiniętym (próg w2) muszą przechodzić nad krawędzią
paska, gdy boczna rozwija się jako panel nakładkowy.

---

## 7. Architektura siatki i punktów łamania

### 7.1. Cztery progi

| Żeton | Wartość | Nazwa robocza | Klasa urządzenia | Rozstrzygnięcie źródłowe |
|---|---:|---|---|---|
| `--dn-bp-w1` | 640 px | w1 | telefon poziomo | widok mobilny (okno Mobile, E11) |
| `--dn-bp-w2` | 960 px | w2 | tablet | boczna nawigacja zwija się do ikon |
| `--dn-bp-w3` | 1280 px | w3 | biurko | pełny kokpit |
| `--dn-bp-w4` | 1600 px | w4 | szerokie biurko | dwa okna komunikacji |

Uzupełniająco: `--dn-tresc-max` 1200 px (maks. szerokość treści dokumentowej),
`--dn-siatka-kolumny` 12, `--dn-siatka-przerwa` 24 px (`--dn-od-6`).

### 7.2. Co robi powłoka na każdym progu

```
  < w1 (640)          w1–w2 (640–959)     w2–w3 (960–1279)   w3–w4 (1280–1599)   ≥ w4 (1600)
  ┌──────────┐        ┌──────────────┐    ┌───────────────┐  ┌────────────────┐  ┌──────────────────┐
  │▓▓ pasek ▓│        │▓▓▓ pasek ▓▓▓▓│    │▓▓▓▓ pasek ▓▓▓▓│  │▓▓▓▓▓ pasek ▓▓▓▓│  │▓▓▓▓▓▓ pasek ▓▓▓▓▓│
  ├──────────┤        ├──────────────┤    ├──┬────────────┤  ├────┬───────────┤  ├────┬─────────────┤
  │  karty   │        │ karty sesji  │    │ik│ karty sesji│  │ bo │ karty ses.│  │ bo │ karty sesji │
  ├──────────┤        ├──────────────┤    │on├────────────┤  │ cz ├───────────┤  │ cz ├──────┬──────┤
  │          │        │              │    │y │            │  │ na ├───────────┤  │ na │ okno │ okno │
  │  jedno   │        │  okno        │    │  │  okno      │  │224 │  okno     │  │224 │ wiod.│ wsp. │
  │  okno    │        │  wiodące     │    │48│  wiodące   │  │ px │  wiodące  │  │ px ├──────┴──────┤
  │  naraz   │        ├──────────────┤    │px├────────────┤  │    ├───────────┤  │    │ Chat  Chat  │
  │          │        │  Chat Window │    │  │Chat Window │  │    │Chat Window│  │    │  A  │   B   │
  └──────────┘        └──────────────┘    └──┴────────────┘  └────┴───────────┘  └────┴─────────────┘
   boczna jako          boczna ukryta       boczna = ikony     boczna pełna       boczna pełna
   szuflada             (wywołanie menu)    48 px              224 px             + druga kolumna
```

### 7.3. Tabela zachowań — element po elemencie

| Element powłoki | < w1 | w1–w2 | w2–w3 | w3–w4 | ≥ w4 |
|---|---|---|---|---|---|
| **Pasek górny** | 48 px, godło + `menu` + AOD | 48 px, bez wyszukiwarki inline | 48 px, wyszukiwarka skrócona | 48 px, komplet | 48 px, komplet |
| **Wyszukiwarka paska** | ukryta — wywołanie ikoną | ukryta — wywołanie ikoną | 220 px | elastyczna do 520 px | elastyczna do 520 px |
| **Pas kart sesji** | przewijany poziomo, bez `+` inline | przewijany poziomo | przewijany, `+` na końcu | pełny | pełny |
| **Boczna nawigacja** | szuflada nakładkowa | ukryta, wywołanie `menu` | **zwinięta do ikon** (etykieta w dymku) | 224 px, etykiety | 224 px, etykiety |
| **Obszar roboczy** | jedno okno naraz, przełączanie zakładkami | okno wiodące + Chat Window | okno wiodące + Chat Window | wiodące + wspomagające w zakładkach + Chat Window | wiodące + wspomagające obok + **dwa Chat Window** |
| **Chat Window** | pełna wysokość, osobny widok | pas dolny, min. 320 px | pas dolny, min. 320 px | pas dolny 320 px / 54% wysokości | **dwa pasy** obok siebie (para koordynator–wykonawca) |
| **Modal** | pełna szerokość minus 32 px | `min(560 px, 100vw − 32 px)` | 560 px | 560 px | 560 px |
| **Toast** | pełna szerokość minus 32 px | 280–420 px, prawy dolny | 280–420 px | 280–420 px | 280–420 px |
| **AOD** | zwinięty do rdzenia | zwinięty do rdzenia | rdzeń + treść | rdzeń + treść | rdzeń + treść |

### 7.4. Kierunek zapytań mediowych

Punkty łamania są zapisane jako `min-width` **w kierunku rosnącym** — projekt
zaczyna od najwęższego układu i rozwija go progami. Wyjątkiem są makiety okien
w pakiecie (E1, E4), które używają `max-width` dla nielicznych korekt siatki
wewnątrz sekcji; to korekty lokalne, nie progi architektoniczne.

| Zapis | Zakres | Zastosowanie |
|---|---|---|
| `@media (min-width: 640px)` | ≥ w1 | wyjście z widoku jednokolumnowego |
| `@media (min-width: 960px)` | ≥ w2 | boczna nawigacja pojawia się jako pas stały (ikony) |
| `@media (min-width: 1280px)` | ≥ w3 | boczna rozwija się do 224 px, wspomagające w zakładkach |
| `@media (min-width: 1600px)` | ≥ w4 | druga kolumna okien, dwa pasy komunikacji |

**Uwaga o żetonach w zapytaniach mediowych:** wartości `--dn-bp-*` nie mogą być
użyte wprost w warunku `@media` (zmienne CSS nie są tam interpretowane).
Żetony pełnią rolę **rejestru wartości** — punktu odniesienia dokumentacyjnego
i wartości dostępnej ze skryptu. W arkuszach progi zapisuje się liczbą,
a zgodność z rejestrem sprawdza przegląd kodu.

### 7.5. Siatka 12 kolumn

Siatka 12-kolumnowa z przerwą 24 px obowiązuje **stronę główną (Centrum
dowodzenia) i okna dokumentowe**, nie powłokę roboczą. Powłoka używa siatki
o kolumnach nazwanych wymiarami (`--dn-wym-boczna` + `1fr`), bo jej podział
jest funkcjonalny, nie proporcjonalny.

| Kontekst | Model siatki | Podstawa |
|---|---|---|
| Centrum dowodzenia — Strefa 1 | 4 kolumny kart środowisk (→ 2 poniżej 1100 px) | makieta E1 |
| Centrum dowodzenia — Strefa 2 | 4 kolumny kafli (→ 2 poniżej 1100 px) | makieta E1 |
| Centrum dowodzenia — Strefa 3 | listwa jednorzędowa, 3 pozycje | makieta E1 |
| Powłoka środowiska | `grid-template-columns: var(--dn-wym-boczna) 1fr` | makieta E2 |
| Para koordynator–wykonawca | `grid-template-columns: 1fr 1fr` | makieta E4 |
| Treść dokumentowa | jedna kolumna, `max-width: var(--dn-tresc-max)` | żeton `--dn-tresc-max` |

---

## 8. Architektura kompozycji okna — powłoka

### 8.1. Cztery pasy powłoki

Dokumentacja (`elementy-okien.md` 3.1, Schemat 4) definiuje powłokę jako
cztery pasy: **pasek górny → karty sesji → boczna nawigacja | obszar roboczy**.
Obszar roboczy zawiera okna operacyjne modułu i Chat Window.

```
 ╔══════════════════════════════════════════════════════════════════════════╗
 ║  PASEK GÓRNY  ── 48 px ── zawsze atramentowy, z-index 100                ║
 ║  [godło+logotyp] [plakietka środowiska] [wyszukiwarka] [dom][Mobile]     ║
 ║                                                        [AOD][motyw]      ║
 ╠══════════════════════════════════════════════════════════════════════════╣
 ║  PAS KART SESJI  ── 36 px ──  [● Studio ✕][Research ✕][+]                ║
 ╠═══════════════════════╤══════════════════════════════════════════════════╣
 ║  BOCZNA NAWIGACJA     │  OBSZAR ROBOCZY                                  ║
 ║  224 px · z-index 200 │  min-height: 0 · 1fr                             ║
 ║                       │  ┌────────────────────────────────────────────┐  ║
 ║  ┌ nagłówek: TalkIn ┐ │  │  OKNO WIODĄCE (np. Studio Editor)          │  ║
 ║  │ Studio         ◄─┼─┼──┤  nagłówek okna + obszar zawartości         │  ║
 ║  │ Workspace        │ │  └────────────────────────────────────────────┘  ║
 ║  │ Browser          │ │  ┌────────────────────────────────────────────┐  ║
 ║  │ Research         │ │  │  CHAT WINDOW — pas komunikacji, min 320 px │  ║
 ║  │ Library          │ │  │  historia · monitor · prompt · przybornik  │  ║
 ║  │ Translate        │ │  └────────────────────────────────────────────┘  ║
 ║  │ Roundtable       │ │                                                  ║
 ║  │ Assistant        │ │                                                  ║
 ║  │ Agents           │ │                                                  ║
 ║  └──────────────────┘ │                                                  ║
 ╚═══════════════════════╧══════════════════════════════════════════════════╝
```

### 8.2. Model grid — trzy poziomy zagnieżdżenia

| Poziom | Selektor roboczy | Deklaracja | Skąd |
|---|---|---|---|
| 1 — powłoka | `.powloka` | `display: grid; grid-template-rows: var(--dn-wym-pasek) var(--dn-wym-pas-kart) 1fr; height: 100dvh` | makieta E2, E4 |
| 2 — ciało | `.cialo` | `display: grid; grid-template-columns: var(--dn-wym-boczna) 1fr; min-height: 0` | makieta E2 |
| 3 — obszar główny | `.glowna` | `display: grid; grid-template-rows: minmax(0, 1fr) minmax(320px, 54%); min-width: 0; min-height: 0` | makieta E2 |

W makiecie E4 (para koordynator–wykonawca) poziom 3 przyjmuje proporcje
`minmax(0, 58%) minmax(0, 42%)`, a wewnątrz górnego wiersza wchodzi czwarty
poziom: `.para { grid-template-columns: 1fr 1fr }`.

### 8.3. Wymiary powłoki — komplet

| Element | Żeton | Wartość zwarta | Wartość przestronna | Wartość dotykowa |
|---|---|---:|---:|---:|
| Pasek górny | `--dn-wym-pasek` | 48 px | 56 px | 48 px |
| Pas kart sesji | `--dn-wym-pas-kart` | 36 px | 40 px | 36 px |
| Boczna nawigacja (szerokość) | `--dn-wym-boczna` | 224 px | 224 px | 224 px |
| Pozycja bocznej (wysokość) | `--dn-wym-kontrolka` | 32 px | 40 px | 40 px |
| Pas komunikacji (minimum) | `--dn-wym-pas-komunikacji` | 320 px | 320 px | 320 px |
| Wiersz tabeli w oknie | `--dn-wym-wiersz` | 36 px | 44 px | 44 px |
| Rytm wewnętrzny panelu | `--dn-odstep-panel` | 12 px | 20 px | 12 px |
| Rytm między sekcjami | `--dn-odstep-sekcji` | 24 px | 32 px | 24 px |

### 8.4. Reguła `min-height: 0` — dlaczego jest architektoniczna

Element siatki (`grid item`) ma domyślnie `min-height: auto`, co oznacza:
**nie skurczy się poniżej rozmiaru swojej treści**. W powłoce o wysokości
`100dvh` to prowadzi do jednego z dwóch błędów:

```
   BEZ min-height: 0                       Z min-height: 0
   ┌──────────────────┐                    ┌──────────────────┐
   │ pasek            │                    │ pasek            │
   ├──────────────────┤                    ├──────────────────┤
   │ karty sesji      │                    │ karty sesji      │
   ├──────────────────┤                    ├──────────────────┤
   │ obszar roboczy   │                    │ obszar roboczy   │
   │ …                │                    │ ┌──────────────┐ │
   │ …treść rośnie…   │                    │ │ przewijanie  │ │
   │ …               ▼│ ← powłoka pęka     │ │ WEWNĄTRZ     │ │
   └──────────────────┘   i wypycha        │ └──────────────┘ │
   (albo cała strona       stopkę poza     └──────────────────┘
    dostaje pasek          ekran            wysokość zamknięta
    przewijania)
```

**Reguła:** każdy pojemnik siatki, który ma przewijać własną treść, deklaruje
`min-height: 0` (oraz `min-width: 0` w osi poziomej). W makietach E2 i E4 reguła
występuje na `.cialo`, `.glowna`, `.okno-operacyjne`, `.okno-komunikacji`, `.para`
i `.dol` — czyli na **każdym** pojemniku pośrednim między `100dvh` a elementem
z `overflow-y: auto`.

Dokumentacja źródłowa nie formułuje tego jako reguły; wynika ona z makiet
i zostaje podniesiona do rangi reguły architektonicznej (rozdz. 14, D-09).

### 8.5. Powłoka MultitaskingAI — różnice

| Cecha | Powłoka modułowa | Powłoka MultitaskingAI |
|---|---|---|
| Pasek górny | `.dn-pasek`, komplet elementów | **identyczny — bez zmian** |
| Pas kart sesji | karta = jeden moduł | karta = jeden **proces orkiestracji** |
| Panel boczny | lista modułów (9/9/8 pozycji) | **panel orkiestracji — 6 sekcji stałych** |
| Zawartość obszaru | okna modułu + Chat Window | zawartość wybranej sekcji panelu |
| Tytuł karty przy zmianie | zmienia się na nazwę modułu | **nie zmienia się** — pozostaje nazwą procesu |
| Poziom izolacji | „karta sesji" (6 z 7) | karta + zagnieżdżony poziom **„Rola"** (najwęższy, z pierwszeństwem) |

Sześć sekcji panelu orkiestracji (kolejność domyślna, konfigurowalna):
**Zespoły · Role · Kolejki · Orkiestracja · Harmonogram i automatyki ·
Monitor procesu.** Sekcja „Role" pokazuje dokładnie cztery karty:
**Executor 1 · Executor 2 · Coordinator · Executor 3 / Validator**.

**Model grid pozostaje ten sam** — zmienia się wyłącznie zawartość drugiej
kolumny poziomu 2. To rozstrzygnięcie architektoniczne: powłoka jest jedna,
panel boczny jest wymienny.

### 8.6. Co przeładowanie obszaru roboczego obejmuje, a czego nie

```
   ZMIANA MODUŁU (Etap 6)
   ┌───────────────────────────────┬───────────────────────────────────────┐
   │  PRZEŁADOWUJE SIĘ             │  NIE ZMIENIA SIĘ                      │
   ├───────────────────────────────┼───────────────────────────────────────┤
   │  • zestaw okien operacyjnych  │  • sama karta sesji                   │
   │  • narzędzia modułu           │  • boczna nawigacja / panel orkiestr. │
   │  • kontekst Chat Window       │  • pozostałe otwarte karty            │
   │    (rekonfiguracja, NIE       │  • funkcje globalne Mobile i AOD      │
   │     zniknięcie okna)          │  • pasek górny                        │
   │  • tytuł karty sesji          │                                       │
   │  • stan wybrany w bocznej     │                                       │
   └───────────────────────────────┴───────────────────────────────────────┘
```

Chat Window **nie znika przy zmianie modułu** — rekonfiguruje kontekst poleceń
i dostępnych operacji. To najważniejsza konsekwencja architektoniczna dla
kompozycji: pas komunikacji jest elementem powłoki, nie elementem modułu.

---

## 9. Model kompozycji okna operacyjnego

### 9.1. Trzy role okien

| Rola | Definicja | Liczność | Pozycja w kompozycji |
|---|---|---|---|
| **Okno wiodące** | okno, które nadaje modułowi tożsamość; otwiera się jako pierwsze po Chat Window | dokładnie 1 na moduł | największa masa, górna część obszaru roboczego |
| **Okna wspomagające** | panele narzędziowe, podglądy, repozytoria, monitory | 1–5 na moduł | zakładki nad oknem wiodącym albo kolumna boczna (≥ w4) |
| **Chat Window** | jedyne okno wspólne wszystkim piętnastu modułom | dokładnie 1, zawsze | pas dolny, minimum 320 px |

### 9.2. Okna wiodące — pełny wykaz

| Moduł | Okno wiodące | Okna wspomagające |
|---|---|---|
| Studio | **Studio Editor** | Tools Panel · Diff/Grep Panel · Session Repository · Preview Window |
| Research | **Research Workspace** | Sources Manager · Findings Panel · Report Builder · Export Panel |
| Library | **Library Explorer** | Tags & Collections · File Preview · Versioning Panel |
| Translate | **Translation Panels** | Source Panel · Glossary Manager |
| Browser | **Browser Window** | Sources Panel · Notes Panel |
| Assistant | **Voice Console** | Actions Monitor · Activity Feed |
| Roundtable | **Model Panels** | Debate Panel · Moderator Panel · Consensus Panel |
| Workspace | **Project Dashboard** | Instructions Panel · Context Memory · Project Library · Agent Manager |
| Automations | **Workflow Builder** | Scheduler · Queue Manager · Orchestrator · Execution Monitor |
| Design | **Design Board** | Assets Panel · Prompt Builder · Preview Window |
| Apps | **Product Builder** | Architecture Designer · Frontend Workspace · Backend Workspace · Deployment Panel |
| Terminal | **Terminal Tabs** | Output Console · Process Monitor |
| Developer | **Code Editor** | Project Tree · Git Panel · Build Output |
| Diagnostics | **Diagnostics Center** | Logs Viewer · Errors Panel · Recommendations Panel |
| Agents | **Agent Builder** | Model Configuration · Skills Manager · Connectors Manager · Permissions Center |

### 9.3. Anatomia okna operacyjnego (wzorzec wspólny)

```
   ┌──────────────────────────────────────────────────────────────────────┐
   │ NAGŁÓWEK OKNA   nazwa okna · przynależność do modułu       [⋯]       │  auto
   ├──────────────────────────────────────────────────────────────────────┤
   │                                                                      │
   │ OBSZAR ZAWARTOŚCI    edytor · lista · monitor · podgląd              │  1fr
   │                      aktualizacja na żywo (kanał WebSocket)          │  overflow: auto
   │                                                                      │
   ├──────────────────────────────────────────────────────────────────────┤
   │ ELEMENTY KONFIGURACJI    każde ustawienie z objaśnieniem `[?]`       │  auto
   └──────────────────────────────────────────────────────────────────────┘
      STAN PROCESU SESJI — po stronie serwera, trwały:
      rozłączenie klienta NIE zamyka okna
```

Deklaracja siatki: `grid-template-rows: auto 1fr` (E2) albo `auto 1fr auto`
(E4, gdy okno ma pas konfiguracji), zawsze z `min-height: 0`.

### 9.4. Chat Window — pięć pasów

Makieta E2 realizuje Chat Window jako siatkę `grid-template-rows: auto auto 1fr auto auto`:

| # | Pas | Wysokość | Zawartość | Klasy |
|---:|---|---|---|---|
| 1 | Nagłówek | auto | pola kontekstu (moduł, model, sesja) + stan po prawej | `.dn-kropka`, `.dn-plakietka` |
| 2 | Kontekst | auto | przypięte źródła, dokumenty, zakres pamięci | `.dn-plakietka`, `.dn-btn--sm` |
| 3 | **Historia rozmowy** | 1fr, `overflow-y: auto` | wpisy trzech klas semantycznych | `.dn-wpis--czlowiek / --inteligencja / --system` |
| 4 | Monitor | auto | pasek postępu bieżącej operacji | `.dn-postep` |
| 5 | Dół | auto | pole promptu + przybornik akcji | `.dn-prompt`, `.dn-przybornik` |

**Wysokość pasa:** żeton `--dn-wym-pas-komunikacji` (320 px) jest **minimum**,
nie wartością stałą. Makieta E2 zapisuje `minmax(320px, 54%)` — pas rośnie
proporcjonalnie na wysokich ekranach, ale nigdy nie schodzi poniżej 320 px.

### 9.5. Dziewięciu nadawców, trzy klasy semantyczne

Rozstrzygnięcie L-P-28 (README pakietu): dziewięć ról nadawców nie dostaje
dziewięciu barw tła. Dostaje **trzy klasy semantyczne**, a rozróżnienie w obrębie
klasy niesie komplet: **ikona medalionu + etykieta nadawcy + plakietka roli**.

| Rola kontraktu | Klasa | Barwa krawędzi lewej | Medalion |
|---|---|---|---|
| `user` | `.dn-wpis--czlowiek` | `--dn-atrament` | ikona `uzytkownik` |
| `assistant` | `.dn-wpis--inteligencja` | `--dn-sygnal-wypelnienie` | ikona `agent` |
| `system` | `.dn-wpis--system` | `--dn-obrys-mocny` | ikona `info` |
| `tool` | `.dn-wpis--system` | `--dn-obrys-mocny` | ikona narzędzia + **plakietka roli = nazwa narzędzia** |

Modyfikator `.dn-wpis--pracuje` dokłada kropkę tętna przy nazwie nadawcy —
oznacza nadawcę aktywnie piszącego.

### 9.6. Kompozycja przy w4 — dwa okna komunikacji

Próg w4 (1600 px) uruchamia układ pary koordynator–wykonawca (makieta E4):
dwa Chat Window obok siebie, każdy z własnym kontekstem roli.

```
  ┌───────────────────────────────┬───────────────────────────────┐
  │  Coordinator Chat             │  Executor Chat (Executor 1)   │  58%
  │  plan etapów · podział pracy  │  wykonanie zadań z kolejki    │
  │  budowa promptów · kolejka    │  Subagent Network (do 15)     │
  ├───────────────────────────────┴───────────────────────────────┤
  │  Results Analyzer             │  kolejka kroków / monitor     │  42%
  │  (Executor 3 / Validator)     │  .dn-kolejka + .dn-krok       │
  └───────────────────────────────┴───────────────────────────────┘
   grid-template-columns: 1fr 1fr          grid-template-columns: 1.25fr 1fr
```

Asymetria dolnego wiersza (1,25fr : 1fr) jest zgodna z pokrętłem
`WARIANCJA_PROJEKTOWA = 4/10`: asymetria wyłącznie tam, gdzie niesie hierarchię
(analiza wyników ma większą masę niż monitor kolejki).

---

## 10. Kontrakt komponent ↔ żeton

### 10.1. Co komponent MOŻE

| # | Uprawnienie | Przykład z `komponenty.css` |
|---|---|---|
| M1 | Czytać żetony semantyczne (W2) | `.dn-karta { background: var(--dn-powierzchnia) }` |
| M2 | Czytać żetony niezależne od motywu (W3) | `.dn-btn { height: var(--dn-wym-kontrolka) }` |
| M3 | Czytać żetony ramy, jeśli jest komponentem ramy | `.dn-btn-ikona--na-ramie { color: var(--dn-rama-tekst-2) }` |
| M4 | Składać żetony arytmetycznie | `calc(100vw - var(--dn-od-8))` w `.dn-modal` |
| M5 | Definiować własne `@keyframes` | `dn-tetno`, `dn-wejscie`, `dn-obrot` |
| M6 | Deklarować stany atrybutami ARIA | `[aria-pressed='true']`, `[aria-busy='true']`, `[aria-selected='true']`, `[aria-current='page']`, `[aria-invalid='true']` |
| M7 | Wprowadzać własne wartości geometryczne bez odpowiednika w żetonach | `width: 12px` dla ikony wewnątrz plakietki; `8px` dla kropki radia |

### 10.2. Czego komponent NIE MOŻE

| # | Zakaz | Dlaczego | Konsekwencja złamania |
|---|---|---|---|
| N1 | Sięgać po prymitywy W1 | łamie oś motywu | komponent nie przełącza się z motywem |
| N2 | Wpisywać wartość barwy wprost | omija źródło prawdy | rozjazd z `kontrasty.json`, niemierzalna dostępność |
| N3 | Definiować własną media query barwną | motyw jest mechanizmem żetonów | dwa źródła prawdy motywu |
| N4 | Definiować własną obsługę `prefers-reduced-motion` | obsługa jest globalna (§14 żetonów) | podwójne wyłączenie albo pominięcie |
| N5 | Używać `disabled` | ADL-017 zero blokad | interfejs przestaje być klikalny |
| N6 | Komunikować stan samą barwą | WCAG + KIERUNEK.md | stan niewidoczny dla części Operatorów |
| N7 | Używać `z-index` jako liczby | rozsypuje skalę warstw | konflikt warstw nie do wyśledzenia |
| N8 | Używać gradientu jako tła przycisku, karty albo sekcji | L-P-09, zakres ilustracyjny | wprowadza drugą barwę do monochromu |
| N9 | Używać `!important` | kaskada jest jednokierunkowa | nadpisanie nie do cofnięcia bez kolejnego `!important` |

### 10.3. Zasada „stan nigdy samym kolorem" — realizacja

| Komponent stanu | Barwa | Nośnik drugi (obowiązkowy) |
|---|---|---|
| `.dn-plakietka--sukces / --ostrzezenie / --blad / --informacja` | tło + obrys + tekst | ikona 12×12 wewnątrz plakietki + etykieta tekstowa |
| `.dn-kropka--sukces / --ostrzezenie / --blad / --neutralna` | wypełnienie | etykieta obok kropki (kropka nie występuje samodzielnie) |
| `.dn-krok--pracuje / --poprawny / --bledy / --wstrzymany` | krawędź lewa | znak kroku (`.dn-krok-znak`) z ikoną + `.dn-krok-meta` z opisem |
| `.dn-toast--sukces / --ostrzezenie / --blad / --informacja` | krawędź lewa 2 px | ikona 16×16 + tytuł tekstowy |
| `.dn-wpis--czlowiek / --inteligencja / --system` | krawędź lewa | medalion z ikoną + nazwa nadawcy + plakietka roli |
| `.dn-btn[aria-invalid='true']` | obrys | atrybut ARIA odczytywany przez technologie wspomagające |

### 10.4. Realizacja zasady zero blokad w warstwie komponentowej

Nagłówek `komponenty.css` zapisuje ADL-017 wprost: *„żaden wariant nie odbiera
klikalności — niegotowość komunikuje się opisem albo komunikatem"*.

| Zamiast | Stosuje się | Klasa / atrybut |
|---|---|---|
| `disabled` na przycisku | komunikat po naciśnięciu (toast) | `.dn-toast--ostrzezenie` |
| `disabled` na przycisku | opis obok kontrolki | `.dn-pole-opis` |
| wyszarzenie pola | `readonly` z widocznym tłem `--dn-powierzchnia-2` | `.dn-pole-kontrolka[readonly]` |
| stan „w trakcie" jako blokada | `aria-busy="true"` + spinner w przycisku | `.dn-btn[aria-busy='true']::after` |
| ukrycie niedostępnej metody przez wyszarzenie | **nierenderowanie segmentu** | brak metody = mniej segmentów |

### 10.5. Wyjątki udokumentowane

| Miejsce | Odstępstwo | Uzasadnienie |
|---|---|---|
| `.dn-awatar--inteligencja` | używa `var(--dn-szary-0)` (prymityw W1) | tekst na gradiencie sygnałowym, identycznym w obu motywach; brak żetonu semantycznego dla tej roli |
| `.dn-plakietka > svg`, `.dn-krok-znak > svg` | `width: 12px` wprost | stopień pośredni między `--dn-wym-ikona-sm` (14) a brakiem żetonu 12 px |
| `.dn-radio::before` | `8px` wprost | wypełnienie radia = połowa `--dn-wym-check`; wyrażenie `calc()` byłoby mniej czytelne niż wartość |
| `.dn-awatar-stan` | kropka 8×8 px | większa niż `--dn-wym-kropka` (6 px), bo nosi obrys odcinający od awatara |

Wyjątki są zamknięte listą. Rozszerzenie listy wymaga wpisu w niniejszym
rozdziale — nie wolno dodać wartości wprost bez odnotowania.

---

## 11. Architektura ikon

### 11.1. Reguły zestawu

| Reguła | Wartość | Źródło |
|---|---|---|
| Siatka źródłowa | 24 × 24 | `manifest.json` → `zasady.siatka` |
| Grubość obrysu | **1,75** | `manifest.json` → `zasady.obrys` |
| Wypełnienie | `none` | `manifest.json` → `zasady.wypelnienie` |
| Barwa | `currentColor` | `manifest.json` → `zasady.barwa` |
| Zakończenia i łączenia | `round` | `manifest.json` → `zasady.zakonczenia`, `.laczenia` |
| Stopnie renderowania | 14 · 16 · 20 · 24 px | `manifest.json` → `zasady.renderowanie` |
| Biblioteka źródłowa | Lucide (ISC), obrys ujednolicony do 1,75 | `manifest.json` → `zrodlo.biblioteka` |
| Ikony własne | emblematy środowisk — ta sama siatka i kreska, jedna wypełniona kropka sygnału | `manifest.json` → `zrodlo.wlasne` |
| Liczba pozycji | **82** | `manifest.json` → `liczba-ikon`; zgodna z liczbą plików w `svg/` |

### 11.2. Dlaczego `currentColor` jest decyzją architektoniczną

Ikona nie zna żetonów. Dziedziczy barwę po elemencie nadrzędnym — a ten czyta
żeton semantyczny. Dzięki temu:

```
  .dn-btn-ikona          color: var(--dn-tekst-2)   ──┐
      └─ <svg stroke="currentColor">                 │  ikona jest szara
                                                     │
  .dn-btn-ikona:hover    color: var(--dn-tekst)     ──┤  ikona rozjaśnia się
      └─ <svg stroke="currentColor">                 │  BEZ własnej reguły
                                                     │
  .dn-btn-ikona--na-ramie color: var(--dn-rama-tekst-2)  ikona na ramie
      └─ <svg stroke="currentColor">                    dostosowuje się sama
```

Konsekwencja: **`komponenty.css` nie zawiera ani jednej reguły barwiącej ikonę**.
Zawiera wyłącznie reguły wymiarujące (`width`, `height`), bo wymiar nie jest
dziedziczony.

### 11.3. Nazewnictwo — polskie, opisowe

Nazwy plików i wpisów manifestu są polskie i opisują **pojęcie**, nie kształt:
`wykonawca` (nie „osoba z kluczem"), `srodowisko-talkin` (nie „dymek z kropką"),
`walidator` (nie „tarcza z ptaszkiem"). To celowe: literówka w nazwie opisowej
jest wychwytywalna przy czytaniu kodu, literówka w nazwie kształtu — nie.

Pełny wykaz 82 nazw znajduje się w KANON §4 i w `manifest.json`.

### 11.4. Manifest jako źródło prawdy

```
   manifest.json  (82 pozycje: nazwa · zrodlo · zastosowanie)
        │
        │  generuje / weryfikuje
        ▼
   ikony/zrodla-ikon.ts   ── typ unii nazw ──►  ikony.ts  ──►  komponenty
        │                                                       │
        │  literówka w nazwie                                   │
        ▼                                                       ▼
   BŁĄD KOMPILACJI  ◄─────────────────────────────────  nazwa spoza unii
```

HANDOFF §2.1 zapisuje wprost: *„Literówka ma dalej zatrzymywać kompilację."*
Mechanizm: `zrodla-ikon.ts` deklaruje typ unii złożony z 82 dosłownych nazw;
`ikony.ts` przyjmuje wyłącznie ten typ. Nazwa spoza wykazu nie jest błędem
w czasie działania (brakująca ikona), tylko **błędem kompilacji**.

### 11.5. Trzy grupy ikon i ich pochodzenie

| Grupa | Liczność | Pochodzenie | Przykłady |
|---|---:|---|---|
| Interfejs ogólny | większość zestawu | Lucide, obrys ujednolicony do 1,75 | `dom`, `menu`, `szukaj`, `zamknij`, `plus`, `ustawienia` |
| Pojęcia domenowe | mniejszość | Lucide, dobrane do pojęcia produktu | `agent`, `agenci`, `kolejka` (przez `warstwy`), `walidator`, `wykonawca`, `terminal` |
| Emblematy środowisk | 4 | **rysunek własny** na siatce zestawu | `srodowisko-talkin`, `srodowisko-workspace`, `srodowisko-codestudio`, `srodowisko-multitaskingai` |

Emblematy środowisk mają cechę wspólną wymuszoną przez kierunek: **dokładnie
jedną wypełnioną kropkę sygnału** (KIERUNEK.md 3.4).

### 11.6. Zasada „nie rysuj, jeśli istnieje"

KIERUNEK.md 3.6 i KANON §4: *„Ikon nie rysuje się ręcznie, jeżeli istnieją
w bibliotece. Dorysowuje się wyłącznie brakujące pojęcia domenowe."*

Procedura dodania ikony:

```
   potrzebna nowa ikona
        │
        ├─ czy istnieje w Lucide? ── TAK ──► pobierz, ujednolić obrys do 1,75,
        │                                     nadaj nazwę polską opisową,
        │                                     dopisz do manifestu
        │
        └─ NIE ──► czy to pojęcie DOMENOWE produktu? ── NIE ──► użyj najbliższej
                          │                                      istniejącej
                          └─ TAK ──► narysuj na siatce 24, obrys 1,75,
                                      zakończenia round, wypełnienie none,
                                      dopisz do manifestu z `zrodlo: "wlasne"`
```

---

## 12. Architektura ruchu

### 12.1. Cztery czasy, jeden ease

| Żeton | Wartość | Przeznaczenie | Przykład zastosowania |
|---|---:|---|---|
| `--dn-czas-1` | 0,10 s | mikroreakcje | najechanie na przycisk, naciśnięcie wiersza tabeli |
| `--dn-czas-2` | 0,16 s | przejścia barw | przełączenie motywu, zmiana tła pozycji nawigacji |
| `--dn-czas-3` | 0,22 s | wejście i wyjście warstw | modal, panel, toast |
| `--dn-czas-tetno` | **2,40 s** | jedyny ruch ciągły | tętno kropki sygnału |
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | **jedyna krzywa systemu** | wszystkie powyższe |

Jedna krzywa dla wszystkich czasów to decyzja spójnościowa: ruch ma być
rozpoznawalny jako „ten sam ruch o różnej długości", nie jako cztery różne
charaktery. Krzywa startuje szybko i zatrzymuje się miękko — zachowanie
instrumentu, nie sprężyny.

### 12.2. Trzy animacje nazwane

| `@keyframes` | Linia | Czas | Komponent | Co robi |
|---|---:|---|---|---|
| `dn-tetno` | 348 | `--dn-czas-tetno` (2,4 s), `infinite` | `.dn-kropka--tetno` | pulsowanie kropki — „tu biegnie praca" |
| `dn-wejscie` | 962 | `--dn-czas-3` (0,22 s), raz | `.dn-modal[open]` | wejście warstwy: `opacity` + przesunięcie |
| `dn-obrot` | 1103 | 0,8 s `linear`, `infinite` | `.dn-spinner`, `.dn-btn[aria-busy]::after` | obrót wskaźnika pracy |

`dn-obrot` jest jedyną animacją poza skalą `--dn-czas-*` — obrót ciągły
o krzywej `linear` nie mieści się w semantyce mikroprzejść. Wartość 0,8 s jest
zapisana w `komponenty.css` wprost i odnotowana jako rozbieżność wobec katalogu
v1.0 (0,7 s); obowiązuje wartość v2.0.

### 12.3. Wzorzec przejścia widoku

KANON §11 (wzorzec animacji, `INTENSYWNOSC_RUCHU = 3/10`):

| Reguła | Wartość |
|---|---|
| Zakres mikroprzejść | 100–220 ms |
| Właściwości animowane przy zmianie widoku | `opacity` + `translateY(4px)` |
| Maksymalny czas przejścia widoku | 220 ms (`--dn-czas-3`) |
| Liczba ruchów znaczących na widok | **jeden** — tętno pracy w tle |
| Zakazane | animacje dekoracyjne, parallax, scrollytelling w oknach roboczych |
| Warunek | ruch zawsze niesie informację o stanie systemu |

### 12.4. Globalna obsługa `prefers-reduced-motion`

Obsługa jest **w żetonach, nie w komponentach** (zetony.css §14, l. 432–448).
Realizuje dwa poziomy jednocześnie:

```
   @media (prefers-reduced-motion: reduce) {
     :root { --dn-czas-1 …4 → 0.01ms }        ← POZIOM 1: żetony
                                                 komponenty czytające żeton
                                                 tracą ruch automatycznie
     *, *::before, *::after {                  ← POZIOM 2: siatka bezpieczeństwa
       animation-duration: 0.01ms !important     wyłapuje wartości zapisane
       animation-iteration-count: 1 !important   wprost (np. 0.8s spinnera)
       transition-duration: 0.01ms !important
       scroll-behavior: auto !important
     }
   }
```

Poziom 2 jest jedynym miejscem w całym systemie, gdzie `!important` jest
dozwolony — bo musi wygrać z każdą regułą komponentu, także przyszłą.

**Zastępstwo dla tętna:** przy ograniczonym ruchu kropka nie miga, tylko
otrzymuje **pierścień statyczny** (L-P-16). `komponenty.css` realizuje to
regułą w bloku `@media (prefers-reduced-motion: reduce)` w sekcji kropki
(l. 353) — to jedyny komponent z własnym blokiem tej media query, i jest to
odstępstwo świadome: chodzi o **zamianę formy**, nie o wyłączenie ruchu.

### 12.5. Ruch a warstwy — co się animuje przy wejściu warstwy

| Warstwa | Animacja wejścia | Czas | Uwaga |
|---|---|---|---|
| Modal | `dn-wejscie` (opacity + przesunięcie) | `--dn-czas-3` | `::backdrop` z `blur(2px)` — jedyne rozmycie w systemie |
| Toast | wejście z prawej dolnej krawędzi | `--dn-czas-3` | pozycja `fixed`, prawy dolny róg |
| Tooltip | `opacity` + drobne przesunięcie | `--dn-czas-1` | wywołanie `:hover` / `:focus-within` |
| Panel boczny (zwijanie) | szerokość + `opacity` etykiet | `--dn-czas-2` | próg w2 |
| AOD | wejście rdzenia + rozwinięcie treści | `--dn-czas-3` | rdzeń pulsuje tętnem po wejściu |
| Zmiana modułu | `opacity` + `translateY(4px)` | `--dn-czas-3` | wzorzec przejścia widoku |

---

## 13. Bramy weryfikacyjne fal wdrożeniowych

### 13.1. Cztery fale (HANDOFF.md §3)

```
   FALA 1 — ŻETONY                FALA 2 — KOMPONENTY
   motyw/ : nowe wartości          komponenty/ + fundament
   + nowe pliki wymiary,           wzorzec odbioru:
     warstwy, rama                 06-okna/komponenty.html
   ────────────────────────        ────────────────────────
   aplikacja rusza na nowej        komponenty przechodzą
   palecie; komponenty stare       na nowe klasy .dn-*
   jeszcze działają
          │                                │
          ▼                                ▼
   FALA 3 — IKONY I MARKA         FALA 4 — OKNA
   ikony/, favicon,               okno-komunikacji + powłoka
   ikony Tauri                    wg makiet E1 / E2 / E4
   ────────────────────────       ────────────────────────
   podmiana całego katalogu       przepięcie --dc-* na .dn-wpis,
   svg/; manifest jako            .dn-prompt, .dn-postep
   źródło prawdy
```

### 13.2. Brama wspólna wszystkich fal

HANDOFF.md §3, zdanie zamykające:

> *„Brama weryfikacyjna każdej fali: porównanie z makietami w obu motywach oraz
> przebieg `kontrasty.json` (pomiar, nie deklaracja)."*

Dwa człony bramy są rozłączne i oba obowiązkowe:

| Człon | Co sprawdza | Jak |
|---|---|---|
| **Porównanie z makietami** | zgodność wizualna | zestawienie zrzutu z makietą źródłową, oba motywy, ta sama szerokość |
| **Przebieg `kontrasty.json`** | zgodność dostępnościowa | obliczenie kontrastu 33 par; wynik `ok: false` zatrzymuje falę |

### 13.3. Tabela bram per fala

| Fala | Zakres | Brama specyficzna | Kryterium przejścia |
|---|---|---|---|
| **1 — żetony** | `motyw/` (13 arkuszy + spinający) | wartości przeniesione **dosłownie** z `zetony.css`; żaden arkusz > 300 linii; kolejność importów zgodna z rozdz. 2.5 | wszystkie 33 pary `kontrasty.json` z `ok: true`; oba motywy renderują się na nowej palecie; brak wartości szesnastkowej poza `motyw/prymitywy.css` i `motyw/stany.css` |
| **2 — komponenty** | `komponenty/` (14 arkuszy) + `fundament.css` | galeria `06-okna/komponenty.html` jako wzorzec odbioru wizualnego; usunięty duplikat fokusu z dotychczasowego `motyw.css` | 119 klas `.dn-*` obecnych i renderujących się; zero wystąpień `disabled`; zero prymitywów W1 poza listą wyjątków (rozdz. 10.5); fokus widoczny na każdej kontrolce |
| **3 — ikony i marka** | `ikony/svg/` (podmiana całego katalogu), `manifest.json`, favicon, ikony Tauri | wykaz nazw w `zrodla-ikon.ts` zgodny z manifestem (82 pozycje) | kompilacja przechodzi; próba użycia nazwy spoza wykazu **zatrzymuje kompilację**; ikony dziedziczą barwę przez `currentColor` w obu motywach; godło zachowuje barwy własne |
| **4 — okna** | `okno-komunikacji/` + powłoka wg makiet E1/E2/E4 | usunięte lokalne `--dc-*`; mapowanie ról kontraktu na klasy `.dn-wpis--*` | pas komunikacji nie schodzi poniżej 320 px; Chat Window nie znika przy zmianie modułu; karty sesji zachowują stan; powłoka zachowuje się poprawnie na progach w1–w4 |

### 13.4. Czego nie wolno zmienić po cichu (HANDOFF §4)

| # | Element chroniony | Skutek zmiany bez decyzji |
|---:|---|---|
| 1 | Wartości żetonów (`zetony.css` / `zetony.json`) | rozjazd z `kontrasty.json` — dostępność przestaje być mierzalna |
| 2 | Geometria znaku i emblematów (krzywe, nie fonty) | znak przestaje być odtwarzalny w skali |
| 3 | Zasada jednego akcentu (sygnał ≤ 5% ekranu, nigdy tło sekcji) | monochrom przestaje być monochromem |
| 4 | Pierścień fokusu i zachowanie `prefers-reduced-motion` | naruszenie warunku wejściowego WCAG |
| 5 | Gęstość zwarta jako domyślna | zmiana decyzji Właściciela (L-P-13) |

### 13.5. Kolejność fal — dlaczego ta, a nie inna

```
   Dlaczego żetony pierwsze?
   ─────────────────────────
   Komponenty stare czytają nazwy żetonów, które w v2.0 istnieją nadal
   (--dn-tlo, --dn-tekst, --dn-obrys…). Po fali 1 aplikacja ma nową paletę,
   ale nie ma jeszcze nowej geometrii. To stan pośredni ZAMIERZONY: pozwala
   odebrać barwę osobno od formy.

   Dlaczego ikony po komponentach?
   ────────────────────────────────
   Ikona jest osadzana w komponencie i dziedziczy po nim barwę. Podmiana
   katalogu przed falą 2 dałaby nowe ikony w starych opakowaniach —
   niemożliwe do odebrania wizualnie.

   Dlaczego okna na końcu?
   ───────────────────────
   Okno komunikacji i powłoka używają WSZYSTKICH trzech poprzednich warstw
   naraz. Wcześniejsze wdrożenie oznaczałoby odbiór na niepełnym fundamencie.
```

---

## 14. Decyzje projektowe

Rozstrzygnięcia podjęte ponad dokumentację źródłową — wraz z uzasadnieniem.

| # | Zagadnienie | Co mówi dokumentacja | Rozstrzygnięcie | Uzasadnienie |
|---|---|---|---|---|
| **D-01** | Sekcja 7 `zetony.css` (siatka i punkty łamania) nie ma pliku docelowego | HANDOFF.md §1 wymienia arkusze dla sekcji 1, 2, 3, 4, 5, 6, 8, 9, 10–11, 12, 13, 14 — **pomija sekcję 7** | wprowadzono arkusz **`motyw/siatka.css`** | sekcja 7 jest samodzielną odpowiedzialnością (progi + siatka + maks. szerokość treści); doklejenie jej do `wymiary.css` złamałoby R1 |
| **D-02** | Liczba poziomów warstw | README pakietu (L-P-10) mówi „10 poziomów"; `zetony.css` §8 i KANON §2 wymieniają **11** | przyjęto **11 poziomów** | źródłem prawdy wartości jest `zetony.css` (HANDOFF §4.1); README podaje liczbę omyłkowo |
| **D-03** | Zawartość `stany.css` | HANDOFF wymienia nazwę pliku bez wskazania sekcji | do `stany.css` trafiają: prymitywy stanów (§1, blok STANY) + żetony stanów semantycznych z §10/§11 + odpowiadające bloki zapasowe | stany są jedyną rodziną barw przecinającą wszystkie trzy warstwy żetonów; trzymanie ich razem jest jedyną wersją zgodną z R1 |
| **D-04** | Cienie w `przestrzen.css` | HANDOFF: „`przestrzen.css` (odstępy+promienie+cienie)"; cienie leżą fizycznie w blokach motywów §10/§11 | cienie wydzielone z bloków motywów do `przestrzen.css`, z zachowaniem selektorów `:root[data-theme='light' \| 'dark']` | zgodnie z literą HANDOFF; selektor musi zostać, bo cienie są jedyną częścią `przestrzen.css` zależną od motywu |
| **D-05** | Kolejność importów w `motyw.css` | HANDOFF: „arkusz spinający importuje całość" — bez kolejności | ustalono kolejność: prymitywy → rama → typografia → przestrzeń → ruch → wymiary → siatka → warstwy → gradienty → semantyczne-jasny → semantyczne-ciemny → stany (rozdz. 2.5) | W1 przed W3 przed W2; `stany.css` ostatni, bo czyta prymitywy i musi nadpisać ewentualne wartości z bloków motywów |
| **D-06** | Kolejność importów w `komponenty.css` | brak wskazania | kolejność **zgodna z kolejnością sekcji w pliku źródłowym** (przycisk → pole → wybór → plakietka → karta → tabela → zakładki → pasek → boczna → wpis → postęp → nakładka → powiadomienie → drobne) | plik źródłowy jest kaskadą jednokierunkową; zmiana kolejności na alfabetyczną mogłaby odwrócić zależności `.dn-suwak` ↔ `.dn-przelacznik` |
| **D-07** | Zachowanie powłoki na progach w1 i w3 | dokumentacja rozstrzyga wyłącznie w2 („boczna zwija się do ikon") i w4 („dwa okna komunikacji") | w1: jedno okno naraz, boczna jako szuflada, Chat Window jako osobny widok. w3: boczna 224 px z etykietami, okna wspomagające w zakładkach | doprecyzowanie wyprowadzone z wymiarów (`--dn-wym-boczna` 224 px nie mieści się poniżej w3 obok okna wiodącego) oraz z inwentarza okien (E11 Mobile = widok w1) |
| **D-08** | Szerokość bocznej w stanie zwiniętym (w2) | brak wartości w dokumentacji | **48 px** | równa wysokości paska górnego (`--dn-wym-pasek`) — zwinięta boczna tworzy kwadrat 48×48 dla ikony 16 px z marginesem, i zachowuje rytm rogu kompozycji |
| **D-09** | Reguła `min-height: 0` | brak w dokumentacji jako reguły; obecna w makietach E2 i E4 na sześciu pojemnikach | podniesiona do rangi **reguły architektonicznej**: każdy pojemnik siatki między `100dvh` a elementem przewijanym deklaruje `min-height: 0` | bez niej powłoka o stałej wysokości pęka przy przyroście treści; to jedyna reguła układu, której brak nie daje błędu, tylko cichą awarię kompozycji |
| **D-10** | Model kompozycji okna operacyjnego | README pakietu odnotowuje pytanie otwarte M-4 („układ okien per moduł") | przyjęto układ makiet: **okno wiodące u góry, Chat Window na dole (min. 320 px), okna wspomagające jako zakładki (w3) albo druga kolumna (w4)** | makiety E2 i E4 są jedynym rozstrzygnięciem wizualnym w pakiecie; oznaczone jako układ dokumentacyjny do potwierdzenia decyzją Właściciela |
| **D-11** | Rozbieżności wartości między katalogiem v1.0 a żetonami v2.0 | katalog v1.0 podaje pasek 56 px, przycisk ikonowy 36 px, stopień bazowy 15 px; `zetony.css` v2.0 podaje 48 / 32 / 13 px | obowiązują **wartości v2.0** | KIERUNEK.md: pakiet v2.0 „zastępuje w całości" warstwę wizualną v1.0; wartości v1.0 pozostają wyłącznie jako zapis historyczny w inwentarzu |
| **D-12** | Punkty łamania w zapytaniach `@media` | żetony `--dn-bp-*` istnieją, ale zmienne CSS nie działają w warunkach `@media` | progi zapisuje się liczbą; żetony pełnią rolę **rejestru wartości** dla dokumentacji i skryptu; zgodność sprawdza przegląd kodu | ograniczenie języka, nie decyzja projektowa; odnotowane, by nie zostało odczytane jako niekonsekwencja |
| **D-13** | `!important` w obsłudze ograniczonego ruchu | zakaz `!important` w warstwie komponentowej vs. jego obecność w `zetony.css` §14 | `!important` dozwolony **wyłącznie** w bloku `@media (prefers-reduced-motion: reduce)` w `ruch.css` | siatka bezpieczeństwa musi wygrać z każdą regułą komponentu, także przyszłą i nieprzewidzianą |
| **D-14** | Odstępstwo `.dn-kropka` od zakazu N4 | zakaz media query ograniczonego ruchu w komponencie vs. `komponenty.css` l. 353 | odstępstwo utrzymane i odnotowane | blok nie **wyłącza** ruchu (to robi warstwa żetonów), tylko **zamienia formę** — tętno na pierścień statyczny; zamiana formy jest odpowiedzialnością komponentu |
| **D-15** | Wyjątki od zakazu wartości wprost | brak listy w dokumentacji | zamknięta lista czterech wyjątków (rozdz. 10.5) | wyjątek bez rejestru staje się precedensem; rejestr zamienia go w decyzję |

---

## Podsumowanie architektury w dziesięciu zdaniach

1. System ma cztery warstwy: prymitywy, semantykę per motyw, warstwę niezależną od motywu i warstwę komponentową — zależność biegnie tylko w górę.
2. Rama kokpitu jest jedynym elementem barwnym poza osią motywu i dlatego mieszka w warstwie niezależnej od motywu.
3. Pakiet design ma osiem katalogów; katalog docelowy `budowa/client/src` ma trzy: `motyw/`, `komponenty/`, `ikony/` — plus `zasoby/marka/`.
4. Regulamin budowy to trzy zasady: jedna odpowiedzialność na plik, arkusz ≤ 300 linii, komponent czyta wyłącznie semantykę.
5. Motyw działa na trzech elementach — `data-theme`, `prefers-color-scheme`, `color-scheme` — a powielenie wartości w bloku zapasowym jest świadome i objęte kontraktem utrzymaniowym.
6. Gęstość ma dwie osie: automatyczną (`pointer: coarse`) i jawną (`data-gestosc`), obie realizowane żetonem, nie wyjątkiem w komponencie.
7. Skala warstw ma jedenaście poziomów i dwa pasma rezerwowe, a Centrum poleceń stoi na szczycie, bo wywołuje się z każdego stanu.
8. Powłoka to trzypoziomowa siatka z regułą `min-height: 0` na każdym pojemniku pośrednim; Chat Window jest elementem powłoki, nie modułu.
9. Ikony dziedziczą barwę przez `currentColor`, a manifest jest źródłem prawdy, którego literówka zatrzymuje kompilację.
10. Wdrożenie biegnie czterema falami — żetony, komponenty, ikony i marka, okna — a bramą każdej z nich jest porównanie z makietami w obu motywach oraz przebieg pomiarów kontrastu.

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
