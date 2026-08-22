# Danaco Console — Model danych

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
| **Tytuł** | Model danych |
| **Klasa dokumentu** | Specyfikacja docelowa |
| **Odbiorcy** | deweloper (główny) |
| **Przeznaczenie** | Ustala pełny model danych platformy zgodny ze schematem przechowywania SQLite — grupy encji, ich pola, typy, relacje i ograniczenia — realizujący zasadę pełnej konfigurowalności platformy |
| **Zakres** | konto i urządzenia, środowiska i moduły, strona główna i jej trzy strefy, karty sesji, komponenty własne, konfiguracja i warstwowość ustawień, tożsamość i kanały modeli, punkty dostępu, punkty izolacji, rozszerzenia, pamięć, automatyki i orkiestracja, dane modułowe, zgodność ze schematem SQLite |
| **Poza zakresem** | kontrakty komunikacji przenoszące te dane między klientem a rdzeniem — [Kontrakty komunikacji](kontrakty-komunikacji.md); model konfiguracji jako mechanizm warstwowy — [Model konfiguracji](model-konfiguracji.md); wygląd i zachowanie interfejsu — [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) |
| **Dokument nadrzędny** | [Architektura techniczna](architektura.md) |
| **Dokumenty powiązane** | [Koncepcja platformy](koncepcja-platformy.md) · [Kontrakty komunikacji](kontrakty-komunikacji.md) · [Model konfiguracji](model-konfiguracji.md) · [Bezpieczeństwo i uwierzytelnianie](bezpieczenstwo-i-uwierzytelnianie.md) · [Izolacja i zależności](izolacja-i-zaleznosci.md) · [Integracja modeli](integracja-modeli.md) · [Rozszerzenia](rozszerzenia.md) |
| **Źródła normatywne** | `budowa/shared/contract.json` (552 struktury · 377 wyliczeń) — schemat bazy danych SQLite odwzorowany z pól struktur oznaczonych kolumną bazy |
| **Zasada nadrzędna** | Każda encja modelu ma dokładnie jedno miejsce definicji w tym dokumencie; struktury kontraktu i model bazy danych pozostają zgodne — pole struktury bez odpowiednika w bazie jest jawnie oznaczone jako przelotowe |

Dokument opisuje pełny model danych produktu Danaco Console. Model jest zgodny ze schematem przechowywania SQLite przyjętym w rozdziale 4 dokumentu Architektury i realizuje zasadę pełnej konfigurowalności platformy. Poniższa tabela wskazuje grupy encji objęte modelem wraz z rozdziałem, w którym każda grupa jest opisana.

| Grupa encji | Rozdział niniejszego dokumentu |
|---|---|
| Konto i urządzenia | 3 |
| Środowiska i moduły | 4 |
| Strona główna i jej trzy strefy | 5 |
| Komponenty własne (wzorzec nadrzędny) | 6 |
| Projekty | 7 |
| Karty sesji i sesje | 8 |
| Pamięć wielopoziomowa | 9 |
| Zadania, kolejki, harmonogram i orkiestracja | 10 |
| Kanały modeli | 11 |
| Agenci | 12 |
| Role środowiska MultitaskingAI | 13 |
| Rozszerzenia | 14 |
| Profile izolacji i punkty izolacji | 15 |
| Artefakty i biblioteka | 16 |
| Konfiguracja warstwowa | 17 |
| Funkcje globalne — dane wspierające | 18 |
| Warstwa komunikacji operacyjnej | 19 |
| Konfiguracja warstw widoczności | 20 |

---

## Spis treści

