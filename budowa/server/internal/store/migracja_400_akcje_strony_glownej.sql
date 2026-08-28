-- Migracja 400 — akcje Strefy 2 i Strefy 3 Strony głównej w katalogu akcji.
--
-- Zaczyn katalogu (migracja 009) wziął komendy nawigacji platformy, sesji,
-- okien, konfiguracji i kanałów oraz pasek promptu każdego modułu. Poza nim
-- zostały dwie strefy Strony głównej, które opracowanie interfejsu użytkownika
-- wymienia POZYCJA PO POZYCJI, wraz z etykietą działania każdej z nich:
--
--   Strefa 2 — kafle komponentów własnych: Automations („Zbuduj automatykę”),
--   Agents („Skonfiguruj agenta”), Workspace („Załóż projekt”), Assistant
--   („Ustaw profil asystenta”), panel wysuwany z listą komponentów zapisanych
--   oraz rozwinięcie `Operacje ▾` z usunięciem komponentu.
--
--   Strefa 3 — listwa ustawień: okno konfiguracji, Mobile i Always On Display.
--
-- Okno konfiguracji ma już pozycję w katalogu (`config.get` z zaczynu), więc
-- z listwy dochodzą dwie funkcje globalne. Nazwy kolumny `nazwa` są etykietami
-- z opracowania, nie tłumaczeniem nazw komend — kafel ma w panelu mówić to samo,
-- co mówi na Stronie głównej.
--
-- ── Czego tu NIE MA i dlaczego ──────────────────────────────────────────────
-- Rozwinięcie `Operacje ▾` kafla wymienia cztery czynności: duplikowanie,
-- eksport, import i usunięcie. Kontrakt ma komendę wyłącznie na usunięcie
-- (`component.delete`), więc wchodzi tylko ono. Pozycja katalogu wskazująca
-- komendę, której w rejestrze rdzenia nie ma, byłaby atrapą w oknie — kontrola
-- stoi w `TestKatalogAkcjiWskazujeKomendyKontraktu`.
--
-- `component.assign` zostaje poza katalogiem świadomie: kontrakt zapisuje przy
-- tej komendzie, że znaczenie przypisania nie zostało rozstrzygnięte. Kontrolka
-- czynności o nierozstrzygniętym skutku jest gorsza od jej braku, bo Operator
-- naciska ją, nie wiedząc, co się stanie. Wiersz dojdzie razem z rozstrzygnięciem.
--
-- ── Kod pozycji ─────────────────────────────────────────────────────────────
-- Cztery kafle wołają JEDNĄ komendę (`component.create`) z różnym rodzajem
-- komponentu, a kolumna `kod` ma warunek UNIQUE. Kod pozycji jest więc nazwą
-- komendy złączoną z wartością `ComponentKind` z kontraktu — pełnymi nazwami
-- obu bytów, bez ani jednego znaku wymyślonej numeracji.
--
-- Warunek dostępności zostaje pusty we wszystkich wierszach: żadna z tych
-- komend nie wymaga w żądaniu środowiska, sesji, okna, kolejki ani kanału.
-- Kolejność podejmuje wykaz zaczynu (poziom globalny kończył się na 13).

WITH katalog(kod, nazwa, opis, ikona, komenda, kolejnosc) AS (
    VALUES
        -- Strefa 2 — kafle komponentów własnych, w kolejności z opracowania.
        ('component.create.automations', 'Zbuduj automatykę',
         'Zakłada komponent własny rodzaju automations — automatykę modułu Automations',
         'automatyzacja',  'component.create', 14),
        ('component.create.agents',      'Skonfiguruj agenta',
         'Zakłada komponent własny rodzaju agents — eksperta modułu Agents',
         'agent',          'component.create', 15),
        ('component.create.workspace',   'Załóż projekt',
         'Zakłada komponent własny rodzaju workspace — projekt przestrzeni roboczej',
         'folder',         'component.create', 16),
        ('component.create.assistant',   'Ustaw profil asystenta',
         'Zakłada komponent własny rodzaju assistant — profil asystenta',
         'uzytkownik',     'component.create', 17),
        ('component.list',               'Komponenty własne',
         'Zwraca komponenty własne zapisane w Strefie 2 Strony głównej',
         'tabela-danych',  'component.list',   18),
        ('component.update',             'Zmień komponent własny',
         'Zmienia komponent własny; pola pominięte zostają bez zmian',
         'olowek',         'component.update', 19),
        ('component.delete',             'Usuń komponent własny',
         'Usuwa komponent własny',
         'kosz',           'component.delete', 20),

        -- Strefa 3 — dwie funkcje globalne listwy ustawień.
        ('mobile.process.list',          'Mobile',
         'Procesy platformy widoczne z urządzenia mobilnego',
         'telefon',        'mobile.process.list', 21),
        ('aod.status.get',               'Always On Display',
         'Stan nakładki globalnego agenta towarzyszącego',
         'aktywnosc',      'aod.status.get',      22)
)
INSERT INTO akcja (kod, nazwa, opis, ikona, poziom_zasiegu_id, klucz_zasiegu,
                   komenda, warunek_dostepnosci, kolejnosc, aktywna)
SELECT katalog.kod, katalog.nazwa, katalog.opis, katalog.ikona, p.id, '',
       katalog.komenda, '', katalog.kolejnosc, 1
  FROM katalog
  JOIN poziom_zasiegu p ON p.kod = 'globalny'
 WHERE true
ON CONFLICT(kod) DO NOTHING;
