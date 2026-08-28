-- Dodaje do tabeli sesja_bramki kolumnę trwanie, zapisującą długość życia sesji, używaną do wyliczenia chwili wygaśnięcia przy każdym odnowieniu.

ALTER TABLE sesja_bramki ADD COLUMN trwanie INTEGER NOT NULL DEFAULT 0;
