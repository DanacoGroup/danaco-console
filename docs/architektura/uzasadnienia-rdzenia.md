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

## szablony_pism_test.go

Sprawdzian wyklucza cztery szkody warsztatu szablonów. Szablon zapisujący
samą treść, bez papieru firmowego, nie jest wzorem pisma. Usunięcie szablonu
fabrycznego, po którym wykazu nie da się odtworzyć bez ponownego wdrożenia,
odbiera warsztatowi podstawę. Pole wymagane bez wartości usunięte z treści
sprawia, że pismo wygląda na kompletne, choć nie jest. Wypełnienie pól
zamianą w napisie treści gubi postać wzorcową pisma — kroje, wcięcia
i granice akapitów.

## budowa/server/internal/core/adapter_narzedzia_obraz_wektor.go

Komenda `image.vectorize` rozstrzyga, co obrysować, w jednym z trzech trybów.
Tryb `outline` odpowiada na pytanie, gdzie kończy się kształt: obraz sprowadza
się do dwóch wartości progiem jasności, a obrysowywany jest obszar ciemny —
tryb właściwy logotypowi, pieczęci, znakowi. Tryb `posterize` odpowiada na
pytanie, z ilu płaszczyzn barwnych składa się obraz: barwy skupiają się w
tylu grupach, ile mówi pole `colors`, a każda grupa jest obrysowywana osobno
— tryb właściwy ilustracji. Tryb `centerline` odpowiada na pytanie, którędy
biegnie kreska: obszar zostaje ścieńczony do linii o grubości piksela i
obrysowany jako kreska — tryb właściwy rysunkowi technicznemu i pismu
odręcznemu, gdzie obrys konturu dałby każdą kreskę jako podwójną pętlę.

Wynik komendy jest zasobem SVG, nie obrazem rastrowym. Idzie do tego samego
magazynu, co każdy inny wytwór rodziny `image.*`, z formatem `svg`. Rdzeń nie
mierzy przy nim wymiarów: dekoder obrazu rastrowego nie rozumie formatu SVG,
a wpisanie tam wymiarów źródła podałoby liczby, których nikt nie zmierzył na
wyniku. Wymiary niosą atrybuty samego dokumentu SVG.

Granica pola obrazu poddawanego obrysowi wynika z kosztu: obrys przechodzi
po każdym pikselu i po każdym jego sąsiedzie, więc koszt rośnie liniowo z
polem, a pamięć maski jest dodatkowa; powyżej szesnastu megapikseli wynik
miałby więcej wierzchołków niż źródło pikseli.

Barwy trybu `posterize` skupia się metodą k-średnich na próbce pikseli, a nie
na komplecie: dla obrazu megapikselowego przejście po wszystkich pikselach w
każdej iteracji kosztuje sekundy, a środki skupień z próbki co dziesiąty
piksel wychodzą praktycznie takie same. Płaszczyzny idą od najciemniejszej,
ponieważ w dokumencie SVG ścieżka późniejsza zasłania wcześniejszą, a
płaszczyzna jasna bywa tłem dla ciemnej.

Środki startowe skupień k-średnich rozkłada się równomiernie po osi jasności
próbki, nie losowo: losowy start dawałby dwa różne wyniki dla dwóch wywołań
na tym samym obrazie, a wektoryzacja ma być powtarzalna.

Wartość tolerancji upraszczania przekłada procent kontraktu na odchylenie w
pikselach: sto procent to odchylenie o osiem pikseli, przy którym z litery
zostaje czworobok, a zero znaczy brak upraszczania — łamana zostaje taka, jak
wyszła z obrysu.

Atrybut `shape-rendering="geometricPrecision"` dokumentu SVG mówi
przeglądarce, żeby nie zaokrąglała wierzchołków do siatki pikseli, bez czego
uproszczona łamana wygląda w podglądzie na bardziej postrzępioną, niż jest.
Zapis współrzędnej pomija zbędne zera, ponieważ dokument SVG z setkami
tysięcy wierzchołków rośnie o megabajty na samych ogonach dziesiętnych.

## budowa/server/internal/core/skutek_obrazu_wektor_test.go

Sprawdzian złożenia obrazu nie kończy się na odpowiedzi i nie ufa polu paths.
Wynik złożenia jest dekodowany z powrotem i mierzony pikselem: sprawdza się,
czy nakładka naprawdę legła w miejscu, w które ją kazano położyć, i czy
krycie naprawdę zadziałało. Wynik wektoryzacji jest czytany jako dokument
SVG: sprawdza się, czy niesie ścieżki, a nie pusty korpus. Żadna z tych
dwóch czynności nie startuje procesu potomnego, więc sprawdzian nie pomija
się przy braku programu zewnętrznego.

TestRozkladNaWarstwyDajeOsobneZasobyZPrzezroczystoscia pomija się bez
silnika segmentacji, ponieważ bez niego komenda odmawia i to jest jej
właściwe zachowanie, sprawdzane osobno w innym miejscu. Pominięcie dotyczy
skutku, którego bez silnika nie ma prawa być.

## budowa/server/internal/core/adapter_narzedzia_obraz_warstwy.go

Segmentacji nie da się policzyć z samego rastra — „gdzie kończy się obiekt

## budowa/server/internal/core/adapter_narzedzia_obraz_warstwy.go

Segmentacji nie da się policzyć z samego rastra — „gdzie kończy się obiekt"
jest pytaniem o znaczenie, nie o piksele. Rozkład idzie dwoma krokami: sieć
segmentująca (rembg, U2-Net) rozdziela obraz na plan pierwszy i tło, zapisując
przynależność w kanale alfa — to krok, którego nie zastąpi żadna arytmetyka,
i jest zależnością tej komendy. Plan pierwszy rozpada się na obszary spójne,
każdy obszar jest jednym obiektem — ten krok liczy rdzeń u siebie, to zwykłe
przejście po tablicy pikseli, więc program zewnętrzny byłby tu zależnością bez
powodu.

Bez silnika segmentacji komenda odmawia, nazywając brak. Nie oddaje całego
obrazu jako jednej warstwy: rozkład, który zwraca to samo, co dostał, jest
atrapą nie do odróżnienia od rozkładu udanego, dopóki liczba warstw nie
zostanie policzona. Kontrakt mówi o tym wprost i ta droga jest tu jedyną.

Obraz, na którym sieć znalazła jeden spójny obiekt, daje jedną warstwę i to nie
jest atrapa: ta warstwa niesie obiekt wycięty z tła, coś, czego w źródle nie
było. Różnica jest sprawdzalna — tło zniknęło.

Wycięcie warstwy samym prostokątem obejmującym wniosłoby do niej kawałek
sąsiedniego obiektu, gdy prostokąty się nachodzą — dlatego przynależność
pikseli jest trzymana osobno od prostokąta.

Przejście po planie pierwszym przy szukaniu obszarów spójnych jest iteracyjne,
z własnym stosem, nie rekurencyjne: obszar megapikselowy przy rekurencji
przepełniłby stos wywołań i przewrócił proces rdzenia, a nie jedno żądanie.

## budowa/server/internal/core/adapter_mowa_wybudzenie.go

Nastawa wybudzania mieszka w konfiguracji, nie w osobnej tabeli: fraza
wybudzająca, tryb nasłuchu, próg detekcji mowy i odszumianie mają poziom
zasięgu i wartość domyślną, zmienianą przez Operatora w oknie konfiguracji.
Osobna tabela dałaby drugie miejsce, w którym mieszka ustawienie, i drugą
drogę jego rozstrzygania, podczas gdy rozstrzyganie po poziomach zasięgu jest
już zbudowane i jest jedno.

Wykonalność jest odpowiedzią, nie awarią: `speech.wake.get` oddaje
`available: false` wraz z powodem, gdy wybudzenia nie da się wykonać, tak
samo jak `speech.availability.get` przy braku silnika. Odmowa kazałaby Voice
Console pokazać błąd tam, gdzie Operator po prostu nie ma jeszcze silnika mowy.

Rdzeń nie ma dostępu do mikrofonu maszyny Operatora i mieć go nie będzie:
mikrofon jest urządzeniem tamtej maszyny, a rdzeń stoi na serwerze. Nasłuch
ciągły jest umową między oknem a rdzeniem, a nie otwarciem urządzenia:
`speech.listen.start` zakłada nasłuch okna i oddaje jego identyfikator, okno
nagrywa u siebie i wysyła kolejne odcinki `speech.audio.upload` wraz
z `windowId`, rdzeń rozpoznaje każdy odcinek silnikiem mowy i ogłasza wynik —
`speech.listen.partial` z tekstem, a gdy w tekście padła fraza wybudzająca,
`speech.wake.detected` — a `speech.listen.stop` nasłuch zamyka. Dzięki temu
podziałowi nasłuch ciągły znaczy dokładnie tyle, ile robi: rdzeń słucha
strumienia, który mu podano, i ogłasza, co usłyszał, bez udawania dostępu do
cudzego mikrofonu. Nasłuch nie jest bramką: rdzeń go sam nie zatrzymuje,
zatrzymuje go Operator.

Metoda `ogloszOdcinekNasluchu` nie zwraca błędu: przyjęcie nagrania powiodło
się niezależnie od tego, czy rozpoznanie odcinka się udało, a odmowa
odebrałaby oknu odnośnik, który już istnieje.

Fraza wybudzająca jest rozpoznawana na rozpoznanym tekście, a nie osobnym
modelem słowa kluczowego: model frazy jest osobną siecią i osobnym plikiem
wag, których instalka nie niesie. Dopasowanie na tekście jest wykonalne
wszędzie tam, gdzie działa rozpoznawanie mowy, i mówi dokładnie to, co robi.

Pole `final` odpowiedzi `speech.listen.partial` jest prawdą, bo odcinek
został rozpoznany w całości: rdzeń dostaje gotowy kawałek nagrania, a nie
strumień w locie. Fałsz obiecywałby oknu poprawkę tego tekstu, która nigdy
nie przyjdzie.

## budowa/server/internal/core/adapter_narzedzia_obraz_zlozenie.go

Cała czynność stoi na bibliotekach image, image/png, image/jpeg i
golang.org/x/image/draw, wkompilowanych w binarium rdzenia. Nie startuje tu ani
jeden proces potomny: złożenie dwóch rastrów to przejście po pikselach, a nie
praca, do której potrzeba silnika zewnętrznego jak ImageMagick — takie wołanie
odebrałoby komendzie prawo do działania na maszynie, która go nie ma, bez
żadnego zysku. Jedna komenda zamyka cztery funkcje opracowania: znak wodny,
branding wsadowy, osadzenie w ramce urządzenia i warstwy rastrowe, bo
wszystkie cztery potrzebują dokładnie tego samego — obrazu położonego na
obrazie. Wynik jest zawsze PNG, bo podstawa bywa fotografią bez kanału alfa,
ale nakładka z przezroczystością wnosi go do wyniku; zapis w JPEG-u zamieniłby
przezroczystość na czarny albo biały prostokąt, a PNG nie traci jakości przy
powtórnym składaniu, którym jest właśnie branding wsadowy.

Nakładka wychodząca poza obszar podstawy nie jest odmową: znak wodny
wypuszczony za krawędź jest przycinany, tak jak przycięłaby go każda inna
warstwa graficzna. Odmowa kazałaby Operatorowi liczyć piksele, zamiast
przesunąć nakładkę.

Przeskalowanie nakładki jeden do jednego byłoby przepróbkowaniem bez powodu, a
każde przepróbkowanie kosztuje ostrość obrazu. Filtr CatmullRom jest brany
zarówno przy pomniejszaniu, jak i powiększaniu, ponieważ znak wodny
pomniejszany najbliższym sąsiadem rozsypuje się na schodki widoczne gołym
okiem.

Mieszanie „prawie zwykłe" byłoby trybem, którego nazwa mówi co innego niż
skutek, dlatego zmieszajNakladke rysuje własną pętlą. Praca na barwie już
przemnożonej przez alfę dałaby mnożenie ciemniejsze przy każdej
półprzezroczystości, dlatego barwy liczone są bez wstępnego mnożenia przez
alfę.

Gałąź domyślna „mieszaj zwykle" w funkcjaMieszania oddałaby złożenie, które
wygląda poprawnie i nie jest tym, o co proszono, dlatego tryb spoza wyliczenia
kontraktu kończy się odmową nazywającą go wprost.

## budowa/server/internal/core/wydanie_formatu_test.go

Szkody, które ten plik ma wykluczyć: wydanie do formatu txt gubiące tabelę w
milczeniu, bo zlecenie nazywa to wprost ciszą zakazaną; wykaz cech pominiętych
oddawany jako zdanie ogólne bez liczb, po którym nie wiadomo, ile stracono;
wydanie wsadowe wstrzymane przez jeden uszkodzony dokument; wydanie oddane
jako udane, a bez zasobu, po którym cienka instalka mogłaby sięgnąć po plik.

## budowa/server/internal/core/widok_nastawy_test.go

Szkody, które ten plik ma wykluczyć: nastawa widoku spoza wyliczenia
kontraktu, którą tabela odrzuciłaby dopiero przy zapisie, po fakcie; różnica
dwóch wersji milcząca o zmianie postaci, gdy litery zostały te same, bo
kontrakt wymaga wprost, żeby zmiana kroju była widoczna jako zmiana; wykaz
różnicy w kolejności wziętej z przebiegu mapy, przez co ten sam dokument
oglądany dwa razy dawałby dwa różne wykazy.

## budowa/server/internal/core/zaleznosci_wykaz_wydruk_test.go

Deklaracja pakietu bez podpowiedzi jest legalna, choć wykaz zależności pilnuje
osobno, żeby jej nie było.

Wstrzymanie silnika kontenerów dotyczy obu jego deklaracji: warsztatu
Developera i modułu Terminal — obie muszą trafić do warstwy decyzyjnej.

Zwykły start rdzenia nie może wpaść w tryb wykazu mimo rozpoznawania obu
postaci znacznika (pojedynczego i podwójnego minusa).

Gdyby reguła rozpoznawania podpowiedzi pytała najpierw o adres GitHuba,
golangci-lint stałby się krokiem ręcznym zamiast wchodzić z podpowiedzi.

## budowa/server/internal/core/adapter_konfiguracja_prowenancja.go

Ten plik wpina dwie komendy rodziny `config.*`, które nie dotykają rejestru
ustawień: `config.explain.get` (prowenancja wywołania modelu) i
`config.window.open` (zakres, na którym otwiera się okno konfiguracji). Stoją
osobno od `handlers_config.go` i `handlers_sesja_konfiguracja.go`, bo tamte
pliki wpinają komendy czytające i piszące rejestr ustawień; te dwie nie
zapisują niczego — wspólny mają wyłącznie przedrostek nazwy.

Prowenancja idzie tą samą drogą co tura, nie drugą: `config.explain.get` nie
ma własnego składacza wywołania, tylko bierze dokładnie te funkcje, którymi
jedzie `message.send` — zapytanieKanalu, nakladkaOkna, uzupelnijSrodowisko,
uzupelnijKonfiguracje, ustawieniaKanaluGlownego, nakladkaKanaluGlownego,
injection.Argumenty. Gdyby prowenancja składała wywołanie po swojemu,
pokazywałaby wiersz, którego tura nigdy nie wykona.

Od wiersza procesu prowenancja różni się świadomie w dwóch miejscach. Po
pierwsze, treść `--settings` i `--mcp-config` jedzie tu jako treść, a nie
jako ścieżka pliku tymczasowego: tura materializuje te napisy na dysku tuż
przed uruchomieniem procesu (`injection/materializacja.go`) i sprząta je po
sobie, więc pliku o tej ścieżce jeszcze nie ma; tak samo rozstrzyga sam pakiet
injection, którego prowenancja niesie treść pierwotną, a argv realną ścieżkę
(`injection/przebieg.go`). Po drugie, `argv[0]` jest programem: kontrakt tej
komendy nie ma osobnego pola `program`, więc bez tego nie byłoby widać, jaki
plik wykonywalny rusza, a `argv` ma być wierszem wywołania, nie samą listą
przełączników.

Rejestr oddaje wyłącznie kanały czynne. Kanał wskazany oknem, a nieobecny
w rejestrze, wywróciłby także turę (`models/wysylka.go`) — odmowa idzie tym
samym rozpoznaniem, zamiast pokazywać wiersz wywołania, którego nikt nie
wykona.

Suma kontrolna liczona w prowenancji dotyczy samej konstytucji, nie całej
nakładki: profil roli i ekspertyza zadaniowa zmieniają się z oknem,
a konstytucja jest warstwą, której niezmienność się sprawdza. Pusta
konstytucja daje pusty skrót, bo nie ma czego sumować.

Zapytanie prowenancji nie niesie treści wypowiedzi ani historii: żadna z nich
nie wchodzi do wiersza wywołania ani do promptu systemowego, a zmyślona treść
pytania byłaby atrapą. Załączników tu nie ma i nic ich nie udaje z tego samego
powodu — wymyślony załącznik byłby ścieżką, której Operator nie dołączył.
Wznowienie rozmowy wchodzi do argv przełącznikiem `--resume`, więc bez niego
wiersz byłby wierszem innej tury niż ta, która pójdzie. Ekspert okna nakłada
model, ustawienia i mosty MCP w tej samej kolejności co w turze; pominięcie go
tutaj pokazywałoby Operatorowi wiersz sprzed wyboru eksperta i nazywałoby go
prowenancją bieżącej tury.

Karta z wieloma oknami jest odmową, nie zgadywaniem przy ustalaniu okna
prowenancji: sesja niesie okna w relacji 1:N, a każde okno ma własny kanał,
własny katalog i własną konfigurację — wybranie „pierwszego z brzegu

## budowa/server/internal/core/adapter_konfiguracja_prowenancja.go

Ten plik wpina dwie komendy rodziny `config.*`, które nie dotykają rejestru
ustawień: `config.explain.get` (prowenancja wywołania modelu) i
`config.window.open` (zakres, na którym otwiera się okno konfiguracji). Stoją
osobno od `handlers_config.go` i `handlers_sesja_konfiguracja.go`, bo tamte
pliki wpinają komendy czytające i piszące rejestr ustawień; te dwie nie
zapisują niczego — wspólny mają wyłącznie przedrostek nazwy.

Prowenancja idzie tą samą drogą co tura, nie drugą: `config.explain.get` nie
ma własnego składacza wywołania, tylko bierze dokładnie te funkcje, którymi
jedzie `message.send` — zapytanieKanalu, nakladkaOkna, uzupelnijSrodowisko,
uzupelnijKonfiguracje, ustawieniaKanaluGlownego, nakladkaKanaluGlownego,
injection.Argumenty. Gdyby prowenancja składała wywołanie po swojemu,
pokazywałaby wiersz, którego tura nigdy nie wykona.

Od wiersza procesu prowenancja różni się świadomie w dwóch miejscach. Po
pierwsze, treść `--settings` i `--mcp-config` jedzie tu jako treść, a nie
jako ścieżka pliku tymczasowego: tura materializuje te napisy na dysku tuż
przed uruchomieniem procesu (`injection/materializacja.go`) i sprząta je po
sobie, więc pliku o tej ścieżce jeszcze nie ma; tak samo rozstrzyga sam pakiet
injection, którego prowenancja niesie treść pierwotną, a argv realną ścieżkę
(`injection/przebieg.go`). Po drugie, `argv[0]` jest programem: kontrakt tej
komendy nie ma osobnego pola `program`, więc bez tego nie byłoby widać, jaki
plik wykonywalny rusza, a `argv` ma być wierszem wywołania, nie samą listą
przełączników.

Rejestr oddaje wyłącznie kanały czynne. Kanał wskazany oknem, a nieobecny
w rejestrze, wywróciłby także turę (`models/wysylka.go`) — odmowa idzie tym
samym rozpoznaniem, zamiast pokazywać wiersz wywołania, którego nikt nie
wykona.

Suma kontrolna liczona w prowenancji dotyczy samej konstytucji, nie całej
nakładki: profil roli i ekspertyza zadaniowa zmieniają się z oknem,
a konstytucja jest warstwą, której niezmienność się sprawdza. Pusta
konstytucja daje pusty skrót, bo nie ma czego sumować.

Zapytanie prowenancji nie niesie treści wypowiedzi ani historii: żadna z nich
nie wchodzi do wiersza wywołania ani do promptu systemowego, a zmyślona treść
pytania byłaby atrapą. Załączników tu nie ma i nic ich nie udaje z tego samego
powodu — wymyślony załącznik byłby ścieżką, której Operator nie dołączył.
Wznowienie rozmowy wchodzi do argv przełącznikiem `--resume`, więc bez niego
wiersz byłby wierszem innej tury niż ta, która pójdzie. Ekspert okna nakłada
model, ustawienia i mosty MCP w tej samej kolejności co w turze; pominięcie go
tutaj pokazywałoby Operatorowi wiersz sprzed wyboru eksperta i nazywałoby go
prowenancją bieżącej tury.

Karta z wieloma oknami jest odmową, nie zgadywaniem przy ustalaniu okna
prowenancji: sesja niesie okna w relacji 1:N, a każde okno ma własny kanał,
własny katalog i własną konfigurację — wybranie „pierwszego z brzegu"
pokazałoby Operatorowi prowenancję cudzego okna i wyglądałoby na odpowiedź.
Rdzeń nie zna tu ogniska klienta, więc mówi wprost, czego mu brakuje.

Kanał bez procesu nie ma wiersza wywołania i nic tego nie udaje: adaptery
`echo` i `api` nie uruchamiają programu, więc dla nich `argv` jest pustą
tablicą, nie brakiem pola i nie zmyślonym wierszem `claude …`, dokładnie tak,
jak rozstrzyga prowenancja tych kanałów w `models/prowenancja.go`. Pozostałe
pola odpowiedzi — prompt systemowy, ustawienia, suma konstytucji — są dla nich
równie prawdziwe jak dla kanału głównego.

Trzy stany treści ustawień w prowenancji są nazwane wprost: treść JSON jedzie
sobą; wartość, która jest ścieżką pliku zapisanego wcześniej przez inną
warstwę, jedzie jako napis, bo surowa wywróciłaby kodowanie odpowiedzi; brak
wartości jedzie jako `null`, co znaczy, że proces nie dostanie przełącznika
`--settings`, a nie że przekazano pusty zestaw reguł.

`OtworzOknoKonfiguracji` nie ma w kontrakcie nic więcej do zrobienia poza
podaniem zakresu, więc robi dokładnie to i nic ponadto: sprawdza, że byt, na
którym okno ma stanąć, istnieje, że wskazany obszar należy do kontraktu i że
wskazany poziom zasięgu ma swój byt, po czym oddaje rozstrzygnięty zakres.
Okno konfiguracji jest bytem klienta; rdzeń nie ma tabeli jego otwarć i nie
zakłada jej pod jedno pole `opened`. To pole nie jest stałą: fałszu ta komenda
nie zwraca, zwraca odmowę — nieistniejące okno albo karta sesji to
`not_found`, obszar spoza kontraktu i poziom bez bytu to `validation_failed`;
odpowiedź `opened: false` bez powodu byłaby ciszą udającą wynik.

Bez wskazania poziomu zasięgu obowiązuje poziom najwęższy, jaki żądanie
potrafi wskazać bytem: okno komunikacji, w jego braku karta sesji, w braku
obu poziom globalny. Wskazanie wprost jest brane, ale nie na wiarę: poziom
okna bez okna i poziom karty bez karty ani okna to zakres, którego nie da się
pokazać. Okno zna swoją kartę, więc wskazanie okna wystarcza także poziomowi
karty sesji, i wystarcza wtedy, gdy klient karty nie podał.

## budowa/server/internal/core/uprzaz_wejscia_test.go

Sprawdzian mierzy skutek wychodzący poza proces: konto zapisuje się w bazie,
a droga potwierdzenia idzie listem — sprawdzian, który mierzy wyłącznie
kopertę odpowiedzi, przepuściłby zarówno list, który poszedł, jak i list,
który przepadł, bo odpowiedź wygląda tak samo. Uprząż zgodności z kontraktem
tego nie mierzy: nie zwraca bazy wołającemu, bo dowodem założenia konta jest
wiersz w tabeli konta właściciela, nie zdanie w odpowiedzi; nie stawia
odbiornika SMTP na pętli zwrotnej, bo dowodem wysłania listu jest list.

Odbiornik testowy jest prawdziwym gniazdem na porcie efemerycznym, nie
zaślepką podstawioną w miejsce funkcji wysyłki: zaślepka sprawdzałaby, czy
rdzeń woła funkcję, gniazdo sprawdza, czy list doszedł, i pozwala przeczytać
jego treść, rozstrzygając, że niesie drogę i nie niesie hasła.

Sieci sprawdzian nie dotyka: nasłuch stoi na pętli zwrotnej, a nastawy idą
z wyłączonym StartTLS, więc rozmowa nie próbuje ani podnieść TLS, ani wyjść
poza maszynę.

Sejf poświadczeń, w którym leży znacznik bramki bez poczty i sekret kotwicy,
jest otwierany nad tym samym katalogiem danych, który dostał rdzeń, więc
czytany jest ten sam plik, do którego rdzeń pisze.

Droga zapisu konta nadawczego w sprawdzianie jest ta sama co w produkcie, nie
skrót przez nastawy montażu: badane jest właśnie to, co dzieje się, gdy poczta
pojawia się po rejestracji.

## budowa/server/internal/core/handlers_developer.go

Plik wpina komendy obszaru `developer.*` modułu Developer wraz z jego oknami
operacyjnymi: Code Editor, Project Tree, Git Panel i Build Output. Kontrakt
daje modułowi jedno zdarzenie, `developer.build.changed`, więc każdy przyrost
budowania — start, kolejny wiersz logu, domknięcie — rozgłasza się przebiegiem
po zmianie; Build Output odświeża się z jednej subskrypcji, a nie z
odpytywania.

Podział dróg w Git Panelu jest jawny: `developer.git.action` wykonuje
czynności zmieniające repozytorium wedle zamkniętego słownika, a
`developer.git.status`, `.diff`, `.log`, `.branch.list` i `.conflict.*`
wyłącznie czytają. Odczyt stoi na bibliotece `go-git` wkompilowanej w rdzeń
i nie startuje ani jednego procesu — Git Panel ma się otwierać także tam,
gdzie programu `git` nie ma, bo instalka go nie niesie.

`PodepnijPrzyrostBudowania` oddaje adapterowi drogę do zdarzenia przyrostu.
Budowanie kończy się poza wykonaniem komendy, czasem minuty później, a log
narasta przez cały ten czas, więc rozgłoszenie nie może iść wyłącznie
z obsługiwacza żądania. Sesja komunikatu przyrostu zostaje pusta: przebieg
należy do okna, a rdzeń rozgłasza go także wtedy, gdy okna nie ma już
w rejestrze — budowanie przeżywa zamknięcie okna, a klient ma prawo zobaczyć
jego koniec.

Odczyt repozytorium w Git Panelu to sześć czynności na bibliotece `go-git`
wkompilowanej w rdzeń, bez ani jednego procesu potomnego.

Historia pliku edytora to wersja robocza zakładana przy zapisie. To nie jest
historia repozytorium: tamta należy do Gita i jedzie `git.log`.

Warstwa językowa Code Editora jako jedyna w module woła programy serwera
(gopls, gofmt, goimports, prettier, golangci-lint), bo program jest tą wiedzą
o kodzie; brak programu wraca polem `serverAvailable`/`linterAvailable`,
a nie odmową całej komendy.

Odczyt okna Build Output daje historię przebiegów, log, wynik testów
i pokrycie. Wynik testów i pokrycie powstają z rozbioru wyjścia w chwili
biegu, a nie z ponownego czytania przyciętego dziennika.

Okno Run & Debug na protokole DAP: punkt przerwania jest trwały i należy do
okna, a sesja debugowania jest żywa i gaśnie razem z procesem adaptera.

Zakładka Data Console stoi na `database/sql` ze sterownikami wkompilowanymi
w rdzeń. Hasło nie leży w opisie połączenia, tylko w sejfie.

Zakładka Containers stoi na Docker SDK Go rozmawiającym z gniazdem silnika.
Silnika nie da się wkompilować: jego brak wraca polem `engineAvailable`.

Zależności, bezpieczeństwo i jakość liczą własne parsery manifestów i własne
reguły skanowania, bez zależności od programu spoza instalki.

Przebieg obciążeniowy stoi przy zapytaniu pojedynczym, bo jest tym samym
zapytaniem powtórzonym pod obciążeniem, z tym samym podstawianiem zmiennych
środowiska kolekcji.

## budowa/server/internal/core/adapter_przejecie_sterowania.go

Plik nie prowadzi biegu: nie liczy obiegów, nie wykrywa braku postępu, nie
rozpoczyna tur — robi to wyłącznie pętla koordynator-wykonawca. Nie zapisuje
też do pozycji kolejki: wznowienie drogą queue.action resume woła krok, a ten
przesuwa pozycję po mapie kroku naprzód, czyli pracę przerwaną w stanie
wykonywana traktuje jak skończoną.

Rejestr steru zna swój stan w chwili obiegu i wydaje go jednorazowo, tak jak
rejestr biegów, z tą różnicą, że pytanie „kto steruje tym zleceniem” pada
z zewnątrz w dowolnej chwili, a rejestr nie prowadzi biegu i nie zatrzymuje go.

Wykaz pusty przy sprzątaniu rejestru steru niczego nie kasuje — brak wiedzy
o oknach nie jest wiedzą o ich zamknięciu.

Rejestr steru jest bytem pakietowym, nie polem struktury, bo klucz rejestru
(identyfikator zewnętrzny okna) jest niepowtarzalny w procesie i dwa rejestry
dałyby dwie odpowiedzi o tym samym oknie. Odpis biegu kontraktu w
core/stan_obiegu.go po ster nie sięga: stan pętli nie ma pola o sterującym,
więc rejestr obsługuje wyłącznie rodzinę control.*.

Repozytorium przekazań pisze i czyta ślad w dzienniku akcji okna. Kolejność
czynności w Przejmij jest treścią, nie stylem: ślad idzie do dziennika
pierwszy, bo tylko on może się nie udać, a przejęcie bez zapisu zostałoby bez
świadka; ster wchodzi do rejestru przed zatrzymaniem, żeby rozgłoszenie stanu
biegu wywołane zatrzymaniem niosło już nowego sterującego; zatrzymanie pętli
staje na końcu i nie kasuje niczego z dorobku Koordynatora — Operator wchodzi
tam, gdzie proces stoi.

