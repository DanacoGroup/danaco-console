# Uzasadnienia rdzenia

Dokument gromadzi uzasadnienia decyzji projektowych, które nie mieszczą się
w nagłówku pliku źródłowego. Rozdział nosi nazwę pliku, którego dotyczy.

## budowa/server/internal/core/skutek_warsztatu_developera_test.go

Sprawdziany skutku warsztatu modułu Developer nie sprawdzają, czy odpowiedź
jest odpowiedzią: koperta ze stanem `ok` i pustym wynikiem jest kopertą udaną
i zarazem kłamiącą. Dlatego każdy sprawdzian tego pliku po udanej odpowiedzi
mierzy niezależnie — własnym zapytaniem SQL do tej samej bazy albo odczytem
pliku z dysku. Zapytanie idzie osobnym połączeniem do pliku bazy, a nie przez
repozytorium rdzenia, ponieważ droga przez repozytorium mierzyłaby to samo,
czym mierzy się rdzeń, i wspólna usterka odczytu zostałaby niewidoczna po obu
stronach. Uprząż tego pliku jest własna, a nie wspólna z uprzężą montującą
cały rdzeń kopertami protokołu, ponieważ moduł Developer potrzebuje okna
z katalogiem roboczym, którego kontrakt komendy nie zakłada; adapter składany
wprost daje dokładnie tę jedną rzecz, a mierzone i tak jest to, co zostało
w bazie i na dysku.

Sprawdzian wyniku testów i pokrycia wchodzi drogą, którą wchodzi rdzeń:
zakłada przebieg, dopisuje mu wiersze wyjścia tak, jak robi to pompa logu,
i domyka pomiar. Mierzy potem bazę własnym zapytaniem, bo to ona jest jedynym
miejscem, z którego wynik testu da się odczytać po zamknięciu okna.

Sprawdzian wykazu kontenerów nie zakłada, czy silnik kontenerów na maszynie
stoi — sprawdza spójność odpowiedzi z tym, co zastała: wykaz niepusty przy
`engineAvailable: false` byłby odpowiedzią wewnętrznie sprzeczną.

Sprawdzian formatowania pomija się, gdy serwer nie ma formatera plików Go:
brak programu jest wtedy stanem serwera, a nie usterką modułu, i mówi o tym
sonda warsztatu, która ma własny sprawdzian.

Sprawdzian zapytania API obejmuje też podstawienie `{{nazwa}}` ze środowiska
kolekcji — bez niego kolekcja miałaby adres wpisany na stałe.

Sprawdzian sesji debugowania mierzy skutek, a nie kopertę: odpowiedź `ok`
z pustym stosem wywołań byłaby dokładnie tym wzorcem szkody, którego ten plik
pilnuje. Dlatego pomiar idzie po wartości zmiennej odczytanej z zatrzymanego
procesu — takiej, której nie da się oddać bez faktycznego zatrzymania
programu.

Sprawdzian analizy statycznej poza Go nie liczy zgłoszeń — liczy, który
program je wystawił. Wykaz niepusty złożony wyłącznie ze zgłoszeń Go byłby
odpowiedzią, która wygląda dobrze i milczy o plikach, których nie sprawdzono.

Sprawdzian odmowy przy braku programu składni podmienia program na nazwę,
której na żadnej maszynie nie ma: brak jest stanem maszyny, więc sprawdzian
ten stan odtwarza, zamiast czekać, aż zastanie go u kogoś.

Sprawdzian nawigacji po TypeScripcie zakłada, że pusty wykaz przy
`serverAvailable: true` znaczyłby „sprawdziłem i nie ma” — zdanie
nieprawdziwe, którego okno nie miałoby jak odróżnić od prawdziwego.

Sprawdzian przebiegu obciążeniowego z odpowiedziami mierzy skutek
niezależnie od odpowiedzi rdzenia: wynik mówiący o tysiącach żądań przy
liczniku serwera sprawdzianu na zerze byłby przebiegiem zmyślonym.

Sprawdzian przebiegu obciążeniowego punktu milczącego zamyka serwer przed
uruchomieniem przebiegu, więc odmowa ma nazwać liczbę błędów, a nie oddać
wynik złożony z samych zer, który czytałby się jak usługa skrajnie wolna
zamiast jak usługa nieobecna.

## budowa/server/internal/core/skutek_budowy_produktu_test.go

Sprawdziany modułu Apps sprawdzają, czy za odpowiedzią rdzenia stoi zapis,
plik albo stojący serwer, a nie samą kopertę: moduł Design w tym produkcie
już meldował `status: ok` z wykazem zasobów, za którymi nie było ani jednego
bajtu. Dlatego każdy sprawdzian tego pliku, po udanej odpowiedzi, schodzi
niżej i mierzy niezależnie — do bazy drugim połączeniem otwartym niezależnie
od rdzenia, do pliku w magazynie treści sprawdzając, czy to naprawdę jest
archiwum, obraz albo dokument, za który się podaje, albo do sieci, wołając
adres, który rdzeń wypuścił jako `previewUrl`. Sprawdziany nie pomijają się
przy braku żadnego programu, bo cały moduł pracuje bibliotekami
wkompilowanymi w rdzeń — `archive/zip`, `image/png`, `crypto/ed25519`,
`net/http` — więc mierzą to, co ma stać u Operatora na cienkiej instalce.

Funkcja pomocnicza `oknoModuluSprawdzianu` zakłada prawdziwe okno tam, gdzie
`apps.deployment.run` bierze treść z przestrzeni roboczej okna i odmawia
oknu, którego rdzeń nie zna; pozostałe komendy obszaru okna w rejestrze nie
wymagają, więc reszta sprawdzianów posługuje się nazwą stałą i mierzy to
samo taniej.

Sprawdzian audytu wydajności Core Web Vitals jest wyjątkiem od zasady tego
pliku: pomija się bez programu pomiarowego, bo audyt wydajności z założenia
idzie programem, a nie biblioteką wkompilowaną w rdzeń; odmowę przy jego
braku mierzy sprawdzian osobny.

## budowa/server/internal/core/skutek_zarzadu_biblioteki_test.go

Sprawdziany rodziny `library.*` nie kończą się na odpowiedzi komendy: moduł
Design w tym produkcie już meldował `status: ok` wraz z wykazem zasobów, za
którymi nie było ani jednego bajtu. Każdy sprawdzian tego pliku schodzi niżej
niż rdzeń — otwiera plik bazy osobnym połączeniem SQL albo czyta bajty
z magazynu treści na dysku — i pyta o to samo, co komenda zameldowała. Miara
jest niezależna, nie druga komenda: gdyby stan czytała inna komenda tego
samego modułu, obie mogłyby mylić się zgodnie, bo adapter oddający wykaz
z tego samego miejsca, w którym zapisał, potwierdziłby sam siebie.

Sprawdzian opisu zasobu mierzy skutek w tabeli opisu, a nie samą odpowiedź
komendy, ponieważ odpowiedź niesie opis po zapisie i sama z siebie zawsze
wygląda pomyślnie.

Sprawdzian weryfikacji integralności mierzy drogę, której żadne okno nigdy
nie przejdzie normalnie: bajtów zasobu nie oddaje żadna komenda kontraktu,
więc treść na dysku podmienia sam sprawdzian.

Pulpit stanu na pustym repozytorium sprawdza przypadek graniczny agregatów:
sumowanie po zbiorze pustym daje w SQL wartość pustą, a pulpit ma pola
liczbowe bez stanu pustego — zbiór pusty ma dać zera, nie odmowę, bo zero
zasobów jest wynikiem, a nie niepowodzeniem odczytu.

Pliki `jpegZeZnacznikiemXmp` i `mp3ZeZnacznikiemId3` składają materiał ręcznie,
bo znacznik XMP i ID3v2 musi być osadzony w bajtach samego pliku, a materiał
wniesiony spoza pliku niczego by nie dowiódł; koder `image/jpeg` nie umie
wstawić własnego segmentu, a biblioteka standardowa nie zapisuje MP3 wcale.

Sprawdzian odczytu metadanych osadzonych przez XMP pomija się z nazwanym
powodem na maszynie bez programu: odczyt XMP, IPTC i ID3 jest poszerzeniem
opisu, a nie jego warunkiem — resztę pól pilnuje sprawdzian odczytu metadanych
technicznych, który idzie zawsze.

Sprawdzian braku programu zeruje ścieżkę wyszukiwania, żeby uczynić z bieżącej
maszyny maszynę bez programu, i mierzy to zdanie wszędzie, nie tylko na
cienkiej instalce.

Zapora pokrycia rodziny `library.*` stoi obok sprawdzianu pokrycia całego
kontraktu, bo ten moduł ma osobno pilnować, żeby żadna z jego 48 komend nie
została bez uchwytu w rejestrze rdzenia.

## budowa/server/internal/core/skutek_designu_test.go

Moduł Design ma w tym produkcie precedens szkody: `design.asset.list` meldował
kiedyś `status: ok` z wykazem zasobów, za którymi nie było ani jednego bajtu.
Koperta była udana i zarazem kłamała, a klient czyta kopertę, nie komentarz
w kodzie. Od tamtej pory każdy sprawdzian tego pliku schodzi po odwołaniu do
magazynu albo do bazy i mierzy niezależnie, nigdy nie kończąc się na treści
odpowiedzi. Stąd dwa narzędzia pomiaru, oba omijające rdzeń: wydanie zasobu
jest dekodowane z powrotem jako obraz i mierzone co do wymiarów po skali, bo
odpowiedź mówiąca „oto png @2x” nie jest dowodem — dowodem jest `image.Decode`,
któremu te bajty wystarczą, i szerokość, która wyszła dwa razy większa; a
kolekcja, wersja kompozycji, adnotacja i zestaw żetonów są odczytywane drugim,
niezależnym połączeniem SQLite do pliku bazy stanowiska, bo rdzeń nie bierze
w tym odczycie udziału — gdyby zapis nie doszedł do pliku, odpowiedź komendy
i tak wyglądałaby tak samo. Sumę kontrolną treści sprawdzian liczy z bajtów,
nie przepisuje jej z odpowiedzi, więc zgodność tych dwóch wartości jest tu
mierzona, a nie założona.

Funkcja pomocnicza `zycieZGniazdemDesignu` podstawia tożsamość połączenia,
bo uprząż woła rdzeń z pominięciem gniazda, a zgłoszenie obecności bez
rozpoznanego klienta jest odmawiane — inaczej sprawdzian mierzyłby wywołanie,
które w produkcie nie zachodzi.

## budowa/server/internal/core/skutek_katalogu_rozszerzen_test.go

Sprawdziany rodziny `extension.*` nie kończą się na tym, że odpowiedź jest
udana — ten sam wzorzec szkody, którego pilnuje sprawdzian modułu Developer,
grozi kopertą `ok` bez pokrycia. Każdy sprawdzian tego pliku schodzi niżej:
do bazy drugim połączeniem i liczy wiersze, do pliku w magazynie treści
i czyta bajty, albo do serwera protokołu podniesionego przez sprawdzian, który
wie, o co go naprawdę zapytano. Serwer MCP sprawdzianu jest prawdziwym
serwerem JSON-RPC nad HTTP: odpowiada na `initialize`, `tools/list`
i `tools/call` i zapamiętuje, co dostał. Mierzona jest droga rdzenia —
powitanie, odkrycie, wywołanie, dziennik ramek, metryka użycia — a nie to,
co odpowiada konkretny serwer.

## budowa/server/internal/core/skutek_debaty_test.go

Wzorzec szkody, którego pilnuje ten plik, wystąpił w tym produkcie: komenda
meldowała `status: ok` z wykazem, za którym nie stał ani jeden bajt. Dlatego
żaden sprawdzian tutaj nie kończy się na tym, że odpowiedź jest udana — każdy
schodzi do bazy drugim połączeniem, otwartym niezależnie od rdzenia, i liczy
wiersze albo czyta bajty z magazynu. Kanałem uczestników jest kanał echo:
odsyła treść zapytania porcjami, tak jak zrobiłby to model, bez sieci, konta
i klucza. Mierzona jest droga rdzenia — rejestr, tura, zapis, analiza,
głosowanie, wydanie — a nie to, co odpowiada konkretny dostawca.

## budowa/server/internal/core/adapter_studio_bezpieczenstwo.go

Rodzina bezpieczeństwa dokumentu modułu Studio — szyfrowanie, czyszczenie
metadanych, redakcja poufności, podpis i jego weryfikacja oraz rozpoznanie
danych wrażliwych — pracuje wyłącznie na bibliotece `pdfcpu` i na bibliotece
standardowej Go (`crypto/x509`, `crypto/rsa`, `crypto/sha256`), bez ani
jednego uruchomienia procesu zewnętrznego: funkcja zależna od programu,
którego instalka nie niesie, jest u Operatora odmową, a nie funkcją, a rodzina
bezpieczeństwa jest tym miejscem, w którym odmowa boli najbardziej, bo Operator
dowiaduje się o niej dopiero po wysłaniu dokumentu, którego nie oczyścił.

Zamalowanie prostokąta zostawia tekst pod spodem: da się go zaznaczyć,
skopiować i odczytać wyszukiwarką — to jest redakcja pozorna. `Zredaguj`
wycina z treści strony bloki tekstu leżące w obszarze, dopiero potem kładzie
prostokąt, więc wynik nie zawiera już wyciętych znaków.

Podpis dokumentu nazywa się wprost podpisem tego produktu, nie podpisem
kwalifikowanym: podpis, o którym Operator sądziłby, że niesie skutek prawny
podpisu kwalifikowanego, byłby gorszy niż brak podpisu.

## skutek_przegladarki_test.go

Sprawdziany tego pliku nie kończą się na odpowiedzi komendy modułu Browser
oznaczonej jako udana. Każdy schodzi o poziom niżej niż ta odpowiedź: albo do
bazy — drugim, niezależnym połączeniem do tego samego pliku, którym jedzie
rdzeń — albo do magazynu treści, gdzie liczy bajty leżące pod odwołaniem.
Powód jest znany z doświadczenia tego produktu: odpowiedź udana z pustym
wykazem wygląda dokładnie tak samo jak odpowiedź udana z wykazem prawdziwym,
a klient czyta samą odpowiedź.

Sprawdziany zależne od zrzutu ekranu albo silnika przeglądarki wymagają
programu Chromium wpisanego do sondy zależności zewnętrznych. Na maszynie bez
tego programu pomijają się, nazywając powód pominięcia; pozostałe sprawdziany
biegną zawsze, ponieważ stoją wyłącznie na bibliotece standardowej.

Test monitora sprawdza różnicę treści miarą liczby wierszy dodanych
i usuniętych, nie samym znacznikiem zmiany wziętym z odpowiedzi komendy bez
pokrycia w treści.

Test audytu dostępności na stronie z naruszeniem obrazka bez tekstu
zastępczego sprawdza, że zgłoszenie niesie selektor węzła — audyt bez
wskazania węzła nie mówi, co poprawić. Test audytu strony zgaszonej sprawdza,
że zero naruszeń na stronie, która przestała odpowiadać po przejściu, wraca
odmową, nie wynikiem zerowym: zero naruszeń na nieistniejącej stronie byłoby
brakiem pomiaru podanym jako pomiar.

Test odmowy bez programu audytującego zwęża ścieżkę wyszukiwania do katalogu
z samą przeglądarką, aby odróżnić brak programu audytującego od braku
przeglądarki — odmowa o przeglądarce mówiłaby o innym braku niż ten, który
sprawdzian bada.

## budowa/server/internal/core/studio_dlugi_domkniecie_test.go

Ten plik mierzy skutek, nie kopertę, dla sześciu długów funkcjonalnych modułu
Studio: paginacji spisu treści liczonej kartką A4 zaszytą na stałe, przez co
numer strony zgadzał się z podglądem, a nie z nastawami dokumentu; cofnięcia
czynności, które odtwarzało treść i postać, a aparat dokumentu zostawiało
w brzmieniu bieżącym; malarza formatów żyjącego w pamięci procesu, przez co
postać zabrana przepadała przy przeładowaniu rdzenia; zmiany postaci zapisanej
bez kodu wykonawcy, po której rozbicie zmian po wykonawcy pokazywało dwóch
agentów jako jednego nienazwanego; sortowania tabeli układającego
„Łukasiewicza” za „Zawadzkim”, bo porównanie szło po Unikodzie, a nie
porządkiem alfabetycznym pisma polskiego; oraz wklejenia „zachowaj postać
źródła”, które zachowywało postać miejsca wklejenia zamiast postaci źródła.

Cofnięcie czynności aparatu jest mierzone po dwóch stronach naraz: z drzewa
postaci i z wiersza aparatu czytanego przez `studio.apparatus.list` — usunięcie
elementu tylko z jednej z tych dwóch stron zostawiłoby czytelnika z przypisem,
którego dokument już nie niesie, albo z wierszem bez pokrycia w drzewie.

Wklejenie sposobem „zachowaj postać źródła” różni się od sposobu „scal
postać”: gdy wpis schowka nie przenosi postaci źródła, oba sposoby dają ten
sam wynik, więc sam fakt powodzenia wklejenia niczego nie dowodzi — dowodem
jest wytłuszczenie fragmentu źródłowego widoczne w miejscu wklejenia.

## budowa/server/internal/core/skutek_warsztatow_designu_test.go

Sprawdziany warsztatów wektora, kroju ikonowego, schematu, makiety i barwy
modułu Design mierzą liczbę wynikową operacji zamiast samego faktu odpowiedzi
bez błędu: operacja logiczna jest mierzona współrzędnymi węzłów i przynależnością
punktów płaszczyzny do kształtu, krój ikonowy liczbą glifów odczytaną z pliku
TTF czytnikiem krojów wraz z liczbą krzywych w konturach, schemat liczbą ścieżek
w dokumencie SVG, barwa współczynnikiem kontrastu przeliczonym niezależnie
wzorem WCAG oraz kątami odcieni harmonii, układ makiety położeniami warstw
odczytanymi z bazy po zapisie, a odszumienie tym, czy krawędź pozostaje krawędzią,
a płaski obszar wygładzeniem.

Pierwsze uruchomienie tego zestawu wykryło defekt tabeli `maxp` składanego
kroju: tabela miała 36 bajtów wobec 32 wymaganych przez format, więc żaden
czytnik krojów nie wczytywał pliku, choć rdzeń oddawał go bez odmowy. Sprawdzian
oparty wyłącznie na fakcie udanej odpowiedzi komendy nie wykrywa takich usterek,
ponieważ operacja oddająca pierwszy kształt zamiast sumy, krój bez glifów albo
schemat wychodzący pustym płótnem są dla niego odpowiedziami udanymi.

## budowa/server/internal/core/skutek_zakresu_eksperta_test.go

Sprawdzian koperty odpowiedzi komendy potwierdza wyłącznie to, że komenda
zwróciła status powodzenia; nie potwierdza, że zawężenie zakresu eksperta
faktycznie działa. Permissions Center może oddawać `status: ok` przy każdym
zawężeniu i wyglądać na działający, podczas gdy ekspert nadal wykonuje
wszystko, co wykonywał wcześniej. Zawężenie, którego nikt nie egzekwuje, jest
gorsze niż jego brak, ponieważ operator widzi ograniczenie w interfejsie
i pracuje w przekonaniu, że ono obowiązuje.

Z tego powodu każdy sprawdzian w tym pliku sięga do bazy własnym zapytaniem
SQL albo mierzy odmowę wydaną przez straż eksperta, a nie sam zapis komendy.

## budowa/server/internal/core/adapter_narzedzia_obraz_czynnosci.go

Priorytet drogi wkompilowanej nad programem pakietu serwera wynika z jednego
rozpoznawalnego błędu: `errBrakRachunkuGoObrazu`. Tylko on przełącza czynność
na program; każdy inny błąd jest odmową wprost, ponieważ obraz uszkodzony ma
zostać nazwany, a nie oddany drugiej drodze, która powie o nim to samo wolniej.
Program wchodzi wyłącznie tam, gdzie rachunku Go nie ma wcale — AVIF bez kodera
i dekodera w Go oraz WEBP stratny, którego `nativewebp` nie zapisuje. Składanie
argumentów programu w `argumentyPrzeksztalcenia`, `argumentyPoprawki`
i `argumentyKonwersji` musi umieć dokładnie to samo, co rachunek wkompilowany,
inaczej ten sam wniosek modelu dawałby dwa różne skutki zależnie od formatu
pliku.

Dogniecenie zapisu programem dogniatającym należy wyłącznie do `image.convert`,
ponieważ tylko ta czynność obiecuje w odpowiedzi `savedBytes`; przekształcenie
i poprawka zmieniają treść obrazu, więc pytanie o oszczędność miejsca nie ma
przy nich sensu. Krok dogniatania nie odmawia: przy braku programu oddaje bajty
bez zmiany, więc zapis udaje się także na maszynie, która go nie niesie.

Każda wartość `ImageTransformKind` i `ImageAdjustKind` ma własną gałąź
w `argumentyPrzeksztalcenia` i `argumentyPoprawki`; wartość spoza wyliczenia
kończy się odmową nazywającą ją wprost, bez gałęzi domyślnej, która oddałaby
zasób identyczny ze źródłem jako rzekomy skutek retuszu.

Brak pola `amount` w żądaniu poprawki bierze wartość domyślną operacji, a nie
zero: zero byłoby poprawką bez skutku, a model proszący o rozjaśnienie bez
podanej liczby dostałby obraz nieodróżnialny od źródła.

`formatyDocelowe` ogranicza `image.convert` do zamkniętego zbioru formatów,
ponieważ nazwa formatu trafia w argument programu drogi zapasowej:
przepuszczenie dowolnego tekstu, na przykład `ephemeral:`, dałoby modelowi
wpływ szerszy niż zamiana formatu.

