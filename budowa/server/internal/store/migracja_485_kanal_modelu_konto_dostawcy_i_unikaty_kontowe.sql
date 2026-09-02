-- Kolumna `konto_id` znaczy w całej bazie konto Operatora — poza `kanal_modelu`,
-- gdzie od migracji 001 wskazywała konto dostawcy z tabeli `konto`. Jedna nazwa
-- na dwa pojęcia znosi granicę konta w tabeli niosącej odwołania do
-- poświadczeń: warunek konta porównywałby wskazanie dostawcy z kontem
-- Operatora. Wskazanie dostawcy dostaje nazwę własną, a `konto_id` znaczy tu to
-- samo, co wszędzie.
ALTER TABLE kanal_modelu RENAME COLUMN konto_id TO konto_dostawcy_id;
DROP INDEX idx_kanal_modelu_konto;
CREATE INDEX idx_kanal_modelu_konto_dostawcy ON kanal_modelu (konto_dostawcy_id);
ALTER TABLE kanal_modelu ADD COLUMN konto_id INTEGER;

-- Wskaźniki jednoznaczne obejmowały całą tabelę, więc drugie konto nie mogło
-- założyć własnego wiersza pod nazwą zajętą przez pierwsze, a zapis z klauzulą
-- ON CONFLICT sięgał wiersza cudzego. Jednoznaczność obowiązuje w koncie.
-- Wskazanie puste wchodzi przez COALESCE jako zero: w indeksie wartości NULL są
-- względem siebie różne, więc wiersze zastane straciłyby jednoznaczność.
DROP INDEX idx_ikona_design_nazwa;
CREATE UNIQUE INDEX idx_ikona_design_nazwa
    ON ikona_design (okno, nazwa, COALESCE(konto_id, 0));

DROP INDEX idx_konto_domyslne_rodzaj;
CREATE UNIQUE INDEX idx_konto_domyslne_rodzaj
    ON konto (rodzaj, COALESCE(konto_id, 0)) WHERE domyslne = 1;

DROP INDEX idx_profil_asystenta_domyslny;
CREATE UNIQUE INDEX idx_profil_asystenta_domyslny
    ON profil_asystenta (domyslny, COALESCE(konto_id, 0)) WHERE domyslny = 1;

DROP INDEX idx_skrzynka_pocztowa_domyslna;
CREATE UNIQUE INDEX idx_skrzynka_pocztowa_domyslna
    ON skrzynka_pocztowa (domyslna, COALESCE(konto_id, 0)) WHERE domyslna = 1;

DROP INDEX idx_terminal_skrypt_alias;
CREATE UNIQUE INDEX idx_terminal_skrypt_alias
    ON terminal_skrypt (alias, COALESCE(konto_id, 0))
    WHERE alias IS NOT NULL AND alias <> '';

DROP INDEX idx_wyciszenie_nakladki_byt;
CREATE UNIQUE INDEX idx_wyciszenie_nakladki_byt
    ON wyciszenie_nakladki (rodzaj, zakres, klucz_zakresu, klasa_zdarzen,
                            COALESCE(konto_id, 0));

DROP INDEX idx_zasada_przechowywania_zakres;
CREATE UNIQUE INDEX idx_zasada_przechowywania_zakres
    ON zasada_przechowywania (zakres, COALESCE(zakres_kod, ''), COALESCE(konto_id, 0));
