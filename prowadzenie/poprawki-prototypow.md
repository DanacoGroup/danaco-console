# Poprawki prototypów

Wykaz roboczy zgłoszeń Właściciela do prototypów etapu 1. Pozycja zdjęta z wykazu
po wykonaniu i przyjęciu. Wykaz znika wraz z zamknięciem etapu.

Prototypy: instalacja · uruchomienie wraz z rejestracją i logowaniem · centrum
dowodzenia · TalkIn i moduł Studio. Cztery prototypy pełnoklikalne, około
dziesięciu makiet.

## Reguły obowiązujące wszystkie prototypy

Zgłoszone raz, przy instalatorze; obowiązują wszędzie i nie będą powtarzane przy
kolejnych krokach.

1. **Cała treść do przebudowy** — nagłówki, akapity wprowadzające, opisy pól,
   etykiety, komunikaty. Bez wyjątków i bez pozostawiania zdań „w zasadzie
   dobrych".
2. **Język formalny i zawodowy.** Zwroty potoczne, poufałe i żartobliwe nie
   wchodzą do produktu. Zdania w rodzaju „tutaj wylądują pliki programu", „jeśli
   nie masz powodu tego zmieniać, zostaw jak jest", „nic tu nie jest decyzją na
   zawsze" są nieakceptowalne.
3. **Okno zachowuje stały rozmiar przez cały przepływ.** Przejście między
   krokami nie zmienia szerokości ani wysokości okna. Dziś okno kurczy się
   między krokiem trzecim a czwartym.
4. **Struktura do naprawy** — wykaz usterek strukturalnych niżej, przy każdym
   prototypie.
5. **Animacja przestrzenna** — do wprowadzenia; zakres wskazany przy prototypie,
   którego dotyczy.

## Instalacja

### Struktura i zachowanie

| Rzecz | Stan | Zgłoszenie |
|---|---|---|
| przyciski `Dalej`, `Wstecz`, `Instaluj` | atrybuty `data-krok-dalej` i `data-krok-wstecz` są w znaczniku, obsługi nie ma — przełącza tylko wykaz po lewej | ścieżka, którą przechodzi użytkownik, ma działać |
| rozmiar okna | zmienia się między krokiem 3 a 4 | jeden rozmiar przez cały przepływ |
| krok wyboru składników | wykaz stały, niezależny od wariantu instalacji | wykaz wynika z wybranego wariantu — patrz pozycja 8 rejestru decyzji |
| animacja przestrzenna | brak | do wprowadzenia |

### Treść — krok 1, zasady korzystania

Pozostaje na pierwszym miejscu. Cała treść do napisania od nowa; poniżej to, co
jest wprost błędne merytorycznie, nie tylko językowo.

- „Program jest przeznaczony do pracy na Twoim komputerze" — **nieprawda**.
  Produkt jest hybrydowy: wariant cienki stawia u Operatora samo okno, a praca
  toczy się na serwerze; wariant pełny natywny pracuje lokalnie. Zasady muszą
  obejmować oba, bez rozdzielania na dwie wersje.
- „Treści, nad którymi pracujesz, pozostają u Ciebie" — zakres zawężony do
  treści. Praca nie sprowadza się do treści.
- „Program dostajesz w takiej postaci, w jakiej jest" — do przebudowy.

### Treść — krok 2, miejsce instalacji

Przedmiot kroku jest właściwy: wskazanie miejsca instalacji. Do przebudowy
wszystko poza tym.

- nagłówek i zdanie wprowadzające
- etykieta i opis pola folderu programu
- pole wyboru dotyczące folderu projektów wraz z opisem

### Treść — krok 3, składniki

- zdanie wprowadzające do przebudowy
- wykaz składników zależny od wariantu instalacji, nie stały

### Treść — krok 6, zakończenie

- „Uruchom Danaco Console teraz" — słowo „teraz" jest obce konwencji
  instalatorów i nie występuje w produktach tej klasy.
