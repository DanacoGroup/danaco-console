# Danaco Console — Always On Display

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
| **Tytuł** | Always On Display |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Pełnozakresowy materiał źródłowy do projektowania i budowy funkcji globalnej Always On Display obowiązującej we wszystkich środowiskach i modułach platformy: co istnieje, gdzie leży, w jakiej formie, kiedy się ujawnia i do czego służy |
| **Zakres** | warstwa funkcji globalnych — Always On Display: postać wizualna, reguły wyzwalania sugestii, katalog sugestii, tor głosowy, zachowanie per środowisko i per moduł, warstwy widoczności, stany i dane, punkty sterowania z okna konfiguracji i z okna ustawień, scenariusze użycia, wykazy normatywne, kryteria odbioru |
| **Poza zakresem** | dialog głosowy i tekstowy wieloetapowy jako taki — [Moduł Assistant](../moduly/assistant.md); środowisko MultitaskingAI, w którym funkcja pełni dodatkowo rolę warstwy centralnej — [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) rozdz. 2 |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) · [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) · [Elementy okien](../interfejs-uzytkownika/elementy-okien.md) · [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) · [Ustawienia](../interfejs-uzytkownika/ustawienia.md) · [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) · [Model danych](../architektura/model-danych.md) · [Model konfiguracji](../architektura/model-konfiguracji.md) · [Moduł Assistant](../moduly/assistant.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) |
| **Prototypy odniesienia** | `design/05-okna/platformowe/always-on-display.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `aod`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `design/zasoby/rama.css`, `design/zasoby/prototyp.css` · `design/05-okna/platformowe/always-on-display.html` |
| **Zasada nadrzędna** | Danaco Console nie narzuca twardych blokad, bram bezpieczeństwa ani wymuszonych zgód; domyślne zachowanie systemu to wykonanie polecenia; żaden element interfejsu nie traci klikalności — niekompletność sygnalizowana jest komunikatem kontekstowym, nie wyłączeniem kontrolki; zakres działania Always On Display, reguły wyzwalania sugestii i obecność awatara są ustawieniami konfiguracyjnymi Operatora ze stanem wyjściowym „pełny dostęp” |

---

## Spis treści

1. [Wprowadzenie](#wprowadzenie)
2. [Przeznaczenie i kontekst funkcji](#1-przeznaczenie-i-kontekst-funkcji)
   - [1.1 Definicja](#11-definicja)
   - [1.2 Charakter przekrojowy](#12-charakter-przekrojowy)
   - [1.3 Zakres dostępu](#13-zakres-dostępu)
   - [1.4 Odróżnienie od sąsiednich bytów platformy](#14-odróżnienie-od-sąsiednich-bytów-platformy)
3. [Postać wizualna i umiejscowienie w układzie](#2-postać-wizualna-i-umiejscowienie-w-układzie)
   - [2.1 Zgodność z układem pionowym](#21-zgodność-z-układem-pionowym)
   - [2.2 Makieta — stan bezczynny](#22-makieta--stan-bezczynny)
   - [2.3 Makieta — powierzchnia interakcji otwarta](#23-makieta--powierzchnia-interakcji-otwarta)
   - [2.4 Awatar — forma, waga i stany](#24-awatar--forma-waga-i-stany)
   - [2.5 Zachowanie przy zwężaniu kolumn](#25-zachowanie-przy-zwężaniu-kolumn)
   - [2.6 Katalog elementów funkcji](#26-katalog-elementów-funkcji)
4. [Reguły wyzwalania proaktywnych sugestii](#3-reguły-wyzwalania-proaktywnych-sugestii)
   - [3.1 Zasada wyzwalania](#31-zasada-wyzwalania)
   - [3.2 Klasy zdarzeń wyzwalających](#32-klasy-zdarzeń-wyzwalających)
   - [3.3 Źródła sygnału](#33-źródła-sygnału)
   - [3.4 Progi i częstotliwość](#34-progi-i-częstotliwość)
   - [3.5 Wyciszanie](#35-wyciszanie)
   - [3.6 Zachowanie przy braku istotnych zdarzeń](#36-zachowanie-przy-braku-istotnych-zdarzeń)
5. [Katalog rodzajów sugestii](#4-katalog-rodzajów-sugestii)
   - [4.1 Doradztwo (`doradztwo`)](#41-doradztwo-doradztwo)
   - [4.2 Sugestia konfiguracji (`konfiguracja`)](#42-sugestia-konfiguracji-konfiguracja)
   - [4.3 Wskazanie problemu (`problem`)](#43-wskazanie-problemu-problem)
   - [4.4 Kolejny krok (`kolejny_krok`)](#44-kolejny-krok-kolejny_krok)
   - [4.5 Waga sugestii i sposób ujawnienia](#45-waga-sugestii-i-sposób-ujawnienia)
6. [Tor głosowy i granica wobec modułu Assistant](#5-tor-głosowy-i-granica-wobec-modułu-assistant)
   - [5.1 Podział zakresów](#51-podział-zakresów)
   - [5.2 Która powierzchnia jest właściwa dla którego polecenia](#52-która-powierzchnia-jest-właściwa-dla-którego-polecenia)
   - [5.3 Przekazanie między torami](#53-przekazanie-między-torami)
   - [5.4 Rozstrzygnięcie pierwszeństwa](#54-rozstrzygnięcie-pierwszeństwa)
7. [Zachowanie per środowisko i per moduł](#6-zachowanie-per-środowisko-i-per-moduł)
   - [6.1 Zachowanie podstawowe](#61-zachowanie-podstawowe)
   - [6.2 Tabela zachowań](#62-tabela-zachowań)
   - [6.3 MultitaskingAI jako przypadek szczególny](#63-multitaskingai-jako-przypadek-szczególny)
8. [Relacja do dwóch kanałów komunikacji operacyjnej](#7-relacja-do-dwóch-kanałów-komunikacji-operacyjnej)
   - [7.1 Miejsce funkcji wobec obu kanałów](#71-miejsce-funkcji-wobec-obu-kanałów)
   - [7.2 Sugestia przeniesiona do Chat Window jako polecenie](#72-sugestia-przeniesiona-do-chat-window-jako-polecenie)
   - [7.3 Sygnały pochodzące z Execution Loop Window](#73-sygnały-pochodzące-z-execution-loop-window)
   - [7.4 Rozgraniczenie odpowiedzialności](#74-rozgraniczenie-odpowiedzialności)
9. [Warstwy widoczności](#8-warstwy-widoczności)
   - [8.1 Przydział warstw](#81-przydział-warstw)
   - [8.2 Zasada jednego kliknięcia](#82-zasada-jednego-kliknięcia)
   - [8.3 Mechanizmy ukrywania zastosowane w funkcji](#83-mechanizmy-ukrywania-zastosowane-w-funkcji)
10. [Stany, dane i powiązania](#9-stany-dane-i-powiązania)
   - [9.1 Stany funkcji](#91-stany-funkcji)
   - [9.2 Tryby obecności](#92-tryby-obecności)
   - [9.3 Encje modelu danych](#93-encje-modelu-danych)
   - [9.4 Cykl życia sugestii](#94-cykl-życia-sugestii)
   - [9.5 Komunikaty kontraktu komunikacji](#95-komunikaty-kontraktu-komunikacji)
   - [9.6 Powiązania](#96-powiązania)
11. [Punkty sterowania z okna konfiguracji i z okna ustawień](#10-punkty-sterowania-z-okna-konfiguracji-i-z-okna-ustawień)
   - [10.1 Punkty sterowania z okna konfiguracji](#101-punkty-sterowania-z-okna-konfiguracji)
   - [10.2 Punkty sterowania z okna Ustawień](#102-punkty-sterowania-z-okna-ustawień)
   - [10.3 Sterowanie poleceniem języka naturalnego](#103-sterowanie-poleceniem-języka-naturalnego)
12. [Scenariusze użycia](#11-scenariusze-użycia)
   - [11.1 Wynik kontroli jakości odrzucony po raz drugi](#111-wynik-kontroli-jakości-odrzucony-po-raz-drugi)
   - [11.2 Powtarzalność czynności ręcznej](#112-powtarzalność-czynności-ręcznej)
   - [11.3 Krótkie polecenie głosowe poza modułem Assistant](#113-krótkie-polecenie-głosowe-poza-modułem-assistant)
   - [11.4 Polecenie głosowe przekazane do modułu Assistant](#114-polecenie-głosowe-przekazane-do-modułu-assistant)
   - [11.5 Punkt decyzyjny pętli podczas pracy ciągłej](#115-punkt-decyzyjny-pętli-podczas-pracy-ciągłej)
   - [11.6 Praca wymagająca ciszy](#116-praca-wymagająca-ciszy)
   - [11.7 Nadzór ciągły nad procesem od uruchomienia po zamknięcie zlecenia](#117-nadzór-ciągły-nad-procesem-od-uruchomienia-po-zamknięcie-zlecenia)
13. [Załącznik A. Skróty klawiszowe i ikonografia](#załącznik-a-skróty-klawiszowe-i-ikonografia)
   - [A.1. Skróty klawiszowe](#a1-skróty-klawiszowe)
   - [A.2. Ikonografia](#a2-ikonografia)
   - [A.3. Komponenty systemu wizualnego](#a3-komponenty-systemu-wizualnego)
14. [Załącznik B. Pełny wykaz komend kontraktu Always On Display](#załącznik-b-pełny-wykaz-komend-kontraktu-always-on-display)
   - [Obszar `aod` — 12 komend](#obszar-aod--12-komend)
15. [Załącznik C. Wykazy normatywne uzupełniające i kryteria odbioru](#załącznik-c-wykazy-normatywne-uzupełniające-i-kryteria-odbioru)
   - [C.1. Etykiety interfejsu z prototypu](#c1-etykiety-interfejsu-z-prototypu)
   - [C.2. Żetony projektowe](#c2-żetony-projektowe)
   - [C.3. Punkty łamania](#c3-punkty-łamania)
   - [C.4. Stany kontrolek](#c4-stany-kontrolek)
   - [C.5. Kryteria odbioru](#c5-kryteria-odbioru)

---

## Wprowadzenie

Always On Display jest globalnym agentem towarzyszącym Operatora. Nie posiada własnego środowiska ani modułu, nie tworzy przestrzeni roboczej i nie zakłada nowych sesji; obecny jest jednocześnie we wszystkich częściach platformy jako nadrzędna warstwa inteligencji ponad środowiskami, modułami, projektami i sesjami ([Koncepcja platformy](../architektura/koncepcja-platformy.md), rozdz. 5.3, 12.2). Kluczowa idea funkcji brzmi: osobisty agent operacyjny obecny zawsze i wszędzie.

Funkcja jest przekrojowa wobec całej platformy. Ten sam awatar, ten sam mechanizm sugestii, ten sam tor głosowy i ta sama powierzchnia interakcji obowiązują w czterech środowiskach (TalkIn, WorkSpace, CodeStudio, MultitaskingAI), w piętnastu modułach oraz na stronie głównej. Środowisko MultitaskingAI jest przypadkiem szczególnym: funkcja pełni tam dodatkowo rolę warstwy centralnej nadzorującej zespół ról i silnik kolejek — specyfikę tej roli opisuje [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) (rozdz. 2), a niniejszy dokument obowiązuje w pozostałym zakresie bez zmian.

Dwa kanały komunikacji operacyjnej są dla Always On Display źródłem sygnału i miejscem skutku. **Chat Window** — kanał Użytkownik ↔ Wykonawca — przyjmuje przeniesione sugestie jako polecenia i pozostaje centralnym punktem pracy oraz podstawowym mechanizmem sterowania wszystkimi procesami platformy. **Execution Loop Window** — kanał Koordynator ↔ Wykonawca — dostarcza stan pętli wykonawczej, punkty decyzyjne oczekujące na zatwierdzenie oraz wyniki kontroli jakości, na których opiera się wyzwalanie sugestii.

Dokument czyta się w czterech warstwach nałożonych na siebie:

| Warstwa dokumentu | Co dostarcza | Rozdziały |
|---|---|---|
| Kontekst | Miejsce funkcji w architekturze platformy i jej charakter przekrojowy | 1 |
| Forma | Postać wizualna, umiejscowienie w układzie kolumnowym, warstwy widoczności | 2, 8 |
| Zachowanie w czasie | Wyzwalanie sugestii, katalog sugestii, tor głosowy, zachowanie per środowisko, relacja do kanałów, stany | 3–7, 9 |
| Sterowanie | Punkty konfiguracji i ustawień, scenariusze, mapa skrótów | 10–11, załącznik A |

---

## 1. Przeznaczenie i kontekst funkcji

### 1.1. Definicja

Always On Display jest funkcją globalną platformy: funkcjonalnością dostępną niezależnie od aktywnego środowiska lub modułu, nietworzącą własnej przestrzeni roboczej. Realizuje trzy zadania jednocześnie:

| Zadanie | Treść |
|---|---|
| Stała obecność | Punkt dostępu do agenta towarzyszącego widoczny w każdym miejscu platformy, niezależny od otwartego okna, karty sesji i modułu |
| Proaktywne doradztwo | Wykrywanie zdarzeń istotnych dla Operatora — próg istotności każdej klasy zdarzeń podaje rozdz. 3.4 — i przedstawianie sugestii działania, konfiguracji, wskazań problemów oraz kolejnego kroku |
| Sterowanie ponadkontekstowe | Przyjęcie polecenia tekstowego i głosowego kierowanego do platformy jako całości, poza kontekstem pojedynczej karty sesji |

### 1.2. Charakter przekrojowy

Przynależność do warstwy funkcji globalnych (Global Features Layer) ma trzy skutki wiążące dla projektu interfejsu:

| Cecha | Skutek projektowy |
|---|---|
| Brak własnej powłoki roboczej | Funkcja nie ma strony, ekranu ani zestawu okien operacyjnych; jej całą powierzchnią jest awatar oraz powierzchnia interakcji otwierana jako rozszerzenie boczne |
| Brak przeładowania | Przełączenie środowiska, modułu lub karty sesji nie przerywa działania funkcji, nie zamyka powierzchni interakcji i nie kasuje kolejki oczekujących sugestii ([Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md), rozdz. 5.6) |
| Brak przypisania do kontekstu | Funkcja nie należy do żadnego środowiska ani modułu; czyta kontekst wszystkich z nich i odnosi się do niego, sama pozostając ponad nim |

### 1.3. Zakres dostępu

| Obszar dostępu | Zakres |
|---|---|
| Środowiska | Wszystkie cztery, wraz ze stanem ich okien operacyjnych |
| Moduły | Wszystkie piętnaście, wraz z bieżącym stanem zadań i zasobów |
| Projekty i sesje | Wszystkie projekty, karty sesji i ich historia |
| Historia rozmów | Pełna historia komunikacji Operatora w obu kanałach komunikacji operacyjnej |
| Agenci | Agenci utworzeni w module Agents, dostępni jako wykonawcy sugerowanych działań |
| Procesy | Pętle wykonawcze, kolejki, automatyki, przebiegi środowiska MultitaskingAI |

### 1.4. Odróżnienie od sąsiednich bytów platformy

| Byt | Odróżnienie |
|---|---|
| Funkcja globalna Mobile | Mobile jest globalnym trybem dostępu do platformy z urządzenia poza stanowiskiem roboczym; Always On Display jest agentem towarzyszącym. Mobile udostępnia kanał wykonania interwencji zainicjowanej przez Always On Display (rozdz. 6.7) |
| Moduł Assistant | Assistant jest pełnoprawnym modułem osadzonym w środowisku, z własną powłoką, kartą sesji, pamięcią i siedmioma oknami operacyjnymi; Always On Display jest warstwą ponad platformą bez sesji własnej. Granicę toru głosowego rozstrzyga rozdz. 5 |
| Chat Window | Chat Window jest kanałem komunikacji operacyjnej wewnątrz bieżącej karty sesji i miejscem wykonania polecenia; Always On Display kieruje do niego polecenia powstałe z sugestii, sam nie zastępując tego kanału |
| Awatar tożsamościowy | Awatar Operatora w pasku stanu (i w stopce szyny nawigacji) identyfikuje konto ([Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md), rozdz. 12.1); awatar Always On Display jest punktem dostępu do funkcji (tamże, rozdz. 12.2) |

---

## 2. Postać wizualna i umiejscowienie w układzie

### 2.1. Zgodność z układem pionowym

Cała platforma stosuje układ pionowy (podział lewa–prawa): okna robocze, okna komunikacji i okna pomocnicze sąsiadują ze sobą poziomo w kolumnach o regulowanej szerokości. Always On Display wpisuje się w ten układ dwoma elementami: pływającym awatarem zakotwiczonym przy prawej krawędzi obszaru roboczego oraz powierzchnią interakcji otwieraną jako kolumna boczna po prawej stronie obszaru roboczego — rozszerzenie boczne o regulowanej szerokości.

| Element funkcji | Waga wizualna | Umiejscowienie |
|---|---|---|
| Awatar Always On Display | Element pływający, waga znikoma w stanie bezczynnym | Zakotwiczony przy prawej krawędzi obszaru roboczego, ponad całą powłoką aplikacji (pozycjonowanie stałe) |
| Powierzchnia interakcji | Kolumna boczna, otwierana jako rozszerzenie boczne | Po prawej stronie obszaru roboczego, na pełną wysokość, szerokość regulowana |
| Dymek kontekstowy sugestii | Popover średniej wielkości, waga niska | Przy awatarze, po jego lewej stronie, bez przesuwania kolumn obszaru roboczego |
| Plakietka powiadomień | Plakietka liczbowa, waga znikoma | Na awatarze |

Otwarcie powierzchni interakcji zwęża kolumny obszaru roboczego, nie przesłania ich i nie zamyka żadnego z okien komunikacji operacyjnej: Chat Window pozostaje w lewej kolumnie, Execution Loop Window w kolumnie sąsiadującej.

### 2.2. Makieta — stan bezczynny

```
Makieta 1 — Always On Display w stanie bezczynnym (warstwa 1 oraz zwinięte wyzwalacze warstw 2–3)

 ══════════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]      [⋮]
 ══════════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu           │ Panel      │
  nawigacja │ Użytkownik ↔         │                                 │ pomocniczy │
  modułów   │ Wykonawca            │                                 │ ▼          │
            │                      │                                 │            │
            │ ──────────────────   │                                 │        ╭───╮
            │ Execution Loop ▼     │                                 │        │ ● │ ◄ awatar AOD
            │ Koordynator ↔        │                                 │        │ ② │ ◄ plakietka
            │ Wykonawca            │                                 │        ╰───╯
            │                      │                                 │            │
            │ [ polecenie… ]       │                                 │            │
 ══════════════════════════════════════════════════════════════════════════════════
