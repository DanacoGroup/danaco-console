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

## budowa/klient-poprzedni/src/moduly/automations/okno-queue-manager.ts

Okno nie ma własnego silnika kolejek ani własnego cyklu życia zlecenia: jeden
silnik obsługuje pętlę sesyjną i MultitaskingAI, a okno posuwa kolejkę
komendami kontraktu i pokazuje to, co rdzeń o niej oddaje.

Czego kontrakt nie oddaje, tego okno nie zmyśla. Byt `Queue` niesie stan,
licznik obiegów, okna, politykę, liczbę zleceń oczekujących i powiązaną
automatykę, lecz nie niesie samych zleceń — te oddaje osobna komenda wykazu
zleceń. Dopóki rdzeń nie ma jej uchwytu, przegląd pokazuje kolejkę jako
całość, a identyfikator zlecenia wpisuje Operator; okno mówi to wprost,
zamiast pokazywać pustą listę udającą przegląd.

O tym, czy rdzeń ma uchwyt danej komendy, okno samo nie orzeka i nie
przepisuje stanu kontraktu do zdań: powód pozycji paska bierze się z wykazu
komend oddanego przez rdzeń (`pokrycie-komend.ts`), a wykaz działań silnika
z wyliczenia kontraktu. Oba przestają być prawdziwe same, bez edycji tego
pliku.

Zdanie potwierdzenia powstaje z oddanej kolejki, nie z etykiety przycisku:
etykieta mówi o żądaniu, a rdzeń może oddać kolejkę w innym stanie, niż
żądanie zapowiadało.

Odpowiedź zasilenia kolejki krokami niesie stan kolejki i liczbę zleceń
oczekujących, lecz nie mówi, które zlecenia w niej stoją — zdanie po zasileniu
mówi więc to, co odpowiedź naprawdę niesie, i nie orzeka, że kroki weszły.

Priorytet zlecenia jedzie przy działaniu „wstrzymaj”, bo jest polem żądania,
nie osobnym działaniem silnika, a wstrzymanie jest działaniem odwracalnym
jednym naciśnięciem „Wznów”.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-zarzadu.ts

Zarząd repozytorium jest wydzielony ze `zrodlo-biblioteki.ts` wzdłuż
odpowiedzialności: tamten plik prowadzi wykaz i podgląd, ten niesie czynności
zarządcze warstw trzeciej i czwartej modułu biblioteki (opis zasobu, słownik
etykiet, tezaurus, reguły, higienę, cykl życia, utrwalenie, udostępnienia
i sugestie). Kształt pozostaje jeden — `ZrodloBiblioteki` — więc okna nadal
widzą jedno źródło.

## budowa/klient-poprzedni/src/moduly/automations/okno-scheduler.ts

Cykliczność w oknie Schedulera ustala się wzorcem, nie samą składnią: Operator wybiera
opis w rodzaju „co tydzień, poniedziałek, 07:00”, a okno składa z tego zapis cron według
kontraktu wzorców z `wzorce-cyklicznosci.ts`. Pole zapisu cron pozostaje przy tym widoczne
i edytowalne — zapis wpisany wprost jest rozpoznawany z powrotem jako wzorzec, a zapis,
którego żaden wzorzec nie obejmuje, jest wzorcem własnym i idzie do rdzenia w całości, bez
udziału kreatora.

Harmonogram obowiązuje dopiero po powiązaniu z automatyką, dlatego okno wprost nazywa ten
warunek zamiast odmawiać zapisu bez wyjaśnienia przyczyny.

Potwierdzenie po zapisie opiera się na polach `enabled` i `nextRunAt` odpowiedzi rdzenia,
nie na treści żądania: gdyby zdania „wstrzymany” i „wznowiony” zależały od tego, o co okno
prosiło, rdzeń, który zapisu nie przyjął zgodnie z żądaniem Operatora, dostawałby od okna
potwierdzenie czynności, która się nie odbyła.

Pola kreatora cykliczności stoją w rzędzie zawsze, niezależnie od wybranego wzorca, zamiast
chować pola nieużywane przez dany wzorzec: pole schowane przy przełączeniu wzorca
przeskakiwałoby układ pod ręką Operatora, a pole nieużywane jest nieszkodliwe, bo jego
wartość po prostu nie wchodzi do zapisu.

Kalendarz uruchomień czyta przebiegi bez wskazania okna, więc rdzeń nie zakłada na oknie
Schedulera obserwacji telemetrii — kalendarz jest jednorazowym zdjęciem stanu, a stała
obserwacja przebiegów należy do Execution Monitora i ma pozostać jej jedynym odbiorcą
w module Automations.

## budowa/klient-poprzedni/src/moduly/apps/etykiety-apps.ts

Teksty tego pliku stoją osobno od plików budujących elementy okien modułu Apps, na wzór
podziału zastosowanego w `sterowanie/etykiety-sterowania.ts`. Pozycja wykazu braków niesie
tylko to, co okno wie samo: nazwę czynności, czego by do niej trzeba (pole `czego`), nazwę
komendy, której wejście do kontraktu znosi brak (pole `komendaZnoszaca` — pozycja znika
wtedy z wykazu sama) oraz cudze drogi do sprawdzenia w wykazie (pole `komendyCudze`).
Zdanie o stanie kontraktu dokłada `braki-kontraktu.ts`, czytając stałą `KOMENDY` przy
składaniu okna.

Nazwy komend proponowanych w tym pliku nie są nazwami kontraktu i nie idą na drut do
rdzenia: pozycja z takim wskazaniem stoi w wykazie braków, a wpisanie komendy do kontraktu
zdejmuje pozycję z wykazu samo. Nazwy stoją w jednym miejscu, żeby definicja wpisana
do kontraktu i zdanie widoczne Operatorowi mówiły o tej samej komendzie. Obszar
nazwy idzie za bytem, którego dotyczy: czynności na katalogu rozszerzeń należą do obszaru
`extension`, bo katalog stoi poziom wyżej niż moduł, a czynności na produkcie budowanym
w module — do obszaru `apps`.

## budowa/klient-poprzedni/src/moduly/apps/okno-integrations-hub.ts

Okno pokazuje dwa rodzaje pozycji katalogu — serwer protokołu MCP i integrację
przez interfejs API — bo tylko te dwa są drogą do usługi zewnętrznej. Wtyczki
i umiejętności nie należą do tego okna: są rozszerzeniami powłoki i sposobami
wykonania zadań, nie drogami do usług zewnętrznych, więc ich miejsce jest
w katalogu aplikacji.

Test połączenia i monitor zdrowia sprawdzają most, przez który dochodzi się do
serwera, a nie usługę stojącą za mostem. Sprawdzenie oddaje stan, czas
i korzenie potwierdzone przez punkt dostępu; czasu odpowiedzi samej integracji
ani liczby jej narzędzi kontrakt nie niesie, więc okno tej wartości nie zmyśla.

Integracja bez punktu dostępu nie jest usterką: pozycja rodzaju API bywa
dostępna wprost adresem, a punkt dostępu opisuje most do maszyny. Okno
odróżnia brak mostu od mostu nieznanego rdzeniowi, ponieważ to dwa różne stany
i drugi z nich jest niespójnością rejestru.

Transport protokołu MCP da się dziś wpisać wyłącznie do nieprzezroczystej
konfiguracji w postaci zapisu JSON. Kontrakt nie nazywa tego pola osobno, więc
rdzeń i okno mogą rozumieć wpisaną wartość odmiennie.

## budowa/klient-poprzedni/src/asystent-plywajacy/stan-dymka.ts

Dymek nie ma własnej warstwy wywołań — bierze `zrodlo-assistant.ts`
i `zrodlo-zaplecza.ts` modułu Assistant, żeby własna warstwa nie rozjechała się
z modułem przy zmianie kształtu odpowiedzi. Żadna ścieżka nie kończy się
milczeniem: ustalenie okna (`window.list`), wydanie polecenia
(`assistant.voice.command`), odczyt dziennika (`assistant.activity.list`)
oraz zamknięcie zlecenia błędem albo anulowaniem zawsze kończy się wypowiedzią
w historii, nigdy cichym `return` — pusta historia czytałaby się jak brak
głosu asystenta, co byłoby nieprawdą. Stan okna rozróżnia sześć wartości,
nie dwie, bo odmowa rdzenia i brak okna asystenta prowadzą do różnych
wniosków Operatora.

Odpowiedź asystenta przychodzi później niż odpowiedź komendy: rdzeń potwierdza
samo przyjęcie zlecenia, a wpis rodzaju `result` dopisuje przy domykaniu
(`adapter_modul_asystent_wykonawca.go` → `domknijZlecenie`), więc zejście
zlecenia z toru, ogłaszane zdarzeniem `assistant.action.changed`, pociąga
odczyt dziennika.

Źródło posunięć jest jedno i wspólne z pasem dolnym (`aplikacja/
pas-posuniec.ts`); druga subskrypcja byłaby drugim rozstrzyganiem sprawcy,
a `utworzRozstrzyganieSprawcy` zużywa odcisk okna przy weryfikacji, więc dwa
egzemplarze wydałyby dwa różne werdykty o jednym zdarzeniu.

Sprawca idzie z koperty zdarzenia: `AssistantActionChangedEvent` niesie
`actor` i `actorClientId` wypełniane przez rdzeń (`core/sprawca.go`).
Rozgłoszenie idzie do całego konta (`transport/rozgloszenie.go`), więc
zlecenie założone poza tym dymkiem wchodzi jako posunięcie (kto, co, w jakim
stanie), ale bez treści dziennika — treść odpowiedzi należy do okna,
w którym padło polecenie.

Odmowa `odnotujOdmowe` trafia do historii oprócz dymka powiadomienia, bo
powiadomienie znika po sekundach, a Operator wraca po przebieg rozmowy do
historii i ma tam znaleźć ślad każdej odmowy.

## budowa/klient-poprzedni/src/dostepy/dodanie-katalogu.ts

Drogą pierwszą i właściwą jest natywne okno wyboru powłoki. Ścieżka wraca z niego
istniejąca i rozwinięta, więc punkt zakładany na niej nie wskazuje miejsca, którego
nie ma.

Drogą drugą jest pole ścieżki wpisywanej z ręki. Nie jest ono atrapą ani zapasem na
gorsze czasy: interfejs bywa otwarty w przeglądarce, bez powłoki, a wtedy natywnego
okna po prostu nie ma. Pole zostaje czynne zawsze, także w powłoce, ponieważ Operator
znający ścieżkę nie musi jej odklikiwać.

Przycisk natywnego wyboru nie jest wygaszany poza powłoką. Naciśnięcie daje odpowiedź:
zdanie o tym, że okno systemowe należy do powłoki, oraz przeniesienie uwagi do pola
ścieżki.

Katalog zakładany tą drogą leży na maszynie bieżącej, czyli na tej, na której stoi
powłoka. Schemat wymaga wskazania urządzenia, a jego identyfikator poda komenda
`device.list`, gdy trafi do kontraktu. Do tego czasu drugi argument funkcji
`zalozKatalogLokalny` zostaje pusty i rdzeń odmawia z powodem, bez atrapy
identyfikatora.

## budowa/klient-poprzedni/src/moduly/agents/podglad-wywolania.ts

Podgląd pyta rdzeń o gotowe wywołanie zamiast składać wiersz polecenia po
swojemu. Gdyby okno składało wywołanie własnymi regułami, pokazywałoby wiersz,
którego tura nigdy nie wykona — rdzeń liczy podgląd tymi samymi funkcjami,
którymi wykonuje wysłanie wiadomości, wraz z nałożeniem eksperta okna.

Wywołanie liczy się dla okna komunikacji, a moduł agentów żadnego okna rozmowy
nie zna: okna modułu i okna rozmowy są odrębnymi bytami. Dlatego identyfikator
okna wpisuje Operator, a bez niego podgląd rdzenia nie pyta, zamiast sięgać po
okno przypadkowe.

Pusty prompt systemowy jest odpowiedzią, nie brakiem odpowiedzi: znaczy, że ani
oś, ani ekspert nie wnoszą treści systemowej. Pole zostaje wtedy widoczne wraz
z takim zdaniem, zamiast się schować i sugerować, że odczyt się nie odbył.

## budowa/klient-poprzedni/src/moduly/developer/okno-warsztatu-kodu.ts

Czynności warsztatu stoją w osobnym oknie zamiast w pasku pływającym nad
zaznaczeniem, ponieważ pole edycji Code Editora jest zwykłym obszarem tekstu
bez modelu dokumentu — zaznaczenie znika przy pierwszym kliknięciu poza nim.
Okno podaje ten sam zestaw czynności w postaci, która nie potrzebuje żywego
zaznaczenia: plik bierze się ze stanu modułu, a zaznaczenie wpisuje się
jawnie w polu formularza.

Refaktoryzacja i operacje kontekstowe modelu wracają jako podgląd zmiany, nie
jako zapis na dysk. Zapis jest osobnym rozstrzygnięciem Operatora, bo to
jedyna chwila, w której da się pracę modelu odrzucić przed jej utrwaleniem.

Nawigacja po symbolach i analiza statyczna niosą w odpowiedzi jawne pole
mówiące, czy serwer języka i program analizy statycznej były osiągalne. Okno
rozróżnia stan, w którym wystąpień nie ma, od stanu, w którym nie było czym
ich szukać — pierwszy naprawia się w kodzie, drugi instalacją programu po
stronie serwera.

## budowa/klient-poprzedni/src/mobile/wywolania-interwencji.ts

Moduł składa cztery drogi interwencji warstwy mobilnej w wywołania komend kontraktu rdzenia i jest jedynym miejscem tej warstwy, które zna nazwy komend zmieniających stan sesji.

Każda droga stoi na określonych komendach kontraktu. Zatwierdzenie kroku wywołuje komendę akcji kolejki z rozstrzygnięciem wznowienia albo powtórzenia. Wstrzymanie stosuje akcję kolejki pauzy, zatrzymanie wiadomości albo zatrzymanie sesji, przy czym wstrzymanie sesji oddaje wykaz okien, w których turę zatrzymano. Nastawienie koordynatora zapisuje konfigurację sesji na zasięgu okna i potwierdza wynik odczytem zwrotnym konfiguracji obowiązującej, który wskazuje źródło nadpisania. Przejęcie sterowania składa trzy wywołania: zatrzymanie wiadomości, zapis konfiguracji sesji z trybem uprawnień ręcznym oraz wysłanie wiadomości z poleceniem Operatora.

Droga przejęcia sterowania zatrzymuje turę, przestawia okno na pytanie o każdy krok i wpuszcza polecenie Operatora. Dedykowanej komendy przejęcia kontrakt nie niesie, dlatego złożenie trzech wywołań zastępuje jedno wywołanie, a ekran pozostaje bez zmian.

Droga zatwierdzenia kroku ma widoczną granicę: zatwierdzenie prowadzi wyłącznie przez kolejkę, więc pozycja bez kolejki nie ma czym zatwierdzić kroku, ponieważ komendy zatwierdzenia pojedynczego kroku kontrakt nie niesie. Funkcja rozstrzygająca dostępność dróg wyraża ten stan zdaniem, a ekran nie rysuje wtedy przycisku zatwierdzenia.

Każde wywołanie oddaje kwit zamiast wartości logicznej czy pustego wyniku: nazwę komendy, rozstrzygnięcie, stan po zmianie wyjęty z odpowiedzi rdzenia albo treść odmowy. Kwit jest jedynym dowodem interwencji, jaki pozostaje po zamknięciu aplikacji mobilnej.

## budowa/klient-poprzedni/src/moduly/apps/odmowa-rdzenia.ts

Powód osobnej ścieżki wywołania jest mechaniczny. Komenda bez uchwytu w rdzeniu wraca kopertą typu
`<obszar>.unknown` z polem `requestedType`. Ta koperta nie niesie pola `status`, a korelacja żądań
rozstrzyga wyłącznie koperty ze statusem, co opisuje `protokol/koperta.ts`. Obietnica zwykłego
wywołania pozostałaby więc nierozstrzygnięta na zawsze, a okno stałoby w stanie ładowania bez końca.

Rozwiązanie nie buduje drugiej drogi do rdzenia. Żądanie idzie tym samym `kanal.wyslij`, a odmowę
odczytuje się z dziennika komunikatów nierozpoznanych, który kanał prowadzi dla wszystkich obszarów
kontraktu w `polaczenie/dziennik-nieznanych.ts`. Moduł nie zna literału `apps.unknown` i nie zakłada
własnej subskrypcji zdarzenia odmowy.

Wynik odmowy jest zwykłym `Wynik` z polem `blad`, dzięki czemu okno pokazuje go tak samo jak każdą
inną odmowę rdzenia i nie potrzebuje osobnej gałęzi widoku.

## budowa/klient-poprzedni/src/moduly/design/modal-kreatora.ts

Cztery drogi zamknięcia są równorzędne i wszystkie obowiązkowe. Dwie z nich są
własnością natywnego elementu `dialog`: klawisz Escape obsługuje przeglądarka,
a kliknięcie w nakładkę rozpoznaje się po tym, że celem zdarzenia jest sam
element dialogu, a nie jego wnętrze — wnętrze przykrywa dialog w całości, więc
kliknięcie w treść nigdy nie dochodzi do elementu nadrzędnego. Ten sam wzorzec
niesie `konfiguracja/okno-konfiguracji.ts`.

Przycisk główny stopki zachowuje klikalność zawsze, także przy brakach w polach
i w czasie ładowania. Powtórzone naciśnięcie rozstrzyga logika akcji, a nie
odebranie klikalności kontrolce: kontrolka nieklikalna nie mówi Operatorowi,
czego brakuje, więc blokada należy do warstwy, która zna powód.

Wywołanie `showModal` daje nakładkę, umieszczenie na stosie okien i obsługę
klawisza Escape. W środowisku, które tej metody nie ma, powłoka otwiera okno
zwykłym ustawieniem stanu otwarcia, ponieważ okno ma się pokazać niezależnie od
dostępności metody.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-izolacji-apps.ts

Granica pojęciowa jest tu istotna i okno musi ją nazwać wprost. Komenda
`isolation.policy.preview` oddaje politykę platformy rozstrzygniętą po ośmiu
poziomach zasięgu, a nie izolację nadaną pojedynczemu rozszerzeniu. Osiem
zakresów technicznych — katalog roboczy, środowisko procesu, dostęp sieciowy,
odczyt i zapis plików, konto i token, model procesu, serwer wykonania oraz
katalog danych modelu — opisuje warunki, w jakich kod się wykonuje, a w tych
samych warunkach wykona się kod rozszerzenia. Deklaracji uprawnień rozszerzenia
kontrakt nie niesie, więc okno mówi o tym wprost, zamiast podstawiać jedno
za drugie.

Poziom zasięgu źródło bierze najwęższy, jaki moduł zna. Podanie okna modułu każe
rdzeniowi rozstrzygnąć dziedziczenie aż do niego, więc odpowiedź opisuje
politykę obowiązującą pracy prowadzonej w tym oknie. Bez okna źródło pyta
o poziom modułu, a wtedy odpowiedź opisuje warunki wspólne wszystkim jego oknom:
zdanie pozostaje prawdziwe, tyle że szersze.

Warstwy żądanie nie podaje. Warstwa pominięta znaczy warstwę obowiązującą,
a wskazanie którejkolwiek innej byłoby cudzym rozstrzygnięciem przebranym
za odczyt stanu.

## budowa/klient-poprzedni/src/moduly/browser/zapisy-przekazania.ts

Moduł skupia ścieżki zapisu wychodzące poza obszar komend `browser.*` i jest
trzymany osobno od odczytu w `zrodlo-browser.ts`, ponieważ dotyczy innych
obszarów kontraktu i ma innego adresata.

Przekazanie międzymodułowe, uruchamiane poleceniami wysłania do innego modułu
oraz tłumaczenia, idzie komendą `context.transfer`. Jest to jedyna droga
przeniesienia kompletu kontekstu między modułami, jaką niesie kontrakt i jaką
obsługuje rdzeń.

Pytanie o zaznaczenie, uruchamiane poleceniem wyjaśnienia na pasku pływającym,
idzie komendą `message.send`. Oknem, do którego pytanie trafia, jest okno modułu
Browser, to samo, którego identyfikator niosą komendy `browser.*`.

Adnotacja spłaszczona do obrazu rastrowego idzie tą samą drogą co pytanie:
`message.send` kładzie treść w rozmowie wskazanego okna i niesie pole
`attachments`, natomiast `context.transfer` wymaga pola `targetModuleId`
i wysłałby rysunek poza moduł, w którym powstał, nie stawiając go w rozmowie
wcale.

Pole `attachments` żądania `message.send` jest polem kontraktu: rdzeń przepisuje
je do zakładanej wiadomości w adapterze rozmowy i oddaje w odpowiedzi, dzięki
czemu okno ocenia dołączenie po wierszu, który wrócił, a nie po tym, że wysłało.

## budowa/klient-poprzedni/src/komponenty/naglowek-okna.ts

Jeden kształt górnego pasa obowiązuje we wszystkich modułach, ponieważ Operator ma otrzymywać tę samą odpowiedź na pytanie o tożsamość oglądanego okna niezależnie od tego, który moduł je otworzył. Granica wobec ramy okna przebiega tak, że rama obejmuje cały kontener okna wraz z ciałem i stanami, natomiast nagłówek odpowiada wyłącznie za górny pas i nie zna treści umieszczonej pod sobą. Nagłówek nie osadza się samodzielnie w drzewie dokumentu, lecz dostaje go rama okna albo plik okna.

Plakietka roli stoi bezpośrednio przy nazwie, ponieważ nazwy okien powtarzają się między modułami: Process Monitor w module Terminal i Execution Monitor w module Automations to dwa odrębne okna. Podpis roli jest tym elementem, który je na ekranie rozróżnia, dlatego napis w plakietce bywa zarówno samym określeniem roli, jak i zdaniem doprecyzowującym przeznaczenie okna.

Wygląd pasa pochodzi w całości z biblioteki: pas stoi na klasie dn-karta-naglowek, nazwa na klasie dn-karta-tytul, a plakietka roli na klasie dn-plakietka--rola, bez klas własnych i bez barw zapisanych w kodzie. Moduł, który potrzebuje odstępstwa, podaje własną klasę polem klasa i styluje nazwę selektorem potomka, zamiast powielać komponent. Wyrównanie kontrolek do prawej krawędzi należy do modułu wołającego, który opakowuje je klasą dn-pasek-prawa.

## budowa/klient-poprzedni/src/moduly/design/okno-tokens-system-panel.ts

Okno Tokens & System Panel stoi jako rozwinięcie warstwy trzeciej, bo opracowanie wywołuje
je z menu kebab obszaru roboczego, a nie stawia go na widoku spoczynkowym. Wartości żetonów
czyta z motywu obowiązującego w danej chwili (`zetony-systemu.ts`), pary i progi kontrastu
z wykazu progów produktu (`kontrast-wcag.ts`), a powiązania żetonu z komponentami — z arkuszy
wczytanych do dokumentu; nic z tego nie jest przepisane do modułu, więc poprawka w motywie
jest widoczna od razu, a rozjazd nazw wychodzi na wierzch zamiast zniknąć.

Okno nie zapisuje zestawu żetonów, bo kontrakt nie zna takiego bytu — wydanie wychodzi
plikiem do przeglądarki Operatora. Nie nadpisuje motywu wczytanym zestawem, bo motyw jest
własnością powłoki, a import kończy się zestawieniem różnicy. Nie wydaje przewodnika stylu
do modułów Library i Studio, bo kontrakt nie zna przekazania tego bytu między modułami.

Przełącznik motywu podglądu nie jest kopią przełącznika powłoki: woła tę samą czynność
motywu produktu, więc drugiego stanu motywu nie ma. Okno nasłuchuje zmiany motywu i przelicza
wszystkie pomiary, bo oba motywy są równoprawne i każdy wymaga własnego pomiaru.

Panel zestawów trwałych czyta żetony motywu w chwili naciśnięcia, a nie z kopii zrobionej
przy otwarciu okna, ponieważ między otwarciem a naciśnięciem Operator mógł przełączyć motyw,
a zestaw ma opisywać stan obowiązujący w chwili działania.

Porównanie wczytanego zestawu z żetonami motywu idzie po nazwach ról, nie po strukturze
pliku: zestaw obcy bywa zagnieżdżony dowolnie, a jedyne, co da się z nim zrobić uczciwie, to
powiedzieć, które role produktu w nim są, które mają inną wartość i których produkt nie zna.

## budowa/klient-poprzedni/src/dostepy/zapisy-nadan.ts

Odczyt i zapis nadań to dwie różne odpowiedzialności. Stan sekcji pilnuje tego, co widok wie
o punktach dostępu i nadaniach, a ten plik pilnuje tego, co dzieje się po zapisie.

Komplet zastępuje zbiór w całości. Wszystkie trzy komendy zwracają nie tylko zmienione nadanie,
lecz komplet nadań okna po zmianie, ponieważ przestawienie kolejności albo oznaczenia głównego
dotyka pozostałych wierszy. Widok bierze komplet zamiast składać go z domysłów.

Nadanie bez okna rozmowy nie idzie do rdzenia, ponieważ nadanie żyje przy oknie i bez okna nie ma
czego nadać. Odmowa przyjmuje kształt wyniku komendy, żeby widok nie potrzebował drugiej ścieżki
obsługi.

## budowa/klient-poprzedni/src/moduly/apps/okno-backend-workspace.ts

Plik jest wiązaniem, a nie drugim widokiem. Formularz warsztatu stoi raz,
w `okno-warsztatu.ts`, natomiast osobny plik daje jedno miejsce w drzewie, po
którym widać, że okno warstwy usług jest zbudowane, a nie tylko wymienione
w słowniku kodów okien.

Kod okna nie pada tu literałem: ramę woła `utworzRameApps(opis.kodOkna, …)` ze
zmiennej, dzięki czemu `KODY_OKIEN` pozostaje jedynym miejscem, w którym kody
okien modułu Apps są zapisane. Okna Studio rozdziela ta sama zasada —
`okno-studio.ts` jest wspólną ramą, a poszczególne okna mają własne pliki.

## budowa/klient-poprzedni/src/moduly/library/tresc-podgladu.ts

Podgląd graficzny oraz podgląd stron dokumentu PDF przychodzą z rdzenia
wyłącznie jako odnośnik w polu `imageRef`. Klient nie ma komendy pobierającej
bajty spod tego odnośnika, więc ciało podglądu pokazuje sam odnośnik i nazywa
brak treści, zamiast udawać obraz, którego nie otrzymał.

Pole `truncated` odpowiedzi mówi o skróceniu podglądu po stronie rdzenia.
Wypisane pod treścią zdanie odróżnia fragment pliku od jego całości, dzięki
czemu odczyt niepełny nie jest brany za odczyt kompletny.

## budowa/klient-poprzedni/src/konfiguracja/kontrolka.ts

Formularz okna konfiguracji powstaje z katalogu pozycji i z założenia nie zna żadnej kontrolki z osobna. Zna wyłącznie wspólny kształt zadeklarowany w tym pliku, co pozwala złożyć ekran z pozycji katalogu bez wiedzy o rodzaju wartości, jaki za nimi stoi.

Konsekwencją takiego rozdziału jest koszt rozszerzenia. Dodanie kolejnego rodzaju wartości do kontraktu sprowadza się do jednego przypadku w module rozdzielającym wybor-kontrolki.ts oraz do jednej funkcji budującej kontrolkę. Nie powstaje przy tym nowy ekran ani zmiana w samym formularzu.

## budowa/klient-poprzedni/src/moduly/assistant/okno-memory-context-manager.ts

Okno zamyka obszar pamięci modułu Assistant i realizuje zasadę jawności: pamięć
asystenta jest w całości widoczna, edytowalna i usuwalna. Cztery pierwsze
zakładki odpowiadają czterem obszarom opracowania modułu — fakty, pamięć
semantyczna, konteksty i baza wiedzy.

Plik odpowiada wyłącznie za skład okna. Wywołania kontraktu mieszkają
w `zrodlo-pamieci.ts`, a każda zakładka ma własny plik obszaru: ustalenia
w `panel-faktow.ts`, wskaźnik znaczenia w `panel-wiedzy.ts`, poziomy pamięci
karty sesji w `panel-kontekstow.ts`.

Fazy odczytu nie ma na poziomie okna, tylko w obszarach. Cztery zakładki czytają
cztery różne byty i wołają je w różnych chwilach — jedna faza dla całego okna
kazałaby odmowie odczytu pamięci przesłonić zakładkę wskaźnika, która o tej
odmowie nic nie wie.

Piąta zakładka „Zestawy i retencja" zamyka trzy obszary, które kontrakt niesie
w całości: nazwane konteksty pamięci `memory.context.*`, zasady retencji
i wygaszania `memory.retention.*` oraz miernik okna kontekstu
`context.usage.get`. Miernik pokazuje liczbę policzoną tokenizatorem rdzenia
wraz z nazwą słownika, którym policzono, a gdy pomiar jest niewykonalny — sam
powód.

## budowa/klient-poprzedni/src/modele/zakladki-sekcji.ts

Cztery obszary sekcji odpowiadają czterem pytaniom o ten sam model: konta rozstrzygają, czym prowadzone jest połączenie, ustawienia osi rozstrzygają, jak model ma działać, tożsamość rozstrzyga, kim model ma być, a podgląd promptu pokazuje, co ostatecznie trafia do modelu. Wszystkie cztery obszary mówią względem tej samej osi, dlatego trzyma je jedna sekcja, a nie cztery osobne sekcje.

Przełącznik nie porzuca obszaru, z którego Operator wychodzi. Element obszaru zostaje w drzewie dokumentu i jedynie przestaje być widoczny, dzięki czemu tekst wpisany w edytorze tożsamości przeżywa zajrzenie do obszaru kont i powraca w niezmienionej postaci.

## budowa/klient-poprzedni/src/moduly/multitasking/kontrolki.ts

Przycisk, pole, wybór, wykaz i pozycja wykazu pochodzą z pliku
`modele/kontrolki-formularza.ts`, a moduł bierze je stamtąd zamiast trzymać
własne kopie. Wiersz klucz–wartość odpowiednika tam nie ma: niosą go wyłącznie
okna ról, gdzie zastępuje tabelę stanu — więź wykonawcy z koordynatorem,
licznik obiegów oraz powód zatrzymania biegu.

