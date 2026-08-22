# Danaco Console — Koncepcja platformy

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
| **Tytuł** | Koncepcja platformy |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (główny) · projektant |
| **Przeznaczenie** | Ustala koncepcyjny, docelowy kształt platformy — charakter platformy, pojęcia podstawowe, warstwy architektury i widoczności, strukturę środowisk i modułów oraz zasady nadrzędne funkcjonalności — na którym opierają się wszystkie pozostałe opracowania zbioru |
| **Zakres** | charakter platformy, pojęcia podstawowe, komunikacja operacyjna, warstwy widoczności interfejsu, warstwy architektury, konfigurowalność zależności i punktów izolacji, strona główna i centrum dowodzenia, nawigacja i sesje, środowiska, macierz dostępności modułów, piętnaście modułów, funkcje globalne, środowisko MultitaskingAI, zasady nadrzędne funkcjonalności |
| **Poza zakresem** | szczegóły techniczne wdrożenia i stosu technologicznego — [Architektura techniczna](architektura.md); pełne kontrakty komunikacji — [Kontrakty komunikacji](kontrakty-komunikacji.md); wygląd i zachowanie pojedynczych okien — [Elementy okien](../interfejs-uzytkownika/elementy-okien.md) |
| **Dokument nadrzędny** | [Spis opracowań](../SPIS-OPRACOWAN.md) |
| **Dokumenty powiązane** | [Architektura techniczna](architektura.md) · [Model danych](model-danych.md) · [Model konfiguracji](model-konfiguracji.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Izolacja i zależności](izolacja-i-zaleznosci.md) · [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) · [Rozszerzenia](rozszerzenia.md) |
| **Prototypy odniesienia** | `design/05-okna/` — komplet 38 prototypów |
| **Źródła normatywne** | `budowa/shared/contract.json` (68 obszarów, przekrój pełny) · `design/05-okna/` (38 prototypów) |
| **Zasada nadrzędna** | Platforma jest jednym, spójnym środowiskiem operacyjnym AI — piętnaście modułów i cztery funkcje globalne dzielą wspólny model komunikacji, danych i konfiguracji; żaden moduł nie definiuje własnego, odrębnego mechanizmu tam, gdzie działa mechanizm wspólny |

Dokument opisuje koncepcję produktu Danaco Console: charakter platformy, pojęcia podstawowe, komunikację operacyjną, warstwy widoczności interfejsu, warstwy architektury, konfigurowalność zależności i punktów izolacji, stronę główną i centrum dowodzenia, nawigację i sesje, środowiska, macierz dostępności modułów, moduły, funkcje globalne, środowisko MultitaskingAI oraz zasady nadrzędne funkcjonalności. Opracowanie przedstawia kompletny, docelowy kształt platformy w formie tabel zestawczych, schematów tekstowych, szablonów konfiguracji i scenariuszy użycia. Nazwy własne pochodzą od Operatora.

---

## Spis treści

