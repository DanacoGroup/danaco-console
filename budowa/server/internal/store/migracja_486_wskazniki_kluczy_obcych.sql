-- Dwadzieścia jeden kolumn klucza obcego stało bez wskaźnika. Zawężenie wiersza
-- dziecka do konta idzie wzorem EXISTS przez klucz obcy do korzenia, a kasowanie
-- kaskadowe przechodzi tą samą drogą — bez wskaźnika każde takie zapytanie czyta
-- tabelę dziecka w całości.
CREATE INDEX IF NOT EXISTS idx_agent_konektor_punkt ON agent_konektor (punkt_dostepu_id);
CREATE INDEX IF NOT EXISTS idx_bieg_orkiestracji_okno_koordynatora ON bieg_orkiestracji (okno_koordynatora_id);
CREATE INDEX IF NOT EXISTS idx_instancja_komponentu_design_kompozycja ON instancja_komponentu_design (kompozycja_id);
CREATE INDEX IF NOT EXISTS idx_kategoria_ustawien_nadrzedna ON kategoria_ustawien (kategoria_nadrzedna_id);
CREATE INDEX IF NOT EXISTS idx_komponent_poziom_zasiegu ON komponent (poziom_zasiegu_id);
CREATE INDEX IF NOT EXISTS idx_nastawa_pracy_studio_dokument ON nastawa_pracy_studio (dokument_id);
CREATE INDEX IF NOT EXISTS idx_orkiestracja_spiecie_okno_roli ON orkiestracja_spiecie_multitasking (okno_roli_id);
CREATE INDEX IF NOT EXISTS idx_powiadomienie_centrum_urzadzenie ON powiadomienie_centrum (urzadzenie_id);
CREATE INDEX IF NOT EXISTS idx_przypisanie_profilu_izolacji_profil ON przypisanie_profilu_izolacji (profil_id);
CREATE INDEX IF NOT EXISTS idx_terminal_host_posredni ON terminal_host (host_posredni_kod);
CREATE INDEX IF NOT EXISTS idx_terminal_tunel_host ON terminal_tunel (host_kod);
CREATE INDEX IF NOT EXISTS idx_wiadomosc_biegu_podagent ON wiadomosc_biegu (podagent_id);
CREATE INDEX IF NOT EXISTS idx_wiadomosc_biegu_pozycja_kolejki ON wiadomosc_biegu (pozycja_kolejki_id);
CREATE INDEX IF NOT EXISTS idx_wpis_pamieci_projektu_poziom_zasiegu ON wpis_pamieci_projektu (poziom_zasiegu_id);
CREATE INDEX IF NOT EXISTS idx_wylaczenie_pamieci_wpis ON wylaczenie_pamieci (wpis_id);
CREATE INDEX IF NOT EXISTS idx_wyzwolenie_automatyki_przebieg ON wyzwolenie_automatyki (przebieg_id);
CREATE INDEX IF NOT EXISTS idx_zlecenie_kolejki_przebieg ON zlecenie_kolejki (przebieg_id);
CREATE INDEX IF NOT EXISTS idx_zlecenie_kolejki_zrodlowe ON zlecenie_kolejki (zlecenie_zrodlowe_id);
CREATE INDEX IF NOT EXISTS idx_zlecenie_pakietu_tlumaczenia_okno ON zlecenie_pakietu_tlumaczenia (okno_id);
CREATE INDEX IF NOT EXISTS idx_zlecenie_przekazania_sesja ON zlecenie_przekazania (sesja_id);
CREATE INDEX IF NOT EXISTS idx_znakowanie_studio_znacznik ON znakowanie_studio (znacznik_nazwa);
