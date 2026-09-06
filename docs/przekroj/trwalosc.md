# Danaco Console — przekrój pionowy: trwałość

Opracowanie opisuje warstwę trwałości w postaci, w jakiej działa w bieżącej
budowie. Obejmuje otwarcie bazy, prowadzenie schematu migracjami oraz kontrolę
spójności kroków.

## Baza

Za trwałość odpowiada pakiet `budowa/server/internal/store`: otwiera plik
SQLite, doprowadza schemat do bieżącej wersji migracjami i kontroluje spójność
bazy. Sterownik jest czystym Go (`modernc.org/sqlite`), bez CGO, przez co rdzeń
buduje się bez zależności od kompilatora C. Dostęp do danych dziedzinowych
ponad warstwą trwałości prowadzi pakiet `budowa/server/internal/dane`.

## Migracje

Schemat prowadzą kroki migracji — pliki `migracja_<numer>_<nazwa>.sql`
w pakiecie `store`, stosowane w kolejności numerów od `migracja_001_fundament.sql`
do `migracja_499_przebieg_powiadomiony.sql`. Krok już zastosowany jest
nietykalny: zmiana jego treści DDL jest odmową startu rdzenia, a rozszerzenie
schematu wchodzi wyłącznie nowym krokiem o kolejnym numerze.

## Kontrola spójności kroków

Spójność kroku mierzy suma SHA-256 z treści znormalizowanej
(`zrodlo_migracji.go`): normalizacja usuwa — poza literałami znakowymi —
komentarze `--`, końcowe białe znaki wiersza i wiersze puste. Dzięki temu
redakcja komentarza w pliku migracji nie unieważnia kroków już zastosowanych,
a każda zmiana samego schematu zmienia sumę i zatrzymuje start na bazie
istniejącej.

## Odwzorowanie kontraktu

Wyliczenia kontraktu mające odpowiednik w schemacie niosą pole `kolumnaBazy`
w zapisie `tabela.kolumna` oraz pole `baza` przy każdej wartości — związek
kontraktu z modelem danych opisuje opracowanie [kontrakt](kontrakt.md).