Wygląd w całości pochodzi z arkusza rodziny `dm-`, przez co plik nie zna ani
jednej barwy i ani jednego odstępu. Przedrostek klas stoi jedną stałą, bo
arkusz rodziny należy do modułu, a trzy wykazy okien ról muszą wskazywać ten
sam arkusz.

Wysokość pól redakcyjnych jest wspólna dla trzech okien ról, ponieważ kreator
promptu, polecenie wykonawcy i uzasadnienie oceny stoją obok siebie na scenie.

## budowa/klient-poprzedni/src/aod/sekcja-kontekstu.ts

Pole `executionParams` ma w kontrakcie typ `unknown`, więc sekcja go nie rozbiera i melduje jedynie
obecność parametrów wykonania. Rozbiór wymagałby założenia o kształcie danych, którego kontrakt nie
gwarantuje.

Komplet pusty jest stanem poprawnym, a nie odmową odczytu, ponieważ wszystkie pola `ContextBundle`
są w kontrakcie opcjonalne. Sekcja mówi to osobnym zdaniem zamiast wypisywać wykaz samych oznaczeń
braku.

## budowa/klient-poprzedni/src/mobile/pasek-kwitu.ts

Pasek czyta kwity z pamięci trwałej urządzenia, a nie ze stanu okna.
Potwierdzenie trzymane w stanie okna ginie razem z oknem, natomiast kwit zapisany
w magazynie przeżywa zamknięcie ekranu i ponowne uruchomienie telefonu, dzięki
czemu wykaz wysłanych decyzji pozostaje kompletny po powrocie do aplikacji.

Wiersz paska powstaje ze zdania złożonego przez `zdanieKwitu` z odpowiedzi
rdzenia, nie z zamiaru klienta. Kwit nieudany stoi na tym samym pasku wraz
z treścią odmowy, ponieważ odmowa jest rozstrzygnięciem równie wiążącym jak
przyjęcie i ma być widoczna w tym samym miejscu.

## budowa/klient-poprzedni/src/mission-control/naglowek-pulpitu.ts

Nagłówek ma jedną odpowiedzialność: podaje tytuł ekranu i jawnie nazywa
pochodzenie liczb pokazywanych niżej. Dopóki odczyt z rdzenia nie nadszedł,
a więc dopóki źródło danych ma wartość `ZrodloDanych.Oczekiwanie`, przy tytule
stoi plakietka oczekiwania wraz ze zdaniem wyjaśniającym. Zdanie jest konieczne,
ponieważ bez niego pusty ekran zostałby wzięty za pomiar mówiący, że pracy
nie ma.

Po pierwszym odczycie plakietka mówi o danych pochodzących z rdzenia. Pulpit nie
zna innych źródeł: liczby albo pochodzą z odczytu rdzenia przez kanał kontraktu,
albo nie ma ich wcale, a sekcje pokazują stany puste.

## budowa/klient-poprzedni/src/moduly/library/wykaz-plikow.ts

Struktura katalogów pochodzi z danych, nie z osobnego zapytania: kontrakt nie ma
komendy katalogu folderów, niesie za to pole ścieżki każdego pliku. Wykaz
katalogów składa się więc z pierwszych członów ścieżek zwróconych przez rdzeń,
a zbiór pusty daje jedną pozycję obejmującą cały zbiór.

Formę prezentacji rozstrzyga widok wybrany w oknie i zbudowany w module widoków
wykazu. Wykaz nie zna żadnej z pięciu form: składa katalogi, oddaje zbiór
widoczny i osadza to, co widok zbudował.

## budowa/klient-poprzedni/src/moduly/multitasking/plan-etapow.ts

Etap planu jest kolejką, a nie osobnym bytem klienta: drugiego rejestru etapów po stronie
okna nie ma, więc wykaz odświeża wyłącznie zdarzenie `queue.changed` przychodzące z rdzenia.

Kontrolka zależności etapów pozostaje widoczna i nieczynna. Komenda
`orchestration.dependency.set` wiąże kroki układów automatyk polami `fromStepId` oraz
`toStepId`, a etap tego planu jest kolejką, nie krokiem automatyki. Powód nieczynności bierze
się z wykazu komend oddanego przez rdzeń, nie ze stałej wpisanej w kodzie okna.

Potwierdzenie założenia etapu składa się z pól odpowiedzi rdzenia, nie z pól żądania:
przycięta nazwa albo nieprzyjęte okna wykonawców mają być dla Operatora widoczne.

Wiersz etapu bieżącego niesie plakietkę zamiast wygaszonego przycisku, a kontrolka wyboru
pojawia się wyłącznie w wierszach, w których jest co wybrać.

## budowa/klient-poprzedni/src/moduly/library/tresc-base64.ts

Odczyt jest wspólny dla wgrania pliku komendą `library.file.upload` oraz dla
dołożenia wersji komendą `library.version.add`, ponieważ obie komendy niosą to
samo pole `contentBase64`.

Nieudany odczyt, czyli zniknięcie pliku albo odrzucenie dostępu przez
urządzenie, wraca jako `null`, a nie wyjątkiem wywracającym widok. Wołający ma
wtedy powiedzieć Operatorowi, że nic nie zostało wysłane.

## budowa/klient-poprzedni/src/moduly/apps/rama-okna.ts

Obudowa okna, na którą składają się nagłówek, plakietka roli, pas akcji i ciało, należy do komponentu ramy okna w katalogu komponenty, a znakowanie fazy do komponentu fazy okna w tym samym katalogu. Plik ramy modułu Apps splata oba byty, ponieważ rama biblioteczna nie zna cyklu życia danych, których sama nie pobiera, a pięć okien modułu ma mieć jedną przesłonę stanu zamiast pięciu wariantów tego samego rozwiązania.

Przesłona stanu nie kasuje treści, tylko ją przykrywa. Powrót do fazy gotowej odsłania to, co Operator już widział, dzięki czemu nieudane odświeżenie nie zabiera wyniku poprzedniego odczytu. Wskaźnik odczytu stoi obok komunikatu, a nigdy zamiast kontrolki, więc treść okna zostaje na miejscu i pozostaje klikalna. Wspólna procedura znakowania fazy tego nie przesądza, ponieważ moduły różnią się między sobą pod tym względem.

Wygląd pochodzi w całości z biblioteki komponentów z przedrostkiem dn oraz z żetonów motywu, dlatego plik nie zawiera ani jednej wartości barwy i ani jednej wartości odstępu.

## budowa/klient-poprzedni/src/aplikacja/widok-pulpitu.ts

Pulpit operacyjny stoi obok strony głównej, a nie zamiast niej. Centrum
dowodzenia jest wejściem przy rozpoczynaniu pracy, pulpit — przy powrocie
do pracy już rozpoczętej, a obydwa pozostają dostępne z przełącznika widoków.

Dane pulpitu pochodzą z rdzenia. Źródło pulpitu odpytuje `session.list`,
`window.list` oraz `channel.list`, a następnie subskrybuje zdarzenia
`session.changed`, `window.changed`, `queue.changed` i `progress.changed`.
Do pierwszego odczytu pulpit pokazuje stany puste, ponieważ liczby wymyślonej
nie poda.

Zamiary Operatora idą do rdzenia tym samym kanałem, którym przyszły dane:
`session.open` oraz `queue.action` przechodzą kanałem kontraktu.

## budowa/klient-poprzedni/src/aplikacja/stan-pary-gniazda.ts

Stan pary ma dwa człony. Skład pary mówi, kto jest koordynatorem, a kto
wykonawcą; położenie pary mówi, czy wykonawca pracuje, czy kolejka stoi.
Pierwszy człon niesie zdarzenie `window.changed`, drugi — `queue.changed`.

Właścicielem pola `windowRole` jest rdzeń, nie scena. Po odtworzeniu okna
z bazy, po `window.update` z innego urządzenia albo po `context.transfer` może
wrócić rola inna niż domyślna nadana przy tworzeniu gniazda. Dlatego każde
`window.changed` dotyczące tego okna przestawia nagłówek, a zmianę na tę samą
rolę gniazdo pomija, żeby widok nie przebudowywał się bez powodu.

Okno poza pętlą nie ma pary, więc jego nagłówek zostaje bez plakietki stanu
zamiast pokazywać stan cudzej kolejki.

Każde gniazdo słucha całej magistrali zdarzeń, więc zdarzenie dotyczące innego
okna jest pomijane. Kolejka bez wykazu okien dotyczy całej sesji i wchodzi do
każdego gniazda tej sceny.

## budowa/klient-poprzedni/src/modele/okno-modeli.ts

Okno stoi na natywnym elemencie okna dialogowego: warstwę tła, stos okien
i zamknięcie klawiszem ucieczki daje przeglądarka, a nie własna nakładka. Wygląd
bierze z biblioteki komponentów, z klasy okna modalnego, więc plik nie ustala
barw.

Rama wyłącznie osadza sekcję. Treść i stan mieszkają w module sekcji modeli, aby
tę samą sekcję dało się wstawić także w widok osadzony bez powielania kodu.

## budowa/klient-poprzedni/src/modele/zrodlo-tozsamosci.ts

Kategorie zasad są danymi, nie kodem: konstytucja, profil roli, ekspertyza,
zasady bezpieczeństwa i zasady harnessu przychodzą komendą
`identity.category.list` wraz z warstwą, porządkiem, obowiązkowością i trybem
proponowanym. Dopisanie kategorii jest wtedy nowym wierszem katalogu, a nie
zmianą kodu klienta.

Nakładka obowiązująca jest odczytem rdzenia, nie sklejeniem w kliencie: klient
nie składa promptu z treści kategorii własnym porządkiem, tylko pyta
`identity.effective.get` i pokazuje to, co trafi do modelu. Drugie składanie po
stronie widoku dałoby podgląd rozjeżdżający się z rdzeniem przy pierwszej
zmianie reguł.

## budowa/klient-poprzedni/src/moduly/design/symulacja-widzenia.ts

Zakres symulacji jest tu ważniejszy niż sam rachunek: obejmuje ona próbki
żetonów, a nie obrazy modułu. Obrazu nie ma czym pobrać, ponieważ pole odsyłacza
zasobu niesie ścieżkę w systemie plików rdzenia, więc nałożenie symulacji na
zasób wizualny nie ma dziś drogi wykonania i panel tego nie udaje.

Rachunek jest przybliżeniem i tak jest nazwany. Macierze odwzorowują trzy
dichromazje wprost w przestrzeni sRGB, bez przejścia przez przestrzeń
długofalową, co jest uproszczeniem przyjętym w narzędziach projektowych.
Wystarcza ono, żeby zobaczyć, które dwie barwy systemu zlewają się w jedną,
i nie wystarcza do orzeczenia medycznego. Panel mówi to Operatorowi wprost,
ponieważ różnica między podglądem a badaniem jest tutaj istotna.

Granica zlania jest odległością w składowych sRGB — miarą zgrubną, dobraną tak,
żeby wskazywała pary wymagające obejrzenia przez człowieka, a nie żeby
rozstrzygała za niego.

## budowa/klient-poprzedni/src/moduly/apps/dymek-objasnienia.ts

Wiersz `mp-pole-z-dymkiem` jest wstawką układu właściwą modułowi: pole rośnie na
całą dostępną szerokość, a znak zapytania stoi przy jego prawej krawędzi.

Klasy własne idą do wspólnej fabryki parametrem, ponieważ znak modułu nie jest
bibliotecznym przyciskiem ikonowym. Klasa `mp-dymek__znak` opisuje obwódkę
o wymiarze `--dn-wym-ikona-sm` i narożniku `--dn-r-pill` ze wskaźnikiem myszy
`help`, natomiast biblioteczna klasa `dn-btn-ikona` opisuje kwadrat o boku
trzydziestu dwóch pikseli. Przekazanie klas parametrem pozwala zachować obie
postacie bez rozgałęziania fabryki.

## budowa/klient-poprzedni/src/dostepy/elementy-karty.ts

Elementy kart sekcji dostępów powstają w jednym miejscu, żeby karta punktu,
wiersz nadania i obszar dodawania katalogu miały spójny wygląd. Plik nie zna
barw, odstępów ani reguł widoku: całość wyglądu niosą klasy biblioteki
`komponenty/` oraz klasa modyfikująca podana przez wywołującego. Plakietka
znaku przyjmuje klasę biblioteki i klasę miejsca osobno, ponieważ pierwsza
opisuje wygląd, a druga położenie w układzie karty.

## budowa/klient-poprzedni/src/moduly/assistant/stan-okna.ts

Plik odpowiada wyłącznie za nośnik komunikatu stanu wraz z miejscem na treść okna. Nazwę fazy oraz znakowanie drzewa dokumentu wnosi komponent fazy okna z katalogu komponenty, więc tutaj zostaje jedynie to, czym moduł Assistant różni się od reszty drzewa: wskaźnik odczytu stoi obok treści, a treść pozostaje widoczna również w fazie ładowania, dzięki czemu kontrolki są klikalne przez cały czas trwania wywołania.

Brak zapytania do rdzenia nie jest fazą okna, lecz fazą źródła danych, którą niesie zapis modułu w polach opisujących zadane pytania. Okno pokazuje ten przypadek jako fazę pustą z własnym zdaniem wyjaśniającym. Stan nie zastępuje treści, tylko ją przesłania, więc po powrocie do fazy gotowej wcześniejsza treść pozostaje nietknięta, a nieudane odświeżenie nie kasuje tego, co było już na ekranie.

## budowa/klient-poprzedni/src/dostepy/ostrzezenie-zapisu.ts

Specyfikacja mostu mcp-danaco-pulpit-console zabrania nadawania trybu zapisu na maszynie danaco-data bez wyraźnej potrzeby. Maszyna ta jest hostem produkcyjnej platformy LEX, a zapis modelu sięga tam zbiorów, z których korzysta cała kancelaria. Zakazu nie egzekwuje jednak kod klienta: nadanie zapisu pozostaje możliwe, ponieważ rozstrzyga o nim Operator. Egzekwuje go jawne, widoczne zdanie postawione przy przełączniku trybu, a nie podpowiedź ukryta pod kursorem.

Wykaz maszyn chronionych stoi w kodzie, a nie w katalogu konfiguracji, ponieważ jest to ostrzeżenie bezpieczeństwa i nie może dać się wyłączyć zapisem w bazie. Rozszerzenie wykazu sprowadza się do jednego wiersza w stałej MASZYNY_CHRONIONE.

## budowa/klient-poprzedni/src/moduly/design/kanaly-obrazowe.ts

Żądanie `design.asset.generate` niesie pole `channelId`, a rdzeń sprawdza rodzaj
wskazanego kanału zanim cokolwiek wyśle — kanał tekstowy odmawia kodem
`validation_failed` w `core/adapter_modul_design_kanal.go`. Wykaz silników
pokazuje więc kanały obrazowe jako wybieralne, a pozostałe wymienia z nazwy
i z powodem, zamiast je ukrywać.

Kolejność rozpoznania jest przepisana z funkcji `KluczAdaptera()`
w `models/definicja.go`: parametr `adapter` ma pierwszeństwo, a rodzaj wiersza
jest wartością zapasową. Do kontraktu obie strony jadą tym samym wierszem
(`core/adapter_kanaly.go`, pola `Kind: k.RodzajKanalu` oraz
`Config: ParametryJSON`), więc `Channel.kind` odpowiada polu `Rodzaj`,
a `Channel.config.adapter` polu `Parametr`. Kluczem adaptera obrazowego jest
napis `obrazy` — stała `AdapterObrazy` w `models/kanal.go`.

Rozpoznanie może rozejść się z rdzeniem: klient czyta `config` jako treść
nieokreśloną (`config?: unknown` kontraktu), więc kanał o parametrach zapisanych
inaczej niż obiektem trafi między nieobrazowe. Ostateczną kontrolę wykonuje
rdzeń — dlatego wskazanie kanału spoza wykazu obrazowego nie jest w oknie
blokowane, tylko opisane.

## budowa/klient-poprzedni/src/moduly/browser/indeks.ts

Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
rozstrzyga, gdzie go postawić. Dzięki temu te same okna wchodzą i w obszar
roboczy powłoki, i w podgląd sprawdzianu.

Wytwórnia oddaje moduł wprost, ponieważ `ModulBrowser` niesie pola `element`
i `wczytaj` żądane przez powłokę oraz czynność `rozlacz`, po którą sięga
sprawdzian modułu.

## budowa/klient-poprzedni/src/konfiguracja/kontrolki-liczbowe.ts

Granice i skok pola liczbowego pochodzą z katalogu ustawień, a nie z kodu
klienta. Katalog, który granic nie podaje, daje pole bez granic: brak metadanej
nie jest błędem i niczego nie blokuje, więc kontrolka nie wymyśla ograniczenia,
którego rdzeń nie postawił.

Pole puste znaczy brak wartości, a nie zero. Odczyt oddaje wtedy `null`, przez
co zapis czyści wartość zamiast wpisywać liczbę, której Operator nie podał.
Rozróżnienie jest konieczne, ponieważ zero bywa wartością obowiązującą
i podstawienie go pod pustkę zmieniałoby ustawienie bez wiedzy Operatora.

## budowa/klient-poprzedni/src/mission-control/pasek-dzialan.ts

Pasek ostatniego działania ma jedną odpowiedzialność: wypisuje, co pulpit właśnie
nadał powłoce. Żaden przycisk pulpitu nie jest wyszarzony, a przycisk czynny bez
widocznego skutku nie mówi Operatorowi, czy zamiar poszedł dalej. Pasek pokazuje
zatem fakt nadania zamiaru, a nie wynik jego wykonania.

Stąd bierze się dobór słów w treści paska. Zdanie mówi „zamiar nadany", nigdy
„wykonano", ponieważ w chwili wypisania powłoka mogła jeszcze nie wysłać komendy,
a rdzeń mógł jej nie potwierdzić.

## budowa/klient-poprzedni/src/ikony/marka.ts

Znak czyta się jako podwójny grot z kropką sygnału ustawioną na linii bazowej. Godło uproszczone
niesie jeden grot i obowiązuje od szesnastu pikseli boku w dół, ponieważ poniżej tej wartości drugi
grot oraz prześwit między grotami przestają być czytelne.

Odmianę znaku dobiera się do podłoża, a nie do motywu interfejsu, ponieważ barwy znaku są wpisane
w plik źródłowy. Pasek kokpitu pozostaje atramentowy w obu motywach, więc leży na nim odmiana
przeznaczona na podłoże ciemne, niezależnie od tego, który motyw jest czynny.

## budowa/klient-poprzedni/src/moduly/apps/warsztaty-apps.ts

Opisy warsztatów stoją osobno od ramy warsztatu, ponieważ rama jest czynnością
budującą okno, a opisy są wykazem danych. Wykaz odwzorowuje warstwy wyliczenia
`AppWorkspaceLayer` z kontraktu, więc trzeci warsztat może powstać dopiero po
rozszerzeniu kontraktu, a nie po zmianie samego widoku.

## budowa/klient-poprzedni/src/aplikacja/polaczenie-z-rdzeniem.ts

Pole nawigacji platformy niesie komendy wejścia do środowiska oraz wykazu
modułów osadzone na wspólnym kanale i służy wyłącznie powłoce: powłoka powstaje
bez połączenia i nie ma jak sięgnąć po kanał sama, więc dostaje te dwie drogi
gotowe. Widoki mające kanał wołają opakowania wprost, widok strony głównej
komendę wejścia na stronę główną, a przestrzeń modułu komendę wejścia do
przestrzeni pracy, bez obiektu pośredniego.

Plik wyłącznie składa warstwy: nie zna ramki, nie buduje koperty i nie zna nazwy
żadnej komendy. Nazwy pochodzą z pakietu wspólnego i żyją w warstwie protokołu.

Złożenie nie otwiera połączenia. Rozpoczęcie łączności należy do cyklu życia,
żeby moment jej nawiązania był jednym miejscem, a nie skutkiem ubocznym budowy
obiektów.

Token sesji bramki jest czytany domknięciem przy każdym powitaniu, a nie raz
przy składaniu, ponieważ po ponownym nawiązaniu połączenia obowiązuje sesja
bieżąca, nie ta sprzed zerwania.

Nawigacja powłoki składa się w tym pliku, a nie w uzgodnieniu, ponieważ
uzgodnienie odpowiada wyłącznie za powitanie, sesję i okno, i o nawigacji nic
nie wie. Korzeń montażu klienta jest jedynym punktem znającym jednocześnie kanał
i odbiorcę tych dróg.

## budowa/klient-poprzedni/src/moduly/apps/okno-architecture-designer.ts

Okno kreatora modułu Apps jest punktem wyjścia procesu tworzenia architektury. Definiowanie komponentów rozwiązania i ustalenie zależności stoją na jednej komendzie kontraktu, która przyjmuje szablon i komplet komponentów wraz z ich zależnościami. Przycisk walidacji uruchamia tę samą komendę, nie osobną drogę: kontrakt nie ma osobnej komendy walidacji, zwraca za to zastrzeżenia w architekturze po zapisie. Przycisk pozostaje klikalny przy pustej kanwie, ponieważ odpowiedź jest opisowa, nie blokadą. Zdanie potwierdzające bierze treść z odpowiedzi rdzenia, nie z zamówienia: wymienia identyfikator, wersję, szablon i komponenty tak, jak oddał je rdzeń, a każdą rozbieżność z zamówieniem i każde zastrzeżenie walidacji ogłasza odmową.

Odczyt architektury zapisanej w rdzeniu nadpisuje kanwę i mówi o tym wprost: architektura oddana przez rdzeń wchodzi wraz ze swoim kompletem komponentów, więc komponenty zestawione na kanwie, a jeszcze niezapisane, po tym kroku znikają. Zdanie wymienia liczbę komponentów, która przyszła. Brak architektury jest odpowiedzią, nie odmową — rdzeń oddaje wtedy wynik udany bez pola, a kanwa zostaje nietknięta.

Funkcja porównująca oddany układ z zamówionym zestawia zbiory identyfikatorów, nie listy. Rdzeń trzyma komponenty i krawędzie zależności w osobnych tabelach i oddaje je w porządku własnym, a powtórzoną krawędź zdejmuje kluczem pierwotnym. Porządek i powtórzenie nie są więc rozbieżnością, tylko odmiennym sposobem przechowywania danych; rozbieżnością jest komponent, który zniknął, albo taki, którego nikt nie zamawiał. Zależność jest wskazaniem identyfikatora innego komponentu, więc jej zgubienie zmienia układ tak samo jak zgubienie komponentu, a kanwa zaciągnięta z odpowiedzi sama z siebie nie pokazuje, że zamówiono więcej.

Zdanie o pustym polu zastrzeżeń walidacji orzeka o odebranej ramce, nie o zachowaniu rdzenia. Rdzeń tego pola nie wylicza: adapter czyta zastrzeżenia z bazy, zapisuje je z powrotem tym samym zapisem scalającym i oddaje w odpowiedzi, więc jedynym pisarzem kolumny jest ten, kto ją przed chwilą przeczytał. Zdanie o braku zastrzeżeń byłoby zapewnieniem o sprawdzeniu, którego nikt nie wykonał. Własny walidator w oknie byłby drugą prawdą o tym samym układzie i zacierał granicę odpowiedzialności: zastrzeżenia widoczne w oknie sugerowałyby, że zna je rdzeń. Rozróżnienie pola nieobecnego od pola obecnego i pustego sprawia, że w chwili, w której rdzeń zacznie zastrzeżenia liczyć, zdanie zmieni treść bez tknięcia tego pliku.

## budowa/klient-poprzedni/src/aplikacja/korzen-dokumentu.ts

Całym dokumentem zajmuje się `gospodarz-dokumentu.ts`, a kształt
`KorzenAplikacji` jest wejściem funkcji `zamontujUkladOkien` z katalogu
`okna-rownolegle/`. Korzeń pozostaje przez to warstwą montażową bez wiedzy
o rdzeniu i bez wiedzy o układzie okien.

Miejsce akcji nie powstaje w korzeniu — przychodzi z paska górnego powłoki
środowiska, żeby uchwyt szuflady sterowania stał w pasku wspólnym z resztą
akcji, a nie w osobnym pasku założonym dla jednego przycisku.

Klasa `dn-sesja` należy do karty sesji w pasie powłoki
(`powloka/karta-sesji.ts`), a jej arkusz zawęża szerokość do pasa nawigacji bez
zawężenia selektora rodzicem. Korzeń sceny nosi własną nazwę `dn-scena-sesji`,
ponieważ w przeciwnym razie wpadłby pod tamtą regułę i zszedł do szerokości
słupka nawigacji.

## budowa/klient-poprzedni/src/moduly/library/indeks.ts

Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
a warstwa składająca rozstrzyga, gdzie go postawić.

Kod `library` odpowiada kolumnie `modul.kod` w rdzeniu, a wpis do rejestru
modułów wiąże widok z modułem właśnie po tym kodzie.

## budowa/klient-poprzedni/src/moduly/agents/zlecenia-agentow.ts

Okno pracuje polami formularza: napisami z kontrolek i wartościami list wyboru.
Kontrakt pracuje polami opcjonalnymi, których pusta wartość ma znaczyć „nie
podano", a nie „podano pustkę". Struktury zleceń są granicą między jednym
a drugim: pole puste zostaje w oknie, a do rdzenia idzie żądanie bez niego.

Poziomy pamięci jadą zawsze, także jako zbiór pusty. Kontrakt rozstrzyga to
wprost: pominięcie pola znaczy „bez zmiany", a lista pusta znaczy „pamięć
wyłączona w całości" — i jest jedynym zapisem wyłączenia, ponieważ `MemoryLevel`
piątej wartości nie ma. Gdyby okno pomijało pole przy wszystkich poziomach
odznaczonych, Operator odznaczyłby cztery pola i nie wyłączyłby niczego.

Odstępstwo od globalnego promptu systemowego jest wartością logiczną, a nie
wartością `IdentityMode`, ponieważ dwie wartości pola `Agent.mode` nie są
równorzędne. Prompt systemowy ustawiany globalnie w oknie konfiguracji
obowiązuje domyślnie, a moduł Agents daje albo instrukcję dopisywaną do niego
(stan domyślny, wartość fałszywa), albo jawne odstępstwo (wartość prawdziwa), po
którym instrukcja eksperta staje się promptem systemowym. Przekład na pole
`mode` wykonuje `zrodlo-agentow.ts`, na granicy formularza i kontraktu.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-pamieci.ts

Wiersz pamięci nie prowadzi żadnego wywołania samodzielnie: zamiar oddaje oknu, które prowadzi wywołania rodziny memory. Ta sama granica obowiązuje wiersz zlecenia w pliku wiersz-zlecenia.ts, dzięki czemu oba wiersze pozostają wymiennymi elementami wykazu, a odpowiedzialność za wywołania skupia się w oknie.

Wpis o pochodzeniu model jest propozycją czekającą na decyzję Operatora. Przyjęcie zapisuje tę samą treść z pochodzeniem operator, a odrzucenie usuwa wpis. Rozróżnienie pochodzenia należy do kontraktu przez typ MemoryEntryOrigin, a nie do okna: bez niego pamięć zapisana przez model byłaby nie do odróżnienia od ustalenia wprowadzonego przez Operatora.

## budowa/klient-poprzedni/src/moduly/agents/prowenancja-rady.ts

Rada doradcy jest jawna i nie wolno pokazywać jej jako własnej odpowiedzi eksperta. Gdyby zdanie o pochodzeniu rady składał widok, pierwszy widok, który by tego zaniechał, pokazałby radę bez źródła. Dlatego zdanie o pochodzeniu powstaje jeden raz w zapisie prowenancji, a widok pobiera treść rady razem z tym zdaniem z jednego miejsca.

Kontrakt nie zawiera pola prowenancji ani komendy zapisującej ją w bazie, więc prowenancja żyje dokładnie tyle, co widok, a wykaz konsultacji mówi o tym wprost. Przeniesienie rady do instrukcji eksperta zabiera ze sobą etykietę oraz oba zdania prowenancji, dzięki czemu po wklejeniu nadal widać, skąd akapit pochodzi.

## budowa/klient-poprzedni/src/konfiguracja/indeks.ts

Okno konfiguracji jest jedno na klienta i trwa między otwarciami z tego samego
powodu, dla którego router nie porzuca widoku opuszczonej trasy: powtórne
otwarcie wraca do kategorii, na której poprzednie się zakończyło, zamiast
zaczynać od początku. Trwałość okna utrzymuje przy okazji subskrypcję zdarzenia
`config.changed`, dzięki czemu zapis dokonany w innym miejscu produktu dociera
również wtedy, gdy okno pozostaje zamknięte.

Kanał podaje się przy pierwszym otwarciu. Wywołanie z innym kanałem, właściwe
dla ponownego połączenia z rdzeniem, buduje okno od nowa, aby komendy nie szły
przez transport, którego już nie ma.

## budowa/klient-poprzedni/src/mobile/kwit-decyzji.ts

Telefon jest kanałem interwencji, a nie miejscem pracy: Operator podejmuje jedną decyzję
i odchodzi od urządzenia. Potwierdzenie wypisane w oknie ginie razem z oknem, a pytanie
o wynik wraca później, bez dostępu do pulpitu. Kwit leży w pamięci trwałej przeglądarki,
więc pierwsze, co telefon pokazuje po ponownym otwarciu, to zdanie o wykonanej interwencji
wraz z godziną, nazwą pozycji i stanem potwierdzonym przez rdzeń.

Kwit niesie zdanie rdzenia, nie widoku. Każdy krok zapisuje nazwę komendy, rozstrzygnięcie
wywołania i stan po zmianie wyjęty z odpowiedzi rdzenia — stan kolejki, zatrzymanie, źródło
konfiguracji. Kwit nieudany zapisuje treść odmowy, bo niepowodzenie także jest wiadomością,
na którą Operator czeka.

Magazyn jest wstrzykiwany, wzorem `uwierzytelnienie/sesja-bramki.ts`: domyślnie pamięć trwała
okna, a w sprawdzianie atrapa. Błąd pamięci nie zatrzymuje drogi interwencji — brak magazynu
znaczy, że kwitu nie da się zapisać, a nie że decyzja się nie wykonała. Stąd kolejność:
najpierw wywołanie rdzenia, potem kwit.

