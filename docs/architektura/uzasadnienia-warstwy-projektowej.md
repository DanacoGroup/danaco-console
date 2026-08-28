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

## design/zasoby/przedsionek.css
Animacja wejścia stref korzysta z trybu `backwards`, nie `both`, ponieważ element z czynną animacją krycia albo przesunięcia staje się kontekstem nakładania. Przy trybie `both` panel rozwinięty w danej strefie schodził pod strefę stojącą niżej w dokumencie — usterka zmierzona na wykazie zapisanych automatyk pod listwą ustawień. Przy `backwards` klatka początkowa działa tylko na czas opóźnienia i ruchu, a po jego zakończeniu element wraca do zwykłego stanu, więc kontekst nakładania znika razem z animacją.

## budowa/klient-poprzedni/src/dostepy/dodanie-katalogu.css
Arkusz wciąga plik punkty-dostepu.css, ponieważ dodanie katalogu kończy wykaz
punktów dostępu i stoi pod nim w tej samej sekcji. Barw własnych tu nie ma —
wyłącznie żetony warstwy motyw/ — ani żadnej reguły wygaszającej kontrolkę.

## budowa/klient-poprzedni/src/dostepy/katalog-roboczy.css
Ramę sekcji prowadzi plik dostepy.css; barwy pochodzą wyłącznie z żetonów
warstwy motyw/.

## budowa/klient-poprzedni/src/dostepy/nadania-okna.css
Okno ma zbiór nadań, nie jedno — arkusz służy liście, nie pojedynczej
wartości. Ramę sekcji prowadzi plik dostepy.css; żadna reguła nie wygasza
kontrolki.

## budowa/klient-poprzedni/src/dostepy/punkty-dostepu.css
Wykaz jest rozdzielony rodzajem punktu dostępu. Ramę sekcji prowadzi plik
dostepy.css, który ten arkusz wciąga; barwy pochodzą wyłącznie z żetonów
warstwy motyw/.

## budowa/klient-poprzedni/src/motyw/gradienty.css
Zakres użycia gradientów jest wyłącznie ilustracyjny: awatary bez zdjęcia,
rdzeń oznaczenia AOD, grafiki brandowe.

## budowa/klient-poprzedni/src/motyw/przestrzen.css
Odstępy i promienie zaokrągleń są niezależne od motywu. Cienie przełączają
się wraz z motywem, ponieważ tonuje je barwa podłoża — stąd komplet wartości
dla obu motywów wraz z zapasem na zapytanie prefers-color-scheme.

## budowa/klient-poprzedni/src/motyw/rama.css
Żetony ramy żyją poza blokami motywów i nie przełączają się wraz z nimi:
pasek górny pozostaje atramentowy również w motywie jasnym.

## budowa/klient-poprzedni/src/motyw/typografia.css
Trzy role dostają trzy kroje pisma: krój nagłówkowy obsługuje nagłówki,
tytuły środowisk i logotyp, krój interfejsu obsługuje treść roboczą, a krój
o stałej szerokości znaku obsługuje dane techniczne, identyfikatory
i terminal. Pliki krojów wczytuje osobny arkusz fontów na licencji OFL 1.1,
w podzbiorach znaków łacińskich i rozszerzonych łacińskich. Wzorce tekstowe
etykiety wersalikowej, etykiety w kroju maszynowym i danych liczbowych
należą do arkusza fundamentu — ta warstwa deklaruje wyłącznie wartości.

## budowa/klient-poprzedni/src/motyw/warstwy.css
Komponent nie wymyśla własnego z-index: sięga po żeton z tej skali albo nie
ustawia warstwy wcale. Skala jest pełna — wymienia wszystkie poziomy systemu
wizualnego, także te, po które kod jeszcze nie sięga, żeby poziom bez użycia
trzymał swoje miejsce w porządku nakładania.

