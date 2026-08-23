-- Migracja 116 — skrzynka Operatora: protokół, źródło nastaw i domyślność.
--
-- Aplikacja nie ma własnego serwera poczty — używa skrzynki, którą Operator już
-- ma skonfigurowaną na urządzeniu albo w chmurze. Ta migracja jest zapisem tego
-- w schemacie: nie zakłada ani jednej tabeli opisującej pocztę (kolejki, aliasy,
-- domeny, konta), a jedynie dopowiada, jak rdzeń łączy się z cudzą skrzynką.
--
-- ── DLACZEGO NIE MA TU NOWEJ TABELI SKRZYNEK ────────────────────────────────
-- Tabela `skrzynka_pocztowa` opisuje dokładnie ten sam byt: adres, host i port
-- odbioru, host i port wysyłki, użytkownik, odwołanie do sejfu, weryfikacja TLS.
-- Druga tabela na skrzynkę Operatora byłaby drugą prawdą o tym, do czego rdzeń
-- się loguje: rodzina `mail.*` widziałaby jedne skrzynki, wyzwalacz automatyk
-- rodzaju `mail` drugie, a Operator, który podpiął skrzynkę komendą, nie mógłby
-- na nią zbudować automatyki. Dokładamy więc sześć kolumn do bytu, który już
-- istnieje.
--
-- ── PROTOKÓŁ STAJE SIĘ POLEM, BO PRZESTAŁ BYĆ STAŁĄ ─────────────────────────
-- Odbiór może iść IMAP-em, JMAP-em albo POP3, więc protokół musi stać w wierszu,
-- a nie być stałą. Domyślną wartością jest `imap`, bo tym mówi rdzeń i tego chce
-- kontrakt („brak bierze imap"); zależność `github.com/emersion/go-imap/v2`
-- niesie odbiór IMAP (uzasadnienie w nagłówku pakietu `internal/poczta`).
--
-- ── ŹRÓDŁO NASTAW JEST FAKTEM O POCHODZENIU, NIE OZDOBĄ ─────────────────────
-- `zrodlo` mówi, skąd rdzeń wziął te nastawy: `operator` (wpisał je sam),
-- `urzadzenie` (odczytane z klienta poczty na maszynie — `mail.account.discover`)
-- albo `chmura`. Bez tej kolumny podpowiedź odczytana z Thunderbirda wyglądałaby
-- po zapisie identycznie jak nastawy wpisane ręcznie, a to są dwie różne rzeczy:
-- pierwsza może być nieaktualna względem tego, co Operator zmienił w kliencie.
-- Wartości są te same, co w wyliczeniu `MailAccountSource` kontraktu.
--
-- ── SZYFROWANIE JEST DWIEMA KOLUMNAMI, BO JEST DWIEMA DECYZJAMI ─────────────
-- Odbiór i wysyłka jadą osobnymi gniazdami do osobnych serwerów: skrzynka
-- z IMAP-em po TLS-ie na 993 i wysyłką przez STARTTLS na 587 jest układem
-- typowym, a nie wyjątkiem. Jedna wspólna kolumna zmuszałaby do zgadywania
-- trybu z numeru portu — i myliłaby się przy każdej skrzynce postawionej na
-- porcie nietypowym. `tls_weryfikacja` zostaje i dotyczy obu:
-- „czy szyfrować" i „czy sprawdzać łańcuch certyfikatu" to różne pytania.
--
-- ── DOMYŚLNOŚĆ ROZSTRZYGA ZA OPERATORA, GDY NIE WSKAZAŁ ─────────────────────
-- Wszystkie komendy rodziny niosą `accountId` jako pole opcjonalne („brak
-- bierze domyślną”). Bez tej kolumny „domyślna" znaczyłoby „pierwsza po
-- kolejności klucza" — czyli ta, którą Operator podpiął najdawniej, co nie ma
-- żadnego związku z tym, której używa. Indeks częściowy pilnuje, że domyślna
-- jest co najwyżej jedna: dwie domyślne skrzynki to pytanie bez odpowiedzi.

ALTER TABLE skrzynka_pocztowa ADD COLUMN protokol TEXT NOT NULL DEFAULT 'imap'
    CHECK(protokol IN ('imap','jmap','pop3'));

ALTER TABLE skrzynka_pocztowa ADD COLUMN zrodlo TEXT NOT NULL DEFAULT 'operator'
    CHECK(zrodlo IN ('operator','urzadzenie','chmura'));

-- Nazwa wyświetlana nadawcy. NULL znaczy „Operator jej nie podał" — list
-- wyjdzie wtedy z samym adresem w polu From, a nie z nazwą zmyśloną z adresu.
ALTER TABLE skrzynka_pocztowa ADD COLUMN nazwa_wyswietlana TEXT;

ALTER TABLE skrzynka_pocztowa ADD COLUMN szyfruj_odbior INTEGER NOT NULL DEFAULT 1
    CHECK(szyfruj_odbior IN (0,1));

ALTER TABLE skrzynka_pocztowa ADD COLUMN szyfruj_wysylke INTEGER NOT NULL DEFAULT 0
    CHECK(szyfruj_wysylke IN (0,1));

ALTER TABLE skrzynka_pocztowa ADD COLUMN domyslna INTEGER NOT NULL DEFAULT 0
    CHECK(domyslna IN (0,1));

CREATE UNIQUE INDEX idx_skrzynka_pocztowa_domyslna
    ON skrzynka_pocztowa(domyslna) WHERE domyslna = 1;
