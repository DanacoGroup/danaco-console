//! Hak paniki — ostatnia linia obrony przed ciszą przy nieudanym starcie.
//!
//! `tauri::App::run` panikuje samodzielnie, gdy hak `setup` (tu: `montaz::zloz`,
//! czyli otwarcie okna i budowa zasobnika) zwróci błąd — działanie
//! udokumentowane wprost w źródłach `tauri` (`app.rs`, funkcja wolna `setup`,
//! wołana z domknięcia pętli zdarzeń przy `RuntimeRunEvent::Ready`):
//! `panic!("Failed to setup app: {e}")`. Dzieje się to wewnątrz `run()`, zanim
//! jakiekolwiek okno powstanie — powłoka nie tworzy okien deklaratywnie
//! w `tauri.conf.json`, więc do chwili tej paniki nie istnieje żadne.
//!
//! Powłoka pracuje w podsystemie okienkowym poza kompilacją debug
//! (`windows_subsystem = "windows"`, `main.rs`) i profil release ma
//! `panic = "abort"` (`Cargo.toml`). Bez własnego haka ta panika kończy
//! proces natychmiast: bez konsoli (nie ma dokąd wypisać), bez okna (jeszcze
//! nie powstało) i bez wpisu w dzienniku (`dziennik::dopisz` z `montaz::zloz`
//! nie zdążyło zapisać przyczyny — panika przerywa w miejscu błędu). Proces
//! znika wtedy bez śladu dla Operatora.
//!
//! Hak paniki uruchamia się zawsze przed odwinięciem/przerwaniem procesu —
//! nawet przy `panic = "abort"` (gwarancja `std::panic`, niezależna od
//! strategii paniki). To jedyne miejsce, w którym powłoka może jeszcze coś
//! powiedzieć Operatorowi, zanim zniknie — dlatego instalacja haka jest
//! pierwszą instrukcją `main`, przed jakąkolwiek inną pracą.

use std::panic;

use crate::rdzen::dziennik;

/// Instaluje hak paniki procesu. Wołać jako pierwszą instrukcję `main`.
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
