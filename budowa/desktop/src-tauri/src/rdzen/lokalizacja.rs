//! Odnalezienie binarki rdzenia, którą powłoka ma uruchomić w tle.
//!
//! Kolejność miejsc biegnie od wskazania jawnego do miejsc wynikających
//! z układu katalogów projektu. Brak binarki nie jest błędem wstrzymującym —
//! okno otwiera się i tak, a przyczyna trafia do opisu stanu rdzenia.

use std::env;
use std::path::{Path, PathBuf};

/// Nazwa pliku wykonywalnego rdzenia na bieżącej platformie.
pub const NAZWA_BINARKI: &str = if cfg!(windows) {
    "danaco-console.exe"
} else {
    "danaco-console"
};

/// Wynik wyszukiwania: ścieżka albo wykaz miejsc, w których szukano.
pub enum Wynik {
    Znaleziona(PathBuf),
    Brak(Vec<PathBuf>),
}

/// Szuka binarki rdzenia. `wskazana` pochodzi ze zmiennej `DANACO_RDZEN`.
pub fn znajdz(wskazana: Option<&Path>) -> Wynik {
    let miejsca = miejsca(wskazana);
    for miejsce in &miejsca {
        if miejsce.is_file() {
            return Wynik::Znaleziona(miejsce.clone());
        }
    }
    Wynik::Brak(miejsca)
}

/// Buduje wykaz miejsc w kolejności rozstrzygania.
fn miejsca(wskazana: Option<&Path>) -> Vec<PathBuf> {
    let mut wykaz: Vec<PathBuf> = Vec::new();
    if let Some(sciezka) = wskazana {
        wykaz.push(sciezka.to_path_buf());
    }
    wykaz.extend(obok_powloki());
    wykaz.extend(w_drzewie_budowy());
    wykaz.push(PathBuf::from(r"C:\DanacoConsole_App").join(NAZWA_BINARKI));
    wykaz
}

/// Miejsca obok pliku wykonywalnego powłoki — układ instalacji na urządzeniu.
fn obok_powloki() -> Vec<PathBuf> {
    let Ok(powloka) = env::current_exe() else {
        return Vec::new();
    };
    let Some(katalog) = powloka.parent() else {
        return Vec::new();
    };
    vec![
        katalog.join(NAZWA_BINARKI),
        katalog.join("rdzen").join(NAZWA_BINARKI),
    ]
}

/// Miejsca w drzewie repozytorium — układ pracy deweloperskiej.
/// `CARGO_MANIFEST_DIR` wskazuje `budowa/desktop/src-tauri`.
fn w_drzewie_budowy() -> Vec<PathBuf> {
    let manifest = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    let Some(budowa) = manifest.parent().and_then(Path::parent) else {
        return Vec::new();
    };
    vec![
        budowa.join(NAZWA_BINARKI),
        budowa.join("bin").join(NAZWA_BINARKI),
        budowa.join("server").join(NAZWA_BINARKI),
    ]
}

/// Zapis wykazu miejsc do jednego wiersza opisu.
pub fn opis_miejsc(miejsca: &[PathBuf]) -> String {
    miejsca
        .iter()
        .map(|m| m.display().to_string())
        .collect::<Vec<_>>()
        .join("; ")
}
