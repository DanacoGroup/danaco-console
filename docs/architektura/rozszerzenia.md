# Danaco Console — Rozszerzenia

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
| **Tytuł** | Rozszerzenia |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (główny) · projektant (widoczność elementów interfejsu wnoszonych przez rozszerzenie) |
| **Przeznaczenie** | Ustala jednolity kontrakt rozszerzenia obowiązujący niezależnie od jego rodzaju i pochodzenia — wtyczek, umiejętności, konektorów i serwerów MCP — oraz cykl życia rozszerzenia od instalacji po konfigurację |
| **Zakres** | dwa źródła rozszerzeń (Danaco Plugin, Personal), cykl życia rozszerzenia, stan włączenia, udział rozszerzeń w warstwie komunikacji operacyjnej, warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie, obszary `extension`, `component` kontraktu |
| **Poza zakresem** | pełny model danych encji Rozszerzenie — [Model danych](model-danych.md) rozdz. 14; moduł Agents jako całość — [Moduł Agents](../moduly/agents.md); okno konfiguracji rozszerzeń jako widok — [Okno konfiguracji](../interfejs-uzytkownika/konfiguracja.md) |
| **Dokument nadrzędny** | [Architektura techniczna](architektura.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](koncepcja-platformy.md) · [Model danych](model-danych.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Model konfiguracji](model-konfiguracji.md) |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `extension`, `component`; struktury `Extension`, `ExtensionDetail`, `ExtensionPermission`; wyliczenia `ExtensionTrustLevel`, `ExtensionAuthKind`, `ExtensionPermissionScope`, `AppValidationSeverity`) |
| **Zasada nadrzędna** | Rozszerzenie ma jeden kontrakt niezależnie od źródła — Danaco Plugin i Personal przechodzą przez ten sam rejestr, ten sam cykl życia i te same komendy kontraktu |

Dokument opisuje warstwę rozszerzeń platformy Danaco Console: wtyczki, umiejętności, konektory oraz serwery MCP. Ustala jednolity kontrakt rozszerzenia obowiązujący niezależnie od jego rodzaju i pochodzenia, opisuje dwa źródła rozszerzeń — Danaco Plugin oraz Personal — cykl życia rozszerzenia od instalacji, przez włączenie, po konfigurację, udział rozszerzeń w warstwie komunikacji operacyjnej oraz warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie. Warstwa rozszerzeń jest elementem architektury platformy, spójnym z rozdziałem 14 Architektury, rozdziałem 14 Modelu danych, rozdziałem 17 Kontraktów komunikacji oraz rozdziałem 11.15 (moduł Agents) Koncepcji platformy.

---

## Spis treści

