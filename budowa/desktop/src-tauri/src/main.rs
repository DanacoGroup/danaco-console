// Punkt wejścia powłoki natywnej Danaco Console.
//
// Wyłącznie kompozycja: odczyt ustawień, uruchomienie rdzenia w tle — tylko gdy
// wskazanie Operatora na to pozwala (trzy stany: wskazanie niezłożone, rdzeń na
// tym urządzeniu, rdzeń na serwerze; zob. `ustawienia.rs` i `wskazanie.rs`) —
// złożenie aplikacji Tauri i praca do jawnego zakończenia.
// Zero logiki, zero typów, zero obsługi zdarzeń — każda odpowiedzialność
// mieszka w osobnym module.
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod aktualizacja;
mod awaria_startu;
mod dialog_katalogu;
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
mod zrodlo_interfejsu;

use ustawienia::Ustawienia;

fn main() {
    // Pierwsza instrukcja z rozmysłem: łapie panikę Tauri z wnętrza `run()`,
    // gdy `montaz::zloz` (okno, zasobnik) zawiedzie — zob. `awaria_startu.rs`.
    awaria_startu::zainstaluj();

    let ustawienia = Ustawienia::ustal();
    // Trzy stany wskazania, trzy zachowania startu:
    //
    //   niezłożone — pierwsze uruchomienie po instalacji. Powłoka NIE stawia
    //     rdzenia: nie wie jeszcze, czy Operator pracuje z rdzeniem na tym
    //     urządzeniu, czy z rdzeniem na serwerze. Pyta o to okno
    //     (`wskazanie_rdzenia`, `wskaz_rdzen`), a proces staje po wskazaniu.
    //   rdzeń lokalny — powłoka stawia proces w tle, jak dotąd.
    //   rdzeń na serwerze — powłoka nie stawia niczego; drugi proces tworzyłby
    //     drugi, zbędny stan.
    let rdzen_w_tle = if !ustawienia.wskazanie_zlozone() {
        rdzen::uruchomienie::oczekuj_na_wskazanie(&ustawienia)
    } else if ustawienia.rdzen_lokalny() {
        rdzen::uruchom_w_tle(&ustawienia)
    } else {
        rdzen::uruchomienie::nie_stawiaj_lokalnie(&ustawienia)
    };

    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .manage(ustawienia)
        .manage(rdzen_w_tle)
        // Lista poleceń jest zamknięta: wchodzi na nią wyłącznie czynność,
        // której przeglądarka nie wykona sama — nie podmieni pliku aplikacji
        // ani nie postawi procesu na nowo. Zatrzymanie rdzenia poleceniem nie
        // jest i zostaje w zasobniku (powód w `polecenia.rs`).
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
