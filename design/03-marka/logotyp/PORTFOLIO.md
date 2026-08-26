# Portfolio logotypu — pełne repozytorium znaku Danaco Console

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment |
| **Opracowanie** | Pełne portfolio logotypu — wszystkie konfiguracje i formaty |
| **Wersja warstwy wizualnej** | v2.0 |
| **Status** | Deweloperski · geometria znaku zatwierdzona |
| **Data** | 2026-08-14 |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Odbiorcy** | projektanci, wdrożeniowcy, dostawcy druku i grawerowania, zespół produktu |
| **Zakres** | 102 plików: 37 × SVG · 62 × PNG · 3 × ICO |
| **Katalog** | `WYNIK/03-marka/logotyp/` |
| **Galeria** | `WYNIK/03-marka/portfolio-logotypu.html` |
| **Źródło prawdy** | kontrakt systemu projektowego · `kierunek systemu projektowego |

> **Zasada nadrzędna tego repozytorium.** Krzywe sygnetu i wektory logotypu są
> **zatwierdzone i nietykalne**. Wolno budować konfiguracje, warianty barwne,
> kadry i formaty — **nie wolno zmieniać ani jednej współrzędnej znaku**.
> Każdy plik w tym katalogu powstał przez złożenie tych samych, niezmienionych
> ścieżek; żaden nie został przerysowany.

---

## Spis treści

1. [Przeznaczenie i budowa repozytorium](#1-przeznaczenie-i-budowa-repozytorium)
2. [Znak — fakty geometryczne](#2-znak--fakty-geometryczne)
3. [Mapa konfiguracji A–H](#3-mapa-konfiguracji-ah)
4. [Katalog SVG — pełna tabela](#4-katalog-svg--pełna-tabela)
5. [Katalog PNG — pełna tabela](#5-katalog-png--pełna-tabela)
6. [Katalog ICO](#6-katalog-ico)
7. [Dobór wariantu — drzewo decyzyjne](#7-dobór-wariantu--drzewo-decyzyjne)
8. [Pole ochronne i rozmiary minimalne](#8-pole-ochronne-i-rozmiary-minimalne)
9. [Barwy znaku i ich żetony](#9-barwy-znaku-i-ich-żetony)
10. [Matryca nośników](#10-matryca-nośników)
11. [Czego nie wolno robić ze znakiem](#11-czego-nie-wolno-robić-ze-znakiem)
12. [Technika wytworzenia i odtworzenia](#12-technika-wytworzenia-i-odtworzenia)
13. [Kontrola jakości](#13-kontrola-jakości)
14. [Decyzje projektowe](#14-decyzje-projektowe)

---

## 1. Przeznaczenie i budowa repozytorium

### 1.1 Po co to jest

Repozytorium odpowiada na jedno pytanie zadawane setki razy w cyklu życia
produktu: **„którego pliku znaku mam użyć tutaj?"**. Zamiast odsyłać do
jednego pliku i liczyć na zdrowy rozsądek odbiorcy, katalog dostarcza
gotową konfigurację dla każdego przewidzianego nośnika — od favicona 16 px
po matrycę do grawerowania.

Repozytorium jest **kompletne i zamknięte**: jeśli nośnik nie ma tu swojego
pliku, znaczy to, że nośnik wymaga decyzji projektowej, a nie improwizacji
na miejscu.

### 1.2 Struktura katalogu

```
WYNIK/03-marka/logotyp/
├── PORTFOLIO.md          ← ten dokument: spis całego repozytorium
├── svg/                  ← 37 plików · wektor, źródło każdego formatu rastrowego
│   ├── sygnet*.svg              A · znak samodzielny (13 konfiguracji)
│   ├── logotyp*.svg             B · sama typografia (4 konfiguracje)
│   ├── logo-poziomy*.svg        C+F+G+H · lockup poziomy i pochodne
│   ├── logo-pionowy*.svg        D+F · lockup pionowy i pochodne
│   └── logo-kompaktowy*.svg     E · sygnet + DANACO
├── png/                  ← 62 plików · rasteryzacje, tło przezroczyste
└── ico/                  ←  3 pliki  · wielorozmiarowe ikony Windows/przeglądarki
```

### 1.3 Reguła nazewnictwa

```
                <konfiguracja> - <wariant> - <rozmiar> . <format>
                      │             │            │         │
   sygnet ────────────┘             │            │         └── svg | png | ico
   logotyp                          │            │
   logo-poziomy                     │            └── szerokość w pikselach (PNG)
   logo-pionowy                     │
   logo-kompaktowy                  └── (brak)          = wariant podstawowy, jasne tło
                                       na-ciemnym       = odwrócenie na powierzchnię atramentową
                                       mono-bialy       = jedna barwa, biel
                                       mono-czarny      = jedna barwa, atrament
                                       uproszczony      = jeden grot (poniżej 24 px)
                                       w-kole           = medalion okrągły
                                       w-kwadracie      = kafel zaokrąglony
                                       deskryptor       = z linią „AI Operating Environment"
                                       do-grawerowania  = jednobarwny, kropka konturem
                                       znak-wodny       = obniżone krycie
```

Nazwy plików są **kebab-case po polsku**. Wyjątkiem są `favicon.ico`
i `ikona-aplikacji.ico`, gdzie nazwa `favicon` jest wymuszona konwencją
przeglądarek.

### 1.4 Pochodzenie plików

| Pochodzenie | Liczba | Znaczenie |
|---|---|---|
| **źródłowy** | 14 | plik przeniesiony **bajt w bajt** z `zasoby/marka/logo/` — wersja zatwierdzona, wymieniona w kontrakcie systemu projektowego |
| **wytworzony** | 23 | konfiguracja złożona w tym opracowaniu z **niezmienionych** ścieżek znaku |

---

## 2. Znak — fakty geometryczne

### 2.1 Sygnet „Delegacja" — siatka 96 × 96

Znak czyta się **„»»."** i opowiada pętlę produktu: zlecenie przechodzi od
koordynatora do wykonawcy, kropka to praca, która właśnie ruszyła.

| Element | Definicja | Uwaga |
|---|---|---|
| **Grot 1** (koordynator) | `M12 26 H24 L44 48 L24 70 H12 L32 48 Z` | szerokość ramienia 12, wierzchołek na osi y = 48 |
| **Grot 2** (wykonawca) | `M40 26 H52 L72 48 L52 70 H40 L60 48 Z` | ten sam kształt przesunięty o 28 w osi x |
| **Kropka sygnału** | `cx = 83 · cy = 63,5 · r = 6,5` | jedyny element w barwie sygnałowej |
| **Pole farby** | x ∈ ⟨12 ; 89,5⟩ · y ∈ ⟨26 ; 70⟩ | 77,5 × 44 — proporcja 1,761 : 1 |
| **Oś optyczna** | x = 50,75 · y = 48 | punkt wyśrodkowania w medalionach |

```
  siatka 96 × 96                        oś y = 48
  ┌───────────────────────────────────────────────────┐
  │                                                   │
  │   ▄▄▄▄▖        ▄▄▄▄▖                              │ 26
  │   █   ▀▚▖      █   ▀▚▖                            │
  │   █     ▚▖     █     ▚▖                           │
  │   █      ▚▖    █      ▚▖                          │
  │   ▀▀▀▀▘   ▚    ▀▀▀▀▘   ▚          ●               │ 48  ← kropka cy 63,5
  │   █      ▞     █      ▞          ▀▀               │
  │   █    ▗▞      █    ▗▞                            │
  │   █  ▗▞        █  ▗▞                              │
  │   ▀▀▀▘         ▀▀▀▘                               │ 70
  │                                                   │
  └───────────────────────────────────────────────────┘
   12   24    44   40  52    72       76,5  89,5
   grot 1 „koordynator"  grot 2 „wykonawca"   kropka „praca ruszyła"
