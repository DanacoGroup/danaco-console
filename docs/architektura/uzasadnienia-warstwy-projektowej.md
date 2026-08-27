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

## design/zasoby/okna/instalator/stany.js
Mapa przechowuje wyłącznie klucze katalogu treści, nigdy gotowy tekst, ponieważ zmiana słowa wchodzi w katalog treści, nie w mapę stanów. Ekran wykrycia procesora przyjmuje trzy stany: architekturę x64 z podtytułem i zaznaczeniem x64, architekturę arm z podtytułem i zaznaczeniem arm oraz stan nierozpoznany, w którym nic nie jest zaznaczone, wskazówka kieruje do sprawdzenia procesora ręcznie, a czynność główna jest wyłączona. Ekran przebiegu zapisu przyjmuje stan przebiegu z paskiem miary, wykazem etapów i paskiem szczegółów, stan wycofywania z paskiem bez miary określonej i bez wykazu oraz z wyłączonymi czynnościami, a także stan błędu bez paska i wykazu, z blokiem błędu i trzema czynnościami. Ekran wyniku instalacji przyjmuje stan ukończenia ze znakiem sukcesu przy tytule oraz stan ukończenia z uwagami ze znakiem ostrzeżenia i banerem listy uwag. Stany pozaekranowe obejmują krok drugi bez zgody, w którym czynność główna jest wygaszona atrybutem dostępności, a fraza nawigacyjna zmienia treść i barwę po próbie przejścia, wybór niezgodny z ostrzeżeniem pod kartami i pytaniem przy próbie przejścia dalej, oraz anulowanie z pytaniem o skutek, gdzie potwierdzenie uruchamia wycofywanie.

## budowa/klient-poprzedni/src/aplikacja/aplikacja.css
Wygląd samych widoków opisują arkusze ich katalogów, więc ten plik ich reguł nie powtarza. Barwy, stopnie pisma, odstępy, promienie i czasy pochodzą z żetonów katalogu `motyw/`, bez wartości szesnastkowych zapisanych wprost, a wymiary są wielokrotnością 4 pikseli.

## budowa/klient-poprzedni/src/aplikacja/pas-posuniec.css
Pasek nie wprowadza wartości szesnastkowych ani odstępów spoza skali motywu —
każda barwa i każdy odstęp pochodzi z żetonów. Ruch asystenta nie jest ani
sukcesem, ani awarią, dlatego pasek nosi barwę informacyjną: para żetonów
`--dn-informacja-tekst` i `--dn-informacja-tlo` przechodzi tę samą kontrolę
kontrastu co pozostałe pary motywu, a pasek nie dokłada barwy własnej.

## budowa/klient-poprzedni/src/aplikacja/powloka.css
Wskaźnik ładowania łączności korzysta ze spinnera biblioteki zmniejszonego do
rozmiaru kropki, którą zastępuje na czas łączenia — dzięki temu plakietka nie
rośnie i nie przeskakuje przy zmianie stanu. Kropka i spinner mają własne
właściwości wyświetlania, dlatego ukrycie każdego z nich wymaga osobnego
zapisu jawnego.

Punkt łamania w2 (960 pikseli) stoi w zapytaniu medialnym liczbą, ponieważ
zapytanie medialne nie odczytuje zmiennych własnych arkusza. Jedynym źródłem
tej wartości pozostaje żeton `--dn-bp-w2` z pliku motyw/wymiary.css — liczba
w zapytaniu musi się z nim zgadzać.

Arkusz powloka.css idzie w wykazie arkuszy stylów zaraz po pliku okno.css.
Kolejność jest konieczna: widok okna niesie własny komplet zmiennych `--dc-*`
o barwach wpisanych na stałe, a nadpisanie ich tutaj sprawia, że okno i komplet
sterowania zmieniają się razem z motywem, bez zmian w katalogu widoku okna.

## budowa/klient-poprzedni/src/moduly/multitasking/zadania-w-tle.css
Rodzina klas panelu ma przedrostek dm-. Plakietka stanu, tabela agentów
i przyciski pochodzą z biblioteki wspólnej; ten arkusz zawiera wyłącznie to,
czego biblioteka nie pokrywa: rozkład karty przepływu, wiersz etapu ze
wskaźnikiem kropkowym oraz znacznik kolumny bez pokrycia w kontrakcie danych.
Atrybut data-bez-pokrycia wycisza kolumnę do koloru tekstu trzeciorzędnego
i ustawia kursor pomocy: kolumna pozostaje widoczna, a powód nieobecności
danych pokazuje atrybut title elementu.

