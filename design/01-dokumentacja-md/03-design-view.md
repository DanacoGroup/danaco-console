# Danaco Console — Design View · widoki i układy projektowe

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) · warstwa funkcjonalna: Danaco Pilot v1.0 |
| **Rodzaj** | Opracowanie merytoryczno-techniczne — warstwa układu (Design View) systemu projektowego |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 · opracowanie A3 (Zespół A — Fundament systemu projektowego) |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | Projektant wizualny (jak zbudować widok) · Deweloper front-end (jak zapisać układ w CSS) · Architekt informacji (jaką masę niesie która strefa) · Kontroler jakości (co sprawdzić na czterech progach) |
| **Zakres** | Taksonomia siedmiu rodzajów widoku · siatka 12-kolumnowa i granica jej obowiązywania · rytm pionowy · pięć wiążących wzorców układu (powłoka środowiska, obszar roboczy z pasem komunikacji, trzy strefy Centrum dowodzenia, asymetria koordynator–wykonawca, panele wielokolumnowe) · zachowanie responsywne na czterech progach · hierarchia wizualna · wzorce pustych stanów, ładowania i błędu na poziomie widoku · gęstość informacji · decyzje projektowe |
| **Czego NIE zawiera** | Definicji żetonów (patrz `zasoby/zetony/zetony.css` i opracowanie A1) · anatomii pojedynczych komponentów `.dn-*` (patrz `zasoby/css/komponenty.css` i opracowanie A2) · treści merytorycznej okien operacyjnych piętnastu modułów (dokumenty modułów) · zawartości trzynastu zakresów Okna Konfiguracji · specyfikacji znaku marki (opracowanie o marce) · kodu produkcyjnego aplikacji klienckiej |

---

## Spis treści

