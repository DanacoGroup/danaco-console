//! Moduł instaluje hak paniki procesu, który przy nieudanym starcie zapisuje
//! przyczynę w dzienniku i pokazuje Operatorowi natywny komunikat błędu, zanim proces zniknie bez śladu.

use std::panic;

use crate::dziennik;

/// Instaluje hak paniki procesu jako pierwszą instrukcję funkcji `main`, zanim jakiekolwiek
/// okno powłoki zdąży powstać.
pub fn zainstaluj() {
    panic::set_hook(Box::new(|informacja| {
        let tresc = format!("KRYTYCZNE — powłoka kończy się awaryjnie: {informacja}");
        dziennik::dopisz(&tresc);
        pokaz_operatorowi(&tresc);
    }));
}

/// Pokazuje Operatorowi natywne okno komunikatu — bez zależności od kontekstu
/// Tauri, którego w chwili tej paniki może jeszcze nie być. `user32.dll` jest
/// częścią systemu Windows, więc wywołanie nie dodaje żadnej zależności
/// w `Cargo.toml`.
#[cfg(windows)]
fn pokaz_operatorowi(tresc: &str) {
    use std::ffi::{c_void, OsStr};
    use std::os::windows::ffi::OsStrExt;

    fn na_wide(tekst: &str) -> Vec<u16> {
        OsStr::new(tekst)
            .encode_wide()
            .chain(std::iter::once(0))
            .collect()
    }

    let tytul = na_wide("Danaco Console — start się nie powiódł");
    let tresc_szeroka = na_wide(tresc);

    const MB_ICONERROR: u32 = 0x0000_0010;
    const MB_OK: u32 = 0x0000_0000;

    extern "system" {
        fn MessageBoxW(hwnd: *mut c_void, tekst: *const u16, tytul: *const u16, typ: u32) -> i32;
    }

    unsafe {
        MessageBoxW(
            std::ptr::null_mut(),
            tresc_szeroka.as_ptr(),
            tytul.as_ptr(),
            MB_ICONERROR | MB_OK,
        );
    }
}

/// Poza Windows nie ma jednego pewnego API na natywny komunikat — zostaje
/// wyjście standardowe błędu, obok wpisu do dziennika powyżej.
#[cfg(not(windows))]
fn pokaz_operatorowi(tresc: &str) {
    eprintln!("{tresc}");
}
