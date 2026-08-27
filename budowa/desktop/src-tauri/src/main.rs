// Punkt wejścia powłoki Danaco Console.
//
// Wyłącznie kompozycja: odczyt ustawień, złożenie aplikacji Tauri i praca do
// jawnego zakończenia. Zero logiki, zero typów, zero obsługi zdarzeń — każda
// odpowiedzialność mieszka w osobnym module.
//
// Powłoka niesie okno wraz z wkompilowanym interfejsem i nie niesie rdzenia.
// Rdzeń stoi na serwerze wdrożenia, więc przy starcie nie ma czego stawiać ani
// na co czekać: powłoka czyta wskazanie, gdzie ten rdzeń szukać, i otwiera okno.
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
    // Pierwsza instrukcja z rozmysłem: łapie panikę Tauri z wnętrza `run()`,
    // gdy `montaz::zloz` (okno, zasobnik) zawiedzie — zob. `awaria_startu.rs`.
    awaria_startu::zainstaluj();

    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .manage(Ustawienia::ustal())
        // Lista poleceń jest zamknięta: wchodzi na nią wyłącznie czynność,
        // której przeglądarka nie wykona sama — natywne okno wyboru katalogu,
        // wskazanie serwera rdzenia i podmiana pliku aplikacji.
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
