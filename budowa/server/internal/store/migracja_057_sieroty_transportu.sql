-- Migracja usuwa dwie tabele-sieroty, do których nikt nie pisze i z których nikt nie
-- czyta: połączenie i proces sesji.

DROP TABLE polaczenie;
DROP TABLE proces_sesji;
