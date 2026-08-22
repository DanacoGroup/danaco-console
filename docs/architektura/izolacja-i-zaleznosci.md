# Danaco Console — Izolacja i konfigurowalność zależności

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
| **Tytuł** | Izolacja i konfigurowalność zależności |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (główny) · projektant (okno konfiguracji punktów izolacji) |
| **Przeznaczenie** | Ustala pełną konfigurowalność izolacji technicznej i izolacji kontekstu platformy — cztery rodzaje izolacji, osiem zakresów technicznych, poziomy zasięgu i ich pierwszeństwo, profile izolacji oraz zasadę zera blokad twardych |
| **Zakres** | zakresy izolacji technicznej i izolacji kontekstu, łączenie wielu zakresów w politykę, poziomy zasięgu i pierwszeństwo, profile izolacji, przypisanie punktów izolacji do sesji/roli/projektu, współdzielenie historii i pamięci, okno konfiguracji punktów izolacji, punkty izolacji warstwy komunikacji operacyjnej, zależności izolacji i warstw widoczności, obszar `isolation` kontraktu |
| **Poza zakresem** | mechanizm warstw widoczności jako taki — [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) rozdz. 14; pełny model danych profilu izolacji — [Model danych](model-danych.md); układ okien interfejsu — [Przepływ okien](../interfejs-uzytkownika/przeplyw-okien.md) |
| **Dokument nadrzędny** | [Architektura techniczna](architektura.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](koncepcja-platformy.md) · [Model konfiguracji](model-konfiguracji.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Model danych](model-danych.md) · [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) |
| **Źródła normatywne** | `budowa/shared/contract.json` (obszar `isolation`; struktury `IsolationPolicy`, `IsolationProfile`, `IsolationScopeLevel`, `IsolationSwitch`, `IsolationTechnicalSwitch`; wyliczenia `IsolationTechnicalScope`, `IsolationLayer`, `ConfigScope`) |
| **Zasada nadrzędna** | Zero blokad twardych — każdy zakres izolacji jest ustawieniem konfigurowalnym na ośmiu poziomach zasięgu wymienionych w rozdz. 6, nigdy wartością wpisaną na stałe |

Dokument stanowi pełną specyfikację konfigurowalności zależności i punktów izolacji w platformie Danaco Console: zakresy izolacji technicznej procesu sesji, zakresy izolacji kontekstu, punkty izolacji warstwy komunikacji operacyjnej (Chat Window — Użytkownik ↔ Wykonawca; Execution Loop Window — Koordynator ↔ Wykonawca), zależności i punkty izolacji warstw widoczności, łączenie wielu zakresów w jedną politykę, profile izolacji, przypisanie punktów izolacji do sesji, roli i projektu, współdzielenie historii i pamięci między kartami sesji, modułami i środowiskami, a także szczegółowy interfejs okna konfiguracji punktów izolacji. Opisy mają postać tabel zestawczych, schematów tekstowych, diagramów decyzyjnych oraz szablonów konfiguracji.

---

## Spis treści

