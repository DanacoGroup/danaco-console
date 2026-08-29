// Skrypt budowania instalatora: przed złożeniem osadza w `interfejs/warstwa/`
// świeżą kopię `design/zasoby/`, żeby okno kreatora wydawało się z warstwy
// projektowej bez odwołań uciekających poza katalog `frontendDist`, którego
// wymaga pakowanie. Warstwa projektowa jest tu wyłącznie źródłem do odczytu —
// ten skrypt do niej nie zapisuje ani jednego bajtu.
use std::io;
use std::path::Path;

fn main() {
    let zrodlo = Path::new("../../../design/zasoby");
    let cel = Path::new("../interfejs/warstwa/zasoby");

    println!("cargo:rerun-if-changed={}", zrodlo.display());
    kopiuj_katalog(zrodlo, cel).unwrap_or_else(|blad| {
        panic!(
            "nie udało się skopiować warstwy projektowej z {} do {}: {blad}",
            zrodlo.display(),
            cel.display()
        )
    });

    tauri_build::build();
}

/// Kopiuje katalog rekurencyjnie, nadpisując cel, żeby kopia zawsze odzwierciedlała
/// bieżący stan `design/zasoby/` w chwili budowy, nigdy poprzedni.
fn kopiuj_katalog(zrodlo: &Path, cel: &Path) -> io::Result<()> {
    if cel.exists() {
        std::fs::remove_dir_all(cel)?;
    }
    std::fs::create_dir_all(cel)?;
    for wpis in std::fs::read_dir(zrodlo)? {
        let wpis = wpis?;
        let typ = wpis.file_type()?;
        let cel_wpisu = cel.join(wpis.file_name());
        if typ.is_dir() {
            kopiuj_katalog(&wpis.path(), &cel_wpisu)?;
        } else {
            std::fs::copy(wpis.path(), &cel_wpisu)?;
        }
    }
    Ok(())
}
