//! Rozstrzyganie pozycji menu zasobnika.
//!
//! Pozycja nieznana nie przerywa pracy powłoki — zostaje pominięta
//! (nieznana nazwa nie zrywa kanału).

use tauri::{AppHandle, Manager};
use tauri_plugin_dialog::{DialogExt, MessageDialogKind};

use crate::dialog_katalogu;
use crate::okno;
use crate::rdzen::UchwytRdzenia;
use crate::zamkniecie;

/// Identyfikatory pozycji menu — jedyne miejsce, w którym te nazwy występują.
pub const POKAZ: &str = "powloka.pokaz";
pub const KATALOG: &str = "powloka.katalog";
pub const STAN: &str = "powloka.stan";
pub const ZATRZYMAJ: &str = "powloka.zatrzymaj-rdzen";
pub const ZAKONCZ: &str = "powloka.zakoncz";

/// Wykonuje czynność przypisaną pozycji menu.
pub fn obsluz(aplikacja: &AppHandle, identyfikator: &str) {
    match identyfikator {
        POKAZ => okno::pokaz(aplikacja),
        // Zasobnik nie ma sprawy, w której wskazuje — bierze napis domyślny
        // i rozgłasza wybór; adresat po stronie interfejsu decyduje, co z nim.
        KATALOG => dialog_katalogu::wybierz_i_rozglos(aplikacja),
        STAN => pokaz_stan(aplikacja),
        ZATRZYMAJ => zatrzymaj_rdzen(aplikacja),
        ZAKONCZ => zamkniecie::zakoncz_powloke(aplikacja),
        _ => {}
    }
}

/// Pokazuje opis stanu rdzenia w natywnym oknie komunikatu.
fn pokaz_stan(aplikacja: &AppHandle) {
    let uchwyt = aplikacja.state::<UchwytRdzenia>();
    let stan = uchwyt.opis();
    let tresc = format!(
        "{}\n\nAdres: {}\nNasłuch: {}\nDziennik: {}",
        stan.opis,
        stan.adres,
        if stan.pracuje { "tak" } else { "nie" },
        stan.dziennik
    );
    komunikat(aplikacja, "Stan rdzenia", &tresc, MessageDialogKind::Info);
}

/// Zatrzymuje rdzeń na jawne polecenie i melduje wynik.
fn zatrzymaj_rdzen(aplikacja: &AppHandle) {
    let wynik = aplikacja.state::<UchwytRdzenia>().zatrzymaj();
    match wynik {
        Ok(tresc) => komunikat(aplikacja, "Rdzeń", &tresc, MessageDialogKind::Info),
        Err(tresc) => komunikat(aplikacja, "Rdzeń", &tresc, MessageDialogKind::Warning),
    }
}

/// Wyświetla komunikat bez wstrzymywania wątku okna.
fn komunikat(aplikacja: &AppHandle, tytul: &str, tresc: &str, rodzaj: MessageDialogKind) {
    aplikacja
        .dialog()
        .message(tresc)
        .title(tytul)
        .kind(rodzaj)
        .show(|_| {});
}
