# Danaco Console — Uzasadnienia komentarzy klienta poprzedniego

Dokument gromadzi uzasadnienia, które przekraczają dopuszczalną długość
nagłówka komentarza w plikach `budowa/klient-poprzedni/` i `budowa/klient/`.
Każdy rozdział nosi nazwę pliku źródłowego, którego uzasadnienie dotyczy.

## budowa/klient-poprzedni/src/moduly/studio/okno-pracy-z-dokumentem.ts

Okno łączy w jednej powierzchni to, co dawniej działało jako cztery osobne
okna nad tą samą treścią: edycję z formatowaniem, podgląd wydania i różnicę
wersji jako tryby jednego widoku, przy czym różnica jest warstwą nakładaną na
treść, a nie drugą kolumną obok niej. Wynik operacji wchodzi do dokumentu jako
zmiana śledzona autora `model`, a jej poprawianie przed decyzją odbywa się tam,
gdzie zmiana stoi — w treści, nie w osobnym buforze.

Tryb wykazu operacji (narzędzia ukryte albo stały panel) jest jedną nastawą na
cały moduł, więc rozstrzyga ją i zapisuje okno pracy, a nie pływak i stały
panel osobno — inaczej te dwa miejsca rozjechałyby się w tym, co pokazują.

Siedem źródeł postaci dokumentu, kontroli pracy i wstawień (rodziny poleceń
`studio.page.*`, `style.*`, `format.*`, `table.*`, `object.*`, `apparatus.*`,
`field.*`, `journal.*`, `markup.*`, `lock.*`, `backup.*`, `template.*`) nie
należą do portu pracy, dlatego wchodzą osobnym bytem `ZrodlaPostaciStudia`,
a nie dopiskiem do `ZrodloPracyStudio`. Byt jest nieobowiązkowy: gdy któregoś
źródła nie podano, okno działa dalej, tylko odpowiednia rodzina poleceń jest
wyłączona — brak leży wtedy po stronie złożenia modułu, nie rdzenia.

Górna granica wielkości jednej kartki (sześć megabajtów) wynika z tego, że
kartka A4 w 96 punktach na cal jako obraz PNG mieści się w rzędzie setek
kilobajtów, więc sześć megabajtów daje zapas także przy nośniku
wielkoformatowym; rdzeń dostaje tę granicę, bo kontrakt każe mu wtedy odmówić
z podaniem zmierzonej wielkości, zamiast oddać treść uciętą, co wyglądałoby na
kartkę zepsutą przez rdzeń. Górna liczba kartek pobieranych jednym renderem
(dwadzieścia cztery) wynika z tego, że pismo dwustustronicowe dałoby dwieście
żądań i dwieście obrazów w pamięci karty na raz, a podgląd wydania służy
sprawdzeniu składu, nie czytaniu całego pisma obrazkami.

Postać dokumentu musi powstać przed powierzchnią, bo powierzchnia bierze od
niej magazyn nastaw widoku i zgłasza jej chwyty linijki; odwrotna kolejność
zostawiłaby chwyty bez odbiorcy. Kartki wyrysowane przez rdzeń tracą
aktualność w chwili, gdy Operator zmieni treść, dlatego są zdejmowane od razu:
pokazywanie starego wyrysu jako podglądu treści bieżącej byłoby nieprawdą
o dokumencie, nie oszczędnością jednego wywołania.

Mikrofon wiersza polecenia dostaje z powłoki okno komunikacji, nie kartę
sesji, bo `speech.audio.upload` przyjmuje puste `sessionId` jako nagranie bez
przypisania do karty — zgadywanie karty z okna byłoby wskazaniem nieprawdziwym.
Licznik użycia operacji jedzie do rdzenia po każdym uruchomieniu, bo z niego
bierze się kolejność czynności na wierzchu pływaka: z użycia, nie z domysłu.

Zakładka założona przyciskiem przybornika jest elementem aparatu dokumentu
i zapisuje się w rdzeniu, więc przeżywa zapis; wykaz w oknie zostaje obok, bo
powrót ma działać od ręki, bez odczytu, ale prawdą o dokumencie jest wiersz
w rdzeniu, nie ten wykaz. Trzeci argument przybornika znakowania jest rdzeniem
znakowania: bez niego przybornik pokazuje znakowania samej sesji, z nim
znakowanie jest trwałe i przechodzi przez zapis dokumentu.

Dziennik, kopie i blokady stoją nakładkami przy powierzchni, bo dotyczą tej
samej treści: dziennik cofa jej czynności, kopie ją ratują, blokady jej
pilnują — każdy przyjmuje skutek tą samą drogą co reszta okna, treścią roboczą
stanu. Postać po czynności czyta się z rdzenia, a nie przyjmuje nieopisanym
ładunkiem z odpowiedzi: postać jest drzewem o kształcie kontraktu, a wartość
niesprawdzonego typu wstawiona w stan okna byłaby drugą, niepewną prawdą
o dokumencie.

Wskaźnik zakresu operacji rozróżnia trzeci przypadek — cały dokument z
zaznaczeniem pominiętym wyborem ręcznym z Tools Panelu — osobno od pozostałych
dwóch, bo bez niego wybór ręczny wyglądałby jak usterka: zaznaczenie widać,
a operacja idzie na całość.

Panel osadzenia źródeł stoi w powłoce treści, nie na całej powierzchni okna,
żeby wstążka i pasek statusu zostały widoczne i Operator nie stracił drogi
powrotu. Pas otwarć trzech paneli (schowek, źródła, przybornik) stoi osobno od
wstążki, bo wstążka niesie formatowanie i widok, a te trzy są narzędziami
bocznymi treści; siedem powierzchni postaci i kontroli pracy w dalszej części
okna idzie tym samym wzorem — przycisk w pasie, treść nakładką, panel
zamknięty nie zajmuje miejsca — jako reguła stopniowego ujawniania, nie
oszczędność powierzchni.

Przerysowanie marginesu dzieje się jednym wywołaniem na trzy byty naraz
(komentarz, propozycja, zmiana śledzona), bo stoją w jednej kolumnie i muszą
się ustawić wobec siebie — dwa osobne przerysowania kładłyby je jedne na
drugich.

Zasoby stron w odpowiedzi renderu stoją na początku wykazu, a dokument
w formacie docelowym dochodzi na jego końcu, więc kartek jest dokładnie tyle,
ile podaje pole `pages`, i do tej liczby wykaz się przycina; rozpoznawanie
kartki po typie treści wymagałoby pobrania także dokumentu, czyli megabajtów
pobranych po to, żeby je odrzucić.

Wstawianie treści w miejsce kursora jest jedną drogą wspólną dla schowka,
dyktowania i wniesienia ze źródła: zaznaczenie zostaje zastąpione tak jak
w każdym edytorze, a brak zaznaczenia znaczy „na końcu treści", bo kursor bez
zaznaczenia stoi tam, gdzie stan modułu ostatnio go widział.

Malarz formatów kopiuje postać akapitu, nie treść, więc nie jedzie schowkiem
rdzenia: rodzina poleceń `clipboard.*` niesie treść i rodzaj wpisu, a nie
arkusz nastaw akapitu, dlatego postać czyta się z nastaw akapitu okna.

## budowa/klient-poprzedni/src/moduly/studio/powierzchnia-dokumentu.ts

Powierzchnia niesie zaznaczenie do rdzenia jako zakres znaków treści w zapisie
markdown, bo tak liczy je komenda kontekstowa operacji, a kursor po
przerysowaniu wraca na miejsce w tekście widocznym bloku, bo tam Operator go
widział; dwie miary służą dwóm różnym odbiorcom i zlanie ich przesuwałoby
kursor o długość znaczników markdown. Podział na kartki jest prawdziwy —
liczą go wysokości bloków zmierzone przez przeglądarkę, nie licznik znaków —
bo blok, który nie mieści się na kartce, ma zejść na następną tak, jak zejdzie
w wydaniu. Rdzeń liczy własną paginację osobno, silnikiem swojego składu,
i oddaje ją obrazami stron; podgląd wydania sięga po te obrazy dopiero, gdy
rdzeń je odda, a do tego czasu pokazuje kartki liczone tutaj i mówi to wprost
paskiem stanu, żeby wyrys nieaktualny nie uchodził za wynik ostateczny.

Kartka rdzenia wchodzi obrazem w zapisie `data:`, bo pole `uri` zasobu
magazynu jest ścieżką w systemie plików rdzenia, niedostępną przeglądarce;
bajty przychodzą odpowiedzią komendy pobrania treści zasobu i stąd zapis
`data:`, bez jednego żądania do sieci. Kartki rdzenia tracą aktualność w
chwili zmiany treści, więc wykaz pusty zdejmuje wyrys rdzenia i wraca do
kartek liczonych w oknie — pokazywanie starego wyrysu jako podglądu treści
bieżącej byłoby nieprawdą o dokumencie.

Magazyn nastaw widoku rozróżnia argument pominięty od podanego jako `null`:
pominięty daje magazyn tego urządzenia, `null` znaczy „bez zapisu" — nastawy
działają przez sesję i nikt nie obiecuje, że przeżyją zamknięcie. Nastawy
strony przestawione chwytem linijki albo paskiem widoku idą do osobnego
wykazu nadpisań, nie tylko do nastaw bieżących, bo okno pcha własną kopię
całych nastaw strony przy każdej swojej zmianie; bez wykazu nadpisań ta kopia
zabierałaby Operatorowi to, co właśnie ustawił chwytem. Gdy okno pcha nowe
nastawy, wygrywa wyłącznie to pole, które okno naprawdę zmieniło — nadpisanie
zostaje, dopóki wartość okna się od niego nie różni.

Nastawa gotowa skali liczy się z pola widoku zmierzonego w chwili wywołania,
bo tylko wymiar zmierzony teraz jest prawdziwy; nastawa nieznana albo pole
jeszcze nieosadzone nie zmienia nic, bo skala zgadnięta byłaby gorsza od
nastawy niewykonanej. Wydruk bierze te kartki, które Operator widzi w
podglądzie, nie drugi wyrys ani surowy tekst źródłowy, więc przed drukiem
powierzchnia schodzi do trybu wydania, oznacza strony poza zakresem druku i
zdejmuje adiustację, gdy Operator drukuje pismo do wysłania — wszystkie trzy
kroki wraca po zamknięciu okna drukarki.

Decyzja o zmianie śledzonej stoi w miejscu zmiany, w treści, a nie w osobnym
oknie, bo Operator rozstrzyga tam, gdzie patrzy; pasek zatwierdzenia pod
blokiem używa tej samej drogi decyzji co znaczniki w treści i wykaz w
okienku recenzowania, nie jest trzecim miejscem decyzji. Kotwica komentarza
jest znakiem w treści, nie wpisem w wykazie, bo komentarz przypięty do
fragmentu ma pokazywać, do którego — komentarz bez zakresu kotwicy nie
dostaje, bo dotyczy dokumentu, a nie miejsca w nim. Panele okna stoją nad
powierzchnią jako nakładka, bo powierzchnia należy do dokumentu — stałe
kolumny zostają trybem do wyboru, a znacznik układu idzie na wspólnego
przodka powierzchni i paneli, ponieważ to on rozkłada kolumny.

Krawędzie kolumn tabeli mierzy powierzchnia z wyrysu, nie zgaduje z liczby
kolumn, bo tabela ustawia szerokości wedle treści komórek, więc krawędź na
linijce ma stać tam, gdzie naprawdę stoi krawędź komórki; krawędź ostatniej
komórki jest krawędzią całej tabeli, nie granicą między kolumnami, więc nie
podlega chwytowi przesunięcia. Położenie punktu w zapisie treści liczy się
sumą zapisów bloków poprzedzających i zapisem bloku bieżącego do miejsca
kliknięcia — tak samo, jak treść dokumentu jedzie do rdzenia — a punkt poza
treścią oddaje `null`, nie zero, bo zero jest położeniem prawdziwym i pomyłka
tutaj wysłałaby operację na początek dokumentu.

## budowa/klient-poprzedni/src/moduly/studio/strona-postac-dokumentu.ts

Warstwa postaci spaja panele okna pracy z komendami rdzenia: panele składają
treść żądania z pól i nic więcej wiedzieć nie muszą, a warstwa dokłada do
każdego żądania identyfikator dokumentu, zakres z bieżącego zaznaczenia oraz
autora czynności, woła rdzeń, czyta jego bilans i mówi Operatorowi, co się
stało — łącznie z tym, czego czynność nie zrobiła i przez którą blokadę.
Zakres dokłada się w warstwie, nie w panelu, bo zaznaczenie jest jedno na
moduł i mieszka w stanie — gdyby panel je czytał sam, powstałoby drugie
miejsce wiedzące, nad czym Operator pracuje, i te dwa rozjechałyby się przy
pierwszym przełączeniu zakładki dokumentu.

Autor czynności jedzie jawnie: czynność Operatora odkłada się jako zmiana
autora `uzytkownik`, czynność modelu jako zmiana autora `model`, inaczej
przełącznik pokazujący pracę modelu nie miałby czego podświetlić. Warstwa
zawsze podaje autorem Operatora, bo to jego panel — droga modelu do tych
samych komend idzie osobno, narzędziami modelu, i pola autora tam nie
zapomina. Brak wczytanego dokumentu jest odmową nazwaną, nie ciszą: każda
z komend przyjmuje identyfikator dokumentu jako pole obowiązkowe, bo żądanie
bez dokumentu nie ma czego dotyczyć, a wysłanie go z polem pustym wróciłoby
błędem walidacji, którego Operator by nie zrozumiał — warstwa mówi więc
wprost, czego brakuje i po czyjej stronie leży brak.

## budowa/klient-poprzedni/src/moduly/research/wywolania-komend.ts

Zdanie odpowiedzi komendy mówi o mierzonym skutku czynności — ile pozycji,
jaki plik, ile luk — a nie o tym, że wywołanie się powiodło: „rdzeń oddał
wynik" jest zdaniem, po którym Operator nadal nie wie, czy coś się stało.
Akcje panelu niosą w polu kodu dokładnie nazwę komendy rdzenia, więc
rozdzielnik komend tego pliku jest odwzorowaniem jeden do jednego i nie ma
w nim ani jednej nazwy pisanej z ręki. Nie ma w nim natomiast żądań, których
nie da się złożyć bez tekstu od Operatora — zapytania wyszukiwania, treści
notatki, uzasadnienia odrzucenia — bo takie pozycje mają mówić wprost, czego
brakuje, zamiast wysyłać żądanie z polem pustym i wracać odmową walidacji,
z której nic nie wynika.

## budowa/klient-poprzedni/src/moduly/terminal/okno-task-schedule.ts

Plan zadania powłoki jest drugą powierzchnią jednej rodziny komend automatyki,
nie drugą rodziną komend: cykliczność wiąże się z automatyką, zapisuje ją
`automation.schedule.set`, czyta `schedule.get`, a rodziny `terminal.schedule.*`
nie ma i nie będzie, bo dwie rodziny na jeden harmonogram byłyby dwiema
prawdami o jednym bycie. Zaplanowanie zadania idzie dwoma krokami tej rodziny:
`automation.workflow.save` zapisuje automatykę o jednym kroku rodzaju
`command` wołającym `terminal.command.exec`, a `automation.schedule.set`
dokłada jej wyrażenie cron; edytora automatyk to okno nie stawia, bo pełna
praca na krokach, wersjach i zmiennych zostaje w module Automations.

Plan nie jest budzikiem: rdzeń wylicza chwilę najbliższego uruchomienia, ale
nie ma czym odpalić automatyki samodzielnie, bo kontrakt nie zna komendy ani
zdarzenia, którym harmonogram zgłaszałby wyzwolenie — zapisany plan jest
zapisem obowiązującym wraz z terminem, a uruchomienie prowadzi Operator
z okien modułu Automations, i potwierdzenie zapisu mówi to wprost. Kolejki
okno nie zakłada, bo `queue.create` żąda identyfikatora karty sesji, którego
rdzeń jeszcze nie wypełnia. Potok jest sekwencjonowaniem po stronie klienta,
nie silnikiem w rdzeniu: okno wysyła krok, czeka na jego domknięcie odczytem
wyjścia i dopiero wtedy decyduje o kroku następnym; zamknięcie okna w trakcie
przerywa sekwencjonowanie, ale nie krok już uruchomiony — ten kończy się
w rdzeniu i widać go w Process Monitorze.

Granica czekania na domknięcie kroku jest dłuższa niż w podglądzie wyjścia
Process Monitora, bo tu czekanie jest warunkiem decyzji: bez kodu wyjścia
kroku nie ma jak rozstrzygnąć, czy wolno ruszyć z krokiem następnym. Krok
automatyki wołający polecenie powłoki niesie nazwę komendy kontraktu i treść
jej żądania, nie polecenie powłoki wprost, dzięki czemu plan przechodzi tę
samą bramę uprawnień i ten sam egzekutor izolacji, co polecenie wydane ręcznie
z karty; krok omijający komendę byłby drugą drogą do powłoki, bez żadnego
z tych sprawdzeń.

## budowa/klient-poprzedni/src/moduly/translate/okno-warsztat-translate.ts

Rodziny komend spoza czterech okien pierwotnych dzielą trzy wskazania — okno
tłumaczenia, panel języka i ścieżkę pliku — więc stoją w jednym warsztacie,
a nie w dwunastu osobnych oknach: rozbicie kazałoby Operatorowi wpisywać te
same trzy wskazania po kilkanaście razy, a każda rodzina prowadziłaby własną
kopię wykazu paneli. Tutaj wskazania stoją raz, na górze, i wchodzą do
każdego żądania. Każda sekcja kończy się zdaniem o mierzonym skutku, nie
słowem „gotowe" — liczbą par wniesionych do pamięci, ścieżką pliku, który
powstał, liczbą pozycji przebiegu — bo komunikat „udało się" bez liczby jest
meldunkiem zamiast skutku.

## budowa/klient-poprzedni/src/moduly/studio/przybornik-znakowania.ts

Znaczenie fragmentu jest czynnością powtarzaną kilkadziesiąt razy na
dokument, a droga przez menu kosztuje przy każdym razie dwa ruchy więcej,
dlatego wszystko, czym się dokument znaczy, stoi w bocznym przyborniku przy
krawędzi treści, nie w menu ani osobnym oknie. Przybornik jest narzędziem do
pracy na fragmentach, nie na dokumencie w całości: gdy zaznaczenia nie ma,
pozycje nie milkną, tylko mówią, że znakowanie obejmie cały dokument, i to
jest odpowiedź, nie odmowa.

Wszystko, czym Operator znaczy dokument, model musi umieć założyć i
przeczytać: komentarz i adnotacja idą komendami rdzenia niosącymi pole
autora, więc znakowanie modelu jest podpisane jako `model`, a Operatora jako
`uzytkownik`. Znacznik własny pola autora w kontrakcie nie ma, więc
przybornik trzyma go u siebie wraz z autorem — to brak nazwany, nie
zatajony. Czynności bez zaplecza w rdzeniu stoją w osobnej części, z
nazwaniem, co dokładnie brakuje — przypis, odsyłacz, odwołanie wzajemne,
wstawienie tabeli i pola należą do aparatu i postaci dokumentu, których
kontrakt jeszcze nie niesie — bo przycisk wychodzący do rdzenia po komendę,
której nie ma, byłby uprzejmą odmową udającą funkcję.

Zaplecze rdzenia znakowania jest nieobowiązkowe, bo przybornik stoi w oknie
pracy z dokumentem, którego montaż powstaje etapami: zaplecze podane znaczy,
że droga do rdzenia jest wpięta i znakowanie jest wtedy trwałe oraz widoczne
modelowi, a zaplecze `null` znaczy, że tej drogi jeszcze nie wpięto, i
przybornik pisze to Operatorowi wprost, zamiast pokazywać przyciski, które
nic nie robią. Przycisk dyktafonu, wzorem pakietu biurowego, jest przyciskiem
dyktowania do treści dokumentu; przybornik go tylko osadza, a nagrywanie
i przepisanie liczy zaplecze mowy, jednym rachunkiem wspólnym z mikrofonem
wiersza polecenia.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-arkusza.test.ts

Sprawdziany panelu arkusza stylów oraz panelu list i znaków mierzą treść zgłoszonego żądania,
a nie samo naciśnięcie przycisku. Trzy zachowania są sprawdzane osobno jako miejsca dawnej szkody
w tym produkcie. Cecha logiczna niesie trzy stany, nie dwa: naniesienie pogrubienia nie zdejmuje
przy okazji kursywy, której operator nie wskazał. Czynność bez ani jednego wypełnionego pola jest
odmawiana, a nie wysyłana jako żądanie, które niczego nie zmieni. Nastawa poziomu listy wymaga
wskazania listy zastanej — bez niej okno nazywa brakujący element zamiast wysyłać żądanie z pustym
identyfikatorem.

## budowa/klient-poprzedni/src/moduly/studio/szablon-panel.ts

Zapis szablonu z bieżącego dokumentu przenosi domyślnie blokady fragmentów
wzorcowych — pole `includeLocks` bez wartości znaczy „tak", ponieważ fragmenty
wzorcowe pisma mają pozostać wzorcowe także w dokumentach założonych z tego
szablonu.

Szablonu fabrycznego nie da się usunąć: rdzeń odmawia i zwraca powód odmowy
w polu `deleted`. Panel zostawia pozycję w wykazie, ponieważ szablon nadal
istnieje — zdjęcie jej z widoku przy jednoczesnym powrocie po kolejnym odczycie
byłoby pokazaniem skutku, który się nie wydarzył.

## budowa/klient-poprzedni/src/moduly/studio/osadzenie-panel.ts

Panel osadzony pokazuje dokument przez cały czas pracy: Operator nie opuszcza
edytora, a wybór położenia panelu — obok treści albo na całej powierzchni —
należy do Operatora. Wniesienie fragmentu z Biblioteki albo ze strony
sieciowej trafia wprost w miejsce kursora w dokumencie, nie do kolejki ani do
zasobów; komenda `studio.ingest.url` oddaje pole tekstowe, które panel
prowadzi do dokumentu, a gdy tekst jeszcze nie istnieje, bo pozycja czeka na
rozpoznanie pisma, panel nazywa ten stan zamiast wnosić pustkę.

