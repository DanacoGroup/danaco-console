*Uzasadnienia komentarzy i decyzji zaplecza, przeniesione z komentarzy kodu przy redakcji; rozdziały nazwane ścieżkami plików.*

# Uzasadnienia — zaplecze

## budowa/server/internal/store/migracja_012_katalog_ustawien.sql

Oś zasięgu `ConfigAxis` (`platform` · `model` · `account`) jest prostopadła do
ośmiu poziomów zasięgu: poziom mówi jak wąsko (okno → … → globalny), oś mówi
dla czego (platforma, model, konto). Klucz rozstrzygania jest złożony: klucz +
poziom + byt poziomu + oś + byt osi, dlatego `ustawienie` niesie kolumny `os`
i `klucz_osi`, a więz jednoznaczności obejmuje obie.

Kolumna `os_zasiegu.pierwszenstwo` porządkuje osie od najszerszej (platforma)
do najwęższej (konto): konto jest bytem konkretnym, model klasą, platforma
tłem. Rozstrzyganie idzie najpierw po poziomie, a dopiero w ramach poziomu po
osi: konto → model → platforma. Ustawienie per konto zapisane globalnie nie
bije więc ustawienia zapisanego na oknie — oś opisuje adresata wartości, nie
jej wagę.

Tabela `ustawienie` nie jest wskazywana kluczem obcym z żadnej innej tabeli,
więc przebudowa (nowa tabela → przepisanie wierszy → podmiana nazwy) nie rusza
niczyich odwołań. Indeks `idx_ustawienie_klucz` wraca pod tą samą nazwą.

Pozycje katalogu ustawień odpowiadają bytom czytanym przez rdzeń:

1. `server/internal/konfig/definicje_wykonania.go` — osiem parametrów
   wykonania okna komunikacji wraz z wartościami domyślnymi pochodzącymi ze
   stałych kontraktu (PermissionMode.manual, WindowRole.standalone,
   ExecutionEnv.local).
2. `server/internal/konfig/definicje_izolacji.go` — jedenaście punktów
   izolacji: trzy wymiary kontekstu (`odrebna`) i osiem zakresów technicznych
   (`wylaczony`).
3. Klucze katalogu roboczego i tożsamości: `katalog.roboczy.podstawa`,
   `katalog.roboczy.wzorzec_sesji`, `tozsamosc.tryb_domyslny`.
4. `server/internal/injection/ustawienia.go` — pola struktury Ustawienia:
   Program (`harness.program_claude`), PlikUstawien (`harness.plik_ustawien`,
   przełącznik --settings), KonfiguracjaMCP (`harness.konfiguracja_mcp`,
   przełącznik --mcp-config).
5. `server/internal/session/obieg.go` — ProgBrakuPostepuDomyslny
   (`petla.prog_braku_postepu`, pole LoopState.threshold kontraktu).
6. Egzekwowanie uwierzytelniania włączane przez Operatora
   (`bezpieczenstwo.egzekwowanie_uwierzytelniania`).
7. `client/src/motyw/motyw.ts` — wybór motywu `light` albo `dark`; brak
   wyboru znaczy preferencję systemu (`personalizacja.motyw`).

Poza katalogiem pozostają parametry startu procesu — port, katalog danych,
katalog klienta, katalog profili, rola procesu. Czyta je
`server/internal/konfiguracja` w chwili startu z warstw: wartość domyślna →
zmienna środowiska → argument wywołania (`.env.example`), a plik bazy jest
dopiero skutkiem tych parametrów. Wiersz katalogu byłby dla nich drugim
źródłem prawdy, którego nikt nie czyta.

Każde wstawienie kończy się klauzulą ON CONFLICT DO NOTHING, więc migracja
przechodzi także na bazie, w której część wierszy już jest. Klauzula
`WHERE true` przed ON CONFLICT jest wymogiem składni SQLite dla
INSERT ... SELECT z upsertem.

