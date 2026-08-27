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