Urządzenie bez pamięci trwałej pokazuje wyłącznie potwierdzenie bieżące, a po zamknięciu
telefonu go nie ma. Taki jest stan faktyczny takiego urządzenia i okno go nie ukrywa.

## budowa/klient-poprzedni/src/moduly/design/wczytanie-pliku.ts

Treść pliku czyta klient, a nie rdzeń. Żądanie `design.asset.upload` zna
wprawdzie pole `sourcePath`, lecz ścieżka każe otworzyć plik rdzeniowi, a ten
stoi na innej maszynie niż przeglądarka, w której Operator plik wskazał.

Wymiary mierzy się wyłącznie dla rastra. Wektor oraz obraz, którego przeglądarka
nie zdekoduje, otrzymują zera znaczące brak pomiaru; oba pola wymiaru są
w kontrakcie opcjonalne, więc zasób bez nich jest poprawny i wgranie ma się
odbyć mimo nieudanego pomiaru.

Bajty odczytuje `readAsDataURL`, ponieważ oddaje zapis base64 bez ręcznego
przepisywania bajtów przez `btoa`, które na treści binarnej wymaga przejścia
przez ciąg znaków jednobajtowych i zawodzi na pierwszym bajcie powyżej 0xFF.

## budowa/klient-poprzedni/src/aplikacja/akcje-tras.ts

Powłoka środowiska ma własny pasek górny wraz z przełącznikiem motywu,
powiadomieniami i awatarem. Grupa akcji tras nie zakłada więc drugiego paska,
lecz dokłada się do grupy akcji paska istniejącego, przy jego prawej krawędzi,
przed akcjami samej powłoki. Grupa niesie tylko to, czego pasek powłoki nie ma:
oznaczenie trasy bieżącej i stan łączności z rdzeniem.

## budowa/klient-poprzedni/src/modele/zrodlo-modeli.ts

Kontrakt nie ma komendy `model.list`, ponieważ model nie jest bytem
rejestrowanym, tylko wartością danych, którą niesie kanał modelu
(`Channel.model`) albo konto (`Account.defaultModel`). Wykaz identyfikatorów
powstaje ze złożenia tych dwóch źródeł.

Wykaz służy wyłącznie jako podpowiedź do pola tekstowego. Model jeszcze
nieużywany wpisuje się identyfikatorem wprost, ponieważ pole nie jest listą
zamkniętą.

Rejestr kanałów, który nie dotarł, daje podpowiedź pustą, a nie pusty formularz.
Adresowanie osi modelu działa wtedy dalej, a nieudany odczyt zostaje odnotowany
w dzienniku przeglądarki zamiast zatrzymywać pracę.

## budowa/klient-poprzedni/src/moduly/apps/tabela-wdrozen.ts

Pusty stan wchodzi w miejsce ciała tabeli, a nie zamiast całej tabeli. Nagłówek
kolumn zostaje, ponieważ Operator ma widzieć, czego wykaz dotyczy, zanim
cokolwiek się w nim znajdzie.

Błąd wdrożenia znakuje plakietka w kolumnie stanu, bez zmiany struktury wiersza.
Stan nigdy nie opiera się na samej barwie: plakietka niesie słowo.

Gdy silnik wdrożeń wpisze do pola `logRef` zdanie o tym, dlaczego przebieg się nie
powiódł, tabela pokazuje jego treść tak, jak przyszła, i nazywa pole, z którego
pochodzi. Samo słowo `failed` zostawiałoby Operatora bez wyjaśnienia, które
przyszło w tej samej ramce.

## budowa/klient-poprzedni/src/moduly/agents/indeks.ts

Powłoka zna z tego katalogu jedną czynność: wytwarza moduł wytwórnią
`utworzModulAgents`, przyjmując kanał rdzenia, stawia oddany element w obszarze
roboczym i wywołuje czynność wczytania z identyfikatorem sesji. Kolejność jest
wiążąca, ponieważ wczytanie żąda danych sesji, a element musi już stać
w dokumencie, gdy odpowiedź wraca.

Moduł nie osadza się sam w dokumencie i nie zna powłoki, dzięki czemu te same
okna wchodzą zarówno w obszar roboczy powłoki, jak i w podgląd sprawdzianu.

Kontekst testowanego eksperta jest wystawiony w interfejsie katalogu, ponieważ
sięga po niego powłoka, a nie moduł. Zmianę testowanego agenta rozgłasza plik
`aplikacja/ulotnosc-okna.ts`, na którym stoi okno rozmowy sesji.

## budowa/klient-poprzedni/src/aplikacja/zamiary-pulpitu.ts

Pulpit sam niczego nie wysyła: nadaje zamiar wyrażony nazwami kontraktu, a kopertę
komendy składa dopiero ta warstwa. Dzięki temu pulpit zna wyłącznie słownik zamiarów,
a wiedza o komendach i o ich wymaganych polach zostaje w jednym miejscu.

Dwa zamiary nie mają dziś drogi przez kontrakt i żaden z nich nie jest udawany.
Podniesienie priorytetu nie ma odpowiednika ani w `QueueAction`, ani wśród komend.
Przekazanie kontekstu wymaga kompletu `ContextBundle`, którego zamiar kolejki nie
niesie, a którego nie wolno złożyć po stronie klienta z wartości domyślnych.

W obu przypadkach warstwa wystawia komunikat mówiący wprost, czego brakuje, zamiast
wysyłać komendę niekompletną albo milczeć.

## budowa/klient-poprzedni/src/mission-control/etykiety-pulpitu.ts

Plik ma jedną odpowiedzialność: przekłada wartości kontraktu na słowo widoczne
dla Operatora. Widok nie zna literałów kontraktu i o nazwę pyta wyłącznie tutaj,
dzięki czemu zmiana słownictwa pulpitu nie sięga do samych widoków.

Każdy słownik jest rekordem zupełnym po dziedzinie swojego typu. Dopisanie
wartości do kontraktu przerywa więc kompilację tego pliku, i jest to skutek
zamierzony: nowy stan ma dostać polską nazwę, zamiast wypaść z widoku po cichu
albo pokazać się Operatorowi surowym literałem kontraktu.

Etykieta braku źródła zastępuje liczbę wszędzie tam, gdzie kontrakt miary dziś
nie niesie. Pulpit nie pokazuje wartości wymyślonej — brak źródła danych jest
nazywany wprost.

## budowa/klient-poprzedni/src/konfiguracja/stany-odczytu.ts

Komunikat blokowy nie znika po naciśnięciu przycisku ponowienia. Naciśnięcie wyzwala kolejny odczyt
i zostawia zdanie na miejscu, dopóki sytuacja nie ustanie; zdejmuje je dopiero odczyt zakończony
powodzeniem.

Stan pustego katalogu ustawień należy do panelu kategorii, a nie do tego pasa, ponieważ mówi
o zawartości katalogu, a nie o przebiegu odczytu.

Zawężenie zależności pasa do fazy odczytu i powodu niepowodzenia czyni z niego jeden byt dla obu
odczytów okna konfiguracji: katalogu ustawień oraz konfiguracji obowiązującej. Stan konfiguracji
spełnia ten kształt bez żadnej zmiany.

## budowa/klient-poprzedni/src/asystent-plywajacy/indeks.ts

Warstwa wystawia na zewnątrz jedno wywołanie `zaczepAsystenta`, przyjmujące
kanał oraz opcjonalne posunięcia i identyfikator klienta. Pozostałe części,
czyli favikon, dymek, stan rozmowy i dostępność głosu, są wnętrzem warstwy
i nie mają powodu być widoczne z aplikacji.

Pierwszym wyjątkiem są `czyGlosDziala` oraz `brakujaceOgniwa`. Stan kanału
głosowego mówi o produkcie, a nie o widoku, więc jest potrzebny każdej warstwie
budującej obsługę głosu.

Drugim wyjątkiem jest `sprawca-zdarzenia.ts`. Odczyt pól `actor`
i `actorClientId` z koperty przydaje się poza tą warstwą, ponieważ
`aplikacja/zrodlo-posuniec.ts` rozpoznaje sprawcę odciskiem okna, mając prawdę
w kopercie. Wystawienie tego odczytu w punkcie zbiorczym zapobiega powstaniu
drugiej, odmiennej realizacji tej samej reguły.

## budowa/klient-poprzedni/src/mission-control/sekcja-aktywnosci.ts

Nazwy miar w kaflach są dosłowne, ponieważ mają nie dopuścić do pomylenia wysycenia kanałów modelu z obciążeniem maszyny. Miara bez źródła w kontrakcie pokazuje kreskę i etykietę braku źródła zamiast wartości zastępczej, aby pulpit nie sugerował danych, których rdzeń nie dostarcza.

Pas decyzji nie należy do tej sekcji. Sekcja udostępnia jedynie miejsce montażu pod rzędem kafli, a sam pas dokłada moduł mission-control.ts, który zna kolejność elementów pulpitu.

## budowa/klient-poprzedni/src/moduly/apps/wykaz-plikow-warsztatu.ts

Wykaz pokazuje wszystkie pliki warsztatu okna, dzięki czemu widać, że wpisana ścieżka należy już do istniejącego pliku, zanim zapis go nadpisze. Wiersz jest przyciskiem, ponieważ wiersz nieklikalny kazałby przepisywać ścieżkę ręcznie, a przepisana ścieżka bywa w rzeczywistości innym plikiem: klucz warsztatu tworzą okno, warstwa i ścieżka rozumiane co do znaku.

Wykaz niczego nie streszcza. Rozmiar i wersja idą tak, jak oddał je rdzeń, a brak pola zostaje nazwany brakiem, zamiast zostać zastąpiony wartością zerową, która wyglądałaby jak wynik pomiaru.

## budowa/klient-poprzedni/src/komponenty/dymek.ts

Znak zapytania jest przyciskiem, nie ozdobą: naciśnięcie prowadzi do niego
ognisko, a ognisko pokazuje objaśnienie. Treść czytają technologie wspomagające
z atrybutu `aria-label` znaku, więc dymek nie potrzebuje identyfikatora i nie
zderza się z dymkami innych okien.

Klasy wywołującego godzą jedną fabrykę z wyglądem właściwym dla miejsca
wywołania. Klasa `.dn-btn-ikona` to kwadrat 32 px, na dotyku 40 px, o narożniku
`--dn-r-sm`, z przezroczystym obrysem i wskaźnikiem `pointer`. Znaki modułów to
obwódki 14–16 px o narożniku `--dn-r-pill`, z widocznym obrysem i wskaźnikiem
`help`. Powłoka niesie z kolei zachowanie w rzędzie — `flex: none`, odstęp
i `vertical-align` — którego `.dn-tooltip` nie zna.

## budowa/klient-poprzedni/src/moduly/library/filtr-wersji.ts

Filtr działa w całości po stronie okna, ponieważ komenda wykazu wersji
biblioteki nie przyjmuje pola zawężającego, a jej odpowiedź niesie komplet pól,
po których dokumentacja każe zawężać, czyli sprawcę zmiany i etykietę wersji.
Wysłanie w tym celu drugiego odczytu byłoby pytaniem o to, co okno już ma.

Sprawca zmiany jest napisem, a nie wyliczeniem kontraktu: rdzeń zapisuje w polu
sprawcy to, co poda wołający. Pozycje zawężające do zmian Operatora oraz do
zmian modelu dopasowują więc po zawartości napisu, a wersja o sprawcy
nienazwanym wchodzi wyłącznie do pozycji obejmującej wszystkie wersje.
Zgadywanie po stronie okna kazałoby takiej wersji trafić do jednej z dwóch grup
bez podstawy.

## budowa/klient-poprzedni/src/moduly/library/sterowanie-widoku.ts

Stopień obrotu i krotność powiększenia trafiają do zbioru danych elementu, a cały
wygląd wynikający z tych wartości bierze się z arkusza stylów modułu. Plik widoku
nie zna zatem żadnej wartości wizualnej: nie podaje wymiarów, odstępów ani barw,
a zmiana wyglądu obrotu i powiększenia odbywa się wyłącznie w arkuszu.

## budowa/klient-poprzedni/src/modele/warstwy-tozsamosci.ts

Warstwa mówi, jak krytyczna jest treść nakładki: konstytucja stoi wyżej niż
profil roli, a profil wyżej niż ekspertyza. Ta kolejność rządzi zarówno układem
wykazu kategorii, jak i porządkiem złożonego promptu. Wartości pochodzą wyłącznie
ze stałych kontraktu.

Tryb rozstrzyga los ustawień fabrycznych: `ZASTAP` podmienia prompt fabryczny
w całości, `DOLACZ` dokłada treść do niego. Wybór decyduje o tym, którym
argumentem uruchomienia pojedzie nakładka, więc każdy przełącznik trybu w tej
sekcji podaje skutek wprost.

## budowa/klient-poprzedni/src/komponenty/odmowa.ts

Kod błędu zostaje przy wiadomości z zamysłem. Kody `validation_failed`
i `not_found` znaczą dla Operatora co innego, a treść wiadomości bywa dla obu
tym samym zdaniem, więc zdanie bez kodu odbierałoby Operatorowi rozróżnienie.
Brak wiadomości nie zostawia pustki: idzie wtedy sam kod, a gdy nie ma i kodu,
zdanie nazywa milczenie rdzenia wprost, żeby Operator nie patrzył na pusty
prostokąt.

Kształt powodu odmowy jest luźniejszy od kształtu `ErrorInfo` z kontraktu,
ponieważ moduły dostają powód z pola `blad` wyniku komendy, które bywa puste,
a pojedyncze okna trzymają własne, częściowe zapisy odmowy.

Nazwa nieznanej komendy obsługuje przypadek osobny, otwarty dla każdego modułu:
rdzeń odpowiedział zdarzeniem mówiącym, że komendy nie zna. Nie jest to awaria
wykonania, a Operator musi widzieć różnicę między usterką a czynnością, której
rdzeń jeszcze nie umie, dlatego nazwa typu idzie na początek zdania.

Zdanie składane z powodu wziętego w całości istnieje po to, żeby wywołanie
po odmowie komendy nie powtarzało w każdym miejscu pary pól `code` i `message`.
Para rozjeżdża się przy przepisywaniu, ponieważ łatwo podać kod z jednego
wyniku, a wiadomość z drugiego.

## budowa/klient-poprzedni/src/moduly/browser/przyciski-browser.ts

Panele pomocnicze i wiersze wykazu modułu Browser potrzebują tego samego:
zbudować przycisk i podpiąć nasłuch naciśnięcia. Czynność ta stoi w jednym
pliku, żeby cztery pliki modułu nie rozjechały się przy pierwszej poprawce.
Nazwa funkcji budującej zawiera wyraz „przycisk”, ponieważ po nim rozpoznają
kontrolkę narzędzia zestawiające etykiety czynności okien; opakowanie nazwane
inaczej chowa etykietę przed takim zestawieniem. Nazwy odmian zastępują napis
klasy powtarzany przy każdym wywołaniu: `glowny` stoi tam, gdzie panel ma jedną
czynność wiodącą, na przykład zapis albo dodanie, `zarys` przy pozostałych
czynnościach panelu, `duch` przy pozycjach wykazu, gdzie przycisków w jednym
wierszu jest wiele.

## budowa/klient-poprzedni/src/modele/panel-kont.ts

Panel pamięta wyłącznie, które konto jest już wypełnione w formularzu. Bez tej pamięci każde zdarzenie zmiany konta kasowałoby treść właśnie wpisywaną przez Operatora, dlatego wypełnienie formularza następuje przy zmianie konta czynnego, a nie przy każdym przeliczeniu stanu.

Ograniczenie rodzaju kont jest częścią żądania wykazu kont kierowanego do rdzenia, a nie filtrem zakładanym w kliencie. Rejestr kont bywa długi, a kolejność rotacji zna rdzeń, więc wybór rodzaju musi trafić do żądania, żeby wykaz pozostał zgodny z rejestrem.

## budowa/klient-poprzedni/src/moduly/diagnostics/okno-logs-viewer.ts

Port Diagnostyka jest w rdzeniu odbiorcą odmów wykonania komend. Każda odmowa
dostaje wpis dziennika, w którym pole źródła niesie nazwę odrzuconej komendy —
także wtedy, gdy zapis wiersza błędu się nie powiedzie i wpis dziennika jest
jedynym śladem tej odmowy. Dlatego filtr źródła i wzorzec wyszukiwania stoją
na pierwszym miejscu paska narzędzi okna.

Przy włączonym przełączniku wyrażenia regularnego okno próbuje zbudować wzorzec
z wpisanego pola, zanim żądanie pójdzie do rdzenia: błąd składni zatrzymuje
żądanie i mówi to wprost, zamiast wysyłać wzorzec, który rdzeń odrzuciłby po
cichszej stronie swojej.

Zakres czasu jest wspólny całemu modułowi. Okno nie prowadzi własnych pól
zakresu: czyta zakres ze stanu modułu i nasłuchuje jego zmiany, żeby
przestawienie zakresu w Diagnostics Center przeliczyło wykaz bez osobnej
czynności okna.

Pola informujące o obcięciu wyniku i o łącznej liczbie spełniających warunki
trafiają do zdania potwierdzenia, ponieważ wykaz ucięty granicą bufora bez
wiedzy czytelnika jest kłamstwem tej samej rodziny co pusty wykaz pokazany
przy odmowie.

Odświeżanie na żądanie znaczy tu odczyt wywołany czynnością, nie strumień na
żywo. Kontrakt nie niesie zdarzenia nowego wpisu dziennika, więc okno nie
odpytuje rdzenia w pętli — odczyt biegnie na żądanie Operatora i na zmianę
zakresu wspólnego modułu.

Rozwijana lista poziomu wpisu korzysta z biblioteki komponentu drzewa menu,
osadzonego przez moduł wyboru z menu w warstwie okien aplikacji, zamiast
z natywnego elementu wyboru.

## budowa/klient-poprzedni/src/moduly/browser/wiersz-notatki.ts

Powiązanie ze źródłem jest widoczne, a nie domyślne. Notatka bez identyfikatora źródła mówi o tym
wprost, ponieważ powiązanie zakłada się z wykazu okna, a jego brak jest informacją, a nie pustką do
przemilczenia.

Klasyfikacja stoi przy wierszu jako lista wyboru, a nie jako ikona. Rodzaj notatki ma być czytelny
bez najeżdżania na znak i bez rozróżniania barw, a zmiana klasyfikacji ma być jednym gestem.

## budowa/klient-poprzedni/src/aplikacja/gospodarz-dokumentu.ts

Gospodarz ma jedną odpowiedzialność: przygotowuje dokument oraz jeden element,
do którego router wstawia widoki tras. Nie zna rdzenia ani żadnego widoku i nie
wie nawet, ile tras ma aplikacja.

Brak kontenera przygotowanego przez `index.html` nie zatrzymuje uruchomienia.
Gospodarz trafia wtedy do ciała dokumentu, a aplikacja rusza normalnie.

Tytuł dokumentu składa nazwę projektu wziętą z opisu okna z nazwą produktu,
ponieważ nazwa projektu pochodzi z konfiguracji budowania i różni się między
osadzeniami, natomiast nazwa produktu jest stała.

## budowa/klient-poprzedni/src/moduly/agents/sklad-zespolu.ts

Kolejność zaznaczania jest treścią, nie porządkiem wyświetlania. Pole `agentIds`
niesie ekspertów w kolejności nadanej przy wyborze; składanie tablicy przejściem
po bibliotece brałoby kolejność z wykazu ekspertów, więc zespół wracałby po zapisie
w innym porządku niż złożony. Wybór trzyma zatem własną listę i dopisuje na jej koniec.

Zespół wczytany może wskazywać eksperta spoza bieżącego wykazu: zarchiwizowanego albo
odciętego frazą zawężającą. Taki wpis zostaje w składzie i pokazuje się wierszem
nazywającym brak, żeby ponowny zapis nie okroił zespołu po cichu.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-zaplecza.ts

Odczyt `window.list` podaje okna komunikacji sesji. Komenda `assistant.voice.command`
wymaga pola `windowId`, a moduł okien sam nie zakłada: bierze to, które rdzeń
przypisał modułowi.

Odczyt `action.list` podaje katalog akcji zasięgu modułu. Siatka szybkich akcji Voice
Console jest widokiem tego katalogu, a nie listą zaszytą w kliencie: nowa akcja jest
nowym wierszem rdzenia, nie zmianą kodu.

Odczyt `component.list` podaje komponenty własne rodzaju `assistant`, czyli profile
asystenta ze strefy drugiej strony głównej. Jest to jedyny odczyt kontraktu, którym da
się wypełnić selektor profilu Voice Console. Komponenty wyłączone wchodzą do wykazu:
profil wyłączony w strefie drugiej wciąż istnieje w rdzeniu i wskazanie go należy do
Operatora, a nie jest pomyłką okna.

Odczyt `speech.availability.get` mówi, czy silnik mowy stoi na tej maszynie. Brak
silnika jest odpowiedzią, a nie awarią, więc okno mówi o nim wprost.

Wywołanie `speech.transcribe` rozpoznaje nagranie wskazane ścieżką na maszynie silnika.
Pusta transkrypcja jest wynikiem prawidłowym, gdy pole `processed` niesie prawdę, ponieważ
nagranie mogło zawierać ciszę albo szum; sprawdzian kształtu pyta zatem o obecność pól,
a nie o niepustą treść.

## budowa/klient-poprzedni/src/moduly/design/panel-warstw.ts

Podział ról między kanwą a panelem jest następujący: kanwa pokazuje położenie
warstwy, panel pokazuje porządek, blokadę i adnotację. Adnotacja jest polem
notatki warstwy w kontrakcie, więc zostaje w kompozycji na trwałe.

Blokada dotyczy wyłącznie położenia warstwy na kanwie. Wiersz warstwy
zablokowanej pozostaje klikalny i edytowalny, a stan blokady niesie atrybut
danych wiersza.

## budowa/klient-poprzedni/src/moduly/apps/wykaz-komponentow-warstwy.ts

Okno Frontend Workspace i okno Backend Workspace korzystają z tego samego zbioru komponentów pochodzącego z Architecture Designera i różnią się wyłącznie tym, które rodzaje komponentów do nich należą. Dlatego oba okna składa jeden moduł wykazu, sparametryzowany warstwą warsztatu.

Przypisanie rodzaju komponentu do warstwy jest wyborem klienta, a nie kontraktu. Ani typ AppWorkspaceLayer, ani typ AppComponentKind nie łączą tych dwóch zbiorów, więc podział został wypisany wprost w stałej RODZAJE_WARSTWY, aby stał w jednym miejscu i nie rozchodził się po oknach.

## budowa/klient-poprzedni/src/moduly/browser/warstwy-widocznosci.ts

Ujawnianie nie jest bramą. Element zwinięty nie jest zablokowany — jest schowany, a każda
jego czynność zostaje osiągalna wyzwalaczem albo skrótem. Dlatego skrót warstwy czwartej
działa również wtedy, gdy tryb administracyjny jest wyłączony: tryb decyduje o obecności
wyzwalacza w pasku, nie o dostępie do funkcji.

Skrót podpina się do elementu modułu, a nie do dokumentu. W powłoce stoi obok siebie kilka
modułów, a nasłuch założony na dokumencie odpowiadałby także na klawisze wciśnięte w cudzym
oknie. Skrót działa więc wtedy, gdy fokus stoi wewnątrz modułu. Powłoka gospodarza może
przechwycić kombinację przed stroną — wtedy pozostaje wyzwalacz w pasku.

Licznik pozycji rozszerzenia jest polem pominiętym, gdy rozszerzenie niczego nie zlicza.
Wartość zero jest nieodróżnialna od braku licznika, więc pominięcie niesie własne znaczenie.

## budowa/klient-poprzedni/src/aplikacja/godlo-aplikacji.ts

Znak marki ma barwy własne, więc nie dziedziczy barwy tekstu i nie zmienia się
wraz z motywem; odmianę znaku dobiera podłoże, a pasek aplikacji jest atramentowy
w obu motywach, stąd podłoże ciemne. Droga przez `elementGodla` niesie ponadto
próg odmiany uproszczonej: przy rozmiarze szesnastu pikseli i mniejszym znak
przechodzi na jeden grot.

Nazwa produktu jest stała, a nazwa projektu przychodzi z opisu okna, czyli
z konfiguracji budowania utrwalonej w `okno-komunikacji/opis-okna.ts`. Blok
godła nie jest kontrolką i nią nie udaje: nie prowadzi nigdzie i nie zmienia
kursora na wskazujący.

## budowa/klient-poprzedni/src/ladowanie/indeks.ts

Poza punktem wejścia katalog nie wystawia niczego. Arkusz stylu sceny i postać
jej bryły pozostają sprawą wewnętrzną, więc zmiana układu ekranu wejścia nie
sięga warstw, które scenę stawiają.

## budowa/klient-poprzedni/src/modele/znak-wykazu.ts

Klasa wyglądu pochodzi z biblioteki — jest to `dn-plakietka` wraz z jej
odmianami — a klasa miejsca należy do wykazu, w którym plakietka stoi. Dlatego
obie stoją w osobnych parametrach: wygląd znaku i jego położenie w wierszu są
dwiema niezależnymi decyzjami wywołującego.

Plik nie wymienia żadnej barwy. Barwa plakietki wynika wyłącznie z odmiany klasy
bibliotecznej, więc zmiana palety nie sięga tego pliku.

## budowa/klient-poprzedni/src/modele/pola-konta.ts

Zestaw pól jest oddzielony od formularza, ponieważ formularz odpowiada za rozgałęzienie
zapisu na dwie komendy kontraktu i za to, co wolno nadpisać w polu już wypełnionym.
Trzymanie obu odpowiedzialności w jednym miejscu wiązałoby kształt danych z przebiegiem
zapisu.

Opis poświadczenia należy do pola, a nie do widoku: kontrakt nie zwraca poświadczenia
żadną komendą, więc pole konta, które poświadczenie ma, i tak wygląda na puste. Opis
jest jedynym miejscem, w którym ten stan zostaje nazwany.

## budowa/klient-poprzedni/src/moduly/browser/pozycje-pokrycia.ts

Powód nieczynnej pozycji rozstrzyga byt wspólny wszystkich modułów, stojący
w pliku `moduly/pokrycie-komend.ts`. Byt pyta rdzeń o wykaz jego komend
powitaniem `connection.hello` i rozróżnia trzy stany: kontrakt danej komendy
nie zna, kontrakt ją zna, lecz rdzeń nie ma uchwytu, albo rdzeń ma uchwyt,
a okno jeszcze go nie wywołuje. Zdanie wpisane w moduł na sztywno mówiłoby
nieprawdę od chwili, w której rdzeń dostanie obsługę, i nic nie wymuszałoby
jego zdjęcia, ponieważ nic nie łączyłoby go z rdzeniem.

Przycisk nieczynnej pozycji nie jest wygaszony i nie milczy: naciśnięcie daje
dymek z powodem oraz wiersz odpowiedzi okna, a żądanie do rdzenia nie wychodzi
i nic nie udaje wykonania. Znaczniki `data-brak-komendy` oraz `data-pokrycie`
niesie kontrolka biblioteczna i po nich sięgają sprawdziany produktu.

Powód jest czytany z atrybutu `title` dopiero przy naciśnięciu, tak samo jak
robi to dymek kontrolki. Byt przerysowuje ten atrybut po każdej zmianie wykazu,
więc domknięcie na treści z chwili budowy podawałoby stan odczytu w toku długo
po odpowiedzi rdzenia.

## budowa/klient-poprzedni/src/moduly/multitasking/obsada-rol.ts

Rola okna pochodzi z pola roli w kontrakcie, nie z domysłu. Koordynatorem jest
okno roli koordynatora, a wykonawcą okno roli wykonawcy, którego wskazanie okna
zarządzającego trafia w tego koordynatora — to jest cała definicja pętli.
Znakowanie ról po tytule albo po kolejności rozjechałoby rdzeń z widokiem przy
pierwszym przepięciu wykonawcy pod innego koordynatora.

Analityk stoi poza pętlą. Okno analizy wyników nie jest wykonawcą: nie kończy
tury, która miałaby wybudzić koordynatora, więc jego rolą kontraktową jest rola
samodzielna. Odróżnia go wskazanie zapisane na poziomie sesji, ponieważ kontrakt
nie ma czwartej wartości roli okna, a wymyślenie jej po stronie klienta
rozjechałoby wykaz ról z bazą.

Kody okien są bezmodułowe: rejestr okien operacyjnych niesie kod rozmowy
koordynatora, dwa kody rozmowy wykonawcy oraz kod analizy wyników, każdy bez
przedrostka modułu. Przynależność do modułu niesie osobna macierz, w której tych
czterech okien nie ma — liczą się tam jako okna pozamodułowe. Kod zapisany
z przedrostkiem nie trafiłby w żaden wiersz rejestru, a nawigacja po kodzie okna
nie miałaby dokąd skoczyć.

## budowa/klient-poprzedni/src/aplikacja/trasy.ts

Aplikacja ma trzy widoki najwyższego rzędu. Trasa `strona-glowna` prowadzi do
Centrum dowodzenia i stanowi wejście do produktu. Trasa `srodowisko` otwiera
powłokę środowiska wraz z kartami sesji i modułami. Trasa `pulpit` otwiera
Mission Control jako widok stojący obok strony głównej, a nie zagnieżdżony w niej.

Nazwa trasy pełni dwie role naraz: jest kluczem katalogu i wartością zapisywaną
w adresie dokumentu. Dzięki temu odświeżenie strony wraca do tego samego widoku,
a nie na początek przepływu.

## budowa/klient-poprzedni/src/moduly/apps/indeks.ts

Moduł nie osadza się sam w dokumencie: oddaje element, a warstwa składająca
rozstrzyga, gdzie go postawić. Dzięki temu te same okna wchodzą i w obszar
roboczy powłoki, i w stanowisko sprawdzianu.

Pole `zamknij` podpina `rozlacz()` złożenia modułu. Bez niego subskrypcje stanu
`apps.build.changed` oraz `progress.changed`, a wraz z nimi subskrypcje okien
pomocniczych, żyją dalej po zejściu modułu ze sceny. Pole `zamknij` jest
w typie `WidokModulu` opcjonalne, więc kompilator jego braku nie zgłosi.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-podagentow.ts