```

### 2.2 Wariant uproszczony — obowiązuje poniżej 24 px

| Element | Definicja |
|---|---|
| **Grot** | `M18 18 H35 L64 48 L35 78 H18 L45 48 Z` |
| **Kropka** | `cx = 79 · cy = 69 · r = 9` |
| **Pole farby** | x ∈ ⟨18 ; 88⟩ · y ∈ ⟨18 ; 78⟩ |

Powód istnienia wariantu: przy 16 px odstęp między grotami pełnego znaku
schodzi poniżej jednego piksela i znak zlewa się w plamę. Jeden grot
z powiększoną kropką zachowuje czytelność „grot → praca" do 16 px włącznie.

### 2.3 Logotyp — blok typograficzny

| Linia | Krój | Stopień | Linia bazowa | Wersalik | Szerokość farby |
|---|---|---|---|---|---|
| `DANACO` | Space Grotesk 700 | 34 px (skala 0,034) | y = 34 | 23,80 | 132,19 |
| `CONSOLE` | IBM Plex Mono 500, rozstrzelone wersaliki | 12,5 px (skala 0,0125) | y = 50 | 8,73 | 73,31 |
| kropka po `CONSOLE` | — | `r = 3,4` | cy = 46 | — | 6,8 |
| deskryptor `AI Operating Environment` | IBM Plex Sans 400 | 10,5 px (skala 0,0105) | wg konfiguracji | 7,33 | 121,91 |

Odstęp linii bazowych `DANACO` → `CONSOLE` wynosi **16** w układzie poziomym
i logotypie samodzielnym, **18** w układzie pionowym.
Rozstrzelenie `CONSOLE`: krok liter 11,25 przy szerokości znaku 7,5 —
światło międzyliterowe 3,75, czyli **0,3 em**.

**Para krojów opowiada produkt: marka + maszyna.**

### 2.4 Kropka sygnału

Kropka jest **elementem sygnaturowym całego systemu**, nie ozdobą znaku.
Ta sama forma wraca w emblematach czterech środowisk, na kartach sesji,
w oknie komunikacji i w rdzeniu awatara Always On Display.
W znaku kropka **nie animuje się** — tętno 2,4 s należy do interfejsu,
nie do godła.

---

## 3. Mapa konfiguracji A–H

| Grupa | Konfiguracja | Plików SVG | Co zawiera | Kiedy używać |
|---|---|:-:|---|---|
| **A** | Sygnet — znak samodzielny | 12 | groty + kropka, bez typografii | znak działa bez nazwy: ikona, awatar, favicon, sygnatura |
| **B** | Logotyp — sama typografia | 4 | `DANACO` + `CONSOLE.` | nazwa działa bez znaku: stopka, tekst prawny, sygnatura dokumentu |
| **C** | Lockup poziomy | 4 | sygnet + logotyp obok | **układ podstawowy** — nagłówki, nośniki szerokie |
| **D** | Lockup pionowy | 4 | sygnet nad logotypem | nośniki wąskie i kwadratowe, plansze tytułowe |
| **E** | Lockup kompaktowy | 2 | sygnet + `DANACO` (bez `CONSOLE`) | pasek górny, stopka, wysokość poniżej 40 px |
| **F** | Znak z deskryptorem | 4 | lockup + `AI Operating Environment` | pierwszy kontakt z produktem — okno startowe, wizytówka |
| **G** | Tłoczenie i grawer | 3 | jedna barwa, kropka jako pierścień | grawer, tłoczenie, haft, pieczęć, matryca |
| **H** | Znak wodny | 4 | mono, krycie 8 % / 10 % | podkład dokumentu, papier firmowy, tło planszy |

```
                          ZNAK DANACO CONSOLE
                                   │
        ┌──────────────────────────┼──────────────────────────┐
        ▼                          ▼                          ▼
   A · SYGNET                B · LOGOTYP              A+B · LOCKUPY
   »». samodzielnie          DANACO CONSOLE.        ┌────┬────┬────┐
        │                          │                ▼    ▼    ▼    ▼
   ┌────┼────┬─────┐               │               C     D    E    F
   ▼    ▼    ▼     ▼               │            poziomy pion. komp. desk.
 pełny upr. medalion kafel         │
        │                          │
        └──────────┬───────────────┘
                   ▼
          G · TŁOCZENIE      H · ZNAK WODNY
          jedna barwa,       mono, krycie 8 %
          kropka konturem
