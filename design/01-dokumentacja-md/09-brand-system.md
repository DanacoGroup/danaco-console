# 09 · System marki — Danaco Console

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna **v2.0**) |
| **Warstwa funkcjonalna** | **Danaco Pilot** — Platforma AI Workspace OS (dokumentacja projektowa **v1.0**) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |
| **Dokument** | `09-brand-system.md` — system marki (opracowanie merytoryczno-techniczne) |
| **Wersja dokumentu** | 2.0 · **Status:** Deweloperski |
| **Data pakietu wizualnego** | 2026-08-11 |
| **Data dokumentacji funkcjonalnej** | 2026-08-06 |
| **Odbiorcy** | Właściciel produktu · zespół projektowy · zespół wdrożeniowy · projektant zewnętrzny · autor treści interfejsu |
| **Zakres** | pozycjonowanie · osobowość · obietnica · filary · architektura marki · elementy konstytutywne · barwa · terytoria · lista kontrolna spójności |
| **Poza zakresem** | szczegółowa księga znaku (pole ochronne, rozmiary minimalne, zakazy) — dokument `10` · żetony i skale — dokument `02` · komponenty — dokument `04` |
| **Źródła wiążące** | `WYNIK/KANON.md` · `design/01-kierunek/KIERUNEK.md` · pliki znaku w `WYNIK/zasoby/marka/**` · `zasoby/zetony/zetony.css` · `zasoby/zetony/kontrasty.json` · `dok/projekt-ui/**` |
| **Dokument towarzyszący** | `WYNIK/02-dokumentacja-html/09-brand-system.html` — interaktywny przewodnik po marce |

> **Zasada nadrzędna tego opracowania.** Nic, co nie ma pokrycia w dokumentacji źródłowej,
> nie pojawia się poniżej. Geometria znaku jest **odczytana z plików SVG**, nie odtworzona
> z pamięci. Rozstrzygnięcia własne — wyłącznie tam, gdzie źródła milczą — zebrano
> w rozdziale 14 „Decyzje projektowe".

---

## Spis treści

