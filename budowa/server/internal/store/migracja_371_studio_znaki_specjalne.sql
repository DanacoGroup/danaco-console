-- Migracja 371 dodaje wykaz znaków ostatnio użytych oraz prawnicze
-- i ułamkowe uzupełnienie zasad autozamiany.

CREATE TABLE znak_ostatnio_uzyty_studio (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    kod       TEXT    NOT NULL UNIQUE,
    znak      TEXT    NOT NULL,
    ile_uzyc  INTEGER NOT NULL DEFAULT 1,
    uzyto     TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_znak_ostatnio_uzyty_studio_kolejnosc
    ON znak_ostatnio_uzyty_studio(uzyto DESC, ile_uzyc DESC);

-- Uzupełnienie wykazu zasad fabrycznych dokłada znaki prawnicze i ułamki
-- pominięte przy zasadach zakładanych wcześniej.
INSERT OR IGNORE INTO autozamiana_znaku_studio (skrot, zamiennik, fabryczna) VALUES
    ('(par)', '§',  1),
    ('(nr)',  '№',  1),
    ('(st)',  '°',  1),
    ('<->',   '↔',  1),
    ('=>',    '⇒',  1),
    ('1/2',   '½',  1),
    ('1/4',   '¼',  1),
    ('3/4',   '¾',  1);
