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
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | System wizualny |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant · deweloper · redaktor treści interfejsu |
| **Przeznaczenie** | Ustala pełny, kompletny wobec źródeł wykaz żetonów `--dn-*` i klas komponentów `.dn-*` systemu wizualnego platformy oraz mechanizm motywu jasnego i ciemnego, tak aby żadna wartość wizualna nie była wpisywana na sztywno w oknach platformy |
| **Zakres** | prymitywy barw, żetony ramy kokpitu, żetony semantyczne obu motywów, typografia, przestrzeń i promienie, cienie i ruch, wymiary, punkty łamania, warstwy, gradienty, mechanizm przełączania motywu i adaptacji, pełna biblioteka komponentów `.dn-*`, ikonografia, stany interakcji, zasada zero blokad, progi kontrastu |
| **Poza zakresem** | klasy ramy aplikacji specyficzne dla belki, szyny i wstążki (`design/zasoby/rama.css`) — [Rama okna aplikacji](rama-okna.md); klasy wyłącznie prototypowe (`design/zasoby/prototyp.css`, prefiks `.pt-*`) — przywoływane punktowo tam, gdzie okno platformowe z nich korzysta, nie katalogowane tu w pełni; katalog komponentów złożonych z wielu klas podstawowych — [Katalog komponentów](katalog-komponentow.md) |
| **Dokument nadrzędny** | [Elementy okien](elementy-okien.md) |
| **Dokumenty powiązane** | [Rama okna aplikacji](rama-okna.md) · [Katalog komponentów](katalog-komponentow.md) · [Strona główna i nawigacja](strona-glowna-i-nawigacja.md) · [Okno instalatora](okno-instalatora.md) · [Okno nowego projektu](okno-nowego-projektu.md) · [Okno historii sesji](okno-historii-sesji.md) · [Okno instrukcji](okno-instrukcji.md) |
| **Prototypy odniesienia** | `design/05-okna/` — komplet prototypów, każdy zbudowany wyłącznie na żetonach i klasach niniejszego dokumentu |
| **Źródła normatywne** | `design/zasoby/zetony/zetony.css` (309 żetonów `--dn-*`) · `design/zasoby/css/komponenty.css` (123 klasy `.dn-*` w 18 rodzinach) · `design/zasoby/ikony/manifest.json` (82 ikony) · `budowa/shared/kontrasty-progi.json` (33 pary zmierzone wobec progu WCAG 2.1 AA) |
| **Zasada nadrzędna** | Komponenty odwołują się wyłącznie do żetonów semantycznych, nigdy do prymitywów ani do wartości zaszytych na sztywno. Oba motywy — jasny i ciemny — są równoprawne i zdefiniowane osobno, nie wywodzone jeden z drugiego. Stan nie jest nigdy niesiony samym kolorem. Zero blokad: żaden komponent nie odbiera klikalności — niegotowość komunikuje się opisem albo komunikatem. |

---

## Spis treści

