# Danaco Console — Specyfikacja biblioteki komponentów

| | |
|---|---|
| **Produkt** | **Danaco Console** — AI Operating Environment (warstwa wizualna v2.0) |
| **Rodzaj** | Specyfikacja biblioteki komponentów interfejsu — opracowanie merytoryczno-techniczne (opracowanie 06 · Komponenty) |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-14 |
| **Odbiorcy** | Projektant (co narysować, w jakiej formie, w jakich stanach) · Deweloper (jaką klasę zaimplementować, z jakimi wariantami i zachowaniem) · Redaktor dokumentacji (jak nazywać komponent w tekście) |
| **Zakres** | Wszystkie 33 komponenty biblioteki `.dn-*` zaimplementowane w `zasoby/css/komponenty.css` (1207 linii), uporządkowane w siedem kategorii; karta każdego komponentu (definicja, anatomia, warianty, komplet stanów, wymiary z żetonów, użyte żetony, zachowanie, dostępność, kiedy używać, przykład kodu); macierz komponent × moduł; ściągawka wymiarów; pełna lista klas; zasada zero blokad; luki i rozbieżności między katalogiem v1.0 a implementacją v2.0 |
| **Czego NIE zawiera** | Wartości żetonów (przedmiot opracowania 04 · Żetony) · reguł CSS warstwy zerowej (opracowanie 05 · Style) · makiet i przepływów okien operacyjnych (opracowania zespołu okien) · modelu danych, kontraktów komunikacji i logiki modułów · gotowych plików wykonawczych aplikacji |

---

## Spis treści

