# Danaco Console — Model konfiguracji

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Opis** | Platforma jest wielośrodowiskowym systemem operacyjnym dla sztucznej inteligencji, integrującym komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie, automatyzacje procesów oraz rozwój oprogramowania. |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Źródło** | Koncepcja platformy i architektura |

Dokument opisuje model konfiguracji produktu Danaco Console: pełną strukturę okna konfiguracji, wszystkie zakresy ustawień wpływających na aplikację, procesy, akcje, zachowanie modeli, tożsamość modeli, prompty systemowe, rozszerzenia, integracje, historię, pamięć, izolację, karty sesji, komponenty własne, oba kanały komunikacji operacyjnej (Chat Window oraz Execution Loop Window) i warstwy widoczności funkcji, a także warstwowość konfiguracji (globalna → środowisko → projekt → sesja) i mechanizm objaśnień kontekstowych towarzyszący każdemu ustawieniu. Dokument realizuje zasadę pełnej konfigurowalności ustanowioną w rozdziale 4 i w rozdziale 12 dokumentu Koncepcja platformy oraz w rozdziałach 5, 9, 10, 12 i 13 dokumentu Architektura.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
1. [Cel i zakres dokumentu](#1-cel-i-zakres-dokumentu)
2. [Zasada nadrzędna: pełna konfigurowalność](#2-zasada-nadrzędna-pełna-konfigurowalność)
3. [Struktura okna konfiguracji](#3-struktura-okna-konfiguracji)
4. [Warstwowość konfiguracji](#4-warstwowość-konfiguracji)
5. [Zakresy ustawień](#5-zakresy-ustawień)
6. [Okno konfiguracji punktów izolacji](#6-okno-konfiguracji-punktów-izolacji)
7. [Warstwy widoczności funkcji](#7-warstwy-widoczności-funkcji)
8. [Profile konfiguracji](#8-profile-konfiguracji)
9. [Komunikacja zmian konfiguracji](#9-komunikacja-zmian-konfiguracji)
10. [Zgodność z zasadami nadrzędnymi platformy](#10-zgodność-z-zasadami-nadrzędnymi-platformy)
11. [Słowniczek pojęć](#11-słowniczek-pojęć)
- [Załącznik A. Zbiorcza macierz zakresów ustawień i warstw](#załącznik-a-zbiorcza-macierz-zakresów-ustawień-i-warstw)
- [Załącznik B. Szablony](#załącznik-b-szablony)
- [Załącznik C. Scenariusze konfiguracji](#załącznik-c-scenariusze-konfiguracji)

---

## Wprowadzenie

Dokument Architektura ustanawia w rozdziale 5.2 istnienie okna konfiguracji jako jednego z dwóch okien realizujących interfejs użytkownika platformy oraz wskazuje jego zakres w formie wyliczenia: aplikacja, procesy, akcje, zachowanie modeli, tożsamość modeli, prompty systemowe, rozszerzenia, integracje, historia, pamięć i izolacja. Dokument Koncepcja platformy ustanawia w rozdziale 4 zasadę pełnej konfigurowalności jako centralną wobec pozostałych zasad nadrzędnych platformy (rozdz. 12, zasada 2), a w rozdziałach 4.4–4.6 opisuje interfejs jednego wycinka tej konfiguracji — okna konfiguracji punktów izolacji.

Niniejszy dokument opisuje pełną strukturę okna konfiguracji jako całości. Do jedenastu zakresów wskazanych w Architekturze dodaje pięć dalszych, wynikających wprost z mechanizmów platformy: karty sesji, komponenty własne, Chat Window, Execution Loop Window i warstwy widoczności funkcji. Interfejs okna izolacji osadzony jest w tej strukturze jako jeden z szesnastu zakresów ustawień. Model konfiguracji rozwija zasady przyjęte w Koncepcji platformy i w Architekturze do poziomu szczegółowości pozwalającego na projektowanie interfejsu i implementację.

Miejsce dokumentu w łańcuchu źródeł i przedmiot każdego z nich przedstawia poniższe zestawienie.

| Dokument źródłowy | Co ustala | Odniesienie w tym dokumencie |
|---|---|---|
| Koncepcja platformy | Zasada pełnej konfigurowalności; okno konfiguracji punktów izolacji; role i mechanizmy MultitaskingAI | rozdz. 2, 4, 6, 7 |
| Architektura | Istnienie i zakres okna konfiguracji; osiem zakresów izolacji technicznej; warstwowość i kanały komunikacji operacyjnej | rozdz. 3, 5, 6, 9 |
| Model danych | Model danych ustawienia; retencja historii; zasoby pamięci; profil zapisany na warstwie | rozdz. 3.3, 4.2, 5.9, 5.10 |

---

## 1. Cel i zakres dokumentu

Celem dokumentu jest ustalenie pełnej struktury okna konfiguracji Danaco Console — jedynego miejsca, z którego konfiguruje się dosłownie wszystkie zależności i zachowania platformy, zgodnie z zasadą pełną konfigurowalności (Koncepcja platformy, rozdz. 6 i rozdz. 14, zasada 2). Zakres dokumentu obejmuje:

| Element zakresu | Opis | Rozdział |
|---|---|---|
| Warstwowość konfiguracji | Cztery poziomy ogólne (globalny, środowisko, projekt, sesja) oraz relacja tej warstwowości do siedmiu poziomów zasięgu właściwych izolacji | rozdz. 4 i rozdz. 6 |
| Szesnaście zakresów ustawień | Konkretne pozycje ustawień, ich zastosowanie warstwowe i wartości domyślne | rozdz. 5 |
| Mechanizm objaśnień kontekstowych | Objaśnienie towarzyszące każdemu ustawieniu | rozdz. 3.3 |
| Interfejs okna konfiguracji punktów izolacji | Trzypanelowe okno łączące izolację kontekstu i izolację techniczną (Koncepcja platformy, rozdz. 6.4–6.6) | rozdz. 6 |
| Warstwy widoczności funkcji | Przypisanie funkcji do warstw 1–4, profile warstw wg roli użytkownika, wyszukiwarka funkcji, skróty klawiszowe, tryb administracyjny | rozdz. 7 |
| Mechanizm profili konfiguracji | Zapis i ponowne użycie raz zbudowanej konfiguracji | rozdz. 8 |
| Komunikacja zmian konfiguracji | Sposób przekazywania zmian między klientem a serwerem zgodnie z kontraktami komunikacji | rozdz. 9 |

Wymiary, odstępy oraz tokeny kolorów i typografii okna konfiguracji określają przewodnik marki i pliki tokenów systemu wizualnego, analogicznie do rozwiązania przyjętego dla strony głównej w rozdziale 7.4 Koncepcji platformy.

---

## 2. Zasada nadrzędna: pełna konfigurowalność

Model konfiguracji opisany w niniejszym dokumencie jest bezpośrednią realizacją zasady centralnej platformy, ustanowionej w rozdziale 14 Koncepcji platformy (zasada 2): **każdą funkcję, każde zachowanie i każdą zależność między elementami platformy konfiguruje się z poziomu okna konfiguracji.** Żadne zachowanie platformy nie jest zaszyte na stałe w sposób niedostępny dla użytkownika.

Z tej zasady wynikają dla modelu konfiguracji cztery konsekwencje wiążące dla każdego z szesnastu zakresów opisanych w rozdziale 5.

| Konsekwencja zasady | Znaczenie dla każdego z szesnastu zakresów (rozdz. 5) | Podstawa lub przykład |
|---|---|---|
| Brak twardych blokad | Żaden zakres ustawień nie zawiera wartości niemożliwej do zmiany; istniejące ograniczenia techniczne są parametrem technicznym, a nie decyzją produktową ukrywającą ustawienie przed użytkownikiem | górny limit piętnastu podagentów w mechanizmie Subagent Network (Koncepcja platformy, rozdz. 13.3) |
| Brak ustawienia = wartość domyślna | Nieskonfigurowanie pozycji nie blokuje pracy i nie wymaga interwencji użytkownika — pozycja przyjmuje wartość domyślną albo dziedziczy wartość z warstwy szerszej | reguła dziedziczenia (rozdz. 4.2) |
| Jawność zależności | Każde powiązanie między środowiskami, modułami, komponentami własnymi i sesjami jest widoczne w oknie konfiguracji jako możliwość do świadomego ustanowienia, a nie jako ukryta reguła platformy | podpunkty „Powiązania" (Koncepcja platformy, rozdz. 11); rozdz. 12, zasada 3 |
| Objaśnienie jako część definicji ustawienia | Każdy element konfiguracji zawiera objaśnienie kontekstowe opisujące jego działanie i wpływ na aplikację; treść objaśnienia jest częścią definicji ustawienia | Architektura, rozdz. 13; rozdz. 3.3 niniejszego dokumentu |

---

## 3. Struktura okna konfiguracji

### 3.1. Dostęp do okna konfiguracji

Okno konfiguracji jest dostępne z dwóch miejsc platformy, odpowiadających dwóm różnym zakresom dokonywanej zmiany.

| Punkt wejścia | Zakres dostępnej zmiany | Źródło |
|---|---|---|
| Strefa 3 strony głównej (listwa ustawień) | Pełny zakres okna konfiguracji, warstwa domyślna i warstwa sesji, wszystkie szesnaście zakresów ustawień | Koncepcja platformy, rozdz. 7.3 |
| Uproszczone menu kontekstowe okna operacyjnego | Szybka zmiana konfiguracji bieżącej sesji — podzbiór ustawień właściwy warstwie sesji | Architektura, rozdz. 13 |

```
DWA PUNKTY WEJŚCIA DO TEGO SAMEGO MODELU KONFIGURACJI

  Strona główna · strefa 3 (listwa ustawień)      Okno operacyjne · menu kontekstowe
              │                                                │
              ▼                                                ▼
     pełne okno konfiguracji                       szybka zmiana bieżącej sesji
     warstwa domyślna + warstwa sesji              tylko warstwa sesji
     wszystkie szesnaście zakresów                 podzbiór ustawień
              │                                                │
              └───────────────►  jeden model konfiguracji  ◄───┘
                                 (nie dwa równoległe mechanizmy)
```

Wejście przez stronę główną otwiera okno konfiguracji w pełnym zakresie i pozwala pracować zarówno na warstwie domyślnej, jak i na warstwie sesji bieżącej (rozdz. 4.3). Uproszczone menu kontekstowe, dostępne bezpośrednio z okna operacyjnego bez opuszczania trwającej sesji, udostępnia skrócony dostęp do tych samych ustawień, ograniczony do warstwy sesji — jest wygodnym skrótem do tego samego modelu konfiguracji, a nie odrębnym, równoległym mechanizmem.

### 3.2. Układ okna i nawigacja wewnętrzna

Okno konfiguracji dzieli się na lewą kolumnę nawigacji, wskazującą jeden z szesnastu zakresów ustawień opisanych w rozdziale 5, oraz prawą kolumnę zawartości, prezentujący pozycje ustawień właściwe wybranemu zakresowi wraz z bieżącą wartością, poziomem warstwy i objaśnieniem kontekstowym. Nawigacja odzwierciedla wprost podział przyjęty w rozdziale 5 niniejszego dokumentu.

```
┌─────────────────────────┬──────────────────────────────────────────────┐
│ NAWIGACJA ZAKRESÓW       │ ZAWARTOŚĆ ZAKRESU                             │
│                          │                                                │
│  Aplikacja               │  [ nazwa ustawienia ]            [?]         │
│  Procesy                 │      wartość bieżąca · poziom warstwy         │
│  Akcje                   │                                                │
│  Zachowanie modeli        │  [ nazwa ustawienia ]            [?]         │
│  Tożsamość modeli         │      wartość bieżąca · poziom warstwy         │
│  Prompty systemowe        │                                                │
│  Rozszerzenia             │  …                                             │
│  Integracje               │                                                │
│  Historia                │  ┌──────────────────────────────────────┐     │
│  Pamięć                  │  │ Warstwa:  ( ) globalna  ( ) sesji     │     │
│  Izolacja  → rozdz. 6     │  └──────────────────────────────────────┘     │
│  Karty sesji              │                                                │
│  Komponenty własne        │                                                │
│  Chat Window              │                                                │
│  Execution Loop Window    │                                                │
│  Warstwy widoczności      │                                                │
│           → rozdz. 7      │                                                │
└─────────────────────────┴──────────────────────────────────────────────┘
```

Pozycja „Izolacja” w nawigacji zakresów otwiera dedykowane, trzypanelowe okno konfiguracji punktów izolacji opisane w rozdziale 6 — jest ono osobno rozwiniętym zakresem ze względu na własną złożoność (dwa rodzaje izolacji, siedem poziomów zasięgu, profile), nie zaś odrębnym mechanizmem poza modelem opisanym w niniejszym rozdziale.

### 3.3. Mechanizm objaśnień kontekstowych

Każda pozycja ustawienia w oknie konfiguracji — niezależnie od zakresu, do którego należy — jest opatrzona objaśnieniem kontekstowym, oznaczonym symbolem `[?]`, zgodnie z zasadą przyjętą w rozdziale 13 Architektury i powtórzoną wprost dla macierzy izolacji w rozdziale 6.4 Koncepcji platformy. Objaśnienie nie jest elementem interfejsu dodawanym niezależnie od ustawienia — jest częścią jego definicji, przechowywaną razem z pozostałymi atrybutami ustawienia (identyfikator, klucz, wartość, poziom warstwy), zgodnie z modelem danych ustawienia opisanym w rozdziale 12.1 dokumentu Model danych.

Treść objaśnienia odpowiada na dwa pytania: co robi dane ustawienie i jaki ma wpływ na działanie aplikacji. Poniższa tabela pokazuje przykładową treść objaśnienia dla trzech ustawień z różnych zakresów, ilustrując jednolity wzorzec stosowany w całym oknie konfiguracji.

| Ustawienie | Zakres | Przykładowa treść objaśnienia `[?]` |
|---|---|---|
| Retencja historii | Historia (5.9) | „Określa, jak długo platforma przechowuje historię rozmów tej sesji. Brak ustawienia oznacza przechowywanie bez limitu. Skrócenie retencji nie usuwa historii natychmiast — usunięcie następuje po upływie wskazanego okresu.” |
| Model procesu | Izolacja techniczna (6.3) | „Określa, czy proces sesji korzysta z odrębnego modelu wykonawczego, czy współdzieli go z innymi sesjami. Włączenie odrębnego modelu procesu zwiększa izolację kosztem współdzielenia zasobów obliczeniowych między sesjami.” |
| Tryb pracy w tle | Akcje (5.3) | „Określa, czy zadanie wykonuje się w tle, pozwalając kontynuować pracę w innej karcie sesji, czy blokuje bieżącą kartę do czasu zakończenia. Ustawienie dotyczy pojedynczego zadania i nie zmienia zachowania innych zadań w kolejce.” |

Mechanizm objaśnień kontekstowych jest identyczny dla wszystkich szesnastu zakresów ustawień — różni się jedynie treścią, nie formą prezentacji ani miejscem w interfejsie.

---

## 4. Warstwowość konfiguracji

### 4.1. Cztery warstwy ogólne

Konfiguracja platformy ma strukturę warstwową, ustanowioną w rozdziale 13 Architektury i powtórzoną w rozdziale 12.2 Modelu danych: wartość nakłada się od poziomu najszerszego do najwęższego.

| Warstwa | Zakres obowiązywania | Miejsce ustalania |
|---|---|---|
| Globalna | Cała platforma; wartość bazowa dla wszystkich pozostałych warstw | Okno konfiguracji, strefa 3 strony głównej |
| Środowisko | Jedno z czterech środowisk (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) | Okno konfiguracji |
| Projekt | Jeden projekt prowadzony w module Workspace | Okno konfiguracji lub okno modułu Workspace |
| Sesja | Jedna karta sesji | Okno konfiguracji lub uproszczone menu kontekstowe okna operacyjnego (rozdz. 3.1) |

Poniższy schemat przedstawia cztery warstwy ogólne jako stos od najszerszej do najwęższej, wraz z kierunkiem dziedziczenia wartości.

```
WARSTWOWOŚĆ KONFIGURACJI — wartość nakłada się od najszerszej do najwęższej

  [ WARSTWA GLOBALNA ]    cała platforma; wartość bazowa dla pozostałych warstw
        │                 ustalana w oknie konfiguracji (listwa ustawień strony głównej)
        ▼  dziedziczy z niej
  [ WARSTWA ŚRODOWISKA ]  TalkIn · WorkSpace · CodeStudio · MultitaskingAI
        │
        ▼  dziedziczy z niej
  [ WARSTWA PROJEKTU ]    jeden projekt prowadzony w module Workspace
        │
        ▼  dziedziczy z niej
  [ WARSTWA SESJI ]       jedna karta sesji — warstwa najwęższa; dostępna także
                          z uproszczonego menu kontekstowego okna operacyjnego (rozdz. 3.1)

  Zasada nakładania: wartość warstwy węższej nakłada się na wartość warstwy szerszej,
  nie zmieniając jej. Brak ustawienia na warstwie oznacza dziedziczenie z warstwy szerszej.
```

Kolejność warstw odpowiada kolejności nakładania się wartości: warstwa sesji dziedziczy z warstwy projektu, projektu — ze środowiska, środowiska — z warstwy globalnej, która stanowi bazę. Ta czterostopniowa warstwowość dotyczy ogólnych ustawień aplikacji, procesów, akcji, modeli, promptów, rozszerzeń, integracji, historii, pamięci, kart sesji, komponentów własnych, obu kanałów komunikacji operacyjnej oraz warstw widoczności funkcji, opisanych w rozdziale 5. Izolacja, ze względu na własną specyfikę wymagającą rozróżnienia dodatkowych poziomów (moduł, para modułów, karta sesji jako zasięg jednorazowy, rola środowiska MultitaskingAI), korzysta z rozszerzonej, siedmiopoziomowej hierarchii zasięgów opisanej w rozdziale 6.2 — hierarchia ta nie zastępuje czterech warstw ogólnych, lecz doprecyzowuje je dla przypadku izolacji. Odpowiedniość obu ujęć przedstawia poniższa tabela.

| Cztery warstwy ogólne (rozdz. 4.1) | Odpowiadające poziomy zasięgu izolacji (rozdz. 6.2) |
|---|---|
| Globalna | Globalny (domyślny) |
| Środowisko | Środowisko |
| — (poziomy pośrednie właściwe wyłącznie izolacji) | Moduł; Para modułów (relacja) |
| Projekt | Projekt |
| Sesja | Karta sesji |
| — (poziom właściwy wyłącznie izolacji) | Rola (MultitaskingAI) |

### 4.2. Reguła dziedziczenia i „brak ustawienia = wartość domyślna”

Ustawienie nieskonfigurowane na danej warstwie dziedziczy wartość z warstwy bezpośrednio szerszej, aż do warstwy globalnej, która — jeżeli i ona pozostaje nieskonfigurowana — dostarcza wartość domyślną wbudowaną w platformę. Poniższy schemat przedstawia rozstrzyganie wartości od warstwy najwęższej do najszerszej: obowiązuje pierwsza warstwa, na której ustawienie zostało określone.

```
ROZSTRZYGNIĘCIE WARTOŚCI USTAWIENIA DLA BIEŻĄCEJ SESJI
(od warstwy najwęższej do najszerszej — pierwsza ustawiona warstwa wygrywa)

  ustawiono na warstwie SESJI?       ──tak──►  użyj wartości warstwy sesji
        │ nie
        ▼
  ustawiono na warstwie PROJEKTU?    ──tak──►  użyj wartości warstwy projektu
        │ nie
        ▼
  ustawiono na warstwie ŚRODOWISKA?  ──tak──►  użyj wartości warstwy środowiska
        │ nie
        ▼
  ustawiono na warstwie GLOBALNEJ?   ──tak──►  użyj wartości warstwy globalnej
        │ nie
        ▼
  wartość domyślna wbudowana w platformę
```

Żaden poziom warstwy nie musi być skonfigurowany, a poziom nieskonfigurowany nie blokuje pracy: przejmuje wartość odziedziczoną. Reguła ta jest bezpośrednim zastosowaniem zasady nadrzędnej „brak ustawienia = wartość domyślna”, przywołanej wprost w rozdziale 1 Modelu danych oraz w rozdziale 6.5 Koncepcji platformy dla przypadku izolacji. Pierwszeństwo poszczególnych warstw i zachowanie przy braku ustawienia zestawia tabela dziedziczenia i pierwszeństwa.

| Warstwa | Pierwszeństwo | Dziedziczy z | Jeżeli nieskonfigurowana | Zakres zmiany dokonanej na tej warstwie |
|---|---|---|---|---|
| Sesja | Najwyższe | warstwa projektu | przejmuje wartość z warstwy projektu | tylko bieżąca karta sesji |
| Projekt | wyższe niż środowisko | warstwa środowiska | przejmuje wartość ze środowiska | jeden projekt modułu Workspace |
| Środowisko | wyższe niż globalna | warstwa globalna | przejmuje wartość globalną | jedno z czterech środowisk |
| Globalna | Najniższe (warstwa bazowa) | wartość domyślna platformy | dostarcza wartość domyślną wbudowaną w platformę | cała platforma |

Zmiana dokonana na warstwie węższej nie modyfikuje wartości warstwy szerszej — nakłada się na nią. Poniższy przykład ilustruje regułę na ustawieniu retencji historii (5.9): ustawienie retencji dla pojedynczej sesji nie zmienia globalnej wartości retencji obowiązującej pozostałe sesje.

| Warstwa | Wartość retencji historii w przykładzie | Co faktycznie obowiązuje |
|---|---|---|
| Globalna | bez limitu (wartość domyślna, rozdz. 5.9) | wszystkie sesje bez własnego ustawienia |
| Środowisko | brak ustawienia | dziedziczy „bez limitu” z warstwy globalnej |
| Projekt | brak ustawienia | dziedziczy „bez limitu” z warstwy szerszej |
| Sesja bieżąca | 30 dni (ustawienie przykładowe) | wyłącznie ta jedna karta sesji |

Zmiana dotyczy wyłącznie tej jednej sesji, dla której została dokonana; wartość globalna oraz konfiguracja pozostałych sesji pozostają nienaruszone.

### 4.3. Warstwa sesji i szybka zmiana konfiguracji z okna operacyjnego

Warstwa sesji jest jedyną warstwą dostępną również poza pełnym oknem konfiguracji — z uproszczonego menu kontekstowego okna operacyjnego (rozdz. 3.1, Architektura rozdz. 13). Zmiany dokonane tą drogą obejmują wyłącznie bieżącą kartę sesji i nie wymagają otwarcia pełnego okna konfiguracji, co pozwala na szybką korektę zachowania modelu, izolacji lub pamięci w trakcie trwającej pracy, bez przerywania jej przejściem do innego okna platformy.

| Cecha | Warstwa domyślna (globalna, środowisko, projekt) | Warstwa sesji |
|---|---|---|
| Miejsce ustalania | wyłącznie pełne okno konfiguracji | pełne okno konfiguracji lub uproszczone menu okna operacyjnego |
| Zasięg zmiany | wszystkie przyszłe sesje objęte warstwą | wyłącznie bieżąca karta sesji |
| Wpływ na warstwę domyślną | ustala wartość bazową | nakłada się na wartość domyślną, nie zmieniając jej |
| Typowe zastosowanie | trwała polityka pracy | jednorazowa korekta w trakcie trwającej pracy |

---

## 5. Zakresy ustawień

Poniższe siedemnaście zakresów wyczerpuje pełną strukturę okna konfiguracji: jedenaście zakresów wskazanych w rozdziale 13 Architektury (aplikacja, procesy, akcje, zachowanie modeli, tożsamość modeli, prompty systemowe, rozszerzenia, integracje, historia, pamięć, izolacja) oraz sześć zakresów wynikających wprost z mechanizmów platformy (karty sesji, komponenty własne, Chat Window, Execution Loop Window, warstwy widoczności funkcji, Always On Display). Każdy zakres jest opisany celem, wykazem pozycji ustawień, zastosowaniem warstwowym oraz wartością domyślną. Pełne, zbiorcze zestawienie wszystkich pozycji ustawień zawiera Załącznik A.

Katalog siedemnastu zakresów w jednej tabeli, wraz z głównym przedmiotem konfiguracji i źródłem mechanizmu:

| Zakres ustawień | Co konfiguruje | Główne źródło mechanizmu |
|---|---|---|
| 5.1 Aplikacja | Platforma jako całość: motyw wizualny, język interfejsu, metody uwierzytelniania, urządzenia połączone, powiadomienia | Architektura, rozdz. 13, 15 |
| 5.2 Procesy | Osiem zakresów izolacji technicznej procesu sesji | Architektura, rozdz. 12; niniejszy dokument, rozdz. 6 |
| 5.3 Akcje | Silnik kolejek i zadania modułu Automations oraz orkiestracja środowiska MultitaskingAI | Koncepcja platformy, rozdz. 11.3, 13.4–13.6 |
| 5.4 Zachowanie modeli | Kanał połączenia z modelem, dobór modelu, parametry wywołania, współpraca wykonawców MultitaskingAI | Architektura, rozdz. 9; Koncepcja platformy, rozdz. 13.3 |
| 5.5 Tożsamość modeli | Tożsamość agentów i wcielenia ról MultitaskingAI, uprawnienia | Koncepcja platformy, rozdz. 11.15, 13.3 |
| 5.6 Prompty systemowe | Instrukcje systemowe projektu i agenta, budowa promptów przez rolę Coordinator | Koncepcja platformy, rozdz. 11.2, 11.15, 13.3 |
| 5.7 Rozszerzenia | Wtyczki, umiejętności, konektory i serwery MCP (Danaco Plugin oraz Personal) | Architektura, rozdz. 14 |
| 5.8 Integracje | Jawne powiązania między środowiskami, modułami i komponentami własnymi | Koncepcja platformy, rozdz. 6.2, 11 |
| 5.9 Historia | Przechowywanie i dostępność historii rozmów i sesji | Model danych, rozdz. 4 |
| 5.10 Pamięć | Zasoby pamięci kontekstowej na czterech poziomach | Model danych, rozdz. 6 |
| 5.11 Izolacja | Izolacja kontekstu i izolacja techniczna procesu sesji (okno trzypanelowe) | niniejszy dokument, rozdz. 6; Architektura, rozdz. 12 |
| 5.12 Karty sesji | Ustawienia pojedynczej karty sesji jako jednostki równoległej pracy | Koncepcja platformy, rozdz. 8 |
| 5.13 Komponenty własne | Tworzenie, zapis i przypisanie automatyki, agenta, projektu i profilu asystenta | Koncepcja platformy, rozdz. 2.3, 7.2 |
| 5.14 Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca: szerokość kolumny, zachowanie strumienia, zatwierdzanie i przerywanie działań, zakres pamięci | Architektura, rozdz. 13; Koncepcja platformy, rozdz. 6.1 |
| 5.15 Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca: równoległość zadań, polityka ponowień, progi kontroli jakości, zakres autonomii, punkty zatwierdzeń, sterowanie przebiegiem | Koncepcja platformy, rozdz. 13.3–13.6 |
| 5.16 Warstwy widoczności funkcji | Przypisanie funkcji do warstw 1–4, profile warstw wg roli użytkownika, wyszukiwarka funkcji, skróty klawiszowe, tryb administracyjny | niniejszy dokument, rozdz. 7 |
| 5.17 Always On Display | Obecność i tryb funkcji globalnej Always On Display, klasy zdarzeń wyzwalających sugestie, progi i częstotliwość, reguły wyciszania, tor głosowy, tryb wobec procesu | `funkcje-globalne/always-on-display.md`, rozdz. 3, 10 |

### 5.1. Aplikacja

**Cel:** ustawienia dotyczące samej platformy jako całości, niezależne od konkretnego środowiska, modułu lub sesji.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Motyw wizualny | Wariant systemu wizualnego marki (granat i złoto) stosowany w oknie klienta | globalna | wariant podstawowy |
| Język interfejsu | Język wyświetlania interfejsu platformy | globalna | polski |
| Metody uwierzytelniania | Aktywne metody logowania: hasło, PIN, Windows Hello, e-mail uwierzytelniający (Architektura, rozdz. 15) | globalna | hasło i e-mail uwierzytelniający |
| Urządzenia połączone | Wykaz urządzeń powiązanych z kontem oraz zarządzanie ich dostępem | globalna | brak ograniczeń liczby urządzeń |
| Powiadomienia | Zakres i kanał powiadomień o zdarzeniach platformy, w tym zdarzeniach nadzorowanych przez funkcję Mobile (Koncepcja platformy, rozdz. 12.1) | globalna, sesja | powiadomienia aktywne |

### 5.2. Procesy

**Cel:** ustawienia procesu sesji po stronie serwera — zakresy izolacji technicznej ustanowione w rozdziale 12 Architektury. Pełny mechanizm konfiguracji tych zakresów, wraz z siedmiopoziomową hierarchią zasięgów, opisano w rozdziale 6 niniejszego dokumentu; poniższa tabela wskazuje wyłącznie same pozycje ustawień, będące przedmiotem tej konfiguracji.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Katalog roboczy sesji | Izolacja katalogu plikowego procesu sesji | rozdz. 6 (7 poziomów) | wyłączony |
| Środowisko procesu | Izolacja zmiennych środowiskowych i kontekstu wykonania procesu | rozdz. 6 (7 poziomów) | wyłączony |
| Katalog danych i konfiguracji modelu | Izolacja katalogu przechowującego dane i konfigurację modelu | rozdz. 6 (7 poziomów) | wyłączony |
| Dostęp sieciowy | Izolacja dostępu procesu sesji do sieci | rozdz. 6 (7 poziomów) | wyłączony |
| Zakres odczytu i zapisu plików | Ograniczenie zakresu plików dostępnych procesowi sesji | rozdz. 6 (7 poziomów) | wyłączony |
| Konto i token per sesja | Odrębne konto i token dostępowy dla procesu sesji | rozdz. 6 (7 poziomów) | wyłączony |
| Model procesu | Odrębny albo współdzielony model wykonawczy procesu | rozdz. 6 (7 poziomów) | wyłączony (współdzielony) |
| Serwer wykonania | Przypisanie procesu sesji do wskazanego serwera wykonawczego | rozdz. 6 (7 poziomów) | wyłączony |

### 5.3. Akcje

**Cel:** ustawienia silnika kolejek i zadań — jednostek pracy definiowanych w module Automations (Koncepcja platformy, rozdz. 11.3) oraz w warstwie orkiestracji środowiska MultitaskingAI (rozdz. 11.4–11.5).

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Akcje kolejki dostępne | Zestaw akcji: enqueue, dequeue, delay, retry, pause, resume, split, merge, route, branch, condition | globalna, moduł | pełny zestaw jedenastu akcji |
| Tryb pracy w tle zadania | Wykonywanie zadania w tle bez blokowania bieżącej karty sesji | sesja, zadanie | tryb w tle aktywny |
| Zasięg kolejki | Poziom, na którym zdefiniowana jest kolejka: globalna, lokalna, modelu, agenta, projektu | moduł, projekt | lokalna |
| Obsługa błędów kolejki | Zachowanie przy błędzie wykonania zadania: retry, route, pause | zadanie | retry |
| Powiązanie z silnikiem kolejek MultitaskingAI | Spięcie silnika kolejek modułu Automations z silnikiem kolejek środowiska MultitaskingAI (Koncepcja platformy, rozdz. 13.6) | globalna, środowisko | wyłączone |

### 5.4. Zachowanie modeli

**Cel:** ustawienia sposobu, w jaki modele AI są wywoływane i jak reagują na polecenia — kanał połączenia, dobór modelu oraz parametry wywołania.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Kanał modelu | Sposób połączenia z modelem: API, CLI, SSH, HTTP (Architektura, rozdz. 9) | sesja, rola | API |
| Model bazowy | Wybór konkretnego modelu udostępnianego przez kanał | sesja, rola, agent | zależna od kanału |
| Relacje między wykonawcami (MultitaskingAI) | Tryb współpracy Executor 1 i Executor 2: praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna (Koncepcja platformy, rozdz. 13.3) | rola | praca niezależna |
| Liczba podagentów Subagent Network | Liczba jednocześnie uruchomionych podagentów wykonawcy, do piętnastu (Koncepcja platformy, rozdz. 13.3) | rola | 0 (nieaktywny) |
| Dane dostępowe kanału | Odwołanie do danych dostępowych kanału modelu, przechowywanych poza bazą danych (Model danych, rozdz. 8.1) | globalna, sesja, rola | brak |

### 5.5. Tożsamość modeli

**Cel:** ustawienia definiujące, kim — w sensie tożsamości operacyjnej — jest dany model lub agent w kontekście pracy. Zakres ten obejmuje zarówno tożsamość agentów tworzonych w module Agents, jak i wcielenia ról środowiska MultitaskingAI.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Nazwa i tożsamość agenta | Nazwa własna i opis tożsamości agenta (Koncepcja platformy, rozdz. 11.15) | komponent własny | brak (wymagane przy tworzeniu) |
| Model bazowy agenta | Model AI leżący u podstaw agenta | komponent własny | brak (wymagane przy tworzeniu) |
| Wcielenie roli Executor 3 / Validator | Konkretne wcielenie czwartej roli MultitaskingAI: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator (Koncepcja platformy, rozdz. 13.3) | rola | Validator |
| Uprawnienia agenta | Zakres działań dozwolonych agentowi, ustalany w Permissions Center (Koncepcja platformy, rozdz. 11.15) | komponent własny | zakres minimalny, rozszerzany świadomie |

### 5.6. Prompty systemowe

**Cel:** ustawienia instrukcji systemowych kierujących zachowaniem modeli — na poziomie projektu, agenta i procesu budowy promptów przez rolę Coordinator.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Instrukcje systemowe projektu | Zestaw instrukcji właściwych projektowi w module Workspace, ustalany w Instructions Panel (Koncepcja platformy, rozdz. 11.2) | projekt | brak |
| Instrukcje systemowe agenta | Zestaw instrukcji definiujących zachowanie agenta, ustalany w Agent Builder (Koncepcja platformy, rozdz. 11.15) | komponent własny | brak |
| Budowa promptów przez Coordinatora | Dynamiczne tworzenie i optymalizacja promptów dla wykonawców przez rolę Coordinator (Koncepcja platformy, rozdz. 13.3) | rola | aktywna, gdy rola Coordinator jest przypisana |
| Zakres poleceń profilu asystenta | Zakres poleceń i obsługiwanych działań profilu asystenta modułu Assistant (Koncepcja platformy, rozdz. 11.10) | komponent własny | zakres podstawowy |

### 5.7. Rozszerzenia

**Cel:** ustawienia wtyczek, umiejętności, konektorów i serwerów MCP realizujących wspólny kontrakt integracji, niezależnie od źródła pochodzenia (Architektura, rozdz. 14).

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Wykaz rozszerzeń | Lista dostępnych rozszerzeń wraz ze źródłem: Danaco Plugin (wbudowane) albo Personal (zainstalowane przez Operatora) | globalna | zestaw wbudowany Danaco Plugin |
| Stan włączenia rozszerzenia | Włączenie albo wyłączenie pojedynczego rozszerzenia | globalna, moduł | zależna od rozszerzenia |
| Instalacja rozszerzenia Personal | Przesłanie rozszerzenia z urządzenia Operatora na serwer, do katalogu użytkownika | globalna | — |
| Skille agenta | Umiejętności przypisane agentowi w Skills Manager (Koncepcja platformy, rozdz. 11.15) | komponent własny | brak |

### 5.8. Integracje

**Cel:** jawne, konfigurowalne powiązania między środowiskami, modułami, komponentami własnymi i funkcjami platformy — mechanizm opisany w rozdziale 6.2 Koncepcji platformy i realizowany w podpunktach „Powiązania” rozdziału 9. Integracja różni się od rozszerzenia (5.7) tym, że nie dotyczy dołączenia nowego narzędzia z zewnątrz, lecz połączenia dwóch elementów już obecnych na platformie.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Powiązanie pary modułów | Współdzielenie mechanizmu między dwoma modułami, na przykład operacji kontekstowych AI między Studio a Translate (Koncepcja platformy, rozdz. 11.1, 11.7) | para modułów | brak powiązania |
| Powiązanie komponentu własnego z modułem | Spięcie automatyki, agenta lub projektu z modułem, w którym ma być wykorzystywany | komponent własny | brak powiązania |
| Integracja MultitaskingAI z Automations | Połączenie orkiestracji środowiska MultitaskingAI z harmonogramami i kolejkami modułu Automations (Koncepcja platformy, rozdz. 13.6) | globalna, środowisko | wyłączona |
| Konektory agenta | Integracje zewnętrzne podłączone do agenta w Connectors Manager (Koncepcja platformy, rozdz. 11.15) | komponent własny | brak |

### 5.9. Historia

**Cel:** ustawienia przechowywania i dostępności historii rozmów i sesji.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Retencja historii | Okres przechowywania historii sesji i rozmów (Model danych, rozdz. 4.3) | globalna, środowisko, projekt, sesja | bez limitu |
| Zapis do powrotu sesji | Zachowanie zapisywania zamkniętej sesji, umożliwiające jej wznowienie (Model danych, rozdz. 4.1) | globalna, sesja | zapis aktywny |
| Współdzielenie historii między kartami | Zakres, w jakim historia jest współdzielona między kartami sesji, modułami lub środowiskami (Koncepcja platformy, rozdz. 6.1) | rozdz. 6 (poziomy zasięgu izolacji) | odrębna dla każdej karty |
| Usunięcie historii | Usunięcie wskazanego wpisu lub całości historii sesji | sesja | — |

### 5.10. Pamięć

**Cel:** ustawienia zasobów pamięci kontekstowej wykorzystywanej przez AI, na poziomach globalnym, środowiska, projektu i sesji (Model danych, rozdz. 6).

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Poziom pamięci | Poziom, na którym zasób pamięci obowiązuje: globalna, środowisko, projekt, sesja | globalna, środowisko, projekt, sesja | sesja |
| Stan włączenia pamięci | Włączenie albo wyłączenie pamięci na danym poziomie | wszystkie cztery warstwy | włączona |
| Odłączenie pamięci w sesji | Odłączenie dostępu do pamięci dla pojedynczej sesji, możliwe przed pierwszym promptem (Model danych, rozdz. 6.2) | sesja | dołączona |
| Zawartość zasobu pamięci | Treść pojedynczego wpisu pamięci wraz ze znacznikiem czasu | wszystkie cztery warstwy | — |

### 5.11. Izolacja

**Cel:** zakres izolacji kontekstu (historia, pamięć, kontekst) i izolacji technicznej procesu sesji (osiem zakresów z rozdziału 12 Architektury), konfigurowany z poziomu dedykowanego, trzypanelowego okna. Ze względu na złożoność — dwa rodzaje izolacji, siedem poziomów zasięgu, profile izolacji — zakres ten opisano w całości w rozdziale 6 niniejszego dokumentu.

### 5.12. Karty sesji

**Cel:** ustawienia właściwe pojedynczej karcie sesji jako jednostce równoległej pracy w oknie środowiska (Koncepcja platformy, rozdz. 8).

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Współdzielenie kontekstu między kartami | Zakres, w jakim karty sesji współdzielą historię, pamięć i kontekst (Koncepcja platformy, rozdz. 8) | rozdz. 6 (poziom „karta sesji”) | odrębny dla każdej karty |
| Migawka układu okien | Zapisany układ okien operacyjnych właściwy karcie sesji (Model danych, rozdz. 4.1) | sesja | układ domyślny modułu |
| Przypisanie projektu do karty | Powiązanie karty sesji z projektem modułu Workspace | sesja | brak przypisania |
| Szybka zmiana konfiguracji | Dostępność uproszczonego menu kontekstowego dla danej karty (rozdz. 3.1) | sesja | dostępne |

### 5.13. Komponenty własne

**Cel:** ustawienia tworzenia, zapisu i przypisania komponentów własnych — automatyki, agenta, projektu i profilu asystenta (Koncepcja platformy, rozdz. 2.3 i rozdz. 7.2). Szablon pól komponentu własnego, wraz z polami swoistymi dla każdego rodzaju, zawiera Załącznik D Koncepcji platformy oraz Załącznik B niniejszego dokumentu.

| Ustawienie | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|
| Nazwa komponentu | Nazwa własna komponentu, nadawana przy tworzeniu | komponent własny | brak (wymagane przy tworzeniu) |
| Miejsce zastosowania | Moduł lub moduły, w których komponent jest wykorzystywany operacyjnie | komponent własny | brak |
| Widoczność komponentu | Zasięg dostępności komponentu: globalny albo projektowy | komponent własny | globalny |
| Powiązania komponentu | Jawne, konfigurowalne powiązania z innymi komponentami lub modułami (rozdz. 5.8) | komponent własny | brak |

### 5.14. Chat Window

**Cel:** ustawienia głównego okna komunikacji między Użytkownikiem a Wykonawcą — okna zajmującego lewą kolumnę obszaru roboczego, stałą i o pełnej wysokości, obecnego w tym samym miejscu układu w każdym module i w każdym środowisku platformy. Okno przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Ustawienia tego zakresu sterują szerokością kolumny, zachowaniem strumienia, trybem zatwierdzania i przerywania działań oraz zakresem pamięci okna.

| Ustawienie | Klucz | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|---|
| Szerokość kolumny Chat Window | `chat.column.width` | Szerokość lewej kolumny obszaru roboczego zajmowanej przez okno komunikacji, wyrażona w procentach szerokości obszaru roboczego; regulacji podlega wyłącznie szerokość — położenie kolumny pozostaje stałe | globalna, środowisko, projekt, sesja | 30% szerokości obszaru roboczego |
| Zakres regulacji szerokości | `chat.column.width.range` | Dolna i górna granica dopuszczalnej szerokości kolumny okna komunikacji | globalna, środowisko | od 20% do 50% |
| Tryb strumienia odpowiedzi | `chat.stream.mode` | Sposób prezentacji odpowiedzi Wykonawcy: przyrostowy strumień znaków albo prezentacja całości po zakończeniu generowania | globalna, środowisko, sesja | strumień przyrostowy |
| Przewijanie strumienia | `chat.stream.autoscroll` | Podążanie widoku za napływającą treścią strumienia albo utrzymanie pozycji wskazanej przez Użytkownika | globalna, sesja | podążanie za strumieniem |
| Widoczność przebiegu pośredniego | `chat.stream.steps` | Prezentowanie w strumieniu kroków pośrednich Wykonawcy — wywołań narzędzi, odczytów i zapisów — obok wyniku końcowego | globalna, środowisko, sesja | kroki pośrednie widoczne |
| Zakres zatwierdzania działań | `chat.approval.scope` | Klasa działań Wykonawcy wymagająca jawnego zatwierdzenia przez Użytkownika przed wykonaniem: wszystkie działania, działania nieodwracalne, działania poza katalogiem roboczym, brak zatwierdzania | globalna, środowisko, projekt, sesja | działania nieodwracalne |
| Zasięg udzielonego zatwierdzenia | `chat.approval.remember` | Zasięg obowiązywania jednorazowo udzielonego zatwierdzenia: pojedyncze działanie, bieżąca karta sesji, bieżący projekt | projekt, sesja | pojedyncze działanie |
| Przerwanie działania | `chat.interrupt.mode` | Zachowanie platformy po przerwaniu działania przez Użytkownika: zatrzymanie natychmiastowe z zachowaniem wyniku cząstkowego albo zatrzymanie po zakończeniu bieżącego kroku | globalna, sesja | zatrzymanie natychmiastowe z zachowaniem wyniku cząstkowego |
| Zakres pamięci okna komunikacji | `chat.memory.scope` | Poziom, z którego okno komunikacji czerpie kontekst i na którym go utrwala: sesja, projekt, środowisko, globalny (rozdz. 5.10) | globalna, środowisko, projekt, sesja | sesja |
| Głębokość kontekstu rozmowy | `chat.memory.depth` | Liczba wcześniejszych wymian rozmowy przekazywanych Wykonawcy wraz z bieżącym poleceniem | globalna, środowisko, projekt, sesja | pełna historia bieżącej karty sesji |
| Odłączenie pamięci okna komunikacji | `chat.memory.detach` | Prowadzenie rozmowy bez dostępu do zasobów pamięci i bez utrwalania kontekstu (Model danych, rozdz. 6.2) | sesja | pamięć dołączona |

### 5.15. Execution Loop Window

**Cel:** ustawienia okna pętli wykonawczej prezentującego komunikację między Koordynatorem a Wykonawcą — okna otwieranego jako kolumna sąsiadująca z Chat Window. Okno prowadzi pętlę wykonawczą, koordynuje zadania, nadzoruje realizację i kontroluje przebieg procesów; zawiera bieżące zlecenie wraz z dekompozycją na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem. Ustawienia tego zakresu sterują równoległością zadań, polityką ponowień, progami kontroli jakości, zakresem autonomii Koordynatora, punktami zatwierdzeń oraz dostępnością operacji sterowania przebiegiem.

| Ustawienie | Klucz | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|---|
| Otwarcie okna pętli wykonawczej | `loop.window.open` | Otwarcie kolumny sąsiadującej z Chat Window przy rozpoczęciu zlecenia wielozadaniowego albo otwarcie na żądanie Użytkownika | globalna, środowisko, sesja | otwarcie przy rozpoczęciu zlecenia wielozadaniowego |
| Szerokość kolumny Execution Loop Window | `loop.column.width` | Szerokość kolumny zajmowanej przez okno pętli wykonawczej, wyrażona w procentach szerokości obszaru roboczego | globalna, środowisko, sesja | 25% szerokości obszaru roboczego |
| Równoległość zadań | `loop.concurrency` | Liczba zadań realizowanych jednocześnie w ramach jednego zlecenia | globalna, środowisko, projekt, sesja, rola | 3 zadania |
| Równoległość wykonawców | `loop.executors.concurrency` | Liczba wykonawców pracujących jednocześnie nad zadaniami jednego zlecenia (Koncepcja platformy, rozdz. 13.3) | środowisko, rola | 2 wykonawców |
| Głębokość kolejki zadań | `loop.queue.depth` | Liczba zadań oczekujących utrzymywanych w kolejce pętli wykonawczej | globalna, środowisko, projekt | bez limitu |
| Liczba ponowień zadania | `loop.retry.count` | Liczba ponowień zadania zakończonego niepowodzeniem przed przekazaniem sprawy Użytkownikowi | globalna, środowisko, projekt, sesja, rola | 2 ponowienia |
| Strategia ponowień | `loop.retry.strategy` | Odstęp między ponowieniami: ponowienie natychmiastowe, odstęp stały, odstęp narastający | globalna, środowisko, rola | odstęp narastający |
| Zachowanie po wyczerpaniu ponowień | `loop.retry.exhausted` | Działanie po wyczerpaniu ponowień: wstrzymanie pętli i przekazanie decyzji Użytkownikowi, przekierowanie zadania do innego wykonawcy, oznaczenie zadania jako nieudanego i kontynuacja pozostałych | globalna, środowisko, projekt | wstrzymanie pętli i przekazanie decyzji Użytkownikowi |
| Próg przyjęcia wyniku | `loop.quality.threshold` | Poziom oceny kontroli jakości, od którego wynik zadania zostaje przyjęty bez ponowienia | globalna, środowisko, projekt, rola | ocena pozytywna roli kontrolnej |
| Rola kontrolna | `loop.quality.role` | Wcielenie czwartej roli MultitaskingAI wykonujące kontrolę jakości wyniku: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator (rozdz. 5.5) | środowisko, rola | Validator |
| Zakres kontroli jakości | `loop.quality.scope` | Zadania poddawane kontroli jakości: każde zadanie, zadania oznaczone jako krytyczne, wynik końcowy zlecenia | globalna, środowisko, projekt | każde zadanie |
| Zakres autonomii Koordynatora | `loop.autonomy.scope` | Zakres decyzji podejmowanych przez Koordynatora bez udziału Użytkownika: dekompozycja zlecenia, przydział zadań wykonawcom, ponowienia, zmiana kolejności zadań, korekta zlecenia | globalna, środowisko, projekt, sesja | dekompozycja zlecenia, przydział zadań i ponowienia |
| Granica autonomii operacyjnej | `loop.autonomy.limit` | Klasa operacji, której Koordynator nie wykonuje bez zatwierdzenia Użytkownika: operacje nieodwracalne, operacje poza katalogiem roboczym, operacje sieciowe | globalna, środowisko, projekt | operacje nieodwracalne |
| Punkty zatwierdzeń | `loop.approval.points` | Momenty pętli wykonawczej wymagające zatwierdzenia Użytkownika: przyjęcie dekompozycji zlecenia, rozpoczęcie realizacji, przejście między etapami zlecenia, przyjęcie wyniku końcowego | globalna, środowisko, projekt, sesja | przyjęcie dekompozycji zlecenia i przyjęcie wyniku końcowego |
| Oczekiwanie na zatwierdzenie | `loop.approval.timeout` | Zachowanie pętli w oczekiwaniu na zatwierdzenie Użytkownika: oczekiwanie bezterminowe albo wstrzymanie po upływie wskazanego okresu | globalna, środowisko, sesja | oczekiwanie bezterminowe |
| Sterowanie przebiegiem | `loop.control.actions` | Operacje sterowania przebiegiem dostępne w oknie pętli wykonawczej: wstrzymanie, wznowienie, przerwanie, korekta zlecenia | globalna, środowisko, sesja | pełny zestaw czterech operacji |
| Zachowanie po przerwaniu pętli | `loop.control.abort` | Stan zadań w toku po przerwaniu pętli: zatrzymanie natychmiastowe wszystkich zadań albo dokończenie zadań rozpoczętych | globalna, środowisko, sesja | zatrzymanie natychmiastowe wszystkich zadań |
| Szczegółowość komunikatów sterujących | `loop.messages.verbosity` | Zakres komunikatów wymienianych między Koordynatorem a Wykonawcą prezentowanych w oknie: pełna wymiana, komunikaty sterujące i wyniki, wyłącznie zmiany stanu zadań | globalna, środowisko, sesja | komunikaty sterujące i wyniki |
| Powiązanie pętli z silnikiem kolejek | `loop.queue.binding` | Spięcie pętli wykonawczej z silnikiem kolejek modułu Automations (rozdz. 5.3; Koncepcja platformy, rozdz. 13.6) | globalna, środowisko | wyłączone |

### 5.16. Warstwy widoczności funkcji

**Cel:** ustawienia stopniowego ujawniania funkcjonalności interfejsu — przypisania każdego elementu interfejsu do jednej z czterech warstw widoczności, profili warstw właściwych roli użytkownika oraz dróg dostępu do funkcji warstw wyższych: wyszukiwarki funkcji, skrótów klawiszowych, trybu administracyjnego i polecenia języka naturalnego. Pełny opis tej grupy ustawień zawiera rozdział 7.

| Ustawienie | Klucz | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|---|
| Przypisanie funkcji do warstwy | `visibility.function.layer` | Warstwa widoczności pojedynczej funkcji interfejsu: 1 — zawsze widoczna, 2 — widoczna na żądanie, 3 — rozwinięcie kontekstowe, 4 — funkcja ekspercka (rozdz. 7.1) | globalna, środowisko, moduł, rola użytkownika | zgodnie z katalogiem warstw platformy |
| Profil warstw roli użytkownika | `visibility.profile.role` | Nazwany zestaw przypisań funkcji do warstw obowiązujący daną rolę użytkownika (rozdz. 7.2) | globalna, środowisko, rola użytkownika | profil użytkownika podstawowego |
| Zwijanie elementów warstwy 2 | `visibility.layer2.autocollapse` | Samoczynne zwinięcie elementu warstwy 2 po jego użyciu | globalna, środowisko, sesja | zwijanie samoczynne |
| Postać rozwinięć warstwy 3 | `visibility.layer3.trigger` | Element wywołujący rozwinięcie kontekstowe: menu kebab, menu hamburger, menu kontekstowe, panel popover, lista rozwijana | globalna, środowisko, moduł | menu kebab |
| Znaczniki kontekstowe | `visibility.context.tags` | Zestaw znaczników prezentowanych w pasku kontekstu: środowisko, repozytorium, projekt, model, wykonawca | globalna, środowisko, moduł, sesja | pełny zestaw pięciu znaczników |
| Wyszukiwarka funkcji | `visibility.search.enabled` | Dostępność wyszukiwarki funkcji jako drogi dostępu do elementów warstw 3 i 4 (rozdz. 7.3) | globalna, środowisko, rola użytkownika | wyszukiwarka dostępna |
| Zakres wyszukiwarki funkcji | `visibility.search.scope` | Zbiór przeszukiwanych elementów: funkcje bieżącego modułu, funkcje wszystkich modułów, funkcje wraz z komponentami własnymi | globalna, środowisko, sesja | funkcje wszystkich modułów |
| Skrót wywołania wyszukiwarki | `visibility.search.shortcut` | Skrót klawiszowy otwierający wyszukiwarkę funkcji | globalna, rola użytkownika | wartość z katalogu skrótów platformy |
| Skróty klawiszowe funkcji | `visibility.shortcuts.map` | Przypisanie skrótów klawiszowych do funkcji platformy (rozdz. 7.3) | globalna, środowisko, moduł, rola użytkownika | katalog skrótów platformy |
| Tryb administracyjny | `visibility.admin.mode` | Ujawnienie funkcji warstwy 4 — trybów administracyjnych i narzędzi diagnostycznych niskiego poziomu (rozdz. 7.4) | globalna, środowisko, rola użytkownika | tryb nieaktywny |
| Zakres trybu administracyjnego | `visibility.admin.scope` | Zbiór funkcji ujawnianych po włączeniu trybu administracyjnego: narzędzia diagnostyczne, konfiguracja ról, operacje na procesach sesji | globalna, środowisko | narzędzia diagnostyczne |
| Dostęp poleceniem języka naturalnego | `visibility.nlcommand.enabled` | Wywoływanie funkcji warstw 3 i 4 poleceniem języka naturalnego w Chat Window (rozdz. 5.14, 7.3) | globalna, środowisko, rola użytkownika | dostęp aktywny |

### 5.17. Always On Display

**Cel:** ustawienia funkcji globalnej Always On Display — agenta towarzyszącego obecnego ponad wszystkimi środowiskami, modułami, projektami i sesjami. Pełny opis funkcji, w tym reguł wyzwalania sugestii i toru głosowego, zawiera `funkcje-globalne/always-on-display.md`.

| Ustawienie | Klucz | Opis | Warstwy | Wartość domyślna |
|---|---|---|---|---|
| Obecność awatara | `aod.presence.visible` | Widoczność pływającego awatara funkcji; ukrycie nie przerywa działania funkcji w tle | globalna, środowisko, sesja | awatar widoczny |
| Tryb obecności | `aod.presence.mode` | Tryb pełny, cichy albo ukryty (`funkcje-globalne/always-on-display.md`, rozdz. 9.2) | globalna, środowisko, sesja | pełny |
| Klasy zdarzeń wyzwalających | `aod.trigger.classes` | Zbiór czynnych klas zdarzeń tworzących sugestie: stan pętli wykonawczej, stan kolejki zadań, wynik kontroli jakości, zdarzenia modułów, harmonogram, kontekst pracy | globalna, środowisko, moduł | wszystkie klasy czynne |
| Progi wyzwalania | `aod.trigger.thresholds` | Progi czasu oczekiwania zadania, liczby ponowień, wypełnienia kolejki i powtarzalności czynności ręcznej | globalna, środowisko, moduł | wartości domyślne katalogu progów |
| Częstotliwość ujawniania | `aod.suggestion.rate` | Liczba sugestii ujawnianych samoczynnie w godzinie oraz odstęp między nimi | globalna, sesja | 3 na godzinę, odstęp 5 minut |
| Waga ujawniana samoczynnie | `aod.suggestion.minweight` | Najniższa waga sugestii otwierającej dymek kontekstowy bez interakcji | globalna, środowisko | wysoka |
| Reguły wyciszania | `aod.mute.rules` | Zakresy i czasy wyciszenia sugestii, wraz z wyjątkiem wagi krytycznej | globalna | komplet zakresów i czasów |
| Czas życia sugestii nieprzyjętej | `aod.suggestion.ttl` | Czas, po którym sugestia nieprzyjęta otrzymuje status odrzuconej | globalna | 24 godziny |
| Tor głosowy funkcji | `aod.voice.enabled` | Dostępność globalnego toru głosowego funkcji, globalnej frazy wybudzającej i syntezy mowy odpowiedzi | globalna, sesja | tor czynny |
| Tryb wobec procesu | `aod.process.mode` | Obserwator albo operator wobec procesów środowiska MultitaskingAI | globalna, projekt, sesja | obserwator |

---

## 6. Okno konfiguracji punktów izolacji

Niniejszy rozdział opisuje szczegółowy interfejs okna konfiguracji punktów izolacji, przyjęty w Koncepcji platformy (rozdz. 6.4–6.6), i osadza go w pełnej strukturze okna konfiguracji. Okno to łączy dwa rodzaje izolacji opisane dotąd w odrębnych źródłach — izolację kontekstu (historia, pamięć, kontekst; Koncepcja platformy, rozdz. 2.4, 6 i 8) oraz izolację techniczną procesu sesji (osiem zakresów z rozdziału 12 Architektury) — w jednym, spójnym miejscu, zgodnie z zasadą, że wszystkie zależności platformy konfiguruje się z jednego okna.

### 6.1. Trzy panele

Okno konfiguracji punktów izolacji otwiera się z pozycji „Izolacja” w nawigacji zakresów (rozdz. 3.2) i dzieli się na trzy panele.

```
┌────────────────┬──────────────────────────────┬─────────────────────┐
│ SELEKTOR       │ MACIERZ IZOLACJI             │ PROFIL I PODGLĄD    │
│ ZASIĘGU        │                              │                     │
│                │ Izolacja kontekstu:          │ [ Zapisz profil ]   │
│ • Globalny     │   historia   [współdz.|odr.] │ [ Wczytaj profil ]  │
│ • Środowisko   │   pamięć     [współdz.|odr.] │ [ Przypisz do… ]    │
│ • Moduł        │   kontekst   [współdz.|odr.] │                     │
│ • Para modułów │                              │ Warstwa:            │
│ • Projekt      │ Izolacja techniczna:         │  ( ) domyślna       │
│ • Karta sesji  │   katalog roboczy   [wł|wył] │  ( ) sesji          │
│ • Rola         │   środowisko proc.  [wł|wył] │                     │
│                │   dane/konfig. model[wł|wył] │ Polityka efektywna: │
│                │   dostęp sieciowy   [wł|wył] │  (podgląd reguł     │
│                │   odczyt/zapis plik.[wł|wył] │   dziedziczonych)   │
│                │   konto i token     [wł|wył] │                     │
│                │   model procesu     [wł|wył] │ [?] objaśnienia     │
│                │   serwer wykonania  [wł|wył] │     kontekstowe     │
└────────────────┴──────────────────────────────┴─────────────────────┘
```

| Panel | Kolumna | Zawartość | Rozdział |
|---|---|---|---|
| Selektor zasięgu | lewa | Wybór poziomu, na którym reguła izolacji ma obowiązywać | rozdz. 6.2 |
| Macierz izolacji | środkowa | Dwie grupy przełączników: izolacja kontekstu (historia, pamięć, kontekst — „współdzielone” albo „odrębne”) i izolacja techniczna (osiem zakresów — „włączony” albo „wyłączony”); każdy przełącznik z objaśnieniem `[?]` (rozdz. 3.3) | rozdz. 6.3 |
| Profil i podgląd | prawa | Zapis, wczytanie i przypisanie profilu izolacji (rozdz. 8.2), wybór warstwy (domyślna albo sesji, rozdz. 4.3) oraz podgląd polityki efektywnej — wynikowego zestawu reguł po uwzględnieniu dziedziczenia | rozdz. 6.4 |

### 6.2. Siedem poziomów zasięgu

Izolację ustawia się na każdym z siedmiu poziomów zasięgu, udostępnianych jednocześnie jako hierarchia wybierana w lewej kolumnie selektora. Ograniczenie konfiguracji do jednego, z góry narzuconego poziomu byłoby rodzajem twardej reguły, sprzecznym z zasadą pełnej konfigurowalności.

| Poziom zasięgu | Przykład zastosowania | Pierwszeństwo |
|---|---|---|
| Globalny (domyślny) | Domyślna polityka izolacji obowiązująca w całej platformie | Najniższe — warstwa bazowa |
| Środowisko | Odmienna polityka dla środowiska CodeStudio niż dla TalkIn | ↑ |
| Moduł | Odrębna pamięć przypisana modułowi Developer | ↑ |
| Para modułów (relacja) | Współdzielenie operacji kontekstowych AI między Studio a Translate | ↑ |
| Projekt | Rozdzielenie kontekstu między dwoma projektami w module Workspace | ↑ |
| Karta sesji | Jednorazowe współdzielenie historii między dwiema otwartymi kartami | ↑ |
| Rola (MultitaskingAI) | Profil izolacji przypisany roli Executor 1 | Najwyższe |

Regułę pierwszeństwa i dziedziczenia w siedmiopoziomowej hierarchii przedstawia poniższy schemat: reguła z poziomu węższego wygrywa z regułą z poziomu szerszego, a brak ustawienia oznacza dziedziczenie z poziomu bezpośrednio szerszego, aż do poziomu globalnego.

```
PIERWSZEŃSTWO REGUŁ IZOLACJI (od najwyższego do najniższego)

  Rola (MultitaskingAI)       ▲  wygrywa z każdym poziomem szerszym
  Karta sesji                 │
  Projekt                     │  reguła z poziomu węższego ma pierwszeństwo
  Para modułów (relacja)      │  brak ustawienia → dziedziczenie z poziomu
  Moduł                       │  bezpośrednio szerszego
  Środowisko                  │
  Globalny (domyślny)         ▼  warstwa bazowa — wartość dziedziczona ostatecznie
```

Reguły z różnych poziomów rozstrzyga zasada pierwszeństwa zasięgu najbardziej szczegółowego: ustawienie określone na poziomie węższym ma pierwszeństwo przed ustawieniem z poziomu szerszego. Jest to zastosowanie tej samej reguły dziedziczenia, która obowiązuje cztery warstwy ogólne opisane w rozdziale 4.2, rozszerzonej tutaj o trzy dodatkowe poziomy pośrednie (moduł, para modułów, rola) właściwe wyłącznie izolacji.

### 6.3. Macierz izolacji

Macierz izolacji obejmuje dwie niezależne grupy ustawień, stosowane łącznie na wybranym poziomie zasięgu.

| Grupa | Pozycje | Stany |
|---|---|---|
| Izolacja kontekstu | historia, pamięć, kontekst | współdzielone \| odrębne |
| Izolacja techniczna | katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania | włączony \| wyłączony |

Osiem pozycji izolacji technicznej odpowiada wprost ośmiu zakresom zdefiniowanym w rozdziale 12 Architektury; macierz izolacji nie wprowadza żadnego zakresu ponad ten zbiór — jedynie udostępnia go w tym samym oknie, w którym konfiguruje się izolację kontekstu, oraz wiąże go z siedmiopoziomową hierarchią zasięgów opisaną w rozdziale 6.2.

### 6.4. Profile izolacji i podgląd polityki efektywnej

Zestaw ustawień izolacji zapisuje się jako nazwany profil izolacji i przypisuje do sesji lub roli, zgodnie z mechanizmem profilu opisanym w rozdziale 12 Architektury. Profil grupuje wybrane wartości izolacji kontekstu i izolacji technicznej dla wskazanego poziomu zasięgu, dzięki czemu raz zbudowana polityka jest stosowana wielokrotnie bez ponownej konfiguracji. Szablon profilu izolacji zawiera Załącznik D Koncepcji platformy oraz Załącznik B niniejszego dokumentu.

Prawa kolumna profilu i podglądu udostępnia wybór warstwy konfiguracji — domyślnej albo sesji, zgodnie z rozdziałem 4.3 — oraz podgląd polityki efektywnej: wynikowego zestawu reguł izolacji po uwzględnieniu dziedziczenia między wszystkimi poziomami zasięgu, od globalnego aż do poziomu aktualnie wybranego w lewej kolumnie selektora. Podgląd polityki efektywnej pozwala użytkownikowi zweryfikować, jaka reguła faktycznie obowiązuje daną sesję lub rolę, zanim dokona zmiany — istotne w sytuacji, gdy reguły ustanowione na kilku poziomach zasięgu jednocześnie oddziałują na ten sam element platformy.

### 6.5. Stan wyjściowy

Zgodnie z zasadą pełnej konfigurowalności oraz z nadrzędną zasadą braku twardych blokad, okno konfiguracji punktów izolacji nigdy nie wymusza izolacji — udostępnia ją jako możliwość. Stan wyjściowy, obowiązujący dopóki użytkownik nie utworzy żadnego profilu ani nie zmieni żadnego ustawienia, zestawia poniższa tabela.

| Element | Stan wyjściowy | Charakter stanu wyjściowego |
|---|---|---|
| Historia, pamięć i kontekst każdej nowej karty sesji | odrębne dla każdej karty | domyślny, uporządkowany podział pracy, a nie ograniczenie możliwości |
| Osiem zakresów izolacji technicznej | żaden niewłączony | proces sesji nie jest technicznie ograniczany — pełna operacyjna swoboda |

Taki stan wyjściowy łączy przejrzysty porządek kontekstu z pełną operacyjną swobodą procesu. Wszelka dalsza izolacja — zarówno poluzowanie podziału kontekstu przez współdzielenie, jak i włączenie zakresów izolacji technicznej — jest świadomą decyzją użytkownika, a nie ustawieniem narzuconym przez platformę.

---

## 7. Warstwy widoczności funkcji

Interfejs Danaco Console ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności. Funkcja niepotrzebna do realizacji aktualnego zadania nie jest widoczna. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji; liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu. Niniejszy rozdział opisuje grupę ustawień sterujących tym mechanizmem, wskazanych w rozdziale 5.16.

### 7.1. Cztery warstwy widoczności

Każdy element interfejsu należy do dokładnie jednej z czterech warstw widoczności. Przypisaniem steruje ustawienie `visibility.function.layer`.

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania; zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie (`visibility.layer2.autocollapse`) |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana (`visibility.layer3.trigger`) |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli; użytkownik podstawowy nie widzi tych elementów |

Mechanizmy ukrywania funkcjonalności, którymi steruje ta grupa ustawień, zestawia poniższa tabela.

| Mechanizm | Działanie | Ustawienie |
|---|---|---|
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`); lista pozycji rozwija się po kliknięciu | `visibility.layer3.trigger` |
| Panele wysuwane | Funkcjonalność umieszczana jest w kolumnach bocznych, panelach wysuwanych i oknach popover; po zamknięciu panel znika całkowicie z przestrzeni roboczej | `visibility.function.layer` |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji | `visibility.layer3.trigger` |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor | `visibility.context.tags` |

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu nie występuje — ukrycie zmniejsza chaos wizualny i nie utrudnia dostępu.

### 7.2. Profile warstw według roli użytkownika

Ustawienie `visibility.profile.role` wiąże rolę użytkownika z nazwanym zestawem przypisań funkcji do warstw. Profil warstw jest odmianą profilu konfiguracji opisanego w rozdziale 8: zapisuje się go, wczytuje i przypisuje tymi samymi trzema operacjami.

| Profil warstw | Zakres ujawnionych warstw | Charakterystyka |
|---|---|---|
| Użytkownik podstawowy | 1–3 | Elementy warstwy 4 pozostają niewidoczne; dostęp do nich prowadzi wyłącznie przez polecenie języka naturalnego w Chat Window |
| Użytkownik zaawansowany | 1–3 oraz wybrane elementy warstwy 4 | Wyszukiwarka funkcji i skróty klawiszowe obejmują wskazane funkcje eksperckie |
| Operator | 1–4 | Pełny zakres funkcji, z trybem administracyjnym wywoływanym z palety poleceń po potwierdzeniu uprawnień (rozdz. 7.4; `interfejs-uzytkownika/konfiguracja.md`, rozdz. 9.6) |

Profil warstw obowiązuje na warstwie globalnej, środowiska oraz roli użytkownika; wartość węższa nakłada się na szerszą zgodnie z regułą dziedziczenia opisaną w rozdziale 4.2.

### 7.3. Wyszukiwarka funkcji, skróty klawiszowe i polecenie języka naturalnego

Trzy drogi dostępu prowadzą do funkcji warstw 3 i 4 z pominięciem hierarchii menu. Każda z nich realizuje zasadę jednego kliknięcia.

| Droga dostępu | Zakres | Ustawienia |
|---|---|---|
| Wyszukiwarka funkcji | Funkcje bieżącego modułu, funkcje wszystkich modułów albo funkcje wraz z komponentami własnymi; paleta poleceń otwiera się skrótem `Ctrl/Cmd + K` i wywołuje funkcję bezpośrednio z listy wyników (`interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 15.8) | `visibility.search.enabled`, `visibility.search.scope`, `visibility.search.shortcut` |
| Skróty klawiszowe | Katalog skrótów platformy przypisany funkcjom; przypisania podlegają zmianie na warstwie globalnej, środowiska, modułu i roli użytkownika | `visibility.shortcuts.map` |
| Polecenie języka naturalnego | Wywołanie funkcji poleceniem wpisanym w Chat Window (rozdz. 5.14) — droga dostępu obowiązująca również dla profilu użytkownika podstawowego | `visibility.nlcommand.enabled` |

### 7.4. Tryb administracyjny

Tryb administracyjny ujawnia funkcje warstwy 4: narzędzia diagnostyczne niskiego poziomu, konfigurację ról oraz operacje na procesach sesji. Drogę wejścia w tryb, model uprawnień, makietę znacznika stanu i przejścia stanu opisuje `interfejs-uzytkownika/konfiguracja.md` (rozdz. 9.6) — miejsce rozstrzygające. Niniejszy rozdział podaje wyłącznie klucze konfiguracji trybu.

| Klucz | Zakres wartości | Warstwy | Wartość domyślna |
|---|---|---|---|
| `visibility.admin.mode` | Zasięg trybu: sesja, urządzenie, rola użytkownika | globalna, środowisko, rola użytkownika | tryb nieaktywny |
| `visibility.admin.scope` | Zbiór ujawnianych funkcji: narzędzia diagnostyczne, konfiguracja ról, operacje na procesach sesji | globalna, środowisko | narzędzia diagnostyczne |

Wartością domyślną jest tryb nieaktywny: platforma pozostaje w stanie spoczynku interfejsu, w którym widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

---

## 8. Profile konfiguracji

### 8.1. Zasada wspólna

Mechanizm profilu — nazwanego, zapisywalnego zestawu wartości ustawień, przypisywalnego wielokrotnie bez ponownej konfiguracji — nie jest wyłączną własnością zakresu izolacji (rozdz. 6.4). Ten sam mechanizm, wynikający z rozdziału 12 Architektury, obejmuje inne zakresy okna konfiguracji, w których raz zbudowana konfiguracja ma wartość powtarzalną.

| Rodzaj profilu | Zawartość | Przypisywalny do | Podstawa |
|---|---|---|---|
| Profil izolacji | Izolacja kontekstu i izolacja techniczna dla wskazanego poziomu zasięgu | sesja, rola | rozdz. 6.4 |
| Zespół (MultitaskingAI) | Presety ról, powiązań i kolejek środowiska MultitaskingAI | środowisko MultitaskingAI | Koncepcja platformy, rozdz. 13.8, sekcja Zespoły panelu orkiestracji |
| Agent (komponent własny) | Model bazowy, tożsamość, instrukcje systemowe, skille, konektory, pamięć, uprawnienia | dowolny moduł, w którym agent jest wykorzystywany | Koncepcja platformy, rozdz. 11.15 |
| Profil asystenta (komponent własny) | Głos i synteza mowy, zakres poleceń, obsługiwane akcje | środowisko TalkIn, moduł Assistant | Koncepcja platformy, rozdz. 11.10 |
| Profil warstw widoczności | Przypisania funkcji do warstw 1–4, skróty klawiszowe, zakres wyszukiwarki funkcji | rola użytkownika, środowisko | rozdz. 7.2 |

### 8.2. Zapis, wczytanie i przypisanie

Każdy rodzaj profilu udostępnia trzy operacje analogiczne do tych opisanych dla profilu izolacji w prawej kolumnie profilu i podglądu (rozdz. 6.1).

```
CYKL ŻYCIA PROFILU KONFIGURACJI

  bieżąca konfiguracja
        │  [ Zapisz profil ]
        ▼
  nazwany profil w pamięci aplikacji
        │  [ Wczytaj profil ]
        ▼
  profil wczytany do bieżącego kontekstu
        │  [ Przypisz do… ]
        ▼
  profil przypisany do elementu platformy
  (sesja | rola | rola użytkownika | środowisko | moduł — zależnie od rodzaju profilu)
```

Operacje te — zapis bieżącej konfiguracji jako nazwanego profilu, wczytanie zapisanego profilu do bieżącego kontekstu oraz przypisanie profilu do wskazanego elementu platformy (sesji, roli, roli użytkownika, środowiska lub modułu, zależnie od rodzaju profilu) — realizowane są przez kontrakty komunikacji opisane w rozdziale 9.

---

## 9. Komunikacja zmian konfiguracji

Zmiany konfiguracji dokonane w oknie konfiguracji przekazywane są między klientem a serwerem tym samym kanałem WebSocket, który obsługuje pozostałe polecenia i zdarzenia platformy (Architektura, rozdz. 11). Poniższa tabela wskazuje polecenia i zdarzenia dokumentu Kontrakty komunikacji właściwe zakresom opisanym w rozdziale 5, oknu izolacji opisanemu w rozdziale 6 oraz warstwom widoczności opisanym w rozdziale 7.

| Zakres | Polecenie (klient → serwer) | Zdarzenie (serwer → klient) |
|---|---|---|
| Ogólne ustawienia (rozdz. 5) | `config.get`, `config.set` | `config.changed` |
| Izolacja (rozdz. 6) | `isolation.profile.set`, `isolation.profile.assign` | `config.changed` |
| Pamięć (5.10) | `memory.list`, `memory.set`, `memory.toggle`, `memory.detach`, `memory.delete` | `memory.changed` |
| Historia (5.9) | `history.load`, `history.delete` | `history.changed` |
| Akcje (5.3) | `task.create`, `task.update`, `queue.action`, `schedule.set`, `orchestration.define` | `task.changed` |
| Zachowanie i tożsamość modeli (5.4, 5.5) | `model.channel.set`, `agent.create`, `agent.update`, `role.assign` | `agent.changed`, `process.status` |
| Rozszerzenia (5.7) | `extension.list`, `extension.install`, `extension.toggle` | — |
| Komponenty własne — projekt (5.13) | `project.create`, `project.update`, `project.delete` | `project.changed` |
| Chat Window (5.14) | `config.get`, `config.set` | `config.changed` |
| Execution Loop Window (5.15) | `orchestration.define`, `queue.action`, `task.update`, `role.assign` | `task.changed`, `process.status` |
| Warstwy widoczności funkcji (5.16) | `config.get`, `config.set` | `config.changed` |

Serwer jest jedynym źródłem prawdy: każda zmiana konfiguracji dokonana na jednym urządzeniu jest rozgłaszana zdarzeniem do wszystkich pozostałych połączonych urządzeń (Architektura, rozdz. 16). Przebieg pojedynczej zmiany konfiguracji przedstawia poniższy schemat sekwencji.

```
SEKWENCJA ZMIANY KONFIGURACJI (serwer jako jedyne źródło prawdy)

  Urządzenie A            Serwer                       Urządzenie B
      │                     │                              │
      │  config.set  ──────►│                              │
      │                     │  zapis stanu konfiguracji    │
      │  ◄── config.changed │── config.changed ───────────►│
      │  (potwierdzenie      │      (rozgłoszenie do        │  (aktualizacja
      │   dla źródła zmiany)  │       pozostałych urządzeń)  │   widoku)
```

Dzięki temu okno konfiguracji prezentuje spójny stan niezależnie od tego, z którego urządzenia Operator z niego korzysta.

---

## 10. Zgodność z zasadami nadrzędnymi platformy

Model konfiguracji opisany w niniejszym dokumencie jest w całości podporządkowany czterem zasadom nadrzędnym platformy ustanowionym w rozdziale 14 Koncepcji platformy.

| Zasada nadrzędna | Realizacja w modelu konfiguracji |
|---|---|
| Pełna kompozycyjność | Szesnaście zakresów ustawień (rozdz. 5) nie ogranicza dostępności żadnej funkcji do wąskiego podzbioru modułów — każdy zakres jest dostępny wszędzie, gdzie ma zastosowanie, zgodnie z macierzą warstw właściwą danemu ustawieniu. |
| Pełna konfigurowalność (zasada centralna) | Cały niniejszy dokument jest rozwinięciem tej zasady — każde ustawienie opisane w rozdziale 5 oraz każdy przełącznik macierzy izolacji (rozdz. 6.3) podlega jawnej konfiguracji, bez wyjątku. |
| Jawność i konfigurowalność zależności | Zakres Integracje (5.8) oraz mechanizm powiązań komponentów własnych (5.13) czynią każdą zależność między elementami platformy widoczną i odwracalną, zgodnie z zasadą, że powiązania są możliwością, a nie regułą wbudowaną. |
| Orkiestracja równoległej pracy modeli i agentów | Zakresy Akcje (5.3), Zachowanie modeli (5.4), Tożsamość modeli (5.5) oraz Execution Loop Window (5.15) obejmują wprost mechanizmy środowiska MultitaskingAI: role, silnik kolejek, Subagent Network, równoległość zadań, politykę ponowień i punkty zatwierdzeń. |

---

## 11. Słowniczek pojęć

**Okno konfiguracji** — jedno z dwóch okien realizujących interfejs użytkownika platformy (obok okna operacyjnego), udostępniające pełny zakres ustawień wpływających na aplikację, procesy, akcje, zachowanie modeli, tożsamość modeli, prompty systemowe, rozszerzenia, integracje, historię, pamięć, izolację, karty sesji i komponenty własne (Architektura, rozdz. 13; niniejszy dokument, rozdz. 5).

**Zakres ustawień** — jeden z szesnastu obszarów tematycznych okna konfiguracji, wskazywany w panelu nawigacji zakresów (rozdz. 3.2) i opisany szczegółowo w rozdziale 5.

**Warstwa konfiguracji** — jeden z czterech poziomów ogólnych, na których obowiązuje wartość ustawienia: globalna, środowisko, projekt, sesja (rozdz. 4.1). Odrębna od poziomu zasięgu izolacji (rozdz. 6.2), który rozszerza tę warstwowość o trzy poziomy pośrednie właściwe wyłącznie izolacji.

**Objaśnienie kontekstowe** — treść oznaczona symbolem `[?]`, towarzysząca każdemu ustawieniu okna konfiguracji, opisująca jego działanie i wpływ na aplikację; część definicji ustawienia (rozdz. 3.3).

**Polityka efektywna** — wynikowy zestaw reguł izolacji obowiązujący dany element platformy po uwzględnieniu dziedziczenia między wszystkimi poziomami zasięgu, prezentowany w panelu profilu i podglądu okna konfiguracji punktów izolacji (rozdz. 6.4).

**Profil konfiguracji** — nazwany, zapisywalny zestaw wartości ustawień właściwych danemu zakresowi (izolacja, zespół MultitaskingAI, agent, profil asystenta), przypisywalny wielokrotnie do sesji, roli, środowiska lub modułu bez ponownej konfiguracji (rozdz. 8).

**Warstwa sesji** — najwęższa z czterech warstw ogólnych oraz jedna z siedmiu poziomów zasięgu izolacji; jedyna warstwa dostępna zarówno z pełnego okna konfiguracji, jak i z uproszczonego menu kontekstowego okna operacyjnego (rozdz. 3.1, 4.3).

---

## Załącznik A. Zbiorcza macierz zakresów ustawień i warstw

Zestawienie porządkuje siedemnaście zakresów ustawień opisanych w rozdziale 5 według głównych warstw, na których dopuszczają konfigurację. Znak „×” oznacza, że dany zakres dopuszcza ustawienia na wskazanej warstwie; symbol „rozdz. 6” oznacza zastosowanie rozszerzonej, siedmiopoziomowej hierarchii zasięgów właściwej izolacji.

| Zakres ustawień | Globalna | Środowisko | Moduł / para modułów | Projekt | Sesja | Rola (MultitaskingAI) | Komponent własny | Rola użytkownika |
|---|---|---|---|---|---|---|---|---|
| 5.1 Aplikacja | × | | | | × | | | |
| 5.2 Procesy | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | | |
| 5.3 Akcje | × | | × | × | × | | | |
| 5.4 Zachowanie modeli | | | | | × | × | × | |
| 5.5 Tożsamość modeli | | | | | | × | × | |
| 5.6 Prompty systemowe | | | | × | | × | × | |
| 5.7 Rozszerzenia | × | | × | | | | × | |
| 5.8 Integracje | × | × | × | | | | × | |
| 5.9 Historia | × | × | | × | × | | | |
| 5.10 Pamięć | × | × | | × | × | | | |
| 5.11 Izolacja | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | rozdz. 6 | | |
| 5.12 Karty sesji | | | | | × | | | |
| 5.13 Komponenty własne | × | | | × | | | × | |
| 5.14 Chat Window | × | × | | × | × | | | |
| 5.15 Execution Loop Window | × | × | | × | × | × | | |
| 5.16 Warstwy widoczności funkcji | × | × | × | | × | | | × |
| 5.17 Always On Display | × | × | × | × | × | | | |

---

## Załącznik B. Szablony

### B.1. Szablon ustawienia

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Identyfikator | Unikalny identyfikator ustawienia | wartość systemowa |
| Klucz | Nazwa techniczna ustawienia | np. `history.retention` |
| Zakres ustawień | Jeden z szesnastu zakresów rozdziału 5 | np. Historia (5.9) |
| Wartość | Bieżąca wartość ustawienia | zależna od typu ustawienia |
| Poziom warstwy | Warstwa, na której wartość jest ustalona | globalna \| środowisko \| moduł \| projekt \| sesja \| rola \| rola użytkownika \| komponent własny |
| Wartość domyślna | Wartość przyjmowana przy braku ustawienia | zgodnie z tabelami rozdziału 5 |
| Objaśnienie kontekstowe | Treść wyświetlana pod symbolem `[?]` | dowolny tekst opisujący działanie i wpływ na aplikację |

Postać blokowa szablonu ustawienia:

```
Ustawienie
  identyfikator            : <wartość systemowa>
  klucz                    : np. history.retention
  zakres ustawień          : jeden z szesnastu zakresów (rozdz. 5)
  wartość                  : <zależna od typu ustawienia>
  poziom warstwy           : globalna | środowisko | moduł | projekt | sesja |
                             rola | rola użytkownika | komponent własny
  wartość domyślna         : zgodnie z tabelami rozdziału 5
  objaśnienie kontekstowe  : tekst opisujący działanie i wpływ na aplikację  [?]
```

### B.2. Szablon profilu konfiguracji

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Nazwa profilu | Nazwa własna profilu | dowolny tekst |
| Rodzaj profilu | Zakres, którego profil dotyczy | izolacja \| zespół MultitaskingAI \| agent \| profil asystenta \| profil warstw widoczności (rozdz. 7.2, 8.1) |
| Zawartość | Zestaw wartości ustawień objętych profilem | zależna od rodzaju profilu |
| Poziom zasięgu (dla profilu izolacji) | Poziom obowiązywania reguły | globalny \| środowisko \| moduł \| para modułów \| projekt \| karta sesji \| rola (rozdz. 6.2) |
| Przypisanie | Element, do którego profil jest przypisany | sesja \| rola \| rola użytkownika \| środowisko \| moduł |
| Warstwa | Warstwa konfiguracji, na której profil obowiązuje | domyślna \| sesji (rozdz. 4.3) |

Postać blokowa szablonu profilu konfiguracji:

```
Profil konfiguracji
  nazwa profilu    : <nazwa własna>
  rodzaj profilu   : izolacja | zespół MultitaskingAI | agent | profil asystenta |
                     profil warstw widoczności
  zawartość        : <zestaw wartości ustawień objętych profilem>
  poziom zasięgu   : globalny | środowisko | moduł | para modułów |
                     projekt | karta sesji | rola      (dla profilu izolacji)
  przypisanie      : sesja | rola | rola użytkownika | środowisko | moduł
  warstwa          : domyślna | sesji
```

---

## Załącznik C. Scenariusze konfiguracji

Scenariusze przedstawiono jako ponumerowane mini-przepływy ilustrujące użycie mechanizmów opisanych w rozdziałach 3–8.

### C.1. Skrócona zmiana kanału modelu w trakcie sesji

**Kontekst:** Operator pracujący w module Developer chce tymczasowo przełączyć kanał modelu z API na CLI dla bieżącej karty sesji, bez przerywania trwającej pracy.

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Otwiera uproszczone menu kontekstowe okna operacyjnego | rozdz. 3.1 |
| 2 | Wybiera zakres Zachowanie modeli | rozdz. 5.4 |
| 3 | Zmienia wartość ustawienia Kanał modelu na warstwie sesji z API na CLI | rozdz. 4.3 |
| 4 | Zmiana obejmuje wyłącznie tę kartę sesji; wartość globalna oraz konfiguracja pozostałych kart pozostają nienaruszone | rozdz. 4.2 |

### C.2. Budowa profilu izolacji dla roli Executor 1

**Kontekst:** Operator konfigurujący środowisko MultitaskingAI do pracy nad projektem zawierającym dane wrażliwe chce ograniczyć proces roli Executor 1, pozostawiając izolację kontekstu współdzieloną z Coordinatorem.

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Otwiera okno konfiguracji punktów izolacji | rozdz. 6 |
| 2 | W panelu selektora wybiera zasięg „Rola” | rozdz. 6.2 |
| 3 | W macierzy izolacji włącza zakresy techniczne Konto i token per sesja oraz Katalog roboczy sesji dla roli Executor 1 | rozdz. 6.3 |
| 4 | Pozostawia izolację kontekstu współdzieloną z Coordinatorem | rozdz. 6.3 |
| 5 | Zapisuje ustawienia jako profil izolacji i przypisuje go roli Executor 1 | rozdz. 6.4, 8.2 |
| 6 | Podgląd polityki efektywnej potwierdza, że reguła roli ma pierwszeństwo przed globalną polityką braku izolacji technicznej | rozdz. 6.2, 6.5 |

### C.3. Powiązanie automatyki z modułem Developer

**Kontekst:** Operator tworzy w strefie 2 strony głównej automatykę cyklicznego przeglądu jakości kodu (komponent własny, rozdz. 5.13) i chce, aby była dostępna w sesjach modułu Developer.

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Tworzy automatykę cyklicznego przeglądu jakości kodu jako komponent własny w strefie 2 strony głównej | rozdz. 5.13 |
| 2 | W oknie konfiguracji, w zakresie Integracje, ustanawia powiązanie komponentu własnego z modułem Developer | rozdz. 5.8 |
| 3 | Od tego momentu automatyka jest wybieralna w sesjach modułu Developer | Koncepcja platformy, rozdz. 2.3 |
| 4 | Powiązanie pozostaje jawne, odwracalne i widoczne w oknie konfiguracji, a nie ukryte w kodzie modułu | rozdz. 2, zasada „Jawność zależności” |

### C.4. Ustawienie pętli wykonawczej dla zlecenia wielozadaniowego

**Kontekst:** Operator prowadzący w środowisku MultitaskingAI zlecenie złożone z kilkunastu zadań ogranicza równoległość pracy i zaostrza kontrolę jakości na czas realizacji tego zlecenia.

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Otwiera okno konfiguracji i wybiera zakres Execution Loop Window | rozdz. 5.15 |
| 2 | Ustawia `loop.concurrency` na warstwie sesji na 2 zadania | rozdz. 4.3, 5.15 |
| 3 | Ustawia `loop.quality.scope` na wartość „każde zadanie” oraz `loop.quality.role` na Reviewer dla roli kontrolnej | rozdz. 5.15 |
| 4 | Rozszerza `loop.approval.points` o zatwierdzenie przejścia między etapami zlecenia | rozdz. 5.15 |
| 5 | Okno pętli wykonawczej otwiera się w kolumnie sąsiadującej z Chat Window i prezentuje kolejkę zadań, wyniki kontroli jakości oraz sterowanie przebiegiem | rozdz. 5.15 |
| 6 | Zmiana obejmuje wyłącznie bieżącą kartę sesji; wartości warstwy środowiska pozostają nienaruszone | rozdz. 4.2 |

### C.5. Ujawnienie funkcji eksperckich dla roli Operatora

**Kontekst:** Operator potrzebuje dostępu do narzędzi diagnostycznych niskiego poziomu, pozostawiając interfejs użytkowników podstawowych bez zmian.

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Otwiera okno konfiguracji i wybiera zakres Warstwy widoczności funkcji | rozdz. 5.16 |
| 2 | Przypisuje roli użytkownika Operator profil warstw ujawniający warstwy 1–4 (`visibility.profile.role`) | rozdz. 7.2 |
| 3 | Włącza `visibility.admin.mode` i ogranicza `visibility.admin.scope` do narzędzi diagnostycznych | rozdz. 7.4 |
| 4 | Funkcje warstwy 4 pozostają dostępne wyłącznie dla tej roli; profil użytkownika podstawowego pozostaje niezmieniony | rozdz. 7.1, 7.2 |
| 5 | Wyszukiwarka funkcji i skróty klawiszowe obejmują ujawnione funkcje, zachowując zasadę jednego kliknięcia | rozdz. 7.3 |

---

*Koniec dokumentu. Danaco Console — Model konfiguracji, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