Format wyniku w `przetworz` równa się formatowi źródła, gdy czynność go nie
zmienia — bez tego wymuszenia zapis nie miałby jak nazwać formatu; format
nierozpoznany oddaje `png` jako wybór bezstratny.

Odmowa rachunku wkompilowanego w `policzWkompilowanym` — kadr poza obrazem,
operacja spoza wyliczenia — wraca wprost i kończy czynność: droga zapasowa nie
ma jej czym naprawić, a jej uruchomienie zamieniłoby odmowę czytelną na odmowę
programu.

## budowa/server/internal/core/postac_skutek_stylu_test.go

Ten plik mierzy skutek trzech wymagań odbioru, których nie mierzył żaden inny
plik: zmiana stylu nazwanego przestawia wszystkie miejsca, które go używają;
zastosowanie postaci do fragmentu zmienia wyłącznie ten fragment, bez
rozlania się na akapit sąsiedni; a czynność wywołana jako narzędzie modelu
odkłada zmianę śledzoną autora `model`, a nie autora `uzytkownik`, bo bez
tego przełącznik pokazujący wszystko, co zrobił model, nie ma czego
podświetlić.

Miara nie idzie z odpowiedzi czynności, bo odpowiedź pisze ten sam kod, który
zmieniał dokument, i potwierdzałaby samą siebie. Idzie dwiema drogami
niezależnymi: drugim połączeniem do pliku bazy własnym zapytaniem SQL,
z pytaniem „czy to naprawdę leży w bazie”, oraz osobnym wywołaniem komendy
odczytu, z pytaniem „czy Operator, który otworzy dokument ponownie, zobaczy
to samo”. Wszystkie czynności idą przez rejestr, a nie po adapterze, bo
narzędzie modelu jedzie dokładnie tą drogą — pomiar autora `model` mierzy
wtedy drogę, którą naprawdę pojedzie model.

## budowa/server/internal/core/skutek_zasobu_designu_test.go

Wzorzec szkody, którego pilnuje ten plik, wydarzył się już w tym produkcie:
`design.asset.generate` meldował `status: ok` z wykazem zasobów rodzaju
`image`, za którymi nie było ani jednego bajtu. Koperta była poprawna, kontrakt
spełniony, a panel zasobów dostawał kafelki bez treści. Z tego powodu żaden
sprawdzian tego pliku nie kończy się na sprawdzeniu, że odpowiedź jest udana:
każdy schodzi po odwołaniu `asset.uri` do magazynu i pyta plik, ile ma bajtów
i czy są tymi bajtami, które wjechały.

Kanał obrazowy jest w tych sprawdzianach prawdziwym kanałem rdzenia, nie
zaślepką: wiersz rejestru rodzaju `api` z parametrem `adapter: obrazy`,
wskazujący adres serwera podniesionego na czas sprawdzianu w tym samym
procesie. Droga mierzona jest dzięki temu tą samą drogą, którą przechodzi
wywołanie operatora — żądanie HTTP, odpowiedź kształtu OpenAI Images, fragment
`image`, magazyn, wiersz — bez żadnego programu spoza maszyny.

Rodzaj zasobu spoza siedmiu wartości kontraktu odbijał się dotąd od warunku
CHECK schematu bazy i wracał kodem `internal_error` z `retryable: true`,
cytując w treści warunek bazy wraz z nazwą kolumny. Żądanie takie nie mogło
się udać przy żadnym ponowieniu, więc klient z pętlą ponowień powtarzał je bez
końca, a operator dostawał zdanie o kolumnie zamiast o swoim żądaniu.

## adapter_studio_pdf.go

Cała rodzina czynności warsztatu PDF pracuje na bibliotece pdfcpu
wkompilowanej w binarium rdzenia. Nie ma tu ani jednego uruchomienia procesu
zewnętrznego i mieć nie będzie — to jest zasada bezwzględna produktu, nie
wybór wygody. Funkcja zależna od programu, którego instalka nie niesie, jest
u odbiorcy odmową, a nie funkcją: maszyna deweloperska ma doinstalowane
wszystko, więc sprawdzian na niej świeci zielono przy czynności, która u
odbiorcy nie ruszy ani razu. Komenda obsługiwana, która pod spodem woła cudzy
program, jest fasadą — kontrakt niesie polecenie, adapter je przyjmuje, a
odbiorca dostaje puste okno.

Żadna z czynności nie zmienia materiału w miejscu. Wynik jest nowym zasobem
pod własną sumą kontrolną, a materiał zostaje nietknięty — dokument zmieniony
w miejscu byłby dokumentem, do którego nie ma jak wrócić. Praca idzie na
bajtach w pamięci, nie na plikach pośrednich: magazyn oddaje treść, biblioteka
przetwarza strumień, wynik wraca do magazynu; plik pośredni byłby trzecim
miejscem, w którym ta sama treść żyje.

Sprawdzanie zgodności biblioteki jest wyłączone, ponieważ dokumenty zastane
bywają niezgodne ze specyfikacją w szczegółach, których nikt nie zmieni,
a odmowa pracy nad takim dokumentem byłaby odmową pracy nad materiałem, który
otwiera każda przeglądarka.

Odchudzanie dokumentu nie przelicza obrazów w dół, więc zysk bywa mniejszy niż
przy narzędziu rasteryzującym. Jest to zamiana świadoma: przeliczenie obrazów
wymagałoby silnika rasteryzacji, czyli programu spoza instalki, a wtedy
czynność przestałaby działać u odbiorcy.

Czynność UlozStrony odmawia wstawienia stron z innego dokumentu, ponieważ
biblioteka potrafi wstawić wyłącznie strony puste — wstawienie cudzej treści
jest scalaniem z wyborem miejsca i ma iść komendą scalania, nie wstawiane
jako puste kartki z meldunkiem powodzenia.

Czynność stemplowania kładzie pieczęć nad treścią strony, nie pod nią —
pieczęć schowana pod treścią byłaby pieczęcią niewidoczną mimo meldunku
o powodzeniu.

Czynność wyciągania obrazów i załączników oddaje liczbę wyciągniętą, nie
zapowiedzianą: dokument bez obrazów oddaje zero i to jest odpowiedź
prawdziwa, nie niepowodzenie.

## budowa/server/internal/core/skutek_wejscia_test.go

Sprawdziany drogi wejścia do aplikacji mierzą skutek trwały w bazie i w skrzynce,
nigdy samą odpowiedź udaną komendy — odpowiedź udana nie jest dowodem niczego,
bo pole stanu potrafi nieść wartość pozytywną przy zapisie, który się nie udał.

Droga rejestracji bez konta nadawczego, czyli bez ustawionego nadajnika poczty,
przechodzi zgodnie z rejestrem decyzji: poczta jest potrzebna do pisania do
innych ludzi, a nie do postawienia bramki na własnym urządzeniu, a nadajnik
ustawia się w oknie Konfiguracji, czyli już za tą bramką. Czwórka sprawdzianów
tej drogi mierzy kolejno: że rejestracja bez poczty zakłada konto i stawia
w sejfie znacznik pamiętający brak adresu nadawczego; że samo hasło otwiera
bramkę mimo braku potwierdzenia adresu, bo znacznik zdejmuje wyłącznie ten jeden
warunek wejścia; że ustawienie nadajnika po rejestracji nie potwierdza adresu
samo z siebie, skoro potwierdza go wyłącznie droga przepisana z listu, a listu
na tej instalce nie było; oraz że znacznik istnieje tylko po tej stronie
granicy, gdzie listu nie było — na drodze z pocztą działającą nie ma prawa
powstać ani przed potwierdzeniem adresu, ani po nim, bo inaczej wyjątek
pierwszego uruchomienia zamieniłby się w trwałe obejście bramki.

Sprawdzian jednorazowości rejestracji mierzy też, czego druga próba nie robi:
nie podmienia konta, nie zakłada drugiej kotwicy hasła i nie wysyła drugiego
listu; kod odmowy `conflict` pozwala klientowi odróżnić konto już istniejące od
awarii wartej ponowienia.

Sprawdzian jednorazowości drogi potwierdzenia mierzy obie jej połowy: droga
z listu wpuszcza raz i przy drugiej próbie odmawia, a druga połowa mierzona jest
na drodze odzyskania konta, ponieważ powtórzone `auth.verify` odbija się o stan
konta już potwierdzonego zanim dojdzie do sprawdzenia samej drogi, więc nie
dowodzi, że zamknięcie wiersza drogi w ogóle działa — jedynym zabezpieczeniem
drugiego użycia drogi odzyskania jest zamknięcie jej wiersza w bazie.

Pole `cel` rozdziela drogę weryfikacji od drogi odzyskania, choć obie wyglądają
tak samo i leżą w jednej tabeli: bez tego rozdziału droga wysłana na prośbę
o nowe hasło potwierdzałaby adres, a droga wysłana przy rejestracji ustawiałaby
hasło. Sprawdzian dowodzi też, że próba użycia drogi poza jej czynnością nie
zużywa jej — odrzucenie i spalenie naraz byłoby gorsze od samej odmowy.

Odpowiedź `auth.recover` nie może być wyrocznią: komenda jest osiągalna przed
zalogowaniem, więc pyta ją każdy, a odpowiedź różniąca się choćby jednym bajtem
zdradzałaby pytającemu, czy dany adres ma konto właściciela. Dla adresu obcego
list nie wychodzi wcale, bo wysłany szedłby do osoby, która o nic nie prosiła.

Wykaz urządzeń oznacza urządzenie bieżące wyłącznie po tożsamości połączenia
przekazanej w kontekście gniazda, nigdy przez zgadywanie po ostatnim wejściu —
inaczej wskazywałby cudzą maszynę jako własną, a Operator odebrałby dostęp nie
tej sesji, co trzeba.

Sprawdzian trwałości drogi potwierdzenia nie ufa nazwom kolumn: czyta cały
wiersz tabeli i szuka w każdej wartości materiału wysłanego listem, żeby kolumna
dołożona kiedyś obok `skrot` nie przeszła pomiaru pytającego wyłącznie o `skrot`.
Gdyby materiał z listu leżał w tabeli jawnie, kopia bazy wystarczyłaby do
potwierdzenia cudzej tożsamości i ustawienia hasła do konta.

## skutek_warsztatu_tlumaczenia_test.go

Wzorzec szkody, którego pilnują sprawdziany tego pliku, ma w tym produkcie
precedens: komenda meldowała odpowiedź udaną z pustym wynikiem, a za
odpowiedzią nie leżało nic. Dlatego żaden sprawdzian tutaj nie kończy się na
udanej odpowiedzi. Każdy schodzi niżej, do jednego z dwóch miejsc, w których
skutek albo jest, albo go nie ma: do bazy — drugim, niezależnym połączeniem do
tego samego pliku SQLite, zapytaniem SQL wprost, z pominięciem całej warstwy
adapterów — albo na dysk — otwarciem pliku, który komenda miała wytworzyć,
i odczytaniem jego treści, nie samego istnienia. Żaden sprawdzian nie woła
modelu ani programu zewnętrznego, więc wszystkie wypadają tak samo u odbiorcy,
jak na maszynie budującej.

## budowa/server/internal/core/adapter_narzedzia_dokument.go

Same czynności leżą osobno wedle odpowiedzialności: zamiana formatu w
`adapter_narzedzia_dokument_konwersja.go`, odczyt treści w
`adapter_narzedzia_dokument_tekst.go`, wpięcie do rejestru w
`handlers_narzedzia_dokument.go`. Obie komendy rodziny są narzędziami modelu,
nie panelem operatora: `document.convert` oddaje operatorowi dokument, a nie
tekst w oknie rozmowy, a `document.text.extract` czyta plik, który operator
dostarczył — PDF, skan, zdjęcie kartki — bo bez tej komendy model nie ma
żadnej drogi do treści pliku leżącego na dysku operatora.

Jedna droga prowadzi do binarium: wszystkie uruchomienia idą przez
`zewnetrzne.Wolaj` — port `session.Uruchamiacz`, brama izolacji okna, objęcie
drzewa potomstwa, obowiązkowa granica czasu. Własnego `exec.Command` ten
moduł nie ma; w całym drzewie stoi dokładnie jedno wywołanie, w
`injection/rozruch.go`.

Katalog uruchomienia zostaje pusty, a ścieżki argumentów są bezwzględne.
Gdyby moduł podał bramie własny katalog tymczasowy, punkt izolacji katalogu
roboczego odrzuciłby uruchomienie jako wyjście poza katalog okna
(`session/izolacja_polecenie.go`). Pusty katalog pozwala bramie wstawić
katalog własnego zasięgu, a bezwzględne ścieżki sprawiają, że wybór katalogu
nie zmienia wyniku pracy.

Zasięg jest zasięgiem platformy. Żądania obu komend nie niosą okna, tak samo
jak w rodzinie `speech.*` (`adapter_modul_mowa.go`), bo narzędzie dokumentowe
jest zdolnością platformy. Zasady i obszar składają się więc dla pustego
`konfig.Kontekst{}`, adresu najszerszego poziomu; podstawienie
`session.Zasady{}` z ręki znaczyłoby brak izolacji niezależnie od tego, co
operator ustawił.

Wynik jest widoczny dla operatora, gdy żądanie poda `windowId`: powstaje
wtedy wiersz zasobu w wykazie okna (`oknoWynikuArsenalu`). Brak `windowId`
nie wstrzymuje czynności — bajty i tak idą do magazynu pod sumą kontrolną,
tylko zasób nie pojawia się w wykazie okna. Rodzaj zasobu rozstrzyga format
wyniku, nie nazwa komendy: PDF jest dokumentem, a konwersja do PNG obrazem
(`rodzajZasobuArsenalu`).

Pole `magazyn` trzyma bajty wyników konwersji pod sumą sha256 i jest tym
samym magazynem zasobów Designu, którym jedzie `design.asset.upload`. Wynik
jest typu `DesignAsset`, więc panel zasobów znajduje go tam, gdzie szuka;
skład zasobu poza bazą jest jeden. Pole `biblioteka` jest drugim magazynem, w
którym `assetId` może wskazywać materiał: kontrakt obiecuje przy tym polu
zasób z magazynu rdzenia, Design albo Library, a plik biblioteki trzyma
bajty tą samą konwencją co zasób designu — bezwzględną ścieżką do bloku pod
sumą sha256, różniącą się wyłącznie korzeniem katalogu.

W `ustalZrodlo` kolejność wejścia kontraktu nie jest dowolna: zasób z
magazynu jest treścią zamrożoną pod sumą kontrolną, ścieżka jest treścią
żywą, którą operator może w międzyczasie nadpisać, a treść wprost jest tym,
co model właśnie napisał. Gdy przyszło więcej niż jedno wskazanie, wygrywa
to pewniejsze; brak każdego wskazania jest odmową, a nie pustym wynikiem, bo
czynność bez materiału nie ma czego przetworzyć, a powodzenie bez treści
byłoby powodzeniem czynności, której nikt nie wykonał.

W `zrodloZZasobu` magazyny są dwa, i kontrakt mówi to wprost przy polu
`assetId`: zasób z magazynu rdzenia, Design albo Library. Oba trzymają bajty
tą samą konwencją i różnią się wyłącznie korzeniem katalogu oraz nazwą
kolumny (`URI` zasobu, `TrescOdwolanie` pliku); szukanie idzie najpierw po
zasobach designu, bo tam trafiają wyniki własnych konwersji, a dopiero potem
po bibliotece. Blob nie ma rozszerzenia — jego nazwą jest suma sha256 — więc
format bierze się z kolumny wiersza, a nie ze ścieżki; zasób bez formatu i
bez wskazania w żądaniu jest odmową, bo zgadnięty format wygląda w wyniku
identycznie jak rozpoznany i nie da się ich odróżnić.

## skutek_terminala_test.go

Różnica wobec sprawdzianu koperty jest tu istotą rzeczy: odpowiedź komendy ze
stanem udanym i pustym wynikiem jest odpowiedzią udaną i zarazem kłamiącą, bo
klient czyta odpowiedź, nie komentarz w kodzie. Każdy sprawdzian tego pliku
mierzy więc świat niezależnie od odpowiedzi rdzenia: wpis książki hostów,
pozycję biblioteki i jej wersje — własnym zapytaniem SQL do bazy, nie ponownym
pytaniem tej samej komendy; klucz SSH — plikiem na dysku, jego prawami
i tym, czy biblioteka kliencka potrafi go odczytać; odczyt pliku — treścią,
którą sprawdzian sam wcześniej zapisał; wstrzymanie procesu — stanem procesu
w systemie, nie polem odpowiedzi; obserwację plików — plikiem, który powstał,
bo wyzwolone polecenie naprawdę się wykonało; tunel — stanem końcowym procesu
ssh odczytanym z bazy.

Test karty powłoki urządzeniowej pilnuje szkody, która byłaby cicha w
najgorszy możliwy sposób. Rdzeń nauczył się czterech powłok sięgających poza
jego maszynę — kontenera, poda, konsoli szeregowej i sesji Telnet — a warunek
kolumny powłoka w schemacie bazy przez pewien czas wymieniał inny zestaw
wartości niż kod. Karta takiego rodzaju powstawała wtedy w pamięci
i działała, ale jej zapis odbijał się od warunku bazy po cichu, bez wywrócenia
czynności — karta znikała po restarcie rdzenia, nic tego nie zapowiadając.
Dlatego sprawdzian mierzy wiersz w bazie, nie samą odpowiedź komendy.

## budowa/server/internal/core/adapter_studio_wyrys.go

Podgląd układu bywa robiony tak, że rdzeń startuje przeglądarkę, otwiera w niej
HTML i robi zrzut. Ta droga jest tu zamknięta z tego samego powodu, co
w warsztacie PDF: przeglądarka bezgłowa nie jest częścią instalki, a funkcja
zależna od programu, którego instalka nie niesie, jest u Operatora odmową,
a nie funkcją. Wyrys idzie więc `image`, `image/png` i `x/image/draw`, a PDF
składa `pdfcpu` — wszystko wkompilowane w binarium rdzenia. Krój wczytany
z katalogu systemowego byłby zależnością tego samego rodzaju: u jednego
Operatora podgląd wyszedłby, u drugiego rozsypałby się na prostokąty.
`gofont/goregular` jedzie wkompilowany i niesie łacinkę rozszerzoną, więc
polskie znaki diakrytyczne wychodzą literami, a nie zastępnikami.

Geometria strony jest bytem liczonym z nastaw sekcji, a nie stałymi A4,
ponieważ numer strony w spisie treści, w indeksie i w polu `pageNumber` liczy
się łamaniem wiersza, a łamanie zależy od szerokości kolumny i od wysokości
kartki. Dopóki oba wymiary były stałymi A4, spis treści dokumentu ustawionego
na A5 albo na marginesach szerokich wskazywał strony podglądu, a nie strony
dokumentu.

W `geometriaZNastawStrony` nastawy niepodane biorą wartość domyślną pole po
polu, nie całością: dokument z ustawionym samym nośnikiem ma dostać ten
nośnik z marginesami domyślnymi, a nie A4 z powodu braku marginesu. Nośnik
nazwany bierze wymiary ze wspólnego wykazu rdzenia; wymiar własny (`widthMm`,
`heightMm`) ma przed nim pierwszeństwo, bo jest wskazaniem wprost, a nie
pozycją katalogu. Margines na oprawę (`gutterMm`) dokłada się do marginesu
wewnętrznego, więc zwęża kolumnę tekstu, inaczej oprawa zjadałaby litery,
a nie miejsce na nią. Marginesy odbicia obrotu stron w wyrysie nie zmieniają:
kolumna tekstu ma wtedy tę samą szerokość na obu stronach kartki. Kolumn ta
geometria nie liczy z zamysłu: wyrys rysuje jedną kolumnę, więc zwężenie
łamania do kolumny bez rysowania kolumn dałoby numer strony niezgodny
z obrazem.

`uzdrowiona` traktuje nastawę, która zjada całą kartkę marginesami, jako
wskazanie do poprawienia, a nie powód, żeby rdzeń dzielił przez zero: kolumna
schodzi wtedy do najmniejszej sensownej szerokości, a wyrys wychodzi widocznie
za wąski.

`wierszyNaStrone` odejmuje dwie interlinie na miejsce nagłówka i stopki, bo
oba stoją w wyrysie poza kolumną tekstu, ale kartkę zajmują; rachunek jest ten
sam, którym wyrys układa wiersze, inaczej spis treści wskazywałby inne strony,
niż pokazuje podgląd.

`wyrysujStronyStudia` łamie wiersz po słowach i po zmierzonej szerokości, nie
po stałej liczbie znaków, bo dokument polski ma słowa różnej długości,
a łamanie po znakach rozcinałoby je w połowie i podgląd pokazywałby układ,
którego żaden eksport by nie powtórzył.

`pustaStronaStudia` liczy wymiary z geometrii wyrysu, nie ze stałej A4,
ponieważ porównanie dokumentu na A5 skalowałoby stronę pustą, zamiast ją
nałożyć.

`pdfZeStronStudia` porównuje dwie strony komórka po komórce po jasności, nie
po składowych barwy, bo podgląd jest czarny na białym, a różnica jasności
mówi wprost, czy w komórce coś przybyło albo ubyło. Współrzędne wychodzą
w punktach strony — tych samych, w których mierzone są wymiary wyrysu.
Komórka mniejsza niż `bokKratkiRoznicyStudia` dałaby wykaz obszarów tak długi,
że nikt by go nie przejrzał; większa scaliłaby zmianę akapitu ze zmianą
marginesu.

