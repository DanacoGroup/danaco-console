# Danaco Console — Uzasadnienia komentarzy klienta poprzedniego

Dokument gromadzi uzasadnienia, które przekraczają dopuszczalną długość
nagłówka komentarza w plikach `budowa/klient-poprzedni/`. Każdy rozdział nosi
nazwę pliku źródłowego, którego uzasadnienie dotyczy.

## budowa/klient-poprzedni/src/moduly/design/czynnosci-warsztatow-designu.ts

Katalog opisuje pięć warsztatów danymi zamiast pięciu odrębnych zestawów
formularzy, wzorem pliku `studio/czynnosci-warsztatu.ts` — okno buduje pola
z tego wykazu i składa żądanie wspólną funkcją `zloz`. Druga, osobna rodzina
typów rozjechałaby się z pierwszą przy każdym nowym rodzaju pola, co
wymuszałoby poprawki w dwóch oknach zamiast w jednym.

Nazwy pól są nazwami kontraktu, ponieważ trafiają wprost do rdzenia. Etykiety
są zdaniem Operatora i z nazwami pól się nie pokrywają — Operator czyta
„proporcje kadru", a rdzeń otrzymuje `aspectRatio`.

Nazwa komendy pochodzi wyłącznie ze stałych kontraktu (`Command.*`), nigdy
z napisu wpisanego ręcznie: napis literowany z pamięci przechodzi sprawdzian
typów, lecz zawodzi dopiero u Operatora, gdy rdzeń takiej komendy nie zna.

Podpowiedź pola wykazu czynności bierze nazwę komendy z tych samych stałych
kontraktu z tego samego powodu — napis wpisany wprost przeżyłby zmianę nazwy
w kontrakcie i podpowiadałby Operatorowi komendę, której rdzeń już nie zna,
a to podpowiedź jest tym, co Operator przepisuje do pola.

## budowa/klient-poprzedni/src/moduly/design/katalog-funkcji-designu.ts

Katalog istnieje po to, żeby stan modułu Design dało się przeczytać, a nie
zgadnąć. Okno pokazuje czynności, które wykonuje; katalog pokazuje komplet
zamierzonych pozycji i przy każdej mówi jedno z trojga: którą komendą
kontraktu jest wykonywana, że wykonuje ją samo okno bez udziału rdzenia, albo
czego brakuje, żeby była. Pozycja bez drogi nie znika z wykazu, ponieważ
zniknięcie byłoby ukryciem braku.

Trzeci stan jest konieczny, bo dwa nie oddają prawdy o tym module. Kanwa,
wyrównanie warstw, siatka pomocnicza, drzewo żetonów i rachunek kontrastu
dzieją się w całości w przeglądarce i komend nie potrzebują; nazwanie ich
brakiem kontraktu byłoby zmyśleniem długu, a nazwanie komendą — zmyśleniem
drogi wykonania, której nie ma.

Nazwy pozycji i podział na grupy nie są tłumaczone ani parafrazowane w żadnym
miejscu, które je wyświetla — pozycję odnajduje się po pełnej nazwie własnej.

Liczba pozycji katalogu nie jest zapisana nigdzie na stałe jako osobna
wartość, tylko liczona z wykazu funkcją length w miejscu użycia, żeby napis
nie mógł rozjechać się z rzeczywistą zawartością wykazu, który stoi obok.

Katalog jest zbiorem danych, nie widokiem — wyszukiwarka funkcji modułu
Design buduje z niego listę do wyświetlenia.

## budowa/klient-poprzedni/src/moduly/apps/narzedzia-apps.ts

Narzędzia są pogrupowane wedle okna modułu Apps, w którym dana czynność się
odbywa: Product Builder prowadzi produkt, etapy, kamienie milowe i oś czasu;
Architecture Designer — walidację układu, wersje, adnotacje i eksport;
warsztaty — podgląd, trasy, motyw, punkty końcowe i schemat; Deployment Panel
— środowiska, zmienne, domenę, skalowanie, kondycję, dzienniki i artefakty;
Publisher Panel — pakowanie, manifest, walidację, podpis i publikację.

