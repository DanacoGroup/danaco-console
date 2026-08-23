-- ══════════════════════════════════════════════════════════════════════════════
-- Zaczyn kanału lokalnego — świeża instalacja rozmawia bez ręcznej konfiguracji
-- ══════════════════════════════════════════════════════════════════════════════
--
-- Powód: bez zaczynu rejestr kanałów po instalacji jest pusty, więc pierwsze
-- okno nie ma czym rozmawiać. Użytkownik musiałby sam założyć kanał, zanim
-- zobaczy cokolwiek działającego — nie mając skąd wiedzieć, jakich wartości się
-- od niego oczekuje.
--
-- IDENTYFIKATOR MODELU ZOSTAJE PUSTY, I JEST TO WYBÓR, NIE NIEDOPATRZENIE.
-- Nie wiemy, do których modeli ma dostęp to konkretne konto, a wpisanie nazwy
-- „na oko" dałoby kanał, który wygląda na skonfigurowany i wywraca się przy
-- pierwszej turze. Argumenty procesu powstają tak, że przełącznik pojawia się
-- wyłącznie dla ustawienia wypełnionego — pusty model znaczy
-- więc „bez --model", czyli „model domyślny programu `claude`". Wybór modelu
-- pozostaje decyzją Operatora w oknie konfiguracji.
--
-- Rodzaj `cli` i brak wskazania programu dają program domyślny `claude`
-- (adapter_kanal_cli.go). Konto zostaje puste: pula pusta jest stanem
-- dopuszczalnym, a program używa wtedy poświadczeń zastanych w systemie.
--
-- ON CONFLICT DO NOTHING pilnuje, żeby powtórzenie migracji ani instalacja na
-- istniejącej bazie nie nadpisały niczego, co Operator już ustawił.

INSERT INTO kanal_modelu (kod, nazwa, dostawca, identyfikator_modelu, rodzaj_kanalu,
                          parametry_json, multimodalny, aktywny, kolejnosc)
VALUES ('lokalny-claude', 'Kanał lokalny (claude)', 'anthropic', '', 'cli',
        '{}', 0, 1, 0)
ON CONFLICT(kod) DO NOTHING;