Komendy `studio.insert.from.library` i `studio.insert.from.web` wnoszą
fragment do dokumentu i oddają zapis pochodzenia — adres albo plik, wersję,
czas sięgnięcia i autora. Panel woła te komendy, gdy dostał źródło wstawień,
i pokazuje pochodzenie oddane przez rdzeń wraz z bilansem czynności. Wiersz
pochodzenia wnoszony do treści dokumentu pozostaje drogą zapasową, ponieważ
przeżywa wydanie dokumentu do formatu, który zapisu pochodzenia nie niesie.

Panel nie rysuje strony sieciowej. Migawka `browser.snapshot.get` oddaje
adres, tytuł, treść renderowaną i źródło strony jako tekst, nie jako obraz;
zrzut ekranu jest odnośnikiem zasobu, a komendy pobierającej jego bajty do
przeglądarki kontrakt nie niesie. Panel pokazuje więc treść, którą rdzeń
naprawdę oddaje, i nazywa wprost to, czego nie oddaje. Ramka z cudzą stroną
wewnątrz okna byłaby drugą przeglądarką, a moduł Browser jest jeden.

## budowa/klient-poprzedni/src/moduly/studio/widok-pasek-widoku.ts

Panele nastaw wchodzą na żądanie i schodzą, gdy nie są używane; stała kolumna
zabierałaby kartce szerokość na stałe, dlatego pasek jest wąskim rzędem przy
krawędzi powierzchni, a nastawy rzadsze — układ kartek, przewijanie, jednostka
linijki, tryb dwóch dokumentów — stoją w nakładce rozwijanej przyciskiem
„Nastawy widoku". Nakładka stoi nad treścią i schodzi naciśnięciem, klawiszem
Escape albo naciśnięciem poza nią, nie zabierając kartce ani milimetra
szerokości.

Nastawy widoku stoją przy powierzchni, a nie wyłącznie na wstążce, ponieważ
dotyczą tego, na co Operator patrzy, i sięga po nie stale — skala oraz skok
o stronę są czynnościami ciągłymi, nie wyprawą na osobną zakładkę. Wstążka
niesie te same nastawy jako gniazdo tego samego paska, więc nie powstają dwa
miejsca, które mogłyby się rozjechać.

Szybkie drukowanie stoi w rzędzie stałym, nie w nakładce: jedno naciśnięcie
uruchamia druk ostatnimi nastawami, bez okna nastaw, tak jak pasek szybkiego
dostępu pakietu biurowego. Bez sterownika druku przycisk pozostaje widoczny
i nazywa powód niedostępności zamiast milczeć.

Tryb źródłowy ze znacznikami i podgląd wydruku wykluczają się: znaczniki nie
mają paginacji, a podgląd jej wymaga. Wyłączony przełącznik mówi to wprost,
zamiast oddawać widok, którego Operator nie zamawiał.

## budowa/klient-poprzedni/src/moduly/studio/strona-panel-nastaw.ts

Numeracja stron i marginesy stoją w tym panelu jako nastawy zapisywane osobnymi
komendami rdzenia, nie jako pojedynczy przełącznik profilu wydania — numeracja
niesie styl, umiejscowienie i numer startowy, a marginesy przeżywają zapis
zamiast żyć wyłącznie przez czas trwania sesji.

Panel jest nakładką otwieraną przyciskiem i domyślnie schowaną, ponieważ
powierzchnia należy do dokumentu — stała kolumna nastaw zajmowałaby ją na
stałe, choć Operator sięga po te nastawy rzadko.

Każda nastawa strony da się ustawić osobno dla sekcji, dlatego wybór sekcji
stoi raz, u góry panelu, i dotyczy wszystkich pól poniżej: pismo z załącznikiem
w orientacji poziomej pozostaje wtedy jednym dokumentem, nie dwoma, a Operator
nie wskazuje sekcji osobno przy każdej nastawie.

Panel nie woła rdzenia i nie zna dokumentu, na którym pracuje Operator — składa
treść żądania z wypełnionych pól i oddaje ją wywołującemu. Identyfikator
dokumentu, autor czynności i odczyt odpowiedzi rdzenia należą do warstwy
wyżej, żeby panel nie stał się drugim miejscem, które śledzi, nad czym
Operator pracuje.

## budowa/klient-poprzedni/src/moduly/translate/katalog-funkcji-translate.ts

Wykaz istnieje po to, żeby stan modułu dało się przeczytać, a nie zgadnąć. Okno
pokazuje czynności, które wykonuje; katalog pokazuje komplet zamierzony
w opracowaniu modułu i przy każdej pozycji mówi jedno z dwojga: którą komendą
kontraktu jest wykonywana albo czego brakuje, żeby była. Pozycja bez komendy
nie znika z wykazu, bo zniknięcie byłoby ukryciem braku.

Nazwy pozycji i podział na grupy pochodzą z opracowania modułu i nie są
tłumaczone ani parafrazowane. Numeracji opracowania katalog nie przenosi:
pozycję odnajduje się po pełnej nazwie, nie po oznaczeniu, które poza
dokumentem nic nie znaczy.

Katalog jest zbiorem danych, nie widokiem — wyszukiwarka funkcji buduje z niego
listę, a okna sięgają po pozycje swojej grupy.

## budowa/klient-poprzedni/src/moduly/terminal/okno-session-manager.ts

Do rdzenia idzie jedna komenda otwarcia karty, z powłoką ssh i adresem celu
podanym zmienną środowiska SSH_TARGET — tak i tylko tak rdzeń przekazuje adres
programowi ssh. Karta powstaje w tym samym stanie modułu, w którym stoją karty
otwarte w oknie Terminal Tabs, więc połączenie z tego okna jest od razu
widoczne w oknie wiodącym i w konsoli wyjścia poleceń.

Książka hostów, klucze SSH, tunele portowe i wykaz kart są bytami rdzenia,
a okno jest ich widokiem. Wpis hosta zapisuje i czyta rodzina komend hosta,
klucze prowadzi rodzina komend klucza, tunele — rodzina komend tunelu, a karty
pokazuje wykaz sesji rdzenia, więc wykaz sięga dalej niż pamięć tego
połączenia: rdzeń odtwarza karty przy starcie i okno je widzi po ponownym
podłączeniu gniazda.

Czytanie i pisanie pliku konfiguracyjnego OpenSSH zostaje, bo służy czemu
innemu niż trwałość: wnosi wpisy z maszyny operatora i wynosi je z powrotem.

Pozycja panelu akcji bez odpowiadającej jej komendy kontraktu — polityka
znanych hostów — stoi jawnie nieczynna wraz z powodem liczonym z odczytu
wykazu komend rdzenia, tak samo jak każda inna pozycja zależna od pokrycia
komend.

## budowa/klient-poprzedni/src/uwierzytelnienie/ekran-logowania.ts

Ekran jest przesłoną nad aplikacją, nie osobną trasą: kładzie się nad
gospodarzem dokumentu, a aplikacja pod nim składa się i łączy z rdzeniem
w tym samym czasie, więc po wejściu przesłona znika i widoczne jest gotowe
Centrum dowodzenia zamiast drugiego ładowania. Trasy aplikacji zostają
nietknięte, ponieważ bramka nie jest miejscem pracy.

Ekran nie zakłada budzika, nie liczy czasu sesji i nie przerywa pracy
pytaniem o tożsamość. Rdzeń nie odcina komend po wygaśnięciu sesji — bramka
jest progiem wejścia, nie strażnikiem każdego żądania — więc wygaśnięcie
ujawnia się wyłącznie przy następnym uruchomieniu, zdaniem nad formularzem.

Który z dwóch formularzy pokazać — wejście hasłem czy pierwsze ustawienie
hasła — rozstrzyga rdzeń, nie domysł klienta. Pole nie jest blokowane,
przycisk nie jest wyszarzany; jedyny sprawdzian po stronie formularza to
zgodność hasła z powtórzeniem przy zakładaniu, ponieważ kontrakt opisuje
powtórzenie jako sprawę formularza klienta.

Wymóg logowania jest nastawą `gateway.requireLogin` — gdy jest wyłączony,
przesłona nie staje.

Stan pola „nie wyloguj mnie” bierze się z miejsca, w którym leży sesja
zapisana poprzednio, a nie ze stałej: raz odznaczone pole ma zostać
odznaczone. Zaznaczone pole kładzie sesję w pamięci trwałej, odznaczone
w pamięci okna, a ta sama wartość idzie do rdzenia i rozstrzyga o trwaniu
sesji.

Odnośnik resetu hasła prowadzi do odzyskania konta, nie do zmiany hasła ze
znanym hasłem bieżącym: naciska go ten, kto hasła nie pamięta, a wtedy
zmiana ze znanym hasłem jest drogą donikąd. Zmiana hasła ze znanym hasłem
jest czynnością Ustawień.

Odnośnik potwierdzenia z listu wprowadza krok drugi rejestracji i pozwala
z niego wyjść. Bez tego odnośnika krok potwierdzenia prowadziłby wyłącznie
z udanej rejestracji w tym samym oknie; odświeżenie strony między listem
a przepisaniem drogi zostawiałoby Operatora przed formularzem wejścia bez
żadnej drogi dalej, ponieważ rejestracja odmawia wtedy kodem `conflict`,
a logowanie oczekiwaniem na potwierdzenie adresu. To samo dotyczy listu
odczytanego na innej maszynie.

## budowa/klient-poprzedni/src/moduly/terminal/okno-process-monitor.ts

Okno Process Monitor pokazuje na żywo rejestr procesów rdzenia terminala, w tym proces
zainicjowany poleceniem sztucznej inteligencji, zgodnie z rejestrem procesów rdzenia serwera.
Aktualizacja na żywo znaczy ze zdarzeń: okno odczytuje wykaz przy wejściu i na wyraźne żądanie,
a każdą późniejszą zmianę przynosi zdarzenie rdzenia. Odpytywanie w pętli dałoby ten sam obraz
drożej i z opóźnieniem, a przy stu procesach zalałoby gniazdo. Filtr stanu i inicjatora jedzie do
rdzenia parametrem komendy, bo tak stanowi kontrakt i tak wynik jest spójny z dziennikiem rdzenia;
grupowanie jest wyłącznie porządkiem wyświetlania i zostaje w oknie. Wyjście na żywo idzie do
wspólnego bufora konsoli wyjścia, ale bufor żyje jedno połączenie: po rozłączeniu i ponownym
podłączeniu ma zero fragmentów, choć rdzeń wciąż oddaje pełną treść. Osobny przycisk pozycji pyta
o wyjście procesu wprost i pokazuje je przy pozycji, a nie w buforze, więc żaden wiersz nie
wchodzi do konsoli dwa razy.

Okno nie ma dziś ani jednej pozycji bez pokrycia w rdzeniu: wstrzymanie procesu było ostatnią
i już stoi przy wierszu wykazu.

Wstrzymanie albo wznowienie procesu w rdzeniu niesie pole odpowiedzi mówiące, czy system maszyny
rdzenia w ogóle zna wstrzymanie obcego drzewa procesów. Fałsz tego pola znaczy, że proces został
nietknięty, ponieważ platforma tego nie umie — co jest czymś innym niż niepowodzenie czynności.
Okno mówi to wprost, zamiast pokazywać powodzenie przy procesie, który dalej zajmuje procesor.

Czas oczekiwania rdzenia na domknięcie procesu przed odczytem wyjścia wynosi trzy sekundy.
Kontrakt dopuszcza sześćdziesiąt tysięcy milisekund, ale czekanie zajmuje zadanie gniazda,
a minuta bez odpowiedzi wygląda jak zawieszone okno. Trzy sekundy wystarczają, aby polecenie
krótkie zdążyło się domknąć i oddało komplet wraz z kodem wyjścia; polecenie długie i tak odda
wyjście dotychczasowe ze stanem biegnącym, bo czekanie nie jest warunkiem odpowiedzi, a kto chce
zobaczyć resztę, odczytuje ponownie.

Zwinięcie podglądu wyjścia nie pyta rdzenia: drugie kliknięcie zdejmuje treść z widoku i mówi to
wprost, inaczej nie dałoby się odróżnić zwinięcia od odczytu, który wrócił pusty.

Zdanie potwierdzenia odczytu wyjścia mierzy treść, nie sam fakt odpowiedzi. Wyjście puste jest
przebiegiem udanym, bo polecenie mogło nic nie wypisać, ale różni się od wyjścia niepustego —
inaczej potwierdzenie znaczyłoby to samo w obu przypadkach. Rozjazd stanu też idzie wprost: wykaz
w oknie jest kopią wcześniejszego odczytu, a odpowiedź na odczyt przychodzi z tej chwili, więc gdy
się różnią, świeższa jest odpowiedź.

Odmowa rdzenia przy braku wyjścia procesu idzie dosłownie, ponieważ sama nazywa dokładny powód:
czy proces nigdy nie ruszył, wypadł z historii, czy rdzeń był uruchomiony ponownie. Parafraza
zgubiłaby wszystkie trzy powody.

Złożenie powierzchni rejestru procesów jest konstrukcją czystą: nie domyka się na stanie okna ani
na rdzeniu, więc dała się wyjąć bez przenoszenia zależności.

## budowa/klient-poprzedni/src/moduly/roundtable/zrodlo-arsenalu.ts

Rozdział komend modułu Roundtable na dwa pliki idzie po roli, nie po
wielkości. Cztery komendy źródła debaty prowadzą samą debatę: dodają
uczestnika, otwierają turę, moderują ją i czytają stanowisko. Komendy tego
pliku pracują nad zapisem debaty — czytają go, analizują, oceniają, wydają.
Okno rozmowy może stać wyłącznie na komendach źródła debaty; rozszerzenia
boczne stoją na komendach tego pliku.

Każda czynność oddaje typ Wynik, nie samą treść, i każda sprawdza kształt
odpowiedzi. Powód jest ten sam co w źródle debaty: obszar odmawia z powodów
zwyczajnych — kanał uczestnika bywa nieczynny, głosowanie bywa zamknięte —
a okno musi odróżnić brak treści od nieudanego zapytania, więc żaden
odczyt nie zastępuje braku odpowiedzi pustą wartością domyślną.

## budowa/klient-poprzedni/src/moduly/roundtable/czynnosci-arsenalu.ts

Plik istnieje po to, żeby żadna komenda obszaru nie została bez drogi z okna.
Wcześniej pozycje bez obsługi stały jako przyciski nazywające brak; obecnie
każda z nich naprawdę woła swoją komendę, a wynik melduje w oknie, z którego
padła.

Dane żądania pochodzą ze stanu debaty wspólnego oknom modułu: okno, tura
bieżąca, skład, wypowiedzi. Czynność, dla której stan nie ma jeszcze wskazania
— nie ma tury, nie ma uczestników, nie ma wypowiedzi — nie idzie do rdzenia po
to, żeby dostać odmowę: melduje brak wskazania od razu i nazywa, czego
brakuje. Rdzeń odmówiłby tak samo, tylko po podróży tam i z powrotem.

Plik nie podstawia wartości domyślnych za Operatora. Tam, gdzie komenda
potrzebuje treści — stanowisko, zdanie odrębne, warianty głosowania — treść
przychodzi z pola okna, a przycisk bez wypełnionego pola melduje, czego
brakuje.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-arkusza.ts

Panel prowadzi trzynaście czynności na fragmencie dokumentu: style nazwane
(wykaz, założenie, zmiana, stosowanie, usunięcie), styl znaku, styl akapitu,
czyszczenie formatowania, wielkość liter, malarz formatów, zaznaczenie wedle
podobnego formatowania, zamianę wraz z postacią oraz tabulator zakładany
liczbą. Praca na fragmentach jest osią całego zamówienia.

Zapis stylu nazwanego przestawia wszystkie miejsca dokumentu, które danego
stylu używają. Panel mówi to wprost przy przycisku i oddaje liczbę
przestawionych miejsc z bilansu odpowiedzi, żeby Operator widział skutek
zapisu, a nie samo potwierdzenie wykonania.

Malarz formatów prowadzi dwie czynności rdzenia: pobranie postaci oddaje
uchwyt, naniesienie stosuje go gdzie indziej. Uchwyt pamięta warstwa wyżej;
panel pokazuje wyłącznie, czy coś jest pobrane.

Panel nie woła rdzenia i nie zna dokumentu ani zaznaczenia. Składa treść
żądania z pól bez identyfikatora dokumentu i bez zakresu — zakres dokłada
warstwa wyżej z bieżącego zaznaczenia, bo ona jedna wie, co Operator
zaznaczył.

## budowa/klient-poprzedni/src/moduly/studio/wstazka-pracy.ts

Czynności wstążki stoją zebrane wedle rodzaju pracy, wzorem pakietu biurowego,
zamiast jednego paska ze stoma czynnościami naraz: taki pasek byłby wykazem,
po którym Operator szuka wzrokiem, a zakładka mówi wprost, gdzie czego szukać,
i nazwana grupa — dlaczego te czynności stoją razem.

„Asystent" jest osobną zakładką wstążki, nie pozycją w menu: niesie suwaki
koncepcyjne, wykaz operacji z pliku kategorii operacji i decyzję o wyniku.
Wykaz operacji nie powstaje w tym module od nowa — pochodzi z tego samego
pliku, co wykaz Tools Panelu, żeby obie powierzchnie pokazywały ten sam zestaw.

Wstążka nie woła rdzenia i nie zna stanu modułu. Buduje kontrolki i zgłasza
naciśnięcia oknu przez czynności wstążki, a złożone kawałki — formularz
wczytania, panel znajdź/zamień, pola różnicy — przyjmuje gotowe jako gniazda.
Dzięki temu jedno miejsce trzyma układ wstążki, a inne prowadzi rozmowę
z rdzeniem.

## budowa/klient-poprzedni/src/moduly/studio/kopie-panel.ts

Panel łączy cztery wymagania trwałości dokumentu w jednym miejscu. Autozapis:
odstęp i zapis przy zdarzeniach okna jako jawne, odwracalne ustawienie;
obejmuje treść i postać dokumentu i idzie osobnym szeregiem wersji, żeby nie
zaśmiecał historii Operatora. Kopia zapasowa: zakładana przed zapisem
i niezależnie od historii wersji, żeby przetrwała awarię procesu i awarię
zapisu. Przywrócenie po nagłym zamknięciu: Studio zgłasza je samo, zdaniem
o niezapisanym dokumencie z podaniem godziny, zamiast czekać, aż Operator się
domyśli; przywrócenie do nowego dokumentu stoi obok przywrócenia na miejsce,
żeby przywracanie nie kasowało tego, co już jest. Powrót do stanu pierwotnego:
jedno polecenie, bez szukania wersji założycielskiej w wykazie; wersje nowsze
zostają, więc powrót jest odwracalny.

Wskaźnik stanu zapisu ma cztery wartości i „nieudany" jest jedną z nich.
Nieudany zapis samoczynny wraca odpowiedzią udaną z polem `saved: false` —
panel pokazuje go jako niepowodzenie wraz z powodem, bo pokazanie „zapisano"
po nieudanym zapisie kosztowałoby Operatora pracę.

Panel woła rdzeń sam; treść i postać do odłożenia bierze z kontekstu, bo
powierzchnia dokumentu jest po stronie okna.

## budowa/klient-poprzedni/src/moduly/studio/styl-panel-list.ts

Wzór numeru numeracji jest polem, nie wyborem z wykazu, ponieważ numeracja
prawnicza wielopoziomowa typu 1.1.2 jest wzorem, nie pozycją wykazu gotowych.
Format `legal` mówi rdzeniowi rodzaj numeracji, a wzór mówi jej postać —
pismo urzędowe wymaga obu wartości naraz i nie da się go zamknąć w wykazie
gotowych wzorów.

Znaki tablicy są znakami ostatnio użytymi, nie ulubionymi, ponieważ kontrakt
oddaje pole `recentlyUsed` przy każdym znaku i to ono rozstrzyga, co stoi pod
ręką. Panel nie zakłada drugiego pojęcia obok tego pola: wykaz jest jeden
i pochodzi z rdzenia, więc znaki ostatnio użyte przeżywają zamknięcie karty.

Panel nie woła rdzenia i nie zna dokumentu — składa treść żądania i oddaje ją
warstwie wyżej, która dokłada dokument i zakres przed wywołaniem komendy.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-wstawien-studio.ts

Siedemnaście komend studio wnoszą treść do dokumentu i wydają go dalej, obejmując wejście,
wydanie, wniesienie ze źródła i cyfryzację. Wszystkie prowadzą tę samą drogę: coś spoza
dokumentu wchodzi do dokumentu albo dokument wychodzi na zewnątrz. Założenie dokumentu zakłada
pustą stronę, wniesienie pliku i PDF wnosi materiał wprost do edytora, wniesienie obrazu stawia
go w miejscu kursora, wniesienie z biblioteki i ze strony sieci niesie fragment wraz z zapisem
pochodzenia, a cyfryzacja prowadzi kolejkę wczytywania po stronie rdzenia. Wydanie do formatu,
wydanie wsadowe, zapis pod nazwą i kopia zamykają drogę w drugą stronę.

Bilans nie jest ozdobą odpowiedzi. Rdzeń oddaje przy wniesieniu bilans odzyskania, przy
czynności bilans zmiany, a przy wydaniu wykaz cech pominiętych przez format docelowy. Zdania
składające te trzy bilanse stoją w jednym miejscu, bo czytają je wszystkie panele wstawień
naraz. Milczące zgubienie tabeli przy wydaniu do tekstu czystego jest ciszą niedopuszczalną,
więc okno nie ma drogi, którą mogłoby bilans pominąć.

Wywołanie uczciwe osłania każdą komendę, ponieważ koperta nieznanej komendy studia nie niesie
pola stanu i bez tej osłony okno stałoby w ładowaniu bez końca.

