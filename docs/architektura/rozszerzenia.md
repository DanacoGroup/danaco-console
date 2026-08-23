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
| **Data** | 2026-08-06 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Źródło** | Koncepcja platformy i architektura |

Dokument opisuje warstwę rozszerzeń platformy Danaco Console: wtyczki, umiejętności, konektory oraz serwery MCP. Ustala jednolity kontrakt rozszerzenia obowiązujący niezależnie od jego rodzaju i pochodzenia, opisuje dwa źródła rozszerzeń — Danaco Plugin oraz Personal — cykl życia rozszerzenia od instalacji, przez włączenie, po konfigurację, udział rozszerzeń w warstwie komunikacji operacyjnej oraz warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie. Warstwa rozszerzeń jest elementem architektury platformy, spójnym z rozdziałem 14 Architektury, rozdziałem 12 Modelu danych, rozdziałem 17 Kontraktów komunikacji oraz rozdziałem 11.15 (moduł Agents) Koncepcji platformy.

---

## Spis treści

- [Wprowadzenie](#wprowadzenie)
1. [Rodzaje rozszerzeń](#1-rodzaje-rozszerzeń)
2. [Źródła rozszerzeń](#2-źródła-rozszerzeń)
3. [Jednolity kontrakt rozszerzenia](#3-jednolity-kontrakt-rozszerzenia)
4. [Model danych rozszerzenia](#4-model-danych-rozszerzenia)
5. [Cykl życia rozszerzenia](#5-cykl-życia-rozszerzenia)
6. [Rozszerzenia w module Agents](#6-rozszerzenia-w-module-agents)
7. [Rozszerzenia w oknie konfiguracji platformy](#7-rozszerzenia-w-oknie-konfiguracji-platformy)
8. [Bezpieczeństwo i zakres działania](#8-bezpieczeństwo-i-zakres-działania)
9. [Komunikacja klient–serwer](#9-komunikacja-klient–serwer)
10. [Zasada pełnej konfigurowalności rozszerzeń](#10-zasada-pełnej-konfigurowalności-rozszerzeń)
11. [Rozszerzenia w warstwie komunikacji operacyjnej](#11-rozszerzenia-w-warstwie-komunikacji-operacyjnej)
12. [Warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie](#12-warstwy-widoczności-elementów-interfejsu-wnoszonych-przez-rozszerzenie)
13. [Słowniczek pojęć](#13-słowniczek-pojęć)
- [Załącznik A. Tabela zestawcza rozszerzeń](#załącznik-a-tabela-zestawcza-rozszerzeń)
- [Załącznik B. Szablony konfiguracji rozszerzenia](#załącznik-b-szablony-konfiguracji-rozszerzenia)
- [Załącznik C. Scenariusze użycia](#załącznik-c-scenariusze-użycia)

---

## Wprowadzenie

Warstwa rozszerzeń obejmuje wtyczki, umiejętności, konektory i serwery MCP. Wszystkie realizują jednolity kontrakt integracji i dzielą się na dwa źródła — Danaco Plugin oraz Personal. Rdzeń serwera nie rozróżnia źródła rozszerzenia; obowiązuje wspólny kontrakt integracji (Architektura, rozdz. 14).

| Zagadnienie | Rozdział |
|---|---|
| Cztery rodzaje rozszerzeń | 1 |
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
OKNO KONFIGURACJI — sekcja „rozszerzenia" (rejestr, rozdz. 7)
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

Platforma udostępnia cztery rodzaje rozszerzeń, zgodnie z rozdziałem 12 dokumentu Architektury. Każdy rodzaj odpowiada innemu sposobowi rozszerzania możliwości platformy i modeli, lecz wszystkie cztery podlegają temu samemu, jednolitemu kontraktowi integracji opisanemu w rozdziale 3.

| Rodzaj rozszerzenia | Definicja | Zastosowanie |
|---|---|---|
| Wtyczka | Rozszerzenie dostarczające platformie lub modułowi funkcję, narzędzie lub akcję spoza zestawu bazowego. | Wnosi narzędzie dostępne agentowi i modułowi, w tym funkcje przetwarzania dokumentu. |
| Umiejętność (skill) | Zdefiniowany sposób wykonania określonego rodzaju zadania — zestaw instrukcji, wzorców postępowania i wiedzy proceduralnej przypisywany agentowi. | Prowadzi określony rodzaj analizy; agent stosuje ją przy realizacji zadań tego rodzaju. |
| Konektor | Integracja z zewnętrznym systemem lub usługą, umożliwiająca agentowi i modułowi odczyt i zapis danych poza platformą. | Łączy platformę z zewnętrznym systemem biznesowym, repozytorium i usługą sieciową. |
| Serwer MCP | Zewnętrzny serwer zgodny ze standardem Model Context Protocol, udostępniający platformie zestaw narzędzi zdefiniowanych po stronie tego serwera. | Udostępnia agentowi narzędzia zewnętrznego dostawcy realizującego protokół MCP. |

Cztery rodzaje dzielą się według charakteru podłączenia na dwie grupy, odpowiadające dwóm oknom modułu Agents (rozdz. 6).

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

---

## 2. Źródła rozszerzeń

Zgodnie z rozdziałem 14 Architektury oraz rozdziałem 9 dokumentu Model danych, każde rozszerzenie — niezależnie od rodzaju wskazanego w rozdziale 1 — pochodzi z jednego z dwóch źródeł. Poniższy schemat przedstawia drogę każdego źródła do wspólnego rejestru rozszerzeń.

```
Danaco Plugin                                   Personal
─────────────                                   ────────
pakiet aplikacyjny serwera              urządzenie Operatora
        │ start serwera                         │ wskazanie pliku w oknie konfiguracji
        ▼                                        ▼ przesłanie (WebSocket)
rejestracja automatyczna                katalog użytkownika na serwerze
        │                                        │ extension.install
        └───────────────► REJESTR ROZSZERZEŃ ◄───┘
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
   │ Wtyczka       │ Umiejętność   │ Konektor      │ Serwer MCP      │  ← rodzaj (rozdz. 1)
   │ Danaco Plugin │ Personal      │ Danaco Plugin │ Personal        │  ← źródło (rozdz. 2)
   └───────┬───────┴───────┬───────┴───────┬───────┴───────┬────────┘
           └───────────────┴──────┬────────┴───────────────┘
                                   ▼
              ╔═══════════════════════════════════════╗
              ║   JEDNOLITY KONTRAKT INTEGRACJI        ║
              ║   (rozdz. 3)                           ║
              ║   • Tożsamość i rejestracja  (rozdz. 4)║
              ║   • Stan i sterowanie        (rozdz. 5)║
              ║   • Zakres działania         (rozdz. 8)║
              ╚═══════════════════╤═══════════════════╝
                                  │  interfejs integracji (Architektura, zasada 4)
                 ┌────────────────┴────────────────┐
                 ▼                                  ▼
           RDZEŃ SERWERA                     MODUŁY / AGENCI
      rejestruje · włącza ·             korzystają przez kontrakt,
      wykorzystuje operacyjnie          nie przez zależność właściwą
      (nie rozróżnia rodzaju            konkretnemu rodzajowi lub źródłu
       ani źródła)
```

Ponieważ kontrakt jest wspólny, wymiana jednego rozszerzenia na inne tego samego rodzaju — na przykład zastąpienie konektora Personal odpowiadającym mu konektorem z zestawu Danaco Plugin — nie wymaga zmian po stronie modułu ani agenta korzystającego z rozszerzenia; zmienia się jedynie wskazanie rozszerzenia w konfiguracji.

---

## 4. Model danych rozszerzenia

Rozdział 9 dokumentu Model danych definiuje encję Rozszerzenie sześcioma atrybutami. Poniższa tabela przenosi tę definicję na poziom opisu produktowego, wskazując znaczenie każdego atrybutu w kontekście niniejszego dokumentu. Szablon redakcyjny rekordu przedstawia Załącznik B.

| Atrybut | Znaczenie | Wartości / zakres |
|---|---|---|
| Identyfikator | Jednoznaczne odwołanie do rozszerzenia w rejestrze, niezależne od źródła i rodzaju. | Wartość jednoznaczna w rejestrze serwera. |
| Nazwa | Nazwa własna rozszerzenia, prezentowana Operatorowi w rejestrze oraz przy podłączaniu do agenta. | Nazwa własna. |
| Źródło | Pochodzenie rozszerzenia zgodnie z rozdziałem 2. | Danaco Plugin albo Personal. |
| Stan włączenia | Wartość logiczna sterująca dostępnością rozszerzenia do wykorzystania operacyjnego; zmienialna poleceniem `extension.toggle` (rozdz. 9). | Włączone albo wyłączone. |
| Konfiguracja | Zbiór parametrów właściwych danemu rozszerzeniu, edytowalny z okna konfiguracji zgodnie z mechanizmem konfiguracji warstwowej (Model danych rozdz. 12; Architektura rozdz. 13). | Parametry swoiste rodzajowi — zob. tabela poniżej. |
| Lokalizacja | Miejsce przechowywania pliku lub pakietu rozszerzenia (rozdz. 2). | Katalog pakietu serwera (Danaco Plugin) albo katalog użytkownika (Personal). |

Atrybut Konfiguracja przechowuje wartości specyficzne dla danego egzemplarza rozszerzenia, odrębne od stanu włączenia i od zakresu uprawnień nadawanego przy podłączeniu do agenta (rozdz. 6, rozdz. 8). Zawartość tego atrybutu zależy od rodzaju rozszerzenia.

| Rodzaj rozszerzenia | Parametry przechowywane w atrybucie Konfiguracja |
|---|---|
| Wtyczka | Parametry działania wtyczki. |
| Umiejętność (skill) | Zakres zastosowania umiejętności. |
| Konektor | Dane dostępowe do zewnętrznego systemu lub usługi. |
| Serwer MCP | Adres serwera oraz dane dostępowe. |

Rozróżnienie między konfiguracją egzemplarza a zakresem uprawnień odpowiada ogólnej zasadzie warstwowości konfiguracji przyjętej dla całej platformy (Architektura, rozdz. 13; Koncepcja platformy, rozdz. 6.6).

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
     └─────┬───────┘                          └──────┬──────┘
           │            extension.toggle             │
           └──────────────────►  ◄───────────────────┘
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

**Danaco Plugin.** Instalacja nie jest operacją wykonywaną przez Operatora — rozszerzenie jest obecne od chwili instalacji lub aktualizacji pakietu aplikacyjnego serwera i rejestruje się automatycznie przy starcie serwera.

**Personal.** Instalacja przebiega w czterech krokach.

1. Operator wskazuje plik lub pakiet rozszerzenia na urządzeniu, z poziomu okna konfiguracji.
2. Urządzenie klienckie przesyła plik kanałem WebSocket do serwera.
3. Serwer zapisuje plik w katalogu użytkownika i tworzy odpowiadający mu rekord encji Rozszerzenie (rozdz. 4), ze źródłem Personal i wskazaną lokalizacją.
4. Serwer potwierdza zakończenie instalacji; rozszerzenie staje się widoczne w rejestrze rozszerzeń okna konfiguracji (rozdz. 7).

```
Urządzenie Operatora                 Serwer
│  wybór pliku rozszerzenia
│──────── przesłanie (WebSocket) ────────►
│                                     │  zapis w katalogu użytkownika
│                                     │  rejestracja encji Rozszerzenie
│                                     │  (źródło: Personal)
│◄─────── potwierdzenie instalacji ───────
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

Oba stany wyjściowe pozostają w pełni zmienialne z okna konfiguracji, zgodnie z zasadą „brak ustawienia = wartość domyślna" (Koncepcja platformy, rozdz. 6.5).

### 5.3. Konfiguracja

Parametry właściwe danemu rozszerzeniu — atrybut Konfiguracja opisany w rozdziale 4 — edytowane są z poziomu okna konfiguracji, w miejscu rejestru rozszerzeń (rozdz. 7). Każdy parametr opatrzony jest objaśnieniem kontekstowym (`[?]`), zgodnie z ogólną zasadą przyjętą dla okna konfiguracji (Architektura, rozdz. 13): treść objaśnienia opisuje działanie parametru i jego wpływ na aplikację i stanowi część definicji parametru.

### 5.4. Aktualizacja

Aktualizacja rozszerzenia Danaco Plugin następuje wraz z aktualizacją pakietu aplikacyjnego serwera i nie jest odrębną operacją Operatora. Aktualizacja rozszerzenia Personal przebiega przez ponowne przesłanie zaktualizowanej wersji z urządzenia, tą samą drogą co pierwotna instalacja (rozdz. 5.1) — nowa wersja zastępuje poprzednią pod tym samym identyfikatorem, z zachowaniem dotychczasowej konfiguracji i stanu włączenia.

**Zakres poleceń kontraktu komunikacji.** Kontrakt komunikacji (rozdz. 9; Kontrakty komunikacji, rozdz. 17) definiuje trzy polecenia obszaru rozszerzeń: `extension.list`, `extension.install` i `extension.toggle`. Wycofanie rozszerzenia z użycia realizowane jest zmianą stanu włączenia poleceniem `extension.toggle`; operacja ta jest odwracalna i nie usuwa rekordu z rejestru.

---

## 6. Rozszerzenia w module Agents

Moduł Agents (Koncepcja platformy, rozdz. 11.15) jest udokumentowanym miejscem operacyjnego wykorzystania rozszerzeń. Przy budowie agenta w oknie Agent Builder Operator dobiera umiejętności w oknie Skills Manager, podłącza integracje zewnętrzne przez okno Connectors Manager, a zakres uprawnień agenta ustala w oknie Permissions Center.

Cztery rodzaje rozszerzeń z rozdziału 1 przypisują się do dwóch okien modułu Agents zgodnie z ich charakterem.

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

Po utworzeniu agent jest, zgodnie z rozdziałem 11.15 Koncepcji platformy, komponentem platformowym dostępnym jako wykonawca zadań w środowiskach TalkIn, WorkSpace i CodeStudio, a po przypisaniu do roli — także w środowisku MultitaskingAI (Koncepcja platformy, rozdz. 9.4, rozdz. 13). Rozszerzenia podłączone do agenta pozostają z nim związane niezależnie od tego, w którym środowisku lub w jakiej roli agent jest w danej chwili wykorzystywany — zgodnie z zasadą, że agent jest niezależny od środowiska (Model danych, rozdz. 13).

---

## 7. Rozszerzenia w oknie konfiguracji platformy

Okno konfiguracji (Architektura, rozdz. 13) udostępnia pełny zakres ustawień wpływających na aplikację, w tym osobno wymienione sekcje „rozszerzenia" oraz „integracje". Rozróżnienie to nie jest przypadkowe — obie sekcje sąsiadują w tym samym oknie, lecz operują na odrębnych encjach.

| Sekcja okna konfiguracji | Encja | Zakres |
|---|---|---|
| „integracje" | Kanał modelu (Model danych, rozdz. 8.1) | Sposób połączenia z samym modelem językowym: API, CLI, SSH, HTTP (Architektura, rozdz. 9). |
| „rozszerzenia" | Rozszerzenie (Model danych, rozdz. 9) | Wtyczki, umiejętności, konektory i serwery MCP opisane w niniejszym dokumencie. |

```
OKNO KONFIGURACJI (Architektura, rozdz. 13)
├── sekcja „integracje"  → encja Kanał modelu (Model danych, rozdz. 8.1)
│        API · CLI · SSH · HTTP   (połączenie z modelem — Architektura, rozdz. 9)
│
└── sekcja „rozszerzenia" → encja Rozszerzenie (Model danych, rozdz. 9)   ── REJESTR
         wtyczki · umiejętności · konektory · serwery MCP
         ├── extension.list          → wykaz rozszerzeń obu źródeł
         ├── extension.install       → instalacja rozszerzenia Personal
         ├── extension.toggle        → zmiana stanu włączenia
         └── config.get / config.set → edycja atrybutu Konfiguracja
```

Sekcja rozszerzeń pełni funkcję rejestru: prezentuje wszystkie rozszerzenia zarejestrowane na serwerze — Danaco Plugin i Personal łącznie, zgodnie z poleceniem `extension.list` (rozdz. 9) — niezależnie od tego, czy dane rozszerzenie jest w danej chwili podłączone do jakiegokolwiek agenta. Rejestracja rozszerzenia i jego wykorzystanie operacyjne są, zgodnie z zasadą jawności i konfigurowalności zależności (Koncepcja platformy, rozdz. 14, zasada 3), dwiema odrębnymi, świadomymi decyzjami Operatora.

| Decyzja | Miejsce | Co ustala |
|---|---|---|
| Rejestracja rozszerzenia | Okno konfiguracji — sekcja „rozszerzenia" (rozdz. 7) | Czy rozszerzenie istnieje w rejestrze serwera; jego stan włączenia i konfigurację. |
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
   ┌─────────────────────────────┬─────────────────────────────────┐
   │ 1. Stan wyjściowy WYŁĄCZONE  │ 2. Zakres uprawnień przy         │
   │    dla rozszerzeń Personal   │    podłączeniu do agenta         │
   │    (rozdz. 5.2)              │    — Permissions Center (rozdz. 6)│
   │                             │                                  │
   │ Operator świadomie włącza    │ dostęp sieciowy, odczyt i zapis  │
   │                             │ plików (Architektura, rozdz. 12)  │
   └─────────────────────────────┴─────────────────────────────────┘
        │
        ▼   profil izolacji roli w środowisku MultitaskingAI
            (Koncepcja, rozdz. 6.5, poziom zasięgu „Rola")
   obejmuje pośrednio rozszerzenia podłączone do agenta pełniącego tę rolę,
   bez odrębnej konfiguracji izolacji per rozszerzenie
```

Zakres uprawnień nadawany w Permissions Center pozostaje spójny z zakresami izolacji technicznej zdefiniowanymi w rozdziale 12 Architektury i rozwiniętymi w oknie konfiguracji punktów izolacji (Koncepcja platformy, rozdz. 6.4–6.6) — w szczególności z zakresami dostępu sieciowego oraz odczytu i zapisu plików, istotnymi dla konektorów i serwerów MCP łączących agenta z systemami poza platformą. Profil izolacji przypisany do roli agenta w środowisku MultitaskingAI obejmuje w ten sposób pośrednio również rozszerzenia podłączone do agenta pełniącego tę rolę, bez potrzeby odrębnej konfiguracji izolacji per rozszerzenie.

---

## 9. Komunikacja klient–serwer

Rozdział 5.9 dokumentu Kontrakty komunikacji definiuje trzy polecenia obszaru rozszerzeń, wymieniane przez klienta i serwer kanałem WebSocket (Architektura, rozdz. 11).

| Polecenie | Kierunek | Opis |
|---|---|---|
| `extension.list` | klient → serwer | Wykaz rozszerzeń zarejestrowanych na serwerze, obu źródeł (Danaco Plugin i Personal). |
| `extension.install` | klient → serwer | Instalacja rozszerzenia Personal przesłanego z urządzenia (rozdz. 5.1). |
| `extension.toggle` | klient → serwer | Włączenie lub wyłączenie rozszerzenia, niezależnie od źródła (rozdz. 5.2). |

Konfiguracja parametrów rozszerzenia (rozdz. 5.3) korzysta z ogólnego mechanizmu konfiguracji platformy — poleceń `config.get` i `config.set` (Kontrakty komunikacji, rozdz. 14) — stosowanego do atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4).

| Polecenie ogólne | Zastosowanie do rozszerzeń |
|---|---|
| `config.get` | Odczyt atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4). |
| `config.set` | Zapis atrybutu Konfiguracja encji Rozszerzenie (rozdz. 4). |

Zmiana stanu rozszerzenia — instalacja, włączenie, wyłączenie lub zmiana konfiguracji — jest, zgodnie z ogólną zasadą komunikacji (Kontrakty komunikacji, rozdz. 1), rozgłaszana przez serwer do wszystkich połączonych urządzeń, tak aby rejestr rozszerzeń pozostawał spójny niezależnie od urządzenia, z którego Operator w danej chwili korzysta.

```
Urządzenie A            Serwer                       Urządzenie B … N
    │ extension.toggle      │                                │
    │──────────────────────►│                                │
    │                       │ zmiana stanu encji Rozszerzenie │
    │◄──── potwierdzenie ────│                                │
    │                       │──────── rozgłoszenie zmiany ───►│
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
 ═══════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │              [znacznik ▼]│ rozszerzenia
            │                      │                       [⋮]│ (rozszerzenie
            │ ─────────────────    │                          │  boczne,
            │ Execution Loop       │                          │  warstwa 3)
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ═══════════════════════════════════════════════════════════════════════════
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

## Załącznik A. Tabela zestawcza rozszerzeń

Zestawienie krzyżowe czterech rodzajów rozszerzeń (rozdz. 1) i dwóch źródeł (rozdz. 2), ze wskazaniem zawartości każdego pola, okna modułu Agents, w którym dane rozszerzenie jest podłączane do agenta (rozdz. 6), oraz warstwy widoczności elementów interfejsu wnoszonych przez rozszerzenie (rozdz. 12).

| Rodzaj | Źródło Danaco Plugin | Źródło Personal | Okno podłączenia w module Agents | Warstwa wnoszonych elementów |
|---|---|---|---|---|
| Wtyczka | Wbudowane funkcje przetwarzania dokumentu, dostarczone z pakietem serwera. | Funkcje przesłane przez Operatora, właściwe jego potrzebom. | Skills Manager | 2–3 |
| Umiejętność (skill) | Wbudowane sposoby prowadzenia analiz określonego rodzaju. | Umiejętności zdefiniowane przez Operatora dla własnych rodzajów zadań. | Skills Manager | 3–4 |
| Konektor | Wbudowane integracje z usługami dostarczanymi razem z platformą. | Integracje z zewnętrznymi systemami Operatora, przesłane z urządzenia. | Connectors Manager | 2–3 |
| Serwer MCP | Wbudowane serwery MCP dostarczone z pakietem serwera. | Serwery MCP wskazane i przesłane przez Operatora. | Connectors Manager | 2–4 |

Niezależnie od pola tabeli, w którym mieści się dane rozszerzenie, obowiązuje ten sam jednolity kontrakt integracji (rozdz. 3): rodzaj i źródło zmieniają zakres zastosowania oraz okno podłączenia, nie zmieniają zaś sposobu rejestracji, sterowania stanem ani konfiguracji.

---

## Załącznik B. Szablony konfiguracji rozszerzenia

Poniższe szablony porządkują pola opisane w rozdziałach 4, 5, 11 i 12, zgodnie z definicją encji Rozszerzenie (Model danych, rozdz. 9). Wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna" (Koncepcja platformy, rozdz. 6, rozdz. 14).

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
| 2 | Profil izolacji przypisany tej roli obejmuje zakres dostępu sieciowego stosowany również do połączeń realizowanych przez konektor. | rozdz. 8; Koncepcja platformy, rozdz. 6.5 (poziom „Rola") |
| 3 | Nie jest wymagana odrębna konfiguracja izolacji dla samego rozszerzenia. | rozdz. 8 |

### C.4. Wywołanie funkcji rozszerzenia poleceniem języka naturalnego

| Krok | Działanie i skutek | Odniesienie |
|---|---|---|
| 1 | Użytkownik formułuje polecenie w oknie Chat Window, bez wskazywania rozszerzenia z listy. | rozdz. 11.3 |
| 2 | Koordynator dekomponuje zlecenie na zadania i przydziela je Wykonawcy dysponującemu wymaganym narzędziem rozszerzenia. | rozdz. 11.1–11.2 |
| 3 | Przebieg zadania, wymiana komunikatów sterujących i wynik kontroli jakości prezentowane są w oknie Execution Loop Window. | rozdz. 11.2 |
| 4 | Parametry przebiegu ujawniają się jako znacznik kontekstowy warstwy 2 wyłącznie na czas wykonania zadania; stan spoczynku interfejsu pozostaje niezmieniony. | rozdz. 12 |
| 5 | Wynik i wyjaśnienie kontekstu wracają do okna Chat Window, gdzie Użytkownik zatwierdza albo przerywa dalsze działania. | rozdz. 11.3 |

---

*Koniec dokumentu. Danaco Console — Rozszerzenia, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