Wznowienie w Oddaj kasuje wyłącznie licznik braku postępu i ostatni odcisk
strumienia: historia obiegów zostaje, więc Koordynator podejmuje bieg, a nie
zaczyna go od nowa. Cicha zgoda na oddanie nieprzejętego zlecenia wznowiłaby
bieg, którego Operator nie zatrzymywał — zmianę stanu, o którą nikt nie
prosił.

Szeroki kontrakt obszaru window.* deklaruje w całości inny plik, dlatego
historia sięga repozytorium wąskim interfejsem. Port, który tej zdolności nie
niesie, dostaje odmowę nazywającą brak, nigdy pusty wykaz udający, że nikt nie
przejmował.

## budowa/server/internal/core/adapter_rozmowa_powierzchnia.go

Plik przekłada obszary tools, permissions, hooks, skills, environment,
provider i mcp konfiguracji sesji na trzy powierzchnie procesu kanału, które
warstwa injection już potrafi złożyć: napis `--settings` (plik ustawień
sesji) z regułami uprawnień i narzędzi, sekcją hooks oraz odmową narzędzia
Skill, gdy obszar skills jest wyłączony; zmienne środowiskowe procesu, czyli
obszar environment oraz adres dostawcy; i osobny `--mcp-config` z wiązaniami
serwerów MCP opisanymi wprost. Ten plik jest jedynym miejscem w drzewie,
które zna kształt pliku ustawień dostawcy i nazwy zmiennych środowiskowych —
model konfiguracji pozostaje dziedzinowy, dopiero tutaj staje się
„permissions.allow” czy „ANTHROPIC_BASE_URL”.

`uprawnieniaCLI` odwzorowuje obszar permissions oraz tools na sekcję
permissions pliku ustawień dostawcy. Trybu domyślnego tu nie ma: tryb
uprawnień jedzie przełącznikiem `--permission-mode` (`z.TrybUprawnien`,
`przelozUprawnienia`), którego słownik kontrakt potwierdza dosłownie —
powtórzenie go w pliku groziłoby rozejściem słownictwa i dwoma źródłami tej
samej decyzji.

`plikUstawienZKonfiguracji` buduje napis `--settings` z obszarów permissions,
tools, skills oraz hooks. Reguły narzędzi i wyłączenie obszaru skills
dokładają się do reguł uprawnień; zaczepy jadą osobną sekcją hooks.

`dodajRegulyUmiejetnosci` przekłada obszar skills na regułę uprawnień. Jedyny
przekład, który powierzchnia pliku ustawień unosi bez atrapy: wyłączenie
obszaru wprost (`Enabled == false`) odmawia narzędzia Skill w sekcji deny —
tą samą drogą, którą obszar tools odmawia narzędzi imiennych. Obszar włączony
albo nieokreślony nie dokłada reguły: umiejętności pozostają wtedy dostępne,
jak przed wpięciem. Dopuszczanie imienne (`AllowedSkillIds`), katalogi
wyszukiwania (`Directories`) i samowykrywanie (`AutoDiscovery`) nie mają pola
na tej powierzchni i jadą do ryzyk — nie ma tu dla nich cichej atrapy.

`hooksZKonfiguracji` buduje sekcję hooks pliku ustawień z obszaru hooks.
Obszar wyłączony wprost (`Enabled == false`) nie daje żadnego zaczepu; obszar
włączony albo nieokreślony przenosi zaczepy czynne. Zaczep bez zdarzenia albo
bez polecenia jest niekompletny i nie jedzie. Zaczepy o tym samym zdarzeniu
i zawężeniu zbierają się w jednej grupie. Brak zaczepów daje nil, więc sekcja
znika z JSON.

`granicaSekund` przelicza granicę czasu zaczepu z milisekund kontraktu na
sekundy pliku ustawień. Wartość niedodatnia nie daje granicy (nil); wartość
dodatnia poniżej sekundy zaokrągla w górę do jednej sekundy — pole timeout
nie wyraża ułamka, a granica poniżej pełnej sekundy nie może zejść do zera
i zamienić się w brak granicy.

`srodowiskoZKonfiguracji` składa zmienne środowiskowe procesu z obszaru
environment oraz z adresu i granicy odpowiedzi obszarów provider i model.
Zmienna tajna (`SecretRef`) nie wchodzi: jej treść zna wyłącznie sejf,
którego ta droga nie ma wpiętego — wpisanie nazwy bez wartości byłoby atrapą.

`mcpZKonfiguracji` buduje napis `--mcp-config` z wiązań obszaru mcp opisanych
wprost (transport stdio z programem albo sse/http z adresem). Wiązania
wskazujące punkt dostępu (`AccessPointId`) pomija: ich adres i poświadczenie
żyją w rejestrze punktów dostępu, którego ta droga nie rozstrzyga — jadą one
drogą nadań okna (mosty).

## budowa/server/internal/core/zaleznosci_zewnetrzne.go

Plik niesie jeden wykaz programów spoza instalki, których rdzeń używa, wraz
z sondą sprawdzającą ich obecność przy starcie. `zewnetrzne.Wolaj` sprawdza
obecność programu przed uruchomieniem i odmawia zdaniem nazywającym brak
(`zewnetrzne.BrakNarzedzia`) — to jest właściwe zachowanie w chwili
czynności, ale za późne jako informacja: Operator dowiaduje się o braku
dopiero po naciśnięciu przycisku, osobno dla każdej funkcji, i nie ma skąd
wiedzieć, ile jeszcze takich niespodzianek przed nim. Sonda startowa odwraca
kolejność — rdzeń mówi na wejściu, czym dysponuje i czego mu brak, jednym
wykazem; to jest ta sama wiedza, ale podana zanim zawiedzie czynność.

Każde narzędzie jest zadeklarowane tam, gdzie jest używane — Pandoc przy
dokumentach, ffmpeg przy nagraniach, Tesseract przy rozpoznaniu pisma. Wykaz
bierze te deklaracje, zamiast wypisywać nazwy po raz drugi: druga lista
rozjechałaby się z pierwszą przy pierwszej zmianie pakietu i Operator
czytałby podpowiedź instalacyjną prowadzącą donikąd.

Kolejność wykazu `zaleznosciZewnetrzne` jest ustalona, żeby dziennik startu
miał się różnić wtedy, gdy zmienił się stan maszyny, a nie wtedy, gdy inaczej
ułożyła się mapa.

Tika i LanguageTool to programy Javy, więc obie pozycje wskazują ten sam plik
wykonywalny (`java`) i stoją w wykazie osobno. Nie jest to powtórzenie:
pozycja wykazu odpowiada na pytanie „co przestaje działać i co z tym zrobić",
a odpowiedzi są tu dwie różne — brak Javy zabiera naraz odczyt plików
i korektę językową, czyli dwa zakresy w dwóch modułach. Wiersz na zakres jest
tym, po co ten wykaz istnieje.

Wiersz zbiorczy `zglosZaleznosci` idzie zawsze — także wtedy, gdy nie brakuje
niczego, bo „wszystkie obecne" jest informacją, nie ciszą. Każdy brak
dostaje własny wiersz z zakresem, który przestaje działać, i podpowiedzią
instalacyjną: Operator ma po starcie wiedzieć, czego nie zrobi, zanim
spróbuje.

## budowa/server/internal/core/adapter_okna.go

Klient zakładał okno z kodem kanału wziętym z wyliczenia rodzajów, którego
pierwsza pozycja to "cli", a rejestr rdzenia niesie kanał o kodzie
lokalny-claude rodzaju "cli". Każde okno świeżej instalacji wskazywało więc
kanał, którego nie ma, i pierwsza wypowiedź Operatora wracała odmową, że kanał
"cli" nie istnieje w rejestrze.

Ta sama pomyłka powtarzała się przy module: klient zakładał okno z modułem
wziętym z wyliczenia znanych identyfikatorów modułu, a pierwsza pozycja tam to
"talkin" — kod środowiska, nie modułu. Okno wskazywało więc moduł, którego
w katalogu nie ma, a panel sterowania pokazywał moduł spoza wykazu.

Więź koordynator-wykonawca ustanawia komenda window.handoff, zapisując ją do
SQLite. Dopóki wykaz czytał koordynatora wyłącznie z pamięci, zaraz po udanym
przekazaniu pokazywał wykonawcę bez koordynatora, choć wiersz w bazie był
poprawny: trwałość działała, a żywy widok jej nie widział.

Wiersz okna powstaje leniwie, przy pierwszej wiadomości, więc okno świeżo
otwarte nie ma go jeszcze wcale — czyszczenie więzi na tej podstawie
gubiłoby informację prawdziwą. Bez repozytorium przekazań adapter pracuje jak
dotąd, na samej pamięci: trwałość jest dodatkiem, nie warunkiem pracy rdzenia.

Wybór eksperta w oknie idzie dalej niż pamięć: rejestr nadzorcy przyjmuje go
od migracji 084, ale wiersz okna dostawał go wyłącznie przy zakładaniu, więc
zmiana eksperta w oknie już utrwalonym ginęła przy restarcie rdzenia. Zapis
stoi po zmianie w rejestrze, a nie przed nią: nie ma po co utrwalać wyboru,
którego pakiet sesji nie przyjął.

## budowa/server/internal/core/adapter_kondycja_pomiar.go

Ten plik niesie sam pomiar sondy kondycji, każdy z pięciu rodzajów mierzy coś
naprawdę: `http` wysyła żądanie pod adres i patrzy na kod odpowiedzi; `tcp`
otwiera połączenie z gniazdem i patrzy, czy się otworzyło; `internal` dotyka
wnętrza rdzenia — bazy stanu albo jego własnej pamięci, co nie jest „zwróć
w porządku": baza odpytana jest bazą, która odpowiedziała, a jej czas obiegu
jest zmierzoną liczbą; `command` uruchamia program i patrzy na jego kod
wyjścia; `modelCall` wysyła krótkie zapytanie kanałem modelu i czeka na
odpowiedź. Definicje, seria i dostępność leżą w `adapter_kondycja.go`.

Czego tu nie ma: gałęzi „nie umiem zmierzyć, więc `up`". Każdy powód, dla
którego pomiar się nie odbył — brak uruchamiacza, brak rejestru kanałów,
nieznany cel sondy wewnętrznej — kończy się stanem `unknown` wraz ze zdaniem
mówiącym, czego brakuje. `unknown` znaczy „nie wiem" i tylko tak wygląda
w wykazie; `up` znaczyłoby „sprawdziłem i jest dobrze", a nikt nie sprawdzał.

Klient `zmierzHttp` jest budowany na jeden przebieg i nie chodzi za
przekierowaniami dalej niż pięć razy: sonda ma zmierzyć adres, który podał
Operator, a nie zwiedzić łańcuch przekierowań do cudzej strony błędu.

`zmierzWnetrze` mierzy sam rdzeń dwoma celami, oba mierzalne: `database`
odpytuje bazę stanu (prawdziwe zapytanie, nie sprawdzenie wskaźnika),
`runtime` czyta liczniki procesu. Cel spoza tych dwóch kończy się stanem
„nie wiem" wraz z wykazem znanych — zgadywanie, o co Operatorowi chodziło,
dałoby pomiar czegoś innego niż prosił.

Program w `zmierzProgram` idzie tą samą drogą co każde inne wołanie arsenału
(`zewnetrzne.Wolaj`): przez port uruchamiacza, bramę izolacji i objęcie
drzewa procesów. Własnego `exec.Command` tu nie ma — proces uruchomiony obok
tej drogi wypada spod nadzoru i zostaje po nim uchwyt.

## budowa/server/internal/core/autozapis_skutek_test.go

Niepowodzenie zapisu jest wywoływane zabraniem rdzeniowi tabeli, w którą
zapisuje postać dokumentu — to niepowodzenie tej samej klasy, co awaria dysku
albo uszkodzony plik bazy, bo zapis wchodzi tą samą drogą i pada w tym samym
miejscu; atrapa repozytorium mierzyłaby wtedy tylko atrapę.

Zabranie całej tabeli w sprawdzianie awarii zapisu mierzyłoby co innego niż
zamierzone: tam pada już odczyt i nie ma nawet czego zapisać, dlatego
wyzwalacz odrzuca tylko zapis postaci, a odczyt zostawia nietknięty.

Odróżnialność wersji szeregu autozapisu od szeregu Operatora wychodzi dziś
zawężeniem wykazu i dwoma licznikami, bo pole szeregu w kontrakcie pozycji
jeszcze nie istnieje — odcinek kontraktu ma je dołożyć, a sprawdzian mierzy
stan prawdziwy, taki, jaki jest dziś zgłoszony.

## budowa/server/internal/core/adapter_narzedzia_wynik.go

Cztery kopie jednej reguły byłyby czterema okazjami do rozjazdu, dlatego
reguła jest jedna: okno wyniku bierze się z nieobowiązkowego pola windowId
żądania, a gdy go nie ma, z okna zasobu źródłowego; rodzaj zasobu rozstrzyga
format wyniku, nie nazwa komendy, bo warunek CHECK kolumny rodzaju poszerza
migracja 113, więc film nie musi jechać jako image, a PDF jako vector;
odłożenie bajtów i założenie wiersza idą jedną sekwencją: zapisz do magazynu,
złóż zasób Designu, zapisz wiersz, przełóż na kontrakt; magazynem bajtów jest
magazyn zasobów Designu, ten sam, którym jedzie design.asset.upload, bo dwa
składy bajtów byłyby dwiema prawdami o tym, gdzie rdzeń trzyma treść poza
bazą.

Odmowy zostają przy rodzinach: treść odmowy nazywa brak właściwy rodzinie
("brakuje pola assetId albo sourcePath — nie ma czego zmierzyć" mówi co innego
niż "nie ma czego spakować"), a wspólna formułka zamieniłaby cztery zdania
mówiące na jedno milczące, dlatego wspólne funkcje przyjmują opakowywacz
odmowy rodziny i nie znają żadnego kodu błędu z własnej głowy.

Pusty wynik oknaWynikuArsenalu znaczy "bez wiersza", a nie "bez okna": kolumna
zasob_design.okno jest NOT NULL i zapis odmawia zasobowi bez okna. Gdy nie
wiadomo, do którego okna wynik należy, wiersz świadomie nie powstaje: bajty
leżą w magazynie pod sumą kontrolną, model dostaje odwołanie, którym plik da
się otworzyć, a Assets Panel nie dostaje kafelka o zmyślonej przynależności.
Kolejność pytań jest zamierzona: najpierw pole żądania, bo model zna własny
zasięg i podaje okno, w którym Operator na wynik czeka; potem okno zasobu
źródłowego, odczytane z wiersza źródła, a nie wymyślone. Nazwa okna, którego
nie ma w rejestrze okien, byłaby przynależnością zmyśloną — półka nazwana
wygląda porządniej niż brak wiersza, ale mówi nieprawdę.

Rozstrzyganie rodzaju po formacie, nie po nazwie komendy, jest istotą tablicy
rodzajeZasobuPoFormacie: media.transcode z czynnością frame daje obraz, a z
extractAudio dźwięk, mimo że komenda jest ta sama; document.convert do png
daje obraz, a nie dokument. Poszerzenie tablicy rodzajów jest tańsze niż
odmowa, bo bajty powstały i model ma prawo dostać do nich odwołanie.

Kolejność "najpierw treść, potem wiersz" w odlozWynikArsenalu jest wspólna
całemu rdzeniowi, bo wiersz wskazujący odwołanie, za którym nic nie leży,
jest w panelu kafelkiem bez zawartości. Zasób bez wiersza jest co do bajtów
zasobem prawdziwym: DesignAsset złożony z ręki niesie uri wskazujące blob pod
sumą sha256, a pusty windowId jest uczciwym oświadczeniem, że wynik do żadnego
okna nie należy — model dostaje odwołanie, Operator nie dostaje w panelu nic,
bo o nic nie prosił.

## budowa/server/internal/core/adapter_doradcy.go

Ten plik niesie adapter doradcy: wykonanie konsultacji tą samą drogą, którą
rdzeń woła każdy inny model. O tym, który model radzi, rozstrzyga sufit
siły, a nie ten plik — model z własnej inicjatywy konsultuje wyłącznie model
równy sobie albo słabszy. Reguła i jej powód stoją w
`podagenci/doradca_wybor.go`; adapter jej nie powtarza, żeby nie było dwóch
miejsc, w których wolno ją poluzować. Komenda `advisor.consult`, jedyny
wołacz konsultacji, stoi w `adapter_doradcy_konsultacja.go`.

Adapter nie zna SQL-a ani protokołu dostawcy. Bierze rejestr kanałów modelu
i dziennik konsultacji pakietu `podagenci`, a pojęcie doradcy — kto nim może
być, jak brzmi pytanie, czym jest rada — zostaje po stronie pakietu
pojęciowego.

Trzy własności pojęcia widać tu jako trzy czynności. Jawność: wszystkie
fragmenty strumienia doradcy idą do ujścia wołającego bez zmiany, a na
koniec dokładany jest blok `Rada.Jawnie` — Operator widzi pytanie, doradcę
i radę w tym samym oknie, w którym pracuje agent; adapter nie oddaje samej
treści rady bez wskazania doradcy. Prowenancja: opis wywołania składa kanał,
tak jak przy każdym innym zapytaniu (`models.NadajProwenancje`) — adapter go
nie wytwarza, tylko przejmuje fragment `provenance` przelatujący strumieniem
i zapisuje ten sam napis w dzienniku; drugiej prowenancji nie ma. Konsultacja,
nie delegacja: zwrócona `Rada` nie zmienia niczego w stanie rdzenia — nie
zakłada pozycji kolejki, nie startuje tury, nie zapisuje wiadomości agenta.
Wołający dostaje radę i sam rozstrzyga, co z nią zrobi.

Odmowa też trafia do dziennika. Konsultacja, która się nie odbyła, jest
zdarzeniem, o które Operator zapyta jako pierwsze — wpis o stanie `odmowa`
niesie powód.

Rejestr okien sesji, przez który komenda `advisor.consult` odmawia gdy nie
jest wpięty (bez gaszenia samej metody `Skonsultuj`), bez tego rejestru
musiałby wierzyć modelowi na słowo, kim jest — a wtedy sufit siły dałoby się
obejść jednym polem żądania.

Nadajnik, gdy niewpięty, nie gasi konsultacji — gasi wyłącznie jej
widoczność w oknie; bez niego jawność konsultacji kończy się na buforze,
którego nikt nie czyta.

Blok jawności idzie strumieniem po odpowiedzi doradcy: Operator widzi
najpierw to, co doradca powiedział, a potem podpis mówiący, że to była rada
cudza i niewiążąca.

Ujście `Skonsultuj` pochodzi od wołającego i jest zwykle tym samym ujściem,
którym płynie odpowiedź agenta — dlatego rada widoczna jest tam, gdzie
pracuje agent, a nie w osobnym, cichym kanale.

Sufit siły stoi w pakiecie pojęciowym (`podagenci/doradca_wybor.go`),
a nie tutaj — adapter go wykonuje, nie ma własnej wersji. Odmowa wraca
wołającemu nazwana i jednocześnie idzie do dziennika jako wpis `odmowa`, bo
pytanie „dlaczego rdzeń nie zapytał mocniejszego modelu" pada po fakcie
i musi mieć odpowiedź w bazie.

`zbierakRadyDoradcy` przepuszcza wszystkie fragmenty, nie tylko tekst.
Gdyby zatrzymywał je u siebie, konsultacja byłaby niewidoczna do chwili jej
zakończenia, a rada pokazana dopiero jako gotowy napis — czyli tak samo jak
odpowiedź własna agenta, wbrew jawności konsultacji.

Kod odmowy `bladDoboruDoradcy` mówi, czy ponawiać. `channel_unavailable`
jest w kontrakcie kodem ponawialnym (`shared.KodyPonawialne`), więc odmowa
trwała nie może nim jechać — model dostałby polecenie ponowienia
rozstrzygnięcia, które się nie zmieni. Stąd podział: kanału nie ma albo jest
wygaszony daje `not_found` (jak nieistniejące okno); sufit, dopuszczenie
albo brak siły daje `validation_failed` jako odmowę trwałą, zmienianą przez
Operatora wpisem do rejestru, nie ponowieniem tego samego żądania; pozostałe
przypadki dają `channel_unavailable`. Odmowa jedzie wołającemu z powodem od
doboru, bo „nie skonsultowano" bez zdania dlaczego jest ciszą tam, gdzie
stała decyzja; podmiana odciętego doradcy na innego byłaby tą samą ciszą,
tylko z radą w tle.

## budowa/server/internal/core/handlers_agent_wersje.go

Historia i archiwum eksperta stoją w osobnym porcie, a nie w porcie Agenci:
siedzą w innych tabelach, mają własne repozytoria i wchodzą do rdzenia jednym
wywołaniem `zarejestrujWersjeEksperta`. Wtopienie ich w interfejs Agenci
rozdęłoby port biblioteki o pięć czynności z biblioteką niezwiązanych.

Nazwy komend i kształty pól pochodzą wyłącznie z pakietu `shared`; powielenie
literału nazwy po którejkolwiek stronie jest błędem.

Kształt `AgentVersion` jest ten sam co w `StudioVersion` i `LibraryVersion`:
`id · agentId · label · suma kontrolna · createdAt`. Migawka tożsamości nie
jedzie w wierszu wykazu — wykaz trzyma sumę kontrolną, a treść pobiera się
osobno przy przywróceniu. `mode` jest polem niewymaganym i niesie trzeci stan
promptu: brak wartości znaczy prompt globalny, `DOLACZ` prompt dopisywany,
`ZASTAP` odstępstwo jawne.

Kontrakt daje modułowi Agents wyłącznie zdarzenie `agent.changed`, więc
przywrócenie wersji, archiwizacja i powrót z archiwum rozgłaszają się
rodzajem `updated` wraz z ekspertem po zmianie.

## budowa/server/internal/core/adapter_narzedzia_obraz_model_silniki.go

Wspólne zaplecze (źródło, pracownia, wołanie binarium, odmowy) stoi
w adapter_narzedzia_obraz_model.go, metody stoją na tym samym adapterze. Oba
silniki liczą na procesorze: Real-ESRGAN w wydaniu ncnn-vulkan jest jednym
plikiem wykonywalnym bez Pythona i bez Torcha, a Vulkana dostaje od sterownika
programowego lavapipe z Mesy — wydanie pythonowe dałoby ten sam wynik za cenę
kilku gigabajtów zależności. Rembg jest Pythonem, ale jego runtime ONNX
Runtime ma tryb procesorowy jako podstawowy, nie awaryjny.

Celowo brak gałęzi "gdy silnika nie ma, przeskaluj ImageMagickiem", bo
rozciągnięcie oddane jako powiększenie jest atrapą, której nie widać do
przybliżenia. Celowo brak poprawiania twarzy w tym samym przebiegu: pole
faces niesie kontrakt, ale przebieg twarzowy robi osobna sieć GFPGAN, której
wydanie ncnn tego silnika nie zawiera — idzie ona drugim przebiegiem nad
wynikiem powiększenia. Bez niej żądanie z faces true kończy się odmową
nazywającą brak, zamiast oddać obraz bez poprawki twarzy jako poprawiony.

Nazwa silnika powiększenia dobiera rdzeń, nie model językowy: kontrakt
image.upscale nie ma pola na jej nazwę, ta sama reguła co przy formatach
image.convert.

Wybór modelu wycinania tła jest wyborem ze zbioru, a nie dowolnym tekstem:
u2net jest domyślną siecią ogólną, isnet-general-use bywa dokładniejsza na
cienkim szczególe (włosy, gałęzie), u2net_human_seg jest uczona na ludziach
i na portrecie bije obie, u2netp jest wersją lekką dla maszyn bez zapasu
pamięci.

Pomocnik wycinania tła puszczony na plik pusty odmówiłby dopiero po starcie
interpretera i wczytaniu wag, czyli o minutę później i mniej zrozumiale,
dlatego odczyt idzie wcześniej.

Mapa w Go chodzi losowo, dlatego znaneModeleWycinania wypisuje zbiór w
kolejności stałej, żeby odmowa czytana dwa razy brzmiała tak samo.

## budowa/server/internal/core/handlers_roundtable.go

Cały moduł Roundtable ma jedno zdarzenie: kontrakt daje mu wyłącznie
`roundtable.debate.changed`, więc każda zmiana debaty — otwarcie tury,
wypowiedź uczestnika, interwencja moderatora, zamknięcie tury — rozgłasza się
turą po zmianie, a wypowiedź dołącza jako pole opcjonalne. Cztery okna
odświeżają się z jednej subskrypcji.

Rozgłoszenie nie idzie z obsługiwacza: tura debaty trwa dłużej niż wykonanie
komendy `roundtable.debate.start`, bo uczestnicy odpowiadają równolegle,
każdy we własnym czasie. Zdarzenia nadaje adapter przez podpiętą drogę
rozgłoszenia (tak samo jak w `handlers_terminal.go`); obsługiwacz nadałby
wyłącznie stan sprzed odpowiedzi.

## budowa/server/internal/core/adapter_alerty_miara.go

Każda z ośmiu miar reguł alertu ma źródło i jest liczona w chwili
odpowiedzi: `errorCount` czyta dziennik błędów rdzenia jako liczbę wpisów
w oknie czasu; `errorRate` czyta ślad wywołań modelu jako udział wywołań
nieudanych; `cost` sumuje koszt w oknie czasu ze śladu wywołań; `tokens`
sumuje żetony w oknie czasu ze śladu wywołań; `callLatency` liczy średnie
opóźnienie ważone liczbą wywołań ze śladu; `budgetPercent` porównuje koszt
w oknie wobec pułapu kosztu z konfiguracji; `probeFailure` liczy pomiary
kondycji nieudane w serii pomiarów; `processFailure` liczy pozycje kolejki
rdzenia zakończone błędem. Miara, której źródła rdzeń nie ma wpiętego, jest
odrzucana już przy zapisie reguły (`sprawdzMiareReguly`) — po to, żeby nie
istniała reguła, która wygląda na czynną, a nigdy nie zawoła.

`kluczPulapuKosztu` jest nastawą, wobec której liczy się miara budżetu. Ta
sama nastawa wstrzymuje turę w warstwie kanału (`injection/pulap.go`) — dwa
odczyty jednej wartości, nie dwie wartości.

`wymiarZuzyciaAlertow` jest osią, po której grupuje się ślad wywołań na
potrzeby miar. Kanał, bo każde wywołanie ma kanał; sumy po wszystkich
wierszach dają wartość globalną, a rozbicie zostaje na przyszłe zawężenie
reguły do zasięgu.

`zmierzZeSladu` liczy jeden odczyt na miarę, nie pięć: sumy po wymiarze
niosą komplet potrzebnych składników, a drugie zapytanie o ten sam okres
dałoby liczby z innej chwili.

`rozglosWyzwolenie` nadaje zdarzenie `alert.triggered`. Zdarzenie jest drogą
alertu do okien — bez niego Operator dowiedziałby się o wyzwoleniu dopiero
przy następnym otwarciu wykazu.

Drugi wynik zwracany przez `zmierzMiare` mówi, czy pomiar się odbył — miara,
której nie zmierzono, nie wyzwala niczego, bo nie ma wartości do porównania
z progiem.

## budowa/server/internal/core/adapter_narzedzia_dokument_formaty.go

Trzy nazwy jednej rzeczy, dlatego tabela jest jedna: markdown kontraktu to
markdown dla Pandoca i .md na dysku; txt to plain dla Pandoca i .txt na dysku
— trzy osobne mapy rozjechałyby się przy pierwszym dołożonym formacie. Format
nieznany jest odmową, nie domysłem: wołający dostaje wykaz tego, co rodzina
umie, bo zgadnięty format kończy się plikiem, którego nikt nie otworzy, a
koperta meldowałaby powodzenie.

Pandoc ani nie czyta, ani nie zapisuje PDF-u bez silnika składu, dlatego
czytaPandoc i piszePandoc są dla PDF-u fałszem: zapis PDF-u Pandokiem wymaga
silnika składu wołanego przez niego samego, czyli procesu poza bramą rdzenia.
Czytaniem PDF-u zajmuje się document.text.extract, a zapisem dwie drogi
rdzenia: skład typstem albo LibreOffice.

Rozpoznanie formatu wyłącznie po rozszerzeniu zrywałoby łańcuch najbardziej
w tej rodzinie naturalny: zamień na PDF, a potem przeczytaj, co wyszło — bo
wynik document.convert leży w magazynie jako blob bez kropki i bez
rozszerzenia. Rozpoznanie po nagłówku jest odczytem, nie domysłem po nazwie:
bajty %PDF- na początku pliku są definicją PDF-u, a nie poszlaką. Materiał,
którego nagłówek nie mówi nic pewnego, oddaje pustkę, bo "nie wiem" jest
odpowiedzią uczciwą, a zgadnięty format kończy się odmową narzędzia w połowie
pracy albo, gorzej, treścią przeczytaną nie tym słownikiem.

## budowa/server/internal/core/skutek_bezpieczenstwa_dokumentu_test.go

Rodzina bezpieczeństwa jest tym miejscem, w którym meldunek `status: ok` bez
skutku kosztuje najwięcej: Operator wysyła dokument w przekonaniu, że został
oczyszczony albo zredagowany. Dlatego żaden sprawdzian tutaj nie kończy się
na odpowiedzi komendy — każdy schodzi do bajtów wyniku i pyta je wprost.
Żaden sprawdzian nie pomija się przy braku programu, bo rodzina nie
uruchamia ani jednego procesu.

`trescStronWyniku` skleja treść wszystkich stron dokumentu leżącego pod
odwołaniem zasobu — to jest miara redakcji, bo tekst wycięty znika właśnie
stąd, a nie z odpowiedzi komendy.

`certyfikatProbny` wytwarza parę klucz–certyfikat w postaci PEM. Certyfikat
powstaje na czas sprawdzianu, a nie leży w drzewie: klucz wniesiony do
repozytorium byłby kluczem, który wyciekł w chwili wniesienia.

