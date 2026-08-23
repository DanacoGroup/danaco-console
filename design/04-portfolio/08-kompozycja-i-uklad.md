# Danaco Console — Kompozycja i układ (plansza portfolio P8)

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Producent** | Danaco Holding Group Sp. z o.o. · Twórca: Dariusz Naharnowicz |
| **Wersja opracowania** | v2.0 · Status: Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | projektant wizualny, projektant interakcji, zespół wdrożeniowy front-end, recenzent portfolio |
| **Zakres** | jednostka rytmu, gęstość i wymiary, anatomia powłoki i obszaru roboczego, katalog wzorców układu, siatka 12 kolumn, punkty łamania, hierarchia wizualna, warstwy z-index |
| **Plik towarzyszący** | `WYNIK/04-portfolio/08-kompozycja-i-uklad.html` |

---

## 0. Czym jest ta plansza

Plansza P8 opowiada **budowę ekranu** Danaco Console. Nie jest dokumentacją techniczną —
ta powstała w `01-dokumentacja-md/` i `02-dokumentacja-html/`. Jest opracowaniem
prezentacyjnym: pokazuje, jak z jednej liczby (**4 px**) wyrasta cała geometria kokpitu,
i dlaczego każdy wymiar jest taki, jaki jest.

Wszystkie liczby na planszy pochodzą z KANON rozdz. 2 (przestrzeń, wymiary, punkty łamania,
warstwy) albo z bezpośredniego odczytu plików `WYNIK/zasoby/zetony/zetony.css`
i `WYNIK/zasoby/prototyp.css`. Wszystkie nazwy okien, modułów, paneli i stref pochodzą
z KANON rozdz. 6–7 oraz z `dok/projekt-ui/przeplyw/elementy-okien.md`.
Żaden układ pokazany na planszy nie jest ilustracyjnym wymysłem — każdy ma wskazane
miejsce występowania w platformie i odpowiadający mu prototyp w `WYNIK/05-okna/`.

---

## 1. Jednostka 4 px — fundament nienegocjowalny

> **4 px.** Wszystkie odstępy są jej wielokrotnością — brak wartości spoza skali.
> `SIATKA-I-UKLAD.md` rozdz. 1 · [NIENEGOCJOWALNE]

### 1.1. Skala odstępów

| Żeton | `od-0` | `od-1` | `od-2` | `od-3` | `od-4` | `od-5` | `od-6` | `od-8` | `od-10` | `od-12` | `od-16` |
|---|---|---|---|---|---|---|---|---|---|---|---|
| px | 0 | 4 | 8 | 12 | 16 | 20 | 24 | 32 | 40 | 48 | 64 |
| krotność 4 | 0× | 1× | 2× | 3× | 4× | 5× | 6× | 8× | 10× | 12× | 16× |

Skala jest **nierównomierna z zamysłem**: gęsto tam, gdzie decyzje zapadają często
(4–24 px, wnętrze paneli i list), rzadko tam, gdzie odstęp rozdziela całe bloki
(32–64 px, przerwy między sekcjami planszy dokumentowej).

### 1.2. Dlaczego akurat 4

| Powód | Konsekwencja w kokpicie |
|---|---|
| Każdy wymiar kontrolki dzieli się przez 4 | 32 = 8×4 · 36 = 9×4 · 48 = 12×4 · 224 = 56×4 · 320 = 80×4 |
| Wiersze list układają się w regularny rytm | wiersz 36 px = pas kart 36 px — kontynuacja rytmu przez pasy powłoki |
| Ikona wpisuje się w kontrolkę bez ułamków | ikona 16 w kontrolce 32 → margines 8 z każdej strony |
| Zero wartości „na oko” | brak `padding: 7px`, brak `margin: 13px` — recenzja jest arytmetyczna, nie estetyczna |

### 1.3. Co się psuje przy odstępstwie

Plansza pokazuje parę porównawczą na wierszu listy modułów bocznej nawigacji
(realny element powłoki TalkIn):

| Cecha | Zgodne z rytmem | Odstępstwo (wartości spoza skali) |
|---|---|---|
| Wysokość wiersza | 36 px | 35 px |
| Odstęp wewnętrzny | `od-3` = 12 px | 11 px |
| Odstęp ikona–tekst | `od-2` = 8 px | 7 px |
| Skutek na 9 pozycjach TalkIn | pas kończy się równo z sąsiadem | narastający dryf 9 px — dolna krawędź listy nie trafia w żadną linię |
| Skutek przy przełączeniu gęstości | wszystko skaluje się żetonem | trzeba poprawiać ręcznie każdy przypadek |

