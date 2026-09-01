-- Licznik prób drogi potwierdzenia oraz jedna postać wskazania konta w bramce.
--
-- Kod z listu ma sześć cyfr, czyli milion wartości, i jest ważny godzinę.
-- Zamknięcie drogi po pierwszym użyciu nie dotyka prób nieudanych, bo kod
-- nietrafiony nie wskazuje żadnego wiersza — bez licznika cały ten zakres
-- przechodzi się w całości. Kolumna zlicza pomyłki popełnione, gdy droga stała
-- czynna, i pozwala ją zamknąć, zanim zgadujący dojdzie do jej wartości.
ALTER TABLE potwierdzenie_tozsamosci ADD COLUMN proby INTEGER NOT NULL DEFAULT 0;

-- Pomyłka dolicza się wszystkim drogom czynnym jednego celu, więc wybór po
-- parze (cel, użyte) idzie przy każdej próbie wpisania kodu.
CREATE INDEX idx_potwierdzenie_cel_uzyte ON potwierdzenie_tozsamosci (cel, uzyte);

-- Wskazanie konta miało dwie postacie znaczące to samo: NULL wiersza zastanego
-- sprzed migracji 406 oraz zero wpisywane, gdy rdzeń konta nie znał. Zapytania
-- bramki czytają NULL jako konto najstarsze, a zera nie czytają jako niczego,
-- więc wiersz z zerem jest niewidoczny nawet dla własnego właściciela: jego
-- urządzenie nie stoi w wykazie, a metody nie da się zdjąć.
UPDATE metoda_uwierzytelnienia  SET konto_id = NULL WHERE konto_id = 0;
UPDATE sesja_bramki             SET konto_id = NULL WHERE konto_id = 0;
UPDATE potwierdzenie_tozsamosci SET konto_id = NULL WHERE konto_id = 0;
