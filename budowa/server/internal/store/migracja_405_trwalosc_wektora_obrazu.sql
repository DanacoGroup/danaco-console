-- Podobieństwo obrazu trzyma wynik osi obrazu policzony raz, żeby kolejne
-- zapytanie o ten sam obraz i to samo zdanie czytało wiersz zamiast wołać
-- pomocnika Pythona od nowa. Model liczy podobieństwo łącznie dla obrazu
-- i zdania w jednym przebiegu (`pomocnik_obrazu.py`), więc wiersz nie niesie
-- wektora obrazu osobno — niesie gotowy wynik tej pary.

CREATE TABLE podobienstwo_obrazu (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Ścieżka pliku obrazu na maszynie rdzenia w chwili liczenia.
    sciezka      TEXT    NOT NULL,
    -- Odcisk pliku (rozmiar i czas modyfikacji) w chwili liczenia; zmiana
    -- obrazu daje inny odcisk i unieważnia wiersz.
    odcisk       TEXT    NOT NULL,
    -- Nazwa modelu osi obrazu (`wiedza_model_obrazu`); wynik innego modelu
    -- opisuje podobieństwo w innej przestrzeni, więc zmiana modelu
    -- unieważnia wiersz — inny model to inny wiersz, nie nadpisanie.
    model        TEXT    NOT NULL,
    -- Zdanie zapytania dosłownie; inne zdanie dla tego samego obrazu to inny
    -- wiersz, nie unieważnienie.
    pytanie      TEXT    NOT NULL,
    -- Kosinus zdania z obrazem w przestrzeni wspólnej modelu.
    podobienstwo REAL    NOT NULL,
    -- Milisekundy epoki, w których wiersz powstał albo ostatnio się zmienił.
    utworzono    INTEGER NOT NULL,

    -- Powtórne liczenie tej samej trójki nadpisuje wiersz zamiast dokładać
    -- drugi komplet.
    UNIQUE (sciezka, model, pytanie)
);