1. [Zasady systemu wizualnego](#1-zasady-systemu-wizualnego)
   - [1.1 Kierunek projektowy](#11-kierunek-projektowy)
   - [1.2 Zasady nienegocjowalne](#12-zasady-nienegocjowalne)
   - [1.3 Trzy warstwy żetonów](#13-trzy-warstwy-żetonów)
2. [Prymitywy barw](#2-prymitywy-barw)
   - [2.1 Skala szarości](#21-skala-szarości)
   - [2.2 Skala sygnału](#22-skala-sygnału)
   - [2.3 Prymitywy stanu](#23-prymitywy-stanu)
3. [Żetony ramy kokpitu — stałe w obu motywach](#3-żetony-ramy-kokpitu--stałe-w-obu-motywach)
4. [Żetony semantyczne — motyw jasny](#4-żetony-semantyczne--motyw-jasny)
5. [Żetony semantyczne — motyw ciemny](#5-żetony-semantyczne--motyw-ciemny)
6. [Typografia](#6-typografia)
   - [6.1 Kroje](#61-kroje)
   - [6.2 Stopnie, wagi, interlinia, liternictwo](#62-stopnie-wagi-interlinia-liternictwo)
7. [Przestrzeń i promienie](#7-przestrzeń-i-promienie)
   - [7.1 Jednostka 4 px](#71-jednostka-4-px)
   - [7.2 Promienie zaokrągleń](#72-promienie-zaokrągleń)
8. [Cienie i ruch](#8-cienie-i-ruch)
   - [8.1 Cienie](#81-cienie)
   - [8.2 Czas i easing](#82-czas-i-easing)
9. [Wymiary](#9-wymiary)
10. [Punkty łamania i siatka](#10-punkty-łamania-i-siatka)
11. [Warstwy (z-index)](#11-warstwy-z-index)
12. [Gradienty](#12-gradienty)
13. [Mechanizm motywu i adaptacji](#13-mechanizm-motywu-i-adaptacji)
   - [13.1 Wybór jawny `data-theme`](#131-wybór-jawny-data-theme)
   - [13.2 Domyślna preferencja systemu](#132-domyślna-preferencja-systemu)
   - [13.3 Dotyk — cele rosną żetonem](#133-dotyk--cele-rosną-żetonem)
   - [13.4 Gęstość przestronna](#134-gęstość-przestronna)
   - [13.5 Ograniczony ruch](#135-ograniczony-ruch)
14. [Biblioteka komponentów](#14-biblioteka-komponentów)
   - [14.1 Przycisk — pełny rozkład stanów](#141-przycisk--pełny-rozkład-stanów)
   - [14.2 Plakietka — role i stany](#142-plakietka--role-i-stany)
   - [14.3 Powiadomienie (dymek) — struktura i cztery warianty](#143-powiadomienie-dymek--struktura-i-cztery-warianty)
   - [14.4 Dymek objaśnienia — mechanizm ujawnienia](#144-dymek-objaśnienia--mechanizm-ujawnienia)
   - [14.5 Awatar, wirnik, stan pusty](#145-awatar-wirnik-stan-pusty)
   - [14.6 Wpis okna komunikacji (rodzina 11)](#146-wpis-okna-komunikacji-rodzina-11)
   - [14.7 Monitor wykonania (rodzina 12)](#147-monitor-wykonania-rodzina-12)
   - [14.8 Pole wpisywania polecenia (rodzina 17)](#148-pole-wpisywania-polecenia-rodzina-17)
   - [14.9 Pole formularza — pełny rozkład stanów](#149-pole-formularza--pełny-rozkład-stanów)
   - [14.10 Wybór: pole wyboru, przełącznik, radio, suwak](#1410-wybór-pole-wyboru-przełącznik-radio-suwak)
   - [14.11 Karta — interaktywność i karta środowiska](#1411-karta--interaktywność-i-karta-środowiska)
   - [14.12 Boczna nawigacja modułów — pozycja bieżąca](#1412-boczna-nawigacja-modułów--pozycja-bieżąca)
   - [14.13 Tabela — nagłówek przyklejony i wiersz interaktywny](#1413-tabela--nagłówek-przyklejony-i-wiersz-interaktywny)
   - [14.14 Zakładki — kreska dolna i wariant pigułkowy](#1414-zakładki--kreska-dolna-i-wariant-pigułkowy)
   - [14.15 Modal (dialog) — struktura, warianty szerokości i trzy weryfikacje niezależne](#1415-modal-dialog--struktura-warianty-szerokości-i-trzy-weryfikacje-niezależne)
15. [Ikonografia](#15-ikonografia)
16. [Stany interakcji i zasada zero blokad](#16-stany-interakcji-i-zasada-zero-blokad)
17. [Dostępność i kontrast](#17-dostępność-i-kontrast)
18. [Zastosowanie zweryfikowane w czterech oknach platformowych](#18-zastosowanie-zweryfikowane-w-czterech-oknach-platformowych)
   - [18.1 Żetony ramy kokpitu (rozdz. 3) w oknach wejściowych](#181-żetony-ramy-kokpitu-rozdz-3-w-oknach-wejściowych)
   - [18.2 Zasada zero blokad (rozdz. 16) w praktyce czterech okien](#182-zasada-zero-blokad-rozdz-16-w-praktyce-czterech-okien)
   - [18.3 Klasy `.pt-*` przywołane obok `.dn-*` w oknach wejściowych](#183-klasy-pt--przywołane-obok-dn--w-oknach-wejściowych)
   - [18.4 Progi łamania lokalne — czwarty potwierdzony wzorzec](#184-progi-łamania-lokalne--czwarty-potwierdzony-wzorzec)
19. [Kryteria odbioru](#19-kryteria-odbioru)
20. [Załącznik A. Pełna tabela żetonów `--dn-*`](#załącznik-a-pełna-tabela-żetonów---dn-)
21. [Załącznik B. Pełna tabela klas `.dn-*`](#załącznik-b-pełna-tabela-klas-dn-)
22. [Załącznik C. Pełny katalog ikon](#załącznik-c-pełny-katalog-ikon)

---

## 1. Zasady systemu wizualnego

### 1.1 Kierunek projektowy

System wizualny platformy jest zapisany w komentarzu nagłówkowym `design/zasoby/zetony/zetony.css`
jako **monochromatyczna precyzja** — dosłownie: „odcienie bieli (motyw jasny), odcienie czerni
(motyw ciemny), jeden błękit sygnałowy”. Nie jest to system dwubarwny ani wielobarwny: osiemnaście
kroków skali `--dn-szary-*` niesie niemal całą powierzchnię interfejsu, a jedynym akcentem
chromatycznym poza szarością jest rodzina `--dn-sygnal-*` — jeden błękit, w ośmiu odcieniach.
Trzy dodatkowe skale prymitywne (`--dn-zielen-*`, `--dn-bursztyn-*`, `--dn-czerwien-*`) istnieją
wyłącznie jako podstawa czterech stanów semantycznych (sukces, ostrzeżenie, błąd, informacja,
rozdz. 4–5) i nigdy nie występują jako akcent dekoracyjny poza tym zastosowaniem.

Plik źródłowy nosi w nagłówku dopisek „v2.0 (projekt od podstaw)” — system wizualny obowiązujący w
platformie jest przebudową całkowitą, nie ewolucją wcześniejszej palety. Każda wartość wizualna w
oknach platformy (rozdz. 14, [Katalog komponentów](katalog-komponentow.md)) pochodzi z tego jednego
pliku; żaden inny zbiór tokenów nie jest źródłem równoległym.

### 1.2 Zasady nienegocjowalne

Komentarz źródłowy `zetony.css` wylicza wprost cztery zasady, którym podlega cały system —
„nienegocjowalne”, dosłownie:

| Zasada | Znaczenie |
|---|---|
| Jednostka 4 px | cała skala przestrzeni (`--dn-od-*`, rozdz. 7.1) jest wielokrotnością 4 px; żadna wartość odstępu nie łamie tej siatki |
| Oba motywy równoprawne | jasny i ciemny są zdefiniowane osobno w kodzie źródłowym (rozdz. 4–5), nie jeden wywodzony algorytmicznie z drugiego — obie wersje mają ten sam komplet tokenów, zmierzony osobno |
| Stan nigdy samym kolorem | każdy komponent niosący stan (błąd, sukces, aktywność, wybór) łączy kolor z kształtem, ikoną albo tekstem — barwa nigdy nie jest jedynym nośnikiem znaczenia (rozdz. 16) |
| Zero blokad | żaden wariant komponentu (`design/zasoby/css/komponenty.css`, w. 7) nie odbiera klikalności — niegotowość komunikuje się opisem albo komunikatem, nigdy odjęciem możliwości działania (rozdz. 16) |

Do tego dochodzi wymóg piąty, zapisany osobno w nagłówku: **WCAG 2.1 AA na wejściu** — każda para
tekst/tło jest zmierzona (rozdz. 17, `budowa/shared/kontrasty-progi.json`), nie szacowana.

### 1.3 Trzy warstwy żetonów

`zetony.css` dzieli żetony na trzy warstwy jawnie nazwane w komentarzu nagłówkowym:

```
┌───────────────────────────────────────────────────────────────────────┐
│  WARSTWA 1 · PRYMITYWY                                                │
│  surowe skale barw — ZAKAZ użycia wprost w komponentach               │
│  --dn-szary-0…1000 · --dn-sygnal-100…800 · --dn-zielen-400/bursztyn/  │
│  czerwien-100…700                                                     │
├───────────────────────────────────────────────────────────────────────┤
│  WARSTWA 2 · SEMANTYCZNE                                              │
│  role per motyw, data-theme="light|dark", oba motywy równoprawne      │
│  --dn-tlo · --dn-tekst · --dn-sygnal · --dn-sukces-tekst · …          │
├───────────────────────────────────────────────────────────────────────┤
│  WARSTWA 3 · NIEZALEŻNE OD MOTYWU                                     │
│  typografia, przestrzeń, ruch, wymiary, warstwy z-index               │
│  --dn-fs-* · --dn-od-* · --dn-czas-* · --dn-wym-* · --dn-z-*          │
└───────────────────────────────────────────────────────────────────────┘
```

Komponenty (`.dn-*`, rozdz. 14) odwołują się **wyłącznie** do warstwy 2 i 3 — nigdy do warstwy 1.
Ten zakaz jest zapisany w nagłówku `komponenty.css` wprost: „Wyłącznie żetony semantyczne z
zetony.css. Zero wartości zaszytych.”

Konsekwencja praktyczna tego podziału trzywarstwowego dla dalszej budowy produktu: nowy komponent
albo nowe okno platformowe nigdy nie deklaruje własnej wartości koloru, wymiaru ani czasu przejścia
wprost w pikselach, procentach czy kodzie szesnastkowym — sięga po żeton warstwy 2 (jeśli wartość
zależy od motywu) albo warstwy 3 (jeśli nie zależy). Cztery okna platformowe zweryfikowane w tym
przejściu redakcyjnym potwierdzają to konsekwentnie: żaden z czterech dokumentów
([Okno instalatora](okno-instalatora.md), [Okno historii sesji](okno-historii-sesji.md),
[Okno nowego projektu](okno-nowego-projektu.md), [Okno instrukcji](okno-instrukcji.md)) nie
odnotowuje ani jednej wartości wpisanej na sztywno w arkuszu lokalnym właściwym temu oknu (`.in-*`,
`.hs-*`, `.np-*`, `.iu-*`) — jedynym miejscem, w którym te cztery arkusze lokalne wprowadzają
wartość liczbową bezpośrednio, są progi punktów łamania własnych (rozdz. 10, 18.4), które z
definicji nie mają odpowiednika w żadnej z trzech warstw żetonów.

---

## 2. Prymitywy barw

### 2.1 Skala szarości

Osiemnaście kroków, czysto neutralnych (równe RGB) — od bieli do niemal-czerni. Wartość
`#000000` jest zakazana wprost komentarzem źródłowym; krańcem skali jest `--dn-szary-1000`
(`#0A0A0A`).

| Żeton | Wartość | Uwaga |
|---|---|---|
| `--dn-szary-0` | `#FFFFFF` | wyłącznie powierzchnie kart i tekst na atramencie |
| `--dn-szary-25` | `#FAFAFA` | |
| `--dn-szary-50` | `#F4F4F4` | |
| `--dn-szary-100` | `#ECECEC` | |
| `--dn-szary-150` | `#E3E3E3` | |
| `--dn-szary-200` | `#D7D7D7` | |
| `--dn-szary-300` | `#C0C0C0` | |
| `--dn-szary-400` | `#9E9E9E` | |
| `--dn-szary-500` | `#7C7C7C` | |
| `--dn-szary-600` | `#616161` | |
| `--dn-szary-700` | `#4A4A4A` | |
| `--dn-szary-750` | `#3A3A3A` | |
| `--dn-szary-800` | `#2A2A2A` | |
| `--dn-szary-850` | `#212121` | |
| `--dn-szary-900` | `#181818` | |
| `--dn-szary-925` | `#131313` | |
| `--dn-szary-950` | `#0F0F0F` | |
| `--dn-szary-1000` | `#0A0A0A` | kraniec skali; `#000000` jest zakazane |

### 2.2 Skala sygnału

Osiem kroków błękitu sygnałowego — jedyna barwa akcentu systemu poza szarością.

| Żeton | Wartość | Uwaga |
|---|---|---|
| `--dn-sygnal-100` | `#EDF3FE` | |
| `--dn-sygnal-200` | `#C9DAFB` | |
| `--dn-sygnal-300` | `#8FB2F5` | |
| `--dn-sygnal-400` | `#5C8CEC` | |
| `--dn-sygnal-500` | `#3B6FE0` | bazowy — kropka sygnału, fokus |
| `--dn-sygnal-600` | `#2457C9` | tekst i wypełnienia na jasnym |
| `--dn-sygnal-700` | `#1D49AF` | |
| `--dn-sygnal-800` | `#173A8C` | |

### 2.3 Prymitywy stanu

Dwanaście żetonów, cztery kroki na trzy skale — podstawa trzech stanów semantycznych z czterema
progami niosącymi (rozdz. 4–5): sukces (zieleń), ostrzeżenie (bursztyn), błąd (czerwień). Rodzina
„informacja” nie ma własnej skali prymitywnej — czerpie z rodziny sygnału (rozdz. 2.2), zgodnie z
komentarzem źródłowym „informacja = rodzina sygnału”.

| Skala | 700 (najciemniejszy) | 400 | 200 | 100 (najjaśniejszy) |
|---|---|---|---|---|
| Zieleń | `--dn-zielen-700` `#1E7A46` | `--dn-zielen-400` `#4FBD85` | `--dn-zielen-200` `#BCE3CD` | `--dn-zielen-100` `#E7F5EE` |
| Bursztyn | `--dn-bursztyn-700` `#8A5B0C` | `--dn-bursztyn-400` `#E0A63C` | `--dn-bursztyn-200` `#EED9A7` | `--dn-bursztyn-100` `#FBF2DE` |
| Czerwień | `--dn-czerwien-700` `#C0362F` | `--dn-czerwien-400` `#EE7168` | `--dn-czerwien-200` `#F3C8C4` | `--dn-czerwien-100` `#FCEDEB` |

---

## 3. Żetony ramy kokpitu — stałe w obu motywach

Siedem żetonów, wyjątek jawny od reguły „oba motywy zdefiniowane osobno” (rozdz. 1.2): rama kokpitu
(belka, szyna nawigacji) niesie tę samą wartość niezależnie od `data-theme`, zgodnie z komentarzem
źródłowym „stała w obu motywach (kierunek systemu projektowego)”.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-rama` | `var(--dn-szary-925)` (`#131313`) | tło belki tytułowej i szyny nawigacji |
| `--dn-rama-grafit` | `var(--dn-szary-850)` (`#212121`) | powierzchnie na atramencie ramy |
| `--dn-rama-tekst` | `var(--dn-szary-100)` (`#ECECEC`) | tekst i ikony ramy |
| `--dn-rama-tekst-2` | `var(--dn-szary-400)` (`#9E9E9E`) | tekst drugorzędny ramy |
| `--dn-rama-hover` | `rgba(255, 255, 255, 0.08)` | wskazanie kursorem pozycji ramy |
| `--dn-rama-obrys` | `rgba(255, 255, 255, 0.10)` | kreska działowa wewnątrz ramy |
| `--dn-rama-obrys-mocny` | `rgba(255, 255, 255, 0.24)` | ramka pozycji szyny nawigacji |

Stałość tych siedmiu żetonów oznacza, że kolumna tożsamości okien wejściowych i belka/szyna ramy
aplikacji wyglądają identycznie w obu motywach — zjawisko zweryfikowane osobno w
[Oknie instalatora](okno-instalatora.md) rozdz. 3.2 i w [Oknie nowego projektu](okno-nowego-projektu.md)
rozdz. 6, gdzie kontrastuje to z żetonami `--dn-powierzchnia`/`--dn-panel`, które **zmieniają się**
z motywem (rozdz. 4–5).

---

## 4. Żetony semantyczne — motyw jasny

Blok `:root[data-theme='light']` niesie czterdzieści żetonów semantycznych, zdefiniowanych
niezależnie od motywu ciemnego (rozdz. 5) — nie jako jego odwrócenie algorytmiczne.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-tlo` | `var(--dn-szary-50)` `#F4F4F4` | tło głównego obszaru roboczego — „papier roboczy” |
| `--dn-powierzchnia` | `var(--dn-szary-0)` `#FFFFFF` | karty, okna |
| `--dn-powierzchnia-2` | `var(--dn-szary-100)` `#ECECEC` | wgłębienia, tory |
| `--dn-panel` | `var(--dn-szary-25)` `#FAFAFA` | panele boczne |
| `--dn-wstazka` | `var(--dn-szary-0)` `#FFFFFF` | belka nakładana paska edycji |
| `--dn-wstazka-taca` | `var(--dn-szary-150)` `#E3E3E3` | taca pod belką nakładaną |
| `--dn-hover` | `rgba(10, 10, 10, 0.05)` | wskazanie kursorem, powierzchnia neutralna |
| `--dn-wcisniecie` | `rgba(10, 10, 10, 0.09)` | wciśnięcie |
| `--dn-nakladka` | `rgba(10, 10, 10, 0.45)` | przyciemnienie pod modalem |
| `--dn-obrys` | `var(--dn-szary-150)` `#E3E3E3` | obrys standardowy |
| `--dn-obrys-subtelny` | `var(--dn-szary-100)` `#ECECEC` | obrys ledwo widoczny |
| `--dn-obrys-mocny` | `var(--dn-szary-300)` `#C0C0C0` | obrysy kontrolek |
| `--dn-tekst` | `var(--dn-szary-900)` `#181818` | tekst główny — 16,1∶1 na tle |
| `--dn-tekst-2` | `var(--dn-szary-600)` `#616161` | tekst drugorzędny — 6,4∶1 na tle |
| `--dn-tekst-3` | `var(--dn-szary-500)` `#7C7C7C` | wyłącznie metadane i tekst ≥ 18,66 px półgruby |
| `--dn-tekst-inv` | `var(--dn-szary-0)` `#FFFFFF` | tekst odwrócony na tle ciemnym |
| `--dn-atrament` | `var(--dn-szary-900)` `#181818` | działanie główne — inwersja atramentu |
| `--dn-atrament-hover` | `var(--dn-szary-1000)` `#0A0A0A` | wskazanie kursorem działania głównego |
| `--dn-atrament-tekst` | `var(--dn-szary-0)` `#FFFFFF` | tekst na działaniu głównym |
| `--dn-sygnal` | `var(--dn-sygnal-600)` `#2457C9` | tekst, odnośniki — 5,5∶1 |
| `--dn-sygnal-mocny` | `var(--dn-sygnal-700)` `#1D49AF` | sygnał wzmocniony |
| `--dn-sygnal-wypelnienie` | `var(--dn-sygnal-600)` `#2457C9` | wypełnienie sygnałowe — biel na nim 5,5∶1 |
| `--dn-sygnal-wypelnienie-hover` | `var(--dn-sygnal-700)` `#1D49AF` | wskazanie kursorem wypełnienia |
| `--dn-sygnal-tlo` | `var(--dn-sygnal-100)` `#EDF3FE` | tło sygnałowe delikatne |
| `--dn-sygnal-obrys` | `var(--dn-sygnal-200)` `#C9DAFB` | obrys sygnałowy |
| `--dn-kropka` | `var(--dn-sygnal-500)` `#3B6FE0` | kropka sygnału (tętno) |
| `--dn-fokus` | `var(--dn-sygnal-500)` `#3B6FE0` | pierścień ogniska klawiatury |
| `--dn-fokus-cien` | `rgba(59, 111, 224, 0.30)` | poświata pierścienia ogniska |
| `--dn-sukces-tekst` | `var(--dn-zielen-700)` | tekst stanu sukcesu |
| `--dn-sukces-tlo` | `var(--dn-zielen-100)` | tło stanu sukcesu |
| `--dn-sukces-obrys` | `var(--dn-zielen-200)` | obrys stanu sukcesu |
| `--dn-ostrzezenie-tekst` | `var(--dn-bursztyn-700)` | tekst stanu ostrzeżenia |
| `--dn-ostrzezenie-tlo` | `var(--dn-bursztyn-100)` | tło stanu ostrzeżenia |
| `--dn-ostrzezenie-obrys` | `var(--dn-bursztyn-200)` | obrys stanu ostrzeżenia |
| `--dn-blad-tekst` | `var(--dn-czerwien-700)` | tekst stanu błędu |
| `--dn-blad-tlo` | `var(--dn-czerwien-100)` | tło stanu błędu |
| `--dn-blad-obrys` | `var(--dn-czerwien-200)` | obrys stanu błędu |
| `--dn-informacja-tekst` | `var(--dn-sygnal-600)` | tekst stanu informacyjnego |
| `--dn-informacja-tlo` | `var(--dn-sygnal-100)` | tło stanu informacyjnego |
| `--dn-informacja-obrys` | `var(--dn-sygnal-200)` | obrys stanu informacyjnego |
| `--dn-cien-1` … `--dn-cien-sygnal` | rozdz. 8.1 | pięć poziomów cienia właściwych motywowi jasnemu |

**Zasada użycia koloru sygnałowego jako tekstu.** Wartość `--dn-sygnal` w motywie jasnym
(`--dn-sygnal-600`, `#2457C9`) jest krokiem ciemniejszym od bazowego `--dn-sygnal-500` — bazowy
błękit na bieli nie osiąga progu kontrastu 4,5∶1 dla drobnego tekstu (rozdz. 17), dlatego tekst i
odnośniki sygnałowe korzystają zawsze z tokenu `--dn-sygnal`, nigdy z prymitywu `--dn-sygnal-500`
wprost.

---

## 5. Żetony semantyczne — motyw ciemny

Blok `:root[data-theme='dark']` niesie ten sam komplet czterdziestu żetonów, zdefiniowany osobno —
wartości różne od motywu jasnego wszędzie tam, gdzie kontrast wymaga innego kroku skali, identyczne
tylko przypadkiem tam, gdzie oba motywy dzielą tę samą wartość (żadna nie jest w tym zestawie).

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-tlo` | `var(--dn-szary-950)` `#0F0F0F` | tło głównego obszaru roboczego |
| `--dn-powierzchnia` | `var(--dn-szary-900)` `#181818` | karty, okna |
| `--dn-powierzchnia-2` | `var(--dn-szary-925)` `#131313` | wgłębienia |
| `--dn-panel` | `var(--dn-szary-850)` `#212121` | panele uniesione |
| `--dn-wstazka` | `var(--dn-szary-800)` `#2A2A2A` | belka nakładana paska edycji |
| `--dn-wstazka-taca` | `var(--dn-szary-950)` `#0F0F0F` | taca pod belką nakładaną |
| `--dn-hover` | `rgba(255, 255, 255, 0.06)` | wskazanie kursorem |
| `--dn-wcisniecie` | `rgba(255, 255, 255, 0.10)` | wciśnięcie |
| `--dn-nakladka` | `rgba(0, 0, 0, 0.62)` | przyciemnienie pod modalem |
| `--dn-obrys` | `var(--dn-szary-800)` `#2A2A2A` | obrys standardowy |
| `--dn-obrys-subtelny` | `var(--dn-szary-850)` `#212121` | obrys ledwo widoczny |
| `--dn-obrys-mocny` | `var(--dn-szary-750)` `#3A3A3A` | obrysy kontrolek |
| `--dn-tekst` | `var(--dn-szary-100)` `#ECECEC` | tekst główny — 15,0∶1 na tle |
| `--dn-tekst-2` | `var(--dn-szary-400)` `#9E9E9E` | tekst drugorzędny — 6,8∶1 na tle |
| `--dn-tekst-3` | `var(--dn-szary-500)` `#7C7C7C` | wyłącznie metadane |
| `--dn-tekst-inv` | `var(--dn-szary-950)` `#0F0F0F` | tekst odwrócony na tle jasnym |
| `--dn-atrament` | `var(--dn-szary-100)` `#ECECEC` | działanie główne — tu: biel |
| `--dn-atrament-hover` | `var(--dn-szary-25)` `#FAFAFA` | wskazanie kursorem działania głównego |
| `--dn-atrament-tekst` | `var(--dn-szary-950)` `#0F0F0F` | tekst na działaniu głównym |
| `--dn-sygnal` | `var(--dn-sygnal-300)` `#8FB2F5` | tekst, odnośniki — 8,0∶1 |
| `--dn-sygnal-mocny` | `var(--dn-sygnal-200)` `#C9DAFB` | sygnał wzmocniony |
| `--dn-sygnal-wypelnienie` | `var(--dn-sygnal-500)` `#3B6FE0` | wypełnienie sygnałowe — biel na nim 4,6∶1 |
| `--dn-sygnal-wypelnienie-hover` | `var(--dn-sygnal-400)` `#5C8CEC` | wskazanie kursorem wypełnienia |
| `--dn-sygnal-tlo` | `rgba(59, 111, 224, 0.16)` | tło sygnałowe delikatne |
| `--dn-sygnal-obrys` | `rgba(92, 140, 236, 0.40)` | obrys sygnałowy |
| `--dn-kropka` | `var(--dn-sygnal-400)` `#5C8CEC` | kropka sygnału (tętno) |
| `--dn-fokus` | `var(--dn-sygnal-400)` `#5C8CEC` | pierścień ogniska klawiatury |
| `--dn-fokus-cien` | `rgba(92, 140, 236, 0.35)` | poświata pierścienia ogniska |
| `--dn-sukces-tekst` | `var(--dn-zielen-400)` | tekst stanu sukcesu |
| `--dn-sukces-tlo` | `rgba(79, 189, 133, 0.14)` | tło stanu sukcesu |
| `--dn-sukces-obrys` | `rgba(79, 189, 133, 0.34)` | obrys stanu sukcesu |
| `--dn-ostrzezenie-tekst` | `var(--dn-bursztyn-400)` | tekst stanu ostrzeżenia |
| `--dn-ostrzezenie-tlo` | `rgba(224, 166, 60, 0.14)` | tło stanu ostrzeżenia |
| `--dn-ostrzezenie-obrys` | `rgba(224, 166, 60, 0.34)` | obrys stanu ostrzeżenia |
| `--dn-blad-tekst` | `var(--dn-czerwien-400)` | tekst stanu błędu |
| `--dn-blad-tlo` | `rgba(238, 113, 104, 0.14)` | tło stanu błędu |
| `--dn-blad-obrys` | `rgba(238, 113, 104, 0.34)` | obrys stanu błędu |
| `--dn-informacja-tekst` | `var(--dn-sygnal-300)` | tekst stanu informacyjnego |
| `--dn-informacja-tlo` | `rgba(59, 111, 224, 0.16)` | tło stanu informacyjnego |
| `--dn-informacja-obrys` | `rgba(92, 140, 236, 0.40)` | obrys stanu informacyjnego |
| `--dn-cien-1` … `--dn-cien-sygnal` | rozdz. 8.1 | pięć poziomów cienia właściwych motywowi ciemnemu, głębsze niż w motywie jasnym |

Zamiary przeciwstawne motywom widoczne w tabelach 4 i 5: w motywie jasnym `--dn-atrament` jest
niemal-czernią (`--dn-szary-900`) z tekstem białym; w motywie ciemnym `--dn-atrament` jest
niemal-bielą (`--dn-szary-100`) z tekstem niemal-czarnym (`--dn-szary-950`) — działanie główne jest
zawsze inwersją tła otoczenia, nie stałym kolorem powtórzonym w obu motywach.

---

## 6. Typografia

### 6.1 Kroje

Trzy role typograficzne, trzy kroje — żaden krój nie jest używany poza swoją rolą przypisaną.

| Żeton | Kroje (w kolejności odwołania) | Rola |
|---|---|---|
| `--dn-ff-naglowek` | `'Space Grotesk', 'IBM Plex Sans', 'Segoe UI', system-ui, sans-serif` | tytuły okien, nagłówki sekcji |
| `--dn-ff-bazowa` | `'IBM Plex Sans', 'Segoe UI', system-ui, -apple-system, sans-serif` | tekst interfejsu, treść bieżąca |
| `--dn-ff-mono` | `'IBM Plex Mono', 'Cascadia Mono', 'Consolas', monospace` | dane liczbowe, kod, etykiety wersalikowe, znaczniki czasu |

### 6.2 Stopnie, wagi, interlinia, liternictwo

Dziewięć stopni, gęstość zwarta, baza 14 px:

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-fs-xs` | `12px` | etykiety wersalikowe, nagłówki kolumn |
| `--dn-fs-sm` | `13px` | metadane, tekst pomocniczy |
| `--dn-fs-base` | `14px` | bazowy — tekst interfejsu |
| `--dn-fs-md` | `15px` | wyróżniona treść, wpisy rozmowy |
| `--dn-fs-lg` | `17px` | nagłówki paneli |
| `--dn-fs-xl` | `21px` | nagłówki okien i modali |
| `--dn-fs-2xl` | `24px` | tytuły sekcji strony głównej |
| `--dn-fs-3xl` | `30px` | tytuły kart środowisk |
| `--dn-fs-display` | `40px` | największy stopień ekspozycyjny |

Cztery wagi, trzy interlinie, trzy liternictwa:

| Rodzina | Żetony | Wartości |
|---|---|---|
| Waga | `--dn-fw-normalna` · `--dn-fw-srednia` · `--dn-fw-polgruba` · `--dn-fw-gruba` | `400` · `500` · `600` · `700` |
| Interlinia | `--dn-lh-ciasny` (nagłówki) · `--dn-lh-bazowy` (interfejs, gęstość zwarta) · `--dn-lh-luzny` (treść ciągła) | `1.25` · `1.45` · `1.6` |
| Liternictwo | `--dn-ls-naglowek` (Space Grotesk w dużych stopniach) · `--dn-ls-wersaliki` (etykiety Plex Sans) · `--dn-ls-mono-wersaliki` (etykiety Plex Mono) | `-0.01em` · `0.08em` · `0.14em` |

**Trzy kroje w praktyce, zweryfikowane w czterech oknach platformowych.** `--dn-ff-naglowek`
(Space Grotesk) nosi tytuł „Instalator aplikacji” w [Oknie instalatora](okno-instalatora.md) i tytuł
płótna „Nowy projekt” w [Oknie nowego projektu](okno-nowego-projektu.md); `--dn-ff-mono` (IBM Plex
Mono) nosi każdy nadtytuł wersalikowy zweryfikowany w tych dokumentach — „Instalacja · krok 3 z 6”,
„WorkSpace · Workspace środowiska”, „Rozdział 1” w [Oknie instrukcji](okno-instrukcji.md) — oraz
każdą wartość liczbową wyrównaną w kolumnie: rozmiary plików (rozdz. 5,
[Okno instalatora](okno-instalatora.md)), znaczniki czasu („12:04”, „wczoraj 16:08”, rozdz. 4.4
[Okna historii sesji](okno-historii-sesji.md)). `--dn-ff-bazowa` (IBM Plex Sans) nosi wszystko
pozostałe — opisy, etykiety pól, treść przycisków — jako krój dominujący ilościowo, mimo że nie jest
przywoływany z nazwy równie często jak dwa pozostałe w tekście tych dokumentów, właśnie dlatego że
jest wartością domyślną nienazywaną przy każdym akapicie.

---

## 7. Przestrzeń i promienie

### 7.1 Jednostka 4 px

Jedenaście kroków przestrzeni, każdy wielokrotnością 4 px — zasada nienegocjowalna (rozdz. 1.2).

| Żeton | Wartość |
|---|---|
| `--dn-od-0` | `0` |
| `--dn-od-1` | `4px` |
| `--dn-od-2` | `8px` |
| `--dn-od-3` | `12px` |
| `--dn-od-4` | `16px` |
| `--dn-od-5` | `20px` |
| `--dn-od-6` | `24px` |
| `--dn-od-8` | `32px` |
| `--dn-od-10` | `40px` |
| `--dn-od-12` | `48px` |
| `--dn-od-16` | `64px` |

Dwa żetony pochodne, rytm zestawów paneli i sekcji: `--dn-odstep-panel` (`var(--dn-od-3)`, 12 px) i
`--dn-odstep-sekcji` (`var(--dn-od-6)`, 24 px).

### 7.2 Promienie zaokrągleń

Sześć promieni, jeden system narożników opisany w komentarzu źródłowym jako „precyzja instrumentu
(małe promienie)”.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-r-xs` | `3px` | plakietki drobne, bloki kodu w wierszu |
| `--dn-r-sm` | `6px` | przyciski, pola, kontrolki |
| `--dn-r-md` | `8px` | karty pojedyncze, dymki |
| `--dn-r-lg` | `10px` | karty, modale, panele |
| `--dn-r-xl` | `14px` | karty środowisk, duże powierzchnie |
| `--dn-r-pill` | `999px` | plakietki stanu, przełączniki, kropki |

---

## 8. Cienie i ruch

### 8.1 Cienie

Pięć poziomów cienia na motyw, dziesięć żetonów łącznie — „tonowane neutralnie, jeden kierunek
światła (z góry)” w motywie jasnym; „głębsze; przełączają się wraz z motywem” w motywie ciemnym.

| Żeton | Motyw jasny | Motyw ciemny |
|---|---|---|
| `--dn-cien-1` | `0 1px 2px rgba(10,10,10,.05)` | `0 1px 2px rgba(0,0,0,.35)` |
| `--dn-cien-2` | `0 2px 8px rgba(10,10,10,.07), 0 1px 2px rgba(10,10,10,.05)` | `0 3px 10px rgba(0,0,0,.42), 0 1px 2px rgba(0,0,0,.30)` |
| `--dn-cien-3` | `0 8px 24px rgba(10,10,10,.12)` | `0 10px 30px rgba(0,0,0,.52)` |
| `--dn-cien-lg` | `0 20px 48px rgba(10,10,10,.18)` | `0 24px 60px rgba(0,0,0,.65)` |
| `--dn-cien-sygnal` | `0 0 0 3px rgba(59,111,224,.30)`* | `0 0 0 3px rgba(92,140,236,.22)` |

*Wartość `--dn-cien-sygnal` w tabeli powyżej pokrywa się z `--dn-fokus-cien` motywu jasnego
(rozdz. 4) — komentarz źródłowy zastrzega ten żeton wyłącznie do dwóch zastosowań: „pierścień pola
aktywnego i poświata kropki”, nie do ogólnego użycia dekoracyjnego.

### 8.2 Czas i easing

Jedna krzywa łagodzenia, cztery czasy przejścia.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | krzywa wspólna wszystkim przejściom |
| `--dn-czas-1` | `0.1s` | mikroreakcje: wskazanie kursorem, naciśnięcie |
| `--dn-czas-2` | `0.16s` | przejścia barw, przełączenie motywu |
| `--dn-czas-3` | `0.22s` | wejście/wyjście warstw: modal, panel, dymek powiadomienia |
| `--dn-czas-tetno` | `2.4s` | tętno kropki sygnału — jedyny ruch ciągły systemu |

„Jedyny ruch ciągły” oznacza dosłownie: poza tętnem kropki sygnału (zweryfikowanym w
[Oknie instalatora](okno-instalatora.md) rozdz. 3.2, klasa `.pt-tetno`, i w
[Oknie historii sesji](okno-historii-sesji.md) rozdz. 3.4) żaden element systemu nie animuje się
w pętli — cztery pozostałe żetony czasu opisują wyłącznie przejścia jednorazowe, wyzwalane
zdarzeniem.

---

## 9. Wymiary

Dwadzieścia osiem żetonów wymiaru, gęstość zwarta („kokpit”) jako stan domyślny — rozdz. 13.3–13.4
opisują dwa warianty adaptacji tej gęstości.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-wym-kontrolka` | `32px` | przyciski, pola, pozycje wyboru |
| `--dn-wym-ikonowy` | `34px` | przycisk ikonowy |
| `--dn-wym-belka` | `48px` | belka tytułowa okna — jej wysokość wyznacza szerokość szyny nawigacji |
| `--dn-wym-stan` | `28px` | pasek stanu (dolny) |
| `--dn-wym-pasek` | `48px` | pasek górny (rama) |
| `--dn-wym-pas-kart` | `36px` | pas kart sesji |
| `--dn-wym-wiersz` | `36px` | wiersz tabeli / listy |
| `--dn-wym-boczna` | `224px` | boczna nawigacja modułów |
| `--dn-wym-pas-komunikacji` | `320px` | wysokość pasa komunikacji (powłoka środowiska) |
| `--dn-wym-modal` | `560px` | maksymalna szerokość modala |
| `--dn-wym-toast-min` | `280px` | szerokość minimalna dymka powiadomienia |
| `--dn-wym-toast-max` | `420px` | szerokość maksymalna dymka powiadomienia |
| `--dn-wym-awatar-sm` | `24px` | awatar mały |
| `--dn-wym-awatar` | `28px` | awatar bazowy |
| `--dn-wym-awatar-lg` | `36px` | awatar duży |
| `--dn-wym-ikona-sm` | `16px` | ikona mała |
| `--dn-wym-ikona` | `18px` | ikona bazowa |
| `--dn-wym-ikona-szyna` | `20px` | ikona pozycji szyny nawigacji (szersza od belki) |
| `--dn-wym-ikona-lg` | `22px` | ikona duża |
| `--dn-wym-ikona-xl` | `26px` | ikona bardzo duża |
| `--dn-wym-przelacznik-szer` | `36px` | szerokość przełącznika |
| `--dn-wym-przelacznik-wys` | `20px` | wysokość przełącznika |
| `--dn-wym-check` | `16px` | pole wyboru |
| `--dn-wym-kropka` | `6px` | kropka sygnału |
| `--dn-wym-wstega` | `2px` | wstęga aktywności karty |
| `--dn-wym-spinner` | `14px` | wirnik ładowania |
| `--dn-wym-fokus` | `2px` | grubość pierścienia fokusu |
| `--dn-wym-fokus-odsuniecie` | `2px` | odsunięcie pierścienia fokusu |

Wartość `--dn-wym-belka` (48 px) jest przywołana dwukrotnie w tabeli — raz jako belka tytułowa, raz
jako `--dn-wym-pasek` (pasek górny ramy) — obie nazwy wskazują tę samą wartość liczbową z dwóch
różnych miejsc arkusza, celowo: wysokość belki wyznacza wprost szerokość szyny nawigacji pionowej
(kwadrat w narożniku, [Rama okna aplikacji](rama-okna.md)).

**Weryfikacja krzyżowa w oknie instalatora.** [Okno instalatora](okno-instalatora.md) rozdz. 3.4
zweryfikowało `--dn-wym-belka` niezależnie: pasek tytułu własny tego okna przedaplikacyjnego
(`.in-belka`) przejmuje dosłownie tę samą wysokość 48px, mimo że instalator nie niesie ramy aplikacji
w rozumieniu [Ramy okna](rama-okna.md) — okno bezramowe wybiera świadomie wysokość belki tożsamą z
ramą pełną, tak by pierwsze okno, jakie widzi Operator, nie skakało wysokością względem okien
kolejnych po zainstalowaniu produktu. Jest to jedyny żeton wymiaru odnaleziony jak dotąd w dwóch
zupełnie różnych klasach okien (okno przedaplikacyjne bez ramy i pełna rama aplikacji) o identycznej
wartości nieprzypadkowej.

---

## 10. Punkty łamania i siatka

Cztery progi nazwane, jeden próg treści dokumentowej, dwa żetony siatki dwunastokolumnowej.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-bp-w1` | `640px` | telefon poziomo — widok mobilny |
| `--dn-bp-w2` | `960px` | tablet — boczna nawigacja zwijana do ikon |
| `--dn-bp-w3` | `1280px` | biurko — pełny kokpit |
| `--dn-bp-w4` | `1600px` | szerokie biurko — dwa okna komunikacji (para koordynator–wykonawca) |
| `--dn-tresc-max` | `1200px` | maksymalna szerokość treści dokumentowej |
| `--dn-siatka-kolumny` | `12` | liczba kolumn siatki |
| `--dn-siatka-przerwa` | `var(--dn-od-6)` (24 px) | przerwa między kolumnami siatki |

**Progi lokalne okien nie zawsze pokrywają się z tymi czterema.** Weryfikacja w czterech oknach
platformowych tego przejścia redakcyjnego wykazała progi własne nienazwane żadnym `--dn-bp-*`:
1000 px i 640 px w [Oknie instalatora](okno-instalatora.md) (rozdz. 3.4), 900 px w
[Oknie historii sesji](okno-historii-sesji.md) (rozdz. 3.5) i w [Oknie instrukcji](okno-instrukcji.md)
(rozdz. 17), 1240 px w [Oknie nowego projektu](okno-nowego-projektu.md) (rozdz. 10). Żaden z tych
czterech progów lokalnych nie jest literalnie równy `--dn-bp-w1`…`--dn-bp-w4` — wzorzec powtarzający
się w całym zbiorze okien platformowych: arkusze lokalne `.hs-*`, `.in-*`, `.np-*`, `.iu-*` ustalają
własne progi wprost w pikselach, niezależnie od skali nazwanej tego rozdziału.

---

## 11. Warstwy (z-index)

Jedenaście poziomów, od podłogi po Command Center — zawsze najwyżej.

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-z-podloga` | `0` | poziom bazowy |
| `--dn-z-przybornik` | `10` | przyklejone nagłówki tabel, przybornik |
| `--dn-z-pasek` | `100` | pasek górny (rama) |
| `--dn-z-boczna` | `200` | boczna nawigacja |
| `--dn-z-pas-komunikacji` | `300` | pas komunikacji |
| `--dn-z-nakladka` | `800` | przyciemnienie pod modalem |
| `--dn-z-modal` | `900` | okno nakładkowe |
| `--dn-z-powiadomienie` | `1000` | dymki powiadomień |
| `--dn-z-tooltip` | `1100` | dymki podpowiedzi |
| `--dn-z-aod` | `1200` | Always On Display |
| `--dn-z-centrum-polecen` | `1300` | Command Center — zawsze najwyżej |

Jedenaście wartości powyżej, po jednej na żeton — kolejność rosnąca odpowiada kolejności nakładania
warstw wizualnych: rama i boczna nawigacja stoją nad treścią bazową, modal i jego przyciemnienie
stoją nad ramą, powiadomienia i podpowiedzi stoją nad modalem, a Always On Display i Command Center
stoją nad wszystkim innym, w tej kolejności.

Trzy z czterech okien platformowych zweryfikowanych w tym przejściu redakcyjnym są, w kategoriach
tej skali, oknami warstwy `--dn-z-modal` (900): [Okno historii sesji](okno-historii-sesji.md) i
[Okno instrukcji](okno-instrukcji.md) wprost jako `.dn-modal`, oraz pośrednio
[Okno instalatora](okno-instalatora.md), które nie korzysta z tej warstwy w ogóle — działa przed
uruchomieniem powłoki aplikacji, poza drzewem DOM, w którym te warstwy w ogóle mają znaczenie.
[Okno nowego projektu](okno-nowego-projektu.md), po korekcie względem redakcji poprzedniej
(rozdz. 1.1 tamtego dokumentu), nie jest modalem — stoi w warstwie treści bazowej przestrzeni
roboczej, nie w warstwie 900.

---

## 12. Gradienty

Dwa gradienty, użycie **wyłącznie ilustracyjne** — zastrzeżenie zapisane wprost w komentarzu
źródłowym.

| Żeton | Wartość | Zastosowanie dozwolone |
|---|---|---|
| `--dn-grad-atrament` | `linear-gradient(180deg, var(--dn-szary-800), var(--dn-szary-925))` | awatary bez zdjęcia, rdzeń Always On Display, grafiki brandowe |
| `--dn-grad-sygnal` | `linear-gradient(135deg, var(--dn-sygnal-400), var(--dn-sygnal-600))` | jak wyżej |

Zastrzeżenie dosłowne z komentarza źródłowego: „Gradient NIGDY nie jest tłem przycisku, karty ani
sekcji.” Dwa gradienty istnieją dla trzech zastosowań nazwanych wprost — nie są ogólną paletą
dekoracyjną do swobodnego użycia w nowych komponentach.

Jedynym zastosowaniem zweryfikowanym wprost w `komponenty.css` jest `.dn-awatar` (rozdz. 14.5): tło
bazowe `--dn-grad-atrament`, wariant jednostki inteligencji `--dn-grad-sygnal`. Żaden z czterech
dokumentów okien platformowych zweryfikowanych w tym przejściu redakcyjnym nie zawiera zrzutu
komponentu awatara wypełnionego gradientem — zastosowanie pozostaje potwierdzone wyłącznie w
źródle CSS, bez ilustracji w prototypie żadnego z czterech okien tej fali.

---

## 13. Mechanizm motywu i adaptacji

### 13.1 Wybór jawny `data-theme`

Atrybut `data-theme="light"` albo `data-theme="dark"` na elemencie `:root` wybiera jawnie jeden z
dwóch bloków semantycznych opisanych w rozdz. 4–5. Deklaracja `color-scheme: light` /
`color-scheme: dark` towarzyszy każdemu blokowi — przekazuje wybór motywu również elementom
natywnym przeglądarki (pola, suwaki), poza zasięgiem żetonów `--dn-*`.

Wybór motywu jest ustawieniem Operatora, nie właściwością pojedynczego okna — żadne z czterech okien
platformowych zweryfikowanych w tym przejściu redakcyjnym nie ustala własnego motywu niezależnie od
reszty aplikacji; wszystkie cztery opisują zachowanie „w obu motywach” jako właściwość swojego
układu, nie jako wybór lokalny.

### 13.2 Domyślna preferencja systemu

Bez jawnego atrybutu `data-theme`, dwa bloki `@media (prefers-color-scheme: light|dark)` targetujące
`:root:not([data-theme])` powtarzają — wartość po wartości — komplet czterdziestu żetonów
semantycznych z rozdz. 4 albo rozdz. 5. Komentarz źródłowy nazywa to powielenie świadomym:
„Wartości muszą być identyczne z blokami motywów; powielenie jest świadome (mechanizm kaskady, nie
drugie źródło prawdy)”. Innymi słowy: nie są to dwa źródła równoległe grożące rozjazdem wartości,
lecz mechanizm kaskady CSS wymagający powtórzenia tych samych wartości w bloku selektorowym o innej
specyficzności — redakcja przyszła zmieniająca wartość semantyczną musi zmienić ją w obu miejscach
jednocześnie (bloku `[data-theme]` i bloku `@media`), inaczej dwa mechanizmy wyboru motywu
rozjadą się w praktyce mimo wspólnego zamysłu.

### 13.3 Dotyk — cele rosną żetonem

`@media (pointer: coarse)` (urządzenie dotykowe) podnosi sześć żetonów wymiaru — komentarz źródłowy:
„cele dotykowe rosną żetonem, nie wyjątkiem”, czyli bez osobnej ścieżki kodu w komponentach, wyłącznie
przez nadpisanie wartości tokenu.

| Żeton | Wartość bazowa (rozdz. 9) | Wartość dotykowa |
|---|---|---|
| `--dn-wym-kontrolka` | 32px | 40px |
| `--dn-wym-ikonowy` | 34px | 40px |
| `--dn-wym-wiersz` | 36px | 44px |
| `--dn-wym-check` | 16px | 20px |
| `--dn-wym-przelacznik-szer` | 36px | 44px |
| `--dn-wym-przelacznik-wys` | 20px | 24px |

### 13.4 Gęstość przestronna

Atrybut `data-gestosc="przestronna"` na `:root` — komentarz źródłowy: „przygotowana, domyślnie
nieaktywna” — powiększa siedem żetonów wobec gęstości zwartej domyślnej (rozdz. 9).

| Żeton | Wartość zwarta (domyślna) | Wartość przestronna |
|---|---|---|
| `--dn-fs-base` | 14px | 14px (bez zmiany) |
| `--dn-lh-bazowy` | 1.45 | 1.5 |
| `--dn-wym-kontrolka` | 32px | 40px |
| `--dn-wym-ikonowy` | 34px | 40px |
| `--dn-wym-wiersz` | 36px | 44px |
| `--dn-wym-pasek` | 48px | 56px |
| `--dn-wym-pas-kart` | 36px | 40px |
| `--dn-odstep-panel` | `var(--dn-od-3)` (12px) | `var(--dn-od-5)` (20px) |
| `--dn-odstep-sekcji` | `var(--dn-od-6)` (24px) | `var(--dn-od-8)` (32px) |

Wariant „przestronna” jest przygotowany w arkuszu, lecz nieaktywny domyślnie — żaden atrybut
`data-gestosc` nie jest ustawiany automatycznie przez system; aktywacja wymaga jawnego ustawienia
tego atrybutu, poza zakresem żetonów samych w sobie. **[DO DECYZJI OPERATORA]** — czy i gdzie okno
konfiguracji udostępnia Operatorowi przełącznik tej gęstości.

### 13.5 Ograniczony ruch

`@media (prefers-reduced-motion: reduce)` sprowadza cztery żetony czasu (rozdz. 8.2) do `0.01ms` i
dodaje regułę globalną `*, *::before, *::after` wymuszającą `animation-duration: 0.01ms !important`,
`animation-iteration-count: 1 !important`, `transition-duration: 0.01ms !important` oraz
`scroll-behavior: auto !important`. Komentarz źródłowy: „globalnie, bez konfiguracji per komponent”
— żaden komponent nie potrzebuje własnej reguły dla ruchu ograniczonego, bo reguła globalna
zerowego czasu obejmuje każdą animację i przejście jednym blokiem.

---

## 14. Biblioteka komponentów

`design/zasoby/css/komponenty.css` niesie sto dwadzieścia trzy klasy `.dn-*` w osiemnastu rodzinach
funkcjonalnych. Nagłówek arkusza powtarza zasadę rozdz. 1.2 wprost: „Każdy komponent interaktywny ma
stany: spoczynek · wskazanie kursorem · naciśnięcie · fokus · wybrany · ładowanie · błąd. Zasada zero
blokad: żaden wariant nie odbiera klikalności — niegotowość komunikuje się opisem albo komunikatem.
Stan nigdy samym kolorem — komponent stanu niesie ikonę albo etykietę.”

| # | Rodzina | Liczba klas | Klasy główne |
|---:|---|---:|---|
| 1 | Przycisk | 11 | `.dn-btn`, warianty `--atrament` `--duch` `--niebezpieczny` `--sm` `--sygnal` `--wybrany` `--zarys`, `.dn-btn-ikona`, `.dn-btn-ikona--na-ramie`, `.dn-btn--lg` |
| 2 | Pole formularza | 6 | `.dn-pole`, `.dn-pole-blad`, `.dn-pole-etykieta`, `.dn-pole-kontrolka`, `.dn-pole-opis`, `.dn-szukaj` |
| 3 | Wybór: pole wyboru · przełącznik · radio | 5 | `.dn-check`, `.dn-przelacznik`, `.dn-radio`, `.dn-suwak`, `.dn-wybor` |
| 4 | Kropka sygnału | 6 | `.dn-kropka`, warianty `--blad` `--neutralna` `--ostrzezenie` `--sukces` `--tetno` |
| 5 | Plakietka | 7 | `.dn-plakietka`, warianty `--blad` `--informacja` `--ostrzezenie` `--rola` `--sukces` `--sygnal` |
| 6 | Karta | 15 | `.dn-karta`, `.dn-karta--klikalna`, `.dn-karta--wybrana`, `.dn-karta-cialo`, `.dn-karta-naglowek`, `.dn-karta-tytul`, `.dn-karta-srodowiska` i cztery elementy pochodne, `.dn-kafel` i trzy elementy pochodne |
| 7 | Tabela | 2 | `.dn-dane`, `.dn-tabela` |
| 8 | Zakładki · karty sesji | 5 | `.dn-karta-sesji`, `.dn-karty-sesji`, `.dn-zakladka`, `.dn-zakladki`, `.dn-zakladki--pigulki` |
| 9 | Rama kokpitu — pasek górny | 7 | `.dn-listwa`, `.dn-listwa-pozycja`, `.dn-pasek`, `.dn-pasek-godlo`, `.dn-pasek-logotyp`, `.dn-pasek-prawa`, `.dn-pasek-szukaj` |
| 10 | Boczna nawigacja modułów | 3 | `.dn-boczna`, `.dn-boczna-naglowek`, `.dn-boczna-pozycja` |
| 11 | Wpis okna komunikacji | 10 | `.dn-wpis`, warianty `--czlowiek` `--inteligencja` `--pracuje` `--system`, `.dn-wpis-godzina`, `.dn-wpis-medalion`, `.dn-wpis-nadawca`, `.dn-wpis-tozsamosc`, `.dn-wpis-tresc` |
| 12 | Monitor wykonania · kolejka kroków | 12 | `.dn-kolejka`, `.dn-krok` z wariantami `--bledy` `--poprawny` `--pracuje` `--wstrzymany`, `.dn-krok-meta`, `.dn-krok-znak`, `.dn-postep` i trzy elementy pochodne |
| 13 | Modal (dialog) · nakładka | 7 | `.dn-modal`, `.dn-modal--szeroki`, `.dn-modal-cialo`, `.dn-modal-naglowek`, `.dn-modal-stopka`, `.dn-modal-tytul`, `.dn-modal-zamknij` |
| 14 | Powiadomienie (dymek) | 8 | `.dn-toast`, warianty `--blad` `--informacja` `--ostrzezenie` `--sukces`, `.dn-toast-tresc`, `.dn-toast-tytul`, `.dn-toasty` |
| 15 | Dymek objaśnienia | 2 | `.dn-tooltip`, `.dn-tooltip-tresc` |
| 16 | Awatar · wirnik · stan pusty | 10 | `.dn-awatar` z wariantami `--inteligencja` `--kwadrat` `--lg` `--sm`, `.dn-awatar-stan`, `.dn-pusty-stan` i dwa elementy pochodne, `.dn-spinner` |
| 17 | Pole wpisywania polecenia | 4 | `.dn-prompt`, `.dn-prompt-grot`, `.dn-prompt-obszar`, `.dn-przybornik` |
| 18 | Always On Display | 3 | `.dn-aod`, `.dn-aod-rdzen`, `.dn-aod-tresc` |

Wykaz pełny, klasa po klasie, stoi w [Załączniku B](#załącznik-b-pełna-tabela-klas-dn-). Pięć
rodzin z pełnym rozbiciem stanów i wariantów poniżej (rozdz. 14.1–14.5), potem trzy rodziny warte
odnotowania w prozie (rozdz. 14.6–14.8).

### 14.1 Przycisk — pełny rozkład stanów

| Stan / wariant | Selektor | Realizacja |
|---|---|---|
| Spoczynek | `.dn-btn` | obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`, tekst `--dn-tekst` |
| Wskazanie kursorem | `.dn-btn:hover` | tło `--dn-hover` |
| Wciśnięcie | `.dn-btn:active` | `transform: translateY(1px)` |
| Ognisko | `.dn-btn:focus-visible` | pierścień `--dn-fokus`, grubość `--dn-wym-fokus`, odsunięcie `--dn-wym-fokus-odsuniecie` |
| Wybrany trwale | `.dn-btn[aria-pressed='true']`, `.dn-btn--wybrany` | tło `--dn-sygnal-tlo`, obrys `--dn-sygnal-obrys`, tekst `--dn-sygnal` |
| Ładowanie | `.dn-btn[aria-busy='true']` | kursor `progress`; wskaźnik `::after` — okrąg `--dn-wym-spinner`, obrót `dn-obrot` 0,8s liniowo; przycisk pozostaje klikalny (zero blokad) |
| Błąd | `.dn-btn[aria-invalid='true']` | obrys `--dn-blad-obrys`, tekst `--dn-blad-tekst` |
| Wariant główny | `.dn-btn--atrament` | tło `--dn-atrament`, obrys `--dn-atrament`, tekst `--dn-atrament-tekst`; wskazanie kursorem → `--dn-atrament-hover` |
| Wariant sygnałowy | `.dn-btn--sygnal` | tło `--dn-sygnal-wypelnienie`, tekst `--dn-szary-0`; zarezerwowany dla jednego działania systemowego na widok |
| Wariant zarys | `.dn-btn--zarys` | tło przezroczyste |
| Wariant duch | `.dn-btn--duch` | tło i obrys przezroczyste, tekst `--dn-tekst-2`; wskazanie kursorem → tło `--dn-hover`, tekst `--dn-tekst` |
| Wariant niebezpieczny | `.dn-btn--niebezpieczny` | tło przezroczyste, obrys `--dn-blad-obrys`, tekst `--dn-blad-tekst`; wskazanie kursorem → tło `--dn-blad-tlo` |
| Rozmiar mały / duży | `.dn-btn--sm` / `.dn-btn--lg` | wysokość `calc(--dn-wym-kontrolka ± --dn-od-1…2)`, stopień `--dn-fs-sm`/`--dn-fs-md` |

Komentarz źródłowy przy wariancie sygnałowym zastrzega jego rolę wprost: „wyłącznie działanie
systemowe „uruchom / zatwierdź plan”; jeden na widok. Nie zastępuje atramentu jako domyślnego
głównego.” Dwa warianty — atrament i sygnał — nie są więc zamienne: atrament jest domyślnym
działaniem głównym każdego widoku, sygnał jest wyjątkiem zarezerwowanym dla jednego, konkretnego
typu działania.

### 14.2 Plakietka — role i stany

| Wariant | Selektor | Tło / obrys / tekst |
|---|---|---|
| Neutralna | `.dn-plakietka` | `--dn-powierzchnia-2` / `--dn-obrys` / `--dn-tekst-2` |
| Sukces | `.dn-plakietka--sukces` | `--dn-sukces-tlo` / `--dn-sukces-obrys` / `--dn-sukces-tekst` |
| Ostrzeżenie | `.dn-plakietka--ostrzezenie` | `--dn-ostrzezenie-tlo` / `--dn-ostrzezenie-obrys` / `--dn-ostrzezenie-tekst` |
| Błąd | `.dn-plakietka--blad` | `--dn-blad-tlo` / `--dn-blad-obrys` / `--dn-blad-tekst` |
| Informacja | `.dn-plakietka--informacja` | `--dn-informacja-tlo` / `--dn-informacja-obrys` / `--dn-informacja-tekst` |
| Sygnał | `.dn-plakietka--sygnal` | `--dn-sygnal-tlo` / `--dn-sygnal-obrys` / `--dn-sygnal` |
| Rola (koordynator/wykonawca/walidator) | `.dn-plakietka--rola` | krój `--dn-ff-mono`, wersaliki, liternictwo `--dn-ls-mono-wersaliki`, zaokrąglenie `--dn-r-xs` zamiast pigułki |

Sześć wariantów barwnych dzieli jeden szkielet (`.dn-plakietka`: dopełnienie `1px var(--dn-od-2)
2px`, zaokrąglenie `--dn-r-pill`, stopień `--dn-fs-xs`) — jedyną różnicą jest para tło/obrys/tekst,
zgodnie z zasadą, że barwa towarzyszy zawsze tekstowi wewnątrz plakietki, nigdy nie jest jedynym
nośnikiem (rozdz. 1.2, 16). Wariant `--rola` łamie kształt pigułkowy na rzecz zaokrąglenia drobnego
(`--dn-r-xs`) i krój monospace — jedyny wariant plakietki odróżniający się kształtem, nie tylko
barwą.

### 14.3 Powiadomienie (dymek) — struktura i cztery warianty

Kontener `.dn-toasty` jest przyklejony (`position: fixed`) do naroża `inset: auto var(--dn-od-4)
var(--dn-od-4) auto` — dolny prawy róg okna, warstwa `--dn-z-powiadomienie` (rozdz. 11). Pojedynczy
dymek `.dn-toast` układa się w siatkę trzykolumnową (ikona, treść, akcja), szerokość między
`--dn-wym-toast-min` (280px) a `--dn-wym-toast-max` (420px), wejście animowane `dn-wejscie` w czasie
`--dn-czas-3`.

| Wariant | Kreska lewa | Kolor ikony |
|---|---|---|
| Sukces | `.dn-toast--sukces`, `border-left: 2px solid var(--dn-sukces-tekst)` | `--dn-sukces-tekst` |
| Ostrzeżenie | `.dn-toast--ostrzezenie` | `--dn-ostrzezenie-tekst` |
| Błąd | `.dn-toast--blad` | `--dn-blad-tekst` |
| Informacja | `.dn-toast--informacja` | `--dn-informacja-tekst` |

Cztery warianty różnią się wyłącznie kreską lewą grubości `--dn-wym-wstega` (2px) i barwą ikony —
tło pozostaje stałe (`--dn-panel`) dla wszystkich czterech, w odróżnieniu od plakietki (rozdz. 14.2),
gdzie tło zmienia się z wariantem. Rozróżnienie oszczędniejsze wizualnie jest zamierzone: dymek
powiadomienia niesie już tytuł (`.dn-toast-tytul`) i treść (`.dn-toast-tresc`) jako nośnik znaczenia
— kreska i ikona są wzmocnieniem, nie jedynym sygnałem.

### 14.4 Dymek objaśnienia — mechanizm ujawnienia

`.dn-tooltip` jest kontenerem pozycjonującym (`position: relative`); `.dn-tooltip-tresc` jest
domyślnie niewidoczny (`opacity: 0`, `visibility: hidden`, `pointer-events: none`) i przesunięty o
`--dn-od-1` w dół względem pozycji docelowej. Ujawnienie następuje przez selektor złożony:
`.dn-tooltip:hover .dn-tooltip-tresc, .dn-tooltip:focus-within .dn-tooltip-tresc` — dymek pojawia się
zarówno przy wskazaniu kursorem, jak i przy ognisku klawiatury wewnątrz kontenera, jednym mechanizmem
CSS bez skryptu. Tło dymka to `--dn-rama` — jedyny komponent spoza samej ramy aplikacji korzystający
z żetonu ramy kokpitu (rozdz. 3), niezależnie od motywu bieżącego.

### 14.5 Awatar, wirnik, stan pusty

`.dn-awatar` bazowy niesie tło `--dn-grad-atrament` (jeden z dwóch gradientów dozwolonych, rozdz.
12) — jedyne udokumentowane zastosowanie gradientu jako tła w całym arkuszu komponentów, zgodne z
zastrzeżeniem „nigdy tłem przycisku, karty ani sekcji”: awatar nie jest żadnym z tych trzech.
Wariant `.dn-awatar--inteligencja` (jednostka inteligencji — agent albo model) zamienia gradient na
`--dn-grad-sygnal`. Znacznik stanu `.dn-awatar-stan` jest kropką 8×8 px w narożniku, tło
`--dn-sukces-tekst`, obwiedziona pierścieniem `box-shadow: 0 0 0 2px var(--dn-powierzchnia)`
oddzielającym ją wizualnie od awatara pod spodem niezależnie od motywu.

`.dn-spinner` jest okręgiem obracającym się (`dn-obrot`, 0,8s liniowo, w nieskończoność) — obrys
`--dn-obrys-mocny` z górnym segmentem `--dn-sygnal-wypelnienie`, wymiar `--dn-wym-spinner` (14px).

`.dn-pusty-stan` centruje ikonę (28×28px), tytuł (`--dn-fs-lg`, `--dn-fw-polgruba`) i opis
ograniczony do 40 znaków szerokości (`max-width: 40ch`) — jedyny komponent arkusza z jawnym limitem
znaków w definicji CSS, nie tylko w konwencji redakcyjnej.

Trzy rodziny warte dodatkowego odnotowania w prozie (rozdz. 14.6–14.8):

### 14.6 Wpis okna komunikacji (rodzina 11)

Komentarz źródłowy przy tej rodzinie ustala rozstrzygnięcie
projektowe wprost: „dziewięciu nadawców […] Trzy klasy semantyczne (nie dziewięć barw) […] Rola w
obrębie klasy różnicuje się ikoną i plakietką roli, nie kolorem.” Cztery warianty klasy `.dn-wpis`
(`--czlowiek`, `--inteligencja`, `--system`, `--pracuje`) niosą więc rozróżnienie nadawcy przez
tożsamość wizualną (medalion, plakietka), nie przez cztery różne kolory tła — zgodne z zasadą „stan
nigdy samym kolorem” (rozdz. 1.2).

### 14.7 Monitor wykonania (rodzina 12)

Komentarz źródłowy: „trzy wyjścia weryfikacji (rozstrzygnięcie
projektowe): poprawny → zieleń · błędy + licznik obiegów → bursztyn · zdarzenie nieoczekiwane →
czerwień + wstrzymanie kolejki. Zawsze ikona + etykieta.” Cztery warianty klasy `.dn-krok`
(`--bledy`, `--poprawny`, `--pracuje`, `--wstrzymany`) odpowiadają czterem stanom rozstrzygniętym
tym komentarzem.

### 14.8 Pole wpisywania polecenia (rodzina 17)

Nazwa rozdziału w arkuszu — „grot ❯ jako sygnatura
wejścia” — wskazuje literalny znak grota jako element identyfikacyjny pola `.dn-prompt`, odróżniający
je wizualnie od pola formularza zwykłego (rodzina 2).

### 14.9 Pole formularza — pełny rozkład stanów

| Stan | Selektor | Realizacja |
|---|---|---|
| Spoczynek | `.dn-pole-kontrolka` | wysokość minimalna `--dn-wym-kontrolka`, obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia` |
| Tekst zastępczy | `.dn-pole-kontrolka::placeholder` | kolor `--dn-tekst-3` |
| Wskazanie kursorem | `.dn-pole-kontrolka:hover` | obrys → `--dn-tekst-3` |
| Ognisko | `.dn-pole-kontrolka:focus` | obrys → `--dn-fokus`; cień `--dn-cien-sygnal` (pierścień poświaty) |
| Błąd | `.dn-pole-kontrolka[aria-invalid='true']` | obrys → `--dn-blad-tekst` |
| Tylko do odczytu | `.dn-pole-kontrolka[readonly]` | tło → `--dn-powierzchnia-2`, tekst → `--dn-tekst-2` |
| Wieloliniowe | `textarea.dn-pole-kontrolka` | wysokość minimalna podwojona, `resize: vertical` |
| Rozwijane | `select.dn-pole-kontrolka` | dopełnienie prawe powiększone pod strzałkę, `cursor: pointer` |
| Pole wyszukiwania | `.dn-szukaj` | ikona bezwzględnie pozycjonowana po lewej (`pointer-events: none`), pole wewnętrzne z dopełnieniem lewym `--dn-od-8` |

Trzy przejścia (`border-color`, `box-shadow`, `background-color`) dzielą wspólny czas `--dn-czas-2`
(rozdz. 8.2) — pole formularza zmienia wygląd między stanami w tym samym rytmie co przycisk (rozdz.
14.1), zapewniając spójność odczuwalną między dwoma najczęściej używanymi rodzinami kontrolek.

### 14.10 Wybór: pole wyboru, przełącznik, radio, suwak

| Komponent | Klasa | Zachowanie |
|---|---|---|
| Pole wyboru | `.dn-check` (rozdz. 1.2 wielu okien platformowych) | kwadrat `--dn-wym-check` (16px), zaznaczone → tło `--dn-atrament`, obrys `--dn-atrament` |
| Radio | `.dn-radio` | jak pole wyboru, zaokrąglenie `--dn-r-pill` zamiast `--dn-r-xs` |
| Przełącznik | `.dn-przelacznik` | tor `--dn-wym-przelacznik-szer`×`--dn-wym-przelacznik-wys` (36×20px); suwak wewnętrzny przesuwa się `left` w czasie `--dn-czas-2`; zaznaczony → tło `--dn-sygnal-wypelnienie`, suwak → `--dn-szary-0` |
| Suwak zakresu | `.dn-suwak` | tor 4px, gradient dwukolorowy sterowany zmienną CSS `--dn-suwak-pozycja`; uchwyt 16px (WebKit) / 12px (Firefox), obrys `--dn-sygnal-wypelnienie` |

Komentarz źródłowy przy `.dn-suwak` nazywa jego zastosowanie wprost: „nakład rozumowania (szybciej
↔ mądrzej)” — kontrolka jest kandydatem naturalnym dla pola `effort` struktury `SessionConfigModel`
albo pola `ReasoningEffort` przypisanego agentowi (`Agent.effort`, zweryfikowane w
[Oknie nowego projektu](okno-nowego-projektu.md) rozdz. 12.5 przy okazji struktury `Agent`), choć
żadne z czterech okien platformowych tego przejścia redakcyjnego nie zawiera zrzutu użycia tego
konkretnego komponentu.

### 14.11 Karta — interaktywność i karta środowiska

| Stan | Selektor | Realizacja |
|---|---|---|
| Spoczynek | `.dn-karta` | obrys `--dn-obrys`, zaokrąglenie `--dn-r-lg`, cień `--dn-cien-1` |
| Klikalna, wskazanie kursorem | `.dn-karta--klikalna:hover` | obrys → `--dn-obrys-mocny`, cień → `--dn-cien-2`, `transform: translateY(-1px)` — karta unosi się wizualnie 1px |
| Klikalna, wciśnięcie | `.dn-karta--klikalna:active` | powrót do `translateY(0)` |
| Wybrana | `.dn-karta--wybrana` | obrys → `--dn-sygnal-obrys`, tło → `--dn-sygnal-tlo` |
| Karta środowiska | `.dn-karta-srodowiska` | komentarz źródłowy: „Spoczynek bez sygnału; najechanie/aktywność: wstęga górna 2 px” — jedyna karta z sygnałem realizowanym jako wstęga (grubość `--dn-wym-wstega`), nie jako zmiana tła całej powierzchni |

Cztery przejścia jednoczesne na `.dn-karta--klikalna` (`border-color`, `box-shadow`, `transform`,
`background-color`) dzielą czas `--dn-czas-2` — karta unosi się i zmienia obrys w tym samym rytmie,
nie w sekwencji rozłącznej.

### 14.12 Boczna nawigacja modułów — pozycja bieżąca

`.dn-boczna` jest kolumną o stałej szerokości `--dn-wym-boczna` (224px), tło `--dn-panel`. Pozycja
bieżąca (`.dn-boczna-pozycja[aria-current='page']`) niesie tło `--dn-sygnal-tlo` **oraz** znacznik
kreskowy dodatkowy — pseudo-element `::before`, szerokość `--dn-wym-wstega` (2px), wysokość 16px,
tło `--dn-kropka` — przy lewej krawędzi pozycji. Komentarz źródłowy nazywa ten znacznik wprost:
„krótka kreska sygnału zamiast pełnej wstęgi” — rozróżnienie świadome od wstęgi pełnowysokościowej
karty środowiska (rozdz. 14.11): pozycja modułu w bocznej nawigacji jest elementem gęstszym,
mniejszym, więc jej znacznik aktywności jest proporcjonalnie krótszy, nie pełnowysokościowy.

### 14.13 Tabela — nagłówek przyklejony i wiersz interaktywny

`.dn-tabela` niesie nagłówek przyklejony (`th { position: sticky; top: 0 }`, warstwa
`--dn-z-przybornik`, rozdz. 11) — inaczej niż nagłówek kolumn zweryfikowany w
[Oknie historii sesji](okno-historii-sesji.md) rozdz. 4.6, gdzie `.hs-naglowek-kolumn` jest
`aria-hidden="true"` i nieinteraktywny: ten komponent ogólny, `.dn-tabela`, oferuje nagłówek
przyklejony do góry przy przewijaniu, dostępny technicznie, choć samo okno historii sesji z niego
nie korzysta (buduje własny wiersz nagłówka spoza tabeli HTML). Wiersz `tbody tr` reaguje na
wskazanie kursorem (tło `--dn-hover`) i na zaznaczenie (`[aria-selected='true']`, tło
`--dn-sygnal-tlo`) — ten sam żeton tła zaznaczenia co karta wybrana (rozdz. 14.11) i pozycja
bieżąca bocznej nawigacji (rozdz. 14.12), spójny wzorzec „wybrane = `--dn-sygnal-tlo`” powtórzony w
trzech niezależnych rodzinach komponentów.

### 14.14 Zakładki — kreska dolna i wariant pigułkowy

Zakładka bazowa (`.dn-zakladka`) niesie kreskę dolną 2px (`::after`, tło `--dn-kropka`) pod pozycją
wybraną (`[aria-selected='true']`) — wzorzec zakładek przeglądarkowych. Wariant pigułkowy
(`.dn-zakladki--pigulki`) zastępuje kreskę wypełnieniem pełnym: komentarz źródłowy uzasadnia to
wprost — „Nie ma dolnej kreski zakładek: wybór niesie wypełniona pigułka. Stan wybrany dostaje
barwę sygnałową i kontrę — sam kształt bez koloru czytałby się jako pole nieaktywne, a to jest
kontrolka wyboru, nie etykieta.” Pozycja wybrana wariantu pigułkowego niesie tło
`--dn-sygnal-wypelnienie` i tekst `--dn-tekst-inv` (odwrócony) — ten sam wariant jest zweryfikowany
w [Oknie historii sesji](okno-historii-sesji.md) rozdz. 5.1 jako „Zakres środowisk”: pięć pozycji
(„Wszystkie”, cztery nazwy środowisk), `role="tablist"`/`role="tab"`, dokładnie klasa
`.dn-zakladki.dn-zakladki--pigulki` przywołana tam wprost.

### 14.15 Modal (dialog) — struktura, warianty szerokości i trzy weryfikacje niezależne

`.dn-modal` bazowy: szerokość `min(var(--dn-wym-modal), calc(100vw - var(--dn-od-8)))` — a więc
560px jako sufit, pomniejszony o `--dn-od-8` (32px) z każdej strony na ekranach węższych; wysokość
maksymalna `min(80dvh, 720px)`; tło `--dn-panel`; cień `--dn-cien-lg`; wejście animowane
`dn-wejscie` (przesunięcie `--dn-od-2` w dół + skala `0.98`→`1`, czas `--dn-czas-3`). Warstwa pod
nakładką (`::backdrop`) korzysta z `--dn-nakladka` z rozmyciem `blur(2px)`. Wariant szeroki
(`.dn-modal--szeroki`) podnosi sufit do `min(1080px, calc(100vw - var(--dn-od-8)))` i wysokość do
`min(86dvh, 780px)`.

| Okno | Szerokość bazowa użyta | Zgodność z dwoma wariantami arkusza |
|---|---|---|
| Bazowy `.dn-modal` (bez modyfikatora) | `min(560px, …)` | wariant standardowy |
| `.dn-modal--szeroki` | `min(1080px, …)` | wariant szeroki |
| [Okno historii sesji](okno-historii-sesji.md) (`.hs-okno`) | własna, nienazwana w arkuszu ogólnym | rozszerza `.dn-modal` lokalnie, poza dwoma wariantami nazwanymi tutaj |
| [Okno instrukcji](okno-instrukcji.md) (`.iu-okno`) | `min(1180px, calc(100vw - var(--dn-od-8)))` | **przekracza** nawet wariant szeroki (1080px) — trzeci poziom szerokości, zweryfikowany w tamtym dokumencie rozdz. 4, nienazwany w `komponenty.css` |

Trzy okna nakładkowe zweryfikowane w tym przejściu redakcyjnym potwierdzają niezależnie: żadne z
nich nie zatrzymuje się na szerokości bazowej 560px — [Okno nowego projektu](okno-nowego-projektu.md)
w ogóle nie jest modalem (rozdz. 1.1 tamtego dokumentu, korekta względem redakcji poprzedniej), a
oba pozostałe (historia sesji, instrukcja) rozszerzają szerokość bazową ponad wariant szeroki
nazwany. Wzorzec zbieżny z rozdz. 10 i rozdz. 18.4 niniejszego dokumentu: okna platformowe ustalają
wymiary własne wprost w arkuszu lokalnym, nie ograniczając się do dwóch wariantów `.dn-modal`
nazwanych w bibliotece ogólnej.

---

## 15. Ikonografia

`design/zasoby/ikony/manifest.json` niesie **osiemdziesiąt dwie ikony**, nie czterdzieści siedem —
liczba jest polem `liczba-ikon` dokumentu, odczytana wprost, nie oszacowana. Zasady wspólne
wszystkim ikonom, pole `zasady` manifestu:

| Zasada | Wartość |
|---|---|
| Siatka | `24×24` |
| Grubość obrysu | `1.75` |
| Wypełnienie | `none` |
| Barwa | `currentColor` — dziedziczy kolor tekstu otoczenia, nigdy własny kolor zaszyty |
| Zakończenia linii | `round` |
| Łączenia linii | `round` |
| Renderowanie | `14, 16, 20, 24` px — cztery rozmiary docelowe, zgodne ze skalą `--dn-wym-ikona-*` (rozdz. 9), pokrywającą `--dn-wym-ikona-sm` (16px) i `--dn-wym-ikona` (18px) w praktyce z zaokrągleniem |

Siedemdziesiąt osiem z osiemdziesięciu dwóch ikon pochodzi z biblioteki zewnętrznej **Lucide**
(pole `zrodlo` w postaci `lucide:{nazwa}`); cztery ikony środowisk (`srodowisko-talkin`,
`srodowisko-workspace`, `srodowisko-codestudio`, `srodowisko-multitaskingai`) są własne (`zrodlo:
wlasna`) — jedyne cztery w całym katalogu niebędące adaptacją zestawu zewnętrznego, każda niesie
motto środowiska jako pole `zastosowanie` zamiast opisu funkcjonalnego (np. „Środowisko TalkIn —
Myśl. Analizuj. Rozumiej.”) — cztery motta odpowiadają dosłownie czterem środowiskom platformy.
Wykaz pełny osiemdziesięciu dwóch pozycji stoi w
[Załączniku C](#załącznik-c-pełny-katalog-ikon).

---

## 16. Stany interakcji i zasada zero blokad

Nagłówek `komponenty.css` (rozdz. 1.2) nazywa siedem stanów, którymi dysponuje każdy komponent
interaktywny: spoczynek, wskazanie kursorem, naciśnięcie, fokus, wybrany, ładowanie,
błąd. Rozdział 4.2 standardu redakcyjnego zbioru dodaje do tego wykazu stan „nieaktywny” i „pusty” —
w systemie Danaco Console pierwszy z nich jest **stanem wyłączonym z definicji**:

```
┌───────────────────────────────────────────────────────────────────┐
│  ZASADA ZERO BLOKAD — komponenty.css, w. 7                        │
├───────────────────────────────────────────────────────────────────┤
│  Żaden wariant komponentu NIE ODBIERA klikalności.                │
│  Niegotowość komunikuje się:                                      │
│    • opisem przy kontrolce                                        │
│    • komunikatem po próbie działania                              │
│  Nigdy: atrybutem `disabled`, odjęciem `pointer-events`,          │
│         szarym, nieklikalnym wariantem wizualnym.                 │
└───────────────────────────────────────────────────────────────────┘
```

Cztery okna platformowe zweryfikowane w tym przejściu redakcyjnym potwierdzają tę zasadę w praktyce,
niezależnie od siebie: [Okno instalatora](okno-instalatora.md) rozdz. 6 zastępuje przycisk wyłączony
komunikatem walidacji przy próbie „Dalej” bez akceptacji licencji; [Okno nowego projektu](okno-nowego-projektu.md)
rozdz. 7 robi to samo dla przycisku „Załóż projekt i otwórz sesję”; [Okno historii sesji](okno-historii-sesji.md)
rozdz. 7 nie niesie ani jednego wariantu „nieaktywny” w tabeli stanów kontrolek, konsekwentnie
zapisując „brak (zasada zero blokad)” tam, gdzie inny system pokazałby przycisk wyszarzony.

Stan „nigdy samym kolorem” (rozdz. 1.2) obowiązuje równolegle: pierścień fokusu (`--dn-fokus`,
rozdz. 4–5) towarzyszy zawsze zmianie kształtu (obrys przerywany na ciągły, grubość 2 px), nie samej
zmianie barwy; kropka sygnału stanu czynnego (rozdz. 2.2, `--dn-kropka`) towarzyszy zawsze tekstowi
albo etykiecie nazywającej stan wprost — zweryfikowane w [Oknie historii sesji](okno-historii-sesji.md)
rozdz. 4.1, gdzie znak `.pt-tetno` nigdy nie stoi bez towarzyszącego słowa „czynna” w kolumnie
sąsiedniej.

---

## 17. Dostępność i kontrast

`budowa/shared/kontrasty-progi.json` niesie trzydzieści trzy pary zmierzone, każda wobec progu
domyślnego **4,5∶1** — kryterium WCAG 2.1 AA dla tekstu zwykłego (kryterium 1.4.3), pole
`prog_zrodlo` dokumentu wprost. Siedemnaście par dotyczy motywu jasnego, szesnaście motywu ciemnego.

| Para zmierzona | Żetony | Motyw |
|---|---|---|
| tekst / tło | `tekst`, `tlo` | jasny i ciemny |
| tekst / powierzchnia | `tekst`, `powierzchnia` | jasny i ciemny |
| tekst-2 / tło | `tekst-2`, `tlo` | jasny i ciemny |
| tekst-2 / powierzchnia | `tekst-2`, `powierzchnia` | jasny i ciemny |
| tekst-2 / powierzchnia-2 | `tekst-2`, `powierzchnia-2` | jasny |
| tekst-3 / tło | `tekst-3`, `tlo` | jasny |
| tekst-3 / powierzchnia | `tekst-3`, `powierzchnia` | jasny i ciemny |
| biel / przycisk atrament | `atrament-tekst`, `atrament` | jasny i ciemny |
| sygnal-tekst / tło | `sygnal`, `tlo` | jasny i ciemny |
| sygnal-tekst / powierzchnia | `sygnal`, `powierzchnia` | jasny i ciemny |
| sygnal-tekst / sygnal-tlo | `sygnal`, `sygnal-tlo` | jasny |
| biel / sygnał wypełnienie | `szary-0`, `sygnal-wypelnienie` | jasny i ciemny |
| obrys kontrolki / tło | `obrys-mocny`, `tlo` | jasny |
| fokus sygnał / tło | `fokus`, `tlo` | jasny i ciemny |
| sukces tekst / tło stanu | `sukces-tekst`, `sukces-tlo` | jasny |
| ostrzeżenie tekst / tło stanu | `ostrzezenie-tekst`, `ostrzezenie-tlo` | jasny |
| błąd tekst / tło stanu | `blad-tekst`, `blad-tlo` | jasny |
| sukces / ostrzeżenie / błąd tekst / powierzchnia | odpowiednie pary `*-tekst`, `powierzchnia` | jasny i ciemny |

Dokument źródłowy niesie pole `pochodzenie` z zastrzeżeniem wprost: plik jest **odtworzeniem**, nie
kopią pierwowzoru, który „leży na maszynie Właściciela … i NIE JEST pod kontrolą wersji”. Wartości
progów (4,5∶1) i lista trzydziestu trzech par są jednak kompletne w pliku dostępnym w repozytorium i
stanowią źródło normatywne wiążące dla niniejszego dokumentu, niezależnie od statusu pierwowzoru.

Trzy wartości kontrastu zmierzone wprost i przywołane liczbowo w tabelach rozdz. 4–5: tekst główny
na tle — 16,1∶1 (jasny) i 15,0∶1 (ciemny); tekst drugorzędny — 6,4∶1 (jasny) i 6,8∶1 (ciemny);
sygnał jako tekst — 5,5∶1 (jasny) i 8,0∶1 (ciemny). Wszystkie sześć przekracza próg 4,5∶1 z zapasem.

Trzydzieści trzy pary są zmierzone, nie oszacowane — pole `dokument` pliku źródłowego nazywa je
„WEJŚCIE POMIARU” dla przyrządu automatycznego (`budowa/scripts/kontrasty.mjs`), nie tabelą
poglądową. Redaktor dokumentacji nie ma dostępu do wyniku tego przyrządu z poziomu niniejszego
przejścia redakcyjnego — rozdział opiera się na wykazie par i progu wejściowym, nie na wyniku
przebiegu pomiaru, którego plik nie niesie.

---

## 18. Zastosowanie zweryfikowane w czterech oknach platformowych

Cztery okna platformowe zredagowane w tym samym przejściu redakcyjnym —
[Okno instalatora](okno-instalatora.md), [Okno historii sesji](okno-historii-sesji.md),
[Okno nowego projektu](okno-nowego-projektu.md), [Okno instrukcji](okno-instrukcji.md) — są
zbudowane wyłącznie na żetonach i klasach opisanych w rozdziałach 2–15 powyżej, potwierdzone przez
inspekcję źródłową każdego z czterech prototypów HTML wymienionych w polu „Prototypy odniesienia”
tych dokumentów. Rozdział zbiera
punkty styku między systemem opisanym abstrakcyjnie tutaj a jego użyciem konkretnym w tych czterech
oknach, jako weryfikację wzajemną dwóch grup dokumentów tego samego przejścia.

### 18.1 Żetony ramy kokpitu (rozdz. 3) w oknach wejściowych

[Okno instalatora](okno-instalatora.md) rozdz. 3.2 i [Okno nowego projektu](okno-nowego-projektu.md)
rozdz. 6 zweryfikowały niezależnie tę samą parę żetonów kontrastowych: `--dn-rama` (stały w obu
motywach, rozdz. 3 niniejszego dokumentu) wobec `--dn-powierzchnia`/`--dn-panel` (zmienne z motywem,
rozdz. 4–5). Obie weryfikacje, wykonane osobno, doszły do tego samego wniosku: kolumna tożsamości
albo panel podsumowania kontrastuje wyraźnie z otoczeniem w motywie jasnym, a w motywie ciemnym
kontrast opiera się wyłącznie na kresce obrysu, bo obie powierzchnie zbliżają się jasnością — zgodne
z komentarzem źródłowym „w motywie ciemnym obie kolumny mają zbliżoną jasność — kreska rozdziela je
niezależnie od motywu” (rozdz. 3).

### 18.2 Zasada zero blokad (rozdz. 16) w praktyce czterech okien

| Okno | Kontrolka | Zachowanie zero blokad zweryfikowane |
|---|---|---|
| Instalator | „Dalej” bez akceptacji licencji | pozostaje klikalny; komunikat walidacji zamiast blokady (rozdz. 6 tamtego dokumentu) |
| Nowy projekt | „Załóż projekt i otwórz sesję” bez nazwy | pozostaje klikalny; komunikat „Nadaj projektowi nazwę” |
| Historia sesji | brak przycisku „Usuń” w ogóle | zamiast wyłączonego przycisku — nieobecność kontrolki, udokumentowana jako luka jawna (rozdz. 6.4 tamtego dokumentu), nie jako przycisk zablokowany |
| Instrukcja | brak kontrolek interaktywnych generujących komunikat | okno w całości statyczne (rozdz. 21 tamtego dokumentu) — zasada zero blokad nie ma tu zastosowania praktycznego, bo nie ma czynności do zablokowania |

### 18.3 Klasy `.pt-*` przywołane obok `.dn-*` w oknach wejściowych

Dwa dokumenty — [Okno instalatora](okno-instalatora.md) rozdz. 3.4 i
[Okno historii sesji](okno-historii-sesji.md) rozdz. 3.4 — odnotowują niezależnie tę samą
obserwację: animacja tętna sesji/kropki czynnej korzysta w praktyce z klasy `.pt-tetno`
(`design/zasoby/prototyp.css`), nie z `.dn-kropka--tetno` (rodzina 4, rozdz. 14 niniejszego
dokumentu) — dwie klasy dają wizualnie zbliżony efekt tętna, lecz są zdefiniowane osobno, w dwóch
różnych arkuszach. Zgodnie z zakresem niniejszego dokumentu (metryka, pole „Poza zakresem”), klasy
`.pt-*` nie są tu katalogowane w pełni — przywołanie to jest wyjątkiem uzasadnionym bezpośrednim
punktem styku z żetonem `--dn-czas-tetno` (rozdz. 8.2) opisanym tutaj.

### 18.4 Progi łamania lokalne — czwarty potwierdzony wzorzec

Rozdział 10 niniejszego dokumentu odnotował już trzy progi lokalne nienazwane (1000/640 px, 900 px,
1240 px). [Okno instrukcji](okno-instrukcji.md) rozdz. 22.2 zestawia czwarty, tego samego wzorca —
900 px własny tego okna, różny liczbowo od progu 900 px [Okna historii sesji](okno-historii-sesji.md),
choć obie wartości są identyczne co do liczby: dwa okna niezależnie wybrały tę samą wartość progu
bez odwołania do wspólnego żetonu `--dn-bp-*`, co sugeruje zbieżność praktyki redakcyjnej, nie
współdzielony mechanizm.

---

## 19. Kryteria odbioru

| Warunek | Sposób sprawdzenia |
|---|---|
| Wszystkie 309 żetonów `--dn-*` niniejszego dokumentu odpowiadają dosłownie `zetony.css` | porównanie Załącznika A z arkuszem źródłowym, wartość po wartości |
| Wszystkie 123 klasy `.dn-*` niniejszego dokumentu odpowiadają dosłownie `komponenty.css` | porównanie Załącznika B z arkuszem źródłowym |
| Wszystkie 82 ikony niniejszego dokumentu odpowiadają dosłownie `manifest.json` | porównanie Załącznika C z polem `ikony` manifestu |
| Żaden komponent nie odwołuje się do prymitywu (rozdz. 2) wprost, wyłącznie do żetonu semantycznego | inspekcja `komponenty.css` — brak `--dn-szary-*`/`--dn-sygnal-*` poza definicjami warstwy 2 w `zetony.css` |
| Oba motywy niosą komplet czterdziestu żetonów semantycznych, zdefiniowanych osobno | porównanie bloków `[data-theme='light']` i `[data-theme='dark']` — te same nazwy, różne wartości gdzie zasadne |
| Siedem żetonów ramy kokpitu (rozdz. 3) ma tę samą wartość w obu motywach | odczyt wartości obliczonej w obu trybach `data-theme` |
| Żaden gradient (rozdz. 12) nie jest tłem przycisku, karty ani sekcji | przegląd `komponenty.css` — brak `--dn-grad-*` poza dwoma zastosowaniami nazwanymi |
| Cztery żetony czasu spadają do `0.01ms` przy `prefers-reduced-motion: reduce` | test z preferencją systemową ustawioną |
| Sześć żetonów wymiaru rośnie przy `pointer: coarse` | test na urządzeniu dotykowym albo emulacji |
| Trzydzieści trzy pary kontrastu z `kontrasty-progi.json` spełniają próg 4,5∶1 | pomiar automatyczny wartości obliczonych względem progu |
| Żaden komponent interaktywny nie niesie atrybutu `disabled` ani stanu odbierającego klikalność | przegląd `komponenty.css` — zero wystąpień `disabled`, zero `pointer-events: none` na komponentach interaktywnych |
| Punkty styku rozdz. 18 między niniejszym dokumentem a czterema oknami platformowymi pozostają zgodne po każdej zmianie któregokolwiek z pięciu dokumentów | porównanie krzyżowe przy każdej redakcji |
| Wariant `.dn-btn--sygnal` pozostaje jednym na widok, nigdy zamiennikiem `.dn-btn--atrament` | przegląd każdego okna platformowego pod kątem liczby przycisków sygnałowych jednocześnie widocznych |

---

## Załącznik A. Pełna tabela żetonów `--dn-*`

Osiemnaście wierszy rodzin nazewniczych, w kolejności występowania w `zetony.css`. Rodziny zgodne
z rozdz. 7.3 standardu redakcyjnego zbioru: `wym`, `sygnal`, `cien`, `szary`, `tekst`, `fs`, `od`
i pozostałe poniżej — standard wylicza łącznie dwadzieścia siedem rodzin żetonów, wszystkie
potwierdzone obecne w arkuszu źródłowym.

| Rodzina | Liczba żetonów | Zakres nazw | Rozdział niniejszego dokumentu |
|---|---:|---|---|
| `szary` | 18 | `--dn-szary-0` … `--dn-szary-1000` | rozdz. 2.1 |
| `sygnal` | 14 | `--dn-sygnal-100…800`, `--dn-sygnal`, `--dn-sygnal-mocny`, `--dn-sygnal-wypelnienie[-hover]`, `--dn-sygnal-tlo`, `--dn-sygnal-obrys` | rozdz. 2.2, 4, 5 |
| `zielen` / `bursztyn` / `czerwien` | 4 + 4 + 4 = 12 (prymitywy) | `-700`, `-400`, `-200`, `-100` każda | rozdz. 2.3 |
| `rama` | 7 | `--dn-rama`, `--dn-rama-grafit`, `--dn-rama-tekst[-2]`, `--dn-rama-hover`, `--dn-rama-obrys[-mocny]` | rozdz. 3 |
| `ff` | 3 | `--dn-ff-naglowek`, `--dn-ff-bazowa`, `--dn-ff-mono` | rozdz. 6.1 |
| `fs` | 9 | `--dn-fs-xs` … `--dn-fs-display` | rozdz. 6.2 |
| `fw` | 4 | `--dn-fw-normalna` … `--dn-fw-gruba` | rozdz. 6.2 |
| `lh` | 3 | `--dn-lh-ciasny`, `--dn-lh-bazowy`, `--dn-lh-luzny` | rozdz. 6.2 |
| `ls` | 3 | `--dn-ls-naglowek`, `--dn-ls-wersaliki`, `--dn-ls-mono-wersaliki` | rozdz. 6.2 |
| `od` (+`odstep`) | 11 + 2 | `--dn-od-0` … `--dn-od-16`, `--dn-odstep-panel`, `--dn-odstep-sekcji` | rozdz. 7.1 |
| `r` | 6 | `--dn-r-xs` … `--dn-r-pill` | rozdz. 7.2 |
| `ease` / `czas` | 1 + 4 | `--dn-ease`, `--dn-czas-1…3`, `--dn-czas-tetno` | rozdz. 8.2 |
| `cien` | 5 | `--dn-cien-1…3`, `--dn-cien-lg`, `--dn-cien-sygnal` | rozdz. 8.1 |
| `wym` | 28 | `--dn-wym-kontrolka` … `--dn-wym-fokus-odsuniecie` | rozdz. 9 |
| `bp` (+`tresc`, `siatka`) | 4 + 1 + 2 | `--dn-bp-w1…w4`, `--dn-tresc-max`, `--dn-siatka-kolumny`, `--dn-siatka-przerwa` | rozdz. 10 |
| `z` | 11 | `--dn-z-podloga` … `--dn-z-centrum-polecen` | rozdz. 11 |
| `grad` | 2 | `--dn-grad-atrament`, `--dn-grad-sygnal` | rozdz. 12 |
| `tlo`/`powierzchnia`/`panel`/`wstazka`/`obrys`/`tekst`/`atrament`/`fokus`/`hover`/`wcisniecie`/`nakladka`/`sukces`/`ostrzezenie`/`blad`/`informacja`/`kropka` (semantyczne) | 34 | rozdz. 4, 5 pełny wykaz | rozdz. 4, 5 |

Suma kontrolna: 18 + 14 + 12 + 7 + 3 + 9 + 4 + 3 + 3 + 13 + 6 + 5 + 5 + 28 + 7 + 11 + 2 + 34 =
**184** żetony — tyle nazw własnych `--dn-*` niesie arkusz `design/zasoby/zetony/zetony.css` i tyle
wyliczają łącznie osiemnaście wierszy rodzin powyżej. Miarą wiążącą dla całego zbioru jest liczba
nazw własnych, nie liczba deklaracji: arkusz zawiera 309 deklaracji, bo bloki motywu jasnego,
motywu ciemnego, gęstości przestronnej oraz wskazywania dotykiem powtarzają te same nazwy
z innymi wartościami. Powtórzenie nazwy w takim bloku nie jest osobnym żetonem — jest inną
wartością tego samego żetonu. Wartości dosłowne każdego żetonu stoją w
rozdziałach 2–13 powyżej — załącznik służy jako mapa rodzin dla szybkiego odnalezienia rozdziału
właściwego, nie jako druga kopia wartości.

---

## Załącznik B. Pełna tabela klas `.dn-*`

Sto dwadzieścia trzy klasy w osiemnastu rodzinach, w kolejności występowania w `komponenty.css`.

| Rodzina | Klasy pełne |
|---|---|
| Przycisk | `.dn-btn` · `.dn-btn--atrament` · `.dn-btn--duch` · `.dn-btn--lg` · `.dn-btn--niebezpieczny` · `.dn-btn--sm` · `.dn-btn--sygnal` · `.dn-btn--wybrany` · `.dn-btn--zarys` · `.dn-btn-ikona` · `.dn-btn-ikona--na-ramie` |
| Pole formularza | `.dn-pole` · `.dn-pole-blad` · `.dn-pole-etykieta` · `.dn-pole-kontrolka` · `.dn-pole-opis` · `.dn-szukaj` |
| Wybór | `.dn-check` · `.dn-przelacznik` · `.dn-radio` · `.dn-suwak` · `.dn-wybor` |
| Kropka sygnału | `.dn-kropka` · `.dn-kropka--blad` · `.dn-kropka--neutralna` · `.dn-kropka--ostrzezenie` · `.dn-kropka--sukces` · `.dn-kropka--tetno` |
| Plakietka | `.dn-plakietka` · `.dn-plakietka--blad` · `.dn-plakietka--informacja` · `.dn-plakietka--ostrzezenie` · `.dn-plakietka--rola` · `.dn-plakietka--sukces` · `.dn-plakietka--sygnal` |
| Karta | `.dn-kafel` · `.dn-kafel-etykieta` · `.dn-kafel-ikona` · `.dn-kafel-opis` · `.dn-karta` · `.dn-karta--klikalna` · `.dn-karta--wybrana` · `.dn-karta-cialo` · `.dn-karta-naglowek` · `.dn-karta-srodowiska` · `.dn-karta-srodowiska-godlo` · `.dn-karta-srodowiska-motto` · `.dn-karta-srodowiska-opis` · `.dn-karta-srodowiska-tytul` · `.dn-karta-tytul` |
| Tabela | `.dn-dane` · `.dn-tabela` |
| Zakładki · karty sesji | `.dn-karta-sesji` · `.dn-karty-sesji` · `.dn-zakladka` · `.dn-zakladki` · `.dn-zakladki--pigulki` |
| Rama kokpitu | `.dn-listwa` · `.dn-listwa-pozycja` · `.dn-pasek` · `.dn-pasek-godlo` · `.dn-pasek-logotyp` · `.dn-pasek-prawa` · `.dn-pasek-szukaj` |
| Boczna nawigacja | `.dn-boczna` · `.dn-boczna-naglowek` · `.dn-boczna-pozycja` |
| Wpis okna komunikacji | `.dn-wpis` · `.dn-wpis--czlowiek` · `.dn-wpis--inteligencja` · `.dn-wpis--pracuje` · `.dn-wpis--system` · `.dn-wpis-godzina` · `.dn-wpis-medalion` · `.dn-wpis-nadawca` · `.dn-wpis-tozsamosc` · `.dn-wpis-tresc` |
| Monitor wykonania | `.dn-kolejka` · `.dn-krok` · `.dn-krok--bledy` · `.dn-krok--poprawny` · `.dn-krok--pracuje` · `.dn-krok--wstrzymany` · `.dn-krok-meta` · `.dn-krok-znak` · `.dn-postep` · `.dn-postep-etykieta` · `.dn-postep-tor` · `.dn-postep-wartosc` |
| Modal · nakładka | `.dn-modal` · `.dn-modal--szeroki` · `.dn-modal-cialo` · `.dn-modal-naglowek` · `.dn-modal-stopka` · `.dn-modal-tytul` · `.dn-modal-zamknij` |
| Powiadomienie | `.dn-toast` · `.dn-toast--blad` · `.dn-toast--informacja` · `.dn-toast--ostrzezenie` · `.dn-toast--sukces` · `.dn-toast-tresc` · `.dn-toast-tytul` · `.dn-toasty` |
| Dymek objaśnienia | `.dn-tooltip` · `.dn-tooltip-tresc` |
| Awatar · wirnik · pusty stan | `.dn-awatar` · `.dn-awatar--inteligencja` · `.dn-awatar--kwadrat` · `.dn-awatar--lg` · `.dn-awatar--sm` · `.dn-awatar-stan` · `.dn-pusty-stan` · `.dn-pusty-stan-opis` · `.dn-pusty-stan-tytul` · `.dn-spinner` |
| Pole wpisywania polecenia | `.dn-prompt` · `.dn-prompt-grot` · `.dn-prompt-obszar` · `.dn-przybornik` |
| Always On Display | `.dn-aod` · `.dn-aod-rdzen` · `.dn-aod-tresc` |

Suma: 11+6+5+6+7+15+2+5+7+3+10+12+7+8+2+10+4+3 = **123** klasy — zgodna z liczbą zweryfikowaną
programowo w `komponenty.css`. Osiemnaście rodzin odpowiada dokładnie osiemnastu blokom
komentarzowym wyznaczającym sekcje arkusza źródłowego, wymienionym w tej samej kolejności co w
pliku — załącznik nie zmienia ani nie grupuje na nowo podziału już ustalonego przez autora arkusza.

---

## Załącznik C. Pełny katalog ikon

Osiemdziesiąt dwie pozycje, w kolejności występowania w polu `ikony` manifestu. Kolumna „Źródło” niesie
`lucide:{nazwa}` dla siedemdziesięciu ośmiu pozycji adaptowanych z biblioteki Lucide i `własna` dla
czterech ikon środowisk.

| Nazwa | Źródło | Zastosowanie |
|---|---|---|
| `dom` | lucide:house | Centrum dowodzenia — powrót do strony głównej |
| `menu` | lucide:menu | Rozwinięcie nawigacji bocznej, menu podręczne |
| `strzalka-lewo` | lucide:arrow-left | Cofnięcie, powrót |
| `strzalka-prawo` | lucide:arrow-right | Przejście dalej, wejście w pozycję |
| `grot-dol` | lucide:chevron-down | Rozwinięcie sekcji, lista rozwijana |
| `grot-gora` | lucide:chevron-up | Zwinięcie sekcji |
| `grot-prawo` | lucide:chevron-right | Wejście w głąb, okruszki nawigacji |
| `wiecej` | lucide:ellipsis | Menu dodatkowych czynności pozycji |
| `link-zewnetrzny` | lucide:external-link | Odesłanie poza bieżący widok |
| `szukaj` | lucide:search | Wyszukiwanie — pasek górny, Command Center |
| `filtr` | lucide:funnel | Zawężenie wykazu, panel warunków |
| `zamknij` | lucide:x | Zamknięcie karty, modala, powiadomienia |
| `karta-okna` | lucide:app-window | Okno, karta przeglądarki, moduł Browser |
| `plus` | lucide:plus | Utworzenie bytu: sesja, okno, kanał |
| `olowek` | lucide:pencil | Zmiana nazwy, edycja |
| `kosz` | lucide:trash-2 | Usunięcie pozycji (z potwierdzeniem) |
| `pobierz` | lucide:download | Pobranie pliku, eksport wyniku |
| `wgraj` | lucide:upload | Wgranie pliku, import |
| `wyslij` | lucide:send-horizontal | Wysłanie komunikatu w oknie komunikacji |
| `odswiez` | lucide:refresh-cw | Ponowienie próby, odświeżenie danych |
| `odpowiedz` | lucide:reply | Odpowiedź, cytowanie wpisu |
| `uruchom` | lucide:play | Uruchomienie procesu, zadania, automatyzacji |
| `zatrzymaj` | lucide:square | Zatrzymanie procesu |
| `wstrzymaj` | lucide:pause | Wstrzymanie kolejki, pauza |
| `spinacz` | lucide:paperclip | Załącznik komunikatu |
| `ustawienia` | lucide:settings | Konfiguracja — okna, kanału, środowiska |
| `oko` | lucide:eye | Podgląd, odsłonięcie wartości ukrytej |
| `gwiazdka` | lucide:star | Wyróżnienie pozycji, presety zespołów |
| `kopiuj` | lucide:copy | Kopiowanie treści, identyfikatora |
| `ptaszek` | lucide:check | Potwierdzenie wykonania, zaznaczenie |
| `ptaszek-kolo` | lucide:circle-check | Powodzenie — proces zakończony poprawnie |
| `blad` | lucide:circle-x | Błąd — niepowodzenie operacji |
| `ostrzezenie` | lucide:triangle-alert | Ostrzeżenie — skutki wymagające uwagi |
| `info` | lucide:info | Wyjaśnienie, podpowiedź kontekstowa |
| `klodka` | lucide:lock | Treść chroniona, dane dostępowe |
| `tarcza` | lucide:shield | Uprawnienia, zakres dostępu, izolacja |
| `zegar` | lucide:clock | Oczekiwanie, znacznik czasu |
| `historia` | lucide:history | Historia sesji, wersje |
| `dzwonek` | lucide:bell | Powiadomienia |
| `aktywnosc` | lucide:activity | Praca w tle, telemetria procesu |
| `dokument` | lucide:file-text | Dokument, opracowanie — moduł Studio |
| `plik` | lucide:file | Plik o nieokreślonym rodzaju |
| `folder` | lucide:folder | Katalog roboczy, grupa zasobów |
| `archiwum` | lucide:archive | Zasoby odłożone, Session Repository |
| `obraz` | lucide:image | Załącznik graficzny, moduł Design (zasoby) |
| `kod` | lucide:code-xml | Fragment kodu, moduł Developer |
| `koperta` | lucide:mail | Wiadomość, korespondencja |
| `kalendarz` | lucide:calendar | Harmonogram, automatyzacja cykliczna |
| `tabela-danych` | lucide:table | Zestawienie, dane tabelaryczne |
| `biblioteka` | lucide:library-big | Moduł Library — repozytorium wiedzy |
| `baza` | lucide:database | Trwałość, magazyn danych |
| `wykres` | lucide:chart-line | Metryki, moduł Research (wyniki) |
| `terminal` | lucide:square-terminal | Moduł Terminal — powłoki wykonawcze |
| `galaz` | lucide:git-branch | Repozytorium, moduł Developer |
| `agent` | lucide:bot | Agent — jednostka inteligencji |
| `agenci` | lucide:users | Zespół, moduł Agents / role |
| `uzytkownik` | lucide:user | Operator, konto, tożsamość |
| `siec` | lucide:network | Orkiestracja, topologia — MultitaskingAI |
| `wezly` | lucide:waypoints | Koordynator — plan i delegacja |
| `wykonawca` | lucide:circle-play | Wykonawca — realizacja zlecenia |
| `walidator` | lucide:badge-check | Walidator — weryfikacja kroku |
| `warstwy` | lucide:layers | Moduł Workspace — przestrzeń projektowa |
| `automatyzacja` | lucide:workflow | Moduł Automations — procesy i kolejki |
| `debata` | lucide:messages-square | Moduł Roundtable — debata modeli |
| `tlumacz` | lucide:languages | Moduł Translate |
| `badanie` | lucide:telescope | Moduł Research — analizy i raporty |
| `aplikacje` | lucide:layout-grid | Moduł Apps — budowa produktów |
| `paleta` | lucide:palette | Moduł Design — materiały wizualne |
| `rozmowa` | lucide:message-square | Okno komunikacji — pas komunikacji |
| `mikrofon` | lucide:mic | Moduł Assistant — interfejs głosowy |
| `diagnostyka` | lucide:stethoscope | Moduł Diagnostics — analiza problemów |
| `monitor` | lucide:monitor | Stanowisko, środowisko wykonania |
| `telefon` | lucide:smartphone | Widok mobilny — pozycja „Mobile” |
| `cpu` | lucide:cpu | Zasoby wykonawcze, serwer rdzenia |
| `polecenie` | lucide:command | Command Center — paleta poleceń |
| `slonce` | lucide:sun | Przełącznik motywu — jasny |
| `ksiezyc` | lucide:moon | Przełącznik motywu — ciemny |
| `globus` | lucide:globe | Sieć, źródło zewnętrzne, danaco-web |
| `srodowisko-talkin` | własna | Środowisko TalkIn — Myśl. Analizuj. Rozumiej. |
| `srodowisko-workspace` | własna | Środowisko WorkSpace — Planuj. Organizuj. Realizuj. |
| `srodowisko-codestudio` | własna | Środowisko CodeStudio — Projektuj. Buduj. Rozwijaj. |
| `srodowisko-multitaskingai` | własna | Środowisko MultitaskingAI — Deleguj. Koordynuj. Nadzoruj. |

---

*Koniec dokumentu. System wizualny — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
