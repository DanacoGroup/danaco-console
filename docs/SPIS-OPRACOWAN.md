# Danaco Console — Spis opracowań

| | |
|---|---|
| **Produkt** | Danaco Console |
| **Rodzaj** | Platforma AI Workspace OS |
| **Producent** | Danaco Holding Group Sp. z o.o. |
| **Twórca** | Dariusz Naharnowicz |
| **Wersja** | v2.0 |
| **Status** | Deweloperski |
| **Data** | 2026-08-20 |

**Informacje szczegółowe dokumentu:**

| | |
|---|---|
| **Tytuł** | Spis opracowań zbioru dokumentacji |
| **Klasa dokumentu** | Nawigacja |
| **Odbiorcy** | deweloper · projektant · redaktor dokumentacji |
| **Przeznaczenie** | Jest katalogiem całego zbioru: wymienia wszystkie opracowania z ich klasą, odbiorcą, przeznaczeniem i objętością oraz podaje ścieżki czytania dla trzech ról |
| **Zakres** | pełny wykaz 51 opracowań pogrupowanych po katalogach, ścieżki czytania, statystyka klas dokumentów, historia redakcji |
| **Poza zakresem** | treść merytoryczna produktu — noszą ją opracowania wymienione niżej; wzór redakcyjny — [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) |
| **Dokument nadrzędny** | [Opis produktu (README)](README.md) |
| **Dokumenty powiązane** | [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) · [Licencja produktu](LICENSE.md) |
| **Prototypy odniesienia** | nie dotyczy — spis porządkuje zbiór dokumentacji, nie opisuje żadnego okna produktu; prototypy przywołują opracowania wymienione w rozdziale 1 |
| **Źródła normatywne** | struktura katalogu `docs/`, polecenie `find docs -name '*.md' \| sort` |
| **Zasada nadrzędna** | Katalog wymienia komplet opracowań zbioru; opracowanie spoza tego spisu nie należy do zbioru |

Zbiór dokumentacji Danaco Console liczy pięćdziesiąt jeden opracowań rozłożonych w korzeniu
katalogu `docs/` oraz w sześciu podkatalogach tematycznych. Niniejszy spis jest jedynym
miejscem, w którym cały zbiór widoczny jest naraz — każde opracowanie wymienione jest
dokładnie raz, z tytułem, klasą, jednozdaniowym przeznaczeniem i objętością liczoną
w znakach pliku źródłowego. Kolejność wewnątrz każdej grupy odpowiada kolejności czytania
właściwej danej grupie, nie kolejności alfabetycznej nazw plików. Dziewięć rozdziałów
niniejszego dokumentu odpowiada dziewięciu odrębnym pytaniom nawigacyjnym: co należy do
zbioru, jak opracowania są ze sobą powiązane, w jakiej kolejności je czytać, do jakiej klasy
każde należy, jak duży jest cały zbiór, jak spis utrzymywać w zgodzie z katalogiem, jakich
pojęć używa, co zmieniło się w kolejnych falach redakcji oraz na jakich warunkach spis jest
przyjmowany.

---

## Spis treści