```

---

## 4. Katalog SVG — pełna tabela

Wszystkie pliki wektorowe. Kolumna **Pochodzenie**: `źródłowy` = przeniesiony
bez zmian z `zasoby/marka/logo/`, `wytworzony` = złożony w tym opracowaniu
z niezmienionych ścieżek.

### 4.1 A · Sygnet — znak samodzielny

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `sygnet.svg` | 96×96 | jasne | źródłowy | 317 B | pasek górny, karta środowiska, nagłówek dokumentu |
| `sygnet-na-ciemnym.svg` | 96×96 | ciemne | źródłowy | 317 B | pasek górny motywu ciemnego, plansza atramentowa |
| `sygnet-mono-bialy.svg` | 96×96 | ciemne | źródłowy | 317 B | nadruk 1/0 na ciemnym, tłoczenie białe, materiały jednobarwne |
| `sygnet-mono-czarny.svg` | 96×96 | jasne | źródłowy | 317 B | stempel, faks, dokument czarno-biały, znak w tekście |
| `sygnet-uproszczony.svg` | 96×96 | jasne | źródłowy | 245 B | favicon, awatar 16–24 px, plakietka listy, ikona pliku |
| `sygnet-uproszczony-na-ciemnym.svg` | 96×96 | ciemne | źródłowy | 245 B | favicon motywu ciemnego, awatar na pasku |
| `sygnet-uproszczony-mono-bialy.svg` | 96×96 | ciemne | wytworzony | 245 B | tłoczenie i nadruk 1/0 w małej skali |
| `sygnet-uproszczony-mono-czarny.svg` | 96×96 | jasne | wytworzony | 245 B | stempel małoformatowy, druk gazetowy |
| `sygnet-w-kole.svg` | 96×96 | jasne | wytworzony | 421 B | awatar konta, znak na fotografii, sygnatura w mediach |
| `sygnet-w-kole-na-ciemnym.svg` | 96×96 | ciemne | wytworzony | 421 B | awatar na powierzchni atramentowej, plansza ciemna |
| `sygnet-w-kwadracie.svg` | 96×96 | jasne | wytworzony | 427 B | ikona aplikacji, kafel systemowy, skrót pulpitu |
| `sygnet-w-kwadracie-na-ciemnym.svg` | 96×96 | ciemne | wytworzony | 427 B | ikona na jasnym pulpicie, kafel odwrócony |

### 4.2 B · Logotyp — sama typografia

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `logotyp.svg` | 149×56 | jasne | źródłowy | 4.3 kB | stopka dokumentu, nagłówek listu, sygnatura tekstowa |
| `logotyp-na-ciemnym.svg` | 149×56 | ciemne | źródłowy | 4.3 kB | stopka na powierzchni atramentowej |
| `logotyp-mono-bialy.svg` | 149×56 | ciemne | wytworzony | 4.3 kB | nadruk 1/0, materiał jednobarwny na ciemnym |
| `logotyp-mono-czarny.svg` | 149×56 | jasne | wytworzony | 4.3 kB | dokument czarno-biały, faks, umowa |

### 4.3 C · Lockup poziomy — układ podstawowy

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `logo-poziomy.svg` | 225×96 | jasne | źródłowy | 4.6 kB | układ podstawowy — nagłówek strony, prezentacja, nośnik szeroki |
| `logo-poziomy-na-ciemnym.svg` | 225×96 | ciemne | źródłowy | 4.6 kB | układ podstawowy na powierzchni atramentowej |
| `logo-poziomy-mono-bialy.svg` | 225×96 | ciemne | źródłowy | 4.6 kB | nadruk 1/0 na ciemnym, materiał promocyjny |
| `logo-poziomy-mono-czarny.svg` | 225×96 | jasne | źródłowy | 4.6 kB | dokument czarno-biały, ogłoszenie prasowe |

### 4.4 D · Lockup pionowy

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `logo-pionowy.svg` | 159×172 | jasne | źródłowy | 4.6 kB | nośnik wąski i kwadratowy, plansza tytułowa, roll-up |
| `logo-pionowy-na-ciemnym.svg` | 159×172 | ciemne | źródłowy | 4.6 kB | plansza tytułowa na atramencie |
| `logo-pionowy-mono-bialy.svg` | 159×172 | ciemne | wytworzony | 4.6 kB | nadruk 1/0 na ciemnym w układzie pionowym |
| `logo-pionowy-mono-czarny.svg` | 159×172 | jasne | wytworzony | 4.6 kB | dokument jednobarwny w układzie pionowym |

### 4.5 E · Lockup kompaktowy — sygnet + DANACO

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `logo-kompaktowy.svg` | 225×96 | jasne | wytworzony | 1.9 kB | pasek górny, stopka, nośnik o niskiej wysokości |
| `logo-kompaktowy-na-ciemnym.svg` | 225×96 | ciemne | wytworzony | 1.9 kB | pasek górny motywu ciemnego, stopka atramentowa |

### 4.6 F · Znak z deskryptorem „AI Operating Environment"

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `logo-poziomy-deskryptor.svg` | 225×96 | jasne | wytworzony | 12.2 kB | materiał wprowadzający produkt, prezentacja, wizytówka |
| `logo-poziomy-deskryptor-na-ciemnym.svg` | 225×96 | ciemne | wytworzony | 12.2 kB | plansza wprowadzająca na atramencie |
| `logo-pionowy-deskryptor.svg` | 159×186 | jasne | wytworzony | 12.2 kB | okno startowe, plansza tytułowa, ekran powitalny |
| `logo-pionowy-deskryptor-na-ciemnym.svg` | 159×186 | ciemne | wytworzony | 12.2 kB | okno startowe motywu ciemnego |

### 4.7 G · Wariant do tłoczenia i grawerowania

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `sygnet-do-grawerowania.svg` | 96×96 | jasne | wytworzony | 351 B | grawer, tłoczenie, haft, pieczęć, matryca jednobarwna |
| `sygnet-do-grawerowania-odwrotny.svg` | 96×96 | ciemne | wytworzony | 351 B | grawer w materiale ciemnym, sitodruk 1/0 |
| `logo-poziomy-do-grawerowania.svg` | 225×96 | jasne | wytworzony | 4.6 kB | tabliczka, grawer laserowy, matryca tłoczna |

### 4.8 H · Znak wodny

| Plik | Wymiar bazowy | Tło | Pochodzenie | Rozmiar | Zastosowanie |
|---|---|---|---|---|---|
| `sygnet-znak-wodny.svg` | 96×96 | jasne | wytworzony | 339 B | podkład strony, papier firmowy, tło planszy |
| `sygnet-znak-wodny-na-ciemnym.svg` | 96×96 | ciemne | wytworzony | 339 B | podkład planszy atramentowej, tło sekcji |
| `logo-poziomy-znak-wodny.svg` | 225×96 | jasne | wytworzony | 4.6 kB | papier firmowy, stopka dokumentu, podkład wydruku |
| `logo-poziomy-znak-wodny-na-ciemnym.svg` | 225×96 | ciemne | wytworzony | 4.6 kB | podkład planszy atramentowej |

---

## 5. Katalog PNG — pełna tabela

Wszystkie rasteryzacje mają **tło przezroczyste** (kanał alfa) i powstały
z plików wektorowych z katalogu `svg/` biblioteką `cairosvg` — bez pośrednictwa
edytora, bez utraty kadru.

Szerokość jest parametrem sterującym; wysokość wynika z proporcji pliku
wektorowego.

#### Sygnet pełny — jasne tło

Źródło wektorowe: `svg/sygnet.svg` · zastosowanie: ikona okna, awatar, znak w interfejsie motywu jasnego

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-pelny-16.png` | 16×16 px | 378 B |
| `sygnet-pelny-24.png` | 24×24 px | 441 B |
| `sygnet-pelny-32.png` | 32×32 px | 816 B |
| `sygnet-pelny-48.png` | 48×48 px | 817 B |
| `sygnet-pelny-64.png` | 64×64 px | 1.5 kB |
| `sygnet-pelny-128.png` | 128×128 px | 2.4 kB |
| `sygnet-pelny-256.png` | 256×256 px | 3.9 kB |
| `sygnet-pelny-512.png` | 512×512 px | 8.0 kB |
| `sygnet-pelny-1024.png` | 1024×1024 px | 30.9 kB |

