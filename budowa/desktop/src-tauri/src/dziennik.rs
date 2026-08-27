//! Dziennik powłoki — jedyne miejsce, w którym powłoka zostawia ślad.
//!
//! Powłoka pracuje bez konsoli (podsystem okienkowy, `main.rs`), więc to, co ma
//! do powiedzenia o swoim starcie, wskazaniu rdzenia i aktualizacji, musi mieć
//! dokąd trafić. Brak możliwości otwarcia pliku nie wstrzymuje startu — zapis
//! jest wtedy pomijany, bo dziennik nie jest bramą.
//!
//! Katalog jest własnością powłoki, nie rdzenia. Zmiennej `DANACO_KATALOG_DANYCH`
//! ten moduł nie czyta: należy ona do rdzenia, a rdzeń stoi na serwerze
//! wdrożenia i wskazuje nią katalog na TAMTEJ maszynie. Powłoka pisze u siebie,
//! obok pliku nastaw (`nastawy.rs`), bo oba pliki są jej własne.

use std::env;
use std::fs::{create_dir_all, OpenOptions};
use std::path::PathBuf;

/// Nazwa pliku dziennika powłoki. Jeden wzór nazw z plikiem nastaw
/// (`nastawy::NAZWA_PLIKU`), bo oba pliki leżą w jednym katalogu i oba należą
/// do powłoki.
pub const NAZWA_PLIKU: &str = "powloka-dziennik.log";

/// Katalog danych powłoki: `%LOCALAPPDATA%\DanacoConsole`.
pub fn katalog_danych() -> PathBuf {
    let baza = env::var("LOCALAPPDATA")
        .or_else(|_| env::var("HOME"))
        .unwrap_or_else(|_| ".".to_string());
    PathBuf::from(baza).join("DanacoConsole")
}

/// Ścieżka pliku dziennika powłoki.
pub fn sciezka() -> PathBuf {
    katalog_danych().join(NAZWA_PLIKU)
}

/// Dopisuje wiersz do dziennika. Niepowodzenie zapisu jest pomijane —
/// dziennik nie jest bramą.
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
