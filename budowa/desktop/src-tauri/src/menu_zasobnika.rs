//! Moduł rozstrzyga pozycję menu zasobnika w czynność powłoki, pomijając bez przerwania pracy pozycję nieznaną.

use tauri::{AppHandle, Manager};
use tauri_plugin_dialog::{DialogExt, MessageDialogKind};

use crate::dialog_katalogu;
use crate::okno;
use crate::rdzen;
use crate::ustawienia::Ustawienia;
use crate::zamkniecie;

/// Identyfikatory pozycji menu zasobnika — jedyne miejsce, w którym te napisy nazw pozycji występują.
pub const POKAZ: &str = "powloka.pokaz";
pub const KATALOG: &str = "powloka.katalog";
pub const STAN: &str = "powloka.stan";
pub const ZAKONCZ: &str = "powloka.zakoncz";

/// Wykonuje czynność przypisaną wskazanej pozycji menu zasobnika, pomijając pozycję o nieznanym identyfikatorze.
pub fn obsluz(aplikacja: &AppHandle, identyfikator: &str) {
    match identyfikator {
        POKAZ => okno::pokaz(aplikacja),
        // Zasobnik bierze napis domyślny; adresat po stronie interfejsu decyduje o rozgłoszonym wyborze.
        KATALOG => dialog_katalogu::wybierz_i_rozglos(aplikacja),
        STAN => pokaz_stan(aplikacja),
        ZAKONCZ => zamkniecie::zakoncz_powloke(aplikacja),
        _ => {}
    }
}

/// Pokazuje opis bieżącego stanu rdzenia w natywnym oknie komunikatu, wraz z adresem i stanem dziennika.
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

/// Wyświetla natywny komunikat o podanym tytule, treści i rodzaju, bez wstrzymywania wątku okna powłoki.
fn komunikat(aplikacja: &AppHandle, tytul: &str, tresc: &str, rodzaj: MessageDialogKind) {
    aplikacja
        .dialog()
        .message(tresc)
        .title(tytul)
        .kind(rodzaj)
        .show(|_| {});
}
