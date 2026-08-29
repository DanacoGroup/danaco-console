// Punkt wejścia instalatora: składa aplikację Tauri, otwiera okno kreatora
// z warstwy projektowej i wystawia polecenia wiążące jego kroki z rzeczywistym
// stanem tej maszyny.
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod katalog;
mod pobranie;
mod stan_maszyny;

use tauri::{WebviewUrl, WebviewWindowBuilder};

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .invoke_handler(tauri::generate_handler![
            stan_maszyny::stan_maszyny,
            katalog::wybierz_katalog,
            katalog::sprawdz_katalog,
            pobranie::pobierz_skladniki,
        ])
        .setup(|aplikacja| {
            let adres = zloz_adres_okna();
            WebviewWindowBuilder::new(aplikacja, "kreator", WebviewUrl::App(adres.into()))
                .title("Danaco Console — Instalator")
                .inner_size(1040.0, 720.0)
                .min_inner_size(900.0, 620.0)
                .center()
                .visible(true)
                .build()?;
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("złożenie okna instalatora");
}

/// Składa adres startowy okna: architekturę tej maszyny (krok 3) oraz — gdy
/// serwer wydań odmówił pliku już teraz — powód i zdanie odmowy dla kroku 5,
/// czytany raz, przy starcie kreatora.
fn zloz_adres_okna() -> String {
    let procesor = stan_maszyny::architektura_kreatora();
    let mut adres = format!("index.html?procesor={procesor}");

    match pobranie::sprawdz_wstepnie(procesor) {
        Ok(_) => {}
        Err(odmowa) => {
            adres.push_str("&stan=blad");
            adres.push_str("&blad_powod=");
            adres.push_str(&koduj_query(&odmowa.powod));
            adres.push_str("&blad_zdanie=");
            adres.push_str(&koduj_query(&odmowa.zdanie));
        }
    }
    adres
}

/// Kodowanie procentowe wartości parametru zapytania — bez dodatkowej zależności,
/// bo instalator potrzebuje go w jednym miejscu i dla napisów po polsku (UTF-8).
fn koduj_query(wartosc: &str) -> String {
    let mut wynik = String::with_capacity(wartosc.len());
    for bajt in wartosc.as_bytes() {
        match bajt {
            b'A'..=b'Z' | b'a'..=b'z' | b'0'..=b'9' | b'-' | b'_' | b'.' | b'~' => {
                wynik.push(*bajt as char)
            }
            _ => wynik.push_str(&format!("%{bajt:02X}")),
        }
    }
    wynik
}
