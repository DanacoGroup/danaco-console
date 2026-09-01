//! Moduł prowadzi dziennik powłoki: jedyne miejsce, w którym powłoka bez konsoli
//! zostawia ślad o swoim starcie, wskazaniu rdzenia i aktualizacji.

use std::env;
use std::fs::{create_dir_all, OpenOptions};
use std::path::PathBuf;
use std::sync::OnceLock;

/// Nazwa pliku dziennika powłoki. Jeden wzór nazw z plikiem nastaw
/// (`nastawy::NAZWA_PLIKU`), bo oba pliki leżą w jednym katalogu i oba należą
/// do powłoki.
pub const NAZWA_PLIKU: &str = "powloka-dziennik.log";

/// Nazwa pliku wskazania katalogu danych, zapisywanego przez instalator
/// w katalogu programu z wyboru kroku 4 (`instalator/src-tauri/src/pobranie.rs`,
/// `NAZWA_WSKAZANIA_KATALOGU_DANYCH` — obie stałe muszą być równe).
pub const NAZWA_WSKAZANIA_KATALOGU: &str = "powloka-katalog-danych.txt";

/// Katalog danych powłoki, własny i osobny od katalogu danych rdzenia:
/// wskazanie instalatora z pliku obok programu, inaczej `%LOCALAPPDATA%\DanacoConsole`.
/// Ustalany raz na proces — dziennik i nastawy pytają o niego przy każdym zapisie.
pub fn katalog_danych() -> PathBuf {
    static KATALOG: OnceLock<PathBuf> = OnceLock::new();
    KATALOG
        .get_or_init(|| wskazany_przez_instalator().unwrap_or_else(katalog_domyslny))
        .clone()
}

/// Katalog wskazany w kroku 4 instalatora; brak pliku, plik pusty albo
/// nieczytelny znaczy brak wskazania.
fn wskazany_przez_instalator() -> Option<PathBuf> {
    let plik = env::current_exe()
        .ok()?
        .parent()?
        .join(NAZWA_WSKAZANIA_KATALOGU);
    let tresc = std::fs::read_to_string(plik).ok()?;
    let sciezka = tresc.trim();
    (!sciezka.is_empty()).then(|| PathBuf::from(sciezka))
}

/// Katalog danych bez wskazania instalatora: `%LOCALAPPDATA%\DanacoConsole`.
fn katalog_domyslny() -> PathBuf {
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