1. [Taksonomia widoków platformy](#1-taksonomia-widoków-platformy)
2. [Siatka: 12 kolumn, przerwa 24 px, treść max 1200 px](#2-siatka-12-kolumn-przerwa-24-px-treść-max-1200-px)
3. [Rytm pionowy: 4 · 12 · 24](#3-rytm-pionowy-4--12--24)
4. [Wzorzec układu A — powłoka środowiska](#4-wzorzec-układu-a--powłoka-środowiska)
5. [Wzorzec układu B — obszar roboczy z pasem komunikacji](#5-wzorzec-układu-b--obszar-roboczy-z-pasem-komunikacji)
6. [Wzorzec układu C — trzy strefy Centrum dowodzenia](#6-wzorzec-układu-c--trzy-strefy-centrum-dowodzenia)
7. [Wzorzec układu D — asymetria koordynator–wykonawca](#7-wzorzec-układu-d--asymetria-koordynatorwykonawca)
8. [Wzorzec układu E — panele wielokolumnowe](#8-wzorzec-układu-e--panele-wielokolumnowe)
9. [Zachowanie responsywne na czterech progach](#9-zachowanie-responsywne-na-czterech-progach)
10. [Hierarchia wizualna: masa, kontrast, pozycja, przestrzeń](#10-hierarchia-wizualna-masa-kontrast-pozycja-przestrzeń)
11. [Puste stany, stany ładowania i stany błędu na poziomie widoku](#11-puste-stany-stany-ładowania-i-stany-błędu-na-poziomie-widoku)
12. [Gęstość informacji — pojemność widoku na progu](#12-gęstość-informacji--pojemność-widoku-na-progu)
13. [Decyzje projektowe](#13-decyzje-projektowe)
- [Załącznik A — katalog zbiorczy wzorców układu](#załącznik-a--katalog-zbiorczy-wzorców-układu)
- [Załącznik B — mapa źródeł](#załącznik-b--mapa-źródeł)
- [Załącznik C — lista sprawdzeń układu przed oddaniem widoku](#załącznik-c--lista-sprawdzeń-układu-przed-oddaniem-widoku)

---

## 1. Taksonomia widoków platformy

### 1.1. Po co taksonomia

Danaco Console nie ma jednego układu. Ma **siedem rodzajów widoku**, z których każdy odpowiada na inne pytanie Operatora i dlatego rządzi się innym prawem układu. Pomylenie rodzaju widoku jest najkosztowniejszym błędem projektowym w kokpicie: widok operacyjny zbudowany jak strona dokumentowa traci połowę pojemności, a widok formularza rozciągnięty na całą szerokość okna przestaje być formularzem.

Taksonomia rozstrzyga trzy rzeczy naraz: **czy obowiązuje siatka 12-kolumnowa**, **czy widok ma miarę czytelności (max 1200 px)** i **czy widok wypełnia całą wysokość okna klienckiego**.

### 1.2. Siedem rodzajów widoku

| # | Rodzaj widoku | Okno źródłowe (KANON rozdz. 7) | Pytanie Operatora | Siatka 12 kol. | Miara treści | Wysokość |
|---|---|---|---|---|:-:|---|
| **W1** | **Widok przejściowy** | Okno startowe (ładowania) | „Czy system żyje?” | nie | oś pionowa, blok centrowany | `100dvh`, treść centrowana |
| **W2** | **Widok formularza** | Okno rejestracji i logowania | „Kim jestem dla systemu?” | nie | `--dn-wym-modal` (560 px) | `100dvh`, blok centrowany |
| **W3** | **Widok rozdzielczy** | Strona główna — Centrum dowodzenia | „Dokąd idę?” | **tak** | `--dn-tresc-max` (1200 px) | przewijalna, minimum `100dvh` |
| **W4** | **Widok powłoki środowiska** | Powłoki TalkIn · WorkSpace · CodeStudio · MultitaskingAI | „W jakim trybie pracuję?” | nie — układ sztywny | pełna szerokość okna | `100dvh`, `overflow: hidden` |
| **W5** | **Widok operacyjny modułu** | Chat Window + okna właściwe modułowi | „Jakim narzędziem realizuję?” | nie — układ sztywny | pełna szerokość obszaru roboczego | wypełnia obszar roboczy |
| **W6** | **Widok nakładki** | Okno Konfiguracji · Okno Ustawień · modale komponentów własnych | „Co zmieniam, nie tracąc kontekstu?” | **tak, wewnątrz nakładki** | 560 px (modal prosty) / 1200 px (nakładka dwudzielna) | `max-height` 90 % okna |
| **W7** | **Widok funkcji globalnej** | Always On Display · Mobile | „Co się dzieje, gdy patrzę gdzie indziej?” | nie (AOD) / jedna kolumna (Mobile) | AOD: `36ch` treści · Mobile: pełna szerokość | AOD: pływający · Mobile: `100dvh` |

### 1.3. Charakterystyka każdego rodzaju

**W1 · Widok przejściowy — okno startowe.**
Stan trwa od ułamka sekundy do kilku sekund. Nie ma paska górnego, nie ma paneli bocznych, nie ma siatki. Jedyna oś to pion: godło marki → wskaźnik stanu → komunikat. Trzy warianty (Łączenie · Powrót z ważnym tokenem · Błąd połączenia) różnią się **wyłącznie zawartością bloku stanu**, nigdy układem — układ musi być stały, aby przejście do stanu błędu nie wywoływało przeskoku treści.

```
        ┌──────────────────── 100dvh ────────────────────┐
        │                                                 │
        │                                                 │
        │                    »». godło                    │  ← blok centrowany
        │                    DANACO CONSOLE               │     w obu osiach
        │                                                 │
        │              ◌  łączenie z serwerem…            │  ← blok stanu — jedyny
        │                                                 │     element zmienny
        │                                                 │
        └─────────────────────────────────────────────────┘
          brak paska · brak paneli · brak siatki 12-kolumnowej
```

**W2 · Widok formularza — okno rejestracji i logowania.**
Rama okna jest kontenerem centrowanym o szerokości `--dn-wym-modal` (560 px). Trzy pod-stany (Rejestracja · Logowanie · Odzyskiwanie konta) i sześć wariantów inwentarza dzielą tę samą ramę — wymienia się wyłącznie wnętrze. Wysokość ramy jest zmienna; **pozycja pionowa ramy jest stała względem środka okna**, żeby dołożenie komunikatu błędu nie przesuwało przycisku pod kursorem.

**W3 · Widok rozdzielczy — Centrum dowodzenia.**
Jedyny widok platformy, w którym siatka 12-kolumnowa obowiązuje w pełnym zakresie. Trzy strefy o malejącej masie w pionie, treść ograniczona do 1200 px, przewijanie pionowe dozwolone. To „przedpokój przed strefą roboczą” — jego zadaniem jest **rozdzielenie ruchu na trzy ścieżki**, nie mieszczenie pracy.

**W4 · Widok powłoki środowiska.**
Układ sztywny, wyprowadzony wprost z żetonów wymiarów. Cztery pasy: pasek górny (48 px) → pas kart sesji (36 px) → boczna nawigacja (224 px) obok obszaru roboczego (`1fr`). Widok wypełnia okno klienckie co do piksela i **nie przewija się jako całość** — przewijają się wyłącznie wnętrza paneli.

**W5 · Widok operacyjny modułu.**
Wnętrze obszaru roboczego powłoki. Zawsze zawiera Chat Window (jedyne okno wspólne wszystkim piętnastu modułom) i od jednego do pięciu okien właściwych modułowi. Rządzi nim wzorzec B (rozdz. 5) albo — dla modułów wielopanelowych — wzorzec E (rozdz. 8).

**W6 · Widok nakładki.**
Dwie odmiany. **Modal prosty** (`.dn-modal`, do 560 px) — potwierdzenia, krótkie formularze, „Pula kont Code CLI”. **Nakładka dwudzielna** — Okno Konfiguracji i Okno Ustawień: lewa kolumna nawigacji zakresów o szerokości `--dn-wym-boczna` (224 px, ta sama miara co boczna nawigacja środowiska) plus prawa kolumna zawartości. Nakładka nigdy nie zastępuje widoku pod spodem — zamknięcie przywraca dokładnie stan sprzed otwarcia.

**W7 · Widok funkcji globalnej.**
Always On Display jest elementem pływającym (`position: fixed`, warstwa `--dn-z-aod` = 1200, zakotwiczony w prawym dolnym rogu z odsunięciem `--dn-od-5`). Nie tworzy przestrzeni roboczej i nie przeładowuje niczego pod spodem. Mobile jest widokiem jednokolumnowym progu `w1`.

### 1.4. Diagram — mapa rodzajów widoku wobec przepływu platformy

```
URUCHOMIENIE
     │
     ▼
 [W1] OKNO STARTOWE ────────────► [W2] REJESTRACJA I LOGOWANIE
     │  token ważny                     │  token wydany / „Pomiń”
     │                                  │
     └──────────────┬───────────────────┘
                    ▼
        [W3] CENTRUM DOWODZENIA  ◄──────────────┐
                    │                            │ powrót ikoną „dom”
     ┌──────────────┼───────────────────┐        │
     ▼ Strefa 1     ▼ Strefa 2          ▼ Strefa 3
[W4] POWŁOKA    [W6] NAKŁADKA      [W6] NAKŁADKA / [W7] FUNKCJA GLOBALNA
 ŚRODOWISKA      komponentu         Konfiguracja ·  Mobile · Always On Display
     │            własnego           Ustawienia
     ▼
[W5] WIDOK OPERACYJNY MODUŁU
     Chat Window + okna właściwe modułowi
     │
     └─► [W7] Always On Display może być przywołany nad KAŻDYM z widoków W3–W6
```

### 1.5. Reguła rozstrzygająca

> **Widok, w którym Operator czyta i wybiera, dostaje siatkę i miarę czytelności.
> Widok, w którym Operator pracuje, dostaje całe okno i sztywny podział na piksele.**

W1, W2, W3, W6 to widoki czytania i wyboru. W4, W5 to widoki pracy. W7 jest warstwą towarzyszącą i nie należy do żadnej z grup — dlatego ma własne prawo (pływanie).

---

## 2. Siatka: 12 kolumn, przerwa 24 px, treść max 1200 px

### 2.1. Definicja wiążąca

| Parametr | Żeton | Wartość | Źródło |
|---|---|---|---|
| Liczba kolumn | `--dn-siatka-kolumny` | **12** | KANON rozdz. 2 · `zetony.css` §7 |
| Przerwa (rynna) | `--dn-siatka-przerwa` | **24 px** (`--dn-od-6`) | KANON rozdz. 2 · `zetony.css` §7 |
| Maksymalna szerokość treści | `--dn-tresc-max` | **1200 px** | KANON rozdz. 2 · `zetony.css` §7 |
| Margines zewnętrzny (w3/w4) | `--dn-od-6` | 24 px | wyprowadzenie z rytmu sekcji |
| Margines zewnętrzny (w1/w2) | `--dn-od-4` | 16 px | decyzja projektowa D-07 |

### 2.2. Arytmetyka kolumny

Przy dwunastu kolumnach i jedenastu rynnach po 24 px szerokość pojedynczej kolumny wynika z prostego rachunku:

```
kolumna = (szerokość_treści − 11 × 24) / 12
```

| Próg | Szerokość okna | Szerokość treści | Kolumna | Uwaga |
|---|---:|---:|---:|---|
| `w4` | 1600 px | 1200 px (limit) | **78,0 px** | marginesy rosną, kolumna stała |
| `w3` | 1280 px | 1200 px (limit) | **78,0 px** | pierwszy próg z limitem czynnym |
| `w2` | 960 px | 912 px (960 − 2 × 24) | **54,0 px** | limit nieaktywny |
| `w1` | 640 px | 608 px (640 − 2 × 16) | **28,7 px** | kolumna poniżej wartości użytkowej — elementy zajmują 12/12 |

**Wniosek projektowy:** powyżej progu `w3` siatka przestaje rosnąć — rośnie wyłącznie margines. To celowe: kokpit „instrumentu pomiarowego” nie rozciąga wiersza tekstu w nieskończoność, tylko utrzymuje stałą miarę czytelności i oddaje nadmiar szerokości na powietrze.

### 2.3. Rozpiętości elementów Centrum dowodzenia

| Element | `w4` | `w3` | `w2` | `w1` |
|---|:-:|:-:|:-:|:-:|
| Karta środowiska (Strefa 1, 4 sztuki) | 3 / 12 | 3 / 12 | 6 / 12 | 12 / 12 |
| Kafel komponentu własnego (Strefa 2, 4 sztuki) | 3 / 12 | 3 / 12 | 6 / 12 | 12 / 12 |
| Listwa ustawień (Strefa 3, 3 pozycje) | 12 / 12 | 12 / 12 | 12 / 12 | 12 / 12 (zawijanie) |
| Nagłówek strefy (etykieta wersalikowa) | 12 / 12 | 12 / 12 | 12 / 12 | 12 / 12 |

**Zasada niezmienności liczby kolumn:** siatka ma dwanaście kolumn na każdym progu. Na progach niższych zmienia się **wyłącznie rozpiętość elementu**, nigdy liczba kolumn ani szerokość rynny. Dzięki temu każdy element pozostaje wyrównany do tej samej osi na wszystkich progach, a projektant ma jeden zestaw linii pomocniczych zamiast czterech.

### 2.4. Kiedy siatka obowiązuje, a kiedy układ jest sztywny

| Widok | Siatka 12 kolumn | Uzasadnienie |
|---|:-:|---|
| **W1** Okno startowe | **nie** | jeden blok na osi pionowej; siatka pozioma nie ma czego porządkować |
| **W2** Logowanie | **nie** | jedna kolumna formularza o stałej mierze 560 px |
| **W3** Centrum dowodzenia | **TAK** | trzy strefy, dwanaście punktów wejścia, treść dokumentowa |
| **W4** Powłoka środowiska | **nie — układ sztywny** | wymiary pochodzą z żetonów (48 · 36 · 224), nie z podziału szerokości |
| **W5** Widok operacyjny modułu | **nie — układ sztywny** | podział wynika z liczby i typologii okien modułu, nie z dwunastu kolumn |
| **W6** Nakładka — modal prosty | **nie** | jedna kolumna, do 560 px |
| **W6** Nakładka — Konfiguracja / Ustawienia | **TAK, wewnątrz prawej kolumny** | prawa kolumna zawartości zakresu jest treścią dokumentową |
| **W7** Always On Display | **nie** | element pływający o szerokości wynikającej z treści (maks. 36 znaków) |
| **W7** Mobile | **nie — jedna kolumna** | próg `w1`; każdy element zajmuje pełną szerokość |

**Reguła graniczna (wiążąca):**

> Siatka 12-kolumnowa nie przekracza progu obszaru roboczego. Wszystko, co jest **ramą kokpitu** (pasek górny, pas kart sesji, boczna nawigacja, panel orkiestracji, pas komunikacji), jest wymiarowane **żetonem**. Wszystko, co jest **treścią dokumentową** (Centrum dowodzenia, zawartość zakresu konfiguracji), jest wymiarowane **siatką**.

### 2.5. Diagram — dwa reżimy wymiarowania

```
 REŻIM ŻETONOWY (kokpit, układ sztywny)     REŻIM SIATKOWY (treść, 12 kolumn)
 ══════════════════════════════════════     ══════════════════════════════════
 ┌──────────────────────────────────┐       ┌──────────────────────────────────┐
 │ pasek górny            48 px      │       │ ◄─── margines ───►               │
 ├──────────────────────────────────┤       │  ┌──┬──┬──┬──┬──┬──┬──┬──┬──┬──┐ │
 │ pas kart sesji         36 px      │       │  │  │  │  │  │  │  │  │  │  │  │ │
 ├────────┬─────────────────────────┤       │  │12 kolumn × 78 px, rynna 24 px │ │
 │ boczna │ obszar roboczy           │       │  └──┴──┴──┴──┴──┴──┴──┴──┴──┴──┘ │
 │ 224 px │ 1fr                      │       │        treść max 1200 px         │
 │        │                          │       │                                  │
 └────────┴─────────────────────────┘       └──────────────────────────────────┘
  wymiar = wartość żetonu                     wymiar = udział w dwunastu kolumnach
  zmiana = zmiana żetonu                      zmiana = zmiana rozpiętości
```

---

## 3. Rytm pionowy: 4 · 12 · 24

### 3.1. Trzy stopnie rytmu

| Stopień | Żeton | Wartość (zwarta) | Wartość (przestronna) | Zastosowanie |
|---|---|---:|---:|---|
| **Jednostka** | `--dn-od-1` | **4 px** | 4 px | najmniejszy dopuszczalny skok; wszystkie inne wartości są jej wielokrotnością |
| **Odstęp panelu** | `--dn-odstep-panel` | **12 px** (`--dn-od-3`) | 20 px (`--dn-od-5`) | rytm wewnętrzny panelu: między wierszem a wierszem, między kontrolką a etykietą |
| **Odstęp sekcji** | `--dn-odstep-sekcji` | **24 px** (`--dn-od-6`) | 32 px (`--dn-od-8`) | rytm między sekcjami widoku: między strefami Centrum dowodzenia, między grupami ustawień |

**Zasada nienegocjowalna:** w całym systemie nie występuje żadna wartość odstępu spoza skali `0 · 4 · 8 · 12 · 16 · 20 · 24 · 32 · 40 · 48 · 64`. Wartość „14 px”, „18 px”, „30 px” jest błędem, nie decyzją.

### 3.2. Rytm pionowy widoku operacyjnego — pełny rachunek

```
 ┌─ krawędź okna klienckiego ────────────────────────────────────────────┐
 │                                                                        │
 │  pasek górny                                              48 px  [żeton]│
 ├────────────────────────────────────────────────────────────────────────┤
 │  pas kart sesji                                           36 px  [żeton]│
 ├────────────────────────────────────────────────────────────────────────┤
 │  ▲ 12 px  margines obszaru roboczego                        (od-3)      │
 │  ┌──────────────────────────────────────────────────────────────────┐ │
 │  │  nagłówek okna operacyjnego            8 + 20 + 8  =  36 px       │ │
 │  ├──────────────────────────────────────────────────────────────────┤ │
 │  │  ▲ 12 px  wewnętrzny odstęp panelu                                │ │
 │  │  wiersz tabeli / listy                          36 px  [żeton]    │ │
 │  │  wiersz tabeli / listy                          36 px             │ │
 │  │  wiersz tabeli / listy                          36 px             │ │
 │  │  ▼ 12 px                                                          │ │
 │  └──────────────────────────────────────────────────────────────────┘ │
 │  ▲ 12 px  odstęp między oknem operacyjnym a pasem komunikacji         │
 │  ┌──────────────────────────────────────────────────────────────────┐ │
 │  │  pas komunikacji — Chat Window        minmax(320px, 54%)          │ │
 │  └──────────────────────────────────────────────────────────────────┘ │
 │  ▼ 12 px                                                               │
 └────────────────────────────────────────────────────────────────────────┘
```

### 3.3. Tabela rytmu — co jakim odstępem

| Para elementów | Odstęp | Żeton |
|---|---:|---|
| Etykieta pola ↔ kontrolka pola | 4 px | `--dn-od-1` |
| Ikona ↔ etykieta w jednym wierszu | 8 px | `--dn-od-2` |
| Wiersz listy ↔ wiersz listy (odstęp wewnętrzny panelu) | 12 px | `--dn-odstep-panel` |
| Panel ↔ panel w obszarze roboczym | 12 px | `--dn-od-3` |
| Nagłówek okna ↔ pierwszy wiersz treści | 12 px | `--dn-odstep-panel` |
| Karta ↔ karta w siatce Centrum dowodzenia | 24 px | `--dn-siatka-przerwa` |
| Strefa ↔ strefa Centrum dowodzenia | 24 px (minimum) / 48 px (zalecane) | `--dn-odstep-sekcji` / `--dn-od-12` |
| Blok treści ↔ stopka widoku | 64 px | `--dn-od-16` |

### 3.4. Dlaczego 12, a nie 16

Odstęp panelu 12 px wynika wprost z pokrętła `GESTOSC_WIZUALNA = 8/10`. Przy odstępie 16 px kokpit traci na progu `w3` około jednego wiersza tabeli w każdym panelu — a to jest dokładnie ta jednostka informacji, dla której Operator otwiera okno. Wariant 20 px jest przygotowany żetonem (`data-gestosc="przestronna"`) i włączany świadomie, nie domyślnie.

---

## 4. Wzorzec układu A — powłoka środowiska

### 4.1. Pełna specyfikacja

```css
/* Powłoka środowiska — TalkIn · WorkSpace · CodeStudio · MultitaskingAI */
.powloka {
  display: grid;
  grid-template-rows: var(--dn-wym-pasek) var(--dn-wym-pas-kart) 1fr;  /* 48px 36px 1fr */
  height: 100dvh;
}
.powloka-cialo {
  display: grid;
  grid-template-columns: var(--dn-wym-boczna) 1fr;                      /* 224px 1fr */
  min-height: 0;                                                        /* warunek przewijania wnętrz */
}
```

| Pas | Wymiar | Żeton | Rola | Przewijanie |
|---|---:|---|---|---|
| Pasek górny | 48 px | `--dn-wym-pasek` | rama kokpitu: godło, logotyp, dom, wyszukiwanie, szybka konfiguracja sesji, Mobile, Always On Display, motyw | nie |
| Pas kart sesji | 36 px | `--dn-wym-pas-kart` | karty sesji + kontrolka „+” | poziomo, gdy kart jest więcej niż mieści pas |
| Boczna nawigacja | 224 px | `--dn-wym-boczna` | lista modułów (9 / 9 / 8) albo Panel orkiestracji (6 sekcji) | pionowo, wewnątrz panelu |
| Obszar roboczy | `1fr` | — | okna operacyjne modułu + Chat Window | pionowo, wewnątrz okien |

### 4.2. Diagram szkieletu

```
 ╔════════════════════════════════════════════════════════════════════════╗
 ║ »». DANACO CONSOLE │ dom │ [szukaj…] │ ustawienia·telefon·AOD·motyw    ║  48 px
 ╠════════════════════════════════════════════════════════════════════════╣
 ║ [● Studio ✕] [ Research ✕] [ + ]                                        ║  36 px
 ╠══════════════════╦═════════════════════════════════════════════════════╣
 ║ TalkIn            ║                                                     ║
 ║ ─────────────────║   OBSZAR ROBOCZY                                    ║
 ║  ▸ Studio         ║   (wzorzec B — rozdz. 5)                            ║  1fr
 ║    Workspace      ║                                                     ║
 ║    Browser        ║                                                     ║
 ║    Research       ║                                                     ║
 ║    Library        ║                                                     ║
 ║    Translate      ║                                                     ║
 ║    Roundtable     ║                                                     ║
 ║    Assistant      ║                                                     ║
 ║    Agents         ║                                                     ║
 ╚══════════════════╩═════════════════════════════════════════════════════╝
   224 px                              1fr
```

### 4.3. Warstwy (z-index) w powłoce

| Warstwa | Żeton | Wartość | Element |
|---|---|---:|---|
| Podłoga | `--dn-z-podloga` | 0 | obszar roboczy, okna operacyjne |
| Przybornik | `--dn-z-przybornik` | 10 | przyklejone nagłówki tabel, przybornik promptu |
| Pasek | `--dn-z-pasek` | 100 | pasek górny (rama) |
| Boczna | `--dn-z-boczna` | 200 | boczna nawigacja / Panel orkiestracji |
| Pas komunikacji | `--dn-z-pas-komunikacji` | 300 | Chat Window |
| Nakładka | `--dn-z-nakladka` | 800 | przyciemnienie pod modalem |
| Modal | `--dn-z-modal` | 900 | Okno Konfiguracji, Okno Ustawień, modale |
| Powiadomienie | `--dn-z-powiadomienie` | 1000 | toasty |
| Tooltip | `--dn-z-tooltip` | 1100 | dymki objaśnień `[?]` |
| Always On Display | `--dn-z-aod` | 1200 | pływający awatar |
| Centrum poleceń | `--dn-z-centrum-polecen` | 1300 | zawsze najwyżej |

### 4.4. Cztery powłoki — jedna mechanika, dwie zawartości

| Powłoka | Pasek górny | Pas kart | Panel boczny | Zawartość obszaru roboczego |
|---|---|---|---|---|
| **TalkIn** | wzorzec wspólny | karta = moduł | boczna nawigacja, **9** pozycji | okna modułu + Chat Window |
| **WorkSpace** | wzorzec wspólny | karta = moduł | boczna nawigacja, **9** pozycji | okna modułu + Chat Window |
| **CodeStudio** | wzorzec wspólny | karta = moduł | boczna nawigacja, **8** pozycji | okna modułu + Chat Window |
| **MultitaskingAI** | wzorzec wspólny — **bez zmian** | karta = **proces orkiestracji** | **Panel orkiestracji**, 6 sekcji | zawartość wybranej sekcji + Chat Window per rola |

**Rozstrzygnięcie źródłowe:** różnica między powłoką modułową a MultitaskingAI leży wyłącznie w **semantyce karty sesji i zawartości panelu bocznego**. Wymiary pasów są identyczne. To oznacza jeden kod układu i cztery konfiguracje treści.

### 4.5. Nagłówek panelu bocznego

| Powłoka | Wiersz 1 (krój nagłówkowy) | Wiersz 2 (krój bazowy) |
|---|---|---|
| TalkIn / WorkSpace / CodeStudio | nazwa środowiska | — |
| MultitaskingAI | `MultitaskingAI` | `Panel orkiestracji` (podtytuł stały) |

Nagłówek zajmuje 36 px (jeden wiersz) albo 56 px (dwa wiersze, MultitaskingAI) i **nie przewija się razem z listą pozycji** — jest przyklejony do górnej krawędzi panelu.

---

## 5. Wzorzec układu B — obszar roboczy z pasem komunikacji

### 5.1. Pełna specyfikacja

```css
/* Obszar roboczy: okna operacyjne modułu nad pasem komunikacji */
.obszar-roboczy {
  display: grid;
  grid-template-rows: minmax(0, 1fr) minmax(320px, 54%);
  min-width: 0;
  min-height: 0;
}
```

| Pas | Wymiar | Minimum | Rola |
|---|---|---:|---|
| Okna operacyjne modułu | `minmax(0, 1fr)` | 0 | od jednego do pięciu okien właściwych modułowi |
| Pas komunikacji (Chat Window) | `minmax(320px, 54%)` | **320 px** (`--dn-wym-pas-komunikacji`) | historia rozmowy + pole poleceń |

### 5.2. Dlaczego 54 %, a nie połowa

Wartość 54 % jest przesunięciem o cztery punkty procentowe na korzyść pasa komunikacji względem podziału równego. Uzasadnienie wynika z kontraktu produktu: Chat Window jest **jedynym oknem operacyjnym wspólnym wszystkim piętnastu modułom** — jest stałym punktem interakcji Operatora, podczas gdy okna modułowe zmieniają się przy każdym przełączeniu pozycji bocznej nawigacji. Podział równy (50/50) czytałby się jako „dwa równorzędne panele”; 54 % czyta się jako „tu rozmawiam, tam patrzę na wynik”. Jednocześnie `minmax(0, 1fr)` na górze gwarantuje, że okno modułowe nigdy nie zniknie — skurczy się, ale pozostanie.

Dolny limit 320 px (`--dn-wym-pas-komunikacji`) to wartość, przy której pas komunikacji mieści komplet: dwa wpisy historii + pole promptu + przybornik. Poniżej tej wartości Chat Window przestaje być użyteczny i wzorzec przechodzi na tryb zakładek (rozdz. 9.4).

### 5.3. Diagram i zachowanie przy zmianie modułu

```
 ┌───────────────────────────────────────────────────────────────────┐
 │  OKNO OPERACYJNE MODUŁU                     minmax(0, 1fr)         │
 │  ┌──────────────────────────────────────────────────────────────┐ │
 │  │ Studio Editor                                       ⋯          │ │  ← nagłówek 36 px
 │  ├──────────────────────────────────────────────────────────────┤ │
 │  │ treść dokumentu (przykładowa sesja robocza)                    │ │
 │  └──────────────────────────────────────────────────────────────┘ │
 ├───────────────────────────────────────────────────────────────────┤
 │  PAS KOMUNIKACJI — Chat Window          minmax(320px, 54%)         │
 │  ┌──────────────────────────────────────────────────────────────┐ │
 │  │ historia rozmowy (przewijana)                                  │ │
 │  ├──────────────────────────────────────────────────────────────┤ │
 │  │ ❯  [ pole poleceń…                       ]  spinacz  wyslij  │ │
 │  └──────────────────────────────────────────────────────────────┘ │
 └───────────────────────────────────────────────────────────────────┘

 PRZY ZMIANIE MODUŁU:
   górny pas  ── wymienia się w całości (okna modułu A → okna modułu B)
   dolny pas  ── NIE znika; rekonfiguruje kontekst poleceń i operacji
```

### 5.4. Odmiany wzorca B

| Odmiana | Kiedy | Podział | Uzasadnienie |
|---|---|---|---|
| **B1 — pionowa (domyślna)** | moduły z jednym oknem wiodącym: Studio, Research, Library, Browser, Developer, Design, Apps | `minmax(0,1fr)` / `minmax(320px,54%)` | wynik nad rozmową — kierunek czytania z góry na dół |
| **B2 — pionowa z pasmem narzędzi** | moduły z panelem pomocniczym: Studio (Tools Panel), Research (Sources Manager) | okno modułu dzieli się poziomo `1fr / 280px` | panel pomocniczy jest wąski i stały; nie konkuruje o wysokość |
| **B3 — zakładkowa** | próg `w2` i niżej | Chat Window i okna modułu jako `.dn-zakladki` | poniżej 320 px pas komunikacji przestaje być użyteczny |
| **B4 — dwukolumnowa** | próg `w4` — „dwa okna komunikacji” | pas komunikacji dzieli się na dwie kolumny `1fr 1fr` | szerokie biurko pozwala prowadzić dwie rozmowy równolegle (Roundtable, MultitaskingAI) |

### 5.5. Warunek techniczny `min-height: 0`

Każdy kontener siatki, w którym dziecko ma przewijać własne wnętrze, **musi** mieć `min-height: 0` (i `min-width: 0` przy podziale kolumnowym). Bez tego domyślna wartość `auto` sprawia, że dziecko rozpycha rodzica zamiast się przewijać, a pas komunikacji wypada poza dolną krawędź okna. To najczęstszy błąd wdrożeniowy tego wzorca.

---

## 6. Wzorzec układu C — trzy strefy Centrum dowodzenia

### 6.1. Malejąca masa jako komunikat

Trzy strefy Centrum dowodzenia różnią się formą prezentacji celowo. Hierarchia skali **karty → kafle → listwa** komunikuje wagę, zanim Operator przeczyta którykolwiek napis. To jedyne miejsce w platformie, w którym anty-domyślne „trzy równe karty funkcji” zostaje zastąpione przez „strefy o malejącej masie”.

### 6.2. Specyfikacja trzech stref

| Strefa | Zawartość | Forma | Liczba | Rozpiętość (`w3`/`w4`) | Wysokość elementu | Krój tytułu |
|---|---|---|---|---|---|---|
| **Strefa 1** | karty środowisk | `.dn-karta-srodowiska` | 4 | 3 / 12 kolumn | wysoka (godło + tytuł + motto + opis) | **Space Grotesk** `--dn-fs-3xl` |
| **Strefa 2** | kafle komponentów własnych | `.dn-kafel` | 4 | 3 / 12 kolumn | średnia (ikona + etykieta + opis) | IBM Plex Sans `--dn-fs-base` półgruby |
| **Strefa 3** | listwa ustawień | `.dn-listwa` + `.dn-listwa-pozycja` | 3 | 12 / 12 kolumn | jeden wiersz | IBM Plex Sans `--dn-fs-base` |

### 6.3. Zawartość źródłowa stref

**Strefa 1 — cztery karty środowisk (kolejność źródłowa, nienaruszalna):**

| # | Środowisko | Opis trybu pracy | Prowadzi do |
|---|---|---|---|
| 1 | **TalkIn** | Wiedza, komunikacja i praca z treścią | powłoka z 9 modułami |
| 2 | **WorkSpace** | Produktywność, organizacja i realizacja projektów | powłoka z 9 modułami |
| 3 | **CodeStudio** | Programowanie | powłoka z 8 modułami |
| 4 | **MultitaskingAI** | Orkiestracja autonomicznej pracy ciągłej | powłoka z Panelem orkiestracji (6 sekcji) |

**Strefa 2 — cztery kafle komponentów własnych:**

| # | Komponent własny | Etykieta działania | Okno docelowe |
|---|---|---|---|
| 1 | **Automations** | „Zbuduj automatykę” | Workflow Builder |
| 2 | **Agents** | „Skonfiguruj agenta” | Agent Builder |
| 3 | **Workspace** | „Załóż projekt” | Project Dashboard |
| 4 | **Assistant** | „Ustaw profil asystenta” | Voice Console |

**Strefa 3 — trzy pozycje listwy:** Okno konfiguracji · Mobile · Always On Display.

### 6.4. Diagram szkieletu

```
 ┌──────────────────────────────── treść max 1200 px ────────────────────────────┐
 │                                                                                │
 │  STREFA 1 · KARTY ŚRODOWISK                             masa: ████████ 100 %   │
 │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐                       │
 │  │ »». godło│  │ »». godło│  │ »». godło│  │ »». godło│    4 × 3/12 kolumn     │
 │  │ TalkIn   │  │WorkSpace │  │CodeStudio│  │Multitask.│    rynna 24 px         │
 │  │ motto    │  │ motto    │  │ motto    │  │ motto    │    Space Grotesk 30 px │
 │  │ opis     │  │ opis     │  │ opis     │  │ opis     │                        │
 │  └──────────┘  └──────────┘  └──────────┘  └──────────┘                       │
 │                                                                                │
 │  ▲ 48 px odstęp międzystrefowy                                                 │
 │                                                                                │
 │  STREFA 2 · KAFLE KOMPONENTÓW WŁASNYCH                  masa: ████ 50 %        │
 │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐                           │
 │  │ ikona   │  │ ikona   │  │ ikona   │  │ ikona   │        4 × 3/12 kolumn    │
 │  │Automat. │  │ Agents  │  │Workspace│  │Assistant│        Plex Sans 13 px    │
 │  │ „Zbuduj”│  │„Skonfig”│  │ „Załóż” │  │ „Ustaw” │                           │
 │  └─────────┘  └─────────┘  └─────────┘  └─────────┘                           │
 │                                                                                │
 │  ▲ 48 px odstęp międzystrefowy                                                 │
 │                                                                                │
 │  STREFA 3 · LISTWA USTAWIEŃ                             masa: ██ 20 %          │
 │  ┌──────────────────────────────────────────────────────────────────────────┐ │
 │  │ [ustawienia] Okno konfiguracji │ [telefon] Mobile │ ◉ Always On Display  │ │
 │  └──────────────────────────────────────────────────────────────────────────┘ │
 │                                                          jeden wiersz, 12/12  │
 └────────────────────────────────────────────────────────────────────────────────┘
```

### 6.5. Miary malejącej masy

Masa wizualna strefy nie jest odczuciem — jest iloczynem czterech mierzalnych składników. Poniższa tabela pokazuje, którym z nich sterujemy w każdej strefie.

| Składnik masy | Strefa 1 | Strefa 2 | Strefa 3 |
|---|---|---|---|
| **Powierzchnia elementu** | duża (godło + 3 wiersze tekstu) | średnia (ikona + 2 wiersze) | mała (jeden wiersz) |
| **Stopień pisma tytułu** | `--dn-fs-3xl` (30 px) | `--dn-fs-base` (13 px) | `--dn-fs-base` (13 px) |
| **Krój tytułu** | Space Grotesk (nagłówkowy) | IBM Plex Sans (bazowy) | IBM Plex Sans (bazowy) |
| **Promień narożnika** | `--dn-r-xl` (14 px) | `--dn-r-lg` (10 px) | `--dn-r-md` (8 px) |
| **Sygnał przy najechaniu** | wstęga górna 2 px + uniesienie | zmiana obrysu i tła | zmiana tła |
| **Udział w polu widzenia** | ≈ 50 % | ≈ 30 % | ≈ 10 % |

### 6.6. Odstęp międzystrefowy

Minimum to `--dn-odstep-sekcji` (24 px). Zalecenie dla Centrum dowodzenia to **48 px** (`--dn-od-12`) — podwojony odstęp sekcji. Uzasadnienie: strefy są trzema **niezależnymi ścieżkami działania**, nie trzema akapitami jednej sekcji. Podwojony odstęp to najtańszy sposób zakomunikowania niezależności bez dokładania linii ani ramek.

---

## 7. Wzorzec układu D — asymetria koordynator–wykonawca

### 7.1. Rozstrzygnięcie fundamentalne

Podział na rolę wykonującą (**Executor**) i rolę zarządzającą (**Coordinator**) jest rozstrzygnięciem fundamentalnym środowiska MultitaskingAI: obie odpowiedzialności są celowo rozdzielone, aby żaden pojedynczy model nie musiał jednocześnie planować i realizować. Układ ma tę asymetrię **pokazać**, a nie zatrzeć symetrycznym podziałem na pół.

### 7.2. Pełna specyfikacja

```css
/* Plansza roli — asymetria pionowa: para okien komunikacji nad pasem sterowania */
.plansza {
  display: grid;
  grid-template-rows: minmax(0, 58%) minmax(0, 42%);
  gap: var(--dn-od-3);
  padding: var(--dn-od-3);
  min-height: 0;
}
/* Para okien: Coordinator obok Executora — równe kolumny, różna waga obrysu */
.para { display: grid; grid-template-columns: 1fr 1fr; gap: var(--dn-od-3); min-height: 0; }
/* Pas dolny: kolejka kroków szersza od panelu sterowania */
.pas-dolny { display: grid; grid-template-columns: 1.25fr 1fr; gap: var(--dn-od-3); min-height: 0; }
```

| Element | Wymiar | Uzasadnienie asymetrii |
|---|---|---|
| Para okien komunikacji | **58 %** wysokości | rozmowa ról jest przedmiotem obserwacji Operatora |
| Pas dolny (kolejka + sterowanie) | **42 %** wysokości | narzędzia towarzyszą, nie dominują |
| Coordinator Chat ↔ Executor Chat | **1fr : 1fr** (szerokość równa) | równa szerokość = równa czytelność strumienia; różnicę niesie **obrys**, nie rozmiar |
| Kolejka kroków ↔ panel sterowania | **1,25fr : 1fr** | kolejka jest dziennikiem przebiegu — wymaga więcej znaków w wierszu |

### 7.3. Gdzie leży asymetria — trzy poziomy

| Poziom | Coordinator | Executor | Nośnik różnicy |
|---|---|---|---|
| **Obrys górny okna** | 2 px `--dn-atrament` | 2 px `--dn-sygnal-wypelnienie` | barwa + pozycja (wstęga górna) |
| **Medalion nadawcy** | ikona `wykonawca` w wariancie koordynatora | ikona `agent` | ikona, nie barwa tła |
| **Plakietka roli** | `.dn-plakietka--rola` z nazwą `Coordinator` | `.dn-plakietka--rola` z nazwą `Executor 1` | etykieta tekstowa |

**Zasada:** stan nigdy samym kolorem. Różnica ról jest zakodowana **trzykrotnie** — obrysem, ikoną i etykietą — więc czyta się także w skali szarości i przy zaburzeniach postrzegania barw.

### 7.4. Szyna przekazania

Moment przekazania zlecenia od Coordinatora do Executora jest widoczny **w obu oknach jednocześnie**: pozioma szyna zakotwiczona na styku dwóch kolumn, na warstwie `--dn-z-przybornik` (10), z kropką sygnału w tętnie 2,4 s. To jedyny ruch ciągły w tym widoku i jedyne miejsce, w którym para okien komunikuje się wizualnie ponad podziałem kolumnowym.

```
 ┌──────────────────────────┬──────────────────────────┐
 │  Coordinator             │  Executor 1              │
 │  ▔▔▔▔▔ obrys atrament    │  ▔▔▔▔▔ obrys sygnał      │
 │                ╭─────────┴─────────╮                │
 │                │ ● przekazanie #128 │  ← szyna       │  58 %
 │                ╰─────────┬─────────╯    (tętno 2,4 s)│
 │  historia rozmowy         │  historia rozmowy        │
 │  ❯ [ polecenie…    ] ➤    │  ❯ [ polecenie…    ] ➤   │
 ├──────────────────────────┴──────────────────────────┤
 │  KOLEJKA KROKÓW    1,25fr    │  STEROWANIE     1fr    │  42 %
 │  ○ etap 1 · poprawny         │  start · stop · pauza  │
 │  ◐ etap 2 · pracuje          │  wznowienie ·          │
 │  ○ etap 3 · oczekuje         │  przekazanie ·         │
 │                              │  powtórzenie · walid.  │
 └──────────────────────────────┴────────────────────────┘
```

### 7.5. Cztery role i piąta pozycja zespołu

| Rola | Okno robocze | Miejsce w układzie |
|---|---|---|
| **Executor 1** | Executor Chat | kolumna prawa pary (tor główny) |
| **Executor 2** | Executor Chat | druga plansza — równoległy tor pracy |
| **Coordinator** | Coordinator Chat | kolumna lewa pary (stała) |
| **Executor 3 / Validator** | Results Analyzer | plansza wywoływana po zakończeniu tury |
| **Subagent Network** | rozwinięcie pod kartą roli | zagnieżdżone pozycje `.dn-karta--pozycja`, **nie osobna sekcja panelu** |

Subagent Network jest rozwijany **wewnątrz** karty roli Executor 1 lub Executor 2 — do 15 podagentów. Układ nie tworzy dla nich osobnej planszy, bo w hierarchii są podrzędne wobec roli, a nie równorzędne wobec niej.

---

## 8. Wzorzec układu E — panele wielokolumnowe

### 8.1. Dwa moduły, jeden wzorzec

| Moduł | Okno wielokrotne | Punkt odniesienia | Panel towarzyszący |
|---|---|---|---|
| **Roundtable** | **Model Panels** (N instancji) | Debate Panel (poniżej, pełna szerokość) | Moderator Panel · Consensus Panel |
| **Translate** | **Translation Panels** (N instancji) | **Source Panel** (powyżej, pełna szerokość) | Glossary Manager |

W obu przypadkach liczba instancji zależy od decyzji Operatora: od liczby wybranych uczestników debaty (Roundtable) albo liczby wybranych języków (Translate). Układ musi więc być **parametryczny względem N**, a nie zaprojektowany pod trzy panele.

### 8.2. Specyfikacja

```css
/* Siatka równoległa paneli — parametryczna względem liczby instancji */
.panele-rownolegle {
  display: grid;
  grid-template-columns: repeat(var(--liczba-paneli, 3), minmax(240px, 1fr));
  gap: var(--dn-od-3);
  min-height: 0;
  overflow-x: auto;          /* powyżej pojemności progu — przewijanie poziome */
}
```

| Parametr | Wartość | Uzasadnienie |
|---|---:|---|
| Szerokość minimalna panelu | **240 px** | poniżej tej wartości nagłówek panelu (kod języka + status + selektor tonu) łamie się na dwa wiersze |
| Szerokość maksymalna panelu | `1fr` | panele dzielą dostępną szerokość równo — żaden uczestnik debaty nie ma przywileju |
| Przerwa między panelami | 12 px (`--dn-od-3`) | rytm paneli obszaru roboczego, nie rynna siatki dokumentowej |
| Zachowanie powyżej pojemności | przewijanie poziome | liczba paneli jest decyzją Operatora — układ nie może jej ograniczać (zasada zero blokad) |

### 8.3. Pojemność siatki równoległej na progu

Rachunek: szerokość obszaru roboczego = szerokość okna − boczna nawigacja (224 px) − marginesy (2 × 12 px). Pojemność = `floor((obszar + 12) / (240 + 12))`.

| Próg | Szerokość okna | Obszar roboczy | Paneli bez przewijania | Zachowanie powyżej |
|---|---:|---:|:-:|---|
| `w4` | 1600 px | 1352 px | **5** | przewijanie poziome |
| `w3` | 1280 px | 1032 px | **4** | przewijanie poziome |
| `w2` | 960 px | 880 px (boczna zwinięta do 56 px) | **3** | przewijanie poziome |
| `w1` | 640 px | 608 px (boczna ukryta) | **2** | przewijanie poziome + zakładki panelu |

### 8.4. Diagram — Roundtable

```
 ┌─ obszar roboczy ─────────────────────────────────────────────────────┐
 │  [ + Dodaj uczestnika ▾ ]                                             │
 │  ┌────────────┬────────────┬────────────┐                            │
 │  │ Model      │ Model      │ Model      │   ◄ Model Panels           │
 │  │ Panel A    │ Panel B    │ Panel C    │     repeat(N, minmax(240px,1fr))
 │  └────────────┴────────────┴────────────┘                            │
 │  ┌──────────────────────────────────────────────────────────────────┐│
 │  │  Debate Panel — rejestr wymiany argumentów      12/12 szerokości  ││
 │  └──────────────────────────────────────────────────────────────────┘│
 │  ┌───────────────────────┐        ┌───────────────────────┐          │
 │  │ Moderator Panel        │        │ Consensus Panel        │          │
 │  │ (panel boczny, wąski)  │        │ (wywoływany po turze)  │          │
 │  └───────────────────────┘        └───────────────────────┘          │
 ├──────────────────────────────────────────────────────────────────────┤
 │  Chat Window — pas komunikacji                                        │
 └──────────────────────────────────────────────────────────────────────┘
```

### 8.5. Diagram — Translate

```
 ┌─ obszar roboczy ─────────────────────────────────────────────────────┐
 │  ┌──────────────────────────────────────────────────────────────────┐│
 │  │  Source Panel — tekst źródłowy, punkt odniesienia                 ││  ← pas górny
 │  └──────────────────────────────────────────────────────────────────┘│     stały
 │  [+ Dodaj język ▾]              Zsynchronizowane przewijanie: [✓]     │
 │  ┌──────────────┬──────────────┬──────────────┐                      │
 │  │ Translation  │ Translation  │ Translation  │  ◄ Translation Panels│
 │  │ Panel — EN   │ Panel — DE   │ Panel — FR   │    siatka równoległa │
 │  │ [1] segment  │ [1] segment  │ [1] segment  │    segmenty w RÓWNYCH│
 │  │ [2] segment  │ [2] segment  │ [2] segment  │    wierszach — wspólna
 │  │ [3] segment  │ [3] segment  │ [3] segment  │    linia bazowa      │
 │  └──────────────┴──────────────┴──────────────┘                      │
 ├──────────────────────────────────────────────────────────────────────┤
 │  Chat Window — pas komunikacji                                        │
 └──────────────────────────────────────────────────────────────────────┘
```

### 8.6. Warunek wyrównania segmentów

W Translate segment o numerze `[n]` musi leżeć **w tym samym wierszu we wszystkich panelach**, bo kliknięcie segmentu w Source Panel przewija wszystkie Translation Panels do odpowiadającego segmentu. Wynika z tego wymóg układu: wysokość wiersza segmentu jest wspólna dla całej siatki i równa najwyższemu wystąpieniu segmentu w danym numerze. Rozwiązanie: `display: grid` na poziomie **siatki paneli**, nie każdego panelu z osobna, z jednym rzędem na segment (`grid-auto-rows: auto` i `subgrid` w kolumnach paneli tam, gdzie silnik prezentacji to wspiera; zapasowo — synchronizacja wysokości wierszy skryptem).

---

## 9. Zachowanie responsywne na czterech progach

### 9.1. Cztery progi

| Próg | Żeton | Szerokość | Kontekst | Widok wiodący |
|---|---|---:|---|---|
| **w1** | `--dn-bp-w1` | 640 px | telefon poziomo | funkcja globalna **Mobile** |
| **w2** | `--dn-bp-w2` | 960 px | tablet | powłoka z **boczną zwiniętą do ikon** |
| **w3** | `--dn-bp-w3` | 1280 px | biurko | **pełny kokpit** — próg odniesienia |
| **w4** | `--dn-bp-w4` | 1600 px | szerokie biurko | **dwa okna komunikacji** |

**Próg odniesienia to `w3`.** Każdy wzorzec układu jest projektowany na 1280 px, a pozostałe progi opisują odstępstwo od niego: `w4` dokłada, `w2` zwija, `w1` wymienia widok.

### 9.2. Tabela zbiorcza — co znika, co się zwija, co się przesuwa

| Element | `w4` 1600 | `w3` 1280 | `w2` 960 | `w1` 640 |
|---|---|---|---|---|
| **Pasek górny — godło + logotyp** | pełny (`DANACO` + `CONSOLE`) | pełny | **zwija się** do sygnetu + `DANACO` | **zwija się** do sygnetu uproszczonego |
| **Pasek górny — pole wyszukiwania** | pełne, do 520 px | pełne, do 400 px | **zwija się** do przycisku ikonowego `szukaj` | **zwija się** do przycisku ikonowego |
| **Pasek górny — przyciski prawe** | 5 ikon (ustawienia · Mobile · AOD · gęstość · motyw) | 5 ikon | 3 ikony + `wiecej` | 2 ikony + `wiecej` |
| **Pas kart sesji** | wszystkie karty widoczne | przewijanie poziome powyżej pojemności | przewijanie poziome; tytuł skracany | **zwija się** do listy rozwijanej „karty sesji (n)” |
| **Boczna nawigacja / Panel orkiestracji** | 224 px, tekst + ikona | 224 px, tekst + ikona | **zwija się do 56 px** — same ikony + dymek | **znika**; wywoływana przyciskiem `menu` jako nakładka |
| **Nagłówek panelu bocznego** | pełny (nazwa + podtytuł) | pełny | tylko emblemat środowiska | w nakładce — pełny |
| **Obszar roboczy — okna modułu** | do 5 paneli równolegle | do 4 paneli | do 3 paneli, przewijanie poziome | 1 panel, pozostałe jako `.dn-zakladki` |
| **Pas komunikacji (Chat Window)** | **dzieli się na dwie kolumny** (B4) | `minmax(320px, 54%)` | `minmax(320px, 54%)`, przy niedoborze → zakładki (B3) | zakładka „Rozmowa” obok zakładki „Wynik” |
| **Strefa 1 — karty środowisk** | 4 × 3/12 | 4 × 3/12 | 2 × 6/12 (dwa rzędy) | 1 × 12/12 (cztery rzędy) |
| **Strefa 2 — kafle** | 4 × 3/12 | 4 × 3/12 | 2 × 6/12 | 1 × 12/12 |
| **Strefa 3 — listwa** | 3 pozycje w wierszu | 3 pozycje w wierszu | 3 pozycje w wierszu | **zawija się** do trzech wierszy |
| **Modal prosty** | 560 px, wyśrodkowany | 560 px | 560 px | pełna szerokość − 2 × 16 px, przyklejony do dołu |
| **Nakładka dwudzielna (Konfiguracja)** | 1200 px, układ 224 + 1fr | 1200 px, układ 224 + 1fr | pełna szerokość, układ 56 + 1fr | **rozdziela się** na dwa ekrany: lista zakresów → zawartość zakresu |
| **Always On Display** | pływający, prawy dolny róg | pływający | pływający, zwężony do rdzenia + jedna linia | **zwija się** do samego rdzenia (36 px) |
| **Toasty** | prawy górny, 420 px | prawy górny, 420 px | prawy górny, 320 px | pełna szerokość − 2 × 16 px, u dołu |
| **Dymki objaśnień `[?]`** | najechanie + fokus | najechanie + fokus | najechanie + fokus | **dotknięcie** — dymek trwały do zamknięcia |
| **Siatka równoległa paneli** | 5 paneli | 4 panele | 3 panele | 2 panele + przewijanie |

### 9.3. Trzy rodzaje zmiany — słownik

| Rodzaj | Definicja | Kiedy dozwolony |
|---|---|---|
| **Zwinięcie** | element pozostaje obecny i dostępny, traci etykietę albo szerokość | gdy funkcja jest potrzebna, ale nie w pełnej formie |
| **Przesunięcie** | element zmienia miejsce w układzie, zachowując formę | gdy kolejność czytania na wąskim ekranie jest inna niż na szerokim |
| **Zniknięcie** | element przestaje być renderowany w tym progu | **wyłącznie** wtedy, gdy pozostaje wywoływalny innym mechanizmem (przycisk `menu`, lista rozwijana, zakładka) |

**Zasada bezwzględna:** żadne zniknięcie nie może odebrać dostępu do funkcji. To bezpośrednia konsekwencja ADL-017 (zero blokad) przeniesiona na warstwę układu: element może zniknąć z widoku, nie może zniknąć z systemu.

### 9.4. Reguła przejścia B1 → B3 (pas komunikacji na zakładki)

Przejście nie jest wywoływane szerokością okna, tylko **wysokością dostępną dla obszaru roboczego**:

```
wysokość_obszaru = wysokość_okna − 48 (pasek) − 36 (pas kart) − 24 (marginesy)

jeżeli   wysokość_obszaru  <  320 (pas komunikacji) + 200 (minimum okna modułu)
wówczas  wzorzec B1  →  wzorzec B3 (zakładki)
```

Próg wynosi więc **628 px wysokości okna**. Poniżej tej wartości Chat Window i okno modułu dzielą tę samą powierzchnię jako dwie zakładki `.dn-zakladki`, a przełączanie między nimi nie gubi stanu żadnej z nich.

### 9.5. Punkt dotykowy — nakładka na każdym progu

`@media (pointer: coarse)` nie jest progiem szerokości i **nakłada się na próg bieżący**, nie zastępuje go. Podnosi wyłącznie cele dotykowe: kontrolka i przycisk ikonowy do 40 px, wiersz do 44 px, pole wyboru do 20 px, przełącznik do 44 × 24 px. Konsekwencja układowa: na tablecie (`w2`) w trybie dotykowym pojemność wiersza spada o około 18 %, co jest ujęte w rachunku rozdz. 12.

---

## 10. Hierarchia wizualna: masa, kontrast, pozycja, przestrzeń

### 10.1. Cztery narzędzia i ich budżet

Kokpit ma cztery — i tylko cztery — narzędzia budowania hierarchii. Każde ma określony budżet, bo `WARIANCJA_PROJEKTOWA = 4/10` nie pozwala na nieograniczone różnicowanie.

| Narzędzie | Czym się steruje | Budżet w kokpicie | Kiedy stosować |
|---|---|---|---|
| **Masa** | powierzchnia, stopień pisma, krój, promień | **3 poziomy** (karty → kafle → listwa) | rozdzielanie ścieżek działania (W3) |
| **Kontrast** | `--dn-tekst` vs `--dn-tekst-2` vs `--dn-tekst-3`; obrys vs powierzchnia | **3 poziomy tekstu + 3 obrysów** | oddzielenie treści od metadanych |
| **Pozycja** | kolejność w pionie, przypięcie do krawędzi, warstwa z-index | **11 warstw** (żetonowe) | stałość ramy wobec zmienności treści |
| **Przestrzeń** | 12 px / 24 px / 48 px | **3 poziomy** | grupowanie i rozdzielanie bez linii |

### 10.2. Kontrast — trzy poziomy tekstu i ich prawo użycia

| Żeton | Motyw jasny | Motyw ciemny | Kontrast | Dozwolone użycie |
|---|---|---|---:|---|
| `--dn-tekst` | `szary-900` | `szary-100` | 16,1 : 1 / 15,0 : 1 | treść główna, tytuły, wartości |
| `--dn-tekst-2` | `szary-600` | `szary-400` | 6,4 : 1 / 6,8 : 1 | opisy, treść drugorzędna, etykiety kontrolek |
| `--dn-tekst-3` | `szary-500` | `szary-500` | poniżej 4,5 : 1 na tekst zwykły | **wyłącznie** metadane i tekst ≥ 18,66 px półgruby |

**Konsekwencja układowa:** `--dn-tekst-3` nie może nieść informacji krytycznej dla decyzji. Etykieta wersalikowa strefy, godzina wpisu w Chat Window, identyfikator przebiegu w Monitorze procesu — tak. Status zadania, nazwa modułu, wartość ustawienia — nie.

### 10.3. Pozycja — co jest stałe, a co zmienne

```
 ZAWSZE STAŁE (rama kokpitu — nie przewija się, nie zmienia przy zmianie modułu)
 ├─ pasek górny                              zawsze atramentowy w OBU motywach
 ├─ pas kart sesji
 ├─ boczna nawigacja / Panel orkiestracji
 └─ pas komunikacji (Chat Window)            rekonfiguruje kontekst, nie znika

 ZMIENNE (treść — wymienia się przy zmianie modułu, sekcji albo karty)
 ├─ okna operacyjne modułu
 ├─ zawartość sekcji Panelu orkiestracji
 └─ zawartość zakresu w nakładce konfiguracji

 WARSTWA TOWARZYSZĄCA (ponad wszystkim, nie należy do żadnej z grup)
 ├─ Always On Display                        z-index 1200
 ├─ toasty                                   z-index 1000
 └─ Centrum poleceń                          z-index 1300
```

Stałość ramy jest tu instrumentem hierarchii, nie tylko wygody: to, co się nie porusza przy zmianie kontekstu, Operator przestaje sprawdzać wzrokiem. Cała uwaga trafia do jednego pola — obszaru roboczego.

### 10.4. Przestrzeń zamiast linii

| Problem | Rozwiązanie linią (**odrzucone**) | Rozwiązanie przestrzenią (**przyjęte**) |
|---|---|---|
| Rozdzielenie trzech stref Centrum dowodzenia | trzy separatory poziome | odstęp 48 px między strefami |
| Grupowanie ustawień w zakresie konfiguracji | ramka wokół grupy | odstęp 24 px + etykieta wersalikowa |
| Oddzielenie metadanych wpisu od treści | linia pod nagłówkiem wpisu | odstęp 4 px + zmiana stopnia i barwy |
| Rozdzielenie kart sesji | pionowe kreski między kartami | odstęp 4 px + tło karty aktywnej |

Linia jest w tym systemie zarezerwowana dla **granic strukturalnych** — dolna krawędź paska górnego, prawa krawędź bocznej nawigacji, górna krawędź pasa komunikacji, dolna krawędź nagłówka okna operacyjnego. Wszystko inne rozdziela przestrzeń.

### 10.5. Jedna asymetria na widok

`WARIANCJA_PROJEKTOWA = 4/10` przekłada się na regułę operacyjną: **każdy widok ma prawo do jednej asymetrii, i musi ona nieść hierarchię.**

| Widok | Jedyna asymetria | Co niesie |
|---|---|---|
| W3 Centrum dowodzenia | malejąca masa trzech stref | „ta ścieżka jest ważniejsza od tamtej” |
| W4 Powłoka środowiska | boczna 224 px vs obszar `1fr` | „nawigacja jest narzędziem, praca jest celem” |
| W5 Widok operacyjny | pas komunikacji 54 % vs okno modułu `1fr` | „rozmowa jest stała, wynik jest zmienny” |
| W5 MultitaskingAI | plansza 58 % / 42 % | „obserwuję role, narzędzia towarzyszą” |
| W6 Nakładka dwudzielna | nawigacja 224 px vs zawartość `1fr` | „zakres jest wyborem, ustawienia są pracą” |

Drugiej asymetrii w tym samym widoku nie wprowadza się. Jeśli układ jej potrzebuje, to znaczy, że w widoku są dwa widoki.

---

## 11. Puste stany, stany ładowania i stany błędu na poziomie widoku

### 11.1. Trzy stany, jedna geometria

Puste stany, stany ładowania i stany błędu to **stany widoku, nie komponentu**. Ich wspólna zasada geometryczna: **stan zajmuje dokładnie tę samą powierzchnię co treść, którą zastępuje**, i jest w niej wyśrodkowany w obu osiach. Dzięki temu przejście stan → treść nie wywołuje przeskoku układu.

### 11.2. Puste stany — katalog wg widoku

| Widok / miejsce | Wyzwalacz | Ikona | Tytuł (dokładny) | Opis | Działanie |
|---|---|---|---|---|---|
| Obszar roboczy, nowa karta sesji | kliknięcie „+” | `karta-okna` | Nowa karta sesji — pusta | „Wybierz moduł z bocznej nawigacji, aby rozpocząć pracę” | wskazanie bocznej nawigacji |
| Obszar roboczy MultitaskingAI, nowy proces | kliknięcie „+” | `wezly` | Nowy proces orkiestracji | oczekiwanie na konfigurację zespołu i ról w Panelu orkiestracji | wskazanie Panelu orkiestracji |
| Chat Window, nowy moduł | pierwsze wejście do modułu | `rozmowa` | Historia rozmowy jest pusta | polecenie wpisane poniżej otworzy strumień odpowiedzi | fokus na polu poleceń |
| Sekcja „Zespoły” Panelu orkiestracji | brak zapisanych konfiguracji | `gwiazdka` | Brak zapisanych zespołów | zapisany zespół można wczytać do nowego procesu | przycisk zapisu bieżącej konfiguracji |
| Sekcja „Kolejki” | brak zdefiniowanych kolejek | `filtr` | Brak kolejek w tym procesie | kolejka powstaje przy pierwszym `enqueue` | odesłanie do Queue Manager |
| Model Panels (Roundtable) | brak uczestników | `debata` | Brak uczestników debaty | dodanie uczestnika otwiera pierwszy panel | przycisk „+ Dodaj uczestnika” |
| Translation Panels | brak wybranych języków | `tlumacz` | Brak paneli tłumaczenia | wybór języka otwiera panel równoległy | przycisk „+ Dodaj język” |
| Wyniki wyszukiwania w pasku górnym | brak trafień | `szukaj` | Brak trafień dla tej frazy | podpowiedź zmiany zakresu wyszukiwania | wyczyszczenie frazy |

**Struktura pustego stanu (`.dn-pusty-stan`):** ikona 28 px → tytuł (`--dn-fs-lg`, krój nagłówkowy, `--dn-tekst-2`) → opis (`--dn-fs-sm`, `--dn-tekst-3`, maks. 40 znaków w wierszu) → opcjonalne działanie. Odstępy: 8 px między elementami, 40 px odstępu wewnętrznego kontenera.

**Zakaz:** pusty stan nigdy nie jest ilustracją dekoracyjną ani żartem. Jest instrukcją następnego kroku.

### 11.3. Stany ładowania — trzy poziomy

| Poziom | Zasięg | Forma | Czas typowy | Element |
|---|---|---|---|---|
| **L1 — mikro** | pojedyncza kontrolka | `.dn-spinner` 14 px w miejscu etykiety przycisku; przycisk pozostaje klikalny | do 1 s | przycisk „Zaloguj”, „Wyślij”, „Zapisz agenta” |
| **L2 — panel** | jedno okno operacyjne | nakładka na ciało okna: spinner + etykieta w krojach mono; nagłówek okna pozostaje czytelny | 1–5 s | przeładowanie obszaru roboczego przy zmianie modułu |
| **L3 — widok** | cały widok | widok przejściowy (W1): godło + spinner + komunikat stanu | zmienny | okno startowe, ponowne łączenie po utracie sieci |

**Zasada nienaruszalna:** żaden stan ładowania nie odbiera klikalności. Przycisk w stanie `aria-busy="true"` zmienia kursor na `progress` i pokazuje wskaźnik, ale pozostaje aktywny. Jest to wprost konsekwencja ADL-017.

**Zasada geometrii:** stan L2 nie zmienia wysokości okna. Nakładka jest pozycjonowana absolutnie w ciele okna, więc po zakończeniu ładowania treść wchodzi w miejsce już zarezerwowane. Przejście: `opacity` + `translateY(4px)`, czas `--dn-czas-3` (0,22 s).

### 11.4. Stany błędu — trzy poziomy

| Poziom | Zasięg | Forma | Przykład |
|---|---|---|---|
| **B1 — pole** | pojedyncza kontrolka | `.dn-pole-blad` pod polem: ikona `blad` + komunikat, obrys pola w `--dn-blad-tekst` | „Nieprawidłowy login lub hasło. Spróbuj ponownie.” |
| **B2 — panel** | jedno okno operacyjne | pas komunikatu w ciele okna: `--dn-blad-tlo`, obrys lewy 2 px, ikona + treść + działanie naprawcze | „Kolejka wstrzymana — zadanie #130 oczekuje na rozstrzygnięcie” (przykładowe) |
| **B3 — widok** | cały widok | widok przejściowy w wariancie błędu: ikona `ostrzezenie` + komunikat + kontrolka „Spróbuj ponownie” | brak połączenia z serwerem w oknie startowym |

**Zasada:** stan błędu **nigdy nie usuwa treści, którą Operator już wprowadził**. Formularz logowania po błędzie pozostaje wypełniony. Pole promptu po nieudanym wysłaniu zachowuje treść. Kolejka po niepowodzeniu etapu zachowuje pozostałe etapy.

**Zasada barwy:** stan błędu ma zawsze trzy nośniki — barwę `--dn-blad-*`, ikonę `blad` albo `ostrzezenie` oraz etykietę tekstową. Nigdy sam kolor.

### 11.5. Diagram — cykl stanów widoku operacyjnego

```
        ┌──────────────┐  wybór modułu   ┌──────────────┐  odpowiedź   ┌──────────────┐
        │  PUSTY       │ ───────────────►│  ŁADOWANIE   │ ────────────►│  WYPEŁNIONY  │
        │ .dn-pusty-   │                 │  L2 (panel)  │              │  treść okna  │
        │  stan        │                 │              │              │              │
        └──────────────┘                 └──────┬───────┘              └──────┬───────┘
               ▲                                 │ niepowodzenie              │ zmiana modułu
               │ zamknięcie / nowa karta         ▼                            │
               │                          ┌──────────────┐                    │
               └──────────────────────────│  BŁĄD  B2    │◄───────────────────┘
                                          │ + działanie  │
                                          │   naprawcze  │
                                          └──────────────┘
   Wszystkie cztery stany zajmują tę samą powierzchnię — przejście nie przesuwa układu.
```

---

## 12. Gęstość informacji — pojemność widoku na progu

### 12.1. Założenia rachunku

Rachunek pojemności jest wyprowadzony wyłącznie z żetonów wymiarów. Przyjęte wysokości okna klienckiego dla każdego progu (wartości typowe dla proporcji ekranu w danej klasie urządzenia):

| Próg | Szerokość | Przyjęta wysokość | Uzasadnienie |
|---|---:|---:|---|
| `w1` | 640 px | 360 px | telefon w orientacji poziomej |
| `w2` | 960 px | 640 px | tablet |
| `w3` | 1280 px | 800 px | monitor biurkowy |
| `w4` | 1600 px | 900 px | szerokie biurko |

Wysokość dostępna dla obszaru roboczego = wysokość okna − 48 (pasek górny) − 36 (pas kart sesji).

### 12.2. Podział wysokości obszaru roboczego (gęstość zwarta)

| Próg | Obszar roboczy | Pas komunikacji `minmax(320, 54%)` | Okna modułu `minmax(0,1fr)` | Wzorzec czynny |
|---|---:|---:|---:|---|
| `w4` | 816 px | **441 px** (54 %) | 375 px | B1 (lub B4 — dwie kolumny) |
| `w3` | 716 px | **387 px** (54 %) | 329 px | B1 |
| `w2` | 556 px | **320 px** (limit dolny) | 236 px | B1 na granicy |
| `w1` | 276 px | — | — | **B3 — zakładki** (próg 628 px niespełniony) |

### 12.3. Pojemność w wierszach tabeli / listy

Rachunek: `wiersze = floor((wysokość_okna_modułu − 12 margines − 36 nagłówek okna − 32 nagłówek tabeli) / 36)`.

| Próg | Gęstość zwarta (wiersz 36 px) | Gęstość przestronna (wiersz 44 px) | Tryb dotykowy (wiersz 44 px) |
|---|:-:|:-:|:-:|
| `w4` | **8 wierszy** | 6 wierszy | 6 wierszy |
| `w3` | **6 wierszy** | 5 wierszy | 5 wierszy |
| `w2` | **4 wiersze** | 3 wiersze | 3 wiersze |
| `w1` | 2 wiersze (w zakładce) | 1 wiersz | 1 wiersz |

**Wniosek:** przejście z gęstości zwartej na przestronną kosztuje **dwa wiersze na progu `w3`** — jedną trzecią pojemności okna. To jest cena, którą Operator płaci świadomie, i uzasadnienie, dla którego gęstość zwarta jest domyślna.

### 12.4. Pojemność w kartach i kaflach (Centrum dowodzenia)

| Element | `w4` | `w3` | `w2` | `w1` |
|---|:-:|:-:|:-:|:-:|
| Karty środowisk w rzędzie | 4 | 4 | 2 | 1 |
| Rzędów kart środowisk | 1 | 1 | 2 | 4 |
| Kafle komponentów w rzędzie | 4 | 4 | 2 | 1 |
| Pozycje listwy w wierszu | 3 | 3 | 3 | 1 (zawijanie) |
| Szerokość karty środowiska | 282 px | 282 px | 444 px | 608 px |
| Szerokość kafla | 282 px | 282 px | 444 px | 608 px |

Szerokość karty przy rozpiętości 3/12: `3 × 78 + 2 × 24 = 282 px`. Przy rozpiętości 6/12 na progu `w2`: `6 × 54 + 5 × 24 = 444 px`.

### 12.5. Pojemność w panelach równoległych

| Próg | Model Panels / Translation Panels bez przewijania | Szerokość panelu |
|---|:-:|---:|
| `w4` | **5** | 260 px |
| `w3` | **4** | 246 px |
| `w2` | **3** | 285 px (boczna zwinięta) |
| `w1` | **2** | 298 px (boczna ukryta) |

### 12.6. Pojemność pasa kart sesji

Karta sesji ma szerokość zmienną: minimum 120 px, maksimum 200 px, tytuł skracany wielokropkiem. Kontrolka „+” zajmuje 32 px, odstępy po 4 px.

| Próg | Szerokość pasa | Kart przy 180 px | Kart przy 120 px (skrócone) | Powyżej pojemności |
|---|---:|:-:|:-:|---|
| `w4` | 1568 px | **8** | 12 | przewijanie poziome |
| `w3` | 1248 px | **6** | 10 | przewijanie poziome |
| `w2` | 928 px | **5** | 7 | przewijanie poziome |
| `w1` | 608 px | 3 | 4 | lista rozwijana „karty sesji (n)” |

### 12.7. Pojemność Chat Window w wpisach

Wpis `.dn-wpis` ma wysokość zmienną. Przyjmując wpis typowy o trzech wierszach treści: medalion + nadawca + godzina (20 px) + 3 × 19 px treści + odstępy 16 px ≈ **93 px**. Pole promptu z przybornikiem zajmuje 96 px.

| Próg | Pas komunikacji | Wysokość historii | Wpisów widocznych |
|---|---:|---:|:-:|
| `w4` | 441 px | 345 px | **3** |
| `w3` | 387 px | 291 px | **3** |
| `w2` | 320 px | 224 px | **2** |

**Wniosek projektowy:** Chat Window pokazuje trzy wpisy jednocześnie na progu odniesienia. Oznacza to, że **strumień odpowiedzi modelu musi automatycznie przewijać historię do dołu**, bo w przeciwnym razie Operator nie zobaczy powstającej odpowiedzi. Automatyczne przewijanie zatrzymuje się w chwili, gdy Operator sam przewinie w górę.

---

## 13. Decyzje projektowe

Poniższe rozstrzygnięcia wykraczają poza literalny zapis dokumentacji źródłowej. Każde jest wyprowadzone z katalogu komponentów, z żetonów albo z kontraktu kierunku — i każde jest tu jawnie odnotowane.

| # | Rozstrzygnięcie | Co mówiło źródło | Decyzja i uzasadnienie |
|---|---|---|---|
| **D-01** | **Taksonomia siedmiu rodzajów widoku** | Dokumentacja wymienia okna (KANON rozdz. 7), nie klasyfikuje ich wg praw układu | Wprowadzono siedem rodzajów (W1–W7) jako warstwę porządkującą **nad** istniejącym inwentarzem okien. Żadne nowe okno nie powstało — każdy rodzaj wskazuje okna już zinwentaryzowane. |
| **D-02** | **Granica obowiązywania siatki 12-kolumnowej** | KANON podaje siatkę (12 / 24 px / 1200 px), nie wskazuje, gdzie obowiązuje | Rozstrzygnięto: siatka obowiązuje w widokach czytania i wyboru (W1–W3, W6), nie obowiązuje w widokach pracy (W4, W5). Uzasadnienie: wymiary ramy kokpitu pochodzą z żetonów (48 · 36 · 224), nie z podziału szerokości — mieszanie dwóch reżimów wymiarowania w jednym widoku produkuje wartości spoza skali 4 px. |
| **D-03** | **Szerokość ramy widoku formularza = 560 px** | Źródło mówi „duży panel/karta `.dn-modal`-podobny, wyśrodkowany”, bez wartości | Przyjęto `--dn-wym-modal` (560 px), bo źródło samo odsyła do modala. Zero nowych wartości. |
| **D-04** | **Szerokość nawigacji zakresów w nakładce konfiguracji = 224 px** | Źródło mówi „dwudzielny układ (nawigacja + zawartość)”, bez wymiaru | Przyjęto `--dn-wym-boczna` (224 px) — tę samą miarę co boczna nawigacja środowiska. Operator uczy się jednej szerokości nawigacji, nie dwóch. |
| **D-05** | **Zwinięcie bocznej nawigacji na progu `w2` = 56 px** | KANON: „w2 960 (tablet — boczna zwija się do ikon)”, bez wartości | 56 px = przycisk ikonowy 32 px + 2 × 12 px odstępu panelu. Wartość w całości wyprowadzona z żetonów. |
| **D-06** | **Liczba kolumn siatki jest niezmienna na wszystkich progach** | Źródła nie rozstrzygają zachowania siatki poniżej `w3` | Rozstrzygnięto: dwanaście kolumn i rynna 24 px obowiązują zawsze; zmienia się wyłącznie rozpiętość elementu (3/12 → 6/12 → 12/12). Jeden zestaw linii pomocniczych zamiast czterech. |
| **D-07** | **Margines zewnętrzny 16 px na progach `w1`–`w2`** | Brak w dokumencie | Przyjęto `--dn-od-4` (16 px) zamiast `--dn-od-6` (24 px). Wartość ze skali, uzasadniona odzyskaniem 16 px szerokości użytecznej na wąskim ekranie. |
| **D-08** | **Odstęp międzystrefowy Centrum dowodzenia = 48 px** | Źródło mówi o „malejącej masie” stref, nie podaje odstępu | Przyjęto `--dn-od-12` (48 px) — podwojony `--dn-odstep-sekcji`. Strefy są trzema niezależnymi ścieżkami, nie trzema akapitami; podwojony odstęp komunikuje niezależność bez linii. |
| **D-09** | **Szerokość minimalna panelu w siatce równoległej = 240 px** | Źródła (Roundtable, Translate) opisują „siatkę równoległą”, bez wymiaru | 240 px to szerokość, przy której nagłówek panelu (kod języka + status + selektor tonu) mieści się w jednym wierszu. Wartość wyprowadzona z realnej zawartości nagłówka opisanej w dokumentacji modułu Translate. |
| **D-10** | **Próg przejścia B1 → B3 = 628 px wysokości okna** | Brak w dokumencie | Wyprowadzone z żetonów: 48 + 36 + 24 + 320 + 200. Przejście wywołuje **wysokość**, nie szerokość — bo to wysokość decyduje o użyteczności pasa komunikacji. |
| **D-11** | **Automatyczne przewijanie historii Chat Window** | Źródło mówi o strumieniu odpowiedzi na żywo, nie rozstrzyga przewijania | Rozstrzygnięto: historia przewija się automatycznie do dołu w czasie strumienia i zatrzymuje przewijanie, gdy Operator sam przewinie w górę. Wynika z rachunku pojemności (rozdz. 12.7): przy trzech widocznych wpisach bez auto-przewijania odpowiedź powstaje poza polem widzenia. |
| **D-12** | **Katalog pustych stanów per widok** | Źródła opisują `.dn-pusty-stan` i jeden przykład („Wybierz moduł z bocznej nawigacji…”) | Rozszerzono na osiem miejsc występowania. Każdy tytuł i opis zbudowany wyłącznie z pojęć istniejących w dokumentacji (moduł, karta sesji, kolejka, uczestnik debaty, panel tłumaczenia). Zero nowych funkcji. |
| **D-13** | **Trzy poziomy stanu ładowania (L1/L2/L3) i błędu (B1/B2/B3)** | Źródła opisują `.dn-spinner`, `.dn-pole-blad`, komunikat błędu połączenia — bez klasyfikacji zasięgu | Wprowadzono klasyfikację zasięgu, bo bez niej ta sama sytuacja bywa projektowana raz jako toast, raz jako nakładka. Formy pozostały te z biblioteki komponentów. |
| **D-14** | **Rozbieżność źródeł: wymiary z pakietu v1.0** | `SIATKA-I-UKLAD.md` (v1.0) podaje pasek górny 56 px, przycisk 36 px, promienie 4/8/10/14/20 px, cień złoty | Rozstrzygnięto na rzecz **KANON i `zetony.css` v2.0**: pasek 48 px, kontrolka 32 px, promienie 3/6/8/10/14 px, brak cienia złotego. Pakiet v1.0 opisuje zastąpioną warstwę wizualną (granat + złoto), o czym KIERUNEK.md mówi wprost. Wartości v1.0 nie są w tym opracowaniu używane. |
| **D-15** | **Jedna asymetria na widok** | KIERUNEK.md: `WARIANCJA_PROJEKTOWA = 4/10`, „asymetria tylko gdy niesie hierarchię” | Przełożono pokrętło na regułę operacyjną: dokładnie jedna asymetria na widok, wskazana w tabeli rozdz. 10.5. Druga asymetria oznacza, że w widoku są dwa widoki. |
| **D-16** | **Zniknięcie elementu wymaga alternatywnego wywołania** | ADL-017 dotyczy blokad funkcjonalnych, nie układu | Rozszerzono zasadę zero blokad na warstwę układu: element może zniknąć z widoku (próg `w1`), nie może zniknąć z systemu — musi pozostać wywoływalny (`menu`, lista rozwijana, zakładka). |
| **D-17** | **Wyrównanie segmentów w Translation Panels** | Źródło: „kliknięcie segmentu przewija wszystkie Translation Panels do odpowiadającego segmentu” | Z wymagania funkcjonalnego wyprowadzono wymóg układowy: wspólna wysokość wiersza segmentu w całej siatce paneli (siatka na poziomie kontenera paneli, nie pojedynczego panelu). |
| **D-18** | **Wartości pojemności w rozdz. 12** | Brak w dokumencie | Wszystkie liczby są **rachunkiem** z żetonów przy jawnie podanych założeniach wysokości okna (rozdz. 12.1), nie pomiarem ani szacunkiem. Zmiana żetonu zmienia wynik — dlatego podano wzory, nie same liczby. |

---

## Załącznik A — katalog zbiorczy wzorców układu

| Kod | Wzorzec | Zapis CSS (rdzeń) | Gdzie obowiązuje | Rozdział |
|---|---|---|---|:-:|
| **A** | Powłoka środowiska | `grid-template-rows: 48px 36px 1fr` + `grid-template-columns: 224px 1fr` | 4 powłoki środowisk | 4 |
| **B1** | Obszar roboczy, podział pionowy | `grid-template-rows: minmax(0,1fr) minmax(320px,54%)` | domyślny dla 15 modułów | 5 |
| **B2** | Obszar roboczy z pasmem narzędzi | okno modułu: `grid-template-columns: 1fr 280px` | Studio (Tools Panel), Research (Sources Manager) | 5.4 |
| **B3** | Obszar roboczy zakładkowy | `.dn-zakladki` zamiast podziału | `w2` i niżej; okno < 628 px wysokości | 5.4 · 9.4 |
| **B4** | Dwa okna komunikacji | pas komunikacji: `grid-template-columns: 1fr 1fr` | próg `w4` | 5.4 · 9.2 |
| **C** | Trzy strefy malejącej masy | siatka 12 kol.: 4 × 3/12 → 4 × 3/12 → 1 × 12/12 | Centrum dowodzenia | 6 |
| **D** | Asymetria koordynator–wykonawca | `grid-template-rows: minmax(0,58%) minmax(0,42%)`; para `1fr 1fr`; pas dolny `1.25fr 1fr` | MultitaskingAI, okna robocze ról | 7 |
| **E** | Panele wielokolumnowe | `repeat(N, minmax(240px, 1fr))` | Roundtable (Model Panels), Translate (Translation Panels) | 8 |
| **F** | Nakładka dwudzielna | `grid-template-columns: 224px 1fr` w `.dn-modal` do 1200 px | Okno Konfiguracji, Okno Ustawień | 1.3 · 9.2 |
| **G** | Blok centrowany | `place-content: center`; szerokość 560 px (formularz) albo `auto` (widok przejściowy) | Okno startowe, Okno rejestracji i logowania | 1.3 |

---

## Załącznik B — mapa źródeł

| Rozdział opracowania | Źródło wiążące |
|---|---|
| 1. Taksonomia widoków | KANON rozdz. 6–7 · `dok/projekt-ui/przeplyw/przeplyw-okien.md` · `brief/INWENTARZ-OKIEN.md` rozdz. 1 |
| 2. Siatka | KANON rozdz. 2 · `zasoby/zetony/zetony.css` §7 · `system-wizualny/pakiet/SIATKA-I-UKLAD.md` (z zastrzeżeniem D-14) |
| 3. Rytm pionowy | KANON rozdz. 2 · `zetony.css` §4, §6 |
| 4. Powłoka środowiska | `dok/projekt-ui/przeplyw/elementy-okien.md` rozdz. 3.1–3.4, 4.1–4.4 · `design/06-okna/e2-srodowisko.html` |
| 5. Obszar roboczy z pasem komunikacji | `elementy-okien.md` rozdz. 3.5, 3.5.1 · `design/06-okna/e2-srodowisko.html` |
| 6. Trzy strefy | `elementy-okien.md` rozdz. 2.2–2.5 · `INWENTARZ-OKIEN.md` rozdz. 4 · `SIATKA-I-UKLAD.md` rozdz. 8 |
| 7. Asymetria koordynator–wykonawca | `dok/projekt-ui/srodowiska/multitaskingai.md` rozdz. 3 · `design/06-okna/e4-koordynator-wykonawca.html` · KIERUNEK.md rozdz. 4 |
| 8. Panele wielokolumnowe | `dok/projekt-ui/moduly/roundtable.md` rozdz. 2, 3.2 · `dok/projekt-ui/moduly/translate.md` rozdz. 2, 3.3 |
| 9. Zachowanie responsywne | KANON rozdz. 2 (punkty łamania) · `zetony.css` §7, §12 |
| 10. Hierarchia wizualna | KIERUNEK.md rozdz. 2–4 · KANON rozdz. 1, 9 · `zasoby/zetony/kontrasty.json` |
| 11. Puste stany, ładowanie, błąd | `elementy-okien.md` rozdz. 3.5 (Tabela 12) · `INWENTARZ-OKIEN.md` rozdz. 2.3, 5.5 · `zasoby/css/komponenty.css` (`.dn-pusty-stan`, `.dn-spinner`, `.dn-pole-blad`) |
| 12. Gęstość informacji | rachunek z `zetony.css` §6–§7 przy założeniach rozdz. 12.1 |
| 13. Decyzje projektowe | niniejsze opracowanie |

---

## Załącznik C — lista sprawdzeń układu przed oddaniem widoku

- [ ] Rodzaj widoku (W1–W7) rozpoznany i zapisany w komentarzu na początku pliku
- [ ] Reżim wymiarowania zgodny z rodzajem widoku (siatka vs żeton) — bez mieszania
- [ ] Wszystkie odstępy pochodzą ze skali `0 · 4 · 8 · 12 · 16 · 20 · 24 · 32 · 40 · 48 · 64`
- [ ] Wszystkie wymiary ramy pochodzą z żetonów `--dn-wym-*`
- [ ] Zero wartości szesnastkowych wpisanych wprost — wyłącznie `var(--dn-*)`
- [ ] Każdy kontener siatki z przewijanym dzieckiem ma `min-height: 0` (i `min-width: 0` przy kolumnach)
- [ ] Pas komunikacji ma `minmax(320px, 54%)` albo jawnie udokumentowane odstępstwo
- [ ] Widok ma dokładnie **jedną** asymetrię i niesie ona hierarchię (rozdz. 10.5)
- [ ] Zachowanie sprawdzone na czterech progach: `w1` 640 · `w2` 960 · `w3` 1280 · `w4` 1600
- [ ] Żaden element, który znika na progu `w1`, nie traci alternatywnego wywołania
- [ ] Puste stany, stany ładowania i stany błędu zajmują tę samą powierzchnię co treść
- [ ] Żaden stan ładowania nie odbiera klikalności (`disabled` nie występuje w pliku)
- [ ] Stan komunikowany co najmniej dwoma nośnikami (barwa + ikona albo etykieta)
- [ ] Pasek górny atramentowy w obu motywach; oba motywy sprawdzone
- [ ] Pierścień fokusu widoczny na każdej kontrolce (2 px + odsunięcie 2 px)
- [ ] Warstwy z-index wyłącznie z żetonów `--dn-z-*`
- [ ] Gęstość przestronna (`data-gestosc="przestronna"`) sprawdzona — układ nie pęka
- [ ] Tryb dotykowy (`pointer: coarse`) sprawdzony — cele 40/44 px mieszczą się w rytmie
- [ ] Polskie diakrytyki poprawne w całym pliku

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