#### Sygnet pełny — ciemne tło

Źródło wektorowe: `svg/sygnet-na-ciemnym.svg` · zastosowanie: znak w interfejsie motywu ciemnego

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-pelny-na-ciemnym-16.png` | 16×16 px | 367 B |
| `sygnet-pelny-na-ciemnym-24.png` | 24×24 px | 396 B |
| `sygnet-pelny-na-ciemnym-32.png` | 32×32 px | 671 B |
| `sygnet-pelny-na-ciemnym-48.png` | 48×48 px | 745 B |
| `sygnet-pelny-na-ciemnym-64.png` | 64×64 px | 1.4 kB |
| `sygnet-pelny-na-ciemnym-128.png` | 128×128 px | 2.2 kB |
| `sygnet-pelny-na-ciemnym-256.png` | 256×256 px | 3.7 kB |
| `sygnet-pelny-na-ciemnym-512.png` | 512×512 px | 7.7 kB |
| `sygnet-pelny-na-ciemnym-1024.png` | 1024×1024 px | 28.6 kB |

#### Sygnet uproszczony — jasne tło

Źródło wektorowe: `svg/sygnet-uproszczony.svg` · zastosowanie: favicon, awatar 16–24 px, plakietka listy

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-uproszczony-16.png` | 16×16 px | 340 B |
| `sygnet-uproszczony-24.png` | 24×24 px | 547 B |
| `sygnet-uproszczony-32.png` | 32×32 px | 638 B |
| `sygnet-uproszczony-48.png` | 48×48 px | 969 B |
| `sygnet-uproszczony-64.png` | 64×64 px | 1.2 kB |
| `sygnet-uproszczony-128.png` | 128×128 px | 2.2 kB |
| `sygnet-uproszczony-256.png` | 256×256 px | 4.1 kB |
| `sygnet-uproszczony-512.png` | 512×512 px | 10.1 kB |
| `sygnet-uproszczony-1024.png` | 1024×1024 px | 25.3 kB |

#### Sygnet uproszczony — ciemne tło

Źródło wektorowe: `svg/sygnet-uproszczony-na-ciemnym.svg` · zastosowanie: favicon i awatar motywu ciemnego

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-uproszczony-na-ciemnym-16.png` | 16×16 px | 315 B |
| `sygnet-uproszczony-na-ciemnym-24.png` | 24×24 px | 461 B |
| `sygnet-uproszczony-na-ciemnym-32.png` | 32×32 px | 641 B |
| `sygnet-uproszczony-na-ciemnym-48.png` | 48×48 px | 922 B |
| `sygnet-uproszczony-na-ciemnym-64.png` | 64×64 px | 1.1 kB |
| `sygnet-uproszczony-na-ciemnym-128.png` | 128×128 px | 2.1 kB |
| `sygnet-uproszczony-na-ciemnym-256.png` | 256×256 px | 3.9 kB |
| `sygnet-uproszczony-na-ciemnym-512.png` | 512×512 px | 9.1 kB |
| `sygnet-uproszczony-na-ciemnym-1024.png` | 1024×1024 px | 24.7 kB |

#### Medalion okrągły

Źródło wektorowe: `svg/sygnet-w-kole.svg` · zastosowanie: awatar konta, znak na fotografii

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-w-kole-256.png` | 256×256 px | 8.3 kB |
| `sygnet-w-kole-512.png` | 512×512 px | 17.0 kB |

#### Medalion okrągły odwrócony

Źródło wektorowe: `svg/sygnet-w-kole-na-ciemnym.svg` · zastosowanie: awatar na powierzchni atramentowej

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-w-kole-na-ciemnym-256.png` | 256×256 px | 7.7 kB |

#### Kafel zaokrąglony

Źródło wektorowe: `svg/sygnet-w-kwadracie.svg` · zastosowanie: ikona aplikacji, kafel systemowy, sklep z aplikacjami

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-w-kwadracie-512.png` | 512×512 px | 11.1 kB |
| `sygnet-w-kwadracie-1024.png` | 1024×1024 px | 34.2 kB |

#### Kafel zaokrąglony odwrócony

Źródło wektorowe: `svg/sygnet-w-kwadracie-na-ciemnym.svg` · zastosowanie: ikona na jasnym pulpicie

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `sygnet-w-kwadracie-na-ciemnym-512.png` | 512×512 px | 10.7 kB |

#### Lockup poziomy

Źródło wektorowe: `svg/logo-poziomy.svg` · zastosowanie: nagłówek strony, prezentacja, nośnik szeroki

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-poziomy-400.png` | 400×171 px | 9.3 kB |
| `lockup-poziomy-800.png` | 800×341 px | 20.5 kB |
| `lockup-poziomy-1600.png` | 1600×683 px | 44.3 kB |

#### Lockup poziomy na ciemnym

Źródło wektorowe: `svg/logo-poziomy-na-ciemnym.svg` · zastosowanie: nagłówek na powierzchni atramentowej

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-poziomy-na-ciemnym-400.png` | 400×171 px | 9.2 kB |
| `lockup-poziomy-na-ciemnym-800.png` | 800×341 px | 20.2 kB |
| `lockup-poziomy-na-ciemnym-1600.png` | 1600×683 px | 43.9 kB |

#### Lockup poziomy z deskryptorem

Źródło wektorowe: `svg/logo-poziomy-deskryptor.svg` · zastosowanie: materiał wprowadzający produkt

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-poziomy-deskryptor-800.png` | 800×341 px | 28.9 kB |

#### Lockup poziomy z deskryptorem na ciemnym

Źródło wektorowe: `svg/logo-poziomy-deskryptor-na-ciemnym.svg` · zastosowanie: plansza wprowadzająca na atramencie

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-poziomy-deskryptor-na-ciemnym-800.png` | 800×341 px | 28.8 kB |

#### Lockup pionowy

Źródło wektorowe: `svg/logo-pionowy.svg` · zastosowanie: nośnik wąski, plansza tytułowa

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-pionowy-400.png` | 400×433 px | 14.6 kB |
| `lockup-pionowy-800.png` | 800×865 px | 36.3 kB |

#### Lockup pionowy na ciemnym

Źródło wektorowe: `svg/logo-pionowy-na-ciemnym.svg` · zastosowanie: plansza tytułowa na atramencie

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-pionowy-na-ciemnym-400.png` | 400×433 px | 14.4 kB |
| `lockup-pionowy-na-ciemnym-800.png` | 800×865 px | 35.3 kB |

#### Lockup kompaktowy

