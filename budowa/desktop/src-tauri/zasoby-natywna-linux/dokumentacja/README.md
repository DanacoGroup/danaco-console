# Danaco Console

**Platforma AI Workspace OS — system operacyjny do pracy ze sztuczną inteligencją.**

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Wersja** | v2.0 |
| **Producent** | Danaco Holding Group Sp. z o.o., ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Twórca** | Dariusz Naharnowicz |
| **Kontakt** | support@danaco-group.pl |

Danaco Console jest wielośrodowiskowym systemem operacyjnym dla sztucznej
inteligencji. Integruje komunikację, zarządzanie wiedzą, tworzenie treści,
projektowanie, automatyzację procesów oraz rozwój oprogramowania w jednej
platformie — tak, aby praca z modelami AI przestała być rozmową w oknie czatu
i stała się zorganizowaną, nadzorowaną pracą cyfrowej organizacji.

---

## Spis treści

1. [Czym jest Danaco Console](#1-czym-jest-danaco-console)
2. [Dlaczego nie kolejna aplikacja czatowa](#2-dlaczego-nie-kolejna-aplikacja-czatowa)
3. [Hierarchia pracy: środowisko → moduł → okno operacyjne](#3-hierarchia-pracy-środowisko--moduł--okno-operacyjne)
4. [Dwa kanały komunikacji operacyjnej](#4-dwa-kanały-komunikacji-operacyjnej)
5. [Cztery środowiska](#5-cztery-środowiska)
6. [Piętnaście modułów](#6-piętnaście-modułów)
7. [MultitaskingAI — praca ciągła 24/7/365](#7-multitaskingai--praca-ciągła-247365)
8. [Agenci i komponenty własne](#8-agenci-i-komponenty-własne)
9. [Funkcje globalne: Mobile i Always On Display](#9-funkcje-globalne-mobile-i-always-on-display)
10. [Modele i kanały integracji](#10-modele-i-kanały-integracji)
11. [Rozszerzenia](#11-rozszerzenia)
12. [Pełna konfigurowalność i punkty izolacji](#12-pełna-konfigurowalność-i-punkty-izolacji)
13. [Interfejs: warstwy widoczności](#13-interfejs-warstwy-widoczności)
14. [Architektura i model wdrożenia](#14-architektura-i-model-wdrożenia)
15. [Bezpieczeństwo i dostęp](#15-bezpieczeństwo-i-dostęp)
16. [Wymagania](#16-wymagania)
17. [Zasady nadrzędne platformy](#17-zasady-nadrzędne-platformy)
18. [Dokumentacja produktu](#18-dokumentacja-produktu)
19. [Słowniczek](#19-słowniczek)

---

## 1. Czym jest Danaco Console

Danaco Console zarządza wieloma równoległymi przestrzeniami roboczymi, każdą
z własnym stanem, kontekstem i zestawem narzędzi — tak jak system operacyjny
zarządza wieloma aplikacjami i procesami. Określenie „system operacyjny" jest
tu opisem funkcjonalnym, nie metaforą marketingową.

Fundament platformy tworzy pięć elementów układających się w jeden łańcuch
organizacji pracy:

| Element | Rola w łańcuchu organizacji pracy |
|---|---|
| **Środowiska** | Definiują sposób pracy — wyznaczają ogólny kontekst, w jakim odbywa się zadanie |
| **Moduły** | Definiują obszary funkcjonalne — uszczegóławiają kontekst środowiska do konkretnego zadania |
| **Okna operacyjne** | Właściwa przestrzeń wykonywania zadań — dostarczają narzędzi do ich realizacji |
| **Czat** | Uniwersalna warstwa komunikacji — pośredniczy w rozmowie z AI na każdym poziomie łańcucha |
| **Agenci** | Warstwa inteligencji ponad całym ekosystemem — działa niezależnie od położenia użytkownika w strukturze |

Wejście do łańcucha prowadzi przez stronę główną — centrum dowodzenia, z którego
Operator wybiera środowisko, buduje własne komponenty albo przechodzi do ustawień.
Każde przełączenie środowiska lub modułu tworzy nową, wyspecjalizowaną przestrzeń
roboczą z własnym układem, historią, kontekstem i zestawem narzędzi, zachowując
spójny sposób komunikacji z AI. Zmiana przestrzeni jest zawsze świadomym,
widocznym przejściem, nigdy ukrytą zmianą stanu wewnątrz jednego okna.

**Co platforma daje w praktyce.** Po spięciu środowiska MultitaskingAI z modułem
Automations i przypisaniu ról własnym agentom powstaje w pełni autonomiczna,
wieloagentowa pętla pracy zdolna działać 24 godziny na dobę, 7 dni w tygodniu —
z planowaniem, kolejkowaniem, kontrolą jakości i nadzorem Operatora dostępnym
także z telefonu.

---

## 2. Dlaczego nie kolejna aplikacja czatowa

Rozróżnienie między platformą a aplikacją czatową jest fundamentem produktu
i przesądza każdą decyzję architektoniczną.

| Wymiar | Klasyczna aplikacja czatowa | Danaco Console |
|---|---|---|
| Punkt wyjścia | Jedno, względnie stałe okno rozmowy | Zadanie, dla którego dobierana jest przestrzeń robocza |
| Dokładanie możliwości | Kolejne funkcje w tym samym oknie: załączniki, wtyczki, tryby, panele | Osobna, dopasowana przestrzeń robocza dla każdego rodzaju pracy |
| Kontekst konwersacyjny | Wspólny i niezmienny dla wszystkich zadań | Rekonfigurowany przy każdej zmianie kontekstu pracy |
| Rola czatu | Stały komponent interfejsu | Uniwersalna warstwa komunikacji — obecna wszędzie, lecz nie stała |
| Co zmienia się wraz z zadaniem | Wyłącznie treść rozmowy | Układ okien, dostępne narzędzia, historia sesji, pamięć kontekstowa, sposób działania AI |

Czat w module Developer „wie" o kodzie i repozytorium, czat w module Research —
o zebranych źródłach, choć w obu przypadkach jest to ten sam mechanizm. Zakres
współdzielenia historii i pamięci między modułami, sesjami i środowiskami ustala
Operator w oknie konfiguracji: ciągła, wspólna historia jest ustawieniem, nie
stanem narzuconym.

---

## 3. Hierarchia pracy: środowisko → moduł → okno operacyjne

```
ZADANIE OPERATORA
        │  (dobór właściwej przestrzeni roboczej)
        ▼
ŚRODOWISKO        „w jakim trybie pracuję?"     TalkIn · WorkSpace · CodeStudio · MultitaskingAI
        │
        ▼
MODUŁ             „jakie zadanie wykonuję?"     Studio · Research · Developer · … (15 modułów)
        │
        ▼
OKNO OPERACYJNE   właściwa przestrzeń wykonania Studio Editor · Code Editor · Report Builder · …

CZAT — uniwersalna warstwa komunikacji, rekonfigurowana na każdym poziomie hierarchii
```

Struktura trójstopniowa zamiast płaskiego okna czatu wynika z czterech przesłanek:

| Przesłanka | Na czym polega |
|---|---|
| **Dopasowanie przestrzeni do zadania** | Praca nad kodem wymaga edytora, drzewa projektu i terminala; praca badawcza — zarządzania źródłami i budowy raportu; automatyzacja — harmonogramu, kolejki i monitora wykonania. Płaski czat nie pomieści tej różnorodności bez przeciążenia |
| **Konfigurowalne punkty izolacji** | Środowiska i moduły dają odrębne przestrzenie, które można utrzymywać oddzielnie albo świadomie łączyć — kontekst jednego zadania nie zanieczyszcza innego wbrew woli Operatora |
| **Skalowalność** | Nowe rodzaje pracy dochodzą jako kolejne moduły wewnątrz istniejących środowisk, bez naruszania spójności przestrzeni już używanych |
| **Miejsce dla warstwy nadrzędnej** | Hierarchia tworzy naturalne miejsce dla warstwy agentowej i funkcji globalnych działających *ponad* strukturą, z dostępem do wszystkich środowisk jednocześnie |

Pracą równoległą zarządzają **karty sesji** w oknie środowiska — mechanika
analogiczna do kart przeglądarki. W obrębie jednej karty widoczny jest interfejs
jednego modułu; zmiana modułu przeładowuje tę kartę, a praca równoległa odbywa się
przez otwarcie kilku kart obok siebie. Zakres, w jakim karty współdzielą kontekst,
historię i pamięć, jest konfigurowalny.

---

## 4. Dwa kanały komunikacji operacyjnej

Pracą platformy steruje się dwoma kanałami komunikacji. Oba są elementami
pierwszoplanowymi architektury — nie dodatkami pojedynczych modułów — i występują
w każdym środowisku, module, oknie operacyjnym i przepływie pracy.

| Kanał | Okno | Uczestnicy | Rola |
|---|---|---|---|
| Pierwszy | **Chat Window** | Użytkownik ↔ Wykonawca | Centralny punkt pracy i podstawowy mechanizm sterowania wszystkimi procesami platformy |
| Drugi | **Execution Loop Window** | Koordynator ↔ Wykonawca | Pętla wykonawcza: koordynacja zadań, nadzór, orkiestracja, kontrola realizacji |

**Chat Window** przyjmuje polecenia w języku naturalnym, prezentuje strumień
odpowiedzi i artefaktów na żywo, pozwala zatwierdzić proponowane działanie albo
przerwać je w toku, a także zażądać wyjaśnienia wyniku i przyjętych założeń.
Zajmuje stałe miejsce w układzie — lewa kolumna obszaru roboczego, pełna wysokość,
szerokość regulowana przez Operatora.

**Execution Loop Window** pokazuje bieżące zlecenie wraz z jego dekompozycją na
zadania, kolejkę i stan każdego zadania, komunikaty sterujące między Koordynatorem
a Wykonawcą, wyniki kontroli jakości i decyzje o ponowieniu oraz wskaźniki
przebiegu pętli. Z tego okna Operator wstrzymuje, wznawia, przerywa i koryguje
zlecenie. Otwiera się jako kolumna sąsiadująca z Chat Window i pozostaje otwarte
przez cały czas trwania pętli.

| Rola | Charakter | Odpowiedzialność |
|---|---|---|
| **Użytkownik** | Osoba prowadząca pracę | Zleca zadania, zatwierdza wyniki, przerywa i koryguje działania |
| **Koordynator** | Komponent orkiestrujący platformy | Dekomponuje zlecenie, przydziela zadania, nadzoruje pętlę, kontroluje realizację |
| **Wykonawca** | Model AI, agent lub system wykonawczy | Realizuje zadania i zwraca wyniki |

Cały interfejs rozmieszczony jest w układzie pionowym — kolumny sąsiadują poziomo,
regulacji podlega wyłącznie ich szerokość. Okna pomocnicze otwierają się jako
kolejne kolumny po prawej stronie obszaru roboczego.

---

## 5. Cztery środowiska

Wybór środowiska jest pierwszą decyzją nawigacyjną Operatora i określa, które
moduły są dostępne w bocznej nawigacji.

| Środowisko | Tryb pracy | Przeznaczenie |
|---|---|---|
| **TalkIn** | Wiedza, komunikacja i praca z treścią | Rozmowy z AI, praca z dokumentami, analizy, tłumaczenia, badania, raporty, współpraca wielu modeli |
| **WorkSpace** | Produktywność, organizacja i realizacja projektów | Zarządzanie projektami, automatyzacja, prowadzenie procesów, aplikacje biznesowe, projektowanie, praca operacyjna |
| **CodeStudio** | Programowanie | Tworzenie oprogramowania, praca z kodem, debugowanie, terminale, budowanie aplikacji, procesy developerskie |
| **MultitaskingAI** | Orkiestracja autonomicznej pracy ciągłej | Zespół modeli i agentów w rolach; po spięciu z automatyką — pełna pętla pracy ciągłej 24/7/365 |

TalkIn, WorkSpace i CodeStudio mają w bocznej nawigacji listę modułów.
MultitaskingAI organizuje nie zadania modułowe, lecz zespół modeli i agentów —
dlatego w jego bocznej nawigacji znajduje się **panel orkiestracji**.

---

## 6. Piętnaście modułów

Moduł nie należy do jednego środowiska — jest dostępny w wybranych środowiskach.
Środowisko jest profilem widoczności modułów, a nie pojemnikiem, do którego moduł
przypisano na stałe.

| Moduł | Cel | Dostępność |
|---|---|---|
| **Studio** | Zaawansowana praca z tekstem, dokumentami i treścią: tworzenie, redagowanie, przekształcanie | TalkIn, WorkSpace |
| **Workspace** | Izolowane środowisko projektowe — własne instrukcje, pamięć, biblioteka i konfiguracja modeli | TalkIn, WorkSpace, CodeStudio |
| **Automations** | Budowanie i wykonywanie procesów automatycznych: harmonogramy, kolejki, pętle wielomodelowe | strona główna |
| **Browser** | Współdzielone przeglądanie internetu — Operator i AI pracują nad tą samą treścią jednocześnie | TalkIn, WorkSpace |
| **Research** | Realizacja badań, analiz i opracowań ze wielu źródeł, z budową raportu końcowego | TalkIn, WorkSpace |
| **Library** | Centralne repozytorium wiedzy i plików — trwała pamięć zewnętrzna wspólna dla środowisk | TalkIn, WorkSpace |
| **Translate** | Wielojęzyczne tłumaczenia w trybie wielozadaniowym, z zarządzaniem terminologią | TalkIn |
| **Roundtable** | Współpraca wielu modeli AI nad wspólnym problemem: debata, krytyka, konsensus | TalkIn, WorkSpace, CodeStudio |
| **Design** | Projektowanie i tworzenie zasobów wizualnych — od ilustracji po makiety interfejsu | WorkSpace, CodeStudio |
| **Assistant** | Naturalna komunikacja głosowa z AI i realizacja poleceń wieloetapowych | TalkIn |
| **Terminal** | Praca z konsolami i środowiskami wykonawczymi bez opuszczania platformy | CodeStudio |
| **Developer** | Tworzenie i rozwój kodu: edytor, kontrola wersji, wynik budowania, wsparcie AI w kontekście repozytorium | CodeStudio |
| **Diagnostics** | Analiza i usuwanie problemów technicznych: logi, błędy, diagnostyka wydajności | CodeStudio |
| **Apps** | Budowa kompletnych produktów cyfrowych — od architektury, przez frontend i backend, po wdrożenie | WorkSpace, CodeStudio |
| **Agents** | Tworzenie i zarządzanie własnymi agentami AI: tożsamość, umiejętności, uprawnienia | TalkIn, WorkSpace, CodeStudio |

Każdy moduł udostępnia oba kanały komunikacji operacyjnej (Chat Window,
Execution Loop Window) oraz własne okna robocze — na przykład Studio Editor,
Diff/Grep Panel i Session Repository w module Studio; Code Editor, Project Tree,
Git Panel i Build Output w module Developer; Research Workspace, Sources Manager,
Findings Panel i Report Builder w module Research. Pełny katalog okien wszystkich
modułów zawiera Instrukcja użytkowania.

Moduły łączą się ze sobą jawnie i na warunkach ustalonych przez Operatora: źródła
zebrane w Browser mogą zasilić Research, wyniki Studio — trafić do Library,
projekty Workspace — zostać spięte z procesami Automations. Żadne z tych powiązań
nie jest domyślne ani ukryte.

---

## 7. MultitaskingAI — praca ciągła 24/7/365

MultitaskingAI umożliwia równoczesną pracę wielu instancji AI w ramach jednego
projektu, procesu lub zadania. Każdy model pełni odrębną rolę operacyjną
i współpracuje z pozostałymi przez konfigurowalne mechanizmy przekazywania zadań,
wyników i kontekstu. Przedmiotem konfiguracji przestaje być pojedyncza interakcja
z modelem — staje się nim cały zespół.

**Role zespołu:**

| Rola | Odpowiedzialność |
|---|---|
| **Executor 1** — główny wykonawca | Realizuje pracę: dokumenty, analizy, kod, projekty, aplikacje, przetwarzanie materiałów. Nie zarządza procesem |
| **Subagent Network** | Mechanizm uruchamiany przez wykonawcę: do 15 wyspecjalizowanych podagentów (interfejs, backend, API, bezpieczeństwo, baza danych, testy); wyniki agreguje wykonawca |
| **Coordinator** — koordynator | Nie tworzy produktu, lecz projektuje i kontroluje sposób jego wytwarzania: etapy, zależności, harmonogramy, podział pracy, budowa promptów, sterowanie i kolejka |
| **Executor 2** — wykonawca równoległy | Drugi niezależny model wykonawczy: praca niezależna, przekazywanie wyników, praca naprzemienna albo iteracyjna |
| **Executor 3 / Validator** | Rola kontrolna lub doradcza domykająca zespół — wcielenia: Validator, Reviewer, Security Auditor, Architect, Product Owner, QA Lead, Arbitrator |

Rozdzielenie Koordynatora od Wykonawcy jest celowe: żaden model nie musi
jednocześnie planować i realizować.

**Silnik kolejek** pozwala definiować kolejki globalne, lokalne, modeli, agentów
i projektów, z pełnym zestawem akcji: `enqueue`, `dequeue`, `delay`, `retry`,
`pause`, `resume`, `split`, `merge`, `route`, `branch`, `condition`.

**Orkiestracja** pilnuje zależności: zadanie drugiego wykonawcy nie rozpocznie się
przed zakończeniem etapu pierwszego, jeśli tak ustalono, a wynik trafi do właściwej
kolejki albo roli w kolejnym kroku procesu.

**Po spięciu z modułem Automations** powstaje pełny system autonomiczny:

```
                    Always On Display   (nadzór / operator procesu)
                             │
                        Coordinator
              (plan · podział pracy · budowa promptów · sterowanie · kolejka)
                             │
         ┌───────────────────┼───────────────────┐
         ▼                   ▼                   ▼
    Executor 1          Executor 2        Executor 3 / Validator
   (wykonanie)      (wykonanie równoległe)     (kontrola jakości)
         │                   │
  Subagent Network     Subagent Network
  (do 15 podagentów)   (do 15 podagentów)
         │                   │
         └────────►  Silnik kolejek  ◄────────┘
                             │
                       Orkiestracja
                             │
                  Integracja z Automations
          (harmonogram · cykliczność → praca ciągła 24/7/365)
```

Panel orkiestracji, zastępujący w tym środowisku listę modułów, ma sześć sekcji:
**Zespoły** (zapisane presety ról, powiązań i kolejek), **Role** (cztery okna
robocze i przypisanie agentów), **Kolejki**, **Orkiestracja**, **Harmonogram
i automatyki** oraz **Monitor procesu**. Kolejność i widoczność sekcji podlega
konfiguracji.

---

## 8. Agenci i komponenty własne

Komponent własny jest nazwanym wytworem Operatora, tworzonym w oknie konfiguracji
i wpinanym później do pracy tam, gdzie ma zastosowanie.

| Komponent własny | Tworzony w | Przykład zastosowania |
|---|---|---|
| **Automatyka** | Automations | Wpięta w sesji modułu Developer |
| **Agent** | Agents | Wykonawca zadań w dowolnym środowisku, także w roli MultitaskingAI |
| **Projekt** | Workspace | Izolowana przestrzeń projektowa z własnym kontekstem |
| **Profil asystenta** | Assistant | Wykorzystywany operacyjnie w środowisku TalkIn |

Agent powstaje w module Agents: wybór modelu bazowego, definicja tożsamości
i instrukcji systemowych, dodanie umiejętności, wtyczek i konektorów, zarządzanie
pamięcią oraz konfiguracja uprawnień w Centrum uprawnień. Uprawnienia startują
z zakresu minimalnego i rozszerzane są świadomą decyzją Operatora. Raz zbudowany
agent działa ponad wszystkimi środowiskami — nie wymaga definiowania kontekstu
od nowa przy każdym zadaniu.

---

## 9. Funkcje globalne: Mobile i Always On Display

Funkcje globalne nie tworzą własnej przestrzeni roboczej i nie podlegają
przełączaniu — pozostają dostępne w tej samej formie niezależnie od aktywnego
środowiska i modułu.

### Mobile — zarządzanie platformą z dowolnego miejsca

| Możliwość | Zakres |
|---|---|
| Dostęp mobilny | Obsługa platformy z telefonu i tabletu |
| Oba kanały operacyjne | Pełny Chat Window oraz pełny nadzór nad pętlą wykonawczą |
| Monitoring procesów | Projekty, automatyzacje, agenci, aplikacje, sesje, przebiegi MultitaskingAI, kolejki |
| Zarządzanie zadaniami | Uruchamianie, zatrzymywanie, zatwierdzanie, wstrzymywanie, wznawianie, modyfikowanie |
| Decyzje procesu | Zatwierdzanie punktów decyzyjnych, odrzucenie wyniku, decyzja o ponowieniu |
| Powiadomienia wypychane | Odbiór, przegląd i działanie bezpośrednio z powiadomienia |

Mobile odpowiada na potrzebę zachowania kontroli nad procesami długotrwałymi poza
stanowiskiem roboczym. Praca konstrukcyjna — zakładanie sesji, okna operacyjne
modułów, konfiguracja platformy — pozostaje przy stanowisku stacjonarnym;
urządzenie przenośne steruje tym, co już działa.

### Always On Display — osobisty agent operacyjny

| Możliwość | Zakres |
|---|---|
| Proaktywne doradztwo | Propozycje działań, sugestie konfiguracji, wskazywanie problemów, rekomendacje kolejnych kroków |
| Pomoc kontekstowa | Rozumienie aktywnego modułu, środowiska, projektu i historii pracy |
| Komunikacja głosowa i tekstowa | Mowa, synteza głosu, szybkie polecenia, dymki kontekstowe |
| Pływający avatar | Stale obecny element interfejsu; po aktywacji natychmiastowa interakcja |
| Dostęp do platformy | Wszystkie środowiska, projekty, sesje, historie rozmów i agenci Operatora |

W środowisku MultitaskingAI Always On Display pełni dodatkowo rolę obserwatora
albo operatora nadzorującego przebieg złożonych, wieloetapowych procesów.

---

## 10. Modele i kanały integracji

Rdzeń platformy zna jedną abstrakcję: „wyślij zapytanie, odbierz odpowiedź".
Nie wie, czy zapytanie trafia do modelu kluczem dostępu, narzędziem wiersza
poleceń, połączeniem z własnym hostem czy przez interfejs przeglądarkowy.

| Kanał | Sposób połączenia | Właściwe zastosowanie |
|---|---|---|
| **API** | Żądanie HTTP uwierzytelnione kluczem dostępu | Dostawcy z interfejsem programistycznym; podstawowy kanał pracy operacyjnej |
| **CLI** | Wywołanie narzędzia wiersza poleceń dostawcy | Modele udostępniane przede wszystkim jako narzędzie wiersza poleceń |
| **SSH** | Połączenie z własnym modelem pomocniczym na hoście zdalnym | Modele na własnej infrastrukturze; pełna kontrola nad środowiskiem wykonawczym |
| **HTTP** | Interakcja z interfejsem webowym dostawcy | Modele dostępne wyłącznie przez przeglądarkę |

Konsekwencje tej konstrukcji są trzy: rdzeń pozostaje niezależny od kanału, nowy
sposób połączenia oznacza dostarczenie adaptera, a odpowiedź każdego z czterech
kanałów dociera do interfejsu jednolitym strumieniem na żywo.

Model przypisuje się do sesji albo do roli. Kanał, model bazowy, parametry
wywołania i dane dostępowe są ustawieniami konfiguracji — na warstwie globalnej,
środowiska, projektu albo pojedynczej sesji. Dane dostępowe kanałów przechowywane
są poza bazą danych.

---

## 11. Rozszerzenia

Platforma przyjmuje cztery rodzaje rozszerzeń — wtyczki, umiejętności, konektory
i serwery MCP — wszystkie w jednolitym kontrakcie: tożsamość i rejestracja, stan
włączenia zmienialny w każdej chwili, konfiguracja właściwa rodzajowi oraz zakres
uprawnień przypisany przy podłączeniu do agenta lub modułu.

| Cecha | **Danaco Plugin** | **Personal** |
|---|---|---|
| Pochodzenie | Dostarczane z platformą | Instalowane przez Operatora |
| Aktualizacja | Wraz z aktualizacją platformy | Ponowne przesłanie przez Operatora |
| Stan wyjściowy | Włączone | Wyłączone — wymaga świadomego włączenia |

Rozszerzenia biorą udział w warstwie komunikacji operacyjnej i mogą wnosić własne
elementy interfejsu, przypisane do warstw widoczności na tych samych zasadach co
elementy platformy.

---

## 12. Pełna konfigurowalność i punkty izolacji

**Nic nie jest zaszyte na stałe.** Każdą funkcję, każde zachowanie i każdą
zależność między elementami platformy konfiguruje się z okna konfiguracji. To
zasada centralna produktu, nadrzędna wobec pozostałych.

Konfiguracja ma strukturę warstwową — wartość nakłada się od najszerszej do
najwęższej, a warstwa węższa nie zmienia szerszej:

```
[ GLOBALNA ]  cała platforma
      ▼
[ ŚRODOWISKO ]  TalkIn · WorkSpace · CodeStudio · MultitaskingAI
      ▼
[ PROJEKT ]  jeden projekt modułu Workspace
      ▼
[ SESJA ]  jedna karta sesji — dostępna też z menu kontekstowego okna operacyjnego

Brak ustawienia na warstwie = dziedziczenie z warstwy szerszej.
```

Okno konfiguracji obejmuje siedemnaście zakresów: aplikacja, procesy, akcje,
zachowanie modeli, tożsamość modeli, prompty systemowe, rozszerzenia, integracje,
historia, pamięć, izolacja, karty sesji, komponenty własne, Chat Window,
Execution Loop Window, warstwy widoczności funkcji oraz Always On Display. Każda
pozycja ustawienia ma objaśnienie kontekstowe `[?]` mówiące, co robi i jaki ma
wpływ na działanie platformy.

**Punkty izolacji** ustala Operator, nie sztywna reguła systemu. Macierz izolacji
łączy dwie niezależne grupy:

| Grupa | Pozycje | Stany |
|---|---|---|
| Izolacja kontekstu | historia, pamięć, kontekst | współdzielone \| odrębne |
| Izolacja techniczna | katalog roboczy sesji, środowisko procesu, katalog danych i konfiguracji modelu, dostęp sieciowy, zakres odczytu i zapisu plików, konto i token per sesja, model procesu, serwer wykonania | włączony \| wyłączony |

Izolację ustawia się na siedmiu poziomach zasięgu — globalnym, środowiska,
modułu, pary modułów, projektu, karty sesji i roli MultitaskingAI — z regułą
pierwszeństwa poziomu najbardziej szczegółowego. Gotowe układy zapisuje się jako
profile izolacji z podglądem polityki efektywnej.

---

## 13. Interfejs: warstwy widoczności

Zasada nadrzędna interfejsu: **jeżeli funkcja nie jest potrzebna do realizacji
aktualnego zadania, nie jest widoczna.** Złożoność platformy istnieje
w architekturze i pozostaje niewidoczna do chwili, w której funkcja jest potrzebna.
Liczba modułów, agentów, przepływów i ustawień nie wpływa na postrzeganą prostotę
interfejsu.

| Warstwa | Zawartość | Sposób dostępu |
|---|---|---|
| **1 — zawsze widoczna** | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania; ponad 80% powierzchni interfejsu | Bez interakcji |
| **2 — widoczna na żądanie** | Wybór modelu, wykonawcy, środowiska, trybu pracy, poziom wysiłku, parametry przepływu | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się sam |
| **3 — rozwinięcia kontekstowe** | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu `⋮`, menu `☰`, menu kontekstowe, panel, lista rozwijana |
| **4 — funkcje eksperckie** | Operacje zaawansowane, tryby administracyjne, diagnostyka niskiego poziomu | Polecenie języka naturalnego, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny |

**Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem
klawiszowym albo jednym poleceniem w języku naturalnym.** Zagnieżdżanie funkcji
głęboko w menu jest w tym produkcie zabronione — ukrycie zmniejsza chaos wizualny,
nie utrudnia dostępu.

System wizualny marki opiera się na granacie i złocie, ze złotem rezerwowanym dla
stanu aktywnego i najechania. Interfejs dostępny jest w wariancie jasnym
i ciemnym, w wybranym języku.

---

## 14. Architektura i model wdrożenia

Platforma działa w modelu hybrydowym: **przetwarzanie po stronie serwera, dostęp
przez cienkiego klienta** uruchamiającego natywne okno na urządzeniu.

```
   Przetwarzanie, logika, stan   ─────►   SERWER WYKONAWCZY
                                          (aplikacja serwerowa)
   Dostęp, prezentacja, obraz    ─────►   CIENKI KLIENT URZĄDZENIA
                                          (natywne okno na urządzeniu Operatora)
```

| Cecha | Serwer wykonawczy | Cienki klient urządzenia |
|---|---|---|
| Zawartość | Rdzeń platformy i procesy sesji | Lekki pakiet uruchamiający natywne okno |
| Logika biznesowa | Tak | Nie |
| Praca obliczeniowa | Tak | Nie |
| Stan trwały | Tak — jedyne źródło prawdy | Nie |
| Rola w połączeniu | Przyjmuje połączenia klientów | Łączy się z serwerem |

Konsekwencje praktyczne: **stanowisko Operatora nie wymaga lokalnego układu GPU
ani instalowania żadnych dodatkowych programów**, aktualizacja odbywa się
centralnie po stronie serwera, a jedno konto obsługuje wiele urządzeń — komputery,
telefony, tablety.

**Warstwy systemu i technologie:**

| Komponent | Technologia |
|---|---|
| Rdzeń serwera | Go — orkiestracja, współbieżność, zarządzanie procesami sesji |
| Okno klienta | Natywna powłoka aplikacji (warstwa w języku Rust) |
| Interfejs | Aplikacja webowa w TypeScript, wspólna dla wszystkich platform |
| Kanał komunikacji | WebSocket, komunikaty JSON — połączenie stałe i dwukierunkowe |
| Przechowywanie danych | SQLite — trwałość transakcyjna |

Komunikacja klient–serwer biegnie jednym stałym kanałem w trzech kategoriach:
**polecenia** (klient → serwer), **zdarzenia** (serwer → klient) oraz
**strumienie** — odpowiedź modelu przekazywana na żywo w miarę generowania. Tym
samym kanałem przenoszone są komunikaty obu kanałów komunikacji operacyjnej.
Zmiana dokonana na jednym urządzeniu jest rozgłaszana przez serwer do pozostałych;
synchronizacja jest stałą właściwością platformy, nie opcją do włączenia.

---

## 15. Bezpieczeństwo i dostęp

Platforma realizuje model **jednego konta właściciela i wielu urządzeń**.
Rejestracja jest jednorazowym utworzeniem konta przy pierwszym uruchomieniu;
każde kolejne urządzenie dochodzi przez logowanie albo parowanie i otrzymuje
własny token dostępu. Odwołanie dostępu jednemu urządzeniu nie odcina pozostałych.

| Metoda uwierzytelniania | Charakter | Aktywna domyślnie |
|---|---|---|
| Hasło | Bazowa, obecna zawsze | Tak |
| Adres e-mail uwierzytelniający | Weryfikacja, logowanie, odzyskiwanie konta | Tak |
| PIN | Szybkie potwierdzenie na urządzeniu już uwierzytelnionym | Nie |
| Windows Hello | Uwierzytelnienie wiązane z konkretnym urządzeniem | Nie |

Materiał uwierzytelniający — skrót hasła, materiał PIN, materiał Windows Hello —
przechowywany jest poza bazą danych. Dane dostępowe kanałów modeli również leżą
poza bazą. Kanały komunikacji operacyjnej mają własną autoryzację: granice
autonomii Wykonawcy, punkty obowiązkowego zatwierdzenia przez Użytkownika przed
działaniem nieodwracalnym oraz dziennik audytu. Warstwy widoczności funkcji pełnią
przy tym rolę mechanizmu kontroli dostępu — profile warstw przypisuje się do roli
użytkownika, a tryb administracyjny odsłania funkcje eksperckie.

---

## 16. Wymagania

| Element | Wymaganie |
|---|---|
| Urządzenie Operatora | Komputer, telefon albo tablet zdolny wyświetlić natywne okno klienta |
| Układ GPU na stanowisku | Niewymagany |
| Dodatkowe programy na stanowisku | Żadne — pakiet kliencki jest samowystarczalny |
| Połączenie | Dostęp sieciowy do serwera platformy |
| Instalacja klienta | Jednorazowa na urządzenie |
| Liczba urządzeń | Bez ograniczenia |

---

## 17. Zasady nadrzędne platformy

| Zasada | Istota |
|---|---|
| **Pełna kompozycyjność** | W każdym module dostępne są wszystkie funkcje i ich kompozycje: prompt, zadanie, akcja, orkiestracja, subagenci, pełna pętla z Koordynatorem, sekwencje. Żaden moduł nie ogranicza funkcji do wąskiego podzbioru |
| **Pełna konfigurowalność** *(zasada centralna)* | Każdą funkcję, zachowanie i zależność konfiguruje się z okna konfiguracji; nic nie jest zaszyte na stałe poza zasięgiem Operatora |
| **Jawność i konfigurowalność zależności** | Powiązania między funkcjami są jawne i ustanawiane przez Operatora — platforma nie narzuca gotowego zestawu ukrytych zależności |
| **Pełna orkiestracja** | Platforma prowadzi równoległą pracę wielu modeli i agentów oraz udostępnia pełny zestaw funkcji operacyjnych sterujących tą pracą |
| **Warstwy widoczności interfejsu** | Funkcja niepotrzebna do bieżącego zadania jest niewidoczna; każda ukryta pozostaje osiągalna jednym krokiem |

---

## 18. Dokumentacja produktu

| Dokument | Zawartość |
|---|---|
| **`INSTALACJA-I-KONFIGURACJA.md`** | Instalacja pakietu klienckiego, pierwsze uruchomienie, rejestracja i logowanie, podłączanie kolejnych urządzeń, okno Ustawień, okno konfiguracji, kanały i konta modeli, rozszerzenia, izolacja, diagnostyka połączenia |
| **`INSTRUKCJA-UZYTKOWANIA.md`** | Praca z platformą: strona główna, nawigacja, karty sesji, oba kanały komunikacji operacyjnej, środowiska, wszystkie moduły, MultitaskingAI, komponenty własne, funkcje globalne, scenariusze pracy |
| **`LICENSE.md`** | Warunki licencyjne, prawa producenta i twórcy, zakres dozwolonego użycia |
| **`CHANGELOG.md`** | Historia zmian wersji produktu |

---

## 19. Słowniczek

| Pojęcie | Znaczenie |
|---|---|
| **Operator** | Właściciel konta prowadzący pracę na platformie |
| **Środowisko** | Najwyższy poziom organizacji pracy — odpowiada na pytanie „w jakim trybie pracuję?" |
| **Moduł** | Wyspecjalizowany obszar roboczy wewnątrz środowiska — „jakie zadanie wykonuję?" |
| **Okno operacyjne** | Właściwa przestrzeń wykonania zadania wewnątrz modułu |
| **Karta sesji** | Jednostka pracy równoległej w oknie środowiska, z własnym układem, historią i kontekstem |
| **Komponent własny** | Nazwany wytwór Operatora: automatyka, agent, projekt albo profil asystenta |
| **Chat Window** | Okno kanału Użytkownik ↔ Wykonawca — centralny punkt sterowania platformą |
| **Execution Loop Window** | Okno kanału Koordynator ↔ Wykonawca — pętla wykonawcza i nadzór |
| **Wykonawca** | Model AI, agent albo system wykonawczy realizujący zadania |
| **Koordynator** | Komponent orkiestrujący: dekompozycja zlecenia, przydział zadań, nadzór pętli |
| **Punkt izolacji** | Ustawienie rozstrzygające, co jest współdzielone, a co odrębne między przestrzeniami pracy |
| **Warstwa widoczności** | Poziom, na którym element interfejsu jest ujawniany Operatorowi (1–4) |
| **Kanał modelu** | Skonfigurowany sposób osiągnięcia modelu: API, CLI, SSH albo HTTP |

---

## Kontakt

| | |
|---|---|
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Adres** | ul. Gen. J. Hallera 81, 43-400 Cieszyn |
| **Twórca** | Dariusz Naharnowicz |
| **Wsparcie i zgłoszenia** | support@danaco-group.pl |

---

*Danaco Console — Platforma AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz*