Dryf jest niewidoczny na jednym elemencie i nieodwracalny na dziewięciu — dlatego
reguła jest bezwzględna, nie zalecana.

---

## 2. Gęstość zwarta (8/10) — pełna tabela wymiarów

`GESTOSC_WIZUALNA = 8/10` (KANON rozdz. 1). Kokpit dowodzenia pokazuje dużo stanu naraz;
gęstość jest narzędziem czytelności, nie oszczędności.

### 2.1. Wymiary wiążące

| Element | Żeton | Zwarta (domyślna) | Przestronna | Dotyk (`pointer: coarse`) |
|---|---|---|---|---|
| Kontrolka (przycisk, pole) | `--dn-wym-kontrolka` | **32 px** | 40 px | 40 px |
| Przycisk ikonowy | `--dn-wym-ikonowy` | **32 px** | 40 px | 40 px |
| Pasek górny | `--dn-wym-pasek` | **48 px** | 56 px | 48 px |
| Pas kart sesji | `--dn-wym-pas-kart` | **36 px** | 40 px | 40 px |
| Wiersz listy / tabeli | `--dn-wym-wiersz` | **36 px** | 44 px | 44 px |
| Boczna nawigacja | `--dn-wym-boczna` | **224 px** | 224 px | 224 px |
| Pas komunikacji (min.) | `--dn-wym-pas-komunikacji` | **320 px** | 320 px | 320 px |
| Modal (maks. szerokość) | `--dn-wym-modal` | **560 px** | 560 px | 560 px |
| Awatar | `--dn-wym-awatar-sm/-/-lg` | 24 / 28 / 36 px | — | — |
| Ikona | `--dn-wym-ikona-sm/-/-lg/-xl` | 14 / 16 / 20 / 24 px | — | — |
| Przełącznik | `--dn-wym-przelacznik-*` | 36 × 20 px | 44 × 24 px | 44 × 24 px |
| Pole wyboru | `--dn-wym-check` | 16 px | 20 px | 20 px |
| Kropka sygnału | `--dn-wym-kropka` | **6 px** | — | — |
| Wstęga wyboru | `--dn-wym-wstega` | 2 px | — | — |
| Wskaźnik ładowania | `--dn-wym-spinner` | 14 px | — | — |
| Pierścień fokusu | `--dn-wym-fokus` + odsunięcie | 2 px + 2 px | — | — |

Wszystkie odczytane bezpośrednio z `WYNIK/zasoby/zetony/zetony.css`
(bloki `:root`, `:root[data-gestosc='przestronna']`, `@media (pointer: coarse)`).

### 2.2. Zasada przełączania gęstości

Gęstość zmienia się **jednym atrybutem** na elemencie głównym (`data-gestosc`),
nie interwencją na poziomie komponentu. Żaden komponent nie zna słowa „przestronna” —
sięga wyłącznie po żeton. Ta sama zasada obowiązuje przy `pointer: coarse`:
**cele dotykowe rosną żetonem, nie wyjątkiem** (KANON rozdz. 9).

### 2.3. Rozbieżność wobec pakietu v1.0 — rozstrzygnięcie

`SIATKA-I-UKLAD.md` (v1.0, 2026-08-11) podaje: przycisk ~36 px, przycisk ikonowy 36×36,
przełącznik 40×22, pasek górny 56 px. KANON v2.0 podaje odpowiednio 32 / 32 / 36×20 / 48.
**Wiążący jest KANON v2.0** — wartości v1.0 pozostają jako zapis wcześniejszego etapu
i odpowiadają dzisiejszemu wariantowi „przestronna” (kontrolka 40, pasek 56).
Rozbieżność odnotowana jawnie, zgodnie z KANON rozdz. 10 pkt 5.

---

## 3. Anatomia powłoki środowiska

Powłoka środowisk modułowych (TalkIn · WorkSpace · CodeStudio) ma **cztery pasy**
(`elementy-okien.md` rozdz. 3.1, Schemat 4). Realizacja w `prototyp.css` klasą `.pt-powloka`:

```
grid-template-rows: 48px 36px minmax(0, 1fr);   /* pasek · karty sesji · ciało */
grid-template-columns: 224px minmax(0, 1fr);    /* boczna · obszar roboczy (.pt-cialo) */
```