Źródło jest jednym z trzech źródeł modułu, obok `zrodlo-biegu.ts` niosącego bieg
pracy i `zrodlo-okien.ts` niosącego obsadę ról. Stoi osobno od biegu pracy,
ponieważ bieg porusza oknami ról — tura, przekazanie, kolejka etapu, telemetria
— a podagent jest bytem pod oknem wykonawcy, o którym kontrakt nie rozstrzyga,
czy jest oknem, procesem, czy pozycją kolejki.

Telemetria zostaje w biegu pracy: panel bierze `monitor.status` ze źródła biegu,
które tę komendę już niesie. Drugie wywołanie tej samej komendy w drugim źródle
byłoby kopią bez powodu.

Powołanie oddaje podagentów tak, jak założył je rdzeń, a nie tak, jak prosił
formularz: liczba powołanych pochodzi z odpowiedzi, ponieważ rdzeń ma prawo
powołać ich mniej, a górna granica piętnastu należy do rdzenia, nie do widoku.
Zawężenie wykazu po stanie również wykonuje rdzeń — żądanie niesie pole
`status`, więc sito po stronie klienta byłoby drugim, rozjeżdżającym się sitem.

Zebranie wyników oddaje komplet podagentów wraz z polem `complete`, którego
panel nie zgaduje ze stanów pojedynczych podagentów, tylko bierze z odpowiedzi.
Zdarzenie zmiany stanu rdzeń rozgłasza przy powołaniu, wejściu w bieg,
zakończeniu i zatrzymaniu, więc wykaz nadąża bez ręcznego odpytywania.

Zatrzymanie idzie komendą `subagent.stop` wprost, a nie obejściem przez
`queue.action` na kolejce o nazwie podagenta: podagent powołany bez wpiętej
kolejki nie miałby wtedy czego zatrzymać. Odpowiedź rozróżnia `stopped` od
`notRunning`, ponieważ podagent już zakończony nie jest błędem, tylko innym
stanem, o którym panel mówi osobno.

## budowa/klient-poprzedni/src/aplikacja/ustawienia-okna-sledzone.ts

Migawka ustawień wraz z wykazem różnic mówi, co się zmieniło. Bez nich pas
meldunkowy podawałby samo zdanie o zmianie okna przy każdym dotknięciu, zamiast
nazwać dobrany model, rolę albo zasięg wykonania.

Odcisk okna mówi, czy zmiana jest własna: zdarzenie przynoszące dokładnie ten
stan okna, co świeża odpowiedź na własną komendę, jest skutkiem tej komendy.
Do odcisku nie wchodzi czas zmiany zapisany w polu `updatedAt`, ponieważ rdzeń
nadaje go osobno w zdarzeniu i osobno w odpowiedzi; doklejenie czasu sprawiłoby,
że własna czynność nigdy nie zrównałaby się sama ze sobą.

Plik nie zna ani widoku, ani kanału. Niesie sam przekład pól okna na wartości
porównywalne oraz na zdania przeznaczone dla czytającego.

## budowa/klient-poprzedni/src/moduly/diagnostics/niepowodzenie-odczytu.ts

Wynik oznaczony jako nieudany powstaje w warstwie protokołu na więcej niż jeden sposób, a tylko jeden z nich jest odmową rdzenia. Pierwsza droga to odmowa rdzenia, czyli koperta ze statusem błędu i opisem błędu z kontraktu. Druga droga to odpowiedź nieczytelna: rdzeń odpowiedział powodzeniem, lecz treść nie ma kształtu przewidzianego kontraktem, więc sprawdzanie kształtu odpowiedzi zamienia ją na niepowodzenie z kodem nieudanej walidacji.

Nazwanie odpowiedzi nieczytelnej odmową byłoby oskarżeniem rdzenia o czyn, którego nie popełnił, ponieważ rdzeń odpowiedział i odpowiedział powodzeniem, a nie zrozumiał go klient. Okno nie ma prawa twierdzić o braku, którego rdzeń nie orzekł, tak samo jak nie ma prawa zamienić odmowy w pustkę.

Stan odpowiedzi udanej bez treści pełni rolę straży, a nie drogi osiągalnej w obecnym układzie. Wynik cząstkowy oddaje niepowodzenie bez pola błędu, gdy odpowiedź jest udana, lecz pusta, jednak do przeniesienia taka odpowiedź nie dochodzi, bo sprawdzanie kształtu przechwytuje ją wcześniej i okno mówi wówczas zdanie o odpowiedzi nieczytelnej. Stan zostaje w słowniku na wypadek zmiany warstwy protokołu.

Zdania stoją w słowniku, a nie w oknach, z tego samego powodu, dla którego stan treści stoi w osobnym komponencie: cztery okna mówiące o tej samej ciszy czterema różnymi zdaniami rozjeżdżają się przy pierwszej poprawce.

## budowa/klient-poprzedni/src/moduly/apps/panel-akcji.ts

Panel akcji modułu powstaje wyłącznie z rejestru rdzenia. Nowa akcja modułu jest nowym wierszem rejestru, a nie zmianą kodu klienta, dlatego w tym pliku nie ma ani jednej pozycji katalogu zapisanej na stałe.

Wiersz rejestru niesie nazwę komendy, którą akcja wywołuje, jednak generycznej drogi wywołania window.action nie woła żaden widok klienta. Naciśnięcie przycisku pozycji nazywa więc komendę wiersza i stwierdza, że droga generyczna czeka na konsumenta. Takie rozwiązanie przyjęto zamiast milczącego przycisku oraz zamiast pozorowania wykonania akcji.

## budowa/klient-poprzedni/src/mission-control/indeks.ts

Katalog pozwala powłoce zamontować pulpit i wpiąć źródło danych rdzenia jednym
ruchem:

    import {
      pustyKomplet,
      utworzMissionControl,
      utworzZrodloPulpitu,
    } from './mission-control/indeks';

    const pulpit = utworzMissionControl(pustyKomplet());
    const zrodlo = utworzZrodloPulpitu(kanal, (dane) => pulpit.odswiez(dane));
    zrodlo.uruchom();

Komplet danych pulpitu pochodzi wyłącznie z odczytów `session.list`,
`window.list` i `channel.list` oraz ze zdarzeń `session.changed`,
`window.changed`, `queue.changed` i `progress.changed`. Do pierwszego odczytu
pulpit pokazuje stany puste zamiast wartości zastępczych. Sam katalog nie
zawiera logiki: jest wykazem tego, co warstwa wystawia na zewnątrz.

## budowa/klient-poprzedni/src/moduly/browser/wykazy-zebranego.ts

Zaciągnięcie składa się z dwóch odczytów rdzenia: `browser.source.list`
i `browser.note.list`. Wynik obu przekłada się na zbiór zebranego w sesji.

Odmowa „okna nie ma" nie jest błędem okna. Rdzeń odpowiada kodem `not_found`,
gdy moduł Browser nie widział jeszcze wskazanego okna, a to normalny stan okna
świeżo otwartego, w którym nikt niczego nie zebrał. Rozpoznanie kończy się wtedy
stanem `nietkniete`: wykaz zostaje pusty, powód jest podany, ale okno nie wchodzi
w stan błędu. Każda inna odmowa jest błędem i mówi swoim powodem. Okno nieznane
zgłasza się przy obu wykazach naraz, więc jest to jedno rozpoznanie, a nie dwa
osobne niepowodzenia.

Stanu `odczyt` nie wpisuje tu nikt: wynik oddawany jest dopiero po obu komendach,
więc stan oczekiwania ustawia wołający w `stan-przegladania.ts` przed czekaniem.
Wartość stoi mimo to w jednym typie z pozostałymi, ponieważ panele czytają jeden
stan wykazu, a nie dwa niezależne.

## budowa/klient-poprzedni/src/moduly/multitasking/stany-okna.ts

Stan błędu jest rozłączny ze stanem pustki. Okno roli, które po odmowie odczytu
pokazałoby wyzerowany licznik zamiast błędu, sugerowałoby zatrzymany bieg tam,
gdzie bieg trwa, a odmowa odczytu i brak pozycji w wykazie są dla Operatora
dwiema różnymi wiadomościami.

## budowa/klient-poprzedni/src/moduly/apps/okno-frontend-workspace.ts

Plik jest wiązaniem, nie drugim widokiem: formularz warsztatu stoi raz,
w `okno-warsztatu.ts`, a stąd dostaje wyłącznie stan produktu i opis warsztatu
frontendu.

Kod okna nie pada tu wprost napisem. Ramę modułu Apps woła wspólne okno
warsztatu, biorąc kod ze zmiennej opisu, dzięki czemu wykaz `KODY_OKIEN`
pozostaje jedynym miejscem, w którym kody okien tego modułu są wypisane.

## budowa/klient-poprzedni/src/ikony/zrodla/srodowiska.ts

Nazwy i kolejność pozycji tego pliku odpowiadają grupie środowisk w wykazie
`ikony/manifest.json`, dzięki czemu wykaz ikon i wykaz źródeł dają się porównać
pozycja po pozycji.

Emblematy środowisk są znakami własnymi budowy i nie pochodzą z zestawu Lucide,
więc nie da się ich podmienić aktualizacją biblioteki ikon.

## budowa/klient-poprzedni/src/moduly/developer/stany-okna.ts

Stan błędu jest rozłączny ze stanem pustki. Okno, które po odmowie odczytu
pokazuje puste drzewo, mówi Operatorowi, że katalog jest pusty, zamiast
powiedzieć, że drzewa nie udało się odczytać.

Pustka bywa w tym module stanem poprawnym: repozytorium bez zmian do
zatwierdzenia nie jest usterką, więc pusty wykaz nie może być pokazywany jako
niepowodzenie.

## budowa/klient-poprzedni/src/aod/sekcja-rozmowy.ts

Ognisko okien zna rdzeń i mogło się ono zmienić między odczytem stanu a wysyłką,
więc sekcja nie podstawia w pole `windowId` okna odczytanego wcześniej. Puste pole
oznacza wprost „okno ogniskowane w chwili przyjęcia wiadomości". Odpowiedź niesie
`messageId` oraz `windowId` okna, które wiadomość przyjęło, i meldunek pokazuje oba.

Przycisk wysyłki pozostaje czynny także przy pustym polu treści. Pustą treść ocenia
rdzeń i to on zwraca odmowę `validation_failed`; blokada po stronie nakładki
rozdzielałaby ocenę wiadomości na dwa miejsca.

## budowa/klient-poprzedni/src/dostepy/okno-dostepow.ts

Sekcja dostępów sama nie zakłada, gdzie zostanie osadzona: oddaje element. Rama
okna jest jednak potrzebna od razu, ponieważ pierwszym gospodarzem sekcji jest
listwa ustawień Centrum dowodzenia, a ta otwiera widoki jako okna modalne, tak
samo jak okno konfiguracji.

Okno stoi na natywnym elemencie okna dialogowego, więc warstwę tła, stos okien
i zamknięcie klawiszem ucieczki daje przeglądarka, a nie własna nakładka. Wygląd
bierze z biblioteki komponentów, z klasy okna modalnego, więc plik nie zna ani
jednej barwy.

Okno otwiera się natychmiast, przed odpowiedzią rdzenia, a wykazy dojeżdżają do
niego odpowiedzią.

## budowa/klient-poprzedni/src/mission-control/kolumna-procesow.ts

Wykaz procesów w tle buduje się wyłącznie ze zdarzeń telemetrii, ponieważ kontrakt nie ma komendy
odczytu procesów bieżących. Do nadejścia pierwszego zdarzenia kolumna pokazuje stan pusty, który
nazywa tę okoliczność wprost, zamiast udawać brak procesów.

Stan procesu jest wypisany słowem obok paska postępu. Barwa paska wspiera odczyt, ale go nie
zastępuje, dzięki czemu wykaz pozostaje czytelny bez rozróżniania barw.

## budowa/klient-poprzedni/src/moduly/multitasking/powierzchnia-sekcji.ts

Sekcje panelu — Zespoły, Kolejki, Orkiestracja, Harmonogram i Monitor — mają ten sam
szkielet. Kolejność i widoczność podsekcji prowadzi wzorzec powłoki
`powloka/sekcje-panelu.ts` przez komendy `panel.sections.*`; ten plik wzorzec woła
i sam niczego nie przestawia.

Komendy `panel.sections.*` adresują parę okna i panelu, a panel orkiestracji należy do
środowiska, nie do okna. Adresem zostaje więc okno, przy którym sekcja pracuje: okno
koordynatora, a gdy obsada go nie ma — pierwsze okno sesji; bez żadnego okna układ
zostaje miejscowy. Brak zapisanego układu znaczy układ domyślny, więc odmowa
`panel.sections.get` zostawia sekcję kompletną.

Identyfikatorem panelu jest klucz sekcji, ten sam, którym boczna nawigacja wskazuje
sekcję. Napis składany z tytułu rozjechałby się przy pierwszej zmianie tytułu, a układ
zapisany w rdzeniu przestałby mieć odpowiednik na ekranie.

## budowa/klient-poprzedni/src/moduly/design/indeks.ts

Moduł nie osadza się sam w dokumencie i nie zna powłoki: oddaje element,
a warstwa składająca rozstrzyga, gdzie go postawić.

Rozłączenie zostaje wewnątrz modułu, ponieważ typ `WidokModulu` powłoki nie ma
czynności odpięcia. Subskrypcje kanału odpina `rozlacz()` widoku modułu.

## budowa/klient-poprzedni/src/modele/indeks.ts

Okno modeli jest jedno na klienta i trwa między otwarciami z tego samego powodu,
dla którego router nie porzuca widoku opuszczonej trasy: powtórne otwarcie wraca
do zakładki i kategorii, na której poprzednie się zakończyło. Trwałość okna
utrzymuje subskrypcję zdarzeń `account.changed` oraz `identity.changed`, dzięki
czemu zmiana dokonana na innym urządzeniu konta dociera również wtedy, gdy okno
pozostaje zamknięte.

Wywołanie z innym kanałem, właściwe dla ponownego połączenia z rdzeniem, buduje
okno od nowa, aby komendy nie szły przez transport, którego już nie ma. Sekcję
można też osadzić bez ramy okna: `utworzSekcjeModeli` oddaje ten sam byt bez
elementu `dialog` wokół niego, czyli jedna implementacja w dwóch oprawach.

## budowa/klient-poprzedni/src/moduly/automations/stany-okna.ts

Nośnik stanu treści okien modułu Automations różni się od nośnika wspólnego
zachowaniem metod `blad()` i `pusto()`: obie zdejmują wcześniejsze potwierdzenie,
zanim postawią własny komunikat. Nośnik wspólny (`komponenty/stan-tresci.ts`)
przestawia wyłącznie pas komunikatu, więc bez tego zdjęcia odmowa rdzenia stoi
obok potwierdzenia czynności poprzedniej, a zielone zdanie nad odmową potwierdza
czynność, która się nie odbyła.

Różnica siedzi w module, a nie w nośniku wspólnym, ponieważ ten sam nośnik
obsługuje okna pozostałych modułów. Przedrostek klas `da` należy do arkusza
modułu, dlatego wiązanie nośnika z przedrostkiem stoi po stronie modułu.

## budowa/klient-poprzedni/src/moduly/assistant/zrodlo-pamieci.ts

Wzorem jest `zrodlo-assistant.ts`: żadne wywołanie nie rzuca wyjątkiem, odmowa
wraca polem `blad` wyniku, a okno pokazuje ją w swoim stanie błędu.

Komenda `memory.detach` nie jest tu wystawiona. Kontrakt oznacza znaczenie
odpięcia jako nierozstrzygnięte, więc okno nie nadaje mu własnego sensu — tak
samo postępuje moduł Workspace w pliku `moduly/workspace/pamiec-pozycja.ts`.

Rodziny są dwie, ponieważ mówią o dwóch różnych bytach. Rodzina `memory.*`
prowadzi ustalenia, czyli zdania, które Operator kazał zapamiętać. Rodzina
`knowledge.*` prowadzi wskaźnik znaczenia zbudowany z treści już istniejących:
biblioteki, historii rozmów i plików przestrzeni roboczej. Zlanie ich w jedno
źródło zatarłoby, co jest ustaleniem, a co odnalezionym fragmentem.

## budowa/klient-poprzedni/src/asystent-plywajacy/sprawca-zdarzenia.ts

Napis nie bywa pewniejszy niż dowód, na którym stoi, dlatego brak pola `actor`
w kopercie zdarzenia daje rozstrzygnięcie „nieznany", a nie domysł.

Plik jest osobny od `aplikacja/rozstrzyganie-sprawcy.ts`, który rozwiązuje inne
zadanie: odpowiada na pytanie, czy czynność wykonało bieżące połączenie, i robi
to dla zdarzeń, które pola `actor` w ogóle nie niosą.

## budowa/klient-poprzedni/src/mission-control/pas-decyzji.ts

Pas jest jedynym miejscem pulpitu, w którym praca stoi, dlatego jako jedyny nosi
akcent: wstęgę trzech pikseli u góry, tło żetonu `--dn-akcent-tlo` oraz wezwanie
z cieniem akcentu. Akcent obejmuje wstęgę i przycisk, nie całą powierzchnię.

Kontrakt nie niesie odczytu przepływów wstrzymanych do decyzji: komplet danych
podaje wtedy wartość pustą, pas mówi o braku źródła i chowa wezwanie, ponieważ
przycisk wzywający do rozstrzygnięcia nieodczytanego wykazu byłby atrapą. Gdy
źródło jest, a nic nie czeka, pas zmienia treść na zdanie o braku oczekujących
przepływów, a wezwanie pozostaje czynne.

## budowa/klient-poprzedni/src/aod/sterowanie-obecnoscia.ts

Menu kebab nie jest tu budowane po swojemu: to ten sam komponent
`wyciszenie-menu.ts`, który stoi przy awatarze. Dwie kopie jednego menu czytają
jeden stan, więc wyciszenie założone przy awatarze widać w nagłówku bez żadnego
zszywania.

Stan nie jest tu przechowywany: obie kontrolki czytają i przestawiają
`StanObecnosci` z pliku `tryb-obecnosci.ts`, a rysują się z powrotem jego
powiadomieniem. Dzięki temu skrót klawiszowy i menu nigdy się nie rozjeżdżają.

## budowa/klient-poprzedni/src/moduly/apps/filtr-katalogu.ts

Filtrowanie idzie po stronie klienta, a nie w żądaniu, ponieważ komenda `extension.list`
zawęża wyłącznie rodzajem i stanem zainstalowania. Pozostałych faset — źródła, stanu
włączenia i tekstu — kontrakt w żądaniu nie ma, a cztery odczyty po jednym na fasetę dałyby
cztery migawki z czterech różnych chwil.

Zakres szukania jest ograniczony i okno musi to powiedzieć: przeszukiwane są nazwa, kod
i opis, bo tyle niesie pozycja katalogu. Szukanie po udostępnianych narzędziach
i znacznikach stoi w wykazie braków kontraktu, zamiast udawać, że pusty wynik znaczy brak
takiego rozszerzenia.

Porównanie tekstu idzie bez rozróżnienia wielkości liter, ale bez normalizacji znaków
diakrytycznych. Napisy różniące się wyłącznie znakami diakrytycznymi zostają dwoma różnymi
napisami, bo zrównanie ich w oknie kazałoby mu twierdzić coś o dopasowaniu, czego rdzeń
przy własnym szukaniu nie potwierdzi.

Zdanie podsumowujące wykaz jest konieczne, ponieważ pusty wykaz przy czynnym filtrze czyta
się jak pusty rejestr. Zdanie wymienia zakres szukania, bo tekst nieznaleziony w opisie bywa
nazwą narzędzia, której pozycja katalogu nie niesie.

## budowa/klient-poprzedni/src/mission-control/nadajnik.ts

Pulpit operacyjny nie sięga po magistralę połączenia, ponieważ jego zdarzenia
są zdarzeniami widoku, a nie kontraktu; to powłoka rozstrzyga, którą komendę
kontraktu z takiego zdarzenia zbuduje. Rozesłanie ładunku biegnie po kopii
zbioru odbiorców i w bloku ochronnym, więc odbiorca zgłaszający usterkę nie
odbiera zdarzenia pozostałym, a odpięcie w trakcie rozgłoszenia nie narusza
przebiegu pętli.

## budowa/klient-poprzedni/src/komponenty/okno-aplikacji.ts

Znakowanie służy oknom, które mają wiersz w katalogu okien operacyjnych, lecz nie
noszą ramy operacyjnej. Okno Konfiguracji i Okno Ustawień są modalami opartymi na
elemencie `dialog`, a ich wygląd niesie klasa `dn-modal`. Strona główna jest widokiem
pełnoekranowym. Okno Rozmowy dostaje ramę od gniazda układu równoległego
`okna-rownolegle/gniazdo-okna.ts`. W każdym z tych przypadków nagłówek, plakietka
roli i pas akcji z ramy operacyjnej byłyby elementem zbędnym albo powtórzonym.

Przypisanie kodu wprost do zbioru danych elementu, rozsiane po tych plikach,
dawałoby kilka miejsc, które mogą się rozejść. Funkcja nadaje przy tej samej
okazji etykietę dostępności z tej samej nazwy, więc kod okna i treść czytana przez
czytnik ekranu pozostają zgodne.

Funkcja zwraca ten sam element i ten sam typ, ponieważ modal potrzebuje typu
`HTMLDialogElement` do wywołania metod `showModal` oraz `close`. Zwracanie typu
`HTMLElement` odbierałoby te metody i zmuszało wywołujących do rzutowania.

## budowa/klient-poprzedni/src/aktualizacja/indeks.ts

Adres kanału wydań, porównanie wersji i most do powłoki są własnością tego
katalogu i na zewnątrz nie wychodzą; punkt wejścia aplikacji zna wyłącznie to,
co jest tu wymienione.

Widok wykazu wydań wychodzi stąd osobno od banera. Baner przywołuje go
przyciskiem, ale widok jest gotowy do osadzenia także w stałym miejscu, więc
jedno użycie nie przesądza o drugim.

## budowa/klient-poprzedni/src/moduly/assistant/indeks.ts

Kod modułu stoi wewnątrz opisu modułu, więc moduł sam mówi, którym jest modułem.
Rozjazd między nazwą zapisaną w `aplikacja/rejestr-modulow.ts` a rzeczywistością
staje się przez to niemożliwy.

Moduł nie zna powłoki: oddaje element, a warstwa składająca rozstrzyga, gdzie go
postawić i czy wpisać go do rejestru modułów.

## budowa/klient-poprzedni/src/aplikacja/arkusze-stylow.ts

Żetony motywu dają barwy, typografię, przestrzeń i ruch, a biblioteka `.dn-*`
stoi na nich. Arkusz `powloka.css` sprowadza zmienne widoku okna do żetonów
motywu, więc idzie po arkuszu tego widoku, a `aplikacja.css` osadza reguły
w obszarze roboczym powłoki i dlatego jest ostatni.

Arkusz kompletu sterowania `sterowanie/sterowanie.css` wciąga sam moduł panelu.
Trafia przez to do pakietu za tym wykazem, jako ostatni styl widoku, i w wykazie
nie figuruje.

## budowa/klient-poprzedni/src/moduly/diagnostics/stany-okna.ts

Rozłączność stanu błędu i stanu pustki waży w tym module szczególnie. Port
diagnostyki jest w rdzeniu odbiorcą odmów wykonania komend, więc okno błędów
pokazujące po nieudanym odczycie pusty wykaz ukrywałoby odmowy dwa razy: cudze
i swoją własną.

Pustka bywa tu zaś stanem poprawnym, ponieważ instalacja bez ani jednego błędu
w zadanym zakresie czasu nie jest usterką.

## budowa/klient-poprzedni/src/konfiguracja/dymek-objasnienia.ts

Sam dymek objaśnienia mieszka w bibliotece komponentów, a stąd jest wyłącznie
re-eksportowany, żeby okno konfiguracji miało jedno wejście.

Zdania objaśnienia są dopisywane, a nie zastępowane. Pozycja bez opisu nadal ma
czym objaśnić swój klucz, a pozycja z opisem zyskuje to, czego opis nie mówi:
jednostkę, poziom domyślny i wymóg ponownego uruchomienia.

## budowa/klient-poprzedni/src/dostepy/dialog-katalogu.ts

Most do powłoki natywnej mieszka w katalogu `powloka/`, ponieważ dotyczy powłoki
natywnej, a nie dostępów. Tutaj zostaje wyłącznie zapis wyniku przyjęty
w widokach dostępów: napis pusty znaczy, że katalogu nie wskazano.

## budowa/klient-poprzedni/src/dostepy/obszar-katalogu-roboczego.ts

Katalog roboczy modelu jest ustawieniem niezależnym od dostępów, ponieważ obie rzeczy odpowiadają na inne pytanie. Model może mieć wgląd w cudzy katalog i niczego w nim nie zapisywać, a własne katalogi sesyjne trzymać w miejscu instalacji aplikacji. Zlanie tych dwóch ustawień w jedno kazałoby otworzyć zapis wszędzie tam, gdzie model ma wyłącznie czytać.

Obszar nie jest kopią okna konfiguracji. Pokazuje dwa klucze, ponieważ bez nich sekcja dostępów byłaby niepełna: po nadaniu modelowi dostępu do katalogów następne pytanie dotyczy miejsca, w którym model będzie pisał. Pełny wybór poziomów zasięgu i osi zostaje w oknie konfiguracji, które prowadzi jedyny rachunek dziedziczenia wartości.

Wczytanie danych nie czeka na rdzeń i nie blokuje osadzenia obszaru. Pola stoją od razu z wartościami domyślnymi, a odpowiedź rdzenia nanosi na nie wartości obowiązujące.

## budowa/klient-poprzedni/src/moduly/design/tor-komendy.ts

Słownik zdań o torze komend stoi osobno od słownika etykiet modułu Design,
ponieważ tamten plik jest słownikiem nazw i braków, a ten niesie wiedzę o drodze
żądania: o kolejności kroków rdzenia, o tym, co zastępuje stan poprzedni,
i o tym, czego kontrakt nie przewiduje wcale.

Zdanie o torze nie orzeka o powodzie odmowy. Powód niesie opis odmowy wyjęty
z odpowiedzi rdzenia, a słownik dokłada wyłącznie to, czego rdzeń o sobie nie
mówi.

## budowa/klient-poprzedni/src/moduly/research/badanie-ksiazka-kodow.ts

Komenda zapisu książki kodów zapisuje ją w całości: pole kodów niesie stan po
zmianie, więc żądanie złożone z samych kodów wpisanych teraz wymazałoby wszystkie
dotychczasowe. Dlatego zapis idzie po odczycie i jest złączeniem, a nie
podmianą. Zdanie odpowiedzi mówi osobno, ile kodów zostało zachowanych i ile
dołożonych, ponieważ to jedyny sposób, żeby Operator rozpoznał wymazanie, gdyby
rdzeń zapisał co innego.

Kod istniejący rozpoznaje się po nazwie, bez względu na wielkość liter, ponieważ
nazwa jest tym, co Operator wpisuje; identyfikator kodu nadaje rdzeń i Operator
go nie zna. Kod dopisany jedzie z pustym identyfikatorem, tak samo jak nowe
pytanie badawcze w komendzie zapisu pytania warsztatu.

Definicja stoi w wierszu za znakiem rozdzielającym. Znak jest wyborem okna
i okno mówi o nim wprost przy chwycie pola. Nazwa kodu bez definicji jest
poprawna, ponieważ kontrakt ma pole opisu nieobowiązkowe.

Zdjęcie kodu z książki jest osobną czynnością i okno mówi Operatorowi, że dziś
jej nie ma. Kod zastany, którego Operator teraz nie wpisał, zostaje: wymazanie
reszty byłoby skutkiem, którego nikt nie zamówił.

## budowa/klient-poprzedni/src/mission-control/sekcja-utworz.ts

Sekcja tworzenia stoi na pulpicie, ponieważ bez niej jedyną drogą do pracy jest kliknięcie sesji już biegnącej w matrycy. Rząd kafli otwiera wejście tam, gdzie nic jeszcze nie biegnie, dlatego każdy kafel jest czynny od razu, żaden nie czeka na spełnienie warunku, a podpis kafla nazywa byt, który powstanie po naciśnięciu.

Etykiety kafli idą krojem bazowym półgrubym, a nie krojem szeryfowym, ponieważ różnica krojów niesie w tym miejscu znaczenie: krój szeryfowy oznacza wejście do środowiska, a krój bazowy oznacza zbudowanie komponentu.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-okna-modulu.ts

Każde żądanie obszaru niesie pole `windowId`, które kontrakt nazywa oknem
modułu Apps, więc identyfikator musi pochodzić z rdzenia, a nie z wyliczenia po
stronie widoku. Podaje go komenda `window.list` zawężona do sesji i do okien
modułu Apps, przy czym pole `Window.moduleId` niesie to, co rdzeń przypisał
oknu przy wejściu w obszar roboczy komendą `workspace.enter`.

Brak okna jest stanem poprawnym, ponieważ sesja bywa jeszcze nieznana w chwili
wejścia w moduł. Wykaz wraca wtedy pusty, a okna modułu pokazują stan pusty
nazywający brakujący warunek, zamiast wysyłać żądanie bez identyfikatora.

Wybór okna bierze pierwsze okno o module zgodnym z kodem modułu. Gdy rdzeń nie
przypisał żadnego, wynikiem jest pusty łańcuch, a nie okno przypadkowe.

## budowa/klient-poprzedni/src/konfiguracja/widocznosc-pol.ts

Pozycja katalogu może zależeć od innej: pole „host wykonania" ma sens dopiero
przy zasięgu zdalnym, a pole „ścieżka klucza" dopiero przy moście SSH. Warunek
żyje w katalogu, nie w kodzie interfejsu, więc zmiana zależności jest nowym
wierszem katalogu, a nie nową gałęzią w kliencie.

Warunek wskazujący klucz, którego katalog nie zna, nie chowa pola. Ukrycie pola
z powodu braku metadanej byłoby cichym odebraniem Operatorowi ustawienia,
dlatego brak klucza warunkującego zostaje odnotowany w dzienniku klienta, a pole
pozostaje widoczne.

## budowa/klient-poprzedni/src/mission-control/zdarzenia-pulpitu.ts

Pulpit nie wysyła komend samodzielnie: nadaje zamiar, a kopertę składa z niego
powłoka. Dzięki temu nazwa komendy pojawia się w jednym miejscu i pochodzi
z importu ze wspólnego kontraktu, zamiast być przepisana z pamięci do każdego
wywołania.