## budowa/klient-poprzedni/src/motyw/wymiary.css
Gęstość zwarta jest domyślna: kontrolki mają wysokość 32 pikseli, wiersze
tabel 36 pikseli, pasek górny 48 pikseli. Progi punktów łamania stoją jako
wartości nazwane, ponieważ warunek zapytania medialnego nie przyjmuje
zmiennej — żeton daje liczbie jedno miejsce do sprawdzenia, a warstwa
skryptów czyta ten sam próg co arkusze, nie własną kopię.

## budowa/klient-poprzedni/src/motyw/prymitywy.css
Kierunek skali jest monochromatyczną precyzją: neutralne są czysto neutralne,
o równych składowych RGB, a jedyną barwą akcentu systemu jest błękit
sygnałowy. Po stopniu #7C7C7C nie sięga już żadna rola tekstu: na papierze
#F4F4F4 stopień ten daje kontrast 3,80:1, za jasny na tekst i na obrys,
a na powierzchni #181818 — 4,25:1, za ciemny na tekst drugoplanowy.

## budowa/klient-poprzedni/src/komponenty/ikona.css
Rozmiar ikony niosą atrybuty width i height nadawane przez ikony/ikony.ts
z parametru rozmiar. Klasa niesie znaczenie, nie wygląd: po niej selekcjonują
sprawdziany i zaczepiają się motywy. Nośnik jest elementem wierszowym i nie
kurczy się w kontenerze elastycznym.

## budowa/klient-poprzedni/src/komponenty/indeks.css
Warunkiem użycia jest wcześniejsze wczytanie warstwy żetonów z katalogu
motyw/: kolejność importów jest kolejnością kaskady. Preferencja ograniczonego
ruchu jest respektowana globalnie w warstwie żetonów, więc biblioteka nie
powtarza tej reguły osobno — wyjątkiem jest zamiana tętna kropki na pierścień
statyczny w pliku plakietka.css.

## budowa/klient-poprzedni/src/komponenty/postep.css
Weryfikacja kroku ma trzy wyjścia: poprawny w zieleni (przeszedł), błędy
w bursztynie (błędy i licznik obiegów) oraz wstrzymany w czerwieni (zdarzenie
nieoczekiwane, kolejka wstrzymana). Stan nigdy nie jest niesiony samym
kolorem: znak kroku niesie zawsze ikonę albo etykietę.

## budowa/klient-poprzedni/src/komponenty/powiadomienie.css
Stos powiadomień stoi na osobnej warstwie z-index. Stan nigdy nie jest
niesiony samym kolorem: obok krawędzi stanu każdy wariant niesie ikonę
w tej samej barwie oraz tytuł.

## budowa/klient-poprzedni/src/komponenty/suwak.css
Zmienna pozycji suwaka nie jest żetonem warstwy motyw/ — to zmienna ustawiana
w locie na elemencie jako procent wypełnienia toru. Wartość zapasowa wynosi
50%, więc bez ustawienia tor pokazuje połowę.

## budowa/klient-poprzedni/src/konfiguracja/konfiguracja.css
Barw własnych nie ma — wszystko pochodzi z żetonów warstwy motyw/:
monochromatyczna precyzja, jeden błękit sygnałowy, gęstość zwarta. Sygnał
jest rzadki: błękit pojawia się tylko w trzech miejscach — kategorii czynnej
w kolumnie, plakietce pochodzenia wartości nadpisanej i wierszu
dziedziczenia, który obowiązuje. Pole bez zastosowania znika w całości
zamiast leżeć szare i nieczynne.

## budowa/klient-poprzedni/src/konfiguracja/pola-konfiguracji.css
Sygnał w barwie błękitnej występuje tu dokładnie dwa razy: plakietka
pochodzenia wartości nadpisanej oraz wiersz dziedziczenia, który obowiązuje.

## budowa/klient-poprzedni/src/modele/modele.css
Błękit oznacza pozycję czynną wykazu oraz plakietkę konta domyślnego;
bursztyn jest zarezerwowany dla ostrzeżenia, że tryb ZASTĄP zdejmuje prompt
fabryczny. Pole bez zastosowania — katalog konfiguracji przy koncie API, byt
osi przy osi platformy — znika w całości, zamiast stać wygaszone.

