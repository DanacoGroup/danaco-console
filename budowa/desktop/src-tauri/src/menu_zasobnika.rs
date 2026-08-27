//! Rozstrzyganie pozycji menu zasobnika.
//!
//! Pozycja nieznana nie przerywa pracy powłoki — zostaje pominięta
//! (nieznana nazwa nie zrywa kanału).
//!
//! Czego w tym menu nie ma: zatrzymania rdzenia. Rdzeń stoi na serwerze
//! wdrożenia, powłoka go nie postawiła i nie ma czym go wygasić — pozycja
//! obiecywałaby władzę, której powłoka nie ma.

use tauri::{AppHandle, Manager};
use tauri_plugin_dialog::{DialogExt, MessageDialogKind};

use crate::dialog_katalogu;
use crate::okno;
use crate::rdzen;
use crate::ustawienia::Ustawienia;
use crate::zamkniecie;

/// Identyfikatory pozycji menu — jedyne miejsce, w którym te nazwy występują.
pub const POKAZ: &str = "powloka.pokaz";
pub const KATALOG: &str = "powloka.katalog";
pub const STAN: &str = "powloka.stan";
pub const ZAKONCZ: &str = "powloka.zakoncz";

/// Wykonuje czynność przypisaną pozycji menu.
pub fn obsluz(aplikacja: &AppHandle, identyfikator: &str) {
    match identyfikator {
        POKAZ => okno::pokaz(aplikacja),
        // Zasobnik nie ma sprawy, w której wskazuje — bierze napis domyślny
        // i rozgłasza wybór; adresat po stronie interfejsu decyduje, co z nim.
        KATALOG => dialog_katalogu::wybierz_i_rozglos(aplikacja),
        STAN => pokaz_stan(aplikacja),
        ZAKONCZ => zamkniecie::zakoncz_powloke(aplikacja),
        _ => {}
    }
}

/// Pokazuje opis stanu rdzenia w natywnym oknie komunikatu.
fn pokaz_stan(aplikacja: &AppHandle) {
    let ustawienia = aplikacja.state::<Ustawienia>();
    let stan = rdzen::opisz(&ustawienia);
    let tresc = format!(
        "{}\n\nAdres: {}\nOdpowiada: {}\nDziennik: {}",
        stan.opis,
        stan.adres.as_deref().unwrap_or("nie wskazano"),
        if stan.pracuje { "tak" } else { "nie" },
        stan.dziennik
    );
    komunikat(aplikacja, "Stan rdzenia", &tresc, MessageDialogKind::Info);
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
