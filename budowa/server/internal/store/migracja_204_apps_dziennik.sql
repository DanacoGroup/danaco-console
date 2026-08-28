-- Migracja 204 — moduł Apps: wiersze dziennika usług i wdrożeń.
--
-- Jedna tabela na oba dzienniki, bo jest to jeden byt. `apps.service.log.read`
-- i `apps.deployment.log.read` mają identyczny kształt wyniku (`lines`,
-- `total`, `streaming`) i różnią się wyłącznie zawężeniem: pierwsza po
-- komponencie, druga po wdrożeniu. Dwie tabele o tych samych kolumnach byłyby
-- dwiema prawdami o wierszu dziennika.
--
-- Wiersz powstaje z pracy, nie z żądania odczytu. Zapisuje go silnik wykonania
-- wdrożenia (`adapter_modul_aplikacje_wdrozenie_bieg.go`) przy każdym kroku
-- przebiegu oraz serwer podglądu przy podniesieniu i zatrzymaniu. Odczyt
-- niczego nie dopisuje — dziennik pokazujący własne odczyty byłby dziennikiem
-- o sobie.
--
-- `chwila` niesie milisekundy epoki, bo `since` w obu żądaniach jest liczbą
-- milisekund; porównanie z kolumną tekstową wymagałoby przekładu przy każdym
-- odczycie.
--
-- `wdrozenie_kod` i `komponent_kod` są kodami zewnętrznymi, nie więzami obcymi:
-- wiersz dziennika przeżywa byt, którego dotyczy — po to się dziennik prowadzi.

CREATE TABLE wiersz_dziennika_apps (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    okno          TEXT    NOT NULL,
    wdrozenie_kod TEXT,
    komponent_kod TEXT,
    chwila        INTEGER NOT NULL,
    tresc         TEXT    NOT NULL
);
CREATE INDEX idx_wiersz_dziennika_apps_okno ON wiersz_dziennika_apps(okno, chwila, id);
CREATE INDEX idx_wiersz_dziennika_apps_wdrozenie ON wiersz_dziennika_apps(wdrozenie_kod, chwila, id);