Źródło wektorowe: `svg/logo-kompaktowy.svg` · zastosowanie: pasek górny, stopka

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-kompaktowy-400.png` | 400×171 px | 7.3 kB |
| `lockup-kompaktowy-800.png` | 800×341 px | 16.3 kB |

#### Lockup kompaktowy na ciemnym

Źródło wektorowe: `svg/logo-kompaktowy-na-ciemnym.svg` · zastosowanie: pasek górny motywu ciemnego

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `lockup-kompaktowy-na-ciemnym-400.png` | 400×171 px | 7.1 kB |
| `lockup-kompaktowy-na-ciemnym-800.png` | 800×341 px | 16.1 kB |

#### Logotyp

Źródło wektorowe: `svg/logotyp.svg` · zastosowanie: stopka dokumentu, sygnatura tekstowa

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `logotyp-400.png` | 400×150 px | 11.1 kB |
| `logotyp-800.png` | 800×301 px | 22.8 kB |

#### Logotyp na ciemnym

Źródło wektorowe: `svg/logotyp-na-ciemnym.svg` · zastosowanie: stopka na atramencie

| Plik PNG | Rozdzielczość | Waga |
|---|---|---|
| `logotyp-na-ciemnym-400.png` | 400×150 px | 11.0 kB |
| `logotyp-na-ciemnym-800.png` | 800×301 px | 22.7 kB |

---

## 6. Katalog ICO

Pliki `.ico` są **wielorozmiarowe** — jeden plik zawiera trzy mapy bitowe
(16 × 16, 32 × 32, 48 × 48), a system albo przeglądarka wybiera właściwą.

| Plik | Zawarte rozmiary | Wariant znaku | Zastosowanie |
|---|---|---|---|
| `favicon.ico` | 16 · 32 · 48 px | sygnet uproszczony, jasne tło | favicon witryny i prototypów, zakładka przeglądarki |
| `favicon-na-ciemnym.ico` | 16 · 32 · 48 px | sygnet uproszczony, odwrócony | favicon dla powłok o ciemnym pasku zakładek |
| `ikona-aplikacji.ico` | 16 · 32 · 48 px | kafel zaokrąglony, pole atramentowe | skrót aplikacji na pulpicie Windows, pasek zadań |

**Dlaczego wariant uproszczony w faviconie:** favicon renderuje się w 16 px —
to dokładnie zakres, w którym kontrakt systemu projektowego nakazuje wariant uproszczony.
Użycie pełnego sygnetu w faviconie jest **błędem wdrożeniowym**, nie kwestią
gustu.

---

## 7. Dobór wariantu — drzewo decyzyjne

```
START: gdzie stanie znak?
  │
  ├─ Czy odbiorca zna już markę i widzi nazwę obok?
  │    TAK ─────────────────────────────► A · SYGNET
  │    │                                   └─ rozmiar < 24 px?  ──► uproszczony
  │    │                                   └─ własne pole / zdjęcie? ──► w-kole
  │    │                                   └─ ikona systemu?      ──► w-kwadracie
  │    NIE
  │     │
  ├─ Czy nośnik ma wysokość < 40 px?
  │    TAK ─────────────────────────────► E · LOCKUP KOMPAKTOWY
  │    NIE
  │     │
  ├─ Czy to pierwszy kontakt z produktem?
  │    TAK ─────────────────────────────► F · Z DESKRYPTOREM
  │    │                                   └─ nośnik szeroki  ──► poziomy
  │    │                                   └─ nośnik wąski    ──► pionowy
  │    NIE
  │     │
  ├─ Jaka jest proporcja pola?
  │    szerokie (≥ 2 : 1) ──────────────► C · LOCKUP POZIOMY   ← DOMYŚLNY
  │    kwadratowe lub wąskie ──────────► D · LOCKUP PIONOWY
  │
  └─ Ograniczenie technologii nośnika?
       jedna barwa, brak rastra ───────► G · TŁOCZENIE / GRAWER
       podkład pod treścią ────────────► H · ZNAK WODNY
       tylko tekst, brak grafiki ──────► B · LOGOTYP
```

**Reguła tła — dwa pytania, w tej kolejności:**

1. Czy powierzchnia jest ciemniejsza niż `--dn-szary-500` (#7C7C7C)?
   TAK → wariant `-na-ciemnym`. NIE → wariant podstawowy.
2. Czy technologia dopuszcza dwie barwy?
   NIE → wariant `mono-*` albo `do-grawerowania`.

**Godło zachowuje barwy własne w obu motywach interfejsu.** Przełączenie
motywu strony **nie zmienia** pliku znaku — zmiana pliku następuje wyłącznie
wtedy, gdy zmienia się **powierzchnia**, na której znak stoi.

---

## 8. Pole ochronne i rozmiary minimalne

### 8.1 Jednostka pola ochronnego

Jednostką jest **x = ½ wysokości sygnetu**. Wokół znaku pozostaje wolne pole
**≥ ½ x z każdej strony** — bez tekstu, innych znaków i krawędzi nośnika.

```
   ┌ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┐   ← pole ochronne ½x
   │                                          │
   │   ┌──────────────────────────────────┐   │
   │   │                                  │   │
   │   │      »»●   DANACO                │   │  ← znak
   │   │            CONSOLE●              │   │
   │   │                                  │   │
   │   └──────────────────────────────────┘   │
   │                                          │
   └ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ┘
      ½x                                  ½x        x = ½ wysokości sygnetu
```

### 8.2 Rozmiary minimalne

| Nośnik | Lockup pełny (C/D) | Lockup kompaktowy (E) | Sam sygnet (A) |
|---|---|---|---|
| **Ekran** | ≥ 120 px szerokości | ≥ 96 px szerokości | ≥ 24 px (16 px — wariant uproszczony) |
| **Druk** | ≥ 28 mm | ≥ 22 mm | ≥ 9 mm |
| **Grawer / haft / tłoczenie** | ≥ 40 mm | ≥ 32 mm | ≥ 15 mm |

Poniżej rozmiaru minimalnego **nie zmniejsza się znaku** — zmienia się
konfigurację na prostszą (C → E → A pełny → A uproszczony).

### 8.3 Progi przejścia między wariantami

```
  1024 px ┤ ████████████████████████████  sygnet pełny · kafel · medalion
   512 px ┤ ████████████████████████████
   256 px ┤ ████████████████████████████
   128 px ┤ ████████████████████████████
    64 px ┤ ████████████████████████████
    48 px ┤ ████████████████████████████
    32 px ┤ ████████████████████████████
    24 px ┤ ████████████████████████████  ← ostatni rozmiar pełnego sygnetu
    16 px ┤ ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  ← OBOWIĄZUJE wariant uproszczony
          └────────────────────────────────────────────────────────────
