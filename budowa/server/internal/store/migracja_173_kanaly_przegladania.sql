-- Migracja 173 — kanały RSS/Atom/JSON Feed i ich wpisy (`browser.feed.*`).
--
-- Wpis kanału jest bytem osobnym, nie polem kanału: ma własny adres, własny
-- czas publikacji i własne oznaczenie przeczytania, a wykaz `browser.feed.list`
-- pyta o kanały z wpisami albo o same kanały (`includeEntries`). Wpisy zapisane
-- kolumną JSON w wierszu kanału nie dałyby się oznaczyć pojedynczo bez
-- przepisywania całej kolumny przy każdym przeczytanym wpisie.
--
-- Wpis ma więz obcy do kanału z kasowaniem kaskadowym, bo kontrakt mówi wprost:
-- `browser.feed.remove` zdejmuje subskrypcję „wraz z jej wpisami". Kaskada
-- w schemacie jest tu jedyną gwarancją, że zdjęcie kanału nie zostawia wpisów
-- bez rodzica — sprzątanie w kodzie pominęłoby je przy pierwszym błędzie.
--
-- Ten sam adres w tym samym oknie nie zakłada drugiej subskrypcji: warunek
-- UNIQUE(okno, url) czyni z ponownego wywołania `browser.feed.subscribe`
-- odświeżenie zastanego kanału, a nie jego duplikat.
CREATE TABLE kanal_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    okno                     TEXT    NOT NULL,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    postac                   TEXT,
    interwal_sekund          INTEGER,
    pobrano                  TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(okno, url)
);
CREATE INDEX idx_kanal_przegladania_okno ON kanal_przegladania(okno, utworzono DESC, id);

CREATE TABLE wpis_kanalu_przegladania (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    kanal_zewnetrzny_id      TEXT    NOT NULL REFERENCES kanal_przegladania(identyfikator_zewnetrzny) ON DELETE CASCADE,
    url                      TEXT    NOT NULL,
    tytul                    TEXT,
    streszczenie             TEXT,
    przeczytany              INTEGER NOT NULL DEFAULT 0 CHECK(przeczytany IN (0,1)),
    opublikowano             TEXT,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    UNIQUE(kanal_zewnetrzny_id, url)
);
CREATE INDEX idx_wpis_kanalu_przegladania ON wpis_kanalu_przegladania(kanal_zewnetrzny_id, opublikowano DESC, id);
