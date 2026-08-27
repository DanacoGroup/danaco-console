# Danaco Console — Okno Konfiguracji

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
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) i Deweloper (co zbudować) |
| **Przeznaczenie** | Specyfikacja projektowa okna Konfiguracji — panelu programującego nakładkę na model bazowy, pełnego zestawu pól sterujących wywołaniem modelu, konfiguracji obu kanałów komunikacji operacyjnej, konfiguracji warstw widoczności oraz Panelu prowenancji wywołania |
| **Źródło** | Koncepcja platformy (autorytatywna), Model konfiguracji, Architektura techniczna, Integracja modeli, Specyfikacja agentów, Izolacja i konfigurowalność zależności, Specyfikacja okien operacyjnych, Bezpieczeństwo i uwierzytelnianie, System wizualny, Rozszerzenia |
| **Autor koncepcji** | Operator (Dariusz Naharnowicz, Danaco Group) |
| **Opracowanie** | Danaco Console — projekt UI |

Dokument specyfikuje okno Konfiguracji Danaco Console — okno o statusie serca produktu, ponieważ to z niego programuje się nakładkę uruchamianą ponad modelem bazowym przy każdym jego wywołaniu. Okno skupia w jednym miejscu trzy zakresy nakładki ustalone w Modelu konfiguracji („Zachowanie modeli” 5.4, „Tożsamość modeli” 5.5, „Prompty systemowe” 5.6), scalając je w jeden dedykowany obszar roboczy **Nakładka i wywołanie modelu**, oraz udostępnia **Panel prowenancji**, dający Operatorowi dokładny, weryfikowalny podgląd tego, co faktycznie zostało przekazane modelowi przy danym wywołaniu.

---

## Spis treści