```

---

## 9. Barwy znaku i ich żetony

Pliki SVG znaku to **jedyne miejsce w całym systemie**, gdzie wartość
szesnastkowa jest zapisana wprost — barwa marki musi działać poza kontekstem
arkusza stylów (poczta, druk, grawer, plik przesłany dostawcy).
W dokumentach HTML i CSS obowiązuje bezwzględnie `var(--dn-*)`.

| Rola w znaku | Wartość | Żeton systemu | Występuje w |
|---|---|---|---|
| Atrament — groty i typografia na jasnym tle | `#181818` | `--dn-szary-900` | wszystkie warianty podstawowe i `mono-czarny` |
| Atrament jasny — groty i typografia na ciemnym tle | `#ECECEC` | `--dn-szary-100` | warianty `-na-ciemnym` wytworzone w tym opracowaniu |
| Biel źródłowa — j.w. w plikach źródłowych | `#F4F4F4` | `--dn-szary-50` | 6 plików przeniesionych z `zasoby/marka/logo/` |
| Kropka sygnału na jasnym tle | `#3B6FE0` | `--dn-sygnal-500` | warianty podstawowe |
| Kropka sygnału na ciemnym tle | `#5C8CEC` | `--dn-sygnal-400` | warianty `-na-ciemnym` |
| Pole medalionu atramentowego | `#131313` | `--dn-szary-925` | `sygnet-w-kole`, `sygnet-w-kwadracie` |
| Pole medalionu odwróconego | `#F4F4F4` | `--dn-szary-50` | `*-w-kole-na-ciemnym`, `*-w-kwadracie-na-ciemnym` |
| Deskryptor na jasnym tle | `#616161` | `--dn-szary-600` | konfiguracja F |
| Deskryptor na ciemnym tle | `#9E9E9E` | `--dn-szary-400` | konfiguracja F |

**Czerń absolutna (RGB 0, 0, 0) nie występuje w żadnym pliku repozytorium.**
Wariant „mono czarny" używa atramentu `#181818` — czerni systemu, nie czerni absolutnej.

### 9.1 Kontrast kropki względem tła

Wartości zmierzone metodą WCAG 2.1 (współczynnik kontrastu luminancji
względnej). Próg dla elementów graficznych i komponentów interfejsu wynosi
**3 : 1**; dla tekstu — 4,5 : 1.

| Powierzchnia | Element znaku | Kontrast | Ocena |
|---|---|---:|---|
| `#F4F4F4` papier roboczy | kropka `#3B6FE0` | 4,21 : 1 | element graficzny — **spełnia** |
| `#FFFFFF` powierzchnia karty | kropka `#3B6FE0` | 4,63 : 1 | **spełnia** |
| `#0F0F0F` podłoże motywu ciemnego | kropka `#5C8CEC` | 5,86 : 1 | **spełnia** |
| `#131313` rama i pole medalionu | kropka `#5C8CEC` | 5,68 : 1 | **spełnia** |
| `#F4F4F4` papier roboczy | atrament `#181818` | 16,14 : 1 | **spełnia z zapasem** |
| `#0F0F0F` podłoże motywu ciemnego | atrament jasny `#ECECEC` | 16,23 : 1 | **spełnia z zapasem** |
| `#0F0F0F` podłoże motywu ciemnego | biel źródłowa `#F4F4F4` | 17,43 : 1 | **spełnia z zapasem** |
| `#131313` pole medalionu | atrament jasny `#ECECEC` | 15,73 : 1 | **spełnia z zapasem** |
| `#F4F4F4` papier roboczy | deskryptor `#616161` | 5,63 : 1 | tekst — **spełnia** |
| `#0F0F0F` podłoże motywu ciemnego | deskryptor `#9E9E9E` | 7,15 : 1 | tekst — **spełnia** |

Kropka nigdy nie niesie informacji **samym kolorem** — w znaku jest formą
zamykającą kompozycję, nie stanem.

---

## 10. Matryca nośników

| Nośnik | Plik | Format | Uwaga |
|---|---|---|---|
| Pasek górny prototypu (motyw jasny) | `logo-kompaktowy.svg` | SVG inline | wysokość paska 48 px → sygnet 24 px |
| Pasek górny prototypu (motyw ciemny) | `logo-kompaktowy-na-ciemnym.svg` | SVG inline | pasek jest zawsze atramentowy → wariant odwrócony |
| Karta środowiska na stronie głównej | `sygnet.svg` / `sygnet-na-ciemnym.svg` | SVG inline | znak w barwie własnej |
| Zakładka przeglądarki | `ico/favicon.ico` + `svg/sygnet-uproszczony.svg` | ICO + SVG | SVG jako `rel="icon"` z przełączaniem motywu systemu |
| Ikona aplikacji na pulpicie | `png/sygnet-w-kwadracie-1024.png` | PNG | promień 21/96 = 21,9 % — zgodny z siatką ikon systemowych |
| Ikona aplikacji Windows | `ico/ikona-aplikacji.ico` | ICO | trzy mapy bitowe w jednym pliku |
| Awatar konta | `png/sygnet-w-kole-256.png` | PNG | pole atramentowe, znak nie dotyka krawędzi |
| Okno startowe / plansza tytułowa | `logo-pionowy-deskryptor-na-ciemnym.svg` | SVG | pierwszy kontakt — deskryptor obecny |
| Prezentacja, slajd tytułowy | `png/lockup-poziomy-1600.png` | PNG 1600 | tło przezroczyste, znak wstawiany na dowolną planszę |
| Papier firmowy — nagłówek | `logo-poziomy-mono-czarny.svg` | SVG | druk 1/0 |
| Papier firmowy — podkład | `logo-poziomy-znak-wodny.svg` | SVG | krycie 8 % |
| Dokument PDF — stopka | `logotyp.svg` | SVG | sama typografia, bez sygnetu |
| Stopka wiadomości e-mail | `png/lockup-poziomy-400.png` | PNG 400 | raster — klienty poczty nie renderują SVG |
| Tabliczka grawerowana | `logo-poziomy-do-grawerowania.svg` | SVG | jedna barwa, kropki jako pierścienie |
| Pieczęć / stempel | `sygnet-do-grawerowania.svg` | SVG | kropka konturem — czytelna w tuszu |
| Haft na odzieży | `sygnet-mono-bialy.svg` / `sygnet-mono-czarny.svg` | SVG | min. 15 mm |
| Materiał promocyjny 1/0 na ciemnym | `logo-poziomy-mono-bialy.svg` | SVG | jedna barwa |
| Ogłoszenie prasowe czarno-białe | `logo-poziomy-mono-czarny.svg` | SVG | jedna barwa |

---

## 11. Czego nie wolno robić ze znakiem

