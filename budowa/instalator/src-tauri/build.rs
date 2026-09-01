// Skrypt budowania instalatora: przed złożeniem osadza w `interfejs/warstwa/`
// świeżą kopię `design/zasoby/`, żeby okno kreatora wydawało się z warstwy
// projektowej bez odwołań uciekających poza katalog `frontendDist`, którego
// wymaga pakowanie. Warstwa projektowa jest tu wyłącznie źródłem do odczytu —
// ten skrypt do niej nie zapisuje ani jednego bajtu.
use std::io;
use std::path::{Path, PathBuf};

fn main() {
    let zrodlo = Path::new("../../../design/zasoby");
    let cel = Path::new("../interfejs/warstwa/zasoby");

    println!("cargo:rerun-if-changed={}", zrodlo.display());
    osadz_warstwe(zrodlo, cel).unwrap_or_else(|blad| {
        panic!(
            "nie udało się osadzić warstwy projektowej z {} w {}: {blad}",
            zrodlo.display(),
            cel.display()
        )
    });

    tauri_build::build();
}

/// Osadza kopię warstwy w katalogu przejściowym i podmienia nią kopię zastaną
/// dopiero po zakończeniu przepisywania. Budowa przerwana w połowie zostawia
/// wtedy katalog przejściowy, a nie warstwę niekompletną, z której następna
/// budowa złożyłaby okno bez części składników.
fn osadz_warstwe(zrodlo: &Path, cel: &Path) -> io::Result<()> {
    let przejsciowy = obok(cel, "nowy");
    let zastany = obok(cel, "stary");
    for pozostalosc in [&przejsciowy, &zastany] {
        if pozostalosc.exists() {
            std::fs::remove_dir_all(pozostalosc)?;
        }
    }

    kopiuj_katalog(zrodlo, &przejsciowy)?;
    if cel.exists() {
        std::fs::rename(cel, &zastany)?;
    }
    std::fs::rename(&przejsciowy, cel)?;
    if zastany.exists() {
        std::fs::remove_dir_all(&zastany)?;
    }
    Ok(())
}

/// Katalog roboczy obok celu, w tym samym katalogu nadrzędnym: `rename` działa
/// wyłącznie w obrębie jednego systemu plików, a katalog tymczasowy systemu
/// bywa na innym.
fn obok(cel: &Path, przyrostek: &str) -> PathBuf {
    let nazwa = cel.file_name().unwrap_or_default().to_string_lossy();
    match cel.parent() {
        Some(rodzic) => rodzic.join(format!("{nazwa}.{przyrostek}")),
        None => PathBuf::from(format!("{nazwa}.{przyrostek}")),
    }
}

/// Kopiuje katalog rekurencyjnie do celu pustego.
fn kopiuj_katalog(zrodlo: &Path, cel: &Path) -> io::Result<()> {
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