## budowa/klient-poprzedni/src/modele/tozsamosc.css
Tryb ZASTĄP zdejmuje prompt fabryczny w całości, więc jego ostrzeżenie stoi
w bursztynie na pełnej szerokości i na stałe. Tryb DOŁĄCZ dostaje barwę
informacyjną.

## budowa/klient-poprzedni/src/modele/wykazy.css
Pozycja wykazu opisująca byt nieczynny nie jest wygaszona ani nieklikalna:
nieczynność jest stanem danych, nie blokadą interfejsu.

## budowa/klient-poprzedni/src/moduly/design/design.css
Kontrolki formularza modułu pochodzą z modele/kontrolki-formularza, więc
moduł wciąga też arkusz modele.css, inaczej pola dm-* byłyby bez oprawy.
Arkusz jest podzielony na cztery pliki wzdłuż odpowiedzialności: rama modułu
i okien wraz ze stanami stoi tutaj; kompozycja.css niesie kanwę, warstwy
i panele boczne; zasoby.css niesie wykaz zasobów i kreator promptu;
podglad.css niesie płytę Preview Window i porównanie wariantów.

## budowa/klient-poprzedni/src/moduly/design/podglad.css
Płyta sygnalizuje kreską brak treści obrazu: obrys ciągły należy się treści,
którą rdzeń oddał, obrys kreskowany — miejscu na treść, której nie ma. Ta
sama zasada obowiązuje w polu podglądu karty zasobu i w miniaturze paska
postępu, żeby brak wyglądał wszędzie w module tak samo.

## budowa/klient-poprzedni/src/moduly/design/tokeny.css
Okno pokazuje barwy produktu, więc próbka musi nieść barwę wpisaną w locie —
to jedyne miejsce w module, w którym barwa trafia do stylu z kodu jako dana
odczytana z motywu, nie jako decyzja projektowa. Sam arkusz nie zna ani
jednej wartości szesnastkowej poza próbkami.

## budowa/klient-poprzedni/src/moduly/developer/okno-project-tree.css
Reszta wyglądu okna, poza menu kontekstowym, idzie z arkusza modułu
developer.css. Plik nie zna barwy dosłownej ani rozmiaru dosłownego —
wyłącznie żetony motywu.

## budowa/klient-poprzedni/src/okna-pomocnicze/pomocnicze.css
Rodzina klas pasa jest wspólna oknu Developer i oknu Diagnostics: ten sam pas
stoi w obu, więc jego wygląd ma jedno miejsce zamiast dwóch kopii do
rozjechania się. Klasy modułowe zostają przy stanach treści okna, bo tam
pokrycie w arkuszach obu okien już jest. Pozycja, której jeszcze nie ma, ma
być czytelna, bo po niej poznaje się, czego platforma nie potrafi — wiersz
braku nie dostaje ani przezroczystości, ani przekreślenia.

## budowa/klient-poprzedni/src/okna-rownolegle/lacznosc.css
Tutaj stoi wyłącznie to, czego biblioteka komponentów nie niesie: miejsce
plakietki w nagłówku gniazda i rozmiar przycisku ponowienia próby.

## budowa/klient-poprzedni/src/okna-rownolegle/panele.css
Rozkład sceny opisuje osobny arkusz układu, wygląd gniazda osobny arkusz
gniazda; ten arkusz nie sięga do żadnego z nich. Gęstość jest wymaganiem,
nie skutkiem ciasnoty: nagłówek panelu ma wysokość jednego wiersza i mieści
tytuł oraz rząd ikon. Nie ma tu reguły dla drugiego rzędu ani dla paska
narzędzi, ponieważ panel z drugim rzędem jest zbudowany źle, nawet jeśli
działa.

## budowa/klient-poprzedni/src/okna-rownolegle/przeciaganie.css
Przedrostek klas jest przedrostkiem obszaru, nie modułu — gest przenoszenia
nie należy do żadnego pojedynczego okna.

