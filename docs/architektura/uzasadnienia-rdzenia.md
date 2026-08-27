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