| # | Zakaz | Dlaczego |
|---|---|---|
| 1 | **Zmiana krzywych sygnetu** — nawet o setną jednostki | geometria jest zatwierdzona; każdy plik repozytorium wywodzi się z tych samych ścieżek |
| 2 | **Rozciąganie i ścinanie** (skalowanie nieproporcjonalne) | groty tracą kąt 45°, znak czyta się jako niechlujstwo |
| 3 | **Obracanie** | groty niosą kierunek „od koordynatora do wykonawcy"; obrót niszczy znaczenie |
| 4 | **Przebarwianie grotów na sygnałowy błękit** | kropka jest jedynym nośnikiem sygnału; barwny grot znosi tę hierarchię |
| 5 | **Zamiana kropki na inną barwę** (zieleń, bursztyn, czerwień) | barwy stanów należą do interfejsu, nie do godła |
| 6 | **Gradient w znaku** | gradienty są wyłącznie ilustracyjne (kontrakt systemu projektowego) i nigdy nie są tłem elementu |
| 7 | **Poświata, cień rzucony, obrys dodatkowy** | znak jest płaski; głębia należy do warstw interfejsu |
| 8 | **Pełny sygnet poniżej 24 px** | groty zlewają się; obowiązuje wariant uproszczony |
| 9 | **Znak na tle o kontraście < 3 : 1** | znak przestaje być rozpoznawalny |
| 10 | **Wstawianie znaku w ramkę, tarczę lub inne pole poza `w-kole` / `w-kwadracie`** | pole znaku jest zdefiniowane i skończone |
| 11 | **Rekonstrukcja logotypu z zainstalowanego kroju** | pliki zawierają krzywe; skład tekstowy da inne światła i inne rozstrzelenie |
| 12 | **Zamiana `CONSOLE` na nazwę modułu lub środowiska** | logotyp identyfikuje produkt, nie jego część |
| 13 | **Emoji, ikona zastępcza lub tekst „logo" w miejscu znaku** | zakaz systemowy (kontrakt systemu projektowego) |
| 14 | **Czerń absolutna (RGB 0, 0, 0) w jakimkolwiek pliku znaku** | zakaz systemowy — atramentem jest `#181818` |

---

## 12. Technika wytworzenia i odtworzenia

### 12.1 Skąd wzięły się wektory

- **Sygnet** — ścieżki wpisane wprost z definicji zatwierdzonej (kontrakt systemu projektowego,
  kierunek systemu projektowego). Kopiowane do każdej konfiguracji **bez modyfikacji**.
- **Logotyp** — 13 ścieżek liter odczytanych z `zasoby/marka/logo/logotyp.svg`
  i przeniesionych znak w znak. Skład tekstowy krojem **nie jest stosowany**.
- **Deskryptor** — jedyny nowy tekst w repozytorium; wektory wygenerowane
  z pliku `zasoby/fonty/ibm-plex-sans-latin-400-normal.woff2` (licencja OFL),
  osadzone jako krzywe. Dzięki temu plik nie zależy od kroju zainstalowanego
  u odbiorcy.

### 12.2 Rasteryzacja

```bash
python3 -c "import cairosvg; cairosvg.svg2png(
    url='svg/sygnet.svg',
    write_to='png/sygnet-pelny-512.png',
    output_width=512)"
```

Tło pozostaje przezroczyste — `background_color` nie jest ustawiane.
Wysokość wynika z proporcji pliku wektorowego, nie jest wymuszana.

### 12.3 Złożenie pliku ICO

```python
from PIL import Image
kadry = [Image.open(f'_ico_{s}.png').convert('RGBA') for s in (16, 32, 48)]
kadry[-1].save('ico/favicon.ico', format='ICO',
               sizes=[(16, 16), (32, 32), (48, 48)])
```

### 12.4 Weryfikacja pliku przed przekazaniem dostawcy

| Sprawdzenie | Jak |
|---|---|
| Czy krzywe sygnetu są nienaruszone | wyszukaj w pliku `M12 26 H24 L44 48.0` — musi wystąpić dosłownie |
| Czy nie ma czerni absolutnej | wyszukaj `000000` w pliku → brak trafień |
| Czy nie ma osadzonego rastra | `grep -c 'data:image' plik.svg` → 0 |
| Czy tekst jest krzywymi | plik nie zawiera znacznika `<text>` |
| Czy proporcja jest zachowana | `viewBox` i `width`/`height` mają tę samą proporcję |

---

## 13. Kontrola jakości