`dopasujStroneStudia` skaluje stronę tylko wtedy, gdy nastawy wyrysu dwóch
porównywanych wersji się różnią — inaczej porównanie pikselowe stron
o różnych wymiarach pokazywałoby jako zmianę samo przesunięcie.

## budowa/server/internal/core/kompozycja.go

Porty modułowe rozdzielają się od siebie według zależności i rodziny danych,
nie według nazwy modułu w interfejsie. Prowenancja i Zuzycie, Kondycja
i Alerty stoją portami przekrojowymi, bo ich odczyt sięga do kilku okien
naraz, a nie do jednego modułu; Schowek, SkrotyTekstowe i KontekstyPamieci
obsługują tekst poza jednym oknem z tego samego powodu.

WarsztatPdf i BezpieczenstwoDokumentu stoją obok portu Studio, ale osobno,
bo mają inny komplet zależności albo inną rodzinę awarii: warsztat pracuje
biblioteką wkompilowaną w rdzeń, a otwarcie i zamknięcie dokumentu
w bezpieczeństwie nie mają wspólnego powodu do odmowy. Podobny podział
rządzi portami wywoływanymi jako narzędzia modelu z rozmowy — obraz, media,
dokumenty, archiwum, poczta — które stoją osobno od okien Operatora
(Design, Asystent), mimo że część z nich ląduje we wspólnym magazynie
zasobów Designu.

Rzutowanie opcjonalnych rozszerzeń portu (pamięć i planowanie Workspace,
wyposażenie Terminala) idzie dwuwartościowo z rozmysłem: rdzeń złożony
z portem niepełnym ma pracować dalej w pozostałych czynnościach, a brak
rozszerzenia widać w wykazie komend powitania, tak samo jak przy porcie
całkiem pustym.

Pętla wykonawcza Studia wchodzi osobnym wpięciem, bo niesie stan własny —
magazyn rozkładów — i czyta nastawy tym samym obsługiwaczem konfiguracji,
którym czyta je moduł Agents, więc drugiego magazynu nastaw nie zakłada.
Zapora blokad fragmentu owija to, co w rejestrze już stoi, więc jedzie po
obu wpięciach Studia, nie przed nimi — owinięcia nie da się założyć na
komendę, której w rejestrze jeszcze nie ma. Siatka śladu autora owija
zaporę z zewnątrz, żeby czynność zatrzymana blokadą nie zostawiła śladu
w dzienniku.

## budowa/server/internal/core/skutek_rodzin_tekstowych_test.go

Plik mierzy skutek rodzin obsługi tekstu i mowy poza jednym oknem: historii
schowka, słownika skrótów, kontekstów pamięci, nagrań mowy, wywoływacza
poleceń, wyróżnienia wpisu dziennika i pomiaru zajętości okna kontekstu.
Wzorzec sprawdzianu jest ten sam co przy warsztacie PDF: żaden sprawdzian nie
kończy się na odczytaniu odpowiedzi komendy. Każdy schodzi własnym zapytaniem
SQL do tabeli albo otwiera plik na dysku i mierzy go niezależnie od tego, co
komenda zameldowała.

## budowa/server/internal/core/skutek_badania_test.go

Ten plik mierzy skutek modułu Research: czy za odpowiedzią komendy leży byt,
który da się zmierzyć niezależnie od tej odpowiedzi. Żaden sprawdzian nie
kończy się na `status: ok` — po każdej udanej komendzie schodzi do pliku
bazy albo do magazynu bajtów i mierzy tam stan własnym zapytaniem: wiersz
źródła, wiązanie kodu z ustaleniem, wiersz sprzeczności, plik eksportu
o niezerowej długości.

Rodziny zależne od sieci (odkrywanie, rozstrzyganie identyfikatorów) i od
modelu (streszczenie, weryfikacja, operacje kontekstowe) nie wchodzą tutaj
świadomie: ich skutek zależy od świata poza tą maszyną, a sprawdzian ma
mierzyć rdzeń, nie łącze. Wchodzi za to wszystko, co rdzeń robi sam.

## budowa/server/internal/core/montaz_porty.go

`skladPortow` niesie komplet bytów, z których powstają porty rdzenia, w postaci
struktury zamiast dwudziestu pozycyjnych argumentów: przy wywołaniu pozycyjnym
o tylu polach każda zmiana kolejności byłaby cichą pomyłką nie do wyłapania
przez kontrolę typów, bo połowa tych pól ma ten sam typ. Pole `aplikacje` trzeba
zwalniać osobno, bo `apps.preview.start` podnosi nasłuch HTTP, który bez
zamknięcia przeżyłby zatrzymanie rdzenia i zostawił zajęty port. Pole
`dolozenia` i pole `zakresyNarzedzi` niosą adaptery współdzielone z zestawem
narzędzi tury — jeden byt na dwóch czytelników w obu przypadkach.

W `nastawyNadajnika` szyfrowanie jest włączone, dopóki Operator jawnie go nie
zdejmie: wartość domyślna ma chronić, a nie ułatwiać. Zejście do rozmowy
otwartym tekstem ma sens wyłącznie dla przekaźnika na tej samej maszynie
i wymaga jawnego zapisu.

`zlozPorty` przekłada zestaw bytów na porty kontraktu; rozdział względem
`Zmontuj` idzie po odpowiedzialności — `Zmontuj` składa byty i wiąże je ze
sobą, a ta funkcja wyłącznie przekłada je na porty. Katalog danych rdzenia jest
odczytywany raz i przekazywany do obu składów trzymających stan poza bazą:
sejfu poświadczeń i magazynu treści biblioteki. Wartość niesie przełącznik
`-dane` albo zmienną `DANACO_KATALOG_DANYCH`; przy braku wskazania
`konfiguracja.Domyslna()` wstawia katalog domyślny, a konstruktory adapterów
wołane bez tej wartości stoją na `konfiguracja.KatalogDanychDomyslny()` — więc
pominięcie jej rozdzieliłoby stan rdzenia między katalog wskazany a domyślny.

Jeden sejf poświadczeń obsługuje konta i punkty dostępu, zbudowany nad
skonfigurowanym katalogiem danych, a nie domyślnym: jedna instancja pod jednym
zamkiem, bo dwa sejfy nad tym samym plikiem ścigałyby się o zapis. Ten sam sejf
idzie do warstwy modeli — kanał API rozwiązuje odwołanie „sejf:<byt>" przez
uchwyt pakietowy `models.UstawSejfPoswiadczen`, bo kanały powstają fabryką
z samego wiersza rejestru i nie mają jak dostać sejfu argumentem; montaż jest
jedyny w procesie i biegnie przed obsługą pierwszego żądania, więc zapis
uchwytu wyprzedza wszystkie odczyty.

Katalog akcji i adapter przenoszenia kontekstu powstają przed literałem portów,
bo mają po dwóch czytelników: własny port niżej i nakładkę AOD, gdzie
podpowiedzi biorą się z katalogu akcji, a `aod.context.get` z magazynu kompletu
okna — druga instancja każdego z nich byłaby drugą prawdą o tym samym bycie.
Centrum powiadomień powstaje przed literałem portów z tego samego powodu: ma
dwóch czytelników, port `CentrumPowiadomien` i drogę wewnętrzną `Zglos`, którą
rdzeń wnosi do rejestru zdarzenia, gdy coś zaszło; nastawy idą tym samym
adapterem ustawień, co konto nadawcze, bo o tym, czy klasa zdarzenia w ogóle
powiadamia, rozstrzyga sekcja „Powiadomienia". Adapter Asystenta powstaje
wcześniej analogicznie: port `Asystent` i nakładka AOD wołają ten sam moduł,
a składacz mostów idzie ten sam, co do rozmowy, bo wykaz narzędzi jest jeden —
bez `ZMostami` tura zlecenia idzie bez wpisu serwera narzędzi w konfiguracji
MCP, a model prowadzący zlecenie nie ma ani jednego narzędzia kontraktu.
Nakładka AOD dostaje te same byty, którymi jedzie reszta rdzenia: port rozmowy,
moduł Assistant, katalog akcji, magazyn kompletu kontekstu okna i katalog
urządzeń — bez kompletu wpięć komendy rodziny `aod.*` odmawiają mimo
rejestracji.

Jeden adapter okien zasila dwa porty, Okna i Role, bo obie rodziny komend
muszą czytać okno tym samym kompletem wpięć: adapter zbudowany bez `ZWieziami`
ma `wiezie == nil`, więc `dolozWiezi` (adapter_okna.go) kończy się na pierwszym
warunku i `role.list` nie widzi więzi koordynator–wykonawca zapisanej w bazie.
Warsztat PDF stoi w zmiennej, bo bierze go także port bezpieczeństwa — dwa
adaptery na jednym magazynie zasobów, a nie dwa magazyny.

Kondycja mierzy adres, gniazdo, program i kanał modelu — stąd trzy zależności
ponad magazyn sond: uruchamiacz procesów jako sonda programowa, izolacja wraz
z katalogiem roboczym tą samą drogą, co każde inne wołanie arsenału, oraz
rejestr kanałów jako sonda wywołania modelu. Alerty liczą miary z magazynów,
o których warstwa alertu nie ma prawa wiedzieć sama: ślad wywołań modelu,
dziennik błędów i nastawa pułapu kosztu — bez któregokolwiek reguła na tej
mierze nie powstanie, adapter odmawia jej zapisu zamiast milczeć przy
ewaluacji. Wywoływacz trzyma nastawę skrótu w konfiguracji, tam gdzie mieszka
każde inne ustawienie, i czyta deklaracje zdolności klientów, żeby nie
obiecywać skrótu, którego nie ma kto przechwycić. Zajętość okna kontekstu
liczy się z treści, która naprawdę pojedzie do modelu: prompt systemowy
z portu tożsamości, historia rozmowy okna i wpisy pamięci okna; granicę okna
podaje parametr kanału. Kanały biorą sejf poświadczeń wyłącznie dla
`channel.credential.status` i widzą z niego sam odczyt, bo stan poświadczenia
jest pytaniem, czy coś pod odwołaniem leży, a nie prośbą o treść. Warstwa
mobilna jest rozszerzeniem monitora, nie osobnym portem: rodzina `mobile.*`
czyta procesy z tego samego rejestru telemetrii, którym jedzie
`monitor.status` — bez tego ogniwa asercja portu w `handlers_mobile.go` nie
przechodzi i wszystkie trzy komendy odmawiają, choć adapter jest napisany.

Moduł Workspace bierze instrukcje warstwowe z rozstrzygacza jednego na całą
platformę, a bibliotekę projektu z tego samego ustalacza katalogu roboczego,
którym jedzie sesja. Wydobycie treści pliku projektu idzie tym samym
warsztatem, co `document.text.extract`: warstwa tekstowa dokumentu, a po jej
braku rozpoznanie pisma — drugiego czytnika dokumentów w rdzeniu nie ma.
Kolejność ogniw budowy adaptera nie jest dowolna: `ZPamiecia` oddaje adapter
opakowany o rodzinę `memory.*` i musi stać ostatnie, bo ogniwo dopięte po nim
oddawałoby adapter wewnętrzny i port straciłby pamięć — rdzeń nie miałby wtedy
pięciu komend `memory.*`. `ZWylaczeniami` idzie po `ZPamiecia` z tego samego
powodu: inna kolejność zostawiłaby rodzinę `memory.disable.*` bez magazynu.

Ślad wywołań i rozliczenie zużycia stoją na jednym magazynie — dwa porty,
jedno repozytorium; Prowenancja bierze ponadto rejestr kanałów wyłącznie dla
powtórzenia wywołania (`provenance.call.replay`), które jest nowym wywołaniem
kanału, a nie odczytem śladu. Warsztat PDF nie bierze uruchamiacza procesów
ani zasad izolacji, bo pracuje biblioteką wkompilowaną w rdzeń i nie startuje
ani jednego procesu potomnego; bezpieczeństwo dokumentu sięga po ten sam
warsztat, bo materiał wchodzi tą samą drogą, oraz po repozytorium Studia, bo
rozpoznanie danych wrażliwych czyta treść dokumentu, a nie zasób magazynu.

Porty Developer i Debata muszą być wypełnione oba: strażniki
`if d == nil { return }` w `zarejestrujDevelopera` i `zarejestrujDebate`
wychodzą przed rejestracją, więc port pusty daje odpowiedź `*.unknown` na
komendy, których klient ma komplet. Debata dostaje ponadto katalog danych pod
magazyn wydanych transkryptów, grafów i nagrań oraz arsenał pod dwie czynności
wymagające programu serwerowego: zamianę transkryptu na dokument biurowy
(Pandoc) i odsłuch debaty silnikiem mowy — reszta modułu, analiza i głosowanie,
nie startuje ani jednego procesu. Moduł Automations buduje definicje, a
wykonuje je tym samym adapterem kolejek, którym jedzie domena kolejek, bo
silnik jest jeden; instancja powstaje w złożeniu modułów (montaz_moduly.go),
bo dzieli budzik harmonogramu — gdyby port i budzik miały osobne adaptery,
przypięcia obserwatorów przebiegów rozjechałyby się na dwie mapy.

Biblioteka, Studio, Przeglądarka i Design dostają wyłącznie swoje repozytoria
i nie dzielą stanu. Biblioteka bierze ponadto katalog danych, ten sam, nad
którym stoi sejf poświadczeń, bo baza trzyma wyłącznie odwołanie do treści
i bajty wgranych plików muszą leżeć tam, gdzie reszta stanu rdzenia; arsenał
i rejestr kanałów są jej potrzebne w dwóch miejscach, pomiar czasu trwania
nagrania przy odczycie metadanych osadzonych oraz klasyfikacja wsadowa. Studio
dzieli rejestr kanałów z oknem rozmowy i modułem Roundtable, bo drugiego
silnika modelu nie ma nigdzie, a magazyn zasobów jest ten sam, którym jedzie
warsztat PDF i moduł Design — archiwum historii, paczka redakcyjna i strony
podglądu są zasobami tej samej platformy i leżą w jednym miejscu; repozytorium
Library wchodzi wyłącznie do odczytu, bo `studio.diff.source` zestawia
dokument roboczy z materiałem wejściowym. Przeglądarka dostaje ponad własne
repozytorium katalog danych — magazyn bajtów zrzutów, archiwów i rejestrów
sieciowych, ten sam korzeń co magazyn biblioteki i zasobów Designu — oraz
uruchamiacz procesów wraz z izolacją dla silnika Chromium prowadzonego
protokołem CDP, i nic więcej: moduł nie zna innych modułów.

Design bierze rejestr kanałów tym samym sposobem co Roundtable i Studio: most
do modelu dla pracy tekstowej, budowy promptu i opisu zasobu — rdzeń nie
generuje obrazów, rejestr służy wyłącznie operacjom słownym modułu. Design
dostaje ponadto katalog danych, ten sam, którym jadą sejf poświadczeń
i magazyn biblioteki, bo `design.asset.upload` trzyma bajty zasobów poza bazą;
sejf, ten sam co Konta i PunktyDostepu, bo klucze darmowych baz zdjęciowych
(`design.stock.*`) leżą w nim pod bytem `design.stock.<dostawca>` i moduł
wyłącznie je czyta; oraz drogę odczytu pisma (`ZOdczytemPisma`) — uruchamiacz
i bramę izolacji dla jednej czynności, `design.mockup.import`, która czyta
treść napisów ze zrzutu programem pakietu serwera, bo czytnika liter w czystym
Go nie ma. Reszta modułu procesów nie startuje i pilnuje tego zapora.

Assistant dostaje wykonawcę zleceń: rejestr kanałów jako droga modelu,
nadzorcę sesji jako okno zlecenia i nadajnik z kontekstem życia dla strumienia
tury i `assistant.action.changed` — bez nich zlecenie zostaje `queued` zamiast
się wykonać; rozpoznanie mowy (`ZMowa`) idzie tym samym adapterem, który
wypełnia port `Mowa`, bo silnik jest jeden.

Apps dostaje rejestr okien, bo wdrożenie bez istniejącego okna nie ma
przestrzeni roboczej, z której miałoby cokolwiek wziąć — bez rejestru komenda
zakłada przebieg dla okna, którego nie ma, i kończy go powodem o pustym
warsztacie zamiast o braku okna. Apps dostaje ponadto katalog danych, magazyn
bajtów eksportu, artefaktów i pakietów oraz sejf, z którego bierze się klucz
wydawcy, i rejestr pozycji katalogu — `apps.package.publish` publikuje do
tego rejestru, nie do drugiego obok niego. Komponenty własne dostają swój
rejestr oraz trzy magazyny modułowe, w których `component.create` zakłada byt
docelowy: projekty, ekspertów i automatyki; czwartego, profili asystenta,
platforma nie ma, więc rodzaj `assistant` odmawia z powodem, zamiast zakładać
kafel wskazujący na nic.

Translate dostaje rejestr kanałów jako most do modelu, którym tłumaczenie
faktycznie woła kanał zamiast zakładać puste panele; wybór kanału niesie już
żądanie kontraktu (pole `channelId`), a kanał wskazany, nieznany albo
nieczynny kończy się odmową nazwaną, nie cichym zejściem na kanał domyślny.
`ZSynteza` wpina silnik syntezy mowy dla `translate.speech.synthesize`: ten
sam uruchamiacz i te same dwa źródła izolacji, którymi jadą Terminal,
Developer i silnik rozpoznawania mowy, plus katalog danych rdzenia jako
miejsce na nagrania — ten sam, który dostają sejf poświadczeń i magazyn
treści biblioteki. `ZWytworami` wpina drogę `translate.artifact.publish`:
repozytorium biblioteki jako wiersz pliku i katalog danych rdzenia jako
magazyn bajtów pod sumą kontrolną — bez niej wytwór nie miałby gdzie leżeć ani
czym się zgłosić reszcie platformy, więc sama ta komenda odmawia nazywając
brak.

Przekazanie okna zakłada pozycję kolejki, więc bierze ten sam adapter kolejek,
którym jedzie domena kolejek i moduł Automations, bo silnik jest jeden — bez
niego `window.handoff` odmawia wprost. Cztery repozytoria wchodzą, bo
przekazanie okna domyka lukę między kontraktem a schematem: kontrakt niesie
identyfikatory zewnętrzne sesji i okien, a schemat wiąże klucze wewnętrzne
i słowniki modułów oraz kanałów — bez tych czterech repozytoriów adapter nie
umiałby ani odnaleźć okna wskazanego przez klienta, ani oddać okna w kształcie
kontraktu.

Sejf Uwierzytelnienia jest tym samym co Konta i PunktyDostepu — jedna
instancja pod jednym zamkiem, bez tego ogniwa sekret bramki nie ma gdzie
leżeć i logowanie odmawia zawsze. Wiązanie z kontem właściciela niesie login,
adres uwierzytelniający i drogi potwierdzenia, bez tego ogniwa rejestracja
odmawia, bo konta nie ma gdzie zapisać. Nadajnik jest kontem nadawczym
platformy, nie skrzynką Operatora — nim idą dwa listy systemowe, potwierdzenie
adresu i droga odzyskania konta; nastawy idą tym samym adapterem ustawień,
którym idzie każda inna nastawa platformy, więc konto nadawcze zapisane
w oknie Konfiguracji przesłania to ze startu i pomyłkę w adresie serwera
poczty naprawia się bez zatrzymywania rdzenia. Urządzenia stoją na tym samym
repozytorium, bo urządzeniem konta jest to, które weszło przez bramkę — wykaz
bierze się z sesji bramki.

Katalog rozszerzeń — warstwa danych wystawia go metodą, nie polem
(dane/extension.go), bo rejestr nie trzyma stanu poza wskaźnikiem na wspólną
pamięć zapytań; tak samo role okien, sekcje paneli i skrzynki poczty niżej.
Katalog dostępu i biblioteka ekspertów idą razem z rozgłoszeniem, bo bez nich
adapter odmawia `internal_error` każdemu żądaniu niosącemu `accessPointId`
albo `agentId`, a wskazanie puste przechodzi; katalog danych wchodzi tą samą
drogą, co do modułu Apps, bo paczka przesłana instalacją Personal i dziennik
piaskownicy są bajtami na dysku. Role okien bierze cały zestaw, nie samo
repozytorium, bo adapter musi dosięgnąć także przekazań (więź koordynatora)
i adresów okien — stoi na tym samym adapterze okien, który wypełnia port
Okna, bo drugi byłby drugą prawdą o oknie. Historia rozmowy okna bierze
jedno repozytorium, swoje: okien ani sesji nie dobiera, wiąże je zapytanie po
identyfikatorze kontraktowym, a drugi czytelnik tych samych wierszy byłby
drugą prawdą o wypowiedzi.

Podagenci biorą komplet wiązań z jednego miejsca
(adapter_modul_orkiestracja_zlozenie.go: zlozPodagentow): repozytorium
podagentów i okien, biegi orkiestracji, ten sam adapter kolejek co domena
kolejek i Automations oraz ocenę żywotności z rejestru procesów sesji; montaż
nie zna pakietu `podagenci`, wiedzę o nim trzyma adapter, który jako jedyny go
używa. Zespoły biorą własne repozytorium oraz dziennik rdzenia, ten sam,
którym mówi reszta montażu — skład, z którego wypadł ekspert usunięty albo
zarchiwizowany, wraca do klienta krótszy, bo kontrakt nie ma pola, którym
dałoby się o tym powiedzieć. Historia i archiwum eksperta biorą trzy
repozytoria: dwa widoki wersji oraz bibliotekę ekspertów, która służy
wyłącznie oddaniu eksperta w kształcie kontraktu po zmianie wersji. Zakres
działania eksperta bierze pięć źródeł: własne repozytorium zakresu,
bibliotekę dla kształtu kontraktu po zmianie, historię wersji dla podglądu
migawki, katalog modułów dla sprawdzenia wskazania i katalog punktów dostępu
dla konfiguracji instancji konektora — rozstrzygacz dokłada dziedziczenie
izolacji do polityki efektywnej, ten sam, którym jedzie okno konfiguracji
punktów izolacji. Zakresy narzędzi stoją na własnym repozytorium i na
katalogu profili asystenta — profil pusty w żądaniu bierze profil domyślny.