| Strefa | Wymiar | Rola | Zawartość (dokumentacja) |
|---|---|---|---|
| Pasek górny | wys. 48 px, pełna szerokość | tożsamość i funkcje globalne | blok marki, ikona „dom”, pole wyszukiwania (do 520 px), szybka konfiguracja sesji, Mobile, Always On Display, przełącznik motywu |
| Pas kart sesji | wys. 36 px, pełna szerokość | równoległość pracy | karty sesji jak karty przeglądarki + przycisk „+” |
| Boczna nawigacja | szer. 224 px | wybór modułu | TalkIn 9 pozycji · WorkSpace 9 · CodeStudio 8; MultitaskingAI: panel orkiestracji, 6 sekcji |
| Obszar roboczy | reszta | praca | okna operacyjne modułu + Chat Window (pas komunikacji) |

**Pasek górny jest atramentowy w obu motywach** — jedyny element nieprzełączający się
z motywem (KANON rozdz. 2, „Rama kokpitu”).

### 3.1. Warianty powłoki w prototypie

| Klasa | Wiersze | Gdzie występuje |
|---|---|---|
| `.pt-powloka` | 48 · 36 · 1fr | powłoki czterech środowisk |
| `.pt-powloka--bez-kart` | 48 · 1fr | okna platformowe bez kart sesji (Konfiguracja, Ustawienia) |
| `.pt-powloka--jednolita` | 1fr | okno startowe, rejestracja i logowanie, Mobile |
| `.pt-cialo--zwinieta` | kolumny 32 · 1fr | boczna zwinięta do ikon (poniżej w2) |
| `.pt-cialo--bez-bocznej` | kolumny 1fr | Centrum dowodzenia, Always On Display |

---

## 4. Obszar roboczy z pasem komunikacji

Chat Window jest **oknem wspólnym wszystkim modułom** (KANON rozdz. 7.4). W obszarze
roboczym zajmuje dolny pas:

```
.pt-obszar--z-komunikacja {
  grid-template-rows: minmax(0, 1fr) minmax(320px, 46%);
}
```

| Człon | Znaczenie |
|---|---|
| `minmax(0, 1fr)` — góra | okna operacyjne modułu; `0` jako minimum pozwala treści kurczyć się bez rozpychania siatki |
| `minmax(320px, 46%)` — dół | pas komunikacji: nigdy węższy niż 320 px, nigdy większy niż 46 % wysokości obszaru |
| `.pt-obszar--komunikacja-zwinieta` | pas zwinięty do samego pola poleceń (`auto`) — praca w oknie operacyjnym na pełnej wysokości |

**Dlaczego 46 %, a nie 50 %:** równy podział czyta się jak dwa równoważne okna.
Kokpit ma jeden środek ciężkości — okno operacyjne modułu. Przewaga 54 : 46 jest
minimalna, ale jednoznaczna, i mieści się w `WARIANCJA_PROJEKTOWA 4/10`
(asymetria wyłącznie tam, gdzie niesie hierarchię).

**Dlaczego 320 px twardego minimum:** poniżej tej wartości wpis rozmowy
(`.dn-wpis` — medalion 28 px + nadawca + treść) łamie się na dwa wiersze przy
pierwszym dłuższym zdaniu, a historia przestaje być czytelna jako rozmowa.

---

## 5. Katalog wzorców układu

Osiem wzorców — każdy z realnym miejscem występowania w platformie i klasą prototypu.

| # | Wzorzec | Geometria | Klasa | Gdzie występuje (dokumentacja) |
|---|---|---|---|---|
| 1 | Trzy strefy o malejącej masie | pion: siatka kart → siatka kafli → listwa | — (układ własny okna) | **Centrum dowodzenia**: Strefa 1 (4 karty środowisk) · Strefa 2 (4 kafle) · Strefa 3 (listwa 3 pozycji) |
| 2 | Boczny lewy | `minmax(180px, 260px) / minmax(0, 1fr)` | `.pt-panele--boczny-lewy` | **Developer** (Project Tree + Code Editor), **Library** (Library Explorer + File Preview), **Research** (Sources Manager) |
| 3 | Boczny prawy | `minmax(0, 1fr) / minmax(200px, 300px)` | `.pt-panele--boczny-prawy` | **Research** (Findings Panel), **Konfiguracja** (Panel prowenancji), **Browser** (Notes Panel) |
| 4 | Trójpodział | `minmax(160px,220px) / minmax(0,1fr) / minmax(200px,300px)` | `.pt-panele--trojpodzial` | **Developer** (Project Tree + Code Editor + Git Panel), **Apps**, **Design**, **Workspace**, **Automations** |
| 5 | Panele równoległe 2 / 3 / 4 | `repeat(n, minmax(0, 1fr))` | `.pt-panele--2/--3/--4` | **Roundtable** (Model Panels — do 4), **Translate** (Source Panel + Translation Panels — 2/3/4), **Studio** (Diff/Grep Panel) |
| 6 | Asymetria koordynator–wykonawca | siatka kart ról o nierównej masie | `.dn-karta` w obszarze roboczym | **MultitaskingAI**, sekcja „Role”: Executor 1 · Executor 2 · **Coordinator** · Executor 3 / Validator |
| 7 | Nakładka modalna | nakładka 800 + modal 900, maks. 560 px | `.dn-modal` | **Konfiguracja** (modal „Pula kont Code CLI”), **Ustawienia**, okno punktów izolacji |
| 8 | Warstwa centralna | warstwa 1200 nad całą powłoką | `.dn-aod` | **Always On Display** — rdzeń awatara, monitor procesu |

