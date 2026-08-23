//! Plik dziennika rdzenia uruchomionego przez powłokę.
//!
//! Powłoka pracuje bez konsoli (podsystem okienkowy), więc wyjście rdzenia
//! musi mieć dokąd trafić. Zamiast je porzucać, dopisujemy je do pliku
//! w katalogu danych — ten sam katalog, którego używa rdzeń.
//! Brak możliwości otwarcia pliku nie wstrzymuje startu: wyjście idzie
//! wtedy do urządzenia pustego.

use std::env;
use std::fs::{create_dir_all, OpenOptions};
use std::path::PathBuf;
use std::process::Stdio;

/// Nazwa pliku dziennika rdzenia prowadzonego przez powłokę.
pub const NAZWA_PLIKU: &str = "rdzen-powloki.log";

/// Zmienna katalogu danych — własność rdzenia, powłoka wyłącznie czyta.
const ZMIENNA_KATALOG_DANYCH: &str = "DANACO_KATALOG_DANYCH";

/// Katalog danych rdzenia: wskazany zmienną albo `%LOCALAPPDATA%\DanacoConsole`.
pub fn katalog_danych() -> PathBuf {
    if let Ok(wskazany) = env::var(ZMIENNA_KATALOG_DANYCH) {
        if !wskazany.trim().is_empty() {
            return PathBuf::from(wskazany.trim());
        }
    }
    let baza = env::var("LOCALAPPDATA")
        .or_else(|_| env::var("HOME"))
        .unwrap_or_else(|_| ".".to_string());
    PathBuf::from(baza).join("DanacoConsole")
}

/// Ścieżka pliku dziennika rdzenia.
pub fn sciezka() -> PathBuf {
    katalog_danych().join(NAZWA_PLIKU)
}

/// Dopisuje wiersz powłoki do tego samego dziennika, w którym leży wyjście
/// rdzenia. Niepowodzenie zapisu jest pomijane — dziennik nie jest bramą.
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

/// Otwiera dziennik do dopisywania i zwraca dwa strumienie — wyjście i błąd.
/// Przy niepowodzeniu zwraca strumienie puste, nie przerywając startu.
pub fn strumienie() -> (Stdio, Stdio) {
    let sciezka = sciezka();
    if let Some(katalog) = sciezka.parent() {
        let _ = create_dir_all(katalog);
    }
    let otworz = || {
        OpenOptions::new()
            .create(true)
            .append(true)
            .open(&sciezka)
            .ok()
    };
    match (otworz(), otworz()) {
        (Some(wyjscie), Some(blad)) => (Stdio::from(wyjscie), Stdio::from(blad)),
        _ => (Stdio::null(), Stdio::null()),
    }
}