Zdanie o skutku jest w każdym narzędziu inne, bo skutek jest inny: liczba
pozycji, nazwa bytu, kod odpowiedzi, odwołanie do pliku, wynik weryfikacji
podpisu. Wspólny komunikat „gotowe" byłby meldunkiem, z którego nic nie
wynika, a to jest dokładnie wzorzec, którego moduł ma nie powtarzać.

Wartości startowe pól są przykładami z domeny produktu, nie wartościami
wymuszonymi: Operator zmienia je przed wykonaniem, a puste pole nieobowiązkowe
oznacza brak zawężenia.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-apps.ts

Obszar apps.* jest dwukierunkowy: obok trzech komend zapisu (architektura,
plik warsztatu, wdrożenie) stoją trzy komendy odczytu — apps.deployment.list,
apps.architecture.get, apps.workspace.list. Bez nich moduł po odświeżeniu
okna przeglądarki zaczynałby od zera, bo wypełniałyby go tylko zdarzenia
bieżącej sesji gniazda, a historia wdrożeń trzymana w bazie znikałaby z oczu.
Każdy odczyt oddaje ten sam kształt, którym odpowiada jego komenda zapisu.

Zdarzenie apps.workspace.changed jest drogą rozgłoszenia dla
apps.workspace.update: bez niego plik zapisany w jednym oknie nie docierałby
do drugiego, dopóki plik nie zostałby odczytany ręcznie. Wchodzi tą samą
bramą co apps.build.changed — subskrypcją źródła, nie własnym gniazdem.

Źródło nie ma własnego stanu: jest warstwą wywołań i sprawdzianu kształtu
odpowiedzi. Stan produktu mieszka osobno, żeby pięć okien patrzyło na jeden
zbiór, a nie na pięć kopii.

Komendy zapisu mają uchwyt w rdzeniu, więc wywołanie wraca zwykłą
odpowiedzią albo odmową merytoryczną, nie kopertą apps.unknown. Droga przez
warstwę odmowy rdzenia zabezpiecza ścieżkę fail-open: gdyby uchwyt zniknął
albo rdzeń nie rozpoznał którejś komendy, obietnica wywołania ma się czym
rozstrzygnąć zamiast wisieć bez końca, a okno nazywa odmowę zamiast ją ukryć.

Pola tożsamości idą do rdzenia tak, jak je wpisano, bez przycięcia po drodze
po stronie klienta. Przycięcie w przeglądarce i przycięcie w rdzeniu nie są
tym samym działaniem i rozjeżdżają się na dwóch znakach białych, których
jedna strona zdejmuje, a druga zostawia. Ma to znaczenie na ścieżce
warsztatu, bo klucz złożony z okna, warstwy i ścieżki czyni ją tożsamością
pliku: przycięcie po stronie klienta zapisałoby plik pod ścieżką inną niż
wpisana albo odrzuciłoby ścieżkę, którą rdzeń przyjmuje. Rdzeń przycina te
pola po swojemu i to on jest tu jedyną władzą; o pustce pola rozstrzyga
więc pustka dosłowna, a okno zestawia potem wpisane z oddanym i ogłasza
różnicę.

Sprawdzian kształtu odpowiedzi jest przy każdej komendzie osobny i pyta
o pole, którego okno naprawdę używa. Sprawdzanie, czy cokolwiek wróciło,
przepuściłoby odpowiedź o kształcie innym niż kontraktowy, a okno wywróciłoby
się dopiero przy rysowaniu. Pustka nie jest tu uszkodzonym kształtem: wykaz
pusty i pole opcjonalne bez wartości są odpowiedziami prawdziwymi i znaczą,
że w tym oknie danej rzeczy jeszcze nie ma — dlatego przy nich sprawdzian
pyta o tablicę albo przepuszcza wszystko, zamiast żądać obiektu.

## budowa/klient-poprzedni/src/moduly/apps/narzedzia-rozszerzen.ts

Narzędzia są pogrupowane wedle okna strony rozszerzeń: App Catalog prowadzi
wyszukiwarkę, kartę szczegółów, kolekcje i rejestr organizacji; Installed
Apps Manager — aktualizacje, przesyłkę paczki, wersjonowanie, instalację
zestawu, dziennik cyklu życia i tryb administracyjny; Integrations Hub i MCP
and Connector Console — transport, poświadczenie, odkrywanie narzędzi, próbne
wywołanie, log protokołu, piaskownicę, import definicji, webhooki,
odwzorowania, metryki i kondycję; Permissions and Trust Center — uprawnienia,
podpis, skaner i sekrety.

