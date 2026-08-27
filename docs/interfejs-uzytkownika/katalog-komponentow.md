*Dokument stanowi katalog komponentów interfejsu wielokrotnego użytku platformy Danaco Console wraz z ich formą, wagą wizualną, stanami i warstwami widoczności.*

# Danaco Console — Katalog komponentów interfejsu

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
| **Odbiorcy** | Designer (co, gdzie, w jakiej formie, do czego) · Deweloper (co zbudować) |
| **Źródło** | Koncepcja platformy (autorytatywna); Specyfikacja modułów; Specyfikacja okien operacyjnych; Strona główna i nawigacja; System wizualny; Model konfiguracji; Specyfikacja agentów; Specyfikacja MultitaskingAI; Izolacja i zależności; Rozszerzenia; `tokens.css`, `components.css`, `interfejs-uzytkownika/system-wizualny.md` |

Dokument jest jednym rejestrem wszystkich komponentów interfejsu wielokrotnego użytku platformy Danaco Console — elementów, które nie należą do jednego okna czy modułu, lecz powtarzają się w dziesiątkach miejsc platformy w niezmienionej formie. Odpowiada na pytanie, które biblioteka `components.css` oraz katalog okien operacyjnych zadają z dwóch różnych stron: **co to jest, do czego służy, jaką ma formę i wagę wizualną, jakie przyjmuje stany oraz gdzie dokładnie w platformie występuje.** Dokument porządkuje w jednym miejscu to, co jest rozproszone po bibliotece komponentów, specyfikacji okien operacyjnych, specyfikacji modułów i specyfikacji nawigacji, uzupełniając to o jawne rozstrzygnięcie formy i wagi oraz o przypisanie każdego komponentu do jednej z czterech warstw widoczności (rozdz. 3a).

**Jak korzystać z tego dokumentu.** Designer używa go, aby dla każdego elementu ekranu wiedzieć z góry, czy ma to być mała ikonka (18 px), przycisk (36 px wysokości), zakładka, pole, pigułka czy duży panel — bez zgadywania i bez ryzyka narysowania dużego kwadratu tam, gdzie w bibliotece komponentów jest mała ikonka. Deweloper używa go jako mapy między nazwą produktową komponentu (np. „karta sesji”) a jego klasą w `components.css` (np. `.dn-karta`, `.dn-karta--pozycja`), jego wariantami i dokładnym zestawem stanów do zaimplementowania.

---

## Spis treści

