# Plan etapów

Kolejność budowy Danaco Console. Podstawą jest pozycja 4
[rejestru decyzji](decyzje.md): powstaje jedno okno modułowe, a następnie pełny
przekrój pionowy aplikacji na tym jednym module. Komplet prototypów i komplet
opracowań przed budową nie powstaje.

Etap kolejny startuje po zamknięciu bramki etapu poprzedniego przez Właściciela.
Zasady bramek opisuje [ustrój budowy](ustroj-budowy.md).

| Etap | Przedmiot | Prowadzi | Stan |
|---|---|---|---|
| 1 | Łańcuch od instalacji po okno modułowe Studio | Właściciel, poprawki zgłaszane pojedynczo | w toku |
| 2 | Aplikacja pełna kształtem, działająca w module Studio | wykonawcy w terenach | wstrzymany do zamknięcia etapu 1 |
| 3 | Kolejne moduły na działającym przekroju | wykonawcy w terenach | wstrzymany do zamknięcia etapu 2 |

## Etap 1 — Łańcuch od instalacji po okno modułowe

**Cel.** Doprowadzić do postaci przyjętej ciąg okien od instalacji po pierwsze
okno modułowe. Ciąg ten rozstrzyga kompozycję, gęstość i zachowanie interfejsu
platformy — jest wzorcem dla wszystkich następnych okien nie dlatego, że
zostanie skopiowany, lecz dlatego, że ustala reguły, których pozostałe będą się
trzymać.

**Prowadzi Właściciel.** Kompozycja okna nie jest pracą odtwórczą i nie da się
jej przekazać zleceniem — zlecenie opisuje zasadę i kryterium sprawdzalne,
a wyglądu sprawdzić się w ten sposób nie da. Orkiestracja wieloagentowa w tym
etapie nie jest stosowana.

**Zakres — pięć plików, około dziesięciu makiet.** Makieta to jeden stan okna,
nie osobny plik; przepływ niesie swoje stany wewnątrz jednego pliku klikalnego.
Zakres jest zamknięty: okno spoza tego wykazu nie powstaje w tym etapie.

| Plik | Makiety |
|---|---|
| `platformowe/instalator.html` | kroki instalacji wraz z konserwacją pakietu |
| `przeplyw/przeplyw-wejscia.html` | łączenie, rejestracja, logowanie, odzyskiwanie, przygotowanie środowiska |
| `przeplyw/centrum-dowodzenia.html` | strona główna wraz z animacją przejścia do środowiska |
| `srodowiska/talkin-przedsionek.html` | przedsionek środowiska TalkIn |
| `moduly/studio.html` | widok pojedynczy z czatem; widok dzielony czat–edytor wraz z drugim czatem w oknie bocznym; przejścia między układami |

Instalator wraz z przepływem wejścia daje około pięciu makiet, centrum
z animacją i przedsionkiem trzy, Studio dwie.

**Prototypy są klikalne i przechodzą przez układy.** Prototyp Studio nie
przedstawia dwóch obrazów, lecz pozwala przejść wszystkie ułożenia widoku:
zestawienie kart w widok dzielony i rozłączenie go, przeniesienie karty do okna
bocznego i powrót, przełączanie między parami. Mechanizm ma być pokazany
w działaniu, nie zapowiedziany.

**Tryb pracy.** Właściciel zgłasza poprawki pojedynczo sesji kontroli jakości
designu i ocenia wynik, aż do wersji przyjętej. Nie jest to orkiestracja
wieloagentowa ani zlecenie hurtowe — poprzednie podejście wykazało, że
kompozycji nie da się przekazać zleceniem opisującym zasadę.

**Rozstrzygnięcia z pętli poprawek.** Poprawka jednorazowa zostaje w pliku
prototypu. Rozstrzygnięcie obowiązujące kolejne okna — podział wstążek, piętra
kart, przypisanie wskaźników — idzie do [rejestru decyzji](decyzje.md), bo
transkrypt sesji jest ulotny, a kolejne okna muszą wiedzieć, co ustalono przy
poprzednich.

**Miejsce pracy.** Gałąź `teren/prototypy`, drzewo `~/robocze/prototypy`.
Prototypy okien leżą w `design/05-okna/`. Zrzuty stanu wyjściowego stoją poza
repozytorium w `~/robocze/zrzuty-probki/`.

**Ograniczenie warstwy wspólnej.** Warstwa o zasięgu ogólnym — żetony, fundament,
komponenty, rama, stanowisko, karty okna — jest wspólna dla wszystkiego, co
powstanie później. Zmiana w niej jest decyzją, nie poprawką okna. Potrzeba takiej
zmiany, ujawniona przy pracy nad Studio, zostaje odnotowana w
[rejestrze terenów](rejestr-terenow.md) w sekcji zgłoszeń i wchodzi do etapu 2
wraz z kodem, który z tej warstwy wyrośnie.

**Kryteria zamknięcia etapu.** Rozstrzyga Właściciel oceną okna. Jedno kryterium
jest jednak sprawdzalne maszynowo i obowiązuje: okno nie zawiera wartości
projektowych wpisanych liczbowo — barwa, odstęp i krój pochodzą z żetonów.

## Etap 2 — Aplikacja pełna kształtem, działająca w jednym module

**Cel.** Zbudować aplikację o pełnym kształcie, w której **działa jeden moduł**.
Produkt ma dać się przejść w całości — cztery środowiska, przedsionki, szyna,
okna platformowe — ale rzeczywistą pracę od interfejsu do danych wykonuje
wyłącznie Studio.

**Dwa stany okna — trzeciego nie ma.** Okno albo działa, albo jest wyłącznie
zapowiedziane. Stanu pośredniego — okna, które się otwiera i pokazuje kształt
bez działania — **nie wprowadzamy**.