Pole roli niesie nazwę kolejki widoczną dla Operatora. Kontrakt roli kolejki
nie niesie, więc zamiar jej nie podstawia i nie udaje, że pochodzi z rdzenia.

Podniesienie priorytetu nie ma odpowiednika ani w wyliczeniu `QueueAction`, ani
wśród komend kontraktu. Zamiar niesie wtedy rodzaj priorytetowy z komendą pustą,
zamiast napisu wymyślonego po stronie klienta, którego rdzeń nie rozpozna.

## budowa/klient-poprzedni/src/mission-control/kolumna-zespolu.ts

Liczba podagentów nie ma źródła w kontrakcie. Wiersz mówi to wprost, zamiast pokazać
wartość zastępczą, ponieważ liczba zmyślona wygląda tak samo jak odczytana i wprowadza
w błąd przy ocenie obciążenia zespołu.

## budowa/klient-poprzedni/src/aplikacja/pasek-aplikacji.ts

Pasek ma jedną odpowiedzialność: złożenie elementów paska. Stoi nad widokami,
które własnego paska nie mają — Centrum dowodzenia oraz Mission Control.
Powłoka środowiska ma pasek własny w `powloka/pasek-gorny.ts` i drugiego nie
dostaje; przełącznik tras oraz przełącznik motywu wchodzą wtedy w jej grupę
akcji.

Motyw jest przełączalny z każdego widoku, ponieważ pasek niesie ten sam
przełącznik, którego używa powłoka. Wartości motywu ani żadnej barwy pasek nie
zna — całą pracę wykonuje warstwa `motyw/`.

Poza godłem, przełącznikiem tras i motywem pasek niesie menu aplikacji przy
godle, ikonę ustawień oraz menu Operatora w grupie akcji. Są to te same wejścia,
które ma powłoka środowiska w `powloka/akcje-paska.ts`, żeby produkt zachowywał
się jednakowo na każdej trasie. Żadne z nich nie zakłada nowej czynności ani nie
woła komendy spoza wykazu: pozycje są te same, co w listwie strony głównej,
a skutek jeden i wspólny — wykaz skutków stoi w `akcje-ustawien.ts`.

Trasy stoją w grupie akcji, a nie osobno pośrodku paska: znaki tras i kontrolki
platformy tworzą jeden rząd ikon przy prawej krawędzi, rozdzielony kreską, więc
pasek ma dwie grupy zamiast trzech.

Ikona ustawień jest skrótem, a nie drugą drogą: otwiera dokładnie tę pozycję,
którą listwa strony głównej niesie pod kodem `konfiguracja`, przez ten sam wykaz
skutków. Stoi w pasku, ponieważ ustawienia platformy wykonuje się z każdego
miejsca, a listwa widoczna jest wyłącznie na stronie głównej.

## budowa/klient-poprzedni/src/moduly/browser/zrodlo-izolacji.ts

Panel modułu jest wyłącznie odczytem, dlatego źródło wystawia same komendy
odczytu. Zapis przełączników oraz przypisanie profilu należą do okna
konfiguracji punktów izolacji, ponieważ dwa miejsca zapisujące tę samą politykę
dawałyby dwa różne zdania o tym, co obowiązuje.

Rozdzielenie odpowiada podziałowi punktów izolacji: moduł przeglądania wskazuje
punkty właściwe przeglądaniu, czyli dostęp sieciowy procesu sesji, kontenery
tożsamości oraz zakres pętli, a personalizacja tych punktów odbywa się w oknie
konfiguracji.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-znaczenia.ts

Rodzina komend `knowledge.*` ma drugiego klienta w module Wiedza, z którego moduł Library świadomie
nie korzysta. Powód pierwszy jest własnościowy: moduł nie sięga do plików innego modułu, ponieważ
wiązałby swoje okna z cudzym cyklem życia. Powód drugi jest techniczny i ważniejszy: tamto źródło
woła kanał wprost, a rodzina `knowledge.*` nie ma własnego zdarzenia odmowy, więc rdzeń bez wpiętego
portu Wiedzy odpowiada kopertą `connection.unknown`, pozbawioną pola `status`. Library idzie zatem
przez straż odmów, która wiąże taką odpowiedź z żądaniem po identyfikatorze i zamienia ją w zwykły
wynik z błędem.

Rodzina leży poza obszarem `library`, ale czyta dokładnie ten zbiór. Czytelnikiem zakresu `library`
jest repozytorium modułu Library, a pole `sourceId` trafienia niesie identyfikator pliku biblioteki,
co ustala przejściówka rdzenia `adapter_modul_wiedza_zrodla.go`. Dzięki temu trafienie daje się
odwzorować na wiersz wykazu bez drugiego odczytu.

Rozdział zdolności między rodzinami jest wiążący dla okna. Komenda `library.file.search` dopasowuje
słowa poprzez indeks pełnotekstowy repozytorium, a `knowledge.search` dopasowuje znaczenie poprzez
wskaźnik osadzeń. Tryb hybrydowy okna nie jest komendą kontraktu, lecz złożeniem obu odpowiedzi po
stronie klienta.

Wskaźnik osadzeń nie odświeża się przy wgraniu pliku. Rdzeń buduje go wyłącznie na żądanie
`knowledge.index`, ponieważ osadzanie potrafi trwać minutę i sięga po wagi modelu. Przeliczenie
wskaźnika jest więc osobną czynnością w zakładce higieny, a nie skutkiem ubocznym wgrania pliku.

## budowa/klient-poprzedni/src/aplikacja/przelacznik-motywu.ts

Oba motywy są równoprawne i kontrolka nie wskazuje żadnego z nich jako
domyślnego. Dopóki nie zapadł wybór własny, obowiązuje preferencja systemu,
a kontrolka powrotu stoi w stanie wciśniętym: pokazuje, że rozstrzygnięcie należy
do systemu.

Wartości motywu ani żadnej barwy ten plik nie zna — całą pracę wykonuje warstwa
`motyw/`, a tutaj zostaje wyłącznie obsługa kontrolek. Zmiana preferencji systemu
przy braku wyboru własnego dochodzi tą samą drogą zdarzenia co przełączenie
ręczne, więc kontrolka nadąża bez przeładowania okna.

## budowa/klient-poprzedni/src/konfiguracja/kontrolki-wyboru.ts

Definicja nastawy bez podanych opcji nie daje pustego ekranu bez wyjaśnienia: obie listy
niosą wtedy ostrzeżenie przy polu, a lista jednokrotna zostaje przy samej pozycji pustej.
Operator widzi wtedy, że wartości nie podał katalog nastaw, a nie że kontrolka się nie
zbudowała.

## budowa/klient-poprzedni/src/moduly/library/pasek-kontekstu.ts

Okno komunikacji wybierane jest po oznaczeniu modułu, a nie z brzegu wykazu, ponieważ komenda wykazu okien oddaje okna wszystkich modułów sesji. Sięgnięcie po pierwsze okno z wykazu wysyłałoby kontekst cudzego okna, dlatego pasek szuka okna należącego do modułu Library.

Okno komunikacji jest zarazem oknem źródłowym przeniesienia kontekstu, więc jego brak pasek opisuje wprost, zamiast podstawiać okno innego modułu. Tytuł okna pasek pokazuje razem z liczbą wiadomości, natomiast moduł okna i tryb uprawnień zostają w panelu sterowania okna, a tytuł powtórzony jest już w karcie sesji.

## budowa/klient-poprzedni/src/moduly/design/brak-drogi.ts

Przycisk czynności bez komendy w kontrakcie zachowuje się według wzoru
z `aplikacja/zamiary-pulpitu.ts`: pozostaje aktywny i podpisany, a powód braku
podaje dopiero po naciśnięciu. Kontrolka wyłączona albo pusta nie niesie tej
informacji, więc Operator nie odróżniłby czynności niedostępnej chwilowo od
czynności, której kontrakt w ogóle nie zna.

Znak `aria-description` niesie tę samą wiedzę technologiom wspomagającym, zanim
przycisk zostanie naciśnięty: zapowiedź płynie z etykiety, a powód z dymka
otwieranego naciśnięciem.

## budowa/klient-poprzedni/src/moduly/automations/kody-okien.ts

Rejestr okien operacyjnych rdzenia zna kody `workflow-builder`, `scheduler`,
`queue-manager`, `orchestrator` i `execution-monitor`, zakładane migracją
`migracja_030_rejestr_okien_operacyjnych.sql`; powiązanie tych kodów z modułem
`automations` niesie `migracja_031_okna_modulow.sql`. Kody stoją w jednym
miejscu, ponieważ rdzeń nie sprawdza identyfikatora okna podawanego
w `automation.execution.subscribe`, tylko zapamiętuje go jako klucz obserwacji.
Kod z przedrostkiem modułu nie wywołałby więc odmowy: trafiłby do
`Queue.windowIds` i do telemetrii postępu jako okno nieistniejące.

Znacznik wykazu gotowych pętli stoi poza wykazem kodów okien, bo służy wyłącznie
znacznikowi `data-okno` w układzie modułu, czyli skokom nawigacji i sprawdzianom
widoku. Okno wykazu nie zakłada obserwacji telemetrii i nie podaje `windowIds`
przy zakładaniu kolejki, więc jego znacznik nigdy nie opuszcza klienta.

## budowa/klient-poprzedni/src/modele/zapisy-kont.ts

Odczyt i zapis rejestru kont są rozdzielone celowo. Stan rejestru pilnuje tego, co widok wie o kontach, a moduł zapisów pilnuje tego, co dzieje się po zapisie: które wiersze wolno nanieść wprost z odpowiedzi rdzenia, a kiedy trzeba przeczytać wykaz od nowa. Rozdział pozwala zmienić regułę uzgodnienia bez dotykania ścieżki odczytu.

Komenda wskazująca konto domyślne przestawia dwa wiersze naraz, a odpowiedź rdzenia niesie tylko jeden z nich. Dlatego po jej powodzeniu wykaz zostaje przeczytany ponownie: bez tego konto tracące oznaczenie zostałoby w widoku jako drugie konto domyślne swojego rodzaju.

## budowa/klient-poprzedni/src/aplikacja/menu-aplikacji.ts

Menu nosi pasek Centrum dowodzenia oraz pasek Mission Control. Pozycje pochodzą z wykazu pozycji ustawień, czyli z tego samego wykazu, z którego powstaje listwa ustawień strony głównej, ponieważ drugi wykaz rozjechałby się z pierwszym przy dopisaniu pozycji. Skutek naciśnięcia również jest wspólny i rozstrzyga go wykaz skutków ustawień, więc samo menu nie wie, co stanie się po wyborze pozycji.

Rozwijanie wnosi komponent menu rozwijanego z katalogu okien równoległych wraz z obsługą klawiatury, znacznikiem zapowiadającym menu podręczne, zamykaniem po wskazaniu poza obszarem oraz warstwą przybornika. Menu nie zastępuje listwy strony głównej, lecz daje te same pozycje tam, gdzie listwy nie ma: na Mission Control oraz po zwinięciu strony w dół.

## budowa/klient-poprzedni/src/moduly/assistant/okno-command-tools-hub.ts

Okno zamyka obszary szybkich akcji, wywołań narzędzi oraz proaktywności modułu.
Każda z czterech pierwszych zakładek stoi nad rzeczywistą komendą kontraktu: zakładka
Akcje nad `action.list`, czyli katalogiem akcji zasięgu modułu; zakładka Makra nad
`automation.workflow.save`, jedynym trwałym miejscem dla sekwencji kroków; zakładka
Narzędzia i MCP nad `tools.catalog.list` wraz z rodziną `session.tool.*`; zakładka
Rutyny nad `automation.workflow.list`, `schedule.get` oraz `automation.schedule.set`.

Zakładki „Umiejętności" osobno nie ma z rozstrzygnięcia, nie z przeoczenia: kontrakt
trzyma narzędzia i umiejętności w jednym katalogu i rozróżnia je polem `kind` oraz
przedrostkiem źródła. Druga zakładka nad tą samą komendą udawałaby drugie źródło,
a rozróżnienie robi filtr rodzaju w zakładce narzędzi.

Piąta zakładka stoi nad rodzinami `clipboard.*` (trwała historia schowka), `snippet.*`
(słownik skrótów rozwijanych we wszystkich polach platformy) oraz `launcher.hotkey.*`
(skrót globalny wywoływacza). Schowka maszyny Operatora rdzeń nie czyta i zakładka tego
nie udaje: podział ról jest widoczny na ekranie, w pliku `panel-schowka.ts`.

Plik odpowiada wyłącznie za skład okna. Wywołania mieszkają w `zrodlo-narzedzi.ts`,
a każda zakładka ma własny plik obszaru.

Siatka akcji jest tym samym widokiem katalogu, który stoi w Voice Console: drugim
widokiem, a nie drugim źródłem. Obie sięgają `action.list` zasięgu modułu, a naciśnięcie
kafla wypełnia to samo pole polecenia.

Trzy odczyty składania okna idą równolegle, ponieważ katalog akcji, katalog narzędzi
i wykaz rutyn dotyczą trzech różnych rodzin komend i żaden nie warunkuje pozostałych.
Każdy nazywa swoje niepowodzenie w swoim obszarze.

## budowa/klient-poprzedni/src/aplikacja/akcje-ustawien.ts

Pozycja listwy pozbawiona własnego widoku pozostaje klikalna, a jej naciśnięcie
mówi wprost, że ekran jeszcze nie powstał, zamiast otwierać atrapę. Powierzchnia
Always On Display nie ma odrębnego okna i wchodzi jako rozszerzenie boczne,
dlatego jej pozycja prowadzi do powierzchni interakcji, a nie do okna.

## budowa/klient-poprzedni/src/aod/wyciszenie-aod.ts

Wyciszenie Always On Display ma pięć rodzajów: ten plik trzyma trzy z nich jako
stan (czasowe, kontekstowe, klasy zdarzeń); czwarty — tryb cichy — jest trybem
obecności i mieszka w `tryb-obecnosci.ts`; piąty — wyjątek wagi krytycznej —
nie jest wyciszeniem, tylko regułą przebijającą wszystkie pozostałe, i stoi
w regule ujawniania. Plik nie dotyka kanału ani rdzenia: kontrakt nie niesie
wyciszenia nakładki ani jednym polem (`grep -i wycisz shared/contract.go`
znajduje wyłącznie wyciszenie uczestnika tury w module Roundtable i wyciszone
wyzwolenia reguł alarmowych, nic z rodziny `aod.*`), więc stan wyciszenia
jest dziś stanem okna. Magazyn jest podawany wołaczowi (`MagazynWyciszen`),
nie brany z globalnej przestrzeni na sztywno — wzorem
`moduly/studio/widok-nastawy-operatora.ts` — żeby przełożenie zapisu na
rdzeń nie ruszyło ani jednego wołacza. Braki kontraktu nazywa wprost
`wyciszenie-braki-kontraktu.ts`.

Klasa zdarzeń każdego powodu rozpoznanego przez nakładkę przypisuje „pętlę
wstrzymaną i przerwaną" klasie stanu pętli wykonawczej, a „kolejkę
zatrzymaną, zadanie w stanie błędu i zadanie oczekujące dłużej niż próg" —
klasie stanu kolejki zadań. Pozostałe cztery klasy zdarzeń nie mają dziś
w kontrakcie nośnika sygnału (nie ma zdarzenia kontroli jakości,
harmonogramu ani powtarzalności czynności Operatora), więc ich wyciszenie
zapisze się i zadziała z chwilą, w której sygnał wejdzie; menu mówi to
wprost, zamiast udawać, że wycisza coś, co i tak milczy.

Sugestia opisana wobec wyciszenia niesie moduł i sesję jako pola opcjonalne,
bo nie każda sugestia je zna: telemetria rdzenia niesie okno i sesję, modułu
nie niesie wcale, a `AodSuggestion` kontraktu niesie samo okno. Funkcja
rozstrzygająca, które wyciszenie obejmuje sugestię, nie zna wyjątku wagi
krytycznej — wyjątek nie znosi wyciszenia, tylko przepuszcza sugestię
plakietką mimo niego, a rozstrzyga to reguła ujawniania w
`tryb-obecnosci.ts`, w jednym miejscu dla wszystkich trzech rodzajów.

## budowa/klient-poprzedni/src/moduly/assistant/braki-kontraktu.ts

Okno nie stawia w miejscach bez drogi kontraktu przycisku wygaszonego ani
milczącego. Kontrolka zostaje klikalna i odpowiada zdaniem mówiącym, czego
w kontrakcie nie ma; tak samo postępuje `aplikacja/zamiary-pulpitu.ts`.

Każde zdanie wymienia komendę albo pole, którego brakuje, żeby zgłoszenie braku
dało się napisać wprost z ekranu, bez sięgania do innego opracowania.

## budowa/klient-poprzedni/src/moduly/multitasking/strumien-wykonawcy.ts

Bieg naprawczy prowadzi rdzeń w typie `session.StrumienWykonawcy`. Po stronie
klienta leży wyłącznie odkładanie fragmentów; klient niczego nie wybudza ani nie
rozpoczyna.

Odkładane są wszystkie rodzaje fragmentów, nie sama odpowiedź: tok rozumowania,
wywołania narzędzi, ich wyniki oraz pliki. Bez wywołań narzędzi koordynator nie
odróżnia wykonawcy pracującego od zatrzymanego.

Tekst narasta w obrębie jednej tury. Nowa wiadomość zaczyna wynik od nowa,
inaczej przekazanie wyniku niosłoby sklejkę dwóch odpowiedzi wykonawcy.

## budowa/klient-poprzedni/src/moduly/browser/wypis-zrodel.ts

Wykaz źródeł przychodzi do okna komendą `browser.source.list`, a notacja jest
wyłącznie sposobem jego zapisania, dlatego komendy eksportu bibliografii kontrakt
nie niesie i nie musi.

Rodzajem wpisu jest zasób elektroniczny — moduł zbiera strony internetowe i tylko
o nich może tak zaświadczyć. Datą jest czas dodania źródła do wykazu okna,
ponieważ daty publikacji strony rdzeń nie oddaje.

## budowa/klient-poprzedni/src/aod/sekcja-glosu.ts

Pola `audioRef` i `transcript` kontraktu są opcjonalne, a nakładka nagrań nie
tworzy, więc polecenie wydaje się samą transkrypcją.

Sekcja rozmowy jest osobną drogą. Komenda `aod.chat.send` niesie zdanie do okna
rozmowy, natomiast `aod.voice.command` kieruje polecenie do asystenta i ma własne
pole `speak` na odpowiedź syntezą mowy.

## budowa/klient-poprzedni/src/moduly/library/zrodlo-wersji.ts

Podział wobec pliku `zrodlo-biblioteki.ts` biegnie wzdłuż odpowiedzialności:
historia dokumentu należy tutaj, natomiast plik, wyszukiwanie i kolekcje należą
tam. Kształt pozostaje jeden, `ZrodloBiblioteki`, aby okna widziały jedno źródło,
a nie dwa.

## budowa/klient-poprzedni/src/aplikacja/komunikaty.ts

Żaden przycisk w aplikacji nie jest wyszarzany, więc każde naciśnięcie musi dać
Operatorowi odpowiedź. Gdy zamiar nie ma odpowiednika w komendzie kontraktu albo
rdzeń odmawia wykonania, dymek mówi to wprost, zamiast udawać czynność wykonaną.
Milczenie po naciśnięciu byłoby dla Operatora nieodróżnialne od powodzenia.

## budowa/klient-poprzedni/src/moduly/library/ster-modulu.ts

Wybór modułu docelowego jest wyborem z katalogu, a nie polem tekstowym. Komenda
przeniesienia z nieznanym kodem modułu wraca powodzeniem, a rdzeń zakłada okno
z dokładnie tym kodem, więc literówka kończy się oknem, do którego nikt nie
wejdzie. Katalog pozycji pochodzi z komendy wykazu modułów.

Wybór zmienia nastawę we wspólnym stanie modułu. Nastawa jest jedna dla całego
modułu Library, więc zmiana dokonana w oknie eksploratora jest widoczna w oknie
podglądu pliku i odwrotnie.

Pozycja pusta zostaje w wykazie: bierze ona moduł wytwórcy przeniesionego pliku,
a ster tę drogę nazywa wprost, zamiast kazać jej się domyślać z pustego pola.

Gdy wykaz modułów wróci odmową, ster ma samą pozycję pustą, a pod nim stoi powód
odmowy. Milcząca lista z jedną pozycją wyglądałaby jak platforma z jednym
modułem, czyli mówiłaby Operatorowi nieprawdę o stanie rdzenia.

Moduł Library wypada z własnego wykazu pozycji, ponieważ przeniesienie kompletu
do modułu, w którym Operator już stoi, założyłoby drugie okno tego samego modułu
i nie otworzyłoby niczego nowego. Tak samo zawęża swój katalog moduł Design.

Przy odświeżeniu wartość bierze się ze stanu modułu, a nie z kontrolki: funkcja
ustawiająca pozycje utrzymuje poprzedni wybór, o ile pozycja nadal istnieje, ale
prawdą o nastawie jest stan, także wtedy, gdy zmieniło ją drugie okno. Nastawa
wskazująca moduł, którego katalog już nie zna, zniknęłaby po cichu na pozycję
pustą, czyli przeniesienie poszłoby gdzie indziej, niż mówił ster przed chwilą.

## budowa/klient-poprzedni/src/ikony/zrodla-ikon.ts

Wykaz nazw jest w tym pliku jedynym źródłem prawdy — typ `NazwaIkony` wywodzi się
z rejestru, a nie z osobnej listy. Same wiązania nazw z plikami mieszkają w katalogu
`zrodla/`, w podziale na grupy zastosowań zgodnym z porządkiem `ikony/manifest.json`;
ten plik wyłącznie je scala.

Godło marki stoi poza zestawem 82 ikon, bo jest jedynym znakiem o barwach własnych.
Pozostałe odmiany znaku niesie `ikony/marka.ts`.

Nazwy wycofane pakietem design v2.0 wskazują na następców z zestawu — nie wnoszą
żadnego nowego rysunku ani pliku spoza pakietu. Znikają wraz z ostatnim wywołaniem.

## budowa/klient-poprzedni/src/moduly/developer/braki-kontraktu.ts

Wykaz `KOMENDY` z `shared/contract.ts` jest wykazem komend kontraktu dostępnym
w czasie działania, wytwarzanym z `contract.json`. Zdanie powodu powstaje
z odczytu tego wykazu przy składaniu okna, więc dopisanie komendy do kontraktu
przepisuje je samo. Zdanie wpisane na stałe przestałoby być prawdziwe w dniu
takiej zmiany i nikt nie miałby po czym tego rozpoznać.

Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje. Na to pytanie
odpowiada `katalog-komend.ts`, który pyta rdzeń o jego własny rejestr. Tutaj
brak jest brakiem po stronie kontraktu i tylko o kontrakcie zdanie mówi.

Bliźniaczy mechanizm stoi w module Diagnostics. Biblioteka `komponenty/` nie ma
dziś wspólnego bytu dla obu, a modułowi nie wolno sięgać do wnętrza sąsiada,
więc każdy z nich składa zdanie powodu u siebie.

## budowa/klient-poprzedni/src/aplikacja/ulotnosc-okna.ts

Informację o pamięci sesyjnej modułu bierze `politykaModulu` wyłącznie z profilu
modułu. Powłoka dokłada do niej jedną rzecz, której warstwa rozmowy znać nie
może: co w danym module jest kontekstem roboczym, którego zmiana kończy rozmowę.

Taki kontekst ma jeden moduł — Agents, gdzie kontekstem jest testowany ekspert.
Powłoka jest jedynym miejscem, w którym obie strony są widoczne naraz: widok
modułu powstaje w `przestrzen-modulu`, rozmowa okna w `wiazanie-gniazda`, a żadna
z tych warstw nie zna drugiej.

Moduł nieobjęty tą funkcją również przechodzi przez `politykaModulu` i dostaje
politykę wprost z własnego profilu, bez kontekstu roboczego. Odjęcie pamięci
sesyjnej kolejnemu modułowi czyni jego okno ulotnym bez zmiany w tym pliku.

## budowa/klient-poprzedni/src/aod/indeks.ts

Reszta aplikacji zna stąd dwie czynności. `zaczepAod(kanal, opis?)` buduje
warstwę do osadzenia w powłoce i to ona wnosi pływający awatar widoczny bez
interakcji. `otworzPowierzchnieAod(kanal, opis?)` otwiera powierzchnię interakcji
jako rozszerzenie boczne i jest wołana z listwy ustawień.

Warstwa jest jedna na klienta i żyje między otwarciami powierzchni: kolejka
decyzji dosypuje się ze zdarzeń `progress.changed` oraz `window.state.changed`
także wtedy, gdy kolumna stoi zwinięta, więc plakietka awatara jest prawdziwa
jeszcze przed otwarciem powierzchni.

Kanał podaje się przy pierwszym zaczepieniu. Wywołanie z innym kanałem, po
ponownym połączeniu z rdzeniem, buduje warstwę na nowo, żeby przyszłe komendy
`aod.*` nie szły przez transport, którego już nie ma.

Powierzchnia nie czeka na rdzeń: odczyty `aod.status.get`, `aod.context.get`
i `aod.suggestion` dojeżdżają do otwartej kolumny odpowiedzią, a czynności
`aod.observe.attach`, `aod.observe.detach`, `aod.chat.send` oraz
`aod.voice.command` jadą z sekcji na żądanie. Warstwa niezaczepiona w powłoce
osadza się w korzeniu dokumentu, ponieważ pozycja „Always On Display" listwy
ustawień ma otwierać powierzchnię, a nie odmawiać z powodu kolejności montażu.

## budowa/klient-poprzedni/src/moduly/diagnostics/braki-kontraktu.ts

Plik `shared/contract.ts` niesie wykaz `KOMENDY`, czyli komendy kontraktu dostępne
w czasie działania, wytwarzane z pliku `contract.json`. Zdanie o powodzie powstaje
z odczytu tego wykazu przy składaniu okna. Zdanie wpisane na stałe przestałoby być
prawdziwe w dniu dopisania komendy do kontraktu i nikt by tego nie zauważył.

Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje; jest to pytanie osobne.
Brak jest tu brakiem po stronie kontraktu i wyłącznie o kontrakcie zdanie mówi.

Treść czynności podaje okno wywołujące, ponieważ to okno zna przeznaczenie danej
pozycji inwentarza kontrolek.

## budowa/klient-poprzedni/src/aod/pola-wykazu.ts

Cztery sekcje okna — stan, podpowiedzi, obecność i kontekst — wypisują pola tym
samym wykazem `dl`. Wspólne miejsce trzyma pętlę `dt`/`dd` w jednym egzemplarzu,
dzięki czemu zmiana układu pola obowiązuje we wszystkich sekcjach naraz.

Plik nie zna barw ani odstępów. Nadaje wyłącznie klasy z przedrostkiem `ao-`,
a ich wygląd pokrywa `aod.css` żetonami `--dn-*`.

## budowa/klient-poprzedni/src/moduly/design/polecenie-z-odmowy.ts

Komenda `design.asset.generate` kończy się błędem, gdy w rejestrze brakuje kanału
obrazowego, kanał jest nieczynny albo wskazany kanał okazuje się tekstowy, gdy
brakuje poświadczenia, gdy odpowiedź nie niesie obrazu oraz gdy bajtów nie da się
pobrać. Każda taka odmowa niesie gotową treść polecenia w polu
`details.polecenie`, ponieważ rdzeń składa prompt w całości, zanim cokolwiek
wyśle. Okno pokazuje ten tekst, aby praca nad promptem nie ginęła wraz z odmową
i dała się użyć poza produktem, a zdanie towarzyszące mówi wprost, że obrazu nie
ma i skąd tekst pochodzi.

Kształt pola `details` jest sprawdzany, a nie zakładany: kontrakt opisuje je jako
`unknown`, więc rzutowanie na własny typ byłoby obietnicą bez pokrycia. Brak pola
albo pole innego kształtu daje pusty wynik, nie wyjątek.

## budowa/klient-poprzedni/src/aplikacja/otwarcie-okna.ts

Pierwsze okno komunikacji otwiera uzgodnienie połączenia, prowadzone w pliku
`protokol/uzgodnienie`. Kolejne okna powstają dopiero wtedy, gdy Operator
wprowadza je na scenę, i przechodzą tą samą drogą kontraktu: mają własny
identyfikator, własną historię oraz własne ustawienia.

Nazwa komendy i kształt żądania pochodzą z pliku `shared/contract`, a przekład
opisu okna na treść żądania wykonuje `zamowienieOkna`. Dzięki temu plik
zamówienia nie zna ani jednej nazwy pola kontraktu i nie rozjeżdża się z nim
przy zmianie kształtu żądania.

Rola przekazana w wywołaniu nadpisuje rolę zapisaną w opisie okna, ponieważ
o roli okna na scenie rozstrzyga jego miejsce w figurze koordynator–wykonawca,
a nie konfiguracja budowania pakietu.

## budowa/klient-poprzedni/src/moduly/multitasking/format-zadan.ts

Zapis liczb i czasów stoi w osobnym pliku, ponieważ reguła zapisu czasu trwania jest potrzebna w trzech miejscach naraz: w wierszu podsumowania przepływu, w kolumnie czasu w tabeli agentów oraz w liczniku pozycji zwiniętych. Trzy kopie tej samej reguły rozjechałyby się przy pierwszej poprawce.

Podagent bez znacznika startu otrzymuje brak, a nie czas zerowy. Zero znaczyłoby, że praca ruszyła i nie trwała ani chwili, co jest twierdzeniem innym niż brak pomiaru.

## budowa/klient-poprzedni/src/moduly/browser/wiersz-zrodla.ts

Każda akcja wiersza jest naciskalna — również usunięcie, którego kontrakt nie
niesie. Przycisk odpowiada wtedy nazwaniem brakującej komendy, zamiast znikać
albo gasnąć, ponieważ kontrolka wygaszona bez wyjaśnienia nie mówi, czego
brakuje.

Wiersz nie usuwa pozycji z widoku na własną rękę. Zniknięcie pozycji bez zapisu
w rdzeniu byłoby udawaniem wykonania czynności, która się nie odbyła.