1. [Katalog opracowań](#1-katalog-opracowań)
   - [1.1 Korzeń zbioru](#11-korzeń-zbioru)
   - [1.2 architektura/](#12-architektura)
   - [1.3 specyfikacje/](#13-specyfikacje)
   - [1.4 interfejs-uzytkownika/](#14-interfejs-uzytkownika)
   - [1.5 moduly/](#15-moduly)
   - [1.6 srodowiska/](#16-srodowiska)
   - [1.7 funkcje-globalne/](#17-funkcje-globalne)
   - [1.8 Metoda weryfikacji krzyżowej zbioru](#18-metoda-weryfikacji-krzyżowej-zbioru)
2. [Mapa hierarchii dokumentu nadrzędnego](#2-mapa-hierarchii-dokumentu-nadrzędnego)
3. [Ścieżki czytania](#3-ścieżki-czytania)
   - [3.1 Deweloper — od koncepcji do budowy](#31-deweloper--od-koncepcji-do-budowy)
   - [3.2 Projektant — od interfejsu do zachowań](#32-projektant--od-interfejsu-do-zachowań)
   - [3.3 Redaktor dokumentacji — przed każdą edycją](#33-redaktor-dokumentacji--przed-każdą-edycją)
4. [Klasy dokumentów](#4-klasy-dokumentów)
5. [Statystyka objętości zbioru](#5-statystyka-objętości-zbioru)
6. [Zasady utrzymania spisu](#6-zasady-utrzymania-spisu)
   - [6.1 Dodawanie nowego opracowania](#61-dodawanie-nowego-opracowania)
   - [6.2 Aktualizacja objętości](#62-aktualizacja-objętości)
   - [6.3 Weryfikacja automatyczna](#63-weryfikacja-automatyczna)
7. [Słowniczek pojęć nawigacyjnych](#7-słowniczek-pojęć-nawigacyjnych)
8. [Historia redakcji](#8-historia-redakcji)
9. [Kryteria odbioru spisu](#9-kryteria-odbioru-spisu)

---

## 1. Katalog opracowań

Zbiór liczy **51 opracowań** w sześciu katalogach tematycznych oraz korzeniu. Liczbę
potwierdza polecenie `find docs -name '*.md' | sort`, uruchamiane ponownie po każdej fali
redakcji — rozbieżność między jego wynikiem a liczbą wierszy tego rozdziału jest usterką
niniejszego spisu, nie zbioru. Każdy wiersz odsyła do pliku klikalnie; kolumna „Klasa”
rozstrzyga, czy dokument opisuje stan docelowy (Specyfikacja docelowa), stan faktyczny kodu
(Stan wdrożenia), warunki korzystania (Akt prawny) czy porządek zbioru (Nawigacja) — cztery
klasy zdefiniowane w rozdziale 2.2 Standardu redakcyjnego i językowego.

### 1.1 Korzeń zbioru

Sześć opracowań korzenia dzieli się na dwie grupy o odmiennej naturze. Cztery opracowania
Właściciela — README, Instrukcja użytkowania, Instalacja i konfiguracja, Licencja — niosą
stan faktyczny kodu bieżącego wydania i warunki prawne korzystania z produktu; ich objętość
nie może ulec zmniejszeniu przy żadnej redakcji (rozdz. 9 Standardu redakcyjnego
i językowego). Dwa opracowania własne zbioru — niniejszy spis i Standard redakcyjny
i językowy — porządkują sam zbiór dokumentacji, nie produkt. Znajomość tej granicy
rozstrzyga, czy poprawka w opracowaniu korzenia jest zmianą treści merytorycznej produktu,
czy zmianą porządku samego zbioru — rozróżnienie istotne przy każdej redakcji.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Platforma AI Workspace OS](README.md) | Stan wdrożenia | Główna karta techniczno-produktowa Danaco Console — stan faktyczny wdrożenia platformy w bieżącym wydaniu: architektura zrealizowana, moduły dostępne, stos technologiczny uruchomiony. | 134 555 |
| [Spis opracowań zbioru dokumentacji](SPIS-OPRACOWAN.md) | Nawigacja | Niniejszy dokument — katalog całego zbioru: wymienia wszystkie opracowania z ich klasą, odbiorcą, przeznaczeniem i objętością oraz podaje ścieżki czytania dla trzech ról. | 45 887 |
| [Standard redakcyjny i językowy zbioru dokumentacji oraz konwencje dalszej budowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) | Nawigacja | Ustala jeden obowiązujący kształt każdego opracowania zbioru, jeden słownik terminologiczny oraz konwencje nazewnicze wiążące w dalszej budowie produktu. | 45 879 |
| [Instalacja i konfiguracja](INSTALACJA-I-KONFIGURACJA.md) | Stan wdrożenia | Kompletny fundament techniczny wdrożenia w bieżącej wersji: kontrakt, rdzeń, warstwa danych, transport i powłoka natywna — krok po kroku od pobrania pakietu do pierwszego uruchomienia. | 187 662 |
| [Instrukcja użytkowania](INSTRUKCJA-UZYTKOWANIA.md) | Stan wdrożenia | Przewodnik po codziennej pracy z aplikacją w stanie faktycznym bieżącego wydania — uruchamianie, środowiska, moduły, komponenty własne i rozwiązywanie typowych trudności. | 160 251 |
| [Licencja produktu Danaco Console](LICENSE.md) | Akt prawny | Warunki, na jakich Danaco Holding Group Sp. z o.o. udostępnia produkt Danaco Console — zakres uprawnień, ograniczenia i odpowiedzialność. | 144 095 |

### 1.2 architektura/

Dziewięć opracowań fundamentu technicznego platformy. Koncepcja platformy stoi na szczycie
hierarchii dokumentu nadrzędnego całego zbioru — każde pozostałe opracowanie klasy
Specyfikacja docelowa wskazuje ją jako swój dokument nadrzędny, bezpośrednio albo przez
opracowanie pośredniczące. Katalog domyka osiem opracowań rozwijających poszczególne
mechanizmy przekrojowe: model wdrożenia, dane, kontrakty, konfigurację, bezpieczeństwo,
izolację, integrację modeli i rozszerzenia. Kolejność tabeli 1.2 odpowiada kolejności
czytania rozdziału 3.1 — od koncepcji, przez model wdrożenia, do mechanizmów szczegółowych.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Koncepcja platformy](architektura/koncepcja-platformy.md) | Specyfikacja docelowa | Ustala koncepcyjny, docelowy kształt platformy — charakter platformy, pojęcia podstawowe, warstwy architektury i widoczności, strukturę środowisk i modułów oraz zasady nadrzędne funkcjonalności. | 115 097 |
| [Architektura techniczna](architektura/architektura.md) | Specyfikacja docelowa | Ustala model wdrożenia, stos technologiczny, warstwy systemu i wszystkie mechanizmy przekrojowe: komunikację, sesje, obliczenia, dane, izolację, rozszerzenia, uwierzytelnianie i synchronizację. | 78 984 |
| [Model danych](architektura/model-danych.md) | Specyfikacja docelowa | Ustala pełny model danych platformy zgodny ze schematem przechowywania SQLite — grupy encji, ich pola, typy, relacje i ograniczenia — realizujący zasadę pełnej konfigurowalności. | 139 424 |
| [Kontrakty komunikacji](architektura/kontrakty-komunikacji.md) | Specyfikacja docelowa | Ustala pełne kontrakty komunikacji między klientem a rdzeniem platformy — kopertę protokołu, wykaz 68 obszarów kontraktu, schematy ładunku poleceń i zdarzeń oraz kody błędów. | 102 065 |
| [Model konfiguracji](architektura/model-konfiguracji.md) | Specyfikacja docelowa | Ustala pełną strukturę okna konfiguracji oraz warstwowość ustawień — globalna, środowisko, projekt, sesja — obowiązującą we wszystkich zakresach konfigurowalnych platformy. | 88 738 |
| [Bezpieczeństwo i uwierzytelnianie](architektura/bezpieczenstwo-i-uwierzytelnianie.md) | Specyfikacja docelowa | Ustala pełny model bezpieczeństwa platformy: uwierzytelnianie jako jedyny mechanizm kontroli dostępu, model jednego właściciela i wielu urządzeń oraz autoryzację w obu kanałach komunikacji. | 88 305 |
| [Izolacja i konfigurowalność zależności](architektura/izolacja-i-zaleznosci.md) | Specyfikacja docelowa | Ustala pełną konfigurowalność izolacji technicznej i izolacji kontekstu — cztery rodzaje izolacji, osiem zakresów technicznych, poziomy zasięgu i zasadę zera blokad twardych. | 75 894 |
| [Integracja modeli](architektura/integracja-modeli.md) | Specyfikacja docelowa | Ustala warstwę dostawcy modelu: cztery kanały integracji, mechanizm adaptera, encję kanału modelu, tożsamość i personę modelu oraz przypisanie modelu per sesja i per rola. | 71 902 |
| [Rozszerzenia](architektura/rozszerzenia.md) | Specyfikacja docelowa | Ustala jednolity kontrakt rozszerzenia obowiązujący niezależnie od jego rodzaju i pochodzenia — wtyczek, umiejętności, konektorów i serwerów MCP — oraz cykl życia rozszerzenia. | 64 944 |

### 1.3 specyfikacje/

Trzy opracowania łączące ustalenia architektury z konkretnym zakresem budowy: agentami jako
komponentem platformowym, piętnastoma modułami platformy oraz katalogiem wszystkich typów
okien operacyjnych. Każde z trzech schodzi o jeden poziom szczegółowości niżej niż Koncepcja
platformy i wyżej niż pojedyncze opracowanie modułu z katalogu `moduly/` — stanowi pomost
między zamysłem a wykonaniem. Trzy opracowania katalogu są jedynym miejscem zbioru, w którym
potwierdzana jest zgodność nazw modułów, agentów i okien z rzeczywistymi obszarami
`contract.json`.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Specyfikacja agentów](specyfikacje/specyfikacja-agentow.md) | Specyfikacja docelowa | Ustala model agentów platformy: definicję agenta, siedem komponentów jego definicji, Centrum uprawnień, dwa kanały komunikacji operacyjnej i pełny przebieg pętli wykonawczej. | 88 319 |
| [Specyfikacja modułów](specyfikacje/specyfikacja-modulow.md) | Specyfikacja docelowa | Ustala operacyjną specyfikację piętnastu modułów platformy — cel, okna operacyjne z warstwami widoczności, funkcjonalności, przebieg pracy oraz powiązania konfigurowalne. | 89 675 |
| [Specyfikacja okien operacyjnych](specyfikacje/specyfikacja-okien-operacyjnych.md) | Specyfikacja docelowa | Ustala katalog wszystkich typów okien operacyjnych piętnastu modułów oraz czterech ról środowiska MultitaskingAI — cel, zawartość, zachowanie, warstwę widoczności każdego okna. | 95 439 |

### 1.4 interfejs-uzytkownika/

Dwanaście opracowań niosących warstwę wizualną i nawigacyjną platformy, niezależną od
podziału na moduły: ramę okna obecną w każdym widoku, system tokenów i komponentów, stronę
główną, przepływ nawigacji między oknami oraz sześć odrębnych okien nakładkowych
(Konfiguracja, Ustawienia, historia sesji, nowy projekt, instalator, instrukcja użytkowania)
niemieszczących się w opracowaniu żadnego pojedynczego modułu. Sześć z dwunastu opracowań
tego katalogu opisuje pojedyncze okno nakładkowe wywoływane spoza przepływu głównego,
pozostałe sześć — mechanizmy i elementy wspólne wielu oknom naraz.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Strona główna i nawigacja](interfejs-uzytkownika/strona-glowna-i-nawigacja.md) | Specyfikacja docelowa | Szczegółowa specyfikacja strony głównej platformy (centrum dowodzenia) oraz nawigacji: układu kolumnowego, trzech stref wyboru, warstw widoczności i przepływu strona główna → środowisko → moduł. | 78 809 |
| [Rama okna aplikacji — belka tytułowa, szyna nawigacji, pasek edycji, pasek stanu i specyfikacja menu rozwijanych](interfejs-uzytkownika/rama-okna.md) | Specyfikacja docelowa | Materiał źródłowy dla czterech pasów ramy obecnych w każdym oknie platformy: co zawierają, jak się zachowują, jakie niosą czynności i co kryje każde menu rozwijane. | 86 391 |
| [Elementy okien przepływu głównego](interfejs-uzytkownika/elementy-okien.md) | Specyfikacja docelowa | Pełny inwentarz elementów każdego okna przepływu głównego bez własnego odrębnego opracowania — okna startowego, strony głównej i powłok czterech środowisk — wraz ze stanami. | 122 674 |
| [Przepływ okien: od uruchomienia aplikacji do okna roboczego](interfejs-uzytkownika/przeplyw-okien.md) | Specyfikacja docelowa | Ciągły przebieg nawigacyjny od uruchomienia aplikacji do okna roboczego, złożony z opracowań źródłowych — bez wprowadzania nowych funkcji, modułów, okien ani mechanizmów. | 148 700 |
| [System wizualny](interfejs-uzytkownika/system-wizualny.md) | Specyfikacja docelowa | System projektowy platformy na bazie przewodnika marki: tokeny kolorów w motywie jasnym i ciemnym, typografia, cienie, siatka kolumnowa, warstwy widoczności i ikonografia. | 73 855 |
| [Katalog komponentów interfejsu](interfejs-uzytkownika/katalog-komponentow.md) | Specyfikacja docelowa | Jeden rejestr wszystkich komponentów interfejsu wielokrotnego użytku — elementów, które nie należą do jednego okna czy modułu, lecz powtarzają się w dziesiątkach miejsc platformy. | 184 751 |
| [Okno Konfiguracji](interfejs-uzytkownika/konfiguracja.md) | Specyfikacja docelowa | Specyfikacja projektowa okna Konfiguracji — panelu programującego nakładkę na model bazowy: pola sterujące wywołaniem modelu, konfiguracja kanałów komunikacji i warstw widoczności. | 149 119 |
| [Okno Ustawień (poziom aplikacji)](interfejs-uzytkownika/ustawienia.md) | Specyfikacja docelowa | Projekt okna Ustawień — okna poziomu aplikacji odrębnego od okna Konfiguracji — gromadzącego wszystko, co dotyczy Operatora: konto, uwierzytelnianie, wygląd, urządzenia, powiadomienia. | 124 861 |
| [Okno historii sesji](interfejs-uzytkownika/okno-historii-sesji.md) | Specyfikacja docelowa | Kompletna budowa, zawartość i zachowanie okna historii sesji — nakładki okna ustawień otwartej na sekcji rejestru sesji — tak, aby deweloper zbudował je bez rozstrzygania niczego samodzielnie. | 44 364 |
| [Okno nowego projektu](interfejs-uzytkownika/okno-nowego-projektu.md) | Specyfikacja docelowa | Kompletna budowa i przebieg okna zakładania projektu w przestrzeni Workspace środowiska, tak aby deweloper zbudował je bez rozstrzygania czegokolwiek samodzielnie. | 17 661 |
| [Okno instalatora](interfejs-uzytkownika/okno-instalatora.md) | Specyfikacja docelowa | Kompletna budowa i przebieg okna instalatora aplikacji — okna przedaplikacyjnego prowadzącego przez sześć kroków od warunków licencji do pierwszego uruchomienia. | 63 758 |
| [Okno instrukcji użytkowania — nakładka wywoływana z menu Pomoc](interfejs-uzytkownika/okno-instrukcji.md) | Specyfikacja docelowa | Pełny ciąg budowy okna nakładkowego „Instrukcja użytkowania”: układ, nawigacja po rozdziałach, dosłowna treść dziewięciu rozdziałów, makiety osadzone w treści i kryteria odbioru. | 76 283 |

### 1.5 moduly/

Piętnaście opracowań — po jednym na każdy moduł platformy wymieniony w rozdziale 4
Specyfikacji modułów, w tej samej kolejności. Każde opracowanie schodzi z poziomu
specyfikacji operacyjnej modułu do poziomu gotowego do budowy: komplet okien operacyjnych,
katalog funkcji i narzędzi, przebieg pracy krok po kroku oraz punkty sterowania osiągalne
z okna konfiguracji. Objętość pojedynczego opracowania modułu — powyżej 85 000 znaków —
odzwierciedla komplet wykazów normatywnych wymaganych rozdziałem 4 Standardu redakcyjnego
i językowego, nie długość dla samej długości.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Moduł Studio — pełnozakresowa dokumentacja projektowa](moduly/studio.md) | Specyfikacja docelowa | Zaawansowana praca z tekstem, dokumentami oraz treścią — komplet okien operacyjnych, katalog funkcji i narzędzi, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 100 599 |
| [Moduł WorkSpace — dokumentacja projektowa](moduly/workspace.md) | Specyfikacja docelowa | Izolowana przestrzeń projektowa gromadząca w jednym miejscu zadania, pliki, notatki, wiedzę, instrukcje, pamięć kontekstową i przypisanych wykonawców AI jednego przedsięwzięcia. | 86 598 |
| [Moduł Automations — dokument projektowy](moduly/automations.md) | Specyfikacja docelowa | Budowanie i wykonywanie procesów automatycznych — harmonogramów, kolejek i pętli pracy ciągłej — komplet okien operacyjnych oraz punkty sterowania z okna konfiguracji. | 86 124 |
| [Moduł Browser — pełnozakresowa dokumentacja projektowa](moduly/browser.md) | Specyfikacja docelowa | Współdzielone przeglądanie internetu wspierane przez AI — katalog funkcji i narzędzi, komplet okien operacyjnych, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 84 129 |
| [Moduł Research — pełnozakresowa dokumentacja projektowa](moduly/research.md) | Specyfikacja docelowa | Realizacja badań, analiz i opracowań — komplet okien operacyjnych, zarządzanie źródłami i ustaleniami, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 117 328 |
| [Moduł Library — pełnozakresowa dokumentacja projektowa](moduly/library.md) | Specyfikacja docelowa | Centralne repozytorium wiedzy i plików platformy — komplet okien operacyjnych, katalog funkcji, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 90 631 |
| [Moduł Translate — pełnozakresowa dokumentacja projektowa](moduly/translate.md) | Specyfikacja docelowa | Wielojęzyczne tłumaczenia w trybie wielozadaniowym — katalog funkcji i narzędzi, komplet okien operacyjnych, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 93 990 |
| [Moduł Roundtable — pełnozakresowa dokumentacja projektowa](moduly/roundtable.md) | Specyfikacja docelowa | Współpraca wielu modeli AI nad wspólnym problemem — konfrontacja perspektyw i wypracowanie wspólnego stanowiska — komplet okien operacyjnych i katalog funkcji. | 87 110 |
| [Moduł Design — dokument projektowy](moduly/design.md) | Specyfikacja docelowa | Projektowanie i tworzenie zasobów wizualnych — komplet okien operacyjnych, katalog funkcji i narzędzi, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 99 495 |
| [Moduł Assistant — pełnozakresowa dokumentacja projektowa](moduly/assistant.md) | Specyfikacja docelowa | Naturalna komunikacja głosowa z AI, realizowana operacyjnie w środowisku TalkIn — komplet okien operacyjnych oraz punkty sterowania z okna konfiguracji. | 86 531 |
| [Moduł Terminal — pełnozakresowa dokumentacja projektowa](moduly/terminal.md) | Specyfikacja docelowa | Praca z konsolami i środowiskami wykonawczymi — komplet okien operacyjnych, katalog funkcji i narzędzi, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 105 886 |
| [Moduł Developer — pełnozakresowa dokumentacja projektowa](moduly/developer.md) | Specyfikacja docelowa | Tworzenie i rozwój kodu — komplet okien operacyjnych, katalog funkcji i narzędzi, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 101 923 |
| [Moduł Diagnostics — pełnozakresowa dokumentacja projektowa](moduly/diagnostics.md) | Specyfikacja docelowa | Analiza oraz usuwanie problemów technicznych — katalog funkcji i narzędzi, komplet okien operacyjnych, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 99 926 |
| [Moduł Apps — dokument projektowy](moduly/apps.md) | Specyfikacja docelowa | Budowa kompletnych produktów cyfrowych — komplet okien operacyjnych modułu, katalog funkcji i narzędzi, przebieg pracy oraz punkty sterowania z okna konfiguracji. | 94 217 |
| [Moduł Agents](moduly/agents.md) | Specyfikacja docelowa | Ustala interfejs modułu Agents: komplet okien operacyjnych, katalog elementów każdego okna, diagram cyklu tworzenia eksperta, stany oraz punkty sterowania z okna konfiguracji. | 117 024 |

### 1.6 srodowiska/

Cztery opracowania — po jednym na każde z czterech środowisk platformy: TalkIn, WorkSpace,
CodeStudio i MultitaskingAI. Środowisko jest ogólnym trybem pracy nadrzędnym wobec modułu;
opracowania tego katalogu opisują powłokę środowiska, jego przedsionek oraz mechanikę
wejścia, pozostawiając zawartość poszczególnych modułów katalogowi `moduly/`. Środowiska
TalkIn i WorkSpace udostępniają boczną nawigację modułów, CodeStudio zawęża ją do modułów
właściwych pracy inżynierskiej, a MultitaskingAI zastępuje ją panelem orkiestracji
czterech ról.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Środowisko TalkIn — dokumentacja projektowa: przeznaczenie, moduły, komplet okien, specyfikacja okien, przepływy pracy, stany i powiązania, punkty sterowania, scenariusze użycia](srodowiska/talkin.md) | Specyfikacja docelowa | Pełnozakresowy materiał źródłowy do projektowania i budowy interfejsu środowiska TalkIn: co istnieje, gdzie leży, w jakiej formie i do czego służy — dla każdego okna i elementu. | 71 464 |
| [Środowisko WorkSpace — powłoka środowiska, komplet okien, warstwy widoczności, nawigacja modułów, karty sesji, układy okien, współdzielenie kontekstu, punkty sterowania](srodowiska/workspace.md) | Specyfikacja docelowa | Materiał wykonawczy do projektowania i budowy interfejsu środowiska WorkSpace: co powstaje, gdzie leży, w jakiej formie i do czego służy — dla każdego okna i elementu powłoki. | 79 531 |
| [Środowisko CodeStudio](srodowiska/codestudio.md) | Specyfikacja docelowa | Pełnozakresowy materiał źródłowy do projektowania i budowy interfejsu środowiska programistycznego CodeStudio: co istnieje, gdzie leży w układzie, w jakiej formie i do czego służy. | 101 170 |
| [Środowisko MultitaskingAI — warstwa centralna, kanały komunikacji, role, silnik kolejek, panel orkiestracji, powłoka środowiska, integracja z Automations, makiety i katalog interfejsu](srodowiska/multitaskingai.md) | Specyfikacja docelowa | Pełnozakresowy, gęsty wizualnie materiał źródłowy do projektowania i budowy interfejsu środowiska MultitaskingAI: co ma powstać, gdzie ma leżeć i do czego służy — dla każdej roli. | 181 985 |

### 1.7 funkcje-globalne/

Dwa opracowania funkcji działających ponad podziałem na środowiska i moduły — Always On
Display oraz Mobile. Obie funkcje są dostępne niezależnie od aktywnego środowiska i modułu,
zgodnie z rozdziałem 5.3 Koncepcji platformy, i z tego powodu nie znalazły miejsca w żadnym
z pozostałych sześciu katalogów tematycznych. Miejsce w interfejsie, w którym każda z dwóch
funkcji się ujawnia, oraz jej zachowanie per moduł i per środowisko pozostaje przedmiotem
opracowania własnego funkcji, nie opracowań środowisk ani modułów.

| Dokument | Klasa | Przeznaczenie | Objętość (znaków) |
|---|---|---|---|
| [Always On Display — dokumentacja projektowa funkcji globalnej](funkcje-globalne/always-on-display.md) | Specyfikacja docelowa | Materiał źródłowy do projektowania i budowy funkcji globalnej Always On Display obowiązującej we wszystkich środowiskach i modułach: postać wizualna, reguły wyzwalania, tor głosowy. | 68 963 |
| [Funkcja globalna Mobile](funkcje-globalne/mobile.md) | Specyfikacja docelowa | Materiał źródłowy do projektowania interfejsu funkcji globalnej Mobile oraz implementacji protokołu parowania urządzeń: komplet okien, oba kanały komunikacji, praca bez połączenia. | 104 052 |

### 1.8 Metoda weryfikacji krzyżowej zbioru

Poprawność nawigacyjna zbioru — brak odsyłacza martwego i brak opracowania osieroconego —
nie jest własnością pojedynczego pliku, lecz własnością całego zbioru naraz, dlatego
sprawdzana jest osobnym przebiegiem, niezależnym od redakcji poszczególnych katalogów
tematycznych. Przebieg obejmuje dwa sprawdzenia uruchamiane po każdej fali redakcji.

**Odsyłacze martwe.** Dla każdego pliku `.md` w `docs/` wyodrębniane są odsyłacze markdown
`[tytuł](ścieżka)` spoza bloków kodu, a każda ścieżka względna rozwiązywana jest względem
katalogu pliku źródłowego i sprawdzana na istnienie. Odsyłacz do kotwicy rozdziału
(`#rozdział`) sprawdzany jest wyłącznie na poziomie pliku docelowego — zgodność samej
kotwicy z nagłówkiem pozostaje przedmiotem lektury, nie automatu.

**Opracowania osierocone.** Plik jest osierocony, jeżeli żadne inne opracowanie zbioru —
poza niniejszym spisem i, gdy dotyczy, Standardem redakcyjnym i językowym — nie odsyła do
niego odsyłaczem markdown. Bycie wskazanym wyłącznie jako „Dokument nadrzędny” samego siebie
nie zwalnia z tego sprawdzenia w drugą stronę: opracowanie potomne jest również odnajdywalne
z opracowania nadrzędnego innego niż wyłącznie niniejszy spis, zgodnie z regułą
dwukierunkowości odsyłaczy (rozdz. 6 Standardu redakcyjnego i językowego).

Wynik obu sprawdzeń nie jest treścią niniejszego spisu — spis pozostaje migawką katalogu,
nie dziennikiem usterek — lecz przedmiotem raportu redakcyjnego przekazywanego redaktorom
właściwych katalogów tematycznych.

---

## 2. Mapa hierarchii dokumentu nadrzędnego

Pole „Dokument nadrzędny” metryki (rozdz. 2.1 Standardu redakcyjnego i językowego) buduje
jedno drzewo obejmujące cały zbiór — każde opracowanie klasy Specyfikacja docelowa wskazuje
dokładnie jednego rodzica, aż do korzenia. Poniższy schemat odtwarza to drzewo z pól
rzeczywiście wypełnionych w chwili niniejszej redakcji.

```
README.md  (Stan wdrożenia — karta produktowa, bez rodzica)
 └─ SPIS-OPRACOWAN.md  (Nawigacja)
     ├─ STANDARD-REDAKCYJNY-I-JEZYKOWY.md
     └─ architektura/koncepcja-platformy.md
         │
         ├─ architektura/architektura.md
         │   ├─ architektura/bezpieczenstwo-i-uwierzytelnianie.md
         │   ├─ architektura/integracja-modeli.md
         │   ├─ architektura/izolacja-i-zaleznosci.md
         │   ├─ architektura/kontrakty-komunikacji.md
         │   ├─ architektura/model-danych.md
         │   ├─ architektura/model-konfiguracji.md
         │   └─ architektura/rozszerzenia.md
         │
         ├─ specyfikacje/specyfikacja-agentow.md
         ├─ specyfikacje/specyfikacja-modulow.md
         │   └─ moduly/{agents · apps · library · research · roundtable ·
         │              studio · terminal · translate · workspace}.md
         ├─ specyfikacje/specyfikacja-okien-operacyjnych.md
         │
         ├─ interfejs-uzytkownika/elementy-okien.md
         │   ├─ interfejs-uzytkownika/okno-historii-sesji.md
         │   ├─ interfejs-uzytkownika/okno-instalatora.md
         │   ├─ interfejs-uzytkownika/okno-instrukcji.md
         │   ├─ interfejs-uzytkownika/okno-nowego-projektu.md
         │   └─ interfejs-uzytkownika/rama-okna.md
         │
         ├─ funkcje-globalne/mobile.md
         ├─ srodowiska/codestudio.md
         └─ srodowiska/multitaskingai.md
```

**Gałęzie nieukończone.** Pole „Dokument nadrzędny” pozostaje niewypełnione w części
opracowań w chwili niniejszej redakcji — stan odnotowany poniżej jako pozycja do
uzupełnienia przez redaktora właściwego katalogu, nie jako rozstrzygnięcie tego spisu
o brzmieniu pola.

| Katalog | Opracowania z niewypełnionym polem „Dokument nadrzędny” |
|---|---|
| `interfejs-uzytkownika/` | `katalog-komponentow.md` · `konfiguracja.md` · `przeplyw-okien.md` · `strona-glowna-i-nawigacja.md` · `system-wizualny.md` · `ustawienia.md` |
| `moduly/` | `assistant.md` · `automations.md` · `browser.md` · `design.md` · `developer.md` · `diagnostics.md` |
| `srodowiska/` | `talkin.md` · `workspace.md` |
| `funkcje-globalne/` | `always-on-display.md` |

Domniemany rodzic każdej pozycji powyżej wynika jednoznacznie z jej katalogu: opracowania
`moduly/` — [Specyfikacja modułów](specyfikacje/specyfikacja-modulow.md); opracowania
`interfejs-uzytkownika/` — [Elementy okien](interfejs-uzytkownika/elementy-okien.md) albo
[Koncepcja platformy](architektura/koncepcja-platformy.md), zależnie od tego, czy
opracowanie dotyczy przepływu głównego czy odrębnego okna; `srodowiska/talkin.md`
i `srodowiska/workspace.md` — [Koncepcja platformy](architektura/koncepcja-platformy.md),
analogicznie do `codestudio.md` i `multitaskingai.md`;
`funkcje-globalne/always-on-display.md` —
[Koncepcja platformy](architektura/koncepcja-platformy.md), analogicznie do `mobile.md`.

---

## 3. Ścieżki czytania

Zbiór czyta się różnie zależnie od roli — deweloper szuka ciągu budowy, projektant szuka
zachowań interfejsu, redaktor szuka wzoru formy. Poniżej trzy sekwencje; każda podaje
kolejność, dokument właściwego kroku i powód, dla którego krok poprzedza kolejne.

### 3.1 Deweloper — od koncepcji do budowy

```
Koncepcja platformy  ──►  Architektura techniczna  ──►  Model danych
        │                        │                          │
        │                        ▼                          ▼
        │              Kontrakty komunikacji        Model konfiguracji
        ▼                        │                          │
  Specyfikacja modułów ◄─────────┴──────────────────────────┘
        │
        ├──►  Specyfikacja agentów  (jeśli budowany element korzysta z warstwy agentowej)
        ├──►  Specyfikacja okien operacyjnych  (katalog wszystkich typów okna)
        ▼
  Opracowanie modułu, nad którym pracuje  ──►  Katalog komponentów (warstwa UI)
```

| Krok | Dokument | Po co |
|---|---|---|
| 1 | [Koncepcja platformy](architektura/koncepcja-platformy.md) | ustala pojęcia i granice, bez których reszta jest nieczytelna |
| 2 | [Architektura techniczna](architektura/architektura.md) | model wdrożenia, stos i mechanizmy przekrojowe |
| 3 | [Model danych](architektura/model-danych.md) · [Kontrakty komunikacji](architektura/kontrakty-komunikacji.md) | encje i komendy, na których stoi każda funkcja |
| 4 | [Model konfiguracji](architektura/model-konfiguracji.md) | warstwowość ustawień obowiązująca w każdym module |
| 5 | [Specyfikacja modułów](specyfikacje/specyfikacja-modulow.md) | zakres funkcjonalny piętnastu modułów do implementacji |
| 6 | [Specyfikacja agentów](specyfikacje/specyfikacja-agentow.md) · [Specyfikacja okien operacyjnych](specyfikacje/specyfikacja-okien-operacyjnych.md) | model agentów i pełny katalog typów okna, gdy budowany element ich dotyczy |
| 7 | opracowanie modułu z katalogu [`moduly/`](#15-moduly) | pełny ciąg budowy jednego modułu — okna, funkcje, przebieg pracy, punkty sterowania |
| 8 | [Katalog komponentów](interfejs-uzytkownika/katalog-komponentow.md) | warstwa wizualna: żetony `--dn-*` i klasy `.dn-*` |

### 3.2 Projektant — od interfejsu do zachowań

| Krok | Dokument | Po co |
|---|---|---|
| 1 | [Strona główna i nawigacja](interfejs-uzytkownika/strona-glowna-i-nawigacja.md) | punkt wejścia i mapa nawigacji |
| 2 | [Rama okna](interfejs-uzytkownika/rama-okna.md) · [Elementy okien](interfejs-uzytkownika/elementy-okien.md) | stała rama i repertuar elementów przepływu głównego |
| 3 | [System wizualny](interfejs-uzytkownika/system-wizualny.md) | żetony, typografia, motywy jasny/ciemny |
| 4 | [Przepływ okien](interfejs-uzytkownika/przeplyw-okien.md) | stany i przejścia między oknami od uruchomienia do pracy |
| 5 | opracowania środowisk z katalogu [`srodowiska/`](#16-srodowiska) | przedsionki i przestrzenie robocze czterech środowisk |
| 6 | [Specyfikacja okien operacyjnych](specyfikacje/specyfikacja-okien-operacyjnych.md) | warstwy widoczności i sposób wywołania każdego typu okna |
| 7 | [Okno Konfiguracji](interfejs-uzytkownika/konfiguracja.md) · [Okno Ustawień](interfejs-uzytkownika/ustawienia.md) | dwa okna poziomu aplikacji odrębne od okien modułowych |

### 3.3 Redaktor dokumentacji — przed każdą edycją

| Krok | Dokument | Po co |
|---|---|---|
| 1 | [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) | wiążący wzór formy, języka i konwencji budowy — rozdziały 2–8 |
| 2 | niniejszy spis, rozdział 1 | ustalenie klasy, katalogu i powiązań redagowanego dokumentu |
| 3 | niniejszy spis, rozdział 4 | próg objętości właściwy klasie redagowanego dokumentu |
| 4 | Standard redakcyjny i językowy, rozdział 12 | polecenia kontrolne uruchamiane przed i po redakcji |

---

## 4. Klasy dokumentów

Cztery klasy zdefiniowane w rozdziale 2.2 Standardu redakcyjnego i językowego wyczerpują
cały zbiór — każde z 51 opracowań należy do dokładnie jednej z nich. Klasa rozstrzyga
o dwóch rzeczach: czy w opracowaniu występują oznaczenia stanu wdrożenia (rozdz. 2.6
Standardu) oraz jaki próg objętości je wiąże (rozdz. 4 niniejszego spisu).

| Klasa | Co opisuje | Ile w zbiorze | Oznaczenia stanu wdrożenia | Próg objętości |
|---|---|---|---|---|
| Specyfikacja docelowa | stan, który ma powstać — zakres budowy | 45 | wyłącznie `[DO DECYZJI OPERATORA]` | ≥ 85 000 znaków |
| Stan wdrożenia | stan faktyczny kodu w bieżącym wydaniu | 3 | pełny zestaw sześcioelementowy | objętość nie mniejsza niż stan zastany |
| Akt prawny | warunki korzystania z produktu | 1 | nie dotyczy | objętość nie mniejsza niż stan zastany |
| Nawigacja | porządek zbioru | 2 | nie dotyczy | ≥ 45 000 znaków |

**Rozmieszczenie klas po katalogach.** Podział jest ścisły: katalog `docs/` w korzeniu miesza
wszystkie cztery klasy (rozdz. 1.1), natomiast sześć podkatalogów tematycznych —
`architektura/`, `specyfikacje/`, `interfejs-uzytkownika/`, `moduly/`, `srodowiska/`,
`funkcje-globalne/` — nosi wyłącznie opracowania klasy Specyfikacja docelowa. Rozbieżność,
czyli dokument klasy odmiennej w jednym z sześciu podkatalogów, jest usterką kwalifikacji,
nie wyjątkiem dopuszczalnym.

**Oznaczenia stanu wdrożenia dopuszczalne w każdej klasie.** Rozdział 2.6 Standardu
redakcyjnego i językowego wiąże sześcioelementowy zestaw oznaczeń wyłącznie do opracowań
klasy Stan wdrożenia; opracowania klasy Specyfikacja docelowa dopuszczają z tego zestawu
wyłącznie jedno oznaczenie.

| Klasa | Repertuar dopuszczalny |
|---|---|
| Stan wdrożenia | `[DZIAŁA]` · `[DZIAŁA CZĘŚCIOWO]` · `[NIEZINTEGROWANE]` · `[ATRAPA]` · `[BRAK]` · `[DO DECYZJI OPERATORA]` |
| Specyfikacja docelowa | wyłącznie `[DO DECYZJI OPERATORA]` |
| Akt prawny, Nawigacja | żadne — oznaczenia stanu wdrożenia nie dotyczą treści prawnej ani porządkującej |

---

## 5. Statystyka objętości zbioru

Rozdział zestawia objętość zbioru liczoną w znakach pliku źródłowego (funkcja `len` języka
Python, ze znacznikami Markdown), zgodnie z miarą przyjętą w rozdziale 9 Standardu
redakcyjnego i językowego. Dane odświeżane są przy każdej fali redakcji zbioru.

| Wskaźnik | Wartość |
|---|---|
| Liczba opracowań | 51 |
| Suma znaków całego zbioru | 5 765 847 |
| Średnia objętość opracowania | 113 054 |
| Najobszerniejsze opracowanie | [Środowisko MultitaskingAI](srodowiska/multitaskingai.md) — 204 803 znaki |
| Najzwięźlejsze opracowanie | [Standard redakcyjny i językowy](STANDARD-REDAKCYJNY-I-JEZYKOWY.md) — 45 879 znaków |

**Objętość według katalogu.**

| Katalog | Liczba opracowań | Suma znaków | Średnia na opracowanie |
|---|---|---|---|
| Korzeń zbioru | 6 | 718 329 | 119 722 |
| architektura/ | 9 | 927 277 | 103 031 |
| specyfikacje/ | 3 | 279 173 | 93 058 |
| interfejs-uzytkownika/ | 12 | 1 385 636 | 115 470 |
| moduly/ | 15 | 1 752 111 | 116 807 |
| srodowiska/ | 4 | 500 101 | 125 025 |
| funkcje-globalne/ | 2 | 203 255 | 101 628 |

**Zgodność z progami rozdziału 9 Standardu redakcyjnego i językowego.** Trzy progi
obowiązują trzy odmienne grupy plików: cztery opracowania Właściciela w korzeniu nie mogą
się skurczyć względem stanu zastanego; opracowania merytoryczne sześciu katalogów
tematycznych wiąże próg 85 000 znaków; opracowania własne zbioru (niniejszy spis, Standard)
wiąże próg 45 000 znaków.

| Grupa | Próg | Spełnia próg | Poniżej progu |
|---|---|---|---|
| Opracowania merytoryczne (6 katalogów tematycznych) | ≥ 85 000 znaków | 45 z 45 | — |
| Opracowania własne zbioru (niniejszy spis, Standard) | ≥ 45 000 znaków | kontrolowane przy każdej redakcji tych dwóch plików | — |
| Cztery opracowania Właściciela w korzeniu | brak zmniejszenia względem stanu zastanego | kontrolowane przy każdej redakcji tych czterech plików | — |

**Objętość według klasy dokumentu.**

| Klasa | Liczba opracowań | Suma znaków | Średnia na opracowanie |
|---|---|---|---|
| Specyfikacja docelowa | 45 | 5 047 553 | 112 168 |
| Stan wdrożenia | 3 | 482 468 | 160 823 |
| Akt prawny | 1 | 144 095 | 144 095 |
| Nawigacja | 2 | 91 766 | 45 883 |

**Odległość od progu 85 000 znaków.** W chwili niniejszej redakcji żadne opracowanie
merytoryczne sześciu katalogów tematycznych nie stoi poniżej progu rozdziału 9 Standardu
redakcyjnego i językowego — wszystkie czterdzieści pięć przekracza 85 000 znaków, a najmniej
liczne z nich, [Okno nowego projektu](interfejs-uzytkownika/okno-nowego-projektu.md),
ma 85 046 znaków.
Tabela odległości od progu pozostaje pusta; wypełnia się ją ponownie dopiero wtedy, gdy
pomiar po fali redakcji wykaże opracowanie poniżej progu.

---

## 6. Zasady utrzymania spisu

Niniejszy spis traci wartość nawigacyjną w chwili, gdy przestaje odpowiadać rzeczywistej
zawartości katalogu `docs/`. Rozdział ustala trzy reguły utrzymania, obowiązujące każdego
redaktora dodającego, przenoszącego lub usuwającego opracowanie.

### 6.1 Dodawanie nowego opracowania

Nowe opracowanie wchodzi do zbioru w tej samej czynności redakcyjnej, w której powstaje jego
plik źródłowy — nie jako osobny krok „później”. Kolejność czynności:

1. Plik powstaje w katalogu tematycznym właściwym jego treści (rozdz. 1.2–1.7); opracowanie
   spoza sześciu katalogów tematycznych i korzenia nie należy do zbioru.
2. Metryka pliku wypełnia komplet jedenastu pól rozdziału 2.1 Standardu redakcyjnego
   i językowego, w tym pole „Dokument nadrzędny” wskazujące istniejące opracowanie zbioru.
3. Wiersz katalogu (rozdz. 1 niniejszego spisu) dodawany jest do tabeli właściwej grupy,
   w miejscu odpowiadającym kolejności czytania tej grupy, nie na końcu tabeli.
4. Wiersz „Mapa hierarchii dokumentu nadrzędnego” (rozdz. 2) rozszerzany jest o nową gałąź,
   zgodnie z polem „Dokument nadrzędny” z kroku 2.
5. Liczba „51 opracowań” przywoływana w rozdziale 1 oraz w metryce niniejszego dokumentu
   aktualizowana jest łącznie w obu miejscach — rozbieżność między nimi jest usterką tego
   samego rzędu co odsyłacz martwy.

### 6.2 Aktualizacja objętości

Kolumna „Objętość (znaków)” w rozdziale 1 oraz statystyki rozdziału 5 są migawką z chwili
ostatniej redakcji niniejszego spisu, nie wartością obliczaną automatycznie przy każdym
odczycie. Po każdej fali redakcji zbioru wartości te wymagają odświeżenia poleceniem podanym
w rozdziale 6.3 — opracowanie zredagowane bez odświeżenia tej kolumny pozostawia spis
niespójny z rzeczywistą objętością pliku.

### 6.3 Weryfikacja automatyczna

Poniższe polecenie sprawdza jeden warunek: czy liczba plików `.md` w `docs/` odpowiada
liczbie zadeklarowanej w rozdziale 1. Warunek drugi — czy każdy plik ma odpowiednik
w tabelach niniejszego spisu — sprawdzany jest porównaniem listy wypisanej tym samym
poleceniem bez `wc -l` z listą odsyłaczy rozdziału 1. Polecenie uruchamiane jest z katalogu
`docs/`.

```bash
find . -name '*.md' | sort | wc -l
```

Wynik różny od 51 oznacza, że rozdział 1 niniejszego spisu wymaga aktualizacji przed
przyjęciem bieżącej fali redakcji — dodania brakującego wiersza albo usunięcia wiersza
wskazującego plik już nieistniejący. Sprawdzenie, który plik odpowiada za rozbieżność,
wykonuje porównanie listy z powyższego polecenia z listą odsyłaczy w tabelach rozdziału 1.

---

## 7. Słowniczek pojęć nawigacyjnych

Pojęcia używane w niniejszym spisie, a nienależące do słownika wiążącego rozdziału 5.1
Standardu redakcyjnego i językowego — właściwe wyłącznie porządkowaniu zbioru, nie treści
produktu.

| Pojęcie | Definicja |
|---|---|
| Katalog tematyczny | Jeden z sześciu podkatalogów `docs/` grupujących opracowania klasy Specyfikacja docelowa według przedmiotu: `architektura/`, `specyfikacje/`, `interfejs-uzytkownika/`, `moduly/`, `srodowiska/`, `funkcje-globalne/` |
| Korzeń zbioru | Katalog `docs/` sam w sobie, poza sześcioma katalogami tematycznymi — miejsce sześciu opracowań o naturze odmiennej od katalogów tematycznych (rozdz. 1.1) |
| Opracowanie Właściciela | Jedno z czterech opracowań korzenia niosących stan faktyczny kodu albo warunki prawne — README, Instrukcja użytkowania, Instalacja i konfiguracja, Licencja — którego objętość nie może się zmniejszyć przy redakcji (rozdz. 9 Standardu) |
| Opracowanie własne zbioru | Jedno z dwóch opracowań korzenia porządkujących sam zbiór dokumentacji, nie produkt — niniejszy spis i Standard redakcyjny i językowy |
| Dokument nadrzędny | Pole metryki (rozdz. 2.1 Standardu) wskazujące opracowanie wyżej w hierarchii; komplet pól całego zbioru składa się w jedno drzewo odtworzone w rozdziale 2 |
| Ścieżka czytania | Uporządkowana sekwencja opracowań właściwa jednej roli — deweloperowi, projektantowi albo redaktorowi — podana w rozdziale 3 |
| Objętość | Liczba znaków pliku źródłowego opracowania, licząc znaki znaczników Markdown, mierzona funkcją `len` języka Python (rozdz. 9 Standardu) |
| Próg objętości | Minimalna (albo, dla opracowań Właściciela, niemalejąca) wartość objętości wiążąca daną klasę dokumentu, zestawiona w rozdziale 4 |
| Fala redakcji | Jeden przebieg redakcyjny obejmujący wiele opracowań zbioru jednocześnie, odnotowywany jako osobny wiersz rozdziału 8 |
| Weryfikacja krzyżowa | Sprawdzenie poprawności nawigacyjnej całego zbioru naraz — odsyłaczy martwych i opracowań osieroconych — metodą opisaną w rozdziale 1.8 |

---

## 8. Historia redakcji

| Data | Zakres fali |
|---|---|
| 2026-08-20 | Ustanowienie [Standardu redakcyjnego i językowego](STANDARD-REDAKCYJNY-I-JEZYKOWY.md); ujednolicenie tożsamości produktu do „Platforma AI Workspace OS · v2.0” w całym zbiorze; naprawa martwych odsyłaczy `[LICENSE](LICENSE)` i odwołań do nieistniejącego drzewa `budowa/docs/`; ujednolicenie cudzysłowów i terminologii; uzupełnienie wykazów komend kontraktu; przemianowanie `INDEKS.md` → `SPIS-OPRACOWAN.md`, `instrukcja-uzytkowania.md` → `okno-instrukcji.md`, `instalacja-i-wdrozenie.md` → `okno-instalatora.md`. |
| 2026-08-20 (fala kolejna) | Redakcja trzech opracowań katalogu `specyfikacje/`: pełna metryka dwutabelowa w Specyfikacji modułów i Specyfikacji okien operacyjnych, zamiana nagłówków „Workflow modułu” na „Przebieg pracy modułu”, potwierdzenie piętnastu obszarów kontraktu i trzech obszarów agentowych w `contract.json`, potwierdzenie liczby 38 prototypów okien, dodanie rozdziałów Kryteria odbioru, maszyn stanów i macierzy uprawnień w trzech specyfikacjach; przebudowa niniejszego spisu na pełny katalog 51 opracowań z objętością i przeznaczeniem każdej pozycji oraz ścieżkami czytania. |
| 2026-08-20 (fala korzenia) | Redakcja sześciu opracowań korzenia według Standardu redakcyjnego i językowego: metryka dwutabelowa w README, Licencji, Instrukcji użytkowania oraz Instalacji i konfiguracji; stopka trzyblokowa o treści dosłownej w czterech opracowaniach Właściciela; ujednolicenie sześcioelementowego zestawu oznaczeń stanu wdrożenia w trzech opracowaniach klasy Stan wdrożenia; usunięcie kodów wymyślonych i dywizu nierozdzielającego; ujednolicenie liczby wierszy katalogu okien operacyjnych do siedemdziesięciu dziewięciu i liczby pozycji nawigacji czterech środowisk do trzydziestu dwóch; separatory rozdziałów i numeracja spisu treści w opracowaniu instalacyjnym; dodanie rozdziału Kryteria odbioru spisu; odświeżenie statystyki objętości zbioru. |

---

## 9. Kryteria odbioru spisu

Rozdział zamykający podaje warunki sprawdzalne, których łączne spełnienie
oznacza, że niniejszy spis jest zgodny z rzeczywistą zawartością katalogu `docs/`
i może zostać przyjęty po fali redakcji. Warunek niespełniony jest usterką spisu,
nie zbioru.

| Warunek | Sposób sprawdzenia |
|---|---|
| Liczba plików `.md` katalogu `docs/` równa liczbie zadeklarowanej w rozdziale 1 | polecenie rozdziału 6.3; wynik różny od 51 wskazuje wiersz do dodania albo do usunięcia |
| Każdy plik zbioru wymieniony w rozdziale 1 dokładnie raz | porównanie listy plików z listą odsyłaczy tabel rozdziału 1 (rozdz. 6.3) |
| Każdy odsyłacz tabel rozdziału 1 wskazuje plik istniejący | kontrola odsyłaczy martwych opisana w rozdziale 1.8 |
| Liczba opracowań zgodna w trzech miejscach: metryce, rozdziale 1 i rozdziale 5 | odczyt trzech miejsc i porównanie wartości |
| Suma kolumny „Liczba opracowań” tabeli „Objętość według katalogu” równa liczbie opracowań zbioru | dodanie siedmiu wierszy tabeli rozdziału 5 |
| Suma kolumny „Liczba opracowań” tabeli „Objętość według klasy dokumentu” równa liczbie opracowań zbioru | dodanie czterech wierszy tabeli rozdziału 5 |
| Każde opracowanie klasy Specyfikacja docelowa ma w rozdziale 2 gałąź albo pozycję w tabeli gałęzi nieukończonych | zestawienie schematu rozdziału 2 z tabelami rozdziału 1 |
| Objętość niniejszego pliku nie mniejsza niż próg opracowania własnego zbioru | pomiar znaków pliku wobec progu z rozdziału 4 |
| Każda fala redakcji odnotowana wierszem rozdziału 8 | odczyt rozdziału 8 po zamknięciu fali |

---

*Koniec dokumentu. Spis opracowań — Nawigacja, wersja 2.0, 2026-08-20.*

---
*Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.*
*Warunki korzystania: [Licencja produktu](LICENSE.md). Kontakt: support@danaco-group.pl*
