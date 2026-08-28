// Skrypt budowania powłoki: osadza zasoby okna, ikonę i manifest platformy.
// Cała treść pochodzi z `tauri.conf.json`; tutaj nie ma żadnej konfiguracji.
use std::path::Path;

fn main() {
    oglos_zaleznosc_od_pakietu_klienta();
    tauri_build::build();
}

/// Ogłasza cargo, że wynik budowy zależy od pakietu klienta, ponieważ frontendDist osadza client/dist w chwili budowy, a cargo śledzi wyłącznie pliki źródłowe i manifest, nie katalog zasobów.
fn oglos_zaleznosc_od_pakietu_klienta() {
    let pakiet = Path::new("../../client/dist");
    println!("cargo:rerun-if-changed={}", pakiet.display());

    // Katalog assets/ niesie treść pakietu — bez zejścia weń zmiana zasobów nie wznowi budowy.
    if let Ok(wpisy) = std::fs::read_dir(pakiet) {
        for wpis in wpisy.flatten() {
            println!("cargo:rerun-if-changed={}", wpis.path().display());
        }
    }
}