Wykaz źródeł, formularz dodania oraz przekazanie do modułu Research mieszkają
w oknie modułu, a nie w wierszu. Wiersz odpowiada wyłącznie za własną treść
i za wywołanie czynności podanych w `AkcjeZrodla`.

## budowa/klient-poprzedni/src/dostepy/usuniecie-punktu.ts

Usunięcie punktu dostępu unieważnia wszystkie nadania, które się na ten punkt
powoływały, także w oknach, których Operator w danej chwili nie widzi. Zasięg
skutku jest więc szerszy od tego, co widać na ekranie, i to on rozstrzyga
o dwustopniowym przebiegu czynności: pierwsze naciśnięcie nazywa skutek, drugie
go wywołuje.

Potwierdzenie nie korzysta z okna `confirm` przeglądarki, ponieważ okno takie
blokuje wątek dokumentu i nie da się go ubrać w warstwę wizualną platformy.
Potwierdzenie mieszka więc w samym przycisku.

Zamiar wygasa samoczynnie po upływie czasu trwania. Przycisk zostawiony w stanie
potwierdzenia byłby pułapką dla kolejnego naciśnięcia w to samo miejsce,
a wygaśnięciu towarzyszy komunikat mówiący, że punkt pozostał nietknięty.

W drzewie obecne są definicje typów zarówno przeglądarki, jak i środowiska Node,
dlatego typ zegara bierze się z wyniku samej funkcji `setTimeout`, a nie z typu
liczbowego.

## budowa/klient-poprzedni/src/moduly/apps/kanwa-komponentow.ts

Linia zależności jest nazwana, a nie narysowana. Kontrakt niesie zależność jako
wykaz identyfikatorów komponentów, bez współrzędnych i bez kierunku przepływu,
więc rysunek strzałek wymagałby danych, których w kontrakcie nie ma. Kafel podaje
zależność nazwami komponentów, a identyfikator wskazujący komponent spoza kanwy
zostaje w treści wprost, żeby niespójność układu była widoczna zamiast zniknąć.

## budowa/klient-poprzedni/src/mission-control/stan-zrodla.ts

Plik ustala kształt i wartość początkową wsadu złożenia danych pulpitu. Pole
`odczytano` mówi, czy rdzeń już się odezwał, i tym odróżnia pulpit pusty
z powodu braku odpowiedzi od pulpitu pustego z powodu braku danych.

Nazwa kolumny matrycy środowisk pochodzi z odczytu `environment.list`, a nie
z kopii katalogu utrzymywanej w pulpicie. Dzięki temu kolumna znika i pojawia
się wraz ze środowiskiem po stronie rdzenia, bez osobnej pielęgnacji kopii.

## budowa/klient-poprzedni/src/aod/kolumna-aod.ts

Funkcja Always On Display nie ma własnego okna. Całą jej powierzchnią jest
awatar oraz kolumna otwierana jako rozszerzenie boczne po prawej stronie
obszaru roboczego — na pełną wysokość, o regulowanej szerokości. Otwarcie
zwęża kolumny obszaru roboczego, nie przesłania ich i nie zamyka żadnego
z okien komunikacji operacyjnej.

Stąd trzy różnice wobec modala: brak przyciemnienia tła, brak pułapki
ogniska i brak zabranej pracy pod spodem. Kolumna leży na własnej warstwie
graficznej, wyższej niż warstwa modala.

Ten plik wyłącznie składa powierzchnię: rama z biblioteki, stan treści
z biblioteki, sekcje z osobnych plików, sterowanie obecnością z osobnego
modułu i dwa źródła komend — jedno rodziny funkcji, drugie komend cudzych
rodzin potrzebnych kolejce decyzji.

Kolumna subskrybuje dwa zdarzenia rdzenia i zdejmuje subskrypcje przy
rozłączeniu: zmianę postępu, rozgłaszaną do każdego połączenia i niosącą
proces, etap, procent, stan oraz okno; oraz zmianę stanu okna, jedyny żywy
nośnik stanu pętli, ponieważ zmiana postępu nie wypełnia tego pola. Stan
sprzed otwarcia kolumny nadrabia osobny odczyt stanu nakładki.

Subskrypcje żyją między zbudowaniem warstwy a rozłączeniem, nie między
otwarciem a zamknięciem kolumny: proces zmienia stan niezależnie od tego,
czy Operator patrzy. Dzięki temu plakietka awatara jest prawdziwa także
wtedy, gdy kolumna stoi zamknięta.

## budowa/klient-poprzedni/src/main.ts

Punkt wejścia jest wyłącznie kompozycją: logika, typy i obsługiwacze mieszkają
w modułach. Kolejno pyta powłokę natywną o adres rdzenia (`powloka/most-rdzenia`),
wczytuje warstwy wizualne (`aplikacja/arkusze-stylow`), składa aplikację
(`aplikacja/aplikacja`) i uruchamia ją: najpierw widok, potem łączność.

Adres gniazda rdzenia ustala się z lokalizacji dokumentu (`polaczenie/adres-rdzenia.ts`),
a w oknie powłoki natywnej z pakietem osadzonym wychodzi z tego adres `tauri.localhost`,
pod którym nikt nie nasłuchuje. Powłoka zna adres prawdziwy i podaje go poleceniem
`adres_rdzenia`. Pytanie musi paść przed złożeniem aplikacji, ponieważ złożenie zakłada
transport na adresie; zadane później zastałoby gniazdo już otwarte pod adresem ślepym.
Poza powłoką natywną odpowiedzią jest brak, a pytanie kosztuje jedno sprawdzenie
środowiska.

Kolejność dwóch ostatnich wywołań uruchomienia ma znaczenie. Router pokazuje trasę
zapisaną w adresie dokumentu, zanim ruszy transport, więc Operator widzi Centrum
dowodzenia od pierwszej klatki, a nie puste tło czekające na rdzeń. Uzgodnienia nie
rozpoczyna się tutaj: robi to przepływ komunikatów w chwili, gdy transport zgłosi stan
połączenia. Wywołanie stąd dałoby drugie powitanie i podwojenie historii otwarcia.

Uruchomienie bramki idzie po `polacz()`, ponieważ rozpoznanie bramki wysyła komendy
`auth.*` i ekran potrzebuje kanału, który już rusza; wcześniej pokazywałby formularz,
którego nie ma jak wysłać. Ekran montuje się sam, wywołaniem `document.body.append`
wewnątrz `uruchom()`, dlatego punkt wejścia go nie osadza — inaczej przesłona wisiałaby
w dwóch miejscach naraz.

Wywołanie zwrotne `naWejscie` wiąże świeżą sesję z żywym połączeniem: `connection.hello`
przyjmuje token, a rdzeń zapamiętuje, które gniazdo należy do której sesji bramki. To
powitanie nie dubluje powitania z uzgodnienia, ponieważ tamto pada przy nawiązaniu, zanim
Operator poda hasło, więc niesie token sesji poprzedniej albo żaden. Bez powtórzenia po
wejściu rdzeń znałby gniazdo jako niezwiązane przez całe uruchomienie, a komenda
`auth.password.reset` rozłączałaby tego, kto właśnie zmienił hasło. Powitanie jest czystym
odczytem: nie zakłada sesji pracy ani okna.

Scena wejścia staje w tej samej chwili, w której bramka schodzi. Kolejność jest wiążąca:
scena idzie do dokumentu przed zdjęciem przesłony, ponieważ bramka woła `naWejscie`, zanim
usunie swój element; między jednym ekranem a drugim nie mignie wtedy strona główna z pustym
jeszcze wykazem środowisk. Scena stoi warstwę niżej niż bramka, wedle `ladowanie.css`,
i schodzi, gdy pierwszy odczyt strony `home.enter` się domknie albo gdy minie jej własny
kres czekania, żeby rdzeń milczący nie zamienił jej w zasłonę nad produktem.

Pas aktualizacji pojawia się dopiero wtedy, gdy kanał wydań ma wersję nowszą niż
zainstalowana, a przy pierwszym uruchomieniu pyta o to z opóźnieniem, żeby nie konkurować
z uruchomieniem o łącze.

## budowa/klient-poprzedni/src/moduly/agents/archiwum-kontrolki.ts

Plik trzyma wyłącznie budulec panelu archiwum: wiersz wykazu wraz z kontrolką
pozycji oraz przekład odmowy rdzenia na zdanie czytelne dla Operatora. Składanie
panelu i jego czynności należą do `archiwum-ekspertow.ts`, dzięki czemu kontrolki
pozostają wolne od wiedzy o przebiegu pracy panelu.

## budowa/klient-poprzedni/src/mission-control/kafel-liczby.ts

Liczba w kaflu idzie krojem technicznym `--dn-ff-mono`, ponieważ jest daną
techniczną, a nie nagłówkiem. Krój szeryfowy zostaje zastrzeżony dla tytułów
sekcji, przez co czytelnik odróżnia pomiar od nazwy części widoku samym
kształtem znaków.

Wyróżnienie kafla jest zastrzeżone dla miary krytycznej. Etykieta nazywa miarę
wprost i nie bywa zastępowana samym kolorem, dzięki czemu kafel pozostaje
czytelny bez rozróżniania barw.

## budowa/klient-poprzedni/src/aplikacja/menu-operatora.ts

W menu nie ma nazwy własnej, adresu, awatara ani przełączania kont, ponieważ
platforma kont nie prowadzi: komenda `auth.register` konta nie zakłada i adresu
poczty nie przyjmuje. Menu nie udaje więc profilu osobowego, którego pod spodem
nie ma.

Nie ma również pozycji wylogowania. Rdzeń nie odcina komend po wygaśnięciu
sesji, więc taka pozycja mogłaby jedynie skasować zapis sesji i przeładować
aplikację, obiecując Operatorowi skutek, którego nie osiąga.

Obie pozostałe pozycje pochodzą z tego samego wykazu ustawień, co listwa strony
głównej, i idą tą samą drogą skutku. Dwa wykazy rozjechałyby się przy pierwszej
zmianie nazwy albo trasy pozycji.

## budowa/klient-poprzedni/src/moduly/multitasking/wcielenia-analizy.ts

Wybrane wcielenie utrwala komenda `config.set` na poziomie okna analityka. Rdzeń
nie rejestruje komendy `role.update`, więc profilu wcielenia nie zna; pokrycie
kontraktu mierzy `braki-kontraktu.ts`.

Kontrakt nie ma komendy utrwalającej werdykt oceny jako stan: `monitor.status`
czyta stan procesów i niczego nie zapisuje. Ocena analityka jest poleceniem dla
koordynatora, dlatego zgłoszenie niezgodności idzie zwykłym `message.send`.

## budowa/klient-poprzedni/src/moduly/multitasking/zrodlo-zespolow.ts

Zespół i ekspert stoją w jednym źródle, ponieważ zespół kontraktu jest wyłącznie
wykazem identyfikatorów ekspertów. Widok zespołu odcięty od wykazu ekspertów
pokazywałby Operatorowi ciąg identyfikatorów bez nazw, więc obie rodziny komend
obsługuje to samo źródło.

Zespół nie ma pola ról ani powiązania z kolejkami: komenda zapisu przyjmuje
nazwę, opis i skład, i nic poza tym. Gotowy układ ról, powiązań i kolejek daje
się zatem zapisać wyłącznie jako skład, a reszta nie ma w kontrakcie gdzie
zamieszkać. Sekcja mówi o tym Operatorowi wprost, zamiast wpychać te dane w opis
zespołu, gdzie żaden odbiorca by ich nie odczytał.

## budowa/klient-poprzedni/src/modele/stany-odczytu.ts

Pas stanów stoi pod paskiem osi, a nad zakładkami, czyli tam, gdzie widać go
niezależnie od wybranej zakładki. Zawartość zakładek zostaje na miejscu
i pozostaje czynna także przy niepowodzeniu odczytu; podmiana panelu na komunikat
zablokowałaby pracę w pozostałych obszarach sekcji.

Komunikat błędu nie znika po naciśnięciu. Przycisk ponowienia wyzwala odczyt
i zostawia zdanie na miejscu, dopóki sytuacja nie ustanie, dzięki czemu ponowna
próba nie jest brana za powodzenie.

Stan pusty rejestru niesie sam wykaz kont, ponieważ mówi o zawartości rejestru,
a nie o przebiegu odczytu.

## budowa/klient-poprzedni/src/moduly/browser/wykaz-automatyk.ts

Podział odpowiedzialności jest tu taki sam jak w panelach pomocniczych modułu przeglądarki, czyli w plikach wiersz-zrodla.ts oraz wiersz-notatki.ts. Okno składa formularz i prowadzi rozmowę z rdzeniem, a moduł wykazu zamienia wynik tej rozmowy w wiersze widoczne na ekranie. Dzięki temu postać wykazu można zmienić bez dotykania obsługi komend.

## budowa/klient-poprzedni/src/aod/warstwa-aod.ts

Warstwa Always On Display osadza się w powłoce na wysokiej warstwie graficznej, poza obszarem podmienianym przez moduł, więc przełączenie środowiska ani modułu nie zabiera awatara i nie przeładowuje kolumny. Jest jedynym miejscem, które zna naraz cztery rzeczy i wiąże je jedną regułą: kolejkę decyzji, czyli ile sugestii czeka i o jakiej wadze; stan obecności, czyli tryb pełny, cichy albo ukryty i trwające wyciszenie; progi ujawniania, czyli wagę, limit godzinowy i odstęp między dymkami; oraz układ, czyli szerokość otwartej kolumny, o którą odsuwa się awatar i dymek. Skróty klawiszowe są zaczepione w tym samym module, ponieważ tylko tutaj widać wszystkie trzy elementy naraz.

Kolumnę zwężoną poniżej progu czytelności powierzchnia interakcji zwija do dymka kontekstowego, a dymek skraca treść do jednego zdania z działaniem rozwinięcia.

Menu kebab na powierzchni interakcji awatara daje dodatkową akcję wyciszenia bezpośrednio z pozycji awatara. Menu jest tym samym komponentem, który stoi w nagłówku kolumny, a nie nowym wzorem, i pozostaje osobnym wyzwalaczem: kliknięcie pojedyncze awatara nadal otwiera dymek, podwójne nadal otwiera powierzchnię interakcji, a wyciszenie nie wymaga przejścia do kolumny.

Reguła samoczynnego ujawnienia otwiera dymek wyłącznie przy wadze wysokiej, przy nienaruszonym limicie godzinowym i odstępie, poza trybem cichym i poza wyciszeniem. Sugestia krytyczna wyciszona ujawnia się plakietką, bez dymka, a robi to funkcja odświeżająca wygląd, nie reguła ujawnienia.

Decyzja przekazywana regule wyciszenia niesie klasę zdarzenia wynikającą z powodu rozpoznania, kartę sesji z telemetrii oraz moduł z odczytu bieżącego okna w kontekście wyciszenia. Modułu, którego nakładka nie rozpoznaje, funkcja nie podstawia: sugestia bez modułu nie wpada w wyciszenie modułu.

Plakietka liczy sugestie ujawniane, a nie wszystkie: liczba obejmująca sugestie wstrzymane wyciszeniem wybiórczym obiecywałaby coś, czego wyciszenie nie pokaże. Że coś milczy, mówi stan awatara i menu wyciszania. Wyciszenie czasowe chowa plakietkę, z wyjątkiem sugestii krytycznej, która ujawnia się mimo każdego wyciszenia; wtedy plakietka liczy same sugestie krytyczne z tego samego powodu.

Zwężenie obszaru roboczego o szerokość otwartej kolumny idzie wyściółką rodzica warstwy, w której warstwa siedzi, a nie zmianą w katalogu powłoki: funkcja globalna dokłada się do układu, a nie przepisuje go. Klasa zwężenia i szerokość jadą razem, więc powłoka bez otwartej kolumny nie nosi po niej śladu.

## budowa/klient-poprzedni/src/moduly/multitasking/zwiniete-zakonczone.ts

Przepływów zakończonych bywa kilkadziesiąt i wypisane w całości spychałyby poza
krawędź panelu pracę trwającą. Zwinięcie nie jest ukryciem: licznik podaje ich
liczbę, a rozwinięcie oddaje każdy przepływ po nazwie i stanie wraz
z podsumowaniem czasu odcinka oraz liczby podagentów.

## budowa/klient-poprzedni/src/mobile/indeks.ts

Katalog wystawia reszcie aplikacji jedną czynność — `otworzOknoMobile(kanal)` —
oraz typ okna. Okno jest jedno na klienta i żyje między otwarciami, więc
powtórne otwarcie ponawia odczyt zamiast budować okno od nowa.

Kanał podaje się przy pierwszym otwarciu. Wywołanie z innym kanałem, czyli po
ponownym połączeniu z rdzeniem, rozłącza okno dotychczasowe i buduje je na nowo,
żeby przyszłe komendy `mobile.*` nie szły przez transport, którego już nie ma.

## budowa/klient-poprzedni/src/konfiguracja/adres-ustawienia.ts

Ten sam kształt adresu służy dwóm rolom, dlatego mieszka w jednym module.
Pierwszą rolą jest punkt widzenia, czyli miejsce, z którego oglądana jest
konfiguracja i względem którego liczone jest dziedziczenie. Drugą rolą jest
adres zapisu, czyli miejsce, w którym komenda `config.set` zapisze wartość.

Byt pusty znaczy poziom bez bytu: na osi poziomów jest to poziom globalny, a na
osi rozstrzygania — platforma. Dzięki temu adres nie potrzebuje osobnego
znacznika nieobecności bytu.

## budowa/klient-poprzedni/src/moduly/automations/wykaz-petli-widok.ts

Przycisk uruchomienia otrzymuje każda pozycja wykazu, również pętla wyłączona.
Stan pętli stoi w opisie pozycji, a nie w blokadzie wiersza, ponieważ decyzja
o uruchomieniu należy do okna, a nie do widoku wykazu.

Zawężanie wykazu ma dwie drogi: napis szukania oraz przełącznik obejmujący
wyłącznie pętle czynne, który zawęża żądanie do rdzenia polem `enabledOnly`.
Szukanie przebiega po stronie klienta, ponieważ komenda `automation.workflow.list`
przyjmuje jedynie `enabledOnly` oraz `limit`, a pełny wykaz jest już w oknie.
Dopasowanie bierze kolejno nazwę, identyfikator i opis pętli.

Widok oddaje osobne zdanie dla każdego rodzaju pustki: rdzeń bez ani jednej
zapisanej pętli jest innym brakiem niż wykaz zawężony napisem, który do niczego
nie pasuje.

## budowa/klient-poprzedni/src/aod/sekcja-stanu.ts

Wykaz procesów przypiętych niesie sekcja obecności, a nie sekcja stanu,
ponieważ to w sekcji obecności Operator przypina i odpina proces. Wykaz stoi
przy czynności, która go zmienia, dzięki czemu skutek czynności jest widoczny
w tym samym miejscu, w którym została ona wykonana.

Pola puste zostają puste: widok stawia oznaczenie braku zamiast wartości
zmyślonej po stronie widoku. Świeża nakładka bez okna ogniskowanego jest
stanem poprawnym, więc pusta karta sesji i puste okno nie są sygnałem błędu.

Odmowa jednego odczytu jest faktem o jednym odczycie, a nie o całej nakładce.
Pozostałe sekcje stoją na innych komendach i zostają widoczne, przez co odmowa
komendy `aod.status.get` nie wygasza całego widoku.

## budowa/klient-poprzedni/src/moduly/developer/okno-build-output.ts

Build Output i Run & Debug tworzą jedno okno o dwóch częściach zajmujących tę samą kolumnę,
przełączanych pasem zakładek w jej nagłówku — stąd jedna rama, dwie zakładki i jeden kod
okna `build-output` z katalogu rdzenia; Run & Debug nie jest osobnym oknem i osobnego
wiersza katalogu nie dostaje.

Aktualizacja na żywo idzie ze zdarzenia `developer.build.changed`, które okno subskrybuje
bezpośrednio u źródła, mimo że `stan-developer.ts` subskrybuje to samo zdarzenie i budzi
okna przez `stan.naZmiane(...)`: stan tylko rozgłasza, że coś się zmieniło, po filtrze
`windowId`, nie niesie pola `logLine` ani przebiegu z treści zdarzenia, a kontrakt nie ma
komendy, którą dałoby się dogonić stan inaczej. Ta druga subskrypcja bierze więc co innego
niż subskrypcja stanu i nie jest powieleniem tej samej pracy.

Zgłoszenie przebiegu prowadzi do pliku: zgłoszenia niosą pola `path` i `line`, a każde ze
ścieżką ma przejście „Otwórz w edytorze”, które woła `stan.wskazPlik`, a Code Editor otwiera
plik sam.

Log jest ucięty przy otwarciu okna w trakcie przebiegu i rośnie wyłącznie z pola `logLine`
kolejnych zdarzeń, więc okno otwarte po starcie przebiegu — albo przebiegu uruchomionego
przez inne okno tego konta — widzi tylko ogon logu. Ucięcie jest oznaczone wprost w treści
okna, nie zamaskowane.

Szukanie w logu i zawężanie zgłoszeń wagą są czynnościami wyłącznie klienckimi nad
materiałem już zebranym, bo kontrakt nie ma komendy, którą dałoby się dopytać rdzeń
o wiersze pominięte, i widok mówi to wprost. Waga zgłoszenia pochodzi z pola `severity`
oddanego przez rdzeń, nie z rozpoznawania treści wiersza logu.

Zdanie stanu pustego o tym, czym rozporządza rdzeń, składa `katalog-komend.ts` z rejestru
komend wziętego z odpowiedzi `connection.hello` — z tego, co rdzeń rzeczywiście
zarejestrował, a nie z napisu wpisanego na stałe w kliencie.

Powtórzenie uruchomienia komendą `developer.build.run` znaczy to samo żądanie, a nie to,
co aktualnie stoi w polach kreatora — od tego jest osobny przycisk uruchomienia. Gdy okno
nie wysłało jeszcze niczego samo, bo przebieg zaczęło inne okno tego konta, powtarza samo
zadanie z migawki rdzenia i mówi wprost, że parametrów zadania nie zna, ponieważ
`DeveloperBuild` ich nie niesie.

## budowa/klient-poprzedni/src/moduly/design/zrodlo-warsztatow-designu.ts

Metoda `wykonaj` jest jedną drogą na siedemdziesiąt pięć komend, ponieważ
wszystkie idą tak samo: nazwa komendy ze stałych kontraktu oraz treść żądania
złożona z formularza. Nazwa nigdy nie jest napisem wpisanym wprost — przychodzi
z katalogu czynności `czynnosci-warsztatow-designu.ts`, a ten bierze ją ze
stałych kontraktu.

Kształt żądania sprawdza rdzeń i odsyła odmowę walidacji wraz z nazwą pola.
Klient tego nie zastępuje: kontrakt rozstrzyga po stronie rdzenia, a drugie
sprawdzenie po stronie okna byłoby drugą prawdą o tym, co wolno wysłać,
i rozjechałoby się z pierwszą przy najbliższej zmianie kontraktu.

Materiał pochodzi z jednego magazynu zasobów, czytanego komendą
`design.asset.list` — tego samego, z którego czyta Assets Panel. Drugi wykaz
materiału byłby drugim miejscem, w którym ta sama treść żyje.

## budowa/klient-poprzedni/src/konfiguracja/zrodlo-obszarow-sesji.ts

Komenda `config.get` oddaje pojedyncze klucze katalogu ustawień, czyli płaskie
wpisy `ConfigEntry` ze wszystkich poziomów naraz, więc dziedziczenie jednego
klucza rozstrzyga klient w pliku `rozstrzygniecie.ts`. Komenda
`config.effective.get` oddaje obszary konfiguracji sesji już rozstrzygnięte przez
rdzeń, wraz z pochodzeniem obszaru i z rozejściem katalogu roboczego. Klient tego
rachunku wykonać nie może: nie zna pełnej ścieżki bytów ani rejestrów spoza
rodziny `config.*` — rejestru kont, rejestru kanałów, katalogu tożsamości
i nadań dostępu — z których treść obszaru jest składana.

Wynik obu czynności wraca opakowany, więc okno odróżnia obszar bez zapisu od
nieudanego zapytania.

Kształt odpowiedzi `config.effective.get` powstaje w rdzeniu w pliku
`server/internal/core/sesja_konfiguracja_skladanie.go`. Pola `config`, `origins`
i `workingDirectory` wychodzą zawsze, pole `capabilities` wyłącznie na żądanie,
a `unsupportedFields` tą ścieżką nie wychodzi. Sprawdzian obejmuje trzy pola
obowiązkowe i ani jednego więcej, ponieważ pole nieobowiązkowe w sprawdzianie
zamieniłoby zdrową odpowiedź w fałszywą odmowę. Tak samo w odpowiedzi
`config.session.set` pole `unsupportedFields` nie jest wymagane do uznania zapisu
za udany, choć okno pokazuje je, gdy przyjdzie.

Poziom i byt idą w żądaniu tym samym tłumaczeniem adresu, co komenda
`config.set`.

## budowa/klient-poprzedni/src/konfiguracja/rozstrzygniecie.ts

Rachunek dziedziczenia jest w całości po stronie klienta i opiera się wyłącznie
na tym, co przyszło z rdzenia komendą `config.get`: rdzeń oddaje wpisy wraz
z poziomem, bytem poziomu, osią i bytem osi, więc dziedziczenie da się odtworzyć
bez drugiej komendy.

Klient nie zna pełnej ścieżki bytów od środowiska po okno, więc jej nie zgaduje:
punkt widzenia wskazuje pasek u góry okna. Łańcuch takiego punktu składa się
z wpisów globalnych oraz z wpisów zapisanych dokładnie na wskazanym poziomie
i dla wskazanego bytu; zapis na tym samym poziomie, lecz dla innego bytu, dotyczy
kogoś innego i do rachunku nie wchodzi.

Brak zapisu nie jest błędem ani blokadą: obowiązuje wartość domyślna z katalogu,
a wskaźnik mówi to wprost.

## budowa/klient-poprzedni/src/dostepy/sekcja-dostepow.ts

Sekcja dotyczy dwóch bytów: punktu dostępu wraz z nadaniem, czyli tego, do jakich
maszyn i katalogów model sięga w obrębie jednego okna rozmowy, oraz katalogu
roboczego, w którym model zostawia swoje pliki. Środowiska, czyli profilu
widoczności modułów w bocznej nawigacji, sekcja nie dotyka i nie zapisuje do
niego niczego.

Układ jest dwuczęściowy: po lewej wykaz punktów wraz z dodawaniem katalogu
z „Mój komputer", po prawej zbiór nadań okna. Pod nimi stoi obszar katalogu
roboczego, oddzielony, ponieważ mówi o czym innym.

Sekcja pojawia się natychmiast, przed odpowiedzią rdzenia; puste wykazy niosą
zdanie, nie pusty prostokąt.

## budowa/klient-poprzedni/src/moduly/diagnostics/prowenancja-zrodlo.ts

Źródło jest osobne, zamiast dołożenia metod do `zrodlo-diagnostics.ts`, ponieważ
tamten plik należy do rodziny `diagnostics.*`, a prowenancja jest własną rodziną
komend kontraktu. Rozstrzygnięcie odpowiedzi jest jednak to samo — funkcja
`rozstrzygnij` z tamtego pliku — ponieważ rozróżnienie odmowy rdzenia od
odpowiedzi nieczytelnej i od odpowiedzi bez treści ma tu tę samą wagę: zakładka
poświęcona jawności pracy modeli nie może zamilczeć własnego potknięcia.

Piąta komenda rodziny, `provenance.call.replay`, do tego źródła nie należy:
powtórzenie wywołania wydaje pieniądze Operatora i jest czynnością sprawczą
osobnego odcinka, nie odczytem. Zakładka mówi o niej wprost, zamiast wołać ją po
cichu.

Odcinki są w kontrakcie obowiązkowe i puste dla wywołania bez drzewa, więc ich
brak jest odpowiedzią nieczytelną, a nie wywołaniem prostym. Treść wydania śladu
bywa pusta legalnie, gdy zakres nie obejmuje żadnego wywołania, ale format
i licznik muszą przyjść.

## budowa/klient-poprzedni/src/moduly/research/badanie-zdjecie-adnotacji.ts

Komenda `research.annotation.remove` usuwa podświetlenie albo notatkę, a rdzeń
nie ma komendy, która by je przywróciła: adnotacja niesie cytat, komentarz
i kotwicę pozycji, a po zdjęciu nie ma z czego ich odtworzyć. Dlatego czynność
mówi to przed wykonaniem i wymaga drugiego naciśnięcia tego samego chwytu.

Uzbrojenie jest wiedzą modułu, a nie okna: adnotacja wskazana w widoku czytania
i chwyt naciśnięty w panelu akcji tego samego okna dotyczą jednego badania,
a moduł prowadzi jedno badanie na sesję — stan trzyma `stan-badania.ts`. Osobne
uzbrojenie dla każdego okna pozwoliłoby uzbroić w jednym oknie, a zdjąć
w drugim, czyli bez ostrzeżenia w miejscu naciśnięcia.

Zapowiedź nie jest okienkiem dialogowym i nie odbiera klikalności: pierwsze
naciśnięcie oddaje zdanie odmowy z powodem, drugie wykonuje. Adnotacja wskazana
inna niż uzbrojona uzbraja od nowa, więc zdjęciu zawsze towarzyszy ostrzeżenie
nazywające rzecz zdejmowaną.

## budowa/klient-poprzedni/src/modele/sekcja-modeli.ts

Oś wskazania jest wspólna dla całej sekcji. Ustawienia bytu, tożsamość i podgląd
promptu mówią o tym samym bycie: wskazanym modelu albo wskazanym koncie. Trzy
osobne wybory osi w trzech panelach dałyby trzy rozbieżne wskazania, dlatego
pasek osi jest jeden i stoi nad zakładkami.

Wykaz kont zasila wybór osi: konta znane rdzeniowi stają się podpowiedziami bytu
osi `account`, a identyfikatory modeli składa osobne źródło z rejestru kanałów
oraz z modeli domyślnych kont. Podpowiedź nie zamyka pola — byt spoza wykazu
wolno wpisać wprost. Sekcja otwiera się przed odpowiedzią rdzenia, a każdy
z czterech obszarów ma własny komunikat na wypadek braku danych i żaden nie
blokuje pozostałych.