Miara szyfrowania jest niezależna od odpowiedzi: dokument zaszyfrowany nie
daje się otworzyć bez hasła i to jest jedyny dowód, że szyfrowanie naprawdę
leży w bajtach.

`TestCzyszczenieMetadanychUsuwaOpisZBajtow` wykazuje, że opis znika
z dokumentu, a nie tylko z odpowiedzi.

`TestRedakcjaWycinaTekstZTresciStrony` wykazuje, że redakcja usuwa, a nie
zasłania: po czynności tekst nie leży już w treści strony.

Sam „podpis poprawny" niczego by nie dowiódł w teście przeżycia podpisu przez
odczyt i unieważnienia go zmianą treści — podpis, który zawsze mówi „tak",
nie jest podpisem.

## budowa/server/internal/core/dziennik_cofanie_test.go

Ten plik sprawdzianów wyklucza trzy szkody rachunku cofania z dziennika
czynności: cofnięcie czynności ze środka dziennika, które wgrywa stan sprzed
niej w miejsce dokumentu i zabiera ze sobą wszystko, co po niej weszło —
czyli jest przywróceniem wersji podanym jako cofnięcie pojedyncze; cofnięcie
czynności, na której stoi późniejsza, wykonane po cichu — dokument zostaje
wtedy w stanie niespójnym, a Operator o tym nie wie; i cofnięcie postaci,
które nie działa — pomyłkowa zmiana kroju nie rusza ani jednej litery, więc
treść sprzed jej nie odtworzy.

`TestDziennikCofnieciePostaciNieRuszaLiter` jest sednem wymagania
Właściciela: pomyłkowa zmiana kroju cofa się tak samo jak skasowany akapit,
choć nie ruszyła ani jednej litery.

Stan wpisu dziennika to wyłącznie `active` albo `reverted`, a rodzaj
czynności to `StudioActionKind`, nie `StudioChangeKind`. Podanie jednego
w miejsce drugiego wywraca się dopiero na ograniczeniu tabeli — po zapisaniu
zmiany, czyli w najgorszym możliwym miejscu. Rodzaj „zmiana pola"
(`fieldChange`) doszedł do słownika przy dobudowie, bo bez niego pola
odkładały się jako `objectChange` i nie dawały się cofnąć osobno.

`dziennikZawieraNapis` mówi, czy napis niesie podnapis. Własny pomocnik
z przedrostkiem odcinka, bo sprawdziany tego odcinka pytają o treść odmowy,
a nie o sam fakt odmowy.

## budowa/server/internal/core/adapter_kolejki.go

Kolumnę `kolejka.sesja_id` wypełnia dowiązanie wpięte przez `ZSesjami`:
tabela `sesja` ma kolumnę `identyfikator_zewnetrzny`
(`migracja_006_identyfikatory_zewnetrzne.sql`), repozytorium sesji ma
`PoIdentyfikatorze`, a wiersz sesji powstaje wraz z pierwszą utrwaloną
wiadomością okna. Gdy wiersz sesji istnieje, kolejka zapisuje się z kluczem
obcym, a odczyt kolejki spoza pamięci odtwarza identyfikator sesji z wiersza.
Kolejka założona wcześniej niż pierwsza wiadomość sesji zapisuje się bez
dowiązania, bo wiersza sesji jeszcze nie ma — brak konfiguracji ani brak
wiersza nie mogą odmówić założenia kolejki. Pamięć powiązań zostaje jako
źródło identyfikatora sesji dla takich kolejek.

`nowyAdapterKolejek` wiąże port z repozytorium kolejek i z silnikiem
wykonania pozycji (`kolejka_silnik.go`). Adapter przekłada kontrakt; cykl
życia zlecenia prowadzi silnik i nikt poza nim.

`ZSesjami` dokłada repozytorium sesji, którym adapter dowiązuje kolejkę do
wiersza sesji. Bez niego adapter pracuje jak dotąd — dowiązanie jest
dodatkiem, nie warunkiem.

`ZWykonawcaModelu`: pozycję posuwa wtedy wyłącznie działanie Operatora,
pętli albo MultitaskingAI. Rejestr kanałów i nadajnik strumienia idą tą samą
drogą co tura okna.

`UstawUjscieWyniku` woła adapter podagentów przy własnym montażu
(`zWykonaniem`), więc zapis biegnie przed obsługą pierwszego żądania — pole
silnika jest wartościowe i podmienia się w miejscu, dokładnie jak przy
`ZWykonawcaModelu`.

`rozwiazKanalPozycji`: ani schemat `pozycja_kolejki`, ani kontrakt nie mają
pola kanału. Brak kanału czynnego znaczy „nie ma czym wykonać kroku" —
wykonawca odda błąd, a silnik pokaże pozycję jako `bledna`, zamiast udawać
wykonanie.

`ZTelemetria` dokłada producenta telemetrii postępu. Kolejka jest procesem:
ma etapy równe pozycjom, stan i bieg naprawczy bez limitu obiegów, więc
zasila Process Monitor tym samym zdarzeniem, co tura okna.

`Utworz` czyta wykaz zleceń przed założeniem kolejki: ładunek uszkodzony ma
zakończyć się odmową, a nie kolejką założoną w połowie.

`wierszSesji`: brak repozytorium, brak wiersza i błąd odczytu dają kolejkę
bez dowiązania — założenie kolejki nie może zależeć od tego, czy sesja
zdążyła się utrwalić.

`sesjaZWiersza` odtwarza identyfikator sesji rdzenia z dowiązanego wiersza —
tą drogą kolejka zapisana przed ponownym uruchomieniem rdzenia wraca
z sesją, choć pamięć powiązań jest już pusta.

`Wykonaj`: powtórzenie kroku nie ma limitu obiegów — adapter niczego nie
zlicza i niczego nie odmawia.

`etapDzialania` nazywa etap telemetrii odpowiadający działaniu na kolejce.
Powtórzenie kroku jest biegiem naprawczym i tak też się nazywa.

`kolejkaKontraktu`: pole `Queue.Cycle` opisuje bieg naprawczy zlecenia,
a ten liczy się na pozycji, nie na kolejce.

`Kolejka` oddaje kolejkę kontraktu po jej identyfikatorze. Służy czynnościom
rodziny `queue.item.*`, których wynikiem jest kolejka po zmianie, a nie samo
zlecenie — drugiego składania kolejki w rdzeniu nie ma.

## budowa/server/internal/core/adapter_zespolow.go

Zespół jest składem, a nie drugą biblioteką ekspertów: adapter przenosi
wyłącznie kody, którymi posługuje się kontrakt (`Agent.id` = `agent.kod`),
i nie sprawdza tożsamości eksperta. Dostępność eksperta ustala warstwa
`dane` jednym złączeniem.

Brak eksperta nie przerywa wczytania: ekspert usunięty albo zarchiwizowany
wypada ze składu oddawanego klientowi, a zespół wraca kompletny w pozostałej
części. Fakt pominięcia trafia do dziennika rdzenia, ponieważ struktura
`Team` kontraktu nie ma pola na taką wiadomość, a milczenie ukryłoby zmianę
składu.

Kopia nie zmienia źródła: `team.duplicate` czyta zespół źródłowy i zakłada
obok nowy wiersz z tym samym składem; źródło nie przechodzi ani przez
`Zapisz`, ani przez żadne inne polecenie zmieniające.

`dolozStanZastany` uzupełnia zapis zespołu o to, czego żądanie nieść nie
mogło. Opis: `description` jest w kontrakcie polem nieobowiązkowym, a brak
pola znaczy „nie zmieniaj opisu" — tak zapowiada panel składu zespołu i tak
zachowuje się klient, który pustego napisu nie wysyła; opis zmienia się
wtedy i tylko wtedy, gdy pole przyszło, także wypełnione pustką. Skład:
odczyt zespołu okrawa skład z ekspertów zarchiwizowanych i usuniętych, więc
żądanie złożone z tego, co widać w oknie, ich nie niesie; adapter dokleja
ich z powrotem — tak samo jak kopia zespołu (`Skopiuj`) — dzięki czemu
wracają do składu po przywróceniu do biblioteki, a ekspert widoczny
i wypisany ze składu wypada normalnie.

## budowa/server/internal/core/zapora_procesow_rdzenia_test.go

Zapora procesów CAŁEGO rdzenia — trzecia po zaporze warsztatu dokumentu
(`zapora_warsztatu_pdf_test.go`) i zaporze fotografii
(`zapora_fotografii_test.go`), i pierwsza, która nie ogranicza się do jednego
obszaru.

Dlaczego powstała i czego nie powtarza: dwie starsze zapory zabraniają procesu
w swoich obszarach, bo tam każdą pracę wykonuje w całości biblioteka Go. Tej
reguły nie da się rozciągnąć na rdzeń bez kłamstwa: rozpoznanie pisma, archiwa,
mowa, nagrania i silniki neuronowe nie mają biblioteki czysto-Go, więc ich
programy są składnikiem pakietu serwera i wolno je wołać. Granicą nie jest
„proces czy biblioteka", a: 1. czy program jest wołany JEDYNĄ dozwoloną drogą
(`zewnetrzne.Wolaj`), która sprawdza obecność, obejmuje drzewo procesów
i odmawia zdaniem nazywającym brak — a nie własnym `exec.Command`, który
żadnej z tych rzeczy nie robi; 2. czy pakiet serwera ten program NIESIE, czyli
czy stoi w wykazie zależności (`zaleznosci_zewnetrzne.go`). Program wołany bez
wpisu w wykazie to cichy wymóg wobec wdrożenia: u Operatora funkcja odmawia,
a przy starcie nikt nie powiedział, czego brakuje. Trzecia rzecz, której
zapora pilnuje, jest odwrotna do dwóch pierwszych: rachunek, który JUŻ jest
wkompilowany, nie ma prawa wrócić do procesu. Cztery czynności obrazu modelu
(`image.inspect`, `image.transform`, `image.adjust`, `image.convert`) liczyły
się kiedyś programem, choć biblioteka wystarcza — i właśnie dlatego ta zapora
wymienia ich pliki po nazwie.

plikiWlasnegoExecuUzasadnione: klucz jest ścieżką liczoną od `internal/`, bo
zapora obejmuje cały ten katalog, a sama nazwa pliku nie mówi, w którym
pakiecie wywołanie stoi. Wykaz jest ZAMKNIĘTY: dopisanie do niego pliku wymaga
powodu tej samej wagi, co pozycje poniżej, i jest zmianą rozstrzygnięcia, nie
porządkowaniem listy. `core/adapter_modul_extension_protokol.go`
i `core/adapter_modul_extension_integracje.go` uruchamiają program WSKAZANY
PRZEZ ROZSZERZENIE, a nie program produktu. `zewnetrzne.Wolaj` sprawdza
obecność narzędzia z deklaracji rdzenia — tu deklaracji nie ma, bo program
przychodzi z manifestu rozszerzenia. `core/adapter_modul_developer_narzedzia.go`
pyta narzędzie o jego własną wersję ścieżką już rozwiązaną; to sonda obecności,
nie czynność Operatora. `injection/rozruch.go` JEST drogą, do której zapora
odsyła: `zewnetrzne.Wolaj` idzie portem `session.Uruchamiacz`, a ten port
kończy się tutaj. Pozycja nie jest wyjątkiem od reguły, tylko jej dnem — bez
niej reguła nie ma się o co oprzeć, a `zewnetrzne/wolanie.go` nazywa ten plik
jedynym `exec.Command` w drzewie. `zdalne/pliki.go` przenosi plik programem
`scp` i nie ma dziś drzwi, przez które mógłby przejść. Arsenał
(`zewnetrzne.Wolaj`) zbiera całe wyjście do pamięci pod obowiązkową granicą
czasu — przenosiny pliku dowolnego rozmiaru nie mają uczciwej granicy, a ich
wyjście nie jest wynikiem do zebrania. Spawner platformy stoi po stronie
kanału, a zależność biegnie od kanału do toru i nigdy odwrotnie
(`zdalne/polecenie.go`), więc pakiet `zdalne` go nie zaimportuje. Pozycja stoi
tu po to, żeby brak drzwi był widoczny zamiast niewidoczny — nagłówek
`zdalne/zdalne.go` głosi, że pakiet procesów nie uruchamia, a ten plik je
uruchamia.

TestCzteryCzynnosciObrazuLiczaSieWkompilowane: mierzy dwie rzeczy naraz, że
wołanie programu stoi wyłącznie w dwóch funkcjach drogi zapasowej i że każda
z czterech czynności przechodzi przez rachunek wkompilowany. Sprawdzian liczy
funkcje, a nie samą liczbę wystąpień, bo trzecie wołanie dopisane do
istniejącej funkcji jest tą samą szkodą.

korzenZapory: reguła jednej drogi do procesu jest regułą rdzenia, a nie regułą
jednego pakietu: wywołanie przeniesione o katalog dalej wychodzi spod niej,
choć szkodę robi tę samą. Zapora czytająca własny katalog nie widziałaby ani
przenosin plików torem zdalnym, ani żadnego następnego wywołania spoza tego
katalogu.

TestRdzenUruchamiaProcesyJednaDroga: własny `exec.Command` pomija trzy rzeczy
naraz: sprawdzenie obecności programu, bramę izolacji okna i objęcie drzewa
procesów. Pierwsza zamienia brak programu w niezrozumiały błąd zamiast zdania
nazywającego brak, druga wypuszcza czynność poza zasięg okna, trzecia zostawia
sieroty po granicy czasu.

TestProgramyRdzeniaStojaWWykazieZaleznosci: to jest właściwa miara gotowości
funkcji opartej o program: nie „czy w kodzie jest exec", a „czy pakiet serwera
to niesie i czy Operator dowie się o braku przy starcie". Program wołany bez
wpisu w wykazie jest cichym wymogiem wobec wdrożenia — na maszynie
deweloperskiej, gdzie ktoś doinstalował go ręcznie, wygląda jak funkcja gotowa.

## budowa/server/internal/core/skutek_adnotacji_studia_test.go

Szkoda, którą ten plik ma wykluczyć: decyzja bez skutku. Przyjęcie zmiany
śledzonej albo propozycji, które przestawia sam znacznik wiersza i zostawia
treść dokumentu nietkniętą, wraca kopertą `ok` i wygląda jak praca wykonana —
a Operator po otwarciu dokumentu widzi tekst sprzed decyzji.

`trescDokumentu` odczytuje treść dokumentu osobnym wywołaniem — to jest
miara właściwa, bo Operator otworzy dokument ponownie, a nie przeczyta
odpowiedź.

## budowa/server/internal/core/adapter_rozmowa_zestaw.go

Rdzeń nie zna nazw narzędzi eksperta i znać ich nie ma. Przełożenie kodu
z definicji eksperta na pozycje wykazu (albo na całą ich grupę) należy do
serwera narzędzi — `narzedzia/ekspert_wykaz.go` — bo tylko on trzyma wykaz
kontraktu. Gdyby rdzeń wyliczał listę sam, powstałby drugi wykaz obok
kontraktowego. Stąd wychodzi więc wskazanie: kod eksperta i nazwy dołożeń;
rozwinięcie wskazania w wykaz robi strona przeciwna.

Konfigurację MCP składa `mostyOkna.tekstZNarzedziami` — jedyna droga, którą
jedzie rozmowa. Ten plik jej nie powtarza: bierze jej wynik i dokłada do
gotowego wpisu `danaco` argumenty zestawu. To ten sam wzorzec, którym okno
asystenta dokłada sobie zasięg klawiatury (`zZasiegiemKlawiatury`
w `adapter_modul_asystent_sterowanie.go`). Tekst nieczytelny albo wpis
`danaco` nieobecny (brak binarium serwera narzędzi — `most_narzedzi.go`
melduje ten powód osobno) zostaje nietknięty: tura idzie z zestawem, jaki
jest, zamiast paść na składaniu konfiguracji.

Dokładanie argumentów jest tu wybiórcze, a nie hurtowe, i każdy nowy
argument wymaga umowy, a nie założenia, ponieważ `cmd/danaco-narzedzia/main.go`
woła `flag.Parse()` na domyślnym `flag.CommandLine`, czyli z `ExitOnError`:
przełącznik, którego tamta strona nie zna, kończy proces serwera narzędzi.
Skutek nie jest wtedy „tura bez kilku narzędzi", tylko „tura bez wykazu
w całości" — i wygląda dla Operatora jak awaria modelu, nie jak rozjazd
dwóch pakietów.

Oba argumenty składane w tym pliku mają odczyt po drugiej stronie:
`--zasieg ekspert --ekspert <kod>` w `narzedzia/zasieg_eksperta.go`,
`ekspert_definicja.go` i `ekspert_wykaz.go`; `--dolozenia <lista>`
w `narzedzia.RozbijDolozenia` i `WykazEksperta.ZDolozeniami` — wartość
składa `injection.PrzelacznikDolozen`, czyli jedna definicja nazwy po obu
stronach.

Brak źródła dołożeń nie odbiera modelowi narzędzi: tura jedzie podstawą,
jaką ma, i wpisuje do dziennika rdzenia powód. Meldunek idzie raz na powód,
nie raz na turę — tak samo jak meldunek o braku binarium serwera narzędzi
(`most_narzedzi.go`) i z tego samego powodu: tur bywa kilkaset dziennie,
a powód się między nimi nie zmienia.

Port `DolozeniaNarzedziSesji` jest wąski z rozmysłu: pyta o jedno i oddaje
jedno. Sesja wskazana jest identyfikatorem kontraktowym, tym samym, który
niesie `session.Okno.IdSesji` — przekład na klucz wiersza należy do strony
odpowiadającej, tak jak przy każdym innym repozytorium. Dołożenia są listą
jednorodną i nie niosą znacznika zawężenia, bo dołożenie z natury dokłada:
nie potrafi zawęzić czegoś, czego samo nie ustanowiło.

Drugiego wskazania eksperta w `zestawTury` nie ma i mieć nie może: dwa
wskazania to dwie prawdy o tym, kim tura jest.

`argumentyZestawu` niesie całą decyzję tego pliku: co wolno dołożyć,
a czego nie wolno, i dlaczego. Przepisanie tekstu niżej jest już tylko
mechaniką. Serwer narzędzi ma dziś jedno miejsce, w którym dołożenie może
dojść do wykazu — złożenie wykazu eksperta (`narzedzia/ekspert_wykaz.go`).
Okno bez eksperta nie przechodzi tą drogą wcale, więc argument dołożeń nie
miałby tam czego dołożyć; cichy odrzut byłby tu gorszy niż brak, bo
Operator widziałby narzędzie na wykazie sesji i nie widziałby go w turze.

Zasięgu roli funkcja `zZestawemTury` nie tyka. Rola okna (`--zasieg
klawiatura`) jedzie osobną drogą okna asystenta i dokłada, a nie zawęża;
zasięg eksperta zawęża i wchodzi tutaj. Gdyby oba spotkały się kiedyś na
jednym wpisie, rozstrzyga `narzedzia.RozpoznajZasieg`, czytając ostatnie
wystąpienie przełącznika — i to jest rozstrzygnięcie tamtej strony, nie tej.

`zKonfiguracjaZestawu` jest całą drogą złożoną w jedno wywołanie, żeby
miejsce wpięcia w `adapter_rozmowa_srodowisko.go` pozostało jedną linią.

Kluczem dziennika `zglosZestaw` jest sam powód, bez identyfikatorów okna
i sesji: te idą do treści meldunku, ale nie do klucza — inaczej każde nowe
okno wywoływałoby ten sam meldunek od nowa. Dziennikiem jest dziennik
składacza mostów — ten sam, do którego idzie meldunek o braku binarium
serwera narzędzi. Jedna sprawa, jedno miejsce. Dziennik niewskazany nie
zmienia przebiegu tury; znika wyłącznie meldunek.

## budowa/server/internal/core/skutek_narzedzi_obrazu_test.go

Dlaczego ten sprawdzian nie pomija się przy braku programu: te cztery
czynności liczył wcześniej program pakietu serwera i sprawdzian skutku
wymagałby jego obecności — na maszynie bez niego świeciłby na zielono jako
„pominięty", czyli nie mierzyłby niczego. Rachunek stoi teraz wkompilowany
w binarium (`adapter_narzedzia_obraz_wkompilowany.go`), więc sprawdzian idzie
zawsze i mierzy PIKSELE wyniku, a nie pola odpowiedzi. Program zostaje
wyłącznie drogą zapasową dla AVIF-a i WEBP-a stratnego — tych dwóch wyjść ten
plik nie mierzy, bo nie ma czym: kodera czysto-Go dla nich nie ma i dlatego
właśnie tamta droga istnieje.

## budowa/server/internal/core/kolejka_silnik.go

Protokół działań: kontrakt daje sześć działań (`QueueAction`) i nie ma
osobnego działania „zamknij pozycję werdyktem". Silnik czyta więc działania
dosłownie tak, jak są nazwane, i posuwa pozycje po tabeli przejść
`krokNaprzod`: `start`/`resume` to krok naprzód (oczekuje → wykonywana →
do_weryfikacji → ukonczona); `retry` to bieg naprawczy (licznik obiegów +1,
werdykt do_poprawy, powrót do wykonywana, bez limitu i bez warunku); `stop`
to przerwanie (pozycja wskazana albo wszystkie czynne → anulowana); `pause`
to wstrzymanie kolejki (pozycje zostają, gdzie były); `clear` to opróżnienie
(pozycje czynne → anulowana, dziennik zostaje — przejrzystość zamiast
kasowania śladu).

Przyjęcie wyniku kroku wyraża się wyborem działania: krok naprzód znaczy
przyjęcie, retry znaczy odesłanie do poprawy. Silnik nie wystawia werdyktu,
którego nie wywołało działanie Operatora, i nie zmyśla postępu. Działanie
wskazujące pozycję (`itemId`) dotyczy tej pozycji; bez wskazania dotyczy
pozycji, na której kolejka stoi.

Most do realnego wykonania: sam przebieg stanów nie wykonuje pracy — pozycja
wchodząca w stan `wykonywana` musi zostać naprawdę wykonana, a jej wynik,
nie klik Operatora, przesuwa ją dalej — powodzenie do `do_weryfikacji`,
niepowodzenie do `bledna`, czyli do realnego stanu błędu zamiast cichego
ukończenia. Robi to wpięty `wykonawca`. Silnik bez wykonawcy zostaje czystą
maszyną stanów: pozycję posuwa wtedy działanie Operatora, pętli sesyjnej
albo MultitaskingAI, tak jak przed wpięciem mostu.

`ZWykonawca` zwraca silnik przez wartość, bo `silnikKolejki` trzymany jest
w adapterze kolejek jako pole wartościowe, a nie wskaźnik — montaż podmienia
je w miejscu.

`ZUjsciemWyniku` wpina odbiorcę zebranej treści tury. Silnik nie wie, kto
odbiera — dziś jest to wiersz podagenta (`adapter_modul_orkiestracja.go`),
ale silnik zna wyłącznie pozycję; pozycja bez odbiorcy przechodzi bez śladu
treści.

`Zasil` zakłada zlecenia początkowe kolejki. Wykaz pusty zostawia kolejkę
bez pozycji — to poprawny stan, nie awaria.

Stan kolejki oddawany przez `Wykonaj` jest wyprowadzony z pozycji, nie
zadeklarowany: dopóki jest co robić, kolejka pracuje; gdy nie ma — jest
wyczerpana.

`Pozycje` zwraca zlecenia kolejki. Błąd odczytu daje wykaz pusty —
odpowiedź o kolejce nie ma znikać z powodu jednego zapytania pobocznego.

`Cykl` podaje licznik obiegów pozycji, na której stoi kolejka — pole
`Queue.Cycle` kontraktu. Kolejka wyczerpana pokazuje licznik pozycji
ostatniej, kolejka pusta nie pokazuje żadnego.

Wejście w stan `wykonywana` nie kończy kroku `krok`: pozycja jest wtedy
naprawdę wykonywana, a jej wynik przesuwa ją dalej. Zawrócona pozycja
`biegNaprawczy` wchodzi w `wykonywana`, więc jest wykonywana od nowa —
powtórzenie kroku ma powtórzyć pracę, nie samo przełożyć etykietę stanu.

