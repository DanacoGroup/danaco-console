# Danaco Console — przekrój pionowy: kanał

Opracowanie opisuje kanał łączący klienta z rdzeniem w postaci, w jakiej działa
w bieżącej budowie. Obejmuje gniazdo, nawiązanie połączenia, bramkę dostępu,
pętle wymiany oraz rozgłaszanie zdarzeń.

## Gniazdo

Kanał jest gniazdem WebSocket obsługiwanym przez pakiet
`budowa/server/internal/transport` biblioteką `github.com/coder/websocket`.
Gniazdo stoi na ścieżce `/ws` (`SciezkaGniazdaDomyslna` w `ustawienia.go`),
a port pochodzi z konfiguracji rdzenia. Treść komunikatów niesie koperta
kontraktu w formacie JSON.

## Limity ramki

Limit odczytu ramki wynosi 16 MiB (`LimitOdczytu` w `nawiazanie.go`), ponieważ
domyślny limit biblioteki — 32 KiB — zrywałby połączenie przy obszernym
ładunku, a zerwanie z powodu długości wiadomości byłoby bramą, której kontrakt
nie przewiduje. Gniazdo, które nie przeszło jeszcze bramki, obowiązuje limit
niższy (`LimitOdczytuPrzedBramka`), bo komendy wejścia mieszczą się w kilku
kilobajtach.

## Nawiązanie i sekret

Nawiązanie połączenia rozstrzyga się przed przyjęciem gniazda
(`nawiazanie.go`). Strona łącząca okazuje sekret nawiązania — wartość losową
jednego uruchomienia powłoki, przekazywaną w zapytaniu adresu gniazda —
albo poświadczenie, które ten sekret zastępuje. Porównanie sekretu przebiega
w czasie stałym (`crypto/subtle`); niezgodność kończy się odpowiedzią HTTP 403
i wpisem do dziennika. Wytworzenie sekretu opisuje opracowanie
[powłoka](powloka.md).

## Bramka

Bramka (`bramka.go`) rozstrzyga, czy żądanie z danego gniazda wolno oddać
rdzeniowi, na podstawie adresu nasłuchu i wskazania Operatora. Bez wskazania
rdzeń nasłuchuje na pętli zwrotnej. Komenda `auth.verify` — potwierdzenie
adresu drogą z listu — przechodzi bramkę zawsze.

## Pętle wymiany i rejestr połączeń

Każde połączenie prowadzi pętlę odbioru (`petla_odbioru.go`) i pętlę wysyłki
(`petla_wysylki.go`); wykonanie komendy spina `wykonanie.go`. Otwarte
połączenia trzyma rejestr (`rejestr_polaczen.go`) wraz z tożsamością strony
(`tozsamosc.go`); rozłączenie sprząta `rozlaczenie.go`.

## Rozgłaszanie zdarzeń

Rozgłos (`rozgloszenie.go`) wysyła kopertę do wszystkich połączeń wskazanego
konta i zwraca liczbę urządzeń, które ją przyjęły; puste wskazanie konta
oznacza wszystkie połączenia rdzenia. Tą drogą rdzeń doręcza zdarzenia
kontraktu do okien klienta.
