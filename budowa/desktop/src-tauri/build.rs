// Skrypt budowania powłoki: osadza zasoby okna, ikonę i manifest platformy.
// Cała treść pochodzi z `tauri.conf.json`; tutaj nie ma żadnej konfiguracji.
use std::path::Path;

fn main() {
    oglos_zaleznosc_od_pakietu_klienta();
    tauri_build::build();
}

/// Ogłasza cargo, że wynik budowy zależy od pakietu klienta.
///
/// `frontendDist` osadza `client/dist` w binarium w chwili budowy, ale cargo
/// tego nie śledzi: pilnuje wyłącznie plików `.rs` i manifestu. Bez tego
/// ogłoszenia przebudowa po samej zmianie w kliencie kończy się kodem zero,
/// nie przebudowując niczego, a binarium zostaje z poprzednim pakietem.
fn oglos_zaleznosc_od_pakietu_klienta() {
    let pakiet = Path::new("../../client/dist");
    println!("cargo:rerun-if-changed={}", pakiet.display());

    // Sam katalog wystarcza cargo tylko dla wpisów pierwszego poziomu, a pakiet
    // trzyma treść w `assets/`. Bez zejścia głębiej zmiana w skrypcie albo
    // arkuszu stylów nie ruszyłaby przebudowy.
    if let Ok(wpisy) = std::fs::read_dir(pakiet) {
        for wpis in wpisy.flatten() {
            println!("cargo:rerun-if-changed={}", wpis.path().display());
        }
    }
}