`queue.action stop` wydane w trakcie tury zapisuje pozycji stan `anulowana`.
Bez sprawdzenia w `domknijPoTurze` zapis nadpisałby go chwilę później stanem
`do_weryfikacji`, jak gdyby przerwania nie było — łamiąc własny protokół
silnika („stop · przerwanie: pozycja … → anulowana") i zamykając jedyną
kontraktową drogę zatrzymania pracy podagenta, bo komendy `subagent.stop`
kontrakt nie ma. Przerwanie ma pierwszeństwo przed spóźnionym werdyktem
tury; sama treść odpowiedzi trafiła już do ujścia wyniku, więc nic z pracy
nie ginie po cichu.

Odczyt w `domknijPoTurze` idzie listą pozycji kolejki, bo repozytorium nie
ma odczytu jednej pozycji — a dorabianie go dla tego jednego miejsca byłoby
drugim zapytaniem o to samo. Pozycja nieodnaleziona przechodzi na zapis
wprost: lepiej zapisać stan wynikający z tury, niż zgubić go z powodu błędu
odczytu.

`anuluj` zamyka pozycje przerwaniem. Wiersz zostaje razem z dziennikiem —
ślad przerwanego zlecenia jest częścią przejrzystości pętli.

TestObrotIdzieZgodnieZeWskazowkamiZegara: kierunek odwrotny daje obraz
o tych samych wymiarach, więc pomiar samych boków by go nie zauważył —
mierzymy więc, gdzie wylądowała czerwona połowa. Obrót o dziewięćdziesiąt
stopni zgodnie z ruchem wskazówek zegara przenosi lewą połowę na górę. Obrót
przeciwny położyłby ją na dole.

TestKonwersjaDoWebpBezstratnegoNieWymagaProgramu: sprawdzian dekoduje wynik
bibliotecznym czytnikiem WEBP: gdyby zapis szedł programem, na maszynie bez
niego czynność by odmówiła, a tu ma przejść zawsze.

TestKonwersjaDoAvifBezProgramuOdmawiaNazywajacBrak: pustą ścieżką wyszukiwania
sprawdzian czyni z tej maszyny maszynę bez programu, więc odmowę mierzy każdy
bieg, nie tylko bieg na cienkiej instalce.

## budowa/server/internal/core/skutek_migawki_przegladarki_test.go

Migawka jest tym, co model i Operator widzą zamiast strony. Szkoda, którą ten
plik ma wykluczyć, polegała na uciszeniu przekroczenia rozmiaru: rdzeń czytał
tyle, ile mieściła granica, i podawał ucięty początek jako pełną treść strony.
Odpowiedź była udana, migawka istniała, a model wnioskował z połowy dokumentu,
nie wiedząc, że to połowa.

Dlatego sprawdziany tego pliku nie pytają, czy migawka powstała. Pytają, czy
niesie koniec strony, a przy stronie ponad granicą — czy rdzeń odmówił i nie
zostawił po sobie migawki, którą `browser.snapshot.get` podałby dalej jako
bieżący stan strony.

Zrzut ekranu wymaga silnika przeglądarki spoza biblioteki standardowej, więc
rdzeń go nie wypełnia. To nie jest brak do zmierzenia sprawdzianem skutku,
tylko granica nazwana wprost: pole `screenshotRef` ma zostać puste, i to
właśnie sprawdzian niżej stwierdza — pustka jest tu prawdą, a odsyłacz
wskazujący nic byłby drugą postacią tej samej szkody.

TestMigawkaNiesieKoniecPobranejStronyANieJejPoczatek: sprawdzenie długości nie
wystarczyłoby — ucięcie zawsze daje jakąś długość — więc miarą jest znacznik
stojący na samym końcu ciała.

TestStronaPonadGranicaRozmiaruNieZostawiaMigawkiOgryzka: migawka ucięta, raz
zapisana, jest odtąd podawana przez `browser.snapshot.get` jako bieżący stan
strony i nic już nie mówi o tym, że jest połową.

TestStronaDeklarujacaRozmiarPonadGranicaNieJestPobierana: gdy witryna sama
zapowiada rozmiar większy niż granica, rdzeń nie ma po co ciągnąć ani bajta.

TestMigawkaOddajeZrodloStronyDopieroNaZadanie: brak HTML-a przy `includeHtml`
niewskazanym jest oszczędnością, nie brakiem treści — źródło ucięte byłoby tą
samą szkodą co tekst ucięty, tylko w drugim polu.

TestSnapshotGetOddajeMigawkeNajswiezszegoPrzejscia: wydanie migawki poprzedniej
pokazywałoby modelowi stronę, z której Operator już wyszedł.

## budowa/server/internal/core/handlers_automatyka_petla.go

Wybudzenie to nie pauza. Pauza czeka na czas (krok `wait`) albo na Operatora
(stan `paused`) — w obu razach wiadomo, kiedy bieg ruszy. Wybudzenie czeka na
świat: na reakcję, która może przyjść za godzinę, za trzy dni albo nigdy.

Bieg czekający wiecznie jest wyciekiem, dlatego oczekiwanie z terminem
trafia pod zegar, a `po_terminie` mówi, co zrobić, gdy termin minie: `wznow`
rusza dalej tak, jakby sygnał przyszedł — domyślne, bo produkt ma pracować
dalej, a nie stawać; `ponow` wykonuje krok oczekiwania jeszcze raz, `proba`
rośnie; `przerwij` kończy bieg stanem `stopped` z jawnym powodem.

Oczekiwanie bez terminu pod zegar nie trafia: bywają reakcje, na które czeka
się bez zegara, więc taki bieg czeka, aż przyjdzie sygnał albo aż Operator
go zamknie.

Sygnał doręczony biegowi, który nie czeka, nie ma adresata — to krótsza
lista odbiorców, a nie odmowa.

Interfejs `repozytoriumWybudzen` jest składany w tym pliku, a nie w `dane`,
bo łączy dwa zakresy w jeden widok jednego odbiorcy — repozytorium
automatyk spełnia go w całości.

`nowySilnikWybudzen`: repozytorium, które nie niesie bytów pętli, daje
silnik pusty zamiast awarii — rdzeń wstaje, automatyka pracuje, a wybudzeń
po prostu nie ma.

`Wybudz` doręcza sygnał wszystkim biegom, które na niego czekają, i zwraca
liczbę wybudzonych. Sygnał bez adresata nie jest błędem: świat zewnętrzny
nie wie, które biegi czekają, więc zero wybudzonych to poprawna odpowiedź.

`ZamknijRecznie` wybudza bieg na żądanie operatora, bez czekania na sygnał
ani na termin — trzecia, obok sygnału i terminu, droga wyjścia biegu ze
stanu oczekiwania.

`rozstrzygnijTermin` stosuje politykę zapisaną przy zakładaniu oczekiwania.
Awaria jednego biegu nie zatrzymuje zegara — pętla idzie dalej.

`wznowZProba` wznawia bieg, licząc próbę. `proba` nie ma górnej granicy:
licznik rośnie i jest widoczny, ale sam nie zatrzymuje pracy.

`przerwij` kończy bieg, zostawiając powód. Komunikat nazywa sygnał i krok —
bez tego Operator zobaczyłby bieg przerwany bez przyczyny.

## budowa/server/internal/core/handlers_extension.go

Schemat katalogu rozszerzeń leży w
`store/migracja_070_katalog_rozszerzen.sql`. Rodzina niesie pięć komend:
`extension.list`, `extension.install`, `extension.configure`,
`extension.toggle` i `extension.uninstall`. Pozycja `extension.unknown`
o tym przedrostku jest odpowiedzią na komendę nieznaną obszaru, a nie
komendą: nie ma pary żądanie/wynik i nie rejestruje się jej w rejestrze
komend.

Rozgłoszenie `extension.changed` jest wpięte po stronie adaptera:
`adapter_modul_extension_rozgloszenie.go` niesie dokładkę `ZRozgloszeniem`,
a cztery komendy wołają ją po udanym zapisie — `extension.install`
(`created`, także przy przywróceniu pozycji odinstalowanej),
`extension.configure` i `extension.toggle` (`updated`), `extension.uninstall`
(`deleted`). `extension.list` nie rozgłasza niczego. Gdy port złożono bez
nadajnika (`montaz_porty.go`), `a.rozgloszenie` jest nilem i cztery wywołania
milkną.

Port niewypełniony nie rejestruje niczego: komendy odpowiadają wtedy
`extension.unknown`, a pozostałe domeny pracują bez zmian.

Jeden port na całą rodzinę: pięć komend obsługuje jeden byt — pozycję
katalogu — i jedną maszynerię. Osobny port „instalacji" obok portu
„katalogu" byłby dwiema prawdami o jednej tabeli.

## budowa/server/internal/core/zgodnosc_kontraktu_test.go

Zgodność rdzenia z kontraktem. Kontrakt jest jedynym źródłem prawdy nazw, ale
sam z siebie niczego nie wymusza: nazwa może stać w `contract.json`,
wygenerować się do `contract.go` i nie mieć po stronie rdzenia ani jednego
obsługiwacza. Kompilacja tego nie wychwyci — stała jest użyta w kontrakcie,
więc nie jest martwa. Grep tego też nie rozstrzygnie. Część rejestracji idzie
przez zmienną (`r.Zarejestruj(n.Przejecie, …)`, `r.Zarejestruj(nazwa,
obsluga)`), więc wyliczenie literałów w źródle zaniża wynik i nie wiadomo o ile.
Jedynym pomiarem, który mówi prawdę, jest zmontowany rejestr.

granicaKomendySprawdzianu: wartość wychodzi z granicy warstwy, nie z czasu
pomiaru. Najdłuższa czynność mierzona uprzężą jest czynnością skanera i sama
stoi pod granicą `granicaWykazuUrzadzen` (45 s, `urzadzenia_skaner.go`):
komenda, której urządzenie nie odpowiada, wraca odmową dopiero po tym czasie.
Uprząż ciaśniejsza od tej granicy urywa komendę przed jej własną odmową
i melduje usterkę rdzenia tam, gdzie zwisło urządzenie — tak chwiał się
sprawdzian przy granicy 15 s, podczas gdy czynność skanera dochodzi na tej
maszynie do ~14,8 s. Zapas ponad granicę warstwy to 15 s: tyle trwa montaż
rdzenia i droga koperty wokół samej czynności, a jest to zarazem
czterokrotność najdłuższego zmierzonego wywołania. Granica pozostaje o rząd
wielkości niższa od granicy pojedynczego przebiegu skanera (5 min), więc zwis
rdzenia nadal wychodzi w minutach, nie w godzinach.

komendyBezObslugiwacza: wykaz jest zaporą, nie zgodą: sprawdzian wypada
niepomyślnie zarówno wtedy, gdy pojawi się brak spoza wykazu, jak i wtedy, gdy
brak z wykazu zostanie uzupełniony, a wiersz zostanie. Dług nie rośnie po cichu
i nie znika po cichu. Klient nie zobaczy tych komend w powitaniu, bo powitanie
oddaje wykaz z rejestru rdzenia. Wołanie ich wraca zdarzeniem `*.unknown`
z kodem `not_found` — odmową nazwaną, nie zerwaniem połączenia. Wykaz jest
dziś PUSTY: każda komenda kontraktu ma w rdzeniu obsługiwacza. Pustego wykazu
nie zwijamy do usunięcia zmiennej — obie zapory niżej stoją na niej i mają
działać dalej, a wiersz dopisany tu w przyszłości ma być decyzją widoczną
w przeglądzie, nie skutkiem ubocznym.

TestRejestrPokrywaKomendyKontraktu: komenda bez obsługiwacza nie jest błędem
zrywającym, ale jest funkcją zapowiedzianą i niedostarczoną, czyli dokładnie
tym, czego wykaz braków nie widzi.

TestRejestrNieMaNazwSpozaKontraktu: nazwa taka byłaby funkcją
nieudokumentowaną — klient nie miałby jak jej wywołać, bo bindingi powstają
wyłącznie z kontraktu.

TestPowitanieOddajeWykazZRejestru: rozjazd tych dwóch zbiorów oznacza, że
klient odblokowuje okna funkcji, których rdzeń nie ma — albo ukrywa te, które
ma.

TestKomendaSpozaKontraktuWracaJakoNieznana: połączenie nie jest zrywane —
sprawdzian dowodzi tego wywołaniem kolejnej komendy po odmowie.

powitanieSprawdzianu: powitanie jest tu narzędziem, nie przedmiotem pomiaru —
oba sprawdziany powyżej pytają o wykaz komend i o to, czy rdzeń pracuje po
odmowie. Treść niepełna mierzyłaby w tym miejscu bramę kontraktu zamiast tego,
o co sprawdzianom idzie.

TestKazdaKomendaZnosiPustyLadunek: sprawdzian nie ocenia treści odpowiedzi —
ocenia, że obsługiwacz nie przerywa wykonania i że odmowa jest odmową nazwaną:
koperta ze stanem i kodem kontraktu. Ładunek pusty jest tu przypadkiem
granicznym najtańszym do wywołania i najczęstszym w praktyce: tak wygląda
żądanie klienta z niewypełnionym formularzem oraz wywołanie narzędzia przez
model, który pominął parametr.

## budowa/server/internal/core/skutek_centrum_powiadomien_test.go

Uprząż sprawdzianów jest ta sama, co dla pozostałych sprawdzianów skutku
(`zmontujDoPomiaruSkutku`): świeża baza, pełny montaż, komendy przez
rejestr.

Waga zdarzenia bierze się z taksonomii klasy, a nie ze zgłoszenia —
zgłoszenie jej nie podało.

Klasa „zakończenie" ma w katalogu kanał Mobile włączony domyślnie, więc
zdarzenie idzie obiema drogami; nastawy rozstrzygają o kanale, nie kod
adaptera.

Klasa „wzmianka" nie ma kanału dodatkowego w katalogu, więc zdarzenie
zostaje w samym centrum — sprawdzian odłożenia nie dotyka kolejki doręczeń.

## budowa/server/internal/core/brama_kontraktu_test.go

Straże bramy kontraktu: rdzeń wykonuje żądanie dopiero po sprawdzeniu go wobec
kontraktu. Sprawdziany w tym pliku pilnują dwóch rzeczy naraz, bo obie znoszą
się nawzajem: 1. żądanie niezgodne z kontraktem ma wrócić odmową, a nie
powodzeniem z wartością dobraną przez rdzeń; 2. odmowa ma nazwać pole, o które
idzie — odmowa mówiąca „coś jest nie tak" zostawia wołającego dokładnie tam,
gdzie zostawiało go milczenie. Ładunki są tu wypisane mapą, nie strukturą
kontraktu: struktura Go niesie pole wymagane zawsze (znacznik bez
`omitempty`), więc żądania bez pola nie da się nią złożyć. Żądanie niepełne
przychodzi z drutu i tylko mapą da się je odtworzyć.

TestBrakPolaWymaganegoNazywaPoleNieUsterkeWewnetrzna: pilnuje czterech dróg, na
których brak wartości wracał jako `internal_error`: wołający dostawał tekst
zapytania SQL albo zdanie o kolumnie bazy zamiast nazwy pola, którego nie
wypełnił, a kod błędu kazał mu ponowić żądanie, które nie ma prawa się udać.

TestPustaWartoscPolaWymaganegoNazywaPole: pole obecne, lecz puste, przechodzi
bramę kontraktu — kontrakt żąda obecności pola, nie jego niepustości — i
rozstrzyga o nim dziedzina. Bez tego sprawdzenia pusty rodzaj kanału dojeżdżał
do więzu schematu i wracał treścią zapytania SQL.

TestPowitanieNiepelnePrzechodziBrameIOddajeWersjeProtokolu pilnuje jedynego
wyjątku spod bramy (wpis rejestru decyzji). Czym się to łamie: brama
objęła `connection.hello`, którego trzy pola kontrakt oznacza jako wymagane.
Klient sprzed wprowadzenia pola `clientId` dostawał więc odmowę zamiast wersji
protokołu — a powitanie jest jedynym miejscem, z którego klient tę wersję
czyta, czyli jedynym, w którym rozpoznaje, że jest starszy niż rdzeń. Im
starszy klient, tym pewniej nie dowiadywał się, dlaczego został odrzucony.
Ładunek idzie mapą pustą, bo struktura kontraktu niesie pola wymagane zawsze
i żądania bez nich nie da się nią złożyć — a niepełne powitanie przychodzi
właśnie z drutu.

TestBrakiPowitaniaIdaDoDziennikaRdzenia pilnuje drugiej połowy pozycji 10:
braki pól są odnotowane, a nie przemilczane. Dziennik jest jedynym miejscem,
do którego mogą dojść: `ConnectionHelloResponse` nie ma pola na wykaz braków,
a dołożenie takiego pola jest zmianą kontraktu. Sprawdzian mierzy więc zapis,
a nie treść odpowiedzi — i tym samym pilnuje, żeby wyjątek nie zamienił się
w milczenie.

TestWyjatekObejmujeWylaczniePowitanie pilnuje granicy wyjątku od drugiej
strony: sąsiadka powitania w tej samej rodzinie kontraktu przechodzi bramę bez
ustępstw i odmawia, nazywając brakujące pole. Bez tego sprawdzianu wyjątek
dopisany dla powitania mógłby rozlać się na całą rodzinę `connection.*` albo
na komendy wołane przed zalogowaniem, a nikt by tego nie zobaczył — odmowa
zniknięta wygląda jak działanie.

## budowa/server/internal/core/zapora_fotografii_test.go

Obróbkę zdjęć i wydania drukarskie robi się zwykle programami zewnętrznymi.
Sięgnięcie po nie jest tu jednym `exec.Command` i wygląda w kodzie
niewinnie, a kosztuje całą funkcję u Operatora: arsenał produktu stoi
wkompilowany w binarium serwera, a u Operatora leży cienka instalka — samo
okno. Funkcja zależna od programu, którego instalka nie niesie, jest
u niego odmową, nie funkcją. Na maszynie deweloperskiej, gdzie te programy
bywają doinstalowane ręcznie, sprawdzian takiej funkcji świeciłby zielono
i nikt by się nie dowiedział.

Zapora pilnuje dwóch rzeczy: że pliki rodzin `design.photo.*`
i `design.print.*` nie wołają procesu, ani przez `exec.Command`, ani przez
pomocnika drzewa `zewnetrzne.Wolaj`; i że żaden plik obszaru Design nie
wymienia nazw silników obrazu i obrysowywania konturów, po które sięga się
przy takiej pracy — także w komentarzu, bo nazwa w komentarzu jest
wskazówką dla następnego wykonawcy, a wskazówka w tę stronę jest wskazówką
ku szkodzie.

Zapory nie wolno osłabić ani przestawić. Gdy nowa funkcja potrzebuje czegoś,
czego rdzeń nie ma, drogą jest biblioteka Go wkompilowana przez `go.mod`
albo zgłoszenie braku — nigdy program w miejsce biblioteki, która istnieje.
Ten sam warunek pilnuje warsztatu dokumentu Studio
(`zapora_warsztatu_pdf_test.go`) i z tego samego powodu.

Wyjątek wolno dopisać w jednym przypadku: gdy dla danej pracy nie istnieje
żadna biblioteka czysto-Go, a program jest składnikiem pakietu serwera i stoi
w wykazie zależności. Dziś taki wyjątek jest jeden — odczyt liter przy
wniesieniu materiału źródłowego (`plikiDesignuZOdczytemPisma`) — i jest
ograniczony do jednego pliku oraz jednego programu. Wyjątek bez tych dwóch
ograniczeń nie jest wyjątkiem, tylko zdjęciem zapory.

Wykaz w `silnikiSpozaInstalki` jest nazwami programów, nie bibliotek:
`pdfcpu`, `tdewolff/canvas`, `disintegration/imaging` i `x/image` są
bibliotekami Go wkompilowanymi w binarium i wolno ich używać wszędzie.
Zakaz dotyczy tego, co trzeba by uruchomić jako proces.

Reguła produktu ma dwie połowy: gdy jest biblioteka Go, bierze się
bibliotekę, bo `exec` jest regresem — tego pilnuje cała ta zapora; gdy nie
ma biblioteki Go, program jest składnikiem pakietu serwera. Rozpoznanie
pisma jest właśnie drugim przypadkiem: czytnika liter w czystym Go nie ma,
a Tesseract stoi w wykazie zależności pakietu serwera
(`zaleznosci_zewnetrzne.go`) i jest wołany tą samą drogą, co w Studiu
i Translate. `design.mockup.import` bez odczytu liter rozpoznaje układ
obszarów i to, które z nich są liniami tekstu, ale treści napisów nie
czyta. Zakaz procesu bez tego wyjątku nie chronił tu instalki — odbierał
funkcję, której nie ma czym zastąpić. Wyjątek obowiązuje jeden plik
i jeden program; własny `exec` zostaje zabroniony i tu, a silniki obrazu
z `silnikiSpozaInstalki` zostają zabronione w całym obszarze Design, ten
plik włącznie — bo dla nich biblioteka Go istnieje.

Pliki w `TestWarsztatFotografiiNieWolaProcesu` rozpoznaje się po nazwie, bo
tak leżą w drzewie: warsztat fotografii w
`adapter_modul_design_fotografia*.go`, część drukarska w
`adapter_modul_design_druk*.go`, a rachunek wektora, ikon i wykresów —
w plikach, które te dwie rodziny wołają.

Zakres `TestObszarDesignNieWymieniaSilnikowObrazuSpozaInstalki` jest
obszarem Design (`adapter_modul_design*.go`), nie całym rdzeniem, i jest to
pomiar granicy odpowiedzialności, nie ustępstwo: warsztat obrazu modułu
Library i narzędzia obrazowe modelu mają własne nagłówki, w których te
nazwy stoją jako zapis decyzji ich autorów. Poszerzenie tej zapory na cały
rdzeń wymagałoby przepisania cudzych plików — a to jest osobna praca
i osobna zgoda. Zapora obszaru Design jest warunkiem, który ten obszar
spełnia w całości i który da się utrzymać.

## budowa/server/internal/core/adapter_narzedzia_obraz_model_twarze.go

Osobny przebieg poprawiania twarzy w `image.upscale` — pole `faces` kontraktu.
Powiększanie samo stoi w `adapter_narzedzia_obraz_model_silniki.go` i to ono
woła ten plik; wspólne zaplecze (pracownia, wołanie binarium, odmowy) —
w `adapter_narzedzia_obraz_model.go`.

Dlaczego pomocnik pythonowy, a nie wydanie ncnn: reszta tej rodziny to binaria
`ncnn` bez Pythona i bez Torcha, więc pomocnik pythonowy jest tu wyłomem
i wymaga powodu. Sieć twarzowa wydana jest jako wagi PyTorcha
(`GFPGANv1.4.pth`), a wydania `ncnn` tej sieci nie publikuje jej autor —
chodzące po sieci przeróbki niosą wagi przeliczone przez osoby trzecie, więc
rdzeń liczyłby nie tym modelem, który leży na dysku, tylko czyjąś kopią
o nieustalonym pochodzeniu. Pomocnik na wagach stojących liczy dokładnie tym,
co Operator ma u siebie, i tą samą drogą, co wektory znaczenia
(`internal/wiedza/pomocnik_osadzen.py`).

Dlaczego przebieg jest drugi, a nie jeden wspólny: kolejność jest zamierzona —
najpierw Real-ESRGAN powiększa CAŁY obraz, potem pomocnik odnajduje twarze
w wyniku i podmienia same wycinki. Dzięki temu wymiary odpowiedzi pochodzą
wyłącznie z powiększenia, a `faces: false` i `faces: true` różnią się
dokładnie tym, co obiecuje opis pola — twarzami, nie rozmiarem.

Czego tu celowo nie ma: gałęzi „gdy pomocnika nie ma, oddaj samo
powiększenie". Obraz bez poprawki twarzy podany jako poprawiony jest tą samą
atrapą, co rozciągnięcie podane jako powiększenie — rdzeń odmawia, nazywając
brak i drogę naprawy.

skryptPomocnikaTwarzy: skrypt jedzie w binarium, a nie leży obok niego,
z tego samego powodu, co pomocnik osadzeń — wdrożenie, w którym ktoś
przeniósł samo binarium, ma działać. Wykładany jest do katalogu jednego
przebiegu (`pracowniaObrazu`), bo znika razem z nim.

granicaOdtwarzaniaTwarzy: bez karty graficznej jedna twarz liczy się na
procesorze kilka sekund, a zdjęcie grupowe niesie ich kilkanaście; do tego
dochodzi start interpretera i wczytanie trzech zestawów wag. Dziesięć minut
znaczy „coś stanęło", a nie „to długo trwa".

katalogWagTwarzyLinux: odbiega od `/usr/local/share/<silnik>`, którym idą
wagi powiększania i wycinania tła, bo te wagi nie są składnikiem pakietu
żadnego programu — są wydaniem modelu pobieranym osobno.

wagiWykrywaniaTwarzy: bez niego nie ma czego odtwarzać: sieć twarzowa pracuje
na wycinku wyrównanym do pięciu punktów charakterystycznych, a te punkty
wskazuje właśnie ten model.

wagiPodzialuTwarzy: z niej powstaje maska wklejenia — bez maski wycinek
wraca do obrazu prostokątem o widocznej krawędzi.

narzedzieOdtwarzaniaTwarzy: wołamy opakowanie `/usr/local/bin/danaco-twarze`,
a nie plik z wnętrza środowiska pythonowego — tak samo jak przy `rembg`.
`zewnetrzne.Wolaj` nie dziedziczy środowiska rdzenia, a biblioteki pomocnika
szukają katalogu pamięci podręcznej i katalogu domowego; opakowanie ustawia
je samo, więc pomocnik jest samowystarczalny niezależnie od tego, kto go woła.

katalogWagTwarzy: zależy od systemu z tego samego powodu, co katalogi wag
powiększania i wycinania tła: na Linuksie drzewo modeli stoi pod `/opt`,
a w wydaniu natywnym Windows jedzie obok rdzenia w `pomocniki/`.

sprawdzWagiTwarzy: sprawdzane są wszystkie naraz, bo przebieg potrzebuje
każdego z nich, a odmowa po dwóch minutach startu interpretera z powodu
trzeciego pliku byłaby czasem straconym. Katalog przychodzi argumentem,
a nie jest brany ze stałej, żeby sprawdzian mógł zmierzyć samą odmowę na
katalogu bez wag.

wylozPomocnikaTwarzy: zapis idzie przez plik tymczasowy i przemianowanie, bo
skrypt obcięty w połowie wystartowałby i wywrócił się komunikatem o składni,
którego nikt nie powiąże z przerwanym zapisem.

poprawTwarze: wejściem jest plik wyniku Real-ESRGAN-a, a nie źródło żądania —
przebieg twarzowy pracuje na tym, co powiększenie już wytworzyło. Wyjście
idzie do osobnego pliku w tej samej pracowni, żeby wynik powiększenia został
nietknięty na wypadek odmowy pomocnika.

opisPrzebieguTwarzy: liczba jest zmierzona przez pomocnika, a nie założona —
zdjęcie bez rozpoznanej twarzy przechodzi przebieg nietknięte i opis ma to
mówić wprost.

## budowa/server/internal/core/adapter_prowenancja_powtorzenie.go

Cztery pozostałe komendy rodziny są odczytem śladu: nie ruszają kanału, nie
kosztują ani grosza i nie zmieniają niczego. Powtórzenie jest czymś innym —
jest nowym wywołaniem kanału modelu, z własnym kosztem i własnym wierszem
w śladzie. Dlatego wchodzi razem z warstwą, która kanały prowadzi.

Wiersz powtórzenia jest zwykłym wierszem prowenancji: wskazuje pierwowzór
jako rodzica, więc drzewo śladu pokazuje, że jedno wzięło się z drugiego.
Powtórzenie ukryte przed śladem byłoby wywołaniem, za które ktoś zapłacił,
a którego rozliczenie nie widzi.

Rdzeń nie powtórzy wywołania, którego treści nie zapisano. Ślad bywa
prowadzony bez treści (`TrescZapisana` fałszywe) albo zredagowany, i wtedy
nie ma czego wysłać po raz drugi. Odpowiedź mówi to wprost polem
`contentAvailable`, zamiast wysyłać pusty prompt i zestawiać jego odpowiedź
z pierwowzorem jak gdyby nigdy nic.

`ZKanalami` wpina rejestr kanałów — jedyną drogę, którą rdzeń wykonuje
wywołanie modelu. Zależność opcjonalna: bez niej cztery komendy odczytu
pracują bez zmian, a powtórzenie odmawia, nazywając brak.

Stan zapisany z ręki odbiłby się od warunku kolumny i zamienił udane
powtórzenie w awarię zapisu — dlatego stan idzie przez odwzorowanie
kontraktu, tak samo jak czyta go `stanKontraktuWywolania`.

`zasiegiPowtorzenia` odtwarza kontekst pierwowzoru, żeby powtórzenie poszło
tam, gdzie poszło oryginalne wywołanie. Zasięg wzięty z powietrza dałby
wywołanie z innymi parametrami wykonania, czyli nieporównywalne.

`modelPowtorzenia` rozstrzyga model zapisywany w śladzie powtórzenia:
wskazany w żądaniu, a w jego braku model pierwowzoru.

## budowa/server/internal/core/adapter_kontekst_zajetosc.go

Tokenizator jest podsystemem rdzenia, nie oszacowaniem: kontrakt tej komendy
stawia sprawę wprost, bez tokenizatora liczba żetonów byłaby wartością
wziętą znikąd. Rdzeń ma więc tokenizator wkompilowany
(`server/internal/tokenizator`) i liczy nim naprawdę — słownikiem BPE, tym
samym podziałem, którym liczy model. Odpowiedź niesie nazwę słownika
(`ContextUsage.tokenizer`), więc czytelnik wie, czym zmierzono. Przybliżenie
„znaki podzielone przez cztery" byłoby tu gorsze niż brak odpowiedzi: myli
się na polszczyźnie o kilkadziesiąt procent, a Operator czytający „w
normie" traci turę na przepełnieniu okna.

Granica okna jest odczytana, nie zmyślona: bierze się z parametru kanału
(`contextWindow` w konfiguracji wiersza rejestru kanałów). Kanał, który jej
nie podaje, daje odpowiedź `available: false` wraz z powodem i wskazaniem
naprawy — bo pasek zajętości wobec granicy wziętej z głowy pokazywałby
„w normie" albo „prawie pełne" zależnie od tego, co rdzeń akurat zgadł.

Prompt systemowy, historia rozmowy i pamięć są liczone osobno, każda swoim
tekstem. Suma jest sumą tych trzech, a nie osobnym pomiarem — inaczej części
nie sumowałyby się do całości i okno pokazywałoby dwie prawdy naraz.

`parametrGranicyOkna` to nazwa parametru konfiguracji kanału niosącego
wielkość okna kontekstu modelu w żetonach. Kontrakt `channel.add` nie ma
osobnego pola na tę liczbę, więc jedzie ona parametrem — tą samą drogą co
`credentialRef`.

Drugi składacz promptu systemowego dałby drugą liczbę żetonów warstwy
systemowej i rozjechał się z pierwszym przy pierwszej zmianie warstw —
dlatego `tozsamosc` oddaje ten sam prompt, który pojedzie do modelu.

## budowa/server/internal/core/adapter_narzedzia_obraz_wektor_slad.go

Odpowiedzialność pliku: zamiana maski rastrowej na ścieżki — obrys konturu,
upraszczanie łamanej i ścieńczanie do linii środkowej. To czysta arytmetyka na
tablicy pikseli; wołający (`adapter_narzedzia_obraz_wektor.go`) rozstrzyga,
skąd maska pochodzi i co z gotowymi ścieżkami zrobić.

Dlaczego własna arytmetyka, a nie program zewnętrzny: zamiana rastra na
ścieżki bywa robiona programem `potrace`. Program ten nie stoi na serwerze,
a instalka Operatora nie niesie żadnego programu — czynność oparta na nim
byłaby u odbiorcy odmową, a nie funkcją. Obrys konturu, upraszczanie
Ramera-Douglasa-Peuckera i ścieńczanie Zhanga-Suena to algorytmy opisane
i skończone; wkompilowane w binarium działają wszędzie tam, gdzie działa
rdzeń. Wynikiem jest łamana, nie krzywa Beziera. To rozstrzygnięcie, nie brak:
łamana po uproszczeniu opisuje kontur wiernie i przewidywalnie, a dopasowanie
krzywych wprowadza odchylenie, którego Operator nie kontroluje żadnym polem
kontraktu. Pole `simplify` steruje właśnie odchyleniem łamanej i mówi wprost,
ile go wolno.

kierunkiObrysu: kolejność ma znaczenie — śledzenie konturu metodą Moore'a
chodzi po sąsiadach właśnie w tym porządku i to on rozstrzyga, że kontur
zewnętrzny wychodzi zgodnie z ruchem wskazówek.

obrysyMaski: metoda jest klasycznym śledzeniem sąsiedztwa Moore'a: znajdź
piksel brzegowy, obejdź obszar dookoła, wróć do punktu wyjścia. Piksele już
objęte konturem są znakowane, żeby ten sam obszar nie dał dwóch identycznych
ścieżek. Obszary mniejsze niż `najmniejszyObszar` są pomijane: pojedyncze
piksele szumu dałyby setki ścieżek o wielkości kropki, przez które wynik jest
cięższy od źródła i nie do otwarcia w edytorze wektorowym.

obejdzObszar: granica liczby kroków chroni przed obrazem, którego kontur
z jakiegoś powodu nie domyka się w punkcie wyjścia: pętla bez granicy
zawiesiłaby żądanie na zawsze, a odmowa po granicy jest odpowiedzią.

scienczMaske: algorytm Zhanga-Suena zasila obrys linii środkowej
(`centerline`): rysunek kreskowy obrysowany po konturze dałby każdą kreskę
jako podwójną pętlę, a obrysowany po linii środkowej — jako jedną kreskę.
Algorytm chodzi naprzemiennie dwoma podprzebiegami, aż przestanie cokolwiek
zdejmować. Granica przebiegów chroni przed układem, który oscyluje.

czyZdejmowalny: sąsiedzi liczeni są w kolejności zegarowej od północy.
Warunki są cztery: liczba sąsiadów mieści się w 2..6 (piksel nie jest ani
końcem, ani wnętrzem), przejść z tła do obszaru jest dokładnie jedno (zdjęcie
nie rozerwie linii), oraz dwie pary sąsiadów zależne od podprzebiegu — to one
na przemian ścinają obszar z dwóch przeciwnych stron, żeby linia wyszła
pośrodku, a nie przy jednej krawędzi.

## budowa/server/internal/core/adapter_narzedzia_media_argumenty.go

Odpowiedzialność pliku: układanie wiersza wywołania `ffmpeg` dla pięciu
czynności wyliczenia `MediaOperationKind` oraz dobór kontenera. Rozdział
z `adapter_narzedzia_media_przetworzenie.go` idzie po odpowiedzialności:
tamten plik prowadzi przebieg komendy (źródło, pomiar, binarium, zasób),
a ten mówi wyłącznie językiem `ffmpeg`. Każda odmowa wychodzi stąd przed
uruchomieniem programu. Nazwany brak parametru („brakuje pól startMs
i endMs") jest dla wołającego czymś zupełnie innym niż diagnostyka binarium,
które dostało wiersz bez sensu i odmówiło po swojemu.

argumentyPrzetworzeniaMediow: `-y` stoi przy każdej czynności, bo plik
wynikowy leży w świeżym katalogu tymczasowym i nadpisać może wyłącznie
samego siebie; bez tego przełącznika `ffmpeg` czeka na odpowiedź człowieka,
którego przy nim nie ma, i kończy się dopiero granicą czasu.

argumentyWycieciaMediow składa wycięcie fragmentu. Brak obu granic jest
odmową nazywającą brak. Jedna granica wystarczy i znaczy dokładnie tyle, ile
mówi: sam `startMs` to „od tego miejsca do końca", sam `endMs` to „od
początku do tego miejsca". To nie jest domyślanie się całości, tylko
odczytanie tego, co wskazano. `-ss` i `-to` stoją po wejściu z zamysłem:
przed wejściem `ffmpeg` przeskakuje do najbliższej klatki kluczowej i
granica przesuwa się o ułamek sekundy, po wejściu jest dokładna. `-c copy`
przepisuje strumienie bez ponownego kodowania — fragment ma być tym samym
materiałem, tylko krótszym.

argumentyDzwiekuMediow składa wyodrębnienie ścieżki dźwiękowej. `-vn`
odrzuca obraz — to jest cała treść tej czynności. Wynik nie ma wymiarów
i mieć ich nie będzie; kontrakt przewiduje to wprost polami opcjonalnymi.
Strumień jest kopiowany wtedy, gdy kontener wyprowadzono z kodeka (brak
`format`): dźwięk trafia do nośnika, który zna ten kodek, więc ponowne
kodowanie pogorszyłoby materiał bez powodu. Przy formacie wskazanym wybór
kodeka zostaje przy `ffmpeg`u — kopia mogłaby do wskazanego kontenera nie
pasować, a odmowa binarium byłaby wtedy karą za spełnienie prośby
wołającego.

argumentyRozmiaruMediow składa zmianę rozdzielczości. Brak obu wymiarów jest
odmową: „zmień rozmiar" bez podania rozmiaru nie niesie żadnego polecenia.
Wymiar niepodany idzie jako `-2`, a nie jako liczba wyliczona samodzielnie:
`-2` znaczy dla filtra „dobierz z proporcji źródła, zaokrąglając do liczby
parzystej" — proporcje zostają nietknięte, a parzystość jest wymogiem
kodeków obrazu, które próbkują chrominancję co dwa piksele i wysokości
nieparzystej wprost odmawiają. Dźwięk jest przepisywany bez kodowania
(`-c:a copy`): zmiana rozdzielczości dotyczy obrazu i nie ma prawa dotknąć
ścieżki dźwiękowej.

argumentyKlatkiMediow składa zrzut pojedynczej klatki. `-ss` stoi przed
wejściem, odwrotnie niż przy wycięciu, i to jest wybór: przeskok do klatki
kluczowej jest tu tani i szybki, a różnica ułamka sekundy nie ma znaczenia
dla zrzutu poglądowego — przy wycięciu miałaby, bo przesuwa granice
fragmentu. Brak `startMs` znaczy początek materiału: zrzut klatki ma sens od
pierwszej klatki, a „pierwsza" jest wskazaniem tak samo jednoznacznym jak
każde inne. `-frames:v 1` ogranicza wynik do jednej klatki, `-an` odrzuca
dźwięk, którego obraz nie uniesie.

kontenerDzwiekuMediow dobiera nośnik do zmierzonego kodeka dźwięku, tak żeby
ścieżkę dało się przepisać bez ponownego kodowania. Kodek spoza wykazu
oddaje pustkę, a wołający robi z niej odmowę proszącą o wskazanie formatu —
bo nośnik dobrany na chybił trafił kończy się odmową samego binarium.

kontenerZrodlaMediow: nazwa bywa wykazem, nie pojedynczym słowem — jeden
zestaw procedur czyta całą rodzinę kontenerów i `ffprobe` oddaje wtedy
wszystkie naraz („mov,mp4,m4a,3gp,3g2,mj2"). Wybór idzie po kolejności
pierwszeństwa tego pliku, a nie po kolejności w wykazie programu: dla
materiału MP4 wykaz zaczyna się od „mov" i wynik nosiłby rozszerzenie,
którego nikt nie zamawiał, choć „mp4" stoi w tym samym wykazie o jedną
pozycję dalej.

## budowa/server/internal/core/zaleznosci_wykaz_wydruk.go

Dziennik startu mówi Operatorowi, czego brakuje na tej maszynie — jest
diagnozą stanu zastanego. Prowizjonowanie potrzebuje czego innego: pełnego
wykazu pakietów do postawienia na serwerze docelowym, niezależnie od tego,
co stoi na maszynie budującej. To ta sama wiedza (te same deklaracje
narzędzi), ale wyprowadzona kompletnie i w postaci, którą skrypt rozbierze
na pola.

`zaleznosci_zewnetrzne.go` w nagłówku ostrzega: druga lista rozjedzie się
z pierwszą przy pierwszej zmianie pakietu. Skrypt prowizjonowania
(`scripts/arsenal-serwera.sh`) nie przepisuje nazw — woła binarium rdzenia
w tym trybie i konsumuje wynik. Zmiana pola `Pakiet` w deklaracji narzędzia
dojeżdża więc i do sondy startowej, i do prowizjonowania jednym ruchem.

Silnik kontenerów jest wstrzymany i nie może zostać postawiony milcząco.
Gdyby ten rozdział robił skrypt dopasowaniem napisów w powłoce, byłby
drugą regułą obok deklaracji — nietypowaną, niesprawdzalną i cichą przy
pomyłce.

Silnik kontenerów rozpoznaje się po programie, nie po podpowiedzi
instalacyjnej: w wykazie stoją dwie jego deklaracje („Docker" warsztatu
Developera i „Docker (klient wiersza poleceń)" modułu Terminal), niosą
różne podpowiedzi, a obie mają trafić do warstwy decyzyjnej. Pozostałe
warstwy bierze się z treści podpowiedzi, bo ona już dziś mówi, czym program
dociągnąć. Kolejność pytań jest istotna: podpowiedź `go install
github.com/...` niesie adres GitHuba, a nie jest krokiem ręcznym.

Wiersze komentarza wykazu zależności zaczynają się od `#` — skrypt
konsumujący je pomija. Kolejność jest ta sama, którą ustala
`zaleznosciZewnetrzne` (alfabetyczna po nazwie czytelnej), więc dwa kolejne
wywołania dają ten sam wykaz.

Wykaz zależności niesie z mowy tylko syntezator zapasowy (eSpeak NG), bo
tylko on jest zwykłym programem na ścieżce. Reszta arsenału mowy to piper
wraz z plikami głosów, biblioteka pythonowa rozpoznawania (faster-whisper)
i wagi jej modelu — rzeczy stawiane inaczej niż pakietem dystrybucji, a bez
nich mikrofon i odsłuch odmawiają Operatorowi tak samo. Wartości pochodzą
ze stałych rdzenia (nazwy programów, miejsca arsenału, zmienne wskazania)
oraz z ustawień silnika mowy (`mowa.ModelDomyslny`) — nie są tu wpisane po
raz drugi.

Plik zależności pomocnika transkrypcji jest wskazywany ścieżką wyliczoną
z położenia skryptu, a nie wypisaną wprost:
`pomocniki/transkrypcja/wymagania.txt` jest jedynym miejscem, w którym stoi
nazwa i wersja biblioteki rozpoznawania, więc prowizjonowanie ma go
zainstalować przez `-r`, zamiast powtarzać nazwę pakietu u siebie.

## budowa/server/internal/core/wstawienia_skutek_test.go

Skutek wstawień: czy obiekt osadzony w dokumencie niesie pochodzenie i czy
odmowa jest nazwana tam, gdzie rdzeń nie ma czym wykonać czynności. Szkody,
które ten plik ma wykluczyć: 1. obraz wstawiony bez zapisanego pochodzenia —
za tydzień nikt nie odtworzy, na czym pismo się opiera; 2. obiekt o zerowym
rozmiarze, czyli niewidzialny, oddany jako wstawiony; 3. uprzejma odmowa bez
nazwania braku — Operator ma wiedzieć, po czyjej stronie brakuje i którą
drogą czynność jest wykonalna; 4. usunięcie obiektu, po którym wykaz nadal go
pokazuje.

## budowa/server/internal/core/adapter_rozmowa_zalaczniki.go

Odwołania z `MessageSendRequest.attachments` idą dwiema drogami: do wiersza
wiadomości (`adapter_rozmowa.go`) oraz do treści zapytania kanału. Drugiej
drogi nie da się ominąć, bo `zapytanieKanalu` niesie samą treść, a kanał CLI
nie ma osobnego pola na załącznik.

Kontrakt nazywa elementy listy „odwołaniami do załączników" i nie zawęża
ich kształtu. Rdzeń rozstrzyga trzy przypadki: URI danych
(`data:image/png;base64,…`) niesie bajty ze sobą — lądują w magazynie
treści rdzenia, a do modelu idzie ścieżka bloba, tak posyła adnotacje
moduł Browser (`client/src/moduly/browser/warstwa-adnotacji.ts`); ścieżka
bezwzględna istniejącego pliku idzie dalej bez kopiowania — tego samego
rodzaju odwołanie zapisuje moduł Library (`adapter_modul_library_magazyn.go`);
pozostałe odwołania wracają nazwane jako niedoręczone, bez zgadywania.

Ładunek base64 nie wchodzi do treści: model nie zdekoduje kilobajtów
napisu, a okno kontekstu za nie zapłaci. Wchodzi wyłącznie ścieżka, po którą
model sięga narzędziem odczytu pliku.

Plik nie zmienia kształtu kontraktu i nie sprząta magazynu: blob adresowany
sumą kontrolną żyje tak samo długo jak wiersz wiadomości, który go
wymienia.

Pole `Sciezka` puste w `zalacznikTury` znaczy odwołanie nierozwiązane —
wtedy `Powod` podaje przyczynę, a treść wiadomości niesie to zdanie do
modelu i do operatora.

Blok odwołań idzie w treści, a nie osobnym polem, bo kanał główny przyjmuje
od rdzenia dokładnie jedno wejście rozmowy: `injection.Zapytanie.Tekst`,
wysyłane na stdin jako JSON-lines (`injection/ustawienia.go`,
`injection/przebieg.go`). Ścieżka wpleciona w treść jest więc jedyną drogą
odwołania do modelu, a droga ta jest skuteczna: model czyta plik narzędziem
odczytu. Odwołania niedoręczone stoją w tym samym bloku, aby model wiedział,
że coś pokazano i czego nie dostał.

TestIkonaNieznanaOdmawiaWskazaniemKatalogu: nazwa szukana jest umyślnie bez
ani jednego słowa z etykiet katalogu — katalog Designu dopasowuje po
zawieraniu w obie strony, więc nazwa niosąca „katalog" albo „folder"
trafiłaby we wzór i sprawdzian mierzyłby coś innego.

## budowa/server/internal/core/usterki_wejscia_test.go

Straże drogi wejścia: trzy sprawdziany trzymają zachowania, których zerwanie
zamyka Operatorowi drogę do platformy albo otwiera ją komuś, kto nie powinien
wejść. Trzy razem, a nie osobno, bo trzymają się nawzajem: bramka zamknięta
do potwierdzenia adresu bez drogi wyjścia z nieudanego nadania zamieniałaby
zatrzymany serwer poczty w trwałą utratę produktu, a wykaz urządzeń bez
wejścia hasłem nie miałby czego pokazać w oknie odbierania dostępu.

TestBramkaZamknietaDoPotwierdzeniaAdresu: rejestracja nie zakładała sesji, ale
zakładała kotwicę, a `auth.login` nie pytał o stan potwierdzenia w ogóle.
Operator wołał więc logowanie zaraz po rejestracji i dostawał pełny token, nie
zaglądając do skrzynki. Adres jest jedyną drogą odzyskania konta — adres
niesprawdzony (literówka, cudza skrzynka, domena bez rekordu) wychodziłby na
jaw dopiero w dniu, w którym trzeba nim odzyskać dostęp, a rejestracji nie da
się powtórzyć.

Odmowa ma prowadzić do naprawy: mówić, czego brakuje i czym to zrobić. „Nie
wolno" bez drogi dalszej zostawia Operatora przed zamkniętą bramką bez
klucza.

TestNieudaneNadanieListuSchodziNaDrogeBezPoczty: sprawdzana bywała wyłącznie
obecność nastaw konta nadawczego, przed zapisem. Samo nadanie idzie ostatnie,
już po zapisaniu konta, kotwicy i drogi, a nastawa wskazana nie znaczy, że
serwer odpowiada: przekaźnik bywa zatrzymany, zapora zamknięta, a nazwa hosta
wpisana z literówką. Cofnięcie rejestracji było ratunkiem przed platformą nie
do otwarcia, dopóki bramkę zamykał brak potwierdzenia adresu — konto
zostawione po nieudanym nadaniu nie miało czym wejść. Odkąd konto bez
potwierdzonego adresu wchodzi hasłem, ratunek jest zbędny, a sam był pułapką:
literówka w nazwie hosta zamykała pierwsze uruchomienie równie szczelnie jak
brak poczty w ogóle. Rejestracja schodzi więc na drogę bez poczty i kończy
się tym samym stanem, co instalka, która nadajnika nie ma wcale. Rejestracja
wykonuje się raz: gdyby nieudane nadanie cofało ją bez otwarcia drogi powrotu
albo zostawiało konto bez klucza, pierwszy Operator tracił platformę na jedną
niedostępność serwera poczty.

Droga potwierdzenia zostaje po nieudanym nadaniu i to nie jest usterka: leży
jako sam skrót materiału, który do nikogo nie dojechał, i wygasa po godzinie.
Kasowanie jej wymagałoby czwartej czynności repozytorium dla stanu, który sam
się kończy.

TestWykazUrzadzenWidziWejscieHaslem: sesja brała urządzenie z wiersza metody,
nie z żądania. Kotwica hasła urządzenia nie ma i mieć nie może — hasło nie
jest materiałem jednej maszyny — więc `deviceId` z `auth.login` był
porzucany. Maszyna nie pojawiała się w `device.list`, a `device.revoke` z jej
identyfikatorem wracał `revoked: false` i zostawiał token czynny, w oknie,
które istnieje po to, żeby Operator dostęp odbierał.

## budowa/server/internal/core/adapter_wywolywacz.go

Skrót globalny przechwytuje powłoka programu okiennego na maszynie
Operatora, nie rdzeń: rdzeń stoi na serwerze i klawiatury tamtej maszyny
nie widzi. Podział jest więc taki: rdzeń trzyma nastawę, dzięki czemu
skrót jest ten sam na każdej maszynie tego samego Operatora i przeżywa
ponowne zainstalowanie okna; powłoka rejestruje skrót u siebie i to ona
wie, czy się udało — skrót zajęty przez inny program zajmie go dalej,
cokolwiek rdzeń o tym sądzi.

`supported` mówi, czy po drugiej stronie stoi powłoka, która w ogóle umie
zarejestrować skrót globalny. Rdzeń wie to z jednego miejsca: z powitania.
Klient deklaruje w nim swoje zdolności (`connection.hello`, pole
`capabilities`), a rdzeń zapamiętuje deklarację. Przeglądarka takiej
zdolności nie zadeklaruje i wtedy odpowiedź mówi wprost, że skrótu nie ma
kto przechwycić — zamiast obiecywać skrót, który nikogo nie obudzi.

`registered` nie jest zgadywane: jest prawdą wtedy i tylko wtedy, gdy skrót
jest niepusty oraz stoi powłoka deklarująca zdolność. Rdzeń nie twierdzi,
że rejestracja się powiodła, gdy nie ma komu jej wykonać.

Zapis nastawy idzie pierwszy, a odpowiedź mówi osobno o zapisie i osobno
o rejestracji — kontrakt mówi to wprost i tak jest tutaj. Skrót zajęty
przez inny program nie jest błędem zapisu — Operator zwolni go później
i nie będzie musiał wpisywać nastawy od nowa.

## budowa/server/internal/core/montaz_zrodla.go

Odpowiedzialność pliku: budowniczowie rejestrów i źródeł, z których korzysta
montaż rdzenia — rejestr kanałów modelu wraz z pulą kont rotacji, źródło
wierszy konfiguracji spod adresu złożonego, rejestr definicji z katalogu
ustawień oraz katalog akcji. Wszystkie są sterowane danymi: nowy kanał, nowa
pozycja okna konfiguracji i nowa akcja to nowy wiersz, nie nowa gałąź
w kodzie. Niepowodzenie pierwszego odczytu nigdy nie przerywa startu — rdzeń
rusza z zawartością uboższą i odbudowuje ją przy kolejnym odczycie. Osobno od
montaz.go, bo montaż mówi, co z czym się wiąże, a ten plik — jak powstaje
każde źródło.

Kanał główny w rejestrKanalow obsługuje oba rodzaje procesu lokalnego: „cli"
(program code CLI) i „lokalny" (proces lokalny rozmawiający strumieniem).
Obie fabryki dzielą jedną pulę kont — rotacja po wyczerpaniu limitu ma sens
tylko przy wspólnej pamięci wyczerpania. Rodzaj „sdk" nie ma tu fabryki: bez
dostawcy SDK kanał tego rodzaju nie ma jak działać, więc rejestr nie pokaże
go jako czynnego, zamiast udawać gotowość.

pulaKont: katalog kont jest źródłem pierwszym, gdy repozytorium istnieje —
niezależnie od tego, czy na starcie ma już wpisy. Pula buduje się z bieżącego
wykazu (może być pusty) i zawsze wpina w nią źródło oraz utrwalacz:
wyczerpanie zapisuje się do bazy, a bieżący wykaz odczytuje się na progu tury
przez OdświeżZeŹródła, więc konto dodane komendą account.* wchodzi do rotacji
bez restartu — także gdy pula wystartowała pusta. Ścieżka zapasowa z dysku
zostaje tylko dla instalacji bez katalogu kont.

## budowa/server/internal/core/skutek_ukladu_orkiestracji_test.go

Skutek dopełnień układu zależności obejmuje: czy bramka, grupa, kompensacja
i spięcie z MultitaskingAI zostawiają po sobie wiersz — i czy spięcie
naprawdę przestawia kolejki, zamiast być znacznikiem, na który nikt nie
patrzy. Każdy sprawdzian schodzi do bazy własnym zapytaniem, bo odpowiedź
komendy oddaje układ po zapisie i wyglądałaby tak samo, gdyby zapis nie
doszedł.

nazwaKontaZKatalogu: kontrakt identyfikuje konto liczbą (`Account.id` = klucz
główny wiersza), a pula rotacji kodem (`konto.nazwa`) — wskazanie z obszaru
account konfiguracji sesji i z kolumny `kanal_modelu.konto_id` przychodzi
więc w innym słowniku niż ten, którym mówi pula. Wskazanie nieliczbowe jedzie
dosłownie: Operator mógł wpisać nazwę wprost.

## budowa/server/internal/core/pola_skutek_test.go

Skutek pól dokumentu: czy pole wchodzi policzone, czy pole obliczane liczy
naprawdę i czy pole bez czego policzyć mówi to wprost. Szkody, które ten plik
ma wykluczyć: 1. pole wstawione z pustą wartością i odpowiedzią „wstawiono" —
Operator widziałby w dokumencie puste miejsce; 2. pole obliczane liczone
w kliencie, a w rdzeniu udawane; 3. właściwość dokumentu, której rdzeń nie
zna, przyjęta w ciszy; 4. wzór daty przyjmowany układem odniesienia
biblioteki, którego Operator nie wpisze.

## budowa/server/internal/core/handlers_queue_zlecenia.go

Ten plik stoi osobno od `handlers_queue.go` (cykl życia kolejki) i
`handlers_queue_wiazania.go` (wykaz i wiązanie), bo osobny jest przedmiot:
tam bytem jest kolejka, tutaj zlecenie. Port jest rozszerzeniem portu
`Kolejki`, nie drugim portem — silnik kolejek pętli sesyjnej, MultitaskingAI
i modułu Automations jest jeden.

Rozgłoszenie idzie `queue.changed` wszędzie tam, gdzie zmienia się
zawartość albo polityka kolejki: Queue Manager rysuje wykaz zleceń
i wskaźnik głębokości z tego, co o kolejce wie, więc dołożone, zdjęte,
odłożone, podzielone, scalone, skierowane, rozgałęzione i uwarunkowane
zlecenie czyni jego obraz nieaktualnym. Trzy odczyty — wykaz zleceń,
zadania martwe i głębokość — nie rozgłaszają niczego.

Rodzajem zmiany jest `updated`, nie `created`: bytem zdarzenia
`queue.changed` jest kolejka, a kolejka przy dołożeniu zlecenia nie
powstaje, tylko się zmienia. `created` opisywałoby założenie samej kolejki
i tak jest używane w `queue.create`.

`zarejestrujCzynnosciZlecen` wpina sześć czynności, których wynikiem jest
zlecenie (albo ich wykaz), a nie kolejka — rozgłoszenie dobiera więc kolejkę
osobno, po identyfikatorze z żądania.

`rozglosKolejkeZlecenia` dobiera kolejkę po identyfikatorze i rozgłasza jej
zmianę. Nieudany dobór kończy wyłącznie rozgłoszenie — komenda już się
powiodła i odmawianie jej z powodu zdarzenia byłoby odwróceniem porządku.

## budowa/server/internal/core/adapter_schowek.go

Odpowiedzialność pliku: rodzina `clipboard.*` — trwała historia schowka
Operatora.

Czyj jest schowek i czego rdzeń nie robi: schowek należy do maszyny
Operatora. Rdzeń stoi na serwerze i schowka tej maszyny nie widzi — ani go
nie czyta, ani do niego nie pisze. Ta rodzina nie udaje inaczej: nie ma tu
ani jednej ścieżki, która sięgałaby po cudzy schowek, i nie ma komendy
„wklej", bo wklejenie jest czynnością okna. Podział ról jest taki: Operator
kopiuje u siebie, okno oddaje skopiowaną treść rdzeniowi komendą
`clipboard.push` — rdzeń dostaje treść, a nie dostęp. Rdzeń daje tej treści
trwałość: historia przestaje ginąć razem z kartą i jest ta sama na każdej
maszynie tego samego Operatora. Operator wybiera wpis z wykazu
(`clipboard.list`), okno wstawia go u siebie — do pola, do schowka
systemowego, gdziekolwiek. Dzięki temu podziałowi komenda robi dokładnie to,
co obiecuje jej nazwa, i ani kroku więcej.

Powtórzenie nie mnoży wpisów: ta sama treść skopiowana drugi raz podnosi
wpis zastany na czoło wykazu. Rozstrzyga to warunek UNIQUE na odcisku
treści, a nie odczyt-i-zapis w adapterze: dwa okna kopiujące naraz
rozjechałyby się na odczycie.

## budowa/server/internal/core/handlers_auth.go

Zdarzenie `auth.changed` rozgłaszają trzy komendy, każda swoim powodem:
założenie metody, jej zdjęcie i zmiana hasła. `auth.login` nie rozgłasza —
wejście nie zmienia ani składu metod, ani hasła, a sekcja Uwierzytelnianie
dostaje wykaz metod wprost w odpowiedzi. `auth.register` też nie: przy
pierwszym uruchomieniu nie ma komu rozgłosić zmiany. `auth.token.refresh`
nie zmienia stanu uwierzytelnienia w ogóle.

Bez zdjęcia metody nie ma zmiany, więc nie ma czego rozgłaszać.

`rozglosZmianeBramki` wysyła `auth.changed`. Zdarzenie idzie bez sesji
komunikatu: bramka nie należy do żadnej karty sesji — to stan platformy,
a nie stan pracy. Nadajnik niepodłączony nie jest błędem.

## budowa/server/internal/core/skutek_warsztatu_pdf_test.go

Wzorzec szkody, którego pilnuje ten plik, ma w tym produkcie precedens
w module Design: komenda meldowała `status: ok` z wykazem zasobów, za
którymi nie było ani jednego bajtu. Dlatego żaden sprawdzian tutaj nie
kończy się na sprawdzeniu, że odpowiedź jest udana: każdy schodzi po
odwołaniu do magazynu, otwiera plik i liczy jego strony z bajtów, a nie
z odpowiedzi komendy.

Sprawdzian nie pomija się przy braku żadnego programu, bo warsztat PDF nie
uruchamia ani jednego procesu: pracuje biblioteką wkompilowaną w rdzeń.
Materiał powstaje tą samą biblioteką i to jest świadome — mierzone jest
działanie warsztatu, nie zgodność dwóch bibliotek między sobą.

## budowa/server/internal/core/adapter_kanaly.go

adapterKanalow: zapis idzie do tabeli rejestru, odczyt do rejestru kanałów
zbudowanego z tej samej tabeli. Po każdej zmianie rejestr jest odświeżany, więc
dopisanie wiersza natychmiast daje działający kanał — bez zmiany w kodzie
i bez restartu. Identyfikatorem kanału w kontrakcie jest kod wiersza, nie
numer wiersza: kod przeżywa przeniesienie bazy i jest tym, co widzi okno
komunikacji. Pole `sejf` obsługuje wyłącznie `channel.credential.status`
(`adapter_kanaly_sprawdzenie.go`) i widzi z sejfu tylko odczyt. Zależność
opcjonalna: bez niej stan poświadczenia mówi „nieustawione" wraz
z odwołaniem, pod którym rdzeń szukał.

kluczOdwolaniaKanalu: nazwa parametru konfiguracji kanału niosącego odwołanie
do danych dostępowych — nazwę zmiennej środowiskowej, pod którą Operator
trzyma klucz kanału API. Sama nazwa nie jest sekretem, więc jedzie
w parametrach; wartość mieszka poza bazą, a kanał API czyta ją przy wysyłce.

odwolaniePoswiadczeniaKanalu: bez tej drogi kolumna poswiadczenie_odwolanie
zostaje pusta i kanał API nie ma skąd wziąć nazwy zmiennej z kluczem. Brak
parametru albo pusta wartość znaczy kanał bez uwierzytelnienia. Kontrakt nie
ma osobnego pola na to odwołanie, więc jedzie ono parametrem konfiguracji.

kluczKontaKanalu: nazwa parametru konfiguracji kanału niosącego powiązanie
z kontem — identyfikator wiersza rejestru kont (kanal_modelu.konto_id).
Kontrakt ChannelAdd/Update nie ma pola accountId, więc powiązanie jedzie
parametrem konfiguracji, tą samą drogą co credentialRef.

kontoKanalu: wartość jest identyfikatorem wiersza rejestru kont, przyjmuje
postać liczby albo napisu liczbowego (JSON koduje liczby jako float64). Brak
parametru, wartość pusta albo nieliczbowa znaczy kanał bez powiązania — dana
pomocnicza nie może wywrócić zapisu.

brakiWierszaKanalu: wiersz bez rodzaju kanału nie przechodził dotąd więzu
schematu i wracał jako usterka wewnętrzna z treścią zapytania SQL — Operator
dostawał nazwę kolumny bazy zamiast nazwy pola, którego nie wypełnił. Brak
samego pola w treści żądania odsiewa brama kontraktu; tutaj rozstrzyga się
wartość pusta.

## budowa/server/internal/core/handlers_panel_sekcje.go

Odpowiedzialność pliku: układ sekcji panelu okna, `panel.sections.get`
i `panel.sections.set`. Nośnikiem układu jest tabela z migracji
`migracja_105_sekcje_paneli.sql`.

UstawUklad: pole `sections` jest wymagane. Brak pola to brak, nie polecenie:
wraca odmowa, żeby żądanie milczące o układzie nie skasowało układu
zapisanego. Wykaz pusty (`"sections": []`) pozostaje czynnością poprawną —
zdejmuje układ własny i przywraca układ domyślny widoku.

Odmowy: każda odmowa mówi trzy rzeczy — co odmówiło (przedrostek), dlaczego
(powód nazwany) i czym czytelnik to zmieni (czynność po średniku).

UstawUklad: pole `order` rozstrzyga porządek, ale kolejne pozycje nadaje
rdzeń jako 1..N — układ z dziurami („1, 7, 9") wracałby w kolejności
zależnej od bazy.

## budowa/server/internal/core/handlers_role_wykaz.go

Odpowiedzialność pliku: wykaz i zdjęcie nadań ról (`role.list`,
`role.remove`) wraz z rozgłoszeniem zdarzenia `role.changed`. Rodzina
`role.*` ma jedno repozytorium i jeden port: `role.list` i `role.remove`
stoją na tym samym adapterze `adapterRolOkien`, tym samym rejestrze okien
i tych samych dwóch kolumnach `okno_komunikacji`, co `role.assign`
i `role.update`. Port `RoleWykaz` osadza `RoleOkien` tak samo, jak tamten
osadza `Okna`. Układ sekcji panelu okna prowadzi osobny plik
`handlers_panel_sekcje.go`.

WykazRol: wykaz jest widokiem na okna, nie na drugą tabelę: nadanie roli to
para pól okna (`windowRole`, `coordinatorWindowId`), więc wykaz składa się
z tego samego wykazu okien, co `window.list`, wraz z więzią doczytaną
z bazy. Własne zapytanie po rolach dałoby drugą odpowiedź na pytanie o rolę
okna. Rolę ma każde okno, także samodzielne — stąd komplet.

Wcielenie leży tam, gdzie zapisuje je `role.update`. Odczyt nieudany nie
kończy wykazu: wykaz bez wcieleń bije odmowę.

ZdejmijRole: zdjęcie roli to powrót do roli samodzielnej, a nie skasowanie
pola — katalog kontraktu nie ma wartości „brak", a okno poza pętlą jest
stanem wyjściowym pakietu sesji (`session.RolaDomyslna`). Więź znika razem
z rolą, w tej samej czynności: koordynatora niesie wyłącznie okno wykonawcy.
Okno bez roli nie jest odmową — wraca `removed=false`, bo stan docelowy już
obowiązuje.

Wcielenie odchodzi razem z rolą: `persona` opisuje rolę, a nie okno, więc
zostawione przy oknie samodzielnym wychodziłoby w `role.list` i w wykazie
nadań (`client/src/moduly/multitasking/wykaz-nadan-rol.ts`) jako wcielenie
roli, której już nie ma. Niepowodzenie zapisu kończy komendę: pół zdjęcia
roli zostawiłoby ten sam rozjazd, tyle że po odmowie.

rozwiazWiezWykonawcow: więź ma dwa końce. Zdjęcie roli samemu koordynatorowi
zostawiłoby wykonawców z więzią do okna, które koordynatorem już nie jest.
Zamiast odmawiać zdjęcia roli, drugi koniec doprowadza się do stanu
zgodnego: wykonawca zostaje wykonawcą, ale bez koordynatora — stan
dopuszczalny. Niepowodzenie nie cofa zdjęcia roli: rola została zdjęta,
a więź jest następstwem, nie warunkiem.

## budowa/server/internal/core/adapter_narzedzia_obraz_kompresja.go

Dogniatanie zapisu obrazu po rachunku wkompilowanym: `image.convert` oddaje
plik mniejszy o tyle, ile potrafią zdjąć optipng, jpegoptim, pngquant
i cwebp.

Dlaczego to jest osobny krok, a nie inny koder: kodery Go zapisują obraz
poprawnie, ale nie szukają najlepszego zapisu — `image/png` bierze jeden
filtr na wiersz i jeden przebieg deflate, `image/jpeg` zapisuje domyślne
tablice Huffmana, `nativewebp` nie stroi predyktorów. Wymienione programy
robią dokładnie jedną rzecz — przeliczają ten sam obraz na krótszy strumień
bajtów — i robią to lepiej, bo na to je napisano.

Ulepszenie, nie warunek: rodzina `image.*` liczy się biblioteką wkompilowaną
i ma się liczyć dalej na maszynie, na której żaden z tych czterech programów
nie stoi. Dlatego ten plik nie odmawia nigdy: brak programu, niezerowy kod
wyjścia, wynik pusty albo wynik większy od źródła — każdy z tych przypadków
oddaje bajty wejściowe bez zmiany. Obraz ma się zapisać także wtedy, gdy nie
ma czym go dogniatać, a odmowa w tym miejscu zamieniłaby ulepszenie w wymóg
wobec wdrożenia. Z tego samego powodu nie ma tu odmowy nazywającej brak,
jaką niesie `zewnetrzne.BrakNarzedzia`: brak programu dogniatającego nie
jest czymś, o czym Operator ma się dowiedzieć w chwili zapisu obrazu —
dowiaduje się przy starcie, z sondy wykazu zależności.

Bezstratnie znaczy bezstratnie: trzy z czterech programów przeliczają zapis
bez ruszania pikseli — optipng szuka filtrów i dłuższego deflate, jpegoptim
przelicza tablice Huffmana, cwebp w trybie `-lossless` stroi predyktory
WebP. `pngquant` jest inny — sprowadza obraz do palety, więc piksele
zmienia. Wchodzi wyłącznie wtedy, gdy żądanie wprost prosi o zapis stratny
(`lossless: false`); przy braku pola i przy `lossless: true` nie jest
w ogóle wołany. Pomylenie tych dwóch rzeczy oddałoby model prosząc o zapis
bezstratny obraz o zmienionych barwach. Z tej samej strony patrzy druga
reguła: dogniatanie nie dokłada pokolenia kompresji stratnej. WEBP stratny
powstaje już programem, bo kodera stratnego w Go nie ma, więc przy
`lossless: false` ten plik zostawia go nietkniętym — drugi przebieg kodera
stratnego odjąłby jakość, nie bajty.

narzedziePngquant: zapis jest stratny, sprowadza obraz do palety zamiast
przeliczać zapis bez ruszania pikseli, jak pozostałe trzy programy.

dogniecZapisObrazu: nie zwraca błędu żadną drogą — to jest istota tego kroku.
Wołający nie ma tu czego obsłużyć: zapis obrazu już się udał, a ten krok może
go wyłącznie skrócić.

Paleta idzie po optipng: pngquant oddaje plik palety, który optipng jeszcze
skraca, a odwrotna kolejność marnuje pierwszy przebieg.

`--dest` żąda katalogu, nie pliku, i zachowuje nazwę źródła — dlatego nazwa
wejścia i wyjścia jest ta sama, a różni je katalog.

WEBP stratny powstaje już programem (`policzProgramem`), bo kodera stratnego
w Go nie ma. Ponowne przepuszczenie gotowego pliku przez koder stratny byłoby
drugą stratą na tych samych pikselach — dogniecenie ma skracać zapis, a nie
dokładać pokolenie kompresji.

Formaty bez programu dogniatającego (avif, tiff, gif) wychodzą takie, jakie
przyszły. Milczenie jest tu właściwe: nie ma czego zgłaszać.

przezPlik: pliki pośrednie są konieczne z tego samego powodu, co przy
silnikach neuronowych (`pracowniaObrazu`) — optipng, pngquant i cwebp żądają
ścieżki wyniku i nie umieją pisać na standardowe wyjście. jpegoptim by umiał,
ale idzie tą samą drogą — jedna droga zamiast dwóch jest tu warta jednego
zapisu na dysk.

Wejście i wyjście mają tę samą nazwę w dwóch katalogach, bo jpegoptim
wskazuje wynik katalogiem, a nie nazwą pliku; pozostałym trzem programom
jest to obojętne.

Kod niezerowy znaczy tu najczęściej „nie umiem tego skrócić" (pngquant
kończy tak zapis, którego nie da się sprowadzić do palety w zadanej
jakości). Zapis pierwotny jest wtedy właściwą odpowiedzią.

Program, który zakończył się powodzeniem i nie zostawił pliku krótszego, nie
miał czego skrócić. Oddanie jego wyniku mimo to powiększyłoby zasób w imię
jego zmniejszenia.

## budowa/server/internal/core/adapter_doradcy_konsultacja.go
Adapter doradcy wykonuje konsultację tą samą drogą, którą przechodzi każdy inny wołacz kanału; ten plik zawiera wyłącznie to, co należy do komendy: odczyt okna, złożenie pytania z danych okna, strumień jawności i przełożenie wyniku na kontrakt.

Jawność nie działa jak bufor: ujście konsultacji jest nadawcą strumienia, więc fragmenty trafiają kopertami stream.chunk do okna pytającego, tą samą drogą, którą płynie odpowiedź agenta, debata Roundtable i wyjście terminala. Bez tego stanowisko widziałoby wyłącznie skrót rady w zdarzeniu advisor.consulted, a pytanie, doradca i rada nie stałyby w jednym oknie.

Kim jest pytający, rozstrzyga okno, a nie żądanie: model podaje wyłącznie sprawę, kontekst i prośbę o doradcę, więc pole pytającego nie da się podać samym żądaniem.

Okno zamknięte nie konsultuje, ponieważ w oknie zamkniętym nie pracuje agent, a zdarzenie konsultacji rozgłaszałoby pracę okna, które stanęło. Sprawdzenie jest tym samym, które wykonują pozostałe moduły oknowe. Kod odpowiedzi mówi prawdę o powodzie: stan zasobu, nie brak kanału, i nie jest ponawialny, bo okno samo się nie otworzy.

Identyfikator strumienia jest tożsamością jednej konsultacji: wchodzi do koperty jako identyfikator żądania i do fragmentów jako identyfikator wiadomości, więc klient wie, że fragmenty należą do jednego wywołania.

Wartości doboru w kontrakcie są dwie, ale wskazanie Operatora nie ma jeszcze w produkcie danych źródłowych: pole kanału modelu okna mówi, którym modelem pracuje okno, nie kogo Operator wyznaczył na doradcę. Podstawienie jednego pod drugie ogłaszałoby wskazanie Operatora przy każdej konsultacji bez prośby.

Pole niosące identyfikator strumienia zasila identyfikator wiadomości wszystkich fragmentów kanału; bez niego fragmenty rady jechałyby do okna bez wskazania, do czego należą.

Ujście konsultacji zbiera cały tekst, łącznie z blokiem jawności doklejanym przy konsultacji, ponieważ fragment domykający ma być tym, co Operator widzi w oknie po konsultacji, czyli radą wraz z podpisem, kto jej udzielił.

## budowa/server/internal/core/handlers_konfiguracja_osi.go
Oś jest prostopadła do poziomu zasięgu: poziom określa, jak wąsko obowiązuje
wartość, od zasięgu globalnego do okna komunikacji, a oś określa, dla czego
wartość obowiązuje — dla platformy, dla modelu albo dla konta. Osie niosą
konfigurację odrębną dla każdego modelu i każdego konta.

Adapter osi nie powtarza adaptera podstawowego ustawień, tylko go owija: droga
osi platformy schodzi do adaptera podstawowego bez zmiany, a osobną drogę ma
wyłącznie oś modelu i oś konta. Przekład wiersza repozytorium na wpis kontraktu
należy w całości do pakietu konfiguracji — rdzeń nie prowadzi tu własnego
kodowania wartości.

## budowa/server/internal/core/adapter_kanal_cli.go
Rejestr modeli zna wyłącznie interfejs kanału i klucz adaptera będący wartością danych; pakiet injection zna wyłącznie swój proces, pulę kont i strumień. Żaden z nich nie importuje drugiego, łączy je ten adapter, mieszkający w warstwie składania.

Pula kont żyje dłużej niż jedno wywołanie, więc jest wspólna dla wszystkich wierszy tego rodzaju: rotacja konta po wyczerpaniu limitu ma sens tylko wtedy, gdy pamięć wyczerpania jest jedna.

Fragmenty tury idą w kolejności nadania, pierwszy jest fragment prowenancji, który składa kanał. Błąd tury wraca wynikiem, a nie drugim fragmentem: fragment błędu nadał już kanał, a rejestr modeli nie ma powielać tej samej przyczyny.

Pusta pula i pula bez wpiętego źródła zostają nietknięte przy odświeżeniu z katalogu na progu tury.

Wskazanie kontraktowe konta, identyfikator liczbowy, tłumaczy na kod konta puli resolver montażu; bez tego tor kanału głównego jechałby zawsze rotacją puli.

## budowa/server/internal/core/handlers_konfiguracja_prowenancja.go
Kontrakt niesie w rodzinie config.* dziewięć komend; pozostałe siedem wpinają
odrębne pliki obsługi, w tym cztery komendy jednolitego modelu konfiguracji
sesji. Port jest rozszerzeniem portu ustawień, nie drugim portem: konfiguracja
ma w rdzeniu jedną bramę, a prowenancja osadza port ustawień zamiast wchodzić
osobnym polem, tak jak port pamięci osadza port przestrzeni roboczej. Żadna
z dwóch komend tego pliku nie zapisuje wartości i żadna nie rozgłasza zdarzenia
zmiany konfiguracji — pierwsza liczy prowenancję na świeżo, druga rozstrzyga
wyłącznie zakres.
## budowa/server/internal/core/nieznana.go
Ścieżka obsługi nieznanej komendy jest fail-open: nieznana nazwa nie stanowi błędu protokołu i nie zamyka niczego. Połączenie pozostaje otwarte, sesja pracuje dalej, kolejne żądania są przyjmowane, a klient otrzymuje zdarzenie obszaru żądanego typu z nazwą, której nie rozpoznano, więc widzi przyczynę zamiast ciszy. Nazwę zdarzenia wyznacza kontrakt na podstawie obszaru żądanego typu — rdzeń nie dobiera jej sam i nie zawiera żadnego literału nazwy zdarzenia.

Komunikat niepoprawny strukturalnie nie ma typu ani identyfikatora, więc nie da się zbudować dla niego żądania ani skorelować odpowiedzi z żądaniem. Mimo to odpowiedź zostaje wysłana — zerwanie połączenia byłoby nieproporcjonalną karą za jeden zepsuty bajt. Typ odpowiedzi wyznacza kontrakt: nierozpoznany komunikat bez obszaru należy do obszaru połączenia, więc wraca jego zdarzeniem `*.unknown`, a kod błędu informuje klienta, że przyczyną jest niepoprawna treść, a nie nieznana nazwa.

## budowa/server/internal/core/handlers_konta.go
Konto jest profilem uwierzytelnienia, kanał jest definicją rozmowy z modelem.
Jeden kanał wskazuje konto preferowane, jedno konto obsługuje wiele kanałów,
a pula rotacji bierze konta tego samego rodzaju; skasowanie konta odłącza
kanały, ale ich nie kasuje. Poświadczenie wchodzi żądaniem i nie wychodzi
żadną drogą — ani wykazem, ani odpowiedzią na zapis, ani zdarzeniem zmiany;
odpowiedź niesie wyłącznie znacznik obecności poświadczenia.

## budowa/server/internal/core/adapter_kanaly_sprawdzenie.go
Sprawdzenie jest narzędziem pomocniczym, nie bramką: wynik niczego nie warunkuje, nie wstrzymuje zapisu eksperta, nie blokuje wysłania tury i nie wyłącza kanału. Operator pyta, czy kanał odpowiada, dostaje odpowiedź i sam decyduje, co z nią zrobić; kanał, który nie odpowiedział minutę temu, bywa sprawny teraz i odwrotnie.

Sprawdzenie jest prawdziwym wywołaniem: rdzeń wysyła krótkie zapytanie tym samym rejestrem kanałów, którym jedzie okno rozmowy, i mierzy czas odpowiedzi. Sprawdzenie, które oglądałoby wyłącznie wiersz rejestru, mówiłoby jedynie, że kanał jest skonfigurowany, a Operator czyta z niego, że działa.

Stan poświadczenia jest zawsze stanem, nigdy treścią: oddaje, czy poświadczenie jest ustawione, jakiego jest rodzaju i kiedy zostało zmienione. W pliku nie ma ścieżki, którą sekret wychodzi do klienta; wartość z sejfu służy wyłącznie do rozstrzygnięcia, czy jest niepusta.

Interfejs sejfu poświadczeń zamiast pełnego typu istnieje po to, by adapter widział wyłącznie odczyt: węższy widok jest tu granicą, a nie ozdobą, ponieważ adapter, który nie umie zapisać, nie nadpisze cudzego klucza przy pomyłce.

Kod odpowiedzi kanału bez wskazania jest błędem żądania Operatora, nie błędem rdzenia.

## budowa/server/internal/core/handlers_konta_adapter.go
Repozytorium oddaje wyłącznie znacznik obecności poświadczenia; metoda
wynosząca odwołanie do sekretu nie jest w adapterze wywoływana ani razu,
a struktura kontraktu konta nie ma pola na sekret, więc nie da się go wynieść
nawet przez pomyłkę. Katalog danych obowiązujący dla sejfu poświadczeń wchodzi
montażem: budowa rdzenia stawia jeden sejf nad katalogiem danych, który
właściciel instalacji może przestawić przełącznikiem uruchomieniowym albo
zmienną środowiskową, i wpina go metodą ZSejfem tą samą drogą, którą dostaje
magazyn treści biblioteki. Sejf domyślny zostaje ustanowiony w konstruktorze
dla wywołania bez montażu, żeby konstruktor nigdy nie oddał adaptera bez sejfu.
## budowa/server/internal/core/nosniki_druku_wspolne.go
Wykaz nośników druku jest jeden i stoi poza modułami, ponieważ czytają go dziś dwa moduły: Design wykazem `design.print.paper.list` i wydaniem do druku, Studio nastawami strony (`studio.page.paper.list`, `studio.page.setup.set`) oraz podglądem wydania. Osobne wykazy w każdym module się rozjeżdżały: koperta DL miała w jednym miejscu 99 na 210 mm, a w drugim 220 na 110 mm, a szereg B urywał się w jednym wykazie na B3, w drugim sięgał B5. Plik nie należy więc do żadnego modułu i nie woła ich funkcji: jest wiedzą rdzenia o materiale, tak samo jak przelicznik milimetrów na cale.

Wymiary są zapisane w układzie pionowym, bo nośnik ma jeden wymiar własny, a obrót jest rozstrzygnięciem materiału, nie nośnika. Koperta zapisana szerokością większą od wysokości niosłaby obrót w samej definicji i nie dałoby się określić, czy dany format stoi, czy leży. Wszystkie pozycje są więc pionowe (szerokość nie większa od wysokości), a orientację nakłada nastawa strony.

Koperty serii C mają wymiary normy, a nie zaokrąglone, ponieważ norma ISO 269 wiąże je z szeregiem A: C4 bierze arkusz A4 bez zginania, C5 arkusz zgięty raz, C6 zgięty dwa razy — stąd wymiary 229 na 324, 162 na 229 oraz 114 na 162 milimetrów.

Pozycje B4 i B5 stały wcześniej wyłącznie w wykazie okna Studia, poza wykazem rdzenia, co prowadziło do odmowy rozpoznania nośnika wybranego z widocznego wykazu; obecnie należą do tego samego szeregu ISO B co B1–B3.

Koperta DL ma wymiar normy ISO 269 (110 na 220 mm) zapisany pionowo; wcześniej dwa wykazy niosły dla niej różne liczby — 99 na 210 oraz 220 na 110 — dla tego samego nośnika.

## budowa/server/internal/core/handlers_konta_poswiadczenia.go
Poświadczenie wchodzi żądaniem przy koncie i przy punkcie dostępu i nie
wychodzi nigdy: kolumny na treść sekretu nie ma w schemacie bazy w ogóle,
a odpowiedź kontraktu niesie co najwyżej znacznik obecności poświadczenia
albo odwołanie do niego. Zamiana sekretu na odwołanie należy do portu sejfu
poświadczeń wypełnianego przy montażu. Gdy sejfu nie wpięto, poświadczenie
jest odrzucane w tym samym wywołaniu — nie trafia do bazy, do dziennika ani
do odpowiedzi, a byt powstaje bez odwołania i pracuje dalej; cicha utrata
sekretu byłaby gorsza od jawnego braku sejfu, więc warstwa wyżej widzi to po
pustym odwołaniu i po fałszywym znaczniku obecności poświadczenia.

## budowa/server/internal/core/adapter_narzedzia_archiwum_sciezki.go
Blob magazynu nazywa się swoją sumą kontrolną i nie ma rozszerzenia. Gdyby format brać z nazwy, rozpakowanie własnego wytworu pakowania musiałoby albo odmówić, albo zaufać polu format wiersza, czyli etykiecie wpisanej kiedyś ręcznie, a nie zawartości pliku. Narzędzie zewnętrzne rozpoznaje zip i 7z po sygnaturze bez względu na nazwę, więc etykieta jest zbędna: format czytany jest z sygnatury, bo to wiedza o pliku, nie o jego opisie.

Sygnatury są krótkie i jednoznaczne: każdy z trzech formatów otwiera własna sekwencja bajtów na początku pliku. Plik krótszy niż sygnatura albo o sygnaturze nieznanej jest odmową, a nie domysłem, że to pewnie zip: rozpakowywanie czegoś nierozpoznanego kończyłoby się komunikatem narzędzia zewnętrznego w obcym języku zamiast zdaniem o tym, co jest nie tak.

Format domyślny to zip, ponieważ otwiera się dwukrotnym kliknięciem w każdym systemie, więc Operator dostający archiwum nie potrzebuje niczego doinstalowywać.

Sprawdzenia ścieżki są dwa, bo są dwa sposoby ucieczki z katalogu roboczego okna. Ścieżka bezwzględna omija katalog wprost. Ścieżka względna wychodzi z niego członem wskazującym katalog nadrzędny, i tego nie widać po samym napisie, dlatego liczona jest ścieżka oczyszczona i sprawdzane, czy nadal leży wewnątrz katalogu.

Brak binarium zewnętrznego dostaje inny kod niż niepowodzenie programu: brak narzędzia Operator usuwa jedną instalacją, a wywrócenie się programu jest usterką przetwarzania.

## budowa/server/internal/core/handlers_krok_zlecenia.go
Kontrakt nie niesie dziś żadnej komendy dotyczącej pozycji kolejki ani jej
struktury, więc rejestracja pyta kontrakt o każdą nazwę komendy rodziny
queue.step.* i wpina uchwyt wyłącznie wtedy, gdy nazwa do kontraktu należy.
Dopóki kontrakt tych nazw nie niesie, klient wołający komendę wstrzymania
kroku dostaje odpowiedź komendy nieznanej — rdzeń nie ogłasza zdolności, które
kontrakt nie opisuje; po wniesieniu komend do kontraktu uchwyty wpinają się
bez zmiany tego pliku, bo warunek rejestracji zadaje to samo pytanie, co
sprawdzian zgodności rejestru z kontraktem. Literały nazw komend stoją tu
wyjątkowo jako klucz wyszukania w kontrakcie, a nie jako deklaracja komendy;
po wniesieniu komend do kontraktu mają ustąpić stałym kontraktu, a struktury
pliku mają stać się aliasami typów kontraktu.

Krok widziany z zewnątrz niesie stan pracy i stan sterowania jako dwa odrębne
pola, bo są to dwa różne fakty — gdzie krok stoi w pracy i czego oczekuje od
człowieka; sklejenie ich w jedno pole zmusiłoby do wyboru, który z faktów
zataić. Odpowiedź komendy decyzji o kroku niesie osobno fakt doręczenia
decyzji wykonawcy kroku, bo bez tego pola brak efektu decyzji trzeba by
odgadywać. Port sterowania krokiem bez wpięcia w kolejkach nie jest odmową:
dopóki kontrakt nazw komend nie niesie, odmowa dotyczyłaby komendy, której
i tak nikt nie może zawołać.

## budowa/server/internal/core/handlers_macierz.go
Kontrakt nie definiuje ani jednej komendy macierzy, więc nie ma czego wpiąć;
rejestracja nazwy spoza kontraktu byłaby ogłoszeniem zdolności, której
kontrakt nie opisuje. Macierz dociera do klienta wyłącznie jako pola bytów
nawigacji obsługiwane przez plik nawigacji. Port istnieje, bo macierz ma
w rdzeniu czytelnika — nawigację — a ta bierze odwzorowanie stąd zamiast
sięgać po repozytorium wprost; odwzorowanie modułu na środowiska, w których
jest widoczny, ma dzięki temu jedno miejsce, a gdy komendy macierzy powstaną
w kontrakcie, będzie już co zarejestrować. Moduł nieobecny w wyniku metody
odwzorowania nie ma okna modułowego w żadnym środowisku.

## budowa/server/internal/core/adapter_narzedzia_media_pomiar.go
Pomiar jest komendą osobną od przetworzenia, bo wycięcie fragmentu bez znajomości czasu trwania daje pusty plik, a zmiana rozdzielczości bez znajomości proporcji daje rozciągnięty obraz.

Odpowiedź programu pomiaru jest pytana w zapisie strukturalnym, bo kształt tej odpowiedzi jest zobowiązaniem programu, a wydruk domyślny bywa zmieniany między wydaniami; ucisza się też banner i ostrzeżenia, żeby na wyjściu stał sam zapis strukturalny.

Materiał bez zapisanego czasu trwania, strumień żywy albo kontener bez nagłówka czasu, daje pole czasu trwania zerowe: zero rozpoznawalnie znaczy brak wartości.

Reszty odpowiedzi programu pomiaru plik nie odczytuje: pola, którego nikt nie czyta, nie trzeba utrzymywać przy zmianie wydania programu. Typowany rozbiór strumieni idzie osobno, na tych samych bajtach co zapis surowy.

Czas trwania i rozmiar przychodzą z programu pomiaru jako tekst, sekundy z ułamkiem i bajty, i tak są brane, zamiast wymuszać typ, którego program nie obiecuje.

Odpowiedź nieczytelna programu pomiaru daje odmowę, nie pusty wynik: zakończenie powodzeniem bez poprawnego zapisu strukturalnego znaczy plik nie będący materiałem albo program w wydaniu, którego rdzeń nie rozumie.

Brak wymiarów strumienia oddaje dwie pustki, bo zasób bez wymiarów ma nie nieść wartości liczbowej; zero znaczyłoby zmierzone zero pikseli.

## budowa/server/internal/core/handlers_message.go
Odpowiedź na komendę wysłania wiadomości potwierdza wyłącznie przyjęcie
wiadomości; treść modelu wraca osobno strumieniem fragmentów przez nadajnik.
Rozdzielenie jest celowe, bo strumień trwa dłużej niż wykonanie komendy,
a klient musi dostać potwierdzenie od razu. Komenda zatrzymania odpowiada
zawsze, niezależnie od stanu pętli — przycisk zatrzymania jest czynny bez
względu na to, czy było co zatrzymywać; gdy nie było, wynik niesie znacznik
zatrzymania równy fałsz, a nie błąd.
## budowa/server/internal/core/obsluga.go
Powtarzalna praca dyspozycji dzieje się w jednym miejscu: odczytanie ładunku w kształcie żądania kontraktu, wywołanie czynności domeny, zamiana wyniku albo błędu na odpowiedź protokołu. Dzięki temu pliki obsługi komend zawierają wyłącznie wiązanie nazwy kontraktu z czynnością, bez powielonej obsługi błędów. Typ żądania i typ wyniku pochodzą z pakietu kontraktu, więc zmiana kontraktu przerywa kompilację obsługiwacza zamiast rozjeżdżać się z nim po cichu.

Zgodność z kontraktem sprawdza się po odczytaniu ładunku i przed czynnością domeny, ponieważ odczyt orzeka o kształcie treści, a sprawdzenie zgodności o jej zawartości. Bez tego sprawdzenia czynność domeny dostawała żądanie niepełne i uzupełniała brak wartością domyślną, meldując powodzenie. Kontekst niesie tę operację, bo przepuszczenie żądania mimo braków zostawia wpis w dzienniku rdzenia, a dziennik jedzie właśnie kontekstem.

Tożsamość żądania jedzie kontekstem, bo ładunek jej nie niesie: pole identyfikatora mieszka w kopercie, a czynność domeny dostaje wyłącznie rozpakowaną treść. Bez tego wpisu nadawca strumienia nie miałby czym powtórzyć identyfikatora zadania w kopertach odpowiedzi strumieniowej, choć kontrakt każe mu go powtarzać.

## budowa/server/internal/core/odtworzenie_stanu.go
Rejestr nadzorcy jest pamięcią jednego uruchomienia rdzenia. Bez odtworzenia stanu z bazy rozmowa sprzed restartu wraca wyłącznie historią wiadomości, a sam byt sesji i okna znika — klient trzyma identyfikator, pod którym nie ma już czego wskazać. Odtworzenie wnosi sesje i okna z bazy pod tymi samymi identyfikatorami zewnętrznymi, więc komendy powiązania sesji, odczytu stanu okna i wysyłki wiadomości trafiają w te same byty. Odtworzenie nie startuje procesów okien: proces ginie wraz z rdzeniem, a okno wraca jako byt bez procesu — rejestr procesów zgłosi wtedy stan oczekujący, a pierwsza tura uruchomi proces zwykłą drogą.

Odtworzenie sesji po identyfikatorze rdzenia obsługuje też powrót sesji z kosza: usunięcie zdejmuje sesję z rejestru żywego, a wykaz startowy sesji ładowany z bazy przy rozruchu jej wtedy nie widzi, więc przywrócenie musi wnieść ją do rejestru samodzielnie, tą samą drogą, którą wnosi start rdzenia.

## budowa/server/internal/core/okna_utrwalone.go
Rejestr nadzorcy zna wyłącznie okna bieżącego uruchomienia rdzenia. Po jego restarcie okna wskazane przez klienta istnieją już tylko wierszami w bazie, i wtedy ten odczyt odpowiada na pytanie o okna otwarte przed restartem. Wiersz niesie klucze obce, kontrakt niesie kody, więc przekład dokłada słowniki modułów i kanałów modelu.
## budowa/server/internal/core/pamiec_sesji_kolejek.go
Powiązanie kolejki z sesją i oknami mieszka w pamięci procesu, a nie w bazie ani w kontrakcie, z powodu tymczasowego: kolumna łącząca kolejkę z sesją wskazuje wiersz sesji, a sesja żyje w pamięci pakietu sesji pod identyfikatorem tekstowym, dla którego wiersza nie ma. Kontrakt jest w tym względzie w porządku i nie wymaga zmiany. Wprowadzenie trwałości sesji zdejmie ten plik w całości, a powiązanie wróci wtedy do kolumny bazy.

## budowa/server/internal/core/pokrycie_kontraktu_test.go
Pokrycie kontraktu liczono wcześniej czytaniem źródeł wyrażeniem dopasowującym wywołania rejestracji komend, a ta miara myliła się w obie strony: komendy wpinane przez parametr, a nie literałem nazwy, wyrażenie omijało, więc wychodziły z pomiaru jako komendy bez obsługi, choć rdzeń od dawna na nie odpowiada; a wywołanie rejestracji w gałęzi kodu, do której montaż nigdy nie dochodzi, pomiar liczył jako pokrycie. Rejestr zna prawdę, bo to on rozstrzyga, czy komenda dostanie uchwyt, czy odpowiedź o nieznanej komendzie, więc sprawdzian pyta jego, a nie źródeł.

## budowa/server/internal/core/handlers_mobile.go
Warstwa mobilna pyta o stan warstwy wspólnej — sesje, okna, procesy, kolejki —
a tę obsługuje port nawigacji; rejestr telemetrii postępu ma w rdzeniu jednego
właściciela. Wzorem jest rozszerzenie monitora, które tą samą drogą rozszerza
ten sam port, oraz rozszerzenie pamięci modułu przestrzeni roboczej. Brak
portu nie jest ciszą, tak samo jak w monitorze: obsługiwacze rejestrują się
zawsze, a port, który nie niesie warstwy mobilnej, odpowiada błędem
wewnętrznym z nazwą brakującego bytu, żeby właściciel instalacji zobaczył, że
warstwa mobilna nie jest wpięta, a nie pusty wykaz procesów i wyzerowany stan
platformy. Rodzina komend mobilnych nie rozgłasza zdarzeń, mimo że kontrakt
zna zdarzenie zmiany procesu mobilnego, bo proces mobilny jest procesem
telemetrii postępu, a telemetria ma jednego producenta zdarzeń. Rozgłaszanie
z tego pliku dałoby zdarzenie wyłącznie po zmianie zleconej z telefonu,
a milczałoby przy zmianie zleconej z pulpitu — drugą, niepełną prawdę
o procesie; dlatego rejestracja warstwy mobilnej nie przyjmuje nadajnika
zdarzeń.

## budowa/server/internal/core/adapter_narzedzia_media_przetworzenie.go
Każda czynność wytwarza nowy zasób, źródło zostaje nietknięte: narzędzie zewnętrzne nigdy nie dostaje ścieżki źródła jako celu zapisu. Wejściem jest blob magazynu, plik pod sumą kontrolną, którego nadpisanie zerwałoby tożsamość wszystkich zasobów o tej treści, wyjściem plik katalogu tymczasowego.

Pomiar idzie przed przetworzeniem, bo daje trzy rozstrzygnięcia, których inaczej trzeba by zgadywać: kontener domyślny, gdy wołający nie podał formatu, kodek dźwięku przy wyodrębnianiu ścieżki oraz to, czy materiał w ogóle niesie dźwięk. Zgadnięte, każde z nich kończy się plikiem pustym albo odmową samego binarium, czyli odmową bez nazwy braku.

Parametru, którego czynność wymaga, plik się nie domyśla: operacja przycięcia bez granic daje odmowę nazywającą brak, a nie kopię materiału pod nazwą fragmentu.

Sprawdzenie rozmiaru wyniku wyprzedza pomiar wyniku, bo narzędzie zewnętrzne potrafi zakończyć się powodzeniem, zostawiając plik pusty, gdy wskazany fragment leży poza materiałem, a zasobu bez bajtów magazyn nie przyjmuje.

Zrzut klatki jest obrazem: jego wymiary czyta nagłówek pliku, tą samą drogą, którą mierzy je wniesienie zasobu, a pole czasu trwania zostaje puste, bo obraz nie trwa. Każdy inny wynik jest materiałem czasowym: mierzy go program pomiaru, a wymiary są brane ze strumienia obrazu, którego wyodrębniony dźwięk nie ma.

Rozstrzygnięcie kontenera wyniku przy braku wskazania wołającego: operacja konwersji bez formatu jest odmową, bo bez wskazania program przepisałby materiał do tego samego kontenera i oddał kopię pod nazwą przekształcenia; zrzut klatki spada na PNG; wyodrębnienie dźwięku spada na kontener wyprowadzony z kodeka zmierzonego w materiale, więc ścieżkę da się przepisać bez ponownego kodowania; przycięcie i zmiana rozmiaru zostają w kontenerze źródła, bo żadna z tych czynności nie jest zmianą formatu.

## budowa/server/internal/core/handlers_model.go
Rodzina model.* weszła do kontraktu osobno i osobno się wpina, dokładnie tak,
jak rodzina komend pamięci wpina się osobno od rodziny komend przestrzeni
roboczej, jadąc na tej samej maszynerii rejestracji. Zdarzenie rozgłaszane po
zmianie kanału jest tym samym zdarzeniem, które rozgłasza aktualizacja okna:
zmiana kanału jest zmianą okna, więc panel sterowania odświeża się z tej
samej subskrypcji, a drugiego zdarzenia dla tego samego faktu kontrakt nie
ma. Żądanie wskazujące kartę sesji dotyka wielu okien naraz i rozgłasza tyle
zmian, ile okien naprawdę zmieniło kanał. Asercja portu okien przy rejestracji
jest celowo twarda, wzorem rodziny pamięci: cicha nieobecność uchwytu jest
usterką widoczną dopiero na uruchomionym produkcie, więc rdzeń ma paść głośno
przy montażu, zamiast po cichu zostawić komendę nieznaną.

## budowa/server/internal/core/handlers_monitor.go
Monitor procesów jest oknem warstwy wspólnej obsługiwanym przez port
nawigacji, więc telemetria ma w rdzeniu jednego właściciela; wpięcie idzie
osobno, bo rodzina monitor.* weszła do kontraktu osobno, tak samo jak
rozszerzenie pamięci wpina się osobno od portu modułu przestrzeni roboczej.
Pozostałe rodziny komend przy porcie niewypełnionym nie rejestrują niczego
i rdzeń odpowiada komendą nieznaną; komendy monitora rejestrują się zawsze,
a port bez telemetrii odpowiada błędem wewnętrznym z nazwą brakującego bytu,
żeby odpowiedź brzmiała „monitor nie jest wpięty", a nie „nie znam takiej
komendy" ani, najgorzej, pusty wykaz procesów. Rodzina nie rozgłasza zdarzeń,
bo telemetria postępu ma w rdzeniu jednego producenta zdarzenia zmiany
postępu; odczyt stanu niczego nie zmienia, a zapis obserwacji zmienia pamięć
rdzenia, dla której kontrakt zdarzenia nie ma.

## budowa/server/internal/core/adapter_narzedzia_obraz_model_port.go
Port jest osobny od portu ImageMagicka mimo wspólnego przedrostka nazw komend. Tamten stoi wyłącznie na ImageMagicku i wywraca się na braku jednego binarium; ten stoi na dwóch silnikach neuronowych, z których każdy może być nieobecny osobno. Wpięcie tych metod do tamtego interfejsu związałoby dostępność sześciu komend w jedno wspólne rozstrzygnięcie dostępności, choćby maszyna z ImageMagickiem i bez sieci powiększającej straciła retusz razem z powiększaniem.

Ta para komend nie ma zdarzeń, więc port nie bierze nadajnika: wytworzony zasób jest zasobem modułu Design i to jego rodzina zdarzeń o zasobach mówi.

Port niewypełniony nie rejestruje niczego: obie komendy odpowiedzą wtedy odmową nieznanej komendy, a pozostałe domeny pracują bez zmian.

Adapter wypełnia port w całości; gdyby port i adapter się rozjechały, kompilacja stanie w tym miejscu, a nie dopiero na martwej komendzie.

## budowa/server/internal/core/handlers_mowa.go
Adapter wraz z rozstrzygnięciami, skąd bierze się okno, zasady i obszar oraz
jak typowane odmowy pakietu mowy przekładają się na kody kontraktu, leży
w adapterze modułu mowy; sam silnik leży w pakiecie mowy rdzenia. Dziewiąta
pozycja o przedrostku speech, odpowiedź komendy nieznanej obszaru, nie jest
komendą: nie ma pary żądanie-wynik i w rejestrze się nie zjawia. Wszystkie
osiem komend stoi na jednym porcie, bo stoją na jednym adapterze i na tym
samym silniku; rozdzielenie ich na dwa porty rozdzieliłoby też stan, którego
rozdzielić nie wolno, bo nasłuch ciągły rozpoznaje odcinki przysłane komendą
przyjęcia nagrania, więc obie muszą widzieć ten sam rejestr nasłuchów.
Synteza mowy tu nie należy: komenda syntezy tłumaczenia idzie w drugą
stronę, z tekstu na dźwięk, i obsługuje ją port tłumaczenia — wspólny
przedrostek w nazwie komend jest zbieżnością słowa, nie jednej maszynerii,
bo rdzeń rozpoznaje mowę i nadal jej nie syntezuje. Rodzina ma zdarzenia
wyłącznie w nasłuchu ciągłym: zdarzenie niesie rozpoznany odcinek albo
wykrycie frazy wybudzającej, rozgłaszane przez adapter przy przyjęciu
odcinka nadajnikiem wpiętym osobną metodą, a nie przez obsługę komend, bo
odpowiedź komendy przyjęcia nagrania mówi wyłącznie o przyjęciu bajtów, a nie
o tym, co w nich usłyszano. Port niewypełniony nie rejestruje niczego: obie
komendy sprawdzenia gotowości i transkrypcji odpowiedzą wtedy komendą
nieznaną, a pozostałe domeny pracują bez zmian.

## budowa/server/internal/core/adapter_narzedzia_obraz_raster.go
Dekodery są rejestrowane importem pobocznym: PNG, JPEG i GIF stoją w bibliotece standardowej, a WEBP dokłada bibliotekę rozszerzeń obrazu. To ten sam zestaw, którym rdzeń mierzy wymiary wnoszonego zasobu; dwa różne zestawy dekoderów oznaczałyby, że panel mierzy plik, którego złożenie odmawia.

Zdjęcie przemnożenia jest konieczne, nie kosmetyczne: kolor w bibliotece standardowej oddaje składowe już przemnożone przez alfę, a tryby mieszania są zdefiniowane na barwie własnej piksela. Mnożenie barw przemnożonych dałoby wynik ciemniejszy przy każdej półprzezroczystości, błąd niewidoczny na krawędziach, a wyraźny na dużej płaszczyźnie znaku wodnego.

Przycięcie składowej do zakresu jest konieczne, bo tryby mieszania screen i overlay potrafią wyjść nieznacznie poza jedynkę na zaokrągleniach.

## budowa/server/internal/core/handlers_narzedzia_archiwum.go
Adapter leży w osobnym pliku adaptera narzędzi archiwum; pakowanie
i rozpakowanie leżą w plikach czynności, a wyrok o zawartości archiwum leży
w pliku spisu. Trzecia pozycja o przedrostku archive, odpowiedź komendy
nieznanej obszaru, nie jest komendą: nie ma pary żądanie-wynik i w rejestrze
komend się nie zjawia. Rodzina nie ma zdarzeń, więc port nie bierze nadajnika:
archiwum wytworzone przez komendę spakowania jest zasobem modułu projektowania
i to jego rodzina zdarzeń o zasobach mówi. Port niewypełniony nie rejestruje
niczego — obie komendy odpowiedzą wtedy komendą nieznaną, a pozostałe domeny
pracują bez zmian; rdzeń niczym nie warunkuje startu.

## budowa/server/internal/core/handlers_narzedzia_dokument.go
Adapter wraz z rozstrzygnięciami, skąd biorą się zasady i obszar uruchomienia,
którym binarium jedzie która droga i dlaczego znacznik użycia rozpoznawania
znaków mówi prawdę, leży w plikach adaptera narzędzi dokumentu. Są to
narzędzia modelu, nie panel właściciela konta — model wykonuje je sam w
trakcie tury, dlatego żadna nie niesie okna i żadna niczego nie rozgłasza.
Rodzina nie ma zdarzeń: kontrakt nie zna zdarzenia zmiany dokumentu, więc port
nie bierze nadajnika, bo zdarzenie spoza kontraktu dałoby klientowi nazwę,
której nie zna nikt poza rdzeniem. Port niewypełniony nie rejestruje niczego:
obie komendy odpowiedzą wtedy komendą nieznaną, a pozostałe domeny pracują
bez zmian.

## budowa/server/internal/core/handlers_narzedzia_media.go
Adapter wraz z rozstrzygnięciami leży w plikach adaptera narzędzi mediów,
a jedyna droga wołania binarium leży w pakiecie wywołań zewnętrznych. Trzeciej
komendy rodziny nie ma i rdzeń jej nie wymyśli: nazwa spoza kontraktu byłaby
nazwą, której nie zna nikt poza rdzeniem. Rodzina nie ma zdarzeń: kontrakt
nie zna zdarzenia zmiany mediów, więc żadna z komend niczego nie rozgłasza
i port nie bierze nadajnika, mimo że przetworzenie zakłada zasób —
rozgłoszenie własnego zdarzenia dałoby klientowi kopertę, której nie zna jego
strona kontraktu. Port niewypełniony nie rejestruje niczego: obie komendy
odpowiedzą wtedy komendą nieznaną, a pozostałe domeny pracują bez zmian.

## budowa/server/internal/core/handlers_narzedzia_obraz.go
Adapter wraz z rozstrzygnięciami, skąd bierze się źródło, jak woła się
binarium i gdzie ląduje wynik, leży w pliku adaptera narzędzi obrazu, a
składanie argumentów każdej operacji leży w osobnym pliku czynności. Piąta
pozycja o przedrostku image, odpowiedź komendy nieznanej obszaru, nie jest
komendą: nie ma pary żądanie-wynik i w rejestrze się nie zjawia. Rodzina nie
ma zdarzeń, więc port nie bierze nadajnika: wytworzony zasób jest zasobem
modułu projektowania i mówi o nim jego rodzina zdarzeń. Port niewypełniony
nie rejestruje niczego — komendy odpowiedzą wtedy komendą nieznaną, a
pozostałe domeny pracują bez zmian, rdzeń niczym nie warunkuje startu. Dwie
ostatnie komendy rodziny, złożenie i wektoryzacja, nie wołają ani jednego
programu zewnętrznego: pracują bibliotekami wkompilowanymi w binarium
rdzenia, lecz stoją na tym samym porcie, co reszta rodziny, bo dzielą z nią
wszystko poza sposobem liczenia — osobny port dałby drugą prawdę o tym, gdzie
rdzeń odkłada bajty obrazu.

## budowa/server/internal/core/handlers_narzedzia_zakresy.go
Port jest osobny od portu narzędzi sesji: tamten opisuje dołożenie narzędzia na czas sesji i wykaz
po ukośniku, ten zakres i limit wywołań właściwy profilowi — dwa różne pytania nad tym samym
katalogiem. Kontrakt nie daje rodzinie komend zakresu żadnego zdarzenia, więc rdzeń go nie
wymyśla, a okno odświeża wykaz po odpowiedzi.

Straż stoi po stronie odbiorcy: pyta ją rdzeń przed skierowaniem komendy do obsługiwacza. Straż
niewpięta nie zmienia niczego, bo stanem wyjściowym platformy jest pełny dostęp bez granicy.

Pozycja bez wiersza zakresu przechodzi bez zapytania i bez rachunku, bo stanem wyjściowym jest
pełny dostęp, a rachunek prowadzony dla wszystkich pozycji dopisywałby wiersz przy każdym
wywołaniu narzędzia w produkcie.

## budowa/server/internal/core/handlers_nawigacja.go
Ciąg jest jeden: strona główna, wykaz i wejście środowiska, wykaz modułów, wejście przestrzeni
roboczej, stan okna komunikacji. Klient może wejść w dowolnym miejscu, na przykład wprost do
środowiska zapamiętanego z poprzedniej pracy. Brak podłączonej domeny nie wywraca rdzenia: komenda
nawigacji odpowie wtedy kodem nieznanej pozycji, a pozostałe domeny pracują dalej.

## budowa/server/internal/core/handlers_orchestration.go
Rodzina wpina się osobno od uchwytów modułu Automations, choć jedzie na tej samej maszynerii
układu zależności: port rozszerza port Automatyki, bo układ zależności ma w rdzeniu jednego
właściciela. Rozdział miejsc emisji zdarzeń układu: zapis i usunięcie zależności emitują zdarzenie
tutaj; przepisanie całego układu automatyki emituje zdarzenie aktualizacji bez wskazania
zależności, bo pole zależności zostaje wtedy puste; ustawienie harmonogramu automatyki emituje
zdarzenie zmiany powiązania z identyfikatorem automatyki równym identyfikatorowi przebiegu. Dwa
ostatnie przypadki obsługuje moduł Automations.

Nieudany odczyt układu przed zapisem zależności nie wstrzymuje niczego — komenda wykonuje się tak
samo, a rodzaj zmiany schodzi wtedy do wartości „zaktualizowano", bo stwierdzenie że układ się
zmienił jest prawdziwe w obu przypadkach.

Puste wskazanie łuku po rozgłoszeniu zmiany układu jest stanem zamierzonym: pole zależności jest
nieobowiązkowe, bo przepisanie całego układu nie dotyczy żadnego jednego łuku — okno czyta wtedy
zdarzenie jako polecenie przeliczenia układu od nowa.

## budowa/server/internal/core/nosniki_druku_wspolne_test.go

Sprawdziany pilnują, żeby wykaz nośników druku Designu i wykaz nośników druku
Studia nie rozjechały się z wykazem wspólnym. Wcześniej wykazy stały obok
siebie osobno i rozjechały się już raz: koperta DL miała w wykazie Designu
99 na 210 mm (wymiar wkładki), a w wykazie okna Studia 220 na 110 mm, podczas
gdy koperta DL według normy ISO 269 ma 110 na 220 mm. Pilnowanie w tym pliku
nie jest przewidywaniem ryzyka, tylko zapisem tego, co się już zdarzyło.

TestWykazDesignuJestPrzekladem sprawdza zgodność liczby pozycji, wymiarów
i rodziny między wykazem wspólnym a przekładem używanym przez obszar druku
Designu. Poprzednik tego sprawdzianu porównywał dwa osobne wykazy stojące
obok siebie i wykrył przy pierwszym przebiegu dwa rozjazdy: kopertę DL
(99 na 210 mm zamiast 110 na 220 mm z normy ISO 269) oraz wizytówkę zapisaną
poziomo wbrew założeniu pionowego zapisu. Wykaz własny Designu został potem
usunięty na rzecz przekładu z wykazu wspólnego, więc test sprawdza teraz, że
przekład oddaje każdą pozycję wykazu wspólnego z tymi samymi wymiarami
i tą samą rodziną, a nie tylko podzbiór pozycji. Pomiar liczby pozycji ma
znaczenie praktyczne: przekład, który gubi pozycję, odbiera ją korzystającemu
z interfejsu po cichu, bo wykaz nadal wygląda na kompletny.

Koperty C4, C5, C6, B4 i B5 zostały dodane do wykazu wspólnego i muszą dojść
również do wykazu Designu — wcześniej ich tam nie było, co było powodem tego
przestawienia na wspólne źródło.

## budowa/server/internal/core/handlers_orkiestracja.go
Rodzina ma port osobny, nie rozszerzenie portu Automatyk: podagent nie jest krokiem automatyki ani
jej przebiegiem, wisi na oknie wykonawcy, a nie na zapisanej definicji, i powstaje w rozmowie, nie
w konstruktorze przebiegów. Zdarzenie zmiany podagenta nie rozgłasza ten plik, tylko adapter
rozgłoszenia: podagent zmienia stan głównie poza żądaniem — wejście w stan działania, zakończenie
i niepowodzenie dzieją się w pracy puszczonej w tle, której uchwyt komendy nie widzi. Dlatego
funkcja wpinająca nie bierze nadawcy.

## budowa/server/internal/core/przeklad.go

Pakiet session jest właścicielem pojęcia okna komunikacji i jego cyklu życia,
więc trzyma własne struktury opisujące sesję i okno. Kontrakt trzyma własne,
niezależne struktury o tym samym znaczeniu. Ten plik jest jedynym miejscem,
w którym jedna postać przechodzi w drugą — przekład nie jest rozproszony po
obsługiwaczach żądań.

Pole eksperta w ustawieniach okna wraca do klienta wyłącznie wtedy, gdy jest
ustawione: pole puste oznacza model surowy, a wskaźnik na pusty napis
wyrażałby to samo znaczenie drugim sposobem, co wprowadzałoby dwie
reprezentacje tego samego stanu.

Pole eksperta w zmianie okna przechodzi do struktury zmiany wybiórczej
pakietu sesji, więc wybór eksperta zapisuje się tak samo jak każde inne
ustawienie okna, zamiast być tracone po drodze przy przekładzie.

Funkcja listaKatalogow zamienia wycinek pusty (nil) na tablicę pustą, ponieważ
kontrakt zapowiada pole `workingDirs` bezwarunkowo, a wycinek pusty w Go
koduje się do wartości `null`. Okno bez katalogów roboczych jest stanem
poprawnym, więc kontrakt nie może wymagać od odbiorcy przygotowania się na
brak pola — klient czytający to pole bez osłony wywróciłby wczytywanie modułu
przy wartości `null`. Wysłanie pustej tablicy usuwa tę różnicę.

## budowa/server/internal/core/handlers_poczta.go
Siedem komend kontrakt wystawia modelowi jako narzędzia: wykaz skrzynek, wykaz folderów,
odnalezienie listu, odczytanie go w całości, zapisanie szkicu, wysłanie i oznaczenie. Podpięcie,
rozpoznanie i odpięcie skrzynki narzędziami nie są, bo rozstrzygają, do czego platforma ma dostęp —
wpinają się tak samo jak reszta, różnica leży w tym, że nie ma ich w wykazie narzędzi kontraktu.
Zdarzenie zmiany idzie po trzech komendach, nie po dziesięciu: wykazy niczego nie zmieniają,
a podpięcie i odpięcie skrzynki zmieniają katalog skrzynek, dla którego kontrakt osobnego zdarzenia
nie ma; rozgłaszanie ich zdarzeniem o wiadomości donosiłoby oknu o zmianie bytu, który się nie
zmienił.

Wtopienie poczty w Asystenta odebrałoby ją każdemu innemu oknu, bo Asystent prowadzi zlecenie,
a skrzynka nie jest jego bytem.

## budowa/server/internal/core/przenoszenie_komplet.go

Komplet kontekstu ma siedem składników: polecenie wyjściowe, dokumenty,
projekt, agentów, historię rozmowy, źródła wiedzy i parametry wykonania.
Żaden z nich nie może zginąć po drodze, więc komplet powstaje z trzech
warstw nakładanych w tej kolejności: to, co przyniosło żądanie, potem
komplet zapisany przy oknie źródłowym, a na końcu realny stan okna i sesji
źródłowej. Warstwa wcześniejsza wygrywa: składnik wskazany wprost w żądaniu
nie zostaje nadpisany tym, co system odczytał sam.

## budowa/server/internal/core/handlers_prowenancja.go
Rodziny są dwie, magazyn jeden. Provenance Explorer pyta o pojedyncze wywołanie, rozliczenie liczy
sumy po wymiarze, ale obie odpowiedzi powstają z tych samych wierszy — druga tabela z tymi samymi
liczbami rozjechałaby się z pierwszą przy pierwszej korekcie cennika. Port jest osobny od portu
modułu Diagnostics, mimo że to jego okna po niego sięgają: rdzeń nie ma prawa wiedzieć, że istnieje
moduł Diagnostics. Ślad wywołania czyta też okno pętli wykonania, a alerty sięgają po niego z
ekranu stałej obecności.

Powtórzenie wywołania czekało na warstwę, która kanały prowadzi, i wchodzi razem z nią: adapter
bierze ten sam rejestr kanałów, którym jedzie okno rozmowy. Rejestr niewpięty daje odmowę nazywającą
brak, a nie pusty uchwyt udający zdolność.

## budowa/server/internal/core/przenoszenie_magazyn.go

Komplet kontekstu nie ma własnej tabeli, bo nie ma jej też żaden z jego
składników: dokument, agent i źródło wiedzy są w kontrakcie identyfikatorami,
a schemat bazy ich nie zna. Magazyn kładzie więc komplet tam, gdzie okno już
ma swój stan trwały — na najwęższym poziomie zasięgu konfiguracji, pod jednym
kluczem i w rodzaju `json`. Ta sama droga czyni komplet czytelnym dla klienta
zwykłym odczytem konfiguracji na poziomie okna, bez drugiej komendy i bez
drugiego magazynu. Pamięć procesu działa jako bufor: odpowiada bez odpytywania
bazy i przejmuje magazyn, gdy zapis albo odczyt trwałości zawiedzie, ponieważ
awaria trwałości nie ma prawa odmówić przeniesienia kontekstu.

## budowa/server/internal/core/handlers_przejecie_sterowania.go
Kształty stoją tutaj, a nie w warstwie kontraktu, bo generowany kontrakt nie zna jeszcze żadnej
z tych trzech nazw. Po wniesieniu komend do kontraktu typy stają się aliasami kontraktowymi,
a katalog wartości wraca do kontraktu. Nazwy komend są argumentem, a nie literałem: rejestr nie
zawiera nazw własnych, wstrzykuje je montaż. Nazwa pusta niczego nie rejestruje, więc dopóki
kontrakt nie niesie tych komend, rdzeń nie ogłasza zdolności, której kontrakt nie zna. Rodzina nie
bierze nadajnika i nie rozgłasza własnego zdarzenia: stan biegu wychodzi już rodzinami zmiany stanu
okna i postępu oraz komendą odczytu stanu okna, a przejęcie zmienia właśnie ten stan — funkcja
przejęcia zatrzymuje pętlę, funkcja oddania ją wznawia, a każda z nich rozgłasza stan biegu
obserwatorom pętli.

Zamiast odpowiedzi nieznanej komendy albo cichej zgody na przejęcie, które się nie odbyło.

## budowa/server/internal/core/przenoszenie_parametry.go

Kontrakt zostawia pole executionParams surowym JSON-em, bo nie każdy moduł
wykonuje pracę tak samo, podczas gdy okno komunikacji ma ustalony zestaw
parametrów wykonania i to on tędy jedzie.

Pole agenta jedzie razem z kanałem modelu, bo przekazanie kontekstu przenosi
całe stanowisko, nie samą pracę. Okno docelowe ma pracować tą samą
tożsamością — inaczej ta sama wypowiedź trafiłaby do modelu z innym
promptem systemowym i innym modelem bazowym, a korzystający z interfejsu
zobaczyłby rozjazd bez przyczyny.

Przy zmianie okna docelowego rola okna pozostaje nietknięta, ponieważ rola
wiąże okno z koordynatorem pętli, a przekazanie kontekstu nie jest zmianą
układu pętli.

## budowa/server/internal/core/handlers_queue.go
Działanie „powtórz" nie ma limitu obiegów. Rdzeń nie zlicza prób i nie odmawia po którejś z kolei —
przerwanie należy do użytkownika.

## budowa/server/internal/core/handlers_queue_wiazania.go
Zdarzenie rozgłasza wyłącznie wiązanie kolejki. Wiązanie zmienia kolejkę, więc idzie tym samym
zdarzeniem zmiany kolejki, co założenie i działanie — rodzajem zmiany „zaktualizowano". Wykaz
niczego nie zmienia i niczego nie rozgłasza; rozgłoszony jako zmiana byłby zdarzeniem bez faktu.

Port kolejek bez tych dwóch czynności to co innego: kolejki są, a rdzeń nie umie ich oddać. Wtedy
komendy zostają wpięte i odmawiają wprost błędem wewnętrznym, bo odpowiedź nieznanej komendy
wskazywałaby na brak kolejek, a nie na usterkę montażu.

## budowa/server/internal/core/handlers_role.go
Zmiana roli rozgłasza się dwoma zdarzeniami i każde ma innego odbiorcę: zmiana okna odświeża okno
w wykazie Mission Control, zmiana roli niesie samo nadanie wraz z więzią koordynatora — tego
drugiego panel ról nie złoży z pierwszego, bo okno nie niesie wcielenia roli.

## budowa/server/internal/core/rdzen.go

Pole niepowodzenia rdzenia: brak podłączonego obserwatora nie zmienia
zachowania rdzenia, znika wyłącznie zapis odmowy w module Diagnostics.

Pole straży zakresów narzędzi trzyma zakres zapisany dla profilu asystenta
w pliku handlers_narzedzia_zakresy.go. Niewpięta straż nie zmienia niczego:
stanem wyjściowym platformy jest pełny dostęp bez granicy, a zawężenia po
prostu wtedy nie ma.

Pole więzi połączeń jest opisane w pliku wiez_polaczenia.go; wartość zerowa
znosi się sama, bo wszystkie metody więzi przyjmują odbiornik zerowy.

WykonajZadanie to droga, którą wchodzi transport po samodzielnym odkodowaniu
koperty. Rozpoznanie warstwy niższej opiera się na całym kontrakcie, a rdzeń
obsługuje tylko to, co ma wpięte, dlatego komenda bez obsługiwacza jest tu
rozpoznawana ponownie — inaczej odpowiedź „*.unknown" wróciłaby pod nazwą
komendy zamiast pod nazwą zdarzenia obszaru.

zDziennikiemRdzenia wprowadza kontekst z dziennikiem raz, w jedynym gardle
każdego żądania, więc warstwy niższe nie muszą sobie go podawać ręcznie.

odmowaZakresu pyta straż wyłącznie o rękę modelu. Zakres opisuje, jak szeroko
działa asystent w imieniu Operatora, a nie co wolno samemu Operatorowi —
zasada opisana w pliku sprawca.go. Praca własna rdzenia i połączenie, które
się nie przywitało, też przechodzą, bo zawężenie nałożone na przemiatanie
zatrzymywałoby platformę bez decyzji Operatora.

## budowa/server/internal/core/handlers_schedule.go
Odczyt niczego nie rozgłasza: komenda odczytu harmonogramu nie zmienia stanu, więc nie dostaje
emitera i nie wysyła zdarzenia. Zmianę harmonogramu niesie komenda ustawienia harmonogramu
automatyki i to ona odpowiada za rozgłoszenie — zdarzenie po odczycie byłoby szumem, na który okna
reagowałyby odświeżeniem bez powodu.

Trzy pozostałe komendy — historia wyzwoleń, adres webhooka i odczyt harmonogramów — niczego nie
zmieniają i niczego nie rozgłaszają. Uruchomienie wsteczne rozgłasza to samo zdarzenie zmiany
powiązania, a nie stan przebiegu: zakłada wiele przebiegów naraz, a każdy z nich rozgłasza swój
stan sam, drogą kolejki.

## budowa/server/internal/core/rejestr.go

Komenda spoza zbioru rzeczywiście obsługiwanych dostaje odpowiedź
„*.unknown" zamiast błędu zrywającego — dotyczy to również komendy, która
jest w kontrakcie, lecz nie ma jeszcze obsługiwacza.

## budowa/server/internal/core/handlers_schowek_skroty.go
Trzy porty, nie jeden: schowek, słownik skrótów i konteksty pamięci mają osobne magazyny i osobne
powody do awarii. Wspólny port związałby ich dostępność w jedno „jest albo nie ma" — maszyna ze
słownikiem i bez historii schowka straciłaby rozwijanie skrótów razem z historią.

Zdarzeń żadna z tych rodzin nie ma. Kontrakt zna zdarzenie zmiany pamięci, ale dotyczy ono wpisów
pamięci przestrzeni roboczej i rozgłasza je ta domena; kontekst jest zestawem wskazań, a nie
wpisem, więc rozgłaszanie go pod tą samą nazwą kazałoby oknu odświeżyć wykaz faktów po zmianie,
która żadnego faktu nie dotknęła.

## budowa/server/internal/core/handlers_sesja_konfiguracja.go
To jest ta sama konfiguracja co rodzina config.*, nie drugi rejestr obok niej: konfiguracja sesji
jest innym kształtem wartości w tym samym rejestrze ośmiu poziomów zasięgu i trzech osi. Dlatego
port konfiguracji sesji wchodzi w skład portu Ustawienia, a nie obok niego — rdzeń ma jedną bramę
do konfiguracji i jedno miejsce rozgłaszania zdarzenia zmiany. Zapis rozgłasza zmianę obszaru.
Każdy dotknięty obszar wraca wpisem rezolwera i idzie zdarzeniem zmiany konfiguracji — obszar
zapisany jako zmieniony, obszar wyczyszczony jako usunięty, bo skasowanie zapisu przywraca
dziedziczenie. Klient nie musi odpytywać poziomu, żeby dowiedzieć się o zmianie.

## budowa/server/internal/core/handlers_session.go
Okna komunikacji sesji mają własną domenę. Brak podłączonej domeny nie wywraca rdzenia — komendy
sesji odpowiedzą wtedy kodem nieznanej pozycji, a pozostałe domeny pracują dalej.

Rejestracja historii sesji stoi osobno od rejestracji domeny sesji, bo to inny rodzaj czynności:
tamta prowadzi sesję w pracy bieżącej, ta porządkuje historię. Każda rozgłasza zmianę, żeby wykaz
na pozostałych urządzeniach konta przestawił się bez odpytywania.

## budowa/server/internal/core/rejestrator_blokow.go

Jedynym miejscem, które widzi każdy fragment strumienia przed spakowaniem
w kopertę, jest ujście tury w adapterze rozmowy (funkcja prowadzTure
w pliku adapter_rozmowa.go). Dziennik rozmowy i utrwalacz widzą dopiero
domkniętą wiadomość, czyli sam sklejony tekst; wpięcie rejestracji bloków
w tym miejscu wymagałoby przenoszenia bloków przez typ Message, którego
kontrakt nie modeluje. Rejestrator jest więc portem adaptera, wzorem
odbiornika zdarzeń wykonawczych: adapter woła jedną metodę, montaż podaje
całość.

Rejestrator zapisuje wszystkie rodzaje fragmentów poza tekstem: rozumowanie,
wywołania narzędzi z wynikami, obraz, dźwięk, błąd, prowenancję i metadane
konta. Tekst pomija świadomie, bo treść tekstową domyka dziennik rozmowy
w polu wiadomosc.tresc, a druga kopia byłaby drugą prawdą.

Zapis bloków korzysta z tabeli blok_wiadomosci wprowadzonej migracją
store/migracja_096_tresc_rozmowy.sql. Okno, którego zapis bloków zawiódł,
schodzi z rejestracji do końca życia procesu — wzorem degradacji dziennika
rozmowy: jedno zgłoszenie do dziennika rdzenia, zero powtórek.

## budowa/server/internal/core/handlers_terminal.go
Jedno zdarzenie na cały moduł. Kontrakt daje modułowi wyłącznie zdarzenie zmiany procesu, więc
każda zmiana stanu procesu — uruchomienie, zakończenie, ubicie — rozgłasza się tym samym
zdarzeniem. Monitor procesów odświeża się z jednej subskrypcji, a nie z odpytywania. Wyjście
procesu nie idzie tą drogą: strumień konsoli wyjścia jedzie zdarzeniem fragmentu strumienia, bo
trwa dłużej niż wykonanie komendy i bo kontrakt nie ma zapowiadanego zdarzenia strumienia wyjścia
terminala. Składa go adapter strumienia terminala.

Rozgłoszeń po wykonaniu komendy i zabiciu procesu nie ma tutaj z zamysłem: zdarzenie nadaje
adapter, bo tylko on wie, kiedy proces naprawdę zmienił stan. Podwójne rozgłoszenie z obsługiwacza
dałoby monitorowi procesów ten sam wiersz dwa razy.

## budowa/server/internal/core/sesja_konfiguracja.go

Obszar jest jednostką zapisu konfiguracji sesji. Komendy config.session.*
biorą wykaz obszarów, nie wykaz pól, więc jeden obszar odpowiada dokładnie
jednemu wierszowi tabeli ustawienie pod adresem złożonym z poziomu zasięgu
i osi. Dzięki temu obszar zapisany na poziomie węższym przykrywa obszar
poziomu szerszego w całości, a obszar spoza wykazu obszarów pozostaje
nietknięty.

Wykaz obszarów nie jest listą stałych w kodzie — powstaje z nazw pól
kontraktu shared.SessionConfig. Dopisanie obszaru do kontraktu wystarcza:
rdzeń pozna go bez zmiany ani jednej gałęzi kodu. Kolejność obszarów jest
kolejnością pól kontraktu, więc to samo wejście daje ten sam wynik przy
każdym wywołaniu. Rozbiór i złożenie konfiguracji idą przez kodowanie JSON
kontraktu, a nie przez ręczne przypisania pól, więc drugiego opisu obszarów
w rdzeniu nie ma.

## budowa/server/internal/core/handlers_terminal_wyjscie.go
Żadna z tych komend nie rozgłasza zdarzenia. Kontrakt daje modułowi jedno zdarzenie zmiany
procesu, i dotyczy ono procesu, nie zapisu na strumień. Wiersze wyjścia jadą zdarzeniem fragmentu
strumienia prosto z dziennika zbiorczego, a nie z obsługiwacza żądania: proces pisze długo po tym,
jak odpowiedź na komendę już wróciła.

## budowa/server/internal/core/handlers_terminal_wyposazenie.go
Port jest rozszerzeniem, nie drugim portem — tak samo jak wyjście: osadza port Terminal, bo
wyposażenie modułu ma tego samego właściciela co karty i procesy, a adapter wypełniający jeden
wypełnia wszystkie. Żadna z tych komend nie rozgłasza własnego zdarzenia i nie jest to przeoczenie:
kontrakt daje modułowi jedno zdarzenie zmiany procesu, dotyczące procesu. Zmiany, które procesu
dotyczą — wstrzymanie, zakończenie procesów zamykanej karty, wyzwolenie obserwacji — idą tym
zdarzeniem z adaptera. Zmiany wpisu książki hostów czy pozycji biblioteki nie mają w kontrakcie
zdarzenia, więc klient odświeża wykaz po własnym zapisie, zamiast dostawać rozgłoszenie nazwą,
której kontrakt nie zna.

## budowa/server/internal/core/sesja_konfiguracja_skladanie.go

Kolejność adresów bierze pakiet konfig (Kontekst.Adresy) — najpierw poziom,
w ramach poziomu oś. Wygrywa pierwszy adres, pod którym obszar jest zapisany,
a obszary wychodzą w kolejności pól kontraktu. Ta sama zawartość rejestru
daje więc zawsze ten sam wynik; rdzeń nie ma tu ani jednego rozstrzygnięcia
zależnego od kolejności mapy. Pierwszeństwo w ustalaniu kolejności poziomów
i osi należy do pakietu konfig — rdzeń tej kolejności nie zna i jej nie
powtarza, tylko pyta o wykaz adresów i czyta po kolei.

Rozejście między katalogiem ustawionym a faktycznie używanym składane jest
z samej konfiguracji, bez dotykania dysku: odczyt konfiguracji obowiązującej
ma być powtarzalny i nie ma prawa zakładać katalogów. Sprawdzenie dysku
należy do funkcji KatalogRoboczy na drodze uruchomienia tury, dlatego pole
checkedAt pozostaje zerowe — nic tu nie było sprawdzane.

KonfiguracjaSesjiOkna służy drodze tury: adapter rozmowy odczytuje nią
obowiązującą konfigurację i tłumaczy jej obszary na wejście procesu przez
warstwę wstrzykiwania w pliku adapter_rozmowa_konfiguracja.go. Osie modelu
i konta są pomijane w rozstrzyganiu, ponieważ obszary model i account
tłumaczą się na wybór modelu i konta wywołania — inaczej powstałaby pętla,
w której konto zależy od konfiguracji, która zależy od konta.

W kontekstKonfiguracjiSesji pole puste znaczy, że poziom albo oś nie dotyczy
tego wywołania i jest pomijana przy rozstrzyganiu.

## budowa/server/internal/core/handlers_tozsamosc.go
Silnik nakładki rozdziela tryb zastąpienia od trybu dopisania, a jego test wymusza ten rozdział.
Wartością domyślną klucza rozstrzygającego tryb podania nakładki jest tryb zastąpienia.

Zdarzenie niesie treść, bo okno konfiguracji otwarte na drugim urządzeniu ma ją pokazać bez
dopytywania. Nie jest to sekret: prompt systemowy jest zasadą pracy modelu, nie poświadczeniem.

## budowa/server/internal/core/sesja_konfiguracja_zapis.go

Dziedziczenie ośmiu poziomów zasięgu nie zachodzi w tym pliku — odpowiada za
nie plik sesja_konfiguracja_skladanie.go. Konfiguracja sesji nie jest drugim
rejestrem ustawień, tylko innym kształtem wartości w rejestrze istniejącym,
dlatego obszary nie mają własnej tabeli.

Obszar wskazany w polu areas żądania, dla którego żądanie nie niesie treści,
zostaje usunięty z danego poziomu i wtedy obowiązuje poziom szerszy. Obszar
spoza pola areas nie jest ruszany przy zapisie.

## budowa/server/internal/core/handlers_tozsamosc_osie.go
Okno wskazuje kanał modelu, a kanał identyfikator modelu i konto preferowane. Oś modelu bierze
więc identyfikator modelu kanału, a oś konta — konto kanału. Okno bez kanału, kanał bez konta
i okno nieznane dają osie puste: nakładka schodzi wtedy na samą oś platformy zamiast nie powstać
wcale.

## budowa/server/internal/core/sesja_konfiguracja_zdolnosci.go

Deklaracja zdolności może pochodzić wyłącznie od adaptera, który tłumaczy
model konfiguracji na powierzchnię dostawcy. Żaden adapter takiego wykazu
jeszcze nie wystawia — nie ma ani funkcji wejściowej, ani tablicy odwzorowań
pole-powierzchnia. Rdzeń oddaje więc deklarację pustą: wykaz obszarów i wykaz
pól są puste, bo rdzeń nie zna ani jednego rozstrzygnięcia adaptera; pole
probedAt pozostaje zerowe, bo nic nie zostało sprawdzone, a czas sprawdzenia,
którego nie było, byłby zapisem czynności niewykonanej; pola channelId
i accountId pozostają puste, bo nie sprawdzono żadnego kanału ani konta.

Wypełnienie tych wykazów zgadywanką byłoby gorsze od pustki: okno
konfiguracji pokazałoby korzystającemu z interfejsu, że pole jedzie do
dostawcy albo że nie jedzie, na podstawie niczego. Pusta prawda jest tu
jedyną dopuszczalną odpowiedzią, a brak deklaracji niczego nie wstrzymuje —
zapis konfiguracji idzie dalej.

Identyfikator adaptera jest stałą tego pakietu, bo powierzchnia CLI nie ma
w drzewie osobnego pakietu — przekład konfiguracji sesji na wejście procesu
jedzie warstwą wstrzykiwania na drodze tury, w plikach adapter_rozmowa_*.go.
Transport spoza kanału CLI nie ma w drzewie odpowiednika, więc jego
identyfikator pozostaje pusty.

Pole unsupported deklaracji zdolności nie jest wypełniane w tym miejscu, bo
znaczyłoby „dostawca nie ma odpowiednika", a nie „nie zbadano" — te dwa stany
nie mogą być mylone.

## budowa/server/internal/core/handlers_urzadzenia.go
Operator odbierający dostęp maszynie stojącej obok ma zobaczyć skutek na obu ekranach naraz —
inaczej drugi ekran pokazywałby dostęp, którego już nie ma, aż do następnej komendy.

Port, który więzi nie przyjmie, po prostu nie oznaczy bieżącego urządzenia; wykaz działa dalej.

Wykaz do zdarzenia bierze się z tego samego źródła co odpowiedź wykazu urządzeń — pozostałe ekrany
dostają stan po zmianie, a nie polecenie odpytania jeszcze raz.

## budowa/server/internal/core/skutek_materialu_z_biblioteki_test.go

Opis pola assetId w kontrakcie obiecuje w dziewięciu komendach zasób
z magazynu rdzenia, Design albo Library. Rozwiązywanie identyfikatora miało
iść przez oba magazyny: repozytorium zasobów designu oraz magazyn plików
biblioteki, tak żeby narzędzie dokumentu mogło wydobyć tekst z pliku
biblioteki i zamienić jego format, niezależnie od tego, z którego magazynu
zasób pochodzi.

Sprawdzian mierzy skutek, nie samą kopertę odpowiedzi: wnosi plik do
biblioteki, po czym pyta narzędzie dokumentu o jego treść i porównuje ją
z tym, co naprawdę weszło. Odpowiedź udana z pustym tekstem byłaby tą samą
szkodą co odmowa.

## budowa/server/internal/core/handlers_wiedza.go
Port nie rozgłasza żadnego zdarzenia zmiany, bo nie ma okna, które by je odebrało. Wskaźnik nie
odświeża się sam przy wgraniu pliku: osadzenie dokumentu to sekundy pracy procesora, a pierwsze
pobiera rząd gigabajta wag, więc wpięcie go w komendę wgrania pliku zamieniłoby ją w operację,
która czasem trwa minutę i czasem odmawia z powodu braku sieci. Budowanie wskaźnika jest czynnością
osobną i świadomą.

## budowa/server/internal/core/handlers_workspace.go
Komenda wejścia do przestrzeni nie należy do tego obszaru: przeładowuje przestrzeń roboczą karty
sesji na dowolny moduł i wpina ją nawigacja — nazwa jest wspólna, obszar nie.

Jedno zdarzenie na cały moduł. Kontrakt daje modułowi wyłącznie zdarzenie zmiany projektu, więc
każda zmiana stanu projektu — instrukcje, wpis pamięci, przypisanie eksperta — rozgłasza się tym
samym zdarzeniem. Okna modułu odświeżają się z jednej subskrypcji, a nie z czterech. Odczyty
niczego nie rozgłaszają poza jednym przypadkiem: wejście na pulpit projektu, którego jeszcze nie
było, zakłada go — a założenie bytu jest zmianą.

## budowa/server/internal/core/skutek_odczytu_pisma_designu_test.go

Do wprowadzenia tego sprawdzianu rdzeń rozpoznawał, które obszary zrzutu są
liniami tekstu, ale treści napisów nie czytał. Sprawdzian kompilacji
przechodził nad tym zielono, bo komenda oddawała ramkę i warstwy — brakowało
jedynie odczytu treści napisanej na zrzucie.

Odczyt idzie programem pakietu serwera, więc sprawdzian rozgałęzia się według
tego, czy program na maszynie stoi, i żadna z gałęzi nie jest pominięciem:
gdy program stoi, mierzony jest odczyt — adnotacja warstwy ma nieść treść
napisu, który sprawdzian sam wpisał w zrzut; gdy programu nie ma, mierzone są
dwie odmowy — żądanie z odczytem wskazanym wprost dostaje odmowę nazwaną
wraz z naprawą, a żądanie bez wskazania dostaje układ obszarów, w którym
każda linia tekstu mówi w adnotacji, że treści nie odczytano i dlaczego.
Dzięki temu plik świeci zielono zarówno na serwerze z programem, jak i na
maszynie bez niego, a w obu przypadkach mierzy zachowanie, nie samą
kompilację.

Sprawdzian wolno pytać wprost, czy program rozpoznający pismo jest na
maszynie, bo zapora obszaru Design pilnuje plików rdzenia, nie sprawdzianów,
a rozgałęzienie bez tego pomiaru musiałoby zgadywać, którą odpowiedź uznać
za poprawną.

Napis w zrzucie jest rysowany krojem wkompilowanym, tym samym, którym rdzeń
podpisuje wykresy, więc sprawdzian nie zależy od krojów zainstalowanych na
maszynie. Rozmiar pisma 20 jest zmierzony, nie dobrany na oko: przy nim
wyraz schodzi na obszar o wysokości 32 punktów, wewnątrz granicy linii
tekstu rdzenia (do sześciu kratek) i o szerokości grubo ponad trzykrotność
wysokości. Pismo większe rozsypuje się na kratce po jednej literze na
obszar i żadna nie jest już linią tekstu.

Adnotacja „linia tekstu 1: " sama w sobie, bez słowa ze zrzutu, byłaby
odczytem pustym udającym odczyt — dlatego sprawdzian szuka w adnotacji
konkretnego słowa, nie samej obecności dwukropka.

## budowa/server/internal/core/handlers_workspace_pamiec.go
Jedna zmiana wpisu rozgłasza dwa zdarzenia i nie jest to powtórzenie: zdarzenie zmiany projektu
niesie projekt i mówi oknu Workspace, że coś się w nim ruszyło, a zdarzenie zmiany pamięci niesie
sam wpis wraz z rodzajem zmiany i zasila okno pamięci kontekstu bez odpytywania wykazu. To dwa
różne byty tej samej czynności — tym samym wzorem, co zdarzenie kolejki obok stanu wykonania
automatyki. Przestawienie pamięci nie rozgłasza niczego: przestawia konfigurację karty sesji, nie
stan wpisu ani projektu, a zdarzenia dla tego bytu kontrakt nie ma. Rodzaj zmiany idzie za tym, co
się stało: zapis bez wskazania wpisu zakłada wpis, ze wskazaniem zmienia go; odpięcie zwęża zasięg
wpisu, który żyje dalej; usunięcie kasuje. Odesłanie wszystkiego jako aktualizacji kazałoby
klientowi zgadywać, czy wpis dopisać do wykazu, czy z niego zdjąć.

## budowa/server/internal/core/skutek_odtwarzania_twarzy_test.go

Sprawdzian mierzy odmowę zamiast skutku, bo skutek przebiegu twarzowego
wymaga zdjęcia twarzy, sieci liczącej minutami i trzech zestawów wag
ważących pół gigabajta. Sprawdzian tego rodzaju byłby na maszynie bez wag
pominięty, czyli świeciłby na zielono, nie mierząc niczego. Odmowa natomiast
jest zachowaniem, które ma działać wszędzie i daje się zmierzyć na pustym
katalogu. Skutek sieci na zdjęciu wykazuje się osobno, uruchomieniem na
maszynie z wagami: pole faces ustawione na false i na true nad tym samym
źródłem dają obrazy różne, a różnicę podaje się liczbą.

Odmowa mówiąca samo „brak wag" zostawiałaby korzystającego z interfejsu
z pytaniem, na które kod zna odpowiedź — dlatego test wymaga w treści
odmowy nazwy pliku, ścieżki i miejsca pochodzenia.

Program wołany przez rdzeń, ale nieobecny w wykazie zależności zewnętrznych,
jest brakiem, o którym dowiaduje się dopiero po naciśnięciu przycisku — wykaz
istnieje właśnie po to, żeby się nie dowiadywał tą drogą.

## budowa/server/internal/core/handlers_workspace_planowanie.go
Powód jest ten sam co przy pamięci projektu: rodzina liczy trzynaście komend i wpisanie ich do
portu wspólnego zrobiłoby z niego wykaz wszystkiego, co moduł umie, zamiast wykazu tego, czym jest
projekt. Zdarzeniem modułu jest jedno zdarzenie zmiany projektu — każda zmiana planu rozgłasza się
projektem. Okna huba odświeżają się z tej jednej subskrypcji, a nie z trzynastu.

## budowa/server/internal/core/skutek_terminala_wstrzymanie_test.go

Plik jest osobny, bo osobna jest jego platforma: wstrzymanie drzewa procesów
mają systemy uniksowe, a Windows nie ma dla obcego procesu odpowiednika,
zgodnie z plikiem session/wstrzymanie_windows.go. Sprawdzian pyta o skutek
w systemie, nie w rdzeniu: czyta stan procesu z /proc/<pid>/stat, gdzie
litera T znaczy zatrzymany sygnałem. Uwierzenie polu supported w odpowiedzi
byłoby uwierzeniem mierzonemu, że zrobił to, co miał zrobić.

## budowa/server/internal/core/skutek_zakresow_narzedzi_test.go

Sprawdzenie, że zakres realnie powstrzymuje wywołanie, jest sednem tego
sprawdzianu: zakres zapisany, którego rdzeń nie czyta przy wykonaniu, byłby
tylko suwakiem w oknie — wyłączenie pozycji w interfejsie nie przeszkadzałoby
modelowi wywołać jej dalej.

Pozycja bez wiersza zakresu przechodzi, ponieważ platforma niczego nie
zawęża z góry — zawężenie powstaje wyłącznie z zapisanego zakresu.

## budowa/server/internal/core/handlers_wywolywacz.go
Rodzina jest przekrojowa: wywoływacz otwiera się z dowolnego miejsca platformy, a nie z jednego
okna. Nastawa jest więc własnością rdzenia, choć samo przechwycenie klawiszy należy do powłoki
programu okiennego — adapter rozstrzyga ten podział i mówi o nim wprost w odpowiedzi. Zdarzeń
rodzina nie ma: zmiana skrótu jest zmianą nastawy, a o zmianach nastaw mówi rodzina config.*.
Osobne zdarzenie byłoby drugą drogą tej samej wiadomości. Port niewypełniony nie rejestruje
niczego: obie komendy odpowiedzą wtedy kodem nieznanej pozycji, a pozostałe domeny pracują bez
zmian.

## budowa/server/internal/core/handlers_zespoly.go
Kontrakt daje zespołom wyłącznie jedno zdarzenie zmiany, więc zapisanie nowego składu, zmiana
istniejącego i skopiowanie rozgłaszają się tą samą drogą, różniąc się wyłącznie rodzajem zmiany.
Okno składu odświeża się z jednej subskrypcji. Rodzaj zmiany rozstrzyga żądanie, nie odpowiedź:
zapis bez identyfikatora zakłada zespół, z identyfikatorem zmienia istniejący — odpowiedź w obu
przypadkach niesie ten sam kształt, więc po niej samej rozróżnić się tego nie da. Odczyt nie
rozgłasza: wykaz i wczytanie zespołu niczego nie zmieniają, więc nie mają czego ogłaszać — zdarzenie
po odczycie byłoby zawiadomieniem o zmianie, której nie było.

## budowa/server/internal/core/skutek_znacznika_bez_poczty_test.go

Znacznik auth:bramka-bez-poczty zdejmuje jeden warunek wejścia i tylko
jeden: bramki nie zamyka potwierdzenie, którego platforma nie miała czym
wysłać. Wyjątek ma trwać dokładnie tak długo, jak trwa jego powód, a powód
kończy się w chwili, w której adres zostaje potwierdzony. Znacznik, który
przeżyje potwierdzenie, przestaje być wyjątkiem pierwszego uruchomienia
i staje się trwałym obejściem bramki: od tej chwili konto cofnięte do stanu
niepotwierdzonego wchodziłoby hasłem mimo działającej poczty. Ten sprawdzian
mierzy stronę, której nie mierzy granica istnienia znacznika w pliku
skutek_wejscia_test.go: tamta prowadzi drogę z pocztą, gdzie znacznik nie
powstaje w ogóle, więc przechodzi także wtedy, gdy rdzeń znacznika nie
zdejmuje.

Bramka wpuszczająca hasłem i bramka trzymana wierszem konta oddają
w kopercie odpowiedzi to samo, dlatego dowodem w teście jest odczyt sejfu po
potwierdzeniu, nie samo pole verified.

Bramka otwarta po potwierdzeniu adresu nie zamyka wejścia, które przed nim
działało.

wydajDrogeWeryfikacji idzie warstwą danych, nie komendą, bo na instalce bez
poczty żadna komenda tej drogi nie wydaje: rejestracja wykonuje się raz
i list wyszedłby tylko z niej, a komenda auth.recover wydaje drogę do innego
celu, której auth.verify nie przyjmuje — wykazuje to test
TestBezPocztyPotwierdzenieAdresuCzekaNaDrogeZListu. Brak tej drogi nie jest
przedmiotem tego pomiaru — mierzone jest to, co rdzeń robi, gdy adres
zostaje potwierdzony. Sprawdzian trzyma sam materiał drogi, bo z bazy
odczytać się go nie da — leży tam wyłącznie jego skrót.

## budowa/server/internal/core/harmonogram_budzik.go
Okno Scheduler zapisuje cykliczność, a wyliczenie następnego uruchomienia liczy najbliższy termin.
Budzik cyklicznie pyta repozytorium o harmonogramy należne — czynne, z terminem minionym — i każdy
odpala tą samą drogą, którą automatykę rusza Operator. Przesunięcie terminu idzie przed odpaleniem:
najpierw wyliczamy i zapisujemy następny termin, dopiero potem ruszamy automatykę. Dzięki temu ten
sam harmonogram nie odpali się w kółko, gdyby odpalenie trwało dłużej niż takt zegara albo gdyby
uruchomienie zawiodło — awaria jednego przebiegu nie zawiesza budzika.

## budowa/server/internal/core/sprawca.go

Rękę sprawcy rozpoznaje się po faktach gniazda, nie po treści żądania i nie
po przedrostku napisu. Cztery rodzaje sprawcy:

- operator — gniazdo, które przedstawiło się identyfikatorem klienta
  w powitaniu i nie jest serwerem narzędzi modelu; tak wygląda okno
  interfejsu Operatora i nic innego tak nie wygląda;
- assistant — serwer narzędzi okna o roli klawiatury, gdzie rola przychodzi
  z wpisu MCP ułożonego przez rdzeń dla okna modułu Assistant (funkcja
  w pliku adapter_modul_asystent_sterowanie.go, zasięg klawiatura), więc
  jest faktem rdzenia, nie deklaracją modelu;
- model — serwer narzędzi każdego innego okna, czyli model roboczy;
- core — czynność powołana przez sam rdzeń: przemiatanie, harmonogram,
  odtworzenie stanu. Tego nie da się wyprowadzić z gniazda, bo gniazda tam
  nie ma, więc te miejsca znaczą się same przez funkcję zSprawcaRdzenia,
  zamiast być domyślane z pustki.

Brak gniazda nie znaczy core. Poza pracą własną rdzenia bez gniazda woła się
też z próby, z sondy stdio i z biegu wewnętrznego adaptera — nazwanie tego
wszystkiego rdzeniem byłoby zgadywaniem. Kontrakt mówi o polu actor wprost:
brak znaczy, że rdzeń nie potrafił tego rozstrzygnąć, więc brak jest tu
odpowiedzią, a nie luką. Gniazdo bez identyfikatora klienta też zostaje bez
sprawcy: połączenie, które jeszcze się nie przywitało, może być czymkolwiek.

Sprawca jest opisem i niczego nie rozstrzyga o tym, czy wolno. W tym pliku
nie ma ani jednej odmowy, a wynik nie wchodzi do żadnego warunku poza
wypełnieniem pola zdarzenia. Strażą rdzenia zostaje bramka wejścia, opisana
w pliku transport/bramka.go.

Znak kontekstu pracy rdzenia stawia się ręcznie, w miejscu, które wie, czym
jest — wywiedzenie core z samego braku gniazda byłoby nieprawdziwe, bo brak
gniazda ma więcej niż jeden powód.

## budowa/server/internal/core/izolacja.go
Rozstrzyganie należy do pakietu konfiguracji i tutaj się go nie powtarza; ten plik wyłącznie łączy
wynik rozstrzygnięcia z postacią wykonawczą, której używają strażnicy plików, sieci, kontekstu,
polecenia i przydziału.

Podkatalog danych modelu leży wewnątrz własnego katalogu okna, więc izolacja katalogu danych nie
wymaga drugiej podstawy.

Brak rozstrzygacza daje stan wyjściowy platformy — kontekst odrębny, żaden zakres techniczny
niewłączony — a nie odmowę wykonania.

Ustalenie bez ścieżki daje obszar pusty. Obszar pusty przy izolacji włączonej jest naruszeniem
rozstrzyganym w egzekutorze polecenia, a nie milczącym przejściem: okno bez własnego katalogu
poszłoby do katalogu wspólnego.

## budowa/server/internal/core/izolacja_kontekst.go
Wymiar odrębny znaczy, że treść należąca do innego okna nie wchodzi do tury tego okna: nowa karta
zaczyna z pustą historią, a pamięć jednego zasięgu pozostaje niewidoczna w innym. Egzekucja polega
na odrzuceniu źródła cudzego, nie na cichym pominięciu go — Operator ma wiedzieć, że tura miała
zaciągnąć treść spoza okna. Wymiar współdzielony nie ogranicza niczego; jest jawną decyzją
Operatora o tym, że ta sama treść zasila kilka zasięgów.

## budowa/server/internal/core/izolacja_pliki.go
Nadanie niesie także tryb, więc korzeń nadany do odczytu nie staje się korzeniem do zapisu. Straż
egzekwuje ścieżki, po które rdzeń sięga w imieniu okna, oraz wykaz korzeni podawany procesowi
modelu przy uruchomieniu. Nie zabroni uruchomionemu procesowi otworzyć pliku samodzielnie — to
leży poza zasięgiem rdzenia i wymaga środków systemu operacyjnego.

Zawężenie korzeni nadania do korzeni punktu robi funkcja korzeni nadania — reguła zawężenia ma
w drzewie jedno miejsce.

## budowa/server/internal/core/stan_obiegu_rejestr.go

Rejestr biegów nie prowadzi biegu, nie zatrzymuje go i nie liczy obiegów —
to robi wyłącznie pętla pakietu sesji, wywoływana przy wejściu na stronę
główną, przy stanie okna i przy telemetrii postępu. Rejestr wypełnia
interfejs session.ObserwatorObiegu, więc wpina się w pętlę tą samą drogą,
którą wpina się ślad dziennika.

Sprzątanie w metodzie Zachowaj idzie od strony okien żywych, a nie od
zdarzenia zamknięcia okna, bo rejestr biegów nie ma prawa trzymać okna przy
życiu dłużej niż rejestr nadzorcy. Wykaz pusty niczego nie kasuje, ponieważ
brak wiedzy o oknach nie jest wiedzą o ich zamknięciu.

## budowa/server/internal/core/stan_obiegu_test.go

Pole LoopState.stopReason jest jedynym polem, którym rdzeń mówi, dlaczego
bieg stanął, więc od niego zależy, czy układ złożony przez korzystającego
z interfejsu pozna maszynowo wynik pracy. Dopóki pole niosło trzy wartości,
każda znaczyła przerwane i wynik był nieodróżnialny od porzucenia.

## budowa/server/internal/core/izolacja_siec.go
Wykaz maszyn bierze się z nadań dostępu okna, nie z osobnego ustawienia — drugiego miejsca,
w którym Operator wskazywałby maszyny, w produkcie nie ma. Granica egzekucji: straż rozstrzyga
adresy, po które sięga rdzeń, oraz wykaz mostów podawany procesowi modelu w konfiguracji.
Gniazda otwierane przez sam proces modelu pozostają poza jej zasięgiem — na to potrzeba zapory
albo przestrzeni nazw sieci, czyli środka systemu, nie rdzenia.

## budowa/server/internal/core/stan_sesji.go

Start zawsze prowadzi na stronę główną, a komenda session.bind nie jest
ruchem otwierającym, tylko powrotem do sesji trwającej w tle. Powrót ma
sens wyłącznie wtedy, gdy strona główna wie, dokąd wracać i co tam się
dzieje: w którym środowisku sesja stoi, ile okien ma otwartych, w ilu trwa
strumień odpowiedzi, co widzą przez nadania dostępu i czy pracuje w nich
bieg naprawczy koordynatora. Sesje i okna zna nadzorca, bieg zna pętla,
nadania zna warstwa danych, a punkt pracy zna telemetria.
