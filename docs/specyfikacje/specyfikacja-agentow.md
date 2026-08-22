# Danaco Console — Specyfikacja agentów

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
| **Tytuł** | Specyfikacja agentów |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper · projektant |
| **Przeznaczenie** | Ustala model agentów platformy: definicję agenta, siedem komponentów jego definicji, Centrum uprawnień, dwa kanały komunikacji operacyjnej i pełny przebieg pętli wykonawczej, mechanizm przypisywania agentów do ról w środowisku MultitaskingAI i do projektów w module Workspace |
| **Zakres** | model agenta jako komponentu własnego, okna operacyjne modułu Agents i cykl życia agenta, Centrum uprawnień, model danych i kontrakty komunikacji agenta, scenariusze użycia |
| **Poza zakresem** | projekt interfejsu modułu Agents — makiety okien, katalogi elementów z warstwą widoczności i sposobem wywołania, stany okien, punkty sterowania z okna konfiguracji — [Moduł Agents](../moduly/agents.md) |
| **Dokument nadrzędny** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](../architektura/koncepcja-platformy.md) · [Architektura techniczna](../architektura/architektura.md) · [Model danych](../architektura/model-danych.md) · [Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) · [Moduł Agents](../moduly/agents.md) |
| **Prototypy odniesienia** | `design/05-okna/moduly/agents.html` |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `agent`, `subagent`, `advisor`) · [Koncepcja platformy](../architektura/koncepcja-platformy.md) rozdz. 11.15 · [Architektura techniczna](../architektura/architektura.md) rozdz. 9, 12, 14, 15 · [Model danych](../architektura/model-danych.md) rozdz. 12–13 |
| **Zasada nadrzędna** | Zakres działania agenta jest przedmiotem jawnej, świadomej konfiguracji Użytkownika w Centrum uprawnień — stanem wyjściowym jest pełny dostęp operacyjny, a nie ograniczenie |

Dokument specyfikuje moduł Agents oraz warstwę agentową platformy Danaco Console: sposób tworzenia agenta (model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć, uprawnienia), status agenta jako komponentu własnego działającego ponad środowiskami, dwa kanały komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca) — wraz z pełnym przebiegiem pętli wykonawczej, warstwy widoczności w zarządzaniu agentami, mechanizm przypisywania agentów do ról w środowisku MultitaskingAI i do projektów w module Workspace, szablon definicji agenta oraz Centrum uprawnień. Układ interfejsu jest pionowy — okna rozmieszczone są w kolumnach sąsiadujących poziomo.

**Rozgraniczenie wobec dokumentu modułowego.** Niniejszy dokument opisuje model agentów: definicję agenta, jego komponenty, role, cykl życia oraz pełny przebieg pętli wykonawczej. Projekt interfejsu modułu Agents — makiety okien, katalogi elementów z warstwą widoczności i sposobem wywołania, stany okien oraz punkty sterowania z okna konfiguracji — zawiera [Moduł Agents](../moduly/agents.md).

---

## Spis treści