Identyfikatory modeli trwają między odczytami, ponieważ pochodzą z rejestru
kanałów czytanego osobną komendą. Zmiana rejestru kont odświeża podpowiedzi kont
i zostawia modele bez zmian: podanie pustej listy skasowałoby podpowiedź bez
powodu.

Odczyt sekcji nie rozsyła osi ponownie. Oś zmienia wyłącznie pasek u góry,
a jego zmiana już ją rozesłała; powtórzenie rozesłania przy odczycie kazałoby
rdzeniowi przeczytać te same zapisy i tę samą nakładkę drugi raz pod rząd.

## budowa/klient-poprzedni/src/dostepy/nazwy-dostepow.ts

Plik jest jedynym miejscem, w którym wartość wyliczenia kontraktu zamienia się w zdanie
po polsku. Wartości pochodzą wyłącznie ze stałych `shared/contract`, a ten plik dokłada
do nich warstwę językową i nic więcej. Rozsypanie tych zdań po widokach dałoby dwie nazwy
tego samego stanu.

Wartość spoza wyliczenia nie jest błędem i niczego nie wygasza: wraca jako własny napis,
żeby widoczne było to, co przysłał rdzeń, zamiast pustego miejsca.

## budowa/klient-poprzedni/src/dostepy/klucze-katalogu.ts

Dostęp to nie katalog roboczy. Model może mieć wgląd w katalog wskazany nadaniem
dostępu, a swoje pliki zostawiać w katalogu sesyjnym wewnątrz miejsca
instalacji. Są to dwa niezależne ustawienia i dwa osobne obszary tej sekcji,
mylenie ich prowadziłoby do nadania dostępu tam, gdzie potrzebny jest katalog
roboczy.

Katalog roboczy jest zwykłym ustawieniem: idzie komendami odczytu, zapisu
i przywrócenia wartości domyślnej konfiguracji, przez ten sam rezolwer zasięgów
co reszta ustawień. Osobnej komendy dla niego nie ma.

Katalog roboczy powstaje domyślnie w miejscu instalacji aplikacji głównej.
Klient tej ścieżki nie zna i jej nie zgaduje: mówi Operatorowi, skąd się ona
bierze, zamiast pokazywać zmyśloną ścieżkę jako fakt.

Brak wiersza w katalogu ustawień nie może zabrać możliwości ustawienia katalogu
roboczego, dlatego obszar niesie własne definicje zastępcze. Definicja
z katalogu ma pierwszeństwo, a zastępcza wchodzi wyłącznie wtedy, gdy katalog
milczy, i jest zbudowana z tych samych wartości co odpowiadające jej definicje
rdzenia.

Zasięgi katalogu roboczego są węższe niż komplet zasięgów konfiguracji: katalog
roboczy jest własnością instalacji, środowiska, projektu, sesji oraz okna, a nie
pary modułów.

## budowa/klient-poprzedni/src/moduly/agents/stan-agentow.ts

Pięć okien modułu — Agent Builder, Model Configuration, Skills Manager,
Connectors Manager, Permissions Center — pracuje na jednym zbiorze stanu. Gdyby
każde okno prowadziło własny wykaz i własny wybór, zmiana modelu bazowego nie
odświeżałaby biblioteki, a przypisanie umiejętności dotyczyłoby innego eksperta
niż ten pokazany w edytorze. Poza własnym działaniem stan odświeża wyłącznie
zdarzenie rdzenia; odpytywania w pętli tu nie ma.

Ogłoszenie zmiany stanu ustawia najpierw eksperta testowanego, dopiero potem
woła słuchaczy okien. Czat modułu jest czatem testowym bez pamięci sesyjnej:
rozmowa znika przy zamknięciu okna i przy zmianie testowanego eksperta, a okno
rozmowy stoi obok widoku modułu, więc o przełączeniu dowiaduje się wyłącznie
tym wywołaniem. Kolejność jest treścią: rozmowa zdąży się wyczyścić i nazwać
powód, zanim okna modułu przerysują się na nowego eksperta. Odwrotna kolejność
pokazywałaby przez moment czat poprzedniego eksperta w oknach już opisanych
nazwą następnego.

## budowa/klient-poprzedni/src/moduly/library/panel-akcji.ts

Panel nie jest zaszytym wykazem: pozycje przychodzą z katalogu rdzenia komendą
`action.list` o zasięgu modułu, a każda niesie własny kod, którym wykonuje ją
`window.action`. Zaszycie pozycji po stronie klienta byłoby drugą kopią
katalogu, rozjeżdżającą się z rdzeniem przy każdej zmianie po jego stronie.

Czynności, którym kontrakt nie przypisał komendy — eksport, archiwizacja, kosz,
wykrywanie duplikatów, udostępnienie odnośnikiem oraz porównanie — stoją obok
jako przyciski klikalne. Naciśnięcie takiego przycisku nic nie wysyła i nazywa
brak, zamiast zostawiać Operatora przy kontrolce wygaszonej albo milczącej.

Powodzenie wywołania oznacza wynik akcji, a nie samo przyjęcie zgłoszenia:
akcję, której rdzeń nie wykonuje, rdzeń odrzuca wprost kodem `conflict` wraz
z kodem komendy do wywołania, a panel pokazuje tę odmowę.

## budowa/klient-poprzedni/src/moduly/agents/okno-connectors-manager.ts

Wykaz mostów protokołu MCP pochodzi z tego samego katalogu punktów dostępu,
z którego rdzeń składa plik konfiguracji procesu modelu. Okno nie zakłada
drugiego rejestru mostów i nie zna ani jednego adresu maszyny: podaje kod
punktu dostępu, resztę wie rdzeń.

Katalogiem rozszerzeń zarządza osobne okno modułu, obsługujące komplet pięciu
komend tego katalogu; tutaj stan tych komend bierze się na żywo z bytu
pokrycia, a samo okno wykonuje wyłącznie podłączenie konektora i wtyczki.

Licznik narzędzi stoi także w tym oknie, ponieważ serwer narzędzi czyta
identyfikatory umiejętności i konektorów eksperta jako jeden zbiór kodów —
podłączenie konektora zmienia tę samą liczbę, którą pokazuje Agent Builder.

Wtyczki stoją osobno od konektorów: wykaz wtyczek ma własny panel i własne
komendy podłączenia oraz usunięcia. Konektor jest drogą do usługi zewnętrznej,
wtyczka katalogiem rozszerzeń powłoki — to dwa różne mechanizmy rozszerzania
eksperta.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-dziennika.ts

Wszystkie trzy czynności wiersza mają drogę w kontrakcie i wszystkie trzy naprawdę
coś robią. „Odtwórz przebieg" oddaje zamiar oknu, „Odsłuchaj nagranie" pobiera bajty
spod odnośnika wpisu komendą `speech.audio.fetch` i odtwarza je w karcie, a „Oznacz
jako ważne" zapisuje wyróżnienie w rdzeniu komendą `assistant.activity.flag`. Znacznik
przeżywa odświeżenie wykazu, ponieważ ma gdzie zamieszkać po stronie rdzenia.

Wpis bez odnośnika nagrania nie jest brakiem produktu, tylko wpisem tekstowym. Przycisk
odsłuchu mówi to wprost, zamiast milczeć albo znikać z wiersza.

## budowa/klient-poprzedni/src/moduly/multitasking/widok-strumienia.ts

Okno koordynatora pokazuje cztery rodzaje treści, a nie samą odpowiedź modelu:
tok rozumowania, wywołania narzędzi, wyniki oraz pliki. Pliki są wyróżnione
znacznikiem osobnym, ponieważ stoją obok pozostałych trzech rodzajów jako
czwarta rzecz widoczna dla koordynatora, a nie jako odmiana wyniku narzędzia.

Rodzaje fragmentów pochodzą z kontraktu przez katalog warstwy rozmowy, więc
moduł nie wprowadza własnego nazewnictwa rodzajów ani własnych wartości.

Wiersz licznika niesie powód zatrzymania biegu wprost z pola `LoopState.stopReason`.
Bez tego wiersza bieg zatrzymany byłby nieodróżnialny od biegu bezczynnego,
a koordynator nie miałby przesłanki do rozstrzygnięcia o wznowieniu.

## budowa/klient-poprzedni/src/moduly/assistant/wiersz-zlecenia.ts

Wszystkie sterowania wiersza zostają klikalne niezależnie od stanu zlecenia. Wstrzymania
zlecenia już wykonanego widok nie blokuje — odpowiada na nie rdzeń, a odpowiedź trafia do
wiersza odpowiedzi okna. Wygaszona kontrolka kazałaby zgadywać, czy przycisk nie działa, czy
tylko nie odpowiada.

Pole priorytetu wysyła sterowanie `none`, choć rdzeń zapisuje priorytet wyłącznie przy
sterowaniu innym niż `none`. Zmiana kolejności obsługi nie jest zmianą stanu zlecenia,
a kontrakt nie zna wartości `AssistantActionControl` znaczącej sam zapis priorytetu — są
wyłącznie `none`, `pause`, `resume`, `cancel` oraz `retry`, a komenda
`assistant.action.status` jest jedyną komendą obszaru niosącą pole `priority`.

Podszycie się pod `pause` albo `retry` po to, żeby przemycić priorytet, przestawiłoby stan
zlecenia, o który nikt nie prosił. Pole zostaje więc przy `none`, a rozbieżność nazywa wiersz
odpowiedzi okna składany przez `skutek-sterowania.ts`.

## budowa/klient-poprzedni/src/moduly/apps/zrodlo-punktow-dostepu.ts

Serwer MCP jest w kontrakcie pozycją katalogu rozszerzeń, a most, którym się do
niego dochodzi — punktem dostępu `AccessPoint`. Oba byty łączy pole
`Extension.accessPointId`, więc okno integracji potrzebuje obu wykazów:
z pierwszego bierze rodzaj i stan włączenia, z drugiego adres mostu, korzenie
oraz wynik ostatniego sprawdzenia.

Komenda `access.point.check` jest jedyną drogą kontraktu do testu połączenia
i do monitora zdrowia: odpowiada, czy punkt odpowiada, jakie korzenie potwierdza
i co poszło nie tak. Nie jest to sprawdzenie kondycji samego rozszerzenia —
sprawdzeniu podlega most, a nie usługa stojąca za nim. Korzenie i szczegół
bywają puste także przy sprawdzeniu udanym, dlatego całą odpowiedź niesie pole
stanu.

Źródło nie ma własnego stanu ani nie zna okna. Wykazy trzyma `stan-rozszerzen.ts`,
żeby cztery okna patrzyły na jeden zbiór danych.

## budowa/klient-poprzedni/src/konfiguracja/zasiegi.ts

Poziom zasięgu mówi, jak wąsko obowiązuje wartość: okno jest najwęższe i wygrywa,
globalny najszerszy i przegrywa z każdym innym.

Oś mówi, dla czego wartość obowiązuje: dla platformy, dla wskazanego modelu albo dla
wskazanego konta. Oś pominięta we wpisie znaczy `platform`.

## budowa/klient-poprzedni/src/moduly/library/wiersz-wersji.ts

Osiągalność treści bierze się z odpowiedzi rdzenia opracowanej w `dostepnosc-tresci.ts`,
a nie z obecności sumy kontrolnej w opisie wersji. Pole `checksum` bywa obecne przy wersji,
której treści rdzeń nie oddaje, i nieobecne przy wersji, którą oddaje w całości, więc jako
przesłanka osiągalności myli.

Odpowiedź rdzenia dotyczy treści, którą dokument niesie jako bieżącą, dlatego werdykt siada
wyłącznie na wierszu bieżącym: kontrakt nie ma komendy pytającej o treść wersji niebieżącej.
Wiersz niebieżący nie jest zatem „nieznany" z braku odpowiedzi — po prostu nie był pytany.

Werdykt `brak` znaczy, że nie ma do czego wracać. Werdykt `odmowa` znaczy, że rdzeń treści
nie oddał, ale też o niej nie orzekł, więc nie odbiera przycisku przywrócenia. Werdykt
`odwolanie` znaczy wskazanie miejsca zamiast bajtów.

## budowa/klient-poprzedni/src/konfiguracja/panel-obszarow-sesji.ts

Zasięg jest widoczny w dwóch miejscach, bo mówi o dwóch rzeczach. Nagłówek
sekcji mówi, dla kogo liczona jest konfiguracja i pod jaki adres pójdzie zapis;
oba biorą się z punktu widzenia ustawionego paskiem u góry okna, więc panel nie
stawia drugiego selektora zasięgu. Plakietka przy obszarze mówi natomiast, skąd
wartość przyszła: z którego rejestru, z którego poziomu i z której osi.

Ładowanie i odmowa idą tym samym pasem stanów odczytu, którym idą stany
katalogu, a pustka idzie stanem pustym panelu kategorii. Żadnego z tych trzech
stanów panel nie buduje sam.

Zapis wybranych obszarów zawsze kończy się zdaniem odpowiedzi. Wybór pusty nie
idzie do rdzenia, ponieważ komenda z pustym wykazem obszarów nie zapisałaby
niczego, a milczący przycisk wyglądałby jak awaria.

## budowa/klient-poprzedni/src/moduly/library/administracja-repozytorium.ts

Siedem rodzin zdolności eksperckich stoi w jednym panelu, bo wszystkie są
sterowaniem repozytorium jako całością, a nie pracą nad wskazanym plikiem.
Rozbicie ich na siedem osobnych paneli dałoby siedem miejsc, w których
Operator szukałby tego samego: gdzie ustawia się zachowanie biblioteki.

Każda czynność panelu mówi wynikiem, co zaszło. Wykazy odczytują się po
wykonaniu czynności, a nie w pętli — repozytorium nie odpytuje się samo.

## budowa/klient-poprzedni/src/moduly/agents/zrodlo-zespolow.ts

Źródło jest osobne od `zrodlo-agentow.ts`, ponieważ zespół nie jest polem
eksperta ani jego wersją: jest własnym bytem kontraktu `Team`, z własnym
zdarzeniem `team.changed` i własnym cyklem życia. Włączenie tych czterech
komend do źródła obszaru `agent.*` dałoby jeden plik o dwóch
odpowiedzialnościach.

Nazwy `Command.*` kończą się w tym pliku. Panel oraz widok składu wołają
czynności interfejsu nazwane po polsku, bytami dziedziny, i kontraktu nie
dotykają. Przemianowanie komendy w `contract.json` przerywa wtedy kompilację
w jednym pliku, a nie w trzech widokach.

Żadna czynność źródła nie rzuca wyjątku. Odmowa wraca polem `blad` wyniku
i pokazuje się w wierszu odpowiedzi tego panelu, w którym czynność wywołano.

Rozróżnienie między założeniem nowego zespołu a zmianą istniejącego opiera się
wyłącznie na obecności identyfikatora, więc pole idzie do żądania tylko wtedy,
gdy naprawdę jest. Wysłanie pustego napisu byłoby dla rdzenia wskazaniem
zespołu o nazwie pustej, a nie brakiem wskazania.

Opis jest w kontrakcie polem opcjonalnym, więc pusty nie trafia do żądania:
przy zmianie zespołu wysłanie pustego napisu skasowałoby opis już zapisany,
choć formularz o skasowanie nie prosił.

Zdanie o zespole podaje liczebność składu, ponieważ jest ona jedyną rzeczą, po
której Operator poznaje na wykazie różnicę między zespołem a jego kopią świeżo
powieloną.

## budowa/klient-poprzedni/src/dostepy/przelacznik-trybu.ts

Ostrzeżenie o zapisie jest częścią przełącznika, ponieważ dotyczy skutku wyboru,
a nie sąsiedztwa na ekranie. Zapis na maszynie chronionej pokazuje pełne zdanie pod
przełącznikiem, a nie w podpowiedzi pod kursorem: skutek zapisu na produkcyjnej
platformie LEX jest nieodwracalny, więc zdanie musi być widoczne bez najeżdżania
wskaźnikiem.

Treść ostrzeżenia liczy się z nazwy maszyny podanej w punkcie dostępu, dlatego
nanoszenie nowego punktu przelicza je od nowa.

## budowa/klient-poprzedni/src/modele/rodzaje-kont.ts

Plik nie wymienia żadnego dostawcy. Dostawca jest w kontrakcie wartością danych,
a nie typem kodu, i wpisuje się go w polu tekstowym formularza konta. Rodzaj
konta pozostaje natomiast wyliczeniem kontraktu, więc jego wartości pochodzą
wyłącznie ze stałych `AccountKind`.

Rodzaj spoza kontraktu nie gaśnie — nazwa gotowa do wydruku pokazuje wtedy
własny kod wartości, żeby było widać, co przyszło z rdzenia. Ta sama zasada
obowiązuje kontrolkę wyboru: napis pusty znaczy brak ograniczenia, a napis
spoza kontraktu również, ponieważ zawężenie wykazu do wartości, której rdzeń nie
zna, dałoby wykaz pusty bez powodu.

## budowa/klient-poprzedni/src/moduly/assistant/stan-probny.ts

Sprawdzianów modułu są dwa — stany okien oraz wykonywanie pracy — a atrapa stanu modułu jest
jedna. Skopiowana do obu plików rozjeżdżałaby się przy każdym nowym polu interfejsu i jeden
ze sprawdzianów badałby wtedy stan, którego moduł już nie ma.

Plik nie należy do produktu: sięgają po niego wyłącznie pliki sprawdzianów, więc punkt wejścia
aplikacji go nie wciąga.

Atrapa nie sięga do rdzenia i nie udaje, że sięga: każda droga wywołania rzuca wyjątkiem,
dopóki sprawdzian jej nie obsadzi. Sprawdzian, który przypadkiem wywoła komendę, dostaje przez
to błąd, a nie ciszę. Dotyczy to również trzech źródeł dobudowanych obok rdzenia modułu —
atrapa oddająca pustą odpowiedź udawałaby wynik, którego nie ma.

## budowa/klient-poprzedni/src/moduly/agents/pola-zalezne.ts

Adapter dostawcy deklaruje dla każdego pola obszaru `model`, czy obsłuży je w całości,
częściowo, czy wcale, i podaje powód. Okno pokazuje tę deklarację wprost, więc widać,
że kanał wiersza poleceń zignoruje `samplingTemperature`, zanim wartość zostanie wpisana.

Wykaz zbudowany z deklaracji, a nie z listy wpisanej w kodzie, nadąża za dostawcą:
nowa wersja jego programu zmienia zestaw pól bez wydania nowego klienta.

## budowa/klient-poprzedni/src/moduly/multitasking/indeks.ts

Scena modułu obsadza cztery okna trzech typów. Okno koordynatora planuje bieg
i steruje nim, lecz nie tworzy produktu końcowego. Dwa okna wykonawców pracują
równolegle i nie zarządzają procesem. Okno analityka zestawia wyniki obu
wykonawców. Wszystkie cztery okna stoją na jednym stanie i jednym źródle biegu,
ponieważ przy osobnych źródłach koordynator widziałby inny przebieg niż
wykonawca, którym steruje.

Wykaz komend rdzenia powstaje raz na scenę i jest wspólny czterem oknom. Cztery
osobne zapytania o ten sam wykaz dałyby cztery odpowiedzi, które mogłyby się
rozejść. Zapytanie wychodzi przed pierwszym odświeżeniem okien, dzięki czemu
kontrolki bez pokrycia w wykazie oznaczają odpowiedź w drodze, a nie brak
komendy w rdzeniu.

Panel obsady stoi przed oknami ról, ponieważ role są przypisaniem okien sesji,
a nie osobnymi bytami, i bez przydziału nie ma czego otworzyć. Panel jest
jedynym miejscem nadania tych ról, a jego meldunki idą do jego własnego stanu
treści, więc odmowa rdzenia nie ginie po drodze. Układ sekcji paneli pochodzi
z powłoki, ponieważ dotyczy okna, a nie dziedziny modułu.

Panel zadań w tle stoi pod oknami ról, ponieważ opisuje pracę już przez nie
zleconą, a podagenci należą do okien wykonawców. Panel jest jeden na kartę
sesji, ponieważ komenda wykazu podagentów przyjmuje identyfikator sesji i dwa
panele pytałyby dwukrotnie o to samo. Oba okna wykonawców korzystają z tego
samego źródła podagentów co panel, więc drugie źródło byłoby kopią.

Nasłuchy odpina się przed zdjęciem sceny, ponieważ zdarzenie przyjęte po
usunięciu elementu odświeżałoby widok już nieistniejący.

## budowa/klient-poprzedni/src/moduly/diagnostics/telemetria-procesow.ts

Zakres czasu nie dotyczy tego odczytu: `monitor.status` i `monitor.subscribe`
nie mają w kontrakcie pól `fromTime` ani `toTime`, więc telemetria jest zawsze
stanem bieżącym rejestru procesów, a zawężanie jej zakresem po stronie okna
byłoby zawężaniem pozorowanym. Obserwacja zakłada się raz, odczyt powtarza
się: pierwsze pytanie idzie `monitor.subscribe` (zapisuje okno na telemetrię),
każde następne `monitor.status` (czyta i niczego nie zapisuje) — pole
`windowId` znaczy w tych dwóch komendach co innego: w pierwszej jest oknem
obserwatora, w drugiej sitem procesów, więc do odczytu nie idzie wcale.

Na żywo idzie `progress.changed`. Zdarzenie niesie stan, etap i stopień
ukończenia, lecz nie niesie czasu zmiany, dlatego wiersz odświeżony
zdarzeniem pokazuje liczby ze zdarzenia, a czas nadal ten z odczytu, wraz ze
zdaniem, skąd się bierze — podstawienie czasu klienta w miejsce czasu rdzenia
byłoby wpisaniem do telemetrii wartości, której rdzeń nie orzekł.

Eksport telemetrii stoi przy danych, a nie w pasku akcji okna: potwierdzenie
eksportu wypisuje się w tej samej zakładce, w której widać eksportowany
materiał, bo potwierdzenie na zakładce zasłoniętej nie byłoby potwierdzeniem.
Wytwórnia sama nie czyta — pierwszy odczyt zleca złożenie modułu wywołaniem
`odswiez()`, tak samo jak w panelu rekomendacji, żeby odczyt w wytwórni obok
odczytu ze złożenia nie dał dwóch żądań na jedno zmontowanie okna.

Po odmowie odczytu wykaz poprzedni też odchodzi: skoro odczyt się nie udał,
wykaz przestaje być bieżący, a eksport oddałby go, jakby wciąż nim był.

## budowa/klient-poprzedni/src/aplikacja/rejestr-modulow.ts

Nawigacja bierze wykaz modułów z rdzenia komendami `module.list`
oraz `environment.enter`, więc powłoka poznaje kody modułów dopiero w czasie
działania i nie może o nich nic zakładać. Bez rejestru każdy moduł musiałby
wpinać się sam, a przestrzeń robocza trzymałaby własną kopię katalogu modułów.
Moduł nieobecny w rejestrze nie jest błędem: rdzeń zna moduły, których klient
jeszcze nie zbudował, a taki moduł dostaje stan pusty z paska uczciwości
przestrzeni roboczej zamiast martwego kliknięcia.

Element widoku i wczytanie są rozdzielone, ponieważ sesja bywa nieznana w chwili
wyboru pozycji nawigacji: przestrzeń robocza odkłada wtedy wejście i ponawia je
po otwarciu okna przez rdzeń. Gdyby wytwórnia oddawała sam element, każdy moduł
obsługiwałby tę samą zwłokę u siebie.

Metoda `zamknij` jest szwem rozbiórki widoku i pozostaje opcjonalna. Powłoka jej
nie woła, ponieważ `przestrzen-modulu.ts` woła `utworzWidok` raz na kod modułu
i oddaje ten sam widok przy każdym powrocie, a `obszar-roboczy.ts` nie zdejmuje
elementu z drzewa, tylko przestawia atrybut `hidden`. Liczba żywych subskrypcji
jest przez to ograniczona z góry liczbą modułów i nie rośnie wraz z liczbą
przełączeń. Gdy w powłoce pojawi się pierwsza granica życia widoku, wołanie
`zamknij` musi towarzyszyć zdejmowaniu elementu w `obszar-roboczy.ts`.

Wpis MultitaskingAI niesie kod `multitaskingai`, który jest kodem środowiska,
a nie modułu. Nawigacja tego środowiska podaje sekcje orkiestracji ze stałej
`SEKCJE_ORKIESTRACJI` w `powloka/srodowiska.ts`, więc wyszukanie opisu modułu
po kodzie środowiska nie daje wyniku.

Moduł Diagnostics montuje się od razu, ponieważ większość jego komend nie
wymaga identyfikatora okna.

## budowa/klient-poprzedni/src/moduly/developer/zrodlo-probne.ts

Jedna atrapa dla całego modułu zastępuje atrapę przepisywaną w każdym sprawdzianie z osobna:
rozjazd z umową `ZrodloDeveloper` przerywa wtedy kompilację, zamiast rozjeżdżać sprawdziany
po cichu. Plik służy wyłącznie sprawdzianom tego modułu i nie jest wciągany przez żadne okno.

## budowa/klient-poprzedni/src/moduly/developer/stan-developer.ts

Cztery okna modułu pracują nad tym samym repozytorium: drzewo projektu wskazuje
plik, edytor kodu go otwiera i zapisuje, panel repozytorium podstawia ścieżkę
wskazaną do pola czynności i sam wskazuje ścieżki ze swojego wyniku, a okno
wyniku budowania wskazuje pliki ze zgłoszeń przebiegu. Gdyby każde okno trzymało
własną ścieżkę i własny korzeń, wskazanie pliku w drzewie nie dotarłoby do
edytora.

Identyfikator okna modułu należy do stanu wspólnego, ponieważ kontrakt wymaga go
w każdej komendzie obszaru, a okna mają go podawać identycznie. Rozbieżność
rozdzieliłaby ich pracę między dwa katalogi robocze rdzenia.

Stan nie wywołuje komend w zastępstwie okien: przechowuje wybór i rozgłasza
zmianę, a odczyt i zapis pozostają w oknach, ponieważ tylko one prowadzą stan
ładowania i obsługę odmowy.

Ostatni plik oddany przez rdzeń nie jest tym samym co ścieżka wskazana:
wskazanie biegnie natychmiast po wskazaniu węzła w drzewie, a plik przychodzi
dopiero odpowiedzią rdzenia. Zgodność obu rozstrzyga porównanie ścieżki pliku ze
ścieżką wskazaną i dlatego plik niesie własną ścieżkę. Edytor kodu czerpie stąd
identyfikator wersji sprzed zapisu, a drzewo projektu odróżnia węzeł wskazany od
węzła wczytanego do edytora.

Treść pliku pochodząca z rdzenia nie jest treścią pola edycji. Pole bywa
zmienione i niezapisane, a rdzeń po zapisie nie zawsze odsyła treść, więc
znacznik czystości pola należy do edytora kodu.

Wskazanie nowej ścieżki nie zeruje pliku: plik niesie własną ścieżkę, więc nie
sposób pomylić go ze wskazaniem, a wyzerowanie usunęłoby jedyną wiedzę o tym, co
leży w polu edytora.

Przyrost budowania przychodzi także z pracy innego okna tego konta, dlatego stan
przyjmuje wyłącznie przyrost dotyczący własnego okna modułu. Budowanie tworzy
pliki, więc po nim drzewo i plik w edytorze bywają nieaktualne i wymagają
ponownego odczytu.

## budowa/klient-poprzedni/src/moduly/browser/material-sesji.ts

Zbiór materiału nie jest odbiciem stanu rdzenia, ponieważ odbijać nie ma czego:
zapis pozycji jako wytworu sesji ma komendę `browser.artifact.add`, natomiast
komendy odczytu wykazu wytworów okna kontrakt nie niesie. Pozycją jest migawka,
którą rdzeń oddał na żądanie, a po przeładowaniu karty wykaz zaczyna się od nowa;
same migawki zostają w rdzeniu pod swoimi identyfikatorami. Okno mówi o tym
wprost, zamiast obiecywać trwałość, której nie ma.

Rodzaj pozycji rozstrzyga zawartość odpowiedzi, a nie zamówienie okna. Rdzeń
pobiera stronę bez uruchamiania przeglądarki (`przegladarka_pobieranie.go`), więc
migawka zamówiona ze zrzutem potrafi przyjść bez niego — pozycja nazywa się
wtedy archiwum albo treścią. Rozpoznanie po odpowiedzi nazwie ją zrzutem samo,
gdy obraz zacznie przychodzić.

Monitor zmian stoi na adresie, a nie na migawce: sprawdzenie polega na ponownym
pobraniu tej samej strony i zestawieniu jej treści z zapamiętaną. Adres jest więc
jedynym rozsądnym kluczem monitora, bo dwa monitory tej samej strony pilnowałyby
dokładnie tego samego. Założenie monitora na adresie już pilnowanym oddaje
wartość fałszywą, żeby okno powiedziało o tym wprost zamiast meldować czynność,
która się nie odbyła.

## budowa/klient-poprzedni/src/aod/zrodlo-decyzji.ts

Źródło stoi osobno od `zrodlo-komend.ts`, ponieważ tamto niesie wyłącznie rodzinę
komend nakładki. Sklejenie obu zatarłoby granicę między komendami własnymi nakładki
a komendami cudzych obszarów, po które nakładka sięga.

Odczyt procesów idzie bez argumentów świadomie: nakładka pyta o wszystkie procesy,
które rdzeń zna, bo proces czekający na decyzję nie musi należeć do sesji tego klienta.

Konfiguracja okna otwierana jest na zasięgu `window`, czyli poziomie najwęższym,
który wygrywa z każdym szerszym. Konfiguracja koordynatora dotyczy tego jednego okna
i nie ma sięgać ustawień sesji ani konta.

## budowa/klient-poprzedni/src/moduly/agents/panel-zakresu-eksperta.ts

Grupy dostępu do rozszerzeń i dostępu do modułów są wartościami wyliczenia
uprawnień i stoją wyżej jako zwykłe wiersze Permissions Center. W tym panelu
mieszkają trzy rzeczy, które wartością logiczną nie są: moduły zastosowania
jako wykaz kodów, gdzie wykaz pusty znaczy brak ograniczenia i jest stanem
wyjściowym; izolacja techniczna jako macierz ośmiu suwaków, identyczna
z macierzą okna konfiguracji punktów izolacji, bo macierz jest jedna;
oraz para ustawień sieci podagentów — czynność i górna liczba podagentów.