Silnik mowy powstaje w złożeniu modułów (montaz_moduly.go), bo dzieli
uruchamiacz procesów i oba źródła izolacji z Terminalem i Developerem, a tę
samą instancję bierze moduł Assistant — drugi silnik byłby drugą prawdą o tym,
czy platforma rozpoznaje mowę. Wskaźnik znaczenia stoi na tym samym
uruchamiaczu i tych samych dwóch źródłach izolacji, co Terminal, Developer
i silnik mowy, bo pomocnik osadzeń jest procesem drzewa jak każdy inny;
katalog danych jedzie ten sam, którym jadą sejf, magazyn biblioteki i zasoby
Designu, bo tam wykłada się pomocnik i tam lądują wagi modelu. Składnica
wektorów bierze `*sql.DB` z montażu, tak samo jak dziennik transkrypcji
i dziennik doradcy, bo `dane.Zestaw` uchwytu bazy nie wystawia — stan trzyma
tabela `fragment_wiedzy`; źródła treści idą przez repozytoria modułu Library
i historii, bo drugiej drogi do tych wierszy nie ma.

Narzędzia obrazu, media, dokumenty i archiwum dzielą jedną trójkę zależności:
repozytorium Designu, bo wynik każdej z tych czynności jest zasobem tej samej
platformy; katalog danych jako magazyn bajtów wyniku, ten sam, którym jadą
sejf, magazyn biblioteki i zasoby Designu; oraz arsenał wpięty tym samym
uruchamiaczem i tymi samymi dwoma źródłami izolacji, co Terminal, Developer
i silniki mowy. Narzędzia dokumentowe biorą repozytorium Designu wyłącznie do
odczytu, żeby rozwiązać `assetId` żądania na bajty. Silniki neuronowe obrazu
(`image.upscale`, `image.background.remove`) stoją na tym samym zapleczu, co
rodzina `image.*`: zaplecze jest bezstanowe, więc druga instancja niczego nie
rozdwaja, dokładają się do niego wyłącznie wagi sieci i granice czasu liczone
w minutach. Narzędzia archiwum biorą katalog roboczy podwójnie: raz jako
obszar wołania `7z`, raz jako jedyny układ odniesienia dla ścieżek żądania —
bez niego `archive.*` nie miałyby względem czego rozstrzygać `targetPath`.

Poczta bierze repozytorium skrzynek wystawione metodą (dane/poczta_skrzynki.go)
oraz trzy dalsze wiązania będące tymi samymi bytami, którymi jedzie reszta
rdzenia: ten sam sejf poświadczeń co konta, punkty dostępu i bramka, bo
poświadczenie skrzynki nie ma innej drogi niż sejf; ten sam katalog danych,
co magazyn biblioteki i zasoby Designu; i to samo repozytorium Designu, do
którego pisze `design.asset.upload`, bo załącznik listu wciągnięty osobną
drogą byłby zasobem, którego arsenał obrazu i dokumentów nie widzi.

Doraźne dołożenia sesji: adapter powstaje w montażu (montaz.go), bo ten sam
byt wnosi dołożenia do zestawu narzędzi tury — port oddaje je Operatorowi,
składacz tury oddaje je modelowi, a prawda ma być jedna. Trzy wiązania niesie
już konstruktor: dołożenia sesji (migracja 122), wykaz sesji, bo dołożenie
żyje w stanie sesji i bez przekładu identyfikatora kontraktowego nie ma gdzie
usiąść, oraz katalog rozszerzeń, z którego bierze się druga połowa wykazu po
ukośniku. Doradca powstaje w złożeniu modułów (montaz_moduly.go) razem z mową:
stoi na tym samym rejestrze kanałów, którym jedzie okno rozmowy, i na
dzienniku bazy.

## urzadzenia_skaner.go

Rozdzielenie warstwy skanera warunkiem budowy per system, zamiast rozstrzygania
po runtime.GOOS w jednym pliku, wygląda porządniej, ale ma jedną wadę
rozstrzygającą: budowa na maszynie budującej pod Linuksem nie skompilowałaby
ani razu drogi WIA. Droga Windows przestałaby się kompilować przy pierwszej
zmianie sąsiedniego pliku i nikt by tego nie zobaczył aż do wydania instalki
natywnej — dokładnie w chwili, w której nie ma już czasu na naprawę.
Rozstrzygnięcie po runtime.GOOS trzyma oba warianty pod jednym sprawdzianem
kompilacji kosztem kilku bajtów martwego kodu w wydaniu; to ta sama droga,
którą rdzeń rozróżnia system już w innych miejscach, nie druga jej odmiana.
Warunek budowy będzie właściwy dopiero wtedy, gdy droga Windows sięgnie po
bibliotekę wołającą COM z Go — dziś woła PowerShell, więc nie ma czego
chronić warunkiem budowy: kod jest zwykłym napisem i kompiluje się wszędzie.

WIA idzie przez PowerShell, nie TWAIN, ponieważ TWAIN wymaga okna i pętli
komunikatów, a z procesu serwera bez pulpitu nie wystartuje. WIA jest warstwą
systemową Windows dostępną przez COM, a jedyną drogą do COM, którą rdzeń ma
bez wkompilowanej biblioteki, jest PowerShell — już wpisany do wykazu
zależności rdzenia, więc sonda startowa mówi o jego braku sama. Zasada
bezwzględna: program zewnętrzny idzie wyłącznie przez warstwę wołania
procesów rdzenia, nigdy wywołaniem bezpośrednim.

Funkcja odmowaSkanuSane rozstrzyga, czy skan nie doszedł do skutku z braku
urządzenia czy z innej przyczyny, i dopiero wtedy nazywa brak. Bez tego
rozstrzygnięcia obie sytuacje wychodziły jednym zdaniem: program scanimage
kończy się tym samym kodem przy każdej przyczynie, a odmowa arsenału
przekładała kod niezerowy na usterkę wewnętrzną wraz ze zrzutem procesu —
odbiorca bez podłączonego skanera dostawał więc kod mówiący o usterce
rdzenia zamiast zdania o tym, czego brakuje. Droga WIA rozróżnia te dwie
rzeczy od początku, i to samo należy się drodze SANE, bo stan maszyny jest
ten sam. Rozstrzygnięcie idzie pytaniem o wykaz, nie czytaniem diagnostyki
programu: zdanie, które program mówi o sobie, jest napisem obcego programu,
a rdzeń nie ma prawa opierać kodu odmowy na tym, że napis nie zmieni się przy
następnym wydaniu. Wykaz pusty znaczy brak urządzenia; wykaz niepusty albo
niedostępny zostawia odmowę pierwotną, bo rdzeń nie wie wtedy nic ponad to,
co powiedział program, a odmowa zgadnięta byłaby gorsza od surowej. Pytanie
idzie wyłącznie po nieudanym skanie, więc droga udana nie płaci za nie ani
jednym wywołaniem.

## budowa/server/internal/core/adapter_kolejki_zlecenia.go

Ten plik nie jest drugim silnikiem kolejek: cykl życia kolejki jako całości
prowadzi `kolejka_silnik.go` przez `queue.action`, i tylko on. Tutaj żyje
zlecenie — byt, którego kolejka jako całość nie zna: ładunek strukturalny,
priorytet, termin wykonania, klucz idempotencji, warunek przetworzenia.

Idempotencja jest sprawdzana przed założeniem zlecenia i pole `duplicate`
mówi o tym wprost. Zlecenie o kluczu już użytym nie zakłada drugiego
wiersza i nie jest odmową: wywołanie przychodzące powtórzone przez nadawcę
dostaje odpowiedź „to już jest”, a nie błąd, na który nadawca odpowie
kolejnym powtórzeniem.

Zlecenie zdjęte, scalone i podzielone nie znika z bazy — dostaje stan
końcowy. Historia kolejki ma pokazywać, co się z ładunkiem stało, a wiersz
skasowany nie pokazuje niczego.

## budowa/server/internal/core/adapter_narzedzia_obraz_wkompilowany.go

Te cztery czynności liczył wcześniej program zewnętrzny, choć każdą z nich
wykonuje w całości biblioteka Go wkompilowana w binarium. Program zewnętrzny
wołany tam, gdzie biblioteka wystarcza, jest regresem: kosztuje uruchomienie
procesu, wiąże funkcję z wersją cudzego wydania i przy niekompletnym serwerze
zamienia retusz w odmowę. Rachunek stoi więc w tym pliku i idzie w procesie:
`disintegration/imaging` (skalowanie Lanczosem, kadr, obrót, odbicia, korekcje
barwne, rozmycie, wyostrzenie), `golang.org/x/image` (dekodery WEBP, TIFF, BMP
oraz kodery TIFF i BMP), `HugoSmits86/nativewebp` (zapis WEBP bezstratnego),
`rwcarlsen/goexif` (odczyt metadanych EXIF) i rachunek własny na odszumianie
medianą oraz rozciągnięcie poziomów.

Dwa wyjścia nie mają w Go kodera i nie da się ich tu policzyć: AVIF, którego
kodera czysto-Go nie ma wcale, oraz WEBP stratny, bo `nativewebp` zapisuje
wyłącznie bezstratny VP8L. Tak samo AVIF nie ma dekodera, więc obraz w tym
formacie nie wchodzi. Te przypadki oddają `errBrakRachunkuGoObrazu`, a
czynność sięga wtedy po program pakietu serwera — jedyna droga, jaka zostaje.
Zapora `zapora_narzedzi_obrazu_test.go` pilnuje, żeby ta droga została
wyjątkiem nazwanym, a nie wróciła jako droga podstawowa.

Nazwy przestrzeni barw w `przestrzenBarwArsenalu` zostają te, które produkt
wypisywał do tej pory, bo czyta je model w treści odpowiedzi — zmiana
słownika byłaby zmianą kontraktu przy okazji zmiany rachunku.

`przeskalujObrazArsenalu` traktuje obie miary podane przy zachowanych
proporcjach jako „zmieść się w tej ramce" (`Fit`), a nie „rozciągnij do niej",
zgodnie z tym, co mówi kontrakt: model prosi zwykle o „szerokość 800", a nie
o rozciągnięcie zdjęcia.

Brak pola `amount` w `poprawObrazArsenalu` bierze wartość domyślną operacji,
nie zero: zero byłoby poprawką bez skutku, a model proszący o rozjaśnienie
bez liczby dostałby obraz nieodróżnialny od źródła.

`sigmaZSilyArsenalu` dla zera i wartości ujemnych daje najmniejszą sigmę
o skutku widocznym, bo sigma zerowa nie zrobiłaby nic, a czynność ma robić to,
o co poproszono.

`odszumMedianaArsenalu` liczy medianą, nie rozmyciem, bo szum pojedynczych
punktów (sól i pieprz z matrycy przy wysokiej czułości) jest wartością
odstającą, a mediana odstających nie bierze. Krawędzie zostają ostre, bo po
obu ich stronach mediana wskazuje wartość strony liczniejszej — czego rozmycie
Gaussa nie robi. Kanał alfa liczy się tą samą medianą co barwy, osobno, żeby
mieszanie z barwami nie zmieniło przezroczystości na krawędziach wycięcia.

`rozciagnijPoziomyArsenalu` rozciąga tylko wtedy, gdy zakres jest węższy niż
pełny: obraz już rozciągnięty przeszedłby przez mnożenie bez zmiany,
a dzielenie przez zero przy obrazie jednobarwnym wywróciłoby rachunek.

Odmowa dekodera w `odczytajObrazArsenalu` znaczy „nie ma czym tego przeczytać
w procesie", a nie „plik jest zepsuty": pod tą samą odmową kryje się AVIF,
którego dekodera w Go nie ma. Rozstrzyga to wołający, przechodząc na program
pakietu serwera — gdyby plik był naprawdę uszkodzony, tamta droga powie to
wprost.

## skutek_przestrzeni_roboczej_test.go

Wzorzec szkody, którego pilnują sprawdziany tego pliku, ma w produkcie
precedens: odpowiedź udana z pustym wynikiem przechodzi każdy sprawdzian
zgodności z kontraktem i nie mówi nic o tym, czy cokolwiek zostało. Dlatego
żaden sprawdzian tutaj nie kończy się na odczytaniu odpowiedzi komendy: każdy
schodzi własnym zapytaniem SQL do pliku bazy albo do pliku na dysku i mierzy
niezależnie od tego, co rdzeń zameldował.

## budowa/server/internal/core/przegladarka_pobieranie.go

`browser.navigate` sięga po stronę realnym HTTP GET-em (biblioteka
standardowa `net/http`, bez zależności), z limitem czasu i rozmiaru, i
wydobywa z HTML-a tytuł oraz tekst renderowany, którymi wypełnia migawkę
opisaną w `adapter_modul_przegladarka.go`. Strona może milczeć, ale rdzeń nie
udaje, że odpowiedziała: błąd transportu, stan nienoszący treści (4xx, 5xx,
204/205, przekierowanie bez `Location`) albo treść nietekstowa wraca uczciwym
błędem — migawka powstaje tylko z faktycznie pobranej strony, nie z pustki
podszytej pod sukces. Zrzut ekranu wymaga silnika przeglądarki poza
`net/http` i nie jest tu wypełniany.

Rdzeń stoi na maszynie operatora, więc adres wskazany przez model nie jedzie
w świat bez granic. Cztery granice, każda odmawiająca zdaniem po polsku, nie
ciszą i nie wyjątkiem: protokół (tylko `http`/`https`, żadnego `file:`,
`ftp:` ani `data:`), czas całego pobrania, rozmiar odpowiedzi i długość
łańcucha przekierowań. Odmowa opisuje brak, nie zakaz postawiony operatorowi.

Adres jest obcinany raz, na wejściu funkcji `pobierzStrone`, i to jest cała
prawda o nim dalej: gdyby odstępy zdejmował sam sprawdzian protokołu, a do
złożenia żądania szedł łańcuch nieobcięty, adres z otaczającymi spacjami
przeszedłby granicę protokołu i rozbiłby się dopiero o
`http.NewRequestWithContext` zdaniem o pierwszym segmencie ścieżki — dwie
prawdy o jednym adresie w jednej funkcji.

`http.Client` zawija powód niepowodzenia w `*url.Error`; rozwinięcie sprawia,
że zdanie odmowy niesie samą przyczynę, nie powtórzony adres z prefiksem
„Get". Odmowa własna z `CheckRedirect` wraca tą samą drogą co usterka
transportu, ale nią nie jest: zły protokół kolejnego skoku i zbyt długi
łańcuch to wada adresu, nie milczenie gospodarza.

Deklarowany rozmiar sprawdzamy przed czytaniem: gdy serwer sam mówi, że
przysyła więcej niż granica, nie ma po co ciągnąć ani bajta. Czytanie idzie
o bajt więcej niż granica, bo nadmiarowy bajt jest dowodem, że strona się nie
zmieściła — bez niego ucięcie byłoby nie do odróżnienia od strony, która ma
dokładnie tyle treści, a wtedy migawka kłamałaby po cichu.

`sprawdzStanOdpowiedzi` istnieje, bo sam warunek „poza 200-399 odmawiaj"
przepuszczałby dwie odpowiedzi, które treści strony nie niosą wcale: 3xx,
które dojechało aż tutaj tylko wtedy, gdy przekierowanie nie wskazało
`Location`, oraz 204/205, które z definicji nie mają ciała. Obie skończyłyby
się stanem powodzenia i migawką bez tytułu, tekstu i HTML-a, a
`browser.snapshot.get` oddałby tę pustkę modelowi jako bieżący stan strony,
a historia nawigacji zapisałaby przejście, którego nie było.

W `pilnujPrzekierowan` granica znaczy dokładnie to, co mówi odmowa: pole
niosące żądania już wysłane zawiera przy pierwszym przekierowaniu żądanie
pierwotne, więc przy k-tym skoku długość wynosi k. Warunek ponad pięć
odrzucałby zatem skok piąty, czyli przepuszczał tylko cztery, podczas gdy
zdanie odmowy mówi „dalej niż pięć razy" — miara i zdanie muszą być tą samą
prawdą, więc odmowa pada dopiero za granicą.

Kanały wchodzą tą samą drogą co strony: `browser.feed.subscribe` pobiera
dokument RSS, Atom albo JSON Feed tym samym pobraniem, a serwery kanałów
deklarują je własnymi typami (`application/rss+xml`, `application/atom+xml`,
`application/feed+json`, `application/json`). Bez tych typów subskrypcja
odmawiałaby zdaniem „zasób nie jest stroną do odczytu" przy dokumencie, który
jest dokładnie tym, o co poprosił operator.

Wydobycie tekstu renderowanego biblioteką standardową, bez zależności
`x/net/html`, nie jest pełnym silnikiem renderującym — to uczciwe wydobycie
treści czytelnej dla modelu: usunięcie bloków skryptów i stylów, zdjęcie
pozostałych znaczników, odkodowanie encji i zbicie odstępów.

## skutek_wydania_studia_test.go

Sprawdzian zgodności z kontraktem pyta, czy odpowiedź jest odpowiedzią; ten
plik pyta o co innego i tylko o to: czy za odpowiedzią coś zostało. Dlatego
archiwum jest tu rozpakowywane i czytane wpis po wpisie, a nie liczone z pola
odpowiedzi, które komenda wypełnia sama; wyrys strony jest dekodowany
z powrotem jako obraz, bo plik nazwany „png" i plik będący obrazem to dwie
różne rzeczy; gałąź i odwołanie do wersji są odczytywane drugim, niezależnym
połączeniem do pliku bazy — rdzeń, który melduje zapis, a wiersza nie
zakłada, przechodzi każdy sprawdzian pytający sam siebie.

Test skanowania z urządzenia sprawdza, że nazwanie braku niesie trzy rzeczy
naraz: czym rdzeń szukał, po czyjej stronie leży brak — maszyna odbiorcy, nie
usterka rdzenia — i jaka droga działa mimo niego. Odmowa bez tych trzech
członów zostawia odbiorcę tam, gdzie zostawiał go pusty wykaz. Mierzony jest
jeden stan maszyny: warstwa skanera jest, a urządzenia nie ma; oba warunki
sprawdzają się przed pomiarem, bo maszyna bez programu skanującego i maszyna
z podłączonym skanerem prowadzą tę czynność innymi drogami, więc pomiar
wykonany w niewłaściwym miejscu mierzyłby coś innego i meldował to jako wynik.

## budowa/server/internal/core/skutek_narzedzi_tresci_pisanej_test.go

Ten plik mierzy skutek sześciu programów treści pisanej: składu PDF-u,
odczytu formatu spoza słownika rdzenia, korekty językowej dwoma silnikami
pisowni, analizy prozy i obróbki wstępnej skanu. Żaden sprawdzian tego pliku
nie kończy się na kopercie udanej, bo programy zewnętrzne mają w tym
produkcie własną odmianę szkody „ok przy pustym wyniku”: każdy z nich
potrafi skończyć się kodem zero i nie zrobić nic — typst zapisze PDF bez
treści, Tika odda pustkę dla rodzaju pliku, którego nie zna, vale bez
konfiguracji zamelduje jeden „błąd wykonania” zamiast ustaleń, a unpaper
przepisze obraz bez żadnego filtru. Dlatego każdy sprawdzian pyta o skutek:
czy w PDF-ie da się odczytać zdanie, które do niego weszło, czy ustalenie
niesie rodzaj i propozycję, czy tekst rozpoznany ze skanu niesie słowa
materiału.

## budowa/server/internal/core/adapter_narzedzia_obraz.go

Same czynności czterech narzędzi obrazu leżą w
`adapter_narzedzia_obraz_czynnosci.go`, rachunek na pikselach —
w `adapter_narzedzia_obraz_wkompilowany.go`, port i wpięcie —
w `handlers_narzedzia_obraz.go`. Pracę wykonuje głównie biblioteka
wkompilowana w binarium; program pakietu serwera (ImageMagick) zostaje drogą
zapasową wyłącznie dla dwóch wyjść, których żaden koder czysto-Go nie
zapisze — AVIF i WEBP stratny — oraz dla pomiaru pliku AVIF, którego nie ma
czym zdekodować. Wykaz zależności pakietu serwera niesie ImageMagicka z tym
właśnie zakresem, a osobny sprawdzian pilnuje, żeby ta droga nie wróciła jako
droga podstawowa.

Wołanie binarium idzie wyłącznie przez `zewnetrzne.Wolaj`: stamtąd prowadzi
port `session.Uruchamiacz`, brama izolacji okna i objęcie drzewa procesów.
Własnego `exec.Command` w tym pliku nie ma — odstępstwo od tej sekwencji
kończy się wyciekiem procesu albo uchwytu.

Bajty wyniku lądują w tym samym magazynie zasobów, co `design.asset.upload`
(blob pod sumą sha256), a wiersz — w tej samej tabeli zasobów Designu. Drugi
magazyn byłby drugą prawdą o tym, gdzie rdzeń trzyma bajty poza bazą, a
Assets Panel przestałby widzieć połowę zasobów, które sam wytworzył.

Źródło zostaje nietknięte: wynikiem każdej z trzech czynności zmieniających
jest nowy zasób. Blob źródła leży pod swoją sumą kontrolną i nikt go tu nie
otwiera do zapisu — rachunek czyta go, a bajty wyniku składa osobno; na
drodze zapasowej program dostaje ścieżkę do odczytu, a wynik oddaje na
standardowe wyjście.