```

Legenda: w stanie bezczynnym widoczne są wyłącznie elementy warstwy 1 — awatar oraz
plakietka liczby oczekujących sugestii — a także zwinięte wyzwalacze warstw 2–3
(`▼`, `⋮`, `☰`). Dymek sugestii, powierzchnia interakcji, przycisk mikrofonu
i przełącznik trybu pozostają niewidoczne do chwili wywołania.

### 2.3. Makieta — powierzchnia interakcji otwarta

```
Makieta 2 — Powierzchnia interakcji jako rozszerzenie boczne (po aktywacji awatara)

 ═══════════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window       │ Obszar roboczy      │ Always On Display          │
  nawigacja │ Użytkownik ↔      │ modułu              │ ─────────────────────────  │
  modułów   │ Wykonawca         │ (kolumny zwężone)   │ [ tryb ▼ ]  [ 🎙 ]   [ ⋮ ]  │
            │                   │                     │                            │
            │ ───────────────   │                     │  strumień rozmowy z AOD    │
            │ Execution Loop ▼  │                     │  …                         │
            │ Koordynator ↔     │                     │                            │
            │ Wykonawca         │                     │  ── oczekujące sugestie ── │
            │                   │                     │  ▸ kontrola jakości        │
            │ [ polecenie… ]    │                     │  ▸ kolejka wstrzymana      │
            │                   │                     │  [ polecenie do AOD… ]     │
 ═══════════════════════════════════════════════════════════════════════════════════
