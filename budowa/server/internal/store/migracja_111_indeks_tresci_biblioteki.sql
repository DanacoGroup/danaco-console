-- Tworzy indeks pełnotekstowy FTS5 nad wyciągiem treści bieżącej plików repozytorium wiedzy, wspierający wyszukiwanie modułu Library bez odczytu blobów magazynu.

CREATE VIRTUAL TABLE indeks_tresci_biblioteki USING fts5(
    tresc,
    -- unicode61 usuwa znaki diakrytyczne, również polskie, lecz traktuje literę ł jako odrębną literę.
    tokenize = 'unicode61 remove_diacritics 2'
);
