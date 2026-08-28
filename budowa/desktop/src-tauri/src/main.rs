// Plik stanowi punkt wejścia powłoki: czyta ustawienia, składa aplikację Tauri i pracuje do jawnego zakończenia procesu.
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod aktualizacja;
mod awaria_startu;
mod dialog_katalogu;
mod dziennik;
mod menu_zasobnika;
mod montaz;
mod nastawy;
mod okno;
mod polecenia;
mod rdzen;
mod ustawienia;
mod wskazanie;
mod zamkniecie;
mod zasobnik;

use ustawienia::Ustawienia;

fn main() {
    // Instaluje hak paniki jako pierwszą instrukcję, by przechwycić panikę składania aplikacji.
    awaria_startu::zainstaluj();

    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .manage(Ustawienia::ustal())
        // Lista poleceń jest zamknięta do czynności, których przeglądarka nie wykona samodzielnie.
        .invoke_handler(tauri::generate_handler![
            polecenia::wybierz_katalog_roboczy,
            polecenia::stan_rdzenia,
            polecenia::adres_rdzenia,
            polecenia::wskazanie_rdzenia,
            polecenia::wskaz_rdzen,
            polecenia::wykonaj_aktualizacje,
        ])
        .setup(montaz::zloz)
        .build(tauri::generate_context!())
        .expect("złożenie powłoki Danaco Console")
        .run(|_aplikacja, _zdarzenie| {});
}