```

### 2.4. Awatar — forma, waga i stany

| Cecha | Wartość |
|---|---|
| Komponent | `.dn-aod` (wrapper pozycjonowania pływającego) zawierający `.dn-awatar`/`.dn-aod-rdzen` w wariancie `--lg` (48 px) oraz `.dn-aod-tresc` (tekst towarzyszący w stanach z treścią) — `design/zasoby/css/komponenty.css` |
| Waga wizualna w stanie bezczynnym | Znikoma — pojedyncze koło, bez etykiety tekstowej, bez ramki kontenera |
| Waga wizualna po aktywacji | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Zachowanie przy przełączeniu środowiska lub modułu | Bez zmian — awatar nie znika, nie zmienia pozycji i nie podlega przeładowaniu |

| Stan awatara | Postać wizualna | Znaczenie |
|---|---|---|
| Bezczynny | Awatar w barwie bazowej, bez animacji | Brak oczekującej sugestii, gotowość do interakcji |
| Sugestia oczekująca | Plakietka liczbowa `.dn-plakietka` z liczbą pozycji | Jedna lub więcej sugestii oczekuje na zapoznanie |
| Sugestia o wadze wysokiej | Plakietka pulsująca w akcencie złotym | Zdarzenie wymagające decyzji Operatora — punkt decyzyjny pętli albo niepowodzenie kontroli jakości |
| Mówi lub pisze | Obwódka w akcencie złotym, animacja przebiegu | Trwa synteza mowy albo strumieniowanie odpowiedzi |
| Nasłuch | Pierścień pulsujący, ikona mikrofonu na awatarze | Trwa odbiór polecenia głosowego |
| Wyciszony | Awatar w barwie przygaszonej, plakietka ukryta | Wyciszenie sugestii aktywne (rozdz. 3.5) |
| Ukryty | Awatar poza polem widzenia, funkcja czynna w tle | Obecność awatara wyłączona ustawieniem; dostęp przez szybki wybór szyny nawigacji, skrót i polecenie języka naturalnego |

### 2.5. Zachowanie przy zwężaniu kolumn

| Sytuacja | Zachowanie |
|---|---|
| Zwężenie kolumn obszaru roboczego | Awatar zachowuje wymiar i pozycję względem prawej krawędzi obszaru roboczego; nie skaluje się wraz z kolumnami |
| Powierzchnia interakcji otwarta, kolumny zwężone poniżej progu czytelności | Powierzchnia interakcji zwija się do dymka kontekstowego; strumień rozmowy pozostaje dostępny po ponownym rozszerzeniu kolumny |
| Kolizja awatara z panelem pomocniczym otwartym po prawej | Awatar przesuwa się na krawędź nowej, najbardziej wysuniętej kolumny bocznej i pozostaje w całości widoczny |
| Dymek sugestii przy wąskiej kolumnie obszaru roboczego | Dymek otwiera się po lewej stronie awatara i skraca treść do jednego zdania z działaniem „Rozwiń” otwierającym powierzchnię interakcji |

### 2.6. Katalog elementów funkcji

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Gdzie występuje | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Awatar Always On Display | Pływający awatar agenta towarzyszącego, `.dn-awatar` w pozycjonowaniu stałym | Stały punkt dostępu do funkcji globalnej, niezależny od środowiska i modułu | Koło 48 px, waga znikoma | spoczynek · sugestia oczekująca · sugestia o wadze wysokiej · mówi/pisze · nasłuch · wyciszony · ukryty | Kliknięcie otwiera dymek kontekstowy z bieżącą sugestią; kliknięcie podwójne otwiera powierzchnię interakcji | Nad całą powłoką aplikacji, przy prawej krawędzi obszaru roboczego | 1 | Widoczny bez interakcji |
| Plakietka powiadomień | Plakietka liczbowa, `.dn-plakietka` + `.dn-kropka` | Sygnalizacja liczby oczekujących sugestii | Bardzo mała plakietka na awatarze | ukryta (zero) · z liczbą · pulsująca (waga wysoka) | Kliknięcie otwiera listę oczekujących sugestii w powierzchni interakcji | Na awatarze | 1 | Widoczna bez interakcji |
| Dymek kontekstowy sugestii | Dymek kontekstowy z treścią sugestii i działaniami, `ao-dymek` (selektor lokalny prototypu `design/05-okna/platformowe/always-on-display.html`) | Przekazanie sugestii bez przerywania pracy | Dymek średniej wielkości, waga niska | ukryty · widoczny · z działaniami · zwinięty do skrótu (wąska kolumna) | Wybór działania wykonuje skutek z rozdz. 4 i zamyka dymek; brak wyboru zamyka dymek po odejściu fokusu | Przy awatarze | 2 | Kliknięcie awatara; ujawnia się samoczynnie przy sugestii o wadze wysokiej; po użyciu zwija się samoczynnie |
| Powierzchnia interakcji | Kolumna boczna ze strumieniem rozmowy, listą sugestii i polem polecenia — komponent bez odrębnej klasy `.dn-*` w bieżącym stanie `komponenty.css`; **[DO DECYZJI OPERATORA]** doprecyzowania klasy dla kolumny bocznej rozszerzalnej funkcji globalnej | Rozmowa z agentem towarzyszącym i przegląd oczekujących sugestii | Kolumna boczna, otwierana jako rozszerzenie boczne, szerokość regulowana | zamknięta · otwarta · strumieniowanie odpowiedzi · nasłuch głosowy · zwinięta do dymka | Zamknięcie przywraca stan bezczynny bez utraty historii rozmowy | Po prawej stronie obszaru roboczego | 2 | Kliknięcie podwójne awatara; pozycja „Always On Display” listwy ustawień strefy 3; szybki wybór szyny nawigacji; skrót `Ctrl/Cmd + Shift + A` |
| Lista oczekujących sugestii | Lista pozycji, `.dn-karta--klikalna` | Przegląd sugestii, które nie zostały jeszcze przyjęte ani odrzucone | Lista w powierzchni interakcji, waga średnia | pusta · wypełniona · pozycja rozwinięta | Rozwinięcie pozycji odsłania pełną treść i komplet działań sugestii | Powierzchnia interakcji | 2 | Kliknięcie plakietki powiadomień; pozycja „Sugestie” powierzchni interakcji |
| Przycisk mikrofonu | Przycisk ikonowy toru głosowego, `.dn-btn-ikona` | Wydanie krótkiego polecenia głosowego skierowanego do platformy jako całości | Mała ikona w nagłówku powierzchni interakcji | domyślny · nasłuch · przetwarzanie · błąd urządzenia wejściowego | Aktywacja rozpoczyna nasłuch; zakończenie wysyła `aod.voice.command` | Nagłówek powierzchni interakcji; dymek kontekstowy | 2 | Ikona `mikrofon`; skrót `Ctrl/Cmd + Shift + V`; fraza wybudzająca |
| Pole polecenia do Always On Display | Pole tekstowe, `.dn-pole-kontrolka` | Wydanie polecenia tekstowego poza kontekstem karty sesji | Pole formularza w stopce powierzchni interakcji | domyślny · fokus · wysyłanie · błąd kanału | Wysyła `aod.chat.send`; odpowiedź trafia do strumienia powierzchni interakcji | Stopka powierzchni interakcji | 2 | Widoczne po otwarciu powierzchni interakcji |
| Przełącznik trybu obecności | Kontrolka trzystanowa, `.dn-zakladki--pigulki` | Wybór trybu obecności funkcji: pełny, cichy, ukryty (rozdz. 9.2) | Mały przełącznik segmentowy | pełny · cichy · ukryty | Zmiana obowiązuje natychmiast na wszystkich urządzeniach Operatora | Nagłówek powierzchni interakcji; sekcja „Always On Display” okna Ustawień | 2 | Znacznik `[ tryb ▼ ]` w nagłówku powierzchni interakcji |
| Wyciszenie sugestii | Działanie czasowe w menu kebab | Wstrzymanie ujawniania sugestii na wskazany czas lub w wskazanym zakresie | Pozycja menu z listą wartości czasu | brak wyciszenia · wyciszenie czasowe · wyciszenie zakresu | Uruchamia regułę wyciszenia z rozdz. 3.5; awatar przechodzi w stan wyciszony | Menu kebab powierzchni interakcji | 3 | Menu kebab (⋮) powierzchni interakcji |
| Konfiguracja reguł wyzwalania | Zestaw ustawień progów, częstotliwości i klas zdarzeń | Zmiana warunków, w jakich funkcja przedstawia sugestie | Panel ustawień w oknie Konfiguracji, sekcja „Always On Display” okna Ustawień | domyślna · zmieniona przez Operatora · profil roli | Zapis obowiązuje natychmiast i rozgłasza `config.changed` | Okno konfiguracji; okno Ustawień | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, skrót klawiszowy, tryb administracyjny |
| Przełącznik trybu wobec procesu | Kontrolka dwustanowa obserwator/operator | Wybór zakresu działania wobec procesu środowiska MultitaskingAI | Przełącznik segmentowy, `.dn-suwak` | obserwator · operator | Rozszerza lub zawęża komplet działań interwencyjnych | Sekcja Monitor procesu panelu orkiestracji środowiska MultitaskingAI | 2 | Znacznik `AOD ▼` w panelu stanu procesu ([Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2.6) |

---

## 3. Reguły wyzwalania proaktywnych sugestii

Rozdział obowiązuje we wszystkich środowiskach i modułach platformy. Środowisko MultitaskingAI stosuje te same reguły, rozszerzone o zdarzenia zespołu ról i silnika kolejek opisane w opracowaniu [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) (rozdz. 2).

### 3.1. Zasada wyzwalania

Always On Display nie komentuje pracy w sposób ciągły. Sugestia powstaje wyłącznie wtedy, gdy zdarzenie należy do jednej z klas wyzwalających z rozdz. 3.2, przekracza próg z rozdz. 3.4 i nie podlega wyciszeniu z rozdz. 3.5. Zdarzenie niespełniające tych warunków zostaje odnotowane w kontekście funkcji i nie ujawnia się w interfejsie.

### 3.2. Klasy zdarzeń wyzwalających

| Klasa | Zdarzenie wyzwalające | Rodzaj powstającej sugestii (rozdz. 4) |
|---|---|---|
| Stan pętli wykonawczej | Pętla wstrzymana, pętla przerwana, zadanie ponawiane powyżej progu, zlecenie oczekujące na zatwierdzenie Użytkownika | Punkt decyzyjny, wskazanie problemu |
| Stan kolejki zadań | Kolejka zatrzymana, zadanie w stanie błędu, zadanie oczekujące dłużej niż próg czasu, kolejka wypełniona powyżej progu | Wskazanie problemu, kolejny krok |
| Wynik kontroli jakości | Negatywny wynik kontroli, powtórzone niepowodzenie tego samego zadania, rozbieżność wyniku z zleceniem | Wskazanie problemu, kolejny krok |
| Zdarzenia modułów | Zakończenie długiego zadania modułu, dostępny wynik do przeglądu, konflikt zasobu, brak konfiguracji potrzebnej do wykonania polecenia | Doradztwo, sugestia konfiguracji |
| Harmonogram | Nadejście terminu reguły czasowej, cykliczne uruchomienie automatyki, przypomnienie ustanowione przez Operatora | Kolejny krok, doradztwo |
| Kontekst pracy Operatora | Powtarzalność ręcznie wykonywanej czynności, praca w module bez skonfigurowanego wykonawcy, otwarty zasób powiązany z zadaniem oczekującym | Doradztwo, sugestia konfiguracji |

### 3.3. Źródła sygnału

| Źródło | Nośnik | Zakres odczytu |
|---|---|---|
| Execution Loop Window | Zdarzenia pętli wykonawczej kanału Koordynator ↔ Wykonawca | Stan pętli, dekompozycja zlecenia, punkty decyzyjne, wyniki kontroli jakości, decyzje o ponowieniu |
| Silnik kolejek | Zdarzenia zadań i kolejek | Stan zadania, czas oczekiwania, liczba ponowień, wypełnienie kolejki |
| Moduły platformy | Zdarzenia okien operacyjnych modułów | Zakończenie zadania, dostępny wynik, konflikt zasobu, brak konfiguracji |
| Harmonogram i automatyki | Reguły czasowe modułu Automations | Terminy, cykle, uruchomienia i ich wynik |
| Chat Window | Strumień kanału Użytkownik ↔ Wykonawca | Treść zlecenia, zatwierdzenia, przerwania — jako kontekst rozumienia sytuacji |
| Historia Operatora | Zapis działań w sesjach i projektach | Powtarzalność czynności, wcześniejsze decyzje wobec analogicznych sugestii |

### 3.4. Progi i częstotliwość

| Parametr | Wartość domyślna | Skutek |
|---|---|---|
| Waga sugestii ujawnianej samoczynnie | wysoka | Sugestia o wadze wysokiej otwiera dymek kontekstowy bez interakcji Operatora; pozostałe czekają pod plakietką |
| Liczba sugestii ujawnianych samoczynnie w godzinie | 3 | Po przekroczeniu limitu kolejne sugestie trafiają do listy oczekujących bez otwierania dymka |
| Odstęp między dymkami | 5 minut | Dwie sugestie oddzielone krótszym odstępem łączą się w jedną pozycję zbiorczą |
| Próg czasu oczekiwania zadania w kolejce | 15 minut | Przekroczenie tworzy sugestię klasy „stan kolejki zadań” |
| Próg liczby ponowień zadania | 2 | Trzecie ponowienie tego samego zadania tworzy sugestię o wadze wysokiej |
| Próg wypełnienia kolejki | 80% pojemności | Przekroczenie tworzy sugestię klasy „stan kolejki zadań” |
| Próg powtarzalności czynności ręcznej | 3 wystąpienia w karcie sesji | Przekroczenie tworzy sugestię konfiguracji (automatyka, makro, agent) |
| Czas życia sugestii nieprzyjętej | 24 godziny | Po upływie sugestia otrzymuje status odrzuconej i znika z listy oczekujących |

Wszystkie progi są ustawieniami konfiguracyjnymi zakresu „Always On Display” (rozdz. 10) i podlegają zmianie z okna konfiguracji, z okna Ustawień oraz poleceniem języka naturalnego w Chat Window.

### 3.5. Wyciszanie

| Rodzaj wyciszenia | Zakres | Zachowanie funkcji |
|---|---|---|
| Wyciszenie czasowe | 15 minut, 1 godzina, do końca dnia | Sugestie gromadzą się w liście oczekujących; dymek nie otwiera się; plakietka pozostaje ukryta |
| Wyciszenie kontekstowe | Bieżący moduł albo bieżąca karta sesji | Sugestie dotyczące wskazanego kontekstu nie ujawniają się; pozostałe zachowują pełne działanie |
| Wyciszenie klasy zdarzeń | Wskazana klasa z rozdz. 3.2 | Klasa nie tworzy sugestii do chwili przywrócenia; pozostałe klasy działają bez zmian |
| Tryb cichy | Cała platforma | Funkcja nie ujawnia sugestii samoczynnie i nie stosuje syntezy mowy; lista oczekujących pozostaje dostępna po otwarciu powierzchni interakcji |
| Wyjątek wagi krytycznej | Punkt decyzyjny pętli wykonawczej wstrzymujący proces | Ujawnia się mimo wyciszenia czasowego i kontekstowego, w postaci plakietki bez dymka; tryb cichy zachowuje ten wyjątek bez syntezy mowy |

Wyciszenie jest odwracalne w każdej chwili — jednym kliknięciem w menu kebab powierzchni interakcji, jednym skrótem klawiszowym albo poleceniem języka naturalnego w Chat Window.

### 3.6. Zachowanie przy braku istotnych zdarzeń

Przy braku zdarzeń spełniających warunki wyzwalania funkcja pozostaje w stanie bezczynnym: awatar w barwie bazowej, plakietka ukryta, dymek zamknięty, powierzchnia interakcji zwinięta. Funkcja nie wypełnia ciszy treścią zastępczą, nie przypomina o swojej obecności i nie podsumowuje pracy bez wywołania. Otwarcie powierzchni interakcji w tym stanie prezentuje strumień rozmowy z pustym stanem listy sugestii oraz gotowe pole polecenia.

---

## 4. Katalog rodzajów sugestii

Rodzaj sugestii odpowiada wprost polu `rodzaj` encji `sugestia_aod` ([Model danych](../architektura/model-danych.md), rozdz. 18.2): `doradztwo`, `konfiguracja`, `problem`, `kolejny_krok`. Każda sugestia niesie treść, komplet działań i skutek każdego działania.

### 4.1. Doradztwo (`doradztwo`)

| Cecha | Treść |
|---|---|
| Kiedy powstaje | Zakończenie długiego zadania modułu, dostępny wynik do przeglądu, praca nad zasobem powiązanym z zadaniem oczekującym |
| Treść sugestii | Zdanie opisujące zaobserwowany stan i zalecane działanie, na przykład: „Zadanie badawcze w module Research zakończyło się — wynik oczekuje na przegląd.” |

| Działanie | Skutek |
|---|---|
| Przejdź do wyniku | Otwiera moduł i okno operacyjne, w którym wynik powstał; karta sesji pozostaje ta sama |
| Przenieś do Chat Window | Wstawia treść sugestii do pola polecenia Chat Window jako polecenie gotowe do wysłania (rozdz. 7.2) |
| Odłóż | Pozycja wraca do listy oczekujących, dymek się zamyka, status pozostaje `nowa` |
| Odrzuć | Status `odrzucona`; analogiczne zdarzenie nie tworzy sugestii do końca karty sesji |

### 4.2. Sugestia konfiguracji (`konfiguracja`)

| Cecha | Treść |
|---|---|
| Kiedy powstaje | Powtarzalność czynności ręcznej powyżej progu, brak konfiguracji potrzebnej do wykonania polecenia, praca w module bez przypisanego wykonawcy |
| Treść sugestii | Zdanie wskazujące ustawienie, automatykę albo agenta rozwiązującego powtarzalność, na przykład: „Ta sama sekwencja poleceń powtórzyła się trzykrotnie — utworzenie automatyki wykona ją jednym wywołaniem.” |

| Działanie | Skutek |
|---|---|
| Otwórz ustawienie | Otwiera okno konfiguracji na wskazanej pozycji, z objaśnieniem kontekstowym `[?]` |
| Utwórz automatykę | Otwiera moduł Automations z formularzem wypełnionym rozpoznaną sekwencją |
| Przypisz agenta | Otwiera wybór agenta z modułu Agents dla bieżącego kontekstu pracy |
| Odrzuć | Status `odrzucona`; funkcja nie powtarza tej sugestii dla tego samego kontekstu |

### 4.3. Wskazanie problemu (`problem`)

| Cecha | Treść |
|---|---|
| Kiedy powstaje | Negatywny wynik kontroli jakości, zadanie w stanie błędu, powtórzone ponowienie tego samego zadania, kolejka zatrzymana, pętla przerwana |
| Treść sugestii | Zdanie nazywające problem i jego umiejscowienie, na przykład: „Kontrola jakości odrzuciła wynik zadania po raz drugi — przyczyną jest niezgodność wyniku ze zleceniem.” |

| Działanie | Skutek |
|---|---|
| Otwórz Execution Loop Window | Otwiera okno pętli wykonawczej na zadaniu, którego dotyczy problem (rozdz. 7.3) |
| Ponów zadanie | Wysyła do Koordynatora decyzję o ponowieniu wskazanego zadania |
| Skoryguj zlecenie | Wstawia treść zlecenia do pola polecenia Chat Window w trybie edycji, do poprawy i ponownego wysłania |
| Wstrzymaj proces | Wywołuje wstrzymanie kolejki albo procesu; działanie dostępne w trybie operatora, w trybie obserwatora wyświetla komunikat kontekstowy ze wskazaniem przełączenia trybu |
| Odrzuć | Status `odrzucona`; problem pozostaje widoczny w Execution Loop Window, sugestia nie wraca |

### 4.4. Kolejny krok (`kolejny_krok`)

| Cecha | Treść |
|---|---|
| Kiedy powstaje | Zlecenie oczekujące na zatwierdzenie, nadejście terminu reguły harmonogramu, zakończony etap pracy z jednoznacznym następstwem |
| Treść sugestii | Zdanie nazywające krok następny, na przykład: „Zlecenie oczekuje na zatwierdzenie kroku — po zatwierdzeniu pętla przejdzie do kontroli jakości.” |

| Działanie | Skutek |
|---|---|
| Zatwierdź krok | Przekazuje zatwierdzenie punktu decyzyjnego; pętla przechodzi do kroku następnego. W trybie obserwatora przycisk pozostaje klikalny i wyświetla komunikat kontekstowy ze wskazaniem przełączenia w tryb operatora |
| Wznów proces | Wywołuje wznowienie wstrzymanej kolejki albo procesu |
| Przenieś do Chat Window | Wstawia krok jako polecenie do pola polecenia Chat Window |
| Odłóż | Pozycja wraca do listy oczekujących; punkt decyzyjny pozostaje otwarty w Execution Loop Window |

### 4.5. Waga sugestii i sposób ujawnienia

| Waga | Przykład | Sposób ujawnienia |
|---|---|---|
| Wysoka | Punkt decyzyjny wstrzymujący proces, trzecie ponowienie zadania, kolejka zatrzymana | Dymek kontekstowy otwierany samoczynnie, plakietka pulsująca |
| Średnia | Zakończone zadanie z wynikiem do przeglądu, przekroczony czas oczekiwania zadania | Plakietka liczbowa; dymek po kliknięciu awatara |
| Niska | Sugestia konfiguracji wynikająca z powtarzalności, przypomnienie harmonogramu | Pozycja listy oczekujących; brak plakietki pulsującej |

---

## 5. Tor głosowy i granica wobec modułu Assistant

Platforma prowadzi dwa tory głosowe o rozłącznym zakresie: tor globalny Always On Display i tor sesyjny okna Voice Console modułu Assistant ([Moduł Assistant](../moduly/assistant.md), rozdz. 3.3). Rozdział rozstrzyga podział bez pozostawiania obszaru wspólnego.

### 5.1. Podział zakresów

| Cecha | Tor Always On Display | Tor Voice Console (moduł Assistant) |
|---|---|---|
| Umiejscowienie | Warstwa funkcji globalnych, ponad środowiskami i modułami | Okno operacyjne modułu Assistant, prawa kolumna obszaru roboczego |
| Kontekst | Platforma jako całość; brak przypisania do karty sesji | Karta sesji modułu Assistant, profil asystenta, pamięć modułu |
| Długość interakcji | Pojedyncze, krótkie polecenie i krótka odpowiedź | Rozmowa wieloturowa, dyktowanie, łańcuch poleceń w jednej wypowiedzi |
| Aktywacja | Kliknięcie awatara, skrót `Ctrl/Cmd + Shift + V`, fraza wybudzająca globalna | Duży przycisk mikrofonu, przytrzymanie, fraza wybudzająca profilu, tryb ciągłego nasłuchu |
| Rozpoznawanie mowy | Rozpoznawanie krótkiej frazy, bez trybu ciągłego nasłuchu | Pełny tor STT: detekcja mowy, redukcja szumu, wskaźnik pewności, automatyczne wykrycie języka |
| Synteza mowy | Krótkie potwierdzenie i treść sugestii; barwa i tempo wspólne dla całej platformy | Pełny tor TTS: wybór głosu, tempa i tonu wypowiedzi profilu asystenta |
| Dyktowanie | Poza zakresem | Dyktowanie tekstu do pola aktywnego wraz z komendami formatującymi |
| Makra i polecenia szybkie | Poza zakresem | Biblioteka poleceń szybkich, siatka szybkich akcji, makra głosowe |
| Zapis | `polecenie_glosowe` z pustym `profil_asystenta_id` | `polecenie_glosowe` z wypełnionym `profil_asystenta_id` |
| Trwałość | Rozmowa poza kartą sesji, bez wpływu na pamięć modułu | Rozmowa w karcie sesji, zapisywana w pamięci i historii modułu |

### 5.2. Która powierzchnia jest właściwa dla którego polecenia

| Rodzaj polecenia głosowego | Powierzchnia właściwa | Uzasadnienie |
|---|---|---|
| „Pokaż stan pętli”, „ile zadań w kolejce” | Always On Display | Odczyt stanu platformy niezależny od modułu |
| „Zatwierdź krok”, „wstrzymaj proces”, „wznów kolejkę” | Always On Display | Sterowanie procesem z dowolnego miejsca platformy |
| „Otwórz moduł Studio”, „przełącz na CodeStudio” | Always On Display | Nawigacja między środowiskami i modułami |
| „Co się wydarzyło w projekcie od rana” | Always On Display | Zapytanie o kontekst obejmujący wiele sesji i modułów |
| „Podyktuj akapit do dokumentu” | Voice Console modułu Assistant | Dyktowanie należy do toru sesyjnego |
| „Sprawdź pocztę i podsumuj najważniejsze wiadomości” | Voice Console modułu Assistant | Polecenie wieloetapowe prowadzone w pętli wykonawczej sesji modułu |
| „Zmień głos asystenta na spokojniejszy” | Voice Console modułu Assistant | Ustawienie profilu asystenta |
| „Uruchom makro «raport dzienny»” | Voice Console modułu Assistant | Makra należą do biblioteki modułu |
| Polecenie wydane, gdy moduł Assistant jest oknem wiodącym | Voice Console modułu Assistant | Tor sesyjny ma pierwszeństwo w obrębie własnego modułu (rozdz. 5.4) |

### 5.3. Przekazanie między torami

Przekazanie następuje w jednym kierunku dla poleceń wykraczających poza zakres toru globalnego i w drugim dla poleceń wykraczających poza sesję modułu.

```
Schemat 1 — Przekazanie polecenia głosowego między torami

  Operator mówi
        │
        ▼
  ┌──────────────────────────────┐        polecenie krótkie, ponadkontekstowe
  │ Always On Display            │ ─────► wykonanie w torze globalnym
  │ (tor globalny)               │        (aod.voice.command)
  └──────────────────────────────┘
        │  polecenie wymagające dyktowania, rozmowy wieloturowej,
        │  makra, profilu głosu albo pamięci modułu
        ▼
  przekazanie: otwarcie modułu Assistant, uruchomienie okna Voice Console,
  przeniesienie rozpoznanej treści jako pierwszej wypowiedzi rozmowy
        │
        ▼
  ┌──────────────────────────────┐        rozmowa prowadzona dalej w torze
  │ Voice Console (tor sesyjny)  │ ─────► sesyjnym; pamięć i historia modułu
  └──────────────────────────────┘
        │  polecenie dotyczące platformy jako całości
        │  (nawigacja, sterowanie procesem poza sesją)
        ▼
  przekazanie zwrotne: wykonanie w torze globalnym, bez zamykania Voice Console