### 5.1. Wzorzec 1 — trzy strefy o malejącej masie

Anty-domyślne KANON rozdz. 1 blokuje wprost „trzy równe karty funkcji”. Centrum
dowodzenia rozwiązuje to malejącą masą:

| Strefa | Element | Liczba | Masa wizualna |
|---|---|---|---|
| 1 | karta środowiska (`.dn-karta-srodowiska`) — emblemat, tytuł, motto, opis | 4 | największa — środek ciężkości strony |
| 2 | kafel komponentu własnego (`.dn-kafel`) — ikona, etykieta, opis | 4 | pośrednia |
| 3 | listwa ustawień (`.dn-listwa`) — pas jednowierszowy | 3 pozycje | najmniejsza |

Kolejność i liczebność są źródłowe (`elementy-okien.md` rozdz. 2.2–2.5,
`SIATKA-I-UKLAD.md` rozdz. 8). Wymiary stref są decyzją projektową — dokumentacja
świadomie ich nie przesądza.

### 5.2. Wzorzec 6 — asymetria koordynator–wykonawca

Cztery role MultitaskingAI nie są czterema równymi kaflami. **Coordinator** rozdaje
pracę, **Executory** ją wykonują, **Executor 3 / Validator** ją sprawdza. Układ to pokazuje:
karta koordynatora ma większą masę (szerszy kafel, kropka tętna przy pracy w tle),
karty wykonawców tworzą pod nią rytm równych jednostek, karta walidatora zamyka pętlę.
Piąta pozycja zespołu — **Subagent Network** — jest rozwinięciem zagnieżdżonym pod kartą
wykonawcy, nie osobną sekcją panelu (`elementy-okien.md` Tabela 17).

---

## 6. Siatka 12 kolumn — kiedy obowiązuje

| Parametr | Wartość | Żeton |
|---|---|---|
| Kolumny | 12 | `--dn-siatka-kolumny` |
| Przerwa | 24 px | `--dn-siatka-przerwa` (= `od-6`) |
| Maks. szerokość treści | 1200 px | `--dn-tresc-max` |

| Rodzaj powierzchni | Czy obowiązuje siatka 12 | Przykład |
|---|---|---|
| Treść dokumentowa | **tak** | Konfiguracja (13 zakresów + obszar 14.), Ustawienia, plansze portfolio, panel „O tym opracowaniu” |
| Strona główna — Centrum dowodzenia | tak, w obrębie stref | Strefa 1: 4 karty (3 kolumny każda) · Strefa 2: 4 kafle |
| Kokpit — powłoka środowiska | **nie** | 224 px bocznej i 320 px pasa komunikacji to **wymiary bezwzględne**, nie ułamki 12 kolumn |
| Okna operacyjne modułów | nie | proporcje paneli opisane w `minmax()`, nie w kolumnach |

**Uzasadnienie rozdziału:** siatka proporcjonalna służy czytaniu — treść ma być tej samej
szerokości niezależnie od monitora. Kokpit służy sterowaniu — boczna nawigacja ma mieć
tę samą szerokość, żeby ręka trafiała w tę samą pozycję na każdym ekranie. Rozciąganie
bocznej do 1/12 szerokości 2560-pikselowego monitora dałoby pas 213 px na jednym sprzęcie
i 107 px na drugim — czyli różne narzędzie na każdym biurku.

---

## 7. Cztery punkty łamania

