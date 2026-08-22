# Danaco Console — Integracja modeli

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
| **Tytuł** | Integracja modeli |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (główny) · projektant (fragmenty dotyczące okna konfiguracji i znaczników kontekstowych) |
| **Przeznaczenie** | Ustala warstwę dostawcy modelu: cztery kanały integracji, mechanizm adaptera, encję kanału modelu, tożsamość i personę modelu, przypisanie modelu per sesja i per rola oraz obsługę modelu w obu kanałach komunikacji operacyjnej |
| **Zakres** | warstwa dostawcy modelu, adaptery kanałów, kanał modelu jako encja, tożsamość i persona modelu, przypisanie modelu, konfiguracja kanału, model w warstwie komunikacji operacyjnej, warstwy widoczności, relacja z izolacją i z modułem Agents, bezpieczeństwo danych dostępowych, wykaz komend kontraktu obszarów `channel`, `identity`, `model` i `account` |
| **Poza zakresem** | pełna konfiguracja agenta — [Model konfiguracji](model-konfiguracji.md); model danych encji kanału modelu w pełnym zakresie pól — [Model danych](model-danych.md) rozdz. 11; okno Agent Builder i moduł Agents jako całość — [Moduł Agents](../moduly/agents.md) |
| **Dokument nadrzędny** | [Architektura techniczna](architektura.md) |
| **Dokumenty powiązane** | [Architektura techniczna](architektura.md) · [Model danych](model-danych.md) · [Model konfiguracji](model-konfiguracji.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Izolacja i zależności](izolacja-i-zaleznosci.md) · [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszary `channel`, `identity`, `model`, `account`; struktury `Channel`, `Account`, `IdentityCategory`, `IdentityDocument`, `IdentityLayerContent`, `ChannelCredentialStatus`; wyliczenia `ProviderTransport`, `ReasoningEffort`, `ConfigAxis`, `IdentityLayer`, `IdentityMode`, `AccountKind`) |
| **Zasada nadrzędna** | Rdzeń zna wyłącznie jednolity kontrakt warstwy dostawcy modelu — „wyślij zapytanie, odbierz odpowiedź” — i nie zależy od szczegółów technicznych żadnego dostawcy ani sposobu połączenia z nim |

Dokument opisuje warstwę dostawcy modelu w platformie Danaco Console: cztery kanały integracji modeli sztucznej inteligencji — API, CLI, SSH i HTTP — mechanizm adapterów realizujących te kanały, sposób przypisania modelu do sesji lub roli, mechanizm tożsamości i persony modelu oraz miejsce konfiguracji kanału w oknie konfiguracji. Dokument opisuje ponadto obsługę obu kanałów komunikacji operacyjnej platformy — Użytkownik ↔ Wykonawca w oknie Chat Window oraz Koordynator ↔ Wykonawca w oknie Execution Loop Window — wraz z warstwami widoczności funkcji wyboru i konfiguracji modelu. Podstawę stanowią rozdział 9 (Integracja modeli) oraz rozdziały 5 i 12 dokumentu [Architektura techniczna](architektura.md), a także rozdział 11 (Kanały modeli) dokumentu [Model danych](model-danych.md).

---

## Spis treści