```

| Cecha przekazania | Zasada |
|---|---|
| Zachowanie treści | Rozpoznana treść polecenia przechodzi w całości; Operator nie powtarza wypowiedzi |
| Zachowanie nasłuchu | Nasłuch toru źródłowego kończy się z chwilą przekazania; tor docelowy przejmuje nasłuch bez ponownej aktywacji |
| Kontekst | Przekazanie nie łączy kontekstów: tor globalny nie zapisuje rozmowy w pamięci modułu, tor sesyjny nie przejmuje kolejki sugestii |
| Widoczność | Przekazanie jest jawne — powierzchnia interakcji funkcji odnotowuje przekazanie i wskazuje okno docelowe |
| Odwracalność | Zamknięcie okna Voice Console nie przerywa działania funkcji globalnej; awatar wraca do stanu bezczynnego |

### 5.4. Rozstrzygnięcie pierwszeństwa

| Sytuacja | Tor obsługujący |
|---|---|
| Moduł Assistant jest oknem wiodącym, Operator używa przycisku mikrofonu Voice Console | Voice Console |
| Moduł Assistant jest oknem wiodącym, Operator klika mikrofon w powierzchni interakcji funkcji globalnej | Always On Display; polecenie wykraczające poza zakres globalny przechodzi do Voice Console (rozdz. 5.3) |
| Operator pracuje w dowolnym innym module i wypowiada globalną frazę wybudzającą | Always On Display |
| Fraza wybudzająca profilu asystenta wypowiedziana poza modułem Assistant | Always On Display otwiera moduł Assistant i przekazuje nasłuch do Voice Console |
| Trwa synteza mowy toru globalnego, Operator uruchamia Voice Console | Synteza toru globalnego milknie; nasłuch przejmuje Voice Console |

Poza wymienionymi sytuacjami tory nie działają jednocześnie: w danej chwili nasłuch prowadzi dokładnie jedna powierzchnia.

---

## 6. Zachowanie per środowisko i per moduł

### 6.1. Zachowanie podstawowe

Zachowaniem podstawowym jest zestaw z rozdziałów 2–5: awatar przy prawej krawędzi obszaru roboczego, powierzchnia interakcji jako rozszerzenie boczne, komplet klas zdarzeń wyzwalających i progów domyślnych, tor głosowy o zakresie globalnym. Tabela poniżej podaje wyłącznie różnice względem tego zachowania.

### 6.2. Tabela zachowań

| Środowisko / moduł | Umiejscowienie | Rodzaje sugestii właściwe dla kontekstu | Różnice względem zachowania podstawowego |
|---|---|---|---|
| Strona główna | Awatar przy prawej krawędzi strefy roboczej strony; pozycja „Always On Display” listwy ustawień strefy 3 | Doradztwo, konfiguracja | Sugestie dotyczą wyboru środowiska i komponentów własnych; brak sugestii klasy pętli i kolejek |
| Środowisko TalkIn | Zachowanie podstawowe | Doradztwo, kolejny krok | Sugestie odnoszą się do treści krążącej między modułami — odłożonych zasobów, wyników badań, materiałów do przeglądu |
| Środowisko WorkSpace | Zachowanie podstawowe | Doradztwo, konfiguracja, kolejny krok | Sugestie odnoszą się do projektu i jego zasobów; próg powtarzalności czynności ręcznej odnosi się do projektu, nie do karty sesji |
| Środowisko CodeStudio | Zachowanie podstawowe | Problem, kolejny krok | Sugestie odnoszą się do wyniku kontroli jakości, przebiegu zadań budowy i stanu repozytorium; waga sugestii klasy „wynik kontroli jakości” podniesiona do wysokiej |
| Środowisko MultitaskingAI | Awatar przy prawej krawędzi obszaru roboczego; dodatkowo wskaźnik roli w sekcji Monitor procesu panelu orkiestracji | Problem, kolejny krok, doradztwo, konfiguracja | Przypadek szczególny — funkcja pełni rolę warstwy centralnej z trybami obserwatora i operatora; komplet oprzyrządowania, ról i kolejek opisuje [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2 |
| Moduł Assistant | Zachowanie podstawowe | Doradztwo, kolejny krok | Tor głosowy ustępuje pierwszeństwa oknu Voice Console (rozdz. 5.4); sugestie nie dotyczą treści rozmowy prowadzonej w module |
| Moduł Automations | Zachowanie podstawowe | Problem, kolejny krok, konfiguracja | Klasa „harmonogram” jest podstawowym źródłem sygnału; sugestie odnoszą się do przebiegów automatyk i ich wyników |
| Moduł Agents | Zachowanie podstawowe | Konfiguracja, doradztwo | Sugestie dotyczą przypisania agenta do roli i zakresu uprawnień; klasa „stan kolejki zadań” nie tworzy sugestii |
| Moduł Workspace | Zachowanie podstawowe | Doradztwo, konfiguracja | Sugestie odnoszą się do izolacji projektu i jego zależności |
| Moduł Developer | Zachowanie podstawowe | Problem, kolejny krok | Waga sugestii wynikających z negatywnej kontroli jakości podniesiona do wysokiej |
| Moduł Terminal | Awatar przy prawej krawędzi obszaru roboczego; dymek zwinięty do skrótu | Problem | Sugestie ograniczone do przerwanych i zakończonych błędem przebiegów; dymek nie otwiera się samoczynnie w trakcie strumieniowania wyjścia |
| Moduł Studio | Zachowanie podstawowe | Doradztwo, kolejny krok | Sugestie odnoszą się do wersji dokumentu i zadań redakcyjnych oczekujących na przegląd |
| Moduł Research | Zachowanie podstawowe | Doradztwo, kolejny krok | Sugestie odnoszą się do zakończonych zadań badawczych i wyników do przeglądu |
| Moduł Library | Zachowanie podstawowe | Doradztwo, konfiguracja | Sugestie odnoszą się do katalogowania artefaktów powstałych w innych modułach |
| Moduł Roundtable | Zachowanie podstawowe | Kolejny krok | Sugestie odnoszą się do zakończonej debaty i decyzji oczekującej na Operatora; w trakcie debaty funkcja nie ujawnia dymka |
| Moduł Diagnostics | Zachowanie podstawowe | Problem | Wszystkie sugestie klasy „wynik kontroli jakości” otrzymują wagę wysoką |
| Moduł Browser | Zachowanie podstawowe | Doradztwo | Sugestie odnoszą się do zakończonych zadań pozyskiwania treści |
| Moduł Design | Zachowanie podstawowe | Doradztwo, kolejny krok | Sugestie odnoszą się do wariantów projektu oczekujących na wybór |
| Moduł Translate | Zachowanie podstawowe | Doradztwo | Sugestie odnoszą się do zakończonych zadań tłumaczenia |
| Moduł Apps | Zachowanie podstawowe | Problem, kolejny krok | Sugestie odnoszą się do przebiegu budowy aplikacji i jej wyniku |
| Okno konfiguracji i okno Ustawień | Awatar przy prawej krawędzi obszaru roboczego | Konfiguracja | Sugestie ograniczone do wyjaśnienia skutku zmienianego ustawienia; klasy zdarzeń procesowych nie tworzą dymka w trakcie edycji ustawień |

### 6.3. MultitaskingAI jako przypadek szczególny

W środowisku MultitaskingAI Always On Display zachowuje wszystkie właściwości opisane w niniejszym dokumencie i uzyskuje dodatkowo rolę warstwy centralnej nadrzędnej wobec zespołu ról. Dwa tryby względem procesu — obserwator i operator — pełne oprzyrządowanie nadzoru, pokrycie sekcji Monitor procesu, współdziałanie z funkcją Mobile w interwencji zdalnej oraz katalog elementów interfejsu funkcji w tym środowisku opisuje [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2. Reguły wyzwalania sugestii, katalog rodzajów sugestii, tor głosowy i warstwy widoczności obowiązują tam w brzmieniu z niniejszego dokumentu.

---

## 7. Relacja do dwóch kanałów komunikacji operacyjnej

### 7.1. Miejsce funkcji wobec obu kanałów

Always On Display nie jest trzecim kanałem komunikacji operacyjnej i nie zastępuje żadnego z dwóch istniejących. Czyta sygnał z obu kanałów i kieruje do nich skutek: Chat Window pozostaje miejscem wydania polecenia i podstawowym mechanizmem sterowania procesami platformy, Execution Loop Window — miejscem prowadzenia i nadzoru pętli wykonawczej.

```
Schemat 2 — Przepływ sygnału i skutku między funkcją a kanałami komunikacji

  Execution Loop Window                              Chat Window
  Koordynator ↔ Wykonawca                            Użytkownik ↔ Wykonawca
        │                                                   ▲
        │  sygnał: stan pętli, punkt decyzyjny,             │  skutek: sugestia
        │  wynik kontroli jakości, decyzja o ponowieniu     │  przeniesiona jako
        ▼                                                   │  polecenie
  ┌───────────────────────────────────────────────────────────────┐
  │  Always On Display  — warstwa funkcji globalnych              │
  │  wyzwalanie sugestii · lista oczekujących · tor głosowy       │
  └───────────────────────────────────────────────────────────────┘
        │
        │  skutek bezpośredni w trybie operatora:
        ▼  zatwierdzenie kroku · wstrzymanie · wznowienie
  Execution Loop Window