1. [Streszczenie zarządcze](#streszczenie-zarządcze)
2. [Wprowadzenie](#wprowadzenie)
   - [Pozycjonowanie produktu](#pozycjonowanie-produktu)
   - [Uzasadnienie architektury: środowisko → moduł → okno](#uzasadnienie-architektury-środowisko--moduł--okno)
3. [Charakter platformy](#1-charakter-platformy)
4. [Pojęcia podstawowe](#2-pojęcia-podstawowe)
   - [2.1 Środowisko](#21-środowisko)
   - [2.2 Moduł](#22-moduł)
   - [2.3 Komponent własny](#23-komponent-własny)
   - [2.4 Czat](#24-czat)
5. [Komunikacja operacyjna](#3-komunikacja-operacyjna)
   - [3.1 Kanał pierwszy — Chat Window (Użytkownik ↔ Wykonawca)](#31-kanał-pierwszy--chat-window-użytkownik--wykonawca)
   - [3.2 Kanał drugi — Execution Loop Window (Koordynator ↔ Wykonawca)](#32-kanał-drugi--execution-loop-window-koordynator--wykonawca)
   - [3.3 Role uczestników komunikacji](#33-role-uczestników-komunikacji)
   - [3.4 Rozmieszczenie kanałów w obszarze roboczym](#34-rozmieszczenie-kanałów-w-obszarze-roboczym)
6. [Warstwy widoczności interfejsu](#4-warstwy-widoczności-interfejsu)
   - [4.1 Cztery warstwy widoczności](#41-cztery-warstwy-widoczności)
   - [4.2 Mechanizmy ukrywania funkcjonalności](#42-mechanizmy-ukrywania-funkcjonalności)
   - [4.3 Zasada jednego kliknięcia](#43-zasada-jednego-kliknięcia)
   - [4.4 Stan spoczynku interfejsu](#44-stan-spoczynku-interfejsu)
7. [Warstwy architektury](#5-warstwy-architektury)
   - [5.1 Warstwa środowisk (Environment Layer)](#51-warstwa-środowisk-environment-layer)
   - [5.2 Warstwa modułów (Workspace Layer)](#52-warstwa-modułów-workspace-layer)
   - [5.3 Warstwa funkcji globalnych (Global Features Layer)](#53-warstwa-funkcji-globalnych-global-features-layer)
   - [5.4 Schemat hierarchii warstw](#54-schemat-hierarchii-warstw)
8. [Konfigurowalność zależności i punktów izolacji](#6-konfigurowalność-zależności-i-punktów-izolacji)
   - [6.1 Punkty izolacji jako decyzja użytkownika](#61-punkty-izolacji-jako-decyzja-użytkownika)
   - [6.2 Zależności jako możliwość, nie wbudowana reguła](#62-zależności-jako-możliwość-nie-wbudowana-reguła)
   - [6.3 Konsekwencje dla porządku pracy](#63-konsekwencje-dla-porządku-pracy)
   - [6.4 Okno konfiguracji punktów izolacji](#64-okno-konfiguracji-punktów-izolacji)
   - [6.5 Poziomy zasięgu reguł izolacji](#65-poziomy-zasięgu-reguł-izolacji)
   - [6.6 Profile izolacji i warstwy konfiguracji](#66-profile-izolacji-i-warstwy-konfiguracji)
9. [Strona główna i Centrum dowodzenia](#7-strona-główna-i-centrum-dowodzenia)
   - [7.1 Strefa 1 — wybór środowiska](#71-strefa-1--wybór-środowiska)
   - [7.2 Strefa 2 — komponenty własne](#72-strefa-2--komponenty-własne)
   - [7.3 Strefa 3 — ustawienia](#73-strefa-3--ustawienia)
   - [7.4 Forma prezentacji](#74-forma-prezentacji)
10. [Nawigacja i sesje](#8-nawigacja-i-sesje)
11. [Środowiska](#9-środowiska)
   - [9.1 TalkIn](#91-talkin)
   - [9.2 WorkSpace](#92-workspace)
   - [9.3 CodeStudio](#93-codestudio)
   - [9.4 MultitaskingAI](#94-multitaskingai)
12. [Macierz dostępności modułów](#10-macierz-dostępności-modułów)
13. [Moduły](#11-moduły)
   - [11.1 Studio](#111-studio)
   - [11.2 Workspace](#112-workspace)
   - [11.3 Automations](#113-automations)
   - [11.4 Browser](#114-browser)
   - [11.5 Research](#115-research)
   - [11.6 Library](#116-library)
   - [11.7 Translate](#117-translate)
   - [11.8 Roundtable](#118-roundtable)
   - [11.9 Design](#119-design)
   - [11.10 Assistant](#1110-assistant)
   - [11.11 Terminal](#1111-terminal)
   - [11.12 Developer](#1112-developer)
   - [11.13 Diagnostics](#1113-diagnostics)
   - [11.14 Apps](#1114-apps)
   - [11.15 Agents](#1115-agents)
14. [Funkcje globalne](#12-funkcje-globalne)
   - [12.1 Mobile](#121-mobile)
   - [12.2 Always On Display](#122-always-on-display)
15. [MultitaskingAI — rozwinięcie środowiska](#13-multitaskingai--rozwinięcie-środowiska)
   - [13.1 Definicja](#131-definicja)
   - [13.2 Warstwa centralna](#132-warstwa-centralna)
   - [13.3 Role](#133-role)
   - [13.4 Silnik kolejek](#134-silnik-kolejek)
   - [13.5 Orkiestracja](#135-orkiestracja)
   - [13.6 Integracja z Automations](#136-integracja-z-automations)
   - [13.7 Kluczowa idea](#137-kluczowa-idea)
   - [13.8 Boczna nawigacja środowiska MultitaskingAI (panel orkiestracji)](#138-boczna-nawigacja-środowiska-multitaskingai-panel-orkiestracji)
   - [13.9 Schemat przepływu orkiestracji](#139-schemat-przepływu-orkiestracji)
16. [Zasady nadrzędne funkcjonalności](#14-zasady-nadrzędne-funkcjonalności)
17. [Słowniczek pojęć](#15-słowniczek-pojęć)
18. [Kryteria odbioru](#16-kryteria-odbioru)
19. [Załącznik A. Macierz okien operacyjnych per moduł](#załącznik-a-macierz-okien-operacyjnych-per-moduł)
20. [Załącznik B. Zestawienie ról środowiska MultitaskingAI](#załącznik-b-zestawienie-ról-środowiska-multitaskingai)
21. [Załącznik C. Macierz funkcji globalnych](#załącznik-c-macierz-funkcji-globalnych)
22. [Załącznik D. Szablony konfiguracji](#załącznik-d-szablony-konfiguracji)
   - [D.1. Szablon komponentu własnego](#d1-szablon-komponentu-własnego)
   - [D.2. Szablon roli (środowisko MultitaskingAI)](#d2-szablon-roli-środowisko-multitaskingai)
   - [D.3. Szablon kolejki](#d3-szablon-kolejki)
   - [D.4. Szablon profilu izolacji](#d4-szablon-profilu-izolacji)
23. [Załącznik E. Scenariusze użycia](#załącznik-e-scenariusze-użycia)
   - [E.1. Autonomiczna budowa aplikacji w środowisku MultitaskingAI](#e1-autonomiczna-budowa-aplikacji-w-środowisku-multitaskingai)
   - [E.2. Przepływ badawczy między modułami](#e2-przepływ-badawczy-między-modułami)
   - [E.3. Konfiguracja punktu izolacji między projektami](#e3-konfiguracja-punktu-izolacji-między-projektami)
   - [E.4. Praca dwujęzyczna Studio i Translate](#e4-praca-dwujęzyczna-studio-i-translate)

---

## Streszczenie zarządcze

Danaco Console jest systemem operacyjnym do pracy z modelami sztucznej inteligencji — nie aplikacją czatową, lecz środowiskiem organizującym pracę wokół dedykowanych przestrzeni roboczych dopasowanych do rodzaju zadania. Punktem wejścia jest strona główna (centrum dowodzenia), z której użytkownik wybiera jedno z czterech środowisk, konfiguruje własne komponenty lub przechodzi do ustawień. Zamiast jednego uniwersalnego okna rozmowy platforma udostępnia hierarchię **środowisko → moduł → okno operacyjne**, ponad którą działają funkcje globalne i warstwa agentowa. Pracą platformy steruje się dwoma kanałami komunikacji operacyjnej — Chat Window (Użytkownik ↔ Wykonawca) oraz Execution Loop Window (Koordynator ↔ Wykonawca) — rozmieszczonymi w układzie pionowym obszaru roboczego (rozdz. 3), a interfejs ujawnia możliwości systemu stopniowo, w czterech warstwach widoczności (rozdz. 4). Bloki konstrukcyjne platformy zbiera poniższa tabela.

| Blok konstrukcyjny | Czym jest | Elementy | Rozdział |
|---|---|---|---|
| Komunikacja operacyjna | Dwa kanały komunikacji sterujące pracą platformy | Chat Window (Użytkownik ↔ Wykonawca), Execution Loop Window (Koordynator ↔ Wykonawca) | 3 |
| Warstwy widoczności interfejsu | Stopniowe ujawnianie możliwości systemu | Cztery warstwy widoczności, mechanizmy ukrywania, zasada jednego kliknięcia | 4 |
| Środowiska | Cztery odrębne tryby pracy z AI | TalkIn, WorkSpace, CodeStudio, MultitaskingAI | 9 |
| Moduły | Wyspecjalizowane obszary robocze wewnątrz środowisk | Piętnaście modułów — od edycji dokumentów po budowę agentów | 11 |
| Funkcje globalne | Warstwa działająca ponad środowiskami i modułami | Mobile ([Funkcje mobilne](../funkcje-globalne/mobile.md)), Always On Display | 12 |
| MultitaskingAI | Najbardziej zaawansowane konfiguracyjnie środowisko | Role (wykonawca, koordynator, walidator), silnik kolejek, orkiestracja, pętla pracy ciągłej 24/7/365 | 13 |
| Zasady nadrzędne | Wspólny standard całej platformy | Pełna kompozycyjność, pełna konfigurowalność (centralna), jawność zależności, pełna orkiestracja, warstwy widoczności interfejsu | 14 |

Zasadą centralną jest pełna konfigurowalność: wszystkie zależności platformy — między środowiskami, modułami, komponentami własnymi i sesjami — konfiguruje się z okna konfiguracji, a punkty izolacji kontekstu ustala użytkownik, nie sztywna reguła systemu. Po spięciu środowiska MultitaskingAI z modułem Automations i przypisaniu ról własnym agentom powstaje w pełni autonomiczna, wieloagentowa pętla pracy ciągłej, zdolna działać 24 godziny na dobę, 7 dni w tygodniu.

---

## Wprowadzenie

### Pozycjonowanie produktu

Danaco Console nie jest aplikacją czatową rozszerzoną o dodatkowe funkcje — to rozróżnienie jest fundamentem koncepcji i determinuje wszystkie decyzje architektoniczne dokumentu. Dwa modele pracy zestawia poniższa tabela.

| Wymiar | Klasyczna aplikacja czatowa | Danaco Console |
|---|---|---|
| Punkt wyjścia | Jedno, względnie stałe okno rozmowy | Zadanie użytkownika, dla którego dobierana jest przestrzeń robocza |
| Dokładanie możliwości | Kolejne funkcje w tym samym oknie (załączniki, wtyczki, tryby, panele) | Osobna, dopasowana przestrzeń robocza dla każdego rodzaju pracy |
| Kontekst konwersacyjny | Wspólny i niezmienny dla wszystkich zadań | Rekonfigurowany przy każdej zmianie kontekstu pracy |
| Rola czatu | Stały komponent interfejsu | Uniwersalna warstwa komunikacji, obecna wszędzie, lecz nie stała |
| Co zmienia się wraz z zadaniem | Wyłącznie treść rozmowy | Układ okien, dostępne narzędzia, historia sesji, pamięć kontekstowa, sposób działania AI |

### Uzasadnienie architektury: środowisko → moduł → okno

Koncepcja opiera się na trójstopniowej hierarchii **środowisko → moduł → okno operacyjne**. Wybór tej struktury zamiast płaskiego modelu jednego okna czatu wynika z czterech przesłanek.

| Przesłanka | Na czym polega |
|---|---|
| Dopasowanie przestrzeni do zadania | Różne rodzaje pracy z AI wymagają różnych układów interfejsu i narzędzi: praca nad kodem — edytora, drzewa projektu i terminala; praca badawcza — zarządzania źródłami i budowania raportu; praca nad automatyzacją — harmonogramu, kolejki i monitora wykonania. Płaski czat nie pomieści tej różnorodności bez przeciążenia; hierarchia daje każdemu rodzajowi pracy dokładnie potrzebną przestrzeń. |
| Konfigurowalne punkty izolacji | Środowiska i moduły dają odrębne przestrzenie robocze (własny układ, historia, zestaw narzędzi), które użytkownik może utrzymywać oddzielnie albo świadomie łączyć z okna konfiguracji, tak aby kontekst jednego zadania nie zanieczyszczał innego wbrew jego woli. |
| Skalowalność | Nowe rodzaje pracy dodaje się jako kolejne moduły wewnątrz istniejących środowisk, bez naruszania spójności istniejących przestrzeni i bez rozrostu jednego, współdzielonego interfejsu. |
| Miejsce dla warstwy nadrzędnej | Hierarchia tworzy naturalne miejsce dla warstwy agentowej i funkcji globalnych (rozdz. 12 i 13) działających *ponad* strukturą, z dostępem do wszystkich środowisk i modułów jednocześnie — co byłoby pojęciowo niespójne w modelu jednego, płaskiego okna czatu. |

Poniższy schemat przedstawia tę hierarchię w formie tekstowej.

```
ZADANIE UŻYTKOWNIKA
        │  (dobór właściwej przestrzeni roboczej)
        ▼
ŚRODOWISKO        „w jakim trybie pracuję?”     TalkIn · WorkSpace · CodeStudio · MultitaskingAI
        │
        ▼
MODUŁ             „jakie zadanie wykonuję?”      Studio · Research · Developer · … (15 modułów)
        │
        ▼
OKNO OPERACYJNE   właściwa przestrzeń wykonania  Studio Editor · Code Editor · Report Builder · …

CZAT — uniwersalna warstwa komunikacji, rekonfigurowana na każdym poziomie hierarchii
```

Cztery powyższe przesłanki stanowią łączne uzasadnienie architektury opisanej w pozostałej części dokumentu.

---

## 1. Charakter platformy

Danaco Console jest systemem operacyjnym dla pracy z modelami AI (AI Workspace OS), a nie aplikacją czatową rozszerzoną o funkcje. Określenie „system operacyjny” jest opisem funkcjonalnym: platforma zarządza wieloma równoległymi przestrzeniami roboczymi, każdą z własnym stanem, kontekstem i zestawem narzędzi — analogicznie do systemu operacyjnego zarządzającego wieloma aplikacjami i procesami. Mechanizm tej równoległości (karty sesji w oknie środowiska) opisano w rozdziale 8.

Fundament platformy stanowi pięć elementów tworzących spójny łańcuch organizacji pracy.

| Element | Rola w łańcuchu organizacji pracy |
|---|---|
| Środowiska | Definiują sposób pracy — wyznaczają ogólny kontekst, w jakim odbywa się zadanie. |
| Moduły | Definiują obszary funkcjonalne — uszczegóławiają kontekst środowiska do konkretnego zadania biznesowego. |
| Okna operacyjne | Właściwa przestrzeń wykonywania zadań — dostarczają narzędzi do jego realizacji. |
| Czat | Uniwersalna warstwa komunikacji — pośredniczy w rozmowie z AI na każdym poziomie łańcucha. |
| Agenci | Warstwa inteligencji działająca ponad całym ekosystemem — obserwuje i działa niezależnie od położenia użytkownika w strukturze. |

Wejście do łańcucha prowadzi przez stronę główną (rozdz. 7), z której użytkownik wybiera środowisko. Każde przełączenie środowiska lub modułu tworzy nową wyspecjalizowaną przestrzeń roboczą z własnym układem, historią, kontekstem i zestawem narzędzi, zachowując spójny sposób komunikacji z AI. Zmiana przestrzeni jest zawsze świadomym, widocznym przejściem, a nie ukrytą zmianą stanu wewnątrz jednego okna.

---

## 2. Pojęcia podstawowe

Cztery pojęcia — Środowisko, Moduł, Komponent własny i Czat — stanowią słownik bazowy całej koncepcji; ich precyzyjne rozumienie jest warunkiem poprawnej interpretacji rozdziałów 3–15.

| Pojęcie | Czym jest | Charakter | Pytanie, na które odpowiada |
|---|---|---|---|
| Środowisko | Najwyższy poziom organizacji pracy w platformie | Element platformy | „W jakim trybie pracuję?” |
| Moduł | Wyspecjalizowany obszar roboczy działający wewnątrz środowiska | Element platformy | „Jakie konkretne zadanie wykonuję?” |
| Komponent własny | Nazwany wytwór użytkownika (automatyka, agent, projekt, profil asystenta) | Wytwór użytkownika | „Co sam buduję i wpinam do pracy?” |
| Czat | Uniwersalny mechanizm komunikacji z AI | Element platformy (nie stały komponent) | „Jak komunikuję się z AI?” |

```
ŚRODOWISKO
   └── MODUŁ  (wiele modułów; dostępność wg macierzy — rozdz. 10)
          └── OKNA OPERACYJNE

KOMPONENT WŁASNY — wytwarzany przez użytkownika, wpinany do modułu, w którym ma zastosowanie
CZAT             — obecny na każdym poziomie, rekonfigurowany przy każdej zmianie kontekstu
```

### 2.1. Środowisko

Środowisko jest najwyższym poziomem organizacji pracy. Definiuje: układ interfejsu, dostępne moduły, typ pracy, sposób działania AI, dostępne narzędzia, model nawigacji, historię sesji, kontekst użytkownika.

Odpowiada na pytanie „w jakim trybie pracuję?”, a nie „jakie zadanie wykonuję?” — to rozróżnia je od modułu. Wybór środowiska jest decyzją najbardziej ogólną: określa, czy użytkownik pracuje w trybie wiedzy i komunikacji, produktywności projektowej, programistycznym, czy orkiestracji autonomicznej pracy ciągłej. Przełączenie środowiska zmienia całą przestrzeń roboczą — dostępne narzędzia, model nawigacji oraz historię sesji przypisaną do tego trybu pracy.

### 2.2. Moduł

Moduł jest wyspecjalizowanym obszarem roboczym działającym wewnątrz środowiska. Definiuje: konkretne zadanie biznesowe, zestaw okien operacyjnych, dedykowane narzędzia, dedykowany przebieg pracy, typ danych wykorzystywanych przez AI.

Moduł doprecyzowuje ogólny tryb środowiska do konkretnego zadania — w środowisku TalkIn moduł Studio odpowiada za pracę z dokumentami, a Research za prowadzenie badań. Przełączenie modułu przeładowuje przestrzeń: znikają okna poprzedniego modułu, a pojawia się zestaw okien dopasowany do nowego kontekstu. W obrębie jednej karty sesji widoczny jest interfejs jednego modułu; praca równoległa odbywa się przez wiele kart sesji otwartych w oknie środowiska (rozdz. 8), a zakres współdzielenia kontekstu użytkownik ustala w oknie konfiguracji (rozdz. 6).

Moduł nie jest tożsamy z komponentem własnym (2.3): moduł jest elementem platformy dostępnym dla wszystkich, komponent własny jest nazwanym wytworem konkretnego użytkownika.

### 2.3. Komponent własny

Komponent własny jest nazwanym wytworem użytkownika, tworzonym i zapisywanym w oknie konfiguracji właściwym danemu rodzajowi komponentu — w odróżnieniu od modułu (2.2), będącego elementem platformy. Cztery rodzaje komponentów własnych zestawia poniższa tabela.

| Rodzaj komponentu własnego | Tworzony w oknie konfiguracji modułu | Przykład zastosowania operacyjnego |
|---|---|---|
| Automatyka | Automations (11.3) | Wpięta w sesji modułu Developer |
| Agent | Agents (11.15) | Wykorzystywany jako wykonawca zadań w dowolnym środowisku |
| Projekt | Workspace (11.2) | Izolowana przestrzeń projektowa z własnym kontekstem |
| Profil asystenta | Assistant (11.10) | Wykorzystywany operacyjnie w środowisku TalkIn |

Po utworzeniu komponent własny trafia do pamięci aplikacji jako zasób użytkownika i pozostaje wybieralny później, w trakcie sesji, w module, w którym ma zastosowanie. Tworzenie komponentów własnych na stronie głównej opisano w rozdziale 7.

### 2.4. Czat

Czat jest wspólnym mechanizmem komunikacji obecnym we wszystkich środowiskach i modułach; nie jest stałym komponentem interfejsu. Przy każdej zmianie środowiska lub modułu jego wygląd, możliwości, historia, dostępne narzędzia i kontekst są rekonfigurowane — użytkownik komunikuje się z AI przez nową instancję czatu działającą w kontekście wybranego środowiska i modułu.

Ta cecha odróżnia go od czatu w klasycznej aplikacji konwersacyjnej: domyślnie nie jest to jedna, niezmienna historia rozmowy, lecz warstwa dopasowująca zawartość i możliwości do bieżącego modułu. Zakres współdzielenia historii i pamięci między modułami, sesjami lub środowiskami użytkownik ustala w oknie konfiguracji (rozdz. 6) — ciągła, współdzielona historia jest jednym z ustawień konfiguracyjnych, a nie stanem domyślnym. Dzięki temu czat w module Developer „wie” o kodzie i repozytorium, a czat w module Research „wie” o zebranych źródłach, mimo że w obu przypadkach jest to ten sam mechanizm. Realizacją interfejsową czatu jest Chat Window — główne okno komunikacji kanału Użytkownik ↔ Wykonawca, opisane wraz z drugim kanałem operacyjnym w rozdziale 3.

---

## 3. Komunikacja operacyjna

Praca platformy prowadzona jest dwoma kanałami komunikacji operacyjnej. Oba kanały stanowią elementy pierwszoplanowe architektury Danaco Console — nie funkcje dodatkowe pojedynczych modułów — i występują w każdym środowisku, module, oknie operacyjnym oraz przepływie pracy.

| Kanał | Nazwa okna | Uczestnicy | Rola w architekturze | Położenie w układzie |
|---|---|---|---|---|
| Pierwszy | Chat Window | Użytkownik ↔ Wykonawca | Centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami platformy | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Drugi | Execution Loop Window | Koordynator ↔ Wykonawca | Mechanizm pętli wykonawczej, koordynacji zadań, nadzoru, orkiestracji i kontroli realizacji procesów | Kolumna sąsiadująca z Chat Window, otwierana |

### 3.1. Kanał pierwszy — Chat Window (Użytkownik ↔ Wykonawca)

Chat Window jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agentem lub systemem wykonawczym). Stanowi centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami realizowanymi przez platformę.

| Funkcja okna | Zakres |
|---|---|
| Przyjmowanie poleceń | Polecenia formułowane w języku naturalnym, kierowane do Wykonawcy |
| Strumień odpowiedzi | Prezentacja odpowiedzi, wyników cząstkowych i artefaktów wytworzonych w toku pracy |
| Zatwierdzanie i przerywanie | Akceptacja proponowanych działań oraz przerwanie działania w toku |
| Wyjaśnianie | Wyjaśnienie wyniku, przyjętych założeń i kontekstu, w jakim zadanie zostało zrealizowane |

Każdy moduł i każde środowisko udostępnia Chat Window w tym samym miejscu układu — w lewej kolumnie obszaru roboczego, o stałej szerokości regulowanej przez użytkownika i pełnej wysokości obszaru roboczego. Chat Window należy do warstwy widoczności 1 (rozdz. 4): jest widoczne bez interakcji, niezależnie od aktywnego modułu i wykonywanej czynności.

### 3.2. Kanał drugi — Execution Loop Window (Koordynator ↔ Wykonawca)

Execution Loop Window jest oknem pętli wykonawczej prezentującym komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów.

| Zawartość okna | Zakres |
|---|---|
| Zlecenie i dekompozycja | Bieżące zlecenie wraz z jego rozłożeniem na zadania składowe |
| Kolejka i stan zadań | Kolejka zadań oraz stan realizacji każdego z nich |
| Komunikaty sterujące | Wymiana komunikatów sterujących między Koordynatorem a Wykonawcą |
| Kontrola jakości | Wyniki kontroli jakości oraz decyzje o ponowieniu zadania |
| Wskaźniki przebiegu | Wskaźniki przebiegu pętli wykonawczej |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie oraz korekta zlecenia |

Okno otwierane jest jako kolumna sąsiadująca z Chat Window i pozostaje otwarte przez cały czas trwania pętli wykonawczej. Wyzwalacz otwarcia okna oraz wskaźnik stanu pętli należą do warstwy widoczności 1, zawartość szczegółowa — do warstwy 2 (rozdz. 4).

### 3.3. Role uczestników komunikacji

| Rola | Charakter | Odpowiedzialność |
|---|---|---|
| Użytkownik | Osoba prowadząca pracę | Zleca zadania, zatwierdza wyniki, przerywa i koryguje działania |
| Koordynator | Komponent orkiestrujący platformy | Dekomponuje zlecenie na zadania, przydziela je, nadzoruje przebieg pętli wykonawczej i kontroluje realizację |
| Wykonawca | AI, agent lub system wykonawczy | Realizuje zadania przydzielone przez Użytkownika i Koordynatora oraz zwraca wyniki |

Rozdzielenie ról jest podstawą obu kanałów: kanał pierwszy łączy Użytkownika z Wykonawcą, kanał drugi — Koordynatora z Wykonawcą. Wcielenia ról w środowisku MultitaskingAI (Executor, Coordinator, Validator) opisano w rozdziale 13.3.

### 3.4. Rozmieszczenie kanałów w obszarze roboczym

Oba kanały zajmują stałe miejsce w układzie pionowym obszaru roboczego (podział lewa–prawa). Kolumny sąsiadują poziomo; regulacji podlega wyłącznie ich szerokość.

```
 ════════════════════════════════════════════════════════════════════════════
  Boczna    │ Chat Window          │ Obszar roboczy modułu    │ Panel
  nawigacja │ Użytkownik ↔         │                          │ pomocniczy
  modułów   │ Wykonawca            │  (okna edycyjne,         │ (rozszerzenie
            │                      │   podglądu, monitory)    │  boczne)
            │ ─────────────────    │                          │
            │ Execution Loop       │                          │
            │ Koordynator ↔        │                          │
            │ Wykonawca            │                          │
 ════════════════════════════════════════════════════════════════════════════
```

Okna pomocnicze i panele otwierane są jako kolejne kolumny boczne po prawej stronie obszaru roboczego, jako rozszerzenia boczne. Żaden element interfejsu nie jest rozmieszczany w podziale poziomym obszaru roboczego.

---

## 4. Warstwy widoczności interfejsu

Zasadą nadrzędną interfejsu Danaco Console jest reguła: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Interfejs ujawnia możliwości systemu stopniowo — zależnie od kontekstu, roli użytkownika i wykonywanej czynności — zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli, ustawień i funkcji administracyjnych nie wpływa na postrzeganą prostotę interfejsu.

### 4.1. Cztery warstwy widoczności

Każdy element interfejsu należy do dokładnie jednej z czterech warstw widoczności. Przy opisie okna, panelu, paska narzędzi i pojedynczej funkcji podaje się jej warstwę.

| Warstwa | Nazwa | Zawartość | Sposób dostępu |
|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, lista rozwijana |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów |

### 4.2. Mechanizmy ukrywania funkcjonalności

| Mechanizm | Działanie |
|---|---|
| Menu progresywne | Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`); lista pozycji rozwija się po kliknięciu |
| Panele wysuwane | Funkcjonalność umieszczana jest w panelach bocznych, panelach wysuwanych, oknach popover i panelach kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej |
| Grupowanie logiczne akcji | Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`), którego rozwinięcie zawiera pełną listę akcji |
| Znaczniki kontekstowe | Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu w postaci `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera selektor wartości tego znacznika |

### 4.3. Zasada jednego kliknięcia

Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu.

### 4.4. Stan spoczynku interfejsu

W stanie spoczynku interfejsu widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki kontekstowe, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisywane są w tabelach elementów okna z podaniem warstwy i sposobu wywołania. Ten sam stan przedstawiają makiety interfejsu zamieszczone w dokumentacji.

---

## 5. Warstwy architektury

Architektura platformy dzieli się na trzy warstwy, z których każda odpowiada za inny aspekt funkcjonowania systemu i porządkuje pojęcia z rozdziału 2.

| Warstwa | Co definiuje | Elementy | Podlega przełączaniu i przeładowaniu |
|---|---|---|---|
| Warstwa środowisk (Environment Layer) | Sposób pracy użytkownika | TalkIn, WorkSpace, CodeStudio, MultitaskingAI | Tak |
| Warstwa modułów (Workspace Layer) | Wyspecjalizowane przestrzenie robocze (własna powłoka, układ okien, historia sesji, kontekst, przebieg pracy) | Piętnaście modułów (rozdz. 11) | Tak |
| Warstwa funkcji globalnych (Global Features Layer) | Funkcje działające ponad wszystkimi środowiskami, modułami i projektami | Mobile, Always On Display | Nie |

```
╔══════════════════════════════════════════════════════════════════════════╗
║  WARSTWA FUNKCJI GLOBALNYCH  ·  Global Features Layer                    ║
║  Mobile · Always On Display                                              ║
║  — ponad całą strukturą; bez własnej powłoki; bez przeładowania          ║
╠══════════════════════════════════════════════════════════════════════════╣
║  WARSTWA MODUŁÓW  ·  Workspace Layer                                     ║
║  Studio · Workspace · Automations · … · Agents  (15 modułów)             ║
║  — własna powłoka, układ okien, historia sesji, kontekst, przebieg pracy ║
╠══════════════════════════════════════════════════════════════════════════╣
║  WARSTWA ŚRODOWISK  ·  Environment Layer   (warstwa najwyższego rzędu)   ║
║  TalkIn · WorkSpace · CodeStudio · MultitaskingAI                        ║
║  — definiuje sposób pracy użytkownika                                    ║
╚══════════════════════════════════════════════════════════════════════════╝
    MultitaskingAI: zamiast bocznej nawigacji modułów — panel orkiestracji (rozdz. 13.8)
```

### 5.1. Warstwa środowisk (Environment Layer)

Definiuje sposób pracy użytkownika; środowiska: TalkIn, WorkSpace, CodeStudio, MultitaskingAI. Jest warstwą najwyższego rzędu — ustala ogólne ramy pracy. Każde z czterech środowisk (rozdz. 9) stanowi odrębny kontekst nawigacyjny i funkcjonalny, wewnątrz którego działają moduły warstwy niższej — z wyjątkiem MultitaskingAI, które nie ma bocznej nawigacji modułów (rozdz. 9.4, 10).

### 5.2. Warstwa modułów (Workspace Layer)

Definiuje wyspecjalizowane przestrzenie robocze z własną powłoką, układem okien, historią sesji, kontekstem oraz przebiegiem pracy.

**Uwaga terminologiczna.** Nazwa warstwy — Workspace Layer — oraz nazwa modułu Workspace (11.2) są bytami odrębnymi, celowo zbieżnymi jedynie brzmieniowo. Warstwa modułów jest pojęciem architektonicznym obejmującym wszystkie moduły platformy (Studio, Automations, Research, Developer i pozostałe z rozdz. 11). Moduł Workspace (11.2) jest jednym konkretnym modułem tej warstwy, odpowiedzialnym za izolowane środowisko projektowe. W dalszej części dokumentu nazwa „Workspace” bez kwalifikatora odnosi się zawsze do modułu z rozdziału 11.2, a odniesienia do warstwy są opatrzone pełną nazwą „warstwa modułów” lub „Workspace Layer”.

### 5.3. Warstwa funkcji globalnych (Global Features Layer)

Funkcje globalne (Mobile, Always On Display) nie posiadają własnej powłoki roboczej, nie przeładowują środowiska i nie tworzą nowych sesji — działają ponad wszystkimi środowiskami, modułami i projektami. Warstwa ta jest architektonicznie odrębna, bo nie podlega mechanizmowi przełączania i przeładowania właściwemu środowiskom i modułom: pozostaje dostępna w niezmienionej formie niezależnie od położenia użytkownika. Ten brak zależności od kontekstu czyni ją „globalną”.

### 5.4. Schemat hierarchii warstw

Poniższy schemat przedstawia trójstopniową hierarchię (środowisko → moduł → okno operacyjne) wraz z warstwą funkcji globalnych i komponentów działających ponad tą strukturą; stanowi punkt odniesienia dla dalszej lektury.

```
STRONA GŁÓWNA — Centrum dowodzenia (rozdz. 7)
│
├── Strefa 1 · wybór środowiska
│       TalkIn · WorkSpace · CodeStudio · MultitaskingAI
│
├── Strefa 2 · komponenty własne
│       Automations · Agents · Workspace · Assistant
│
└── Strefa 3 · ustawienia
        Okno konfiguracji · Mobile · Always On Display

         │  (wybór środowiska w strefie 1)
         ▼
ŚRODOWISKO  (TalkIn / WorkSpace / CodeStudio / MultitaskingAI)
│
├── boczna nawigacja modułów        → TalkIn, WorkSpace, CodeStudio
│   albo panel orkiestracji          → MultitaskingAI (rozdz. 13.8)
│
├── MODUŁ  (Studio, Research, Developer … — rozdz. 11)
│       └── OKNA OPERACYJNE  (Chat Window, Execution Loop Window, Studio Editor … — Załącznik A)
│
└── KARTY SESJI  (równoległe — rozdz. 8)

PONAD CAŁĄ STRUKTURĄ
    Warstwa funkcji globalnych:  Mobile · Always On Display  (rozdz. 12)
    Agenci (komponent własny modułu Agents) — dostępni we wszystkich środowiskach
```

---

## 6. Konfigurowalność zależności i punktów izolacji

Platforma nie narzuca twardej izolacji historii sesyjnej ani pamięci i nie ma blokad wbudowanych na stałe. Wszystkie zależności platformy — między środowiskami, modułami, komponentami własnymi, sesjami i projektami — są konfigurowalne z okna konfiguracji. Zasada ma charakter fundamentalny: określa, jak rozumieć każdy opisany w dokumencie mechanizm domyślnego zachowania, izolacji lub powiązania.

### 6.1. Punkty izolacji jako decyzja użytkownika

To użytkownik ustanawia punkty izolacji, a nie platforma. Środowiska, moduły i komponenty własne dają domyślny, uporządkowany podział pracy — z osobnym układem okien, historią i pamięcią dla każdej nowej karty sesji (rozdz. 8) — ale jest to punkt wyjścia, a nie sztywna reguła. Z okna konfiguracji użytkownik podejmuje cztery rodzaje rozstrzygnięć; pełny wykaz pozycji podlegających konfiguracji zawiera [Izolacja i zależności](izolacja-i-zaleznosci.md):

- zdecydować, czy i w jakim zakresie historia oraz pamięć są współdzielone między kartami sesji, modułami lub środowiskami,
- połączyć moduł jako funkcję z innym modułem — współdzielić mechanizm operacji kontekstowych AI między Studio i Translate,
- powiązać komponent własny z innym komponentem lub z modułem — spiąć automatykę ze środowiskiem MultitaskingAI (decyzja podejmowana w oknie konfiguracji oraz w ustawieniach okna modułu Automations),
- utrzymać pełne rozdzielenie kontekstu tam, gdzie jest pożądane — między równolegle prowadzonymi projektami w module Workspace.

### 6.2. Zależności jako możliwość, nie wbudowana reguła

Podpunkty „Powiązania” opisane przy każdym module w rozdziale 11 wskazują możliwe powiązania między modułami i funkcjami — nie są to zależności wbudowane na stałe, lecz możliwości, które użytkownik może, ale nie musi, ustanowić. Ten sam mechanizm dotyczy relacji między środowiskiem MultitaskingAI a modułem Automations (rozdz. 13.6) oraz relacji między silnikami kolejek z rozdziałów 11 (Automations) i 13 (MultitaskingAI).

### 6.3. Konsekwencje dla porządku pracy

Domyślny podział na środowiska, moduły i karty sesji jest wystarczający dla porządku pracy i nie wymaga dodatkowej konfiguracji. Konfigurowalność zależności i punktów izolacji jest możliwością dostępną wtedy, gdy użytkownik świadomie chce połączyć konteksty, a nie wymaganiem. Opisane w dalszej części dokumentu domyślne zachowanie modułu, funkcji lub środowiska należy rozumieć jako stan wyjściowy, możliwy do zmiany zgodnie z niniejszym rozdziałem oraz zasadą pełnej konfigurowalności (rozdz. 14, zasada 2).

### 6.4. Okno konfiguracji punktów izolacji

Rozdział 6 ustala zasadę pełnej konfigurowalności; niniejszy podpunkt określa interfejs jej przeprowadzania. Okno konfiguracji punktów izolacji jest jednym miejscem, w którym schodzą się cztery rodzaje izolacji dotąd opisane w odrębnych źródłach.

| Rodzaj izolacji | Czego dotyczy | Liczba pozycji | Podstawa |
|---|---|---|---|
| Izolacja kontekstu | Historia, pamięć i kontekst współdzielone między kartami sesji, modułami, środowiskami i projektami | 3 | Niniejszy rozdział oraz rozdz. 2.4 i 8 |
| Izolacja techniczna procesu sesji | Osiem zakresów z rozdziału 12 dokumentu Architektury: katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu (odrębny lub współdzielony), serwer wykonania | 8 | Rozdz. 12 dokumentu Architektury |
| Izolacja warstwy komunikacji operacyjnej | Historia i pamięć okna Chat Window oraz zlecenia i zadania, kolejka i rejestr przebiegu pętli okna Execution Loop Window (rozdz. 3.1–3.2) | 5 | [Izolacja i zależności](izolacja-i-zaleznosci.md) rozdz. 14 |
| Izolacja warstw widoczności | Konfiguracja warstw widoczności według roli użytkownika, rejestr funkcji wyszukiwarki funkcji oraz zasięg przypisania elementu interfejsu do warstwy (rozdz. 4.1) | 3 | [Izolacja i zależności](izolacja-i-zaleznosci.md) rozdz. 15 |

Połączenie czterech rodzajów izolacji w jednym oknie odpowiada zasadzie konfiguracji wszystkich zależności z jednego miejsca. Polityka izolacji jednego poziomu zasięgu jest zatem kombinacją dziewiętnastu niezależnych ustawień. Okno dzieli się na trzy panele.

| Panel | Położenie | Zawartość |
|---|---|---|
| Selektor zasięgu | Lewy | Wybór poziomu obowiązywania reguły izolacji (rozdz. 6.5): globalny, środowisko, moduł, para modułów, projekt, karta sesji lub rola środowiska MultitaskingAI |
| Macierz izolacji | Środkowy | Cztery grupy przełączników: dla izolacji kontekstu każda z trzech pozycji (historia, pamięć, kontekst) przyjmuje stan „współdzielone” albo „odrębne”; dla izolacji technicznej każdy z ośmiu zakresów przyjmuje stan „włączony” albo „wyłączony”; dla izolacji warstwy komunikacji operacyjnej każda z pięciu pozycji oraz dla izolacji warstw widoczności każda z trzech pozycji przyjmuje stan „współdzielone” albo „odrębne” |
| Profil i podgląd | Prawy | Zapis, wczytanie i przypisanie profilu izolacji (rozdz. 6.6), wybór warstwy konfiguracji (domyślna albo sesji) oraz podgląd polityki efektywnej — wynikowego zestawu reguł po uwzględnieniu dziedziczenia między poziomami zasięgu |

Poniższy schemat przedstawia układ okna w formie tekstowej.

```
┌────────────────┬────────────────────────────────────┬─────────────────────┐
│ SELEKTOR       │ MACIERZ IZOLACJI                   │ PROFIL I PODGLĄD    │
│ ZASIĘGU        │                                    │                     │
│                │ Izolacja kontekstu:                │ [ Zapisz profil ]   │
│ • Globalny     │   historia         [współdz.|odr.] │ [ Wczytaj profil ]  │
│ • Środowisko   │   pamięć           [współdz.|odr.] │ [ Przypisz do… ]    │
│ • Moduł        │   kontekst         [współdz.|odr.] │                     │
│ • Para modułów │                                    │ Warstwa:            │
│ • Projekt      │ Izolacja techniczna:               │  ( ) domyślna       │
│ • Karta sesji  │   katalog roboczy         [wł|wył] │  ( ) sesji          │
│ • Rola         │   środowisko proc.        [wł|wył] │                     │
│                │   dane/konfig. mod.       [wł|wył] │ Polityka efektywna: │
│                │   dostęp sieciowy         [wł|wył] │  (podgląd reguł     │
│                │   odczyt/zapis pl.        [wł|wył] │   dziedziczonych)   │
│                │   konto i token           [wł|wył] │                     │
│                │   model procesu           [wł|wył] │ [?] objaśnienia     │
│                │   serwer wykonania        [wł|wył] │     kontekstowe     │
│                │                                    │                     │
│                │ Izolacja komunikacji operacyjnej:  │                     │
│                │   historia Chat W. [współdz.|odr.] │                     │
│                │   pamięć Chat W.   [współdz.|odr.] │                     │
│                │   zlecenia/zadania [współdz.|odr.] │                     │
│                │   kolejka pętli    [współdz.|odr.] │                     │
│                │   rejestr przebiegu[współdz.|odr.] │                     │
│                │                                    │                     │
│                │ Izolacja warstw widoczności:       │                     │
│                │   warstwy wg roli  [współdz.|odr.] │                     │
│                │   rejestr funkcji  [współdz.|odr.] │                     │
│                │   przypisanie el.  [współdz.|odr.] │                     │
└────────────────┴────────────────────────────────────┴─────────────────────┘
```

Legenda: `proc.` — procesu; `mod.` — modelu; `pl.` — plików; `Chat W.` — Chat Window;
`el.` — elementu interfejsu; `współdz.` — współdzielone; `odr.` — odrębne; `wł` — włączony;
`wył` — wyłączony.

Każdy przełącznik macierzy opatrzony jest objaśnieniem kontekstowym (znak zapytania) opisującym działanie ustawienia i jego wpływ na aplikację — zgodnie z zasadą przyjętą dla okna konfiguracji w rozdziale 13 dokumentu Architektury. Zgodnie z zasadą pełnej konfigurowalności (rozdz. 14, zasada 2) oraz nadrzędną zasadą braku twardych blokad okno nigdy nie wymusza izolacji — udostępnia ją jako możliwość, a jego stan wyjściowy opisano w rozdziale 6.6.

### 6.5. Poziomy zasięgu reguł izolacji

Izolację ustawia się jednocześnie na wszystkich poziomach — od globalnego po parę modułów, projekt, kartę sesji, rolę i okno komunikacji — jako hierarchię zasięgów wybieraną w panelu selektora (rozdz. 6.4). Ograniczenie konfiguracji do jednego narzuconego poziomu stanowiłoby twardą regułę sprzeczną z zasadą pełnej konfigurowalności. Platforma udostępnia osiem poziomów, od najogólniejszego do najbardziej szczegółowego.

| Poziom zasięgu | Przykład zastosowania | Pierwszeństwo |
|---|---|---|
| Globalny (domyślny) | Domyślna polityka izolacji obowiązująca w całej platformie. | Najniższe — warstwa bazowa |
| Środowisko | Odmienna polityka dla środowiska CodeStudio niż dla TalkIn. | ↑ |
| Moduł | Odrębna pamięć przypisana modułowi Developer. | ↑ |
| Para modułów (relacja) | Współdzielenie operacji kontekstowych AI między Studio a Translate. | ↑ |
| Projekt | Rozdzielenie kontekstu między dwoma projektami w module Workspace. | ↑ |
| Karta sesji | Jednorazowe współdzielenie historii między dwiema otwartymi kartami. | ↑ |
| Rola (MultitaskingAI) | Profil izolacji przypisany roli Executor 1. | ↑ |
| Okno komunikacji | Odrębna pamięć i odrębne uprawnienia okna Execution Loop Window wobec okna Chat Window tej samej karty sesji. | Najwyższe — poziom najwęższy |

Reguły z różnych poziomów rozstrzyga zasada pierwszeństwa zasięgu najbardziej szczegółowego: ustawienie z poziomu węższego ma pierwszeństwo przed ustawieniem z poziomu szerszego (pierwszeństwo rośnie ku dołowi tabeli). Brak ustawienia na danym poziomie oznacza dziedziczenie wartości z poziomu bezpośrednio szerszego, aż do poziomu globalnego domyślnego. Jest to zastosowanie zasady nadrzędnej „brak ustawienia = wartość domyślna”: żaden poziom nie musi być skonfigurowany, a poziom nieskonfigurowany nie blokuje pracy, lecz przejmuje wartość odziedziczoną.

### 6.6. Profile izolacji i warstwy konfiguracji

Zestaw ustawień izolacji można zapisać jako nazwany profil izolacji i przypisać go do sesji lub roli — zgodnie z mechanizmem profilu z rozdziału 12 dokumentu Architektury. Profil grupuje wybrane wartości izolacji kontekstu i izolacji technicznej dla wskazanego poziomu zasięgu, dzięki czemu raz zbudowana polityka może być stosowana wielokrotnie. Szablon profilu izolacji przedstawia Załącznik D.

Konfiguracja izolacji ma strukturę warstwową, zgodną z rozdziałem 13 dokumentu Architektury.

| Warstwa konfiguracji | Zakres obowiązywania | Relacja do pozostałych |
|---|---|---|
| Warstwa domyślna | Ustalana w oknie konfiguracji; obowiązuje przy każdej nowej sesji | Warstwa bazowa |
| Warstwa sesji | Zmiany dokonane dla sesji bieżącej | Nakłada się na warstwę domyślną, nie modyfikując jej wartości |

Dzięki warstwowości jednorazowe poluzowanie lub zaostrzenie izolacji w konkretnej sesji nie zmienia domyślnego zachowania platformy. Stanem wyjściowym — obowiązującym, dopóki użytkownik nie utworzy profilu ani nie zmieni ustawienia — jest stan z rozdziału 6.1: każda nowa karta sesji otrzymuje odrębną historię, pamięć i kontekst, przy jednoczesnym braku aktywnej izolacji technicznej (żaden z zakresów z rozdziału 12 dokumentu Architektury nie jest domyślnie włączony, więc proces sesji nie jest technicznie ograniczany). Taki stan łączy przejrzysty porządek kontekstu z pełną swobodą procesu; wszelka dalsza izolacja — zarówno współdzielenie kontekstu, jak i włączenie zakresów izolacji technicznej — jest świadomą decyzją użytkownika.

---

## 7. Strona główna i Centrum dowodzenia

Główne okno aplikacji nie otwiera od razu przestrzeni roboczej żadnego środowiska. Pełni funkcję przedpokoju przed strefą roboczą — centrum dowodzenia, z którego użytkownik wybiera środowisko, konfiguruje własne komponenty lub przechodzi do ustawień. Dopiero wybór środowiska w strefie 1 otwiera właściwą przestrzeń roboczą (rozdz. 9). Strona główna dzieli się na trzy strefy.

```
STRONA GŁÓWNA — Centrum dowodzenia
├─ Strefa 1 · wybór środowiska    → Karty środowisk        (waga główna)
│     TalkIn · WorkSpace · CodeStudio · MultitaskingAI
├─ Strefa 2 · komponenty własne   → Kafle komponentów      (waga pośrednia)
│     Automations · Agents · Workspace · Assistant
└─ Strefa 3 · ustawienia          → Listwa ustawień        (waga najniższa)
      Okno konfiguracji · Mobile · Always On Display
```

### 7.1. Strefa 1 — wybór środowiska

Punkt wejścia do czterech środowisk: TalkIn, WorkSpace, CodeStudio i MultitaskingAI (rozdz. 9). Wybór środowiska jest pierwszą decyzją nawigacyjną użytkownika, opisaną szerzej w rozdziale 8.

### 7.2. Strefa 2 — komponenty własne

Punkt wejścia do konfiguracji czterech komponentów własnych: Automations, Agents, Workspace i Assistant. Wejście w daną pozycję otwiera okno konfiguracji, w którym użytkownik tworzy i zapisuje nazwany komponent własny (2.3) — konkretną automatykę, agenta, projekt lub profil asystenta.

Utworzony komponent trafia do pamięci aplikacji jako zasób własny i pozostaje wybieralny później, w trakcie sesji, w module, w którym ma zastosowanie — automatykę wpina się w sesji modułu Developer. Zakres okna modułowego dostępnego dodatkowo wewnątrz środowisk dla tych komponentów opisano w rozdziale 10.

### 7.3. Strefa 3 — ustawienia

Okno konfiguracji i ustawień platformy oraz dwie funkcje globalne z rozdziału 12: Mobile i Always On Display.

### 7.4. Forma prezentacji

Forma prezentacji trzech stref wynika z ich odmiennej funkcji i wagi. Centrum dowodzenia nie jest jednorodną siatką równoważnych elementów, lecz przestrzenią o wyraźnej hierarchii: jedna strefa jest głównym punktem wejścia do pracy, druga — zapleczem twórczym budowy komponentów własnych, trzecia — warstwą narzędziową ustawień. Trzy strefy otrzymują zatem trzy różne formy prezentacji.

| Strefa | Zawartość | Forma prezentacji | Waga wizualna |
|---|---|---|---|
| Strefa 1 — wybór środowiska | TalkIn, WorkSpace, CodeStudio, MultitaskingAI | Karty środowisk — cztery duże karty wejścia | Główna (środek ciężkości) |
| Strefa 2 — komponenty własne | Automations, Agents, Workspace, Assistant | Kafle komponentów własnych — siatka czterech kafli | Pośrednia |
| Strefa 3 — ustawienia | Okno konfiguracji, Mobile, Always On Display | Listwa ustawień — zwarty pasek narzędziowy | Najniższa |

- **Strefa 1 (wybór środowiska)** — cztery duże karty środowisk, po jednej na środowisko, z tytułem złożonym krojem nagłówkowym, jednozdaniowym opisem trybu pracy i godłem środowiska. Karty są wizualnym środkiem ciężkości strony; akcent złoty rezerwowany dla stanu aktywnego i najechania, zgodnie z zasadą oszczędnego użycia złota w systemie wizualnym marki. Afordancja: „wejdź do przestrzeni roboczej”.
- **Strefa 2 (komponenty własne)** — cztery jednorodne kafle w siatce, wizualnie lżejsze od kart strefy 1, z etykietą zorientowaną na działanie twórcze (utworzenie automatyki, agenta, projektu lub profilu asystenta). Odmienna forma oddziela „wejście do środowiska” od „zbudowania komponentu”.
- **Strefa 3 (ustawienia)** — zwarty pasek z wejściem do okna konfiguracji oraz dwiema funkcjami globalnymi (Mobile, Always On Display). Sygnalizuje warstwę narzędziową, stale dostępną, o najniższej wadze wizualnej.

Wybór opiera się na czterech przesłankach: (1) metaforze centrum dowodzenia, które ma hierarchię, a nie jednorodną siatkę; (2) zasadach systemu wizualnego marki — oszczędnym użyciu złota jako akcentu oraz kroju nagłówkowym dla tytułów kart; (3) zgodności formy z rodzajem akcji (wejście, utworzenie, ustawienie); (4) czytelnym rytmie strony — dwie siatki po cztery elementy o różnej skali, dopełnione listwą. Szczegółową realizację wizualną (wymiary, odstępy, tokeny kolorów i typografii) określa system projektowy marki (granat i złoto).

---

## 8. Nawigacja i sesje

Nawigacja w platformie przebiega w czterech krokach.

| Krok | Przejście | Opis |
|---|---|---|
| 1 | Strona główna → środowisko | Ze strony głównej (rozdz. 7, strefa 1) użytkownik wchodzi do wybranego środowiska. |
| 2 | Boczna nawigacja → moduł | Wewnątrz środowiska moduły dostępne w danym środowisku (rozdz. 10) wybiera się ze stałej bocznej nawigacji, widocznej niezależnie od otwartego modułu. Wyjątkiem jest MultitaskingAI, które ma nie listę modułów, lecz panel orkiestracji (rozdz. 9.4, 10, 13.8). |
| 3 | Moduł → okno kontekstowe | Wybrany moduł otwiera właściwe mu okno kontekstowe/robocze — dopiero tu zaczyna się praca właściwa, opisana przy każdym module w rozdziale 11. |
| 4 | Karty sesji | Aktywne sesje trzymane są w kartach sesji w oknie środowiska. Karta aktywnej sesji jest kartą środowiska sesyjnego; przełączanie kart oznacza przełączanie między sesjami — mechanika analogiczna do kart w przeglądarce internetowej. |

Ten mechanizm rozstrzyga pozorną sprzeczność między zarządzaniem wieloma równoległymi przestrzeniami (rozdz. 1) a przeładowaniem przestrzeni przy zmianie modułu (rozdz. 2.2): w obrębie jednej karty sesji widoczny jest interfejs jednego modułu, a przeładowanie następuje przy zmianie modułu w tej samej karcie; praca równoległa odbywa się przez otwarcie kilku kart sesji obok siebie. Zakres, w jakim karty sesji współdzielą kontekst, historię i pamięć, jest konfigurowalny (rozdz. 6).

```
Strona główna
    │  wybór w strefie 1
    ▼
Środowisko  ──(boczna nawigacja / panel orkiestracji)──►  Moduł  ──►  Okno kontekstowe / robocze
    │                                                                        │
    │  ◄──────────────────── karty sesji         ────────────────────────────┘
    ▼
[ Karta sesji A · Moduł X ]  [ Karta sesji B · Moduł Y ]  [ + nowa karta ]
        │                            │
   własny układ,               własny układ,
   historia, kontekst          historia, kontekst
   (współdzielenie konfigurowalne — rozdz. 6)
```

---

## 9. Środowiska

Platforma udostępnia cztery środowiska, odpowiadające czterem odrębnym trybom pracy z AI. Wybór środowiska w strefie 1 strony głównej (rozdz. 7) jest pierwszą decyzją nawigacyjną i określa, które moduły są dostępne w bocznej nawigacji (rozdz. 10).

| Środowisko | Tryb pracy | Przeznaczenie | Boczna nawigacja |
|---|---|---|---|
| TalkIn | Wiedza, komunikacja i praca z treścią | Rozmowy z AI, praca z dokumentami, analizy, tłumaczenia, badania, raporty, współpraca wielu modeli | Lista modułów (rozdz. 10) |
| WorkSpace | Produktywność, organizacja i realizacja projektów | Zarządzanie projektami, automatyzacja, prowadzenie procesów, tworzenie aplikacji biznesowych, projektowanie, praca operacyjna | Lista modułów (rozdz. 10) |
| CodeStudio | Programowanie | Tworzenie oprogramowania, praca z kodem, debugowanie, wykorzystanie terminali, budowanie aplikacji, procesy developerskie | Lista modułów (rozdz. 10) |
| MultitaskingAI | Orkiestracja autonomicznej pracy ciągłej | Zespół modeli i agentów w rolach; po spięciu z automatyką — pełna pętla pracy ciągłej 24/7/365 z własnymi akcjami i harmonogramem | Panel orkiestracji (rozdz. 13.8) |

### 9.1. TalkIn

Środowisko wiedzy, komunikacji i pracy z treścią — właściwe dla zadań, w których centralnym przedmiotem pracy jest treść: jej tworzenie, analiza, przekształcanie lub pozyskiwanie ze źródeł zewnętrznych. Pełną listę modułów przedstawia macierz dostępności (rozdz. 10).

### 9.2. WorkSpace

Środowisko produktywności, organizacji i realizacji projektów, zorientowane na przekształcanie zamierzeń w zorganizowane działania i produkty. Istotne są struktura projektu, powtarzalność procesów oraz koordynacja zadań w czasie — w odróżnieniu od jednorazowej, konwersacyjnej pracy z treścią właściwej TalkIn. Pełną listę modułów przedstawia macierz dostępności (rozdz. 10).

### 9.3. CodeStudio

Środowisko programistyczne, dedykowane pracy inżynierskiej — obszarom, których wspólnym mianownikiem jest kod źródłowy oraz środowiska wykonawcze. Pełną listę modułów przedstawia macierz dostępności (rozdz. 10).

### 9.4. MultitaskingAI

Środowisko orkiestracji autonomicznej pracy ciągłej. Poza czterema oknami roboczymi (rolami) daje, po spięciu z automatyką, pełną orkiestrację autonomicznej pętli pracy ciągłej z własnymi akcjami i harmonogramem (praca 24/7/365). Po przypisaniu ról własnym agentom z komponentu Agents (rozdz. 7, strefa 2) powstaje wieloagentowa pętla z kolejkowaniem, rolami, akcjami i hierarchią decyzji.

Jest najbardziej zaawansowane konfiguracyjnie spośród czterech środowisk — organizuje nie pracę pojedynczego użytkownika, lecz zespół modeli i agentów realizujących wspólny proces. W odróżnieniu od pozostałych trzech środowisk nie ma bocznej nawigacji modułów (rozdz. 10); zawartość jego bocznej nawigacji — panel orkiestracji — opisano w rozdziale 13.8. Pełny opis (definicja, warstwa centralna, role wykonawcze, silnik kolejek, orkiestracja, integracja z Automations) znajduje się w rozdziale 13.

---

## 10. Macierz dostępności modułów

Moduł nie należy do jednego środowiska — jest dostępny w wybranych środowiskach. Środowisko jest profilem widoczności modułów w bocznej nawigacji (rozdz. 8), a nie pojemnikiem, do którego moduł należy wyłącznie.

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|---|---|---|
| Studio (11.1) | TAK | TAK | NIE |
| Workspace (11.2) | TAK | TAK | TAK |
| Automations (11.3) | — | — | — |
| Browser (11.4) | TAK | TAK | NIE |
| Research (11.5) | TAK | TAK | NIE |
| Library (11.6) | TAK | TAK | NIE |
| Translate (11.7) | TAK | NIE | NIE |
| Roundtable (11.8) | TAK | TAK | TAK |
| Design (11.9) | NIE | TAK | TAK |
| Assistant (11.10) | TAK | NIE | NIE |
| Terminal (11.11) | NIE | NIE | TAK |
| Developer (11.12) | NIE | NIE | TAK |
| Diagnostics (11.13) | NIE | NIE | TAK |
| Apps (11.14) | NIE | TAK | TAK |
| Agents (11.15) | TAK | TAK | TAK |

Automations nie ma okna modułowego w żadnym środowisku — działa wyłącznie ze strony głównej (rozdz. 7, strefa 2). Do sesji trafiają gotowe automatyki jako komponenty własne (2.3); moduł nie otwiera własnego okna w bocznej nawigacji, stąd „—” we wszystkich trzech kolumnach.

Automations, Agents, Workspace i Assistant są konfigurowane na stronie głównej, w strefie 2 (rozdz. 7). Agents i Workspace mają dodatkowo okno modułowe w środowiskach wskazanych w macierzy. Assistant jest konfigurowany jako komponent własny na stronie głównej, a wykorzystywany operacyjnie w środowisku TalkIn — stąd jedyne „TAK” w jego wierszu. Środowisko MultitaskingAI nie ma bocznej nawigacji modułów, dlatego nie występuje jako kolumna macierzy; zawartość jego bocznej nawigacji — panel orkiestracji — rozstrzyga rozdział 13.8.

---

## 11. Moduły

Poniżej katalog piętnastu modułów platformy. Każdy moduł opisano w formie tabelarycznej (cel, okna operacyjne, funkcjonalności, powiązania) uzupełnionej scenariuszem użycia. Środowiska, w których dany moduł jest dostępny, wskazuje macierz dostępności (rozdz. 10). Zbiorcze zestawienie okien wszystkich modułów zawiera Załącznik A.

| Moduł | Cel (w skrócie) | Dostępność (środowiska) |
|---|---|---|
| Studio (11.1) | Zaawansowana praca z tekstem, dokumentami i treścią | TalkIn, WorkSpace |
| Workspace (11.2) | Izolowane środowisko projektowe | TalkIn, WorkSpace, CodeStudio |
| Automations (11.3) | Budowanie i wykonywanie procesów automatycznych | Strona główna (bez okna modułowego) |
| Browser (11.4) | Współdzielone przeglądanie internetu | TalkIn, WorkSpace |
| Research (11.5) | Realizacja badań, analiz i opracowań | TalkIn, WorkSpace |
| Library (11.6) | Centralne repozytorium wiedzy i plików | TalkIn, WorkSpace |
| Translate (11.7) | Wielojęzyczne tłumaczenia w trybie wielozadaniowym | TalkIn |
| Roundtable (11.8) | Współpraca wielu modeli AI nad wspólnym problemem | TalkIn, WorkSpace, CodeStudio |
| Design (11.9) | Projektowanie i tworzenie zasobów wizualnych | WorkSpace, CodeStudio |
| Assistant (11.10) | Naturalna komunikacja głosowa z AI | TalkIn (konfiguracja na stronie głównej) |
| Terminal (11.11) | Praca z konsolami i środowiskami wykonawczymi | CodeStudio |
| Developer (11.12) | Tworzenie i rozwój kodu | CodeStudio |
| Diagnostics (11.13) | Analiza oraz usuwanie problemów technicznych | CodeStudio |
| Apps (11.14) | Budowa kompletnych produktów cyfrowych | WorkSpace, CodeStudio |
| Agents (11.15) | Tworzenie i zarządzanie własnymi agentami AI | TalkIn, WorkSpace, CodeStudio (konfiguracja na stronie głównej) |

### 11.1. Studio

| Pole | Treść |
|---|---|
| Cel | Zaawansowana praca z tekstem, dokumentami i treścią — tworzenie, redagowanie i przekształcanie materiałów pisanych, od notatek po obszerne dokumenty. Jedno miejsce trwale integrujące edycję treści, wsparcie AI i porównywanie wersji, zamiast rozproszenia między osobny edytor a osobne okno rozmowy. |
| Okna | Chat Window, Execution Loop Window, Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window |
| Funkcjonalności | edycja dokumentów, obsługa PDF, DOCX, TXT, Markdown, tłumaczenia, korekta, analiza treści, przepisywanie, zmiana stylu, streszczenia, rozwijanie treści, Diff, Grep, operacje kontekstowe AI |

**Typowy scenariusz:**
1. Wczytanie dokumentu do Studio Editor.
2. Zlecenie AI korekty lub zmiany stylu wybranego fragmentu przez Chat Window.
3. Porównanie wersji przed i po zmianie w Diff/Grep Panel.
4. Akceptacja wyniku; powrót do wcześniejszych wersji przez Session Repository bez utraty historii zmian.

**Powiązania:**
- Operacje kontekstowe AI są mechanizmem właściwym modułowi Studio; ich wspólne wykorzystanie z modułem Translate (11.7) przy pracy dwujęzycznej jest konfigurowalne przez użytkownika, nie jest zaś współdzielone domyślnie.
- Studio korzysta z Library (11.6) jako repozytorium źródłowego dla dokumentów wejściowych i miejsca docelowego dla wygenerowanych artefaktów.
- Wyniki pracy w Studio mogą, decyzją użytkownika, być dalej przetwarzane w module Research (11.5) przy budowie raportów końcowych.

### 11.2. Workspace

| Pole | Treść |
|---|---|
| Cel | Izolowane środowisko projektowe dla użytkowników prowadzących równolegle wiele odrębnych projektów; pozwala skonfigurować rozdzielenie kontekstu, instrukcji i pamięci między nimi, tak aby — zgodnie z ustawieniami okna konfiguracji — ustalenia jednego projektu nie przenikały do innego, gdy izolacja jest pożądana. Moduł jest jednym z elementów warstwy modułów (Workspace Layer), nie jej odpowiednikiem (rozdz. 5.2). Konfigurowany jako komponent własny na stronie głównej (rozdz. 7, strefa 2), dodatkowo dostępny jako okno modułowe w środowiskach z macierzy (rozdz. 10). |
| Okna | Chat Window, Execution Loop Window, Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager |
| Funkcjonalności | separacja projektów, własne instrukcje, własna pamięć, własna biblioteka, własna konfiguracja modeli, zarządzanie agentami |

**Typowy scenariusz:**
1. Założenie nowego projektu w module Workspace.
2. Zdefiniowanie odrębnego zestawu instrukcji systemowych w Instructions Panel.
3. Zbudowanie dedykowanej pamięci kontekstowej w Context Memory.
4. Przypisanie do projektu konkretnych agentów przez Agent Manager (domyślnie w zakresie tego projektu; wpływ na inne, równolegle prowadzone projekty użytkownik określa w oknie konfiguracji).

**Powiązania:**
- Moduł może korzystać z agentów skonfigurowanych w module Agents (11.15), udostępniając ich — decyzją użytkownika — jako wykonawców zadań w ramach danego projektu.
- Project Library pełni funkcję analogiczną do modułu Library (11.6), w zakresie ograniczonym do jednego projektu.
- Projekty prowadzone w Workspace mogą zostać powiązane, decyzją użytkownika, z procesami automatycznymi tworzonymi w module Automations (11.3).

### 11.3. Automations

| Pole | Treść |
|---|---|
| Cel | Budowanie i wykonywanie procesów automatycznych dla użytkowników przekształcających powtarzalne czynności (cykliczne raporty, regularne przetwarzanie danych, zaplanowane wywołania modeli) w procesy działające bez stałego nadzoru — uruchamiane według harmonogramu lub zdarzenia, a nie w jednorazowej interakcji. Konfigurowany wyłącznie ze strony głównej (rozdz. 7, strefa 2) — bez własnego okna w bocznej nawigacji (rozdz. 10); gotowe automatyki trafiają do sesji jako komponenty własne (2.3). |
| Okna | Chat Window, Execution Loop Window, Workflow Builder, Scheduler, Queue Manager, Orchestrator, Execution Monitor |
| Funkcjonalności | harmonogramy, zadania cykliczne, przebieg pracy, kolejki, automatyczne wywołania modeli, pętle wielomodelowe, orkiestracja |

**Typowy scenariusz:**
1. Zaprojektowanie procesu w Workflow Builder.
2. Ustalenie jego cykliczności w Scheduler.
3. Obserwacja kolejnych uruchomień w Execution Monitor, gdzie widoczny jest status każdego przebiegu oraz ewentualne błędy wymagające interwencji.

**Powiązania:**
- Po skonfigurowaniu integracji opisanej w rozdziale 13.6, Automations może stanowić operacyjne zaplecze harmonogramów i kolejek dla środowiska MultitaskingAI (rozdz. 9.4, 13) — połączenie nie jest domyślne i wymaga decyzji użytkownika podjętej w oknie konfiguracji oraz w ustawieniach okna modułu Automations.
- Silnik kolejek modułu Automations może zostać w ten sposób powiązany z silnikiem kolejek opisanym w punkcie 13.4.
- Procesy zdefiniowane w Automations mogą obejmować zadania z dowolnego innego modułu platformy, zgodnie z zakresem powiązań ustanowionym przez użytkownika.

### 11.4. Browser

| Pole | Treść |
|---|---|
| Cel | Współdzielone przeglądanie internetu — dla sytuacji, w których użytkownik i AI muszą pracować nad tą samą treścią internetową w tym samym czasie (analiza strony, porównanie ofert, wspólna weryfikacja informacji), zamiast przekazywać sobie linki i fragmenty tekstu. |
| Okna | Chat Window, Execution Loop Window, Browser Window, Sources Panel, Notes Panel |
| Funkcjonalności | analiza stron, wyszukiwanie informacji, wspólny podgląd użytkownika i AI, analiza dokumentów online |

**Typowy scenariusz:**
1. Otwarcie strony w Browser Window.
2. AI — mając ten sam podgląd — odpowiada na pytania dotyczące treści i wyszukuje powiązane informacje.
3. Odnotowanie istotnych fragmentów w Notes Panel.
4. Budowa listy źródeł wykorzystanych w toku pracy w Sources Panel.

**Powiązania:**
- Źródła zebrane w module Browser mogą, decyzją użytkownika, zasilać dalszą pracę badawczą w module Research (11.5).
- Odnotowane materiały mogą trafiać do Library (11.6) w celu trwałego przechowania.

### 11.5. Research

| Pole | Treść |
|---|---|
| Cel | Realizacja badań, analiz i opracowań dla użytkowników prowadzących pracę wymagającą systematycznego zbierania, porządkowania i syntezowania informacji z wielu źródeł (od analiz rynkowych po opracowania konkurencyjne), w sposób bardziej ustrukturyzowany niż pojedyncza rozmowa z modelem. |
| Okna | Chat Window, Execution Loop Window, Research Workspace, Sources Manager, Findings Panel, Report Builder, Export Panel |
| Funkcjonalności | raporty, analizy rynku, benchmarki, analiza konkurencji, wieloźródłowe badania |

**Typowy scenariusz:**
1. Zgromadzenie źródeł w Sources Manager.
2. Odnotowywanie ustaleń cząstkowych w Findings Panel w miarę postępu analizy.
3. Skompletowanie raportu końcowego w Report Builder.
4. Eksport przez Export Panel do formatu wymaganego przez odbiorcę.

**Powiązania:**
- Research może korzystać ze źródeł zebranych w module Browser (11.4) oraz z dokumentów przechowywanych w Library (11.6), zgodnie z powiązaniami skonfigurowanymi przez użytkownika.
- Wyniki pracy badawczej mogą być dalej redagowane w module Studio (11.1) lub prezentowane wielomodelowo w module Roundtable (11.8) przed opracowaniem ostatecznych wniosków.

### 11.6. Library

| Pole | Treść |
|---|---|
| Cel | Centralne repozytorium wiedzy i plików — jedno, uporządkowane miejsce przechowywania materiałów wykorzystywanych i wytwarzanych w toku pracy; trwała pamięć zewnętrzna wspólna dla wielu środowisk i modułów. |
| Okna | Chat Window, Execution Loop Window, Library Explorer, Tags & Collections, File Preview, Versioning Panel |
| Funkcjonalności | przechowywanie plików, katalogowanie wiedzy, odbiór artefaktów generowanych przez AI, zarządzanie dokumentami |

**Typowy scenariusz:**
1. Porządkowanie zgromadzonych materiałów za pomocą Tags & Collections.
2. Przegląd zawartości bez opuszczania modułu dzięki File Preview.
3. Śledzenie kolejnych wersji tego samego dokumentu w Versioning Panel, w tym wersji wygenerowanych automatycznie przez AI w innych modułach.

**Powiązania:**
- Library może pełnić rolę centralnego repozytorium dla materiałów wykorzystywanych w Studio (11.1), Research (11.5) i Browser (11.4), zgodnie z powiązaniami ustanowionymi przez użytkownika.
- Project Library w module Workspace (11.2) stanowi jej odpowiednik ograniczony do zakresu pojedynczego projektu.

### 11.7. Translate

| Pole | Treść |
|---|---|
| Cel | Wielojęzyczne tłumaczenia w trybie wielozadaniowym dla użytkowników pracujących równolegle z treścią w wielu językach (zespoły tłumaczeniowe, lokalizacyjne, użytkownicy indywidualni obsługujący wielojęzyczną komunikację), którym pojedyncze, sekwencyjne tłumaczenie w oknie czatu nie wystarcza. |
| Okna | Chat Window, Execution Loop Window, Source Panel, Translation Panels, Glossary Manager |
| Funkcjonalności | równoległe tłumaczenia, obsługa wielu języków, zarządzanie terminologią, lokalizacja treści |

**Typowy scenariusz:**
1. Umieszczenie tekstu źródłowego w Source Panel.
2. Równoczesne tłumaczenie na wiele języków widocznych w osobnych Translation Panels.
3. Zapewnienie spójności terminologii przez Glossary Manager, w którym użytkownik definiuje preferowane odpowiedniki kluczowych pojęć.

**Powiązania:**
- Translate może zostać skonfigurowany do korzystania z operacji kontekstowych AI właściwych modułowi Studio (11.1) przy pracy dwujęzycznej oraz może pełnić rolę etapu pośredniego w pracy prowadzonej w tym module, gdy dokument wymaga wersji wielojęzycznej.
- Mechanizm operacji kontekstowych AI pozostaje natywnie mechanizmem modułu Studio; jego wspólne wykorzystanie ustanawia użytkownik w oknie konfiguracji.

### 11.8. Roundtable

| Pole | Treść |
|---|---|
| Cel | Współpraca wielu modeli AI nad wspólnym problemem — dla zadań, w których pojedyncza odpowiedź jednego modelu jest niewystarczająca, a wartość wynika z konfrontacji perspektyw, wzajemnej krytyki i wypracowania wspólnego stanowiska. |
| Okna | Chat Window, Execution Loop Window, Model Panels, Debate Panel, Moderator Panel, Consensus Panel |
| Funkcjonalności | współpraca modeli, debaty, wymiana argumentów, porównywanie odpowiedzi, wypracowywanie konsensusu |

**Typowy scenariusz:**
1. Kilka modeli widocznych równolegle w Model Panels odpowiada na to samo zagadnienie.
2. Debate Panel rejestruje wymianę argumentów między nimi.
3. Moderator Panel pozwala użytkownikowi ukierunkować dyskusję.
4. Consensus Panel gromadzi finalne, uzgodnione stanowisko.

**Powiązania:**
- Roundtable stanowi funkcjonalne rozwinięcie idei wielomodelowości realizowanej w pełniejszej, zorientowanej na role formie przez środowisko MultitaskingAI (rozdz. 9.4, 13) — różnica polega na tym, że Roundtable koncentruje się na debacie i konsensusie, a nie na podziale ról wykonawczych i orkiestracji procesu.

### 11.9. Design

| Pole | Treść |
|---|---|
| Cel | Projektowanie i tworzenie zasobów wizualnych dla użytkowników generujących i edytujących materiały graficzne — od ilustracji i elementów brandingowych po makiety interfejsu — z bezpośrednim wsparciem AI na każdym etapie procesu twórczego. |
| Okna | Chat Window, Execution Loop Window, Design Board, Assets Panel, Prompt Builder, Preview Window |
| Funkcjonalności | generowanie grafiki, edycja obrazów, ilustracje, branding, UI/UX, materiały marketingowe |

**Typowy scenariusz:**
1. Precyzyjne formułowanie poleceń generujących grafikę w Prompt Builder.
2. Gromadzenie wygenerowanych zasobów w Assets Panel.
3. Zestawianie zasobów i dalsza praca koncepcyjna nad spójną kompozycją wizualną w Design Board.

**Powiązania:**
- Zasoby wygenerowane w module Design mogą, w ramach powiązań ustanowionych przez użytkownika, być wykorzystywane w module Studio (11.1) przy redagowaniu dokumentów oraz w module Apps (11.14) przy budowie interfejsu produktu.

### 11.10. Assistant

| Pole | Treść |
|---|---|
| Cel | Naturalna komunikacja głosowa z AI dla scenariuszy, w których interakcja tekstowa jest mniej wygodna niż mowa (praca w ruchu, sterowanie zadaniami bez klawiatury, preferencje użytkownika). Konfigurowany jako komponent własny — profil asystenta — na stronie głównej (rozdz. 7, strefa 2), a wykorzystywany operacyjnie w środowisku TalkIn. |
| Okna | Chat Window, Execution Loop Window, Voice Console, Actions Monitor, Activity Feed |
| Funkcjonalności | komunikacja głosowa, wykonywanie poleceń, sterowanie zadaniami, obsługa aplikacji, realizacja działań wieloetapowych |

**Typowy scenariusz:**
1. Wydawanie poleceń głosowych przez Voice Console.
2. Podgląd bieżącego statusu ich realizacji w Actions Monitor.
3. Chronologiczny zapis wykonanych działań w Activity Feed, pozwalający odtworzyć przebieg wieloetapowego zlecenia zrealizowanego głosowo.

**Powiązania:**
- Assistant współdzieli charakter komunikacji głosowej z funkcją globalną Always On Display (12.2), różniąc się od niej tym, że działa jako pełnoprawny moduł osadzony w konkretnym środowisku, a nie jako warstwa obecna ponad całą platformą.

### 11.11. Terminal

| Pole | Treść |
|---|---|
| Cel | Praca z konsolami i środowiskami wykonawczymi dla użytkowników technicznych potrzebujących bezpośredniego dostępu do powłok systemowych i narzędzi wiersza polecenia z poziomu platformy, bez przełączania się do zewnętrznej aplikacji terminala. |
| Okna | Chat Window, Execution Loop Window, Terminal Tabs, Output Console, Process Monitor |
| Funkcjonalności | PowerShell, CMD, Bash, Node.js, Python, narzędzia CLI |

**Typowy scenariusz:**
1. Prowadzenie równolegle wielu sesji w różnych powłokach w Terminal Tabs.
2. Gromadzenie ich wyniku w Output Console.
3. Obserwacja i kontrola uruchomionych procesów w Process Monitor, w tym tych zainicjowanych poleceniem wydanym przez AI.

**Powiązania:**
- Terminal dostarcza warstwy wykonawczej dla poleceń wydawanych z poziomu modułu Developer (11.12) oraz Diagnostics (11.13), z którymi jest najściślej powiązany funkcjonalnie w ramach środowiska CodeStudio.

### 11.12. Developer

| Pole | Treść |
|---|---|
| Cel | Tworzenie i rozwój kodu — dla pracy programistycznej wymagającej pełnego zestawu narzędzi inżynierskich (edytora kodu, kontroli wersji, podglądu wyniku budowania) zintegrowanych ze wsparciem AI działającym bezpośrednio w kontekście danego repozytorium. |
| Okna | Chat Window, Execution Loop Window, Code Editor, Project Tree, Git Panel, Build Output |
| Funkcjonalności | generowanie kodu, refaktoryzacja, analiza architektury, dokumentacja, testowanie |

**Typowy scenariusz:**
1. Nawigacja po strukturze projektu w Project Tree.
2. Edycja kodu w Code Editor przy wsparciu AI udzielanym przez Chat Window.
3. Zarządzanie zmianami przez Git Panel.
4. Obserwacja wyniku kompilacji lub budowania w Build Output.

**Powiązania:**
- Developer korzysta z Terminala (11.11) jako warstwy wykonawczej oraz z modułu Diagnostics (11.13) przy analizie błędów wykrytych w toku pracy.
- Stanowi jeden z komponentów wykorzystywanych przy budowie kompletnych produktów w module Apps (11.14).

### 11.13. Diagnostics

| Pole | Treść |
|---|---|
| Cel | Analiza oraz usuwanie problemów technicznych — dla sytuacji, w których punktem wyjścia jest nie tworzenie nowej funkcjonalności, lecz zrozumienie przyczyny istniejącego błędu lub spadku wydajności na podstawie logów, komunikatów błędów i obserwowalnych objawów. |
| Okna | Chat Window, Execution Loop Window, Diagnostics Center, Logs Viewer, Errors Panel, Recommendations Panel |
| Funkcjonalności | debugowanie, analiza logów, analiza błędów, diagnostyka wydajności |

**Typowy scenariusz:**
1. Zebranie materiału źródłowego w Logs Viewer i Errors Panel.
2. Agregacja go w spójny obraz stanu systemu w Diagnostics Center.
3. Przedstawienie sugerowanych kroków naprawczych w Recommendations Panel, wypracowanych przez AI na podstawie zebranych danych.

**Powiązania:**
- Diagnostics współpracuje z modułem Developer (11.12) przy wdrażaniu poprawek wynikających z analizy oraz z Terminalem (11.11) przy odtwarzaniu i weryfikacji objawów błędu w rzeczywistym środowisku wykonawczym.

### 11.14. Apps

| Pole | Treść |
|---|---|
| Cel | Budowa kompletnych produktów cyfrowych — realizacja pełnych projektów aplikacyjnych (od architektury, przez warstwę kliencką i warstwę serwerową, po wdrożenie) w jednym, spójnym procesie obejmującym wszystkie etapy powstawania produktu, a nie pojedyncze fragmenty kodu. |
| Okna | Chat Window, Execution Loop Window, Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel |
| Funkcjonalności | aplikacje webowe, mobilne, desktopowe, API, warstwa serwerowa, warstwa kliencka, repozytoria, dokumentacja techniczna |

**Typowy scenariusz:**
1. Zaprojektowanie architektury rozwiązania w Architecture Designer.
2. Równoległa praca nad Frontend Workspace i Backend Workspace.
3. Zarządzanie końcowym etapem udostępnienia gotowego produktu przez Deployment Panel.

**Powiązania:**
- Apps integruje funkcjonalności modułów Developer (11.12) i Terminal (11.11) w ramach jednego, ustrukturyzowanego procesu budowy produktu.
- Zasoby wizualne dla warstwy klienckiej mogą pochodzić z modułu Design (11.9).
- Realizacja rozbudowanych projektów w module Apps jest naturalnym zastosowaniem środowiska MultitaskingAI (rozdz. 9.4, 13), w szczególności podziału na Executor 1 i Executor 2 przy równoległej pracy nad warstwą serwerową i warstwą kliencką.

### 11.15. Agents

| Pole | Treść |
|---|---|
| Cel | Tworzenie i zarządzanie własnymi agentami AI — dla użytkowników chcących skonfigurować trwałe, wyspecjalizowane jednostki AI (z określoną tożsamością, zestawem umiejętności i uprawnień), wykorzystywane następnie jako wykonawcy zadań w innych częściach platformy, zamiast każdorazowego definiowania kontekstu od nowa. Konfigurowany jako komponent własny na stronie głównej (rozdz. 7, strefa 2), dodatkowo dostępny jako okno modułowe w środowiskach z macierzy (rozdz. 10). |
| Okna | Chat Window, Execution Loop Window, Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center |
| Funkcjonalności | wybór modelu bazowego, definiowanie tożsamości, konfiguracja instrukcji systemowych, dodawanie umiejętności, wtyczek, konektorów, zarządzanie pamięcią, konfiguracja uprawnień |

**Typowy scenariusz:**
1. Utworzenie agenta w Agent Builder.
2. Wybór modelu bazowego w Model Configuration.
3. Dobór umiejętności w Skills Manager.
4. Podłączenie integracji zewnętrznych przez Connectors Manager.
5. Precyzyjne ustalenie zakresu uprawnień w Permissions Center, zanim agent zostanie udostępniony do wykorzystania operacyjnego.

**Charakterystyka:** Agenci są komponentem platformowym działającym ponad wszystkimi środowiskami. Po utworzeniu mogą być wykorzystywani jako wykonawcy zadań w TalkIn, WorkSpace i CodeStudio, a po przypisaniu do roli — także w środowisku MultitaskingAI (rozdz. 9.4, 13).

**Powiązania:**
- Agenci skonfigurowani w tym module mogą być przypisywani do projektów w module Workspace (11.2) przez Agent Manager.
- Mogą pełnić role wykonawcze (Executor, Coordinator, Validator i inne) w środowisku MultitaskingAI (rozdz. 9.4, 13).

---

## 12. Funkcje globalne

Funkcje globalne, w odróżnieniu od środowisk i modułów, nie tworzą własnej, izolowanej przestrzeni roboczej. Nie podlegają mechanizmowi przełączania i przeładowania z rozdziału 5.3 — pozostają dostępne w tej samej formie niezależnie od aktywnego środowiska lub modułu. Platforma udostępnia dwie funkcje globalne: Mobile oraz Always On Display. Zbiorcze zestawienie porównawcze zawiera Załącznik C.

### 12.1. Mobile

Globalny tryb dostępu do platformy, niezależny od aktywnego środowiska.

| Możliwość | Zakres |
|---|---|
| Dostęp mobilny | Pełna obsługa platformy z urządzeń mobilnych |
| Monitoring procesów | Podgląd projektów, automatyzacji, agentów, aplikacji, sesji |
| Zarządzanie zadaniami | Uruchamianie, zatrzymywanie, zatwierdzanie, modyfikowanie procesów |
| Nadzór operacyjny | Kontrola pracy AI wykonywanej w tle |

Pełne opracowanie funkcji — zakres funkcjonalny, układ interfejsu, komplet ekranów, oba kanały komunikacji operacyjnej, praca bez połączenia, powiadomienia wypychane i protokół parowania urządzenia — zawiera dokument [Funkcje mobilne](../funkcje-globalne/mobile.md).

Mobile odpowiada na potrzebę zachowania pełnej kontroli nad platformą poza stanowiskiem roboczym — w szczególności nad procesami długotrwałymi, uruchomionymi w tle automatyzacjami i projektami środowiska MultitaskingAI, które mogą wymagać zatwierdzenia, wstrzymania lub interwencji niezależnie od lokalizacji użytkownika.

**Kluczowa idea:** zarządzanie platformą z dowolnego miejsca.

### 12.2. Always On Display

Globalny agent towarzyszący. Nie posiada własnego środowiska ani modułu. Obecny we wszystkich częściach platformy jednocześnie. Stanowi nadrzędną warstwę inteligencji ponad środowiskami, modułami, projektami i sesjami.

| Możliwość | Zakres |
|---|---|
| Proaktywne doradztwo | Proponowanie działań, sugerowanie konfiguracji, wskazywanie problemów, rekomendacje kolejnych kroków |
| Pomoc kontekstowa | Rozumienie aktualnego modułu, środowiska, projektu i historii użytkownika |
| Komunikacja głosowa | Mowa, synteza głosu, szybkie polecenia |
| Komunikacja tekstowa | Dymki kontekstowe |
| Pływający avatar | Stale widoczny element interfejsu; po aktywacji umożliwia natychmiastową interakcję |
| Dostęp do platformy | Wszystkie środowiska, projekty, sesje, historie rozmów oraz agenci użytkownika |

Always On Display nie jest przypisany do żadnego konkretnego kontekstu pracy — jego rolą jest właśnie brak takiego przypisania. Dzięki stałej obecności i dostępowi do pełnej historii działania użytkownika może pełnić funkcję doradczą wykraczającą poza możliwości pojedynczego modułu, a w kontekście środowiska MultitaskingAI (rozdz. 13.2) może dodatkowo pełnić rolę obserwatora lub operatora nadzorującego przebieg złożonych, wieloetapowych procesów.

**Kluczowa idea:** osobisty agent operacyjny obecny zawsze i wszędzie.

Pełną dokumentację projektową funkcji — postać wizualną, reguły wyzwalania proaktywnych sugestii, katalog rodzajów sugestii, tor głosowy i jego granicę wobec modułu Assistant, zachowanie w poszczególnych środowiskach i modułach, warstwy widoczności oraz punkty sterowania — zawiera [Always-on Display](../funkcje-globalne/always-on-display.md).

---

## 13. MultitaskingAI — rozwinięcie środowiska

Rozdział rozwija środowisko MultitaskingAI przedstawione w rozdziale 9.4. Ze względu na zakres i złożoność funkcjonalną środowisko to wykracza poza format opisu pozostałych trzech środowisk i piętnastu modułów i wymaga odrębnego omówienia: definicji, warstwy centralnej, ról wykonawczych, silnika kolejek, orkiestracji, integracji z modułem Automations, bocznej nawigacji oraz schematu przepływu.

### 13.1. Definicja

MultitaskingAI jest zaawansowanym środowiskiem wielomodelowym umożliwiającym równoczesną pracę wielu instancji AI w ramach jednego projektu, procesu lub zadania. Każdy model może pełnić odrębną rolę operacyjną i współpracować z pozostałymi przez konfigurowalne mechanizmy przekazywania zadań, wyników i kontekstu. Środowisko pozwala budować zarówno proste układy współpracy, jak i złożone, wieloetapowe systemy autonomiczne.

W odróżnieniu od piętnastu modułów z rozdziału 11, które organizują pracę pojedynczego użytkownika wspieranego przez AI, MultitaskingAI zmienia jednostkę organizacji pracy — przedmiotem konfiguracji przestaje być pojedyncza interakcja z modelem, a staje się cały zespół modeli, z których każdy pełni określoną rolę we wspólnym procesie.

### 13.2. Warstwa centralna

**Always On Display** ([Always-on Display](../funkcje-globalne/always-on-display.md)) — globalny agent użytkownika, nadrzędna warstwa nad całym systemem. Posiada dostęp do wszystkich środowisk, projektów, agentów, sesji, historii i procesów. Może współpracować z MultitaskingAI jako obserwator lub operator procesu.

Obecność Always On Display jako warstwy centralnej wynika bezpośrednio z jego charakterystyki (rozdz. 12.2): ponieważ ma dostęp do pełnego kontekstu działania użytkownika i nie jest przypisany do pojedynczego modułu, jest naturalnym kandydatem do roli nadzorczej nad procesem angażującym wiele modeli i ról jednocześnie.

### 13.3. Role

MultitaskingAI opiera się na podziale pracy między precyzyjnie zdefiniowane role, z których każda odpowiada za inny aspekt procesu. Poniższa tabela zestawia role wraz z definicją i kluczowymi możliwościami; szczegółową macierz odpowiedzialności zawiera Załącznik B.

| Rola | Definicja | Kluczowe możliwości |
|---|---|---|
| Executor 1 — główny wykonawca | Realizuje pracę (tworzenie dokumentów, analiza danych, programowanie, projektowanie, budowa aplikacji, przetwarzanie materiałów); nie zarządza procesem. Rola podstawowa i punkt wyjścia każdej konfiguracji | Wykonywanie zadań z kolejki, realizacja poleceń użytkownika i Koordynatora, korzystanie z pamięci projektu, uruchamianie podagentów |
| Subagent Network | Mechanizm uruchamiany przez wykonawcę, nie odrębna rola zespołu | Do 15 wyspecjalizowanych podagentów; wcielenia wskazane w Załączniku B: Agent UI, Agent Backend, Agent API, Agent Security, Agent Database, Agent Testing; wyniki agreguje wykonawca |
| Coordinator — koordynator | Nie tworzy końcowego produktu; projektuje i kontroluje sposób jego wytwarzania | Planowanie (etapy, zależności, harmonogramy, logika przepływu), podział pracy, budowa promptów, sterowanie procesem (start, stop, pauza, wznowienie, przekazanie, powtórzenie, walidacja), zarządzanie kolejką |
| Executor 2 — równoległy wykonawca | Drugi niezależny model wykonawczy umożliwiający procesy równoległe — warstwę serwerową obok warstwy klienckiej, badania obok raportu, kod obok dokumentacji | Tryby współpracy z pierwszym wykonawcą: praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna |
| Executor 3 / Validator — czwarty model | Dowolna rola kontrolna lub doradcza, nieograniczona do funkcji wykonawczej; domyka strukturę zespołu | Wcielenia: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator |

Rozróżnienie między Coordinatorem a Executorem jest fundamentalne: Executor odpowiada za wykonanie, Coordinator za jego organizację — odpowiedzialności celowo rozdzielone, aby żaden pojedynczy model nie musiał jednocześnie planować i realizować.

### 13.4. Silnik kolejek

Pozwala definiować kolejki globalne, lokalne, dla modeli, dla agentów oraz dla projektów. Akcje: enqueue, dequeue, delay, retry, pause, resume, split, merge, route, branch, condition. Silnik kolejek jest mechanizmem technicznym leżącym u podstaw zarządzania kolejką realizowanego przez Coordinatora — udostępnia operacje, za pomocą których zadania są dodawane, wstrzymywane, dzielone, łączone lub kierowane warunkowo między wykonawcami i podagentami, niezależnie od poziomu struktury (globalnego, projektowego, modelu czy agenta), na którym kolejka została zdefiniowana. Szablon definicji kolejki zawiera Załącznik D.

### 13.5. Orkiestracja

MultitaskingAI pozwala definiować zależności pomiędzy modelami, agentami, zadaniami, kolejkami, automatyzacjami oraz projektami. Warstwa orkiestracji nadaje sens funkcjonowaniu ról i kolejek jako spójnej całości — odpowiada za to, że zależności zdefiniowane przez Coordinatora są respektowane w toku wykonania: zadanie Executora 2 nie rozpocznie się przed zakończeniem etapu Executora 1, jeśli taka zależność została ustalona, a wynik pracy trafi do właściwej kolejki lub roli w kolejnym kroku procesu.

### 13.6. Integracja z Automations

Po skonfigurowaniu połączenia między środowiskiem MultitaskingAI a modułem Automations — decyzją użytkownika w oknie konfiguracji oraz w ustawieniach okna modułu Automations (rozdz. 11.3) — możliwe jest zbudowanie pełnego autonomicznego systemu realizacji projektów. Poniższa tabela wskazuje rolę każdego elementu w tak skonfigurowanej pętli.

| Element | Rola w autonomicznym systemie |
|---|---|
| Coordinator | Planuje proces |
| Automations | Zarządza harmonogramami i kolejkami |
| Wykonawcy (Executor 1, Executor 2) | Realizują zadania |
| Validator (Executor 3) | Kontroluje jakość |
| Always On Display | Nadzoruje całość z perspektywy użytkownika |

Integracja łączy możliwości środowiska MultitaskingAI (role, orkiestracja, silnik kolejek) z operacyjnym zapleczem modułu Automations (11.3) — harmonogramowaniem i trwałym zarządzaniem cyklicznym wykonaniem. Efektem jest system zdolny do działania ciągłego, zgodnie z zasadą konfigurowalności zależności z rozdziału 6.

### 13.7. Kluczowa idea

MultitaskingAI zarządza zespołem modeli AI, agentów, kolejek i procesów, umożliwiając realizację wieloetapowych, autonomicznych projektów złożonych z dziesiątek lub setek powiązanych działań. Jest wirtualnym zespołem AI zarządzanym przez użytkownika, w którym każdy model pełni określoną rolę, a cały proces może działać częściowo lub całkowicie autonomicznie.

### 13.8. Boczna nawigacja środowiska MultitaskingAI (panel orkiestracji)

Trzy środowiska modułowe — TalkIn, WorkSpace i CodeStudio — udostępniają w bocznej nawigacji listę modułów (rozdz. 8, 10). MultitaskingAI organizuje jednak nie zadania modułowe, lecz zespół modeli i agentów realizujących wspólny proces (rozdz. 13.1), dlatego jego boczna nawigacja zawiera nie listę modułów, lecz panel orkiestracji — zestaw sekcji odpowiadających warstwom sterowania tego środowiska (rozdz. 13.2–13.6). Panel orkiestracji jest odpowiednikiem bocznej nawigacji modułów: nawiguje po tym, czym w tym środowisku faktycznie się steruje — po rolach, kolejkach i orkiestracji.

Panel orkiestracji zawiera sześć sekcji.

| Sekcja | Zawartość | Podstawa | Powiązane mechanizmy i okna |
|---|---|---|---|
| Zespoły | Zapisane konfiguracje zespołu — presety ról, powiązań i kolejek; zapis, wczytanie, duplikowanie | rozdz. 13.1 | konfiguracje ról (13.3) |
| Role | Cztery okna robocze: Executor 1, Executor 2, Coordinator, Executor 3 / Validator; przypisanie agentów z modułu Agents do ról; Subagent Network | rozdz. 13.3, 11.15 | Agent Builder, Permissions Center (11.15) |
| Kolejki | Definicje kolejek globalnych, lokalnych, modeli, agentów i projektów oraz ich akcje | rozdz. 13.4 | silnik kolejek (13.4), Queue Manager (11.3) |
| Orkiestracja | Zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami | rozdz. 13.5 | Orchestrator (11.3) |
| Harmonogram i automatyki | Harmonogram pracy ciągłej oraz wpięte automatyki z modułu Automations | rozdz. 13.6, 11.3 | Scheduler, Execution Monitor (11.3) |
| Monitor procesu | Podgląd przebiegu pętli, statusy przebiegów, hierarchia decyzji; miejsce nadzoru Always On Display | rozdz. 13.2, 9.4, 12.2 | Always On Display (12.2), Mobile (12.1) |

Zgodnie z zasadą pełnej konfigurowalności (rozdz. 6, 14) kolejność i widoczność sekcji panelu orkiestracji podlegają konfiguracji; przedstawiony zestaw sześciu sekcji jest zestawem domyślnym, obowiązującym przy braku odmiennego ustawienia.

### 13.9. Schemat przepływu orkiestracji

Poniższy schemat przedstawia w formie tekstowej przepływ pracy w środowisku MultitaskingAI po skonfigurowaniu pełnej pętli: nadzór Always On Display, planowanie przez Coordinatora, równoległe wykonanie przez wykonawców wraz z Subagent Network, kontrolę przez Executora 3 / Validatora, obsługę przez silnik kolejek i orkiestrację oraz integrację z modułem Automations zapewniającą ciągłość pracy.

```
                    Always On Display   (warstwa centralna — nadzór / operator)
                             │
                        Coordinator
              (plan · podział pracy · budowa promptów · sterowanie · kolejka)
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
    Executor 1          Executor 2        Executor 3 / Validator
   (wykonanie)      (wykonanie równoległe)  (kontrola: Validator,
         │                   │               Reviewer, Security Auditor,
         ▼                   ▼               Architect, QA Lead, Arbitrator …)
  Subagent Network     Subagent Network
  (do 15 podagentów)   (do 15 podagentów)
         │                   │
         └────────►  Silnik kolejek  ◄────────┘
          (enqueue · dequeue · delay · retry · pause · resume
           · split · merge · route · branch · condition)
                             │
                       Orkiestracja
        (respektowanie zależności między rolami, kolejkami, zadaniami)
                             │
                  Integracja z Automations
          (harmonogram · cykliczność → praca ciągła 24/7/365)
```

---

## 14. Zasady nadrzędne funkcjonalności

Poniższe pięć zasad ma charakter nadrzędny wobec wszystkich elementów opisanych we wcześniejszych rozdziałach — obowiązują jednolicie we wszystkich czterech środowiskach (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) oraz we wszystkich piętnastu modułach platformy, wyznaczając wspólny standard, jakiemu podlega cała platforma.

| Zasada | Istota | Kluczowa uwaga |
|---|---|---|
| 1. Pełna kompozycyjność | W każdym module dostępne są wszystkie funkcje i ich kompozycje; Operator komponuje własne układy z komponentów: prompt, zadanie, akcja, orkiestracja, subagenci, pełna pętla z Koordynatorem, pętla wykonawcza, sekwencje i kompilacje | Żaden moduł nie ogranicza funkcji do wąskiego, z góry ustalonego podzbioru; przykład — złożenie pełnej pętli Coordinator-Executor z Subagent Network (rozdz. 13.3) |
| 2. Pełna konfigurowalność (zasada centralna) | Każdą funkcję, każde zachowanie i każdą zależność między elementami platformy konfiguruje się z poziomu okna konfiguracji; nic nie jest zaszyte na stałe w sposób niedostępny dla użytkownika | Zasada centralna wobec pozostałych; rozstrzyga stopień powiązania funkcji oraz odrębności kontekstów i historii sesji; pełne rozwinięcie z punktami izolacji w rozdz. 6 |
| 3. Jawność i konfigurowalność zależności | Powiązania między funkcjami nie są zaszyte na stałe — są jawne i konfigurowalne przez użytkownika (zgodnie z zasadą 2 i rozdz. 6); platforma nie narzuca gotowego zestawu zależności | Chroni przed narastaniem ukrytych, niejawnych powiązań — nie przed liczbą powiązań jako taką |
| 4. Pełna orkiestracja | Platforma prowadzi równoległą pracę wielu modeli i agentów oraz udostępnia pełny zestaw funkcji operacyjnych sterujących tą pracą | Zasada obowiązuje we wszystkich środowiskach; najpełniejsze wcielenie — środowisko MultitaskingAI (rozdz. 9.4, 13) |
| 5. Warstwy widoczności interfejsu | Jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna; interfejs ujawnia możliwości stopniowo, zależnie od kontekstu, roli użytkownika i wykonywanej czynności, zachowując maksymalną moc funkcjonalną przy minimalnej złożoności wizualnej | Każdy element interfejsu należy do jednej z czterech warstw widoczności; każda ukryta funkcja pozostaje osiągalna jednym kliknięciem, skrótem klawiszowym albo poleceniem języka naturalnego; pełne rozwinięcie w rozdz. 4 |

Zasada druga — pełna konfigurowalność — ma charakter centralny: rozstrzyga, w jakim stopniu funkcje platformy są ze sobą powiązane oraz w jakim stopniu konteksty i historie sesji pozostają odrębne. Pozostałe zasady dopełniają ten fundament, wyznaczając dostępność funkcji (kompozycyjność), sposób ustanawiania powiązań (jawność), zakres równoległej pracy modeli i agentów (pełna orkiestracja) oraz sposób prezentacji funkcji w interfejsie (warstwy widoczności).

---

## 15. Słowniczek pojęć

Poniższa tabela zbiera pojęcia platformy wraz ze zwięzłą definicją i odesłaniem do rozdziału, w którym opisano je szczegółowo.

| Pojęcie | Definicja |
|---|---|
| Środowisko | Najwyższy poziom organizacji pracy; definiuje układ interfejsu, dostępne moduły, typ pracy, sposób działania AI, dostępne narzędzia, model nawigacji, historię sesji i kontekst użytkownika. Cztery środowiska: TalkIn, WorkSpace, CodeStudio i MultitaskingAI. |
| Moduł | Wyspecjalizowany obszar roboczy działający wewnątrz środowiska, odpowiadający za konkretne zadanie biznesowe; definiuje zestaw okien operacyjnych, dedykowane narzędzia, przebieg pracy oraz typ danych wykorzystywanych przez AI. Piętnaście modułów wymieniono w rozdziale 11, każdy dostępny w wybranych środowiskach zgodnie z macierzą (rozdz. 10). Element platformy — w odróżnieniu od komponentu własnego. |
| Komponent własny | Nazwany wytwór użytkownika (konkretna automatyka, agent, projekt lub profil asystenta), tworzony i zapisywany w oknie konfiguracji na stronie głównej (rozdz. 7, strefa 2) lub w oknie modułu, a następnie wybieralny w trakcie sesji w module, w którym ma zastosowanie. Pojęcie odrębne od modułu (rozdz. 2.3). |
| Okno operacyjne | Pojedynczy element interfejsu wchodzący w skład zestawu okien danego modułu — Studio Editor, Workflow Builder, Code Editor i pozostałe okna wymienione w Załączniku A; właściwa przestrzeń wykonywania zadań przypisana do konkretnej funkcjonalności modułu. Zbiorcze zestawienie okien zawiera Załącznik A. |
| Użytkownik | Osoba prowadząca pracę w platformie; zleca zadania, zatwierdza wyniki, przerywa i koryguje działania. Komunikuje się z Wykonawcą kanałem pierwszym — Chat Window (rozdz. 3.1). |
| Koordynator | Komponent orkiestrujący platformy; dekomponuje zlecenie na zadania, przydziela je, nadzoruje przebieg pętli wykonawczej i kontroluje realizację procesów. Komunikuje się z Wykonawcą kanałem drugim — Execution Loop Window (rozdz. 3.2). W środowisku MultitaskingAI rolę tę pełni Coordinator (rozdz. 13.3). |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania przydzielone przez Użytkownika i Koordynatora oraz zwracający wyniki. W środowisku MultitaskingAI rolę tę pełnią Executor 1 i Executor 2 (rozdz. 13.3). |
| Chat Window | Główne okno komunikacji kanału Użytkownik ↔ Wykonawca; centralny punkt pracy użytkownika i podstawowy mechanizm sterowania wszystkimi procesami platformy. Zajmuje lewą, stałą kolumnę o pełnej wysokości obszaru roboczego i należy do warstwy widoczności 1 (rozdz. 3.1, 4.1). |
| Execution Loop Window | Okno pętli wykonawczej kanału Koordynator ↔ Wykonawca; mechanizm koordynacji zadań, nadzoru, orkiestracji i kontroli realizacji procesów. Otwierane jako kolumna sąsiadująca z Chat Window (rozdz. 3.2). |
| Warstwa widoczności | Jedna z czterech warstw, do której należy każdy element interfejsu: zawsze widoczna, widoczna na żądanie, rozwinięcie kontekstowe lub funkcja ekspercka; wyznacza sposób dostępu do elementu (rozdz. 4.1). |
| Układ pionowy (podział lewa–prawa) | Standard rozmieszczenia okien platformy: wszystkie okna robocze, okna komunikacji, przestrzenie modułowe i okna pomocnicze zajmują kolumny sąsiadujące poziomo; regulacji podlega wyłącznie ich szerokość, a okna pomocnicze otwierane są jako rozszerzenia boczne (rozdz. 3.4). |
| Czat | Uniwersalny mechanizm komunikacji użytkownika z AI, obecny we wszystkich środowiskach i modułach, lecz nie będący stałym, niezmiennym komponentem interfejsu; przy każdej zmianie środowiska lub modułu jego wygląd, historia, dostępne narzędzia i kontekst dopasowują się do bieżącego kontekstu, o ile użytkownik nie skonfiguruje ich współdzielenia (rozdz. 6). |
| Sesja | Pojedyncza aktywność robocza prowadzona w module, reprezentowana kartą sesji w oknie środowiska (rozdz. 8); posiada własny układ okien, historię i kontekst, których zakres współdzielenia z innymi sesjami jest konfigurowalny (rozdz. 6). |
| Projekt | Jednostka organizacji pracy w module Workspace (11.2), grupująca instrukcje, pamięć kontekstową, bibliotekę i przypisanych agentów właściwe jednemu przedsięwzięciu; rozdzielenie kontekstu między projektami jest konfigurowalne (rozdz. 6). |
| Centrum dowodzenia | Inna nazwa strony głównej aplikacji (rozdz. 7); przedpokój przed strefą roboczą, z trzema strefami: wyborem środowiska, konfiguracją komponentów własnych oraz ustawieniami platformy. |
| Warstwa środowisk (Environment Layer) | Warstwa architektury odpowiedzialna za definiowanie sposobu pracy użytkownika poprzez cztery dostępne środowiska. |
| Warstwa modułów (Workspace Layer) | Warstwa architektury obejmująca wszystkie moduły platformy, definiująca wyspecjalizowane przestrzenie robocze z własną powłoką, układem okien, historią sesji, kontekstem i przebiegiem pracy. Pojęcie odrębne od modułu Workspace (11.2) — zob. rozdz. 5.2. |
| Warstwa funkcji globalnych (Global Features Layer) | Warstwa architektury obejmująca funkcje działające ponad wszystkimi środowiskami, modułami i projektami, nieposiadające własnej powłoki roboczej i nietworzące nowych sesji. |
| Funkcja globalna | Funkcjonalność platformy dostępna niezależnie od aktywnego środowiska lub modułu, nietworząca własnej przestrzeni roboczej. Dwie funkcje globalne: Mobile i Always On Display. |
| Agent | Jednostka AI o zdefiniowanej tożsamości, zestawie umiejętności, konektorów i uprawnień, tworzona i konfigurowana jako komponent własny w module Agents (11.15); po utworzeniu może być wykorzystywana jako wykonawca zadań w dowolnym środowisku, w tym jako uczestnik procesu w środowisku MultitaskingAI. |
| Always On Display | Globalny agent towarzyszący, nieposiadający własnego środowiska ani modułu, obecny jednocześnie we wszystkich częściach platformy; nadrzędna warstwa inteligencji ponad środowiskami, modułami, projektami i sesjami, zdolna do proaktywnego doradztwa oraz pełnienia roli obserwatora lub operatora procesów środowiska MultitaskingAI. |
| MultitaskingAI | Czwarte, najbardziej zaawansowane konfiguracyjnie środowisko platformy, umożliwiające — poza czterema oknami roboczymi (rolami) — pełną orkiestrację autonomicznej, wieloagentowej pętli pracy ciągłej, z kolejkowaniem, rolami, akcjami, harmonogramem i hierarchią decyzji; opisane szczegółowo w rozdziale 13. |
| Executor | Rola w środowisku MultitaskingAI odpowiedzialna za faktyczne wykonanie pracy (tworzenie dokumentów, analiza danych, programowanie, projektowanie, budowa aplikacji, przetwarzanie materiałów), bez odpowiedzialności za zarządzanie procesem. Równoległą pracę prowadzą dwaj wykonawcy (Executor 1, Executor 2). |
| Subagent Network | Mechanizm pozwalający Executorowi uruchomić jednocześnie do 15 wyspecjalizowanych podagentów realizujących wąsko zdefiniowane zadania cząstkowe, których wyniki są następnie agregowane przez Executora. |
| Coordinator (Koordynator) | Rola w środowisku MultitaskingAI odpowiedzialna za planowanie, podział pracy, budowę promptów dla wykonawców, sterowanie procesem oraz zarządzanie kolejką; nie tworzy końcowego produktu, lecz projektuje i kontroluje sposób jego wytwarzania. |
| Validator | Jedna z przykładowych ról czwartego modelu (Executor 3) w środowisku MultitaskingAI, odpowiedzialna za kontrolę jakości pracy pozostałych modeli; może przyjmować również inne wcielenia: Reviewer, Security Auditor, Architect, Product Owner, QA Lead lub Arbitrator. |
| Silnik kolejek | Mechanizm środowiska MultitaskingAI pozwalający definiować kolejki globalne, lokalne, dla modeli, dla agentów oraz dla projektów, obsługujący akcje: enqueue, dequeue, delay, retry, pause, resume, split, merge, route, branch i condition. Może zostać powiązany z silnikiem kolejek modułu Automations (11.3) decyzją użytkownika (rozdz. 6). |
| Orkiestracja | Mechanizm środowiska MultitaskingAI definiujący zależności pomiędzy modelami, agentami, zadaniami, kolejkami, automatyzacjami oraz projektami, zapewniający respektowanie tych zależności w toku wykonania procesu. |
| Panel orkiestracji | Boczna nawigacja środowiska MultitaskingAI, zastępująca listę modułów obecną w pozostałych środowiskach; złożona z sześciu sekcji (Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu) odpowiadających własnym warstwom sterowania tego środowiska. Opisany w rozdziale 13.8. |
| Izolacja kontekstu | Rozdzielenie historii, pamięci i kontekstu między kartami sesji, modułami, środowiskami lub projektami; domyślnie odrębna dla każdej nowej karty sesji, z możliwością skonfigurowania współdzielenia (rozdz. 6). |
| Izolacja techniczna | Ograniczenie procesu sesji obejmujące zakresy zdefiniowane w rozdziale 9 dokumentu Architektury (katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania); domyślnie żaden zakres nie jest aktywny (rozdz. 6.6). |
| Poziom zasięgu (reguły izolacji) | Poziom, na którym obowiązuje reguła izolacji: globalny, środowisko, moduł, para modułów, projekt, karta sesji, rola lub okno komunikacji; reguły z różnych poziomów rozstrzyga zasada pierwszeństwa zasięgu najbardziej szczegółowego (rozdz. 6.5). |
| Profil izolacji | Nazwany, zapisywalny zestaw ustawień izolacji kontekstu i izolacji technicznej dla wskazanego poziomu zasięgu, przypisywalny do sesji lub roli (rozdz. 6.6, Załącznik D). |

---

## 16. Kryteria odbioru

Koncepcja platformy jest zrealizowana zgodnie z niniejszym opracowaniem, gdy spełnione są łącznie poniższe warunki. Każdy warunek jest sprawdzalny bez odwołania do intencji autora — wyłącznie przez odczyt kontraktu, konfiguracji albo zachowania obserwowalnego z poziomu interfejsu.

| Warunek | Sposób sprawdzenia |
|---|---|
| Platforma udostępnia dokładnie cztery środowiska (rozdz. 9) | `environment.enter` obsługuje wyłącznie identyfikatory TalkIn, WorkSpace, CodeStudio, MultitaskingAI |
| Piętnaście modułów jest rozmieszczonych zgodnie z macierzą dostępności (rozdz. 10–11), bez modułu dostępnego w środowisku nieprzewidzianym macierzą | boczna nawigacja środowiska nie prezentuje modułu spoza wiersza tej macierzy odpowiadającego środowisku |
| Strona główna otwiera się jako pierwszy widok po uwierzytelnieniu, nie środowisko (rozdz. 7) | `home.enter` jest wywoływane przed jakąkolwiek komendą obszaru `environment`, chyba że klient jawnie zgłasza powrót do otwartej karty sesji |
| Oba kanały komunikacji operacyjnej — Chat Window i Execution Loop Window — są dostępne we wszystkich piętnastu modułach jako okna wspólne (rozdz. 3, Załącznik A) | każdy wiersz Załącznika A niesie niepustą wartość w kolumnach Chat Window i Execution Loop Window |
| Warstwy widoczności interfejsu (rozdz. 4) nie mają wyjątków niezgodnych z przypisaniem warstwy w klasach `.dn-*` | przegląd komponentów `design/zasoby/css/komponenty.css` względem deklarowanej warstwy elementu |
| Izolacja kontekstu i izolacja techniczna (rozdz. 6) zachowują zasadę „brak ustawienia = wartość domyślna” | ustawienie nieskonfigurowane na żadnym poziomie zasięgu dziedziczy wartość globalną, nie zwraca błędu |
| Komponent własny jest tworzony i zapisywany wyłącznie w strefie 2 strony głównej albo w oknie modułu, nigdy jako element trwały modułu (rozdz. 2.3) | usunięcie komponentu własnego komendą `component.delete` nie usuwa modułu ani jego okien operacyjnych |
| Środowisko MultitaskingAI udostępnia sześć sekcji panelu orkiestracji zamiast listy modułów (rozdz. 13.8) | boczna nawigacja środowiska MultitaskingAI nie zawiera pozycji identycznej z listą modułów pozostałych środowisk |
| Dwie funkcje globalne (Mobile, Always On Display) działają niezależnie od aktywnego środowiska i modułu, bez tworzenia własnej sesji (rozdz. 12) | otwarcie funkcji globalnej nie zmienia identyfikatora bieżącej karty sesji ani środowiska |

---

## Załącznik A. Macierz okien operacyjnych per moduł

Zestawienie zbiera w jednym miejscu okna operacyjne wymienione przy poszczególnych modułach w rozdziale 11. Chat Window i Execution Loop Window występują w każdym z piętnastu modułów jako okna dwóch kanałów komunikacji operacyjnej (rozdz. 3) i podano je jako pierwsze; pozostałe okna wymieniono w kolejności z rozdziału 11.

| Moduł | Okna operacyjne |
|---|---|
| Studio (11.1) | Chat Window, Execution Loop Window, Studio Editor, Tools Panel, Diff/Grep Panel, Session Repository, Preview Window |
| Workspace (11.2) | Chat Window, Execution Loop Window, Project Dashboard, Instructions Panel, Context Memory, Project Library, Agent Manager |
| Automations (11.3) | Chat Window, Execution Loop Window, Workflow Builder, Scheduler, Queue Manager, Orchestrator, Execution Monitor |
| Browser (11.4) | Chat Window, Execution Loop Window, Browser Window, Sources Panel, Notes Panel |
| Research (11.5) | Chat Window, Execution Loop Window, Research Workspace, Sources Manager, Findings Panel, Report Builder, Export Panel |
| Library (11.6) | Chat Window, Execution Loop Window, Library Explorer, Tags & Collections, File Preview, Versioning Panel |
| Translate (11.7) | Chat Window, Execution Loop Window, Source Panel, Translation Panels, Glossary Manager |
| Roundtable (11.8) | Chat Window, Execution Loop Window, Model Panels, Debate Panel, Moderator Panel, Consensus Panel |
| Design (11.9) | Chat Window, Execution Loop Window, Design Board, Assets Panel, Prompt Builder, Preview Window |
| Assistant (11.10) | Chat Window, Execution Loop Window, Voice Console, Actions Monitor, Activity Feed |
| Terminal (11.11) | Chat Window, Execution Loop Window, Terminal Tabs, Output Console, Process Monitor |
| Developer (11.12) | Chat Window, Execution Loop Window, Code Editor, Project Tree, Git Panel, Build Output |
| Diagnostics (11.13) | Chat Window, Execution Loop Window, Diagnostics Center, Logs Viewer, Errors Panel, Recommendations Panel |
| Apps (11.14) | Chat Window, Execution Loop Window, Product Builder, Architecture Designer, Frontend Workspace, Backend Workspace, Deployment Panel |
| Agents (11.15) | Chat Window, Execution Loop Window, Agent Builder, Model Configuration, Skills Manager, Connectors Manager, Permissions Center |

**Uwagi.** Chat Window i Execution Loop Window są oknami wspólnymi dla wszystkich piętnastu modułów — pozostają tymi samymi kanałami komunikacji operacyjnej (rozdz. 3), rekonfigurowanymi do kontekstu bieżącego modułu (rozdz. 2.4). Chat Window zajmuje lewą kolumnę obszaru roboczego, Execution Loop Window — kolumnę z nią sąsiadującą; pozostałe okna modułu wypełniają prawą, dominującą część obszaru roboczego, a panele pomocnicze otwierane są jako rozszerzenia boczne. Preview Window występuje w dwóch modułach (Studio, Design), pełniąc w obu funkcję podglądu wyniku pracy. Nazwy zbliżone brzmieniowo oznaczają okna odrębne, właściwe różnym modułom: Process Monitor (Terminal), Execution Monitor (Automations) i Actions Monitor (Assistant) są trzema różnymi oknami; podobnie Project Library (okno modułu Workspace) i Library Explorer (okno modułu Library).

---

## Załącznik B. Zestawienie ról środowiska MultitaskingAI

Zestawienie porządkuje role opisane w rozdziale 13.3 według ich odpowiedzialności i relacji do procesu. Ostatni wiersz opisuje Subagent Network, który jest mechanizmem uruchamianym przez wykonawcę, a nie odrębną rolą zespołu.

| Rola | Zakres odpowiedzialności | Tworzy produkt końcowy | Zarządza procesem | Uruchamia Subagent Network | Przykładowe wcielenia / tryby |
|---|---|---|---|---|---|
| Executor 1 (główny wykonawca) | Faktyczna realizacja pracy; wykonywanie zadań z kolejki i poleceń Koordynatora | Tak | Nie | Tak (do 15 podagentów) | — |
| Executor 2 (równoległy wykonawca) | Drugi niezależny tor wykonania, umożliwiający procesy równoległe | Tak | Nie | Tak (do 15 podagentów) | tryby: praca niezależna, przekazywanie wyników, praca naprzemienna, praca iteracyjna |
| Coordinator (koordynator) | Planowanie, podział pracy, budowa promptów, sterowanie procesem, zarządzanie kolejką | Nie | Tak | Nie | — |
| Executor 3 / Validator (czwarty model) | Funkcja kontrolna lub doradcza wobec pracy pozostałych modeli | Zależnie od wcielenia | Nie | — | Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator |
| Subagent Network (mechanizm) | Do 15 wyspecjalizowanych podagentów uruchamianych przez wykonawcę; wyniki agreguje wykonawca | Nie (wynik agreguje wykonawca) | Nie | — | Agent UI, Agent Backend, Agent API, Agent Security, Agent Database, Agent Testing |

Warstwą centralną nadrzędną wobec wszystkich ról jest Always On Display, pełniący rolę obserwatora lub operatora procesu (rozdz. 13.2).

---

## Załącznik C. Macierz funkcji globalnych

Zestawienie porównuje dwie funkcje globalne opisane w rozdziale 12 według cech wynikających z ich przynależności do warstwy funkcji globalnych (rozdz. 5.3).

| Cecha | Mobile | Always On Display |
|---|---|---|
| Charakter | Globalny tryb dostępu do platformy | Globalny agent towarzyszący |
| Własna powłoka robocza | Nie | Nie |
| Przeładowuje środowisko | Nie | Nie |
| Tworzy nowe sesje | Nie | Nie |
| Obecność | Dostęp z urządzeń mobilnych, niezależny od aktywnego środowiska | Jednoczesna we wszystkich częściach platformy (pływający avatar) |
| Główne funkcje | Dostęp mobilny, monitoring procesów, zarządzanie zadaniami, nadzór operacyjny | Proaktywne doradztwo, pomoc kontekstowa, komunikacja głosowa i tekstowa, dostęp do pełnej historii |
| Rola w środowisku MultitaskingAI | Zatwierdzanie, wstrzymywanie i interwencja w procesy z dowolnego miejsca (rozdz. 12.1) | Warstwa centralna — obserwator lub operator procesu (rozdz. 13.2) |
| Kluczowa idea | Zarządzanie platformą z dowolnego miejsca | Osobisty agent operacyjny obecny zawsze i wszędzie |

---

## Załącznik D. Szablony konfiguracji

Poniższe szablony porządkują pola konfiguracji elementów opisanych w koncepcji. Są to szablony redakcyjne wspierające implementację; wszystkie pola podlegają zasadzie „brak ustawienia = wartość domyślna” (rozdz. 6, rozdz. 14). Każdy szablon podano najpierw w formie skróconej struktury pól (blok konfiguracji), a następnie w formie tabeli objaśniającej znaczenie i dopuszczalne wartości każdego pola.

### D.1. Szablon komponentu własnego

```
Komponent własny
  nazwa:                 <tekst>                                  # przykład: „Raport tygodniowy”
  rodzaj:                automatyka | agent | projekt | profil asystenta
  miejsce konfiguracji:  strefa 2 strony głównej (rozdz. 7.2) | okno modułu
  przeznaczenie:         <tekst>
  moduły zastosowania:   <lista modułów>                          # przykład: Developer, Workspace, TalkIn
  powiązania:            <lista jawnych, konfigurowalnych powiązań>   # rozdz. 6.2
  widoczność / zasięg:   globalny | projektowy

  # Pola swoiste wg rodzaju:
  automatyka:            przebieg pracy, harmonogram, kolejka, wywołania modeli          # rozdz. 11.3
  agent:                 model bazowy, tożsamość, instrukcje systemowe, umiejętności,
                         wtyczki, konektory, pamięć, uprawnienia                   # rozdz. 11.15
  projekt:               instrukcje, pamięć kontekstowa, biblioteka projektu,
                         konfiguracja modeli, przypisani agenci                    # rozdz. 11.2
  profil asystenta:      głos i synteza mowy, zakres poleceń, obsługiwane akcje    # rozdz. 11.10
```

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Nazwa | Nazwa własna komponentu | dowolny tekst; przykład: „Raport tygodniowy” |
| Rodzaj | Typ komponentu własnego | automatyka \| agent \| projekt \| profil asystenta |
| Miejsce konfiguracji | Okno, w którym komponent powstaje | strefa 2 strony głównej (rozdz. 7.2) lub okno modułu |
| Przeznaczenie | Krótki opis zadania komponentu | dowolny tekst |
| Moduł(y) zastosowania | Miejsce operacyjnego wykorzystania | przykład: Developer, Workspace, TalkIn |
| Powiązania | Jawne, konfigurowalne powiązania z innymi komponentami lub modułami | lista powiązań (rozdz. 6.2) |
| Widoczność / zasięg | Zakres dostępności komponentu | globalny \| projektowy |

Pola swoiste dla poszczególnych rodzajów komponentu własnego:

| Rodzaj | Pola swoiste (podstawa) |
|---|---|
| Automatyka | przebieg pracy, harmonogram, kolejka, wywołania modeli (rozdz. 11.3) |
| Agent | model bazowy, tożsamość, instrukcje systemowe, umiejętności, wtyczki, konektory, pamięć, uprawnienia (rozdz. 11.15) |
| Projekt | instrukcje, pamięć kontekstowa, biblioteka projektu, konfiguracja modeli, przypisani agenci (rozdz. 11.2) |
| Profil asystenta | głos i synteza mowy, zakres poleceń, obsługiwane akcje (rozdz. 11.10) |

### D.2. Szablon roli (środowisko MultitaskingAI)

```
Rola (środowisko MultitaskingAI)
  rola:                       Executor 1 | Executor 2 | Coordinator | Executor 3 / Validator
  wcielenie:                  Validator | Reviewer | Security Auditor | Architect |
                              Product Owner | QA Lead | Arbitrator     # dla Executor 3 / Validator
  przypisany agent lub model: agent z modułu Agents (11.15) | model bazowy
  kanał modelu:               API | CLI | SSH | HTTP                   # rozdz. 8 dokumentu Architektury
  zakres odpowiedzialności:   <tekst>
  dostęp do pamięci projektu: tak | nie
  subagent network:           tak (do 15) | nie
  relacje z innymi rolami:    niezależna | przekazywanie | naprzemienna | iteracyjna
  profil izolacji:            <odwołanie do profilu — D.4>
```

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Rola | Rola w zespole | Executor 1 \| Executor 2 \| Coordinator \| Executor 3 / Validator |
| Wcielenie | Konkretne wcielenie roli (dla Executor 3 / Validator) | Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator |
| Przypisany agent lub model | Wykonawca roli | agent z modułu Agents (11.15) albo model bazowy |
| Kanał modelu | Sposób połączenia modelu | API \| CLI \| SSH \| HTTP (rozdz. 8 dokumentu Architektury) |
| Zakres odpowiedzialności | Zadania przypisane roli | dowolny tekst |
| Dostęp do pamięci projektu | Czy rola korzysta z pamięci projektu | tak \| nie |
| Subagent Network | Możliwość uruchomienia podagentów | tak (do 15) \| nie |
| Relacje z innymi rolami | Tryb współpracy między wykonawcami | niezależna \| przekazywanie \| naprzemienna \| iteracyjna |
| Profil izolacji | Przypisany profil izolacji roli | odwołanie do profilu (D.4) |

### D.3. Szablon kolejki

```
Kolejka
  nazwa:                    <tekst>
  zasięg:                   globalna | lokalna | modelu | agenta | projektu    # rozdz. 13.4
  źródło zadań:             rola | moduł | automatyka
  obsługiwane akcje:        enqueue, dequeue, delay, retry, pause, resume,
                            split, merge, route, branch, condition
  reguły warunkowe:         <definicje branch / condition>
  powiązanie z Automations: tak | nie                              # konfigurowalne, rozdz. 6, 13.6
  priorytet:                <wartość porządkująca>
  obsługa błędów:           retry | route | pause
```

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Nazwa | Nazwa własna kolejki | dowolny tekst |
| Zasięg | Poziom, na którym kolejka jest zdefiniowana | globalna \| lokalna \| modelu \| agenta \| projektu (rozdz. 13.4) |
| Źródło zadań | Skąd pochodzą zadania kolejki | rola, moduł, automatyka |
| Obsługiwane akcje | Operacje dostępne na kolejce | enqueue, dequeue, delay, retry, pause, resume, split, merge, route, branch, condition |
| Reguły warunkowe | Warunki rozgałęzień i kierowania | definicje branch / condition |
| Powiązanie z Automations | Spięcie z silnikiem kolejek modułu Automations | tak \| nie (konfigurowalne, rozdz. 6, rozdz. 13.6) |
| Priorytet | Kolejność obsługi | wartość porządkująca |
| Obsługa błędów | Zachowanie przy błędzie | retry \| route \| pause |

### D.4. Szablon profilu izolacji

```
Profil izolacji
  nazwa profilu:         <tekst>
  poziom zasięgu:        globalny | środowisko | moduł | para modułów |
                         projekt | karta sesji | rola |
                         okno komunikacji                            # rozdz. 6.5
  izolacja kontekstu:
    historia:            współdzielona | odrębna
    pamięć:              współdzielona | odrębna
    kontekst:            współdzielony | odrębny
  izolacja techniczna (osiem zakresów — rozdz. 9 dokumentu Architektury):
    katalog roboczy sesji:                włączony | wyłączony
    środowisko procesu:                   włączony | wyłączony
    katalog danych i konfiguracji modelu: włączony | wyłączony
    dostęp sieciowy:                      włączony | wyłączony
    zakres odczytu i zapisu plików:       włączony | wyłączony
    konto i token per sesja:              włączony | wyłączony
    model procesu:                        włączony | wyłączony
    serwer wykonania:                     włączony | wyłączony
  warstwa:               domyślna | sesji                           # rozdz. 6.6
  przypisanie:           sesja | rola
  wartość domyślna:      dziedziczenie z poziomu szerszego; brak aktywnej izolacji technicznej
```

| Pole | Opis | Wartości / przykład |
|---|---|---|
| Nazwa profilu | Nazwa własna profilu | dowolny tekst |
| Poziom zasięgu | Poziom obowiązywania reguły | globalny \| środowisko \| moduł \| para modułów \| projekt \| karta sesji \| rola \| okno komunikacji (rozdz. 6.5) |
| Izolacja kontekstu — historia | Współdzielenie historii | współdzielona \| odrębna |
| Izolacja kontekstu — pamięć | Współdzielenie pamięci | współdzielona \| odrębna |
| Izolacja kontekstu — kontekst | Współdzielenie kontekstu | współdzielony \| odrębny |
| Izolacja techniczna | Osiem zakresów z rozdziału 9 dokumentu Architektury | każdy zakres: włączony \| wyłączony |
| Warstwa | Warstwa konfiguracji | domyślna \| sesji (rozdz. 6.6, rozdz. 10 dokumentu Architektury) |
| Przypisanie | Element, do którego profil jest przypisany | sesja \| rola |
| Wartość domyślna | Zachowanie przy braku ustawienia | dziedziczenie z poziomu szerszego; brak aktywnej izolacji technicznej |

---

## Załącznik E. Scenariusze użycia

Scenariusze ilustrują kompozycję funkcji opisanych w koncepcji. Nie wprowadzają nowych funkcji — pokazują współdziałanie istniejących modułów, funkcji globalnych i środowiska MultitaskingAI.

### E.1. Autonomiczna budowa aplikacji w środowisku MultitaskingAI

Użytkownik buduje kompletny produkt cyfrowy w module Apps (11.14) z wykorzystaniem środowiska MultitaskingAI. Przebieg pętli przedstawia poniższa tabela.

| Krok | Element | Działanie |
|---|---|---|
| 1 | Coordinator | Planuje etapy i przypisuje pracę wykonawcom |
| 2 | Executor 1 | Realizuje warstwę serwerową; uruchamia Subagent Network (Agent Backend, Agent API, Agent Database) |
| 3 | Executor 2 | Realizuje równolegle warstwę kliencką (tryb pracy niezależnej); uruchamia Subagent Network (Agent UI) |
| 4 | Executor 3 | Przyjmuje wcielenie Security Auditor, a w kolejnym przebiegu QA Lead — kontroluje bezpieczeństwo i jakość |
| 5 | Silnik kolejek (13.4) | Rozdziela i łączy zadania |
| 6 | Orkiestracja (13.5) | Pilnuje zależności między warstwą serwerową a warstwą kliencką |
| 7 | Integracja z Automations (13.6) | Po spięciu proces działa cyklicznie, bez stałego nadzoru |
| 8 | Always On Display (13.2) | Nadzoruje całość jako operator procesu |
| 9 | Mobile (12.1) | Użytkownik zatwierdza kluczowe kroki z urządzenia mobilnego |

Zasoby wizualne warstwy klienckiej pochodzą z modułu Design (11.9), zgodnie z powiązaniem ustanowionym przez użytkownika.

### E.2. Przepływ badawczy między modułami

Użytkownik prowadzi badanie rynkowe, komponując cztery moduły w jeden przepływ pracy (rozdz. 14, zasada 1).

```
Browser (11.4)            Research (11.5)                 Studio (11.1)
wspólna analiza stron →  Sources Manager, Findings   →  redakcja i operacje
lista źródeł w             Panel, Report Builder          kontekstowe AI
Sources Panel
        │                        │                             │
        └────────────────────────┴──────────────┬──────────────┘
                                                 ▼
                                          Library (11.6)
                                 trwałe repozytorium źródeł i wersji raportu
```

1. W module Browser (11.4) użytkownik i AI wspólnie analizują strony, budując listę źródeł w Sources Panel.
2. Źródła — decyzją użytkownika — zasilają moduł Research (11.5), gdzie w Sources Manager i Findings Panel powstaje analiza, a w Report Builder raport końcowy.
3. Raport trafia do redakcji w module Studio (11.1), gdzie operacje kontekstowe AI dopracowują styl.
4. Materiały źródłowe i wersje raportu przechowuje moduł Library (11.6) jako trwałe repozytorium.

Powiązania między modułami są jawne i ustanowione przez użytkownika (rozdz. 6.2).

### E.3. Konfiguracja punktu izolacji między projektami

Użytkownik prowadzi równolegle dwa projekty w module Workspace (11.2) i chce, by ich konteksty pozostały rozdzielone, ale współdzieliły dostęp do wspólnej biblioteki.

1. W oknie konfiguracji punktów izolacji (6.4) wybiera poziom zasięgu „projekt” (6.5).
2. W macierzy izolacji ustawia historię, pamięć i kontekst jako odrębne, a dostęp do repozytorium plików pozostawia współdzielony.
3. Ustawienia zapisuje jako profil izolacji (6.6, Załącznik D) i przypisuje go obu projektom.
4. Reguła projektowa ma pierwszeństwo przed regułą globalną, zgodnie z zasadą pierwszeństwa zasięgu najbardziej szczegółowego (6.5).

### E.4. Praca dwujęzyczna Studio i Translate

Użytkownik redaguje dokument wymagający wersji w dwóch językach. Domyślnie operacje kontekstowe AI są mechanizmem modułu Studio (11.1), a Translate (11.7) działa niezależnie.

1. Użytkownik w oknie konfiguracji ustanawia powiązanie na poziomie pary modułów (6.5), współdzieląc operacje kontekstowe AI między Studio a Translate.
2. Tekst źródłowy redagowany w Studio Editor jest tłumaczony równolegle w Translation Panels.
3. Spójność terminologii pilnuje Glossary Manager.

Powiązanie jest jawne i odwracalne (rozdz. 14, zasada 3).

---

*Koniec dokumentu. Koncepcja platformy — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
