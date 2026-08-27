# Uzasadnienia modułów

Dokument gromadzi uzasadnienia decyzji projektowych, które nie mieszczą się
w nagłówku pliku źródłowego. Rozdział nosi nazwę pliku, którego dotyczy.

## budowa/server/internal/core/adapter_modul_extension_integracje.go

Rozmowa z serwerem MCP przy odkrywaniu narzędzi, próbnym wywołaniu i sprawdzeniu
kondycji jest rzeczywista: prowadzi ją klient MCP wkompilowany w rdzeń, po
transporcie ustawionym dla pozycji. Każda ramka rozmowy trafia do dziennika
protokołu, a każde wywołanie narzędzia — do dziennika użycia; metryki i audyt
liczą się z tych wierszy, nie z licznika w pamięci procesu.

Import definicji API czyta opis zamiast zgadywać operacje. Opis OpenAPI 3
rozkłada biblioteka `getkin/kin-openapi`, wkompilowana w binarium; schemat
GraphQL rozkłada czytnik własny, wyodrębniający pola typów `Query`
i `Mutation`. Żadna z dwóch dróg nie woła programu zewnętrznego.

Piaskownica uruchamia wyłącznie to, co integracja ma podłączone: pozycja
z transportem procesu lokalnego dostaje wejście na standardowe wejście procesu
i oddaje jego wyjście; pozycja rozmawiająca protokołem MCP dostaje wywołanie
`tools/call`. Pozycja bez jednego i drugiego wraca odmową z powodem, a nie
wynikiem udającym przebieg.

Odwołanie do poświadczenia integracji wchodzi także do rejestru referencji
sekretów rozszerzeń, ponieważ centrum uprawnień pyta o nie osobno — referencja
znana wyłącznie jednej integracji byłaby niewidoczna dla rejestru rotacji
sekretów.

Adres zgody OAuth2 powstaje wyłącznie wtedy, gdy integracja ma ustawiony adres,
do którego można pokierować przeglądarkę — adres wymyślony prowadziłby
donikąd.

Polecenie uruchamiane w piaskownicy pochodzi wprost od integracji podłączonej
rdzeniowi jako rozszerzenie; rdzeń nie składa go z danych zewnętrznych, stąd
wyłączenie ostrzeżenia gosec G204.

Log przebiegu piaskownicy ląduje w magazynie treści rdzenia, ponieważ kontrakt
oddaje do niego wyłącznie odwołanie, za którym musi leżeć zapisany plik.

Pole `output` odpowiedzi piaskownicy przyjmuje wyjście wyłącznie wtedy, gdy
jest poprawnym JSON-em, zgodnie z typem `json` zadeklarowanym w kontrakcie.

Operacje zaimportowane z opisu API lądują w tym samym wykazie narzędzi
pozycji, który wypełnia odkrywanie MCP — konektor zbudowany z opisu API jest
z punktu widzenia eksperta tym samym co narzędzie odkryte protokołem.

Adres nasłuchu webhooka przychodzącego składa rdzeń, nie integracja
zewnętrzna, ponieważ tylko rdzeń zna ścieżkę, pod którą odbiera zdarzenia.

Okno domyślne metryk użycia obejmuje ostatnią dobę, aby uniknąć sumowania
wywołań od początku istnienia pozycji.

Sprawdzenie kondycji bez wskazanego kodu pozycji obejmuje wyłącznie pozycje
włączone, aby uniknąć ruchu sieciowego do integracji wyłączonych.

Stan `degraded` sprawdzenia kondycji oznacza serwer, który odpowiedział na
powitanie, lecz nie oddał wykazu narzędzi — stan pośredni między dostępnością
a awarią.

Uprawnienia nadane pozycji rozszerzenia audyt czyta raz na pozycję, nie raz na
wpis, ponieważ audyt jednej pozycji obejmuje zwykle setki wywołań.

## budowa/server/internal/core/adapter_modul_developer_jezyk.go

Wiedzę o kodzie repozytorium niesie program serwera, nie logika rdzenia: gopls
zna typy repozytorium i graf odwołań, gofmt i prettier znają styl,
golangci-lint zna reguły analizy. Wołanie tych programów jest jedynym
miejscem tego modułu, które sięga po proces serwera — odtworzenie tej wiedzy
w kodzie rdzenia oznaczałoby napisanie drugiego kompilatora. Wszystkie stoją
na serwerze i wchodzą do sondy zależności.

Brak programu na serwerze nie zamienia odpowiedzi w odmowę. Kontrakt niesie
pola `serverAvailable` i `linterAvailable` właśnie po to: pusty wykaz symboli
przy `serverAvailable: false` znaczy brak narzędzia do sprawdzenia, a przy
`true` — sprawdzenie bez wyniku. Pierwszy stan naprawia instalacja programu na
serwerze, drugi — poprawka w kodzie.

