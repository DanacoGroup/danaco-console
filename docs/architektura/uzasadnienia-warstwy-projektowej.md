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