Zdanie o bilansie czynności traktuje zero zmian jako wynik prawdziwy, nie jako powodzenie:
czynność, która nie tknęła ani jednego miejsca, jest odpowiedzią „nic się nie stało" i tak ma
być przeczytana. Pominięcie z powodu blokady nazywa tę blokadę, aby było wiadomo, która zapora
zatrzymała czynność, zamiast pozostawiać wrażenie częściowego powodzenia.

Zdanie o bilansie wniesienia odzwierciedla to, że odzyskanie z formatu PDF jest odtworzeniem,
nie odczytem: liczby stron z warstwą tekstową i bez niej, tabel rozpoznanych i nierozpoznanych,
obrazów osadzonych i pominiętych stoją obok siebie. Skierowanie na rozpoznanie pisma jest
zdaniem osobnym, ponieważ to inna decyzja niż samo dalsze wniesienie materiału.

Zdanie o wyniku wydania nazywa wprost stratę cech, gdy format docelowy jest uboższy niż
dokument: format uboższy jest sytuacją normalną, ale przemilczenie straty nie jest, więc
wydanie, które zgubiło tabelę albo przypisy, mówi to wprost zamiast pokazywać samo „zapisano".

## budowa/klient-poprzedni/src/moduly/studio/obiekt-panel.ts

Panel nie rysuje ani jednego kształtu i nie prowadzi biblioteki ikon. Komenda
studio.object.insert przyjmuje rodzaj kształtu i nazwę ikony, a rachunek stoi
w module Design, który zakłada kształt od razu jako węzły ścieżki albo
wyszukuje ikonę. Wskazanie węzła Designu idzie osobnym polem — to jedyna droga
osadzenia kształtu złożonego w Designie.

Rodzaj „wykres" stoi w kontrakcie, ale rdzeń go odmawia, ponieważ rachunku
wykresu po stronie Studia nie ma. Odmowa jest widoczna zanim Operator naciśnie
kontrolkę i mówi, co zrobić zamiast tego: złożyć wykres w module Design
i osadzić go w dokumencie jako obiekt wskazany węzłem. Kontrolka kończąca się
odmową rdzenia byłaby obietnicą bez pokrycia, a jej ukrycie zabrałoby
Operatorowi wiedzę, że taka droga istnieje.

## budowa/klient-poprzedni/src/moduly/terminal/okno-terminal-tabs.ts

Do rdzenia idą dwie komendy: otwarcie karty, obejmujące nową kartę i duplikat,
oraz wysłanie polecenia. Podział widoku, rozmiar pisma, schemat barw,
przypięcie, nazwa karty i eksport transkryptu są czynnościami widoku — kontrakt
nie ma dla nich komend, bo nie zmieniają niczego po stronie rdzenia. Zamknięcie
karty jest wyjątkiem opisanym niżej.

Zdanie potwierdzenia po otwarciu karty albo uruchomieniu polecenia składa się
wyłącznie z pól odpowiedzi rdzenia — stan karty, katalog roboczy, stan procesu,
polecenie — i rozstrzyga o powodzeniu tym, co rdzeń oddał, a nie tym, że
żądanie poszło. Czynności czysto widokowe mówią wprost, że rdzeń o nich nie
wie.

Zamknięcie karty w rdzeniu ma dziś komendę kontraktu, ale okno jej jeszcze nie
wywołuje. Pozycja mówi o tym zdaniem liczonym z odczytu wykazu komend rdzenia,
a nie stałym napisem: napis orzekałby o stanie rdzenia z chwili pisania kodu
i nie przestałby go orzekać w dniu, w którym rdzeń dostanie odpowiednią
obsługę.

Zamknięcie karty zachodzi najpierw w rdzeniu, potem w widoku. Kolejność jest
rozstrzygnięciem: karta zdjęta z ekranu przed odpowiedzią rdzenia znikałaby
Operatorowi także wtedy, gdy rdzeń zamknięcia odmówił, a wtedy powłoka
biegłaby dalej bez żadnego widoku na siebie. Przypięcie jest bramą, nie
ozdobą: pasek kart zapowiada, że karty przypiętej nie zamyka się jednym
kliknięciem, więc panel akcji odmawia tak samo jak krzyżyk na pasku. Procesy
karty zostają biegnące — zamknięcie zakładki nie ma prawa przerwać
budowania, które trwa trzecią minutę; zakończenie procesów jest osobną
czynnością.

Rozmiar pisma i schemat barw stawiane są na węźle modułu, a nie na węźle
okna: pismo konsoli ma być jedno w całym module, bo podgląd wyjścia
w monitorze procesów niesie ten sam strumień co ogon karty. Węzeł modułu
bierze się z drzewa, bo okno nie zna swojego gospodarza w chwili budowy.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-warsztatu-translate.ts

Warsztat modułu translate jest jednym źródłem dla wszystkich rodzin komend, które nie mieszczą
się w czterech oknach pierwotnych: pamięci jako bytu widocznego dla operatora, segmentacji,
terminologii, korekty i spójności, profili kontroli jakości, obiegu zatwierdzeń, dokumentów,
lokalizacji oprogramowania, napisów i dubbingu, silników, polityki pivota, przebiegu
pakietowego oraz wymiany zewnętrznej. Jedno źródło zamiast jednego pliku na okno wynika
z tego, że wszystkie te komendy idą tą samą drogą wywołania i różnią się wyłącznie nazwą
i ładunkiem; rozbicie ich na dwanaście plików dałoby dwanaście kopii tej samej obudowy, z których
każda musiałaby osobno pamiętać o odmowie komendy nieznanej.

Sprawdzian kształtu towarzyszy każdej komendzie i nie jest formalnością: odróżnia rdzeń, który
odpowiedział wynikiem, od rdzenia, który odpowiedział pustą kopertą. Okno przyjmujące pustą
kopertę jako wynik pokazałoby pustkę jako skutek, czyli dokładnie tę szkodę, przed którą stoi
cały moduł.

## budowa/klient/src/wejscie/wejscie.test.ts

Sprawdzian obejmuje przebieg bez rdzenia uruchomionego: każda odsłona, w tym
odsłona błędu połączenia i odsłona zwłoki nałożonej przez rdzeń, jest tu
osiągalna naprawdę, a nie tylko opisana. Rozmowę z rdzeniem naprawdę
uruchomionym mierzy osobny sprawdzian, który rdzenia wymaga.

## budowa/klient-poprzedni/src/moduly/studio/blokada-panel.ts

Panel łączy blokady fragmentów, zajęcia wykonawców i nastawy pracy kilku agentów naraz, bo
wszystkie trzy odpowiadają na jedno pytanie: czego modelowi nie wolno tknąć i kto teraz pisze
po którym akapicie. Blokada jest trwała i skierowana przeciw modelowi; zdejmuje ją wyłącznie
operator konta, a zasięg obejmujący także operatora jest osobnym, jawnym ustawieniem, nie
zachowaniem domyślnym. Zajęcie fragmentu jest czasowe i skierowane przeciw drugiemu wykonawcy;
odmowa nazywa wykonawcę i czas, bo cicha odmowa kazałaby zgadywać, dlaczego fragment nie
drgnął. Nastawy rozstrzygają, czy pętla wykonawcza i praca wielu agentów w ogóle stoją, ilu
wykonawców pracuje naraz i co się dzieje przy spięciu; oba narzędzia są domyślnie wyłączone
i włącza je operator konta, nie okno. Panel woła rdzeń sam, bo skutkiem każdej z tych czynności
jest wykaz albo odmowa nazwana — jedno i drugie jest treścią widoczną, nie wartością pośrednią
przekazywaną wyżej.

## budowa/klient-poprzedni/src/moduly/terminal/okno-script-library.ts

Do rdzenia idą dwie komendy uruchomienia: jedna wykonuje treść skryptu,
a druga oddaje wynik sprawdzenia składni. Rdzeń podaje treść programowi
powłoki jako pojedynczy argument, więc skrypt wieloliniowy wykonuje się bez
zapisywania go do pliku — dlatego biblioteka nie potrzebuje ani komendy
zapisu pliku, ani ścieżki na dysku serwera.

Każdy zapis pozycji zakłada kolejną wersję, więc poprzednie brzmienie treści
zostaje w dzienniku wersji; usunięcie zabiera pozycję wraz ze wszystkimi jej
wersjami naraz.

Analizę treści prowadzi program leżący na maszynie rdzenia, osobny od
kontroli wstępnej okna. Odpowiedź niesie pole mówiące, czy program analizy
był dostępny, i okno je pokazuje, bo pusty wykaz uwag przy braku programu
znaczyłby fałszywie „treść bez zastrzeżeń". Kontrola wstępna okna zostaje
obok analizy, nie zamiast niej: sprawdza to, co da się rozstrzygnąć bez
żadnego programu, i robi to w chwili pisania.

## budowa/klient-poprzedni/src/moduly/studio/nastawy-strony.ts

Kształt nastaw strony nie jest wymyślony: to `StudioPageSetup` z kontraktu,
ten sam, który jedzie w profilu wydania i w renderze podglądu. Okno pracy
pokazuje więc kartkę o tych samych nastawach, którymi rdzeń wyda dokument.
Rozmiar nośnika w milimetrach stoi w tym pliku, bo kontrakt niesie samo
oznaczenie, a nie jego wymiary; wykaz jest zamknięty i nazwany, a nośnik
spoza wykazu bierze wymiary A4 i okno mówi o tym wprost, zamiast rysować
kartkę o zgadniętym rozmiarze. Plik nie zna DOM: wejściem są nastawy,
wyjściem liczby i zdania.

Źródłem prawdy o wymiarach nośników jest rdzeń: komenda właściwa temu
obszarowi oddaje szereg A, szereg B, Letter, Legal, Tabloid oraz koperty
DL, C4, C5 i C6 wraz z wymiarami, czytając je z jednego wspólnego wykazu
rdzenia. Okno pracy woła tę komendę przy wczytaniu dokumentu i wchłania
odpowiedź, więc kartka rysuje się wymiarami rdzenia, a nie wymiarami
wbudowanymi. Wykaz wbudowany stoi w tym pliku wyłącznie dlatego, że kartkę
trzeba narysować, zanim odpowiedź rdzenia dojedzie — rysowanie w rozmiarze
zgadniętym byłoby nieprawdą o nośniku. Poprawki wymiarów wnosi się do
wykazu rdzenia; wpis wbudowany ma się do niego równać, a nie odwrotnie.
Sprawdzian pokrywający ten plik pilnuje, żeby wchłonięcie naprawdę
nadpisywało wymiary wbudowane.

Rodzaj nośnika przy wchłanianiu wykazu rdzenia bierze się z pola `kind`
odpowiedzi, gdy ta go niesie, a nie z wpisu wbudowanego — inaczej koperta
dołożona kiedyś do wykazu rdzenia wchodziłaby do okna jako arkusz i nadruk
koperty nie miałby się na czym wykonać.

Margines na oprawę nie ma pola w kontrakcie — `StudioPageSetup` niesie
cztery marginesy i nic poza nimi — więc oprawa żyje przez sesję okna
i jest zgłoszona jako brak pozycji kontraktu; rysowana jest prawdziwie,
bo pas oprawy zabiera pole pisania. Marginesy odbicia włączone znaczą, że
margines lewy jest marginesem wewnętrznym: na stronie nieparzystej stoi po
lewej, na parzystej po prawej — bez tego oprawa wypadałaby raz w rowku, raz
na krawędzi.

Wysokość bloku przy rozdzielaniu stron podaje wołający, bo zmierzyć ją
potrafi wyłącznie przeglądarka, a rachunek podziału ma dać się sprawdzić
bez niej. Blok wyższy od całej strony zostaje na stronie własnej, ponieważ
dzielenie go w środku wymagałoby łamania wiersza wewnątrz akapitu, czego ta
warstwa nie umie i czego nie udaje.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-redakcji.ts

Katalog jest osobny od katalogu czynności warsztatu, choć typy dzieli. Powód
leży w materiale, na którym te dwa zbiory pracują: warsztat bierze dokument
PDF z magazynu okna, redakcja bierze dokument Studia wczytany w edytorze.
Jeden wspólny wykaz kazałby Operatorowi wybierać czynność, która nie ma na
czym pracować, i dowiadywać się o tym dopiero z odmowy.

Nazwy pól są nazwami kontraktu i jadą do rdzenia bez zmiany; etykiety są
zdaniem Operatora i z nazwami się nie pokrywają.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-kontroli-studio.ts

Źródło pracy niesie zmiany śledzone, profile wydania, podgląd, szablony
i komentarze — to czynności redakcji. To źródło niesie czynności kontroli nad
tym, co się z dokumentem stało i czego modelowi nie wolno: odwracalny dziennik
czynności z cofnięciem pojedynczym i ponowieniem; wszystko, co zrobił model,
wraz z licznikiem, skakaniem i cofnięciem z zachowaniem pracy Operatora;
różnicę postaci dwóch wersji i przeniesienie pojedynczego fragmentu do stanu
bieżącego; zapis samoczynny osobnym szeregiem, kopie zapasowe i powrót do
wersji założycielskiej; znakowanie fragmentu w rdzeniu wraz z rodzajami
znaczników własnych Operatora; blokadę fragmentu obowiązującą w rdzeniu;
zajęcie fragmentu przez wykonawcę wraz z nastawami pracy kilku wykonawców
naraz; oraz schowek dokumentu, który sprawdza blokady i odkłada wpis
dziennika, więc wklejenie da się cofnąć pojedynczo.

Każde wywołanie idzie drogą odporną na odmowę: komenda bez uchwytu w rdzeniu
wraca kopertą bez pola stanu, której korelacja nie rozstrzyga, więc zwykłe
wywołanie zostawiłoby okno w ładowaniu bez końca. Droga odporna zamienia taki
brak na zwykłe niepowodzenie nazywające komendę, żeby okno mogło pokazać
Operatorowi, po czyjej stronie jest brak.

Źródło nie ma stanu i nie buduje elementu — sprawdza kształt odpowiedzi
i oddaje ją oknu. Nazwy metod noszą wspólny przedrostek, więc suma
z pozostałymi źródłami modułu nie ma kolizji nazw.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-techniczny.ts

Obszar zakresu technicznego niesie osiem punktów izolacji technicznej: katalog roboczy sesji,
środowisko procesu, katalog danych modelu, dostęp sieciowy, odczyt i zapis plików, konto i token,
model procesu oraz serwer wykonania. Wartości kluczy są zawsze domyślnie wyłączone. Różnica
wobec obszaru kontekstu jest zasadnicza: tam wartość domyślna opisuje zakres widoczności i obie
wartości coś znaczą, tu wyłączony nie jest słabszym wariantem ochrony. Egzekutor rdzenia bierze
wyłącznie klucze włączone i odrzuca wykonanie, które by je naruszyło; klucz wyłączony nie jest
sprawdzany w ogóle, więc nikt nie pilnuje tego zakresu, dopóki nie zostanie świadomie włączony.
Stan wyjściowy platformy to pełna swoboda operacyjna, a widok nazywa ten skutek przy każdym
z ośmiu kluczy osobno, nie tylko raz na górze, żeby przewijanie wykazu nie zgubiło ostrzeżenia.

Odczyt woła komendę odczytu zakresu technicznego, a przełączenie klucza woła komendę zapisu
z pełnym zestawem ośmiu przełączników naraz, bo kontrakt niesie cały zestaw, nie pojedynczy
klucz. Wartość na ekranie zmienia się wyłącznie po potwierdzeniu z rdzenia: stan modułu to
zawsze ostatnia odpowiedź odczytu albo zapisu, nigdy przewidywanie kliknięcia. Odmowa zapisu nie
idzie do paska ogólnego obszaru, tylko do komórki zmiany wiersza klucza, którego dotyczy, aby
było wiadomo, który z ośmiu kluczy się nie zapisał; wartość wyświetlona zostaje ta sprzed próby,
bo rdzeń jej nie potwierdził.

Plakietka wartości ma trzy stany i trzy wyglądy. Stan nieznany, gdy rdzeń nie zwrócił klucza,
i stan wyłączony, gdy rdzeń zwrócił, że nikt nie pilnuje, prowadzą do przeciwnych wniosków
o bezpieczeństwie maszyny, więc nie mogą wyglądać tak samo. Warianty plakietki są zdefiniowane
w bibliotece komponentów: sukces, ostrzeżenie, błąd, informacja, sygnał, rola, atrament — poza
tym zbiorem plakietka nie dostaje żadnego wyróżnienia.

## budowa/klient-poprzedni/src/moduly/studio/okno-session-repository.ts

Historia wersji stała dotąd wąską kolumną na dole modułu i zabierała miejsce
także wtedy, gdy nikt do niej nie zaglądał. Element wyzwalacza jest teraz
wąskim paskiem ze znacznikiem liczby wersji i przyciskiem; cała tabela
otwiera się nakładką nad treścią i schodzi naciśnięciem albo klawiszem
Escape — powierzchnia dokumentu wraca do dokumentu, gdy historia nie jest
używana.

Repozytorium narasta i niczego nie usuwa samo: przywrócenie wersji
przestawia treść bez usuwania wersji nowszych, więc samo cofnięcie jest
odwracalne. „Usuń z wykazu" jest czynnością miejscową i jawną — schowaniem
pozycji w tym oknie, nie usunięciem wersji w rdzeniu; okno mówi to wprost,
żeby Operator nie sądził, że stracił wersję.

Zapisy samoczynne są domyślnie ukryte, z przełącznikiem pokazującym też je.
Gdyby wchodziły do jednego wykazu z wersjami nazwanymi, historia zasypałaby
się w kilka minut. Rozróżnienie bierze się z pola oznaczającego kamień
milowy i z etykiety wersji — drugiego pojęcia okno nie zakłada.

Czynności rozgałęzienia debaty to okno dziś nie pokazuje: rodzina komend
rozgałęzienia jest zbudowana i pracuje dalej, po prostu nie dostaje w tej
turze własnego miejsca w interfejsie.

Zaplecze czynności obejmuje etykietowanie wersji, eksport historii i paczkę
przekazania. Bez zaplecza trzy czynności stoją jako brak nazwany wraz
z komendą, która czeka gotowa w rdzeniu, zamiast jako przycisk milczący po
naciśnięciu.

## budowa/klient-poprzedni/src/moduly/studio/strona-postac-dokumentu.test.ts

Miarą odbioru tego odcinka jest numeracja stron konfigurowana przez Operatora
— styl, umiejscowienie i numer początkowy — oraz marginesy zmienialne.
Sprawdziany mierzą więc treść żądania, które wychodzi z okna do rdzenia, a nie
samo to, że przycisk dał się nacisnąć: nastawa, która nie dojedzie do rdzenia,
nie przeżyje zapisu i jest usterką, nie uproszczeniem.

Druga rzecz mierzona wprost: pole puste znaczy „nie ruszaj tej cechy", a nie
zero. Bez tego każde naciśnięcie przycisku zerowałoby nastawy, których Operator
nie tknął.

## budowa/klient-poprzedni/src/moduly/studio/okno-ingest-ocr-panel.ts

Panel stał na jednym wywołaniu wydobycia tekstu: jeden język, bez silnika,
bez progu pewności, z kolejką prowadzoną w oknie i ginącą z jego
odświeżeniem. Prowadzi teraz rodzinę komend cyfryzacji: kolejkę po stronie
rdzenia, rozpoznanie z pełnym sterowaniem, poprawkę rozpoznanego słowa przed
przyjęciem, przyjęcie wyniku jako dokumentu wraz z pierwszą wersją i wykaz
urządzeń wejściowych.

Poprawianie słów wymaga pracy na obrazie obok tekstu, nie ciasnego paska:
panel stawia obok siebie wykaz słów wraz z ich położeniem na stronie
i pewnością rozpoznania oraz warstwę tekstową pozycji. Obrazu skanu tu nie
ma i panel mówi to wprost — komendy pobierającej bajty zasobu do
przeglądarki kontrakt nie niesie, więc położenie słowa jest podane liczbami,
a nie zaznaczone na obrazku.

Narzędziownia nie jest stałą kolumną: wchodzi przyciskiem jako nakładka
i schodzi. Kolejka z wieloma pozycjami potrzebuje miejsca na wykaz, więc
dostaje nakładkę, a nie pasek, który zabierałby szerokość także wtedy, gdy
nikt nic nie cyfryzuje.

## budowa/klient/src/wejscie/montaz.ts

Okno jest siatką o trzech wierszach: belka, korpus, pas działań. W wierszu
pasa stoją dwie rzeczy — pas przez całą szerokość i nota w kolumnie
tożsamości. Montaż tutaj zostaje wyłącznie przełożeniem stanu na węzły
dokumentu i zdarzeń na wywołania przebiegu, dzięki czemu każda odsłona jest
osiągalna w sprawdzianie bez przeglądarki — sprawdzian prowadzi przebieg, nie
okno.

Okno przygotowania nie ma belki systemowej, bo stoi już wewnątrz ramy
aplikacji, a ta niesie własną; nie ma też odsłon do przełączania.

Wpisanie stałej czasu zwłoki byłoby obietnicą odmierzania, którego nikt nie
mierzył — dlatego licznik dostaje wyłącznie wartość zmierzoną przez przebieg.

Pole, którego dotyczy usterka, niesie atrybut `aria-invalid`, żeby czytnik
ekranu dowiedział się tego samego co oko. Przy kilku usterkach naraz wykaz
niesie same rozpoznania, bez wskazówek co robić — wskazówkę niesie zaznaczone
pole, a jej powtórzenie kilka razy wypchnęłoby formularz poza okno.

Kontrolka niegotowa nie jest wyłączana atrybutem disabled, tylko oznaczona
jako zajęta atrybutem aria-busy, ponieważ przycisk wygaszony nie mówi,
dlaczego nie działa.

Licznik wypisuje czas we wszystkich odsłonach, także ukrytych — inaczej byłby
pusty przez sekundę po pokazaniu odsłony. Ubywa natomiast tylko licznik
widoczny: czas schodzący za plecami doprowadziłby do tego, że użytkownik
zastaje zero, choć odsłonę zobaczył przed chwilą.

