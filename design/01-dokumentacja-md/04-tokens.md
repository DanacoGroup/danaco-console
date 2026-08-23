# Danaco Console — Żetony systemu wizualnego (design tokens)

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — pełna dokumentacja żetonów (design tokens) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | prowadzący system projektowy · projektanci interfejsu · zespół wdrożeniowy (front-end) · osoba odpowiedzialna za dostępność · audyt zgodności wizualnej |
| **Zakres** | Kompletny rejestr **177 unikalnych nazw żetonów** systemu: trzy warstwy modelu, konwencja nazewnicza, 38 prymitywów barwnych, 5 żetonów ramy kokpitu, 22 żetony typografii, 17 żetonów przestrzeni, 5 żetonów ruchu, 27 żetonów wymiarów (+ dwa nadpisania kontekstowe), 7 żetonów siatki i punktów łamania, 11 warstw `z-index`, 2 gradienty, po 43 żetony semantyczne w każdym z dwóch motywów, komplet 33 pomiarów kontrastu, procedura dodania żetonu |
| **Czego NIE zawiera** | specyfikacji komponentów `.dn-*` (opracowanie osobne), makiet okien, katalogu ikon, księgi znaku, wytycznych redakcyjnych treści, wartości warstwy funkcjonalnej Danaco Pilot v1.0 (zastąpionych niniejszym pakietem) |

**Źródła wiążące:** `WYNIK/KANON.md` · `zasoby/zetony/zetony.css` (448 linii) · `zasoby/zetony/zetony.json` · `zasoby/zetony/kontrasty.json` (33 pomiary) · `zasoby/css/fundament.css` (156 linii) · `zasoby/css/komponenty.css` (1207 linii) · `design/opracowania/design/01-kierunek/KIERUNEK.md` · `zasoby/ikony/manifest.json` · `brief/INWENTARZ-KOMPONENTOW.md` · `brief/INWENTARZ-OKIEN.md`.

---

## Spis treści