1. [Wprowadzenie](#wprowadzenie)
2. [Rodzaje rozszerzeń](#1-rodzaje-rozszerzeń)
3. [Źródła rozszerzeń](#2-źródła-rozszerzeń)
   - [2.1 Danaco Plugin](#21-danaco-plugin)
   - [2.2 Personal](#22-personal)
   - [2.3 Wspólny kontrakt niezależny od źródła](#23-wspólny-kontrakt-niezależny-od-źródła)
4. [Jednolity kontrakt rozszerzenia](#3-jednolity-kontrakt-rozszerzenia)
5. [Model danych rozszerzenia](#4-model-danych-rozszerzenia)
6. [Cykl życia rozszerzenia](#5-cykl-życia-rozszerzenia)
   - [5.1 Instalacja](#51-instalacja)
   - [5.2 Włączenie i wyłączenie](#52-włączenie-i-wyłączenie)
   - [5.3 Konfiguracja](#53-konfiguracja)
   - [5.4 Aktualizacja](#54-aktualizacja)
7. [Rozszerzenia w module Agents](#6-rozszerzenia-w-module-agents)
8. [Rozszerzenia w oknie konfiguracji platformy](#7-rozszerzenia-w-oknie-konfiguracji-platformy)
9. [Bezpieczeństwo i zakres działania](#8-bezpieczeństwo-i-zakres-działania)
   - [8.1 Weryfikacja podpisu i skan manifestu](#81-weryfikacja-podpisu-i-skan-manifestu)
10. [Komunikacja klient–serwer](#9-komunikacja-klientserwer)
11. [Zasada pełnej konfigurowalności rozszerzeń](#10-zasada-pełnej-konfigurowalności-rozszerzeń)
12. [Rozszerzenia w warstwie komunikacji operacyjnej](#11-rozszerzenia-w-warstwie-komunikacji-operacyjnej)
   - [11.1 Rejestracja narzędzi dostępnych Wykonawcy](#111-rejestracja-narzędzi-dostępnych-wykonawcy)
   - [11.2 Udział w pętli wykonawczej Koordynator ↔ Wykonawca](#112-udział-w-pętli-wykonawczej-koordynator--wykonawca)
   - [11.3 Wywołanie funkcji rozszerzenia poleceniem języka naturalnego](#113-wywołanie-funkcji-rozszerzenia-poleceniem-języka-naturalnego)
13. [Warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie](#12-warstwy-widoczności-elementów-interfejsu-wnoszonych-przez-rozszerzenie)
14. [Słowniczek pojęć](#13-słowniczek-pojęć)
15. [Kryteria odbioru](#14-kryteria-odbioru)
16. [Załącznik A. Tabela zestawcza rozszerzeń](#załącznik-a-tabela-zestawcza-rozszerzeń)
17. [Załącznik B. Szablony konfiguracji rozszerzenia](#załącznik-b-szablony-konfiguracji-rozszerzenia)
   - [B.1. Szablon rekordu rozszerzenia](#b1-szablon-rekordu-rozszerzenia)
   - [B.2. Szablon pola Konfiguracja swoistego rodzajowi](#b2-szablon-pola-konfiguracja-swoistego-rodzajowi)
   - [B.3. Szablon poleceń obszaru rozszerzeń](#b3-szablon-poleceń-obszaru-rozszerzeń)
18. [Załącznik C. Scenariusze użycia](#załącznik-c-scenariusze-użycia)
   - [C.1. Podłączenie personalnego serwera MCP do agenta](#c1-podłączenie-personalnego-serwera-mcp-do-agenta)
   - [C.2. Czasowe wyłączenie wbudowanej wtyczki](#c2-czasowe-wyłączenie-wbudowanej-wtyczki)
   - [C.3. Konektor wykorzystywany przez rolę w środowisku MultitaskingAI](#c3-konektor-wykorzystywany-przez-rolę-w-środowisku-multitaskingai)
   - [C.4. Wywołanie funkcji rozszerzenia poleceniem języka naturalnego](#c4-wywołanie-funkcji-rozszerzenia-poleceniem-języka-naturalnego)
   - [C.5. Instalacja rozszerzenia Personal o niskim poziomie zaufania](#c5-instalacja-rozszerzenia-personal-o-niskim-poziomie-zaufania)
19. [Załącznik D. Struktury danych pełne](#załącznik-d-struktury-danych-pełne)
   - [D.1. `Extension` — nagłówek pozycji katalogu](#d1-extension--nagłówek-pozycji-katalogu)
   - [D.2. `ExtensionDetail` — pełna metryka pozycji](#d2-extensiondetail--pełna-metryka-pozycji)
   - [D.3. `ExtensionPermission` — uprawnienie deklarowane lub nadane](#d3-extensionpermission--uprawnienie-deklarowane-lub-nadane)
   - [D.4. `ExtensionSignature` — wynik weryfikacji podpisu](#d4-extensionsignature--wynik-weryfikacji-podpisu)
   - [D.5. `ExtensionHealth` — wynik sprawdzenia kondycji integracji](#d5-extensionhealth--wynik-sprawdzenia-kondycji-integracji)
   - [D.6. Struktury pomocnicze cyklu życia i integracji](#d6-struktury-pomocnicze-cyklu-życia-i-integracji)
   - [D.7. `Component` — komponent własny](#d7-component--komponent-własny)
   - [D.8. Wyliczenia obszaru `extension` i `component`](#d8-wyliczenia-obszaru-extension-i-component)
20. [Załącznik E. Kody błędów kontraktu rozszerzeń](#załącznik-e-kody-błędów-kontraktu-rozszerzeń)
21. [Załącznik — pełny wykaz komend kontraktu warstwy rozszerzeń](#załącznik--pełny-wykaz-komend-kontraktu-warstwy-rozszerzeń)
   - [Obszar `extension` — 37 komend](#obszar-extension--37-komend)
   - [Obszar `component` — 5 komend](#obszar-component--5-komend)

---

## Wprowadzenie

Warstwa rozszerzeń obejmuje wtyczki, umiejętności, konektory i serwery MCP — cztery rodzaje funkcyjne — oraz dwie wartości drugiej grupy, opisujące konfigurację połączenia: `api` i `ssh` (rozdz. 1). Wszystkie realizują jednolity kontrakt integracji i dzielą się na dwa źródła — Danaco Plugin oraz Personal. Rdzeń serwera nie rozróżnia źródła rozszerzenia; obowiązuje wspólny kontrakt integracji (Architektura, rozdz. 14).

| Zagadnienie | Rozdział |
|---|---|
| Rodzaje rozszerzeń w dwóch grupach | 1 |
| Dwa źródła — różnice i podobieństwa | 2 |
| Treść jednolitego kontraktu | 3 |
| Model danych encji Rozszerzenie | 4 |
| Pełny cykl życia (instalacja, włączenie, konfiguracja, aktualizacja) | 5 |
| Miejsca operacyjnego wykorzystania — moduł Agents | 6 |
| Rejestracja w oknie konfiguracji | 7 |
| Udział w warstwie komunikacji operacyjnej | 11 |
| Warstwy widoczności elementów interfejsu | 12 |

Rozszerzenia nie stanowią odrębnej warstwy architektonicznej równorzędnej środowiskom, modułom czy funkcjom globalnym opisanym w rozdziale 5 Koncepcji platformy. Są zasobem — analogicznym do komponentu własnego (Koncepcja platformy, rozdz. 2.3) — rejestrowanym centralnie i wykorzystywanym operacyjnie w miejscach wskazanych w rozdziałach 6, 11 i 12. Poniższy schemat porządkuje pozycję rozszerzenia względem rejestracji i wykorzystania.

```
OKNO KONFIGURACJI — sekcja „rozszerzenia” (rejestr, rozdz. 7)
        │
        │  Decyzja 1: rejestracja — co jest dostępne na serwerze
        ▼
   ROZSZERZENIE  (zasób — analogiczny do komponentu własnego, Koncepcja rozdz. 2.3)
        │
        │  Decyzja 2: podłączenie — kto z niego korzysta (moduł Agents, rozdz. 6)
        ▼
   AGENT / MODUŁ — wykorzystanie operacyjne
        │
        ▼
   Środowiska: TalkIn · WorkSpace · CodeStudio · MultitaskingAI
```

Zgodnie z zasadą pełnej konfigurowalności (Koncepcja platformy, rozdz. 14, zasada 2) oraz z nadrzędną zasadą braku twardych blokad, żadne rozszerzenie nie jest platformie narzucone poza zestawem wbudowanym Danaco Plugin, a stan każdego rozszerzenia — włączone lub wyłączone, skonfigurowane lub domyślne — pozostaje w każdej chwili zmienialny z poziomu okna konfiguracji.

---

## 1. Rodzaje rozszerzeń

Rodzaje rozszerzeń dzielą się na dwie grupy. Grupę pierwszą — **rozszerzenia funkcyjne** — tworzą cztery rodzaje wskazane w rozdziale 14 dokumentu Architektury: serwer MCP (`mcp`), wtyczka (`plugin`), umiejętność (`skill`) i konektory (`connectors`). Każdy z nich odpowiada innemu sposobowi rozszerzania możliwości platformy i modeli, lecz wszystkie cztery podlegają temu samemu, jednolitemu kontraktowi integracji opisanemu w rozdziale 3. Grupę drugą — **rozszerzenia / konfiguracja połączenia** — tworzą wartości `api` i `ssh`, opisujące nie funkcję wniesioną do platformy, lecz sposób zestawienia połączenia z zasobem zewnętrznym; wartość `api` nie jest zatem rodzajem funkcyjnym. Poniższa tabela wymienia rodzaje funkcyjne; grupę konfiguracji połączenia omawia dalsza część rozdziału.

| Rodzaj rozszerzenia | Definicja | Zastosowanie |
|---|---|---|
| Wtyczka | Rozszerzenie dostarczające platformie lub modułowi funkcję, narzędzie lub akcję spoza zestawu bazowego. | Wnosi narzędzie dostępne agentowi i modułowi, w tym funkcje przetwarzania dokumentu. |
| Umiejętność (skill) | Zdefiniowany sposób wykonania określonego rodzaju zadania — zestaw instrukcji, wzorców postępowania i wiedzy proceduralnej przypisywany agentowi. | Prowadzi określony rodzaj analizy; agent stosuje ją przy realizacji zadań tego rodzaju. |
| Konektor | Integracja z zewnętrznym systemem lub usługą, umożliwiająca agentowi i modułowi odczyt i zapis danych poza platformą. | Łączy platformę z zewnętrznym systemem biznesowym, repozytorium i usługą sieciową. |
| Serwer MCP | Zewnętrzny serwer zgodny ze standardem Model Context Protocol, udostępniający platformie zestaw narzędzi zdefiniowanych po stronie tego serwera. | Udostępnia agentowi narzędzia zewnętrznego dostawcy realizującego protokół MCP. |

Cztery rodzaje funkcyjne dzielą się według charakteru podłączenia na dwie podgrupy, odpowiadające dwóm oknom modułu Agents (rozdz. 6).

```
ROZSZERZENIE  (jednolity kontrakt — rozdz. 3)
│
├── Rozszerza działanie agenta ─────────────► Skills Manager (rozdz. 6)
│     ├── Wtyczka       — nowa funkcja, narzędzie lub akcja
│     └── Umiejętność   — zdefiniowany sposób wykonania zadania
│
└── Otwiera dostęp poza platformę ──────────► Connectors Manager (rozdz. 6)
      ├── Konektor      — integracja dedykowana usłudze
      └── Serwer MCP    — integracja wg protokołu Model Context Protocol
```

Serwer MCP jest, na poziomie technicznym, szczególnym rodzajem konektora — integracją realizującą standaryzowany protokół zamiast integracji dedykowanej pojedynczej usłudze — lecz Architektura wymienia go jako rodzaj odrębny ze względu na odmienny sposób odkrywania i udostępniania narzędzi. Poniższa tabela zestawia obie integracje w cechach, w których się różnią.

| Cecha rozróżniająca | Konektor | Serwer MCP |
|---|---|---|
| Charakter integracji | Integracja dedykowana pojedynczej usłudze. | Integracja realizująca standaryzowany protokół Model Context Protocol. |
| Definicja zestawu narzędzi | Po stronie autora konektora. | Po stronie samego serwera MCP. |
| Sposób odkrywania narzędzi | Ustalony przez konektor. | Udostępniany przez serwer zgodnie z protokołem. |
| Okno podłączenia w module Agents | Connectors Manager. | Connectors Manager. |

Rozróżnienie to jest zachowane w dalszej części dokumentu, w tym w mapowaniu na okna modułu Agents (rozdz. 6).

**Rozszerzenia / konfiguracja połączenia.** Druga grupa liczy dwie wartości. `api` opisuje połączenie z usługą zewnętrzną przez interfejs programowy dostawcy, `ssh` — połączenie z zasobem zdalnym kanałem SSH. Obie wskazują sposób zestawienia połączenia, nie funkcję wniesioną do platformy, dlatego nie występują ani w tabeli rodzajów funkcyjnych powyżej, ani w przypisaniu do okien modułu Agents (rozdz. 6); rozszerzenie takiej postaci konfiguruje się z okna konfiguracji platformy (rozdz. 7), a jego parametry mieszczą się w atrybucie Konfiguracja (rozdz. 4). Wyliczenie `ExtensionKind` (Załącznik D.8) niesie obie wartości jako postać docelową. Plik `budowa/shared/contract.json` zna dziś wartość `api`; **wdrożenie wartości `ssh` w tym pliku jest osobnym zadaniem w kodzie** — kod jej nie zna.

---

## 2. Źródła rozszerzeń

Zgodnie z rozdziałem 14 Architektury oraz rozdziałem 14 dokumentu Model danych, każde rozszerzenie — niezależnie od rodzaju wskazanego w rozdziale 1 — pochodzi z jednego z dwóch źródeł. Poniższy schemat przedstawia drogę każdego źródła do wspólnego rejestru rozszerzeń.

```
Danaco Plugin                                   Personal
─────────────                                   ────────
pakiet aplikacyjny serwera              urządzenie Operatora
        │ start serwera                         │ wskazanie pliku w oknie konfiguracji
        ▼                                       ▼ przesłanie (WebSocket)
rejestracja automatyczna                katalog użytkownika na serwerze
        │                                       │ extension.install
        └───────────────► REJESTR ROZSZERZEŃ ◄──┘
                          (rozdz. 7 — oba źródła łącznie)
```

### 2.1. Danaco Plugin

Zestaw wbudowany, dostarczany razem z pakietem aplikacyjnym serwera. Rozszerzenia tego źródła są obecne na serwerze od chwili instalacji lub aktualizacji pakietu serwera, nie wymagają odrębnej instalacji przez Operatora i rejestrują się automatycznie przy starcie serwera. Aktualizacja następuje wraz z aktualizacją pakietu serwera, a nie odrębną operacją z okna konfiguracji.

### 2.2. Personal

Rozszerzenia instalowane przez Operatora z urządzenia. Plik lub pakiet wskazywany jest w oknie konfiguracji na urządzeniu klienckim, przesyłany kanałem WebSocket do serwera (Architektura, rozdz. 11) i zapisywany w katalogu użytkownika na serwerze (Architektura, rozdz. 14). Zgodnie z zasadą jednego źródła prawdy (Architektura, rozdz. 1.3, zasada 6) urządzenie klienckie nie przechowuje trwałej kopii rozszerzenia po zakończeniu przesyłania — stan trwały pozostaje wyłącznie po stronie serwera (Architektura, rozdz. 10).

### 2.3. Wspólny kontrakt niezależny od źródła

Rdzeń serwera nie rozróżnia pochodzenia rozszerzenia przy jego rejestracji, włączaniu ani wykorzystaniu operacyjnym — obowiązuje wspólny kontrakt integracji opisany w rozdziale 3. Poniższa tabela zestawia oba źródła w zakresie cech, w których rzeczywiście się różnią.

| Cecha | Danaco Plugin | Personal |
|---|---|---|
| Pochodzenie | Dostarczane z pakietem aplikacyjnym serwera. | Instalowane przez Operatora z urządzenia. |
| Sposób dostarczenia | Wbudowane, obecne od instalacji lub aktualizacji serwera. | Przesyłane kanałem WebSocket, zapisywane w katalogu użytkownika. |
| Lokalizacja przechowywania | Katalog pakietu aplikacyjnego serwera. | Katalog użytkownika na serwerze. |
| Kopia na urządzeniu klienckim | Nie dotyczy — rozszerzenie jest wbudowane na serwerze. | Brak trwałej kopii po przesłaniu; stan trwały wyłącznie na serwerze (Architektura, rozdz. 10). |
| Aktualizacja | Wraz z aktualizacją pakietu serwera. | Ponowne przesłanie zaktualizowanej wersji przez Operatora. |
| Rejestracja encji Rozszerzenie | Automatyczna, przy starcie serwera. | Po zakończeniu przesyłania z urządzenia (`extension.install`). |
| Widoczność w rejestrze | Zawsze widoczne w rejestrze rozszerzeń. | Widoczne po zakończonej instalacji. |
| Stan wyjściowy | Włączone (rozdz. 5.2). | Wyłączone (rozdz. 5.2). |
| Kontrakt integracji, stan włączenia, konfiguracja | Wspólne, zgodnie z rozdziałem 3. | Wspólne, zgodnie z rozdziałem 3. |

---

## 3. Jednolity kontrakt rozszerzenia

Kontrakt rozszerzenia jest zbiorem cech i zachowań, jakie musi spełniać każde rozszerzenie — niezależnie od rodzaju (rozdz. 1) i źródła (rozdz. 2) — aby rdzeń mógł je zarejestrować, włączyć i wykorzystać operacyjnie. Jednolitość kontraktu jest bezpośrednią konsekwencją zasady rozdzielenia przez kontrakty (Architektura, rozdz. 1.3, zasada 4): moduły i rdzeń łączą się z rozszerzeniem przez zdefiniowany interfejs integracji, a nie przez zależność właściwą konkretnemu rodzajowi lub źródłu.

Kontrakt obejmuje trzy grupy elementów, opisane szczegółowo w rozdziałach 4, 5 i 8.

| Grupa elementów kontraktu | Treść | Opisana w |
|---|---|---|
| Tożsamość i rejestracja | Każde rozszerzenie posiada identyfikator, nazwę i wskazane źródło, niezależnie od sposobu trafienia do rejestru. | rozdz. 4 |
| Stan i sterowanie | Każde rozszerzenie posiada stan włączenia, zmienialny w dowolnym momencie, oraz konfigurację właściwą jego rodzajowi. | rozdz. 5 |
| Zakres działania | Każde rozszerzenie działa w zakresie uprawnień przypisanym przy podłączeniu do agenta lub modułu, zgodnie z zasadami izolacji technicznej. | rozdz. 8 |

Poniższy diagram przedstawia kontrakt jako wspólny interfejs, przez który rdzeń, moduły i agenci obcują z rozszerzeniem niezależnie od jego rodzaju i źródła.

```
                          REJESTR ROZSZERZEŃ (rozdz. 7)
                                   │
   ┌───────────────┬───────────────┼───────────────┬────────────────┐
   │ Wtyczka       │ Umiejętność   │ Konektor      │ Serwer MCP     │  ← rodzaj (rozdz. 1)
   │ Danaco Plugin │ Personal      │ Danaco Plugin │ Personal       │  ← źródło (rozdz. 2)
   └───────┬───────┴───────┬───────┴───────┬───────┴───────┬────────┘
           └───────────────┴───────┬───────┴───────────────┘
                                   ▼
              ╔════════════════════════════════════════╗
              ║   JEDNOLITY KONTRAKT INTEGRACJI        ║
              ║   (rozdz. 3)                           ║
              ║   • Tożsamość i rejestracja  (rozdz. 4)║
              ║   • Stan i sterowanie        (rozdz. 5)║
              ║   • Zakres działania         (rozdz. 8)║
              ╚════════════════════╤═══════════════════╝
                                   │  interfejs integracji (Architektura, zasada 4)
                  ┌────────────────┴────────────────┐
                  ▼                                 ▼
           RDZEŃ SERWERA                     MODUŁY / AGENCI
      rejestruje · włącza ·             korzystają przez kontrakt,
      wykorzystuje operacyjnie          nie przez zależność właściwą
      (nie rozróżnia rodzaju            konkretnemu rodzajowi lub źródłu
       ani źródła)
```

Ponieważ kontrakt jest wspólny, wymiana jednego rozszerzenia na inne tego samego rodzaju — na przykład zastąpienie konektora Personal odpowiadającym mu konektorem z zestawu Danaco Plugin — nie wymaga zmian po stronie modułu ani agenta korzystającego z rozszerzenia; zmienia się jedynie wskazanie rozszerzenia w konfiguracji.

---

## 4. Model danych rozszerzenia

Rozdział 14 dokumentu Model danych definiuje encję Rozszerzenie sześcioma atrybutami. Poniższa tabela przenosi tę definicję na poziom opisu produktowego, wskazując znaczenie każdego atrybutu w kontekście niniejszego dokumentu. Szablon redakcyjny rekordu przedstawia Załącznik B.

| Atrybut | Znaczenie | Wartości / zakres |
|---|---|---|
| Identyfikator | Jednoznaczne odwołanie do rozszerzenia w rejestrze, niezależne od źródła i rodzaju. | Wartość jednoznaczna w rejestrze serwera. |
| Nazwa | Nazwa własna rozszerzenia, prezentowana Operatorowi w rejestrze oraz przy podłączaniu do agenta. | Nazwa własna. |
| Źródło | Pochodzenie rozszerzenia zgodnie z rozdziałem 2. | Danaco Plugin albo Personal. |
| Stan włączenia | Wartość logiczna sterująca dostępnością rozszerzenia do wykorzystania operacyjnego; zmienialna poleceniem `extension.toggle` (rozdz. 9). | Włączone albo wyłączone. |
| Konfiguracja | Zbiór parametrów właściwych danemu rozszerzeniu, edytowalny z okna konfiguracji zgodnie z mechanizmem konfiguracji warstwowej (Model danych rozdz. 17; Architektura rozdz. 13). | Parametry swoiste rodzajowi — zob. tabela poniżej. |
| Lokalizacja | Miejsce przechowywania pliku lub pakietu rozszerzenia (rozdz. 2). | Katalog pakietu serwera (Danaco Plugin) albo katalog użytkownika (Personal). |

Atrybut Konfiguracja przechowuje wartości specyficzne dla danego egzemplarza rozszerzenia, odrębne od stanu włączenia i od zakresu uprawnień nadawanego przy podłączeniu do agenta (rozdz. 6, rozdz. 8). Zawartość tego atrybutu zależy od rodzaju rozszerzenia.

| Rodzaj rozszerzenia | Parametry przechowywane w atrybucie Konfiguracja |
|---|---|
| Wtyczka | Parametry działania wtyczki. |
| Umiejętność (skill) | Zakres zastosowania umiejętności. |
| Konektor | Dane dostępowe do zewnętrznego systemu lub usługi. |
| Serwer MCP | Adres serwera oraz dane dostępowe. |

Rozróżnienie między konfiguracją egzemplarza a zakresem uprawnień odpowiada ogólnej zasadzie warstwowości konfiguracji przyjętej dla całej platformy (Architektura, rozdz. 13; Koncepcja platformy, rozdz. 6.6).

Poniższy schemat pokazuje relacje między encją Rozszerzenie i strukturami pomocniczymi wykorzystywanymi przez komendy obszaru `extension` (Załącznik D pełne pola).

```
                         ┌─────────────────────┐
                         │      Extension      │  nagłówek pozycji katalogu
                         │  id · code · kind   │  (rozdz. 4, Załącznik D.1)
                         │  installed · enabled│
                         └──────────┬──────────┘
                                    │  1
             ┌──────────────────────┼───────────────────────┬───────────────────────┐
             │ 0..1                 │ 0..1                  │ 0..n                  │ 0..n
             ▼                      ▼                       ▼                       ▼
    ┌──────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐ ┌─────────────────────┐
    │ ExtensionDetail  │ │ ExtensionSignature  │ │ ExtensionPermission │ │ ExtensionHealth     │
    │ changelog · tools│ │ trustLevel          │ │ scope · mode        │ │ status · latencyMs  │
    │ dependencies     │ │ verified            │ │ (D.3)               │ │ (D.5)               │
    │ (D.2)            │ │ (D.4)               │ └─────────────────────┘ └─────────────────────┘
    └──────────────────┘ └─────────────────────┘
                                    │
                                    │  Component (D.7) jest encją odrębną —
                                    │  reprezentuje byt magazynu modułowego,
                                    │  nie rozszerzenie; powiązanie łączy je
                                    │  wyłącznie przez wspólny mechanizm
                                    │  komponentów własnych (Koncepcja, rozdz. 2.3)
```

Legenda: krotność `1` oznacza dokładnie jedno wystąpienie na rozszerzenie; `0..1` — najwyżej jedno; `0..n` — dowolnie wiele. `ExtensionPermission` występuje zarówno jako pole deklarowane manifestem (`ExtensionDetail.permissions`), jak i jako pole nadane (`extension.permission.grant`), zgodnie z rozróżnieniem `declared` / `granted` komendy `extension.permission.list` (Załącznik, obszar `extension`).

---

## 5. Cykl życia rozszerzenia

Cykl życia rozszerzenia obejmuje instalację, włączenie i wyłączenie, konfigurację oraz aktualizację. Przebieg poszczególnych etapów różni się między źródłami wyłącznie tam, gdzie wynika to z różnicy opisanej w rozdziale 2 — sam mechanizm sterowania stanem i konfiguracją jest, zgodnie z rozdziałem 3, wspólny.

```
   Personal: extension.install            Danaco Plugin: start serwera
            │                                        │
            ▼                                        ▼
     ┌─────────────┐                          ┌─────────────┐
     │ WYŁĄCZONE   │  ◄── stan wyjściowy      │  WŁĄCZONE   │  ◄── stan wyjściowy
     │             │      (Personal)          │             │      (Danaco Plugin)
     └──────┬──────┘                          └──────┬──────┘
            │           extension.toggle             │
            └─────────────────►  ◄───────────────────┘
                     (natychmiastowe, odwracalne — bez restartu serwera)
                                    │
                                    ▼
                 KONFIGURACJA  ·  config.get / config.set
                 (atrybut Konfiguracja encji Rozszerzenie, rozdz. 4)
                                    │
                                    ▼
                 AKTUALIZACJA
                 Danaco Plugin: wraz z pakietem serwera
                 Personal: ponowne przesłanie (ten sam identyfikator,
                           zachowane konfiguracja i stan włączenia)
```

Zestawienie przebiegu etapów w podziale na źródła.

| Etap cyklu życia | Danaco Plugin | Personal |
|---|---|---|
| Instalacja | Automatyczna przy starcie serwera; nie jest operacją Operatora. | Cztery kroki: wskazanie pliku → przesłanie (WebSocket) → zapis i rejestracja → potwierdzenie (rozdz. 5.1). |
| Stan wyjściowy | Włączone. | Wyłączone. |
| Włączenie i wyłączenie | `extension.toggle` — natychmiastowe, odwracalne. | `extension.toggle` — natychmiastowe, odwracalne. |
| Konfiguracja | `config.get` / `config.set`; każdy parametr z objaśnieniem kontekstowym `[?]`. | `config.get` / `config.set`; każdy parametr z objaśnieniem kontekstowym `[?]`. |
| Aktualizacja | Wraz z aktualizacją pakietu serwera. | Ponowne przesłanie zaktualizowanej wersji (ten sam identyfikator, zachowane konfiguracja i stan). |

### 5.1. Instalacja

**Znaczenie instalacji zależy od rodzaju rozszerzenia** i rozstrzyga się według podziału na dwie grupy z rozdziału 1:

| Grupa rodzajów | Co robi instalacja | Co robi odinstalowanie |
|---|---|---|
| Rozszerzenia funkcyjne pobierane jako paczka: wtyczka (`plugin`), umiejętność (`skill`) | Pobiera paczkę składników i zapisuje ją na dysku serwera | Usuwa pobraną paczkę składników z dysku |
| Rozszerzenia zestawiające połączenie: serwer MCP (`mcp`), konektory (`connectors`), `api`, `ssh` | Rejestruje połączenie — adres i punkt dostępu; nic się nie ściąga | Usuwa zapis połączenia, nie kasując niczego po stronie usługi zdalnej |

Kroki opisane niżej dotyczą rodzajów pobieranych jako paczka. Dla rodzajów zestawiających połączenie krok przesłania pliku nie występuje — Operator podaje w oknie konfiguracji adres i punkt dostępu, a serwer zapisuje rekord encji Rozszerzenie (rozdz. 4) z tymi parametrami w atrybucie Konfiguracja. Komendy realizujące oba warianty to `extension.install` i `extension.uninstall` ([Moduł Apps](../moduly/apps.md) rozdz. 10.2).

**Danaco Plugin.** Instalacja nie jest operacją wykonywaną przez Operatora — rozszerzenie jest obecne od chwili instalacji lub aktualizacji pakietu aplikacyjnego serwera i rejestruje się automatycznie przy starcie serwera.

**Personal.** Instalacja przebiega w czterech krokach.

1. Operator wskazuje plik lub pakiet rozszerzenia na urządzeniu, z poziomu okna konfiguracji.
2. Urządzenie klienckie przesyła plik kanałem WebSocket do serwera.
3. Serwer zapisuje plik w katalogu użytkownika i tworzy odpowiadający mu rekord encji Rozszerzenie (rozdz. 4), ze źródłem Personal i wskazaną lokalizacją.
4. Serwer potwierdza zakończenie instalacji; rozszerzenie staje się widoczne w rejestrze rozszerzeń okna konfiguracji (rozdz. 7).

```
Urządzenie Operatora                 Serwer
│  wybór pliku rozszerzenia
│────── przesłanie (WebSocket) ──────►
│                                     │  zapis w katalogu użytkownika
│                                     │  rejestracja encji Rozszerzenie
│                                     │  (źródło: Personal)
│◄───── potwierdzenie instalacji ─────
│                                     │
▼                                     ▼
rozszerzenie widoczne w rejestrze okna konfiguracji (rozdz. 7)
```

### 5.2. Włączenie i wyłączenie

Każde rozszerzenie — Danaco Plugin lub Personal — posiada stan włączenia, zmieniany poleceniem `extension.toggle` (rozdz. 9) z poziomu okna konfiguracji. Zmiana stanu jest natychmiastowa i nie wymaga ponownego uruchomienia serwera ani odłączenia sesji korzystających z innych rozszerzeń. Zgodnie z zasadą braku twardych blokad wyłączenie rozszerzenia nie usuwa go z rejestru ani z katalogu — jest odwracalną decyzją Operatora, nie operacją niszczącą.

| Źródło | Stan wyjściowy | Uzasadnienie |
|---|---|---|
| Danaco Plugin | Włączone | Zestaw wbudowany jest częścią bazowej funkcjonalności platformy i pozostaje dostępny od chwili instalacji serwera, chyba że Operator świadomie go wyłączy. |
| Personal (nowo zainstalowane) | Wyłączone | Operator włącza je świadomie, po przejrzeniu konfiguracji (rozdz. 5.3) i zakresu działania (rozdz. 8), zanim rozszerzenie zostanie wykorzystane operacyjnie. |

Oba stany wyjściowe pozostają w pełni zmienialne z okna konfiguracji, zgodnie z zasadą „brak ustawienia = wartość domyślna” (Koncepcja platformy, rozdz. 6.5).

### 5.3. Konfiguracja

Parametry właściwe danemu rozszerzeniu — atrybut Konfiguracja opisany w rozdziale 4 — edytowane są z poziomu okna konfiguracji, w miejscu rejestru rozszerzeń (rozdz. 7). Każdy parametr opatrzony jest objaśnieniem kontekstowym (`[?]`), zgodnie z ogólną zasadą przyjętą dla okna konfiguracji (Architektura, rozdz. 13): treść objaśnienia opisuje działanie parametru i jego wpływ na aplikację i stanowi część definicji parametru.

### 5.4. Aktualizacja

Aktualizacja rozszerzenia Danaco Plugin następuje wraz z aktualizacją pakietu aplikacyjnego serwera i nie jest odrębną operacją Operatora. Aktualizacja rozszerzenia Personal przebiega przez ponowne przesłanie zaktualizowanej wersji z urządzenia, tą samą drogą co pierwotna instalacja (rozdz. 5.1) — nowa wersja zastępuje poprzednią pod tym samym identyfikatorem, z zachowaniem dotychczasowej konfiguracji i stanu włączenia.

**Zakres poleceń kontraktu komunikacji.** Kontrakt komunikacji (rozdz. 9; Kontrakty komunikacji, rozdz. 17) wymienia trzy polecenia podstawowe obszaru rozszerzeń: `extension.list`, `extension.install` i `extension.toggle`; pełny wykaz trzydziestu siedmiu komend obszaru `extension` zawiera Załącznik z wykazem komend kontraktu warstwy rozszerzeń. Czasowe wycofanie rozszerzenia z użycia realizowane jest zmianą stanu włączenia poleceniem `extension.toggle`; operacja ta jest odwracalna i nie usuwa rekordu z rejestru. Trwałe usunięcie rekordu z rejestru realizuje odrębna komenda `extension.uninstall`, a zapis parametrów rozszerzenia — komenda `extension.configure`.

---

## 6. Rozszerzenia w module Agents

Moduł Agents (Koncepcja platformy, rozdz. 11.15) jest udokumentowanym miejscem operacyjnego wykorzystania rozszerzeń. Przy budowie agenta w oknie Agent Builder Operator dobiera umiejętności w oknie Skills Manager, podłącza integracje zewnętrzne przez okno Connectors Manager, a zakres uprawnień agenta ustala w oknie Permissions Center.

Cztery rodzaje funkcyjne z rozdziału 1 przypisują się do dwóch okien modułu Agents zgodnie z ich charakterem; wartości grupy konfiguracji połączenia (`api`, `ssh`) nie mają przypisania do okna modułu Agents.

| Okno modułu Agents | Rodzaje rozszerzeń | Charakter podłączenia |
|---|---|---|
| Skills Manager | Umiejętność (skill), wtyczka | Rozszerzenie działania agenta o zdefiniowany sposób wykonania zadania lub o nową funkcję dostępną w toku pracy agenta. |
| Connectors Manager | Konektor, serwer MCP | Podłączenie agenta do zewnętrznego systemu, usługi lub serwera zgodnego z protokołem Model Context Protocol. |

Poniższy schemat przedstawia rozmieszczenie okien modułu Agents istotnych dla rozszerzeń. Pełny zestaw okien modułu obejmuje ponadto Chat Window i Model Configuration (Koncepcja platformy, rozdz. 11.15), niezwiązane bezpośrednio z rozszerzeniami.

```
MODUŁ AGENTS ─ okno Agent Builder (Koncepcja platformy, rozdz. 11.15)
┌────────────────────┬────────────────────┬─────────────────────┐
│ Skills Manager     │ Connectors Manager │ Permissions Center  │
│                    │                    │                     │
│ • Umiejętność      │ • Konektor         │ zakres uprawnień    │
│ • Wtyczka          │ • Serwer MCP       │ agenta do zasobów   │
│                    │                    │ (rozdz. 8)          │
│ zmienia sposób     │ otwiera dostęp     │ dostęp sieciowy,    │
│ działania agenta   │ do świata poza     │ odczyt i zapis      │
│                    │ platformą          │ plików              │
└─────────┬──────────┴─────────┬──────────┴──────────┬──────────┘
          └────────────────────┴─────────────────────┘
                    wspólny rejestr rozszerzeń
                 i jednolity kontrakt (rozdz. 3)
```

Rozdzielenie na dwa okna odzwierciedla różnicę charakteru między rozszerzeniami zmieniającymi sposób działania agenta (Skills Manager) a rozszerzeniami otwierającymi agentowi dostęp do świata poza platformą (Connectors Manager) — nie stanowi odrębnego kontraktu: oba okna operują na tym samym rejestrze rozszerzeń i tym samym jednolitym kontrakcie opisanym w rozdziale 3.

Po utworzeniu agent jest, zgodnie z rozdziałem 11.15 Koncepcji platformy, komponentem platformowym dostępnym jako wykonawca zadań w środowiskach TalkIn, WorkSpace i CodeStudio, a po przypisaniu do roli — także w środowisku MultitaskingAI (Koncepcja platformy, rozdz. 9.4, rozdz. 13). Rozszerzenia podłączone do agenta pozostają z nim związane niezależnie od tego, w którym środowisku lub w jakiej roli agent jest w danej chwili wykorzystywany — zgodnie z zasadą, że agent jest niezależny od środowiska (Model danych, rozdz. 12).

---

## 7. Rozszerzenia w oknie konfiguracji platformy

Okno konfiguracji (Architektura, rozdz. 13) udostępnia pełny zakres ustawień wpływających na aplikację, w tym osobno wymienione sekcje „rozszerzenia” oraz „integracje”. Rozróżnienie to nie jest przypadkowe — obie sekcje sąsiadują w tym samym oknie, lecz operują na odrębnych encjach.

| Sekcja okna konfiguracji | Encja | Zakres |
|---|---|---|
| „integracje” | Kanał modelu (Model danych, rozdz. 11.1) | Sposób połączenia z samym modelem językowym: API, CLI, SSH, HTTP (Architektura, rozdz. 9). |
| „rozszerzenia” | Rozszerzenie (Model danych, rozdz. 14) | Wtyczki, umiejętności, konektory i serwery MCP opisane w niniejszym dokumencie. |

```
OKNO KONFIGURACJI (Architektura, rozdz. 13)
├── sekcja „integracje”  → encja Kanał modelu (Model danych, rozdz. 11.1)
│        API · CLI · SSH · HTTP   (połączenie z modelem — Architektura, rozdz. 9)
│
└── sekcja „rozszerzenia” → encja Rozszerzenie (Model danych, rozdz. 14)  ── REJESTR
         wtyczki · umiejętności · konektory · serwery MCP
         ├── extension.list          → wykaz rozszerzeń obu źródeł
         ├── extension.install       → instalacja rozszerzenia Personal
         ├── extension.toggle        → zmiana stanu włączenia
         └── config.get / config.set → edycja atrybutu Konfiguracja
```

Sekcja rozszerzeń pełni funkcję rejestru: prezentuje wszystkie rozszerzenia zarejestrowane na serwerze — Danaco Plugin i Personal łącznie, zgodnie z poleceniem `extension.list` (rozdz. 9) — niezależnie od tego, czy dane rozszerzenie jest w danej chwili podłączone do jakiegokolwiek agenta. Rejestracja rozszerzenia i jego wykorzystanie operacyjne są, zgodnie z zasadą jawności i konfigurowalności zależności (Koncepcja platformy, rozdz. 14, zasada 3), dwiema odrębnymi, świadomymi decyzjami Operatora.

| Decyzja | Miejsce | Co ustala |
|---|---|---|
| Rejestracja rozszerzenia | Okno konfiguracji — sekcja „rozszerzenia” (rozdz. 7) | Czy rozszerzenie istnieje w rejestrze serwera; jego stan włączenia i konfigurację. |
| Podłączenie do agenta | Moduł Agents — Skills Manager lub Connectors Manager (rozdz. 6) | Który agent korzysta z rozszerzenia i w jakim zakresie uprawnień. |

Z poziomu rejestru Operator instaluje nowe rozszerzenia Personal (rozdz. 5.1), zmienia stan włączenia (rozdz. 5.2) oraz edytuje konfigurację (rozdz. 5.3). Podłączenie zarejestrowanego rozszerzenia do konkretnego agenta następuje odrębnie, w module Agents (rozdz. 6).

---

## 8. Bezpieczeństwo i zakres działania

Instalacja i włączenie rozszerzenia — obu źródeł — nie podlega żadnej twardej blokadzie: platforma nie ogranicza z góry liczby ani rodzaju rozszerzeń Personal, jakie Operator może zainstalować, zgodnie z nadrzędną zasadą braku blokad wbudowanych na stałe (Koncepcja platformy, rozdz. 6). Kontrola nad zakresem działania rozszerzenia odbywa się nie przez blokowanie instalacji, lecz przez dwa mechanizmy skonfigurowalne z poziomu okna konfiguracji.

| Mechanizm kontroli | Miejsce konfiguracji | Działanie |
|---|---|---|
| Stan wyjściowy wyłączony dla rozszerzeń Personal | Okno konfiguracji (rozdz. 5.2) | Nowo zainstalowane rozszerzenie nie jest od razu wykorzystywane operacyjnie; Operator świadomie je włącza. |
| Zakres uprawnień nadawany przy podłączeniu do agenta | Permissions Center, moduł Agents (rozdz. 6) | Ustala, do jakich zasobów agent korzystający z danego rozszerzenia ma dostęp. |

```
Instalacja / włączenie rozszerzenia — BRAK TWARDEJ BLOKADY (Koncepcja, rozdz. 6)
        │
        ▼   kontrola nie przez blokadę instalacji, lecz przez dwa mechanizmy:
   ┌──────────────────────────────┬───────────────────────────────────┐
   │ 1. Stan wyjściowy WYŁĄCZONE  │ 2. Zakres uprawnień przy          │
   │    dla rozszerzeń Personal   │    podłączeniu do agenta          │
   │    (rozdz. 5.2)              │    — Permissions Center (rozdz. 6)│
   │                              │                                   │
   │ Operator świadomie włącza    │ dostęp sieciowy, odczyt i zapis   │
   │                              │ plików (Architektura, rozdz. 12)  │
   └──────────────────────────────┴───────────────────────────────────┘
        │
        ▼   profil izolacji roli w środowisku MultitaskingAI
            (Koncepcja, rozdz. 6.5, poziom zasięgu „Rola”)
   obejmuje pośrednio rozszerzenia podłączone do agenta pełniącego tę rolę,
   bez odrębnej konfiguracji izolacji per rozszerzenie
```

Zakres uprawnień nadawany w Permissions Center pozostaje spójny z zakresami izolacji technicznej zdefiniowanymi w rozdziale 12 Architektury i rozwiniętymi w oknie konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4–6.6) — w szczególności z zakresami dostępu sieciowego oraz odczytu i zapisu plików, istotnymi dla konektorów i serwerów MCP łączących agenta z systemami poza platformą. Profil izolacji przypisany do roli agenta w środowisku MultitaskingAI obejmuje w ten sposób pośrednio również rozszerzenia podłączone do agenta pełniącego tę rolę, bez potrzeby odrębnej konfiguracji izolacji per rozszerzenie.

### 8.1. Weryfikacja podpisu i skan manifestu

Trzeci mechanizm, odrębny od stanu wyjściowego i od zakresu uprawnień, dostarcza Operatorowi informację o wiarygodności pozycji katalogu przed jej wykorzystaniem operacyjnym — bez wstrzymywania instalacji ani włączenia. Mechanizm składa się z trzech komend: `extension.signature.verify` sprawdza podpis cyfrowy i sumę kontrolną paczki oraz ustala poziom zaufania wydawcy (wyliczenie `ExtensionTrustLevel`: `danacoPlugin`, `verifiedPublisher`, `unverifiedPersonal`); `extension.manifest.scan` przegląda manifest pod kątem szerokich uprawnień, nieznanego wydawcy i braku podpisu, zwracając wykaz spostrzeżeń o wadze `info`, `warning` albo `error` (wyliczenie `AppValidationSeverity`); `extension.audit.list` zestawia, które rozszerzenia mają jakie uprawnienia oraz kiedy i przez którego agenta były użyte, tworząc trwały ślad wykorzystania operacyjnego.

Poniższy diagram sekwencji przedstawia przebieg weryfikacji wywołany z rejestru rozszerzeń okna konfiguracji po zakończeniu instalacji.

```
Operator              Powłoka                            Rdzeń
   │                     │                                 │
   │  otwiera kartę      │                                 │
   │  pozycji katalogu   │                                 │
   ├────────────────────►│                                 │
   │                     │  extension.detail.get           │
   │                     ├────────────────────────────────►│
   │                     │  detail:ExtensionDetail         │
   │                     │◄────────────────────────────────┤
   │                     │  extension.signature.verify     │
   │                     ├────────────────────────────────►│
   │                     │  signature:ExtensionSignature   │
   │                     │◄────────────────────────────────┤
   │                     │  extension.manifest.scan        │
   │                     ├────────────────────────────────►│
   │                     │  findings:ExtensionScanFinding[]│
   │                     │◄────────────────────────────────┤
   │  karta z poziomem   │                                 │
   │  zaufania i wykazem │                                 │
   │  spostrzeżeń        │                                 │
   │◄────────────────────┤                                 │
   │                     │                                 │
   │  decyzja: włącza mimo ostrzeżenia (zasada zero blokad)│
   ├────────────────────►│  extension.toggle               │
   │                     ├────────────────────────────────►│
   │                     │  extension:Extension            │
   │                     │◄────────────────────────────────┤
```

Legenda: `extension.detail.get`, `extension.signature.verify`, `extension.manifest.scan`, `extension.toggle` — komendy obszaru `extension`; żadna z odpowiedzi nie niesie pola blokującego wykonanie kolejnego kroku — wynik `error` w `findings` pozostaje spostrzeżeniem, zgodnie z zasadą zero blokad (rozdz. 10). Dziennik wykorzystania rozszerzenia po podłączeniu do agenta czyta się komendą `extension.audit.list`, poza przebiegiem tej karty.

---

## 9. Komunikacja klient–serwer

Rozdział 17 dokumentu Kontrakty komunikacji wymienia trzy polecenia podstawowe obszaru rozszerzeń, wymieniane przez klienta i serwer kanałem WebSocket (Architektura, rozdz. 11). Pełny wykaz trzydziestu siedmiu komend obszaru `extension` oraz pięciu komend obszaru `component` zawiera Załącznik z wykazem komend kontraktu warstwy rozszerzeń.

| Polecenie | Kierunek | Opis |
|---|---|---|
| `extension.list` | klient → serwer | Wykaz rozszerzeń zarejestrowanych na serwerze, obu źródeł (Danaco Plugin i Personal). |
| `extension.install` | klient → serwer | Instalacja rozszerzenia Personal przesłanego z urządzenia (rozdz. 5.1). |
| `extension.toggle` | klient → serwer | Włączenie lub wyłączenie rozszerzenia, niezależnie od źródła (rozdz. 5.2). |

Konfiguracja parametrów rozszerzenia (rozdz. 5.3) korzysta z ogólnego mechanizmu konfiguracji platformy — poleceń `config.get` i `config.set` (Kontrakty komunikacji, rozdz. 14) — stosowanego do atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4); obszar `extension` udostępnia ponadto własną komendę zapisu kompletu parametrów.

| Polecenie | Zastosowanie do rozszerzeń |
|---|---|
| `config.get` | Odczyt atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4). |
| `config.set` | Zapis atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4). |
| `extension.configure` | Zapis kompletu parametrów rozszerzenia jednym wywołaniem obszaru `extension` (rozdz. 5.3). |

Zmiana stanu rozszerzenia — instalacja, włączenie, wyłączenie lub zmiana konfiguracji — jest, zgodnie z ogólną zasadą komunikacji (Kontrakty komunikacji, rozdz. 1), rozgłaszana przez serwer do wszystkich połączonych urządzeń, tak aby rejestr rozszerzeń pozostawał spójny niezależnie od urządzenia, z którego Operator w danej chwili korzysta.

```
Urządzenie A            Serwer                       Urządzenie B … N
    │ extension.toggle      │                                │
    │──────────────────────►│                                │
    │                       │ zmiana stanu encji Rozszerzenie│
    │◄──── potwierdzenie ───│                                │
    │                       │─────── rozgłoszenie zmiany ───►│
    │                       │   (Kontrakty komunikacji, rozdz. 1)
    ▼                       ▼                                ▼
        rejestr rozszerzeń spójny na wszystkich urządzeniach
```

---

## 10. Zasada pełnej konfigurowalności rozszerzeń

Warstwa rozszerzeń podlega w całości czterem zasadom nadrzędnym funkcjonalności platformy (Koncepcja platformy, rozdz. 14).

| Zasada nadrzędna | Zastosowanie do warstwy rozszerzeń | Przykład |
|---|---|---|
| Pełna kompozycyjność | Rozszerzenia podłączone do agenta łączą się swobodnie z pozostałymi komponentami kompozycji w środowisku MultitaskingAI (Koncepcja platformy, rozdz. 13.3; Model danych, rozdz. 10.2). | Umiejętność lub konektor podłączony do agenta pełniącego rolę Executor pozostaje dostępny również wtedy, gdy Executor uruchamia Subagent Network. |
| Pełna konfigurowalność | Stan włączenia, konfiguracja i zakres uprawnień każdego rozszerzenia są w każdej chwili zmienialne z okna konfiguracji (rozdz. 5, rozdz. 7). | Bez wyjątków wynikających z rodzaju lub źródła rozszerzenia. |
| Jawność i konfigurowalność zależności | Podłączenie rozszerzenia do agenta jest jawną, widoczną w module Agents decyzją Operatora (rozdz. 6). | Nie jest zależnością wbudowaną na stałe w agenta lub w moduł. |
| Rozszerzenie orkiestracji | Rejestr rozszerzeń przyjmuje dowolną liczbę pozycji obu źródeł bez ograniczenia wbudowanego (Koncepcja platformy, rozdz. 14). | Serwery MCP i konektory obejmują zakres współpracy agentów środowiska MultitaskingAI z systemami zewnętrznymi. |

Konsekwencją tych zasad jest brak jakiejkolwiek kombinacji rodzaju, źródła, modułu docelowego lub roli, dla której platforma z góry wykluczałaby zastosowanie rozszerzenia — ograniczenia wynikają wyłącznie ze świadomej konfiguracji zakresu uprawnień przez Operatora (rozdz. 8), nigdy z reguły wbudowanej na stałe w platformę.

---

## 11. Rozszerzenia w warstwie komunikacji operacyjnej

Warstwa komunikacji operacyjnej platformy obejmuje dwa kanały: Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca). Rozszerzenie uczestniczy w obu kanałach — rejestruje narzędzia dostępne Wykonawcy, wnosi typy zadań do pętli wykonawczej i podlega kontroli jakości wyniku prowadzonej przez Koordynatora, a jego funkcje są wywoływalne poleceniem języka naturalnego z głównego okna komunikacji.

### 11.1. Rejestracja narzędzi dostępnych Wykonawcy

Włączone rozszerzenie podłączone do agenta rejestruje w rdzeniu zestaw narzędzi, którymi Wykonawca dysponuje przy realizacji zadania. Zestaw jest właściwy rodzajowi rozszerzenia (rozdz. 1) i mieści się w zakresie uprawnień nadanym przy podłączeniu (rozdz. 8).

| Rodzaj rozszerzenia | Narzędzia rejestrowane Wykonawcy |
|---|---|
| Wtyczka | Funkcje, narzędzia i akcje wnoszone przez wtyczkę do zestawu operacyjnego Wykonawcy. |
| Umiejętność (skill) | Procedury wykonania zadania określonego rodzaju, stosowane przez Wykonawcę jako wzorzec postępowania. |
| Konektor | Operacje odczytu i zapisu danych w zewnętrznym systemie lub usłudze. |
| Serwer MCP | Narzędzia odkrywane po stronie serwera zgodnie z protokołem Model Context Protocol. |

Rejestracja narzędzi następuje z chwilą włączenia rozszerzenia i podłączenia go do agenta; wyłączenie rozszerzenia (rozdz. 5.2) wycofuje jego narzędzia z zestawu dostępnego Wykonawcy, bez usuwania rekordu z rejestru.

### 11.2. Udział w pętli wykonawczej Koordynator ↔ Wykonawca

Narzędzia wniesione przez rozszerzenie stanowią podstawę typów zadań, na jakie Koordynator dekomponuje zlecenie. Przebieg tych zadań, ich kolejka, stan oraz wynik kontroli jakości prezentowane są w oknie Execution Loop Window.

| Element pętli wykonawczej | Udział rozszerzenia |
|---|---|
| Dekompozycja zlecenia | Koordynator uwzględnia typy zadań odpowiadające narzędziom zarejestrowanym przez włączone rozszerzenia. |
| Przydział zadania | Zadanie kierowane jest do Wykonawcy dysponującego narzędziem wymaganym przez ten typ zadania. |
| Wymiana komunikatów sterujących | Wywołanie narzędzia rozszerzenia i jego wynik przechodzą przez komunikaty pętli wykonawczej. |
| Kontrola jakości wyniku | Koordynator ocenia wynik zadania zrealizowanego narzędziem rozszerzenia i podejmuje decyzję o przyjęciu albo ponowieniu. |
| Wskaźniki przebiegu | Stan zadań realizowanych narzędziami rozszerzeń prezentowany jest wraz z pozostałymi zadaniami pętli. |

```
Chat Window (Użytkownik ↔ Wykonawca)
   │ polecenie języka naturalnego
   ▼
Koordynator ── dekompozycja na zadania ──► Execution Loop Window
   │                                        (Koordynator ↔ Wykonawca)
   │ przydział zadania wg wymaganego narzędzia
   ▼
Wykonawca ──► narzędzia zarejestrowane przez rozszerzenie (rozdz. 11.1)
   │ wynik zadania
   ▼
Koordynator ── kontrola jakości ── przyjęcie albo ponowienie
   │
   ▼
Chat Window — prezentacja wyniku i wyjaśnienie kontekstu
```

### 11.3. Wywołanie funkcji rozszerzenia poleceniem języka naturalnego

Funkcje rozszerzenia wywoływane są poleceniem języka naturalnego kierowanym z głównego okna komunikacji Użytkownik ↔ Wykonawca. Użytkownik nie wskazuje rozszerzenia z listy ani nie otwiera odrębnego okna — dobór narzędzia następuje po stronie Koordynatora i Wykonawcy na podstawie treści polecenia i zestawu narzędzi zarejestrowanych przez włączone rozszerzenia. Zatwierdzanie i przerywanie działań realizowanych narzędziami rozszerzeń odbywa się w tym samym oknie, w tym samym trybie co pozostałe działania platformy.

Polecenie, którego treść nie została dopasowana do żadnej funkcji zarejestrowanej przez włączone rozszerzenie, kończy się kodem `command_not_understood` — dziewiątym kodem zamkniętego zbioru kodów kontraktu (Załącznik E; Kontrakty komunikacji, rozdz. 7.2 i Załącznik D), odrębnym od `validation_failed`, aby interfejs odpowiadał na niezrozumiane polecenie inaczej niż na błąd walidacji. Plik `budowa/shared/contract.json` zna osiem kodów — **dopisanie dziewiątego jest osobnym zadaniem w kodzie**.

---

## 12. Warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie

Rozszerzenie wnoszące element interfejsu deklaruje dla każdego takiego elementu warstwę widoczności (1–4) oraz sposób jego wywołania. Deklaracja jest częścią jednolitego kontraktu rozszerzenia (rozdz. 3) i podlega tym samym regułom stopniowego ujawniania funkcjonalności, co elementy własne platformy.

Nadrzędna zasada interfejsu obowiązuje rozszerzenia bez wyjątku: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**.

| Warstwa | Zawartość wnoszona przez rozszerzenie | Sposób wywołania |
|---|---|---|
| 1 | Wskaźnik stanu wykonania zadania realizowanego narzędziem rozszerzenia, prezentowany w oknie komunikacji i w pętli wykonawczej. | Widoczny bez interakcji, w ramach elementów warstwy 1 należących do platformy. |
| 2 | Wybór wariantu działania rozszerzenia, wskazanie zasobu zewnętrznego, parametry przebiegu. | Znacznik kontekstowy, ikona lub przełącznik; po użyciu element zwija się samoczynnie. |
| 3 | Zestawy akcji rozszerzenia, ustawienia szybkie. | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana. |
| 4 | Operacje zaawansowane, narzędzia diagnostyczne rozszerzenia, tryby administracyjne. | Polecenie języka naturalnego w oknie komunikacji, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. |

Reguły wiążące dla elementów interfejsu wnoszonych przez rozszerzenia:

- **Stała liczba elementów w stanie spoczynku.** Elementy rozszerzeń nie zwiększają liczby elementów widocznych w stanie spoczynku interfejsu. Rejestracja i włączenie dowolnej liczby rozszerzeń pozostawia stan spoczynku niezmieniony; elementy warstw 2–4 ujawniają się wyłącznie w kontekście zadania, które ich wymaga.
- **Zasada jednego kliknięcia.** Każdy element wniesiony przez rozszerzenie jest osiągalny jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji rozszerzenia głęboko w hierarchii menu jest zabronione.
- **Rozmieszczenie kolumnowe.** Panel wnoszony przez rozszerzenie otwiera się jako okno pomocnicze w kolumnie bocznej, po prawej stronie obszaru roboczego, zgodnie z układem pionowym (podział lewa–prawa). Po zamknięciu panel znika całkowicie z przestrzeni roboczej.
- **Deklaracja w rejestrze.** Warstwa widoczności i sposób wywołania każdego elementu są częścią opisu rozszerzenia prezentowanego w rejestrze okna konfiguracji (rozdz. 7), obok stanu włączenia i konfiguracji.

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │              [znacznik ▼]│ rozszerzenia
            │                      │                       [⋮]│ (rozszerzenie
            │ ─────────────────    │                          │  boczne,
            │ Execution Loop       │                          │  warstwa 3)
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ════════════════════════════════════════════════════════════════════════════
   stan spoczynku: warstwa 1 oraz zwinięte wyzwalacze warstw 2–3
```

---

## 13. Słowniczek pojęć

| Pojęcie | Znaczenie |
|---|---|
| Rozszerzenie | Wtyczka, umiejętność, konektor lub serwer MCP, rejestrowany w oknie konfiguracji platformy i wykorzystywany operacyjnie, w szczególności przez agentów w module Agents; realizuje jednolity kontrakt integracji niezależny od rodzaju i źródła (rozdz. 3). |
| Wtyczka | Rozszerzenie dostarczające nową funkcję, narzędzie lub akcję, niedostępną w zestawie bazowym platformy lub modułu (rozdz. 1). |
| Umiejętność (skill) | Rozszerzenie definiujące sposób wykonania określonego rodzaju zadania przez agenta — zestaw instrukcji i wzorców postępowania (rozdz. 1). |
| Konektor | Rozszerzenie realizujące integrację z zewnętrznym systemem lub usługą (rozdz. 1). |
| Serwer MCP | Rozszerzenie będące zewnętrznym serwerem zgodnym ze standardem Model Context Protocol, udostępniającym zestaw narzędzi zdefiniowany po stronie tego serwera (rozdz. 1). |
| Danaco Plugin | Źródło rozszerzeń wbudowanych, dostarczanych z pakietem aplikacyjnym serwera (rozdz. 2.1). |
| Personal | Źródło rozszerzeń instalowanych przez Operatora z urządzenia, przechowywanych w katalogu użytkownika na serwerze (rozdz. 2.2). |
| Kontrakt rozszerzenia | Jednolity zbiór cech i zachowań (tożsamość i rejestracja, stan i sterowanie, zakres działania) obowiązujący każde rozszerzenie niezależnie od rodzaju i źródła (rozdz. 3). |
| Katalog użytkownika | Lokalizacja na serwerze, w której przechowywane są rozszerzenia źródła Personal (rozdz. 2.2, rozdz. 4). |
| Rejestr rozszerzeń | Sekcja okna konfiguracji prezentująca wszystkie rozszerzenia zarejestrowane na serwerze, obu źródeł łącznie (rozdz. 7). |
| Narzędzie Wykonawcy | Funkcja rejestrowana przez włączone rozszerzenie i dostępna Wykonawcy przy realizacji zadania pętli wykonawczej (rozdz. 11.1). |
| Warstwa widoczności | Jedna z czterech warstw stopniowego ujawniania funkcjonalności, deklarowana przez rozszerzenie dla każdego wnoszonego elementu interfejsu (rozdz. 12). |

---

## 14. Kryteria odbioru

Warstwa rozszerzeń jest zgodna z niniejszym opracowaniem, gdy spełnione są łącznie poniższe warunki. Każdy warunek jest sprawdzalny bez odwołania do intencji autora — wyłącznie przez odczyt kontraktu, konfiguracji albo zachowania obserwowalnego z poziomu interfejsu.

| Warunek | Sposób sprawdzenia |
|---|---|
| Każde rozszerzenie — niezależnie od rodzaju (rozdz. 1) i źródła (rozdz. 2) — realizuje jednolity kontrakt rozszerzenia (rozdz. 3) | odczyt struktury `Extension` w `budowa/shared/contract.json`: pola `id`, `code`, `name`, `kind`, `installed`, `enabled`, `origin` obecne niezależnie od wartości `kind` i `origin` |
| Wszystkie 37 komend obszaru `extension` i 5 komend obszaru `component` odpowiadają wykazowi Załącznika bez nazw wymyślonych | `grep -oE '"typ": *"(extension\|component)\.[a-zA-Z.]+"' budowa/shared/contract.json \| sort -u` zwraca dokładnie 42 pozycje zgodne z Załącznikiem |
| Stan włączenia rozszerzeń Personal po instalacji jest wyłączony; Danaco Plugin — włączony (rozdz. 5.2) | nowo zainstalowana pozycja źródła `personal` ma pole `enabled` równe `false`; pozycja źródła `danaco` — `true` |
| Zmiana stanu włączenia komendą `extension.toggle` jest natychmiastowa i odwracalna, bez usunięcia rekordu z rejestru | wywołanie `extension.toggle` nie zmienia pola `id` ani `code`; rekord pozostaje odczytywalny komendą `extension.list` po wyłączeniu |
| Brak nadania uprawnienia w Permissions Center nie blokuje instalacji ani włączenia rozszerzenia (zasada zero blokad, rozdz. 10) | instalacja i włączenie rozszerzenia bez uprzedniego wywołania `extension.permission.grant` kończy się powodzeniem |
| Wynik `extension.manifest.scan` nie wstrzymuje instalacji ani włączenia niezależnie od wagi spostrzeżenia | rozszerzenie ze spostrzeżeniem wagi `error` w `findings` pozostaje instalowalne i włączalne |
| Warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie (rozdz. 12) odpowiadają wartościom zadeklarowanym, bez wyjątków twardo zakodowanych w kliencie | przegląd komponentów `.dn-*` wnoszonych przez rozszerzenie względem warstwy zadeklarowanej w konfiguracji rozszerzenia |
| Treść poświadczenia powiązanego komendą `extension.credential.bind` nigdy nie opuszcza rdzenia | odpowiedź na `extension.credential.bind` i na `extension.secret.list` niesie wyłącznie pole `ref`, nigdy treści sekretu |
| Wszystkie dziewięć kodów błędów kontraktu (Załącznik E) jest obsługiwanych jednolicie w warstwie klienckiej rozszerzeń; `budowa/shared/contract.json` zna osiem, a dziewiąty kod `command_not_understood` jest do wdrożenia w kodzie osobnym zadaniem | próba wywołania komendy obszaru `extension` z nieistniejącym `extensionId` zwraca `not_found`, nie błąd nieobsłużony po stronie klienta |

---

## Załącznik A. Tabela zestawcza rozszerzeń

Zestawienie krzyżowe czterech rodzajów funkcyjnych (rozdz. 1) i dwóch źródeł (rozdz. 2), ze wskazaniem zawartości każdego pola, okna modułu Agents, w którym dane rozszerzenie jest podłączane do agenta (rozdz. 6), oraz warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie (rozdz. 12).

| Rodzaj | Źródło Danaco Plugin | Źródło Personal | Okno podłączenia w module Agents | Warstwa wnoszonych elementów |
|---|---|---|---|---|
| Wtyczka | Wbudowane funkcje przetwarzania dokumentu, dostarczone z pakietem serwera. | Funkcje przesłane przez Operatora, właściwe jego potrzebom. | Skills Manager | 2–3 |
| Umiejętność (skill) | Wbudowane sposoby prowadzenia analiz określonego rodzaju. | Umiejętności zdefiniowane przez Operatora dla własnych rodzajów zadań. | Skills Manager | 3–4 |
| Konektor | Wbudowane integracje z usługami dostarczanymi razem z platformą. | Integracje z zewnętrznymi systemami Operatora, przesłane z urządzenia. | Connectors Manager | 2–3 |
| Serwer MCP | Wbudowane serwery MCP dostarczone z pakietem serwera. | Serwery MCP wskazane i przesłane przez Operatora. | Connectors Manager | 2–4 |

Niezależnie od pola tabeli, w którym mieści się dane rozszerzenie, obowiązuje ten sam jednolity kontrakt integracji (rozdz. 3): rodzaj i źródło zmieniają zakres zastosowania oraz okno podłączenia, nie zmieniają zaś sposobu rejestracji, sterowania stanem ani konfiguracji.

---

## Załącznik B. Szablony konfiguracji rozszerzenia

Poniższe szablony porządkują pola opisane w rozdziałach 4, 5, 11 i 12, zgodnie z definicją encji Rozszerzenie (Model danych, rozdz. 14). Wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna” (Koncepcja platformy, rozdz. 6, rozdz. 14).

### B.1. Szablon rekordu rozszerzenia

```
Rekord rozszerzenia (encja Rozszerzenie — rozdz. 4)
  Identyfikator:    <jednoznaczne odwołanie w rejestrze serwera>
  Nazwa:            <nazwa własna prezentowana Operatorowi>
  Źródło:           Danaco Plugin | Personal
  Stan włączenia:   włączone | wyłączone
                    (stan wyjściowy: Danaco Plugin = włączone,
                     Personal = wyłączone — rozdz. 5.2)
  Konfiguracja:     { parametry właściwe rodzajowi rozszerzenia — B.2 }
  Lokalizacja:      katalog pakietu serwera (Danaco Plugin)
                    | katalog użytkownika (Personal)
  Narzędzia:        { zestaw narzędzi rejestrowanych Wykonawcy — rozdz. 11.1 }
  Typy zadań:       { typy zadań pętli wykonawczej — rozdz. 11.2 }
  Elementy UI:      { element → warstwa widoczności 1–4
                      → sposób wywołania — rozdz. 12 }
```

### B.2. Szablon pola Konfiguracja swoistego rodzajowi

Zawartość atrybutu Konfiguracja zależy od rodzaju rozszerzenia (rozdz. 4). Poniższy szablon wskazuje charakter parametrów przechowywanych dla każdego rodzaju; konkretne pola należą do definicji danego rozszerzenia i są edytowane z rejestru okna konfiguracji, z objaśnieniem kontekstowym `[?]` przy każdym parametrze (rozdz. 5.3).

```
Konfiguracja — pole swoiste rodzajowi (atrybut Konfiguracja, rozdz. 4)
  Wtyczka:
      parametry działania wtyczki
  Umiejętność (skill):
      zakres zastosowania umiejętności
  Konektor:
      dane dostępowe do zewnętrznego systemu lub usługi
  Serwer MCP:
      adres serwera
      dane dostępowe
```

### B.3. Szablon poleceń obszaru rozszerzeń

```
Polecenia obszaru rozszerzeń (Kontrakty komunikacji, rozdz. 17)
  extension.list      klient → serwer   wykaz rozszerzeń obu źródeł
  extension.install   klient → serwer   instalacja rozszerzenia Personal
  extension.toggle    klient → serwer   zmiana stanu włączenia

Konfiguracja parametrów (mechanizm ogólny, Kontrakty komunikacji, rozdz. 14)
  config.get          odczyt atrybutu Konfiguracja encji Rozszerzenie
  config.set          zapis  atrybutu Konfiguracja encji Rozszerzenie

Rozgłaszanie: każda zmiana stanu jest rozsyłana do wszystkich
              połączonych urządzeń (Kontrakty komunikacji, rozdz. 1)
```

---

## Załącznik C. Scenariusze użycia

Scenariusze ilustrują cykl życia i wykorzystanie rozszerzeń opisane w rozdziałach 5, 6, 8, 11 i 12, pokazując współdziałanie tych mechanizmów.

### C.1. Podłączenie personalnego serwera MCP do agenta

| Krok | Działanie Operatora | Odniesienie |
|---|---|---|
| 1 | Instaluje z okna konfiguracji serwer MCP źródła Personal, wskazując plik na urządzeniu. | rozdz. 5.1 |
| 2 | Po zakończonej instalacji rozszerzenie pozostaje w stanie wyłączonym. | rozdz. 5.2 |
| 3 | Przegląda konfigurację serwera MCP i włącza rozszerzenie. | rozdz. 5.2–5.3 |
| 4 | W module Agents, w oknie Connectors Manager, podłącza serwer MCP do wybranego agenta. | rozdz. 6 |
| 5 | W oknie Permissions Center ustala zakres dostępu sieciowego agenta do tego serwera. | rozdz. 8 |

### C.2. Czasowe wyłączenie wbudowanej wtyczki

| Krok | Działanie i skutek | Odniesienie |
|---|---|---|
| 1 | Operator wyłącza w rejestrze rozszerzeń wybraną wtyczkę źródła Danaco Plugin, bez usuwania jej z rejestru. | rozdz. 5.2 |
| 2 | Agenci, do których wtyczka była podłączona, tracą do niej dostęp na czas wyłączenia. | rozdz. 6 |
| 3 | Ponowne włączenie przywraca dostęp bez potrzeby ponownego podłączania wtyczki w Skills Manager. | rozdz. 5.2, rozdz. 6 |

### C.3. Konektor wykorzystywany przez rolę w środowisku MultitaskingAI

| Krok | Działanie i skutek | Odniesienie |
|---|---|---|
| 1 | Agent z podłączonym konektorem (Connectors Manager) zostaje przypisany do roli Executor 1 w środowisku MultitaskingAI. | rozdz. 6; Koncepcja platformy, rozdz. 13.3 |
| 2 | Profil izolacji przypisany tej roli obejmuje zakres dostępu sieciowego stosowany również do połączeń realizowanych przez konektor. | rozdz. 8; Koncepcja platformy, rozdz. 6.5 (poziom „Rola”) |
| 3 | Nie jest wymagana odrębna konfiguracja izolacji dla samego rozszerzenia. | rozdz. 8 |

### C.4. Wywołanie funkcji rozszerzenia poleceniem języka naturalnego

| Krok | Działanie i skutek | Odniesienie |
|---|---|---|
| 1 | Użytkownik formułuje polecenie w oknie Chat Window, bez wskazywania rozszerzenia z listy. | rozdz. 11.3 |
| 2 | Koordynator dekomponuje zlecenie na zadania i przydziela je Wykonawcy dysponującemu wymaganym narzędziem rozszerzenia. | rozdz. 11.1–11.2 |
| 3 | Przebieg zadania, wymiana komunikatów sterujących i wynik kontroli jakości prezentowane są w oknie Execution Loop Window. | rozdz. 11.2 |
| 4 | Parametry przebiegu ujawniają się jako znacznik kontekstowy warstwy 2 wyłącznie na czas wykonania zadania; stan spoczynku interfejsu pozostaje niezmieniony. | rozdz. 12 |
| 5 | Wynik i wyjaśnienie kontekstu wracają do okna Chat Window, gdzie Użytkownik zatwierdza albo przerywa dalsze działania. | rozdz. 11.3 |

### C.5. Instalacja rozszerzenia Personal o niskim poziomie zaufania

| Krok | Działanie i skutek | Odniesienie |
|---|---|---|
| 1 | Operator przesyła z urządzenia paczkę serwera MCP nieznanego wydawcy komendą `extension.package.upload`, po czym instaluje ją komendą `extension.install`. | rozdz. 5.1 |
| 2 | `extension.signature.verify` zwraca poziom zaufania `unverifiedPersonal` — paczka nie jest podpisana albo podpis nie został potwierdzony. | rozdz. 8.1 |
| 3 | `extension.manifest.scan` zwraca spostrzeżenie wagi `error` z powodu żądania szerokiego zakresu `fileWrite` bez zawężenia `target`. | rozdz. 8.1 |
| 4 | Zgodnie z zasadą zero blokad żadne z powyższych nie wstrzymuje instalacji ani włączenia; Operator podejmuje świadomą decyzję, widząc oba wyniki na karcie pozycji. | rozdz. 8, rozdz. 10 |
| 5 | Operator ogranicza zakres nadania w Permissions Center do pojedynczego katalogu zamiast zakresu żądanego przez manifest, po czym podłącza serwer do agenta. | rozdz. 6, rozdz. 8 |
| 6 | Każde użycie narzędzia tego serwera przez agenta zostaje odnotowane w `ExtensionAuditEntry`, odczytywalnym komendą `extension.audit.list`. | rozdz. 8.1 |

---

## Załącznik D. Struktury danych pełne

Wykaz przenosi z `budowa/shared/contract.json` pełne definicje pól struktur używanych przez komendy obszarów `extension` i `component` (Załącznik z pełnym wykazem komend, poniżej). Kolumna „Wymagane” podaje wartość pola `wymagane` struktury źródłowej; pole nieoznaczone jako wymagane jest opcjonalne w każdym wystąpieniu struktury.

### D.1. `Extension` — nagłówek pozycji katalogu

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator pozycji katalogu |
| `code` | `string` | tak | Kod pozycji, stały między wydaniami |
| `name` | `string` | tak | Nazwa rozszerzenia |
| `kind` | `ExtensionKind` | tak | Rodzaj rozszerzenia |
| `description` | `string` | nie | Opis rozszerzenia |
| `version` | `string` | nie | Wersja rozszerzenia |
| `installed` | `bool` | tak | Czy rozszerzenie jest zainstalowane |
| `enabled` | `bool` | tak | Czy rozszerzenie jest włączone |
| `accessPointId` | `string` | nie | Punkt dostępu obsługujący rozszerzenie |
| `config` | `json` | nie | Konfiguracja rozszerzenia |
| `origin` | `ExtensionOrigin` | tak | Źródło pochodzenia pozycji katalogu — rozstrzyga wyłącznie stan wyjściowy przy rejestracji (rozdz. 5.2); poza tym jest faktem prezentowanym Operatorowi |
| `updatedAt` | `int64` | tak | Czas ostatniej zmiany w milisekundach epoki |

### D.2. `ExtensionDetail` — pełna metryka pozycji

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `extensionId` | `string` | tak | Pozycja katalogu |
| `changelog` | `string` | nie | Dziennik zmian pozycji |
| `tools` | `ExtensionToolEntry[]` | nie | Narzędzia i zasoby udostępniane przez pozycję |
| `permissions` | `ExtensionPermission[]` | nie | Uprawnienia deklarowane w manifeście |
| `dependencies` | `string[]` | nie | Kody pozycji, których ta pozycja wymaga |
| `signature` | `ExtensionSignature` | nie | Wynik weryfikacji podpisu |
| `publisher` | `string` | nie | Wydawca pozycji |
| `homepageUrl` | `string` | nie | Adres strony pozycji |
| `tags` | `string[]` | nie | Znaczniki pozycji |

### D.3. `ExtensionPermission` — uprawnienie deklarowane lub nadane

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `scope` | `ExtensionPermissionScope` | tak | Zakres dostępu |
| `target` | `string` | nie | Byt objęty zakresem — domena, korzeń katalogu, nazwa zasobu |
| `mode` | `AccessMode` | nie | Tryb dostępu do bytu |
| `explanation` | `string` | nie | Objaśnienie uprawnienia dla znacznika kontekstowego `[?]` |
| `grantedAt` | `int64` | nie | Czas nadania w milisekundach epoki |

### D.4. `ExtensionSignature` — wynik weryfikacji podpisu

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `signed` | `bool` | tak | Czy paczka jest podpisana |
| `verified` | `bool` | tak | Czy podpis został potwierdzony |
| `trustLevel` | `ExtensionTrustLevel` | tak | Poziom zaufania wydawcy |
| `algorithm` | `string` | nie | Algorytm podpisu |
| `checksumSha256` | `string` | nie | Suma kontrolna paczki |
| `publisher` | `string` | nie | Wydawca wskazany w podpisie |
| `detail` | `string` | nie | Szczegół niepowodzenia weryfikacji; puste przy powodzeniu |

### D.5. `ExtensionHealth` — wynik sprawdzenia kondycji integracji

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `extensionId` | `string` | tak | Integracja |
| `status` | `ExtensionHealthStatus` | tak | Stan integracji |
| `latencyMs` | `int` | nie | Czas odpowiedzi w milisekundach |
| `toolCount` | `int` | nie | Liczba narzędzi udostępnianych przez integrację |
| `handshakeError` | `string` | nie | Błąd powitania; puste przy powodzeniu |
| `checkedAt` | `int64` | tak | Czas sprawdzenia w milisekundach epoki |

### D.6. Struktury pomocnicze cyklu życia i integracji

| Struktura | Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|---|
| `ExtensionRejection` | `extensionId` | `string` | nie | Pozycja odrzucona, jeśli była znana |
| `ExtensionRejection` | `code` | `string` | tak | Kod pozycji |
| `ExtensionRejection` | `reason` | `string` | tak | Powód odrzucenia |
| `ExtensionScanFinding` | `severity` | `AppValidationSeverity` | tak | Waga spostrzeżenia |
| `ExtensionScanFinding` | `code` | `string` | tak | Kod spostrzeżenia |
| `ExtensionScanFinding` | `message` | `string` | tak | Treść spostrzeżenia |
| `ExtensionScanFinding` | `permissionScope` | `ExtensionPermissionScope` | nie | Zakres uprawnienia, którego spostrzeżenie dotyczy |
| `ExtensionCollection` | `id` | `string` | tak | Identyfikator kolekcji |
| `ExtensionCollection` | `name` | `string` | tak | Nazwa kolekcji |
| `ExtensionCollection` | `extensionIds` | `string[]` | tak | Pozycje wchodzące w skład kolekcji |
| `ExtensionUpdate` | `currentVersion` | `string` | tak | Wersja zainstalowana |
| `ExtensionUpdate` | `availableVersion` | `string` | tak | Wersja dostępna w rejestrze |
| `ExtensionUpdate` | `breaking` | `bool` | nie | Czy zmiana łamie zgodność semantyczną |
| `ExtensionHistoryEntry` | `action` | `ExtensionLifecycleAction` | tak | Czynność odnotowana we wpisie |
| `ExtensionHistoryEntry` | `fromVersion` / `toVersion` | `string` | nie | Wersje przed i po czynności |
| `ExtensionUsage` | `calls` / `failures` | `int` | tak | Liczba wywołań i wywołań nieudanych |
| `ExtensionUsage` | `rateLimit` / `rateRemaining` | `int` | nie | Limit szybkości deklarowany przez usługę i pozostały budżet |
| `ExtensionWebhook` | `direction` | `ExtensionWebhookDirection` | tak | Kierunek wywołania zwrotnego |
| `ExtensionWebhook` | `secretRef` | `string` | nie | Odwołanie do sekretu podpisu HMAC |
| `SecretRef` | `ref` | `string` | tak | Klucz jawny odwołania |
| `SecretRef` | `sharedWithExtensionIds` / `sharedWithRoleIds` | `string[]` | nie | Rozszerzenia i role dopuszczone do odwołania |
| `ExtensionToolEntry` | `kind` | `ExtensionToolKind` | tak | Rodzaj wpisu — `tool`, `resource` albo `prompt` |
| `ExtensionToolEntry` | `inputSchema` | `json` | nie | Schemat wejścia w formacie JSON Schema |
| `ExtensionMapping` | `id` | `string` | tak | Identyfikator odwzorowania |
| `ExtensionMapping` | `extensionId` | `string` | tak | Integracja, której odwzorowanie dotyczy |
| `ExtensionMapping` | `name` | `string` | tak | Nazwa odwzorowania |
| `ExtensionMapping` | `rules` | `json` | tak | Reguły odwzorowania i transformacji pól |
| `ExtensionMapping` | `updatedAt` | `int64` | tak | Czas ostatniej zmiany w milisekundach epoki |
| `ProtocolFrame` | `id` | `string` | tak | Identyfikator ramki |
| `ProtocolFrame` | `direction` | `ProtocolFrameDirection` | tak | Kierunek ramki |
| `ProtocolFrame` | `method` | `string` | nie | Nazwa metody; puste dla odpowiedzi |
| `ProtocolFrame` | `correlationId` | `string` | nie | Korelacja żądania z odpowiedzią |
| `ProtocolFrame` | `payload` | `json` | nie | Treść ramki |
| `ProtocolFrame` | `errorCode` | `string` | nie | Kod błędu protokołu; puste przy powodzeniu |
| `ProtocolFrame` | `occurredAt` | `int64` | tak | Czas ramki w milisekundach epoki |
| `ExtensionAuditEntry` | `extensionId` | `string` | tak | Rozszerzenie, którego wpis dotyczy |
| `ExtensionAuditEntry` | `agentId` | `string` | nie | Agent, który rozszerzenia użył |
| `ExtensionAuditEntry` | `toolName` | `string` | nie | Nazwa użytego narzędzia |
| `ExtensionAuditEntry` | `permissions` | `ExtensionPermission[]` | nie | Uprawnienia obowiązujące w chwili użycia |
| `ExtensionAuditEntry` | `usedAt` | `int64` | tak | Czas użycia w milisekundach epoki |

### D.7. `Component` — komponent własny

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator komponentu |
| `kind` | `ComponentKind` | tak | Rodzaj komponentu |
| `name` | `string` | tak | Nazwa komponentu |
| `description` | `string` | nie | Opis komponentu |
| `targetId` | `string` | nie | Byt magazynu modułowego, który komponent reprezentuje |
| `enabled` | `bool` | tak | Czy komponent jest czynny |
| `config` | `json` | nie | Konfiguracja komponentu przekazywana adapterowi modułu |
| `createdAt` | `int64` | tak | Czas utworzenia w milisekundach epoki |
| `updatedAt` | `int64` | tak | Czas ostatniej zmiany w milisekundach epoki |

### D.8. Wyliczenia obszaru `extension` i `component`

| Wyliczenie | Wartości | Znaczenie |
|---|---|---|
| `ExtensionKind` | rozszerzenia funkcyjne: `mcp` · `plugin` · `skill` · `connectors`<br>rozszerzenia / konfiguracja połączenia: `api` · `ssh` | Serwer MCP, wtyczka, umiejętność dostępna agentowi, konektory — cztery rodzaje funkcyjne; oraz dwie wartości opisujące konfigurację połączenia: interfejs programowy usługi i kanał SSH (rozdz. 1). `budowa/shared/contract.json` zna dziś `api`; **wdrożenie wartości `ssh` jest osobnym zadaniem w kodzie** |
| `ExtensionTrustLevel` | `danacoPlugin` · `verifiedPublisher` · `unverifiedPersonal` | Zestaw wbudowany dostarczany z pakietem serwera, wydawca o potwierdzonym podpisie, pozycja Operatora bez potwierdzonego podpisu; oznaczenie, nie bramka instalacji (rozdz. 8.1) |
| `ExtensionAuthKind` | `oauth2` · `apiKey` · `token` · `basic` · `none` | Sposób uwierzytelnienia integracji: OAuth2 z przekierowaniem zgody, klucz API, token nosiciela, uwierzytelnienie podstawowe, bez uwierzytelnienia; wprowadzenie poświadczenia pozostaje po stronie Operatora |
| `ProtocolFrameDirection` | `outgoing` · `incoming` | Ramka wysłana przez platformę i ramka odebrana od serwera MCP |
| `ExtensionOrigin` | `danaco` · `personal` | Danaco Plugin (stan wyjściowy: włączone) · Personal (stan wyjściowy: wyłączone, rozdz. 5.2) |
| `ExtensionPermissionScope` | `network` · `fileRead` · `fileWrite` · `processSpawn` · `secretRead` · `modelCall` | Dostęp sieciowy, odczyt plików, zapis plików, uruchamianie procesów, odczyt referencji sekretów, wywołanie modelu |
| `ExtensionToolKind` | `tool` · `resource` · `prompt` | Narzędzie wywoływane przez `tools/call`, zasób czytany przez `resources/read`, prompt pobierany przez `prompts/get` |
| `ExtensionHealthStatus` | `unknown` · `healthy` · `degraded` · `unavailable` | Stan integracji przy sprawdzeniu kondycji; stan nierozpoznany nie wygasza kontrolki |
| `ExtensionBulkAction` | `enable` · `disable` · `uninstall` | Czynność zbiorcza na wielu pozycjach rejestru naraz |
| `ExtensionLifecycleAction` | `installed` · `updated` · `enabled` · `disabled` · `rolledBack` · `uninstalled` · `configured` | Czynność odnotowana w dzienniku cyklu życia rozszerzenia |
| `ExtensionDefinitionFormat` | `openapi3` · `graphql` | Format opisu API, z którego budowana jest integracja |
| `ExtensionWebhookDirection` | `inbound` · `outbound` | Wywołanie zwrotne przychodzące (weryfikowane podpisem HMAC) i wychodzące (wysyłające zdarzenia platformy) |
| `McpTransport` | `stdio` · `sse` · `http` · `streamableHttp` | Proces lokalny strumieniem standardowym, strumień zdarzeń po HTTP, żądanie i odpowiedź po HTTP, Streamable HTTP |
| `ComponentKind` | `automations` · `agents` · `workspace` · `assistant` | Automatyka modułu Automations, agent modułu Agents, projekt przestrzeni roboczej, profil asystenta |
| `AppValidationSeverity` | `info` · `warning` · `error` | Waga spostrzeżenia walidacji; żadna wartość nie wstrzymuje czynności (zasada zero blokad) |

---

## Załącznik E. Kody błędów kontraktu rozszerzeń

Dziewięć kodów błędów kontraktu jest wspólnych dla całego kontraktu komunikacji; poniższa tabela wskazuje ich zastosowanie do komend obszarów `extension` i `component`. Kolumna „Ponawialny” przenosi pole `retryable` źródła — kod ponawialny uzasadnia automatyczne powtórzenie zadania przez klienta, kod nieponawialny wymaga zmiany żądania albo działania Operatora. Pole `kodyBledow` w `budowa/shared/contract.json` zawiera dziś osiem pierwszych pozycji — **dopisanie dziewiątego kodu `command_not_understood` jest osobnym zadaniem w kodzie**.

| Kod | Ponawialny | Znaczenie | Zastosowanie w warstwie rozszerzeń |
|---|---|---|---|
| `validation_failed` | nie | Treść zadania niezgodna z kontraktem | Pole `kind` poza wyliczeniem `ExtensionKind`; `format` poza `ExtensionDefinitionFormat` w `extension.definition.import`; brak `code` przy `extension.install` |
| `not_found` | nie | Wskazany byt nie istnieje | `extensionId` bez odpowiadającej pozycji rejestru; `collectionId` nieznany w `extension.collection.list`; `componentId` nieznany w `component.update`, `component.delete`, `component.assign` |
| `not_authenticated` | nie | Brak uwierzytelnienia | Każda komenda obszaru wymaga tokenu dostępu ważnego dla konta właściciela (Bezpieczeństwo i uwierzytelnianie, rozdz. 8) |
| `permission_denied` | nie | Uwierzytelniony, lecz bez uprawnienia do czynności | `extension.credential.bind` bez referencji dopuszczonej dla wywołującego; `extension.secret.share` rozszerzającego zakres poza nadany |
| `conflict` | nie | Stan bytu wyklucza czynność | `extension.uninstall` rozszerzenia wymaganego przez `dependencies` innej zainstalowanej pozycji; `extension.version.rollback` do wersji nieobecnej w `ExtensionHistoryEntry` |
| `channel_unavailable` | tak | Kanał modelu niedostępny; błąd dotyczy tylko bieżącego wywołania | `extension.tool.call` kierujące do integracji, której kanał transportu (`McpTransport`) jest chwilowo nieosiągalny |
| `rate_limited` | tak | Ograniczenie tempa po stronie kanału lub rdzenia | `extension.tool.call`, `extension.sandbox.run` przy przekroczeniu `rateLimit` odnotowanego w `ExtensionUsage` |
| `internal_error` | tak | Błąd wewnętrzny rdzenia | Dowolna komenda obszaru przy awarii rdzenia niezwiązanej z treścią żądania |
| `command_not_understood` | nie | Treść polecenia języka naturalnego nie została dopasowana do żadnej funkcji | Polecenie języka naturalnego kierowane do funkcji wniesionej przez rozszerzenie (rozdz. 11.3), którego treść nie została dopasowana; kod odrębny od `validation_failed`, aby interfejs odpowiadał na niezrozumiane polecenie inaczej niż na błąd walidacji. Kod jest dziewiąty w zamkniętym zbiorze kodów kontraktu; `budowa/shared/contract.json` zna osiem kodów — **dopisanie dziewiątego jest osobnym zadaniem w kodzie**, do chwili wdrożenia sytuacja ta nie ma w kodzie odrębnego kodu odpowiedzi |

Zasada zero blokad (rozdz. 10) ogranicza zastosowanie `permission_denied` w tym obszarze: brak nadania uprawnienia (`extension.permission.grant`) nie blokuje instalacji ani włączenia rozszerzenia — kod ten dotyczy wyłącznie operacji na referencjach sekretów i poświadczeniach, gdzie kontrola dostępu jest zamierzona i jawna.

---

## Załącznik — pełny wykaz komend kontraktu warstwy rozszerzeń

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `extension` — 37 komend

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> **kluczami konfiguracji**, nie komendami kontraktu. Klucz wskazuje miejsce wartości w modelu
> konfiguracji; komendy kontraktu, którymi się go odczytuje i zapisuje, to `config.get` i
> `config.set` (obszar `config` w `budowa/shared/contract.json`).

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `extension.list` | Zwraca katalog rozszerzeń: konektory, wtyczki, serwery MCP i umiejętności | `kind:ExtensionKind` (opc)<br>`agentId:string` (opc)<br>`installedOnly:bool` (opc) | `extensions:Extension[]` (wym) |
| `extension.install` | Instaluje pozycję katalogu rozszerzeń. **[DO DECYZJI OPERATORA]** — znaczenie instalacji (pobranie paczki, zarejestrowanie adresu czy zapis punktu dostępu) nie zostało rozstrzygnięte | `code:string` (wym)<br>`kind:ExtensionKind` (wym)<br>`source:string` (opc)<br>`accessPointId:string` (opc)<br>`origin:ExtensionOrigin` (opc)<br>`config:json` (opc) | `extension:Extension` (wym) |
| `extension.configure` | Zapisuje konfigurację rozszerzenia | `extensionId:string` (wym)<br>`config:json` (wym)<br>`accessPointId:string` (opc) | `extension:Extension` (wym) |
| `extension.toggle` | Włącza albo wyłącza rozszerzenie bez odinstalowania go | `extensionId:string` (wym)<br>`enabled:bool` (wym) | `extension:Extension` (wym) |
| `extension.uninstall` | Odinstalowuje rozszerzenie. **[DO DECYZJI OPERATORA]** — znaczenie odinstalowania jest związane z nierozstrzygniętym znaczeniem instalacji | `extensionId:string` (wym) | `uninstalled:bool` (wym) |
| `extension.search` | Szuka w katalogu rozszerzeń po nazwie, opisie, kategorii, udostępnianych narzędziach i znacznikach. Dopełnia `extension.list`, który zawężał wyłącznie rodzajem i stanem zainstalowania; dopasowanie idzie także środkiem nazwy, bo pozycje niosą przedrostki źródła | `query:string` (wym)<br>`kind:ExtensionKind` (opc)<br>`origin:ExtensionOrigin` (opc)<br>`installedOnly:bool` (opc)<br>`limit:int` (opc)<br>`offset:int` (opc) | `extensions:Extension[]` (wym)<br>`total:int` (wym)<br>`suggestions:string[]` (opc) |
| `extension.detail.get` | Zwraca pełną metrykę pozycji katalogu: dziennik zmian, wykaz udostępnianych narzędzi i zasobów, wymagane uprawnienia, zależności, pochodzenie i podpis. Struktura `Extension` niesie sam nagłówek pozycji, więc karta szczegółów wymaga osobnej struktury `ExtensionDetail` | `extensionId:string` (wym) | `detail:ExtensionDetail` (wym) |
| `extension.collection.list` | Zwraca kolekcje kuratorskie — nazwane zestawy rozszerzeń instalowane i aktywowane grupowo | `collectionId:string` (opc) | `collections:ExtensionCollection[]` (wym)<br>`total:int` (wym) |
| `extension.collection.save` | Zapisuje kolekcję kuratorską; puste `collectionId` zakłada nową, podane zmienia istniejącą | `collectionId:string` (opc)<br>`name:string` (wym)<br>`description:string` (opc)<br>`colorTag:string` (opc)<br>`extensionIds:string[]` (wym) | `collection:ExtensionCollection` (wym) |
| `extension.collection.apply` | Instaluje i włącza pozycje kolekcji jednym zleceniem. Wynik jest bilansem, nie potwierdzeniem: pozycja odrzucona wraca z powodem, zamiast znikać po cichu | `collectionId:string` (wym)<br>`enable:bool` (opc) | `applied:Extension[]` (wym)<br>`rejected:ExtensionRejection[]` (wym) |
| `extension.registry.list` | Zwraca pozycje prywatnego rejestru organizacji, prezentowane w katalogu obok Danaco Plugin. Wyliczenie `ExtensionOrigin` ma dziś dwie wartości, więc miejsce rejestru organizacji w katalogu pozostaje otwarte | `query:string` (opc)<br>`kind:ExtensionKind` (opc) | `extensions:Extension[]` (wym)<br>`total:int` (wym)<br>`registryUrl:string` (opc) |
| `extension.update.check` | Sprawdza, czy rejestr zna nowszą wersję pozycji. Struktura `Extension` niesie własną wersję, ale nie wersję dostępną — porównanie realizuje osobna struktura `ExtensionUpdate` | `extensionId:string` (opc) | `updates:ExtensionUpdate[]` (wym)<br>`checkedAt:int64` (wym) |
| `extension.package.upload` | Przyjmuje paczkę rozszerzenia przesłaną z urządzenia Operatora i zapisuje ją w katalogu użytkownika. `extension.install` przyjmuje kod i deklarowane źródło, ale treści pliku kontrakt nie przenosi tą komendą | `fileName:string` (wym)<br>`contentBase64:string` (wym)<br>`checksumSha256:string` (opc) | `uploadRef:string` (wym)<br>`sizeBytes:int64` (wym) |
| `extension.version.pin` | Przypina pozycję do wskazanej wersji albo zdejmuje przypięcie | `extensionId:string` (wym)<br>`version:string` (opc) | `extension:Extension` (wym)<br>`pinnedVersion:string` (opc) |
| `extension.version.rollback` | Cofa pozycję do wcześniej zainstalowanej wersji. Wykonuje się od razu; ustawienie `extension.rollback.confirm` włącza potwierdzenie, które nie warunkuje wykonania | `extensionId:string` (wym)<br>`targetVersion:string` (wym) | `extension:Extension` (wym)<br>`rolledBackFromVersion:string` (wym) |
| `extension.bundle.install` | Instaluje i konfiguruje zestaw rozszerzeń z jednego pliku definicji, odtwarzając środowisko. Wynik jest bilansem: pozycja odrzucona wraca z powodem | `manifest:json` (wym)<br>`enable:bool` (opc) | `installed:Extension[]` (wym)<br>`rejected:ExtensionRejection[]` (wym) |
| `extension.history.list` | Zwraca chronologię zdarzeń instalacji, aktualizacji, włączeń, wyłączeń i cofnięć wersji pozycji. Rejestr niesie stan bieżący, nie drogę, która do niego doprowadziła | `extensionId:string` (opc)<br>`limit:int` (opc)<br>`since:int64` (opc) | `entries:ExtensionHistoryEntry[]` (wym)<br>`total:int` (wym) |
| `extension.admin.bulk` | Wykonuje czynność na wielu pozycjach rejestru naraz — włączenie, wyłączenie albo odinstalowanie. Wynik jest bilansem przyjętych i odrzuconych, nie pojedynczym potwierdzeniem | `extensionIds:string[]` (wym)<br>`action:ExtensionBulkAction` (wym) | `affected:Extension[]` (wym)<br>`rejected:ExtensionRejection[]` (wym) |
| `extension.tool.list` | Zwraca narzędzia, zasoby i prompty udostępniane przez serwer MCP albo integrację wraz ze schematami wejścia — odpowiednik operacji `tools/list`, `resources/list` i `prompts/list` protokołu Model Context Protocol, ujęty strukturą `ExtensionToolEntry` | `extensionId:string` (wym)<br>`kind:ExtensionToolKind` (opc)<br>`refresh:bool` (opc) | `entries:ExtensionToolEntry[]` (wym)<br>`discoveredAt:int64` (wym)<br>`protocolVersion:string` (opc) |
| `extension.tool.call` | Wykonuje próbne wywołanie narzędzia integracji i oddaje odpowiedź surową oraz sformatowaną. Wywołanie jest próbne: przebiega poza kontekstem agenta i niczego mu nie przypisuje | `extensionId:string` (wym)<br>`toolName:string` (wym)<br>`arguments:json` (opc)<br>`timeoutMs:int` (opc) | `ok:bool` (wym)<br>`raw:json` (opc)<br>`text:string` (opc)<br>`durationMs:int` (wym)<br>`errorDetail:string` (opc) |
| `extension.protocol.log.list` | Zwraca ramki JSON-RPC wymienione z serwerem MCP: zadania, odpowiedzi i notyfikacje wraz z korelacją | `extensionId:string` (wym)<br>`limit:int` (opc)<br>`since:int64` (opc) | `frames:ProtocolFrame[]` (wym)<br>`total:int` (wym) |
| `extension.sandbox.run` | Uruchamia umiejętność albo wtyczkę na przykładowym wejściu w piaskownicy izolacyjnej, bez podłączania jej do agenta produkcyjnego | `extensionId:string` (wym)<br>`input:json` (wym)<br>`timeoutMs:int` (opc) | `ok:bool` (wym)<br>`output:json` (opc)<br>`logRef:string` (opc)<br>`durationMs:int` (wym) |
| `extension.definition.import` | Buduje integrację z opisu API — OpenAPI 3 albo schematu GraphQL — wraz z wykryciem operacji | `format:ExtensionDefinitionFormat` (wym)<br>`source:string` (opc)<br>`content:string` (opc)<br>`code:string` (wym) | `extension:Extension` (wym)<br>`operations:ExtensionToolEntry[]` (wym) |
| `extension.transport.set` | Zapisuje transport serwera MCP i sprawdza go | `extensionId:string` (wym)<br>`transport:McpTransport` (wym)<br>`endpoint:string` (opc)<br>`command:string` (opc)<br>`probe:bool` (opc) | `extension:Extension` (wym)<br>`probeStatus:ExtensionHealthStatus` (opc) |
| `extension.credential.bind` | Wiąże integrację z referencją poświadczenia w warstwie sekretów. Komenda operuje na referencjach; treść poświadczenia nigdy nie opuszcza rdzenia i nie wchodzi w to zadanie | `extensionId:string` (wym)<br>`authKind:ExtensionAuthKind` (wym)<br>`credentialRef:string` (wym)<br>`scopes:string[]` (opc) | `extension:Extension` (wym)<br>`authorizationUrl:string` (opc) |
| `extension.webhook.list` | Zwraca wywołania zwrotne przychodzące i wychodzące związane z integracją | `extensionId:string` (opc)<br>`direction:ExtensionWebhookDirection` (opc) | `webhooks:ExtensionWebhook[]` (wym)<br>`total:int` (wym) |
| `extension.webhook.save` | Zapisuje wywołanie zwrotne; puste `webhookId` zakłada nowe, podane zmienia istniejące | `webhookId:string` (opc)<br>`extensionId:string` (wym)<br>`direction:ExtensionWebhookDirection` (wym)<br>`url:string` (opc)<br>`eventTypes:string[]` (opc)<br>`secretRef:string` (opc)<br>`enabled:bool` (opc) | `webhook:ExtensionWebhook` (wym) |
| `extension.mapping.save` | Zapisuje odwzorowanie pól między systemem zewnętrznym a encjami platformy wraz z transformacjami | `extensionId:string` (wym)<br>`mappingId:string` (opc)<br>`name:string` (wym)<br>`rules:json` (wym) | `mapping:ExtensionMapping` (wym)<br>`validationIssues:string[]` (opc) |
| `extension.usage.get` | Zwraca liczniki wywołań, limity szybkości i szacunkowy koszt integracji w oknie czasu | `extensionId:string` (opc)<br>`since:int64` (opc)<br>`until:int64` (opc) | `usage:ExtensionUsage[]` (wym)<br>`windowStart:int64` (wym)<br>`windowEnd:int64` (wym) |
| `extension.health.check` | Sprawdza kondycję integracji: dostępność, czas odpowiedzi, liczbę udostępnianych narzędzi i błędy powitania. Komenda `access.point.check` bada most dostępu, nie usługę za nim, więc kondycję samej integracji mierzy ta komenda | `extensionId:string` (opc) | `results:ExtensionHealth[]` (wym)<br>`checkedAt:int64` (wym) |
| `extension.secret.list` | Zwraca referencje sekretów używane przez rozszerzenia wraz z ich ważnością. Oddaje wyłącznie odwołania i metrykę; treść sekretu nie opuszcza rdzenia | `extensionId:string` (opc)<br>`expiringWithinDays:int` (opc) | `secrets:SecretRef[]` (wym)<br>`total:int` (wym) |
| `extension.secret.share` | Ustala, które rozszerzenia i role mają dostęp do wskazanej referencji sekretu | `secretRef:string` (wym)<br>`extensionIds:string[]` (opc)<br>`roleIds:string[]` (opc) | `secret:SecretRef` (wym) |
| `extension.permission.list` | Zwraca uprawnienia deklarowane przez rozszerzenie w jego manifeście wraz z objaśnieniem każdego | `extensionId:string` (wym) | `declared:ExtensionPermission[]` (wym)<br>`granted:ExtensionPermission[]` (wym)<br>`excessive:string[]` (opc) |
| `extension.permission.grant` | Nadaje rozszerzeniu zakres dostępu. Nadanie jest jedyną kontrolą; brak nadania nie wstrzymuje instalacji ani włączenia (zasada zero blokad) | `extensionId:string` (wym)<br>`permissions:ExtensionPermission[]` (wym)<br>`agentId:string` (opc) | `granted:ExtensionPermission[]` (wym) |
| `extension.signature.verify` | Sprawdza podpis cyfrowy i sumę kontrolną paczki rozszerzenia oraz ustala poziom zaufania jego wydawcy | `extensionId:string` (wym) | `signature:ExtensionSignature` (wym)<br>`checkedAt:int64` (wym) |
| `extension.manifest.scan` | Przegląda manifest rozszerzenia pod kątem szerokich uprawnień, nieznanego wydawcy i braku podpisu. Wynik jest **sygnałem**, nie bramą: nie wstrzymuje instalacji ani włączenia | `extensionId:string` (wym) | `findings:ExtensionScanFinding[]` (wym)<br>`scannedAt:int64` (wym) |
| `extension.audit.list` | Zwraca zestawienie, które rozszerzenia mają jakie uprawnienia oraz kiedy i przez którego agenta były użyte | `extensionId:string` (opc)<br>`agentId:string` (opc)<br>`since:int64` (opc)<br>`limit:int` (opc) | `entries:ExtensionAuditEntry[]` (wym)<br>`total:int` (wym) |

**Zdarzenia obszaru `extension` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `extension.changed` | Zmiana katalogu rozszerzeń | — |
| `extension.protocol.frame` | Ramka JSON-RPC wymieniona z serwerem MCP. Podgląd dziennika na żywo prowadzi okno MCP & Connector Console ([Moduł Apps](../moduly/apps.md), rozdz. 3.12); odczyt zaległych ramek realizuje komenda `extension.protocol.log.list` | — |

### Obszar `component` — 5 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `component.list` | Zwraca komponenty własne strefy 2 strony głównej | `kind:ComponentKind` (opc)<br>`includeDisabled:bool` (opc) | `components:Component[]` (wym) |
| `component.create` | Zakłada komponent własny w magazynie właściwym jego rodzajowi | `kind:ComponentKind` (wym)<br>`name:string` (wym)<br>`description:string` (opc)<br>`config:json` (opc) | `component:Component` (wym) |
| `component.update` | Zmienia komponent własny; pola pominięte zostają bez zmian | `componentId:string` (wym)<br>`name:string` (opc)<br>`description:string` (opc)<br>`enabled:bool` (opc)<br>`config:json` (opc) | `component:Component` (wym) |
| `component.delete` | Usuwa komponent własny | `componentId:string` (wym) | `deleted:bool` (wym) |
| `component.assign` | Przypisuje komponent własny do poziomu zasięgu — środowiska, projektu, karty sesji albo okna. **[DO DECYZJI OPERATORA]** — znaczenie przypisania (widoczność wyłączna czy współdzielona na wskazanym poziomie) nie zostało rozstrzygnięte | `componentId:string` (wym)<br>`scope:ConfigScope` (wym)<br>`scopeId:string` (opc) | `component:Component` (wym)<br>`assigned:bool` (wym) |

**Zdarzenia obszaru `component` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `component.changed` | Zmiana komponentu własnego | — |

Razem w wykazie: **42 komendy** z 2 obszarów kontraktu.

---

*Koniec dokumentu. Rozszerzenia — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
