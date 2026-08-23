# Plan budowy klienta na przejętym rdzeniu

Wykaz roboczy. Znika wraz z otwarciem budowy właściwej — jego treść przechodzi
wtedy do zleceń wykonawczych.

## Podstawa

Rdzeń wersji poprzedniej jest przejmowany. Podstawą są pomiary, nie założenia:

| Pomiar | Wynik |
|---|---|
| kompilacja rdzenia | bez błędu |
| testy rdzenia | 2003 przechodzą, 4 nie |
| komendy kontraktu odpowiadające | 1005 z 1006 zbadanych |
| testy klienta | 508 przechodzi |
| granica rdzeń–klient | czysta — własny klient trzydziestu wierszy rozmawia z rdzeniem bez zmiany w rdzeniu |

Warstwa widoku poprzedniego klienta nie jest przejmowana. Powstaje nowa.

## Co rdzeń daje klientowi

Dziewięć komend drogi wejścia wraz z wymaganymi polami:

| Komenda | Pola wymagane |
|---|---|
| `connection.hello` | `clientId`, `clientVersion`, `protocolVersion` |
| `auth.register` | `login`, `email`, `password` |
| `auth.verify` | `token` |
| `auth.login` | `method` |
| `environment.enter` | `environmentId`, `clientId` |
| `module.list` | — |
| `session.create` | — |
| `window.create` | `sessionId`, `moduleId`, `modelChannelId`, `workingDirs`, `executionEnv`, `permissionMode`, `windowRole` |
| `message.send` | `windowId`, `content` |

Do tego **75 zdarzeń**, które klient odbiera — od `session.changed` i `stream.chunk`
po zdarzenia właściwe każdemu modułowi.

## Rozjazd nazewnictwa — do świadomego przyjęcia

| Kontrakt rdzenia | Model produktu |
|---|---|
| `Session` — zawiera wiele `Window`, należy do projektu | okno robocze |
| `Window` — własny kanał modelu, katalogi, tryb uprawnień, rola | karta czatu |

Struktura odpowiada modelowi; przesunięte są nazwy. Przemianowanie dotknęłoby
1077 komend i 542 struktur naraz, więc jest osobnym rozstrzygnięciem, nie pracą
przy okazji.

## Fale

Fala kończy się czymś, co daje się uruchomić i zmierzyć. Tereny wewnątrz fali są
rozłączne i biegną równolegle.

### Fala 1 — fundament klienta · **może ruszyć natychmiast**

Nie zależy od prototypów, bo nie dotyczy wyglądu.

| Teren | Przedmiot | Kryterium odbioru |
|---|---|---|
| typy z kontraktu | wytworzenie `contract.ts` z `contract.json`; rozjazd z Go wykrywany poleceniem | wytworzenie powtarzalne, plik bajtowo ten sam przy dwóch przebiegach |
| warstwa połączenia | gniazdo, koperta, korelacja żądanie–odpowiedź, wznowienie po zerwaniu, przeciwciśnienie | zerwanie połączenia w trakcie strumienia i powrót bez utraty zdarzeń — wykazane próbą |
| warstwa stanu | odbiór 75 zdarzeń, jedno źródło stanu, brak stanu w widoku | każde zdarzenie kontraktu ma obsługę albo jawne pominięcie z powodem |
| warstwa zakresu | kaskada czterech kondygnacji, dwie osie, dziedziczenie | rozstrzygnięcie zakresu dla karty zgodne z pozycją 6 rejestru — wykazane tabelą przypadków |

### Fala 2 — droga wejścia · czeka na prototypy

| Teren | Przedmiot |
|---|---|
| instalacja | powłoka Tauri, warianty cienki i pełny, wybór składników zależny od wariantu |
| uruchomienie i uwierzytelnienie | łączenie z rdzeniem, rejestracja, logowanie, odzyskiwanie, przygotowanie środowiska |

Pozycja otwarta rejestru decyzji — pierwsze uruchomienie bez poczty — wiąże ten
teren i musi być rozstrzygnięta przed jego otwarciem.

### Fala 3 — rama aplikacji · czeka na prototypy

| Teren | Przedmiot |
|---|---|
| lewa strona | zarząd sesji i projektów, jedyne miejsce z całością dorobku, zwijanie |
| prawa strona | samouczek jako zawartość domyślna, własne pasmo kart, zwijanie |
| obszar roboczy | przełącznik okien roboczych, przejmowanie uwolnionej przestrzeni |
| szyna szybkiego dostępu | druga droga do modułu, mieszczenie stref, menu nadmiaru |

### Fala 4 — okno robocze · czeka na prototypy

| Teren | Przedmiot |
|---|---|
| pasmo kart | karty jako funkcje, grupowanie, wspólna ramka pary |
| wstążka narzędziowa | kontrakt wstążki wraz z trzema strefami i menu nadmiaru |
| widok dzielony | zestawienie dwóch kart, rozłączenie, wiele par |
| okno boczne | przeniesienie karty i powrót |
| trwałość układu | układ paneli jako stan sesji, odtworzenie po restarcie |

Układ paneli i karty inne niż czat **nie mają struktur w kontrakcie** — trzeba je
dobudować po stronie rdzenia. To jedyne miejsce, w którym fala klienta wymaga
zmiany w rdzeniu.

### Fala 5 — moduł Studio · czeka na prototypy

| Teren | Przedmiot |
|---|---|
| karta czatu | rozmowa, strumień, sterowanie uprawnieniami i dostępem przy rozmowie |
| karta edytora | praca na treści pliku |
| przybornik przeglądarki | inspektor elementów, znakowanie, menu wskazanego elementu |

## Co blokuje co

Fala 1 nie jest blokowana przez nic i powinna ruszyć natychmiast — jej wynik
jest potrzebny każdej następnej.

Fale 2 do 5 czekają na prototypy, bo bez nich nie ma czego zbudować; zlecenie
opisujące wygląd słowami zostało już raz sprawdzone i nie działa.

Fala 4 wymaga dobudowy struktur w rdzeniu — to jedyna zależność zwrotna.