- [x] Geometria sygnetu identyczna we wszystkich 37 plikach SVG
- [x] Zero plików z czernią absolutną (RGB 0, 0, 0)
- [x] Zero znaczników `<text>` — cała typografia jako krzywe
- [x] Zero osadzonych rastrów w plikach SVG
- [x] Każdy wariant ma parę: jasne tło ↔ ciemne tło (albo uzasadnienie braku pary)
- [x] Wszystkie PNG mają kanał alfa i przezroczyste tło
- [x] ICO są wielorozmiarowe (16 / 32 / 48)
- [x] Favicon używa wariantu uproszczonego (reguła „poniżej 24 px")
- [x] Nazwy plików kebab-case, po polsku
- [x] Wartości szesnastkowe wyłącznie w plikach SVG znaku — nigdzie indziej
- [x] Polskie znaki diakrytyczne poprawne w całym dokumencie
- [x] Liczba plików: 102 (wymagane minimum: 60)

---

## 14. Decyzje projektowe

Rozstrzygnięcia podjęte w tym opracowaniu tam, gdzie dokumentacja źródłowa
nie rozstrzygała formy. Każde jest odwracalne i każde ma uzasadnienie.

### Barwa atramentu na ciemnym tle: `#ECECEC`

**Rozstrzygnięcie.** Warianty wytworzone w tym opracowaniu używają `#ECECEC`
(`--dn-szary-100`), czyli dokładnie wartości żetonu `--dn-tekst` motywu
ciemnego.

**Stan zastany.** Sześć plików przeniesionych z `zasoby/marka/logo/` używa
`#F4F4F4` (`--dn-szary-50`). Pliki te przeniesiono **bajt w bajt**, ponieważ
są wersją zatwierdzoną wymienioną w kontrakcie systemu projektowego — podmiana wartości bez
decyzji Właściciela byłaby zmianą znaku.

**Skutek.** Różnica jasności między `#F4F4F4` a `#ECECEC` to 3 punkty w skali
0–255 — niedostrzegalna gołym okiem, bez wpływu na kontrast (oba przekraczają
13 : 1 na podłożu `#0F0F0F`). **Rekomendacja:** przy najbliższej rewizji plików
źródłowych ujednolicić do `#ECECEC`, tak aby znak i tekst interfejsu miały
w motywie ciemnym tę samą wartość.

### Medalion jako pole własne znaku, nie jako powtórzenie tła

**Rozstrzygnięcie.** `sygnet-w-kole.svg` i `sygnet-w-kwadracie.svg` mają pole
**atramentowe** (`#131313`), mimo że nazwa nie zawiera przyrostka
`-na-ciemnym`. Warianty z przyrostkiem `-na-ciemnym` mają pole **jasne**.

**Uzasadnienie.** Medalion istnieje po to, by dać znakowi **własne pole** tam,
gdzie tło jest nieprzewidywalne (fotografia, powierzchnia barwna, kafel
systemowy). Pole w barwie tła byłoby niewidoczne i bezużyteczne. Przyrostek
`-na-ciemnym` czytamy więc konsekwentnie jako **„do stosowania na ciemnym
tle"**, a nie „o ciemnym wypełnieniu" — tak samo jak we wszystkich pozostałych
konfiguracjach.

### Skala znaku w medalionach: 0,58 (koło) i 0,62 (kafel)

**Rozstrzygnięcie.** Znak zajmuje 46,8 % szerokości koła i 50,0 % szerokości
kafla.

**Uzasadnienie.** Skala 0,62 dla kafla odtwarza dokładnie proporcję
zatwierdzonego pliku `zasoby/marka/ikona-aplikacji/ikona-aplikacji.svg`
(znak zajmuje tam 50 % szerokości pola 1024). Koło ma mniejszą powierzchnię
użytkową przy tej samej średnicy, dlatego znak zmniejszono o 6,5 % —
tak, by masa optyczna obu medalionów była równa.

### Deskryptor w IBM Plex Sans 400, w barwie tekstu drugiego rzędu

**Rozstrzygnięcie.** `AI Operating Environment` składa się krojem interfejsu
(IBM Plex Sans 400), w stopniu 10,5 px względem siatki znaku, w barwie
`--dn-szary-600` (jasne tło) / `--dn-szary-400` (ciemne tło).

**Uzasadnienie.** kontrakt systemu projektowego przypisuje trzy kroje trzem rolom: Space
Grotesk to marka, Plex Mono to maszyna, Plex Sans to **treść**. Deskryptor
jest zdaniem o produkcie, czyli treścią — nie nazwą i nie danymi.
Hierarchię buduje **barwa**, nie stopień: obniżenie do tekstu drugiego rzędu
podporządkowuje deskryptor logotypowi bez zmniejszania go do granicy
czytelności.

### Kropka w wariancie do grawerowania jako pierścień

**Rozstrzygnięcie.** W konfiguracji G kropka jest pierścieniem o grubości 2,5
(sygnet) i 1,4 (kropka po `CONSOLE`); **krawędź zewnętrzna pierścienia
zachowuje promień oryginału** (r = 6,5 i r = 3,4).

**Uzasadnienie.** Technologie jednobarwne (grawer, tłoczenie, pieczęć) nie
odróżniają kropki od grotów barwą. Pierścień przywraca to rozróżnienie
formą, a zachowanie promienia zewnętrznego utrzymuje kompozycję i pole
ochronne bez zmiany.

### Lockup kompaktowy w kadrze lockupu poziomego

**Rozstrzygnięcie.** `logo-kompaktowy.svg` ma kadr 225 × 96 — identyczny
z `logo-poziomy.svg`. Sygnet stoi w tym samym miejscu i w tej samej skali;
`DANACO` jest wyśrodkowane w pionie na osi optycznej sygnetu (linia bazowa
y = 59,90).

**Uzasadnienie.** Jednakowy kadr pozwala zamienić lockup pełny na kompaktowy
**bez przeliczania układu** — znak nie przeskoczy i nie zmieni rozmiaru.
Kompaktowość dotyczy treści (brak linii `CONSOLE`), nie kadru.

### Wyrównanie deskryptora

**Rozstrzygnięcie.** W układzie poziomym lewa krawędź **farby** deskryptora
jest zrównana z lewą krawędzią farby `DANACO` (x = 77,56 w kadrze 225 × 96).
W układzie pionowym farba deskryptora jest wyśrodkowana na **osi farby**
`DANACO` (x = 79,66 w kadrze 159 × 186).

**Uzasadnienie.** Wyrównanie do szerokości znaku (advance width) zamiast do
farby dawało przesunięcie 1,3 px w lewo — widoczne przy skali prezentacyjnej.
Wyrównanie liczone na skrajnych punktach krzywych jest niezależne od świateł
bocznych kroju.

### Krycie znaku wodnego: 8 % na jasnym, 10 % na ciemnym

**Rozstrzygnięcie.** `sygnet-znak-wodny.svg` i `logo-poziomy-znak-wodny.svg`
mają krycie 0,08; warianty na ciemnym — 0,10.

**Uzasadnienie.** Przy tym samym kryciu znak na ciemnym podłożu czyta się
słabiej (rozjaśnienie działa na mniejszy zakres percepcyjny niż
przyciemnienie). Podniesienie o 2 punkty procentowe wyrównuje odczuwalną
obecność podkładu w obu motywach. Znak wodny nigdy nie przekracza 12 % —
powyżej tej wartości konkuruje z treścią.

### PNG bez wariantów mono

**Rozstrzygnięcie.** Rasteryzacje obejmują warianty podstawowe i `-na-ciemnym`;
warianty `mono-*` i `do-grawerowania` pozostają wyłącznie wektorowe.

**Uzasadnienie.** Odbiorcami wariantów jednobarwnych są drukarnie, grawernie
i hafciarnie — wszyscy przyjmują wektor i wszyscy odrzucają raster.
Rasteryzacja tych plików tworzyłaby zasób, którego użycie byłoby błędem.

### Trzy pliki ICO zamiast jednego

**Rozstrzygnięcie.** Repozytorium zawiera `favicon.ico`, `favicon-na-ciemnym.ico`
i `ikona-aplikacji.ico`.

**Uzasadnienie.** `favicon.ico` obsługuje przeglądarki bez wsparcia dla
`favicon.svg` z zapytaniem `prefers-color-scheme`; wariant odwrócony pokrywa
powłoki o ciemnym pasku zakładek; `ikona-aplikacji.ico` to inna konfiguracja
znaku (kafel z polem), a nie inny rozmiar tego samego pliku.

### `PORTFOLIO.md` jako jedyny spis

**Rozstrzygnięcie.** Repozytorium nie zawiera pliku manifestu maszynowego
(JSON/YAML). Spisem jest ten dokument.

**Uzasadnienie.** Odbiorcą repozytorium jest człowiek podejmujący decyzję
„którego pliku użyć". Manifest maszynowy dublowałby tę treść i natychmiast
rozjechałby się z dokumentem. Gdy pojawi się proces budowania korzystający
z manifestu, należy go **wygenerować** z tego katalogu, nie prowadzić równolegle.

---

## Podsumowanie liczbowe

| Kategoria | Liczba plików |
|---|---:|
| SVG — konfiguracje wektorowe | 37 |
| PNG — rasteryzacje z kanałem alfa | 62 |
| ICO — ikony wielorozmiarowe | 3 |
| **Razem** | **102** |

| Grupa konfiguracji | SVG |
|---|---:|
| A · sygnet | 12 |
| B · logotyp | 4 |
| C · lockup poziomy | 4 |
| D · lockup pionowy | 4 |
| E · lockup kompaktowy | 2 |
| F · znak z deskryptorem | 4 |
| G · tłoczenie i grawer | 3 |
| H · znak wodny | 4 |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