Żadne pole nie niesie treści poświadczenia. Pola sekretów przyjmują wyłącznie
klucz jawny warstwy sekretów — nazwę, po której rdzeń wydaje wartość.
Wprowadzenie hasła, tokenu czy klucza API zostaje po stronie Operatora, poza
tą drogą i poza kontraktem.

Pozycja, na której narzędzie ma pracować, bierze się ze wskazania w oknie
funkcją stan.wybrane(), a nie z pola wpisywanego przy każdym przycisku. Brak
wskazania kończy się odmową z powodem — Operator ma najpierw wybrać pozycję
z listy.

## budowa/klient-poprzedni/src/komponenty/menu-drzewo.ts

Mechanizm stoi w bibliotece komponentów, a nie wewnątrz paska zlecenia, bo
pasek jest tylko jego największym odbiorcą: takie samo menu stawiają
nagłówki okien, panele i ekrany konfiguracji. Mechanizm zamknięty w pasku
byłby dla nich nieosiągalny i zostałby skopiowany, a kopie rozjeżdżają się
w szczegółach — haczyk wyboru pojawiałby się w jednym menu, a w drugim nie,
przy jednakowym wymogu.

To jest mechanizm docelowy dla sterów nastawy: dopisanie obok niego kolejnego
mechanizmu rozwijania odtwarza dokładnie ten rozjazd, który on likwiduje.
Trwałym wyjątkiem zostaje mechanizm menu czynności z osobnym plikiem, bo
obsługuje menu, którego uchwyt jest ikoną, a treść wchodzi jako gotowe
elementy, a nie ster niosący bieżącą wartość — to osobny wzorzec, nie rozjazd
do scalenia.

Mechanizm niesie sam sześć cech, których odbiorcy nie budują u siebie:
zagnieżdżenie gałęzi bez ograniczenia głębokości; znacznik bieżącego wyboru
na liściu i ślad tego wyboru na gałęzi, która ten liść niesie, żeby Operator
widział wybór bez wchodzenia w gałąź; opis jako pole każdej pozycji;
przełącznik dwustanowy wewnątrz gałęzi; grupowanie pozycji po źródle lub
rodzinie; drogę do rejestru na dole menu. Treść wchodzi jako dane, a formę
nadaje mechanizm; przyjmowanie gotowych elementów od odbiorcy zniosłoby
jednolitość tych sześciu cech.

Mechanizm nie wie, co jest w drzewie, i nie wykonuje wyboru — oddaje klucz
pozycji wołającemu; nie zna kontraktu, komendy ani stanu okna. Żadna pozycja
nie dostaje stanu wyłączonego: drzewo puste nie jest błędem, uchwyt otwiera
się i mówi zdaniem, że wykaz jest pusty, zamiast przestać reagować.

Pozycja niesie równocześnie nazwę pełną i nazwę skróconą. W wykazie płaskim
na setki pozycji nazwa pełna musi nieść źródło, bo inaczej dwie komendy
o tej samej nazwie z dwóch wtyczek są nie do rozróżnienia. To samo źródło
czyni jednak wykaz nieczytelnym, gdy sto pozycji zaczyna się tym samym
przedrostkiem — dlatego pozycja niesie obie nazwy, pełną ze źródłem i
skróconą do samej rzeczy. Szukanie obejmuje obie nazwy i opis; trafienie
w nazwę skróconą liczy się wyżej niż w pełną, bo to ona jest tym, czego
Operator naprawdę szukał.

Opis pozycji stoi wprost pod nazwą, a nie za dymkiem wywoływanym najechaniem,
bo w menu o kilkudziesięciu liściach najeżdżanie na każdy liść z osobna nie
jest drogą; dymek zostaje formą dla kontrolek stojących pojedynczo.