Serwer języka jest wołany w trybie jednorazowym na każde pytanie, a nie jako
utrzymywana sesja: komenda kontraktu jest pytaniem i odpowiedzią bez stanu
między nimi, a proces utrzymywany na czas okna wymagałby własnego rejestru
uchwytów, sprzątania i przerywania — trzeciego takiego mechanizmu w module
obok budowania i debugowania.

Serwer języka TypeScriptu nie ma trybu jednorazowego wywołania, jaki ma gopls
— rozmawia wyłącznie sesją protokołu LSP na strumieniu wejścia: uzgodnienie,
otwarcie dokumentu, pytanie, zamknięcie. Jedyna droga rdzenia do procesu
zewnętrznego z zamysłu nie pisze na jego wejście, ponieważ jednoczesne pisanie
i czytanie tego samego procesu jest klasą zakleszczenia, której ten pakiet ma
nie mieć. Dlatego pliki TypeScriptu dostają odpowiedź nazywającą brak serwera
zamiast wyniku z serwera Go, który tego pliku nie rozumie; program pozostaje
zadeklarowany w wykazie zależności, więc sonda startowa zgłasza jego brak.

Formater dobiera się po rozszerzeniu pliku, bo styl należy do języka. Treść
z żądania wygrywa z treścią na dysku, ponieważ edytor formatuje bufor,
którego jeszcze nie zapisano, i to on jest przedmiotem czynności. Dla plików
Go pierwszeństwo ma goimports: robi to samo co gofmt, a dodatkowo porządkuje
wykaz importów, którego ręczne pilnowanie jest najczęstszym powodem
niekompilującego się pliku po refaktoryzacji; gdy goimports nie ma na
serwerze, formatuje gofmt.

Plik roboczy formatowania powstaje obok pliku źródłowego, nie w katalogu
tymczasowym systemu, ponieważ reguły formatowania zależą od położenia
w repozytorium — plik przeniesiony gdzie indziej sformatowałby się wedle
innych reguł niż ten, którego dotyczy czynność.

Zawężenie wyniku do zakresu bierze z formatowanej treści wyłącznie wskazane
wiersze, ponieważ formatery pracują na całym pliku, a kontrakt dopuszcza
zakres — reszta pliku zostaje nietknięta, bo zaznaczenie dotyczyło zakresu,
nie całego pliku.

Repozytorium bywa wielojęzyczne, a każdy analizator zna jeden język:
golangci-lint czyta Go, Ruff — Pythona, Stylelint — arkusze CSS, typos szuka
literówek niezależnie od języka. Zgłoszenia wszystkich programów idą do
jednego wykazu, ponieważ kontrakt niesie w każdym wpisie pole `source`.
Pole `linterAvailable` pochodzi z czasu, gdy program analizy był jeden,
i znaczy: choć jeden program odpowiedział. Fałsz tego pola jest więc stanem,
w którym pusty wykaz zgłoszeń kłamałby, bo nie sprawdzono niczym; które
programy milczały, ta odpowiedź nie rozróżnia bez zmiany kontraktu.

Dobór programów analizy idzie po rozszerzeniu plików repozytorium, nie po
tym, co stoi na maszynie, ponieważ program uruchomiony tam, gdzie nie ma ani
jednego pliku jego języka, nie milczy, tylko odmawia — golangci-lint
w repozytorium bez plików Go kończy się błędem o braku plików do analizy.
Odmowa jednego programu nie ma prawa zabrać analizy pozostałym: repozytorium
samego Pythona albo samego frontendu zostałoby wtedy bez analizy w ogóle.
Żądanie bez wskazania plików obejmuje całe drzewo i każdy program dostaje
swój katalog.

Nonzero exit code programów analizy jest tu wynikiem, nie usterką: Ruff,
Stylelint i typos kończą się niezerowo dokładnie wtedy, gdy mają co zgłosić.
Miarą przeprowadzenia analizy jest odczytana odpowiedź, a nie kod wyjścia ani
strumień, na który program ją napisał — program, który przewrócił się bez
zwrócenia czytelnego wyniku, analizy nie przeprowadził, i jego milczenie nie
ma prawa wyglądać jak brak zastrzeżeń.

Waga zgłoszenia bez wagi nadanej przez linter jest ostrzeżeniem, nie błędem:
literówka w identyfikatorze bywa nazwą celowo skróconą, a program nie ma jak
tego rozstrzygnąć — podniesienie wagi kazałoby poprawiać rzeczy, których
nikt tak nie oznaczył.

Podgląd jest domyślnym trybem refaktoryzacji. Refaktoryzacja semantyczna
dotyka wielu plików naraz, więc zmiana zapisuje się na dysk wyłącznie na
wyraźne żądanie z wyłączonym podglądem — w przeciwnym razie wraca różnica do
zaakceptowania.

Różnica zunifikowana jest jedynym kształtem, w którym serwer języka mówi
o zmianie wielu plików naraz; kontrakt niesie ją jako wykaz zmian tekstu,
więc rozbiór wyniku idzie po nagłówkach plików i fragmentach zmiany.