| Próg | Szerokość | Nazwa | Co się dzieje |
|---|---|---|---|
| `w1` | 640 px | telefon poziomo / Mobile | boczna nawigacja znika całkowicie; panele boczne i równoległe składają się do jednej kolumny (widoczny wyłącznie pierwszy panel); `.pt-dokument` traci szerokie marginesy (`od-8` → `od-4`) |
| `w2` | 960 px | tablet | boczna **zwija się do ikon** (224 → 32 px); etykiety pozycji i nagłówek nawigacji znikają; trójpodział traci lewą kolumnę i staje się układem dwupanelowym |
| `w3` | 1280 px | biurko | pełny kokpit — wszystkie pasy i panele w wymiarach bazowych |
| `w4` | 1600 px | szerokie biurko | **dwa okna komunikacji** równolegle |

Reguły odczytane z `prototyp.css` rozdz. 7 (`@media (max-width: 960px)`,
`@media (max-width: 640px)`); progi z KANON rozdz. 2.

**Zasada porządku zwijania:** znika najpierw to, co najtańsze do odtworzenia jednym
kliknięciem. Kolejność: panele pomocnicze → etykiety bocznej nawigacji → sama boczna
nawigacja. Nigdy nie znika pasek górny (tożsamość + funkcje globalne) ani pas
komunikacji (Chat Window jest oknem wspólnym wszystkim modułom).

---

## 8. Hierarchia wizualna w kokpicie — cztery narzędzia

Monochrom odbiera projektantowi barwę jako narzędzie hierarchii. Zostają cztery.

| Narzędzie | Czym się posługuje | Słaba hierarchia | Mocna hierarchia |
|---|---|---|---|
| **Masa** | rozmiar elementu i jego wagi typograficznej | cztery karty środowisk równe kaflom komponentów własnych | karta środowiska większa od kafla; kafel większy od pozycji listwy |
| **Kontrast** | stopień jasności tekstu wobec tła (`--dn-tekst` / `-2` / `-3`) | wszystko w `--dn-tekst` — oko nie wie, co czytać najpierw | tytuł `--dn-tekst`, opis `--dn-tekst-2`, metadana `--dn-tekst-3` |
| **Pozycja** | miejsce w porządku czytania | akcja główna wmieszana między akcje poboczne | akcja główna zawsze na końcu stopki modala, po prawej |
| **Przestrzeń** | odstęp jako granica grupy | jednakowy odstęp między wszystkim — jedna kasza | wewnątrz grupy `od-2`, między grupami `od-6` |

Para porównawcza na planszy zbudowana jest na realnym układzie: **Strefa 1 Centrum
dowodzenia** (masa), **wiersz Monitora procesu MultitaskingAI** (kontrast),
**stopka modala „Pula kont Code CLI”** (pozycja), **panel Role** (przestrzeń).

Ograniczenie bezwzględne (KANON rozdz. 9): **stan nigdy samym kolorem** — każde
narzędzie hierarchii działa równolegle z ikoną albo etykietą.

---

## 9. Jedenaście warstw z-index

| Poziom | Wartość | Żeton | Przykład elementu |
|---|---|---|---|
| Podłoga | 0 | `--dn-z-podloga` | obszar roboczy, okna operacyjne modułu |
| Przybornik | 10 | `--dn-z-przybornik` | `.dn-przybornik`, lokalne paski narzędzi paneli |
| Pasek | 100 | `--dn-z-pasek` | pasek górny + pas kart sesji |
| Boczna | 200 | `--dn-z-boczna` | boczna nawigacja modułów, panel orkiestracji |
| Pas komunikacji | 300 | `--dn-z-pas-komunikacji` | Chat Window |
| Nakładka | 800 | `--dn-z-nakladka` | przyciemnienie tła modala (jedyne miejsce rozmycia) |
| Modal | 900 | `--dn-z-modal` | „Pula kont Code CLI”, okno punktów izolacji |
| Powiadomienie | 1000 | `--dn-z-powiadomienie` | `.dn-toast` |
| Tooltip | 1100 | `--dn-z-tooltip` | `.dn-tooltip` |
| Always On Display | 1200 | `--dn-z-aod` | rdzeń awatara, monitor procesu |
| Centrum poleceń | 1300 | `--dn-z-centrum-polecen` | wywołanie polecenia ponad wszystkim |

**Dlaczego AOD nad powiadomieniem:** Always On Display jest agentem towarzyszącym
obecnym ponad bieżącą przestrzenią roboczą. Toast, który przykryłby jego rdzeń,
przerwałby jedyny stały wskaźnik pracy w tle.

