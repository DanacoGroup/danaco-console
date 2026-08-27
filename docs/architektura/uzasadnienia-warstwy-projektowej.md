# Uzasadnienia warstwy projektowej

Uzasadnienia i zastrzeżenia projektowe przeniesione z komentarzy plików warstwy
projektowej, których treść przekracza dopuszczalną długość nagłówka. Rozdziały
noszą nazwę pliku, którego dotyczą.

### design/zasoby/powloka.js

Plik jest jedynym źródłem prawdy dla powłoki prototypów i niesie markup stałych
elementów okna: szynę nawigacji, belkę tytułową, pasek edycji (wstążkę), pasek
stanu oraz okno „Dostosuj paski”. Każdy prototyp podaje wyłącznie treść okna
centralnego w `<template id="dn-tresc-okna">` oraz punkt montażu
`<div id="dn-powloka-montaz">` — skrypt wstrzykuje powłokę wokół tej treści,
odtwarzając strukturę okna jeden do jednego. Działa wprost z dysku przez
`file://`, bez zapytań sieciowych.

### design/zasoby/tresci/licencja.js

Treść jest wynikiem działania generatora `zbuduj-licencje.py`, uruchamianego
na źródłowym tekście umowy licencyjnej. Plik nie jest przeznaczony do
bezpośredniej edycji — zmiana wchodzi w źródle i przechodzi przez generator
ponownie. Treść stoi w skrypcie, a nie w osobnym pliku do pobrania, ponieważ
okno bywa otwierane wprost z dysku, a przeglądarka blokuje wówczas pobieranie
plików towarzyszących.

## budowa/klient-poprzedni/src/moduly/research/research.css

Plik nie zna ani jednej barwy zapisanej wprost — wszystkie wartości pochodzą
z żetonów motywu. Klasy `dn-*` (pole, przycisk, plakietka, modal, pusty stan,
spinner, dymek) wnosi biblioteka komponentów; tutaj leży wyłącznie rozkład
okien modułu i to, czego biblioteka nie ma: pasy układu, wykazy źródeł
i ustaleń oraz znakowanie faz stanu.

Kontrolki formularza pochodzą z `modele/kontrolki-formularza`, więc moduł
wciąga arkusz `modele.css` — inaczej pola `dm-*` byłyby bez oprawy.

## budowa/klient-poprzedni/src/moduly/multitasking/formularz-powolania.css
Barwy i odstępy formularza pochodzą wyłącznie z żetonów motywu, więc plik działa poprawnie w obu motywach bez odrębnego kodu. Wygląd pól i przycisku pochodzi z biblioteki komponentów wspólnych; ten plik niesie wyłącznie rozkład, którego biblioteka nie dostarcza.

## budowa/klient-poprzedni/src/aktualizacja/aktualizacja.css
Barwy pasa pochodzą wyłącznie z żetonów motywu, bez żadnej wartości barwnej wpisanej wprost, dzięki czemu podmiana żetonu wystarcza i pas nie wypada z pomiaru kontrastu. Arkusz nie powtarza biblioteki: tor postępu korzysta z klasy `.dn-postep` z pliku `komponenty/postep.css`, chronologia z klasy `.dn-tabela` z pliku `komponenty/tabela.css`, a przycisk przywołania z klasy `.dn-btn` z pliku `komponenty/przycisk.css`. W tym pliku pozostaje wyłącznie układ pasa i miejsca styku.

## budowa/klient-poprzedni/src/moduly/multitasking/panel-orkiestracji.css
Barwy i odstępy pochodzą wyłącznie z żetonów warstwy motywu, arkusz nie wprowadza własnych decyzji wyglądu ponad układ. Arkusz nie wygasza kontrolek: sekcja bez powiązania z modułem Automations sygnalizuje to zdaniem, a nie odcieniem szarości.

## budowa/klient-poprzedni/src/moduly/roundtable/analiza-debaty.css
Oba okna dzielą jeden nośnik wykazu, bo pokazują to samo w kształcie: pozycję ze znacznikiem mówcy, nazwą i wierszami metryki. Dwa osobne nośniki o tym samym układzie rozjechałyby się przy pierwszej poprawce. Barwy i rozmiary czcionki pochodzą wyłącznie z żetonów motywu i z biblioteki komponentów.

## budowa/klient-poprzedni/src/moduly/roundtable/debata.css
Rozkład modułu, stany obowiązkowe, wykaz pozycji i skład debaty stoją w pliku roundtable.css. Barwy i rozmiary czcionki pochodzą wyłącznie z żetonów motywu i z biblioteki komponentów.

## budowa/klient-poprzedni/src/moduly/roundtable/panel-debaty.css
Panel wciąga swoje arkusze sam, bo gospodarzem bywa moduł, którego arkusz o Roundtable nic nie wie; wydzielenie sprawia, że gospodarz dostaje wyłącznie reguły panelu. Obok tego pliku panel wciąga roundtable.css i wyłącznie z niego bierze cztery reguły stanu treści, bo nośnik stanu treści modułu ma jeden zapis. Wypowiedzi kilku modeli idą kolumną jedna pod drugą, a odróżnia je wstęga w wariancie motywu i tożsamość w nagłówku głosu. Plik nie zapisuje żadnej barwy ani odstępu wprost — wszystko pochodzi z żetonów motywu.

## budowa/klient-poprzedni/src/moduly/studio/dziennik-kontrola.css
Arkusz wchodzi z plików obsługi dziennika, kopii zapasowych i blokad, a nie z pliku studio.css, ponieważ ten należy do innego odcinka prac. Wciągnięcie arkusza przez moduł kliencki daje potrzebny skutek bez zmiany w cudzym pliku. Żadna barwa nie jest zapisana wprost — wszystkie wartości pochodzą z żetonów motywu, więc oba motywy obsługują się same.