`wartosc_domyslna` trzymana jest tak samo jak `ustawienie.wartosc`: napisem
w postaci właściwej dla rodzaju. Brak wiersza w tabeli `ustawienie` znaczy
właśnie tę wartość, a brak wiersza w katalogu nie jest awarią — rezolwer
schodzi wtedy na rejestr wbudowany rdzenia.

## budowa/server/internal/store/migracja_016_stan_kolejki_wyczerpana.sql

Kontrakt zna pięć stanów kolejki: QueueStatus.idle, running, paused, stopped,
done. Mapa przekładu `WartosciBazyQueueStatus` wiąże je kolejno z kolumnami
'bezczynna', 'pracuje', 'wstrzymana', 'zatrzymana', 'wyczerpana'. Więz CHECK
wprowadzony w `migracja_003_kolejki.sql` dopuszczał tylko cztery pierwsze, więc
zapis stanu 'wyczerpana' kończył się odmową schematu, mimo że kontrakt ten
stan zna. Kolejka, w której nie została ani jedna pozycja czynna, musi mieć
jak się nazwać, inaczej silnik wykonania nie ma stanu końcowego dla kolejki
i po ostatniej pozycji błędnie zgłasza dalszą pracę.

SQLite nie udostępnia polecenia zdejmującego więz CHECK z istniejącej kolumny,
więc tabela `kolejka` przechodzi przebudowę: nowa tabela z poprawionym więzem,
przepisanie wierszy, podmiana nazwy — tak samo jak `ustawienie` w
`migracja_012_katalog_ustawien.sql`. Odmienność wynika z tego, że `kolejka`
jest wskazywana kluczem obcym przez `pozycja_kolejki` i przez
`log_akcji_kolejki`, w obu przypadkach z ON DELETE CASCADE. Przy włączonym
więzie kluczy obcych (`store/baza.go`, pragma `foreign_keys(1)` w DSN
każdego połączenia) polecenie DROP TABLE na tabeli rodzica wykonuje niejawne
DELETE wszystkich jej wierszy, co uruchamia kaskadę i kasuje wszystkie
pozycje kolejek oraz cały dziennik akcji. Pragmy nie da się wyłączyć na czas
kroku, ponieważ migracja biegnie w transakcji (`store/migracje.go`), a
`PRAGMA foreign_keys` w transakcji jest bez skutku. Z tego powodu przebudowa
obejmuje wszystkie trzy tabele obszaru w porządku bezpiecznym dla danych:
kopie z poprawionym więzem i przepisaniem wierszy co do kolumny; skasowanie
tabel starych od dziecka do rodzica, tak że w chwili kasowania `kolejka`
żaden wiersz jej już nie wskazuje i kaskada nie ma czego zabrać; podmiana
nazw poleceniem ALTER TABLE ... RENAME, które przepisuje odwołania kluczy
obcych w tabelach pozostałych; odtworzenie indeksów pod nazwami z
`migracja_003_kolejki.sql`, wolnymi po skasowaniu tabeli, do której należały.

Słownik stanów pozycji kolejki pozostaje bez zmiany: siedem wartości z
`migracja_003_kolejki.sql` wystarcza silnikowi wykonania na pełny cykl życia
zlecenia. Migracja naprawia więz stanu kolejki, nie rozszerza model danych.

## budowa/server/internal/store/migracja_192_roundtable_graf.sql

Graf argumentów jest trwały, a nie liczony przy każdym odczycie, ponieważ węzeł
da się oznaczyć jako kluczowy (`roundtable.argument.pin`). Oznaczenie postawione
na węźle wyliczanym w locie znikałoby przy następnym odczycie, bo węzeł
dostawałby wtedy nowy identyfikator.