Brak binarium na drodze zapasowej jest odmową nazwaną: pakiet zewnętrzny
oddaje wtedy sygnał braku narzędzia niosący nazwę programu i pakiet do
doinstalowania. Rdzeń przekłada go na kod niedostępności kanału — ten sam,
co przy braku silnika mowy — a nie na cichy zasób bez zmian, który
wyglądałby jak udany retusz. Dotyczy to wyłącznie AVIF-a i WEBP-a stratnego:
pozostałe czynności rodziny nie mają czego zabraknąć.

ImageMagick siódmy stoi jednym plikiem `magick`, szósty — osobnymi
`convert` i `identify`. Adapter pyta po kolei, który program stoi
w systemie, i pierwszy obecny wygrywa; gdy nie ma żadnego, odmowa nazywa
tryb siódmy, bo to on jest wskazówką do instalacji. Nazwa podpolecenia
(na przykład `identify`) wchodzi na początek argumentów tylko wtedy, gdy
wołany jest plik `magick`; wersja szósta ma na to osobne binarium
i podpolecenia nie rozumie.

## budowa/server/internal/core/skutek_kondycji_i_alertow_test.go

Rodziny kondycji i alertów mają najostrzejszy warunek prawdziwości w całym
produkcie: komenda zdrowia, która oddaje stan sprawności nie zmierzywszy
niczego, jest gorsza niż jej brak — to fasada, przez którą awaria przechodzi
niezauważona. Dlatego żaden sprawdzian w tym pliku nie kończy się na
odczytaniu odpowiedzi komendy. Każdy schodzi własnym zapytaniem SQL do tabel
`sonda_kondycji`, `wynik_sondy_kondycji`, `regula_alertu`
i `wyzwolenie_alertu` i sprawdza, czy za odpowiedzią stoi wiersz — a w
wierszu wartość, której nikt nie wpisał z ręki.

## skutek_petli_wykonawczej_test.go

Sprawdziany tego pliku mierzą skutek pętli wykonawczej modułu Studio, czyli
stan, który zostaje po komendzie, a nie kopertę odpowiedzi: stan zadań
odczytany osobnym wywołaniem `studio.plan.get`, bilans przebiegu i powód
odmowy. Każdy sprawdzian wyklucza jedną konkretną szkodę, którą ten produkt
już poniósł: plan zameldowany jako uruchomiony, który nie wykonał ani jednego
zadania i nie powiedział o tym ani słowa; pętlę puszczoną przy wyłączonej
nastawie, która albo cicho nic nie robi, albo melduje powodzenie, którego nie
ma; wsad na wielu dokumentach, meldujący jedną liczbą przyjętych i milczący
o losie każdego dokumentu osobno; zadanie pominięte przez blokadę fragmentu,
o którym operator nie dowiaduje się z kolejki, bo powód został przemilczany;
oraz zatrzymanie, po którym zadanie w biegu znika, zamiast zostać opisane
jako przerwane.

`TestPetlaPrzyWylaczonejNastawieOdmawiaNazywajacBrak` sprawdza dwie miary
naraz, bo pojedyncza dałaby się obejść: powód odmowy pokazuje operatorowi,
czego brakuje, a stan zadań odczytany osobnym wywołaniem pokazuje, że pętla
naprawdę nie tknęła kolejki. Odpowiedź „started: false” przy zadaniach
przestawionych na „done” byłaby kłamstwem w drugą stronę.

`TestPetlaUruchomionaWykonujeZadaniaAlboNazywaBrak` prowadzi zadanie rodzaju
`export` rodziną wydania, która wykonuje się do skutku, oraz zadania treści
kanałem modelu, którego stanowisko sprawdzianu nie ma; to czyni je miarą
najostrzejszą, bo pętla ma je zamknąć stanem `failed` z powodem nazywającym
brak, a nie zameldować gotowość.

`TestZatrzymaniePetliZostawiaSladPrzerwaniaINiegubiPracy` sprawdza to samo
zatrzymanie od strony pracy już wykonanej: zatrzymanie, które kasowałoby
zadania domknięte przed nim, byłoby gorsze niż brak zatrzymania, bo operator
naciskający przerwanie traciłby to, co pętla zrobiła dobrze przed komendą
przerwania.

`TestBlokadaFragmentuWidocznaWKolejceZadan` sprawdza przy okazji, że pętla
wykonawcza sięga do zapory blokad przez rejestr, a nie wywołaniem adaptera
wprost: pętla omijająca rejestr przepisałaby zablokowany fragment, a
sprawdzian oparty tylko na kopercie odpowiedzi by tego nie zobaczył, bo
zadanie skończyłoby się jako gotowe.

## budowa/server/internal/core/adapter_narzedzia_dokument_tekst.go

Pole `usedOcr` odpowiedzi mówi modelowi, czy wolno mu zacytować zwrócony
tekst jako fakt. Wartość fałszywa znaczy, że tekst pochodzi z warstwy
tekstowej dokumentu — to te same znaki, które wpisał autor, więc cyfra „0”
nie zamieni się w literę „O”, a kwota nie zgubi przecinka, i model może
cytować dosłownie. Wartość prawdziwa znaczy, że tekst odczytano z pikseli:
rozpoznanie pisma myli znaki podobne, gubi kolumny i wymyśla spacje, więc
model ma traktować taki tekst jak relację świadka i przy cytowaniu zaznaczyć,
skąd treść pochodzi. Wartość `usedOcr` wynika wyłącznie z tego, która droga
dała treść, nigdy z długości odczytanego tekstu. Dlatego dla PDF-u kolejność
jest jedna i nieodwracalna: najpierw warstwa tekstowa, dopiero po jej braku
albo na wyraźne żądanie Operatora rasteryzacja i rozpoznanie pisma. Kolejność
odwrotna byłaby szybsza do napisania i kłamliwa w skutkach — dokument
z doskonałą warstwą tekstową wracałby jako odczyt z pikseli, a model bez
potrzeby przestałby ufać własnemu materiałowi. Skan bez warstwy tekstowej, na
maszynie bez Tesseracta, jest dokumentem, którego rdzeń nie umie przeczytać:
wraca odmowa nazywająca ten brak, nie pusty tekst z `usedOcr: false`, bo
pusty tekst znaczy, że dokument jest pusty, a to jest zdanie o dokumencie,
nie o rdzeniu.

Słownik `formatyDokumentu` zna dziewięć formatów i wymienia to, co rdzeń
umie zamieniać; odczyt jest czymś innym niż zamiana, więc materiał spoza tej
dziewiątki — arkusz, prezentacja, wiadomość poczty, plik biurowy — idzie do
Apache Tiki, biblioteki, której zadaniem jest rozpoznanie rodzaju pliku
i wydobycie z niego tekstu. Tika nie wypiera żadnej z istniejących dróg:
format znany słownikowi jedzie drogą znaną słownikowi, bo Pandoc i poppler
znają jego strukturę lepiej.

Apache Tika nie jest plikiem wykonywalnym — jest archiwum Javy, które
uruchamia maszyna wirtualna. Rdzeń rozdziela więc rozpoznanie programu od
rozpoznania archiwum, tak samo jak przy silniku mowy, gdzie osobno stoi
binarium programu, a osobno plik głosu. Programem jest `java`, i jego
dotyczy deklaracja narzędzia oraz sonda obecności — ścieżka wyszukiwania
systemu odpowiada na pytanie o niego wprost. Archiwum `tika-app-*.jar` wraz
z bibliotekami wydania programem nie jest, więc sonda obecności nie ma o co
je zapytać: jego brak jest osobną odmową z osobną naprawą, bo naprawa
„zainstalować Javę” niczego by tu nie załatwiła. Wersja archiwum nie jest
wpisana w kod, ponieważ numer wydania niesie nazwa pliku — numer wpisany na
stałe rozjechałby się z pierwszą aktualizacją Tiki i objawił odmową u
Operatora. Katalogi bibliotek dokładają się do ścieżki klas dlatego, że
wydania Tiki bywają dwojakie: archiwum samowystarczalne, niosące zależności
w sobie, albo archiwum cienkie obok katalogu `lib`. Rdzeń nie zgaduje, które
ma przed sobą — dokłada każdy katalog `lib`, jaki stoi w katalogu wydania,
a gdy nie stoi żaden, ścieżka klas zostaje samym archiwum i wydanie
samowystarczalne rusza tak samo.

## budowa/server/internal/core/przelot_baz_zdjeciowych_designu_test.go

Sprawdzian dowodzi, że dla każdego z ośmiu dostawców droga idzie do końca:
rdzeń składa adres, wysyła żądanie po HTTP, nosi klucz tam, gdzie dostawca go
żąda, rozkłada odpowiedź i wypełnia z niej wspólną postać zasobu. Serwer
próbny stoi w uprzęży i odpowiada kształtami, które rdzeń zakłada.

Sprawdzian nie dowodzi, że kształt odpowiedzi zgadza się z tym, co dostawca
naprawdę wysyła. Kształty pochodzą z dokumentacji, nie z pomiaru na jego API.
Ta połowa braku zostaje otwarta i jest tak nazwana — inaczej zielony wynik
tego pliku czytałoby się jako „dostawcy zmierzeni", a zmierzona jest droga.

Rdzeń nie zmienia się na potrzeby sprawdzianu: każda droga dostawcy
przyjmuje klienta HTTP jako argument, więc sprawdzian podstawia własnego —
z przekładnią, która przepisuje gospodarza adresu na serwer próbny
i zapisuje, o co rdzeń naprawdę poprosił. Rdzeń nie dostaje ani jednego pola
„adres na potrzeby sprawdzianu", bo pole takie żyłoby w produkcie i dałoby
się nim wskazać serwer obcy.

Defekt Smithsonian Open Access, na który jeden ze sprawdzianów odpowiada:
rdzeń szukał mediów pod ścieżką `content.descriptiveNonRepeating.media`,
a odpowiedź `api.si.edu` niesie je pod
`…descriptiveNonRepeating.online_media.media`. Wiersze przychodziły, żaden
nie dawał się złożyć w zasób, komenda oddawała wykaz pusty bez błędu — więc
dostawca nie wracał ani w wykazie, ani w bilansie dostawców nieudanych,
a odczyt sugerował, że fraza nie ma zdjęć. Sprawdzian mierzy obie strony:
kształt właściwy daje zasób, a wiersze bez mediów dają odmowę nazwaną wraz
z liczbą wierszy.

## skutek_fotografii_designu_test.go

Sprawdziany tego pliku mierzą skutek warsztatu fotografii i części drukarskiej
modułu Design w pikselach i w stronach, nie w kopercie odpowiedzi. Żaden
sprawdzian nie kończy się na stwierdzeniu, że odpowiedź jest udana: każdy
schodzi po odwołaniu zasobu do magazynu, rozkłada plik i pyta go o rzeczy,
których koperta nie zna — ile ma pikseli, czy niesie kanał krycia i jakie ma
w nim wartości, o ile przesunęła się średnia jasność, ile kafli powstało
i o jakich wymiarach, ile stron ma wydany plik PDF. Powód leży w historii tego
produktu: `design.asset.generate` meldował `status: ok` z wykazem zasobów,
za którymi nie było ani jednego bajtu, a sprawdzian zaglądający wyłącznie
w pole `status` uznawał to za powodzenie.

Liczenie stron pliku PDF w `liczbaStronPdfSprawdzianu` jest rachunkiem
własnym, celowo bez użycia biblioteki, którą rdzeń plik złożył: sprawdzian
liczący tą samą biblioteką mierzyłby zgodność biblioteki z samą sobą, a błąd
w składaniu pliku i błąd w jego odczycie zniosłyby się wzajemnie. Rozpakowanie
strumieni jest konieczne, bo `pdfcpu` zapisuje katalog obiektów w strumieniach
skompresowanych metodą Flate od wersji formatu PDF 1.5 wzwyż, więc napis
`/Type /Page` nie stoi wtedy w pliku wprost; pakiet `compress/zlib` biblioteki
standardowej wystarcza, bo to ta sama kompresja. Wzorzec bez ukośnika po
`Page` odróżnia stronę od drzewa stron (`/Type /Pages`), które w pliku
występuje raz.

## budowa/server/internal/core/adapter_kondycja.go

Zasada rodziny health.*: każda wartość pochodzi z pomiaru. Komenda zdrowia,
która oddaje stan sprawności nie zmierzywszy niczego, jest gorsza niż jej
brak — to fasada, przez którą awaria przechodzi niezauważona. W tym pliku
nie ma ani jednej ścieżki oddającej stan `up` bez wykonanego pomiaru: sonda,
której rdzeń nie ma czym wykonać, kończy się stanem `unknown` wraz z
powodem, nigdy stanem `up` dlatego, że nic nie zawiodło.

Dostępność liczy się z wierszy, nie z licznika: `health.uptime.get` nie
czyta żadnej kolumny „dostępność", tylko bierze serię pomiarów z zakresu
czasu i liczy udział wyników udanych. Licznik podnoszony przy zapisie
rozjechałby się z serią przy pierwszym usunięciu wyników albo zmianie
zakresu, a rozjazd byłby niewidoczny — obie liczby wyglądają tak samo.

Sam przebieg sondy, czyli co właściwie mierzy każdy jej rodzaj, leży
w `adapter_kondycja_pomiar.go`; port i wpięcie — w `handlers_kondycja.go`.

Trzy zależności ponad repozytorium — uruchamiacz procesów, rejestr kanałów
i izolacja wraz z katalogiem roboczym — obsługują trzy rodzaje pomiaru,
których rdzeń nie wykona sam: sondę `command` (uruchamiacz), sondę
`modelCall` (rejestr kanałów) oraz drogę wspólną z wszystkimi wołaniami
arsenału (izolacja). Brak którejkolwiek nie psuje montażu: psuje jeden
rodzaj sondy, która wtedy oddaje `unknown` wraz z powodem.

## budowa/server/internal/core/adapter_krok_zlecenia.go

Krokiem jest pozycja kolejki — ten sam byt, którym operuje silnik kolejki,
posuwany po tabeli przejść krokNaprzod. Działanie pause rodziny queue.action
wstrzymuje całą kolejkę, a nie krok: pozycja zostaje w stanie wykonywana,
a najbliższe resume przesuwa ją do do_weryfikacji, czyli traktuje pracę
przerwaną jak skończoną. Sterowanie pojedynczym krokiem prowadzi do innego
skutku — decyzja o kroku stojącym w wykonywana powoduje ponowne wykonanie
pracy z decyzją doklejoną do treści zlecenia, zamiast potraktować przerwaną
turę modelu jako gotową. Krok w stanie końcowym jest odmawiany na wejściu,
i przy wstrzymaniu, i przy decyzji, bo bez obu tych warunków wznowienie
zostawiłoby zlecenie w stanie pracy nad krokiem, który nie ma jak tego stanu
opuścić. Metody sterowania krokiem siedzą na adapterKolejek, ponieważ muszą
widzieć te same pozycje, ten sam silnik i tego samego wykonawcę, co
queue.action.

Wartość stanu zamkniety w wyliczeniu stanów widzianych przez Operatora jest
konieczna, bo bez niej krok ukończony, błędny albo anulowany przedstawiałby
się jako czekający na decyzję.

Zapisanie decyzji przed jej zastosowaniem — kolumny zdecydowano_o
i zastosowano_o — sprawia, że przerwanie rdzenia między jednym a drugim
zostawia decyzję widoczną i czekającą na dokończenie, a nie zgubioną.
Odrzucenie zamyka krok stanem anulowana z werdyktem odrzucone.

Stan sterowania kroku wystawiany w odpowiedzi na decyzję różni się celowo od
stanu, jaki pokazałby wykaz kroków odczytany chwilę później: liczy się
z wstrzymania czynnego, więc dopóki epizod nie jest jeszcze zastosowany,
znaczniki decyzji zostają w odpowiedzi, żeby Operator zobaczył, co się
właśnie stało; po zastosowaniu epizod jest już historią i krok wraca pod
stan swojej pracy.

## budowa/server/internal/core/skutek_automatyzacji_test.go

Żaden sprawdzian tego pliku nie kończy się na tym, że odpowiedź komendy jest
udana. Każdy schodzi do bazy własnym zapytaniem SQL i mierzy niezależnie, bo
odpowiedź komendy oddaje na przykład automatykę z etykietami i wyglądałaby
dokładnie tak samo, gdyby zapis do tabeli towarzyszącej nie doszedł. Wzorzec
szkody, którego pilnuje ten plik, ma w produkcie precedens: komenda meldowała
powodzenie z wykazem, za którym nie było ani jednego bajtu w bazie.

Sprawdzian skarbca i audytu mierzy dwie rzeczy naraz. Pierwsza: czy w bazie
nie ma wartości poświadczenia — kolumny na nią nie ma, więc sprawdzian
przeszukuje cały wiersz, i gdyby ktoś dołożył kolumnę i zapisał w niej sekret,
ten sprawdzian by upadł. Druga: czy dziennik audytu zapełnia się sam, przy
okazji czynności, bo audyt, który trzeba jawnie zawołać, jest audytem,
o którym się zapomina.

## skutek_drogi_neuronowej_designu_test.go

Sprawdziany tego pliku mierzą wywołanie cudzego silnika żądaniem wychodzącym do kanału
obrazowego, nie samą odpowiedzią rdzenia, bo odpowiedź poprawna kształtem potrafi
ukrywać żądanie, które nie niesie materiału. Przy każdej z czterech czynności —
powiększeniu, usunięciu tła, domalowaniu i rozszerzeniu kadru — mierzone są trzy rzeczy:
żądanie do kanału niesie zdjęcie Operatora, a przy domalowaniu i rozszerzeniu kadru
także maskę, bo wywołanie bez materiału każe silnikowi wygenerować obraz nowy zamiast
przetworzyć zdjęcie wniesione przez Operatora; plik wyniku ma wymiar, który czynność
obiecała, a nie ten, który akurat oddał kanał; pole `computedBy` niesie wartość
`kanalModelu`, a łańcuch edycji zapisuje tę samą drogę.

## budowa/server/internal/core/urzadzenia_druk.go

Warstwa druku lokalnego różni się od wydania plikiem: wydanie kończy pracę
plikiem gotowym do drukarni, z przestrzenią barw, spadami i znacznikami
cięcia, i kładzie go w magazynie jako zasób. Tego pliku nikt jednak nie
wydrukował na drukarce stojącej obok Operatora — to jest dziura, którą
zamyka ten plik: zasób wydany wcześniej, albo dowolny inny plik widziany
przez rdzeń, idzie tu na kolejkę druku systemu.

Kontrakt nie ma jeszcze komendy wysłania na drukarkę ani komendy wykazu
drukarek; rodzina komend druku ma nastawy profilu, kontrolę przeddrukową,
wydanie i podział wielkoformatowy, i na tym się kończy. Warstwa stoi tu
gotowa i czeka na dwie komendy, które kontrakt musi wnieść. Do tego czasu
funkcje tego pliku są wystawione poza pakiet, żeby adapter wołający mógł je
wziąć jedną linią w dniu, w którym komendy powstaną, bez przepisywania
warstwy.

Rozdzielenie drogi po systemie operacyjnym opiera się na sprawdzeniu
w czasie działania, nie na warunku budowy: budowa całego drzewa na jednym
systemie nie skompilowałaby gałęzi drugiego systemu ani razu, więc zepsułaby
się niezauważona aż do wydania instalki natywnej. Na Linuksie drogą jest
CUPS, tą samą, którą druku używa cały system; na Windowsie drogą jest
PowerShell, gdzie wykaz drukarek oddaje polecenie systemowe, a wysłanie idzie
przez .NET albo przez czasownik powłoki systemu, zależnie od rodzaju pliku:
obraz drukuje biblioteka .NET, tekst idzie przez bufor wydruku wprost, a
dokument złożony, którego rdzeń sam nie umie odczytać, potrzebuje programu
zarejestrowanego w systemie pod tym czasownikiem. Gdy taki program nie jest
zarejestrowany, odmowa mówi to wprost zamiast milczeć — plik wysłany
w nicość wygląda jak wydruk, który się nie pojawił.

## adapter_rozmowa.go

Osłona gorutyny tury w `prowadzTure` jest rejestrowana pierwsza, więc przy panice biegnie
ostatnia — dopiero po tym, jak strumień domknie się awaryjnie i bieg zostanie zapomniany;
klient dostaje więc znacznik końca, a wpis rozmowy wychodzi ze stanu `strumien`, zanim
panika przestanie lecieć. Bez tej osłony każda usterka adaptera w torze tury gasi cały
rdzeń, a nie jedną turę: panika w gorutynie nie ma kto przechwycić i proces ginie z sesją,
kolejką i połączeniem Operatora naraz — tak wywracał go pusty wskaźnik portu doraźnych
dołożeń, jedna wada jednego adaptera zabierała całą pracę. Osłona nie jest zgodą na
usterki i niczego nie ucisza: powód idzie do dziennika ze śladem stosu, żeby wada została
zgłoszona jako wada, a nie zniknęła w ciszy; miarą jest to, że Operator traci jedną
odpowiedź zamiast całej sesji.

## budowa/server/internal/core/przegladarka_silnik.go

Silnik przeglądarki uzupełnia pobranie realizowane w
`przegladarka_pobieranie.go` o wszystko, co jest własnością strony
URUCHOMIONEJ, nie jej źródła: zrzut ekranu, drzewo DOM po zbudowaniu przez
skrypty, komunikaty konsoli, rejestr żądań sieciowych, emulację urządzenia
i przewinięcie. Rozmowa z uruchomioną stroną idzie protokołem Chrome
DevTools.