| Stan | Co znaczy | Które okna |
|---|---|---|
| **działające** | pełny przekrój od interfejsu do zapisu i z powrotem | łańcuch etapu 1: instalacja, wejście, centrum dowodzenia, przedsionek TalkIn, Studio |
| **zapowiedziane** | ikona i pozycja w nawigacji są; okno nie powstaje, a próba wejścia odpowiada komunikatem i nie prowadzi dalej | wszystkie pozostałe: trzy przedsionki, czternaście modułów, okna platformowe spoza łańcucha |

**Jedno środowisko i jeden moduł wystarczą, żeby sprawdzić pion.** Studio należy
do TalkIn i do WorkSpace; przekrój buduje się na drodze TalkIn → Studio i tylko
na niej. Dopiero gdy pion działa, a budowa prowadzona jest sprawnie, wchodzi
droga przez WorkSpace i kolejne moduły. Sprawdzenie jednej drogi rozstrzyga
o rdzeniu, kontrakcie i kanale — powtarzanie tego na drugiej przed sprawdzeniem
pierwszej niczego nie ujawnia, a mnoży pracę.

Centrum dowodzenia niesie **komplet ikon modułów**, żeby produkt miał swój
prawdziwy kształt. Przejście prowadzi wyłącznie do Studio. Każde inne
zatrzymuje się na komunikacie mówiącym, co się w tym miejscu wydarzy.

**Zapora przed atrapą — reguła bezwzględna.** Okno zapowiedziane **nie powstaje
jako pusta skorupa i nie przedstawia danych, których nie ma**. Żadnych
przykładowych sesji, liczników, nazw plików ani metryk wpisanych po to, żeby
okno wyglądało na pełne.

Powód jest praktyczny, nie estetyczny: **wypełniacza nikt później nie
podmienia**. Raz wpuszczony zostaje do końca życia produktu, a z czasem
przestaje być rozpoznawalny jako wypełniacz i zaczyna uchodzić za wymaganie.
Poprzednie podejście upadło między innymi na tym — strona główna pokazywała
„siedem sesji, pięć czynnych" przy zerze sesji rzeczywistych, obok banera
pierwszego wejścia; dwa wykluczające się stany, oba zmyślone.

Sposób odmowy jest już rozstrzygnięty w warstwie projektowej: nagłówek
`design/zasoby/css/komponenty.css` stanowi „zasada zero blokad — żaden wariant
nie odbiera klikalności; niegotowość komunikuje się opisem albo komunikatem".
Ikona modułu niegotowego pozostaje klikalna i odpowiada zdaniem — nie jest
wyszarzona i nie prowadzi do pustego okna.

**Zakres.** Kontrakt jako jedyne źródło prawdy typów i komunikatów; rdzeń Go;
kanał komunikacyjny wraz z korelacją żądanie–odpowiedź i wznowieniem po zerwaniu;
klient TypeScript budujący okno Studio z kontraktu; powłoka Tauri wraz z cyklem
życia rdzenia; trwałość w SQLite; uwierzytelnienie w zakresie, którego wymaga
przepływ wejścia.

**Kolejność wewnątrz etapu.** Model danych i kontrakt, potem pionowy przepływ
krytyczny jednego działania modułu od interfejsu do zapisu, potem kolejne
działania, na końcu wykończenie. Interfejs powstaje na działającym przepływie
danych, nie odwrotnie.

**Dokumentacja powstaje wraz z przekrojem.** Nie przed nim i nie po nim. Opisuje
to, co działa, a nie to, co zamierzone — i obejmuje wyłącznie zakres, który
przekrój obsługuje. Skład tej dokumentacji rozstrzyga Właściciel przed otwarciem
etapu; jest to pozycja otwarta rejestru decyzji.

**Kryteria zamknięcia etapu.**

1. Aplikacja uruchamia się jako produkt złożony, nie jako zestaw procesów
   deweloperskich.
2. Jedno pełne działanie modułu Studio przechodzi od kliknięcia w oknie do
   zapisu w bazie i z powrotem, z przytoczonym wynikiem uruchomienia.
3. Zerwanie połączenia i powrót nie gubią stanu sesji — wykazane próbą.
4. Typy po stronie rdzenia i klienta pochodzą z kontraktu i rozjazd między nimi
   jest wykrywany maszynowo.
5. Okno Studio w produkcie odpowiada prototypowi z etapu 1 co do kompozycji,
   gęstości i zachowania.
6. Każde okno w stanie „obecne" nazywa swoją niegotowość i nie przedstawia
   żadnej wartości bez pokrycia w danych — sprawdzone przeglądem, nie deklaracją.
7. Produkt da się przejść od instalacji przez wejście i wybór środowiska do
   pracy w Studio bez ślepego zaułka.

## Etap 3 — Kolejne moduły

Moduły wchodzą po jednym, na przekroju, który działa. Każdy przechodzi tę samą
drogę: okno, potem podłączenie do istniejącego przekroju. Podział na tereny
i ewentualna orkiestracja wieloagentowa mają sens dopiero tutaj — bo dopiero
tutaj istnieje wzorzec, względem którego można sprawdzić wynik.

Harmonogram i dokumentacja zleceń wykonawczych dla tego etapu powstają po
zamknięciu etapu 2, na podstawie tego, co przekrój faktycznie pokazał — nie na
podstawie zamierzeń.

Zakres poczty transakcyjnej jest **zamknięty**: wszystkie siedem listów dostawy
ma przebieg, który je nadaje. Zmiana adresu konta wchodzi parą listów wraz
z drogą wycofania, list o zakończonym przebiegu ma wyzwalacz, adresata
i nastawę częstotliwości — [plan rozbudowy poczty](plan-rozbudowy-poczty.md).