## budowa/klient-poprzedni/src/okna-rownolegle/przepelnienie.css
Przedrostek klas jest przedrostkiem obszaru okien równoległych i pomocniczych,
nie modułu. Docelowym domem tego wzorca jest biblioteka komponentów — biblioteka
jest cudza i nietykalna, więc wniosek o przeniesienie idzie wpięciem, a nie
trzecim miejscem na tę samą rzecz.

## budowa/klient-poprzedni/src/okna-rownolegle/uchwyt.css
Widać krechę, chwyta się pas: gdyby pole chwytu było równie wąskie jak
krecha, trafienie w nie wymagałoby celowania, więc pas chwytu wychodzi poza
obrys elementu w obie strony zamiast rozszerzać sam element i zabierać
piksele kolumnom.

## budowa/klient-poprzedni/src/okno-komunikacji/powierzchnia-modulu.css
Historia wątku i pole wypowiedzi mieszkają w arkuszach warstwy rozmowy, nie
w tym arkuszu. Źródłem wartości są wyłącznie żetony motywu i klasy biblioteki
komponentów — ani jednej wartości szesnastkowej, ani jednego odstępu spoza
skali czterech pikseli.

## budowa/klient-poprzedni/src/okno-komunikacji/przelacznik-srodowiska.css
Wygląd uchwytu i wykazu niesie arkusz menu drzewa, wygląd zdania odmowy —
plakietka błędu, a kreskę, tło i odstępy rzędu — arkusz paska zlecenia. Tutaj
zostaje wyłącznie ustawienie uchwytu i odmowy w kolumnie; wyściółka i tło
panelu dałyby w tym miejscu pasek w pasku.

## budowa/klient-poprzedni/src/powloka/wykaz-wynikow.css
Wygląd wykazu niesie w całości arkusz menu drzewa, a wygląd pustki arkusz
drobnych elementów. Nośnik ma wysokość zerową, bo siedzi wewnątrz pola
poleceń paska górnego — wysokość własna rozepchnęłaby pasek w chwili
pojawienia się wyników.

## budowa/klient-poprzedni/src/punkty-izolacji/punkty-izolacji.css
Obudowa okna — nagłówek, opis, akcje, ciało — niesie arkusz ramy okna. Stan
treści każdego obszaru (ładowanie, pustka, błąd, gotowe) niesie arkusz stanu
treści przez klasy tego okna, nie duplikat.

## budowa/klient-poprzedni/src/komponenty/boczna.css
Znacznik aktywności nie jest pełną wstęgą ani tłem sekcji, zgodnie z zasadą
jednego akcentu: sygnał zajmuje nie więcej niż pięć procent ekranu.

## budowa/klient-poprzedni/src/motyw/fonty.css
Pliki licencji krojów stoją w katalogu fonty. Podzbiory latin i latin-ext
niosą pełne polskie znaki diakrytyczne.

## budowa/klient-poprzedni/src/komponenty/karta.css
Sygnał karty środowiska to wstęga górna o wysokości 2 pikseli, sterowana
żetonem wymiaru wstęgi. Kartę środowiska buduje strona-glowna/karta-srodowiska.ts
na rodzinie dn-karta--akcent i dn-strona__karta, nie na klasie karty zwykłej;
kafel komponentu buduje strona-glowna/kafel-komponentu.ts na rodzinie
dn-karta, nie na rodzinie kafla.

## budowa/klient-poprzedni/src/motyw/motyw.css
Mechanizm motywu: atrybut data-theme ustawiony na wartość light albo dark
jest wyborem jawnym; brak atrybutu rozstrzyga zapytanie
prefers-color-scheme. Oba motywy są równoprawne — przełączenie motywu
zmienia wyłącznie wartości żetonów, nie reguły komponentów. Właściwość
color-scheme ustawiają arkusze semantyczne, osobno dla każdego motywu.
Gęstość zwarta jest domyślna; wariant przestronny czeka pod atrybutem
data-gestosc="przestronna".

