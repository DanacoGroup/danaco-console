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
