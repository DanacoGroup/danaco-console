# Moduł Roundtable — dokumentacja projektowa

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
| **Tytuł** | Moduł Roundtable — pełnozakresowa dokumentacja projektowa |
| **Przeznaczenie dokumentu** | Źródło wykonawcze dla Designera (co, gdzie, w jakiej formie) i Dewelopera (co zbudować) |
| **Środowiska dostępności** | TalkIn, WorkSpace, CodeStudio |
| **Forma udostępnienia** | Okno modułowe w bocznej nawigacji środowisk; moduł nie tworzy komponentu własnego |
| **Data opracowania** | 2026-08-06 |
| **Źródła** | Koncepcja platformy (rozdz. 2, 4, 9.8, 11, 12, Załącznik A) · Specyfikacja modułów 1.1 (rozdz. 4.8) · Specyfikacja okien operacyjnych 1.1 (rozdz. 5, 6.8, 7, 10) · System wizualny (rozdz. 6, 8, 10) · Integracja modeli 1.1 (rozdz. 5, Załącznik B.3) · Model danych 1.1 (rozdz. 8, 11, 17) · Model konfiguracji 1.0 (rozdz. 5, 6, 7) · Izolacja i zależności (rozdz. 1–6) |
| **Zasada nadrzędna** | Pełna kompozycyjność i pełna konfigurowalność (Koncepcja, rozdz. 14); zero blokad, klucze jawne, personalizacja z okna konfiguracji; domyślne zachowanie modułu = wykonanie |

---

## Spis treści