1. [Streszczenie](#streszczenie)
2. [Wprowadzenie](#wprowadzenie)
3. [Agent w strukturze platformy](#1-agent-w-strukturze-platformy)
   - [1.1 Agent jako komponent własny](#11-agent-jako-komponent-własny)
   - [1.2 Agent jako komponent platformowy ponad środowiskami](#12-agent-jako-komponent-platformowy-ponad-środowiskami)
   - [1.3 Moduł Agents w macierzy dostępności modułów](#13-moduł-agents-w-macierzy-dostępności-modułów)
   - [1.4 Relacja modułu Agents do warstwy modułów i do MultitaskingAI](#14-relacja-modułu-agents-do-warstwy-modułów-i-do-multitaskingai)
4. [Dwa kanały komunikacji operacyjnej i pętla wykonawcza](#2-dwa-kanały-komunikacji-operacyjnej-i-pętla-wykonawcza)
   - [2.1 Rozmieszczenie kanałów w układzie interfejsu](#21-rozmieszczenie-kanałów-w-układzie-interfejsu)
   - [2.2 Chat Window — kanał Użytkownik ↔ Wykonawca](#22-chat-window--kanał-użytkownik--wykonawca)
   - [2.3 Execution Loop Window — kanał Koordynator ↔ Wykonawca](#23-execution-loop-window--kanał-koordynator--wykonawca)
   - [2.4 Przebieg pętli wykonawczej](#24-przebieg-pętli-wykonawczej)
   - [2.5 Kolejka i stan zadań](#25-kolejka-i-stan-zadań)
   - [2.6 Komunikaty sterujące](#26-komunikaty-sterujące)
   - [2.7 Kontrola jakości i decyzja o ponowieniu](#27-kontrola-jakości-i-decyzja-o-ponowieniu)
   - [2.8 Sterowanie przebiegiem zlecenia](#28-sterowanie-przebiegiem-zlecenia)
   - [2.9 Koordynacja wielu Wykonawców pracujących równolegle](#29-koordynacja-wielu-wykonawców-pracujących-równolegle)
5. [Warstwy widoczności w zarządzaniu agentami](#3-warstwy-widoczności-w-zarządzaniu-agentami)
   - [3.1 Wybór Wykonawcy jako znacznik kontekstowy](#31-wybór-wykonawcy-jako-znacznik-kontekstowy)
   - [3.2 Role, uprawnienia i limity zasobów jako funkcje warstwy 4](#32-role-uprawnienia-i-limity-zasobów-jako-funkcje-warstwy-4)
6. [Okna operacyjne modułu Agents i cykl życia agenta](#4-okna-operacyjne-modułu-agents-i-cykl-życia-agenta)
   - [4.1 Okna operacyjne i okna komunikacji](#41-okna-operacyjne-i-okna-komunikacji)
   - [4.2 Cykl życia agenta](#42-cykl-życia-agenta)
7. [Tworzenie agenta — komponenty definicji](#5-tworzenie-agenta--komponenty-definicji)
   - [5.1 Model bazowy i kanał modelu](#51-model-bazowy-i-kanał-modelu)
   - [5.2 Tożsamość](#52-tożsamość)
   - [5.3 Instrukcje systemowe](#53-instrukcje-systemowe)
   - [5.4 Umiejętności](#54-umiejętności)
   - [5.5 Wtyczki i konektory](#55-wtyczki-i-konektory)
   - [5.6 Pamięć](#56-pamięć)
   - [5.7 Uprawnienia](#57-uprawnienia)
   - [5.8 Zestawienie komponentów definicji](#58-zestawienie-komponentów-definicji)
8. [Centrum uprawnień (Permissions Center)](#6-centrum-uprawnień-permissions-center)
   - [6.1 Zasada wyjściowa: uwierzytelnienie a uprawnienia agenta](#61-zasada-wyjściowa-uwierzytelnienie-a-uprawnienia-agenta)
   - [6.2 Zakres konfigurowany w Centrum uprawnień](#62-zakres-konfigurowany-w-centrum-uprawnień)
   - [6.3 Relacja do okna konfiguracji punktów izolacji](#63-relacja-do-okna-konfiguracji-punktów-izolacji)
   - [6.4 Stan wyjściowy i macierz zakresów izolacji technicznej](#64-stan-wyjściowy-i-macierz-zakresów-izolacji-technicznej)
9. [Agent jako komponent własny ponad środowiskami](#7-agent-jako-komponent-własny-ponad-środowiskami)
   - [7.1 Wykorzystanie w TalkIn, WorkSpace i CodeStudio](#71-wykorzystanie-w-talkin-workspace-i-codestudio)
   - [7.2 Wykorzystanie w środowisku MultitaskingAI](#72-wykorzystanie-w-środowisku-multitaskingai)
10. [Przypisanie agenta do ról w środowisku MultitaskingAI](#8-przypisanie-agenta-do-ról-w-środowisku-multitaskingai)
   - [8.1 Cztery role i Subagent Network](#81-cztery-role-i-subagent-network)
   - [8.2 Mechanika przypisania](#82-mechanika-przypisania)
   - [8.3 Kanał modelu per rola](#83-kanał-modelu-per-rola)
   - [8.4 Relacje między Wykonawcami](#84-relacje-między-wykonawcami)
   - [8.5 Profil izolacji roli](#85-profil-izolacji-roli)
11. [Przypisanie agenta do projektów w module Workspace](#9-przypisanie-agenta-do-projektów-w-module-workspace)
   - [9.1 Agent Manager](#91-agent-manager)
   - [9.2 Zakres domyślny i rozszerzenie](#92-zakres-domyślny-i-rozszerzenie)
   - [9.3 Powiązanie z modelem danych](#93-powiązanie-z-modelem-danych)
12. [Szablon definicji agenta](#10-szablon-definicji-agenta)
13. [Model danych i kontrakty komunikacji agenta](#11-model-danych-i-kontrakty-komunikacji-agenta)
   - [11.1 Encja Agent](#111-encja-agent)
   - [11.2 Encja Rola w MultitaskingAI](#112-encja-rola-w-multitaskingai)
   - [11.3 Komunikaty warstwy komunikacji](#113-komunikaty-warstwy-komunikacji)
   - [11.4 Wykaz pełny komend obszarów `agent`, `subagent`, `advisor`](#114-wykaz-pełny-komend-obszarów-agent-subagent-advisor)
14. [Konfigurowalność i zasada braku twardych blokad](#12-konfigurowalność-i-zasada-braku-twardych-blokad)
15. [Scenariusze użycia](#13-scenariusze-użycia)
   - [13.1 Budowa agenta specjalistycznego i przypisanie do projektu](#131-budowa-agenta-specjalistycznego-i-przypisanie-do-projektu)
   - [13.2 Zespół agentów w rolach środowiska MultitaskingAI](#132-zespół-agentów-w-rolach-środowiska-multitaskingai)
   - [13.3 Agent z ograniczonymi uprawnieniami technicznymi dla pracy autonomicznej](#133-agent-z-ograniczonymi-uprawnieniami-technicznymi-dla-pracy-autonomicznej)
16. [Słowniczek pojęć](#14-słowniczek-pojęć)
17. [Kryteria odbioru](#15-kryteria-odbioru)
18. [Załącznik A. Przykład wypełnionej definicji agenta](#załącznik-a-przykład-wypełnionej-definicji-agenta)
19. [Załącznik B. Macierz elementów definicji agenta](#załącznik-b-macierz-elementów-definicji-agenta)

---

## Streszczenie

Agent jest komponentem własnym platformy Danaco Console — nazwaną, skonfigurowaną przez Użytkownika jednostką AI o określonej tożsamości, zestawie umiejętności i uprawnień — tworzonym i zarządzanym w module Agents. W odróżnieniu od pozostałych czternastu modułów platformy, moduł Agents nie tylko udostępnia własną przestrzeń roboczą, lecz wytwarza komponent działający ponad całą strukturą środowisk: raz utworzony agent jest Wykonawcą dostępnym w środowiskach TalkIn, WorkSpace i CodeStudio, przypisywalnym do projektów w module Workspace przez okno Agent Manager oraz przypisywalnym do ról wykonawczych w środowisku MultitaskingAI.

| Zagadnienie | Ustalenie |
|---|---|
| Czym jest agent | Komponent własny — nazwana, skonfigurowana jednostka AI o określonej tożsamości, umiejętnościach i uprawnieniach |
| Gdzie powstaje | Moduł Agents, okno Agent Builder — siedem komponentów definicji oraz Centrum uprawnień |
| Gdzie działa | Ponad środowiskami — Wykonawca w TalkIn, WorkSpace i CodeStudio; rola w środowisku MultitaskingAI; przypisanie do projektu w module Workspace |
| Czym sterują uprawnienia | Zakresem operacyjnym agenta (rozszerzenia, moduły i zasoby, izolacja techniczna, MultitaskingAI) — niezależnie od uwierzytelnienia Użytkownika |
| Stan wyjściowy uprawnień | Pełny dostęp operacyjny; każde zawężenie jest świadomą, jawną decyzją Użytkownika |

Praca Wykonawcy przebiega w dwóch kanałach komunikacji operacyjnej. Chat Window (Użytkownik ↔ Wykonawca) jest centralnym punktem pracy Użytkownika i podstawowym mechanizmem sterowania procesami platformy. Execution Loop Window (Koordynator ↔ Wykonawca) jest oknem pętli wykonawczej: przyjęcia zlecenia, dekompozycji na zadania, kolejki i stanu zadań, wymiany komunikatów sterujących, wykonania, kontroli jakości, decyzji o ponowieniu, zamknięcia zlecenia oraz sterowania przebiegiem. Oba okna zajmują sąsiadujące kolumny układu pionowego (rozdz. 2).

Definicja agenta obejmuje siedem komponentów budowanych w oknie Agent Builder: model bazowy wraz z kanałem połączenia, tożsamość, instrukcje systemowe, umiejętności, wtyczki i konektory, konfigurację pamięci oraz konfigurację uprawnień ustalaną w Centrum uprawnień (Permissions Center). Centrum uprawnień określa zakres operacyjny agenta niezależnie od uwierzytelnienia Użytkownika, które w architekturze platformy pozostaje jedynym mechanizmem kontroli dostępu do samej aplikacji. Zgodnie z zasadą nadrzędną pełnej konfigurowalności, żaden zakres uprawnień nie jest domyślnie zawężony — brak ustawienia w Centrum uprawnień oznacza pełny dostęp operacyjny agenta, a nie ograniczenie.

---

## Wprowadzenie

Niniejszy dokument rozwija rozdział 11.15 (moduł Agents) oraz powiązane fragmenty Koncepcji platformy do postaci specyfikacji operacyjnej, opierając się dodatkowo na dokumentach Architektury, Modelu danych i Kontraktów komunikacji. Tam, gdzie Koncepcja platformy ustala cel i zakres funkcjonalny modułu Agents, niniejszy dokument ustala sposób jego realizacji: kolejność i zawartość kroków tworzenia agenta, zawartość Centrum uprawnień oraz mechanikę przypisywania agenta do ról i projektów.

Poniższa tabela porządkuje podstawę merytoryczną dokumentu — dokumenty źródłowe, ich rozdziały oraz to, co każdy z nich wnosi do specyfikacji agentów.

| Dokument źródłowy | Rozdziały | Co wnosi do specyfikacji agentów |
|---|---|---|
| Koncepcja platformy | 11.15 | Moduł Agents — cel, okna operacyjne, charakterystyka, powiązania |
| Koncepcja platformy | 2.3, 7.2 | Agent jako komponent własny; strefa 2 strony głównej |
| Koncepcja platformy | 6 | Konfigurowalność zależności i punktów izolacji |
| Koncepcja platformy | 10 | Macierz dostępności modułów |
| Koncepcja platformy | 11.2 | Moduł Workspace — Agent Manager, przypisanie do projektu |
| Koncepcja platformy | 13 | Środowisko MultitaskingAI — role wykonawcze |
| Architektura | 9 | Integracja modeli; kanały API, CLI, SSH, HTTP |
| Architektura | 12 | Izolacja — osiem zakresów izolacji technicznej |
| Architektura | 13 | Konfiguracja — struktura warstwowa |
| Architektura | 14 | Rozszerzenia — źródła Danaco Plugin i Personal |
| Architektura | 15 | Uwierzytelnianie |
| Model danych | 12, 13 | Modele i agenci — encja Agent, encja Rola w MultitaskingAI |
| Kontrakty komunikacji | 10.15 | Komunikaty warstwy komunikacji dla modeli i agentów |

Wszędzie, gdzie niniejszy dokument doprecyzowuje mechanizm opisany w jednym z dokumentów źródłowych, odwołanie do właściwego rozdziału podano w tekście.

---

## 1. Agent w strukturze platformy

### 1.1. Agent jako komponent własny

Agent jest jednym z czterech rodzajów komponentu własnego przewidzianych przez Koncepcję platformy (rozdz. 2.3 i rozdz. 7.2). Jak każdy komponent własny, agent jest nazwanym wytworem konkretnego Użytkownika — konkretną, zindywidualizowaną instancją utworzoną i zapisaną w oknie Agent Builder, a po zapisaniu dostępną jako zasób Użytkownika w pamięci aplikacji.

| Komponent własny | Moduł | Zakres wykorzystania operacyjnego |
|---|---|---|
| Agent | Agents | Ponad wszystkimi środowiskami — Wykonawca w TalkIn, WorkSpace i CodeStudio; rola w środowisku MultitaskingAI |
| Automatyka | Automations | Procesy cykliczne; komponent wpinany w sesji modułu |
| Projekt | Workspace | Jednostka organizacji pracy właściwa jednemu modułowi |
| Profil asystenta | Assistant | Środowisko TalkIn |

Moduł Agents jest konfigurowany z dwóch punktów wejścia prowadzących do tego samego okna Agent Builder i tej samej definicji agenta.

| Punkt wejścia | Kontekst rozpoczęcia pracy |
|---|---|
| Strona główna, strefa 2 (kafle komponentów własnych, rozdz. 7.2 Koncepcji platformy) | Tworzenie i zarządzanie agentami jako takimi |
| Okno modułowe środowiska (rozdz. 1.3) | Tworzenie lub edycja agenta w toku bieżącej pracy nad zadaniem |

### 1.2. Agent jako komponent platformowy ponad środowiskami

Rozdział 1 Koncepcji platformy wymienia agentów jako piąty — obok środowisk, modułów, czatu i okien operacyjnych — fundamentalny element platformy, warstwę inteligencji działającą ponad całym ekosystemem. Charakterystyka modułu Agents (rozdz. 11.15 Koncepcji platformy) precyzuje: agenci są komponentem platformowym działającym ponad wszystkimi środowiskami.

| Środowisko | Tryb wykorzystania agenta | Mechanizm |
|---|---|---|
| TalkIn | Wykonawca wybierany w kontekście modułu | Wybór komponentu własnego w sesji |
| WorkSpace | Wykonawca; przypisanie do projektu | Agent Manager (moduł Workspace) |
| CodeStudio | Wykonawca wybierany w kontekście modułu | Wybór komponentu własnego w sesji |
| MultitaskingAI | Rola wykonawcza (po przypisaniu) | Sekcja Role panelu orkiestracji |

Ta cecha odróżnia agenta od pozostałych komponentów własnych, których zastosowanie operacyjne jest ograniczone do jednego kontekstu pracy. Zależność tę potwierdza dokument Modelu danych: agent jest komponentem platformowym ponad wszystkimi środowiskami (rozdz. 12.1), a w rozdziale 13.2 tego dokumentu stwierdzono wprost — agent jest niezależny od środowiska i przypisywany jako Wykonawca.

### 1.3. Moduł Agents w macierzy dostępności modułów

Macierz dostępności modułów (rozdz. 10 Koncepcji platformy) wskazuje moduł Agents jako jeden z trzech modułów dostępnych we wszystkich trzech środowiskach modułowych.

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|---|---|---|
| Agents (11.15) | TAK | TAK | TAK |
| Workspace (11.2) | TAK | TAK | TAK |
| Roundtable (11.8) | TAK | TAK | TAK |

Pełna dostępność modułu Agents we wszystkich trzech środowiskach modułowych jest bezpośrednią konsekwencją jego charakteru: skoro agent ma być wykorzystywany jako Wykonawca niezależnie od kontekstu pracy, okno modułowe pozwalające go utworzyć lub zmodyfikować musi być dostępne w każdym środowisku, w którym taki Wykonawca może zostać powołany do pracy. Środowisko MultitaskingAI nie ma bocznej nawigacji modułów (rozdz. 9.4 i rozdz. 10 Koncepcji platformy) — agenci nie są w nim dostępni przez okno modułowe, lecz przez sekcję Role panelu orkiestracji (rozdz. 8 niniejszego dokumentu).

### 1.4. Relacja modułu Agents do warstwy modułów i do MultitaskingAI

Moduł Agents należy do warstwy modułów (Workspace Layer) opisanej w rozdziale 5.2 Koncepcji platformy, nie jest jednak modułem w pełni analogicznym do pozostałych czternastu — jego wytwór przekracza granice warstwy modułów.

| Aspekt | Ustalenie |
|---|---|
| Przynależność warstwowa | Warstwa modułów (Workspace Layer) — mechanizm przełączania i przeładowania, własny zestaw okien operacyjnych (rozdz. 4 niniejszego dokumentu) |
| Zasięg wytworu (agenta) | Ponad warstwą modułów — agent jest wykorzystywany jak element o zasięgu platformowym |
| Status funkcji globalnej | Agent nie jest funkcją globalną w rozumieniu rozdziału 5.3 Koncepcji platformy (nie ma statusu Mobile ani Always On Display); jest komponentem własnym o zasięgu platformowym (rozdz. 1.2) |
| Relacja do MultitaskingAI | Agent pełni role wykonawcze (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) po jawnym przypisaniu w sekcji Role panelu orkiestracji — powiązanie jawne i konfigurowalne (rozdz. 6.2 Koncepcji platformy) |

---

## 2. Dwa kanały komunikacji operacyjnej i pętla wykonawcza

Warstwa agentowa platformy Danaco Console pracuje w dwóch kanałach komunikacji operacyjnej. Oba kanały są elementami pierwszoplanowymi architektury modułu Agents i każdego środowiska, w którym Wykonawca realizuje zadania.

| Kanał | Okno | Uczestnicy | Przedmiot wymiany |
|---|---|---|---|
| Pierwszy | Chat Window | Użytkownik ↔ Wykonawca | Zlecenia w języku naturalnym, strumień odpowiedzi i wyników, zatwierdzenia, przerwania, wyjaśnienia wyniku i kontekstu |
| Drugi | Execution Loop Window | Koordynator ↔ Wykonawca | Zlecenie i jego dekompozycja, kolejka i stan zadań, komunikaty sterujące, wyniki kontroli jakości, decyzje o ponowieniu, sterowanie przebiegiem |

Role uczestników wymiany rozdzielone są jednoznacznie.

| Rola | Zakres |
|---|---|
| Użytkownik | Zleca pracę i zatwierdza jej wyniki |
| Koordynator | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela zadania, nadzoruje ich realizację i zamyka zlecenie |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania przydzielone przez Koordynatora |

### 2.1. Rozmieszczenie kanałów w układzie interfejsu

Układ interfejsu jest pionowy — okna rozmieszczone są w kolumnach sąsiadujących poziomo. Chat Window zajmuje lewą kolumnę, stałą, o pełnej wysokości obszaru roboczego. Execution Loop Window otwierane jest jako kolumna sąsiadująca z Chat Window. Obszar roboczy modułu Agents — Agent Builder wraz z czterema oknami konfiguracyjnymi — zajmuje kolumnę dominującą po prawej stronie. Okna pomocnicze otwierane są jako rozszerzenia boczne, w kolejnych kolumnach po prawej stronie obszaru roboczego. Regulacji podlega wyłącznie szerokość kolumn.

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │ Agents                   │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │                      │  Agent Builder           │  boczne)
            │ [Agent ▼] [Model ▼]  │  Model Configuration     │
            │                      │  Skills Manager          │  Podgląd
            │ ─────────────────    │  Connectors Manager      │  definicji
            │ Execution Loop       │  Permissions Center      │  agenta
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                     [⋮]  │
            │                      │                          │
 ════════════════════════════════════════════════════════════════════════════
```

### 2.2. Chat Window — kanał Użytkownik ↔ Wykonawca

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą, centralnym punktem pracy Użytkownika i podstawowym mechanizmem sterowania wszystkimi procesami realizowanymi przez platformę. Okno występuje w tym samym miejscu układu w każdym module i w każdym środowisku — w lewej kolumnie obszaru roboczego — i należy do warstwy widoczności 1.

| Funkcja Chat Window | Realizacja w module Agents |
|---|---|
| Przyjmowanie poleceń w języku naturalnym | Zlecenie utworzenia agenta, zmiany jego definicji, przypisania do projektu lub do roli |
| Prezentacja strumienia odpowiedzi i wyników | Podgląd budowanej definicji, wyniki testowego uruchomienia agenta, komunikaty o zapisaniu komponentu własnego |
| Zatwierdzanie działań | Zatwierdzenie zapisu definicji agenta, zawężenia uprawnień, przypisania agenta do roli wykonawczej |
| Przerywanie działań | Przerwanie trwającego zlecenia realizowanego przez Wykonawcę |
| Wyjaśnianie wyniku i kontekstu | Wyjaśnienie zakresu uprawnień agenta, obowiązującej polityki izolacji i przyczyny ponowienia zadania |
| Wywoływanie funkcji warstwy 4 | Definiowanie ról, uprawnień i limitów zasobów poleceniem języka naturalnego (rozdz. 3) |

Sterowanie procesami z Chat Window obejmuje cały cykl życia zlecenia: złożenie zlecenia, wgląd w jego stan, wstrzymanie, wznowienie, przerwanie oraz korektę treści zlecenia w toku wykonania. Polecenia wydane w tym oknie są przekazywane Koordynatorowi i ujawniają się w Execution Loop Window jako komunikaty sterujące.

### 2.3. Execution Loop Window — kanał Koordynator ↔ Wykonawca

Execution Loop Window jest oknem pętli wykonawczej. Prezentuje komunikację między Koordynatorem a Wykonawcą i odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Okno otwierane jest jako kolumna sąsiadująca z Chat Window i należy do warstwy widoczności 1 w chwili, gdy zlecenie jest realizowane.

| Sekcja okna | Zawartość |
|---|---|
| Zlecenie bieżące | Treść zlecenia przyjętego od Użytkownika wraz z identyfikatorem, stanem i przypisanym Wykonawcą lub zespołem Wykonawców |
| Dekompozycja | Rozbicie zlecenia na zadania, zależności między zadaniami, kryteria uznania zadania za zakończone |
| Kolejka zadań | Zadania oczekujące, przydzielone, w realizacji, wstrzymane, zakończone i odrzucone wraz z ich stanem |
| Wymiana komunikatów sterujących | Strumień komunikatów Koordynator → Wykonawca oraz Wykonawca → Koordynator |
| Kontrola jakości | Wyniki oceny rezultatów zadań oraz decyzje o ponowieniu, przekazaniu dalej lub zamknięciu |
| Wskaźniki przebiegu pętli | Liczba zadań w każdym stanie, liczba ponowień, czas realizacji, obciążenie poszczególnych Wykonawców |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie oraz korekta zlecenia |

### 2.4. Przebieg pętli wykonawczej

Pętla wykonawcza obejmuje osiem faz — od przyjęcia zlecenia do jego zamknięcia. Faza kontroli jakości zamyka obieg: jej wynik prowadzi albo do ponowienia zadania, albo do zamknięcia zlecenia.

| Faza | Uczestnik prowadzący | Przebieg | Wynik fazy |
|---|---|---|---|
| 1. Przyjęcie zlecenia | Koordynator | Zlecenie złożone przez Użytkownika w Chat Window trafia do Koordynatora, który potwierdza jego przyjęcie i nadaje mu identyfikator | Zlecenie w stanie „przyjęte” |
| 2. Dekompozycja na zadania | Koordynator | Zlecenie zostaje rozbite na zadania elementarne wraz z zależnościami i kryteriami zakończenia | Zbiór zadań w stanie „oczekujące” |
| 3. Przydział zadań | Koordynator | Każde zadanie otrzymuje Wykonawcę odpowiadającego jego charakterowi, zakresowi uprawnień i dostępności | Zadania w stanie „przydzielone” |
| 4. Wymiana komunikatów sterujących | Koordynator ↔ Wykonawca | Koordynator przekazuje treść zadania, kontekst i kryteria; Wykonawca potwierdza przyjęcie, zgłasza postęp, pytania i przeszkody | Uzgodniony zakres wykonania |
| 5. Wykonanie | Wykonawca | Realizacja zadania z wykorzystaniem umiejętności, wtyczek, konektorów i pamięci zapisanych w definicji agenta, w granicach uprawnień z Permissions Center | Rezultat zadania |
| 6. Kontrola jakości | Koordynator, Wykonawca w roli kontrolnej | Rezultat oceniany jest względem kryteriów zakończenia ustalonych w fazie dekompozycji | Ocena pozytywna albo negatywna |
| 7. Decyzja o ponowieniu | Koordynator | Ocena negatywna kieruje zadanie do ponownego wykonania z uzupełnionym kontekstem; przekroczenie ustalonej liczby ponowień eskaluje sprawę do Użytkownika w Chat Window | Zadanie w stanie „ponowione” albo eskalacja |
| 8. Zamknięcie zlecenia | Koordynator | Po pozytywnej ocenie wszystkich zadań rezultaty zostają scalone i przekazane Użytkownikowi do zatwierdzenia w Chat Window | Zlecenie w stanie „zamknięte” |

Poniższy schemat przedstawia przebieg pętli w formie tekstowej.

```
  UŻYTKOWNIK ──(zlecenie, Chat Window)──►  KOORDYNATOR
                                                │
                                                ▼
                                     1. przyjęcie zlecenia
                                     2. dekompozycja na zadania
                                     3. przydział zadań
                                                │
                                                ▼
                                     ┌─── KOLEJKA ZADAŃ ───┐
                                     │ oczekujące          │
                                     │ przydzielone        │
                                     │ w realizacji        │
                                     │ wstrzymane          │
                                     │ zakończone          │
                                     └──────────┬──────────┘
                                                │
                     4. komunikaty sterujące    │
        KOORDYNATOR ◄──────────────────────────►│──► WYKONAWCA 1
                                                │──► WYKONAWCA 2
                                                │──► WYKONAWCA 3
                                                │
                                                ▼
                                        5. wykonanie
                                                │
                                                ▼
                                        6. kontrola jakości
                                          │             │
                           ocena negatywna│             │ocena pozytywna
                                          ▼             ▼
                                7. ponowienie      8. zamknięcie zlecenia
                                          │             │
                                          └──►(kolejka) ▼
                                                  UŻYTKOWNIK
                                            (zatwierdzenie, Chat Window)
```

### 2.5. Kolejka i stan zadań

Kolejka zadań jest podstawą pracy pętli wykonawczej. Każde zadanie ma w każdej chwili dokładnie jeden stan; przejścia między stanami zapisywane są w Execution Loop Window.

| Stan zadania | Znaczenie | Przejście następne |
|---|---|---|
| Oczekujące | Zadanie powstało w dekompozycji i czeka na przydział lub na spełnienie zależności | Przydzielone |
| Przydzielone | Zadanie ma wskazanego Wykonawcę i przekazany kontekst | W realizacji |
| W realizacji | Wykonawca realizuje zadanie | Zakończone, wstrzymane albo odrzucone |
| Wstrzymane | Przebieg wstrzymany poleceniem Użytkownika lub decyzją Koordynatora | W realizacji albo odrzucone |
| Zakończone | Rezultat przekazany do kontroli jakości | Zatwierdzone albo ponowione |
| Ponowione | Ocena negatywna; zadanie wraca do kolejki z uzupełnionym kontekstem | Przydzielone |
| Zatwierdzone | Rezultat spełnia kryteria zakończenia | Wejście do scalenia rezultatów |
| Odrzucone | Zadanie przerwane lub usunięte w wyniku korekty zlecenia | Stan końcowy |

### 2.6. Komunikaty sterujące

Wymiana w kanale Koordynator ↔ Wykonawca opiera się na komunikatach sterujących o rozłącznych funkcjach.

| Kierunek | Komunikat | Funkcja |
|---|---|---|
| Koordynator → Wykonawca | Przydział zadania | Przekazanie treści zadania, kontekstu i kryteriów zakończenia |
| Koordynator → Wykonawca | Korekta kontekstu | Uzupełnienie lub zmiana kontekstu zadania w toku realizacji |
| Koordynator → Wykonawca | Polecenie ponowienia | Skierowanie zadania do ponownego wykonania wraz z uzasadnieniem oceny negatywnej |
| Koordynator → Wykonawca | Polecenie wstrzymania, wznowienia, przerwania | Sterowanie przebiegiem pojedynczego zadania lub całego zlecenia |
| Wykonawca → Koordynator | Potwierdzenie przyjęcia | Potwierdzenie objęcia zadania realizacją |
| Wykonawca → Koordynator | Raport postępu | Bieżący stan realizacji zadania |
| Wykonawca → Koordynator | Zgłoszenie przeszkody | Brak uprawnienia, brak danych, konflikt zależności |
| Wykonawca → Koordynator | Przekazanie rezultatu | Przekazanie wyniku zadania do kontroli jakości |

Komunikaty warstwy komunikacji, którymi wymiana ta jest realizowana technicznie, opisano w rozdziale 11.3.

### 2.7. Kontrola jakości i decyzja o ponowieniu

Kontrola jakości jest odrębną fazą pętli, oddzieloną od wykonania. Ocenę prowadzi Koordynator albo Wykonawca pełniący funkcję kontrolną — w środowisku MultitaskingAI rolę tę pełni Executor 3 / Validator (rozdz. 8.1).

| Element kontroli jakości | Ustalenie |
|---|---|
| Podstawa oceny | Kryteria zakończenia zadania ustalone w fazie dekompozycji |
| Wynik oceny | Ocena pozytywna — rezultat przechodzi do scalenia; ocena negatywna — zadanie wraca do kolejki |
| Zakres ponowienia | Ponowieniu podlega pojedyncze zadanie, nie całe zlecenie |
| Kontekst ponowienia | Uzasadnienie oceny negatywnej dołączane jest do zadania jako uzupełnienie kontekstu |
| Granica ponowień | Liczba ponowień zadania jest ustawieniem konfiguracyjnym; jej przekroczenie eskaluje sprawę do Użytkownika w Chat Window |

### 2.8. Sterowanie przebiegiem zlecenia

Użytkownik steruje przebiegiem pętli z Chat Window, a stan sterowania odzwierciedla Execution Loop Window. Te same operacje dostępne są bezpośrednio w oknie pętli wykonawczej.

| Operacja | Skutek dla pętli | Skutek dla zadań |
|---|---|---|
| Wstrzymanie | Pętla zatrzymuje przydzielanie kolejnych zadań | Zadania w realizacji przechodzą w stan „wstrzymane”; kolejka zachowuje zawartość |
| Wznowienie | Pętla podejmuje przydzielanie zadań od miejsca wstrzymania | Zadania wstrzymane wracają do stanu „w realizacji” |
| Przerwanie | Pętla kończy pracę nad zleceniem | Zadania niezakończone przechodzą w stan „odrzucone”; rezultaty zatwierdzone pozostają zachowane |
| Korekta zlecenia | Koordynator powtarza dekompozycję dla zmienionej treści zlecenia | Zadania zgodne z nową dekompozycją zachowują stan; zadania bezprzedmiotowe przechodzą w stan „odrzucone”; zadania nowe wchodzą do kolejki jako „oczekujące” |

### 2.9. Koordynacja wielu Wykonawców pracujących równolegle

Koordynator prowadzi pętlę dla dowolnej liczby Wykonawców pracujących równolegle. Każdy Wykonawca ma własny strumień komunikatów sterujących, własny zestaw zadań i własny zakres uprawnień wynikający z Permissions Center jego definicji.

| Mechanizm koordynacji | Ustalenie |
|---|---|
| Przydział względem kompetencji | Zadanie trafia do Wykonawcy, którego definicja — umiejętności, konektory, kanał modelu i uprawnienia — odpowiada charakterowi zadania |
| Zależności między zadaniami | Zadanie zależne rozpoczyna się dopiero po zatwierdzeniu rezultatu zadania poprzedzającego |
| Tory równoległe | Zadania wzajemnie niezależne realizowane są jednocześnie przez odrębnych Wykonawców |
| Wymiana rezultatów pośrednich | Rezultat jednego Wykonawcy wchodzi jako kontekst do zadania drugiego, gdy zależność taka została ustalona w dekompozycji |
| Scalanie rezultatów | Koordynator scala zatwierdzone rezultaty wszystkich torów w jeden wynik zlecenia |
| Rozstrzyganie rozbieżności | Rozbieżne rezultaty równoległych torów kieruje do Wykonawcy pełniącego funkcję kontrolną, a przy braku rozstrzygnięcia — do Użytkownika |
| Podagenci | Wykonawca uruchamia własny Subagent Network (do 15 podagentów) i agreguje ich wyniki przed przekazaniem rezultatu Koordynatorowi (rozdz. 6.2) |

Tryby współpracy dwóch Wykonawców pracujących równolegle w środowisku MultitaskingAI — praca niezależna, przekazywanie wyników, praca naprzemienna i praca iteracyjna — opisano w rozdziale 8.4.

---

## 3. Warstwy widoczności w zarządzaniu agentami

Interfejs Danaco Console ujawnia funkcje stopniowo: funkcja niepotrzebna do realizacji aktualnego zadania nie jest widoczna. Zasada ta obowiązuje warstwę agentową w całości — liczba zapisanych agentów, ról, uprawnień i limitów zasobów nie wpływa na postrzeganą prostotę interfejsu. Każdy element zarządzania agentami należy do dokładnie jednej z czterech warstw widoczności.

| Warstwa | Elementy zarządzania agentami | Sposób dostępu |
|---|---|---|
| 1 — zawsze widoczna | Chat Window, Execution Loop Window w toku realizacji zlecenia, obszar roboczy Agent Builder, pasek kontekstu ze znacznikiem aktywnego Wykonawcy, wskaźniki stanu wykonania zadań | Widoczna bez interakcji |
| 2 — widoczna na żądanie | Wybór Wykonawcy, wybór modelu bazowego, wybór kanału modelu, wybór środowiska, poziom wysiłku | Znacznik kontekstowy i menu progresywne; po użyciu element zwija się samoczynnie |
| 3 — rozwinięcia kontekstowe | Operacje na zapisanej definicji agenta: duplikowanie, eksport, przypisanie do projektu, przypisanie do roli, ustawienia szybkie pamięci | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel nakładkowy, lista rozwijana |
| 4 — funkcje eksperckie | Definiowanie ról, konfiguracja uprawnień w Permissions Center, limity zasobów, konfiguracja ośmiu zakresów izolacji technicznej, diagnostyka pętli wykonawczej | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

### 3.1. Wybór Wykonawcy jako znacznik kontekstowy

Wybór Wykonawcy należy do warstwy 2. Aktywny Wykonawca występuje jako lekki znacznik w pasku kontekstu, obok znaczników środowiska, repozytorium, projektu i modelu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`. Kliknięcie znacznika otwiera selektor Wykonawcy.

Lista zapisanych agentów prezentowana jest jako menu progresywne — jeden element zwinięty `Agent ▼`, którego rozwinięcie zawiera pełny wykaz komponentów własnych rodzaju „agent”. Liczba zapisanych agentów nie zwiększa liczby elementów widocznych w stanie spoczynku interfejsu. Po dokonaniu wyboru menu zwija się samoczynnie, a wybrany Wykonawca pozostaje widoczny wyłącznie jako znacznik kontekstowy.

### 3.2. Role, uprawnienia i limity zasobów jako funkcje warstwy 4

Definiowanie ról, konfiguracja uprawnień i ustalanie limitów zasobów są funkcjami eksperckimi warstwy 4. Nie zajmują miejsca w stanie spoczynku interfejsu i są osiągalne czterema drogami dostępu.

| Droga dostępu | Zastosowanie |
|---|---|
| Polecenie języka naturalnego w Chat Window | Zlecenie zmiany zakresu uprawnień, przypisania roli albo ustalenia limitu zasobów wyrażone zdaniem |
| Skrót klawiszowy | Bezpośrednie otwarcie Permissions Center lub sekcji Role panelu orkiestracji |
| Wyszukiwarka funkcji | Odnalezienie funkcji po nazwie bez znajomości jej umiejscowienia w interfejsie |
| Tryb administracyjny | Praca nad zestawem definicji agentów, ról i limitów zasobów w jednym widoku |

| Funkcja warstwy 4 | Okno realizujące | Rozdział |
|---|---|---|
| Definiowanie ról wykonawczych i przypisanie Wykonawcy do roli | Sekcja Role panelu orkiestracji | 8 |
| Konfiguracja uprawnień: rozszerzenia, moduły i zasoby | Permissions Center | 6.2 |
| Limity zasobów: zakresy izolacji technicznej, dostęp sieciowy, zakres odczytu i zapisu plików, serwer wykonania | Permissions Center | 6.4 |
| Limit podagentów uruchamianych przez Wykonawcę | Permissions Center | 6.2 |
| Profil izolacji przypisany roli | Panel profilu i podglądu okna konfiguracji punktów izolacji | 8.5 |

Każda z tych funkcji jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu nie występuje — ukrycie zmniejsza chaos wizualny i nie utrudnia dostępu.

---
## 4. Okna operacyjne modułu Agents i cykl życia agenta

### 4.1. Okna operacyjne i okna komunikacji

Moduł Agents udostępnia sześć okien operacyjnych wymienionych w rozdziale 11.15 i w Załączniku A Koncepcji platformy oraz oba okna komunikacji operacyjnej opisane w rozdziale 2. Kolumna okien komunikacji sąsiaduje z kolumną obszaru roboczego modułu; regulacji podlega wyłącznie szerokość kolumn.

| Okno | Funkcja | Kolumna | Warstwa widoczności |
|---|---|---|---|
| Chat Window | Kanał Użytkownik ↔ Wykonawca — centralny punkt pracy Użytkownika i podstawowy mechanizm sterowania procesami; w module Agents obsługuje budowę, testowanie i zmianę definicji agenta (rozdz. 2.2). | Lewa, stała, pełna wysokość obszaru roboczego | 1 |
| Execution Loop Window | Kanał Koordynator ↔ Wykonawca — okno pętli wykonawczej: dekompozycja zlecenia, kolejka i stan zadań, komunikaty sterujące, kontrola jakości, sterowanie przebiegiem (rozdz. 2.3). | Kolumna sąsiadująca z Chat Window | 1 |
| Agent Builder | Główne okno tworzenia i edycji definicji agenta; spina pozostałe cztery okna konfiguracyjne w jeden proces. | Obszar roboczy modułu | 1 |
| Model Configuration | Wybór modelu bazowego i kanału jego połączenia (rozdz. 5.1 niniejszego dokumentu). | Obszar roboczy modułu | 2 |
| Skills Manager | Dobór i konfiguracja umiejętności agenta (rozdz. 5.4 niniejszego dokumentu). | Obszar roboczy modułu | 3 |
| Connectors Manager | Podłączanie wtyczek i konektorów zewnętrznych (rozdz. 5.5 niniejszego dokumentu). | Obszar roboczy modułu | 3 |
| Permissions Center | Konfiguracja zakresu uprawnień agenta (rozdz. 6 niniejszego dokumentu). | Kolumna boczna, otwierana jako rozszerzenie boczne | 4 |

Agent Builder jest oknem nadrzędnym wobec czterech okien konfiguracyjnych — każde z nich odpowiada za jeden komponent definicji, a Agent Builder scala ich wynik w pojedynczy, zapisywalny komponent własny.

Makiety tych okien, pełne katalogi ich elementów wraz z warstwą widoczności i sposobem wywołania każdego elementu oraz stany szczegółowe każdego okna zawiera [Moduł Agents](../moduly/agents.md), rozdz. 2 i 5–9.

```
                         AGENT BUILDER
        (okno nadrzędne — scala definicję w komponent własny)
                              │
   ┌──────────────┬───────────┴───────────┬──────────────┐
   ▼              ▼                       ▼              ▼
 Model         Skills                 Connectors      Permissions
 Configuration Manager                Manager         Center
 (model bazowy,(umiejętności)         (wtyczki,       (uprawnienia)
  kanał modelu)                        konektory)

 Chat Window (Użytkownik ↔ Wykonawca) — lewa kolumna obszaru roboczego
 Execution Loop Window (Koordynator ↔ Wykonawca) — kolumna sąsiadująca
```

### 4.2. Cykl życia agenta

Cykl życia agenta przebiega w trzech etapach: utworzenie definicji, zapis jako komponent własny, wykorzystanie operacyjne.

| Etap cyklu życia | Co się dzieje | Okno / miejsce |
|---|---|---|
| 1. Utworzenie definicji | Konfiguracja siedmiu komponentów definicji oraz uprawnień | Agent Builder i cztery okna konfiguracyjne |
| 2. Zapis jako komponent własny | Agent trafia do pamięci aplikacji jako zasób Użytkownika | Agent Builder |
| 3. Wykorzystanie operacyjne | Wybór jako Wykonawca, przypisanie do projektu lub przypisanie do roli | TalkIn, WorkSpace, CodeStudio; Agent Manager; sekcja Role |

Poniższy schemat przedstawia ten cykl w formie tekstowej.

```
STREFA 2 STRONY GŁÓWNEJ (rozdz. 5.2)          OKNO MODUŁOWE AGENTS (środowisko)
              │                                            │
              └───────────────────┬────────────────────────┘
                                  ▼
                            AGENT BUILDER
        ┌───────────────┬─────────┴─────┬───────────────┐
        ▼               ▼               ▼               ▼
  Model Configuration  Skills Manager  Connectors     Permissions
  (model bazowy,       (umiejętności)  Manager        Center
   kanał modelu)                       (wtyczki,      (uprawnienia)
                                        konektory)
        │               │               │               │
        └───────────────┴─────────┬─────┴───────────────┘
                                  ▼
                    + tożsamość, instrukcje systemowe,
                      konfiguracja pamięci  (Agent Builder)
                                  ▼
                     ZAPIS JAKO KOMPONENT WŁASNY
                     (zasób Użytkownika w pamięci aplikacji)
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        ▼                         ▼                         ▼
  Wykonawca                 Agent Manager              Rola w panelu
  w TalkIn / WorkSpace /    modułu Workspace            orkiestracji
  CodeStudio (rozdz. 7)     — przypisanie do            środowiska
                            projektu (rozdz. 9)          MultitaskingAI
                                                          (rozdz. 8)
```

Zapis definicji agenta jako komponentu własnego jest punktem, od którego agent staje się zasobem wybieralnym w dalszej pracy — zgodnie z ogólną zasadą cyklu życia komponentu własnego (rozdz. 2.3 Koncepcji platformy). Trzy gałęzie wykorzystania operacyjnego nie wykluczają się wzajemnie: ten sam agent może jednocześnie być dostępny jako Wykonawca w środowisku TalkIn, przypisany do projektu w module Workspace oraz pełnić rolę Executora 1 w zespole środowiska MultitaskingAI.

---

## 5. Tworzenie agenta — komponenty definicji

Rozdział 11.15 Koncepcji platformy wymienia elementy definiowane przy tworzeniu agenta w Agent Builder: model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory oraz uprawnienia; Załącznik D.1 Koncepcji platformy dodaje pamięć jako pole swoiste rodzaju „agent”. Niniejszy dokument ujmuje ten sam zbiór jako siedem komponentów definicji, opisanych w podrozdziałach 5.1–5.7: wtyczki i konektory — jako dwie postacie tego samego rozszerzenia podłączanego w Connectors Manager — liczone są łącznie jako jeden komponent, a pamięć stanowi komponent odrębny. Uprawnienia, ze względu na zakres, otrzymują odrębny rozdział 6.

### 5.1. Model bazowy i kanał modelu

Model Configuration ustala, jaki model AI stanowi bazę agenta oraz w jaki sposób platforma się z nim łączy. Zgodnie z rozdziałem 9 Architektury wybór modelu i kanału dokonywany jest per sesja lub per rola — dla agenta kanał zapisany w jego definicji jest wartością domyślną, nadpisywalną w momencie przypisania do roli (rozdz. 8.3 niniejszego dokumentu).

| Kanał | Sposób połączenia | Zastosowanie typowe dla agenta |
|---|---|---|
| API | Klucz dostępu, żądanie HTTP. | Modele dostawców zewnętrznych wykorzystywane jako Wykonawca standardowy. |
| CLI | Token narzędzia wiersza poleceń. | Agenci pracujący w module Developer, Terminal lub Diagnostics. |
| SSH | Własne modele pomocnicze na hoście zdalnym. | Agenci wykorzystujący infrastrukturę obliczeniową Użytkownika. |
| HTTP | Modele przeglądarkowe. | Agenci obsługiwani przez moduł Browser lub interfejsy webowe modeli. |

Rdzeń serwera operuje na abstrakcji „wyślij zapytanie, odbierz odpowiedź” niezależnie od wybranego kanału (rozdz. 9 Architektury) — zmiana kanału modelu nie wymaga zmiany pozostałych sześciu elementów definicji.

### 5.2. Tożsamość

Tożsamość jest tym elementem definicji, który odróżnia agenta jako nazwaną jednostkę AI od samego modelu bazowego. Obejmuje ona — zgodnie z ogólnym szablonem komponentu własnego (Załącznik D.1 Koncepcji platformy) — cztery pola.

| Element tożsamości | Opis |
|---|---|
| Nazwa własna | Nazwa odróżniająca agenta jako nazwaną jednostkę AI |
| Opis przeznaczenia | Krótki opis zadania agenta |
| Moduł(y) zastosowania | Środowiska i moduły wykorzystania operacyjnego |
| Zasięg widoczności | Globalny lub projektowy |

Dwóch agentów może korzystać z tego samego modelu i kanału, różniąc się tożsamością, instrukcjami systemowymi, umiejętnościami, konektorami, pamięcią i uprawnieniami — a więc pełnić w praktyce zupełnie odmienne funkcje operacyjne.

### 5.3. Instrukcje systemowe

Instrukcje systemowe są tekstowym opisem sposobu działania agenta, przekazywanym modelowi bazowemu przy każdym wywołaniu. Instrukcje systemowe agenta są komponentem odrębnym od instrukcji systemowych projektu.

| Zestaw instrukcji | Gdzie definiowany | Kiedy obowiązuje |
|---|---|---|
| Instrukcje systemowe agenta | Agent Builder (moduł Agents) | Przy każdym wywołaniu modelu bazowego agenta |
| Instrukcje systemowe projektu | Instructions Panel (moduł Workspace, rozdz. 11.2 Koncepcji platformy) | Gdy agent przypisany do projektu przez Agent Manager |

Gdy agent zostaje przypisany do projektu przez Agent Manager (rozdz. 9 niniejszego dokumentu), oba zestawy instrukcji obowiązują łącznie, zgodnie z zasadą pełnej kompozycyjności (rozdz. 14 Koncepcji platformy, zasada 1).

### 5.4. Umiejętności

Umiejętności dobierane w Skills Manager rozszerzają zakres działania agenta o wyspecjalizowane procedury lub wiedzę proceduralną. Są, obok wtyczek i konektorów, jedną z trzech kategorii rozszerzeń, których wspólny kontrakt integracji opisuje rozdział 14 Architektury.

| Kategoria rozszerzenia | Okno konfiguracji | Co wnosi do agenta | Wspólny kontrakt integracji |
|---|---|---|---|
| Umiejętności | Skills Manager | Wyspecjalizowane procedury lub wiedza proceduralna | rozdz. 14 Architektury |
| Wtyczki | Connectors Manager | Rozszerzenia funkcjonalne podłączane do definicji agenta | rozdz. 14 Architektury |
| Konektory | Connectors Manager | Integracje zewnętrzne podłączane do definicji agenta | rozdz. 14 Architektury |

### 5.5. Wtyczki i konektory

Connectors Manager podłącza do definicji agenta wtyczki i konektory zewnętrzne. Zgodnie z rozdziałem 14 Architektury rozszerzenia dzielą się na dwa źródła.

| Źródło | Charakterystyka |
|---|---|
| Danaco Plugin | Zestaw wbudowany, dostarczany z pakietem aplikacyjnym serwera. |
| Personal | Rozszerzenia instalowane przez Użytkownika z urządzenia; przesyłane na serwer i przechowywane w katalogu użytkownika. |

Rdzeń serwera nie rozróżnia źródła rozszerzenia — obowiązuje wspólny kontrakt integracji (rozdz. 14 Architektury) — więc agent korzysta z wtyczki lub konektora źródła Personal dokładnie tak samo, jak z rozszerzenia źródła Danaco Plugin. Podłączenie rozszerzenia do definicji agenta i decyzja o zakresie jego wykorzystania są rozstrzygane w dwóch różnych oknach.

| Pytanie | Rozstrzyga |
|---|---|
| Czy rozszerzenie jest podłączone do definicji agenta | Connectors Manager (umiejętności — Skills Manager) |
| Czy i w jakim zakresie agent może je wykorzystać (odczyt, zapis, wywołanie akcji) | Permissions Center (rozdz. 6.2 niniejszego dokumentu) |

### 5.6. Pamięć

Konfiguracja pamięci agenta korzysta z ogólnego mechanizmu pamięci opisanego w rozdziale 9 Modelu danych oraz w rozdziale 2.4 i rozdziale 6 Koncepcji platformy: pamięć działa na konfigurowalnych poziomach, podlega pełnej konfiguracji z okna konfiguracji i można ją wyłączyć w całości lub przypisać dowolnie.

| Poziom pamięci | Charakterystyka dla agenta |
|---|---|
| Globalna | Pamięć obejmująca całą platformę |
| Środowisko | Pamięć w zakresie środowiska |
| Projekt | Pamięć w zakresie projektu |
| Sesja | Pamięć w zakresie pojedynczej sesji |
| Wyłączona | Pamięć wyłączona w całości |

Definicja agenta wskazuje, z których poziomów pamięci agent korzysta domyślnie. Ustawienie to jest wartością wyjściową, nadpisywalną w momencie przypisania agenta do projektu (rozdz. 9.2 niniejszego dokumentu) lub do roli (rozdz. 8.5 niniejszego dokumentu), zgodnie z zasadą pierwszeństwa zasięgu najbardziej szczegółowego (rozdz. 6.5 Koncepcji platformy).

### 5.7. Uprawnienia

Uprawnienia agenta konfigurowane są w Permissions Center. Ze względu na zakres tego elementu definicji oraz jego powiązanie z oknem konfiguracji punktów izolacji (rozdz. 6.4–6.6 Koncepcji platformy), Centrum uprawnień opisano odrębnie w rozdziale 6 niniejszego dokumentu.

### 5.8. Zestawienie komponentów definicji

| Komponent | Okno konfiguracji | Podstawa w Koncepcji platformy | Podstawa w Architekturze |
|---|---|---|---|
| Model bazowy i kanał | Model Configuration | rozdz. 11.15 | rozdz. 9 |
| Tożsamość | Agent Builder | rozdz. 11.15, Załącznik D.1 | — |
| Instrukcje systemowe | Agent Builder | rozdz. 11.15 | — |
| Umiejętności | Skills Manager | rozdz. 11.15 | rozdz. 14 |
| Wtyczki i konektory | Connectors Manager | rozdz. 11.15 | rozdz. 14 |
| Pamięć | Agent Builder | rozdz. 11.15, rozdz. 6, rozdz. 2.4, rozdz. 8 | rozdz. 10, rozdz. 12 |
| Uprawnienia | Permissions Center | rozdz. 11.15, rozdz. 6.4–6.6 | rozdz. 12, rozdz. 15 |

---

## 6. Centrum uprawnień (Permissions Center)

### 6.1. Zasada wyjściowa: uwierzytelnienie a uprawnienia agenta

Rozdział 15 Architektury ustala, że uwierzytelnianie jest jedynym mechanizmem kontroli dostępu w systemie, a platforma obsługuje jedno konto właściciela — Użytkownika. Ustalenie to dotyczy dostępu do samej aplikacji. Centrum uprawnień rozstrzyga zagadnienie odrębne — zakres działania agenta, gdy Użytkownik już go uruchomił jako Wykonawcę zadania.

| Zagadnienie | Uwierzytelnienie (rozdz. 15 Architektury) | Centrum uprawnień (Permissions Center) |
|---|---|---|
| Czego dotyczy | Dostęp do samej aplikacji | Zakres działania agenta jako Wykonawcy |
| Podmiot | Użytkownik — jedno konto właściciela | Agent — nie jest podmiotem uwierzytelnianym |
| Rozstrzygane pytanie | Kto może zalogować się i korzystać z platformy | Jak szeroki jest zakres działań agenta w imieniu Użytkownika |

Ta rozdzielność jest bezpośrednim zastosowaniem zasady pełnej konfigurowalności (rozdz. 14 Koncepcji platformy, zasada 2) do warstwy agentowej: zakres działania agenta — jednostki zdolnej do samodzielnego wykonywania zadań, w tym w trybie autonomicznym w środowisku MultitaskingAI — pozostaje przedmiotem jawnej, świadomej konfiguracji Użytkownika, a nie domyślnego ograniczenia wbudowanego w platformę.

### 6.2. Zakres konfigurowany w Centrum uprawnień

Centrum uprawnień ustala cztery grupy zakresu operacyjnego agenta.

| Grupa | Przedmiot konfiguracji | Powiązane okno lub mechanizm |
|---|---|---|
| Dostęp do rozszerzeń | Które z podłączonych w Connectors Manager wtyczek i konektorów agent może rzeczywiście wykorzystać oraz w jakim zakresie (odczyt, zapis, wywołanie akcji). | Connectors Manager (rozdz. 5.5), rozdz. 14 Architektury |
| Dostęp do modułów i zasobów | W których modułach platformy agent może działać jako Wykonawca oraz do jakich zasobów modułowych — zapisu plików w module Library i wykonywania poleceń w module Terminal — ma dostęp. | Macierz dostępności modułów (rozdz. 10 Koncepcji platformy) jako ramy działania agenta |
| Zakresy izolacji technicznej | Domyślne ustawienie ośmiu zakresów izolacji technicznej właściwe temu agentowi. | Okno konfiguracji punktów izolacji (rozdz. 6.4 Koncepcji platformy), rozdz. 12 Architektury |
| Zakres działania w MultitaskingAI | Możliwość uruchamiania Subagent Network (do 15 podagentów) oraz udział w kolejkach i orkiestracji, gdy agent pełni rolę wykonawczą. | Panel orkiestracji (rozdz. 13.8 Koncepcji platformy), rozdz. 8 niniejszego dokumentu |

### 6.3. Relacja do okna konfiguracji punktów izolacji

Centrum uprawnień i okno konfiguracji punktów izolacji (rozdz. 6.4–6.6 Koncepcji platformy) są dwoma punktami wejścia do tej samej macierzy izolacji technicznej opisanej w rozdziale 12 Architektury. Ustawienia zapisane w Centrum uprawnień pełnią funkcję wartości wyjściowej właściwej danemu agentowi — obowiązują wszędzie, gdzie agent działa, dopóki nie zostaną nadpisane przez regułę zapisaną na bardziej szczegółowym poziomie zasięgu. Rozdział 6.5 Koncepcji platformy przewiduje siedem poziomów zasięgu, uporządkowanych według pierwszeństwa.

| Poziom zasięgu | Pierwszeństwo | Przykład zastosowania dla agenta |
|---|---|---|
| Globalny | Najniższe — warstwa bazowa | Domyślna polityka izolacji obowiązująca w całej platformie |
| Środowisko | ↑ | Odmienna polityka dla środowiska CodeStudio niż dla TalkIn |
| Moduł | ↑ | Odrębna pamięć przypisana modułowi Developer |
| Para modułów | ↑ | Współdzielenie mechanizmu między dwoma modułami |
| Projekt | ↑ | Rozdzielenie kontekstu między dwoma projektami |
| Karta sesji | ↑ | Jednorazowe współdzielenie historii między kartami |
| Rola (MultitaskingAI) | Najwyższe | Profil izolacji przypisany roli, którą pełni agent |

Zmiana dokonana na poziomie bardziej szczegółowym nakłada się na ustawienie z Centrum uprawnień, nie modyfikując zapisanej w nim definicji agenta — zgodnie z zasadą nienaruszania warstwy domyślnej. Poniższy schemat przedstawia tę relację w formie tekstowej.

```
DEFINICJA AGENTA                          OKNO KONFIGURACJI PUNKTÓW IZOLACJI
Permissions Center                        (rozdz. 6.4–6.6 Koncepcji platformy)
(wartość wyjściowa
 właściwa agentowi)
        │                                            │
        │            zasięg: Globalny ───────────────┤  (najniższe pierwszeństwo)
        │            zasięg: Środowisko ─────────────┤
        │            zasięg: Moduł ──────────────────┤
        │            zasięg: Para modułów ───────────┤
        │            zasięg: Projekt ────────────────┤
        │            zasięg: Karta sesji ────────────┤
        └──────────► zasięg: Rola (MultitaskingAI) ──┤  (najwyższe pierwszeństwo)
                                                     ▼
                                              POLITYKA EFEKTYWNA
                                        (podgląd — panel profilu i podglądu)
```

### 6.4. Stan wyjściowy i macierz zakresów izolacji technicznej

Zgodnie z rozdziałem 6.6 Koncepcji platformy żaden z ośmiu zakresów izolacji technicznej nie jest domyślnie aktywny — stan ten dotyczy również agenta jako komponentu własnego. Stanem wyjściowym Centrum uprawnień jest pełny dostęp operacyjny: agent może korzystać ze wszystkich podłączonych w jego definicji rozszerzeń, działać we wszystkich modułach, w których został wykorzystany jako Wykonawca, i nie podlega żadnemu z ośmiu zakresów izolacji technicznej.

| Zakres izolacji technicznej (rozdz. 12 Architektury) | Stan wyjściowy | Przykład świadomego zawężenia przez Użytkownika |
|---|---|---|
| Katalog roboczy sesji | Wyłączony (pełny dostęp) | Ograniczenie do wskazanego katalogu dla pracy autonomicznej |
| Środowisko procesu | Wyłączony (pełny dostęp) | Odrębne środowisko procesu agenta |
| Katalog danych i konfiguracji modelu | Wyłączony (pełny dostęp) | Odrębny katalog danych i konfiguracji modelu |
| Dostęp sieciowy | Wyłączony (pełny dostęp) | Ograniczenie dostępu sieciowego agenta pracującego bez nadzoru |
| Zakres odczytu i zapisu plików | Wyłączony (pełny dostęp) | Wyłącznie odczyt dla agenta w roli kontrolnej |
| Konto i token per sesja | Wyłączony (pełny dostęp) | Odrębne konto i token dla sesji agenta |
| Model procesu | Wyłączony (pełny dostęp) | Odrębny model procesu |
| Serwer wykonania | Wyłączony (pełny dostęp) | Wskazany serwer wykonania |

Stan wyjściowy odpowiada nadrzędnej zasadzie „brak ustawienia = wartość domyślna” (rozdz. 6.5 Koncepcji platformy) oraz zasadzie braku twardych blokad wbudowanych na stałe (rozdz. 6 Koncepcji platformy). Centrum uprawnień nie wymusza ograniczenia, lecz udostępnia je jako ustawienie konfiguracyjne — na przykład dla agenta wykonującego zadania w trybie autonomicznym 24 godziny na dobę, 7 dni w tygodniu w środowisku MultitaskingAI (rozdz. 13.6 Koncepcji platformy), gdzie ograniczenie dostępu sieciowego lub zakresu zapisu plików jest decyzją operacyjną Użytkownika, a nie wymaganiem platformy.

---

## 7. Agent jako komponent własny ponad środowiskami

Po zapisaniu definicji agent staje się Wykonawcą wybieralnym w środowiskach TalkIn, WorkSpace i CodeStudio (rozdz. 11.15 Koncepcji platformy), a po przypisaniu do roli — Wykonawcą w środowisku MultitaskingAI. Sposób wykorzystania różni się jakościowo między tymi dwoma trybami.

| Środowisko | Sposób wykorzystania agenta | Charakter | Rozdział |
|---|---|---|---|
| TalkIn | Wybierany jako Wykonawca w kontekście modułu | Doraźny, per zadanie | 7.1 |
| WorkSpace | Wybierany jako Wykonawca; przypisywany do projektu | Doraźny oraz przypisanie do projektu | 7.1, 9 |
| CodeStudio | Wybierany jako Wykonawca w kontekście modułu | Doraźny, per zadanie | 7.1 |
| MultitaskingAI | Przypisywany do trwałej roli w zespole modeli | Trwały, per rola | 7.2, 8 |

### 7.1. Wykorzystanie w TalkIn, WorkSpace i CodeStudio

Sposób wykorzystania jest spójny z ogólnym mechanizmem komponentu własnego (rozdz. 2.3 Koncepcji platformy): agent nie otwiera własnego, odrębnego okna operacyjnego w tych trzech środowiskach — jest wybierany jako Wykonawca w kontekście modułu, w którym Użytkownik aktualnie pracuje, analogicznie do sposobu, w jaki automatykę wpina się w sesji modułu Developer.

### 7.2. Wykorzystanie w środowisku MultitaskingAI

Wykorzystanie agenta w środowisku MultitaskingAI różni się jakościowo od pozostałych trzech środowisk: agent nie jest tam wybierany doraźnie jako Wykonawca pojedynczego zadania, lecz przypisywany do trwałej roli w zespole modeli (rozdz. 13.3 Koncepcji platformy). Ten tryb wykorzystania opisano w rozdziale 8 niniejszego dokumentu.

---

## 8. Przypisanie agenta do ról w środowisku MultitaskingAI

### 8.1. Cztery role i Subagent Network

Środowisko MultitaskingAI opiera się na czterech rolach oraz na mechanizmie Subagent Network uruchamianym przez Wykonawcę (rozdz. 13.3 i Załącznik B Koncepcji platformy). Sekcja Role panelu orkiestracji (rozdz. 13.8 Koncepcji platformy) udostępnia cztery okna robocze odpowiadające tym rolom oraz mechanizm przypisania agentów z modułu Agents do każdej z nich.

| Rola | Zakres odpowiedzialności | Typowe przypisanie agenta |
|---|---|---|
| Executor 1 | Główny Wykonawca — faktyczna realizacja pracy. | Agent wyspecjalizowany w typie zadania dominującym w danym procesie (np. programowanie, redakcja treści). |
| Executor 2 | Równoległy Wykonawca — drugi, niezależny tor wykonania. | Agent wyspecjalizowany w komplementarnym typie zadania (np. frontend przy Executorze 1 realizującym backend). |
| Coordinator | Planowanie, podział pracy, budowa promptów, sterowanie procesem, zarządzanie kolejką. | Agent skonfigurowany pod kątem organizacji procesu, nie wytwarzania produktu końcowego. |
| Executor 3 / Validator | Funkcja kontrolna lub doradcza wobec pracy pozostałych modeli (Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator). | Agent wyspecjalizowany we wcieleniu właściwym danemu procesowi. |

Rolę pełni agent zapisany jako komponent własny albo model bazowy wskazany bezpośrednio, bez pośrednictwa agenta (Załącznik D.2 Koncepcji platformy). Wybór między tymi dwiema możliwościami jest ustawieniem konfiguracyjnym sekcji Role panelu orkiestracji.

### 8.2. Mechanika przypisania

Przypisanie odbywa się w sekcji Role panelu orkiestracji według poniższego przepływu.

```
1. Użytkownik otwiera sekcję Role panelu orkiestracji.
2. Dla danej roli wskazuje jedną z dwóch możliwości:
      • agenta zapisanego jako komponent własny w module Agents, albo
      • model bazowy bez pośrednictwa agenta.
3. Przypisanie agenta niesie ze sobą całą jego definicję:
      tożsamość · instrukcje systemowe · umiejętności · konektory · pamięć · uprawnienia
      (z Centrum uprawnień)
4. Definicja stanowi punkt wyjścia dla działania w tej roli, nadpisywalny przez
      ustawienia zapisane na poziomie zasięgu „Rola” w oknie konfiguracji
      punktów izolacji (rozdz. 6.3 niniejszego dokumentu).
```

### 8.3. Kanał modelu per rola

Zgodnie z rozdziałem 9 Architektury wybór modelu i kanału dokonywany jest per sesja lub per rola. Kanał modelu zapisany w definicji agenta (rozdz. 5.1 niniejszego dokumentu) może zostać dla danej roli nadpisany.

| Kontekst | Kanał domyślny (definicja agenta) | Kanał nadpisany (rola) |
|---|---|---|
| Agent poza rolą, Wykonawca standardowy | API | — |
| Ten sam agent w roli Executor 1 wymagającej wiersza polecenia w module Developer lub Terminal | API | CLI |

### 8.4. Relacje między Wykonawcami

Relacje między Executorem 1 a Executorem 2 są w pełni konfigurowalne (rozdz. 13.3 Koncepcji platformy). Tryb współpracy pozostaje w gestii Użytkownika lub Coordinatora i jest ustawiany niezależnie od tego, czy dana rola jest pełniona przez agenta, czy przez model bazowy.

| Tryb współpracy Executor 1 — Executor 2 | Opis |
|---|---|
| Praca niezależna | Dwa niezależne tory wykonania działające równolegle |
| Przekazywanie wyników | Wymiana wyników pośrednich między Wykonawcami |
| Praca naprzemienna | Naprzemienne przekazywanie zadań |
| Praca iteracyjna | Wzajemne poprawianie rezultatów w kolejnych przebiegach |

### 8.5. Profil izolacji roli

Rola, do której przypisano agenta, może otrzymać własny profil izolacji zapisany i przypisany z poziomu panelu profilu i podglądu w oknie konfiguracji punktów izolacji (rozdz. 6.6 Koncepcji platformy). Ponieważ poziom zasięgu „Rola” ma najwyższe pierwszeństwo spośród siedmiu poziomów (rozdz. 6.5 Koncepcji platformy), profil przypisany roli nadpisuje zarówno ustawienia zapisane w Centrum uprawnień agenta, jak i wszelkie reguły z poziomów szerszych (globalny, środowisko, moduł, para modułów, projekt, karta sesji).

---

## 9. Przypisanie agenta do projektów w module Workspace

### 9.1. Agent Manager

Moduł Workspace udostępnia okno Agent Manager (rozdz. 11.2 i Załącznik A Koncepcji platformy), przez które agenci skonfigurowani w module Agents są przypisywani do konkretnego projektu jako Wykonawcy w jego ramach (rozdz. 11.15 Koncepcji platformy). Przypisanie jest jawne i konfigurowalne — Użytkownik wskazuje w Agent Manager, którzy z zapisanych agentów mają zastosowanie w danym projekcie.

### 9.2. Zakres domyślny i rozszerzenie

Zgodnie z rozdziałem 11.2 Koncepcji platformy agenci przypisani do projektu przez Agent Manager działają domyślnie w zakresie tego projektu. Poniższa tabela porządkuje warstwy konfiguracji obowiązujące agenta pracującego w projekcie.

| Warstwa konfiguracji | Źródło | Zakres obowiązywania |
|---|---|---|
| Instrukcje systemowe projektu | Instructions Panel (moduł Workspace) | Łącznie z instrukcjami systemowymi agenta |
| Pamięć kontekstowa projektu | Context Memory (moduł Workspace) | Łącznie z pamięcią zapisaną w definicji agenta |
| Instrukcje systemowe i pamięć agenta | Definicja agenta (rozdz. 5.3 i 5.6) | Zawsze |
| Rozdzielenie od innego projektu | Okno konfiguracji punktów izolacji, poziom „Projekt” lub „Para modułów” | Domyślnie pełne rozdzielenie kontekstu (rozdz. 6.1 Koncepcji platformy) |

To, czy i w jakim stopniu konfiguracja może wpływać na inny, równolegle prowadzony projekt, Użytkownik określa samodzielnie w oknie konfiguracji punktów izolacji — domyślnym stanem pozostaje pełne rozdzielenie kontekstu między projektami.

### 9.3. Powiązanie z modelem danych

Rozdział 7.1 Modelu danych opisuje encję Projekt jako jednostkę organizacji pracy niosącą przypisany profil izolacji, instrukcje systemowe i zakres pamięci — pełny wykaz pól encji zawiera [Model danych](../architektura/model-danych.md) rozdz. 7.1. Rozdział 13.2 Modelu danych potwierdza, że agent jest niezależny od środowiska i przypisywany jako Wykonawca. Koncepcja platformy przewiduje dwa mechanizmy tego przypisania.

| Mechanizm przypisania agenta | Okno | Rozdział niniejszego dokumentu |
|---|---|---|
| Przypisanie do projektu | Agent Manager (moduł Workspace) | 9 |
| Przypisanie do roli | Sekcja Role panelu orkiestracji (środowisko MultitaskingAI) | 8 |

---

## 10. Szablon definicji agenta

Poniższy szablon rozwija dla rodzaju „agent” ogólny szablon komponentu własnego z Załącznika D.1 Koncepcji platformy, uwzględniając pola swoiste wskazane w tym załączniku oraz komponenty wyspecyfikowane w rozdziałach 5 i 6 niniejszego dokumentu. Wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna” (rozdz. 6 i rozdz. 14 Koncepcji platformy).

| Pole | Opis | Okno konfiguracji | Wartości / przykład |
|---|---|---|---|
| Nazwa | Nazwa własna agenta | Agent Builder | dowolny tekst (np. „Agent Backend”) |
| Przeznaczenie | Krótki opis zadania agenta | Agent Builder | dowolny tekst |
| Model bazowy | Model AI stanowiący podstawę agenta | Model Configuration | nazwa modelu |
| Kanał modelu | Sposób połączenia z modelem bazowym | Model Configuration | API \| CLI \| SSH \| HTTP (rozdz. 9 Architektury) |
| Instrukcje systemowe | Tekstowy opis sposobu działania agenta | Agent Builder | dowolny tekst |
| Umiejętności | Wyspecjalizowane procedury lub wiedza proceduralna | Skills Manager | lista umiejętności |
| Wtyczki i konektory | Podłączone rozszerzenia zewnętrzne | Connectors Manager | lista rozszerzeń; źródło: Danaco Plugin \| Personal (rozdz. 14 Architektury) |
| Pamięć | Poziomy pamięci, z których agent korzysta domyślnie | Agent Builder | globalna \| środowisko \| projekt \| sesja; wyłączona (rozdz. 9 Modelu danych) |
| Uprawnienia — rozszerzenia | Zakres wykorzystania podłączonych konektorów | Permissions Center | pełny dostęp (domyślnie) \| zawężony per konektor |
| Uprawnienia — moduły i zasoby | Moduły i zasoby modułowe dostępne agentowi | Permissions Center | pełny dostęp (domyślnie) \| zawężony per moduł |
| Uprawnienia — izolacja techniczna | Domyślne ustawienie ośmiu zakresów izolacji technicznej | Permissions Center | każdy zakres: włączony \| wyłączony (domyślnie wyłączony, rozdz. 12 Architektury) |
| Uprawnienia — MultitaskingAI | Możliwość uruchamiania Subagent Network i udziału w orkiestracji | Permissions Center | tak (do 15 podagentów) \| nie |
| Moduł(y) zastosowania | Środowiska i moduły, w których agent jest wykorzystywany operacyjnie | Agent Builder | TalkIn, WorkSpace, CodeStudio; rola w MultitaskingAI |
| Powiązania | Przypisania do projektów i ról | Agent Manager (rozdz. 9.2), sekcja Role panelu orkiestracji (rozdz. 8) | lista przypisań |
| Widoczność / zasięg | Zakres dostępności agenta | Agent Builder | globalny \| projektowy (Załącznik D.1 Koncepcji platformy) |

Ten sam szablon w formie zwartej struktury konfiguracyjnej, wygodnej przy implementacji:

```
AGENT
  nazwa:                 <nazwa własna agenta>
  przeznaczenie:         <krótki opis zadania agenta>

  model_bazowy:          <nazwa modelu>
  kanał_modelu:          API | CLI | SSH | HTTP        # domyślny; nadpisywalny per rola

  instrukcje_systemowe:  <tekst>

  umiejętności:          [ <umiejętność>, … ]          # Skills Manager
  rozszerzenia:                                        # Connectors Manager
    - nazwa:             <wtyczka lub konektor>
      źródło:            Danaco Plugin | Personal

  pamięć:
    poziomy:             [ globalna | środowisko | projekt | sesja ]   # lub: wyłączona

  uprawnienia:                                         # Permissions Center
    rozszerzenia:        pełny dostęp | zawężony per konektor
    moduły_i_zasoby:     pełny dostęp | zawężony per moduł
    izolacja_techniczna:                               # osiem zakresów; domyślnie wyłączone
      katalog_roboczy_sesji:              włączony | wyłączony
      środowisko_procesu:                 włączony | wyłączony
      katalog_danych_i_konfiguracji:      włączony | wyłączony
      dostęp_sieciowy:                    włączony | wyłączony
      odczyt_i_zapis_plików:              włączony | wyłączony
      konto_i_token_per_sesja:            włączony | wyłączony
      model_procesu:                      włączony | wyłączony
      serwer_wykonania:                   włączony | wyłączony
    multitaskingai:
      subagent_network:  tak (do 15 podagentów) | nie

  moduły_zastosowania:   [ TalkIn, WorkSpace, CodeStudio, rola w MultitaskingAI ]
  powiązania:            [ <projekt / rola> ]          # Agent Manager, sekcja Role
  widoczność:            globalny | projektowy
```

---

## 11. Model danych i kontrakty komunikacji agenta

### 11.1. Encja Agent

Rozdział 12.1 Modelu danych definiuje encję Agent jako definicję agenta AI, komponent platformowy ponad wszystkimi środowiskami. Zestaw atrybutów odpowiada siedmiu komponentom definicji z rozdziału 5 oraz Centrum uprawnień z rozdziału 6.

| Atrybut encji Agent | Odpowiednik w definicji agenta |
|---|---|
| identyfikator | — (klucz encji) |
| nazwa | Tożsamość — nazwa własna (rozdz. 5.2) |
| model bazowy | Model Configuration (rozdz. 5.1) |
| tożsamość | Tożsamość (rozdz. 5.2) |
| instrukcje systemowe | Instrukcje systemowe (rozdz. 5.3) |
| umiejętności | Skills Manager (rozdz. 5.4) |
| wtyczki | Connectors Manager (rozdz. 5.5) |
| konektory | Connectors Manager (rozdz. 5.5) |
| konfiguracja pamięci | Pamięć (rozdz. 5.6) |
| konfiguracja uprawnień | Permissions Center (rozdz. 6) |

### 11.2. Encja Rola w MultitaskingAI

Rozdział 13.2 Modelu danych definiuje encję Rola w MultitaskingAI jako przypisanie roli w środowisku wielomodelowym.

| Atrybut encji Rola w MultitaskingAI | Wartość |
|---|---|
| identyfikator | — (klucz encji) |
| sesja lub orkiestracja | Kontekst, w którym rola jest przypisana |
| rodzaj roli | Executor 1 \| Executor 2 \| Coordinator \| Executor 3 / Validator |
| przypisany model | Model bazowy albo agent (pośrednik niosący pełną definicję z rozdz. 11.1) |
| konfiguracja współpracy | Tryb współpracy między Wykonawcami (rozdz. 8.4) |

Atrybut „przypisany model” obejmuje zarówno przypisanie bezpośrednie modelu bazowego, jak i przypisanie agenta — pełniącego w takim wypadku funkcję pośrednika niosącego pełną definicję, zgodnie z mechanizmem z rozdziału 8.2 niniejszego dokumentu.

### 11.3. Komunikaty warstwy komunikacji

Rozdział 10.15 Kontraktów komunikacji definiuje komunikaty właściwe modelom i agentom, wymieniane kanałem WebSocket (rozdz. 11 Architektury).

| Komunikat | Funkcja |
|---|---|
| `model.channel.set` | Konfiguracja kanału modelu (rozdz. 5.1 niniejszego dokumentu). |
| `agent.create` | Utworzenie agenta. |
| `agent.update` | Zmiana agenta, w tym zmiana definicji dokonana w Agent Builder. |
| `agent.list` | Wykaz agentów zapisanych jako komponenty własne. |
| `role.assign` | Przypisanie roli w środowisku wielomodelowym (rozdz. 8.2 niniejszego dokumentu). |

Konfiguracja uprawnień agenta korzysta dodatkowo z komunikatów i zdarzenia rozgłaszanego przez serwer.

| Komunikat / zdarzenie | Funkcja | Rozdział Kontraktów komunikacji |
|---|---|---|
| `agent.changed` | Zdarzenie rozgłaszane po każdej zmianie definicji agenta; niesie stan wystarczający do odświeżenia widoku klienta bez dodatkowego zapytania. | rozdz. 10.15 |
| `extension.list` | Wykaz rozszerzeń podłączonych w Connectors Manager. | rozdz. 17 |
| `extension.install` | Instalacja rozszerzenia. | rozdz. 17 |
| `extension.toggle` | Włączenie lub wyłączenie rozszerzenia. | rozdz. 17 |
| `isolation.profile.save` | Zapis profilu izolacji. | rozdz. 15 |
| `isolation.profile.assign` | Przypisanie profilu izolacji roli (rozdz. 8.5 niniejszego dokumentu). | rozdz. 15.3 |

### 11.4. Wykaz pełny komend obszarów `agent`, `subagent`, `advisor`

Rozdział 18 Architektury i `budowa/shared/contract.json` wyznaczają dla warstwy agentowej trzy obszary kontraktu: `agent` (30 komend), `subagent` (4 komendy) i `advisor` (1 komenda). Poniższe wykazy niosą pełną treść każdej komendy — pola żądania i pola wyniku wraz z wymagalnością — zgodnie z regułą weryfikowalności (rozdz. 4.1 Standardu redakcyjnego i językowego). Każda komenda oznaczona `(wym)` jest polem wymaganym, `(opc)` — polem opcjonalnym.

**Obszar `agent` — 30 komend.**

| Komenda | Pola żądania | Pola wyniku |
|---|---|---|
| `agent.create` | `name:string(wym)`, `description:string(opc)`, `systemPrompt:string(opc)`, `channelId:string(opc)`, `model:string(opc)`, `displayName:string(opc)`, `favicon:string(opc)`, `mode:IdentityMode(opc)`, `memoryLevels:MemoryLevel[](opc)`, `visibility:AgentVisibility(opc)` | `agent:Agent(wym)` |
| `agent.update` | `agentId:string(wym)`, `name:string(opc)`, `description:string(opc)`, `systemPrompt:string(opc)`, `enabled:bool(opc)`, `displayName:string(opc)`, `favicon:string(opc)`, `mode:IdentityMode(opc)`, `memoryLevels:MemoryLevel[](opc)`, `visibility:AgentVisibility(opc)` | `agent:Agent(wym)` |
| `agent.list` | `query:string(opc)`, `enabledOnly:bool(opc)`, `limit:int(opc)`, `projectId:string(opc)` | `agents:Agent[](wym)`, `total:int(opc)` |
| `agent.delete` | `agentId:string(wym)` | `agentId:string(wym)`, `deleted:bool(wym)` |
| `agent.model.set` | `agentId:string(wym)`, `channelId:string(wym)`, `model:string(opc)`, `transport:ProviderTransport(opc)`, `parameters:json(opc)`, `effort:ReasoningEffort(opc)` | `agent:Agent(wym)` |
| `agent.skill.add` | `agentId:string(wym)`, `skillId:string(wym)` | `agent:Agent(wym)` |
| `agent.skill.remove` | `agentId:string(wym)`, `skillId:string(wym)` | `agent:Agent(wym)` |
| `agent.connector.add` | `agentId:string(wym)`, `name:string(wym)`, `kind:AgentConnectorKind(wym)`, `accessPointId:string(opc)`, `config:json(opc)` | `connector:AgentConnector(wym)` |
| `agent.connector.list` | `agentId:string(wym)` | `connectors:AgentConnector[](wym)` |
| `agent.connector.remove` | `agentId:string(wym)`, `connectorId:string(wym)` | `agent:Agent(wym)` |
| `agent.connector.configure` | `agentId:string(wym)`, `connectorId:string(wym)`, `accessPointId:string(opc)`, `config:json(opc)`, `enabled:bool(opc)` | `connector:AgentConnector(wym)` |
| `agent.plugin.add` | `agentId:string(wym)`, `name:string(wym)`, `source:string(opc)`, `version:string(opc)` | `plugin:AgentPlugin(wym)` |
| `agent.plugin.remove` | `agentId:string(wym)`, `pluginId:string(wym)` | `agent:Agent(wym)` |
| `agent.plugin.list` | `agentId:string(wym)` | `plugins:AgentPlugin[](wym)` |
| `agent.permission.set` | `agentId:string(wym)`, `group:AgentPermissionGroup(wym)`, `granted:bool(wym)`, `scope:string(opc)` | `permissions:AgentPermission[](wym)` |
| `agent.permission.remove` | `agentId:string(wym)`, `group:AgentPermissionGroup(opc)`, `scope:string(opc)` | `permissions:AgentPermission[](wym)`, `removed:int(wym)` |
| `agent.layer.set` | `agentId:string(wym)`, `layer:IdentityLayer(wym)`, `content:string(wym)`, `enabled:bool(opc)` | `agent:Agent(wym)` |
| `agent.layer.remove` | `agentId:string(wym)`, `layer:IdentityLayer(wym)` | `agent:Agent(wym)` |
| `agent.version.list` | `agentId:string(wym)`, `limit:int(opc)` | `agentId:string(wym)`, `versions:AgentVersion[](wym)`, `total:int(wym)` |
| `agent.version.get` | `agentId:string(wym)`, `versionId:string(wym)` | `version:AgentVersion(wym)`, `snapshot:AgentVersionSnapshot(wym)` |
| `agent.version.restore` | `agentId:string(wym)`, `versionId:string(wym)` | `agent:Agent(wym)`, `version:AgentVersion(wym)` |
| `agent.archive` | `agentId:string(wym)` | `agentId:string(wym)`, `archived:bool(wym)`, `agent:Agent(opc)` |
| `agent.restore` | `agentId:string(wym)` | `agentId:string(wym)`, `restored:bool(wym)`, `agent:Agent(opc)` |
| `agent.archive.list` | `offset:int(opc)`, `limit:int(opc)` | `agents:Agent[](wym)`, `total:int(wym)` |
| `agent.assignment.list` | `agentId:string(opc)`, `kind:AgentAssignmentKind(opc)` | `assignments:AgentAssignment[](wym)`, `total:int(wym)` |
| `agent.modules.set` | `agentId:string(wym)`, `moduleCodes:string[](wym)` | `agent:Agent(wym)` |
| `agent.isolation.get` | `agentId:string(wym)` | `agentId:string(wym)`, `switches:IsolationTechnicalSwitch[](wym)` |
| `agent.isolation.set` | `agentId:string(wym)`, `switches:IsolationTechnicalSwitch[](wym)` | `switches:IsolationTechnicalSwitch[](wym)` |
| `agent.policy.get` | `agentId:string(wym)`, `windowId:string(opc)` | `policy:AgentPolicy(wym)` |
| `agent.subagent.set` | `agentId:string(wym)`, `enabled:bool(wym)`, `limit:int(opc)` | `agent:Agent(wym)` |

Zdarzenie zwrotne obszaru: `agent.changed` — rozgłaszane po każdej zmianie agenta zapisanego jako komponent własny (rozdz. 11.3).

Pola `agent.modules.set.moduleCodes` i `agent.update.memoryLevels` niosą znaczenie własne dla listy pustej `[]`: pusty zbiór oznacza brak ograniczenia, czyli pełną dostępność, i jest stanem wyjściowym; pominięcie pola oznacza brak zmiany wobec stanu zapisanego. Komenda `agent.permission.remove` przywraca stan „brak ustawienia = wartość domyślna” (rozdz. 6.4 niniejszego dokumentu) — sam wpis zakresu pozostaje w wykazie uprawnień, wartość wraca do przyznanej.

**Obszar `subagent` — 4 komendy.** Mechanizm Subagent Network uruchamiany przez Wykonawcę (rozdz. 2.9 i rozdz. 8.1 niniejszego dokumentu); górna granica podagentów wynosi piętnaście.

| Komenda | Pola żądania | Pola wyniku |
|---|---|---|
| `subagent.spawn` | `windowId:string(wym)`, `task:string(wym)`, `name:string(opc)`, `count:int(opc)`, `modelChannelId:string(opc)` | `subagents:Subagent[](wym)` |
| `subagent.list` | `windowId:string(opc)`, `sessionId:string(opc)`, `status:SubagentStatus(opc)` | `subagents:Subagent[](wym)` |
| `subagent.result.collect` | `windowId:string(opc)`, `subagentIds:string[](opc)`, `waitForAll:bool(opc)` | `subagents:Subagent[](wym)`, `complete:bool(wym)` |
| `subagent.stop` | `windowId:string(opc)`, `subagentIds:string[](opc)` | `stopped:string[](wym)`, `notRunning:string[](wym)`, `subagents:Subagent[](wym)` |

Zdarzenie zwrotne obszaru: `subagent.changed` — zmiana stanu podagenta okna wykonawcy (rozdz. 11.2 Kontraktów komunikacji).

**Obszar `advisor` — 1 komenda.** Konsultacja u doradcy: model pyta inny model o radę w trakcie tury. Rada nie jest wiążąca i nie wykonuje pracy za pytającego — odpowiedzialność za wynik zostaje przy agencie, który pyta.

| Komenda | Pola żądania | Pola wyniku |
|---|---|---|
| `advisor.consult` | `windowId:string(wym)`, `question:string(wym)`, `context:string(opc)`, `requestedAdvisor:string(opc)` | `advice:string(wym)`, `advisorChannel:string(wym)`, `advisorModel:string(opc)`, `selection:AdvisorSelection(wym)` |

Zdarzenie zwrotne obszaru: `advisor.consulted` — odbyta konsultacja; niesie okno, doradcę, podstawę doboru i skrót rady. Pole `requestedAdvisor` jest prośbą podlegającą sufitowi siły: prośba o model silniejszy od kanału pytającego bez wcześniejszego wskazania Operatora jest odmawiana, nie zamieniana po cichu na słabszego doradcę.

Wszystkie trzydzieści pięć komend trzech obszarów korzysta ze wspólnego wykazu dziewięciu kodów błędów kontraktu ([Kontrakty komunikacji](../architektura/kontrakty-komunikacji.md) rozdz. 19, Załącznik D): `validation_failed`, `not_found`, `not_authenticated`, `permission_denied`, `conflict`, `command_not_understood` (żaden z nich nie jest ponawialny) oraz `channel_unavailable`, `rate_limited`, `internal_error` (ponawialne — dopuszczają automatyczne ponowienie zadania przez klienta). Pole `kodyBledow` w `budowa/shared/contract.json` zna osiem pierwszych pozycji; dopisanie dziewiątego kodu `command_not_understood` jest osobnym zadaniem w kodzie.

---

## 12. Konfigurowalność i zasada braku twardych blokad

Cztery zasady nadrzędne funkcjonalności platformy (rozdz. 14 Koncepcji platformy) obowiązują warstwę agentową na równi z pozostałymi elementami platformy.

| Zasada | Zastosowanie do warstwy agentowej |
|---|---|
| Pełna kompozycyjność | Agent może zostać złożony z dowolną kombinacją umiejętności, konektorów i uprawnień; może uczestniczyć w pełnej pętli Coordinator–Executor z Subagent Network (rozdz. 13.3 Koncepcji platformy). |
| Pełna konfigurowalność (zasada centralna) | Każdy element definicji agenta — model, tożsamość, instrukcje, umiejętności, konektory, pamięć, uprawnienia — konfigurowany jest z poziomu okien modułu Agents i okna konfiguracji punktów izolacji; żaden element nie jest zaszyty na stałe. |
| Jawność i konfigurowalność zależności | Przypisanie agenta do projektu (Agent Manager) i do roli (panel orkiestracji) jest jawną decyzją Użytkownika, a nie zależnością wbudowaną w platformę. |
| Rozszerzenie orkiestracji | Warstwa agentowa — w szczególności Subagent Network i wieloagentowa pętla środowiska MultitaskingAI — prowadzi pracę wielu modeli i Wykonawców równolegle w ramach jednej pętli wykonawczej. |

Zasada braku twardych blokad, wyrażona w rozdziale 6 Koncepcji platformy („Platforma nie narzuca twardej izolacji historii sesyjnej ani pamięci. Nie ma w niej blokad wbudowanych na stałe.”), obowiązuje w Centrum uprawnień w postaci opisanej w rozdziale 6.4 niniejszego dokumentu: stanem wyjściowym jest pełny dostęp operacyjny agenta, a każde zawężenie jest świadomą decyzją Użytkownika podjętą z poziomu okna konfiguracji.

---

## 13. Scenariusze użycia

Scenariusze ilustrują kompozycję mechanizmów opisanych w niniejszym dokumencie — pokazują współdziałanie modułu Agents z modułem Workspace i środowiskiem MultitaskingAI, rozwijając scenariusz E.1 Załącznika E Koncepcji platformy w odniesieniu do warstwy agentowej.

### 13.1. Budowa agenta specjalistycznego i przypisanie do projektu

Użytkownik tworzy w Agent Builder agenta „Agent Redaktor” i przypisuje go do dwóch równolegle prowadzonych projektów redakcyjnych.

```
1. Agent Builder          → utworzenie agenta „Agent Redaktor”
2. Model Configuration    → model bazowy, kanał API
3. Agent Builder          → tożsamość i instrukcje ukierunkowane na redakcję dokumentów
4. Skills Manager         → skille korekty i streszczeń
5. Connectors Manager     → konektor do zewnętrznego słownika terminologicznego
6. Agent Builder          → pamięć ograniczona do poziomu projektu
7. Permissions Center     → uprawnienia w stanie wyjściowym (pełny dostęp)
8. Zapis                  → agent zapisany jako komponent własny
9. Agent Manager          → przypisanie do dwóch projektów redakcyjnych (moduł Workspace)
```

W każdym z dwóch projektów agent korzysta z osobnych instrukcji systemowych projektu i osobnej pamięci kontekstowej, zgodnie z rozdziałem 9.2 niniejszego dokumentu, ponieważ Użytkownik nie ustanowił żadnego powiązania na poziomie zasięgu „Para modułów” ani „Projekt” w oknie konfiguracji punktów izolacji.

### 13.2. Zespół agentów w rolach środowiska MultitaskingAI

Użytkownik buduje w module Apps kompletną aplikację z wykorzystaniem środowiska MultitaskingAI, rozwijając scenariusz E.1 Koncepcji platformy o warstwę agentową. W sekcji Role panelu orkiestracji przypisuje agentów do czterech ról.

| Rola | Przypisany agent | Kanał modelu / uprawnienia |
|---|---|---|
| Executor 1 | Agent Backend | Kanał CLI; uprawnienia w Permissions Center zawężone do modułów Developer i Terminal |
| Executor 2 | Agent Frontend | Kanał API |
| Coordinator | Agent ogólnego przeznaczenia | Bez zawężonych uprawnień |
| Executor 3 / Validator | Agent wcielający się w Security Auditor | Uprawnienia w Centrum uprawnień ograniczone wyłącznie do odczytu — bez dostępu do zapisu plików |

Każdy z dwóch Wykonawców uruchamia własny Subagent Network. Profil izolacji przypisany roli Security Auditor, ustawiony na poziomie zasięgu „Rola”, nadpisuje ustawienia zapisane w Centrum uprawnień tego agenta w każdym innym kontekście jego wykorzystania, zgodnie z rozdziałem 8.5 niniejszego dokumentu.

### 13.3. Agent z ograniczonymi uprawnieniami technicznymi dla pracy autonomicznej

Użytkownik konfiguruje w module Automations proces działający cyklicznie 24 godziny na dobę, wykorzystujący w środowisku MultitaskingAI agenta jako Executora 1. Ze względu na nienadzorowany charakter pracy Użytkownik włącza w Permissions Center tego agenta trzy z ośmiu zakresów izolacji technicznej.

| Zakres izolacji technicznej | Stan w tym scenariuszu |
|---|---|
| Katalog roboczy sesji | Włączony |
| Dostęp sieciowy | Włączony |
| Zakres odczytu i zapisu plików | Włączony |
| Środowisko procesu | Stan wyjściowy (wyłączony) |
| Katalog danych i konfiguracji modelu | Stan wyjściowy (wyłączony) |
| Konto i token per sesja | Stan wyjściowy (wyłączony) |
| Model procesu | Stan wyjściowy (wyłączony) |
| Serwer wykonania | Stan wyjściowy (wyłączony) |

Zawężenie to jest świadomą decyzją Użytkownika podjętą wyłącznie dla tego agenta i tej konfiguracji — pozostali agenci oraz ten sam agent wykorzystywany poza tym procesem zachowują stan wyjściowy pełnego dostępu, zgodnie z rozdziałem 6.4 niniejszego dokumentu.

---

## 14. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Użytkownik | Uczestnik pracy, który zleca zadania i zatwierdza ich wyniki; jedyne konto właściciela platformy. |
| Koordynator | Komponent orkiestrujący platformy: przyjmuje zlecenie, dekomponuje je na zadania, przydziela je Wykonawcom, nadzoruje realizację, prowadzi kontrolę jakości i zamyka zlecenie. |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania przydzielone przez Koordynatora. |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca; lewa kolumna obszaru roboczego, warstwa widoczności 1. |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca; kolumna sąsiadująca z Chat Window. |
| Pętla wykonawcza | Obieg ośmiu faz: przyjęcie zlecenia, dekompozycja na zadania, przydział zadań, wymiana komunikatów sterujących, wykonanie, kontrola jakości, decyzja o ponowieniu, zamknięcie zlecenia. |
| Warstwa widoczności | Jedna z czterech warstw ujawniania funkcjonalności interfejsu (zawsze widoczna, widoczna na żądanie, rozwinięcia kontekstowe, funkcje eksperckie), przypisywana każdemu elementowi interfejsu (rozdz. 3). |
| Agent | Jednostka AI o zdefiniowanej tożsamości, modelu bazowym, instrukcjach systemowych, umiejętnościach, wtyczkach, konektorach, pamięci i uprawnieniach, tworzona i konfigurowana jako komponent własny w module Agents; komponent platformowy działający ponad wszystkimi środowiskami, wykorzystywany jako Wykonawca w TalkIn, WorkSpace i CodeStudio, przypisywalny do projektów w module Workspace i do ról w środowisku MultitaskingAI. |
| Agent Builder | Nadrzędne okno operacyjne modułu Agents, w którym powstaje i jest edytowana definicja agenta, spinające Model Configuration, Skills Manager, Connectors Manager i Permissions Center w jeden zapisywalny komponent własny. |
| Model Configuration | Okno operacyjne modułu Agents ustalające model bazowy agenta oraz kanał jego połączenia (API, CLI, SSH, HTTP). |
| Skills Manager | Okno operacyjne modułu Agents służące do doboru i konfiguracji umiejętności agenta. |
| Connectors Manager | Okno operacyjne modułu Agents służące do podłączania wtyczek i konektorów, pochodzących ze źródła Danaco Plugin lub Personal. |
| Permissions Center (Centrum uprawnień) | Okno operacyjne modułu Agents ustalające zakres uprawnień agenta: dostęp do rozszerzeń, dostęp do modułów i zasobów, domyślne ustawienie ośmiu zakresów izolacji technicznej oraz zakres działania w środowisku MultitaskingAI; stan wyjściowy to pełny dostęp operacyjny. |
| Agent Manager | Okno operacyjne modułu Workspace, przez które agenci skonfigurowani w module Agents są przypisywani do projektu jako Wykonawcy w jego ramach. |
| Sekcja Role | Jedna z sześciu sekcji panelu orkiestracji środowiska MultitaskingAI, udostępniająca cztery okna robocze ról oraz mechanizm przypisania agentów z modułu Agents do tych ról. |
| Kanał modelu | Sposób połączenia platformy z modelem AI (API, CLI, SSH lub HTTP), wybierany per sesja lub per rola, zapisywany jako wartość domyślna w definicji agenta w Model Configuration. |
| Zakres izolacji technicznej | Jeden z ośmiu zakresów zdefiniowanych w rozdziale 12 Architektury (katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania), konfigurowalny dla agenta w Permissions Center i dla dowolnego poziomu zasięgu w oknie konfiguracji punktów izolacji. |

---

## 15. Kryteria odbioru

Opracowanie jest przyjmowane do budowy, gdy poniższe warunki są spełnione łącznie.

| Warunek sprawdzalny | Sposób sprawdzenia |
|---|---|
| Siedem komponentów definicji agenta opisanych w rozdziale 5 | Podrozdziały 5.1–5.7 obecne, po jednym na komponent, zgodnie z rozdz. 11.15 i Załącznikiem D.1 Koncepcji platformy |
| Obszary kontraktu `agent`, `subagent`, `advisor` istnieją w `contract.json` | Odczyt pola `nazwa` z listy `obszary` w `budowa/shared/contract.json` obejmuje wszystkie trzy wartości |
| Pętla wykonawcza opisana jako osiem faz kolejnych | Tabela rozdziału 2.4 zawiera osiem wierszy ponumerowanych 1–8, zamkniętych fazą kontroli jakości |
| Centrum uprawnień ma stan wyjściowy pełnego dostępu | Rozdział 6 i rozdział 12 stwierdzają wprost brak domyślnego zawężenia; zero sformułowań sugerujących odwrotny stan wyjściowy |
| Cztery warstwy widoczności przypisane każdemu elementowi zarządzania agentami | Tabela rozdziału 3 zawiera warstwy 1–4 z kompletem elementów i sposobem dostępu |
| Przypisanie agenta do roli i do projektu opisane jako mechanizmy odrębne i jawne | Rozdziały 8 i 9 wskazują osobne okna (sekcja Role, Agent Manager) oraz warunek jawnej decyzji użytkownika |
| Szablon definicji agenta zgodny z komponentami rozdziału 5 | Rozdział 10 wymienia komplet pól odpowiadających podrozdziałom 5.1–5.7 oraz Centrum uprawnień |
| Formy wizualne zajmują co najmniej 30% wierszy dokumentu | Kontrola ręczna zgodnie z rozdziałem 3.1 Standardu redakcyjnego i językowego |
| Objętość dokumentu nie mniejsza niż 85 000 znaków | liczba znaków pliku źródłowego (rozdz. 9 Standardu redakcyjnego i językowego), licznik uruchomiony na `docs/specyfikacje/specyfikacja-agentow.md` |

---

## Załącznik A. Przykład wypełnionej definicji agenta

Poniższy przykład wypełnia szablon z rozdziału 10 dla agenta ze scenariusza 13.2 niniejszego dokumentu.

| Pole | Wartość |
|---|---|
| Nazwa | Agent Backend |
| Przeznaczenie | Realizacja warstwy backendowej aplikacji budowanej w module Apps w roli Executor 1 środowiska MultitaskingAI. |
| Model bazowy | model ogólnego przeznaczenia o wysokiej jakości generowania kodu |
| Kanał modelu | CLI |
| Instrukcje systemowe | „Realizujesz zadania backendowe zgodnie z architekturą przekazaną przez Coordinatora. Zgłaszaj wyniki do kolejki po zakończeniu każdego etapu.” |
| Umiejętności | wzorce architektury API, konwencje kodowania repozytorium |
| Wtyczki i konektory | konektor repozytorium kodu (Danaco Plugin) |
| Pamięć | poziom projektu |
| Uprawnienia — rozszerzenia | pełny dostęp do konektora repozytorium kodu |
| Uprawnienia — moduły i zasoby | zawężone do modułów Developer i Terminal |
| Uprawnienia — izolacja techniczna | stan wyjściowy (brak aktywnych zakresów) |
| Uprawnienia — MultitaskingAI | Subagent Network: tak (Agent Backend, Agent API, Agent Database) |
| Moduł(y) zastosowania | CodeStudio; rola Executor 1 w MultitaskingAI |
| Powiązania | rola Executor 1 w zespole projektu „Budowa aplikacji X” (panel orkiestracji) |
| Widoczność / zasięg | globalny |

---

## Załącznik B. Macierz elementów definicji agenta

Zestawienie porządkuje siedem komponentów definicji agenta oraz Centrum uprawnień według okna konfiguracji, w którym powstają, i miejsca, z którego pochodzi ich podstawa merytoryczna.

| Element definicji | Rozdział niniejszego dokumentu | Rozdział Koncepcji platformy | Rozdział Architektury | Rozdział Modelu danych |
|---|---|---|---|---|
| Model bazowy i kanał | 5.1 | 11.15 | 9 | 11.1, 12.1 |
| Tożsamość | 5.2 | 11.15, Załącznik D.1 | — | 12.1 |
| Instrukcje systemowe | 5.3 | 11.15 | — | 12.1 |
| Umiejętności | 5.4 | 11.15 | 14 | 12.1, 12.2 |
| Wtyczki i konektory | 5.5 | 11.15 | 14 | 12.1, 12.3 |
| Pamięć | 5.6 | 11.15, 6, 2.4, 8 | 10, 12 | 9, 12.1 |
| Uprawnienia (Centrum uprawnień) | 6 | 11.15, 6.4–6.6 | 12, 15 | 12.4 |
| Przypisanie do roli | 8 | 13.3, 13.8 | 8, 10 | 13.2 |
| Przypisanie do projektu | 9 | 11.2, 11.15 | — | 7.1, 6.5 |

---

*Koniec dokumentu. Specyfikacja agentów — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