```

### 7.2. Sugestia przeniesiona do Chat Window jako polecenie

| Cecha | Zasada |
|---|---|
| Wyzwolenie | Działanie „Przenieś do Chat Window” dostępne przy każdym rodzaju sugestii (rozdz. 4) |
| Postać | Treść sugestii przekształcona w polecenie języka naturalnego wstawione do pola polecenia Chat Window bieżącej karty sesji |
| Stan przed wysłaniem | Polecenie pozostaje edytowalne; Operator poprawia je przed wysłaniem albo wysyła bez zmian |
| Wykonanie | Wysłanie polecenia uruchamia zadanie w kanale Użytkownik ↔ Wykonawca na zasadach ogólnych; wynik pozostaje w strumieniu rozmowy |
| Status sugestii | Wysłanie polecenia ustawia status sugestii na `przyjeta`; pozycja znika z listy oczekujących |
| Rozgraniczenie | Funkcja nie wykonuje polecenia zamiast Chat Window i nie tworzy równoległej ścieżki wykonania — przenosi treść i oddaje sterowanie kanałowi pierwszemu |

### 7.3. Sygnały pochodzące z Execution Loop Window

| Sygnał | Co niesie | Reakcja funkcji |
|---|---|---|
| Stan pętli | Pętla bezczynna, w toku, wstrzymana, ponowienie zadania, przerwana, błąd | Zmiana stanu awatara; sugestia klasy „stan pętli wykonawczej” przy przejściu do stanu wstrzymanej, przerwanej albo błędu |
| Punkt decyzyjny | Krok zlecenia oczekujący na zatwierdzenie | Sugestia rodzaju `kolejny_krok` o wadze wysokiej; działanie „Zatwierdź krok” przekazuje decyzję do pętli |
| Wynik kontroli jakości | Wynik pozytywny albo negatywny wraz z przyczyną | Wynik negatywny tworzy sugestię rodzaju `problem`; wynik pozytywny nie tworzy sugestii |
| Decyzja o ponowieniu | Liczba ponowień zadania | Przekroczenie progu ponowień tworzy sugestię rodzaju `problem` o wadze wysokiej |
| Dekompozycja zlecenia | Podział zlecenia na zadania i ich kolejka | Kontekst rozumienia sytuacji; samodzielnie nie tworzy sugestii |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Działania sugestii wywołują te same operacje sterujące, którymi dysponuje okno pętli |

### 7.4. Rozgraniczenie odpowiedzialności

| Pytanie | Rozstrzygnięcie |
|---|---|
| Gdzie powstaje polecenie | W Chat Window — także wtedy, gdy jego treść pochodzi z sugestii |
| Gdzie prowadzona jest pętla | W Execution Loop Window — funkcja nie prowadzi własnej pętli |
| Gdzie zapada decyzja o zatwierdzeniu kroku | W punkcie decyzyjnym pętli; funkcja przekazuje decyzję Operatora, nie podejmuje jej samodzielnie |
| Co wnosi funkcja | Obecność ponad kontekstem, wykrycie zdarzenia istotnego, sformułowanie sugestii i skrócenie drogi do właściwego okna |

---

## 8. Warstwy widoczności

### 8.1. Przydział warstw

Always On Display stosuje regułę stopniowego ujawniania funkcjonalności obowiązującą w całej platformie: jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna. Funkcja w postaci minimalnej należy do warstwy 1, rozwinięcia sugestii do warstw 2–3, konfiguracja reguł wyzwalania do warstwy 4.

| Warstwa | Zawartość funkcji | Sposób wywołania |
|---|---|---|
| 1 — zawsze widoczna | Awatar w postaci minimalnej (koło 48 px) oraz plakietka liczby oczekujących sugestii | Widoczne bez interakcji |
| 2 — widoczna na żądanie | Dymek kontekstowy sugestii, powierzchnia interakcji, lista oczekujących sugestii, pole polecenia, przycisk mikrofonu, przełącznik trybu obecności, przełącznik trybu wobec procesu | Kliknięcie awatara, kliknięcie plakietki, szybki wybór szyny nawigacji, pozycja listwy strefy 3, skrót klawiszowy; po użyciu element zwija się samoczynnie |
| 3 — rozwinięcia kontekstowe | Komplet działań sugestii, wyciszanie czasowe i kontekstowe, wybór zakresu wyciszenia, menu akcji pozycji listy sugestii, historia sugestii przyjętych i odrzuconych | Menu kebab (⋮) powierzchni interakcji, menu kontekstowe pozycji listy, panel popover |
| 4 — funkcje eksperckie | Konfiguracja reguł wyzwalania: klasy zdarzeń, progi, częstotliwość, wagi sugestii, źródła sygnału, reguły wyciszania klas, profil warstw według roli użytkownika | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli; użytkownik podstawowy nie widzi tych elementów |

### 8.2. Zasada jednego kliknięcia

Każdy ukryty element funkcji jest osiągalny jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window. Powierzchnia interakcji otwiera się jednym kliknięciem podwójnym awatara, jednym skrótem `Ctrl/Cmd + Shift + A` albo jedną pozycją listwy ustawień strefy 3 strony głównej. Konfiguracja reguł wyzwalania, mimo przynależności do warstwy 4, pozostaje osiągalna jednym poleceniem języka naturalnego.

### 8.3. Mechanizmy ukrywania zastosowane w funkcji

| Mechanizm | Zastosowanie |
|---|---|
| Menu progresywne | Przełącznik trybu obecności `[ tryb ▼ ]` w nagłówku powierzchni interakcji |
| Panele wysuwane | Powierzchnia interakcji jako rozszerzenie boczne, dymek kontekstowy jako popover — po zamknięciu znikają całkowicie z przestrzeni roboczej |
| Grupowanie logiczne akcji | Menu kebab (⋮) powierzchni interakcji zbierające wyciszanie, historię sugestii i wejście do konfiguracji reguł |
| Znaczniki kontekstowe | Znacznik `AOD ▼` w panelu stanu procesu środowiska MultitaskingAI otwierający przełącznik trybu wobec procesu |

---

## 9. Stany, dane i powiązania

### 9.1. Stany funkcji

| Stan | Warunek | Postać wizualna |
|---|---|---|
| Bezczynny | Brak oczekujących sugestii, powierzchnia interakcji zamknięta | Awatar w barwie bazowej, plakietka ukryta |
| Sugestia oczekująca | Co najmniej jedna sugestia o statusie `nowa` | Plakietka liczbowa na awatarze |
| Doradza | Dymek kontekstowy otwarty z treścią sugestii | Dymek przy awatarze, awatar w akcencie |
| Rozmowa | Powierzchnia interakcji otwarta, trwa wymiana z agentem | Kolumna boczna otwarta, strumień rozmowy aktywny |
| Nasłuch | Tor głosowy przyjmuje polecenie | Pierścień pulsujący, ikona mikrofonu na awatarze |
| Mówi | Trwa synteza mowy | Obwódka w akcencie złotym, animacja przebiegu |
| Wyciszony | Aktywne wyciszenie czasowe, kontekstowe albo klasy zdarzeń | Awatar przygaszony, plakietka ukryta |
| Interweniuje | Wykonywane działanie sterujące procesem (zatwierdzenie, wstrzymanie, wznowienie) | Wskaźnik przebiegu przy działaniu w dymku lub w powierzchni interakcji |

Maszyna stanów funkcji — przejścia między ośmioma stanami z tabeli powyżej. Wyciszenie i wygaśnięcie wyciszenia są przejściami dostępnymi z każdego stanu poza „Interweniuje” (rozdz. 3.5 — zdarzenia wagi krytycznej ujawniają się mimo wyciszenia).

```
                              zdarzenie klasy wyzwalającej
                              przekracza próg (rozdz. 3.4)
        ┌──────────────┐  ───────────────────────────►  ┌─────────────────────┐
        │   Bezczynny  │                                │ Sugestia oczekująca │
        └──────┬───────┘  ◄───────────────────────────  └──────────┬──────────┘
               │             sugestia rozstrzygnięta               │
               │             (przyjęta/odrzucona/wygasła)          │ kliknięcie awatara
               │                                                   ▼
       fraza wybudzająca                                  ┌─────────────────┐
       albo skrót głosowy                                 │     Doradza     │
               │                                          └────┬─────────┬──┘
               ▼                                    kliknięcie │         │ brak wyboru,
        ┌──────────────┐                              podwójne │         │ odejście fokusu
        │    Nasłuch   │                                       ▼         ▼
        └──────┬───────┘                                ┌─────────────┐  Bezczynny
               │ transkrypcja                           │  Rozmowa    │
               │ gotowa                                 └──────┬──────┘
               ▼                                               │ wybór działania
        ┌──────────────┐    odpowiedź syntezowana              │ sterującego procesem
        │     Mówi     │ ◄─────────────────────────────────────┘
        └──────┬───────┘                                       ▼
               │ synteza                              ┌─────────────────┐
               │ zakończona                           │  Interweniuje   │
               ▼                                      └────────┬────────┘
           Bezczynny                                           │ działanie wykonane
                                                               ▼
                                                           Bezczynny

        Z każdego stanu poza „Interweniuje”:
     ─── polecenie wyciszenia (rozdz. 3.5) ──► Wyciszony ─── wygaśnięcie ──► Bezczynny