## budowa/klient-poprzedni/src/moduly/studio/dziennik-panel.ts

Historia wersji ma swoje okno, a zmiany śledzone swoje miejsce w treści. Ten
panel dokłada to, czego ani jedno, ani drugie nie ma: dziennik czynności
z cofnięciem pojedynczym i nie po kolei wraz z ponowieniem, przy czym
zależności wypisane przy wpisie oznaczają, że odmowa nie jest zaskoczeniem;
zmiany modelu z licznikiem, skakaniem i cofnięciem wszystkiego albo tylko
odhaczonych, z zachowaniem pracy Operatora; różnicę postaci dwóch wersji wraz
z przeniesieniem pojedynczego fragmentu do stanu bieżącego; oraz schowek
dokumentu — odłożenie i wklejenie fragmentu drogą rdzenia, która sprawdza
blokady i odkłada wpis dziennika, więc wklejenie da się cofnąć pojedynczo.

Panel woła rdzeń sam, bo każda z tych czynności oddaje bilans, a bilans jest
treścią dla Operatora, nie wartością pośrednią: co przeszło, co stanęło
i przez którą blokadę. Przepuszczanie go przez okno nadrzędne kosztowałoby
jedno przełożenie na każdej z dziesięciu dróg, a bilans musi trafić na widok
w całości. Skutek czynności — nową treść i postać — panel oddaje oknu
wywołaniem zwrotnym, bo powierzchnia dokumentu jest po jego stronie.

Panel nie zna treści dokumentu ani zaznaczenia: bierze je z kontekstu, który
podaje okno. Dzięki temu ten sam panel obsługuje dokument w zakładce i drugi
w podziale powierzchni.

## budowa/klient-poprzedni/src/uwierzytelnienie/zrodlo-auth.ts

Źródło komend bramki jest jedynym miejscem znającym nazwy komend kontraktu rodziny auth i kształt
ich żądań. Ekran logowania woła cztery z siedmiu komend tej rodziny: wejście, rejestrację,
przedłużenie tokenu — tyle, ile trzeba, żeby wejść — oraz zmianę hasła, bo ta nie wymaga
zalogowania i stoi za odnośnikiem resetu hasła. Dodanie i usunięcie metody należą do czynności
ustawień wykonywanych po zalogowaniu. Potwierdzenie rejestracji nie jest wołane stąd, bo uchwytu
w rdzeniu nie ma; wywołanie wróciłoby kopertą nieznaną i niczym więcej. Odpowiedź niesie typ
koperty, nie tylko wynik, ponieważ rdzeń odpowiada na komendę bez uchwytu kopertą nieznaną
z tym samym kodem, którym wejście odmawia przy bramce nieustawionej — te dwie odmowy znaczą co
innego, więc wysyłka oddaje typ koperty odpowiedzi razem z treścią.

Metoda wejścia widziana przez operatora ogranicza się do hasła i PIN-u, bo tyle rdzeń dziś
naprawdę otwiera; metoda logowania systemu operacyjnego odmawia z powodu pochodzenia dokumentu,
nie z braku kodu, więc nie jest tu metodą — segment wyszarzony byłby bramą, a segment czynny
kłamstwem.

Powitanie połączenia jest podsłuchane z koperty przechodzącej przez kanał, nie wysyłane osobno,
ponieważ powitanie i tak leci przy każdym nawiązaniu połączenia, a subskrypcja staje przy
składaniu źródła, zanim gniazdo się otworzy, więc odpowiedź nie ma jak przepaść. Gdy powitanie
nie wróci w zadanym czasie, wynikiem jest powitanie puste, czyli rdzeń nie wie, i ekran schodzi
na sondę zamiast czekać bez końca; czas oczekiwania zerowy znaczy, że trzeba wziąć to, co już
jest.

Sonda PIN-u idzie tą samą drogą co sonda hasła i z tego samego powodu: kontrakt nie ma komendy
odczytu metod, a segment PIN-u nie ma prawa stanąć na ekranie, jeżeli nie ma czego otworzyć.
Rdzeń szuka metody po parze rodzaj i urządzenie, zanim spojrzy na sekret, więc żądanie bez
sekretu odpowiada wprost: brak wyniku znaczy, że PIN-u tu nie ma, a niepowodzenie walidacji
znaczy, że PIN jest i zabrakło tylko sekretu. Sonda nie podnosi dławika prób, bo zwłokę
zwiększa wyłącznie sekret niezgodny, a tu sekretu nie ma.

Wejście przez bramkę hasłem albo PIN-em niesie do rdzenia także przełącznik przedłużenia sesji,
nie tylko do magazynu przeglądarki, bo bez tego zaznaczenie przedłużałoby wyłącznie życie zapisu
lokalnego, a sesja gasłaby po dobie roboczej.

Rejestracja zakłada jedyne konto właściciela przy pierwszym uruchomieniu i nie zakłada sesji.
Konto powstaje niepotwierdzone, a token dostępu wydaje dopiero potwierdzenie po drodze przysłanej
listem, dlatego odpowiedź rejestracji nie niesie sesji i ekran przechodzi wtedy do kroku
potwierdzenia.

Prośba o drogę odzyskania konta idzie na adres uwierzytelniający, a odpowiedź jest zawsze taka
sama niezależnie od tego, czy adres pasuje do konta — ekran nie ma z czego wywnioskować, czy
konto istnieje, i tak ma być.

Zmiana hasła ze znanym hasłem dotychczasowym należy do ekranu wejścia, choć wygląda na czynność
ustawień: sprawdza wyłącznie zgodność hasła dotychczasowego z zapisem bramki i nie pyta o sesję.
Zmiana jest więc wykonalna przed zalogowaniem i odnośnik resetu hasła ma co wywołać. Odzyskaniem
hasła zapomnianego to nie jest, bo takiej drogi kontrakt nie ma.

Wysyłka jednej komendy odczytuje całą kopertę odpowiedzi. Subskrypcja staje przed wysłaniem,
a dopasowanie idzie po identyfikatorze żądania, dokładnie tak jak koreluje kanał, więc między
subskrypcją a wysyłką nie ma okna, w którym odpowiedź mogłaby przepaść. Transport kolejkuje
ramki do chwili otwarcia połączenia, więc wywołanie przed nawiązaniem łączności czeka, zamiast
przepaść.

Sonda stanu bramki wysyła wejście metodą hasła bez sekretu, bo kontrakt nie ma osobnej komendy
odczytu stanu, a ekran musi wiedzieć, czy pokazać wejście, czy pierwsze ustawienie hasła. Adapter
bramki najpierw szuka kotwicy, a dopiero potem patrzy na sekret, więc odmowa mówi wprost: brak
wyniku znaczy, że hasło bramki nie istnieje i bramki jeszcze nie ustawiono, a niepowodzenie
walidacji znaczy, że kotwica jest i zabrakło tylko sekretu. Sonda jest drogą zapasową, nie
pierwszą: odpowiedź na to samo pytanie niesie powitanie, które nie wyprowadza stanu z odmowy
i nie przechodzi przez dławik prób; sonda zostaje na wypadek milczenia rdzenia, bo pole puste
znaczy, że rdzeń nie wie, a wtedy trzeba zapytać inaczej, nie zgadnąć. Sonda niczego nie zmienia
w rdzeniu i nie niesie żadnego sekretu.

Identyfikator urządzenia idzie tylko przy PIN-ie, bo metoda hasła kontraktu go nie wymaga —
hasło otwiera bramkę z każdej maszyny, a PIN jest właściwy maszynie i bez identyfikatora rdzeń
odmawia wprost. Tożsamość klienta z powitania urządzeniem nie jest, bo nadaje się ją na czas
uruchomienia, więc PIN chodzi po zapisie trwałym urządzenia.

## budowa/klient-poprzedni/src/moduly/studio/czynnosci-warsztatu.ts

Okno buduje pola z tego wykazu i składa żądanie funkcją zloz czynności.
Piętnaście osobnych formularzy dawałoby piętnaście miejsc, w których nazwa
pola kontraktu bywa wpisana z pamięci; w tym pliku każde pole stoi raz.
Nazwy pól są nazwami kontraktu, bo to one jadą do rdzenia. Etykiety są
zdaniem Operatora i z nazwami się nie pokrywają — Operator czyta „zakres
stron", a rdzeń dostaje pole `pages`.

## budowa/klient-poprzedni/src/moduly/studio/przybornik-znakowania.test.ts

Sprawdziany odcinka znakowania, asystenta, schowka i osadzenia mierzą rzeczy, które da się
zmierzyć bez stawiania okna, i tylko te, w których pomyłka jest cicha: rozróżnienie trzech
bytów marginesu, zawężenie wykazu znakowań, kolejność czynności pływaka liczoną z użycia,
zaporę przed pokazaniem tej samej operacji dwa razy w katalogu oraz zapis pochodzenia. Czego
tu nie ma — sprawdzianu, że przycisk wywołuje komendę — jest tak samo istotne: sprawdzian
takiego kształtu mierzyłby atrapę, którą sam stawia, a nie skutek na dokumencie.

## budowa/klient-poprzedni/src/moduly/studio/tabela-panel.ts

Siatka jest wejściem pierwszym, a pola liczbowe stoją obok jako droga druga —
dla tabeli większej niż siatka i dla pracy z klawiatury. Oba wejścia prowadzą
do tej samej komendy wstawienia, więc nie ma dwóch zachowań.

Każda czynność tabeli oddaje bilans i panel go wypisuje zawsze, nie tylko przy
pominięciu. Scalenie komórek w zablokowanym fragmencie wraca odpowiedzią
pomyślną z pominięciem w bilansie — bez wypisania bilansu wyglądałoby to na
scalenie wykonane. Wpis dziennika idzie do tego samego zdania, bo bez niego
Operator nie wie, co ma cofnąć.

Panel nie liczy sam ani szerokości kolumn, ani wyniku sortowania, ani zamiany
tekstu na tabelę. Wszystko to robi rdzeń; panel składa żądanie, czyta tabelę
z odpowiedzi i pokazuje szerokości policzone, a nie założone.

## budowa/klient-poprzedni/src/moduly/studio/wczytanie-dokumentu.ts

Otwarcie dokumentu przyjmuje trzy wzajemnie wykluczające się pola:
identyfikator dokumentu otwartego wcześniej w sesji, identyfikator pliku
repozytorium albo ścieżkę na urządzeniu. Formularz pyta o źródło wprost,
zamiast zgadywać, bo pomyłka kończyłaby się odmową rdzenia o powodzie
trudnym do odczytania. Wybierak źródła nie jest wykazem plików: wykaz
zasobów Library należy do okna Library Explorer, a to pole przyjmuje sam
identyfikator. Wykaz formatów pochodzi z kontraktu, nie z listy zapisanej
w widoku; format wybiera rdzeń przy wczytaniu, okno go tylko pokazuje.

Wniesienie pliku różni się od otwarcia dokumentu tym, że wnosi postać pliku
Operatora — arkusz stylów, sekcje, tabele i obrazy — zamiast otwierać
pozycję, którą rdzeń już prowadzi. Trzy komendy to robią i dlatego stoją
osobno: wniesienie pliku wnosi dokumenty i szablony biurowe wraz
z rozpoznaniem zapisu znaków, wniesienie PDF-u odzyskuje z niego tekst,
akapity, tabele i obrazy na tyle, na ile PDF je niesie, a wniesienie obrazu
wnosi obraz wprost w miejsce kursora.

Odzyskanie z PDF-u jest odtworzeniem, nie odczytem: PDF nie niesie
struktury akapitu ani tabeli wprost. Panel wypisuje więc bilans za każdym
razem — strony z warstwą tekstową i bez niej, tabele rozpoznane
i nierozpoznane, obrazy osadzone i pominięte. PDF ze samych skanów kieruje
na rozpoznanie tekstu, a panel mówi to wprost, zamiast oddać pustą kartkę
jako gotowy dokument.

## budowa/klient-poprzedni/src/moduly/workspace/hub-planowania.ts

Cztery widoki huba planowania dzielą jeden zbiór zadań. Przełącznik widoku nie zmienia danych,
tylko pytanie zadawane rdzeniowi: lista pyta o wykaz zadań, tablica o stan tablicy kanban, oś
czasu o harmonogram, kalendarz o wpisy kalendarza. Zadanie zmienione w jednym widoku jest tym
samym zadaniem w pozostałych trzech, bo drugiego zapisu zadania w module nie ma. Zadanie
wskazane w wykazie staje się przedmiotem czynności paska akcji; wskazanie idzie polem, a nie
stanem ukrytym, żeby było widać, czego dotyczy usunięcie, zanim się je naciśnie.

## budowa/klient-poprzedni/src/punkty-izolacji/obszar-efektywna.ts

Wartość zapisana świadomie i wartość domyślna muszą wyglądać inaczej, ponieważ
stan wyjściowy platformy to pełna swoboda operacyjna — egzekutor odrzuca
naruszenie tylko przy punkcie włączonym, a wartość wyłączona niczego nie
ogranicza. Pokazanie domyślnej tak samo jak zapisu Operatora sugerowałoby
ochronę, której nikt nie włączył.

Pochodzenie czytamy z odpowiedzi rdzenia na isolation.policy.preview: obecne
pole origin znaczy zapis na wskazanym poziomie, brak origin znaczy wartość
z rejestru definicji. Stan „nieznane" wchodzi wtedy, gdy odpowiedź w ogóle nie
niesie przełącznika danego punktu.

Żądanie podglądu musi nieść punkt widzenia — sesję, warstwę, zasięg i okno.
Przy pustym żądaniu rdzeń schodzi kolejką do poziomu globalnego, więc okno
pokazałoby politykę całej platformy pod nazwą „efektywna".

Czwarty stan „brak-odczytu" jest wyłącznie kliencki i nie wolno go zlewać
z wartością domyślną: oba prowadzą do przeciwnych wniosków o bezpieczeństwie.

Jedenaście punktów izolacji odpowiada definicjom rdzenia w module konfiguracji
izolacji. Objaśnienie każdego punktu jest opisem klucza — co znaczy i jaki
jest stan wyjściowy platformy — a nie wynikiem odczytu, dlatego stoi w pliku
klienta: odpowiedź policy.preview niesie wartości, nie objaśnienia.

Miejscem styku granicy asystenta jest komenda zapisu konfiguracji sesji z roli
klawiatury: pisze przedmioty tych samych punktów — katalog roboczy, katalogi
dodatkowe, środowisko, uruchomienie, narzędzia, uprawnienia i konto —
przełączników izolacji jednak nie rusza, więc egzekutor zostaje na miejscu
i granicy nie przekracza.

## budowa/klient-poprzedni/src/moduly/studio/okno-warsztatu-dokumentu.ts

Wstążka jest kontekstowa: wchodzi po naciśnięciu przycisku, a grupa Ochrona
wstążki okna pracy woła przeniesienie ogniska do niej; wchodzi też
samoczynnie, gdy dokument w pracy jest PDF-em. Nie zajmuje miejsca, kiedy
Operator nad PDF-em nie pracuje.

Grupy Strony, Nakładanie, Treść, Bezpieczeństwo i Narzędziownia cyfryzacji
zawężają wykaz czynności do swojego obszaru, a formularz pod nimi jest
jeden: wybór czynności przestawia pola. Piętnaście osobnych formularzy
dałoby wstążkę trudną do przejrzenia i do utrzymania. Pola pochodzą
z katalogu czynności warsztatu, więc nazwa pola kontraktu stoi w drzewie
raz.

Materiał wchodzi z magazynu okna i wynik do niego wraca. Żadna czynność nie
zmienia materiału w miejscu — okno mówi to przy polu materiału, bo Operator
ma wiedzieć, że pomyłka nie kosztuje go dokumentu źródłowego. Odpowiedź
jest opisana skutkiem, nie słowem gotowe: liczbą stron, liczbą części,
identyfikatorem nowego zasobu.

Kolejność grup wstążki jest kolejnością pracy nad plikiem: najpierw strony,
potem to, co się na nich kładzie, potem treść, na końcu bezpieczeństwo
przed wysyłką. Narzędziownia cyfryzacji stoi grupą osobną, bo pracuje przed
dokumentem — na materiale, z którego dokument powstaje.

## budowa/klient-poprzedni/src/moduly/wiedza/okno-wyszukiwania-znaczenia.ts

Wyszukiwanie po słowach ma osobną drogę; to okno prowadzi wyszukiwanie po znaczeniu oraz
przebudowę wskaźnika. Źródło jest treścią wyniku, nie ozdobą przy nim, bo trafienie bez
wskazania źródła jest bezwartościowe: każdy fragment niesie źródło w nagłówku pozycji, nie
w dymku, źródło jest klikalne przyciskiem kopiującym wskazanie do schowka, bo to jedyne, czym
da się po źródło sięgnąć w innym oknie, a fragment bez wskazania dostaje zdanie o braku
wskazania zamiast martwego przycisku. Trafność też jest widoczna jako plakietka w setnych przy
każdej pozycji, bo bez niej wykaz posortowany malejąco wygląda tak samo przy trafieniu bliskim
i przypadkowym. Przebudowa wskaźnika stoi w tym samym oknie, bo pusty wynik szukania ma
dokładnie dwie przyczyny: nie ma czego znaleźć albo wskaźnik nie został zbudowany, a przycisk
przebudowy obok pustego wyniku pozwala rozstrzygnąć to jednym kliknięciem.

Wiersz ze sterem zamiast pola różni się od wiersza biblioteki kontrolek tym, że tamten buduje
etykietę wiążącą się z pierwszym potomkiem dającym się etykietować, czyli z uchwytem menu —
kliknięcie w podpis otwierałoby wtedy menu, którym podpis nie jest, a nazwa dostępna uchwytu
wchodziłaby w spór z opisem, który mechanizm ustawia sam i który niesie bieżącą wartość nastawy.

Zdanie opisujące każdy zakres wiedzy jest tu potrzebne bardziej niż przy większości nastaw:
wszystko i pliki przestrzeni brzmią podobnie, a przeszukują dwa różne zbiory, i to od nich
zależy, czy pusty wynik znaczy, że nie ma czego znaleźć, czy że szuka się nie tam.

Uchwyt do źródła kopiuje wskazanie zamiast skakać do niego, bo kontrakt nie ma komendy
otwierającej źródło fragmentu — samo wskazanie jest jedynym sposobem sięgnięcia po nie w oknie
biblioteki albo w rozmowie z modelem.

Trafność bez miary widocznej sprawiałaby, że wykaz posortowany malejąco wyglądałby tak samo
przy trafieniu bliskim i przy przypadkowym, a to rozstrzygnięcie należy do czytającego, nie do
samego wykazu; rdzeń ma prawo trafności nie podać, a wtedy plakietka mówi o braku miary, zamiast
zmyślać zero.

## budowa/klient-poprzedni/src/moduly/studio/konwersja-dokumentu.ts

Wydania żąda i pasek narzędzi Studio Editora, i Preview Window; czynność jest
jedna i stoi w jednym miejscu, bo dwie kopie rozjechałyby się przy pierwszej
poprawce zdania o wyniku.

Przedmiotem zamiany jest treść zaakceptowana modułu, nie plik na dysku.
Dokument Studia mieszka w rdzeniu, a nie w systemie plików Operatora, więc
komenda zamiany formatu dostaje treść wprost polem. Wynik jest zasobem
magazynu rdzenia — komenda oddaje jego identyfikator i rozmiar, a nie bajty do
pobrania.

Format dokumentu ma cztery wartości, a zamiana formatu przyjmuje nazwy
własnego słownika. Trzy z czterech mają w nim odpowiednik wprost; PDF
odpowiednika użytecznego nie ma, bo treść jedzie napisem, a nie plikiem —
napis nie jest PDF-em, choćby rdzeń tak nazywał dokument, z którego powstał.
Zwraca się wtedy tekst czysty i mówi się o tym Operatorowi osobnym zdaniem.
Bez tego zdania Operator widziałby format dokumentu z rdzenia i zamawiałby
wydanie z formatu, którego rdzeń w tej drodze nie czyta.

Komenda zamiany treści napisem stoi obok czterech komend wydania dokumentu po
stronie rdzenia, nie zamiast nich: pierwsza nie ma czym przenieść postaci
dokumentu prowadzonego jako PDF albo DOCX, bo treścią jest tam sam tekst.
Cztery komendy niżej pracują na dokumencie po stronie rdzenia i przenoszą
arkusz stylów, sekcje, tabele, obrazy i aparat — zapisują pod nową nazwą,
wydają do formatu, wydają wiele dokumentów naraz albo przestawiają format
dokumentu Studia. Starsza droga zostaje dla treści, która nie jest dokumentem
rdzenia.

Format uboższy niż dokument jest normalną sytuacją; przemilczenie straty nie
jest. Wydanie do tekstu czystego, które zgubiło tabelę i przypisy, mówi to
wprost, a wydanie wsadowe wypisuje bilans dokument po dokumencie — odmowa
jednego nie ukrywa się za liczbą wydanych.

## budowa/klient-poprzedni/src/moduly/roundtable/okno-model-panels.ts

Skład uczestników idzie ze stanu debaty, nie z własnej listy okna: stan subskrybuje zmianę
debaty i dokłada uczestników z przyrostu tury oraz z czynności moderatora, a to okno tylko
czyta i rysuje. Odpowiedzi rosną na żywo z fragmentów strumienia wypowiedzi, a rysowanie
składu wraz z odpowiedziami leży w module składu modeli; subskrypcji na zmianę stanu okno nie
zakłada samo, bo jedna na całe złożenie modułu wystarcza.

Dwie tożsamości jednego kanału dają dwa panele. Kluczem panelu jest identyfikator uczestnika,
nigdy identyfikator kanału, bo tak stanowi kontrakt dodania modelu. Liczba dotychczasowych
wystąpień kanału mówi, ile razy kanał już wystąpił — od drugiego wystąpienia wzwyż panel niesie
oznaczenie powtórzenia, żeby dwa identyczne z pozoru panele nie wyglądały na usterkę.