**Dlaczego skoki, nie kolejne liczby:** przerwa 10 → 100 → 200 zostawia miejsce na
warstwy pośrednie bez przenumerowania stosu. Skok 300 → 800 oddziela **warstwy powłoki**
(zawsze obecne) od **warstw przywołanych** (pojawiają się na żądanie).

---

## 10. Decyzje projektowe tej planszy

1. **Linijka rytmu jako nakładka, nie obrazek.** Rytm 4 px pokazany jest nakładką
   `repeating-linear-gradient` włączaną przyciskiem na żywym układzie, nie zrzutem
   ekranu z narysowanymi liniami — czytelnik widzi, że rytm jest własnością układu.
2. **Każdy wymiar w skali 1:1.** Tabela wymiarów jest równocześnie makietą: obok liczby
   stoi element o tej dokładnie wysokości z linią wymiarową. Liczba bez okazu nie uczy nic.
3. **Przełącznik gęstości przebudowuje podgląd, nie tylko tabelę.** Trzy warianty
   (zwarta / przestronna / dotyk) zmieniają atrybut `data-gestosc` i `--dn-wym-*`
   w obrębie podglądu — dokładnie tak, jak robi to platforma.
4. **Symulator punktów łamania zamiast zrzutów.** Suwak szerokości steruje szerokością
   kontenera szkieletu; szkielet reaguje zapytaniami kontenerowymi odwzorowującymi
   progi w1–w4. Dzięki temu zachowanie jest **sprawdzalne**, a nie deklarowane.
5. **Rozbieżność v1.0 / v2.0 pokazana, nie zamieciona.** Sekcja wymiarów odnotowuje
   różnicę wobec `SIATKA-I-UKLAD.md` i wskazuje, że dawne wartości odpowiadają
   dzisiejszemu wariantowi „przestronna”.
6. **Stos z-index w perspektywie 3D.** Jedenaście płaszczyzn przechylonych transformacją
   CSS — jedyne miejsce planszy, gdzie forma jest przestrzenna, bo pojęcie jest
   przestrzenne. Ruch pozostaje w budżecie 3/10: przechył jest statyczny, animowane
   jest wyłącznie wyróżnienie wybranej warstwy (160 ms).
7. **Pary porównawcze „słaba vs mocna” zbudowane z realnych układów.** Nie ma szarych
   prostokątów zastępczych — porównanie działa tylko wtedy, gdy czytelnik rozpoznaje
   w nim ekran, który zna.
8. **Zero wartości szesnastkowych, zero `disabled`.** Cała warstwa lokalna stylu
   korzysta wyłącznie z `var(--dn-*)`; żadna kontrolka planszy nie jest wyłączana —
   przycisk już wybranego wariantu pozostaje klikalny i potwierdza wybór komunikatem.

---

## 11. Źródła

| Zakres | Plik / rozdział |
|---|---|
| Przestrzeń, wymiary, punkty łamania, warstwy, siatka | `WYNIK/KANON.md` rozdz. 2 |
| Inwentarz okien i powłok | `WYNIK/KANON.md` rozdz. 7 |
| Architektura pojęciowa, moduły, środowiska | `WYNIK/KANON.md` rozdz. 6 |
| Zero blokad (ADL-017) | `WYNIK/KANON.md` rozdz. 8 |
| Anatomia powłoki, pasek górny, karty sesji, boczna nawigacja, obszar roboczy | `dok/projekt-ui/przeplyw/elementy-okien.md` rozdz. 3.1–3.5 |
| Centrum dowodzenia — trzy strefy | `dok/projekt-ui/przeplyw/elementy-okien.md` rozdz. 2.2–2.5 |
| Panel orkiestracji (6 sekcji), karty ról, Subagent Network | `dok/projekt-ui/przeplyw/elementy-okien.md` rozdz. 4.4–4.5 |
| Jednostka 4 px, skala odstępów, promienie, rytm strony głównej | `design/opracowania/system-wizualny/pakiet/SIATKA-I-UKLAD.md` rozdz. 1–8 |
| Klasy układu powłoki i paneli, responsywność | `WYNIK/zasoby/prototyp.css` rozdz. 1–3, 7 |
| Wartości żetonów wymiarów i warstw | `WYNIK/zasoby/zetony/zetony.css` |
| Realne zastosowania wzorców | `WYNIK/05-okna/srodowiska/*.html`, `WYNIK/05-okna/moduly/*.html`, `WYNIK/05-okna/platformowe/*.html` |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o.*