## budowa/klient-poprzedni/src/motyw/semantyczne-ciemny.css
Skala nie sięga czystej czerni: kończy się na odcieniu #0A0A0A, a tło motywu
to #0F0F0F. Budowa pliku jest identyczna jak w arkuszu motywu jasnego: wybór
jawny przez atrybut, potem zapas rozstrzygany preferencją systemu. Drabina
wag tekstu jest przeliczona dla tego motywu osobno, a nie odbita z jasnego:
stopień trzeci stoi na #9E9E9E, ponieważ #7C7C7C na powierzchni #181818 daje
tylko kontrast 4,25:1. Tekst drugoplanowy ustępuje mu miejsca na #C0C0C0,
inaczej obie wagi zlałyby się w jedną.

## budowa/klient-poprzedni/src/komponenty/pole.css
Stan tylko do odczytu realizuje atrybut readonly. Błąd niesie atrybut
aria-invalid oraz opis w klasie dn-pole-blad.

## budowa/klient-poprzedni/src/komponenty/wybor.css
Suwak wartości ciągłej prowadzi osobny arkusz suwak.css.

## budowa/klient-poprzedni/src/komponenty/zakladki.css
Pas kart sesji używa tej samej mechaniki co zakładki, dokładając wstęgę
aktywności.

## budowa/klient-poprzedni/src/komponenty/okno.css
Wygląd ramy należy do biblioteki, nie do modułu — moduł dokłada swoją klasę
tylko wtedy, gdy naprawdę ma czym, przez parametr przedrostek komponentu
rama-okna. Barwa, obrys i cień pochodzą z klasy dn-karta — tutaj wyłącznie
układ i odstępy, wszystkie z żetonów warstwy motyw/.

## budowa/klient-poprzedni/src/komponenty/pasek.css
Pasek sięga wyłącznie po żetony rama z pliku motyw/rama.css, które stoją poza
blokami motywów. Tło pola wyszukiwania w pasku jest półprzezroczystą bielą
sześcioprocentową, a nie żetonem: to warstwa na atramencie ramy, niezależna
od motywu.

## design/zasoby/zetony/zetony.css
Powielenie wartości cienia sygnału i pozostałych żetonów między blokiem `prefers-color-scheme` i blokami motywów jest świadome: stanowi mechanizm kaskady na wypadek braku jawnego wyboru motywu, a nie drugie źródło prawdy dla tych wartości. Sekcja czternasta, dotycząca ograniczonego ruchu bez konfiguracji per komponent, nie ma jeszcze treści — zachowanie `prefers-reduced-motion` obsługuje w całości plik `zetony/ruch.css`.

## budowa/klient-poprzedni/src/komponenty/drobne.css
Klatki nazwane dn-obrot definiuje ten arkusz i wykorzystuje je również
przycisk w stanie ładowania; klatki dn-tetno definiuje plakietka.css.

## budowa/klient-poprzedni/src/motyw/stany.css
Rola -tekst niesie tekst i ikonę stanu, rola -tlo tło plakietki lub alertu,
rola -obrys obrys plakietki lub alertu. Stan informacyjny jest rodziną
sygnału: system ma jedną barwę akcentu, więc stan informacyjny nie
wprowadza kolejnej.

## budowa/klient-poprzedni/src/komponenty/plakietka.css
Klatki nazwane dn-tetno definiuje ten arkusz; korzystają z nich także wpis
pracujący w pliku wpis.css oraz plik drobne.css. Plakietkę roli koordynatora
niesie dziś klasa dn-plakietka--sygnal, budowana przez
okna-rownolegle/plakietka-roli.ts.

## budowa/klient-poprzedni/src/komponenty/wpis.css
Mapowanie ról kontraktu na klasy wpisu: rola user na klasę dn-wpis--czlowiek,
rola assistant na klasę dn-wpis--inteligencja, rola system na klasę
dn-wpis--system, rola tool na klasę dn-wpis--system wraz z plakietką roli
niosącą nazwę narzędzia. Wpis w trakcie pracy dokłada klasę
dn-wpis--pracuje; klatki nazwane dn-tetno definiuje plakietka.css.