1. [Pozycjonowanie marki](#1-pozycjonowanie-marki)
2. [Producent, twórca i dwie warstwy produktu](#2-producent-twórca-i-dwie-warstwy-produktu)
3. [Idea przewodnia: instrument pomiarowy](#3-idea-przewodnia-instrument-pomiarowy)
4. [Osobowość marki — pięć cech](#4-osobowość-marki--pięć-cech)
5. [Obietnica marki i pętla produktu w znaku „»».”](#5-obietnica-marki-i-pętla-produktu-w-znaku-)
6. [Filary marki — cztery środowiska, cztery obietnice](#6-filary-marki--cztery-środowiska-cztery-obietnice)
7. [Zero blokad (ADL-017) jako wartość marki](#7-zero-blokad-adl-017-jako-wartość-marki)
8. [Architektura marki](#8-architektura-marki)
9. [Elementy konstytutywne marki](#9-elementy-konstytutywne-marki)
10. [Kropka sygnału — element sygnaturowy](#10-kropka-sygnału--element-sygnaturowy)
11. [Barwa marki](#11-barwa-marki)
12. [Terytoria marki](#12-terytoria-marki)
13. [Spójność — lista kontrolna zgodności z marką](#13-spójność--lista-kontrolna-zgodności-z-marką)
14. [Decyzje projektowe](#14-decyzje-projektowe)
15. [Załącznik A — inwentarz plików marki](#załącznik-a--inwentarz-plików-marki)
16. [Załącznik B — słowniczek marki](#załącznik-b--słowniczek-marki)

---

## 1. Pozycjonowanie marki

### 1.1. Zdanie pozycjonujące

> **Danaco Console to platforma operacyjna — kokpit dowodzenia — dla zawodowego
> Operatora zarządzającego cyfrową organizacją, zbudowana w języku monochromatycznej
> precyzji: odcienie bieli w motywie jasnym, odcienie czerni w motywie ciemnym,
> jeden chłodny sygnał; z ciążeniem ku własnemu systemowi „instrumentu pomiarowego":
> zwarta gęstość, dane krojem mono, zero dekoracji bez funkcji.**

Zdanie pochodzi z `KIERUNEK.md`, rozdz. 1, i jest w tym pakiecie **kontraktem**:
każda decyzja marki musi się z niego wyprowadzać. Dokument, który da się z tego zdania
wyprowadzić, jest zgodny z marką. Dokument, którego się nie da — nie jest.

### 1.2. Adresat: Operator

Marka nie mówi do „użytkownika". Mówi do **Operatora** — osoby, która nie rozmawia
z jednym modelem, tylko **zarządza cyfrową organizacją** złożoną z modeli, agentów,
projektów, procesów i automatyzacji. Ta różnica jest fundamentem całego systemu
wizualnego i wszystkich tekstów interfejsu.

| Wymiar | Adresat Danaco Console | Adresat, którym **nie jest** |
|---|---|---|
| Rola | Operator — dowodzi pracą wielu modeli i agentów | użytkownik czatu prowadzący jedną rozmowę |
| Sposób pracy | równolegle, w wielu kartach sesji i wielu oknach operacyjnych | sekwencyjnie, w jednym wątku |
| Oczekiwanie wobec narzędzia | ma być **czytelne pod obciążeniem** i przewidywalne | ma być efektowne i „przyjazne" |
| Reakcja na blokadę | traktuje ją jako wadę narzędzia | traktuje ją jako opiekę |
| Stosunek do gęstości | chce **więcej informacji na ekranie** | chce mniej i większe |
| Stosunek do danych | czyta identyfikatory, ścieżki i liczby | ich nie czyta |

### 1.3. Cztery dowody, które interfejs składa bez słów

Pozycjonowanie nie jest deklaracją w stopce — jest **udowadniane układem**.
Cztery dowody, wyprowadzone z architektury platformy (`KANON.md`, rozdz. 6–7):

| # | Dowód | Nośnik w interfejsie |
|---|---|---|
| 1 | **To nie jest czat** | rama kokpitu (pasek górny 48 px, atramentowy w obu motywach), pas kart sesji, okna operacyjne modułów, pas komunikacji jako jeden z paneli — nie jako całe okno |
| 2 | **Praca biegnie w tle** | kropka sygnału tętniąca (2,4 s) na kartach sesji i przy nadawcy piszącym |
| 3 | **Modele pracują dla siebie nawzajem** | asymetria koordynator–wykonawca; role MultitaskingAI: Executor 1, Executor 2, Coordinator, Executor 3 / Validator |
| 4 | **Wszystko jest konfigurowalne** | Okno konfiguracji (13 zakresów + „Nakładka i wywołanie modelu"), profile izolacji, Permissions Center — wszystkie jako **opcje**, nie bramy |

### 1.4. Czym marka nie jest — mapa odróżnień

Pozycjonowanie definiuje się także przez **odrzucone odruchy**. Katalog anty-domyślnych
(`KIERUNEK.md`, rozdz. 4) jest częścią pozycjonowania marki, nie tylko regulaminem CSS:

```
   ODRZUCONY ODRUCH                    POZYCJA MARKI
   ─────────────────────────────────────────────────────────────────────
   fioletowy gradient „AI"      ──►    monochrom + jeden błękit sygnałowy
   Inter + slate-900            ──►    Plex Sans + czysto neutralna skala
   trzy równe karty funkcji     ──►    strefy o malejącej masie
   wyszarzone przyciski-bramy   ──►    ZERO BLOKAD (ADL-017)
   stan samym kolorem           ──►    zawsze ikona albo etykieta
   emoji jako ikony             ──►    wyłącznie SVG z zestawu
   glassmorfizm wszędzie        ──►    powierzchnie kryjące
   #000000 jako tło             ──►    #0F0F0F / #F4F4F4
   Lorem ipsum, zmyślone dane   ──►    treści operacyjne z domeny produktu
```

Każdy z tych wierszy jest **testowalny**: da się sprawdzić w pliku, czy marka
została naruszona. To celowa cecha systemu — marka Danaco Console jest zapisana
tak, żeby dało się ją zweryfikować maszynowo, a nie tylko „poczuć".

---

## 2. Producent, twórca i dwie warstwy produktu

### 2.1. Podmiot i autorstwo

| Pole | Wartość |
|---|---|
| Producent | **Danaco Holding Group Sp. z o.o.** |
| Twórca | **Dariusz Naharnowicz** |
| Kontakt | support@danaco-group.pl |
| Status produktu | Deweloperski |
| Nota stopki (wiążąca) | `Danaco Console — AI Operating Environment · v2.0` / `© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz` |

### 2.2. Dwie warstwy jednego produktu

Produkt istnieje w dwóch warstwach dokumentacyjnych, wydanych w różnym czasie
i o różnym zakresie. **Nie są to dwa produkty ani dwie marki.**

| Warstwa | Nazwa w dokumentacji | Wersja | Data | Co rozstrzyga |
|---|---|---|---|---|
| **Funkcjonalna** | **Danaco Pilot** — Platforma AI Workspace OS | v1.0 | 2026-08-06 | środowiska, moduły, okna operacyjne, przepływ, role, mechanika kolejek, katalog komponentów, zasada zero blokad |
| **Wizualna** | **Danaco Console** — AI Operating Environment | v2.0 | 2026-08-11 | znak, barwa, typografia, żetony, gęstość, ruch, ikonografia, ton wizualny |

```
   ┌──────────────────────────────────────────────────────────────────┐
   │  DANACO PILOT  ·  warstwa funkcjonalna  ·  v1.0  ·  2026-08-06   │
   │  co system robi:  4 środowiska · 15 modułów · okna operacyjne ·  │
   │  przepływ E1–E12 · role MultitaskingAI · ADL-017 zero blokad     │
   └───────────────────────────────┬──────────────────────────────────┘
                                   │  obowiązuje nadal w całości
                                   ▼
   ┌──────────────────────────────────────────────────────────────────┐
   │  DANACO CONSOLE  ·  warstwa wizualna  ·  v2.0  ·  2026-08-11     │
   │  jak system wygląda i brzmi:  znak „Delegacja" · monochrom +     │
   │  jeden sygnał · Space Grotesk / Plex Sans / Plex Mono · żetony   │
   │  --dn-* · gęstość zwarta · tętno 2,4 s                           │
   └──────────────────────────────────────────────────────────────────┘
```

### 2.3. Zasada nazewnicza [WIĄŻĄCE]

1. Produkt nazywamy **Danaco Console**.
2. Nazwa **Danaco Pilot** występuje wyłącznie tam, gdzie cytujemy dokumentację
   funkcjonalną v1.0 — i wtedy jako **nazwa warstwy funkcjonalnej**, nie jako
   alternatywna nazwa produktu.
3. **Zakaz mieszania obu nazw w jednym zdaniu bez wyjaśnienia relacji.**
4. Nazwy własne środowisk, modułów, okien i ról pozostają **w brzmieniu
   dokumentacji** (`Studio Editor`, `Workflow Builder`, `Chat Window`,
   `Coordinator`, `Permissions Center`) — zakaz tłumaczenia i parafrazowania.
5. Podpis wersji na materiałach marki: **v2.0** (warstwa wizualna).

### 2.4. Co zostało zastąpione

Warstwa wizualna v2.0 **zastępuje w całości** poprzednią warstwę wizualną —
parę granat + złoto, godło „tarcza z belką wagi" i parę krojów Cormorant
Garamond + Inter — uznaną przez Właściciela za zastępczą i nieobowiązującą
(`KIERUNEK.md`, nagłówek „Zastępuje"). Poprzedni przewodnik marki
(`system-wizualny/pakiet/PRZEWODNIK-MARKI.md`, v1.0) zachowuje wartość
**wyłącznie historyczną i funkcjonalną** (opis relacji do rodziny Danaco);
jego wartości wizualne nie obowiązują.

| Element | v1.0 (nieobowiązujące) | v2.0 (obowiązujące) |
|---|---|---|
| Para barw | granat `#26395C` + złoto `#C9A24B` | monochrom + **jeden** sygnał `#3B6FE0` |
| Godło | tarcza kancelaryjna, monogram „D", belka wagi | sygnet **„Delegacja"** — podwójny grot + kropka |
| Nagłówki | Cormorant Garamond (szeryfowy) | **Space Grotesk** |
| Interfejs | Inter | **IBM Plex Sans** |
| Dane | JetBrains Mono | **IBM Plex Mono** |
| Charakter | „kancelaryjno-nowoczesny" | **instrument pomiarowy** |
| Godło pełnokolorowe | tak (gradienty w tarczy i monogramie) | **nie** — znak jest jednobarwny + kropka |

Ciągłość zachowały trzy rzeczy: **kokpit zamiast okna rozmowy**, **oba motywy
równoprawne** oraz **zasada zero blokad**. To one, a nie barwy, były nośnikiem
tożsamości.

---

## 3. Idea przewodnia: instrument pomiarowy

### 3.1. Metafora

Marka nie odwołuje się do metafory „asystenta", „mózgu" ani „iskry".
Odwołuje się do **instrumentu pomiarowego**: urządzenia, które ma jedno zadanie —
**pokazać stan rzeczy wiernie**, bez upiększania i bez opóźnienia.

Instrument pomiarowy ma trzy cechy, z których wyprowadzony jest cały system:

| Cecha instrumentu | Przełożenie na system wizualny | Żeton / wartość |
|---|---|---|
| **Skala jest gęsta** — instrument mieści dużo odczytów na małej powierzchni | gęstość zwarta 8/10: kontrolka 32 px, wiersz 36 px, pasek górny 48 px, odstępy paneli 12–16 px | `--dn-wym-kontrolka: 32px` · `--dn-wym-wiersz: 36px` · `--dn-wym-pasek: 48px` |
| **Odczyt jest maszynowy** — liczba ma być porównywalna, nie ładna | dane krojem mono z `tabular-nums`; identyfikatory, ścieżki, czasy, wersje | `--dn-ff-mono: IBM Plex Mono` |
| **Nic nie jest ozdobne** — każdy element tarczy coś znaczy | zero dekoracji bez funkcji: żaden gradient jako tło, żaden cień „domyślny", żaden ruch bez informacji | `--dn-grad-*` wyłącznie ilustracyjne |

### 3.2. Trzy pokrętła jako miara idei

Idea „instrumentu" jest w tym pakiecie **skwantyfikowana** — nie zostawiona
intuicji. Trzy pokrętła z `KIERUNEK.md`, rozdz. 2, są miarą marki:

| Pokrętło | Wartość | Co znaczy dla marki |
|---|---|---|
| `WARIANCJA_PROJEKTOWA` | **4/10** | marka jest rozpoznawalna przez **powtarzalność**, nie przez wariacje; asymetria wolno tylko tam, gdzie niesie hierarchię |
| `INTENSYWNOSC_RUCHU` | **3/10** | marka porusza się **raz na widok** i zawsze po coś; ruch jest informacją o stanie systemu |
| `GESTOSC_WIZUALNA` | **8/10** | marka ufa Operatorowi: pokazuje dużo, bo Operator umie czytać dużo |

```
   WARIANCJA      ├────●─────────────────────┤   4 / 10   przewidywalność
   RUCH           ├──●───────────────────────┤   3 / 10   jeden ruch znaczący
   GĘSTOŚĆ        ├───────────────────●──────┤   8 / 10   kokpit, nie strona
                  0                        10
```

### 3.3. Test idei — trzy pytania do każdego elementu

Element wchodzi do systemu marki, jeżeli przechodzi wszystkie trzy:

1. **Czy niesie informację?** Jeżeli nie — jest dekoracją i odpada.
2. **Czy zachowuje się tak samo w obu motywach?** Jeżeli nie — nie jest żetonem,
   tylko wyjątkiem, i wymaga uzasadnienia w „Decyzjach projektowych".
3. **Czy da się go zmierzyć?** Barwa ma mieć pomiar kontrastu, odstęp ma być
   wielokrotnością 4 px, czas ma być jednym z `--dn-czas-*`. Wartość spoza skali
   jest naruszeniem marki, nawet jeżeli wygląda dobrze.

---

## 4. Osobowość marki — pięć cech

Każda z pięciu cech jest **wyprowadzona z konkretnego rozdziału `KIERUNEK.md`**
i ma przypisany dowód w pliku. Cecha bez dowodu w pliku nie jest cechą marki,
tylko przymiotnikiem.

### 4.1. Precyzyjni — nie ozdobni

**Źródło:** `KIERUNEK.md` 1 („zero dekoracji bez funkcji"), 3.1 (neutralne czysto
neutralne), 3.6 (obrys ikon ujednolicony do **1,75** — „precyzja instrumentu").

| Jesteśmy | Nie jesteśmy |
|---|---|
| Obrys ikony **1,75** — jedna wartość w całym zestawie | zestaw ikon o obrysach 1,5 / 2 / 2,5 „dla urozmaicenia" |
| Neutralne o **równych składowych RGB** (`#7C7C7C`, `#616161`) | neutralne podbarwione na slate albo beż |
| Gradient **wyłącznie ilustracyjny** (awatar bez zdjęcia, rdzeń AOD, grafiki brandowe) | gradient jako tło przycisku, karty albo sekcji |
| Cień z dwóch kompletów tonowanych per motyw | jeden „domyślny szary cień" pod wszystkim |
| Wierzchołki grotów sygnetu **ostre** — jedyny celowo ostry element systemu | zaokrąglenie znaku „żeby był milszy" |

**Dowód w pliku:** `srodowiska/*.svg` — `stroke-width="1.75"` w każdym z czterech
emblematów; `logo/sygnet.svg` — ścieżki grotów bez ani jednego łuku.

### 4.2. Zwarci — nie ciaśni

**Źródło:** `KIERUNEK.md` 2 (`GESTOSC_WIZUALNA` 8/10), 3.7 (gęstość zwarta;
`pointer: coarse` → kontrolki 40 px „żeton, nie wyjątek").

| Jesteśmy | Nie jesteśmy |
|---|---|
| Kontrolka **32 px**, wiersz **36 px** — dużo odczytów na ekranie | kontrolka 48 px i trzy pola na widok |
| Odstępy paneli **12–16 px**, jednostka **4 px** [NIENEGOCJOWALNE] | odstępy „na oko", 13 px, 18 px, 22 px |
| Na dotyku kontrolka rośnie do **40 px**, wiersz do **44 px** — **żetonem** | „wersja mobilna" jako osobny zestaw reguł |
| Cel dotykowy zawsze osiągalny; check 16 → 20 px, przełącznik 36×20 → 44×24 | gęstość okupiona nietrafialnym celem |
| Wariant przestronny **przygotowany** w żetonach (`data-gestosc`) | gęstość jako jedyna możliwość |

Granica między „zwarty" a „ciasny" jest w tym systemie **mierzalna**: zwarte jest to,
co mieści się w skali 4 px i zachowuje cel dotykowy; ciasne jest to, co łamie
którykolwiek z tych dwóch warunków.

### 4.3. Sprawczy — nie hałaśliwi

**Źródło:** `KIERUNEK.md` 2 (`INTENSYWNOSC_RUCHU` 3/10 — „jeden ruch znaczący:
tętno pracy w tle"), 3.1 („czerń działa, sygnał wskazuje"), 3.4 (tętno jako jedyny
ruch ciągły).

| Jesteśmy | Nie jesteśmy |
|---|---|
| **Jeden ruch ciągły** w całym interfejsie — tętno kropki 2,4 s | pulsujące, przesuwające się i migające elementy jednocześnie |
| Mikroprzejścia **100–220 ms**, `cubic-bezier(0.2, 0, 0, 1)` | animacja wejścia 600 ms „dla efektu" |
| Przycisk główny = **inwersja atramentu** — działa, nie krzyczy | przycisk główny w kolorze akcentu na każdej karcie |
| Sygnał zajmuje **≤ 5 %** powierzchni — wskazuje, gdzie patrzeć | akcent na wszystkim, czyli na niczym |
| Ruch **niesie informację o stanie systemu** | parallaks, scrollytelling, animacje dekoracyjne w oknie roboczym |
| Przy `prefers-reduced-motion` tętno zastępuje **pierścień statyczny** — informacja zostaje | wyłączenie animacji kosztem utraty informacji |

### 4.4. Techniczni — nie zimni

**Źródło:** `KIERUNEK.md` 3.2 (Plex Sans — „doskonała czytelność w małych stopniach
i wzorowe polskie diakrytyki"), 3.6 (nazewnictwo ikon **polskie, opisowe**),
rozdz. 9 KANON-u (WCAG 2.1 AA jako **warunek wejściowy**, nie cel).

| Jesteśmy | Nie jesteśmy |
|---|---|
| Mono dla **danych** — identyfikatorów, ścieżek, czasów, wersji | mono dla całej treści „bo konsola" |
| Plex Sans dla **pracy** — pełne `latin-ext`, `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` sprawdzone | krój bez polskich diakrytyków „poprawiony" syntezą |
| Nazwy ikon **po polsku, opisowo** (`wykonawca`, `ptaszek-kolo`, `tabela-danych`) | `icon-check-circle-alt-2` |
| **33 zmierzone pary** kontrastu — pomiar jest wejściem do palety | „wygląda na czytelne" |
| Stan zawsze z **ikoną albo etykietą** — nie tylko barwą | poleganie na rozróżnianiu barw |
| Komunikat mówi, **co się stało i co dalej** | „Wystąpił błąd" |

Techniczność marki nie polega na estetyce terminala, tylko na **rygorze pomiaru**.
Ciepło bierze się stąd, że system jest zrozumiały, opisany po polsku i nie zostawia
Operatora bez odpowiedzi.

### 4.5. Przewidywalni — nie nudni

**Źródło:** `KIERUNEK.md` 2 (`WARIANCJA_PROJEKTOWA` 4/10 — „kokpit wymaga
przewidywalności; asymetria tylko tam, gdzie niesie hierarchię"), 3.3 (rama kokpitu
jako stały dom godła), 3.5 (znak opowiadający pętlę produktu).

| Jesteśmy | Nie jesteśmy |
|---|---|
| Rama kokpitu **atramentowa w obu motywach** — stałe stanowisko dowodzenia | pasek zmieniający charakter przy każdym przełączeniu motywu |
| Ten sam komponent zachowuje się **tak samo w każdym oknie** | wariant przycisku właściwy jednemu ekranowi |
| Asymetria **tylko** tam, gdzie niesie hierarchię: strefy strony głównej, para koordynator–wykonawca | asymetria dla „dynamiki kompozycji" |
| Znak, który **opowiada pętlę produktu** — „»»." zamiast litery w tarczy | monogram bez treści |
| Emblemat środowiska: jedna wypełniona kropka — **reguła, nie wyjątek** | cztery emblematy według czterech pomysłów |

Przewidywalność jest tu **funkcją zaufania**: Operator, który wie, gdzie co jest,
pracuje szybciej. Nienudność bierze się z **treści znaku**, nie z wariacji formy.

### 4.6. Zestawienie cech — jedna tabela

| # | Cecha | Jesteśmy | Nie jesteśmy | Twardy dowód |
|---|---|---|---|---|
| 1 | Precyzja | jedna wartość obrysu, czyste neutralne | ozdobność, gradient jako tło | `stroke-width="1.75"` |
| 2 | Zwartość | 32 / 36 / 48 px, jednostka 4 px | ciasnota, nietrafialny cel | `--dn-od-*` = 0/4/8/12/16… |
| 3 | Sprawczość | jeden ruch ciągły, sygnał ≤ 5 % | hałas wizualny | `--dn-czas-tetno: 2.4s` |
| 4 | Techniczność | pomiar, mono dla danych, polskie nazwy | zimno, ozdobna „konsolowość" | `kontrasty.json` — 33 pomiary |
| 5 | Przewidywalność | stała rama, stałe komponenty | monotonia formy zamiast treści | `--dn-rama` stała w obu motywach |

---

## 5. Obietnica marki i pętla produktu w znaku „»».”

### 5.1. Obietnica

> **Wydajesz polecenie — praca rusza. Widzisz, że biegnie. Nic cię nie zatrzymuje.**

Obietnica ma trzy człony i każdy ma nośnik w systemie:

| Człon obietnicy | Nośnik wizualny | Nośnik funkcjonalny |
|---|---|---|
| **Wydajesz polecenie** | grot promptu „»" — sygnet, `.dn-prompt-grot`, pole wpisywania | Chat Window wspólne wszystkim modułom |
| **Praca rusza i biegnie** | **kropka sygnału** + tętno 2,4 s | karty sesji w tle, kolejka, monitor procesu |
| **Nic cię nie zatrzymuje** | brak stanu wyłączonego w całej bibliotece komponentów | ADL-017 — zero blokad |

### 5.2. Znak jako zapis obietnicy

Sygnet **„Delegacja"** czyta się dosłownie jako trzy znaki: **„»»."**

```
   siatka 96 × 96

        26 ─────  ▲          ▲                         górna krawędź grotów
                 ╱│         ╱│
                ╱ │        ╱ │
        48 ─── ◄  │  ◄     │  ◄        ●               oś symetrii · kropka
                ╲ │        ╲ │        (83; 63,5) r 6,5
                 ╲│         ╲│
        70 ─────  ▼          ▼         ─────────       linia bazowa = spód kropki

              GROT 1      GROT 2      KROPKA
              12 … 44     40 … 72     76,5 … 89,5
              ramię 12    ramię 12    Ø 13
              ├── rytm 28 ──┤
```

| Element znaku | Znaczenie | Kto to jest w produkcie |
|---|---|---|
| **Grot 1** | polecenie wydane | Operator / koordynator |
| **Grot 2** | polecenie przyjęte i przekazane dalej | wykonawca / model |
| **Kropka** | praca, która właśnie ruszyła | proces w tle, sesja czynna |

Znak nie przedstawia litery ani przedmiotu — **przedstawia zdarzenie**.
To rozstrzygnięcie z `KIERUNEK.md` 3.5: „znak opowiada pętlę, nie literę".

### 5.3. Pętla produktu

```
   ┌───────────────┐   polecenie   ┌───────────────┐   wykonanie   ┌──────────────┐
   │   OPERATOR    │ ────  »  ───► │  KOORDYNATOR  │ ────  »  ───► │  WYKONAWCA   │
   │  (kokpit)     │               │  (Coordinator)│               │ (Executor)   │
   └───────────────┘               └───────────────┘               └──────┬───────┘
           ▲                                                              │
           │                        wynik + dowód pracy                   │
           └──────────────────────────────  ●  ◄───────────────────────────┘
                                       kropka sygnału
                                    „praca ruszyła / biegnie"
```

Pętla jest jednocześnie: **treścią znaku**, **mechaniką środowiska MultitaskingAI**
(role Coordinator / Executor / Validator) i **wzorcem interakcji** w każdym oknie
komunikacji. Jeden rysunek opisuje markę, architekturę i interfejs — to jest
kryterium, po którym poznaje się znak wyprowadzony z produktu, a nie dobrany do niego.

### 5.4. Wariant uproszczony a obietnica

Poniżej **24 px** obowiązuje wariant uproszczony: **jeden grot + większa kropka**.
Znika grot drugi, czyli człon „przekazanie". Zostaje to, co jest nieredukowalne:
**polecenie i praca, która ruszyła**. To rozstrzygnięcie hierarchii znaczeń —
gdyby usunąć kropkę, znak przestałby być znakiem tej marki.

| Wariant | Ścieżka grotu | Kropka | Stosunek kropki do wysokości znaku |
|---|---|---|---|
| Pełny (≥ 24 px) | dwa groty, 12…44 i 40…72, wysokość 44 | `cx 83 · cy 63,5 · r 6,5` (Ø 13) | 13 / 44 ≈ **0,30** |
| Uproszczony (< 24 px) | jeden grot, 18…64, wysokość 60 | `cx 79 · cy 69 · r 9` (Ø 18) | 18 / 60 = **0,30** |

Stosunek jest **identyczny** — 0,30. Wariant uproszczony nie jest innym znakiem,
tylko tym samym znakiem o zmienionej liczbie członów.

---

## 6. Filary marki — cztery środowiska, cztery obietnice

Marka nie ma abstrakcyjnych filarów. Ma **cztery środowiska**, a każde jest
obietnicą złożoną Operatorowi. Motta i opisy pochodzą z makiety Centrum dowodzenia
(E1) oraz z `dok/projekt-ui/przeplyw/przeplyw-okien.md`, rozdz. 5 — **nie są
wymyślone na potrzeby tego dokumentu**.

### 6.1. TalkIn — „Myśl. Analizuj. Rozumiej."

| Pole | Wartość |
|---|---|
| Motto (E1) | **Myśl. Analizuj. Rozumiej.** |
| Opis na karcie (E1) | Wiedza, treść, dokumenty, badania i tłumaczenia. |
| Opis trybu (dok. przepływu) | Wiedza, komunikacja i praca z treścią |
| Powłoka | pasek górny · pas kart sesji · boczna nawigacja **9 modułów** · obszar roboczy |
| Moduły | Studio · Research · Library · Translate · Browser · Assistant · Roundtable · Workspace · Agents |
| Emblemat | `srodowisko-talkin.svg` — dymek rozmowy, dwie linie tekstu, kropka sygnału `cx 16,2 · cy 10,8 · r 1,5` |
| **Obietnica marki** | **Zrozumiesz materiał, zanim podejmiesz decyzję.** Środowisko zdejmuje z Operatora pracę czytania, porównywania i tłumaczenia — a nie podejmuje decyzji za niego. |

### 6.2. WorkSpace — „Planuj. Organizuj. Realizuj."

| Pole | Wartość |
|---|---|
| Motto (E1) | **Planuj. Organizuj. Realizuj.** |
| Opis na karcie (E1) | Projekty, procesy, automatyzacje i produkty. |
| Opis trybu (dok. przepływu) | Produktywność, organizacja i realizacja projektów |
| Powłoka | jak wyżej · boczna nawigacja **9 modułów** |
| Moduły | Studio · Research · Library · Browser · Roundtable · Workspace · Design · Apps · Agents |
| Emblemat | `srodowisko-workspace.svg` — cztery zaokrąglone pola 7,2 × 7,2, kropka sygnału w polu prawym dolnym `cx 16,9 · cy 16,9 · r 1,6` |
| **Obietnica marki** | **Praca nie zginie między narzędziami.** Projekt, jego kontekst, biblioteka i automatyki są jednym obiektem, nie czterema. |

Emblemat WorkSpace niesie dodatkową treść zgodną z ideą: **cztery pola, jedno
z kropką** — praca biegnie w jednym z nich, reszta czeka. To ta sama figura
myślowa, co „cztery środowiska jednego systemu".

### 6.3. CodeStudio — „Projektuj. Buduj. Rozwijaj."

| Pole | Wartość |
|---|---|
| Motto (E1) | **Projektuj. Buduj. Rozwijaj.** |
| Opis na karcie (E1) | Programowanie, terminale i architektura systemów. |
| Opis trybu (dok. przepływu) | Programowanie |
| Powłoka | jak wyżej · boczna nawigacja **8 modułów** |
| Moduły | Roundtable · Workspace · Design · Apps · Terminal · Developer · Diagnostics · Agents |
| Emblemat | `srodowisko-codestudio.svg` — okno 18 × 16 (r 2,5), grot `>` , linia wiersza, kropka sygnału `cx 18,4 · cy 9,9 · r 1,5` |
| **Obietnica marki** | **Zobaczysz, co system naprawdę robi.** Terminal, build, logi i diagnostyka są widoczne wprost — narzędzie nie ukrywa własnego działania. |

Emblemat CodeStudio zawiera **ten sam grot**, co sygnet marki — w skali ikony.
To jedyny emblemat, który cytuje znak nadrzędny; uzasadnienie: CodeStudio jest
środowiskiem, w którym polecenie jest wydawane najbardziej dosłownie.

### 6.4. MultitaskingAI — „Deleguj. Koordynuj. Nadzoruj."

| Pole | Wartość |
|---|---|
| Motto (E1) | **Deleguj. Koordynuj. Nadzoruj.** |
| Opis na karcie (E1) | Orkiestracja zespołów modeli i procesy autonomiczne. |
| Opis trybu (dok. przepływu) | Orkiestracja autonomicznej pracy ciągłej |
| Powłoka | pasek górny · pas kart sesji · **panel orkiestracji (6 sekcji)** zamiast bocznej nawigacji modułów · obszar roboczy |
| Role zespołu | Executor 1 · Executor 2 · Coordinator · Executor 3 / Validator · Subagent Network |
| Mechanika | silnik kolejek — **11 akcji** · integracja Automations 24/7 · warstwa centralna Always On Display |
| Emblemat | `srodowisko-multitaskingai.svg` — węzeł centralny **wypełniony** `r 2,6`, trzy węzły obrysowe `r 1,9`, trzy połączenia |
| **Obietnica marki** | **Możesz oddać pracę i nadal nad nią panować.** Delegacja nie oznacza utraty widoczności — rdzeń AOD i monitor procesu pokazują, co dzieje się bez Operatora. |

MultitaskingAI to **środowisko, które jest znakiem marki rozwiniętym do systemu**:
sygnet pokazuje jedno przekazanie, panel orkiestracji pokazuje ich wiele naraz.
Dlatego w tym emblemacie kropką sygnału jest **węzeł centralny** — koordynator —
a nie punkt na obrzeżu.

### 6.5. Cztery obietnice — zestawienie

| Środowisko | Motto | Obietnica w jednym zdaniu | Emblemat: gdzie jest kropka |
|---|---|---|---|
| **TalkIn** | Myśl. Analizuj. Rozumiej. | Zrozumiesz materiał, zanim podejmiesz decyzję | koniec drugiej linii tekstu w dymku |
| **WorkSpace** | Planuj. Organizuj. Realizuj. | Praca nie zginie między narzędziami | pole prawe dolne z czterech |
| **CodeStudio** | Projektuj. Buduj. Rozwijaj. | Zobaczysz, co system naprawdę robi | prawa górna część okna, obok grotu |
| **MultitaskingAI** | Deleguj. Koordynuj. Nadzoruj. | Możesz oddać pracę i nadal nad nią panować | węzeł centralny |

**Reguła emblematu [WIĄŻĄCE]:** siatka 24 × 24 · obrys **1,75** · wypełnienie `none`
· barwa `currentColor` · zakończenia i łączenia `round` · **dokładnie jedna wypełniona
kropka sygnału**. Emblemat, który ma zero kropek albo dwie, nie należy do systemu.

---

## 7. Zero blokad (ADL-017) jako wartość marki

### 7.1. Zapis zasady

> **Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych
> zgód. Domyślne zachowanie systemu to wykonanie polecenia.**

Izolacja techniczna, izolacja kontekstu, uprawnienia, profile izolacji
i Permissions Center są **wyłącznie opcjami konfigurowalnymi** przez Operatora,
ze stanem wyjściowym „wyłączone / pełny dostęp".

### 7.2. Dlaczego to jest wartość marki, a nie reguła UI

Zasada rozstrzyga **stosunek producenta do Operatora**. System, który blokuje,
zakłada, że użytkownik jest zagrożeniem. System, który wykonuje i opisuje
konsekwencje, zakłada, że użytkownik jest **zawodowcem odpowiadającym za swoją
pracę**. Danaco Console przyjmuje drugie założenie — i to jest zdanie o marce,
nie o CSS.

Z tego wynika mierzalna konsekwencja: **atrybut `disabled` nie występuje
w produkcie**. Nie jest odradzany — nie występuje.

### 7.3. Objawy zasady w trzech warstwach

#### A. W interfejsie

| Sytuacja | Zachowanie zgodne z marką | Zachowanie naruszające markę |
|---|---|---|
| Warunek formularza niespełniony | przycisk klikalny; po naciśnięciu **komunikat, czego brakuje** | przycisk wyszarzony |
| Trwa walidacja | kontrolka fokusowalna, `aria-busy`, opis obok | zablokowany formularz |
| Brak uprawnienia | opis w miejscu użycia (dymek, tekst pomocniczy, plakietka) | usunięcie kontrolki z obiegu interakcji |
| Odliczanie „Wyślij ponownie" | licznik **informacyjny**, kliknięcie działa | licznik blokujący |
| Metoda uwierzytelniania wyłączona w konfiguracji | **w ogóle nierenderowana** — mniej segmentów | segment zablokowany |
| Akcja o dużej wadze | wariant wizualny sygnalizujący ciężar; ewentualny modal potwierdzenia jako **wzorzec opcjonalny** | obowiązkowa bramka |
| Wartość odziedziczona z szerszego poziomu izolacji | **jawnie nadpisywalna** w każdym widoku | zablokowana do edycji |
| Przycisk „Pomiń" w oknie logowania | obecny i klikalny **w obu fazach wdrożenia** | ukryty albo zablokowany |

Stany realizuje się atrybutami ARIA — `aria-pressed`, `aria-busy`, `aria-invalid`,
`aria-selected`, `aria-current` — **nigdy przez `disabled`**.

#### B. W języku

Zasada rządzi tonem komunikatów. Zdanie interfejsu ma **opisywać stan i wskazywać
następny krok**, nie odmawiać.

| Zamiast | Marka mówi |
|---|---|
| „Nie masz uprawnień do tej operacji." | „Operacja wymaga dostępu do katalogu roboczego. Włącz go w Oknie konfiguracji → Dostępy." |
| „Formularz zawiera błędy." | „Pole „Nazwa projektu" jest puste — zapis pominie tę wartość." |
| „Poczekaj 30 s." | „Kolejna wiadomość może zostać wysłana za 30 s. Wysyłka teraz jest możliwa." |
| „Funkcja niedostępna w tym środowisku." | „Terminal działa w CodeStudio. Otworzyć CodeStudio z tą sesją?" |
| „Czy na pewno chcesz kontynuować?" | „Usunięcie skasuje 3 pliki repozytorium sesji. Usuwam." + odwołanie |

Trzy reguły języka wyprowadzone z ADL-017:

1. **Nie odmawiaj — opisz warunek i drogę.** Komunikat kończy się czynnością,
   nie zakazem.
2. **Nie strasz.** Waga akcji jest opisana faktem (co zostanie zmienione), nie
   przymiotnikiem („nieodwracalne!", „uwaga!").
3. **Nie proś o zgodę, której nie potrzebujesz.** Potwierdzenie pojawia się tylko
   tam, gdzie skutek jest nieodwracalny — i wtedy podaje **zakres skutku liczbowo**.

#### C. W komunikacji zewnętrznej

Marka nie obiecuje „bezpieczeństwa przez ograniczenie". Obiecuje **kontrolę
przez konfigurowalność i widoczność**. W materiałach oznacza to:

- nie używamy sformułowań „system nie pozwoli", „zabezpieczy przed",
- używamy: „Operator decyduje", „wartość wyjściowa: pełny dostęp",
  „każda izolacja jest opcją, nie warunkiem",
- funkcje bezpieczeństwa opisujemy jako **narzędzia Operatora** (Permissions
  Center, punkty izolacji, Panel prowenancji), nie jako strażników produktu.

### 7.4. Test zgodności z ADL-017

Opracowanie jest zgodne, jeżeli w pliku:

- [ ] nie występuje ciąg `disabled`,
- [ ] nie występuje `aria-disabled="true"` jako substytut blokady,
- [ ] nie występuje `pointer-events: none` na kontrolce,
- [ ] nie występuje `cursor: not-allowed`,
- [ ] każdy komunikat odmowy zawiera **następny krok**,
- [ ] każdy licznik czasu jest opisany jako **informacyjny**.

---

## 8. Architektura marki

### 8.1. Model hierarchii

Marka ma **cztery poziomy**, o malejącym zakresie i rosnącej liczbie wystąpień.
Poziom niższy nigdy nie zastępuje wyższego; poziom wyższy nigdy nie schodzi do roli
niższego.

```
   POZIOM 0 · ZNAK NADRZĘDNY
   ┌──────────────────────────────────────────────────────────────────────┐
   │  sygnet „Delegacja"  +  logotyp  DANACO / CONSOLE                    │
   │  siatka 96 × 96 · barwy własne · nie przełącza się z motywem         │
   │  warianty: poziomy · pionowy · sam logotyp · sam sygnet · mono ×2    │
   └───────────────┬──────────────────────────────────────────────────────┘
                   │  dziedziczy: kropka sygnału, monochrom, jeden akcent
   ┌───────────────┴──────────────────────────────────────────────────────┐
   │  POZIOM 1 · EMBLEMATY ŚRODOWISK                                      │
   │  siatka 24 × 24 · obrys 1,75 · currentColor · JEDNA wypełniona kropka│
   │  TalkIn · WorkSpace · CodeStudio · MultitaskingAI                    │
   └───────────────┬──────────────────────────────────────────────────────┘
                   │  dziedziczy: siatka 24, obrys 1,75, currentColor
   ┌───────────────┴──────────────────────────────────────────────────────┐
   │  POZIOM 2 · ZESTAW IKON INTERFEJSU (Lucide, 82 nazwy)                │
   │  siatka 24 × 24 · obrys 1,75 · fill none · currentColor              │
   │  BEZ kropki sygnału — kropka jest zarezerwowana dla marki i stanu    │
   └───────────────┬──────────────────────────────────────────────────────┘
                   │  dziedziczy: geometrię sygnetu, barwy własne
   ┌───────────────┴──────────────────────────────────────────────────────┐
   │  POZIOM 3 · TOŻSAMOŚCI SYSTEMOWE                                     │
   │  ikona aplikacji (kafel 1024, r 224) · ikona maskowalna (r 0)        │
   │  favicon SVG + ICO 16/32/48 · apple-touch-icon · manifest            │
   └──────────────────────────────────────────────────────────────────────┘
```

### 8.2. Reguły dziedziczenia

| Reguła | Treść |
|---|---|
| **R1 · Jedna kropka** | Kropka sygnału występuje **raz** w znaku i **raz** w każdym emblemacie środowiska. W zestawie ikon interfejsu nie występuje w ogóle. |
| **R2 · Barwy własne w górę** | Poziom 0 i poziom 3 mają **barwy własne** i nie reagują na motyw (poza wariantami przeznaczonymi na ciemne tło). Poziomy 1 i 2 dziedziczą barwę tekstu przez `currentColor`. |
| **R3 · Siatka niżej** | Poziom 0 rysuje się na siatce **96**, poziomy 1–2 na siatce **24**. Znaku nie skaluje się do siatki ikon — poniżej 24 px wchodzi wariant uproszczony. |
| **R4 · Emblemat nie zastępuje znaku** | Emblemat środowiska nigdy nie występuje jako godło produktu (np. w pasku górnym zamiast sygnetu). |
| **R5 · Znak nie zastępuje ikony** | Sygnetu nie stosuje się jako ikony akcji ani jako punktora listy. |
| **R6 · Tożsamości systemowe pochodzą ze znaku** | Ikona aplikacji i favicon są **wycinkami znaku**, nigdy osobnymi rysunkami. |

### 8.3. Skąd który plik pochodzi — łańcuch pochodzenia

```
   logo/sygnet.svg  (96×96, groty #181818, kropka #3B6FE0)
     ├─► logo/sygnet-na-ciemnym.svg      groty #F4F4F4 · kropka #5C8CEC
     ├─► logo/sygnet-mono-czarny.svg     wszystko #181818   (druk, grawer, pieczęć)
     ├─► logo/sygnet-mono-bialy.svg      wszystko #F4F4F4   (kontra, tłoczenie)
     ├─► logo/sygnet-uproszczony.svg     jeden grot, kropka r 9   (< 24 px)
     │      └─► favicon/favicon.svg      + @media prefers-color-scheme
     │             └─► favicon-16/32/48.png ─► favicon.ico
     ├─► logo/logo-poziomy.svg           sygnet ×0,6667 + logotyp (225×96)
     ├─► logo/logo-pionowy.svg           sygnet + logotyp pod spodem (159×172)
     ├─► logo/logotyp.svg                sam logotyp (149×56)
     └─► ikona-aplikacji/ikona-aplikacji.svg   kafel 1024, r 224, tło #131313
            ├─► ikona-maskowalna.svg     r 0, znak mniejszy (strefa bezpieczna)
            └─► ikona-180/192/512/1024.png · apple-touch-icon.png
```

### 8.4. Kiedy którego wariantu użyć

| Kontekst | Wariant | Uzasadnienie |
|---|---|---|
| Pasek górny aplikacji (48 px) | **sygnet-na-ciemnym** + logotyp tekstem | rama jest atramentowa w obu motywach |
| Nagłówek dokumentu, papier firmowy, stopka | **logo-poziomy** / `-na-ciemnym` | pełna nazwa produktu w jednej linii |
| Format kwadratowy, roll-up, awatar duży | **logo-pionowy** | proporcja bliższa 1 : 1 |
| Gdy sygnet już występuje obok | **logotyp** | zakaz podwajania sygnetu w jednym układzie |
| Druk jednobarwny, grawer, pieczęć, tłoczenie | **sygnet-mono-czarny** / **-bialy** | brak drugiej farby |
| Kafel aplikacji, PWA, pulpit | **ikona-aplikacji** (r 224) | pełny znak na własnym tle |
| Android maskable | **ikona-maskowalna** (r 0) | system sam przycina kształt |
| Karta przeglądarki, favicon, awatar < 24 px | **sygnet-uproszczony** | próg czytelności dwóch grotów |

### 8.5. Rozstrzygnięcie: ikona aplikacji ma własne tło

Ikona aplikacji nie jest przezroczysta. Stoi na polu `#131313` — czyli na wartości
`--dn-szary-925`, **tej samej, co rama kokpitu**. To celowe: kafel na pulpicie jest
zapowiedzią paska górnego. Promień `224 / 1024 ≈ 21,9 %` daje kształt zbliżony
do systemowego kafla, bez naśladowania konkretnego systemu operacyjnego.

| Plik | Pole | Promień | Skala znaku | Zajętość kafla |
|---|---|---|---|---|
| `ikona-aplikacji.svg` | 1024 × 1024, `#131313` | 224 | ×6,6133 (offset 194,6) | ≈ **62 %** |
| `ikona-maskowalna.svg` | 1024 × 1024, `#131313` | 0 | ×5,5467 (offset 245,8) | ≈ **52 %** — mieści się w strefie bezpiecznej |

---

## 9. Elementy konstytutywne marki

Sześć elementów. Usunięcie któregokolwiek sprawia, że opracowanie przestaje być
rozpoznawalne jako Danaco Console.

### 9.1. Sygnet

| Cecha | Wartość odczytana z pliku |
|---|---|
| Siatka | 96 × 96, `viewBox="0 0 96 96"` |
| Grot 1 | `M12 26 H24 L44 48 L24 70 H12 L32 48 Z` |
| Grot 2 | `M40 26 H52 L72 48 L52 70 H40 L60 48 Z` |
| Kropka | `cx 83 · cy 63,5 · r 6,5` |
| Wysokość grotu | 44 (26 → 70) · oś symetrii `y = 48` |
| Ramię grotu | 12 (poziomo) · wcięcie karbu 20 (od 12 do 32) |
| Rytm grotów | **28** jednostek (12 → 40) |
| Odstęp grot 2 → kropka | 4,5 (72 → 76,5) |
| Styczność | spód kropki `63,5 + 6,5 = 70` = **linia bazowa grotów** |
| Barwy (jasne tło) | groty `#181818` · kropka `#3B6FE0` |
| Barwy (ciemne tło) | groty `#F4F4F4` · kropka `#5C8CEC` |
| Krzywe | **zero łuków w grotach** — wyłącznie odcinki proste; kropka jest idealnym kołem |

**Zakaz bezwzględny:** krzywych sygnetu nie wolno modyfikować. Wolno budować
warianty, konfiguracje i zastosowania — nie wolno zmieniać ani jednej współrzędnej.

### 9.2. Logotyp

| Cecha | Wartość |
|---|---|
| Wers 1 | `DANACO` — **Space Grotesk 700** |
| Wers 2 | `CONSOLE` — **IBM Plex Mono 500**, rozstrzelone wersaliki |
| Zakończenie wersu 2 | **kropka sygnału** `r 3,4` (`#3B6FE0`) |
| Plik samodzielny | `logotyp.svg` — 149 × 56 |
| W układzie poziomym | logotyp od `x = 76`; wers 1 na `y = 51`, wers 2 na `y = 67` |
| Znaczenie pary krojów | **marka + maszyna** — Space Grotesk niesie tożsamość, Plex Mono niesie działanie |

Para krojów w logotypie jest miniaturą całego systemu typograficznego: to samo
rozróżnienie („nagłówek = wejście i tożsamość, mono = maszyna") rządzi każdym oknem.

### 9.3. Kropka sygnału

Opisana osobno w rozdziale 10 — jest elementem sygnaturowym, nie tylko składnikiem
znaku.

### 9.4. Para krojów

| Rola | Krój | Wagi | Żeton | Gdzie |
|---|---|---|---|---|
| Nagłówki, tytuły środowisk, logotyp | **Space Grotesk** | 500–700 | `--dn-ff-naglowek` | wejście do środowiska, tożsamość marki |
| Interfejs i treść | **IBM Plex Sans** | 400–700 | `--dn-ff-bazowa` | praca |
| Dane, identyfikatory, terminal | **IBM Plex Mono** | 400–600 | `--dn-ff-mono` | maszyna |

Wszystkie trzy: licencja **OFL 1.1**, podzbiory `latin` + `latin-ext`, pliki lokalne
w `zasoby/fonty/` (bez CDN). Liczby w mono zawsze `tabular-nums`.
Test diakrytyków: **ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ**.

**Reguła [WIĄŻĄCE]:** różnica krojów niesie znaczenie, nie ozdobę. Zamiana krojów
między rolami zaciera rozróżnienie architektoniczne i jest naruszeniem marki.

### 9.5. Monochrom + jeden sygnał

| Rola | Zakres | Udział ekranu | Ograniczenie |
|---|---|---|---|
| Neutralne (18 kroków, `#FFFFFF` → `#0A0A0A`) | tła, powierzchnie, tekst, obrysy | ≈ **95 %** | zakaz `#000000`; `#FFFFFF` wyłącznie jako powierzchnia karty i tekst na atramencie |
| Atrament (działanie) | przycisk główny — **inwersja** względem motywu | punktowo | — |
| Sygnał `#3B6FE0` i rodzina | fokus · wybór · odnośniki · postęp · kropka pracy | **≤ 5 %** | nigdy tło sekcji, nigdy przycisk główny |
| Stany | zieleń · bursztyn · czerwień · informacja = sygnał | punktowo | zawsze z ikoną albo etykietą |

### 9.6. Rama atramentowa

Pasek górny (48 px) jest **atramentowy w obu motywach** — `--dn-rama:
var(--dn-szary-925)`. To jedyny element interfejsu, który nie przełącza się
z motywem (obok godła).

Rama daje trzy rzeczy naraz:
1. **efekt stanowiska dowodzenia** — ciemna obwoluta wokół przełączalnej treści,
2. **stały dom dla godła**, wyszukiwarki i wskaźników stanu,
3. **stały punkt orientacji** — Operator zawsze wie, gdzie jest góra systemu.

Żetony ramy są odrębne od żetonów treści: `--dn-rama-tekst` (`szary-100`),
`--dn-rama-tekst-2` (`szary-400`), `--dn-rama-hover` (`rgba(255,255,255,.08)`),
`--dn-rama-obrys` (`rgba(255,255,255,.10)`). Komponent na ramie sięga po nie,
nie po żetony treści — inaczej zniknie w motywie jasnym.

---

## 10. Kropka sygnału — element sygnaturowy

### 10.1. Definicja

**Wypełniony punkt błękitu sygnałowego oznaczający „tu biegnie praca".**
Jedna rzecz, po której platformę się zapamięta (`KIERUNEK.md` 3.4).

### 10.2. Pięć miejsc wystąpienia

| # | Miejsce | Postać | Co komunikuje |
|---|---|---|---|
| 1 | **Godło** | kropka zamykająca podwójny grot `»»` | zlecenie przekazane — praca ruszyła |
| 2 | **Emblematy czterech środowisk** | dokładnie jedna wypełniona kropka na siatce 24 | to środowisko jest miejscem, gdzie praca biegnie |
| 3 | **Karty sesji** | kropka **pulsująca** (`.dn-kropka--tetno`) | proces w tle tej sesji |
| 4 | **Okno komunikacji** | kropka przy nadawcy aktywnie piszącym | ten rozmówca właśnie pracuje |
| 5 | **Always On Display** | rdzeń awatara | system pracuje, choć nie patrzysz |

Piąte wystąpienie jest najważniejsze dla marki: **AOD to kropka powiększona
do rozmiaru warstwy**. Element sygnaturowy przestaje być detalem i staje się
całym widokiem.

### 10.3. Tętno — 2,4 s

| Parametr | Wartość |
|---|---|
| Żeton | `--dn-czas-tetno: 2.4s` |
| Status | **jedyny ruch ciągły w całym interfejsie** |
| Krzywa | `--dn-ease: cubic-bezier(0.2, 0, 0, 1)` |
| Zachowanie przy `prefers-reduced-motion` | tętno wyłączone, zastąpione **pierścieniem statycznym** |
| Wymiar bazowy kropki | `--dn-wym-kropka: 6px` |

```
   krycie / skala
       │        ╭──╮                    ╭──╮
       │      ╭─╯  ╰─╮                ╭─╯  ╰─╮
       │   ╭──╯      ╰──╮          ╭──╯      ╰──╮
       │ ──╯            ╰──────────╯            ╰────►  czas
       └────────────────┬────────────────┬───────────
                      2,4 s            4,8 s
       jedno uderzenie na 2,4 s — rytm spokojnego pulsu, nie migotanie
```

Dlaczego 2,4 s: rytm szybszy odbierany jest jako alarm, wolniejszy — jako awaria
animacji. 2,4 s jest na tyle powolne, że **nie rozprasza przy pracy wielogodzinnej**,
i na tyle wyraźne, że **widać je kątem oka**. Przy `prefers-reduced-motion`
informacja nie znika — zmienia nośnik z ruchu na kształt (pierścień).

### 10.4. Czego kropka nie oznacza

| Kropka **nie** oznacza | Właściwy nośnik |
|---|---|
| błędu | `.dn-kropka--blad` + ikona + etykieta |
| powodzenia | `.dn-kropka--sukces` + ikona + etykieta |
| powiadomienia nieprzeczytanego | plakietka licznikowa |
| ozdoby, punktora, separatora | `·` w treści albo obrys |
| „element aktywny w nawigacji" | `aria-current` + podświetlenie |

Kropka **sygnałowa** (błękit) ma jedno znaczenie: **praca biegnie**. Warianty
stanowe (`--sukces`, `--ostrzezenie`, `--blad`, `--neutralna`) są osobnym
mechanizmem i zawsze występują z ikoną albo etykietą — nigdy samą barwą.

---

## 11. Barwa marki

### 11.1. Dlaczego monochrom

| Powód | Konsekwencja |
|---|---|
| **Kokpit pracuje długo** — Operator patrzy godzinami | neutralne tło nie męczy i nie konkuruje z treścią |
| **Treść jest kolorowa sama z siebie** — wykresy, podglądy, obrazy, diffy | interfejs musi być tłem dla treści, nie jej konkurentem |
| **Stany potrzebują wolnego kanału barwy** | jeżeli interfejs jest barwny, ostrzeżenie przestaje być widoczne |
| **Neutralne czysto neutralne** (równe składowe RGB) | brak podbarwienia slate/beż — „odcienie bieli i czerni" dosłownie |
| **Oba motywy równoprawne** | żaden nie jest „przyciemnioną" wersją drugiego; wartości definiowane osobno |

### 11.2. Dlaczego jeden sygnał

Jeden akcent to **jedna reguła czytania**: „błękit = system coś robi albo coś
wskazuje". Dwa akcenty wymagałyby od Operatora pamiętania, który jest który —
koszt poznawczy bez zysku informacyjnego.

| Rola sygnału | Żeton | Wartość |
|---|---|---|
| Fokus | `--dn-fokus` | `sygnal-500` (jasny) / `sygnal-400` (ciemny) |
| Tekst, odnośnik | `--dn-sygnal` | `sygnal-600` (jasny) / `sygnal-300` (ciemny) |
| Wypełnienie (postęp, wybór) | `--dn-sygnal-wypelnienie` | `sygnal-600` / `sygnal-500` |
| Tło informacyjne | `--dn-sygnal-tlo` | `sygnal-100` / `rgba(59,111,224,.16)` |
| Kropka sygnału | wartość znaku | `#3B6FE0` / `#5C8CEC` |

**Celowe scalenie:** rodzina stanu „informacja" **jest** rodziną sygnału.
System nie wprowadza szóstej barwy dla komunikatu informacyjnego — informacja
i „system pracuje" to ten sam rejestr.

### 11.3. Dlaczego działanie główne to inwersja atramentu

Przycisk główny jest **czarny na jasnym motywie i biały na ciemnym**.

| Zaleta | Wyjaśnienie |
|---|---|
| Najwyższy możliwy kontrast | 17,76 : 1 (biel na `#181818`) · 16,23 : 1 (`#0F0F0F` na `#ECECEC`) — pomiar z `kontrasty.json` |
| Brak drugiej barwy | monochromatyczna inwersja jest najsilniejszym akcentem, jaki system może wykonać, nie wydając ani jednego kanału barwy |
| Sygnał zostaje wolny | gdyby przycisk główny był błękitny, błękit oznaczałby jednocześnie „działanie" i „praca biegnie" |
| Odróżnienie od odruchu | „przycisk główny w kolorze marki" jest domyślnym odruchem, którego katalog anty-domyślnych zakazuje |

**Reguła [WIĄŻĄCE]:** czerń działa, sygnał wskazuje.

### 11.4. Pomiar zamiast deklaracji

Paleta powstała **z pomiaru**: wartość, która nie przechodzi progu WCAG 2.1 AA,
nie wchodzi do żetonów. Komplet: **33 pomiary** w `zasoby/zetony/kontrasty.json`.

| Para (wybór) | Wartości | Kontrast | Próg |
|---|---|---|---|
| tekst / tło (jasny) | `#181818` / `#F4F4F4` | **16,14 : 1** | 4,5 |
| tekst / tło (ciemny) | `#ECECEC` / `#0F0F0F` | **16,23 : 1** | 4,5 |
| biel / przycisk atrament (jasny) | `#FFFFFF` / `#181818` | **17,76 : 1** | 4,5 |
| sygnał-tekst / tło (jasny) | `#2457C9` / `#F4F4F4` | **5,82 : 1** | 4,5 |
| sygnał-tekst / tło (ciemny) | `#8FB2F5` / `#0F0F0F` | **9,00 : 1** | 4,5 |
| fokus sygnał / tło (jasny) | `#3B6FE0` / `#F4F4F4` | **4,21 : 1** | 3,0 |
| fokus sygnał / tło (ciemny) | `#5C8CEC` / `#0F0F0F` | **5,86 : 1** | 3,0 |

**Reguła `tekst-3`:** żeton metadanych `#7C7C7C` jest zarezerwowany dla metadanych,
znaczników czasu i etykiet pomocniczych — treść czytelna używa `--dn-tekst-2`.

### 11.5. Zakazy barwne

| Zakaz | Powód |
|---|---|
| `#000000` jako tło i jako tekst | zabija głębię, powoduje halację na OLED; skala kończy się na `#0A0A0A` |
| `#FFFFFF` jako tło całej strony | tło strony to `#F4F4F4`; biel jest **powierzchnią** karty |
| Gradient jako tło przycisku, karty, sekcji | gradienty są wyłącznie ilustracyjne (awatar bez zdjęcia, rdzeń AOD, grafiki brandowe) |
| Wartość szesnastkowa wpisana wprost w CSS | wyłącznie `var(--dn-*)`; wyjątek: **pliki SVG znaku** |
| Komponent sięgający po prymitywy (`--dn-szary-500`) | komponent sięga wyłącznie po żetony **semantyczne** |
| Stan wyrażony samym kolorem | zawsze ikona albo etykieta |

---

## 12. Terytoria marki

### 12.1. Cztery terytoria

| Terytorium | Co obejmuje | Nośnik znaku | Reguła szczególna |
|---|---|---|---|
| **A · Aplikacja** | pasek górny, karty środowisk, AOD, okno startowe, Mobile | sygnet + logotyp tekstem na ramie; emblematy w kartach i bocznej nawigacji | godło zachowuje barwy własne; rama atramentowa w obu motywach |
| **B · Dokumentacja** | opracowania MD i HTML, prototypy okien, indeks | nagłówek metrykowy + stopka; favicon w `<head>` | zero wartości hex wprost; wyłącznie `var(--dn-*)` |
| **C · Materiały** | grafika OG, banner, tapeta, tło slajdu, papier firmowy | `logo-poziomy` / `logo-pionowy` | gradienty brandowe dopuszczone **tylko tutaj** |
| **D · Ikony systemowe** | favicon, ICO, apple-touch-icon, kafel PWA, manifest | wariant uproszczony (favicon) / pełny (kafel) | tożsamości systemowe pochodzą ze znaku, nie są rysowane osobno |

### 12.2. Mapa wystąpień w aplikacji

```
   ┌──────────────────────────────────────────────────────────────────────┐
   │ ▸▸· DANACO CONSOLE      TalkIn        [ szukaj ]        ● DN         │ ← A: rama
   ├──────────┬───────────────────────────────────────────────────────────┤
   │ ● TalkIn │  ┌── karta sesji ──┬── karta sesji ●──┬── + ──┐           │ ← A: kropka (3)
   │   Studio │  │                                                        │
   │   Resea… │  │   OKNO OPERACYJNE MODUŁU                               │
   │   Libra… │  │                                                        │
   │   Trans… │  ├────────────────────────────────────────────────────────┤
   │   Brows… │  │   CHAT WINDOW      ● nadawca pisze                     │ ← A: kropka (4)
   │   Assis… │  │   ▸ prompt…                                            │
   └──────────┴───────────────────────────────────────────────────────────┘
          ▲                                                    ▲
     emblemat środowiska (poziom 1)                    AOD — rdzeń awatara
     jedna wypełniona kropka                           A: kropka (5)
```

### 12.3. Reguły terytorialne

| # | Reguła |
|---|---|
| T1 | **Znak nie wchodzi w treść.** Sygnet występuje w ramie, w nagłówku dokumentu i w materiałach — nigdy jako element treści okna operacyjnego. |
| T2 | **Jedno wystąpienie znaku na widok.** Sygnet i logotyp poziomy nie występują jednocześnie w jednym układzie. |
| T3 | **Emblemat mieszka w swoim środowisku.** Emblemat TalkIn nie pojawia się w kontekście CodeStudio inaczej niż jako odnośnik do TalkIn. |
| T4 | **Gradient tylko na terytorium C** i w dwóch wyjątkach terytorium A: awatar bez zdjęcia oraz rdzeń AOD. |
| T5 | **Dokumentacja dziedziczy żetony bez modyfikacji.** Opracowanie nie definiuje własnych barw — wpina `zetony.css` i używa `var(--dn-*)`. |
| T6 | **Favicon jest zawsze wariantem uproszczonym.** Dwa groty w 16 px zlewają się w plamę. |

### 12.4. Komplet ikon dokumentu

Wpięcie w `<head>` (plik źródłowy: `zasoby/marka/favicon/naglowek-snippet.html`):

```html
<link rel="icon" href="/favicon.svg" type="image/svg+xml">
<link rel="icon" href="/favicon.ico" sizes="16x16 32x32 48x48">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">
<link rel="manifest" href="/site.webmanifest">
<meta name="theme-color" content="#F4F4F4" media="(prefers-color-scheme: light)">
<meta name="theme-color" content="#131313" media="(prefers-color-scheme: dark)">
```

Manifest deklaruje: `name` **Danaco Console** · `short_name` **Danaco** ·
`display` **standalone** · `background_color` `#0F0F0F` · `theme_color` `#131313` ·
trzy ikony, w tym jedna `maskable`.

---

## 13. Spójność — lista kontrolna zgodności z marką

Opracowanie jest **zgodne z marką Danaco Console**, jeżeli przechodzi wszystkie
pozycje poniżej. Lista jest podzielona na sześć bloków i jest przeznaczona do
odhaczania przed oddaniem pliku.

### 13.1. Blok I — znak

- [ ] Sygnet użyty w wariancie właściwym dla nośnika (rozdz. 8.4)
- [ ] Krzywe sygnetu **niezmienione co do współrzędnej**
- [ ] Kropka sygnału obecna i **w błękicie** (nie w atramencie — poza wariantami mono)
- [ ] Poniżej 24 px zastosowany **wariant uproszczony**
- [ ] Godło zachowuje **barwy własne** w obu motywach
- [ ] Znak nie jest rozciągnięty, obrócony, przebarwiony ani złożony z fontu
- [ ] Jedno wystąpienie znaku na widok

### 13.2. Blok II — barwa

- [ ] Zero wartości szesnastkowych wprost w CSS — wyłącznie `var(--dn-*)`
- [ ] Zero `#000000`
- [ ] Komponent sięga po żetony **semantyczne**, nie po prymitywy
- [ ] Sygnał zajmuje **≤ 5 %** powierzchni i nie jest tłem sekcji
- [ ] Przycisk główny jest **inwersją atramentu**, nie błękitem
- [ ] Pasek górny **atramentowy w obu motywach**
- [ ] Gradient nieużyty jako tło przycisku, karty ani sekcji
- [ ] Oba motywy sprawdzone wzrokowo, nie tylko zadeklarowane

### 13.3. Blok III — typografia i treść

- [ ] Nagłówki **Space Grotesk**, treść **IBM Plex Sans**, dane **IBM Plex Mono**
- [ ] Liczby w mono z `tabular-nums`
- [ ] Polskie diakrytyki poprawne w całym pliku (`ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ`)
- [ ] `lang="pl"` w dokumencie
- [ ] Nazwy własne modułów, okien i ról **dokładnie jak w dokumentacji**
- [ ] Zero Lorem ipsum, zmyślonych nazwisk, zmyślonych metryk, fikcyjnych firm
- [ ] Dane przykładowe pochodzą z domeny produktu i są oznaczone jako przykładowe

### 13.4. Blok IV — zero blokad

- [ ] Zero `disabled`
- [ ] Zero `pointer-events: none` na kontrolce
- [ ] Zero `cursor: not-allowed`
- [ ] Każda odmowa zawiera **następny krok**
- [ ] Liczniki czasu opisane jako informacyjne
- [ ] Metoda wyłączona w konfiguracji **nie jest renderowana**, nie jest blokowana

### 13.5. Blok V — ruch i dostępność

- [ ] Ruch ciągły **wyłącznie** jako tętno kropki (2,4 s)
- [ ] Mikroprzejścia w zakresie 100–220 ms, `--dn-ease`
- [ ] `prefers-reduced-motion` obsłużone **żetonami**, nie per komponent
- [ ] Stan nigdy samym kolorem — ikona albo etykieta towarzyszy
- [ ] Fokus widoczny na każdej kontrolce (2 px + odsunięcie 2 px)
- [ ] Nawigacja klawiaturą działa; role i `aria-*` obecne
- [ ] Cele dotykowe rosną żetonem przy `pointer: coarse`

### 13.6. Blok VI — struktura opracowania

- [ ] Nagłówek metrykowy obecny (produkt, wersja, status, data, odbiorcy, zakres)
- [ ] Stopka: `Danaco Console — AI Operating Environment · v2.0` + nota © 
- [ ] W HTML: przełącznik motywu `[data-przelacz-motyw]` **działa**
- [ ] W HTML: sekcja `<details>` „O tym opracowaniu" obecna
- [ ] W HTML: plik działa po otwarciu z `file://`, cały JS inline, zero CDN
- [ ] Ikony wyłącznie z `zasoby/ikony/svg/`, wklejane inline
- [ ] Zero emoji jako ikon
- [ ] Zero elementów bez pokrycia w dokumentacji

### 13.7. Wynik

| Bloki zaliczone | Werdykt |
|---|---|
| 6 / 6 | zgodne z marką — do oddania |
| 5 / 6 | do poprawy — wskaż blok i pozycję |
| ≤ 4 / 6 | niezgodne — opracowanie wymaga przeróbki, nie korekty |

---

## 14. Decyzje projektowe

Rozstrzygnięcia podjęte w tym opracowaniu tam, gdzie dokumentacja źródłowa
nie wypowiada się wprost. Każde jest odwracalne i podlega akceptacji Właściciela.

| # | Decyzja | Uzasadnienie | Alternatywa odrzucona |
|---|---|---|---|
| **D1** | Osobowość marki opisano **pięcioma cechami w parach „jesteśmy / nie jesteśmy"**, każda wyprowadzona z konkretnego rozdziału `KIERUNEK.md` i z dowodem w pliku | cecha bez dowodu w pliku jest nieweryfikowalna; para „jesteśmy / nie jesteśmy" czyni ją operacyjną | lista przymiotników bez odniesień |
| **D2** | Obietnica marki sformułowana jako **trzy człony** („wydajesz polecenie — praca rusza — nic cię nie zatrzymuje"), po jednym na każdy element znaku i zasadę ADL-017 | obietnica musi mieć nośnik; trzy człony pokrywają się z trzema znakami „»»." | jedno zdanie ogólne bez odwzorowania w znaku |
| **D3** | Obietnice czterech środowisk **wyprowadzone z ich mott i opisów w dokumentacji** (E1 + `przeplyw-okien.md`), nie ułożone na nowo | zakaz projektowania treści bez pokrycia; motta są zatwierdzone w makiecie E1 | wymyślenie haseł marketingowych dla środowisk |
| **D4** | ADL-017 opisane w **trzech warstwach**: interfejs, język, komunikacja zewnętrzna | zasada jest wartością marki, więc musi mieć konsekwencje poza CSS-em | pozostawienie ADL-017 wyłącznie jako reguły komponentów |
| **D5** | Architektura marki ujęta jako **cztery poziomy z regułami dziedziczenia R1–R6** | hierarchia bez reguł dziedziczenia nie rozstrzyga sporów granicznych (emblemat vs ikona, znak vs favicon) | prosta lista plików bez modelu |
| **D6** | Rozstrzygnięto, że **zestaw ikon interfejsu nie zawiera kropki sygnału** | kropka jest zarezerwowana dla marki (poziom 0–1) i dla stanu; ikona z kropką konkurowałaby ze wskaźnikiem pracy | dopuszczenie kropki w ikonach domenowych |
| **D7** | Odnotowano, że w plikach znaku na ciemnym tle groty mają wartość **`#F4F4F4`** (`szary-50`), a nie `#ECECEC` (`szary-100` = żeton tekstu w motywie ciemnym) | zgodnie z zasadą „znak ma barwy własne", biel znaku jest wartością marki, nie żetonem tekstu; wartość odczytana z plików źródłowych i pozostawiona bez zmian | ujednolicenie znaku do żetonu tekstu — naruszałoby zasadę barw własnych |
| **D8** | Kropka w wariancie uproszczonym opisana przez **stały stosunek 0,30** do wysokości znaku | stosunek, a nie wartość bezwzględna, tłumaczy dlaczego oba warianty czyta się jako ten sam znak | opis wyłącznie wartościami `r 6,5` / `r 9` |
| **D9** | Emblemat **CodeStudio** wskazano jako jedyny cytujący grot znaku nadrzędnego; **MultitaskingAI** — jako jedyny, w którym kropką jest węzeł centralny | obie cechy są faktem geometrycznym w plikach i niosą znaczenie; wymagają odnotowania, żeby nie zostały „ujednolicone" przy przyszłej edycji | pominięcie różnic jako przypadkowych |
| **D10** | Ikona aplikacji opisana jako **zapowiedź paska górnego** (tło `#131313` = `--dn-szary-925` = `--dn-rama`) | tożsamość barwy kafla i ramy jest faktem w plikach; nazwanie jej czyni ją regułą, nie zbiegiem okoliczności | opis kafla jako „ciemne tło" |
| **D11** | Lista kontrolna spójności podzielona na **sześć bloków z progiem 6/6** | próg binarny wymusza rozstrzygnięcie zamiast negocjacji; podział na bloki wskazuje, gdzie leży wada | jedna płaska lista bez progu |
| **D12** | Relację Console ↔ Pilot opisano tabelą **„co zostało zastąpione"** z wyszczególnieniem trzech elementów ciągłości (kokpit, dwa motywy, zero blokad) | wersja v2.0 zastępuje warstwę wizualną, ale nie funkcjonalną — bez tego rozróżnienia czytelnik uzna całą dokumentację v1.0 za nieaktualną | wzmianka jednozdaniowa |
| **D13** | W dokumencie **nie powtórzono** księgi znaku (pole ochronne, rozmiary minimalne, katalog zakazów) | te treści należą do opracowania `10` — powielenie groziłoby rozjechaniem się dwóch źródeł | włączenie skróconej księgi do tego pliku |

---

## Załącznik A — inwentarz plików marki

Ścieżka bazowa: `WYNIK/zasoby/marka/`

### A.1. Logo

| Plik | Wymiar | Zawartość |
|---|---|---|
| `logo/sygnet.svg` | 96 × 96 | dwa groty `#181818` + kropka `#3B6FE0` |
| `logo/sygnet-na-ciemnym.svg` | 96 × 96 | dwa groty `#F4F4F4` + kropka `#5C8CEC` |
| `logo/sygnet-mono-czarny.svg` | 96 × 96 | całość `#181818` |
| `logo/sygnet-mono-bialy.svg` | 96 × 96 | całość `#F4F4F4` |
| `logo/sygnet-uproszczony.svg` | 96 × 96 | jeden grot + kropka `r 9` |
| `logo/sygnet-uproszczony-na-ciemnym.svg` | 96 × 96 | jw., wersja kontrowa |
| `logo/logo-poziomy.svg` | 225 × 96 | sygnet ×0,6667 + logotyp dwuwierszowy |
| `logo/logo-poziomy-na-ciemnym.svg` | 225 × 96 | jw., wersja kontrowa |
| `logo/logo-poziomy-mono-czarny.svg` | 225 × 96 | jednobarwny `#181818` |
| `logo/logo-poziomy-mono-bialy.svg` | 225 × 96 | jednobarwny `#F4F4F4` |
| `logo/logo-pionowy.svg` | 159 × 172 | sygnet nad logotypem |
| `logo/logo-pionowy-na-ciemnym.svg` | 159 × 172 | jw., wersja kontrowa |
| `logo/logotyp.svg` | 149 × 56 | sam logotyp + kropka `r 3,4` |
| `logo/logotyp-na-ciemnym.svg` | 149 × 56 | jw., wersja kontrowa |
| `logo/png/*` | — | `sygnet-512`, `logo-poziomy@2x`, `logo-pionowy@2x`, `logotyp@2x` i wersje na ciemnym |

### A.2. Emblematy środowisk

| Plik | Kropka sygnału | Elementy |
|---|---|---|
| `srodowiska/srodowisko-talkin.svg` | `cx 16,2 · cy 10,8 · r 1,5` | dymek rozmowy, dwie linie |
| `srodowiska/srodowisko-workspace.svg` | `cx 16,9 · cy 16,9 · r 1,6` | cztery pola 7,2 × 7,2 (r 2) |
| `srodowiska/srodowisko-codestudio.svg` | `cx 18,4 · cy 9,9 · r 1,5` | okno 18 × 16 (r 2,5), grot, linia |
| `srodowiska/srodowisko-multitaskingai.svg` | `cx 12 · cy 12 · r 2,6` (węzeł centralny) | trzy węzły `r 1,9`, trzy połączenia |

### A.3. Tożsamości systemowe

| Plik | Przeznaczenie |
|---|---|
| `favicon/favicon.svg` | wariant uproszczony + `@media (prefers-color-scheme: dark)` |
| `favicon/favicon.ico` | 16 / 32 / 48 px |
| `favicon/favicon-16.png` · `-32` · `-48` | źródła ICO |
| `favicon/apple-touch-icon.png` | 180 px |
| `favicon/icon-192.png` · `icon-512.png` | PWA |
| `favicon/site.webmanifest` | manifest aplikacji |
| `favicon/naglowek-snippet.html` | gotowe wpięcie do `<head>` |
| `ikona-aplikacji/ikona-aplikacji.svg` | kafel 1024, r 224, tło `#131313` |
| `ikona-aplikacji/ikona-maskowalna.svg` | r 0, strefa bezpieczna |
| `ikona-aplikacji/ikona-180/192/512/1024.png` | rastry kafla |

### A.4. Koncepcje i akceptacja

| Plik | Rola |
|---|---|
| `logo/alternatywy/plansza-koncepcji.html` | pięć koncepcji: A1 Kwadranty · A2 D-obwód · A3 Monogram DC · A4 Warstwy · **A5 Delegacja** |
| `logo/alternatywy/a5-final-akceptacja.html` | wersja ostateczna do akceptacji — testy skali 48/24/16 px, lockupy, kontekst ramy, favicon |
| `logo/alternatywy/a1-kwadranty*.svg` · `a4-warstwy*.svg` | koncepcje odrzucone, zachowane jako materiał dowodowy procesu |
| `logo/alternatywy/a5-final*.svg` | pliki zatwierdzonej koncepcji, tożsame z plikami produkcyjnymi |

**Zatwierdzono koncepcję A5 „Delegacja"** — z uzasadnieniem: „znak opowiada pętlę,
nie literę". Poprawka wprowadzona przed akceptacją: **kropka przeniesiona na dół,
na linię bazową grotów**, dzięki czemu znak czyta się „»».".

---

## Załącznik B — słowniczek marki

| Pojęcie | Definicja |
|---|---|
| **Operator** | Adresat marki: zawodowiec zarządzający cyfrową organizacją złożoną z modeli, agentów, projektów i automatyzacji. Nie „użytkownik czatu". |
| **Kokpit** | Metafora układu: rama atramentowa wokół przełączalnej treści, dająca efekt stanowiska dowodzenia. |
| **Instrument pomiarowy** | Idea przewodnia: gęsta skala, maszynowy odczyt, zero elementów bez funkcji. |
| **Sygnet „Delegacja"** | Znak nadrzędny: podwójny grot promptu + kropka sygnału, czytany „»».". |
| **Kropka sygnału** | Element sygnaturowy: wypełniony punkt błękitu oznaczający „tu biegnie praca". |
| **Tętno** | Animacja kropki, 2,4 s — jedyny ruch ciągły w interfejsie. |
| **Rama** | Pasek górny, atramentowy w obu motywach; jedyny element nieprzełączający się z motywem (obok godła). |
| **Barwy własne** | Barwy znaku niezależne od motywu; przeciwieństwo `currentColor`. |
| **Inwersja atramentu** | Przycisk główny: czarny na jasnym, biały na ciemnym — najsilniejszy akcent bez drugiej barwy. |
| **Sygnał** | Jedyna barwa akcentu (`#3B6FE0` i rodzina); wskazuje, nie działa. |
| **Zero blokad (ADL-017)** | Wartość marki: domyślne zachowanie systemu to wykonanie polecenia; izolacja i uprawnienia są opcjami. |
| **Emblemat środowiska** | Znak poziomu 1: siatka 24, obrys 1,75, `currentColor`, dokładnie jedna wypełniona kropka. |
| **Wariant uproszczony** | Jeden grot + większa kropka; obowiązuje poniżej 24 px. |
| **Żeton semantyczny** | Zmienna `--dn-*` o nazwie opisującej rolę; jedyny dopuszczalny sposób odwołania się do barwy w komponencie. |
| **Terytorium marki** | Jeden z czterech obszarów wystąpienia: aplikacja · dokumentacja · materiały · ikony systemowe. |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
