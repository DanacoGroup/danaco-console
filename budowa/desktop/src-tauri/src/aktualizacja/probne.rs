//! Moduł zakłada dla testu aktualizacji odosobniony katalog tymczasowy, którego
//! ścieżka niesie nazwę testu i numer procesu, i sprząta go przy zwolnieniu wartości.

use std::path::PathBuf;

/// Struktura opisuje katalog tymczasowy jednego testu: powstaje pusty, a implementacja
/// `Drop` usuwa go wraz z zawartością w chwili zwolnienia wartości.
pub struct KatalogProbny {
    /// Ścieżka założonego katalogu.
    sciezka: PathBuf,
}

impl KatalogProbny {
    /// Zakłada świeży katalog próbny o podanej nazwie, usuwając pozostałość po przerwanym biegu.
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

    /// Ścieżka pliku o podanej nazwie wewnątrz katalogu, bez założenia pliku na dysku.
    pub fn plik(&self, nazwa: &str) -> PathBuf {
        self.sciezka.join(nazwa)
    }

    /// Tworzy w katalogu plik o podanej treści, dopuszczając treść pustą, i oddaje jego ścieżkę.
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
    /// Sprząta katalog wraz z zawartością, przemilczając niepowodzenie, by nie zasłonić paniki testu.
    fn drop(&mut self) {
        let _ = std::fs::remove_dir_all(&self.sciezka);
    }
}