## budowa/klient-poprzedni/src/motyw/podglad-zetonow.html
Strona działa bez rdzenia, kanału kontraktu i sesji, więc nie mówi nic
o działaniu produktu ani o tym, czy dany żeton jest przez aplikację
rzeczywiście używany.

## budowa/klient-poprzedni/src/motyw/fundament.css
Fundament wymaga wcześniejszego wczytania krojów i żetonów motywu; kolejność
wczytania ustala arkusz motywu. Domyślny styl przeglądarki dla atrybutu
hidden ma tę samą wagę co reguła klasy, więc komponent z własnym display
(na przykład pole w układzie flex) przesłania go w kaskadzie i element
z atrybutem hidden zostaje widoczny — ważność przywraca ukrycie.

## budowa/klient-poprzedni/src/komponenty/przycisk.css
Wariant sygnałowy jest zarezerwowany dla jednego działania systemowego na
widok i nie zastępuje atramentu jako działania domyślnego. Żaden wariant nie
odbiera klikalności — ładowanie i błąd komunikuje atrybut ARIA oraz
wskaźnik, nigdy blokada. Klatki nazwane dn-obrot dla wskaźnika ładowania
definiuje plik drobne.css przy spinnerze; reguła @keyframes działa globalnie
niezależnie od kolejności arkuszy.

## budowa/klient-poprzedni/src/powiadomienia/powiadomienia.css
Kolumna centrum powiadomień jest regulowana wyłącznie na szerokość. Plik nie
zna ani jednej barwy zapisanej wprost: wszystko pochodzi z żetonów motywu,
a każdy selektor zaczyna się od przedrostka po-.

## budowa/klient-poprzedni/src/powloka/powloka.css
Wygląd każdego z czterech pasów powłoki niesie własny arkusz obok modułu,
który go buduje: pasek górny, karty sesji, nawigacja modułów i obszar
roboczy mają osobne arkusze. Ten arkusz wciąga sam moduł powłoki, więc
trafia do pakietu razem z nim; wszystkie barwy, stopnie pisma, odstępy,
promienie i czasy pochodzą z żetonów motywu, a wymiary układu są
wielokrotnością czterech pikseli. Szerokość pełnej kolumny nawigacji niesie
system wizualny; zwężony wariant jest właściwością samej powłoki — system
nazywa zachowanie (nawigacja schodzi do ikon), ale szerokości szyny ikon nie
ustala. Próg zawężenia jest punktem łamania w2 systemu wizualnego, tym
samym, przy którym system nazywa zwinięcie bocznej nawigacji do ikon.
Wartość stoi tu liczbą, ponieważ zapytanie medialne nie czyta zmiennych
własnych; żeton progu w arkuszu wymiarów pozostaje jej jedynym źródłem
i musi się z nią zgadzać.

## budowa/klient-poprzedni/src/komponenty/rozmowa.css
Klasa dn-rozmowa nie jest klasą dn-okno z pliku okno.css: dn-okno to rama
okna operacyjnego modułu z pionem i gniazdami akcji, a dn-rozmowa to
pięciowierszowa siatka rozmowy. Obudowa nie nosi cienia, inaczej niż karta —
okno rozmowy stoi wewnątrz środowiska i uniesienia nie deklaruje; kto go
potrzebuje, dokłada kartę na rodzicu. Wnętrze wierszy niosą osobne arkusze:
wpis.css dla historii, postep.css dla monitora, drobne.css dla pola
wpisywania i przybornika.

Pasek nagłówka tożsamości przewija się bez widocznego suwaka, ponieważ pola
tożsamości nigdy się nie łamią. Panel kontekstu nie ma własnego tła, bo
kontekst należy do powierzchni okna, nie do panelu. Tło historii jest
zadeklarowane wprost jako żeton tła, ponieważ obudowa wyżej ustawia żeton
powierzchni — bez tej deklaracji historia dostałaby barwę karty zamiast
papieru roboczego. Monitor wykonania niesie pasek postępu w wariancie
rozciągniętym oraz przycisk zatrzymania. W wariancie zwartym monitora nie
ma, więc dół okna niesie własną kreskę górną zamiast kreski monitora.