1. [Konwencje modelu](#1-konwencje-modelu)
   - [1.1 Zasada nadrzędna](#11-zasada-nadrzędna)
   - [1.2 Wzorzec komponentu własnego](#12-wzorzec-komponentu-własnego)
   - [1.3 Typy danych i ich odwzorowanie w SQLite](#13-typy-danych-i-ich-odwzorowanie-w-sqlite)
   - [1.4 Nazewnictwo](#14-nazewnictwo)
   - [1.5 Bezpieczeństwo danych dostępowych](#15-bezpieczeństwo-danych-dostępowych)
2. [Mapa encji](#2-mapa-encji)
3. [Tożsamość i dostęp](#3-tożsamość-i-dostęp)
   - [3.1 Konto (`konto`)](#31-konto-konto)
   - [3.2 Urządzenie (`urzadzenie`)](#32-urządzenie-urzadzenie)
   - [3.3 Połączenie (`polaczenie`)](#33-połączenie-polaczenie)
   - [3.4 Poświadczenie urządzenia (`poswiadczenie_urzadzenia`)](#34-poświadczenie-urządzenia-poswiadczenie_urzadzenia)
   - [3.5 Kod parowania (`kod_parowania`)](#35-kod-parowania-kod_parowania)
   - [3.6 Działanie lokalne (`dzialanie_lokalne`)](#36-działanie-lokalne-dzialanie_lokalne)
4. [Środowiska i moduły](#4-środowiska-i-moduły)
   - [4.1 Środowisko (`srodowisko`)](#41-środowisko-srodowisko)
   - [4.2 Moduł (`modul`)](#42-moduł-modul)
   - [4.3 Macierz dostępności modułów (`srodowisko_modul`)](#43-macierz-dostępności-modułów-srodowisko_modul)
   - [4.4 Okno operacyjne (`okno_operacyjne`)](#44-okno-operacyjne-okno_operacyjne)
   - [4.5 Stan okna w sesji (`sesja_okno`)](#45-stan-okna-w-sesji-sesja_okno)
5. [Strona główna i centrum dowodzenia](#5-strona-główna-i-centrum-dowodzenia)
   - [5.1 Strefa strony głównej (`strefa_glowna`)](#51-strefa-strony-głównej-strefa_glowna)
   - [5.2 Pozycja strefy (`pozycja_strefy`)](#52-pozycja-strefy-pozycja_strefy)
6. [Komponenty własne](#6-komponenty-własne)
   - [6.1 Komponent własny (`komponent_wlasny`)](#61-komponent-własny-komponent_wlasny)
   - [6.2 Automatyka (`automatyka`)](#62-automatyka-automatyka)
   - [6.3 Profil asystenta (`profil_asystenta`)](#63-profil-asystenta-profil_asystenta)
   - [6.4 Powiązanie komponentu (`powiazanie_komponentu`)](#64-powiązanie-komponentu-powiazanie_komponentu)
   - [6.5 Przypisanie agentów do projektu (`projekt_agent`)](#65-przypisanie-agentów-do-projektu-projekt_agent)
7. [Projekty](#7-projekty)
   - [7.1 Projekt (`projekt`)](#71-projekt-projekt)
   - [7.2 Widoczność projektu w środowiskach (`projekt_srodowisko`)](#72-widoczność-projektu-w-środowiskach-projekt_srodowisko)
8. [Karty sesji, sesje i procesy](#8-karty-sesji-sesje-i-procesy)
   - [8.1 Karta sesji (`karta_sesji`)](#81-karta-sesji-karta_sesji)
   - [8.2 Sesja (`sesja`)](#82-sesja-sesja)
   - [8.3 Proces sesji (`proces_sesji`)](#83-proces-sesji-proces_sesji)
   - [8.4 Wiadomość (`wiadomosc`)](#84-wiadomość-wiadomosc)
   - [8.5 Załącznik wiadomości (`zalacznik_wiadomosci`)](#85-załącznik-wiadomości-zalacznik_wiadomosci)
   - [8.6 Retencja historii](#86-retencja-historii)
9. [Pamięć wielopoziomowa](#9-pamięć-wielopoziomowa)
   - [9.1 Zasób pamięci (`zasob_pamieci`)](#91-zasób-pamięci-zasob_pamieci)
   - [9.2 Konfiguracja pamięci sesji (`konfiguracja_pamieci_sesji`)](#92-konfiguracja-pamięci-sesji-konfiguracja_pamieci_sesji)
10. [Zadania, kolejki, harmonogram i orkiestracja](#10-zadania-kolejki-harmonogram-i-orkiestracja)
   - [10.1 Zadanie (`zadanie`)](#101-zadanie-zadanie)
   - [10.2 Komponent kompozycji (`komponent_kompozycji`)](#102-komponent-kompozycji-komponent_kompozycji)
   - [10.3 Kolejka (`kolejka`)](#103-kolejka-kolejka)
   - [10.4 Pozycja kolejki (`pozycja_kolejki`)](#104-pozycja-kolejki-pozycja_kolejki)
   - [10.5 Log akcji kolejki (`log_akcji_kolejki`)](#105-log-akcji-kolejki-log_akcji_kolejki)
   - [10.6 Harmonogram (`harmonogram`)](#106-harmonogram-harmonogram)
   - [10.7 Zależność orkiestracji (`zaleznosc_orkiestracji`)](#107-zależność-orkiestracji-zaleznosc_orkiestracji)
11. [Kanały modeli](#11-kanały-modeli)
   - [11.1 Kanał modelu (`kanal_modelu`)](#111-kanał-modelu-kanal_modelu)
12. [Agenci](#12-agenci)
   - [12.1 Agent (`agent`)](#121-agent-agent)
   - [12.2 Umiejętność agenta (`agent_umiejetnosc`)](#122-umiejętność-agenta-agent_umiejetnosc)
   - [12.3 Rozszerzenie agenta (`agent_rozszerzenie`)](#123-rozszerzenie-agenta-agent_rozszerzenie)
   - [12.4 Uprawnienie agenta (`agent_uprawnienie`)](#124-uprawnienie-agenta-agent_uprawnienie)
13. [Środowisko MultitaskingAI](#13-środowisko-multitaskingai)
   - [13.1 Zespół (`zespol_multitasking`)](#131-zespół-zespol_multitasking)
   - [13.2 Rola MultitaskingAI (`rola_multitasking`)](#132-rola-multitaskingai-rola_multitasking)
   - [13.3 Subagent (`subagent`)](#133-subagent-subagent)
   - [13.4 Sekcja panelu orkiestracji (`sekcja_panelu_orkiestracji`)](#134-sekcja-panelu-orkiestracji-sekcja_panelu_orkiestracji)
   - [13.5 Przebieg procesu (`przebieg_multitasking`)](#135-przebieg-procesu-przebieg_multitasking)
   - [13.6 Decyzja procesu (`decyzja_procesu`)](#136-decyzja-procesu-decyzja_procesu)
14. [Rozszerzenia](#14-rozszerzenia)
   - [14.1 Rozszerzenie (`rozszerzenie`)](#141-rozszerzenie-rozszerzenie)
15. [Profile izolacji i punkty izolacji](#15-profile-izolacji-i-punkty-izolacji)
   - [15.1 Poziom zasięgu (`poziom_zasiegu`)](#151-poziom-zasięgu-poziom_zasiegu)
   - [15.2 Profil izolacji (`profil_izolacji`)](#152-profil-izolacji-profil_izolacji)
   - [15.3 Reguła izolacji kontekstu (`regula_izolacji_kontekstu`)](#153-reguła-izolacji-kontekstu-regula_izolacji_kontekstu)
   - [15.4 Reguła izolacji technicznej (`regula_izolacji_technicznej`)](#154-reguła-izolacji-technicznej-regula_izolacji_technicznej)
   - [15.5 Przypisanie profilu izolacji (`przypisanie_profilu_izolacji`)](#155-przypisanie-profilu-izolacji-przypisanie_profilu_izolacji)
   - [15.6 Para modułów (`para_modulow`)](#156-para-modułów-para_modulow)
16. [Artefakty i biblioteka](#16-artefakty-i-biblioteka)
   - [16.1 Artefakt (`artefakt`)](#161-artefakt-artefakt)
   - [16.2 Wersja artefaktu (`wersja_artefaktu`)](#162-wersja-artefaktu-wersja_artefaktu)
   - [16.3 Kolekcja (`kolekcja`)](#163-kolekcja-kolekcja)
   - [16.4 Przypisanie kolekcji (`kolekcja_artefakt`)](#164-przypisanie-kolekcji-kolekcja_artefakt)
   - [16.5 Etykieta artefaktu (`etykieta_artefaktu`)](#165-etykieta-artefaktu-etykieta_artefaktu)
17. [Konfiguracja warstwowa](#17-konfiguracja-warstwowa)
   - [17.1 Ustawienie (`ustawienie`)](#171-ustawienie-ustawienie)
   - [17.2 Objaśnienie kontekstowe jako część definicji](#172-objaśnienie-kontekstowe-jako-część-definicji)
18. [Funkcje globalne — dane wspierające](#18-funkcje-globalne--dane-wspierające)
   - [18.1 Mobile](#181-mobile)
   - [18.2 Sugestia Always On Display (`sugestia_aod`)](#182-sugestia-always-on-display-sugestia_aod)
   - [18.3 Polecenie głosowe (`polecenie_glosowe`)](#183-polecenie-głosowe-polecenie_glosowe)
   - [18.4 Powiadomienie (`powiadomienie`)](#184-powiadomienie-powiadomienie)
19. [Warstwa komunikacji operacyjnej](#19-warstwa-komunikacji-operacyjnej)
   - [19.1 Rozmowa (`rozmowa`)](#191-rozmowa-rozmowa)
   - [19.2 Wiadomość rozmowy (`wiadomosc`)](#192-wiadomość-rozmowy-wiadomosc)
   - [19.3 Zlecenie (`zlecenie`)](#193-zlecenie-zlecenie)
   - [19.4 Zadanie i kolejka zadań](#194-zadanie-i-kolejka-zadań)
   - [19.5 Przebieg pętli wykonawczej (`przebieg_petli`)](#195-przebieg-pętli-wykonawczej-przebieg_petli)
   - [19.6 Komunikat sterujący (`komunikat_sterujacy`)](#196-komunikat-sterujący-komunikat_sterujacy)
   - [19.7 Wynik kontroli jakości (`wynik_kontroli_jakosci`)](#197-wynik-kontroli-jakości-wynik_kontroli_jakosci)
   - [19.8 Ponowienie (`ponowienie`)](#198-ponowienie-ponowienie)
   - [19.9 Indeksy warstwy komunikacji operacyjnej](#199-indeksy-warstwy-komunikacji-operacyjnej)
   - [19.10 Retencja danych komunikacji operacyjnej](#1910-retencja-danych-komunikacji-operacyjnej)
20. [Konfiguracja warstw widoczności](#20-konfiguracja-warstw-widoczności)
   - [20.1 Warstwa widoczności (`warstwa_widocznosci`)](#201-warstwa-widoczności-warstwa_widocznosci)
   - [20.2 Element interfejsu (`element_interfejsu`)](#202-element-interfejsu-element_interfejsu)
   - [20.3 Rola użytkownika (`rola_uzytkownika`)](#203-rola-użytkownika-rola_uzytkownika)
   - [20.4 Konfiguracja warstwy roli (`konfiguracja_warstwy_roli`)](#204-konfiguracja-warstwy-roli-konfiguracja_warstwy_roli)
   - [20.5 Funkcja platformy (`funkcja_platformy`)](#205-funkcja-platformy-funkcja_platformy)
   - [20.6 Indeksy i retencja grupy warstw widoczności](#206-indeksy-i-retencja-grupy-warstw-widoczności)
21. [Zbiorcze powiązania encji](#21-zbiorcze-powiązania-encji)
22. [Zgodność ze schematem SQLite](#22-zgodność-ze-schematem-sqlite)
   - [22.1 Klucze i tryb pracy](#221-klucze-i-tryb-pracy)
   - [22.2 Indeksy](#222-indeksy)
   - [22.3 Pola JSON](#223-pola-json)
   - [22.4 Odniesienia polimorficzne](#224-odniesienia-polimorficzne)
   - [22.5 Tabele słownikowe a tabele danych Operatora](#225-tabele-słownikowe-a-tabele-danych-operatora)
23. [Kryteria odbioru](#23-kryteria-odbioru)
24. [Załącznik A. Słownik enumeracji](#załącznik-a-słownik-enumeracji)
25. [Załącznik B. Pełny wykaz tabel](#załącznik-b-pełny-wykaz-tabel)
26. [Załącznik C. Komponenty własne a moduły platformowe](#załącznik-c-komponenty-własne-a-moduły-platformowe)
27. [Załącznik D. Diagram powiązań encji](#załącznik-d-diagram-powiązań-encji)
28. [Załącznik E. Szablony konfiguracji pól JSON](#załącznik-e-szablony-konfiguracji-pól-json)
29. [Załącznik F. Scenariusze danych](#załącznik-f-scenariusze-danych)
   - [F.1. Autonomiczna budowa aplikacji w środowisku MultitaskingAI](#f1-autonomiczna-budowa-aplikacji-w-środowisku-multitaskingai)
   - [F.2. Przepływ badawczy między modułami](#f2-przepływ-badawczy-między-modułami)
   - [F.3. Konfiguracja punktu izolacji między projektami](#f3-konfiguracja-punktu-izolacji-między-projektami)
   - [F.4. Praca dwujęzyczna Studio i Translate](#f4-praca-dwujęzyczna-studio-i-translate)

---

## 1. Konwencje modelu

### 1.1. Zasada nadrzędna

Model danych realizuje zasadę pełnej konfigurowalności: każdy atrybut zachowania platformy jest konfigurowalny z poziomu okna konfiguracji, a brak ustawienia oznacza wartość domyślną, nie brak możliwości. W modelu danych zasada ta wyraża się w dwóch mechanizmach powtarzających się w wielu encjach.

| Mechanizm | Encje nośne | Reguła |
|---|---|---|
| Dziedziczenie warstwowe | `ustawienie` (rozdz. 17), `profil_izolacji` (rozdz. 15) | Wspólny słownik poziomów zasięgu (`poziom_zasiegu`, rozdz. 15.1) i jedna reguła rozstrzygania: pierwszeństwo poziomu najbardziej szczegółowego, dziedziczenie z poziomu szerszego przy braku ustawienia, poziom globalny jako warstwa bazowa. |
| Jawność powiązań | `powiazanie_komponentu`, `zaleznosc_orkiestracji` | Połączenia między modułami, komponentami własnymi i środowiskiem MultitaskingAI — stanowiące konfigurowalne powiązania, a nie reguły wbudowane — są rekordami w tabelach powiązań, nigdy logiką zaszytą w kodzie. |

### 1.2. Wzorzec komponentu własnego

Cztery rodzaje komponentów własnych — automatyka, agent, projekt, profil asystenta — dzielą wspólny zestaw atrybutów oraz zestaw atrybutów swoistych dla każdego rodzaju. Model realizuje to jako wzorzec nadrzędna/podrzędna, przedstawiony niżej.

```
komponent_wlasny   (atrybuty wspólne: nazwa, przeznaczenie, miejsce konfiguracji,
     │              widoczność, powiązania)
     │  rozróżnienie: kolumna rodzaj ∈ {automatyka, agent, projekt, profil_asystenta}
     │  relacja jeden-do-jednego: kolumna komponent_wlasny_id z ograniczeniem UNIQUE
     ├──> automatyka           (rozdz. 6.2)
     ├──> agent                (rozdz. 12.1)
     ├──> projekt              (rozdz. 7.1)
     └──> profil_asystenta     (rozdz. 6.3)
```

Encja `komponent_wlasny` przechowuje atrybuty wspólne, a każdy rodzaj ma odrębną tabelę szczegółów powiązaną relacją jeden-do-jednego, każda z kolumną `komponent_wlasny_id` objętą ograniczeniem `UNIQUE`. Wzorzec ten jest opisany w rozdziale 6 i wykorzystywany w rozdziałach 7 (Projekty) i 12 (Agenci).

### 1.3. Typy danych i ich odwzorowanie w SQLite

| Kategoria pola | Typ SQLite | Konwencja |
|---|---|---|
| Identyfikator | `INTEGER PRIMARY KEY AUTOINCREMENT` | Alias `rowid`; jedyny wzorzec klucza głównego w modelu (rozdział 22.1). |
| Tekst | `TEXT` | Kodowanie UTF‑8, zgodnie z domyślnym trybem SQLite. |
| Liczba całkowita | `INTEGER` | Liczniki, kolejność, pierwszeństwo. |
| Liczba rzeczywista | `REAL` | Stosowana wyjątkowo — wagi w regułach warunkowych kolejki. |
| Wartość logiczna | `INTEGER` z `CHECK (pole IN (0,1))` | SQLite nie ma natywnego typu logicznego. |
| Data i czas | `TEXT` w formacie ISO 8601 (`RRRR-MM-DDTGG:MM:SSZ`) | SQLite nie ma natywnego typu daty; format tekstowy pozwala na porównania leksykograficzne. |
| Dane ustrukturyzowane zmiennej postaci | `TEXT` (JSON) | Zgodne ze składnią funkcji JSON1 SQLite; walidacja `json_valid()`. |
| Wyliczenie (enum) | `TEXT` z `CHECK` na dozwolony zbiór wartości | Pełny słownik wyliczeń w Załączniku A. |
| Klucz obcy | `INTEGER REFERENCES tabela(id)` | Reguła kasowania (`CASCADE` / `SET NULL` / `RESTRICT`) podana przy każdej relacji w rozdziale 21. |
| Odniesienie polimorficzne | para kolumn `typ_odniesienia TEXT` + `odniesienie_id INTEGER` | Stosowana tam, gdzie jedna kolumna musi wskazywać wiersze różnych tabel — poziomy izolacji oraz cele powiązań. Egzekwowana na poziomie aplikacji, nie ograniczeniem `FOREIGN KEY`. |

### 1.4. Nazewnictwo

| Reguła nazewnicza | Postać przyjęta w modelu |
|---|---|
| Język | Polski. |
| Format nazw tabel i pól | `snake_case`. |
| Liczba w nazwie encji | Pojedyncza (`sesja`, nie `sesje`). |
| Zapis w treści dokumentu | Encja z wielkiej litery („Sesja”); nazwa tabeli czcionką maszynową (`sesja`). |
| Skróty literowo‑cyfrowe | Nieużywane w nazwach ani w opisach — zgodnie z regulaminem redakcyjnym projektu. |

### 1.5. Bezpieczeństwo danych dostępowych

Zgodnie z architekturą (rozdział 15) dane dostępowe nie są przechowywane w bazie SQLite w postaci jawnej — baza przechowuje wyłącznie skrót lub odwołanie do wartości utrzymywanej poza bazą.

| Dana wrażliwa | Co przechowuje baza SQLite | Miejsce wartości jawnej |
|---|---|---|
| Klucze API, tokeny CLI, dane SSH kanałów modeli | Odwołanie (`kanal_modelu.dane_dostepowe_odwolanie`) | Magazyn sekretów po stronie serwera |
| Token urządzenia | Odwołanie (`urzadzenie.token_dostepu_odwolanie`) | Poza bazą |
| Poświadczenie urządzenia przenośnego | Odwołanie (`poswiadczenie_urzadzenia.wartosc_odwolanie`) | Poza bazą |
| Kod parowania urządzenia | Odwołanie (`kod_parowania.wartosc_odwolanie`) | Poza bazą |
| Token połączenia WebSocket | Skrót lub odwołanie | Poza bazą |
| Hasło i PIN konta | Skrót (`konto.haslo_hash`, `konto.pin_hash`) | Nigdy w postaci jawnej |

---

## 2. Mapa encji

Poniższy schemat porządkuje encje modelu według warstw architektury platformy i wskazuje rozdział niniejszego dokumentu, w którym każda grupa jest opisana.

```
KONTO / URZĄDZENIE / POŁĄCZENIE                                    (rozdz. 3)
│     POŚWIADCZENIE URZĄDZENIA / KOD PAROWANIA / DZIAŁANIE LOKALNE
│
├── ŚRODOWISKO ◄──► MODUŁ  (macierz dostępności modułów)            (rozdz. 4)
│         │
│         └── OKNO OPERACYJNE (katalog okien per moduł)
│
├── STREFA STRONY GŁÓWNEJ → POZYCJA STREFY                          (rozdz. 5)
│
├── KOMPONENT WŁASNY  (wzorzec nadrzędny)                           (rozdz. 6)
│     ├── AUTOMATYKA
│     ├── AGENT                                                     (rozdz. 12)
│     ├── PROJEKT                                                   (rozdz. 7)
│     └── PROFIL ASYSTENTA
│
├── KARTA SESJI → SESJA → PROCES SESJI → ROZMOWA → WIADOMOŚĆ        (rozdz. 8)
│
├── ZASÓB PAMIĘCI (poziomy: globalna/środowisko/moduł/projekt/sesja) (rozdz. 9)
│
├── ZADANIE → KOLEJKA → POZYCJA KOLEJKI → HARMONOGRAM → ORKIESTRACJA (rozdz. 10)
│
├── KANAŁ MODELU                                                    (rozdz. 11)
│
├── ŚRODOWISKO MULTITASKINGAI
│     ├── ZESPÓŁ (preset)
│     ├── ROLA MULTITASKINGAI → SUBAGENT
│     ├── SEKCJA PANELU ORKIESTRACJI
│     └── PRZEBIEG PROCESU → DECYZJA PROCESU                        (rozdz. 13)
│
├── ROZSZERZENIE                                                    (rozdz. 14)
│
├── PROFIL IZOLACJI → REGUŁA IZOLACJI → PRZYPISANIE PROFILU
│     (osiem poziomów zasięgu)                                     (rozdz. 15)
│
├── ARTEFAKT → WERSJA ARTEFAKTU → KOLEKCJA / ETYKIETA                (rozdz. 16)
│
├── USTAWIENIE (warstwy: globalna … rola)                           (rozdz. 17)
│
├── ROZMOWA → WIADOMOŚĆ            (Chat Window: Użytkownik ↔ Wykonawca)
│   ZLECENIE → PRZEBIEG PĘTLI → KOMUNIKAT STERUJĄCY
│                             → WYNIK KONTROLI JAKOŚCI → PONOWIENIE
│             (Execution Loop Window: Koordynator ↔ Wykonawca) (rozdz. 19)
│
└── WARSTWA WIDOCZNOŚCI → ELEMENT INTERFEJSU → KONFIGURACJA WARSTWY ROLI
    FUNKCJA PLATFORMY (rejestr wyszukiwarki funkcji)                (rozdz. 20)

FUNKCJE GLOBALNE — dane wspierające: SUGESTIA AOD, POLECENIE GŁOSOWE  (rozdz. 18)
POWIADOMIENIE (centrum powiadomień, kanał Mobile)                    (rozdz. 18.4)
```

---

## 3. Tożsamość i dostęp

Diagram relacji (klucze obce) grupy tożsamości i dostępu:

```
konto ──1:N──> urzadzenie ──1:N──> polaczenie
  │                │                    │ N:1 (aktywne powiązanie, NULL dopuszczalny)
  │                ├──1:N──> poswiadczenie_urzadzenia
  │                ├──1:N──> kod_parowania
  │                └──1:N──> dzialanie_lokalne
  └──1:N──> karta_sesji  <──────────────┘
```

### 3.1. Konto (`konto`)

Jedyne konto właściciela platformy (Operator), zgodnie z modelem wdrożenia jednego użytkownika opisanym w architekturze (rozdział 2).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator konta. | PK |
| login | TEXT | Login logowania. | UNIQUE, NOT NULL |
| email | TEXT | Adres e‑mail uwierzytelniający, pełniący także funkcję odzyskiwania konta. | UNIQUE, NOT NULL |
| haslo_hash | TEXT | Skrót hasła (nigdy wartość jawna). | NOT NULL |
| pin_hash | TEXT | Skrót kodu PIN, jeśli metoda aktywna. | NULL dopuszczalny |
| windows_hello_aktywny | INTEGER | Czy metoda Windows Hello jest aktywna. | CHECK (0,1) |
| email_logowania_aktywny | INTEGER | Czy logowanie kodem e‑mail jest aktywną metodą dodatkową. | CHECK (0,1) |
| data_utworzenia | TEXT | Data rejestracji. | NOT NULL |
| data_ostatniego_logowania | TEXT | Data ostatniego udanego logowania. | — |

*Uwaga projektowa.* Konto jest właścicielem wszystkich `urzadzenie` oraz pośrednio wszystkich `karta_sesji` (przez urządzenie tworzące kartę). Tryb działania bez kontroli dostępu (architektura, rozdział 15) nie zmienia struktury tabeli — przełącznik trybu działa na poziomie aplikacji, nie schematu danych.

### 3.2. Urządzenie (`urzadzenie`)

Urządzenie połączone z platformą — komputer, telefon lub tablet (architektura, rozdział 2).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator urządzenia. | PK |
| konto_id | INTEGER | Właściciel urządzenia. | FK → konto.id |
| nazwa | TEXT | Nazwa własna urządzenia. | NOT NULL |
| typ | TEXT | `komputer` \| `telefon` \| `tablet`. | CHECK enum |
| token_dostepu_odwolanie | TEXT | Odwołanie do tokenu przechowywanego poza bazą (rozdział 1.5). | NOT NULL |
| tryb_mobile_aktywny | INTEGER | Czy urządzenie korzysta z funkcji globalnej Mobile. | CHECK (0,1) |
| data_polaczenia | TEXT | Data pierwszego połączenia (instalacja pakietu klienta). | NOT NULL |
| data_ostatniego_polaczenia | TEXT | Data ostatniej aktywności. | — |

*Relacje.* Urządzenie inicjuje wiele `polaczenie` w czasie; jedno urządzenie może utrzymywać wiele otwartych `karta_sesji`.

### 3.3. Połączenie (`polaczenie`)

Pojedyncze połączenie WebSocket urządzenia z serwerem, odwzorowujące cykl życia opisany w dokumencie Kontraktów komunikacji (rozdział 2): nawiązanie, uwierzytelnienie, powiązanie z sesją, praca, rozłączenie.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator połączenia. | PK |
| urzadzenie_id | INTEGER | Urządzenie inicjujące połączenie. | FK → urzadzenie.id |
| karta_sesji_id | INTEGER | Karta sesji, z którą połączenie jest aktualnie powiązane. | FK → karta_sesji.id, NULL dopuszczalny |
| wersja_kontraktu | TEXT | Wersja kontraktu komunikacji uzgodniona przy nawiązaniu (Kontrakty, rozdział 9). | NOT NULL |
| data_nawiazania | TEXT | Znacznik czasu otwarcia kanału. | NOT NULL |
| data_rozlaczenia | TEXT | Znacznik czasu zamknięcia kanału; puste dla połączeń aktywnych. | NULL dopuszczalny |

*Uwaga projektowa.* Rozłączenie `polaczenie` nie kończy powiązanej `sesja` — proces sesji działa dalej po stronie serwera (architektura, rozdział 7), zgodnie z zasadą, że stan trwały przechowuje wyłącznie serwer.

### 3.4. Poświadczenie urządzenia (`poswiadczenie_urzadzenia`)

Poświadczenie wydane urządzeniu w protokole parowania (funkcje-globalne/mobile.md, rozdział 9.5). Każde urządzenie ma własne poświadczenie — odwołanie jednego nie odcina pozostałych. Wartość jawna pozostaje poza bazą, zgodnie z rozdziałem 1.5.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator poświadczenia. | PK |
| urzadzenie_id | INTEGER | Urządzenie, któremu poświadczenie wydano. | FK → urzadzenie.id |
| wartosc_odwolanie | TEXT | Odwołanie do wartości przechowywanej poza bazą (rozdział 1.5). | NOT NULL |
| data_wydania | TEXT | Data wydania poświadczenia. | NOT NULL |
| data_waznosci | TEXT | Data wygaśnięcia poświadczenia. | NOT NULL |
| stan | TEXT | `aktywne` \| `wygasle` \| `odwolane`. | CHECK enum |
| data_odwolania | TEXT | Data odwołania dostępu; puste dla poświadczeń nieodwołanych. | NULL dopuszczalny |

*Relacje.* Urządzenie ma w danej chwili dokładnie jedno poświadczenie w stanie `aktywne`; poświadczenia wcześniejsze pozostają w tabeli ze stanem `wygasle` albo `odwolane` jako podstawa dziennika audytu.

### 3.5. Kod parowania (`kod_parowania`)

Kod jednorazowy wydany przy inicjowaniu parowania z okna Ustawień (funkcje-globalne/mobile.md, rozdział 9.4). Wartość jawna pozostaje poza bazą.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator kodu. | PK |
| konto_id | INTEGER | Konto, dla którego kod wydano. | FK → konto.id |
| urzadzenie_inicjujace_id | INTEGER | Urządzenie, z którego zainicjowano parowanie. | FK → urzadzenie.id |
| wartosc_odwolanie | TEXT | Odwołanie do wartości przechowywanej poza bazą (rozdział 1.5). | NOT NULL |
| data_wydania | TEXT | Data wydania kodu. | NOT NULL |
| data_waznosci | TEXT | Termin ważności kodu. | NOT NULL |
| stan | TEXT | `wydany` \| `uzyty` \| `uniewazniony` \| `wygasly`. | CHECK enum |
| urzadzenie_sparowane_id | INTEGER | Urządzenie, które kod wykorzystało; puste do chwili użycia. | FK → urzadzenie.id, NULL dopuszczalny |

### 3.6. Działanie lokalne (`dzialanie_lokalne`)

Pozycja kolejki działań wykonanych na urządzeniu przenośnym bez połączenia, oczekujących na uzgodnienie stanu z serwerem (funkcje-globalne/mobile.md, rozdział 7.2).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator działania. | PK |
| urzadzenie_id | INTEGER | Urządzenie, na którym działanie wykonano. | FK → urzadzenie.id |
| rodzaj | TEXT | `polecenie_czatu` \| `zatwierdzenie` \| `przerwanie` \| `wstrzymanie` \| `wznowienie` \| `ustawienie` \| `powiadomienie`. | CHECK enum |
| przedmiot_typ | TEXT | `sesja` \| `zadanie` \| `przebieg_petli` \| `automatyka` \| `ustawienie` \| `powiadomienie`. | CHECK enum |
| przedmiot_id | INTEGER | Odniesienie polimorficzne do przedmiotu działania. | NULL dopuszczalny |
| tresc | TEXT | Ładunek działania zapisany w chwili wykonania. | NOT NULL |
| wersja_stanu | TEXT | Wersja stanu przedmiotu, na której Użytkownik podejmował decyzję. | NOT NULL |
| data_zapisu | TEXT | Data wykonania działania na urządzeniu. | NOT NULL |
| stan_uzgodnienia | TEXT | `oczekuje` \| `przyjete` \| `bezprzedmiotowe` \| `odrzucone`. | CHECK enum |
| przyczyna_odrzucenia | TEXT | Przyczyna odrzucenia działania w uzgadnianiu. | NULL dopuszczalny |

*Uwaga projektowa.* Serwer pozostaje jedynym źródłem prawdy — kolejka nie jest stanem równoległym, lecz zapisem intencji Użytkownika oczekujących na weryfikację wersji stanu przy najbliższym połączeniu.

---

## 4. Środowiska i moduły

Diagram relacji (klucze obce) grupy środowisk i modułów:

```
srodowisko ──1:N──> srodowisko_modul <──N:1── modul ──1:N──> okno_operacyjne
                                                                  ▲
sesja ──1:N──> sesja_okno ────────────────────────────────────N:1┘
```

### 4.1. Środowisko (`srodowisko`)

Najwyższy poziom organizacji pracy. Tabela słownikowa obejmująca cztery stałe wiersze: TalkIn, WorkSpace, CodeStudio, MultitaskingAI.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator środowiska. | PK |
| kod | TEXT | `TALKIN` \| `WORKSPACE` \| `CODESTUDIO` \| `MULTITASKINGAI`. | UNIQUE, CHECK enum |
| nazwa | TEXT | Nazwa własna środowiska. | NOT NULL |
| opis | TEXT | Krótki opis trybu pracy. | — |
| ma_nawigacje_modulow | INTEGER | `0` wyłącznie dla MultitaskingAI, które ma panel orkiestracji zamiast listy modułów. | CHECK (0,1) |

*Uwaga.* Środowisko MultitaskingAI nie występuje w macierzy `srodowisko_modul` jako kolumna nawigacji modułów — jego boczna nawigacja jest opisana w rozdziale 13.4 (Sekcja panelu orkiestracji).

### 4.2. Moduł (`modul`)

Wyspecjalizowany obszar roboczy wewnątrz środowiska. Tabela słownikowa obejmująca piętnaście stałych wierszy.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator modułu. | PK |
| kod | TEXT | Kod modułu — komplet piętnastu wartości: `STUDIO`, `WORKSPACE`, `AUTOMATIONS`, `BROWSER`, `RESEARCH`, `LIBRARY`, `TRANSLATE`, `ROUNDTABLE`, `DESIGN`, `ASSISTANT`, `TERMINAL`, `DEVELOPER`, `DIAGNOSTICS`, `APPS`, `AGENTS`. | UNIQUE, CHECK enum |
| nazwa | TEXT | Nazwa własna modułu. | NOT NULL |
| cel | TEXT | Opis celu modułu. | — |
| ma_okno_w_nawigacji_srodowisk | INTEGER | `0` dla Automations — moduł nie ma własnego okna w bocznej nawigacji żadnego środowiska. | CHECK (0,1) |
| konfigurowany_na_stronie_glownej | INTEGER | `1` dla modułów będących zarazem rodzajem komponentu własnego (Workspace, Agents) lub konfigurowanych wyłącznie na stronie głównej (Automations, Assistant). | CHECK (0,1) |

**Uwaga terminologiczna.** Wiersz o kodzie `WORKSPACE` w tej tabeli odpowiada modułowi Workspace i jest bytem odrębnym od pojęcia architektonicznego „warstwa modułów” (Workspace Layer), które w niniejszym modelu nie jest osobną encją, lecz określeniem zbioru wszystkich wierszy tabeli `modul`.

### 4.3. Macierz dostępności modułów (`srodowisko_modul`)

Tabela łącząca realizująca macierz dostępności modułów: moduł nie należy do jednego środowiska na wyłączność, lecz jest dostępny w wybranych środowiskach.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator wiersza macierzy. | PK |
| srodowisko_id | INTEGER | Środowisko. | FK → srodowisko.id |
| modul_id | INTEGER | Moduł. | FK → modul.id |
| dostepny | INTEGER | `1` jeśli moduł jest widoczny w bocznej nawigacji tego środowiska. | CHECK (0,1) |

Ograniczenie `UNIQUE(srodowisko_id, modul_id)`. Wartości zasiewające tabelę odpowiadają macierzy dostępności modułów:

| Moduł | TalkIn | WorkSpace | CodeStudio |
|---|---|---|---|
| Studio | TAK | TAK | NIE |
| Workspace | TAK | TAK | TAK |
| Automations | NIE | NIE | NIE |
| Browser | TAK | TAK | NIE |
| Research | TAK | TAK | NIE |
| Library | TAK | TAK | NIE |
| Translate | TAK | NIE | NIE |
| Roundtable | TAK | TAK | TAK |
| Design | NIE | TAK | TAK |
| Assistant | TAK | NIE | NIE |
| Terminal | NIE | NIE | TAK |
| Developer | NIE | NIE | TAK |
| Diagnostics | NIE | NIE | TAK |
| Apps | NIE | TAK | TAK |
| Agents | TAK | TAK | TAK |

Dla Automations wartość „—” macierzy jest odwzorowana jako `dostepny = 0` we wszystkich trzech wierszach (moduł konfigurowany wyłącznie w strefie 2 strony głównej).

### 4.4. Okno operacyjne (`okno_operacyjne`)

Katalog okien roboczych właściwych każdemu modułowi.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator okna. | PK |
| modul_id | INTEGER | Moduł, do którego okno należy. | FK → modul.id |
| nazwa | TEXT | Nazwa własna okna wzięta z Załącznika A Koncepcji platformy — „Studio Editor”, „Workflow Builder”. | NOT NULL |
| kolejnosc | INTEGER | Kolejność prezentacji w module. | NOT NULL |
| czy_okno_wspolne | INTEGER | `1` dla Chat Window oraz Execution Loop Window — dwóch okien kanałów komunikacji operacyjnej wspólnych dla wszystkich piętnastu modułów (Koncepcja platformy, rozdz. 3 i Załącznik A). | CHECK (0,1) |

### 4.5. Stan okna w sesji (`sesja_okno`)

Migawka układu okien danej sesji, pozwalająca odtworzyć rozmieszczenie i stan poszczególnych okien operacyjnych po ponownym połączeniu.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator wiersza. | PK |
| sesja_id | INTEGER | Sesja, do której migawka należy. | FK → sesja.id, ON DELETE CASCADE |
| okno_operacyjne_id | INTEGER | Rodzaj okna. | FK → okno_operacyjne.id |
| widoczne | INTEGER | Czy okno jest aktualnie wyświetlane. | CHECK (0,1) |
| kolejnosc | INTEGER | Kolejność / pozycja okna w bieżącym układzie. | — |
| stan | TEXT | Stan specyficzny dla okna — otwarty plik w Code Editor — zapis JSON. | JSON |

Szablon zawartości pola `stan` (kolumna JSON) dla okna Code Editor:

```json
{
  "otwarty_plik": "rdzen/sesja.go",
  "pozycja_kursora": { "wiersz": 128, "kolumna": 4 },
  "przewiniecie": 1840,
  "panele_zwiniete": ["Tools Panel"]
}
```

---

## 5. Strona główna i centrum dowodzenia

Model odzwierciedla zarówno stały podział na trzy strefy o różnej wadze wizualnej, jak i konfigurowalność kolejności oraz widoczności pozycji wewnątrz każdej strefy, zgodną z zasadą pełnej konfigurowalności.

Diagram relacji (klucze obce) grupy strony głównej:

```
strefa_glowna ──1:N──> pozycja_strefy
pozycja_strefy ──(odniesienie polimorficzne: typ_odniesienia + kod_odniesienia)──>
      srodowisko  |  rodzaj komponentu własnego  |  funkcja ustawień
```

### 5.1. Strefa strony głównej (`strefa_glowna`)

Tabela słownikowa, trzy stałe wiersze.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator strefy. | PK |
| numer | INTEGER | `1` \| `2` \| `3`. | UNIQUE |
| nazwa | TEXT | „Wybór środowiska” \| „Komponenty własne” \| „Ustawienia”. | NOT NULL |
| forma_prezentacji | TEXT | `karta_srodowiska` \| `kafel_komponentu` \| `listwa_ustawien`. | CHECK enum |
| waga_wizualna | TEXT | `glowna` \| `posrednia` \| `najnizsza`. | CHECK enum |

Domyślne wiersze słownika odpowiadają formie prezentacji stref:

| Numer | Nazwa | Forma prezentacji | Waga wizualna |
|---|---|---|---|
| 1 | Wybór środowiska | `karta_srodowiska` | `glowna` |
| 2 | Komponenty własne | `kafel_komponentu` | `posrednia` |
| 3 | Ustawienia | `listwa_ustawien` | `najnizsza` |

### 5.2. Pozycja strefy (`pozycja_strefy`)

Pojedynczy element widoczny w danej strefie — karta środowiska, kafel komponentu własnego lub pozycja listwy ustawień. Kolejność i widoczność są konfigurowalne; zestaw domyślny podano niżej.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator pozycji. | PK |
| strefa_id | INTEGER | Strefa nadrzędna. | FK → strefa_glowna.id |
| typ_odniesienia | TEXT | `srodowisko` (strefa 1) \| `rodzaj_komponentu_wlasnego` (strefa 2) \| `funkcja_ustawien` (strefa 3). | CHECK enum |
| kod_odniesienia | TEXT | Kod celu: kod środowiska, rodzaj komponentu własnego (`automatyka`\|`agent`\|`projekt`\|`profil_asystenta`) lub kod funkcji ustawień (`okno_konfiguracji`\|`mobile`\|`always_on_display`). | NOT NULL |
| etykieta | TEXT | Etykieta widoczna na pozycji; dla strefy 2 zorientowana na działanie twórcze, w brzmieniu „Utwórz automatykę”. | — |
| kolejnosc | INTEGER | Kolejność wyświetlania w obrębie strefy. | NOT NULL |
| widoczna | INTEGER | Widoczność pozycji (konfigurowalna, zestaw domyślny = wszystkie widoczne). | CHECK (0,1) |

*Uwaga projektowa.* `kod_odniesienia` jest odniesieniem symbolicznym, nie kluczem obcym w sensie ścisłym — pozwala jednej tabeli obsłużyć trzy rodzajowo różne cele bez trzech osobnych tabel łączących (rozdział 1.3). Wejście w pozycję strefy 2 otwiera okno konfiguracji właściwe rodzajowi komponentu i tworzy wiersz `komponent_wlasny` (rozdział 6); wejście w pozycję strefy 1 tworzy `karta_sesji` w wybranym środowisku (rozdział 8).

---

## 6. Komponenty własne

Diagram relacji (klucze obce) grupy komponentów własnych:

```
komponent_wlasny ──1:1 (wg kolumny rodzaj)──> automatyka | agent | projekt | profil_asystenta
komponent_wlasny ──(powiazanie_komponentu, N:M)──> modul | komponent_wlasny | srodowisko
automatyka ──N:1──> harmonogram, kolejka
projekt ──1:N──> projekt_agent <──N:1── agent
```

### 6.1. Komponent własny (`komponent_wlasny`)

Encja nadrzędna wzorca opisanego w rozdziale 1.2, realizująca pojęcie komponentu własnego: nazwanego wytworu Operatora, odrębnego od modułu platformy.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator komponentu. | PK |
| nazwa | TEXT | Nazwa własna nadana przez Operatora. | NOT NULL |
| rodzaj | TEXT | `automatyka` \| `agent` \| `projekt` \| `profil_asystenta`. | CHECK enum |
| miejsce_konfiguracji | TEXT | `strefa_2_strony_glownej` \| `okno_modulu`. | CHECK enum |
| przeznaczenie | TEXT | Krótki opis zadania komponentu. | — |
| widocznosc | TEXT | `globalny` \| `projektowy`. | CHECK enum |
| data_utworzenia | TEXT | — | NOT NULL |
| data_modyfikacji | TEXT | — | — |

*Relacje.* Każdy wiersz `komponent_wlasny` ma dokładnie jeden wiersz szczegółowy w tabeli odpowiadającej wartości `rodzaj`: `automatyka` (6.2), `agent` (rozdz. 12.1), `projekt` (rozdz. 7.1) lub `profil_asystenta` (6.3). Powiązania jawne z modułami i innymi komponentami rejestruje `powiazanie_komponentu` (6.4), zgodnie z zasadą jawności i konfigurowalności zależności.

### 6.2. Automatyka (`automatyka`)

Szczegóły komponentu własnego rodzaju „automatyka”.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator. | PK |
| komponent_wlasny_id | INTEGER | Komponent nadrzędny. | FK → komponent_wlasny.id, UNIQUE |
| definicja_workflow | TEXT | Definicja procesu zbudowanego w Workflow Builder. | JSON |
| harmonogram_id | INTEGER | Powiązany harmonogram, jeśli ustanowiony. | FK → harmonogram.id, NULL dopuszczalny |
| kolejka_id | INTEGER | Powiązana kolejka operacyjna. | FK → kolejka.id, NULL dopuszczalny |
| wywolania_modeli | TEXT | Konfiguracja zaplanowanych wywołań modeli. | JSON |

Szablon zawartości pól `definicja_workflow` i `wywolania_modeli` (kolumny JSON):

```json
{
  "definicja_workflow": {
    "kroki": [
      { "nazwa": "Pobranie danych", "okno": "Workflow Builder", "akcja": "odczyt" },
      { "nazwa": "Wywołanie modelu", "akcja": "prompt", "kanal_modelu_id": 3 },
      { "nazwa": "Zapis artefaktu", "akcja": "zapis", "cel": "Library" }
    ]
  },
  "wywolania_modeli": [
    { "kanal_modelu_id": 3, "cyklicznosc": "cyklicznie", "harmonogram_id": 12 }
  ]
}
```

### 6.3. Profil asystenta (`profil_asystenta`)

Szczegóły komponentu własnego rodzaju „profil asystenta”, konfigurowanego na stronie głównej i wykorzystywanego operacyjnie w środowisku TalkIn.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator. | PK |
| komponent_wlasny_id | INTEGER | Komponent nadrzędny. | FK → komponent_wlasny.id, UNIQUE |
| glos | TEXT | Wybrany głos syntezy mowy. | — |
| synteza_mowy_aktywna | INTEGER | Czy synteza mowy jest włączona. | CHECK (0,1) |
| zakres_polecen | TEXT | Dozwolony zakres poleceń głosowych. | JSON |
| obslugiwane_akcje | TEXT | Lista obsługiwanych akcji wieloetapowych. | JSON |

Szablon zawartości pól `zakres_polecen` i `obslugiwane_akcje` (kolumny JSON):

```json
{
  "zakres_polecen": {
    "dozwolone": ["uruchom_zadanie", "wstrzymaj_zadanie", "odczytaj_status"],
    "wymaga_potwierdzenia": ["zatwierdz_przebieg"]
  },
  "obslugiwane_akcje": [
    { "nazwa": "Raport dzienny", "etapy": ["zebranie", "synteza", "wysyłka"] }
  ]
}
```

### 6.4. Powiązanie komponentu (`powiazanie_komponentu`)

Generyczna encja rejestrująca jawne, konfigurowalne powiązania między modułami i komponentami własnymi opisane w specyfikacji modułów jako powiązania jawne.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator powiązania. | PK |
| zrodlo_typ | TEXT | `modul` \| `komponent_wlasny` \| `srodowisko`. | CHECK enum |
| zrodlo_id | INTEGER | Identyfikator elementu źródłowego (odniesienie polimorficzne). | NOT NULL |
| cel_typ | TEXT | `modul` \| `komponent_wlasny` \| `srodowisko`. | CHECK enum |
| cel_id | INTEGER | Identyfikator elementu docelowego (odniesienie polimorficzne). | NOT NULL |
| rodzaj_powiazania | TEXT | Opis powiązania, w brzmieniu „współdzielenie operacji kontekstowych AI”. | NOT NULL |
| aktywne | INTEGER | Czy powiązanie jest aktualnie aktywne. | CHECK (0,1) |
| data_ustanowienia | TEXT | Data ustanowienia powiązania przez Operatora. | NOT NULL |

Przykłady powiązań odwzorowywane tą tabelą:

| Powiązanie | `zrodlo_typ` → `cel_typ` |
|---|---|
| Współdzielenie operacji kontekstowych AI między Studio a Translate | `modul` → `modul` |
| Spięcie automatyki ze środowiskiem MultitaskingAI | `komponent_wlasny` → `srodowisko` |
| Wykorzystanie Library jako repozytorium dla Studio i Research | `modul` → `modul` |

*Uwaga projektowa.* Wiersz tej tabeli jest zawsze skutkiem jawnej decyzji Operatora podjętej w oknie konfiguracji — model nie przewiduje wierszy tworzonych automatycznie ani domyślnie aktywnych bez interwencji Operatora, co odzwierciedla zasadę, że domyślny podział pracy jest wystarczający, a łączenie kontekstów jest możliwością, nie wymaganiem.

### 6.5. Przypisanie agentów do projektu (`projekt_agent`)

Strukturalne (nie swobodne) powiązanie agenta z projektem przez Agent Manager, wydzielone z `powiazanie_komponentu` ze względu na częstość użycia i jednoznaczną semantykę.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator. | PK |
| projekt_id | INTEGER | Projekt. | FK → projekt.id, ON DELETE CASCADE |
| agent_id | INTEGER | Agent. | FK → agent.id |
| rola_w_projekcie | TEXT | Opisowa rola agenta w projekcie. | NULL dopuszczalny |
| data_przypisania | TEXT | — | NOT NULL |

Ograniczenie `UNIQUE(projekt_id, agent_id)`.

---

## 7. Projekty

Diagram relacji (klucze obce) grupy projektów:

```
komponent_wlasny ──1:1──> projekt ──1:N──> projekt_srodowisko ──N:1──> srodowisko
                             │──1:N──> projekt_agent ──N:1──> agent
                             │──1:N──> artefakt, kolekcja        (rozdz. 16)
                             └──N:1──> profil_izolacji (NULL dopuszczalny, rozdz. 15)
```

### 7.1. Projekt (`projekt`)

Szczegóły komponentu własnego rodzaju „projekt” — izolowane środowisko projektowe konfigurowane jako komponent własny na stronie głównej, dostępne dodatkowo jako okno modułowe w środowiskach wskazanych w macierzy dostępności modułów.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator. | PK |
| komponent_wlasny_id | INTEGER | Komponent nadrzędny. | FK → komponent_wlasny.id, UNIQUE |
| opis | TEXT | — | — |
| instrukcje_systemowe | TEXT | Instrukcje właściwe temu projektowi (Instructions Panel). | — |
| katalog_roboczy | TEXT | Ścieżka katalogu roboczego projektu. | — |
| tryb_dostepu_historii | TEXT | `pelny` \| `ograniczony` \| `brak`. | CHECK enum |
| tryb_incognito | INTEGER | Pełne odcięcie historii, gdy aktywne. | CHECK (0,1) |
| profil_izolacji_id | INTEGER | Przypisany profil izolacji (rozdział 15.2). | FK → profil_izolacji.id, NULL dopuszczalny |
| data_utworzenia | TEXT | — | NOT NULL |
| data_modyfikacji | TEXT | — | — |

*Odwzorowanie okien modułu Workspace.* Pamięć kontekstowa projektu to wiersze `zasob_pamieci` z `poziom = 'projekt'` (rozdział 9). Biblioteka projektu to wiersze `artefakt` i `kolekcja` z ustawionym `projekt_id` (rozdział 16), odpowiadające oknu Project Library — funkcjonalnie ograniczonemu odpowiednikowi modułu Library.

### 7.2. Widoczność projektu w środowiskach (`projekt_srodowisko`)

Dowolny zbiór spośród TalkIn, WorkSpace, CodeStudio, w których projekt jest widoczny.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| projekt_id | INTEGER | Projekt. | FK → projekt.id, ON DELETE CASCADE |
| srodowisko_id | INTEGER | Środowisko, w którym projekt jest widoczny. | FK → srodowisko.id |

Ograniczenie `UNIQUE(projekt_id, srodowisko_id)`.

---

## 8. Karty sesji, sesje i procesy

Diagram relacji (klucze obce) grupy kart sesji, sesji i procesów:

```
karta_sesji ──1:N──> sesja ──1:1──> proces_sesji
                       │
                       ├──1:N──> rozmowa ──1:N──> wiadomosc ──1:N──> zalacznik_wiadomosci
                       ├──1:N──> sesja_okno ──N:1──> okno_operacyjne
                       └──0:N──> zadanie (NULL dopuszczalny)
sesja ──N:1──> srodowisko / modul / projekt
```

### 8.1. Karta sesji (`karta_sesji`)

Karta pozioma w oknie środowiska. Karta utrzymuje tożsamość niezależnie od przełączania modułu wewnątrz niej — przełączenie modułu przeładowuje przestrzeń roboczą karty i tworzy nowy wiersz `sesja` powiązany z tą samą kartą.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator karty. | PK |
| konto_id | INTEGER | Właściciel karty. | FK → konto.id |
| srodowisko_id | INTEGER | Aktualnie aktywne środowisko karty. | FK → srodowisko.id |
| nazwa | TEXT | Etykieta widoczna na karcie, nadawana lub generowana. | NULL dopuszczalny |
| kolejnosc | INTEGER | Pozycja karty w poziomym pasku kart. | NOT NULL |
| stan | TEXT | `aktywna` \| `w_tle` \| `zamknieta`. | CHECK enum |
| data_utworzenia | TEXT | — | NOT NULL |
| data_ostatniej_aktywnosci | TEXT | — | — |

*Relacje.* Karta gromadzi w czasie wiele wierszy `sesja` (jeden aktywny, pozostałe historyczne po przełączeniu modułu). Zakres współdzielenia kontekstu między kartami konfiguruje `profil_izolacji` na poziomie zasięgu `karta_sesji` (rozdział 15).

### 8.2. Sesja (`sesja`)

Instancja przestrzeni roboczej modułu (lub, w środowisku MultitaskingAI, zestawu ról) w obrębie karty sesji.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator sesji. | PK |
| karta_sesji_id | INTEGER | Karta, do której sesja należy. | FK → karta_sesji.id, ON DELETE CASCADE |
| srodowisko_id | INTEGER | Środowisko sesji. | FK → srodowisko.id |
| modul_id | INTEGER | Moduł sesji; puste w środowisku MultitaskingAI, gdzie przestrzeń definiują role (rozdział 13), nie moduł. | FK → modul.id, NULL dopuszczalny |
| projekt_id | INTEGER | Projekt, w kontekście którego sesja działa. | FK → projekt.id, NULL dopuszczalny |
| stan | TEXT | `aktywna` \| `wstrzymana` \| `zakonczona`. | CHECK enum |
| kontekst | TEXT | Bieżący kontekst przekazywany do AI. | JSON |
| zapis_do_powrotu | INTEGER | Czy zamknięta sesja jest zapisywana do wznowienia (konfigurowalne z okna konfiguracji). | CHECK (0,1) |
| data_utworzenia | TEXT | — | NOT NULL |
| data_zakonczenia | TEXT | — | NULL dopuszczalny |

*Relacje.* Sesja gromadzi `rozmowa` (rozdział 19.1) wraz z należącymi do nich `wiadomosc` (8.4), ma migawkę układu okien `sesja_okno` (4.5) i dokładnie jeden powiązany `proces_sesji` (8.3). Sesja może być źródłem `zadanie` (rozdział 10.1) oraz celem `zasob_pamieci` na poziomie sesji (rozdział 9).

### 8.3. Proces sesji (`proces_sesji`)

Rejestr aktywnych procesów serwera prowadzony przez rdzeń (architektura, rozdział 7): każda sesja działa jako odrębny, izolowany proces po stronie serwera, z własnym katalogiem roboczym i środowiskiem, niezależny od stanu połączenia klienta.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator procesu. | PK |
| sesja_id | INTEGER | Sesja, którą proces obsługuje. | FK → sesja.id, UNIQUE, ON DELETE CASCADE |
| katalog_roboczy | TEXT | Katalog roboczy procesu. | NOT NULL |
| zmienne_srodowiskowe | TEXT | Środowisko procesu. | JSON |
| stan | TEXT | `uruchomiony` \| `wstrzymany` \| `zakonczony`. | CHECK enum |
| profil_izolacji_efektywny_id | INTEGER | Profil izolacji technicznej faktycznie zastosowany do procesu, po rozstrzygnięciu dziedziczenia (rozdział 15.5). | FK → profil_izolacji.id, NULL dopuszczalny |
| data_uruchomienia | TEXT | — | NOT NULL |
| data_zakonczenia | TEXT | — | NULL dopuszczalny |

*Uwaga projektowa.* Pole `stan` odpowiada zdarzeniu `process.status` z dokumentu Kontraktów komunikacji (rozdział 6). Rozłączenie `polaczenie` (3.3) nie zmienia stanu procesu — zgodnie z architekturą (rozdział 7), rozłączenie klienta nie kończy sesji.

### 8.4. Wiadomość (`wiadomosc`)

Pojedynczy wpis w historii rozmowy prowadzonej w głównym oknie komunikacji (Chat Window, kanał Użytkownik ↔ Wykonawca). Warstwę komunikacji operacyjnej, do której encja należy, opisuje rozdział 19.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator wiadomości. | PK |
| sesja_id | INTEGER | Sesja, do której wiadomość należy. | FK → sesja.id, ON DELETE CASCADE |
| rozmowa_id | INTEGER | Rozmowa, do której wiadomość należy (rozdział 19.1). | FK → rozmowa.id, ON DELETE CASCADE |
| rola | TEXT | `uzytkownik` \| `wykonawca`. | CHECK enum |
| tresc | TEXT | Treść wiadomości. | NOT NULL |
| kanal_modelu_id | INTEGER | Kanał modelu, przez który Wykonawca wygenerował odpowiedź (dla `rola = 'wykonawca'`). | FK → kanal_modelu.id, NULL dopuszczalny |
| znacznik_czasu | TEXT | — | NOT NULL |

### 8.5. Załącznik wiadomości (`zalacznik_wiadomosci`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| wiadomosc_id | INTEGER | Wiadomość, do której załącznik należy. | FK → wiadomosc.id, ON DELETE CASCADE |
| nazwa_pliku | TEXT | — | NOT NULL |
| sciezka_pliku | TEXT | Ścieżka pliku na serwerze; treść obszerna przechowywana jako plik, z odwołaniem w bazie (architektura, rozdział 10). | NOT NULL |
| typ_mime | TEXT | — | — |

### 8.6. Retencja historii

Okres przechowywania historii sesji i rozmów nie jest odrębną encją — jest wierszem `ustawienie` (rozdział 17.1), zgodnie z zasadą „brak ustawienia = wartość domyślna”.

| Zagadnienie | Odwzorowanie w modelu |
|---|---|
| Okres przechowywania historii sesji i rozmów | Wiersz `ustawienie` (rozdział 17.1) o kluczu `retencja_historii`, ustawialny na dowolnym poziomie zasięgu |
| Brak ustawienia | Przechowywanie bez limitu |

Pełny zestaw kluczy retencji obejmujących rozmowy, zlecenia, komunikaty sterujące i wyniki kontroli jakości podano w rozdziale 19.10.

---

## 9. Pamięć wielopoziomowa

Diagram relacji (klucze obce) grupy pamięci:

```
zasob_pamieci ──(odniesienie polimorficzne: poziom + odniesienie_id)──>
      srodowisko | modul | projekt | sesja        (poziom 'globalna' → brak odniesienia)
konfiguracja_pamieci_sesji ──1:1──> sesja
```

### 9.1. Zasób pamięci (`zasob_pamieci`)

Pojedynczy wpis pamięci kontekstowej, przypisany do jednego z pięciu poziomów.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator zasobu. | PK |
| poziom | TEXT | `globalna` \| `srodowisko` \| `modul` \| `projekt` \| `sesja`. | CHECK enum |
| odniesienie_id | INTEGER | Identyfikator elementu wskazanego poziomu (odniesienie polimorficzne); puste dla poziomu `globalna`. | NULL dopuszczalny |
| stan_wlaczenia | INTEGER | Czy zasób jest aktywny. | CHECK (0,1) |
| tresc | TEXT | Treść zapamiętana. | NOT NULL |
| znacznik_czasu | TEXT | Data utworzenia wpisu. | NOT NULL |
| data_modyfikacji | TEXT | — | — |

Odwzorowanie kolumny `odniesienie_id` dla każdego poziomu pamięci:

| `poziom` | Cel `odniesienie_id` |
|---|---|
| globalna | — (puste) |
| srodowisko | `srodowisko.id` |
| modul | `modul.id` |
| projekt | `projekt.id` |
| sesja | `sesja.id` |

*Przykład.* Wskazanie „odrębna pamięć przypisana modułowi Developer” — odpowiada wierszowi z `poziom = 'modul'` i `odniesienie_id` wskazującym moduł Developer.

### 9.2. Konfiguracja pamięci sesji (`konfiguracja_pamieci_sesji`)

Ustawienia pamięci swoiste dla pojedynczej sesji: z poziomu pojedynczej sesji odłącza się dostęp do pamięci przed pierwszym promptem.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| sesja_id | INTEGER | Sesja. | FK → sesja.id, UNIQUE, ON DELETE CASCADE |
| dostep_do_pamieci | INTEGER | Czy sesja korzysta z pamięci (domyślnie `1`). | CHECK (0,1) |
| odlaczono_przed_pierwszym_promptem | INTEGER | Znacznik jednorazowego odłączenia dostępu przed pierwszą wiadomością. | CHECK (0,1) |

---

## 10. Zadania, kolejki, harmonogram i orkiestracja

Diagram relacji (klucze obce) grupy zadań, kolejek i orkiestracji:

```
zadanie ──N:1──> kolejka ──1:N──> pozycja_kolejki ──N:1──> zadanie
   │                │
   │                └──1:N──> log_akcji_kolejki
   ├──N:1──> zlecenie ──1:N──> przebieg_petli        (rozdz. 19)
   ├──N:1──> przebieg_petli
   ├──1:N──> wynik_kontroli_jakosci ──1:N──> ponowienie
   ├──N:1──> harmonogram
   ├──1:N──> komponent_kompozycji ──N:1──> modul
   └──N:1──> sesja | rola_multitasking            (NULL dopuszczalny)
zaleznosc_orkiestracji ──(polimorficzne)──> kanal_modelu | agent | zadanie | kolejka | automatyka | projekt
```

### 10.1. Zadanie (`zadanie`)

Jednostka pracy do wykonania.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator zadania. | PK |
| typ | TEXT | `prompt` \| `akcja` \| `orkiestracja` \| `petla`. | CHECK enum |
| tresc | TEXT | Treść lub definicja zadania. | — |
| sesja_id | INTEGER | Sesja, w której zadanie powstało. | FK → sesja.id, NULL dopuszczalny |
| kolejka_id | INTEGER | Kolejka, do której zadanie należy. | FK → kolejka.id, NULL dopuszczalny |
| harmonogram_id | INTEGER | Powiązany harmonogram. | FK → harmonogram.id, NULL dopuszczalny |
| rola_multitasking_id | INTEGER | Rola MultitaskingAI, jeśli zadanie należy do procesu wielomodelowego. | FK → rola_multitasking.id, NULL dopuszczalny |
| zlecenie_id | INTEGER | Zlecenie, z którego dekompozycji zadanie powstało (rozdział 19.3). | FK → zlecenie.id, ON DELETE CASCADE, NULL dopuszczalny |
| przebieg_petli_id | INTEGER | Przebieg pętli wykonawczej realizujący zadanie (rozdział 19.5). | FK → przebieg_petli.id, ON DELETE SET NULL, NULL dopuszczalny |
| zadanie_nadrzedne_id | INTEGER | Zadanie nadrzędne w dekompozycji zlecenia (samoreferencja). | FK → zadanie.id, NULL dopuszczalny |
| liczba_ponowien | INTEGER | Licznik wykonanych ponowień zadania (rozdział 19.8). | NOT NULL |
| status | TEXT | `oczekujace` \| `w_trakcie` \| `wstrzymane` \| `zakonczone` \| `blad`. | CHECK enum |
| tryb_pracy_w_tle | INTEGER | Konfigurowalny z okna konfiguracji. | CHECK (0,1) |
| data_utworzenia | TEXT | — | NOT NULL |

### 10.2. Komponent kompozycji (`komponent_kompozycji`)

Element, z którego Operator komponuje własne układy pracy.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| modul_id | INTEGER | Moduł zarządzania pracą, w którym komponent jest dostępny (każdy rodzaj dostępny w każdym takim module). | FK → modul.id |
| zadanie_id | INTEGER | Zadanie realizowane przez ten komponent, jeśli dotyczy. | FK → zadanie.id, NULL dopuszczalny |
| rodzaj | TEXT | `prompt` \| `zadanie` \| `akcja` \| `orkiestracja` \| `subagenci` \| `petla_koordynator` \| `petla_wykonawcza` \| `sekwencja` \| `kompilacja`. | CHECK enum |
| definicja | TEXT | Definicja skomponowanego układu. | JSON |

Odwzorowanie dziewięciu rodzajów komponentów kompozycji:

| Rodzaj (`rodzaj`) | Element kompozycji |
|---|---|
| prompt | Prompt pojedynczy |
| zadanie | Zadanie |
| akcja | Akcja |
| orkiestracja | Orkiestracja |
| subagenci | Subagenci |
| petla_koordynator | Pełna pętla z Koordynatorem |
| petla_wykonawcza | Pętla wykonawcza |
| sekwencja | Sekwencja |
| kompilacja | Kompilacja |

### 10.3. Kolejka (`kolejka`)

Struktura porządkująca wykonanie zadań.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator kolejki. | PK |
| nazwa | TEXT | — | NOT NULL |
| zasieg | TEXT | `globalna` \| `lokalna` \| `modelu` \| `agenta` \| `projektu`. | CHECK enum |
| zrodlo_zadan | TEXT | Opis źródła zadań (rola, moduł, automatyka). | — |
| priorytet | INTEGER | Kolejność obsługi. | — |
| obsluga_bledow | TEXT | `retry` \| `route` \| `pause`. | CHECK enum |
| powiazanie_automations | INTEGER | Czy kolejka jest spięta z silnikiem kolejek modułu Automations (ustawienie konfiguracyjne). | CHECK (0,1) |
| zespol_id | INTEGER | Zespół MultitaskingAI, jeśli kolejka jest częścią zapisanego presetu (rozdział 13.1). | FK → zespol_multitasking.id, NULL dopuszczalny |
| zlecenie_id | INTEGER | Zlecenie, któremu kolejka jest dedykowana (rozdział 19.3). | FK → zlecenie.id, ON DELETE SET NULL, NULL dopuszczalny |

### 10.4. Pozycja kolejki (`pozycja_kolejki`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kolejka_id | INTEGER | Kolejka. | FK → kolejka.id, ON DELETE CASCADE |
| zadanie_id | INTEGER | Zadanie w kolejce. | FK → zadanie.id |
| kolejnosc | INTEGER | — | NOT NULL |
| stan | TEXT | `oczekuje` \| `przetwarzana` \| `zakonczona` \| `wstrzymana`. | CHECK enum |

### 10.5. Log akcji kolejki (`log_akcji_kolejki`)

Rejestr akcji silnika kolejek.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kolejka_id | INTEGER | Kolejka, na której wykonano akcję. | FK → kolejka.id, ON DELETE CASCADE |
| pozycja_kolejki_id | INTEGER | Pozycja objęta akcją, jeśli dotyczy. | FK → pozycja_kolejki.id, NULL dopuszczalny |
| rodzaj_akcji | TEXT | `enqueue` \| `dequeue` \| `delay` \| `retry` \| `pause` \| `resume` \| `split` \| `merge` \| `route` \| `branch` \| `condition`. | CHECK enum |
| parametry | TEXT | Parametry akcji. | JSON |
| znacznik_czasu | TEXT | — | NOT NULL |

Jedenaście akcji silnika kolejek (`rodzaj_akcji`):

| Akcja | Znaczenie |
|---|---|
| enqueue | Dodanie zadania do kolejki. |
| dequeue | Pobranie zadania z kolejki do wykonania. |
| delay | Opóźnienie zadania w kolejce. |
| retry | Ponowienie zadania zakończonego błędem. |
| pause | Wstrzymanie kolejki lub pozycji. |
| resume | Wznowienie wstrzymanej kolejki lub pozycji. |
| split | Podział zadania na zadania cząstkowe. |
| merge | Połączenie zadań w jedno. |
| route | Skierowanie zadania do innej kolejki lub roli. |
| branch | Rozgałęzienie przepływu według reguły. |
| condition | Warunkowe rozstrzygnięcie dalszego przepływu. |

Szablon zawartości pola `parametry` (kolumna JSON) dla akcji `route`:

```json
{
  "akcja": "route",
  "z_kolejki_id": 4,
  "do_kolejki_id": 9,
  "kryterium": { "pole": "typ", "wartosc": "orkiestracja" }
}
```

### 10.6. Harmonogram (`harmonogram`)

Definicja wykonania zadań w czasie.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator harmonogramu. | PK |
| regula_czasu | TEXT | Wyrażenie definiujące moment lub cykl uruchomienia. | NOT NULL |
| cyklicznosc | TEXT | `jednorazowo` \| `cyklicznie`. | CHECK enum |
| stan | TEXT | `aktywny` \| `wstrzymany` \| `zakonczony`. | CHECK enum |
| nastepne_uruchomienie | TEXT | — | NULL dopuszczalny |

### 10.7. Zależność orkiestracji (`zaleznosc_orkiestracji`)

Definicja zależności między modelami, agentami, zadaniami, kolejkami, automatyzacjami i projektami.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| nazwa | TEXT | Nazwa opisowa zależności. | NULL dopuszczalny |
| zrodlo_typ | TEXT | `kanal_modelu` \| `agent` \| `zadanie` \| `kolejka` \| `automatyka` \| `projekt`. | CHECK enum |
| zrodlo_id | INTEGER | Odniesienie polimorficzne. | NOT NULL |
| cel_typ | TEXT | Ta sama enumeracja co `zrodlo_typ`. | CHECK enum |
| cel_id | INTEGER | Odniesienie polimorficzne. | NOT NULL |
| typ_zaleznosci | TEXT | `poprzedza` \| `blokuje` \| `wyzwala` \| `warunkuje`. | CHECK enum |
| warunek | TEXT | Definicja warunku, dla `typ_zaleznosci = 'warunkuje'`. | JSON, NULL dopuszczalny |
| zespol_id | INTEGER | Zespół MultitaskingAI, jeśli zależność jest częścią zapisanego presetu. | FK → zespol_multitasking.id, NULL dopuszczalny |

Odwzorowanie przykładu — „zadanie Executora 2 nie rozpocznie się przed zakończeniem etapu przypisanego Executorowi 1”:

| Element zależności | Wartość |
|---|---|
| `zrodlo_typ` / `zrodlo_id` | `zadanie` — etap przypisany Executorowi 1 |
| `cel_typ` / `cel_id` | `zadanie` — zadanie Executora 2 |
| `typ_zaleznosci` | `poprzedza` |

Szablon zawartości pola `warunek` (kolumna JSON) dla `typ_zaleznosci = 'warunkuje'`:

```json
{
  "warunek": {
    "pole_zrodla": "status",
    "operator": "rowne",
    "wartosc": "zakonczone"
  }
}
```

---

## 11. Kanały modeli

Diagram relacji (klucze obce) grupy kanałów modeli:

```
kanal_modelu ──(przypisanie_typ + przypisanie_id, polimorficzne)──> sesja | rola_multitasking
wiadomosc ──N:1──> kanal_modelu       (odpowiedź Wykonawcy, rozdz. 8.4)
agent ──N:1──> kanal_modelu           (model bazowy, rozdz. 12.1)
```

### 11.1. Kanał modelu (`kanal_modelu`)

Sposób połączenia z modelem (architektura, rozdział 9).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator kanału. | PK |
| typ | TEXT | `API` \| `CLI` \| `SSH` \| `HTTP`. | CHECK enum |
| nazwa_modelu | TEXT | Nazwa modelu udostępnianego przez kanał. | NOT NULL |
| dane_dostepowe_odwolanie | TEXT | Odwołanie do danych dostępowych przechowywanych poza bazą (rozdział 1.5). | NOT NULL |
| przypisanie_typ | TEXT | `sesja` \| `rola`. | CHECK enum |
| przypisanie_id | INTEGER | Odniesienie polimorficzne do `sesja.id` lub `rola_multitasking.id`. | NOT NULL |
| aktywny | INTEGER | — | CHECK (0,1) |

Odwzorowanie pary `przypisanie_typ` / `przypisanie_id`:

| `przypisanie_typ` | Cel `przypisanie_id` |
|---|---|
| sesja | `sesja.id` |
| rola | `rola_multitasking.id` |

*Uwaga projektowa.* Wybór modelu i kanału jest dokonywany per sesja lub per rola (architektura, rozdział 9) — stąd para `przypisanie_typ` / `przypisanie_id` zamiast dwóch osobnych, częściowo pustych kolumn kluczy obcych.

---

## 12. Agenci

Diagram relacji (klucze obce) grupy agentów:

```
komponent_wlasny ──1:1──> agent ──N:1──> kanal_modelu (model bazowy)
                            │──1:N──> agent_umiejetnosc
                            │──1:N──> agent_rozszerzenie ──N:1──> rozszerzenie
                            │──1:N──> agent_uprawnienie
                            ├──(rola_multitasking.przypisany_agent_id)──> rola wykonawcza
                            └──(projekt_agent)──> projekt
```

### 12.1. Agent (`agent`)

Szczegóły komponentu własnego rodzaju „agent” — komponent platformowy działający ponad wszystkimi środowiskami.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator agenta. | PK |
| komponent_wlasny_id | INTEGER | Komponent nadrzędny. | FK → komponent_wlasny.id, UNIQUE |
| model_bazowy_kanal_id | INTEGER | Kanał modelu bazowego agenta. | FK → kanal_modelu.id, NULL dopuszczalny |
| tozsamosc | TEXT | Opis tożsamości agenta. | — |
| instrukcje_systemowe | TEXT | — | — |
| konfiguracja_pamieci | TEXT | Konfiguracja pamięci agenta. | JSON |
| konfiguracja_uprawnien | TEXT | Zestaw uprawnień w formie skróconej (pełne uprawnienia — tabela 12.4). | JSON |
| data_utworzenia | TEXT | — | NOT NULL |

Szablon zawartości pól `konfiguracja_pamieci` i `konfiguracja_uprawnien` (kolumny JSON):

```json
{
  "konfiguracja_pamieci": {
    "korzysta_z_pamieci_projektu": true,
    "poziomy": ["projekt", "sesja"]
  },
  "konfiguracja_uprawnien": {
    "pliki": "odczyt",
    "siec": "brak",
    "moduly": ["Developer", "Terminal"]
  }
}
```

*Relacje.* Agent może pełnić rolę wykonawczą w środowisku MultitaskingAI przez `rola_multitasking.przypisany_agent_id` (rozdział 13.2) oraz być przypisany do projektu przez `projekt_agent` (rozdział 6.5).

### 12.2. Umiejętność agenta (`agent_umiejetnosc`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| agent_id | INTEGER | Agent. | FK → agent.id, ON DELETE CASCADE |
| nazwa | TEXT | Nazwa umiejętności. | NOT NULL |
| konfiguracja | TEXT | — | JSON |

### 12.3. Rozszerzenie agenta (`agent_rozszerzenie`)

Łączy agenta z wtyczką, umiejętnością lub konektorem katalogowanym w `rozszerzenie` (rozdział 14).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| agent_id | INTEGER | Agent. | FK → agent.id, ON DELETE CASCADE |
| rozszerzenie_id | INTEGER | Rozszerzenie. | FK → rozszerzenie.id |
| konfiguracja | TEXT | Konfiguracja specyficzna dla tego przypisania. | JSON |

Ograniczenie `UNIQUE(agent_id, rozszerzenie_id)`.

### 12.4. Uprawnienie agenta (`agent_uprawnienie`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| agent_id | INTEGER | Agent. | FK → agent.id, ON DELETE CASCADE |
| zakres | TEXT | Zakres uprawnienia: dostęp do plików, dostęp do sieci albo dostęp do wskazanego modułu. | NOT NULL |
| poziom | TEXT | `brak` \| `odczyt` \| `zapis` \| `pelny`. | CHECK enum |

---

## 13. Środowisko MultitaskingAI

Rozdział opisuje zawartość bocznej nawigacji środowiska MultitaskingAI: panel orkiestracji złożony z sześciu sekcji odpowiadających własnym warstwom sterowania środowiska.

Diagram relacji (klucze obce) grupy środowiska MultitaskingAI:

```
zespol_multitasking ──1:N──> rola_multitasking ──1:N──> subagent
      │                            │──N:1──> agent | kanal_modelu (wykonawca)
      │                            └──N:1──> profil_izolacji
      ├──1:N──> kolejka (zespol_id)
      ├──1:N──> zaleznosc_orkiestracji (zespol_id)
      └──1:N──> przebieg_multitasking ──1:N──> decyzja_procesu ──(samoreferencja)
sekcja_panelu_orkiestracji   (tabela słownikowa: sześć sekcji)
```

### 13.1. Zespół (`zespol_multitasking`)

Zapisana konfiguracja zespołu — preset ról, powiązań i kolejek.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator zespołu. | PK |
| nazwa | TEXT | — | NOT NULL |
| opis | TEXT | — | — |
| data_utworzenia | TEXT | — | NOT NULL |
| data_modyfikacji | TEXT | — | — |

*Relacje.* Zespół grupuje wiersze `rola_multitasking` (13.2), `kolejka` (10.3) i `zaleznosc_orkiestracji` (10.7) mające ustawione `zespol_id` na ten zespół — zapisany preset jest więc zbiorem konkretnych, uprzednio skonfigurowanych ról, kolejek i zależności, wczytywalnym i duplikowalnym bez ponownej konfiguracji.

### 13.2. Rola MultitaskingAI (`rola_multitasking`)

Przypisanie roli w środowisku wielomodelowym.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator roli. | PK |
| zespol_id | INTEGER | Zespół, do którego rola należy jako część zapisanego presetu. | FK → zespol_multitasking.id, NULL dopuszczalny |
| sesja_id | INTEGER | Sesja, w której rola jest aktywną instancją procesu. | FK → sesja.id, NULL dopuszczalny |
| rodzaj_roli | TEXT | `executor_1` \| `executor_2` \| `coordinator` \| `executor_3_validator`. | CHECK enum |
| wcielenie | TEXT | Dla `executor_3_validator`: `validator` \| `reviewer` \| `security_auditor` \| `architect` \| `product_owner` \| `qa_lead` \| `arbitrator` \| inne (tekst dowolny). | NULL dopuszczalny |
| przypisany_agent_id | INTEGER | Agent pełniący rolę, jeśli wskazany. | FK → agent.id, NULL dopuszczalny |
| przypisany_kanal_modelu_id | INTEGER | Model bazowy pełniący rolę, jeśli wskazany bezpośrednio. | FK → kanal_modelu.id, NULL dopuszczalny |
| zakres_odpowiedzialnosci | TEXT | — | — |
| dostep_do_pamieci_projektu | INTEGER | — | CHECK (0,1) |
| subagent_network_aktywny | INTEGER | Czy rola (Executor 1 lub 2) ma aktywną sieć podagentów. | CHECK (0,1) |
| relacja_z_innymi | TEXT | `niezalezna` \| `przekazywanie` \| `naprzemienna` \| `iteracyjna` (relacja między Executor 1 a Executor 2). | NULL dopuszczalny, CHECK enum |
| profil_izolacji_id | INTEGER | Profil izolacji przypisany roli (rozdział 15.2). | FK → profil_izolacji.id, NULL dopuszczalny |

Cztery rodzaje ról (`rodzaj_roli`):

| Rodzaj roli | Odpowiedzialność | Tworzy produkt końcowy | Zarządza procesem | Subagent Network |
|---|---|---|---|---|
| executor_1 | Główny wykonawca — faktyczna realizacja pracy | Tak | Nie | Tak (do 15 podagentów) |
| executor_2 | Równoległy wykonawca — drugi niezależny tor wykonania | Tak | Nie | Tak (do 15 podagentów) |
| coordinator | Planowanie, podział pracy, budowa promptów, sterowanie, zarządzanie kolejką | Nie | Tak | Nie |
| executor_3_validator | Funkcja kontrolna lub doradcza (konkretne wcielenie w kolumnie `wcielenie`) | Zależnie od wcielenia | Nie | — |

*Relacje.* Rola jest wykonawcą przypisanym albo przez `przypisany_agent_id` (agent z modułu Agents), albo przez `przypisany_kanal_modelu_id` (model bazowy bez pośrednictwa agenta) — pola wzajemnie wykluczające się, każde dopuszczające wartość NULL, egzekwowane na poziomie aplikacji. Dokładnie jeden wiersz `rodzaj_roli = 'coordinator'` na aktywny proces jest oczekiwany logiką aplikacji, nie ograniczeniem schematu.

### 13.3. Subagent (`subagent`)

Podagent uruchomiony przez wykonawcę w ramach Subagent Network — do piętnastu jednocześnie na jednego wykonawcę.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| rola_multitasking_id | INTEGER | Wykonawca nadrzędny (Executor 1 lub Executor 2). | FK → rola_multitasking.id, ON DELETE CASCADE |
| nazwa | TEXT | Nazwa podagenta, w brzmieniu „Agent Backend”, „Agent Security”. | NOT NULL |
| agent_id | INTEGER | Agent bazowy, jeśli podagent jest oparty na zdefiniowanym agencie. | FK → agent.id, NULL dopuszczalny |
| specjalizacja | TEXT | Wąski, dobrze zdefiniowany zakres podagenta. | — |
| status | TEXT | `aktywny` \| `zakonczony` \| `blad`. | CHECK enum |
| kolejnosc | INTEGER | — | — |

Limit piętnastu podagentów na wykonawcę jest egzekwowany na poziomie aplikacji — wyzwalaczem sprawdzającym liczbę aktywnych wierszy przy wstawieniu; SQLite nie wyraża tego ograniczenia deklaratywnie w definicji tabeli.

### 13.4. Sekcja panelu orkiestracji (`sekcja_panelu_orkiestracji`)

Boczna nawigacja środowiska MultitaskingAI — sześć sekcji, z konfigurowalną kolejnością i widocznością; zestaw poniższy jest zestawem domyślnym.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kod | TEXT | `zespoly` \| `role` \| `kolejki` \| `orkiestracja` \| `harmonogram_automatyki` \| `monitor_procesu`. | UNIQUE, CHECK enum |
| nazwa | TEXT | — | NOT NULL |
| kolejnosc | INTEGER | Kolejność w panelu. | NOT NULL |
| widoczna | INTEGER | Widoczność sekcji (domyślnie wszystkie widoczne). | CHECK (0,1) |

Sześć domyślnych sekcji panelu orkiestracji:

| Kolejność | `kod` | Nazwa sekcji |
|---|---|---|
| 1 | zespoly | Zespoły |
| 2 | role | Role |
| 3 | kolejki | Kolejki |
| 4 | orkiestracja | Orkiestracja |
| 5 | harmonogram_automatyki | Harmonogram i automatyki |
| 6 | monitor_procesu | Monitor procesu |

### 13.5. Przebieg procesu (`przebieg_multitasking`)

Pojedynczy przebieg pętli autonomicznej — treść sekcji „Monitor procesu” panelu orkiestracji.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator przebiegu. | PK |
| zespol_id | INTEGER | Zespół realizujący przebieg. | FK → zespol_multitasking.id, NULL dopuszczalny |
| sesja_id | INTEGER | Sesja, w której przebieg jest obserwowany. | FK → sesja.id, NULL dopuszczalny |
| harmonogram_id | INTEGER | Harmonogram, który wyzwolił przebieg (praca cykliczna sterowana modułem Automations). | FK → harmonogram.id, NULL dopuszczalny |
| numer_przebiegu | INTEGER | Kolejny numer przebiegu w ramach zespołu. | — |
| status | TEXT | `w_trakcie` \| `zakonczony` \| `blad` \| `wstrzymany`. | CHECK enum |
| data_rozpoczecia | TEXT | — | NOT NULL |
| data_zakonczenia | TEXT | — | NULL dopuszczalny |
| podsumowanie | TEXT | — | NULL dopuszczalny |

### 13.6. Decyzja procesu (`decyzja_procesu`)

Hierarchia decyzji podejmowanych w toku przebiegu, obserwowana przez Always On Display jako operatora lub obserwatora procesu.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| przebieg_id | INTEGER | Przebieg, do którego decyzja należy. | FK → przebieg_multitasking.id, ON DELETE CASCADE |
| rola_multitasking_id | INTEGER | Rola, która podjęła decyzję. | FK → rola_multitasking.id, NULL dopuszczalny |
| nadrzedna_decyzja_id | INTEGER | Decyzja nadrzędna w hierarchii, jeśli dotyczy. | FK → decyzja_procesu.id (samoreferencja), NULL dopuszczalny |
| typ_decyzji | TEXT | Rodzaj decyzji: planowanie, sterowanie albo walidacja. | — |
| tresc | TEXT | — | — |
| znacznik_czasu | TEXT | — | NOT NULL |

---

## 14. Rozszerzenia

Diagram relacji (klucze obce) grupy rozszerzeń:

```
rozszerzenie ──1:N──> agent_rozszerzenie ──N:1──> agent
```

### 14.1. Rozszerzenie (`rozszerzenie`)

Wtyczka, umiejętność, konektor lub serwer MCP (architektura, rozdział 14).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator rozszerzenia. | PK |
| nazwa | TEXT | — | NOT NULL |
| rodzaj | TEXT | `wtyczka` \| `umiejetnosc` \| `konektor` \| `serwer_mcp`. | CHECK enum |
| zrodlo | TEXT | `danaco_plugin` \| `personal`. | CHECK enum |
| stan_wlaczenia | INTEGER | — | CHECK (0,1) |
| konfiguracja | TEXT | — | JSON |
| lokalizacja | TEXT | Katalog na serwerze; dla źródła `personal` — katalog użytkownika, dla `danaco_plugin` — katalog pakietu aplikacyjnego. | NOT NULL |
| data_instalacji | TEXT | — | NOT NULL |

Rodzaje rozszerzeń i ich źródła (kolumny `rodzaj` i `zrodlo`, architektura rozdział 14):

| Kolumna | Wartość | Znaczenie |
|---|---|---|
| `rodzaj` | wtyczka | Wtyczka rozszerzająca funkcje platformy. |
| `rodzaj` | umiejetnosc | Umiejętność (skill) możliwa do przypisania agentowi. |
| `rodzaj` | konektor | Konektor do usługi zewnętrznej. |
| `rodzaj` | serwer_mcp | Serwer Model Context Protocol. |
| `zrodlo` | danaco_plugin | Rozszerzenie z katalogu pakietu aplikacyjnego. |
| `zrodlo` | personal | Rozszerzenie z katalogu użytkownika. |

*Uwaga projektowa.* Rdzeń nie rozróżnia źródła rozszerzenia przy wykonaniu (architektura, rozdział 14) — kolumna `zrodlo` służy wyłącznie prezentacji i zarządzaniu, nie zmienia sposobu integracji. Rozszerzenia są przypisywane agentom przez `agent_rozszerzenie` (rozdział 12.3).

---

## 15. Profile izolacji i punkty izolacji

Rozdział opisuje szczegółowy interfejs okna konfiguracji punktów izolacji: trzy panele — selektor zasięgu, macierz izolacji, profil i podgląd — odwzorowane jako słownik poziomów zasięgu, encja profilu z dwiema grupami reguł oraz encja przypisania.

Diagram relacji (klucze obce) grupy izolacji:

```
poziom_zasiegu ──1:N──> profil_izolacji ──1:N──> regula_izolacji_kontekstu
                                        │──1:N──> regula_izolacji_technicznej
                                        └──1:N──> przypisanie_profilu_izolacji ──(polimorficzne)──>
                                              srodowisko | modul | para_modulow | projekt |
                                              karta_sesji | rola_multitasking
para_modulow ──N:1──> modul (a) + modul (b)
```

### 15.1. Poziom zasięgu (`poziom_zasiegu`)

Tabela słownikowa, osiem stałych wierszy uporządkowanych według pierwszeństwa.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kod | TEXT | `globalny` \| `srodowisko` \| `modul` \| `para_modulow` \| `projekt` \| `karta_sesji` \| `rola` \| `okno`. | UNIQUE, CHECK enum |
| nazwa | TEXT | — | NOT NULL |
| pierwszenstwo | INTEGER | `1` (globalny, warstwa bazowa) do `8` (okno komunikacji, najwyższe pierwszeństwo). | NOT NULL |

Osiem stałych wierszy słownika, od najogólniejszego do najbardziej szczegółowego:

| Pierwszeństwo | `kod` | Nazwa | Przykład zastosowania |
|---|---|---|---|
| 1 (warstwa bazowa) | globalny | Globalny (domyślny) | Domyślna polityka izolacji obowiązująca w całej platformie. |
| 2 | srodowisko | Środowisko | Odmienna polityka dla środowiska CodeStudio niż dla TalkIn. |
| 3 | modul | Moduł | Odrębna pamięć przypisana modułowi Developer. |
| 4 | para_modulow | Para modułów | Współdzielenie operacji kontekstowych AI między Studio a Translate. |
| 5 | projekt | Projekt | Rozdzielenie kontekstu między dwoma projektami w module Workspace. |
| 6 | karta_sesji | Karta sesji | Jednorazowe współdzielenie historii między dwiema otwartymi kartami. |
| 7 | rola | Rola (MultitaskingAI) | Profil izolacji przypisany roli Executor 1. |
| 8 (najwyższe) | okno | Okno komunikacji | Odrębna pamięć i odrębne uprawnienia okna Execution Loop Window wobec okna Chat Window tej samej karty sesji. |

Wartość `okno` jest postacią docelową słownika przyjętą rozstrzygnięciem Właściciela: kod aplikacji jej jeszcze nie zna — wyliczenie `ConfigScope` w `budowa/shared/contract.json` zawiera wprawdzie wartość `window`, lecz warstwa danych stosuje dziś siedem wierszy słownika. Dopisanie ósmego wiersza wraz z ograniczeniem `CHECK` i zakresem pola `pierwszenstwo` jest osobnym zadaniem w kodzie. Żadna z dotychczasowych wartości nie zmienia znaczenia ani pierwszeństwa.

Ten sam słownik obsługuje zarówno `profil_izolacji` (15.2), jak i `ustawienie` (rozdział 17.1), zapewniając jeden, spójny mechanizm dziedziczenia dla obu rodzajów konfiguracji.

### 15.2. Profil izolacji (`profil_izolacji`)

Nazwany, zapisywalny zestaw ustawień izolacji kontekstu i izolacji technicznej dla wskazanego poziomu zasięgu.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator profilu. | PK |
| nazwa | TEXT | — | NOT NULL |
| poziom_zasiegu_id | INTEGER | Poziom, na którym profil domyślnie obowiązuje. | FK → poziom_zasiegu.id |
| warstwa | TEXT | `domyslna` \| `sesji` (architektura, rozdział 13). | CHECK enum |
| data_utworzenia | TEXT | — | NOT NULL |

Szablon profilu izolacji rozłożonego na wiersze tabel podrzędnych (przykład — dwa projekty rozdzielone, ze współdzieloną biblioteką):

```
profil_izolacji: { nazwa: "Projekty rozdzielone", poziom_zasiegu: projekt, warstwa: domyslna }
  regula_izolacji_kontekstu:
      historia = odrebna
      pamiec   = odrebna
      kontekst = odrebna
  regula_izolacji_technicznej:
      (brak wierszy → żaden zakres techniczny nie jest włączony)
  przypisanie_profilu_izolacji:
      typ_celu = projekt, cel_id = <projekt pierwszy>
      typ_celu = projekt, cel_id = <projekt drugi>
```

### 15.3. Reguła izolacji kontekstu (`regula_izolacji_kontekstu`)

Trzy pozycje macierzy izolacji, grupa „izolacja kontekstu”.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| profil_izolacji_id | INTEGER | Profil, do którego reguła należy. | FK → profil_izolacji.id, ON DELETE CASCADE |
| wymiar | TEXT | `historia` \| `pamiec` \| `kontekst`. | CHECK enum |
| wartosc | TEXT | `wspoldzielona` \| `odrebna`. | CHECK enum |

Ograniczenie `UNIQUE(profil_izolacji_id, wymiar)`. Trzy wymiary izolacji kontekstu:

| Wymiar (`wymiar`) | Znaczenie | Wartości (`wartosc`) |
|---|---|---|
| historia | Historia rozmów sesji | wspoldzielona \| odrebna |
| pamiec | Pamięć kontekstowa | wspoldzielona \| odrebna |
| kontekst | Kontekst przekazywany do AI | wspoldzielona \| odrebna |

### 15.4. Reguła izolacji technicznej (`regula_izolacji_technicznej`)

Osiem zakresów zdefiniowanych w architekturze (rozdział 12), grupa „izolacja techniczna” macierzy izolacji.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| profil_izolacji_id | INTEGER | Profil, do którego reguła należy. | FK → profil_izolacji.id, ON DELETE CASCADE |
| zakres | TEXT | `katalog_roboczy_sesji` \| `srodowisko_procesu` \| `katalog_danych_modelu` \| `dostep_sieciowy` \| `odczyt_zapis_plikow` \| `konto_i_token` \| `model_procesu` \| `serwer_wykonania`. | CHECK enum |
| stan | TEXT | `wlaczony` \| `wylaczony`. | CHECK enum |

Ograniczenie `UNIQUE(profil_izolacji_id, zakres)`. Osiem zakresów izolacji technicznej (architektura, rozdział 12):

| Zakres (`zakres`) | Znaczenie |
|---|---|
| katalog_roboczy_sesji | Odrębny katalog roboczy procesu sesji. |
| srodowisko_procesu | Odrębne środowisko (zmienne) procesu. |
| katalog_danych_modelu | Odrębny katalog danych i konfiguracji modelu. |
| dostep_sieciowy | Ograniczenie dostępu sieciowego procesu. |
| odczyt_zapis_plikow | Ograniczenie zakresu odczytu i zapisu plików. |
| konto_i_token | Odrębne konto i token per sesja. |
| model_procesu | Model procesu odrębny albo współdzielony. |
| serwer_wykonania | Serwer wykonania procesu. |

Stan wyjściowy platformy — brak wierszy tej tabeli dla profilu domyślnego, co odpowiada „żaden z zakresów nie jest domyślnie włączony”.

### 15.5. Przypisanie profilu izolacji (`przypisanie_profilu_izolacji`)

Realizuje panel selektora zasięgu: przypisanie profilu do konkretnego elementu na jednym z ośmiu poziomów.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| profil_izolacji_id | INTEGER | Profil przypisywany. | FK → profil_izolacji.id |
| typ_celu | TEXT | Ta sama enumeracja co `poziom_zasiegu.kod`. | CHECK enum |
| cel_id | INTEGER | Odniesienie polimorficzne do elementu wskazanego poziomu (`srodowisko.id`, `modul.id`, `para_modulow.id`, `projekt.id`, `karta_sesji.id`, `rola_multitasking.id`); puste dla `typ_celu = 'globalny'`. | NULL dopuszczalny |
| data_przypisania | TEXT | — | NOT NULL |

**Poziom `okno` bez encji celu — pozycja do rozstrzygnięcia.** Ósmy poziom zasięgu
(`poziom_zasiegu.kod = 'okno'`, rozdz. 15.1) jest dopuszczalną wartością kolumny `typ_celu`,
lecz wyliczenie odniesień polimorficznych kolumny `cel_id` nie wskazuje dla niego żadnej encji
modelu — okno komunikacji nie ma w niniejszym opracowaniu tabeli, do której `cel_id` mógłby
się odnieść, tak jak odnosi się `srodowisko.id` dla poziomu środowiska albo `karta_sesji.id`
dla poziomu karty sesji. Model pozostaje w tym miejscu bez zmian: rozstrzygnięcie, czy okno
komunikacji otrzymuje własną encję trwałą, czy `cel_id` odnosi się dla tego poziomu do
identyfikatora okna prowadzonego wyłącznie w stanie operacyjnym rdzenia, należy do Właściciela
i nie jest przesądzane tutaj. **[DO DECYZJI OPERATORA]** — czy poziom `okno` otrzymuje encję
trwałą w schemacie przechowywania, a jeśli tak, jaka jest jej nazwa i zakres pól?

Reguła rozstrzygania (polityka efektywna) — algorytm odczytu reguły dla danego elementu, wykonywany osobno dla każdego wymiaru izolacji (historia, pamięć, kontekst oraz każdy z ośmiu zakresów technicznych):

1. Zbierz wszystkie wiersze `przypisanie_profilu_izolacji` z poziomów obejmujących dany element.
2. Wybierz wiersz o najwyższej wartości `poziom_zasiegu.pierwszenstwo`.
3. Brak przypisania na danym poziomie oznacza przejście do poziomu bezpośrednio szerszego.
4. Powtarzaj aż do poziomu globalnego, który stanowi warstwę bazową.

Ten sam algorytm, zaimplementowany raz na poziomie aplikacji, obsługuje panel podglądu polityki efektywnej oraz odczyt `ustawienie` w rozdziale 17.

### 15.6. Para modułów (`para_modulow`)

Element docelowy poziomu zasięgu „para modułów” — relacja Studio–Translate.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| modul_a_id | INTEGER | Pierwszy moduł pary. | FK → modul.id |
| modul_b_id | INTEGER | Drugi moduł pary. | FK → modul.id |

Ograniczenie `UNIQUE(modul_a_id, modul_b_id)`.

---

## 16. Artefakty i biblioteka

Diagram relacji (klucze obce) grupy artefaktów:

```
artefakt ──1:N──> wersja_artefaktu
   │──1:N──> kolekcja_artefakt ──N:1──> kolekcja
   │──1:N──> etykieta_artefaktu
   └──N:1──> sesja | projekt          (co najmniej jedno ustawione)
```

### 16.1. Artefakt (`artefakt`)

Plik utworzony w trakcie pracy, wykrywany automatycznie.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator artefaktu. | PK |
| sesja_id | INTEGER | Sesja pochodzenia. | FK → sesja.id, NULL dopuszczalny |
| projekt_id | INTEGER | Projekt pochodzenia (Project Library) lub przechowywania (Library). | FK → projekt.id, NULL dopuszczalny |
| nazwa_pliku | TEXT | — | NOT NULL |
| sciezka_pliku | TEXT | Treść przechowywana jako plik, z odwołaniem w bazie (architektura, rozdział 10). | NOT NULL |
| typ | TEXT | Typ / kategoria pliku. | — |
| wersja_biezaca | INTEGER | Numer aktualnej wersji. | NOT NULL |
| data_utworzenia | TEXT | — | NOT NULL |

Co najmniej jedno z pól `sesja_id`, `projekt_id` powinno być ustawione — reguła egzekwowana na poziomie aplikacji.

### 16.2. Wersja artefaktu (`wersja_artefaktu`)

Historia wersji, wykorzystywana przez Versioning Panel oraz Diff/Grep Panel.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| artefakt_id | INTEGER | Artefakt macierzysty. | FK → artefakt.id, ON DELETE CASCADE |
| numer_wersji | INTEGER | — | NOT NULL |
| sciezka_pliku | TEXT | — | NOT NULL |
| autor | TEXT | `uzytkownik` \| `wykonawca`. | CHECK enum |
| data_utworzenia | TEXT | — | NOT NULL |

Ograniczenie `UNIQUE(artefakt_id, numer_wersji)`.

### 16.3. Kolekcja (`kolekcja`)

Grupowanie w Tags & Collections modułu Library, lub w bibliotece projektu.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| nazwa | TEXT | — | NOT NULL |
| projekt_id | INTEGER | Zakres Project Library, jeśli kolekcja jest ograniczona do projektu. | FK → projekt.id, NULL dopuszczalny |

### 16.4. Przypisanie kolekcji (`kolekcja_artefakt`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kolekcja_id | INTEGER | — | FK → kolekcja.id, ON DELETE CASCADE |
| artefakt_id | INTEGER | — | FK → artefakt.id, ON DELETE CASCADE |

Ograniczenie `UNIQUE(kolekcja_id, artefakt_id)`.

### 16.5. Etykieta artefaktu (`etykieta_artefaktu`)

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| artefakt_id | INTEGER | — | FK → artefakt.id, ON DELETE CASCADE |
| etykieta | TEXT | — | NOT NULL |

---

## 17. Konfiguracja warstwowa

Diagram relacji (klucze obce) grupy konfiguracji warstwowej:

```
poziom_zasiegu ──1:N──> ustawienie ──(polimorficzne: poziom_odniesienie_id)──>
      srodowisko | modul | para_modulow | projekt | karta_sesji | rola
      (poziom globalny → poziom_odniesienie_id puste)
```

### 17.1. Ustawienie (`ustawienie`)

Pojedynczy parametr zachowania, konfigurowalny z okna konfiguracji (architektura, rozdział 13).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator ustawienia. | PK |
| klucz | TEXT | Klucz parametru — `retencja_historii`, `tryb_pracy_w_tle`. | NOT NULL |
| wartosc | TEXT | Wartość parametru; dopuszczalny zapis JSON dla wartości złożonych. | — |
| poziom_zasiegu_id | INTEGER | Poziom obowiązywania (współdzielony słownik z rozdziałem 15.1). | FK → poziom_zasiegu.id |
| poziom_odniesienie_id | INTEGER | Odniesienie polimorficzne do elementu wskazanego poziomu; puste dla `globalny`. | NULL dopuszczalny |
| objasnienie | TEXT | Treść objaśnienia kontekstowego (`[?]`) opisująca działanie ustawienia i jego wpływ na aplikację (architektura, rozdział 13). | NOT NULL |
| data_modyfikacji | TEXT | — | — |

Ograniczenie `UNIQUE(klucz, poziom_zasiegu_id, poziom_odniesienie_id)`.

Przykład rozstrzygnięcia ustawienia `retencja_historii` na kolejnych poziomach zasięgu — wartość globalna, następnie środowisko, projekt i sesja:

| Poziom zasięgu | Wartość `retencja_historii` (przykład) |
|---|---|
| globalny | bez limitu |
| srodowisko (CodeStudio) | 90 dni |
| projekt | 30 dni |
| karta_sesji | 7 dni (nadpisanie jednorazowe) |

*Uwaga projektowa.* Warstwowość ogólnej konfiguracji (architektura, rozdział 13: warstwa domyślna i warstwa sesji) jest przypadkiem szczególnym reguły ośmiu poziomów zasięgu z rozdziału 15.5; poziomy pośrednie (`srodowisko`, `modul`, `para_modulow`, `projekt`) są dostępne tam, gdzie ustawienie ma sens na danym poziomie. Ta sama reguła rozstrzygania z rozdziału 15.5 stosuje się bez zmian.

### 17.2. Objaśnienie kontekstowe jako część definicji

Zgodnie z architekturą (rozdział 13), treść objaśnienia (`objasnienie`) jest częścią definicji ustawienia, nie treścią pomocniczą przechowywaną osobno — stąd pole `NOT NULL` w tabeli `ustawienie`: każdy wiersz musi nieść własne objaśnienie w chwili zapisu.

---

## 18. Funkcje globalne — dane wspierające

Mobile i Always On Display nie tworzą własnej przestrzeni roboczej ani nowych sesji, dlatego nie wymagają rozbudowanego zestawu własnych encji. Poniższe tabele obejmują wyłącznie dane, których nie da się odwzorować istniejącymi encjami rozdziałów 8–13, oraz encję powiadomienia — wspólną dla centrum powiadomień każdego środowiska i dla powiadomień wypychanych funkcji Mobile.

Diagram relacji (klucze obce) grupy funkcji globalnych:

```
sugestia_aod ──(kontekst_typ + kontekst_id, polimorficzne)──> srodowisko | modul | projekt | sesja
powiadomienie ──(zrodlo_typ + zrodlo_id, polimorficzne)──> srodowisko | modul | projekt | sesja | zadanie | przebieg_petli | automatyka
powiadomienie ──N:1──> srodowisko | sesja | urzadzenie          (NULL dopuszczalny)
polecenie_glosowe ──N:1──> profil_asystenta | sesja   (NULL dopuszczalny)
```

### 18.1. Mobile

Funkcja Mobile korzysta z istniejących encji rozdziałów 8–13 oraz z trzech encji tożsamości i dostępu opisanych w rozdziale 3: `poswiadczenie_urzadzenia`, `kod_parowania` i `dzialanie_lokalne`. Pełne opracowanie funkcji zawiera dokument funkcje-globalne/mobile.md (rozdział 11.3).

| Aspekt funkcji Mobile | Odwzorowanie w encjach |
|---|---|
| Rozróżnienie telefon / tablet | `urzadzenie.typ` |
| Aktywność funkcji Mobile | `urzadzenie.tryb_mobile_aktywny` |
| Rejestr sparowanych urządzeń | `urzadzenie` w powiązaniu z `poswiadczenie_urzadzenia` (rozdz. 3.2, 3.4) |
| Protokół parowania — kod jednorazowy | `kod_parowania` (rozdz. 3.5) |
| Poświadczenie urządzenia, wygaśnięcie i odwołanie dostępu | `poswiadczenie_urzadzenia` (rozdz. 3.4) |
| Kolejka działań wykonanych bez połączenia | `dzialanie_lokalne` (rozdz. 3.6) |
| Monitoring i zarządzanie procesami | Odczyt i zmiana stanu `sesja`, `zadanie`, `kolejka`, `przebieg_multitasking`, `automatyka` z urządzenia mobilnego |
| Powiadomienia wypychane na urządzenie | `powiadomienie` z `kanal_dostarczenia = 'centrum_i_push'` i ustawionym `urzadzenie_id` (rozdz. 18.4) |

### 18.2. Sugestia Always On Display (`sugestia_aod`)

Proaktywne doradztwo — sugestie działań, sugestie konfiguracji, wskazania problemów, rekomendacje kolejnych kroków. Reguły wyzwalania sugestii, katalog ich rodzajów i cykl życia opisuje [Always-on Display](../funkcje-globalne/always-on-display.md) (rozdz. 3–4, 9.4).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| kontekst_typ | TEXT | `srodowisko` \| `modul` \| `projekt` \| `sesja`. | CHECK enum |
| kontekst_id | INTEGER | Odniesienie polimorficzne do elementu, którego sugestia dotyczy. | NULL dopuszczalny |
| rodzaj | TEXT | `doradztwo` \| `konfiguracja` \| `problem` \| `kolejny_krok`. | CHECK enum |
| tresc | TEXT | — | NOT NULL |
| status | TEXT | `nowa` \| `przyjeta` \| `odrzucona`. | CHECK enum |
| znacznik_czasu | TEXT | — | NOT NULL |

Odwzorowanie pary `kontekst_typ` / `kontekst_id`:

| `kontekst_typ` | Cel `kontekst_id` |
|---|---|
| srodowisko | `srodowisko.id` |
| modul | `modul.id` |
| projekt | `projekt.id` |
| sesja | `sesja.id` |

### 18.3. Polecenie głosowe (`polecenie_glosowe`)

Log Voice Console / Actions Monitor / Activity Feed modułu Assistant, stosowany również do komunikacji głosowej Always On Display. Rozgraniczenie obu torów głosowych — sesyjnego toru modułu Assistant i globalnego toru Always On Display, rozpoznawalnych po wartości pola `profil_asystenta_id` — podaje [Always-on Display](../funkcje-globalne/always-on-display.md) (rozdz. 5).

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| profil_asystenta_id | INTEGER | Profil asystenta obsługujący polecenie; puste dla poleceń wydanych bezpośrednio przez Always On Display. | FK → profil_asystenta.id, NULL dopuszczalny |
| sesja_id | INTEGER | Sesja, w kontekście której polecenie zostało wydane. | FK → sesja.id, NULL dopuszczalny |
| tresc_polecenia | TEXT | — | NOT NULL |
| status | TEXT | `w_trakcie` \| `zakonczone` \| `blad`. | CHECK enum |
| wynik | TEXT | — | NULL dopuszczalny |
| znacznik_czasu | TEXT | — | NOT NULL |

### 18.4. Powiadomienie (`powiadomienie`)

Pojedyncze zdarzenie zarejestrowane przez centrum powiadomień — jeden mechanizm powiadamiania obowiązujący w całej platformie, opisany jako karta komponentu w [Katalog komponentów](../interfejs-uzytkownika/katalog-komponentow.md) (rozdz. 11.6). Encja obsługuje zarówno prezentację zdarzenia w centrum powiadomień, jak i jego dostarczenie na urządzenie mobilne przez funkcję Mobile.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator powiadomienia. | PK |
| klasa | TEXT | Klasa zdarzenia: `zakonczenie` \| `decyzja` \| `blad` \| `wzmianka` \| `termin` \| `automatyka` \| `system`. | CHECK enum |
| waga | TEXT | `informacyjna` \| `normalna` \| `wymagajaca_decyzji`. | CHECK enum |
| zrodlo_typ | TEXT | `srodowisko` \| `modul` \| `projekt` \| `sesja` \| `zadanie` \| `przebieg_petli` \| `automatyka`. | CHECK enum |
| zrodlo_id | INTEGER | Odniesienie polimorficzne do elementu, którego zdarzenie dotyczy. | NULL dopuszczalny |
| srodowisko_id | INTEGER | Środowisko, w którego centrum powiadomień zdarzenie jest prezentowane. | FK → srodowisko.id, NULL dopuszczalny |
| sesja_id | INTEGER | Sesja pochodzenia zdarzenia. | FK → sesja.id, NULL dopuszczalny |
| tresc | TEXT | Treść zdarzenia prezentowana w centrum powiadomień. | NOT NULL |
| akcja | TEXT | Definicja działań dostępnych z poziomu powiadomienia (przejście do źródła, zatwierdzenie, wstrzymanie, ponowienie, odłożenie). | JSON |
| stan | TEXT | `nowe` \| `odczytane` \| `obsluzone` \| `odlozone`. | CHECK enum |
| kanal_dostarczenia | TEXT | `centrum` \| `centrum_i_push`. | CHECK enum |
| urzadzenie_id | INTEGER | Urządzenie, na które zdarzenie zostało wypchnięte przez funkcję Mobile. | FK → urzadzenie.id, NULL dopuszczalny |
| znacznik_czasu | TEXT | — | NOT NULL |

Odwzorowanie pary `zrodlo_typ` / `zrodlo_id`:

| `zrodlo_typ` | Cel `zrodlo_id` |
|---|---|
| srodowisko | `srodowisko.id` |
| modul | `modul.id` |
| projekt | `projekt.id` |
| sesja | `sesja.id` |
| zadanie | `zadanie.id` |
| przebieg_petli | `przebieg_petli.id` |
| automatyka | `automatyka.id` |

*Uwaga projektowa.* Zakres klas zdarzeń zgłaszanych przez centrum powiadomień oraz kanał dostarczenia są ustawieniami konfiguracyjnymi warstwy globalnej i warstwy środowiska (rozdz. 17.1), zgodnie z sekcją powiadomień okna Ustawień ([Ustawienia](../interfejs-uzytkownika/ustawienia.md), rozdz. 7). Powiadomienie o wadze `wymagajaca_decyzji` pozostaje w stanie `nowe` do chwili podjęcia decyzji w Chat Window albo w Execution Loop Window. Sugestie Always On Display (18.2) są odrębną encją — doradztwo proaktywne nie jest zdarzeniem powiadomienia.

---

## 19. Warstwa komunikacji operacyjnej

Rozdział opisuje encje obsługujące dwa kanały komunikacji operacyjnej platformy: Chat Window (kanał Użytkownik ↔ Wykonawca) oraz Execution Loop Window (kanał Koordynator ↔ Wykonawca). Oba kanały są pierwszoplanowymi elementami architektury danych — każda praca prowadzona w dowolnym module powstaje jako rozmowa w kanale pierwszym i jest realizowana jako zlecenie prowadzone w kanale drugim.

Diagram relacji (klucze obce) warstwy komunikacji operacyjnej:

```
sesja ──1:N──> rozmowa ──1:N──> wiadomosc ──1:N──> zalacznik_wiadomosci
                  │                                  (Użytkownik ↔ Wykonawca)
                  └──1:N──> zlecenie
                              │
                              ├──1:N──> zadanie ──N:1──> kolejka ──1:N──> pozycja_kolejki
                              │            │
                              │            ├──1:N──> wynik_kontroli_jakosci ──1:N──> ponowienie
                              │            └──1:N──> komunikat_sterujacy
                              └──1:N──> przebieg_petli
                                           │        (Koordynator ↔ Wykonawca)
                                           ├──1:N──> komunikat_sterujacy
                                           ├──1:N──> wynik_kontroli_jakosci
                                           └──N:1──> przebieg_multitasking (rozdz. 13.5)
przebieg_petli ──N:1──> rola_multitasking (Koordynator)
komunikat_sterujacy ──N:1──> rola_multitasking (nadawca wykonawczy)
```

### 19.1. Rozmowa (`rozmowa`)

Wątek głównego okna komunikacji (Chat Window) prowadzony między Użytkownikiem a Wykonawcą w obrębie sesji. Sesja gromadzi wiele rozmów; rozmowa jest jednostką retencji, wyszukiwania i przenoszenia kontekstu między modułami.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator rozmowy. | PK |
| sesja_id | INTEGER | Sesja, w której rozmowa jest prowadzona. | FK → sesja.id, ON DELETE CASCADE |
| projekt_id | INTEGER | Projekt, w kontekście którego rozmowa działa. | FK → projekt.id, NULL dopuszczalny |
| tytul | TEXT | Etykieta rozmowy, nadawana lub generowana. | NULL dopuszczalny |
| stan | TEXT | `aktywna` \| `wstrzymana` \| `zarchiwizowana`. | CHECK enum |
| kanal_modelu_id | INTEGER | Kanał modelu obsługujący Wykonawcę w tej rozmowie. | FK → kanal_modelu.id, NULL dopuszczalny |
| kontekst | TEXT | Kontekst rozmowy przekazywany Wykonawcy. | JSON |
| data_utworzenia | TEXT | — | NOT NULL |
| data_ostatniej_wiadomosci | TEXT | — | NULL dopuszczalny |

*Relacje.* Rozmowa gromadzi `wiadomosc` (19.2) i jest źródłem `zlecenie` (19.3). Zakres współdzielenia historii rozmów między sesjami i projektami rozstrzyga `profil_izolacji` (rozdział 15) w wymiarze `historia`.

### 19.2. Wiadomość rozmowy (`wiadomosc`)

Pojedynczy wpis kanału Użytkownik ↔ Wykonawca. Struktura pól encji `wiadomosc` jest opisana w rozdziale 8.4; kolumna `rozmowa_id` przypisuje wpis do wątku rozmowy, a kolumna `rola` rozróżnia nadawcę (`uzytkownik` \| `wykonawca`).

| Aspekt kanału pierwszego | Odwzorowanie w modelu |
|---|---|
| Polecenie w języku naturalnym | `wiadomosc` z `rola = 'uzytkownik'` |
| Strumień odpowiedzi i wyników Wykonawcy | `wiadomosc` z `rola = 'wykonawca'` i wypełnionym `kanal_modelu_id` |
| Zatwierdzenie lub przerwanie działania | `komunikat_sterujacy` (19.6) o rodzaju `zatwierdzenie`, `przerwanie` |
| Materiał przekazany do rozmowy | `zalacznik_wiadomosci` (rozdział 8.5) |
| Wynik utrwalony jako plik | `artefakt` (rozdział 16.1) |

### 19.3. Zlecenie (`zlecenie`)

Jednostka pracy przyjęta od Użytkownika i przekazana Koordynatorowi do dekompozycji. Zlecenie jest encją nadrzędną wobec zadań i przebiegów pętli wykonawczej.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator zlecenia. | PK |
| rozmowa_id | INTEGER | Rozmowa, z której zlecenie powstało. | FK → rozmowa.id, ON DELETE SET NULL, NULL dopuszczalny |
| wiadomosc_id | INTEGER | Wiadomość będąca treścią pierwotną zlecenia. | FK → wiadomosc.id, NULL dopuszczalny |
| sesja_id | INTEGER | Sesja realizująca zlecenie. | FK → sesja.id, NULL dopuszczalny |
| projekt_id | INTEGER | Projekt, w kontekście którego zlecenie jest realizowane. | FK → projekt.id, NULL dopuszczalny |
| koordynator_rola_id | INTEGER | Rola pełniąca funkcję Koordynatora zlecenia. | FK → rola_multitasking.id, NULL dopuszczalny |
| tresc | TEXT | Treść zlecenia w postaci przekazanej Koordynatorowi. | NOT NULL |
| kryteria_akceptacji | TEXT | Kryteria, wobec których prowadzona jest kontrola jakości. | JSON |
| priorytet | INTEGER | Kolejność obsługi zleceń. | — |
| status | TEXT | `przyjete` \| `w_dekompozycji` \| `w_realizacji` \| `wstrzymane` \| `zakonczone` \| `przerwane` \| `blad`. | CHECK enum |
| limit_ponowien | INTEGER | Dopuszczalna liczba ponowień zadania w ramach zlecenia. | NOT NULL |
| data_utworzenia | TEXT | — | NOT NULL |
| data_zakonczenia | TEXT | — | NULL dopuszczalny |

*Relacje.* Zlecenie gromadzi `zadanie` (19.4) i `przebieg_petli` (19.5). Korekta zlecenia wydana z okna pętli wykonawczej jest wierszem `komunikat_sterujacy` o rodzaju `korekta_zlecenia` (19.6).

### 19.4. Zadanie i kolejka zadań

Dekompozycja zlecenia na jednostki wykonawcze korzysta z encji `zadanie` (rozdział 10.1) i `kolejka` (rozdział 10.3). Encja `zadanie` niesie kolumny wiążące ją z warstwą komunikacji operacyjnej:

| Kolumna encji `zadanie` | Opis | Klucz |
|---|---|---|
| zlecenie_id | Zlecenie, z którego dekompozycji zadanie powstało. | FK → zlecenie.id, ON DELETE CASCADE, NULL dopuszczalny |
| przebieg_petli_id | Przebieg pętli wykonawczej, w którym zadanie jest realizowane. | FK → przebieg_petli.id, ON DELETE SET NULL, NULL dopuszczalny |
| zadanie_nadrzedne_id | Zadanie nadrzędne w dekompozycji (samoreferencja). | FK → zadanie.id, NULL dopuszczalny |
| liczba_ponowien | Licznik wykonanych ponowień zadania. | NOT NULL, domyślnie `0` |

Kolejka zadań (`kolejka`, rozdział 10.3) porządkuje wykonanie zadań jednego lub wielu zleceń; pozycja zadania w kolejce jest wierszem `pozycja_kolejki` (rozdział 10.4), a każda operacja silnika kolejek — wierszem `log_akcji_kolejki` (rozdział 10.5). Kolumna `kolejka.zlecenie_id` wiąże kolejkę dedykowaną jednemu zleceniu:

| Kolumna encji `kolejka` | Opis | Klucz |
|---|---|---|
| zlecenie_id | Zlecenie, któremu kolejka jest dedykowana. | FK → zlecenie.id, ON DELETE SET NULL, NULL dopuszczalny |

### 19.5. Przebieg pętli wykonawczej (`przebieg_petli`)

Pojedyncza iteracja pętli wykonawczej prowadzonej przez Koordynatora nad Wykonawcami — treść okna Execution Loop Window. Kolejne przebiegi tego samego zlecenia różnią się numerem iteracji.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator przebiegu. | PK |
| zlecenie_id | INTEGER | Zlecenie realizowane w przebiegu. | FK → zlecenie.id, ON DELETE CASCADE |
| sesja_id | INTEGER | Sesja, w której przebieg jest obserwowany. | FK → sesja.id, NULL dopuszczalny |
| koordynator_rola_id | INTEGER | Rola prowadząca pętlę. | FK → rola_multitasking.id, NULL dopuszczalny |
| przebieg_multitasking_id | INTEGER | Przebieg środowiska MultitaskingAI, gdy pętla jest prowadzona przez zespół wielomodelowy. | FK → przebieg_multitasking.id, NULL dopuszczalny |
| numer_iteracji | INTEGER | Kolejny numer iteracji w ramach zlecenia. | NOT NULL |
| status | TEXT | `w_trakcie` \| `wstrzymany` \| `zakonczony` \| `przerwany` \| `blad`. | CHECK enum |
| wskazniki | TEXT | Wskaźniki przebiegu pętli — liczba zadań, zadania zakończone, ponowienia, czas trwania. | JSON |
| podsumowanie | TEXT | Podsumowanie iteracji. | NULL dopuszczalny |
| data_rozpoczecia | TEXT | — | NOT NULL |
| data_zakonczenia | TEXT | — | NULL dopuszczalny |

Ograniczenie `UNIQUE(zlecenie_id, numer_iteracji)`.

Szablon zawartości pola `wskazniki` (kolumna JSON):

```json
{
  "zadania_lacznie": 12,
  "zadania_zakonczone": 9,
  "zadania_w_trakcie": 2,
  "ponowienia": 1,
  "czas_trwania_sekundy": 348
}
```

### 19.6. Komunikat sterujący (`komunikat_sterujacy`)

Pojedynczy wpis wymiany między Koordynatorem a Wykonawcą oraz sterowanie przebiegiem wydawane z okna pętli wykonawczej.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator komunikatu. | PK |
| przebieg_petli_id | INTEGER | Przebieg, w którym komunikat powstał. | FK → przebieg_petli.id, ON DELETE CASCADE |
| zlecenie_id | INTEGER | Zlecenie, którego komunikat dotyczy. | FK → zlecenie.id, NULL dopuszczalny |
| zadanie_id | INTEGER | Zadanie, którego komunikat dotyczy. | FK → zadanie.id, ON DELETE SET NULL, NULL dopuszczalny |
| nadawca | TEXT | `koordynator` \| `wykonawca` \| `uzytkownik`. | CHECK enum |
| rola_multitasking_id | INTEGER | Rola będąca nadawcą komunikatu. | FK → rola_multitasking.id, NULL dopuszczalny |
| rodzaj | TEXT | `przydzial_zadania` \| `raport_postepu` \| `zapytanie_o_kontekst` \| `odpowiedz_kontekstowa` \| `zatwierdzenie` \| `korekta_zlecenia` \| `wstrzymanie` \| `wznowienie` \| `przerwanie` \| `zakonczenie`. | CHECK enum |
| tresc | TEXT | Treść komunikatu. | NOT NULL |
| ladunek | TEXT | Dane strukturalne komunikatu (parametry przydziału, wskaźniki postępu, zakres korekty). | JSON |
| znacznik_czasu | TEXT | — | NOT NULL |

Dziesięć rodzajów komunikatu sterującego (`rodzaj`):

| Rodzaj | Znaczenie |
|---|---|
| przydzial_zadania | Koordynator przydziela zadanie Wykonawcy. |
| raport_postepu | Wykonawca zgłasza stan realizacji zadania. |
| zapytanie_o_kontekst | Wykonawca żąda uzupełnienia kontekstu lub decyzji. |
| odpowiedz_kontekstowa | Koordynator uzupełnia kontekst lub rozstrzyga wątpliwość. |
| zatwierdzenie | Zatwierdzenie wyniku zadania lub zlecenia. |
| korekta_zlecenia | Zmiana treści lub kryteriów zlecenia w trakcie realizacji. |
| wstrzymanie | Wstrzymanie przebiegu pętli. |
| wznowienie | Wznowienie wstrzymanego przebiegu. |
| przerwanie | Przerwanie przebiegu bez ukończenia zlecenia. |
| zakonczenie | Zamknięcie przebiegu po spełnieniu kryteriów akceptacji. |

Szablon zawartości pola `ladunek` (kolumna JSON) dla rodzaju `przydzial_zadania`:

```json
{
  "rodzaj": "przydzial_zadania",
  "zadanie_id": 214,
  "wykonawca_rola_id": 6,
  "kanal_modelu_id": 3,
  "termin_sekundy": 600
}
```

### 19.7. Wynik kontroli jakości (`wynik_kontroli_jakosci`)

Ocena rezultatu zadania wobec kryteriów akceptacji zlecenia, wystawiana przez rolę walidującą i rozstrzygająca o zamknięciu zadania albo o ponowieniu.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator wyniku. | PK |
| zadanie_id | INTEGER | Zadanie objęte kontrolą. | FK → zadanie.id, ON DELETE CASCADE |
| przebieg_petli_id | INTEGER | Przebieg, w którym kontrola została przeprowadzona. | FK → przebieg_petli.id, ON DELETE CASCADE |
| rola_multitasking_id | INTEGER | Rola przeprowadzająca kontrolę. | FK → rola_multitasking.id, NULL dopuszczalny |
| ocena | TEXT | `przyjete` \| `do_poprawy` \| `odrzucone`. | CHECK enum |
| kryteria | TEXT | Wynik oceny w rozbiciu na poszczególne kryteria akceptacji. | JSON |
| uzasadnienie | TEXT | Opis uzasadniający ocenę. | — |
| znacznik_czasu | TEXT | — | NOT NULL |

Szablon zawartości pola `kryteria` (kolumna JSON):

```json
{
  "kryteria": [
    { "nazwa": "zgodnosc_z_trescia_zlecenia", "spelnione": true },
    { "nazwa": "kompletnosc_wyniku", "spelnione": false, "uwaga": "brak sekcji podsumowania" }
  ]
}
```

### 19.8. Ponowienie (`ponowienie`)

Zapis powtórnego skierowania zadania do realizacji po ocenie negatywnej, błędzie wykonania lub korekcie zlecenia.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator ponowienia. | PK |
| zadanie_id | INTEGER | Ponawiane zadanie. | FK → zadanie.id, ON DELETE CASCADE |
| wynik_kontroli_id | INTEGER | Wynik kontroli jakości, który wywołał ponowienie. | FK → wynik_kontroli_jakosci.id, ON DELETE SET NULL, NULL dopuszczalny |
| przebieg_petli_id | INTEGER | Przebieg, w którym ponowienie zostało zlecone. | FK → przebieg_petli.id, ON DELETE CASCADE |
| numer_ponowienia | INTEGER | Kolejny numer ponowienia zadania. | NOT NULL |
| przyczyna | TEXT | `negatywna_kontrola_jakosci` \| `blad_wykonania` \| `przekroczony_czas` \| `korekta_zlecenia`. | CHECK enum |
| zakres_poprawki | TEXT | Opis zakresu poprawki przekazany Wykonawcy. | — |
| znacznik_czasu | TEXT | — | NOT NULL |

Ograniczenie `UNIQUE(zadanie_id, numer_ponowienia)`. Osiągnięcie wartości `zlecenie.limit_ponowien` zamyka zadanie statusem `blad` i wywołuje komunikat sterujący rodzaju `zapytanie_o_kontekst` kierowany do Użytkownika.

### 19.9. Indeksy warstwy komunikacji operacyjnej

| Tabela | Kolumny objęte jawnym indeksem |
|---|---|
| `rozmowa` | `sesja_id`, `data_ostatniej_wiadomosci` |
| `wiadomosc` | `rozmowa_id` |
| `zlecenie` | `rozmowa_id`, `status` |
| `przebieg_petli` | `zlecenie_id`, `status` |
| `komunikat_sterujacy` | `przebieg_petli_id`, `zadanie_id`, `znacznik_czasu` |
| `wynik_kontroli_jakosci` | `zadanie_id`, `przebieg_petli_id` |
| `ponowienie` | `zadanie_id` |
| `zadanie` | `zlecenie_id`, `przebieg_petli_id` |

### 19.10. Retencja danych komunikacji operacyjnej

Okresy przechowywania są wierszami `ustawienie` (rozdział 17.1) i podlegają regule ośmiu poziomów zasięgu z rozdziału 15.5.

| Klucz `ustawienie` | Zakres danych | Wartość przy braku ustawienia |
|---|---|---|
| `retencja_historii` | `rozmowa`, `wiadomosc`, `zalacznik_wiadomosci` | Przechowywanie bez limitu |
| `retencja_zlecen` | `zlecenie`, `przebieg_petli` | Przechowywanie bez limitu |
| `retencja_komunikatow_sterujacych` | `komunikat_sterujacy` | 90 dni od zamknięcia przebiegu |
| `retencja_kontroli_jakosci` | `wynik_kontroli_jakosci`, `ponowienie` | Okres równy retencji zlecenia nadrzędnego |

Usunięcie rozmowy po upływie okresu retencji kasuje kaskadowo jej wiadomości i załączniki, a zlecenia powstałe z rozmowy zachowują się z pustą wartością `rozmowa_id` do czasu upływu własnego okresu retencji. Artefakty (rozdział 16.1) nie podlegają retencji komunikacji — ich usunięcie następuje wyłącznie na polecenie Operatora.

---

## 20. Konfiguracja warstw widoczności

Rozdział opisuje dane sterujące stopniowym ujawnianiem funkcjonalności interfejsu: przypisanie elementu interfejsu do jednej z czterech warstw widoczności, konfigurację warstw według roli użytkownika oraz rejestr funkcji zasilający wyszukiwarkę funkcji.

Diagram relacji (klucze obce) grupy warstw widoczności:

```
warstwa_widocznosci ──1:N──> element_interfejsu ──N:1──> modul | srodowisko | okno_operacyjne
                                    │
                                    ├──1:N──> konfiguracja_warstwy_roli ──N:1──> rola_uzytkownika
                                    └──1:N──> funkcja_platformy ──N:1──> warstwa_widocznosci
```

### 20.1. Warstwa widoczności (`warstwa_widocznosci`)

Tabela słownikowa, cztery stałe wiersze.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator warstwy. | PK |
| numer | INTEGER | Numer warstwy `1`–`4`. | UNIQUE, CHECK (1,2,3,4) |
| kod | TEXT | `zawsze_widoczna` \| `na_zadanie` \| `rozwiniecia_kontekstowe` \| `funkcje_eksperckie`. | UNIQUE, CHECK enum |
| nazwa | TEXT | Nazwa warstwy prezentowana w oknie konfiguracji. | NOT NULL |
| opis | TEXT | Opis zawartości warstwy. | — |

Cztery stałe wiersze słownika:

| `numer` | `kod` | `nazwa` | Sposób dostępu |
|---|---|---|---|
| 1 | zawsze_widoczna | Zawsze widoczna | Widoczna bez interakcji |
| 2 | na_zadanie | Widoczna na żądanie | Ikona, przycisk, przełącznik, znacznik kontekstowy |
| 3 | rozwiniecia_kontekstowe | Rozwinięcia kontekstowe | Menu kebab, menu hamburger, menu kontekstowe, panel popover, lista rozwijana |
| 4 | funkcje_eksperckie | Funkcje eksperckie | Polecenie języka naturalnego, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli |

### 20.2. Element interfejsu (`element_interfejsu`)

Pojedynczy element interfejsu — okno, panel, pasek narzędzi, menu, znacznik kontekstu lub funkcja — wraz z przypisaną warstwą widoczności i sposobem wywołania.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator elementu. | PK |
| kod | TEXT | Kod elementu, niepowtarzalny w skali platformy. | UNIQUE, NOT NULL |
| nazwa | TEXT | Nazwa elementu prezentowana w interfejsie. | NOT NULL |
| rodzaj | TEXT | `okno` \| `panel` \| `pasek_narzedzi` \| `menu` \| `znacznik_kontekstu` \| `funkcja`. | CHECK enum |
| warstwa_id | INTEGER | Warstwa widoczności elementu. | FK → warstwa_widocznosci.id |
| srodowisko_id | INTEGER | Środowisko, w którym element występuje; puste dla elementów wspólnych. | FK → srodowisko.id, NULL dopuszczalny |
| modul_id | INTEGER | Moduł, w którym element występuje; puste dla elementów wspólnych. | FK → modul.id, NULL dopuszczalny |
| okno_operacyjne_id | INTEGER | Okno operacyjne, do którego element należy. | FK → okno_operacyjne.id, NULL dopuszczalny |
| element_nadrzedny_id | INTEGER | Element nadrzędny, gdy element jest pozycją menu lub panelu (samoreferencja). | FK → element_interfejsu.id, NULL dopuszczalny |
| sposob_dostepu | TEXT | `bez_interakcji` \| `ikona` \| `przycisk` \| `przelacznik` \| `znacznik_kontekstowy` \| `menu_kebab` \| `menu_hamburger` \| `menu_kontekstowe` \| `panel_popover` \| `lista_rozwijana` \| `polecenie_jezyka_naturalnego` \| `skrot_klawiszowy` \| `wyszukiwarka_funkcji` \| `tryb_administracyjny` \| `konfiguracja_roli`. | CHECK enum |
| zwija_sie_samoczynnie | INTEGER | Czy element wraca do stanu zwiniętego po użyciu. | CHECK (0,1) |
| kolejnosc | INTEGER | Kolejność elementu w obrębie elementu nadrzędnego. | — |

*Uwaga projektowa.* Element należy do dokładnie jednej warstwy widoczności. Elementy warstwy 1 są widoczne w stanie spoczynku interfejsu; elementy warstw 2–4 są wywoływane sposobem wskazanym w kolumnie `sposob_dostepu`, przy czym każdy z nich jest osiągalny jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Głębokość hierarchii `element_nadrzedny_id` jest ograniczona do jednego poziomu zagnieżdżenia i egzekwowana przez rdzeń serwera.

### 20.3. Rola użytkownika (`rola_uzytkownika`)

Tabela słownikowa ról użytkownika platformy, sterujących zestawem widocznych elementów interfejsu. Encja jest odrębna od `rola_multitasking` (rozdział 13.2), która opisuje rolę wykonawczą w środowisku wielomodelowym.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator roli. | PK |
| kod | TEXT | `podstawowy` \| `zaawansowany` \| `administrator`. | UNIQUE, CHECK enum |
| nazwa | TEXT | Nazwa roli prezentowana w oknie konfiguracji. | NOT NULL |
| najwyzsza_warstwa_id | INTEGER | Najwyższa warstwa widoczności dostępna roli. | FK → warstwa_widocznosci.id |

Trzy stałe wiersze słownika:

| `kod` | `nazwa` | Najwyższa dostępna warstwa |
|---|---|---|
| podstawowy | Użytkownik podstawowy | 3 |
| zaawansowany | Użytkownik zaawansowany | 4 |
| administrator | Administrator | 4 |

### 20.4. Konfiguracja warstwy roli (`konfiguracja_warstwy_roli`)

Nadpisanie warstwy widoczności elementu interfejsu dla wskazanej roli użytkownika. Brak wiersza oznacza obowiązywanie warstwy z `element_interfejsu.warstwa_id`.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | — | PK |
| rola_uzytkownika_id | INTEGER | Rola, której konfiguracja dotyczy. | FK → rola_uzytkownika.id, ON DELETE CASCADE |
| element_interfejsu_id | INTEGER | Element objęty nadpisaniem. | FK → element_interfejsu.id, ON DELETE CASCADE |
| warstwa_id | INTEGER | Warstwa obowiązująca dla tej roli. | FK → warstwa_widocznosci.id |
| widoczny | INTEGER | Czy element jest udostępniany roli. | CHECK (0,1) |
| data_modyfikacji | TEXT | — | — |

Ograniczenie `UNIQUE(rola_uzytkownika_id, element_interfejsu_id)`.

Reguła rozstrzygania warstwy efektywnej elementu dla zalogowanej roli:

1. Wiersz `konfiguracja_warstwy_roli` dla pary rola–element; jeżeli `widoczny = 0`, element nie jest prezentowany.
2. Przy braku wiersza — wartość `element_interfejsu.warstwa_id`.
3. Warstwa wyższa niż `rola_uzytkownika.najwyzsza_warstwa_id` powoduje ukrycie elementu przed daną rolą.

### 20.5. Funkcja platformy (`funkcja_platformy`)

Rejestr funkcji zasilający wyszukiwarkę funkcji — jedno z narzędzi dostępu do elementów warstwy 4. Rejestr obejmuje wszystkie funkcje platformy niezależnie od warstwy, w której są prezentowane.

| Pole | Typ | Opis | Klucz |
|---|---|---|---|
| id | INTEGER | Identyfikator funkcji. | PK |
| kod | TEXT | Kod funkcji, niepowtarzalny w skali platformy. | UNIQUE, NOT NULL |
| nazwa | TEXT | Nazwa funkcji prezentowana w wynikach wyszukiwania. | NOT NULL |
| opis | TEXT | Opis działania funkcji. | — |
| slowa_kluczowe | TEXT | Zbiór określeń, po których funkcja jest odnajdywana. | JSON |
| element_interfejsu_id | INTEGER | Element interfejsu uruchamiający funkcję. | FK → element_interfejsu.id, ON DELETE SET NULL, NULL dopuszczalny |
| warstwa_id | INTEGER | Warstwa widoczności funkcji. | FK → warstwa_widocznosci.id |
| srodowisko_id | INTEGER | Środowisko, w którym funkcja działa; puste dla funkcji globalnych. | FK → srodowisko.id, NULL dopuszczalny |
| modul_id | INTEGER | Moduł, w którym funkcja działa; puste dla funkcji globalnych. | FK → modul.id, NULL dopuszczalny |
| polecenie_jezyka_naturalnego | TEXT | Wzorzec polecenia wywołującego funkcję z okna komunikacji. | NULL dopuszczalny |
| skrot_klawiszowy | TEXT | Skrót klawiszowy wywołujący funkcję. | NULL dopuszczalny |
| widoczna_w_wyszukiwarce | INTEGER | Czy funkcja jest zwracana w wynikach wyszukiwarki funkcji. | CHECK (0,1) |

Szablon zawartości pola `slowa_kluczowe` (kolumna JSON):

```json
{
  "slowa_kluczowe": ["izolacja", "punkt izolacji", "profil", "zasięg"],
  "synonimy": ["rozdzielenie kontekstu"]
}
```

### 20.6. Indeksy i retencja grupy warstw widoczności

| Tabela | Kolumny objęte jawnym indeksem |
|---|---|
| `element_interfejsu` | `warstwa_id`, `modul_id`, `srodowisko_id` |
| `konfiguracja_warstwy_roli` | `rola_uzytkownika_id`, `element_interfejsu_id` |
| `funkcja_platformy` | `warstwa_id`, `modul_id`, `widoczna_w_wyszukiwarce` |

Tabele `warstwa_widocznosci`, `element_interfejsu`, `rola_uzytkownika` i `funkcja_platformy` są tabelami słownikowymi (rozdział 22.5) i nie podlegają retencji czasowej. Tabela `konfiguracja_warstwy_roli` przechowuje dane Operatora i jest kasowana wyłącznie wraz z elementem lub rolą, do których się odnosi.

---

## 21. Zbiorcze powiązania encji

Poniższa tabela zbiera w jednym miejscu relacje między głównymi encjami modelu, z regułą kasowania kaskadowego tam, gdzie encja zależna traci sens bez encji nadrzędnej.

| Encja nadrzędna | Encja zależna | Charakter relacji | Reguła kasowania |
|---|---|---|---|
| `konto` | `urzadzenie` | jeden‑do‑wielu | CASCADE |
| `urzadzenie` | `polaczenie` | jeden‑do‑wielu | CASCADE |
| `urzadzenie` | `poswiadczenie_urzadzenia`, `dzialanie_lokalne` | jeden‑do‑wielu | CASCADE |
| `konto` | `kod_parowania` | jeden‑do‑wielu | CASCADE |
| `konto` | `karta_sesji` | jeden‑do‑wielu | CASCADE |
| `karta_sesji` | `sesja` | jeden‑do‑wielu (historyczne przełączenia modułu) | CASCADE |
| `sesja` | `proces_sesji` | jeden‑do‑jednego | CASCADE |
| `sesja` | `wiadomosc`, `sesja_okno`, `zadanie` (NULL dopuszczalny) | jeden‑do‑wielu | CASCADE dla `wiadomosc`, `sesja_okno`; `SET NULL` dla `zadanie.sesja_id` |
| `sesja` | `rozmowa` | jeden‑do‑wielu | CASCADE |
| `rozmowa` | `wiadomosc`, `zlecenie` | jeden‑do‑wielu | CASCADE dla `wiadomosc`; `SET NULL` dla `zlecenie.rozmowa_id` |
| `wiadomosc` | `zalacznik_wiadomosci` | jeden‑do‑wielu | CASCADE |
| `zlecenie` | `przebieg_petli`, `zadanie`, `kolejka` (NULL dopuszczalny) | jeden‑do‑wielu | CASCADE dla `przebieg_petli`, `zadanie`; `SET NULL` dla `kolejka.zlecenie_id` |
| `przebieg_petli` | `komunikat_sterujacy`, `wynik_kontroli_jakosci`, `ponowienie`, `zadanie` (NULL dopuszczalny) | jeden‑do‑wielu | CASCADE dla trzech pierwszych; `SET NULL` dla `zadanie.przebieg_petli_id` |
| `zadanie` | `wynik_kontroli_jakosci`, `ponowienie`, `komunikat_sterujacy` | jeden‑do‑wielu | CASCADE dla `wynik_kontroli_jakosci`, `ponowienie`; `SET NULL` dla `komunikat_sterujacy.zadanie_id` |
| `wynik_kontroli_jakosci` | `ponowienie` | jeden‑do‑wielu | SET NULL |
| `warstwa_widocznosci` | `element_interfejsu`, `konfiguracja_warstwy_roli`, `funkcja_platformy`, `rola_uzytkownika` | jeden‑do‑wielu | RESTRICT (encja słownikowa) |
| `element_interfejsu` | `konfiguracja_warstwy_roli`, `funkcja_platformy` | jeden‑do‑wielu | CASCADE dla `konfiguracja_warstwy_roli`; `SET NULL` dla `funkcja_platformy.element_interfejsu_id` |
| `rola_uzytkownika` | `konfiguracja_warstwy_roli` | jeden‑do‑wielu | CASCADE |
| `srodowisko` | `srodowisko_modul`, `sesja`, `karta_sesji`, `projekt_srodowisko` | jeden‑do‑wielu | RESTRICT (encje słownikowe) |
| `modul` | `srodowisko_modul`, `okno_operacyjne`, `sesja` | jeden‑do‑wielu | RESTRICT / SET NULL (`sesja.modul_id`) |
| `komponent_wlasny` | `automatyka` \| `agent` \| `projekt` \| `profil_asystenta` (dokładnie jedna, wg `rodzaj`) | jeden‑do‑jednego | CASCADE |
| `projekt` | `projekt_srodowisko`, `projekt_agent`, `sesja` (NULL dopuszczalny), `artefakt`, `kolekcja` | jeden‑do‑wielu | CASCADE dla `projekt_srodowisko`, `projekt_agent`; `SET NULL` dla `sesja.projekt_id`, `artefakt.projekt_id` |
| `agent` | `agent_umiejetnosc`, `agent_rozszerzenie`, `agent_uprawnienie`, `projekt_agent` | jeden‑do‑wielu | CASCADE |
| `zadanie` | `pozycja_kolejki`, `komponent_kompozycji` | jeden‑do‑wielu | CASCADE / SET NULL |
| `kolejka` | `pozycja_kolejki`, `log_akcji_kolejki` | jeden‑do‑wielu | CASCADE |
| `zespol_multitasking` | `rola_multitasking`, `kolejka`, `zaleznosc_orkiestracji`, `przebieg_multitasking` | jeden‑do‑wielu (NULL dopuszczalny) | SET NULL |
| `rola_multitasking` | `subagent`, `decyzja_procesu`, `zadanie` (NULL dopuszczalny) | jeden‑do‑wielu | CASCADE dla `subagent`; SET NULL dla pozostałych |
| `przebieg_multitasking` | `decyzja_procesu` | jeden‑do‑wielu | CASCADE |
| `profil_izolacji` | `regula_izolacji_kontekstu`, `regula_izolacji_technicznej`, `przypisanie_profilu_izolacji` | jeden‑do‑wielu | CASCADE |
| `poziom_zasiegu` | `profil_izolacji`, `ustawienie` | jeden‑do‑wielu | RESTRICT (encja słownikowa) |
| `artefakt` | `wersja_artefaktu`, `kolekcja_artefakt`, `etykieta_artefaktu` | jeden‑do‑wielu | CASCADE |
| `rozszerzenie` | `agent_rozszerzenie` | jeden‑do‑wielu | CASCADE |
| `strefa_glowna` | `pozycja_strefy` | jeden‑do‑wielu | CASCADE |

---

## 22. Zgodność ze schematem SQLite

### 22.1. Klucze i tryb pracy

| Aspekt | Rozstrzygnięcie |
|---|---|
| Klucz główny | `id INTEGER PRIMARY KEY AUTOINCREMENT` (alias `rowid`) — jedyny wzorzec w modelu. |
| Nadawanie identyfikatorów | Wyłącznie serwer (architektura, rozdziały 2 i 10); klienci nie tworzą rekordów offline ani nie przechowują stanu trwałego. |
| Wymuszalność kluczy obcych | `PRAGMA foreign_keys = ON` ustawiane dla każdego połączenia rdzenia z bazą. |

*Uzasadnienie.* Alias `rowid` jest najwydajniejszym dostępnym kluczem głównym i nie wymaga identyfikatorów globalnie unikalnych między urządzeniami, ponieważ serwer jest jedynym źródłem prawdy i jedynym miejscem nadawania identyfikatorów.

### 22.2. Indeksy

Poza indeksami niejawnymi wynikającymi z ograniczeń `UNIQUE` i `PRIMARY KEY` warstwa danych zakłada jawne indeksy na kolumnach klucza obcego najczęściej filtrowanych w zapytaniach interfejsu:

| Tabela | Kolumny objęte jawnym indeksem |
|---|---|
| `sesja` | `karta_sesji_id` |
| `wiadomosc` | `sesja_id`, `rozmowa_id` |
| `zadanie` | `kolejka_id`, `zlecenie_id`, `przebieg_petli_id` |
| `pozycja_kolejki` | `kolejka_id` |
| `zasob_pamieci` | `poziom`, `odniesienie_id` |
| `ustawienie` | `klucz`, `poziom_zasiegu_id`, `poziom_odniesienie_id` |
| `przypisanie_profilu_izolacji` | `typ_celu`, `cel_id` |
| `rozmowa` | `sesja_id`, `data_ostatniej_wiadomosci` |
| `zlecenie` | `rozmowa_id`, `status` |
| `przebieg_petli` | `zlecenie_id`, `status` |
| `komunikat_sterujacy` | `przebieg_petli_id`, `zadanie_id`, `znacznik_czasu` |
| `wynik_kontroli_jakosci` | `zadanie_id`, `przebieg_petli_id` |
| `ponowienie` | `zadanie_id` |
| `element_interfejsu` | `warstwa_id`, `modul_id`, `srodowisko_id` |
| `konfiguracja_warstwy_roli` | `rola_uzytkownika_id`, `element_interfejsu_id` |
| `funkcja_platformy` | `warstwa_id`, `modul_id`, `widoczna_w_wyszukiwarce` |

Rdzeń serwera (Go) zakłada te indeksy przy inicjalizacji schematu bazy.

### 22.3. Pola JSON

Kolumny oznaczone typem „JSON” są kolumnami `TEXT` przechowującymi tekst zgodny ze składnią JSON, odczytywane i zapisywane funkcjami rozszerzenia JSON1 wbudowanego w SQLite (`json_extract`, `json_set`, `json_valid`). Rozwiązanie stosuje się do pól o zmiennej, niestandaryzowanej strukturze, których sztywne rozbicie na kolumny relacyjne ograniczałoby kompozycyjność bez korzyści dla integralności danych. Pełny wykaz kolumn JSON w modelu:

| Tabela.Kolumna | Przeznaczenie |
|---|---|
| `sesja.kontekst` | Bieżący kontekst przekazywany do AI. |
| `sesja_okno.stan` | Stan specyficzny dla okna operacyjnego. |
| `proces_sesji.zmienne_srodowiskowe` | Środowisko procesu sesji. |
| `automatyka.definicja_workflow` | Definicja procesu z Workflow Builder. |
| `automatyka.wywolania_modeli` | Konfiguracja zaplanowanych wywołań modeli. |
| `profil_asystenta.zakres_polecen` | Dozwolony zakres poleceń głosowych. |
| `profil_asystenta.obslugiwane_akcje` | Lista obsługiwanych akcji wieloetapowych. |
| `komponent_kompozycji.definicja` | Definicja skomponowanego układu pracy. |
| `log_akcji_kolejki.parametry` | Parametry akcji silnika kolejek. |
| `zaleznosc_orkiestracji.warunek` | Warunek dla `typ_zaleznosci = 'warunkuje'`. |
| `agent.konfiguracja_pamieci` | Konfiguracja pamięci agenta. |
| `agent.konfiguracja_uprawnien` | Uprawnienia agenta w formie skróconej. |
| `agent_umiejetnosc.konfiguracja` | Konfiguracja umiejętności. |
| `agent_rozszerzenie.konfiguracja` | Konfiguracja przypisania rozszerzenia. |
| `rozszerzenie.konfiguracja` | Konfiguracja rozszerzenia. |
| `ustawienie.wartosc` | Wartość parametru (JSON dopuszczalny dla wartości złożonych). |
| `rozmowa.kontekst` | Kontekst rozmowy przekazywany Wykonawcy. |
| `zlecenie.kryteria_akceptacji` | Kryteria oceny wyniku zlecenia. |
| `przebieg_petli.wskazniki` | Wskaźniki przebiegu pętli wykonawczej. |
| `komunikat_sterujacy.ladunek` | Dane strukturalne komunikatu sterującego. |
| `wynik_kontroli_jakosci.kryteria` | Wynik oceny w rozbiciu na kryteria akceptacji. |
| `funkcja_platformy.slowa_kluczowe` | Określenia i synonimy dla wyszukiwarki funkcji. |

### 22.4. Odniesienia polimorficzne

Wymienione niżej tabele wykorzystują parę kolumn `typ_* / *_id` zamiast ograniczenia `FOREIGN KEY`, ponieważ jedna kolumna musi wskazywać wiersze różnych tabel w zależności od wartości kolumny towarzyszącej.

| Tabela | Para kolumn polimorficznych | Możliwe cele |
|---|---|---|
| `pozycja_strefy` | `typ_odniesienia` / `kod_odniesienia` | srodowisko, rodzaj komponentu własnego, funkcja ustawień |
| `powiazanie_komponentu` | `zrodlo_typ` / `zrodlo_id`, `cel_typ` / `cel_id` | modul, komponent_wlasny, srodowisko |
| `zaleznosc_orkiestracji` | `zrodlo_typ` / `zrodlo_id`, `cel_typ` / `cel_id` | kanal_modelu, agent, zadanie, kolejka, automatyka, projekt |
| `zasob_pamieci` | `poziom` / `odniesienie_id` | srodowisko, modul, projekt, sesja |
| `przypisanie_profilu_izolacji` | `typ_celu` / `cel_id` | srodowisko, modul, para_modulow, projekt, karta_sesji, rola_multitasking |
| `ustawienie` | `poziom` (przez `poziom_zasiegu`) / `poziom_odniesienie_id` | jak w `przypisanie_profilu_izolacji` |
| `sugestia_aod` | `kontekst_typ` / `kontekst_id` | srodowisko, modul, projekt, sesja |
| `powiadomienie` | `zrodlo_typ` / `zrodlo_id` | srodowisko, modul, projekt, sesja, zadanie, przebieg_petli, automatyka |
| `kanal_modelu` | `przypisanie_typ` / `przypisanie_id` | sesja, rola_multitasking |

SQLite nie wspiera kluczy obcych polimorficznych deklaratywnie — integralność tych odniesień jest egzekwowana na poziomie rdzenia serwera (Go), w warstwie rdzenia opisanej w architekturze (rozdział 4), a nie ograniczeniem schematu bazy.

### 22.5. Tabele słownikowe a tabele danych Operatora

Model rozróżnia dwie kategorie tabel:

| Kategoria | Tabele | Charakter zawartości |
|---|---|---|
| Tabele słownikowe | `srodowisko`, `modul`, `okno_operacyjne`, `strefa_glowna`, `poziom_zasiegu`, `sekcja_panelu_orkiestracji`, `warstwa_widocznosci`, `rola_uzytkownika`, `element_interfejsu`, `funkcja_platformy` | Zasiewane przy inicjalizacji bazy wartościami ustalonymi dla platformy; zbiór wierszy stały, natomiast pola konfigurowalne — `dostepny`, `widoczna`, `kolejnosc` — w pełni edytowalne przez Operatora. |
| Tabele danych Operatora | Pozostałe | Rosną w trakcie użytkowania: sesje, wiadomości, komponenty własne, zadania, profile izolacji i inne encje tworzone przez Operatora lub przez AI w jego imieniu. |

*Polityka migracji.* Rozróżnienie to determinuje politykę migracji: zmiana zbioru wierszy tabeli słownikowej — dodanie nowego środowiska — wymaga migracji schematu, natomiast zmiana wartości pól konfigurowalnych w tych tabelach jest zwykłą operacją zapisu, nie migracją.

---

## 23. Kryteria odbioru

Model danych jest zrealizowany zgodnie z niniejszym opracowaniem, gdy spełnione są łącznie poniższe warunki. Każdy warunek jest sprawdzalny bez odwołania do intencji autora — wyłącznie przez odczyt schematu bazy danych albo kontraktu.

| Warunek | Sposób sprawdzenia |
|---|---|
| Każda tabela danych Operatora używa klucza głównego `id INTEGER PRIMARY KEY AUTOINCREMENT` (rozdz. 22.1) | odczyt schematu SQLite — `PRAGMA table_info` nie zwraca tabeli danych Operatora z kluczem głównym odmiennego typu |
| Identyfikatory nadaje wyłącznie serwer; klient nie tworzy rekordów offline (rozdz. 22.1) | brak w kliencie ścieżki kodu zapisującej wiersz bezpośrednio do pliku bazy danych |
| `PRAGMA foreign_keys = ON` jest ustawione dla każdego połączenia rdzenia z bazą (rozdz. 22.1) | odczyt inicjalizacji połączenia w kodzie rdzenia (Go) |
| Kolumny oznaczone typem JSON są kolumnami `TEXT` zgodnymi ze składnią funkcji JSON1 SQLite, z walidacją `json_valid()` (rozdz. 22.3) | próba zapisu wartości niezgodnej ze składnią JSON do kolumny z wykazu 22.3 kończy się odmową na poziomie ograniczenia `CHECK` |
| Odniesienia polimorficzne (rozdz. 22.4) są egzekwowane w warstwie rdzenia, nie ograniczeniem `FOREIGN KEY` schematu | schemat SQLite nie deklaruje `FOREIGN KEY` na kolumnach `*_typ` / `*_id` wymienionych w tabeli 22.4 |
| Zbiór wierszy tabeli słownikowej (rozdz. 22.5) zmienia się wyłącznie migracją schematu, nigdy operacją zapisu z poziomu interfejsu | interfejs Operatora nie udostępnia komendy dodającej wiersz do tabeli `srodowisko`, `modul`, `warstwa_widocznosci` ani innej tabeli słownikowej wykazu 22.5 |
| Każda grupa encji wymieniona w tabeli wprowadzającej (rozdz. 1 wprowadzenia) jest opisana w rozdziale wskazanym w tej tabeli, bez grupy pominiętej milczeniem | liczba rozdziałów 3–20 odpowiada liczbie wierszy tabeli wprowadzającej |
| Pole struktury kontraktu oznaczone jako przelotowe (nieprzenoszone do bazy) jest jawnie tak oznaczone w opisie encji, nie pomijane milczeniem | przegląd pól struktur `budowa/shared/contract.json` bez odpowiednika kolumny bazy w niniejszym dokumencie kończy się odnalezieniem adnotacji „poza bazą” albo „przelotowe” |

---

## Załącznik A. Słownik enumeracji

Zestawienie porządkuje wszystkie pola typu wyliczeniowego (`TEXT` z ograniczeniem `CHECK`) użyte w modelu, wraz z tabelą i kolumną, w której występują.

| Tabela.Kolumna | Dozwolone wartości |
|---|---|
| `urzadzenie.typ` | komputer, telefon, tablet |
| `poswiadczenie_urzadzenia.stan` | aktywne, wygasle, odwolane |
| `kod_parowania.stan` | wydany, uzyty, uniewazniony, wygasly |
| `dzialanie_lokalne.rodzaj` | polecenie_czatu, zatwierdzenie, przerwanie, wstrzymanie, wznowienie, ustawienie, powiadomienie |
| `dzialanie_lokalne.przedmiot_typ` | sesja, zadanie, przebieg_petli, automatyka, ustawienie, powiadomienie |
| `dzialanie_lokalne.stan_uzgodnienia` | oczekuje, przyjete, bezprzedmiotowe, odrzucone |
| `srodowisko.kod` | TALKIN, WORKSPACE, CODESTUDIO, MULTITASKINGAI |
| `modul.kod` | STUDIO, WORKSPACE, AUTOMATIONS, BROWSER, RESEARCH, LIBRARY, TRANSLATE, ROUNDTABLE, DESIGN, ASSISTANT, TERMINAL, DEVELOPER, DIAGNOSTICS, APPS, AGENTS |
| `strefa_glowna.forma_prezentacji` | karta_srodowiska, kafel_komponentu, listwa_ustawien |
| `strefa_glowna.waga_wizualna` | glowna, posrednia, najnizsza |
| `pozycja_strefy.typ_odniesienia` | srodowisko, rodzaj_komponentu_wlasnego, funkcja_ustawien |
| `komponent_wlasny.rodzaj` | automatyka, agent, projekt, profil_asystenta |
| `komponent_wlasny.miejsce_konfiguracji` | strefa_2_strony_glownej, okno_modulu |
| `komponent_wlasny.widocznosc` | globalny, projektowy |
| `powiazanie_komponentu.zrodlo_typ` / `.cel_typ` | modul, komponent_wlasny, srodowisko |
| `projekt.tryb_dostepu_historii` | pelny, ograniczony, brak |
| `karta_sesji.stan` | aktywna, w_tle, zamknieta |
| `sesja.stan` | aktywna, wstrzymana, zakonczona |
| `wiadomosc.rola` | uzytkownik, wykonawca |
| `proces_sesji.stan` | uruchomiony, wstrzymany, zakonczony |
| `sesja_okno` — brak enumeracji (pole `stan` jest JSON) | — |
| `zasob_pamieci.poziom` | globalna, srodowisko, modul, projekt, sesja |
| `zadanie.typ` | prompt, akcja, orkiestracja, petla |
| `zadanie.status` | oczekujace, w_trakcie, wstrzymane, zakonczone, blad |
| `komponent_kompozycji.rodzaj` | prompt, zadanie, akcja, orkiestracja, subagenci, petla_koordynator, petla_wykonawcza, sekwencja, kompilacja |
| `kolejka.zasieg` | globalna, lokalna, modelu, agenta, projektu |
| `kolejka.obsluga_bledow` | retry, route, pause |
| `pozycja_kolejki.stan` | oczekuje, przetwarzana, zakonczona, wstrzymana |
| `log_akcji_kolejki.rodzaj_akcji` | enqueue, dequeue, delay, retry, pause, resume, split, merge, route, branch, condition |
| `harmonogram.cyklicznosc` | jednorazowo, cyklicznie |
| `harmonogram.stan` | aktywny, wstrzymany, zakonczony |
| `zaleznosc_orkiestracji.zrodlo_typ` / `.cel_typ` | kanal_modelu, agent, zadanie, kolejka, automatyka, projekt |
| `zaleznosc_orkiestracji.typ_zaleznosci` | poprzedza, blokuje, wyzwala, warunkuje |
| `kanal_modelu.typ` | API, CLI, SSH, HTTP |
| `kanal_modelu.przypisanie_typ` | sesja, rola |
| `agent_uprawnienie.poziom` | brak, odczyt, zapis, pelny |
| `rola_multitasking.rodzaj_roli` | executor_1, executor_2, coordinator, executor_3_validator |
| `rola_multitasking.wcielenie` | validator, reviewer, security_auditor, architect, product_owner, qa_lead, arbitrator, (tekst dowolny dla innych wcieleń) |
| `rola_multitasking.relacja_z_innymi` | niezalezna, przekazywanie, naprzemienna, iteracyjna |
| `subagent.status` | aktywny, zakonczony, blad |
| `sekcja_panelu_orkiestracji.kod` | zespoly, role, kolejki, orkiestracja, harmonogram_automatyki, monitor_procesu |
| `przebieg_multitasking.status` | w_trakcie, zakonczony, blad, wstrzymany |
| `rozszerzenie.rodzaj` | wtyczka, umiejetnosc, konektor, serwer_mcp |
| `rozszerzenie.zrodlo` | danaco_plugin, personal |
| `poziom_zasiegu.kod` | globalny, srodowisko, modul, para_modulow, projekt, karta_sesji, rola, okno (wartość docelowa, jeszcze nieobecna w kodzie — rozdz. 15.1) |
| `profil_izolacji.warstwa` | domyslna, sesji |
| `regula_izolacji_kontekstu.wymiar` | historia, pamiec, kontekst |
| `regula_izolacji_kontekstu.wartosc` | wspoldzielona, odrebna |
| `regula_izolacji_technicznej.zakres` | katalog_roboczy_sesji, srodowisko_procesu, katalog_danych_modelu, dostep_sieciowy, odczyt_zapis_plikow, konto_i_token, model_procesu, serwer_wykonania |
| `regula_izolacji_technicznej.stan` | wlaczony, wylaczony |
| `przypisanie_profilu_izolacji.typ_celu` | globalny, srodowisko, modul, para_modulow, projekt, karta_sesji, rola |
| `wersja_artefaktu.autor` | uzytkownik, wykonawca |
| `sugestia_aod.kontekst_typ` | srodowisko, modul, projekt, sesja |
| `sugestia_aod.rodzaj` | doradztwo, konfiguracja, problem, kolejny_krok |
| `sugestia_aod.status` | nowa, przyjeta, odrzucona |
| `polecenie_glosowe.status` | w_trakcie, zakonczone, blad |
| `powiadomienie.klasa` | zakonczenie, decyzja, blad, wzmianka, termin, automatyka, system |
| `powiadomienie.waga` | informacyjna, normalna, wymagajaca_decyzji |
| `powiadomienie.zrodlo_typ` | srodowisko, modul, projekt, sesja, zadanie, przebieg_petli, automatyka |
| `powiadomienie.stan` | nowe, odczytane, obsluzone, odlozone |
| `powiadomienie.kanal_dostarczenia` | centrum, centrum_i_push |
| `rozmowa.stan` | aktywna, wstrzymana, zarchiwizowana |
| `zlecenie.status` | przyjete, w_dekompozycji, w_realizacji, wstrzymane, zakonczone, przerwane, blad |
| `przebieg_petli.status` | w_trakcie, wstrzymany, zakonczony, przerwany, blad |
| `komunikat_sterujacy.nadawca` | koordynator, wykonawca, uzytkownik |
| `komunikat_sterujacy.rodzaj` | przydzial_zadania, raport_postepu, zapytanie_o_kontekst, odpowiedz_kontekstowa, zatwierdzenie, korekta_zlecenia, wstrzymanie, wznowienie, przerwanie, zakonczenie |
| `wynik_kontroli_jakosci.ocena` | przyjete, do_poprawy, odrzucone |
| `ponowienie.przyczyna` | negatywna_kontrola_jakosci, blad_wykonania, przekroczony_czas, korekta_zlecenia |
| `warstwa_widocznosci.kod` | zawsze_widoczna, na_zadanie, rozwiniecia_kontekstowe, funkcje_eksperckie |
| `element_interfejsu.rodzaj` | okno, panel, pasek_narzedzi, menu, znacznik_kontekstu, funkcja |
| `element_interfejsu.sposob_dostepu` | bez_interakcji, ikona, przycisk, przelacznik, znacznik_kontekstowy, menu_kebab, menu_hamburger, menu_kontekstowe, panel_popover, lista_rozwijana, polecenie_jezyka_naturalnego, skrot_klawiszowy, wyszukiwarka_funkcji, tryb_administracyjny, konfiguracja_roli |
| `rola_uzytkownika.kod` | podstawowy, zaawansowany, administrator |

---

## Załącznik B. Pełny wykaz tabel

| Tabela | Rozdział | Przeznaczenie w jednym zdaniu |
|---|---|---|
| `konto` | 3.1 | Jedyne konto właściciela platformy. |
| `urzadzenie` | 3.2 | Urządzenie połączone z platformą. |
| `polaczenie` | 3.3 | Pojedyncze połączenie WebSocket urządzenia z serwerem. |
| `poswiadczenie_urzadzenia` | 3.4 | Poświadczenie wydane urządzeniu w protokole parowania. |
| `kod_parowania` | 3.5 | Kod jednorazowy wydany przy inicjowaniu parowania. |
| `dzialanie_lokalne` | 3.6 | Pozycja kolejki działań wykonanych bez połączenia. |
| `srodowisko` | 4.1 | Cztery środowiska platformy. |
| `modul` | 4.2 | Piętnaście modułów platformy. |
| `srodowisko_modul` | 4.3 | Macierz dostępności modułów. |
| `okno_operacyjne` | 4.4 | Katalog okien roboczych per moduł. |
| `sesja_okno` | 4.5 | Migawka układu okien danej sesji. |
| `strefa_glowna` | 5.1 | Trzy strefy strony głównej. |
| `pozycja_strefy` | 5.2 | Pozycja (karta/kafel/element listwy) w strefie strony głównej. |
| `komponent_wlasny` | 6.1 | Encja nadrzędna czterech rodzajów komponentów własnych. |
| `automatyka` | 6.2 | Szczegóły komponentu własnego rodzaju automatyka. |
| `profil_asystenta` | 6.3 | Szczegóły komponentu własnego rodzaju profil asystenta. |
| `powiazanie_komponentu` | 6.4 | Jawne powiązania między modułami i komponentami własnymi. |
| `projekt_agent` | 6.5 | Przypisanie agenta do projektu. |
| `projekt` | 7.1 | Szczegóły komponentu własnego rodzaju projekt. |
| `projekt_srodowisko` | 7.2 | Widoczność projektu w środowiskach. |
| `karta_sesji` | 8.1 | Karta pozioma w oknie środowiska. |
| `sesja` | 8.2 | Instancja przestrzeni roboczej modułu w karcie. |
| `proces_sesji` | 8.3 | Rejestr procesu serwera obsługującego sesję. |
| `wiadomosc` | 8.4 | Wpis w historii rozmowy sesji. |
| `zalacznik_wiadomosci` | 8.5 | Załącznik do wiadomości. |
| `zasob_pamieci` | 9.1 | Wpis pamięci kontekstowej na jednym z pięciu poziomów. |
| `konfiguracja_pamieci_sesji` | 9.2 | Ustawienia dostępu sesji do pamięci. |
| `zadanie` | 10.1 | Jednostka pracy do wykonania. |
| `komponent_kompozycji` | 10.2 | Element kompozycji układów pracy Operatora. |
| `kolejka` | 10.3 | Struktura porządkująca wykonanie zadań. |
| `pozycja_kolejki` | 10.4 | Zadanie w obrębie kolejki. |
| `log_akcji_kolejki` | 10.5 | Rejestr akcji silnika kolejek. |
| `harmonogram` | 10.6 | Definicja wykonania zadań w czasie. |
| `zaleznosc_orkiestracji` | 10.7 | Zależność między elementami pracy. |
| `kanal_modelu` | 11.1 | Sposób połączenia z modelem AI. |
| `agent` | 12.1 | Szczegóły komponentu własnego rodzaju agent. |
| `agent_umiejetnosc` | 12.2 | Umiejętność przypisana agentowi. |
| `agent_rozszerzenie` | 12.3 | Rozszerzenie przypisane agentowi. |
| `agent_uprawnienie` | 12.4 | Uprawnienie agenta w danym zakresie. |
| `zespol_multitasking` | 13.1 | Zapisany preset zespołu środowiska MultitaskingAI. |
| `rola_multitasking` | 13.2 | Przypisanie roli w środowisku wielomodelowym. |
| `subagent` | 13.3 | Podagent uruchomiony przez wykonawcę. |
| `sekcja_panelu_orkiestracji` | 13.4 | Sekcja bocznej nawigacji środowiska MultitaskingAI. |
| `przebieg_multitasking` | 13.5 | Pojedynczy przebieg pętli autonomicznej. |
| `decyzja_procesu` | 13.6 | Decyzja w hierarchii decyzji przebiegu. |
| `rozszerzenie` | 14.1 | Wtyczka, umiejętność, konektor lub serwer MCP. |
| `poziom_zasiegu` | 15.1 | Osiem poziomów zasięgu reguł izolacji i konfiguracji. |
| `profil_izolacji` | 15.2 | Nazwany zestaw ustawień izolacji. |
| `regula_izolacji_kontekstu` | 15.3 | Reguła współdzielenia historii, pamięci lub kontekstu. |
| `regula_izolacji_technicznej` | 15.4 | Reguła jednego z ośmiu zakresów izolacji technicznej. |
| `przypisanie_profilu_izolacji` | 15.5 | Przypisanie profilu do elementu na danym poziomie zasięgu. |
| `para_modulow` | 15.6 | Para modułów jako cel poziomu zasięgu „para modułów”. |
| `artefakt` | 16.1 | Plik utworzony w trakcie pracy. |
| `wersja_artefaktu` | 16.2 | Historia wersji artefaktu. |
| `kolekcja` | 16.3 | Grupowanie artefaktów. |
| `kolekcja_artefakt` | 16.4 | Przypisanie artefaktu do kolekcji. |
| `etykieta_artefaktu` | 16.5 | Etykieta artefaktu. |
| `ustawienie` | 17.1 | Pojedynczy parametr zachowania platformy. |
| `sugestia_aod` | 18.2 | Proaktywna sugestia Always On Display. |
| `polecenie_glosowe` | 18.3 | Log polecenia głosowego (Assistant, Always On Display). |
| `powiadomienie` | 18.4 | Zdarzenie centrum powiadomień, z kanałem dostarczenia na urządzenie mobilne. |
| `rozmowa` | 19.1 | Wątek Chat Window w kanale Użytkownik ↔ Wykonawca. |
| `zlecenie` | 19.3 | Jednostka pracy przekazana Koordynatorowi do dekompozycji. |
| `przebieg_petli` | 19.5 | Pojedyncza iteracja pętli wykonawczej Koordynator ↔ Wykonawca. |
| `komunikat_sterujacy` | 19.6 | Wpis wymiany i sterowania w oknie pętli wykonawczej. |
| `wynik_kontroli_jakosci` | 19.7 | Ocena rezultatu zadania wobec kryteriów akceptacji. |
| `ponowienie` | 19.8 | Powtórne skierowanie zadania do realizacji. |
| `warstwa_widocznosci` | 20.1 | Cztery warstwy widoczności interfejsu. |
| `element_interfejsu` | 20.2 | Element interfejsu wraz z warstwą i sposobem wywołania. |
| `rola_uzytkownika` | 20.3 | Rola użytkownika sterująca zestawem widocznych elementów. |
| `konfiguracja_warstwy_roli` | 20.4 | Nadpisanie warstwy elementu dla roli użytkownika. |
| `funkcja_platformy` | 20.5 | Rejestr funkcji dla wyszukiwarki funkcji. |

Model liczy siedemdziesiąt dwie tabele.

---

## Załącznik C. Komponenty własne a moduły platformowe

Zestawienie porządkuje rozróżnienie między modułem (elementem platformy) a komponentem własnym (nazwanym wytworem Operatora), wprowadzone w rozdziale 1.2 niniejszego dokumentu.

| Cecha | Moduł (`modul`) | Komponent własny (`komponent_wlasny`) |
|---|---|---|
| Liczba wierszy | Stała, piętnaście (tabela słownikowa) | Rośnie w trakcie użytkowania platformy |
| Kto tworzy | Element platformy, dostępny dla każdego użytkownika | Operator, w oknie konfiguracji właściwym rodzajowi |
| Miejsce powstania | — (istnieje od instalacji platformy) | Strefa 2 strony głównej lub okno modułu |
| Odniesienie w tym modelu | `srodowisko_modul`, `okno_operacyjne` | `automatyka`, `agent`, `projekt`, `profil_asystenta` (1:1 z `komponent_wlasny`) |
| Przykład | Moduł Studio (`STUDIO`) | Konkretna automatyka „Raport tygodniowy” |

Moduły Workspace i Agents są jednocześnie elementami warstwy modułów (obecne w macierzy dostępności modułów, `srodowisko_modul`) oraz punktem wejścia do tworzenia komponentów własnych odpowiednio rodzaju `projekt` i `agent` — te dwie role są rozdzielone w modelu na dwie odrębne tabele (`modul` i `projekt`/`agent`) połączone jedynie pośrednio, przez fakt, że okno modułu Workspace jest jednym z dwóch miejsc konfiguracji (`komponent_wlasny.miejsce_konfiguracji = 'okno_modulu'`) obok strefy 2 strony głównej.

---

## Załącznik D. Diagram powiązań encji

Poniższy schemat tekstowy przedstawia główne powiązania klucza obcego między encjami modelu, w układzie zgodnym z mapą encji z rozdziału 2.

```
konto ──< urzadzenie ──< polaczenie ──> karta_sesji
konto ──< kod_parowania ──> urzadzenie
urzadzenie ──< poswiadczenie_urzadzenia
urzadzenie ──< dzialanie_lokalne
konto ──< karta_sesji ──< sesja ──> proces_sesji (1:1)
                            │
                            ├──< rozmowa ──< wiadomosc ──< zalacznik_wiadomosci
                            ├──< sesja_okno ──> okno_operacyjne
                            ├──< zadanie
                            └──> projekt / modul / srodowisko

srodowisko ──< srodowisko_modul >── modul ──< okno_operacyjne

strefa_glowna ──< pozycja_strefy ──(odniesienie polimorficzne)──> srodowisko
                                                                 → rodzaj komponentu własnego
                                                                 → funkcja ustawień

komponent_wlasny ──> automatyka | agent | projekt | profil_asystenta   (1:1, wg `rodzaj`)
komponent_wlasny ──< powiazanie_komponentu >── modul / komponent_wlasny / srodowisko

projekt ──< projekt_srodowisko >── srodowisko
projekt ──< projekt_agent >── agent
projekt ──< artefakt, kolekcja

zasob_pamieci ──(odniesienie polimorficzne)──> srodowisko / modul / projekt / sesja

rozmowa ──< zlecenie ──< przebieg_petli ──< komunikat_sterujacy
                    │                    ──< wynik_kontroli_jakosci ──< ponowienie
                    ├──< zadanie
                    └──> kolejka (zlecenie_id)
zadanie ──< wynik_kontroli_jakosci ──< ponowienie
przebieg_petli ──> przebieg_multitasking / rola_multitasking (Koordynator)

zadanie ──> kolejka ──< pozycja_kolejki ──> zadanie
zadanie ──> harmonogram
zadanie ──< komponent_kompozycji ──> modul
kolejka ──< log_akcji_kolejki
zaleznosc_orkiestracji ──(odniesienie polimorficzne)──> kanal_modelu / agent / zadanie / kolejka / automatyka / projekt

agent ──< agent_umiejetnosc, agent_rozszerzenie >── rozszerzenie, agent_uprawnienie
agent ──> kanal_modelu (model bazowy)

zespol_multitasking ──< rola_multitasking ──< subagent
                     ──< kolejka (zespol_id)
                     ──< zaleznosc_orkiestracji (zespol_id)
                     ──< przebieg_multitasking ──< decyzja_procesu
rola_multitasking ──> agent | kanal_modelu (wykonawca)
rola_multitasking ──> profil_izolacji

poziom_zasiegu ──< profil_izolacji ──< regula_izolacji_kontekstu
                                    ──< regula_izolacji_technicznej
                                    ──< przypisanie_profilu_izolacji ──(polimorficzne)──> srodowisko / modul
                                                                                        / para_modulow / projekt
                                                                                        / karta_sesji / rola_multitasking
poziom_zasiegu ──< ustawienie ──(polimorficzne)──> (jak wyżej)

artefakt ──< wersja_artefaktu
artefakt ──< kolekcja_artefakt >── kolekcja
artefakt ──< etykieta_artefaktu

warstwa_widocznosci ──< element_interfejsu ──> modul / srodowisko / okno_operacyjne
                     ──< rola_uzytkownika ──< konfiguracja_warstwy_roli >── element_interfejsu
                     ──< funkcja_platformy ──> element_interfejsu
```

---

## Załącznik E. Szablony konfiguracji pól JSON

Szablony redakcyjne wspierające implementację kolumn JSON (rozdział 22.3). Szablony pól `sesja_okno.stan`, `automatyka.definicja_workflow`, `automatyka.wywolania_modeli`, `profil_asystenta.zakres_polecen`, `profil_asystenta.obslugiwane_akcje`, `log_akcji_kolejki.parametry`, `zaleznosc_orkiestracji.warunek`, `agent.konfiguracja_pamieci`, `agent.konfiguracja_uprawnien`, `przebieg_petli.wskazniki`, `komunikat_sterujacy.ladunek`, `wynik_kontroli_jakosci.kryteria` i `funkcja_platformy.slowa_kluczowe` podano przy właściwych encjach w rozdziałach 4–20. Poniżej trzy szablony uzupełniające.

Pole `sesja.kontekst` (bieżący kontekst przekazywany do AI):

```json
{
  "podsumowanie_biezace": "Praca nad modułem rozliczeń.",
  "wskazniki": { "otwarte_pliki": 3, "aktywny_agent_id": 7 },
  "zrodla": ["Library/umowa.pdf", "Research/analiza-rynku.md"]
}
```

Pole `proces_sesji.zmienne_srodowiskowe` (środowisko procesu sesji):

```json
{
  "KATALOG_ROBOCZY": "dane/sesje/1024",
  "MODEL_KANAL": "api-podstawowy",
  "TRYB": "produkcyjny"
}
```

Pole `komponent_kompozycji.definicja` (definicja skomponowanego układu pracy — pełna pętla z Koordynatorem):

```json
{
  "rodzaj": "petla_koordynator",
  "role": ["coordinator", "executor_1", "executor_2", "executor_3_validator"],
  "kolejka_id": 5,
  "warunek_zakonczenia": { "pole": "status", "wartosc": "zakonczone" }
}
```

---

## Załącznik F. Scenariusze danych

Scenariusze ilustrują, jakie wiersze powstają lub zmieniają się w modelu podczas kompozycji funkcji platformy. Nie wprowadzają nowych encji — pokazują współdziałanie tabel opisanych w rozdziałach 3–20.

### F.1. Autonomiczna budowa aplikacji w środowisku MultitaskingAI

Scenariusz obejmuje budowę produktu w module Apps z wykorzystaniem środowiska MultitaskingAI.

| Krok | Wiersze utworzone lub zmienione |
|---|---|
| 1. Zapis presetu zespołu | `zespol_multitasking` |
| 2. Przypisanie ról | `rola_multitasking`: `executor_1` (warstwa serwerowa), `executor_2` (warstwa kliencka, `relacja_z_innymi = 'niezalezna'`), `coordinator`, `executor_3_validator` (`wcielenie = 'security_auditor'`, w kolejnym przebiegu `'qa_lead'`) |
| 3. Uruchomienie sieci podagentów | `subagent`: pod Executorem 1 — Agent Backend, Agent API, Agent Database; pod Executorem 2 — Agent UI |
| 4. Kolejkowanie i rozdział pracy | `kolejka`, `pozycja_kolejki`, `log_akcji_kolejki` (`split`, `merge`) |
| 5. Zależność warstwy serwerowej i klienckiej | `zaleznosc_orkiestracji` (`typ_zaleznosci = 'poprzedza'`) |
| 6. Spięcie z automatyką i harmonogramem | `automatyka` powiązana przez `powiazanie_komponentu` ze środowiskiem; `harmonogram` (`cyklicznosc = 'cyklicznie'`) |
| 7. Praca ciągła i decyzje | `przebieg_multitasking` (kolejne przebiegi), `decyzja_procesu` (hierarchia) |
| 8. Prowadzenie pętli wykonawczej | `rozmowa` i `wiadomosc` (polecenie Użytkownika), `zlecenie`, `przebieg_petli`, `komunikat_sterujacy` (`przydzial_zadania`, `raport_postepu`), `wynik_kontroli_jakosci`, `ponowienie` |
| 9. Zasoby wizualne z modułu Design | `powiazanie_komponentu` (`modul` Design → `modul` Apps) |
| 10. Nadzór i zatwierdzenia | `sugestia_aod` (Always On Display jako operator), `polecenie_glosowe` (zatwierdzenie z urządzenia przez funkcję Mobile) |

### F.2. Przepływ badawczy między modułami

Scenariusz obejmuje przepływ Browser → Research → Studio → Library.

| Krok | Wiersze utworzone lub zmienione |
|---|---|
| 1. Wspólna analiza stron | `sesja` w module Browser; źródła jako `artefakt` |
| 2. Zasilenie modułu Research | `powiazanie_komponentu` (`modul` Browser → `modul` Research) |
| 3. Analiza i raport | `sesja` w module Research; `artefakt` (raport), `wersja_artefaktu` |
| 4. Redakcja w module Studio | `powiazanie_komponentu` (`modul` Research → `modul` Studio); nowe `wersja_artefaktu` |
| 5. Trwałe przechowanie | `artefakt`, `kolekcja`, `kolekcja_artefakt`, `etykieta_artefaktu` (moduł Library) |

### F.3. Konfiguracja punktu izolacji między projektami

Scenariusz obejmuje dwa projekty rozdzielone, ze współdzieloną biblioteką.

| Krok | Wiersze utworzone lub zmienione |
|---|---|
| 1. Utworzenie profilu na poziomie projektu | `profil_izolacji` (`poziom_zasiegu` = projekt, `warstwa = 'domyslna'`) |
| 2. Ustawienie izolacji kontekstu | `regula_izolacji_kontekstu`: historia = odrebna, pamiec = odrebna, kontekst = odrebna |
| 3. Repozytorium plików współdzielone | `regula_izolacji_technicznej`: brak wiersza `odczyt_zapis_plikow` (zakres niewłączony) |
| 4. Przypisanie obu projektom | `przypisanie_profilu_izolacji`: dwa wiersze `typ_celu = 'projekt'` |

Reguła projektowa ma pierwszeństwo przed regułą globalną, zgodnie z algorytmem rozstrzygania z rozdziału 15.5.

### F.4. Praca dwujęzyczna Studio i Translate

Scenariusz obejmuje współdzielenie operacji kontekstowych AI między Studio a Translate.

| Krok | Wiersze utworzone lub zmienione |
|---|---|
| 1. Zdefiniowanie pary modułów | `para_modulow` (Studio, Translate) |
| 2. Profil na poziomie pary modułów | `profil_izolacji` (`poziom_zasiegu` = para_modulow) |
| 3. Współdzielenie kontekstu | `regula_izolacji_kontekstu`: kontekst = wspoldzielona |
| 4. Jawne powiązanie modułów | `powiazanie_komponentu` (`modul` Studio → `modul` Translate) — jawne i odwracalne |
| 5. Przypisanie pary | `przypisanie_profilu_izolacji`: `typ_celu = 'para_modulow'` |

---

*Koniec dokumentu. Model danych — Specyfikacja docelowa, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](../LICENSE.md). Kontakt: support@danaco-group.pl*
