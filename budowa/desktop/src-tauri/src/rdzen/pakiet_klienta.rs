//! Wskazanie rdzeniowi pakietu interfejsu (`client/dist`).
//!
//! Rdzeń serwuje pakiet klienta obok gniazda WebSocket
//! (`server/internal/transport/statyka.go`), a jego ścieżkę czyta ze zmiennej
//! `DANACO_KATALOG_KLIENTA`; bez niej szuka `client\dist` względem katalogu
//! bieżącego. Powłoka stawia rdzeń w jego własnym katalogu, więc ścieżkę
//! podaje wprost. Zmienna ustawiona przez Operatora ma pierwszeństwo —
//! powłoka jej nie nadpisuje.

use std::env;
use std::path::{Path, PathBuf};

/// Zmienna wskazująca rdzeniowi katalog pakietu interfejsu.
pub const ZMIENNA: &str = "DANACO_KATALOG_KLIENTA";

/// Szuka katalogu `client/dist`; zwraca pierwszy istniejący.
/// `rdzen` wskazuje binarkę rdzenia — obok niej leży pakiet w instalacji.
pub fn znajdz(rdzen: &Path) -> Option<PathBuf> {
    if let Ok(wskazany) = env::var(ZMIENNA) {
        if !wskazany.trim().is_empty() {
            return Some(PathBuf::from(wskazany.trim()));
        }
    }
    miejsca(rdzen).into_iter().find(|miejsce| miejsce.is_dir())
}

/// Wykaz miejsc pakietu interfejsu w kolejności rozstrzygania.
fn miejsca(rdzen: &Path) -> Vec<PathBuf> {
    let mut wykaz: Vec<PathBuf> = Vec::new();
    if let Some(katalog) = rdzen.parent() {
        wykaz.push(katalog.join("client").join("dist"));
        wykaz.push(katalog.join("dist"));
    }
    let manifest = PathBuf::from(env!("CARGO_MANIFEST_DIR"));
    if let Some(budowa) = manifest.parent().and_then(Path::parent) {
        wykaz.push(budowa.join("client").join("dist"));
    }
    wykaz
}
