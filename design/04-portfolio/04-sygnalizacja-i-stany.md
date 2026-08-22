# Danaco Console — Plansza portfolio 04: System sygnalizacji i stanów

| | |
|---|---|
| **Produkt** | Danaco Console — AI Operating Environment (warstwa wizualna v2.0) |
| **Warstwa funkcjonalna** | Danaco Console — Platforma AI Workspace OS (dokumentacja projektowa v1.0) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Rodzaj opracowania** | Plansza portfolio design identity (nie dokumentacja techniczna) |
| **Odbiorcy** | Właściciel · Designer · Deweloper · odbiorca portfolio |
| **Zakres** | Kropka sygnału jako element sygnaturowy · cztery rodziny stanów z pomiarem kontrastu · zasada „stan nigdy samym kolorem" · katalog nośników stanu · zero blokad jako sygnalizacja · katalog wszystkich modeli stanów platformy · budżet akcentu |
| **Czego NIE zawiera** | Wnętrz okien operacyjnych (te mają prototypy w `05-okna/`), tożsamości środowisk (plansza 01), tożsamości modułów (plansza 02), ról i orkiestracji (plansza 03) |
| **Plik towarzyszący** | `04-portfolio/04-sygnalizacja-i-stany.html` |

---

## Spis treści

1. [Czym jest ta plansza](#1-czym-jest-ta-plansza)
2. [Kropka sygnału — element sygnaturowy](#2-kropka-sygnału--element-sygnaturowy)
3. [Cztery rodziny stanów](#3-cztery-rodziny-stanów)
4. [Zasada „stan nigdy samym kolorem"](#4-zasada-stan-nigdy-samym-kolorem)
5. [Nośniki stanu — katalog](#5-nośniki-stanu--katalog)
6. [Zero blokad jako sygnalizacja](#6-zero-blokad-jako-sygnalizacja)
7. [Katalog modeli stanów platformy](#7-katalog-modeli-stanów-platformy)
8. [Budżet akcentu](#8-budżet-akcentu)
9. [Decyzje projektowe planszy](#9-decyzje-projektowe-planszy)
10. [Źródła](#10-źródła)
11. [Kontrola jakości](#11-kontrola-jakości)

---

## 1. Czym jest ta plansza

Danaco Console jest **monochromatyczna z zasady**. Odcienie bieli w motywie jasnym, odcienie czerni
w motywie ciemnym, jeden chłodny błękit sygnałowy — i nic więcej. W takim systemie **każde użycie
barwy jest zdarzeniem**. Nie dekoracją, nie „wyróżnieniem sekcji", nie sposobem na rozbicie monotonii.
Barwa w tym interfejsie znaczy: *coś się dzieje* albo *coś się stało*.

Plansza 04 pokazuje, jak ten jeden mechanizm jest zbudowany i jak działa w komplecie.

Trzy twierdzenia, które plansza demonstruje na żywo:

1. **Jest jeden element sygnaturowy — kropka sygnału.** Występuje w pięciu miejscach platformy
   i w każdym z nich znaczy dokładnie to samo: *tu biegnie praca*. Jej tętno (2,4 s) jest
   **jedynym ruchem ciągłym** w całym interfejsie.
2. **Są cztery rodziny stanów, nie dziewięć.** Sukces, ostrzeżenie, błąd, informacja — przy czym
   rodzina informacyjna **jest** rodziną sygnału (celowe scalenie, nie oszczędność). Każda para
   tekst/tło zmierzona wzorem WCAG, nie zadeklarowana.
3. **Stan nigdy nie jest niesiony samym kolorem.** Zawsze towarzyszy mu ikona albo etykieta.
   Konsekwencja: interfejs pozostaje czytelny przy widzeniu monochromatycznym.

Czwarte twierdzenie jest osobliwością tej platformy i wymaga własnego rozdziału:
**niedostępność nie jest komunikowana wyszarzeniem.** Danaco Console nie stosuje atrybutu
wyłączającego kontrolkę. Przycisk jest zawsze klikalny; niegotowość sygnalizuje się
**komunikatem po naciśnięciu** albo **opisem obok**. To Zasada zero blokad, opisana w rozdziale 6.

### 1.1. Liczby planszy

| Wielkość | Wartość | Źródło |
|---|---|---|
| Rodziny stanów | **4** (sukces · ostrzeżenie · błąd · informacja) | kontrakt systemu projektowego |
| Miejsca wystąpienia kropki sygnału | **5** | kontrakt systemu projektowego · kierunek systemu projektowego |
| Średnica kropki sygnału | **6 px** (`--dn-wym-kropka`) | `zetony.css` rozdz. 6 |
| Czas tętna | **2,4 s** (`--dn-czas-tetno`) | `zetony.css` rozdz. 5 |
| Pomiary kontrastu w pakiecie | **33** | `zasoby/zetony/kontrasty.json` |
| Nośniki stanu w katalogu `.dn-*` | **8** | kontrakt systemu projektowego · `komponenty.css` |
| Warianty nośników stanu | **23** | `komponenty.css` |
| Modele stanów w dokumentacji | **19** | rozdz. 7 tej planszy |
| Nazwane stany łącznie | **105** | rozdz. 7 tej planszy |
| Budżet akcentu | **≤ 5 % powierzchni ekranu** | opracowanie o przekazaniu, ruchu i dostępności |
| Intensywność ruchu | **3/10** | kontrakt systemu projektowego |

### 1.2. Siedem zagadnień planszy

| # | Sekcja | Co demonstruje | Źródło |
|---|---|---|---|
| 01 | Kropka sygnału | pięć miejsc wystąpienia, każde na żywo; tętno; zachowanie przy ograniczonym ruchu | kontrakt systemu projektowego · kierunek systemu projektowego |
| 02 | Cztery rodziny stanów | żetony tekst/tło/obrys w obu motywach + kontrast liczony na żywo | kontrakt systemu projektowego · kierunek systemu projektowego · `kontrasty.json` |
| 03 | Stan nigdy samym kolorem | pary porównawcze + symulacja widzenia monochromatycznego | kontrakt systemu projektowego · kierunek systemu projektowego |
| 04 | Nośniki stanu | 8 nośników, 23 warianty, wszystkie na żywo | `komponenty.css` |
| 05 | Zero blokad | dwa wzorce komunikowania niedostępności + kontrprzykład | kontrakt systemu projektowego (zasada zero blokad) |
| 06 | Katalog modeli stanów | 19 modeli, 105 stanów, tabela sortowalna i filtrowana | rozdz. 5.1/6.2 modułów · multitaskingai 12 |
| 07 | Budżet akcentu | kalkulator udziału sygnału w powierzchni ekranu | opracowanie o przekazaniu, ruchu i dostępności |

---

## 2. Kropka sygnału — element sygnaturowy

### 2.1. Definicja

> Wypełniony punkt błękitu sygnałowego oznaczający „**tu biegnie praca**".
> — kontrakt systemu projektowego

Kropka nie jest ikoną i nie należy do zestawu ikon. Jest **elementem geometrycznym systemu**:
koło o średnicy `--dn-wym-kropka` (6 px) wypełnione żetonem `--dn-kropka`
(`--dn-sygnal-500` w motywie jasnym, `--dn-sygnal-400` w ciemnym — jaśniejszy stopień na ciemnym
tle utrzymuje tę samą wagę optyczną).

### 2.2. Pięć miejsc wystąpienia

| # | Miejsce | Postać | Znaczenie w tym miejscu | Nośnik techniczny |
|---|---|---|---|---|
| 1 | **Godło** | kropka zamykająca podwójny grot „»»." | zlecenie przekazane — praca ruszyła | `marka/logo/sygnet.svg`, okrąg `cx=83 cy=63.5 r=6.5` |
| 2 | **Emblematy czterech środowisk** | dokładnie jedna wypełniona kropka w każdym emblemacie | punkt aktywności właściwy środowisku | `marka/srodowiska/srodowisko-*.svg` |
| 3 | **Karty sesji** | kropka pulsująca przy tytule karty | proces biegnie w tle tej sesji | `.dn-kropka.dn-kropka--tetno` |
| 4 | **Okno komunikacji** | kropka przy nadawcy aktywnie piszącym | ten nadawca właśnie generuje odpowiedź | `.dn-wpis--pracuje` (pseudoelement `::after` przy nadawcy) |
| 5 | **Always On Display** | rdzeń awatara z pulsującym pierścieniem | proces nadzorowany jest żywy | `.dn-aod-rdzen` (pseudoelement `::after`) |

Wspólny mianownik pięciu miejsc: **kropka nigdy nie oznacza statusu jako takiego** —
oznacza *trwanie pracy*. Status (powodzenie, ostrzeżenie, błąd) niosą warianty barwne
`.dn-kropka--sukces / --ostrzezenie / --blad / --neutralna`, które **nie pulsują**.
Ruch zarezerwowany jest dla jednego znaczenia.

### 2.3. Położenie kropki w czterech emblematach

| Środowisko | Plik | Element wypełniony | Położenie w kompozycji |
|---|---|---|---|
| TalkIn | `srodowisko-talkin.svg` | `circle cx=16.2 cy=10.8 r=1.5` | koniec drugiego wiersza dymka — miejsce, w którym rozmowa trwa |
| WorkSpace | `srodowisko-workspace.svg` | `circle cx=16.9 cy=16.9 r=1.6` | wnętrze czwartego kafla siatki — kafel, nad którym trwa praca |
| CodeStudio | `srodowisko-codestudio.svg` | `circle cx=18.4 cy=9.9 r=1.5` | prawa górna część okna edytora — wskaźnik uruchomienia |
| MultitaskingAI | `srodowisko-multitaskingai.svg` | `circle cx=12 cy=12 r=2.6` | **środek** grafu — koordynator; trzy węzły zewnętrzne pozostają obrysowe |

Reguła kompozycyjna: kropka jest **jedynym elementem wypełnionym** w każdym z czterech
emblematów. Reszta rysunku to obrys 1,75 na siatce 24×24, zgodnie z zasadami zestawu ikon.
MultitaskingAI odwraca kierunek — kropka jest w środku, nie na obrzeżu — bo w tym środowisku
praca nie „dzieje się gdzieś", tylko jest z jednego miejsca rozdzielana.

### 2.4. Tętno — jedyny ruch ciągły

| Parametr | Wartość | Uzasadnienie |
|---|---|---|
| Czas cyklu | 2,4 s (`--dn-czas-tetno`) | poniżej częstotliwości, przy której ruch staje się natrętny; powyżej progu, przy którym przestaje być zauważalny |
| Krzywa | `--dn-ease` = `cubic-bezier(0.2, 0, 0, 1)` | ta sama krzywa co wszystkie przejścia systemu |
| Postać | rozchodzący się pierścień `box-shadow` 0 → 5 px, barwa `--dn-fokus-cien` | pierścień nie zmienia geometrii kropki — nie powoduje przeskoków układu |
| Liczba jednocześnie | jeden ruch znaczący na widok | kontrakt systemu projektowego — wzorzec animacji |

Wszystkie inne animacje systemu są **przejściami** (100–220 ms) i kończą się.
Tętno jest jedyną animacją **nieskończoną**. To jest cała różnica: gdy coś w interfejsie
Danaco Console porusza się bez końca, to znaczy, że gdzieś biegnie praca.

### 2.5. Zachowanie przy `prefers-reduced-motion`

Preferencja obsłużona **globalnie w żetonach**, nie per komponent
(`zetony.css` rozdz. 12, kontrakt systemu projektowego):

```
@media (prefers-reduced-motion: reduce) {
  :root { --dn-czas-1: 0.01ms; --dn-czas-2: 0.01ms;
          --dn-czas-3: 0.01ms; --dn-czas-tetno: 0.01ms; }
}
```

Dodatkowo `komponenty.css` podmienia tętno na **pierścień statyczny**:

```
@media (prefers-reduced-motion: reduce) {
  .dn-kropka--tetno { box-shadow: 0 0 0 2px var(--dn-fokus-cien); }
}
```

Informacja nie ginie — zmienia nośnik. Kropka z pierścieniem nadal odróżnia się od kropki bez
pierścienia, więc „proces w tle" pozostaje odczytywalny bez ruchu. Plansza HTML udostępnia
**przełącznik symulacji** tej preferencji, żeby oba warianty można było zestawić obok siebie
bez zmiany ustawień systemu operacyjnego.

---

## 3. Cztery rodziny stanów

### 3.1. Skład rodzin

| Rodzina | Barwa | Prymitywy | Ikona towarzysząca | Znaczenie w platformie |
|---|---|---|---|---|
| **Sukces** | zieleń | `--dn-zielen-700 / -400 / -200 / -100` | `ptaszek` | operacja zakończona zgodnie z zamierzeniem |
| **Ostrzeżenie** | bursztyn | `--dn-bursztyn-700 / -400 / -200 / -100` | `ostrzezenie` | problem zasygnalizowany — **nigdy nie blokuje** wykonania |
| **Błąd** | czerwień | `--dn-czerwien-700 / -400 / -200 / -100` | `blad` | awaria techniczna poza kontrolą formularza |
| **Informacja** | **rodzina sygnału** | `--dn-sygnal-600 / -300 / -200 / -100` | `info` | stan neutralny wymagający uwagi; przebieg w toku |

### 3.2. Żetony semantyczne — oba motywy

| Żeton | Motyw jasny | Motyw ciemny |
|---|---|---|
| `--dn-sukces-tekst` | `--dn-zielen-700` (#1E7A46) | `--dn-zielen-400` (#4FBD85) |
| `--dn-sukces-tlo` | `--dn-zielen-100` (#E7F5EE) | `rgba(79, 189, 133, 0.14)` |
| `--dn-sukces-obrys` | `--dn-zielen-200` (#BCE3CD) | `rgba(79, 189, 133, 0.34)` |
| `--dn-ostrzezenie-tekst` | `--dn-bursztyn-700` (#8A5B0C) | `--dn-bursztyn-400` (#E0A63C) |
| `--dn-ostrzezenie-tlo` | `--dn-bursztyn-100` (#FBF2DE) | `rgba(224, 166, 60, 0.14)` |
| `--dn-ostrzezenie-obrys` | `--dn-bursztyn-200` (#EED9A7) | `rgba(224, 166, 60, 0.34)` |
| `--dn-blad-tekst` | `--dn-czerwien-700` (#C0362F) | `--dn-czerwien-400` (#EE7168) |
| `--dn-blad-tlo` | `--dn-czerwien-100` (#FCEDEB) | `rgba(238, 113, 104, 0.14)` |
| `--dn-blad-obrys` | `--dn-czerwien-200` (#F3C8C4) | `rgba(238, 113, 104, 0.34)` |
| `--dn-informacja-tekst` | `--dn-sygnal-600` (#2457C9) | `--dn-sygnal-300` (#8FB2F5) |
| `--dn-informacja-tlo` | `--dn-sygnal-100` (#EDF3FE) | `rgba(59, 111, 224, 0.16)` |
| `--dn-informacja-obrys` | `--dn-sygnal-200` (#C9DAFB) | `rgba(92, 140, 236, 0.40)` |

Konstrukcja motywu ciemnego jest inna niż jasnego i **nie jest z niego wywiedziona**:
w motywie jasnym tło stanu jest **osobnym stopniem skali** (100), w motywie ciemnym jest
**przezroczystością barwy stanu nałożoną na powierzchnię** (alfa 0,14 dla tła, 0,34 dla obrysu).
Powód jest praktyczny: na ciemnym tle stały stopień 100 byłby jaśniejszy od powierzchni karty
i zachowywałby się jak plama światła. Przezroczystość utrzymuje tło stanu w tej samej rodzinie
głębi co powierzchnia, na której leży.

### 3.3. Kontrasty zmierzone

Plansza HTML **nie odczytuje liczb z pliku** — liczy je na żywo wzorem WCAG 2.1
z faktycznie wyrenderowanych wartości żetonów, dla obu motywów jednocześnie
(chwilowe przełączenie `data-theme` na czas pomiaru, bez zapisu preferencji).
Wartości alfa są kompozytowane na `--dn-powierzchnia` przed pomiarem.

Wzór:

```
L = 0,2126·R + 0,7152·G + 0,0722·B          (składowe zlinearyzowane)
c = (c₂₅₅/255 ≤ 0,03928) ? c/12,92 : ((c+0,055)/1,055)^2,4
K = (L_jaśniejszy + 0,05) / (L_ciemniejszy + 0,05)
```

Wartości odniesienia z `zasoby/zetony/kontrasty.json` (33 pomiary), które plansza odtwarza:

| Para | Kontrast | Próg | Wynik |
|---|---|---|---|
| sukces tekst / tło stanu (jasny) | 4,76 | 4,5 | spełnia |
| ostrzeżenie tekst / tło stanu (jasny) | 5,26 | 4,5 | spełnia |
| błąd tekst / tło stanu (jasny) | 4,84 | 4,5 | spełnia |
| sygnał-tekst / sygnał-tło (jasny) | 5,74 | 4,5 | spełnia |
| sukces tekst / powierzchnia (jasny) | 5,35 | 4,5 | spełnia |
| ostrzeżenie tekst / powierzchnia (jasny) | 5,86 | 4,5 | spełnia |
| błąd tekst / powierzchnia (jasny) | 5,51 | 4,5 | spełnia |
| sukces tekst / powierzchnia (ciemny) | 7,57 | 4,5 | spełnia |
| ostrzeżenie tekst / powierzchnia (ciemny) | 8,19 | 4,5 | spełnia |
| błąd tekst / powierzchnia (ciemny) | 6,09 | 4,5 | spełnia |
| sygnał-tekst / powierzchnia (ciemny) | 8,33 | 4,5 | spełnia |
| fokus sygnał / tło (jasny) | 4,21 | 3,0 | spełnia |
| fokus sygnał / tło (ciemny) | 5,86 | 3,0 | spełnia |

**Obrysy stanów mierzone są, lecz nie podlegają progowi tekstowemu.** Obrys plakietki jest
nośnikiem **pomocniczym** — stan niesie tekst i ikona, więc obrys nie jest elementem koniecznym
do rozpoznania komponentu. Plansza podaje pomiar i oznacza rolę, nie wystawia werdyktu
zaliczenia. Dla porównania: `kontrasty.json` przyjmuje dla obrysu kontrolki próg dostrzegalności
**1,6** (pomiar „obrys kontrolki / tło (jasny)" = 1,65).

### 3.4. Dlaczego informacja jest rodziną sygnału

Rozstrzygnięcie z kierunku systemu projektowego: **rodzina informacyjna stanów = rodzina sygnału
(celowe scalenie)**.

Uzasadnienie projektowe:

1. **Nie ma czwartej barwy do wydania.** System deklaruje monochrom + jeden sygnał.
   Dodanie osobnego błękitu „informacyjnego" obok błękitu sygnałowego złamałoby deklarację
   przy zerowym zysku semantycznym.
2. **Znaczenia się pokrywają.** Sygnał znaczy „tu biegnie praca"; stan informacyjny znaczy
   „to jest w toku, zwróć uwagę". To ten sam komunikat na dwóch poziomach szczegółowości.
3. **Odróżnienie i tak istnieje** — niesie je nośnik, nie barwa: kropka pulsująca to praca w tle,
   plakietka `--informacja` z ikoną `info` to opisany stan pozycji, toast `--informacja`
   to zdarzenie jednorazowe.

Konsekwencja praktyczna: `.dn-plakietka--informacja` i `.dn-plakietka--sygnal` mają barwy
z tej samej rodziny, lecz **różne żetony** (`--dn-informacja-tekst` vs `--dn-sygnal`)
i różne przeznaczenie — pierwsza opisuje stan bytu, druga wyróżnia bieżący wybór.

---

## 4. Zasada „stan nigdy samym kolorem"

### 4.1. Reguła

> **Stan nigdy samym kolorem** — zawsze ikona albo etykieta.
> — kontrakt systemu projektowego (Dostępność) i rozdz. 1 (katalog anty-domyślnych)

W katalogu anty-domyślnych zapisana jako zablokowany odruch:

| Odruch | Zamiast tego |
|---|---|
| Stan samym kolorem | **zawsze** ikona lub etykieta |
| Dziewięć kolorów tła dla dziewięciu nadawców | trzy klasy semantyczne + ikona + etykieta + plakietka roli |

Druga pozycja jest tą samą regułą zastosowaną do okna komunikacji: dziewięciu nadawców nie
rozróżnia się dziewięcioma barwami, tylko trzema klasami semantycznymi
(`--czlowiek`, `--inteligencja`, `--system`) uzupełnionymi o medalion z ikoną, nazwę nadawcy
i plakietkę roli.

### 4.2. Pary porównawcze — cztery przypadki

Plansza HTML zestawia każdy przypadek w dwóch wersjach obok siebie: **niezgodnej** (tylko kolor)
i **zgodnej** (kolor + ikona + etykieta). Obie wersje są renderowane z tych samych żetonów.

| # | Kontekst (dane przykładowe z domeny produktu) | Wersja niezgodna | Wersja zgodna |
|---|---|---|---|
| 1 | Kolumna „Stan" w tabeli kolejki zadań (silnik kolejek MultitaskingAI) | sama kropka barwna w komórce | kropka + etykieta stanu („W realizacji", „Do powtórzenia") |
| 2 | Plakietka przy przebiegu automatyki (Execution Monitor) | plakietka z samą barwą tła, bez treści | plakietka z ikoną `ostrzezenie` + tekst „Błędna" |
| 3 | Powiadomienie o zakończeniu operacji (toast) | pasek barwny bez ikony i bez tytułu rodzaju | ikona `ptaszek` + tytuł + treść |
| 4 | Pole formularza z problemem (Agent Builder — nazwa zduplikowana) | sam obrys bursztynowy pola | obrys + komunikat pod polem z ikoną |

### 4.3. Symulacja widzenia monochromatycznego

Plansza stosuje filtr CSS `grayscale(1)` do obszaru porównawczego. Efekt jest jednoznaczny:

| Wersja | Po odjęciu barwy | Czy stan pozostaje odczytywalny |
|---|---|---|
| Niezgodna (tylko kolor) | cztery kropki różnią się wyłącznie jasnością szarości; plakietki stają się identyczne | **nie** |
| Zgodna (kolor + ikona + etykieta) | kształt ikony i treść etykiety nie zależą od barwy | **tak** |

To nie jest test dla wąskiej grupy odbiorców. Ten sam efekt daje wydruk czarno-biały,
zrzut ekranu w skali szarości w zgłoszeniu błędu, ekran o wypaczonej kalibracji i praca
przy silnym świetle bocznym. Reguła „stan nigdy samym kolorem" broni interfejsu w każdym
z tych przypadków przy zerowym koszcie wizualnym.

### 4.4. Konsekwencja dla nośników

Każdy nośnik stanu w katalogu `.dn-*` ma wbudowane miejsce na drugi kanał informacji:

| Nośnik | Drugi kanał obok barwy |
|---|---|
| `.dn-plakietka` | treść tekstowa (obowiązkowa) + opcjonalny `<svg>` 12 px |
| `.dn-kropka` | etykieta obok (komponent nigdy nie stoi sam — komentarz w `komponenty.css`) |
| `.dn-toast` | ikona rodzaju wstrzykiwana przez `prototyp.js` + tytuł |
| `.dn-krok` | `.dn-krok-znak` z ikoną lub numerem + treść kroku + `.dn-krok-meta` |
| `.dn-postep` | `.dn-postep-etykieta` z wartością liczbową krojem mono |
| `.dn-spinner` | etykieta tekstowa obok (np. „Pracuję…") |
| `.dn-wpis--pracuje` | nazwa nadawcy przed kropką |

Komentarz w `komponenty.css` przy wariantach `.dn-kropka` brzmi wprost:
„*Warianty stanu (zawsze obok etykiety albo ikony — nigdy same)*". Reguła jest zapisana
w kodzie systemu, nie tylko w dokumencie.

---

## 5. Nośniki stanu — katalog

Osiem nośników, 23 warianty. Wszystkie pochodzą z `zasoby/css/komponenty.css`
i z listy wiążącej w kontrakcie systemu projektowego.

### 5.1. `.dn-plakietka` — 6 wariantów

| Wariant | Żetony | Zastosowanie w dokumentacji |
|---|---|---|
| `--sukces` | `--dn-sukces-tlo / -obrys / -tekst` | „Aktywna" (wersja eksperta, Agents 4.1); proces zakończony sukcesem (MTAI 12.1) |
| `--ostrzezenie` | `--dn-ostrzezenie-*` | projekt wstrzymany (Workspace 6.2); proces wstrzymany (MTAI 12.1) |
| `--blad` | `--dn-blad-*` | automatyka błędna (Automations 6.2); proces zakończony błędem (MTAI 12.1) |
| `--informacja` | `--dn-informacja-*` | etap oczekujący na decyzję (MTAI 12.1, z ikoną `oko`) |
| `--sygnal` | `--dn-sygnal-tlo / -obrys / --dn-sygnal` | wyróżnienie bieżącego wyboru |
| `--rola` | krój mono, wersaliki, `--dn-r-xs` | plakietka roli okna: koordynator / wykonawca / walidator (kontrakt systemu projektowego) |

Wariant `--rola` różni się od pozostałych **nie barwą, lecz typografią** — IBM Plex Mono,
wersaliki, rozstrzelenie `--dn-ls-mono-wersaliki`, mniejszy promień. To jest wizualne
rozróżnienie klasy komunikatu: pozostałe pięć mówi o *stanie*, ten jeden o *tożsamości*.

### 5.2. `.dn-kropka` — 5 wariantów

| Wariant | Barwa | Znaczenie |
|---|---|---|
| bazowa | `--dn-kropka` (sygnał) | tu biegnie praca |
| `--sukces` | `--dn-sukces-tekst` | rola pracuje poprawnie (MTAI 12.2) |
| `--ostrzezenie` | `--dn-ostrzezenie-tekst` | oczekiwanie na zależność (MTAI 12.2); niezapisane zmiany (Agents 4.1) |
| `--blad` | `--dn-blad-tekst` | błąd wykonania roli (MTAI 12.2) |
| `--neutralna` | `--dn-tekst-3` | rola bezczynna (MTAI 12.2) |

Modyfikator `--tetno` jest **ortogonalny** do wariantów barwnych — dokłada ruch, nie barwę.
W praktyce systemu łączy się wyłącznie z kropką bazową: pulsuje praca, nie status.

### 5.3. `.dn-toast` — 4 warianty

| Wariant | Wstęga lewa | Ikona (wstrzykiwana przez `prototyp.js`) |
|---|---|---|
| `--sukces` | 2 px `--dn-sukces-tekst` | `ptaszek` |
| `--ostrzezenie` | 2 px `--dn-ostrzezenie-tekst` | `ostrzezenie` |
| `--blad` | 2 px `--dn-blad-tekst` | `blad` |
| `--informacja` | 2 px `--dn-informacja-tekst` | `info` |

Toast **nie barwi całego tła** — powierzchnia pozostaje neutralna, barwa ogranicza się do
wstęgi 2 px i do ikony. To jest budżet akcentu w praktyce: powiadomienie o błędzie zajmuje
kilkaset pikseli kwadratowych barwy, a nie kilkanaście tysięcy.

Warstwa: `--dn-z-powiadomienie` = 1000, nad modalem (900), pod dymkiem (1100).
Czas życia domyślny 3200 ms, dla potwierdzenia kopiowania 2000 ms.

### 5.4. `.dn-krok` — 4 warianty

Kolejka kroków realizuje rozstrzygnięcie projektowe zapisane w komentarzu `komponenty.css`:
**trzy wyjścia weryfikacji**.

| Wariant | Wstęga lewa | Znak | Znaczenie |
|---|---|---|---|
| `--pracuje` | `--dn-sygnal-wypelnienie` | obrys sygnałowy | krok w toku |
| `--poprawny` | `--dn-sukces-obrys` | ikona `ptaszek` | weryfikacja pozytywna |
| `--bledy` | `--dn-ostrzezenie-obrys` | ikona `ostrzezenie` + licznik obiegów | błędy odnalezione, kolejka biegnie dalej |
| `--wstrzymany` | `--dn-blad-obrys` | ikona `blad` | zdarzenie nieoczekiwane — kolejka wstrzymana |

Rozróżnienie `--bledy` (bursztyn) od `--wstrzymany` (czerwień) jest istotą modelu:
**znalezione błędy nie są awarią**. Bursztyn oznacza pracę, która wykryła to, po co ją
uruchomiono. Czerwień oznacza, że kolejka nie może biec dalej.

### 5.5. `.dn-postep`

| Element | Żeton | Uwaga |
|---|---|---|
| `.dn-postep-tor` | `--dn-powierzchnia-2`, wysokość 4 px | wgłębienie, nie ramka |
| `.dn-postep-wartosc` | `--dn-sygnal-wypelnienie` | jedyne miejsce, w którym sygnał wypełnia powierzchnię ciągłą |
| `.dn-postep-etykieta` | `--dn-ff-mono`, `--dn-fs-xs` | wartość liczbowa — drugi kanał obok długości paska |

Sterowanie: `data-postep-do` (procent) + `data-postep-czas` (ms). Animacja liczona w
`requestAnimationFrame` z krzywą wyhamowania; przy `prefers-reduced-motion` pasek ustawia się
na wartość docelową **natychmiast** (`prototyp.js` sprawdza `matchMedia` przed animacją).

### 5.6. `.dn-spinner`

Średnica `--dn-wym-spinner` = 14 px, obrys 2 px `--dn-obrys-mocny` z górną krawędzią
`--dn-sygnal-wypelnienie`, obrót 0,8 s liniowo. Stosowany tam, gdzie **nie znamy postępu**:
lista modeli w Model Configuration, rejestr rozszerzeń przy pierwszym otwarciu Skills Manager
i Connectors Manager (Agents 4.2), akcje silnika kolejek w toku (MTAI 12.6).

Rozstrzygnięcie: **spinner tam, gdzie postęp nieznany; `.dn-postep` tam, gdzie znany.**
Nigdy oba naraz.

### 5.7. `.dn-wpis--pracuje`

Modyfikator wpisu okna komunikacji. Dokłada pulsującą kropkę **za nazwą nadawcy**
(pseudoelement `::after` na `.dn-wpis-nadawca`), nie w treści. Znaczenie: ten nadawca
właśnie generuje odpowiedź.

Łączy się z klasą semantyczną nadawcy:

| Klasa | Rola kontraktu | Medalion | Przykład nadawcy |
|---|---|---|---|
| `--czlowiek` | `user` | atrament | Operator |
| `--inteligencja` | `assistant` | sygnał | Executor 1, Executor 2, Coordinator |
| `--system` | `system` / `tool` | neutralny, tło wycofane | silnik kolejek, Process Monitor |

### 5.8. Komunikat blokowy

Dokumentacja przewiduje **komunikat kontekstowy w postaci bloku** — Agents rozdz. 4.2
(„komunikatem, ostrzeżeniem lub dymkiem — albo opisowo obok kontrolki") oraz
MultitaskingAI rozdz. 12.6 („komunikat kontekstowy" przy stanach Ostrzeżenie i Błąd).
Katalog `.dn-*` w kontrakcie systemu projektowego **nie zawiera dla niego klasy**.

Rozstrzygnięcie planszy: komunikat blokowy pokazany jest jako wzorzec lokalny
zbudowany wyłącznie z żetonów rodziny stanu (tło `--dn-{rodzina}-tlo`,
obrys `--dn-{rodzina}-obrys`, tekst `--dn-{rodzina}-tekst`, ikona rodziny, promień `--dn-r-md`).
Propozycja nazwy do katalogu: `.dn-komunikat` + `--sukces --ostrzezenie --blad --informacja`.
Odnotowane w decyzjach projektowych (rozdz. 9).

### 5.9. Tabela doboru nośnika

| Sytuacja | Nośnik | Dlaczego nie inny |
|---|---|---|
| Stan trwały bytu (projekt, ekspert, automatyka) | `.dn-plakietka` | toast znika; plakietka trwa tak długo jak stan |
| Praca biegnąca w tle sesji | `.dn-kropka--tetno` | plakietka „w toku" zajmuje miejsce w każdej karcie; kropka zajmuje 6 px |
| Zdarzenie jednorazowe (zapisano, skopiowano, nie udało się) | `.dn-toast` | plakietka po zdarzeniu nie ma czego opisywać |
| Etap procesu wielokrokowego | `.dn-krok` | plakietki nie pokazują kolejności ani zależności |
| Operacja o znanym postępie | `.dn-postep` | spinner nie niesie informacji ile zostało |
| Operacja o nieznanym postępie | `.dn-spinner` | pasek postępu bez wartości kłamie |
| Nadawca generujący odpowiedź | `.dn-wpis--pracuje` | toast dla każdej odpowiedzi byłby hałasem |
| Problem wymagający wyjaśnienia | komunikat blokowy | plakietka nie zmieści zdania; toast zniknie zanim przeczytany |

---

## 6. Zero blokad jako sygnalizacja

### 6.1. Zasada zero blokad

> **Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód.
> Domyślne zachowanie systemu to wykonanie polecenia.**
> — kontrakt systemu projektowego

Konsekwencje projektowe zapisane w kontrakcie systemu projektowego jako bezwzględne:

| # | Konsekwencja | Sygnalizacja zastępcza |
|---|---|---|
| 1 | Nie stosuje się atrybutu wyłączającego kontrolkę — przycisk jest zawsze klikalny | komunikat po naciśnięciu albo opis obok |
| 2 | Przycisk „Pomiń" w oknie logowania jest obecny i klikalny **w obu fazach wdrożenia** | — |
| 3 | Odliczanie przy „Wyślij ponownie" ma charakter **wyłącznie informacyjny** | tekst odliczania obok przycisku; kliknięcie działa |
| 4 | Metoda uwierzytelniania wyłączona w konfiguracji **nie jest w ogóle renderowana** | „brak metody = mniej segmentów, nie zablokowany segment" |
| 5 | Permissions Center konfiguruje możliwości, nie kontroluje dostępu | stan wyjściowy: pełny dostęp operacyjny |

Ta sama zasada wraca w dokumentacji modułowej dwoma niezależnymi zapisami:

- **Agents rozdz. 4.2**: „*żaden przycisk akcji […] nie występuje w stanie zablokowanym jako
  sposób wymuszenia kompletności czy kolejności. Każdy pozostaje w pełni klikalny w każdej chwili;
  niegotowość lub brak danych sygnalizowana jest dopiero po kliknięciu*".
- **MultitaskingAI rozdz. 12.6**, stan „**Niegotowy (aktywny, z komunikatem)**":
  „*Wygląd tożsamy ze stanem domyślnym, bez przyciemnienia i bez blokady kursora —
  element pozostaje w pełni klikalny*".

### 6.2. Dwa wzorce

| Wzorzec | Kiedy | Realizacja techniczna | Przykład z dokumentacji |
|---|---|---|---|
| **A — komunikat po naciśnięciu** | warunek zależy od stanu, którego Operator może nie znać | `data-komunikat` + `data-komunikat-tytul` + `data-komunikat-rodzaj` w `prototyp.js` → toast | „Scal wyniki" (`merge`) przed zakończeniem wszystkich podagentów (MTAI 12.6) |
| **B — opis obok** | warunek jest stały i przewidywalny; opis oszczędza kliknięcie | tekst `--dn-tekst-2` / `--dn-fs-sm` przy kontrolce, bez zmiany jej wyglądu | odliczanie „Wyślij ponownie" (kontrakt systemu projektowego) |

Wybór między A a B jest decyzją projektową o **koszcie kliknięcia**: jeżeli warunek da się
opisać jednym zdaniem i nie zmienia się w czasie, opis obok jest tańszy dla Operatora.
Jeżeli warunek zależy od stanu procesu (ilu podagentów skończyło, czy zapis się powiódł),
komunikat po naciśnięciu jest jedynym uczciwym rozwiązaniem — opis obok byłby nieaktualny.

### 6.3. Para porównawcza

| | Wzorzec **zakazany** | Wzorzec **właściwy** |
|---|---|---|
| Wygląd | przycisk przyciemniony, kursor zmieniony na „niedozwolone" | przycisk w wyglądzie domyślnym, bez zmian |
| Kliknięcie | brak reakcji | toast z konkretnym powodem i liczbą |
| Klawiatura | element wypada z kolejności fokusu — użytkownik klawiatury nie dowie się, że istnieje | element w kolejności fokusu, komunikat czytany przez `aria-live` kontenera toastów |
| Czytnik ekranu | „przycisk niedostępny", bez powodu | pełna treść komunikatu |
| Czego uczy Operatora | „system mi nie ufa" | „system wie więcej niż ja i mówi mi co" |

Kontener `.dn-toasty` tworzony przez `prototyp.js` ma `role="status"` i `aria-live="polite"` —
komunikat wzorca A jest ogłaszany bez przerywania bieżącej czynności.

### 6.4. Jedyny wyjątek — „nieaktywny kontekstowo"

Agents rozdz. 4.2 dopuszcza jeden stan przygaszony, ściśle ograniczony:

> **Nieaktywny kontekstowo (pole bez zastosowania — nie: brak zgody)**
> `opacity: 0.5`, bez kursora „niedozwolone" — pole pominięte w kolejności fokusu, nie zablokowane.
> Przykład: pole tokenu CLI, dopóki kanał „CLI" nie jest wybrany.

Trzy ograniczenia tego wyjątku, wszystkie zapisane w dokumentacji:

1. Dotyczy **wyłącznie pól danych**, nigdy przycisków akcji.
2. Dotyczy pola **bez zastosowania w bieżącym kontekście**, nigdy pola „czekającego na zgodę".
3. Nie odbiera interakcji jako mechanizm kontroli dostępu (Agents rozdz. 9.1).

### 6.5. Cztery przykłady z dokumentacji

| Kontrolka | Warunek | Wzorzec | Treść komunikatu (wg dokumentacji) |
|---|---|---|---|
| „Scal wyniki" (`merge`) w panelu Subagent Network | nie wszystkie podzadania zakończone | A | wykonanie częściowe z ostrzeżeniem albo komunikat kontekstowy (MTAI 12.6) |
| „Zapisz eksperta" w Agent Builder | pole nazwy puste | A | ostrzeżenie: zapis pod nazwą roboczą; walidacja formy **nigdy nie blokuje zapisu** (Agents 4.1, 4.2) |
| „Wyślij ponownie" w oknie rejestracji | trwa odliczanie | B | odliczanie informacyjne obok przycisku (kontrakt systemu projektowego) |
| „Pomiń" w oknie logowania | faza wdrożenia | — | przycisk obecny i klikalny w obu fazach (kontrakt systemu projektowego) |

---

## 7. Katalog modeli stanów platformy

### 7.1. Tabela zbiorcza — 19 modeli, 105 stanów

| # | Model | Liczba stanów | Stany | Źródło |
|---|---|---|---|---|
| 1 | Stan sesji modułu (Studio) | 5 | Pusta sesja · Edycja · Oczekiwanie na odpowiedź · Podgląd zmiany (Diff/Grep) · Nowa wersja | `moduly/studio.md` 5.1 |
| 2 | Stan sesji modułu (Developer) | 6 | Brak repo · Przegląd · Oczekiwanie na odpowiedź · Podgląd zmiany w edytorze · Zmiana niezatwierdzona · Commit zapisany | `moduly/developer.md` 5.1 |
| 3 | Stan zadania w kolejce | 10 | W kolejce · Odroczone · W realizacji · Podzielone · Scalone · Skierowane · Rozgałęzione · Wstrzymane · Do powtórzenia · Zakończone | `srodowiska/multitaskingai.md` 12.3 |
| 4 | Stan roli / okna roboczego | 6 | Bezczynna · Pracuje · Oczekuje na zależność · Oczekuje na ocenę · Błąd wykonania · Brak przypisania | `srodowiska/multitaskingai.md` 12.2 |
| 5 | Stan przebiegu procesu | 7 | Nieuruchomiony · Uruchomiony · Wstrzymany · Oczekujący na decyzję · Zakończony sukcesem · Zakończony błędem · Zatrzymany trwale | `srodowiska/multitaskingai.md` 12.1 |
| 6 | Stan zlecenia głosowego | 5 | Nasłuch · Rozpoznano · Realizacja wieloetapowa · Zakończone · Wymaga uwagi | `moduly/assistant.md` 5.1 |
| 7 | Stan segmentu tłumaczenia | 4 | Nieprzetłumaczony · Przetłumaczony automatycznie · Zweryfikowany · Wymaga rewizji | `moduly/translate.md` 5.1 |
| 8 | Stan badania | 5 | Zakres zdefiniowany · Zbieranie źródeł · Analiza i synteza · Redakcja raportu · Raport wyeksportowany | `moduly/research.md` 5.1 |
| 9 | Stan rekomendacji | 7 | Zebrany materiał · Analiza w toku · Rekomendacja zaproponowana · Odrzucona · Inna propozycja · Zastosowana · Skuteczność potwierdzona / błąd powrócił | `moduly/diagnostics.md` 5.1 |
| 10 | Stan zasobu wizualnego | 6 | W generowaniu · Nowy · Zaakceptowany · Odrzucony · W kompozycji · Przekazany | `moduly/design.md` 6.2 |
| 11 | Stan produktu | 5 | W architekturze · W budowie · Gotowy do wdrożenia · Wdrożony · Wycofany | `moduly/apps.md` 6.2 |
| 12 | Stan automatyki | 5 | Robocza · Aktywna · Wstrzymana · W toku wykonania · Błędna | `moduly/automations.md` 6.2 |
| 13 | Stan projektu | 4 | Aktywny · Wstrzymany · Zarchiwizowany · Współdzielony (częściowo) | `moduly/workspace.md` 6.2 |
| 14 | Stan eksperta | 5 | Szkic · Niezapisane zmiany · Zapisany / Aktywny · Przypisany · Zarchiwizowany | `okna/agents.md` 4.1 |
| 15 | Stan procesu (Terminal) | 5 | Karta utworzona · Uruchomiony · Zakończony sukcesem · Zakończony błędem · Zakończony przez Operatora | `moduly/terminal.md` 5.1 |
| 16 | Stan repozytorium (Library) | 5 | Plik poza repozytorium · Nieuporządkowany · Skatalogowany · Wersjonowany · Zarchiwizowany | `moduly/library.md` 5.1 |
| 17 | Stan sesji przeglądania (Browser) | 5 | Pusta karta · Przeglądanie · Analiza wspólna z AI · Kolejne źródło dodane automatycznie · Materiał gotowy do przekazania | `moduly/browser.md` 5.1 |
| 18 | Stan sesji Roundtable | 5 | Skład nieustalony · Gotowa do rozpoczęcia · Debata w toku · Generowanie konsensusu · Stanowisko zaakceptowane | `moduly/roundtable.md` 5.1 |
| 19 | Stan Always On Display | 5 | Obserwator/Bezczynny · Obserwator/Doradza · Operator/Bezczynny · Operator/Doradza · Operator/Interweniuje | `srodowiska/multitaskingai.md` 12.5 |

**Razem: 19 modeli, 105 nazwanych stanów.**

Zadanie planszy wymieniało 15 modeli. Cztery pozycje dołożone (2, 17, 18, 19) mają
identyczne umocowanie w dokumentacji — są to modele stanów opisane w rozdziałach 5.1
modułów Developer, Browser i Roundtable oraz w rozdziale 12.5 dokumentu MultitaskingAI.
Model „stan sesji modułu" występuje w dokumentacji **dwukrotnie i z różną liczbą stanów**
(Studio 5, Developer 6), dlatego rozbity jest na dwa wiersze zamiast być uśredniony.

### 7.2. Mapowanie stanów na cztery rodziny

Dokumentacja nazywa stany, lecz przypisania barw dokonuje **wprost tylko dla MultitaskingAI**
(rozdz. 12.1 i 12.2, kolumna „Sygnalizacja wizualna"). Dla pozostałych modeli mapowanie
jest **decyzją projektową tej planszy**, zbudowaną z jednej reguły:

> Rodzinę stanu wyznacza **skutek dla pracy Operatora**, nie ładunek emocjonalny nazwy.
> Praca biegnie → informacja. Praca skończona zgodnie z zamierzeniem → sukces.
> Praca wymaga decyzji lub napotkała problem, który jej nie zatrzymuje → ostrzeżenie.
> Praca nie może biec dalej bez interwencji → błąd. Stan spoczynkowy → neutralny.

| Rodzina | Nośnik | Przykłady stanów z tabeli 7.1 |
|---|---|---|
| **neutralna** (bez barwy stanu) | `.dn-kropka--neutralna`, plakietka bazowa | Pusta sesja · Brak repo · Nieuruchomiony · Bezczynna · Zakres zdefiniowany · Robocza · Szkic · Zarchiwizowany · Karta utworzona · Skład nieustalony |
| **informacja** (= sygnał) | `.dn-kropka` (tętno), `.dn-plakietka--informacja`, `.dn-krok--pracuje` | Oczekiwanie na odpowiedź · W realizacji · Uruchomiony · Pracuje · Realizacja wieloetapowa · Analiza w toku · W generowaniu · W budowie · W toku wykonania · Debata w toku · Oczekujący na decyzję · Oczekuje na ocenę |
| **sukces** | `.dn-plakietka--sukces`, `.dn-krok--poprawny`, `.dn-kropka--sukces` | Nowa wersja · Commit zapisany · Zakończone · Zakończony sukcesem · Zweryfikowany · Raport wyeksportowany · Zaakceptowany · Wdrożony · Aktywna · Aktywny · Zapisany/Aktywny · Skatalogowany · Stanowisko zaakceptowane |
| **ostrzeżenie** | `.dn-plakietka--ostrzezenie`, `.dn-krok--bledy`, `.dn-kropka--ostrzezenie` | Wymaga uwagi · Wymaga rewizji · Do powtórzenia · Wstrzymane · Wstrzymany · Oczekuje na zależność · Niezapisane zmiany · Odrzucony · Wycofany · Nieuporządkowany |
| **błąd** | `.dn-plakietka--blad`, `.dn-krok--wstrzymany`, `.dn-kropka--blad` | Zakończony błędem · Błąd wykonania · Błędna · Zakończony błędem (Terminal) |

Uwaga o „Wstrzymany": dokumentacja MultitaskingAI 12.1 przypisuje mu wprost
**plakietkę ostrzeżenia**, nie błędu — wstrzymanie jest decyzją, nie awarią.
Ta sama logika przenosi się na wstrzymany projekt (Workspace 6.2) i wstrzymaną
automatykę (Automations 6.2).

Uwaga o „Odrzucony" i „Wycofany": nie są awariami — są **decyzjami odwracalnymi**
(wariant zachowany w historii, wdrożenie cofnięte do wcześniejszej wersji).
Bursztyn, nie czerwień.

### 7.3. Arytmetyka katalogu

| Wielkość | Wartość |
|---|---|
| Modele stanów | 19 |
| Nazwane stany łącznie | 105 |
| Najliczniejszy model | stan zadania w kolejce — 10 stanów |
| Najmniej liczne modele | segment tłumaczenia i projekt — po 4 stany |
| Średnia liczba stanów na model | 5,5 |
| Modele z jawnym przypisaniem barwy w dokumentacji | 3 (MTAI 12.1, 12.2, 12.5) |
| Modele opisane diagramem przejść | 11 |
| Modele opisane tabelą | 8 |

Wniosek projektowy: **105 stanów obsługuje 4 rodziny barwne i 8 nośników.**
Ten stosunek jest sednem systemu — gdyby każdy model dostał własną paletę,
platforma miałaby kilkadziesiąt barw i żadnej z nich Operator by nie zapamiętał.
Rodzina jest wspólna, znaczenie precyzuje etykieta.

### 7.4. Konwencja stanów interakcji

Niezależnie od stanów domenowych obowiązuje jednolita konwencja stanów **elementów
interfejsu**, zapisana dwukrotnie: MultitaskingAI 12.6 i Agents 4.2.

| Stan | Reguła wizualna | Żeton |
|---|---|---|
| Domyślny | wygląd spoczynkowy | — |
| Najechanie | subtelna zmiana tła | `--dn-hover` |
| Aktywny / wciśnięty | `translateY(1px)` albo zaznaczenie trwałe | `--dn-wcisniecie` |
| Wybrany | akcent sygnałowy: obrys, podkreślenie, tło | `--dn-sygnal-tlo`, `--dn-sygnal-obrys` |
| Fokus | pierścień 2 px + odsunięcie 2 px, wyłącznie `:focus-visible` | `--dn-fokus`, `--dn-wym-fokus` |
| Ładowanie | `.dn-spinner` albo `.dn-postep` | `--dn-sygnal-wypelnienie` |
| Ostrzeżenie | obramowanie bursztynowe + komunikat; **nie blokuje** | `--dn-ostrzezenie-*` |
| Błąd | obramowanie/tło czerwone + komunikat; wyłącznie awarie techniczne | `--dn-blad-*` |
| Niegotowy (aktywny, z komunikatem) | wygląd tożsamy z domyślnym, pełna klikalność | — |

Dwie pozycje zasługują na podkreślenie, bo odróżniają ten system od domyślnych odruchów:

- **Ostrzeżenie nie blokuje.** Bursztynowe obramowanie pola oznacza „sprawdź to",
  nie „popraw, zanim pozwolę dalej". Walidacja formy jest informacyjna (Agents 4.1).
- **Błąd to awaria techniczna, nie wynik walidacji.** Czerwień rezerwowana jest dla
  utraty połączenia, nieprawidłowego adresu serwera MCP, braku dostępu do mikrofonu —
  zdarzeń poza kontrolą formularza (Agents 4.2, MTAI 12.6).

---

## 8. Budżet akcentu

### 8.1. Reguła

> Zasada jednego akcentu (**sygnał ≤ 5 % ekranu**, nigdy tła sekcji).
> — opracowanie o przekazaniu, ruchu i dostępności, „Czego nie wolno zmienić po cichu"

Reguła zapisana jest w tym samym akapicie co wartości żetonów i geometria znaku —
w kategorii zmian, których nie wolno wprowadzić bez decyzji projektowej. To nie jest
wskazówka stylistyczna, tylko parametr systemu.

### 8.2. Nośniki akcentu i ich koszt powierzchni

Kalkulator planszy liczy udział sygnału w powierzchni ekranu odniesienia
**1280 × 800 px = 1 024 000 px²** (punkt łamania `w3` — biurko, pełny kokpit).

| Nośnik | Geometria | Powierzchnia jednostkowa | Uwaga |
|---|---|---|---|
| Kropka sygnału | koło ⌀ 6 px | ≈ 28,3 px² | `π · 3²` |
| Wstęga aktywnej karty sesji | 2 × 160 px | 320 px² | `--dn-wym-wstega` |
| Pierścień fokusu | obwód kontrolki 120 × 32, grubość 2 px | ≈ 640 px² | jednocześnie zawsze jeden |
| Pasek postępu | 4 × 240 px | 960 px² | tylko wypełniona część |
| Przycisk sygnałowy `.dn-btn--sygnal` | 120 × 32 px | 3 840 px² | najdroższy pojedynczy nośnik |
| Tło pozycji wybranej `--dn-sygnal-tlo` | 224 × 32 px | 7 168 px² | boczna nawigacja |

Kompozycja odniesienia (typowy kokpit z jedną pracą w tle):
3 kropki + 1 wstęga + 1 pierścień fokusu + 1 pasek postępu + 1 przycisk sygnałowy
+ 2 pozycje wybrane = **≈ 20 200 px² ≈ 1,97 %**. Budżet zachowany z zapasem 2,5-krotnym.

### 8.3. Co zjada budżet

| Odruch | Koszt | Dlaczego zakazany |
|---|---|---|
| Sygnał jako tło sekcji | dziesiątki tysięcy px² | jedno użycie zjada cały budżet; kontrakt systemu projektowego: gradient/sygnał nigdy tłem sekcji |
| Sygnał jako barwa przycisku głównego | 3 840 px² × liczba przycisków | działanie główne to **inwersja atramentu**, nie sygnał (kierunek systemu projektowego) |
| Barwne tło całego toastu | ≈ 12 000 px² na powiadomienie | toast barwi wstęgę 2 px i ikonę, nie powierzchnię |
| Dziewięć barw dla dziewięciu nadawców | cała powierzchnia okna komunikacji | anty-domyślne kontrakt systemu projektowego |

Zdanie, które podsumowuje regułę: **czerń działa, sygnał wskazuje.**
Przycisk główny jest czarny na jasnym motywie i biały na ciemnym — monochromatyczna
inwersja jest najsilniejszym możliwym akcentem i **nie wydaje ani jednego piksela**
z budżetu barwy.

### 8.4. Werdykt kalkulatora

| Udział | Werdykt | Nośnik werdyktu |
|---|---|---|
| ≤ 5,00 % | w budżecie | `.dn-plakietka--sukces` + ikona `ptaszek` |
| > 5,00 % | budżet przekroczony | `.dn-plakietka--ostrzezenie` + ikona `ostrzezenie` + toast |

Werdykt nigdy nie jest niesiony samym kolorem paska — towarzyszy mu plakietka z ikoną
i wartość procentowa krojem mono. Ta sama reguła co w całej planszy.

---

## 9. Decyzje projektowe planszy

| # | Decyzja | Uzasadnienie | Umocowanie |
|---|---|---|---|
| 1 | Kontrasty liczone **na żywo w JS**, nie odczytywane z `kontrasty.json` | plansza ma dowodzić, nie cytować; pomiar z faktycznie wyrenderowanych żetonów wyklucza rozejście się dokumentu z kodem | kontrakt systemu projektowego („każda para **zmierzona**") |
| 2 | Pomiar obu motywów jednocześnie przez chwilowe przełączenie `data-theme` | oba motywy są równoprawne, więc tabela musi pokazywać oba naraz; przełączenie nie zapisuje preferencji, bo omija `wspolne.js` | kontrakt systemu projektowego 11 |
| 3 | Próbki barwne obu motywów malowane **wartościami zmierzonymi**, wpisywanymi z JS | pozwala zestawić motyw jasny i ciemny obok siebie na jednej stronie; zero wartości szesnastkowych w `<style>` | kontrakt systemu projektowego |
| 4 | Obrysy stanów mierzone, lecz **bez werdyktu zaliczenia** | obrys jest nośnikiem pomocniczym — stan niesie tekst i ikona; wystawienie progu 4,5 dla obrysu byłoby fałszywym niepowodzeniem | kontrakt systemu projektowego · `kontrasty.json` (próg obrysu 1,6) |
| 5 | Kropka w emblematach pokazana w dwóch wariantach: wiernym (`currentColor`) i z wyróżnieniem | pliki emblematów rysują kropkę w `currentColor`; wyróżnienie służy **wyłącznie dydaktyce planszy** i jest jawnie oznaczone jako podgląd | kontrakt systemu projektowego · pliki `marka/srodowiska/*.svg` |
| 6 | Symulacja `prefers-reduced-motion` jako przełącznik planszy | ustawienie systemowe nie da się przełączyć ze strony; bez symulacji nie da się zestawić tętna i pierścienia statycznego obok siebie | kontrakt systemu projektowego 11 |
| 7 | Symulacja widzenia monochromatycznego filtrem `grayscale(1)` | jedyny sposób pokazania skutku zasady „stan nigdy samym kolorem" bez opisywania go słowami | kontrakt systemu projektowego |
| 8 | Kontrprzykład wyszarzonej kontrolki zbudowany **bez atrybutu wyłączającego** | plansza pokazuje wygląd zakazanego wzorca, sama go nie stosuje; kontrolka pozostaje klikalna i po kliknięciu wyjaśnia, dlaczego ten wzorzec jest zakazany | kontrakt systemu projektowego 14 |
| 9 | Komunikat blokowy jako wzorzec lokalny z propozycją nazwy `.dn-komunikat` | dokumentacja przewiduje komunikat kontekstowy, katalog `.dn-*` nie ma dla niego klasy; luka odnotowana zamiast przemilczana | kontrakt systemu projektowego 10 pkt 5 |
| 10 | Katalog modeli rozszerzony z 15 do 19 pozycji | cztery dołożone modele mają identyczne umocowanie dokumentacyjne; pominięcie ich dałoby obraz niepełny | kontrakt systemu projektowego pkt 5 |
| 11 | „Stan sesji modułu" rozbity na Studio i Developer | dokumentacja opisuje dwa różne modele o różnej liczbie stanów; scalenie ich wymagałoby uśrednienia, czyli wymyślenia | kontrakt systemu projektowego pkt 2 |
| 12 | Mapowanie stanów na rodziny barwne oznaczone jako decyzja planszy | dokumentacja przypisuje barwy wprost tylko dla MultitaskingAI; reszta jest rozstrzygnięciem projektowym i musi być tak nazwana | kontrakt systemu projektowego pkt 5 |
| 13 | Kalkulator budżetu akcentu na ekranie odniesienia 1280 × 800 | punkt łamania `w3` jest zdefiniowany jako „biurko — pełny kokpit"; procent bez powierzchni odniesienia byłby liczbą bez znaczenia | `zetony.css` rozdz. 7 |
| 14 | Dane przykładowe wyłącznie z domeny produktu, jawnie oznaczone | zakaz Lorem ipsum, zmyślonych nazwisk i metryk; identyfikatory zadań i nazwy sesji budowane z realnych nazw okien i modułów | kontrakt systemu projektowego pkt 3, 4 |
| 15 | Wszystkie ikony wklejone inline z realnych plików `zasoby/ikony/svg/` | `currentColor` wymaga inline; jednocześnie gwarantuje zgodność z zestawem (obrys 1,75, siatka 24×24) | kontrakt systemu projektowego 11 pkt 9 |

---

## 10. Źródła

### 10.1. Kanon i kierunek

| Zakres | Plik | Rozdział |
|---|---|---|
| Żetony, rodziny stanów, wymiary, ruch | kontrakt systemu projektowego | 2 |
| Znak marki, kropka sygnału, pięć miejsc | kontrakt systemu projektowego | 3 |
| Ikony — zestaw i zasady | kontrakt systemu projektowego | 4 |
| Katalog komponentów `.dn-*` | kontrakt systemu projektowego | 5 |
| Zero blokad | kontrakt systemu projektowego | 8 |
| Dostępność — warunek wejściowy | kontrakt systemu projektowego | 9 |
| Standard techniczny prototypów | kontrakt systemu projektowego | 11 |
| Barwa: monochrom + jeden sygnał | kierunek systemu projektowego | 3.1 |
| Element sygnaturowy: kropka sygnału | kierunek systemu projektowego | 3.4 |
| Budżet akcentu ≤ 5 % | `design/opracowania/opracowanie o przekazaniu, ruchu i dostępności | 4 |

### 10.2. Pliki systemu wizualnego

| Zakres | Plik |
|---|---|
| Żetony (prymitywy, semantyka, oba motywy, `prefers-reduced-motion`) | `WYNIK/zasoby/zetony/zetony.css` |
| Pomiary kontrastu (33 pary) | `WYNIK/zasoby/zetony/kontrasty.json` |
| Komponenty `.dn-*` | `WYNIK/zasoby/css/komponenty.css` |
| Warstwa prototypu (układ, `.pt-tetno`, panel opracowania) | `WYNIK/zasoby/prototyp.css` |
| Zachowanie prototypu (toasty, widoki, tabele, filtry, `data-komunikat`) | `WYNIK/zasoby/prototyp.js` |
| Przełącznik motywu | `WYNIK/zasoby/wspolne.js` |
| Ikony | `WYNIK/zasoby/ikony/svg/*.svg` |
| Emblematy środowisk | `WYNIK/zasoby/marka/srodowiska/*.svg` |

### 10.3. Modele stanów w dokumentacji produktu

| Model | Plik | Rozdział |
|---|---|---|
| Stan sesji modułu (Studio) | `docs/moduly/studio.md` | 5.1 |
| Stan sesji modułu (Developer) | `docs/moduly/developer.md` | 5.1 |
| Stan zlecenia głosowego | `docs/moduly/assistant.md` | 5.1 |
| Stan segmentu tłumaczenia | `docs/moduly/translate.md` | 5.1 |
| Stan badania | `docs/moduly/research.md` | 5.1 |
| Stan rekomendacji | `docs/moduly/diagnostics.md` | 5.1 |
| Stan repozytorium | `docs/moduly/library.md` | 5.1 |
| Stan procesu (Terminal) | `docs/moduly/terminal.md` | 5.1 |
| Stan sesji przeglądania | `docs/moduly/browser.md` | 5.1 |
| Stan sesji Roundtable | `docs/moduly/roundtable.md` | 5.1 |
| Stan zasobu wizualnego | `docs/moduly/design.md` | 6.2 |
| Stan produktu | `docs/moduly/apps.md` | 6.2 |
| Stan automatyki | `docs/moduly/automations.md` | 6.2 |
| Stan projektu | `docs/moduly/workspace.md` | 6.2 |
| Stan eksperta · stany interfejsu | `docs/moduly/agents.md` | 4.1, 4.2 |
| Stan przebiegu procesu · roli · zadania · AOD · konwencja | `docs/srodowiska/multitaskingai.md` | 12.1–12.6 |
| Jedenaście akcji silnika kolejek | `docs/srodowiska/multitaskingai.md` | 4.2 |
| Cykl życia zadania w kolejce | `docs/srodowiska/multitaskingai.md` | 4.3 |
| Wskaźnik stanu na karcie sesji | `docs/interfejs-uzytkownika/elementy-okien.md` | katalog elementów |

### 10.4. Plansze siostrzane

| Plansza | Plik | Zakres |
|---|---|---|
| 1 | `04-portfolio/01-tozsamosc-srodowisk.html` | cztery środowiska, emblematy, powłoki |
| 2 | `04-portfolio/02-tozsamosc-modulow.html` | piętnaście modułów, okna operacyjne |
| 3 | `04-portfolio/03-role-i-orkiestracja.md` | role MultitaskingAI, silnik kolejek |
| **4** | **`04-portfolio/04-sygnalizacja-i-stany.html`** | **ta plansza** |

---

## 11. Kontrola jakości

| Sprawdzenie | Wynik |
|---|---|
| Nazwy własne zgodne z dokumentacją (zero parafraz) | wszystkie nazwy modułów, okien, ról i stanów cytowane dosłownie |
| Zero elementów bez pokrycia w dokumentacji | każdy stan, nośnik i liczba ma wskazane źródło |
| Zero wartości szesnastkowych w `<style>` | warstwa lokalna planszy wyłącznie `var(--dn-*)` |
| Zero `#000000` | skala kończy się na `--dn-szary-1000` (#0A0A0A) |
| Zero atrybutu wyłączającego kontrolkę | kontrprzykład zbudowany wyglądem, nie atrybutem |
| Zero emoji jako ikon | ikony wklejone inline z `zasoby/ikony/svg/` |
| Zero Lorem ipsum / zmyślonych nazwisk / zmyślonych metryk | dane przykładowe z domeny produktu, oznaczone |
| Oba motywy działają | przełącznik `data-przelacz-motyw` na pasku górnym |
| Pasek górny atramentowy w obu motywach | `--dn-rama` = `--dn-szary-925` |
| Stan nigdy samym kolorem | reguła jest tematem rozdziału 4 i obowiązuje w całej planszy |
| Fokus widoczny na każdej kontrolce | pierścień `--dn-fokus` 2 px + odsunięcie 2 px |
| Interakcje działają z `file://` | wyłącznie `prototyp.js`, `wspolne.js` i skrypt inline |
| Panel „O tym opracowaniu" obecny | `<details class="pt-o-opracowaniu">` na końcu planszy |
| Polskie diakrytyki poprawne | plik w UTF-8, test `ŁÓDŹ ŻÓŁĆ GĘŚLĄ JAŹŃ` |
| Odnośnik powrotu do indeksu | pasek górny i stopka spisu sekcji → `../INDEKS.html` |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