1. [Czym jest biblioteka komponentów](#1-czym-jest-biblioteka-komponentów)
2. [Taksonomia — siedem kategorii](#2-taksonomia--siedem-kategorii)
3. [Tabela zbiorcza wszystkich komponentów](#3-tabela-zbiorcza-wszystkich-komponentów)
4. [Konwencja nazewnicza i mechanika stanów](#4-konwencja-nazewnicza-i-mechanika-stanów)
5. [Kategoria A — Przyciski i akcje](#5-kategoria-a--przyciski-i-akcje)
6. [Kategoria B — Pola formularzy](#6-kategoria-b--pola-formularzy)
7. [Kategoria C — Nawigacja](#7-kategoria-c--nawigacja)
8. [Kategoria D — Dane i treść](#8-kategoria-d--dane-i-treść)
9. [Kategoria E — Informacja zwrotna i nakładki](#9-kategoria-e--informacja-zwrotna-i-nakładki)
10. [Kategoria F — Tożsamość](#10-kategoria-f--tożsamość)
11. [Kategoria G — Komunikacja](#11-kategoria-g--komunikacja)
12. [Macierz komponent × moduł](#12-macierz-komponent--moduł)
13. [Zasada zero blokad w komponentach](#13-zasada-zero-blokad-w-komponentach)
14. [Załącznik A — ściągawka wymiarów](#załącznik-a--ściągawka-wymiarów)
15. [Załącznik B — luki i rozbieżności katalog ↔ implementacja](#załącznik-b--luki-i-rozbieżności-katalog--implementacja)
16. [Załącznik C — pełna lista klas CSS](#załącznik-c--pełna-lista-klas-css)
17. [Decyzje projektowe](#17-decyzje-projektowe)

---

## 1. Czym jest biblioteka komponentów

Biblioteka komponentów to **warstwa trzecia** systemu wizualnego Danaco Console: zbiór gotowych, nazwanych klas `.dn-*`, z których składane są wszystkie okna platformy. Komponent jest jedynym miejscem, w którym wartość wizualna spotyka się ze strukturą HTML — nad nim są już tylko wzorce złożone i okna operacyjne, pod nim wyłącznie żetony.

### 1.1. Siedem warstw — od prymitywu do środowiska

```
 WARSTWA 1 · PRYMITYWY              --dn-szary-900 · --dn-sygnal-500 · --dn-zielen-700
   surowe skale barwne                   ZAKAZ użycia wprost w komponencie
        │
        ▼
 WARSTWA 2 · ŻETONY SEMANTYCZNE     --dn-tlo · --dn-powierzchnia · --dn-tekst
   rola zależna od motywu               --dn-sygnal-wypelnienie · --dn-blad-obrys
        │                                zasoby/zetony/zetony.css
        ▼
 WARSTWA 3 · KOMPONENTY .dn-*       .dn-btn · .dn-pole · .dn-tabela · .dn-modal
   33 klasy bazowe + modyfikatory       zasoby/css/komponenty.css  ← NINIEJSZY DOKUMENT
        │
        ▼
 WARSTWA 4 · WZORCE ZŁOŻONE         menu kontekstowe · panel · plansza
   złożenie komponentów warstwy 3       bez własnej klasy
        │
        ▼
 WARSTWA 5 · OKNO OPERACYJNE        Studio Editor · Workflow Builder · Chat Window
        │
        ▼
 WARSTWA 6 · MODUŁ                  Studio · Automations · Agents  (15 modułów)
        │
        ▼
 WARSTWA 7 · ŚRODOWISKO             TalkIn · WorkSpace · CodeStudio · MultitaskingAI
```

**Reguła łańcucha:** każda warstwa odwołuje się wyłącznie do warstwy bezpośrednio niższej. Okno nigdy nie koduje barwy — sięga po komponent. Komponent nigdy nie koduje wartości — sięga po żeton semantyczny. Żeton semantyczny nigdy nie jest wartością przypadkową — wywodzi się z prymitywu.

### 1.2. Pięć zasad wiążących bibliotekę

| # | Zasada | Zapis w pliku źródłowym |
|---|---|---|
| 1 | **Wyłącznie żetony semantyczne.** Zero wartości szesnastkowych zaszytych w komponencie | Nagłówek `komponenty.css`: „Wyłącznie żetony semantyczne z `zetony.css`. Zero wartości zaszytych." |
| 2 | **Komplet stanów.** Każdy komponent interaktywny ma stany: spoczynek · najechanie · naciśnięcie · fokus · wybrany · ładowanie · błąd | Nagłówek `komponenty.css` |
| 3 | **Zero blokad.** Żaden wariant nie odbiera klikalności — niegotowość komunikuje się opisem albo komunikatem | Nagłówek `komponenty.css`; kontrakt systemu projektowego |
| 4 | **Stan nigdy samym kolorem.** Komponent stanu niesie ikonę albo etykietę | Nagłówek `komponenty.css`; kontrakt systemu projektowego |
| 5 | **Nazwy klas są wiążące.** Nie wolno zmieniać nazw istniejących klas; wolno rozszerzać o nowe modyfikatory | kontrakt systemu projektowego |

### 1.3. Trzy poziomy zapisu klasy

| Poziom | Zapis | Przykład | Znaczenie |
|---|---|---|---|
| Klasa bazowa | `.dn-<rzeczownik>` | `.dn-btn`, `.dn-karta` | Komponent w postaci domyślnej — działa bez żadnego modyfikatora |
| Modyfikator | `.dn-<baza>--<cecha>` | `.dn-btn--atrament` | Wariant wyglądu, rozmiaru albo klasy semantycznej; zawsze łącznie z bazą |
| Element wewnętrzny | `.dn-<baza>-<część>` | `.dn-karta-tytul` | Część składowa komponentu; pojedynczy dywiz, nie podwójny |

Rozróżnienie jednego dywizu od dwóch jest kluczowe przy czytaniu arkusza: `.dn-karta-sesji` to **osobny komponent** (karta w pasie sesji), a `.dn-karta--wybrana` to **modyfikator** karty. Ta sama zasada odróżnia `.dn-btn-ikona` (osobny komponent) od `.dn-btn--sm` (modyfikator rozmiaru).

---

## 2. Taksonomia — siedem kategorii

```
BIBLIOTEKA KOMPONENTÓW DANACO CONSOLE · 33 komponenty · 7 kategorii
│
├── A. PRZYCISKI I AKCJE                                        2 komponenty
│     .dn-btn (+7 wariantów, +1 modyfikator stanu)
│     .dn-btn-ikona (+1 wariant)
│
├── B. POLA FORMULARZY                                          7 komponentów
│     .dn-pole (4 elementy)     .dn-wybor      .dn-check
│     .dn-radio                 .dn-przelacznik
│     .dn-suwak                 .dn-szukaj
│
├── C. NAWIGACJA                                                6 komponentów
│     .dn-pasek (4 elementy)    .dn-boczna (2 elementy)
│     .dn-karty-sesji / .dn-karta-sesji
│     .dn-zakladki / .dn-zakladka
│     .dn-listwa (1 element)    .dn-przybornik
│
├── D. DANE I TREŚĆ                                            10 komponentów
│     .dn-karta (+2, 3 elementy)   .dn-karta-srodowiska (4 elementy)
│     .dn-kafel (3 elementy)       .dn-tabela
│     .dn-dane                     .dn-plakietka (+6)
│     .dn-kropka (+5)              .dn-pusty-stan (2 elementy)
│     .dn-postep (3 elementy)      .dn-kolejka / .dn-krok (+4, 2 elementy)
│
├── E. INFORMACJA ZWROTNA I NAKŁADKI                            4 komponenty
│     .dn-modal (4 elementy)       .dn-toasty / .dn-toast (+4, 2 elementy)
│     .dn-tooltip (1 element)      .dn-spinner
│
├── F. TOŻSAMOŚĆ                                                2 komponenty
│     .dn-awatar (+4, 1 element)   .dn-aod (2 elementy)
│
└── G. KOMUNIKACJA                                              2 komponenty
      .dn-wpis (+4, 5 elementów)   .dn-prompt (2 elementy)
```

**Kryterium podziału.** Kategoria odpowiada na pytanie **co komponent robi w interfejsie**, nie w którym module występuje. Ten sam `.dn-plakietka` oznacza status przebiegu w Execution Monitor i etykietę zasobu w Library Explorer — kategoria pozostaje ta sama (D — dane i treść), bo funkcja jest ta sama: opisać jednostkę treści.

**Rozkład masy wizualnej w kategoriach.**

| Kategoria | Waga wizualna komponentów | Rola w kokpicie |
|---|---|---|
| A | znikoma (32 px) | sprawczość — tu Operator działa |
| B | znikoma do średniej | wprowadzanie danych i konfiguracja |
| C | duża (48 / 224 / 36 px, pełne krawędzie) | rama kokpitu — stała, nie przełącza się z treścią |
| D | zmienna od znikomej (6 px kropka) po dużą (karta środowiska) | treść operacyjna — zdecydowana większość powierzchni |
| E | znikoma do dużej (560 px modal) | reakcja systemu — pojawia się i znika |
| F | znikoma do małej (24–36 px) | kto jest kim — odpowiedź na pytanie „czyja to praca" |
| G | duża (860 px wpis, 320 px pas komunikacji) | rozmowa — komponent wspólny wszystkim modułom |

---

## 3. Tabela zbiorcza wszystkich komponentów

| # | Nazwa produktowa (PL) | Klasa `.dn-*` | Kategoria | Warianty | Stany | Wymiar bazowy | Gdzie występuje |
|---|---|---|---|---|---|---|---|
| **1** | Przycisk | `.dn-btn` | Przyciski i akcje | `--atrament` `--sygnal` `--zarys` `--duch` `--niebezpieczny` `--sm` `--lg` (+`--wybrany` jako stan) | spoczynek · najechanie · naciśnięcie · fokus · wybrany · ładowanie · błąd | wys. 32 px, padding 0/16, promień 6 | Wszystkie okna operacyjne, stopki modali, stopki kart, wiersze tabel |
| **2** | Przycisk ikonowy | `.dn-btn-ikona` | Przyciski i akcje | `--na-ramie` | spoczynek · najechanie · naciśnięcie · fokus · wybrany | 32×32 px, ikona 16–20 | Pasek górny, zamknięcie karty sesji, akcje wierszy tabel, przybornik promptu |
| **3** | Pole formularza | `.dn-pole` (`-etykieta` `-kontrolka` `-opis` `-blad`) | Pola formularzy | jednowierszowe · wieloliniowe (`textarea`) · rozwijane (`select`) | spoczynek · najechanie · fokus · błąd · tylko do odczytu · podpowiedź | wys. 32 px (textarea min. 64) | Agent Builder, Instructions Panel, Source Panel, Prompt Builder, Okno Ustawień |
| **4** | Wiersz wyboru | `.dn-wybor` | Pola formularzy | — | dziedziczy stany kontrolki wewnętrznej | min. wys. 32 px, odstęp 8 | Okno Konfiguracji, Permissions Center, Connectors Manager |
| **5** | Pole wyboru | `.dn-check` | Pola formularzy | — | niezaznaczone · zaznaczone · najechanie · fokus | 16×16 px (dotyk 20) | Wybór warstwy konfiguracji, zaznaczenia wielokrotne w Sources Manager |
| **6** | Opcja jednokrotna | `.dn-radio` | Pola formularzy | — | niewybrane · wybrane · najechanie · fokus | 16×16 px, kropka 8 | Wybór warstwy: domyślna / sesji (Okno Konfiguracji) |
| **7** | Przełącznik | `.dn-przelacznik` | Pola formularzy | — | wyłączony · włączony · najechanie · fokus · odziedziczony | tor 36×20 (dotyk 44×24) | Macierz izolacji (11 przełączników na poziom), Permissions Center, Skills Manager |
| **8** | Suwak zakresu | `.dn-suwak` | Pola formularzy | — | spoczynek · przeciąganie · fokus | tor 4 px, kciuk 16 | Nakład rozumowania (szybciej ↔ mądrzej) — Model Configuration |
| **9** | Pole wyszukiwania | `.dn-szukaj` | Pola formularzy | — | spoczynek · najechanie · fokus · z treścią | wys. 32 px, wcięcie 32 | Library Explorer, Sources Manager, Agent Manager, Glossary Manager |
| **10** | Pasek górny | `.dn-pasek` (`-godlo` `-logotyp` `-szukaj` `-prawa`) | Nawigacja | — (kontener) | statyczny; stany należą do elementów wewnątrz | wys. 48 px (przestronna 56) | Powłoka każdego z czterech środowisk |
| **11** | Boczna nawigacja | `.dn-boczna` (`-naglowek` `-pozycja`) | Nawigacja | lista modułów · panel orkiestracji | spoczynek · najechanie · bieżąca · fokus | szer. 224 px, pozycja 32 px | TalkIn (9), WorkSpace (9), CodeStudio (8), MultitaskingAI (6 sekcji) |
| **11a** | Belka tytułowa okna | `.dn-belka` (`-marka` `-nazwa` `-tytul` `-okno` `-btn`) | Nawigacja | `--zamknij` na kontrolce zamknięcia | kontrolka: spoczynek · najechanie · fokus | pas 36 px, kontrolka 46 × 36 | Rama każdego okna platformy — jedyny dom godła |
| **11b** | Pasek narzędzi okna | `.dn-narzedzia` (`-grupa` `-srodek` `-szukaj`) + `.dn-nrz-btn` + `.dn-etykietka` | Nawigacja | — | kontrolka: spoczynek · najechanie (uniesienie 1 px) · naciśnięcie · fokus · rozwinięta · dwustanowa | pas 48 px, kontrolka 32 × 32 | Rama każdego okna platformy |
| **12** | Pas kart sesji / karta sesji | `.dn-karty-sesji` / `.dn-karta-sesji` | Nawigacja | — | w tle · aktywna · najechanie · praca w tle · nowa (pusta) | pas 36 px (przestronna 40) | Powłoka każdego środowiska, pod paskiem górnym |
| **13** | Zakładki / zakładka | `.dn-zakladki` / `.dn-zakladka` | Nawigacja | — | spoczynek · najechanie · wybrana · fokus | wys. 32 px, wskaźnik 2 px | Terminal Tabs, Model Panels (Roundtable), Okno Konfiguracji |
| **14** | Listwa ustawień | `.dn-listwa` (`-pozycja`) | Nawigacja | — | spoczynek · najechanie · fokus | pozycja 28 px, promień 10 | Strefa 3 strony głównej — Okno Konfiguracji, Mobile, Always On Display |
| **15** | Przybornik | `.dn-przybornik` | Nawigacja | — | statyczny; stany należą do przycisków wewnątrz | odstęp 4 px, zawijany | Pole promptu Chat Window, paski narzędzi paneli |
| **16** | Karta | `.dn-karta` (`-naglowek` `-tytul` `-cialo`) | Dane i treść | `--klikalna` `--wybrana` | spoczynek · najechanie · naciśnięcie · wybrana · fokus · ładowanie · błąd | promień 10, cień 1 | Session Repository, Sources Manager, Assets Panel, karty ról MultitaskingAI |
| **17** | Karta środowiska | `.dn-karta-srodowiska` (`-godlo` `-tytul` `-motto` `-opis`) | Dane i treść | — | spoczynek · najechanie · bieżąca · fokus | padding 24, promień 14, godło 40 | Strefa 1 strony głównej — cztery karty środowisk |
| **18** | Kafel komponentu własnego | `.dn-kafel` (`-ikona` `-etykieta` `-opis`) | Dane i treść | — | spoczynek · najechanie · fokus | padding 16, ikona 36 | Strefa 2 strony głównej — Automations, Agents, Workspace, Assistant |
| **19** | Tabela | `.dn-tabela` | Dane i treść | — | wiersz spoczynek · najechanie · wybrany · pusta | wiersz 36 px (dotyk 44) | Execution Monitor, Queue Manager, Process Monitor, Permissions Center |
| **20** | Komórka danych | `.dn-dane` | Dane i treść | — | statyczna | krój mono, stopień 12 | Każda kolumna liczbowa i identyfikatorowa tabel |
| **21** | Plakietka | `.dn-plakietka` | Dane i treść | `--sukces` `--ostrzezenie` `--blad` `--informacja` `--sygnal` `--rola` | statyczna (nośnik stanu innej jednostki) | wys. ok. 18 px, promień pełny | Kolumny statusu monitorów, Tags & Collections, plakietki ról |
| **22** | Kropka sygnału | `.dn-kropka` | Dane i treść | `--tetno` `--sukces` `--ostrzezenie` `--blad` `--neutralna` | statyczna · tętno (2,4 s) · pierścień statyczny (ograniczony ruch) | 6×6 px | Karty sesji, wpisy komunikacji, monitory, godło |
| **23** | Pusty stan | `.dn-pusty-stan` (`-tytul` `-opis`) | Dane i treść | — | statyczny (sam jest stanem panelu) | padding 40/20, ikona 28 | Sources Manager, Findings Panel, Queue Manager przed pierwszą konfiguracją |
| **24** | Pasek postępu | `.dn-postep` (`-etykieta` `-tor` `-wartosc`) | Dane i treść | — | 0% · w toku · 100% · nieokreślony | tor 4 px | Execution Monitor, Build Output, Deployment Panel |
| **25** | Kolejka / krok | `.dn-kolejka` / `.dn-krok` (`-znak` `-meta`) | Dane i treść | `--pracuje` `--poprawny` `--bledy` `--wstrzymany` | oczekuje · pracuje · poprawny · błędy · wstrzymany | krok 36 px, znak 20 | Queue Manager, Orchestrator, Workflow Builder, Results Analyzer |
| **26** | Modal | `.dn-modal` (`-naglowek` `-tytul` `-cialo` `-stopka`) | Informacja zwrotna | — | zamknięty · otwierany · otwarty · ładowanie · błąd | maks. 560 px, maks. 80 dvh | Pula kont Code CLI, potwierdzenia, kreatory strefy 2 |
| **27** | Powiadomienie | `.dn-toasty` / `.dn-toast` (`-tytul` `-tresc`) | Informacja zwrotna | `--sukces` `--ostrzezenie` `--blad` `--informacja` | pojawienie · widoczny · znikanie | 280–420 px, prawy dolny róg | Zapis profilu izolacji, zakończenie automatyki, komunikaty Mobile |
| **28** | Dymek objaśnienia | `.dn-tooltip` (`-tresc`) | Informacja zwrotna | — | ukryty · widoczny (najechanie / fokus) | maks. 260 px, odsunięcie 8 | Oznaczenia `[?]` w Oknie Konfiguracji, ikony paska górnego |
| **29** | Wskaźnik pracy | `.dn-spinner` | Informacja zwrotna | — | obecny · nieobecny · statyczny (ograniczony ruch) | 14×14 px, obrót 0,8 s | Wnętrze przycisku, Chat Window, Build Output, panele w trakcie pobierania |
| **30** | Awatar | `.dn-awatar` (`-stan`) | Tożsamość | `--sm` `--lg` `--kwadrat` `--inteligencja` | domyślny · ze wskaźnikiem stanu · ładowanie obrazu | 28 px (sm 24, lg 36) | Chat Window, Model Panels, karty ról, pasek górny |
| **31** | Always On Display | `.dn-aod` (`-rdzen` `-tresc`) | Tożsamość | — | spoczynek (tętno rdzenia) · rozwinięty · nadzór | rdzeń 36 px, warstwa 1200 | Ponad całą powłoką; sekcja Monitor procesu MultitaskingAI |
| **32** | Wpis komunikacji | `.dn-wpis` (`-medalion` `-nadawca` `-tozsamosc` `-godzina` `-tresc`) | Komunikacja | `--czlowiek` `--inteligencja` `--system` `--pracuje` | spoczynek · pracuje (tętno) | maks. 860 px, medalion 24 | Chat Window — wspólny wszystkim 15 modułom |
| **33** | Pole promptu | `.dn-prompt` (`-grot` `-obszar`) | Komunikacja | — | spoczynek · fokus · z treścią · wysyłanie | obszar 40–160 px | Chat Window, Executor Chat, Coordinator Chat, Voice Console |

**Rachunek:** 33 klasy bazowe · 26 modyfikatorów · 43 elementy wewnętrzne = **102 nazwy klasowe w bibliotece komponentów** (pełny wykaz w Załączniku C; łącznie z 17 klasami warstwy fundamentu — 119 nazw `.dn-*` w arkuszach projektu).

---

## 4. Konwencja nazewnicza i mechanika stanów

### 4.1. Stan wyrażony atrybutem, nie klasą

Biblioteka v2.0 nie zawiera ani jednej klasy stanu w rodzaju `.is-active` czy `.dn-btn--aktywny`. Stan trwały niosą **atrybuty dostępności**, stan chwilowy — **pseudoklasy CSS**. Konsekwencja jest podwójna: czytnik ekranu i arkusz stylów odczytują tę samą prawdę, a stan nie może się „rozjechać" między warstwą wizualną a semantyczną.

| Stan | Nośnik | Selektor w arkuszu | Komponenty |
|---|---|---|---|
| Najechanie | pseudoklasa | `:hover` | `.dn-btn` `.dn-btn-ikona` `.dn-pole-kontrolka` `.dn-karta--klikalna` `.dn-karta-srodowiska` `.dn-kafel` `.dn-zakladka` `.dn-karta-sesji` `.dn-boczna-pozycja` `.dn-listwa-pozycja` `.dn-tabela tbody tr` |
| Naciśnięcie | pseudoklasa | `:active` | `.dn-btn` `.dn-btn-ikona` `.dn-karta--klikalna` |
| Fokus klawiaturowy | pseudoklasa | `:focus-visible` | wszystkie kontrolki (globalnie w `fundament.css`) |
| Fokus pola | pseudoklasa | `:focus`, `:focus-within` | `.dn-pole-kontrolka` `.dn-pasek-szukaj` `.dn-prompt` `.dn-tooltip` |
| Wybrany trwale (przycisk) | atrybut | `[aria-pressed='true']` | `.dn-btn` `.dn-btn-ikona` |
| Wybrany trwale (klasa) | modyfikator | `.dn-btn--wybrany` `.dn-karta--wybrana` | `.dn-btn` `.dn-karta` |
| Wybrany w zestawie | atrybut | `[aria-selected='true']` | `.dn-zakladka` `.dn-karta-sesji` `.dn-tabela tbody tr` |
| Bieżąca pozycja nawigacji | atrybut | `[aria-current='page']`, `[aria-current='true']` | `.dn-boczna-pozycja` `.dn-karta-srodowiska` |
| Ładowanie | atrybut | `[aria-busy='true']` | `.dn-btn` (wskaźnik `::after`) |
| Błąd | atrybut | `[aria-invalid='true']` | `.dn-btn` `.dn-pole-kontrolka` |
| Tylko do odczytu | atrybut | `[readonly]` | `.dn-pole-kontrolka` |
| Zaznaczenie kontrolki | pseudoklasa | `:checked` | `.dn-check` `.dn-radio` `.dn-przelacznik` |
| Otwarcie nakładki | atrybut natywny | `[open]` | `.dn-modal` (element `<dialog>`) |

**Atrybutu `disabled` nie ma w bibliotece ani jednego wystąpienia.** To nie jest przeoczenie — to Zasada zero blokad (rozdz. 13).

### 4.2. Siedem stanów kontraktowych

Nagłówek `komponenty.css` deklaruje siedem stanów obowiązujących każdy komponent interaktywny. Poniżej mapa: który komponent który stan realizuje i czym.

| Komponent | spoczynek | najechanie | naciśnięcie | fokus | wybrany | ładowanie | błąd |
|---|---|---|---|---|---|---|---|
| `.dn-btn` | ● | tło `--dn-hover` | `translateY(1px)` | pierścień 2 px | `aria-pressed` | `aria-busy` + wskaźnik | `aria-invalid` |
| `.dn-btn-ikona` | ● | tło + barwa | `translateY(1px)` | pierścień 2 px | `aria-pressed` | podmiana ikony na `.dn-spinner` | zmiana ikony |
| `.dn-pole-kontrolka` | ● | obrys `--dn-tekst-3` | — | obrys + poświata | `[readonly]` | wskaźnik obok pola | `aria-invalid` + `.dn-pole-blad` |
| `.dn-check` / `.dn-radio` | ● | kursor | — | pierścień | `:checked` | nie dotyczy | komunikat pola nadrzędnego |
| `.dn-przelacznik` | ● | kursor | — | pierścień | `:checked` = sygnał | nie dotyczy | komunikat wiersza |
| `.dn-karta--klikalna` | ● | obrys + cień + uniesienie | `translateY(0)` | pierścień | `--wybrana` | wskaźnik w ciele | plakietka błędu w ciele |
| `.dn-zakladka` | ● | barwa + tło | — | pierścień | `aria-selected` + wskaźnik 2 px | wskaźnik w panelu | plakietka w panelu |
| `.dn-karta-sesji` | ● | tło | — | pierścień | `aria-selected` + wstęga | kropka tętna | plakietka w obszarze roboczym |
| `.dn-tabela tbody tr` | ● | tło `--dn-hover` | — | pierścień w komórce | `aria-selected` | wskaźnik zamiast wierszy | plakietka w kolumnie statusu |
| `.dn-krok` | oczekuje | — | — | pierścień znaku | — | `--pracuje` | `--bledy` / `--wstrzymany` |
| `.dn-modal` | zamknięty | — | — | przechwycenie fokusu | `[open]` | wskaźnik w ciele | komunikat nad stopką |

### 4.3. Trzy klasy semantyczne zamiast dziewięciu barw

Okno komunikacji ma dziewięciu możliwych nadawców. Zamiast dziewięciu barw tła biblioteka wprowadza **trzy klasy semantyczne** i różnicuje rolę wewnątrz klasy ikoną, etykietą i plakietką roli.

| Klasa semantyczna | Modyfikator | Barwa krawędzi | Kto należy |
|---|---|---|---|
| Człowiek | `.dn-wpis--czlowiek` | atrament | Operator |
| Inteligencja | `.dn-wpis--inteligencja` | sygnał | model · agent · Coordinator · Executor 1 · Executor 2 · Executor 3 / Validator |
| System | `.dn-wpis--system` | neutralna | Automation · Always On Display · wynik narzędzia |

Mapowanie ról kontraktu komunikacji: `user → --czlowiek` · `assistant → --inteligencja` · `system → --system` · `tool → --system` z plakietką roli równą nazwie narzędzia.

---

## 5. Kategoria A — Przyciski i akcje

### 5.1. Przycisk — `.dn-btn`

| | |
|---|---|
| **Klasa bazowa** | `.dn-btn` |
| **Kategoria** | A — Przyciski i akcje |
| **Definicja** | Element sprawczy o wysokości kontrolki, wyzwalający pojedynczą, nazwaną czynność systemu |
| **Do czego służy** | Wszędzie tam, gdzie Operator ma **coś zrobić** — zapisać, uruchomić, zatwierdzić, przejść dalej, przypisać |
| **Kiedy używać** | Akcja ma nazwę czasownikową i skutek w systemie; akcja wymaga etykiety tekstowej, bo sama ikona byłaby wieloznaczna |
| **Kiedy NIE używać** | Wybór z listy (→ `.dn-pole` z `select`); przełączenie właściwości dwuwartościowej (→ `.dn-przelacznik`); przełączenie widoku (→ `.dn-zakladka`); akcja jednoznaczna z ikony w gęstym pasku (→ `.dn-btn-ikona`); nawigacja do innego okna (→ `.dn-boczna-pozycja` albo `.dn-karta--klikalna`) |

**Anatomia.** Jeden wiersz `inline-flex`: opcjonalna ikona 16 px → etykieta tekstowa → opcjonalna plakietka lub skrót klawiszowy. Odstęp między częściami `--dn-od-2` (8 px). Przycisk nigdy nie zawija etykiety do dwóch wierszy.

```
┌──────────────────────────────────────┐   min-height 32 px
│  [svg 16]  Uruchom przebieg          │   padding 0 16 px
└──────────────────────────────────────┘   promień 6 px, obrys 1 px
   ▲ ikona     ▲ etykieta (13 px, 500)
```

**Warianty.**

| Wariant | Klasa | Wygląd | Kiedy stosować |
|---|---|---|---|
| Domyślny | `.dn-btn` | Powierzchnia + obrys mocny | Akcja neutralna w panelu, stopka karty |
| Atramentowy (główny) | `--atrament` | Inwersja atramentu: czerń na jasnym, biel na ciemnym | **Akcja podstawowa okna** — jedna na widok |
| Sygnałowy | `--sygnal` | Wypełnienie błękitem sygnałowym, tekst biały | Wyłącznie działanie systemowe „uruchom / zatwierdź plan"; jeden na widok; **nie zastępuje atramentu jako domyślnego głównego** |
| Zarys | `--zarys` | Tło przezroczyste, obrys zachowany | Akcje drugorzędne równorzędne — „Wczytaj profil", „Anuluj" |
| Duch | `--duch` | Bez tła i obrysu, tekst drugorzędny | Akcje trzeciorzędne, paski narzędzi, gęste wiersze |
| Niebezpieczny | `--niebezpieczny` | Obrys i tekst w barwie błędu, tło przezroczyste | Czynności nieodwracalne — sygnalizuje wagę, **nie blokuje** |
| Mały | `--sm` | Wys. 28 px, padding 0/12, stopień 12 | Wiersze tabel, gęste przyborniki |
| Duży | `--lg` | Wys. 40 px, padding 0/20, stopień 14 | Ekrany powitalne, pojedyncza akcja o dużej wadze |

**Stany.**

| Stan | Nośnik | Realizacja wizualna |
|---|---|---|
| Spoczynek | — | Wygląd wg wariantu |
| Najechanie | `:hover` | Tło `--dn-hover`; warianty atramentowy i sygnałowy — przejście na odcień `-hover` |
| Naciśnięcie | `:active` | `transform: translateY(1px)` — wgniecenie, czas `--dn-czas-1` |
| Fokus klawiaturowy | `:focus-visible` | Pierścień `--dn-wym-fokus` (2 px) barwy `--dn-fokus`, odsunięcie 2 px |
| Wybrany trwale | `[aria-pressed='true']` lub `--wybrany` | Tło `--dn-sygnal-tlo`, obrys `--dn-sygnal-obrys`, tekst `--dn-sygnal` |
| Ładowanie | `[aria-busy='true']` | Wskaźnik `::after` 14 px obok etykiety, obrót 0,8 s, krycie 0,7; kursor `progress`; **przycisk pozostaje klikalny** |
| Błąd | `[aria-invalid='true']` | Obrys `--dn-blad-obrys`, tekst `--dn-blad-tekst`; treść przycisku niesie ikonę albo etykietę błędu |
| Warunek niespełniony | — | **Brak stanu wizualnego.** Przycisk wygląda i działa normalnie; po naciśnięciu system zwraca komunikat wskazujący brak |

**Wymiary i żetony.**

| Wymiar | Żeton | Wartość zwarta | Wartość dotykowa / przestronna |
|---|---|---|---|
| Wysokość | `--dn-wym-kontrolka` | 32 px | 40 px |
| Wysokość `--sm` | `calc(--dn-wym-kontrolka - --dn-od-1)` | 28 px | 36 px |
| Wysokość `--lg` | `calc(--dn-wym-kontrolka + --dn-od-2)` | 40 px | 48 px |
| Wcięcie poziome | `--dn-od-4` / `--dn-od-3` / `--dn-od-5` | 16 / 12 / 20 px | bez zmian |
| Odstęp ikona–etykieta | `--dn-od-2` | 8 px | bez zmian |
| Promień | `--dn-r-sm` | 6 px | bez zmian |
| Stopień pisma | `--dn-fs-base` / `-sm` / `-md` | 14 / 13 / 15 px | 15 px w gęstości przestronnej |
| Waga pisma | `--dn-fw-srednia` | 500 | bez zmian |

**Żetony użyte.** `--dn-obrys-mocny` · `--dn-powierzchnia` · `--dn-tekst` · `--dn-hover` · `--dn-atrament` · `--dn-atrament-hover` · `--dn-atrament-tekst` · `--dn-sygnal-wypelnienie` · `--dn-sygnal-wypelnienie-hover` · `--dn-sygnal-tlo` · `--dn-sygnal-obrys` · `--dn-sygnal` · `--dn-tekst-2` · `--dn-blad-obrys` · `--dn-blad-tekst` · `--dn-fokus` · `--dn-czas-1` · `--dn-czas-2` · `--dn-ease` · `--dn-wym-spinner` · `--dn-r-pill`.

**Zachowanie.** Naciśnięcie wyzwala akcję natychmiast; nie ma stanu pośredniego poza ładowaniem. Przy `aria-busy='true'` ponowne naciśnięcie musi być obsłużone idempotentnie w logice akcji — nigdy blokadą kontrolki. Wariant `--niebezpieczny` sygnalizuje wagę skutku; ewentualne potwierdzenie modalem jest wzorcem opcjonalnym dobieranym do konkretnej akcji, nie regułą sztywną.

**Dostępność.**

| Aspekt | Wymaganie |
|---|---|
| Element | `<button type="button">` albo `<a class="dn-btn">` dla przejścia pod adres |
| Rola | domyślna `button`; nie nadpisywać |
| ARIA | `aria-pressed` dla przełącznika trwałego · `aria-busy` dla ładowania · `aria-invalid` dla błędu · `aria-label` gdy etykieta jest wyłącznie ikoną |
| Klawiatura | `Tab` — fokus · `Enter` / `Spacja` — wyzwolenie · fokus zawsze widoczny pierścieniem 2 px |
| Kontrast | wariant atramentowy 16,1:1 (jasny) / 15,0:1 (ciemny); sygnałowy 5,5:1 (jasny) / 4,6:1 (ciemny) — komplet pomiarów w `kontrasty.json` |

**Przykład kodu.**

```html
<!-- Akcja podstawowa okna Workflow Builder (treść przykładowa) -->
<button class="dn-btn dn-btn--atrament" type="button">
  <svg aria-hidden="true" width="16" height="16" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
    <path d="M5 5a2 2 0 0 1 3.008-1.728l11.997 6.998a2 2 0 0 1 .003 3.458l-12 7A2 2 0 0 1 5 19z"/>
  </svg>
  Uruchom przebieg
</button>

<!-- Ten sam przycisk w trakcie pracy — pozostaje klikalny -->
<button class="dn-btn dn-btn--atrament" type="button" aria-busy="true">Uruchamiam przebieg</button>

<!-- Akcja nieodwracalna — sygnalizuje wagę, nie blokuje -->
<button class="dn-btn dn-btn--niebezpieczny" type="button">Usuń profil izolacji</button>
```

---

### 5.2. Przycisk ikonowy — `.dn-btn-ikona`

| | |
|---|---|
| **Klasa bazowa** | `.dn-btn-ikona` |
| **Kategoria** | A — Przyciski i akcje |
| **Definicja** | Kwadratowy element sprawczy o boku kontrolki, niosący wyłącznie ikonę |
| **Do czego służy** | Akcja jednoznaczna z ikony w miejscu, gdzie etykieta tekstowa zajęłaby zbędną przestrzeń |
| **Kiedy używać** | Paski narzędzi, wiersze tabel, nagłówki kart i paneli, pasek górny, przybornik promptu; akcje z zestawu: `zamknij` · `plus` · `wiecej` · `olowek` · `kosz` · `pobierz` · `uruchom` · `zatrzymaj` · `odswiez` |
| **Kiedy NIE używać** | Akcja podstawowa okna (→ `.dn-btn--atrament`); akcja wieloznaczna bez etykiety; ikona bez `aria-label` — taki przycisk nie istnieje dla czytnika ekranu |

**Anatomia.** Kwadrat `--dn-wym-ikonowy` × `--dn-wym-ikonowy` z wyśrodkowaną ikoną 16 px (na pasku górnym 20 px). Obrys przezroczysty w spoczynku — przycisk „pojawia się" dopiero przy najechaniu.

**Warianty.**

| Wariant | Klasa | Wygląd | Gdzie |
|---|---|---|---|
| Domyślny | `.dn-btn-ikona` | Tło przezroczyste, ikona `--dn-tekst-2` | Panele, tabele, karty, przybornik |
| Na ramie | `--na-ramie` | Ikona `--dn-rama-tekst-2`, najechanie `--dn-rama-hover` | **Wyłącznie pasek górny** — powierzchnia atramentowa w obu motywach |

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Tło przezroczyste, ikona drugorzędna |
| Najechanie | `:hover` | Tło `--dn-hover` (na ramie `--dn-rama-hover`), ikona `--dn-tekst` (na ramie `--dn-rama-tekst`) |
| Naciśnięcie | `:active` | `translateY(1px)` |
| Fokus | `:focus-visible` | Pierścień 2 px + odsunięcie 2 px |
| Wybrany | `[aria-pressed='true']` | Tło `--dn-sygnal-tlo`, ikona `--dn-sygnal`; na ramie tło `--dn-rama-hover`, ikona `--dn-sygnal-300` |
| Ładowanie | — | Podmiana ikony na `.dn-spinner` na czas trwania akcji; przycisk pozostaje klikalny |
| Błąd | — | Podmiana ikony (`uruchom` → `blad`) w warstwie ikonografii, nie w warstwie komponentu |

**Wymiary i żetony.** Bok `--dn-wym-ikonowy` = 32 px (dotyk i gęstość przestronna 40 px) · ikona `--dn-wym-ikona` 16 px, na pasku `--dn-wym-ikona-lg` 20 px · promień `--dn-r-sm` 6 px · obrys 1 px przezroczysty · przejścia `--dn-czas-2` / `--dn-czas-1`.

**Żetony użyte.** `--dn-wym-ikonowy` · `--dn-r-sm` · `--dn-tekst-2` · `--dn-tekst` · `--dn-hover` · `--dn-sygnal-tlo` · `--dn-sygnal` · `--dn-rama-tekst` · `--dn-rama-tekst-2` · `--dn-rama-hover` · `--dn-fokus` · `--dn-czas-1` · `--dn-czas-2` · `--dn-ease`.

**Zachowanie.** Naciśnięcie wyzwala akcję natychmiast. W wierszach tabel występuje w zestawach po dwa–trzy (`olowek` edytuj · `kosz` usuń · `pobierz` eksportuj). W pasie kart sesji realizuje kontrolkę zamknięcia (`zamknij`) i kontrolkę nowej karty (`plus`).

**Dostępność.** Zawsze `<button type="button">` z `aria-label` opisującym akcję po polsku („Zamknij kartę sesji", „Odśwież monitor przebiegu"). Ikona zawsze `aria-hidden="true"` — etykieta niesie znaczenie. Przy stanie przełącznika: `aria-pressed`. Dymek objaśnienia (`.dn-tooltip`) jest zalecany wszędzie tam, gdzie ikona nie jest oczywista; nie zastępuje jednak `aria-label`.

**Przykład kodu.**

```html
<!-- Pasek górny — wariant na ramie -->
<button class="dn-btn-ikona dn-btn-ikona--na-ramie" type="button" aria-label="Powiadomienia">
  <svg aria-hidden="true" width="20" height="20" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
    <path d="M10.268 21a2 2 0 0 0 3.464 0"/>
    <path d="M3.262 15.326A1 1 0 0 0 4 17h16a1 1 0 0 0 .74-1.673C19.41 13.956 18 12.499 18 8A6 6 0 0 0 6 8c0 4.499-1.411 5.956-2.738 7.326"/>
  </svg>
</button>

<!-- Wiersz tabeli Queue Manager — zestaw akcji (treść przykładowa) -->
<button class="dn-btn-ikona" type="button" aria-label="Ponów zadanie kolejki">…</button>
<button class="dn-btn-ikona" type="button" aria-pressed="true" aria-label="Przypnij zadanie">…</button>
```

---

## 6. Kategoria B — Pola formularzy

### 6.1. Pole formularza — `.dn-pole`

| | |
|---|---|
| **Klasa bazowa** | `.dn-pole` (kontener) + `.dn-pole-etykieta` · `.dn-pole-kontrolka` · `.dn-pole-opis` · `.dn-pole-blad` |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Kompletna jednostka wprowadzania danych: etykieta, kontrolka, opis pomocniczy i komunikat błędu w jednej kolumnie |
| **Do czego służy** | Przyjęcie tekstu, liczby, adresu, instrukcji systemowej albo wyboru z zamkniętej listy |
| **Kiedy używać** | Zawsze, gdy Operator wpisuje albo wybiera wartość zapisywaną w konfiguracji, profilu agenta, projekcie lub sesji |
| **Kiedy NIE używać** | Wpisywanie polecenia do modelu (→ `.dn-prompt`); wyszukiwanie w obrębie panelu (→ `.dn-szukaj`); wartość dwuwartościowa (→ `.dn-przelacznik`); wartość na skali ciągłej (→ `.dn-suwak`) |

**Anatomia.**

```
.dn-pole                       kolumna, odstęp 4 px
├── .dn-pole-etykieta          12 px, waga 500, --dn-tekst-2
├── .dn-pole-kontrolka         input | textarea | select — wys. 32 px
├── .dn-pole-opis              12 px, --dn-tekst-3 — objaśnienie stałe
└── .dn-pole-blad              12 px, --dn-blad-tekst + ikona — wyłącznie przy błędzie
```

**Warianty.**

| Wariant | Zapis | Charakterystyka |
|---|---|---|
| Jednowierszowy | `<input class="dn-pole-kontrolka">` | Nazwa, adres, wartość krótka; wys. 32 px |
| Wieloliniowy | `<textarea class="dn-pole-kontrolka">` | Opis, instrukcja systemowa; min. wys. 64 px, `resize: vertical` |
| Rozwijany | `<select class="dn-pole-kontrolka">` | Wybór z zamkniętego zbioru; wcięcie prawe 24 px na grot, kursor wskazujący |

Biblioteka v2.0 świadomie **nie mnoży klas** na `input` / `textarea` / `select` — jedna klasa `.dn-pole-kontrolka` obsługuje trzy elementy natywne, a różnicę niesie sam znacznik HTML.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`, tekst `--dn-tekst` |
| Podpowiedź | `::placeholder` | Barwa `--dn-tekst-3`, widoczna wyłącznie przy pustym polu |
| Najechanie | `:hover` | Obrys `--dn-tekst-3` |
| Fokus / edycja | `:focus` | Obrys `--dn-fokus` + poświata `--dn-cien-sygnal` (3 px); własny `outline` wyłączony, bo pierścień zastępuje poświata |
| Błąd | `[aria-invalid='true']` | Obrys `--dn-blad-tekst`; pod polem `.dn-pole-blad` z ikoną — **pole pozostaje w pełni edytowalne** |
| Tylko do odczytu | `[readonly]` | Tło `--dn-powierzchnia-2`, tekst `--dn-tekst-2`; pole nadal fokusowalne i zaznaczalne |
| Walidacja w toku | — | `.dn-spinner` obok pola; pole pozostaje edytowalne |
| Warunek niespełniony | — | Brak stanu wizualnego; ograniczenie opisane w `.dn-pole-opis` albo w dymku |

**Wymiary i żetony.** Wysokość `--dn-wym-kontrolka` 32 px (dotyk / przestronna 40) · `textarea` min. `calc(--dn-wym-kontrolka * 2)` = 64 px · wcięcie `0 --dn-od-3` (12 px), w `textarea` `--dn-od-2 --dn-od-3` · promień `--dn-r-sm` 6 px · stopień `--dn-fs-base` 14 px · etykieta i opis `--dn-fs-sm` 13 px · odstęp w kolumnie `--dn-od-1` 4 px.

**Żetony użyte.** `--dn-obrys-mocny` · `--dn-powierzchnia` · `--dn-powierzchnia-2` · `--dn-tekst` · `--dn-tekst-2` · `--dn-tekst-3` · `--dn-fokus` · `--dn-cien-sygnal` · `--dn-blad-tekst` · `--dn-r-sm` · `--dn-fs-base` · `--dn-fs-sm` · `--dn-fw-srednia` · `--dn-czas-2` · `--dn-ease`.

**Zachowanie.** Kliknięcie albo `Tab` ustawia fokus. Walidacja ma charakter **ostrzegawczy** — pojawienie się `.dn-pole-blad` nie przerywa wpisywania i nie blokuje zapisu formularza. Pole „tylko do odczytu" pozostaje zaznaczalne, żeby wartość dało się skopiować.

**Dostępność.** `<label class="dn-pole-etykieta" for="…">` powiązana z `id` kontrolki — nigdy sam wizualny tekst nad polem. Opis i komunikat błędu wskazane przez `aria-describedby`. Błąd: `aria-invalid="true"` na kontrolce; komunikat w kontenerze `role="alert"` albo `aria-live="polite"`. Klawiatura: `Tab` / `Shift+Tab` między polami, `Enter` zatwierdza formularz, w `select` strzałki zmieniają wartość.

**Przykład kodu.**

```html
<div class="dn-pole">
  <label class="dn-pole-etykieta" for="agent-nazwa">Nazwa agenta</label>
  <input class="dn-pole-kontrolka" id="agent-nazwa" type="text"
         value="Walidator wyników" aria-describedby="agent-nazwa-opis">
  <span class="dn-pole-opis" id="agent-nazwa-opis">Nazwa widoczna w Agent Manager i na kartach ról.</span>
</div>

<div class="dn-pole">
  <label class="dn-pole-etykieta" for="agent-kanal">Kanał wykonania</label>
  <select class="dn-pole-kontrolka" id="agent-kanal">
    <option>API</option><option>CLI</option><option>SSH</option><option>HTTP</option>
  </select>
</div>

<div class="dn-pole">
  <label class="dn-pole-etykieta" for="sciezka-repo">Ścieżka repozytorium</label>
  <input class="dn-pole-kontrolka" id="sciezka-repo" type="text"
         value="D:\danaco\console\ui" aria-invalid="true" aria-describedby="sciezka-repo-blad">
  <span class="dn-pole-blad" id="sciezka-repo-blad" role="alert">
    <svg aria-hidden="true" width="14" height="14" viewBox="0 0 24 24" fill="none"
         stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
      <circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/>
    </svg>
    Ścieżka nie istnieje na tym urządzeniu — zapis jest nadal możliwy.
  </span>
</div>
```

---

### 6.2. Wiersz wyboru — `.dn-wybor`

| | |
|---|---|
| **Klasa bazowa** | `.dn-wybor` |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Poziome opakowanie kontrolki wyboru i jej etykiety, o wysokości kontrolki systemowej |
| **Do czego służy** | Związanie pola wyboru, opcji jednokrotnej albo przełącznika z etykietą w jeden cel kliknięcia |
| **Kiedy używać** | Zawsze, gdy w widoku występuje `.dn-check`, `.dn-radio` albo `.dn-przelacznik` — kontrolka wyboru nigdy nie stoi sama |
| **Kiedy NIE używać** | Wiersz tabeli z przełącznikiem w osobnej kolumnie — tam etykietę niesie komórka, a nie opakowanie |

**Anatomia.** `inline-flex` o `min-height` równym kontrolce: kontrolka (16 px albo 36×20 px) → odstęp 8 px → etykieta 13 px. Kursor wskazujący obejmuje całą szerokość — kliknięcie w etykietę przełącza kontrolkę.

**Warianty.** Brak modyfikatorów. Jedna forma dla trzech kontrolek.

**Stany.** Komponent nie ma stanów własnych — przekazuje je kontrolce wewnętrznej. Najechanie zmienia kursor; fokus dotyczy kontrolki, nie opakowania.

**Wymiary i żetony.** `min-height: --dn-wym-kontrolka` (32 px, dotyk 40) · odstęp `--dn-od-2` 8 px · stopień `--dn-fs-base` 14 px.

**Zachowanie.** Element `<label>` z zagnieżdżoną kontrolką — kliknięcie w dowolne miejsce wiersza przełącza wartość. Nie stosować `for` przy zagnieżdżeniu; przy kontrolce poza etykietą — `for` obowiązkowe.

**Dostępność.** Element `<label>`; przy grupie opcji — `<fieldset>` z `<legend>`; przy grupie przełączników w konfiguracji — `role="group"` z `aria-labelledby` wskazującym nagłówek sekcji.

```html
<label class="dn-wybor">
  <input class="dn-check" type="checkbox" checked>
  Zachowaj historię sesji po zamknięciu karty
</label>

<label class="dn-wybor">
  <input class="dn-przelacznik" type="checkbox" role="switch" aria-checked="false">
  Izolacja kontekstu — pamięć projektu
</label>
```

---

### 6.3. Pole wyboru — `.dn-check`

| | |
|---|---|
| **Klasa bazowa** | `.dn-check` |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Kwadratowa kontrolka zaznaczenia niezależnego, o boku żetonu `--dn-wym-check` |
| **Do czego służy** | Zaznaczenie jednej lub wielu opcji niezależnych od siebie |
| **Kiedy używać** | Zaznaczenia wielokrotne na listach (wybór źródeł do eksportu w Sources Manager), zgody i opcje sesji |
| **Kiedy NIE używać** | Wybór dokładnie jednej opcji z grupy (→ `.dn-radio`); włączenie/wyłączenie właściwości systemu (→ `.dn-przelacznik`) |

**Anatomia.** `appearance: none` — kontrolka rysowana w całości żetonami. Wnętrze: znacznik 10×10 px o kształcie ptaszka (`clip-path`), skalowany od 0 do 1 przy zaznaczeniu.

**Warianty.** Brak modyfikatorów.

**Stany.**

| Stan | Realizacja |
|---|---|
| Niezaznaczone | Obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`, znacznik w skali 0 |
| Zaznaczone | Tło i obrys `--dn-atrament`, znacznik `--dn-atrament-tekst` w skali 1 |
| Najechanie | Kursor wskazujący (reguła w `.dn-wybor`) |
| Fokus | Pierścień 2 px `--dn-fokus` z globalnej reguły `:focus-visible` |
| Nieokreślone (`indeterminate`) | Stan natywny — dopuszczalny dla zaznaczenia częściowego listy; wymaga ustawienia z poziomu skryptu |

**Wymiary i żetony.** Bok `--dn-wym-check` 16 px (dotyk 20) · promień `--dn-r-xs` 3 px · znacznik 10×10 px · przejście `--dn-czas-1` 0,1 s. Żetony: `--dn-obrys-mocny` · `--dn-powierzchnia` · `--dn-atrament` · `--dn-atrament-tekst` · `--dn-r-xs` · `--dn-czas-1` · `--dn-ease`.

**Zachowanie.** Kliknięcie odwraca zaznaczenie niezależnie od pozostałych pól grupy. Zmiana jest natychmiastowa i lokalna.

**Dostępność.** `<input type="checkbox">` — rola natywna. Klawiatura: `Spacja` przełącza. Grupa: `<fieldset>` + `<legend>`. Zaznaczenie nigdy nie jest jedynym nośnikiem znaczenia — etykieta zawsze obecna.

```html
<label class="dn-wybor"><input class="dn-check" type="checkbox"> Eksportuj ustalenia</label>
<label class="dn-wybor"><input class="dn-check" type="checkbox" checked> Eksportuj źródła</label>
```

---

### 6.4. Opcja jednokrotna — `.dn-radio`

| | |
|---|---|
| **Klasa bazowa** | `.dn-radio` (nakładka na `.dn-check`) |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Okrągła kontrolka wyboru dokładnie jednej opcji z grupy wzajemnie wykluczającej się |
| **Do czego służy** | Wybór jednej wartości z krótkiego, zamkniętego zbioru widocznego w całości |
| **Kiedy używać** | Wybór warstwy konfiguracji (domyślna / sesji); wybór trybu, gdy opcji jest dwie do czterech i wszystkie mają być widoczne naraz |
| **Kiedy NIE używać** | Więcej niż cztery opcje (→ `select` w `.dn-pole`); opcje niezależne (→ `.dn-check`); przełączenie widoku (→ `.dn-zakladka`) |

**Anatomia.** Ta sama podstawa co `.dn-check`, z dwiema różnicami: promień `--dn-r-pill` i znacznik w postaci kropki 8×8 px zamiast ptaszka (`clip-path: none`).

**Warianty.** Brak. Różnicę względem pola wyboru niesie atrybut `type="radio"` i wspólna nazwa `name` grupy.

**Stany.** Niewybrane · wybrane (`:checked` — tło atramentowe, kropka w barwie tekstu odwróconego) · najechanie · fokus (pierścień 2 px). Stanu „nieokreślony" nie ma.

**Wymiary i żetony.** Bok `--dn-wym-check` 16 px (dotyk 20) · promień `--dn-r-pill` · kropka 8×8 px. Żetony jak w `.dn-check` plus `--dn-r-pill`.

**Zachowanie.** Kliknięcie wybiera opcję i automatycznie odznacza pozostałe o tej samej nazwie `name`.

**Dostępność.** `<input type="radio" name="…">`. Klawiatura: strzałki przemieszczają wybór wewnątrz grupy, `Tab` wchodzi do grupy i wychodzi z niej (zachowanie natywne — nie nadpisywać). Grupa zawsze w `<fieldset>` z `<legend>` nazywającym pytanie.

```html
<fieldset style="border:none; padding:0;">
  <legend class="dn-pole-etykieta">Warstwa zapisu ustawienia</legend>
  <label class="dn-wybor"><input class="dn-check dn-radio" type="radio" name="warstwa" checked> Warstwa domyślna</label>
  <label class="dn-wybor"><input class="dn-check dn-radio" type="radio" name="warstwa"> Warstwa sesji</label>
</fieldset>
```

---

### 6.5. Przełącznik — `.dn-przelacznik`

| | |
|---|---|
| **Klasa bazowa** | `.dn-przelacznik` |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Kontrolka dwuwartościowa w formie toru z suwakiem, o wymiarach żetonów `--dn-wym-przelacznik-*` |
| **Do czego służy** | Włączenie albo wyłączenie pojedynczej właściwości systemu, ze skutkiem natychmiastowym |
| **Kiedy używać** | Macierz izolacji (11 przełączników na poziom zasięgu — największe skupisko w platformie), Permissions Center, Connectors Manager, Skills Manager, rejestr rozszerzeń |
| **Kiedy NIE używać** | Wartość zatwierdzana dopiero przyciskiem „Zapisz" (→ `.dn-check`); wybór z więcej niż dwóch wartości (→ `.dn-radio` albo `select`) |

**Anatomia.** Tor `36×20 px` (promień pełny) z suwakiem `::before` o boku `calc(wysokość − 6 px)` = 14 px, odsuniętym o 2 px od krawędzi. Po włączeniu suwak przesuwa się do `calc(100% − wysokość + 4px)`.

```
   wyłączony                    włączony
 ┌────────────────┐          ┌────────────────┐
 │ ●              │          │              ○ │   tor 36×20, suwak 14
 └────────────────┘          └────────────────┘
  tło powierzchnia-2          tło sygnał-wypełnienie
  suwak --dn-tekst-2          suwak biel
```

**Warianty.** Brak modyfikatorów — jeden kształt i jeden rozmiar w całej platformie.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Wyłączony (wartość wyjściowa) | — | Tor `--dn-powierzchnia-2`, obrys `--dn-obrys-mocny`, suwak `--dn-tekst-2` po lewej |
| Włączony | `:checked` | Tor i obrys `--dn-sygnal-wypelnienie`, suwak `--dn-szary-0` po prawej |
| Najechanie | `:hover` | Kursor wskazujący — bez odrębnej reguły barwnej |
| Fokus | `:focus-visible` | Pierścień 2 px `--dn-fokus` z reguły globalnej |
| Odziedziczony | plakietka obok | Wartość przejęta z poziomu szerszego — **sygnalizowana wyłącznie plakietką „odziedziczone"**, nigdy `disabled`; przełącznik pozostaje w pełni klikalny, kliknięcie tworzy jawne nadpisanie na tym poziomie |
| Ładowanie | — | Nie dotyczy — zmiana jest natychmiastowa i lokalna; zapis potwierdza powiadomienie |

**Wymiary i żetony.** Tor `--dn-wym-przelacznik-szer` × `--dn-wym-przelacznik-wys` = 36×20 px (dotyk 44×24) · suwak 14 px (dotyk 18) · promień `--dn-r-pill` · przejście `--dn-czas-2` 0,16 s. Żetony: `--dn-powierzchnia-2` · `--dn-obrys-mocny` · `--dn-tekst-2` · `--dn-sygnal-wypelnienie` · `--dn-szary-0` (jedyny dopuszczony prymityw — biel suwaka na wypełnieniu sygnałowym) · `--dn-r-pill` · `--dn-czas-2`.

**Zachowanie.** Kliknięcie w dowolne miejsce toru odwraca wartość natychmiast, bez potwierdzenia pośredniego. Stanem wyjściowym wszystkich ośmiu przełączników izolacji technicznej jest „wyłączony" — zgodnie z zasadą pełnego dostępu na starcie.

**Dostępność.** `<input type="checkbox" role="switch">` z `aria-checked` odzwierciedlającym wartość. Etykieta obowiązkowa (`.dn-wybor` albo komórka tabeli powiązana `aria-labelledby`). Klawiatura: `Spacja` przełącza. **Stan nigdy nie jest komunikowany samą barwą toru** — obok zawsze etykieta „włączone / wyłączone" albo plakietka.

```html
<tr>
  <td id="izo-siec">Izolacja techniczna — dostęp do sieci</td>
  <td><input class="dn-przelacznik" type="checkbox" role="switch"
             aria-checked="false" aria-labelledby="izo-siec"></td>
  <td><span class="dn-plakietka">odziedziczone z poziomu Środowisko</span></td>
</tr>
```

---

### 6.6. Suwak zakresu — `.dn-suwak`

| | |
|---|---|
| **Klasa bazowa** | `.dn-suwak` |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Pozioma kontrolka wartości na skali ciągłej, z torem 4 px i kciukiem obrysowanym sygnałem |
| **Do czego służy** | Ustawienie nakładu rozumowania modelu — skala „szybciej ↔ mądrzej" |
| **Kiedy używać** | Wartość ma sens porządkowy i przybliżony, a Operator myśli o niej jako o natężeniu, nie o liczbie |
| **Kiedy NIE używać** | Wartość musi być wpisana dokładnie (→ `.dn-pole` z `input type="number"`); zbiór wartości jest nazwany i zamknięty (→ `select`) |

**Anatomia.** Tor wypełniany gradientem dwustopniowym sterowanym zmienną `--dn-suwak-pozycja` (procent), kciuk 16×16 px (WebKit) / 12×12 px (Gecko) z obrysem 2 px w barwie sygnału i cieniem `--dn-cien-1`.

**Warianty.** Brak modyfikatorów.

**Stany.** Spoczynek · przeciąganie (kursor wskazujący) · fokus (pierścień z odsunięciem `--dn-od-1`). Stanu błędu nie ma — wartość z zakresu jest zawsze poprawna.

**Wymiary i żetony.** Tor wys. 4 px, promień `--dn-r-pill` · kciuk 16 px / 12 px · odsunięcie pierścienia `--dn-od-1` 4 px. Żetony: `--dn-sygnal-wypelnienie` · `--dn-powierzchnia-2` · `--dn-powierzchnia` · `--dn-cien-1` · `--dn-r-pill` · `--dn-od-1`.

**Zachowanie.** Przeciągnięcie kciuka albo kliknięcie w tor ustawia wartość. Zmienna `--dn-suwak-pozycja` musi być aktualizowana przy każdej zmianie — inaczej wypełnienie toru rozejdzie się z pozycją kciuka. Wartość zawsze pokazywana obok suwaka tekstem, nigdy samą pozycją.

**Dostępność.** `<input type="range">` z `min`, `max`, `step`, `aria-valuetext` podającym znaczenie słowne („nakład wysoki"). Klawiatura: strzałki zmieniają o `step`, `Home` / `End` skaczą do krańców.

```html
<div class="dn-pole">
  <label class="dn-pole-etykieta" for="naklad">Nakład rozumowania</label>
  <input class="dn-suwak" id="naklad" type="range" min="0" max="4" step="1" value="2"
         style="--dn-suwak-pozycja:50%" aria-valuetext="wyważony">
  <span class="dn-pole-opis">szybciej ← <strong>wyważony</strong> → mądrzej</span>
</div>
```

---

### 6.7. Pole wyszukiwania — `.dn-szukaj`

| | |
|---|---|
| **Klasa bazowa** | `.dn-szukaj` (opakowanie ikony i kontrolki) |
| **Kategoria** | B — Pola formularzy |
| **Definicja** | Pole tekstowe z wtopioną ikoną lupy po lewej stronie i powiększonym wcięciem lewym |
| **Do czego służy** | Zawężenie widocznego zbioru pozycji w panelu, tabeli albo bibliotece |
| **Kiedy używać** | Panel z listą dłuższą niż jeden ekran: Library Explorer, Sources Manager, Agent Manager, Glossary Manager, Tags & Collections |
| **Kiedy NIE używać** | Wyszukiwanie globalne na pasku górnym (→ `.dn-pasek-szukaj`, osobna klasa dla powierzchni atramentowej); wprowadzanie wartości zapisywanej (→ `.dn-pole`) |

**Anatomia.** Kontener `position: relative`; ikona `szukaj` 16 px osadzona absolutnie w odległości 12 px od lewej krawędzi, wyłączona z obsługi zdarzeń; kontrolka `.dn-pole-kontrolka` o szerokości 100% i wcięciu lewym `--dn-od-8` (32 px).

**Warianty.** Brak modyfikatorów. Wariantem funkcjonalnym jest wyszukiwanie na ramie — realizowane osobną klasą `.dn-pasek-szukaj`.

**Stany.** Dziedziczy komplet stanów `.dn-pole-kontrolka`: spoczynek · najechanie · fokus (obrys + poświata) · z treścią · tylko do odczytu. Ikona nie zmienia barwy przy fokusie — pozostaje `--dn-tekst-3`, żeby nie konkurować z pierścieniem.

**Wymiary i żetony.** Wysokość 32 px · ikona `--dn-wym-ikona` 16 px w odległości `--dn-od-3` 12 px · wcięcie lewe `--dn-od-8` 32 px. Żetony: jak `.dn-pole-kontrolka` plus `--dn-tekst-3` · `--dn-wym-ikona` · `--dn-od-3` · `--dn-od-8`.

**Zachowanie.** Filtrowanie następuje w trakcie pisania, bez zatwierdzania. Brak wyników zastępuje listę komponentem `.dn-pusty-stan` — nigdy pustą powierzchnią bez komunikatu.

**Dostępność.** `<input type="search">` z etykietą (widoczną albo `.dn-sr-only`). Liczbę wyników ogłasza obszar `aria-live="polite"`. Ikona `aria-hidden="true"`.

```html
<div class="dn-szukaj">
  <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor"
       stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
    <path d="m21 21-4.34-4.34"/><circle cx="11" cy="11" r="8"/>
  </svg>
  <input class="dn-pole-kontrolka" type="search" placeholder="Szukaj w bibliotece…"
         aria-label="Szukaj zasobu w Library Explorer">
</div>
```

---

## 7. Kategoria C — Nawigacja

### 7.1. Pasek górny — `.dn-pasek`

| | |
|---|---|
| **Klasa bazowa** | `.dn-pasek` + `.dn-pasek-godlo` · `.dn-pasek-logotyp` · `.dn-pasek-szukaj` · `.dn-pasek-prawa` |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Stały poziomy pas na szczycie okna, o wysokości żetonu `--dn-wym-pasek`, na powierzchni atramentowej niezależnej od motywu |
| **Do czego służy** | Niesie tożsamość marki, wyszukiwanie globalne i stały dostęp do funkcji narzędziowych |
| **Kiedy używać** | Powłoka każdego z czterech środowisk; nagłówek każdego opracowania dokumentacyjnego |
| **Kiedy NIE używać** | Nagłówek panelu wewnątrz okna (→ `.dn-karta-naglowek`); listwa niskiej masy na stronie głównej (→ `.dn-listwa`) |

**Anatomia.**

```
┌─ .dn-pasek · 48 px · tło --dn-rama (atrament w OBU motywach) ────────────┐
│ [godło ».] DANACO CONSOLE │ [ szukaj ]        [ikona] [ikona] (awatar)  │
│  .dn-pasek-godlo           │ .dn-pasek-szukaj  .dn-pasek-prawa           │
│  + .dn-pasek-logotyp       │ flex:1, maks. 420 px, wys. 30 px            │
└──────────────────────────────────────────────────────────────────────────┘
```

`.dn-pasek-logotyp` składa się z dwóch krojów: `DANACO` krojem nagłówkowym w wadze 700 i `CONSOLE` w `<small>` krojem mono, wersalikami rozstrzelonymi żetonem `--dn-ls-mono-wersaliki` — para krojów opowiada produkt (marka + maszyna).

**Warianty.** Brak modyfikatorów CSS. Odmianą funkcjonalną jest listwa niskiej masy realizowana osobną klasą `.dn-listwa` (rozdz. 7.5).

**Stany.** Pasek jest kontenerem statycznym. Stany należą do elementów wewnątrz: `.dn-btn-ikona--na-ramie`, `.dn-pasek-szukaj` (fokus: obrys `--dn-fokus` + poświata `--dn-cien-sygnal`), `.dn-awatar`.

**Wymiary i żetony.** Wysokość `--dn-wym-pasek` 48 px (gęstość przestronna 56) · wcięcie `0 --dn-od-3` 12 px · odstęp `--dn-od-2` 8 px · pole wyszukiwania: `flex:1`, maks. 420 px, min. wys. 30 px, tło `rgba(255,255,255,.06)`, obrys `--dn-rama-obrys` · logotyp `--dn-fs-md` 15 px, podpis `--dn-fs-xs` 12 px.

**Żetony użyte.** `--dn-rama` · `--dn-rama-tekst` · `--dn-rama-tekst-2` · `--dn-rama-obrys` · `--dn-rama-hover` · `--dn-wym-pasek` · `--dn-od-2` · `--dn-od-3` · `--dn-ff-naglowek` · `--dn-ff-mono` · `--dn-fs-md` · `--dn-fs-xs` · `--dn-fw-gruba` · `--dn-ls-mono-wersaliki` · `--dn-fokus` · `--dn-cien-sygnal` · `--dn-r-sm`.

**Zachowanie.** Pasek nie przewija się z treścią — pozostaje przyklejony u góry (`position: sticky`, warstwa `--dn-z-pasek` = 100). Jest **jedynym elementem interfejsu, który nie przełącza się z motywem**: w motywie jasnym i ciemnym ma tę samą powierzchnię atramentową. Godło prowadzi do Centrum dowodzenia.

**Dostępność.** Element `<header>` albo `role="banner"`. Pole wyszukiwania z etykietą ukrytą `.dn-sr-only`. Kontrast treści na ramie: `--dn-rama-tekst` na `--dn-rama` — 14,6:1; `--dn-rama-tekst-2` na `--dn-rama` — 6,6:1 (pomiary z `kontrasty.json`). Kolejność fokusu: godło → wyszukiwanie → grupa prawa.

```html
<header class="dn-pasek">
  <a class="dn-pasek-godlo" href="#">
    <svg width="28" height="28" viewBox="0 0 96 96" role="img" aria-label="Danaco Console">
      <path fill="currentColor" d="M12 26 H24 L44 48 L24 70 H12 L32 48 Z"/>
      <path fill="currentColor" d="M40 26 H52 L72 48 L52 70 H40 L60 48 Z"/>
      <circle cx="83" cy="63.5" r="6.5" fill="var(--dn-sygnal-400)"/>
    </svg>
    <span class="dn-pasek-logotyp">DANACO<small>CONSOLE</small></span>
  </a>
  <input class="dn-pasek-szukaj" type="search" placeholder="Szukaj w środowisku…"
         aria-label="Wyszukiwanie globalne">
  <div class="dn-pasek-prawa">
    <button class="dn-btn-ikona dn-btn-ikona--na-ramie" type="button" aria-label="Ustawienia">…</button>
    <span class="dn-awatar dn-awatar--sm">OP</span>
  </div>
</header>
```

---

### 7.2. Boczna nawigacja — `.dn-boczna`

| | |
|---|---|
| **Klasa bazowa** | `.dn-boczna` + `.dn-boczna-naglowek` · `.dn-boczna-pozycja` |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Pionowa kolumna o stałej szerokości `--dn-wym-boczna`, zawierająca listę modułów albo sekcji sterowania |
| **Do czego służy** | Wskazuje moduły dostępne w bieżącym środowisku i przełącza obszar roboczy karty sesji |
| **Kiedy używać** | Powłoka każdego środowiska; nawigacja po 13 zakresach Okna Konfiguracji; selektor poziomów zasięgu izolacji |
| **Kiedy NIE używać** | Przełączanie widoku wewnątrz jednego okna (→ `.dn-zakladki`); lista pozycji do wyboru jako dane (→ `.dn-karta--klikalna`) |

**Anatomia.**

```
┌ .dn-boczna · 224 px ─────────┐
│ .dn-boczna-naglowek           │  nazwa środowiska, wersaliki
│ ┌───────────────────────────┐ │
│ │▌[ikona 16] Studio         │ │  .dn-boczna-pozycja[aria-current='page']
│ └───────────────────────────┘ │  ▌ kreska sygnału 2×16 px, promień pełny
│   [ikona 16] Research         │
│   [ikona 16] Library          │
└───────────────────────────────┘
```

**Warianty.**

| Wariant | Zawartość | Środowisko |
|---|---|---|
| Lista modułów | 9 pozycji | TalkIn, WorkSpace |
| Lista modułów | 8 pozycji | CodeStudio |
| Panel orkiestracji | 6 sekcji: Zespoły · Role · Kolejki · Orkiestracja · Harmonogram i automatyki · Monitor procesu | MultitaskingAI |

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Tło przezroczyste, tekst `--dn-tekst-2` |
| Najechanie | `:hover` | Tło `--dn-hover`, tekst `--dn-tekst` |
| Bieżąca | `[aria-current='page']` | Tło `--dn-sygnal-tlo`, tekst `--dn-tekst`, kreska aktywności 2×16 px przy lewej krawędzi |
| Fokus | `:focus-visible` | Pierścień 2 px |
| Nieobecna w środowisku | — | **Pozycja nie jest renderowana wcale** — brak modułu to mniej pozycji, nie pozycja wyszarzona |
| Ładowanie | — | Dotyczy obszaru roboczego, nie nawigacji — wskaźnik pojawia się po prawej stronie |

**Wymiary i żetony.** Szerokość `--dn-wym-boczna` 224 px · wcięcie `--dn-od-3 --dn-od-2` · pozycja: min. wys. `--dn-wym-kontrolka` 32 px, wcięcie `0 --dn-od-3`, odstęp ikona–etykieta `--dn-od-3` 12 px, promień `--dn-r-sm` · ikona `--dn-wym-ikona` 16 px · kreska aktywności `--dn-wym-wstega` 2 px × 16 px. Żetony: `--dn-panel` · `--dn-obrys` · `--dn-tekst-2` · `--dn-tekst` · `--dn-hover` · `--dn-sygnal-tlo` · `--dn-kropka` · `--dn-fs-base` · `--dn-fw-srednia`.

**Zachowanie.** Kliknięcie pozycji przeładowuje obszar roboczy bieżącej karty sesji do układu okien nowego modułu. Sama nawigacja pozostaje niezmieniona. Przy szerokości poniżej punktu `w2` (960 px) kolumna zwija się do ikon — etykiety chowane, szerokość redukowana do kwadratu pozycji.

**Dostępność.** `<nav aria-label="Moduły środowiska TalkIn">` z listą `<ul>`; pozycja jako `<a>` albo `<button>` z `aria-current="page"` dla otwartej. Klawiatura: `Tab` wchodzi do nawigacji, strzałki góra/dół przemieszczają w obrębie listy (obsługa skryptem), `Enter` otwiera. Znacznik aktywności to kreska **plus** tło **plus** `aria-current` — nigdy sam kolor.

```html
<nav class="dn-boczna" aria-label="Moduły środowiska TalkIn">
  <p class="dn-boczna-naglowek dn-etykieta-wersalikowa">TalkIn</p>
  <button class="dn-boczna-pozycja" type="button" aria-current="page">
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
      <path d="M6 22a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h8a2.4 2.4 0 0 1 1.704.706l3.588 3.588A2.4 2.4 0 0 1 20 8v12a2 2 0 0 1-2 2z"/>
      <path d="M14 2v5a1 1 0 0 0 1 1h5"/>
    </svg>
    Studio
  </button>
  <button class="dn-boczna-pozycja" type="button">Research</button>
  <button class="dn-boczna-pozycja" type="button">Library</button>
</nav>
```

---

### 7.3. Pas kart sesji i karta sesji — `.dn-karty-sesji` / `.dn-karta-sesji`

| | |
|---|---|
| **Klasy bazowe** | `.dn-karty-sesji` (pas) · `.dn-karta-sesji` (pojedyncza karta) |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Poziomy pas o wysokości `--dn-wym-pas-kart` z kartami reprezentującymi samodzielne przestrzenie robocze |
| **Do czego służy** | Równoległa praca w kilku modułach albo kilku miejscach tego samego środowiska, z przełączaniem jednym kliknięciem |
| **Kiedy używać** | Wyłącznie w powłoce środowiska, bezpośrednio pod paskiem górnym |
| **Kiedy NIE używać** | Przełączanie widoku wewnątrz okna (→ `.dn-zakladki`); lista sesji jako dane do przeglądania (→ `.dn-karta` w Session Repository) |

**Anatomia.**

```
┌ .dn-karty-sesji · 36 px · tło --dn-powierzchnia-2 ─────────────────────┐
│ ▔▔▔▔▔▔▔▔▔▔▔▔▔▔                                                        │
│ │● Studio      ✕│  │ Research    ✕│  │ Workflow Builder ✕│  │ + │     │
│  ▲ wstęga 2 px      ▲ karta w tle                                       │
│  ▲ kropka tętna = praca w tle                                           │
└─────────────────────────────────────────────────────────────────────────┘
```

Anatomia pojedynczej karty: kropka stanu (opcjonalna) → tytuł modułu → przycisk ikonowy `zamknij`. Na końcu pasa przycisk ikonowy `plus` otwierający nową kartę.

**Warianty.** Brak modyfikatorów — jedna forma różnicowana wyłącznie stanem i obecnością wskaźnika pracy.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| W tle | — | Tło przezroczyste, obrys przezroczysty, tekst `--dn-tekst-2` |
| Najechanie | `:hover` | Tło `--dn-hover`, tekst `--dn-tekst` |
| Aktywna | `[aria-selected='true']` | Tło `--dn-powierzchnia`, obrys `--dn-obrys`, wstęga górna 2 px w barwie `--dn-kropka` |
| Praca w tle | `.dn-kropka--tetno` wewnątrz | Kropka 6 px z tętnem 2,4 s przy tytule |
| Nowa (pusta) | — | Karta bez tytułu modułu — oczekuje wyboru z bocznej nawigacji |
| Fokus | `:focus-visible` | Pierścień 2 px |
| Wyłączona | — | **Nie istnieje** — karta albo jest, albo została zamknięta |

**Wymiary i żetony.** Pas: wys. `--dn-wym-pas-kart` 36 px (przestronna 40), wcięcie `0 --dn-od-2`, odstęp `--dn-od-1`, przewijanie poziome bez widocznego paska · Karta: margines górny `--dn-od-1` 4 px, wcięcie `0 --dn-od-2 0 --dn-od-3`, promień `--dn-r-sm --dn-r-sm 0 0`, stopień `--dn-fs-sm` 13 px, wstęga `--dn-wym-wstega` 2 px.

**Żetony użyte.** `--dn-wym-pas-kart` · `--dn-powierzchnia-2` · `--dn-powierzchnia` · `--dn-obrys` · `--dn-hover` · `--dn-tekst-2` · `--dn-tekst` · `--dn-kropka` · `--dn-wym-wstega` · `--dn-r-sm` · `--dn-fs-sm` · `--dn-czas-2`.

**Zachowanie.**

| Operacja | Wyzwalacz | Skutek |
|---|---|---|
| Otwarcie karty | przycisk `plus` | Nowa, pusta przestrzeń robocza; boczna nawigacja oczekuje wyboru modułu |
| Przełączenie | kliknięcie karty poza `✕` | Obszar roboczy przełącza się na układ, historię i kontekst zapisane w tej karcie |
| Zamknięcie | przycisk `zamknij` | Kończy widok karty; los sesji leżącej u podstaw jest funkcją modelu sesji |
| Przepełnienie pasa | — | Przewijanie poziome (`overflow-x: auto`), pasek przewijania ukryty |

**Dostępność.** Pas jako `role="tablist"` z `aria-label`; karta jako `role="tab"` z `aria-selected`; obszar roboczy jako `role="tabpanel"` powiązany `aria-controls` / `aria-labelledby`. Przycisk zamknięcia jest osobnym `<button>` wewnątrz karty — nie może przechwytywać kliknięcia przełączającego. Klawiatura: strzałki lewo/prawo przemieszczają między kartami, `Delete` zamyka bieżącą (obsługa skryptem).

```html
<div class="dn-karty-sesji" role="tablist" aria-label="Karty sesji środowiska">
  <button class="dn-karta-sesji" type="button" role="tab" aria-selected="true">
    <span class="dn-kropka dn-kropka--tetno" aria-hidden="true"></span>
    Studio
    <span class="dn-sr-only">— w tle biegnie praca</span>
  </button>
  <button class="dn-karta-sesji" type="button" role="tab" aria-selected="false">Research</button>
  <button class="dn-btn-ikona" type="button" aria-label="Otwórz nową kartę sesji">…</button>
</div>
```

---

### 7.4. Zakładki — `.dn-zakladki` / `.dn-zakladka`

| | |
|---|---|
| **Klasy bazowe** | `.dn-zakladki` (grupa) · `.dn-zakladka` (pozycja) |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Poziomy rząd przełączników widoku w obrębie jednego okna, ze wskaźnikiem 2 px pod pozycją wybraną |
| **Do czego służy** | Podmiana zawartości panelu bez zmiany szerszego kontekstu pracy |
| **Kiedy używać** | Terminal Tabs (równoległe sesje powłoki), Model Panels w Roundtable, przełączanie zakresów w Oknie Konfiguracji, wybór warstwy konfiguracji |
| **Kiedy NIE używać** | Przełączanie całej przestrzeni roboczej (→ `.dn-karta-sesji`); nawigacja między modułami (→ `.dn-boczna`); więcej niż siedem pozycji (→ lista pionowa) |

**Anatomia.** Rząd `flex` z odstępem 4 px i kreską dolną 1 px pod całą grupą. Pozycja: opcjonalna ikona 16 px → etykieta → opcjonalna plakietka licznika. Wskaźnik wybranej pozycji to `::after` o wysokości 2 px w barwie `--dn-kropka`, przyklejony do dolnej krawędzi grupy (`bottom: -1px`) — przykrywa kreskę grupy.

**Warianty.** Brak modyfikatorów w implementacji v2.0. Wariant pigułkowy opisany w katalogu v1.0 nie ma odpowiednika (→ Załącznik B, pozycja „Zakładki pigułkowe").

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Tekst `--dn-tekst-2`, tło przezroczyste |
| Najechanie | `:hover` | Tekst `--dn-tekst`, tło `--dn-hover` |
| Wybrana | `[aria-selected='true']` | Tekst `--dn-tekst` + wskaźnik 2 px `--dn-kropka` |
| Fokus | `:focus-visible` | Pierścień 2 px |
| Niedostępna zawartość | — | Zakładka pozostaje klikalna; po przełączeniu panel wyjaśnia stan komunikatem albo pustym stanem |

**Wymiary i żetony.** Pozycja: min. wys. `--dn-wym-kontrolka` 32 px, wcięcie `0 --dn-od-3` 12 px, odstęp `--dn-od-2` 8 px, stopień `--dn-fs-base` 14 px, waga 500 · Grupa: odstęp `--dn-od-1` 4 px, kreska dolna `--dn-obrys` · Wskaźnik 2 px `--dn-kropka`.

**Zachowanie.** Kliknięcie czyni zakładkę wybraną i podmienia zawartość panelu bez przeładowania okna. Zmiana natychmiastowa, przejście barwy w czasie `--dn-czas-2`.

**Dostępność.** `role="tablist"` / `role="tab"` / `role="tabpanel"`; `aria-selected` na pozycjach; `aria-controls` wiążące zakładkę z panelem. Klawiatura: strzałki lewo/prawo, `Home` / `End`; pozycja niewybrana ma `tabindex="-1"`, wybrana `tabindex="0"` (wzorzec pojedynczego punktu wejścia).

```html
<div class="dn-zakladki" role="tablist" aria-label="Sesje powłoki — Terminal Tabs">
  <button class="dn-zakladka" role="tab" aria-selected="true" aria-controls="p-ps" id="z-ps">PowerShell</button>
  <button class="dn-zakladka" role="tab" aria-selected="false" aria-controls="p-cmd" id="z-cmd" tabindex="-1">CMD</button>
  <button class="dn-zakladka" role="tab" aria-selected="false" aria-controls="p-bash" id="z-bash" tabindex="-1">Bash</button>
</div>
<div id="p-ps" role="tabpanel" aria-labelledby="z-ps">…</div>
```

---

### 7.5. Listwa ustawień — `.dn-listwa`

| | |
|---|---|
| **Klasy bazowe** | `.dn-listwa` (kontener) · `.dn-listwa-pozycja` |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Poziomy pasek o najniższej masie wizualnej w systemie, na powierzchni wgłębionej, z pozycjami niższymi od kontrolki |
| **Do czego służy** | Strefa 3 strony głównej — wejścia o najmniejszej wadze: Okno Konfiguracji, Mobile, Always On Display |
| **Kiedy używać** | Zestaw dwóch–czterech wejść pomocniczych, które nie mogą konkurować z kartami środowisk ani kaflami komponentów własnych |
| **Kiedy NIE używać** | Nawigacja główna (→ `.dn-boczna`); akcje w oknie operacyjnym (→ `.dn-btn--duch` w przyborniku) |

**Anatomia.** Kontener z obrysem, promieniem 10 px i tłem `--dn-powierzchnia-2`; wewnątrz pozycje `inline-flex`: ikona 16 px → etykieta 12 px.

**Warianty.** Brak modyfikatorów.

**Stany.** Spoczynek (tekst `--dn-tekst-2`, tło przezroczyste) · najechanie (tło `--dn-hover`, tekst `--dn-tekst`) · fokus (pierścień 2 px). Pozycja nie ma stanu wybranego — listwa prowadzi dalej, nie przełącza widoku w miejscu.

**Wymiary i żetony.** Kontener: wcięcie `--dn-od-2 --dn-od-3`, odstęp `--dn-od-2`, promień `--dn-r-lg` 10 px · Pozycja: min. wys. `calc(--dn-wym-kontrolka − --dn-od-1)` = 28 px, wcięcie `0 --dn-od-3`, promień `--dn-r-sm`, stopień `--dn-fs-sm` 13 px. Żetony: `--dn-powierzchnia-2` · `--dn-obrys` · `--dn-tekst-2` · `--dn-tekst` · `--dn-hover` · `--dn-r-lg` · `--dn-r-sm`.

**Zachowanie.** Kliknięcie otwiera powiązane okno platformowe. Listwa nie zmienia własnego wyglądu po kliknięciu — skutkiem jest zmiana widoku, nie stan pozycji.

**Dostępność.** `<nav aria-label="Ustawienia i funkcje globalne">`; pozycje jako `<button>` albo `<a>` z etykietą tekstową obok ikony (ikona nigdy sama). Klawiatura: `Tab` po pozycjach, `Enter` otwiera.

```html
<nav class="dn-listwa" aria-label="Strefa 3 — ustawienia i funkcje globalne">
  <button class="dn-listwa-pozycja" type="button">
    <svg aria-hidden="true" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/></svg>
    Okno Konfiguracji
  </button>
  <button class="dn-listwa-pozycja" type="button">Mobile</button>
  <button class="dn-listwa-pozycja" type="button">Always On Display</button>
</nav>
```

---

### 7.6. Przybornik — `.dn-przybornik`

| | |
|---|---|
| **Klasa bazowa** | `.dn-przybornik` |
| **Kategoria** | C — Nawigacja |
| **Definicja** | Zawijany rząd narzędzi o odstępie 4 px, bez własnego tła i obrysu |
| **Do czego służy** | Grupuje przyciski ikonowe i plakietki towarzyszące polu promptu: załączniki, biblioteka, agent, model |
| **Kiedy używać** | Pod polem promptu w Chat Window; jako pasek narzędzi kontekstowych panelu |
| **Kiedy NIE używać** | Akcje główne okna (→ stopka karty albo modala); nawigacja (→ `.dn-zakladki`) |

**Anatomia.** Czysty kontener układu: `display: flex`, `align-items: center`, `gap: --dn-od-1`, `flex-wrap: wrap`. Zawartość: `.dn-btn-ikona`, `.dn-btn--duch --sm`, `.dn-plakietka`.

**Warianty.** Brak modyfikatorów.

**Stany.** Kontener bez stanów własnych; komplet stanów należy do elementów wewnątrz. Przy zawężeniu okna przybornik zawija się do drugiego wiersza — nigdy nie chowa narzędzi za przewijaniem poziomym.

**Wymiary i żetony.** Odstęp `--dn-od-1` 4 px. Żetony pozostałe dziedziczone przez zawartość.

**Dostępność.** `role="toolbar"` z `aria-label`, gdy przybornik grupuje narzędzia jednego pola. Klawiatura: strzałki przemieszczają wewnątrz przybornika, `Tab` przeskakuje cały przybornik (wzorzec paska narzędzi).

```html
<div class="dn-przybornik" role="toolbar" aria-label="Narzędzia promptu">
  <button class="dn-btn-ikona" type="button" aria-label="Dodaj załącznik">…</button>
  <button class="dn-btn-ikona" type="button" aria-label="Wstaw z biblioteki">…</button>
  <span class="dn-plakietka dn-plakietka--sygnal">Agent: Walidator wyników</span>
  <span class="dn-plakietka">Nakład: wyważony</span>
</div>
```

---

## 8. Kategoria D — Dane i treść

### 8.1. Karta — `.dn-karta`

| | |
|---|---|
| **Klasa bazowa** | `.dn-karta` + `.dn-karta-naglowek` · `.dn-karta-tytul` · `.dn-karta-cialo` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Powierzchnia z obrysem, promieniem 10 px i cieniem pierwszego stopnia, grupująca powiązaną treść w jedną jednostkę |
| **Do czego służy** | Wydzielenie panelu, pozycji listy albo bloku ustawień jako samodzielnej całości |
| **Kiedy używać** | Panel z nagłówkiem i treścią; pozycja listy w Session Repository, Sources Manager, Assets Panel; karta roli w panelu orkiestracji |
| **Kiedy NIE używać** | Wejście do środowiska (→ `.dn-karta-srodowiska`); kafel komponentu własnego (→ `.dn-kafel`); zbiór jednorodnych wierszy z kolumnami (→ `.dn-tabela`) |

**Anatomia.**

```
.dn-karta
├── .dn-karta-naglowek     wcięcie 12/16, kreska dolna --dn-obrys-subtelny
│   └── .dn-karta-tytul    krój nagłówkowy, 16 px, waga 600
└── .dn-karta-cialo        wcięcie --dn-odstep-panel / --dn-od-4
```

**Warianty.**

| Wariant | Klasa | Charakterystyka | Zastosowanie |
|---|---|---|---|
| Bazowa | `.dn-karta` | Obrys, promień 10 px, cień 1 | Panel bierny, blok treści |
| Klikalna | `--klikalna` | Kursor wskazujący, przy najechaniu obrys mocniejszy, cień 2 i uniesienie 1 px | Pozycja listy prowadząca do szczegółów |
| Wybrana | `--wybrana` | Obrys `--dn-sygnal-obrys`, tło `--dn-sygnal-tlo` | Pozycja zaznaczona w liście, rola wybrana w panelu orkiestracji |

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Powierzchnia + obrys + cień 1 |
| Najechanie | `:hover` (tylko `--klikalna`) | Obrys `--dn-obrys-mocny`, cień 2, `translateY(-1px)` |
| Naciśnięcie | `:active` | `translateY(0)` — karta wraca na miejsce |
| Wybrana | `--wybrana` | Obrys i tło sygnałowe |
| Fokus | `:focus-visible` | Pierścień 2 px z zachowaniem cienia własnego |
| Ładowanie | — | `.dn-spinner` w miejscu treści ciała |
| Błąd | — | `.dn-plakietka--blad` w nagłówku albo komunikat w ciele; karta nie zmienia obrysu |

**Wymiary i żetony.** Promień `--dn-r-lg` 10 px · obrys 1 px `--dn-obrys` · cień `--dn-cien-1`, przy najechaniu `--dn-cien-2` · nagłówek: wcięcie `--dn-od-3 --dn-od-4`, odstęp `--dn-od-2`, kreska `--dn-obrys-subtelny` · tytuł `--dn-fs-lg` 17 px, krój `--dn-ff-naglowek`, waga `--dn-fw-polgruba` 600, odstęp liter `--dn-ls-naglowek` · ciało: wcięcie `--dn-odstep-panel --dn-od-4` (12/16 px; w gęstości przestronnej 20/16).

**Zachowanie.** Karta bazowa jest bierna. Karta klikalna otwiera powiązany widok kliknięciem w dowolne miejsce poza akcjami w nagłówku. Karty w jednym panelu mają jednakową szerokość i zmienną wysokość — nie wyrównuje się ich sztucznie.

**Dostępność.** Karta klikalna musi być `<button>` albo `<a>`, nie `<div onclick>`. Tytuł jako nagłówek właściwego poziomu (`<h2>`…`<h4>`) — `.dn-karta-tytul` stylizuje, nie zastępuje semantyki. Karta wybrana: `aria-pressed` (przełącznik) albo `aria-current` (pozycja bieżąca).

```html
<article class="dn-karta">
  <header class="dn-karta-naglowek">
    <h3 class="dn-karta-tytul">Executor 1</h3>
    <span class="dn-plakietka dn-plakietka--rola">Wykonawca</span>
  </header>
  <div class="dn-karta-cialo">
    <p>Zadania w kolejce: <span class="dn-dane">7</span></p>
  </div>
</article>

<button class="dn-karta dn-karta--klikalna dn-karta--wybrana" type="button" aria-current="true">
  <span class="dn-karta-cialo">Sesja: Workflow Builder — przebieg nocny</span>
</button>
```

---

### 8.2. Karta środowiska — `.dn-karta-srodowiska`

| | |
|---|---|
| **Klasa bazowa** | `.dn-karta-srodowiska` + `-godlo` · `-tytul` · `-motto` · `-opis` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Duża powierzchnia wejściowa o promieniu 14 px, z godłem 40 px, tytułem w stopniu 30 px i wstęgą sygnału ujawnianą przy najechaniu |
| **Do czego służy** | Wejście do jednego z czterech środowisk platformy ze strefy 1 strony głównej |
| **Kiedy używać** | **Wyłącznie** dla czterech środowisk: TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| **Kiedy NIE używać** | Komponent własny strefy 2 (→ `.dn-kafel`); moduł (→ `.dn-boczna-pozycja`); pozycja listy (→ `.dn-karta`) |

**Anatomia.**

```
┌ .dn-karta-srodowiska · promień 14 px · wcięcie 24 px ─────────┐
│▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔▔ wstęga 2 px (krycie 0 → 1)          │
│  [godło 40×40]                                                 │
│  TalkIn                          ← -tytul 30 px, krój nagł.    │
│  rozmowa i wiedza                ← -motto 12 px, krój mono     │
│  Dziewięć modułów pracy z treścią ← -opis 13 px, maks. 44 znaki │
└────────────────────────────────────────────────────────────────┘
```

**Warianty.** Brak modyfikatorów. Różnicę między czterema kartami niesie wyłącznie treść: emblemat środowiska, nazwa własna, motto i opis.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Obrys `--dn-obrys`, wstęga o kryciu 0 — **bez sygnału** |
| Najechanie | `:hover` | Obrys `--dn-obrys-mocny`, cień 2, `translateY(-2px)`, wstęga o kryciu 1 |
| Bieżąca | `[aria-current='true']` | Jak najechanie, utrzymane trwale |
| Fokus | `:focus-visible` | Pierścień 2 px + odsunięcie |
| Ładowanie | — | Wskaźnik w miejscu opisu na czas przygotowania powłoki środowiska |

**Wymiary i żetony.** Wcięcie `--dn-od-6` 24 px · promień `--dn-r-xl` 14 px · godło 40×40 px, margines dolny `--dn-od-4` · tytuł `--dn-fs-3xl` 30 px, waga 600, interlinia ciasna · motto `--dn-fs-sm` 13 px krojem mono, barwa `--dn-tekst-3` · opis `--dn-fs-base` 14 px, `--dn-tekst-2`, maks. 44 znaki w wierszu · wstęga `--dn-wym-wstega` 2 px w barwie `--dn-kropka` · uniesienie 2 px w czasie `--dn-czas-2`.

**Zachowanie.** Kliknięcie otwiera powłokę środowiska i przenosi Operatora do warstwy modułów. Wstęga jest **jedynym miejscem sygnału** na tej karcie — reszta pozostaje monochromatyczna. Karta nie animuje się w spoczynku.

**Dostępność.** Element `<a>` albo `<button>`; tytuł jako nagłówek `<h2>`. Emblemat `aria-hidden="true"` — nazwa środowiska jest w tytule. Cztery karty w strefie 1 tworzą listę `<ul>` z `aria-label="Środowiska platformy"`.

```html
<a class="dn-karta-srodowiska" href="#" aria-current="true">
  <svg class="dn-karta-srodowiska-godlo" aria-hidden="true" viewBox="0 0 24 24" fill="none"
       stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">…</svg>
  <h2 class="dn-karta-srodowiska-tytul">CodeStudio</h2>
  <p class="dn-karta-srodowiska-motto">kod i wytwarzanie</p>
  <p class="dn-karta-srodowiska-opis">Osiem modułów: Terminal, Developer, Diagnostics, Apps, Design, Roundtable, Workspace, Agents.</p>
</a>
```

---

### 8.3. Kafel komponentu własnego — `.dn-kafel`

| | |
|---|---|
| **Klasa bazowa** | `.dn-kafel` + `-ikona` · `-etykieta` · `-opis` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Pozioma powierzchnia na tle panelu, z ikoną w kwadracie 36 px i dwuwierszowym opisem, świadomie lżejsza od karty środowiska |
| **Do czego służy** | Wejście do budowania komponentu własnego ze strefy 2 strony głównej |
| **Kiedy używać** | **Wyłącznie** dla czterech komponentów własnych: Automations, Agents, Workspace, Assistant |
| **Kiedy NIE używać** | Wejście do środowiska (→ `.dn-karta-srodowiska`); pozycja listy w panelu (→ `.dn-karta--klikalna`) |

**Anatomia.** Rząd `flex` z odstępem 12 px: kwadrat ikony 36×36 px (obrys, promień 8 px, powierzchnia) → kolumna z etykietą 13 px w wadze 600 i opisem 12 px w barwie `--dn-tekst-3`.

**Warianty.** Brak modyfikatorów.

**Stany.** Spoczynek (tło `--dn-panel`, obrys `--dn-obrys`) · najechanie (obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia` — kafel „wychodzi" z panelu) · fokus (pierścień 2 px). Kafel nie unosi się jak karta środowiska — różnica masy jest celowa.

**Wymiary i żetony.** Wcięcie `--dn-od-4` 16 px · odstęp `--dn-od-3` 12 px · promień `--dn-r-lg` 10 px · ikona: `--dn-wym-awatar-lg` 36 px, promień `--dn-r-md` 8 px · etykieta `--dn-fs-base` 14 px waga 600 · opis `--dn-fs-sm` 13 px.

**Zachowanie.** Kliknięcie otwiera kreator komponentu własnego — zwykle w postaci modala (Automations, Agents, Workspace, Assistant).

**Dostępność.** `<button>` z tekstem etykiety; ikona `aria-hidden="true"`; opis powiązany `aria-describedby`, gdy jest istotny dla decyzji.

```html
<button class="dn-kafel" type="button" aria-describedby="kafel-auto-opis">
  <span class="dn-kafel-ikona">
    <svg aria-hidden="true" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
      <rect width="8" height="8" x="3" y="3" rx="2"/><path d="M7 11v4a2 2 0 0 0 2 2h4"/>
      <rect width="8" height="8" x="13" y="13" rx="2"/>
    </svg>
  </span>
  <span>
    <span class="dn-kafel-etykieta">Automations</span>
    <span class="dn-kafel-opis" id="kafel-auto-opis">Automatyka wpinana w sesję dowolnego modułu</span>
  </span>
</button>
```

---

### 8.4. Tabela — `.dn-tabela`

| | |
|---|---|
| **Klasa bazowa** | `.dn-tabela` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Zwarta tabela robocza o wysokości wiersza `--dn-wym-wiersz`, z przyklejonym nagłówkiem kolumn krojem mono |
| **Do czego służy** | Prezentacja zbioru jednorodnych pozycji: przebiegów, zadań kolejki, uprawnień, konektorów, zdarzeń |
| **Kiedy używać** | Okna typu monitor i zarządca: Execution Monitor, Queue Manager, Process Monitor, Permissions Center, Connectors Manager, Agent Manager, macierz izolacji |
| **Kiedy NIE używać** | Pozycje o różnej strukturze (→ lista `.dn-karta`); dwie–trzy pozycje (→ lista opisowa); treść ciągła (→ `.dn-karta-cialo`) |

**Anatomia.**

```
┌ .dn-tabela ────────────────────────────────────────────────────┐
│ ZADANIE      MODUŁ        STATUS        CZAS     AKCJE         │ th 32 px, mono 11 px
├────────────────────────────────────────────────────────────────┤   wersaliki, przyklejony
│ QUE-0142     Automations  ● powodzenie  00:04:12 [olowek][kosz]│ td 36 px
│ QUE-0143     Automations  ● w toku      00:00:37 [olowek][kosz]│ wiersz wybrany:
└────────────────────────────────────────────────────────────────┘   tło --dn-sygnal-tlo
```

**Warianty.** Brak modyfikatorów w implementacji v2.0. Klasa pomocnicza `.dn-dane` na komórce ustawia stopień 12 px i barwę drugorzędną dla danych maszynowych (rozdz. 8.5).

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Wiersz spoczynek | — | Kreska dolna `--dn-obrys-subtelny` |
| Wiersz najechanie | `tbody tr:hover` | Tło `--dn-hover`, przejście `--dn-czas-1` |
| Wiersz wybrany | `[aria-selected='true']` | Tło `--dn-sygnal-tlo` |
| Wiersz z błędem | — | `.dn-plakietka--blad` w kolumnie statusu — **struktura tabeli bez zmian** |
| Ładowanie | — | Wskaźnik w miejscu ciała tabeli albo wiersze szkieletowe |
| Pusta | — | `.dn-pusty-stan` w miejsce ciała tabeli |
| Sortowanie | `aria-sort` na `th` | Grot kierunku przy nazwie kolumny |

**Wymiary i żetony.** Nagłówek: wys. `calc(--dn-wym-wiersz − --dn-od-1)` = 32 px, tło `--dn-panel`, krój `--dn-ff-mono`, stopień `--dn-fs-xs` 12 px, wersaliki, odstęp liter `--dn-ls-mono-wersaliki` 0,14 em, barwa `--dn-tekst-3`, przyklejony (`position: sticky`, warstwa `--dn-z-przybornik` = 10) · Komórka: wys. `--dn-wym-wiersz` 36 px (dotyk 44, przestronna 44), wcięcie `0 --dn-od-3`, wyrównanie do środka w pionie · Stopień treści `--dn-fs-base` 14 px.

**Zachowanie.** Nagłówek pozostaje widoczny przy przewijaniu ciała. Wiersze monitorów aktualizują się na żywo, bez działania Operatora. Kolumna akcji zawiera `.dn-btn-ikona`; kolumna właściwości binarnej — `.dn-przelacznik`.

**Dostępność.** Znacznik `<table>` z `<caption>` (może być `.dn-sr-only`), `<thead>`, `<tbody>`, `scope="col"` na nagłówkach. Sortowanie: `aria-sort="ascending|descending|none"` na `th` + przycisk w nagłówku. Wybór wiersza: `aria-selected` na `<tr>` przy `role="row"` w siatce interaktywnej. Wiersze aktualizowane na żywo: obszar `aria-live="polite"` z podsumowaniem, nie ogłaszanie każdego wiersza.

```html
<table class="dn-tabela">
  <caption class="dn-sr-only">Kolejka zadań — Queue Manager (dane przykładowe)</caption>
  <thead>
    <tr>
      <th scope="col" aria-sort="ascending">Zadanie</th>
      <th scope="col">Moduł</th><th scope="col">Status</th><th scope="col">Czas</th>
    </tr>
  </thead>
  <tbody>
    <tr aria-selected="true">
      <td class="dn-dane">QUE-0142</td>
      <td>Automations</td>
      <td><span class="dn-plakietka dn-plakietka--sukces">powodzenie</span></td>
      <td class="dn-dane">00:04:12</td>
    </tr>
  </tbody>
</table>
```

---

### 8.5. Komórka danych — `.dn-dane`

| | |
|---|---|
| **Klasa bazowa** | `.dn-dane` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Klasa pomocnicza nadająca treści krój mono, liczby tabelaryczne i stopień 12 px |
| **Do czego służy** | Odróżnienie danych maszynowych (identyfikatory, czasy, ścieżki, wersje) od treści redakcyjnej |
| **Kiedy używać** | Kolumny identyfikatorów i liczb w tabelach; wartości w kartach i wpisach; wszędzie tam, gdzie liczby mają się wyrównywać w pionie |
| **Kiedy NIE używać** | Tekst zdaniowy; nazwy modułów i okien (te są treścią, nie danymi) |

**Anatomia.** Brak struktury własnej — modyfikuje krój i wariant liczbowy elementu, w którym jest użyta. Definicja rozdziela się na dwa arkusze: `fundament.css` nadaje `--dn-ff-mono` i `font-variant-numeric: tabular-nums`, `komponenty.css` uzupełnia stopień i barwę w kontekście komórki tabeli (`.dn-tabela td.dn-dane`).

**Warianty.** Brak. Klasa bliźniacza `.dn-liczba` (fundament) daje sam krój mono i liczby tabelaryczne, bez zmiany stopnia.

**Stany.** Statyczna.

**Wymiary i żetony.** Stopień `--dn-fs-sm` 13 px w kontekście tabeli · barwa `--dn-tekst-2` · krój `--dn-ff-mono` · `font-variant-numeric: tabular-nums`.

**Zachowanie.** Liczby o równej szerokości znaku nie „skaczą" przy odświeżaniu monitora — to warunek czytelności okien aktualizowanych na żywo.

**Dostępność.** Klasa czysto wizualna. Identyfikatory techniczne wymagają rozwinięcia w `aria-label` albo w `<abbr>` tylko wtedy, gdy skrót nie jest wyjaśniony nagłówkiem kolumny.

```html
<td class="dn-dane">QUE-0142</td>
<span class="dn-dane">00:04:12</span>
<span class="dn-dane">v2.0</span>
```

---

### 8.6. Plakietka — `.dn-plakietka`

| | |
|---|---|
| **Klasa bazowa** | `.dn-plakietka` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Etykieta w linii tekstu o promieniu pełnym, stopniu 11 px, z opcjonalną ikoną 12 px |
| **Do czego służy** | Oznaczenie statusu, kategorii, roli albo etykiety własnej jednostki treści |
| **Kiedy używać** | Kolumna statusu w monitorach; etykiety w Tags & Collections; rola okna w MultitaskingAI; źródło rozszerzenia w rejestrze; oznaczenie „odziedziczone" w macierzy izolacji |
| **Kiedy NIE używać** | Akcja do kliknięcia (→ `.dn-btn--sm`); status bez miejsca na etykietę (→ `.dn-kropka` z etykietą kolumny); długi komunikat (→ treść panelu) |

**Warianty.**

| Wariant | Klasa | Barwy | Znaczenie |
|---|---|---|---|
| Neutralna | `.dn-plakietka` | Tło `--dn-powierzchnia-2`, obrys `--dn-obrys`, tekst `--dn-tekst-2` | Etykieta bez wydźwięku: kategoria, źródło, tag |
| Sukces | `--sukces` | Rodzina `--dn-sukces-*` | Zakończone powodzeniem |
| Ostrzeżenie | `--ostrzezenie` | Rodzina `--dn-ostrzezenie-*` | Wymaga uwagi; zakończone z błędami |
| Błąd | `--blad` | Rodzina `--dn-blad-*` | Niepowodzenie, wstrzymanie |
| Informacja | `--informacja` | Rodzina `--dn-informacja-*` (rodzina sygnału) | Stan neutralny informacyjny |
| Sygnał | `--sygnal` | `--dn-sygnal-tlo` / `-obrys` / `--dn-sygnal` | Wyróżnienie systemowe: aktywny model, przypisany agent |
| Rola | `--rola` | Dziedziczy barwy wariantu; krój mono, wersaliki, odstęp liter 0,14 em, promień 3 px | Rola okna: KOORDYNATOR · WYKONAWCA · WALIDATOR |

**Stany.** Plakietka jest **nośnikiem stanu innej jednostki**, nie ma stanów interakcji własnych. Jeśli pełni funkcję filtru klikalnego, musi być `<button>` i przejmuje stany przycisku.

**Wymiary i żetony.** Wcięcie `1px --dn-od-2 2px` · promień `--dn-r-pill` (wariant `--rola`: `--dn-r-xs` 3 px) · stopień `--dn-fs-xs` 12 px, waga 500 · ikona 12×12 px · odstęp `--dn-od-1` 4 px · `white-space: nowrap`.

**Zachowanie.** Nie zawija się i nie skraca treści — jeśli etykieta jest długa, należy skrócić treść, nie plakietkę.

**Dostępność.** **Barwa nigdy nie jest jedynym nośnikiem znaczenia** — plakietka zawsze niesie tekst, a dla statusów krytycznych dodatkowo ikonę. Przy aktualizacji na żywo: kontener `aria-live="polite"`. Kontrast każdej pary tekst/tło zmierzony w `kontrasty.json`.

```html
<span class="dn-plakietka dn-plakietka--sukces">
  <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor"
       stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
  powodzenie
</span>
<span class="dn-plakietka dn-plakietka--rola">Coordinator</span>
<span class="dn-plakietka dn-plakietka--sygnal">Agent: Walidator wyników</span>
```

---

### 8.7. Kropka sygnału — `.dn-kropka`

| | |
|---|---|
| **Klasa bazowa** | `.dn-kropka` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Wypełniony punkt o średnicy `--dn-wym-kropka`, element sygnaturowy systemu wizualnego |
| **Do czego służy** | Oznacza „tu biegnie praca" — proces w tle, aktywność nadawcy, stan pozycji |
| **Kiedy używać** | Karta sesji z procesem w tle; wpis komunikacji nadawcy piszącego; kolumna statusu z nazwanym nagłówkiem; godło i emblematy środowisk |
| **Kiedy NIE używać** | Status wymagający nazwy (→ `.dn-plakietka`); dekoracja bez znaczenia — kropka zawsze coś znaczy |

**Warianty.**

| Wariant | Klasa | Barwa | Znaczenie |
|---|---|---|---|
| Sygnał (bazowa) | `.dn-kropka` | `--dn-kropka` | Praca w toku, aktywność |
| Tętno | `--tetno` | jak bazowa + animacja | **Jedyny ruch ciągły interfejsu** — pierścień pulsujący w cyklu 2,4 s |
| Sukces | `--sukces` | `--dn-sukces-tekst` | Stan pozytywny |
| Ostrzeżenie | `--ostrzezenie` | `--dn-ostrzezenie-tekst` | Stan wymagający uwagi |
| Błąd | `--blad` | `--dn-blad-tekst` | Stan negatywny |
| Neutralna | `--neutralna` | `--dn-tekst-3` | Bezczynność, brak procesu |

**Stany.**

| Stan | Realizacja |
|---|---|
| Statyczna | Wypełniony punkt bez ruchu |
| Tętno | `box-shadow` rozchodzący się od 0 do 5 px w cyklu `--dn-czas-tetno` 2,4 s, krzywa `--dn-ease` |
| Ograniczony ruch | Przy `prefers-reduced-motion: reduce` tętno zamiera, kropkę wzmacnia **pierścień statyczny** 2 px — informacja zostaje, ruch znika |

**Wymiary i żetony.** Średnica `--dn-wym-kropka` 6 px · promień `--dn-r-pill` · pierścień tętna do 5 px w barwie `--dn-fokus-cien` · czas `--dn-czas-tetno` 2,4 s.

**Zachowanie.** Tętno jest zarezerwowane dla pracy biegnącej w tle — nie wolno go stosować dekoracyjnie. Na jeden widok przypada jeden ruch znaczący (pokrętło `INTENSYWNOSC_RUCHU` = 3/10).

**Dostępność.** Kropka jest `aria-hidden="true"`; znaczenie niesie tekst obok albo etykieta ukryta `.dn-sr-only`. **Nigdy nie jest jedynym nośnikiem stanu.**

```html
<span class="dn-kropka dn-kropka--tetno" aria-hidden="true"></span>
<span class="dn-sr-only">Proces w tle: trwa</span>

<span class="dn-kropka dn-kropka--sukces" aria-hidden="true"></span> powodzenie
```

---

### 8.8. Pusty stan — `.dn-pusty-stan`

| | |
|---|---|
| **Klasa bazowa** | `.dn-pusty-stan` + `-tytul` · `-opis` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Wyśrodkowany blok z ikoną 28 px, tytułem krojem nagłówkowym i opisem o szerokości maks. 40 znaków |
| **Do czego służy** | Komunikuje, że brak treści jest stanem oczekiwanym i naturalnym, a nie awarią |
| **Kiedy używać** | Panel, tabela albo lista przed pierwszym wypełnieniem; wynik filtrowania bez trafień |
| **Kiedy NIE używać** | Błąd pobierania danych (→ komunikat błędu z akcją ponowienia); trwające ładowanie (→ `.dn-spinner`) |

**Anatomia.** Kolumna wyśrodkowana: ikona 28 px → tytuł 16 px krojem nagłówkowym w barwie `--dn-tekst-2` → opis 12 px w barwie `--dn-tekst-3` → opcjonalny pojedynczy przycisk pierwszej akcji.

**Warianty.** Brak modyfikatorów — różnicuje wyłącznie treść: ikona, tytuł, opis właściwe kontekstowi.

**Stany.** Pusty stan sam jest stanem panelu. Osadzony przycisk przejmuje pełny zestaw stanów `.dn-btn`.

**Wymiary i żetony.** Wcięcie `--dn-od-10 --dn-od-5` (40/20 px) · odstęp `--dn-od-2` 8 px · ikona 28×28 px · tytuł `--dn-fs-lg` 17 px, krój `--dn-ff-naglowek`, waga 600 · opis `--dn-fs-sm` 13 px, maks. 40 znaków w wierszu · barwa całości `--dn-tekst-3`.

**Zachowanie.** Bierny widok informacyjny. Jeżeli zawiera akcję — dokładnie jedną, nazywającą pierwszy krok („Dodaj pierwsze źródło", „Zaplanuj pierwsze zadanie").

**Dostępność.** Tytuł jako nagłówek właściwego poziomu. Przy pojawieniu się po filtrowaniu: kontener `aria-live="polite"`, żeby czytnik ogłosił brak trafień. Ikona `aria-hidden="true"`.

```html
<div class="dn-pusty-stan">
  <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor"
       stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
    <rect width="20" height="5" x="2" y="3" rx="1"/>
    <path d="M4 8v11a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8"/><path d="M10 12h4"/>
  </svg>
  <h3 class="dn-pusty-stan-tytul">Kolejka jest pusta</h3>
  <p class="dn-pusty-stan-opis">Queue Manager czeka na pierwsze zadanie. Dodaj je z Workflow Buildera albo z harmonogramu.</p>
  <button class="dn-btn dn-btn--zarys dn-btn--sm" type="button">Zaplanuj pierwsze zadanie</button>
</div>
```

---

### 8.9. Pasek postępu — `.dn-postep`

| | |
|---|---|
| **Klasa bazowa** | `.dn-postep` + `-etykieta` · `-tor` · `-wartosc` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Poziomy wskaźnik ukończenia: tor 4 px o promieniu pełnym z wypełnieniem w barwie sygnału i etykietą krojem mono |
| **Do czego służy** | Pokazuje postęp operacji o **znanym** zakresie — kroków przebiegu, plików do przetworzenia, procent budowania |
| **Kiedy używać** | Execution Monitor, Build Output, Deployment Panel, Results Analyzer |
| **Kiedy NIE używać** | Operacja o nieznanym czasie i zakresie (→ `.dn-spinner`); ciąg nazwanych kroków (→ `.dn-kolejka` / `.dn-krok`) |

**Anatomia.** Rząd `flex` z odstępem 12 px: tor rozciągliwy (`flex: 1`) z wypełnieniem `.dn-postep-wartosc` → etykieta krojem mono 11 px, bez zawijania.

**Warianty.** Brak modyfikatorów.

**Stany.**

| Stan | Realizacja |
|---|---|
| 0% | Tor `--dn-powierzchnia-2`, wypełnienie o szerokości 0 |
| W toku | Wypełnienie `--dn-sygnal-wypelnienie`, szerokość animowana w czasie `--dn-czas-3` 0,22 s |
| 100% | Wypełnienie na całej szerokości; etykieta podaje wynik |
| Nieokreślony | Nie jest stanem tego komponentu — dla nieznanego zakresu stosuje się `.dn-spinner` |

**Wymiary i żetony.** Tor: wys. 4 px, promień `--dn-r-pill`, tło `--dn-powierzchnia-2`, przycinanie zawartości · Wypełnienie: `--dn-sygnal-wypelnienie`, przejście szerokości `--dn-czas-3 --dn-ease` · Etykieta: krój `--dn-ff-mono`, stopień `--dn-fs-xs` 12 px, barwa `--dn-tekst-2` · Odstęp `--dn-od-3` 12 px.

**Zachowanie.** Szerokość wypełnienia zmienia się płynnie; wartość liczbowa w etykiecie zmienia się skokowo. Etykieta zawsze podaje wartość słownie i liczbowo („krok 3 z 7"), nigdy sam procent bez kontekstu.

**Dostępność.** `role="progressbar"` z `aria-valuenow`, `aria-valuemin`, `aria-valuemax` i `aria-label`. Wartość zawsze widoczna także tekstem — pasek nie jest jedynym nośnikiem informacji.

```html
<div class="dn-postep" role="progressbar" aria-valuenow="3" aria-valuemin="0" aria-valuemax="7"
     aria-label="Postęp przebiegu automatyki">
  <span class="dn-postep-tor"><span class="dn-postep-wartosc" style="width:43%"></span></span>
  <span class="dn-postep-etykieta">krok 3 z 7</span>
</div>
```

---

### 8.10. Kolejka i krok — `.dn-kolejka` / `.dn-krok`

| | |
|---|---|
| **Klasy bazowe** | `.dn-kolejka` (kontener) · `.dn-krok` + `-znak` · `-meta` |
| **Kategoria** | D — Dane i treść |
| **Definicja** | Pionowa lista nazwanych kroków o wysokości wiersza, z kreską klasyfikującą po lewej i znakiem stanu w kole 20 px |
| **Do czego służy** | Pokazuje przebieg wieloetapowy: kolejność, stan każdego kroku i trzy wyjścia weryfikacji |
| **Kiedy używać** | Queue Manager, Orchestrator, Workflow Builder, Execution Monitor, Results Analyzer, plan wykonania w oknach ról MultitaskingAI |
| **Kiedy NIE używać** | Jednorodne dane w kolumnach (→ `.dn-tabela`); pojedyncza operacja o znanym zakresie (→ `.dn-postep`) |

**Anatomia.**

```
.dn-kolejka
├── .dn-krok--poprawny    ▌(✓) Pobranie źródeł            00:00:12  ← -znak / treść / -meta
├── .dn-krok--pracuje     ▌(▸) Analiza dokumentów         00:00:37
├── .dn-krok--bledy       ▌(!) Walidacja wyników  obieg 2 00:00:04
└── .dn-krok--wstrzymany  ▌(✕) Zapis do repozytorium      —
   ▲ kreska 2 px w barwie klasy stanu
```

Siatka kroku: `20px 1fr auto` — znak stanu, nazwa kroku, metadane krojem mono.

**Warianty (trzy wyjścia weryfikacji).**

| Wariant | Klasa | Kreska i znak | Znaczenie |
|---|---|---|---|
| Oczekuje | `.dn-krok` | Obrys neutralny, znak z numerem krojem mono | Krok jeszcze nieuruchomiony |
| Pracuje | `--pracuje` | Kreska i znak w barwie sygnału | Krok w toku |
| Poprawny | `--poprawny` | Rodzina `--dn-sukces-*` | Weryfikacja przeszła |
| Błędy | `--bledy` | Rodzina `--dn-ostrzezenie-*` | Weryfikacja wykryła błędy — krok wraca do obiegu, licznik obiegów w metadanych |
| Wstrzymany | `--wstrzymany` | Rodzina `--dn-blad-*` | Zdarzenie nieoczekiwane — kolejka wstrzymana do decyzji Operatora |

**Stany.** Stan kroku **jest** jego wariantem — nie ma osobnej warstwy stanów. Znak stanu zawsze niesie ikonę albo numer; kreska barwna jest wzmocnieniem, nie jedynym nośnikiem.

**Wymiary i żetony.** Krok: min. wys. `--dn-wym-wiersz` 36 px, wcięcie `0 --dn-od-3`, siatka `20px 1fr auto`, odstęp `--dn-od-3` · Kreska lewa 2 px · Znak: 20×20 px, promień `--dn-r-pill`, obrys 1 px, ikona 12×12 px, krój mono `--dn-fs-xs` · Metadane: krój mono, `--dn-fs-xs` 12 px, barwa `--dn-tekst-3`.

**Zachowanie.** Kroki aktualizują się na żywo. Wstrzymanie kolejki (wariant `--wstrzymany`) zatrzymuje kroki następne, ale **nie odbiera Operatorowi żadnej kontrolki** — akcje „Ponów", „Pomiń", „Zakończ" pozostają klikalne.

**Dostępność.** Kolejka jako `<ol>` — kolejność ma znaczenie. Stan kroku w tekście (widocznym albo `.dn-sr-only`), nie tylko w barwie kreski. Aktualizacja na żywo: `aria-live="polite"` na kontenerze kolejki.

```html
<ol class="dn-kolejka">
  <li class="dn-krok dn-krok--poprawny">
    <span class="dn-krok-znak" aria-hidden="true">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"
           stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
    </span>
    <span>Pobranie źródeł <span class="dn-sr-only">— krok poprawny</span></span>
    <span class="dn-krok-meta">00:00:12</span>
  </li>
  <li class="dn-krok dn-krok--bledy">
    <span class="dn-krok-znak" aria-hidden="true">!</span>
    <span>Walidacja wyników <span class="dn-sr-only">— zakończony z błędami</span></span>
    <span class="dn-krok-meta">obieg 2 · 00:00:04</span>
  </li>
</ol>
```

---

## 9. Kategoria E — Informacja zwrotna i nakładki

### 9.1. Modal — `.dn-modal`

| | |
|---|---|
| **Klasa bazowa** | `.dn-modal` + `-naglowek` · `-tytul` · `-cialo` · `-stopka` |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Definicja** | Natywne okno `<dialog>` o szerokości maks. `--dn-wym-modal`, z nakładką przyciemniającą i rozmyciem 2 px |
| **Do czego służy** | Przerywa bieżący widok na czas decyzji, konfiguracji albo kreatora |
| **Kiedy używać** | Modal „Pula kont Code CLI"; kreatory komponentów własnych strefy 2; potwierdzenia czynności nieodwracalnych; konfiguracja uruchamiana z okna operacyjnego |
| **Kiedy NIE używać** | Potwierdzenie wyniku, który już nastąpił (→ `.dn-toast`); objaśnienie elementu (→ `.dn-tooltip`); treść, którą Operator ma porównywać z tłem (→ panel boczny) |

**Anatomia.**

```
::backdrop  ── --dn-nakladka + rozmycie 2 px
┌ .dn-modal · maks. 560 px · maks. min(80dvh, 720px) ─────────┐
│ .dn-modal-naglowek   [tytuł 20 px]              [✕]         │ wcięcie 16/20
├──────────────────────────────────────────────────────────────┤
│ .dn-modal-cialo      treść przewijana niezależnie            │ wcięcie 20
├──────────────────────────────────────────────────────────────┤
│ .dn-modal-stopka                     [Anuluj] [Zatwierdź]    │ tło --dn-powierzchnia-2
└──────────────────────────────────────────────────────────────┘
```

**Warianty.** Brak modyfikatorów — jedna forma okna różnicowana treścią (formularz, potwierdzenie, kreator wieloetapowy złożony z zakładek albo kroków).

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Zamknięty | brak `[open]` | Brak śladu w układzie |
| Otwierany | `[open]` + `dn-wejscie` | Krycie 0 → 1, `translateY(8px) scale(0.98)` → stan docelowy, czas `--dn-czas-3` 0,22 s |
| Otwarty | `[open]` | Przechwytuje fokus; nakładka przyciemnia i rozmywa tło |
| Ładowanie | — | Wskaźnik w miejscu ciała na czas przygotowania danych |
| Błąd | — | Komunikat osadzony w ciele, nad stopką — nie zastępuje stopki |
| Próba zatwierdzenia z brakami | — | Przycisk stopki **pozostaje klikalny**; braki ujawnia komunikat przy polach, nigdy blokada przycisku |

**Wymiary i żetony.** Szerokość `min(--dn-wym-modal, calc(100vw − --dn-od-8))` = maks. 560 px · wysokość `min(80dvh, 720px)` · promień `--dn-r-lg` 10 px · cień `--dn-cien-lg` (największy w systemie) · nakładka `--dn-nakladka` + `backdrop-filter: blur(2px)` — **jedyne miejsce rozmycia w całym interfejsie** · nagłówek i stopka: wcięcie `--dn-od-4 --dn-od-5` · ciało: wcięcie `--dn-od-5` 20 px · tytuł `--dn-fs-xl` 21 px krojem nagłówkowym · warstwa `--dn-z-modal` = 900.

**Zachowanie.** Otwarcie przez `showModal()` (albo atrybut `command="show-modal"` z zapasem w `wspolne.js`). Zamknięcie: przycisk zamknięcia w nagłówku, klawisz `Escape`, kliknięcie w nakładkę, akcja w stopce. Nakładka wstrzymuje interakcję z tłem wyłącznie na czas otwarcia — droga wyjścia jest zawsze dostępna.

**Dostępność.** Natywny `<dialog>` daje rolę `dialog`, pułapkę fokusu i obsługę `Escape` bez skryptu. Tytuł wskazany przez `aria-labelledby`. Fokus po otwarciu trafia na pierwszą kontrolkę albo na tytuł; po zamknięciu wraca na element wywołujący. Stopka: akcja główna po prawej, wycofanie na lewo od niej.

```html
<button class="dn-btn dn-btn--zarys" type="button" command="show-modal" commandfor="m-pula">
  Pula kont Code CLI
</button>

<dialog class="dn-modal" id="m-pula" aria-labelledby="m-pula-tytul">
  <header class="dn-modal-naglowek">
    <h2 class="dn-modal-tytul" id="m-pula-tytul">Pula kont Code CLI</h2>
    <button class="dn-btn-ikona" type="button" style="margin-left:auto"
            command="close" commandfor="m-pula" aria-label="Zamknij okno">…</button>
  </header>
  <div class="dn-modal-cialo">
    <p>Konta używane rotacyjnie przez kanał CLI. Kolejność wyznacza priorytet przydziału.</p>
  </div>
  <footer class="dn-modal-stopka">
    <button class="dn-btn dn-btn--zarys" type="button" command="close" commandfor="m-pula">Anuluj</button>
    <button class="dn-btn dn-btn--atrament" type="button">Zapisz pulę</button>
  </footer>
</dialog>
```

---

### 9.2. Powiadomienie — `.dn-toasty` / `.dn-toast`

| | |
|---|---|
| **Klasy bazowe** | `.dn-toasty` (stos) · `.dn-toast` + `-tytul` · `-tresc` |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Definicja** | Pływający pasek 280–420 px w prawym dolnym rogu, z ikoną, tytułem, treścią i kreską semantyczną po lewej |
| **Do czego służy** | Potwierdza wynik czynności, która **już się dokonała**, bez przerywania pracy |
| **Kiedy używać** | Zapis profilu izolacji, zakończenie przebiegu automatyki, wynik operacji w tle, komunikaty Mobile |
| **Kiedy NIE używać** | Decyzja wymagana od Operatora (→ `.dn-modal`); stan trwały, który ma pozostać widoczny (→ komunikat w układzie panelu); objaśnienie kontrolki (→ `.dn-tooltip`) |

**Anatomia.** Stos `.dn-toasty` przyklejony do prawego dolnego rogu (`inset: auto 16px 16px auto`), warstwa 1000, odstęp 8 px. Pojedyncze powiadomienie: siatka `auto 1fr auto` — ikona 16 px → kolumna z tytułem i treścią → przycisk zamknięcia.

**Warianty.**

| Wariant | Klasa | Kreska lewa 2 px | Ikona |
|---|---|---|---|
| Neutralny | `.dn-toast` | brak | wg treści |
| Sukces | `--sukces` | `--dn-sukces-tekst` | `ptaszek` |
| Ostrzeżenie | `--ostrzezenie` | `--dn-ostrzezenie-tekst` | `ostrzezenie` |
| Błąd | `--blad` | `--dn-blad-tekst` | `blad` |
| Informacja | `--informacja` | `--dn-informacja-tekst` | `info` |

**Stany.** Pojawienie (animacja `dn-wejscie`, 0,22 s — ta sama co modal) · widoczny · znikanie (po czasie ekspozycji albo po kliknięciu zamknięcia) · najechanie na przycisk zamknięcia. Powiadomienie **nie ma** stanu ładowania ani błędu własnego — wariant `--blad` opisuje cudzą operację.

**Wymiary i żetony.** Szerokość `--dn-wym-toast-min` 280 px do `--dn-wym-toast-max` 420 px · wcięcie `--dn-od-3 --dn-od-4` · promień `--dn-r-md` 8 px · tło `--dn-panel`, obrys `--dn-obrys`, cień `--dn-cien-3` · ikona `--dn-wym-ikona` 16 px · tytuł `--dn-fs-base` 14 px waga 600 · treść `--dn-fs-sm` 13 px, barwa `--dn-tekst-2` · warstwa `--dn-z-powiadomienie` 1000.

**Zachowanie.** Pojawia się samoczynnie w reakcji na zdarzenie systemowe. Znika po czasie ekspozycji albo natychmiast po zamknięciu. Kilka powiadomień układa się pionowo, najnowsze na dole stosu. Nie wstrzymuje interakcji z resztą widoku.

**Dostępność.** Stos jako `aria-live="polite"` (`assertive` wyłącznie dla wariantu `--blad`) i `role="status"`. Ikona `aria-hidden="true"` — znaczenie niesie tytuł. Przycisk zamknięcia z `aria-label`. Czas ekspozycji nie może być jedyną drogą zniknięcia — zamknięcie ręczne zawsze dostępne.

```html
<div class="dn-toasty" role="status" aria-live="polite">
  <div class="dn-toast dn-toast--sukces">
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M20 6 9 17l-5-5"/></svg>
    <div>
      <p class="dn-toast-tytul">Profil izolacji zapisany</p>
      <p class="dn-toast-tresc">Zakres: Środowisko CodeStudio. Zmiany działają od nowej sesji.</p>
    </div>
    <button class="dn-btn-ikona" type="button" aria-label="Zamknij powiadomienie">…</button>
  </div>
</div>
```

---

### 9.3. Dymek objaśnienia — `.dn-tooltip`

| | |
|---|---|
| **Klasy bazowe** | `.dn-tooltip` (wyzwalacz) · `.dn-tooltip-tresc` |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Definicja** | Mały dymek na powierzchni atramentowej, ujawniany nad wyzwalaczem przy najechaniu albo fokusie |
| **Do czego służy** | Wyjaśnia działanie pojedynczego elementu bez zajmowania miejsca w układzie — nośnik oznaczenia `[?]` |
| **Kiedy używać** | Każdy element konfiguracji; przyciski ikonowe paska górnego; skróty i identyfikatory wymagające rozwinięcia; wyjaśnienie, dlaczego czynność nie przyniesie skutku w bieżącym stanie |
| **Kiedy NIE używać** | Treść krytyczna dla decyzji (→ `.dn-pole-opis`); treść dłuższa niż dwa zdania (→ modal albo panel pomocy); jedyne źródło nazwy przycisku ikonowego (→ `aria-label`) |

**Anatomia.** Wyzwalacz `position: relative` z zagnieżdżoną treścią pozycjonowaną absolutnie: `bottom: calc(100% + 8px)`, wyśrodkowana poziomo, maks. 260 px.

**Warianty.** Brak modyfikatorów. Pozycja dostosowywana do dostępnej przestrzeni w warstwie wykonawczej, bez zmiany klasy.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Ukryty | — | Krycie 0, `visibility: hidden`, `pointer-events: none`, przesunięcie 4 px w dół |
| Widoczny | `:hover`, `:focus-within` | Krycie 1, widoczność, przesunięcie do 0 — czas `--dn-czas-2` 0,16 s |

**Wymiary i żetony.** Maks. szerokość 260 px · wcięcie `--dn-od-2 --dn-od-3` (8/12 px) · promień `--dn-r-sm` 6 px · tło `--dn-rama`, tekst `--dn-rama-tekst` — powierzchnia atramentowa w obu motywach · cień `--dn-cien-3` · stopień `--dn-fs-sm` 13 px · odsunięcie `--dn-od-2` 8 px · warstwa `--dn-z-tooltip` 1100.

**Zachowanie.** Pojawia się przy najechaniu myszą i przy fokusie klawiatury (`:focus-within` — warunek dostępności, nie ozdoba). Znika po utracie obu. Nie wymaga kliknięcia ani zamykania.

**Dostępność.** Treść z `role="tooltip"` i `id` wskazanym przez `aria-describedby` wyzwalacza. Dymek **nie zastępuje** `aria-label` przycisku ikonowego — uzupełnia go. Klawisz `Escape` powinien ukryć dymek wywołany fokusem (obsługa skryptem tam, gdzie dymek zasłania treść).

```html
<span class="dn-tooltip">
  <button class="dn-btn-ikona" type="button" aria-label="Objaśnienie ustawienia"
          aria-describedby="d-izo">
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"
         stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/>
      <path d="M12 16v-4"/><path d="M12 8h.01"/></svg>
  </button>
  <span class="dn-tooltip-tresc" role="tooltip" id="d-izo">
    Izolacja techniczna ogranicza dostęp procesu do zasobów urządzenia. Stan wyjściowy: wyłączona.
  </span>
</span>
```

---

### 9.4. Wskaźnik pracy — `.dn-spinner`

| | |
|---|---|
| **Klasa bazowa** | `.dn-spinner` |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Definicja** | Pierścień 14 px o obrysie 2 px z jednym bokiem w barwie sygnału, obracający się w tempie 0,8 s |
| **Do czego służy** | Sygnalizuje operację o **nieznanym** czasie trwania |
| **Kiedy używać** | Wnętrze przycisku w trakcie akcji; Chat Window w trakcie odbierania odpowiedzi; panel albo tabela w trakcie pobierania danych; Build Output w trakcie budowania |
| **Kiedy NIE używać** | Operacja o znanym zakresie (→ `.dn-postep`); praca biegnąca w tle, o której Operator ma tylko wiedzieć (→ `.dn-kropka--tetno`); pełny ekran ładowania — wskaźnik jest komponentem osadzanym |

**Anatomia.** Pojedynczy element bez zawartości: koło o obrysie `--dn-obrys-mocny` z górnym bokiem w barwie `--dn-sygnal-wypelnienie`.

**Warianty.** Brak modyfikatorów. Wewnątrz przycisku wskaźnik jest realizowany pseudoelementem `.dn-btn[aria-busy='true']::after` w barwie `currentColor` — dziedziczy barwę wariantu przycisku.

**Stany.** Obecny (obrót ciągły) · nieobecny · przy `prefers-reduced-motion: reduce` obrót skraca się do 0,01 ms — wskaźnik pozostaje widoczny statycznie zamiast wirować bez końca.

**Wymiary i żetony.** Średnica `--dn-wym-spinner` 14 px · obrys 2 px · promień `--dn-r-pill` · animacja `dn-obrot` 0,8 s liniowo, bez końca. Żetony: `--dn-obrys-mocny` · `--dn-sygnal-wypelnienie` · `--dn-wym-spinner` · `--dn-r-pill`.

**Zachowanie.** Bierny — nie reaguje na najechanie ani kliknięcie. Pojawia się na czas operacji i ustępuje miejsca wynikowi albo komunikatowi o niepowodzeniu.

**Dostępność.** Sam wskaźnik `aria-hidden="true"`; stan komunikuje `aria-busy="true"` na kontenerze operacji plus tekst („Pobieram dane…") w obszarze `aria-live="polite"`. Wskaźnik nigdy nie jest jedynym komunikatem dla czytnika ekranu.

```html
<button class="dn-btn dn-btn--atrament" type="button" aria-busy="true">Zapisuję profil</button>

<div aria-busy="true" aria-live="polite">
  <span class="dn-spinner" aria-hidden="true"></span>
  Pobieram rejestr konektorów…
</div>
```

---

## 10. Kategoria F — Tożsamość

### 10.1. Awatar — `.dn-awatar`

| | |
|---|---|
| **Klasa bazowa** | `.dn-awatar` + `.dn-awatar-stan` |
| **Kategoria** | F — Tożsamość |
| **Definicja** | Koło (albo kwadrat) o boku `--dn-wym-awatar` z inicjałem albo obrazem, na gradiencie ilustracyjnym |
| **Do czego służy** | Odpowiada na pytanie „kto" — Operator, model, agent, rola |
| **Kiedy używać** | Chat Window, Model Panels w Roundtable, karty ról panelu orkiestracji, pasek górny (profil Operatora), Agent Manager |
| **Kiedy NIE używać** | Ikona funkcji (→ ikona z zestawu w `.dn-btn-ikona`); emblemat środowiska (→ `.dn-karta-srodowiska-godlo`); wskaźnik pracy (→ `.dn-kropka`) |

**Anatomia.** Koło z wyśrodkowanym inicjałem krojem nagłówkowym w wadze 600 albo obrazem przyciętym do kształtu; opcjonalny wskaźnik stanu `.dn-awatar-stan` 8×8 px w prawym dolnym rogu, z pierścieniem 2 px w barwie powierzchni.

**Warianty.**

| Wariant | Klasa | Wymiar / wygląd | Zastosowanie |
|---|---|---|---|
| Bazowy | `.dn-awatar` | 28 px, gradient atramentowy | Operator, tożsamość ogólna |
| Mały | `--sm` | 24 px, stopień 11 px | Wpis komunikacji, gęsty wiersz listy |
| Duży | `--lg` | 36 px, stopień 14 px | Nagłówek profilu, karta roli |
| Kwadratowy | `--kwadrat` | promień 6 px | **Identyfikacja nie-osobowa** — agent, w odróżnieniu od modelu bazowego |
| Inteligencja | `--inteligencja` | Gradient sygnałowy, tekst biały | Model, agent, rola MultitaskingAI |

Gradienty (`--dn-grad-atrament`, `--dn-grad-sygnal`) występują tu w jednym z trzech dopuszczonych zakresów ilustracyjnych — awatar bez zdjęcia. Gradient nigdy nie jest tłem przycisku, karty ani sekcji.

**Stany.**

| Stan | Realizacja |
|---|---|
| Domyślny | Inicjał albo obraz wg wariantu |
| Ze wskaźnikiem stanu | `.dn-awatar-stan` — kropka 8 px w barwie `--dn-sukces-tekst` z pierścieniem `--dn-powierzchnia` |
| Ładowanie obrazu | Inicjał widoczny do czasu wczytania obrazu docelowego |
| Błąd obrazu | Powrót do inicjału — awatar nigdy nie zostaje pusty |
| Jako przycisk | Gdy awatar otwiera menu, przejmuje komplet stanów `.dn-btn-ikona` |

**Wymiary i żetony.** `--dn-wym-awatar` 28 px · `--dn-wym-awatar-sm` 24 px · `--dn-wym-awatar-lg` 36 px · promień `--dn-r-pill` (kwadrat: `--dn-r-sm` 6 px) · stopień `--dn-fs-sm` / `-xs` / `-md` · krój `--dn-ff-naglowek`, waga 600 · wskaźnik stanu 8×8 px, pierścień 2 px.

**Zachowanie.** W Chat Window i panelach wielomodelowych pełni wyłącznie funkcję identyfikacyjną. W pasku górnym otwiera menu Operatora — wtedy musi być przyciskiem.

**Dostępność.** Awatar identyfikacyjny: `aria-hidden="true"`, bo nazwa nadawcy stoi obok. Awatar jako przycisk: `<button>` z `aria-label` („Menu Operatora"). Obraz: `alt` pusty przy powtórzeniu nazwy obok, opisowy gdy stoi sam. Wskaźnik stanu wymaga etykiety tekstowej — barwa nie wystarcza.

```html
<span class="dn-awatar dn-awatar--sm" aria-hidden="true">OP</span>

<span class="dn-awatar dn-awatar--lg dn-awatar--kwadrat dn-awatar--inteligencja" aria-hidden="true">
  WW
  <span class="dn-awatar-stan"></span>
</span>
<span class="dn-sr-only">Agent Walidator wyników — aktywny</span>
```

---

### 10.2. Always On Display — `.dn-aod`

| | |
|---|---|
| **Klasa bazowa** | `.dn-aod` + `-rdzen` · `-tresc` |
| **Kategoria** | F — Tożsamość |
| **Definicja** | Pływająca pigułka zakotwiczona w prawym dolnym rogu okna, z rdzeniem 36 px na gradiencie sygnałowym i pierścieniem tętna |
| **Do czego służy** | Punkt dostępu do funkcji globalnej Always On Display — agenta obecnego ponad środowiskami, modułami i sesjami |
| **Kiedy używać** | Ponad całą powłoką aplikacji, niezależnie od otwartego środowiska; sekcja Monitor procesu panelu orkiestracji MultitaskingAI |
| **Kiedy NIE używać** | Identyfikacja nadawcy w rozmowie (→ `.dn-awatar--inteligencja`); powiadomienie o wyniku (→ `.dn-toast`) |

**Anatomia.** Pigułka `flex`: rdzeń 36 px (gradient sygnałowy, pierścień `::after` z tętnem) → treść 12 px o szerokości maks. 36 znaków, z wyróżnieniem `<strong>` w barwie tekstu podstawowego.

**Warianty.** Brak modyfikatorów. Odmiana funkcjonalna: poza MultitaskingAI działa doradczo, wewnątrz MultitaskingAI dodatkowo jako obserwator procesu — różnicę niesie treść i plakietka roli, nie klasa.

**Stany.**

| Stan | Realizacja |
|---|---|
| Spoczynek | Pigułka widoczna, rdzeń z pierścieniem tętna 2,4 s |
| Rozwinięty | Powierzchnia interakcji otwarta obok rdzenia (forma zależna od kontekstu wywołania) |
| Nadzór (MultitaskingAI) | Plakietka roli przy nazwie w sekcji Monitor procesu |
| Ograniczony ruch | Tętno zamiera; pierścień pozostaje statyczny |
| Wyłączony | **Nie istnieje jako stan interfejsu** — ukrycie Always On Display jest ustawieniem, nie blokadą |

**Wymiary i żetony.** Pozycja `right/bottom: --dn-od-5` 20 px · warstwa `--dn-z-aod` 1200 · promień `--dn-r-pill` · wcięcie `--dn-od-2 --dn-od-4 --dn-od-2 --dn-od-2` · rdzeń `--dn-wym-awatar-lg` 36 px na `--dn-grad-sygnal` · pierścień `inset: -3px`, obrys 1 px `--dn-sygnal-obrys`, animacja `dn-tetno` 2,4 s · treść `--dn-fs-sm` 13 px, maks. 36 znaków · tło `--dn-panel`, cień `--dn-cien-3`.

**Zachowanie.** Widoczny stale, niezależnie od karty sesji i modułu — nie podlega przeładowaniu obszaru roboczego. Aktywacja otwiera powierzchnię natychmiastowej interakcji; zamknięcie przywraca spoczynek bez utraty kontekstu.

**Dostępność.** Kontener `role="complementary"` z `aria-label="Always On Display"`. Rdzeń jako `<button>` z etykietą. Warstwa 1200 leży nad wszystkim poza Centrum poleceń (1300) — nie może zasłaniać stopki modala, dlatego przy otwartym modalu jest odsuwana albo wyciszana. Treść aktualizowana na żywo: `aria-live="polite"`.

```html
<aside class="dn-aod" role="complementary" aria-label="Always On Display">
  <button class="dn-aod-rdzen" type="button" aria-label="Otwórz Always On Display">
    <svg aria-hidden="true" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor"
         stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
      <path d="M12 8V4H8"/><rect width="16" height="12" x="4" y="8" rx="2"/>
      <path d="M2 14h2"/><path d="M20 14h2"/><path d="M15 13v2"/><path d="M9 13v2"/>
    </svg>
  </button>
  <p class="dn-aod-tresc" aria-live="polite">
    <strong>Monitor procesu</strong> — Coordinator prowadzi 3 kolejki, 1 krok czeka na decyzję.
  </p>
</aside>
```

---

## 11. Kategoria G — Komunikacja

### 11.1. Wpis okna komunikacji — `.dn-wpis`

| | |
|---|---|
| **Klasa bazowa** | `.dn-wpis` + `-medalion` · `-tozsamosc` · `-nadawca` · `-godzina` · `-tresc` |
| **Kategoria** | G — Komunikacja |
| **Definicja** | Blok rozmowy o szerokości maks. 860 px, z medalionem nadawcy 24 px i kreską klasyfikującą 2 px po lewej |
| **Do czego służy** | Prezentuje jedną wypowiedź w Chat Window — oknie wspólnym wszystkim piętnastu modułom |
| **Kiedy używać** | Chat Window każdego modułu; Executor Chat, Coordinator Chat, Voice Console; dziennik zdarzeń automatyki |
| **Kiedy NIE używać** | Dane tabelaryczne (→ `.dn-tabela`); komunikat systemu o wyniku operacji (→ `.dn-toast`); log techniczny (→ `.dn-kod`) |

**Anatomia.**

```
┌ .dn-wpis --inteligencja · maks. 860 px ────────────────────────┐
▌│ [medalion 24]  COORDINATOR  [plakietka roli]        14:07     │  kreska 2 px = klasa
 │                                                                │
 │ Treść wypowiedzi, stopień 14 px, interlinia luźna 1,6.         │
 └────────────────────────────────────────────────────────────────┘
```

Siatka: `24px 1fr` — medalion w pierwszej kolumnie, tożsamość i treść w drugiej. Godzina odsunięta do prawej krawędzi krojem mono.

**Warianty — trzy klasy semantyczne, nie dziewięć barw.**

| Wariant | Klasa | Kreska i medalion | Kto |
|---|---|---|---|
| Człowiek | `--czlowiek` | Atrament — medalion wypełniony | Operator |
| Inteligencja | `--inteligencja` | Sygnał — medalion na tle sygnałowym | model · agent · Coordinator · Executor 1 · Executor 2 · Executor 3 / Validator |
| System | `--system` | Neutralna, tło wycofane (przezroczyste) | Automation · Always On Display · wynik narzędzia |
| Pracuje | `--pracuje` | Kropka tętna przy nazwie nadawcy | Nadawca aktywnie pracujący — łączy się z klasą semantyczną |

Rozróżnienie roli **wewnątrz** klasy niesie ikona w medalionie, nazwa nadawcy i plakietka roli — nigdy dodatkowa barwa tła.

**Stany.** Spoczynek · pracuje (kropka tętna 2,4 s przy nadawcy, dodawana pseudoelementem) · ograniczony ruch (tętno zamiera). Wpis nie ma stanów najechania ani wyboru — to treść, nie kontrolka. Akcje wpisu (kopiuj, odpowiedz) należą do przybornika pojawiającego się obok.

**Wymiary i żetony.** Maks. szerokość 860 px · siatka `--dn-wym-awatar-sm 1fr`, odstęp `--dn-od-2 --dn-od-3` · wcięcie `--dn-od-3 --dn-od-4` · obrys 1 px `--dn-obrys-subtelny`, lewa krawędź 2 px · promień `--dn-r-md` 8 px · medalion 24×24 px, promień `--dn-r-sm`, ikona `--dn-wym-ikona-sm` 14 px · nadawca `--dn-fs-xs` 12 px, waga 600, wersaliki, odstęp liter `--dn-ls-wersaliki` 0,08 em · godzina krojem mono `--dn-fs-xs` · treść `--dn-fs-md` 15 px, interlinia `--dn-lh-luzny` 1,6, `white-space: pre-wrap`.

**Zachowanie.** Wpisy układają się chronologicznie w pasie komunikacji o wysokości `--dn-wym-pas-komunikacji` 320 px. Treść odbierana strumieniowo narasta w miejscu; wpis pozostaje oznaczony jako pracujący do zakończenia strumienia.

**Dostępność.** Lista wpisów jako `<ol>` z `aria-live="polite"` — czytnik ogłasza nową wypowiedź bez przerywania. Nazwa nadawcy w treści, nie tylko w barwie kreski. Godzina jako `<time datetime="…">`. Medalion `aria-hidden="true"`.

```html
<li class="dn-wpis dn-wpis--inteligencja dn-wpis--pracuje">
  <span class="dn-wpis-medalion" aria-hidden="true">
    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"
         stroke-linecap="round" stroke-linejoin="round">
      <path d="M9 9.003a1 1 0 0 1 1.517-.859l4.997 2.997a1 1 0 0 1 0 1.718l-4.997 2.997A1 1 0 0 1 9 14.996z"/>
      <circle cx="12" cy="12" r="10"/>
    </svg>
  </span>
  <span class="dn-wpis-tozsamosc">
    <span class="dn-wpis-nadawca">Executor 1</span>
    <span class="dn-plakietka dn-plakietka--rola">Wykonawca</span>
    <time class="dn-wpis-godzina" datetime="2026-08-14T14:07">14:07</time>
  </span>
  <p class="dn-wpis-tresc">Pobrałem 12 dokumentów z Library Explorer. Rozpoczynam analizę zgodnie z planem kroku 2.</p>
</li>
```

---

### 11.2. Pole promptu — `.dn-prompt`

| | |
|---|---|
| **Klasa bazowa** | `.dn-prompt` + `-grot` · `-obszar` |
| **Kategoria** | G — Komunikacja |
| **Definicja** | Kontener wejścia z grotem `❯` w barwie sygnału, rozciągliwym obszarem tekstu i akcją wysyłki po prawej |
| **Do czego służy** | Przyjmuje polecenie kierowane do modelu, agenta albo roli |
| **Kiedy używać** | Chat Window każdego modułu; Executor Chat i Coordinator Chat; Voice Console (obok wejścia głosowego); Prompt Builder |
| **Kiedy NIE używać** | Wartość zapisywana w konfiguracji (→ `.dn-pole`); wyszukiwanie (→ `.dn-szukaj`); wpisanie polecenia powłoki (→ okno Terminala) |

**Anatomia.**

```
┌ .dn-prompt · siatka auto 1fr auto · wcięcie 8 px ──────────────┐
│ ❯ │ Wpisz polecenie…                              │ [wyślij]   │
│ ▲ grot mono, barwa --dn-kropka                                 │
│   ▲ .dn-prompt-obszar: 40–160 px, bez obrysu i tła własnego    │
└────────────────────────────────────────────────────────────────┘
  .dn-przybornik pod polem: załączniki · biblioteka · agent · model
```

**Warianty.** Brak modyfikatorów.

**Stany.**

| Stan | Nośnik | Realizacja |
|---|---|---|
| Spoczynek | — | Obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia`, podpowiedź `--dn-tekst-3` |
| Fokus | `:focus-within` | Obrys `--dn-fokus` + poświata `--dn-cien-sygnal` — **cały kontener**, nie samo pole tekstowe |
| Z treścią | — | Obszar rośnie od 40 px do 160 px, potem przewija się wewnętrznie |
| Wysyłanie | — | Przycisk wysyłki z `aria-busy`; pole pozostaje edytowalne — Operator może pisać kolejne polecenie |

**Wymiary i żetony.** Wcięcie `--dn-od-2` 8 px · odstęp `--dn-od-2` · promień `--dn-r-md` 8 px · obszar: min. 40 px, maks. 160 px, wcięcie `--dn-od-2 0`, stopień `--dn-fs-md` 15 px, interlinia bazowa, bez zmiany rozmiaru ręcznego · grot: krój mono, waga 600, barwa `--dn-kropka`, wcięcie górne `--dn-od-2`, niezaznaczalny.

**Zachowanie.** Obszar rośnie wraz z treścią do 160 px. Grot jest sygnaturą wejścia — powtarza motyw godła (podwójny grot i kropka), nie jest kontrolką. Wysłanie polecenia czyści obszar i dodaje wpis do listy rozmowy.

**Dostępność.** Obszar jako `<textarea>` z etykietą ukrytą. Klawiatura: `Enter` wysyła, `Shift+Enter` łamie wiersz (zachowanie do zadeklarowania w podpowiedzi pola). Grot `aria-hidden="true"`. Przybornik jako `role="toolbar"` poniżej.

```html
<form class="dn-prompt">
  <span class="dn-prompt-grot" aria-hidden="true">❯</span>
  <label class="dn-sr-only" for="prompt-chat">Polecenie do modelu</label>
  <textarea class="dn-prompt-obszar" id="prompt-chat" rows="1"
            placeholder="Wpisz polecenie — Enter wysyła, Shift+Enter łamie wiersz"></textarea>
  <button class="dn-btn-ikona" type="submit" aria-label="Wyślij polecenie">
    <svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75"
         stroke-linecap="round" stroke-linejoin="round">
      <path d="M3.714 3.048a.498.498 0 0 0-.683.627l2.843 7.627a2 2 0 0 1 0 1.396l-2.842 7.627a.498.498 0 0 0 .682.627l18-8.5a.5.5 0 0 0 0-.904z"/>
      <path d="M6 12h16"/>
    </svg>
  </button>
</form>
```

---

## 12. Macierz komponent × moduł

Legenda: `●●●` największe w platformie nasilenie komponentu · `●●` znaczące nasilenie · `●` obecny · `—` nie występuje w sposób charakterystyczny.

### 12.1. Komponenty × piętnaście modułów

| Komponent | Studio | Research | Library | Translate | Browser | Assistant | Roundtable | Workspace | Automations | Design | Apps | Terminal | Developer | Diagnostics | Agents |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `.dn-btn` | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ● | ● | ● | ●● |
| `.dn-btn-ikona` | ●● | ● | ●● | ● | ●● | ● | ● | ● | ●● | ● | ● | ●● | ●● | ● | ● |
| `.dn-pole` | ● | ● | ● | ●● | ● | ● | ● | ●● | ●● | ●● | ●● | — | ● | ● | ●●● |
| `.dn-wybor` | — | ● | ● | — | — | — | ● | ● | ●● | — | ● | — | — | ● | ●● |
| `.dn-check` | — | ●● | ● | — | — | — | ● | — | ● | — | — | — | — | ● | ● |
| `.dn-radio` | — | — | — | ● | — | — | ● | — | ● | — | — | — | — | — | ● |
| `.dn-przelacznik` | — | — | ● | ● | — | ● | — | ● | ●● | — | ● | — | — | ● | ●●● |
| `.dn-suwak` | — | — | — | — | — | — | ●● | — | — | ● | — | — | — | — | ●● |
| `.dn-szukaj` | ● | ●● | ●●● | ●● | ● | — | — | ● | ● | ●● | ● | — | ●● | ●● | ●● |
| `.dn-pasek` | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● |
| `.dn-boczna` | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● |
| `.dn-karty-sesji` | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | — | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● |
| `.dn-zakladki` | ● | ● | ● | ●● | ● | — | ●● | ● | ● | ● | ●● | ●●● | ●● | ● | ●● |
| `.dn-listwa` | — | — | — | — | — | — | — | — | — | — | — | — | — | — | — |
| `.dn-przybornik` | ●● | ● | ● | ● | ● | ●● | ● | ● | ● | ●● | ● | ● | ● | ● | ● |
| `.dn-karta` | ●● | ●● | ●● | ● | ●● | ● | ●● | ●● | ●● | ●● | ●● | ● | ●● | ●● | ●● |
| `.dn-karta-srodowiska` | — | — | — | — | — | — | — | — | — | — | — | — | — | — | — |
| `.dn-kafel` | — | — | — | — | — | ● | — | ● | ● | — | — | — | — | — | ● |
| `.dn-tabela` | ● | ●● | ●● | ● | — | ●● | — | ● | ●●● | — | ● | ● | ● | ●●● | ●● |
| `.dn-dane` | ● | ● | ●● | — | ● | ● | ● | ● | ●●● | — | ● | ●● | ●● | ●●● | ● |
| `.dn-plakietka` | ● | ●● | ●●● | ● | ● | ● | ●● | ● | ●● | ● | ● | ● | ● | ●● | ●● |
| `.dn-kropka` | ● | ● | ● | ● | ● | ●● | ●● | ● | ●●● | ● | ● | ●● | ● | ●● | ● |
| `.dn-pusty-stan` | ● | ●● | ●● | ● | ● | ● | ● | ● | ●● | ●● | ● | ● | ● | ●● | ● |
| `.dn-postep` | ● | ● | ● | ● | — | ● | — | ● | ●●● | ●● | ●● | ● | ●● | ● | — |
| `.dn-kolejka` / `.dn-krok` | ● | ●● | — | — | — | ●● | ●● | ● | ●●● | ● | ●● | ● | ●● | ●● | ● |
| `.dn-modal` | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ● | ● | ● | ●●● |
| `.dn-toast` | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ● | ● | ● | ● |
| `.dn-tooltip` | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ● | ● | ● | ●● |
| `.dn-spinner` | ● | ●● | ● | ●● | ● | ●● | ●● | ● | ●● | ●● | ●● | ● | ●● | ● | ● |
| `.dn-awatar` | ● | ● | ● | ● | ● | ●● | ●●● | ● | ● | ● | ● | ● | ● | ● | ●● |
| `.dn-aod` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| `.dn-wpis` | ●●● | ●●● | ●● | ●●● | ●● | ●●● | ●●● | ●● | ●● | ●● | ●● | ● | ●● | ●● | ●● |
| `.dn-prompt` | ●●● | ●●● | ●● | ●●● | ●● | ●●● | ●●● | ●● | ●● | ●●● | ●● | ● | ●● | ●● | ●● |

**Odczyt macierzy.** Pięć komponentów występuje w każdym module bez wyjątku (`.dn-pasek`, `.dn-boczna`, `.dn-btn`, `.dn-aod`, `.dn-karta`) — to szkielet powłoki. Dwa komponenty nie występują w żadnym module (`.dn-listwa`, `.dn-karta-srodowiska`) — należą wyłącznie do strony głównej. `.dn-karty-sesji` nie występuje w Automations, bo Automations jest jedynym modułem bez okna w bocznej nawigacji — działa jako komponent własny wpięty w sesję innego modułu.

### 12.2. Komponenty × okna platformowe i wspólne

| Komponent | Okno startowe | Rejestracja i logowanie | Strona główna | Powłoki środowisk | Chat Window | MultitaskingAI | Okno Konfiguracji | Okno Ustawień | Always On Display | Mobile |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| `.dn-btn` | ● | ●●● | ● | ● | ● | ● | ●● | ●● | ● | ● |
| `.dn-btn-ikona` | — | ● | ● | ●●● | ●● | ●● | ● | ● | ● | ●● |
| `.dn-pole` | — | ●●● | — | — | — | ● | ●●● | ●●● | — | ● |
| `.dn-wybor` | — | ● | — | — | — | ● | ●●● | ●● | — | ● |
| `.dn-check` | — | ● | — | — | — | — | ●● | ● | — | ● |
| `.dn-radio` | — | — | — | — | — | ● | ●● | ● | — | — |
| `.dn-przelacznik` | — | — | — | — | — | ●● | ●●● | ●● | ● | ● |
| `.dn-suwak` | — | — | — | — | ● | ●● | ●● | ● | — | — |
| `.dn-szukaj` | — | — | — | ●● | ● | ● | ●● | ● | — | ● |
| `.dn-pasek` | ● | ● | ●● | ●●● | — | ●●● | ●● | ●● | — | ●● |
| `.dn-boczna` | — | — | — | ●●● | — | ●●● | ●●● | ●● | — | — |
| `.dn-karty-sesji` | — | — | — | ●●● | — | ●●● | — | — | — | ● |
| `.dn-zakladki` | — | ●● | — | ● | ● | ●● | ●●● | ●● | — | ● |
| `.dn-listwa` | — | — | ●●● | — | — | — | — | — | ● | — |
| `.dn-przybornik` | — | — | — | ● | ●●● | ●● | ● | — | ● | ●● |
| `.dn-karta` | ● | ●● | ●● | ●● | ● | ●●● | ●● | ●● | ● | ●● |
| `.dn-karta-srodowiska` | — | — | ●●● | — | — | — | — | — | — | ● |
| `.dn-kafel` | — | — | ●●● | — | — | — | — | — | — | ● |
| `.dn-tabela` | — | — | — | ● | — | ●●● | ●●● | ●● | ● | ● |
| `.dn-dane` | ● | ● | ● | ● | ●● | ●●● | ●● | ● | ●● | ● |
| `.dn-plakietka` | ● | ●● | ● | ●● | ●● | ●●● | ●●● | ●● | ●● | ●● |
| `.dn-kropka` | ●● | ● | ● | ●●● | ●●● | ●●● | ● | ● | ●●● | ●● |
| `.dn-pusty-stan` | — | — | — | ●● | ●● | ●● | ● | — | ● | ● |
| `.dn-postep` | ●●● | ● | — | ● | ● | ●●● | ● | ● | ●● | ● |
| `.dn-kolejka` / `.dn-krok` | — | ●● | — | ● | ●● | ●●● | ● | — | ●● | ● |
| `.dn-modal` | ● | ●● | ●● | ● | ● | ●● | ●●● | ●● | ● | ● |
| `.dn-toast` | ● | ●● | ● | ● | ● | ●● | ●●● | ●● | ●● | — |
| `.dn-tooltip` | — | ● | ● | ●● | ● | ●● | ●●● | ●● | ● | ● |
| `.dn-spinner` | ●●● | ●● | ● | ● | ●●● | ●● | ● | ● | ●● | ●● |
| `.dn-awatar` | — | ● | ● | ●● | ●●● | ●●● | ● | ●● | ●●● | ●● |
| `.dn-aod` | — | — | ● | ●● | ● | ●●● | ● | ● | ●●● | ●● |
| `.dn-wpis` | — | — | — | ● | ●●● | ●●● | — | — | ●● | ●●● |
| `.dn-prompt` | — | — | — | ● | ●●● | ●●● | — | — | ●● | ●●● |

**Odczyt macierzy.** Okno Konfiguracji jest największym w platformie skupiskiem kontrolek formularza — `.dn-przelacznik` (11 przełączników na jeden poziom zasięgu w macierzy izolacji), `.dn-wybor`, `.dn-tooltip` (objaśnienie `[?]` przy każdym ustawieniu) i `.dn-toast` (potwierdzenie zapisu profilu). Środowisko MultitaskingAI skupia komponenty orkiestracji: `.dn-kolejka`, `.dn-tabela`, `.dn-awatar`, `.dn-aod`. Chat Window jest jedynym miejscem `.dn-wpis` i głównym miejscem `.dn-prompt`. Mobile nie używa `.dn-toast` — okno `mobile.html` buduje własny `.mo-push` poza biblioteką, dlatego kolumna Mobile ma dla powiadomienia wartość „—".

---

## 13. Zasada zero blokad w komponentach

> **Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód. Domyślne zachowanie systemu to wykonanie polecenia.** (zasada zero blokad)

### 13.1. Co to znaczy w warstwie komponentu

| Odruch zakazany | Realizacja obowiązująca |
|---|---|
| `disabled` na przycisku, polu, przełączniku | Kontrolka zawsze aktywna; po naciśnięciu system zwraca komunikat wskazujący brak |
| `aria-disabled="true"` jako brama | Nie stosuje się; niegotowość opisuje `.dn-pole-opis`, `.dn-plakietka` albo `.dn-tooltip` |
| Krycie 50–60 % z kursorem „niedozwolone" | Nie występuje w bibliotece — nie ma reguły `opacity` na stanie niedostępności |
| Odebranie fokusu (`tabindex="-1"` jako blokada) | `tabindex="-1"` wyłącznie w obrębie wzorca pojedynczego punktu wejścia (zakładki, lista) |
| Wyszarzenie pozycji nawigacji spoza zakresu | Pozycja **nie jest renderowana** — brak modułu to mniej pozycji, nie pozycja martwa |
| Blokada przycisku zatwierdzenia przy niepełnym formularzu | Przycisk klikalny; braki ujawnia komunikat przy polach po naciśnięciu |
| Odliczanie blokujące („Wyślij ponownie za 30 s") | Odliczanie **informacyjne** — kliknięcie działa przez cały czas |
| Wyłączona metoda uwierzytelniania jako martwy segment | Metoda wyłączona w konfiguracji **nie jest renderowana wcale** |

### 13.2. Trzy dopuszczone sposoby komunikowania niegotowości

| Sposób | Komponent nośny | Kiedy |
|---|---|---|
| **Opis obok** | `.dn-pole-opis` · `.dn-kafel-opis` · `.dn-karta-cialo` | Ograniczenie znane z góry, stałe dla widoku |
| **Objaśnienie na żądanie** | `.dn-tooltip` | Ograniczenie wymagające wyjaśnienia, ale niewarte miejsca w układzie |
| **Komunikat po naciśnięciu** | `.dn-toast` · `.dn-pole-blad` · plakietka w wierszu tabeli | Warunek sprawdzalny dopiero w chwili próby |

### 13.3. Stan „odziedziczony" — wzorzec wzorcowy

Najtrudniejszy przypadek zasady zero blokad to macierz izolacji, w której wartość może pochodzić z poziomu szerszego. W innych systemach byłby to klasyczny `disabled`. Tutaj:

```
┌ wiersz macierzy izolacji ──────────────────────────────────────────┐
│ Izolacja techniczna — dostęp do sieci    [ ○——— ]  ⟨odziedziczone⟩ │
│                                             ▲          ▲            │
│                        przełącznik w pełni klikalny     plakietka   │
│                        kliknięcie tworzy jawne          wyjaśnia    │
│                        nadpisanie na tym poziomie       pochodzenie │
└─────────────────────────────────────────────────────────────────────┘
```

Wartość odziedziczona jest **stanem informacyjnym**, nie ograniczeniem. Przełącznik pozostaje interaktywny na każdym z siedmiu poziomów zasięgu; kliknięcie tworzy jawne nadpisanie. Rozróżnienie ważne dla dewelopera: „odziedziczony" (pochodzenie wartości) to co innego niż „wyłączony" (wartość logiczna `off`) — ta druga jest zwykłą, w pełni interaktywną wartością wyjściową.

### 13.4. Lista sprawdzeń dla komponentu

- [ ] Brak atrybutu `disabled` w znaczniku i brak selektora `:disabled` w arkuszu
- [ ] Każda kontrolka osiągalna klawiszem `Tab` i widoczna pierścieniem fokusu
- [ ] Niegotowość opisana tekstem, plakietką albo dymkiem — nie kryciem
- [ ] Stan komunikowany ikoną albo etykietą, nie samą barwą
- [ ] Odliczanie i limity mają charakter informacyjny
- [ ] Pozycja niedostępna w kontekście jest pominięta, nie wyszarzona
- [ ] Akcja o dużej wadze sygnalizowana wariantem `--niebezpieczny`, nie blokadą

---

## Załącznik A — ściągawka wymiarów

### A.1. Wymiary bazowe komponentów

| Komponent | Wymiar zwarty | Dotyk (`pointer: coarse`) | Gęstość przestronna | Żeton |
|---|---|---|---|---|
| `.dn-btn` | wys. 32 px | 40 px | 40 px | `--dn-wym-kontrolka` |
| `.dn-btn--sm` | wys. 28 px | 36 px | 36 px | `calc(kontrolka − od-1)` |
| `.dn-btn--lg` | wys. 40 px | 48 px | 48 px | `calc(kontrolka + od-2)` |
| `.dn-btn-ikona` | 32×32 px | 40×40 px | 40×40 px | `--dn-wym-ikonowy` |
| `.dn-pole-kontrolka` | wys. 32 px | 40 px | 40 px | `--dn-wym-kontrolka` |
| `.dn-pole-kontrolka` (textarea) | min. 64 px | 80 px | 80 px | `calc(kontrolka × 2)` |
| `.dn-check` / `.dn-radio` | 16×16 px | 20×20 px | 16×16 px | `--dn-wym-check` |
| `.dn-przelacznik` | 36×20 px | 44×24 px | 36×20 px | `--dn-wym-przelacznik-*` |
| `.dn-suwak` | tor 4 px, kciuk 16 px | bez zmian | bez zmian | — |
| `.dn-pasek` | wys. 48 px | 48 px | 56 px | `--dn-wym-pasek` |
| `.dn-boczna` | szer. 224 px | 224 px | 224 px | `--dn-wym-boczna` |
| `.dn-karty-sesji` | wys. 36 px | 36 px | 40 px | `--dn-wym-pas-kart` |
| `.dn-tabela td` | wys. 36 px | 44 px | 44 px | `--dn-wym-wiersz` |
| `.dn-tabela th` | wys. 32 px | 40 px | 40 px | `calc(wiersz − od-1)` |
| `.dn-krok` | min. 36 px | 44 px | 44 px | `--dn-wym-wiersz` |
| `.dn-modal` | maks. 560 px szer. | bez zmian | bez zmian | `--dn-wym-modal` |
| `.dn-toast` | 280–420 px szer. | bez zmian | bez zmian | `--dn-wym-toast-min/-max` |
| `.dn-awatar` | 28 px (sm 24, lg 36) | bez zmian | bez zmian | `--dn-wym-awatar*` |
| `.dn-spinner` | 14×14 px | bez zmian | bez zmian | `--dn-wym-spinner` |
| `.dn-kropka` | 6×6 px | bez zmian | bez zmian | `--dn-wym-kropka` |
| `.dn-aod-rdzen` | 36×36 px | bez zmian | bez zmian | `--dn-wym-awatar-lg` |
| Pas komunikacji | wys. 320 px | bez zmian | bez zmian | `--dn-wym-pas-komunikacji` |
| Wstęga aktywności | 2 px | bez zmian | bez zmian | `--dn-wym-wstega` |
| Pierścień fokusu | 2 px + odsunięcie 2 px | bez zmian | bez zmian | `--dn-wym-fokus*` |

### A.2. Promienie narożników → komponenty

| Żeton | Wartość | Komponenty |
|---|---|---|
| `--dn-r-xs` | 3 px | `.dn-check`, `.dn-plakietka--rola`, kod w wierszu |
| `--dn-r-sm` | 6 px | `.dn-btn`, `.dn-btn-ikona`, `.dn-pole-kontrolka`, `.dn-zakladka`, `.dn-karta-sesji`, `.dn-boczna-pozycja`, `.dn-listwa-pozycja`, `.dn-awatar--kwadrat`, `.dn-wpis-medalion`, `.dn-tooltip-tresc` |
| `--dn-r-md` | 8 px | `.dn-toast`, `.dn-wpis`, `.dn-prompt`, `.dn-kafel-ikona` |
| `--dn-r-lg` | 10 px | `.dn-karta`, `.dn-modal`, `.dn-kafel`, `.dn-listwa` |
| `--dn-r-xl` | 14 px | `.dn-karta-srodowiska` |
| `--dn-r-pill` | 999 px | `.dn-plakietka`, `.dn-kropka`, `.dn-przelacznik`, `.dn-suwak`, `.dn-postep-tor`, `.dn-awatar`, `.dn-krok-znak`, `.dn-aod`, `.dn-spinner` |

### A.3. Stopnie pisma → komponenty

| Żeton | Wartość | Komponenty |
|---|---|---|
| `--dn-fs-xs` | 12 px | `.dn-plakietka`, `.dn-tabela th`, `.dn-krok-meta`, `.dn-postep-etykieta`, `.dn-wpis-nadawca`, `.dn-wpis-godzina`, `.dn-pasek-logotyp small`, `.dn-awatar--sm` |
| `--dn-fs-sm` | 13 px | `.dn-btn--sm`, `.dn-pole-etykieta`, `.dn-pole-opis`, `.dn-pole-blad`, `.dn-karta-sesji`, `.dn-listwa-pozycja`, `.dn-toast-tresc`, `.dn-tooltip-tresc`, `.dn-pusty-stan-opis`, `.dn-aod-tresc`, `.dn-dane`, `.dn-karta-srodowiska-motto` |
| `--dn-fs-base` | 14 px | `.dn-btn`, `.dn-pole-kontrolka`, `.dn-wybor`, `.dn-zakladka`, `.dn-boczna-pozycja`, `.dn-tabela td`, `.dn-krok`, `.dn-kafel`, `.dn-toast-tytul`, `.dn-karta-srodowiska-opis` |
| `--dn-fs-md` | 15 px | `.dn-btn--lg`, `.dn-wpis-tresc`, `.dn-prompt-obszar`, `.dn-pasek-logotyp`, `.dn-awatar--lg` |
| `--dn-fs-lg` | 17 px | `.dn-karta-tytul`, `.dn-pusty-stan-tytul` |
| `--dn-fs-xl` | 21 px | `.dn-modal-tytul` |
| `--dn-fs-3xl` | 30 px | `.dn-karta-srodowiska-tytul` |

### A.4. Czasy i krzywe

| Żeton | Wartość | Zastosowanie w komponentach |
|---|---|---|
| `--dn-czas-1` | 0,1 s | Naciśnięcie (`translateY`), najechanie wiersza tabeli, zaznaczenie `.dn-check` |
| `--dn-czas-2` | 0,16 s | Przejścia barw i obrysów: przycisk, pole, karta, zakładka, przełącznik |
| `--dn-czas-3` | 0,22 s | Wejście modala i powiadomienia, zmiana szerokości `.dn-postep-wartosc` |
| `--dn-czas-tetno` | 2,4 s | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, pierścień `.dn-aod-rdzen` |
| `dn-obrot` | 0,8 s liniowo | `.dn-spinner`, wskaźnik w `.dn-btn[aria-busy]` |
| `--dn-ease` | `cubic-bezier(0.2, 0, 0, 1)` | Wszystkie przejścia i animacje bez wyjątku |

### A.5. Warstwy (`z-index`) używane przez komponenty

| Żeton | Wartość | Komponent |
|---|---|---|
| `--dn-z-przybornik` | 10 | `.dn-tabela th` (nagłówek przyklejony) |
| `--dn-z-pasek` | 100 | `.dn-pasek` |
| `--dn-z-boczna` | 200 | `.dn-boczna` |
| `--dn-z-pas-komunikacji` | 300 | Pas komunikacji z `.dn-wpis` i `.dn-prompt` |
| `--dn-z-nakladka` | 800 | `::backdrop` modala |
| `--dn-z-modal` | 900 | `.dn-modal` |
| `--dn-z-powiadomienie` | 1000 | `.dn-toasty` |
| `--dn-z-tooltip` | 1100 | `.dn-tooltip-tresc` |
| `--dn-z-aod` | 1200 | `.dn-aod` |
| `--dn-z-centrum-polecen` | 1300 | Centrum poleceń — zawsze najwyżej |

---

## Załącznik B — luki i rozbieżności katalog ↔ implementacja

Katalog komponentów v1.0 (`docs/interfejs-uzytkownika/katalog-komponentow.md`, 1500 linii) opisuje bibliotekę `system-wizualny/components.css`. Implementacja v2.0 (`zasoby/css/komponenty.css`, 1207 linii) jest **inną, nowszą generacją** o częściowo rozbieżnym nazewnictwie. Rozbieżności wypisane są tu jawnie, nie rozstrzygane milcząco. Reguła nadrzędna: **rozbieżność między opisem prozą a plikiem CSS rozstrzyga się na korzyść pliku CSS.**

### B.1. Ta sama rola, inna nazwa

| # | Rola | Katalog v1.0 | Implementacja v2.0 | Rozstrzygnięcie |
|---|---|---|---|---|
| 1 | Przycisk główny | `.dn-btn--glowny` | `.dn-btn--atrament` | v2.0 — nazwa opisuje mechanikę (inwersja atramentu), nie barwę |
| 2 | Przycisk CTA | `.dn-btn--zloty` | `.dn-btn--sygnal` | v2.0 — złoto nie istnieje w palecie v2.0 |
| 3 | Akcja ryzykowna | `.dn-btn--blad` | `.dn-btn--niebezpieczny` | v2.0 — „błąd" to stan, „niebezpieczny" to waga akcji |
| 4 | Kontrolka pola | `.dn-input` / `.dn-textarea` / `.dn-select` | jedna `.dn-pole-kontrolka` | v2.0 — jedna klasa, trzy znaczniki natywne |
| 5 | Tekst pomocy / błędu | `.dn-pomoc` / `.dn-pomoc--blad` | `.dn-pole-opis` / `.dn-pole-blad` | v2.0 |
| 6 | **Przełącznik (toggle)** | **`.dn-suwak`** | **`.dn-przelacznik`** | v2.0 — **kolizja nazwy o odwróconym znaczeniu**, patrz pozycja „Suwak zakresu" w Załączniku B |
| 7 | **Suwak zakresu (range)** | brak w katalogu | **`.dn-suwak`** | v2.0 — nazwa `.dn-suwak` znaczy w v2.0 co innego niż w v1.0; przy migracji kodu wymaga uwagi |
| 8 | Karta interaktywna | `.dn-karta--interaktywna` | `.dn-karta--klikalna` | v2.0 |
| 9 | Karta wyróżniona | `.dn-karta--akcent` | `.dn-karta--wybrana` | v2.0 — akcent v1.0 był wstęgą złotą; v2.0 rozdziela wyróżnienie od wyboru |
| 10 | Karta środowiska | `.dn-karta--interaktywna --akcent` | osobna `.dn-karta-srodowiska` | v2.0 — osobny komponent zamiast złożenia modyfikatorów |
| 11 | Kafel komponentu własnego | `.dn-karta--interaktywna` | osobna `.dn-kafel` | v2.0 |
| 12 | Treść karty / modala | `.dn-karta-tresc` / `.dn-modal-tresc` | `.dn-karta-cialo` / `.dn-modal-cialo` | v2.0 |
| 13 | Nakładka modala | klasa `.dn-nakladka` | `::backdrop` natywnego `<dialog>` | v2.0 — mniej kodu, pułapka fokusu i `Escape` za darmo |
| 14 | Boczna nawigacja | `.dn-karta--pozycja` w kolumnie | dedykowane `.dn-boczna*` | v2.0 |
| 15 | Karta sesji | `.dn-karta--pozycja` w pasku | dedykowane `.dn-karty-sesji`, `.dn-karta-sesji` | v2.0 |
| 16 | Listwa ustawień | wariant „lekki" `.dn-pasek` | dedykowane `.dn-listwa*` | v2.0 |
| 17 | Marka i rozpychacz paska | `.dn-pasek-marka`, `.dn-rozpychacz` | `.dn-pasek-godlo`, `.dn-pasek-logotyp`, `.dn-pasek-prawa` | v2.0 |
| 18 | Always On Display | `.dn-awatar` pływający | dedykowane `.dn-aod*` | v2.0 |
| 19 | Awatar jednostki AI | `--zloto` | `--inteligencja` | v2.0 |
| 20 | Warianty semantyczne | `--ostrz`, `--info` | `--ostrzezenie`, `--informacja` | v2.0 — pełne słowa polskie, bez skrótów |

### B.2. Komponenty katalogu bez odpowiednika w implementacji (luki)

| # | Komponent katalogu | Klasa v1.0 | Stan w v2.0 | Zalecane postępowanie |
|---|---|---|---|---|
| 1 | **Komunikat blokowy (Alert)** | `.dn-alert` + 4 warianty | **Brak** | Zastępczo: `.dn-pole-blad` przy polu · plakietka semantyczna w wierszu · `.dn-krok--wstrzymany` w kolejce · `.dn-toast` przy zdarzeniu. **Rekomendacja: dodać `.dn-alert` w kolejnej rewizji biblioteki** — wzorzec „stan trwały w układzie" nie ma dziś jednego nośnika |
| 2 | **Blok kodu** | `.dn-kod` | Obecny, ale w `fundament.css` (warstwa 0), nie w bibliotece | Bez zmian — blok kodu jest wzorcem tekstowym, nie komponentem interaktywnym |
| 3 | **Ikona systemowa** | `.dn-ikona` + `--sm` `--lg` | **Brak klasy** — ikony wklejane inline jako `<svg>` z `zasoby/ikony/svg/` | Bez zmian — `currentColor` i rozmiar ustawiane w miejscu użycia; klasa nie wnosiłaby nic ponad `width`/`height` |
| 4 | Zakładki pigułkowe | `.dn-zakladki--pigulki` | **Brak** | Do rozstrzygnięcia przy oknie punktów izolacji: wybór warstwy realizuje dziś para `.dn-radio` |
| 5 | Tabela naprzemienna | `.dn-tabela--paski` | **Brak** | Zgodne z kierunkiem: gęstość 8/10 i kreska subtelna wystarczą; naprzemienne tło dodałoby szumu |
| 6 | Owijka tabeli | `.dn-tabela-owijka` | **Brak** | Owijka realizowana lokalnie w oknie (przewijanie + obrys), nie w bibliotece |
| 7 | Pozycja listy | `.dn-karta--pozycja` | **Brak** | Zastąpione trzema dedykowanymi komponentami: `.dn-boczna-pozycja`, `.dn-karta-sesji`, `.dn-listwa-pozycja` |
| 8 | Stopka karty | `.dn-karta-stopka` | **Brak** | Akcje karty umieszcza się w `.dn-karta-cialo` albo w `.dn-karta-naglowek`; stopka istnieje tylko w modalu |
| 9 | Plakietka złota / wersalikowa / ze stanem | `--zloto` `--wersaliki` `--stan` | **Brak** | `--sygnal` zastępuje złotą; wersaliki daje `--rola`; „ze stanem" realizuje plakietka z kropką w środku |
| 10 | Awatar złoty | `.dn-awatar--zloto` | **Brak** | Zastąpione przez `--inteligencja` |
| 11 | Panel jako klasa | `.dn-panel` | **Brak klasy** — `--dn-panel` istnieje jako żeton tła | Panel pozostaje wzorcem złożonym, nie komponentem |
| 12 | Menu kontekstowe | `.dn-karta` + `.dn-karta--pozycja` pływające | **Brak wzorca w bibliotece** | Do zbudowania z `.dn-karta` + `.dn-listwa-pozycja` w warstwie okna; **rekomendacja: rozważyć `.dn-menu` w kolejnej rewizji** |

### B.3. Komponenty implementacji nieobecne w katalogu (nadmiar)

| # | Klasa v2.0 | Rola | Uzasadnienie istnienia |
|---|---|---|---|
| 1 | `.dn-wybor` | Wiersz etykieta + kontrolka wyboru | Ujednolica wysokość i cel kliknięcia trzech kontrolek |
| 2 | `.dn-radio` | Osobna klasa opcji jednokrotnej | Katalog łączył z `.dn-check`; rozdzielenie upraszcza arkusz |
| 3 | `.dn-suwak` (range) | Nakład rozumowania | Pojęcie produktowe nieobecne w katalogu v1.0 |
| 4 | `.dn-szukaj` | Pole wyszukiwania z ikoną | Wzorzec powtarzalny w kilkunastu panelach |
| 5 | `.dn-wpis` + 9 klas | Wpis okna komunikacji | Rozstrzygnięcie dziewięciu nadawców przez trzy klasy semantyczne |
| 6 | `.dn-postep` + 3 klasy | Pasek postępu | Monitory wykonania wymagają wskaźnika o znanym zakresie |
| 7 | `.dn-kolejka` / `.dn-krok` + 7 klas | Kolejka kroków | Trzy wyjścia weryfikacji: poprawny · błędy · wstrzymany |
| 8 | `.dn-prompt` + 2 klasy | Pole promptu z grotem | Sygnatura wejścia powtarzająca motyw godła |
| 9 | `.dn-przybornik` | Przybornik promptu | Grupowanie narzędzi pod polem wejścia |
| 10 | `.dn-plakietka--rola` | Plakietka roli okna | Coordinator · Executor · Validator |
| 11 | `.dn-kropka--tetno` | Jedyny ruch ciągły | Element sygnaturowy systemu |
| 12 | `.dn-btn--wybrany` | Stan wybrany trwale | Przyciski przełączające w przybornikach |
| 13 | `.dn-dane` | Komórka danych maszynowych | Liczby tabelaryczne w monitorach |
| 14 | `.dn-karta-srodowiska` + 4 klasy | Karta środowiska | Osobny komponent zamiast złożenia modyfikatorów |
| 15 | `.dn-kafel` + 3 klasy | Kafel komponentu własnego | Jak wyżej |
| 16 | `.dn-aod` + 2 klasy | Always On Display | Jak wyżej |

### B.4. Rozbieżności wartości liczbowych

| # | Wielkość | Katalog v1.0 | Implementacja v2.0 | Rozstrzygnięcie |
|---|---|---|---|---|
| 1 | Wysokość kontrolki | 36 px | **32 px** (`--dn-wym-kontrolka`) | v2.0 — gęstość zwarta 8/10 |
| 2 | Przycisk ikonowy | 36×36, ikona 18 | **32×32, ikona 16** | v2.0 |
| 3 | Pasek górny | 56 px | **48 px** | v2.0 (56 px pozostaje w gęstości przestronnej) |
| 4 | Przełącznik | tor 40×22, suwak 18 | **tor 36×20, suwak 14** | v2.0 |
| 5 | Awatar | 36 / 28 / 48 px | **28 / 24 / 36 px** | v2.0 |
| 6 | Wskaźnik pracy | 16 px, obrót 0,7 s | **14 px, obrót 0,8 s** | v2.0 |
| 7 | Modal — wysokość | 92 % ekranu | **`min(80dvh, 720px)`** | v2.0 — `dvh` obsługuje paski przeglądarek mobilnych |
| 8 | Powiadomienie — pozycja | dół ekranu, wyśrodkowane | **prawy dolny róg** | v2.0 — nie zasłania pasa komunikacji |
| 9 | Powiadomienie — tło | „zawsze ciemne" | **`--dn-panel`** (przełącza się z motywem) | v2.0 — ciemne pozostają wyłącznie pasek i dymek |
| 10 | Kropka statusu awatara | 11 px | **8 px** | v2.0 |
| 11 | Pusty stan — opis | 42 znaki | **40 znaków (`40ch`)** | v2.0 |
| 12 | Przycisk — stopień i waga | 13 px, `semibold` | **13 px (`--dn-fs-base`), waga 500** | v2.0 — skala v2.0 ma inną bazę niż v1.0 |
| 13 | Skala typografii | 12·13·15·16·18·22·28·36·46 | **12·13·14·15·17·21·24·30·40** | v2.0 — skala zwarta kokpitu |
| 14 | Wartości pikselowe w arkuszu | podawane wprost | **wyłącznie żetony** | v2.0 — „zero wartości zaszytych" |
| 15 | Tętno kropki | brak w katalogu | **2,4 s** (`--dn-czas-tetno`) | v2.0 — element sygnaturowy |

### B.5. Zgodność bez rozbieżności

Zasada zero blokad z rozdziału 15 katalogu v1.0 jest w implementacji v2.0 zapisana wprost jako **zasadę zero blokad** w nagłówku arkusza. Obie generacje zgadzają się co do: braku stanu wyłączonego, opisowego komunikowania niegotowości, zakazu komunikowania stanu samą barwą, pełnej kompozycyjności komponentów i pierwszeństwa plików CSS przed opisem prozą.

---

## Załącznik C — pełna lista klas CSS

### C.1. Klasy biblioteki komponentów (`zasoby/css/komponenty.css`)

| Komponent | Klasa bazowa | Modyfikatory | Elementy wewnętrzne | Rozdział |
|---|---|---|---|---|
| Przycisk | `.dn-btn` | `--atrament` `--sygnal` `--zarys` `--duch` `--niebezpieczny` `--wybrany` `--sm` `--lg` | — | 5.1 |
| Przycisk ikonowy | `.dn-btn-ikona` | `--na-ramie` | — | 5.2 |
| Pole formularza | `.dn-pole` | — | `-etykieta` `-kontrolka` `-opis` `-blad` | 6.1 |
| Wiersz wyboru | `.dn-wybor` | — | — | 6.2 |
| Pole wyboru | `.dn-check` | — | — | 6.3 |
| Opcja jednokrotna | `.dn-radio` | — | — | 6.4 |
| Przełącznik | `.dn-przelacznik` | — | — | 6.5 |
| Suwak zakresu | `.dn-suwak` | — | — | 6.6 |
| Pole wyszukiwania | `.dn-szukaj` | — | — | 6.7 |
| Pasek górny | `.dn-pasek` | — | `-godlo` `-logotyp` `-szukaj` `-prawa` | 7.1 |
| Boczna nawigacja | `.dn-boczna` | — | `-naglowek` `-pozycja` | 7.2 |
| Pas kart sesji | `.dn-karty-sesji` | — | — | 7.3 |
| Karta sesji | `.dn-karta-sesji` | — | — | 7.3 |
| Zakładki | `.dn-zakladki` | — | — | 7.4 |
| Zakładka | `.dn-zakladka` | — | — | 7.4 |
| Listwa ustawień | `.dn-listwa` | — | `-pozycja` | 7.5 |
| Przybornik | `.dn-przybornik` | — | — | 7.6 |
| Karta | `.dn-karta` | `--klikalna` `--wybrana` | `-naglowek` `-tytul` `-cialo` | 8.1 |
| Karta środowiska | `.dn-karta-srodowiska` | — | `-godlo` `-tytul` `-motto` `-opis` | 8.2 |
| Kafel | `.dn-kafel` | — | `-ikona` `-etykieta` `-opis` | 8.3 |
| Tabela | `.dn-tabela` | — | — | 8.4 |
| Komórka danych | `.dn-dane` | — | — | 8.5 |
| Plakietka | `.dn-plakietka` | `--sukces` `--ostrzezenie` `--blad` `--informacja` `--sygnal` `--rola` | — | 8.6 |
| Kropka sygnału | `.dn-kropka` | `--tetno` `--sukces` `--ostrzezenie` `--blad` `--neutralna` | — | 8.7 |
| Pusty stan | `.dn-pusty-stan` | — | `-tytul` `-opis` | 8.8 |
| Pasek postępu | `.dn-postep` | — | `-etykieta` `-tor` `-wartosc` | 8.9 |
| Kolejka | `.dn-kolejka` | — | — | 8.10 |
| Krok | `.dn-krok` | `--pracuje` `--poprawny` `--bledy` `--wstrzymany` | `-znak` `-meta` | 8.10 |
| Modal | `.dn-modal` | — | `-naglowek` `-tytul` `-cialo` `-stopka` | 9.1 |
| Stos powiadomień | `.dn-toasty` | — | — | 9.2 |
| Powiadomienie | `.dn-toast` | `--sukces` `--ostrzezenie` `--blad` `--informacja` | `-tytul` `-tresc` | 9.2 |
| Dymek objaśnienia | `.dn-tooltip` | — | `-tresc` | 9.3 |
| Wskaźnik pracy | `.dn-spinner` | — | — | 9.4 |
| Awatar | `.dn-awatar` | `--sm` `--lg` `--kwadrat` `--inteligencja` | `-stan` | 10.1 |
| Always On Display | `.dn-aod` | — | `-rdzen` `-tresc` | 10.2 |
| Wpis komunikacji | `.dn-wpis` | `--czlowiek` `--inteligencja` `--system` `--pracuje` | `-medalion` `-tozsamosc` `-nadawca` `-godzina` `-tresc` | 11.1 |
| Pole promptu | `.dn-prompt` | — | `-grot` `-obszar` | 11.2 |

### C.2. Klasy warstwy fundamentu (`zasoby/css/fundament.css`) — poza biblioteką

| Klasa | Rola | Uwaga |
|---|---|---|
| `.dn-naglowek` | Krój nagłówkowy dla elementu spoza `h1`–`h3` | Warstwa 0 |
| `.dn-etykieta-wersalikowa` | Etykieta metadanych krojem bazowym | Warstwa 0 |
| `.dn-etykieta-mono` | Etykieta metadanych krojem mono | Warstwa 0 |
| `.dn-dane` · `.dn-liczba` | Krój mono + liczby tabelaryczne | Uzupełniane w bibliotece dla komórek tabeli |
| `.dn-kod` · `.dn-kod--wiersz` | Blok kodu i kod w wierszu | Odpowiednik `.dn-kod` z katalogu v1.0 (luka odnotowana w Załączniku B) |
| `.dn-kbd` | Klawisz | Nieobecny w katalogu v1.0 |
| `.dn-separator` · `.dn-separator--pionowy` | Separator | Odpowiednik separatora z katalogu v1.0 |
| `.dn-sr-only` | Ukrycie dostępne — treść dla czytnika ekranu | Warunek dostępności komponentów |

### C.3. Animacje zdefiniowane w bibliotece

| Nazwa | Czas | Komponenty |
|---|---|---|
| `dn-tetno` | `--dn-czas-tetno` 2,4 s | `.dn-kropka--tetno`, `.dn-wpis--pracuje`, `.dn-aod-rdzen::after` |
| `dn-wejscie` | `--dn-czas-3` 0,22 s | `.dn-modal[open]`, `.dn-toast` |
| `dn-obrot` | 0,8 s liniowo | `.dn-spinner`, `.dn-btn[aria-busy='true']::after` |

### C.4. Rachunek klas

| Warstwa | Klasy bazowe | Modyfikatory | Elementy | Razem |
|---|---:|---:|---:|---:|
| Biblioteka komponentów | 33 | 26 | 43 | **102** |
| Fundament | 8 | 2 | — | **10** |
| Klasy pomocnicze rozliczone podwójnie | −1 (`.dn-dane`) | — | — | **−1** |
| **Suma nazw `.dn-*` w arkuszach projektu** | | | | **111** |

Wykaz alfabetyczny wszystkich selektorów w `komponenty.css` obejmuje 119 tokenów klasowych — różnica względem 111 wynika z ujęcia w wykazie klas występujących wyłącznie jako selektory złożone (np. `.dn-tabela td.dn-dane`).

---

## 17. Decyzje projektowe

Rozstrzygnięcia podjęte ponad źródła — każde z uzasadnieniem i wskazaniem, co rozstrzygało.

| # | Decyzja | Uzasadnienie | Co rozstrzygało |
|---|---|---|---|
| **1** | **Kategoria G nazwana „Komunikacja", nie „Elementy pomocnicze"** | Katalog v1.0 umieszczał w G ikonę systemową i separator; oba przeszły w v2.0 do warstwy fundamentu (`.dn-separator`) albo zniknęły jako klasa (ikony inline). Miejsce po nich zajmują dwa komponenty rozmowy, których katalog nie znał: `.dn-wpis` i `.dn-prompt` | kontrakt systemu projektowego grupuje je pod nagłówkiem „Komunikacja" |
| **2** | **`.dn-btn` opisany jako 7 wariantów + 1 modyfikator stanu** | Arkusz definiuje osiem modyfikatorów, ale `--wybrany` nie jest wariantem wyglądu — jest zapisem stanu równoważnym `aria-pressed='true'` (obie reguły dzielą jeden blok CSS) | Struktura arkusza: `--wybrany` stoi w bloku stanów, przed wariantami |
| **3** | **`.dn-radio` opisane jako osobna karta komponentu** | Katalog v1.0 łączył checkbox i radio w jeden komponent B4. Implementacja v2.0 ma dwie osobne klasy o różnym kształcie i różnej mechanice wyboru; jedna karta zaciemniałaby różnicę w obsłudze klawiatury (spacja vs strzałki) | Arkusz: `.dn-radio` to osobny selektor z własnymi regułami |
| **4** | **`.dn-toasty` opisany wewnątrz karty `.dn-toast`, nie jako osobny komponent** | Stos nie ma sensu bez powiadomienia i nie występuje samodzielnie; rozdzielenie dałoby kartę bez treści merytorycznej | Analogia do `.dn-kolejka`/`.dn-krok`, gdzie oba opisano łącznie |
| **5** | **`.dn-kolejka` i `.dn-krok` opisane w jednej karcie** | Krok nie występuje poza kolejką; stany kroku są jednocześnie wariantami — rozbicie na dwie karty powielałoby tabelę wariantów | Struktura arkusza: wspólny blok „Monitor wykonania · kolejka kroków" |
| **6** | **`.dn-przybornik` zaliczony do kategorii C (nawigacja), mimo sąsiedztwa z `.dn-prompt` w arkuszu** | Przybornik jest kontenerem narzędzi, nie elementem rozmowy; występuje także poza polem promptu (paski narzędzi paneli) | kontrakt systemu projektowego wymienia go w grupie „Nawigacja" |
| **7** | **Rekomendacja dodania `.dn-alert` (komunikat blokowy) w kolejnej rewizji** | Katalog v1.0 opisuje komunikat blokowy wywoływany z dziewięciu paneli; v2.0 nie ma jego odpowiednika. Dziś rolę pełnią zastępczo cztery różne komponenty, co rozprasza wzorzec „stan trwały w układzie" | Luka odnotowana w Załączniku B; rekomendacja **nie zmienia** stanu biblioteki — jest odnotowana, nie wprowadzona |
| **8** | **Rekomendacja rozważenia `.dn-menu` (menu kontekstowe) w kolejnej rewizji** | Katalog v1.0 opisuje menu kontekstowe obecne „w całej platformie"; v2.0 nie ma dla niego ani klasy, ani wskazanego złożenia | Luka odnotowana w Załączniku B; jak wyżej — odnotowana, nie wprowadzona |
| **9** | **Brak klasy `.dn-ikona` uznany za rozstrzygnięcie, nie brak** | Ikony wklejane inline jako `<svg>` z `currentColor` nie potrzebują klasy; rozmiar ustawia atrybut w miejscu użycia. Klasa dodałaby warstwę pośrednią bez korzyści | Praktyka wszystkich arkuszy i makiet projektu; kontrakt systemu projektowego |
| **10** | **Wymiary podawane w wartościach gęstości zwartej z jawną kolumną wariantów** | Żetony zmieniają wartość w trzech kontekstach (dotyk, gęstość przestronna, ograniczony ruch); podanie jednej liczby bez kontekstu wprowadzałoby w błąd | Żetony rozdz. 12–13 pliku `zetony.css` |
| **11** | **Stany wymuszane w galerii atrybutem `data-stan`, nie klasami modyfikującymi** | Biblioteka nie ma klas stanu; dodanie ich na potrzeby galerii złamałoby regułę „nie zmieniać nazw i nie mnożyć klas". Atrybut należy do warstwy dokumentacyjnej `.dok-*` i nie zanieczyszcza biblioteki | Zasada wiążących nazw klas (rozdz. 1.2) |
| **12** | **Kolejność kategorii A→G przyjęta za katalogiem v1.0, mimo zmiany zawartości G** | Zachowanie kolejności pozwala czytać oba dokumenty równolegle; zmiana kolejności zerwałaby odsyłacze między nimi | Katalog v1.0 rozdz. 4 |
| **13** | **Macierz rozbita na dwie tabele (moduły osobno, okna platformowe osobno)** | Jedna tabela o 25 kolumnach byłaby nieczytelna; podział przebiega po naturalnej granicy: moduł ma boczną nawigację, okno platformowe jej nie ma | Inwentarz okien, rozdz. 7.1–7.5 |
| **14** | **Nasilenie w macierzy oznaczone czterostopniowo (`—` `●` `●●` `●●●`)** | Skala przejęta z katalogu v1.0 rozdz. 14 dla zachowania porównywalności; wartości dla komponentów nieobecnych w katalogu ustalono na podstawie inwentarza okien i specyfikacji modułów | Katalog v1.0, legenda rozdz. 14 |
| **15** | **Treści przykładowe wyłącznie z domeny produktu** | Nazwy sesji, identyfikatory zadań (`QUE-0142`), nazwy ról (Coordinator, Executor 1), nazwy okien (Workflow Builder, Queue Manager) i ścieżki repozytoriów pochodzą z dokumentacji platformy i są oznaczone jako przykładowe | kontrakt systemu projektowego |
| **16** | **Kolizja nazwy `.dn-suwak` odnotowana jako ryzyko migracji, nie naprawiona** | Zmiana nazwy klasy złamałaby zasadę wiążących nazw klas („nazwy klas są wiążące"). Ryzyko jest realne: ten sam zapis znaczy w v1.0 przełącznik, a w v2.0 suwak zakresu | Rozbieżność odnotowana w Załączniku B; decyzja o nienaprawianiu należy do właściciela biblioteki |

---

*Danaco Console — AI Operating Environment · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. — Dariusz Naharnowicz*