1. [Czym jest żeton — trzy warstwy i reguła nieprzekraczalna](#1-czym-jest-żeton--trzy-warstwy-i-reguła-nieprzekraczalna)
2. [Konwencja nazewnicza `--dn-<grupa>-<wariant>`](#2-konwencja-nazewnicza---dn-grupa-wariant)
3. [Prymitywy — skala szarości (18 stopni)](#3-prymitywy--skala-szarości-18-stopni)
4. [Prymitywy — sygnał (8 stopni)](#4-prymitywy--sygnał-8-stopni)
5. [Prymitywy — stany (zieleń · bursztyn · czerwień, po 4 stopnie)](#5-prymitywy--stany-zieleń--bursztyn--czerwień-po-4-stopnie)
6. [Rama kokpitu — 5 żetonów stałych w obu motywach](#6-rama-kokpitu--5-żetonów-stałych-w-obu-motywach)
7. [Typografia — 22 żetony](#7-typografia--22-żetony)
8. [Przestrzeń — 11 odstępów i 6 promieni](#8-przestrzeń--11-odstępów-i-6-promieni)
9. [Ruch — krzywa i cztery czasy](#9-ruch--krzywa-i-cztery-czasy)
10. [Wymiary — 27 żetonów, warianty dotykowy i przestronny](#10-wymiary--27-żetonów-warianty-dotykowy-i-przestronny)
11. [Siatka i punkty łamania](#11-siatka-i-punkty-łamania)
12. [Warstwy — jedenaście poziomów `z-index`](#12-warstwy--jedenaście-poziomów-z-index)
13. [Gradienty — dwa, zakres i zakazy](#13-gradienty--dwa-zakres-i-zakazy)
14. [Motyw jasny — pełna tabela żetonów semantycznych](#14-motyw-jasny--pełna-tabela-żetonów-semantycznych)
15. [Motyw ciemny — pełna tabela żetonów semantycznych](#15-motyw-ciemny--pełna-tabela-żetonów-semantycznych)
16. [Tabela kontrastów — komplet 33 pomiarów](#16-tabela-kontrastów--komplet-33-pomiarów)
17. [Żetony a JSON — jak czytać `zetony.json`](#17-żetony-a-json--jak-czytać-zetonyjson)
18. [Jak dodać żeton — procedura](#18-jak-dodać-żeton--procedura)
19. [Decyzje projektowe](#19-decyzje-projektowe)

---

## 1. Czym jest żeton — trzy warstwy i reguła nieprzekraczalna

### 1.1. Definicja robocza

**Żeton** (ang. *design token*) to nazwana, pojedyncza decyzja wizualna zapisana raz i przywoływana wszędzie. W Danaco Console żeton jest **zmienną własną CSS** o prefiksie `--dn-`, zadeklarowaną w `zasoby/zetony/zetony.css` i powieloną maszynowo w `zasoby/zetony/zetony.json`.

Żeton nie jest zmienną pomocniczą programisty. Jest **jednostką umowy** między projektem a wdrożeniem: zmiana wartości żetonu jest zmianą systemu, nie zmianą pliku.

> Zasada z KANON-u, rozdz. 2: *nigdy nie wpisuj wartości szesnastkowej wprost. Zawsze `var(--dn-*)`.*

### 1.2. Trzy warstwy modelu

```
   ┌───────────────────────────────────────────────────────────────┐
   │ WARSTWA 1 — PRYMITYWY                                         │
   │ surowe skale: --dn-szary-500, --dn-sygnal-600, --dn-zielen-400│
   │ 38 barw · nie mają znaczenia, mają tylko wartość              │
   │ ZAKAZ użycia wprost w komponentach                            │
   └───────────────────────────┬───────────────────────────────────┘
                               │ mapowanie per motyw
   ┌───────────────────────────▼───────────────────────────────────┐
   │ WARSTWA 2 — SEMANTYCZNE (rola)                                │
   │ --dn-tlo, --dn-tekst, --dn-obrys, --dn-sygnal, --dn-blad-tlo  │
   │ 43 nazwy × 2 motywy · niosą znaczenie, nie wartość            │
   │ JEDYNE wolno wołać z poziomu komponentu                       │
   └───────────────────────────┬───────────────────────────────────┘
                               │ użycie
   ┌───────────────────────────▼───────────────────────────────────┐
   │ WARSTWA 3 — NIEZALEŻNE OD MOTYWU                              │
   │ typografia · przestrzeń · ruch · wymiary · siatka · warstwy   │
   │ 134 nazwy · te same w obu motywach                            │
   └───────────────────────────────────────────────────────────────┘
```

Nazewnictwo warstw pochodzi z nagłówka `zetony.css`: *„1. prymitywy — surowe skale barw (ZAKAZ użycia wprost w komponentach) · 2. semantyczne — role per motyw (`data-theme="light|dark"`, oba równoprawne) · 3. niezależne od motywu — typografia, przestrzeń, ruch, wymiary, warstwy"*.

### 1.3. Rachunek żetonów

| Warstwa / grupa | Liczba nazw | Uwaga |
|---|---:|---|
| Prymitywy — szarość | 18 | skala od `#FFFFFF` do `#0A0A0A` |
| Prymitywy — sygnał | 8 | jedyna rodzina akcentu |
| Prymitywy — stany | 12 | 3 rodziny × 4 stopnie |
| Rama kokpitu | 5 | stałe w obu motywach |
| Typografia | 22 | 3 kroje + 9 stopni + 4 wagi + 3 interlinie + 3 odstępy liter |
| Przestrzeń | 17 | 11 odstępów + 6 promieni |
| Ruch | 5 | 1 krzywa + 4 czasy |
| Wymiary | 27 | 25 `--dn-wym-*` + 2 `--dn-odstep-*` |
| Siatka i łamanie | 7 | 4 punkty łamania + max treści + kolumny + przerwa |
| Warstwy `z-index` | 11 | od podłogi do Centrum poleceń |
| Gradienty | 2 | wyłącznie ilustracyjne |
| **Razem niezależne od motywu** | **134** | deklarowane raz, w bloku `:root` |
| Semantyczne — motyw jasny | 43 | blok `:root[data-theme='light']` |
| Semantyczne — motyw ciemny | 43 | blok `:root[data-theme='dark']` |
| **Razem unikalnych nazw** | **177** | 134 + 43 (nazwy semantyczne są wspólne obu motywom) |

Liczba **deklaracji** w pliku jest wyższa (298), ponieważ komplet semantyczny jest świadomie powielony w blokach `@media (prefers-color-scheme: …)` dla przypadku braku jawnego wyboru motywu. Komentarz w `zetony.css` nazywa to wprost: *„powielenie jest świadome (mechanizm kaskady, nie drugie źródło prawdy)"*.

### 1.4. Reguła nieprzekraczalna

> **Komponent sięga wyłącznie po żetony semantyczne. Nigdy po prymitywy.**

Uzasadnienie: prymityw nie wie, w jakim motywie się znajdzie. `--dn-szary-500` w motywie jasnym jest tekstem metadanych, a w ciemnym — również tekstem metadanych, ale to zbieg okoliczności, nie reguła. `--dn-szary-900` w motywie jasnym jest **tekstem**, a w ciemnym **powierzchnią**. Komponent, który sięgnie po prymityw, przestaje działać po przełączeniu motywu — i nie zgłosi tego błędem.

**Wyjątki jawne — sześć wystąpień w całej bibliotece komponentów (1207 linii):**

| # | Miejsce | Żeton prymitywny | Dlaczego wyjątek jest zasadny |
|---|---|---|---|
| 1 | `.dn-btn--sygnal` | `--dn-szary-0` | tekst leży na `--dn-sygnal-wypelnienie`; wypełnienie jest ciemnym błękitem w obu motywach, więc biel jest właściwa niezależnie od motywu (pomiary 12 i 33 w `kontrasty.json`) |
| 2 | `.dn-btn-ikona--na-ramie[aria-pressed='true']` | `--dn-sygnal-300` | rama kokpitu jest zawsze atramentowa; sygnał na niej musi być jasnym stopniem także w motywie jasnym |
| 3 | `.dn-przelacznik:checked::before` | `--dn-szary-0` | suwak leży na torze `--dn-sygnal-wypelnienie` |
| 4 | `.dn-awatar` | `--dn-szary-100` | inicjał leży na `--dn-grad-atrament` — gradient stały w obu motywach |
| 5 | `.dn-awatar--inteligencja` | `--dn-szary-0` | znak leży na `--dn-grad-sygnal` |
| 6 | `.dn-aod-rdzen` | `--dn-szary-0` | rdzeń Always On Display leży na `--dn-grad-sygnal` |

**Wspólny mianownik wszystkich sześciu:** element leży na powierzchni, która **nie przełącza się z motywem** (wypełnienie sygnałowe, rama, gradient). Wyjątek jest dopuszczalny wyłącznie w tej sytuacji. Poza nią — nie ma wyjątku.

### 1.5. Reguła kierunku odwołań

```
prymityw ──► semantyczny ──► komponent ──► wzorzec złożony ──► okno operacyjne
     ▲                                                              │
     └──────────────────── ZAKAZ przeskoku ─────────────────────────┘
```

Każda warstwa odwołuje się **wyłącznie do warstwy bezpośrednio niższej**. Reguła jest zbieżna z modelem siedmiu warstw systemu opisanym w `brief/INWENTARZ-KOMPONENTOW.md` (rozdz. 3): *prymitywy → tokeny semantyczne → komponenty zdefiniowane `.dn-*` → wzorce złożone → okno operacyjne → moduł → środowisko*.

---

## 2. Konwencja nazewnicza `--dn-<grupa>-<wariant>`

### 2.1. Budowa nazwy

```
--dn-sygnal-600
 │   │      └── wariant: stopień skali, rola szczegółowa albo modyfikator
 │   └───────── grupa: rodzina żetonów
 └───────────── prefiks systemu (Danaco), niezmienny
```

Nazwy są **polskie i opisowe** — spójnie z konwencją nazw ikon (`zasoby/ikony/manifest.json`) i klas komponentów. Wyjątki od polszczyzny to skróty techniczne o ustalonym znaczeniu (`fs`, `fw`, `lh`, `ls`, `bp`, `r`, `od`, `z`, `wym`, `grad`, `ff`).

### 2.2. Pełna tabela prefiksów

| Grupa | Wzór nazwy | Liczba | Warstwa | Znaczenie | Przykład |
|---|---|---:|---|---|---|
| `szary` | `--dn-szary-<0…1000>` | 18 | prymityw | neutralna skala achromatyczna | `--dn-szary-925` |
| `sygnal` | `--dn-sygnal-<100…800>` | 8 | prymityw | rodzina błękitu sygnałowego | `--dn-sygnal-500` |
| `zielen` | `--dn-zielen-<100…700>` | 4 | prymityw | rodzina sukcesu | `--dn-zielen-400` |
| `bursztyn` | `--dn-bursztyn-<100…700>` | 4 | prymityw | rodzina ostrzeżenia | `--dn-bursztyn-700` |
| `czerwien` | `--dn-czerwien-<100…700>` | 4 | prymityw | rodzina błędu | `--dn-czerwien-200` |
| `rama` | `--dn-rama[-rola]` | 5 | stały | powierzchnia paska górnego, niezmienna | `--dn-rama-tekst-2` |
| `ff` | `--dn-ff-<rola>` | 3 | niezależny | rodzina kroju (*font family*) | `--dn-ff-mono` |
| `fs` | `--dn-fs-<stopień>` | 9 | niezależny | stopień pisma (*font size*) | `--dn-fs-base` |
| `fw` | `--dn-fw-<waga>` | 4 | niezależny | waga kroju (*font weight*) | `--dn-fw-polgruba` |
| `lh` | `--dn-lh-<gęstość>` | 3 | niezależny | interlinia (*line height*) | `--dn-lh-ciasny` |
| `ls` | `--dn-ls-<zastosowanie>` | 3 | niezależny | odstęp liter (*letter spacing*) | `--dn-ls-wersaliki` |
| `od` | `--dn-od-<krotność>` | 11 | niezależny | odstęp = krotność 4 px | `--dn-od-6` |
| `r` | `--dn-r-<rozmiar>` | 6 | niezależny | promień narożnika | `--dn-r-pill` |
| `ease` | `--dn-ease` | 1 | niezależny | krzywa czasowa — jedna w systemie | `--dn-ease` |
| `czas` | `--dn-czas-<1\|2\|3\|tetno>` | 4 | niezależny | czas trwania przejścia | `--dn-czas-3` |
| `wym` | `--dn-wym-<element>` | 25 | niezależny | wymiar elementu interfejsu | `--dn-wym-kontrolka` |
| `odstep` | `--dn-odstep-<zakres>` | 2 | niezależny | rytm kompozycyjny (alias na `od`) | `--dn-odstep-panel` |
| `bp` | `--dn-bp-w<1…4>` | 4 | niezależny | punkt łamania (*breakpoint*) | `--dn-bp-w3` |
| `tresc` | `--dn-tresc-max` | 1 | niezależny | maksymalna szerokość treści | `--dn-tresc-max` |
| `siatka` | `--dn-siatka-<cecha>` | 2 | niezależny | kolumny i przerwa siatki | `--dn-siatka-przerwa` |
| `z` | `--dn-z-<warstwa>` | 11 | niezależny | poziom `z-index` | `--dn-z-modal` |
| `grad` | `--dn-grad-<rodzina>` | 2 | niezależny | gradient ilustracyjny | `--dn-grad-sygnal` |
| `tlo`, `powierzchnia`, `panel` | nazwa roli | 4 | semantyczny | powierzchnie dokumentu | `--dn-powierzchnia-2` |
| `hover`, `wcisniecie`, `nakladka` | nazwa roli | 3 | semantyczny | nakładki interakcji | `--dn-wcisniecie` |
| `obrys` | `--dn-obrys[-moc]` | 3 | semantyczny | linie rozdzielające | `--dn-obrys-mocny` |
| `tekst` | `--dn-tekst[-stopień\|-inv]` | 4 | semantyczny | hierarchia tekstu | `--dn-tekst-2` |
| `atrament` | `--dn-atrament[-rola]` | 3 | semantyczny | działanie główne (inwersja) | `--dn-atrament-tekst` |
| `sygnal` (sem.) | `--dn-sygnal[-rola]` | 6 | semantyczny | role akcentu | `--dn-sygnal-wypelnienie` |
| `kropka`, `fokus` | nazwa roli | 3 | semantyczny | element sygnaturowy i fokus | `--dn-fokus-cien` |
| `sukces`, `ostrzezenie`, `blad`, `informacja` | `--dn-<stan>-<tekst\|tlo\|obrys>` | 12 | semantyczny | komplet stanu | `--dn-blad-obrys` |
| `cien` | `--dn-cien-<1\|2\|3\|lg\|sygnal>` | 5 | semantyczny | głębia per motyw | `--dn-cien-lg` |

### 2.3. Reguły nazewnicze — co wolno, czego nie wolno

| Reguła | Treść |
|---|---|
| **R1** | Prefiks `--dn-` jest obowiązkowy i niezmienny. Żeton bez tego prefiksu nie należy do systemu. |
| **R2** | Nazwa grupy jest rzeczownikiem albo skrótem technicznym o ustalonym znaczeniu (tabela 2.2). Nowa grupa wymaga decyzji prowadzącego system. |
| **R3** | Skale liczbowe rosną wraz z ciemnością (`szary-0` = biel → `szary-1000` = niemal czerń) oraz z nasyceniem (`sygnal-100` najjaśniejszy → `sygnal-800` najciemniejszy). Kierunek jest jednakowy we wszystkich rodzinach barwnych. |
| **R4** | Żeton semantyczny **nie zawiera liczby stopnia**. `--dn-tekst-2` to druga rola tekstu, nie drugi stopień skali. |
| **R5** | Sufiks `-hover` oznacza wariant przy najechaniu tej samej roli (`--dn-atrament-hover`, `--dn-sygnal-wypelnienie-hover`, `--dn-rama-hover`). |
| **R6** | Sufiks `-inv` oznacza inwersję względem tła bieżącego motywu (`--dn-tekst-inv`). |
| **R7** | Żeton stanu ma zawsze trójkę `-tekst` / `-tlo` / `-obrys`. Nie wolno wprowadzić stanu z niepełną trójką. |
| **R8** | Nie tworzy się żetonów jednorazowych „pod jeden komponent". Jeżeli wartość występuje raz — należy do komponentu, nie do systemu. |
| **R9** | Nie tworzy się aliasów semantycznych dla tej samej roli (`--dn-tekst-glowny` obok `--dn-tekst` jest zabroniony). |
| **R10** | Wartość żetonu semantycznego jest zawsze `var(--dn-<prymityw>)` albo `rgba(...)` wyprowadzoną z prymitywu. Nigdy nowy, niezadeklarowany hex. |

---

## 3. Prymitywy — skala szarości (18 stopni)

Skala jest **czysto neutralna**: wszystkie trzy składowe RGB są równe. KIERUNEK.md, rozdz. 3.1: *„Żadnego podbarwienia slate/niebieskiego ani ciepłego beżu — »odcienie bieli i czerni« dosłownie."*

| Żeton | Wartość | Kontrast do `#F4F4F4` (tło jasne) | Kontrast do `#0F0F0F` (tło ciemne) | Zastosowanie systemowe | Zakazy |
|---|---|---:|---:|---|---|
| `--dn-szary-0` | `#FFFFFF` | 1,10 | 19,17 | jasny: `--dn-powierzchnia`, `--dn-tekst-inv`, `--dn-atrament-tekst`; biel na wypełnieniu sygnałowym i na gradientach | **nigdy jako tło całej strony** (KIERUNEK 3.1) |
| `--dn-szary-25` | `#FAFAFA` | 1,05 | 18,36 | jasny: `--dn-panel`; ciemny: `--dn-atrament-hover` | nie na tekst w motywie jasnym (1,05 wobec tła) |
| `--dn-szary-50` | `#F4F4F4` | 1,00 | 17,43 | jasny: `--dn-tlo` — papier roboczy kokpitu | nie jako powierzchnia karty (znika na tle) |
| `--dn-szary-100` | `#ECECEC` | 1,07 | 16,23 | jasny: `--dn-powierzchnia-2`, `--dn-obrys-subtelny`; ciemny: `--dn-tekst`, `--dn-atrament`; rama: `--dn-rama-tekst` | nie jako obrys widoczny w motywie jasnym (1,18 do bieli) |
| `--dn-szary-150` | `#E3E3E3` | 1,17 | 14,94 | jasny: `--dn-obrys` — podstawowa linia rozdzielająca | nie jako tekst w żadnym motywie jasnym |
| `--dn-szary-200` | `#D7D7D7` | 1,31 | 13,32 | **stopień rezerwowy** — brak mapowania semantycznego; utrzymywany dla ciągłości skali i wykresów danych | nie wolno użyć wprost w komponencie; wprowadzenie do użycia wymaga nadania roli semantycznej |
| `--dn-szary-300` | `#C0C0C0` | 1,65 | 10,54 | jasny: `--dn-obrys-mocny` — obrysy kontrolek, kciuk paska przewijania | nie jako tekst (1,65 wobec tła) |
| `--dn-szary-400` | `#9E9E9E` | 2,44 | 7,15 | ciemny: `--dn-tekst-2`; rama: `--dn-rama-tekst-2` | **nie jako tekst w motywie jasnym** (2,44 — poniżej progu) |
| `--dn-szary-500` | `#7C7C7C` | 3,80 | 4,59 | oba motywy: `--dn-tekst-3` | **wyłącznie metadane i tekst ≥ 18,66 px półgruby** (KANON, rozdz. 9); nigdy tekst ciągły |
| `--dn-szary-600` | `#616161` | 5,63 | 3,10 | jasny: `--dn-tekst-2` — tekst drugorzędny | nie jako tekst w motywie ciemnym (3,10) |
| `--dn-szary-700` | `#4A4A4A` | 8,06 | 2,16 | **stopień rezerwowy** — brak mapowania semantycznego; przewidziany dla obrysów o wysokim kontraście na powierzchni jasnej | nie wolno użyć wprost w komponencie |
| `--dn-szary-750` | `#3A3A3A` | 10,34 | 1,69 | ciemny: `--dn-obrys-mocny` — obrysy kontrolek, kciuk paska przewijania | nie jako powierzchnia (za jasny na tło ciemne) |
| `--dn-szary-800` | `#2A2A2A` | 13,05 | 1,34 | ciemny: `--dn-obrys`; punkt startowy `--dn-grad-atrament` | nie jako powierzchnia panelu (kolizja z `--dn-panel` #212121) |
| `--dn-szary-850` | `#212121` | 14,64 | 1,19 | ciemny: `--dn-panel`, `--dn-obrys-subtelny` | nie jako tło strony |
| `--dn-szary-900` | `#181818` | 16,14 | 1,08 | jasny: `--dn-tekst`, `--dn-atrament`; ciemny: `--dn-powierzchnia` | najlepszy przykład odwrócenia roli między motywami — nie wolno traktować jako „koloru tekstu" |
| `--dn-szary-925` | `#131313` | 16,89 | 1,03 | **`--dn-rama` w obu motywach**; ciemny: `--dn-powierzchnia-2`; punkt końcowy `--dn-grad-atrament` | nie jako tło treści przełączalnej |
| `--dn-szary-950` | `#0F0F0F` | 17,43 | 1,00 | ciemny: `--dn-tlo`, `--dn-tekst-inv`, `--dn-atrament-tekst` | nie jako tekst w motywie ciemnym |
| `--dn-szary-1000` | `#0A0A0A` | 18,00 | 1,03 | jasny: `--dn-atrament-hover`; źródło alfy dla `--dn-hover`, `--dn-wcisniecie`, `--dn-nakladka` i wszystkich cieni motywu jasnego | **kraniec skali — `#000000` jest zakazane** (KIERUNEK 3.1: „czysta czerń zabija głębię i powoduje halację na OLED") |

**Zakaz nadrzędny skali:** `#000000` nie występuje w systemie w żadnej postaci barwy powierzchni ani tekstu. Skala kończy się na `#0A0A0A`.

**Dwa stopnie rezerwowe** (`-200`, `-700`) nie mają dziś mapowania semantycznego. Nie są błędem: utrzymują równomierność skali, dzięki czemu każdy przyszły żeton semantyczny znajdzie stopień o właściwej jasności bez dokładania nowej wartości.

---

## 4. Prymitywy — sygnał (8 stopni)

Sygnał jest **jedyną barwą akcentu w systemie**. KIERUNEK.md 3.1: *„Sygnał nie jest barwą przycisków głównych — czerń działa, sygnał wskazuje."* Rodzina informacyjna stanów jest **scalona z rodziną sygnału** — to decyzja celowa, nie oszczędność.

| Żeton | Wartość | / `#F4F4F4` | / `#FFFFFF` | / `#0F0F0F` | / `#181818` | Rola systemowa | Zakazy |
|---|---|---:|---:|---:|---:|---|---|
| `--dn-sygnal-100` | `#EDF3FE` | 1,01 | 1,11 | 17,21 | 15,94 | jasny: `--dn-sygnal-tlo`, `--dn-informacja-tlo` | nie jako tekst; nie w motywie ciemnym (tło informacyjne ciemne buduje się alfą) |
| `--dn-sygnal-200` | `#C9DAFB` | 1,28 | 1,41 | 13,61 | 12,60 | jasny: `--dn-sygnal-obrys`, `--dn-informacja-obrys`; ciemny: `--dn-sygnal-mocny` | nie jako wypełnienie przycisku |
| `--dn-sygnal-300` | `#8FB2F5` | 1,94 | 2,13 | 9,00 | 8,33 | ciemny: `--dn-sygnal`, `--dn-informacja-tekst`; sygnał na ramie (`.dn-btn-ikona--na-ramie`) | **nie jako tekst w motywie jasnym** (1,94) |
| `--dn-sygnal-400` | `#5C8CEC` | 2,98 | 3,27 | 5,86 | 5,43 | ciemny: `--dn-kropka`, `--dn-fokus`, `--dn-sygnal-wypelnienie-hover`; punkt startowy `--dn-grad-sygnal` | nie jako tekst w motywie jasnym |
| `--dn-sygnal-500` | `#3B6FE0` | 4,21 | 4,63 | 4,14 | 3,84 | **bazowy** — jasny: `--dn-kropka`, `--dn-fokus`; ciemny: `--dn-sygnal-wypelnienie`; źródło alfy `--dn-fokus-cien` (jasny), `--dn-sygnal-tlo` (ciemny) | nie jako tekst na powierzchni jasnej w rozmiarze bazowym (4,63 — ledwie ponad progiem; do tekstu służy `-600`) |
| `--dn-sygnal-600` | `#2457C9` | 5,82 | 6,40 | 3,00 | 2,77 | jasny: `--dn-sygnal`, `--dn-sygnal-wypelnienie`, `--dn-informacja-tekst`; punkt końcowy `--dn-grad-sygnal` | **nie jako tekst w motywie ciemnym** (2,77) |
| `--dn-sygnal-700` | `#1D49AF` | 7,28 | 8,00 | 2,40 | 2,22 | jasny: `--dn-sygnal-mocny`, `--dn-sygnal-wypelnienie-hover` | nie w motywie ciemnym w żadnej roli tekstowej |
| `--dn-sygnal-800` | `#173A8C` | 9,45 | 10,39 | 1,84 | 1,71 | **stopień rezerwowy** — najciemniejszy stopień rodziny; przewidziany dla stanu wciśnięcia wypełnienia sygnałowego | nie wolno użyć wprost w komponencie |

### 4.1. Reguła budżetu akcentu

Sygnał ma **zamknięty katalog wystąpień**: fokus, stan aktywny/wybrany, odnośnik, postęp, wskaźnik pracy w tle (kropka sygnału), tło i obrys komunikatu informacyjnego, przycisk `.dn-btn--sygnal` (jeden na widok), medalion jednostki inteligencji, rdzeń Always On Display. Poza tym katalogiem sygnał nie występuje.

### 4.2. Sygnał a inwersja atramentu

```
   działanie główne                 wskazanie
   ────────────────                 ─────────
   --dn-atrament                    --dn-sygnal
   jasny:  #181818 (czerń)          jasny:  #2457C9
   ciemny: #ECECEC (biel)           ciemny: #8FB2F5
   „to jest przycisk, kliknij"      „tu biegnie praca / tu jesteś"
```

Monochromatyczna inwersja jest najsilniejszym akcentem systemu i **nie wprowadza drugiej barwy**. Sygnał pozostaje wolny dla informacji o stanie.

---

## 5. Prymitywy — stany (zieleń · bursztyn · czerwień, po 4 stopnie)

Trzy rodziny barwne po cztery stopnie. Czwarty stan systemu — **informacja** — nie ma własnej rodziny: korzysta z rodziny sygnału (`zetony.css`, rozdz. 1: *„STANY — prymitywy (informacja = rodzina sygnału)"*).

### 5.1. Zieleń — sukces

| Żeton | Wartość | / `#FFFFFF` | / `#181818` | Rola | Zakazy |
|---|---|---:|---:|---|---|
| `--dn-zielen-700` | `#1E7A46` | 5,35 | 3,32 | jasny: `--dn-sukces-tekst` | nie jako tekst w motywie ciemnym (3,32) |
| `--dn-zielen-400` | `#4FBD85` | 2,35 | 7,57 | ciemny: `--dn-sukces-tekst`; źródło alfy `--dn-sukces-tlo` (0,14) i `--dn-sukces-obrys` (0,34) | **nie jako tekst w motywie jasnym** (2,35) |
| `--dn-zielen-200` | `#BCE3CD` | 1,40 | 12,69 | jasny: `--dn-sukces-obrys` | nie jako tekst |
| `--dn-zielen-100` | `#E7F5EE` | 1,12 | 15,81 | jasny: `--dn-sukces-tlo` | nie jako tło w motywie ciemnym (tam tło buduje alfa) |

### 5.2. Bursztyn — ostrzeżenie

| Żeton | Wartość | / `#FFFFFF` | / `#181818` | Rola | Zakazy |
|---|---|---:|---:|---|---|
| `--dn-bursztyn-700` | `#8A5B0C` | 5,86 | 3,03 | jasny: `--dn-ostrzezenie-tekst` | nie jako tekst w motywie ciemnym (3,03) |
| `--dn-bursztyn-400` | `#E0A63C` | 2,17 | 8,19 | ciemny: `--dn-ostrzezenie-tekst`; źródło alfy tła (0,14) i obrysu (0,34) | **nie jako tekst w motywie jasnym** (2,17) |
| `--dn-bursztyn-200` | `#EED9A7` | 1,39 | 12,78 | jasny: `--dn-ostrzezenie-obrys` | nie jako tekst |
| `--dn-bursztyn-100` | `#FBF2DE` | 1,11 | 15,94 | jasny: `--dn-ostrzezenie-tlo` | nie jako tło w motywie ciemnym |

### 5.3. Czerwień — błąd

| Żeton | Wartość | / `#FFFFFF` | / `#181818` | Rola | Zakazy |
|---|---|---:|---:|---|---|
| `--dn-czerwien-700` | `#C0362F` | 5,51 | 3,22 | jasny: `--dn-blad-tekst` | nie jako tekst w motywie ciemnym (3,22) |
| `--dn-czerwien-400` | `#EE7168` | 2,92 | 6,09 | ciemny: `--dn-blad-tekst`; źródło alfy tła (0,14) i obrysu (0,34) | **nie jako tekst w motywie jasnym** (2,92) |
| `--dn-czerwien-200` | `#F3C8C4` | 1,51 | 11,73 | jasny: `--dn-blad-obrys` | nie jako tekst |
| `--dn-czerwien-100` | `#FCEDEB` | 1,14 | 15,60 | jasny: `--dn-blad-tlo` | nie jako tło w motywie ciemnym |

### 5.4. Reguła nadrzędna stanów

> **Stan nigdy samym kolorem — zawsze ikona albo etykieta.** (KANON rozdz. 1 i 9; `zetony.json`: *„stan nigdy samym kolorem — zawsze ikona albo etykieta [NIENEGOCJOWALNE]"*.)

Barwa stanu jest wzmocnieniem, nie nośnikiem. Komplet ikon towarzyszących pochodzi z zestawu: `ptaszek-kolo` (sukces), `ostrzezenie` (ostrzeżenie), `blad` (błąd), `info` (informacja).

**Symetria trzech rodzin** — każda ma dokładnie ten sam układ: stopień `-700` (tekst na jasnym), `-400` (tekst na ciemnym oraz źródło alfy), `-200` (obrys na jasnym), `-100` (tło na jasnym). Nowa rodzina stanu musiałaby powtórzyć ten układ w całości (reguła R7).

---

## 6. Rama kokpitu — 5 żetonów stałych w obu motywach

KIERUNEK.md 3.3: *„Pasek górny jest zawsze atramentowy (`szary-925`) — w obu motywach. Ciemna rama wokół przełączalnej treści daje efekt stanowiska dowodzenia i stały dom dla godła, wyszukiwarki i wskaźników."*

| Żeton | Wartość | Mapuje się na | Kontrast na ramie | Rola |
|---|---|---|---:|---|
| `--dn-rama` | `#131313` | `var(--dn-szary-925)` | — | powierzchnia paska górnego; **jedyny element interfejsu nieprzełączający się z motywem** (poza godłem) |
| `--dn-rama-tekst` | `#ECECEC` | `var(--dn-szary-100)` | 15,73 | logotyp, etykiety główne, ikona aktywna |
| `--dn-rama-tekst-2` | `#9E9E9E` | `var(--dn-szary-400)` | 6,94 | ikony w spoczynku, tekst pomocniczy, podpowiedź w polu wyszukiwania |
| `--dn-rama-hover` | `rgba(255,255,255,0.08)` | biel z kryciem 8% | — | tło przycisku ikonowego przy najechaniu i w stanie wciśniętym |
| `--dn-rama-obrys` | `rgba(255,255,255,0.10)` | biel z kryciem 10% | — | obrys pola wyszukiwania i separatorów na ramie |

**Konsekwencje projektowe:**

1. Komponent położony na ramie **nie może** używać `--dn-tekst`, `--dn-tekst-2`, `--dn-hover` ani `--dn-obrys` — te przełączają się z motywem i w motywie jasnym zniknęłyby na atramencie. Służy do tego wariant `.dn-btn-ikona--na-ramie`.
2. Wskazanie stanu aktywnego na ramie używa `--dn-sygnal-300` (`#8FB2F5`, 8,72 na ramie), a nie `--dn-sygnal` — bo `--dn-sygnal` w motywie jasnym to `#2457C9` (2,90 na ramie, poniżej progu).
3. Rama ma tylko dwa stopnie tekstu, nie trzy. Odpowiednika `--dn-tekst-3` na ramie **nie ma** i nie wolno go improwizować kryciem.

```
 ┌───────────────────────────────────────────────────────────┐
 │  RAMA #131313 — stała                                     │  ← --dn-rama
 │  ● godło   DANACO CONSOLE   [szukaj]  [motyw] [ustawienia]│
 ├───────────────────────────────────────────────────────────┤
 │                                                           │
 │      TREŚĆ PRZEŁĄCZALNA — --dn-tlo                        │
 │      jasny #F4F4F4   ⇄   ciemny #0F0F0F                   │
 │                                                           │
 └───────────────────────────────────────────────────────────┘
```

---

## 7. Typografia — 22 żetony

### 7.1. Trzy kroje, trzy role

| Żeton | Wartość (stos) | Rola | Wagi w zestawie | Uzasadnienie (KIERUNEK 3.2) |
|---|---|---|---|---|
| `--dn-ff-naglowek` | `'Space Grotesk', 'IBM Plex Sans', 'Segoe UI', system-ui, sans-serif` | nagłówki, tytuły środowisk, logotyp | 500 · 600 · 700 | geometryczno-techniczny charakter; wyraźna osobowość bez ozdobności; wersaliki do kokpitu |
| `--dn-ff-bazowa` | `'IBM Plex Sans', 'Segoe UI', system-ui, -apple-system, sans-serif` | interfejs i treść | 400 · 500 · 600 · 700 | krój projektowany do środowisk inżynierskich; czytelność w małych stopniach; wzorowe polskie diakrytyki |
| `--dn-ff-mono` | `'IBM Plex Mono', 'Cascadia Mono', 'Consolas', monospace` | dane, identyfikatory, terminal | 400 · 500 · 600 | naturalna para dla Plex Sans; `tabular-nums`; DNA konsoli |

Znaczenie różnicy krojów: **Space Grotesk = wejście do środowiska i tożsamość marki · Plex Sans = praca · Plex Mono = maszyna.** Licencja obu rodzin: OFL 1.1, podzbiory `latin` + `latin-ext` (pełne polskie diakrytyki). Pliki: `zasoby/fonty/*.woff2` (22 pliki), deklaracje `@font-face` w `zasoby/zetony/fonty.css`.

**Test diakrytyków obowiązujący dla każdego kroju i wagi:** `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`.

### 7.2. Dziewięć stopni pisma

| Żeton | Wartość | Typowe zastosowanie (komentarz źródłowy) | Krój domyślny |
|---|---:|---|---|
| `--dn-fs-xs` | 11 px | etykiety wersalikowe, nagłówki kolumn tabeli, plakietki | Plex Sans / Plex Mono |
| `--dn-fs-sm` | 12 px | metadane, tekst pomocniczy, opis pola, godzina wpisu | Plex Sans |
| `--dn-fs-base` | **13 px** | **BAZOWY** — tekst interfejsu, etykieta przycisku, pole formularza | Plex Sans |
| `--dn-fs-md` | 14 px | wyróżniona treść, wpisy rozmowy w Chat Window | Plex Sans |
| `--dn-fs-lg` | 16 px | nagłówki paneli, tytuł karty | Space Grotesk |
| `--dn-fs-xl` | 20 px | nagłówki okien i modali | Space Grotesk |
| `--dn-fs-2xl` | 24 px | tytuły sekcji strony głównej (Centrum dowodzenia) | Space Grotesk |
| `--dn-fs-3xl` | 30 px | tytuły kart środowisk (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) | Space Grotesk |
| `--dn-fs-display` | 40 px | największy stopień ekspozycyjny | Space Grotesk |

**Kształt skali:** 11 → 12 → 13 → 14 (przyrost 1 px, strefa robocza kokpitu) → 16 → 20 → 24 → 30 → 40 (przyrost mnożnikowy ≈ 1,25, strefa ekspozycyjna). Skala jest celowo dwuczęściowa: gęstość zwarta wymaga drobnych, precyzyjnych różnic w strefie danych i wyraźnych skoków w strefie nagłówków.

### 7.3. Cztery wagi

| Żeton | Wartość | Zastosowanie |
|---|---:|---|
| `--dn-fw-normalna` | 400 | tekst ciągły, treść wpisu, wartości w tabeli |
| `--dn-fw-srednia` | 500 | etykieta przycisku, pozycja nawigacji, logotyp `CONSOLE` |
| `--dn-fw-polgruba` | 600 | nagłówki `h1`–`h3` (`fundament.css`), tytuł karty, etykieta wersalikowa, nadawca wpisu |
| `--dn-fw-gruba` | 700 | logotyp `DANACO`, tytuł ekspozycyjny |

### 7.4. Trzy interlinie

| Żeton | Wartość | Zastosowanie |
|---|---:|---|
| `--dn-lh-ciasny` | 1,25 | nagłówki (ustawiane globalnie w `fundament.css` dla `h1, h2, h3, .dn-naglowek`) |
| `--dn-lh-bazowy` | 1,45 | interfejs — wartość `body`; zwarta gęstość |
| `--dn-lh-luzny` | 1,6 | treść ciągła, opisy, akapity dokumentacyjne |

W gęstości przestronnej `--dn-lh-bazowy` rośnie do **1,5** (patrz rozdz. 10.3).

### 7.5. Trzy odstępy liter

| Żeton | Wartość | Zastosowanie |
|---|---:|---|
| `--dn-ls-naglowek` | −0,01 em | Space Grotesk w dużych stopniach — kompensacja rozstrzelenia geometrycznego |
| `--dn-ls-wersaliki` | 0,08 em | etykiety wersalikowe Plex Sans (`.dn-etykieta-wersalikowa`) |
| `--dn-ls-mono-wersaliki` | 0,14 em | etykiety wersalikowe Plex Mono (`.dn-etykieta-mono`), logotyp `CONSOLE` |

### 7.6. Reguła liczb

Wszystkie liczby prezentowane w interfejsie — identyfikatory zadań, czasy, liczniki kolejek, rozmiary plików — składane są krojem **IBM Plex Mono** z `font-variant-numeric: tabular-nums`. Wzorzec jest zamknięty w `fundament.css` klasami `.dn-dane` i `.dn-liczba`, dzięki czemu kolumny liczbowe nie „skaczą" przy aktualizacji na żywo (Execution Monitor, Process Monitor, Queue Manager).

---

## 8. Przestrzeń — 11 odstępów i 6 promieni

### 8.1. Jedenaście odstępów — jednostka 4 px [NIENEGOCJOWALNE]

| Żeton | Wartość | Krotność | Typowe zastosowanie |
|---|---:|---:|---|
| `--dn-od-0` | 0 | 0× | zerowanie odstępu w kompozycji zwartej |
| `--dn-od-1` | 4 px | 1× | odstęp ikona–etykieta, wewnętrzny rytm plakietki |
| `--dn-od-2` | 8 px | 2× | odstęp między kontrolkami w rzędzie, padding pionowy przycisku ikonowego |
| `--dn-od-3` | 12 px | 3× | **rytm wewnętrzny paneli** (`--dn-odstep-panel`), padding poziomy pola |
| `--dn-od-4` | 16 px | 4× | padding karty, odstęp między grupami pól |
| `--dn-od-5` | 20 px | 5× | padding ciała modala i bloku demonstracyjnego |
| `--dn-od-6` | 24 px | 6× | **rytm między sekcjami** (`--dn-odstep-sekcji`), przerwa siatki |
| `--dn-od-8` | 32 px | 8× | odstęp między blokami strony, wcięcie pola wyszukiwania |
| `--dn-od-10` | 40 px | 10× | margines górny stopki, wysokość kontrolki dotykowej |
| `--dn-od-12` | 48 px | 12× | odstęp między sekcjami dokumentu, wysokość paska górnego |
| `--dn-od-16` | 64 px | 16× | margines dolny strony, największy odstęp kompozycyjny |

Brakujące krotności (7×, 9×, 11×, 13×, 14×, 15×) są **celowo pominięte**: skala rzadnie tam, gdzie różnica przestaje być czytelna. Zamiast dodawać stopień — wybiera się sąsiedni.

**Dwa aliasy kompozycyjne** (zaliczone do grupy wymiarów, rozdz. 10): `--dn-odstep-panel` = `var(--dn-od-3)` · `--dn-odstep-sekcji` = `var(--dn-od-6)`. Ich sens: komponent pyta o **rytm**, nie o liczbę pikseli — dzięki temu przełączenie gęstości zmienia oddech wszystkich paneli jednym żetonem.

```
 0    4    8    12    16    20    24        32        40        48              64
 │    │    │    │     │     │     │         │         │         │               │
 od-0 od-1 od-2 od-3  od-4  od-5  od-6      od-8      od-10     od-12           od-16
           ▲          ▲           ▲                             ▲
           │          │           │                             │
      kontrolki   karta      sekcja / siatka              pasek górny
```

### 8.2. Sześć promieni

| Żeton | Wartość | Zastosowanie (komentarz źródłowy) |
|---|---:|---|
| `--dn-r-xs` | 3 px | plakietki drobne, bloki kodu w wierszu, pierścień fokusu |
| `--dn-r-sm` | 6 px | przyciski, pola, kontrolki, karty sesji |
| `--dn-r-md` | 8 px | karty pojedyncze, dymki, toasty |
| `--dn-r-lg` | 10 px | karty, modale, panele |
| `--dn-r-xl` | 14 px | karty środowisk, duże powierzchnie |
| `--dn-r-pill` | 999 px | plakietki stanu, przełączniki, kropki, awatary, tory postępu |

Komentarz źródłowy nazywa zasadę wprost: *„jeden system narożników: precyzja instrumentu (małe promienie)"*. Promień rośnie wraz z powierzchnią elementu — nigdy odwrotnie. Element wewnątrz elementu ma promień **mniejszy albo równy** promieniowi rodzica.

---

## 9. Ruch — krzywa i cztery czasy

Pokrętło `INTENSYWNOSC_RUCHU` = **3/10**. Ruch niesie informację o stanie systemu; nie ma ruchu dekoracyjnego.

| Żeton | Wartość | Przeznaczenie | Przykład zastosowania |
|---|---|---|---|
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | **jedyna krzywa czasowa systemu** — szybki start, miękkie dojście | wszystkie przejścia bez wyjątku |
| `--dn-czas-1` | 0,10 s | mikroreakcje | najechanie na wiersz tabeli, naciśnięcie przycisku |
| `--dn-czas-2` | 0,16 s | przejścia barw | przełączenie motywu (`body` w `fundament.css`), zmiana obrysu pola przy fokusie |
| `--dn-czas-3` | 0,22 s | wejście i wyjście warstw | modal (`dn-wejscie`), panel, toast, zmiana widoku |
| `--dn-czas-tetno` | **2,4 s** | **jedyny ruch ciągły w interfejsie** | tętno kropki sygnału (`.dn-kropka--tetno`, `@keyframes dn-tetno`) |

### 9.1. Wzorzec przejścia widoku

`opacity` + `translateY(4 px)`, maksymalnie 220 ms, krzywa `--dn-ease`. Bez parallaxu, bez scrollytellingu, bez animacji wejścia elementów listy.

### 9.2. Ograniczony ruch — obsługa globalna

`prefers-reduced-motion: reduce` jest obsłużone **w żetonach, nie per komponent** (`zetony.css`, rozdz. 14):

```css
@media (prefers-reduced-motion: reduce) {
  :root { --dn-czas-1: 0.01ms; --dn-czas-2: 0.01ms;
          --dn-czas-3: 0.01ms; --dn-czas-tetno: 0.01ms; }
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}
```

Konsekwencja: **komponent nigdy nie definiuje własnej obsługi ograniczonego ruchu.** Tętno kropki zastępuje wówczas **pierścień statyczny** (`zetony.json`: *„tetno kropki zastepuje pierscien statyczny"*).

---

## 10. Wymiary — 27 żetonów, warianty dotykowy i przestronny

### 10.1. Gęstość zwarta (domyślna) — pełna lista

| # | Żeton | Wartość | Element | Uwaga |
|---:|---|---:|---|---|
| 1 | `--dn-wym-kontrolka` | 32 px | przyciski, pola, pozycje wyboru | podstawa gęstości zwartej |
| 2 | `--dn-wym-ikonowy` | 32 px | przycisk ikonowy (kwadrat) | równy kontrolce — rząd nie „faluje" |
| 3 | `--dn-wym-pasek` | 48 px | pasek górny (rama) | `--dn-od-12` |
| 4 | `--dn-wym-pas-kart` | 36 px | pas kart sesji pod paskiem | |
| 5 | `--dn-wym-wiersz` | 36 px | wiersz tabeli i listy | równy pasowi kart |
| 6 | `--dn-wym-boczna` | 224 px | boczna nawigacja modułów | 56 × 4 px |
| 7 | `--dn-wym-pas-komunikacji` | 320 px | wysokość pasa komunikacji (Chat Window) | |
| 8 | `--dn-wym-modal` | 560 px | maksymalna szerokość modala | |
| 9 | `--dn-wym-toast-min` | 280 px | minimalna szerokość powiadomienia | |
| 10 | `--dn-wym-toast-max` | 420 px | maksymalna szerokość powiadomienia | |
| 11 | `--dn-wym-awatar-sm` | 24 px | awatar mały — wiersz listy | |
| 12 | `--dn-wym-awatar` | 28 px | awatar podstawowy — medalion wpisu | |
| 13 | `--dn-wym-awatar-lg` | 36 px | awatar duży — nagłówek, rdzeń Always On Display | |
| 14 | `--dn-wym-ikona-sm` | 14 px | ikona w plakietce, w medalionie wpisu | |
| 15 | `--dn-wym-ikona` | 16 px | ikona podstawowa — przycisk, nawigacja | |
| 16 | `--dn-wym-ikona-lg` | 20 px | ikona w przycisku ikonowym paska | |
| 17 | `--dn-wym-ikona-xl` | 24 px | ikona w pustym stanie, w kaflu | |
| 18 | `--dn-wym-przelacznik-szer` | 36 px | tor przełącznika | |
| 19 | `--dn-wym-przelacznik-wys` | 20 px | wysokość toru przełącznika | suwak = `wys − 6 px` |
| 20 | `--dn-wym-check` | 16 px | pole wyboru i opcja jednokrotna | |
| 21 | `--dn-wym-kropka` | **6 px** | kropka sygnału — element sygnaturowy | rozstrzygnięcie L-P-16 |
| 22 | `--dn-wym-wstega` | 2 px | wstęga aktywności karty i pozycji nawigacji | |
| 23 | `--dn-wym-spinner` | 14 px | wskaźnik pracy w przycisku | |
| 24 | `--dn-wym-fokus` | 2 px | grubość pierścienia fokusu | |
| 25 | `--dn-wym-fokus-odsuniecie` | 2 px | odsunięcie pierścienia od krawędzi | |
| 26 | `--dn-odstep-panel` | 12 px (`--dn-od-3`) | wewnętrzny rytm paneli | alias kompozycyjny |
| 27 | `--dn-odstep-sekcji` | 24 px (`--dn-od-6`) | rytm między sekcjami | alias kompozycyjny |

### 10.2. Wariant dotykowy — `@media (pointer: coarse)`

Cele dotykowe rosną **żetonem, nie wyjątkiem** (rozstrzygnięcie L-P-14). Sześć nadpisań; reguły komponentów pozostają nietknięte.

| Żeton | Zwarta | Dotyk | Przyrost |
|---|---:|---:|---:|
| `--dn-wym-kontrolka` | 32 px | **40 px** | +8 px |
| `--dn-wym-ikonowy` | 32 px | **40 px** | +8 px |
| `--dn-wym-wiersz` | 36 px | **44 px** | +8 px |
| `--dn-wym-check` | 16 px | **20 px** | +4 px |
| `--dn-wym-przelacznik-szer` | 36 px | **44 px** | +8 px |
| `--dn-wym-przelacznik-wys` | 20 px | **24 px** | +4 px |

### 10.3. Wariant przestronny — `[data-gestosc='przestronna']`

Przygotowany w żetonach, **domyślnie nieaktywny** (rozstrzygnięcie L-P-13). Dziewięć nadpisań, w tym dwa spoza grupy wymiarów (stopień bazowy i interlinia bazowa).

| Żeton | Zwarta | Przestronna | Grupa |
|---|---:|---:|---|
| `--dn-fs-base` | 13 px | **14 px** | typografia |
| `--dn-lh-bazowy` | 1,45 | **1,5** | typografia |
| `--dn-wym-kontrolka` | 32 px | **40 px** | wymiary |
| `--dn-wym-ikonowy` | 32 px | **40 px** | wymiary |
| `--dn-wym-wiersz` | 36 px | **44 px** | wymiary |
| `--dn-wym-pasek` | 48 px | **56 px** | wymiary |
| `--dn-wym-pas-kart` | 36 px | **40 px** | wymiary |
| `--dn-odstep-panel` | 12 px | **20 px** (`--dn-od-5`) | wymiary |
| `--dn-odstep-sekcji` | 24 px | **32 px** (`--dn-od-8`) | wymiary |

**Reguła kompozycji trybów:** dotyk i gęstość przestronna mogą wystąpić jednocześnie. Kolejność kaskady w `zetony.css` sprawia, że blok `[data-gestosc='przestronna']` (rozdz. 13 pliku) rozstrzyga wartości wspólne, ponieważ znajduje się **po** bloku `@media (pointer: coarse)` (rozdz. 12 pliku) przy tej samej specyficzności. Obie ścieżki prowadzą do tych samych 40/44 px dla kontrolki i wiersza, więc kolizja nie daje różnicy widocznej.

---

## 11. Siatka i punkty łamania

### 11.1. Cztery punkty łamania

| Żeton | Wartość | Nazwa | Zachowanie kokpitu |
|---|---:|---|---|
| `--dn-bp-w1` | 640 px | telefon poziomo | widok mobilny (okno **Mobile**, funkcja globalna) |
| `--dn-bp-w2` | 960 px | tablet | boczna nawigacja modułów zwija się do ikon |
| `--dn-bp-w3` | 1280 px | biurko | pełny kokpit: pasek + pas kart sesji + boczna + obszar roboczy |
| `--dn-bp-w4` | 1600 px | szerokie biurko | dwa okna komunikacji obok siebie (powłoka środowiska) |

Zamknięcie luki **L-P-11**.

### 11.2. Siatka treści

| Żeton | Wartość | Rola |
|---|---:|---|
| `--dn-tresc-max` | 1200 px | maksymalna szerokość treści dokumentowej |
| `--dn-siatka-kolumny` | 12 | liczba kolumn siatki |
| `--dn-siatka-przerwa` | 24 px (`--dn-od-6`) | przerwa między kolumnami |

Zamknięcie luki **L-P-12**.

```
 ◄── 640 ──►◄──── 960 ────►◄──── 1280 ────►◄──── 1600 ────►
   w1          w2              w3              w4
 Mobile      ikony         pełny kokpit    dwa okna komunikacji
```

**Reguła:** punkt łamania jest żetonem, ale zapytania medialne **nie mogą** go użyć bezpośrednio (`@media (min-width: var(--dn-bp-w3))` nie działa w CSS). Wartość w zapytaniu medialnym powtarza się liczbowo, a żeton pozostaje źródłem prawdy i punktem odniesienia dla warstwy skryptowej (`matchMedia`). To jedyne miejsce w systemie, gdzie liczba pojawia się poza żetonem — i jedyne, w którym jest to dopuszczone.

---

## 12. Warstwy — jedenaście poziomów `z-index`

Zamknięcie luki **L-P-10**. Skala jest rzadka (skoki 100 i większe), aby każde piętro miało zapas na warstwy pośrednie bez renumeracji.

| Żeton | Wartość | Warstwa | Co leży na tym piętrze |
|---|---:|---|---|
| `--dn-z-podloga` | 0 | podłoga | treść okna operacyjnego |
| `--dn-z-przybornik` | 10 | przybornik | przyklejone nagłówki tabel, przybornik okna |
| `--dn-z-pasek` | 100 | pasek górny | rama kokpitu |
| `--dn-z-boczna` | 200 | boczna nawigacja | lista modułów albo panel orkiestracji |
| `--dn-z-pas-komunikacji` | 300 | pas komunikacji | Chat Window |
| `--dn-z-nakladka` | 800 | nakładka | przyciemnienie pod modalem |
| `--dn-z-modal` | 900 | modal | okno modalne |
| `--dn-z-powiadomienie` | 1000 | powiadomienie | toasty |
| `--dn-z-tooltip` | 1100 | dymek | dymek kontekstowy `[?]` |
| `--dn-z-aod` | 1200 | Always On Display | warstwa centralna AOD |
| `--dn-z-centrum-polecen` | 1300 | Centrum poleceń | **zawsze najwyżej** |

```
 1300  Centrum poleceń      ██████████████████████
 1200  Always On Display    ███████████████████
 1100  dymek                ████████████████
 1000  powiadomienie        █████████████
  900  modal                ██████████
  800  nakładka             ███████
  300  pas komunikacji      █████
  200  boczna nawigacja     ████
  100  pasek górny          ███
   10  przybornik           ██
    0  podłoga              █
```

**Reguła:** komponent nigdy nie zapisuje liczby `z-index` wprost. Nowa warstwa oznacza nowy żeton i decyzję prowadzącego system — nie `z-index: 9999`.

---

## 13. Gradienty — dwa, zakres i zakazy

Zamknięcie luki **L-P-09**.

| Żeton | Wartość | Skład | Rola |
|---|---|---|---|
| `--dn-grad-atrament` | `linear-gradient(180deg, var(--dn-szary-800), var(--dn-szary-925))` | `#2A2A2A` → `#131313`, pionowo | tło awatara bez zdjęcia (`.dn-awatar`) |
| `--dn-grad-sygnal` | `linear-gradient(135deg, var(--dn-sygnal-400), var(--dn-sygnal-600))` | `#5C8CEC` → `#2457C9`, po przekątnej | medalion jednostki inteligencji (`.dn-awatar--inteligencja`), rdzeń Always On Display (`.dn-aod-rdzen`) |

### 13.1. Zakres dopuszczalny

1. awatary bez zdjęcia (człowiek → gradient atramentowy, jednostka inteligencji → gradient sygnałowy),
2. rdzeń Always On Display,
3. grafiki brandowe (materiały poza interfejsem roboczym).

### 13.2. Zakazy bezwzględne

> **Gradient NIGDY nie jest tłem przycisku, karty ani sekcji.** (`zetony.css` rozdz. 9, `zetony.json`, KANON rozdz. 2.)

Dodatkowo z katalogu anty-domyślnych: **fioletowy gradient „AI"** jest pierwszym zablokowanym odruchem systemu. Oba dopuszczone gradienty są wyprowadzone z zadeklarowanych prymitywów — nie wolno wprowadzić trzeciego bez decyzji prowadzącego system.

### 13.3. Konsekwencja dla kontrastu

Gradient jest **stały w obu motywach**, dlatego tekst na nim jest jedynym miejscem, w którym komponent legalnie sięga po prymityw (`--dn-szary-0` / `--dn-szary-100`, wyjątki 4–6 z rozdz. 1.4). Najsłabszy punkt gradientu atramentowego to `#2A2A2A`; `--dn-szary-100` na nim daje 12,15 — z zapasem ponad progiem AAA. Najsłabszy punkt gradientu sygnałowego to `#5C8CEC`; `--dn-szary-0` na nim daje 3,27 — dlatego na tym gradiencie umieszcza się wyłącznie inicjał albo znak graficzny w dużym stopniu, nigdy tekst czytany.

---

## 14. Motyw jasny — pełna tabela żetonów semantycznych

Blok `:root[data-theme='light']`, `color-scheme: light`. **43 żetony.** Motyw jest definiowany osobno, nie wywodzony z ciemnego.

### 14.1. Powierzchnie i nakładki (7)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-tlo` | `#F4F4F4` | `--dn-szary-50` | tło dokumentu — papier roboczy | tekst na nim 16,14 |
| `--dn-powierzchnia` | `#FFFFFF` | `--dn-szary-0` | karty, okna, pola | tekst na niej 17,76 |
| `--dn-powierzchnia-2` | `#ECECEC` | `--dn-szary-100` | wgłębienia, tory postępu, bloki kodu | tekst-2 na niej 5,24 |
| `--dn-panel` | `#FAFAFA` | `--dn-szary-25` | panele boczne, nagłówki tabel | tekst na nim 17,01 |
| `--dn-hover` | `rgba(10,10,10,0.05)` | `--dn-szary-1000` @ 5% | nakładka najechania | na powierzchni daje `#F3F3F3` |
| `--dn-wcisniecie` | `rgba(10,10,10,0.09)` | `--dn-szary-1000` @ 9% | nakładka naciśnięcia | na powierzchni daje `#E9E9E9` |
| `--dn-nakladka` | `rgba(10,10,10,0.45)` | `--dn-szary-1000` @ 45% | przyciemnienie pod modalem | rozmycie 2 px wyłącznie tutaj |

### 14.2. Obrysy (3)

| Żeton | Wartość | Mapuje się na | Rola | Kontrast do powierzchni |
|---|---|---|---|---:|
| `--dn-obrys` | `#E3E3E3` | `--dn-szary-150` | podstawowa linia rozdzielająca kart i paneli | 1,28 |
| `--dn-obrys-subtelny` | `#ECECEC` | `--dn-szary-100` | linia wewnątrz komponentu, separator | 1,18 |
| `--dn-obrys-mocny` | `#C0C0C0` | `--dn-szary-300` | obrys kontrolki, kciuk paska przewijania | 1,82 |

### 14.3. Tekst (4)

| Żeton | Wartość | Mapuje się na | Rola | Na `--dn-tlo` | Na `--dn-powierzchnia` | Werdykt |
|---|---|---|---|---:|---:|---|
| `--dn-tekst` | `#181818` | `--dn-szary-900` | tekst główny | 16,14 | 17,76 | AAA |
| `--dn-tekst-2` | `#616161` | `--dn-szary-600` | tekst drugorzędny, opisy | 5,63 | 6,19 | AA |
| `--dn-tekst-3` | `#7C7C7C` | `--dn-szary-500` | **wyłącznie metadane i tekst ≥ 18,66 px półgruby** | 3,80 | 4,17 | AA dla dużego tekstu |
| `--dn-tekst-inv` | `#FFFFFF` | `--dn-szary-0` | tekst na powierzchni odwróconej | 17,76 na atramencie | — | AAA |

### 14.4. Działanie główne — inwersja atramentu (3)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-atrament` | `#181818` | `--dn-szary-900` | tło przycisku głównego | biel na nim 17,76 |
| `--dn-atrament-hover` | `#0A0A0A` | `--dn-szary-1000` | przycisk główny przy najechaniu | biel na nim 19,80 |
| `--dn-atrament-tekst` | `#FFFFFF` | `--dn-szary-0` | tekst na przycisku głównym | 17,76 |

### 14.5. Sygnał (7)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-sygnal` | `#2457C9` | `--dn-sygnal-600` | tekst sygnałowy, odnośniki | 5,82 na tle · 6,40 na powierzchni |
| `--dn-sygnal-mocny` | `#1D49AF` | `--dn-sygnal-700` | odnośnik przy najechaniu | 7,28 na tle |
| `--dn-sygnal-wypelnienie` | `#2457C9` | `--dn-sygnal-600` | wypełnienie przycisku sygnałowego, tor przełącznika | biel na nim 6,40 |
| `--dn-sygnal-wypelnienie-hover` | `#1D49AF` | `--dn-sygnal-700` | jw. przy najechaniu | biel na nim 8,00 |
| `--dn-sygnal-tlo` | `#EDF3FE` | `--dn-sygnal-100` | tło stanu wybranego, zaznaczenie tekstu | sygnał na nim 5,74 |
| `--dn-sygnal-obrys` | `#C9DAFB` | `--dn-sygnal-200` | obrys stanu wybranego | — |
| `--dn-kropka` | `#3B6FE0` | `--dn-sygnal-500` | **kropka sygnału** — element sygnaturowy | 4,21 na tle |

### 14.6. Fokus (2)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-fokus` | `#3B6FE0` | `--dn-sygnal-500` | pierścień fokusu 2 px + odsunięcie 2 px | 4,21 na tle (próg 3,0 — element nietekstowy) |
| `--dn-fokus-cien` | `rgba(59,111,224,0.30)` | `--dn-sygnal-500` @ 30% | poświata pola aktywnego | — |

### 14.7. Stany (12)

| Stan | `-tekst` | `-tlo` | `-obrys` | Kontrast tekst/tło stanu | Ikona towarzysząca |
|---|---|---|---|---:|---|
| sukces | `#1E7A46` (`zielen-700`) | `#E7F5EE` (`zielen-100`) | `#BCE3CD` (`zielen-200`) | 4,76 | `ptaszek-kolo` |
| ostrzeżenie | `#8A5B0C` (`bursztyn-700`) | `#FBF2DE` (`bursztyn-100`) | `#EED9A7` (`bursztyn-200`) | 5,26 | `ostrzezenie` |
| błąd | `#C0362F` (`czerwien-700`) | `#FCEDEB` (`czerwien-100`) | `#F3C8C4` (`czerwien-200`) | 4,84 | `blad` |
| informacja | `#2457C9` (`sygnal-600`) | `#EDF3FE` (`sygnal-100`) | `#C9DAFB` (`sygnal-200`) | 5,74 | `info` |

### 14.8. Cienie (5)

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-cien-1` | `0 1px 2px rgba(10,10,10,0.05)` | karta w spoczynku, plakietka uniesiona |
| `--dn-cien-2` | `0 2px 8px rgba(10,10,10,0.07), 0 1px 2px rgba(10,10,10,0.05)` | karta przy najechaniu, dymek |
| `--dn-cien-3` | `0 8px 24px rgba(10,10,10,0.12)` | menu kontekstowe, toast |
| `--dn-cien-lg` | `0 20px 48px rgba(10,10,10,0.18)` | modal — największy cień systemu |
| `--dn-cien-sygnal` | `0 0 0 3px rgba(59,111,224,0.18)` | **ZAREZERWOWANY**: pierścień pola aktywnego i poświata kropki sygnału |

Jeden kierunek światła — z góry. Wszystkie cienie tonowane neutralnie z `--dn-szary-1000`.

---

## 15. Motyw ciemny — pełna tabela żetonów semantycznych

Blok `:root[data-theme='dark']`, `color-scheme: dark`. **43 żetony**, te same nazwy, inne wartości. Motyw definiowany osobno — nie jest inwersją motywu jasnego.

### 15.1. Powierzchnie i nakładki (7)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-tlo` | `#0F0F0F` | `--dn-szary-950` | tło dokumentu | tekst na nim 16,23 |
| `--dn-powierzchnia` | `#181818` | `--dn-szary-900` | karty, okna, pola | tekst na niej 15,03 |
| `--dn-powierzchnia-2` | `#131313` | `--dn-szary-925` | wgłębienia, tory, bloki kodu | wartość równa `--dn-rama` |
| `--dn-panel` | `#212121` | `--dn-szary-850` | panele uniesione, nagłówki tabel | tekst na nim 13,63 |
| `--dn-hover` | `rgba(255,255,255,0.06)` | biel @ 6% | nakładka najechania | na powierzchni daje `#262626` |
| `--dn-wcisniecie` | `rgba(255,255,255,0.10)` | biel @ 10% | nakładka naciśnięcia | na powierzchni daje `#2F2F2F` |
| `--dn-nakladka` | `rgba(0,0,0,0.62)` | czerń @ 62% | przyciemnienie pod modalem | czerń wyłącznie jako krycie, nigdy jako barwa powierzchni |

**Uwaga porządkowa:** w motywie ciemnym powierzchnia jest **jaśniejsza** od tła (`#181818` > `#0F0F0F`), a panel jaśniejszy od powierzchni (`#212121`). Kierunek uniesienia jest odwrotny niż w motywie jasnym (tam `#FFFFFF` > `#F4F4F4`, ale panel `#FAFAFA` leży **poniżej** powierzchni). Ta asymetria jest celowa: w motywie jasnym panel boczny cofa się, w ciemnym — wysuwa.

### 15.2. Obrysy (3)

| Żeton | Wartość | Mapuje się na | Rola | Kontrast do powierzchni |
|---|---|---|---|---:|
| `--dn-obrys` | `#2A2A2A` | `--dn-szary-800` | podstawowa linia rozdzielająca | 1,24 |
| `--dn-obrys-subtelny` | `#212121` | `--dn-szary-850` | linia wewnątrz komponentu | 1,10 |
| `--dn-obrys-mocny` | `#3A3A3A` | `--dn-szary-750` | obrys kontrolki, kciuk paska przewijania | 1,56 |

### 15.3. Tekst (4)

| Żeton | Wartość | Mapuje się na | Rola | Na `--dn-tlo` | Na `--dn-powierzchnia` | Werdykt |
|---|---|---|---|---:|---:|---|
| `--dn-tekst` | `#ECECEC` | `--dn-szary-100` | tekst główny | 16,23 | 15,03 | AAA |
| `--dn-tekst-2` | `#9E9E9E` | `--dn-szary-400` | tekst drugorzędny | 7,15 | 6,63 | AAA / AA |
| `--dn-tekst-3` | `#7C7C7C` | `--dn-szary-500` | **wyłącznie metadane** | 4,59 | 4,25 | AA dla dużego tekstu |
| `--dn-tekst-inv` | `#0F0F0F` | `--dn-szary-950` | tekst na powierzchni odwróconej | 16,23 na atramencie | — | AAA |

### 15.4. Działanie główne — inwersja atramentu (3)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-atrament` | `#ECECEC` | `--dn-szary-100` | tło przycisku głównego — **biel** | atrament-tekst na nim 16,23 |
| `--dn-atrament-hover` | `#FAFAFA` | `--dn-szary-25` | przycisk główny przy najechaniu | atrament-tekst na nim 18,36 |
| `--dn-atrament-tekst` | `#0F0F0F` | `--dn-szary-950` | tekst na przycisku głównym | 16,23 |

To najczystszy przykład inwersji: **ten sam żeton `--dn-atrament` jest czernią w motywie jasnym i bielą w ciemnym.** Komponent `.dn-btn--atrament` nie wie o motywie i nie musi wiedzieć.

### 15.5. Sygnał (7)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-sygnal` | `#8FB2F5` | `--dn-sygnal-300` | tekst sygnałowy, odnośniki | 9,00 na tle · 8,33 na powierzchni |
| `--dn-sygnal-mocny` | `#C9DAFB` | `--dn-sygnal-200` | odnośnik przy najechaniu | 13,61 na tle |
| `--dn-sygnal-wypelnienie` | `#3B6FE0` | `--dn-sygnal-500` | wypełnienie przycisku sygnałowego, tor przełącznika | biel na nim 4,63 |
| `--dn-sygnal-wypelnienie-hover` | `#5C8CEC` | `--dn-sygnal-400` | jw. przy najechaniu | biel na nim 3,27 |
| `--dn-sygnal-tlo` | `rgba(59,111,224,0.16)` | `--dn-sygnal-500` @ 16% | tło stanu wybranego | na powierzchni daje `#1E2638`; sygnał na nim 7,09 |
| `--dn-sygnal-obrys` | `rgba(92,140,236,0.40)` | `--dn-sygnal-400` @ 40% | obrys stanu wybranego | na powierzchni daje `#33466D` |
| `--dn-kropka` | `#5C8CEC` | `--dn-sygnal-400` | **kropka sygnału** | 5,86 na tle |

**Zasada tła stanu w motywie ciemnym:** tło nie jest osobnym prymitywem, lecz **barwą tekstu stanu z niskim kryciem**. Dzięki temu tło automatycznie dopasowuje się do powierzchni, na której leży (powierzchnia `#181818` albo panel `#212121`) — bez mnożenia żetonów.

### 15.6. Fokus (2)

| Żeton | Wartość | Mapuje się na | Rola | Pomiar |
|---|---|---|---|---|
| `--dn-fokus` | `#5C8CEC` | `--dn-sygnal-400` | pierścień fokusu 2 px + odsunięcie 2 px | 5,86 na tle (próg 3,0) |
| `--dn-fokus-cien` | `rgba(92,140,236,0.35)` | `--dn-sygnal-400` @ 35% | poświata pola aktywnego | — |

Fokus w motywie ciemnym jest **o stopień jaśniejszy** niż w jasnym (`sygnal-400` zamiast `sygnal-500`) — inaczej pierścień zlałby się z ciemnym tłem.

### 15.7. Stany (12) — zamknięcie luki L-P-19

| Stan | `-tekst` | `-tlo` (alfa) | `-obrys` (alfa) | Tło po złożeniu na `#181818` | Kontrast tekst/tło | Ikona |
|---|---|---|---|---|---:|---|
| sukces | `#4FBD85` (`zielen-400`) | `rgba(79,189,133,0.14)` | `rgba(79,189,133,0.34)` | `#202F27` | 5,98 | `ptaszek-kolo` |
| ostrzeżenie | `#E0A63C` (`bursztyn-400`) | `rgba(224,166,60,0.14)` | `rgba(224,166,60,0.34)` | `#342C1D` | 6,36 | `ostrzezenie` |
| błąd | `#EE7168` (`czerwien-400`) | `rgba(238,113,104,0.14)` | `rgba(238,113,104,0.34)` | `#362423` | 5,02 | `blad` |
| informacja | `#8FB2F5` (`sygnal-300`) | `rgba(59,111,224,0.16)` | `rgba(92,140,236,0.40)` | `#1E2638` | 7,09 | `info` |

Komentarz źródłowy w `zetony.css`: *„stany — pełny komplet dla motywu ciemnego (luka L-P-19 zamknięta)"*.

### 15.8. Cienie (5)

| Żeton | Wartość | Zastosowanie |
|---|---|---|
| `--dn-cien-1` | `0 1px 2px rgba(0,0,0,0.35)` | karta w spoczynku |
| `--dn-cien-2` | `0 3px 10px rgba(0,0,0,0.42), 0 1px 2px rgba(0,0,0,0.30)` | karta przy najechaniu, dymek |
| `--dn-cien-3` | `0 10px 30px rgba(0,0,0,0.52)` | menu kontekstowe, toast |
| `--dn-cien-lg` | `0 24px 60px rgba(0,0,0,0.65)` | modal |
| `--dn-cien-sygnal` | `0 0 0 3px rgba(92,140,236,0.22)` | **ZAREZERWOWANY**: pierścień pola aktywnego, poświata kropki |

Cienie motywu ciemnego są **głębsze i większe** (rozmycie 24–60 px wobec 20–48 px), bo cień na ciemnym tle wymaga większego zasięgu, aby był w ogóle czytelny.

---

## 16. Tabela kontrastów — komplet 33 pomiarów

Źródło: `zasoby/zetony/kontrasty.json`. Metoda: **luminancja względna WCAG 2.1** — `(L_jaśniejszy + 0,05) / (L_ciemniejszy + 0,05)`, gdzie `L = 0,2126·R + 0,7152·G + 0,0722·B` po linearyzacji sRGB. Paleta odniesienia w pliku: 38 barw.

**Legenda werdyktu:**

| Werdykt | Warunek | Znaczenie |
|---|---|---|
| **AAA** | kontrast ≥ 7,0 przy progu 4,5 | WCAG 2.1 kryterium 1.4.6 — tekst dowolnego rozmiaru |
| **AA** | kontrast ≥ próg | WCAG 2.1 kryterium 1.4.3 (tekst) albo 1.4.11 (element nietekstowy) |
| **próg systemowy** | próg 1,6 | rozdzielenie powierzchni; **nie jest progiem WCAG** |
| **nie przechodzi** | kontrast < próg | w komplecie nie występuje ani razu |

| # | Para | Pierwszy plan | Tło | Kontrast | Próg | Zapas | Werdykt |
|---:|---|---|---|---:|---:|---:|---|
| 1 | tekst / tło (jasny) | `szary-900` `#181818` | `szary-50` `#F4F4F4` | 16,14 | 4,5 | +11,64 | **AAA** |
| 2 | tekst / powierzchnia (jasny) | `szary-900` `#181818` | `szary-0` `#FFFFFF` | 17,76 | 4,5 | +13,26 | **AAA** |
| 3 | tekst-2 / tło (jasny) | `szary-600` `#616161` | `szary-50` `#F4F4F4` | 5,63 | 4,5 | +1,13 | **AA** |
| 4 | tekst-2 / powierzchnia (jasny) | `szary-600` `#616161` | `szary-0` `#FFFFFF` | 6,19 | 4,5 | +1,69 | **AA** |
| 5 | tekst-2 / powierzchnia-2 (jasny) | `szary-600` `#616161` | `szary-100` `#ECECEC` | 5,24 | 4,5 | +0,74 | **AA** |
| 6 | tekst-3 / tło (jasny) | `szary-500` `#7C7C7C` | `szary-50` `#F4F4F4` | 3,80 | 3,0 | +0,80 | **AA** (duży tekst / metadane) |
| 7 | tekst-3 / powierzchnia (jasny) | `szary-500` `#7C7C7C` | `szary-0` `#FFFFFF` | 4,17 | 3,0 | +1,17 | **AA** (duży tekst / metadane) |
| 8 | biel / przycisk atrament (jasny) | `szary-0` `#FFFFFF` | `szary-900` `#181818` | 17,76 | 4,5 | +13,26 | **AAA** |
| 9 | sygnał-tekst / tło (jasny) | `sygnal-600` `#2457C9` | `szary-50` `#F4F4F4` | 5,82 | 4,5 | +1,32 | **AA** |
| 10 | sygnał-tekst / powierzchnia (jasny) | `sygnal-600` `#2457C9` | `szary-0` `#FFFFFF` | 6,40 | 4,5 | +1,90 | **AA** |
| 11 | sygnał-tekst / sygnał-tło (jasny) | `sygnal-600` `#2457C9` | `sygnal-100` `#EDF3FE` | 5,74 | 4,5 | +1,24 | **AA** |
| 12 | biel / sygnał wypełnienie | `szary-0` `#FFFFFF` | `sygnal-600` `#2457C9` | 6,40 | 4,5 | +1,90 | **AA** |
| 13 | obrys kontrolki / tło (jasny) | `szary-300` `#C0C0C0` | `szary-50` `#F4F4F4` | 1,65 | 1,6 | +0,05 | **próg systemowy** |
| 14 | fokus sygnał / tło (jasny) | `sygnal-500` `#3B6FE0` | `szary-50` `#F4F4F4` | 4,21 | 3,0 | +1,21 | **AA** (element nietekstowy) |
| 15 | sukces tekst / tło stanu (jasny) | `sukces-700` `#1E7A46` | `sukces-100` `#E7F5EE` | 4,76 | 4,5 | +0,26 | **AA** |
| 16 | ostrzeżenie tekst / tło stanu (jasny) | `ostrz-700` `#8A5B0C` | `ostrz-100` `#FBF2DE` | 5,26 | 4,5 | +0,76 | **AA** |
| 17 | błąd tekst / tło stanu (jasny) | `blad-700` `#C0362F` | `blad-100` `#FCEDEB` | 4,84 | 4,5 | +0,34 | **AA** |
| 18 | sukces tekst / powierzchnia (jasny) | `sukces-700` `#1E7A46` | `szary-0` `#FFFFFF` | 5,35 | 4,5 | +0,85 | **AA** |
| 19 | ostrzeżenie tekst / powierzchnia (jasny) | `ostrz-700` `#8A5B0C` | `szary-0` `#FFFFFF` | 5,86 | 4,5 | +1,36 | **AA** |
| 20 | błąd tekst / powierzchnia (jasny) | `blad-700` `#C0362F` | `szary-0` `#FFFFFF` | 5,51 | 4,5 | +1,01 | **AA** |
| 21 | tekst / tło (ciemny) | `szary-100` `#ECECEC` | `szary-950` `#0F0F0F` | 16,23 | 4,5 | +11,73 | **AAA** |
| 22 | tekst / powierzchnia (ciemny) | `szary-100` `#ECECEC` | `szary-900` `#181818` | 15,03 | 4,5 | +10,53 | **AAA** |
| 23 | tekst-2 / tło (ciemny) | `szary-400` `#9E9E9E` | `szary-950` `#0F0F0F` | 7,15 | 4,5 | +2,65 | **AAA** |
| 24 | tekst-2 / powierzchnia (ciemny) | `szary-400` `#9E9E9E` | `szary-900` `#181818` | 6,63 | 4,5 | +2,13 | **AA** |
| 25 | tekst-3 / powierzchnia (ciemny) | `szary-500` `#7C7C7C` | `szary-900` `#181818` | 4,25 | 3,0 | +1,25 | **AA** (duży tekst / metadane) |
| 26 | atrament / przycisk biel (ciemny) | `szary-950` `#0F0F0F` | `szary-100` `#ECECEC` | 16,23 | 4,5 | +11,73 | **AAA** |
| 27 | sygnał-tekst / tło (ciemny) | `sygnal-300` `#8FB2F5` | `szary-950` `#0F0F0F` | 9,00 | 4,5 | +4,50 | **AAA** |
| 28 | sygnał-tekst / powierzchnia (ciemny) | `sygnal-300` `#8FB2F5` | `szary-900` `#181818` | 8,33 | 4,5 | +3,83 | **AAA** |
| 29 | fokus sygnał / tło (ciemny) | `sygnal-400` `#5C8CEC` | `szary-950` `#0F0F0F` | 5,86 | 3,0 | +2,86 | **AA** (element nietekstowy) |
| 30 | sukces tekst / powierzchnia (ciemny) | `sukces-400` `#4FBD85` | `szary-900` `#181818` | 7,57 | 4,5 | +3,07 | **AAA** |
| 31 | ostrzeżenie tekst / powierzchnia (ciemny) | `ostrz-400` `#E0A63C` | `szary-900` `#181818` | 8,19 | 4,5 | +3,69 | **AAA** |
| 32 | błąd tekst / powierzchnia (ciemny) | `blad-400` `#EE7168` | `szary-900` `#181818` | 6,09 | 4,5 | +1,59 | **AA** |
| 33 | biel / sygnał wypełnienie (ciemny przycisk) | `szary-0` `#FFFFFF` | `sygnal-500` `#3B6FE0` | 4,63 | 3,0 | +1,63 | **AA** |

### 16.1. Rozkład wyniku

| Wskaźnik | Wartość |
|---|---:|
| Pomiarów ogółem | 33 |
| Pomiarów z werdyktem **AAA** | 12 |
| Pomiarów z werdyktem **AA** | 20 |
| Pomiarów o progu systemowym (1,6) | 1 |
| Pomiarów, które **nie przechodzą** | **0** |
| Pomiary o progu 4,5 (tekst) | 26 |
| Pomiary o progu 3,0 (duży tekst / element nietekstowy) | 6 |
| Najmniejszy zapas ponad progiem (poza progiem systemowym) | +0,26 (pomiar 15: sukces na tle stanu) |
| Największy zapas | +13,26 (pomiary 2 i 8) |

### 16.2. Pary najciaśniejsze — do obserwacji przy zmianach

| Pomiar | Para | Zapas | Ryzyko |
|---|---|---:|---|
| 15 | sukces tekst / tło stanu (jasny) | +0,26 | każde rozjaśnienie `--dn-zielen-700` albo przyciemnienie `--dn-zielen-100` zejdzie poniżej progu |
| 17 | błąd tekst / tło stanu (jasny) | +0,34 | jw. dla rodziny czerwieni |
| 13 | obrys kontrolki / tło (jasny) | +0,05 | próg systemowy — patrz decyzja D8 |
| 5 | tekst-2 / powierzchnia-2 (jasny) | +0,74 | `--dn-tekst-2` na wgłębieniu; nie wolno zejść poniżej `szary-600` |
| 16 | ostrzeżenie tekst / tło stanu (jasny) | +0,76 | jw. dla rodziny bursztynu |

---

## 17. Żetony a JSON — jak czytać `zetony.json`

### 17.1. Rola pliku

`zetony.json` jest **maszynowym odpowiednikiem** `zetony.css`, nie drugim źródłem prawdy. Rozstrzyga plik CSS (jest ładowany przez aplikację); JSON służy generatorom, walidatorom, eksportowi do innych warstw i narzędziom projektowym.

### 17.2. Struktura pliku

```
zetony.json
├── _meta                      metryka: produkt, wersja, status, data, kierunek,
│                              źródło, 6 zasad nienegocjowalnych
├── prymitywy
│   ├── szarosc                18 par nazwa→hex
│   ├── sygnal                 8 par
│   └── stany
│       ├── zielen             4 pary
│       ├── bursztyn           4 pary
│       └── czerwien           4 pary
├── rama-kokpitu               5 par + _uwaga
├── semantyczne
│   ├── jasny                  26 par (bez stanów i cieni)
│   └── ciemny                 26 par
├── stany-semantyczne
│   ├── jasny                  4 stany × {tekst, tlo, obrys}
│   ├── ciemny                 4 stany × {tekst, tlo, obrys}
│   └── _uwaga                 „stan nigdy samym kolorem"
├── typografia                 kroje (+ _role, _licencja) · stopnie ·
│                              grubosci · interlinie · tracking
├── przestrzen                 odstepy (11) · promienie (6)
├── ruch                       ease · 4 czasy · zapis o prefers-reduced-motion
├── wymiary                    _gestosc · 27 par ·
│                              „dotyk (pointer coarse)" · „gestosc-przestronna"
├── siatka-i-lamanie           4 punkty · tresc-max · kolumny · przerwa · _uwaga
├── warstwy                    11 poziomów (liczby, nie napisy) · _uwaga
├── gradienty                  2 pary + _zakres
└── cienie                     jasny (5) · ciemny (5) · _uwaga
```

### 17.3. Reguły czytania

| Reguła | Treść |
|---|---|
| **J1** | Klucz zaczynający się od podkreślnika (`_meta`, `_uwaga`, `_zakres`, `_role`, `_licencja`, `_gestosc`) jest **komentarzem**, nie żetonem. Generator ma go pominąć albo przenieść do dokumentacji. |
| **J2** | Wartości semantyczne w JSON są **rozwiązane do wartości końcowej** (`"--dn-tlo": "#F4F4F4"`), podczas gdy w CSS zapisane są jako odwołanie (`var(--dn-szary-50)`). JSON traci informację o mapowaniu — mapowanie odczytuje się z CSS albo z rozdz. 14–15 niniejszego dokumentu. |
| **J3** | Warstwy `z-index` są **liczbami**, nie napisami (`"--dn-z-modal": 900`). Pozostałe wartości są napisami z jednostką. |
| **J4** | Klucze `"dotyk (pointer coarse)"` i `"gestosc-przestronna (data-gestosc)"` w gałęzi `wymiary` opisują **nadpisania kontekstowe**, nie osobne żetony. |
| **J5** | Gałąź `semantyczne` **nie zawiera** stanów ani cieni — te mają własne gałęzie (`stany-semantyczne`, `cienie`). Suma 26 + 12 + 5 = 43 żetony semantyczne na motyw. |
| **J6** | Gradienty w JSON zapisane są z **rozwiniętymi wartościami hex** (`linear-gradient(180deg, #2A2A2A, #131313)`), a w CSS przez `var()`. Przy zmianie prymitywu trzeba zaktualizować oba pliki. |

### 17.4. Sześć zasad z `_meta.zasady` — zapis maszynowy

1. jednostka bazowa 4 px,
2. oba motywy równoprawne — definiowane osobno,
3. stan nigdy samym kolorem,
4. zero blokad (ADL-017),
5. WCAG 2.1 AA na wejściu,
6. zakaz `#000000`; `#FFFFFF` wyłącznie jako powierzchnia kart i tekst na atramencie.

### 17.5. Kolejność aktualizacji przy zmianie wartości

```
1. zetony.css      ← źródło prawdy dla aplikacji; zmiana najpierw tutaj
2. zetony.json     ← odzwierciedlenie wartości (uwaga na J2 i J6)
3. kontrasty.json  ← ponowny pomiar każdej pary dotkniętej zmianą
4. dokumentacja    ← 04-tokens.md + 04-tokens.html (rozdz. 3–6, 14–16)
5. komponenty.css  ← wyłącznie jeśli zmiana wprowadza/usuwa nazwę żetonu
```

Pominięcie kroku 3 jest **najczęstszym możliwym błędem**: wartość zmieniona bez ponownego pomiaru unieważnia deklarację zgodności WCAG całego pakietu.

---

## 18. Jak dodać żeton — procedura

### 18.1. Bramka wstępna — trzy pytania

| Pytanie | Odpowiedź „nie" oznacza |
|---|---|
| Czy wartość wystąpi w **co najmniej dwóch** niezależnych miejscach systemu? | to nie jest żeton — to wartość komponentu (reguła R8) |
| Czy istniejący żeton **nie** obsługuje tej roli po nadaniu mu właściwej nazwy? | należy poprawić nazwę istniejącego, nie dodawać nowego (reguła R9) |
| Czy rola jest **semantyczna** (co to znaczy), a nie opisowa (jak to wygląda)? | nazwa `--dn-niebieski-jasny` jest błędna; `--dn-informacja-tlo` jest poprawna |

Trzy odpowiedzi „tak" otwierają procedurę.

### 18.2. Siedem kroków

| Krok | Czynność | Plik / miejsce | Kryterium zakończenia |
|---:|---|---|---|
| 1 | **Ustal warstwę.** Prymityw czy semantyczny? Jeżeli semantyczny — czy potrzebuje nowego prymitywu, czy wystarczy istniejący stopień skali? | decyzja projektowa | warstwa zapisana; przy nowej rodzinie barwnej wymagana pełna czwórka stopni (reguła R7) |
| 2 | **Nadaj nazwę** wg konwencji `--dn-<grupa>-<wariant>` i reguł R1–R10 (rozdz. 2.3). | — | nazwa nie koliduje z żadną z istniejących 177 |
| 3 | **Zadeklaruj w `zetony.css`** w bloku właściwej warstwy: `:root` (niezależny) albo **obu** blokach motywów (semantyczny). Żeton semantyczny zadeklarowany tylko w jednym motywie jest błędem krytycznym. | `zasoby/zetony/zetony.css` | wartość zapisana jako `var(--dn-<prymityw>)` albo `rgba()` wyprowadzona z prymitywu (reguła R10) |
| 4 | **Powiel do bloków `@media (prefers-color-scheme: …)`** — dotyczy wyłącznie żetonów semantycznych. Wartości muszą być identyczne z blokami motywów. | `zasoby/zetony/zetony.css`, rozdz. 11 pliku | oba bloki medialne uzupełnione |
| 5 | **Zmierz kontrast** każdej nowej pary tekst/tło oraz element/tło i dopisz pomiar do `kontrasty.json` (`para`, `fg`, `bg`, `kontrast`, `prog`, `ok`). Przy alfie mierzy się barwę **po złożeniu** na powierzchni docelowej. | `zasoby/zetony/kontrasty.json` | `ok: true` przy progu 4,5 (tekst) albo 3,0 (element nietekstowy / duży tekst) |
| 6 | **Odzwierciedl w `zetony.json`** w odpowiedniej gałęzi, z zachowaniem reguł J2, J3 i J6. | `zasoby/zetony/zetony.json` | liczba żetonów w JSON zgadza się z liczbą w CSS |
| 7 | **Udokumentuj** — dopisz wiersz do właściwej tabeli w `01-dokumentacja-md/04-tokens.md` i kartę w `02-dokumentacja-html/04-tokens.html`; jeżeli żeton wchodzi do komponentu, uzupełnij też dokumentację komponentu. | dokumentacja | rachunek żetonów (rozdz. 1.3) zaktualizowany |

### 18.3. Lista kontrolna odbioru nowego żetonu

- [ ] Nazwa zgodna z konwencją `--dn-<grupa>-<wariant>` i regułami R1–R10.
- [ ] Warstwa ustalona; komponent nie sięga po prymityw (poza sześcioma wyjątkami z rozdz. 1.4).
- [ ] Żeton semantyczny obecny w **obu** blokach motywów i w **obu** blokach `prefers-color-scheme`.
- [ ] Wartość zapisana przez `var()` albo `rgba()` z prymitywu — zero nowych hexów.
- [ ] Brak `#000000` w wartości i w łańcuchu odwołań.
- [ ] Nowa rodzina stanu ma pełną trójkę `-tekst` / `-tlo` / `-obrys` **w obu motywach**.
- [ ] Pomiar kontrastu dopisany do `kontrasty.json`; `ok: true`.
- [ ] Jeżeli żeton dotyczy stanu — wskazana ikona towarzysząca (stan nigdy samym kolorem).
- [ ] Jeżeli żeton dotyczy ruchu — mieści się w przedziale 100–220 ms albo jest tętnem 2,4 s.
- [ ] Jeżeli żeton dotyczy wymiaru — wartość jest krotnością 4 px.
- [ ] `zetony.json` zaktualizowany; rachunek żetonów zgadza się z rozdz. 1.3.
- [ ] Dokumentacja MD i HTML uzupełniona; oba motywy sprawdzone wizualnie.

### 18.4. Czego nie wolno zrobić

| Zakaz | Dlaczego |
|---|---|
| dodać żeton wyłącznie do motywu ciemnego „bo w jasnym się nie przyda" | komponent nie wie o motywie; brak wartości = pusta deklaracja `var()` i awaria wizualna |
| użyć nowego hexa zamiast stopnia skali | skale są zamknięte; nowa wartość rozmywa monochromatyczną precyzję |
| wprowadzić drugą barwę akcentu | KIERUNEK 3.1: jeden sygnał; rodzina informacyjna jest scalona z sygnałem celowo |
| zapisać `z-index` wprost w komponencie | warstwy są zamkniętym rejestrem 11 poziomów |
| dodać gradient do przycisku, karty albo sekcji | zakaz bezwzględny (rozdz. 13.2) |
| obsłużyć `prefers-reduced-motion` w komponencie | obsługa jest globalna, w żetonach |
| zastosować `disabled` jako sposób sygnalizacji stanu | ADL-017 — zero blokad; komunikat po naciśnięciu albo opis obok |

---

## 19. Decyzje projektowe

Sekcja obejmuje rozstrzygnięcia podjęte **ponad** to, co dosłownie rozstrzygają źródła, wraz z uzasadnieniem i podstawą.

| # | Kwestia | Rozstrzygnięcie | Podstawa |
|---:|---|---|---|
| **D1** | Rachunek żetonów systemu nie był nigdzie podany liczbowo | Ustalono i zapisano: **177 unikalnych nazw** = 134 niezależne od motywu + 43 semantyczne. Liczby wyprowadzono z `zetony.css` przez zliczenie deklaracji w blokach `:root`, `:root[data-theme='light']` i `:root[data-theme='dark']`. | pomiar na pliku źródłowym; zgodność sumy grup z sumą całkowitą potwierdzona |
| **D2** | Prymitywy `--dn-szary-200` i `--dn-szary-700` nie mają mapowania semantycznego | Nazwano je **stopniami rezerwowymi** i opisano funkcję (ciągłość skali, zapas dla przyszłych ról). Nie usuwa się ich ani nie oznacza jako błędu. | zliczenie użyć w `zetony.css`, `fundament.css`, `komponenty.css` — zero wystąpień poza deklaracją |
| **D3** | To samo dotyczy `--dn-sygnal-800` | Oznaczono jako stopień rezerwowy z przewidzianą rolą (wciśnięcie wypełnienia sygnałowego). Rola jest **propozycją**, nie stanem obowiązującym — do decyzji prowadzącego system. | analiza rodziny sygnału: `-500` bazowy, `-600` wypełnienie, `-700` wypełnienie-hover; `-800` bez roli |
| **D4** | Reguła „komponent nie sięga po prymityw" ma w bibliotece odstępstwa | Przeprowadzono inwentaryzację i ustalono **zamkniętą listę sześciu wyjątków** ze wspólnym uzasadnieniem: element leży na powierzchni nieprzełączającej się z motywem. Lista jest zamknięta — siódmy wyjątek wymaga decyzji. | `komponenty.css`, linie 90, 163, 307, 1073, 1083, 1196 |
| **D5** | `komponenty.css` zawiera jedną wartość szesnastkową wprost: `var(--dn-sygnal-300, #8FB2F5)` | Zakwalifikowano jako **wartość zapasowa `var()`**, nie naruszenie zakazu — hex jest identyczny z wartością żetonu i nie tworzy drugiego źródła prawdy. Odnotowano do rozważenia usunięcia przy najbliższym przeglądzie biblioteki. | `komponenty.css`, linia 163; KANON rozdz. 2 |
| **D6** | Komentarze w `zetony.css` podają wartości kontrastu zaokrąglone i miejscami rozbieżne z `kontrasty.json` (np. `--dn-tekst-2` opisany jako „6,4:1 na tle", pomiar 3 daje 5,63 na tle i 6,19 na powierzchni; `--dn-sygnal` opisany jako „5,5:1", pomiary 9–10 dają 5,82 i 6,40; `--dn-tekst` w motywie ciemnym opisany „15,0:1 na tle", pomiar 21 daje 16,23 na tle i 15,03 na powierzchni) | **Wiążący jest `kontrasty.json`** — zgodnie z KANON rozdz. 9 („komplet w `zasoby/zetony/kontrasty.json`"). Komentarze CSS potraktowano jako skrót redakcyjny, w którym pomylono odniesienie „tło" z „powierzchnią". Rozbieżność zgłoszono do korekty redakcyjnej pliku; **żadna wartość barwy nie wymaga zmiany**. | porównanie komentarzy `zetony.css` z pomiarami `kontrasty.json` oraz przeliczenie obu par formułą WCAG |
| **D7** | Wartości kontrastu potrzebne w tabelach rozdz. 3–6, 14–15 wykraczają poza 33 pomiary z `kontrasty.json` | Doliczono brakujące pary **tą samą formułą WCAG 2.1** (luminancja względna). W rozdz. 16 prezentowany jest wyłącznie komplet 33 pomiarów źródłowych; wartości doliczone występują tylko w tabelach opisowych i są tam jednoznacznie przypisane do pary. | metoda identyczna ze źródłem; brak konfliktu wartości tam, gdzie pary się pokrywają |
| **D8** | Pomiar 13 (obrys kontrolki / tło, jasny) ma zadeklarowany próg **1,6**, który nie jest progiem WCAG | Werdykt nazwano **„próg systemowy"** zamiast AA/AAA i opisano jawnie: obrys `--dn-obrys-mocny` pełni funkcję **rozdzielenia powierzchni**, a nie jedynego nośnika rozpoznania kontrolki — kontrolkę identyfikuje wypełnienie, etykieta i pierścień fokusu (`--dn-fokus`, pomiar 14: 4,21 przy progu 3,0). Kwestia formalnego stosunku do kryterium 1.4.11 pozostaje **otwarta i wskazana do decyzji** prowadzącego system. | `kontrasty.json`, pomiar 13; `komponenty.css` — `.dn-pole-kontrolka`, `.dn-btn--zarys`; KANON rozdz. 9 |
| **D9** | `--dn-nakladka` motywu ciemnego ma wartość `rgba(0,0,0,0.62)` przy zakazie `#000000` | Zakwalifikowano jako **zgodne z zakazem**: czerń występuje wyłącznie jako krycie warstwy przyciemniającej, nigdy jako barwa powierzchni ani tekstu. Zapisano jako regułę czytania zakazu. | `zetony.css` rozdz. 11; KIERUNEK 3.1 („zakaz czystego `#000000` **jako tła i tekstu**") |
| **D10** | Kolejność rozstrzygania przy jednoczesnym `pointer: coarse` i `data-gestosc="przestronna"` nie jest w źródłach opisana | Ustalono na podstawie kolejności bloków w pliku: przy równej specyficzności rozstrzyga blok **późniejszy**, czyli `[data-gestosc='przestronna']`. Odnotowano, że kolizja nie daje różnicy widocznej (obie ścieżki dają 40 px kontrolki i 44 px wiersza). | `zetony.css`, rozdz. 12 i 13 pliku |
| **D11** | Punkty łamania są żetonami, ale CSS nie pozwala użyć `var()` w zapytaniu medialnym | Nazwano to **jedynym dopuszczonym miejscem powtórzenia liczby poza żetonem** i wskazano `matchMedia` jako ścieżkę dla warstwy skryptowej. | ograniczenie techniczne CSS; `zetony.css` rozdz. 7 |
| **D12** | Układ tła stanów w motywie ciemnym (alfa zamiast prymitywu) nie był opisany jako zasada | Sformułowano zasadę: **tło i obrys stanu w motywie ciemnym to barwa tekstu stanu z kryciem 0,14 i 0,34**, dzięki czemu dostosowują się do powierzchni, na której leżą. Podano wartości po złożeniu na `--dn-powierzchnia`. | `zetony.css` rozdz. 11; przeliczenie kompozycji alfa |
| **D13** | Asymetria uniesienia powierzchni między motywami (`--dn-panel` jaśniejszy od powierzchni w ciemnym, ciemniejszy w jasnym) mogła wyglądać na niekonsekwencję | Opisano jako **decyzję celową**: w motywie jasnym panel boczny cofa się, w ciemnym wysuwa. Obie ścieżki dają ten sam efekt hierarchii — panel odróżnia się od powierzchni roboczej. | zestawienie wartości `--dn-panel`, `--dn-powierzchnia`, `--dn-tlo` w obu motywach |
| **D14** | Kolejność aktualizacji plików przy zmianie wartości nie była nigdzie zapisana | Ustalono pięciostopniową kolejność (rozdz. 17.5) z wyróżnieniem kroku 3 (ponowny pomiar kontrastu) jako krytycznego. | KANON rozdz. 9 i 14; struktura pakietu |
| **D15** | Bramka wstępna przed dodaniem żetonu (trzy pytania) nie istniała w źródłach | Wprowadzono jako narzędzie zapobiegające rozrostowi rejestru; reguły R8 i R9 wyprowadzono z zasady „zero dekoracji bez funkcji" i z zamkniętego charakteru skal. | KIERUNEK rozdz. 1 i 4 |
| **D16** | Skala stopni pisma jest dwuczęściowa (przyrost 1 px do 14 px, mnożnikowy powyżej) — źródła podają wartości bez opisu kształtu | Opisano kształt skali i jego uzasadnienie (strefa robocza wymaga drobnych różnic, strefa ekspozycyjna wyraźnych skoków). Wartości pozostają bez zmian. | `zetony.css` rozdz. 3; KIERUNEK 3.7 (gęstość zwarta) |

### 19.1. Kwestie otwarte — do decyzji poza niniejszym opracowaniem

| Kwestia | Stan | Kto rozstrzyga |
|---|---|---|
| Formalny stosunek pomiaru 13 do kryterium WCAG 1.4.11 | próg zadeklarowany jako systemowy (1,6); wymaga potwierdzenia albo podniesienia obrysu | prowadzący system + osoba odpowiedzialna za dostępność |
| Rola `--dn-sygnal-800` | propozycja: wciśnięcie wypełnienia sygnałowego; brak deklaracji w źródłach | prowadzący system |
| Korekta komentarzy kontrastowych w `zetony.css` (D6) | rozbieżność redakcyjna zidentyfikowana; wartości barw poprawne | prowadzący system |
| Usunięcie zapasowego hexa w `komponenty.css` linia 163 (D5) | do rozważenia przy przeglądzie biblioteki | prowadzący system |
| Gęstość przestronna jako wybór Operatora w Oknie Ustawień | żetony gotowe; przełącznik niezadeklarowany w inwentarzu okien | prowadzący system + Właściciel |
| Wprowadzenie stopni rezerwowych (`szary-200`, `szary-700`) do użycia | wymaga nadania roli semantycznej | prowadzący system |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