Świeżo otwarte okno nie zna jeszcze składu debaty. Odczyt składu i odczyt całej debaty są
w kontrakcie, ale to okno żadnego z nich jeszcze nie wywołuje, bo obsługi nie zbudowano; do
tego czasu skład narasta z przyrostu zdarzenia i z uczestników zwracanych przez polecenie
moderatora wprost. Stan pusty na starcie jest stanem poprawnym, opisanym wprost, nie udawaną
pustką.

Kontrakt nie niesie górnego limitu instancji tego okna i okno nie narzuca własnego — skład
rośnie tak, jak rośnie w rdzeniu. Ile uczestników mieści debata, mówi nota nad formularzem,
a nie bramka odmowy: uczestnik ponad zwyczajowy skład przechodzi przez rdzeń bez odmowy, więc
przycisk dodania uczestnika nie zna żadnego progu. Nota podaje, ilu uczestników debata ma, czego
nie liczy ani rdzeń, ani kontrakt, i gdzie leży sufit gniazd sceny okien równoległych.

Dwie pozycje paska akcji mają komendę w kontrakcie i nie mają jeszcze obsługi: usunięcie
uczestnika oraz oznaczenie uczestnika jako kluczowego. Wyciszenia w pasku nie ma wcale, bo jest
czynnością moderatora i należy do panelu moderatora, nie do tego okna. Pozostałe czynności
uczestnika bez obsługi siedzą w zestawie akcji warstwy trzeciej: siedem pozycji ma dziś komendę
w kontrakcie, a ósma — liczba wariantów jednej persony — ma pole w uczestniku i nie ma komendy,
która by je zapisała.

## budowa/klient-poprzedni/src/moduly/studio/zapis-formatowany.ts

Okno pracy z dokumentem pokazuje pismo, a nie jego zapis: nagłówek jest
większy, cytat wcięty, lista wypunktowana, tabela ma siatkę. Zapisem pod
tym wyglądem zostaje markdown, bo to jego niesie pole treści dokumentu —
kontrakt nie ma ani pola stylów, ani pola układu, więc drugi zapis, własny
i bogatszy, nie miałby gdzie dojechać do rdzenia i ginąłby przy pierwszym
zapisie dokumentu. Plik nie zna okna ani stanu modułu: wejściem jest napis,
wyjściem bloki albo element, i odwrotnie — dzięki temu przekształcenia
sprawdza się bez stawiania powierzchni.

Trzy drogi i wszystkie trzy muszą się zgadzać: czytanie napisu na bloki
przy czytaniu dokumentu, rysowanie bloku jako element widoczny przy
pokazaniu dokumentu i odczyt elementu z powrotem na napis przy pisaniu
w oknie. Droga trzecia jest odwrotnością drugiej i pierwszej razem: to, co
Operator wpisze w widoku formatowanym, wraca do markdown tą samą składnią,
którą widok odczytał.

Rodzaj bloku przy jego wyrysowaniu idzie do atrybutu danych elementu, bo
z niego czyta go zarówno arkusz stylów, jak i odczyt powierzchni z powrotem
na napis. Wyprowadzanie rodzaju z samej nazwy elementu działałoby tylko
dopóty, dopóki przeglądarka nie wstawi własnego elementu przy naciśnięciu
klawisza Enter.

Podkreślenia markdown nie ma, więc jedzie znacznikiem HTML — markdown
przepuszcza go nietkniętym, a odczyt powierzchni odkłada go z powrotem tak
samo. Odczyt powierzchni edycji jest odwrotnością wyrysu i musi znieść to,
co dokłada przeglądarka samodzielnie: naciśnięcie klawisza Enter zakłada
nowy blok, wklejenie treści dokłada element ze stylami, a pogrubienie
z klawiatury bywa innym elementem niż ten, którym pisze wyrys.

Długość zapisu markdown do wskazanego punktu jest potrzebna do przełożenia
zaznaczenia w widoku formatowanym na zakres znaków w treści dokumentu —
zakres, który jedzie do rdzenia jako granice zaznaczenia. Liczenie go po
samym tekście widocznym dałoby wartość mniejszą od prawdziwej o długość
znaczników składni.

## budowa/klient-poprzedni/src/moduly/studio/linijka-pozioma.ts

Linijka chwyta margines lewy i prawy, wcięcie pierwszego wiersza, wcięcie
lewe i prawe osobnymi znacznikami, tabulatory zakładane naciśnięciem wraz
z rodzajem i znakiem wiodącym oraz szerokości kolumn tabeli. Każdy chwyt
przyjmuje wartość także z klawiatury, w milimetrach albo calach, bo
przestawianie marginesu pisma urzędowego wyłącznie myszą odcięłoby część
Operatorów.

Wcięcie pierwszego wiersza liczy się względem wcięcia lewego, bo tak zachowuje
się linijka pakietu biurowego: przeciągnięcie wcięcia lewego zabiera pierwszy
wiersz ze sobą. Gdyby oba wcięcia były liczone od krawędzi pola, Operator przy
każdej zmianie wcięcia lewego poprawiałby drugie osobno.

Linijka nie zmienia dokumentu i nie woła rdzenia: oddaje nastawę temu, kto ją
zbudował, a nastawy strony i wcięcia trzyma powierzchnia. Inaczej powstałyby
dwa źródła prawdy o marginesie, które rozjechałyby się przy pierwszym profilu
wydania wziętym z rdzenia.

Wybieracz rodzaju tabulatora pokazuje rodzaj na przycisku, żeby Operator
wiedział, co założy, zanim naciśnie podziałkę; jedno naciśnięcie przestawia
rodzaj następnego zakładanego tabulatora.

Krawędź kolumny tabeli nie może przejść przez sąsiednią, bo kolumna
o szerokości ujemnej nie jest kolumną. Granicą jest sąsiad odsunięty
o najmniejsze pole, a nie sam sąsiad — dwie krawędzie w jednym miejscu dałyby
kolumnę zerową.

## budowa/klient/src/wejscie/przebieg.ts

Przebieg nie zna żadnego składnika okna: rozstrzyga wyłącznie, który etap
i która odsłona obowiązuje oraz co wysłać do rdzenia. Dzięki temu każda
odsłona, także odsłona błędu i wstrzymania, jest osiągalna w sprawdzianie bez
przeglądarki. Nazwy komend i kształty ich treści pochodzą wyłącznie
z kontraktu; przebieg nie powtarza ani jednego literału nazwy komendy.

Odsłona `konto-bez-potwierdzenia` obsługuje rejestrację bez konta nadawczego:
w tym wariancie wejście działa hasłem, a okno musi nazwać adres, którego
nikt nie potwierdził, bo adres jest jedyną drogą odzyskania konta. Z kontem
nadawczym list poszedł i okno prowadzi do jego przepisania zamiast do tej
odsłony.

Progu prób logowania nie ma. Kontrakt stanowi wprost, że po stronie rdzenia
nie ma i nie będzie odmowy „za dużo prób" — nieudana próba nakłada na
następną rosnącą zwłokę, nic więcej. Pomiar to potwierdza: osiem kolejnych
nieudanych prób oddaje osiem razy odmowę uwierzytelnienia, a po nich hasło
poprawne wpuszcza. Odsłona z prototypu zostaje i obsługuje zwłokę, nie
zaporę. Ten sam pomiar zwłoki — czas trwania ostatniego wywołania, bo rdzeń
nakłada zwłokę przed odpowiedzią — nazywa też zwłokę żądań odzyskania konta:
progu wysyłek też nie ma, zmierzone siedmioma kolejnymi wysłaniami przyjętymi
bez odmowy.

Miejsce przechowania tokenu bramki nie jest rozstrzygnięte przez ten teren —
należy do powłoki wołającej. Magazyn działający wyłącznie w pamięci procesu
jest więc wartością domyślną: token ginie wraz z procesem, co jest
zachowaniem uczciwym, bo magazyn udający trwałość obiecywałby rozpoznanie
urządzenia, którego nie ma.

Liczba etapów przygotowania środowiska równa się liczbie komend, jakie droga
wejścia potrafi zmierzyć — uwierzytelnienie i przywrócenie kart sesji. Etap
bez komendy w kontrakcie nie jest etapem czekającym, tylko obietnicą, której
nikt nie wykona, a postęp liczony razem z nim nie dobiegłby końca nigdy.

## budowa/klient-poprzedni/src/moduly/terminal/zrodlo-terminala.ts

Każda czynność oddaje typ Wynik, nie samą treść. W terminalu odmowa jest
zjawiskiem zwykłym, nie wyjątkiem: tryb uprawnień okna albo punkt izolacji
zatrzymuje polecenie kodem odmowy uprawnień, a źródło nie zamienia takiej
odmowy w pustą listę.

Wyjście procesu przychodzi zdarzeniem strumienia fragmentów, w którym
identyfikator komunikatu niesie identyfikator procesu z rejestru rdzenia.
Komenda odczytu strumienia zostaje niewpięta, choć rdzeń ją obsługuje: ogon
historii, który ta komenda oddaje, to co do wiersza ta sama treść, która
płynie do wspólnego bufora okna przez zdarzenie strumienia, a w buforze
nie ma klucza, po którym dałoby się odsiać wiersz przyjęty drugi raz.
Komenda odczytu wyjścia jednego procesu jest wpięta, bo zawęża się do
jednego procesu, a jej odpowiedź trafia do własnego widoku pozycji, którego
bufor nie widzi wcale.

Odczyt sięga dalej niż strumień: bufor żyje jedno połączenie, a rejestr
rdzenia jeden bieg rdzenia. Po ponownym podłączeniu gniazda bufor okna jest
pusty, a wykaz procesów oraz odczyt wyjścia wciąż oddają proces wraz z jego
pełnym wyjściem.

Zaplanowanie zadania powłoki jest zapisem automatyki o jednym kroku
rodzaju polecenia wołającym wykonanie w terminalu. Druga rodzina komend po
tej stronie dałaby dwie prawdy o jednym harmonogramie — i to jest jedyny
powód, dla którego okno składa plan z komend cudzego modułu, a nie ze
swoich.

## budowa/klient-poprzedni/src/okna-rownolegle/uklad-okien.ts

Rozsyłanie odczytu łączności idzie z jednego miejsca układu, a nie z każdego gniazda osobno,
bo łącze jest jedno na całego klienta: gniazdo pierwsze i czwarte muszą po jednej zmianie
transportu pokazać dokładnie to samo. Wejście zostaje publiczne, bo scena podglądu układu
pracuje bez rdzenia i podaje odczyty ręcznie.

Kod okna wykonania powstaje dopiero po uzgodnieniu z rdzeniem, już po zmontowaniu układu, bo
panele pomocnicze wołają komendy żądające tego kodu. Dopóki kodu nie ma, panel nazywa wprost,
czego mu brakuje, zamiast pokazywać pustkę.

Subskrypcja strumienia przeżywa usunięcie węzła z drzewa dokumentu, więc zamknięcie panelu musi
odpiąć subskrypcje rdzenia wprost — scena zdjęta bez tego wywołania zostawiałaby żywe nasłuchy.

Pole opcji o wbudowanej rozmowie jest domyślnie wyłączone, bo scena sesji osadza właściwy widok
z zewnątrz; stanowisko podglądu układu włącza je jawnie, bo nie ma z zewnątrz niczego do
osadzenia. Kanał do rdzenia bywa pominięty — wtedy menu paneli jest puste i mówi to wprost,
bo stanowisko podglądu układu pracuje bez rdzenia i właśnie tak ma wyglądać. Transport bywa
pominięty analogicznie — wtedy plakietki łączności milczą, bo scena nie zna łącza.

Moduł bez rozmowy nie zostawia sceny pustej: scena trzyma minimalną liczbę okien, a zdanie
przełącznika mówi wprost, że rozmowy w tym module nie ma, żeby pusta scena nie czytała się jak
awaria. Pytanie, czy dostawienie gniazda cokolwiek zmieni, idzie tą samą rachubą, którą potem
wykona ustawienie liczby, bo sufit składa się z figury modułu i z górnej granicy, a figura
przestawia się z rdzenia — osobne porównanie rozjechałoby się z tamtym przycięciem przy
pierwszym module węższym niż scena. Odpowiedź przecząca znaczy krótszą listę w menu, a nie
wiersz wygaszony.

Znacznik przestawiania ról przez sam układ rozróżnia rolę z zewnątrz, która jest decyzją
i zostaje, od roli nadanej przez zmianę liczby okien, która jest domyślna i ustępuje następnej;
nadawanie ról jest synchroniczne, więc znacznik nie przecieka poza swoją pętlę.

Pozycja otwarcia nowego okna idzie tą samą drogą co przełącznik liczby: podniesienie liczby
wprowadza gniazdo na scenę, a wejście gniazda na scenę zamawia dla niego okno rdzenia. Okno
założone z boku byłoby niewidoczne, więc tu stoją same wywołania — układ nie wie, że po drugiej
stronie jest wiersz menu.

Moduł sceny idzie z gniazda pierwszego, bo to ono niesie okno uzgodnione z rdzeniem i na nie
przychodzi zmiana modułu, także po wejściu do przestrzeni, które przestawia moduł okna zamiast
zakładać drugie. Scena z oknami w różnych modułach bierze figurę z gniazda pierwszego.

Zdanie sceny o liczbie okien modułu składa do trzech zdań, każde tylko wtedy, gdy jest
prawdziwe: czym jest figura modułu, że pozycja wyższa okna nie dokłada, oraz że scena trzyma
okien więcej, niż moduł prowadzi. Ostatnie bierze się stąd, że układ okien nie zdejmuje ich
sam — zdjęcie okna zamyka je również w rdzeniu i jest decyzją operatora.

Żądanie liczby ponad figurę modułu nie jest odmową: pozycja przełącznika zostaje czynna, scena
bierze tyle okien, ile moduł prowadzi, a zdanie pod grupą mówi, dlaczego wyżej się nie da. Cicha
zgoda na okno ponad figurę modułu byłaby sukcesem udawanym.

Przestawienie sceny na moduł nie zdejmuje okien, nawet gdy jest ich więcej, niż moduł prowadzi,
bo zdjęcie okna ze sceny zamyka je również w rdzeniu, a zrobione samoczynnie przy przestawieniu
modułu byłoby zamknięciem okna, o które nikt nie prosił. Scena mówi o nadmiarze wprost
i zostawia zdjęcie operatorowi.

## budowa/klient-poprzedni/src/moduly/studio/okno-petli-wykonawczej.ts

Powierzchnia należy do dokumentu: panele nie zajmują stałych kolumn, wchodzą
na żądanie i schodzą, gdy nie są używane. Okno pętli wchodzi więc znacznikiem
przebiegu i schodzi po zamknięciu, a stała kolumna jest trybem do wyboru,
jawnym, odwracalnym i pamiętanym — nigdy postacią domyślną. Samoczynne
otwarcie zostaje przy zleceniu wielozadaniowym, bo okno wchodzi wtedy, gdy
naprawdę jest używane, i to jest dokładnie ta sama reguła zastosowana wprost.

Kolejka zadań, obsada wykonawców, warsztat łańcucha i tryb wsadowy działają
na tym samym stanie. Krok łańcucha i dokument wsadu są zadaniami tej samej
kolejki, bo kontrakt mówi wprost: przebieg łańcucha i wsadu prowadzi pętla
wykonawcza okna. Druga maszyneria obok byłaby drugą prawdą o tym samym
przebiegu.

Okno nie przestawia nastaw pętli. Nastawy — czy pętla jest czynna, ilu
wykonawców naraz, co przy spięciu — idą osobną rodziną komend. Okno je czyta
i pokazuje wprost, a gdy pętla jest wyłączona, mówi to zdaniem, zamiast
pozwalać naciskać przycisk, który i tak odmówi.

## budowa/klient-poprzedni/src/moduly/studio/zrodlo-postaci-studio.ts

Nastawy strony, marginesy i numeracja stały dotąd wyłącznie w kliencie, prowadzone
przez sesję okna, a jedyną drogą trwałości był profil wydania — osiem pól bez
oprawy, bez sekcji i bez formatu numeru. Rdzeń niesie dziś rodziny komend
obejmujące stronę, sekcję, styl, formatowanie, listę, znak specjalny, widok,
tabulator, tekst, postać dokumentu i pochodzenie fragmentów — i to one są
jedyną drogą, którą nastawa Operatora przeżywa zapis.

Metody biorą treść żądania z kontraktu w całości, a nie rozłożone argumenty.
Rodzina postaci dokumentu ma pola liczone dziesiątkami — same nastawy strony
to czternaście — i połowa z nich jest opcjonalna w znaczeniu „nie ruszaj tej
cechy". Rozkładanie ich na argumenty zmuszałoby wołacza do podawania wartości
tam, gdzie chce ciszy, a „podano zero" i „nie podano" znaczą tu różne rzeczy.

Wywołanie idzie przez osłonę, bo komenda bez uchwytu w rdzeniu wraca kopertą
ogólną bez pola stanu, którego korelacja klienta nie rozstrzyga — okno stałoby
w ładowaniu bez końca. Osłona zamienia to w zwykłą odmowę nazywającą komendę;
rodzina komend jest świeża, więc rdzeń starszy od klienta jest tu przypadkiem
realnym, nie teoretycznym.

## budowa/klient-poprzedni/src/moduly/studio/aparat-panel.ts

Spis treści, spisy ilustracji i tabel, przypisy dolne i końcowe, podpisy,
zakładki, odwołania wzajemne, odsyłacze, powołania, bibliografia, hasła
indeksu i indeks — wszystkie są elementami wyliczanymi z dokumentu. Pola
dokumentu (numer strony, liczba stron, data, właściwość, pole obliczane) mają
tę samą naturę. Wspólne jest to, że po zmianie treści stają się nieświeże,
a odświeżenie liczy je od nowa.

Panel prowadzi jeden wykaz „do odświeżenia" dla obu rodzin i pisze liczbą, ile
elementów rozjechało się z dokumentem. Spis treści pokazany bez znaku
nieświeżości kłamałby o dokumencie, którego nagłówki się zmieniły.

Numer przypisu i powołania nadaje rdzeń przy odświeżeniu, a panel go
wyłącznie pokazuje: okno nie wymyśla żadnego kodu i nie numeruje niczego samo.

## budowa/klient-poprzedni/src/moduly/terminal/katalog-funkcji.ts

Moduł Terminal obiecuje zastąpić emulator terminala, multiplekser sesji,
klienta SSH i SFTP, harmonogram zadań, uruchamianie skryptów, menedżera
sekretów, konsolę szeregową i klienta kontenerów. Część tych zdolności ma
pokrycie w kontrakcie, część nie, a część wymaga programu, którego instalka
Danaco Console nie niesie. Bez tego wykazu Operator poznaje różnicę dopiero
w chwili, w której czegoś potrzebuje — a wtedy milczenie interfejsu czyta się
jak usterka platformy.

Wykaz jest przepisem rozdziału „Katalog funkcji i narzędzi" dokumentacji
modułu, nie zbiorem pomysłów; liczba jego pozycji jest zmierzona z długości
wykazu, nie wpisana osobno.

## budowa/klient-poprzedni/src/moduly/roundtable/katalog-funkcji.ts

Katalog istnieje po to, żeby okno nie wyglądało na kompletne tam, gdzie nie
jest. Opracowanie modułu wymienia siedemdziesiąt cztery narzędzia — bez wykazu
Operator poznawałby rozjazd między nimi a stanem modułu dopiero po naciśnięciu
każdego przycisku z osobna.

Po scaleniu kontraktu rozjazd zmienił naturę. Obszar komend modułu niósł
cztery komendy i niesie ich dziś czterdzieści sześć, więc zdanie „tego nie
wykona żadna komenda obszaru" przestało być prawdziwe co do czterdziestu
dwóch narzędzi naraz. Brakuje ich obsługi, nie kontraktu — i to rozróżnienie
niesie stan pozycji.

Nazwy narzędzi stoją w brzmieniu opracowania modułu, bez numeracji
rozdziałów — numer rozdziału nie jest identyfikatorem niczego w produkcie.
Warstwa pozycji jest warstwą widoczności z opracowania, przypisaną po
elemencie interfejsu, którym narzędzie się otwiera: pierwsza to widoczność
bez interakcji, druga to znacznik kontekstowy, przycisk albo przełącznik,
trzecia to menu zestawu akcji okna, a czwarta to nastawa dostępna poleceniem,
wyszukiwarką funkcji albo trybem administracyjnym — w oknie nie ma jej
w stanie spoczynku.

Stan wykonania narzędzia niesie pięć wartości. „Częściowa" nie jest stanem
pośrednim między porażką a sukcesem: znaczy, że narzędzie działa w części
zakresu, a nie działa w reszcie, i zdanie pozycji mówi, gdzie przebiega
granica. „Bez obsługi" znaczy, że komenda jest w kontrakcie, a moduł jej nie
wywołuje — to najliczniejszy stan modułu i najważniejszy do odróżnienia: nie
brakuje uzgodnienia, brakuje pracy, a zdanie pozycji nazywa komendę, którą
narzędzie wejdzie. „Bez pokrycia" znaczy, że kontrakt nie ma czym narzędzia
wykonać — nie ma komendy albo ma komendę, lecz brakuje pola w jej żądaniu; po
scaleniu kontraktu stan ten zszedł do pojedynczych pozycji i każda nazywa
brakujące pole. „Poza modułem" znaczy, że narzędzie należy do okna wspólnego
platformy.

## budowa/klient-poprzedni/src/moduly/studio/schowek-zrodlo.ts
Schowek przeglądarki żyje tylko w obrębie karty i nie ma historii, ponieważ
`navigator.clipboard` oddaje jedną, ostatnią treść. Wymaganie produktu
obejmuje wykaz wpisów sprzed kilku ruchów oraz wpisy przypięte na stałe,
więc historia leży w rdzeniu poprzez rodzinę komend `clipboard.*`: `push`
dopisuje treść (powtórzenie identycznej treści nie mnoży wpisów, tylko
podnosi zastany wpis na czoło i oddaje `alreadyPresent`), `list` oddaje
historię wraz z przypiętymi, `pin` przypina wpis (przypięty nie wygasa wraz
z retencją), `delete` usuwa wpis albo całą historię nieprzypiętą.

