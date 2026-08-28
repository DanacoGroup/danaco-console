-- Indeks treści modułu Library: pełnotekstowe przeszukiwanie zawartości plików
-- repozytorium wiedzy (`library.file.search`).
--
-- Bajty każdego wgrania leżą na dysku pod swoją sumą sha256
-- (`core/adapter_modul_library_magazyn.go`). Bez tego indeksu każde żądanie
-- musiałoby otworzyć i przeczytać wszystkie bloby repozytorium, więc koszt
-- wyszukiwania rósłby z rozmiarem całej biblioteki, a nie z liczbą trafień.
--
-- FTS5, a nie `LIKE` po wyciągu tekstu: budowa niesie SQLite z `ENABLE_FTS5`,
-- więc pełnotekstowy indeks nie dokłada ani zależności, ani procesu.
-- `LIKE '%słowo%'` byłby skanem każdego wiersza przy każdym żądaniu i nie
-- umiałby dopasować słowa niezależnie od wielkości liter i znaków
-- diakrytycznych bez własnego składania form.
--
-- Baza nie przechowuje pliku: nie ma kolumny BLOB, `tresc_odwolanie` pozostaje
-- jedyną drogą do bajtów, a wersje czytają treść wyłącznie z magazynu.
-- W indeksie leży wyciąg tekstowy treści bieżącej, obcięty granicą
-- (`core/adapter_modul_library_indeks.go`) — byt wtórny, odtwarzalny z blobów,
-- po którego utracie biblioteka traci wyłącznie trafność wyszukiwania.
--
-- Rowid indeksu równa się `plik_biblioteki.id`. Tabela FTS5 nie ma kluczy
-- obcych ani kaskady, więc wiązanie idzie przez rowid nadawany wprost przy
-- zapisie: jeden wiersz indeksu na jeden plik, dla treści bieżącej — tej samej,
-- którą oddaje podgląd. Wersje historyczne nie są indeksowane, bo dawałyby
-- trafienia w pliki, w których szukanego słowa już nie ma.

CREATE VIRTUAL TABLE indeks_tresci_biblioteki USING fts5(
    tresc,
    -- unicode61 z usuwaniem znaków diakrytycznych: „ŚRODEK", „środek"
    -- i „srodek" trafiają w to samo słowo. Wariant 2 zdejmuje diakrytyki także
    -- z liter spoza łacińskiego zakresu podstawowego (polskie ą, ę, ó, ż).
    -- „ł" zostaje literą własną — unicode61 nie zna go jako „l ze znakiem",
    -- więc tego jednego złożenia indeks nie sprowadza do „l".
    tokenize = 'unicode61 remove_diacritics 2'
);
