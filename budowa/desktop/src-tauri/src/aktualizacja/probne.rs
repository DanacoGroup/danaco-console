//! Katalog próbny dla sprawdzianów aktualizacji — pliki na czas jednego testu.
//!
//! Testy `droga.rs` pracują na prawdziwych plikach: tworzą je, podmieniają
//! przez `rename` i sprawdzają, co zostało na dysku. Ten pomocnik daje każdemu
//! testowi własny katalog w katalogu tymczasowym systemu i sprząta go, gdy
//! wartość wychodzi z zakresu — także wtedy, gdy test padł, bo `Drop` biegnie
//! również przy odwijaniu paniki.
//!
//! Katalogi dwóch testów nie mogą na siebie wpadać: testy biegną równolegle
//! w jednym procesie, a te same testy potrafią biec jednocześnie w dwóch
//! drzewach roboczych. Stąd w ścieżce stoi i nazwa podana przez test
//! (rozdziela testy jednego procesu), i numer procesu (rozdziela procesy).

use std::path::PathBuf;

/// Katalog na pliki jednego testu. Powstaje pusty, znika przy `Drop`.
pub struct KatalogProbny {
    /// Ścieżka założonego katalogu.
    sciezka: PathBuf,
}

impl KatalogProbny {
    /// Zakłada świeży katalog próbny o podanej nazwie.
    ///
    /// Katalog pozostały po poprzednim, przerwanym biegu zostaje najpierw
    /// usunięty — test ma zaczynać od stanu pustego, nie od cudzych resztek.
    /// Niepowodzenie jest tu paniką, nie odmową: bez katalogu nie ma czego
    /// testować, a panika wskazuje wprost, który test nie miał gdzie pracować.
    pub fn nowy(nazwa: &str) -> Self {
        let sciezka = std::env::temp_dir().join(format!(
            "danaco-katalog-probny-{}-{nazwa}",
            std::process::id()
        ));
        if sciezka.exists() {
            std::fs::remove_dir_all(&sciezka).unwrap_or_else(|blad| {
                panic!(
                    "nie udało się usunąć pozostałości katalogu próbnego {}: {blad}",
                    sciezka.display()
                )
            });
        }
        std::fs::create_dir_all(&sciezka).unwrap_or_else(|blad| {
            panic!(
                "nie udało się założyć katalogu próbnego {}: {blad}",
                sciezka.display()
            )
        });
        Self { sciezka }
    }

    /// Ścieżka pliku o podanej nazwie wewnątrz katalogu — bez tworzenia pliku.
    ///
    /// Do testów, które pytają o plik nieistniejący.
    pub fn plik(&self, nazwa: &str) -> PathBuf {
        self.sciezka.join(nazwa)
    }

    /// Tworzy w katalogu plik o podanej treści i oddaje jego ścieżkę.
    ///
    /// Treść pusta jest treścią jak każda inna — testy dolnej granicy
    /// sprawdzają właśnie plik zerowej długości.
    pub fn plik_z_trescia(&self, nazwa: &str, tresc: &[u8]) -> PathBuf {
        let sciezka = self.plik(nazwa);
        std::fs::write(&sciezka, tresc).unwrap_or_else(|blad| {
            panic!(
                "nie udało się zapisać pliku próbnego {}: {blad}",
                sciezka.display()
            )
        });
        sciezka
    }
}

impl Drop for KatalogProbny {
    /// Sprząta katalog wraz z zawartością. Niepowodzenie sprzątania jest
    /// przemilczane celowo: panika w `Drop` podczas odwijania paniki testu
    /// przerywa cały proces testowy i zasłania właściwą przyczynę.
    fn drop(&mut self) {
        let _ = std::fs::remove_dir_all(&self.sciezka);
    }
}
