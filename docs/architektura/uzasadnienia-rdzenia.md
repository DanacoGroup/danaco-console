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