Zasada produktu mówi, że żadna funkcja nie zależy od programu, którego
instalka nie niesie. Cała aplikacja z arsenałem stoi na serwerze, u Operatora
jest samo okno — Chromium jest więc programem serwerowym, tak samo jak
ffmpeg czy Tesseract, i jak one stoi w sondzie zależności zewnętrznych.
Silnika przeglądarki nie da się wkompilować w binarium Go; wyjątek na
biblioteki wkompilowane obejmuje PDF i kryptografię, nie renderowanie stron.

`zewnetrzne.Wolaj` prowadzi uruchomienie do końca i oddaje bajty po
zakończeniu programu. Przeglądarka ma żyć, dopóki trwa rozmowa: startuje,
przyjmuje polecenia protokołem i dopiero potem gaśnie, więc proces startuje
w tym pliku, a nie przez `Wolaj`. Sekwencja jest jednak ta sama co przy
innych programach zewnętrznych: port `session.Uruchamiacz`, brama izolacji
`session.SprawdzPolecenie`, objęcie drzewa procesów
`session.PrzejmijDrzewo` — Chromium rozgałęzia procesy renderowania i sieci,
a przerwana sesja bez objęcia drzewa zostawiłaby je na maszynie Operatora.

Strona, która nie kończy wczytywania, jest zjawiskiem codziennym. Każde
otwarcie ma granicę czasu; po jej przekroczeniu sesja oddaje to, co zdążyła
zebrać, albo odmawia — nigdy nie czeka bez końca.

## budowa/server/internal/core/skutek_zdolnosci_wyszukiwania_test.go

Część sprawdzianów żąda wag modeli na dysku i jest pomijana tam, gdzie wag
nie ma. Pominięcie jest tu jedyną uczciwą odpowiedzią: wagi ważą łącznie
blisko cztery gigabajty, więc sprawdzian, który by je pobierał, zamieniałby
bieg sprawdzianów w pobieranie modeli, a sprawdzian, który by ich nie
potrzebował, mierzyłby atrapę i milczałby dokładnie wtedy, gdy zdolność
przestanie działać.

Sprawdziany odmowy wag celowo nie żądają wag: brak silnika ma być
odpowiedzią nazywającą brak na każdej maszynie, więc mierzy się go tam,
gdzie modelu nie ma z samego założenia.

## budowa/server/internal/core/izolacja_test.go

Straż transportu sprawdza wyłącznie to, czy gniazdo w ogóle weszło. Straż
dostępu opisana tym plikiem pilnuje granicy drugiej i odrębnej: co wolno
oknu, które już weszło. Nadanie żyje per okno, niesie tryb dostępu i zawęża
się do korzeni punktu — pomyłka w którymkolwiek z tych trzech elementów
oddaje procesowi modelu katalog albo maszynę, których operator mu nie dał.

Zakres wyłączony przepuszcza wszystko celowo: pełny dostęp w ramach uprawnień
operatora jest stanem wyjściowym platformy. Z tego powodu każdy sprawdzian
w pliku włącza zakres wprost, zamiast liczyć na wartość domyślną — sprawdzian
zapomniany o tym mierzyłby ciszę zamiast granicy.

## skutek_dogniecenia_obrazu_test.go

Dogniecenie zapisu przy konwersji obrazu jest ulepszeniem, nie warunkiem: przy
braku programu zewnętrznego konwersja ma oddać ten sam obraz zapisany dłuższym
strumieniem. Sprawdzian mierzący zysk na rozmiarze mierzy więc obecność
programu tak samo jak jego pracę, dlatego na maszynie bez niego pomija się
z nazwanym powodem zamiast zawieść albo przejść bez zmierzenia niczego.
Sprawdzian TestKonwersjaUdajeSieNiezaleznieOdProgramuDogniatajacego mierzy
odwrotną połowę reguły i dlatego nigdy się nie pomija.

Obraz sprawdzianu dogniecenia niesie gradient o łagodnym przebiegu, bo koder
wkompilowany zapisuje go z zapasem miejsca do skrócenia; obraz jednolity
nie nadawałby się do pomiaru, ponieważ koder Go zapisuje go już blisko
granicy i zysk bywa zerowy, co czytałoby się jak brak dogniecenia.

Sprawdzian PNG i sprawdzian JPEG porównują zapis programu zewnętrznego
z zapisem kodera Go policzonym na miejscu, a nie ze stałą liczbą bajtów:
stała rozjechałaby się przy pierwszej zmianie biblioteki i zaczęłaby mierzyć
jej wydanie zamiast pracy programu.

Sprawdzian zapisu stratnego mierzy model barw wyniku, a nie samą liczbę
bajtów: plik palety niesie `color.Palette`, a zapis pełnobarwny — `NRGBA`.
Model palety dowodzi, że wynik programu został wzięty, bo koder wkompilowany
palety nie zapisuje. Obraz tego sprawdzianu niesie dwieście barw rozrzuconych
bez ładu, ponieważ zapis pełnobarwny nie ma czego przewidzieć i płaci trzy
bajty za punkt, a paleta mieści wszystkie barwy co do jednej; gradient
nadałby się gorzej, bo pngquant odmawia sprowadzenia płynnego przejścia do
palety, gdy nie mieści się w progu jakości, i sprawdzian mierzyłby wtedy
odmowę zamiast pracy programu.

Sprawdzian niezależności od programu wytwarza maszynę bez programów na
miejscu, pustą ścieżką wyszukiwania, więc reguła jest mierzona wszędzie,
a nie tylko tam, gdzie programów akurat nie zainstalowano.

## budowa/server/internal/core/wejscie_wniesienie_test.go

Plik mierzy skutek wejścia dokumentu Operatora do edytora: czy postać pliku
wchodzi bez utraty formy, czy zapis znaków jest rozpoznawany i czy kopia
dokumentu jest bytem osobnym od oryginału. Miara jest zawsze taka sama i nie
ocenia koperty odpowiedzi komendy — po czynności dokument czyta się ponownie
przez `PostacDokumentu`, bo Operator otworzy go ponownie, a nie przeczyta
odpowiedzi komendy wprost.

Sprawdziany wykluczają cztery rodzaje szkody. Pierwsza to plik wniesiony jako
treść płaska, z arkuszem stylów i tabelą zgubionymi po drodze, czyli postać
ginąca przy wniesieniu. Druga to plik w stronie kodowej innej niż UTF-8
wczytany jako ciąg nieczytelnych znaków, bez żadnego zdania o tym w bilansie
czynności. Trzecia to kopia dokumentu założona jako drugie odwołanie do tego
samego bytu, po której zmiana w kopii rusza oryginał. Czwarta to wniesienie
oddane jako udane, a bez zapisu pochodzenia dokumentu, po którym nie da się
odtworzyć, na czym pismo się opiera.

## budowa/server/internal/core/montaz.go

Kolejność montażu w `Zmontuj` jest wymuszona zależnościami, nie upodobaniem:
repozytoria dają źródła konfiguracji i kanałów, transport daje nadajnik
zdarzeń, a rdzeń powstaje na końcu, bo dopiero wtedy ma czym wypełnić porty.
Brak elementu opcjonalnego nie przerywa montażu: pusty rejestr kanałów, brak
profili kanału głównego i brak pakietu klienta zostawiają rdzeń zdolny do
pracy w pozostałym zakresie.

Rozpoznanie maszyny bieżącej zakłada wiersz urządzenia, na którym stoi rdzeń.
Nieudane rozpoznanie idzie do dziennika i nie przerywa montażu. Degradacja
katalogu roboczego do lokalizacji zastępczej także idzie do dziennika, tak by
Operator wiedział, że pracuje gdzie indziej, niż ustawił.

Przejmowanie procesów powstaje przed kanałami, bo kanał wkłada haczyk do
każdej tury, a wiązane jest po nadzorcy, bo to on ma rejestr procesów. Ta
jedna pośredniczka domyka różnicę kolejności. Nadzorca nie dostaje tu żadnego
wypełnienia: proces tury startuje kanał modelu własną drogą, a sesja obejmuje
go uchwytem przez Przejmij. Drugiej drogi startu procesu nie ma.

Sprzątanie startowe stanu trwałego idzie przed odtworzeniem rejestru: start
jest jedynym momentem, w którym rdzeń i tak czyta stan trwały. Opróżnienie
kosza sesji po terminie i przemiecenie retencji historii to dwie czynności
jednej drogi sprzątania; kolejność jest istotna, bo sesja skasowana z kosza
zabiera swoje okna wraz z wypowiedziami.

Telemetria postępu czyta szynę zdarzeń i port rozmowy, więc powstaje przed
rdzeniem i owija nadajnik transportu. Żywy stan sesji stoi pod telemetrią:
producent telemetrii rozgłasza `progress.changed` wprost tym nadajnikiem,
który dostał, więc nasłuch obecności musi być tym nadajnikiem. Tak domyka się
łańcuch producent → szyna → odbiorca kontrolki powrotu do sesji.

Doraźne dołożenia narzędzi mają dwóch czytelników: port `NarzedziaSesji`
rdzenia (rodzina `session.tool.*` i `tools.catalog.list`) oraz składacz
zestawu narzędzi tury, który dokłada je do wykazu eksperta
(`adapter_rozmowa_zestaw.go`). Instancja powstaje przed składaniem portów, bo
oba mają dostać ten sam adapter — drugi byłby drugą prawdą o tym, czym model
w sesji dysponuje.

Zakresy narzędzi mają dwóch czytelników: port `ZakresyNarzedzi` rdzenia
(rodzina `tools.scope.*`) oraz straż, którą rdzeń pyta przed skierowaniem
komendy ręki modelu. Jedna instancja na obie drogi — druga byłaby drugą
prawdą o tym, co profilowi wolno. Straż zakresu eksperta wpina się w dwa
miejsca odmowy: nałożenie eksperta na okno i powołanie podagentów
(`straz_eksperta.go`).

Wymóg logowania w `ustawieniaTransportu` ma dwa źródła i jedno pierwszeństwo.
Wskazanie ze startu (przełącznik wiersza poleceń, zmienna środowiska) wygrywa
zawsze. Dopiero jego brak oddaje głos nastawie poziomu `aplikacja` z tabeli
`ustawienie`, a brak i jej — adresowi nasłuchu. To ten sam rozstrzygacz i ta
sama tabela, którą widzi `config.get`.

## budowa/server/internal/core/adapter_narzedzia_dokument_konwersja.go

Do PDF-u prowadzą dwie drogi, bo żadna pojedyncza nie wystarcza. Pandoc
zamienia struktury tekstowe (markdown, html, docx, odt, rtf, epub, csv, tekst
czysty), ale PDF-u sam nie zapisze: wywołanie z docelowym formatem pdf woła
silnik składu jako własne potomstwo, czyli proces poza bramą rdzenia — rdzeń
woła więc silnik osobno. Dokument, który LibreOffice otwiera wprost (docx,
odt, rtf, html, csv, txt), idzie LibreOffice'em bez okna; ta droga niesie
własny układ dokumentu — style, tabele i podziały stron zapisane w pliku,
których żadne przepisanie przez format pośredni nie odtworzy. Materiał, którego
LibreOffice wprost nie otwiera (markdown, epub), idzie składem: Pandoc
zamienia go na źródło typsta, a typst składa PDF. Wcześniej ta droga szła
przez fragment HTML-a rysowany procesorem tekstu — droga, która działa, ale
składem nie jest; gdy typst nie stoi na maszynie, ta droga zostaje jako
zapasowa, bo odmowa byłaby regresem względem stanu, w którym PDF z markdownu
powstawał i bez typsta. Komenda nie czyta PDF-u jako źródła: jego treść jest
ciągiem instrukcji rysowania, z którego Pandoc nie złoży struktury dokumentu;
żądanie zamiany z formatu pdf kończy się odmową wskazującą, że treść PDF-u
wyciąga osobna komenda tekstowa, a jej wynik da się przekonwertować dalej.
Wyniku zastępczego komenda nie podstawia: każda droga bez bajtów kończy się
błędem, a pusty plik na wyjściu binarium też jest odmową, bo dokument
o zerowej długości wygląda w panelu jak dokument.

Silnik typst wchodzi w miejsce, w którym rdzeń nie miał żadnego: jest jednym
plikiem wykonywalnym bez własnego drzewa zasobów, więc mieści się w płaskim
układzie programów pomocniczych, którym jedzie pakowanie produktu. Silnik
TeX-owy tego układu nie przyjmuje, ponieważ jest drzewem formatów, czcionek
i ścieżki wyszukiwania, a nie pojedynczym plikiem.

Rodzina czcionek narzucana składowi typst nie jest ozdobą, tylko warunkiem
uruchomienia: szablon typsta, który wypuszcza Pandoc, podaje silnikowi
rodzinę czcionek ze zmiennej szablonu, a przy jej braku podaje wykaz pusty,
co typst odrzuca błędem pustej listy czcionek zapasowych, i PDF nie powstaje
wcale. Rodzina nierozpoznana nie jest odmową: typst wtedy zgłasza nieznaną
rodzinę, schodzi na własną czcionkę zastępczą i składa dokument dalej —
sprawdzone uruchomieniem na tej maszynie. Rodzina DejaVu Serif stoi w każdej
instalacji niosącej pakiet czcionek DejaVu i pokrywa komplet polskich znaków
diakrytycznych.

LibreOffice trzyma stan w swoim profilu, a drugie równoległe uruchomienie na
wspólnym profilu kończy się cichym zwarciem — jeden przebieg oddaje plik,
drugi nic. Dlatego każde wywołanie dostaje własny profil w katalogu roboczym
czynności, który znika razem z nim.

Dwa uruchomienia zamiast jednego w drodze składu typstem, choć Pandoc umie
zawołać silnik składu sam. Powód jest ten, dla którego w całym drzewie stoi
jedno bezpośrednie wywołanie procesu na krok: silnik zawołany przez Pandoca
byłby jego potomstwem, więc ominąłby sprawdzenie obecności, bramę izolacji
okna i własną granicę czasu, a jego odmowa dochodziłaby do Operatora zwinięta
w ogólny komunikat Pandoca, bez zdania, które powiedział sam silnik. Wołany
osobno — mówi sam za siebie.

Kolejność zapisu wyniku jest zamierzona: najpierw bajty w magazynie, potem
wiersz zasobu. Wiersz wskazujący odwołanie, za którym nic nie leży, byłby
dokumentem nie do otwarcia, a model zacytowałby go jako gotowy. Okno puste
w wyniku znaczy brak wiersza, a nie błąd — zgodnie z tym, jak rozstrzyga to
wspólny mechanizm odkładania wyniku dla całego arsenału.

## skutek_cyfryzacji_studia_test.go

Szkoda, którą ten plik ma wykluczyć, ma w produkcie postać znaną: odpowiedź
`ok` przy pustym wyniku. Rozpoznanie pisma jest na nią szczególnie podatne,
bo Tesseract kończy się powodzeniem także wtedy, gdy nie odczytał ani
jednego słowa, więc koperta udana nie mówi nic o tym, czy Operator dostał
tekst. Dlatego żaden sprawdzian tego pliku nie kończy się na odpowiedzi
udanej — każdy pyta o to, co zostało: czy pozycja jest w kolejce przy
kolejnym odczycie, czy tekst niesie słowa z obrazu, czy korekta zmieniła
treść przyjmowaną do edytora i czy dokument założony z cyfryzacji ma tę
treść po ponownym otwarciu.

Materiał obrazowy sprawdzianów powstaje na miejscu, a nie leży w drzewie
jako plik binarny: plik w repozytorium starzeje się bez śladu, a tu chodzi
o to, żeby tekst na obrazie i tekst oczekiwany pochodziły z jednego zapisu.

TestPrzyjeciePozycjiZakladaDokumentZTrescia mierzy przyjęcie osobnym
odczytem dokumentu, bo odpowiedź na przyjęcie mogłaby nieść treść, której
baza nie przyjęła. Materiałem sprawdzianu jest obraz, a nie plik tekstowy:
rozpoznanie pisma czyta piksele, więc plik tekstowy podany jako materiał
kończyłby się odmową i sprawdzian pomijałby się zawsze — sprawdzian, który
zawsze się pomija, niczego nie pilnuje. Autor wersji założonej z cyfryzacji
niesie wartość modelu, nie Operatora, ponieważ historia dokumentu ma
rozróżniać, kto wniósł którą wersję treści.

## budowa/server/internal/core/skutek_kontroli_wiernosci_test.go

Kontrola wierności polega na porównaniu dwóch tekstów, które muszą się różnić:
przekładu w języku docelowym i jego tłumaczenia zwrotnego na język źródłowy.
Szkoda, którą sprawdziany tego pliku wykluczają, odbierała kontroli właśnie
tę różnicę — wywołanie tłumaczenia zwrotnego przepisywało panel w to samo
miejsce, więc po przebiegu obie strony porównania stawały się jednym tekstem
i kontrola udawała się zawsze. Odpowiedź wyglądała na udaną, a wynik był
bezwartościowy.

Dlatego sprawdziany tego pliku po każdym przebiegu pytają o obie strony
osobno: co niesie tłumaczenie zwrotne i co dalej niesie panel. Drugi przebieg
z inną odpowiedzią modelu ma dać inny wynik niż pierwszy — wynik, który się
nie zmienia mimo zmiany wejścia, nie jest wynikiem pomiaru.

Model jest w tych sprawdzianach prawdziwym kanałem rdzenia wskazującym punkt
końcowy podniesiony na czas sprawdzianu w tym samym procesie, a nie zaślepką
portu. Program spoza maszyny nie jest do tego potrzebny; potrzebny jest
natomiast do syntezy mowy, którą woła zewnętrzny silnik, więc tej komendy
sprawdziany tego pliku nie dotykają.

## budowa/server/internal/core/handlers_narzedzia_sesji.go

Ekspert dostaje dobrany podzbiór narzędzi; komenda po ukośniku wstrzykuje
jedno narzędzie na żądanie, nie ruszając definicji eksperta ani biegu
rozmowy. Jest to jedyne miejsce, w którym zestaw narzędzi rośnie w trakcie
pracy, więc `session.tool.attach` i `session.tool.detach` nie są narzędziami
modelu: narzędzie, którym model dokłada sobie narzędzia, byłoby drugą
prawdą o tym, czym model dysponuje. Odczyt model ma — `session.tool.list`
i `tools.catalog.list`.

`rozglosDolozenieNarzedzia` niesie w zdarzeniu całą pozycję — nazwę pełną
ze źródłem, skróconą i opis — bo sama nazwa skrócona nic Operatorowi nie
mówi. Sprawcę bierze z kontekstu, tak samo jak pozostałe zdarzenia rdzenia:
bez niego dołożenie zrobione ręką asystenta wyglądałoby jak własne.

`rozglosZdjecieNarzedzia` niesie całą pozycję z tego samego powodu:
rozgłoszenie samego dołożenia zostawiłoby sąsiednie okno z wierszem
narzędzia, którego model już nie ma, a zestaw pokazany szerszym, niż jest,
kłamie tak samo jak poszerzony po cichu. Okno usuwa wiersz po nazwie
pełnej, a gdy dołożenia jeszcze nie widziało, ma z czego złożyć wpis
dziennika mówiący Operatorowi, co dokładnie zeszło.

`bladNosnikaNarzedziSesji` rozróżnia dwie przyczyny: brak wpiętej
trwałości (dołożenie nie doszłoby do bazy i zniknęłoby przy pierwszym
rozłączeniu klienta, choć ma przetrwać) oraz usterkę zapisu, po której
Operator powtórzy czynność i sprawdzi dziennik rdzenia.

`trafnoscPozycji` dopasowuje środkiem nazwy, nie od początku: przy
przedrostkach źródła każda pozycja zaczyna się tak samo, więc szukanie od
początku nazwy pełnej byłoby bezużyteczne. Stopnie są cztery, bo `ski` ma
trafić w `skill-creator` wyżej niż w `template-skill`, a oba wyżej niż
pozycję, która ma `ski` wyłącznie w opisie.

## budowa/server/internal/core/adapter_centrum_powiadomien.go

Centrum powiadomień jest jedynym mechanizmem powiadamiania platformy.
Zdarzenie wchodzi do rejestru wyłącznie od strony rdzenia, funkcją `Zglos` —
tylko rdzeń wie, że coś zaszło. Kontrakt niesie sam odczyt i zmianę stanu;
komendy zgłaszającej nie ma, inaczej klient wpisywałby do rejestru zdarzenia,
które nigdy nie zaszły.

O tym, czy klasa zdarzenia w ogóle wchodzi do rejestru i czy idzie dalej na
telefon albo listem, rozstrzygają nastawy sekcji „Powiadomienia" okna
Ustawień. Adapter ich nie powtarza i nie zna ani jednej wartości domyślnej —
czyta je tym samym rozstrzygaczem, którym idzie każde inne ustawienie
platformy.

Silnik `zdalne.Zglos` stał zbudowany i nieużywany; ten plik niesie jego
jedyne wywołanie. Zdarzenie klasy dopuszczonej do kanału mobilnego wchodzi do
rejestru centrum i zaraz potem do kolejki doręczeń. Niepowodzenie kolejki nie
cofa zapisu w rejestrze — zdarzenie zaszło niezależnie od tego, czy telefon
je odebrał, a rejestr centrum jest kanałem podstawowym.

## budowa/server/internal/core/adapter_narzedzia_archiwum_rozpakowanie.go