```

Legenda: przejścia oznaczone wyzwalaczem nad albo pod strzałką; „Sugestia oczekująca” i „Doradza” różnią się obecnością otwartego dymka — kliknięcie awatara przenosi z pierwszego do drugiego, dymek zamyka się samoczynnie po wyborze działania albo po odejściu fokusu bez wyboru.

### 9.2. Tryby obecności

| Tryb | Zachowanie |
|---|---|
| Pełny | Awatar widoczny, sugestie ujawniane zgodnie z progami, synteza mowy czynna |
| Cichy | Awatar widoczny, sugestie gromadzone bez dymka, synteza mowy wyłączona, wyjątek wagi krytycznej zachowany |
| Ukryty | Awatar poza polem widzenia, funkcja czynna w tle; dostęp przez szybki wybór szyny nawigacji, skrót klawiszowy, wyszukiwarkę funkcji i polecenie języka naturalnego |

### 9.3. Encje modelu danych

| Encja | Rozdział źródła | Zastosowanie w funkcji |
|---|---|---|
| `sugestia_aod` | [Model danych](../architektura/model-danych.md), rozdz. 18.2 | Zapis pojedynczej sugestii: `kontekst_typ`, `kontekst_id`, `rodzaj`, `tresc`, `status`, `znacznik_czasu` |
| `polecenie_glosowe` | [Model danych](../architektura/model-danych.md), rozdz. 18.3 | Zapis polecenia głosowego toru globalnego; `profil_asystenta_id` pozostaje pusty dla poleceń wydanych bezpośrednio do funkcji |
| `ustawienie` | [Model danych](../architektura/model-danych.md), rozdz. 17 | Zapis ustawień zakresu „Always On Display”: obecność, tryb, progi, reguły wyciszania |
| `sesja`, `zadanie`, `kolejka`, `przebieg_multitasking`, `automatyka` | [Model danych](../architektura/model-danych.md), rozdz. 8–13 | Odczyt stanu procesów stanowiących źródło sygnału (rozdz. 3.3) |

Odwzorowanie kontekstu sugestii: `kontekst_typ` przyjmuje wartości `srodowisko`, `modul`, `projekt`, `sesja`; `kontekst_id` wskazuje odpowiednio `srodowisko.id`, `modul.id`, `projekt.id` albo `sesja.id`.

Pola encji `sugestia_aod` ([Model danych](../architektura/model-danych.md), rozdz. 18.2):

| Pole | Typ | Opis |
|---|---|---|
| `id` | INTEGER | Klucz główny |
| `kontekst_typ` | TEXT | `srodowisko` \| `modul` \| `projekt` \| `sesja` |
| `kontekst_id` | INTEGER | Odniesienie polimorficzne do elementu, którego sugestia dotyczy; NULL dopuszczalny |
| `rodzaj` | TEXT | `doradztwo` \| `konfiguracja` \| `problem` \| `kolejny_krok` (rozdz. 4) |
| `tresc` | TEXT | Treść sugestii prezentowana w dymku i na liście oczekujących |
| `status` | TEXT | `nowa` \| `przyjeta` \| `odrzucona` (rozdz. 9.4) |
| `znacznik_czasu` | TEXT | Moment powstania sugestii |

Pola encji `polecenie_glosowe` ([Model danych](../architektura/model-danych.md), rozdz. 18.3) — tor globalny funkcji rozpoznawalny po pustym `profil_asystenta_id`:

| Pole | Typ | Opis |
|---|---|---|
| `id` | INTEGER | Klucz główny |
| `profil_asystenta_id` | INTEGER | Puste dla poleceń wydanych bezpośrednio przez Always On Display — obecne wyłącznie po przekazaniu do modułu Assistant (rozdz. 5.3) |
| `sesja_id` | INTEGER | Sesja, w kontekście której polecenie zostało wydane; NULL dopuszczalny dla poleceń ponadkontekstowych |
| `tresc_polecenia` | TEXT | Transkrypcja polecenia głosowego |
| `status` | TEXT | `w_trakcie` \| `zakonczone` \| `blad` |
| `wynik` | TEXT | Odpowiedź zwrócona Operatorowi; NULL dopuszczalny |
| `znacznik_czasu` | TEXT | Moment wydania polecenia |

### 9.4. Cykl życia sugestii

```
Schemat 3 — Cykl życia sugestii

  zdarzenie źródłowe (rozdz. 3.2)
        │  próg spełniony, brak wyciszenia
        ▼
  sugestia utworzona   status: nowa   ──► plakietka / dymek (zależnie od wagi)
        │
   ┌────┴───────────────┬──────────────────┬────────────────────────┐
   ▼                    ▼                  ▼                        ▼
 działanie          „Odłóż”            „Odrzuć”            brak reakcji 24 h
 wykonane        pozycja wraca      status: odrzucona     status: odrzucona
 status: przyjeta  do listy          bez powtórzenia        pozycja znika
