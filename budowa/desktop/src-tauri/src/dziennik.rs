//! Moduł prowadzi dziennik powłoki: jedyne miejsce, w którym powłoka bez konsoli
//! zostawia ślad o swoim starcie, wskazaniu rdzenia i aktualizacji.

use std::env;
use std::fs::{create_dir_all, OpenOptions};
use std::path::PathBuf;

/// Nazwa pliku dziennika powłoki. Jeden wzór nazw z plikiem nastaw
/// (`nastawy::NAZWA_PLIKU`), bo oba pliki leżą w jednym katalogu i oba należą
/// do powłoki.
pub const NAZWA_PLIKU: &str = "powloka-dziennik.log";

/// Katalog danych powłoki, własny i osobny od katalogu danych rdzenia: `%LOCALAPPDATA%\DanacoConsole`.
pub fn katalog_danych() -> PathBuf {
    let baza = env::var("LOCALAPPDATA")
        .or_else(|_| env::var("HOME"))
        .unwrap_or_else(|_| ".".to_string());
    PathBuf::from(baza).join("DanacoConsole")
}

/// Ścieżka pliku dziennika powłoki, złożona z katalogu danych powłoki i nazwy pliku dziennika powłoki.
pub fn sciezka() -> PathBuf {
    katalog_danych().join(NAZWA_PLIKU)
}

/// Dopisuje wiersz do dziennika powłoki, pomijając niepowodzenie zapisu bez wstrzymania startu powłoki.
pub fn dopisz(wiersz: &str) {
    let sciezka = sciezka();
    if let Some(katalog) = sciezka.parent() {
        let _ = create_dir_all(katalog);
    }
    if let Ok(mut plik) = OpenOptions::new().create(true).append(true).open(&sciezka) {
        use std::io::Write;
        let _ = writeln!(plik, "powłoka: {wiersz}");
    }
}