Komendy schowka nie niosą pola autora, więc wpis odłożony przez model i wpis
operatora są w rdzeniu nierozróżnialne. Odkładając fragment za modelem, okno
zapisuje pole `sourceWindowId` — jedyne pole pochodzenia, które wpis niesie.
Rozróżnienie autora wpisu schowka nie wchodzi w zakres kontraktu.

## budowa/klient-poprzedni/src/moduly/workspace/stany-okna.ts
Przedrostek klas CSS „dw-" jest ustalany lokalnie, bo nośnik stanu treści jest współdzielony
między modułami, a każdy moduł nadaje własny przedrostek. Stan błędu przenosi kod i komunikat
z kontraktu, aby okno po niepowodzeniu zapytania nie wyglądało jak okno z pustą, ale poprawną listą.

## budowa/klient-poprzedni/src/moduly/research/dymek-badania.ts
Znak dymku Research zachowuje trzy własności wspólnej fabryki bez zmian: pokazanie na najechaniu albo ustawieniu ogniska bez kliknięcia i bez osobnego zamykania, prowadzenie ogniska przyciskiem z odpowiedzią na każde naciśnięcie oraz treść w atrybucie opisującym, dzięki czemu dymek nie potrzebuje identyfikatora i nie zderza się między oknami. Znak Research różni się od biblioteki wyglądem — jest pierścieniem o średnicy czternaście do szesnastu pikseli ze wskaźnikiem pomocy, a nie kwadratowym przyciskiem biblioteki — dlatego klasa znaku zastępuje klasę biblioteczną, a klasa powłoki dokłada się do klasy dymka.

## budowa/klient-poprzedni/src/sterowanie/zrodlo-kanalow.ts
Zakładanie stoi w tym samym źródle co usuwanie, ponieważ świeża baza ma dokładnie jeden kanał, a okno komunikacji kanału wymaga; panel z samym usuwaniem pozwoliłby nieodwracalnie skasować jedyny kanał produktu bez drogi powrotu. Wykazu kanałów to źródło nie pobiera — rejestr kanałów jest czytającą pamięcią podręczną nad odczytem wykazu, wspólną całemu klientowi, i po udanym zapisie panel odświeża tę pamięć. Każda czynność oddaje wynik z kodem odmowy, nigdy samej treści, aby odmowa rdzenia dojechała do widoku, a nie wyglądała jak wykonanie. Pole rodzaju kanału nie występuje w opisie zmiany, ponieważ kontrakt zmiany rodzaju kanału po założeniu nie przewiduje, a formularz komunikuje to wprost zamiast przyjmować wpis, który rdzeń zignoruje. Sprawdzenie łączności jest narzędziem pomocniczym, nie bramką — jego wynik nie wstrzymuje zapisu, nie wyłącza wiersza ani nie blokuje wysłania tury, bo kanał niedostępny przed chwilą bywa sprawny teraz. Odpowiedź o stanie poświadczenia mówi, czy poświadczenie jest ustawione, jakiego jest rodzaju i gdzie mieszka, lecz nigdy nie niesie samego sekretu.

## budowa/klient-poprzedni/src/uwierzytelnienie/sesja-bramki.ts

Transportem jest jedno gniazdo WebSocket, nie seria żądań HTTP, więc ciasteczka
nie ma: token wraca z `auth.login` i tylko klient może go przechować między
uruchomieniami. Bez zapisu każdy start aplikacji wymagałby hasła, a bramka ma
stać przy wejściu raz, nie przy każdym otwarciu okna.

W zapisie leży sesja w kształcie kontraktu (`AuthSession`): token, czas
wygaśnięcia, metoda. Token jest poświadczeniem na okaziciela — rdzeń trzyma
w bazie wyłącznie jego skrót, więc jedyną kopią tokenu jest właśnie ten zapis.
`localStorage` interfejsu podawanego z `127.0.0.1` i z powłoki natywnej jest
magazynem lokalnym tej samej maszyny, czyli tym samym progiem zaufania, na
którym stoi plik sejfu rdzenia obok bazy.

Magazyny są dwa, bo wybór trwałości logowania jest jawny. Pole „Nie wyloguj
mnie" na ekranie logowania rozstrzyga, gdzie sesja wyląduje: zaznaczone kładzie
ją w `localStorage` i sesja przeżywa zamknięcie aplikacji, niezaznaczone —
w `sessionStorage`, skąd ginie razem z oknem.

Skutek jest dwustronny: ta sama wartość idzie żądaniem `auth.login` jako
`keepSignedIn` i rozstrzyga w rdzeniu o trwaniu sesji — 365 dni zamiast doby
roboczej. Trwanie jest zapisane przy wierszu sesji, więc `auth.token.refresh`
go nie ścina. Zapis w przeglądarce bez tego byłby obietnicą, której rdzeń nie
dotrzymuje.

Odczyt pyta obu magazynów, w kolejności od trwałego. Inaczej zmiana
rozstrzygnięcia między jednym a drugim uruchomieniem zostawiałaby sesję
niewidoczną, a wyglądałoby to jak jej wygaśnięcie. Kasowanie czyści oba
z tego samego powodu.

Żaden błąd pamięci nie zatrzymuje uruchomienia: brak magazynu, zapis nieczytelny
i kształt spoza kontraktu znaczą to samo — sesji zapisanej nie ma i wejście
wymaga hasła. O tym, czy sesja jeszcze żyje, rozstrzyga wyłącznie rdzeń
odpowiedzią na `auth.token.refresh`; ten plik nie porównuje czasów i niczego
nie unieważnia sam.

Zapisana sesja leży wyłącznie w jednym magazynie naraz: pole „Nie wyloguj mnie"
nie pamięta własnego zaznaczenia osobnym zapisem, tylko czyta miejsce, w którym
sesja naprawdę leży. Osobny zapis rozjechałby się z magazynem przy pierwszym
czyszczeniu pamięci i pokazywałby zaznaczone pole nad sesją, która ginie
z oknem.

Czas ważności idzie z odpowiedzi rdzenia (`expiresAt`), nie ze stałej klienta —
długość życia sesji zna wyłącznie rdzeń. Data pełna pojawia się tylko wtedy,
gdy wygaśnięcie wypada innego dnia; w dniu bieżącym wystarcza godzina.

## budowa/klient-poprzedni/src/rozmowa/blok-zwijany.ts
Zwinięcie bloku opiera się na natywnym zachowaniu przeglądarki dla elementów details i summary, dzięki czemu blok działa klawiaturą, ma poprawną semantykę dla technologii wspomagających i nie wymaga obsługiwaczy zdarzeń. Brak treści bloku prowadzi do jego ukrycia, a nie wyszarzenia.

Sterowanie rozwinięciem z zewnątrz jest wołane wyłącznie przy zmianie trybu widoku, nigdy przy odświeżeniu treści: blok przestawiony ręcznie pozostaje w stanie, w jakim go zostawiono, także gdy w trakcie tury przychodzą kolejne fragmenty.

## budowa/klient-poprzedni/src/moduly/studio/stan-okna-studio.ts
Rozróżnienie czterech faz jest treścią: stan sprzed pierwszego wywołania, stan
w trakcie wywołania i stan odmowy rdzenia są osobnymi fazami, ponieważ zlanie
ich zmusza czytającego do zgadywania, czy czekać, czy działać. W Studio różnica
jest ostra, bo rdzeń oddaje tu i treść, i pustkę, i odmowę: otwarcie dokumentu
ze ścieżką oddaje dokument pusty, wykaz repozytorium przed pierwszym zapisem
oddaje wykaz pusty, porównanie różnic bez wskazanej strony oddaje pustkę,
a otwarcie nieznanego dokumentu odmawia kodem błędu. Pusty widok bez powodu
byłby nie do odróżnienia od żadnej z tych czterech sytuacji.

Komunikat nie kasuje treści, tylko ją przesłania; wyjątkiem jest faza
ładowania, w której treść ustępuje miejsca wskaźnikowi odczytu. Stan pusty ma
jedną formę złożoną ze znaku, tytułu i opisu; znak dochodzi z zestawu ikon,
a nie z własnego rysunku, bo forma stanu pustego wariantów stylu nie ma.

Zestaw faz i ich znakowanie pochodzą ze wspólnego modułu faz okna: ta sama
nazwa fazy trafia do atrybutu danych w każdym module, więc wspólny selektor
i wspólny sprawdzian mają się o co oprzeć. W tym module zostaje wyłącznie to,
czym Studio różni się świadomie: wskaźnik odczytu i chowanie treści.

## budowa/klient-poprzedni/src/moduly/research/formularz-zrodla.ts
Formularz wydzielono z okna zarządzania źródłami, ponieważ katalogowanie źródła i ocena jego wiarygodności to osobne funkcje operatora z własnymi polami w kontrakcie (rodzaj, pochodzenie, adres, dokument repozytorium, ocena wiarygodności). Każde pole niesie dymek objaśnienia, bo każde jest elementem konfiguracji zlecenia. Rodzaj źródła i ocena wiarygodności idą przez wspólny wybór nastawy — obsadę bibliotecznego menu drzewa, gdzie uchwyt pokazuje wartość bieżącą, a natywna lista rozwijana pokazywałaby ją dopiero po rozwinięciu; pola tekstowe zostają polami tekstowymi, bo nie ma w nich czego rozwijać.
Rzutowanie wartości sterowania na wyliczenie kontraktu odbywa się wyłącznie w funkcji pomocniczej: ster oddaje napis, którego pozycje pochodzą z wyliczeń kontraktu, więc innej wartości wydać nie może, a pusty wybór wraca wartością zastępczą, nie napisem pustym.

## budowa/klient-poprzedni/src/moduly/wiedza/indeks.ts
Katalog wiedzy stoi osobno od katalogu poczty, bo to dwie różne dziedziny rdzenia: poczta sięga
po skrzynkę Operatora stojącą poza urządzeniem, wiedza po treści własne leżące w rdzeniu. Wspólny
katalog zlepiłby dwa słowniki pojęć w jeden. Rejestr modułów wiąże widok z kodem modułu rdzenia,
a tabela modułów nie niesie osobnego kodu ani dla poczty, ani dla wiedzy — drugi wpis rejestru
czekałby na kod, którego nawigacja nigdy nie poda. Okno wchodzi do sceny przez złożenie modułu
poczty, które ma w rejestrze jeden wpis. Dlatego ten plik nie wystawia stałej modułu.

## budowa/klient-poprzedni/src/moduly/translate/wiersz-terminu.ts
Jeden wiersz przypada na plik, tak samo jak wiersz biblioteki ekspertów w module agentów. Ten sam
wiersz nadaje się zarówno do wykazu terminów, jak i do zestawienia wystąpień, ponieważ nie wywołuje
działań samodzielnie, tylko oddaje naciśnięcia oknu, które wie, co z nimi zrobić.

## budowa/klient-poprzedni/src/moduly/roundtable/zestawienie-udzialu.ts
Moduł mierzy wyłącznie udział: liczbę wypowiedzi i znaków na uczestnika w turze bieżącej. Okno opisuje głosowanie, ocenę rubrykową, ocenę modelami-sędziami i ranking akumulowany między sesjami; odpowiadające im komendy kontraktu (roundtable.vote.start, roundtable.vote.cast, roundtable.vote.get, roundtable.rubric.set, roundtable.judge.run, roundtable.leaderboard.get) istnieją, ale moduł żadnej z nich jeszcze nie wywołuje. Wynik liczony jako udział nie jest oceną jakości wypowiedzi ani progiem quorum — próg zgody niesie pole quorum żądania roundtable.vote.start, którego okno nie wywołuje.
Wypowiedź otwartą ostatnio liczy się z tekstu widocznego w oknie (strumień), nie z pola content: rdzeń rozgłasza wypowiedź najpierw pustą i dopisuje treść strumieniem, więc sam zapis dawałby zero znaków przez cały czas mówienia modelu. Wypowiedzi wcześniejsze są już utrwalone i liczą się z zapisu.

## budowa/klient-poprzedni/src/okna-pomocnicze/indeks.ts
Plik był punktem zbiorczym eksportów całego obszaru okien pomocniczych, więc pierwotny
nagłówek opisywał po kolei skład obszaru, powód wydzielenia go poza katalog modułów
oraz ryzyko rozejścia się kopii przy duplikacji w dwóch modułach. Redakcja scaliła te
trzy wątki w jedno zdanie nazywające zawartość i wspólność dla modułów Developer
i Diagnostics, zachowując wymienione elementy obszaru bez warstwy uzasadniającej wybór
lokalizacji pliku w drzewie źródeł.

## budowa/klient-poprzedni/src/moduly/research/indeks.ts
Powłoka zna stąd jedną rzecz — opis modułu. Kod modułu stoi wewnątrz opisu, bo moduł sam mówi, którym jest modułem, więc rozjazd między nazwą w rejestrze a rzeczywistością staje się niemożliwy. Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element, a warstwa składająca decyduje, gdzie go postawić, dzięki czemu te same okna wchodzą i w obszar roboczy powłoki, i w stanowisko sprawdzianu. Funkcja zamknięcia podpina rozłączenie modułu, bez którego subskrypcja stanu badania i kreator raportu żyłyby dalej po zejściu modułu ze sceny; pole jest w opisie widoku modułu opcjonalne, więc kompilator braku nie zgłosi.

## budowa/klient-poprzedni/src/okna-rownolegle/uchwyt-szerokosci.ts
Obsługa klawiaturą jest obowiązkowa, ponieważ uchwyt osiągalny wyłącznie
wskaźnikiem odcinałby od jedynego sposobu ustawienia szerokości kolumny; stąd
rola separatora, indeks tabulacji, strzałki, Home/End i pełny komplet
atrybutów wartości ARIA. Uchwyt nie zna paneli ani rozmowy — oddaje wołającemu
jedną liczbę w pikselach bez własnych założeń, a przycięcie do minimów
wykonuje osobny moduł odpowiedzialny za szerokości gniazda. Uchwyt stoi po
lewej stronie sterowanej kolumny zgodnie z układem gniazda: rozmowa, uchwyt,
panele. Ruch w lewo poszerza kolumnę, ruch w prawo ją zwęża; strzałki
klawiatury podążają tą samą logiką kierunku, żeby obsługa ręką i klawiaturą
dawały ten sam skutek. Krok strzałki wynosi wielokrotność czterech pikseli
siatki interfejsu — krok jednopikselowy wymagałby wielokrotnego naciskania
strzałki, a stupikselowy uniemożliwiałby dojście do wartości osiąganej
wskaźnikiem bez wysiłku.

## budowa/klient-poprzedni/src/moduly/studio/stan-studio.ts
Pięć okien modułu Studio — edytor, panel narzędzi, panel różnic i wyszukiwania,
podgląd oraz repozytorium sesji — pracuje na tym samym dokumencie czynnym.
Gdyby każde okno prowadziło własną kopię, przywrócenie wersji w repozytorium
nie przestawiłoby treści edytora, a różnica dotyczyłaby innego dokumentu niż
podgląd. Poza własnym działaniem okien jedynym źródłem odświeżenia stanu jest
zdarzenie zmiany dokumentu — zmiana dokonana w innym oknie albo na innym
urządzeniu konta dociera nim, bez odpytywania w pętli.

Zakres wysyłany do rdzenia operacją kontekstową bierze wskaźnik panelu, pasek
zaznaczenia edytora i żądanie operacji z tej samej drogi, żeby istniała jedna
nastawa i jedna prawda o niej. Powodu odmowy ostatniej operacji kontekstowej
pytają okna stojące obok panelu narzędzi, żeby brak propozycji po odmowie nie
wyglądał u nich jak sytuacja sprzed pierwszej operacji.

Parę wersji do porównania wskazuje repozytorium sesji przyciskiem porównania
przy wersji, a czyta panel różnic przy składaniu żądania; nastawa jest jedna
i idzie przez stan współdzielony, ponieważ dotyczy dwóch okien naraz.

Migawkę pól dokumentu czynnego bierze okno pracy przy przełączaniu zakładek:
dwa dokumenty naraz mają niezależne zaznaczenie, bufor i propozycję, a stan
modułu prowadzi jeden dokument czynny. Migawka jest jedyną drogą do tego —
druga kopia stanu modułu rozjechałaby się z nim przy pierwszym zdarzeniu
rdzenia. Nadejście wyniku operacji unieważnia odmowę poprzednią, ponieważ
komunikat o odmowie i komunikat o wyniku nie mogą stać obok siebie prawdziwe.

## budowa/klient-poprzedni/src/rozmowa/etykiety-rozmowy.ts
Napis nigdy nie stoi obok samej barwy: każdy stan niesie słowo, bo stan sygnalizowany wyłącznie kolorem jest niedopuszczalny.

Podpowiedź pod polem wpisywania mówi wprost, że pole jest zarazem filtrem wykazu, i wymienia klawisze, które przy otwartym wykazie znaczą co innego niż zwykle.

Napis wykazu wytworów mówi „narzędzia", a nie „pliki": rdzeń nie nadaje spisu wytworzonych plików, więc wykaz wymienia nazwy wywołanych narzędzi.

Opis stanu pustego trybu Streszczenie tłumaczy, czym streszczenie jest i skąd się bierze, zamiast meldować brak danych.

Odpowiedź na Zatrzymaj bez tury dotyczy przypadku, w którym rdzeń oddaje pole stopped równe fałsz, a kliknięcie i tak musi zostawić ślad na ekranie.

## budowa/klient-poprzedni/src/strona-glowna/czynnosci-sesji.ts
Plik nie zna kanału ani kontraktu — wykonanie czynności podaje wpięcie. Warunki dostępności czynności stoją w tym pliku, a nie w karcie czy w menu, ponieważ karta buduje wiersz, menu buduje przyciski, i żadne z nich nie rozstrzyga, czy na przykład wznowienie sesji czynnej ma sens; jedno miejsce daje jedną odpowiedź i jedno miejsce do poprawienia, gdy rdzeń zmieni stany sesji. Brak czynności zdejmuje przycisk z menu, zamiast go wyszarzać, ponieważ przycisk widoczny, którego naciśnięcie nic nie robi, jest atrapą — wiersz bez czynności pozostaje czysto informacyjny. Zatrzymanie tur ma sens wyłącznie przy turze w biegu, bo zero okien strumieniujących znaczy, że komenda nie miałaby czego zatrzymać. Wznowienie dotyczy sesji, która nie biegnie — wstrzymanej albo zakończonej — ponieważ sesja czynna nie ma czego wznawiać.

## budowa/klient-poprzedni/src/okna-rownolegle/uwaga-podgladu.ts
Komunikat jest blokiem pełnej szerokości rodzica, oznaczonym wyłącznie cienką
kreską po lewej krawędzi, bez tła i obwódki, żeby nie czytał się jako druga
karta. Forma stoi w arkuszu stylów tego widoku pod nazwą własną, ponieważ
biblioteka współdzielona nie niesie gotowej klasy komunikatu blokowego.
Komunikat pozostaje widoczny na scenie, dopóki nie ustanie jego przyczyna —
brak komendy przekazania w kontrakcie — dzięki czemu odpowiada na pytanie
o brak działania także po czasie.

## budowa/klient-poprzedni/src/moduly/research/kody-okien.ts
Kody okien trzymane w jednym miejscu sprawiają, że dopisanie okna po stronie rdzenia rozjeżdża się z klientem widocznie, a nie w siedmiu plikach naraz. Dwa kody, discovery-panel i reading-view, nie mają jeszcze wiersza w katalogu rdzenia — migracje 030 i 031 zakładają dla modułu Research pięć okien, a opracowanie modułu wylicza siedem własnych obok dwóch wspólnych platformie. Kody stoją w brzmieniu, w jakim mają wejść do katalogu, dzięki czemu atrybut okna obu okien jest już dziś tym samym napisem, którym będzie po dobudowie rdzenia, i nawigacja wewnątrzmodułowa nie wymaga zmiany.
Wykaz okien spoza katalogu jest zaporą, nie zgodą: okna są zbudowane i czynne, a ich kody wchodzą do katalogu akcji tak samo jak kody okien znanych. Rdzeń odmówi im na katalogu okien, zanim dojdzie do katalogu akcji, i nazywa okno operatorowi wprost, zamiast pokazywać puste miejsce. Wykaz znika, gdy migracje 030 i 031 dostaną wiersze obu okien.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-okna-translate.ts
Komendy wykazu okien sesji i odczytu parametrów wykonania wymagają identyfikatora okna, który
bierze się z rdzenia, nigdy z literału po stronie klienta. Idą zwykłą drogą wywołania, nie drogą
odmowy, ponieważ rdzeń je obsługuje i odpowiedź przychodzi kopertą ze statusem, którą korelacja
rozpoznaje.

## budowa/klient-poprzedni/src/strona-glowna/etykiety-sesji.ts
Żadna funkcja tego pliku nie wymyśla danych: każda przekłada wartość oddaną przez rdzeń, a wartość nieobecną oddaje jako `undefined`, dzięki czemu karta pomija cały fragment zamiast pokazać wartość zmyśloną. Nazw środowisk plik nie niesie — daje je osobne źródło zasilane wykazem z rdzenia, żeby na jednym ekranie nie stały dwa źródła tej samej nazwy.