1. [Pozycjonowanie warstwy integracji modeli](#1-pozycjonowanie-warstwy-integracji-modeli)
2. [Warstwa dostawcy modelu](#2-warstwa-dostawcy-modelu)
3. [Cztery kanały integracji](#3-cztery-kanały-integracji)
   - [3.1 Zestawienie czterech kanałów](#31-zestawienie-czterech-kanałów)
   - [3.2 Kryteria doboru kanału](#32-kryteria-doboru-kanału)
4. [Adaptery kanałów](#4-adaptery-kanałów)
   - [4.1 Adapter kanału a protokół, uwierzytelnienie i cel połączenia](#41-adapter-kanału-a-protokół-uwierzytelnienie-i-cel-połączenia)
   - [4.2 Właściwości mechanizmu adaptera](#42-właściwości-mechanizmu-adaptera)
   - [4.3 Uwaga terminologiczna — adapter kanału a rozszerzenie](#43-uwaga-terminologiczna--adapter-kanału-a-rozszerzenie)
5. [Kanał modelu jako encja](#5-kanał-modelu-jako-encja)
6. [Tożsamość i persona modelu](#6-tożsamość-i-persona-modelu)
7. [Przypisanie modelu per sesja i per rola](#7-przypisanie-modelu-per-sesja-i-per-rola)
   - [7.1 Przypisanie per sesja](#71-przypisanie-per-sesja)
   - [7.2 Przypisanie per rola](#72-przypisanie-per-rola)
8. [Konfiguracja kanału w oknie konfiguracji](#8-konfiguracja-kanału-w-oknie-konfiguracji)
   - [8.1 Kroki konfiguracji kanału](#81-kroki-konfiguracji-kanału)
   - [8.2 Komunikaty warstwy komunikacji](#82-komunikaty-warstwy-komunikacji)
   - [8.3 Warstwowość konfiguracji kanału](#83-warstwowość-konfiguracji-kanału)
9. [Model językowy w warstwie komunikacji operacyjnej](#9-model-językowy-w-warstwie-komunikacji-operacyjnej)
   - [9.1 Kanał Użytkownik ↔ Wykonawca — obsługa w Chat Window](#91-kanał-użytkownik--wykonawca--obsługa-w-chat-window)
   - [9.2 Kanał Koordynator ↔ Wykonawca — obsługa w Execution Loop Window](#92-kanał-koordynator--wykonawca--obsługa-w-execution-loop-window)
10. [Warstwy widoczności funkcji integracji modeli](#10-warstwy-widoczności-funkcji-integracji-modeli)
11. [Relacja z izolacją techniczną](#11-relacja-z-izolacją-techniczną)
12. [Relacja z modułem Agents](#12-relacja-z-modułem-agents)
13. [Bezpieczeństwo danych dostępowych](#13-bezpieczeństwo-danych-dostępowych)
14. [Zgodność z zasadami nadrzędnymi](#14-zgodność-z-zasadami-nadrzędnymi)
15. [Wykaz komend kontraktu — obszary `channel`, `identity` i `model`](#15-wykaz-komend-kontraktu--obszary-channel-identity-i-model)
   - [15.1 Struktury i wyliczenia normatywne](#151-struktury-i-wyliczenia-normatywne)
   - [15.2 Diagram sekwencji — zapis kanału i odczyt tożsamości obowiązującej](#152-diagram-sekwencji--zapis-kanału-i-odczyt-tożsamości-obowiązującej)
   - [15.3 Struktury danych pełne](#153-struktury-danych-pełne)
   - [15.4 Kody błędów zastosowane do warstwy dostawcy modelu](#154-kody-błędów-zastosowane-do-warstwy-dostawcy-modelu)
16. [Rozbieżność modelu czterech kanałów z kontraktem](#16-rozbieżność-modelu-czterech-kanałów-z-kontraktem)
17. [Słowniczek pojęć](#17-słowniczek-pojęć)
18. [Kryteria odbioru](#18-kryteria-odbioru)
19. [Załącznik A. Szablon konfiguracji kanału modelu](#załącznik-a-szablon-konfiguracji-kanału-modelu)
   - [A.1. Pola konfiguracji kanału modelu](#a1-pola-konfiguracji-kanału-modelu)
   - [A.2. Szablon kanału API](#a2-szablon-kanału-api)
   - [A.3. Szablon kanału CLI](#a3-szablon-kanału-cli)
   - [A.4. Szablon kanału SSH](#a4-szablon-kanału-ssh)
   - [A.5. Szablon kanału HTTP](#a5-szablon-kanału-http)
20. [Załącznik B. Scenariusze użycia](#załącznik-b-scenariusze-użycia)
   - [B.1. Przypisanie modelu bezpośrednio do roli w MultitaskingAI](#b1-przypisanie-modelu-bezpośrednio-do-roli-w-multitaskingai)
   - [B.2. Model pomocniczy na hoście własnym przez kanał SSH](#b2-model-pomocniczy-na-hoście-własnym-przez-kanał-ssh)
   - [B.3. Ten sam model bazowy, dwie odrębne tożsamości](#b3-ten-sam-model-bazowy-dwie-odrębne-tożsamości)
   - [B.4. Zmiana poświadczenia konta bez wpływu na kanały powiązane](#b4-zmiana-poświadczenia-konta-bez-wpływu-na-kanały-powiązane)
21. [Załącznik C. Pełny wykaz komend kontraktu warstwy modeli](#załącznik-c-pełny-wykaz-komend-kontraktu-warstwy-modeli)
   - [Obszar `channel` — 6 komend](#obszar-channel--6-komend)
   - [Obszar `identity` — 5 komend](#obszar-identity--5-komend)
   - [Obszar `model` — 1 komenda](#obszar-model--1-komenda)
   - [Obszar `account` — 6 komend](#obszar-account--6-komend)

---

## 1. Pozycjonowanie warstwy integracji modeli

Danaco Console jest systemem operacyjnym do pracy z modelami sztucznej inteligencji — wszystkie środowiska, moduły i role wykonują pracę przez wywołania modelu, a sposób, w jaki te wywołania docierają do konkretnego dostawcy, jest przedmiotem niniejszego dokumentu. Poniższa tabela zbiera ustalenia dotyczące miejsca warstwy integracji modeli w architekturze platformy.

| Zagadnienie | Ustalenie |
|---|---|
| Charakter platformy | System operacyjny do pracy z modelami sztucznej inteligencji; każda praca — od modułu Studio po role Executor 1 i Coordinator w środowisku MultitaskingAI — wykonywana jest przez wywołania modelu |
| Miejsce w architekturze | Integracja modeli jest odpowiedzialnością warstwy Rdzeń (rozdz. 4 dokumentu Architektury) |
| Forma realizacji | Jednolita warstwa dostawcy modelu, oddzielona od reszty rdzenia kontraktem |
| Zasada architektoniczna | „Rozdzielenie przez kontrakty” (rozdz. 1 dokumentu Architektury) |
| Wiedza rdzenia | Wyłącznie jednolity kontrakt warstwy dostawcy (rozdz. 2) — bez szczegółów technicznych dostawcy ani sposobu połączenia z nim |
| Widoczność dla użytkownika | Brak osobnego modułu ani środowiska; warstwa usługowa |
| Punkty styku użytkownika | Okno konfiguracji (rozdz. 8), przypisanie modelu do sesji lub roli (rozdz. 7) oraz oba kanały komunikacji operacyjnej — Chat Window i Execution Loop Window (rozdz. 9) — bez dedykowanego okna operacyjnego |

Poniższy schemat wskazuje, kto korzysta z warstwy dostawcy modelu oraz przez jakie punkty styka się z nią użytkownik.

```
KORZYSTAJĄCY Z WARSTWY DOSTAWCY MODELU
  • Moduły platformy            (Studio, Research, Developer, Roundtable … — piętnaście modułów)
  • Role środowiska             (Executor 1, Executor 2, Coordinator, Executor 3 / Validator)
    MultitaskingAI
  • Moduł Agents                (przy budowie agentów)
                         │
                         ▼
             WARSTWA DOSTAWCY MODELU        (warstwa usługowa, niewidoczna jako osobny moduł)
                         │
        styk użytkownika:  • okno konfiguracji                         (rozdz. 8)
                           • przypisanie modelu do sesji lub roli       (rozdz. 7)
```

Dzięki oddzieleniu kontraktem rdzeń nie musi znać szczegółów technicznych żadnego dostawcy ani sposobu połączenia z nim — zna wyłącznie jednolity kontrakt warstwy dostawcy, opisany w rozdziale 2.

---

## 2. Warstwa dostawcy modelu

Rdzeń operuje na jednej, wspólnej dla wszystkich kanałów abstrakcji — „wyślij zapytanie, odbierz odpowiedź” — i jest to jedyny kontrakt, jaki zna. Kontrakt ten nie zależy od tego, czy zapytanie trafia do modelu przez klucz API, token narzędzia wiersza poleceń, połączenie SSH z hostem zdalnym, czy przez model przeglądarkowy kanałem HTTP. Poniższa tabela zestawia trzy konsekwencje tej zasady.

| Konsekwencja zasady | Treść |
|---|---|
| Rdzeń niezależny od kanału | Logika sesji, ról, kolejek i orkiestracji operuje wyłącznie na pojęciu „model przypisany do sesji lub roli”, bez odwołań do konkretnego kanału integracji |
| Dodanie kanału = dostarczenie adaptera | Nowy sposób połączenia nie wymaga zmian w rdzeniu, lecz wyłącznie nowego adaptera realizującego jednolity kontrakt (rozdz. 4); jest to zastosowanie zasady „Konfigurowalność zamiast wartości wpisanych na stałe” (rozdz. 1 dokumentu Architektury) |
| Jednolity strumień odpowiedzi | Odpowiedź przekazywana jest na żywo kanałem WebSocket do klienta (rozdz. 11 dokumentu Architektury); adapter przekształca odpowiedź dostawcy do jednolitego formatu strumienia rdzenia, identycznie dla każdego z czterech kanałów |

Poniższy schemat przedstawia miejsce warstwy dostawcy modelu w architekturze platformy.

```
                    RDZEŃ (Go) — orkiestracja sesji, procesy, polityki
                                        │
                         abstrakcja: „wyślij zapytanie, odbierz odpowiedź”
                                        │
                        WARSTWA DOSTAWCY MODELU
                                        │
        ┌───────────────┬───────────────┼───────────────┬───────────────┐
        ▼               ▼               ▼               ▼
   Adapter API      Adapter CLI     Adapter SSH     Adapter HTTP
   (klucz, żądanie  (token           (host zdalny,   (model
    HTTP)            wiersza          własny model    przeglądarkowy)
                      poleceń)         pomocniczy)
        │               │               │               │
        ▼               ▼               ▼               ▼
   Dostawca API     Narzędzie CLI    Host zdalny      Model
                                      (model           przeglądarkowy
                                      pomocniczy)
```

---

## 3. Cztery kanały integracji

Platforma udostępnia cztery kanały integracji modeli, ustalone w rozdziale 9 dokumentu Architektury. Każdy kanał odpowiada innemu sposobowi, w jaki model jest fizycznie osiągalny przez warstwę dostawcy.

### 3.1. Zestawienie czterech kanałów

| Kanał | Definicja | Sposób połączenia | Pośrednik między platformą a dostawcą | Dane uwierzytelniające | Właściwe zastosowanie |
|---|---|---|---|---|---|
| API | Połączenie z modelem przez żądanie HTTP uwierzytelnione kluczem dostępu | Bezpośrednie żądanie HTTP do punktu końcowego dostawcy | Brak — komunikacja bezpośrednia z punktem końcowym HTTP | Klucz dostępu (identyfikuje konto użytkownika u dostawcy i autoryzuje żądanie) | Dostawcy udostępniający interfejs programistyczny; podstawowy kanał pracy operacyjnej w sesjach modułów i w rolach środowiska MultitaskingAI |
| CLI | Połączenie z modelem przez narzędzie wiersza poleceń dostawcy | Wywołanie narzędzia wiersza poleceń | Narzędzie wiersza poleceń zainstalowane lokalnie lub zdalnie, samodzielnie zarządzające komunikacją z dostawcą | Token narzędzia wiersza poleceń | Modele i narzędzia udostępniane przez dostawcę wyłącznie lub przede wszystkim w formie narzędzia wiersza poleceń |
| SSH | Połączenie z własnym modelem pomocniczym uruchomionym na hoście zdalnym, przez protokół SSH | Połączenie SSH z hostem zdalnym | Host zdalny uruchamiany i utrzymywany samodzielnie przez użytkownika | Dane dostępowe SSH (host, konto, klucz lub hasło) | Modele nieudostępniane przez zewnętrznego dostawcę w modelu API lub CLI, uruchamiane na własnej infrastrukturze; zadania wyspecjalizowane lub wymagające pełnej kontroli nad środowiskiem wykonawczym |
| HTTP | Połączenie z modelem przeglądarkowym przez interfejs webowy dostawcy | Interakcja z interfejsem webowym | Interfejs webowy dostawcy | Sesja / dane logowania właściwe interfejsowi webowemu | Modele dostępne wyłącznie przez interfejs przeglądarkowy, nieudostępniające dedykowanego kanału API ani narzędzia wiersza poleceń |

### 3.2. Kryteria doboru kanału

| Sytuacja lub potrzeba | Właściwy kanał |
|---|---|
| Dostawca udostępnia interfejs programistyczny; pożądana przewidywalność kosztów, jakość dokumentacji interfejsu i stabilność połączenia | API |
| Model lub narzędzie dostawcy dostępne wyłącznie lub przede wszystkim jako narzędzie wiersza poleceń | CLI |
| Własny model pomocniczy uruchamiany na własnej infrastrukturze; potrzeba pełnej kontroli nad środowiskiem wykonawczym lub uniknięcia kosztów zewnętrznego dostawcy | SSH |
| Model pożądany nie udostępnia kanału API ani CLI, a jedyną drogą dostępu jest jego interfejs przeglądarkowy | HTTP |

Kanał integracji jest atrybutem kanału modelu (rozdz. 5), a nie atrybutem sesji ani roli wprost — sesja lub rola odwołuje się do skonfigurowanego kanału modelu, który dopiero określa, którym z czterech kanałów model jest osiągany.

---

## 4. Adaptery kanałów

Każdy z czterech kanałów integracji realizuje adapter jednolitej warstwy dostawcy modelu. Adapter jest komponentem technicznym tłumaczącym abstrakcję rdzenia „wyślij zapytanie, odbierz odpowiedź” na konkretny protokół i sposób uwierzytelnienia właściwy danemu kanałowi. Poniższy schemat przedstawia warstwę adapterów: to samo zapytanie rdzenia jest kierowane do właściwego adaptera, który tłumaczy je na protokół kanału, a odpowiedź dostawcy wraca przekształcona do jednolitego formatu strumienia.

```
                                RDZEŃ
              jednolite zapytanie: „wyślij zapytanie, odbierz odpowiedź”
                                  │   (ten sam kontrakt dla wszystkich kanałów)
                                  ▼
        ┌─────────────────────────────────────────────────────────────────┐
        │                  WARSTWA DOSTAWCY MODELU                        │
        │        wybór adaptera właściwego skonfigurowanemu kanałowi      │
        └─────────────────────────────────────────────────────────────────┘
             │                 │                 │                 │
             ▼                 ▼                 ▼                 ▼
        Adapter API       Adapter CLI       Adapter SSH       Adapter HTTP
        tłumaczy na:      tłumaczy na:      tłumaczy na:      tłumaczy na:
        żądanie HTTP      wywołanie         połączenie SSH    interakcję z
        + klucz           narzędzia         z hostem          interfejsem
        dostępu           wiersza poleceń   zdalnym           webowym
                          + token
             │                 │                 │                 │
             ▼                 ▼                 ▼                 ▼
        Dostawca API      Narzędzie CLI     Host zdalny       Model
                                            (model            przeglądarkowy
                                            pomocniczy)
             │                 │                 │                 │
             └─────────────────┴────────┬────────┴─────────────────┘
                                        ▼
                    odpowiedź dostawcy → adapter przekształca ją
                    do JEDNOLITEGO formatu strumienia rdzenia
                                        │
                                        ▼
              RDZEŃ → strumień WebSocket do klienta (rozdz. 2, rozdz. 11 Architektury)
```

Rdzeń nie rozróżnia adapterów między sobą poza wyborem właściwego adaptera dla skonfigurowanego kanału modelu — każdy adapter przyjmuje to samo zapytanie i zwraca odpowiedź w tym samym, jednolitym formacie strumienia. Poniższa sekwencja ujmuje przebieg pojedynczego wywołania, identyczny dla każdego z czterech kanałów.

```
Sekwencja pojedynczego wywołania (dowolny kanał):

  Rdzeń  ──(jednolite zapytanie)──►  Adapter kanału
  Adapter kanału  ──(protokół kanału + dane dostępowe)──►  Dostawca
  Dostawca  ──(odpowiedź w formacie dostawcy)──►  Adapter kanału
  Adapter kanału  ──(jednolity format strumienia)──►  Rdzeń  ──►  klient (WebSocket)
```

### 4.1. Adapter kanału a protokół, uwierzytelnienie i cel połączenia

| Adapter | Protokół / sposób wywołania | Uwierzytelnienie | Cel połączenia |
|---|---|---|---|
| Adapter API | Żądanie HTTP do punktu końcowego dostawcy | Klucz dostępu | Dostawca API |
| Adapter CLI | Wywołanie narzędzia wiersza poleceń | Token narzędzia wiersza poleceń | Narzędzie wiersza poleceń (lokalne lub zdalne) |
| Adapter SSH | Połączenie SSH z hostem zdalnym | Dane dostępowe SSH (host, konto, klucz lub hasło) | Host zdalny z modelem pomocniczym |
| Adapter HTTP | Interakcja z interfejsem webowym | Dane logowania interfejsu webowego | Model przeglądarkowy |

### 4.2. Właściwości mechanizmu adaptera

| Właściwość | Konsekwencja |
|---|---|
| Rozszerzalność | Każdy kolejny kanał integracji realizowany jest wyłącznie przez dostarczenie nowego adaptera, bez zmian w logice sesji, ról, kolejek i orkiestracji |
| Niezależność adapterów | Błąd lub niedostępność jednego kanału nie wpływa na działanie pozostałych trzech |
| Wymienność kanału | Zmiana kanału dla danej sesji lub roli — z API na SSH — jest wyłącznie zmianą konfiguracji kanału modelu (rozdz. 5), bez zmiany logiki modułu ani środowiska |

### 4.3. Uwaga terminologiczna — adapter kanału a rozszerzenie

Adapter kanału jest pojęciem odrębnym od rozszerzenia (wtyczki, umiejętności, konektora lub serwera MCP) opisanego w rozdziale 14 dokumentu Architektury. Poniższa tabela zestawia oba mechanizmy.

| Cecha | Adapter kanału (niniejszy rozdział) | Rozszerzenie (Danaco Plugin lub Personal) |
|---|---|---|
| Kategoria | Komponent techniczny warstwy dostawcy modelu | Element funkcjonalny dodawany do agenta lub sesji |
| Przykładowe postacie | Adapter API, Adapter CLI, Adapter SSH, Adapter HTTP | Wtyczka, umiejętność, konektor lub serwer MCP |
| Sposób dostarczenia | Wraz z pakietem aplikacyjnym platformy | Dodawane i konfigurowane przez użytkownika |
| Widoczność dla użytkownika | Niewidoczny jako osobny element do zainstalowania | Wybieralne i konfigurowalne w module Agents (rozdz. 12) oraz w innych modułach platformy |
| Realizowany kontrakt | Jednolity kontrakt warstwy dostawcy modelu | Jednolity kontrakt w obrębie własnej kategorii |

Oba mechanizmy realizują jednolity kontrakt w obrębie swojej kategorii, lecz odpowiadają za różne warstwy platformy.

---

## 5. Kanał modelu jako encja

Kanał modelu jest encją modelu danych reprezentującą skonfigurowane połączenie z konkretnym modelem przez jeden z czterech kanałów integracji. Atrybuty kanału modelu, ustalone w rozdziale 11.1 dokumentu Modelu danych, są następujące.

| Atrybut | Opis |
|---|---|
| Identyfikator | Jednoznaczny identyfikator kanału modelu w bazie danych |
| Typ | Jeden z czterech kanałów: API, CLI, SSH, HTTP |
| Odwołanie do danych dostępowych | Wskazanie na dane uwierzytelniające przechowywane poza bazą (rozdz. 13) |
| Nazwa modelu | Techniczna nazwa modelu bazowego udostępnianego przez dostawcę, wykorzystywana przy wywołaniu |
| Przypisanie | Sesja albo rola, do której kanał modelu jest przypisany (rozdz. 7) |

Kanał modelu jest jednostką konfiguracji niższego rzędu niż agent (rozdz. 12) — opisuje wyłącznie sposób techniczny dotarcia do modelu, jak zestawia poniższa tabela.

| Kanał modelu opisuje | Kanał modelu nie obejmuje |
|---|---|
| Sposób techniczny dotarcia do modelu: typ kanału, dane dostępowe, nazwę modelu bazowego | Tożsamości, instrukcji systemowych, umiejętności, wtyczek, konektorów, pamięci ani uprawnień właściwych agentowi |

Kanał modelu może istnieć samodzielnie, przypisany bezpośrednio do sesji lub roli (rozdz. 7), albo stanowić podstawę techniczną, na której zbudowany jest agent w module Agents.

---

## 6. Tożsamość i persona modelu

Poza warstwą techniczną opisaną w rozdziale 5, każdy model przypisany do sesji lub roli może otrzymać własną tożsamość — nazwę i personę — niezależną od technicznej nazwy modelu bazowego. Mechanizm ten realizuje zasadę kompletności dostępu z rozdziału 13 dokumentu Architektury, zgodnie z którą nie istnieje ustawienie niedostępne z okna konfiguracji — dotyczy to również zachowania modeli, tożsamości modeli i promptów systemowych. Tożsamość modelu obejmuje dwa elementy.

| Element tożsamości | Opis | Zależność |
|---|---|---|
| Nazwa | Dowolna nazwa nadawana modelowi przez użytkownika; służy rozpoznawalności modelu w interfejsie — w panelach ról środowiska MultitaskingAI (rozdz. 7.2) oraz w oknach modułu Roundtable, gdzie kilka modeli pracuje równolegle nad tym samym zagadnieniem | Niezależna od technicznej nazwy modelu bazowego przechowywanej w atrybucie kanału modelu (rozdz. 5) |
| Prompt systemowy | Instrukcja wgrywana dla danego przypisania modelu; określa sposób jego działania, ton wypowiedzi, zakres kompetencji lub rolę pełnioną w danej sesji lub roli; wgrywany bezpośrednio przy konfiguracji przypisania | Obowiązuje dla tego konkretnego przypisania modelu do sesji lub roli |

Mechanizm tożsamości i persony modelu jest warstwą lżejszą niż pełna konfiguracja agenta w module Agents (rozdz. 12). Poniższa tabela zestawia obie warstwy.

| Kryterium | Tożsamość i persona modelu (rozdz. 6) | Pełna konfiguracja agenta (rozdz. 12) |
|---|---|---|
| Zakres | Nazwa i prompt systemowy | Model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć, uprawnienia |
| Trwałość | Właściwa danemu przypisaniu; nie tworzy trwałego, nazwanego komponentu własnego | Trwały, nazwany komponent własny budowany w Agent Builder |
| Miejsce konfiguracji | Okno konfiguracji, przy przypisaniu do sesji lub roli | Moduł Agents, Agent Builder (rozdz. 11.15 koncepcji platformy) |
| Właściwe zastosowanie | Przypisania doraźne lub jednorazowe, w których wystarcza rozpoznawalna nazwa i określony sposób działania | Trwałe, wielokrotnie wykorzystywane jednostki AI |
| Odzwierciedlenie w szablonie roli | Pole „Przypisany agent lub model” (Załącznik D.2 koncepcji platformy) — ścieżka „model bazowy” | Pole „Przypisany agent lub model” — ścieżka „agent” |

Agent posiada własne atrybuty tożsamości i instrukcji systemowych jako część trwałej konfiguracji komponentu własnego, natomiast tożsamość i persona modelu dotyczą modelu przypisanego bezpośrednio, bez pośrednictwa agenta. Szablon roli środowiska MultitaskingAI odzwierciedla ten wybór wprost w polu „Przypisany agent lub model”, dopuszczającym obie ścieżki.

---

## 7. Przypisanie modelu per sesja i per rola

Wybór modelu i kanału dokonywany jest per sesja lub per rola — zasada ustalona w rozdziale 9 dokumentu Architektury. Poniższe dwa podpunkty opisują obie ścieżki przypisania.

### 7.1. Przypisanie per sesja

W środowiskach modułowych (TalkIn, WorkSpace, CodeStudio) każda karta sesji może mieć przypisany własny kanał modelu.

| Aspekt przypisania per sesja | Ustalenie |
|---|---|
| Zakres | Każda karta sesji (rozdz. 8 koncepcji platformy) w środowiskach TalkIn, WorkSpace i CodeStudio może mieć przypisany własny kanał modelu |
| Skutek | Przypisanie określa, który model odpowiada na wiadomości wysyłane przez czat (rozdz. 2.4 koncepcji platformy) w obrębie danej karty |
| Stan domyślny | Odrębny kontekst każdej karty (rozdz. 6 koncepcji platformy) oznacza domyślnie odrębne przypisanie modelu; użytkownik może jednak powiązać przypisania między kartami tam, gdzie jest to pożądane |
| Warstwa konfiguracji | Zmiana przypisania jest zmianą warstwy sesji (rozdz. 13 dokumentu Architektury), nakładającą się na warstwę domyślną bez jej trwałej modyfikacji |

### 7.2. Przypisanie per rola

W środowisku MultitaskingAI (rozdz. 9.4 i rozdz. 13 koncepcji platformy) przypisanie modelu dokonywane jest na poziomie roli, nie sesji — każda z czterech ról (Executor 1, Executor 2, Coordinator, Executor 3 / Validator) posiada własne przypisanie, niezależne od pozostałych. Szablon roli (Załącznik D.2 koncepcji platformy) ujmuje to przypisanie w polu „Przypisany agent lub model”, dopuszczającym dwie wartości.

| Wartość pola „Przypisany agent lub model” | Co niesie | Podstawa |
|---|---|---|
| Agent | Komponent własny utworzony w module Agents, niosący pełną konfigurację: model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć i uprawnienia | rozdz. 12 |
| Model bazowy | Bezpośrednie przypisanie kanału modelu wraz z tożsamością i personą ustawianymi w konfiguracji, bez pośrednictwa agenta | rozdz. 5, rozdz. 6 |

Pole „Kanał modelu” szablonu roli wskazuje ponadto, którym z czterech kanałów opisanych w rozdziale 3 dany model jest osiągany — niezależnie od tego, czy rola korzysta z agenta, czy z modelu bazowego, ponieważ agent sam w sobie opiera się na modelu bazowym połączonym przez jeden z czterech kanałów integracji. Subagent Network, uruchamiany przez Executora (rozdz. 13.3 koncepcji platformy), dziedziczy mechanizm przypisania modelu — każdy z do piętnastu podagentów może mieć własny kanał modelu i własną tożsamość, konfigurowane analogicznie do roli nadrzędnej.

Poniższy schemat przedstawia obie ścieżki przypisania w formie tekstowej.

```
KARTA SESJI (TalkIn / WorkSpace / CodeStudio)
    │
    └── Kanał modelu  ──►  Tożsamość i persona (rozdz. 6)

ROLA (MultitaskingAI: Executor 1 / Executor 2 / Coordinator / Executor 3 / Validator)
    │
    ├── Agent (moduł Agents)  ──►  model bazowy · tożsamość · instrukcje systemowe
    │                                · umiejętności · wtyczki · konektory · pamięć · uprawnienia
    │
    └── Model bazowy (bezpośrednio)  ──►  Kanał modelu ──► Tożsamość i persona (rozdz. 6)
```

---

## 8. Konfiguracja kanału w oknie konfiguracji

Konfiguracja kanałów modeli odbywa się w oknie konfiguracji, otwieranym z okna operacyjnego (rozdz. 13 dokumentu Architektury). Poniższa tabela wskazuje miejsce konfiguracji kanału wśród obszarów tego okna.

| Obszar okna konfiguracji | Powiązanie z integracją modeli |
|---|---|
| Zachowanie modeli, tożsamość modeli, prompty systemowe, integracje | Pełny zakres ustawień udostępnianych przez okno konfiguracji (rozdz. 13 dokumentu Architektury) |
| Konfiguracja kanału modelu | Obszar opisany w niniejszym rozdziale |
| Okno konfiguracji punktów izolacji | Osobny obszar (rozdz. 6.4–6.6 koncepcji platformy), opisany w rozdziale 11 niniejszego dokumentu |

### 8.1. Kroki konfiguracji kanału

| Krok | Działanie | Szczegół |
|---|---|---|
| 1 | Wybór typu kanału | API, CLI, SSH lub HTTP (rozdz. 3) |
| 2 | Wprowadzenie danych dostępowych właściwych wybranemu typowi | Klucz dostępu (API), token narzędzia wiersza poleceń (CLI), dane dostępowe SSH (SSH), dane logowania interfejsu webowego (HTTP); dane trafiają do przechowywania poza bazą (rozdz. 13), a okno konfiguracji zapisuje w encji kanału modelu wyłącznie odwołanie do nich |
| 3 | Wskazanie nazwy modelu udostępnianego przez wybrany kanał | Techniczna nazwa modelu bazowego u dostawcy |
| 4 | Nadanie tożsamości i persony | Nazwa dowolna oraz prompt systemowy wgrywany dla danego przypisania (rozdz. 6) |
| 5 | Przypisanie skonfigurowanego kanału modelu do sesji lub roli | rozdz. 7 |

Każdy element konfiguracji kanału jest opatrzony objaśnieniem kontekstowym (oznaczonym znakiem zapytania), opisującym jego działanie i wpływ na aplikację — zgodnie z zasadą przyjętą dla całego okna konfiguracji w rozdziale 13 dokumentu Architektury.

### 8.2. Komunikaty warstwy komunikacji

| Czynność | Komunikat | Podstawa |
|---|---|---|
| Zapis konfiguracji kanału | `model.channel.set` | rozdz. 5.8 dokumentu Kontraktów komunikacji |
| Przypisanie modelu do roli w środowisku MultitaskingAI | `role.assign` | rozdz. 5.8 dokumentu Kontraktów komunikacji |

### 8.3. Warstwowość konfiguracji kanału

Skonfigurowany kanał modelu, podobnie jak inne elementy konfiguracji platformy, ma strukturę warstwową (rozdz. 13 dokumentu Architektury).

| Warstwa konfiguracji | Kiedy obowiązuje | Trwałość |
|---|---|---|
| Warstwa domyślna | Kanał ustalony w oknie konfiguracji; obowiązuje przy każdej nowej sesji | Trwała, dopóki użytkownik jej nie zmieni |
| Warstwa sesji | Zmiana kanału dokonana w oknie operacyjnym dla sesji bieżącej | Nakłada się na warstwę domyślną bez jej trwałej modyfikacji |

---

## 9. Model językowy w warstwie komunikacji operacyjnej

Model językowy nie jest w Danaco Console samodzielnym uczestnikiem pracy, lecz zasobem, z którego korzysta **Wykonawca** — AI, agent lub system wykonawczy realizujący zadania. Każde wywołanie modelu wychodzi z jednego z dwóch kanałów komunikacji operacyjnej platformy i wraca do niego jako strumień odpowiedzi. Warstwa dostawcy modelu (rozdz. 2) obsługuje oba kanały tym samym kontraktem i tym samym zestawem adapterów (rozdz. 4).

| Kanał komunikacji | Uczestnicy | Okno | Rola warstwy dostawcy modelu |
|---|---|---|---|
| Kanał pierwszy | Użytkownik ↔ Wykonawca | Chat Window — lewa kolumna, stała, pełna wysokość obszaru roboczego | Realizacja wywołań wynikających z poleceń Użytkownika, z kanału modelu przypisanego do karty sesji (rozdz. 7.1) |
| Kanał drugi | Koordynator ↔ Wykonawca | Execution Loop Window — kolumna sąsiadująca z Chat Window | Realizacja wywołań wynikających z zadań przydzielonych przez Koordynatora, z kanałów modeli przypisanych do ról (rozdz. 7.2) |

Oba okna rozmieszczone są w układzie pionowym (podział lewa–prawa); obszar roboczy modułu zajmuje kolumnę prawą, dominującą.

### 9.1. Kanał Użytkownik ↔ Wykonawca — obsługa w Chat Window

Chat Window jest podstawowym mechanizmem sterowania wywołaniami modelu: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i pozwala przerwać generowanie w dowolnym momencie.

| Czynność | Przebieg w warstwie integracji modeli |
|---|---|
| Polecenie | Polecenie Użytkownika zamieniane jest na jednolite zapytanie rdzenia i kierowane do adaptera kanału modelu przypisanego do bieżącej karty sesji (rozdz. 5, rozdz. 7.1); tożsamość i persona przypisania (rozdz. 6) wchodzą do wywołania jako prompt systemowy |
| Strumień odpowiedzi | Adapter przekształca odpowiedź dostawcy do jednolitego formatu strumienia rdzenia, identycznie dla każdego z czterech kanałów integracji (rozdz. 4); rdzeń przekazuje strumień kanałem WebSocket do Chat Window, gdzie odpowiedź narasta w miarę napływu fragmentów |
| Przerwanie generowania | Polecenie przerwania zamyka strumień i kończy wywołanie po stronie adaptera; fragment odebrany do chwili przerwania pozostaje w historii karty sesji jako wynik zamknięty, a karta wraca do stanu gotowości do przyjęcia kolejnego polecenia |
| Zmiana modelu w toku pracy | Wybór innego kanału modelu obowiązuje od kolejnego polecenia i stanowi zmianę warstwy sesji (rozdz. 8.3), bez modyfikacji warstwy domyślnej |
| Wyjaśnienie wyniku | Zapytanie o uzasadnienie lub kontekst wyniku realizowane jest jako kolejne wywołanie tego samego kanału modelu, z zachowaniem kontekstu karty |

Aktywny model i aktywny wykonawca prezentowane są w Chat Window jako znaczniki kontekstowe (rozdz. 10).

### 9.2. Kanał Koordynator ↔ Wykonawca — obsługa w Execution Loop Window

Execution Loop Window prezentuje pętlę wykonawczą: dekompozycję zlecenia na zadania, przydział zadań Wykonawcom, kontrolę jakości wyników i sterowanie przebiegiem. Warstwa integracji modeli obsługuje tę pętlę wywołaniami prowadzonymi równolegle dla wielu zadań.

| Mechanizm pętli | Realizacja w warstwie integracji modeli |
|---|---|
| Dobór modelu do zadania | Koordynator przydziela zadanie Wykonawcy wraz z kanałem modelu wynikającym z przypisania roli (rozdz. 7.2); reguła doboru — model o wyższej zdolności rozumowania do zadań analitycznych, model o krótszym czasie odpowiedzi do zadań prostych — zapisana jest jako ustawienie konfiguracyjne profilu doboru i podlega zmianie z okna konfiguracji (rozdz. 8) |
| Równoległe wywołania dla wielu zadań | Każde zadanie kolejki realizowane jest własnym wywołaniem własnego kanału modelu; adaptery działają niezależnie od siebie (rozdz. 4.2), więc zadania przydzielone różnym rolom i różnym dostawcom przebiegają równocześnie, a Execution Loop Window prezentuje stan każdego wywołania osobno |
| Kontrola jakości wyniku | Wynik zadania kierowany jest do roli Executor 3 / Validator, której kanał modelu jest niezależny od kanału wykonawcy zadania; ocena wraca do Koordynatora jako wynik kontroli jakości prezentowany w oknie pętli |
| Ponowienie z korektą | Wynik odrzucony w kontroli jakości powoduje ponowne wywołanie tego samego zadania z poleceniem uzupełnionym o treść uwag; licznik ponowień zadania prowadzony jest przez Koordynatora i widoczny w oknie pętli |
| Limity i budżet wywołań | Liczba wywołań, liczba ponowień pojedynczego zadania, liczba wywołań prowadzonych równolegle oraz budżet wywołań w ramach jednego zlecenia są ustawieniami konfiguracyjnymi; osiągnięcie wartości granicznej wstrzymuje pętlę i wystawia w oknie komunikat wymagający decyzji Użytkownika |

Ustawienia konfiguracyjne sterujące wywołaniami w ramach jednego zlecenia zestawia poniższa tabela.

> **Zapis kluczy nastaw.** Nazwy w postaci `obszar.grupa.nastawa` użyte w tym rozdziale są
> **kluczami konfiguracji**, nie komendami kontraktu. Klucz wskazuje miejsce wartości w modelu
> konfiguracji; komendy kontraktu, którymi się go odczytuje i zapisuje, to `config.get` i
> `config.set` (obszar `config` w `budowa/shared/contract.json`).


| Ustawienie konfiguracyjne | Znaczenie | Zakres obowiązywania |
|---|---|---|
| `model.call.max` | Ograniczenie liczby wywołań modelu w ramach jednego zlecenia. **Domyślnie brak limitu** — klucz służy Operatorowi do nałożenia ograniczenia, jeśli tego chce | Zlecenie |
| `model.retry.max` | Ograniczenie liczby ponowień pojedynczego zadania po odrzuceniu w kontroli jakości. **Domyślnie brak limitu** — ponowienia nie mają granicy narzuconej z góry | Zadanie |
| `model.parallel.max` | Ograniczenie liczby wywołań prowadzonych równolegle w jednej pętli wykonawczej. **Domyślnie brak limitu** | Zlecenie |
| `model.budget.order` | Budżet wywołań przypisany zleceniu, rozliczany łącznie dla wszystkich kanałów modeli biorących udział w pętli | Zlecenie |
| `model.timeout` | Czas oczekiwania na odpowiedź dostawcy, po którym adapter zamyka wywołanie i zgłasza je Koordynatorowi jako nieudane | Wywołanie |
| `model.fallback` | Kanał modelu zapasowego, do którego adapter kieruje wywołanie po niepowodzeniu kanału podstawowego | Kanał modelu |

Wszystkie powyższe ustawienia podlegają warstwowości konfiguracji (rozdz. 8.3): wartość warstwy domyślnej obowiązuje do chwili nałożenia wartości warstwy sesji, a brak jawnego ustawienia oznacza wartość domyślną.

**Warstwa modeli nie narzuca limitów.** Dla kluczy `model.call.max`, `model.retry.max` i `model.parallel.max` wartością domyślną jest brak limitu: platforma nie przerywa pętli po osiągnięciu z góry przyjętej liczby wywołań, ponowień ani wywołań równoległych. Klucz istnieje wyłącznie po to, by Operator mógł nałożyć ograniczenie, jeśli tego chce — nie po to, by je z góry ustanawiać. Ponieważ ponowienia nie mają granicy, przebieg nie ma stanu „wyczerpania prób": przerwanie następuje **na polecenie Operatora**, wydane z Execution Loop Window albo z Chat Window, i musi zadziałać także w trakcie trwającego ponowienia — Koordynator przyjmuje polecenie przerwania między próbami oraz w czasie odczekiwania narastającego odstępu, nie zwlekając z nim do końca bieżącej próby.

---

## 10. Warstwy widoczności funkcji integracji modeli

Funkcje warstwy integracji modeli podlegają regule stopniowego ujawniania funkcjonalności: w stanie spoczynku interfejsu widoczna jest wyłącznie informacja o aktywnym modelu i aktywnym wykonawcy, a pełny aparat konfiguracji pozostaje ukryty do chwili wystąpienia potrzeby jego użycia.

| Element interfejsu | Warstwa | Sposób dostępu | Zachowanie |
|---|---|---|---|
| Znacznik kontekstowy modelu i wykonawcy w postaci `[Fable 5] [Ultra]` | 2 | Widoczny w pasku kontekstu okna komunikacji; kliknięcie znacznika otwiera selektor wartości tego znacznika | Po dokonaniu wyboru selektor zwija się do znacznika |
| Menu progresywne `Model ▼` | 2 | Kliknięcie zwiniętego elementu rozwija listę skonfigurowanych kanałów modeli | Zwija się samoczynnie po wyborze |
| Menu progresywne `Agent ▼` | 2 | Kliknięcie zwiniętego elementu rozwija listę agentów i modeli bazowych dostępnych jako wykonawcy | Zwija się samoczynnie po wyborze |
| Parametry próbkowania: temperatura, top-p, długość odpowiedzi | 3 | Menu kontekstowe znacznika modelu, panel popover | Panel znika po zamknięciu |
| Limity i budżet wywołań: `model.call.max`, `model.retry.max`, `model.parallel.max`, `model.budget.order`, `model.timeout` | 3 | Menu kontekstowe okna pętli wykonawczej, wyszukiwarka funkcji, skrót klawiszowy | Panel znika po zamknięciu |
| Model zapasowy `model.fallback` | 3 | Menu kontekstowe znacznika modelu, wyszukiwarka funkcji | Panel znika po zamknięciu |
| Konfiguracja dostawców: typ kanału, punkt końcowy, dane dostępowe, nazwa modelu bazowego | 4 | Okno konfiguracji (rozdz. 8) otwierane z wyszukiwarki funkcji, skrótem klawiszowym, poleceniem języka naturalnego w oknie komunikacji lub w trybie administracyjnym | Niewidoczna dla użytkownika podstawowego |
| Adaptery kanałów i diagnostyka połączenia z dostawcą | 4 | Tryb administracyjny, narzędzia diagnostyczne | Niewidoczne dla użytkownika podstawowego |

Każdy z powyższych elementów osiągalny jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Poniższa makieta przedstawia stan spoczynku interfejsu — widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw wyższych.

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │                          │ (rozszerzenie
            │ [Fable 5] [Ultra] ⋮  │                          │  boczne)
            │ Model ▼  Agent ▼     │                          │
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
            │ zadania · ponowienia │                          │
            │ budżet wywołań    ⋮  │                          │
 ════════════════════════════════════════════════════════════════════════════
```

---

## 11. Relacja z izolacją techniczną

Okno konfiguracji punktów izolacji (rozdz. 6.4–6.6 koncepcji platformy) obejmuje osiem zakresów izolacji technicznej ustalonych w rozdziale 12 dokumentu Architektury. Trzy z nich dotyczą bezpośrednio warstwy integracji modeli opisanej w niniejszym dokumencie.

| Zakres izolacji technicznej | Relacja z integracją modeli | Stan domyślny |
|---|---|---|
| Katalog danych i konfiguracji modelu | Określa, czy dane i konfiguracja modelu przypisanego do sesji lub roli są przechowywane w katalogu odrębnym, czy współdzielonym z innymi sesjami lub rolami | Nieaktywny |
| Konto i token per sesja | Określa, czy dane dostępowe kanału modelu (klucz API, token CLI, dane SSH, dane logowania HTTP) są odrębne dla każdej sesji, czy współdzielone | Nieaktywny |
| Model procesu | Określa, czy proces obsługujący wywołania modelu jest odrębny dla danej sesji lub roli, czy współdzielony z innymi | Nieaktywny |

Podobnie jak wszystkie zakresy izolacji technicznej, powyższe trzy zakresy są domyślnie nieaktywne (rozdz. 6.6 koncepcji platformy oraz rozdz. 12 dokumentu Architektury) — kanał modelu przypisany do sesji lub roli działa bez dodatkowych ograniczeń technicznych, dopóki użytkownik świadomie nie włączy wybranego zakresu. Włączenie tych zakresów jest szczególnie istotne w środowisku MultitaskingAI, gdzie kilka ról korzysta z kanałów modeli równolegle — zakres „konto i token per sesja” zapewnia, że Executor 1 i Executor 2 korzystają z odrębnych danych dostępowych, nawet jeśli oba kanały modelu wskazują tego samego dostawcę. Profil izolacji (rozdz. 6.6 koncepcji platformy, Załącznik D.4) może grupować powyższe zakresy razem z pozostałymi ustawieniami izolacji kontekstu i izolacji technicznej oraz być przypisany bezpośrednio do roli — pole „Profil izolacji” w szablonie roli (Załącznik D.2 koncepcji platformy) odwołuje się wprost do tego mechanizmu.

---

## 12. Relacja z modułem Agents

Moduł Agents (rozdz. 11.15 koncepcji platformy) jest miejscem, w którym powstaje agent — trwały, nazwany komponent własny łączący model bazowy z pełną konfiguracją tożsamości, instrukcji systemowych, umiejętności, wtyczek, konektorów, pamięci i uprawnień. Budowa agenta w oknie Agent Builder obejmuje wybór modelu bazowego w oknie Model Configuration, które operuje bezpośrednio na warstwie dostawcy modelu: wybór typu kanału, dostawcy i modelu dokonywany przy tworzeniu agenta jest tym samym mechanizmem konfiguracji kanału modelu opisanym w rozdziale 8. Relację między agentem a kanałem modelu porządkuje poniższe zestawienie.

| Poziom | Co zawiera | Gdzie się konfiguruje | Trwałość |
|---|---|---|---|
| Kanał modelu | Typ kanału, dane dostępowe, nazwa modelu bazowego | Okno konfiguracji (rozdz. 8) lub Model Configuration w Agent Builder | Konfiguracja techniczna, wykorzystywana przez wyższe poziomy |
| Tożsamość i persona modelu | Nazwa dowolna, prompt systemowy | Okno konfiguracji, przy przypisaniu do sesji lub roli (rozdz. 6) | Właściwa danemu przypisaniu; nie tworzy trwałego komponentu własnego |
| Agent | Model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć, uprawnienia | Moduł Agents, Agent Builder (rozdz. 11.15 koncepcji platformy) | Trwały, nazwany komponent własny, wielokrotnego użycia w dowolnym środowisku |

Poniższy schemat przedstawia te trzy poziomy jako hierarchię — od najniższego (technicznego) do najwyższego (trwałego komponentu własnego).

```
Poziomy konfiguracji modelu (od najniższego do najwyższego):

  Kanał modelu            typ kanału · dane dostępowe · nazwa modelu bazowego
       │                  (podstawa techniczna)
       ▼
  Tożsamość i persona     nazwa nadana · prompt systemowy
       │                  (dla danego przypisania)
       ▼
  Agent                   model bazowy · tożsamość · instrukcje systemowe
                          · umiejętności · wtyczki · konektory · pamięć · uprawnienia
                          (trwały, nazwany komponent własny)
```

Agent, po utworzeniu, jest komponentem platformowym działającym ponad wszystkimi środowiskami (rozdz. 11.15 koncepcji platformy) — może zostać wykorzystany jako wykonawca zadań w TalkIn, WorkSpace i CodeStudio oraz jako wykonawca roli w środowisku MultitaskingAI. W obu przypadkach agent niesie ze sobą własny kanał modelu, ustalony przy jego budowie, niezależnie od kanałów modeli przypisanych bezpośrednio do innych sesji lub ról tej samej karty lub tego samego procesu.

---

## 13. Bezpieczeństwo danych dostępowych

Dane dostępowe kanału modelu — klucz dostępu API, token narzędzia wiersza poleceń CLI, dane dostępowe SSH, dane logowania interfejsu webowego HTTP — przechowywane są poza bazą danych platformy, zgodnie z zasadami bezpieczeństwa ustalonymi w rozdziale 1.5 dokumentu Modelu danych. Baza SQLite (rozdz. 10 dokumentu Architektury) przechowuje wyłącznie odwołanie do danych dostępowych w ramach encji kanału modelu (rozdz. 5), nigdy same dane w postaci jawnej. Poniższy schemat przedstawia to rozdzielenie.

```
        Baza danych SQLite (rozdz. 10 dokumentu Architektury)
        └── encja: Kanał modelu
                └── Odwołanie do danych dostępowych ──┐
                                                      │  (tylko wskazanie, nigdy dane jawne)
                                                      ▼
        Przechowywanie poza bazą (rozdz. 1.5 dokumentu Modelu danych)
        └── Dane dostępowe w postaci właściwej kanałowi:
                klucz dostępu API · token CLI · dane SSH · dane logowania HTTP
                                                      ▲
                                                      │  pobranie w momencie wywołania
        Adapter kanału (rozdz. 4) ─────────────────────┘
```

Konsekwencje tej zasady dla warstwy integracji modeli zestawia poniższa tabela.

| Konsekwencja | Treść |
|---|---|
| Pobranie danych przez adapter | Adapter kanału (rozdz. 4) pobiera dane dostępowe w momencie wywołania modelu, na podstawie odwołania przechowanego w bazie, a nie z samej bazy wprost |
| Eksport i kopia bazy | Eksport, kopia zapasowa lub podgląd zawartości bazy danych nie ujawnia kluczy, tokenów ani danych logowania przypisanych kanałom modeli |
| Usunięcie kanału modelu | Usunięcie kanału modelu z bazy nie pozostawia danych dostępowych w formie odzyskiwalnej z warstwy danych platformy |

Zasada ta jest spójna z ogólną zasadą architektoniczną „Pojedyncze źródło prawdy” (rozdz. 1 dokumentu Architektury) w tym sensie, że serwer pozostaje jedynym miejscem przechowywania odwołania do danych dostępowych — rozdzielenie samych danych od ich odwołania jest doprecyzowaniem tej zasady na poziomie bezpieczeństwa, a nie odstępstwem od niej.

---

## 14. Zgodność z zasadami nadrzędnymi

Warstwa integracji modeli jest podporządkowana zasadom nadrzędnym funkcjonalności platformy (rozdz. 14 koncepcji platformy) w sposób zestawiony poniżej.

| Zasada nadrzędna | Zastosowanie w warstwie integracji modeli |
|---|---|
| Pełna kompozycyjność | Kanał modelu, tożsamość i persona modelu oraz agent (rozdz. 5, 6, 12) są komponentami swobodnie zestawianymi: model przypisany bezpośrednio do roli, agent zbudowany na innym kanale lub połączenie obu podejść dla różnych ról tego samego zespołu w środowisku MultitaskingAI |
| Pełna konfigurowalność | Wybór kanału, danych dostępowych, modelu, tożsamości i przypisania dokonywany jest w całości z poziomu okna konfiguracji (rozdz. 8); żaden kanał, model ani sposób przypisania nie jest zaszyty na stałe w sposób niedostępny dla użytkownika |
| Jawność i konfigurowalność zależności | Relacja między kanałem modelu a izolacją techniczną (rozdz. 11) oraz relacja między kanałem modelu a agentem (rozdz. 12) są jawne i konfigurowalne; żadna nie jest wymuszana przez platformę jako zależność wbudowana na stałe |
| Rozszerzenie orkiestracji | Mechanizm adaptera (rozdz. 4) utrzymuje liczbę dostępnych kanałów integracji niezależną od rdzenia, a pętla wykonawcza (rozdz. 9.2) prowadzi wywołania wielu modeli i agentów równolegle (rozdz. 14 koncepcji platformy) |

Zgodnie z zasadą zera blokad twardych obowiązującą w projekcie, okno konfiguracji nigdy nie wymusza wyboru konkretnego kanału ani dostawcy modelu — udostępnia wszystkie cztery kanały jako równorzędne możliwości, a brak jawnego ustawienia na danym poziomie skutkuje dziedziczeniem konfiguracji z warstwy domyślnej (rozdz. 8), zgodnie z ogólną regułą „brak ustawienia = wartość domyślna”.

---

## 15. Wykaz komend kontraktu — obszary `channel`, `identity` i `model`

Trzy obszary kontraktu komunikacji (`budowa/shared/contract.json`) niosą komendy warstwy dostawcy modelu opisanej w rozdziałach 2–8: `channel` (rejestr kanałów, sześć komend), `identity` (tożsamość i persona modelu, pięć komend) oraz `model` (przypisanie kanału do okna albo karty sesji, jedna komenda). Do zestawienia dołączono `agent.model.set` z obszaru `agent`, ponieważ ustala model bazowy i parametry wywołania eksperta budowanego w module Agents (rozdz. 12). Poniższy wykaz jest pełny dla trzynastu komend. Czwarty obszar warstwy modeli — `account` (konta modeli i konta programów wiersza poleceń, sześć komend i jedno zdarzenie; struktura `Account` opisana w rozdz. 15.3) — wraz z pełnym brzmieniem wszystkich osiemnastu komend warstwy niesie Załącznik C.

| Komenda | Obszar | Pola żądania (nazwa: typ, wymagalność) | Pola wyniku |
|---|---|---|---|
| `channel.add` | `channel` | `name: string, wymagane` · `kind: string, wymagane` · `model: string, opcjonalne` · `enabled: bool, opcjonalne` · `config: json, opcjonalne` | `channel: Channel, wymagane` |
| `channel.update` | `channel` | `channelId: string, wymagane` · `name, model, enabled, config: opcjonalne` | `channel: Channel, wymagane` |
| `channel.remove` | `channel` | `channelId: string, wymagane` | `channelId: string, wymagane` |
| `channel.list` | `channel` | `enabledOnly: bool, opcjonalne` | `channels: Channel[], wymagane` |
| `channel.check` | `channel` | `channelId: string, wymagane` · `transport: ProviderTransport, opcjonalne` | `reachable: bool, wymagane` · `checkedAt: int64, wymagane` · `latencyMs: int, opcjonalne` · `detail: string, opcjonalne` |
| `channel.credential.status` | `channel` | `channelId: string, wymagane` | `status: ChannelCredentialStatus, wymagane` |
| `identity.category.list` | `identity` | `layer: IdentityLayer, opcjonalne` · `includeDisabled: bool, opcjonalne` | `categories: IdentityCategory[], wymagane` |
| `identity.document.get` | `identity` | `categoryId: string, opcjonalne` · `axis: ConfigAxis, opcjonalne` · `axisId: string, opcjonalne` | `documents: IdentityDocument[], wymagane` |
| `identity.document.set` | `identity` | `categoryId: string, wymagane` · `content: string, wymagane` · `mode: IdentityMode, opcjonalne` · `axis: ConfigAxis, opcjonalne` · `axisId: string, opcjonalne` · `enabled: bool, opcjonalne` | `document: IdentityDocument, wymagane` |
| `identity.document.remove` | `identity` | `documentId: string, wymagane` | `removed: bool, wymagane` |
| `identity.effective.get` | `identity` | `windowId: string, opcjonalne` · `axis: ConfigAxis, opcjonalne` · `axisId: string, opcjonalne` | `mode: IdentityMode, wymagane` · `layers: IdentityLayerContent[], wymagane` · `prompt: string, wymagane` · `promptHash: string, wymagane` · `missingRequiredCategoryIds: string[], opcjonalne` |
| `model.channel.set` | `model` | `modelChannelId: string, wymagane` · `windowId: string, opcjonalne` · `sessionId: string, opcjonalne` · `agentId: string, opcjonalne` | `channel: Channel, wymagane` · `windowId: string, opcjonalne` · `sessionId: string, opcjonalne` |
| `agent.model.set` | `agent` | `agentId: string, wymagane` · `channelId: string, wymagane` · `model: string, opcjonalne` · `transport: ProviderTransport, opcjonalne` · `parameters: json, opcjonalne` · `effort: ReasoningEffort, opcjonalne` | `agent: Agent, wymagane` |

Zdarzeń zwrotnych ani kodów błędów właściwych wyłącznie tym trzynastu komendom kontrakt nie wyodrębnia jako osobnej rodziny — obowiązuje wspólny zestaw dziewięciu kodów błędów kontraktu, identyczny dla wszystkich obszarów ([Kontrakty komunikacji](kontrakty-komunikacji.md) rozdz. 19, Załącznik D). Pole `kodyBledow` w `budowa/shared/contract.json` zna dziś osiem kodów; dopisanie dziewiątego, `command_not_understood`, jest osobnym zadaniem w kodzie.

### 15.1. Struktury i wyliczenia normatywne

| Byt | Rodzaj | Pola albo wartości |
|---|---|---|
| `Channel` | struktura | `id, name, kind: string` · `model, accountId: string, opcjonalne` · `enabled: bool` · `config: json, opcjonalne` · `createdAt: int64` |
| `ChannelCredentialStatus` | struktura | `channelId: string` · `present: bool` · `kind: string, opcjonalne` · `updatedAt: int64, opcjonalne` · `managedBy: string, opcjonalne` — struktura celowo nie niesie pola na wartość poświadczenia |
| `IdentityCategory` | struktura | `id, name: string` · `description: string, opcjonalne` · `layer: IdentityLayer` · `order: int` · `required: bool` · `defaultMode: IdentityMode` · `enabled: bool` |
| `IdentityDocument` | struktura | `id, categoryId: string` · `axis: ConfigAxis` · `axisId: string, opcjonalne` · `mode: IdentityMode` · `content: string` · `contentHash: string, opcjonalne` · `enabled: bool` · `updatedAt: int64` |
| `IdentityLayerContent` | struktura | `layer: IdentityLayer` · `categoryId, name, content: string` · `order: int` · `axis: ConfigAxis` · `axisId: string, opcjonalne` |
| `ProviderTransport` | wyliczenie | `cli` (program wiersza poleceń lokalny) · `api` (punkt końcowy HTTP dostawcy) · `sdk` (biblioteka dostawcy wołana z rdzenia) · `ssh` (program wiersza poleceń na maszynie zdalnej) · `local` (model uruchamiany na urządzeniu bez pośrednictwa dostawcy) |
| `ReasoningEffort` | wyliczenie | `low` (nakład najniższy) · `medium` (nakład pośredni) · `high` (nakład najwyższy) |
| `ConfigAxis` | wyliczenie | `platform` (domyślna, ustawienie niezależne od modelu i konta) · `model` (ustawienie dla wskazanego modelu) · `account` (ustawienie dla wskazanego konta z rejestru kont) |
| `IdentityLayer` | wyliczenie | `constitution` (warstwa najwyższa) · `profile` (profil roli) · `expertise` (warstwa najniższa) — kolejność krytyczności złożonego promptu systemowego |
| `IdentityMode` | wyliczenie | `ZASTAP` (tożsamość zastępuje prompt fabryczny w całości; wartość domyślna klucza `tozsamosc.tryb_domyslny`) · `DOLACZ` (tożsamość dopisuje się do promptu fabrycznego) |

### 15.2. Diagram sekwencji — zapis kanału i odczyt tożsamości obowiązującej

```
Operator            Klient (okno konfiguracji)         Rdzeń
   │                       │                              │
   │  dane kanału          │                              │
   │  (name, kind, model)  │                              │
   ├──────────────────────►│                              │
   │                       │  channel.add                 │
   │                       ├─────────────────────────────►│
   │                       │  channel: Channel            │
   │                       │◄─────────────────────────────┤
   │  treść kategorii      │                              │
   │  zasad (konstytucja,  │                              │
   │  profil, ekspertyza)  │                              │
   ├──────────────────────►│  identity.document.set       │
   │                       │  (categoryId, content, mode) │
   │                       ├─────────────────────────────►│
   │                       │  document: IdentityDocument  │
   │                       │◄─────────────────────────────┤
   │  wybór kanału         │                              │
   │  dla okna             │                              │
   ├──────────────────────►│  model.channel.set           │
   │                       ├─────────────────────────────►│
   │                       │  channel, windowId           │
   │                       │◄─────────────────────────────┤
   │                       │  identity.effective.get      │
   │                       │  (windowId)                  │
   │                       ├─────────────────────────────►│
   │                       │                              │  składanie warstw wg
   │                       │                              │  IdentityLayer: constitution
   │                       │                              │  → profile → expertise
   │                       │  mode, layers, prompt,       │
   │                       │  promptHash                  │
   │                       │◄─────────────────────────────┤
   │  model gotowy         │                              │
   │  do wywołania         │                              │
   │◄──────────────────────┤                              │
```

Legenda: `identity.effective.get` składa złożony prompt systemowy z warstw `IdentityLayer` w kolejności krytyczności (`constitution` przed `profile` przed `expertise`); pole `missingRequiredCategoryIds` sygnalizuje kategorie wymagane bez treści, lecz brak treści nie wstrzymuje uruchomienia modelu.

### 15.3. Struktury danych pełne

Rozdział 15.1 przenosi pola struktur w postaci skróconej. Poniższy wykaz podaje pełną definicję każdego pola, wraz z wymagalnością (`wymagane` struktury źródłowej w `budowa/shared/contract.json`), dla struktur używanych przez trzynaście komend rozdziału 15.

**`Channel` — wiersz rejestru kanałów modelu:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator kanału |
| `name` | `string` | tak | Nazwa kanału |
| `kind` | `string` | tak | Rodzaj kanału; wartość danych, nie typ kodu |
| `model` | `string` | nie | Identyfikator modelu |
| `accountId` | `string` | nie | Konto z rejestru kont, którego poświadczenia obsługują kanał — kanał mówi CZYM się połączyć, konto CZYIM poświadczeniem |
| `enabled` | `bool` | tak | Czy kanał czynny |
| `config` | `json` | nie | Parametry kanału; odwołania do danych dostępowych, nigdy ich treść |
| `createdAt` | `int64` | tak | Czas utworzenia w milisekundach epoki |

**`Account` — wiersz rejestru kont modeli i kont code CLI:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator konta |
| `name` | `string` | tak | Nazwa konta |
| `kind` | `AccountKind` | tak | Rodzaj konta — `api`, `cli` albo `sdk` |
| `provider` | `string` | tak | Dostawca; wartość danych, nie typ kodu |
| `externalId` | `string` | nie | Identyfikator konta u dostawcy |
| `defaultModel` | `string` | nie | Model domyślny konta |
| `baseUrl` | `string` | nie | Adres punktu końcowego dostawcy |
| `configDir` | `string` | nie | Katalog konfiguracji programu code CLI |
| `hasCredential` | `bool` | tak | Czy poświadczenie jest zapisane — odpowiedź mówi tylko tyle, nigdy więcej |
| `isDefault` | `bool` | tak | Czy konto domyślne swojego rodzaju |
| `enabled` | `bool` | tak | Czy konto czynne |
| `order` | `int` | tak | Kolejność konta na wykazie, liczona od 1 |
| `createdAt` / `updatedAt` | `int64` | tak | Czas utworzenia i czas ostatniej zmiany w milisekundach epoki |

**`IdentityCategory` — wiersz katalogu kategorii zasad:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Kod kategorii |
| `name` | `string` | tak | Nazwa kategorii w oknie konfiguracji |
| `description` | `string` | nie | Opis przeznaczenia kategorii |
| `layer` | `IdentityLayer` | tak | Warstwa nakładki, do której kategoria wchodzi |
| `order` | `int` | tak | Kolejność kategorii w obrębie warstwy, liczona od 1 |
| `required` | `bool` | tak | Czy kategoria musi mieć treść |
| `defaultMode` | `IdentityMode` | tak | Tryb proponowany przy zapisie treści |
| `enabled` | `bool` | tak | Czy wiersz katalogu czynny |

**`IdentityDocument` — zapis treści kategorii dla wskazanej osi:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator zapisu |
| `categoryId` | `string` | tak | Kategoria zasad |
| `axis` | `ConfigAxis` | tak | Oś, dla której treść obowiązuje |
| `axisId` | `string` | nie | Byt osi: identyfikator modelu albo konta; puste dla osi `platform` |
| `mode` | `IdentityMode` | tak | Tryb podania treści |
| `content` | `string` | tak | Treść kategorii |
| `contentHash` | `string` | nie | Odcisk treści; ten sam odcisk znaczy tę samą nakładkę |
| `enabled` | `bool` | tak | Czy zapis czynny |
| `updatedAt` | `int64` | tak | Czas ostatniej zmiany w milisekundach epoki |

**`IdentityLayerContent` — warstwa złożonego promptu systemowego w wyniku `identity.effective.get`:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `layer` | `IdentityLayer` | tak | Warstwa nakładki |
| `categoryId` | `string` | tak | Kategoria zasad, z której pochodzi treść |
| `name` | `string` | tak | Nazwa kategorii |
| `content` | `string` | tak | Treść warstwy |
| `order` | `int` | tak | Kolejność warstwy w nakładce, liczona od 1 |
| `axis` | `ConfigAxis` | tak | Oś, z której treść została wzięta |
| `axisId` | `string` | nie | Byt osi; puste dla osi `platform` |

**`ChannelCredentialStatus` — stan poświadczenia kanału, bez treści poświadczenia:**

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `channelId` | `string` | tak | Kanał, którego poświadczenie opisano |
| `present` | `bool` | tak | Czy poświadczenie jest ustawione |
| `kind` | `string` | nie | Rodzaj poświadczenia w brzmieniu właściwym kanałowi — klucz API albo token programu wiersza poleceń; puste, gdy poświadczenia nie ma |
| `updatedAt` | `int64` | nie | Czas ostatniej zmiany poświadczenia w milisekundach epoki; puste, gdy poświadczenia nie ma |
| `managedBy` | `string` | nie | Skąd poświadczenie pochodzi: konto platformy albo zapis własny kanału — Operator ma wiedzieć, gdzie je zmienić |

### 15.4. Kody błędów zastosowane do warstwy dostawcy modelu

Kody błędów kontraktu są wspólne dla całego kontraktu ([Architektura techniczna](architektura.md) rozdz. 19; [Kontrakty komunikacji](kontrakty-komunikacji.md) Załącznik D) i liczą dziewięć pozycji, z których pole `kodyBledow` w `budowa/shared/contract.json` zna osiem — dziewiąty kod `command_not_understood` czeka na wdrożenie w kodzie osobnym zadaniem. Poniższa tabela wskazuje zastosowanie ośmiu kodów wdrożonych do komend obszarów `channel`, `identity`, `model` i do `agent.model.set`.

| Kod | Ponawialny | Zastosowanie w warstwie dostawcy modelu |
|---|---|---|
| `validation_failed` | nie | `kind` poza wartością obsługiwaną przez rejestr kanałów w `channel.add`; `axis` niezgodna z `ConfigAxis` w `identity.document.set` |
| `not_found` | nie | `channelId` bez odpowiadającego wiersza rejestru; `categoryId` nieznany w `identity.document.set`; `accountId` nieznany w `account.update` |
| `not_authenticated` | nie | Każda komenda obszaru wymaga tokenu dostępu ważnego dla konta właściciela (Bezpieczeństwo i uwierzytelnianie, rozdz. 8) |
| `permission_denied` | nie | Rezerwowany dla przyszłego rozróżnienia ról przy współdzielonym dostępie do rejestru kont — dziś rejestr jest jednowłaścicielski |
| `conflict` | nie | `account.default.set` na rodzaj konta bez żadnego konta założonego; usunięcie kategorii wymaganej (`required: true`) bez zastępczej treści |
| `channel_unavailable` | tak | `channel.check` i wywołanie modelu przez kanał chwilowo nieosiągalny (rozdz. 6) |
| `rate_limited` | tak | Wywołanie modelu przez kanał przy przekroczeniu limitu tempa dostawcy |
| `internal_error` | tak | Dowolna komenda obszaru przy awarii rdzenia niezwiązanej z treścią żądania |

---

## 16. Rozbieżność modelu czterech kanałów z kontraktem

Rozdziały 2–4 niniejszego dokumentu opisują warstwę dostawcy modelu jako cztery kanały integracji — API, CLI, SSH, HTTP. Kontrakt komunikacji niesie w tym miejscu dwie wartości niosące pokrewne, lecz nietożsame rozróżnienie:

| Element kontraktu | Wartości | Relacja do modelu „cztery kanały” |
|---|---|---|
| Pole `kind` struktury `Channel` | typ `string` — „wartość danych, nie typ kodu” (opis pola w `contract.json`) | Rodzaj kanału nie jest zamkniętym wyliczeniem czterech wartości w warstwie danych; jest polem otwartym, wypełnianym przez okno konfiguracji |
| Wyliczenie `ProviderTransport` | `cli` · `api` · `sdk` · `ssh` · `local` (pięć wartości) | Cztery z pięciu wartości pokrywają się częściowo z API/CLI/SSH (jako `api`, `cli`, `ssh`); wartości `sdk` i `local` nie mają odpowiednika w zestawie API/CLI/SSH/HTTP, a kanał HTTP (model przeglądarkowy) nie ma odpowiednika wprost w `ProviderTransport` |

**[DO DECYZJI OPERATORA]** Czy pole `kind` encji `Channel` ma być w praktyce ograniczone do czterech wartości „api” | „cli” | „ssh” | „http” zgodnie z modelem opisanym w rozdziałach 2–4, czy model czterech kanałów jest uproszczeniem koncepcyjnym, a rzeczywisty dobór drogi wywołania dostawcy prowadzi wyłącznie wyliczenie `ProviderTransport` (`cli`, `api`, `sdk`, `ssh`, `local`) — a jeśli tak, jaka jest relacja między kanałem przeglądarkowym HTTP opisanym w rozdz. 3.1 a wartościami `sdk` i `local` tego wyliczenia. Rozstrzygnięcie warunkuje treść rozdz. 3–4 niniejszego dokumentu oraz Załącznika A w kolejnym wydaniu.

Poniższy schemat porządkuje wizualnie, które wartości obu stron rozbieżności mają jednoznaczny odpowiednik, a które pozostają nierozstrzygnięte.

```
Model koncepcyjny (rozdz. 2–4)        Wyliczenie ProviderTransport (kontrakt)
┌──────────────────┐                   ┌───────────────────┐
│  API             │ ───── zgodność ──►│  api              │
│  CLI             │ ───── zgodność ──►│  cli              │
│  SSH             │ ───── zgodność ──►│  ssh              │
│  HTTP            │ ─── [DO DECYZJI  ►│  sdk              │
│                  │      OPERATORA]   │  local            │
└──────────────────┘                   └───────────────────┘
                                        pole `kind` struktury `Channel`
                                        jest `string` otwarty — żadna
                                        z dwóch kolumn nie jest wymuszona
                                        przez typ danych kontraktu
```

Legenda: strzałka ciągła oznacza jednoznaczne odwzorowanie nazwy; strzałka opisana `[DO DECYZJI OPERATORA]` prowadzi do dwóch wartości (`sdk`, `local`) bez ustalonego odpowiednika po stronie modelu koncepcyjnego czterech kanałów, oraz od kanału HTTP bez ustalonego odpowiednika po stronie wyliczenia `ProviderTransport`.

---

## 17. Słowniczek pojęć

| Pojęcie | Znaczenie |
|---|---|
| Kanał modelu | Skonfigurowane połączenie z konkretnym modelem, reprezentowane strukturą `Channel`, przypisywalne per sesja albo per rola, niezależnie od agenta (rozdz. 5, rozdz. 7, rozdz. 15.3). |
| Adapter kanału | Komponent rdzenia tłumaczący jednolite zapytanie warstwy dostawcy modelu na protokół i uwierzytelnienie właściwe jednemu z kanałów integracji (rozdz. 4). |
| Konto | Wiersz rejestru kont, reprezentowany strukturą `Account`, niosący dane dostawcy i wskaźnik obecności poświadczenia (`hasCredential`), nigdy samo poświadczenie (rozdz. 13, rozdz. 15.3, Załącznik B.4). |
| Tożsamość modelu | Zestaw kategorii zasad (`IdentityCategory`) o treści zapisanej dla wskazanej osi (`IdentityDocument`), składany w złożony prompt systemowy przez `identity.effective.get` (rozdz. 6, rozdz. 15). |
| Persona modelu | Nazwa dowolna nadana modelowi przez Operatora, odróżniająca dwa przypisania korzystające z tego samego kanału i modelu bazowego (rozdz. 6, Załącznik B.3). |
| Warstwa nakładki (`IdentityLayer`) | Jeden z trzech poziomów krytyczności złożonego promptu systemowego: `constitution`, `profile`, `expertise` — w tej kolejności rozstrzygania (rozdz. 15.1–15.2). |
| Oś konfiguracji (`ConfigAxis`) | Wymiar, dla którego treść tożsamości obowiązuje: `platform` (domyślna), `model` albo `account` — prostopadły do warstwy nakładki (rozdz. 15.1). |
| Tryb podania (`IdentityMode`) | Sposób łączenia treści tożsamości z promptem fabrycznym modelu: `ZASTAP` (zastąpienie w całości) albo `DOLACZ` (dopisanie) — wartości dosłowne wyliczenia kontraktu (rozdz. 15.1). |
| Agent | Trwały, nazwany komponent własny łączący model bazowy z pełną konfiguracją tożsamości, umiejętności, konektorów, pamięci i uprawnień, budowany w module Agents (rozdz. 12). |
| Rozbieżność modelu czterech kanałów | Niespójność między modelem koncepcyjnym API/CLI/SSH/HTTP a wyliczeniem `ProviderTransport` kontraktu, oznaczona **[DO DECYZJI OPERATORA]** (rozdz. 16). |

---

## 18. Kryteria odbioru

Warstwa integracji modeli jest zgodna z niniejszym opracowaniem, gdy spełnione są łącznie poniższe warunki, sprawdzalne bez odwołania do intencji autora.

| Warunek | Sposób sprawdzenia |
|---|---|
| Rdzeń operuje na jednolitej abstrakcji „wyślij zapytanie, odbierz odpowiedź”, niezależnej od kanału | przegląd kodu warstwy dostawcy modelu — brak gałęzi logiki sesji rozróżniającej kanał poza warstwą adaptera |
| Każdy z czterech kanałów integracji realizowany jest przez odrębny adapter (rozdz. 4) | test integracyjny wywołania przez każdy adapter zwraca odpowiedź w jednolitym formacie strumienia |
| Kanał modelu przechowuje wyłącznie odwołanie do danych dostępowych, nigdy ich treść (rozdz. 13) | odczyt encji `Channel` z bazy — pole `config` nie zawiera treści poświadczenia; `channel.credential.status` nie zwraca treści |
| Tożsamość modelu składa się z warstw `IdentityLayer` w kolejności krytyczności `constitution` → `profile` → `expertise` (rozdz. 15.1) | wywołanie `identity.effective.get` zwraca `layers` uporządkowane rosnąco wg pola `order` w obrębie warstwy |
| Przypisanie modelu działa niezależnie per sesja i per rola (rozdz. 7) | zmiana `model.channel.set` dla jednej karty sesji nie zmienia kanału innej karty ani innej roli |
| Rola środowiska MultitaskingAI przyjmuje zamiennie agenta albo model bazowy (rozdz. 7.2, rozdz. 12) | pole „Przypisany agent lub model” szablonu roli waliduje obie ścieżki bez błędu |
| Okno konfiguracji nie wymusza wyboru kanału ani dostawcy (rozdz. 14) | brak jawnego ustawienia kanału skutkuje dziedziczeniem z warstwy domyślnej, nie błędem |
| Rozbieżność rozdz. 16 pozostaje jawnie oznaczona do chwili rozstrzygnięcia przez Właściciela | obecność oznaczenia `[DO DECYZJI OPERATORA]` w rozdz. 16 do czasu wydania rozstrzygającego |
| Każda z dwunastu komend obszarów `channel`, `identity` i `model` oraz pozycja `agent.model.set` odpowiadają wykazowi rozdziału 15, a komplet osiemnastu komend warstwy — wykazowi Załącznika C, bez nazwy wymyślonej i bez pominięcia | `grep -oE '"typ": *"(channel\|identity)\.[a-zA-Z.]+"' budowa/shared/contract.json \| sort -u` zwraca dokładnie 11 pozycji zgodnych z rozdz. 15, uzupełnione o pojedyncze pozycje `model.channel.set` obszaru `model` i `agent.model.set` obszaru `agent` |
| Każde z pięciu pól struktury `ChannelCredentialStatus` (rozdz. 15.3) jest zgodne z definicją `budowa/shared/contract.json` co do nazwy, typu i wymagalności | porównanie wykazu pól rozdz. 15.3 z polem `struktury` kontraktu dla nazwy `ChannelCredentialStatus` nie ujawnia różnicy |

---

## Załącznik A. Szablon konfiguracji kanału modelu

Poniższe szablony porządkują pola konfiguracji kanału modelu opisane w niniejszym dokumencie; wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna”. Zapis `odwołanie → …` oznacza, że w encji kanału modelu przechowywane jest wyłącznie odwołanie do danych dostępowych, a same dane trafiają do przechowywania poza bazą (rozdz. 13); wartości w nawiasach ostrokątnych są miejscami do uzupełnienia.

### A.1. Pola konfiguracji kanału modelu

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Typ kanału | Sposób połączenia z modelem | API \| CLI \| SSH \| HTTP (rozdz. 3) |
| Dane dostępowe | Odwołanie do danych uwierzytelniających, przechowywanych poza bazą | klucz API \| token CLI \| dane SSH \| dane logowania HTTP (rozdz. 13) |
| Nazwa modelu | Techniczna nazwa modelu bazowego udostępnianego przez dostawcę | dowolny tekst zgodny z nazewnictwem dostawcy |
| Nazwa nadana (persona) | Dowolna nazwa nadana modelowi przez użytkownika | dowolny tekst; przykład: „Recenzent” |
| Prompt systemowy | Instrukcja wgrywana dla danego przypisania | dowolny tekst lub plik |
| Przypisanie | Element, do którego kanał modelu jest przypisany | sesja \| rola (rozdz. 7) |
| Warstwa | Warstwa konfiguracji | domyślna \| sesji (rozdz. 8 niniejszego dokumentu, rozdz. 13 dokumentu Architektury) |
| Profil izolacji powiązany | Profil izolacji technicznej obejmujący kanał modelu | odwołanie do profilu (rozdz. 11, Załącznik D.4 koncepcji platformy) |

### A.2. Szablon kanału API

```
Kanał modelu — API
  Typ kanału:          API
  Punkt końcowy:       <adres punktu końcowego dostawcy>
  Dane dostępowe:      odwołanie → klucz dostępu             (przechowywany poza bazą, rozdz. 13)
  Nazwa modelu:        <techniczna nazwa modelu bazowego u dostawcy>
  Nazwa nadana:        <nazwa persony — przykład: „Recenzent”>       (rozdz. 6)
  Prompt systemowy:    <instrukcja przypisania>                      (rozdz. 6)
  Przypisanie:         sesja | rola                                  (rozdz. 7)
  Warstwa:             domyślna | sesji                              (rozdz. 8)
  Profil izolacji:     <odwołanie do profilu>                        (rozdz. 11)
```

### A.3. Szablon kanału CLI

```
Kanał modelu — CLI
  Typ kanału:          CLI
  Narzędzie:           <nazwa narzędzia wiersza poleceń dostawcy>
  Lokalizacja:         lokalna | zdalna
  Dane dostępowe:      odwołanie → token narzędzia wiersza poleceń   (przechowywany poza bazą, rozdz. 13)
  Nazwa modelu:        <techniczna nazwa modelu udostępnianego przez narzędzie>
  Nazwa nadana:        <nazwa persony>                               (rozdz. 6)
  Prompt systemowy:    <instrukcja przypisania>                      (rozdz. 6)
  Przypisanie:         sesja | rola                                  (rozdz. 7)
  Warstwa:             domyślna | sesji                              (rozdz. 8)
  Profil izolacji:     <odwołanie do profilu>                        (rozdz. 11)
```

### A.4. Szablon kanału SSH

```
Kanał modelu — SSH
  Typ kanału:          SSH
  Host zdalny:         <adres hosta>
  Konto:               <konto na hoście zdalnym>
  Dane dostępowe:      odwołanie → dane dostępowe SSH (klucz lub hasło)   (przechowywane poza bazą, rozdz. 13)
  Model pomocniczy:    <nazwa własnego modelu pomocniczego na hoście>
  Nazwa nadana:        <nazwa persony>                                    (rozdz. 6)
  Prompt systemowy:    <instrukcja przypisania>                           (rozdz. 6)
  Przypisanie:         sesja | rola                                       (rozdz. 7)
  Warstwa:             domyślna | sesji                                   (rozdz. 8)
  Profil izolacji:     <odwołanie do profilu>                             (rozdz. 11)
```

### A.5. Szablon kanału HTTP

```
Kanał modelu — HTTP
  Typ kanału:          HTTP
  Interfejs webowy:    <adres interfejsu webowego dostawcy>
  Dane dostępowe:      odwołanie → dane logowania interfejsu webowego   (przechowywane poza bazą, rozdz. 13)
  Nazwa modelu:        <nazwa modelu przeglądarkowego>
  Nazwa nadana:        <nazwa persony>                                  (rozdz. 6)
  Prompt systemowy:    <instrukcja przypisania>                         (rozdz. 6)
  Przypisanie:         sesja | rola                                     (rozdz. 7)
  Warstwa:             domyślna | sesji                                 (rozdz. 8)
  Profil izolacji:     <odwołanie do profilu>                           (rozdz. 11)
```

---

## Załącznik B. Scenariusze użycia

Scenariusze ilustrują sposób wykorzystania czterech kanałów i mechanizmu tożsamości modelu oraz współdziałanie ustaleń opisanych w rozdziałach 3–10. Poniższa tabela zestawia trzy scenariusze; kolejne podpunkty przedstawiają je jako ponumerowane mini-przepływy.

| Scenariusz | Kanał | Miejsce | Wyróżnik |
|---|---|---|---|
| B.1 | API | Rola Executor 2 w środowisku MultitaskingAI | Model bazowy z personą zamiast pełnego agenta |
| B.2 | SSH | Sesja w module Diagnostics | Własny model pomocniczy na hoście zdalnym |
| B.3 | API | Dwa okna Model Panels w module Roundtable | Jeden kanał i jeden model bazowy, dwie tożsamości |
| B.4 | API | Sekcja Konta modeli Okna Ustawień | Rotacja poświadczenia konta bez rekonfiguracji kanałów powiązanych |

### B.1. Przypisanie modelu bezpośrednio do roli w MultitaskingAI

Konfiguracja roli Executor 2 w środowisku MultitaskingAI z bezpośrednim przypisaniem modelu bazowego, bez pełnego agenta z modułu Agents.

1. Wybór kanału API w oknie konfiguracji, z kluczem dostępu do wybranego dostawcy.
2. Wskazanie nazwy modelu bazowego.
3. Nadanie modelowi nazwy „Frontend” i wgranie promptu systemowego ukierunkowującego jego pracę na warstwę kliencką równoległego procesu, zgodnie z podziałem Executor 1 / Executor 2 (rozdz. 13.3 koncepcji platformy).
4. Przypisanie do roli Executor 2, zapisywane komunikatem `role.assign` (rozdz. 5.8 dokumentu Kontraktów komunikacji).

### B.2. Model pomocniczy na hoście własnym przez kanał SSH

Model pomocniczy uruchomiony na własnym serwerze zdalnym, wykorzystany w module Diagnostics przez kanał SSH, bez ponoszenia kosztów zewnętrznego dostawcy API.

1. Uruchomienie własnego modelu pomocniczego na serwerze zdalnym.
2. Konfiguracja kanału SSH w oknie konfiguracji, ze wskazaniem danych dostępowych do hosta; dane trafiają do przechowywania poza bazą (rozdz. 13), a w encji kanału modelu zapisywane jest wyłącznie odwołanie do nich.
3. Przypisanie kanału modelu bezpośrednio do sesji w module Diagnostics (rozdz. 11.13 koncepcji platformy).
4. Model pomocniczy wspiera analizę logów bez kosztów zewnętrznego dostawcy API.

### B.3. Ten sam model bazowy, dwie odrębne tożsamości

Jeden kanał modelu i jeden model bazowy, wykorzystane z dwiema odrębnymi tożsamościami w module Roundtable (rozdz. 11.8 koncepcji platformy).

1. Dwa okna Model Panels wykorzystują ten sam kanał API i ten sam model bazowy.
2. Pierwszemu modelowi nadawana jest nazwa „Zwolennik” oraz prompt systemowy ukierunkowujący go na argumentację za rozpatrywaną tezą.
3. Drugiemu modelowi nadawana jest nazwa „Oponent” oraz prompt systemowy ukierunkowujący go na krytykę tej tezy.
4. Debate Panel rejestruje wymianę argumentów między obiema tożsamościami, mimo że technicznie oba przypisania korzystają z tego samego kanału modelu.

### B.4. Zmiana poświadczenia konta bez wpływu na kanały powiązane

Zmiana klucza dostępu dostawcy przechowywanego w rejestrze kont, wykonana z sekcji Konta modeli Okna Ustawień, bez konieczności ponownej konfiguracji kanałów modelu odwołujących się do tego konta.

1. Operator otwiera sekcję Konta modeli i wybiera konto API, którego klucz dostępu utracił ważność.
2. `account.test` zwraca `ok: false` z polem `detail` wskazującym przyczynę — poświadczenie odrzucone przez dostawcę.
3. Operator zapisuje nowy klucz komendą `account.update`, podając wyłącznie pole `credential`; pola pominięte (`name`, `provider`, `defaultModel`) zostają bez zmiany.
4. Kanały modelu odwołujące się do tego konta polem `accountId` struktury `Channel` (rozdz. 15.3) nie wymagają odrębnej aktualizacji — kanał mówi CZYM się połączyć, konto CZYIM poświadczeniem, a te dwa fakty są rozdzielone.
5. Ponowne `account.test` zwraca `ok: true`; kolejne wywołanie modelu przez którykolwiek z powiązanych kanałów korzysta z nowego poświadczenia bez restartu sesji.

Poniższy diagram sekwencji przedstawia ten przebieg.

```
Operator              Powłoka                    Rdzeń
   │                     │                          │
   │  otwiera konto      │                          │
   │  z kluczem          │                          │
   │  wygasłym           │                          │
   ├────────────────────►│                          │
   │                     │  account.test            │
   │                     ├─────────────────────────►│
   │                     │  ok:false, detail        │
   │                     │◄─────────────────────────┤
   │  wpisuje nowy klucz │                          │
   ├────────────────────►│  account.update          │
   │                     │  (accountId, credential) │
   │                     ├─────────────────────────►│
   │                     │  account:Account         │
   │                     │◄─────────────────────────┤
   │                     │  account.test            │
   │                     ├─────────────────────────►│
   │                     │  ok:true                 │
   │                     │◄─────────────────────────┤
   │  potwierdzenie      │                          │
   │◄────────────────────┤                          │
```

Legenda: `account.test`, `account.update` — komendy obszaru `account`; żaden z kanałów odwołujących się do konta polem `accountId` nie jest wywoływany w tym przebiegu — zmiana poświadczenia jest lokalna dla encji `Account`.

---

## Załącznik C. Pełny wykaz komend kontraktu warstwy modeli

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `channel` — 6 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `channel.add` | Dodaje wiersz do rejestru kanałów modelu | `name:string` (wym)<br>`kind:string` (wym)<br>`model:string` (opc)<br>`enabled:bool` (opc)<br>`config:json` (opc) | `channel:Channel` (wym) |
| `channel.update` | Zmienia wiersz rejestru kanałów | `channelId:string` (wym)<br>`name:string` (opc)<br>`model:string` (opc)<br>`enabled:bool` (opc)<br>`config:json` (opc) | `channel:Channel` (wym) |
| `channel.remove` | Usuwa wiersz rejestru kanałów | `channelId:string` (wym) | `channelId:string` (wym) — usunięty kanał |
| `channel.list` | Zwraca rejestr kanałów modelu | `enabledOnly:bool` (opc) | `channels:Channel[]` (wym) |
| `channel.check` | Sprawdza, czy kanał modelu odpowiada, i oddaje wynik sprawdzenia. Narzędzie pomocnicze, nie bramka: wynik nie warunkuje zapisu eksperta ani wysłania tury | `channelId:string` (wym)<br>`transport:ProviderTransport` (opc) | `reachable:bool` (wym)<br>`checkedAt:int64` (wym)<br>`latencyMs:int` (opc)<br>`detail:string` (opc) |
| `channel.credential.status` | Oddaje stan poświadczenia kanału: czy jest ustawione, kiedy zostało zmienione i jakiego jest rodzaju. Nie oddaje treści poświadczenia i oddawać jej nie może | `channelId:string` (wym) | `status:ChannelCredentialStatus` (wym) |

### Obszar `identity` — 5 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `identity.category.list` | Zwraca katalog kategorii zasad i tożsamości modelu | `layer:IdentityLayer` (opc)<br>`includeDisabled:bool` (opc) | `categories:IdentityCategory[]` (wym) |
| `identity.document.get` | Odczytuje treść kategorii zasad zapisaną dla wskazanej osi | `categoryId:string` (opc)<br>`axis:ConfigAxis` (opc)<br>`axisId:string` (opc) | `documents:IdentityDocument[]` (wym) |
| `identity.document.set` | Zapisuje treść kategorii zasad dla wskazanej osi wraz z trybem podania | `categoryId:string` (wym)<br>`content:string` (wym)<br>`mode:IdentityMode` (opc)<br>`axis:ConfigAxis` (opc)<br>`axisId:string` (opc)<br>`enabled:bool` (opc) | `document:IdentityDocument` (wym) |
| `identity.document.remove` | Usuwa treść kategorii zasad; brak zapisu znaczy treść kategorii z osi szerszej | `documentId:string` (wym) | `removed:bool` (wym) |
| `identity.effective.get` | Zwraca nakładkę obowiązującą: tryb, warstwy w kolejności krytyczności i złożony prompt systemowy | `windowId:string` (opc)<br>`axis:ConfigAxis` (opc)<br>`axisId:string` (opc) | `mode:IdentityMode` (wym)<br>`layers:IdentityLayerContent[]` (wym)<br>`prompt:string` (wym)<br>`promptHash:string` (wym)<br>`missingRequiredCategoryIds:string[]` (opc) |

**Zdarzenia obszaru `identity` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `identity.changed` | Zmiana treści kategorii zasad i tożsamości modelu | — |

### Obszar `model` — 1 komenda

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `model.channel.set` | Ustala kanał modelu obsługujący okno albo kartę sesji. Kanały zakłada i zmienia rodzina `channel.*`; ta komenda wybiera jeden z założonych | `modelChannelId:string` (wym)<br>`windowId:string` (opc)<br>`sessionId:string` (opc)<br>`agentId:string` (opc) | `channel:Channel` (wym)<br>`windowId:string` (opc)<br>`sessionId:string` (opc) |

### Obszar `account` — 6 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `account.add` | Zakłada konto modelu albo konto programu code CLI | `name:string` (wym)<br>`kind:AccountKind` (wym)<br>`provider:string` (wym)<br>`externalId:string` (opc)<br>`defaultModel:string` (opc)<br>`baseUrl:string` (opc)<br>`configDir:string` (opc)<br>`credential:string` (opc)<br>`makeDefault:bool` (opc)<br>`enabled:bool` (opc) | `account:Account` (wym) |
| `account.list` | Zwraca wykaz kont modeli i kont code CLI; poświadczeń wykaz nie niesie | `kind:AccountKind` (opc)<br>`enabledOnly:bool` (opc) | `accounts:Account[]` (wym)<br>`total:int` (wym) |
| `account.update` | Zmienia konto; pominięte pole zostaje bez zmiany | `accountId:string` (wym)<br>`name:string` (opc)<br>`provider:string` (opc)<br>`externalId:string` (opc)<br>`defaultModel:string` (opc)<br>`baseUrl:string` (opc)<br>`configDir:string` (opc)<br>`credential:string` (opc)<br>`enabled:bool` (opc) | `account:Account` (wym) |
| `account.remove` | Usuwa konto; kanały powołujące się na nie tracą powiązanie, nie znikają | `accountId:string` (wym) | `removed:bool` (wym)<br>`detachedChannelIds:string[]` (opc) |
| `account.default.set` | Wskazuje konto domyślne swojego rodzaju; domyślne jest dokładnie jedno na rodzaj | `accountId:string` (wym)<br>`kind:AccountKind` (wym) | `account:Account` (wym)<br>`previousDefaultId:string` (opc) |
| `account.test` | Próbne sprawdzenie konta modelu bez zmiany jego zapisu — działanie „Testuj połączenie” sekcji Konta modeli Okna Ustawień. Sprawdza to, co da się sprawdzić z serwera platformy: obecność poświadczenia w sejfie, dla konta code CLI istnienie katalogu konfiguracji i programu wiersza poleceń w pakiecie serwera, dla konta API osiągalność adresu punktu końcowego. Wynik jest wynikiem sprawdzenia, nie wywołaniem modelu — nie zużywa limitu konta | `accountId:string` (wym) | `ok:bool` (wym)<br>`detail:string` (wym)<br>`checkedAt:int64` (wym) |

**Zdarzenia obszaru `account` — 1:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `account.changed` | Zmiana konta modelu albo konta code CLI; poświadczeń zdarzenie nie niesie | — |

Razem w wykazie: **18 komend** z czterech obszarów kontraktu.

---

*Koniec dokumentu. Integracja modeli — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