- [Wprowadzenie i pozycja dokumentu](#wprowadzenie-i-pozycja-dokumentu)
1. [Okno Konfiguracji jako serce produktu](#1-okno-konfiguracji-jako-serce-produktu)
2. [Warstwowa konfiguracja nakładki: Konstytucja → Profil/Rola → Ekspertyza zadaniowa](#2-warstwowa-konfiguracja-nakładki-konstytucja--profilrola--ekspertyza-zadaniowa)
3. [Pełny zestaw pól sterujących wywołaniem modelu](#3-pełny-zestaw-pól-sterujących-wywołaniem-modelu)
4. [Profile ekspertów](#4-profile-ekspertów)
5. [Dane dostępowe i zmienne środowiskowe](#5-dane-dostępowe-i-zmienne-środowiskowe)
6. [Zależności i izolacja jako ustawienia konfiguracyjne](#6-zależności-i-izolacja-jako-ustawienia-konfiguracyjne)
7. [Uwierzytelnianie jako ustawienie konfiguracyjne](#7-uwierzytelnianie-jako-ustawienie-konfiguracyjne)
8. [Konfiguracja kanałów komunikacji: Chat Window i Execution Loop Window](#8-konfiguracja-kanałów-komunikacji-chat-window-i-execution-loop-window)
9. [Konfiguracja warstw widoczności](#9-konfiguracja-warstw-widoczności)
10. [Panel prowenancji — „co poszło do modelu”](#10-panel-prowenancji--co-poszło-do-modelu)
11. [Import i eksport konfiguracji](#11-import-i-eksport-konfiguracji)
12. [Makieta tekstowa okna](#12-makieta-tekstowa-okna)
13. [Katalog elementów interfejsu](#13-katalog-elementów-interfejsu)
14. [Stany](#14-stany)
15. [Zachowanie i komunikacja klient–serwer](#15-zachowanie-i-komunikacja-klientserwer)
16. [Zgodność z zasadami nadrzędnymi platformy](#16-zgodność-z-zasadami-nadrzędnymi-platformy)
17. [Scenariusze użycia](#17-scenariusze-użycia)
18. [Słowniczek pojęć](#18-słowniczek-pojęć)
- [Załącznik A. Zbiorczy szablon konfiguracji nakładki](#załącznik-a-zbiorczy-szablon-konfiguracji-nakładki)
- [Załącznik B. Katalog ikon zastosowanych w oknie](#załącznik-b-katalog-ikon-zastosowanych-w-oknie)
- [Załącznik C. Mapowanie pól na komunikaty WebSocket](#załącznik-c-mapowanie-pól-na-komunikaty-websocket)

---

## Wprowadzenie i pozycja dokumentu

Okno Konfiguracji jest, obok okna operacyjnego, jednym z dwóch okien realizujących interfejs użytkownika platformy (Architektura, rozdz. 13) i jedynym miejscem, z którego konfiguruje się wszystkie zależności i zachowania platformy (Koncepcja platformy, rozdz. 6 i rozdz. 14, zasada 2). Trzynaście zakresów ustawień tego okna oraz jego warstwowość (globalna → środowisko → projekt → sesja) ustala Model konfiguracji; niniejszy dokument osadza w nich obszar roboczy, który programuje model przy każdym jego wywołaniu, wraz z konfiguracją obu kanałów komunikacji operacyjnej i konfiguracją warstw widoczności interfejsu.

| Dokument źródłowy | Co wnosi do niniejszej specyfikacji | Rozdział wykorzystany tutaj |
|---|---|---|
| Koncepcja platformy | Zasada pełnej konfigurowalności; zasada braku twardych blokad; komponent własny Agent; role środowiska MultitaskingAI | rozdz. 4, 9.15, 11.3, 12 |
| Model konfiguracji | Struktura okna Konfiguracji, trzynaście zakresów, warstwowość ogólna, mechanizm objaśnień kontekstowych, wzorzec okna trzykolumnowego (Izolacja) | rozdz. 3–6 |
| Architektura techniczna | Cztery kanały integracji (API, CLI, SSH, HTTP), osiem zakresów izolacji technicznej, warstwowość konfiguracji, uwierzytelnianie | rozdz. 6, 9, 10, 12 |
| Integracja modeli | Warstwa dostawcy modelu, mechanizm adapterów, kanał modelu jako encja, tożsamość i persona modelu, bezpieczeństwo danych dostępowych | rozdz. 2–8, 11 |
| Specyfikacja agentów | Siedem komponentów definicji agenta, Permissions Center, precedencja definicja agenta ↔ poziomy zasięgu | rozdz. 3, 4, 6 |
| Izolacja i konfigurowalność zależności | Osiem zakresów izolacji technicznej z objaśnieniami, siedem poziomów zasięgu, profile izolacji | rozdz. 3, 6, 7 |
| Bezpieczeństwo i uwierzytelnianie | Wymóg logowania jako ustawienie Operatora, token urządzenia, relacja uwierzytelnianie ↔ izolacja | rozdz. 8–10, 13 |
| Rozszerzenia | Cztery rodzaje rozszerzeń, rejestr rozszerzeń, relacja rozszerzenie ↔ uprawnienia | rozdz. 1, 6–8 |
| System wizualny | Biblioteka komponentów (`.dn-`), tokeny stanów, wzorzec zastosowania w oknie izolacji | rozdz. 6–8, 12 |

### Zakres dokumentu

| Element zakresu | Ujęty | Nieujęty |
|---|---|---|
| Warstwowa konfiguracja nakładki (Konstytucja → Profil/Rola → Ekspertyza zadaniowa) wraz z pełnym zestawem pól treści nakładki | Tak — rozdz. 2 | — |
| Pełny zestaw pól sterujących wywołaniem modelu: system prompt jako podgląd wynikowy, harness, hooki, narzędzia dozwolone i zabronione, model, effort, kanał, parametry próbkowania, pula kont | Tak — rozdz. 3 | — |
| Profile ekspertów wraz z wersjonowaniem, porównaniem wersji i galerią | Tak — rozdz. 4 | — |
| Dane dostępowe, rejestr kluczy, zmienne środowiskowe procesu sesji | Tak — rozdz. 5 | — |
| Zależności i izolacja jako ustawienia konfiguracyjne | Tak — rozdz. 6 (odesłanie do okna Izolacji) | Ponowne wyprowadzenie mechanizmu izolacji — przedmiot Modelu konfiguracji rozdz. 6 i dokumentu Izolacja i konfigurowalność zależności |
| Uwierzytelnianie jako ustawienie konfiguracyjne | Tak — rozdz. 7 | Tożsamość Operatora, metody i urządzenia — przedmiot okna Ustawień |
| Konfiguracja obu kanałów komunikacji operacyjnej: Chat Window i Execution Loop Window | Tak — rozdz. 8 | — |
| Konfiguracja warstw widoczności interfejsu | Tak — rozdz. 9 | — |
| Panel prowenancji wraz z narzędziami analizy, porównania i eksportu | Tak — rozdz. 10 | — |
| Import i eksport konfiguracji, migawki stanu | Tak — rozdz. 11 | — |
| Makieta tekstowa, katalog elementów, stany | Tak — rozdz. 12–14 | — |
| Wymiary, tokeny kolorów i typografii, ostateczny układ graficzny | — | Przedmiot prac wykonawczych, zgodnie z zastrzeżeniem rozdz. 1 Modelu konfiguracji |
| Pozostałe zakresy okna Konfiguracji (Aplikacja, Procesy, Akcje, Rozszerzenia, Integracje, Historia, Pamięć, Izolacja, Karty sesji, Komponenty własne) | Przywołane w zakresie relacji do nakładki (rozdz. 6) | Pełna specyfikacja pozostaje przedmiotem Modelu konfiguracji |

Odbiorcami dokumentu są Designer, dla którego rozdziały 12–14 rozstrzygają co, gdzie i w jakiej formie znajduje się na ekranie, oraz Deweloper, dla którego rozdziały 2–11 i 15 oraz załączniki rozstrzygają, jaki mechanizm i jaki kontrakt komunikacji zostaje zbudowany.

Zasada wiążąca dla całości okna: żadna pozycja konfiguracji nie wprowadza twardej blokady. Brak ustawienia oznacza dziedziczenie albo wartość domyślną, przyciski pozostają aktywne, walidacja ma charakter wyłącznie ostrzegawczy, dane dostępowe są domyślnie jawne, a uwierzytelnianie i izolacja są ustawieniami konfiguracyjnymi, nie wymogami.

---

## 1. Okno Konfiguracji jako serce produktu

### 1.1. Dlaczego serce produktu

Danaco Console jest systemem operacyjnym do pracy z modelami sztucznej inteligencji (Koncepcja platformy, rozdz. 1) — każde zadanie wykonywane w dowolnym środowisku, module czy roli sprowadza się ostatecznie do wywołania modelu. Okno Konfiguracji jest jedynym miejscem, w którym Operator ustala, **czym ten model jest** w chwili wywołania: jaką tożsamość otrzymuje, jakimi zasadami się kieruje, jakich narzędzi może użyć, którym kanałem jest osiągany i ile wysiłku obliczeniowego angażuje. Poniższa tabela zestawia dwa okna aplikacji klienckiej i wskazuje, dlaczego to okno Konfiguracji — a nie okno operacyjne — jest miejscem programowania nakładki.

| Okno aplikacji klienckiej | Rola | Relacja do nakładki na model |
|---|---|---|
| Okno operacyjne (Architektura, rozdz. 13) | Główna przestrzeń pracy — Chat Window, Execution Loop Window, obszar roboczy modułu | Konsumuje efekt nakładki przy każdym wywołaniu modelu; udostępnia uproszczony, kontekstowy skrót do warstwy sesji (rozdz. 3.1 Modelu konfiguracji) |
| Okno Konfiguracji (Architektura, rozdz. 13) | Jedyne miejsce ustalania wszystkich zależności i zachowań platformy | Programuje nakładkę — miejsce, w którym powstaje to, co okno operacyjne wykonuje |

### 1.2. Trzynaście zakresów okna Konfiguracji — struktura

Model konfiguracji (rozdz. 5) ustala trzynaście zakresów ustawień okna Konfiguracji.

| # | Zakres ustawień | Co konfiguruje |
|---|---|---|
| 5.1 | Aplikacja | Platforma jako całość: motyw, język, uwierzytelnianie, urządzenia, powiadomienia, warstwy widoczności |
| 5.2 | Procesy | Osiem zakresów izolacji technicznej procesu sesji |
| 5.3 | Akcje | Silnik kolejek i zadania modułu Automations, orkiestracja MultitaskingAI |
| 5.4 | Zachowanie modeli | Kanał połączenia z modelem, dobór modelu, parametry wywołania |
| 5.5 | Tożsamość modeli | Tożsamość agentów i wcielenia ról MultitaskingAI |
| 5.6 | Prompty systemowe | Instrukcje systemowe projektu i agenta, budowa promptów przez Koordynatora |
| 5.7 | Rozszerzenia | Wtyczki, umiejętności, konektory, serwery MCP |
| 5.8 | Integracje | Jawne powiązania między środowiskami, modułami, komponentami własnymi |
| 5.9 | Historia | Przechowywanie i dostępność historii rozmów i sesji |
| 5.10 | Pamięć | Zasoby pamięci kontekstowej na czterech poziomach |
| 5.11 | Izolacja | Izolacja kontekstu i izolacja techniczna — okno trzykolumnowe |
| 5.12 | Karty sesji | Ustawienia pojedynczej karty sesji |
| 5.13 | Komponenty własne | Tworzenie i przypisanie automatyki, agenta, projektu, profilu asystenta |

### 1.3. Nakładka i wywołanie modelu — scalony obszar roboczy

Trzy spośród trzynastu zakresów — 5.4, 5.5 i 5.6 — opisują wspólnie jedną rzecz: to, czym model jest i jak zostaje osiągnięty w chwili wywołania. Rozdzielenie ich na trzy osobne pozycje nawigacji jest poprawne porządkowo, lecz operacyjnie niewygodne — Operator programujący zachowanie konkretnego agenta lub roli ustala wszystkie trzy jednocześnie i widzi wynik ich złożenia. Model konfiguracji rozwiązuje analogiczny problem dla zakresu „Izolacja” (5.11), wydzielając dla niego własne, trzykolumnowe okno (rozdz. 6 Modelu konfiguracji) — ze względu na własną złożoność, a nie dlatego, że przestaje być jednym z trzynastu zakresów.

Ten sam zabieg obejmuje zakresy 5.4–5.6, scalając je w jeden obszar roboczy: **Nakładka i wywołanie modelu**. Zabieg nie tworzy czternastego zakresu ustawień ponad trzynaście ustalonych — jest ich operacyjnym rozwinięciem, analogicznym do rozdziału 6 Modelu konfiguracji.

```
ANALOGIA STRUKTURALNA

  Model konfiguracji, rozdz. 6                Niniejszy dokument
  ──────────────────────────────              ──────────────────────────────
  Zakres 5.11 „Izolacja”                       Zakresy 5.4 + 5.5 + 5.6
      │  (własna złożoność:                        │  (własna złożoność:
      │   2 rodzaje izolacji,                       │   3 warstwy nakładki,
      │   7 poziomów zasięgu,                       │   pola wywołania,
      │   profile)                                  │   3 tryby wywołania,
      ▼                                             │   Panel prowenancji)
  Okno konfiguracji punktów izolacji                ▼
  (trzykolumnowe — rozdz. 6.1 MK)               Nakładka i wywołanie modelu
                                                 (trzykolumnowe — rozdz. 12.2)
```

Pozycja w nawigacji zakresów (rozdz. 3.2 Modelu konfiguracji) pozostaje potrójna — Operator wchodzi osobno w „Zachowanie modeli”, „Tożsamość modeli” lub „Prompty systemowe” — lecz wszystkie trzy wejścia otwierają ten sam obszar roboczy, z zaznaczoną odpowiednią sekcją, dokładnie jak pozycja „Izolacja” otwiera trzykolumnowe okno punktów izolacji.

---

## 2. Warstwowa konfiguracja nakładki: Konstytucja → Profil/Rola → Ekspertyza zadaniowa

### 2.1. Trzy warstwy nakładki

Nakładka jest tekstem instrukcji systemowej i towarzyszącym mu zestawem ustawień wywołania, faktycznie przekazywanym modelowi bazowemu przy każdym wywołaniu. Powstaje przez złożenie trzech warstw, uporządkowanych od najszerszej i najbardziej stabilnej do najwęższej i najbardziej ulotnej.

| Warstwa | Czym jest | Stabilność | Pytanie, na które odpowiada |
|---|---|---|---|
| Konstytucja | Trwałe, firmowe zasady pracy AI, obowiązujące każde wywołanie modelu na platformie, niezależnie od środowiska, modułu, roli czy agenta | Najwyższa — zmienia się rzadko, świadomie | „Jakimi zasadami kieruje się AI w tej organizacji, zawsze?” |
| Profil / Rola | Tożsamość Wykonawcy: agent (komponent własny modułu Agents) albo rola środowiska MultitaskingAI albo profil asystenta modułu Assistant | Średnia — zmienia się przy budowie lub edycji agenta albo przy przeobsadzeniu roli | „Kim — w sensie operacyjnym — jest Wykonawca tego zadania?” |
| Ekspertyza zadaniowa | Wiedza proceduralna i instrukcje właściwe konkretnemu zadaniu, projektowi lub krokowi procesu — statyczna (Instructions Panel, skille) albo budowana dynamicznie przez Koordynatora | Najniższa — zmienia się per zadanie, czasem per wiadomość | „Co dokładnie trzeba wiedzieć, aby wykonać to zadanie?” |

Trzy warstwy są interpretacyjną ramą nadaną trzem zakresom (5.4–5.6) oraz Konstytucji, która otrzymuje w tym oknie swoje właściwe miejsce. Poniższa tabela mapuje warstwy na zakresy i mechanizmy źródłowe.

| Warstwa nakładki | Zakres okna Konfiguracji | Mechanizm źródłowy |
|---|---|---|
| Konstytucja | Prompty systemowe (5.6) — poziom nadrzędny wobec instrukcji projektu i agenta | Edytor własny okna Konfiguracji |
| Profil / Rola | Tożsamość modeli (5.5) | Tożsamość i persona modelu (Integracja modeli, rozdz. 6); definicja agenta (Specyfikacja agentów, rozdz. 5) |
| Ekspertyza zadaniowa | Prompty systemowe (5.6) | Instrukcje systemowe projektu (Specyfikacja agentów, rozdz. 5.3); skille (rozdz. 3.4); budowa promptów przez Koordynatora (Model konfiguracji, 5.6) |

### 2.2. Gdzie każda warstwa się ustala

Zasada pełnej kompozycyjności (Koncepcja platformy, rozdz. 14, zasada 1) nakazuje nie duplikować edytorów tam, gdzie już istnieją. Konstytucja ma w tym oknie swój jedyny dom. Profil/Rola i Ekspertyza zadaniowa mają domy ustalone we wcześniejszych dokumentach; okno Nakładka i wywołanie modelu pokazuje je w podglądzie tylko do odczytu, z bezpośrednim skrótem do właściwego edytora, oraz dodaje pola wywołania (rozdz. 3) właściwe wybranej warstwie zasięgu.

| Warstwa | Gdzie się tworzy i edytuje treść | Jak widoczna w oknie Nakładka i wywołanie modelu |
|---|---|---|
| Konstytucja | Tutaj — okno Konfiguracji, obszar Nakładka i wywołanie modelu, warstwa zasięgu „Globalna” lub „Środowisko” | Pełny edytor tekstowy (pole własne tego okna) |
| Profil / Rola | Agent Builder (moduł Agents) — dla agenta; sekcja Role panelu orkiestracji — dla roli MultitaskingAI bez agenta; Voice Console / konfiguracja profilu — dla Assistant | Podgląd tożsamości i instrukcji systemowych, tylko do odczytu, ze skrótem „Edytuj w Agent Builder →” / „Edytuj w sekcji Role →” |
| Ekspertyza zadaniowa (statyczna) | Instructions Panel (moduł Workspace) — dla projektu; Skills Manager (moduł Agents) — dla agenta | Podgląd, tylko do odczytu, ze skrótem „Edytuj w Instructions Panel →” |
| Ekspertyza zadaniowa (dynamiczna) | Budowana w czasie rzeczywistym przez Koordynatora (Model konfiguracji, 5.6) | Widoczna retrospektywnie, w Panelu prowenancji (rozdz. 10) danego wywołania |

### 2.3. Kompozycja nakładki w jeden prompt systemowy

Trzy warstwy łączą się w jeden, ciągły tekst przekazywany modelowi, zgodnie z zasadą pełnej kompozycyjności zastosowaną do instrukcji projektu i agenta („gdy agent zostaje przypisany do projektu… oba zestawy instrukcji obowiązują łącznie”, Specyfikacja agentów, rozdz. 5.3). Konstytucja nie jest nadpisywana przez warstwy węższe — jest do niej dokładana kolejna treść, aż do gotowego, wysyłanego tekstu.

```
KOMPOZYCJA NAKŁADKI (kolejność składania — od najszerszej do najwęższej)

  ┌─────────────────────────────────────────────────────────────┐
  │ 1. KONSTYTUCJA                                                │
  │    (globalna albo per środowisko; przy braku treści —          │
  │     sekcja pomijana w całości, nie wstawiana pusta)            │
  └───────────────────────────┬─────────────────────────────────┘
                              │  dokładana, nie nadpisywana
                              ▼
  ┌─────────────────────────────────────────────────────────────┐
  │ 2. PROFIL / ROLA                                              │
  │    tożsamość + instrukcje systemowe agenta lub roli           │
  │    + styl i ton + negatyw + słowniczek terminów               │
  └───────────────────────────┬─────────────────────────────────┘
                              │  dokładana, nie nadpisywana
                              ▼
  ┌─────────────────────────────────────────────────────────────┐
  │ 3. EKSPERTYZA ZADANIOWA                                       │
  │    instrukcje projektu + skille + format wyjścia + tekst      │
  │    zbudowany dynamicznie przez Koordynatora dla tego kroku    │
  └───────────────────────────┬─────────────────────────────────┘
                              │
                              ▼
                  PROMPT SYSTEMOWY WYSŁANY DO MODELU
             (dokładny zapis — Panel prowenancji, rozdz. 10)
```

Sekcja pominięta w całości przy braku ustawienia (a nie wstawiona jako pusty nagłówek) jest zastosowaniem zasady „brak ustawienia = wartość domyślna” (Koncepcja platformy, rozdz. 6.5) do treści promptu: agent bez przypisanej Konstytucji działa dokładnie tak, jakby ta warstwa nie istniała, a nie tak, jakby otrzymał pusty, ale obecny nagłówek.

### 2.4. Precedencja i dziedziczenie

Definicja agenta (Permissions Center, Agent Builder) niesie ze sobą wartość wyjściową warstwy Profil/Rola — dokładnie tak, jak niesie wartość wyjściową dla ośmiu zakresów izolacji technicznej (Specyfikacja agentów, rozdz. 6.3). Poziomy zasięgu konfiguracji ogólnej (rozdz. 4.1 Modelu konfiguracji) nadpisują tę wartość wyjściową w kolejności rosnącego pierwszeństwa, aż do poziomu „Rola”, który ma pierwszeństwo najwyższe — mechanizm identyczny z ustalonym dla Centrum uprawnień.

```
DEFINICJA AGENTA                              NAKŁADKA I WYWOŁANIE MODELU
Profil / Rola (Agent Builder)                 (niniejsze okno)
(wartość wyjściowa
 właściwa agentowi lub roli)
        │
        │            zasięg: Globalna (Konstytucja) ─────────────┤  (najniższe pierwszeństwo)
        │            zasięg: Środowisko (Konstytucja per środ.) ──┤
        │            zasięg: Projekt (Ekspertyza zadaniowa) ─────────┤
        │            zasięg: Sesja (doraźna zmiana) ───────────────────┤
        └───────────► zasięg: Rola (MultitaskingAI, nadpisanie) ───────────┤  (najwyższe pierwszeństwo)
                                                                     ▼
                                                          NAKŁADKA EFEKTYWNA
                                                   (podgląd — Panel prowenancji, rozdz. 10)
```

Brak ustawienia na danym poziomie zasięgu oznacza dziedziczenie z poziomu bezpośrednio szerszego, aż do wartości wyjściowej niesionej przez definicję agenta lub roli — zastosowanie tej samej reguły „brak ustawienia = wartość domyślna”, którą Model konfiguracji ustanawia dla wszystkich trzynastu zakresów (rozdz. 4.2).

### 2.5. Relacja do ogólnej warstwowości platformy

Cztery warstwy ogólne (globalna, środowisko, projekt, sesja — rozdz. 4.1 Modelu konfiguracji) pozostają wspólnym szkieletem.

| Warstwa ogólna (rozdz. 4.1 MK) | Konstytucja | Profil / Rola | Ekspertyza zadaniowa |
|---|---|---|---|
| Globalna | Tak — poziom właściwy | — | — |
| Środowisko | Tak — doprecyzowanie per środowisko | — | — |
| Projekt | — | — | Tak — Instructions Panel |
| Sesja | Tak — doraźna zmiana warstwy sesji | Tak — doraźna zmiana | Tak — doraźna zmiana |
| Rola (MultitaskingAI) | — | Tak — nadpisanie definicji agenta | Tak — budowa dynamiczna przez Koordynatora |
| Komponent własny (Agent) | — | Tak — wartość wyjściowa (Agent Builder) | Tak — skille agenta |

### 2.6. Diagram warstw nakładki

```
╔══════════════════════════════════════════════════════════════════════╗
║  WARSTWA 1 · KONSTYTUCJA                                               ║
║  zasięg: globalna · środowisko                                         ║
║  „jakimi zasadami kieruje się AI w tej organizacji, zawsze?”           ║
║  edytor: TUTAJ (okno Nakładka i wywołanie modelu)                      ║
╠══════════════════════════════════════════════════════════════════════╣
║  WARSTWA 2 · PROFIL / ROLA                                             ║
║  zasięg: komponent własny (agent) · rola MultitaskingAI                ║
║  „kim jest Wykonawca tego zadania?”                                    ║
║  edytor: Agent Builder · sekcja Role panelu orkiestracji               ║
╠══════════════════════════════════════════════════════════════════════╣
║  WARSTWA 3 · EKSPERTYZA ZADANIOWA                                      ║
║  zasięg: projekt · sesja · (dynamicznie) krok procesu Wykonawcy        ║
║  „co dokładnie trzeba wiedzieć, aby wykonać to zadanie?”                ║
║  edytor: Instructions Panel · Skills Manager · budowa przez Koordynatora║
╚══════════════════════════════════════════════════════════════════════╝
   Kompozycja (rozdz. 2.3) → PROMPT SYSTEMOWY → Panel prowenancji (rozdz. 10)
```

### 2.7. Pola treści nakładki

Warstwy nakładki składają się z nazwanych pól sterujących — każde z własnym zasięgiem warstwowym i objaśnieniem kontekstowym `[?]`. Wszystkie pola wchodzą do kompozycji promptu systemowego (rozdz. 2.3) albo do bloku `settings` Panelu prowenancji (rozdz. 10.2) i wszystkie podlegają dziedziczeniu warstw.

| Pole | Co ustala | Efekt | Zasięgi |
|---|---|---|---|
| Negatyw (czego unikać) | Osobna sekcja „nie rób” dokładana do nakładki, jawnie oddzielona od zasad pozytywnych | Ogranicza powtarzalne błędy bez powiększania treści Konstytucji; podlega odrębnemu wersjonowaniu | globalna, środowisko, projekt, rola, agent |
| Styl i ton wypowiedzi | Zwięzły zestaw ustawień stylu („zwięźle, urzędowo, bez emoji”) składany jako oddzielna warstwa tonu | Rozdziela „co robić” od „jak brzmieć”; zmiana tonu nie dotyka merytoryki | środowisko, projekt, sesja, rola, agent |
| Słowniczek terminów projektu | Lista par „termin → znaczenie” wstrzykiwana do nakładki (skróty kancelarii, nazwy spraw) | Model używa terminologii organizacji od pierwszego promptu | projekt, agent |
| Format i ograniczenia wyjścia | Format odpowiedzi (tekst / Markdown / JSON wg schematu), limit długości, wymagane sekcje | Powtarzalna, przewidywalna struktura odpowiedzi | sesja, rola, agent |
| Język odpowiedzi | Język wyjścia niezależny od języka promptu; wartość domyślna dziedziczy z języka interfejsu | Spójność językowa w środowiskach wielojęzycznych | globalna, środowisko, sesja, rola |
| Budżet kontekstu i tokenów | Miękkie budżety: docelowa długość kontekstu, limit tokenów odpowiedzi — parametr, nie bramka | Kontrola kosztu i czasu; przekroczenie skutkuje ostrzeżeniem w Panelu prowenancji, nie odrzuceniem wywołania | sesja, rola, agent |
| Kolejność składania warstw | Podgląd i zmiana kolejności sekcji nakładki; wartość domyślna zgodna z rozdz. 2.3 | Elastyczność przy zachowaniu czytelnego porządku domyślnego | rola, agent |

Każde z pól ma treść objaśnienia `[?]` opisującą działanie i wpływ ustawienia — treść objaśnienia jest częścią definicji ustawienia (Model konfiguracji, rozdz. 3.3), nie elementem dodatkowym.

### 2.8. Sterowanie zasięgiem pola

Mechanizmy wspólne dla wszystkich pól nakładki i pól wywołania — realizują warstwowość (Model konfiguracji, rozdz. 4) i regułę „pierwsza ustawiona warstwa wygrywa”.

| Mechanizm | Co robi | Efekt |
|---|---|---|
| Znacznik warstwy per pole | Przy każdym polu widoczna plakietka „ustawione na: <warstwa>” albo „dziedziczone z: <warstwa>”, spójna ze stanami rozdz. 14.2 | Operator zawsze wie, skąd pochodzi bieżąca wartość |
| Operacje „Przypnij / Odepnij / Przywróć dziedziczenie” | Trzy akcje przy polu: przypięcie wartości do bieżącej warstwy, usunięcie nadpisania, przywrócenie wartości dziedziczonej | Jawne, odwracalne sterowanie zasięgiem; odpięcie nie kasuje warstwy szerszej |
| Podgląd łańcucha dziedziczenia | Rozwijany widok „skąd bierze się ta wartość”: globalna → środowisko → projekt → sesja → rola/agent, z zaznaczeniem warstwy zwycięskiej | Diagnostyka konfliktów warstw przed zmianą — odpowiednik „polityki efektywnej” okna izolacji |
| Tryb różnicowy warstwy | Widok pokazujący wyłącznie pola nadpisane na bieżącej warstwie względem dziedziczenia | Audyt „co ta sesja lub rola faktycznie zmienia” |
| Warstwa „Domyślne platformy” | Selektor warstwy zawiera najniższą, tylko do odczytu pozycję „Domyślne platformy”, prezentującą wartości wbudowane | Jawny stan bazowy; wzmocnienie reguły „brak ustawienia = wartość domyślna” |

---

## 3. Pełny zestaw pól sterujących wywołaniem modelu

### 3.0. Zestawienie ogólne

Poza treścią nakładki (rozdz. 2) — dostępną w podglądzie jako pole „System prompt” (rozdz. 3.1) — wywołanie modelu wymaga siedmiu dalszych pól. Ten sam zestaw siedmiu pól występuje identycznie w bloku „settings” Panelu prowenancji (rozdz. 10.2), w diagramie przepływu danych (rozdz. 10.4) i w zestawieniu zgodności z zasadami nadrzędnymi (rozdz. 16). Wszystkie podlegają zasadzie „brak ustawienia = wartość domyślna” (Koncepcja platformy, rozdz. 6.5) i mechanizmowi objaśnień kontekstowych `[?]` (rozdz. 3.9).

| Pole | Co ustala | Warstwy | Zależność od kanału | Wartość domyślna |
|---|---|---|---|---|
| Harness | Warstwa pętli agentowej opakowująca wywołania modelu | sesja, rola, agent | Wartość i edytowalność zależą od kanału | Zależna od kanału |
| Hooki | Punkty zaczepienia wykonujące automatykę w toku cyklu wywołania | sesja, rola, agent | Dostępne wprost dla Code CLI i Agent SDK; dla API realizowane przez rdzeń | Brak podłączonych hooków |
| Narzędzia dozwolone i zabronione | Zestaw narzędzi i akcji dostępnych modelowi w toku pracy | sesja, rola, agent (dziedziczy z Permissions Center) | Wspólne dla wszystkich kanałów | Pełny dostęp (zgodnie z Permissions Center, Specyfikacja agentów 4.4) |
| Model | Konkretny model bazowy udostępniany przez wybrany kanał | sesja, rola, agent | Lista modeli zależna od kanału | Zależna od kanału |
| Effort | Poziom wysiłku obliczeniowego (głębokości rozumowania) modelu | sesja, rola, agent | Dostępne dla modeli i kanałów wspierających tryb rozszerzonego rozumowania | wyłączony (dziedziczy wartość modelu) |
| Kanał | Sposób technicznego dotarcia do modelu: Code CLI, Agent SDK, API | sesja, rola, agent | — | Zależna od dostawcy |
| Pula kont | Zbiór odrębnie uwierzytelnionych kont Code CLI, między którymi rozdzielane są wywołania | sesja, rola, agent | Dostępne wyłącznie dla kanału Code CLI | Jedno, domyślnie skonfigurowane konto |

### 3.1. System prompt

Pole nie jest edytowane bezpośrednio jako całość — jest **wynikiem** kompozycji opisanej w rozdziale 2. W oknie Nakładka i wywołanie modelu pole „System prompt” występuje w dwóch postaciach: trzy oddzielne edytory i podglądy warstw (rozdz. 2.2) oraz jeden, nieedytowalny podgląd wynikowy — identyczny z blokiem „system prompt” Panelu prowenancji (rozdz. 10.2) — pozwalający zweryfikować złożenie przed zapisaniem zmiany.

> **[?] Objaśnienie kontekstowe.** „System prompt to tekst przekazywany modelowi przed każdą wiadomością w tej sesji lub roli — wynik złożenia Konstytucji, Profilu/Roli i Ekspertyzy zadaniowej. Zmiana dowolnej z trzech warstw zmienia wynikowy tekst; podgląd w kolumnie prawej pokazuje efekt złożenia, zanim zostanie wysłany do modelu.”

### 3.2. Harness

Harness jest warstwą oprogramowania, która przyjmuje polecenie, przekazuje je modelowi wraz ze skomponowaną nakładką, obsługuje zgłoszone przez model wywołania narzędzi, przekazuje ich wyniki z powrotem do modelu i powtarza ten cykl aż do odpowiedzi końcowej. Harness nie jest jednym, jednolitym komponentem platformy — jego pochodzenie i zakres edytowalności zależą od wybranego kanału (rozdz. 3.7).

| Kanał | Kto dostarcza harness | Co jest tutaj edytowalne |
|---|---|---|
| Code CLI | Narzędzie wiersza poleceń dostawcy — pętla agentowa, obsługa narzędzi i hooków są jego wewnętrzną częścią | Argumenty startowe przekazywane narzędziu (rozdz. 10.2, blok argv); sam harness nie jest tu konfigurowalny wprost |
| Agent SDK | Biblioteka dostawcy osadzona w rdzeniu platformy | Parametry udostępniane przez bibliotekę: limit kroków pętli, strategia ponawiania, format komunikatów |
| API | Rdzeń platformy Danaco Console (Go) — pętla wywołań narzędzi budowana i utrzymywana samodzielnie (Architektura, rozdz. 10) | Parametry pętli własnej rdzenia: limit kroków, limit czasu pojedynczego kroku, zachowanie przy błędzie narzędzia |

> **[?] Objaśnienie kontekstowe.** „Harness to mechanizm spinający pojedyncze wywołanie modelu w wieloetapową pętlę: model → narzędzie → wynik → model, aż do odpowiedzi końcowej. Jego pochodzenie zależy od wybranego kanału — dla Code CLI jest wbudowany w narzędzie, dla Agent SDK pochodzi z biblioteki dostawcy, dla API jest budowany przez rdzeń platformy. Zmiana kanału zmienia dostępne tu parametry harness.”

### 3.3. Hooki

Hooki są punktami zaczepienia w cyklu wywołania modelu, w których Operator podłącza własne działanie — najczęściej odwołanie do automatyki utworzonej w module Automations (Koncepcja platformy, rozdz. 11.3) — wykonywane automatycznie w wybranym momencie cyklu. Hooki nie są mechanizmem walidacji blokującej — są punktem rozszerzalności: hook dokłada działanie (zapis do dziennika, powiadomienie, wywołanie automatyki), nie warunkuje kontynuacji wywołania.

| Punkt zaczepienia | Moment cyklu | Typowe zastosowanie |
|---|---|---|
| Przed startem sesji | Zanim pierwsze wywołanie modelu zostanie wysłane | Wczytanie dodatkowego kontekstu, powiadomienie o starcie procesu |
| Przed wywołaniem narzędzia | Zanim model wykona zgłoszone wywołanie narzędzia | Zapis do dziennika działań, wzbogacenie parametrów wywołania |
| Po wywołaniu narzędzia | Po otrzymaniu wyniku wywołania narzędzia, przed przekazaniem go modelowi | Przekształcenie lub wzbogacenie wyniku, zapis do Session Repository |
| Po odpowiedzi modelu | Po wygenerowaniu pełnej odpowiedzi, przed jej wyświetleniem | Powiadomienie, przekazanie wyniku do kolejki modułu Automations |
| Po zakończeniu sesji | Przy zamknięciu karty sesji lub zakończeniu roli | Podsumowanie, archiwizacja, wyzwolenie kolejnej automatyki w łańcuchu |

| Kanał | Dostępność hooków |
|---|---|
| Code CLI | Bezpośrednia — narzędzie udostępnia własny mechanizm punktów zaczepienia |
| Agent SDK | Bezpośrednia — biblioteka udostępnia własny mechanizm punktów zaczepienia |
| API | Pośrednia — te same punkty cyklu realizowane przez rdzeń platformy, który wywołuje wskazaną automatykę bezpośrednio z orkiestracji sesji |

Pięć punktów cyklu tworzy w interfejsie listę z akcją „+ dodaj punkt zaczepienia”. Każdy hook wskazuje automatykę modułu Automations — mechanizm nie dubluje silnika automatyk, lecz korzysta z niego wprost. Biblioteka gotowych hooków zasięgu globalnego zawiera zestawy startowe: „Zapis do repertorium”, „Powiadom po zakończeniu”, „Log wywołań narzędzi”, „Wzbogać wynik o znacznik czasu”. W jednym punkcie cyklu działa wiele hooków z ustalaną kolejnością i indywidualnym przełącznikiem włączenia (zasięgi: rola, agent) — wyłączenie hooka zachowuje jego definicję.

> **[?] Objaśnienie kontekstowe.** „Hook to dodatkowe działanie — najczęściej odwołanie do automatyki modułu Automations — wykonywane automatycznie w wybranym punkcie cyklu wywołania modelu. Brak podłączonych hooków oznacza cykl bez żadnego dodatkowego działania; podłączenie hooka nie wstrzymuje ani nie warunkuje przebiegu wywołania, wyłącznie dokłada do niego wybrane działanie.”

### 3.4. Narzędzia dozwolone i zabronione

Pole wskazuje, z jakich narzędzi — rozszerzeń zarejestrowanych zgodnie z dokumentem Rozszerzenia (wtyczki, umiejętności, konektory, serwery MCP) oraz wbudowanych akcji modułowych — model korzysta w toku danego wywołania. Kontrakt nakładki niesie dwie odrębne listy: `AllowedTools` (dozwolone) i `DisallowedTools` (zabronione), każda z własnym edytorem wielokrotnego wyboru zasilanym rejestrem rozszerzeń (Rozszerzenia, rozdz. 7).

| Kontekst wywołania | Źródło wartości pola | Edytowalność tutaj |
|---|---|---|
| Wywołanie związane z Agentem | Dziedziczy z Permissions Center agenta (Specyfikacja agentów, rozdz. 6.2) | Podgląd, tylko do odczytu, ze skrótem „Edytuj w Permissions Center →” |
| Wywołanie modelu bazowego bez pośrednictwa agenta | Brak nadrzędnej definicji do odziedziczenia | Edytowalne bezpośrednio tutaj — lista rozszerzeń z rejestru (Rozszerzenia, rozdz. 7), zaznaczanych do udostępnienia |

| Mechanizm | Działanie | Zasięgi |
|---|---|---|
| Grupy narzędzi | Nazwane zestawy uprawnień („tylko odczyt plików”, „sieć wyłączona”, „pełny developer”) zaznaczane jednym kliknięciem | rola, agent |
| Podgląd efektywnego zestawu narzędzi | Wynikowa lista po złożeniu list dozwolone−zabronione i dziedziczenia, widoczna przed zapisem | wszystkie warstwy |
| Stan wyjściowy | Puste listy oznaczają pełny dostęp operacyjny | wszystkie warstwy |

Stan wyjściowy — zgodnie z zasadą braku twardych blokad (Specyfikacja agentów, rozdz. 6.4) — to pełny dostęp operacyjny: model korzysta ze wszystkich podłączonych rozszerzeń i wbudowanych akcji modułu, w którym działa. Zawężenie jest zawsze świadomą decyzją Operatora, dodawaną do stanu bazowego.

> **[?] Objaśnienie kontekstowe.** „Narzędzia dozwolone i zabronione określają, z jakich rozszerzeń i akcji model korzysta w tym wywołaniu. Brak zawężenia oznacza pełny dostęp do wszystkich rozszerzeń podłączonych do agenta lub dostępnych w module. Zawężenie ogranicza wyłącznie to jedno przypisanie — nie zmienia uprawnień agenta w innych zastosowaniach.”

### 3.5. Model

Wybór konkretnego modelu bazowego udostępnianego przez wybrany kanał (rozdz. 3.7) — pole tożsame z „Nazwą modelu” encji Kanał modelu (Integracja modeli, rozdz. 5) oraz „Modelem bazowym” definicji agenta (Specyfikacja agentów, rozdz. 5.1). Lista dostępnych wartości pochodzi z hosta i zależy od wybranego kanału i trybu wywołania; obok listy zawsze dostępne jest pole wpisania nazwy modelu wprost, dzięki czemu niedostępność listy nie zatrzymuje pracy.

Pole „Model zapasowy” wskazuje model i kanał uruchamiany przy niedostępności podstawowego (zasięgi: rola, agent). Przełączenie na model zapasowy jest widoczne w Panelu prowenancji jako jawna wartość bloku „settings”.

> **[?] Objaśnienie kontekstowe.** „Model bazowy, do którego trafia skomponowana nakładka. Lista dostępnych modeli zależy od wybranego kanału — zmiana kanału zmienia zestaw modeli możliwych do wyboru. Nazwę modelu można również wpisać wprost.”

### 3.6. Effort i parametry próbkowania

Effort ustala poziom wysiłku obliczeniowego — głębokości rozumowania — jaki model angażuje przy generowaniu odpowiedzi, tam gdzie wybrany model i kanał to wspierają. Pole ma, jak każde ustawienie warstwowe platformy, dwa atrybuty: czy dana warstwa nadpisuje wartość odziedziczoną, oraz — jeśli tak — jaką wartość przyjmuje. Przy każdym poziomie prezentowany jest szacunkowy wpływ na czas odpowiedzi i koszt wywołania.

| Wartość | Charakterystyka |
|---|---|
| szybki | Najkrótszy czas odpowiedzi, najniższy koszt, rozumowanie ograniczone do niezbędnego minimum |
| standardowy | Wartość domyślna dla zadań typowych |
| dogłębny | Dłuższe rozumowanie przed odpowiedzią, właściwe zadaniom złożonym (analiza wielowątkowej sprawy) |
| maksymalny | Najgłębsze dostępne rozumowanie, najdłuższy czas odpowiedzi i najwyższy koszt |

```
effort:
  wlaczony: true | false     # czy ta warstwa nadpisuje wartość odziedziczoną
  wartosc:  szybki | standardowy | dogłębny | maksymalny
```

Parametry próbkowania — temperatura, top-p, kara za powtórzenia, sekwencje stop, `seed` — są polami warstwowymi o zasięgach sesja, rola, agent, z wartością wyjściową „wyłączony = wartość domyślna modelu”. Ustawienie `seed` daje powtarzalność wyniku, pozostałe parametry sterują zmiennością odpowiedzi bez modyfikacji kodu.

Zestaw pól (model + effort + parametry próbkowania + kanał) zapisuje się jako nazwany profil parametrów wywołania, dostępny na wszystkich warstwach i przywoływany jednym kliknięciem.

> **[?] Objaśnienie kontekstowe.** „Effort określa, ile rozumowania model angażuje przed wygenerowaniem odpowiedzi. Wyższy poziom wydłuża czas odpowiedzi i koszt wywołania, w zamian za głębszą analizę. Gdy pole jest wyłączone, wartość dziedziczona jest z warstwy szerszej lub z domyślnego zachowania wybranego modelu.”

### 3.7. Kanał i tryb wywołania

Cztery kanały integracji modeli — API, CLI, SSH, HTTP — są ogólnym mechanizmem platformy, ustalonym w dokumencie Integracja modeli (rozdz. 3) i zbudowanym z myślą o rozszerzalności: „dodanie kanału = dostarczenie nowego adaptera, bez zmian w rdzeniu” (Integracja modeli, rozdz. 4.2). Dla dostawców udostępniających pełną rodzinę form integracji — narzędzie wiersza poleceń, zestaw programistyczny (SDK) i punkt końcowy API jednocześnie — pole „Kanał” w oknie Nakładka i wywołanie modelu doprecyzowuje wybór do trzech trybów wywołania.

| Tryb wywołania | Adapter bazowy (Integracja modeli, rozdz. 3–4) | Sposób połączenia | Dane dostępowe |
|---|---|---|---|
| Code CLI | Adapter CLI | Wywołanie narzędzia wiersza poleceń dostawcy, lokalnie lub zdalnie | Token / konto narzędzia — odwołanie przechowywane poza bazą (rozdz. 3.8, 5) |
| Agent SDK | Adapter API — wariant programistyczny | Wywołanie biblioteki dostawcy osadzonej w rdzeniu platformy | Klucz dostępu — odwołanie przechowywane poza bazą |
| API | Adapter API — wywołanie bezpośrednie | Bezpośrednie żądanie HTTP do punktu końcowego dostawcy, bez pośrednictwa narzędzia ani biblioteki | Klucz dostępu — odwołanie przechowywane poza bazą |

Rozdział 3.2 zestawia harness właściwy każdemu trybowi; poniższa tabela dopełnia zestawienie o hooki (rozdz. 3.3) i rotację kont (rozdz. 3.8).

| Tryb wywołania | Harness | Hooki | Pula kont / rotacja |
|---|---|---|---|
| Code CLI | Wbudowany w narzędzie | Bezpośrednie, wbudowane w narzędzie | Dostępna — rozdz. 3.8 |
| Agent SDK | Dostarczany przez bibliotekę | Bezpośrednie, wbudowane w bibliotekę | Niedostępna — jedno konto na przypisanie |
| API | Budowany przez rdzeń platformy | Pośrednie, przez orkiestrację rdzenia | Niedostępna — jedno konto na przypisanie |

Zakres „Kanały modelu” udostępnia katalog wszystkich skonfigurowanych encji kanału modelu (typ, model bazowy, odwołanie do danych dostępowych, przypisania) z akcjami dodania, edycji, usunięcia i powielenia pozycji, na warstwach globalnej i środowiska. Edytor parametrów właściwych typowi kanału obejmuje zasięgi sesji, roli i agenta.

| Typ kanału | Parametry edytowalne |
|---|---|
| API | Punkt końcowy, nagłówki żądania, wersja interfejsu |
| CLI | Nazwa narzędzia, ścieżka, argumenty startowe |
| SSH | Host, konto, klucz lub hasło, model pomocniczy |
| HTTP | Adres interfejsu webowego |

Akcja „Testuj połączenie” wykonuje próbne wywołanie i przedstawia wynik: sukces, błąd uwierzytelnienia albo niedostępność modelu. Podgląd porównawczy uruchamia to samo polecenie przez dwa lub trzy kanały i modele jednocześnie, zestawiając odpowiedzi obok siebie w kolumnach (zasięg: sesja) — dobór modelu do zadania opiera się na realnym wyniku.

> **[?] Objaśnienie kontekstowe.** „Kanał ustala, w jaki sposób platforma technicznie dociera do modelu. Code CLI wywołuje narzędzie wiersza poleceń dostawcy i jako jedyny udostępnia pulę kont z rotacją; Agent SDK korzysta z biblioteki programistycznej osadzonej w rdzeniu; API wysyła żądanie bezpośrednio do punktu końcowego dostawcy. Zmiana kanału zmienia dostępny model, harness i sposób podłączenia hooków.”

### 3.8. Pula kont i rotacja — kanał Code CLI

Dla trybu Code CLI Operator konfiguruje pulę kilku odrębnie uwierzytelnionych kont lub miejsc narzędzia, między którymi platforma rozdziela wywołania sesji lub ról według wybranej strategii.

| Pole | Opis | Wartości |
|---|---|---|
| Pula kont | Lista kont Code CLI skonfigurowanych przez Operatora, każde z własnym odwołaniem do danych dostępowych, katalogiem profilu, stanem, ostatnim użyciem, bieżącym obciążeniem i limitem | lista nazwanych pozycji („Konto 1”, „Konto 2” …) |
| Strategia przydziału | Sposób wyboru konta z puli dla danego wywołania | kolejno (round-robin) \| najmniej obciążone \| ważona (udziały procentowe) \| wg limitu (do wyczerpania kwoty) \| lepka (to samo konto dla tej samej sprawy) \| ręcznie wskazane per sesja/rola |
| Stan konta | Bieżąca dostępność konta w puli | dostępne \| zajęte \| błąd uwierzytelnienia |
| Limit i kwota | Miękki licznik dzienny lub miesięczny wywołań albo tokenów, z ostrzeżeniem przy zbliżaniu się do progu | liczba \| brak limitu |
| Katalog profilu konta | Ścieżka lokalnego katalogu profilu konta (dane logowania, cache narzędzia), odrębna dla każdego konta | ścieżka |

Zastosowanie mechanizmu:

| Zastosowanie | Mechanizm |
|---|---|
| Równoległa przepustowość | Kilka ról środowiska MultitaskingAI lub kilka kart sesji pracuje jednocześnie, każda obsłużona przez inne konto z puli, bez wzajemnego oczekiwania |
| Rozdzielenie kont między sprawami klientów kancelarii | Różne konta przypisane do różnych spraw, spójne z zakresem izolacji technicznej „konto i token per sesja” (Architektura, rozdz. 12) |
| Przypięcie konta do profilu eksperta, sprawy lub roli | Jawne wiązanie konta z profilem eksperta (rozdz. 4), projektem albo rolą MultitaskingAI |
| Ochrona przed wyczerpaniem kwoty | Przekroczenie miękkiego limitu przełącza wywołania na kolejne konto puli, nie zatrzymuje pracy |

Wiersz konta udostępnia akcje: „Testuj połączenie”, „Popraw” (ponowne uwierzytelnienie przy błędzie), „Ustaw jako aktywne” i „Usuń” z możliwością cofnięcia. Podgląd obciążenia przedstawia licznik wywołań per konto w bieżącej sesji i w dobie wraz z plakietką stanu. Rejestr zdarzeń kont zapisuje przełączenia aktywnego konta, błędy uwierzytelnienia i przekroczenia limitów, ze znacznikiem czasu i urządzeniem — na tej podstawie Operator ustala, dlaczego dane wywołanie poszło z konkretnego konta.

Każde konto w puli działa w ramach własnego, odrębnego uwierzytelnienia u dostawcy narzędzia Code CLI i podlega warunkom korzystania właściwym temu kontu — pula rozdziela wywołania między osobno skonfigurowane konta, nie scala ich ani nie udostępnia jednego konta wielu równoległym wywołaniom jednocześnie. Przełączenie aktywnego konta jest jawną akcją albo jawnie skonfigurowaną regułą, nigdy ukrytą zmianą stanu. Jedno konto w pełni wystarcza do pracy; pula rozszerza przepustowość i rozdziela konteksty. Dane dostępowe każdego konta przechowywane są poza bazą danych (rozdz. 5) — encja kanału modelu przechowuje odwołanie do nich, a Operator ma zawsze możliwość podglądu pełnej wartości własnego klucza lub tokenu; maskowanie jest ustawieniem włączanym świadomie przez Operatora, nie zachowaniem wymuszonym.

> **[?] Objaśnienie kontekstowe.** „Pula kont pozwala skonfigurować kilka odrębnie uwierzytelnionych kont narzędzia Code CLI i rozdzielać między nie wywołania sesji lub ról — przy równoległej pracy kilku ról oraz przy oddzielaniu kont między sprawami klientów. Każde konto działa w ramach własnego uwierzytelnienia; pula nie łączy ani nie udostępnia jednego konta kilku wywołaniom naraz.”

### 3.9. Objaśnienia kontekstowe — tabela zbiorcza

Poniższa tabela zbiera w jednym miejscu pola opisane w rozdziale 3, z treścią objaśnienia `[?]`, zgodnie z wymogiem, że treść objaśnienia jest częścią definicji ustawienia (Architektura, rozdz. 13; Model konfiguracji, rozdz. 3.3).

| Pole | Stan domyślny | Treść `[?]` |
|---|---|---|
| System prompt | Wynik złożenia trzech warstw (rozdz. 2.3) | „Tekst przekazywany modelowi przed każdą wiadomością — wynik złożenia Konstytucji, Profilu/Roli i Ekspertyzy zadaniowej.” |
| Harness | Zależny od kanału | „Mechanizm spinający wywołanie modelu w pętlę narzędziową; pochodzenie zależy od wybranego kanału.” |
| Hooki | Brak podłączonych | „Dodatkowe działanie — zwykle automatyka — wykonywane automatycznie w wybranym punkcie cyklu wywołania; nie warunkuje przebiegu.” |
| Narzędzia dozwolone i zabronione | Pełny dostęp | „Zestaw rozszerzeń i akcji dostępnych modelowi w tym wywołaniu; brak zawężenia oznacza pełny dostęp.” |
| Model | Zależna od kanału | „Model bazowy przyjmujący nakładkę; lista zależy od wybranego kanału, nazwę można też wpisać wprost.” |
| Effort | wyłączony (dziedziczy) | „Poziom rozumowania angażowanego przed odpowiedzią; wyższy poziom wydłuża czas i koszt wywołania.” |
| Parametry próbkowania | wyłączone (wartość modelu) | „Temperatura, top-p, kara za powtórzenia, sekwencje stop i seed; wyłączenie oznacza zachowanie domyślne modelu.” |
| Kanał | Zależna od dostawcy | „Sposób technicznego dotarcia do modelu: Code CLI, Agent SDK albo API.” |
| Pula kont (Code CLI) | Jedno konto domyślne | „Kilka odrębnie uwierzytelnionych kont, między którymi rozdzielane są wywołania; każde konto działa w ramach własnego uwierzytelnienia.” |
| Konstytucja | Brak treści (sekcja pomijana) | „Trwałe zasady pracy AI obowiązujące każde wywołanie na platformie albo w tym środowisku; brak treści oznacza pominięcie tej warstwy.” |
| Negatyw | Brak treści | „Sekcja »czego nie robić« dokładana do nakładki, oddzielona od zasad pozytywnych.” |
| Styl i ton | Dziedziczony | „Sposób formułowania wypowiedzi, niezależny od treści merytorycznej nakładki.” |
| Budżet kontekstu i tokenów | Brak limitu | „Miękki budżet długości kontekstu i odpowiedzi; przekroczenie skutkuje ostrzeżeniem w Panelu prowenancji.” |

---

## 4. Profile ekspertów

### 4.1. Profil eksperta jako byt konfiguracji

Profil eksperta jest nazwanym, wersjonowanym bytem konfiguracji nakładki — obok profili izolacji, zespołu, agenta i asystenta (Model konfiguracji, rozdz. 7.1). Profil eksperta zawiera spójny zestaw warstw nakładki i pól wywołania odpowiadający jednej osobowości roboczej, przełączany jednym gestem: Konstytucja lokalna, Profil/Rola, Ekspertyza zadaniowa, styl, negatyw, narzędzia, model, effort, kanał i hooki. Jeden przełącznik zmienia całą tożsamość Wykonawcy pracującego z Operatorem — „dziś Roman (redaktor), jutro Ignacy (analityk prawny)”.

| Element | Działanie | Efekt |
|---|---|---|
| Szybki przełącznik eksperta | Lista rozwijana aktywnych ekspertów w pasku kontekstu okna oraz w uproszczonym menu okna operacyjnego (warstwa sesji) | Zmiana eksperta w trakcie pracy, bez wchodzenia w pełną konfigurację |
| Ekspert per zasięg | Ekspert domyślny globalny, odrębny per środowisko i projekt, doraźny per sesja, przypisany per rola MultitaskingAI | Spójność z warstwowością: „Ignacy globalnie, Roman w projekcie Redakcja” |
| Galeria ekspertów | Widok kart wszystkich profili z podglądem, stanem (aktywny / nieaktywny), ostatnim użyciem i akcją „Aktywuj” | Zarządzanie zestawem osobowości roboczych w jednym miejscu |
| Harmonogram eksperta | Powiązanie profilu z regułą czasową modułu Automations („w dni robocze 8–16 aktywny Ignacy”) | Zmiana osobowości roboczej według kalendarza, z możliwością ręcznego nadpisania |

Brak aktywnego profilu eksperta oznacza pracę na warstwach nakładki bez kapsuły. Aktywacja eksperta nakłada się na warstwy, nie kasuje ich; Operator nadpisuje dowolne pole eksperta na warstwie węższej w każdej chwili.

### 4.2. Wersjonowanie profilu eksperta

| Element | Działanie | Efekt |
|---|---|---|
| Wersja profilu | Każdy zapis tworzy wersję: numer, data, urządzenie autora, notatka zmiany; pełna historia wersji profilu | Odtwarzalność stanu profilu z dowolnego momentu |
| Przywrócenie i rozgałęzienie | Operacje „Przywróć wersję N” oraz „Utwórz wariant z wersji N” („Roman — wariant sądowy”) | Iteracyjne dopracowywanie eksperta bez utraty sprawdzonych wersji |
| Porównanie wersji | Zestawienie dwóch wersji pole po polu: zmiany w Konstytucji, zestawie narzędzi, modelu | Pełny obraz zmiany przed aktywacją |
| Hash wersji eksperta | Skrót SHA-256 całej kapsuły wersji, widoczny w Panelu prowenancji obok hasha Konstytucji (rozdz. 10.2) | Weryfikowalne wskazanie, która wersja profilu obsłużyła daną wiadomość |

---

## 5. Dane dostępowe i zmienne środowiskowe

Wszystkie dane dostępowe — klucze API, tokeny narzędzi wiersza poleceń, dane SSH i HTTP — oraz zmienne środowiskowe wstrzykiwane do procesu sesji mają jedno, jawne miejsce sterowania, spójne z zasadą „dane poza bazą, odwołanie w bazie” (Integracja modeli, rozdz. 11; Bezpieczeństwo, rozdz. 10) oraz z zasadą jawności kluczy.

| Element | Działanie | Efekt |
|---|---|---|
| Rejestr kluczy i danych dostępowych | Widok wszystkich odwołań: rodzaj (klucz API / token CLI / dane SSH / dane HTTP), właściciel (kanał, konto), stan, ostatnie użycie; wartość domyślnie jawna, z przełącznikiem „Pokaż / Ukryj” | Operator ma zawsze wgląd we własne klucze; maskowanie jest ustawieniem, nie przymusem |
| Edytor zmiennych środowiskowych procesu sesji | Tabela par `KLUCZ = wartość` wstrzykiwanych do środowiska procesu wykonawczego, z importem z pliku `.env` i eksportem do pliku `.env` | Standardowy format; przenoszenie konfiguracji między środowiskami |
| Zasięg zmiennych środowiskowych | Każda zmienna ma warstwę (globalna → środowisko → sesja → rola) i podlega dziedziczeniu jak każde ustawienie | Odrębne klucze dla różnych spraw i ról bez duplikowania całości konfiguracji |
| Wskazanie wartości przez odwołanie albo wprost | Pole wskazuje wpis magazynu danych dostępowych albo niesie wartość jawną wpisaną wprost | Wybór między odwołaniem a wartością roboczą należy do Operatora |
| Podgląd miejsc użycia klucza | Przy kluczu lista miejsc użycia: kanały, konta i role, które go odczytują | Świadomość zależności przed zmianą lub usunięciem klucza |
| Maskowanie danych dostępowych | Globalny przełącznik „maskuj dane dostępowe w podglądach” oraz lokalny „Pokaż / Ukryj” per pole; stan domyślny — wartości jawne | Kontrola widoczności przy pracy na współdzielonym ekranie, bez odbierania dostępu |
| Dane dostępowe w Panelu prowenancji | W bloku `argv` klucz i token są domyślnie jawne; maskowanie obejmuje wyłącznie znaki środkowe | Weryfikowalność wysłanego tokenu z możliwością ukrycia przy zrzutach ekranu |

Same wartości pozostają przechowywane poza bazą danych; baza trzyma odwołanie, przez co eksport bazy nie ujawnia kluczy (Integracja modeli, rozdz. 11). Wartości jawne wpisane wprost w zmiennych środowiskowych procesu są traktowane jak dane dostępowe i również trafiają do magazynu poza bazą.

---

## 6. Zależności i izolacja jako ustawienia konfiguracyjne

Zgodnie z zasadą nadrzędną platformy (Koncepcja platformy, rozdz. 6) żadna zależność ani żaden zakres izolacji opisany w niniejszym rozdziale nie jest wymogiem — każdy jest ustawieniem konfiguracyjnym, które Operator ustanawia świadomie. Okno Nakładka i wywołanie modelu nie duplikuje mechanizmów izolacji, integracji i rozszerzeń opisanych w dokumentach źródłowych — poniższe podrozdziały wskazują punkty styku.

### 6.1. Relacja do okna konfiguracji punktów izolacji

Trzy z ośmiu zakresów izolacji technicznej (Architektura, rozdz. 12) dotyczą bezpośrednio warstwy opisanej w niniejszym dokumencie (Integracja modeli, rozdz. 9).

| Zakres izolacji technicznej | Relacja do nakładki i wywołania modelu | Stan domyślny |
|---|---|---|
| Katalog danych i konfiguracji modelu | Czy dane i konfiguracja kanału modelu przypisanego sesji lub roli są przechowywane odrębnie, czy współdzielone | Wyłączony (współdzielony) |
| Konto i token per sesja | Czy dane dostępowe kanału (w tym poszczególne konta puli Code CLI, rozdz. 3.8) są odrębne dla każdej sesji, czy współdzielone | Wyłączony (współdzielony) |
| Model procesu | Czy proces obsługujący wywołania modelu jest odrębny, czy współdzielony z innymi | Wyłączony (współdzielony) |

Skrót „Otwórz okno konfiguracji punktów izolacji →” dostępny z kolumny prawej (rozdz. 12.2) prowadzi bezpośrednio do trzykolumnowego okna izolacji (Model konfiguracji, rozdz. 6; Izolacja i konfigurowalność zależności, rozdz. 10), z wstępnie zaznaczonym poziomem zasięgu odpowiadającym warstwie aktualnie wybranej w oknie Nakładka i wywołanie modelu. Żadne pole opisane w rozdziale 3 nie wymusza włączenia jakiegokolwiek zakresu izolacji technicznej — pula kont (3.8) działa niezależnie od tego, czy zakres „konto i token per sesja” jest włączony, choć jego włączenie dodatkowo separuje techniczny proces obsługujący każde konto.

### 6.2. Relacja do zakresu Integracje

Zakres Integracje (5.8 Modelu konfiguracji) obejmuje jawne, konfigurowalne powiązania między elementami platformy — w tym powiązanie komponentu własnego (agenta) z modułem, w którym ma zastosowanie. Nakładka korzysta z tego mechanizmu wprost: przypisanie agenta do roli (sekcja Role) lub do projektu (Agent Manager) jest tą samą, jawną decyzją integracyjną, opisaną w rozdziale 5.8 — okno Nakładka i wywołanie modelu jej nie powtarza, lecz odsyła do niej skrótem widocznym przy warstwie Profil/Rola (rozdz. 2.2).

### 6.3. Relacja do Rozszerzeń i Permissions Center

Pole „Narzędzia dozwolone i zabronione” (rozdz. 3.4) czerpie treść z rejestru rozszerzeń (Rozszerzenia, rozdz. 7) oraz — dla wywołań związanych z agentem — z Permissions Center (Specyfikacja agentów, rozdz. 6). Instalacja i włączanie rozszerzeń pozostaje wyłączną odpowiedzialnością rejestru rozszerzeń; okno Nakładka i wywołanie modelu odczytuje jego zawartość i, dla wywołań bez agenta, pozwala zaznaczyć podzbiór do udostępnienia.

### 6.4. Stan wyjściowy

| Element | Stan wyjściowy | Charakter |
|---|---|---|
| Konstytucja | Bez treści | Model działa bez tej warstwy — stan w pełni funkcjonalny |
| Narzędzia dozwolone i zabronione | Pełny dostęp | Zgodnie z zasadą braku twardych blokad (Specyfikacja agentów, rozdz. 6.4) |
| Hooki | Brak podłączonych | Cykl wywołania bez dodatkowego działania |
| Pula kont Code CLI | Jedno, domyślnie skonfigurowane konto | Rotacja jest ustawieniem dodawanym świadomie |
| Profil eksperta | Brak aktywnego | Praca na warstwach nakładki bez kapsuły profilu |
| Zakresy izolacji technicznej powiązane z modelem (6.1) | Wszystkie trzy wyłączone | Pełna operacyjna swoboda procesu, zgodnie z rozdz. 6.6 Koncepcji platformy |

---

## 7. Uwierzytelnianie jako ustawienie konfiguracyjne

Okno Konfiguracji steruje wymogiem logowania i metodami dostępu — egzekwowanie logowania jest ustawieniem Operatora, nie twardą bramką (Bezpieczeństwo, rozdz. 8 i 13). Tożsamość Operatora — konto, metody, urządzenia — pozostaje przedmiotem okna Ustawień; tutaj konfigurowane jest wyłącznie zachowanie platformy wobec dostępu.

| Element | Działanie | Efekt |
|---|---|---|
| Przełącznik wymogu logowania | Pozycja zakresów „Aplikacja” i „Procesy”: wymóg logowania włączony lub wyłączony, z objaśnieniem `[?]` | Operator włącza kontrolę dostępu świadomie, w dowolnej chwili |
| Podgląd stanu uwierzytelniania | Blok informacyjny: czy okno logowania wymusza dostęp, które metody są aktywne, ilu urządzeń dotyczy — ze skrótem „Zarządzaj w Ustawieniach →” | Świadomość stanu bez opuszczania kontekstu konfiguracji |
| Polityka tokenu urządzenia | Konfigurowalny czas ważności i rotacji tokenu urządzenia; wartość domyślna — ważny do unieważnienia (Bezpieczeństwo, rozdz. 8.1) | Wzmocnienie kontroli dostępu jako ustawienie warstwowe |
| Powiązanie z izolacją „konto i token per sesja” | Skrót do okna izolacji z zaznaczonym zakresem „konto i token per sesja” dla roli lub sesji | Jawna, konfigurowalna zależność między uwierzytelnianiem a izolacją (Bezpieczeństwo, rozdz. 8.4) |
| Ostrzeżenie zamiast blokady | Wyłączenie ostatniej metody logowania albo samego wymogu wywołuje komunikat ostrzegawczy, nigdy uniemożliwienie operacji | Konsekwentny brak twardych blokad także w obszarze kontroli dostępu |

---

## 8. Konfiguracja kanałów komunikacji: Chat Window i Execution Loop Window

Platforma prowadzi dwa kanały komunikacji operacyjnej, oba pierwszoplanowe w architekturze. Okno Konfiguracji steruje zachowaniem obu, na wszystkich warstwach zasięgu.

| Kanał | Role | Miejsce w układzie |
|---|---|---|
| Chat Window | Użytkownik ↔ Wykonawca | Lewa kolumna obszaru roboczego, stała, pełna wysokość |
| Execution Loop Window | Koordynator ↔ Wykonawca | Kolumna sąsiadująca z Chat Window, otwierana |

### 8.1. Konfiguracja Chat Window

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą i podstawowym mechanizmem sterowania procesami platformy: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań i wyjaśnianie wyniku. Każdy moduł i każde środowisko udostępnia je w tym samym miejscu układu.

| Grupa ustawień | Ustawienie | Wartości | Warstwy |
|---|---|---|---|
| Wykonawca i model | Wykonawca obsługujący kanał, profil eksperta, model, effort, kanał wywołania | wybór z katalogu (rozdz. 3, 4) | globalna, środowisko, projekt, sesja, rola |
| Strumień odpowiedzi | Prezentacja odpowiedzi: strumieniowa albo po zakończeniu; widoczność kroków pośrednich i wywołań narzędzi | strumieniowo \| po zakończeniu; kroki widoczne \| zwinięte | globalna, sesja, rola |
| Zatwierdzanie i przerywanie | Sterowanie przebiegiem z okna komunikacji: zatwierdzenie działania, przerwanie, korekta polecenia | włączone \| wyłączone dla wskazanych klas działań | globalna, środowisko, sesja, rola |
| Punkty zatwierdzeń | Klasy działań wymagające potwierdzenia Użytkownika przed wykonaniem (zapis do repozytorium, wysyłka, operacja nieodwracalna) | lista klas działań | globalna, środowisko, projekt, rola |
| Wyjaśnianie wyniku | Dostępność objaśnienia kontekstu i uzasadnienia wyniku bezpośrednio w strumieniu | włączone \| wyłączone | globalna, sesja |
| Szerokość kolumny | Szerokość lewej kolumny obszaru roboczego | wartość w jednostkach siatki | globalna, środowisko, sesja |
| Historia i pamięć kanału | Zakres historii wczytywanej do kontekstu, zasób pamięci kontekstowej przypisany kanałowi | odwołanie do zakresów 5.9 i 5.10 | wszystkie warstwy |

### 8.2. Konfiguracja Execution Loop Window

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą: bieżące zlecenie i jego dekompozycję na zadania, kolejkę i stan zadań, wymianę komunikatów sterujących, wyniki kontroli jakości i decyzje o ponowieniu, wskaźniki przebiegu pętli oraz sterowanie przebiegiem. Okno otwiera się jako kolumna sąsiadująca z Chat Window.

| Grupa ustawień | Ustawienie | Wartości | Warstwy |
|---|---|---|---|
| Pętla wykonawcza | Limit kroków pętli, limit czasu pojedynczego kroku, zachowanie przy błędzie zadania | liczba kroków; limit czasu; kontynuuj \| wstrzymaj \| przekaż do decyzji Użytkownika | środowisko, projekt, sesja, rola |
| Równoległość zadań | Liczba zadań realizowanych jednocześnie, liczba Wykonawców w zespole, przydział kont z puli (rozdz. 3.8) | liczba \| bez ograniczenia | środowisko, projekt, sesja, rola |
| Polityka ponowień | Liczba ponowień zadania, odstęp między ponowieniami, warunki ponowienia (błąd narzędzia, wynik poniżej progu jakości), przełączenie na model zapasowy | liczba; odstęp; lista warunków | środowisko, projekt, rola, agent |
| Progi kontroli jakości | Próg akceptacji wyniku zadania, sposób oceny (kontrola przez Koordynatora, kontrola porównawcza, sprawdzenie testem), zachowanie przy wyniku poniżej progu | wartość progu; metoda; ponów \| przekaż Użytkownikowi \| zaakceptuj z adnotacją | środowisko, projekt, rola, agent |
| Zakres autonomii | Zakres działań realizowanych przez Koordynatora bez potwierdzenia Użytkownika: dekompozycja zlecenia, przydział zadań, ponowienie, zmiana modelu, uruchomienie narzędzia zewnętrznego | lista działań objętych autonomią | globalna, środowisko, projekt, rola |
| Punkty zatwierdzeń | Momenty pętli wymagające zatwierdzenia Użytkownika: przyjęcie planu dekompozycji, zakończenie etapu kontroli jakości, wynik końcowy | lista punktów | globalna, środowisko, projekt, rola |
| Sterowanie przebiegiem | Dostępność operacji: wstrzymanie, wznowienie, przerwanie, korekta zlecenia | włączone \| wyłączone per operacja | globalna, sesja, rola |
| Prezentacja pętli | Widoczność kolejki zadań, komunikatów sterujących, wskaźników przebiegu i wyników kontroli jakości | pełna \| zwinięta do wskaźników | globalna, sesja |
| Szerokość kolumny | Szerokość kolumny okna pętli wykonawczej | wartość w jednostkach siatki | globalna, środowisko, sesja |

### 8.3. Relacja obu kanałów do nakładki

Oba kanały korzystają z tej samej nakładki i tych samych pól wywołania (rozdz. 2–3), rozstrzyganych według warstwy zasięgu właściwej kanałowi. Wywołania modelu inicjowane z Chat Window oraz wywołania inicjowane przez Koordynatora w pętli wykonawczej trafiają do Panelu prowenancji (rozdz. 10) jako odrębne wpisy, rozróżniane polem kanału komunikacji — Operator ustala, czy dany prompt systemowy pochodzi z polecenia Użytkownika, czy z zadania przydzielonego przez Koordynatora.

---

## 9. Konfiguracja warstw widoczności

Interfejs Danaco Console ujawnia funkcje stopniowo: funkcja niepotrzebna do realizacji aktualnego zadania nie jest widoczna. Okno Konfiguracji steruje przypisaniem funkcji do warstw widoczności, konfiguracją warstw według roli użytkownika oraz mechanizmami dostępu do funkcji ukrytych.

### 9.1. Cztery warstwy widoczności

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór Wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

### 9.2. Przypisanie funkcji do warstw

| Ustawienie | Działanie | Warstwy zasięgu |
|---|---|---|
| Macierz przypisania funkcji | Tabela wszystkich funkcji interfejsu z kolumną warstwy widoczności (1–4) i kolumną sposobu wywołania; zmiana warstwy przenosi funkcję między poziomami ujawniania | globalna, środowisko, rola |
| Przywrócenie przypisania domyślnego | Akcja przywracająca warstwę wbudowaną dla wskazanej funkcji albo dla całej grupy funkcji | globalna, środowisko, rola |
| Reguła jednego kliknięcia | Kontrola spójności: każda funkcja warstw 2–4 pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego; przypisanie naruszające regułę wywołuje ostrzeżenie | globalna |
| Mechanizm ukrycia | Sposób ukrycia funkcji: menu progresywne, panel wysuwany, grupowanie logiczne akcji, znacznik kontekstowy | globalna, środowisko, rola |

### 9.3. Konfiguracja warstw według roli użytkownika

| Ustawienie | Działanie |
|---|---|
| Profil warstw dla roli | Nazwany zestaw przypisań warstw właściwy roli użytkownika (użytkownik podstawowy, użytkownik zaawansowany, administrator); role dziedziczą przypisania z warstwy globalnej i nadpisują wybrane pozycje |
| Zakres widoczności roli | Rozstrzygnięcie, które warstwy są dostępne roli; użytkownik podstawowy nie widzi elementów warstwy 4 |
| Przypisanie roli do użytkownika i urządzenia | Powiązanie profilu warstw z kontem użytkownika oraz z urządzeniem, spójnie z zakresem Aplikacja (5.1) |
| Podgląd interfejsu w wybranej roli | Podgląd stanu spoczynku interfejsu widzianego przez wybraną rolę, bez zmiany roli bieżącej |

### 9.4. Wyszukiwarka funkcji

| Ustawienie | Działanie |
|---|---|
| Zakres wyszukiwania | Wyszukiwarka obejmuje wszystkie funkcje i wszystkie ustawienia wszystkich zakresów okna Konfiguracji, wraz z treścią objaśnień `[?]` |
| Wywołanie | Skrót klawiszowy oraz pole „szukaj ustawienia” w pasku kontekstu okna |
| Wynik | Lista pozycji z nazwą funkcji, warstwą widoczności, sposobem wywołania i przejściem wprost do miejsca ustawienia |
| Funkcje warstwy 4 | Wyszukiwarka udostępnia funkcje warstwy 4 rolom, którym ta warstwa jest dostępna |

### 9.5. Skróty klawiszowe

| Ustawienie | Działanie |
|---|---|
| Katalog skrótów | Lista wszystkich skrótów z przypisaną funkcją, warstwą widoczności i możliwością zmiany kombinacji klawiszy |
| Skróty okna Konfiguracji | Przełączenie profilu eksperta, podgląd na żywo Panelu prowenancji, zapis Konstytucji, przełączenie warstwy zasięgu, otwarcie Panelu prowenancji, otwarcie wyszukiwarki funkcji |
| Skróty kanałów komunikacji | Przeniesienie fokusu do Chat Window, otwarcie i zamknięcie Execution Loop Window, wstrzymanie i wznowienie pętli wykonawczej |
| Konflikty | Przypisanie kombinacji już zajętej wywołuje ostrzeżenie ze wskazaniem funkcji zajmującej skrót; zapis pozostaje możliwy |

### 9.6. Tryb administracyjny

Niniejsza sekcja jest miejscem rozstrzygającym dla trybu administracyjnego całej platformy: ustala jedną drogę wejścia, jeden model uprawnień, postać znacznika stanu oraz przejścia stanu. `architektura/model-konfiguracji.md` (rozdz. 7.4) podaje klucze konfiguracji trybu, `architektura/bezpieczenstwo-i-uwierzytelnianie.md` (rozdz. 14.2) — model uprawnień i rejestrację w dzienniku audytu.

**Droga wejścia.** Tryb administracyjny wywołuje się wyłącznie z palety poleceń (`Ctrl/Cmd + K`, `interfejs-uzytkownika/katalog-komponentow.md`, rozdz. 15.8), pozycją „Tryb administracyjny”. Okno konfiguracji nie zawiera przełącznika trybu — prezentuje znacznik stanu trybu oraz rejestr wejść. Pozycja palety jest widoczna wyłącznie dla roli, której konfiguracja obejmuje warstwę 4.

**Model uprawnień.**

| Element modelu | Treść |
|---|---|
| Warunek widoczności pozycji | Rola użytkownika obejmująca warstwę 4 (`visibility.layer.role`) |
| Warunek wejścia | Potwierdzenie uprawnień — uwierzytelnienie właściciela konta metodą ustawioną w oknie Ustawień (`interfejs-uzytkownika/ustawienia.md`, rozdz. 3) |
| Zakres ujawnianych funkcji | Wartość `visibility.admin.scope`: narzędzia diagnostyczne niskiego poziomu, konfiguracja ról, operacje na procesach sesji |
| Zasięg trybu | Sesja, urządzenie albo rola użytkownika — wartość `visibility.admin.mode` |
| Czas trwania | Tryb wygasa po okresie bezczynności właściwym poświadczeniu urządzenia oraz przy rozłączeniu urządzenia |
| Rejestracja | Wejście, wyjście i wygaśnięcie trybu oraz każda zmiana konfiguracji uprawnień odnotowywane są w dzienniku audytu (`architektura/bezpieczenstwo-i-uwierzytelnianie.md`, rozdz. 13.5) |

**Makieta znacznika stanu trybu.** Znacznik zajmuje skrajnie prawą pozycję paska kontekstu i jest widoczny przez cały czas działania trybu; kliknięcie znacznika otwiera rozwinięcie z pozostałym czasem trwania i akcją wyjścia.

```
Makieta — znacznik stanu trybu administracyjnego w pasku kontekstu

 ═══════════════════════════════════════════════════════════════════════════════
  Pasek kontekstu:  [Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]   [⛨ Tryb administracyjny ▼]
 ═══════════════════════════════════════════════════════════════════════════════

  Rozwinięcie znacznika (panel popover, warstwa 3):

   ┌──────────────────────────────────────────┐
   │ ⛨ Tryb administracyjny — aktywny         │
   │ Zasięg:            sesja                 │
   │ Zakres funkcji:    narzędzia diagnost.   │
   │ Wygaśnięcie za:    24 min                │
   │ ────────────────────────────────────────  │
   │ [ Wyjdź z trybu ]   [ Rejestr wejść ]     │
   └──────────────────────────────────────────┘

  Poza trybem administracyjnym znacznik nie występuje w pasku kontekstu
  i nie pozostawia śladu w układzie.
```

| Element znacznika | Warstwa | Sposób wywołania |
|---|---|---|
| Znacznik stanu trybu (`.dn-plakietka--stan` z ikoną `⛨`) | 1 w czasie działania trybu; nieobecny poza trybem | Widoczny bez interakcji od chwili wejścia w tryb |
| Rozwinięcie znacznika | 3 | Kliknięcie znacznika |
| Akcja „Wyjdź z trybu” | 3 | Pozycja rozwinięcia znacznika |
| Rejestr wejść | 4 | Pozycja rozwinięcia znacznika albo paleta poleceń `Ctrl/Cmd + K` |

**Diagram przejść stanu.**

```
 STAN POZA TRYBEM (spoczynek interfejsu)
   widoczna warstwa 1 i zwinięte wyzwalacze warstw 2–3; znacznik trybu nieobecny
        │
        │  paleta poleceń `Ctrl/Cmd + K` → pozycja „Tryb administracyjny”
        ▼
 POTWIERDZENIE UPRAWNIEŃ
   uwierzytelnienie właściciela konta; rola obejmująca warstwę 4
        │                                   │
        │  potwierdzenie przyjęte            │  potwierdzenie odrzucone
        ▼                                   └──────────────► STAN POZA TRYBEM
 PRACA W TRYBIE ADMINISTRACYJNYM                              (komunikat blokowy
   funkcje warstwy 4 ujawnione w zakresie                       przy polu, bez blokady
   `visibility.admin.scope`; znacznik stanu                     dostępu do palety)
   widoczny w pasku kontekstu; wpis w dzienniku audytu
        │                        │
        │  akcja „Wyjdź          │  bezczynność do progu wygaśnięcia
        │   z trybu”             │  albo rozłączenie urządzenia
        ▼                        ▼
 WYJŚCIE Z TRYBU            WYGAŚNIĘCIE TRYBU
   wpis w dzienniku audytu    wpis w dzienniku audytu; komunikat o wygaśnięciu
        │                        │
        └────────────┬───────────┘
                     ▼
             STAN POZA TRYBEM
```

Wyjście z trybu i jego wygaśnięcie przywracają stan spoczynku interfejsu: funkcje warstwy 4 przestają być widoczne, znacznik znika z paska kontekstu, a praca prowadzona w oknach operacyjnych trwa bez przerwania.

### 9.7. Ergonomia i personalizacja okna

| Ustawienie | Działanie | Efekt |
|---|---|---|
| Widok prosty i widok pełny | Widok prosty prezentuje najczęściej zmieniane pola (ekspert, model, effort); widok pełny — wszystkie pola zakresu | Niższy próg wejścia bez ograniczania dostępu do funkcji |
| Pola przypięte | Operator przypina najczęściej zmieniane pola na początek kolumny środkowej | Personalizacja układu pod własny sposób pracy |
| Objaśnienia kontekstowe `[?]` | Każde ustawienie okna niesie treść objaśnienia opisującą działanie i wpływ (Model konfiguracji, rozdz. 3.3) | Objaśnienie jest częścią definicji ustawienia |
| Wskaźnik zmian niezapisanych | Pola selektorowe zapisują się natychmiast; pola tekstowe prezentują stan „niezapisane”, bez blokowania nawigacji | Przewidywalność zapisu przy zachowaniu swobody nawigacji |
| Podgląd wpływu zmiany | Dla zmiany o szerokim zasięgu (Konstytucja globalna) — informacja „dotyczy N środowisk / ról / sesji” przed zapisem | Świadomość konsekwencji zmiany na szerokiej warstwie |

---

## 10. Panel prowenancji — „co poszło do modelu”

### 10.1. Cel i pozycjonowanie

Panel prowenancji daje Operatorowi dokładny, weryfikowalny podgląd tego, co faktycznie zostało — albo zostałoby — przekazane modelowi przy danym wywołaniu, po uwzględnieniu całego dziedziczenia warstw opisanego w rozdziałach 2 i 3. Jest odpowiednikiem „podglądu polityki efektywnej” (Model konfiguracji, rozdz. 6.4; Izolacja i konfigurowalność zależności, rozdz. 10.2) — ten sam wzorzec przejrzystości, zastosowany nie do polityki izolacji, lecz do faktycznej treści wywołania modelu.

| Mechanizm pokrewny | Domena | Panel prowenancji |
|---|---|---|
| Polityka efektywna (okno izolacji) | Wynikowy zestaw reguł izolacji po dziedziczeniu | Wynikowa treść i ustawienia wywołania modelu po dziedziczeniu |
| Miejsce w oknie | Kolumna prawa okna izolacji | Kolumna prawa okna Nakładka i wywołanie modelu |
| Forma | `.dn-kod` — blok tekstowy formatowany | `.dn-kod` — blok tekstowy formatowany, pięć sekcji |

Panel jest miejscem wglądu w treść budowaną dynamicznie: Ekspertyza zadaniowa komponowana przez Koordynatora nie ma stałego miejsca edycji, a Panel prowenancji pokazuje dokładny tekst dołożony do nakładki dla danego kroku procesu.

### 10.2. Bloki zawartości

| Blok | Zawartość | Zależność od trybu wywołania (rozdz. 3.7) |
|---|---|---|
| argv | Dosłowny zapis parametrów technicznie przekazanych przy wywołaniu | Code CLI: lista argumentów wiersza poleceń · Agent SDK: parametry wywołania biblioteki · API: nagłówki i parametry żądania HTTP |
| system prompt | Pełny, skomponowany tekst nakładki (rozdz. 2.3), z wyraźnie oznaczonymi granicami trzech warstw, wyróżnieniem sekcji zbudowanej przez Koordynatora i wskazaniem poziomu zasięgu każdej warstwy | Niezależne od trybu |
| settings | Wynikowe wartości pól z rozdziału 3 (harness, hooki, narzędzia dozwolone i zabronione, model, effort, parametry próbkowania, kanał, pula kont) po uwzględnieniu dziedziczenia warstw | Niezależne od trybu — treść pól zależna od wybranego trybu |
| warstwy i pochodzenie | Dla każdego pola nakładki i wywołania: warstwa zwycięska oraz wartość dziedziczona, z której pole zostało nadpisane | Niezależne od trybu |
| hash konstytucji i wersji eksperta | Skrót kryptograficzny (SHA-256) treści Konstytucji oraz kapsuły wersji profilu eksperta obowiązujących w chwili wywołania, w postaci skróconej z akcją skopiowania pełnej wartości | Niezależne od trybu; pole puste, gdy Konstytucja lub profil eksperta bez treści |

Wartości pochodzące z odwołania do danych dostępowych (rozdz. 3.7–3.8, 5) — token narzędzia, klucz dostępu, dane logowania — wyświetlane są w bloku argv domyślnie w postaci jawnej i pełnej, zgodnie z zasadą, że Operator ma zawsze wgląd we własne dane dostępowe; maskowanie jest ustawieniem włączanym świadomie.

Format bloku „settings”:

```
SETTINGS — wywołanie: rola Executor 1, zespół „Sprawa X”      (podgląd w .dn-kod)
  harness           : wbudowany w narzędzie Code CLI          (poziom: rola)
  hooki             : po odpowiedzi modelu → automatyka „Zapis do repertorium”
                                                                (poziom: agent)
  narzędzia         : pełny dostęp                            (dziedziczone: Permissions Center)
  model             : <nazwa modelu bazowego>                  (poziom: agent)
  effort            : dogłębny                                 (poziom: rola)
  kanał             : Code CLI                                 (poziom: rola)
  pula kont         : „Konto 2” (Sprawa X)                      (poziom: rola)
  kanał komunikacji : Execution Loop Window (Koordynator ↔ Wykonawca)
```

### 10.3. Tryby podglądu i narzędzia analizy

| Tryb | Kiedy dostępny | Zawartość |
|---|---|---|
| Podgląd na żywo | Przed wysłaniem — z okna operacyjnego (uproszczone menu kontekstowe) i z okna Konfiguracji | Co zostałoby wysłane przy najbliższym wywołaniu, bez faktycznego wywołania modelu |
| Zapis archiwalny | Po każdym faktycznym wywołaniu | Komplet bloków zapisany razem z wpisem historii (5.9), dostępny retrospektywnie dla danej wiadomości lub tury rozmowy |

Zapis archiwalny jest ustawieniem włączonym domyślnie i konfigurowalnym niezależnie dla każdej warstwy zasięgu (rozdz. 2.4) — wyłączenie go dla wybranej sesji lub roli nie usuwa już zapisanych wpisów, wyłącznie wstrzymuje tworzenie nowych.

| Narzędzie | Działanie | Efekt |
|---|---|---|
| Porównanie dwóch wywołań | Zestawienie prowenancji dwóch wpisów archiwalnych pole po polu, z podświetleniem różnic | Diagnoza, co zmieniło się między wywołaniami o różnym wyniku |
| Oś czasu prowenancji sesji | Chronologiczny wykaz wywołań w sesji lub roli, z hashami i znacznikami czasu; kliknięcie ładuje pełny zapis | Nawigacja po historii wywołań |
| Filtr i wyszukiwanie | Filtrowanie osi czasu po kanale wywołania, kanale komunikacji, modelu, koncie, profilu eksperta, hashu Konstytucji i obecności hooka | Odnalezienie konkretnego wywołania |
| Weryfikacja hasha | Pole „wklej hash → sprawdź wersję” zestawiające skrót z historią wersji Konstytucji i profilu eksperta (rozdz. 4.2) | Potwierdzenie, że sesja korzystała z zamierzonej wersji |
| Eksport prowenancji | Zapis kompletu bloków do pliku Markdown lub JSON | Przenośna kopia zapisu „co dokładnie wysłano” |
| Przekazanie do automatyki | Akcja „Załącz do automatyki / kolejki” przekazująca zapis do modułu Automations (utworzenie zadania korekty) | Domknięcie pętli obserwacja → działanie |
| Maskowanie w eksporcie | Przy eksporcie i przekazaniu — ustawienie maskujące tokeny w bloku `argv`; stan domyślny: wartości jawne | Udostępnianie prowenancji bez ujawniania kluczy |

### 10.4. Diagram przepływu danych

```
WYWOŁANIE MODELU (dowolny kanał — rozdz. 3.7; Chat Window albo Execution Loop Window)
   │
   ├── kompozycja nakładki (rozdz. 2.3) ──────────────► blok „system prompt”
   ├── rozstrzygnięcie pól wywołania (rozdz. 3, dziedziczenie warstw) ─► bloki „settings” i „warstwy i pochodzenie”
   ├── zbudowanie parametrów technicznych adaptera (Integracja modeli, rozdz. 4) ─► blok „argv”
   └── odczyt hasha Konstytucji i wersji eksperta ────► blok „hash konstytucji i wersji eksperta”
                                                              │
                                                              ▼
                                                   PANEL PROWENANCJI
                                          ┌──────────────────┴──────────────────┐
                                          ▼                                     ▼
                                Podgląd na żywo                       Zapis archiwalny
                                (przed wysłaniem,                     (po wywołaniu, powiązany
                                 na żądanie)                           z wpisem historii — 5.9)
```

---

## 11. Import i eksport konfiguracji

Konfiguracja nakładki, profili ekspertów, kanałów i profili parametrów przenosi się między instancjami i urządzeniami jako pliki w jawnym, czytelnym formacie.

| Element | Działanie | Efekt |
|---|---|---|
| Eksport konfiguracji do pliku | Zapis wybranego zakresu (warstwa, profil eksperta, kanał, profil parametrów, zestaw hooków) do pliku JSON lub YAML, spójnego z szablonem Załącznika A | Kopia zapasowa, przenoszenie między instancjami, dokumentacja |
| Import konfiguracji z pliku | Wczytanie pliku z podglądem różnic wobec stanu bieżącego przed zatwierdzeniem | Wdrożenie wcześniejszej lub cudzej konfiguracji z pełną wiedzą o skutkach |
| Format pliku | Struktura pliku odwzorowuje szablon „Nakładka i wywołanie modelu” (Załącznik A): `nakladka`, `pola_wywolania`, `kanaly_komunikacji`, `warstwy_widocznosci`, `panel_prowenancji`; dane dostępowe zapisywane jako odwołania | Przejrzystość, edytowalność ręczna, brak wycieku danych dostępowych |
| Eksport profilu eksperta | Spakowanie profilu wraz z wersją, hashem i notatką do jednego pliku | Współdzielenie sprawdzonych osobowości roboczych między urządzeniami i instancjami |
| Import selektywny | Wybór pól i sekcji do zaimportowania, z decyzją „nadpisz / dołóż / pomiń” dla każdej pozycji | Łączenie konfiguracji zamiast wymiany całości |
| Migawki konfiguracji | Nazwana migawka pełnego stanu okna Konfiguracji z datą i notatką; lista migawek z akcją „Przywróć” | Punkty przywracania przed większymi zmianami |
| Zakres eksportu danych dostępowych | Eksport obejmuje odwołania; ustawienie „dołącz wartości jawne” udostępnia wartości dla przenosin bez połączenia sieciowego | Zgodność z zasadą „dane poza bazą” przy zachowanej przenośności |
| Zgodność wersji schematu | Plik niesie numer wersji schematu konfiguracji; różnica wersji wywołuje ostrzeżenie i mapowanie pól, nie odrzucenie pliku | Trwałość plików konfiguracji |

Import i eksport korzystają z istniejącego kanału komunikacji (`config.get`, `config.set`) i operacji serializacji po stronie klienta; migawki i wersje są rozgłaszane zdarzeniem `config.changed` do wszystkich połączonych urządzeń.

---

## 12. Makieta tekstowa okna

Wszystkie makiety przedstawiają układ pionowy — podział lewa–prawa. Kolumny sąsiadują poziomo, okna pomocnicze otwierają się jako rozszerzenia boczne po prawej stronie obszaru roboczego, a regulacji podlega wyłącznie szerokość kolumn. Makiety rysowane są w stanie spoczynku interfejsu: widoczne są elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisuje katalog elementów (rozdz. 13) z podaniem warstwy i sposobu wywołania.

### 12.1. Pozycja okna Konfiguracji w układzie obszaru roboczego

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │ — okno Konfiguracji      │ pomocniczy
  modułów   │ Wykonawca            │   (rozdz. 12.2–12.4)     │ (rozszerzenie
            │                      │                          │  boczne)
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
```

### 12.2. Powłoka okna — nawigacja zakresów i zawartość

```
┌──────────────────────────┬────────────────────────────────────────────┐
│ NAWIGACJA ZAKRESÓW       │ ZAWARTOŚĆ ZAKRESU                          │
│                          │                                            │
│  Aplikacja               │  [ nazwa ustawienia ]                  [?] │
│  Procesy                 │      wartość bieżąca · poziom warstwy      │
│  Akcje                   │                                            │
│▸ Nakładka i wywołanie    │  ← rozwinięte tutaj; treść pełna w 12.3    │
│   modelu                 │                                            │
│   • Zachowanie modeli    │  Warstwa: [ globalna ▼ ]                   │
│   • Tożsamość modeli     │  Ekspert: [ Roman ▼ ]           [ ⋮ ]      │
│   • Prompty systemowe    │                                            │
│  Kanały komunikacji      │                                            │
│  Warstwy widoczności     │                                            │
│  Rozszerzenia            │                                            │
│  Integracje              │                                            │
│  Historia                │                                            │
│  Pamięć                  │                                            │
│  Izolacja  → rozdz. 6 MK │                                            │
│  Karty sesji             │                                            │
│  Komponenty własne       │                                            │
│                          │                                            │
│  [ szukaj ustawienia ]   │                                            │
└──────────────────────────┴────────────────────────────────────────────┘
```

### 12.3. Widok szczegółowy — Nakładka i wywołanie modelu

Pełny, trzykolumnowy układ obszaru roboczego, wywoływany z dowolnej z trzech pozycji nawigacji (Zachowanie modeli / Tożsamość modeli / Prompty systemowe). Kolumny sąsiadują poziomo; kolumna prawa mieści Panel prowenancji.

```
┌──────────────────────┬────────────────────────────────────────────────┬───────────────────────────────┐
│ SELEKTOR WARSTWY     │ NAKŁADKA I POLA WYWOŁANIA                      │ PANEL PROWENANCJI             │
│ NAKŁADKI             │                                                │ „co poszło do modelu”         │
│                      │                                                │                               │
│ • Domyślne platformy │ ── WARSTWA: KONSTYTUCJA (Globalna) ──          │ tryb: Code CLI                │
│ • Globalna           │                                                │ zasięg: rola Executor 1       │
│   (Konstytucja)      │ ┌──────────────────────────────────────┐       │ kanał: Chat Window            │
│ • Środowisko         │ │ [ edytor tekstu — treść Konstytucji ]│[?]    │                               │
│ • Projekt            │ └──────────────────────────────────────┘       │ ▸ argv                        │
│ • Sesja              │ [ Zapisz ]  hash: 8f3a2c1d…  [ kopiuj ]        │ ▸ system prompt               │
│ • Rola               │                                                │ ▸ settings                    │
│   (MultitaskingAI)   │ ── WARSTWA: PROFIL / ROLA ──                   │ ▸ warstwy i pochodzenie       │
│ • Agent              │ Ekspert: [ Roman ▼ ]  wersja 7   [ ⋮ ]         │ ▸ hash konstytucji i wersji   │
│   (komponent własny) │ Agent: „Agent Redaktor”                        │                               │
│                      │ [ podgląd tożsamości i instrukcji — ]          │   ┌─────────────────────────┐ │
│ zaznaczenie:         │ [ tylko do odczytu ]                           │   │ .dn-kod — zawartość     │ │
│  tło akcentu         │ Edytuj w Agent Builder →                       │   │ wybranego bloku,        │ │
│                      │                                                │   │ dane dostępowe          │ │
│                      │ ── WARSTWA: EKSPERTYZA ZADANIOWA ──            │   │ w postaci jawnej        │ │
│                      │ Projekt: „Sprawa X”                            │   └─────────────────────────┘ │
│                      │ [ podgląd instrukcji projektu — ]              │                               │
│                      │ [ tylko do odczytu ]  Edytuj w                 │ [ Podgląd na żywo ]           │
│                      │ Instructions Panel →                           │ [ Porównaj wywołania ]        │
│                      │ Budowa dynamiczna przez Koordynatora [wł.]     │ Zapis archiwalny:     [wł.]   │
│                      │                                                │                               │
│                      │ ── POLA WYWOŁANIA (ta warstwa) ──              │                               │
│                      │ Harness      [ pochodny od kanału ]    [?]     │                               │
│                      │ Hooki        [ + dodaj punkt zaczepienia ] [?] │                               │
│                      │              po odpowiedzi → „Zapis do         │                               │
│                      │              repertorium”                      │                               │
│                      │ Narzędzia    [ pełny dostęp ▼ ]        [?]     │                               │
│                      │ Model        [ wybór modelu bazowego ▼][?]     │                               │
│                      │ Effort       [wł.] [ dogłębny        ▼][?]     │                               │
│                      │ Kanał        [ Code CLI ▼ ]            [?]     │                               │
│                      │  Code CLI →  Pula kont: [ zarządzaj… ]         │                               │
│                      │              Strategia: [ round-robin ▼ ]      │                               │
│                      │ Parametry próbkowania         [ ⋮ ]            │                               │
│                      │                                                │                               │
│                      │ [ Otwórz okno konfiguracji punktów             │                               │
│                      │   izolacji → ]                                 │                               │
└──────────────────────┴────────────────────────────────────────────────┴───────────────────────────────┘
```

### 12.4. Widok — Kanały komunikacji

Zakres „Kanały komunikacji” prezentuje ustawienia obu kanałów w dwóch sąsiadujących kolumnach, w tej samej kolejności, w jakiej występują w obszarze roboczym: Chat Window po lewej, Execution Loop Window w kolumnie sąsiadującej.

```
┌────────────────────────────────────┬────────────────────────────────────┬──────────────────────┐
│ CHAT WINDOW                        │ EXECUTION LOOP WINDOW              │ PODGLĄD USTAWIEŃ      │
│ Użytkownik ↔ Wykonawca             │ Koordynator ↔ Wykonawca            │ EFEKTYWNYCH           │
│                                    │                                    │                       │
│ Wykonawca      [ Agent Redaktor ▼ ]│ Limit kroków pętli   [ 25 ]        │ ▸ chat window         │
│ Ekspert        [ Roman ▼ ]         │ Limit czasu kroku    [ 120 s ]     │ ▸ execution loop      │
│ Model          [ wybór ▼ ]    [?]  │ Równoległość zadań   [ 3 ]    [?]  │ ▸ warstwa źródłowa    │
│ Effort         [wł.] [ dogłębny ▼ ]│ Polityka ponowień    [ 2 × ]  [?]  │                       │
│ Strumień       [ strumieniowo ▼ ]  │ Próg kontroli jakości[ 0,80 ] [?]  │   ┌─────────────────┐ │
│ Kroki pośrednie[ zwinięte ▼ ]      │ Poniżej progu        [ ponów ▼ ]   │   │ .dn-kod —       │ │
│ Zatwierdzanie  [wł.]          [?]  │ Zakres autonomii     [ ⋮ ]    [?]  │   │ wartości i      │ │
│ Punkty zatwierdzeń       [ ⋮ ]     │ Punkty zatwierdzeń   [ ⋮ ]         │   │ poziom warstwy  │ │
│ Wyjaśnianie wyniku  [wł.]          │ Sterowanie przebiegiem [ ⋮ ]       │   └─────────────────┘ │
│ Szerokość kolumny [ 3 jedn. ▼ ]    │ Prezentacja pętli    [ pełna ▼ ]   │                       │
│ Historia i pamięć        [ ⋮ ]     │ Szerokość kolumny    [ 3 jedn. ▼ ] │ [ Podgląd na żywo ]   │
└────────────────────────────────────┴────────────────────────────────────┴──────────────────────┘
```

### 12.5. Widok — Warstwy widoczności

```
┌──────────────────────┬────────────────────────────────────────────────┬───────────────────────┐
│ PROFILE WARSTW       │ MACIERZ PRZYPISANIA FUNKCJI                    │ PODGLĄD W ROLI        │
│                      │                                                │                       │
│ • Domyślny platformy │ [ szukaj funkcji ]                             │ rola: [ podstawowy ▼ ]│
│ • Użytkownik         │                                                │                       │
│   podstawowy         │ Funkcja              Warstwa  Sposób wywołania │ ┌───────────────────┐ │
│ • Użytkownik         │ Chat Window            1      widoczna         │ │ stan spoczynku    │ │
│   zaawansowany       │ Wybór modelu           2      znacznik kontek. │ │ interfejsu        │ │
│ • Administrator      │ Zestaw operacji        3      menu ⋮           │ │ widziany przez    │ │
│                      │ Tryb administracyjny   4      paleta poleceń   │ │ wybraną rolę      │ │
│ [ + Nowy profil ]    │ Diagnostyka procesu    4      wyszukiwarka     │ └───────────────────┘ │
│                      │                                                │                       │
│                      │ [ Przywróć przypisanie domyślne ]              │ [ Katalog skrótów ]   │
│                      │ Tryb administracyjny: wejście z palety poleceń │ [ Paleta poleceń ]    │
└──────────────────────┴────────────────────────────────────────────────┴───────────────────────┘
```

### 12.6. Widok — Pula kont Code CLI

```
┌──────────────────────────────────────────────────────────┐
│  PULA KONT CODE CLI                                   [x] │
├──────────────────────────────────────────────────────────┤
│  Konto 1 — „Ogólne”        ● dostępne   47/500  [ ⋮ ]     │
│  Konto 2 — „Sprawa X”      ● dostępne  312/500  [ ⋮ ]     │
│  Konto 3 — „Sprawa Y”      ● zajęte    128/500  [ ⋮ ]     │
│  Konto 4 — „Zapasowe”      ● błąd uwierz.       [ Popraw ]│
│                                                            │
│  [ + Dodaj konto ]                                        │
│                                                            │
│  Strategia przydziału:  [ kolejno (round-robin) ▼ ]       │
│  Katalog profilu konta: [ ścieżka ]                       │
│  Rejestr zdarzeń kont:  [ ⋮ ]                             │
│                                                            │
│                                  [ Anuluj ]  [ Zapisz ]   │
└──────────────────────────────────────────────────────────┘
```

### 12.7. Widok — Panel prowenancji rozwinięty

Rozwinięcie Panelu prowenancji poszerza kolumnę prawą kosztem szerokości kolumn sąsiednich; podział pozostaje lewa–prawa.

```
┌───────────────────────────────────────────────────────────┐
│ PANEL PROWENANCJI                                            │
│ „co poszło do modelu”                                        │
│                                                                │
│ tryb: Code CLI · zasięg: rola Executor 1 · zespół „Sprawa X”  │
│ kanał komunikacji: Execution Loop Window                      │
│ znacznik czasu wywołania: 2026-08-06 14:32:07                │
│                                                                │
│ ┌── argv ─────────────────────────────────────────────────┐ │
│ │ narzędzie: <nazwa narzędzia Code CLI>                     │ │
│ │ --konto=Konto 2 --model=<nazwa modelu>                    │ │
│ │ --token=9f83a91f                                          │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                                │
│ ┌── system prompt ───────────────────────────────────────┐ │
│ │ # KONSTYTUCJA (poziom: globalna)                          │ │
│ │ …                                                          │ │
│ │ # PROFIL / ROLA — Agent Redaktor (poziom: agent)           │ │
│ │ …                                                          │ │
│ │ # EKSPERTYZA ZADANIOWA (poziom: projekt + sekcja            │ │
│ │ zbudowana przez Koordynatora dla kroku 3)                  │ │
│ │ …                                                          │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                                │
│ ┌── settings ────────────────────────────────────────────┐ │
│ │ (format — rozdz. 10.2)                                     │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                                │
│ ┌── warstwy i pochodzenie ───────────────────────────────┐ │
│ │ model   : agent      (dziedziczone: globalna)              │ │
│ │ effort  : rola       (dziedziczone: agent)                 │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                                │
│ ┌── hash konstytucji i wersji eksperta ──────────────────┐ │
│ │ konstytucja: 8f3a2c1d9e0b…        [ kopiuj ]              │ │
│ │ ekspert Roman, wersja 7: 41c7be02…[ kopiuj ]              │ │
│ └────────────────────────────────────────────────────────┘ │
│                                                                │
│ [ Podgląd na żywo ] [ Porównaj ] [ Eksportuj ]  Zapis: [wł.]  │
└───────────────────────────────────────────────────────────┘
```

---

## 13. Katalog elementów interfejsu

Katalog opisuje element po elemencie wszystkie elementy okna. Kolumna „Forma i waga” rozstrzyga skalę elementu, tak aby żaden element pomocniczy nie został zbudowany jako duża kolumna, a żaden element główny — jako mała ikona. Kolumny „Warstwa” i „Sposób wywołania” wskazują warstwę widoczności elementu (rozdz. 9.1) i sposób jego ujawnienia.

| # | Element | Co to jest / do czego służy | Warstwa | Sposób wywołania | Forma i waga | Stany | Zachowanie po interakcji |
|---|---|---|---|---|---|---|---|
| 1 | Pozycja nawigacji zakresów | Wejście do jednego z zakresów okna Konfiguracji | 1 | Widoczna bez interakcji | Pozycja listy pionowej (`.dn-karta--pozycja`), pełna szerokość kolumny lewej | domyślny, najechanie, zaznaczony | Kliknięcie ładuje odpowiedni zakres w kolumnie zawartości |
| 2 | Grupa „Nakładka i wywołanie modelu” | Nadrzędna pozycja nawigacji scalająca trzy zakresy (5.4–5.6) | 1 | Widoczna bez interakcji | Pozycja listy z trzema podpozycjami, wcięcie hierarchiczne | domyślny, rozwinięty, zaznaczony | Kliknięcie pozycji nadrzędnej albo dowolnej podpozycji otwiera ten sam obszar roboczy (12.3) |
| 3 | Selektor warstwy nakładki | Wybór warstwy zasięgu, dla której edytowane są pola kolumny środkowej | 1 | Widoczny bez interakcji | Lista pionowa pozycji (`.dn-karta--pozycja`), kolumna lewa okna 12.3 | domyślny, najechanie, zaznaczony | Kliknięcie przełącza zawartość kolumny środkowej i prawej na wybraną warstwę |
| 4 | Edytor tekstu Konstytucji | Pole edycji pełnej treści Konstytucji dla wybranego zasięgu | 1 | Widoczny bez interakcji | Duże pole tekstowe wieloliniowe (`.dn-textarea`), pełna szerokość sekcji | domyślny, edycja, zapisywanie, zapisano, błąd zapisu, pusty | Wpisywanie tekstu; przycisk „Zapisz” zatwierdza zmianę i przelicza hash |
| 5 | Przycisk „Zapisz” (Konstytucja) | Zatwierdzenie edytowanej treści Konstytucji | 1 | Widoczny bez interakcji | Przycisk drugorzędny (`.dn-btn--zarys`), mały | domyślny, najechanie, wciśnięty, ładowanie, błąd | Zawsze aktywny; wysyła `config.set`; kliknięcie bez zmian nie wywołuje ruchu sieciowego; po potwierdzeniu serwera toast sukcesu i odświeżenie hasha |
| 6 | Wyświetlacz hasha Konstytucji | Skrócony skrót SHA-256 treści obowiązującej Konstytucji | 1 | Widoczny bez interakcji | Etykieta monospaced z przyciskiem ikonowym kopiowania | domyślny, po kopiowaniu, pusty | Kliknięcie ikony kopiuje pełną wartość do schowka, pokazuje toast „Skopiowano” |
| 7 | Selektor profilu eksperta | Wybór aktywnego profilu eksperta dla wybranej warstwy | 2 | Menu progresywne `Ekspert ▼` w kolumnie środkowej i w pasku kontekstu | Lista rozwijana (`.dn-select`) z etykietą wersji | domyślny, otwarta, wybrany, brak aktywnego | Wybór natychmiast aktywuje profil i odświeża pola nakładki oraz Panel prowenancji |
| 8 | Menu profilu eksperta | Operacje na profilu: wersje, przywrócenie, wariant, porównanie wersji, galeria, harmonogram, eksport | 3 | Menu kebab (⋮) przy selektorze eksperta | Menu kontekstowe (`.dn-menu`) | domyślny, otwarte | Wybór pozycji otwiera właściwy widok operacji |
| 9 | Karta podglądu Profil/Rola | Podgląd tożsamości i instrukcji systemowych przypisanego agenta lub roli, tylko do odczytu | 1 | Widoczna bez interakcji | Karta (`.dn-karta`), średnia waga, treść przewijalna | domyślny, ładowanie, pusty | Nieinteraktywna poza przewijaniem — edycja przez skrót nr 10 |
| 10 | Skrót „Edytuj w Agent Builder →” / „Edytuj w sekcji Role →” | Przejście do właściwego edytora danej warstwy | 1 | Widoczny bez interakcji | Łącze tekstowe akcentowane, małe | domyślny, najechanie, fokus | Kliknięcie nawiguje do wskazanego okna modułu Agents lub panelu orkiestracji |
| 11 | Karta podglądu Ekspertyzy zadaniowej | Podgląd instrukcji systemowych projektu i skilli, tylko do odczytu | 1 | Widoczna bez interakcji | Karta (`.dn-karta`), średnia waga | domyślny, ładowanie, pusty | Nieinteraktywna poza przewijaniem |
| 12 | Przełącznik „Budowa dynamiczna przez Koordynatora” | Włącza dokładanie ekspertyzy zadaniowej przez Koordynatora dla Wykonawców | 2 | Widoczny po rozwinięciu sekcji Ekspertyza zadaniowa | Przełącznik (`.dn-suwak`) z etykietą | włączony, wyłączony | Kliknięcie przełącza stan natychmiast (`config.set`); bez Koordynatora w zespole ustawienie zapisuje się z komunikatem o warunku jego działania |
| 13 | Pole „Negatyw”, „Styl i ton”, „Słowniczek”, „Format wyjścia”, „Język odpowiedzi”, „Budżet” | Pola treści nakładki (rozdz. 2.7) | 2 | Zwinięta sekcja „Pola treści nakładki” rozwijana jednym kliknięciem | Pola tekstowe i listy rozwijane (`.dn-textarea`, `.dn-select`) | domyślny, edycja, dziedziczony, ustawiony jawnie | Zmiana wysyła `config.set`; wartość i warstwa widoczne w plakietce pola |
| 14 | Plakietka warstwy pola | Wskazanie warstwy, na której ustawiono wartość, albo warstwy dziedziczenia | 1 | Widoczna bez interakcji przy każdym polu | Plakietka (`.dn-plakietka`), mała | ustawione jawnie, dziedziczone | Kliknięcie otwiera podgląd łańcucha dziedziczenia (element 15) |
| 15 | Podgląd łańcucha dziedziczenia | Rozwinięcie pokazujące wszystkie warstwy pola z zaznaczeniem warstwy zwycięskiej | 3 | Rozwinięcie kontekstowe (popover) z plakietki warstwy | Panel popover (`.dn-popover`) | domyślny, ładowanie | Zawiera akcje „Przypnij / Odepnij / Przywróć dziedziczenie” |
| 16 | Pole „Harness” | Podgląd i edycja parametrów pętli agentowej właściwej wybranemu kanałowi | 2 | Widoczne w sekcji „Pola wywołania” | Pole tekstowe lub lista rozwijana (`.dn-select`) zależna od kanału | domyślny, edytowalny, tylko do odczytu (Code CLI), zapisywanie | Zmiana wysyła `config.set`; dla Code CLI pole pozostaje klikalne i pokazuje wartość tylko do odczytu z objaśnieniem `[?]` |
| 17 | Lista hooków | Zestaw podłączonych punktów zaczepienia wraz z przypisaną automatyką | 2 | Widoczna w sekcji „Pola wywołania” | Lista pozycji (`.dn-karta--pozycja`), każda z etykietą punktu i akcji | domyślny, pusty, edycja pozycji | Kliknięcie pozycji otwiera edycję przypisanej automatyki; ikona usunięcia odłącza hook |
| 18 | Przycisk „+ dodaj punkt zaczepienia” | Dodanie nowego hooka | 2 | Widoczny przy liście hooków | Przycisk drugorzędny (`.dn-btn--zarys`), mały, z ikoną plusa | domyślny, najechanie, wciśnięty | Otwiera formularz wyboru punktu cyklu i wskazania automatyki |
| 19 | Biblioteka gotowych hooków | Zestawy startowe hooków zasięgu globalnego | 3 | Menu kebab (⋮) przy liście hooków | Lista pozycji w panelu popover | domyślny, wybrany | Wybór pozycji podłącza hook do wskazanego punktu cyklu |
| 20 | Pole „Narzędzia dozwolone i zabronione” | Wybór zestawu rozszerzeń i akcji dostępnych modelowi | 2 | Widoczne w sekcji „Pola wywołania” | Dwie listy wielokrotnego wyboru (`.dn-select` wariant multi) albo etykieta „pełny dostęp” tylko do odczytu | domyślny (pełny dostęp), zawężony, tylko do odczytu, ładowanie rejestru | Zaznaczenie pozycji aktualizuje zestaw; dla wywołań z agentem pole pokazuje wartość dziedziczoną ze skrótem „Edytuj w Permissions Center →” |
| 21 | Grupy narzędzi | Nazwane zestawy uprawnień zaznaczane jednym kliknięciem | 3 | Lista rozwijana przy polu narzędzi | Lista rozwijana (`.dn-select`) | domyślny, wybrany | Wybór grupy podmienia zawartość obu list narzędzi |
| 22 | Podgląd efektywnego zestawu narzędzi | Wynikowa lista narzędzi po złożeniu list i dziedziczenia | 3 | Rozwinięcie kontekstowe przy polu narzędzi | Blok kodu (`.dn-kod`), przewijalny | domyślny, ładowanie | Tylko do odczytu; odświeża się przy każdej zmianie list |
| 23 | Pole „Model” | Wybór modelu bazowego z listy właściwej kanałowi albo wpisanie nazwy wprost | 2 | Widoczne w sekcji „Pola wywołania” | Lista rozwijana z polem wpisu (`.dn-select`) | domyślny, otwarta, wybrany, ładowanie listy, lista niedostępna | Wybór natychmiast aktualizuje ustawienie i odświeża Panel prowenancji |
| 24 | Pole „Model zapasowy” | Wskazanie modelu i kanału uruchamianego przy niedostępności podstawowego | 3 | Rozwinięcie kontekstowe przy polu „Model” | Lista rozwijana (`.dn-select`) | domyślny, wybrany, pusty | Wybór zapisuje ustawienie dla warstw rola i agent |
| 25 | Przełącznik i pole „Effort” | Włączenie nadpisania wartości effort na tej warstwie oraz wybór poziomu | 2 | Widoczne w sekcji „Pola wywołania” | Przełącznik (`.dn-suwak`) + lista rozwijana (`.dn-select`) w jednym wierszu | wyłączony (dziedziczy), włączony z wartością, ładowanie | Włączenie przełącznika odsłania listę poziomów wraz z szacunkiem czasu i kosztu |
| 26 | Parametry próbkowania | Temperatura, top-p, kara za powtórzenia, sekwencje stop, `seed` | 3 | Menu kebab (⋮) w sekcji „Pola wywołania” | Panel popover z polami liczbowymi i tekstowymi | wyłączone (wartość modelu), ustawione jawnie | Zmiana zapisuje się natychmiast, z plakietką warstwy przy każdym parametrze |
| 27 | Pole „Kanał” | Wybór trybu wywołania: Code CLI, Agent SDK, API | 2 | Menu progresywne `Kanał ▼` w sekcji „Pola wywołania” | Lista rozwijana (`.dn-select`) | domyślny, otwarta, wybrany | Wybór natychmiast przełącza zestaw pól zależnych (harness, model, pula kont) |
| 28 | Katalog kanałów modelu | Widok wszystkich encji kanału modelu z akcjami dodania, edycji, usunięcia, powielenia | 3 | Rozwinięcie kontekstowe przy polu „Kanał” | Lista kart (`.dn-karta`) w panelu wysuwanym | domyślny, ładowanie, pusty | Wybór pozycji otwiera edytor parametrów właściwych typowi kanału |
| 29 | Przycisk „Testuj połączenie” | Próbne wywołanie kanału z prezentacją wyniku | 3 | Widoczny w edytorze kanału i w wierszu konta puli | Przycisk drugorzędny (`.dn-btn--zarys`) | domyślny, ładowanie, sukces, błąd uwierzytelnienia, model niedostępny | Wynik prezentowany plakietką stanu przy polu |
| 30 | Podgląd porównawczy modeli | Uruchomienie tego samego polecenia przez dwa lub trzy kanały jednocześnie | 4 | Polecenie języka naturalnego w Chat Window, skrót klawiszowy albo wyszukiwarka funkcji | Kolumny odpowiedzi zestawione obok siebie w oknie pomocniczym | domyślny, ładowanie, gotowy | Zestawia odpowiedzi obok siebie; wybór odpowiedzi ustawia model w polu „Model” |
| 31 | Przycisk „zarządzaj…” (Pula kont) | Otwarcie widoku zarządzania pulą kont Code CLI | 2 | Widoczny w podsekcji Code CLI | Łącze tekstowe, małe | domyślny, najechanie | Otwiera widok „Pula kont Code CLI” (12.6) |
| 32 | Pole „Strategia przydziału” | Wybór strategii rozdziału wywołań między konta puli | 2 | Widoczne w podsekcji Code CLI | Lista rozwijana (`.dn-select`) | domyślny, otwarta, wybrany | Wybór natychmiast aktualizuje ustawienie |
| 33 | Wiersz konta puli | Pojedyncze konto: etykieta, stan, obciążenie, limit, katalog profilu | 2 | Widoczny w widoku puli kont | Pozycja listy (`.dn-karta--pozycja`) z plakietką stanu (`.dn-plakietka--stan`) | dostępne, zajęte, błąd uwierzytelnienia, limit bliski wyczerpania | Menu (⋮) wiersza zawiera „Testuj połączenie”, „Ustaw jako aktywne”, „Usuń”; „Popraw” przy błędzie otwiera ponowne uwierzytelnienie; usunięcie odwracalne akcją „Cofnij” |
| 34 | Przycisk „+ Dodaj konto” | Dodanie nowego konta do puli | 2 | Widoczny w widoku puli kont | Przycisk drugorzędny (`.dn-btn--zarys`) | domyślny, najechanie | Otwiera formularz uwierzytelnienia nowego konta narzędzia Code CLI |
| 35 | Rejestr zdarzeń kont | Dziennik przełączeń konta, błędów uwierzytelnienia i przekroczeń limitów | 4 | Menu kebab (⋮) widoku puli kont albo wyszukiwarka funkcji | Tabela przewijalna, znacznik czasu i urządzenie | domyślny, pusty, filtrowany | Tylko do odczytu; pozycje filtrowane po koncie i rodzaju zdarzenia |
| 36 | Rejestr kluczy i danych dostępowych | Widok wszystkich odwołań do danych dostępowych z wartością jawną i przełącznikiem widoczności | 3 | Pozycja zakresu „Aplikacja”, rozwinięcie kontekstowe | Tabela pozycji (`.dn-karta--pozycja`) | domyślny, wartość ukryta, wartość jawna, błąd | Przełącznik „Pokaż / Ukryj” zmienia widoczność wartości; pozycja wskazuje miejsca użycia klucza |
| 37 | Edytor zmiennych środowiskowych | Tabela par `KLUCZ = wartość` wstrzykiwanych do procesu sesji | 3 | Rozwinięcie kontekstowe zakresu „Procesy” | Tabela edytowalna z akcjami importu i eksportu pliku `.env` | domyślny, edycja, import, eksport, błąd odczytu pliku | Zmiana zapisuje się per warstwa; import prezentuje różnice przed zatwierdzeniem |
| 38 | Skrót „Otwórz okno konfiguracji punktów izolacji →” | Przejście do trzykolumnowego okna izolacji, z zaznaczonym poziomem zasięgu | 2 | Widoczny na końcu kolumny środkowej | Przycisk drugorzędny (`.dn-btn--zarys`), pełna szerokość sekcji | domyślny, najechanie | Kliknięcie otwiera okno izolacji (Model konfiguracji, rozdz. 6) |
| 39 | Ikona objaśnienia kontekstowego `[?]` | Wywołanie treści objaśniającej działanie i wpływ ustawienia | 2 | Najechanie kursorem albo fokus klawiaturą | Mała ikona (`.dn-ikona--sm`, 14 px) w dymku (`.dn-tooltip`) | domyślny, dymek widoczny | Prezentuje treść z rozdz. 3.9 i z definicji ustawienia |
| 40 | Nagłówek trybu (Panel prowenancji) | Wskazanie trybu wywołania, zasięgu i kanału komunikacji bieżącego podglądu | 1 | Widoczny bez interakcji | Etykieta tekstowa, mała, nad blokami treści | domyślny | Nieinteraktywny — aktualizuje się wraz z wyborem w kolumnach lewej i środkowej |
| 41 | Blok „argv” | Dosłowny zapis parametrów technicznych wywołania, z danymi dostępowymi domyślnie jawnymi | 2 | Nagłówek sekcji `▸` w Panelu prowenancji | Blok kodu (`.dn-kod`), przewijalny, monospaced | domyślny, ładowanie, pusty | Rozwijany i zwijany nagłówkiem sekcji; zawartość tylko do odczytu, z przyciskiem kopiowania |
| 42 | Blok „system prompt” | Pełny, skomponowany tekst nakładki z oznaczonymi granicami warstw i sekcją Koordynatora | 2 | Nagłówek sekcji `▸` w Panelu prowenancji | Blok kodu (`.dn-kod`), przewijalny, z wyróżnieniem nagłówków sekcji | domyślny, ładowanie, pusty | Rozwijany i zwijany; kopiowanie całości albo pojedynczej warstwy |
| 43 | Blok „settings” | Wynikowe wartości pól wywołania po dziedziczeniu, z oznaczeniem poziomu każdej wartości | 2 | Nagłówek sekcji `▸` w Panelu prowenancji | Blok kodu (`.dn-kod`) | domyślny, ładowanie | Rozwijany i zwijany; tylko do odczytu |
| 44 | Blok „warstwy i pochodzenie” | Warstwa zwycięska i wartość dziedziczona dla każdego pola | 3 | Nagłówek sekcji `▸` w Panelu prowenancji | Blok kodu (`.dn-kod`) | domyślny, ładowanie | Tylko do odczytu; kliknięcie pozycji przenosi do pola w kolumnie środkowej |
| 45 | Blok „hash konstytucji i wersji eksperta” | Skrócona i pełna wartość skrótów kryptograficznych | 2 | Nagłówek sekcji `▸` w Panelu prowenancji | Blok kodu (`.dn-kod`) dwuwierszowy + przyciski kopiowania | domyślny, pusty | Kopiowanie pełnej wartości do schowka |
| 46 | Przycisk „Podgląd na żywo” | Przeliczenie i pokazanie Panelu prowenancji dla stanu bieżącego, bez wysyłania wywołania | 1 | Widoczny bez interakcji w Panelu prowenancji | Przycisk główny (`.dn-btn--glowny`) | domyślny, najechanie, wciśnięty, ładowanie | Przelicza i odświeża bloki kolumny prawej |
| 47 | Przycisk „Porównaj wywołania” | Zestawienie prowenancji dwóch wpisów archiwalnych pole po polu | 3 | Widoczny w Panelu prowenancji, rozwija widok porównania | Przycisk drugorzędny (`.dn-btn--zarys`) | domyślny, ładowanie, gotowy | Otwiera dwie kolumny porównania z podświetleniem różnic |
| 48 | Oś czasu prowenancji | Chronologiczny wykaz wywołań sesji lub roli z hashami i znacznikami czasu | 3 | Rozwinięcie kontekstowe Panelu prowenancji | Lista pozycji przewijalna, z polem filtra | domyślny, filtrowany, pusty | Kliknięcie pozycji ładuje pełny zapis archiwalny |
| 49 | Przycisk „Eksportuj” (prowenancja) | Zapis kompletu bloków do pliku Markdown lub JSON | 3 | Widoczny w Panelu prowenancji | Przycisk drugorzędny (`.dn-btn--zarys`) | domyślny, ładowanie | Otwiera wybór formatu i ustawienia maskowania danych dostępowych |
| 50 | Przełącznik „Zapis archiwalny” | Włącza automatyczny zapis kompletu bloków przy każdym faktycznym wywołaniu | 2 | Widoczny w Panelu prowenancji | Przełącznik (`.dn-suwak`) z etykietą | włączony (domyślny), wyłączony | Kliknięcie przełącza stan natychmiast (`config.set`) |
| 51 | Pola kanałów komunikacji | Ustawienia Chat Window i Execution Loop Window (rozdz. 8) | 1–2 | Kolumny zakresu „Kanały komunikacji”; ustawienia szczegółowe pod `⋮` | Pola tekstowe, listy rozwijane, przełączniki | domyślny, dziedziczony, ustawiony jawnie, ładowanie | Zmiana zapisuje się natychmiast i odświeża podgląd ustawień efektywnych |
| 52 | Macierz przypisania funkcji do warstw | Tabela funkcji interfejsu z kolumnami warstwy i sposobu wywołania | 3 | Pozycja nawigacji „Warstwy widoczności” | Tabela przewijalna z polem wyszukiwania | domyślny, filtrowany, zmieniony, ostrzeżenie reguły jednego kliknięcia | Zmiana warstwy zapisuje się per profil roli; ostrzeżenie przy naruszeniu reguły jednego kliknięcia |
| 53 | Profil warstw roli | Nazwany zestaw przypisań warstw dla roli użytkownika | 3 | Kolumna lewa zakresu „Warstwy widoczności” | Lista pozycji (`.dn-karta--pozycja`) | domyślny, zaznaczony, dziedziczony z globalnej | Wybór profilu ładuje macierz przypisań tej roli |
| 54 | Podgląd interfejsu w roli | Podgląd stanu spoczynku interfejsu widzianego przez wybraną rolę | 3 | Kolumna prawa zakresu „Warstwy widoczności” | Karta podglądu (`.dn-karta`) | domyślny, ładowanie | Zmiana roli w selektorze przerysowuje podgląd, bez zmiany roli bieżącej |
| 55 | Wyszukiwarka funkcji i ustawień | Dotarcie do dowolnej funkcji i dowolnego ustawienia po nazwie i treści objaśnienia | 4 | Skrót klawiszowy albo pole „szukaj ustawienia” w pasku kontekstu | Pole tekstowe z listą wyników (`.dn-select` wariant wyszukiwania) | domyślny, wyniki, brak wyników | Wybór wyniku otwiera zakres i zaznacza ustawienie |
| 56 | Katalog skrótów klawiszowych | Lista skrótów z przypisaną funkcją i możliwością zmiany kombinacji | 4 | Wyszukiwarka funkcji albo skrót klawiszowy | Tabela edytowalna | domyślny, edycja, konflikt kombinacji | Konflikt wywołuje ostrzeżenie ze wskazaniem funkcji zajmującej skrót; zapis pozostaje możliwy |
| 57 | Znacznik stanu trybu administracyjnego | Sygnalizacja działania trybu administracyjnego i wyjście z niego | 1 w czasie działania trybu | Widoczny bez interakcji od chwili wejścia w tryb; wejście wyłącznie z palety poleceń `Ctrl/Cmd + K` (rozdz. 9.6) | Plakietka stanu (`.dn-plakietka--stan`, ikona `⛨`) w pasku kontekstu z rozwinięciem | aktywny, rozwinięty, wygasający | Rozwinięcie prezentuje zasięg, zakres i czas do wygaśnięcia oraz akcje „Wyjdź z trybu” i „Rejestr wejść” |
| 58 | Widok migawek i importu konfiguracji | Eksport, import, migawki i przywracanie stanu konfiguracji | 3 | Menu kebab (⋮) w pasku kontekstu okna | Lista pozycji z akcjami, panel różnic przed zatwierdzeniem | domyślny, ładowanie, różnice, konflikt wersji schematu | Import prezentuje różnice; konflikt wersji schematu wywołuje ostrzeżenie i mapowanie pól |
| 59 | Toast potwierdzenia | Krótki komunikat potwierdzający zapis, skopiowanie lub błąd | 1 | Pojawia się po akcji | Powiadomienie (`.dn-toast`), warianty sukces i błąd | pojawienie się, zanikanie | Nieinteraktywny, znika automatycznie lub po kliknięciu „x” |

### 13.1. Warstwy widoczności w oknie

Okno Konfiguracji stosuje regułę stopniowego ujawniania funkcjonalności do samego siebie: w stanie spoczynku widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3.

| Warstwa | Elementy okna Konfiguracji | Sposób wywołania |
|---|---|---|
| 1 | Nawigacja zakresów, selektor warstwy nakładki, edytor Konstytucji z przyciskiem zapisu i hashem, podglądy Profil/Rola i Ekspertyza zadaniowa, nagłówek i przycisk „Podgląd na żywo” Panelu prowenancji, toast potwierdzenia | Widoczne bez interakcji |
| 2 | Selektor profilu eksperta, pola treści nakładki, pola wywołania (harness, hooki, narzędzia, model, effort, kanał, pula kont), przełącznik zapisu archiwalnego, bloki Panelu prowenancji, skrót do okna izolacji, objaśnienia `[?]` | Znacznik kontekstowy, menu progresywne `▼`, nagłówek sekcji `▸`, najechanie kursorem |
| 3 | Menu profilu eksperta, podgląd łańcucha dziedziczenia, parametry próbkowania, katalog kanałów modelu, model zapasowy, grupy narzędzi, podgląd efektywnego zestawu narzędzi, biblioteka hooków, rejestr kluczy, edytor zmiennych środowiskowych, oś czasu i porównanie prowenancji, eksport i migawki, macierz przypisania warstw | Menu kebab (⋮), menu hamburger (☰), panel popover, panel wysuwany |
| 4 | Podgląd porównawczy modeli, rejestr zdarzeń kont, wyszukiwarka funkcji i ustawień, katalog skrótów klawiszowych, tryb administracyjny | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

Znaczniki kontekstowe okna — środowisko, projekt, warstwa zasięgu, profil eksperta, model, kanał — występują w pasku kontekstu jako lekkie elementy w postaci `[Danaco Console] [Ubuntu] [Sprawa X] [Roman] [Code CLI]`; kliknięcie znacznika otwiera odpowiedni selektor. Każda funkcja warstw 2–4 pozostaje osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego.

---

## 14. Stany

### 14.1. Stany okna (globalne)

| Stan | Wyzwalacz | Zachowanie interfejsu |
|---|---|---|
| Domyślny (przegląd) | Otwarcie okna lub zakresu | Wyświetla bieżące, efektywne wartości dla ostatnio wybranej warstwy zasięgu |
| Edycja pola | Fokus na polu edytowalnym | Pole aktywne (obrys `--dn-focus`); pozostałe elementy bez zmian |
| Zapisywanie | Wysłanie `config.set` po edycji pola tekstowego (Konstytucja) | Przycisk „Zapisz” pokazuje stan ładowania; pole tymczasowo nieedytowalne |
| Zapisano | Otrzymanie `config.changed` z potwierdzeniem | Toast sukcesu; wartości bieżące i hash odświeżone |
| Błąd zapisu | Serwer odrzucił zmianę albo utracono połączenie w trakcie zapisu | Toast błędu z treścią przyczyny; pole wraca do ostatniej znanej wartości, zawartość wpisana przez Operatora pozostaje w polu do ponowienia |
| Zmiana z innego urządzenia | Odebranie `config.changed` niewywołanego lokalnie | Widok odświeża się na żywo; nienachalny toast „Zaktualizowano z innego urządzenia” — bez blokowania bieżącej edycji innego pola |
| Rozłączony | Utrata połączenia WebSocket | Komunikat „Połączenie przywracane…” w pasku kontekstu; pola pozostają klikalne i edytowalne — zmiany wprowadzone w tym stanie są buforowane lokalnie i wysyłane po przywróceniu połączenia; sesja po stronie serwera pozostaje aktywna (Architektura, rozdz. 5) |
| Ładowanie zakresu | Przełączenie pozycji nawigacji lub warstwy zasięgu | Kolumna zawartości pokazuje szkielet ładowania, zastąpiony treścią po odebraniu danych |
| Tryb administracyjny aktywny | Wejście w tryb administracyjny z palety poleceń `Ctrl/Cmd + K` po potwierdzeniu uprawnień (rozdz. 9.6) | Odsłonięte funkcje warstwy 4; znacznik trybu widoczny w pasku kontekstu przez cały czas jego działania |

### 14.2. Stany pojedynczego pola / ustawienia

| Stan | Realizacja wizualna (System wizualny, rozdz. 6) | Zastosowanie |
|---|---|---|
| Domyślny | Bez wyróżnienia | Wartość dziedziczona lub wartość wyjściowa |
| Najechanie (hover) | Subtelna zmiana tła `--dn-hover` | Wszystkie elementy interaktywne |
| Fokus | Pierścień złoty 2 px, wyzwalany przez `:focus-visible` | Wszystkie kontrolki interaktywne, obowiązkowo |
| Aktywny / zaznaczony | Akcent złoty — tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt` | Pozycja nawigacji, pozycja selektora warstwy, zakładka |
| Kontekstowo nieadekwatny | Etykieta lub komunikat pomocniczy obok pola (`.dn-tekst-pomocniczy`) albo po interakcji (toast); pole pozostaje w pełni klikalne | Pole, którego wartość nie ma zastosowania w bieżącym kontekście — Harness dla Code CLI pokazuje podgląd tylko do odczytu, z objaśnieniem `[?]` |
| Ustawiony jawnie na tej warstwie | Pełne wypełnienie przełącznika, pogrubiona etykieta wartości, plakietka „ustawione na: <warstwa>” | Pola warstwowe rozdziałów 2–3 i 8 |
| Dziedziczony (nieustawiony) | Stan przełącznika przejrzysty, z plakietką „dziedziczone z: <warstwa>” | Wszystkie pola warstwowe |
| Niezapisany | Znacznik „niezapisane” przy polu tekstowym; nawigacja pozostaje możliwa | Pola tekstowe zapisywane przyciskiem |
| Ładowanie | Wskaźnik ładowania w miejscu wartości | Lista modeli po zmianie kanału, blok Panelu prowenancji podczas przeliczania |
| Ostrzeżenie | Plakietka ostrzegawcza z treścią przyczyny; wartość zapisuje się | Przekroczenie budżetu tokenów, konflikt skrótu klawiszowego, naruszenie reguły jednego kliknięcia, różnica wersji schematu przy imporcie |
| Błąd | Obrys w kolorze statusu błędu, plakietka `.dn-plakietka--blad` z treścią przyczyny | Błąd uwierzytelnienia konta puli, kanał niedostępny, błąd zapisu |
| Pusty | Komunikat pomocniczy zamiast pustego pola | Brak treści Konstytucji, brak podłączonych hooków, brak zapisanego wywołania w Panelu prowenancji |

### 14.3. Stany Panelu prowenancji

| Stan | Opis |
|---|---|
| Pusty | Żadne wywołanie nie zostało jeszcze wykonane ani podglądnięte dla wybranego zasięgu; bloki pokazują komunikat „Brak zarejestrowanego wywołania — użyj »Podgląd na żywo«” |
| Podgląd na żywo — ładowanie | Po użyciu przycisku „Podgląd na żywo”; bloki pokazują wskaźnik ładowania |
| Podgląd na żywo — gotowy | Bloki wypełnione bieżącym, przeliczonym stanem |
| Zapis archiwalny — wybrany wpis historii | Panel pokazuje zapisany komplet bloków powiązany z wybraną wiadomością z historii (5.9), z etykietą znacznika czasu |
| Porównanie dwóch wywołań | Dwie kolumny zapisów zestawione obok siebie, z podświetleniem różnic pole po polu |
| Częściowy (Konstytucja bez treści) | Blok „system prompt” pomija sekcję Konstytucji bez błędu; blok hasha pokazuje „—” zamiast wartości |

### 14.4. Diagram przejść stanów (pole warstwowe)

```
                    ┌──────────────┐
                    │ DZIEDZICZONE │◄────────────────────────┐
                    │ (domyślny)   │                          │
                    └──────┬───────┘                          │
                           │ Operator zmienia wartość           │
                           ▼                                  │
                    ┌──────────────┐                          │
                    │ EDYCJA        │                          │
                    └──────┬───────┘                          │
                           │ zapis (natychmiastowy dla pól     │
                           │ selektorowych; „Zapisz” dla         │
                           │ pól tekstowych)                    │
                           ▼                                  │
                    ┌──────────────┐        błąd              │
                    │ ZAPISYWANIE   │───────────────► ┌────────────┐
                    └──────┬───────┘                  │ BŁĄD ZAPISU │
                           │ config.changed            └─────┬──────┘
                           ▼                                 │ ponowienie
                    ┌──────────────┐                          │
                    │ USTAWIONE     │                          │
                    │ JAWNIE        │──────────────────────────┘
                    └──────┬───────┘
                           │ Operator usuwa wartość na tej warstwie
                           ▼
                    (powrót do DZIEDZICZONE)
```

---

## 15. Zachowanie i komunikacja klient–serwer

Zmiany dokonane w oknie Konfiguracji przekazywane są tym samym kanałem WebSocket, który obsługuje pozostałe polecenia i zdarzenia platformy (Model konfiguracji, rozdz. 8; Architektura, rozdz. 11).

| Zakres | Polecenie (klient → serwer) | Zdarzenie (serwer → klient) |
|---|---|---|
| Nakładka — Konstytucja, Profil/Rola, Ekspertyza zadaniowa, pola treści nakładki | `config.get`, `config.set` | `config.changed` |
| Pola wywołania (harness, hooki, narzędzia, model, effort, parametry próbkowania, kanał) | `model.channel.set`, `config.set` (Integracja modeli, rozdz. 8.2) | `config.changed`, `process.status` |
| Pula kont Code CLI | `config.set` (atrybut puli w encji kanału modelu) | `config.changed` |
| Profile ekspertów — zapis, wersja, przywrócenie, wariant | `expert.set`, `expert.version.restore` | `expert.changed` |
| Kanały komunikacji — Chat Window, Execution Loop Window | `config.set` | `config.changed` |
| Warstwy widoczności — przypisania funkcji, profile ról, skróty, tryb administracyjny | `config.set` | `config.changed` |
| Dane dostępowe i zmienne środowiskowe | `config.set` (odwołanie w bazie, wartość w magazynie poza bazą) | `config.changed` |
| Panel prowenancji — podgląd na żywo | `provenance.preview` | `provenance.preview.result` |
| Panel prowenancji — zapis archiwalny, porównanie, eksport | `provenance.archive.toggle`, `provenance.compare`, `provenance.export` | `provenance.changed`, `provenance.compare.result` |
| Import, eksport i migawki konfiguracji | `config.get`, `config.set` | `config.changed` |
| Przypisanie roli / agenta (Profil/Rola) | `role.assign`, `agent.update` (Specyfikacja agentów, rozdz. 11.3) | `agent.changed` |

### 15.1. Polecenia obsługujące Panel prowenancji i profile ekspertów

| Polecenie / zdarzenie | Kierunek | Funkcja |
|---|---|---|
| `provenance.preview` | klient → serwer | Żądanie przeliczenia i zwrócenia bloków (rozdz. 10.2) dla wskazanej warstwy zasięgu, bez wykonywania faktycznego wywołania modelu |
| `provenance.preview.result` | serwer → klient | Zwraca skomponowane bloki argv, system prompt, settings, warstwy i pochodzenie oraz hashe |
| `provenance.archive.toggle` | klient → serwer | Włącza lub wyłącza automatyczny zapis archiwalny dla wskazanej warstwy zasięgu |
| `provenance.compare` | klient → serwer | Żądanie zestawienia dwóch wpisów archiwalnych pole po polu |
| `provenance.compare.result` | serwer → klient | Zwraca zestawienie różnic obu zapisów |
| `provenance.export` | klient → serwer | Żądanie serializacji zapisu prowenancji do formatu Markdown lub JSON |
| `provenance.changed` | serwer → klient | Rozgłasza zmianę stanu zapisu archiwalnego do wszystkich połączonych urządzeń |
| `expert.set` | klient → serwer | Zapis profilu eksperta jako nowej wersji, z notatką zmiany |
| `expert.version.restore` | klient → serwer | Przywrócenie wskazanej wersji profilu albo utworzenie z niej wariantu |
| `expert.changed` | serwer → klient | Rozgłasza zmianę profilu eksperta i jego aktywnej wersji |

### 15.2. Sekwencja podglądu na żywo

```
Urządzenie A                    Serwer                          Urządzenie A
    │  provenance.preview          │                                  │
    │──────────────────────────────►│                                  │
    │                               │  kompozycja nakładki (rozdz. 2.3)│
    │                               │  rozstrzygnięcie pól (rozdz. 3)  │
    │                               │  zbudowanie argv (adapter kanału)│
    │                               │  odczyt hashy (rozdz. 10.2)      │
    │  ◄─── provenance.preview.result ──────────────────────────────── │
    │       (bloki — rozdz. 10.2)                                       │
```

### 15.3. Spójność z zasadą jednego źródła prawdy

Serwer pozostaje jedynym źródłem prawdy (Architektura, rozdz. 1.3, zasada 6) również dla nakładki, profili ekspertów, konfiguracji kanałów komunikacji, warstw widoczności i Panelu prowenancji: każda zmiana dokonana na jednym urządzeniu jest rozgłaszana zdarzeniem `config.changed` do wszystkich pozostałych połączonych urządzeń — według sekwencji ustalonej w rozdziale 8 Modelu konfiguracji.

---

## 16. Zgodność z zasadami nadrzędnymi platformy

| Zasada nadrzędna (Koncepcja platformy, rozdz. 14) | Realizacja w oknie Konfiguracji |
|---|---|
| Pełna kompozycyjność | Trzy warstwy nakładki (rozdz. 2), pola wywołania (rozdz. 3), profile ekspertów (rozdz. 4) i ustawienia obu kanałów komunikacji (rozdz. 8) łączą się swobodnie — dowolna kombinacja Konstytucji, Profilu/Roli i Ekspertyzy zadaniowej, dowolny kanał, dowolny zestaw hooków |
| Pełna konfigurowalność (zasada centralna) | Każde pole opisane w rozdziałach 2–11 podlega jawnej konfiguracji na poziomach zasięgu (rozdz. 2.4), bez wyjątku; brak ustawienia nigdy nie blokuje wywołania — przejmuje wartość dziedziczoną |
| Jawność i konfigurowalność zależności | Relacje do izolacji (6.1), integracji (6.2), rozszerzeń (6.3) i uwierzytelniania (rozdz. 7) są widoczne wprost jako skróty w interfejsie; Panel prowenancji czyni jawnym również tekst budowany przez Koordynatora, a rejestr kluczy — miejsca użycia każdego klucza |
| Stopniowe ujawnianie funkcjonalności | Warstwy widoczności (rozdz. 9) obejmują wszystkie funkcje okna; stan spoczynku prezentuje warstwę 1 i zwinięte wyzwalacze warstw 2–3, a każda funkcja pozostaje osiągalna jednym kliknięciem, skrótem klawiszowym albo poleceniem języka naturalnego |
| Dwa kanały komunikacji operacyjnej | Chat Window i Execution Loop Window mają pełne, warstwowe zestawy ustawień (rozdz. 8) i odrębne wpisy w Panelu prowenancji |
| Rozszerzenie orkiestracji | Pola wywołania obejmują wprost mechanizmy MultitaskingAI: rolę, pulę kont, równoległość zadań, politykę ponowień, progi kontroli jakości, zakres autonomii i punkty zatwierdzeń |

Zgodnie z nadrzędną zasadą braku twardych blokad (Koncepcja platformy, rozdz. 6) żaden element opisany w niniejszym dokumencie nie warunkuje ani nie wstrzymuje wywołania modelu — Konstytucja, hooki, ograniczenie narzędzi, uwierzytelnianie i izolacja techniczna są ustawieniami dokładanymi świadomie przez Operatora. Stanem wyjściowym pozostaje pełna operacyjna swoboda: model bez ustanowionej Konstytucji, bez podłączonych hooków, z pełnym dostępem do narzędzi i bez aktywnej izolacji technicznej działa w pełni sprawnie — każde zawężenie jest dodawane, nie odejmowane od stanu bazowego. Walidacja ma charakter ostrzegawczy: komunikat informuje o skutku, nie odbiera możliwości zapisu.

---

## 17. Scenariusze użycia

Scenariusze ilustrują współdziałanie mechanizmów opisanych w rozdziałach 2–15.

### 17.1. Ustanowienie Konstytucji i weryfikacja jej przez hash

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator otwiera zakres „Nakładka i wywołanie modelu”, wybiera w selektorze warstwę „Globalna” | rozdz. 12.3 |
| 2 | Wpisuje treść Konstytucji w edytor i używa przycisku „Zapisz” | rozdz. 3.1, 13 (elementy 4–5) |
| 3 | Po potwierdzeniu serwera pole hasha pokazuje nowy skrót; Operator kopiuje go do dokumentacji wewnętrznej kancelarii | rozdz. 3.1, element 6 |
| 4 | Otwierając Panel prowenancji dla dowolnej sesji, Operator porównuje wyświetlony hash z zapisanym — potwierdzenie, że sesja korzysta z zamierzonej wersji Konstytucji | rozdz. 10.2, 12.7 |

### 17.2. Rola Executor 1 z kanałem Code CLI i pulą kont dla dwóch spraw

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator konfiguruje w środowisku MultitaskingAI dwa zespoły, każdy z rolą Executor 1 pracującą nad inną sprawą | Koncepcja platformy, rozdz. 13.3 |
| 2 | Dla obu ról w polu „Kanał” wybiera Code CLI | rozdz. 3.7 |
| 3 | Otwiera widok „Pula kont Code CLI” i dodaje dwa konta: „Sprawa X” i „Sprawa Y” | rozdz. 3.8, 12.6 |
| 4 | Ustawia strategię „lepka” i przypina konto do każdej sprawy | rozdz. 3.8 |
| 5 | W zakresie „Kanały komunikacji” ustawia równoległość zadań na 2 i próg kontroli jakości dla obu ról | rozdz. 8.2 |
| 6 | Obie role pracują równolegle, każda na osobnym, odrębnie uwierzytelnionym koncie — bez wzajemnego oczekiwania i bez mieszania kontekstu spraw | rozdz. 3.8, 6.1 |

### 17.3. Audyt nietypowej odpowiedzi przez Panel prowenancji

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator zauważa nietypową odpowiedź modelu w karcie sesji modułu Developer | — |
| 2 | Z uproszczonego menu kontekstowego okna operacyjnego otwiera zapis archiwalny Panelu prowenancji dla tej wiadomości | rozdz. 10.3, 15 |
| 3 | Zestawia ten zapis z zapisem wywołania o poprawnym wyniku i odczytuje różnice pole po polu | rozdz. 10.3 |
| 4 | W bloku „system prompt” odczytuje treść sekcji Ekspertyza zadaniowa zbudowanej przez Koordynatora dla tego kroku | rozdz. 10.2, 2.2 |
| 5 | Stwierdza, że przyczyną było nieaktualne odniesienie w instrukcjach projektu; poprawia je w Instructions Panel i przekazuje zapis prowenancji do automatyki korekty | rozdz. 2.2, 10.3 |

### 17.4. Przełączenie trybu wywołania z API na Agent SDK

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Agent „Agent Backend” korzysta z kanału API | Specyfikacja agentów, Załącznik A |
| 2 | Operator stwierdza, że zadanie wymaga częstszych, drobnych wywołań narzędzi w jednej turze — zmienia „Kanał” na Agent SDK dla warstwy sesji | rozdz. 3.7, 2.4 |
| 3 | Pole „Harness” przełącza się na parametry biblioteki SDK; pole „Model” odświeża listę do modeli dostępnych tym kanałem | rozdz. 3.2, element 23 |
| 4 | Zmiana obowiązuje wyłącznie tę kartę sesji; definicja agenta i pozostałe jego zastosowania pozostają przy kanale API | rozdz. 2.4 |

### 17.5. Przełączenie profilu eksperta i weryfikacja wersji

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator w pasku kontekstu zmienia profil eksperta z „Roman” na „Ignacy” dla warstwy sesji | rozdz. 4.1, element 7 |
| 2 | Pola nakładki i pola wywołania przyjmują wartości z kapsuły profilu; plakietki warstwy wskazują pochodzenie każdej wartości | rozdz. 2.8, 4.1 |
| 3 | Operator zestawia wersję 7 z wersją 6 profilu i odczytuje różnice pole po polu | rozdz. 4.2 |
| 4 | Po wywołaniu Panel prowenancji pokazuje hash wersji eksperta, potwierdzający, którą wersją profilu obsłużono wiadomość | rozdz. 10.2 |

### 17.6. Konfiguracja pętli wykonawczej dla zlecenia wieloetapowego

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator otwiera zakres „Kanały komunikacji” i wybiera kolumnę Execution Loop Window dla roli Koordynatora | rozdz. 8.2, 12.4 |
| 2 | Ustawia równoległość zadań, limit kroków pętli i politykę ponowień z przełączeniem na model zapasowy | rozdz. 8.2, 3.5 |
| 3 | Ustala próg kontroli jakości i zachowanie przy wyniku poniżej progu na „przekaż Użytkownikowi” | rozdz. 8.2 |
| 4 | Zawęża zakres autonomii Koordynatora i ustanawia punkt zatwierdzenia na przyjęciu planu dekompozycji zlecenia | rozdz. 8.2 |
| 5 | W trakcie pracy śledzi kolejkę zadań i wskaźniki przebiegu w kolumnie Execution Loop Window, a wyniki wywołań weryfikuje w Panelu prowenancji, filtrując wpisy po kanale komunikacji | rozdz. 8.3, 10.3 |

### 17.7. Dostosowanie warstw widoczności dla roli użytkownika

| Krok | Działanie | Odniesienie |
|---|---|---|
| 1 | Operator otwiera zakres „Warstwy widoczności” i wybiera profil „Użytkownik podstawowy” | rozdz. 9.3, 12.5 |
| 2 | W macierzy przypisania przenosi wybór modelu z warstwy 2 do warstwy 3, pozostawiając go dostępnym z menu kontekstowego | rozdz. 9.2 |
| 3 | Kontrola spójności potwierdza zachowanie reguły jednego kliknięcia | rozdz. 9.2 |
| 4 | W kolumnie prawej Operator ogląda stan spoczynku interfejsu widziany przez tę rolę, bez zmiany roli bieżącej | rozdz. 9.3 |
| 5 | Funkcje warstwy 4 pozostają dostępne administratorowi przez paletę poleceń `Ctrl/Cmd + K` i tryb administracyjny | rozdz. 9.4–9.6 |

---

## 18. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Chat Window | Główne okno komunikacji w kanale Użytkownik ↔ Wykonawca, umieszczone w lewej kolumnie obszaru roboczego; konfigurowane w rozdz. 8.1 |
| Execution Loop Window | Okno pętli wykonawczej w kanale Koordynator ↔ Wykonawca, otwierane jako kolumna sąsiadująca z Chat Window; konfigurowane w rozdz. 8.2 |
| Użytkownik | Rola zlecająca i zatwierdzająca działania platformy |
| Koordynator | Komponent orkiestrujący platformy: dekomponuje zlecenie, przydziela i nadzoruje zadania, prowadzi pętlę wykonawczą |
| Wykonawca | AI, agent albo system wykonawczy realizujący zadania |
| Nakładka | Tekst instrukcji systemowej i towarzyszący mu zestaw ustawień wywołania, faktycznie przekazywany modelowi bazowemu; powstaje przez złożenie trzech warstw: Konstytucji, Profilu/Roli i Ekspertyzy zadaniowej (rozdz. 2) |
| Konstytucja | Najwyższa, najbardziej stabilna warstwa nakładki — trwałe, firmowe zasady pracy AI obowiązujące każde wywołanie modelu na platformie albo w danym środowisku; edytowana wyłącznie w oknie Konfiguracji (rozdz. 2.1–2.2) |
| Profil / Rola | Środkowa warstwa nakładki — tożsamość Wykonawcy: agent, rola środowiska MultitaskingAI albo profil asystenta; edytowana w Agent Builder lub w sekcji Role, prezentowana tutaj w podglądzie (rozdz. 2.1–2.2) |
| Ekspertyza zadaniowa | Najwęższa warstwa nakładki — wiedza proceduralna właściwa konkretnemu zadaniu lub projektowi, statyczna albo budowana dynamicznie przez Koordynatora (rozdz. 2.1–2.2) |
| Nakładka efektywna | Wynikowy, skomponowany tekst i zestaw ustawień po uwzględnieniu dziedziczenia między warstwami zasięgu, prezentowany w Panelu prowenancji (rozdz. 2.4, 10) |
| Profil eksperta | Nazwana, wersjonowana kapsuła konfiguracji nakładki i pól wywołania, odpowiadająca jednej osobowości roboczej, przełączana jednym gestem (rozdz. 4) |
| Harness | Warstwa oprogramowania spinająca wywołanie modelu w pętlę agentową — model, narzędzie, wynik, model — aż do odpowiedzi końcowej; jej pochodzenie zależy od wybranego kanału (rozdz. 3.2) |
| Hook | Punkt zaczepienia w cyklu wywołania modelu, w którym Operator podłącza dodatkowe działanie — najczęściej automatykę modułu Automations — bez warunkowania przebiegu wywołania (rozdz. 3.3) |
| Narzędzia dozwolone i zabronione | Dwie listy rozszerzeń i wbudowanych akcji rozstrzygające, z czego model korzysta w toku danego wywołania (rozdz. 3.4) |
| Effort | Poziom wysiłku obliczeniowego — głębokości rozumowania — angażowanego przez model przed wygenerowaniem odpowiedzi (rozdz. 3.6) |
| Tryb wywołania | Doprecyzowanie kanału integracji modelu (Integracja modeli, rozdz. 3) do jednej z trzech wartości: Code CLI, Agent SDK, API (rozdz. 3.7) |
| Pula kont Code CLI | Zbiór kilku odrębnie uwierzytelnionych kont narzędzia Code CLI, między którymi platforma rozdziela wywołania sesji lub ról według wybranej strategii (rozdz. 3.8) |
| Warstwa widoczności | Jeden z czterech poziomów ujawniania funkcji interfejsu, od zawsze widocznej po ekspercką; każdy element interfejsu należy do dokładnie jednej warstwy (rozdz. 9) |
| Tryb administracyjny | Stan sesji odsłaniający pełen zestaw funkcji warstwy 4, wywoływany z palety poleceń `Ctrl/Cmd + K` po potwierdzeniu uprawnień i oznaczony znacznikiem w pasku kontekstu (rozdz. 9.6) |
| Panel prowenancji | Panel dający dokładny, weryfikowalny podgląd tego, co faktycznie zostało — albo zostałoby — przekazane modelowi przy danym wywołaniu: argv, system prompt, settings, warstwy i pochodzenie, hashe (rozdz. 10) |
| argv | Blok Panelu prowenancji zawierający dosłowny zapis parametrów technicznie przekazanych przy wywołaniu, z danymi dostępowymi domyślnie w postaci jawnej (rozdz. 10.2) |
| Hash konstytucji | Skrót kryptograficzny (SHA-256) treści Konstytucji obowiązującej w chwili wywołania, pozwalający zweryfikować jej wersję bez porównywania pełnego tekstu (rozdz. 10.2) |
| Hash wersji eksperta | Skrót kryptograficzny kapsuły wersji profilu eksperta obowiązującej w chwili wywołania (rozdz. 4.2, 10.2) |
| Zapis archiwalny | Automatyczny zapis kompletu bloków Panelu prowenancji przy każdym faktycznym wywołaniu, powiązany z wpisem historii (rozdz. 10.3) |
| Podgląd na żywo | Przeliczenie i wyświetlenie bloków Panelu prowenancji przed wysłaniem wywołania, bez jego faktycznego wykonania (rozdz. 10.3) |
| Migawka konfiguracji | Nazwany zapis pełnego stanu okna Konfiguracji z datą i notatką, przywracany jedną akcją (rozdz. 11) |

---

## Załącznik A. Zbiorczy szablon konfiguracji nakładki

Szablon redakcyjny wspierający implementację, łączący pola z rozdziałów 2–11. Wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna”.

```
NAKŁADKA I WYWOŁANIE MODELU — szablon konfiguracji
  zasięg:                    globalna | środowisko | projekt | sesja | rola | agent

  nakladka:
    konstytucja:
      treść:                 <tekst>                       # tylko: globalna, środowisko
      hash:                  <skrót SHA-256, obliczany automatycznie>
    profil_rola:
      źródło:                agent | rola_bezpośrednio
      odwołanie:             <agent z Agent Builder | definicja roli>
    ekspertyza_zadaniowa:
      statyczna:              <Instructions Panel | skille agenta>
      dynamiczna:
        budowa_przez_koordynatora: włączona | wyłączona     # tylko: rola MultitaskingAI
    pola_tresci:
      negatyw:               <tekst>
      styl_i_ton:            <tekst>
      slowniczek:            [ { termin: <termin>, znaczenie: <opis> }, … ]
      format_wyjscia:        tekst | markdown | json:<schemat>
      jezyk_odpowiedzi:      <kod języka>                   # domyślnie: język interfejsu
      budzet:
        kontekst:            <liczba tokenów>
        odpowiedz:           <liczba tokenów>
      kolejnosc_warstw:      [ konstytucja, profil_rola, ekspertyza_zadaniowa ]

  profil_eksperta:
    nazwa:                   <etykieta>
    wersja:                  <numer>
    hash:                    <skrót SHA-256 kapsuły wersji>
    harmonogram:             <reguła czasowa modułu Automations>

  pola_wywolania:
    harness:
      źródło:                wbudowany_w_narzędzie | biblioteka_SDK | rdzeń_platformy
      parametry:             <zależne od kanału — rozdz. 3.2>
    hooki:                   [ { punkt: <punkt_cyklu>, akcja: <automatyka>, kolejność: <n>, włączony: true|false }, … ]
    narzedzia:
      dozwolone:             pełny_dostęp | [ <lista rozszerzeń> ]
      zabronione:            [ <lista rozszerzeń> ]
      grupa:                 <nazwa grupy narzędzi>
    model:                   <nazwa modelu bazowego>
    model_zapasowy:          <nazwa modelu bazowego + kanał>
    effort:
      włączony:               true | false
      wartość:                szybki | standardowy | dogłębny | maksymalny
    probkowanie:
      temperatura:            <wartość> | wyłączona
      top_p:                  <wartość> | wyłączony
      kara_za_powtorzenia:    <wartość> | wyłączona
      sekwencje_stop:         [ <tekst>, … ]
      seed:                   <liczba> | wyłączony
    kanał:
      tryb:                   code_cli | agent_sdk | api
      dane_dostępowe:         odwołanie → <token | klucz>          # przechowywane poza bazą
      parametry:              <API: punkt końcowy, nagłówki, wersja | CLI: nazwa, ścieżka, argumenty
                               | SSH: host, konto, klucz | HTTP: adres>
      pula_kont:                                                   # tylko: code_cli
        - nazwa:               <etykieta konta>
          odwołanie:           odwołanie → <token>
          katalog_profilu:     <ścieżka>
          limit:               <liczba> | brak
          stan:                 dostępne | zajęte | błąd_uwierzytelnienia
        strategia:             round_robin | najmniej_obciążone | ważona | wg_limitu
                               | lepka | ręcznie_per_zasięg

  kanaly_komunikacji:
    chat_window:
      wykonawca:              <agent | rola | profil asystenta>
      strumien:               strumieniowo | po_zakonczeniu
      kroki_posrednie:        widoczne | zwiniete
      zatwierdzanie:          włączone | wyłączone
      punkty_zatwierdzen:     [ <klasa działania>, … ]
      wyjasnianie_wyniku:     włączone | wyłączone
      szerokosc_kolumny:      <jednostki siatki>
    execution_loop_window:
      limit_krokow:           <liczba>
      limit_czasu_kroku:      <sekundy>
      rownoleglosc_zadan:     <liczba> | bez_ograniczenia
      polityka_ponowien:
        liczba:               <liczba>
        odstep:               <sekundy>
        warunki:              [ blad_narzedzia, wynik_ponizej_progu, … ]
        model_zapasowy:       włączony | wyłączony
      kontrola_jakosci:
        prog:                 <wartość>
        metoda:               kontrola_koordynatora | kontrola_porownawcza | test
        ponizej_progu:        ponow | przekaz_uzytkownikowi | zaakceptuj_z_adnotacja
      zakres_autonomii:       [ dekompozycja, przydzial_zadan, ponowienie, zmiana_modelu,
                                narzedzie_zewnetrzne ]
      punkty_zatwierdzen:     [ plan_dekompozycji, koniec_kontroli_jakosci, wynik_koncowy ]
      sterowanie_przebiegiem: [ wstrzymanie, wznowienie, przerwanie, korekta_zlecenia ]
      szerokosc_kolumny:      <jednostki siatki>

  warstwy_widocznosci:
    profil_roli:              <użytkownik podstawowy | zaawansowany | administrator>
    przypisania:              [ { funkcja: <nazwa>, warstwa: 1|2|3|4, wywolanie: <sposób> }, … ]
    skroty_klawiszowe:        [ { funkcja: <nazwa>, kombinacja: <klawisze> }, … ]
    tryb_administracyjny:     włączony | wyłączony
    zasieg_trybu:             sesja | urzadzenie | rola

  dane_dostepowe:
    maskowanie:               włączone | wyłączone            # domyślnie: wyłączone
    zmienne_srodowiskowe:     [ { klucz: <KLUCZ>, wartość: <wartość | odwołanie>, warstwa: <poziom> }, … ]

  panel_prowenancji:
    zapis_archiwalny:         włączony | wyłączony            # domyślnie: włączony
    maskowanie_w_eksporcie:   włączone | wyłączone            # domyślnie: wyłączone
```

---

## Załącznik B. Katalog ikon zastosowanych w oknie

Zestaw ikon zgodny z zestawem 47 ikon systemu wizualnego (rozdz. 7.2 Systemu wizualnego); okno nie wprowadza nowych ikon.

| Ikona | Zastosowanie w tym oknie |
|---|---|
| `klodka` | Wiersz danych dostępowych w bloku argv; stan konta w widoku puli kont; rejestr kluczy |
| `kod` | Nagłówek bloków „system prompt”, „settings”, „warstwy i pochodzenie” w Panelu prowenancji |
| `dokument` | Ikona przy edytorze Konstytucji i przy podglądach Profil/Rola, Ekspertyza zadaniowa |
| `plik` | Skrót do Instructions Panel i Agent Builder; eksport i import konfiguracji |
| `ustawienia` | Ikona pola „Harness” i parametrów próbkowania |
| `globus` | Ikona pola „Kanał”, gdy tryb API; ustawienie języka odpowiedzi |
| `uruchom` | Przycisk „Podgląd na żywo”, akcja „Testuj połączenie” |
| `kopiuj` | Przyciski kopiowania hashy i bloków Panelu prowenancji |
| `plus` | Przyciski „+ dodaj punkt zaczepienia”, „+ Dodaj konto”, „+ Nowy profil” |
| `kosz` | Akcja „Usuń” w widoku puli kont i w rejestrze kluczy |
| `szukaj` | Wyszukiwarka funkcji i ustawień, filtr osi czasu prowenancji |
| `warstwy` | Selektor warstwy nakładki, macierz przypisania warstw widoczności |
| `zegar` | Oś czasu prowenancji, harmonogram profilu eksperta, wersje profilu |
| `znak zapytania` | Ikona objaśnienia kontekstowego `[?]` przy każdym polu |

---

## Załącznik C. Mapowanie pól na komunikaty WebSocket

| Pole | Polecenie zapisu | Zdarzenie potwierdzające |
|---|---|---|
| Konstytucja (treść) | `config.set` | `config.changed` |
| Pola treści nakładki (negatyw, styl i ton, słowniczek, format wyjścia, język, budżet, kolejność warstw) | `config.set` | `config.changed` |
| Profil / Rola (przypisanie) | `role.assign`, `agent.update` | `agent.changed` |
| Ekspertyza zadaniowa (statyczna) | `config.set` (Instructions Panel) | `config.changed` |
| Budowa dynamiczna przez Koordynatora | `config.set` | `config.changed` |
| Profil eksperta (zapis, wersja, wariant) | `expert.set`, `expert.version.restore` | `expert.changed` |
| Harness | `config.set` | `config.changed` |
| Hooki | `config.set` | `config.changed` |
| Narzędzia dozwolone i zabronione | `config.set` (bezpośrednio) lub dziedziczone z `agent.update` (Permissions Center) | `config.changed`, `agent.changed` |
| Model, model zapasowy | `model.channel.set` | `agent.changed`, `process.status` |
| Effort, parametry próbkowania | `config.set` | `config.changed` |
| Kanał / tryb wywołania, parametry kanału | `model.channel.set` | `agent.changed`, `process.status` |
| Pula kont Code CLI, strategia, limity | `config.set` | `config.changed` |
| Ustawienia Chat Window | `config.set` | `config.changed` |
| Ustawienia Execution Loop Window | `config.set` | `config.changed` |
| Przypisania warstw widoczności, profile ról, skróty, tryb administracyjny | `config.set` | `config.changed` |
| Dane dostępowe, zmienne środowiskowe | `config.set` | `config.changed` |
| Zapis archiwalny (Panel prowenancji) | `provenance.archive.toggle` | `provenance.changed` |
| Podgląd na żywo (Panel prowenancji) | `provenance.preview` | `provenance.preview.result` |
| Porównanie wywołań, eksport prowenancji | `provenance.compare`, `provenance.export` | `provenance.compare.result` |
| Import, eksport, migawki konfiguracji | `config.get`, `config.set` | `config.changed` |

---

*Koniec dokumentu. Danaco Console — Okno Konfiguracji, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