Podgląd polityki efektywnej pokazuje wynikowy zestaw ustawień po uwzględnieniu
wszystkich trzech grup. Jest to podgląd, nie bramka — komenda niczego nie
zapisuje i niczego nie rozstrzyga.

Każde zawężenie ustawione w tym panelu jest zawężeniem egzekwowanym: rdzeń
czyta te zapisy przy nakładaniu eksperta na okno i przy powołaniu podagentów,
i na ich podstawie odmawia. Panel mówi o tym Operatorowi wprost, żeby żaden
suwak nie został wzięty za ozdobę.

## budowa/klient-poprzedni/src/moduly/assistant/panel-kontekstow.ts

Kontrakt niesie jeden mechanizm przełączania pamięci w trakcie pracy: komenda
`memory.toggle` ustala, które poziomy zasięgu karta sesji widzi i czy zapis
pamięci jest czynny. To jest kontekst, który rzeczywiście da się przełączyć,
i tak też okno go nazywa.

Nazwanych zestawów pamięci okno nie udaje. Kontrakt nie ma bytu, w którym taki
zestaw miałby zamieszkać: obszar `memory.*` zna wpis oraz poziom zasięgu, a nie
nazwany profil pamięci. Brak nazywa przycisk obok, zamiast pola, którego rdzeń
nie zapisze.

Przestawienie idzie jednym wywołaniem obejmującym oba pola naraz. Rdzeń
odczytuje pole `levels` jako stan po przestawieniu, a pusta lista zostawia
poziomy bez zmian, więc wysyłanie samego jednego przełącznika zamieniałoby
resztę wykazu w brak zmian przy każdym kliknięciu, a Operator odznaczający
ostatni poziom nie miałby jak wyłączyć wszystkich.

Stan pokazany po przestawieniu bierze się z odpowiedzi rdzenia: kontrolki
odbijające samo zamówienie okna kłamałyby przy każdym takim wywołaniu.

## budowa/klient-poprzedni/src/moduly/assistant/skutek-sterowania.ts

Plik odpowiada wyłącznie za przekład odpowiedzi rdzenia na zdanie dla
czytającego. Stoi poza oknem monitora, ponieważ okno składa tabelę, a to jest
ocena tego, co rdzeń oddał.

Reguła jest jedna: porównaj zamówienie ze zleceniem, które wróciło. Odpowiedź
niesie zlecenie po zmianie, a przekład sterowania na stan stoi w rdzeniu
w jednej funkcji, więc sprawdzenie nie kosztuje nic i obowiązuje tak samo dla
priorytetu, jak dla czterech przycisków panelu akcji. Zamówienie i odpowiedź
rozchodzą się z dwóch powodów. Komenda zapisuje priorytet wyłącznie przy
sterowaniu innym niż odczyt, a sama zmiana kolejności idzie torem odczytu
i wraca ze zleceniem bez zmian. Wykonawca zlecenia domyka je natomiast własnym
zapisem, a rdzeń oddaje wiersz odczytany po zapisie, więc zapis wykonawcy
potrafi wejść między zapis sterowania a jego odczyt.

Zdanie odmowy powstaje wyłącznie z rozbieżności między zamówieniem
a odpowiedzią rdzenia. Okno nie orzeka o braku, którego rdzeń nie pokazał, tak
samo jak nie potwierdza skutku, którego rdzeń nie oddał.

Słownik stanów zamówionych jest kluczowany stałymi kontraktu, więc dopisanie
sterowania przerwie kompilację w tym miejscu, zamiast po cichu wpaść w gałąź
braku zamówienia. Zatwierdzenie bramy potwierdzeń puszcza zlecenie w bieg, bo
czekało na rękę operatora, a nie na zasób; odmowa zdejmuje je z kolejki i jest
anulowaniem z powodem, nie osobnym stanem.

## budowa/klient-poprzedni/src/moduly/design/porownanie-wariantow.ts

Porównanie zestawia pola opisowe zasobów, a nie ich obrazy. Pole `uri` zasobu
wskazuje ścieżkę w systemie plików rdzenia, a droga po bajty obrazu, choć opisana
w kontrakcie, nie ma jeszcze uchwytu w rdzeniu, więc przeglądarka nie ma czym
wczytać żadnej ze stron porównania. Suwak przejścia między stroną „przed"
a stroną „po" nabiera znaczenia dopiero po dobudowaniu tej obsługi.

## budowa/klient-poprzedni/src/moduly/browser/ster-wyboru.ts

Trzy nastawy modułu Browser — rodzaj wyodrębnienia w pasku dolnym oraz powiązane
źródło i moduł docelowy w formularzu notatki — korzystają z jednego mechanizmu
rozwijania z komponentu menu-drzewa. Rozwijanie, filtrowanie, haczyk przy pozycji
i wędrówka klawiszami należą do tego mechanizmu, a ster wyłącznie go obsadza
pozycjami i odbiera wybór.

Wartość nastawy przechowuje jedno pole. Napis na uchwycie oraz haczyk przy
pozycji biorą się z tego pola przy każdym przerysowaniu, a odczyt wartości oddaje
to samo pole, więc uchwyt i wykaz nie mogą pokazać dwóch różnych nastaw.

Podpis nad sterem jest blokiem tekstu, a nie etykietą formularza. Etykieta bez
wskazania pola związałaby się z pierwszym potomkiem dającym się etykietować,
czyli z uchwytem menu. Kliknięcie w podpis otwierałoby wtedy wykaz, a nazwa
dostępna uchwytu konkurowałaby z opisem dostępnym, który mechanizm menu składa
sam z nazwy nastawy i wartości bieżącej.

Podpis pojawia się na ekranie tam, gdzie ster sąsiaduje z polami formularza,
ponieważ sama wartość nie mówi, czego dotyczy. Przy pasku dolnym nazwa nastawy
trafia wyłącznie do opisu dostępnego uchwytu.

Ustawienie nastawy z pominięciem Operatora oddaje informację o przyjęciu.
Wartość spoza wykazu nie zostaje na uchwycie: ster wraca do pozycji pierwszej
i oddaje odpowiedź odmowną. Milczące przyjęcie wartości nieobecnej w wykazie
dałoby uchwyt pokazujący nastawę, której nie da się wybrać, więc rozstrzygnięcie
o powiadomieniu Operatora należy do wołającego. Z tego samego powodu wymiana
wykazu utrzymuje wybór tylko wtedy, gdy jego pozycja nadal w wykazie stoi,
a w przeciwnym razie nastawą staje się pierwsza pozycja nowego wykazu.

## budowa/klient-poprzedni/src/aplikacja/aplikacja.ts

Gotowość Centrum dowodzenia jest domknięciem pierwszego odczytu strony. Punkt
wejścia podaje ją scenie wejścia w katalogu `ladowanie/`, żeby ekran między
bramką a produktem schodził w chwili, gdy strona ma czym stanąć. Nieodwiedzona
trasa strony głównej znaczy gotowość natychmiastową: nie ma na co czekać, skoro
widoku nikt nie zbudował.

Motyw rusza przed montażem czegokolwiek, żeby dokument dostał żetony motywu,
zanim pojawi się pierwszy element. Router buduje widok leniwie, przy pierwszym
wejściu na trasę, i już go nie porzuca — karty opuszczonego środowiska trwają
w tle i wracają z pełnym stanem. Środowisko przygotowuje się od razu, ponieważ to
ono podpina przepływ komunikatów.

Uzgodnienia z rdzeniem złożenie aplikacji nie rozpoczyna: robi to przepływ
komunikatów, gdy transport zgłosi stan połączenia, a wywołanie stąd dałoby drugie
powitanie. Tożsamość klienta pochodzi z powitania, ponieważ komendy `session.bind`
i `session.focus` żądają identyfikatora klienta przedstawionego rdzeniowi raz.

Pas kart sesji powstaje w głębi powłoki, która o rdzeniu nie wie, bo buduje ją
także stanowisko podglądu; drogę do rdzenia podaje jej złożenie aplikacji.

## budowa/klient-poprzedni/src/aod/progi-aod.ts

Wszystkie progi funkcji globalnej Always On Display są ustawieniami zasięgu tej
funkcji i podlegają zmianie z okna konfiguracji, z okna Ustawień oraz poleceniem
języka naturalnego. Kontrakt nie niesie ani kategorii ustawień tej funkcji, ani
komendy zapisującej te wartości — dlatego wartości domyślne stoją w pliku jako
stałe, a nie jako odczyt komendy `config.get`. Rozjazd jest zgłoszony, nie
zasypany atrapą.

Chwilę końca wyciszenia do końca dnia liczy zegar maszyny Operatora, ponieważ
funkcja stoi na jego stanowisku, a nie w strefie czasowej rdzenia.

## budowa/klient-poprzedni/src/moduly/design/stan-okna.ts

Komunikat stanu przesłania treść okna, a nie zastępuje jej: nieudane odświeżenie
nie kasuje tego, co Operator już widział, i treść wraca nietknięta, gdy okno
wróci do fazy gotowej. Odmowa rdzenia jest widoczna w oknie, w którym Operator ją
wywołał, a nie wyłącznie w konsoli.

Wygląd bierze się w całości z biblioteki `komponenty/`, ze stanu pustego oraz ze
wskaźnika obrotowego, więc plik nie zna ani jednej barwy. Nazwy faz i znakowanie
powłoki są wspólne dla wszystkich modułów w pliku `komponenty/faza-okna.ts`;
tutaj zostaje wyłącznie to, czym moduł Design różni się świadomie: wskaźnik
odczytu i chowanie treści na czas ładowania.

Kreator Prompt Buildera stoi na natywnym elemencie `dialog` otwartym wywołaniem
`showModal`, a modal otwarty czyni resztę dokumentu bezwładną. Schowanie treści
okna na czas ładowania schowałoby też ten modal — bezwładność by została, a strona
nie miałaby ani jednej drogi wyjścia. Dlatego treści schować nie wolno, gdy w niej
stoi otwarty modal.

## budowa/klient-poprzedni/src/moduly/diagnostics/zrodlo-obserwowalnosci.ts

Rodzina `monitor.*` jest jedynym źródłem obserwowalności, jakie kontrakt niesie
poza obszarem `diagnostics.*`. Metryki wydajności i kontrola stanu samej platformy
wywodzą się z encji `proces_sesji`, a struktura `MonitorStatus` jest dokładnie jej
odwzorowaniem w kontrakcie. Zakładki Metrics & Performance oraz Health & Uptime
jadą więc jednym odczytem: są to dwie perspektywy na ten sam materiał, a nie dwa
niezależne pomiary, które rozejdą się po pierwszej zmianie.

Odmowa komendy `monitor.status` bywa tu zjawiskiem zwykłym: rdzeń bez wpiętego
rejestru telemetrii odmawia głośno w pliku `core/handlers_monitor.go`, a wykaz
pusty ukryłby ten fakt pod zdaniem o braku procesów.

Podział na dwie komendy jest podziałem ról, nie wygodą: `monitor.status` czyta
stan bieżący i niczego nie zapisuje, a `monitor.subscribe` dodatkowo zapisuje okno
na telemetrię. Pole `subscribed` odpowiedzi mówi, czy zapis doszedł do skutku —
żądanie bez identyfikatora okna jest zwykłym odczytem, obsłużonym w pliku
`core/adapter_modul_monitor.go`, a moduł Diagnostics montuje się bez okna.

Zdarzenie `progress.changed` dochodzi do wszystkich połączeń konta niezależnie od
tego, czy komenda `monitor.subscribe` zapisała okno; zapis mówi rdzeniowi
wyłącznie, które okno których procesów pilnuje. Dlatego wykaz zmienia się na żywo
także wtedy, gdy pole zapisu niesie wartość fałszywą.

## budowa/klient-poprzedni/src/mission-control/pas-relacji.ts

Kontrakt nie niesie odczytu powiązań między sesjami. Komplet danych pulpitu podaje
w tym miejscu `null`, a pas mówi o braku źródła. Jest to co innego niż pusty wykaz:
zdanie „brak powiązań" byłoby twierdzeniem, którego nie da się odczytać z rdzenia,
więc pas rozdziela oba stany i nazywa je osobnymi zdaniami.

Znak relacji jest podwójny — strzałka `⇄` albo `→` oraz słowo w podpowiedzi i w treści
czytanej. Kierunek powiązania nie zależy wtedy od samego kształtu znaku, który dla
odczytu ekranowego jest niedostępny.

## budowa/klient-poprzedni/src/moduly/agents/poziomy-pamieci.ts

Konfiguracja pamięci jest czwartym z siedmiu komponentów definicji eksperta.
Poziomy są cztery i wybiera się je niezależnie, ponieważ ekspert korzysta naraz
z kilku: pamięci globalnej Operatora, pamięci projektu, pamięci karty sesji oraz
pamięci środowiska.

Piątego poziomu oznaczającego wyłączenie nie ma i kontrolka go nie dorabia:
wyłączeniem jest zbiór pusty, czyli cztery pola odznaczone. Kontrakt rozstrzyga
to wprost przy wyliczeniu `MemoryLevel` — piąta wartość obok czterech poziomów
pozwalałaby zapisać stan sprzeczny, w którym sesja jest zaznaczona równocześnie
z wyłączeniem.

Dlatego pod grupą stoi zdanie czytające stan bieżący, a nie stała treść:
Operator ma widzieć, że odznaczenie wszystkiego jest wyłączeniem pamięci, zanim
naciśnie zapis. Zdanie zmienia się przy każdym kliknięciu.

## budowa/klient-poprzedni/src/komponenty/faza-okna.ts

Zestaw faz nie zna wartości spoczynku. Zdanie o tym, że rdzeń nie był jeszcze
pytany, opisuje stan źródła danych, a nie stan okna, i ma własny typ
`FazaOdczytu` w `dostepy/stan-dostepow.ts` oraz w stanach pozostałych modułów.
Moduły odwzorowują go na okienne `puste` z osobną treścią komunikatu — tak
robią `moduly/design/okno-assets-panel.ts` oraz
`moduly/library/okno-library-explorer.ts`.

Widoczność wskaźnika odczytu i chowanie treści na czas ładowania zostają
w modułach i różnią się między nimi celowo: Assistant oraz Apps zostawiają treść
widoczną w ładowaniu, Design i Research ją chowają.

Rola `alert` należy się wyłącznie odmowie, ponieważ przerywa czytnikowi ekranu
bieżącą wypowiedź; pozostałe fazy idą jako `status`. Pas komunikatu znika tylko
w fazie `gotowe` — stan przesłania treść, a nie kasuje jej, więc po powrocie do
fazy `gotowe` widać to, co stało na ekranie przed nieudanym odświeżeniem.

## budowa/klient-poprzedni/src/konfiguracja/wybor-adresu.ts

Ta sama kontrolka obsługuje dwie role okna konfiguracji. Pasek u góry ustawia
punkt widzenia, względem którego liczone jest dziedziczenie pól, a panel przy
polu ustawia adres zapisu wartości nadpisującej. Różni je wyłącznie wykaz
dopuszczalnych poziomów i osi: pasek podaje wszystkie, panel pola tylko te,
które dopuszcza katalog ustawienia.

Pola bytu znikają tam, gdzie poziom albo oś bytu nie mają, czyli przy poziomie
globalnym i przy osi platformy.

## budowa/klient-poprzedni/src/moduly/design/karta-zasobu.ts

Karta zasobu jest przyciskiem, a nie prostokątem, ponieważ wybór zasobu
przestawia naraz panel metadanych, kanwę Design Board oraz pole obrazu
referencyjnego Prompt Buildera. Element musi być z tego powodu osiągalny
klawiaturą. Zasób niesie pole `uri` tylko wtedy, gdy rdzeń zna jego położenie;
bez tego pola karta pokazuje pole zastępcze z rodzajem zasobu zamiast miniatury.

Słownik nazw rodzajów pokrywa cały typ kontraktu, więc brak wartości zatrzymuje
kompilację. Wynik pracy modelu wchodzi do wykazu zasobów tą samą drogą, co zasób
wniesiony ręcznie, dlatego słownik obejmuje także rodzaje wytwarzane przez
model — bez nich karta pokazałaby pustkę w miejscu nazwy.

Zawężanie wykazu frazą jest miejscowe i nie jest polem żądania. Komenda
`design.asset.list` zawęża polami identyfikatora okna, rodzaju, etykiet,
znacznika ulubionego oraz ograniczenia liczby wyników, a frazy wśród nich nie
ma. Predykat zawężania stoi w jednej kopii, ponieważ czytają go dwa widoki tego
samego zbioru: okno Assets Panel oraz panel `zasoby-designu` stosu paneli
pomocniczych. Dwie kopie dawałyby na tę samą frazę różne wyniki.

Etykieta zasobu nie niesie stanu, więc wchodzi plakietką bazową. Barwna odmiana
plakietki pozostaje zarezerwowana dla plakietek znaczących stan.

## budowa/klient-poprzedni/src/moduly/automations/stan-automatyki.ts

Pięć okien modułu pracuje nad tą samą automatyką: Workflow Builder buduje jej kroki,
Orchestrator układa między nimi zależności, Scheduler nadaje jej cykliczność, Queue
Manager uruchamia jej kolejkę, a Execution Monitor pokazuje jej przebiegi. Wskazanie
trzymane osobno w każdym oknie rozjeżdżałoby się przy pierwszej zmianie automatyki.

Kolejka bieżąca stoi obok automatyki, ponieważ Queue Manager musi wiedzieć, którą
kolejkę posuwa, a Execution Monitor — po której przyszedł stan przebiegu. Kolejka jest
bytem silnika kolejek, a nie definicji, więc nie należy do struktury `AutomationWorkflow`
i nie da się jej z niej odczytać.

Stan przebiegu przychodzi także z pracy innego okna albo innego urządzenia tego konta.
Gdy dotyczy automatyki bieżącej, definicja trzymana w oknie jest już nieaktualna
i okna odczytują ją ponownie.

## budowa/klient-poprzedni/src/moduly/apps/formularz-komponentu.ts

Identyfikator komponentu składa się z jego nazwy, o czym mówi objaśnienie pola. Kontrakt
wymaga pola `id` w każdym komponencie żądania, a rdzeń nadaje własne oznaczenie dopiero
w odpowiedzi. Identyfikator roboczy jest więc kluczem zależności wewnątrz jednego zapisu,
a nie obietnicą trwałości.

Rodzaj komponentu wybiera się rozwijaniem z biblioteki przez obsadę `wybor-z-menu.ts`,
a nie kontrolką natywną przeglądarki, żeby wykaz rodzajów wyglądał tak samo jak pozostałe
wykazy wyboru w oknie.

## budowa/klient-poprzedni/src/dostepy/indeks.ts

Reszta aplikacji zna z tego katalogu jedną czynność — `otworzOknoDostepow`,
wołaną z kanałem rdzenia i identyfikatorem okna rozmowy. Gospodarz, który ma
własny obszar na ekranie i nie chce okna modalnego, bierze zamiast tego samą
sekcję przez `utworzSekcjeDostepow` i osadza jej element u siebie; czynność
`ustawOkno` wiąże ją wtedy z oknem rozmowy.

Okno jest jedno na klienta i żyje między otwarciami, tak samo jak okno
konfiguracji. Powtórne otwarcie wraca do stanu, na którym się skończyło, i nie
gubi subskrypcji zdarzeń `access.*.changed`. Wywołanie z innym kanałem, po
ponownym połączeniu z rdzeniem, buduje okno od nowa, żeby komendy nie szły przez
transport, którego już nie ma.

Nadanie dostępu żyje w obrębie okna, więc identyfikator okna podaje się przy
otwarciu. Otwarcie bez niego jest poprawne: wykaz punktów i katalog roboczy nie
zależą od okna, a w miejscu nadań sekcja mówi, na co czeka.

## budowa/klient-poprzedni/src/moduly/library/archiwum-pakowanie.ts

Spakowanie archiwum stoi w obszarze archiwum panelu metadanych obok utrwalenia
i paczki migracyjnej, lecz jest od nich oddzielone osobnym zdaniem, ponieważ
pracuje na innym zbiorze. Komendy `library.preservation.run`
i `library.package.export` obejmują zasoby repozytorium biblioteki i oddają wynik
jako zasób biblioteki, natomiast `archive.pack` pakuje katalog albo plik ze
stanowiska operatora i oddaje wynik jako zasób magazynu Designu. Nazwanie jednej
czynności drugą byłoby obietnicą, że spakowana została biblioteka.

Wykaz zasobów nie jedzie tą drogą z tego samego powodu, dla którego nie jedzie
komendami obszaru multimediów: identyfikator pliku biblioteki wraca z magazynu
Designu odmową braku zasobu, a przycisk prowadzący do pewnej odmowy nie pełni
żadnej funkcji.

Postaci archiwum kontrakt nie zamyka wyliczeniem: pole jest napisem, a wartość
pusta oznacza postać domyślną rdzenia. Wykaz trzech dostępnych postaci pochodzi
z opisu komendy w kontrakcie, a nie z domysłu.

Archiwum o zerowej liczbie pozycji jest poprawną odpowiedzią, a nie awarią: tak
wraca spakowany katalog pusty. Zdanie odpowiedzi nazywa to wprost, ponieważ plik
powstał i jego pustka jest informacją należną operatorowi.

## budowa/klient-poprzedni/src/aod/odmowy-aod.ts

Zdanie odmowy nakładki składa się z trzech części. Nazwę nieudanej czynności
i powód podany przez rdzeń składa wspólny komponent odmowy. Trzecia część, czyli
rada dla Operatora, zależy od czynności nakładki i wspólny komponent nie ma jak
jej znać, dlatego wykaz rad stoi w warstwie nakładki.

Część odmów rdzenia w rodzinie komend nakładki jest zachowaniem poprawnym, a nie
usterką. Odpięcie procesu spoza wykazu przypiętych zwraca odmowę „nie znaleziono",
a puste wskazanie procesu przy przypięciu albo odpięciu zwraca odmowę
walidacyjną. Bez zdania trzeciego Operator odczytuje obie odpowiedzi jak awarię.

Nakładka nie stawia bramki przed wywołaniem: nie blokuje pustego pola, nie
wygasza przycisku i nie pyta o potwierdzenie. Rozstrzygnięcie należy do rdzenia,
a nakładka nazywa jego odpowiedź językiem zrozumiałym dla Operatora.

Wykaz rad obejmuje obszary czynności nakładki — przypięcie, odpięcie, rozmowę,
głos, podpowiedzi, stan i kontekst — oraz obszary kolejki decyzji: procesy,
wstrzymanie, konfigurację i przejęcie. Dla kodu odmowy bez osobnego zdania
podawana jest rada domyślna.

## budowa/klient-poprzedni/src/moduly/agents/zasieg-okien.ts

Najwęższym poziomem zasięgu uprawnień jest okno komunikacji i tryb uprawnień
żyje właśnie tam: pole `Window.permissionMode` niesie wartość przełącznika
`--permission-mode` kanału głównego, a zmienia je komenda `window.update`.

Ta część nie powiela ani okna izolacji, ani panelu sterowania: wykaz pokazuje
stan wszystkich okien sesji naraz, czego żadne z nich nie robi.

Lista wyboru pokazuje tryb okna, a nie wybór wskazany przez Operatora.
Przeglądarka przestawia listę natychmiast, a zapis idzie dopiero po niej, więc
po odmowie rdzenia lista stałaby na wartości, której okno nie ma. Dlatego każda
odmowa cofa listę do trybu ostatnio potwierdzonego przez rdzeń.

## budowa/klient-poprzedni/src/konfiguracja/lancuch-zasiegow.ts

Wiersz, z którego pochodzi wartość obowiązująca w punkcie widzenia okna, jest oznaczony
sygnałem. Pozostałe zapisy zostają widoczne, ponieważ Operator ma wiedzieć nie tylko to,
co obowiązuje, lecz także to, co czeka na innym poziomie.

Ostatni wiersz należy do wartości domyślnej katalogu. Brak zapisu nie jest dziurą: jest
wartością domyślną i tak został nazwany.

Każdy zapis wolno usunąć. Usunięcie nie jest niszczeniem ustawienia, tylko zdjęciem
nadpisania, po którym wartość wraca do poziomu szerszego.

## budowa/klient-poprzedni/src/moduly/multitasking/sterowanie-koordynatora.ts

Pas sterowania koordynatora zestawia siedem przycisków, z których pięć — Start, Stop, Pauza, Wznów i Powtórz — dzieli jeden silnik kolejek: idą akcją kolejki na kolejce etapu, tej samej, którą obsługuje pętla sesyjna i moduł Automations. Koordynator nie ma własnego wykonawcy zleceń i nie prowadzi biegu naprawczego: bieg prowadzi rdzeń, a licznik obiegów przychodzi z odczytu stanu okna.

Przycisk „Przekaż” wydaje zlecenie wykonawcy. Adresata i treść rozstrzyga tryb współpracy; nośnikiem jest wysłanie wiadomości, ponieważ okno wykonawcy jest oknem komunikacji. Droga modelu, czyli wywołanie narzędzia platformy, biegnie obok i tej kontrolki nie zastępuje.

Przycisk „Waliduj” sprawdza przebieg, a nie ocenia wynik: kontrakt nie ma komendy oceny, więc przycisk odczytuje stan koordynatora i wykonawców i mówi, czy tura trwa, ile obiegów zliczono i czy bieg stoi wraz z powodem. Ocena treści wyniku należy do modułu Results Analyzer.

Zatrzymanie biegu gasi kolejkę etapu i tury wykonawców naraz. Kolejność ma znaczenie: najpierw gaśnie tura, która właśnie zużywa czas modelu, a dopiero potem kolejka, która mogłaby ją wznowić.

Przekazanie zlecenia to dwie czynności i obie muszą się odbyć. Doręczenie treści oknu wykonawcy pozwala mu ruszyć do pracy. Utrwalenie w bazie rdzenia zapisuje, kto komu co zlecił, i zakłada pozycję kolejki; bez niego więź ginie z zamknięciem przeglądarki, a moduł Mission Control nie ma czego pokazać po ponownym uruchomieniu. Kolejność ma znaczenie: najpierw zapis, potem doręczenie — odwrotna zostawiałaby wykonawcę pracującego nad zleceniem, którego nikt nie odnotował. Nieudany zapis nie wstrzymuje doręczenia i nie liczy się jako niedoręczenie: dwa rodzaje niepowodzenia wracają osobnymi wykazami, bo pierwszy mówi, że wykonawca nie dostał pracy, a drugi, że dostał, ale po ponownym uruchomieniu nikt tego nie odtworzy.

## budowa/klient-poprzedni/src/moduly/apps/karta-rozszerzenia.ts

Pola, których pozycja katalogu nie niesie — dziennik zmian, wykaz narzędzi oraz
wymagane uprawnienia — nie mają na karcie zaślepki. Stoją w wykazie braków okna,
żeby nieobecność danych była widoczna, a nie udawana.

Stan pozycji nigdy nie opiera się na samej barwie: plakietka niesie słowo, a nie
tylko odmianę stylu. Cztery rozróżnialne stany (dostępna, zainstalowana
wyłączona, zainstalowana włączona, dostępna aktualizacja) składa się z dwóch pól
kontraktu, więc czwartego z nich karta nie podaje — rejestr nie niesie wersji
dostępnej do zestawienia z wersją zainstalowaną.

Instalacja i przełączenie mają w rdzeniu osobne komendy, dlatego karta stawia
dwa osobne przyciski. Jeden przycisk o zmiennym napisie kazałby zgadywać, którą
czynność zaraz zleca. Przycisk uprawnień otwiera Permissions & Trust Center na
pozycji karty i pozostaje czynny również dla pozycji niezainstalowanej, ponieważ
zakres dostępu ogląda się przed włączeniem.

## budowa/klient-poprzedni/src/aplikacja/przelacznik-tras.ts

Żaden przycisk trasy nie jest bramą i żaden nie zostaje wyszarzony. Do
środowiska można wejść także wprost, ponieważ środowisko domyślne istnieje od
pierwszej chwili. Trasa bieżąca jest oznaczona atrybutem `aria-current`, a nie
odebraniem klikalności.

Przyciski niosą sam znak, bez napisu, ponieważ trasy stoją w grupie akcji paska
obok ustawień, motywu i menu konta, które również są ikonami, a nazwa widoku
stoi tuż pod paskiem jako tytuł strony. Nazwę niesie atrybut `title` oraz
etykieta dostępności przycisku.

## budowa/klient-poprzedni/src/moduly/design/historia-promptow.ts

Odczyt historii i zapis szablonu mają już komendy w kontrakcie, brakuje im
natomiast uchwytów w rdzeniu. Okno tej drogi jeszcze nie wywołuje, więc
kontrolka zapisu szablonu nazywa stan, zamiast wykonywać zapis.

Porównanie wymienia pola kontraktu, którymi dwa prompty się różnią. Prompt jest
zbiorem pól, więc różnica liczona na samym tekście nie niosłaby informacji
o tym, co się zmieniło.

## budowa/klient-poprzedni/src/mission-control/przyciski-transportu.ts

Żaden przycisk paska nie dostaje `disabled`, `aria-disabled` ani klasy wygaszającej:
blokada nie jest dozwolonym sposobem informowania o stanie. Stan kolejki jest wypisany
słowem przy jej nazwie, a przycisk zatrzymania pozostaje czynny w każdym stanie.

Osiem działań wyliczenia `QueueAction` leżących poniżej sterowania biegiem — wstawienie
i zdjęcie z kolejki, odłożenie w czasie, rozdzielenie, scalenie, skierowanie,
rozgałęzienie i warunek — należy do układania przebiegu, a nie do jego prowadzenia.
Pasek transportu pulpitu ich nie pokazuje, ale mapa ikon musi je znać, ponieważ
wyliczenie kontraktu je niesie, a mapa niepełna nie skompilowałaby się przy pierwszym
ich użyciu.

## budowa/klient-poprzedni/src/dostepy/wybor-korzeni.ts

Reguła pustego zbioru stoi przy kontrolce zdaniem, a nie w podpowiedzi, ponieważ
Operator, który odznaczy wszystko, ma wiedzieć, że właśnie nadał wszystko.

Nadanie nie sięga poza korzenie punktu, więc kontrolka nie przyjmuje ścieżki wpisanej
z ręki: pokazuje korzenie punktu i pozwala je zaznaczyć. Poszerzenie obszaru wymaga
zmiany punktu, nie nadania.
