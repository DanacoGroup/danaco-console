# Danaco Console — Moduł Assistant

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
| **Tytuł** | Moduł Assistant |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | projektant (co, gdzie, w jakiej formie) · deweloper (co zbudować) |
| **Przeznaczenie** | Ustala interfejs modułu Assistant: pełny zakres funkcji, narzędzi, okien operacyjnych i punktów sterowania komponentu własnego Profil asystenta |
| **Zakres** | okna operacyjne modułu wykorzystywane w środowisku TalkIn, katalog elementów interfejsu, przepływy pracy, komendy kontraktu obszaru `assistant` |
| **Poza zakresem** | model danych Profilu asystenta i pamięci semantycznej jako bytów platformy — [Model danych](../architektura/model-danych.md) |
| **Dokument nadrzędny** | [Specyfikacja modułów](../specyfikacje/specyfikacja-modulow.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Specyfikacja okien operacyjnych](../specyfikacje/specyfikacja-okien-operacyjnych.md) · [System wizualny](../interfejs-uzytkownika/system-wizualny.md) · [Model danych](../architektura/model-danych.md) · [Izolacja i zależności](../architektura/izolacja-i-zaleznosci.md) · [TalkIn](../srodowiska/talkin.md) · [Always On Display](../funkcje-globalne/always-on-display.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/assistant.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `assistant`) · `design/zasoby/zetony/zetony.css` · `design/zasoby/css/komponenty.css`, `rama.css`, `prototyp.css` |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); izolacja i uprawnienia są ustawieniami konfiguracyjnymi Operatora, nie wymogiem; zero blokad w interfejsie; klucze jawne; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
   - [1.1 Definicja](#11-definicja)
   - [1.2 Obszary pokrycia modułu](#12-obszary-pokrycia-modułu)
   - [1.3 Granica tematyczna modułu](#13-granica-tematyczna-modułu)
   - [1.4 Dla kogo](#14-dla-kogo)
   - [1.5 Po co — wartość modułu](#15-po-co--wartość-modułu)
   - [1.6 Miejsce w architekturze platformy](#16-miejsce-w-architekturze-platformy)
   - [1.7 Dostępność i forma udostępnienia](#17-dostępność-i-forma-udostępnienia)
2. [Komplet okien operacyjnych modułu](#2-komplet-okien-operacyjnych-modułu)
   - [2.1 Warstwy widoczności w module](#21-warstwy-widoczności-w-module)
3. [Specyfikacja okien operacyjnych](#3-specyfikacja-okien-operacyjnych)
   - [3.1 Chat Window (okno wspólne)](#31-chat-window-okno-wspólne)
   - [3.2 Execution Loop Window (okno wspólne)](#32-execution-loop-window-okno-wspólne)
   - [3.3 Voice Console](#33-voice-console)
   - [3.4 Actions Monitor](#34-actions-monitor)
   - [3.5 Activity Feed](#35-activity-feed)
   - [3.6 Memory & Context Manager](#36-memory--context-manager)
   - [3.7 Command & Tools Hub](#37-command--tools-hub)
4. [Przepływy pracy](#4-przepływy-pracy)
   - [4.1 Przepływ podstawowy — polecenie głosowe wieloetapowe](#41-przepływ-podstawowy--polecenie-głosowe-wieloetapowe)
   - [4.2 Przepływ rozszerzony — polecenie obejmujące inny moduł platformy](#42-przepływ-rozszerzony--polecenie-obejmujące-inny-moduł-platformy)
   - [4.3 Przepływ rozszerzony — korekta błędnie rozpoznanego polecenia](#43-przepływ-rozszerzony--korekta-błędnie-rozpoznanego-polecenia)
   - [4.4 Przepływ rozszerzony — zadanie o skutkach ubocznych](#44-przepływ-rozszerzony--zadanie-o-skutkach-ubocznych)
   - [4.5 Przepływ rozszerzony — praca z pamięcią i kontekstem](#45-przepływ-rozszerzony--praca-z-pamięcią-i-kontekstem)
   - [4.6 Przepływ rozszerzony — rozgraniczenie toru głosowego wobec Always On Display](#46-przepływ-rozszerzony--rozgraniczenie-toru-głosowego-wobec-always-on-display)
5. [Stany, dane i powiązania](#5-stany-dane-i-powiązania)
   - [5.1 Model stanów zlecenia głosowego](#51-model-stanów-zlecenia-głosowego)
   - [5.2 Model danych wykorzystywany przez moduł](#52-model-danych-wykorzystywany-przez-moduł)
   - [5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Assistant](#53-izolacja-i-konfigurowalność--punkty-właściwe-modułowi-assistant)
   - [5.4 Powiązania z innymi modułami](#54-powiązania-z-innymi-modułami)
6. [Scenariusze użycia](#6-scenariusze-użycia)
7. [Katalog funkcji i narzędzi](#7-katalog-funkcji-i-narzędzi)
   - [7.1 Rozmowa i konwersacja (rdzeń)](#71-rozmowa-i-konwersacja-rdzeń)
   - [7.2 Głos — rozpoznawanie, synteza, wybudzanie](#72-głos--rozpoznawanie-synteza-wybudzanie)
   - [7.3 Dyktowanie i produktywność mowy](#73-dyktowanie-i-produktywność-mowy)
   - [7.4 Pamięć i konteksty](#74-pamięć-i-konteksty)
   - [7.5 Szybkie akcje, wywoływacz i schowek](#75-szybkie-akcje-wywoływacz-i-schowek)
   - [7.6 Wywołania narzędzi i integracje](#76-wywołania-narzędzi-i-integracje)
   - [7.7 Sterowanie zadaniami wieloetapowymi](#77-sterowanie-zadaniami-wieloetapowymi)
   - [7.8 Proaktywność, przypomnienia i rutyny](#78-proaktywność-przypomnienia-i-rutyny)
   - [7.9 Personalizacja profili asystenta](#79-personalizacja-profili-asystenta)
   - [7.10 Wielojęzyczność, dostępność i tryby](#710-wielojęzyczność-dostępność-i-tryby)
   - [7.11 Prywatność, bezpieczeństwo i higiena](#711-prywatność-bezpieczeństwo-i-higiena)
   - [7.12 Zależności techniczne — zestawienie](#712-zależności-techniczne--zestawienie)
8. [Punkty sterowania z okna konfiguracji](#8-punkty-sterowania-z-okna-konfiguracji)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)
10. [Komendy kontraktu obszaru `assistant`](#10-komendy-kontraktu-obszaru-assistant)
   - [10.1 Obszar `assistant` — 4 komendy](#101-obszar-assistant--4-komendy)
11. [Kryteria odbioru](#11-kryteria-odbioru)

---

## 1. Przeznaczenie i kontekst

### 1.1 Definicja

Assistant jest modułem osobistego asystenta ogólnego — jednego, trwale obecnego rozmówcy AI, który prowadzi rozmowę głosem i tekstem, utrzymuje pamięć i konteksty, wykonuje szybkie akcje oraz wywołuje narzędzia, realizując polecenia wieloetapowe obejmujące inne moduły platformy. Głos pozostaje filarem wyróżniającym modułu: obsługuje scenariusze, w których mowa jest wygodniejsza niż pisanie — pracę w ruchu, sterowanie zadaniami bez użycia klawiatury oraz preferencję głosowej formy komunikacji.

W odróżnieniu od pozostałych czternastu modułów platformy, Assistant jest jednocześnie komponentem własnym — Profilem asystenta, konfigurowanym w strefie 2 strony głównej — i modułem operacyjnym osadzonym w środowisku TalkIn, w którym skonfigurowany profil faktycznie działa i wykonuje polecenia.

### 1.2 Obszary pokrycia modułu

Moduł pokrywa sześć obszarów, z których każdy poza platformą wymaga osobnego programu:

| Obszar | Zakres w module Assistant |
|---|---|
| Rozmowa ogólna | Konwersacja wieloturowa głosem i tekstem, strumień odpowiedzi, wątkowanie i rozgałęzianie, persony, załączniki, tryby odpowiedzi |
| Głos | Rozpoznawanie mowy, synteza mowy, fraza wybudzająca, detekcja mowy, redukcja szumu, wielojęzyczność, dyktowanie |
| Pamięć i konteksty | Pamięć krótko- i długoterminowa, pamięć semantyczna, fakty o użytkowniku, konteksty przełączane, wstrzykiwanie wiedzy profilu |
| Szybkie akcje | Wywoływacz poleceń, paleta poleceń, makra głosowe i tekstowe, rozwijanie skrótów tekstowych, menedżer schowka, przechwytywanie notatek |
| Wywołania narzędzi | Katalog narzędzi, Model Context Protocol (MCP), umiejętności, sterowanie modułami platformy i aplikacjami zewnętrznymi |
| Proaktywność i zadania | Przypomnienia, harmonogram, wyzwalacze zdarzeniowe, rutyny, briefy dzienne, powiadomienia i potwierdzenia |

### 1.3 Granica tematyczna modułu

| Poza granicą | Gdzie to należy |
|---|---|
| Budowa i długotrwałe wykonywanie bezobsługowych procesów cyklicznych | Moduł **Automations** (Assistant inicjuje i nadzoruje pojedyncze zlecenia, nie projektuje pełnych potoków procesów) |
| Definiowanie i zarządzanie autonomicznymi agentami wielozadaniowymi | Moduł **Agents** (Profil asystenta korzysta z uprawnień analogicznych do okna Permissions Center, lecz nie jest fabryką agentów) |
| Debata i konsensus wielu modeli nad jednym problemem | Moduł **Roundtable** (Assistant rozmawia jednym głosem profilu) |
| Trwałe katalogowanie repozytorium wiedzy i plików | Moduł **Library** (Assistant zapisuje tam artefakty, nie zastępuje repozytorium) |
| Zaawansowana redakcja dokumentów, tłumaczenie jako usługa, badania wieloźródłowe | Moduły **Studio / Translate / Research** (Assistant wywołuje je jako kroki polecenia) |
| Warstwa głosowa ponad całą platformą, poza obrębem sesji | Funkcja globalna **Always On Display** — krótkie polecenia ponadkontekstowe, nawigacja między środowiskami i modułami, sterowanie procesem poza sesją; rozgraniczenie obu torów głosowych podaje [Always On Display](../funkcje-globalne/always-on-display.md), rozdz. 5 |

Granica jest wykonawcza, nie licencyjna: powiązania są jawne i konfigurowalne, a nie zamknięte. Assistant jest dyrygentem i rozmówcą, a nie miejscem trwałego przechowywania ani budowania automatyk.

### 1.4 Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Assistant |
|---|---|
| Użytkownicy pracujący w ruchu | Wydawanie poleceń bez dostępu do klawiatury |
| Operatorzy nadzorujący wiele procesów jednocześnie | Szybkie polecenia głosowe bez przerywania innej pracy wzrokowej |
| Osoby preferujące komunikację głosową | Naturalna forma interakcji zamiast pisania |
| Użytkownicy zlecający zadania wieloetapowe | Wydanie złożonego polecenia głosem i śledzenie jego realizacji bez ręcznego nadzoru każdego kroku |
| Użytkownicy prowadzący długą współpracę z asystentem | Trwała pamięć ustaleń, przełączane konteksty pracy i własna baza wiedzy profilu |

### 1.5 Po co — wartość modułu

| Problem interakcji wyłącznie tekstowej | Rozwiązanie w Assistant |
|---|---|
| Pisanie polecenia wymaga pełnej uwagi wzrokowej i rąk | Voice Console przyjmuje polecenie mową, bez przerywania innej czynności |
| Brak wglądu w postęp długiego zlecenia wydanego głosem | Actions Monitor pokazuje bieżący status realizacji na żywo, a Execution Loop Window — przebieg pętli wykonawczej |
| Trudność odtworzenia przebiegu złożonego zlecenia głosowego po fakcie | Activity Feed zachowuje chronologiczny zapis każdego kroku i jego wyniku |
| Konieczność każdorazowego opisywania preferencji AI od nowa | Profil asystenta zapamiętuje charakterystykę głosu, tempo i preferencje raz skonfigurowane |
| Rozproszenie pamięci, skrótów i narzędzi po wielu osobnych programach | Memory & Context Manager i Command & Tools Hub gromadzą pamięć, konteksty, makra i katalog narzędzi w jednym module |

### 1.6 Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │
        ▼
Strefa 2 · komponenty własne  ──────────────────────────────
        │  konfiguracja Profilu asystenta (komponent własny)
        ▼
Pamięć aplikacji — profil asystenta zapisany jako zasób użytkownika
        │
        │  wybór profilu w trakcie sesji
        ▼
ŚRODOWISKO: TalkIn ── boczna nawigacja modułów
        │
        ▼
MODUŁ: ASSISTANT (wykorzystanie operacyjne) ─────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Voice Console ·
  Actions Monitor · Activity Feed ·
  Memory & Context Manager · Command & Tools Hub
```

### 1.7 Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Miejsce konfiguracji | Strona główna, strefa 2 (komponenty własne) — tworzenie i zapis Profilu asystenta |
| Środowisko wykorzystania operacyjnego | Wyłącznie TalkIn |
| Środowiska WorkSpace, CodeStudio | Moduł niedostępny operacyjnie |
| Charakter | Jedyny moduł platformy będący jednocześnie komponentem własnym i oknem modułowym w jednym środowisku |
| Liczba okien operacyjnych | 7 (łącznie z Chat Window i Execution Loop Window) |
| Liczba profili asystenta | Dowolna — użytkownik tworzy wiele profili (różne głosy, różne zestawy uprawnień) i przełącza się między nimi |

---

## 2. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja (wspólne wszystkim modułom) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji przez cały czas pracy modułu | Główne okno komunikacji Użytkownik ↔ Wykonawca — centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu |
| 2 | Execution Loop Window | Komunikacja (wspólne wszystkim modułom) | Kolumna sąsiadująca z Chat Window, otwierana na żądanie | 2 — zawsze na żądanie; drugie okno komunikacji nigdy nie jest stałym elementem modułu (rozdz. 2.1) | Wstążka narzędziowa, ikona szybkiego dostępu w prawym górnym rogu okna centralnego aplikacji albo komponent pierwszego okna komunikacji — znacznik stanu pętli `[Pętla ▼]` i przycisk „Pętla wykonawcza” w Chat Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca — dekompozycja zlecenia głosowego na zadania, nadzór i kontrola realizacji |
| 3 | Voice Console | Okno główne (specjalne — interfejs głosowy) | Prawa kolumna, dominująca | 1 | Widoczne bez interakcji jako aktywne okno wiodące modułu | Wydawanie poleceń głosowych i odbiór odpowiedzi mową |
| 4 | Actions Monitor | Okno monitorów i wskaźników | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Wskaźnik postępu zlecenia w pasku kontekstu, przycisk „Szczegóły w Actions Monitor” w dymku odpowiedzi | Bieżący status realizacji poleceń |
| 5 | Activity Feed | Okno monitorów i wskaźników | Kolumna boczna, otwierana jako rozszerzenie boczne | 3 | Menu kebab (⋮) okna modułu, polecenie języka naturalnego „pokaż historię działań” | Chronologiczny zapis wykonanych działań |
| 6 | Memory & Context Manager | Okno zarządzania zasobem | Kolumna boczna szeroka lub obszar roboczy pełny | 3 | Znacznik kontekstu `[Kontekst ▼]` w pasku kontekstu, menu hamburger (☰) modułu | Jawna pamięć, konteksty przełączane i baza wiedzy profilu |
| 7 | Command & Tools Hub | Okno zarządzania zasobem | Kolumna boczna szeroka lub obszar roboczy pełny | 3 | Menu hamburger (☰) modułu, wywoływacz poleceń, skrót klawiszowy | Szybkie akcje, makra, katalog narzędzi i MCP, rutyny |

```
 Makieta zbiorcza — Moduł Assistant       Dostępność: TalkIn (komponent własny)
 Stan spoczynku interfejsu — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window                  │ Voice Console              │ ⋮
  nawigacja │ Użytkownik ↔ Wykonawca       │ [Profil ▼] [Model ▼]  [☰]  │
  modułów   │                              │                            │
            │ [Profil ▼] [Tryb ▼]     [⋮]  │      ((( mikrofon )))      │
            │                              │                            │
            │ transkrypcja i strumień      │  przytrzymaj i mów         │
            │ odpowiedzi                   │                            │
            │                              │  transkrypcja na żywo      │
            │ ───────────────────────────  │                            │
            │ [ 🎙 ] [ 📎 ] pole poleceń   │  [ Szybkie akcje ▼ ]       │
 ═══════════════════════════════════════════════════════════════════════════
  Pasek kontekstu: [TalkIn] [Asystent biurowy] [PL] [Pętla ▼] [Postęp ▼]
 ═══════════════════════════════════════════════════════════════════════════
```

W stanie spoczynku widoczne są wyłącznie Chat Window i Voice Console (warstwa 1) oraz zwinięte wyzwalacze: znaczniki kontekstowe paska kontekstu, menu kebab (⋮) i menu hamburger (☰). Execution Loop Window, Actions Monitor, Activity Feed, Memory & Context Manager i Command & Tools Hub pozostają niewidoczne do chwili wywołania i po zamknięciu znikają całkowicie z przestrzeni roboczej.

Regulacji podlega wyłącznie szerokość kolumn. Okna pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego.

### 2.1 Warstwy widoczności w module

Moduł Assistant realizuje zasadę nadrzędną interfejsu Danaco Console: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Pełny arsenał modułu — siedem okien operacyjnych, katalog narzędzi i serwerów MCP, umiejętności, makra, rutyny, pamięć semantyczna i reguły retencji — istnieje w architekturze i pozostaje poza polem widzenia do chwili wystąpienia potrzeby użycia. Wydanie polecenia głosem nie wymaga widoku żadnego z okien zarządzania zasobem.

**Reguła drugiego okna komunikacji.** Drugie okno komunikacji — w tym module Execution Loop Window — uruchamiane jest zawsze na żądanie i nigdy nie jest stałym elementem otwieranego modułu. Wyjątkiem są moduł okrągłego stołu (Roundtable) oraz środowisko MultitaskingAI, gdzie drugie okno komunikacji jest elementem stałym; moduł Assistant do wyjątków nie należy. Otwarcie następuje z komponentu stojącego w jednym z trzech miejsc: we wstążce narzędziowej, w ikonie szybkiego dostępu w prawym górnym rogu okna centralnego aplikacji albo w komponencie pierwszego okna komunikacji (Chat Window).

| Warstwa | Zakres w module Assistant | Sposób wywołania |
|---|---|---|
| 1 — zawsze widoczna | Chat Window z polem poleceń i przyciskiem mikrofonu, Voice Console z dużym przyciskiem mikrofonu i transkrypcją na żywo, pasek kontekstu ze znacznikami profilu, języka i stanu zlecenia. Zajmuje ponad 80% powierzchni obszaru roboczego modułu | Widoczna bez interakcji |
| 2 — widoczna na żądanie | Selektor Profilu asystenta, selektor persony i trybu odpowiedzi, selektor języka i modelu rozpoznawania, przełącznik trybu cichego, przełącznik „głośne otoczenie”, Execution Loop Window, Actions Monitor, filtr statusu zleceń | Znacznik kontekstowy (`[Profil ▼]`, `[Model ▼]`, `[Pętla ▼]`, `[Postęp ▼]`), ikona, przełącznik; po użyciu element zwija się samoczynnie. Execution Loop Window otwiera się wyłącznie na żądanie — ze wstążki narzędziowej, z ikony szybkiego dostępu w prawym górnym rogu okna centralnego aplikacji albo z komponentu Chat Window |
| 3 — rozwinięcia kontekstowe | Akcje dymka odpowiedzi (Popraw polecenie, Wstaw jako zadanie, Zapisz do Library, Kopiuj), siatka szybkich akcji, sterowanie przebiegiem pętli, Activity Feed, Memory & Context Manager, Command & Tools Hub, zakładki obszarów pamięci, katalog narzędzi i MCP, rutyny | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe dymka, panel popover, lista rozwijana `[ Szybkie akcje ▼ ]` |
| 4 — funkcje eksperckie | Edytor makra w `JSON`/`YAML`, reguły retencji i wygaszania pamięci, znaczniki wrażliwości, zakres uprawnień narzędzia i serwera MCP, limity wywołań, audyt wywołań narzędzi, eksport dziennika aktywności, przekazanie rutyny do modułu Automations, tryb ciągłego nasłuchu i konfiguracja frazy wybudzającej | Polecenie języka naturalnego w Chat Window lub Voice Console, skrót klawiszowy, wywoływacz poleceń, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

**Zasada jednego kliknięcia.** Każda ukryta funkcja modułu jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego wypowiedzianym do asystenta. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu. Głos stanowi w module Assistant równoprawną ścieżkę dostępu do warstw 2–4: polecenie „pokaż historię działań”, „otwórz pamięć”, „wstrzymaj zlecenie” osiąga funkcję bez odnajdywania jej w interfejsie.

---

## 3. Specyfikacja okien operacyjnych

### 3.1 Chat Window (okno wspólne)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — główne okno komunikacji Użytkownik ↔ Wykonawca |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 5.3) |

Chat Window jest centralnym punktem pracy użytkownika w module i podstawowym mechanizmem sterowania wszystkimi procesami realizowanymi przez moduł. Przyjmuje polecenia w języku naturalnym — mową i tekstem, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Okno zajmuje to samo miejsce układu we wszystkich modułach i środowiskach platformy.

**Zawartość i pełny arsenał funkcji.**

- Selektor Profilu asystenta oraz persony i trybu odpowiedzi (zwięzły, wyjaśniający, krok-po-kroku, ekspert, burza mózgów).
- Transkrypcja poleceń głosowych wyświetlana na żywo jako tekst, wraz z odpowiedziami asystenta, w jednym ciągłym wątku.
- Wydanie polecenia tekstem, gdy komunikacja głosowa jest niewygodna (otoczenie głośne lub wymagające ciszy).
- Edycja błędnie rozpoznanego polecenia głosowego bezpośrednio w transkrypcji, z ponownym wysłaniem poprawionej wersji.
- Strumień odpowiedzi tekstowej równoległy do syntezy mowy, wraz z cytowaniami źródeł wprowadzonych do rozmowy.
- Wątkowanie i rozgałęzianie rozmowy od dowolnego punktu bez utraty oryginału.
- Regeneracja odpowiedzi i prezentacja wariantów obok siebie.
- Wrzutnia załączników: plik, obraz, zrzut ekranu jako kontekst wypowiedzi.
- Historia poleceń, zarówno głosowych, jak i tekstowych, w jednym ciągłym zapisie.
- Akcje dymka odpowiedzi: Popraw polecenie, Szczegóły w Actions Monitor, Wstaw jako zadanie, Zapisz do Library, Kopiuj.
- Przełącznik trybu cichego (praca wyłącznie tekstowa, bez syntezy mowy).

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3

┌─ Chat Window (lewa kolumna, pełna wysokość) ─────────────────────┐
│ [Profil ▼] [Tryb ▼]                                        [ ⋮ ] │
├──────────────────────────────────────────────────────────────────┤
│  ┌─ Użytkownik (głos → tekst) ───────────────────┐               │
│  │ „umów spotkanie z zespołem na jutro na 10:00” │               │
│  └───────────────────────────────────────────────┘               │
│  ┌─ Wykonawca (tekst + synteza mowy) ─────────────────┐          │
│  │ „zaplanowano spotkanie, wysyłam zaproszenia” 🔊    │      [⋮] │
│  └────────────────────────────────────────────────────┘          │
│                                                                  │
├──────────────────────────────────────────────────────────────────┤
│ [ 🎙 przytrzymaj, by mówić ] [ 📎 ]  Pole poleceń ……… [ Wyślij ] │
└──────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Selektor profilu asystenta | Rozwijana lista w nagłówku | Wybór aktywnego Profilu asystenta bieżącej karty | Mały przycisk z etykietą | domyślny · rozwinięty | Przełącza charakterystykę głosu, uprawnień i pamięci | Nagłówek Chat Window | 2 | Znacznik kontekstowy `[Profil ▼]` w nagłówku; po wyborze lista zwija się samoczynnie |
| Selektor persony i trybu odpowiedzi | Rozwijana lista w nagłówku | Zmiana stylu odpowiedzi bez zmiany profilu głosu | Mały przycisk z etykietą | domyślny · rozwinięty | Podmienia warstwę promptu systemowego rozmowy | Nagłówek Chat Window | 2 | Znacznik kontekstowy `[Tryb ▼]` w nagłówku; po wyborze lista zwija się samoczynnie |
| Transkrypcja polecenia głosowego | Dymek wiadomości z ikoną mikrofonu | Tekstowy zapis rozpoznanej mowy | Średni blok tekstu z ikoną mikrofonu | rozpoznano poprawnie · niska pewność rozpoznania (podświetlenie) | Kliknięcie „Popraw polecenie” pozwala edytować tekst i wysłać ponownie | Historia rozmowy | 1 | Widoczna bez interakcji w strumieniu rozmowy |
| Ikona syntezy mowy | Mała ikona głośnika przy odpowiedzi | Sygnalizacja, że odpowiedź jest równolegle odczytywana na głos | Bardzo mała ikona (🔊) | odtwarzanie · wyciszone | Kliknięcie wycisza i wznawia odczyt danej odpowiedzi | Przy dymku odpowiedzi Wykonawcy | 1 | Widoczna bez interakcji przy dymku odpowiedzi |
| Przycisk „Szczegóły w Actions Monitor” | Mały przycisk `--zarys` | Przejście do statusu realizacji zlecenia | Mały przycisk tekstowy | domyślny | Otwiera Actions Monitor na pozycji odpowiadającej poleceniu | Pasek akcji dymka odpowiedzi | 3 | Menu kontekstowe dymka odpowiedzi (⋮); otwiera Actions Monitor jednym kliknięciem |
| Przycisk „Wstaw jako zadanie” | Mały przycisk `--zarys` | Zamiana fragmentu odpowiedzi w zadanie zlecenia | Mały przycisk tekstowy | domyślny | Tworzy pozycję zadania w Execution Loop Window i Actions Monitor | Pasek akcji dymka odpowiedzi | 3 | Menu kontekstowe dymka odpowiedzi (⋮) |
| Wrzutnia załączników | Ikona spinacza przy polu poleceń | Dodanie pliku, obrazu lub zrzutu jako kontekstu | Mała ikona `.dn-btn-ikona` | domyślny · przeciąganie pliku | Dołącza `artefakt` do wypowiedzi i przekazuje do ekstrakcji treści | Pole poleceń | 1 | Widoczna bez interakcji jako ikona spinacza przy polu poleceń |
| Przycisk mikrofonu „przytrzymaj, by mówić” | Duży przycisk ikonowy | Szybki dostęp do nagrywania z poziomu pola poleceń | Średni przycisk okrągły | domyślny · nagrywanie (podświetlony, pulsujący) | Przytrzymanie aktywuje nasłuch, puszczenie kończy i wysyła polecenie | Lewa strona pola poleceń | 1 | Widoczny bez interakcji przy polu poleceń |
| Przełącznik trybu cichego | Mały przełącznik w nagłówku | Praca wyłącznie tekstowa bez syntezy mowy | Mały przełącznik | głos włączony · tryb cichy | Wyłącza TTS dla bieżącej karty sesji | Nagłówek Chat Window | 2 | Menu kebab (⋮) nagłówka okna; po przełączeniu menu zwija się samoczynnie |

---

### 3.2 Execution Loop Window (okno wspólne)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja — okno pętli wykonawczej Koordynator ↔ Wykonawca |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana na żądanie (warstwa 2) |
| Izolacja domyślna | Pętla wykonawcza właściwa bieżącej karcie sesji; strumień aktualizowany kanałem WebSocket |

**Otwarcie na żądanie.** Okno należy do warstwy 2 i uruchamiane jest zawsze na żądanie — nigdy nie jest stałym elementem otwieranego modułu Assistant (reguła drugiego okna komunikacji, rozdz. 2.1). Otwarcie następuje z komponentu stojącego w jednym z trzech miejsc: we wstążce narzędziowej, w ikonie szybkiego dostępu w prawym górnym rogu okna centralnego aplikacji albo w komponencie pierwszego okna komunikacji — znaczniku stanu pętli `[Pętla ▼]` i przycisku „Pętla wykonawcza” w Chat Window. Rozpoczęcie zlecenia wieloetapowego zmienia wyłącznie stan znacznika pętli w Chat Window i nie otwiera okna samo z siebie. Po zamknięciu okno znika całkowicie z przestrzeni roboczej.

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. W module Assistant pętla obejmuje zlecenia wydane głosem i tekstem: Koordynator rozkłada wypowiedziane polecenie na zadania, przydziela je Wykonawcy, nadzoruje wywołania narzędzi i sterowanie modułami platformy, sprawdza wynik każdego zadania i decyduje o ponowieniu kroku nieudanego.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie głosowe lub tekstowe wraz z jego dekompozycją na zadania.
- Kolejka i stan zadań: oczekujące, w realizacji, zakończone, wymagające uwagi.
- Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą, z podaniem narzędzia lub modułu wykorzystywanego w zadaniu.
- Wyniki kontroli jakości każdego zadania i decyzje o ponowieniu.
- Wskaźniki przebiegu pętli: liczba iteracji, czas realizacji, liczba wywołań narzędzi, koszt kontekstu.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia w trakcie realizacji.
- Brama potwierdzeń zadań o skutkach ubocznych (wysyłka, publikacja, usunięcie) przed przekazaniem ich Wykonawcy.

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3

┌─ Execution Loop Window (kolumna sąsiadująca) ────────────────────┐
│ Zlecenie: „przygotuj notatkę ze spotkania i zapisz w bibliotece” │
│ Iteracja 2 · zadania 3 · wywołania narzędzi 5                    │
├──────────────────────────────────────────────────────────────────┤
│ Koordynator → Wykonawca                                          │
│   zadanie 1 · transkrypcja nagrania spotkania      ✓ przyjęte    │
│   zadanie 2 · redakcja notatki (Studio)            ▸ w realizacji│
│   zadanie 3 · zapis artefaktu (Library)            · oczekuje    │
├──────────────────────────────────────────────────────────────────┤
│ Wykonawca → Koordynator                                          │
│   „notatka zredagowana, oczekuję na potwierdzenie zapisu”        │
│   kontrola jakości: zgodna ✓                                     │
├──────────────────────────────────────────────────────────────────┤
│ [ Sterowanie ▼ ]                                                 │
└──────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok z treścią zlecenia i licznikami pętli | Identyfikacja realizowanego zlecenia i skali pętli | Średni blok tekstowy | zlecenie aktywne · zlecenie zakończone | Kliknięcie rozwija pełną treść polecenia i jego źródło (głos, tekst) | Nagłówek okna | 1 | Widoczny bez interakcji po otwarciu okna pętli |
| Lista zadań pętli | Lista pozycji z ikoną stanu | Prezentacja dekompozycji zlecenia | Średnia lista `.dn-listwa-pozycja` | oczekujące · w realizacji · zakończone · wymagające uwagi | Kliknięcie pozycji rozwija komunikaty i wynik zadania | Sekcja Koordynator → Wykonawca | 1 | Widoczna bez interakcji po otwarciu okna pętli |
| Strumień komunikatów sterujących | Zapis wymiany komunikatów | Wgląd w koordynację i nadzór realizacji | Średni blok przewijany | strumień aktywny · zatrzymany | — (odczytowy, aktualizowany na żywo) | Sekcja Wykonawca → Koordynator | 1 | Widoczny bez interakcji po otwarciu okna pętli |
| Wynik kontroli jakości | Plakietka przy zadaniu | Informacja o ocenie wyniku zadania | Mała plakietka `.dn-plakietka--sukces` / `.dn-plakietka--ostrzezenie` | zgodna · niezgodna (ponowienie) | Kliknięcie pokazuje kryteria i uzasadnienie decyzji o ponowieniu | Przy pozycji zadania | 2 | Plakietka stanu przy zadaniu; kliknięcie rozwija kryteria i uzasadnienie ponowienia |
| Wskaźniki przebiegu pętli | Rząd liczników | Pomiar iteracji, czasu, wywołań narzędzi i kosztu kontekstu | Małe liczniki tekstowe | aktualizowane na żywo | — (odczytowe) | Nagłówek okna | 1 | Widoczne bez interakcji w nagłówku okna |
| Przyciski sterowania przebiegiem | Rząd przycisków `.dn-btn` | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia | Średnie przyciski | dostępne · niedostępne (pętla zakończona) | Zmieniają stan pętli natychmiast, z zapisem w Activity Feed | Sekcja sterowania okna | 3 | Element zbiorczy `[ Sterowanie ▼ ]` w stopce okna; rozwinięcie zawiera wstrzymanie, wznowienie, przerwanie i korektę. Dostępne także poleceniem języka naturalnego |
| Brama potwierdzeń zadania | Blok `.dn-toast--ostrzezenie` z przyciskiem | Jawne potwierdzenie zadania o skutkach ubocznych | Mały blok z akcją | ukryta · oczekująca na decyzję | Potwierdzenie przekazuje zadanie Wykonawcy, odmowa usuwa je z kolejki | Przy zadaniu wymagającym potwierdzenia | 2 | Ujawnia się samoczynnie przy zadaniu o skutkach ubocznych i znika po decyzji |

---

### 3.3 Voice Console

| Aspekt | Wartość |
|---|---|
| Typologia | Okno główne — interfejs głosowy |
| Waga wizualna | Prawa kolumna, dominująca |
| Izolacja domyślna | Aktywna w trakcie interakcji głosowej; konfiguracja profilu trwała między sesjami |
| Zakres toru głosowego | Tor sesyjny: rozmowa wieloturowa, dyktowanie, makra i polecenia szybkie, profil głosu, pamięć i historia modułu. Krótkie polecenia ponadkontekstowe — nawigacja, odczyt stanu platformy, sterowanie procesem poza sesją — obsługuje tor globalny funkcji Always On Display ([Always On Display](../funkcje-globalne/always-on-display.md), rozdz. 5) |

**Zawartość i pełny arsenał funkcji.**

*Nagrywanie i rozpoznawanie.*
- Aktywacja nasłuchu przez przytrzymanie przycisku, frazę wybudzającą (konfigurowalną) lub tryb ciągłego nasłuchu sterowany ustawieniem konfiguracyjnym.
- Wizualizacja fali dźwiękowej na żywo podczas mówienia — potwierdzenie, że urządzenie odbiera głos.
- Wskaźnik pewności rozpoznania mowy oraz automatyczne wykrycie języka wypowiedzi dla użytkowników wielojęzycznych.
- Detekcja mowy kończąca nagranie po ciszy.
- Redukcja szumów otoczenia — przełącznik trybu „głośne otoczenie”.

*Odpowiedź głosowa.*
- Synteza mowy z wyborem głosu, tempa wypowiedzi i tonu (formalny, swobodny).
- Przerywanie odpowiedzi głosem („dość”, „zatrzymaj”) lub dotykiem w dowolnym momencie odtwarzania.
- Odsłuch ponowny i przewijanie długiej odpowiedzi głosowej.

*Polecenia i makra.*
- Biblioteka poleceń szybkich (voice shortcuts) — zdefiniowane przez użytkownika frazy wyzwalające złożone działania.
- Łańcuchowanie poleceń w jednej wypowiedzi („sprawdź pocztę i podsumuj najważniejsze wiadomości”).
- Polecenia kontekstowe odwołujące się do poprzedniej wypowiedzi („zrób to samo dla przyszłego tygodnia”).
- Siatka szybkich akcji — najczęściej używane polecenia uruchamiane dotykiem, bez mówienia.

*Sterowanie zadaniami i aplikacjami.*
- Wykonywanie poleceń wieloetapowych obejmujących inne moduły platformy (utworzenie dokumentu w Studio, wyszukanie informacji w Browser), inicjowanych głosem i prowadzonych w pętli wykonawczej.
- Obsługa aplikacji zewnętrznych i systemowych w zakresie ustanowionym przez konfigurację profilu.
- Potwierdzenie głosowe wykonania („zrobione”, „gotowe”) po zakończeniu działania.

*Dyktowanie.*
- Dyktowanie tekstu do aktywnego pola okna lub modułu platformy.
- Komendy formatujące głosem: „nowy akapit”, „wypunktuj”, „usuń ostatnie zdanie”, „duża litera”.

*Konfiguracja profilu dostępna z poziomu okna.*
- Wybór i podgląd charakterystyki głosu asystenta, imienia profilu, języka domyślnego.
- Skrót do pełnej konfiguracji Profilu asystenta na stronie głównej (strefa 2).

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3

┌─ Voice Console (prawa kolumna, dominująca) ───────────────────────┐
│ [Profil ▼] [Model ▼]                                        [ ☰ ] │
├───────────────────────────────────────────────────────────────────┤
│                                                                   │
│                  ┌───────────────────────┐                        │
│                  │  ((( wizualizacja ))) │  ◄ ujawnia się         │
│                  │    fali dźwiękowej    │    w trakcie mowy      │
│                  └───────────────────────┘                        │
│                                                                   │
│              [ 🎙  przytrzymaj i mów  ]                            │
│                                                                   │
│  transkrypcja na żywo …                                           │
│                                                                   │
├───────────────────────────────────────────────────────────────────┤
│ [ Szybkie akcje ▼ ]                                               │
└───────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Duży przycisk mikrofonu | Centralny okrągły przycisk | Główny punkt aktywacji nasłuchu głosowego | Bardzo duży przycisk, środek okna | spoczynek · nasłuch (pulsujący, akcent sygnału) · przetwarzanie (`.dn-spinner`) | Przytrzymanie lub kliknięcie rozpoczyna nagrywanie; ponowne kliknięcie kończy | Środek Voice Console | 1 | Widoczny bez interakcji w środku okna |
| Wizualizacja fali dźwiękowej | Animowany wykres amplitudy | Potwierdzenie odbioru głosu w czasie rzeczywistym | Średni element graficzny nad przyciskiem mikrofonu | ukryta (spoczynek) · aktywna (podczas mówienia) | — (odczytowa, reaguje na poziom głośności) | Nad przyciskiem mikrofonu | 1 | Ujawnia się samoczynnie z chwilą rozpoczęcia mowy |
| Transkrypcja na żywo | Linia tekstu przy przycisku | Podgląd rozpoznawanej mowy w czasie rzeczywistym | Średni tekst, wyśrodkowany | pusty · aktualizowany na żywo | — (odczytowa; po zakończeniu trafia do Chat Window) | Pod przyciskiem mikrofonu | 1 | Ujawnia się samoczynnie z chwilą rozpoznania mowy |
| Selektor profilu i języka | Para małych przycisków w nagłówku | Wybór aktywnego profilu i języka rozpoznawania | Małe przyciski z etykietą | domyślny · rozwinięty · auto-wykrycie języka | Zmiana wpływa na głos syntezy i model rozpoznawania | Nagłówek okna | 2 | Znaczniki kontekstowe `[Profil ▼]` i `[Model ▼]` w nagłówku; po wyborze zwijają się samoczynnie |
| Przełącznik „głośne otoczenie” | Mały przełącznik | Włączenie redukcji szumu i tłumienia echa | Mały przełącznik | wyłączony · włączony | Przełącza łańcuch przetwarzania audio na tor odszumiania | Sekcja transkrypcji | 2 | Menu hamburger (☰) nagłówka okna |
| Siatka szybkich akcji | Rząd małych przycisków | Wywołanie najczęstszych poleceń bez mówienia | Małe przyciski `.dn-btn--zarys` w rzędzie | domyślny · najechanie | Kliknięcie uruchamia predefiniowane polecenie tak, jakby zostało wypowiedziane | Sekcja akcji szybkich okna | 3 | Element zbiorczy `[ Szybkie akcje ▼ ]` w stopce okna |
| Przycisk „Ustawienia” | Rozwijane menu w nagłówku | Skrót do konfiguracji głosu, tempa, frazy wybudzającej | Mały przycisk z ikoną (ustawienia) | domyślny · rozwinięte | Otwiera panel szybkiej konfiguracji lub przenosi do pełnej konfiguracji profilu | Prawy róg nagłówka | 3 | Menu hamburger (☰) nagłówka; panel popover szybkiej konfiguracji głosu, tempa i frazy wybudzającej |
| Wskaźnik pewności rozpoznania | Mała plakietka przy transkrypcji | Sygnalizacja niepewności rozpoznanej frazy | Mała plakietka `.dn-plakietka--ostrzezenie` | ukryta (wysoka pewność) · widoczna (niska pewność) | Kliknięcie pozwala poprawić tekst ręcznie przed wysłaniem | Obok transkrypcji na żywo | 2 | Ujawnia się samoczynnie wyłącznie przy niskiej pewności rozpoznania |

---

### 3.4 Actions Monitor

| Aspekt | Wartość |
|---|---|
| Typologia | Okno monitorów i wskaźników |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Status realizacji właściwy bieżącej sesji, aktualizowany na żywo kanałem WebSocket |

**Zawartość i pełny arsenał funkcji.**

- Lista aktywnych i ostatnio zakończonych zleceń wieloetapowych z paskiem postępu każdego kroku.
- Szczegółowy podgląd kroku bieżącego — jaki moduł lub narzędzie jest wykorzystywane do realizacji polecenia.
- Sterowanie realizacją: wstrzymanie, wznowienie, anulowanie, priorytetyzacja i ponowienie kroku.
- Brama potwierdzeń akcji o skutkach ubocznych: wysyłka, płatność, usunięcie.
- Powiadomienie głosowe i wizualne o zakończeniu długotrwałego zlecenia, także gdy Voice Console nie jest aktywnie obserwowana.
- Historia błędów realizacji z możliwością ponowienia kroku, który się nie powiódł, wraz z korektą argumentów.
- Podgląd wyniku pośredniego każdego kroku (wygenerowany dokument, znaleziona informacja) bez opuszczania modułu.
- Filtrowanie listy: w toku, zakończone, nieudane, wymagające uwagi.

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3

┌─ Actions Monitor (kolumna boczna) ────────────────────────────────┐
│ [Filtr ▼]                                                         │
├───────────────────────────────────────────────────────────────────┤
│ ● „umów spotkanie z zespołem”              krok 2/3           [⋮] │
│   ▸ znaleziono wolny termin ✓                                     │
│   ▸ wysyłanie zaproszeń …                                         │
├───────────────────────────────────────────────────────────────────┤
│ ✓ „podsumuj pocztę z dziś”                 zakończone 09:41   [⋮] │
├───────────────────────────────────────────────────────────────────┤
│ ⚠ „zarezerwuj salę konferencyjną”          krok nieudany      [⋮] │
│   przyczyna: brak dostępnych sal                                  │
└───────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Filtr statusu | Rozwijana lista | Zawężenie listy zleceń wg statusu realizacji | Mały przycisk z etykietą | w toku (domyślny) · zakończone · nieudane · wszystkie | Przelicza widoczną listę zleceń | Nagłówek panelu | 2 | Znacznik kontekstowy `[Filtr ▼]` w nagłówku panelu; po wyborze zwija się samoczynnie |
| Pozycja zlecenia | Karta `.dn-karta` z paskiem postępu | Reprezentuje jedno zlecenie wieloetapowe | Średnia karta z kropką stanu `.dn-kropka` | w toku · zakończone (✓) · nieudane (⚠) · wstrzymane | Rozwija listę kroków i dostępne akcje sterujące | Lista główna panelu | 1 | Widoczna bez interakcji po otwarciu panelu |
| Pasek postępu kroków | Lista podpunktów z ikoną statusu | Szczegółowy widok etapów realizacji zlecenia | Mała lista zagnieżdżona w karcie zlecenia | ukończony krok (✓) · bieżący (animowany) · oczekujący | — (odczytowy, aktualizowany na żywo) | Wewnątrz karty zlecenia | 1 | Widoczny bez interakcji w karcie zlecenia |
| Przyciski sterujące (pauza, anuluj) | Małe ikony `.dn-btn-ikona` | Kontrola nad trwającym zleceniem | Bardzo małe ikony przy karcie | dostępne (zlecenie w toku) · niedostępne (zlecenie zakończone) | Wstrzymuje, wznawia lub anuluje realizację | Prawa krawędź karty zlecenia w toku | 3 | Menu kebab (⋮) karty zlecenia; dostępne także poleceniem języka naturalnego |
| Brama potwierdzenia akcji | Blok `.dn-toast--ostrzezenie` z przyciskiem | Jawne potwierdzenie akcji o skutkach ubocznych | Mały blok z akcją | ukryta · oczekująca na decyzję | Potwierdzenie uruchamia krok, odmowa anuluje go z zapisem przyczyny | Karta zlecenia oczekującego na decyzję | 2 | Ujawnia się samoczynnie przy akcji o skutkach ubocznych i znika po decyzji |
| Przycisk „Pokaż wynik” | Mały przycisk `--zarys` | Podgląd rezultatu zakończonego zlecenia | Mały przycisk tekstowy | domyślny | Otwiera podgląd wyniku (dokument, lista, podsumowanie) w nakładce | Karta zlecenia zakończonego | 3 | Menu kebab (⋮) karty zlecenia zakończonego |
| Dymek powiadomienia nieudanego kroku | Blok `.dn-toast--blad` | Informacja o przyczynie niepowodzenia | Mały blok z kreską lewą czerwoną | ukryty · widoczny | Przycisk „Ponów” uruchamia krok ponownie; „Szczegóły” pokazuje pełny log | Karta zlecenia nieudanego | 2 | Ujawnia się samoczynnie wyłącznie przy kroku nieudanym |

---

### 3.5 Activity Feed

| Aspekt | Wartość |
|---|---|
| Typologia | Okno monitorów i wskaźników |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Izolacja domyślna | Zapis chronologiczny właściwy bieżącej sesji, narastający w czasie |

**Zawartość i pełny arsenał funkcji.**

- Chronologiczny, przewijany zapis wszystkich zrealizowanych zleceń wraz z wynikiem końcowym.
- Grupowanie wg dnia i sesji dla orientacji w historii dłuższej pracy z asystentem.
- Wyszukiwanie pełnotekstowe i semantyczne w zapisie działań („kiedy poprosiłem o rezerwację sali”).
- Odtworzenie przebiegu wieloetapowego zlecenia krok po kroku, z powrotem do dowolnego wyniku pośredniego.
- Odsłuch nagrania audio oryginalnego polecenia głosowego, gdy przechowywanie nagrań jest włączone ustawieniem konfiguracyjnym.
- Wznowienie kontekstu z wybranego wpisu — powrót do stanu wcześniejszej rozmowy.
- Audyt wywołań narzędzi: jawny zapis każdego wywołania, jego argumentów i wyniku.
- Oznaczanie wpisów jako ważne.
- Eksport dziennika aktywności jako log tekstowy.
- Filtrowanie wg typu działania: rozmowa, dokument, wyszukiwanie, narzędzie, sterowanie aplikacją zewnętrzną.

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–3

┌─ Activity Feed (kolumna boczna) ─────────────────────────────────┐
│ [ 🔍 szukaj w historii ]              [Filtr ▼]            [ ☰ ] │
├──────────────────────────────────────────────────────────────────┤
│ Dziś                                                             │
│  09:41  „podsumuj pocztę z dziś”         wynik: 6 wiad.     [⋮]  │
│  10:02  „umów spotkanie z zespołem”      wynik: zaplanowano [⋮]  │
├──────────────────────────────────────────────────────────────────┤
│ Wczoraj                                                          │
│  16:30  „zarezerwuj salę konferencyjną”  wynik: nieudane    [⋮]  │
└──────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Pole wyszukiwania | Pole `.dn-pole-kontrolka` z ikoną (szukaj) | Odnalezienie wcześniejszego działania w historii | Małe pole, nagłówek panelu | puste · z wynikami | Filtruje listę na żywo, pełnotekstowo i semantycznie | Nagłówek panelu | 1 | Widoczne bez interakcji w nagłówku panelu |
| Filtr typu działania | Rozwijana lista | Zawężenie historii wg rodzaju zlecenia | Mały przycisk z etykietą | wszystkie (domyślny) · aktywny filtr | Przelicza widoczną listę wpisów | Nagłówek panelu | 2 | Znacznik kontekstowy `[Filtr ▼]` w nagłówku panelu; po wyborze zwija się samoczynnie |
| Nagłówek grupy dnia | Etykieta separatora | Grupowanie wpisów chronologicznie | Mały nagłówek tekstowy | domyślny | — (informacyjny) | Nad grupą wpisów danego dnia | 1 | Widoczny bez interakcji nad grupą wpisów |
| Pozycja działania | Wiersz z godziną, treścią polecenia i wynikiem | Reprezentuje jedno zrealizowane zlecenie | Mały wiersz `.dn-listwa-pozycja` | domyślny · ważne (przypięte, wyróżnienie sygnałowe) · nieudane (oznaczenie ostrzegawcze) | Kliknięcie rozwija szczegóły i dostępne akcje | Lista główna panelu | 1 | Widoczna bez interakcji na liście głównej |
| Przycisk „Odtwórz przebieg” | Mały przycisk `--zarys` | Przegląd wszystkich kroków zlecenia wieloetapowego po fakcie | Mały przycisk tekstowy | domyślny | Otwiera widok krok-po-kroku analogiczny do Execution Loop Window, w trybie historycznym | Przy pozycji działania wieloetapowego | 3 | Menu kebab (⋮) pozycji działania |
| Przycisk „Wznów kontekst” | Mały przycisk `--zarys` | Powrót do rozmowy i stanu z wybranego wpisu | Mały przycisk tekstowy | domyślny | Przywraca kontekst wpisu do Chat Window bieżącej karty | Przy pozycji działania | 3 | Menu kebab (⋮) pozycji działania; dostępny także poleceniem języka naturalnego |
| Ikona nagrania audio | Mała ikona głośnika | Odsłuch oryginalnego nagrania polecenia głosowego | Bardzo mała ikona (🔊) | niedostępna (zapis nagrań wyłączony) · dostępna | Kliknięcie odtwarza zapisane nagranie | Przy pozycji działania głosowego | 2 | Ujawnia się przy pozycji wyłącznie przy włączonym zapisie nagrań |
| Przycisk „Eksportuj dziennik” | Przycisk `--zarys`, pełna szerokość | Pobranie pełnego zapisu aktywności jako logu | Średni przycisk, sekcja końcowa panelu | domyślny · ładowanie | Generuje plik tekstowy dziennika do pobrania | Sekcja końcowa panelu | 4 | Menu hamburger (☰) panelu, wywoływacz poleceń, polecenie języka naturalnego |

---

### 3.6 Memory & Context Manager

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządzania zasobem |
| Waga wizualna | Kolumna boczna szeroka, otwierana jako rozszerzenie boczne; rozwijana do pełnego obszaru roboczego |
| Izolacja domyślna | Pamięć i konteksty przypisane do Profilu asystenta; zakres współdzielenia sterowany konfiguracją |

Okno zamyka obszar pamięci i kontekstów modułu i realizuje zasadę jawności: pamięć asystenta jest w całości widoczna, edytowalna i usuwalna przez użytkownika.

**Zawartość i pełny arsenał funkcji.**

- Jawny edytor pamięci długoterminowej — przegląd, edycja i usuwanie faktów o użytkowniku, z audytem zmian.
- Przeglądarka pamięci semantycznej z wyszukiwaniem po znaczeniu, nie po słowie kluczowym.
- Menedżer kontekstów przełączanych — nazwane zestawy pamięci i instrukcji („praca”, „dom”, „projekt X”) aktywowane jednym poleceniem.
- Baza wiedzy profilu (RAG): wgrywanie dokumentów, indeks, cytowania w odpowiedziach, powiązanie z modułem Library.
- Reguły retencji i wygaszania: co pamiętać trwale, co zapominać po sesji lub czasie (TTL), znaczniki wrażliwości, treści nigdy niezapisywane.
- Podgląd bieżącego okna kontekstu z przypinaniem, zwijaniem i kompresją długiej rozmowy do podsumowań.

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–4

┌─ Memory & Context Manager (kolumna boczna szeroka) ──────────────┐
│ [ Fakty ] [ Pamięć semantyczna ] [ Konteksty ] [ Baza wiedzy ] ☰ │
├──────────────────────────────────────────────────────────────────┤
│ Fakty o użytkowniku                          [ 🔍 szukaj ]       │
│  • „spotkania zespołu we wtorki o 10:00”   trwałe            [⋮] │
│  • „preferuje odpowiedzi zwięzłe”          trwałe            [⋮] │
│  • „projekt Atlas — termin 30.09”          TTL 60 d          [⋮] │
├──────────────────────────────────────────────────────────────────┤
│ [Kontekst ▼]  projekt Atlas                                      │
│ Okno kontekstu: 12 400 / 128 000 tokenów                         │
└──────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Zakładki obszarów pamięci | Rząd zakładek | Przełączanie między faktami, pamięcią semantyczną, kontekstami i bazą wiedzy | Średnie zakładki | aktywna · nieaktywna | Podmienia zawartość panelu głównego | Nagłówek okna | 1 | Widoczne bez interakcji po otwarciu okna |
| Wpis pamięci | Wiersz z treścią faktu i polityką retencji | Reprezentuje jeden zapamiętany fakt | Mały wiersz `.dn-listwa-pozycja` | trwały · z TTL · wrażliwy (oznaczenie) | Edycja zmienia treść, usunięcie kasuje fakt trwale, obie operacje trafiają do audytu | Lista faktów | 1 | Widoczny bez interakcji na liście faktów |
| Wyszukiwarka semantyczna | Pole `.dn-pole-kontrolka` z ikoną (szukaj) | Odnajdywanie zapamiętanych treści po znaczeniu | Małe pole | puste · z wynikami | Zwraca listę fragmentów rozmów i dokumentów wraz ze wskaźnikiem podobieństwa | Zakładka pamięci semantycznej | 1 | Widoczna bez interakcji w zakładce pamięci semantycznej |
| Selektor kontekstu | Lista wyboru pojedynczego | Aktywacja nazwanego zestawu pamięci i instrukcji | Średnia lista wyboru | kontekst aktywny · nieaktywny | Podmienia zestaw pamięci i prompt systemowy bieżącej karty sesji | Sekcja kontekstów | 2 | Znacznik kontekstowy `[Kontekst ▼]` paska kontekstu lub zakładka kontekstów; po wyborze zwija się samoczynnie |
| Wgrywarka dokumentów bazy wiedzy | Obszar przeciągnięcia pliku | Zasilenie odpowiedzi własnymi dokumentami | Średni obszar z ikoną | pusty · indeksowanie (`.dn-spinner`) · zaindeksowany | Ekstrahuje treść, buduje indeks i wektory, udostępnia cytowania | Zakładka bazy wiedzy | 2 | Zakładka bazy wiedzy; przeciągnięcie pliku na obszar roboczy |
| Reguły retencji | Formularz polityki pamięci | Ustalenie, co pamiętać, co wygaszać i po jakim czasie | Średni formularz | domyślny · zmieniony | Zapis reguły obejmuje kolejne zapisy pamięci profilu | Zakładka faktów | 4 | Menu hamburger (☰) zakładki faktów, polecenie języka naturalnego, konfiguracja Profilu asystenta |
| Miernik okna kontekstu | Pasek z licznikiem tokenów | Pomiar zajętości bieżącego kontekstu rozmowy | Mały pasek postępu | w normie · bliski wypełnienia (ostrzeżenie) | Przycisk „Skompresuj rozmowę” zamienia starsze wypowiedzi w podsumowanie | Sekcja końcowa okna | 1 | Widoczny bez interakcji w sekcji końcowej okna |

---

### 3.7 Command & Tools Hub

| Aspekt | Wartość |
|---|---|
| Typologia | Okno zarządzania zasobem |
| Waga wizualna | Kolumna boczna szeroka, otwierana jako rozszerzenie boczne; rozwijana do pełnego obszaru roboczego |
| Izolacja domyślna | Katalog narzędzi, makr i rutyn przypisany do Profilu asystenta, w granicach zakresu uprawnień profilu |

Okno zamyka obszary szybkich akcji, wywołań narzędzi oraz proaktywności modułu.

**Zawartość i pełny arsenał funkcji.**

- Wywoływacz poleceń i paleta poleceń (`Ctrl/Cmd + K`) wraz z konfiguracją siatki.
- Katalog makr głosowych i tekstowych z edytorem kroków (`JSON`/`YAML`), warunkami i harmonogramem.
- Słownik rozwijania skrótów tekstowych z polami szablonu.
- Menedżer schowka z historią, wyszukiwaniem i przypinaniem.
- Katalog narzędzi (function calling) i serwerów MCP z zakresem uprawnień per profil, argumentami i limitami.
- Katalog umiejętności profilu — nazwanych, wielokrokowych zdolności skomponowanych z narzędzi i instrukcji.
- Menedżer przypomnień, rutyn i wyzwalaczy zdarzeniowych wraz z briefem dziennym.
- Przekazanie powtarzalnego makra lub rutyny do modułu Automations przyciskiem „→ Wyślij do Automations”.

**Makieta tekstowa.**

```
Stan spoczynku — warstwa 1 oraz zwinięte wyzwalacze warstw 2–4

┌─ Command & Tools Hub (kolumna boczna szeroka) ───────────────────┐
│ [ Akcje ] [ Makra ] [ Skróty ] [ Schowek ] [ Narzędzia i MCP ]   │
│ [ Umiejętności ] [ Rutyny ]                                  [☰] │
├──────────────────────────────────────────────────────────────────┤
│ Narzędzia dostępne profilowi „Asystent biurowy”                  │
│  ⦿ kalendarz.utworz_wydarzenie   CalDAV     potwierdzenie: tak   │
│  ⦿ poczta.wyslij                 SMTP       potwierdzenie: tak   │
│  ⦿ browser.szukaj                Browser    potwierdzenie: nie   │
│  ⦿ serwer MCP: repozytorium-firmowe  stdio  [ Zakres ▼ ]         │
├──────────────────────────────────────────────────────────────────┤
│ Rutyny:  „poranny brief” 08:00 dni robocze                   [⋮] │
└──────────────────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Stany | Zachowanie po interakcji | Występowanie | Warstwa | Sposób wywołania |
|---|---|---|---|---|---|---|---|---|
| Zakładki obszarów | Rząd zakładek | Przełączanie między akcjami, makrami, skrótami, schowkiem, narzędziami, umiejętnościami i rutynami | Średnie zakładki | aktywna · nieaktywna | Podmienia zawartość panelu głównego | Nagłówek okna | 1 | Widoczne bez interakcji po otwarciu okna |
| Wywoływacz poleceń | Pole poleceń ze skrótem globalnym | Uruchomienie asystenta i polecenia z dowolnego miejsca | Średnie pole z listą podpowiedzi | zwinięty · rozwinięty z wynikami | Wpisane lub wypowiedziane polecenie uruchamia akcję z indeksu | Zakładka akcji | 2 | Skrót klawiszowy globalny; po wykonaniu polecenia pole zwija się samoczynnie |
| Edytor makra | Formularz kroków `JSON`/`YAML` | Definicja frazy wyzwalającej i sekwencji kroków | Duży formularz z podglądem | poprawne · błąd walidacji schematu | Zapis rejestruje makro jako polecenie szybkie profilu | Zakładka makr | 4 | Zakładka makr, przycisk edycji pozycji makra, wywoływacz poleceń. Edycja schematu `JSON`/`YAML` przeznaczona dla użytkownika zaawansowanego |
| Słownik skrótów tekstowych | Tabela skrót → treść | Rozwijanie krótkiego skrótu w dłuższy szablon | Średnia tabela edytowalna | domyślny · edycja wiersza | Zapis udostępnia skrót we wszystkich polach tekstowych platformy | Zakładka skrótów | 3 | Zakładka skrótów okna |
| Historia schowka | Lista skopiowanych treści | Odnalezienie i ponowne użycie wcześniejszej treści | Średnia lista z wyszukiwaniem | domyślny · przypięty wpis | Kliknięcie wkleja treść do aktywnego pola | Zakładka schowka | 3 | Zakładka schowka okna, skrót klawiszowy |
| Pozycja narzędzia lub serwera MCP | Wiersz z nazwą, transportem i polityką potwierdzeń | Zarządzanie dostępnością narzędzia dla profilu | Średni wiersz `.dn-listwa-pozycja` | włączone · wyłączone · wymagające potwierdzenia | Przełącznik zmienia dostępność, „Zakres” otwiera granice uprawnień i limity | Zakładka narzędzi i MCP | 3 | Zakładka narzędzi i MCP; element zbiorczy `[ Zakres ▼ ]` rozwija granice uprawnień i limity (warstwa 4) |
| Pozycja umiejętności | Karta z opisem zdolności i użytymi narzędziami | Włączenie wielokrokowej zdolności dla profilu | Średnia karta | włączona · wyłączona | Włączenie udostępnia umiejętność w pętli wykonawczej | Zakładka umiejętności | 3 | Zakładka umiejętności okna |
| Pozycja rutyny | Wiersz z wyzwalaczem czasowym lub zdarzeniowym | Uruchamianie sekwencji akcji o czasie lub po zdarzeniu | Średni wiersz z ikonami sterującymi | aktywna · wstrzymana | Edycja zmienia wyzwalacz i kroki, wstrzymanie zatrzymuje wyzwalanie | Zakładka rutyn | 3 | Zakładka rutyn okna |
| Przycisk „→ Wyślij do Automations” | Przycisk `--zarys` | Przekazanie makra lub rutyny jako procesu cyklicznego | Średni przycisk | domyślny · przekazano | Tworzy scenariusz w module Automations i zachowuje odnośnik zwrotny | Sekcja końcowa zakładek makr i rutyn | 4 | Menu kebab (⋮) pozycji makra lub rutyny, polecenie języka naturalnego |

---

## 4. Przepływy pracy

### 4.1 Przepływ podstawowy — polecenie głosowe wieloetapowe

```
 Voice Console      Chat Window          Execution Loop      Actions Monitor    Activity Feed
 ─────────────      ───────────          ──────────────      ───────────────    ─────────────
 1. aktywacja
    nasłuchu   ───►
 2. wypowiedź
    polecenia  ────────────►  rozpoznanie
                              i zapis
                              tekstowy
                                  │
                                  ▼
                          3. Koordynator dekomponuje ───►
                             zlecenie na zadania
                             i nadzoruje Wykonawcę
                                  │                     4. postęp kroków ───►
                                  ▼                        i wyniki pośrednie
 5. potwierdzenie  ◄──────────────┘
    głosowe wyniku                                                          │
                                                                            ▼
                                                                 6. zapis w chronologii
                                                                    po zakończeniu
```

### 4.2 Przepływ rozszerzony — polecenie obejmujące inny moduł platformy

```
Voice Console — „przygotuj notatkę ze spotkania i zapisz w bibliotece”
        │
        ▼
Execution Loop Window — dekompozycja na dwa zadania, kontrola jakości każdego
        │
        ├─ zadanie 1: utworzenie dokumentu ─────► Studio › Studio Editor
        │
        └─ zadanie 2: zapis artefaktu      ─────► Library › Library Explorer
        │
        ▼
Actions Monitor — postęp kroków i podgląd wyników pośrednich
        │
        ▼
Potwierdzenie głosowe: „notatka zapisana w bibliotece”
        │
        ▼
Activity Feed — wpis z odnośnikiem do utworzonego dokumentu
```

### 4.3 Przepływ rozszerzony — korekta błędnie rozpoznanego polecenia

```
Voice Console — wypowiedź rozpoznana z niską pewnością
        │
        ▼
Chat Window — transkrypcja oznaczona wskaźnikiem niskiej pewności
        │
        ▼
Użytkownik: „Popraw polecenie” — edycja tekstu ręcznie
        │
        ▼
Ponowne wysłanie poprawionego polecenia → Execution Loop Window → Actions Monitor
```

### 4.4 Przepływ rozszerzony — zadanie o skutkach ubocznych

```
Chat Window / Voice Console — polecenie „wyślij zaproszenia do zespołu”
        │
        ▼
Execution Loop Window — Koordynator oznacza zadanie jako nieodwracalne
        │
        ▼
Brama potwierdzeń — użytkownik zatwierdza wykonanie
        │
        ├─ zatwierdzone → Wykonawca realizuje zadanie
        └─ odrzucone    → zadanie usunięte z kolejki, przyczyna zapisana
        │
        ▼
Activity Feed — jawny zapis decyzji i wywołania narzędzia
```

### 4.5 Przepływ rozszerzony — praca z pamięcią i kontekstem

```
Memory & Context Manager — aktywacja kontekstu „projekt Atlas”
        │
        ▼
Chat Window — rozmowa prowadzona z pamięcią i instrukcjami tego kontekstu
        │
        ▼
Baza wiedzy profilu — odpowiedzi z cytowaniami wgranych dokumentów
        │
        ▼
Memory & Context Manager — nowy fakt zapisany z regułą retencji, widoczny i edytowalny
```

### 4.6 Przepływ rozszerzony — rozgraniczenie toru głosowego wobec Always On Display

Platforma prowadzi dwa rozłączne tory głosowe. Tor sesyjny modułu Assistant obsługuje rozmowę wieloturową, dyktowanie, makra, polecenia szybkie i profil głosu w obrębie karty sesji modułu. Tor globalny funkcji Always On Display obsługuje krótkie polecenia skierowane do platformy jako całości: nawigację między środowiskami i modułami, odczyt stanu procesów oraz sterowanie procesem poza sesją. Pełne rozgraniczenie, wraz z tabelą przypisania rodzajów poleceń i regułą pierwszeństwa, zawiera [Always On Display](../funkcje-globalne/always-on-display.md), rozdz. 5.

```
   Always On Display — tor globalny        │   ASSISTANT › Voice Console — tor sesyjny
   ────────────────────────────────────    │   ─────────────────────────────────────────
   krótkie polecenie ponadkontekstowe      │   rozmowa wieloturowa, dyktowanie,
   nawigacja, odczyt stanu platformy,      │   makra i polecenia szybkie,
   sterowanie procesem poza sesją          │   profil głosu, pamięć i historia modułu
   zapis: polecenie_glosowe bez profilu    │   zapis: polecenie_glosowe z profilem
             │                             │             ▲
             │  polecenie wymagające dyktowania,         │
             │  rozmowy wieloturowej albo makra          │
             └──────────── przekazanie ──────────────────┘
             ◄──────────── przekazanie zwrotne ──────────
                polecenie dotyczące platformy jako całości

   Przekazanie przenosi rozpoznaną treść i nasłuch; nie łączy kontekstów —
   tor globalny nie zapisuje rozmowy w pamięci modułu, tor sesyjny nie
   przejmuje kolejki sugestii funkcji globalnej.
```

Pierwszeństwo w obrębie modułu Assistant ma tor sesyjny: polecenie wydane przyciskiem mikrofonu okna Voice Console obsługuje zawsze moduł. Poza modułem Assistant polecenie głosowe przyjmuje tor globalny, przekazując je do okna Voice Console wtedy i tylko wtedy, gdy wykracza poza zakres globalny.

---

## 5. Stany, dane i powiązania

### 5.1 Model stanów zlecenia głosowego

```
┌────────────────┐  wypowiedź    ┌────────────────┐  dekompozycja  ┌─────────────────┐
│ Nasłuch        │──────────────►│ Rozpoznano     │  na zadania    │ Pętla           │
│ (aktywacja)    │               │ (transkrypcja) │───────────────►│ wykonawcza      │
└────────────────┘               └────────────────┘                └────────┬────────┘
                                                                            │
                                                           wynik pozytywny  │  błąd zadania
                                                    ┌───────────────────────┴─────────┐
                                                    ▼                                 ▼
                                           ┌────────────────┐               ┌──────────────────┐
                                           │ Zakończone     │               │ Wymaga uwagi     │
                                           │ (Activity Feed)│               │ (Actions Monitor)│
                                           └────────────────┘               └──────────────────┘
```

### 5.2 Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Assistant |
|---|---|
| `komponent_wlasny`, `profil_asystenta` | Definicja profilu: głos, tempo, język, uprawnienia, konteksty — tworzona w strefie 2 strony głównej |
| `sesja`, `karta_sesji` | Kontekst operacyjny, w którym profil asystenta jest aktywnie wykorzystywany |
| `wiadomosc` | Transkrypcje poleceń i odpowiedzi w Chat Window, drzewo wątków i wariantów |
| `polecenie_glosowe` | Pojedyncze polecenie głosowe wraz z wynikiem rozpoznania — nośnik danych Voice Console, Actions Monitor i Activity Feed |
| `zadanie`, `komponent_kompozycji` | Dekompozycja zlecenia na zadania widoczne w Execution Loop Window i Actions Monitor |
| `artefakt`, `etykieta_artefaktu` | Załączniki rozmowy, wyniki pośrednie zadań, notatki i transkrypcje przekazywane do Library |
| `prompt_systemowy` | Warstwowane instrukcje profilu, persony i kontekstu wstrzykiwane do rozmowy |
| `kanal_modelu` | Modele wykorzystywane do rozpoznawania mowy, rozumienia polecenia, embeddingów i syntezy odpowiedzi |
| `powiazanie_komponentu` | Powiązania z modułami platformy oraz pokrewieństwo Assistant ═══ Always On Display (bez łączenia kontekstu) |

### 5.3 Izolacja i konfigurowalność — punkty właściwe modułowi Assistant

| Punkt izolacji | Stan wyjściowy (domyślny) | Co można skonfigurować | Gdzie |
|---|---|---|---|
| Przechowywanie nagrań audio | Ustalane przez Operatora — brak ustawienia oznacza wartość domyślną platformy | Włączenie i wyłączenie zapisu nagrań w Activity Feed | Okno konfiguracji, ustawienie modułu |
| Historia i pamięć profilu asystenta | Odrębna per karta sesji | Współdzielenie historii poleceń między kartami korzystającymi z tego samego profilu | Okno konfiguracji, poziom „karta sesji” |
| Pamięć długoterminowa i semantyczna | Zapis jawny, widoczny w Memory & Context Manager | Zakres zapamiętywania, TTL, znaczniki wrażliwości, włączenie magazynu wektorowego | Okno konfiguracji, warstwa pamięci |
| Uprawnienia profilu do sterowania innymi modułami | Zgodnie z zakresem zdefiniowanym przy tworzeniu profilu (analogicznie do okna Permissions Center modułu Agents) | Rozszerzenie lub zawężenie zakresu modułów, narzędzi i serwerów MCP dostępnych profilowi | Konfiguracja Profilu asystenta, strefa 2 strony głównej |
| Bramy potwierdzeń akcji | Akcje nieodwracalne wymagają potwierdzenia | Wskazanie akcji objętych potwierdzeniem i zakresu ich rejestrowania | Okno konfiguracji, warstwa akcji |
| Izolacja techniczna procesu sesji (8 zakresów) | Żaden zakres domyślnie nie jest aktywny | Włączenie odrębnego konta i tokenu dla profilu wykonującego działania w imieniu użytkownika | Okno konfiguracji punktów izolacji, panel macierzy izolacji |

### 5.4 Powiązania z innymi modułami

```
                    ASSISTANT   ═══   Always On Display
        (pokrewieństwo funkcjonalne — wspólny charakter komunikacji
         głosowej; Assistant działa jako moduł osadzony w TalkIn,
         Always On Display jako warstwa ponad całą platformą)

    ASSISTANT — polecenia wieloetapowe obejmują dowolny moduł
    platformy jako zadanie pętli wykonawczej (Studio, Library, Browser,
    Research, Translate, Automations, Agents, Terminal), zgodnie
    z zakresem uprawnień profilu — nie jest to powiązanie stałe między
    modułami, lecz zdolność operacyjna profilu asystenta
```

| Integracja | Co daje | Kierunek | Zależności |
|---|---|---|---|
| **Assistant → Studio** | Utworzenie i redakcja dokumentu jako zadanie polecenia głosowego | jednokierunkowe, konfigurowalne | `powiazanie_komponentu`; `artefakt`; Actions Monitor |
| **Assistant → Library** | Zapis notatek, transkrypcji i artefaktów do trwałego repozytorium | jednokierunkowe | `artefakt`, `etykieta_artefaktu` |
| **Assistant ↔ Translate** | Tłumaczenie wypowiedzi i rozmowy dwustronnej | dwukierunkowe wywołania | Kanał Translate; wykrywanie języka |
| **Assistant → Browser** | Wyszukanie i otwarcie strony oraz przekazanie treści do rozmowy | jednokierunkowe | Integracja Browser; migawka treści |
| **Assistant → Research** | Zlecenie głębszego badania jako zadanie, z powrotem wyniku | jednokierunkowe | Sources/Findings; artefakt raportu |
| **Assistant → Terminal** | Uruchomienie kodu lub zapytania obliczeniowego na potrzeby odpowiedzi | jednokierunkowe | Piaskownica izolacyjna procesu sesji |
| **Assistant → Automations** | Zamiana powtarzalnego makra lub rutyny w proces cykliczny | jednokierunkowe | Scenariusz `JSON`/`YAML`; harmonogram |
| **Assistant → Agents** | Delegacja zadania autonomicznemu agentowi wg uprawnień profilu | jednokierunkowe | Zakres uprawnień; Agent Manager |
| **Assistant ═══ Always On Display** | Wspólny charakter komunikacji głosowej | pokrewieństwo funkcjonalne | `powiazanie_komponentu`, bez współdzielenia pamięci i kontekstu |

---

## 6. Scenariusze użycia

**Scenariusz 1 — polecenie w ruchu.**
Użytkownik, nie mając dostępu do klawiatury, aktywuje Voice Console i wydaje polecenie „umów spotkanie z zespołem na jutro na 10:00”. Execution Loop Window pokazuje dekompozycję na dwa zadania — znalezienie terminu i wysłanie zaproszeń — Actions Monitor prezentuje postęp, a Wykonawca potwierdza wykonanie głosowo.

**Scenariusz 2 — korekta błędnie zrozumianego polecenia.**
System rozpoznaje polecenie z niską pewnością, co Voice Console sygnalizuje wizualnie. Użytkownik poprawia transkrypcję ręcznie w Chat Window i wysyła skorygowaną wersję polecenia bez powtarzania całej wypowiedzi głosem.

**Scenariusz 3 — zlecenie obejmujące inny moduł.**
Użytkownik prosi asystenta o przygotowanie notatki ze spotkania i zapisanie jej w bibliotece. Execution Loop Window prowadzi pętlę dwóch zadań: utworzenie dokumentu (moduł Studio) i zapis artefaktu (moduł Library) — całość zainicjowana wyłącznie głosem.

**Scenariusz 4 — przegląd historii dłuższej współpracy.**
Po tygodniu korzystania z asystenta użytkownik przeszukuje Activity Feed frazą „rezerwacja sali”, odnajduje nieudaną próbę sprzed kilku dni, odtwarza przebieg zlecenia krok po kroku i ustala przyczynę niepowodzenia.

**Scenariusz 5 — wiele profili dla różnych kontekstów.**
Użytkownik konfiguruje na stronie głównej dwa profile asystenta: „Asystent biurowy” o formalnym tonie głosu i szerokich uprawnieniach do modułów pracy oraz „Asystent osobisty” o swobodniejszym tonie i węższym zakresie uprawnień, przełączając się między nimi z poziomu Voice Console zależnie od kontekstu.

**Scenariusz 6 — praca z własną bazą wiedzy.**
Użytkownik wgrywa do Memory & Context Manager zestaw dokumentów projektowych i aktywuje kontekst „projekt Atlas”. Kolejne odpowiedzi asystenta cytują konkretne fragmenty tych dokumentów, a ustalenia zapadające w rozmowie zapisują się jako fakty z regułą wygaszania po zakończeniu projektu.

**Scenariusz 7 — poranny brief jako rutyna.**
W Command & Tools Hub użytkownik definiuje rutynę „poranny brief” uruchamianą o 8:00 w dni robocze: asystent zestawia kalendarz, zadania i najważniejsze wiadomości oraz odczytuje podsumowanie głosem. Po miesiącu użytkownik przekazuje rutynę do modułu Automations jako proces cykliczny.

**Scenariusz 8 — zadanie nieodwracalne pod kontrolą.**
Polecenie „wyślij zaproszenia do zespołu” trafia do pętli wykonawczej, gdzie Koordynator oznacza zadanie wysyłki jako nieodwracalne. Brama potwierdzeń zatrzymuje realizację do momentu jawnej decyzji użytkownika, a decyzja i wywołanie narzędzia zapisują się w Activity Feed.

---

## 7. Katalog funkcji i narzędzi

Każda pozycja: nazwa, opis działania oraz zależności techniczne (biblioteki, formaty, integracje). Nazwy bibliotek Go wskazują realne rozwiązania implementacyjne rdzenia; dla warstwy głosu i modeli podano zarówno wariant lokalny, jak i kanał modelu (`kanal_modelu`).

### 7.1 Rozmowa i konwersacja (rdzeń)

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rozmowa wieloturowa głos+tekst** | Jeden ciągły wątek, w którym polecenie wydaje się głosem lub tekstem naprzemiennie, z zachowaniem kontekstu | `kanal_modelu`; `wiadomosc`; strumień WebSocket |
| **Strumień odpowiedzi z równoległą syntezą** | Odpowiedź pojawia się na żywo jako tekst i jednocześnie jest odczytywana głosem, z możliwością przerwania | Strumień tokenów; synteza mowy (rozdz. 7.2); synchronizacja tekst↔mowa |
| **Wątkowanie i rozgałęzianie rozmowy** | Odgałęzienie rozmowy od dowolnego punktu bez utraty oryginału; wiele równoległych nitek jednego tematu | Drzewo `wiadomosc`; `karta_sesji` |
| **Persony i tryby odpowiedzi** | Przełączane style: zwięzły, wyjaśniający, krok-po-kroku, ekspert, burza mózgów — bez zmiany profilu głosu | Prompt systemowy warstwowany; profil w konfiguracji |
| **Załączniki w rozmowie** | Wrzucenie pliku, obrazu lub zrzutu jako kontekstu wypowiedzi, multimodalnie | `artefakt`; wizja `kanal_modelu`; ekstrakcja tekstu (`ledongthuc/pdf`, `unidoc/unioffice`) |
| **Cytowanie i przypisy źródeł** | Odpowiedź odsyła do konkretnego fragmentu materiału wprowadzonego do rozmowy | Kotwice cytowań; mapowanie na `artefakt` i pamięć |
| **Regeneracja i warianty odpowiedzi** | Ponowne wygenerowanie odpowiedzi, warianty obok siebie, wybór najlepszego | `kanal_modelu`; wersjonowanie `wiadomosc` |
| **Tryb lokalny modelu** | Rozmowa na modelu lokalnym przy braku sieci lub przy wymogu prywatności | Lokalny runtime (`ollama`, `llama.cpp`); wybór w `kanal_modelu` |

### 7.2 Głos — rozpoznawanie, synteza, wybudzanie

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rozpoznawanie mowy (STT)** | Transkrypcja mowy na tekst na żywo, z interpunkcją i znacznikami czasu | whisper.cpp (`ggerganov/whisper.cpp` przez `mutablelogic/go-whisper`), Vosk (`alphacep/vosk-api` Go) lub STT `kanal_modelu`; wejście `WAV/PCM`, `Opus` |
| **Synteza mowy (TTS)** | Odczyt odpowiedzi z wyborem głosu, tempa i tonu; głosy naturalne i neuronowe | Piper (`rhasspy/piper`), Coqui TTS lub TTS `kanal_modelu`; wyjście audio `oto` (`hajimehoshi/oto`), `malgo` (`gen2brain/malgo`) |
| **Fraza wybudzająca (wake word)** | Aktywacja nasłuchu bez dotyku po wypowiedzeniu konfigurowalnej frazy | Porcupine (`Picovoice/porcupine` Go), openWakeWord; model frazy w konfiguracji profilu |
| **Detekcja mowy (VAD)** | Wykrycie początku i końca wypowiedzi — automatyczne kończenie nagrania po ciszy | WebRTC VAD (`maxhawkins/go-webrtcvad`), Silero VAD; strumień PCM |
| **Redukcja szumu i tryb „głośne otoczenie”** | Odszumianie i tłumienie echa dla czystej transkrypcji w hałasie | RNNoise (`xiph/rnnoise`), WebRTC APM; przełącznik w Voice Console |
| **Auto-wykrycie języka wypowiedzi** | Rozpoznanie języka mowy i przełączenie modelu STT/TTS bez ręcznego wyboru | `pemistahl/lingua-go`; wielojęzyczny model Whisper |
| **Wskaźnik pewności i korekta transkrypcji** | Podświetlenie niepewnych fragmentów i szybka korekta przed wysłaniem | Score STT; edycja inline w Chat Window |
| **Głos niestandardowy** | Profil asystenta mówi wybranym, spersonalizowanym głosem, za zgodą i przy jawnych kluczach | XTTS/Piper voice; próbki jako `artefakt`; zgoda w konfiguracji |
| **Odsłuch ponowny i przewijanie odpowiedzi** | Powtórzenie ostatniej wypowiedzi, przewijanie długiej odpowiedzi głosowej | Bufor audio sesji; sterowanie odtwarzaniem |

### 7.3 Dyktowanie i produktywność mowy

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Dyktowanie do dowolnego pola** | Mowa zamieniana na tekst wstawiany do aktywnego okna lub modułu platformy | Rozpoznawanie mowy (rozdz. 7.2); wstawianie do fokusu; komendy formatujące |
| **Komendy formatujące głosem** | „nowy akapit”, „wypunktuj”, „usuń ostatnie zdanie”, „duża litera” sterują tekstem podczas dyktowania | Parser komend; mapowanie na edycję tekstu |
| **Transkrypcja spotkań i nagrań** | Zapis i transkrypcja dłuższej sesji audio z podziałem na mówców, streszczeniem i listą zadań | Whisper; diaryzacja (pyannote); streszczenie `kanal_modelu` |
| **Notatki głosowe → strukturyzacja** | Swobodna dyktowana notatka porządkowana w listę zadań, punktów lub podsumowanie | STT + `kanal_modelu`; zapis jako `artefakt` do Library |

### 7.4 Pamięć i konteksty

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Pamięć długoterminowa (fakty o użytkowniku)** | Trwałe zapamiętywanie preferencji, imion i ustaleń przywoływanych automatycznie w kolejnych rozmowach | Magazyn faktów; `profil_asystenta`; wstrzykiwanie do promptu |
| **Pamięć semantyczna (wektorowa)** | Odnajdywanie i przywoływanie wcześniejszych rozmów i dokumentów po znaczeniu, nie po słowie kluczowym | Embeddingi (`kanal_modelu`); magazyn wektorowy `chromem-go` (`philippgille/chromem-go`), `sqlite-vec` lub `pgvector` |
| **Pamięć krótkoterminowa i okno kontekstu** | Zarządzanie bieżącym kontekstem: przypinanie, zwijanie, kompresja długiej rozmowy do podsumowań | Kompresja `kanal_modelu`; licznik tokenów |
| **Konteksty przełączane** | Nazwane zestawy pamięci i instrukcji („praca”, „dom”, „projekt X”) aktywowane jednym poleceniem | `komponent_wlasny`; zestaw pamięci + prompt; selektor w oknie |
| **Baza wiedzy profilu (RAG)** | Wgranie własnych dokumentów, które asystent przeszukuje i cytuje w odpowiedziach | Indeks Bleve (`blevesearch/bleve`) + wektory; ekstrakcja plików; powiązanie z Library |
| **Edytor pamięci (jawny wgląd)** | Przegląd, edycja i usuwanie zapamiętanych faktów — klucze jawne, żadnej ukrytej pamięci | Widok pamięci w Memory & Context Manager; audyt zmian |
| **Redakcja i wygaszanie pamięci (TTL)** | Reguły: co pamiętać trwale, co zapominać po sesji lub czasie, czego nigdy nie zapisywać | Polityka retencji w konfiguracji; znaczniki wrażliwości |
| **Wznowienie kontekstu międzysesyjnego** | Przywołanie poprzedniej sesji wraz z jej stanem — „wróćmy do tego, o czym mówiliśmy wczoraj” | `sesja`, `karta_sesji`; wyszukiwanie w Activity Feed |

### 7.5 Szybkie akcje, wywoływacz i schowek

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Wywoływacz poleceń** | Globalny pasek i skrót otwierający asystenta z dowolnego miejsca — polecenie wpisane lub wypowiedziane i wykonane | Skrót globalny klienta; paleta poleceń `.dn-paleta` (`Ctrl/Cmd + K`); indeks akcji |
| **Paleta szybkich akcji** | Siatka najczęstszych poleceń uruchamianych dotykiem lub głosem bez pełnej wypowiedzi | Katalog akcji; siatka `.dn-btn--zarys`; profil użytkownika |
| **Makra głosowe i tekstowe (voice shortcuts)** | Zdefiniowana fraza wyzwala złożone, wieloetapowe działanie („poranny brief”) | Definicja makra `JSON`/`YAML`; pętla narzędziowa; harmonogram |
| **Rozwijanie skrótów tekstowych (text expander)** | Krótki skrót rozwijany w dłuższy, konfigurowalny tekst lub szablon z polami | Słownik skrótów; wstawianie do fokusu; zmienne szablonu |
| **Menedżer schowka z historią** | Historia kopiowanych treści z wyszukiwaniem, przypinaniem i wklejaniem przez asystenta | Schowek klienta; indeks Bleve; polityka prywatności |
| **Szybkie przechwytywanie notatek** | Natychmiastowy zrzut notatki bez otwierania edytora, porządkowany później | `artefakt`; kolejka przechwytywania; powiązanie z Library |
| **Kalkulacje i konwersje w locie** | Obliczenia, przeliczniki jednostek i walut oraz działania na datach bez opuszczania paska poleceń | Ewaluator wyrażeń; kursy przez narzędzie (rozdz. 7.6) |
| **Odpowiedzi natychmiastowe** | Mikro-zadania tekstowe („napisz uprzejmą odmowę”, „skróć to”) zwracane od razu do wklejenia | `kanal_modelu`; wstawianie do fokusu |

### 7.6 Wywołania narzędzi i integracje

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Katalog narzędzi (function calling)** | Rejestr wywoływalnych funkcji, które model wybiera i uruchamia z argumentami wg schematu | Definicje `JSON Schema`; pętla obserwacja→wywołanie→wynik (Go) |
| **Klient Model Context Protocol (MCP)** | Podłączanie zewnętrznych serwerów narzędzi i danych standardem MCP — jednolity sposób rozszerzania | Klient MCP (Go); transport `stdio`, `HTTP+SSE`; rejestr serwerów |
| **Umiejętności profilu** | Nazwane, wielokrokowe zdolności skomponowane z narzędzi i instrukcji, włączane per profil | Definicja umiejętności (`JSON`/`YAML`); zakres uprawnień profilu |
| **Sterowanie modułami platformy** | Polecenie realizuje zadanie w Studio, Library, Browser, Research, Translate i pozostałych modułach w ramach uprawnień profilu | `powiazanie_komponentu`; API modułów; Execution Loop Window, Actions Monitor |
| **Obsługa aplikacji i systemu** | Wywołanie aplikacji zewnętrznych i systemowych w zakresie ustanowionym w konfiguracji profilu | Adaptery integracji; klucze jawne; zakres uprawnień |
| **Wyszukiwanie w sieci i pobieranie** | Narzędzie „szukaj / otwórz stronę” zasilające odpowiedź świeżymi danymi | Integracja z modułem **Browser**; `net/http` |
| **Kalendarz, poczta, zadania (konektory)** | Odczyt i zapis wydarzeń, wiadomości i zadań przez jawnie skonfigurowane konektory | Konektory `CalDAV`, `IMAP/SMTP`, API; potwierdzenie akcji nieodwracalnych |
| **Wykonawca kodu i narzędzia obliczeniowe** | Uruchomienie fragmentu kodu lub zapytania do danych na potrzeby jednej odpowiedzi | Integracja z modułem **Terminal**; piaskownica izolacyjna procesu sesji |
| **Zatwierdzanie akcji o skutkach ubocznych** | Wysyłka, płatność i usunięcie wymagają jawnego potwierdzenia użytkownika przed wykonaniem | Brama potwierdzeń akcji; log w Activity Feed |

> **Uwaga bezpieczeństwa (spójna z zasadami platformy):** narzędzia operują na *referencjach* do sekretów (klucze jawne, konfigurowalne). Assistant **nie wykonuje samodzielnie płatności, przelewów ani transakcji finansowych** i nie wprowadza haseł ani danych kart — akcje nieodwracalne (wysyłka, publikacja, usunięcie) wymagają jawnego potwierdzenia użytkownika. Potwierdzenie nie jest blokadą interfejsu, lecz świadomym krokiem realizacji, rejestrowanym jawnie.

### 7.7 Sterowanie zadaniami wieloetapowymi

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Dekompozycja polecenia na zadania** | Złożone zlecenie („przygotuj notatkę i zapisz w bibliotece”) planowane i wykonywane etapami w pętli wykonawczej | Planer `kanal_modelu`; `zadanie`, `komponent_kompozycji` |
| **Łańcuchowanie i polecenia kontekstowe** | Odwołania do poprzedniej wypowiedzi i łączenie poleceń („a potem zrób to samo dla przyszłego tygodnia”) | Pamięć krótkoterminowa; rozwiązywanie odniesień |
| **Podgląd i sterowanie realizacją** | Wstrzymanie, wznowienie, anulowanie, priorytetyzacja i ponowienie zadań na żywo | Execution Loop Window, Actions Monitor; kolejka rdzenia; WebSocket |
| **Wyniki pośrednie każdego zadania** | Podgląd artefaktu lub informacji z każdego etapu bez opuszczania modułu | `artefakt`; nakładka podglądu |
| **Kontrola jakości i ponawianie** | Ocena wyniku zadania, zapis przyczyny niepowodzenia i ponowienie kroku z korektą argumentów | Log błędów; polityka ponowień; wskaźniki pętli |
| **Potwierdzenie głosowe zakończenia** | Komunikat „zrobione” lub „gotowe” po realizacji, także przy nieaktywnym oknie | TTS; powiadomienia platformy |

### 7.8 Proaktywność, przypomnienia i rutyny

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Przypomnienia i alarmy** | Przypomnienia czasowe i kontekstowe z powiadomieniem („przypomnij mi jutro o 9”) | Harmonogram (`robfig/cron`); kanał powiadomień |
| **Rutyny i wyzwalacze zdarzeniowe** | Sekwencja akcji uruchamiana o czasie lub po zdarzeniu („o 8:00 przeczytaj brief”) | Reguły wyzwalaczy; integracja z **Automations** |
| **Brief dzienny** | Zestawienie kalendarza, zadań i wiadomości na start dnia, głosem lub tekstem | Konektory (rozdz. 7.6); szablon briefu; harmonogram |
| **Proaktywne podpowiedzi** | Sugestie asystenta na bazie kontekstu i pamięci, w zakresie ustalonym konfiguracją | Reguły proaktywności; polityka prywatności; jawna zgoda użytkownika |
| **Timery i stopery głosowe** | Timery, stopery i wielokrotne alarmy („ustaw minutnik na 10 minut”) | Harmonogram sesji; sygnał dźwiękowy, TTS |

### 7.9 Personalizacja profili asystenta

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Wiele profili asystenta** | Odrębne profile („biurowy” formalny, „osobisty” swobodny) z własnym głosem, tonem i uprawnieniami | `komponent_wlasny`, `profil_asystenta`; selektor w oknie |
| **Charakterystyka głosu i tonu** | Wybór głosu, tempa, tonu (formalny, swobodny), imienia i języka domyślnego profilu | TTS voice; atrybuty profilu w konfiguracji |
| **Instrukcje systemowe profilu** | Trwałe wytyczne stylu, roli i granic asystenta wstrzykiwane do każdej rozmowy | Prompt warstwowany (`prompt_systemowy`); warstwa profilu |
| **Zakres uprawnień profilu** | Które moduły i narzędzia profil wywołuje — zakres rozszerzalny i zawężalny, jak okno Permissions Center modułu Agents | Zakres uprawnień; macierz izolacji; klucze jawne |
| **Awatar i tożsamość profilu** | Wizualna tożsamość profilu (awatar, kolor akcentu) dla szybkiego rozpoznania | `profil_asystenta`; system wizualny `.dn-*` |

### 7.10 Wielojęzyczność, dostępność i tryby

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Rozmowa wielojęzyczna** | Rozumienie i odpowiadanie w wielu językach, tłumaczenie wypowiedzi w locie | `pemistahl/lingua-go`; integracja z **Translate** |
| **Tryb „tłumacz rozmowy”** | Dwustronne tłumaczenie mowy między dwoma językami w czasie rzeczywistym | STT + TTS dwujęzycznie; silnik Translate |
| **Tryb bez rąk** | Pełne sterowanie głosem dla dostępności; potwierdzenia dźwiękowe, sterowanie tempem odczytu | Fraza wybudzająca; VAD; komendy nawigacyjne |
| **Tryb cichy (tylko tekst)** | Wyłączenie głosu w otoczeniu wymagającym ciszy, pełna praca tekstowa | Przełącznik trybu; wyłączenie TTS |

### 7.11 Prywatność, bezpieczeństwo i higiena

| Funkcja | Co robi | Zależności |
|---|---|---|
| **Przechowywanie nagrań** | Jawny wybór, czy nagrania audio poleceń są zapisywane w Activity Feed | Ustawienie modułu; polityka retencji |
| **Izolacja procesu sesji** | Wydzielone konto, token i dostęp dla profilu działającego w imieniu użytkownika | Proces sesji rdzenia (Go); panel macierzy izolacji |
| **Referencje sekretów (klucze jawne)** | Odwołania do sekretów w warstwie sekretów platformy — bez ukrytych bram, bez wpisywania haseł przez asystenta | Warstwa sekretów; klucze jawne w konfiguracji |
| **Tryb lokalny dla treści wrażliwych** | Cała rozmowa i pamięć na modelu i magazynie lokalnym, bez ruchu sieciowego | Lokalny STT/TTS/LLM; magazyn lokalny |
| **Audyt działań asystenta** | Pełny, jawny log poleceń, wywołań narzędzi i wyników do wglądu i eksportu | Activity Feed; eksport logu tekstowego |

### 7.12 Zależności techniczne — zestawienie

| Obszar | Kluczowa zależność | Realne rozwiązania Go i formaty |
|---|---|---|
| Rozpoznawanie mowy (STT) | Whisper / Vosk | `ggerganov/whisper.cpp` (przez `mutablelogic/go-whisper`), `alphacep/vosk-api`; `WAV/PCM`, `Opus` |
| Synteza mowy (TTS) | Piper / Coqui / kanał modelu | `rhasspy/piper`, Coqui XTTS; wyjście `oto` (`hajimehoshi/oto`), `malgo` (`gen2brain/malgo`) |
| Fraza wybudzająca | Porcupine / openWakeWord | `Picovoice/porcupine` (Go), openWakeWord |
| Detekcja mowy (VAD) | WebRTC / Silero | `maxhawkins/go-webrtcvad`, Silero VAD |
| Redukcja szumu | RNNoise / WebRTC APM | `xiph/rnnoise` |
| Pamięć semantyczna | Embeddingi + wektory | `philippgille/chromem-go`, `sqlite-vec`, `pgvector` |
| Indeks i wyszukiwanie | Pełny tekst | `blevesearch/bleve` |
| Wykrywanie języka | Język wypowiedzi | `pemistahl/lingua-go`, wielojęzyczny Whisper |
| Wywołania narzędzi | Function calling / MCP | `JSON Schema`; klient MCP (`stdio`, `HTTP+SSE`) |
| Harmonogram | Przypomnienia i rutyny | `robfig/cron` |
| Ekstrakcja załączników | PDF / Office / OCR | `ledongthuc/pdf`, `unidoc/unioffice`, `otiai10/gosseract` (`pl+eng`) |
| Model lokalny | LLM lokalny | `ollama`, `llama.cpp` |
| Konektory produktywności | Kalendarz i poczta | `CalDAV`, `IMAP/SMTP`, API zewnętrzne |

---

## 8. Punkty sterowania z okna konfiguracji

Zgodnie z Modelem konfiguracji (warstwy: globalny → środowisko → projekt → sesja; izolacja w dedykowanym oknie o siedmiu poziomach zasięgu) Operator personalizuje moduł Assistant bez blokad — brak ustawienia oznacza wartość domyślną, a wartością domyślną jest wykonanie.

| Zakres | Co Operator personalizuje | Warstwa i miejsce |
|---|---|---|
| **Profil asystenta** | Głos, tempo, ton, imię, język domyślny, awatar, persona i tryb odpowiedzi | Komponent własny (strefa 2) |
| **Kanał modelu** | Model rozmowy, odrębny model STT/TTS/embeddingów, tryb lokalny wobec sieciowego | Zachowanie modeli (`kanal_modelu`) |
| **Aktywacja głosu** | Fraza wybudzająca, próg VAD, ciągły nasłuch wobec przytrzymania, redukcja szumu | Aplikacja / akcja |
| **Instrukcje systemowe** | Trwałe wytyczne stylu, roli i granic profilu | Prompt systemowy (warstwa profilu) |
| **Pamięć** | Co zapamiętywać trwale, TTL i wygaszanie, znaczniki wrażliwości, włączenie pamięci semantycznej | Pamięć / sesja |
| **Konteksty** | Definicje i kontekst domyślny; przełączanie automatyczne wg pory i miejsca | Projekt / sesja |
| **Baza wiedzy (RAG)** | Które dokumenty zasilają odpowiedzi, źródła cytowań, powiązanie z Library | Integracje / komponenty |
| **Narzędzia i MCP** | Które narzędzia i serwery MCP są dostępne profilowi, ich argumenty i limity | Rozszerzenia / integracje |
| **Zakres uprawnień profilu** | Które moduły i aplikacje profil steruje głosowo — rozszerzenie i zawężenie | Uprawnienia (jak okno Permissions Center modułu Agents) |
| **Bramy potwierdzeń** | Które akcje wymagają jawnego potwierdzenia (wysyłka, płatność, usunięcie) | Akcje / bezpieczeństwo |
| **Pętla wykonawcza** | Limit iteracji, próg kontroli jakości, polityka ponowień, zakres komunikatów sterujących w Execution Loop Window | Procesy / sesja |
| **Proaktywność** | Włączenie i zakres podpowiedzi, briefów dziennych, rutyn i wyzwalaczy | Procesy / akcje |
| **Nagrania i prywatność** | Zapis nagrań audio, tryb lokalny dla treści wrażliwych, retencja logów | Izolacja / prywatność |
| **Izolacja procesu** | Odrębne konto i token dla profilu, wydzielony dostęp sieciowy | Okno punktów izolacji (macierz, 7 poziomów) |
| **Szybkie akcje i skróty** | Skrót globalny wywoływacza poleceń, siatka akcji, słownik skrótów tekstowych, menedżer schowka | Aplikacja |
| **Powiązania modułów** | Włączenie i kierunek Assistant → Studio / Library / Browser / Research / Translate / Terminal / Automations / Agents | Komponenty (`powiazanie_komponentu`) |

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| Przytrzymanie spacji (lub gest dotykowy) | Aktywacja nasłuchu | Voice Console |
| Fraza wybudzająca (konfigurowalna) | Aktywacja nasłuchu bez dotyku | Voice Console |
| Skrót globalny wywoływacza poleceń | Otwarcie paska poleceń z dowolnego miejsca | Command & Tools Hub |
| `Esc` / komenda „zatrzymaj” | Przerwanie odtwarzania odpowiedzi | Voice Console |
| Ikona `dzwonek` | Powiadomienie o zakończeniu zlecenia | Actions Monitor |
| Ikona `zegar` | Znacznik czasu działania | Activity Feed |
| Ikona `ostrzezenie` | Zadanie nieudane / niska pewność rozpoznania | Execution Loop Window, Actions Monitor, Voice Console |
| Ikona `zatrzymaj` | Przerwanie pętli wykonawczej lub anulowanie zlecenia | Execution Loop Window, Actions Monitor |
| Ikona `odswiez` | Ponowienie nieudanego zadania | Execution Loop Window, Actions Monitor |
| Ikona `kosz` | Usunięcie faktu z pamięci długoterminowej | Memory & Context Manager |
| Ikona `pobierz` | Eksport dziennika aktywności | Activity Feed |

---

## 10. Komendy kontraktu obszaru `assistant`

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### 10.1 Obszar `assistant` — 4 komendy

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `assistant.action.status` | Zwraca stan zleceń asystenta i wykonuje sterowanie nimi z Actions Monitor | `actionId:string` (opc)<br>`windowId:string` (opc)<br>`control:AssistantActionControl` (opc)<br>`priority:int` (opc)<br>`confirmationNote:string` (opc) | `actions:AssistantAction[]` (wym)<br>`awaiting:AssistantAction[]` (opc) |
| `assistant.activity.flag` | Oznacza wpis dziennika asystenta jako ważny albo zdejmuje to oznaczenie. Dziś wyróżnienie wpisu nie miałoby gdzie zamieszkać: AssistantActivityEntry nie niesie takiego pola, więc okno pokazywałoby wyróżnienie, którego rdzeń nie pamięta | `entryId:string` (wym)<br>`important:bool` (wym)<br>`note:string` (opc) | `entry:AssistantActivityEntry` (wym) |
| `assistant.activity.list` | Zwraca chronologiczny dziennik działań asystenta | `windowId:string` (opc)<br>`actionId:string` (opc)<br>`limit:int` (opc)<br>`query:string` (opc)<br>`importantOnly:bool` (opc)<br>`kind:AssistantActivityKind` (opc)<br>`since:int64` (opc) | `entries:AssistantActivityEntry[]` (wym) |
| `assistant.voice.command` | Przyjmuje polecenie głosowe i prowadzi je przez rozpoznanie mowy, model i syntezę | `windowId:string` (wym)<br>`audioRef:string` (opc)<br>`transcript:string` (opc)<br>`profileId:string` (opc)<br>`speak:bool` (opc)<br>`contextId:string` (opc)<br>`confirmPolicy:bool` (opc) | `transcript:string` (wym)<br>`action:AssistantAction` (wym)<br>`messageId:string` (opc)<br>`speechRef:string` (opc) |

**Zdarzenia obszaru `assistant` — 1:**

| Zdarzenie | Przeznaczenie | Pola ładunku |
|---|---|---|
| `assistant.action.changed` | Zmiana stanu zlecenia asystenta | `change:ChangeKind` (wym)<br>`action:AssistantAction` (wym)<br>`actor:ActorKind` (opc)<br>`actorClientId:string` (opc) |

Razem w wykazie: **4 komendy** z 1 obszaru kontraktu.

---

## 11. Kryteria odbioru

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Profil asystenta konfigurowany w strefie 2 strony głównej jest wykorzystywany operacyjnie w środowisku TalkIn | Utworzenie Profilu asystenta i jego wywołanie w sesji TalkIn |
| Żaden przycisk akcji modułu nie występuje w stanie zablokowanym | Przegląd kontrolek CTA katalogów elementów rozdz. 3, 7 pod kątem obecności `disabled`/`not-allowed` |
| Wszystkie 4 komendy obszaru `assistant` mają pokrycie w interfejsie modułu albo jawne wskazanie miejsca wywołania | Zestawienie rozdz. 10 z katalogami elementów rozdz. 3, 7 |
| Brama potwierdzeń zadania o skutkach ubocznych ujawnia się samoczynnie i nie blokuje kolejki | Test: zadanie z deklarowanym skutkiem ubocznym trafia do bramy potwierdzeń, odmowa usuwa je z kolejki bez przerwania pozostałych zadań |
| Rozpoznanie mowy o niskiej pewności sygnalizuje wskaźnik, nie blokuje wysłania | Test: transkrypcja o niskiej pewności rozpoznania pozwala na ręczną poprawkę przed wysłaniem |
| Stany kontrolek z katalogów elementów są zaimplementowane w komplecie | Przegląd katalogów elementów rozdz. 3, 7 pod kątem kompletu stanów: domyślny, wskazanie kursorem, fokus, ładowanie, ostrzeżenie, błąd |

---

*Koniec dokumentu. Moduł Assistant — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