## budowa/klient-poprzedni/src/moduly/studio/petla-okno.css
Arkusz nie nadaje oknu pętli stałej kolumny w powierzchni dokumentu: postacią
domyślną jest nakładka nad treścią, która wchodzi znacznikiem przebiegu
i schodzi po zwinięciu, nie zabierając dokumentowi ani jednego piksela, gdy
nie jest używana. Stała kolumna istnieje wyłącznie jako tryb wybierany
ręcznie. Odstępy, barwy i promienie pochodzą z warstwy wspólnej żetonów;
arkusz nie wprowadza żadnej wartości surowej poza wymiarami układu, których
warstwa wspólna nie nazywa.

## design/03-marka/zastosowania/html/sygnatura-poczty-jasna.html
Odbiorca kopiuje cały blok tabeli wzorca i wkleja go w ustawieniach klienta poczty elektronicznej. Pola ujęte w nawiasach kwadratowych wymagają uzupełnienia własnymi danymi. Szerokość wzorca wynosi sto procent szerokości dostępnej do pięciuset dwudziestu pikseli, znak zajmuje sto sześćdziesiąt osiem na siedemdziesiąt dwa piksele w podwójnej rozdzielczości pliku, linia reguły ma grubość jednego piksela, a kropka sygnału pozostaje jedynym akcentem barwnym całej sygnatury.

## budowa/klient-poprzedni/src/moduly/studio/przybornik-znakowania.css
Arkusz wchodzi z pliku klienckiego przybornik-znakowania, a nie z arkusza
studia modułu: dopisanie importu do cudzego arkusza byłoby zmianą poza tym
odcinkiem prac, więc wciągnięcie odbywa się przez moduł kliencki. Panel
operacji zwinięty nie zajmuje kolumny pasa wiodącego, ponieważ kolumnę
rezerwuje układ siatki pasa wiodącego z arkusza modułu należącego do innego
odcinka; rozstrzyga to wyłącznie selektor stanu, bez zmiany tamtego arkusza.
Trzy byty marginesu mają różne obramowania, żeby rodzaj wpisu dało się
rozpoznać bez czytania treści dymka. Pływak narzędzi ukrytych mieści szybkie
działania oraz suwaki. Stały panel operacji zajmuje osobną kolumnę pasa
wiodącego modułu studia. Przybornik znakowania mieści również brakujące
pozycje w jednej kolumnie panelu bocznego. Schowek udostępnia zawężenie
wyników i podgląd treści. Osadzenie przeglądarki i biblioteki ma część
nagłówkową, przewijany podgląd treści źródła i listę pochodzeń. Mikrofon
w trakcie nagrywania dodatkowo pulsuje animacją tętna, sygnalizując aktywne
nagrywanie dźwięku.

## budowa/klient-poprzedni/src/moduly/studio/strona-postac.css
Panele nastaw strony są nakładkami otwieranymi z pasa otwarć, domyślnie
schowanymi atrybutem hidden; stałej kolumny nie zajmują i zajmować nie mogą,
więc arkusz nie zawiera ani jednej reguły szerokości kolumny. Żetony barw
i odstępów pochodzą z warstwy wspólnej produktu — arkusz nie zakłada własnej
palety, żeby panel nastaw strony wyglądał jak reszta modułu studia, a nie jak
osobna wyspa wizualna.

## budowa/klient-poprzedni/src/moduly/studio/studio-okna.css
Rama i stany wspólne wszystkim pięciu oknom operacyjnym stoją w arkuszu
wspólnym studia; ten arkusz nie zapisuje żadnej barwy wprost i korzysta
wyłącznie ze zmiennych motywu.

## budowa/klient-poprzedni/src/moduly/studio/studio.css
Biblioteka komponentów wnosi pole, przycisk, plakietkę, tabelę, pusty stan,
spinner i dymek. Okno pracy z dokumentem ma własny arkusz. Kontrolki
formularza modułu pochodzą z arkusza kontrolek formularza modeli, więc moduł
wciąga też arkusz wspólny modeli — inaczej pola formularza byłyby bez oprawy
wizualnej.

## budowa/klient-poprzedni/src/asystent-plywajacy/asystent-plywajacy.css
Komponent nie wymyśla własnego z-index, tylko bierze żeton `--dn-z-modal`
z pliku motyw/warstwy.css, a nie `--dn-z-aod`: kanał AOD jest osobnym bytem
kontraktu (aod.voice.command) i niezbudowaną kontrolką paska górnego, więc
zajęcie jego warstwy zapowiadałoby zderzenie. Warstwa jest niższa od
`--dn-z-powiadomienie`, bo odmowy tego dymka jadą powiadomieniem i mają być
widoczne ponad nim. Położenie jest stałe względem okna dokumentu, w prawym
dolnym rogu, poza obszarem roboczym powłoki, więc zmiana modułu niczego tu
nie przesuwa.

Reguły śladu posunięć asystenta niosą wyłącznie kontrast i znaczenie stanu,
bez nowej ikonografii ani wariantów wizualnych, ponieważ warstwa wizualna
całego komponentu jest tymczasowa.