## budowa/klient-poprzedni/src/moduly/research/kreator-raportu.ts
Kreator raportu ma stany zamknięty, otwierający się, otwarty, próbę z brakami, ładowanie i błąd. Wzorzec czterech dróg zamknięcia (kontrolka w nagłówku, klawisz Escape, kliknięcie w nakładkę, akcja w stopce) powtarza wzorzec okna konfiguracji. Przycisk główny jest czynny od otwarcia; próba z brakami nie gasi kontrolki — kreator zostaje otwarty, a brak wraca komunikatem nad stopką, a ponowne naciśnięcie w trakcie wywołania jest bezpieczne po stronie logiki kreatora, nie przez odebranie klikalności.
Reguła wyboru drogi budowy raportu: żądanie z polem sekcji składa raport z samych sekcji podanych i pomija identyfikatory ustaleń; bez tego pola woła kanał modelu po sekcję nadrzędną. Reguła stoi w kreatorze, bo tu Operator rozstrzyga ją każdym wpisanym znakiem. Droga modelu nie gwarantuje odpowiedzi modelu — bez logowania budowa kończy się powodzeniem, raport zostaje zapisany, a treścią sekcji nadrzędnej jest komunikat procesu kanału; kreator uprzedza o tym przed naciśnięciem, bo po nim klient nie ma już czym rozpoznać pochodzenia treści.

## budowa/klient-poprzedni/src/rozmowa/historia-tur.ts
Komenda message.send zwraca wiadomość użytkownika, a identyfikator odpowiedzi znany jest dopiero z pierwszego fragmentu strumienia. Dlatego wpis odpowiedzi powstaje od razu po wysłaniu, jako wpis oczekujący, a pierwszy fragment go przejmuje zamiast zakładać drugi. Bez tego okno stałoby puste do nadejścia pierwszego znaku, a potem pokazało dwa wpisy zamiast jednego.

Wiązanie wypowiedzi z wiadomością rdzenia dotyczy sytuacji, w której wypowiedź trafia do historii natychmiast, przed komendą, więc nie ma jeszcze identyfikatora wiadomości: message.send oddaje go dopiero w odpowiedzi. Rdzeń rozgłasza tę samą wiadomość zdarzeniem message.changed do wszystkich połączeń konta, także do tego, które ją wysłało. Ponieważ okno pokazuje wypowiedzi roli user, to samo zdanie weszłoby do wątku dwa razy: raz jako echo miejscowe, raz ze zdarzenia. Wiązanie idzie po treści, nie po identyfikatorze, bo identyfikatora w tej chwili nie ma po żadnej stronie: dopasowuje pierwszą niezwiązaną wypowiedź o tej samej treści i nadaje jej identyfikator z rdzenia. Metoda zwraca wpis związany albo wartość pustą, gdy wypowiedzi o tej treści w historii nie ma, czyli gdy wiadomość naprawdę przyszła skądinąd.

Czyszczenie historii nie zeruje licznika kluczy. Klucz jest tożsamością pozycji w widoku; gdyby po wyczyszczeniu zaczął się od nowa, pierwszy wpis nowej rozmowy trafiłby w węzeł pozostały po wpisie rozmowy poprzedniej i zamiast nowej pozycji stanęłaby podmieniona stara.

## budowa/klient-poprzedni/src/moduly/translate/odmowa-translate.ts
Dwie komendy modułu wykonuje rdzeń modelem: przekład tekstu źródłowego na język panelu i rozpoznanie
języka tekstu. Gdy rejestr nie ma czynnego kanału, rdzeń odmawia obu z jednym kodem niedostępności
kanału. Sama wiadomość rdzenia jest powodem, ale nie mówi Operatorowi, co ma zrobić, a jej brzmienie
sugeruje usterkę wewnętrzną tam, gdzie stoi zwykły brak włączonego kanału — zdanie o brakującej rzeczy
dokłada klient, idąc za powodem rdzenia, nie zastępując go. Zdania nie dokłada wspólny tłumacz kodów,
bo ten sam kod niesie w module dwie różne odmowy: powyższą oraz odmowę importu glosariusza, gdzie
rdzeń nie czyta plików z dysku Operatora i zdanie o kanale modelu byłoby nieprawdziwe. Zdanie mówi też,
dlaczego Operator nie ma tu pola wyboru: żądania modelowe nie niosą identyfikatora kanału, więc rdzeń
bierze kanał domyślny czynny, a klient nie wymyśla pól kontraktu.

## budowa/klient-poprzedni/src/moduly/translate/pola-zrodla.ts
Skład pola — etykieta, podpowiedź, dymek, wykaz podpowiadanych wartości — to inna odpowiedzialność
niż to, co dzieje się po naciśnięciu przycisku, dlatego pola są wydzielone z okna, które składa
całość i wiąże zdarzenia. Dymek objaśnienia dostaje każde pole zmieniające treść żądania: język
źródłowy oraz ponowna segmentacja.

## budowa/klient-poprzedni/src/strona-glowna/formularz-komponentu.ts
Formularz stoi zwinięty, ponieważ środkiem ciężkości strony głównej są karty środowisk, a formularz stale rozłożony odbierałby im pas ekranu przy czynności wykonywanej rzadko; rozwija go segment „Dodaj nowy” belki strefy trzeciej. Wykaz rodzajów komponentu odmawia rodzaju bez magazynu profili, a powód stoi na ekranie, nie tylko w komentarzu, inaczej zostałaby sama krótsza lista bez wyjaśnienia. Zakładanie i zmiana komponentu idą jedną szufladą, ponieważ są jedną pracą w dwóch krokach, a osobny uchwyt do drugiego kroku kazałby szukać go po założeniu komponentu; drugi uchwyt do jednej szuflady byłby drugą drogą do jednego bytu. Doklejenie kafla z odpowiedzi zakładania pokazałoby stan, którego rdzeń nie potwierdził drugim odczytem, dlatego wykaz dociąga strefa osobnym zapytaniem.

## budowa/klient-poprzedni/src/okna-rownolegle/widok-pelnoekranowy.ts
Scena dzieli szerokość między trzy obszary: kolumnę rozmowy, kolumnę paneli
i widok pełnoekranowy. Ten trzeci jako jedyny nie ma sąsiada — bierze całą
scenę, bo panel w kolumnie trzyma minimum szerokości, a treść, która się
w nim nie mieści, potrzebuje sceny, nie kolumny. Treść jest przenoszona, nie
kopiowana, żeby uniknąć drugiej subskrypcji rdzenia dla tego samego panelu
i subskrypcji osieroconej po zamknięciu jednego z egzemplarzy; funkcja
pokazująca zapamiętuje rodzica i następnik, a funkcja chowająca odkłada treść
dokładnie tam. Klawisz Escape wychodzi z widoku dzięki nasłuchowi osadzonemu
na elemencie widoku, nie na dokumencie, dzięki czemu sprzątanie nasłuchu jest
zbędne i klawisz nie zabiera działania niczemu innemu na scenie. Wyjście nie
jest zamknięciem panelu: treść wraca do stosu, a panel żyje dalej ze swoją
subskrypcją, natomiast właściwe zamknięcie panelu należy do jego obudowy
w stosie. Widok nie buduje treści i nie rozstrzyga zachowania kolumn sceny —
to należy do układu, który go osadził.

## budowa/klient-poprzedni/src/moduly/research/modul-research.ts
Układ wynika z roli okna i z porządku pracy badawczej: wyszukaj, skataloguj, przeczytaj, odnotuj, złóż, wydaj. Research Workspace jest wiodące i jest punktem wejścia, więc stoi w pasie pierwszym na całą szerokość. Dalej idą pary, w których biegnie wiązanie: Discovery Panel obok Sources Manager, bo pozycja wyniku staje się źródłem; Reading View obok Findings Panel, bo wypis z lektury staje się ustaleniem; Report Builder obok Export Panel, bo dokument staje się plikiem. Badanie jest jedno na cały moduł: wybór źródeł, zaznaczenie ustaleń i wskazanie materiału do lektury przestawiają wszystkie siedem okien naraz, bo stan badania jest jeden. Nawigacja wewnątrzmodułowa jest przeniesieniem ogniska, nie zmianą trasy: okna stoją obok siebie, więc odnośnik do innego okna prowadzi wzrok i ognisko, zamiast wymieniać zawartość obszaru.

## budowa/klient-poprzedni/src/strona-glowna/indeks.ts
Powłoka aplikacji sięga po budowę strony głównej wyłącznie przez ten plik i nie zna podziału widoku na strefy, karty ani kafle, dzięki czemu zmiana wewnętrzna katalogu strony głównej nie dotyka niczego poza nim.

## budowa/klient-poprzedni/src/moduly/research/nasluch-odmow.ts
Odmowa nieznanej komendy wraca kopertą bez pola statusu, a klient uznaje za odpowiedź wyłącznie kopertę ze statusem. Korelacja żądania takiej koperty nie rozstrzyga, więc obietnica wywołania zostałaby nierozstrzygnięta na zawsze, a okno wisiałoby w stanie ładowania — odmowę trzeba więc rozpoznać po ładunku, nie po nazwie zdarzenia. Nasłuch wiąże identyfikator żądania z ładunkiem odmowy nieznanej komendy: zwycięża to, co przyjdzie pierwsze — odpowiedź korelacji albo odmowa nieznanej; połączenia nic to nie kosztuje, rdzeń pracuje dalej.
Kod niepowodzenia not_found mówi prawdę: bytu o tej nazwie w rdzeniu nie ma. Nazwa typu zostaje w polu osobnym, żeby okno mogło ją wypisać wprost, zamiast pokazywać pusty wykaz udający brak danych.

## budowa/klient-poprzedni/src/moduly/research/odczyt-badania.ts
Odczyt ma jedną odpowiedzialność: dwa kroki odczytu i przełożenie ich niepowodzeń na fazę pamięci. Wydzielony jest z pamięci, bo pamięć nie zna kontraktu, a odczyt zna wyłącznie kontrakt — rozdzielenie trzyma obie rzeczy w rozmiarze, w którym dają się przeczytać naraz. Kroki są dwa, bo mówią o dwóch różnych rzeczach: wykaz okien mówi, które okno sesji należy do modułu Research, a stan okna mówi, co w nim jest. Klient nie zgaduje identyfikatora okna — bez wskazania rdzenia zostaje uczciwy stan pusty nazywający brak, nie zaszyta wartość.

## budowa/klient-poprzedni/src/uwierzytelnienie/tozsamosc-urzadzenia.ts

PIN jest właściwy maszynie, nie osobie wchodzącej: rdzeń odnajduje metodę po
parze rodzaj i urządzenie, więc PIN założony pod jednym identyfikatorem nie
otworzy bramki podanym pod innym. Ekran logowania musi znać tę wartość, żeby
wiedzieć, czy segment „PIN" ma w ogóle stanąć, i żeby móc PIN wysłać.

Bramka wyłącznie czyta identyfikator; nadaje go zakładanie PIN-u w Ustawieniach,
czyli czynność wykonana na urządzeniu. Ekran logowania staje przy każdym
uruchomieniu, także pierwszym na świeżej przeglądarce — gdyby nadawał
tożsamość, zapisywałby maszynę, która żadnego PIN-u nie ma. Pusty wynik jest
prawdą: ta przeglądarka jeszcze się nie przedstawiła.

Ten sam klucz niesie plik tożsamości urządzenia w module ustawień i tak samo
rozdziela odczyt od nadania. Bramka nie importuje pliku ustawień, bo warstwa
logowania leży niżej niż okno ustawień i zależność szłaby pod prąd warstw;
scalenie obu plików wymaga warstwy wspólnej poniżej nich.

## budowa/klient-poprzedni/src/okna-rownolegle/wiersz-czynnosci.ts
Wiersz czynności stoi osobno od wiersza panelu, bo wiersz panelu jest
przełącznikiem — niesie stan zaznaczenia i stałe miejsce na ptaszek, ponieważ
panel bywa otwarty albo zamknięty, natomiast czynność sesji nie ma stanu
włączenia, więc ptaszek kłamałby o tym, czym pozycja jest. Układ pozostaje ten
sam co w sekcji paneli: ikona, nazwa i przeznaczenie, skrót wyrównany do
prawej, żeby dwie sekcje jednego menu nie miały dwóch rytmów. Skrót
klawiszowy jest tu prawdziwy, bo nasłuch stoi na sekcji i łapie literę, dopóki
menu jest rozwinięte; skrótu globalnego celowo nie rejestruje się, żeby litery
wciśnięte w polu wypowiedzi pisały tekst, a nie wykonywały czynność. Naciśnięcie
wiersza nie zwija menu — tak samo jak przy przełączaniu paneli — a modal
pytania o nazwę oraz modal potwierdzenia usunięcia są natywnymi oknami
z własną nakładką, stającymi nad menu, więc nie ma czego chować. Wyróżnienie
ostrzegawcze jest wyłącznie barwą wiersza, nie blokadą: czynność usuwająca
zapis bez odwrotu idzie czerwienią, ale wiersz pozostaje tak samo klikalny jak
każdy inny.

## budowa/klient-poprzedni/src/rozmowa/historia-z-rdzenia.ts
Komenda message.list daje wiadomości leżące w bazie, więc odświeżenie okna albo powrót do sesji pokazuje wątek, a nie pustkę. Wpisy odtworzone są domknięte: pochodzą z zapisu, nie ze strumienia, więc nie mają tury w biegu ani fragmentów do doliczenia. Wywołujący dostaje dwie drogi, przyjmij na sukces także z pustą tablicą, i zglosNiepowodzenie na odmowę z gotowym zdaniem złożonym z powodu koperty, i ma obsłużyć obie. Odmowa rdzenia i historia naprawdę pusta muszą wyglądać na ekranie inaczej.

Treść wpisu bierze się z pola content; wiadomość przerwana albo błędna wraca z treścią, którą zdążyła zebrać, bo pokazanie jej jest uczciwsze niż ukrycie tury, która się odbyła. Stan wpisu bierze się z pola status wiadomości, nie z samego faktu zapisu: tura przerwana i tura zamknięta błędem mają po odświeżeniu okna wyglądać tak samo jak w chwili, w której się odbyły.

Odczyt wykazu blocks jest tolerancyjny wzorem odczytu fragmentu: brak obszaru, obszar cudzego kształtu albo pozycja bez rodzaju dają mniej bloków, bo historia ma się wyświetlić, a nie zniknąć od nieczytelnej pozycji.

Wypowiedź operatora nie ma stanu tury: jest tym, co napisał, i nie może wrócić z zapisu jako zamknięta bez odpowiedzi, nawet gdyby zapis był pusty.

## budowa/klient-poprzedni/src/moduly/studio/statystyka-roznicy.ts
Panel różnic i wyszukiwania żąda paska statystyki oraz filtrowania różnic po
rodzaju. Kontrakt nie niesie osobnej komendy liczącej statystykę, ponieważ
porównanie stron oddaje fragmenty wraz z treścią przed i po, więc liczby
wynikają z odpowiedzi, którą panel już ma; drugie wywołanie po te same liczby
byłoby pytaniem o coś, co leży na stole. Moduł nie zna struktury dokumentu:
wejściem są fragmenty kontraktu, wyjściem liczby i wykaz przefiltrowany,
dzięki czemu rachunek sprawdza się bez stawiania okna.

## budowa/klient-poprzedni/src/strona-glowna/kafel-komponentu.ts
Kafel dzieli kartę ze środowiskiem, o mniejszej wadze: wariant komponentu odbiera wstęgę i cień sygnału w spoczynku, zmniejsza skalę ikony i zapisuje etykietę krojem bazowym półgrubym, nie nagłówkowym. Różnica krojów niesie znaczenie: krój nagłówkowy oznacza wejście do środowiska, bazowy oznacza zbudowanie komponentu, dlatego kafel nie dostaje klasy akcentu, którą nosi karta środowiska. Wiersz pusty metadanych wyglądałby jak metadane, których nie odczytano, dlatego kafel rodzaju wiersza w ogóle nie dostaje. Czas formatuje widok, bo rdzeń nie zna strefy czasowej Operatora; data pojawia się dopiero, gdy zmiana wypadła innego dnia, a w dniu bieżącym wystarcza sama godzina.

## budowa/klient-poprzedni/src/moduly/workspace/warstwy-instrukcji.ts
Warstwy instrukcji stoją w osobnym pliku od okna, bo to inna odpowiedzialność: okno prowadzi zapis
i odczyt, tu leży wyłącznie przedstawienie warstwy. Wykaz poziomów i ich nazwy pochodzą z modułu
poziomów zasięgu.

## budowa/klient-poprzedni/src/widok-sterowania/indeks.ts

Katalog nie buduje ani jednej kontrolki, tylko oprawę, w której kontrolki
katalogu stają się widoczne. Punkt wejścia klienta montuje widok dla okna
potwierdzonego przez rdzeń, tworząc rejestr kanałów, odświeżając go i wołając
montaż widoku sterowania z kanałem, oknem, rejestrem kanałów, panelem oraz
akcjami paska. Widok powstaje raz na okno; rejestr kanałów raz na klienta,
bo jest katalogiem wyboru, nie ustawieniem okna.

## budowa/klient-poprzedni/src/okna-rownolegle/wiez-koordynacji.ts
Ustalanie więzi wiąże pierwszy widoczny koordynator z pierwszym widocznym
wykonawcą; trzecie okno w roli wykonawcy pozostaje poza tą więzią, ponieważ
pas relacji pokazuje jedno powiązanie naraz, żeby obraz dał się ogarnąć
wzrokiem. Koordynator niesie napis informujący o zleceniu do wykonawcy,
a wykonawca napis o zleceniach otrzymanych od koordynatora; groty obu napisów
wskazują tę samą stronę sceny — kierunek biegu zlecenia — i stoją po stronie
okna partnera. U koordynatora grot wychodzi z napisu ku wykonawcy, u wykonawcy
wchodzi w napis od strony koordynatora, dzięki czemu kierunek pętli czyta się
bez czytania słów; gniazdo spoza więzi nie dostaje napisu.

## budowa/klient-poprzedni/src/moduly/translate/porownanie-paneli.ts
Zestawienie jest czynnością czysto miejscową: niczego nie liczy i o nic nie pyta rdzenia. Panele
są dokładnie dwa; przy innej liczbie wskazań zestawienie mówi wprost, czego oczekuje, zamiast
pokazać cokolwiek.

## budowa/klient-poprzedni/src/moduly/research/okno-discovery-panel.ts
Okno jest zbudowane, choć katalog okien rdzenia jeszcze go nie zna, bo trzon jego pracy ma dziś pokrycie w kontrakcie: dwa z czterech trybów zapytania idą własną komendą, a przeniesienie pozycji do katalogu idzie tą samą drogą co formularz Sources Manager. Ster trybu niesie wszystkie cztery pozycje, także te dwie, których rdzeń jeszcze nie obsługuje — wybór takiego trybu niczego nie gasi, zapytanie wychodzi pod nazwą swojej komendy, a rdzeń odmawia własnymi słowami; ukrycie tych pozycji zataiłoby zakres modułu, a ich wygaszenie łamałoby zasadę zero blokad. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności odkrywania.
Tryb semantyczny szuka po znaczeniu w bibliotece wiedzy, historii rozmów i przestrzeni roboczej okna badania komendą wyszukiwania wiedzy. Tryb pełnotekstowy szuka po treści w zasobach repozytorium komendą wyszukiwania pliku. Tryb webowy i tryb naukowy dzielą jedną komendę różniącą się wartością pola trybu; rdzeń nie ma dla nich jeszcze uchwytu, więc zapytanie wraca odmową.

## budowa/klient-poprzedni/src/widok-sterowania/montaz-widoku.ts

Kontrolki nie powstają w tym pliku — pochodzą w całości z katalogu sterowania:
środowisko wykonania, host, moduł, model, model zapasowy, nakład rozumowania,
tryb uprawnień, rola okna, katalogi robocze. Ten plik dokłada im miejsce,
nagłówek, podsumowanie wartości i dwa uchwyty rozwijania. Kolumna montuje się
obok sceny okna, a nie w prawym rogu paska górnego — dziewięć kontrolek nie ma
jak się tam zmieścić. Widok powstaje osobno dla każdego okna i domyka się na
jego identyfikatorze, więc dwa okna obok siebie dostają dwie niezależne kolumny
bez ani jednej wspólnej zmiennej. Szuflada startuje rozwinięta, żeby komplet
ustawień był widoczny bez szukania, co nacisnąć.

## budowa/klient-poprzedni/src/strona-glowna/karta-sesji.ts
Każdy napis karty pochodzi z rdzenia: tytuł to tytuł sesji, a gdy rdzeń go nie nadał — identyfikator sesji, nigdy nazwa wymyślona; środowisko, moduł i liczby okien niesie żywy odpis obecności, a przy odpisie nieobecnym cały fragment znika, zamiast pokazywać zero udające odczyt. Przycisk powrotu do sesji pojawia się wyłącznie wtedy, gdy montaż podał czynność powrotu; bez tej czynności wiersz jest czysto informacyjny, bez przycisku wyszarzonego ani martwego. Wykaz środowisk strony jest tym samym bytem, z którego rysuje się karta środowiska w strefie pierwszej — nazwa wzięta z klienckiej stałej rozjeżdżałaby się po cichu z nazwą zapisaną w bazie. Ponowne naciśnięcie przycisku powrotu w trakcie jego zajętości nie dubluje powiązania z sesją.

## budowa/klient-poprzedni/src/widok-sterowania/naglowek-widoku.ts

Tytuł jest jedynym miejscem kroju szeryfowego w tym widoku — krój ten należy
wyłącznie do nagłówków. Rola okna stoi w nagłówku, a nie tylko w podsumowaniu,
ponieważ to ona rozstrzyga, czym okno jest w pętli koordynator–wykonawca: przy
dwóch oknach obok siebie rola musi być widoczna bez rozwijania czegokolwiek.
Plakietka niesie ikonę słowną — nazwę roli — więc stan nie opiera się na samej
barwie.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-dokumentu-translate.ts
Obszar translate nie ma komendy przyjmującej plik: tekst źródłowy wchodzi napisem. Wydobycie tekstu
z dokumentu i rozpoznanie pisma ze skanu leżą w obszarze document i są jedynym wejściem od strony
pliku, więc moduł woła je wprost, tak jak woła wykaz okien i wykaz kanałów, zamiast je kopiować.
Wywołanie idzie zwykłą drogą protokołu, nie ścieżką odmowy właściwą wyłącznie obszarowi modułu,
bo odmowa obszaru dokumentów przychodzi kopertą ze statusem, którą korelacja rozpoznaje bez pomocy.
Ścieżka pliku jest ścieżką po stronie rdzenia — klient dysku nie czyta ani nie zapisuje, tylko podaje
wskazanie i oddaje Operatorowi odpowiedź rdzenia. Brak pola wymuszenia rozpoznania pisma znaczy, że
rdzeń bierze warstwę tekstową dokumentu, gdy ją ma — to zachowanie domyślne kontraktu, nie wybór okna.

