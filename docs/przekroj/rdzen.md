# Danaco Console — przekrój pionowy: rdzeń

Opracowanie opisuje rdzeń aplikacji w postaci, w jakiej działa w bieżącej
budowie. Obejmuje punkt wejścia, układ pakietów oraz drogę komendy od koperty
do odpowiedzi.

## Punkt wejścia

Rdzeń jest programem w Go. Punkt wejścia `budowa/server/cmd/danaco-console/main.go`
zawiera wyłącznie kompozycję: wczytanie konfiguracji, przygotowanie katalogu
danych, otwarcie trwałości, złożenie rdzenia i pracę torów aż do sygnału
zatrzymania. Obok programu głównego katalog `budowa/server/cmd/` niesie
programy pomocnicze `danaco-narzedzia` i `podglad-listu`.

## Układ pakietów

Pakiety rdzenia leżą w `budowa/server/internal/`:

- `core` — obsługa komend kontraktu: adaptery modułów i pliki wpinające;
- `transport` — kanał WebSocket, opisany w opracowaniu [kanał](kanal.md);
- `store` — trwałość SQLite, opisana w opracowaniu [trwałość](trwalosc.md);
- `dane` — dostęp do danych dziedzinowych ponad warstwą trwałości;
- `session` — model sesji; `protocol` — typy protokołu; `models` — integracja
  modeli; `konfig` i `konfiguracja` — konfiguracja;
- `mail` i `poczta` — poczta transakcyjna; `nadajnik` — nadawanie;
- `mowa`, `tokenizator`, `wiedza`, `podagenci`, `repozytorium`, `zdalne`,
  `zewnetrzne`, `injection`, `narzedzia` — pozostałe obszary dziedziny.

## Droga komendy

Komenda przychodzi kopertą kanału WebSocket i trafia do pakietu `core`.
Brama kontraktu (`core/brama_kontraktu.go`) sprawdza żądanie wobec kontraktu —
pola wymagane i wartości wyliczeń, czytane z wytworu kontraktu — zanim
czynność dotknie domeny; powitanie kanału jest jedynym wyjątkiem. Obsługa
obszaru leży w plikach `handlers_<obszar>.go`, które wpinają komendy rodziny do
adapterów `adapter_modul_<obszar>*.go`; adapter wykonuje czynność na warstwie
danych i zwraca wynik w kształcie zapisanym w kontrakcie. Nazwy komend
i zdarzeń pochodzą wyłącznie z wytworu generatora `budowa/shared/contract.go`.

## Weryfikacja

Rdzeń przechodzi `go build ./...` oraz `go vet ./...` bez zastrzeżeń; obie
miary uruchamia drabina weryfikacji `narzedzia/drabina.sh` wraz z kontrolą
formatu `gofmt`.