Nazwa nastawy (na przykład „Model", „Wysiłek", „Urządzenie") nie trafia na
ekran, bo na uchwycie stoi wartość, nie nazwa nastawy. Idzie do etykiety
dostępności uchwytu i listy, żeby czytnik ekranu wiedział, czego dotyczy
wartość, której nazwa sama tego nie mówi — „Opus 5" nie niesie słowa „model".

Próg liczby liści, od którego drzewo stawia pole szukania, liczy liście, nie
wszystkie pozycje: drzewo o trzech gałęziach i dwóch liściach filtra nie
potrzebuje, drzewo o dwóch gałęziach i stu pozycjach potrzebuje go
natychmiast. Próg, a nie stałe pole szukania, bo pole nad wykazem pięciu
pozycji zabierałoby wiersz i nie skracałoby żadnego ruchu.

Kierunek rozwinięcia wykazu nie jest preferencją wizualną. Ster stojący
w pasku u góry okna ma pod sobą całą wysokość sceny i rozwija się w dół.
Ster przy polu wpisywania stoi u dołu okna, a wykaz rozwinięty w dół nie
miałby dokąd pójść — wyszedłby poza krawędź sceny.

Tryb bez uchwytu obsługuje wykaz komend wywoływany ukośnikiem: wykaz ma się
pojawić natychmiast po wpisaniu znaku, a filtrem ma być to samo pole
wpisywania, nie osobne okienko. W tym trybie uchwyt nie wchodzi do
dokumentu, wewnętrzne pole szukania nie powstaje wcale, a wołający steruje
mechanizmem przez rozwinięcie, ustawienie frazy, przesunięcie wyróżnienia
i wybór pozycji wyróżnionej. Ognisko zostaje w polu wołającego — stąd
wyróżnienie jest wirtualne, a nie ogniskiem.

Opis pokazywany tylko przy pozycji wyróżnionej to ustawienie dla wykazu na
setki pozycji: setka opisów naraz nie jest objaśnieniem, tylko ścianą
tekstu, przez którą nie widać już samych nazw — wtedy opis należy się jednej
pozycji, tej, na którą Operator właśnie patrzy.

Metoda ustawiająca całą zawartość naraz podaje wartość widoczną na uchwycie
i całe drzewo. Drzewo przerysowuje się w całości, bo wykaz bywa czytany
z rdzenia i pozycje przychodzą oraz znikają; rozwinięte gałęzie przeżywają
przerysowanie po kluczu, żeby odświeżenie katalogu nie zwijało menu pod ręką
Operatora.

Fraza filtrująca podana z zewnątrz wchodzi do wewnętrznego pola, gdy pole
stoi na ekranie, żeby Operator widział, czym wykaz jest przycięty. Gdy pole
nie stoi (tryb bez uchwytu albo wykaz poniżej progu szukania), fraza i tak
przycina wykaz — próg rozstrzyga, czy mechanizm stawia pole, a nie czy
w ogóle umie filtrować.

Przesunięcie wyróżnienia nie rusza ogniska: przy obsadzie ukośnikiem
strzałki naciska się w polu wpisywania, więc ognisko musi w nim zostać;
przenoszenie go na pozycję wyrwałoby Operatorowi klawiaturę spod palców
w połowie pisania.

Wybór pozycji wyróżnionej oddaje wartość fałsz, gdy nie ma czego wybrać —
wykaz pusty albo przycięty do zera — i wtedy wołający wie, że klawisz Enter
ma zrobić swoje zwykłe zadanie zamiast niczego. Wyróżniona gałąź nie jest
wyborem — Enter na niej otwiera poziom.

## budowa/klient-poprzedni/src/moduly/browser/sekcje-rodzin.ts

Plik składa wyłącznie wywołania rdzenia w zdania o skutku — elementy stawia
`panel-rodzin.ts`, a stan modułu trzyma `stan-przegladania.ts`. Zdanie
zwracane po czynności niesie zawsze liczbę albo identyfikator wzięty wprost
z odpowiedzi rdzenia, nigdy samo potwierdzenie w rodzaju „gotowe" — zdanie
bez pokrycia w odpowiedzi byłoby meldunkiem bez skutku, którego Operator nie
mógłby sprawdzić.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-designu.ts

Zbiór zasobów mieszka w `stan-designu.ts`, żeby trzy okna patrzyły na jeden
zbiór, a nie na osobne kopie.

Zasób wchodzi do modułu dwiema drogami, które się nie zastępują: `design.asset.generate`
daje zasób z pracy modelu i wymaga kanału obrazowego, a jego brak jest odmową
nazywającą brak, nigdy obrazem zastępczym; `design.asset.upload` daje zasób
z pliku wskazanego przez Operatora i działa bez żadnego kanału modelu, wnosząc
treść już istniejącą.

Przy wgraniu Operator wskazuje plik, klient czyta go i oddaje treść w polu
`contentBase64`, a rdzeń zapisuje ją w swoim magazynie pod sumą kontrolną. Od
tej chwili zasób nie zależy już od pliku na dysku, który Operator może
nadpisać albo skasować.

Treść idzie bajtami, nie ścieżką, mimo że kontrakt zna oba pola. `sourcePath`
każe rdzeniowi otworzyć plik pod wskazaną ścieżką, a rdzeń biegnie na innej
maszynie niż przeglądarka Operatora: ścieżka z jego pulpitu nie znaczy tam nic
albo znaczy coś zupełnie innego. Klient czyta plik sam i wysyła to, co
przeczytał.

Odmowa wraca polem `blad`, nie wyjątkiem: rdzeń odmawia czynności na zasobie,
którego nie zna, a okno pokazuje to zdanie i zostaje czynne.

Wskazanie kanału zlecenia generowania jedzie polem `channelId`: rdzeń sprawdza
po nim rodzaj kanału i odmawia wskazania tekstowego.

Pola opisowe zlecenia wgrania — nazwa, format, wymiary — idą tylko wtedy, gdy
klient je zmierzył; podstawienie wartości domyślnej byłoby wpisaniem rdzeniowi
metadanych, których nikt nie sprawdził.

Skala i jakość zlecenia wydania są w kontrakcie opcjonalne: zero wysłane do
rdzenia byłoby żądaniem obrazu o zerowym boku albo zerowym stopniu kompresji
i wróciłoby odmową, dlatego pole zerowe nie wchodzi do żądania.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-rozszerzen-dobudowa.ts

Podział na dwa źródła (`zrodlo-rozszerzen-apps.ts` i to źródło) jest podziałem
obowiązków, nie rozmiaru. Pierwsze prowadzi cykl życia pozycji: wykaz,
instalację, konfigurację, przełącznik i odinstalowanie. To źródło niesie
wszystko, co robi się na pozycji już stojącej — wyszukiwanie, kolekcje,
wersjonowanie, rozmowę protokołem, webhooki, uprawnienia, podpis i sekrety.

Żadna z tych komend nie niesie identyfikatora okna: katalog rozszerzeń jest
bytem rdzenia stojącym poziom wyżej niż moduł, Apps jest jego operacyjnym
frontem, nie drugim źródłem prawdy, więc pozycji nie posiada żadne okno.

Sprawdzian kształtu odpowiedzi pyta przy każdej komendzie o pole, którego okno
naprawdę używa. Pustka nie jest tu uszkodzonym kształtem: wykaz pusty i pole
opcjonalne bez wartości są odpowiedziami prawdziwymi i znaczą, że tego jeszcze
nie ma.

## budowa/klient-poprzedni/src/moduly/automations/zrodlo-dobudowy.ts

Rozdział między tym plikiem a `zrodlo-automations.ts` idzie po roli, nie po
objętości. Tamten plik niesie rdzeń modułu — definicję, harmonogram, kolejkę,
zależności i przebiegi, czyli to, czym okna pracują bez otwierania
jakiegokolwiek panelu. Ten plik niesie czynności paneli i szuflad.

Każda czynność oddaje `Wynik` z całą odpowiedzią, nie z wyciętym polem, bo
większość tych odpowiedzi niesie więcej niż jedną rzecz naraz — wykaz i
znacznik przycięcia, automatykę i numer wersji, zlecenie i znacznik
duplikatu, zmienne i zastrzeżenia. Wycięcie jednego pola gubiłoby drugie,
a okno musi pokazać oba.

Sprawdzian kształtu pilnuje pola obowiązkowego odpowiedzi: rdzeń, który oddał
kopertę powodzenia bez treści, jest dla okna odmową, inaczej widok rysowałby
pustkę i twierdził, że to wynik.