## budowa/klient-poprzedni/src/rozmowa/blok.css
Blok prowenancji jest jedynym miejscem, w którym warstwa rozmowy sięga po
tło akcentu — powierzchnia drobna i zwinięta domyślnie.

## budowa/klient-poprzedni/src/ladowanie/ladowanie.css
Bryła jest prawdziwą bryłą przestrzeni CSS, nie obrazkiem udającym
trójwymiar ani płótnem WebGL, ponieważ sześć ścian i jeden obrót to
dwadzieścia wierszy arkusza, a WebGL wnosi zależność, warstwę sterowników
i ryzyko ekranu, który u odbiorcy nie wstanie. Scena stoi pod żetonem
warstwy bramki, a nad wszystkim innym: przesłona logowania schodzi dopiero
wtedy, gdy scena już stoi, więc żaden ekran nie mignie pomiędzy nimi.

## budowa/klient-poprzedni/src/powloka/usuniecie-sesji.css
Ramę modalu, nagłówek, ciało i stopkę niesie arkusz nakładki, a pas stanu
arkusz drobnych elementów; ten arkusz dokłada wyłącznie to, czego biblioteka
nie ma. Barwy pochodzą z żetonów motywu, więc oba motywy są obsłużone bez
osobnych reguł. Szerokość modalu usunięcia jest liczona tym samym wzorem co
pozostałe modale, więc na wąskim oknie zachowuje się jak każdy modal
produktu.

## budowa/klient-poprzedni/src/mission-control/mission-control.css
Kolejność importów odpowiada kolejności sekcji na ekranie: pulpit.css niesie
ramę, czoło, kafle liczb, pas decyzji i pas działań; matryca.css niesie
kanały modelu, matrycę sesji, pas Relacje i rząd Utwórz; kolumny.css niesie
procesy w tle, kolejki i zespół agentów; stany-pulpitu.css niesie stany
puste i braki źródła w kontrakcie. Warunkiem użycia jest wcześniejsze
wczytanie warstwy żetonów oraz biblioteki komponentów, co zapewnia plik
aplikacja/arkusze-stylow.ts. Arkusz wciąga moduł mission-control.ts, więc
trafia do pakietu razem z pulpitem i nie wymaga wpisu w cudzym wykazie.

## budowa/klient-poprzedni/src/powloka/karty-sesji.css
Wygląd zakładki wraz z podkreśleniem karty czynnej pochodzi z biblioteki
komponentów, tak jak karty rozciągnięte na całą wysokość pasa, żeby
podkreślenie legło dokładnie na jego kresce. Biblioteka niesie pasek ramy
o wysokości 48 pikseli na atramencie marki; pas kart stoi na powierzchni
pracy i jest od niego niższy.

## budowa/klient-poprzedni/src/mission-control/kolumny.css
Arkusz nie ma reguły disabled ani wygaszającej: przyciski transportu
wyglądają tak samo w każdym stanie kolejki, a stan niesie osobna plakietka
ze słowem opisującym stan.

## budowa/klient-poprzedni/src/powloka/pasek-gorny.css
Wysokość 48 pikseli, tło ramy i barwa treści przychodzą z klasy biblioteki
komponentów. Arkusz sięga wyłącznie po żetony motywu i klasy biblioteki;
klasy, których biblioteka nie niesie — kontekst pracy, rozdzielacz i uchwyt
profilu — stoją tutaj pod nazwą własną widoku. Pasek jest atramentowy w obu
motywach: każda kontrolka, która na nim staje, ubiera się żetonami ramy, nie
żetonami treści; wyjątkiem są nakładki opadające pod pasek (wykaz wyników,
panel menu profilu), bo te leżą już na powierzchni pracy i przełączają się
z motywem razem z nią. Uchwyt menu profilu buduje skrypt menu profilu
mechanizmem biblioteki menu drzewa, który normalnie ubiera go żetonami
treści; na powierzchni pracy to jest właściwe, ale na atramencie ramy uchwyt
był jedyną rzeczą w pasku przełączającą się z motywem, więc w motywie jasnym
stawał się jasnym kafelkiem z ciemnym napisem — stąd żetony ramy tutaj, bez
zmiany mechanizmu i bez ruszania biblioteki.

