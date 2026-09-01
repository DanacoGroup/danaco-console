// Punkt wejścia instalatora: składa aplikację Tauri, otwiera okno kreatora
// z warstwy projektowej i wystawia polecenia wiążące jego kroki z rzeczywistym
// stanem tej maszyny.
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod katalog;
mod pobranie;
mod stan_maszyny;

use tauri::{Emitter, WebviewUrl, WebviewWindowBuilder};

fn main() {
    tauri::Builder::default()
        .plugin(tauri_plugin_dialog::init())
        .invoke_handler(tauri::generate_handler![
            stan_maszyny::stan_maszyny,
            katalog::wybierz_katalog,
            katalog::sprawdz_katalog,
            pobranie::pobierz_skladniki,
            pobranie::zaloz_program,
            pobranie::uruchom_program,
        ])
        .setup(|aplikacja| {
            let procesor = stan_maszyny::architektura_kreatora();
            // Belkę tytułową niesie samo okno kreatora (`dn-kreator-belka`),
            // więc okno platformy stoi bez ramy. Wymiar jest ten sam co
            // `--dn-wym-okno-instalatora-*` w żetonach, a zarazem najmniejszy
            // dopuszczalny: okno kreatora nie zmienia rozmiaru między krokami,
            // więc poniżej tej miary nie ma go jak pokazać.
            let okno = WebviewWindowBuilder::new(
                aplikacja,
                "kreator",
                WebviewUrl::App(format!("index.html?procesor={procesor}").into()),
            )
            .title("Danaco Console — Instalator")
            .inner_size(1020.0, 720.0)
            .min_inner_size(1020.0, 720.0)
            .decorations(false)
            .center()
            .visible(true)
            .build()?;

            // Sprawdzenie kanału wydań to żądanie sieciowe do serwera wdrożenia:
            // wykonane przed zbudowaniem okna trzymało ekran pusty przez cały czas
            // odpowiedzi. Idzie więc po oknie, wątkiem blokującym, a wynik wraca
            // do okna zdarzeniem — krok 5 czyta go dopiero wtedy, gdy do niego dojdzie.
            tauri::async_runtime::spawn_blocking(move || {
                let stan = match pobranie::sprawdz_wstepnie(procesor) {
                    Ok(()) => pobranie::StanKanalu {
                        dostepny: true,
                        odmowa: None,
                    },
                    Err(odmowa) => pobranie::StanKanalu {
                        dostepny: false,
                        odmowa: Some(odmowa),
                    },
                };
                let _ = okno.emit(pobranie::ZDARZENIE_KANAL, stan);
            });
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("złożenie okna instalatora");
}
