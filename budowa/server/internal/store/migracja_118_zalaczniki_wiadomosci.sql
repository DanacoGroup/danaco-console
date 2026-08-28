-- Dodaje do tabeli wiadomosc kolumnę zalaczniki, przechowującą tablicę JSON odwołań do plików, z wartością pustą oznaczającą brak wiedzy o załącznikach sprzed migracji.

ALTER TABLE wiadomosc ADD COLUMN zalaczniki TEXT;