Warstwy obrony są trzy, bo żadna pojedyncza nie daje pewności. Pierwsza to
spis przed zapisem: zawartość archiwum jest oglądana spisem `7z l -slt`,
zanim poleci pierwszy bajt, a pozycjom bezwzględnym, pozycjom z członem `..`
i dowiązaniom komenda odmawia już na tym etapie. Druga to kwarantanna:
rozpakowanie idzie do katalogu świeżo założonego i pustego, a nie wprost do
katalogu Operatora, bo w katalogu bez ani jednego pliku nie ma czego
nadpisać; kwarantanna stoi na tym samym nośniku co cel, więc przeniesienie
gotowej treści jest przemianowaniem, a nie drugim kopiowaniem. Trzecia to
przejście po wyniku: po rozpakowaniu komenda obchodzi kwarantannę i sprawdza,
co powstało — przejście po drzewie nie idzie za dowiązaniami i rozpoznaje je
po typie wpisu, więc dowiązanie, które przeszłoby przez spis, zatrzymuje się
tutaj; ta warstwa mierzy też prawdziwy rozmiar wyniku, bo deklaracja
w nagłówku archiwum jest obietnicą jego twórcy. Dopiero po trzeciej warstwie
treść wchodzi do katalogu Operatora — odmowa na każdej z nich zostawia jego
drzewo nietknięte, bo do tej chwili nic w nim nie powstało.

Miejsca docelowe są dwa, bo zamiary są dwa. Kontrakt opisuje pole docelowej
ścieżki jako katalog docelowy w katalogu roboczym okna, którego brak kładzie
zawartość w magazynie: ze ścieżką model rozpakowuje po to, żeby na tych
plikach pracować, bez ścieżki — żeby zawartość przechować pod sumami
kontrolnymi, nie zaśmiecając katalogu Operatora. Plik istniejący nie jest
nadpisywany: przy przenoszeniu z kwarantanny nazwa zajęta w celu jest
odmową, bo nadpisanie cudzej pracy zawartością archiwum jest tą samą szkodą,
przed którą stoi reszta tego pliku.

Pole ścieżek wydanych Operatorowi zostaje puste, gdy zawartość ląduje
w magazynie: tam nie ma ścieżek, są zasoby pod identyfikatorami, a wypisanie
nazw z wnętrza archiwum dałoby napisy wyglądające jak ścieżki na dysku
Operatora, pod którymi nic nie leży. Pole jest w kontrakcie nieobowiązkowe
właśnie dlatego, że jedna z dwóch dróg tej komendy ścieżek nie wytwarza.

## budowa/server/internal/core/adapter_okno_przekazanie.go

Port `PrzekazanieOkna` i komenda `window.action` leżą w `adapter_okno_akcja.go`
i `adapter_okno_przekazanie_uchwyty.go`; ten plik deklaruje wyłącznie typ
adaptera, jego konstruktor i wiązania zależności.

Więź koordynator–wykonawca (kolumna `okno_komunikacji.okno_koordynatora_id`)
bez utrwalenia tutaj żyłaby wyłącznie w pamięci przeglądarki i ginęła z jej
zamknięciem, a Mission Control nie miałby czego pokazać po ponownym
uruchomieniu rdzenia. Komenda `window.handoff` robi trzy zapisy naraz —
wszystkie albo żadną: utrwala więź koordynator–wykonawca, zapisuje zlecenie
wraz z kompletem kontekstu i zakłada pozycję kolejki, której identyfikator
wraca w `QueueItemId`.

Warstwa danych (`dane/przekazanie_okna*.go`) przyjmuje `PozycjaKolejkiID` jako
identyfikator gotowy i sama pozycji nie zakłada. Zakładanie należy więc do
tego adaptera i idzie przez jedyny silnik kolejek — ten sam `adapterKolejek`,
którym pracuje pętla sesyjna i moduł Automations. Wzorem
`adapterAutomatyk.ZKolejkami` (`adapter_modul_automations.go`) konstruktor
bierze wyłącznie repozytorium obszaru, a silnik kolejek dochodzi osobnym
wiązaniem po złożeniu grafu zależności w `montaz_moduly.go`. Bez podpiętego
adaptera kolejek `Przekaz` odmawia wprost, zamiast meldować wykonanie.

Repozytoria okien, sesji, modułów i kanałów stoją tu, bo `window.handoff`
przyjmuje identyfikatory zewnętrzne (tekstowe) okien i sesji, a warstwa danych
obszaru window.* oraz silnik kolejek pracują na kluczach wewnętrznych
(`int64`) tabel `okno_komunikacji` i `sesja`. Rozwiązanie identyfikatora na
wiersz — i odmowa `not_found`, gdy okna nie ma — jest obowiązkiem tego
adaptera, bo warstwa danych przyjmuje klucze już rozwiązane. Moduły i kanały
modelu służą wyłącznie złożeniu odpowiedzi: kontrakt `Window` niesie kody
modułu i kanału, a wiersz okna niesie klucze obce do nich. Wszystkie te
zależności powstają w `montaz_moduly.go`, dlatego adapter przyjmuje je gotowe
przez wiązania tego pliku.

Schemat dopuszcza dokładnie dwa rodzaje kolejki — `sesyjna` i `multitasking` —
bo jeden silnik kolejek obsługuje pętlę sesyjną i MultitaskingAI. Przekazanie
zlecenia z okna koordynatora do okna wykonawcy jest pętlą MultitaskingAI,
więc mieści się w rodzaju już istniejącym; własny rodzaj byłby trzecim
znaczeniem tego samego pojęcia i padłby na warunku CHECK kolumny.

Kolejność trzech zapisów w `Przekaz` jest celowa: pozycja kolejki idzie
pierwsza, bo jest jedynym z trzech zapisów, który silnik kolejek umie cofnąć
samodzielnie — kolejka porzucona bez zlecenia jest stanem nieszkodliwym; więź
i zlecenie idą po niej, gdy wiadomo już, że jest czym wykonać. Wspólnej
transakcji SQL między repozytoriami nie ma, bo warstwa danych obszaru
window.* i silnik kolejek to dwa oddzielne repozytoria. Odmowa wczesna, przed
pierwszym zapisem, pokrywa najczęstszy przypadek: okno albo sesja, których
nie ma. Usterka po pierwszym zapisie zostawia założoną pozycję kolejki bez
zlecenia, co Queue Manager pokazuje jako pozycję do ręcznego domknięcia,
a nie jako ciche zaginięcie zlecenia.

Kolejka w `zalozPozycjeKolejki` idzie przez `a.kolejki.repozytorium` — ten sam
obiekt, którym jedzie `adapter_modul_automations_kolejka.go`; wzorzec
`queue.create` w `adapter_kolejki.go` też zakłada kolejkę za każdym
wywołaniem.

`oknoWynikuKontraktu` korzysta z tego samego przekładu co odczyt utrwalony —
`oknoWierszaKontraktu` w `przeklad_nawigacja.go` — bo identyfikator sesji jest
już znany z żądania i drugi odczyt sesji byłby zapytaniem po to samo.

Wykaz kanałów w `kodSlownika` czyta się bez filtra aktywności, bo okno mogło
zostać założone na kanale, który od tamtej pory wyłączono.

## brama_kontraktu.go

Brama stoi w jednym miejscu, a nie w każdym obsługiwaczu z osobna, ponieważ
kontrakt niesie dla każdej komendy wykaz pól wraz z oznaczeniem „wymagane”,
a dla pól o typie wyliczenia — komplet dopuszczalnych wartości. Rdzeń bez
bramy przyjmował żądanie niepełne i wartość spoza wyliczenia, po czym
uzupełniał brak wartością domyślną i meldował powodzenie; odpowiedź
wyglądająca dobrze jest w tym miejscu groźniejsza od odmowy, bo nie wzywa
wołającego do sprawdzenia — okno dostawało punkt dostępu rodzaju, o który
nie prosiło, i tryb uprawnień, którego nie ustawiło. Brama nie ma własnego
wykazu pól ani własnego wykazu wartości: oba czyta z artefaktu kontraktu,
więc pole dołożone do kontraktu jest pilnowane od razu, bez zmiany w tym
pliku.

Powitanie kanału jest jedyną komendą spod bramy wyjętą i odpowiada zawsze,
także na żądanie niepełne. Powitanie jest jedynym miejscem, w którym klient
odczytuje wersję protokołu rdzenia, czyli jedynym, w którym rozpoznaje, że
jest starszy; brama sprawdzająca powitanie wobec kontraktu zakładałaby, że
obie strony znają już ten sam kontrakt — zakładałaby więc to, co powitanie
ma dopiero ustalić, i klientowi sprzed wprowadzenia pola oddawałaby odmowę
zamiast wersji, po której ten rozpoznałby rozjazd. Braki pól powitania idą
do dziennika rdzenia, nie do treści odpowiedzi: odpowiedź powitania nie ma
pola, w którym mogłyby wrócić wołającemu, a dołożenie takiego pola jest
zmianą kontraktu. Wyjątek jest jeden i pozostaje jeden: wynika z roli
powitania w uzgodnieniu, nie z wygody, więc każda inna komenda przechodzi
bramę bez ustępstw.

Obecność pól wymaganych sprawdzana jest dla każdej komendy, bo oznaczenie
„wymagane” niesie znacznik struktury żądania, a struktury ma każda komenda.
Wartości wyliczeń sprawdzane są węziej: komplet dopuszczalnych wartości
stoi w artefakcie Go wyłącznie przy polach komend wystawionych jako
narzędzia modelu. Wyliczenie samo w sobie ma w artefakcie funkcję
odczytującą jego wartości, ale nie ma odwzorowania typu pola na tę funkcję,
więc pole wyliczeniowe komendy spoza tego zbioru przechodzi bez sprawdzenia
wartości. Pełne pokrycie wymaga tabeli wyprowadzonej z kontraktu przy
generowaniu artefaktu, a nie przepisanej tutaj — drugi wykaz wartości
rozjechałby się z kontraktem przy pierwszej dołożonej wartości.

Sprawdzenie pól wymaganych bada obecność klucza w treści, nie jego
zawartość: kontrakt mówi „pole ma być”, nie „pole ma być niepuste”. Treść
pusta bywa treścią prawdziwą — zapis dokumentu z pustą treścią zapisuje
dokument opróżniony i jest żądaniem poprawnym; wartość pustą, tam gdzie
dziedzina jej nie zniesie, odrzuca obsługiwacz komendy, bo tylko on wie,
czy pustka coś znaczy. Wartość `null` przy polu wymaganym o typie tablicy
przechodzi z tego samego powodu: pole wymagane o typie tablicy wychodzi
z niepustego wykazu pustego jako `null`, tak koduje pusty wycinek biblioteka
standardowa Go, więc odmowa w tym miejscu odrzucałaby żądania składane
przez sam rdzeń.

## budowa/server/internal/core/adapter_narzedzia_archiwum.go

Formaty, ścieżki i odmowy leżą w `adapter_narzedzia_archiwum_sciezki.go`,
pakowanie w `adapter_narzedzia_archiwum_pakowanie.go`, rozpakowanie wraz
z obroną przed ucieczką ze ścieżki w
`adapter_narzedzia_archiwum_rozpakowanie.go`, czytanie spisu archiwum
w `adapter_narzedzia_archiwum_spis.go`, port i wpięcie
w `handlers_narzedzia_archiwum.go`.

Wszystkie formaty obsługuje jedno binarium wołane przez `zewnetrzne.Wolaj` —
przez port `session.Uruchamiacz`, bramę izolacji okna i objęcie drzewa
procesów; własnego `exec.Command` w tych plikach nie ma. Spis archiwum
czyta ten sam program, który potem rozpakowuje, więc wyrok o zawartości
nie rozjeżdża się z rozpakowaniem — osobne `unzip` i `tar` dawałyby trzy
postacie spisu i trzy okazje do obejścia tej kontroli.

Archiwum wytworzone przez `archive.pack` ląduje w tym samym magazynie, co
`design.asset.upload` (`magazynZasobowDesignu`, blob pod sumą sha256), a jego
wiersz — w tabeli `zasob_design`. Drugi magazyn byłby drugą prawdą o tym,
gdzie rdzeń trzyma bajty poza bazą.

`granicaNarzedziArchiwum` istnieje, bo `zewnetrzne.Wolaj` granicy
niedodatniej nie przyjmuje.

Sto tysięcy pozycji graniczne w `granicaRozpakowaniaPozycji` leży
kilkakrotnie powyżej liczby plików całego drzewa tego produktu wraz
z zależnościami. Dwa gibibajty graniczne w `granicaRozpakowaniaBajty` leżą
kilkanaście razy powyżej największego archiwum, jakie ten produkt ma do
wydania, i wiele rzędów wielkości poniżej znanych bomb rozwijających się
do terabajtów; rozpakowanie idzie przez kwarantannę, więc treść leży
przez chwilę w dwóch egzemplarzach, a granica musi zmieścić się na
nośniku dwukrotnie.

`nowyAdapterNarzedziArchiwum` wiąże port z repozytorium modułu Design
i magazynem jego zasobów. Uruchamiacz oraz izolacja wchodzą osobno, przez
`ZArsenalem`, bo montaż zna je dopiero po złożeniu warstwy kanału. Brak tej
zależności nie psuje montażu — obie komendy archiwum odmawiają wtedy,
nazywając brak, zamiast udawać, że archiwum powstało.

`katalogRoboczyOkna` stoi na straży obrony ścieżek: `targetPath` przyjęty
jako ścieżka bezwzględna pozwalałby modelowi rozsypać zawartość archiwum
w dowolnym miejscu maszyny — w `~/.ssh`, w katalogu autostartu, w cudzym
projekcie.

Wyliczenie `DesignAssetKind` kontraktu ma trzy wartości — `image`, `vector`,
`composition` — i ani jednej na archiwum; rodzajem odłożonego zasobu jest
`image`, tak samo jak w rodzinie mediów odkładającej nim dźwięk i film,
bo wartość spoza wyliczenia postawiłaby przed klientem napis, którego jego
typ nie zna. O tym, czym plik jest, mówi zmierzone pole `format`.

## budowa/server/internal/core/aparat_skutek_test.go

Plik mierzy skutek aparatu dokumentu: czy spis treści zgadza się z nagłówkami,
czy przypis przenumerowuje się po wstawieniu przypisu przed nim i czy
znacznik nieświeżości mówi prawdę.

Sprawdziany wykluczają cztery rodzaje szkody. Pierwsza to spis treści oddany
jako odświeżony, a niosący nagłówki sprzed zmiany. Druga to numeracja
przypisów nadawana w kolejności zapisu, po której przypis wstawiony w środek
dokumentu kłamie do końca życia pisma. Trzecia to znacznik nieświeżości
trzymany w drzewie postaci — drzewo zapisuje się bez aparatu, więc znacznik
ginąłby przy pierwszym zapisie. Czwarta to odwołanie do elementu, którego
dokument nie ma, przyjęte jako założone.

## budowa/server/internal/core/tabele_skutek_test.go

Sprawdziany tego pliku wykluczają cztery rodzaje szkody na tabelach dokumentu.
Pierwsza to szerokości kolumn zerowe po scaleniu: tabela wygląda na złożoną,
a w wydaniu ma kolumny niewidzialne, czego wprost zakazuje wymaganie
zlecenia. Druga to sortowanie, które rozrywa wiersze albo przestawia wiersz
nagłówkowy. Trzecia to zamiana tekstu na tabelę, która zostawia tekst
w treści dokumentu i daje dokument niosący tę samą informację dwa razy.
Czwarta to usunięcie ostatniego wiersza albo kolumny, po którym zostaje
tabela o zerowej siatce. Miara jest brana osobnym wywołaniem wykazu tabel,
a nie z odpowiedzi czynności, ponieważ Operator otworzy dokument ponownie,
a nie przeczyta odpowiedź komendy.

## budowa/server/internal/core/adapter_narzedzia_archiwum_pakowanie.go

`sourcePath` rozstrzygamy względem katalogu okna. Pakowanie wygląda na
czynność czytającą, więc ścieżka bezwzględna zdawałaby się nieszkodliwa —
nie jest. Wynikiem `archive.pack` jest zasób w magazynie rdzenia, a zasób
model potrafi odczytać i rozpakować; przyjęcie ścieżki bezwzględnej dałoby
więc drogę: spakuj `~/.ssh`, odłóż jako zasób, rozpakuj u siebie — czyli
wyniesienie dowolnego pliku maszyny Operatora do materiału, którym model
dysponuje. Katalog roboczy okna jest jedynym miejscem, o którym produkt
umówił się z Operatorem, że model tam sięga, i pakowanie tej umowy nie
łamie.

Zasoby idą przez rusztowanie, a nie wprost do `7z`, bo blob magazynu
nazywa się swoją sumą kontrolną: podanie blobów `7z` wprost dałoby
archiwum, w którym pliki nazywają się ciągiem szesnastkowym, a Operator
dostałby nazwy, z których żadna nic nie znaczy. Rusztowanie jest
katalogiem tymczasowym, w którym każdy blob dostaje swoją nazwę czytelną
z wiersza zasobu; dopiero ono trafia do `7z`.

W `Spakuj` liczbę `entries` liczymy z tego, co w archiwum naprawdę jest,
a nie z tego, ile pozycji kazano spakować — liczba wzięta z żądania
byłaby w odpowiedzi nieodróżnialna od policzonej, a rozjechałaby się
z prawdą przy każdym pliku, którego `7z` nie wziął.

W `zrodloPakowania` zasoby mają pierwszeństwo przed ścieżką, tak samo
i z tego samego powodu co w rodzinie obrazu: treść zasobu leży już pod
sumą kontrolną i nie zmieni się między wskazaniem a odczytem, a plik na
dysku jest treścią żywą. Żądanie bez jednego i drugiego jest odmową —
podstawienie katalogu roboczego, skoro model nic nie podał, byłoby
zgadywaniem przedmiotu czynności, a spakowanie całego katalogu zamiast
dwóch plików jest pomyłką kosztowną.

W `rusztowanieZasobow` kopiujemy, a nie dowiązujemy: dowiązanie do bloba
oszczędziłoby bajty, ale `7z` zapisałby wtedy do archiwum dowiązanie albo
poszedł za nim zależnie od przełącznika, a przy rozpakowaniu ta sama
rodzina dowiązania odmawia — archiwum, którego własny produkt nie umie
otworzyć, byłoby rozjazdem wewnątrz jednej rodziny.

W `nazwaZasobuWArchiwum` nazwa zasobu pochodzi od Operatora i mogłaby
nieść ścieżkę, a archiwum z pozycją `../coś` odmówiłaby własna rodzina
przy rozpakowaniu; dwa zasoby o tej samej nazwie nadpisałyby się nawzajem
w rusztowaniu i jeden zniknąłby bez słowa, stąd numerowanie kolizji.

W `zbudujArchiwum` dla `tar.gz` są dwa różne pliki i dwa wywołania: `tar`
jest archiwum dwuwarstwowym, `tar` niesie pozycje, `gzip` opakowuje
całość jednym strumieniem, więc `7z l` na gotowym `.tar.gz` pokazuje
jedną pozycję — zawinięty plik `.tar` — a nie zawartość. Spis czytamy
zatem z warstwy `tar`, zanim ją zawiniemy. Dwa wywołania zamiast potoku,
bo `zewnetrzne.Wolaj` prowadzi jeden proces — i dobrze, bo potok dwóch
programów miałby dwa kody wyjścia i jedną odpowiedź o powodzeniu.

W `wolajPakowanie` przełączniki są wybrane, nie przepisane: `-bd` zdejmuje
pasek postępu, bo w strumieniu bez terminala byłby śmieciem w
diagnostyce, `-y` odpowiada twierdząco na pytania programu, bo proces bez
terminala nie ma komu odpowiadać i czekałby do granicy czasu, a `--`
zamyka listę przełączników — bez niego plik nazwany `-sdel` zostałby
wzięty za polecenie, a nie za nazwę.

## handlers_historia.go

Plik wpina trzy komendy, nie cztery: zdarzenie zmiany historii stoi w dziale
zdarzeń kontraktu i niesie ładunek, nie parę żądanie-wynik, więc do rejestru
komend nie wchodzi. Historia nie jest nowym bytem — pozycją historii jest
wypowiedź okna, wiersz istniejącej tabeli wiadomości; rodzina komend historii
daje jej drugi widok, od najnowszej, kursorem czasu, z możliwością
skasowania, a nie drugą tabelę, i wypełnia ją to samo repozytorium co widok
zwykły. Nowa tabela wchodziłaby wyłącznie pod zasadę przechowywania, bo tego
bytu w schemacie nie było wcale.

Retencja jest egzekwowana, nie tylko zapisana: zasada zapisana, której nikt
nie stosuje, byłaby atrapą — Operator widziałby nastawę i rosnącą historię
naraz. Egzekucja ma dwa miejsca wpięcia, oba na drodze, którą rdzeń i tak
przechodzi: zapis zasady egzekwuje ją zaraz po zapisaniu na wszystkich
oknach zakresu, więc nastawa działa od chwili ustawienia, a nie od
następnego restartu; wczytanie historii egzekwuje ją przed oddaniem wykazu
na czytanym oknie, więc wykaz nigdy nie pokazuje pozycji, która zasadzie już
nie podlega. Czego te dwa miejsca nie dają: okno, którego nikt nie czyta
i na którym nikt nie przestawia zasady, nie zostanie przycięte samo,
ponieważ przemiatania okresowego ani przy starcie rdzenia tu nie ma —
wpięcie takiego przemiatania siedzi poza tym pakietem.

Przy czyszczeniu wielu pozycji zdarzenie zmiany nie niesie pojedynczej
pozycji: kontrakt daje to pole jako niewymagane właśnie na ten przypadek,
bo po skasowaniu setki wypowiedzi nie ma jednej pozycji „po zmianie",
a wysyłanie stu zdarzeń zamieniłoby odświeżenie wykazu w burzę. Idzie jedno
zdarzenie z identyfikatorem okna i rodzajem zmiany „usunięto" — okno
przeładowuje wykaz w całości.

