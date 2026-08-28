-- Migracja 366 — kopie zapasowe, autozapis i nastawy widoku okna pracy.
--
-- ── Dlaczego kopia zapasowa jest osobna od historii wersji ──────────────────
-- Kopia ma przetrwać awarię procesu I awarię zapisu. Wersja leży w repozytorium
-- sesji i zakłada się ją zapisem, który właśnie się nie udał — więc wersja nie
-- ochroni pracy przed nieudanym zapisem. Kopia jest zakładana niezależnie
-- i dlatego ma własny wiersz.
--
-- ── Dlaczego kolumna `udalo_sie` jest obowiązkowa ───────────────────────────
-- Wskaźnik „zapisano" pokazany, gdy zapis się nie udał, jest najgorszym możliwym
-- błędem tego modułu: Operator zamknie okno i straci pracę. Nieudany zapis musi
-- być WIDOCZNY i NAZWANY — stąd `udalo_sie` i `powod_niepowodzenia` przy każdej
-- kopii, a nie tylko wiersz przy kopiach udanych. Wiersz o nieudanym zapisie
-- jest tu najważniejszym wierszem tej tabeli.
--
-- `zmiany_niezapisane` znaczy kopię niosącą pracę, której w dokumencie nie ma.
-- Po tym Studio samo zgłasza „mam niezapisany dokument z godziny X, przywrócić?",
-- zamiast czekać, aż Operator się domyśli.
--
-- ── Dlaczego autozapis odkłada wersje w osobnym szeregu ─────────────────────
-- Wersje nazwane i kluczowe zakłada Operator. Zapisy samoczynne mają być
-- odróżnialne w wykazie i mieć własną zasadę wygasania — inaczej po godzinie
-- pracy historia wersji przestaje być historią decyzji, a staje się dziennikiem
-- naciśnięć klawisza. Rozróżnienie idzie kolumną `szereg`; pojęcia wersji
-- kluczowej Studio ma już w `studio.version.label.set` i drugiego nie zakłada.
--
-- ── Dlaczego nastawy pracy są jedną tabelą na parę okno-dokument ────────────
-- Skala widoku jest pamiętana PRZY DOKUMENCIE (Właściciel wymienia to wprost),
-- a tryb powierzchni — przy oknie. Jedna tabela z nieobowiązkowym dokumentem
-- obsługuje oba: wiersz bez dokumentu jest nastawą okna, wiersz z dokumentem
-- nastawą tego dokumentu. Dwie tabele znaczyłyby dwa odczyty przy każdym
-- otwarciu okna i pytanie, która wygrywa.

CREATE TABLE kopia_zapasowa_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    identyfikator_zewnetrzny TEXT    NOT NULL UNIQUE,
    dokument_id              INTEGER NOT NULL REFERENCES dokument_studio(id) ON DELETE CASCADE,
    powod                    TEXT    NOT NULL DEFAULT 'manual'
                                     CHECK(powod IN ('interval','event','beforeIrreversible','manual')),
    tresc                    TEXT,
    tresc_odwolanie          TEXT,
    postac_json              TEXT,
    rozmiar_bajtow           INTEGER NOT NULL DEFAULT 0,
    udalo_sie                INTEGER NOT NULL DEFAULT 1,
    powod_niepowodzenia      TEXT,
    zmiany_niezapisane       INTEGER NOT NULL DEFAULT 0,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_kopia_zapasowa_studio_dokument
    ON kopia_zapasowa_studio(dokument_id, utworzono DESC, id DESC);
-- Pytanie „czy jest coś do przywrócenia po nagłym zamknięciu" pada przy każdym
-- otwarciu modułu, więc ma własny indeks.
CREATE INDEX idx_kopia_zapasowa_studio_niezapisane
    ON kopia_zapasowa_studio(zmiany_niezapisane, utworzono DESC);

CREATE TABLE nastawa_pracy_studio (
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    okno                     TEXT    NOT NULL,
    dokument_id              INTEGER REFERENCES dokument_studio(id) ON DELETE CASCADE,
    -- Autozapis
    autozapis_czynny         INTEGER NOT NULL DEFAULT 0,
    autozapis_odstep_sekund  INTEGER NOT NULL DEFAULT 120,
    autozapis_przy_odejsciu  INTEGER NOT NULL DEFAULT 1,
    autozapis_przy_zamknieciu INTEGER NOT NULL DEFAULT 1,
    autozapis_przy_przelaczeniu INTEGER NOT NULL DEFAULT 1,
    kopie_ile_zachowac       INTEGER NOT NULL DEFAULT 20,
    kopie_wygasanie_godzin   INTEGER NOT NULL DEFAULT 168,
    ostatni_zapis            TEXT,
    ostatni_zapis_nieudany   INTEGER NOT NULL DEFAULT 0,
    ostatni_powod_niepowodzenia TEXT,
    -- Widok. Wybór Operatora jest PAMIĘTANY: przełączenie trybu powierzchni
    -- niczego nie gubi i nie zamyka, więc tryb musi przetrwać zamknięcie okna.
    tryb_powierzchni         TEXT    NOT NULL DEFAULT 'tabs'
                                     CHECK(tryb_powierzchni IN ('tabs','split')),
    kierunek_podzialu        TEXT    NOT NULL DEFAULT 'vertical'
                                     CHECK(kierunek_podzialu IN ('vertical','horizontal')),
    granica_podzialu         REAL    NOT NULL DEFAULT 0.5,
    tryb_widoku              TEXT    NOT NULL DEFAULT 'edit'
                                     CHECK(tryb_widoku IN ('edit','printPreview','source',
                                                           'diffOverlay','diffColumns')),
    skala_procent            INTEGER NOT NULL DEFAULT 100,
    skala_nastawa            TEXT    NOT NULL DEFAULT 'hundred'
                                     CHECK(skala_nastawa IN ('custom','pageWidth','wholePage',
                                                             'textWidth','hundred')),
    linijki_widoczne         INTEGER NOT NULL DEFAULT 1,
    linijka_jednostka        TEXT    NOT NULL DEFAULT 'millimeter'
                                     CHECK(linijka_jednostka IN ('millimeter','inch')),
    granica_marginesu        INTEGER NOT NULL DEFAULT 1,
    znaki_formatowania       INTEGER NOT NULL DEFAULT 0,
    stron_w_rzedzie          INTEGER NOT NULL DEFAULT 1,
    widok_rozkladowki        INTEGER NOT NULL DEFAULT 0,
    przewijanie              TEXT    NOT NULL DEFAULT 'continuous'
                                     CHECK(przewijanie IN ('continuous','page')),
    podswietlenie_zmian_modelu INTEGER NOT NULL DEFAULT 0,
    przybornik_widoczny      INTEGER NOT NULL DEFAULT 1,
    utworzono                TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    zaktualizowano           TEXT    NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
-- Jeden wiersz na parę okno-dokument. Wiersz bez dokumentu jest nastawą okna;
-- SQLite traktuje NULL jako różne wartości w UNIQUE, więc wiersz okna pilnuje
-- osobny indeks częściowy.
CREATE UNIQUE INDEX idx_nastawa_pracy_studio_dokument
    ON nastawa_pracy_studio(okno, dokument_id) WHERE dokument_id IS NOT NULL;
CREATE UNIQUE INDEX idx_nastawa_pracy_studio_okno
    ON nastawa_pracy_studio(okno) WHERE dokument_id IS NULL;