1. [Cel, odbiorcy i zakres dokumentu](#1-cel-odbiorcy-i-zakres-dokumentu)
2. [Źródła](#2-źródła)
3. [Od tokenu do okna — warstwy systemu](#3-od-tokenu-do-okna--warstwy-systemu)
3a. [Reguła warstw widoczności w bibliotece komponentów](#3a-reguła-warstw-widoczności-w-bibliotece-komponentów)
4. [Taksonomia komponentów](#4-taksonomia-komponentów)
5. [Tabela zbiorcza wszystkich komponentów](#5-tabela-zbiorcza-wszystkich-komponentów)
6. [Diagram kompozycji — jak komponenty składają się w okna](#6-diagram-kompozycji--jak-komponenty-składają-się-w-okna)
7. [Karty komponentów — Przyciski i akcje](#7-karty-komponentów--przyciski-i-akcje)
8. [Karty komponentów — Pola formularzy](#8-karty-komponentów--pola-formularzy)
9. [Karty komponentów — Nawigacja](#9-karty-komponentów--nawigacja)
10. [Karty komponentów — Dane i treść](#10-karty-komponentów--dane-i-treść)
11. [Karty komponentów — Informacja zwrotna i nakładki](#11-karty-komponentów--informacja-zwrotna-i-nakładki)
12. [Karty komponentów — Tożsamość](#12-karty-komponentów--tożsamość)
13. [Elementy pomocnicze](#13-elementy-pomocnicze)
14. [Karty komponentów — Okna komunikacji operacyjnej](#14-karty-komponentów--okna-komunikacji-operacyjnej)
15. [Karty komponentów — Ujawnianie funkcjonalności](#15-karty-komponentów--ujawnianie-funkcjonalności)
16. [Macierz komponent × moduł / okno](#16-macierz-komponent--moduł--okno)
17. [Zgodność z zasadami nadrzędnymi platformy](#17-zgodność-z-zasadami-nadrzędnymi-platformy)
18. [Słowniczek pojęć](#18-słowniczek-pojęć)
- [Załącznik A. Ściągawka wymiarów i tokenów](#załącznik-a-ściągawka-wymiarów-i-tokenów)
- [Załącznik B. Uwaga o rozbieżności źródeł systemu wizualnego](#załącznik-b-uwaga-o-rozbieżności-źródeł-systemu-wizualnego)
- [Załącznik C. Pełna lista klas CSS komponentów](#załącznik-c-pełna-lista-klas-css-komponentów)

---

## 1. Cel, odbiorcy i zakres dokumentu

| Pytanie | Odpowiedź |
|---|---|
| Co obejmuje dokument | Wszystkie komponenty interfejsu wielokrotnego użytku platformy — elementy powtarzalne w wielu oknach, modułach i środowiskach, zbudowane na jednym, wspólnym systemie tokenów (granat + złoto) |
| Czego dokument nie obejmuje | Pojedynczych okien operacyjnych jako całości (przedmiot Specyfikacji okien operacyjnych), modeli danych (przedmiot Modelu danych), kontraktów komunikacji (przedmiot Kontraktów komunikacji), ostatecznych makiet graficznych poszczególnych ekranów (przedmiot dokumentacji wykonawczej ekranów) |
| Dla kogo | Designer — co narysować, w jakiej formie, jakiej wielkości, w jakich stanach. Deweloper — jaką klasę zaimplementować, z jakimi wariantami i jakim zachowaniem |
| Status komponentu wobec źródeł | Komponent nazwany i mający gotową klasę w `components.css` — oznaczony jako **zdefiniowany**. Komponent opisany na poziomie produktowym w specyfikacjach, lecz bez odrębnej klasy CSS — oznaczony jako **wzorzec złożony** (składany z komponentów zdefiniowanych) |

Dokument nie tworzy nowej hierarchii ważności komponentów ani nowych ograniczeń ich stosowania — każdy komponent jest dostępny w każdym miejscu platformy, w którym jego zastosowanie ma sens funkcjonalny, zgodnie z zasadą pełnej kompozycyjności (Koncepcja platformy, rozdział 14, zasada 1).

---

## 2. Źródła

| Dokument | Co wniósł do niniejszego katalogu |
|---|---|
| Koncepcja platformy (autorytatywna) | Definicje środowisk, modułów, funkcji globalnych (Mobile, Always On Display); trzy strefy strony głównej; zasady nadrzędne |
| Specyfikacja modułów | Okna operacyjne piętnastu modułów wraz z ich zawartością — źródło większości pozycji kolumny „Gdzie używany” |
| Specyfikacja okien operacyjnych | Anatomia okna operacyjnego; klasyfikacja okien według charakteru pracy; katalog okien robocze ról MultitaskingAI; okno konfiguracji punktów izolacji |
| Strona główna i nawigacja | Anatomia karty środowiska, kafla komponentu własnego, listwy ustawień; pełna mechanika karty sesji; boczna nawigacja modułów i panel orkiestracji |
| System wizualny | Tokeny kolorów, typografii, odstępów; zastosowanie biblioteki komponentów w typologii okien i w trzech strefach strony głównej |
| Model konfiguracji | Struktura i nawigacja okna konfiguracji; warstwowość ustawień; mechanizm objaśnień kontekstowych `[?]`; uproszczone menu kontekstowe |
| Specyfikacja agentów | Okna modułu Agents; komponenty definicji agenta; Centrum uprawnień |
| Specyfikacja MultitaskingAI | Panel orkiestracji; role i ich okna robocze; warstwa centralna Always On Display |
| Always On Display — funkcja globalna | Postać wizualna awatara i powierzchni interakcji; stany funkcji; warstwy widoczności elementów funkcji (`funkcje-globalne/always-on-display.md`) |
| Izolacja i zależności | Macierz izolacji jako największe w platformie skupisko przełączników |
| Rozszerzenia | Rejestr rozszerzeń jako wzorzec listy z przełącznikiem stanu włączenia |
| `tokens.css` | Źródło prawdy dla wartości pikselowych i barwnych — wykorzystane bezpośrednio w kolumnach „Forma i waga” |
| `components.css` | Źródło prawdy dla klas, wariantów i reguł stanu — wykorzystane bezpośrednio w każdej karcie komponentu |
| `interfejs-uzytkownika/system-wizualny.md` | Zasady marki, dostępność WCAG 2.1 AA, ikonografia |

Wszystkie nazwy własne środowisk, modułów, okien i komponentów przejęto bez zmian ze źródeł. Tam, gdzie dwa źródła różnią się w szczególe (Załącznik B), rozbieżność jest odnotowana jawnie, a nie milcząco rozstrzygnięta na rzecz jednego z nich.

---

## 3. Od tokenu do okna — warstwy systemu

**Diagram — łańcuch warstw systemu wizualnego, od prymitywu koloru do środowiska platformy.**

```
 WARSTWA 1 · PRYMITYWY                 surowe skale kolorów, niewykorzystywane wprost
   granat 900…100 + marka
   złoto  700…100 + ghost
   grafit 900…100 (tło trybu ciemnego)          tokens.css
   szarość 0…900
        │
        ▼
 WARSTWA 2 · TOKENY SEMANTYCZNE        wartość zależna od motywu jasny/ciemny
   --dn-tlo · --dn-powierzchnia · --dn-panel
   --dn-tekst · --dn-akcent · --dn-focus         tokens.css / tokens.json
   --dn-sukces-fg/bg/br · --dn-blad-fg/bg/br …
        │
        ▼
 WARSTWA 3 · KOMPONENTY ZDEFINIOWANI   klasy .dn-*, gotowe do użycia
   .dn-btn · .dn-input · .dn-karta · .dn-tabela
   .dn-modal · .dn-toast · .dn-tooltip …         components.css (rozdz. 7–15
   .dn-suwak · .dn-plakietka · .dn-awatar        niniejszego dokumentu)
        │
        ▼
 WARSTWA 4 · WZORCE ZŁOŻONE            komponenty zdefiniowani złożeni w większą całość
   Chat Window · Execution Loop Window
   Karta sesji · Panel · Plansza · Menu
   kontekstowe · Avatar Always On Display        rozdz. 9–15 niniejszego dokumentu
   (nie mają odrębnej klasy — składają się
   z komponentów warstwy 3)
        │
        ▼
 WARSTWA 5 · OKNO OPERACYJNE           np. Studio Editor, Agent Builder, Execution Monitor
   zestaw komponentów i wzorców złożonych
   realizujący jedną funkcję modułu              Specyfikacja okien operacyjnych
        │
        ▼
 WARSTWA 6 · MODUŁ                     zestaw okien operacyjnych (np. Studio = 6 okien)
        │                                        Specyfikacja modułów
        ▼
 WARSTWA 7 · ŚRODOWISKO                TalkIn · WorkSpace · CodeStudio · MultitaskingAI
                                                  Koncepcja platformy
```

Zasada obowiązująca w całym łańcuchu: każda warstwa odwołuje się wyłącznie do warstwy bezpośrednio niższej — okno operacyjne nigdy nie koduje koloru na sztywno (odwołuje się do komponentu), komponent nigdy nie koduje wartości barwy na sztywno (odwołuje się do tokenu semantycznego), token semantyczny nigdy nie jest wartością przypadkową (odwołuje się do prymitywu). Ta sama warstwa 3 obsługuje wszystkie warstwy 5–7 bez rozgałęziania się na warianty specyficzne dla środowiska — jeden `.dn-btn` wygląda identycznie w Studio, w Agent Builderze i w oknie konfiguracji punktów izolacji.

---

## 3a. Reguła warstw widoczności w bibliotece komponentów

Zasadą nadrzędną interfejsu Danaco Console jest stopniowe ujawnianie funkcjonalności: **jeżeli funkcja nie jest potrzebna do realizacji aktualnego zadania, nie jest widoczna**. Biblioteka komponentów realizuje tę zasadę wprost — każdy komponent katalogu należy do dokładnie jednej z czterech warstw widoczności, a jego karta (rozdziały 7–15) podaje pola „Warstwa widoczności” i „Sposób wywołania”. Złożoność platformy istnieje w architekturze i pozostaje niewidoczna w interfejsie do chwili wystąpienia potrzeby użycia danej funkcji. Liczba modułów, agentów, przepływów pracy, komponentów, narzędzi, paneli i ustawień nie wpływa na postrzeganą prostotę interfejsu.

**Tabela — cztery warstwy widoczności i komponenty katalogu, które je obsługują.**

| Warstwa | Nazwa | Zawartość | Sposób dostępu | Komponenty katalogu |
|---|---|---|---|---|
| 1 | Zawsze widoczna | Chat Window, aktywne okno wiodące, kontekst pracy, podstawowa nawigacja, wskaźniki stanu wykonania. Zajmuje ponad 80% powierzchni interfejsu | Widoczna bez interakcji | Chat Window (14.1), Execution Loop Window (14.2), Pasek górny (9.1), Boczna nawigacja (9.2), Karta sesji (9.3), Zakładka (9.4), Karta (10.1), Panel (10.2), Plansza (10.3), Tabela (10.4), Pigułka i plakietka (10.5), Pasek kontekstu ze znacznikami (15.9) |
| 2 | Widoczna na żądanie | Wybór modelu, wybór wykonawcy, wybór środowiska, wybór trybu pracy, poziom wysiłku, parametry przepływu pracy | Ikona, przycisk, przełącznik, znacznik kontekstowy; po użyciu element zwija się samoczynnie | Menu progresywne (15.2), Znacznik kontekstowy (15.7), Lista rozwijana (8.2), Przełącznik (8.3), Pole wyboru (8.4), Dymek kontekstowy (11.4) |
| 3 | Rozwinięcia kontekstowe | Zestawy akcji, ustawienia szybkie, warianty operacji | Menu kebab (⋮), menu hamburger (☰), menu kontekstowe, panel popover, panel wysuwany, lista rozwijana | Element zbiorczy z rozwinięciem (15.1), Menu kebab (15.3), Menu hamburger (15.4), Panel popover (15.5), Panel wysuwany (15.6), Menu kontekstowe (9.5), Modal (11.1) |
| 4 | Funkcje eksperckie | Najbardziej zaawansowane operacje, tryby administracyjne, narzędzia diagnostyczne niskiego poziomu | Polecenie języka naturalnego w Chat Window, skrót klawiszowy, wyszukiwarka funkcji, tryb administracyjny, konfiguracja roli. Użytkownik podstawowy nie widzi tych elementów | Paleta poleceń (15.8), Chat Window jako kanał poleceń (14.1) |

**Mechanizmy ukrywania funkcjonalności w bibliotece komponentów.**

| Mechanizm | Komponent nośny | Zasada działania |
|---|---|---|
| Menu progresywne | Menu progresywne (15.2) | Zbiór jednorodnych wyborów prezentowany jest jako jeden element zwinięty (`Agent ▼`); lista pozycji rozwija się po kliknięciu i zwija po wyborze |
| Panele wysuwane | Panel wysuwany (15.6), Panel popover (15.5) | Funkcjonalność umieszczana jest w panelach bocznych, wysuwanych i kontekstowych; po zamknięciu panel znika całkowicie z przestrzeni roboczej i nie pozostawia śladu w układzie |
| Grupowanie logiczne akcji | Element zbiorczy z rozwinięciem (15.1), Menu kebab (15.3), Menu hamburger (15.4) | Zamiast zestawu przycisków prezentowany jest jeden element zbiorczy (`Operacje ▼`, ⋮, ☰), którego rozwinięcie zawiera pełną listę akcji |
| Znaczniki kontekstowe | Znacznik kontekstowy (15.7), Pasek kontekstu ze znacznikami (15.9) | Środowisko, repozytorium, projekt, model i wykonawca występują jako lekkie znaczniki w pasku kontekstu, na przykład `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]`; kliknięcie znacznika otwiera odpowiedni selektor |
| Kanał języka naturalnego | Chat Window (14.1) | Funkcje warstwy 4 wywoływane są poleceniem w języku naturalnym, bez własnej reprezentacji graficznej w interfejsie |

**Zasada jednego kliknięcia.** Każda ukryta funkcja jest osiągalna jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Zagnieżdżanie funkcji głęboko w hierarchii menu jest zabronione — ukrycie zmniejsza chaos wizualny, nie utrudnia dostępu. Komponent rozwijający (15.1–15.6) prezentuje pełną listę pozycji na jednym poziomie; komponent wymagający drugiego poziomu rozwinięcia nie występuje w katalogu.

**Zasada rysowania makiet.** Wszystkie makiety ASCII niniejszego dokumentu przedstawiają interfejs w stanie spoczynku: widoczne są wyłącznie elementy warstwy 1 oraz zwinięte wyzwalacze warstw 2–3 (znaczniki, `▼`, `⋮`, `☰`). Elementy warstw 2–4 opisane są w tabelach elementów okna wraz z podaniem warstwy i sposobu wywołania.

---

## 4. Taksonomia komponentów

**Schemat — dziewięć kategorii i trzydzieści osiem komponentów katalogu.**

```
KATALOG KOMPONENTÓW DANACO CONSOLE
│
├── A. PRZYCISKI I AKCJE                         rozdz. 7
│     A1 Przycisk                    A2 Przycisk ikonowy
│
├── B. POLA FORMULARZY                           rozdz. 8
│     B1 Pole tekstowe (Input/Textarea)
│     B2 Lista rozwijana (Select)
│     B3 Przełącznik (Toggle)
│     B4 Pole wyboru / opcja jednokrotna (Checkbox/Radio)
│
├── C. NAWIGACJA                                 rozdz. 9
│     C1 Pasek górny              C2 Boczna nawigacja / Panel orkiestracji
│     C3 Karta sesji              C4 Zakładka
│     C5 Menu kontekstowe
│
├── D. DANE I TREŚĆ                               rozdz. 10
│     D1 Karta                    D2 Panel
│     D3 Plansza                  D4 Tabela
│     D5 Pigułka i plakietka statusu
│     D6 Blok kodu                D7 Pusty stan
│
├── E. INFORMACJA ZWROTNA I NAKŁADKI              rozdz. 11
│     E1 Modal                    E2 Powiadomienie (Toast)
│     E3 Komunikat blokowy (Alert)
│     E4 Dymek kontekstowy (Tooltip)
│     E5 Wskaźnik ładowania (Loader/Spinner)
│     E6 Centrum powiadomień
│
├── F. TOŻSAMOŚĆ                                  rozdz. 12
│     F1 Awatar                   F2 Avatar Always On Display
│
├── G. ELEMENTY POMOCNICZE                        rozdz. 13
│     G1 Ikona systemowa          G2 Separator
│
├── H. OKNA KOMUNIKACJI OPERACYJNEJ               rozdz. 14
│     H1 Chat Window (Użytkownik ↔ Wykonawca)
│     H2 Execution Loop Window (Koordynator ↔ Wykonawca)
│
└── I. UJAWNIANIE FUNKCJONALNOŚCI                 rozdz. 15
      I1 Element zbiorczy z rozwinięciem (`Operacje ▼`)
      I2 Menu progresywne (`Agent ▼`)
      I3 Menu kebab (⋮)            I4 Menu hamburger (☰)
      I5 Panel popover             I6 Panel wysuwany
      I7 Znacznik kontekstowy (tag) otwierający selektor
      I8 Paleta poleceń            I9 Pasek kontekstu ze znacznikami
```

Podział na dziewięć kategorii porządkuje katalog wg funkcji komponentu w interfejsie (co robi), nie wg modułu, w którym akurat występuje — ten sam `.dn-plakietka` służy zarówno do oznaczenia statusu przebiegu w Execution Monitor, jak i etykiety zasobu w module Library.

---

## 5. Tabela zbiorcza wszystkich komponentów

| # | Komponent | Klasa / wzorzec | Kategoria | Forma i waga | Warstwa widoczności | Status wobec źródeł |
|---|---|---|---|---|---|---|
| A1 | Przycisk | `.dn-btn` | Akcja | Mała, w linii tekstu (wys. ok. 36 px domyślnie) | 1 | Zdefiniowany |
| A2 | Przycisk ikonowy | `.dn-btn-ikona` | Akcja | Znikoma, kwadrat 36×36 px | 1 | Zdefiniowany |
| B1 | Pole tekstowe | `.dn-input` / `.dn-textarea` | Formularz | Średnia, pełna szerokość kontenera | 1 | Zdefiniowany |
| B2 | Lista rozwijana | `.dn-select` | Formularz | Średnia, pełna szerokość kontenera | 2 | Zdefiniowany |
| B3 | Przełącznik | `.dn-suwak` | Formularz | Znikoma, 40×22 px | 2 | Zdefiniowany |
| B4 | Pole wyboru / opcja jednokrotna | `.dn-check` | Formularz | Znikoma, 16×16 px | 2 | Zdefiniowany |
| C1 | Pasek górny | `.dn-pasek` | Nawigacja | Duża, pełna szerokość okna, wys. 56 px | 1 | Zdefiniowany |
| C2 | Boczna nawigacja / Panel orkiestracji | `.dn-karta--pozycja` w kolumnie | Nawigacja | Duża, pełna wysokość okna, stała szerokość | 1 | Wzorzec złożony |
| C3 | Karta sesji | `.dn-karta--pozycja` w rzędzie kart | Nawigacja | Mała, w rzędzie kart sesji pod paskiem górnym | 1 | Wzorzec złożony |
| C4 | Zakładka | `.dn-zakladki` / `.dn-zakladka` | Nawigacja | Mała, w poziomym rzędzie | 1 | Zdefiniowany |
| C5 | Menu kontekstowe | `.dn-karta` + `.dn-karta--pozycja` (pływające) | Nawigacja | Mała do średniej, pływająca nakładka przy wyzwalaczu | 3 | Wzorzec złożony |
| D1 | Karta | `.dn-karta` | Dane i treść | Zmienna: od małej pozycji listy po dużą kartę środowiska | 1 | Zdefiniowany |
| D2 | Panel | złożenie `.dn-karta` / `.dn-tabela` / `.dn-pole` w regionie | Dane i treść | Średnia do dużej, region wewnątrz okna operacyjnego | 1 | Wzorzec złożony |
| D3 | Plansza | obszar główny okna, bez stałych wymiarów | Dane i treść | Duża, cały obszar główny okna | 1 | Wzorzec złożony |
| D4 | Tabela | `.dn-tabela` + `.dn-tabela-owijka` | Dane i treść | Duża, pełna szerokość kontenera, wysokość wg liczby wierszy | 1 | Zdefiniowany |
| D5 | Pigułka i plakietka statusu | `.dn-plakietka` / `.dn-kropka` | Dane i treść | Znikoma, pigułka w linii tekstu | 1 | Zdefiniowany |
| D6 | Blok kodu | `.dn-kod` | Dane i treść | Średnia, blok pełnej szerokości kontenera | 1 | Zdefiniowany |
| D7 | Pusty stan | `.dn-pusty-stan` | Dane i treść | Średnia, wypełnia pusty panel lub okno | 1 | Zdefiniowany |
| E1 | Modal | `.dn-nakladka` + `.dn-modal` | Nakładka | Duża, do 560 px szerokości, nakładka pełnoekranowa | 3 | Zdefiniowany |
| E2 | Powiadomienie (Toast) | `.dn-toast` | Nakładka | Mała, 260–440 px szerokości, stos przy prawej krawędzi obszaru roboczego | 1 | Zdefiniowany |
| E3 | Komunikat blokowy (Alert) | `.dn-alert` | Nakładka (osadzona) | Mała do średniej, osadzony w układzie, pełna szerokość rodzica | 1 | Zdefiniowany |
| E4 | Dymek kontekstowy (Tooltip) | `.dn-tooltip` | Nakładka | Znikoma, dymek nad wyzwalaczem | 2 | Zdefiniowany |
| E5 | Wskaźnik ładowania (Loader/Spinner) | `.dn-spinner` | Nakładka (osadzona) | Znikoma, 16×16 px | 1 | Zdefiniowany |
| E6 | Centrum powiadomień | `.dn-panel-wysuwany` + `.dn-karta--pozycja` + `.dn-plakietka--stan` | Nakładka (kolumna boczna) | Duża, kolumna boczna 320–420 px, pełna wysokość obszaru roboczego | 3 | Wzorzec złożony |
| F1 | Awatar | `.dn-awatar` | Tożsamość | Znikoma do małej, koło/kwadrat 28–48 px | 1 | Zdefiniowany |
| F2 | Avatar Always On Display | `.dn-awatar` pływający + powierzchnia interakcji | Tożsamość | Mała w spoczynku, rozszerza się po aktywacji | 1 | Wzorzec złożony |
| G1 | Ikona systemowa | `.dn-ikona` | Pomocniczy | Znikoma, 14–24 px | 1 | Zdefiniowany |
| G2 | Separator | `.dn-separator` | Pomocniczy | Znikoma, linia 1 px | 1 | Zdefiniowany |
| H1 | Chat Window | złożenie `.dn-karta` / `.dn-input` / `.dn-btn-ikona` w kolumnie | Okno komunikacji | Duża, lewa kolumna, stała, pełna wysokość obszaru roboczego | 1 | Wzorzec złożony |
| H2 | Execution Loop Window | złożenie `.dn-karta--pozycja` / `.dn-tabela` / `.dn-plakietka` w kolumnie | Okno komunikacji | Duża, kolumna sąsiadująca z Chat Window, pełna wysokość | 1 | Wzorzec złożony |
| I1 | Element zbiorczy z rozwinięciem | `.dn-btn` + `.dn-karta` (lista rozwinięta) | Ujawnianie | Mała w spoczynku (`Operacje ▼`), rozwinięcie jako lista pływająca | 3 | Wzorzec złożony |
| I2 | Menu progresywne | `.dn-btn--duch` + `.dn-karta--pozycja` | Ujawnianie | Znikoma w spoczynku (`Agent ▼`) | 2 | Wzorzec złożony |
| I3 | Menu kebab | `.dn-btn-ikona` (ikona `wiecej`) + `.dn-karta` | Ujawnianie | Znikoma, kwadrat 36×36 px | 3 | Wzorzec złożony |
| I4 | Menu hamburger | `.dn-btn-ikona` (ikona `menu`) + `.dn-karta` | Ujawnianie | Znikoma, kwadrat 36×36 px | 3 | Wzorzec złożony |
| I5 | Panel popover | `.dn-karta` pływająca zakotwiczona przy wyzwalaczu | Ujawnianie | Mała do średniej, nakładka lekka | 3 | Wzorzec złożony |
| I6 | Panel wysuwany | region `.dn-panel` jako rozszerzenie boczne | Ujawnianie | Średnia do dużej, kolumna boczna otwierana z prawej strony obszaru roboczego | 3 | Wzorzec złożony |
| I7 | Znacznik kontekstowy | `.dn-plakietka--interaktywna` | Ujawnianie | Znikoma, pigułka w linii tekstu | 2 | Zdefiniowany |
| I8 | Paleta poleceń | `.dn-paleta` (`.dn-input` + `.dn-karta--pozycja`) | Ujawnianie | Średnia, nakładka wyśrodkowana, wywoływana skrótem `Ctrl/Cmd + K` | 4 | Wzorzec złożony |
| I9 | Pasek kontekstu ze znacznikami | rząd `.dn-plakietka--interaktywna` | Ujawnianie | Znikoma, jeden wiersz znaczników nad treścią kolumny | 1 | Wzorzec złożony |

Dwadzieścia dwa z trzydziestu ośmiu komponentów mają odrębną klasę CSS gotową do użycia bez dalszej kompozycji („Zdefiniowany”). Szesnaście jest wzorcem złożonym — nazwanym i opisanym na poziomie produktowym w źródłach, lecz zbudowanym ze składników warstwy 3 (rozdz. 3) — dla nich karta komponentu w rozdziałach 9–15 podaje jawnie, z jakich komponentów zdefiniowanych się składa.

---

## 6. Diagram kompozycji — jak komponenty składają się w okna

### 6.1. Powłoka aplikacji — kompozycja makro

**Diagram — powłoka aplikacji w układzie pionowym (podział lewa–prawa), stan spoczynku interfejsu.**

```
 PASEK GÓRNY (C1 · .dn-pasek, 56 px)
   ◈ Danaco   [ szukaj ]                      ikony (A2)    AWATAR (F1)
 KARTY SESJI (C3):  [ ● Studio ✕ ] [ Research ✕ ] [ + ]
 ═══════════════════════════════════════════════════════════════════════════
  BOCZNA      │ CHAT WINDOW (H1)     │ OBSZAR ROBOCZY OKNA  │ PANEL
  NAWIGACJA   │ Użytkownik ↔         │ OPERACYJNEGO         │ WYSUWANY (I6)
  / PANEL     │ Wykonawca            │                      │ (rozszerzenie
  ORKIESTRACJI│                      │ zbudowany z:         │  boczne,
  (C2)        │ pasek kontekstu (I9) │  PANEL (D2)          │  otwierane
              │ [Danaco Console]     │  TABELA (D4)         │  na żądanie)
  .dn-karta   │ [Ubuntu] [Fable 5]   │  KARTA (D1)          │
  --pozycja   │                      │  PLANSZA (D3)        │
  ×N          │ strumień wiadomości  │  pola formularzy (B) │
              │ awatar (F1)          │  pigułki (D5)        │
              │ blok kodu (D6)       │  blok kodu (D6)      │
              │                      │  pusty stan (D7)     │
              │ [ pole B1 ] (A2)     │            ⋮ (I3)    │
              │ Agent ▼ (I2)         │  Operacje ▼ (I1)     │
              │──────────────────────│                      │
              │ EXECUTION LOOP       │                      │
              │ WINDOW (H2)          │                      │
              │ Koordynator ↔        │                      │
              │ Wykonawca            │                      │
 ═══════════════════════════════════════════════════════════════════════════
 Regulacji podlega wyłącznie szerokość kolumn.

WARSTWA NAKŁADEK (ponad całą powłoką, niezależnie od okna operacyjnego)
   MODAL (E1) — środek ekranu, na żądanie
   TOAST (E2) — stos przy prawej krawędzi obszaru roboczego, samoczynnie znika
   DYMEK KONTEKSTOWY (E4) — przy najechaniu/fokusie dowolnego elementu z [?]
   MENU KONTEKSTOWE (C5) · PANEL POPOVER (I5) — przy wyzwalaczu ⋮ / ☰
   PALETA POLECEŃ (I8) — skrót `Ctrl/Cmd + K`, warstwa 4
   AVATAR ALWAYS ON DISPLAY (F2) — pływający, ponad całą powłoką, zawsze
```

Chat Window (rozdz. 14.1) zajmuje w powłoce miejsce stałe — lewą kolumnę obszaru roboczego, na pełnej wysokości, w każdym module i w każdym środowisku. Execution Loop Window (rozdz. 14.2) otwiera się jako kolumna sąsiadująca. Obszar roboczy modułu zajmuje kolumnę dominującą po ich prawej stronie, a okna pomocnicze i panele otwierają się jako kolejne kolumny boczne po prawej stronie obszaru roboczego.

### 6.2. Anatomia okna operacyjnego → mapowanie na komponenty

**Schemat — cztery elementy anatomii okna operacyjnego (Specyfikacja okien operacyjnych, rozdz. 2.4) i komponent, który je realizuje.**

```
┌─ Nagłówek okna ──────────────────────────────────────────────────────
│   nazwa okna                              → tekst nagłówkowy + PASEK/KARTA (D1) nagłówek
├────────────────────────────────────────────────────────────────────
│   OBSZAR ZAWARTOŚCI
│     edytor              → PLANSZA (D3) albo pole tekstowe pełnoekranowe
│     lista                → PANEL (D2) z KARTA--pozycja (D1) albo TABELA (D4)
│     monitor               → TABELA (D4) + PIGUŁKA--stan (D5)
│     podgląd               → PANEL (D2) z treścią wyrenderowaną
├────────────────────────────────────────────────────────────────────
│   ELEMENTY KONFIGURACJI
│     ustawienie + [?]      → pole (B1/B2/B3/B4) + DYMEK KONTEKSTOWY (E4)
├────────────────────────────────────────────────────────────────────
│   STAN PROCESU SESJI
│     wskaźnik pracy w tle  → KROPKA stanu (część D5) na KARCIE SESJI (C3)
├────────────────────────────────────────────────────────────────────
│   KANAŁY KOMUNIKACJI OPERACYJNEJ
│     Użytkownik ↔ Wykonawca   → CHAT WINDOW (H1), lewa kolumna
│     Koordynator ↔ Wykonawca  → EXECUTION LOOP WINDOW (H2), kolumna
│                                 sąsiadująca
└────────────────────────────────────────────────────────────────────
```

### 6.3. Przykład złożony — moduł Studio

**Makieta — moduł Studio złożony z komponentów katalogu (okna wg Specyfikacji okien operacyjnych, rozdz. 6.1).**

```
 PASEK GÓRNY (C1)
 KARTY SESJI (C3): [ ● Studio ✕ ]
 ═══════════════════════════════════════════════════════════════════════════
  BOCZNA      │ CHAT WINDOW (H1)   │ Studio Editor        │ PANEL WYSUWANY
  NAWIGACJA   │ Użytkownik ↔       │  = PLANSZA (D3)      │ (I6)
  (C2)        │ Wykonawca          │    tekstowa: treść   │  Tools Panel:
   Studio ▸   │                    │    dokumentu,        │  KARTA--pozycja
   Workspace  │ AWATAR (F1)        │    zaznaczenie       │  (D1) + PRZYCISK
   Browser    │ strumień wiadomości│    fragmentu         │  (A1) per akcja
   Research   │ blok kodu (D6)     │                      │
   Library    │                    │  Operacje ▼ (I1)     │  Diff/Grep Panel
   Translate  │ [ pole B1 ] (A2)   │                      │  = TABELA (D4)
   Roundtable │ Agent ▼ (I2)       │                      │
   Assistant  │────────────────────│                      │  Session
   Agents     │ EXECUTION LOOP     │                      │  Repository
              │ WINDOW (H2)        │                      │  = PANEL (D2)
              │ Koordynator ↔      │                      │  KARTA--pozycja
              │ Wykonawca          │                      │  (D1) per wersja
 ═══════════════════════════════════════════════════════════════════════════
```

### 6.4. Przykład złożony — Agent Builder

**Makieta — moduł Agents, okno Agent Builder (Specyfikacja agentów, rozdz. 4.1) złożone z komponentów katalogu.**

```
 Agent Builder — okno nadrzędne = MODAL (E1) albo PLANSZA (D3) pełnoekranowa
 ┌─────────────────────────────────────────────────────────────────────┐
 │ Tożsamość:  Nazwa [pole B1]   Opis [pole B1 wieloliniowe]           │
 │             Zasięg widoczności  [ Select B2: globalny / projektowy ]│
 ├─────────────────────────────────────────────────────────────────────┤
 │ ┌─ Model Configuration ────┐ ┌─ Skills Manager ─────────┐           │
 │ │ Model bazowy [B2]        │ │ Umiejętności dobrane:    │           │
 │ │ Kanał: API/CLI/SSH/HTTP  │ │ [pigułka D5] [pigułka]   │           │
 │ │ [zakładki C4 lub B2]     │ │ Dodaj [PRZYCISK A1]      │           │
 │ └──────────────────────────┘ └──────────────────────────┘           │
 │ ┌─ Connectors Manager ─────┐ ┌─ Permissions Center ─────┐           │
 │ │ TABELA (D4): konektor    │ │ TABELA (D4): uprawnienie │           │
 │ │ + PRZEŁĄCZNIK (B3)       │ │ + PRZEŁĄCZNIK (B3)       │           │
 │ │ w wierszu                │ │ w wierszu                │           │
 │ └──────────────────────────┘ └──────────────────────────┘           │
 ├─────────────────────────────────────────────────────────────────────┤
 │                             [ PRZYCISK --zloty A1 „Utwórz agenta” ] │
 └─────────────────────────────────────────────────────────────────────┘
```

### 6.5. Przykład złożony — środowisko MultitaskingAI

**Makieta — panel orkiestracji i cztery okna robocze ról (Specyfikacja MultitaskingAI, rozdz. 9; Specyfikacja okien operacyjnych, rozdz. 7.2).**

```
 ═══════════════════════════════════════════════════════════════════════════
  PANEL          │ CHAT WINDOW (H1)  │ EXECUTION LOOP   │ Executor Chat (1)
  ORKIESTRACJI   │ Użytkownik ↔      │ WINDOW (H2)      │ KARTA--pozycja
  (C2)           │ Wykonawca         │ Koordynator ↔    │ (D1) × do 15
                 │                   │ Wykonawca        │ (Subagent Network)
   Zespoły       │ pasek kontekstu   │                  │
   Role ▸        │ (I9): [Zespół A]  │ zlecenie i jego  │ Executor Chat (2)
    → 4 × KARTA  │ [Fable 5] [Ultra] │ dekompozycja     │ KARTA--pozycja
      (D1) roli, │                   │ kolejka zadań    │ (D1) × do 15
      AWATAR     │ strumień poleceń  │  = TABELA (D4)   │
      (F1)       │ i wyników         │ komunikaty       │ Results Analyzer
      --kwadrat  │                   │ sterujące        │ (Executor 3 /
   Kolejki       │ [ pole B1 ] (A2)  │ wskaźniki pętli  │  Validator)
    → TABELA (D4)│ Agent ▼ (I2)      │  PIGUŁKA--stan   │ TABELA (D4) ocen
   Orkiestracja  │                   │  (D5)            │ + PIGUŁKA--stan
    → KARTA      │                   │ Operacje ▼ (I1): │  (D5)
      --pozycja  │                   │ wstrzymaj,       │
      (D1)       │                   │ wznów, przerwij, │ AVATAR ALWAYS ON
   Harmonogram   │                   │ skoryguj         │ DISPLAY (F2) —
   i automatyki  │                   │                  │ pływający ponad
   Monitor       │                   │                  │ całą powłoką,
   procesu       │                   │                  │ obserwator procesu
 ═══════════════════════════════════════════════════════════════════════════
```

---

## 7. Karty komponentów — Przyciski i akcje

### 7.1. Przycisk

| | |
|---|---|
| **Nazwa** | Przycisk (Button) |
| **Klasa bazowa** | `.dn-btn` |
| **Kategoria** | A — Przyciski i akcje |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w oknie, panelu, karcie lub wierszu tabeli, w którym występuje |

**Przeznaczenie.** Wyzwala pojedynczą akcję — zapis, uruchomienie, potwierdzenie, przejście. Jest podstawową jednostką sprawczą interfejsu: wszędzie tam, gdzie użytkownik ma coś zrobić (nie: coś wybrać z listy czy przełączyć stan), właściwym komponentem jest przycisk.

**Forma i waga.** Mały element w linii tekstu, nie samodzielny blok. Domyślnie: padding 9×16 px, czcionka 13 px (`--dn-fs-sm`), grubość `semibold`, promień narożnika 8 px (`--dn-r-sm`), obrys 1 px. Wariant `--sm`: padding 6×11 px, czcionka 12 px. Wariant `--lg`: padding 12×22 px, czcionka 15 px, promień 10 px (`--dn-r-md`). Zawsze pojedynczy wiersz tekstu — przycisk nigdy nie zawija treści do kilku linii.

**Warianty.**

| Wariant | Klasa | Wygląd | Kiedy stosować |
|---|---|---|---|
| Główny | `--glowny` | Wypełnienie granatowe, tekst biały | Akcja podstawowa okna (np. „Zapisz”) |
| Złoty | `--zloty` | Gradient złoty, cień złoty | Jeden, najważniejszy CTA widoku (np. „Uruchom proces”, „Utwórz agenta”, „Zapisz profil”) — nigdy więcej niż jeden na widok |
| Zarys | `--zarys` | Obrys, wypełnienie tłem akcentu przy najechaniu | Akcje drugorzędne (np. „Wczytaj profil”, „Przypisz do…”) |
| Duch | `--duch` | Bez obrysu, przezroczyste tło | Akcje drugorzędne o najniższej wadze wizualnej |
| Sygnalizujący wagę akcji | `--blad` | Obrys i tekst w kolorze błędu | Akcje o dużej wadze skutku (np. usunięcie, przerwanie procesu) — sygnalizacja wizualna, nie mechanizm blokujący |
| Rozmiar mały | `--sm` | Padding 6×11 px | Wiersze tabel, paski narzędzi gęste |
| Rozmiar duży | `--lg` | Padding 12×22 px | Ekrany powitalne, pojedyncza akcja o dużej wadze |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Wygląd wg wariantu (tabela wyżej) |
| Hover | Tło `--dn-hover` (warianty zarys/duch: tło akcentu złotego), lekkie uniesienie dla wariantu złotego (`translateY(-1px)`) |
| Aktywny (wciśnięcie) | `translateY(1px)` — wizualne „wgniecenie” |
| Fokus klawiaturowy | Pierścień złoty 2 px + odsunięcie 2 px, wyłącznie `:focus-visible` (nawigacja klawiaturą) |
| Warunek niespełniony (np. wymagane pole puste) | Przycisk pozostaje w pełni klikalny i zachowuje hover/pierścień fokusu — nigdy krycie 50% z kursorem „niedozwolone”; niegotowość sygnalizowana opisowo (tekst pomocy przy polu, Dymek kontekstowy) obok przycisku, a kliknięcie w tym stanie zwraca komunikat wskazujący brak zamiast wykonania akcji w ciemno |
| Ładowanie | Realizowane przez umieszczenie Wskaźnika ładowania (rozdz. 11.5) wewnątrz przycisku obok lub zamiast etykiety; przycisk zachowuje klikalność przez cały czas trwania akcji — ponowne kliknięcie w trakcie ładowania jest obsługiwane bezpiecznie (idempotentnie) przez logikę akcji, nie blokadą kontrolki |
| Błąd | Nie dotyczy przycisku samego w sobie — błąd akcji zainicjowanej przyciskiem komunikuje się Powiadomieniem (rozdz. 11.2) lub Komunikatem blokowym (rozdz. 11.3), nie zmianą wyglądu przycisku |

**Zachowanie po interakcji.** Kliknięcie wyzwala akcję natychmiast — przycisk nie ma stanu pośredniego poza opisanym wyżej ładowaniem. Zgodnie z zasadą sygnalizacji wagi akcji (System wizualny, rozdz. 8.1) wariant `--blad` oznacza wizualnie wagę akcji nieodwracalnej lub kosztownej — nie blokuje jej wykonania i nie wymaga dodatkowego potwierdzenia jako reguła sztywna; potwierdzenie przez Modal (rozdz. 11.1) jest wzorcem dostępnym, dobieranym do konkretnej akcji.

**Gdzie używany.** Wszędzie — każde okno operacyjne, każdy modal, każda karta, każdy wiersz tabeli z akcją. Przykłady o szczególnym znaczeniu: „Uruchom proces” (złoty CTA, Execution Monitor); „Zapisz profil” / „Wczytaj profil” / „Przypisz do…” (zarys + jeden złoty, okno konfiguracji punktów izolacji, panel prawy); „Utwórz agenta” (złoty, Agent Builder); stopki wszystkich Modali i Kart.

---

### 7.2. Przycisk ikonowy

| | |
|---|---|
| **Nazwa** | Przycisk ikonowy (IconButton) |
| **Klasa bazowa** | `.dn-btn-ikona` |
| **Kategoria** | A — Przyciski i akcje |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w pasku górnym, nagłówku panelu, wierszu tabeli lub karcie sesji |

**Przeznaczenie.** Wyzwala pojedynczą akcję rozpoznawalną po ikonie, bez towarzyszącej etykiety tekstowej — stosowany tam, gdzie znaczenie akcji jest jednoznaczne z ikony (zamknij, usuń, uruchom, zatrzymaj, odśwież) i gdzie etykieta tekstowa zajęłaby zbędną przestrzeń (paski narzędzi, wiersze tabel, nagłówki kart).

**Forma i waga.** Znikoma — samodzielna ikonka, nigdy duży blok. Kwadrat dokładnie 36×36 px, bez wypełnienia tła w spoczynku, ikona wewnątrz 18×18 px. To jest komponent, którego pomylenie z Przyciskiem (rozdz. 7.1) najczęściej prowadzi do błędu projektowego opisanego we wstępie dokumentu — Przycisk ikonowy nie ma obrysu w spoczynku i nie rośnie wraz z długością etykiety, bo etykiety nie ma.

**Warianty.** Brak wariantów wizualnych odrębnych od klasy bazowej — jedna forma, stosowana konsekwentnie w każdym kontekście. Na powierzchni marki (pasek górny) przyjmuje kolor treści na marce (`--dn-na-marce-2`) zamiast standardowego tekstu drugorzędnego.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Przezroczyste tło, ikona w kolorze tekstu drugorzędnego |
| Hover | Tło `--dn-hover`, ikona w kolorze tekstu podstawowego |
| Aktywny (wciśnięcie) | Tło `--dn-hover` + `translateY(1px)` |
| Fokus klawiaturowy | Pierścień złoty 2 px + odsunięcie 2 px |
| Warunek niespełniony (akcja chwilowo bez zastosowania) | Ikona pozostaje w pełni klikalna — nigdy krycie 50% z kursorem „niedozwolone”; niedostępność sygnalizowana Dymkiem kontekstowym (rozdz. 11.4) wyjaśniającym powód, a kliknięcie zwraca komunikat zamiast wykonania akcji w ciemno |
| Ładowanie | Ikona zastępowana Wskaźnikiem ładowania (rozdz. 11.5) na czas trwania akcji (np. `odswiez` w toku odświeżania); przycisk pozostaje klikalny |
| Błąd | Nie dotyczy bezpośrednio — błąd akcji komunikowany Powiadomieniem lub zmianą samej ikony (np. `uruchom` → `blad`) w warstwie ikonografii, nie w warstwie komponentu |

**Zachowanie po interakcji.** Kliknięcie wyzwala akcję natychmiast, analogicznie do Przycisku. W wierszach tabel i na kartach pozycji zwykle występuje w zestawie kilku sztuk obok siebie (np. `olowek` edytuj, `kosz` usuń, `pobierz` eksportuj).

**Gdzie używany.** Pasek górny (wyszukiwanie, powiadomienia, ustawienia); kontrolka zamknięcia karty sesji (✕); akcje sterujące procesem w oknach monitorów — `uruchom` / `zatrzymaj` / `odswiez` (Execution Monitor, Process Monitor); akcje kolejki — enqueue / dequeue / retry (Queue Manager); akcje wiersza karty pozycji — `olowek` / `kosz` / `pobierz` (Sources Manager, Library Explorer, Assets Panel); listwa ustawień strony głównej (ikona + etykieta tekstowa obok, wyjątkowo połączone).

---

## 8. Karty komponentów — Pola formularzy

### 8.1. Pole tekstowe (Input / Textarea)

| | |
|---|---|
| **Nazwa** | Pole tekstowe |
| **Klasa bazowa** | `.dn-input` (jednowierszowe) / `.dn-textarea` (wieloliniowe) |
| **Kategoria** | B — Pola formularzy |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczne bez interakcji w formularzu, panelu ustawień lub stopce Chat Window |

**Przeznaczenie.** Przyjmuje wprowadzany przez użytkownika tekst — nazwę, opis, instrukcję systemową, treść wyszukiwania, wartość liczbową lub adres. Zawsze towarzyszy mu etykieta `.dn-pole-etykieta` w kontenerze `.dn-pole`.

**Forma i waga.** Średnia — pełna szerokość kontenera (nie kontener sam w sobie), wysokość jednowierszowa dla `.dn-input` (padding 9×12 px, czcionka 15 px), zmienna wysokość dla `.dn-textarea` z minimum 88 px, zmienna ręcznie (`resize: vertical`). Promień narożnika 8 px, obrys 1 px.

**Warianty.**

| Wariant | Klasa | Zastosowanie |
|---|---|---|
| Jednowierszowe | `.dn-input` | Nazwa, wyszukiwanie, wartość krótka |
| Wieloliniowe | `.dn-textarea` | Opis, instrukcja systemowa, treść dłuższa |
| Błąd | `--blad` (na obu powyższych) | Wartość niepoprawna lub odrzucona |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Obrys `--dn-obrys-mocny`, tło `--dn-powierzchnia` |
| Hover | Obrys `--dn-tekst-3` |
| Aktywny / w edycji (fokus) | Obrys w kolorze `--dn-focus`, poświata 3 px `--dn-focus-cien` |
| Niedostępne w danym kontekście | Pole pozostaje fokusowalne i klikalne — nigdy krycie 60% z kursorem „niedozwolone”; ograniczenie sygnalizowane opisowo tekstem pomocy `.dn-pomoc` pod polem lub Dymkiem kontekstowym |
| Ładowanie | Nie dotyczy pola samego — pole pozostaje edytowalne; ewentualna walidacja asynchroniczna sygnalizowana osobnym Wskaźnikiem ładowania obok pola |
| Błąd | Obrys w kolorze błędu (`--dn-blad-fg`), poświata fokusu w kolorze tła błędu; tekst pomocy `.dn-pomoc--blad` pod polem — pole pozostaje w pełni edytowalne |
| Placeholder (tekst podpowiedzi) | Kolor `--dn-tekst-3`, widoczny wyłącznie gdy pole puste |

**Zachowanie po interakcji.** Kliknięcie lub nawigacja klawiaturą (Tab) ustawia fokus i umożliwia wpisywanie. Walidacja i przejście w stan błędu następują wg reguł właściwych konkretnemu polu (nigdy jako blokada uniemożliwiająca dalsze wpisywanie — pole pozostaje edytowalne również w stanie błędu). Pod polem może wystąpić `.dn-pomoc` — tekst pomocniczy stały lub `.dn-pomoc--blad` przy błędzie.

**Gdzie używany.** Pole wyszukiwania paska górnego; Agent Builder (nazwa, opis, instrukcje systemowe); Instructions Panel (instrukcje projektu); Source Panel (tekst źródłowy do tłumaczenia); Prompt Builder (treść polecenia generującego); niemal każde okno kreatora/buildera.

---

### 8.2. Lista rozwijana (Select)

| | |
|---|---|
| **Nazwa** | Lista rozwijana |
| **Klasa bazowa** | `.dn-select` |
| **Kategoria** | B — Pola formularzy |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Kliknięcie pola rozwija listę; po wyborze wartości lista zwija się samoczynnie |

**Przeznaczenie.** Wybór jednej wartości z zamkniętego zbioru nazwanych opcji — tam, gdzie liczba opcji jest zbyt duża dla zestawu Zakładek (rozdz. 9.4) lub grupy pól wyboru, a jednocześnie wybór nie jest przełącznikiem dwuwartościowym (rozdz. 8.3).

**Forma i waga.** Średnia — wizualnie identyczna z Polem tekstowym (padding 9×12 px, ta sama wysokość, ten sam obrys), różni się jedynie strzałką rozwijania po prawej stronie (ikona `strzalka-dol` maskowana kolorem tekstu trzeciorzędnego) i brakiem możliwości wpisywania dowolnego tekstu — wyłącznie wybór z listy.

**Warianty.** Brak wariantów wizualnych odrębnych od klasy bazowej; dziedziczy wariant `--blad` po wspólnej regule pól formularzy.

**Stany.** Analogiczne jak w Polu tekstowym (rozdz. 8.1): domyślny, hover (obrys `--dn-tekst-3`), aktywny/rozwinięty (obrys `--dn-focus` + poświata), niedostępne w danym kontekście (sygnalizowane opisowo, nigdy kryciem 60% z zablokowaną interakcją — zob. rozdz. 8.1), błąd (obrys `--dn-blad-fg`). Stan „ładowanie” dotyczy sytuacji, w której lista opcji jest pobierana asynchronicznie — pole pozostaje klikalne, a rozwinięcie pokazuje Wskaźnik ładowania (rozdz. 11.5) do czasu wypełnienia listy opcji.

**Zachowanie po interakcji.** Kliknięcie rozwija natywną listę opcji przeglądarki/webview (kontrolka natywna systemu, stylizowana tylko w warstwie pola zamkniętego — `color-scheme` dopasowuje wygląd natywnej listy do motywu jasnego/ciemnego, System wizualny rozdz. 2.5). Wybór opcji zamyka listę i wypełnia pole wybraną wartością.

**Gdzie używany.** Model Configuration (wybór modelu bazowego agenta); Export Panel (wybór formatu eksportu); okno konfiguracji, panel nawigacji zakresów (koncepcyjnie — wybór jednego z trzynastu zakresów może być zrealizowany jako lista pionowa `.dn-karta--pozycja`, zob. rozdz. 9.2, a nie jako Select — oba wzorce współistnieją w platformie dla różnych gęstości opcji); wybór zasięgu przypisania profilu izolacji.

---

### 8.3. Przełącznik (Toggle)

| | |
|---|---|
| **Nazwa** | Przełącznik |
| **Klasa bazowa** | `.dn-suwak` |
| **Kategoria** | B — Pola formularzy |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Widoczny po otwarciu panelu ustawień, panelu wysuwanego (15.6) lub rozwinięcia kontekstowego |

**Przeznaczenie.** Przełącza pojedynczą właściwość dwuwartościową — włączone/wyłączone, współdzielone/odrębne. Jest najliczniej występującym pojedynczym typem kontrolki w oknie konfiguracji punktów izolacji, gdzie tworzy macierz do jedenastu przełączników jednocześnie (trzy izolacji kontekstu + osiem izolacji technicznej) na jeden poziom zasięgu.

**Forma i waga.** Znikoma — tor 40×22 px, okrągły suwak 18×18 px wewnątrz. Nigdy nie występuje samodzielnie bez etykiety opisującej, co przełącza — zawsze w parze: etykieta pozycji po lewej, przełącznik po prawej, w wierszu Tabeli (rozdz. 10.4) lub Karty pozycji (rozdz. 10.1).

**Warianty.** Brak wariantów klasy — jeden kształt, jeden rozmiar w całej platformie.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Wyłączony (wartość logiczna „off”, stan domyślny) | Tor w kolorze `--dn-obrys-mocny`, suwak po lewej |
| Włączony (wartość logiczna „on”) | Tor w gradiencie złotym, suwak przesunięty w prawo (`translateX(18px)`) |
| Hover | Brak odrębnej reguły wizualnej — kursor wskazuje klikalność |
| Fokus klawiaturowy | Poświata 3 px `--dn-focus-cien` na torze |
| Odziedziczony (wartość przejęta z poziomu szerszego, brak jawnego nadpisania na tym poziomie) | Tor w przytłumionym kryciu (80%) + Plakietka „odziedziczone” obok etykiety; suwak pozostaje w pełni klikalny na każdym poziomie zasięgu — kliknięcie tworzy jawne nadpisanie na tym poziomie i przełącza wartość |
| Ładowanie | Nie dotyczy — zmiana stanu jest natychmiastowa i lokalna; zapis wartości na serwerze sygnalizowany osobno (np. Powiadomieniem po zapisie profilu) |
| Błąd | Nie dotyczy pojedynczego przełącznika — błąd zapisu całej konfiguracji komunikowany na poziomie okna (Komunikat blokowy lub Powiadomienie), nie zmianą wyglądu przełącznika |

Rozróżnienie nazewnicze ważne dla dewelopera: stan wizualny „odziedziczony” (wartość przejęta z poziomu szerszego, bez jawnego nadpisania w danym widoku) **nie** jest realizowany atrybutem `disabled`/`aria-disabled` — przełącznik pozostaje w pełni interaktywny na każdym poziomie zasięgu, nadpisywalny w dowolnym widoku bez wyjątku. Odziedziczenie jest sygnalizowane wyłącznie opisowo (przytłumione krycie toru + Plakietka „odziedziczone” lub Dymek kontekstowy wyjaśniający pochodzenie wartości), nigdy przez odebranie możliwości kliknięcia. Nie należy tego mylić ze stanem logicznym „wyłączony” (`off`) samego przełącznika, który jest jego zwykłą, w pełni interaktywną wartością wyjściową.

**Zachowanie po interakcji.** Kliknięcie w dowolnym miejscu toru odwraca stan logiczny natychmiast, bez opóźnienia i bez potwierdzenia pośredniego — zgodnie z zasadą, że żaden przełącznik macierzy izolacji nie jest zablokowany do edycji ani nie wymusza stanu włączonego (System wizualny, rozdz. 12.2). Stanem wyjściowym wszystkich ośmiu przełączników izolacji technicznej jest „wyłączony”.

**Gdzie używany.** Macierz izolacji, panel środkowy okna konfiguracji punktów izolacji — jedenaście przełączników na poziom zasięgu (trzy izolacji kontekstu, osiem izolacji technicznej), to największe w platformie skupisko tego komponentu; rejestr rozszerzeń, kolumna stanu włączenia (Connectors Manager, Skills Manager) — jeden przełącznik w wierszu tabeli na rozszerzenie; Permissions Center — jeden przełącznik na uprawnienie; Queue Manager, Connectors Manager — właściwość binarna w wierszu.

---

### 8.4. Pole wyboru / opcja jednokrotna (Checkbox / Radio)

| | |
|---|---|
| **Nazwa** | Pole wyboru (Checkbox) / opcja jednokrotna (Radio) |
| **Klasa bazowa** | `.dn-check` |
| **Kategoria** | B — Pola formularzy |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Widoczne po otwarciu formularza, panelu ustawień lub rozwinięcia kontekstowego |

**Przeznaczenie.** Checkbox: zaznaczenie jednej lub wielu opcji niezależnych od siebie. Radio (ta sama klasa, natywny typ `input` odróżnia zachowanie): wybór dokładnie jednej opcji z zamkniętej, wzajemnie wykluczającej się grupy — w platformie stosowany przede wszystkim do wyboru warstwy konfiguracji (domyślna / sesji).

**Forma i waga.** Znikoma — kwadrat lub koło 16×16 px, zawsze z towarzyszącą etykietą tekstową obok, nigdy samodzielnie.

**Warianty.** Checkbox (kwadrat, zaznaczenie niezależne) i Radio (koło, wybór jednokrotny w grupie) dzielą jedną klasę CSS — różnicuje je wyłącznie natywny atrybut `type` elementu `input` oraz wspólna nazwa `name` dla grupy radio.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (niezaznaczone) | Kontur w kolorze obrysu, tło puste |
| Zaznaczone | Wypełnienie w kolorze akcentu złotego (`accent-color: --dn-zloto-500`) |
| Hover | Kursor wskazuje klikalność, brak odrębnej reguły tła |
| Fokus klawiaturowy | Pierścień złoty wspólny dla kontrolek interaktywnych |
| Niedostępne w danym kontekście | Kontrolka pozostaje klikalna — nigdy krycie 60% z kursorem „niedozwolone”; ograniczenie sygnalizowane opisowo etykietą pomocniczą obok lub Dymkiem kontekstowym |
| Ładowanie / błąd | Nie dotyczy — kontrolka natywna bez pośrednich stanów |

**Zachowanie po interakcji.** Checkbox: kliknięcie odwraca stan zaznaczenia niezależnie od innych pól tej samej grupy. Radio: kliknięcie zaznacza daną opcję i automatycznie odznacza pozostałe opcje tej samej grupy `name`.

**Gdzie używany.** Wybór warstwy konfiguracji — domyślna / sesji (okno konfiguracji punktów izolacji, panel prawy, rozdz. 12.3 Systemu wizualnego); zaznaczenia wielokrotne na listach (np. wybór wielu źródeł do eksportu w Sources Manager).

---

## 9. Karty komponentów — Nawigacja

### 9.1. Pasek górny

| | |
|---|---|
| **Nazwa** | Pasek górny (Top bar / Header) |
| **Klasa bazowa** | `.dn-pasek` |
| **Kategoria** | C — Nawigacja |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji na szczycie okna każdego środowiska |

**Przeznaczenie.** Stały, poziomy pasek na szczycie okna środowiska — niesie tożsamość marki, wyszukiwanie globalne i szybki dostęp do funkcji narzędziowych (powiadomienia, ustawienia, awatar użytkownika). Jest wzorcem wspólnym dla górnego paska okna każdego z czterech środowisk (System wizualny, rozdz. 8.7).

**Forma i waga.** Duża — pełna szerokość okna, wysokość stała 56 px (`--dn-topbar-h`). Tło paska jest zawsze ciemne (granatowo-grafitowe), niezależnie od aktywnego motywu jasnego lub ciemnego całej aplikacji — pasek nie jaśnieje w motywie jasnym (zob. Załącznik B, uwaga o tokenie `--dn-pasek-tlo`).

**Warianty.**

| Wariant | Zastosowanie |
|---|---|
| Pełny (wzorzec bazowy) | Górny pasek okna środowiska — logo, wyszukiwanie, ikony, awatar |
| Lekki (listwa ustawień, strefa 3 strony głównej) | Tło `--dn-powierzchnia-2` zamiast tła marki, wysokość ograniczona, elementy jako Przycisk ikonowy z etykietą tekstową obok |

**Stany.** Pasek sam w sobie jest kontenerem statycznym bez własnych stanów interakcji — stany dotyczą elementów wewnątrz niego (Przycisk ikonowy rozdz. 7.2, Pole tekstowe rozdz. 8.1, Awatar rozdz. 12.1), każdy ze swoim własnym zestawem stanów opisanym w odpowiedniej karcie.

**Zachowanie po interakcji.** Elementy paska działają niezależnie: pole wyszukiwania filtruje/wyszukuje w zależności od kontekstu, ikony otwierają powiązane funkcje (powiadomienia, ustawienia), awatar otwiera menu użytkownika lub profil.

**Makieta — anatomia paska górnego.**

```
┌─ .dn-pasek (56 px) ──────────────────────────────────────────────────┐
│  ◈ Danaco   [ 🔍 szukaj…                    ]        🔔  ⚙  ( AW )  │
│  godło+nazwa   pole wyszukiwania (elastyczne, max 520 px)  ikony  awatar│
└────────────────────────────────────────────────────────────────────┘
```

**Gdzie używany.** Górny pasek okna każdego z czterech środowisk (TalkIn, WorkSpace, CodeStudio, MultitaskingAI); wariant lekki — listwa ustawień strefy 3 strony głównej.

---

### 9.2. Boczna nawigacja modułów / Panel orkiestracji

| | |
|---|---|
| **Nazwa** | Boczna nawigacja modułów (środowiska modułowe) / Panel orkiestracji (MultitaskingAI) |
| **Wzorzec** | `.dn-karta--pozycja` w układzie pionowym |
| **Kategoria** | C — Nawigacja |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji jako pierwsza kolumna od lewej; przełączana ikoną ☰ (15.4) w widoku wąskim |

**Przeznaczenie.** Stały, pionowy element interfejsu środowiska, wskazujący dostępne moduły (w TalkIn, WorkSpace, CodeStudio) lub sześć sekcji sterowania zespołem (w MultitaskingAI: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu). Widoczny niezależnie od tego, który moduł jest aktualnie otwarty w bieżącej karcie sesji — środowisko jest profilem widoczności modułów w tej nawigacji, nie pojemnikiem, do którego moduł należy wyłącznie.

**Forma i waga.** Duża — pełna wysokość okna, szerokość stała (węższa niż obszar roboczy), zawsze widoczna, nigdy nakładka. Lista pozycji zbudowana z `.dn-karta--pozycja`: każda pozycja to wiersz z odstępem `--dn-od-2`/`--dn-od-3`, rozdzielony od kolejnej cienką kreską dolną, nie osobną obwódką.

**Warianty.**

| Wariant | Zawartość | Środowisko |
|---|---|---|
| Lista modułów | 8–9 pozycji wg macierzy dostępności modułów | TalkIn, WorkSpace, CodeStudio |
| Panel orkiestracji | 6 sekcji stałych: Zespoły, Role, Kolejki, Orkiestracja, Harmonogram i automatyki, Monitor procesu | MultitaskingAI |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (pozycja nieotwarta) | Tekst w kolorze podstawowym, tło przezroczyste |
| Hover | Tło `--dn-powierzchnia-2` |
| Aktywny / wybrany (moduł lub sekcja otwarta w bieżącej karcie) | Tło `--dn-akcent-tlo`, tekst `--dn-akcent-txt` — analogicznie do stanu aktywnego Zakładki |
| Wyłączony | Nie dotyczy — pozycja spoza macierzy dostępności modułów danego środowiska nie jest wyświetlana wcale, nie jest wyświetlana jako wyłączona |
| Ładowanie | Krótkotrwały stan podczas przeładowania przestrzeni roboczej karty przy zmianie modułu — realizowany Wskaźnikiem ładowania w obszarze głównym, nie w samej nawigacji |
| Błąd | Nie dotyczy |

**Zachowanie po interakcji.** Kliknięcie pozycji modułu (lub sekcji panelu orkiestracji) przeładowuje przestrzeń roboczą bieżącej karty sesji do układu okien właściwego nowo wybranej pozycji — znikają okna poprzedniego modułu, pojawia się zestaw okien nowego. Sama nawigacja pozostaje widoczna w niezmienionej formie; przeładowaniu podlega wyłącznie obszar roboczy po jej prawej stronie. Kolejność i widoczność sekcji panelu orkiestracji są konfigurowalne z okna konfiguracji; zestaw sześciu sekcji jest zestawem domyślnym.

**Gdzie używany.** Każde z czterech środowisk, po lewej stronie obszaru roboczego, sąsiadująco z każdym oknem operacyjnym; w oknie konfiguracji analogiczny wzorzec pionowej listy realizuje nawigację między trzynastoma zakresami ustawień oraz selektor siedmiu poziomów zasięgu w oknie konfiguracji punktów izolacji.

---

### 9.3. Karta sesji

| | |
|---|---|
| **Nazwa** | Karta sesji (Session Card / Tab) |
| **Wzorzec** | `.dn-karta--pozycja` w pasku poziomym |
| **Kategoria** | C — Nawigacja |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji w rzędzie kart sesji pod paskiem górnym |

**Przeznaczenie.** Reprezentuje jedną, samodzielną przestrzeń roboczą — z własnym układem okien, historią i kontekstem. Mechanika analogiczna do kart przeglądarki internetowej: użytkownik otwiera wiele kart obok siebie, pracuje równolegle w kilku modułach lub kilku miejscach tego samego środowiska, przełączając się między nimi jednym kliknięciem.

**Forma i waga.** Mała — pojedynczy element w rzędzie kart sesji bezpośrednio pod paskiem górnym, szerokości dopasowanej do treści tytułu (nie na pełną szerokość, nie stała). Rząd kart mieści się w jednym wierszu; przy przepełnieniu karty zwężają się do szerokości minimalnej, a nadmiar udostępnia przewijanie w osi rzędu.

**Warianty.** Brak nazwanych wariantów CSS — jedna forma, różnicowana wyłącznie stanem (aktywna / w tle) i obecnością wskaźnika pracy w tle.

**Anatomia.**

| Element karty sesji | Zawartość | Komponent nośny |
|---|---|---|
| Tytuł karty | Nazwa modułu otwartego w karcie (np. „Studio”, „Developer”) | Tekst, aktualizuje się przy przeładowaniu |
| Wskaźnik stanu | Sygnalizacja aktywnej pracy w tle (trwający proces automatyki, aktywna rola MultitaskingAI) | Kropka stanu (rozdz. 10.5) |
| Kontrolka zamknięcia | Zamyka kartę sesji | Przycisk ikonowy (rozdz. 7.2), ikona `zamknij` |
| Kontrolka nowej karty | Otwiera nową, pustą kartę sesji w bieżącym środowisku | Przycisk ikonowy, ikona `plus`, na końcu paska |
| Obszar przełączania | Kliknięcie dowolnego miejsca karty poza kontrolką zamknięcia czyni ją kartą aktywną | Cały obszar karty poza kontrolką „✕” |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny / w tle | Tło neutralne, tekst drugorzędny |
| Aktywna (przednia) | Tło wyróżnione, tekst podstawowy — widok główny okna odpowiada zawartości tej karty |
| Hover | Tło `--dn-powierzchnia-2` |
| Wyłączony | Nie dotyczy — karta sesji nie ma stanu wyłączonego, jedynie istnieje lub jest zamknięta |
| Ładowanie | Karta nowo otwarta oczekuje na wybór modułu z bocznej nawigacji — stan pusty, analogiczny do nowej karty przeglądarki bez załadowanej strony |
| Błąd | Nie dotyczy na poziomie samej karty — błąd procesu wewnątrz sesji sygnalizowany wewnątrz obszaru roboczego tej karty (Komunikat blokowy, Powiadomienie), nie zmianą wyglądu karty |

**Zachowanie po interakcji.**

| Operacja | Wyzwalacz | Skutek |
|---|---|---|
| Otwarcie nowej karty | Kontrolka „+” | Nowa, pusta przestrzeń robocza; nawigacja aktywna, oczekuje wyboru modułu |
| Przełączenie karty | Kliknięcie karty poza „✕” | Obszar roboczy przełącza się na układ, historię i kontekst zapisane w tej karcie |
| Zamknięcie karty | Kontrolka „✕” | Kończy widok karty; los leżącej u podstaw sesji jest funkcją modelu sesji i procesów, poza zakresem niniejszego katalogu |

Każda nowa karta sesji otrzymuje domyślnie odrębny układ okien, odrębną historię i odrębny kontekst — współdzielenie między konkretnymi kartami jest ustawieniem konfiguracyjnym okna konfiguracji, nie stanem wyjściowym.

**Makieta — rząd kart sesji.**

```
  ┌────────────────────────┬────────────────────────┬───────┐
  │  ●  Studio          ✕ │     Research         ✕ │   +   │
  └────────────────────────┴────────────────────────┴───────┘
     ▲ karta aktywna            ▲ karta w tle            ▲ nowa karta
     ●  wskaźnik pracy w tle    ✕  kontrolka zamknięcia
```

**Gdzie używany.** Rząd kart pod paskiem górnym w oknie każdego z czterech środowisk — jedyne miejsce wystąpienia, ale obecne w każdej sesji roboczej platformy.

---

### 9.4. Zakładka

| | |
|---|---|
| **Nazwa** | Zakładka (Tab) |
| **Klasa bazowa** | `.dn-zakladki` (grupa) / `.dn-zakladka` (pozycja) |
| **Kategoria** | C — Nawigacja |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji w obrębie okna operacyjnego, którego widoki przełącza |

**Przeznaczenie.** Przełącza widok w obrębie jednego okna operacyjnego, bez zmiany szerszego kontekstu — w odróżnieniu od Karty sesji (rozdz. 9.3), która przełącza całą przestrzeń roboczą. Zakładka pozostaje w tym samym oknie, zmienia wyłącznie prezentowaną zawartość.

**Forma i waga.** Mała — pozycja w poziomym rzędzie, padding 10×16 px (wariant podkreślony) lub 4×10 px (wariant pigułkowy), czcionka 13 px, grubość `semibold`. Rząd zakładek nigdy nie jest wyższy niż pojedynczy wiersz tekstu plus podkreślenie/tło zaznaczenia.

**Warianty.**

| Wariant | Klasa | Charakterystyka | Zastosowanie |
|---|---|---|---|
| Podkreślony | `.dn-zakladki` (bazowy) | Dolna kreska złota 2 px na zakładce aktywnej, reszta bez tła | Przełączanie widoków w oknie operacyjnym (np. panele modeli w Roundtable) |
| Segmentowany (pigułkowy) | `.dn-zakladki--pigulki` | Cienki rząd, zakładka aktywna z subtelnym tłem, bez obwódki grupy | Przełączniki gęstsze — np. wybór warstwy konfiguracji w oknie punktów izolacji (wskazanie rozdz. 8.5 Systemu wizualnego; w makiecie rozdz. 12.3 tego samego dokumentu ten sam wybór przedstawiono również jako parę pól radiowych — rozdz. 8.4 niniejszego katalogu — oba zapisy dotyczą tego samego wyboru dwuwartościowego) |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (nieaktywna) | Tekst `--dn-tekst-3`, brak podkreślenia/tła |
| Hover | Tekst `--dn-tekst` (podstawowy) |
| Aktywna | Tekst `--dn-akcent-txt`, podkreślenie złote 2 px (wariant podkreślony) lub tło `--dn-powierzchnia-2` (wariant pigułkowy) |
| Wciśnięcie (active) | Tekst `--dn-akcent-txt` |
| Fokus klawiaturowy | Pierścień złoty; dla wariantu pigułkowego dodatkowo zachowany cień pozycji aktywnej |
| Niedostępna w danym kontekście | Zakładka pozostaje w pełni klikalna — nigdy krycie 50% z kursorem „niedozwolone”; niedostępność zawartości pod nią sygnalizowana opisowo (Dymek kontekstowy lub komunikat w treści panelu po przełączeniu), nie odebraniem możliwości kliknięcia samej zakładki |
| Ładowanie / błąd | Nie dotyczy zakładki samej — treść panelu pod zakładką może mieć własny stan ładowania/błędu, niezależny od wyglądu zakładki |

**Zachowanie po interakcji.** Kliknięcie zakładki nieaktywnej czyni ją aktywną i podmienia zawartość panelu pod rzędem zakładek; poprzednio aktywna zakładka wraca do stanu domyślnego. Zmiana jest natychmiastowa, bez przeładowania okna.

**Gdzie używany.** Model Panels w module Roundtable (przełączanie widoku odpowiedzi modeli); Terminal Tabs w module Terminal (równoległe sesje PowerShell, CMD, Bash — każda karta to odrębna sesja powłoki); wybór warstwy konfiguracji w oknie punktów izolacji (wariant pigułkowy).

---

### 9.5. Menu kontekstowe

| | |
|---|---|
| **Nazwa** | Menu kontekstowe |
| **Wzorzec** | `.dn-karta` (pojemnik pływający) + `.dn-karta--pozycja` (pozycje) |
| **Kategoria** | C — Nawigacja |
| **Status** | Wzorzec złożony — brak odrębnej klasy w `components.css`; składany z komponentów zdefiniowanych wskazanych niżej |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie wyzwalacza ⋮ (15.3), ☰ (15.4) lub elementu chrome okna operacyjnego |

**Przeznaczenie.** Udostępnia krótką listę akcji lub ustawień powiązanych z konkretnym elementem (wiersz, karta, okno operacyjne), bez opuszczania bieżącego widoku. Występuje w platformie w dwóch odmiennych zastosowaniach o wspólnej formie: menu akcji pozycji (wyzwalane ikoną `wiecej`) oraz uproszczone menu kontekstowe okna operacyjnego, udostępniające szybką zmianę konfiguracji bieżącej sesji bez otwierania pełnego okna konfiguracji.

**Forma i waga.** Mała do średniej — pływająca nakładka zakotwiczona przy wyzwalaczu (nie na środku ekranu jak Modal), szerokość dopasowana do treści najdłuższej pozycji, wysokość wg liczby pozycji. Zbudowana z tych samych elementów co Karta (rozdz. 10.1): tło `--dn-panel`, cień `--dn-cien-3` (nie `--dn-cien-lg` właściwy Modalowi — waga wizualna mniejsza), promień narożnika `--dn-r-md`.

**Warianty.**

| Wariant | Wyzwalacz | Zawartość |
|---|---|---|
| Menu akcji pozycji | Ikona `wiecej` na wierszu tabeli lub karcie | Lista akcji dotyczących tej jednej pozycji (edytuj, usuń, duplikuj, eksportuj) |
| Uproszczone menu kontekstowe okna operacyjnego | Element chrome okna operacyjnego aplikacji | Podzbiór ustawień okna konfiguracji ograniczony do warstwy sesji — szybka zmiana modelu, izolacji lub pamięci bieżącej karty sesji, bez opuszczania trwającej pracy |
| Otwarcie/zamknięcie bocznej nawigacji | Ikona `menu` | Przełącza widoczność bocznej nawigacji / panelu orkiestracji (zastosowanie typowe dla widoku wąskiego) |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (zamknięte) | Niewidoczne, brak śladu w układzie |
| Otwarte | Widoczna pływająca karta z pozycjami, każda pozycja zachowuje się jak Karta interaktywna (rozdz. 10.1: hover = tło `--dn-powierzchnia-2`) |
| Pozycja z niespełnionym warunkiem | Pozycja menu pozostaje klikalna — nigdy krycie 50% z kursorem „niedozwolone”; niedostępność akcji dla bieżącego stanu sygnalizowana opisowo (tekst pomocniczy w pozycji lub Dymek kontekstowy), a kliknięcie zwraca komunikat zamiast wykonania akcji w ciemno |
| Ładowanie | Pozycje menu zależne od danych pobieranych asynchronicznie mogą wyświetlać Wskaźnik ładowania zamiast listy do czasu jej wypełnienia |
| Błąd | Nie dotyczy samego menu — błąd wykonania akcji wybranej z menu komunikowany Powiadomieniem po zamknięciu menu |

**Zachowanie po interakcji.** Kliknięcie wyzwalacza otwiera menu przy nim zakotwiczone. Kliknięcie poza obszarem menu, klawisz Escape lub wybór pozycji zamyka je. Wybór pozycji wyzwala powiązaną akcję i zamyka menu. Menu nigdy nie przesuwa ani nie przeładowuje reszty widoku — jest nakładką lekką, tymczasową.

**Gdzie używany.** Ikona `wiecej` na wierszach Tabeli i na Kartach pozycji w całej platformie (menu akcji); okno operacyjne aplikacji — szybka zmiana konfiguracji sesji (Model konfiguracji, rozdz. 3.1) jako skrót do tego samego modelu konfiguracji, co pełne okno konfiguracji otwierane ze strony głównej.

---

## 10. Karty komponentów — Dane i treść

### 10.1. Karta

| | |
|---|---|
| **Nazwa** | Karta (Card) |
| **Klasa bazowa** | `.dn-karta` |
| **Kategoria** | D — Dane i treść |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji w panelu, siatce strony głównej lub liście pozycji |

**Przeznaczenie.** Grupuje powiązaną treść w jedną, spójną jednostkę wizualną — od pojedynczej pozycji listy po duży, samodzielny punkt wejścia (kartę środowiska na stronie głównej). Jest najczęściej wykorzystywanym pojedynczym komponentem strukturalnym platformy: bazą zarówno dla kart środowisk i kafli komponentów własnych strony głównej, jak i dla każdej pozycji listy w panelach (Sources Manager, Assets Panel, Library Explorer, selektor zasięgu izolacji, sekcje Zespoły/Role panelu orkiestracji).

**Forma i waga.** Zmienna — jedyny komponent katalogu, którego waga wizualna rozciąga się od znikomej (wiersz listy) po dużą (karta środowiska strony głównej, tytuł krojem nagłówkowym 28–36 px). Karta jest **płaska wewnątrz panelu** — nie ma własnego obrysu ani cienia, gdy występuje jako pozycja listy wewnątrz Panelu (rozdz. 10.2); zyskuje obrys, cień i promień narożnika wyłącznie jako samodzielny, duży element (karta środowiska, kafel).

**Warianty.**

| Wariant | Klasa | Charakterystyka | Waga wizualna |
|---|---|---|---|
| Bazowy | `.dn-karta` | Płaska, bez obrysu i cienia | Zależna od kontekstu |
| Interaktywna | `--interaktywna` | Wskazanie stanu najechania tłem (`--dn-powierzchnia-2`), nie uniesieniem | Zależna od kontekstu |
| Akcentowana | `--akcent` | Wstęga górna 3 px w kolorze złota | Podnosi wagę wizualną — rezerwowana dla stanu aktywnego/najechania kart o dużej wadze |
| Pozycja listy | `--pozycja` | Wiersze rozdzielone cienką kreską dolną zamiast obwódki każdej osobno | Mała — element listy |
| Karta środowiska (strefa 1 strony głównej) | `--interaktywna --akcent` łącznie | Duża, tytuł krojem nagłówkowym 28/36 px, godło, opis trybu pracy, cień złoty na stanie aktywnym/najechania | Duża — środek ciężkości strony głównej |
| Kafel komponentu własnego (strefa 2 strony głównej) | `--interaktywna` bez `--akcent` | Mniejsza, etykieta krojem tekstowym `semibold`, bez godła i bez cienia złotego na spoczynku | Średnia — świadomie lżejsza od karty środowiska |

**Struktura wewnętrzna (schemat).**

```
.dn-karta [--interaktywna] [--akcent]
├── .dn-karta-naglowek          (krój nagłówkowy .dn-karta-tytul)
├── .dn-karta-tresc
└── .dn-karta-stopka            (akcje: Przycisk ×N, wyrównane do prawej)
```

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (spoczynek) | Neutralna powierzchnia, bez akcentu złotego (karty środowisk); tło przezroczyste (pozycje listy) |
| Hover | Tło `--dn-powierzchnia-2` (wariant interaktywny); dla karty środowiska dodatkowo akcent złoty na wstędze i cień `--dn-cien-zloto`, uniesienie `translateY(-2px)` |
| Aktywny / kliknięta | Analogicznie do hover — akcent złoty utrzymany do czasu przejścia do docelowego widoku |
| Fokus klawiaturowy | Pierścień złoty + zachowany cień własny wariantu (np. cień pozycji `--dn-cien-1`) |
| Niedostępna w danym kontekście | Karta pozostaje w pełni klikalna i zachowuje reakcję hover — nigdy krycie 50% z zablokowanym hover; niedostępność docelowego widoku sygnalizowana opisowo (Plakietka lub tekst pomocniczy w treści karty, ewentualnie Dymek kontekstowy), kliknięcie zwraca komunikat zamiast prowadzić do martwego, przygaszonego elementu |
| Ładowanie | Nie jest zdefiniowany jako wariant klasy — realizowany przez umieszczenie Wskaźnika ładowania lub stanu szkieletowego (skeleton) w miejscu treści karty do czasu jej wypełnienia |
| Błąd | Nie dotyczy karty samej — treść błędna wewnątrz karty sygnalizowana Komunikatem blokowym lub Plakietką błędu (rozdz. 10.5) wewnątrz jej treści |

**Zachowanie po interakcji.** Karta interaktywna: kliknięcie w dowolnym miejscu (poza ewentualną akcją w stopce) otwiera powiązany widok — przestrzeń roboczą środowiska (karta środowiska), okno konfiguracji komponentu (kafel), szczegóły pozycji (karta pozycji). Karta nieinteraktywna pełni wyłącznie rolę wizualnego grupowania treści, bez własnej akcji kliknięcia.

**Gdzie używany.** Cztery karty środowisk strefy 1 strony głównej (TalkIn, WorkSpace, CodeStudio, MultitaskingAI); cztery kafle komponentów własnych strefy 2 (Automations, Agents, Workspace, Assistant); każda pozycja listy w panelach — Sources Manager, Sources Panel, Assets Panel, Library Explorer, Project Tree, selektor zasięgu okna konfiguracji punktów izolacji; cztery karty ról panelu orkiestracji MultitaskingAI (Executor 1, Executor 2, Coordinator, Executor 3/Validator), w tym zagnieżdżone karty Subagent Network (do 15 pozycji).

---

### 10.2. Panel

| | |
|---|---|
| **Nazwa** | Panel |
| **Wzorzec** | Region budowany z `.dn-karta` / `.dn-karta--pozycja` / `.dn-tabela` / `.dn-pole`, tło `--dn-panel` |
| **Kategoria** | D — Dane i treść |
| **Status** | Wzorzec złożony — nie jest pojedynczą klasą CSS, lecz strukturalnym regionem wewnątrz okna operacyjnego |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji jako region okna operacyjnego; wariant otwierany na żądanie realizuje Panel wysuwany (15.6) |

**Przeznaczenie.** Wydziela w obrębie okna operacyjnego pomocniczy, samodzielny obszar treści lub narzędzi, towarzyszący obszarowi głównemu — listę źródeł, listę ustaleń, panel narzędzi kontekstowych, pamięć projektu. Jest kategorią o największej liczbie nazwanych wystąpień w całej platformie: Tools Panel, Sources Panel, Notes Panel, Assets Panel, Findings Panel, Instructions Panel, Context Memory, Recommendations Panel, Errors Panel, Diff/Grep Panel, Permissions Center, Glossary Manager, Tags & Collections, Deployment Panel, Export Panel i inne — wszystkie zbudowane na tym samym wzorcu strukturalnym.

**Forma i waga.** Średnia do dużej — prostokątny region wewnątrz okna operacyjnego, zajmujący kolumnę sąsiadującą z obszarem głównym (Planszą lub edytorem) po jego prawej stronie, rzadziej sam będący punktem wejścia modułu (np. Diagnostics Center). W przeciwieństwie do Karty (rozdz. 10.1) Panel nie jest pojedynczym elementem powtarzalnym w siatce — jest kontenerem strukturalnym, którego wewnętrzna treść korzysta z innych komponentów katalogu w zależności od charakteru danych.

**Warianty wg zawartości.**

| Rodzaj treści panelu | Komponent nośny wewnątrz | Przykład |
|---|---|---|
| Lista pozycji | Karta pozycji (`--pozycja`) | Sources Manager, Assets Panel, Session Repository |
| Tabela danych | Tabela (rozdz. 10.4) | Errors Panel, Logs Viewer |
| Formularz ustawień | Pola formularzy (rozdz. 8) | Instructions Panel, Model Configuration |
| Blok treści technicznej | Blok kodu (rozdz. 10.6) | Build Output, Output Console |
| Panel bez zawartości | Pusty stan (rozdz. 10.7) | Sources Manager, Findings Panel, Queue Manager przed pierwszą konfiguracją |

**Anatomia (schemat).**

```
┌─ PANEL ────────────────────────────────────────────────
│  Nagłówek panelu (nazwa, ewentualne akcje zbiorcze)
├──────────────────────────────────────────────────────
│  Treść panelu — jeden z wariantów wg tabeli wyżej
│    (przewijana niezależnie od reszty okna operacyjnego,
│     gdy treść przekracza dostępną wysokość)
└──────────────────────────────────────────────────────
```

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (z zawartością) | Lista/tabela/formularz wg wariantu |
| Aktywny (w fokusie pracy) | Nie ma odrębnej reguły — panel nie „aktywuje się” jako całość, aktywne są jego elementy składowe |
| Niedostępny (w danym kontekście) | Nie dotyczy panelu jako całości — panel może zawierać elementy z sygnalizowaną opisowo niedostępnością (np. pola formularza, zob. rozdz. 8.1) bez zmiany własnego wyglądu; elementy te pozostają klikalne, jak każda kontrolka interaktywna katalogu |
| Ładowanie | Zawartość zastąpiona Wskaźnikiem ładowania (środek panelu) do czasu odebrania danych |
| Błąd | Zawartość zastąpiona lub poprzedzona Komunikatem blokowym wariantu `--blad` (rozdz. 11.3) |
| Pusty | Pusty stan (rozdz. 10.7) — ikona, tytuł, krótki opis, wyśrodkowane |

**Zachowanie po interakcji.** Panel sam w sobie jest bierny — zachowanie należy do jego zawartości (kliknięcie pozycji listy, edycja pola, przewinięcie tabeli). Panele monitorujące proces aktualizują treść na żywo kanałem WebSocket, niezależnie od akcji użytkownika.

**Gdzie używany.** Tools Panel, Diff/Grep Panel, Session Repository, Preview Window (Studio); Instructions Panel, Context Memory, Project Library (Workspace); Sources Panel, Notes Panel (Browser); Sources Manager, Findings Panel (Research); Tags & Collections, File Preview, Versioning Panel (Library); Assets Panel (Design); Skills Manager, Connectors Manager, Permissions Center (Agents); Errors Panel, Logs Viewer, Recommendations Panel (Diagnostics); i odpowiedniki w pozostałych modułach — pełne zestawienie w rozdziale 16 niniejszego dokumentu (Macierz komponent × moduł / okno) oraz w Specyfikacji okien operacyjnych.

---

### 10.3. Plansza

| | |
|---|---|
| **Nazwa** | Plansza (Board / Canvas) |
| **Wzorzec** | Obszar główny okna, bez ograniczenia wymiarów przez komponent nośny |
| **Kategoria** | D — Dane i treść |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji jako obszar główny okna operacyjnego |

**Przeznaczenie.** Duża powierzchnia robocza do swobodnego zestawiania i przekształcania treści — w odróżnieniu od Panelu (rozdz. 10.2), który prezentuje treść w stałym układzie listy/tabeli/formularza, Plansza pozwala na dowolne rozmieszczenie elementów przez użytkownika. Nazwanym wystąpieniem wzorca jest Design Board modułu Design.

**Forma i waga.** Duża — zajmuje cały obszar główny okna operacyjnego, bez stałych wymiarów własnych; skaluje się do szerokości kolumny obszaru roboczego, pozostałej po kolumnie bocznej nawigacji, kolumnie Chat Window, kolumnie Execution Loop Window i kolumnach paneli pomocniczych. Elementy wewnątrz Planszy (zasoby graficzne, karty kompozycji) zachowują własne wymiary i pozycję niezależną od układu siatki reszty okna.

**Warianty.** Jedna forma wzorca, stosowana w oknach o charakterze pracy polegającym na swobodnym zestawianiu elementów wizualnych, w odróżnieniu od pracy liniowej właściwej Panelowi.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (pusta) | Pusty stan (rozdz. 10.7) — Plansza bez zestawionych jeszcze zasobów |
| Z zawartością | Zestawione elementy graficzne rozmieszczone swobodnie |
| Aktywny (element zaznaczony na planszy) | Zależne od implementacji elementu — poza zakresem katalogu komponentów bazowych |
| Wyłączona | Nie dotyczy — Plansza nie ma stanu wyłączonego jako całość |
| Ładowanie | Wskaźnik ładowania na czas odbierania zasobu generowanego przez AI (Prompt Builder → Plansza) |
| Błąd | Komunikat blokowy przy niepowodzeniu generowania lub wczytania zasobu |

**Zachowanie po interakcji.** Odbiera zasoby wygenerowane i zgromadzone w powiązanym Panelu (Assets Panel), umożliwia ich zestawienie i dalszą pracę koncepcyjną nad spójną kompozycją wizualną. Zmiana zaakceptowana na Planszy aktualizuje powiązane okno podglądu (Preview Window).

**Gdzie używany.** Design Board, moduł Design — obszar główny, zasilany przez Assets Panel i Prompt Builder.

---

### 10.4. Tabela

| | |
|---|---|
| **Nazwa** | Tabela |
| **Klasa bazowa** | `.dn-tabela` + `.dn-tabela-owijka` |
| **Kategoria** | D — Dane i treść |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji w panelu lub obszarze głównym okna |

**Przeznaczenie.** Prezentuje zbiór jednorodnych pozycji w układzie wierszy i kolumn — jest podstawowym komponentem okien typu monitor i manager: list przebiegów, zdarzeń, procesów, kolejek, uprawnień i rozszerzeń.

**Forma i waga.** Duża — pełna szerokość kontenera, wysokość zależna od liczby wierszy (przewijana niezależnie, gdy przekracza dostępną wysokość panelu). Nagłówek kolumn wersalikowy (12 px, odstęp liter 0,04 em) na tle `--dn-powierzchnia-2`; komórki 11×16 px padding, czcionka 13 px (`--dn-fs-sm`).

**Warianty.**

| Wariant | Klasa | Charakterystyka |
|---|---|---|
| Bazowa | `.dn-tabela` | Separacja wierszy kreską subtelną, bez naprzemiennego tła |
| Naprzemienna | `--paski` | Naprzemienne tło wierszy parzystych (`--dn-powierzchnia-2`) |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Wiersze rozdzielone kreską `--dn-obrys-subtelny` |
| Hover (wiersz) | Podświetlenie tła `--dn-hover` |
| Wiersz z pozycją niedostępną w danym kontekście | Wiersz i jego akcje pozostają w pełni klikalne — nigdy krycie/blokada całego wiersza; niedostępność sygnalizowana opisowo Plakietką statusu (rozdz. 10.5) w kolumnie statusu lub Dymkiem kontekstowym, a kliknięcie akcji niedostępnej zwraca komunikat zamiast wykonania w ciemno |
| Ładowanie | Wiersze zastąpione Wskaźnikiem ładowania lub stanem szkieletowym do czasu odebrania danych |
| Błąd | Wiersz z pozycją błędną oznaczony Plakietką błędu (rozdz. 10.5) w kolumnie statusu, bez zmiany struktury tabeli |
| Pusta (brak wierszy) | Pusty stan (rozdz. 10.7) w miejsce ciała tabeli |

**Zachowanie po interakcji.** Wiersz zwykle zawiera kolumnę akcji zbudowaną z Przycisków ikonowych (rozdz. 7.2) lub Przełącznika (rozdz. 8.3), gdy zarządzana właściwość jest binarna. Tabele okien monitorujących aktualizują zawartość na żywo kanałem WebSocket, bez akcji użytkownika.

**Gdzie używany.** Execution Monitor, Process Monitor, Diagnostics Center, Build Output, Actions Monitor (monitory procesów — wiersz z Plakietką statusu i Kropką); Queue Manager, Orchestrator, Connectors Manager, Permissions Center, Agent Manager, Glossary Manager (zarządcy — wiersz z akcją lub przełącznikiem); macierz izolacji, panel środkowy okna konfiguracji punktów izolacji (etykieta pozycji | przełącznik); sekcja Kolejki i Monitor procesu panelu orkiestracji MultitaskingAI.

---

### 10.5. Pigułka i plakietka statusu

| | |
|---|---|
| **Nazwa** | Pigułka / Plakietka (Pill / Tag / Badge) i Kropka stanu |
| **Klasa bazowa** | `.dn-plakietka` / `.dn-kropka` |
| **Kategoria** | D — Dane i treść |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji przy jednostce treści, której status opisuje |

**Przeznaczenie.** Oznacza krótką, samodzielną informację o jednostce treści — status, kategorię, etykietę własną nadaną przez użytkownika. Kropka stanu jest jej minimalną, bezetykietową odmianą, stosowaną tam, gdzie kontekst (np. nagłówek kolumny tabeli) już nazywa znaczenie koloru.

**Forma i waga.** Znikoma — pigułka: padding 2×9 px, promień narożnika pełny (999 px — kształt tabletki), czcionka 12 px, zawsze w linii tekstu, nigdy w osobnym bloku. Kropka stanu: koło 9×9 px (7×7 px w wariancie połączonym z pigułką), z poświatą 3 px w kolorze statusu.

**Warianty.**

| Wariant | Klasa | Zastosowanie |
|---|---|---|
| Neutralna | `.dn-plakietka` (bazowa) | Etykieta bez znaczenia statusowego — tag, kategoria |
| Złota | `--zloto` | Wyróżnienie, oznaczenie specjalne |
| Sukces | `--sukces` | Status pozytywny — zakończone powodzeniem |
| Ostrzeżenie | `--ostrz` | Status wymagający uwagi |
| Błąd | `--blad` | Status negatywny — niepowodzenie |
| Informacyjna | `--info` | Status neutralny informacyjny |
| Wersalikowa | `--wersaliki` | Etykieta kapitalikowa, czcionka 9,5 px, dla oznaczeń najmniejszej wagi |
| Połączona ze stanem | `--stan` | Pigułka + Kropka stanu wewnątrz — pełne oznaczenie statusu z kolorem i tekstem jednocześnie |
| Kropka samodzielna | `.dn-kropka` + `--sukces`/`--ostrz`/`--blad` | Oznaczenie statusu bez etykiety tekstowej, np. w nagłówku już nazwanej kolumny |

**Stany.** Pigułka i Kropka są same w sobie nośnikiem stanu innej jednostki (przebiegu, procesu, zasobu) — nie mają własnych stanów interakcji w rozumieniu hover/aktywny/wyłączony, ponieważ nie są klikalne. Jedyny rozróżniany „stan” to wartość semantyczna wskazana wariantem (sukces/ostrzeżenie/błąd/informacja/neutralna). Zgodnie z zasadą dostępności statusy nigdy nie polegają wyłącznie na kolorze — kolor jest zawsze połączony z ikoną lub tekstem etykiety.

**Zachowanie po interakcji.** Domyślnie bierna (informacyjna); w zastosowaniach typu tag może być klikalna do filtrowania listy po tej wartości — w takim wypadku przejmuje stany interakcji Karty interaktywnej (hover: tło `--dn-powierzchnia-2`).

**Gdzie używany.** Tags & Collections, moduł Library (etykiety i kolekcje nadawane zasobom — dosłowne zastosowanie jako tag); kolumna statusu w Execution Monitor, Process Monitor, Monitor procesu MultitaskingAI (wariant `--stan`); status zgodności w sekcji Orkiestracja panelu orkiestracji; rejestr rozszerzeń — źródło rozszerzenia (Danaco Plugin / Personal) jako plakietka neutralna; nagłówki macierzy izolacji.

---

### 10.6. Blok kodu

| | |
|---|---|
| **Nazwa** | Blok kodu |
| **Klasa bazowa** | `.dn-kod` |
| **Kategoria** | D — Dane i treść |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w treści panelu, Chat Window lub okna monitorującego |

**Przeznaczenie.** Prezentuje treść techniczną — polecenie, ścieżkę pliku, fragment kodu, ładunek narzędzia, log — krojem monospace na wyodrębnionej powierzchni. Odróżnia się od zwykłego tekstu powierzchnią (`--dn-powierzchnia-2`) i krojem, celowo bez własnej ramki ani promienia narożnika typowego dla Karty — w przeciwnym razie w oknach osadzających blok kodu wewnątrz już obramowanego panelu powstawałaby rama w ramie.

**Forma i waga.** Średnia — blok pełnej szerokości kontenera, wysokość zależna od treści (`white-space: pre-wrap`, zawijanie długich linii, przewijanie poziome dla treści nie dających się zawinąć). Czcionka mono 12 px (`--dn-fs-xs`), interlinia bazowa.

**Warianty.** Brak nazwanych wariantów CSS — jedna forma, stosowana identycznie niezależnie od rodzaju treści technicznej.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Tło `--dn-powierzchnia-2`, tekst `--dn-tekst` krojem mono |
| Aktualizowany na żywo | Treść narasta w czasie rzeczywistym (logi, wyjście budowania) — bez odrębnej reguły wizualnej poza samą zmianą treści |
| Wyłączony | Nie dotyczy |
| Ładowanie | Blok pusty lub z komunikatem oczekiwania do czasu napłynięcia pierwszej treści |
| Błąd | Fragment treści w kolorze błędu (`--dn-blad-fg`) w obrębie bloku — stosowane wybiórczo do linii błędu w logu, nie do całego bloku |

**Zachowanie po interakcji.** Zwykle nieedytowalny podgląd (read-only); może udostępniać akcję kopiowania treści. W oknach monitorujących przewija się automatycznie do najnowszej linii, chyba że użytkownik przewinął ręcznie wyżej.

**Gdzie używany.** Terminal — Output Console; Developer — cytowania kodu w Chat Window; Diagnostics — Logs Viewer; podgląd polityki efektywnej w oknie konfiguracji punktów izolacji, panel prawy; Build Output.

---

### 10.7. Pusty stan

| | |
|---|---|
| **Nazwa** | Pusty stan (Empty state) |
| **Klasa bazowa** | `.dn-pusty-stan` |
| **Kategoria** | D — Dane i treść |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w miejscu treści panelu, tabeli lub planszy pozbawionych pozycji |

**Przeznaczenie.** Zastępuje treść Panelu, Tabeli lub Planszy, gdy nie zawierają jeszcze żadnej pozycji — komunikuje, że stan jest oczekiwany i naturalny (pierwsze użycie), a nie błędem ładowania.

**Forma i waga.** Średnia — wypełnia dostępną przestrzeń pustego panelu lub okna, treść wyśrodkowana pionowo i poziomo: ikona, tytuł (13 px, `semibold`), krótki opis (12 px, maksymalna szerokość 42 znaków na linię).

**Warianty.** Brak nazwanych wariantów CSS — jedna forma, różnicowana wyłącznie treścią (ikoną, tytułem, opisem) właściwą kontekstowi.

**Stany.** Sam pusty stan jest stanem innego komponentu (Panelu, Tabeli, Planszy) — nie ma własnych stanów interakcji poza ewentualną osadzoną akcją (np. Przycisk „Dodaj pierwszą pozycję”), która przejmuje stany Przycisku (rozdz. 7.1).

**Zachowanie po interakcji.** Bierny widok informacyjny; może zawierać pojedynczy Przycisk zachęcający do pierwszej akcji (np. dodania źródła, utworzenia zadania).

**Gdzie używany.** Sources Manager, Findings Panel, Queue Manager przed pierwszą konfiguracją (wskazane wprost w Systemie wizualnym, rozdz. 8.12); dowolny inny Panel lub Tabela przed pierwszym wypełnieniem danymi — Library Explorer, Assets Panel, Errors Panel (gdy brak zarejestrowanych błędów).

---

## 11. Karty komponentów — Informacja zwrotna i nakładki

### 11.1. Modal

| | |
|---|---|
| **Nazwa** | Modal (okno modalne) |
| **Klasa bazowa** | `.dn-nakladka` (tło) + `.dn-modal` (okno) |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Akcja wywołująca — kliknięcie kafla, przycisku akcji o dużej wadze lub pozycji menu |

**Przeznaczenie.** Przerywa bieżący widok nakładką wymagającą uwagi lub decyzji użytkownika przed powrotem do pracy — konfiguracja uruchamiana z okna operacyjnego, potwierdzenie akcji o dużej wadze, kreator tworzenia komponentu własnego.

**Forma i waga.** Duża — nakładka pełnoekranowa (`.dn-nakladka`, tło półprzezroczyste `rgba` przyciemniające resztę widoku) z wyśrodkowanym oknem o maksymalnej szerokości 560 px i maksymalnej wysokości 92% wysokości ekranu, promień narożnika 14 px (`--dn-r-lg`), cień `--dn-cien-lg` — największy cień w systemie, zarezerwowany dla elementu o najwyższej wadze wizualnej spośród nakładek.

**Warianty.** Brak nazwanych wariantów CSS — jedna forma okna; różnicowana wyłącznie treścią (formularz konfiguracji, potwierdzenie, kreator wielokrokowy zbudowany z Zakładek lub kroków pionowych).

**Struktura wewnętrzna (schemat).**

```
.dn-nakladka (pełny ekran, tło przyciemnione)
  └── .dn-modal (max 560 px, max 92 vh)
        ├── .dn-modal-naglowek     (tytuł krojem nagłówkowym + kontrolka zamknięcia)
        ├── .dn-modal-tresc        (przewijana niezależnie)
        └── .dn-modal-stopka       (akcje: Przycisk ×N, wyrównane do prawej)
```

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (zamknięty) | Niewidoczny, brak śladu w DOM/układzie |
| Otwierający się | Animacja wjazdu — przezroczystość 0→1 i przesunięcie pionowe 12 px→0 w czasie 0,24 s (`--dn-czas-3`), pomijana przy `prefers-reduced-motion` |
| Otwarty | W pełni widoczny, przechwytuje fokus klawiatury |
| Próba potwierdzenia z brakami w polach | Przycisk akcji głównej w stopce pozostaje aktywny i klikalny niezależnie od stanu wypełnienia pól (zgodnie z rozdz. 7.1) — kliknięcie przy niespełnionych wymaganiach ujawnia braki komunikatem (Komunikat blokowy `--blad`, rozdz. 11.3, lub tekst pomocy `.dn-pomoc--blad` przy polu) zamiast zablokowania przycisku |
| Ładowanie | Treść modala zastąpiona Wskaźnikiem ładowania na czas przygotowania danych (np. podgląd przed potwierdzeniem) |
| Błąd | Komunikat blokowy wariantu `--blad` (rozdz. 11.3) osadzony w treści modala, nad stopką z akcjami |

**Zachowanie po interakcji.** Otwiera się w reakcji na akcję wywołującą (kliknięcie kafla, przycisku „Usuń”, wejście w kreator). Zamyka się przez kontrolkę zamknięcia w nagłówku, klawisz Escape, kliknięcie w nakładkę poza oknem lub przez akcję w stopce (np. „Anuluj”/„Zapisz”). Nakładka blokuje interakcję z resztą widoku wyłącznie na czas otwarcia — nie jest to blokada trwała, użytkownik zawsze ma dostępną drogę zamknięcia.

**Gdzie używany.** Okno konfiguracji uruchamiane z okna operacyjnego; potwierdzenia akcji o dużej wadze; kreatory tworzenia komponentu własnego po kliknięciu kafla strefy 2 strony głównej (Automations, Agents, Workspace, Assistant).

---

### 11.2. Powiadomienie (Toast)

| | |
|---|---|
| **Nazwa** | Powiadomienie (Toast) |
| **Klasa bazowa** | `.dn-toast` (w kontenerze `.dn-toasty`) |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Zdarzenie systemowe — pojawia się samoczynnie po zakończeniu operacji, bez interakcji użytkownika |

**Przeznaczenie.** Potwierdza krótko wynik akcji, która już się dokonała — zapis, zakończenie procesu, błąd operacji w tle — bez przerywania bieżącej pracy i bez wymagania reakcji użytkownika.

**Forma i waga.** Mała — pojedynczy pasek zakotwiczony przy prawej krawędzi obszaru roboczego, szerokość 260–440 px, padding 12×16 px. Tło zawsze ciemne (granatowo-grafitowe), niezależnie od motywu aplikacji — analogicznie do Paska górnego (rozdz. 9.1), ponieważ jest to powierzchnia własna komponentu, nie tło strony.

**Warianty.**

| Wariant | Klasa | Kreska lewa (3 px) |
|---|---|---|
| Neutralny | `.dn-toast` (bazowy) | Kolor akcentu złotego |
| Sukces | `--sukces` | Kolor sukcesu (wariant „na ciemnym”, `-dk`) |
| Błąd | `--blad` | Kolor błędu (wariant „na ciemnym”) |
| Ostrzeżenie | `--ostrz` | Kolor ostrzeżenia (wariant „na ciemnym”) |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Pojawienie się | Animacja wjazdu identyczna z Modalem (0,24 s), stos kilku powiadomień układa się pionowo przy prawej krawędzi obszaru roboczego |
| Widoczny | Pełna krycie, czytelny przez czas ekspozycji |
| Znikanie | Automatyczne po upływie czasu ekspozycji (parametr sterowany z okna konfiguracji) lub ręczne przez kontrolkę zamknięcia |
| Hover (nad kontrolką zamknięcia) | Kontrolka zamknięcia przechodzi z `--dn-na-marce-3` na pełną krycie `--dn-na-marce` |
| Wyłączony / ładowanie / błąd formularza | Nie dotyczy — Powiadomienie samo w sobie nie ma tych stanów, jego wariant `--blad` komunikuje błąd innej operacji, nie własny |

**Zachowanie po interakcji.** Pojawia się samoczynnie w reakcji na zdarzenie systemowe (zapis, zakończenie procesu). Znika samoczynnie po czasie ekspozycji lub natychmiast po kliknięciu kontrolki zamknięcia. Nie blokuje interakcji z resztą widoku — jest wyłącznie nakładką informacyjną, nigdy wymogiem reakcji.

**Gdzie używany.** Potwierdzenie zapisu profilu izolacji (System wizualny, rozdz. 12.3, krok 4 przykładowego przepływu); zakończenie procesu automatyki; komunikaty funkcji globalnej Mobile.

---

### 11.3. Komunikat blokowy (Alert)

| | |
|---|---|
| **Nazwa** | Komunikat blokowy (Alert / inline banner) |
| **Klasa bazowa** | `.dn-alert` |
| **Kategoria** | E — Informacja zwrotna i nakładki (osadzona, nie pływająca) |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w układzie okna przez cały czas trwania sytuacji, którą opisuje |

**Przeznaczenie.** Komunikuje informację, ostrzeżenie lub błąd osadzony bezpośrednio w układzie okna — w odróżnieniu od Powiadomienia (rozdz. 11.2), które pływa nad treścią i znika samoczynnie, Komunikat blokowy zajmuje trwałe miejsce w układzie, dopóki sytuacja, którą opisuje, nie ustanie. W implementacji odpowiada rozpoznanemu wzorcowi wywoływanemu z dziewięciu różnych paneli platformy (m.in. niedostępność Subagentów, niepowodzenie wczytania pliku, potwierdzenie w Terminalu, błąd kroku planu).

**Forma i waga.** Mała do średniej — pełna szerokość rodzica, wysokość dopasowana do treści (zwykle jedna do trzech linii tekstu), oznaczony wyłącznie cienką kreską po lewej stronie (2 px) w kolorze semantycznym — celowo bez własnego tła i bez obwódki, aby nie tworzyć „drugiej karty” wewnątrz już obramowanego panelu.

**Warianty.**

| Wariant | Klasa | Kolor kreski |
|---|---|---|
| Informacyjny | `--info` | `--dn-info-fg` |
| Sukces | `--sukces` | `--dn-sukces-fg` |
| Ostrzeżenie | `--ostrz` | `--dn-ostrz-fg` |
| Błąd | `--blad` | `--dn-blad-fg` |

**Stany.** Komunikat blokowy jest sam w sobie reprezentacją stanu innej sytuacji (dostępności zasobu, wyniku operacji) — nie ma własnych stanów hover/aktywny/wyłączony, ponieważ standardowo nie jest klikalny. Jego „stan” to wybór wariantu semantycznego wg tabeli wyżej. Może zawierać osadzoną akcję (Przycisk `--zarys` lub `--duch`) — w takim wypadku ten Przycisk przejmuje własne stany interakcji.

**Zachowanie po interakcji.** Pojawia się i znika wraz ze stanem, który opisuje — nie ma własnej animacji wejścia/wyjścia odrębnej od reszty układu. Jeśli zawiera akcję (np. „Spróbuj ponownie”), kliknięcie wyzwala tę akcję bez usuwania samego komunikatu, dopóki sytuacja nie ustanie.

**Gdzie używany.** Dziewięć paneli zgłaszających stan operacyjny — między innymi niedostępność mechanizmu Subagent Network, niepowodzenie wczytania pliku, potwierdzenie wklejenia w Terminalu, błąd kroku planu w oknach roboczych ról MultitaskingAI; ogólnie wszędzie tam, gdzie stan błędu lub ostrzeżenia musi pozostać widoczny w układzie, a nie zniknąć samoczynnie jak Powiadomienie.

---

### 11.4. Dymek kontekstowy (Tooltip)

| | |
|---|---|
| **Nazwa** | Dymek kontekstowy (Tooltip) |
| **Klasa bazowa** | `.dn-tooltip` |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Najechanie kursorem lub fokus klawiaturowy na wyzwalaczu, w tym na znaczniku `[?]` |

**Przeznaczenie.** Wyjaśnia działanie i wpływ na aplikację pojedynczego elementu interfejsu, bez opuszczania bieżącego widoku i bez trwałego zajmowania miejsca w układzie. Jest bezpośrednim nośnikiem zasady „każdy element konfiguracji zawiera objaśnienie kontekstowe” (Architektura, rozdz. 13) — oznaczenie `[?]` przywołane w niemal każdym dokumencie źródłowym opisuje właśnie ten komponent.

**Forma i waga.** Znikoma — mały dymek nad wyzwalaczem (domyślnie), padding 6×10 px, czcionka 11 px, tło zawsze ciemne (granatowo-grafitowe) niezależnie od motywu, ogonek trójkątny wskazujący wyzwalacz, cień `--dn-cien-3`.

**Warianty.** Brak nazwanych wariantów CSS — jedna forma; pozycja dymka względem wyzwalacza dopasowuje się do dostępnej przestrzeni, bez zmiany klasy.

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (ukryty) | Krycie 0, nieklikalny (`pointer-events: none`), nieobecny wizualnie |
| Widoczny (hover lub fokus wyzwalacza) | Krycie 1, przesunięcie z 4 px do 0 — wjazd delikatny |
| Nie dotyczy (element niedostępny w kontekście) | Dymek towarzyszy elementowi niezależnie od tego, czy jego akcja jest chwilowo niedostępna w danym kontekście — to właśnie wtedy jest najużyteczniejszy: wyjaśnia, dlaczego akcja jest niedostępna, skoro sama kontrolka pozostaje klikalna (rozdz. 15) |
| Ładowanie / błąd | Nie dotyczy — treść dymka jest statyczna, część definicji ustawienia, nie danymi pobieranymi asynchronicznie |

**Zachowanie po interakcji.** Pojawia się przy najechaniu myszą lub przy fokusie klawiatury na wyzwalaczu (`:hover` lub `:focus-within`), znika natychmiast po ich utracie. Nie wymaga kliknięcia ani osobnego zamknięcia.

**Gdzie używany.** Oznaczenie `[?]` przy każdym elemencie konfiguracji w oknie konfiguracji, ze szczególnym nasileniem w macierzy izolacji okna konfiguracji punktów izolacji, gdzie liczba kontrolek interaktywnych opatrzonych objaśnieniem jest największa w całej platformie.

---

### 11.5. Wskaźnik ładowania (Loader / Spinner)

| | |
|---|---|
| **Nazwa** | Wskaźnik ładowania (Loader / Spinner) |
| **Klasa bazowa** | `.dn-spinner` |
| **Kategoria** | E — Informacja zwrotna i nakładki (osadzony) |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji na czas trwania operacji, którą sygnalizuje |

**Przeznaczenie.** Sygnalizuje trwającą operację, której czas trwania nie jest z góry znany — oczekiwanie na odpowiedź modelu, zapis danych, budowanie, odświeżanie listy. Jest komponentem najmniejszym wagowo w całym katalogu, projektowanym do osadzenia wewnątrz innych komponentów (Przycisku, Przycisku ikonowego, Panelu), nie do samodzielnego występowania na pełnym ekranie.

**Forma i waga.** Znikoma — koło 16×16 px, obrys 2 px w bieżącym kolorze tekstu, jeden bok przezroczysty, obrót ciągły w tempie 0,7 s na pełny obrót.

**Warianty.** Brak nazwanych wariantów CSS — jeden rozmiar, jedna forma; kolor dziedziczony z `currentColor` kontekstu, w którym jest osadzony (np. biały wewnątrz Przycisku głównego, granatowy w tekście podstawowym).

**Stany.** Sam Wskaźnik ładowania jest reprezentacją stanu „ładowanie” innego komponentu — nie ma własnych stanów poza obecnością/nieobecnością. Respektuje `prefers-reduced-motion: reduce` — przy tej preferencji systemowej animacja obrotu skraca się do 0,01 ms, a wskaźnik pozostaje widoczny statycznie zamiast wirować w nieskończoność.

**Zachowanie po interakcji.** Bierny — nie reaguje na kliknięcie ani najechanie. Pojawia się na czas trwania operacji asynchronicznej i znika po jej zakończeniu, ustępując miejsca rezultatowi (treści, Powiadomieniu o wyniku) lub Komunikatowi blokowemu przy niepowodzeniu.

**Gdzie używany.** Wewnątrz Chat Window podczas odbierania strumienia odpowiedzi modelu; wewnątrz Przycisku w trakcie trwania wyzwolonej nim akcji (np. „Zapisz” w trakcie zapisu); Build Output podczas budowania; Execution Monitor podczas przebiegu automatyki; dowolny Panel lub Tabela w trakcie pobierania danych.

### 11.6. Centrum powiadomień

| | |
|---|---|
| **Nazwa** | Centrum powiadomień |
| **Wzorzec** | `.dn-panel-wysuwany` (Panel wysuwany, 15.6) + `.dn-karta--pozycja` (pozycje zdarzeń) + `.dn-plakietka--stan` (wyzwalacz z licznikiem) |
| **Kategoria** | E — Informacja zwrotna i nakładki |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie plakietki powiadomień z licznikiem w pasku kontekstu; zamknięcie kontrolką zamknięcia, klawiszem Escape albo kliknięciem poza kolumną |

**Przeznaczenie.** Jest jedynym mechanizmem powiadamiania platformy: gromadzi w jednym miejscu wszystkie zdarzenia wymagające wiedzy albo decyzji użytkownika i pozwala je obsłużyć bez opuszczania bieżącego widoku. Powiadomienie (Toast, 11.2) jest ulotną prezentacją pojedynczego zdarzenia w chwili jego wystąpienia; centrum powiadomień jest trwałym rejestrem tych samych zdarzeń — każde zdarzenie prezentowane jako Toast trafia jednocześnie do centrum i pozostaje w nim do chwili obsłużenia. Obowiązuje w każdym środowisku, w każdym module i w każdej funkcji globalnej w tej samej formie i pod tym samym wyzwalaczem.

**Taksonomia klas zdarzeń.**

| Klasa zdarzenia | Zakres | Waga |
|---|---|---|
| Zakończenie | Zakończenie procesu, zadania, przebiegu pętli wykonawczej, budowy, eksportu | Normalna |
| Decyzja | Krok procesu oczekujący na zatwierdzenie użytkownika, punkt decyzyjny pętli wykonawczej | Wymagająca decyzji |
| Błąd | Niepowodzenie zadania, naruszenie zależności orkiestracji, powtarzające się niepowodzenie walidacji | Wymagająca decyzji |
| Wzmianka | Odwołanie do użytkownika w narzędziu współpracy, komentarz przypisany do artefaktu | Normalna |
| Termin | Termin zadania projektu, zbliżający się cykl harmonogramu | Informacyjna |
| Automatyka | Wpięcie i przebieg automatyki, wynik cyklu harmonogramu | Informacyjna |
| System | Stan połączenia, wygaśnięcie poświadczenia urządzenia, zakończenie synchronizacji | Informacyjna |

**Postać wizualna i umiejscowienie w układzie pionowym.** Kolumna boczna, otwierana jako rozszerzenie boczne — kolejna kolumna po prawej stronie obszaru roboczego, o pełnej wysokości obszaru roboczego, szerokość 320–420 px, regulowana wyłącznie na szerokość. Po zamknięciu kolumna znika całkowicie z przestrzeni roboczej, a zwolniona szerokość przypada kolumnom pozostałym. Wyzwalaczem jest plakietka powiadomień (`.dn-plakietka--stan`) z licznikiem zdarzeń nieodczytanych, umieszczona w pasku kontekstu — element warstwy 1, jedyna spoczynkowa reprezentacja mechanizmu.

**Elementy składowe.**

| Element | Rola |
|---|---|
| Plakietka powiadomień z licznikiem | Wyzwalacz w pasku kontekstu; liczba zdarzeń w stanie „nowe" |
| Filtr klasy i wagi | Zawężenie listy do wybranych klas zdarzeń albo do zdarzeń wymagających decyzji |
| Filtr źródła | Zawężenie listy do wskazanego środowiska, procesu, projektu albo karty sesji |
| Pozycja zdarzenia | Klasa, treść, źródło, znacznik czasu i zestaw działań |
| Działanie zbiorcze | Oznaczenie wszystkich jako odczytane, odłożenie wszystkich informacyjnych |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zamknięte (spoczynek) | Kolumna niewidoczna, bez śladu w układzie; w pasku kontekstu widoczna wyłącznie plakietka |
| Plakietka bez zdarzeń nowych | Plakietka neutralna, bez licznika |
| Plakietka ze zdarzeniami nowymi | Plakietka `--stan` z licznikiem; kropka w kolorze najwyższej wagi zdarzenia na liście |
| Otwarte, lista zapełniona | Kolumna z listą zdarzeń w porządku od najnowszego, pogrupowaną wg klasy |
| Otwarte, lista pusta | Pusty stan (10.7) z opisem wskazującym Chat Window (14.1) jako kanał zapytania o stan procesów |
| Otwarte, lista filtrowana | Lista zawężona; aktywny filtr widoczny jako znacznik kontekstowy (15.7) w nagłówku kolumny |
| Zdarzenie wymagające decyzji | Pozycja wyróżniona kreską w kolorze uwagi; pozostaje na liście do chwili podjęcia decyzji |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu listy do czasu odebrania rejestru zdarzeń |
| Błąd | Komunikat blokowy (11.3) wariantu `--blad` w miejscu listy, z akcją ponowienia |

**Działania dostępne z poziomu powiadomienia.**

| Działanie | Skutek |
|---|---|
| Przejście do źródła | Otwiera kartę sesji, okno operacyjne albo Execution Loop Window z kontekstem zdarzenia |
| Zatwierdzenie | Zatwierdza krok procesu oczekujący na decyzję, bez opuszczania kolumny |
| Wstrzymanie i wznowienie | Steruje przebiegiem pętli wykonawczej, której zdarzenie dotyczy |
| Ponowienie | Uruchamia ponownie zadanie zakończone niepowodzeniem |
| Odłożenie | Przenosi zdarzenie do stanu odłożonego i przywraca je po wskazanym czasie |
| Oznaczenie jako odczytane | Zmienia stan zdarzenia bez podejmowania działania |

**Relacja do powiadomień wypychanych funkcji Mobile.** Centrum powiadomień jest źródłem zdarzeń dla funkcji globalnej Mobile. Zdarzenie o wadze „wymagająca decyzji" jest wypychane na sparowane urządzenie mobilne jako powiadomienie systemowe; obsługa zdarzenia na urządzeniu i obsługa w kolumnie centrum prowadzą do tego samego stanu encji `powiadomienie` (rozdz. 18.4). Zakres klas wypychanych na urządzenie jest ustawieniem konfiguracyjnym sekcji powiadomień okna Ustawień (rozdz. 7).

**Relacja do sugestii Always On Display.** Sugestia Always On Display nie jest zdarzeniem centrum powiadomień: centrum rejestruje zdarzenia, które zaszły, Always On Display proponuje działania, które użytkownik może podjąć. Plakietka powiadomień na awatarze Always On Display prezentuje ten sam licznik zdarzeń nowych co plakietka paska kontekstu i otwiera tę samą kolumnę centrum powiadomień.

**Gdzie używany.** Powłoka każdego środowiska platformy — TalkIn, WorkSpace, CodeStudio, MultitaskingAI; okna operacyjne wszystkich modułów; funkcja globalna Mobile (widok listy zdarzeń urządzenia); funkcja globalna Always On Display (plakietka na awatarze).

---

## 12. Karty komponentów — Tożsamość

### 12.1. Awatar

| | |
|---|---|
| **Nazwa** | Awatar (Avatar) |
| **Klasa bazowa** | `.dn-awatar` |
| **Kategoria** | F — Tożsamość |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji w pasku górnym, Chat Window, karcie roli lub wierszu tabeli |

**Przeznaczenie.** Identyfikuje wizualnie osobę, model AI, rolę lub agenta w miejscach, gdzie rozróżnienie „kto” jest istotne dla zrozumienia treści — uczestnika rozmowy w Chat Window, model w panelu wielomodelowym, agenta przypisanego do roli.

**Forma i waga.** Znikoma do małej — koło (domyślnie) lub kwadrat (wariant `--kwadrat`) o boku/średnicy 36 px domyślnie, 28 px w wariancie małym, 48 px w wariancie dużym. Wnętrze: inicjał lub obraz (`object-fit: cover`, wypełnia całą powierzchnię, przycięty do kształtu).

**Warianty.**

| Wariant | Klasa | Zastosowanie |
|---|---|---|
| Bazowy | `.dn-awatar` | Gradient granatowy — tożsamość ogólna (użytkownik, model bazowy bez przypisanego agenta) |
| Złoty | `--zloto` | Gradient złoty — wyróżnienie szczególne |
| Mały | `--sm` | 28×28 px — konteksty gęste (wiersz listy, wiadomość w wątku) |
| Duży | `--lg` | 48×48 px — konteksty eksponowane (nagłówek profilu) |
| Kwadratowy | `--kwadrat` | Identyfikacja nie-osobowa — agent (odróżnienie od modelu bazowego, który zachowuje kształt kołowy) |
| Ze wskaźnikiem statusu | `-stan` (modyfikator) | Kropka online w prawym dolnym rogu, 11 px, obwiedziona kolorem tła panelu |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Inicjał lub obraz wg wariantu |
| Online / aktywny (wariant ze wskaźnikiem) | Kropka zielona (sukces) w rogu |
| Wyłączony | Nie dotyczy bezpośrednio — awatar nie jest kontrolką interaktywną sam w sobie; jeśli pełni funkcję przycisku (np. w pasku górnym, otwiera menu), przejmuje stany Przycisku ikonowego |
| Ładowanie (obraz) | Placeholder inicjałowy widoczny do czasu wczytania obrazu docelowego |
| Błąd (obraz) | Powrót do placeholderu inicjałowego przy niepowodzeniu wczytania obrazu |

**Zachowanie po interakcji.** W pasku górnym otwiera menu użytkownika lub profil po kliknięciu — w takim zastosowaniu zachowuje się jak Przycisk ikonowy (rozdz. 7.2). W Chat Window i panelach wielomodelowych pełni wyłącznie funkcję identyfikacyjną, bez własnej akcji kliknięcia.

**Gdzie używany.** Chat Window — identyfikacja uczestnika (użytkownik, model, rola), wspólne dla wszystkich piętnastu modułów; Model Panels w Roundtable — identyfikacja modeli uczestniczących; karty ról panelu orkiestracji MultitaskingAI — wariant kwadratowy dla agenta przypisanego z modułu Agents, wariant kołowy dla modelu bazowego bez przypisanego agenta; pasek górny — profil użytkownika.

---

### 12.2. Avatar Always On Display

| | |
|---|---|
| **Nazwa** | Avatar Always On Display |
| **Wzorzec** | `.dn-awatar` w pozycjonowaniu pływającym (fixed) + powierzchnia interakcji po aktywacji |
| **Kategoria** | F — Tożsamość |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji jako element pływający ponad całą powłoką; powierzchnia interakcji rozwija się po kliknięciu |

**Przeznaczenie.** Jest wizualnym zakotwiczeniem funkcji globalnej Always On Display — globalnego agenta towarzyszącego, nieposiadającego własnego środowiska ani modułu, obecnego jednocześnie we wszystkich częściach platformy jako nadrzędna warstwa inteligencji ponad środowiskami, modułami, projektami i sesjami. W odróżnieniu od Awatara ogólnego (rozdz. 12.1), który wyłącznie identyfikuje, ten komponent jest punktem dostępu — jego kluczowa idea to „osobisty agent operacyjny obecny zawsze i wszędzie”.

**Forma i waga.** Mała w spoczynku — pojedynczy pływający awatar (koło, wymiary jak wariant `--lg` Awatara, 48 px), zakotwiczony w stałej pozycji względem okna (fixed), ponad całą powłoką aplikacji niezależnie od otwartego środowiska, modułu czy karty sesji. Po aktywacji rozszerza się o powierzchnię interakcji otwieraną jako kolumna boczna po prawej stronie obszaru roboczego — rozszerzenie boczne o regulowanej szerokości, zbudowane z Panelu wysuwanego (rozdz. 15.6) z osadzonym strumieniem rozmowy.

**Warianty.**

| Kontekst | Zachowanie |
|---|---|
| Poza środowiskiem MultitaskingAI | Funkcja doradcza wykraczająca poza możliwości pojedynczego modułu — dostęp do wszystkich środowisk, projektów, sesji, historii rozmów oraz agentów użytkownika |
| Wewnątrz środowiska MultitaskingAI | Dodatkowo rola obserwatora lub operatora nadzorującego przebieg procesów złożonych z wielu zadań — widoczny w sekcji Monitor procesu panelu orkiestracji |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Spoczynek | Pływający awatar, stale widoczny, bez rozwiniętej powierzchni interakcji |
| Aktywny (po aktywacji) | Rozwinięta powierzchnia natychmiastowej interakcji — kolumna boczna po prawej stronie obszaru roboczego |
| Hover | Analogicznie do Karty interaktywnej lub Przycisku ikonowego — sygnalizacja klikalności |
| Wyłączony | Nie dotyczy w rozumieniu blokady — funkcja globalna pozostaje dostępna niezależnie od aktywnego środowiska lub modułu; wyłączenie widoczności jest ustawieniem konfiguracyjnym okna konfiguracji, nie stanem interfejsu |
| Ładowanie | Wewnątrz rozwiniętej powierzchni interakcji — Wskaźnik ładowania (rozdz. 11.5) na czas oczekiwania na odpowiedź |
| Tryb nadzoru (MultitaskingAI) | Wskaźnik roli (obserwator/operator) w sekcji Monitor procesu — realizowany Plakietką (rozdz. 10.5) przy nazwie Always On Display |

**Zachowanie po interakcji.** Aktywacja (kliknięcie lub gest dedykowany) otwiera natychmiastową powierzchnię interakcji, niezależnie od tego, w którym środowisku lub module znajduje się aktualnie użytkownik — funkcja globalna nie tworzy własnej przestrzeni roboczej i nie podlega mechanizmowi przeładowania właściwemu kartom sesji. Zamknięcie powierzchni interakcji przywraca stan spoczynku (pływający awatar), bez utraty dostępu — Always On Display pozostaje uruchomiony w tle.

**Funkcja zakotwiczona.** Reguły wyzwalania sugestii, katalog ich rodzajów, tor głosowy, zachowanie per środowisko i per moduł, warstwy widoczności oraz punkty sterowania są wspólne całej platformie.

**Gdzie używany.** Pozycja „Always On Display” listwy ustawień strefy 3 strony głównej (jeden z kilku równoważnych punktów dostępu — funkcja pozostaje dostępna również z poziomu każdego środowiska i modułu, nie tylko ze strony głównej); pływający element ponad całą powłoką aplikacji; sekcja Monitor procesu panelu orkiestracji środowiska MultitaskingAI, jako obserwator lub operator procesu.

---

## 13. Elementy pomocnicze

### 13.1. Ikona systemowa

| | |
|---|---|
| **Nazwa** | Ikona systemowa |
| **Klasa bazowa** | `.dn-ikona` |
| **Kategoria** | G — Elementy pomocnicze |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczna bez interakcji wewnątrz komponentu, który ją osadza |

**Przeznaczenie.** Nośnik pojedynczego symbolu graficznego z jednolitego zestawu 47 ikon — nigdy samodzielny cel interakcji (poza wyjątkowym zastosowaniem jako element klikalny z `role="button"`, przejmujący wtedy stany Przycisku ikonowego), zawsze wsparcie znaczeniowe dla tekstu lub kontrolki, którą towarzyszy.

**Forma i waga.** Znikoma — renderowana w 16, 18 lub 24 px (warianty `--sm` 14 px, `--lg` 20 px), rysowana na siatce 24×24 z obrysem 1,7 px, zaokrąglonymi końcami i łączeniami linii. Kolor zawsze dziedziczony z `currentColor` kontekstu — jedna ikona automatycznie zmienia barwę wraz z motywem lub stanem (np. na `--dn-akcent-txt` w stanie aktywnym), bez potrzeby odrębnego pliku na kolor.

**Warianty.** Rozmiar (`--sm`/`--lg`/bazowy) jako jedyny wariant klasy; różnorodność właściwa nie wariantom CSS, lecz wyborowi konkretnego pliku SVG spośród 47 dostępnych (adres podawany funkcjonalnie przez `--dn-ikona-src`, nie stylistycznie).

**Stany.** Ikona sama w sobie jest bierna wizualnie — nie ma stanów hover/aktywny/wyłączony własnych; dziedziczy kolor (a więc pośrednio stan) z elementu, w którym jest osadzona.

**Zachowanie po interakcji.** Bierne, poza wyjątkiem klikalnej ikony z `role="button"`/`tabindex`, która przejmuje pełny zestaw stanów interakcji i pierścień fokusu wspólny dla kontrolek interaktywnych.

**Gdzie używany.** Wewnątrz niemal każdego innego komponentu katalogu — Przycisku, Przycisku ikonowego, Karty pozycji, Plakietki, Paska górnego, listwy ustawień. Pełny katalog 47 ikon wraz z typowym zastosowaniem w platformie znajduje się w Systemie wizualnym, Załącznik C.

---

### 13.2. Separator

| | |
|---|---|
| **Nazwa** | Separator |
| **Klasa bazowa** | `.dn-separator` |
| **Kategoria** | G — Elementy pomocnicze |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji wewnątrz panelu lub karty |

**Przeznaczenie.** Rozdziela wizualnie dwie sekcje treści w obrębie tego samego panelu lub karty, gdy podział ma być sygnalizowany bez wprowadzania nowego nagłówka ani nowej ramki.

**Forma i waga.** Znikoma — pozioma linia 1 px w kolorze obrysu (`--dn-obrys`), pełna szerokość rodzica, z marginesem pionowym 16 px (`--dn-od-4`) powyżej i poniżej.

**Warianty.** Brak — jedna forma.

**Stany.** Statyczny element dekoracyjny bez stanów interakcji.

**Zachowanie po interakcji.** Nie dotyczy — element bierny, niesklikalny.

**Gdzie używany.** Wewnątrz Paneli i Kart wszędzie tam, gdzie potrzebny jest podział sekcji bez wprowadzania dodatkowego Nagłówka.

---

## 14. Karty komponentów — Okna komunikacji operacyjnej

Kategoria H obejmuje dwa okna, które platforma udostępnia w każdym module i w każdym środowisku jako elementy pierwszoplanowe architektury interfejsu, nie jako panele towarzyszące. Oba są komponentami wielokrotnego użytku w rozumieniu niniejszego katalogu: mają stałe miejsce w układzie, stały zestaw elementów składowych, stały zestaw stanów i identyczną formę we wszystkich środowiskach.

**Diagram — miejsce obu okien w kanonicznym układzie kolumnowym.**

```
 ═══════════════════════════════════════════════════════════════════════════
  Boczna     │ CHAT WINDOW (14.1)   │ Obszar roboczy modułu │ Panel
  nawigacja  │ Użytkownik ↔         │                       │ pomocniczy
  modułów    │ Wykonawca            │                       │ (rozszerzenie
             │ stała, pełna wysokość│                       │  boczne, 15.6)
             │ ─────────────────    │                       │
             │ EXECUTION LOOP       │                       │
             │ WINDOW (14.2)        │                       │
             │ Koordynator ↔        │                       │
             │ Wykonawca            │                       │
 ═══════════════════════════════════════════════════════════════════════════
```

---

### 14.1. Chat Window

| | |
|---|---|
| **Nazwa** | Chat Window — główne okno komunikacji (kanał Użytkownik ↔ Wykonawca) |
| **Wzorzec** | Kolumna złożona z `.dn-karta--pozycja` (wiadomości), `.dn-awatar`, `.dn-kod`, `.dn-plakietka`, `.dn-input` / `.dn-textarea`, `.dn-btn-ikona` |
| **Kategoria** | H — Okna komunikacji operacyjnej |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 (kanał wywołania funkcji warstwy 4) |
| **Sposób wywołania** | Widoczne bez interakcji — lewa kolumna obszaru roboczego, w tym samym miejscu układu w każdym module i w każdym środowisku |

**Przeznaczenie.** Jest głównym oknem komunikacji między Użytkownikiem a Wykonawcą (AI, agentem lub systemem wykonawczym) — centralnym punktem pracy użytkownika i podstawowym mechanizmem sterowania wszystkimi procesami realizowanymi przez platformę. Przyjmuje polecenia w języku naturalnym, prezentuje strumień odpowiedzi i wyników, umożliwia zatwierdzanie oraz przerywanie działań, a także wyjaśnianie wyniku i kontekstu. Jest równocześnie kanałem dostępu do funkcji warstwy 4 (rozdz. 3a): operacje najbardziej zaawansowane, tryby administracyjne i narzędzia diagnostyczne niskiego poziomu wywoływane są poleceniem języka naturalnego, bez własnej reprezentacji graficznej w interfejsie.

**Forma i waga.** Duża — lewa kolumna, stała, pełna wysokość obszaru roboczego. Kolumna sąsiaduje bezpośrednio z boczną nawigacją modułów po lewej i z kolumną obszaru roboczego modułu po prawej. Regulacji podlega wyłącznie szerokość kolumny; okno nie zmienia położenia, nie zwija się do ikony i nie ustępuje miejsca zawartości modułu. Wewnątrz kolumny treść układa się w trzech strefach pionowych: pasek kontekstu ze znacznikami (rozdz. 15.9) u wierzchołka kolumny, strumień wiadomości pośrodku (przewijany niezależnie od reszty okna), strefa wprowadzania polecenia przy podstawie kolumny.

**Elementy składowe.**

| Element | Zawartość | Komponent nośny | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Pasek kontekstu | Znaczniki środowiska, repozytorium, projektu, modelu i wykonawcy, np. `[Danaco Console] [Ubuntu] [Worktree] [Fable 5] [Ultra]` | Pasek kontekstu ze znacznikami (15.9) | 1 | Widoczny bez interakcji |
| Selektor przypisany znacznikowi | Lista wartości dla klikniętego znacznika | Znacznik kontekstowy (15.7) + Panel popover (15.5) | 2 | Kliknięcie znacznika; po wyborze selektor zwija się samoczynnie |
| Strumień wiadomości | Wiadomości Użytkownika i Wykonawcy w porządku chronologicznym, z awatarem nadawcy | Karta pozycji (10.1) + Awatar (12.1) | 1 | Widoczny bez interakcji |
| Treść techniczna w wiadomości | Polecenie, ścieżka, fragment kodu, ładunek narzędzia | Blok kodu (10.6) | 1 | Widoczna bez interakcji |
| Wskaźnik pracy Wykonawcy | Sygnalizacja trwającego przetwarzania odpowiedzi | Wskaźnik ładowania (11.5) | 1 | Widoczny na czas trwania operacji |
| Sterowanie przebiegiem | Zatwierdzenie działania, przerwanie działania, ponowienie | Przycisk (7.1) w treści wiadomości wymagającej decyzji | 1 | Widoczne bez interakcji w wiadomości, której dotyczy |
| Pole polecenia | Wprowadzanie polecenia w języku naturalnym | Pole tekstowe (8.1), wariant wieloliniowy | 1 | Widoczne bez interakcji |
| Wysłanie polecenia | Przekazanie polecenia Wykonawcy | Przycisk ikonowy (7.2) | 1 | Widoczne bez interakcji |
| Wybór wykonawcy, modelu i poziomu wysiłku | Zwinięty element `Agent ▼` przy polu polecenia | Menu progresywne (15.2) | 2 | Kliknięcie elementu; po wyborze element zwija się samoczynnie |
| Zestaw operacji na rozmowie | Wyczyszczenie kontekstu, wyeksportowanie zapisu, wyjaśnienie wyniku, przypięcie wiadomości | Element zbiorczy z rozwinięciem `Operacje ▼` (15.1) lub Menu kebab ⋮ (15.3) | 3 | Kliknięcie wyzwalacza |
| Wywołanie funkcji eksperckiej | Tryb administracyjny, narzędzie diagnostyczne niskiego poziomu, operacja zaawansowana | Polecenie języka naturalnego w polu polecenia | 4 | Wpisanie polecenia; brak reprezentacji graficznej |

**Warianty.**

| Wariant | Zawartość i różnica | Gdzie występuje |
|---|---|---|
| Bazowy | Pełny zestaw elementów wg tabeli wyżej | Każdy moduł każdego środowiska |
| Z cytowaniem kodu | Strumień wiadomości z rozbudowanym Blokiem kodu (10.6) i akcją kopiowania w każdym cytowaniu | Developer, Terminal, Diagnostics |
| Z odwołaniami do zasobów | Wiadomość niesie odwołania do pozycji Panelu (10.2) modułu — kliknięcie odwołania podświetla pozycję w kolumnie obszaru roboczego | Research, Library, Workspace |
| Rozmowa wielomodelowa | Strumień rozdzielony na równoległe odpowiedzi kilku Wykonawców, przełączane Zakładką (9.4) | Roundtable |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (rozmowa w toku) | Strumień wiadomości, pole polecenia gotowe do wprowadzenia treści |
| Pusty (rozmowa nierozpoczęta) | Pusty stan (10.7) w miejscu strumienia — ikona, tytuł, krótki opis; pole polecenia pozostaje gotowe |
| Wykonawca przetwarza polecenie | Wskaźnik ładowania (11.5) na końcu strumienia; pole polecenia pozostaje w pełni edytowalne, przycisk wysłania w pełni klikalny |
| Działanie oczekuje na zatwierdzenie | Wiadomość z osadzonymi Przyciskami decyzji (zatwierdź, odrzuć, skoryguj), wyróżniona Komunikatem blokowym (11.3) wariantu `--info` |
| Działanie przerwane przez Użytkownika | Wiadomość zamknięta Plakietką (10.5) wariantu `--ostrz` z oznaczeniem przerwania |
| Błąd wykonania | Komunikat blokowy (11.3) wariantu `--blad` osadzony w strumieniu, z akcją ponowienia; strumień pozostaje przewijalny, pole polecenia edytowalne |
| Fokus klawiaturowy w polu polecenia | Pierścień złoty wspólny dla kontrolek interaktywnych na Polu tekstowym (8.1) |

**Zachowanie po interakcji.** Wprowadzenie polecenia i jego wysłanie przekazuje treść Wykonawcy i dopisuje wiadomość na końcu strumienia; strumień przewija się samoczynnie do najnowszej wiadomości, chyba że użytkownik przewinął ręcznie wyżej. Zatwierdzenie działania uruchamia je natychmiast; przerwanie kończy bieżące działanie i pozostawia w strumieniu zapis stanu, w którym zostało przerwane. Polecenie skierowane do procesu prowadzonego przez Koordynatora ujawnia swój przebieg w Execution Loop Window (rozdz. 14.2) — Chat Window prezentuje wynik i punkty decyzyjne, Execution Loop Window prezentuje przebieg pętli wykonawczej. Zmiana wartości znacznika w pasku kontekstu obowiązuje od następnego polecenia i nie przerywa działań już uruchomionych.

**Gdzie używany.** Lewa kolumna obszaru roboczego każdego modułu każdego z czterech środowisk (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) — jedyne miejsce wystąpienia, obecne w każdej karcie sesji platformy.

---

### 14.2. Execution Loop Window

| | |
|---|---|
| **Nazwa** | Execution Loop Window — okno pętli wykonawczej (kanał Koordynator ↔ Wykonawca) |
| **Wzorzec** | Kolumna złożona z `.dn-karta--pozycja` (zadania i komunikaty), `.dn-tabela` (kolejka), `.dn-plakietka--stan`, `.dn-kropka`, `.dn-btn` |
| **Kategoria** | H — Okna komunikacji operacyjnej |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczne bez interakcji w kolumnie sąsiadującej z Chat Window; szerokość kolumny podlega regulacji |

**Przeznaczenie.** Prezentuje komunikację między Koordynatorem — komponentem orkiestrującym platformy — a Wykonawcą. Odpowiada za prowadzenie pętli wykonawczej, koordynację zadań, nadzór nad realizacją, orkiestrację działań oraz kontrolę realizacji procesów. Jest oknem obserwacji i sterowania przebiegiem, nie oknem zlecania: zlecenie pochodzi z Chat Window (rozdz. 14.1), a Execution Loop Window pokazuje, jak zostało zdekomponowane, komu przydzielone i w jakim jest stanie.

**Forma i waga.** Duża — kolumna sąsiadująca z Chat Window, pełna wysokość obszaru roboczego, wspólna z nim strefa lewej części układu. Regulacji podlega wyłącznie szerokość kolumny. Treść układa się w czterech strefach pionowych: bieżące zlecenie wraz z dekompozycją, kolejka i stan zadań, wymiana komunikatów sterujących, wskaźniki przebiegu pętli wraz ze sterowaniem przebiegiem.

**Elementy składowe.**

| Element | Zawartość | Komponent nośny | Warstwa | Sposób wywołania |
|---|---|---|---|---|
| Bieżące zlecenie | Treść zlecenia przyjętego od Użytkownika i jego dekompozycja na zadania | Karta (10.1) z listą pozycji | 1 | Widoczne bez interakcji |
| Kolejka i stan zadań | Zadania w porządku wykonania, z przypisanym Wykonawcą i stanem | Tabela (10.4) + Pigułka statusu (10.5) | 1 | Widoczna bez interakcji |
| Wymiana komunikatów sterujących | Komunikaty Koordynatora do Wykonawcy i odpowiedzi Wykonawcy | Karta pozycji (10.1) | 1 | Widoczna bez interakcji |
| Wyniki kontroli jakości | Ocena wyniku zadania i decyzja o ponowieniu | Karta pozycji (10.1) + Plakietka (10.5) wariantu `--sukces` / `--ostrz` / `--blad` | 1 | Widoczne bez interakcji |
| Wskaźniki przebiegu pętli | Numer przebiegu, liczba zadań zakończonych i pozostałych, czas trwania | Pigułka statusu i Kropka stanu (10.5) | 1 | Widoczne bez interakcji |
| Sterowanie przebiegiem | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia | Element zbiorczy z rozwinięciem `Operacje ▼` (15.1) | 3 | Kliknięcie elementu zbiorczego |
| Szczegóły pojedynczego zadania | Pełny zapis komunikatów, ładunków narzędzi i wyniku zadania | Panel wysuwany (15.6) z Blokiem kodu (10.6) | 3 | Kliknięcie wiersza zadania; po zamknięciu panel znika z układu |
| Akcje pojedynczego zadania | Ponowienie, pominięcie, przypisanie innemu Wykonawcy | Menu kebab ⋮ (15.3) w wierszu zadania | 3 | Kliknięcie ⋮ |
| Parametry pętli | Liczba ponowień, próg kontroli jakości, poziom wysiłku Wykonawcy | Menu progresywne (15.2) + pola formularzy (rozdz. 8) | 2 | Kliknięcie elementu zwiniętego; po wyborze element zwija się samoczynnie |
| Diagnostyka przebiegu niskiego poziomu | Pełny zapis pętli, ślad wywołań narzędzi, zrzut stanu Koordynatora | Paleta poleceń (15.8) lub polecenie w Chat Window | 4 | Skrót `Ctrl/Cmd + K` albo polecenie języka naturalnego |

**Warianty.**

| Wariant | Zawartość i różnica | Gdzie występuje |
|---|---|---|
| Bazowy | Pełny zestaw czterech stref wg opisu wyżej | Każdy moduł, w którym zlecenie podlega dekompozycji |
| Pętla jednozadaniowa | Strefa kolejki zredukowana do jednego wiersza; wskaźniki przebiegu ograniczone do stanu i czasu | Moduły o pracy liniowej (Translate, Library) |
| Pętla wielowykonawcza | Kolejka rozszerzona o kolumnę Wykonawcy i kolumnę zespołu; wymiana komunikatów grupowana per Wykonawca | MultitaskingAI, Subagent Network |
| Pętla z kontrolą jakości | Strefa wyników kontroli rozbudowana o ocenę i jawną decyzję o ponowieniu | Research, Developer, Diagnostics |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny (pętla w toku) | Kolejka z wierszami w stanach zakończony / w toku / oczekujący, wskaźniki przebiegu aktualizowane na żywo kanałem WebSocket |
| Pusty (brak zlecenia) | Pusty stan (10.7) — okno bez przyjętego zlecenia, z opisem wskazującym Chat Window jako miejsce zlecania |
| Wstrzymana | Pigułka statusu wariantu `--ostrz` przy wskaźniku przebiegu; kolejka zachowana bez zmian, sterowanie przebiegiem udostępnia wznowienie |
| Ponowienie zadania | Wiersz zadania oznaczony Plakietką ponowienia z licznikiem przebiegu; poprzedni wynik pozostaje dostępny w Panelu wysuwanym szczegółów |
| Zakończona powodzeniem | Wskaźnik przebiegu z Pigułką wariantu `--sukces`; kolejka pozostaje widoczna jako zapis przebiegu |
| Błąd pętli | Komunikat blokowy (11.3) wariantu `--blad` nad kolejką, z akcją ponowienia i akcją korekty zlecenia; wszystkie kontrolki pozostają klikalne |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu kolejki do czasu odebrania pierwszego stanu zadań |

**Zachowanie po interakcji.** Okno aktualizuje treść na żywo, niezależnie od akcji użytkownika. Wstrzymanie zatrzymuje przydzielanie kolejnych zadań, pozostawiając zadanie bieżące do naturalnego zakończenia; wznowienie przywraca przydzielanie od pierwszego zadania oczekującego. Przerwanie kończy pętlę i pozostawia kolejkę jako zapis stanu, w którym została przerwana. Korekta zlecenia przekazuje Koordynatorowi zmienioną treść zlecenia i wyzwala ponowną dekompozycję — zadania już zakończone pozostają zachowane. Kliknięcie wiersza zadania otwiera Panel wysuwany ze szczegółami, bez przesuwania i bez przeładowania pozostałych kolumn układu.

**Gdzie używany.** Kolumna sąsiadująca z Chat Window w obszarze roboczym każdego modułu każdego z czterech środowisk; największe nasilenie w środowisku MultitaskingAI, gdzie prowadzi pętle wielowykonawcze zespołów i Subagent Network, oraz w modułach Automations, Developer i Diagnostics.

---

## 15. Karty komponentów — Ujawnianie funkcjonalności

Kategoria I obejmuje komponenty realizujące regułę warstw widoczności (rozdz. 3a) — elementy, których jedynym zadaniem jest utrzymanie funkcji poza polem widzenia do chwili wystąpienia potrzeby jej użycia oraz udostępnienie jej jednym kliknięciem, jednym skrótem klawiszowym albo jednym poleceniem języka naturalnego. Wszystkie komponenty tej kategorii prezentują pełną listę pozycji na jednym poziomie rozwinięcia; drugi poziom rozwinięcia nie występuje w katalogu.

---

### 15.1. Element zbiorczy z rozwinięciem (`Operacje ▼`)

| | |
|---|---|
| **Nazwa** | Element zbiorczy z rozwinięciem |
| **Wzorzec** | `.dn-btn--zarys` z ikoną `strzalka-dol` + `.dn-karta` z `.dn-karta--pozycja` (lista rozwinięta) |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie elementu zwiniętego; rozwinięcie zamyka się po wyborze pozycji, kliknięciu poza obszarem lub klawiszem Escape |

**Przeznaczenie.** Zastępuje zestaw przycisków jednym elementem zbiorczym, którego rozwinięcie zawiera pełną listę akcji — realizuje mechanizm grupowania logicznego akcji (rozdz. 3a). Stosowany wszędzie tam, gdzie liczba akcji dostępnych dla jednego obiektu przekracza dwa przyciski, a żadna z nich nie jest akcją podstawową widoku.

**Forma i waga.** Mała w spoczynku — pojedynczy Przycisk wariantu `--zarys` z etykietą nazywającą zbiór (`Operacje ▼`, `Widok ▼`, `Eksport ▼`) i ikoną `strzalka-dol` 14 px po prawej stronie etykiety. Rozwinięcie: lista pływająca zakotwiczona przy elemencie, szerokość dopasowana do najdłuższej pozycji, wysokość wg liczby pozycji, tło `--dn-panel`, cień `--dn-cien-3`, promień `--dn-r-md`.

**Warianty.**

| Wariant | Etykieta zwinięta | Zawartość rozwinięcia |
|---|---|---|
| Operacje na obiekcie | `Operacje ▼` | Akcje dotyczące bieżącego dokumentu, zadania, przebiegu lub zasobu |
| Sterowanie przebiegiem | `Operacje ▼` | Wstrzymanie, wznowienie, przerwanie, korekta zlecenia (Execution Loop Window, rozdz. 14.2) |
| Widok | `Widok ▼` | Warianty prezentacji treści kolumny obszaru roboczego |
| Eksport | `Eksport ▼` | Formaty i zakresy eksportu |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zwinięty (spoczynek) | Przycisk `--zarys` z ikoną `strzalka-dol`; lista niewidoczna, bez śladu w układzie |
| Hover | Tło akcentu złotego na elemencie zwiniętym, zgodnie z wariantem `--zarys` Przycisku (7.1) |
| Rozwinięty | Ikona `strzalka-dol` obrócona o 180°, lista widoczna, każda pozycja zachowuje się jak Karta interaktywna (10.1) |
| Fokus klawiaturowy | Pierścień złoty na elemencie zwiniętym; klawisze strzałek przenoszą fokus między pozycjami rozwiniętej listy |
| Pozycja z niespełnionym warunkiem | Pozycja pozostaje w pełni klikalna — niedostępność sygnalizowana opisowo tekstem pomocniczym lub Dymkiem kontekstowym (11.4), a kliknięcie zwraca komunikat |
| Ładowanie | Lista pozycji zależnych od danych pobieranych asynchronicznie prezentuje Wskaźnik ładowania (11.5) do czasu wypełnienia |
| Błąd | Nie dotyczy elementu — błąd wykonania wybranej akcji komunikowany Powiadomieniem (11.2) po zamknięciu rozwinięcia |

**Zachowanie po interakcji.** Kliknięcie rozwija listę bez przesuwania i bez przeładowania reszty widoku. Wybór pozycji wykonuje akcję i zwija element z powrotem do postaci spoczynkowej. Rozwinięcie nigdy nie prowadzi do drugiego poziomu listy — akcja złożona otwiera zamiast tego Modal (11.1) lub Panel wysuwany (15.6).

**Gdzie używany.** Chat Window (operacje na rozmowie); Execution Loop Window (sterowanie przebiegiem); nagłówki Paneli i okien operacyjnych we wszystkich modułach; nagłówek Tabeli w oknach typu monitor i manager.

---

### 15.2. Menu progresywne (`Agent ▼`)

| | |
|---|---|
| **Nazwa** | Menu progresywne |
| **Wzorzec** | `.dn-btn--duch` z ikoną `strzalka-dol` + `.dn-karta--pozycja` (lista wyborów) |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Kliknięcie elementu zwiniętego; po wyborze wartości element zwija się samoczynnie |

**Przeznaczenie.** Prezentuje zbiór jednorodnych wyborów jako jeden element zwinięty, w którym widoczna jest wyłącznie wartość aktualnie obowiązująca (`Agent ▼`, `Model ▼`, `Środowisko ▼`, `Wysiłek ▼`). Odróżnia się od Listy rozwijanej (8.2) tym, że nie jest polem formularza: nie ma etykiety nad sobą, nie ma obrysu pola i nie uczestniczy w zapisie formularza — jest lekkim przełącznikiem kontekstu pracy.

**Forma i waga.** Znikoma w spoczynku — Przycisk wariantu `--duch`, czcionka 13 px, bez obrysu i bez tła, z nazwą wybranej wartości i ikoną `strzalka-dol` 14 px. Rozwinięcie: lista pływająca o tej samej formie co rozwinięcie Elementu zbiorczego (15.1), z zaznaczeniem pozycji obowiązującej.

**Warianty.**

| Wariant | Etykieta zwinięta | Zawartość rozwinięcia |
|---|---|---|
| Wybór wykonawcy | `Agent ▼` | Wykonawcy dostępni w bieżącym kontekście pracy |
| Wybór modelu | `Model ▼` | Modele bazowe udostępnione dla wybranego Wykonawcy |
| Wybór środowiska | `Środowisko ▼` | Środowiska uruchomieniowe dostępne dla bieżącej sesji |
| Wybór trybu pracy | `Tryb ▼` | Tryby pracy modułu |
| Poziom wysiłku | `Wysiłek ▼` | Poziomy wysiłku Wykonawcy |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zwinięte (spoczynek) | Nazwa wartości obowiązującej + ikona `strzalka-dol`, tekst `--dn-tekst-3` |
| Hover | Tekst `--dn-tekst` (podstawowy), tło `--dn-hover` |
| Rozwinięte | Lista wartości, pozycja obowiązująca oznaczona ikoną `zaznacz` i tekstem `--dn-akcent-txt` |
| Fokus klawiaturowy | Pierścień złoty; klawisze strzałek przenoszą fokus między wartościami |
| Wartość odziedziczona z poziomu szerszego | Nazwa wartości uzupełniona Plakietką „odziedziczone” (10.5); wybór innej wartości tworzy jawne nadpisanie na bieżącym poziomie |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu listy do czasu odebrania zbioru wartości |
| Błąd | Komunikat w miejscu listy z akcją ponowienia; element zwinięty zachowuje ostatnią obowiązującą wartość |

**Zachowanie po interakcji.** Kliknięcie rozwija listę, wybór wartości zwija ją samoczynnie i podmienia etykietę elementu zwiniętego. Zmiana wartości obowiązuje od następnego polecenia lub następnego uruchomienia procesu i nie przerywa działań już uruchomionych.

**Gdzie używany.** Chat Window, przy polu polecenia (wybór wykonawcy, modelu, poziomu wysiłku); Execution Loop Window (parametry pętli); nagłówek okna operacyjnego w modułach o wielu trybach pracy (Roundtable, Translate, Developer, Terminal).

---

### 15.3. Menu kebab (⋮)

| | |
|---|---|
| **Nazwa** | Menu kebab |
| **Wzorzec** | `.dn-btn-ikona` (ikona `wiecej`, ⋮) + `.dn-karta` z `.dn-karta--pozycja` |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie ikony ⋮ przy elemencie, którego dotyczy |

**Przeznaczenie.** Udostępnia zestaw akcji dotyczących jednej, konkretnej pozycji — wiersza tabeli, karty, zadania, zasobu — bez zajmowania miejsca w spoczynku. Jest podstawowym wyzwalaczem Menu kontekstowego (9.5) w wariancie menu akcji pozycji.

**Forma i waga.** Znikoma — Przycisk ikonowy 36×36 px z ikoną `wiecej` 18×18 px w postaci trzech kropek ułożonych pionowo, bez tła i obrysu w spoczynku, umieszczony przy prawej krawędzi wiersza lub karty.

**Warianty.**

| Wariant | Umiejscowienie | Zawartość rozwinięcia |
|---|---|---|
| Akcje wiersza | Ostatnia kolumna Tabeli (10.4) | Edycja, usunięcie, duplikacja, eksport, ponowienie zadania |
| Akcje karty | Nagłówek Karty (10.1) | Akcje dotyczące pozycji reprezentowanej przez kartę |
| Akcje okna | Nagłówek okna operacyjnego | Akcje dotyczące całego okna, nieujęte w Elemencie zbiorczym (15.1) |

**Stany.** Ikona przejmuje pełny zestaw stanów Przycisku ikonowego (7.2): domyślny, hover, aktywny, fokus klawiaturowy. Rozwinięta lista przejmuje stany Menu kontekstowego (9.5). Pozycja z niespełnionym warunkiem pozostaje w pełni klikalna, a niedostępność jest sygnalizowana opisowo.

**Zachowanie po interakcji.** Kliknięcie otwiera listę zakotwiczoną przy ikonie. Kliknięcie poza obszarem listy, klawisz Escape lub wybór pozycji zamyka ją. Lista nie przesuwa ani nie przeładowuje reszty widoku.

**Gdzie używany.** Wiersze Tabeli we wszystkich oknach typu monitor i manager; karty pozycji Paneli; wiersze kolejki zadań Execution Loop Window (14.2); nagłówki okien operacyjnych.

---

### 15.4. Menu hamburger (☰)

| | |
|---|---|
| **Nazwa** | Menu hamburger |
| **Wzorzec** | `.dn-btn-ikona` (ikona `menu`, ☰) + `.dn-karta` z `.dn-karta--pozycja` |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie ikony ☰ w pasku górnym lub w nagłówku kolumny |

**Przeznaczenie.** Udostępnia zestaw akcji i przejść dotyczących całego widoku lub całej kolumny — w odróżnieniu od Menu kebab (15.3), które dotyczy pojedynczej pozycji. Steruje również widocznością bocznej nawigacji modułów w widoku wąskim.

**Forma i waga.** Znikoma — Przycisk ikonowy 36×36 px z ikoną `menu` 18×18 px w postaci trzech poziomych kresek, umieszczony po lewej stronie paska górnego lub przy lewej krawędzi nagłówka kolumny.

**Warianty.**

| Wariant | Umiejscowienie | Skutek |
|---|---|---|
| Przełącznik bocznej nawigacji | Pasek górny, przy lewej krawędzi | Przełącza widoczność kolumny bocznej nawigacji / panelu orkiestracji (9.2) |
| Menu kolumny | Nagłówek kolumny obszaru roboczego lub kolumny pomocniczej | Akcje i przejścia dotyczące całej kolumny |
| Menu widoku | Nagłówek okna operacyjnego | Przejścia między widokami okna nieujęte w Zakładkach (9.4) |

**Stany.** Ikona przejmuje pełny zestaw stanów Przycisku ikonowego (7.2). W wariancie przełącznika bocznej nawigacji ikona pozostaje w tym samym miejscu niezależnie od tego, czy kolumna jest widoczna — stan kolumny sygnalizowany jest jej obecnością w układzie, nie zmianą ikony. Rozwinięta lista przejmuje stany Menu kontekstowego (9.5).

**Zachowanie po interakcji.** Kliknięcie w wariancie przełącznika ukrywa lub przywraca kolumnę bocznej nawigacji; ukrycie kolumny poszerza pozostałe kolumny układu, nie zmieniając ich kolejności ani proporcji względem siebie. Kliknięcie w wariantach menu otwiera listę zakotwiczoną przy ikonie, zamykaną kliknięciem poza obszarem, klawiszem Escape lub wyborem pozycji.

**Gdzie używany.** Pasek górny okna każdego z czterech środowisk; nagłówki kolumn obszaru roboczego i kolumn pomocniczych; nagłówki okien operacyjnych o wielu widokach.

---

### 15.5. Panel popover

| | |
|---|---|
| **Nazwa** | Panel popover |
| **Wzorzec** | `.dn-karta` pływająca zakotwiczona przy wyzwalaczu, z dowolną treścią komponentową wewnątrz |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie wyzwalacza — znacznika kontekstowego (15.7), ikony, przycisku lub pozycji listy |

**Przeznaczenie.** Udostępnia niewielki, samodzielny zestaw treści lub ustawień zakotwiczony przy elemencie, którego dotyczy — selektor wartości znacznika, ustawienia szybkie, podgląd pozycji, formularz krótki. Odróżnia się od Menu kontekstowego (9.5) tym, że zawiera treść komponentową (pola, tabelę, opis), nie wyłącznie listę akcji; od Modala (11.1) tym, że nie przerywa pracy i nie przechwytuje całego ekranu.

**Forma i waga.** Mała do średniej — pływająca karta zakotwiczona przy wyzwalaczu, szerokość 240–420 px, wysokość wg treści, tło `--dn-panel`, cień `--dn-cien-3`, promień `--dn-r-md`. Po zamknięciu znika całkowicie z przestrzeni roboczej, bez śladu w układzie.

**Warianty.**

| Wariant | Wyzwalacz | Zawartość |
|---|---|---|
| Selektor znacznika | Znacznik kontekstowy (15.7) w pasku kontekstu | Lista wartości dla wymiaru opisanego znacznikiem, z wyszukiwaniem przy zbiorze przekraczającym dwanaście pozycji |
| Ustawienia szybkie | Ikona `ustawienia` w nagłówku kolumny lub okna | Podzbiór ustawień warstwy sesji — pola formularzy (rozdz. 8) |
| Podgląd pozycji | Kliknięcie pozycji listy lub wiersza tabeli | Skrócony opis pozycji, Plakietki statusu, akcja przejścia do pełnego widoku |
| Formularz krótki | Przycisk akcji wymagającej jednej wartości | Pojedyncze Pole tekstowe (8.1) lub Lista rozwijana (8.2) z akcją potwierdzenia |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zamknięty (spoczynek) | Niewidoczny, bez śladu w układzie |
| Otwierający się | Animacja wjazdu — przezroczystość 0→1 i przesunięcie 8 px→0 w czasie 0,16 s, pomijana przy `prefers-reduced-motion` |
| Otwarty | Widoczna karta pływająca; fokus klawiatury przenosi się do pierwszego elementu interaktywnego treści |
| Fokus klawiaturowy | Pierścień złoty na elementach treści; klawisz Escape zamyka panel i przywraca fokus na wyzwalaczu |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu treści do czasu odebrania danych |
| Błąd | Komunikat blokowy (11.3) wariantu `--blad` w treści panelu, z akcją ponowienia |

**Zachowanie po interakcji.** Kliknięcie poza obszarem panelu, klawisz Escape lub zakończenie akcji zamyka panel. Panel nie przesuwa ani nie przeładowuje reszty widoku i nie blokuje interakcji z pozostałymi kolumnami układu.

**Gdzie używany.** Pasek kontekstu Chat Window (selektory znaczników); nagłówki kolumn i okien operacyjnych (ustawienia szybkie); listy pozycji Paneli (podgląd pozycji); wiersze Tabeli w oknach typu manager.

---

### 15.6. Panel wysuwany

| | |
|---|---|
| **Nazwa** | Panel wysuwany |
| **Wzorzec** | Region `.dn-panel` otwierany jako kolumna boczna, z treścią budowaną z komponentów rozdz. 8, 10 i 11 |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 3 |
| **Sposób wywołania** | Kliknięcie pozycji, wiersza, ikony lub pozycji menu otwierającej szczegóły; zamknięcie kontrolką zamknięcia lub klawiszem Escape |

**Przeznaczenie.** Udostępnia rozbudowaną treść pomocniczą — szczegóły pozycji, pełny zapis zadania, dokumentację ustawienia, historię zmian — bez opuszczania bieżącego widoku i bez trwałego zajmowania miejsca w układzie. Jest podstawowym nośnikiem mechanizmu paneli wysuwanych (rozdz. 3a): po zamknięciu znika całkowicie z przestrzeni roboczej.

**Forma i waga.** Średnia do dużej — kolumna boczna otwierana jako rozszerzenie boczne, po prawej stronie obszaru roboczego, na pełną wysokość obszaru roboczego. Szerokość regulowana, wyjściowo 320–480 px. Otwarcie panelu zwęża kolumnę obszaru roboczego, nie przykrywa jej i nie zmienia kolejności pozostałych kolumn układu.

**Warianty.**

| Wariant | Zawartość | Przykład wystąpienia |
|---|---|---|
| Szczegóły pozycji | Pełny opis pozycji listy lub wiersza tabeli, z akcjami w stopce | Execution Loop Window (14.2) — szczegóły zadania |
| Zapis techniczny | Blok kodu (10.6) z pełnym zapisem komunikatów, ładunków narzędzi i wyniku | Diagnostics, Developer, Terminal |
| Panel narzędzi | Lista operacji dostępnych dla treści kolumny obszaru roboczego | Studio — Tools Panel; Design — Assets Panel |
| Dokumentacja ustawienia | Opis ustawienia, jego zakresu i skutków, z odwołaniem do okna konfiguracji | Okno konfiguracji punktów izolacji |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zamknięty (spoczynek) | Niewidoczny, bez śladu w układzie; kolumna obszaru roboczego zajmuje pełną dostępną szerokość |
| Otwierający się | Animacja wysunięcia — szerokość 0→docelowa w czasie 0,24 s (`--dn-czas-3`), pomijana przy `prefers-reduced-motion` |
| Otwarty | Kolumna widoczna, treść przewijana niezależnie od pozostałych kolumn |
| Regulacja szerokości | Krawędź lewa kolumny jest uchwytem regulacji; regulacji podlega wyłącznie szerokość |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu treści do czasu odebrania danych |
| Pusty | Pusty stan (10.7) w miejscu treści |
| Błąd | Komunikat blokowy (11.3) wariantu `--blad` nad treścią, z akcją ponowienia |

**Zachowanie po interakcji.** Otwarcie panelu nie przerywa pracy w pozostałych kolumnach — Chat Window, Execution Loop Window i obszar roboczy pozostają w pełni interaktywne. Zamknięcie usuwa kolumnę z układu i przywraca pierwotną szerokość kolumny obszaru roboczego. Jednocześnie otwarty pozostaje jeden panel wysuwany na kolumnę obszaru roboczego; otwarcie kolejnego podmienia treść dotychczasowego.

**Gdzie używany.** Execution Loop Window (szczegóły zadania); okna operacyjne wszystkich modułów jako Tools Panel, Assets Panel, Session Repository, Diff/Grep Panel, Logs Viewer; powierzchnia interakcji Avatara Always On Display (12.2).

---

### 15.7. Znacznik kontekstowy (tag) otwierający selektor

| | |
|---|---|
| **Nazwa** | Znacznik kontekstowy |
| **Klasa bazowa** | `.dn-plakietka--interaktywna` |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Zdefiniowany |
| **Warstwa widoczności** | 2 |
| **Sposób wywołania** | Kliknięcie znacznika otwiera Panel popover (15.5) z selektorem; po wyborze wartości selektor zwija się samoczynnie |

**Przeznaczenie.** Prezentuje jeden wymiar kontekstu pracy — środowisko, repozytorium, projekt, model, wykonawcę, poziom wysiłku — jako lekki, jednosłowny znacznik, i jednocześnie stanowi punkt wejścia do zmiany jego wartości. Realizuje mechanizm znaczników kontekstowych (rozdz. 3a): zamiast rozbudowanego panelu ustawień kontekst pracy jest widoczny jako ciąg krótkich etykiet, a każdy element tego ciągu prowadzi jednym kliknięciem do właściwego selektora.

**Forma i waga.** Znikoma — pigułka w linii tekstu, padding 2×9 px, promień pełny (999 px), czcionka 12 px, treść w nawiasach kwadratowych lub bez nich, zależnie od wariantu paska kontekstu. Nigdy nie występuje jako blok samodzielny — zawsze w Pasku kontekstu ze znacznikami (15.9).

**Warianty.**

| Wariant | Treść znacznika | Selektor otwierany po kliknięciu |
|---|---|---|
| Środowisko platformy | `[Danaco Console]` | Lista środowisk platformy |
| Środowisko uruchomieniowe | `[Ubuntu]` | Lista środowisk uruchomieniowych dostępnych dla sesji |
| Repozytorium lub projekt | `[Worktree]` | Lista repozytoriów i projektów w zasięgu bieżącej sesji |
| Model | `[Fable 5]` | Lista modeli bazowych udostępnionych dla wybranego Wykonawcy |
| Poziom wysiłku | `[Ultra]` | Lista poziomów wysiłku Wykonawcy |
| Wykonawca | `[Agent]` | Lista Wykonawców dostępnych w bieżącym kontekście pracy |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Pigułka neutralna z nazwą wartości obowiązującej |
| Hover | Tło `--dn-powierzchnia-2`, kursor wskazuje klikalność |
| Selektor otwarty | Pigułka w wariancie `--zloto`, Panel popover (15.5) zakotwiczony pod znacznikiem |
| Fokus klawiaturowy | Pierścień złoty wspólny dla kontrolek interaktywnych |
| Wartość odziedziczona z poziomu szerszego | Pigułka uzupełniona Dymkiem kontekstowym (11.4) wskazującym poziom pochodzenia wartości; kliknięcie tworzy jawne nadpisanie na bieżącym poziomie |
| Wartość niedostępna w bieżącym kontekście | Znacznik pozostaje w pełni klikalny — ograniczenie sygnalizowane opisowo w treści selektora, nie odebraniem klikalności znacznika |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu listy selektora do czasu odebrania zbioru wartości |

**Zachowanie po interakcji.** Kliknięcie otwiera Panel popover z selektorem wartości; wybór wartości zamyka selektor i podmienia treść znacznika. Zmiana obowiązuje od następnego polecenia lub następnego uruchomienia procesu i nie przerywa działań już uruchomionych.

**Gdzie używany.** Pasek kontekstu Chat Window (14.1) w każdym module każdego środowiska; pasek kontekstu Execution Loop Window (14.2); nagłówek okna operacyjnego w modułach o zmiennym kontekście pracy (Developer, Terminal, Translate, Roundtable).

---

### 15.8. Paleta poleceń

| | |
|---|---|
| **Nazwa** | Paleta poleceń |
| **Klasa bazowa** | `.dn-paleta` |
| **Wzorzec** | `.dn-nakladka` + `.dn-paleta` (`.dn-karta`) z `.dn-input` (zapytanie) i `.dn-karta--pozycja` (wyniki) |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 4 |
| **Sposób wywołania** | Skrót klawiszowy `Ctrl/Cmd + K`; zamknięcie klawiszem Escape lub kliknięciem poza obszarem nakładki |

**Przeznaczenie.** Udostępnia jednym krokiem każde miejsce i każdą funkcję platformy — okno operacyjne, moduł, środowisko, kartę sesji, ustawienie, akcję kontekstu bieżącego oraz operację ekspercką, tryb administracyjny i narzędzie diagnostyczne niskiego poziomu, które nie mają własnej reprezentacji graficznej w interfejsie. Jest realizacją zasady jednego kliknięcia (rozdz. 3a) dla warstwy 4: funkcja niewidoczna w interfejsie pozostaje osiągalna jednym skrótem klawiszowym i jednym wpisanym słowem. Paleta poleceń jest jedynym klawiaturowym punktem wejścia tego rodzaju w całej platformie — obowiązuje w każdym module, w każdym środowisku i w każdej funkcji globalnej pod tym samym skrótem `Ctrl/Cmd + K`. Użytkownik podstawowy nie widzi elementów warstwy 4 — paleta udostępnia wyłącznie funkcje objęte konfiguracją jego roli.

**Forma i waga.** Średnia — nakładka wyśrodkowana w obszarze roboczym, szerokość 480–640 px, maksymalna wysokość 60% wysokości ekranu, tło `--dn-panel`, cień `--dn-cien-lg`, promień `--dn-r-lg`. W spoczynku nie zajmuje żadnego miejsca w układzie i nie ma widocznego wyzwalacza graficznego; jedyną dopuszczalną reprezentacją spoczynkową jest podpowiedź skrótu `Ctrl/Cmd + K` w pasku albo kolumnie stanu środowiska (warstwa 1), której kliknięcie otwiera paletę.

**Elementy składowe.**

| Element | Rola | Forma |
|---|---|---|
| Pole zapytania | Przyjmuje frazę wyszukiwania; przyjmuje fokus od chwili otwarcia | `.dn-input` w nagłówku nakładki, pełna szerokość |
| Znacznik zakresu | Wskazuje zakres wyników obowiązujący dla wpisanej frazy | Znacznik kontekstowy (15.7) po lewej stronie pola zapytania |
| Prefiks trybu | Zawęża zakres wyników jednym znakiem wiodącym: `>` polecenie, `@` symbol, `#` plik lub zasób, `:` wiersz, `/` moduł | Znak wiodący we frazie, odzwierciedlony w znaczniku zakresu |
| Lista wyników | Prezentuje dopasowania w kolejności trafności; dopasowanie rozmyte | Lista `.dn-karta--pozycja`, przewijana niezależnie |
| Wiersz wyniku | Nazwa funkcji lub elementu, jej umiejscowienie w platformie i przypisany jej skrót klawiszowy | Pozycja listy z trzema polami tekstowymi |
| Sekcja ostatnio używanych | Lista funkcji użytych ostatnio, prezentowana przy pustym zapytaniu | Lista `.dn-karta--pozycja` w porządku od najnowszej |
| Wiersz podpowiedzi klawiszy | Objaśnienie klawiszy nawigacji: strzałki, Enter, Escape | Wiersz tekstu pomocniczego zamykający nakładkę, poniżej listy wyników |

**Warianty.**

| Wariant | Zakres wyników | Zastosowanie |
|---|---|---|
| Funkcje platformy | Okna operacyjne, moduły, środowiska, ustawienia, akcje globalne | Przejście do dowolnego miejsca platformy |
| Operacje kontekstu bieżącego | Akcje dostępne dla bieżącego okna, zadania, karty sesji lub zasobu | Wykonanie akcji bez odszukiwania jej w menu |
| Wyszukiwarka funkcji — warstwa 4 | Operacje eksperckie, tryb administracyjny, narzędzia diagnostyczne niskiego poziomu, zrzut stanu Koordynatora | Rola z konfiguracją obejmującą warstwę 4 |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Zamknięta (spoczynek) | Niewidoczna, bez śladu w układzie i bez wyzwalacza graficznego |
| Otwarta, zapytanie puste | Lista funkcji używanych ostatnio, w porządku od najnowszej |
| Otwarta, zapytanie wpisane | Lista wyników dopasowanych rozmyto do zapytania, z nazwą funkcji, jej umiejscowieniem i skrótem klawiszowym |
| Wyniki rozmyte | Trafienia częściowe uszeregowane wg trafności, dopasowane znaki wyróżnione akcentem złotym |
| Brak wyników | Pusty stan (10.7) z opisem wskazującym Chat Window (14.1) jako kanał polecenia w języku naturalnym |
| Fokus klawiaturowy | Fokus w polu zapytania od chwili otwarcia; klawisze strzałek przenoszą wskazanie między wynikami, Enter wykonuje wskazany wynik |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu listy wyników do czasu odebrania dopasowań |
| Błąd | Komunikat blokowy (11.3) wariantu `--blad` w miejscu listy, z akcją ponowienia |

**Zachowanie po interakcji.** Otwarcie skrótem `Ctrl/Cmd + K` ustawia fokus w polu zapytania. Wskazanie wyniku i jego zatwierdzenie zamyka nakładkę i wykonuje funkcję albo przenosi do wskazanego miejsca platformy. Nakładka nie przerywa procesów w toku i nie zmienia zawartości kolumn układu do chwili wykonania wybranej funkcji.

**Relacja do wyszukiwarki funkcji.** Wyszukiwarka funkcji nie jest odrębnym komponentem — jest wariantem zakresu wyników palety poleceń obejmującym warstwę 4: operacje eksperckie, tryb administracyjny i narzędzia diagnostyczne niskiego poziomu. Platforma ma jedną nakładkę klawiaturową i jeden skrót `Ctrl/Cmd + K`; nazwa „wyszukiwarka funkcji" oznacza w dokumentacji sposób dostępu do warstwy 4 realizowany przez tę nakładkę, nigdy drugą nakładkę o własnym skrócie i własnej formie. Nazwa komponentu w katalogu, w makietach i w tabelach elementów okna brzmi **Paleta poleceń**.

**Gdzie używany.** Dostępna w każdym module każdego środowiska i w każdej funkcji globalnej, niezależnie od otwartego okna operacyjnego i bieżącej karty sesji: powłoki środowisk TalkIn, WorkSpace, CodeStudio i MultitaskingAI, wszystkie okna operacyjne modułów, okno konfiguracji, okno Ustawień, funkcja globalna Mobile.

---

### 15.9. Pasek kontekstu ze znacznikami

| | |
|---|---|
| **Nazwa** | Pasek kontekstu ze znacznikami |
| **Wzorzec** | Rząd `.dn-plakietka--interaktywna` (Znacznik kontekstowy, 15.7) w jednym wierszu |
| **Kategoria** | I — Ujawnianie funkcjonalności |
| **Status** | Wzorzec złożony |
| **Warstwa widoczności** | 1 |
| **Sposób wywołania** | Widoczny bez interakcji u wierzchołka kolumny, której kontekst opisuje |

**Przeznaczenie.** Zbiera w jednym wierszu wszystkie wymiary kontekstu, w którym pracuje Wykonawca, i czyni je jednocześnie widocznymi oraz zmienialnymi. Zastępuje panel ustawień kontekstu: zamiast otwierania osobnego okna użytkownik odczytuje kontekst wprost z paska, a każdą wartość zmienia jednym kliknięciem odpowiedniego znacznika.

**Forma i waga.** Znikoma — pojedynczy wiersz znaczników u wierzchołka kolumny, wysokość jednego wiersza tekstu z odstępem `--dn-od-2` powyżej i poniżej, znaczniki rozdzielone odstępem `--dn-od-1`. Pasek nie ma własnego tła ani obrysu; od treści kolumny oddziela go Separator (13.2).

**Warianty.**

| Wariant | Zestaw znaczników | Umiejscowienie |
|---|---|---|
| Kontekst rozmowy | `[Środowisko platformy] [Środowisko uruchomieniowe] [Repozytorium] [Model] [Poziom wysiłku]` | Wierzchołek kolumny Chat Window (14.1) |
| Kontekst pętli wykonawczej | `[Zlecenie] [Koordynator] [Wykonawcy] [Przebieg]` | Wierzchołek kolumny Execution Loop Window (14.2) |
| Kontekst okna operacyjnego | Znaczniki właściwe modułowi — projekt, zbiór, język, gałąź | Nagłówek okna operacyjnego |

**Stany.**

| Stan | Realizacja wizualna |
|---|---|
| Domyślny | Rząd znaczników z wartościami obowiązującymi |
| Znacznik w stanie hover lub z otwartym selektorem | Stany wg karty Znacznika kontekstowego (15.7) |
| Przepełnienie wiersza | Znaczniki nadmiarowe zbierane w jeden znacznik zbiorczy `+N ▼`, którego kliknięcie otwiera Panel popover (15.5) z pełną listą; drugi poziom rozwinięcia nie występuje |
| Kontekst niepełny | Znacznik wymiaru bez ustalonej wartości prezentuje nazwę wymiaru zamiast wartości i zachowuje pełną klikalność |
| Ładowanie | Wskaźnik ładowania (11.5) w miejscu wiersza do czasu odebrania kontekstu sesji |
| Błąd | Ostatni znany kontekst zachowany w pasku, błąd odczytu sygnalizowany Powiadomieniem (11.2) |

**Zachowanie po interakcji.** Pasek sam w sobie jest bierny — zachowanie należy do znaczników. Zmiana wartości dowolnego znacznika aktualizuje pasek natychmiast i obowiązuje od następnego polecenia lub następnego uruchomienia procesu.

**Gdzie używany.** Wierzchołek kolumny Chat Window i kolumny Execution Loop Window w każdym module każdego środowiska; nagłówki okien operacyjnych w modułach o zmiennym kontekście pracy.

---

## 16. Macierz komponent × moduł / okno

Poniższa tabela zestawia, w których modułach i oknach globalnych każdy komponent występuje najsilniej — nie jest wyczerpującym spisem (Przycisk i Ikona systemowa występują praktycznie wszędzie), lecz wskazuje miejsca o największym nasileniu lub o znaczeniu wzorcotwórczym dla danego komponentu.

| Komponent | Studio | Workspace | Automations | Browser | Research | Library | Translate | Roundtable | Design | Assistant | Terminal | Developer | Diagnostics | Apps | Agents | MultitaskingAI | Okno konfiguracji |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Przycisk | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| Przycisk ikonowy | ● | — | ●● | — | — | — | — | — | — | ● | ●● | — | ● | — | — | ●● | — |
| Pole tekstowe | ● | ●● | — | — | — | — | ●● | — | ●● | — | — | — | — | — | ●● | — | ● |
| Lista rozwijana | — | — | ● | — | — | — | — | — | — | — | — | — | — | — | ●● | — | ● |
| Przełącznik | — | — | ● | — | — | — | — | — | — | — | — | — | — | — | ●● | — | ●●● |
| Checkbox / Radio | — | — | — | — | — | — | — | — | — | — | — | — | — | — | — | — | ●● |
| Tabela | — | — | ●● | — | — | — | — | — | — | — | ● | — | ●● | — | ●● | ●● | ●● |
| Karta (pozycja listy) | ●● | ●● | — | — | ●● | ●● | — | — | ●● | — | — | ● | — | — | — | ●● | ●● |
| Panel | ●● | ●● | — | ●● | ●● | ●● | ●● | ●● | ●● | ● | — | ● | ●● | ●● | ●● | — | — |
| Plansza | — | — | — | — | — | — | — | — | ●●● | — | — | — | — | — | — | — | — |
| Pigułka / Kropka statusu | — | — | ●● | — | ● | ●●● | — | — | — | — | — | — | ● | — | — | ●● | ● |
| Blok kodu | ● | — | — | — | — | — | — | — | — | — | ●●● | ●● | ●● | ● | — | — | ● |
| Menu kontekstowe | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | — |
| Modal | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●●● | ● | ● |
| Powiadomienie | ● | ● | ●● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●●● |
| Komunikat blokowy | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ●● | ● | ● | ●● | ● |
| Dymek kontekstowy | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●●● |
| Wskaźnik ładowania | ● | ● | ●● | ● | ● | ● | ● | ● | ●● | ● | — | ●● | — | ●● | ● | ●● | — |
| Awatar | ● | ● | ● | ● | ● | ● | ● | ●●● | ● | ● | ● | ● | ● | ● | ●● | ●● | — |
| Zakładka | — | — | — | — | — | — | — | ●● | — | — | ●●● | — | — | — | — | ● | ●● |
| Karta sesji | ●●● | ●●● | — | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | — |
| Pasek górny | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | — |
| Boczna nawigacja / Panel orkiestracji | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● |
| Pusty stan | — | — | ● | — | ●● | ● | — | — | ● | — | — | — | ● | — | — | — | — |
| Avatar Always On Display | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●●● | ● |
| Ikona systemowa | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| Separator | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| Chat Window | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | ●●● | — |
| Execution Loop Window | ●● | ●● | ●●● | ● | ●● | ● | ● | ●● | ● | ● | ●● | ●●● | ●●● | ● | ●● | ●●● | — |
| Element zbiorczy z rozwinięciem | ● | ● | ●● | ● | ● | ● | ● | ● | ● | ● | ●● | ●● | ●● | ● | ● | ●●● | ● |
| Menu progresywne | ●● | ●● | ●● | ● | ●● | ● | ●● | ●●● | ● | ●● | ●● | ●●● | ● | ● | ●●● | ●●● | ● |
| Menu kebab (⋮) | ● | ● | ●● | ● | ●● | ●● | ● | ● | ● | ● | ● | ● | ●● | ● | ●● | ●● | ● |
| Menu hamburger (☰) | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| Panel popover | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ●● | ●● | ●●● |
| Panel wysuwany | ●●● | ●● | ● | ●● | ●● | ●● | ● | ● | ●●● | ● | ●● | ●●● | ●●● | ● | ●● | ●● | ● |
| Znacznik kontekstowy | ●● | ●● | ● | ● | ● | ● | ●● | ●● | ● | ● | ●● | ●●● | ● | ● | ●● | ●●● | ● |
| Paleta poleceń | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ●● | ● | ● | ● | ● |
| Pasek kontekstu ze znacznikami | ●● | ●● | ●● | ●● | ●● | ●● | ●● | ●● | ●● | ●● | ●● | ●●● | ●● | ●● | ●● | ●●● | — |

Legenda: `—` nie występuje w charakterystyczny sposób · `●` obecny · `●●` znaczące nasilenie · `●●●` największe w platformie nasilenie danego komponentu.

---

## 17. Zgodność z zasadami nadrzędnymi platformy

Katalog komponentów pozostaje w zgodności z czterema zasadami nadrzędnymi funkcjonalności platformy (Koncepcja platformy, rozdział 14) oraz z regułą warstw widoczności (rozdz. 3a).

| Zasada nadrzędna | Realizacja w katalogu komponentów |
|---|---|
| Pełna kompozycyjność | Każdy komponent jest dostępny w każdym oknie, module i środowisku, w którym jego zastosowanie ma sens funkcjonalny — katalog nie wprowadza wariantów wykluczających się wzajemnie ani komponentów zastrzeżonych dla pojedynczego modułu |
| Pełna konfigurowalność | Motyw jasny/ciemny, kolejność sekcji nawigacji i wartości wyjściowe przełączników (rozdz. 8.3) są ustawieniami zmienialnymi z okna konfiguracji — żaden komponent katalogu nie koduje zachowania niedostępnego do zmiany tam, gdzie odpowiednia warstwa konfiguracji istnieje |
| Jawność i konfigurowalność zależności | Dymek kontekstowy (rozdz. 11.4) czyni każdy element konfiguracji jawnie opisanym w miejscu użycia |
| Rozszerzenie orkiestracji | Tabela, Plakietka statusu i Karta (rozdz. 10) oraz Execution Loop Window (rozdz. 14.2) skalują się wraz ze wzrostem liczby ról, agentów i procesów równoległych w środowisku MultitaskingAI, bez definiowania nowych wzorców przy zmianie skali |
| Reguła warstw widoczności | Każda karta komponentu (rozdz. 7–15) podaje warstwę widoczności i sposób wywołania; komponenty kategorii I (rozdz. 15) realizują mechanizmy ukrywania funkcjonalności, a zasada jednego kliknięcia obowiązuje w każdym z nich — drugi poziom rozwinięcia nie występuje w katalogu |
| Dwa kanały komunikacji operacyjnej | Chat Window (rozdz. 14.1) i Execution Loop Window (rozdz. 14.2) są komponentami pierwszoplanowymi katalogu, o stałym miejscu w kanonicznym układzie kolumnowym każdego modułu i każdego środowiska |

Zasada zero blokad, obowiązująca dla wszystkich kart komponentów w rozdziałach 7–15: żadna kontrolka interaktywna katalogu (Przycisk, Przycisk ikonowy, pola formularzy, Zakładka, Karta interaktywna, wiersz Tabeli, pozycja Menu kontekstowego, Boczna nawigacja) nie traci klikalności ani fokusowalności w reakcji na niespełniony warunek, brak uprawnienia czy trwający proces walidacji — „stan wyłączony” (krycie 50–60% + kursor „niedozwolone” + zablokowany hover/fokus) nie występuje w katalogu jako mechanizm odbierający interakcję. Niedostępność lub niegotowość jest sygnalizowana wyłącznie opisowo — tekstem pomocniczym, Plakietką, Dymkiem kontekstowym (rozdz. 11.4) lub komunikatem zwracanym po próbie kliknięcia — nigdy przez usunięcie kontrolki z obiegu interakcji; walidacja ma charakter ostrzegawczy i nigdy nie blokuje zapisu ani wykonania. Dotyczy to również stanu błędu: wartość wyjściowa „wyłączony” Przełącznika (rozdz. 8.3) w macierzy izolacji jest punktem startowym możliwym do zmiany na każdym poziomie zasięgu, nie ograniczeniem, a wartość odziedziczona z poziomu szerszego pozostaje jawnie nadpisywalna w dowolnym widoku, nigdy zablokowana do edycji. Komponent Przycisku wariantu sygnalizującego wagę akcji (rozdz. 7.1) oznacza wizualnie ciężar gatunkowy akcji, nie uniemożliwia jej wykonania — potwierdzenie Modalem (rozdz. 11.1) jest wzorcem dobieranym do konkretnej akcji, nie obowiązkową bramką.

---

## 18. Słowniczek pojęć

| Pojęcie | Definicja |
|---|---|
| Komponent zdefiniowany | Komponent katalogu mający odrębną, gotową klasę w `components.css` |
| Wzorzec złożony | Komponent katalogu nazwany i opisany na poziomie produktowym w źródłach, zbudowany ze złożenia kilku komponentów zdefiniowanych, bez własnej odrębnej klasy CSS |
| Token semantyczny | Zmienna CSS/JSON o nazwie opisującej rolę (np. `--dn-tekst`, `--dn-akcent`), a nie surową wartość barwy — jedyny dopuszczalny sposób odwołania się komponentu do koloru |
| Token prymitywny | Surowa wartość skali barwnej (np. `--dn-granat-500`) — podstawa tokenów semantycznych, nigdy niewykorzystywana wprost w komponencie |
| Forma i waga | W niniejszym katalogu: opis wymiarów, kształtu i wizualnej masy komponentu — czy jest to znikoma ikonka, mały przycisk, średnie pole czy duży panel |
| Stan komponentu | Jedna z rozpoznawanych sytuacji wizualnych komponentu: domyślny, hover, aktywny, niedostępny w danym kontekście (sygnalizowany opisowo, bez odbierania klikalności — zob. rozdz. 17), ładowanie, błąd — dokumentowana per komponent w rozdziałach 7–15, z jawnym „nie dotyczy” tam, gdzie dany stan nie ma zastosowania |
| Karta sesji | Zob. rozdz. 9.3 — poziomy element paska kart reprezentujący jedną, samodzielną przestrzeń roboczą |
| Panel orkiestracji | Zob. rozdz. 9.2 — odpowiednik bocznej nawigacji modułów właściwy środowisku MultitaskingAI |
| Always On Display | Globalny agent towarzyszący, zob. rozdz. 12.2 oraz `funkcje-globalne/always-on-display.md` |
| Chat Window | Główne okno komunikacji, kanał Użytkownik ↔ Wykonawca — zob. rozdz. 14.1; lewa kolumna obszaru roboczego, stała, pełna wysokość |
| Execution Loop Window | Okno pętli wykonawczej, kanał Koordynator ↔ Wykonawca — zob. rozdz. 14.2; kolumna sąsiadująca z Chat Window |
| Użytkownik | Rola zlecająca i zatwierdzająca działania |
| Koordynator | Komponent orkiestrujący platformy — dekomponuje zlecenie, przydziela i nadzoruje zadania |
| Wykonawca | AI, agent lub system wykonawczy realizujący zadania |
| Warstwa widoczności | Jedna z czterech warstw ujawniania funkcjonalności (rozdz. 3a), podawana w karcie każdego komponentu: 1 — zawsze widoczna, 2 — widoczna na żądanie, 3 — rozwinięcie kontekstowe, 4 — funkcja ekspercka |
| Sposób wywołania | Pole karty komponentu wskazujące, jaka interakcja ujawnia komponent — brak interakcji, kliknięcie wyzwalacza, skrót klawiszowy albo polecenie języka naturalnego |
| Układ pionowy (podział lewa–prawa) | Kanoniczne rozmieszczenie okien platformy w kolumnach sąsiadujących poziomo; regulacji podlega wyłącznie szerokość kolumn |
| Okno pomocnicze | Okno otwierane jako rozszerzenie boczne — kolejna kolumna po prawej stronie obszaru roboczego (Panel wysuwany, rozdz. 15.6) |
| Paleta poleceń | Jedyna klawiaturowa nakładka platformy, wywoływana skrótem `Ctrl/Cmd + K` — zob. rozdz. 15.8 |
| Wyszukiwarka funkcji | Wariant zakresu wyników palety poleceń obejmujący warstwę 4; sposób dostępu do funkcji eksperckich, nie odrębny komponent — zob. rozdz. 15.8 |
| Centrum powiadomień | Jedyny mechanizm powiadamiania platformy — kolumna boczna z rejestrem zdarzeń, zob. rozdz. 11.6 |
| Magistrala kontekstu | Mechanizm przekazywania artefaktu między modułami — zob. `interfejs-uzytkownika/przeplyw-okien.md`, rozdz. 6a |
| Okno operacyjne | Pojedynczy element zestawu okien modułu — jednostka opisu Specyfikacji okien operacyjnych, nie przedmiot niniejszego katalogu (przedmiotem są komponenty, z których okno jest złożone) |

---

## Załącznik A. Ściągawka wymiarów i tokenów

Wartości poniżej pochodzą z `tokens.css` — źródła prawdy dla warstwy webowej (zob. Załącznik B w sprawie rozbieżności z `tokens.json` i `interfejs-uzytkownika/system-wizualny.md`).

**Tabela — skala typografii.**

| Token | Wartość | Typowe zastosowanie w komponentach katalogu |
|---|---|---|
| `--dn-fs-xs` | 12 px | Plakietka, Dymek kontekstowy, nagłówek kolumny Tabeli, Ikona systemowa (etykieta) |
| `--dn-fs-sm` | 13 px | Przycisk (domyślny), Zakładka, komórka Tabeli, Awatar (inicjał domyślny) |
| `--dn-fs-base` | 15 px | Pole tekstowe, treść Panelu, treść czatu |
| `--dn-fs-md` | 16 px | Tekst wyróżniony |
| `--dn-fs-lg` | 18 px | Tytuł Karty, Awatar (inicjał, wariant duży) |
| `--dn-fs-xl` | 22 px | Tytuł Modala, nagłówek Paska górnego |
| `--dn-fs-2xl` | 28 px | Tytuł Karty środowiska (wariant mniejszy) |
| `--dn-fs-3xl` | 36 px | Tytuł Karty środowiska (wariant większy) |
| `--dn-fs-display` | 46 px | Stopień ekspozycyjny, poza bieżącym katalogiem komponentów |

**Tabela — skala odstępów (jednostka 4 px).**

| Token | `-1` | `-2` | `-3` | `-4` | `-5` | `-6` | `-8` | `-10` | `-12` | `-16` |
|---|---|---|---|---|---|---|---|---|---|---|
| Wartość | 4 px | 8 px | 12 px | 16 px | 20 px | 24 px | 32 px | 40 px | 48 px | 64 px |

**Tabela — promienie narożników i ich przypisanie do komponentów.**

| Token | Wartość | Komponenty katalogu |
|---|---|---|
| `xs` | 4 px | Blok kodu |
| `sm` | 8 px | Przycisk, Przycisk ikonowy, Pole tekstowe, Lista rozwijana, Menu kontekstowe (pozycja) |
| `md` | 10 px | Przycisk (wariant duży), Menu kontekstowe (pojemnik), Karta (mniejsza) |
| `lg` | 14 px | Modal, Karta (duża) |
| `xl` | 20 px | Karta środowiska strony głównej |
| `pill` | 999 px | Pigułka/Plakietka, Przełącznik (tor) |

**Tabela — wymiary stałe kluczowych komponentów.**

| Komponent | Wymiar |
|---|---|
| Przycisk ikonowy | 36×36 px, ikona wewnętrzna 18×18 px |
| Awatar | 36×36 px (domyślny) · 28×28 px (mały) · 48×48 px (duży) |
| Przełącznik | tor 40×22 px, suwak 18×18 px |
| Checkbox / Radio | 16×16 px |
| Pasek górny | wysokość 56 px (`--dn-topbar-h`) |
| Modal | maks. szerokość 560 px, maks. wysokość 92% wysokości ekranu |
| Powiadomienie | szerokość 260–440 px |
| Ikona systemowa | siatka źródłowa 24×24 px, render 14/16/18/20/24 px |

---

## Załącznik B. Uwaga o rozbieżności źródeł systemu wizualnego

Dwie grupy plików źródłowych systemu wizualnego opisują różne wartości w dwóch punktach. Odnotowanie tej rozbieżności ma znaczenie praktyczne: designer i deweloper korzystający z różnych plików źródłowych przyjęliby niezależnie od siebie dwie różne skale typografii lub dwie różne palety trybu ciemnego dla tego samego komponentu.

| Punkt rozbieżności | `tokens.json` oraz `interfejs-uzytkownika/system-wizualny.md` (dokument) | `tokens.css` (plik faktycznie ładowany przez aplikację) |
|---|---|---|
| Skala typografii | `11 · 12 · 14(base) · 15 · 17 · 21 · 26 · 34 · 44` px | `12 · 13 · 15(base) · 16 · 18 · 22 · 28 · 36 · 46` px — skala podniesiona o 1–2 px na każdym stopniu |
| Paleta trybu ciemnego | Oparta na skali granatu (tło `#0B1120`, powierzchnia `#16223C` — te same odcienie granatu co marka) | Oparta na odrębnej skali „grafit” — niemal-czarne odcienie neutralne (tło `#0C0C0E`, powierzchnia `#17171A`), z granatem i złotem zachowanymi wyłącznie jako kolor marki i akcentu, nie jako tło |

`tokens.css` jest — zgodnie z jawnym określeniem w `interfejs-uzytkownika/system-wizualny.md`, rozdz. 1.2 — źródłem prawdy dla warstwy webowej. **Niniejszy katalog przyjmuje wartości `tokens.css` jako obowiązujące** w każdym miejscu, w którym w kartach komponentów (rozdziały 7–15) i w Załączniku A podano konkretną wartość pikselową lub barwną.

Dodatkowa obserwacja porządkowa: Pasek górny (rozdz. 9.1) ma tło zawsze ciemne (`--dn-pasek-tlo`) niezależnie od aktywnego motywu — token ten jest odrębny od `--dn-marka` i w motywie jasnym przyjmuje wartość niemal identyczną z tłem całej aplikacji w motywie ciemnym. Jest to spójne w obu grupach plików źródłowych i nie stanowi rozbieżności — odnotowane tu wyłącznie jako fakt istotny dla poprawnego odczytania kolumny „Forma i waga” Paska górnego.

---

## Załącznik C. Pełna lista klas CSS komponentów

Referencja 1:1 z `components.css`, do użytku dewelopera przy implementacji.

| Komponent (nazwa w katalogu) | Klasa bazowa | Warianty | Rozdział niniejszego dokumentu |
|---|---|---|---|
| Przycisk | `.dn-btn` | `--glowny` `--zloty` `--zarys` `--duch` `--blad` `--sm` `--lg` | 7.1 |
| Przycisk ikonowy | `.dn-btn-ikona` | — | 7.2 |
| Pole tekstowe | `.dn-pole` `.dn-input` `.dn-textarea` | `--blad` | 8.1 |
| Lista rozwijana | `.dn-select` | `--blad` | 8.2 |
| Przełącznik | `.dn-suwak` | — | 8.3 |
| Pole wyboru / opcja jednokrotna | `.dn-check` | — | 8.4 |
| Pasek górny | `.dn-pasek` | `.dn-pasek-marka` `.dn-pasek-szukaj` `.dn-rozpychacz` | 9.1 |
| Boczna nawigacja / Panel orkiestracji | `.dn-karta--pozycja` (kolumna) | — | 9.2 |
| Karta sesji | `.dn-karta--pozycja` (pasek poziomy) | — | 9.3 |
| Zakładka | `.dn-zakladki` `.dn-zakladka` | `--pigulki` | 9.4 |
| Menu kontekstowe | `.dn-karta` + `.dn-karta--pozycja` (pływające) | — | 9.5 |
| Karta | `.dn-karta` `.dn-karta-naglowek` `.dn-karta-tytul` `.dn-karta-tresc` `.dn-karta-stopka` | `--interaktywna` `--akcent` `--pozycja` | 10.1 |
| Panel | region z `.dn-panel` (token tła) | — | 10.2 |
| Plansza | obszar główny, bez klasy dedykowanej | — | 10.3 |
| Tabela | `.dn-tabela` `.dn-tabela-owijka` | `--paski` | 10.4 |
| Pigułka i plakietka statusu | `.dn-plakietka` `.dn-kropka` | `--zloto` `--sukces` `--ostrz` `--blad` `--info` `--wersaliki` `--stan` | 10.5 |
| Blok kodu | `.dn-kod` | — | 10.6 |
| Pusty stan | `.dn-pusty-stan` | — | 10.7 |
| Modal | `.dn-nakladka` `.dn-modal` `.dn-modal-naglowek` `.dn-modal-tytul` `.dn-modal-tresc` `.dn-modal-stopka` | — | 11.1 |
| Powiadomienie | `.dn-toasty` `.dn-toast` | `--sukces` `--blad` `--ostrz` | 11.2 |
| Komunikat blokowy | `.dn-alert` | `--info` `--sukces` `--ostrz` `--blad` | 11.3 |
| Dymek kontekstowy | `.dn-tooltip` `.dn-tooltip-tresc` | — | 11.4 |
| Wskaźnik ładowania | `.dn-spinner` | — | 11.5 |
| Awatar | `.dn-awatar` `.dn-awatar-stan` | `--zloto` `--sm` `--lg` `--kwadrat` | 12.1 |
| Avatar Always On Display | `.dn-awatar` (pozycjonowanie fixed, poza `components.css`) | — | 12.2 |
| Ikona systemowa | `.dn-ikona` | `--sm` `--lg` | 13.1 |
| Separator | `.dn-separator` | — | 13.2 |
| Chat Window | złożenie `.dn-karta--pozycja` `.dn-awatar` `.dn-kod` `.dn-input` `.dn-btn-ikona` (kolumna) | — | 14.1 |
| Execution Loop Window | złożenie `.dn-karta--pozycja` `.dn-tabela` `.dn-plakietka--stan` `.dn-btn` (kolumna) | — | 14.2 |
| Element zbiorczy z rozwinięciem | `.dn-btn--zarys` + `.dn-karta` `.dn-karta--pozycja` | — | 15.1 |
| Menu progresywne | `.dn-btn--duch` + `.dn-karta--pozycja` | — | 15.2 |
| Menu kebab | `.dn-btn-ikona` (ikona `wiecej`) + `.dn-karta` | — | 15.3 |
| Menu hamburger | `.dn-btn-ikona` (ikona `menu`) + `.dn-karta` | — | 15.4 |
| Panel popover | `.dn-karta` (pływająca, zakotwiczona przy wyzwalaczu) | — | 15.5 |
| Panel wysuwany | region `.dn-panel` (kolumna boczna) | — | 15.6 |
| Znacznik kontekstowy | `.dn-plakietka` | `--interaktywna` `--zloto` | 15.7 |
| Paleta poleceń | `.dn-paleta` (`.dn-nakladka` + `.dn-karta` `.dn-input` `.dn-karta--pozycja`) | — | 15.8 |
| Centrum powiadomień | `.dn-panel-wysuwany` + `.dn-karta--pozycja` + `.dn-plakietka--stan` | — | 11.6 |
| Pasek kontekstu ze znacznikami | rząd `.dn-plakietka--interaktywna` | — | 15.9 |

---

*Koniec dokumentu. Danaco Console — Katalog komponentów interfejsu, wersja 2.0.*

---
*Danaco Console — AI Workspace OS · v2.0*
*© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — [LICENSE](LICENSE). Kontakt: support@danaco-group.pl*