Wczytanie wykazu nie przerywa się, gdy egzekucja zasady zawiedzie: Operator
ma dostać historię, jaka jest, zamiast odmowy odczytu, ponieważ
niepowodzenie egzekucji zabiera tylko tę jedną czynność, nie całą komendę.
Pole rozmiaru całej historii w odpowiedzi wykazu panel podstawia pod
ostrzeżenie „ile zniknie" przy czyszczeniu, a kasowanie bez wskazania
pozycji usuwa okno w całości — liczba zawężona kursorem malałaby przy
każdym dociągnięciu strony i obiecywałaby utratę mniejszą niż rzeczywistą.

Zapis zasady przechowywania na zakres, którego byt nie istnieje w bazie,
jest nastawą bez skutku: egzekucja na zakresie okna nie pyta bazy o
istnienie okna, więc zasada zapisana na nieistniejące okno albo sesję
leżałaby w tabeli i nie przycięłaby nigdy niczego, a Operator miałby na nią
potwierdzenie powodzenia. Odmowa w tym miejscu nazywa brak bytu w bazie,
a nie zakaz zapisu zasady.

## budowa/server/internal/core/adapter_narzedzia_media.go

Dwie komendy modelu, `media.inspect` i `media.transcode`, stoją na tym samym
źródle bajtów, tym samym zasięgu izolacji, tym samym wołaniu binarium i tym
samym magazynie wyniku. Pomiar (`ffprobe`) leży w
`adapter_narzedzia_media_pomiar.go`, przetworzenie (`ffmpeg`) w
`adapter_narzedzia_media_przetworzenie.go`, wpięcie komend w
`handlers_narzedzia_media.go`. Film i dźwięk są materiałem, którego model nie
zmierzy ani nie przetworzy bez cudzego programu; kontrakt wpisuje obie komendy
jako narzędzia modelu (`danaco_media_inspect`, `danaco_media_transcode`), więc
ich wołaczem jest model w turze, a nie panel okna.

Binarium wołane jest jedną drogą — `zewnetrzne.Wolaj`. Własnego
`exec.Command` w tym pliku nie ma i być nie może: tamta droga idzie przez
port `session.Uruchamiacz`, bramę izolacji okna i objęcie drzewa procesów.
Ostatnie jest tu ważniejsze niż gdziekolwiek: `ffmpeg` rozgałęzia wątki
dekodera i filtrów, a przerwane transkodowanie bez objęcia drzewa zostawia na
maszynie operatora procesy mielące film w nieskończoność. `ffmpeg` pisze do
pliku w katalogu tymczasowym, a stamtąd bajty wciąga magazyn zasobów pod sumę
sha256 — ten sam magazyn, którym jedzie `design.asset.upload`
(`adapter_modul_design_wgranie.go`). Przetworzenie „w miejscu" byłoby
zniszczeniem materiału operatora przy pierwszej pomyłce w parametrach.

Rodzaj zasobu wyprowadzony przez `odlozWynikMediow` nie jest przybliżeniem.
Kontrakt zna `video`, `audio`, `document` i `archive`, warunek CHECK kolumny
je dopuszcza (`migracja_113_rodzaje_zasobow_arsenalu.sql`), a rodzaj
wyprowadza z formatu wyniku wspólna tablica arsenału. Ma to znaczenie właśnie
w tej rodzinie: `media.transcode` z czynnością `frame` daje obraz,
a `extractAudio` — dźwięk, choć komenda jest ta sama.

`zasiegNarzedziMediow` rozstrzyga pusty `konfig.Kontekst{}` jako poprawny
adres najszerszego z poziomów zasięgu, a nie podstawienie pustych struktur po
cichu — gdy operator włączy punkt izolacji globalnie, brama zadziała tu tak
samo jak dla Terminala. Powód jest ten sam co przy silniku mowy
(`adapter_modul_mowa.go`).

`bladNarzedziMediow` rozróżnia trzy przypadki. Brak narzędzia
(`*zewnetrzne.BrakNarzedzia`) trafia na `channel_unavailable`, bo jest
brakiem po stronie instalacji, który operator usuwa jedną komendą pakietu —
nie wadą żądania i nie usterką rdzenia; kod jest ponawialny, po instalacji
`ffmpeg` to samo żądanie przechodzi bez zmiany. Naruszenie punktu izolacji
(`session.ErrIzolacja`) trafia na `permission_denied`, tak samo jak przy
Terminalu i silniku mowy. Pozostałe usterki — granica czasu, wywrócenie
binarium, brak uruchamiacza — trafiają na `internal_error`, bo treść niesie
już to, co program powiedział o sobie sam.

## budowa/server/internal/core/skutek_pracy_z_dokumentem_test.go

Plik mierzy skutek okna pracy z dokumentem: czy dwie zmiany rdzenia
rzeczywiście coś robią. Sprawdziany wykluczają trzy rodzaje szkody. Pierwsza
to przyjęcie wskazanych fragmentów propozycji, które podmienia treść całą,
podczas gdy Operator wybrał tylko fragment drugi. Druga to zakres zmiany
śledzonej liczony w bajtach, podczas gdy dokument polski niesie litery
dwubajtowe, które taki rachunek rozcinałby w środku. Trzecia to nastawy
suwaków i polecenie Operatora, które nie dojeżdżają do modelu, przez co
suwak przestawiałby pole bez skutku.

Operacja kontekstowa jest mierzona od końca do końca, jednym przebiegiem, nie
po częściach: złożenie polecenia dla modelu i rachunek zakresu mogą być
poprawne osobno, a droga między nimi mimo to zerwana — wynik modelu może nie
dojść do treści, zmiana śledzona może się nie odłożyć, wersja może nie
powstać. Kontrakt obiecuje, że po operacji kontekstowej w dokumencie stoi
zmiana oznaczona autorstwem modelu, i to jest twierdzenie o skutku, nie
o samym rachunku. Miarę umożliwia kanał `echo`: odsyła treść zapytania i nie
sięga do sieci ani do żadnego programu zewnętrznego, jest wkompilowany
w rdzeń. Wynik operacji jest więc znany z góry — jest nim polecenie złożone
przez `trescOperacjiStudia`, a w nim wiersz „Czynnosc: <pozycja rejestru>",
który mógł przyjść wyłącznie od modelu.

Sprawdzian porównania bez wskazania stron pilnuje granicy, na której
`studio.diff.compare` meldował powodzenie kopertą pustą. Koperta pusta ze
stanem `ok` mówi oknu, że porównanie przebiegło i nie ma czego pokazać,
a rdzeń nie porównał niczego: fragmenty różnicy potrzebują dwóch stron,
a wzorzec potrzebuje strony, po której ma szukać. Odmowa nazywająca brakujące
pole jest jedyną odpowiedzią prawdziwą — po pustej kopercie okno nie ma jak
odróżnić zgodność wersji od braku wskazania, co z czym porównać.

## budowa/server/internal/core/adapter_narzedzia_obraz_model.go

`image.upscale` i `image.background.remove` są osobną rodziną od czynności
`image.*`, bo obie zmyślają szczegół, którego w źródle nie ma — piksele
między pikselami przy powiększeniu, granicę obiektu i tła przy wycięciu —
i obie potrzebują do tego sieci neuronowej z wagami na dysku. Wspólny
z rodziną `image.*` zostaje mechanizm: rozwiązanie pary wskazania zasobu
albo ścieżki na plik do odczytu, zasięg izolacji dla wołania binarium
i odłożenie bajtów wyniku w magazynie zasobów Designu; ten mechanizm rodzina
modelu bierze przez zaplecze wspólne, a nie kopiuje go, żeby druga kopia
reguły nie rozjechała się z pierwszą.

Bez silnika na maszynie ta rodzina odmawia, nazywając brak, zamiast
podstawiać przybliżenie: proste rozciągnięcie obrazu oddałoby wynik dwa razy
większy i ani o szczegół bogatszy, a brak silnika wycinającego tło nie może
oddać obrazu bez zmian podanego jako wycięty.

Wagi modelu leżą obok silnika i nie ściągają się w trakcie żądania. Silnik
wycinający tło przy pierwszym uruchomieniu potrafi pobrać model sam, a wtedy
około sto siedemdziesiąt sześć megabajtów wchodzi w czas jednego żądania,
którego granica jest liczona na przetwarzanie, nie na łącze; przy wolnym
łączu albo braku sieci żądanie urywa się w połowie pobierania i wygląda jak
usterka silnika. Wagi rdzeń sprawdza więc przed uruchomieniem, a przy ich
braku odmawia, podając, gdzie mają leżeć i ile ważą.

Ścieżka wag zależy od systemu, bo wagi są składnikiem pakietu, a pakiet jest
inny na serwerze i inny w wersji natywnej Windows: na Linuksie arsenał
serwera stawia je pod stałą ścieżką systemową, w wersji natywnej Windows
jadą obok rdzenia w katalogu programów pomocniczych, więc ścieżkę
bezwzględną rdzeń liczy dopiero w czasie pracy, względem pliku
wykonywalnego. Ścieżka zaszyta po linuksowemu odmawiałaby na Windowsie,
zanim doszłoby do wołania silnika, choćby wagi były w paczce.

Pliki pośrednie w katalogu przebiegu są konieczne, choć rodzina `image.*`
bierze wynik ImageMagicka ze standardowego wyjścia: żaden z dwóch silników
modelu nie umie pisać obrazu na wyjście, oba żądają ścieżki wyniku. Katalog
własny na przebieg nie miesza równoległych żądań i znika jednym usunięciem
niezależnie od tego, czy silnik się udał. Źródło wchodzi do katalogu
przebiegu dowiązaniem, nie kopią: blob zasobu leży pod swoją sumą kontrolną
i bywa wielkim plikiem, więc kopiowanie go tylko po to, żeby zmienić nazwę,
dawałoby drugi egzemplarz zdjęcia przy każdym żądaniu. Oba silniki
rozpoznają format po zawartości, a nie po rozszerzeniu, więc nazwa pliku
wejściowego jest tylko uchwytem, nie deklaracją formatu. Kopia zamiast
dowiązania jest drogą zapasową na wypadek, gdy dowiązanie się nie uda —
magazyn na innym nośniku niż katalog tymczasowy albo system plików bez
dowiązań.

## budowa/server/internal/core/adapter_narzedzia_archiwum_spis.go

Ścieżki wychodzące poza katalog docelowy: spis sprawdza się przed
rozpakowaniem, a nie ścieżki po nim, i to jest wybór, nie skrót.
Sprawdzenie po rozpakowaniu przychodzi za późno z definicji: pozycja
`../../.ssh/authorized_keys` jest już wtedy zapisana, a wykrycie
nadpisania nie odkręca nadpisania — odmowa ma paść, zanim poleci pierwszy
bajt. Spis czyta to samo binarium, które będzie rozpakowywać, i w tym
samym przebiegu wywołań, więc oglądane jest dokładnie to, co `7z` z tego
archiwum odczytuje, a nie własne wyobrażenie o formacie zip; własny
czytnik nagłówków zip w Go byłby drugą prawdą o zawartości archiwum,
a rozjazd między czytnikiem sprawdzającym a rozpakowującym jest
klasycznym sposobem obejścia takiej kontroli.

`7z x` przy rozpakowaniu sam obcina człon `..` i zapisuje pozycję
wewnątrz katalogu docelowego. Ta obrona binarium jest prawdziwa, ale nie
jest tą, na której stoi rdzeń, z dwóch powodów: jest cudzą własnością
i cudzą wersją — archiwum wychodzące poza katalog ma zostać odmówione na
każdej maszynie i przy każdym `7z`, a nie rozpakowane inaczej, niż
zapowiada; a ciche obcięcie członu jest zmianą znaczenia bez powiedzenia
o tym — Operator dostałby drzewo inne niż to, które archiwum opisuje,
i nie dowiedziałby się o tym. Nazwanie braku jest uczciwsze niż
naprawienie archiwum po cichu.

Pewność bierze się więc z trzech warstw naraz: spis odmawia przed
zapisem, rozpakowanie idzie do kwarantanny (pustego katalogu, w którym
nie ma czego nadpisać), a przejście po wyniku sprawdza, co naprawdę
powstało — dopiero potem treść wchodzi do katalogu Operatora.

Dowiązania są odmawiane, nie rozpakowywane: pozycja będąca dowiązaniem
symbolicznym wychodzi poza katalog docelowy inaczej niż członem `..` —
sama w sobie jest niewinna, ale wskazuje na zewnątrz, a kolejna pozycja
archiwum zapisuje przez nią, i zapis ląduje tam, gdzie wskazuje
dowiązanie. `7z l -slt` pokazuje na warstwie `tar` pole `Symbolic Link`,
więc dowiązanie widać w spisie bez zgadywania.

Archiwum puste w odczycie jest odmową, a nie pustym wynikiem: `7z`, który
nie umiał odczytać spisu (archiwum uszkodzone, zaszyfrowany nagłówek),
zostawia rdzeń bez wiedzy o zawartości, a rozpakowanie czegoś
nieobejrzanego omijałoby wszystkie sprawdzenia niżej.

Postać zwykła `7z l` jest tabelą kolumnową z nazwą uciętą do szerokości
kolumny — ścieżki nie da się z niej odczytać wiernie, a wyrok o ścieżce
odczytanej niewiernie nie jest wart nic; stąd wybór formatu `-slt`.

Członu `..` w `sprawdzSciezkePozycji` szuka się po rozbiciu ścieżki na
człony, a nie napisem: `strings.Contains(sciezka, "..")` odmówiłby
uczciwemu plikowi `wersja..txt`, a przepuściłby postacie nieprzewidziane
jako napis. Człon jest jednostką, w której ścieżka naprawdę się
rozstrzyga.

## handlers_automations_dobudowa.go

Plik jest osobny od wpięcia rdzenia modułu Automations, bo osobna jest
odpowiedzialność: rdzeń modułu wpina definicję, harmonogram, kolejkę,
zależności i przebiegi, a ten plik — panele i szuflady okien operacyjnych,
którymi Operator sięga po wersje, szablony, logi, ładunki, alarmy, sekrety
i audyt. Port jest jeden, bo moduł jest jeden — rozszerza port rdzenia
modułu, nie stoi obok niego.

Rozgłasza się tylko to, co zmienia byt widziany przez inne okno. Zapis
definicji, publikacja, udostępnienie, etykiety i przywrócenie wersji
zmieniają automatykę, więc idą zdarzeniem zmiany powiązania automatyki —
tym samym, którym idzie zmiana harmonogramu. Odczyty (wykaz wersji,
porównanie, log, kroki, ładunek, audyt, wykaz szablonów, wykaz reguł, wykaz
sekretów) nie rozgłaszają niczego, bo zdarzenie po odczycie byłoby szumem,
na który okna reagowałyby odświeżeniem bez powodu. Wznowienie i odtworzenie
przebiegu rozgłaszają stan przebiegu, ponieważ obie zmieniają to, co
pokazuje Execution Monitor, i robią to przez ten sam nośnik, którym idzie
działanie na kolejce.

Port modułu bez czynności dobudowy to co innego niż brak portu: moduł jest,
a rdzeń nie umie wykonać części jego pracy. Wtedy komendy dobudowy i tak
zostają wpięte i odmawiają wprost, bo odpowiedź „nieznana komenda"
wskazywałaby na brak modułu, a nie na usterkę montażu.

## budowa/server/internal/core/adapter_alerty.go

Kontrakt nie zna komendy przeliczającej reguły osobno. Ewaluacja idzie więc
tam, gdzie operator pyta o wynik: przy odczycie rejestru wyzwoleń komendą
`alert.trigger.list`. Wykaz powstaje z reguł przeliczonych w chwili
odpowiedzi, a nie z wierszy odłożonych wcześniej przez zegar, którego
kontrakt nie zna. Alert pokazujący stan sprzed godziny jako stan bieżący
byłby tą samą fasadą co sonda oddająca stan poprawny bez pomiaru.

Komenda `alert.rule.save` odmawia zapisu reguły, której rdzeń nie umie
zmierzyć, nazywając wprost brakującą miarę. Zapisanie takiej reguły dałoby
operatorowi wiersz w wykazie i ciszę zamiast alertu, czyli złudzenie
bezpieczeństwa, przed którym cała rodzina komend ma chronić.

Rejestr centrum powiadomień przyjmuje wyzwolony alert jako zdarzenie klasy
błąd, po którym operator ma sięgnąć do platformy. Bez tego wpięcia alarm
żyłby wyłącznie w oknie otwartym w chwili wyzwolenia.

Niepowodzenie ewaluacji reguł przy odczycie rejestru wyzwoleń nie przewraca
odczytu: rejestr zastany jest wartościowszy niż odmowa odpowiedzi, a
wyzwolenie, którego nie dało się policzyć w tej chwili, i tak nie miałoby
wartości obserwowanej.

Zdarzenie `alert.triggered` rozgłasza się bez wskazania karty sesji, ponieważ
wyzwolenie jest bytem przekrojowym, a nie bytem jednej sesji — na wzór
zdarzeń modułu Design dotyczących zasobów. Sprawcy zdarzenie nie niesie i
nieść nie może, ponieważ wyzwolenie powstaje z pomiaru rdzenia, a nie z
działania operatora, i kontrakt nie ma na sprawcę pola. Reguła jedzie razem
z wyzwoleniem, ponieważ okno pokazujące alert musi wiedzieć, czyj to alarm i
jakim progiem został postawiony.

## budowa/server/internal/core/skutek_wersji_biblioteki_test.go

Wersja, po której nie da się odtworzyć zawartości, nie jest wersją. Szkoda,
którą ta rodzina sprawdzianów wyklucza, ma dwie postacie. Pierwsza: wersja
wgrana ścieżką odkładała wskaźnik na cudzy plik, więc treść utrwalona
zmieniała się sama, gdy plik na dysku został nadpisany. Druga: przywrócenie
meldowało powodzenie, nie zmieniając tego, co widzi czytelnik treści.

Dlatego żaden sprawdzian tego pliku nie kończy się na porównaniu samych
identyfikatorów. Każdy pyta czytelnika treści, co widzi po zmianie, i
porównuje to z sumą kontrolną bajtów, które do repozytorium naprawdę weszły.
Podgląd jest tu miarą właściwą, nie wiersz w bazie danych.

TestPrzywrocenieWersjiWracaDoJejTresciAWskaznikNaNiaWskazuje jest sprawdzianem
wprost wymierzonym w szkodę wskaźnika i mierzy trzy rzeczy naraz, z których
każda osobno bywała fałszywa: treść po przywróceniu, sumę kontrolną tej
treści oraz wskazanie na wersję, o którą proszono, a nie na najnowszą.

TestPrzywroceniePilnujeGranicPlikuIZostawiaTrescNietknieta pilnuje wersji
wskazanej kodem z cudzej historii. Sprawdzenie samej odmowy nie wystarcza,
ponieważ przepisanie treści mogłoby zajść przed nią, więc sprawdzian mierzy
też treść po odmowie.

TestWersjaZnacznikNieWymazujeTresciBiezacej pilnuje wersji bez treści,
kamienia milowego zakładanego na tym, co w pliku jest. Wersja pusta
zabrałaby plikowi odwołanie do treści przy przywróceniu, więc oznaczenie
kamienia niszczyłoby zasób, który miało utrwalić.

## adapter_narzedzia_sesji.go

Wykaz pozycji po ukośniku nie jest osobną tabelą, lecz składa się na bieżąco
z dwóch źródeł. Komendy kontraktu tworzą pozycje rodzaju `action`: komenda
wykonuje czynność aplikacji i zestawu narzędzi modelu nie zmienia, a jej opis
pochodzi z deklaracji narzędzi modelu, gdzie jest polem obowiązkowym. Katalog
rozszerzeń tworzy pozycje rodzaju `tool`: powołanie narzędzia albo skilla,
jedyne pozycje poszerzające zestaw modelu na czas sesji. Osobnego katalogu
akcji nie ma, bo każdy jego wiersz wskazywałby komendę już obecną w wykazie
z pierwszego źródła i dublowałby pozycję pod inną nazwą modułową.

Zdejmowanie dołożenia oddaje w drugiej wartości pozycje faktycznie zdjęte,
nazwane pełną nazwą. Zdarzenie zdjęcia narzędzia musi nazwać pozycję pełną
nazwą, ponieważ odpowiedź kontraktu niesie wyłącznie zestaw po czynności
i sama nie pozwala odczytać, co zniknęło; żądanie mogło przy tym przyjść
nazwą skróconą.

Metoda oddająca same nazwy dołożeń sesji dla składania zestawu narzędzi tury
modelu jest osobna od wykazu dla operatora, ponieważ tamten oddaje pozycje
w pełnym kształcie, a ta oddaje wskazanie dla procesu modelu. Nazwą jest
nazwa pełna z przedrostkiem źródła, ponieważ to ona jest tożsamością
dołożenia rozpoznawaną po drugiej stronie składania zestawu tury. Straż
odbiornika zerowego w tej metodzie chroni przed zatrzymaniem całego procesu
rdzenia w gorutynie tury, gdyby adapter trafił do składacza zestawu przez
interfejs z zerowym wskaźnikiem schowanym wewnątrz — takie porównanie
u wołającego przechodzi niezauważone.

Przedrostek źródła pozycji rozszerzenia pochodzi z kodu, gdy kod go niesie
w jednym napisie ze źródłem i nazwą skróconą; w pozostałych wypadkach
przedrostek bierze się z pochodzenia pozycji, aby żadna pozycja wykazu nie
została bez źródła.