1. [Przeznaczenie i kontekst](#1-przeznaczenie-i-kontekst)
2. [Katalog funkcji i narzędzi](#2-katalog-funkcji-i-narzędzi)
3. [Komplet okien operacyjnych modułu](#3-komplet-okien-operacyjnych-modułu)
4. [Specyfikacja okien operacyjnych](#4-specyfikacja-okien-operacyjnych)
5. [Przepływy pracy](#5-przepływy-pracy)
6. [Punkty sterowania z okna konfiguracji](#6-punkty-sterowania-z-okna-konfiguracji)
7. [Stany, dane i powiązania](#7-stany-dane-i-powiązania)
8. [Scenariusze użycia](#8-scenariusze-użycia)
9. [Załącznik — skróty klawiszowe i ikonografia](#9-załącznik--skróty-klawiszowe-i-ikonografia)

---

## 1. Przeznaczenie i kontekst

### 1.1. Definicja i rola modułu

Roundtable jest modułem współpracy wielu modeli AI nad wspólnym problemem — przestrzenią, w której wartość nie wynika z pojedynczej, izolowanej odpowiedzi jednego modelu, lecz z konfrontacji różnych perspektyw, wzajemnej krytyki i wypracowania wspólnego stanowiska. Kilka modeli odpowiada równolegle na to samo zagadnienie, następnie odnosi się nawzajem do swoich odpowiedzi pod kierunkiem użytkownika, aż do wypracowania końcowego, uzgodnionego wniosku.

Roundtable jest **salą obrad wielomodelową** platformy Danaco Console: jednym miejscem, w którym to samo zagadnienie zostaje postawione przed kilkoma modelami i personami jednocześnie, ich odpowiedzi są zestawiane, konfrontowane w turach debaty, oceniane i głosowane, a rozproszone stanowiska syntezowane w jeden, udokumentowany wniosek z zachowanym zdaniem odrębnym. Moduł łączy w jednej przestrzeni porównywanie wielomodelowe, edycję mapy argumentów, ocenę modelami-sędziami, głosowania i rankingi, decyzję wielokryterialną oraz redakcję protokołu ustaleń.

Debata prowadzona jest dwoma kanałami komunikacji operacyjnej. Chat Window jest głównym oknem komunikacji Użytkownik ↔ Wykonawca — centralnym punktem pracy i podstawowym mechanizmem sterowania wszystkimi procesami modułu. Execution Loop Window jest oknem pętli wykonawczej Koordynator ↔ Wykonawca, w którym Koordynator dekomponuje zlecenie debaty na zadania i koordynuje wielu wykonawców uczestniczących w naradzie.

### 1.2. Dla kogo

| Grupa użytkowników | Typowa potrzeba w module Roundtable |
|---|---|
| Osoby podejmujące decyzje strategiczne | Konfrontacja argumentów za i przeciw przed podjęciem decyzji |
| Zespoły oceniające ryzyko | Wielostronna ocena tego samego zagadnienia z różnych perspektyw eksperckich |
| Redaktorzy i badacze | Weryfikacja wniosków badawczych przez kilka niezależnych „opinii” modeli |
| Zespoły programistyczne (CodeStudio) | Ocena architektury lub decyzji technicznej z kilku punktów widzenia jednocześnie |
| Facylitatorzy warsztatów i debat | Symulacja debaty przed rzeczywistą dyskusją zespołową |

### 1.3. Po co — wartość modułu

| Problem pojedynczej odpowiedzi modelu | Rozwiązanie w Roundtable |
|---|---|
| Jedna odpowiedź nie ujawnia swoich słabych punktów | Model Panels prezentują kilka niezależnych odpowiedzi obok siebie |
| Brak śladu, jak modele odniosły się do swoich argumentów | Debate Panel rejestruje pełną wymianę argumentów |
| Dyskusja modeli może zboczyć od pierwotnego pytania | Moderator Panel daje użytkownikowi pełną kontrolę nad kierunkiem debaty |
| Brak nadzoru nad przebiegiem pracy wielu wykonawców | Execution Loop Window prowadzi pętlę wykonawczą i kolejkę zadań narady |
| Struktura argumentacji ginie w zapisie liniowym | Argument Map & Analysis buduje graf tez, argumentów i kontrargumentów |
| Ocena odpowiedzi pozostaje nieformalna | Voting & Evaluation Center prowadzi głosowania, rubryki i rankingi |
| Rozproszone wnioski trudno przełożyć na decyzję | Consensus Panel gromadzi finalne, uzgodnione stanowisko w jednym miejscu |

### 1.4. Zakres tematyczny modułu (granica wewnętrzna)

| Obszar tematyczny | Zakres w module Roundtable |
|---|---|
| Porównanie wielomodelowe | Równoległe odpowiedzi wielu modeli i person na to samo pytanie, zestawienie obok siebie, tryb ślepy |
| Debata i wymiana argumentów | Formaty debaty (swobodna, strukturalna, oksfordzka, tura okrężna, Delphi), przesłuchania krzyżowe, riposty, ustępstwa |
| Persony i role uczestników | Wiele tożsamości na jednym kanale modelu, prompty ról (adwokat diabła, głos ostrożności, ekspert), pule ról |
| Koordynacja wykonawców | Dekompozycja zlecenia debaty na zadania, kolejka i stan zadań, nadzór nad realizacją przez wielu wykonawców |
| Mapowanie i analiza argumentów | Graf teza→argument→kontrargument→riposta, wydobycie argumentów, wykrywanie błędów logicznych i chwytów erystycznych |
| Ocena i głosowanie | Głosowania (aprobata, ranking IRV, Condorcet/Schulze, skala punktowa), rubryki, ocena modelami-sędziami, ranking Elo/TrueSkill |
| Wykrywanie zgody i sporu | Punkty zgody i punkty sporne, macierz zgodności, grupowanie opinii, dryf stanowiska w czasie |
| Synteza i konsensus | Automatyczny szkic stanowiska, zdanie mniejszości, ważenie poparcia, decyzja wielokryterialna |
| Moderacja i facylitacja | Sterowanie turami, timer, interwencje, kolejność głosu, szablony formatu debaty |
| Wydanie wyniku | Transkrypt debaty, protokół ustaleń, zapis decyzji (ADR), przekazanie do Studio/Research/Automations |

### 1.5. Czego moduł celowo nie robi (granica zewnętrzna)

| Obszar | Dlaczego poza Roundtable | Gdzie na platformie |
|---|---|---|
| Podział ról wykonawczych, kolejkowanie i orkiestracja procesu wieloagentowego | Roundtable to debata i konsensus, nie łańcuch wykonawczy | Środowisko **MultitaskingAI** (pokrewieństwo funkcjonalne, rozdz. 13 Koncepcji) |
| Cykliczne, bezobsługowe uruchamianie debat wg harmonogramu lub zdarzenia | To orkiestracja procesu, nie sesja obrad | Moduł **Automations** (Roundtable udostępnia debatę jako krok procesu) |
| Systematyczne zbieranie i porządkowanie źródeł przed debatą | To praca badawcza | Moduł **Research** (Research ──► Roundtable jako materiał wejściowy) |
| Redakcja jednojęzycznego dokumentu ze stanowiska końcowego | To praca autorska | Moduł **Studio** (Consensus Panel ──► Studio Editor) |
| Tworzenie i konfiguracja trwałych agentów jako komponentów własnych | Uczestnik debaty to persona sesji, nie zapisany agent | Moduł **Agents** (agent zostaje wpięty jako uczestnik) |
| Trwałe magazynowanie transkryptów i protokołów | Roundtable wydaje artefakt, nie jest repozytorium | Moduł **Library** |

Zasada rozgraniczenia: **Roundtable odpowiada za konfrontację stanowisk i wypracowanie z nich wspólnego wniosku; sąsiednie moduły odpowiadają za pochodzenie materiału, wykonanie procesu i dalszy obieg wyniku.** Powiązania są jawne i konfigurowalne (rozdz. 6), nigdy wbudowane na stałe.

### 1.6. Miejsce w architekturze platformy

```
STRONA GŁÓWNA (Centrum dowodzenia)
        │  wybór środowiska: TalkIn / WorkSpace / CodeStudio
        ▼
ŚRODOWISKO ── boczna nawigacja modułów
        │
        ▼
MODUŁ: ROUNDTABLE ───────────────────────────────────────────
        │  zestaw okien operacyjnych właściwy modułowi
        ▼
  Chat Window · Execution Loop Window · Model Panels (N) ·
  Debate Panel · Argument Map & Analysis ·
  Voting & Evaluation Center · Moderator Panel · Consensus Panel
        │
        ▼
KARTA SESJI — skład uczestników, przebieg debaty, stanowisko końcowe
  (współdzielenie z innymi kartami konfigurowalne — rozdz. 7.3)

  Pokrewieństwo funkcjonalne (nie łączenie kontekstu):
  Roundtable ═══ środowisko MultitaskingAI (rozdz. 13 Koncepcji) —
  ta sama idea wielomodelowości, w MultitaskingAI rozwinięta
  o podział ról wykonawczych i orkiestrację procesu
```

### 1.7. Dostępność i forma udostępnienia

| Wymiar | Wartość |
|---|---|
| Środowiska, w których moduł jest widoczny w bocznej nawigacji | TalkIn, WorkSpace, CodeStudio — jedyny z siedmiu opisywanych modułów dostępny we wszystkich trzech środowiskach modułowych |
| Komponent własny | Roundtable nie tworzy komponentu własnego — jest wyłącznie oknem modułowym |
| Liczba okien operacyjnych | 8 (łącznie z Chat Window i Execution Loop Window); liczba instancji Model Panels zależna od liczby wybranych uczestników |
| Charakter pracy | Wielomodelowa, turowa — debata przebiega w rundach kierowanych przez użytkownika i prowadzonych przez Koordynatora |

---

## 2. Katalog funkcji i narzędzi

Każda pozycja: **nazwa · co robi · zależności (biblioteki / formaty / integracje)**. Wszystkie funkcje działają domyślnie (zero blokad); klucze dostępowe kanałów modeli zewnętrznych są jawne i wprowadzane w oknie konfiguracji (zakres 5.4 „Zachowanie modeli”, 5.7 „Rozszerzenia”). Odpowiedzi strumieniowane są kanałem WebSocket platformy (Architektura rozdz. 11).

### 2.1. Skład panelu, persony i role uczestników

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.1.1 | **Participant Composer** | Dodaje uczestnika debaty jako parę „kanał modelu + persona” (nazwa, awatar, opis roli, prompt systemowy) | Kanał modelu (Architektura rozdz. 9); model tożsamości (5.5) |
| 2.1.2 | **Multiple Identities per Channel** | Ten sam kanał modelu w dwóch–N tożsamościach o odrębnych promptach (np. „Zwolennik” i „Oponent”) | Integracja modeli, Załącznik B.3 (jeden `kanal_modelu`, wielu uczestników) |
| 2.1.3 | **Role Library** | Biblioteka gotowych ról debaty: adwokat diabła, głos ostrożności, ekspert techniczny, sceptyk, optymista, arbiter, moderator | Presety promptów (5.6); zapis jako profil |
| 2.1.4 | **Persona Card & Profile** | Pełny profil uczestnika: prompt, historia wypowiedzi, styl, przypisany model, statystyki udziału | Model danych `kanal_modelu`/persona; widok profilu |
| 2.1.5 | **Agent-as-Participant** | Wpięcie zapisanego agenta (komponent własny z modułu Agents) jako uczestnika z jego skillami i pamięcią | Integracja Agents ──► Roundtable (5.8); Agent Manager |
| 2.1.6 | **Panel Presets / Teams** | Zapis całego składu (uczestnicy + role + format) jako nazwany zespół do ponownego użycia | Mechanizm profili (Model konfiguracji rozdz. 7); „Zespół” |
| 2.1.7 | **Self-Consistency Ensemble** | Ten sam model odpowiada N razy przy podwyższonej losowości, warianty zestawione jako głosy jednej persony | Parametry wywołania modelu (temperatura, próbki); agregacja |
| 2.1.8 | **Weighting by Expertise** | Nadanie uczestnikowi wagi kompetencji w danym temacie, wpływającej na głosowanie i syntezę | Model wag; wejście do głosowań (2.5) i syntezy (2.7) |

### 2.2. Równoległa odpowiedź i porównanie wielomodelowe

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.2.1 | **Parallel Broadcast** | Kieruje jedno pytanie jednocześnie do wszystkich uczestników i strumieniuje ich odpowiedzi równolegle | Broadcast po WebSocket; Chat Window jako źródło polecenia |
| 2.2.2 | **Side-by-Side Compare** | Zestawia odpowiedzi w siatce z podświetleniem różnic, wspólnych tez i unikalnych punktów każdego modelu | Diff tekstowy (`sergi/go-diff`); dopasowanie zdań |
| 2.2.3 | **Blind Mode** | Ukrywa tożsamość modeli do czasu głosowania, ograniczając uprzedzenie wobec pochodzenia modelu | Warstwa anonimizacji etykiet; ujawnienie po głosowaniu |
| 2.2.4 | **Response Metrics** | Pokazuje długość, czas odpowiedzi, koszt tokenów i szacowaną pewność każdego uczestnika | Licznik tokenów; telemetria kanału; kalibracja pewności |
| 2.2.5 | **Follow-up to One Panel** | Zadaje pytanie doprecyzowujące jednemu uczestnikowi poza turą ogólną, bez zaburzania debaty | Adresowanie pojedynczego panelu; wątek boczny |
| 2.2.6 | **Regenerate & Branch** | Regeneruje odpowiedź uczestnika lub rozgałęzia całą rundę jako wariant „co, gdyby” | Model wersji rundy; rozgałęzianie wątku (Chat Window) |
| 2.2.7 | **Cross-Model Fact-Check** | Krzyżowa weryfikacja twierdzeń — każdy model ocenia faktyczność twierdzeń pozostałych | Ekstrakcja twierdzeń (2.4.2); adresowanie krzyżowe |

### 2.3. Debata, formaty i orkiestracja tur

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.3.1 | **Debate Formats** | Zestaw formatów: swobodna wymiana, strukturalna (teza→riposta→wniosek), oksfordzka (za/przeciw), tura okrężna, panel ekspercki | Silnik tur; szablony formatu (Moderator Panel) |
| 2.3.2 | **Delphi Rounds** | Iteracyjne, anonimowe rundy z podsumowaniem po każdej — uczestnicy zbliżają stanowiska do konwergencji | Anonimizacja (2.2.3); podsumowanie rundy; wykrywanie konwergencji (2.6) |
| 2.3.3 | **Cross-Examination** | Rundy przesłuchania: jeden uczestnik zadaje pytania drugiemu, reszta obserwuje | Adresowanie parami; rejestr powiązań (2.4.1) |
| 2.3.4 | **Turn Orchestrator** | Wysyła stan debaty jako kontekst kolejnej tury, egzekwuje kolejność głosu, limit tur i długości wypowiedzi | Pętla wykonawcza Koordynatora; budowa kontekstu tury; kolejność głosu |
| 2.3.5 | **Devil's Advocate Injection** | Wprowadza kontrgłos, gdy panel zbytnio się zgadza (przeciwdziałanie myśleniu grupowemu) | Detekcja zbieżności (2.6.1); persona kontry (2.1.3) |
| 2.3.6 | **Side Threads** | Otwiera wątek poboczny (pytanie dodatkowe) bez przerywania głównej debaty, scalany później | Gałąź w Debate Panel; scalanie wątków |
| 2.3.7 | **Jury Mode** | Podział ról: część modeli debatuje, część pełni rolę jury oceniającego argumenty | Role (2.1.3); ocena i rubryki (2.5); rozdział uprawnień |
| 2.3.8 | **Position Drift Tracking** | Śledzi, jak stanowisko uczestnika zmieniało się między turami (od–do, powód zmiany) | Wersjonowanie wypowiedzi; porównanie tura-do-tury |

### 2.4. Mapowanie i analiza argumentów

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.4.1 | **Argument Map (graf)** | Interaktywny graf: teza → argument wspierający → kontrargument → riposta, z wcięciami wg relacji odniesienia | Model grafu (`gonum/graph`); układ sił; render SVG/canvas; eksport DOT (`emicklei/dot`) |
| 2.4.2 | **Argument Mining** | Wydobywa z wypowiedzi jednostki argumentacyjne (przesłanka, wniosek) i buduje z nich węzły grafu | Analiza modelem; klasyfikacja wypowiedzi |
| 2.4.3 | **Speech Act Classification** | Klasyfikuje każdą wypowiedź: teza, argument, kontrargument, pytanie, ustępstwo, riposta | Model klasyfikacji; etykiety Debate Panel (rozdz. 4.4) |
| 2.4.4 | **Fallacy Detection** | Wykrywa błędy logiczne i chwyty erystyczne (ad hominem, słomiany kukła, fałszywa alternatywa) | Katalog błędów; analiza modelem; oznaczenie w grafie |
| 2.4.5 | **Steelman / Strawman Check** | Ocenia, czy kontrargument mierzy się z najmocniejszą wersją tezy, czy z jej osłabieniem | Analiza modelem; powiązanie węzłów grafu |
| 2.4.6 | **Evidence & Citation Ledger** | Rejestr dowodów i źródeł przywołanych w argumentach, z oznaczeniem niepopartych twierdzeń | Ekstrakcja cytowań; powiązanie z Research/Browser (5.8) |
| 2.4.7 | **AIF / Argdown Export** | Eksport mapy argumentów w standardach wymiany (Argument Interchange Format, Argdown, GraphML) | JSON AIF; składnia Argdown; GraphML (`encoding/xml`) |
| 2.4.8 | **Key-Argument Pinning** | Oznaczenie fragmentu jako kluczowego argumentu, przekazywanego do syntezy i grafu | Pula argumentów kluczowych; wejście do Consensus (2.7) |

### 2.5. Ocena, głosowanie i ranking

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.5.1 | **Vote Engine (metody)** | Głosowania nad odpowiedziami i stanowiskami wieloma metodami: aprobata, ranking (IRV), Condorcet/Schulze, skala punktowa, kwadratowe | Implementacja metod (własna: IRV, Schulze); agregacja głosów |
| 2.5.2 | **Star & Preference Rating** | Ocena „najbardziej przekonujące” gwiazdkami lub porównaniem parami przez użytkownika | Model ocen; ranking preferencji |
| 2.5.3 | **LLM-as-Judge** | Modele-sędziowie oceniają odpowiedzi wg rubryki (trafność, spójność, dowody), z uzasadnieniem | Prompt sędziego; rubryka (2.5.4); rola jury (2.3.7) |
| 2.5.4 | **Rubric Builder** | Definiuje kryteria oceny z wagami (np. rzetelność 40%, wykonalność 30%, ryzyko 30%) | Model rubryki; edytor kryteriów i wag |
| 2.5.5 | **Model Elo / TrueSkill Leaderboard** | Ranking modeli i person na podstawie wyników pojedynków w debatach i głosowaniach, akumulowany między sesjami | Implementacja Elo/Glicko/TrueSkill; historia pojedynków |
| 2.5.6 | **Multi-Criteria Decision Matrix** | Zestawia warianty i stanowiska w macierzy kryteriów z wagami i wyliczonym wynikiem (styl AHP / ważonej sumy) | Model macierzy decyzyjnej; normalizacja i ważenie |
| 2.5.7 | **Confidence & Calibration** | Zbiera deklarowaną pewność uczestników i konfrontuje ją z trafnością (kalibracja nadmiernej pewności) | Odczyt pewności modelu; metryka kalibracji |
| 2.5.8 | **Quorum & Threshold** | Ustala próg zgody wymagany do uznania konsensusu (np. 2/3 głosów, jednomyślność) | Reguła progu; wejście do stanu „konsensus” (rozdz. 7.1) |

### 2.6. Wykrywanie zgody, sporu i rozbieżności

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.6.1 | **Agreement / Dispute Detection** | Sygnalizuje automatycznie punkty zgody i punkty sporne między uczestnikami | Podobieństwo semantyczne (embeddingi modelu); analiza stanowisk |
| 2.6.2 | **Agreement Matrix (heatmapa)** | Macierz kto-z-kim się zgadza w poszczególnych kwestiach, jako mapa cieplna | Macierz par uczestników; render heatmapy |
| 2.6.3 | **Opinion Clustering** | Grupuje zbliżone stanowiska w klastry opinii (kto tworzy „obóz”) | Embeddingi; grupowanie (k-means/HDBSCAN, `gonum`) |
| 2.6.4 | **Semantic Dedup of Arguments** | Scala powtarzające się argumenty różnych uczestników w jeden węzeł, licząc poparcie | Embeddingi; próg podobieństwa; scalanie węzłów (2.4.1) |
| 2.6.5 | **Consensus Convergence Meter** | Wskaźnik, jak bardzo panel zbliżył się do zgody po każdej turze (dynamika konwergencji) | Metryka rozrzutu stanowisk tura-do-tury |
| 2.6.6 | **Crux Finder** | Wskazuje kluczowy punkt sporny, którego rozstrzygnięcie zmieniłoby wnioski (podwójny crux) | Analiza zależności węzłów grafu; wejście do interwencji (2.8) |

### 2.7. Synteza i konsensus

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.7.1 | **Consensus Generator** | Generuje szkic uzgodnionego stanowiska na podstawie tur i argumentów kluczowych | Argumenty kluczowe (2.4.8); przebieg z Debate Panel; model |
| 2.7.2 | **Editable Position** | Stanowisko końcowe jako edytowalny punkt wyjścia, nie wynik ostateczny — z operacjami kontekstowymi AI | Operacje kontekstowe (natywne Studio, integracja 5.8) |
| 2.7.3 | **Minority Report** | Zachowuje i podpisuje zdanie odrębne uczestnika, który nie dołączył do konsensusu | Model stanowiska mniejszości; blok w eksporcie |
| 2.7.4 | **Weighted Support Tally** | Przypisuje każdemu wnioskowi poparcie (ilu i jak ważnych uczestników go poparło) | Głosowania (2.5); wagi ekspertyzy (2.1.8) |
| 2.7.5 | **Points of Agreement / Dispute Summary** | Zestawia zbiorczo punkty zgody i punkty pozostające sporne jako podsumowanie wyniku | Wykrywanie zgody/sporu (2.6.1); liczniki |
| 2.7.6 | **Position Versioning** | Zapisuje kolejne redakcje stanowiska po kolejnych turach jako odrębne wersje z porównaniem | `artefakt`/`wersja_artefaktu`; diff (`sergi/go-diff`) |
| 2.7.7 | **Decision Record (ADR)** | Formatuje wynik jako zapis decyzji: kontekst, warianty, decyzja, konsekwencje, ryzyko mniejszości | Szablon ADR; render Markdown (`yuin/goldmark`) |

### 2.8. Moderacja i facylitacja

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.8.1 | **Moderator Interventions** | Polecenie do wszystkich uczestników naraz („skupcie się na finansach”, „unikajcie powtórzeń”) wstrzykiwane do kontekstu | Wstrzyknięcie do kontekstu tury; broadcast |
| 2.8.2 | **Turn Timer & Limits** | Timer tury z paskiem postępu, limit czasu i długości wypowiedzi, sugerowane zamknięcie rundy | Zegar tury; reguły limitu (nie blokują) |
| 2.8.3 | **Speaking Order Control** | Ustala i przeciąga kolejność głosu, przełącza między trybem równoległym a sekwencyjnym | Lista przeciągalna; tryb prezentacji (Model Panels) |
| 2.8.4 | **Manual Turn Close** | Ręczne zamknięcie tury przed zakończeniem wypowiedzi wszystkich uczestników | Stan tury; przejście dalej mimo niepełnej rundy |
| 2.8.5 | **Moderation Templates** | Zapis formatu, liczby tur i kolejności głosu jako nazwanego szablonu moderacji | Profil moderacji (Model konfiguracji rozdz. 7) |
| 2.8.6 | **Bias & Tone Guardrails** | Reguły miękkie: sygnalizacja agresji, offtopu, dominacji jednego uczestnika w czasie mówienia | Analiza modelem; udział czasowy per uczestnik |

### 2.9. Pętla wykonawcza i koordynacja wykonawców

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.9.1 | **Debate Brief Decomposition** | Koordynator rozkłada zlecenie debaty na zadania jednostkowe: wypowiedź uczestnika w turze, riposta, analiza argumentów, ocena jury, synteza | Model zlecenia i zadania; Execution Loop Window |
| 2.9.2 | **Multi-Executor Dispatch** | Rozdziela zadania jednej tury równolegle na wielu wykonawców uczestniczących w naradzie i śledzi ich stan | Kolejka zadań; kanały modeli uczestników |
| 2.9.3 | **Executor Health & Retry** | Wykrywa wykonawcę, który nie odpowiedział lub zwrócił wynik odrzucony w kontroli jakości, i ponawia zadanie | Reguły ponowienia; telemetria kanału |
| 2.9.4 | **Loop Quality Gates** | Kontrola jakości wyniku zadania przed dopuszczeniem go do zapisu tury (kompletność, zgodność z rolą, format wypowiedzi) | Reguły kontroli; klasyfikacja wypowiedzi (2.4.3) |
| 2.9.5 | **Loop Control** | Wstrzymanie, wznowienie, przerwanie i korekta zlecenia debaty w toku pętli | Sterowanie przebiegiem; stan tury |
| 2.9.6 | **Coordinator Message Log** | Rejestr komunikatów sterujących Koordynator ↔ Wykonawca dla całej narady | Model wiadomości sterujących; zapis per karta sesji |

### 2.10. Wydanie, raport i integracje

| # | Funkcja | Co robi | Zależności |
|---|---|---|---|
| 2.10.1 | **Transcript Export** | Eksport pełnego zapisu debaty (uczestnik, tura, klasyfikacja) do PDF/DOCX/Markdown/JSON | `unidoc/unioffice` (DOCX), `pdfcpu/pdfcpu`+`go-fitz` (PDF), `goldmark` (MD) |
| 2.10.2 | **Argument Map Export** | Eksport grafu argumentów jako SVG/PNG oraz DOT/GraphML/Argdown | Render SVG; `emicklei/dot`; GraphML (`encoding/xml`) |
| 2.10.3 | **Voting/Scorecard Report** | Raport z głosowań, rubryk i rankingu Elo z wykresami wyników | Dane głosowań (2.5); render wykresów |
| 2.10.4 | **TTS Debate Readback** | Odsłuch przebiegu debaty syntezą mowy z głosem per uczestnik | API TTS (Azure/Google/ElevenLabs/OpenAI); głos per persona (klucze jawne) |
| 2.10.5 | **Send to Studio** | Przekazuje stanowisko końcowe do Studio Editor do dalszej redakcji | Integracja Consensus Panel ──► Studio (5.8) |
| 2.10.6 | **Send to Research** | Przekazuje stanowisko jako ustalenie badawcze do Findings Panel / Report Builder | Integracja Roundtable ──► Research (5.8) |
| 2.10.7 | **Research Intake** | Przyjmuje wstępną syntezę z Research jako pytanie wyjściowe debaty | Integracja Research ──► Roundtable (5.8) |
| 2.10.8 | **Automations Step** | Udostępnia „przeprowadź debatę N modeli i zwróć konsensus” jako krok procesu wsadowego | Kolejka zadań (5.3); tryb w tle |
| 2.10.9 | **Roundtable API / MCP** | Wystawia debatę i głosowanie jako operacje wywoływane z zewnątrz (rozszerzenie, konektor, serwer MCP) | Kontrakt rozszerzeń (Architektura rozdz. 14); MCP/konektor |
| 2.10.10 | **External Model Connectors** | Podłącza modele spoza platformy jako uczestników, z jawnymi kluczami | Kanał modelu HTTP/API (Architektura rozdz. 9); klucze z konfiguracji |

---

## 3. Komplet okien operacyjnych modułu

| # | Okno | Typologia wizualna | Waga wizualna w module | Warstwa | Sposób wywołania | Rola w module |
|---|---|---|---|---|---|---|
| 1 | Chat Window | Komunikacja (kanał Użytkownik ↔ Wykonawca) | Lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Widoczne bez interakcji | Główne okno komunikacji, centralny punkt pracy i podstawowy mechanizm sterowania procesami modułu |
| 2 | Execution Loop Window | Komunikacja (kanał Koordynator ↔ Wykonawca) | Kolumna sąsiadująca z Chat Window, otwierana | 1 | Widoczne bez interakcji po uruchomieniu tury; przełącznik `Pętla ▼` przy braku zlecenia | Pętla wykonawcza narady — dekompozycja zlecenia, kolejka zadań wielu wykonawców, kontrola realizacji |
| 3 | Model Panels | Okno edycyjne / komunikacja (instancja wielokrotna) | Prawa kolumna dominująca, siatka równoległa | 1 | Widoczne bez interakcji | Równoległe odpowiedzi uczestniczących modeli |
| 4 | Debate Panel | Okno monitorów i dashboardów | Prawa kolumna dominująca, widok wiodący | 1 | Widoczne bez interakcji po pierwszej turze | Rejestr wymiany argumentów |
| 5 | Argument Map & Analysis | Okno monitorów i dashboardów | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Przycisk `Mapa argumentów` w Debate Panel; polecenie języka naturalnego w Chat Window | Struktura argumentacji, analiza rzetelności i zgodności |
| 6 | Voting & Evaluation Center | Okno monitorów i dashboardów | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Przycisk `Głosowanie` w Debate Panel; polecenie języka naturalnego w Chat Window | Głosowania, rubryki, ocena modelami-sędziami, rankingi |
| 7 | Moderator Panel | Panel narzędziowy | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Znacznik kontekstowy `[Format: Strukturalna]` w pasku kontekstu; przycisk `Moderacja` | Ukierunkowanie przebiegu dyskusji |
| 8 | Consensus Panel | Okno edycyjne / wynikowe | Kolumna boczna, otwierana jako rozszerzenie boczne | 2 | Przycisk `Konsensus` po dowolnej turze; polecenie języka naturalnego w Chat Window | Finalne, uzgodnione stanowisko |

```
 Makieta zbiorcza — Moduł Roundtable      Dostępność: TalkIn, WorkSpace, CodeStudio
 ═══════════════════════════════════════════════════════════════════════════════
  Boczna     │ Chat Window          │ Execution Loop      │ Obszar roboczy modułu
  nawigacja  │ Użytkownik ↔         │ Koordynator ↔       │ ┌────────┬────────┬────────┐
  modułów    │ Wykonawca            │ Wykonawca           │ │ Model  │ Model  │ Model  │
  (poza      │                      │                     │ │ Panel A│ Panel B│ Panel C│
  zakresem   │ [Roundtable]         │ Zlecenie: debata    │ └────────┴────────┴────────┘
  dokumentu) │ [3 uczestników ▼]    │ Tura 2 · 3 zadania  │ ┌──────────────────────────┐
             │ [Format: Struktur.]  │ ▸ Model A  gotowe   │ │      Debate Panel        │
             │                  ⋮   │ ▸ Model B  w toku   │ │ wymiana argumentów       │
             │                      │ ▸ Zwolennik kolejka │ │                       ⋮  │
             │ …strumień            │                  ⋮  │ └──────────────────────────┘
             │  odpowiedzi…         │                     │
             │                      │                     │
             │ [📎] Pole poleceń ▶  │ [ Sterowanie ▼ ]    │
 ═══════════════════════════════════════════════════════════════════════════════
   Kolumny boczne (rozszerzenia boczne, otwierane): Argument Map & Analysis ·
   Voting & Evaluation Center · Moderator Panel · Consensus Panel
```

### 3.1. Warstwy widoczności w module

Interfejs modułu Roundtable ujawnia funkcje stopniowo: w stanie spoczynku widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 — znaczniki kontekstowe w pasku kontekstu, przełączniki `▼`, menu kebab `⋮` i menu `☰`.

| Warstwa | Zawartość w module Roundtable | Sposób dostępu |
|---|---|---|
| 1 | Chat Window, Execution Loop Window, Model Panels, Debate Panel, pasek kontekstu debaty, wskaźniki stanu tury i zadań | Widoczne bez interakcji |
| 2 | Argument Map & Analysis, Voting & Evaluation Center, Moderator Panel, Consensus Panel, wybór uczestników, wybór formatu debaty, wybór modelu i persony | Znacznik kontekstowy, przycisk, przełącznik; po użyciu element zwija się samoczynnie |
| 3 | Zestawy akcji tury, ustawienia szybkie panelu uczestnika, warianty eksportu, filtry zapisu debaty | Menu kebab `⋮`, menu `☰`, menu kontekstowe, panel popover, lista rozwijana |
| 4 | Rubric Builder, konfiguracja rankingu Elo, katalog błędów logicznych, reguły ponowienia zadań pętli, macierz izolacji per uczestnik, diagnostyka kanałów modeli | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

Każda ukryta funkcja modułu osiągalna jest jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego w Chat Window.

---

## 4. Specyfikacja okien operacyjnych

### 4.1. Chat Window (kanał Użytkownik ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja |
| Waga wizualna | Lewa kolumna, stała, pełna wysokość obszaru roboczego |
| Warstwa | 1 |
| Izolacja domyślna | Odrębna historia i pamięć per karta sesji; współdzielenie konfigurowalne (rozdz. 7.3) |

Chat Window jest głównym oknem komunikacji Użytkownika z Wykonawcą, centralnym punktem pracy w module i podstawowym mechanizmem sterowania wszystkimi procesami debaty: przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu.

**Zawartość i pełny arsenał funkcji.**

- Sformułowanie pytania wyjściowego debaty, kierowanego jednocześnie do wszystkich uczestniczących modeli.
- Strumień odpowiedzi modeli na żywo kanałem WebSocket, równolegle w Model Panels.
- Polecenia meta-poziomu wobec całej debaty: „poproś model B o ustosunkowanie się do argumentu modelu A”, „podsumuj różnice stanowisk”.
- Sterowanie procesami modułu poleceniem języka naturalnego: uruchomienie tury, otwarcie mapy argumentów, uruchomienie głosowania, wygenerowanie konsensusu, wstrzymanie i przerwanie pętli wykonawczej.
- Odwołania do konkretnego modelu lub tury debaty przez wzmiankę w treści polecenia.
- Wstawienie fragmentu dyskusji jako materiału do Consensus Panel.
- Historia poleceń, regeneracja rundy, rozgałęzianie wątku dyskusji.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Chat Window ──────────────────┐
│ [Roundtable] [3 uczestników ▼] │
│ [Format: Strukturalna]      ⋮  │
├────────────────────────────────┤
│ ┌─ Użytkownik ───────────────┐ │
│ │ „czy powinniśmy wejść na   │ │
│ │  rynek azjatycki?”         │ │
│ └────────────────────────────┘ │
│ ┌─ System ───────────────────┐ │
│ │ pytanie przekazane         │ │
│ │ do 3 uczestników        ⋮  │ │
│ └────────────────────────────┘ │
├────────────────────────────────┤
│ [📎] Pole poleceń …… [ Wyślij ▶]│
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Znacznik kontekstowy modułu | Lekki znacznik `[Roundtable]` | Wskazanie aktywnego kontekstu pracy | Mały znacznik w pasku kontekstu | 1 | Widoczny bez interakcji | domyślny | Kliknięcie otwiera selektor modułu | Pasek kontekstu, góra Chat Window |
| Wskaźnik liczby uczestników | Zwinięty przełącznik `[3 uczestników ▼]` | Skrót do zarządzania składem modeli | Mały przełącznik z licznikiem | 2 | Kliknięcie znacznika | zwinięty · rozwinięty | Otwiera selektor uczestników; po wyborze zwija się samoczynnie | Pasek kontekstu |
| Znacznik formatu debaty | Lekki znacznik `[Format: Strukturalna]` | Wskazanie obowiązującego formatu | Mały znacznik | 2 | Kliknięcie znacznika | zwinięty · rozwinięty | Otwiera Moderator Panel jako rozszerzenie boczne | Pasek kontekstu |
| Menu operacji debaty | Menu kebab `⋮` | Zestaw akcji tury i debaty | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę akcji: uruchom turę, otwórz mapę, głosowanie, konsensus, eksport | Nagłówek Chat Window i karty systemowe |
| Pole poleceń | Wieloliniowe pole tekstowe | Wprowadzenie pytania wyjściowego, polecenia meta-poziomu i polecenia sterującego | Duże pole, koniec kolumny | 1 | Widoczne bez interakcji | domyślny · fokus | Enter wysyła pytanie jednocześnie do wszystkich uczestniczących modeli | Chat Window |
| Karta rozgłoszenia pytania | Systemowy komunikat w historii | Potwierdzenie przekazania pytania do wszystkich modeli | Mały blok informacyjny | 1 | Widoczna bez interakcji | domyślny | Pozycja menu `⋮` „Otwórz Debate Panel” przenosi fokus | Historia rozmowy |
| Funkcje eksperckie modułu | Polecenia języka naturalnego | Dostęp do rubryk, rankingu Elo, reguł ponowienia, diagnostyki kanałów | Bez reprezentacji graficznej | 4 | Polecenie w polu poleceń, skrót klawiszowy, wyszukiwarka funkcji | — | Uruchamia funkcję i prezentuje wynik w oknie właściwym | Chat Window |

---

### 4.2. Execution Loop Window (kanał Koordynator ↔ Wykonawca)

| Aspekt | Wartość |
|---|---|
| Typologia | Komunikacja |
| Waga wizualna | Kolumna sąsiadująca z Chat Window, otwierana, pełna wysokość obszaru roboczego |
| Warstwa | 1 |
| Izolacja domyślna | Pętla prowadzona w obrębie karty sesji; rejestr komunikatów sterujących zapisywany per karta |

Execution Loop Window prezentuje komunikację między Koordynatorem a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. W module Roundtable pętla wykonawcza obejmuje koordynację **wielu wykonawców uczestniczących w naradzie** — każdy uczestnik debaty jest odrębnym wykonawcą, a Koordynator prowadzi ich równolegle w obrębie jednej tury.

**Mechanizm koordynacji wielu wykonawców narady.**

1. **Przyjęcie zlecenia.** Pytanie wyjściowe lub polecenie tury z Chat Window trafia do Koordynatora jako zlecenie debaty wraz z kontekstem: skład uczestników, format debaty, kolejność głosu, limity tury.
2. **Dekompozycja.** Koordynator rozkłada zlecenie na zadania jednostkowe — po jednym zadaniu „wypowiedź w turze” na każdego uczestnika, uzupełnione o zadania analityczne (klasyfikacja wypowiedzi, wydobycie argumentów) i zadania jury, gdy format przewiduje ocenę.
3. **Przydział i kolejkowanie.** Zadania trafiają do kolejki i są przydzielane wykonawcom zgodnie z trybem tury: równolegle w trybie broadcastu albo sekwencyjnie w turze okrężnej, z zachowaniem ustalonej kolejności głosu.
4. **Nadzór nad realizacją.** Koordynator śledzi stan każdego zadania (kolejka, w toku, gotowe, odrzucone, ponowione) i udział czasowy każdego wykonawcy; wykonawca, który nie odpowiedział w limicie tury, otrzymuje ponowienie albo zostaje pominięty przy ręcznym zamknięciu tury.
5. **Kontrola jakości.** Wynik zadania przechodzi kontrolę kompletności, zgodności z rolą uczestnika i formatu wypowiedzi przed dopuszczeniem do zapisu tury; wynik odrzucony wraca do wykonawcy z komunikatem sterującym.
6. **Domknięcie tury.** Po zebraniu wyników wszystkich zadań Koordynator scala je w stan tury, przekazuje do Debate Panel i buduje kontekst tury kolejnej.
7. **Sterowanie przebiegiem.** Użytkownik wstrzymuje, wznawia, przerywa i koryguje zlecenie w toku — z Execution Loop Window albo poleceniem języka naturalnego w Chat Window.

**Zawartość i pełny arsenał funkcji.**

- Bieżące zlecenie debaty i jego dekompozycja na zadania.
- Kolejka i stan zadań wszystkich wykonawców uczestniczących w naradzie.
- Wymiana komunikatów sterujących Koordynator ↔ Wykonawca, z pełnym rejestrem.
- Wyniki kontroli jakości i decyzje o ponowieniu zadania.
- Wskaźniki przebiegu pętli: postęp tury, liczba zadań gotowych i oczekujących, czas realizacji per wykonawca.
- Sterowanie przebiegiem: wstrzymanie, wznowienie, przerwanie, korekta zlecenia.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Execution Loop ───────────────┐
│ Koordynator ↔ Wykonawca     ⋮  │
├────────────────────────────────┤
│ Zlecenie: debata — tura 2      │
│ Format: strukturalna           │
├────────────────────────────────┤
│ Zadania tury          3/5      │
│ ▸ Model A      gotowe          │
│ ▸ Model B      w toku          │
│ ▸ Zwolennik    kolejka         │
│ ▸ klasyfikacja gotowe          │
│ ▸ jury         kolejka         │
├────────────────────────────────┤
│ Kontrola jakości: 1 ponowienie │
├────────────────────────────────┤
│ [ Sterowanie ▼ ]               │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek zlecenia | Blok z treścią zlecenia i formatem | Orientacja w bieżącym zleceniu debaty | Mały blok, góra kolumny | 1 | Widoczny bez interakcji | domyślny · zlecenie skorygowane | Kliknięcie rozwija pełną treść zlecenia | Góra okna |
| Lista zadań tury | Lista pozycji z etykietą stanu | Nadzór nad realizacją zadań wielu wykonawców | Średnia lista, ciało okna | 1 | Widoczna bez interakcji | kolejka · w toku · gotowe · odrzucone · ponowione | Kliknięcie pozycji otwiera komunikaty sterujące danego zadania | Ciało okna |
| Wskaźnik postępu tury | Licznik `3/5` z paskiem | Stan realizacji bieżącej tury | Mały element graficzny | 1 | Widoczny bez interakcji | aktualizowany na żywo | — (informacyjny) | Nagłówek listy zadań |
| Blok kontroli jakości | Mały blok `.dn-alert--ostrz` | Sygnalizacja wyników odrzuconych i ponowień | Mały blok | 1 | Widoczny przy wystąpieniu ponowienia | ukryty (brak odrzuceń) · widoczny | Kliknięcie rozwija powód odrzucenia i treść komunikatu zwrotnego | Ciało okna |
| Menu sterowania przebiegiem | Zwinięty element zbiorczy `[ Sterowanie ▼ ]` | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Element zbiorczy, koniec kolumny | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty | Rozwija pełną listę akcji sterujących; po wyborze zwija się samoczynnie | Koniec kolumny |
| Rejestr komunikatów sterujących | Menu kebab `⋮` → panel popover | Podgląd pełnej wymiany Koordynator ↔ Wykonawca | Panel popover | 3 | Kliknięcie `⋮` | zwinięty · rozwinięty | Otwiera rejestr komunikatów z filtrem per wykonawca | Nagłówek okna |
| Reguły ponowienia i limitów pętli | Ustawienia pętli wykonawczej | Definicja liczby ponowień, limitu czasu zadania, progu odrzucenia | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny, okno konfiguracji | — | Zmienia zachowanie pętli dla bieżącej karty sesji | Okno konfiguracji, zakres pętli wykonawczej |

---

### 4.3. Model Panels

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne / komunikacja (instancja wielokrotna) |
| Waga wizualna | Prawa kolumna dominująca, siatka równoległa — jeden panel na każdego uczestniczącego „mówcę” |
| Warstwa | 1 |
| Izolacja domyślna | Każdy panel niezależny w obrębie karty; wspólne pytanie wyjściowe z Chat Window |

**Zawartość i pełny arsenał funkcji.**

*Skład uczestników.*
- Dodanie i usunięcie modelu jako uczestnika debaty — wybór spośród dostępnych kanałów modeli, w tym kanałów modeli zewnętrznych z jawnymi kluczami.
- Ten sam model bazowy uczestniczy w debacie wielokrotnie pod odrębnymi tożsamościami (ten sam kanał modelu, różne persony i prompty systemowe — np. „Zwolennik” i „Oponent”), zgodnie z ustalonym wzorcem integracji modeli.
- Nadanie każdemu uczestnikowi nazwy, awatara i krótkiego opisu roli w debacie (np. „adwokat diabła”, „głos ostrożności”, „ekspert techniczny”), z biblioteki ról albo własnym opisem.
- Przypisanie promptu systemowego ukierunkowującego stanowisko lub perspektywę uczestnika.
- Wpięcie zapisanego agenta z modułu Agents jako uczestnika, z jego skillami i pamięcią.
- Zapis całego składu (uczestnicy, role, format) jako nazwanego zespołu do ponownego użycia.
- Nadanie uczestnikowi wagi kompetencji wpływającej na głosowanie i syntezę.

*Prezentacja odpowiedzi.*
- Równoległe generowanie odpowiedzi wszystkich uczestników na to samo pytanie, widoczne obok siebie.
- Strumieniowanie na żywo każdej odpowiedzi niezależnie — uczestnicy „mówią” równocześnie, nie sekwencyjnie.
- Zestawienie odpowiedzi z podświetleniem różnic, wspólnych tez i unikalnych punktów każdego uczestnika.
- Tryb ślepy — tożsamość modeli pozostaje ukryta do czasu głosowania.
- Metryki odpowiedzi: długość, czas, koszt tokenów, szacowana pewność.
- Ocena i ranking odpowiedzi przez użytkownika (gwiazdki lub głosowanie „najbardziej przekonujące”).
- Ensemble self-consistency — ta sama persona w N wariantach jako głosy do agregacji.

*Interakcja z pojedynczym panelem.*
- Zadanie pytania doprecyzowującego bezpośrednio jednemu uczestnikowi, poza turą ogólną.
- Wyciszenie uczestnika na czas rundy (bez usuwania go z debaty).
- Regeneracja odpowiedzi uczestnika oraz rozgałęzienie całej rundy jako wariantu.
- Podgląd pełnego profilu uczestnika: prompt, historia wypowiedzi, styl, statystyki udziału.
- Oznaczenie fragmentu odpowiedzi jako kluczowego argumentu (przekazywane do Debate Panel, Argument Map & Analysis i Consensus Panel).

**Makieta tekstowa (stan spoczynku).**

```
┌─ Model Panels ──────────────────────────────────────────┐
│ [ Uczestnicy ▼ ]                                     ⋮  │
├───────────────────┬──────────────────┬──────────────────┤
│ 🔵 Model A        │ 🟡 Model B       │ 🟣 „Zwolennik”   │
│ „Ekspert rynku” ⋮ │ „Głos ostroż.” ⋮ │ (Model A×2)    ⋮ │
├───────────────────┼──────────────────┼──────────────────┤
│ argument za       │ argument przeciw │ kontrargument    │
│ wejściem na       │ pochopnej        │ wobec B …        │
│ rynek …           │ decyzji …        │                  │
│                   │                  │                  │
└───────────────────┴──────────────────┴──────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Przełącznik `[ Uczestnicy ▼ ]` | Zwinięty element zbiorczy | Dodanie, usunięcie i konfiguracja uczestnika, wybór zespołu | Mały element zbiorczy | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty | Rozwija formularz wyboru kanału modelu, nazwy, awatara i promptu roli; po zapisie zwija się samoczynnie | Góra okna, nad siatką paneli |
| Nagłówek panelu uczestnika | Pasek z awatarem, nazwą i opisem roli | Identyfikacja uczestnika debaty | Mały pasek, góra każdego panelu | 1 | Widoczny bez interakcji | domyślny · wyciszony · anonimowy (tryb ślepy) | Kliknięcie otwiera profil uczestnika jako panel popover | Góra każdego Model Panel |
| Awatar uczestnika | Element graficzny `.dn-awatar` | Wizualne rozróżnienie uczestników w siatce i w Debate Panel | Mały okrągły element z kolorem lub inicjałem | 1 | Widoczny bez interakcji | domyślny · aktywny (mówi teraz, obwódka pulsująca) | — (identyfikacyjny) | Nagłówek panelu, Debate Panel |
| Treść odpowiedzi | Blok tekstu strumieniowanego | Prezentacja wypowiedzi uczestnika w bieżącej turze | Średni lub duży blok, ciało panelu | 1 | Widoczna bez interakcji | oczekiwanie · strumieniowanie · kompletna | Zaznaczenie fragmentu ujawnia menu kontekstowe z akcją „oznacz jako kluczowy” | Ciało panelu |
| Menu panelu uczestnika | Menu kebab `⋮` | Wyciszenie w turze, regeneracja, rozgałęzienie rundy, pytanie doprecyzowujące, metryki odpowiedzi | Ikona w nagłówku panelu | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę akcji dla wybranego uczestnika | Nagłówek każdego Model Panel |
| Akcja „oznacz jako kluczowy” | Pozycja menu kontekstowego | Wyróżnienie istotnego argumentu do dalszego wykorzystania | Pozycja menu kontekstowego | 3 | Zaznaczenie fragmentu odpowiedzi | nieaktywna · aktywna | Dodaje fragment do puli argumentów kluczowych widocznej w Consensus Panel i w grafie | Menu kontekstowe zaznaczenia |
| Zestawienie różnic | Podświetlenie w treści odpowiedzi | Wskazanie różnic, wspólnych tez i punktów unikalnych | Podświetlenia w blokach treści | 2 | Pozycja `Porównaj odpowiedzi` w menu `⋮` okna | wyłączone · włączone | Podświetla fragmenty rozbieżne i zbieżne w siatce paneli | Ciało paneli |
| Tryb ślepy | Ustawienie konfiguracyjne widoczności tożsamości | Ograniczenie uprzedzenia wobec pochodzenia modelu | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, okno konfiguracji, tryb administracyjny | wyłączony · włączony | Anonimizuje nagłówki paneli do czasu zakończenia głosowania | Okno konfiguracji, presety modułu |
| Ensemble self-consistency | Ustawienie konfiguracyjne wywołania modelu | Zebranie N wariantów jednej persony jako głosów do agregacji | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, okno konfiguracji | — | Uruchamia N wywołań kanału i agreguje wyniki w jednym panelu | Okno konfiguracji, zakres 5.4 |

---

### 4.4. Debate Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno monitorów i dashboardów |
| Waga wizualna | Prawa kolumna dominująca, widok wiodący |
| Warstwa | 1 |
| Izolacja domyślna | Narasta w toku debaty prowadzonej przez pętlę wykonawczą, kierowanej przez użytkownika i Moderator Panel |

**Zawartość i pełny arsenał funkcji.**

- Chronologiczny zapis kolejnych tur debaty — każda wypowiedź powiązana z uczestnikiem i turą, w której padła.
- Wizualne powiązania odniesień między wypowiedziami (która wypowiedź jest odpowiedzią lub kontrargumentem na którą).
- Klasyfikacja wypowiedzi: teza, argument wspierający, kontrargument, pytanie, ustępstwo, riposta.
- Oznaczenie przesłuchań krzyżowych i wątków pobocznych na osi debaty; śledzenie dryfu stanowiska uczestnika między turami.
- Filtrowanie zapisu: tylko wypowiedzi jednego uczestnika, tylko wypowiedzi oznaczone jako kluczowe, tylko punkty sporne.
- Wykrywanie punktów zgody i punktów spornych między uczestnikami (sygnalizacja automatyczna) oraz wskaźnik konwergencji po każdej turze.
- Eksport pełnego zapisu debaty jako transkryptu (PDF, DOCX, Markdown, JSON).
- Uruchomienie kolejnej tury bezpośrednio z okna — zlecenie trafia do Koordynatora i pętli wykonawczej.
- Otwarcie Argument Map & Analysis oraz Voting & Evaluation Center jako rozszerzeń bocznych.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Debate Panel ───────────────────────────────────────┐
│ Tura: 2 z 3                                       ⋮  │
├──────────────────────────────────────────────────────┤
│ Tura 1                                               │
│  🔵 Model A — teza: „wejście na rynek jest korzystne”│
│  🟡 Model B — kontrargument: „ryzyko regulacyjne”    │
│  🟣 Zwolennik — riposta wobec B: „ryzyko ograniczone”│
├──────────────────────────────────────────────────────┤
│ Tura 2 (w toku)                                      │
│  🟡 Model B — ustępstwo częściowe                    │
│  ⚡ punkt sporny: skala inwestycji początkowej       │
├──────────────────────────────────────────────────────┤
│ [ Uruchom turę ]  [ Mapa argumentów ]  [ Głosowanie ]│
└──────────────────────────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Wskaźnik tury | Mała plakietka z licznikiem | Orientacja w postępie debaty | Plakietka `.dn-plakietka--info` | 1 | Widoczny bez interakcji | aktualizowany na żywo | — (informacyjny) | Góra okna |
| Wpis wypowiedzi | Blok tekstu z awatarem i klasyfikacją | Pojedyncza wypowiedź w toku debaty | Mały lub średni blok, wcięty wg relacji odniesienia | 1 | Widoczny bez interakcji | teza · argument · kontrargument · ustępstwo · riposta | Kliknięcie przewija do pełnej treści w Model Panels | Ciało okna |
| Znacznik punktu spornego | Mała ikona błyskawicy z etykietą | Sygnalizacja nierozstrzygniętej kwestii między uczestnikami | Mały blok `.dn-alert--ostrz` | 1 | Widoczny przy wykryciu sporu | ukryty · widoczny | Kliknięcie otwiera Moderator Panel z sugestią doprecyzowania tego punktu | Koniec listy bieżącej tury |
| Przycisk „Uruchom turę” | Przycisk `--zloty` | Rozpoczęcie następnej rundy wymiany argumentów | Średni przycisk CTA | 1 | Widoczny bez interakcji | domyślny · ładowanie | Przekazuje stan debaty Koordynatorowi jako zlecenie kolejnej tury | Koniec kolumny okna |
| Przycisk „Mapa argumentów” | Przycisk `--zarys` | Otwarcie Argument Map & Analysis | Mały przycisk | 2 | Kliknięcie przycisku | domyślny · aktywny | Otwiera okno analizy jako rozszerzenie boczne | Koniec kolumny okna |
| Przycisk „Głosowanie” | Przycisk `--zarys` | Otwarcie Voting & Evaluation Center | Mały przycisk | 2 | Kliknięcie przycisku | domyślny · aktywny | Otwiera okno oceny jako rozszerzenie boczne | Koniec kolumny okna |
| Menu widoku i filtrów | Menu kebab `⋮` | Filtrowanie zapisu, wskaźnik konwergencji, dryf stanowiska, eksport transkryptu | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę filtrów i wariantów eksportu | Nagłówek okna |
| Reguły klasyfikacji wypowiedzi | Katalog etykiet i progów klasyfikacji | Dostosowanie klasyfikacji aktów mowy do specyfiki debaty | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny | — | Zmienia sposób etykietowania wypowiedzi w bieżącej karcie sesji | Okno konfiguracji, presety modułu |

---

### 4.5. Argument Map & Analysis

| Aspekt | Wartość |
|---|---|
| Typologia | Okno monitorów i dashboardów |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 2 |
| Izolacja domyślna | Buduje strukturę na podstawie zawartości Debate Panel w obrębie karty sesji |

Analityczne serce debaty — struktura argumentacji zamiast zapisu liniowego.

**Zawartość i pełny arsenał funkcji.**

- **Argument Map** — interaktywny graf teza→argument→kontrargument→riposta; kliknięcie węzła podświetla wypowiedź w chronologii.
- **Argument Mining** — wydobycie węzłów z wypowiedzi; scalanie powtórzeń (semantic dedup) z licznikiem poparcia.
- **Fallacy Detection** i **Steelman/Strawman Check** — oznaczenie błędów logicznych i rzetelności kontry na węzłach.
- **Evidence Ledger** — rejestr dowodów i cytowań, z flagą twierdzeń niepopartych (powiązanie z Research i Browser).
- **Crux Finder** — wskazanie kluczowego punktu spornego, którego rozstrzygnięcie zmienia wnioski.
- **Opinion Clustering** i **Agreement Matrix** — mapy obozów i macierz zgodności uczestników.
- Eksport: SVG, PNG oraz DOT, GraphML, Argdown, AIF.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Argument Map & Analysis ──────┐
│ Widok: graf                 ⋮  │
├────────────────────────────────┤
│      ┌─ teza ─┐                │
│      │ wejście│                │
│      └───┬────┘                │
│     ┌────┴─────┐               │
│  ┌──▼──┐    ┌──▼──┐            │
│  │ arg │    │kontr│ ⚠          │
│  └──┬──┘    └──┬──┘            │
│     │       ┌──▼──┐            │
│     │       │ripos│            │
│     │       └─────┘            │
├────────────────────────────────┤
│ ⚡ crux: skala inwestycji      │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Graf argumentów | Interaktywny graf węzłów | Wizualizacja struktury tezy, argumentów i kontrargumentów | Duży obszar canvas | 1 (w obrębie okna) | Widoczny po otwarciu okna | pusty stan · wyrenderowany | Kliknięcie węzła podświetla powiązaną wypowiedź w Debate Panel | Ciało okna |
| Węzeł argumentu | Blok grafu z etykietą klasyfikacji | Reprezentacja jednostki argumentacyjnej | Mały blok w grafie | 1 | Widoczny bez interakcji | teza · argument · kontrargument · riposta · scalony | Przeciąganie zmienia układ; kliknięcie otwiera szczegóły węzła | Ciało okna |
| Znacznik błędu logicznego | Ikona ostrzeżenia na węźle | Oznaczenie chwytu erystycznego lub błędu logicznego | Bardzo mała ikona | 1 | Widoczny przy wykryciu | ukryty · widoczny | Kliknięcie rozwija nazwę błędu i uzasadnienie | Węzeł grafu |
| Wskazanie cruxa | Blok `.dn-alert--ostrz` z etykietą | Wskazanie kluczowego punktu spornego | Mały blok | 1 | Widoczny przy wykryciu | ukryty · widoczny | Kliknięcie otwiera Moderator Panel z gotową interwencją | Koniec kolumny okna |
| Przełącznik widoku analizy | Zwinięty element zbiorczy `Widok ▼` | Przejście między grafem, macierzą zgodności i klastrami opinii | Mały element zbiorczy | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty | Zmienia reprezentację tej samej debaty; zwija się samoczynnie | Góra okna |
| Menu analizy i eksportu | Menu kebab `⋮` | Rejestr dowodów, kontrola steelman/strawman, warianty eksportu grafu | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji analitycznych i formatów eksportu | Nagłówek okna |
| Katalog błędów logicznych | Zbiór definicji wykrywanych chwytów | Dostosowanie zakresu wykrywania błędów logicznych | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny | — | Zmienia zestaw wykrywanych błędów dla bieżącej karty sesji | Okno konfiguracji, presety modułu |

---

### 4.6. Voting & Evaluation Center

| Aspekt | Wartość |
|---|---|
| Typologia | Okno monitorów i dashboardów |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 2 |
| Izolacja domyślna | Wyniki głosowań i ocen zapisywane per karta sesji; ranking Elo akumulowany między sesjami |

Ocena, głosowanie i ranking w jednym miejscu.

**Zawartość i pełny arsenał funkcji.**

- **Vote Engine** — głosowania metodami: aprobata, ranking IRV, Condorcet/Schulze, skala punktowa, kwadratowe; próg quorum i konsensusu.
- **Rubric Builder** i **LLM-as-Judge** — kryteria z wagami, ocena modelami-sędziami z uzasadnieniem, tryb jury.
- **Multi-Criteria Decision Matrix** — zestawienie wariantów w macierzy kryteriów (styl AHP) z wynikiem.
- **Model Elo / TrueSkill Leaderboard** — ranking modeli i person akumulowany między sesjami.
- **Confidence & Calibration** — konfrontacja deklarowanej pewności z trafnością.
- **Star & Preference Rating** — ocena gwiazdkami i porównanie parami.
- Raport z głosowań i rubryk z wykresami; eksport scorecardu.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Voting & Evaluation ──────────┐
│ Głosowanie: aprobata        ⋮  │
├────────────────────────────────┤
│ Stanowisko 1  ████████  6 gł.  │
│ Stanowisko 2  ████      3 gł.  │
│ Stanowisko 3  ██        1 gł.  │
├────────────────────────────────┤
│ Quorum 2/3 — osiągnięte ✓      │
├────────────────────────────────┤
│ Ranking uczestników            │
│ 1. Model A   1284              │
│ 2. Zwolennik 1201              │
│ 3. Model B   1176              │
├────────────────────────────────┤
│ [ Uruchom głosowanie ]         │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Zestawienie wyników głosowania | Lista pozycji z paskami wyniku | Prezentacja rozkładu głosów nad stanowiskami | Średni blok, ciało okna | 1 (w obrębie okna) | Widoczne po otwarciu okna | brak głosowania · w toku · rozstrzygnięte | Kliknięcie pozycji rozwija głosy poszczególnych uczestników | Ciało okna |
| Wskaźnik quorum | Plakietka z progiem i stanem | Informacja o osiągnięciu progu konsensusu | Mała plakietka `.dn-plakietka--sukces` / `--ostrz` | 1 | Widoczny bez interakcji | osiągnięte · nieosiągnięte | — (informacyjny) | Ciało okna |
| Tabela rankingu | Lista uczestników z punktacją | Ranking modeli i person między sesjami | Średnia lista | 1 | Widoczna bez interakcji | aktualizowana po rozstrzygnięciu | Kliknięcie pozycji otwiera historię pojedynków uczestnika | Ciało okna |
| Przycisk „Uruchom głosowanie” | Przycisk `--zloty` | Rozpoczęcie głosowania nad bieżącymi stanowiskami | Średni przycisk CTA | 1 | Widoczny bez interakcji | domyślny · ładowanie | Zleca Koordynatorowi zadania głosowania i oceny dla wykonawców | Koniec kolumny okna |
| Selektor metody głosowania | Zwinięty element zbiorczy `Metoda ▼` | Wybór metody: aprobata, IRV, Condorcet/Schulze, skala punktowa, kwadratowa | Mały element zbiorczy | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty | Zmienia sposób agregacji głosów; zwija się samoczynnie | Góra okna |
| Menu oceny i raportu | Menu kebab `⋮` | Tryb jury, ocena modelami-sędziami, macierz decyzyjna, kalibracja pewności, eksport scorecardu | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji oceny i wariantów raportu | Nagłówek okna |
| Rubric Builder | Edytor kryteriów oceny z wagami | Definicja rubryki oceny odpowiedzi i stanowisk | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, wyszukiwarka funkcji, tryb administracyjny | — | Otwiera edytor kryteriów i wag; zapis jako preset modułu | Okno konfiguracji, presety modułu |
| Konfiguracja rankingu Elo | Ustawienia algorytmu rankingowego | Wybór algorytmu (Elo, Glicko, TrueSkill) i parametrów akumulacji | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny | — | Zmienia sposób wyliczania rankingu między sesjami | Okno konfiguracji, presety modułu |

---

### 4.7. Moderator Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Panel narzędziowy |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne |
| Warstwa | 2 |
| Izolacja domyślna | Ustawienia moderacji obowiązują przez cały czas trwania sesji Roundtable |

**Zawartość i pełny arsenał funkcji.**

- Wybór formatu debaty: swobodna wymiana, debata strukturalna (tezy → riposty → wnioski), format oksfordzki (za/przeciw), tura okrężna, panel ekspercki, rundy Delphi.
- Ustalenie liczby tur oraz limitu czasu i długości wypowiedzi na turę.
- Wprowadzenie nowego wątku lub pytania pobocznego w trakcie trwającej debaty.
- Wskazanie kolejności głosu (kto odpowiada pierwszy, kto ostatni w danej turze) oraz przełączenie trybu równoległego i sekwencyjnego.
- Zamknięcie bieżącej tury ręcznie, zanim wszyscy uczestnicy zakończą wypowiedź.
- Interwencja moderująca — polecenie skierowane do wszystkich uczestników naraz („skupcie się na aspekcie finansowym”, „unikajcie powtórzeń”), wstrzykiwane przez Koordynatora do kontekstu zadań tury.
- Timer tury z podglądem czasu pozostałego do zamknięcia rundy.
- Wprowadzenie kontrgłosu przy nadmiernej zgodzie panelu oraz reguły miękkie tonu i dominacji uczestnika.
- Zapisanie konfiguracji moderacji jako szablonu do ponownego użycia w kolejnych debatach.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Moderator Panel ──────────────┐
│ [ Format: Strukturalna ▼ ]  ⋮  │
│ [ Tury: 3 ▼ ]                  │
├────────────────────────────────┤
│ Timer tury: ⏱ 00:42            │
│ [ Zamknij turę teraz ]         │
├────────────────────────────────┤
│ Interwencja:                   │
│ [ …………………… ]        [ Wyślij ] │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Timer tury | Zegar odliczający z paskiem postępu | Wizualizacja czasu pozostałego do zamknięcia rundy | Mały element graficzny | 1 (w obrębie okna) | Widoczny po otwarciu okna | nieaktywny (brak limitu) · odliczanie · czas upłynął | Po upłynięciu czasu sygnalizuje zamknięcie tury Koordynatorowi | Ciało okna |
| Przycisk „Zamknij turę teraz” | Przycisk `--zarys` | Ręczne zakończenie bieżącej rundy | Mały przycisk | 1 | Widoczny bez interakcji | domyślny | Kończy turę niezależnie od stanu zadań poszczególnych wykonawców | Ciało okna |
| Pole interwencji moderującej | Pole `.dn-input` z przyciskiem wysyłki | Polecenie skierowane jednocześnie do wszystkich uczestników | Małe pole tekstowe z przyciskiem `--zloty` | 1 | Widoczne bez interakcji | domyślny · wysłano | Przekazuje instrukcję Koordynatorowi, który wstrzykuje ją do kontekstu zadań tury | Koniec kolumny okna |
| Selektor formatu debaty | Zwinięty element zbiorczy `Format ▼` | Wybór schematu prowadzenia debaty | Mały element zbiorczy | 2 | Kliknięcie przełącznika lub znacznika kontekstowego w Chat Window | zwinięty · rozwinięty | Zmienia sposób dekompozycji zlecenia na zadania kolejnych tur; zwija się samoczynnie | Góra okna |
| Selektor liczby tur | Zwinięty element zbiorczy `Tury ▼` | Ustalenie zakładanej długości debaty | Mały element zbiorczy | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty | Ogranicza lub odblokowuje automatyczne zamknięcie po N turach | Góra okna |
| Menu moderacji | Menu kebab `⋮` | Kolejność głosu, wątek poboczny, zapis szablonu moderacji, reguły tonu i dominacji | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Rozwija listę operacji moderacyjnych | Nagłówek okna |
| Reguły wprowadzania kontrgłosu | Ustawienia detekcji nadmiernej zgody | Definicja progu zbieżności i persony kontry | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny | — | Uruchamia wprowadzenie kontrgłosu po przekroczeniu progu zbieżności | Okno konfiguracji, presety modułu |

---

### 4.8. Consensus Panel

| Aspekt | Wartość |
|---|---|
| Typologia | Okno edycyjne / wynikowe |
| Waga wizualna | Kolumna boczna, otwierana jako rozszerzenie boczne po dowolnej turze lub po zakończeniu debaty |
| Warstwa | 2 |
| Izolacja domyślna | Buduje stanowisko na podstawie zawartości Debate Panel i argumentów oznaczonych jako kluczowe |

**Zawartość i pełny arsenał funkcji.**

- Wygenerowanie szkicu uzgodnionego stanowiska na podstawie przebiegu debaty i argumentów oznaczonych jako kluczowe.
- Edycja ręczna finalnego stanowiska przez użytkownika — treść wygenerowana jest punktem wyjścia, nie wynikiem ostatecznym.
- Zestawienie punktów zgody i punktów pozostających spornych (stanowisko mniejszościowe zachowane i podpisane, nie ukrywane).
- Przypisanie wagi i poparcia poszczególnym wnioskom (które stanowisko poparło ilu i jak ważnych uczestników).
- Operacje kontekstowe AI na treści stanowiska końcowego (korekta, zmiana stylu, skrócenie) — analogiczne do Tools Panel modułu Studio.
- Sformatowanie wyniku jako zapisu decyzji (ADR): kontekst, warianty, decyzja, konsekwencje, ryzyko mniejszości.
- Eksport stanowiska końcowego jako dokumentu, z załącznikiem pełnego transkryptu debaty; odsłuch przebiegu debaty syntezą mowy z głosem per uczestnik.
- Przekazanie stanowiska końcowego do modułu Studio do dalszej redakcji lub do Research jako ustalenie badawcze.
- Wersjonowanie stanowiska — kolejne redakcje po kolejnych turach debaty zapisywane jako odrębne wersje z porównaniem.

**Makieta tekstowa (stan spoczynku).**

```
┌─ Consensus Panel ──────────────┐
│ 3 tury · 4 argumenty kluczowe ⋮│
├────────────────────────────────┤
│ STANOWISKO UZGODNIONE          │
│ „Wejście na rynek azjatycki    │
│  jest zasadne pod warunkiem    │
│  wdrożenia zabezpieczeń…”      │
├────────────────────────────────┤
│ Punkty zgody: 3   Sporne: 1    │
│ ⚡ „skala inwestycji” —        │
│    zdanie odrębne (Model B)    │
├────────────────────────────────┤
│ [ Eksportuj ▼ ]                │
└────────────────────────────────┘
```

**Katalog elementów interfejsu.**

| Element | Co to jest | Do czego służy | Forma i waga | Warstwa | Sposób wywołania | Stany | Zachowanie po interakcji | Występowanie |
|---|---|---|---|---|---|---|---|---|
| Nagłówek źródła stanowiska | Mały tekst informacyjny | Wskazanie, na podstawie ilu tur i argumentów zbudowano stanowisko | Bardzo mały tekst `.dn-tekst-3` | 1 (w obrębie okna) | Widoczny po otwarciu okna | aktualizowany na żywo | Kliknięcie przewija do odpowiadających wypowiedzi w Debate Panel | Góra okna |
| Obszar treści stanowiska | Edytowalny blok tekstu | Redakcja finalnego, uzgodnionego wniosku | Duży blok, ciało okna | 1 | Widoczny bez interakcji | wygenerowany szkic · edytowany · zaakceptowany | Edycja ręczna; zaznaczenie ujawnia menu operacji kontekstowych AI | Ciało okna |
| Licznik punktów zgody i spornych | Para małych plakietek | Podsumowanie liczbowe wyniku debaty | Plakietki `.dn-plakietka--sukces` / `--ostrz` | 1 | Widoczny bez interakcji | aktualizowany na żywo | Kliknięcie rozwija listę odpowiednich punktów | Ciało okna |
| Blok zdania odrębnego | Wyróżniony fragment tekstu | Zachowanie stanowiska uczestnika, który nie dołączył do konsensusu | Mały blok `.dn-alert--info` z podpisem uczestnika | 1 | Widoczny przy braku pełnej zgody | ukryty (pełna zgoda) · widoczny | — (informacyjny, integralna część eksportu) | Pod głównym stanowiskiem |
| Element zbiorczy „Eksportuj” | Zwinięty element zbiorczy `Eksportuj ▼` | Pobranie stanowiska końcowego, z transkryptem lub bez | Średni element zbiorczy | 2 | Kliknięcie przełącznika | zwinięty · rozwinięty · ładowanie | Rozwija listę formatów: PDF, DOCX, Markdown, ADR; po wyborze zwija się samoczynnie | Koniec kolumny okna |
| Menu operacji kontekstowych AI | Menu kontekstowe zaznaczenia | Korekta, zmiana stylu, skrócenie treści stanowiska | Menu kontekstowe | 3 | Zaznaczenie fragmentu treści | ukryte (brak zaznaczenia) · widoczne | Analogiczne do Tools Panel modułu Studio | Obszar treści stanowiska |
| Menu wydania i przekazania | Menu kebab `⋮` | Wysłanie do Studio, wysłanie do Research, wersjonowanie stanowiska, odsłuch syntezą mowy | Ikona | 3 | Kliknięcie `⋮` | zwinięte · rozwinięte | Tworzy dokument w Studio Editor lub wpis w Findings Panel, gdy powiązanie skonfigurowane | Nagłówek okna |
| Szablon zapisu decyzji (ADR) | Definicja struktury zapisu decyzji | Dostosowanie układu protokołu ustaleń | Bez reprezentacji graficznej w stanie spoczynku | 4 | Polecenie języka naturalnego w Chat Window, tryb administracyjny | — | Zmienia strukturę generowanego zapisu decyzji | Okno konfiguracji, presety modułu |

---

## 5. Przepływy pracy

### 5.1. Przepływ podstawowy — od pytania do konsensusu

```
 Chat Window       Execution Loop      Model Panels /        Consensus Panel
 (Użytkownik ↔     (Koordynator ↔      Debate Panel
  Wykonawca)        Wykonawca)
 ────────────      ──────────────      ───────────────       ────────────────
 1. sformułowanie
    pytania    ────►
                   2. dekompozycja
                      zlecenia na
                      zadania tury
                          │
                          ▼
                   3. przydział zadań
                      wielu wykonawcom ────►
                                            4. równoległe
                                               odpowiedzi
                                                    │
                   5. kontrola jakości  ◄───────────┘
                      i domknięcie tury
                          │
                          ▼
                                            6. rejestr wymiany
                                               argumentów, wykrycie
                                               punktów zgody i spornych
                                                    │
                                                    ▼
                                                                7. szkic
                                                                   stanowiska,
                                                                   redakcja
                                                                   i akceptacja
```

### 5.2. Przepływ — pętla wykonawcza narady wielu wykonawców

```
Chat Window — polecenie „uruchom turę 2”
        │
        ▼
Execution Loop Window — Koordynator dekomponuje zlecenie:
        zadanie A: wypowiedź Model A        zadanie D: klasyfikacja wypowiedzi
        zadanie B: wypowiedź Model B        zadanie E: ocena jury
        zadanie C: wypowiedź Zwolennik
        │
        ▼
Przydział równoległy wykonawcom A, B, C · kolejność głosu wg Moderator Panel
        │
        ├──► wykonawca B przekracza limit czasu → ponowienie zadania
        ├──► wynik wykonawcy C odrzucony w kontroli jakości → komunikat zwrotny
        │
        ▼
Domknięcie tury — scalenie wyników, zapis w Debate Panel,
        budowa kontekstu tury kolejnej
        │
        ▼
Chat Window — potwierdzenie zamknięcia tury i podsumowanie stanu debaty
```

### 5.3. Przepływ — debata z dwiema tożsamościami tego samego modelu

```
Model Panels — dodanie uczestnika „Zwolennik” (kanał modelu A, prompt: argumentacja za)
Model Panels — dodanie uczestnika „Oponent”   (ten sam kanał modelu A, prompt: krytyka)
        │
        ▼
Execution Loop Window — Koordynator traktuje obie tożsamości jako odrębnych
                wykonawców z odrębnymi zadaniami w tej samej turze
        │
        ▼
Debate Panel — rejestruje wymianę argumentów między obiema tożsamościami,
                mimo że technicznie oba przypisania korzystają
                z tego samego kanału modelu (Integracja modeli, Załącznik B.3)
        │
        ▼
Consensus Panel — stanowisko uwzględniające obie perspektywy jednego modelu
                    zestawione z perspektywami pozostałych, niezależnych modeli
```

### 5.4. Przepływ — moderacja aktywna w trakcie debaty

```
Moderator Panel — wybór formatu „strukturalna”, 3 tury, timer 60 s na turę
        │
        ▼
Execution Loop Window — tura 1 domknięta po upłynięciu limitu czasu
        │
        ▼
Moderator Panel — wykryty punkt sporny → interwencja: „doprecyzujcie skalę ryzyka”
        │
        ▼
Execution Loop Window — Koordynator wstrzykuje interwencję do kontekstu
                zadań tury 2 dla wszystkich wykonawców
        │
        ▼
Model Panels — tura 2, odpowiedzi uwzględniające interwencję moderatora
        │
        ▼
Moderator Panel — „Zamknij turę teraz” (ręczne zakończenie przed upłynięciem czasu)
```

### 5.5. Przepływ — analiza argumentów i głosowanie

```
Debate Panel — zamknięta tura z pełnym zapisem wypowiedzi
        │
        ├──► Argument Map & Analysis — wydobycie węzłów, scalenie powtórzeń,
        │        oznaczenie błędów logicznych, wskazanie cruxa
        │
        └──► Voting & Evaluation Center — głosowanie metodą Condorcet/Schulze
                 nad stanowiskami, ocena modelami-sędziami wg rubryki,
                 sprawdzenie progu quorum
        │
        ▼
Consensus Panel — stanowisko z poparciem ważonym i zachowanym zdaniem odrębnym
```

### 5.6. Przepływ — wynik badania skonfrontowany wielomodelowo

```
Research › Findings Panel / Report Builder
──────────────────────────────────────────
 wstępna synteza ustaleń      ─────────────►   Roundtable › Chat Window
 (powiązanie konfigurowalne)                    (pytanie wyjściowe debaty)
                                                          │
                                                          ▼
                                              Execution Loop Window
                                              (dekompozycja i koordynacja
                                               wykonawców narady)
                                                          │
                                                          ▼
                                              Model Panels / Debate Panel
                                              (ocena mocnych i słabych stron
                                               wnioskowania z różnych perspektyw)
                                                          │
                                                          ▼
                                              Consensus Panel — stanowisko
                                              końcowe   ─────────────────────►
                                                                                Research › Report Builder
                                                                                (wersja finalna wniosków)
```

---

## 6. Punkty sterowania z okna konfiguracji

Wszystko personalizuje Operator z okna konfiguracji (Model konfiguracji, rozdz. 5–6), na właściwej warstwie (globalna → środowisko → projekt → sesja/karta). Zero blokad — brak ustawienia oznacza wartość domyślną, nie zatrzymanie pracy.

| Zakres konfiguracji (Model konfiguracji) | Co personalizuje Operator w module Roundtable |
|---|---|
| 5.4 Zachowanie modeli | Kanały uczestników (API/CLI/SSH/HTTP), dobór modeli, parametry wywołania (temperatura, próbki dla self-consistency); **klucze API jawne** dla modeli zewnętrznych |
| 5.5 Tożsamość modeli | Persony i role uczestników (nazwa, awatar, opis), wagi ekspertyzy; wpięcie zapisanych agentów jako uczestników z ich uprawnieniami |
| 5.6 Prompty systemowe | Prompty ról (adwokat diabła, głos ostrożności, ekspert), prompt sędziego dla oceny modelami-sędziami, biblioteka ról jako presety |
| — Pętla wykonawcza (Execution Loop Window) | Tryb przydziału zadań tury (równoległy, sekwencyjny), limit czasu zadania, liczba ponowień, progi kontroli jakości wyniku, zakres rejestru komunikatów sterujących Koordynator ↔ Wykonawca, domyślna widoczność kolumny pętli w układzie karty |
| 5.7 Rozszerzenia | Konektory modeli zewnętrznych, silniki TTS (odsłuch debaty), serwery MCP wystawiające i przyjmujące operacje debaty; źródło Danaco Plugin lub Personal |
| 5.8 Integracje | **Research ──► Roundtable** (materiał wejściowy), **Roundtable ──► Studio** i **──► Research** (przekazanie stanowiska), **Agents ──► Roundtable** (agent jako uczestnik), udostępnienie debaty jako kroku **Automations**; webhooki i konektory |
| 5.9 Historia | Współdzielenie historii debaty między kartami tego samego zagadnienia; retencja transkryptów |
| 5.11 Izolacja | Osiem zakresów izolacji technicznej per uczestnik (odrębne konto i token per kanał modelu) dla debat nad materiałem poufnym; izolacja kontekstu między kartami |
| 5.12 Karty sesji | Migawka układu ośmiu okien; przypisanie projektu do karty; współdzielenie składu i zespołu między kartami |
| — Presety modułu (5.1/5.6/7) | Zespoły (skład + role + format), szablony moderacji (format, liczba tur, kolejność), rubryki oceny z wagami, metoda głosowania domyślna, próg quorum i konsensusu, polityka trybu ślepego, reguły wprowadzania kontrgłosu, katalog błędów logicznych, konfiguracja rankingu Elo, szablon zapisu decyzji |

Zasada kluczy jawnych: dane dostępowe kanałów modeli zewnętrznych (dostawcy API, TTS) są widoczne i edytowalne w zakresie 5.4 „Dane dostępowe kanału”, przechowywane poza bazą danych zgodnie z Modelem danych (rozdz. 8.1) — nigdy nie są ukryte ani zaszyte na stałe.

---

## 7. Stany, dane i powiązania

### 7.1. Model stanów sesji Roundtable

```
┌────────────────┐  dodanie          ┌────────────────┐  sformułowanie   ┌────────────────┐
│ Skład           │  uczestników      │ Gotowa do       │  pytania         │ Debata w toku   │
│ nieustalony     │──────────────────►│ rozpoczęcia     │─────────────────►│ (tury pętli     │
└────────────────┘                    └────────────────┘                   │  wykonawczej)   │
                                                                            └───────┬────────┘
                                                                                    │ zamknięcie
                                                                                    │ ostatniej tury
                                                                                    ▼
                                                                           ┌────────────────┐
                                                                           │ Ocena           │
                                                                           │ i głosowanie    │
                                                                           └───────┬────────┘
                                                                                    │ osiągnięty próg
                                                                                    ▼
                                                                           ┌────────────────┐
                                                                           │ Generowanie     │
                                                                           │ konsensusu      │
                                                                           └───────┬────────┘
                                                                                    │ akceptacja
                                                                                    ▼
                                                                           ┌────────────────┐
                                                                           │ Stanowisko      │
                                                                           │ zaakceptowane   │
                                                                           └────────────────┘
```

Debata pozostaje w pełni kompozycyjna — użytkownik wraca w każdej chwili do dodania kolejnego uczestnika, otwiera nowy wątek poboczny w Moderator Panel, uruchamia głosowanie po dowolnej turze albo generuje Consensus Panel przed formalnym zakończeniem debaty.

### 7.2. Model danych wykorzystywany przez moduł

| Encja (model danych) | Rola w module Roundtable |
|---|---|
| `sesja`, `karta_sesji` | Nośnik składu uczestników, przebiegu debaty i stanu pętli wykonawczej |
| `wiadomosc` | Wypowiedzi poszczególnych uczestników, powiązane z `kanal_modelu_id` generującym odpowiedź |
| `zlecenie`, `zadanie` | Zlecenie debaty i jego dekompozycja na zadania przydzielane wykonawcom w Execution Loop Window |
| `kanal_modelu` | Techniczne połączenie z modelem bazowym; jeden kanał wykorzystywany przez wielu uczestników o odrębnych tożsamościach (Integracja modeli, rozdz. 5, Załącznik B.3) |
| `artefakt`, `wersja_artefaktu` | Stanowisko końcowe z Consensus Panel jako dokument, transkrypt debaty i graf argumentów oraz ich kolejne redakcje |
| `kanal_modelu` — tożsamość i persona | Nazwa, awatar i prompt systemowy nadane uczestnikowi niezależnie od kanału technicznego |
| `powiazanie_komponentu` | Pokrewieństwo funkcjonalne Roundtable ═══ MultitaskingAI; konfiguracyjne Research ───► Roundtable, Roundtable ───► Studio/Research, Agents ───► Roundtable |

### 7.3. Izolacja i konfigurowalność — punkty właściwe modułowi Roundtable

| Punkt izolacji | Stan wyjściowy (domyślny) | Co można skonfigurować | Gdzie |
|---|---|---|---|
| Historia debaty | Odrębna per karta sesji | Współdzielenie historii debaty między kartami tego samego zagadnienia | Okno konfiguracji, poziom „karta sesji” |
| Rejestr komunikatów pętli wykonawczej | Odrębny per karta sesji | Zakres zapisu wymiany Koordynator ↔ Wykonawca oraz retencja rejestru | Okno konfiguracji, zakres pętli wykonawczej |
| Materiał wejściowy z Research | Brak automatycznego zasilania | Ustanowienie kanału Research ───► Roundtable dla wstępnej syntezy | Okno konfiguracji, powiązanie komponentu |
| Przekazanie konsensusu dalej | Brak automatycznego przekazania | Ustanowienie Roundtable ───► Studio lub Roundtable ───► Research | Okno konfiguracji, powiązanie komponentu |
| Izolacja techniczna per uczestnik (kanał modelu) | Model procesu współdzielony domyślnie | Przypisanie odrębnego konta i tokenu poszczególnym kanałom modeli uczestniczącym w debacie | Okno konfiguracji punktów izolacji, panel macierzy izolacji |
| Ranking uczestników między sesjami | Akumulowany w obrębie środowiska | Zakres akumulacji rankingu (środowisko, projekt, karta) oraz algorytm rankingowy | Okno konfiguracji, presety modułu |

### 7.4. Powiązania z innymi modułami

```
                       ROUNDTABLE   ═══   Środowisko MultitaskingAI
        (pokrewieństwo funkcjonalne — ta sama idea wielomodelowości;
         Roundtable = debata i konsensus, MultitaskingAI = podział ról
         wykonawczych i orkiestracja procesu, rozdz. 13 Koncepcji)

     Research ───►  ROUNDTABLE  ───►  Studio / Research / Automations
   Agents   ───►                     (przekazanie stanowiska końcowego
  (agent jako                          i udostępnienie debaty jako kroku
   uczestnik)                          procesu wsadowego)
```

| Element pokrewny / moduł docelowy | Charakter powiązania | Typ | Okno źródłowe → okno docelowe |
|---|---|---|---|
| Środowisko MultitaskingAI | Rozwinięcie idei wielomodelowości w formie zorientowanej na role, orkiestrację i kolejkowanie | Pokrewieństwo funkcjonalne | Model Panels / Execution Loop Window ═══ Executor Chat / Coordinator Chat / Results Analyzer |
| Research (jako źródło) | Wielomodelowa prezentacja wyników badawczych przed opracowaniem wniosków końcowych | Konfiguracyjne | Report Builder / Findings Panel → Chat Window |
| Agents (jako źródło) | Wpięcie zapisanego agenta jako uczestnika debaty z jego skillami i pamięcią | Konfiguracyjne | Agent Manager → Model Panels |
| Studio (jako cel) | Dalsza redakcja uzgodnionego stanowiska | Konfiguracyjne | Consensus Panel → Studio Editor |
| Research (jako cel) | Wykorzystanie stanowiska jako ustalenia badawczego | Konfiguracyjne | Consensus Panel → Findings Panel |
| Automations (jako cel) | Udostępnienie debaty wielomodelowej jako kroku procesu wsadowego | Konfiguracyjne | Execution Loop Window → Automations Builder |
| Library (jako cel) | Złożenie transkryptu, grafu argumentów i protokołu ustaleń w repozytorium | Konfiguracyjne | Consensus Panel / Debate Panel → Library |

---

## 8. Scenariusze użycia

**Scenariusz 1 — decyzja strategiczna skonfrontowana wielomodelowo.**
Zarząd rozważa wejście na nowy rynek. Użytkownik otwiera Roundtable, dodaje trzech uczestników — model nastawiony na analizę szans, model nastawiony na analizę ryzyka i model neutralny — i zadaje wspólne pytanie w Chat Window. Execution Loop Window pokazuje dekompozycję zlecenia na zadania trzech wykonawców i stan każdego z nich. Po trzech turach debaty Consensus Panel przedstawia stanowisko warunkowe, z zachowanym zdaniem odrębnym dotyczącym skali inwestycji.

**Scenariusz 2 — debata „Zwolennik kontra Oponent” na jednym modelu.**
Redaktor sprawdza odporność tezy artykułu na krytykę, korzystając z jednego dostawcy modeli. Dodaje dwóch uczestników korzystających z tego samego kanału modelu, ale z odrębnymi tożsamościami i promptami — „Zwolennik” i „Oponent” — i obserwuje w Debate Panel pełną wymianę argumentów mimo wspólnego technicznego źródła. Koordynator prowadzi obie tożsamości jako odrębnych wykonawców tej samej tury.

**Scenariusz 3 — moderowana debata strukturalna z limitem czasu.**
Facylitator warsztatu ustawia w Moderator Panel format „strukturalna”, trzy tury i limit 60 sekund na wypowiedź. Timer dyscyplinuje przebieg, Execution Loop Window sygnalizuje wykonawcę przekraczającego limit i ponawia jego zadanie, a po zakończeniu ostatniej tury facylitator ręcznie domyka rundę przed wygenerowaniem konsensusu.

**Scenariusz 4 — konfrontacja wniosków badawczych z modułu Research.**
Analityk przekazuje wstępną syntezę ustaleń z modułu Research do Roundtable. Trzej uczestnicy oceniają mocne i słabe strony wnioskowania z perspektywy metodologicznej, biznesowej i ryzyka regulacyjnego. Argument Map & Analysis wskazuje kluczowy punkt sporny, a Consensus Panel gromadzi finalne stanowisko, które trafia z powrotem do Report Builder jako sekcja wniosków.

**Scenariusz 5 — decyzja architektoniczna w CodeStudio.**
Zespół programistyczny w środowisku CodeStudio korzysta z Roundtable do oceny wyboru między dwoma podejściami architektonicznymi. Model Panels prezentują argumenty za każdym podejściem, Debate Panel rejestruje wymianę technicznych kontrargumentów, Voting & Evaluation Center przeprowadza głosowanie rankingowe z progiem 2/3, a Consensus Panel formułuje zapis decyzji wraz z odnotowanym ryzykiem mniejszościowym.

**Scenariusz 6 — ocena rubrykowa z udziałem jury modeli.**
Zespół redakcyjny porównuje pięć wariantów komunikatu. Trzej uczestnicy przedstawiają warianty, dwaj kolejni pełnią rolę jury i oceniają je wg rubryki „trafność 40%, spójność 30%, dowody 30%”. Koordynator przydziela zadania oceny jako odrębne pozycje kolejki, a Voting & Evaluation Center prezentuje scorecard z uzasadnieniami i ranking uczestników.

---

## 9. Załącznik — skróty klawiszowe i ikonografia

| Skrót / ikona | Działanie | Okno |
|---|---|---|
| `Ctrl/Cmd + Shift + M` | Dodanie nowego uczestnika/modelu | Model Panels |
| `Tab` | Przejście między panelami uczestników | Model Panels |
| `Ctrl/Cmd + Enter` | Uruchomienie kolejnej tury | Debate Panel, Moderator Panel |
| `Ctrl/Cmd + Shift + L` | Otwarcie i zamknięcie Execution Loop Window | Wszystkie okna modułu |
| `Ctrl/Cmd + .` | Wstrzymanie i wznowienie pętli wykonawczej | Execution Loop Window |
| `Ctrl/Cmd + Shift + G` | Otwarcie Argument Map & Analysis | Debate Panel |
| `Ctrl/Cmd + Shift + V` | Otwarcie Voting & Evaluation Center | Debate Panel |
| Ikona `uzytkownik` | Awatar uczestnika (persona) | Model Panels, Debate Panel |
| Ikona `dzwonek` (przekreślona) | Uczestnik wyciszony w turze | Model Panels |
| Ikona `waga` | Punkt sporny / rozstrzygnięcie | Debate Panel, Argument Map & Analysis |
| Ikona `gwiazdka` | Argument oznaczony jako kluczowy | Model Panels, Debate Panel |
| Ikona `zegar` | Timer tury | Moderator Panel |
| Ikona `petla` | Stan pętli wykonawczej i postęp zadań | Execution Loop Window |
| Ikona `ostrzezenie` | Błąd logiczny lub wynik odrzucony w kontroli jakości | Argument Map & Analysis, Execution Loop Window |
| Ikona `ptaszek` | Punkt zgody / stanowisko zaakceptowane | Consensus Panel |
| Ikona `pobierz` | Eksport transkryptu, grafu lub stanowiska | Debate Panel, Argument Map & Analysis, Consensus Panel |

*Koniec dokumentu. Moduł Roundtable — dokumentacja projektowa, wersja 2.0, 2026-08-06.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