## budowa/klient-poprzedni/src/rozmowa/metadane-konta.ts
Kanał nadaje fragment tego rodzaju także w chwili przełączenia konta w trakcie tury, żeby rotacja była widoczna. Struktura niesie wyłącznie odwołania: kod konta, nazwę profilu, odwołanie do danych dostępowych, a treść poświadczenia nie trafia tu nigdy. Kształt odpowiada strukturze metadanych konta po stronie rdzenia.

## budowa/klient-poprzedni/src/rozmowa/modul-okna.ts
Moduł okna jest czytany z rdzenia, a nie zgadywany po stronie widoku, dwiema drogami z kontraktu. Komenda window.state.get przy złożeniu pozwala oknu zapytać, w czym pracuje, zamiast czekać na pierwszą zmianę; bez tego okno wznowione po odświeżeniu klienta stałoby na module nieustalonym mimo znanego stanu w rdzeniu. Zdarzenie window.changed w toku pracy odzwierciedla sytuację, w której polecenie przestawienia obszaru roboczego przestawia to samo okno na inny moduł zamiast je zamykać, a rdzeń rozgłasza zmianę zdarzeniem. Niepowodzenie zapytania nie jest błędem okna: moduł zostaje nieustalony, a okno pracuje na arsenale wspólnym.

## budowa/klient-poprzedni/src/rozmowa/indeks.ts
Złożenie warstwy rozmowy w punkcie wejścia sprowadza się do wywołania zamontujRozmowe wewnątrz uzgodnienia otwarcia okna, z przekazaniem sceny, kanału, identyfikatora okna oraz persony i roli okna. Nazwy komend i zdarzeń, których warstwa używa — wysłania wiadomości, zatrzymania wiadomości, fragmentu strumienia, zmiany wiadomości, zmiany okna, odczytu stanu okna oraz wykazu akcji — pochodzą wyłącznie ze wspólnego kontraktu. Okno przestawia się na moduł samo: śledzi moduł okna w rdzeniu i przy jego zmianie rekonfiguruje pasek narzędzi promptu, panel akcji i panel kontekstu, zachowując wątek rozmowy. Powłoka może też przestawić okno wprost przez metodę ustawModul zamontowanej rozmowy, gdy zna wynik polecenia wejścia w obszar roboczy.

## budowa/klient-poprzedni/src/moduly/research/okno-export-panel.ts
Export Panel niesie dwie funkcje Operatora: eksport raportu oraz wybór formatu wyjściowego i miejsca docelowego. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności eksportu. Format wyjściowy idzie sterem nastawy: uchwyt niesie wartość bieżącą, a nie nazwę rodzajową, bo natywna lista pokazuje ją dopiero po rozwinięciu. Historia eksportów bieżącej sesji rośnie wyłącznie z odpowiedzi rdzenia, nie z zamówienia, więc wykaz mówi, co rdzeń naprawdę oddał, razem z brakiem ścieżki tam, gdzie jej nie oddał.

## budowa/klient-poprzedni/src/okna-rownolegle/wpis-sesji-okna.ts
Sekcja czynności menu gniazda rozstrzyga o pozycjach na podstawie wpisu sesji
— stanu, tytułu, projektu i liczby okien strumieniujących — a gniazdo zna
wyłącznie identyfikator sesji, więc to źródło prowadzi od identyfikatora do
pełnego wpisu. Źródło sesji w tle się tu nie nadaje, bo obsługuje sesje spoza
bieżącego połączenia, a nie tę, o którą pyta menu, mimo że korzysta z tej
samej komendy z żywym stanem. Źródło nie zna DOM-u, nie wykonuje żadnej
czynności i nie rozstrzyga, co z sesją wolno zrobić — oddaje wpis albo pustą
wartość, która oznacza brak wiedzy, nie zakaz: dopóki rdzeń nie oddał wykazu,
sekcja czynności pozostaje krótsza, bo czynności bez znanego stanu sesji nie
da się uczciwie nazwać.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-zrodla.ts
Tekst źródłowy mieszka w module stanu translate, żeby panele tłumaczenia patrzyły na ten sam tekst,
a nie na własną kopię. Wszystkie trzy komendy rdzeń rejestruje i obsługuje wprost; bez zalogowanego
modelu rdzeń oddaje stan zdegradowany z komunikatem o braku zalogowania — to jest odpowiedź rdzenia,
nie brak uchwytu. Ścieżka odmowy zostaje na wypadek starszego rdzenia albo pośrednika.

## budowa/klient-poprzedni/src/moduly/studio/strona-magazyn-widoku.ts
Moduł nastaw operatora widoku od początku przyjmował magazyn podany z zewnątrz
zamiast sięgać po magazyn przeglądarki wprost, żeby zapis dał się przełożyć na
rdzeń bez zmiany wołających go miejsc. Ten moduł jest tym przełożeniem: zapis
idzie do rdzenia komendami odczytu i zapisu widoku, więc nastawy przestają
ginąć przy przesiadce na inne urządzenie.

Magazyn nastaw widoku pozostaje z natury natychmiastowy: odczyt musi oddać
wartość w chwili składania okna, a rdzeń odpowiada później. Zapis dwustopniowy
rozwiązuje to bez kłamstwa — pamięć podręczna trzyma nastawy zapisane
w rdzeniu, a wczytanie ją napełnia; do chwili odpowiedzi obowiązuje odbicie
w magazynie przeglądarki, gdy jest dostępne, albo nastawy domyślne.

Odbicie w przeglądarce zostaje, ponieważ pierwsza klatka okna rysuje się przed
odpowiedzią rdzenia. Odbicie nie jest drugim źródłem prawdy: rdzeń nadpisuje
je przy każdym odczycie, a zapis idzie do obu naraz. Bez odbicia każde wejście
do modułu zaczynałoby się nastawami domyślnymi, choćby operator ustawił co
innego wcześniej.

Nieudany zapis nastawy widoku nie przerywa pracy, ale i nie milczy: idzie
zdaniem zwrotnym do zaplecza, więc korzystający wie, że wybór nie przeżyje
zamknięcia okna.

## budowa/klient-poprzedni/src/strona-glowna/karta-srodowiska.ts
Karta buduje treść w układzie: godło, tytuł krojem nagłówkowym, motto, jednozdaniowy opis trybu pracy krojem bazowym. W spoczynku powierzchnia jest neutralna, bez akcentu; przy najechaniu i w stanie czynnym pojawia się wstęga górna w błękicie sygnałowym, cień sygnału i uniesienie karty, z wariantu akcentu biblioteki komponentów, a arkusz strony wygasza wstęgę w spoczynku, bo biblioteka pokazuje ją stale. Sygnał jest jedyną barwą akcentu bieżącego systemu wizualnego. Nośnikiem karty jest przycisk, nie warstwa z rolą przycisku, ponieważ karta ma być celem nawigacji klawiaturą z pierwszeństwem natywnym, a pierścień fokusu wnosi wariant klikalny biblioteki komponentów.

## budowa/klient-poprzedni/src/moduly/translate/pasek-kontekstu.ts
Trzy pustki kontekstu są trzema różnymi zdaniami: jeszcze nie pytano, pytanie w toku, rdzeń nie zna
ani jednego okna tej sesji. Pierwsze każe czekać na wejście do modułu, drugie na odpowiedź, trzecie
mówi, że komendy panelu źródła nie mają czym zaadresować żądania. Przycisk ponowienia stoi tylko
tutaj, bo kontekst okna jest jedynym odczytem modułu; resztę wyzwala zapis, którego przycisk zostaje
klikalny także po odmowie. Ponowienie wyzwala odczyt i nie zdejmuje komunikatu — ten znika dopiero,
gdy odczyt się powiedzie. Tryb uprawnień, zasięg wykonania, rola okna i katalogi robocze stoją
w panelu sterowania okna na tym samym ekranie.

## budowa/klient-poprzedni/src/widok-sterowania/nazwy-ustawien.ts

Plik nie zna ani jednej nazwy prezentacyjnej z osobna — wszystkie bierze
z istniejących słowników etykiet okna, etykiet sterowania oraz rejestru
kanałów. Jedyne, co dokłada, to złożenie wartości w wiersz czytelny bez
otwierania kontrolki. Host wykonania nie jest osobnym wierszem, lecz dopiskiem
przy środowisku wykonania — uszczegóławia zasięg, a nie stanowi odrębnego
ustawienia. Nazwa hosta zostaje przez to widoczna bez dziewiątej pozycji.

Identyfikator kanału spoza wykazu nie znika i nie zostaje podmieniony — widok
pokazuje prawdę o oknie, a nie pozycję pierwszą z brzegu. Dopóki wykaz jest
pusty, bo odpowiedź rdzenia jest jeszcze w drodze, dopisek o pozycji spoza
wykazu byłby przedwczesny — wiersz pokazuje sam identyfikator ze wskaźnikiem
ładowania. Rejestr nie odróżnia braku odpowiedzi od pustego rejestru.

## budowa/klient-poprzedni/src/powloka/menu-profilu.ts

Menu profilu Operatora w pasku powłoki środowiska udostępnia siedem pozycji ustawień
pogrupowanych w trzy grupy tematyczne — tożsamość i dostęp, widok, obecność globalna —
oraz stopkę prowadzącą do pełnego rejestru ustawień. Pozycje nie są osobnym wykazem:
powstają z rejestru pozycji ustawień, z którego żyje także listwa strony głównej i menu
aplikacji. W tym module leży wyłącznie podział na grupy oraz umieszczenie pozycji
konfiguracji w stopce, bo stopka jest drogą do rejestru, czyli miejscem rzeczy
pełniejszej niż pozycje nad nią. Zmiana nazwy albo ikony w rejestrze zmienia to menu
samoczynnie. Mechanizm budowy menu pochodzi w całości z biblioteki komponentów drzewa
menu, która niesie grupowanie, nastawę dwustanową, stopkę, znacznik wyboru i wędrówkę
klawiaturą. Uchwyt menu niesie inicjały Operatora zamiast napisu rodzajowego —
pojedyncza litera zastępcza nie oznacza inicjału konkretnej osoby, tylko pierwszą literę
jedynej nazwy roli, jaką produkt zna. Motyw ma jedną prawdę: przełącznik w menu i
przycisk w pasku czytają ten sam stan obowiązującego motywu i nasłuchują tego samego
zdarzenia zmiany, więc zmiana w jednym miejscu przerysowuje drugie. Menu nie zawiera
konta, adresu ani polecenia wylogowania, bo platforma nie prowadzi kont ani profili
osobowych — bramka dostępu jest jedna, a rola Operatora pozostaje bezimienna, więc po
zalogowaniu nie ma bram do przekraczania. Zdanie graniczne tłumaczące ten brak stoi na
ekranie w stopce, tym samym brzmieniem co w menu aplikacji.

Pole szukania nad ośmioma pozycjami zabrałoby wiersz i nie skróciłoby ani jednego ruchu:
próg biblioteki, od którego pole szukania pojawia się samoczynnie, i tak by go nie
pokazał przy tej liczbie pozycji — zapis jawny mówi, że to świadomy wybór, a nie
przypadek. Znacznik wyboru pozycji w drzewie menu zostaje zawsze pusty, bo żadna
pozycja nie jest nastawą jednokrotną: wszystkie otwierają okno, więc menu nie ma czego
zaznaczać jako „ostatnio otwarte". Grupa „Widok" stoi między tożsamością a obecnością
i ma dokładnie jedną pozycję — nastawę motywu, a nie wejście do okna — dlatego składa
się osobno od pozostałych grup budowanych z rejestru.

## budowa/klient-poprzedni/src/moduly/research/okno-findings-panel.ts
Zdarzenie zmiany ustalenia jest już w kontrakcie, ale rdzeń go nie rozgłasza, więc wykaz nie dostanie ustalenia zapisanego na innym urządzeniu konta. Okno mówi o tym w stanie pustym zamiast odpytywać rdzeń w pętli — subskrypcja dojdzie razem z rozgłaszaniem. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności ustaleń.
Trzecia reprezentacja wykazu, kodowanie jakościowe, nie jest tu zbudowana: kody są już w kontrakcie, ale ustalenie oddawane przez rdzeń ich nie niesie i rdzeń nie ma uchwytu książki kodów, więc widok, który nie miałby czego pokazać, tu nie stoi — pozycja jest w panelu akcji, pod nazwą swojej komendy.
Kopia przed sortowaniem w funkcji porządkującej jest konieczna: wykaz pochodzi wprost z pamięci modułu, a sortowanie w miejscu wywróciłoby porządek narastania w pamięci, wspólnej wszystkim oknom.

## budowa/klient-poprzedni/src/moduly/translate/zrodlo-glosariusza.ts
Kontrakt nie ma komendy odczytu glosariusza jako wykazu, choć edycja glosariusza wymaga wczytania.
Źródło nie zmyśla wykazu: okno pokazuje wyłącznie terminy zapisane w bieżącej sesji i mówi wprost,
czego brakuje.

## budowa/klient-poprzedni/src/okna-rownolegle/zrodlo-lacznosci.ts
Doprowadzenie stanu łącza do sceny okien równoległych ma jedną odpowiedzialność:
zamianę subskrypcji transportu na strumień odczytów, które układ rozsyła do
nagłówków gniazd. Nic tu nie jest źródłem prawdy — prawdą jest sam transport,
a ten plik wyłącznie go odpytuje. Licznik kolejki musi być prawdziwy, a nie
zamrożony: transport ogłasza połączenie przed opróżnieniem kolejki, więc odczyt
zrobiony w chwili zmiany stanu zamarzałby na wartości sprzed wysłania, dlatego
po połączeniu odczyt dobija się cyklicznie, aż kolejka spadnie do zera. Tę samą
rachubę prowadzi niezależnie pasek górny aplikacji — nie jako druga prawda,
tylko jako drugi pytający tego samego transportu, ponieważ wspólnego miejsca
dla tej logiki nie ma: pasek zwraca gotowy element, nie strumień odczytów.
Interfejs transportu wystawia sześć podstawowych metod, a numer próby
i zaplanowane ponowienie trzyma prywatnie gniazdo połączenia, więc dojście do
przebiegu ponowienia sprawdza obecność metod jawnie i jednorazowo, zamiast
zakładać ją na sztywno w miejscu wywołania.

## budowa/klient-poprzedni/src/widok-sterowania/obserwator-ustawien.ts

Podsumowanie pokazuje osiem wartości obok siebie także wtedy, gdy szuflada
z kontrolkami jest zwinięta — a kontrolki kompletu sterowania trzymają swój
stan wewnątrz panelu sterowania i nie wystawiają go na zewnątrz. Ten plik jest
odczytem tego samego stanu, a nie drugą jego definicją: buduje stan sterowania,
zmianę okna i zmianę ustawienia z katalogu sterowania i nie dopisuje ani jednej
reguły protokołu.

Obserwator wyłącznie czyta — nie wywołuje czynności zastosowania ani zapisu,
więc nie ma drogi, którą mógłby zmienić okno. Komunikaty o losie zmian pomija,
bo pokazuje je pasek kompletu sterowania; podwójny komunikat byłby szumem. Gdy
panel sterowania wystawi własną migawkę i subskrypcję zmian, ten plik zniknie,
a podsumowanie odczyta stan kompletu wprost.

## budowa/klient-poprzedni/src/rozmowa/lista-wpisow.ts
Wpis rozpoznawany jest po kluczu, nie po pozycji: fragment strumienia odświeża tę samą pozycję, zamiast dokładać kolejną. Dzięki temu tura, która przyniosła prowenancję, dziesiątki fragmentów tekstu, wywołania narzędzi i podsumowanie, pozostaje w historii jednym wpisem. Przewijanie do końca następuje tylko wtedy, gdy lista już stała na końcu, więc czytanie starszej wypowiedzi nie jest przerywane przez nadchodzący strumień. Tryb widoku transkryptu jest stanem listy, nie pojedynczego wpisu: lista trzyma jeden tryb i rozsyła go do widoków, więc wpis założony po przełączeniu rodzi się już w trybie bieżącym. Pamięć wpisów jest od trybu niezależna, przełączenie niczego nie usuwa, a powrót do trybu zwykłego przywraca wątek w całości.

Po zdjęciu wszystkich wpisów stan pusty wraca na wierzch, ale nie zostaje sam: wołający zaraz po wyczyszczeniu dopisuje zdanie o powodzie, żeby w miejscu zniknięcia rozmowy nie stał napis o rozmowie, która się jeszcze nie zaczęła.

Dwie różne pustki dostają dwa różne zdania: to, że wątek się nie zaczął, i to, że tryb Streszczenie nie ma jeszcze czego streścić, to nie ten sam fakt. Drugi wariant tłumaczy, czym streszczenie jest i skąd się bierze.

## budowa/klient-poprzedni/src/moduly/workspace/agenci-czynnosci.ts
Czynności okna zarządcy ekspertów sięgające poza sam wykaz — zmiana uprawnienia eksperta i zdanie
o przypisaniach projektu — stoją w osobnym pliku od okna, bo okno składa kontrolki i prowadzi odczyt,
a tu leży przebieg czynności wraz z tym, co po niej widzi operator, tym samym wzorem co w module
czynności biblioteki. Potwierdzenie mówi to, co zapisał rdzeń, nie to, czego żądało okno: rdzeń
oddaje po zapisie komplet uprawnień eksperta i ten komplet rozstrzyga. Po odpowiedzi okno odrysowuje
wykaz, bo stan ładowania go opróżnia.

## budowa/klient-poprzedni/src/strona-glowna/listwa-ustawien.ts
Segmenty listwy dzieli delikatny separator wewnątrz jednej powierzchni, więc pas czyta się jako jeden byt o kilku wejściach, a nie jako trzecia siatka kart; pozycje listwy nie mają formy kart ani kafli. Segmenty ustawień są trzy: okno konfiguracji, tryb Mobile oraz Always On Display. Czwarty segment nie jest ustawieniem — „Dodaj nowy” otwiera formularz zakładania komponentu własnego i zgłasza to osobnym wywołaniem zwrotnym, nie przez wykaz ustawień, ponieważ pozycja, która ustawieniem nie jest, nie udaje jego kodu; zakładanie komponentu należy do strefy drugiej, więc przeniesienie segmentu jest zgłoszone, a nie wykonane z tego katalogu. Waga wizualna strefy jest najniższa z trzech: segment ma wysokość kontrolki i niesie ikonę oraz nazwę, bez wezwania do działania i bez metadanych. Żadna pozycja nie jest wyszarzona ani pozbawiona klikalności.

## budowa/klient-poprzedni/src/moduly/studio/strona-panel-tresci.ts
Pięć czynności panelu — odczyt treści fragmentu, zmiana brzmienia, odczyt
postaci, zapis postaci i wykaz pochodzenia — dotyczą jednego: treści i postaci
dokumentu jako całości, a nie pojedynczej cechy; rozsypane po osobnych panelach
nastaw byłyby nie do znalezienia. Zapis postaci przyjmuje treść jako pole
opcjonalne kontraktu: brak znaczy brak zmiany treści. Przycisk zapisu odczytuje
postać i zapisuje ją z powrotem, zakładając wersję — utrwala to, co nastawy
strony i style zmieniły, bez dotykania treści, żeby postać przestawała ginąć
nawet wtedy, gdy operator nie pisał.

Brzmienie fragmentu ma pole wielowierszowe, ponieważ poprawa fragmentu bez
przepisywania całości jest osią zamówienia, a fragment bywa akapitem, nie
wyrazem; pole jednowierszowe wymuszałoby wklejanie akapitu w linijkę wysokości
jednego wiersza.

## budowa/klient-poprzedni/src/moduly/research/okno-reading-view.ts
Okno jest zbudowane, choć katalog okien rdzenia jeszcze go nie zna, bo dwie jego czynności mają dziś pokrycie: wczytanie treści dokumentu repozytorium i zamiana zaznaczonego fragmentu w ustalenie z powiązaniem do czytanego źródła. Podświetlenia trwałe, notatki na marginesie i wypisy zbiorcze mają już w kontrakcie własne komendy i własny byt wraz z kotwicą pozycji, ale rdzeń nie ma dla nich uchwytu — okno ich nie udaje, zaznaczenie jest zaznaczeniem przeglądarki, a trwałym staje się dopiero jako ustalenie. Pozostałe operacje na źródle stoją w panelu akcji pod nazwami swoich komend i wracają odmową rdzenia, zamiast znikać z okna. Plik składa widok; zachowanie po naciśnięciu leży w module obsługującym czynności lektury.

## budowa/klient-poprzedni/src/moduly/workspace/okno-rozmowy.ts
Dwie kontrolki modułu wołają przeniesienie kontekstu i obie wymagają identyfikatora okna źródłowego:
udostępnienie zaznaczonych zasobów w bibliotece projektu oraz przejście do budowniczego eksperta
w zarządcy ekspertów. Rdzeń niesie kod modułu w opisie okna, więc okna karty sesji wystarczy zawęzić
do modułu. Moduł montuje się od razu, a nie odracza montażu do chwili znalezienia okna, bo okno
rozmowy jest mu potrzebne do dwóch czynności, nie do istnienia — jego brak nie wygasza pięciu okien
operacyjnych. Kolejność względem wejścia do modułu jest treścią: okno rozmowy jest wtedy już
przestawione na ten moduł, a pytanie zadane wcześniej oddałoby okno modułu poprzedniego. Funkcja nie
zgaduje: gdy rdzeń odmówi albo nie odda okna tego modułu, wraca pusty identyfikator wraz z powodem,
który staje potem w odmowie obu kontrolek. Rdzeń oddaje w wykazie oba stany okna, otwarty i zamknięty.