## budowa/klient-poprzedni/src/mission-control/matryca.css
Wartości pochodzą wyłącznie z żetonów warstwy motyw/: bez zapisów
szesnastkowych i bez odstępów spoza skali czterech pikseli.

## budowa/klient-poprzedni/src/mission-control/pulpit.css
Wartości pochodzą wyłącznie z żetonów warstwy motyw/ i klas biblioteki
komponenty/: arkusz nie zawiera barwy szesnastkowej ani odstępu spoza skali
czterech pikseli. Akcent złoty ogranicza się do wstęgi pasa decyzji i cienia
przycisku wezwania; nigdzie nie wypełnia powierzchni. Arkusz nie zawiera
reguły disabled ani żadnej innej reguły wygaszającej.

## design/zasoby/stanowisko.css
Reakcja na wskazanie i naciśnięcie jest widoczna w każdym oknie tego pakietu, nie tylko w centrum dowodzenia: belka narzędzi, żetony polecenia i pulpit ramy podnoszą się o jeden piksel i ustępują pod naciśnięciem, nośnikiem reakcji jest tło, barwa i przekształcenie, nigdy sama zmiana rozmiaru ikony.

## budowa/klient-poprzedni/src/mission-control/stany-pulpitu.css
Treść stanów pustych buduje plik stan-pusty.ts wraz z widokami poszczególnych
sekcji; ten arkusz niesie wyłącznie ich wygląd.

## budowa/klient-poprzedni/src/mobile/mobile.css
Obudowę okna niesie komponenty/rama-okna.ts, a przyciski komponenty/przycisk.css.
Jeden wymiar jest tu własny — cel dotyku: motyw daje na wskaźnik zgrubny
40 pikseli z pliku motyw/wymiary.css, a ten ekran potrzebuje 44 pikseli, więc
próg podnosi żeton lokalny o przedrostku modułu, tak jak robi to Terminal,
zamiast wyjątku rozsypanego po regułach — podniesienie progu w całej powłoce
należy do motywu. Układ jest jednoręczny: karty pozycji czyta się oczami
u góry, a arkusza dróg i paska kwitu dotyka się kciukiem u dołu. Żadne pismo
nie schodzi poniżej najmniejszego stopnia żetonu, żeby treść pozostała
czytelna bez powiększania.

## budowa/klient-poprzedni/src/moduly/developer/developer.css
Plik nie zna ani jednej barwy dosłownej, ani jednego rozmiaru czcionki
dosłownego: wszystko pochodzi z żetonów motywu i z biblioteki komponenty/.
Arkusz niesie wyłącznie rozkład pięciu okien modułu — edytora kodu, drzewa
projektu, panelu git, logu budowania i zakładek narzędzi — wraz
z odróżnieniem trzech stanów obowiązkowych oraz oznaczeniem wagi zgłoszeń
budowania i konfliktów repozytorium.

## design/zasoby/rama.css
Stopka szyny nawigacji niesie cztery stopnie jasności — menu aplikacji pełną bielą, środowiska i pozycje pracy lekko przydymione, moduły wysunięte ze środowiska mocniej, ustawienia w stopce najmocniej — tak aby różnica jasności mówiła, na którym poziomie wyboru stoi operator. Ikony pracy zostają w pełnej jasności ramy, a ikony stopki schodzą do żetonu --dn-rama-tekst-4, przy kontraście 3,10 do 1 na ramie, czyli progu ustalonego dla grafiki nietekstowej; najechanie i fokus przywracają pełną jasność, żeby cel kliknięcia nie był bledszy od pozostałych pozycji.