Katalog błędów logicznych jest dwuwarstwowy. Definicje wnosi migracja i są
wspólne dla całej platformy — nie należą do okna. Zakres wykrywania
(`roundtable.fallacy.catalog.set`) należy do okna, więc leży w osobnej tabeli
wiążącej kod błędu z oknem; bez tego rozdziału wyłączenie błędu w jednym oknie
wyłączałoby go wszystkim.

## budowa/server/internal/store/migracja_051_aplikacje.sql

Warsztat modułu Apps trzyma treść pliku warstwy frontendu i backendu produktu,
a nie stan projektu. Kontrakt `apps.workspace.update` niesie `Layer`
(frontend/backend), `Path`, `Content` i opcjonalny `ComponentId` — kształt
zgodny z zapisem pliku edytora w module Developer (`developer_wersja_pliku`),
odmienny od tabeli `projekt`, która trzyma nazwę, opis i stan, bez treści
pliku. Współdzielenie jednej tabeli `projekt` obciążałoby każdy jej wiersz
kolumnami treści pliku, których większość wierszy nigdy nie użyje, dlatego
warsztat Apps ma własną tabelę budowaną na wzór `developer_wersja_pliku`.
Commity kodu produktu idą wspólną komendą `developer.git.action`, poza
zasięgiem tej migracji.

Plik warsztatu trzyma jeden wiersz na parę (okno, warstwa, ścieżka), nie
historię wersji: kontrakt `AppsWorkspaceUpdateRequest` nie niesie odpowiednika
`createVersion` z modułu Developer, więc `apps.workspace.update` jest zwykłym
nadpisaniem stanu bieżącego, a tabela odzwierciedla to wprost przez UPSERT po
kluczu (okno, warstwa, ścieżka), bez osobnej tabeli-dziennika.

Architektura trzyma komponenty jako tabelę własną, nie jako zapis w kolumnie
JSON, ponieważ `AppComponent` niesie pola, po których trzeba filtrować
i wiązać przy odczycie: `Kind` wchodzi do walidacji układu, a `DependsOn`
tworzy graf zależności; parsowanie JSON przy każdym odczycie byłoby kosztem
walidacji układu (`AppArchitecture.ValidationIssues`) i zapytań o zależność.
`apps.architecture.define` nadsyła całą listę komponentów na nowo (kontrakt:
`Components []AppComponent`, bez trybu częściowej zmiany), więc zapis jest
zawsze usunięciem komponentów architektury i wstawieniem przysłanych od nowa.
Zależność między komponentami ma z tego samego powodu własną tabelę
złącznikową: `AppComponent.DependsOn []string` jest listą identyfikatorów
zewnętrznych innych komponentów tej samej architektury.

Dziennik wdrożeń powtarza wzór `developer_budowanie`, ponieważ
`apps.deployment.run` i `developer.build.run` mają identyczny kształt zadania:
zlecenie z zewnętrznym kodem, oknem, czasem startu i końca, stanem i śladem
tekstowym. Rdzeń niczego nie wdraża naprawdę — tabela niesie wyłącznie ślad
zlecenia (środowisko, strategia, wersja, notatki, adres po wdrożeniu,
odnośnik do logu) i jego stan, a sam przebieg prowadzi rdzeń w pamięci.
`RollbackToDeploymentId` i `RolledBackFromId` są parą pól tego samego łuku,
więc jedna kolumna samoodwołania wystarcza, bo odczyt idzie zawsze od strony
cofnięcia. Log wdrożenia jest odwołaniem, nie treścią w bazie: kontrakt niesie
`LogRef *string`, więc kolumna `log_odwolanie` przechowuje ten odnośnik wprost.

## budowa/server/internal/store/migracja_013_punkty_dostepu.sql