1. [Wprowadzenie i cel dokumentu](#wprowadzenie-i-cel-dokumentu)
2. [Zasada nadrzędna: konfigurowalność zamiast blokad](#1-zasada-nadrzędna-konfigurowalność-zamiast-blokad)
3. [Cztery rodzaje izolacji](#2-cztery-rodzaje-izolacji)
4. [Zakresy izolacji technicznej](#3-zakresy-izolacji-technicznej)
5. [Zakresy izolacji kontekstu](#4-zakresy-izolacji-kontekstu)
6. [Łączenie wielu zakresów izolacji](#5-łączenie-wielu-zakresów-izolacji)
7. [Poziomy zasięgu reguł izolacji i pierwszeństwo](#6-poziomy-zasięgu-reguł-izolacji-i-pierwszeństwo)
8. [Profile izolacji](#7-profile-izolacji)
   - [7.1 Operacje na profilach](#71-operacje-na-profilach)
   - [7.2 Szablon profilu izolacji](#72-szablon-profilu-izolacji)
   - [7.3 Przykładowe profile izolacji](#73-przykładowe-profile-izolacji)
9. [Przypisanie punktów izolacji do sesji, roli, okna i projektu](#8-przypisanie-punktów-izolacji-do-sesji-roli-okna-i-projektu)
10. [Współdzielenie historii i pamięci między kartami, modułami i środowiskami](#9-współdzielenie-historii-i-pamięci-między-kartami-modułami-i-środowiskami)
11. [Okno konfiguracji punktów izolacji](#10-okno-konfiguracji-punktów-izolacji)
   - [10.1 Układ okna](#101-układ-okna)
   - [10.2 Zawartość trzech paneli](#102-zawartość-trzech-paneli)
   - [10.3 Przebieg konfiguracji](#103-przebieg-konfiguracji)
12. [Warstwy konfiguracji: domyślna i sesji](#11-warstwy-konfiguracji-domyślna-i-sesji)
13. [Stan wyjściowy platformy](#12-stan-wyjściowy-platformy)
14. [Scenariusze konfiguracji punktów izolacji](#13-scenariusze-konfiguracji-punktów-izolacji)
   - [13.1 Rola wykonawcza na zdalnym hoście z osobnym kontem](#131-rola-wykonawcza-na-zdalnym-hoście-z-osobnym-kontem)
   - [13.2 Dwie karty sesji nad tym samym repozytorium](#132-dwie-karty-sesji-nad-tym-samym-repozytorium)
   - [13.3 Profil projektu z jednorazowym wyjątkiem](#133-profil-projektu-z-jednorazowym-wyjątkiem)
   - [13.4 Współdzielona pamięć projektu między Coordinatorem a Executorem](#134-współdzielona-pamięć-projektu-między-coordinatorem-a-executorem)
15. [Punkty izolacji warstwy komunikacji operacyjnej](#14-punkty-izolacji-warstwy-komunikacji-operacyjnej)
   - [14.1 Pozycje izolacji obu kanałów](#141-pozycje-izolacji-obu-kanałów)
   - [14.2 Zasięgi izolacji warstwy komunikacji operacyjnej](#142-zasięgi-izolacji-warstwy-komunikacji-operacyjnej)
   - [14.3 Współdzielenie kontekstu między zadaniami jednego zlecenia](#143-współdzielenie-kontekstu-między-zadaniami-jednego-zlecenia)
   - [14.4 Granice izolacji](#144-granice-izolacji)
16. [Zależności i punkty izolacji warstw widoczności](#15-zależności-i-punkty-izolacji-warstw-widoczności)
   - [15.1 Pozycje izolacji warstw widoczności](#151-pozycje-izolacji-warstw-widoczności)
   - [15.2 Zasięg przypisania elementu do warstwy](#152-zasięg-przypisania-elementu-do-warstwy)
   - [15.3 Zależności między warstwami widoczności a pozostałymi punktami izolacji](#153-zależności-między-warstwami-widoczności-a-pozostałymi-punktami-izolacji)
17. [Słowniczek pojęć](#16-słowniczek-pojęć)
18. [Kryteria odbioru](#17-kryteria-odbioru)
19. [Załącznik A. Macierz zakresów izolacji technicznej z objaśnieniami kontekstowymi](#załącznik-a-macierz-zakresów-izolacji-technicznej-z-objaśnieniami-kontekstowymi)
20. [Załącznik B. Szablon profilu izolacji](#załącznik-b-szablon-profilu-izolacji)
21. [Załącznik C. Tabela poziomów zasięgu i pierwszeństwa](#załącznik-c-tabela-poziomów-zasięgu-i-pierwszeństwa)
22. [Załącznik D. Zestawienie punktów izolacji warstwy komunikacji operacyjnej](#załącznik-d-zestawienie-punktów-izolacji-warstwy-komunikacji-operacyjnej)
23. [Załącznik E. Zestawienie punktów izolacji warstw widoczności](#załącznik-e-zestawienie-punktów-izolacji-warstw-widoczności)
24. [Załącznik F. Struktury danych pełne](#załącznik-f-struktury-danych-pełne)
   - [F.1. `IsolationPolicy` — polityka izolacji obowiązująca poziom zasięgu](#f1-isolationpolicy--polityka-izolacji-obowiązująca-poziom-zasięgu)
   - [F.2. `IsolationProfile` — profil zapisany i przypisywalny](#f2-isolationprofile--profil-zapisany-i-przypisywalny)
   - [F.3. `IsolationScopeLevel`, `IsolationSwitch`, `IsolationTechnicalSwitch`](#f3-isolationscopelevel-isolationswitch-isolationtechnicalswitch)
   - [F.4. Wyliczenia obszaru izolacji](#f4-wyliczenia-obszaru-izolacji)
25. [Załącznik G. Pełny wykaz komend kontraktu punktów izolacji](#załącznik-g-pełny-wykaz-komend-kontraktu-punktów-izolacji)
   - [Obszar `isolation` — 12 komend](#obszar-isolation--12-komend)

---

## Wprowadzenie i cel dokumentu

Platforma nie narzuca twardej izolacji historii sesyjnej ani pamięci, a wszystkie zależności platformy są konfigurowalne z poziomu okna konfiguracji. Techniczny wymiar izolacji obejmuje osiem zakresów dotyczących procesu sesji po stronie serwera oraz warstwową strukturę konfiguracji. Niniejsze opracowanie łączy oba wymiary z izolacją warstwy komunikacji operacyjnej i z izolacją warstw widoczności w jedną, spójną specyfikację funkcjonalną okna konfiguracji, modelu danych i warstwy sesji.

Dokument odpowiada na siedem pytań:

| Pytanie | Opisane w |
|---|---|
| Jakie punkty izolacji istnieją i co dokładnie obejmuje każdy z nich | rozdz. 3–4, 14–15 |
| W jaki sposób punkty izolacji łączy się w jedną politykę | rozdz. 5 |
| Czym jest profil izolacji i jak go zapisać oraz wykorzystać ponownie | rozdz. 7 |
| Jak przypisać punkty izolacji do konkretnej sesji, roli lub projektu | rozdz. 8 |
| Jak skonfigurować współdzielenie historii i pamięci między kartami sesji, modułami i środowiskami | rozdz. 9 |
| Jak izolowane są oba kanały komunikacji operacyjnej — Chat Window i Execution Loop Window | rozdz. 14 |
| Jakie zależności i punkty izolacji obowiązują warstwy widoczności interfejsu | rozdz. 15 |

Rozdział 10 opisuje szczegółowy interfejs okna konfiguracji punktów izolacji.

---

## 1. Zasada nadrzędna: konfigurowalność zamiast blokad

Cała treść niniejszego dokumentu podlega zasadzie centralnej Koncepcji platformy (rozdz. 14, zasada 2): każdą funkcję, każde zachowanie i każdą zależność między elementami platformy konfiguruje się z poziomu okna konfiguracji, a żadne zachowanie nie jest zaszyte na stałe w sposób niedostępny dla użytkownika. W odniesieniu do izolacji zasada ta przyjmuje trzy konkretne postacie:

| Postać zasady | Co oznacza dla izolacji |
|---|---|
| Izolacja jest możliwością, nie regułą wymuszoną | Żaden zakres izolacji opisany w rozdziałach 3, 4, 14 i 15 nie jest aktywny, dopóki użytkownik świadomie go nie włączy — ani globalnie, ani dla pojedynczej sesji. |
| Brak ustawienia oznacza wartość domyślną, a nie blokadę | Poziom zasięgu, na którym reguła izolacji nie została skonfigurowana, dziedziczy wartość z poziomu szerszego (rozdz. 6), aż do poziomu globalnego, który sam w sobie nie zawiera żadnej aktywnej izolacji technicznej (rozdz. 12). |
| Okno konfiguracji punktów izolacji nigdy nie wymusza wyboru | Każdy przełącznik w macierzy izolacji (rozdz. 10) pozostaje dopuszczalny w stanie nieustawionym; system przyjmuje wówczas regułę dziedziczoną, a nie regułę restrykcyjną. |

---

## 2. Cztery rodzaje izolacji

Punkty izolacji w platformie dzielą się na cztery rodzaje o odmiennej naturze, zebrane w jednym oknie konfiguracji (rozdz. 10).

| Rodzaj izolacji | Czego dotyczy | Liczba pozycji |
|---|---|---|
| Izolacja kontekstu | Historia rozmowy, pamięć długoterminowa i bieżący kontekst roboczy współdzielone lub odrębne między kartami sesji, modułami, środowiskami i projektami | 3 (rozdz. 4) |
| Izolacja techniczna | Parametry procesu sesji po stronie serwera: katalog, środowisko, sieć, pliki, konto, model procesu, serwer wykonania | 8 (rozdz. 3) |
| Izolacja warstwy komunikacji operacyjnej | Historia i pamięć głównego okna komunikacji Chat Window (Użytkownik ↔ Wykonawca) oraz zlecenia, zadania, kolejka i rejestr przebiegu pętli wykonawczej Execution Loop Window (Koordynator ↔ Wykonawca) | 5 (rozdz. 14) |
| Izolacja warstw widoczności | Konfiguracja warstw widoczności według roli użytkownika, rejestr funkcji wyszukiwarki funkcji oraz zasięg przypisania elementu interfejsu do warstwy | 3 (rozdz. 15) |

Cztery rodzaje izolacji są od siebie niezależne — świadome połączenie wszystkich, opisane w rozdziale 5, daje pełny obraz polityki obowiązującej dany zasięg.

```
CZTERY NIEZALEŻNE RODZAJE IZOLACJI

  Izolacja kontekstu  ───────────  historia · pamięć · kontekst
        │                          3 pozycje — stan: współdzielone / odrębne
        │
        ▼
  Izolacja techniczna ───────────  katalog · środowisko · dane modelu · sieć ·
        │                          pliki · konto · model procesu · serwer
        │                          8 zakresów — stan: włączony / wyłączony
        │
        ▼
  Izolacja komunikacji ──────────  historia Chat Window · pamięć Chat Window ·
  operacyjnej                      zlecenia i zadania · kolejka pętli ·
        │                          rejestr przebiegu pętli
        │                          5 pozycji — stan: współdzielone / odrębne
        │
        ▼
  Izolacja warstw     ───────────  konfiguracja warstw wg roli ·
  widoczności                      rejestr funkcji wyszukiwarki ·
                                   zasięg przypisania elementu do warstwy
                                   3 pozycje — stan: współdzielone / odrębne

  Włączenie izolacji technicznej dla sesji  ≠  rozdzielenie jej historii/pamięci
  Odrębna historia Chat Window              ≠  odrębny rejestr pętli wykonawczej
  Pełny obraz polityki daje połączenie wszystkich czterech (rozdz. 5)
```

---

## 3. Zakresy izolacji technicznej

Zakresy izolacji technicznej pochodzą wprost z rozdziału 12 dokumentu Architektury i dotyczą procesu sesji uruchamianego po stronie serwera (Architektura, rozdz. 7: „Proces sesji jest izolowany: własny katalog roboczy i środowisko”). Platforma udostępnia osiem zakresów; każdy stosuje się niezależnie i domyślnie żaden nie jest aktywny (rozdz. 12). Poniższa tabela podaje pełny opis funkcjonalny każdego zakresu — czego dotyczy oraz co oznacza jego stan włączony i wyłączony (domyślny). Treść objaśnień kontekstowych `[?]` prezentowanych przy każdym przełączniku zawiera Załącznik A.

| # | Zakres | Czego dotyczy | Stan włączony | Stan wyłączony (domyślny) | Odniesienie źródłowe |
|---|---|---|---|---|---|
| 1 | Katalog roboczy sesji | Fizyczny katalog plików na serwerze, w którym proces danej sesji zapisuje i odczytuje pliki robocze | Sesja lub rola pracuje we własnym katalogu, niewidocznym dla innych sesji tego samego zasięgu | Sesja korzysta ze wspólnego katalogu roboczego, dzielonego z innymi sesjami tego samego zasięgu — z innymi kartami otwartymi w tym samym projekcie modułu Workspace | Architektura, rozdz. 12; Koncepcja, rozdz. 11.2 |
| 2 | Środowisko procesu | Zmienne środowiskowe i kontekst uruchomieniowy procesu serwera przypisanego sesji | Sesja lub rola otrzymuje własny, odrębny zestaw zmiennych środowiskowych | Proces dziedziczy wspólne środowisko uruchomieniowe serwera, współdzielone z innymi procesami sesji | Architektura, rozdz. 12 |
| 3 | Katalog danych i konfiguracji modelu | Miejsce przechowywania danych pomocniczych wykorzystywanych przez kanał modelu przypisany sesji lub roli — konfiguracji, danych tymczasowych, ustawień dostawcy | Każda sesja lub rola ma własny katalog danych modelu | Katalog jest współdzielony między sesjami korzystającymi z tego samego kanału modelu | Architektura, rozdz. 12 i 9 |
| 4 | Dostęp sieciowy | Zakres, w jakim proces sesji nawiązuje połączenia sieciowe wychodzące — do zewnętrznych interfejsów API, do stron przeglądanych w module Browser, do serwerów MCP i rozszerzeń | Sesja otrzymuje odrębny, ograniczony dostęp sieciowy, oddzielony od innych procesów | Proces korzysta ze wspólnego dostępu sieciowego serwera, na równi z pozostałymi sesjami | Koncepcja, rozdz. 11.4; Architektura, rozdz. 14 |
| 5 | Zakres odczytu i zapisu plików | Uprawnienia procesu sesji do odczytu i zapisu plików poza własnym katalogiem roboczym — w module Library, w repozytorium kodu lub w plikach innego projektu | Dostęp ograniczony jest do jawnie dozwolonych ścieżek | Proces ma pełny dostęp do zasobów plikowych dostępnych na serwerze w ramach uprawnień Operatora | Koncepcja, rozdz. 11.6 |
| 6 | Konto i token per sesja | Dane uwierzytelniające — token dostępu, klucz interfejsu programistycznego, dane logowania kanału modelu — przypisane indywidualnie do sesji lub roli, zamiast do platformy jako całości | Sesja korzysta z własnego, dedykowanego konta lub tokenu | Sesje korzystają ze wspólnych, platformowych danych dostępowych | Architektura, rozdz. 12, 9 i 15 |
| 7 | Model procesu (odrębny lub współdzielony) | Czy instancja procesu wykonawczego modelu jest dedykowana danej sesji lub roli, czy współdzielona z innymi; łączy się z wyborem kanału modelu (API, CLI, SSH, HTTP) dokonywanym „per sesja lub per rola” | Odrębny — każda sesja lub rola uruchamia własną, niezależną instancję procesu modelu | Współdzielony (domyślny) — sesje korzystają ze wspólnej puli procesów lub instancji modelu | Architektura, rozdz. 12 i 9 |
| 8 | Serwer wykonania | Fizyczny lub logiczny serwer, na którym wykonywany jest proces sesji — istotny zwłaszcza przy modelach pomocniczych uruchamianych na hoście zdalnym przez kanał SSH | Sesja lub rola przypisana jest do dedykowanego serwera wykonania, odrębnego od pozostałych | Sesja wykonywana jest na współdzielonym serwerze wykonawczym platformy | Architektura, rozdz. 12 i 9 |

---

## 4. Zakresy izolacji kontekstu

Izolacja kontekstu obejmuje trzy pozycje, każdą przyjmującą stan „współdzielone” albo „odrębne” (rozdz. 6.4 Koncepcji). Dotyczą one warstwy komunikacji i pamięci opisanej w rozdziale 2.4 i w rozdziale 8 Koncepcji, a nie procesu technicznego (rozdz. 3 niniejszego dokumentu).

| Pozycja | Definicja | Stan „odrębne” (domyślny) | Stan „współdzielone” |
|---|---|---|---|
| Historia | Zapis wymiany wiadomości między użytkownikiem a AI w danej karcie sesji (Model danych, rozdz. 8.4) | Każda nowa karta sesji zaczyna z pustą historią | Ten sam zapis rozmowy widoczny jest w kilku kartach, modułach lub środowiskach jednocześnie |
| Pamięć | Zasoby pamięci długoterminowej przypisane do poziomu globalnego, projektu, sesji lub środowiska (Model danych, rozdz. 9) | Pamięć przypisana jednemu zasięgowi nie jest widoczna w innym | Ten sam zasób pamięci zasila kilka zasięgów jednocześnie |
| Kontekst | Bieżący stan roboczy sesji: aktywne pliki, wybrany projekt, załączniki i zmienne przekazywane modelowi w danej chwili (Koncepcja, rozdz. 2.4) | Kontekst rekonfigurowany jest od nowa przy każdej zmianie środowiska lub modułu | Kontekst bieżącej pracy przenoszony jest między kartami, modułami lub środowiskami bez przeładowania |

Rozdział 9 opisuje zastosowanie tych trzech pozycji w konkretnych scenariuszach współdzielenia między kartami sesji, modułami i środowiskami. Pozycje izolacji kontekstu dotyczą sesji jako całości; odrębne pozycje przypisane obu kanałom komunikacji operacyjnej opisuje rozdział 14.

---

## 5. Łączenie wielu zakresów izolacji

Zakresy izolacji technicznej stosuje się niezależnie i łącznie — polityka izolacji obowiązująca dany zasięg jest zbiorem aktywnych zakresów (Architektura, rozdz. 12), a nie pojedynczym, wykluczającym się ustawieniem. Ta sama zasada niezależności obowiązuje trzy pozycje izolacji kontekstu (rozdz. 4), pięć pozycji izolacji warstwy komunikacji operacyjnej (rozdz. 14) oraz trzy pozycje izolacji warstw widoczności (rozdz. 15): każdą ustawia się osobno, bez wpływu na pozostałe. W efekcie polityka izolacji dowolnego zasięgu jest zawsze kombinacją do dziewiętnastu niezależnych ustawień, z których każde przyjmuje jedną z dwóch wartości.

```
ANATOMIA POLITYKI IZOLACJI (jeden poziom zasięgu)

  8 zakresów technicznych                3 pozycje kontekstu
  ┌───────────────────────────┐          ┌───────────────────────────┐
  │ katalog roboczy   wł/wył  │          │ historia   współdz./odr.  │
  │ środowisko proc.  wł/wył  │          │ pamięć     współdz./odr.  │
  │ dane modelu       wł/wył  │          │ kontekst   współdz./odr.  │
  │ dostęp sieciowy   wł/wył  │          └───────────────────────────┘
  │ odczyt/zapis pl.  wł/wył  │
  │ konto i token     wł/wył  │          5 pozycji komunikacji operacyjnej
  │ model procesu     wł/wył  │          ┌───────────────────────────┐
  │ serwer wykonania  wł/wył  │          │ historia Chat Window      │
  └───────────────────────────┘          │ pamięć Chat Window        │
                                         │ zlecenia i zadania        │
  3 pozycje warstw widoczności           │ kolejka pętli             │
  ┌───────────────────────────┐          │ rejestr przebiegu pętli   │
  │ konfiguracja warstw / rola│          └───────────────────────────┘
  │ rejestr funkcji wyszukiw. │
  │ zasięg przypisania do w.  │          = do 19 niezależnych ustawień,
  └───────────────────────────┘            każde o dwóch wartościach

                                         Polityka = zbiór aktywnych zakresów,
                                         a nie pojedyncze ustawienie
```

Poniższa tabela zestawia kombinacje powstające przez łączenie kilku zakresów jednocześnie. Kombinacje te nie są zamkniętym katalogiem — każdy zestaw ustawień dopuszczalny przez macierz izolacji (rozdz. 10) jest prawidłową polityką.

| Kombinacja | Aktywne zakresy | Zastosowanie |
|---|---|---|
| Sesja w pełni odizolowana | Wszystkie osiem zakresów technicznych włączonych, historia i pamięć odrębne | Rola w środowisku MultitaskingAI wykonywana na zdalnym hoście przez kanał SSH, niezależna pod każdym względem od pozostałych ról zespołu |
| Współdzielony projekt, odrębne wykonanie | Historia i pamięć współdzielone na poziomie projektu; katalog roboczy sesji i model procesu włączone (odrębne) | Dwóch wykonawców (Executor 1, Executor 2) pracuje równolegle nad tym samym projektem Workspace, korzystając ze wspólnej pamięci projektu, lecz nie nadpisując sobie nawzajem plików roboczych |
| Izolacja bezpieczeństwa bez izolacji kontekstu | Konto i token per sesja oraz dostęp sieciowy włączone; historia i pamięć pozostają współdzielone | Sesja korzysta z osobnego klucza dostawcy modelu, aby nie wyczerpywać wspólnego limitu, ale nadal czerpie z tej samej pamięci projektu co pozostałe karty |
| Odrębna historia bez izolacji technicznej | Historia i kontekst odrębne; wszystkie osiem zakresów technicznych wyłączonych | Stan wyjściowy platformy (rozdz. 12) — porządek pracy bez żadnego technicznego ograniczenia procesu |

Reguła łączenia zakresów obowiązuje wewnątrz jednego poziomu zasięgu (rozdz. 6). Jeżeli różne zakresy tej samej polityki są ustawione na różnych poziomach zasięgu — historia na poziomie projektu, a dostęp sieciowy na poziomie roli — o wyniku decyduje reguła pierwszeństwa opisana w rozdziale 6, stosowana osobno dla każdego zakresu.

---

## 6. Poziomy zasięgu reguł izolacji i pierwszeństwo

Platforma udostępnia osiem poziomów zasięgu, na których ustawia się dowolny z dziewiętnastu zakresów izolacji opisanych w rozdziałach 3, 4, 14 i 15. Poziomy uporządkowane są od najogólniejszego do najbardziej szczegółowego; reguła z poziomu bardziej szczegółowego ma pierwszeństwo przed regułą z poziomu szerszego (Koncepcja, rozdz. 6.5).

| # | Poziom zasięgu | Pierwszeństwo | Dziedziczy z | Przykład zastosowania |
|---|---|---|---|---|
| 1 | Globalny (domyślny) | Najniższe — warstwa bazowa | — | Domyślna polityka izolacji obowiązująca w całej platformie |
| 2 | Środowisko | rośnie ↑ | Globalny | Odmienna polityka dla środowiska CodeStudio niż dla TalkIn |
| 3 | Moduł | rośnie ↑ | Środowisko | Odrębna pamięć przypisana modułowi Developer |
| 4 | Para modułów (relacja) | rośnie ↑ | Moduł | Współdzielenie operacji kontekstowych AI między Studio a Translate |
| 5 | Projekt | rośnie ↑ | Para modułów | Rozdzielenie kontekstu między dwoma projektami w module Workspace |
| 6 | Karta sesji | rośnie ↑ | Projekt | Jednorazowe współdzielenie historii między dwiema otwartymi kartami |
| 7 | Rola (MultitaskingAI) | rośnie ↑ | Karta sesji | Profil izolacji przypisany roli Executor 1 |
| 8 | Okno komunikacji | Najwyższe — poziom najwęższy | Rola | Odrębna pamięć i odrębne uprawnienia okna Execution Loop Window wobec okna Chat Window tej samej karty sesji |

Zasada rozstrzygania brzmi: **pierwszeństwo zasięgu najbardziej szczegółowego.** Brak ustawienia na danym poziomie oznacza dziedziczenie wartości z poziomu bezpośrednio szerszego, aż do poziomu globalnego, który stanowi warstwę bazową. Reguła ta jest bezpośrednim zastosowaniem zasady nadrzędnej „brak ustawienia = wartość domyślna” (rozdz. 1) — żaden poziom nie musi być skonfigurowany, a poziom nieskonfigurowany nie blokuje pracy, lecz przejmuje wartość odziedziczoną. Poniższy diagram decyzyjny przedstawia rozstrzyganie wartości jednego zakresu.

```
ROZSTRZYGANIE WARTOŚCI JEDNEGO ZAKRESU
(pierwszeństwo zasięgu najbardziej szczegółowego — rozdz. 6)

  Czy zakres ustawiony na poziomie OKNO KOMUNIKACJI?
        │ tak ─► użyj wartości z Okno komunikacji        ◄── koniec
        │ nie
        ▼
  Czy zakres ustawiony na poziomie ROLA?
        │ tak ─► użyj wartości z poziomu Rola            ◄── koniec
        │ nie
        ▼
  Czy ustawiony na poziomie KARTA SESJI?
        │ tak ─► użyj wartości z Karta sesji             ◄── koniec
        │ nie
        ▼
  Czy ustawiony na poziomie PROJEKT?
        │ tak ─► użyj wartości z Projekt                 ◄── koniec
        │ nie
        ▼
  Czy ustawiony na poziomie PARA MODUŁÓW?
        │ tak ─► użyj wartości z Para modułów            ◄── koniec
        │ nie
        ▼
  Czy ustawiony na poziomie MODUŁ?
        │ tak ─► użyj wartości z Moduł                   ◄── koniec
        │ nie
        ▼
  Czy ustawiony na poziomie ŚRODOWISKO?
        │ tak ─► użyj wartości ze Środowisko             ◄── koniec
        │ nie
        ▼
  Użyj wartości z poziomu GLOBALNY (warstwa bazowa)
  — domyślnie brak aktywnej izolacji technicznej; kontekst odrębny
```

Dziedziczenie stosuje się osobno dla każdego z dziewiętnastu zakresów (rozdz. 3, 4, 14, 15). Nie ma wymogu, aby cała polityka danego poziomu była w pełni określona — jeden zakres może być ustawiony na poziomie projektu, a pozostałe dziedziczone z poziomu globalnego. Pełne zestawienie z miejscem konfiguracji zawiera Załącznik C.

**Stosunek ośmiu poziomów do wyliczenia `ConfigScope`.** Opis komendy `isolation.scope.list` w `budowa/shared/contract.json` mówi o „ośmiu poziomach zasięgu izolacji”, a wyliczenie `ConfigScope`, którego wartości struktura `IsolationScopeLevel` przenosi w polu `scope`, liczy dziewięć pozycji: `application`, `global`, `environment`, `module`, `modulePair`, `project`, `session`, `role`, `window`. Osiem poziomów wymienionych w tabeli powyżej odpowiada ośmiu spośród tych dziewięciu wartości (`global` → Globalny, `environment` → Środowisko, `module` → Moduł, `modulePair` → Para modułów, `project` → Projekt, `session` → Karta sesji, `role` → Rola, `window` → Okno komunikacji). Poziom `application` nie uczestniczy w rozstrzyganiu polityki izolacji — obowiązuje wyłącznie w konfiguracji aplikacji (Architektura, rozdz. 13) — co daje osiem poziomów właściwych izolacji, zgodnie z liczbą w opisie źródłowym. Parametr opcjonalny `windowId` komendy `isolation.scope.list` wylicza byty poziomu okna komunikacji.

---

## 7. Profile izolacji

Zestaw ustawień izolacji zapisuje się jako nazwany **profil izolacji** i wykorzystuje wielokrotnie bez ponownej konfiguracji (Koncepcja, rozdz. 6.6; Architektura, rozdz. 12: „zbiór aktywnych zakresów tworzy profil izolacji przypisywalny do sesji, roli lub projektu”). Profil grupuje sześć elementów:

| Element profilu | Zawartość |
|---|---|
| Poziom zasięgu | Poziom, dla którego profil jest przeznaczony (rozdz. 6) |
| Izolacja kontekstu | Wartości trzech pozycji: historia, pamięć, kontekst (rozdz. 4) |
| Izolacja techniczna | Stan ośmiu zakresów izolacji technicznej (rozdz. 3) |
| Izolacja komunikacji operacyjnej | Wartości pięciu pozycji obu kanałów komunikacji (rozdz. 14) |
| Izolacja warstw widoczności | Wartości trzech pozycji warstw widoczności (rozdz. 15) |
| Warstwa konfiguracji | Warstwa, do której profil należy — domyślna albo sesji (rozdz. 11) |

Profil jest niezależny od pojedynczego przypisania — ten sam profil o nazwie „Sesja izolowana — host zdalny” przypisuje się jednocześnie do kilku sesji, ról lub projektów, a zmiana profilu odbija się na wszystkich jego przypisaniach. Pełny szablon pól profilu zawiera Załącznik B.

### 7.1. Operacje na profilach

W panelu profilu i podglądu okna konfiguracji punktów izolacji (rozdz. 10) dostępne są cztery operacje na profilach:

| Operacja | Działanie |
|---|---|
| Zapisz profil | Utrwala bieżący stan macierzy izolacji (rozdz. 10) pod nazwaną etykietą |
| Wczytaj profil | Ładuje zapisany wcześniej profil do macierzy izolacji, nadpisując bieżący podgląd — bez automatycznego przypisania |
| Przypisz do… | Wiąże wczytany lub właśnie zapisany profil z konkretną sesją, rolą lub projektem (rozdz. 8) |
| Usuń profil | Usuwa nazwany profil z pamięci aplikacji; nie cofa ustawień już przypisanych sesji, rolom lub projektom, które zachowują ostatnio odziedziczone wartości |

### 7.2. Szablon profilu izolacji

Poniższy szablon konfiguracji porządkuje pola profilu jako strukturę wypełnianą przez użytkownika. Każde pole podlega zasadzie „brak ustawienia = dziedziczenie z poziomu szerszego” (rozdz. 1, rozdz. 6).

```
PROFIL IZOLACJI — struktura

  Nazwa profilu:          < nazwa własna >
  Opis:                   < krótkie uzasadnienie przeznaczenia >
  Poziom zasięgu:         globalny | środowisko | moduł | para modułów |
                          projekt | karta sesji | rola | okno komunikacji

  Izolacja kontekstu:
    historia              współdzielona | odrębna
    pamięć                współdzielona | odrębna
    kontekst              współdzielony | odrębny

  Izolacja techniczna:
    katalog roboczy sesji                 włączony | wyłączony
    środowisko procesu                    włączony | wyłączony
    katalog danych i konfiguracji modelu  włączony | wyłączony
    dostęp sieciowy                       włączony | wyłączony
    zakres odczytu i zapisu plików        włączony | wyłączony
    konto i token per sesja               włączony | wyłączony
    model procesu                         włączony (odrębny) | wyłączony (współdzielony)
    serwer wykonania                      włączony | wyłączony

  Izolacja komunikacji operacyjnej:
    historia Chat Window                  współdzielona | odrębna
    pamięć Chat Window                    współdzielona | odrębna
    zlecenia i zadania pętli              współdzielone | odrębne
    kolejka pętli wykonawczej             współdzielona | odrębna
    rejestr przebiegu pętli               współdzielony | odrębny

  Izolacja warstw widoczności:
    konfiguracja warstw wg roli           współdzielona | odrębna
    rejestr funkcji wyszukiwarki funkcji  współdzielony | odrębny
    zasięg przypisania elementu do warstwy  globalny | środowisko | moduł |
                                            profil | sesja

  Warstwa:                domyślna | sesji
  Przypisanie:            sesja | rola | projekt
  Wartość domyślna:       brak ustawienia = dziedziczenie z poziomu szerszego
```

### 7.3. Przykładowe profile izolacji

Poniższa macierz porównuje cztery przykładowe profile zbudowane z ustaleń opisanych w rozdziałach 5, 12 i 13. Zapis „dziedziczony” oznacza pole pozostawione nieustawionym, przejmujące wartość z poziomu szerszego (w nawiasie podano wartość odziedziczoną z poziomu globalnego).

| Zakres / pozycja | Sesja w pełni odizolowana | Executor zdalny — SSH | Projekt zamknięty | Stan wyjściowy platformy |
|---|---|---|---|---|
| Poziom zasięgu | Rola | Rola (Executor 2) | Projekt | Globalny (domyślny) |
| Historia | odrębna | dziedziczona (odrębna) | odrębna | odrębna |
| Pamięć | odrębna | dziedziczona (odrębna) | odrębna | odrębna |
| Kontekst | odrębny | dziedziczony (odrębny) | dziedziczony (odrębny) | odrębny |
| Katalog roboczy sesji | włączony | dziedziczony (wyłączony) | dziedziczony (wyłączony) | wyłączony |
| Środowisko procesu | włączone | dziedziczone (wyłączone) | dziedziczone (wyłączone) | wyłączone |
| Katalog danych i konfiguracji modelu | włączony | dziedziczony (wyłączony) | dziedziczony (wyłączony) | wyłączony |
| Dostęp sieciowy | włączony | włączony | dziedziczony (wyłączony) | wyłączony |
| Zakres odczytu i zapisu plików | włączony | dziedziczony (wyłączony) | dziedziczony (wyłączony) | wyłączony |
| Konto i token per sesja | włączony | włączony | dziedziczony (wyłączony) | wyłączony |
| Model procesu | włączony (odrębny) | włączony (odrębny) | dziedziczony (współdzielony) | wyłączony (współdzielony) |
| Serwer wykonania | włączony | włączony | dziedziczony (wyłączony) | wyłączony |
| Typowe zastosowanie | Rola MultitaskingAI na zdalnym hoście przez SSH, niezależna pod każdym względem od pozostałych ról | Model pomocniczy Executor 2 na zdalnym hoście przez SSH, z osobnym kontem i tokenem (rozdz. 13.1) | Rozdzielenie kontekstu między projektami w module Workspace (rozdz. 13.3) | Domyślny, uporządkowany podział pracy bez izolacji technicznej (rozdz. 12) |

---

## 8. Przypisanie punktów izolacji do sesji, roli, okna i projektu

Konfigurację izolacji ustanawia się na dwa uzupełniające się sposoby: **bezpośrednio**, ustawiając przełączniki macierzy izolacji dla wybranego poziomu zasięgu w oknie konfiguracji punktów izolacji (rozdz. 10), albo **przez przypisanie** wcześniej zapisanego profilu izolacji (rozdz. 7) do konkretnego obiektu platformy. Cztery poziomy zasięgu — karta sesji, rola, okno komunikacji i projekt — mają, poza centralnym oknem konfiguracji, także własny, kontekstowy punkt przypisania osadzony w miejscu, gdzie dany obiekt jest zarządzany.

| Poziom przypisania | Miejsce dokonania przypisania | Skutek | Odwołanie |
|---|---|---|---|
| Karta sesji | Uproszczone menu kontekstowe karty sesji ([Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md), menu kontekstowe) albo panel profilu okna konfiguracji punktów izolacji, warstwa „sesji” | Profil obowiązuje wyłącznie bieżącą kartę; nakłada się na warstwę domyślną, nie modyfikując jej (rozdz. 11) | Koncepcja, rozdz. 8; niniejszy dokument, rozdz. 6, 11 |
| Rola (MultitaskingAI) | Okno roli — Executor 1, Executor 2, Coordinator albo Executor 3 / Validator — w panelu orkiestracji, pole „Profil izolacji” | Profil obowiązuje daną rolę we wszystkich zespołach, w których rola bierze udział, o ile nie nadpisany na poziomie zespołu lub karty sesji | Koncepcja, rozdz. 13.3, 13.8, Zał. B, Zał. D.2 |
| Okno komunikacji | Komponent konfiguracji okna — Chat Window albo Execution Loop Window (rozdz. 14) — otwierany z komponentu tego okna | Profil obowiązuje wyłącznie wskazane okno; dwa okna tej samej karty sesji mogą mieć różną pamięć i różne uprawnienia | Architektura, rozdz. 5, 12; niniejszy dokument, rozdz. 6, 14 |
| Projekt (Workspace) | Atrybut projektu „przypisany profil izolacji” w oknie modułu Workspace | Profil obowiązuje wszystkie karty sesji otwarte w kontekście danego projektu, z zachowaniem pierwszeństwa poziomu karty sesji, gdy ta ma własne ustawienie | Koncepcja, rozdz. 11.2; Model danych, rozdz. 7 |

Pozostałe cztery poziomy zasięgu — globalny, środowisko, moduł i para modułów (rozdz. 6) — konfiguruje się wyłącznie centralnie, z poziomu panelu selektora zasięgu w oknie konfiguracji punktów izolacji, ponieważ nie odpowiadają one pojedynczemu, zarządzalnemu obiektowi platformy, lecz całej klasie obiektów — wszystkim sesjom otwartym w środowisku CodeStudio.

```
JAK PRZYPISAĆ REGUŁĘ IZOLACJI

  Wybór sposobu przypisania:
    ├─ Bezpośrednio ─ ustaw przełączniki macierzy izolacji dla wybranego
    │                 poziomu w oknie konfiguracji punktów izolacji (rozdz. 10)
    └─ Przez profil ─ przypisz zapisany profil izolacji (rozdz. 7) do obiektu

  Gdzie dokonać przypisania — według poziomu zasięgu:

    Poziom               Punkt przypisania
    ─────────────────    ───────────────────────────────────────────────────
    Karta sesji     ─►   menu kontekstowe karty sesji  LUB  panel profilu
    Rola            ─►   okno roli w panelu orkiestracji, pole „Profil izolacji”
    Projekt         ─►   atrybut projektu w oknie modułu Workspace
    ─────────────────    ───────────────────────────────────────────────────
    Globalny        ┐
    Środowisko      │    wyłącznie centralnie — panel selektora zasięgu
    Moduł           │    w oknie konfiguracji punktów izolacji (klasa obiektów,
    Para modułów    ┘    nie pojedynczy obiekt)
```

Niezależnie od miejsca dokonania przypisania obowiązuje jedna reguła pierwszeństwa (rozdz. 6): ustawienie na poziomie karty sesji ma pierwszeństwo przed ustawieniem roli, ustawienie roli — przed ustawieniem projektu, a ustawienie projektu — przed ustawieniami modułu, pary modułów, środowiska i poziomu globalnego, w tej kolejności zgodnej z tabelą w rozdziale 6.

---

## 9. Współdzielenie historii i pamięci między kartami, modułami i środowiskami

Rozdział ten opisuje zastosowanie trzech pozycji izolacji kontekstu (rozdz. 4) w konkretnych relacjach między obiektami platformy, zgodnie z zasadą, że stan wyjściowy jest uporządkowanym podziałem pracy, a nie ograniczeniem (Koncepcja, rozdz. 6.1).

| Relacja współdzielenia | Poziom zasięgu | Stan domyślny | Przykład zastosowania |
|---|---|---|---|
| Między dwiema kartami sesji tego samego modułu | Karta sesji | Odrębne | Dwie równoległe rozmowy w module Studio, prowadzone nad różnymi dokumentami |
| Między dwoma modułami (para modułów) | Para modułów | Odrębne | Studio i Translate współdzielą mechanizm operacji kontekstowych AI przy pracy dwujęzycznej (Koncepcja, rozdz. 11.1, 11.7, Zał. E.4) |
| Między środowiskami | Środowisko | Odrębne | Historia rozmowy prowadzonej w TalkIn nie przenika domyślnie do CodeStudio |
| Wewnątrz jednego projektu (Workspace) | Projekt | Odrębne od innych projektów; współdzielone wewnątrz tego samego projektu, o ile projekt tak skonfigurowano | Wszystkie karty sesji otwarte w kontekście jednego projektu Workspace mają dostęp do tej samej pamięci projektu (Context Memory) |
| Między rolami środowiska MultitaskingAI | Rola | Zależne od skonfigurowanego profilu roli | Coordinator i Executor 1 współdzielą pamięć projektu, aby Coordinator uwzględniał dotychczasowe ustalenia przy budowie kolejnych poleceń |

Mechanizm techniczny współdzielenia wykorzystuje strukturę pamięci opisaną w Modelu danych (rozdz. 9): zasób pamięci ma przypisany poziom (globalna, projekt, sesja, środowisko); ustawienie danej pozycji kontekstu na „współdzielone” dla wybranej pary obiektów oznacza, że oba obiekty odczytują i zapisują ten sam zasób pamięci lub tę samą historię wiadomości, zamiast każdy prowadzić własną.

```
MECHANIZM WSPÓŁDZIELENIA (Model danych, rozdz. 9)

  „odrębne” (domyślnie):              „współdzielone”:
    Obiekt A ─► własny zasób            Obiekt A ┐
    Obiekt B ─► własny zasób            Obiekt B ┴─► ten sam zasób pamięci
                                                    lub ta sama historia wiadomości

  • Zmiana obowiązuje od chwili aktywacji — nie przenosi istniejącej historii
    retroaktywnie (zgodnie z warstwową strukturą konfiguracji, rozdz. 11).
  • Przełączenie z powrotem na „odrębne” rozdziela dalszy przebieg pracy,
    pozostawiając dotychczas współdzieloną treść w obu miejscach w stanie
    z chwili rozdzielenia — współdzielenie jest zawsze odwracalne
    (Koncepcja, rozdz. 14, zasada 3).
```

---

## 10. Okno konfiguracji punktów izolacji

Niniejszy rozdział opisuje interfejs, w którym użytkownik konfiguruje punkty izolacji opisane w rozdziałach 3–9, 14 i 15.

### 10.1. Układ okna

Okno konfiguracji punktów izolacji jest otwierane z okna konfiguracji (Architektura, rozdz. 13) i stanowi jedno, spójne miejsce, w którym schodzą się izolacja kontekstu i izolacja techniczna. Dzieli się na trzy panele.

```
┌────────────────┬──────────────────────────────┬─────────────────────┐
│ SELEKTOR       │ MACIERZ IZOLACJI             │ PROFIL I PODGLĄD    │
│ ZASIĘGU        │                              │                     │
│                │ Izolacja kontekstu:          │ [ Zapisz profil ]   │
│ • Globalny     │   historia   [współdz.|odr.] │ [ Wczytaj profil ]  │
│ • Środowisko   │   pamięć     [współdz.|odr.] │ [ Przypisz do… ]    │
│ • Moduł        │   kontekst   [współdz.|odr.] │ [ Usuń profil ]     │
│ • Para modułów │                              │                     │
│ • Projekt      │ Izolacja techniczna:         │ Warstwa:            │
│ • Karta sesji  │   katalog roboczy   [wł|wył] │  ( ) domyślna       │
│ • Rola         │   środowisko proc.  [wł|wył] │  ( ) sesji          │
│                │   dane/konfig. model[wł|wył] │                     │
│                │   dostęp sieciowy   [wł|wył] │ Polityka efektywna: │
│                │   odczyt/zapis plik.[wł|wył] │  (podgląd reguł     │
│                │   konto i token     [wł|wył] │   dziedziczonych    │
│                │   model procesu     [wł|wył] │   z poziomów        │
│                │   serwer wykonania  [wł|wył] │   szerszych)        │
│                │                              │                     │
│                │ Komunikacja operacyjna:      │                     │
│                │   hist. Chat Window [w|o]    │                     │
│                │   pam. Chat Window  [w|o]    │                     │
│                │   zlecenia i zadania[w|o]    │                     │
│                │   kolejka pętli     [w|o]    │                     │
│                │   rejestr przebiegu [w|o]    │                     │
│                │                              │                     │
│                │ Warstwy widoczności:         │                     │
│                │   warstwy wg roli   [w|o]    │                     │
│                │   rejestr wyszukiw. [w|o]    │                     │
│                │   zasięg przypisania[▼]      │                     │
│                │                              │                     │
│                │                              │ [?] objaśnienia     │
│                │                              │     kontekstowe     │
└────────────────┴──────────────────────────────┴─────────────────────┘
```

### 10.2. Zawartość trzech paneli

| Panel | Położenie | Zawartość i funkcja |
|---|---|---|
| Selektor zasięgu | lewy | Lista ośmiu poziomów zasięgu opisanych w rozdziale 6, od najogólniejszego do najbardziej szczegółowego. Wybór poziomu ustala, dla jakiego obiektu (całej platformy, danego środowiska, danego modułu, wskazanej pary modułów, wskazanego projektu, wskazanej karty sesji, wskazanej roli lub wskazanego okna komunikacji) obowiązują ustawienia panelu środkowego. Wybór poziomu „para modułów”, „projekt”, „karta sesji”, „rola” lub „okno komunikacji” otwiera dodatkowy selektor wskazujący konkretny obiekt danego poziomu — parę Studio–Translate albo projekt z modułu Workspace. |
| Macierz izolacji | środkowy | Dwie grupy przełączników odpowiadające dwóm rodzajom izolacji (rozdz. 2): izolacja kontekstu — trzy przełączniki (historia, pamięć, kontekst), każdy w stanie „współdzielone” albo „odrębne” (rozdz. 4); izolacja techniczna — osiem przełączników odpowiadających zakresom z rozdziału 3, każdy w stanie „włączony” albo „wyłączony”; izolacja warstwy komunikacji operacyjnej — pięć przełączników (rozdz. 14); izolacja warstw widoczności — trzy pozycje (rozdz. 15). Każdy przełącznik opatrzony jest objaśnieniem kontekstowym `[?]` (Załącznik A). Przełącznik nieustawiony nie blokuje niczego — oznacza dziedziczenie z poziomu szerszego (rozdz. 6) i jest wizualnie odróżniony od przełącznika ustawionego jawnie. |
| Profil i podgląd | prawy | Zarządzanie profilem (zapis, wczytanie, przypisanie i usunięcie — rozdz. 7, rozdz. 8), wybór warstwy konfiguracji („domyślna” albo „sesji” — rozdz. 11) oraz podgląd polityki efektywnej: zestawienie wynikowych wartości wszystkich dziewiętnastu zakresów po uwzględnieniu dziedziczenia (rozdz. 6), z jawnym oznaczeniem, na którym poziomie każda wartość została faktycznie ustalona. |

Objaśnienie kontekstowe `[?]` towarzyszące każdemu przełącznikowi opisuje działanie danego ustawienia i jego wpływ na aplikację — zgodnie z zasadą przyjętą dla okna konfiguracji w rozdziale 13 Architektury:

> „Każde ustawienie posiada objaśnienie opisujące jego działanie i wpływ na aplikację; jest częścią definicji ustawienia i prezentowane przy elemencie w oknie konfiguracji.”

Podgląd polityki efektywnej pozwala użytkownikowi zweryfikować efekt kombinacji ustawień, zanim zmiana zostanie zapisana, zgodnie z zasadą jawności zależności (Koncepcja, rozdz. 14, zasada 3).

### 10.3. Przebieg konfiguracji

```
TYPOWY PRZEBIEG KONFIGURACJI

  [1] Selektor zasięgu ─────► wybór poziomu (i konkretnego obiektu)
            │
  [2] Macierz izolacji ─────► ustawienie przełączników izolacji kontekstu
            │                  lub technicznej (wsparcie: objaśnienia [?])
            │
  [3] Profil i podgląd ─────► zapis jako nazwany profil,
            │                  przypisanie do poziomu, wybór warstwy
            │                  (domyślna / sesji)
            │
  [4] Podgląd polityki ─────► weryfikacja wyniku po uwzględnieniu
      efektywnej              dziedziczenia z poziomów szerszych

  Okno nigdy nie wymusza żadnego kroku — użytkownik może je zamknąć
  na dowolnym etapie, pozostawiając pozycje nieustawione w stanie
  dziedziczonym (rozdz. 1, rozdz. 12).
```

1. Użytkownik wybiera w panelu lewym poziom zasięgu, na którym chce ustalić regułę — „Projekt”, a następnie konkretny projekt z listy.
2. W panelu środkowym ustawia wybrane przełączniki izolacji kontekstu lub izolacji technicznej, wspierając się objaśnieniami `[?]` przy każdej pozycji.
3. W panelu prawym zapisuje ustawienia jako nazwany profil i przypisuje je do wybranego poziomu; wybiera warstwę „domyślna” lub „sesji”, stosownie do trwałości, jaką chce nadać zmianie.
4. Weryfikuje w podglądzie polityki efektywnej, czy wynikowy zestaw reguł — po uwzględnieniu dziedziczenia z poziomów szerszych — odpowiada zamierzonemu efektowi.

Zgodnie z zasadą pełnej konfigurowalności (Koncepcja, rozdz. 14, zasada 2) oraz z nadrzędną zasadą braku twardych blokad okno nigdy nie wymusza żadnego z powyższych kroków — użytkownik może zamknąć okno na dowolnym etapie, pozostawiając nieustawione pozycje w stanie dziedziczonym, opisanym w rozdziale 12.

---

## 11. Warstwy konfiguracji: domyślna i sesji

Konfiguracja izolacji ma strukturę warstwową, zgodną z rozdziałem 13 Architektury.

| Warstwa | Gdzie ustalana | Kiedy obowiązuje | Wpływ na warstwę domyślną |
|---|---|---|---|
| Domyślna | Okno konfiguracji punktów izolacji (rozdz. 10) | Przy każdej nowej sesji, karcie, roli lub projekcie, w zależności od wybranego poziomu zasięgu | — (warstwa bazowa) |
| Sesji | Uproszczone menu kontekstowe okna operacyjnego ([Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md), menu kontekstowe) albo panel profilu okna konfiguracji punktów izolacji przy wybranej warstwie „sesji” (rozdz. 10.2) | Dla sesji bieżącej; kończy się z zamknięciem karty sesji | Nakłada się na warstwę domyślną bez modyfikacji jej wartości |

```
NAKŁADANIE WARSTW KONFIGURACJI

  Warstwa sesji     ▒▒▒▒▒   zmiany dla bieżącej karty (nietrwałe)
                    ─────   nakłada się na ↓ bez modyfikacji wartości
  Warstwa domyślna  █████   obowiązuje każdą nową sesję, kartę, rolę, projekt

  Zamknięcie karty sesji ─► warstwa sesji wygasa ─► kolejna karta w tym samym
  zasięgu dziedziczy wyłącznie warstwę domyślną (chyba że przypisano jej ten
  sam profil izolacji — rozdz. 7).
```

Nowa sesja dziedziczy warstwę domyślną; zmiany warstwy sesji nakładają się na warstwę domyślną, nie modyfikując jej wartości. Dzięki temu jednorazowe poluzowanie lub zaostrzenie izolacji w konkretnej sesji — tymczasowe współdzielenie historii między dwiema kartami na czas jednego zadania — nie zmienia domyślnego zachowania platformy dla kolejnych, nowo otwieranych sesji.

---

## 12. Stan wyjściowy platformy

Stanem wyjściowym — obowiązującym, dopóki użytkownik nie utworzy żadnego profilu ani nie zmieni żadnego ustawienia w oknie konfiguracji punktów izolacji — jest stan łączący dwa niezależne fakty:

| Wymiar stanu wyjściowego | Ustawienie domyślne | Charakter |
|---|---|---|
| Izolacja kontekstu | Każda nowa karta sesji otrzymuje odrębną historię, pamięć i kontekst (rozdz. 4) | Domyślny, uporządkowany podział pracy właściwy hierarchii środowisko → moduł → okno operacyjne (Koncepcja, rozdz. wprowadzenie, rozdz. 6), a nie ograniczenie możliwości użytkownika |
| Izolacja techniczna | Żaden z ośmiu zakresów izolacji technicznej (rozdz. 3) nie jest włączony — proces sesji korzysta ze wspólnego katalogu, środowiska, dostępu sieciowego, uprawnień plikowych, konta, puli procesów modelu i serwera wykonania współdzielonych z pozostałymi sesjami | Pełna operacyjna swoboda procesu, bez technicznego ograniczenia |
| Izolacja warstwy komunikacji operacyjnej | Historia i pamięć Chat Window odrębne dla każdej karty sesji; zlecenia, zadania, kolejka i rejestr przebiegu pętli odrębne dla każdego zlecenia, ze współdzieleniem kontekstu między zadaniami tego samego zlecenia (rozdz. 14) | Uporządkowany podział przebiegu pracy obu kanałów |
| Izolacja warstw widoczności | Konfiguracja warstw widoczności odrębna dla każdej roli użytkownika; rejestr funkcji wyszukiwarki funkcji współdzielony w obrębie środowiska; przypisanie elementu do warstwy o zasięgu globalnym (rozdz. 15) | Stopniowe ujawnianie funkcjonalności właściwe roli |

Ten stan wyjściowy łączy przejrzysty porządek kontekstu pracy z pełną operacyjną swobodą procesu. Wszelka dalsza izolacja — zarówno poluzowanie podziału kontekstu przez świadome włączenie współdzielenia (rozdz. 9), jak i włączenie dowolnego zakresu izolacji technicznej (rozdz. 3) — jest świadomą decyzją użytkownika podjętą w oknie konfiguracji punktów izolacji (rozdz. 10), a nie ustawieniem narzuconym przez platformę.

---

## 13. Scenariusze konfiguracji punktów izolacji

Poniższe scenariusze przedstawiono jako ponumerowane mini-przepływy wraz z tabelą konfiguracji wynikowej. Nie wprowadzają nowych ustaleń — ilustrują zastosowanie mechanizmów opisanych w rozdziałach 3–12.

### 13.1. Rola wykonawcza na zdalnym hoście z osobnym kontem

**Kontekst.** Użytkownik konfiguruje w środowisku MultitaskingAI rolę Executor 2, która ma korzystać z modelu pomocniczego uruchomionego na zdalnym hoście przez kanał SSH (Architektura, rozdz. 9), niezależnie od pozostałych ról zespołu.

1. W oknie konfiguracji punktów izolacji wybiera poziom zasięgu „Rola” i wskazuje Executor 2.
2. W macierzy izolacji technicznej włącza zakresy: serwer wykonania, model procesu, konto i token per sesja oraz dostęp sieciowy.
3. Pozostałe zakresy pozostawia dziedziczone z poziomu globalnego.
4. Zapisuje ustawienia jako profil „Executor zdalny — SSH” i przypisuje go do roli w oknie roli Executor 2 (rozdz. 8), zgodnie z Załącznikiem D.2 Koncepcji.

| Element konfiguracji | Wartość |
|---|---|
| Poziom zasięgu | Rola (Executor 2) |
| Zakresy techniczne włączone | Serwer wykonania, model procesu, konto i token per sesja, dostęp sieciowy |
| Pozostałe zakresy | Dziedziczone z poziomu globalnego |
| Zapisany profil | „Executor zdalny — SSH” |

### 13.2. Dwie karty sesji nad tym samym repozytorium

**Kontekst.** Dwaj współpracownicy pracujący na jednym koncie Operatora otwierają w module Developer dwie karty sesji nad tym samym repozytorium. Celem jest uniknięcie wzajemnego nadpisywania plików.

1. Na poziomie zasięgu „Karta sesji” użytkownik włącza dla drugiej karty zakres „katalog roboczy sesji”.
2. Historię czatu na poziomie „Para modułów” pozostawia nieustawioną — a więc odrębną, zgodnie ze stanem wyjściowym (rozdz. 12).
3. W efekcie każda karta prowadzi własną rozmowę, lecz obie operują na osobnych kopiach katalogu roboczego.

| Element konfiguracji | Wartość |
|---|---|
| Poziom zasięgu | Karta sesji (druga karta) |
| Zakres techniczny włączony | Katalog roboczy sesji |
| Historia (para modułów) | Nieustawiona = odrębna (stan wyjściowy) |
| Wynik | Osobne katalogi robocze, niezależne rozmowy w każdej karcie |

### 13.3. Profil projektu z jednorazowym wyjątkiem

**Kontekst.** Użytkownik przypisuje projektowi w module Workspace profil izolacji „Projekt zamknięty”, w którym pamięć i historia są odrębne od innych projektów (rozdz. 8). Dla jednej, tymczasowej karty chce skonsultować fragment tego projektu z ustaleniami z innego, równolegle prowadzonego projektu.

1. Przypisuje projektowi profil „Projekt zamknięty” (pamięć i historia odrębne od innych projektów).
2. Dla jednej, tymczasowej karty otwiera okno konfiguracji punktów izolacji i wybiera warstwę „sesji” (rozdz. 11).
3. Włącza współdzielenie pamięci wyłącznie dla tej karty.
4. Po zamknięciu karty warstwa sesyjna przestaje obowiązywać, a profil „Projekt zamknięty” pozostaje niezmieniony dla wszystkich pozostałych kart.

| Element konfiguracji | Wartość |
|---|---|
| Profil projektu | „Projekt zamknięty” (pamięć i historia odrębne) |
| Wyjątek jednorazowy | Warstwa „sesji” — współdzielenie pamięci dla jednej karty |
| Trwałość wyjątku | Wygasa z zamknięciem karty; profil projektu bez zmian |

### 13.4. Współdzielona pamięć projektu między Coordinatorem a Executorem

**Kontekst.** W zespole zbudowanym w środowisku MultitaskingAI użytkownik chce, aby Coordinator uwzględniał w budowanych promptach dotychczasowe ustalenia zapisane w pamięci projektu, z którego korzysta Executor 1.

1. W oknie konfiguracji punktów izolacji wybiera poziom zasięgu „Rola”.
2. Dla obu ról ustawia pozycję „pamięć” na „współdzielone”, wskazując jako wspólny zasięg pamięć danego projektu Workspace.
3. Pozostałe zakresy — w tym katalog roboczy i model procesu — pozostawia odrębne dla każdej roli, zgodnie z zasadą niezależnego łączenia zakresów (rozdz. 5).

| Element konfiguracji | Wartość |
|---|---|
| Poziom zasięgu | Rola (Coordinator, Executor 1) |
| Pozycja „pamięć” | Współdzielona — pamięć wskazanego projektu Workspace |
| Pozostałe zakresy | Odrębne dla każdej roli — katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania |

---

## 14. Punkty izolacji warstwy komunikacji operacyjnej

Warstwa komunikacji operacyjnej obejmuje dwa kanały pierwszoplanowe architektury platformy: **Chat Window** — główne okno komunikacji Użytkownik ↔ Wykonawca, umieszczone w lewej kolumnie obszaru roboczego — oraz **Execution Loop Window** — okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane jako kolumna sąsiadująca. Każdy kanał prowadzi własny zapis przebiegu pracy, a każdy z tych zapisów stanowi odrębny punkt izolacji konfigurowany w macierzy izolacji (rozdz. 10).

### 14.1. Pozycje izolacji obu kanałów

| # | Pozycja | Kanał | Czego dotyczy | Stan „odrębne” (domyślny) | Stan „współdzielone” |
|---|---|---|---|---|---|
| 1 | Historia Chat Window | Użytkownik ↔ Wykonawca | Zapis wymiany poleceń w języku naturalnym, strumienia odpowiedzi, zatwierdzeń i przerwań w głównym oknie komunikacji | Każda karta sesji prowadzi własny zapis rozmowy z Wykonawcą | Ten sam zapis rozmowy widoczny jest w kilku kartach, modułach lub środowiskach |
| 2 | Pamięć Chat Window | Użytkownik ↔ Wykonawca | Pamięć długoterminowa zasilająca główne okno komunikacji: ustalenia, preferencje pracy, wnioski z poprzednich rozmów | Pamięć jednej karty sesji nie jest widoczna w innej | Ten sam zasób pamięci zasila główne okno komunikacji w kilku zasięgach |
| 3 | Zlecenia i zadania | Koordynator ↔ Wykonawca | Bieżące zlecenie i jego dekompozycja na zadania, wraz z kontekstem przekazywanym Wykonawcy przy każdym zadaniu | Każde zlecenie prowadzi własny zestaw zadań i własny kontekst | Zadania kilku zleceń korzystają ze wspólnego zestawu ustaleń i wspólnego kontekstu |
| 4 | Kolejka pętli wykonawczej | Koordynator ↔ Wykonawca | Kolejka i stan zadań oczekujących, w realizacji, zakończonych i skierowanych do ponowienia | Każde zlecenie ma własną kolejkę, obsługiwaną niezależnie | Kilka zleceń obsługiwanych jest z jednej kolejki o wspólnym porządku wykonania |
| 5 | Rejestr przebiegu pętli | Koordynator ↔ Wykonawca | Zapis komunikatów sterujących, wyników kontroli jakości, decyzji o ponowieniu, wskaźników przebiegu oraz operacji sterowania (wstrzymanie, wznowienie, przerwanie, korekta zlecenia) | Rejestr prowadzony jest odrębnie dla każdego zlecenia | Jeden rejestr obejmuje przebieg kilku zleceń, modułów lub środowisk |

Pozycje 1–2 są niezależne od pozycji „Historia” i „Pamięć” izolacji kontekstu (rozdz. 4): tamte obejmują sesję jako całość, te — wyłącznie zapis głównego okna komunikacji. Pozycje 3–5 są niezależne wzajemnie: odrębna kolejka nie wymusza odrębnego rejestru przebiegu, a wspólny rejestr nie łączy kontekstu zadań.

### 14.2. Zasięgi izolacji warstwy komunikacji operacyjnej

Pięć pozycji ustawia się na tych samych ośmiu poziomach zasięgu, co pozostałe zakresy izolacji (rozdz. 6), z zachowaniem pierwszeństwa zasięgu najbardziej szczegółowego.

| Zasięg | Chat Window (pozycje 1–2) | Execution Loop Window (pozycje 3–5) |
|---|---|---|
| Karta sesji | Zapis rozmowy i pamięć główne okna komunikacji odrębne dla każdej karty | Zlecenie uruchomione z danej karty prowadzi własną kolejkę i własny rejestr |
| Projekt | Karty sesji jednego projektu współdzielą pamięć głównego okna komunikacji, o ile projekt tak skonfigurowano | Zlecenia jednego projektu współdzielą rejestr przebiegu, zachowując odrębne kolejki |
| Moduł | Historia głównego okna komunikacji nie przenika między modułami | Zlecenia jednego modułu nie widzą zadań ani kolejki zleceń innego modułu |
| Para modułów | Wskazana para modułów współdzieli historię i pamięć głównego okna komunikacji | Wskazana para modułów współdzieli rejestr przebiegu pętli |
| Środowisko | Historia rozmowy prowadzonej w jednym środowisku nie przenika do innego | Kolejka i rejestr pętli wykonawczej jednego środowiska są odrębne od kolejki i rejestru innego |
| Rola (MultitaskingAI) | Pamięć głównego okna komunikacji przypisana roli, niezależnie od zespołu | Koordynator prowadzi rejestr przebiegu wspólny dla nadzorowanych Wykonawców |
| Globalny | Warstwa bazowa: historia i pamięć głównego okna komunikacji odrębne | Warstwa bazowa: zlecenia, kolejka i rejestr odrębne dla każdego zlecenia |

### 14.3. Współdzielenie kontekstu między zadaniami jednego zlecenia

Wewnątrz jednego zlecenia kontekst jest współdzielony: wszystkie zadania powstałe z dekompozycji tego samego zlecenia odczytują ten sam zestaw ustaleń, wynik zadania poprzedzającego i tę samą pamięć zlecenia. Koordynator przekazuje Wykonawcy przy każdym zadaniu kontekst zlecenia w stanie aktualnym na chwilę przydziału.

```
ZASIĘG KONTEKSTU W PĘTLI WYKONAWCZEJ

  ZLECENIE A                              ZLECENIE B
  ┌──────────────────────────────┐        ┌────────────────────────────┐
  │ kontekst zlecenia A          │        │ kontekst zlecenia B        │
  │   ├─ zadanie A1 ─┐           │        │   ├─ zadanie B1            │
  │   ├─ zadanie A2 ─┼─ wspólny  │        │   └─ zadanie B2            │
  │   └─ zadanie A3 ─┘  odczyt   │        │                            │
  │ kolejka A · rejestr A        │        │ kolejka B · rejestr B      │
  └──────────────────────────────┘        └────────────────────────────┘
             ▲                                         ▲
             └──────── granica między zleceniami ──────┘
      przekroczenie granicy wyłącznie przez ustawienie pozycji
      „zlecenia i zadania” lub „rejestr przebiegu pętli” na „współdzielone”
```

### 14.4. Granice izolacji

| Granica | Co obowiązuje | Sposób przekroczenia |
|---|---|---|
| Między zadaniami jednego zlecenia | Kontekst, ustalenia i pamięć zlecenia są wspólne; kolejka porządkuje zadania jednego przebiegu | Granica nie występuje — współdzielenie jest stanem wyjściowym wewnątrz zlecenia |
| Między zleceniami | Zlecenie nie widzi zadań, kolejki ani rejestru innego zlecenia | Ustawienie pozycji „zlecenia i zadania”, „kolejka pętli wykonawczej” lub „rejestr przebiegu pętli” na „współdzielone” dla wskazanego zasięgu |
| Między modułami | Historia głównego okna komunikacji, kolejka i rejestr pętli nie przenikają między modułami | Ustawienie odpowiedniej pozycji na poziomie „moduł” albo „para modułów” (rozdz. 6) |
| Między środowiskami | Oba kanały prowadzą odrębny zapis w każdym środowisku | Ustawienie odpowiedniej pozycji na poziomie „środowisko”; ustawienie z poziomu bardziej szczegółowego ma pierwszeństwo |

Rejestr przebiegu pętli podlega tej samej regule odwracalności, co pozostałe pozycje kontekstu (rozdz. 9): przełączenie na „odrębne” rozdziela dalszy przebieg, pozostawiając treść dotychczas współdzieloną w obu miejscach w stanie z chwili rozdzielenia.

---

## 15. Zależności i punkty izolacji warstw widoczności

Interfejs platformy ujawnia funkcje stopniowo: element, który nie jest potrzebny do realizacji aktualnego zadania, nie jest widoczny. Każdy element interfejsu należy do dokładnie jednej z czterech warstw widoczności — warstwy zawsze widocznej, warstwy widocznej na żądanie, rozwinięć kontekstowych oraz funkcji eksperckich. Przypisanie elementu do warstwy, konfiguracja warstw dla roli użytkownika oraz rejestr funkcji wyszukiwarki funkcji stanowią trzy niezależne punkty izolacji konfigurowane w macierzy izolacji (rozdz. 10).

### 15.1. Pozycje izolacji warstw widoczności

| # | Pozycja | Czego dotyczy | Stan „odrębne” | Stan „współdzielone” |
|---|---|---|---|---|
| 1 | Konfiguracja warstw według roli użytkownika | Zestaw reguł określających, które elementy interfejsu należą do warstwy zawsze widocznej, a które do warstw wywoływanych, dla danej roli użytkownika | Każda rola użytkownika ma własną konfigurację warstw; zmiana dla jednej roli nie zmienia widoczności dla pozostałych | Kilka ról korzysta z jednej konfiguracji warstw; zmiana obowiązuje wszystkie role objęte współdzieleniem |
| 2 | Rejestr funkcji wyszukiwarki funkcji | Indeks funkcji osiągalnych przez wyszukiwarkę funkcji, obejmujący także elementy warstwy funkcji eksperckich, wraz z zapisem funkcji uruchamianych ostatnio | Rejestr prowadzony jest odrębnie dla każdego zasięgu — funkcje modułu nie pojawiają się w wynikach wyszukiwania innego modułu | Jeden rejestr obejmuje kilka modułów lub środowisk; wyszukiwarka zwraca funkcje całego objętego zakresu |
| 3 | Zasięg przypisania elementu do warstwy | Poziom, na którym obowiązuje przypisanie pojedynczego elementu interfejsu do warstwy widoczności | Przypisanie obowiązuje wyłącznie na wskazanym poziomie i nie przenosi się na pozostałe | Przypisanie obowiązuje wszystkie zasięgi objęte poziomem szerszym |

### 15.2. Zasięg przypisania elementu do warstwy

Przypisanie pojedynczego elementu interfejsu do warstwy widoczności przyjmuje jeden z pięciu zasięgów. Zasięgi uporządkowane są od najogólniejszego do najbardziej szczegółowego; przypisanie z zasięgu bardziej szczegółowego ma pierwszeństwo.

| Zasięg przypisania | Zakres obowiązywania | Pierwszeństwo | Przykład |
|---|---|---|---|
| Globalny | Element należy do wskazanej warstwy w całej platformie | Najniższe — warstwa bazowa | Główne okno komunikacji należy do warstwy zawsze widocznej w każdym module i środowisku |
| Środowisko | Przypisanie obowiązuje wszystkie moduły jednego środowiska | rośnie ↑ | Selektor wykonawcy należy do warstwy widocznej na żądanie w środowisku CodeStudio |
| Moduł | Przypisanie obowiązuje jeden moduł | rośnie ↑ | Narzędzia diagnostyczne modułu Developer należą do warstwy funkcji eksperckich |
| Profil | Przypisanie obowiązuje wszystkie sesje korzystające z danego profilu roli użytkownika | rośnie ↑ | Profil roli podstawowej nie ujawnia trybu administracyjnego |
| Sesja | Przypisanie obowiązuje bieżącą kartę sesji i wygasa z jej zamknięciem | Najwyższe | Jednorazowe ujawnienie panelu diagnostycznego na czas jednego zadania |

```
ROZSTRZYGANIE PRZYPISANIA ELEMENTU DO WARSTWY

  Czy element przypisany w zasięgu SESJA?
        │ tak ─► użyj przypisania z zasięgu Sesja        ◄── koniec
        │ nie
        ▼
  Czy przypisany w zasięgu PROFIL?
        │ tak ─► użyj przypisania z zasięgu Profil       ◄── koniec
        │ nie
        ▼
  Czy przypisany w zasięgu MODUŁ?
        │ tak ─► użyj przypisania z zasięgu Moduł        ◄── koniec
        │ nie
        ▼
  Czy przypisany w zasięgu ŚRODOWISKO?
        │ tak ─► użyj przypisania ze Środowisko          ◄── koniec
        │ nie
        ▼
  Użyj przypisania z zasięgu GLOBALNY (warstwa bazowa)
```

### 15.3. Zależności między warstwami widoczności a pozostałymi punktami izolacji

| Zależność | Treść zależności |
|---|---|
| Warstwy widoczności ↔ izolacja kontekstu | Ukrycie elementu w interfejsie nie zmienia zasięgu historii, pamięci ani kontekstu (rozdz. 4) — warstwa steruje wyłącznie widocznością, nie dostępem do danych |
| Warstwy widoczności ↔ izolacja techniczna | Zakresy izolacji technicznej (rozdz. 3) obowiązują niezależnie od tego, czy przełącznik danego zakresu jest widoczny w bieżącej konfiguracji warstw |
| Warstwy widoczności ↔ komunikacja operacyjna | Główne okno komunikacji i wskaźniki przebiegu pętli wykonawczej należą do warstwy zawsze widocznej; sterowanie przebiegiem i parametry pętli należą do warstw wywoływanych (rozdz. 14) |
| Warstwy widoczności ↔ rola użytkownika | Konfiguracja warstw przypisana roli użytkownika rozstrzyga, które elementy warstw wywoływanych są dla tej roli osiągalne; funkcje warstwy eksperckiej pozostają osiągalne przez wyszukiwarkę funkcji, skrót klawiszowy albo polecenie języka naturalnego |
| Warstwy widoczności ↔ zasada jednego kliknięcia | Każdy ukryty element jest osiągalny jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego, niezależnie od zasięgu przypisania |

---

## 16. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Izolacja kontekstu | Rozdzielenie historii, pamięci i kontekstu między kartami sesji, modułami, środowiskami lub projektami; domyślnie odrębna dla każdej nowej karty sesji, z możliwością skonfigurowania współdzielenia (rozdz. 4, rozdz. 9). |
| Izolacja techniczna | Ograniczenie procesu sesji po stronie serwera, obejmujące osiem zakresów zdefiniowanych w rozdziale 12 Architektury; domyślnie żaden zakres nie jest aktywny (rozdz. 3, rozdz. 12). |
| Zakres izolacji | Pojedyncza, niezależnie ustawialna pozycja macierzy izolacji — jedna z ośmiu pozycji izolacji technicznej albo jedna z trzech pozycji izolacji kontekstu (rozdz. 3, rozdz. 4). |
| Polityka izolacji | Zbiór aktywnych zakresów izolacji obowiązujący dany poziom zasięgu; powstaje przez niezależne łączenie zakresów (rozdz. 5). |
| Poziom zasięgu | Poziom, na którym obowiązuje reguła izolacji: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola lub okno komunikacji; reguły z różnych poziomów rozstrzyga zasada pierwszeństwa zasięgu najbardziej szczegółowego (rozdz. 6). |
| Profil izolacji | Nazwany, zapisywalny zestaw ustawień izolacji kontekstu i izolacji technicznej dla wskazanego poziomu zasięgu, przypisywalny do sesji, roli lub projektu (rozdz. 7, Załącznik B). |
| Polityka efektywna | Wynikowy zestaw reguł izolacji po uwzględnieniu dziedziczenia między poziomami zasięgu, prezentowany w panelu profilu i podglądu okna konfiguracji punktów izolacji (rozdz. 10.2). |
| Warstwa domyślna | Poziom konfiguracji izolacji ustalany w oknie konfiguracji punktów izolacji i obowiązujący przy każdej nowej sesji, karcie, roli lub projekcie (rozdz. 11). |
| Warstwa sesji | Poziom konfiguracji izolacji obejmujący zmiany dokonane dla sesji bieżącej, nakładające się na warstwę domyślną bez jej modyfikacji (rozdz. 11). |
| Chat Window | Główne okno komunikacji Użytkownik ↔ Wykonawca, umieszczone w lewej kolumnie obszaru roboczego; jego historia i pamięć stanowią odrębne punkty izolacji (rozdz. 14). |
| Execution Loop Window | Okno pętli wykonawczej Koordynator ↔ Wykonawca, otwierane jako kolumna sąsiadująca; jego zlecenia, zadania, kolejka i rejestr przebiegu stanowią odrębne punkty izolacji (rozdz. 14). |
| Zlecenie | Jednostka pracy przyjęta przez Koordynatora i rozłożona na zadania; wyznacza granicę współdzielenia kontekstu w pętli wykonawczej (rozdz. 14.3). |
| Rejestr przebiegu pętli | Zapis komunikatów sterujących, wyników kontroli jakości, decyzji o ponowieniu i operacji sterowania przebiegiem pętli wykonawczej (rozdz. 14). |
| Warstwa widoczności | Jedna z czterech warstw, do których należy każdy element interfejsu: zawsze widoczna, widoczna na żądanie, rozwinięcia kontekstowe, funkcje eksperckie (rozdz. 15). |
| Zasięg przypisania do warstwy | Poziom obowiązywania przypisania elementu interfejsu do warstwy widoczności: globalny, środowisko, moduł, profil lub sesja (rozdz. 15.2). |
| Okno konfiguracji punktów izolacji | Trzypanelowe okno (selektor zasięgu, macierz izolacji, profil i podgląd), w którym użytkownik konfiguruje izolację kontekstu i izolację techniczną dla dowolnego poziomu zasięgu (rozdz. 10). |

---

## 17. Kryteria odbioru

Mechanizm izolacji i konfigurowalności zależności jest zgodny z niniejszym opracowaniem, gdy spełnione są łącznie poniższe warunki. Każdy warunek jest sprawdzalny bez odwołania do intencji autora — wyłącznie przez odczyt kontraktu, konfiguracji albo zachowania obserwowalnego z poziomu okna konfiguracji punktów izolacji.

| Warunek | Sposób sprawdzenia |
|---|---|
| Osiem zakresów izolacji technicznej (rozdz. 3) odpowiada dokładnie wyliczeniu `IsolationTechnicalScope` z `budowa/shared/contract.json`, bez zakresu dodanego lub pominiętego | `python3 -c "import json; d=json.load(open('budowa/shared/contract.json')); print(len([w for w in d['wyliczenia'] if w['nazwa']=='IsolationTechnicalScope'][0]['wartosci']))"` zwraca `8` |
| Domyślnie żaden zakres izolacji technicznej nie jest aktywny (rozdz. 3, rozdz. 12) | odczyt polityki efektywnej dla sesji bez profilu przypisanego komendą `isolation.policy.preview` zwraca wszystkie przełączniki `IsolationTechnicalSwitch` w stanie wyłączonym |
| Trzy zakresy izolacji kontekstu dla karty sesji są domyślnie odrębne (rozdz. 4, rozdz. 12) | odczyt `isolation.context.get` dla nowej karty sesji bez ustawienia zwraca `separate` dla historii, pamięci i kontekstu |
| Osiem poziomów zasięgu wymienionych w rozdz. 6 rozstrzyga regułę izolacji zgodnie z zasadą pierwszeństwa poziomu najbardziej szczegółowego | zmiana ustawienia na poziomie węższym — karta sesji — przesłania wartość odziedziczoną z poziomu szerszego, czyli środowiska, bez zmiany tej ostatniej |
| Profil izolacji zapisany komendą `isolation.profile.save` jest przypisywalny do sesji, roli lub projektu bez ograniczenia liczby przypisań (rozdz. 7–8) | `isolation.profile.assign` z tym samym `profileId` i różnymi `scopeId` kończy się powodzeniem dla każdego wywołania |
| Brak ustawienia na danym poziomie dziedziczy wartość z poziomu bezpośrednio szerszego, nigdy błędem (zasada „brak ustawienia = wartość domyślna”) | odczyt `isolation.policy.preview` dla poziomu bez żadnego zapisanego ustawienia zwraca politykę odziedziczoną, nie kod błędu |
| Każda z 12 komend obszaru `isolation` odpowiada wykazowi Załącznika G bez nazw wymyślonych | `grep -oE '"typ": *"isolation\.[a-zA-Z.]+"' budowa/shared/contract.json \| sort -u` zwraca dokładnie 12 pozycji zgodnych z Załącznikiem G |
| Żaden zakres izolacji nie blokuje twardo działania — każdy jest ustawieniem odwracalnym na dowolnym poziomie zasięgu (zasada zero blokad) | wyłączenie dowolnego przełącznika technicznego lub kontekstu nie usuwa profilu ani nie wymaga ponownego uruchomienia sesji |
| Historia i pamięć Chat Window oraz zlecenia, zadania, kolejka i rejestr przebiegu Execution Loop Window są punktami izolacji odrębnymi od ośmiu zakresów technicznych (rozdz. 14) | zmiana ustawienia izolacji technicznej nie zmienia stanu współdzielenia historii Chat Window i odwrotnie |
| Poziom `window` uczestniczy w rozstrzyganiu polityki izolacji jako najwęższy, a poziom `application` w nim nie uczestniczy (rozdz. 6) | lista zwracana w polu wyniku `scopes:IsolationScopeLevel[]` komendy `isolation.scope.list` liczy osiem pozycji, zawiera `window` na ostatnim miejscu kolejności rozstrzygania i nie zawiera `application` |

---

## Załącznik A. Macierz zakresów izolacji technicznej z objaśnieniami kontekstowymi

Poniższa tabela zbiera jedenaście zakresów izolacji procesu i kontekstu (osiem technicznych, trzy kontekstu) wraz z treścią objaśnienia kontekstowego `[?]`, zgodnie z wymogiem rozdziału 13 Architektury, że treść objaśnienia jest częścią definicji ustawienia. Objaśnienia pozostałych ośmiu zakresów zawierają Załączniki D i E.

| Zakres | Rodzaj | Stan domyślny | Przykładowa treść objaśnienia `[?]` |
|---|---|---|---|
| Katalog roboczy sesji | Techniczny | Wyłączony | „Gdy włączone, ta sesja pracuje we własnym katalogu plików, niewidocznym dla innych sesji. Gdy wyłączone, katalog jest współdzielony z innymi sesjami tego samego zasięgu.” |
| Środowisko procesu | Techniczny | Wyłączony | „Gdy włączone, proces tej sesji otrzymuje własny zestaw zmiennych środowiskowych. Gdy wyłączone, dzieli środowisko uruchomieniowe z innymi procesami serwera.” |
| Katalog danych i konfiguracji modelu | Techniczny | Wyłączony | „Gdy włączone, dane i konfiguracja wykorzystywanego modelu są przechowywane odrębnie dla tej sesji. Gdy wyłączone, katalog jest współdzielony z innymi sesjami korzystającymi z tego samego modelu.” |
| Dostęp sieciowy | Techniczny | Wyłączony | „Gdy włączone, ta sesja otrzymuje odrębny, ograniczony dostęp sieciowy. Gdy wyłączone, korzysta ze wspólnego dostępu sieciowego serwera.” |
| Zakres odczytu i zapisu plików | Techniczny | Wyłączony | „Gdy włączone, dostęp do plików poza katalogiem roboczym tej sesji ograniczony jest do jawnie dozwolonych ścieżek. Gdy wyłączone, sesja ma pełny dostęp do zasobów plikowych dostępnych na serwerze.” |
| Konto i token per sesja | Techniczny | Wyłączony | „Gdy włączone, ta sesja korzysta z własnego, dedykowanego konta lub tokenu dostępu. Gdy wyłączone, korzysta ze wspólnych, platformowych danych dostępowych.” |
| Model procesu | Techniczny | Wyłączony (współdzielony) | „Gdy włączone, ta sesja uruchamia własną, niezależną instancję procesu modelu. Gdy wyłączone, korzysta ze wspólnej puli procesów modelu.” |
| Serwer wykonania | Techniczny | Wyłączony | „Gdy włączone, ta sesja wykonywana jest na dedykowanym serwerze, odrębnym od pozostałych. Gdy wyłączone, wykonywana jest na współdzielonym serwerze wykonawczym platformy.” |
| Historia | Kontekstu | Odrębna | „Gdy współdzielona, zapis rozmowy jest widoczny jednocześnie w kilku kartach, modułach lub środowiskach. Gdy odrębna, każda nowa karta zaczyna z pustą historią.” |
| Pamięć | Kontekstu | Odrębna | „Gdy współdzielona, ten sam zasób pamięci zasila kilka zasięgów jednocześnie. Gdy odrębna, pamięć przypisana jednemu zasięgowi nie jest widoczna w innym.” |
| Kontekst | Kontekstu | Odrębny | „Gdy współdzielony, bieżący stan roboczy przenoszony jest między kartami, modułami lub środowiskami bez przeładowania. Gdy odrębny, kontekst jest rekonfigurowany od nowa przy każdej zmianie środowiska lub modułu.” |

---

## Załącznik B. Szablon profilu izolacji

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Nazwa profilu | Nazwa własna profilu | dowolny tekst; przykład: „Executor zdalny — SSH” |
| Opis | Krótkie uzasadnienie przeznaczenia profilu | dowolny tekst |
| Poziom zasięgu | Poziom obowiązywania reguły | globalny \| środowisko \| moduł \| para modułów \| projekt \| karta sesji \| rola \| okno komunikacji (rozdz. 6) |
| Izolacja kontekstu — historia | Współdzielenie historii | współdzielona \| odrębna |
| Izolacja kontekstu — pamięć | Współdzielenie pamięci | współdzielona \| odrębna |
| Izolacja kontekstu — kontekst | Współdzielenie kontekstu | współdzielony \| odrębny |
| Izolacja techniczna | Osiem zakresów z rozdziału 3 | każdy zakres: włączony \| wyłączony |
| Izolacja komunikacji operacyjnej | Pięć pozycji z rozdziału 14 | każda pozycja: współdzielona \| odrębna |
| Izolacja warstw widoczności — konfiguracja warstw | Konfiguracja warstw wg roli użytkownika (rozdz. 15) | współdzielona \| odrębna |
| Izolacja warstw widoczności — rejestr wyszukiwarki | Rejestr funkcji wyszukiwarki funkcji (rozdz. 15) | współdzielony \| odrębny |
| Izolacja warstw widoczności — zasięg przypisania | Zasięg przypisania elementu do warstwy (rozdz. 15.2) | globalny \| środowisko \| moduł \| profil \| sesja |
| Warstwa | Warstwa konfiguracji | domyślna \| sesji (rozdz. 11) |
| Przypisanie | Element, do którego profil jest przypisany | sesja \| rola \| projekt (rozdz. 8) |
| Wartość domyślna | Zachowanie przy braku ustawienia | dziedziczenie z poziomu szerszego; brak aktywnej izolacji technicznej (rozdz. 12) |

---

## Załącznik C. Tabela poziomów zasięgu i pierwszeństwa

| Poziom zasięgu | Pierwszeństwo | Miejsce konfiguracji | Dziedziczy z |
|---|---|---|---|
| Globalny | Najniższe — warstwa bazowa | Panel selektora zasięgu, poziom „Globalny” | — |
| Środowisko | ↑ | Panel selektora zasięgu, poziom „Środowisko” | Globalny |
| Moduł | ↑ | Panel selektora zasięgu, poziom „Moduł” | Środowisko |
| Para modułów | ↑ | Panel selektora zasięgu, poziom „Para modułów” | Moduł |
| Projekt | ↑ | Panel selektora zasięgu lub atrybut projektu w module Workspace | Para modułów |
| Karta sesji | ↑ | Uproszczone menu kontekstowe karty sesji lub panel profilu, warstwa „sesji” | Projekt |
| Rola | ↑ | Okno roli w panelu orkiestracji środowiska MultitaskingAI | Karta sesji |
| Okno komunikacji | Najwyższe — poziom najwęższy | Panel selektora zasięgu, poziom „Okno komunikacji”, lub komponent konfiguracji okna Chat Window albo Execution Loop Window | Rola |

---

## Załącznik D. Zestawienie punktów izolacji warstwy komunikacji operacyjnej

| Pozycja | Kanał | Stan domyślny | Treść objaśnienia `[?]` |
|---|---|---|---|
| Historia Chat Window | Użytkownik ↔ Wykonawca | Odrębna | „Gdy współdzielona, ten sam zapis rozmowy z Wykonawcą widoczny jest w kilku kartach, modułach lub środowiskach. Gdy odrębna, każda karta sesji prowadzi własny zapis rozmowy.” |
| Pamięć Chat Window | Użytkownik ↔ Wykonawca | Odrębna | „Gdy współdzielona, ten sam zasób pamięci zasila główne okno komunikacji w kilku zasięgach. Gdy odrębna, pamięć jednej karty sesji nie jest widoczna w innej.” |
| Zlecenia i zadania | Koordynator ↔ Wykonawca | Odrębne | „Gdy współdzielone, zadania kilku zleceń korzystają ze wspólnego zestawu ustaleń i wspólnego kontekstu. Gdy odrębne, każde zlecenie prowadzi własny zestaw zadań; zadania jednego zlecenia zawsze współdzielą jego kontekst.” |
| Kolejka pętli wykonawczej | Koordynator ↔ Wykonawca | Odrębna | „Gdy współdzielona, kilka zleceń obsługiwanych jest z jednej kolejki o wspólnym porządku wykonania. Gdy odrębna, każde zlecenie ma własną kolejkę.” |
| Rejestr przebiegu pętli | Koordynator ↔ Wykonawca | Odrębny | „Gdy współdzielony, jeden rejestr obejmuje przebieg kilku zleceń, modułów lub środowisk. Gdy odrębny, rejestr prowadzony jest osobno dla każdego zlecenia.” |

---

## Załącznik E. Zestawienie punktów izolacji warstw widoczności

| Pozycja | Stan domyślny | Zasięgi ustawienia | Treść objaśnienia `[?]` |
|---|---|---|---|
| Konfiguracja warstw według roli użytkownika | Odrębna dla każdej roli | globalny, środowisko, moduł, profil, sesja | „Gdy współdzielona, kilka ról użytkownika korzysta z jednej konfiguracji warstw widoczności. Gdy odrębna, każda rola ma własny zestaw reguł widoczności elementów interfejsu.” |
| Rejestr funkcji wyszukiwarki funkcji | Współdzielony w obrębie środowiska | globalny, środowisko, moduł, profil, sesja | „Gdy współdzielony, wyszukiwarka funkcji zwraca funkcje wszystkich objętych modułów i środowisk. Gdy odrębny, wyniki ograniczone są do bieżącego zasięgu.” |
| Zasięg przypisania elementu do warstwy | Globalny | globalny, środowisko, moduł, profil, sesja | „Zasięg wskazuje poziom, na którym obowiązuje przypisanie elementu do warstwy widoczności; przypisanie z zasięgu bardziej szczegółowego ma pierwszeństwo przed przypisaniem z zasięgu szerszego.” |

---

## Załącznik F. Struktury danych pełne

Wykaz przenosi z `budowa/shared/contract.json` pełne definicje pól struktur używanych przez komendy obszaru `isolation`. Kolumna „Wymagane” podaje wartość pola `wymagane` struktury źródłowej.

### F.1. `IsolationPolicy` — polityka izolacji obowiązująca poziom zasięgu

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `scope` | `ConfigScope` | tak | Poziom zasięgu, dla którego politykę wyliczono |
| `scopeId` | `string` | nie | Identyfikator bytu poziomu; puste dla `global` |
| `layer` | `IsolationLayer` | tak | Warstwa obowiązująca — `default` albo `session` |
| `contextSwitches` | `IsolationSwitch[]` | tak | Przełączniki izolacji kontekstu |
| `technicalSwitches` | `IsolationTechnicalSwitch[]` | tak | Przełączniki izolacji technicznej |
| `profileId` | `string` | nie | Profil, z którego polityka pochodzi |
| `origin` | `ConfigScope` | nie | Poziom, z którego wartość została odziedziczona |

### F.2. `IsolationProfile` — profil zapisany i przypisywalny

| Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|
| `id` | `string` | tak | Identyfikator profilu |
| `name` | `string` | tak | Nazwa profilu |
| `description` | `string` | nie | Opis profilu |
| `contextSwitches` | `IsolationSwitch[]` | tak | Przełączniki izolacji kontekstu |
| `technicalSwitches` | `IsolationTechnicalSwitch[]` | tak | Przełączniki izolacji technicznej |
| `createdAt` | `int64` | tak | Czas utworzenia w milisekundach epoki |
| `updatedAt` | `int64` | tak | Czas ostatniej zmiany w milisekundach epoki |

### F.3. `IsolationScopeLevel`, `IsolationSwitch`, `IsolationTechnicalSwitch`

| Struktura | Pole | Typ | Wymagane | Znaczenie |
|---|---|---|---|---|
| `IsolationScopeLevel` | `scope` | `ConfigScope` | tak | Poziom zasięgu |
| `IsolationScopeLevel` | `label` | `string` | tak | Nazwa poziomu do wyświetlenia |
| `IsolationScopeLevel` | `order` | `int` | tak | Miejsce w kolejności rozstrzygania; `1` znaczy najszerszy |
| `IsolationScopeLevel` | `narrowest` | `bool` | nie | Czy poziom jest najwęższy i wygrywa rozstrzyganie |
| `IsolationScopeLevel` | `description` | `string` | nie | Objaśnienie poziomu do znacznika kontekstowego `[?]` |
| `IsolationSwitch` | `kind` | `IsolationContextKind` | tak | Rodzaj kontekstu |
| `IsolationSwitch` | `isolated` | `bool` | tak | Czy kontekst jest odcinany |
| `IsolationSwitch` | `explanation` | `string` | nie | Objaśnienie przełącznika do znacznika kontekstowego `[?]` |
| `IsolationTechnicalSwitch` | `scope` | `IsolationTechnicalScope` | tak | Zakres techniczny |
| `IsolationTechnicalSwitch` | `isolated` | `bool` | tak | Czy zakres jest odcinany |
| `IsolationTechnicalSwitch` | `explanation` | `string` | nie | Objaśnienie przełącznika do znacznika kontekstowego `[?]` |

### F.4. Wyliczenia obszaru izolacji

| Wyliczenie | Wartości | Znaczenie |
|---|---|---|
| `IsolationTechnicalScope` | `workingDirectory` · `processEnvironment` · `modelDataDirectory` · `networkAccess` · `fileAccess` · `accountToken` · `processModel` · `executionServer` | Osiem zakresów technicznych izolacji (rozdz. 3): katalog roboczy, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, odczyt i zapis plików, konto i token per sesja, model procesu, serwer wykonania |
| `IsolationLayer` | `default` · `session` | Warstwa, na której obowiązuje ustawienie izolacji: domyślna platformy albo warstwa karty sesji (rozdz. 11) |
| `IsolationContextKind` | `history` · `memory` · `context` | Trzy zakresy izolacji kontekstu (rozdz. 4): historia rozmowy, pamięć projektu, kontekst przekazywany między modułami |
| `ConfigScope` | `application` · `global` · `environment` · `module` · `modulePair` · `project` · `session` · `role` · `window` | Dziewięć poziomów zasięgu konfiguracji całego kontraktu; `window` jest najwęższy i wygrywa rozstrzyganie, `application` najszerszy i ustępuje każdemu innemu. Izolacja stosuje osiem z tych wartości — bez `application` (rozdz. 6) |

Pole `kind` struktury `IsolationSwitch` (F.3) przyjmuje wyłącznie wartości `IsolationContextKind`; każdy z trzech przełączników zwracanych komendą `isolation.context.get` odpowiada dokładnie jednej z tych wartości, a macierz izolacji okna konfiguracji punktów izolacji (rozdz. 10) prezentuje je jako trzy odrębne wiersze niezależnie od ośmiu wierszy izolacji technicznej. Analogicznie pole `scope` struktury `IsolationTechnicalSwitch` przyjmuje wyłącznie wartości `IsolationTechnicalScope` — dwa wyliczenia nigdy nie mieszają się w jednej strukturze przełącznika.

---

## Załącznik G. Pełny wykaz komend kontraktu punktów izolacji

Wykaz wygenerowany wprost z `budowa/shared/contract.json` — jedynego źródła nazw komend, pól ładunków i zdarzeń. Każda pozycja jest wiążąca dla warstwy klienckiej i serwerowej; deweloper nie dodaje ani nie zmienia nazw poza tym wykazem.

### Obszar `isolation` — 12 komend

| Komenda | Przeznaczenie | Pola żądania | Pola wyniku |
|---|---|---|---|
| `isolation.scope.list` | Zwraca osiem poziomów zasięgu izolacji w kolejności rozstrzygania | `sessionId:string` (opc)<br>`windowId:string` (opc) | `scopes:IsolationScopeLevel[]` (wym) |
| `isolation.context.get` | Odczytuje przełączniki izolacji kontekstu zapisane na wskazanym poziomie i warstwie | `scope:ConfigScope` (wym)<br>`scopeId:string` (opc)<br>`layer:IsolationLayer` (opc) | `switches:IsolationSwitch[]` (wym) |
| `isolation.context.set` | Zapisuje przełączniki izolacji kontekstu na wskazanym poziomie i warstwie | `scope:ConfigScope` (wym)<br>`scopeId:string` (opc)<br>`layer:IsolationLayer` (opc)<br>`switches:IsolationSwitch[]` (wym) | `switches:IsolationSwitch[]` (wym) |
| `isolation.technical.get` | Odczytuje przełączniki ośmiu zakresów technicznych zapisane na wskazanym poziomie i warstwie | `scope:ConfigScope` (wym)<br>`scopeId:string` (opc)<br>`layer:IsolationLayer` (opc) | `switches:IsolationTechnicalSwitch[]` (wym) |
| `isolation.technical.set` | Zapisuje przełączniki ośmiu zakresów technicznych na wskazanym poziomie i warstwie | `scope:ConfigScope` (wym)<br>`scopeId:string` (opc)<br>`layer:IsolationLayer` (opc)<br>`switches:IsolationTechnicalSwitch[]` (wym) | `switches:IsolationTechnicalSwitch[]` (wym) |
| `isolation.profile.save` | Zapisuje profil izolacji; puste `profileId` zakłada nowy, podane zmienia istniejący | `profileId:string` (opc)<br>`name:string` (wym)<br>`description:string` (opc)<br>`contextSwitches:IsolationSwitch[]` (opc)<br>`technicalSwitches:IsolationTechnicalSwitch[]` (opc) | `profile:IsolationProfile` (wym) |
| `isolation.profile.list` | Zwraca profile izolacji | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc) | `profiles:IsolationProfile[]` (wym) |
| `isolation.profile.load` | Wczytuje profil izolacji do panelu bez przypisywania go do poziomu | `profileId:string` (wym) | `profile:IsolationProfile` (wym) |
| `isolation.profile.assign` | Przypisuje profil izolacji do wskazanego poziomu zasięgu i warstwy | `profileId:string` (wym)<br>`scope:ConfigScope` (wym)<br>`scopeId:string` (opc)<br>`layer:IsolationLayer` (opc) | `policy:IsolationPolicy` (wym) |
| `isolation.profile.delete` | Usuwa profil izolacji | `profileId:string` (wym) | `deleted:bool` (wym) |
| `isolation.layer.set` | Przełącza warstwę izolacji między domyślną platformy a warstwą karty sesji | `layer:IsolationLayer` (wym)<br>`sessionId:string` (opc)<br>`windowId:string` (opc) | `layer:IsolationLayer` (wym) |
| `isolation.policy.preview` | Zwraca politykę izolacji obowiązującą po rozstrzygnięciu ośmiu poziomów zasięgu, bez zapisu | `scope:ConfigScope` (opc)<br>`scopeId:string` (opc)<br>`sessionId:string` (opc)<br>`windowId:string` (opc)<br>`layer:IsolationLayer` (opc) | `policy:IsolationPolicy` (wym) |

**Zdarzenia obszaru `isolation` — 2:**

| Zdarzenie | Przeznaczenie | Pola |
|---|---|---|
| `isolation.profile.changed` | Zmiana profilu izolacji | — |
| `isolation.policy.changed` | Zmiana polityki izolacji obowiązującej okno | — |

Razem w wykazie: **12 komend** z 1 obszaru kontraktu.

---

*Koniec dokumentu. Izolacja i konfigurowalność zależności — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