```

Odwzorowanie wyniku cyklu na pole `status` encji `sugestia_aod` (rozdz. 9.3): gałąź „działanie wykonane” i gałąź „Odrzuć” zapisują się natychmiast (`przyjeta` albo `odrzucona`); gałąź „Odłóż” pozostawia `status` bez zmiany (`nowa`) — sugestia wraca do listy oczekujących, nie tworzy nowego wiersza; brak reakcji przez 24 godziny zapisuje `odrzucona` bez udziału Operatora, a decyzja rozgłasza się zdarzeniem `aod.suggestion.changed` w każdej z czterech gałęzi.

### 9.5. Komunikaty kontraktu komunikacji

| Komunikat | Kierunek | Zastosowanie |
|---|---|---|
| `aod.status.get` | polecenie | Odczyt stanu aktywacji i widoczności awatara |
| `aod.chat.send` | polecenie | Wiadomość tekstowa skierowana do funkcji, poza kontekstem karty sesji |
| `aod.voice.command` | polecenie | Polecenie głosowe toru globalnego (rozdz. 5) |
| `aod.context.get` | polecenie | Odczyt kontekstu: aktywny moduł, środowisko, projekt, historia Operatora |
| `aod.suggestion` | polecenie | Odczyt sugestii: lista oczekujących albo historia przyjętych i odrzuconych |
| `aod.suggestion.decide` | polecenie | Decyzja Operatora wobec sugestii: przyjęcie, odłożenie albo odrzucenie (rozdz. 9.4) |
| `aod.suggestion.changed` | zdarzenie | Proaktywna sugestia powstała albo zmieniła status: doradztwo, konfiguracja, wskazanie problemu, kolejny krok |
| `aod.observe.attach` | polecenie | Dołączenie funkcji do procesu środowiska MultitaskingAI w trybie `observer` albo `operator` |
| `aod.observe.detach` | polecenie | Odłączenie funkcji od procesu |
| `config.changed` | zdarzenie | Rozgłoszenie zmiany ustawień zakresu „Always On Display” na wszystkie urządzenia Operatora |

Pełne brzmienie kontraktu zawierają [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md), rozdz. 16.2 (funkcje globalne) oraz rozdz. 11.6 (nadzór nad procesem środowiska MultitaskingAI).

### 9.6. Powiązania

| Powiązanie | Treść |
|---|---|
| Chat Window | Odbiorca sugestii przeniesionej jako polecenie (rozdz. 7.2) |
| Execution Loop Window | Źródło sygnału i miejsce skutku działań sterujących (rozdz. 7.3) |
| Funkcja globalna Mobile | Kanał wykonania interwencji zainicjowanej przez funkcję poza stanowiskiem roboczym; powiadomienie o sugestii wagi wysokiej dociera kanałem Mobile zgodnie z ustawieniami powiadomień |
| Moduł Assistant | Tor sesyjny głosu; przekazanie polecenia w obie strony (rozdz. 5.3) |
| Moduł Agents | Źródło agentów przypisywanych działaniem „Przypisz agenta” sugestii konfiguracji |
| Moduł Automations | Odbiorca działania „Utwórz automatykę”; źródło sygnału klasy „harmonogram” |
| Środowisko MultitaskingAI | Rola warstwy centralnej z trybami obserwatora i operatora ([Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2) |

---

## 10. Punkty sterowania z okna konfiguracji i z okna ustawień

### 10.1. Punkty sterowania z okna konfiguracji

| Ustawienie | Opis | Warstwa konfiguracji | Wartość domyślna | Warstwa widoczności |
|---|---|---|---|---|
| Obecność funkcji | Widoczność awatara: widoczny albo ukryty przy zachowaniu działania w tle | globalna, środowisko, sesja | widoczny | 2 |
| Tryb obecności | Pełny, cichy, ukryty (rozdz. 9.2) | globalna, środowisko, sesja | pełny | 2 |
| Klasy zdarzeń wyzwalających | Włączenie i wyłączenie poszczególnych klas z rozdz. 3.2 | globalna, środowisko, moduł | wszystkie klasy czynne | 4 |
| Progi wyzwalania | Wartości progów z rozdz. 3.4 | globalna, środowisko, moduł | wartości z rozdz. 3.4 | 4 |
| Częstotliwość ujawniania | Liczba dymków w godzinie i odstęp między nimi | globalna, sesja | 3 na godzinę, odstęp 5 minut | 4 |
| Waga ujawniana samoczynnie | Najniższa waga sugestii otwierającej dymek bez interakcji | globalna, środowisko | wysoka | 4 |
| Reguły wyciszania | Zakresy i czasy wyciszenia dostępne w menu funkcji | globalna | komplet z rozdz. 3.5 | 4 |
| Tor głosowy funkcji | Włączenie toru głosowego, globalna fraza wybudzająca, synteza mowy odpowiedzi | globalna, sesja | tor czynny, synteza czynna | 2 |
| Tryb wobec procesu | Obserwator albo operator wobec procesów środowiska MultitaskingAI | globalna, projekt, sesja | obserwator | 2 |
| Czas życia sugestii nieprzyjętej | Czas, po którym sugestia otrzymuje status odrzuconej | globalna | 24 godziny | 4 |

Zmiana każdego ustawienia jest odwracalna, zapisywana natychmiast i rozgłaszana zdarzeniem `config.changed` na wszystkie urządzenia Operatora. Każda pozycja niesie objaśnienie kontekstowe `[?]` zgodnie z opracowaniem [Model konfiguracji](../architektura/model-konfiguracji.md), rozdz. 3.3.

### 10.2. Punkty sterowania z okna Ustawień

Sekcja „Always On Display” okna Ustawień ([Ustawienia](../interfejs-uzytkownika/ustawienia.md), rozdz. 11) udostępnia Operatorowi zestaw pozycji o codziennym zastosowaniu, bez wchodzenia w pełny zakres okna konfiguracji.

| Pozycja sekcji | Zawartość | Warstwa | Sposób wywołania |
|---|---|---|---|
| Obecność awatara | Przełącznik widoczności awatara | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Tryb obecności | Wybór trybu: pełny, cichy, ukryty | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |
| Reguły wyzwalania | Tabela klas zdarzeń z przełącznikami i wartościami progów | 4 | Tryb administracyjny, wyszukiwarka funkcji, polecenie języka naturalnego w Chat Window |
| Wyciszanie | Wybór czasu i zakresu wyciszenia oraz podgląd wyciszeń czynnych | 3 | Menu akcji `⋮` w sekcji |
| Tor głosowy | Przełącznik toru głosowego, globalna fraza wybudzająca, synteza mowy odpowiedzi | 2 | Sekcja „Always On Display” w selektorze sekcji `☰` |

### 10.3. Sterowanie poleceniem języka naturalnego

Każdy punkt sterowania jest osiągalny poleceniem języka naturalnego skierowanym do Wykonawcy w Chat Window — na przykład „wycisz sugestie do końca dnia”, „ukryj awatar w tym module”, „podnieś próg ponowień do czterech”, „przełącz agenta towarzyszącego w tryb operatora”. Polecenie wykonuje tę samą operację, którą wykonuje element interfejsu, i podlega tym samym regułom zapisu.

---

## 11. Scenariusze użycia

### 11.1. Wynik kontroli jakości odrzucony po raz drugi

Operator pracuje w module Studio. W środowisku CodeStudio trwa równolegle zlecenie, którego kontrola jakości odrzuca wynik zadania po raz drugi. Zdarzenie należy do klasy „wynik kontroli jakości” i przekracza próg ponowień, więc powstaje sugestia rodzaju `problem` o wadze wysokiej. Dymek kontekstowy otwiera się przy awatarze bez przerywania pracy w Studio. Operator wybiera działanie „Otwórz Execution Loop Window”: okno pętli otwiera się jako kolumna sąsiadująca z Chat Window na zadaniu, którego dotyczy problem. Operator koryguje zlecenie w Chat Window; sugestia otrzymuje status `przyjeta`.

### 11.2. Powtarzalność czynności ręcznej

Operator trzykrotnie wykonuje w module Developer tę samą sekwencję poleceń. Próg powtarzalności zostaje przekroczony i powstaje sugestia rodzaju `konfiguracja` o wadze niskiej — bez dymka, jako pozycja listy oczekujących z plakietką na awatarze. Operator otwiera powierzchnię interakcji, rozwija pozycję i wybiera „Utwórz automatykę”. Moduł Automations otwiera się z formularzem wypełnionym rozpoznaną sekwencją.

### 11.3. Krótkie polecenie głosowe poza modułem Assistant

Operator pracuje w module Terminal i wypowiada globalną frazę wybudzającą, a następnie polecenie „ile zadań czeka w kolejce”. Tor globalny rozpoznaje krótką frazę, odczytuje stan kolejek przez `aod.context.get` i odpowiada syntezą mowy oraz wpisem w strumieniu powierzchni interakcji. Moduł Assistant nie zostaje otwarty, karta sesji modułu Terminal pozostaje bez zmian.

### 11.4. Polecenie głosowe przekazane do modułu Assistant

Operator wypowiada do awatara polecenie „podyktuj akapit do bieżącego dokumentu”. Polecenie wykracza poza zakres toru globalnego, więc następuje przekazanie: moduł Assistant otwiera się z oknem Voice Console, rozpoznana treść trafia do rozmowy jako pierwsza wypowiedź, nasłuch przechodzi do toru sesyjnego bez ponownej aktywacji. Powierzchnia interakcji funkcji odnotowuje przekazanie i wskazuje okno docelowe.

### 11.5. Punkt decyzyjny pętli podczas pracy ciągłej

W środowisku MultitaskingAI trwa proces w trybie pracy ciągłej. Pętla dochodzi do punktu decyzyjnego oczekującego na zatwierdzenie. Powstaje sugestia rodzaju `kolejny_krok` o wadze wysokiej; funkcja działa w trybie operatora, więc dymek zawiera działanie „Zatwierdź krok”. Operator znajduje się poza stanowiskiem roboczym — powiadomienie dociera kanałem funkcji Mobile, a zatwierdzenie zostaje wykonane z urządzenia mobilnego. Pętla przechodzi do kroku następnego. Specyfikę ról i kolejek tego procesu opisuje [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2.

### 11.6. Praca wymagająca ciszy

Operator prowadzi w module Roundtable debatę wymagającą pełnej uwagi i wycisza sugestie na godzinę jednym poleceniem języka naturalnego w Chat Window. Awatar przechodzi w stan wyciszony, plakietka pozostaje ukryta, sugestie gromadzą się w liście oczekujących. Punkt decyzyjny wstrzymujący proces w innym środowisku ujawnia się mimo wyciszenia — jako plakietka bez dymka, zgodnie z wyjątkiem wagi krytycznej. Po upływie godziny wyciszenie wygasa samoczynnie.

### 11.7. Nadzór ciągły nad procesem od uruchomienia po zamknięcie zlecenia

Scenariusz łączy stanowisko stacjonarne, środowisko MultitaskingAI i funkcję Mobile w jednym cyklu procesu, obejmując pełen przebieg od wydania polecenia po zamknięcie zlecenia poza stanowiskiem.

| Krok | Miejsce | Działanie | Stan funkcji |
|---|---|---|---|
| 1 | Chat Window, środowisko MultitaskingAI | Operator wydaje polecenie uruchamiające proces w trybie pracy ciągłej; Coordinator dekomponuje zlecenie na zadania | Bezczynny |
| 2 | Panel orkiestracji | Operator przełącza tryb funkcji z obserwatora na operatora dla tego procesu (`aod.observe.attach`, pole `mode`) | Bezczynny, dołączona do procesu |
| 3 | Stanowisko stacjonarne | Operator opuszcza stanowisko; proces trwa w tle, Executor 1 i Executor 2 realizują zadania równolegle | Bezczynny (obserwacja bierna) |
| 4 | Rdzeń platformy | Kontrola jakości odrzuca wynik jednego z zadań; sygnał odkłada się przez `aod.signal.report` | Sugestia oczekująca |
| 5 | Urządzenie przenośne | Powiadomienie push dociera kanałem funkcji Mobile; Operator otwiera podgląd procesu | Doradza (na urządzeniu przenośnym) |
| 6 | Mobile | Operator zatwierdza ponowienie zadania — działanie sterujące wykonuje się tym samym portem, którym steruje okno pętli na stanowisku | Interweniuje |
| 7 | Rdzeń platformy | Ponowienie kończy się powodzeniem; Coordinator zamyka zlecenie (`loop.order.close`, `outcome: completed`) | Powrót do stanu bezczynnego |
| 8 | Stanowisko stacjonarne | Po powrocie Operator zastaje pełną historię przebiegu w Execution Loop Window, bez luki w zapisie zdarzeń | Bezczynny |

Scenariusz pokazuje trzy własności funkcji działające naraz: rolę warstwy centralnej w środowisku MultitaskingAI (rozdz. 6.3), przekazanie interwencji do funkcji Mobile bez utraty ciągłości sterowania (rozdz. 6.3, [Środowisko MultitaskingAI](../srodowiska/multitaskingai.md) rozdz. 2) oraz zachowanie pełnej historii zdarzeń niezależnie od urządzenia, z którego prowadzony był nadzór.

---

## Załącznik A. Skróty klawiszowe i ikonografia

### A.1. Skróty klawiszowe

| Skrót | Działanie |
|---|---|
| `Ctrl/Cmd + Shift + A` | Otwarcie i zamknięcie powierzchni interakcji Always On Display |
| `Ctrl/Cmd + Shift + V` | Uruchomienie nasłuchu toru głosowego funkcji |
| `Ctrl/Cmd + Shift + Q` | Otwarcie listy oczekujących sugestii |
| `Ctrl/Cmd + Shift + M` | Wyciszenie sugestii na 15 minut i zniesienie wyciszenia |
| `Ctrl/Cmd + Shift + H` | Ukrycie i przywrócenie awatara |
| `Ctrl/Cmd + Enter` | Zatwierdzenie kroku oczekującego na decyzję z poziomu dymka sugestii |
| `Esc` | Zamknięcie dymka kontekstowego bez zmiany statusu sugestii |

### A.2. Ikonografia

| Ikona | Zastosowanie |
|---|---|
| `dzwonek` | Punkt dostępu do funkcji w szybkim wyborze szyny nawigacji i w listwie ustawień strefy 3 strony głównej |
| `uzytkownik` | Znak zastępczy wewnątrz `.dn-aod-rdzen`, gdy nie skonfigurowano obrazu awatara |
| `mikrofon` | Przycisk toru głosowego w powierzchni interakcji i w dymku sugestii |
| `oko` | Tryb obserwatora wobec procesu |
| `ustawienia` | Wejście do konfiguracji reguł wyzwalania |
| `wiecej` | Menu kebab (⋮) powierzchni interakcji |
| `zegar` | Wyciszenie czasowe |
| `ostrzezenie` | Sugestia rodzaju `problem` |
| `strzalka-prawo` | Sugestia rodzaju `kolejny_krok` |

### A.3. Komponenty systemu wizualnego

| Komponent | Klasa | Zastosowanie |
|---|---|---|
| Awatar | `.dn-awatar` (wariant `--lg`, pozycjonowanie stałe) | Awatar funkcji |
| Plakietka stanu | `.dn-plakietka`, `.dn-kropka` | Plakietka liczby oczekujących sugestii |
| Awatar w pozycjonowaniu pływającym | `.dn-aod`, `.dn-aod-rdzen`, `.dn-aod-tresc` | Awatar Always On Display w stanie bezczynnym i w stanach z treścią towarzyszącą |
| Dymek kontekstowy | `ao-dymek`, `ao-dymek-glowa`, `ao-dymek-tresc`, `ao-dymek-stopka` (selektory lokalne prototypu) | Dymek kontekstowy sugestii |
| Karta pozycji | `.dn-karta--klikalna` | Pozycja listy oczekujących sugestii |
| Przycisk ikonowy | `.dn-btn-ikona` | Przycisk mikrofonu, menu kebab |
| Przycisk główny | `.dn-btn--sygnal` | Zatwierdzenie kroku w dymku sugestii |
| Suwak | `.dn-suwak` | Przełącznik trybu wobec procesu |
| Zakładki pigułkowe | `.dn-zakladki--pigulki` | Przełącznik trybu obecności |
| Pole tekstowe | `.dn-pole-kontrolka` | Pole polecenia do funkcji |


---

## Załącznik B. Pełny wykaz komend kontraktu Always On Display

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `aod` — 12 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `aod.status.get` | Zwraca stan nakładki Always On Display | `deviceId:string` (opc) | `status:AodStatus` (wym) |
| `aod.chat.send` | Wiadomość tekstowa skierowana do Always On Display poza kontekstem karty sesji (rozdz. 9.5). Odpowiedź agenta towarzyszącego trafia do strumienia powierzchni interakcji, nie do Chat Window: polecenie odczytu stanu, nawigacji, sterowania procesem i wyciszenia wykonuje się w torze globalnym, a pytanie otwarte trafia do modelu wraz z kontekstem platformy | `text:string` (wym)<br>`windowId:string` (opc)<br>`sessionId:string` (opc)<br>`moduleId:string` (opc) | `messageId:string` (wym)<br>`reply:string` (wym)<br>`windowId:string` (opc)<br>`moduleId:string` (opc)<br>`environmentId:string` (opc)<br>`handedOff:bool` (opc) |
| `aod.voice.command` | Polecenie głosowe toru globalnego Always On Display (rozdz. 5). Krótkie polecenie ponadkontekstowe — odczyt stanu, sterowanie procesem, nawigacja, wyciszenie — wykonuje się w torze globalnym i wraca odpowiedzią; polecenie wymagające dyktowania, rozmowy wieloetapowej, makra albo profilu głosu przechodzi do modułu Assistant jako zlecenie (przekazanie między torami, rozdz. 5.3) | `audioRef:string` (opc)<br>`transcript:string` (opc)<br>`windowId:string` (opc)<br>`sessionId:string` (opc)<br>`speak:bool` (opc) | `transcript:string` (wym)<br>`reply:string` (wym)<br>`handedOff:bool` (wym)<br>`action:AssistantAction` (opc)<br>`moduleId:string` (opc)<br>`environmentId:string` (opc)<br>`speechRef:string` (opc) |
| `aod.context.get` | Pobiera komplet kontekstu okna do nakładki | `windowId:string` (opc)<br>`sessionId:string` (opc)<br>`historyLimit:int` (opc) | `context:ContextBundle` (wym) |
| `aod.suggestion` | Zwraca sugestie Always On Display — listę oczekujących (status `nowa`) albo historię przyjętych i odrzuconych. Sugestie powstają w rdzeniu ze zdarzeń klas wyzwalających po przekroczeniu progu (rozdz. 3); odczyt przelicza je z bieżącej telemetrii i sygnałów, a powstanie oraz decyzja rozgłaszają się zdarzeniem `aod.suggestion.changed` | `windowId:string` (opc)<br>`sessionId:string` (opc)<br>`processId:string` (opc)<br>`status:AodSuggestionStatus` (opc)<br>`limit:int` (opc) | `suggestions:AodSuggestion[]` (wym)<br>`pendingCount:int` (wym) |
| `aod.suggestion.decide` | Decyzja Operatora wobec sugestii — trzy gałęzie cyklu życia z rozdz. 9.4: przyjęcie po wykonanym działaniu, odłożenie (status pozostaje `nowa`) i odrzucenie. Działanie sterujące procesem (zatwierdzenie kroku, wstrzymanie, wznowienie, ponowienie) wykonuje rdzeń tym samym portem, którym steruje okno pętli wykonawczej; działania powłoki (otwarcie okna, przeniesienie do Chat Window) wykonuje klient i zgłasza przyjęcie | `suggestionId:string` (wym)<br>`decision:AodSuggestionDecision` (wym)<br>`actionKind:AodSuggestionActionKind` (opc)<br>`deviceId:string` (opc) | `suggestion:AodSuggestion` (wym)<br>`pendingCount:int` (wym)<br>`reply:string` (opc) |
| `aod.observe.attach` | Dołącza Always On Display do procesu środowiska MultitaskingAI w trybie obserwatora albo operatora ([Środowisko MultitaskingAI](../srodowiska/multitaskingai.md), rozdz. 2.2). Tryb operatora bez procesu w biegu jest dopuszczalny — nadzór obejmie najbliższe uruchomienie | `processId:string` (wym)<br>`mode:AodProcessMode` (opc)<br>`deviceId:string` (opc) | `attachedProcessIds:string[]` (wym)<br>`mode:AodProcessMode` (wym) |
| `aod.observe.detach` | Odpina proces od obserwacji w nakładce | `processId:string` (wym)<br>`deviceId:string` (opc) | `attachedProcessIds:string[]` (wym) |
| `aod.mute.get` | Zwraca wyciszenia nakładki czynne w rdzeniu wraz z rodzajem, zakresem i chwilą końca. Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie potwierdza jego wygaśnięcia | `deviceId:string` (opc) | `mutes:AodMute[]` (wym) |
| `aod.mute.set` | Zakłada wyciszenie nakładki albo je znosi. Zniesienie realizuje ta sama komenda z polem `muted` równym fałszowi — odwracalność jednym ruchem jest wymogiem, nie udogodnieniem. Odpowiedź oddaje wykaz wyciszeń po zmianie, aby powłoka nie musiała pytać drugi raz. Zmiana rozgłasza się zdarzeniem `aod.mute.changed` na pozostałe powłoki Operatora | `muteId:string` (opc)<br>`kind:AodMuteKind` (opc)<br>`scope:AodMuteScope` (opc)<br>`scopeId:string` (opc)<br>`scopeName:string` (opc)<br>`eventClass:AodEventClass` (opc)<br>`endsAt:int64` (opc)<br>`muted:bool` (wym)<br>`deviceId:string` (opc) | `mutes:AodMute[]` (wym)<br>`changed:bool` (wym) |
| `aod.signal.report` | Odkłada w rdzeniu sygnał klasy zdarzeń wyzwalających. Nośnik dla trzech klas rozdz. 3.2, których rdzeń nie obserwuje własną telemetrią: wyniku kontroli jakości, harmonogramu przebiegów automatyk i powtarzalności czynności Operatora. Odpowiedź podaje wprost, czy sygnał trafił w wyciszenie i w które — sygnał wyciszony odkłada się nadal, ponieważ wyciszenie wstrzymuje **ujawnienie**, a nie zapis | `eventClass:AodEventClass` (wym)<br>`text:string` (wym)<br>`moduleId:string` (opc)<br>`sessionId:string` (opc)<br>`occurrenceCount:int` (opc) | `signal:AodSignal` (wym)<br>`suppressed:bool` (wym)<br>`mute:AodMute` (opc) |
| `aod.signal.list` | Zwraca sygnały klas zdarzeń odłożone w rdzeniu, bez tych wstrzymanych wyciszeniem. Odpowiedź nazywa wyciszenia, które sygnały wstrzymały, oraz liczbę wstrzymanych — cisza, po której Operator nie wie, że coś zostało wstrzymane, jest gorsza od braku wyciszenia | `deviceId:string` (opc)<br>`classes:AodEventClass[]` (opc)<br>`moduleId:string` (opc)<br>`sessionId:string` (opc)<br>`limit:int` (opc) | `signals:AodSignal[]` (wym)<br>`suppressedCount:int` (wym)<br>`mutes:AodMute[]` (wym) |

**Zdarzenia obszaru `aod` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `aod.mute.changed` | Wyciszenie nakładki założone albo zniesione. Rozgłoszenie przenosi ciszę na pozostałe powłoki Operatora: bez niego wyciszenie założone w jednej powłoce nie wstrzymywałoby sugestii w drugiej | — |
| `aod.suggestion.changed` | Sugestia Always On Display powstała albo zmieniła status. Proaktywna sugestia — doradztwo, konfiguracja, wskazanie problemu, kolejny krok (rozdz. 9.5) — dociera do każdej powłoki Operatora bez odpytywania, a plakietka awatara liczy pozycje z pola `pendingCount` | — |

Razem w wykazie: **12 komend** z 1 obszaru kontraktu.

---

## Załącznik C. Wykazy normatywne uzupełniające i kryteria odbioru

### C.1. Etykiety interfejsu z prototypu

Wykaz dosłownych brzmień z `design/05-okna/platformowe/always-on-display.html`, nieopisanych już w rozdziałach 1–11.

| Element | Dosłowne brzmienie | Miejsce wystąpienia |
|---|---|---|
| Menu kontekstowe awatara | „Pokaż nakładkę”, „Obserwuj proces”, „Polecenie głosowe”, „Sugestie oczekujące”, „Wycisz nakładkę”, „Ustawienia nakładki” | Menu kebab przy awatarze (rozdz. 2.6) |
| Nagłówek karty procesu obserwowanego | „Pętla ciągła 24/7 · zespół Budowa aplikacji”, „Badanie i redakcja · przebieg 12” | Powierzchnia interakcji w kontekście MultitaskingAI (rozdz. 6.3) |
| Wskaźnik trybu w panelu orkiestracji | „tryb AOD w monitorze”, „Obserwator”, „Operator” | Panel orkiestracji środowiska MultitaskingAI |
| Akcje sterowania procesem | „Uruchom proces w tle”, „Wstrzymaj (pause)”, „Przywróć stan wyjściowy” | Powierzchnia interakcji, karta procesu |
| Nagłówek hierarchii decyzji | „Hierarchia decyzji — poziom bieżący podświetlony”, kolumny „Poziom”, „Podmiot”, „Zakres decyzji” | Załącznik odniesienia rozdz. 8 Koncepcji platformy, zastosowanie w rozdz. 6.3 |
| Katalog oprzyrządowania AOD | „Podsumowanie”, „Oprzyrządowanie”, „Zakres dostępu”, „Stany”, „Katalog” | Sekcje karty referencyjnej funkcji w kontekście MultitaskingAI |
| Skrót stanu procesu obserwowanego | „Skrót stanu procesu obserwowanego przez Always On Display”, kolumny „Element”, „Stan bieżący” | Powierzchnia interakcji, podgląd procesu |
| Katalog zdolności funkcji | „Proaktywne doradztwo”, „Pomoc kontekstowa”, „Komunikacja głosowa”, „Komunikacja tekstowa”, „Pełny dostęp do platformy”, „Podgląd Monitora procesu” | Karta referencyjna zdolności funkcji |
| Katalog akcji interwencji | „Zatwierdzenie kroku”, „Wstrzymanie procesu”, „Wznowienie procesu”, „Interwencja z dowolnego miejsca” | Karta referencyjna akcji sterujących |
| Trzy stany funkcji | „Bezczynny”, „Doradza”, „Interweniuje” wraz z opisami „Brak aktywnej sugestii, gotowość do podglądu.”, „Wyświetla dymek z sugestią lub ostrzeżeniem, bez możliwości interwencji.”, „Aktywnie wykonuje zatwierdzenie, wstrzymanie lub modyfikację procesu, w tym przez Mobile.” | Karta referencyjna stanów — podzbiór stanów wykazanych w rozdz. 2.4 i 9.1 |
| Rola w interwencji zdalnej | „Obserwuje przebieg procesu i inicjuje interwencję — zatwierdzenie kroku, wstrzymanie procesu”, „Udostępnia kanał, przez który interwencja jest faktycznie wykonywana z dowolnego urządzenia, niezależnie od lokalizacji” | Opis współdziałania Always On Display i funkcji Mobile (rozdz. 6.3) |

### C.2. Żetony projektowe

| Miejsce zastosowania | Żeton `--dn-*` | Czego dotyczy |
|---|---|---|
| Awatar w stanie bezczynnym | `--dn-grad-sygnal`, `--dn-szary-0` | Tło gradientowe `.dn-aod-rdzen` i kolor znaku na nim (rozdz. 2.4) |
| Obwódka tętniąca awatara | `--dn-sygnal-obrys`, `--dn-czas-tetno` | Animacja `dn-tetno` sygnalizująca aktywność funkcji (`.dn-aod-rdzen::after`) |
| Plakietka sugestii oczekującej | `--dn-sygnal`, `--dn-sygnal-tlo` | Plakietka liczbowa na awatarze (rozdz. 2.4, 2.6) |
| Plakietka sugestii o wadze wysokiej | `--dn-ostrzezenie-tekst`, `--dn-ostrzezenie-tlo`, `--dn-ostrzezenie-obrys` | Plakietka pulsująca w akcencie złotym — zdarzenie wymagające decyzji (rozdz. 2.4) |
| Stan błędu kanału | `--dn-blad-tekst`, `--dn-blad-tlo`, `--dn-blad-obrys` | Błąd urządzenia wejściowego toru głosowego (rozdz. 2.6), błąd wysyłki `aod.chat.send` |
| Pozycjonowanie pływające | `--dn-od-5`, `--dn-wym-stan`, `--dn-wym-kontrolka`, `--dn-z-aod` | Odsunięcie od krawędzi i warstwa `.dn-aod` nad paskiem stanu (`design/zasoby/css/komponenty.css`) |
| Wymiar awatara | `--dn-wym-awatar-lg` | Średnica koła `.dn-aod-rdzen` (rozdz. 2.4) |
| Promień pigułki | `--dn-r-pill` | Kształt wrappera `.dn-aod` i plakietek stanu |
| Fokus kontrolek powierzchni interakcji | `--dn-fokus` | Obrys fokusu pola polecenia i przycisku mikrofonu |
| Czas przejścia dymka i powierzchni | `--dn-czas-3` | Otwarcie/zamknięcie dymka kontekstowego i powierzchni interakcji |

### C.3. Punkty łamania

Wartości progów z `design/zasoby/zetony/zetony.css`, rozdział „Siatka i punkty łamania”. Funkcja jest przekrojowa wobec wszystkich środowisk (rozdz. 1.2), więc progi obowiązują niezależnie od środowiska aktywnego w tle.

| Żeton | Próg szerokości | Zachowanie funkcji Always On Display |
|---|---|---|
| `--dn-bp-w1` | 640px | Funkcja dostępna w postaci uproszczonej z poziomu funkcji globalnej Mobile ([Funkcja globalna Mobile](mobile.md)) zamiast pływającego awatara |
| `--dn-bp-w2` | 960px | Powierzchnia interakcji otwiera się jako nakładka pełnej wysokości zamiast kolumny bocznej regulowanej szerokości (rozdz. 2.5) |
| `--dn-bp-w3` | 1280px | Zachowanie referencyjne — awatar pływający, dymek kontekstowy i powierzchnia interakcji jako kolumna boczna |
| `--dn-bp-w4` | 1600px | Bez zmian względem `--dn-bp-w3` — funkcja nie zajmuje dodatkowej szerokości przy biurku szerokim |

### C.4. Stany kontrolek

| Kontrolka | Spoczynek | Wskazanie kursorem | Wciśnięcie | Ognisko | Nieaktywny | Ładowanie | Pusty | Błąd |
|---|---|---|---|---|---|---|---|---|
| Awatar (`.dn-aod-rdzen`) | Barwa bazowa, bez animacji | Uniesienie subtelne | Otwarcie dymka kontekstowego | Obrys `--dn-fokus` | nie występuje (zero blokad) — stan „ukryty” to widoczność, nie dezaktywacja | — | — | — |
| Przycisk mikrofonu (`.dn-btn-ikona`) | Ikona spoczynkowa | Podświetlenie subtelne | Rozpoczęcie nasłuchu | Obrys `--dn-fokus` | nie występuje | Pierścień pulsujący w trakcie nasłuchu | — | Plakietka `--dn-blad-tekst`, „błąd urządzenia wejściowego” |
| Pole polecenia (`.dn-pole-kontrolka`) | Obrys `--dn-obrys` | Obrys `--dn-sygnal-obrys` | — | Obrys `--dn-fokus` | nie występuje | Ikona wysyłki zastąpiona wskaźnikiem strumieniowania | Tekst podpowiedzi | Obrys `--dn-blad-obrys`, komunikat błędu kanału |
| Przełącznik trybu obecności (`.dn-zakladki--pigulki`) | Segment aktywny podświetlony | Podświetlenie segmentu wskazanego | Zmiana segmentu | Obrys `--dn-fokus` | nie występuje | — | — | — |
| Pozycja listy oczekujących sugestii (`.dn-karta--klikalna`) | Tło `--dn-tlo` | Tło `--dn-powierzchnia` | Rozwinięcie pozycji | Obrys `--dn-fokus` | nie występuje | — | „Brak oczekujących sugestii” | — |

### C.5. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Awatar jest widoczny nad powłoką niezależnie od środowiska i modułu aktywnego | Przełączenie kolejno między TalkIn, WorkSpace, CodeStudio, MultitaskingAI i stroną główną; potwierdzenie stałej obecności awatara |
| Żaden element funkcji nie jest trwale nieaktywny | Przegląd kontrolek wykazanych w rozdz. 2.6 i C.4 pod kątem stanu `disabled`; zero wystąpień poza stanem „ładowanie” |
| Wyciszenie sugestii nie ukrywa awatara ani nie usuwa dostępu do funkcji | Włączenie wyciszenia (rozdz. 3.5); potwierdzenie, że awatar pozostaje widoczny w stanie „wyciszony” |
| Tor głosowy globalny i tor modułu Assistant nie kolidują | Wydanie krótkiego polecenia głosowego z toru globalnego oraz polecenia wieloetapowego wymagającego przekazania do modułu Assistant (rozdz. 5.3); potwierdzenie właściwego skierowania obu |
| Sugestia zgłoszona zdarzeniem `aod.suggestion.changed` dociera do każdej powłoki Operatora bez odpytywania | Wywołanie `aod.signal.report` na jednym urządzeniu; obserwacja aktualizacji plakietki na drugim urządzeniu tego samego konta |
| Środowisko MultitaskingAI — rola centralna funkcji nie koliduje z rolą globalną | Dołączenie funkcji do procesu (`aod.observe.attach`) w środowisku MultitaskingAI; potwierdzenie zgodności z rozdz. 6.3 |
| Układ odpowiada progom łamania rozdz. C.3 na czterech szerokościach referencyjnych | Zmiana szerokości okna kolejno przez 640px, 960px, 1280px, 1600px; obserwacja zachowań z tabeli rozdz. C.3 |

---

*Koniec dokumentu. Always On Display — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