Punkt dostępu nie jest środowiskiem: kolumna środowisko pozostaje profilem
widoczności modułów w bocznej nawigacji, więc dopisanie maszyny do tamtej
tabeli rozbiłoby nawigację platformy, a niniejsza migracja tamtej tabeli nie
dotyka. Punkt dostępu nie jest też katalogiem roboczym modelu: katalog
roboczy to ustawienie kluczy katalog.roboczy.podstawa i
katalog.roboczy.wzorzec_sesji, mówiące, gdzie model zostawia własne pliki,
podczas gdy punkt dostępu mówi, do czego model ma wgląd — model może czytać
jeden katalog, a pliki zostawiać w zupełnie innym miejscu instalacji,
ponieważ są to dwa niezależne ustawienia.

Odwzorowany most MCP nosi nazwę mcp-danaco-pulpit-console. Skrypt
uruchamiający bierze tryb z argumentu albo ze zmiennej środowiskowej
DANACO_MOST_TRYB; brak obu oznacza tryb odczytu. Proces mostu ogranicza
działanie do korzeni podanych zmienną środowiskową rozdzielonych
dwukropkiem; poza te korzenie most nie wychodzi, a w trybie odczytu
narzędzia zapisu nie są ogłaszane w wykazie narzędzi.

Kolumny rodzaj, tryb_domyslny, tryb i stan niosą wartości kontraktu wprost:
mcpBridge, localDirectory, read, write, unknown, reachable, unreachable.
Kontrakt nie deklaruje przy tych wyliczeniach pola baza, w odróżnieniu od
wyliczenia rodzaju konta, które to pole ma. Własny przekład polsko-angielski
w warstwie trwałości byłby drugim źródłem przekładu obok pliku kontraktu.

Słowa odczyt i zapis są słowami argumentu uruchamiającego dany most, nie
wartościami wyliczenia platformy. Są własnością konkretnego mostu, więc leżą
w danych tabeli argument_trybu_mostu, nie w kodzie: inny most może żądać
innych słów, a to oznacza nowy wiersz danych, nie nową gałąź kodu.

Kolumna identyfikator_zewnetrzny nadania dostępu niesie identyfikator
tekstowy rdzenia; wartość pusta oznacza nadanie założone wprost w bazie
danych, bez odpowiednika w pamięci rdzenia.

## budowa/server/internal/store/migracja_002_okna.sql

Wartości wyliczeniowe kolumn tego obszaru pochodzą z kontraktu —
`shared/contract.json` jest jedynym źródłem prawdy. Kontrakt zapisuje nazwy
po angielsku, model danych po polsku; odwzorowanie jest jeden do jednego i
leży wyłącznie w kontrakcie, w polu `baza` przy każdej wartości wyliczenia.
Generator wytwarza z kontraktu słowniki `WartosciBazy*` i `WartosciKontraktu*`
w `contract.go` i `contract.ts`, a warstwa trwałości sięga po te słowniki
zamiast wpisywać przekład u siebie. Odwzorowanie kolumn na typy kontraktu:
`okno_komunikacji.srodowisko_wykonania` na ExecutionEnv,
`okno_komunikacji.tryb_uprawnien` na PermissionMode,
`okno_komunikacji.rola_okna` na WindowRole,
`okno_komunikacji.stan` na WindowStatus,
`proces_sesji.stan` na ProgressStatus,
`wiadomosc.rola` na MessageRole,
`wiadomosc.stan` na MessageStatus,
`wiadomosc.rodzaj_tresci` na ChunkKind.

Rola i persona wiadomości niosą rozróżnienie odrębne: rola mówi, czym jest
nadawca w rozumieniu kontraktu (MessageRole, cztery wartości), persona mówi,
w jakiej funkcji nadawca wystąpił w danej wypowiedzi. Persona jest warstwą
prezentacji i nie mnoży wartości pola rola; pozostaje pusta dla zwykłej
wypowiedzi użytkownika i modelu. Atrybucję wypowiedzi w pętli
koordynator-wykonawca niesie kolumna okno_zrodlowe_id: okno źródłowe zna
własną rola_okna, więc nie jest potrzebny drugi, równoległy słownik ról.
