//! Okno główne powłoki — jedno okno natywne z interfejsem.
//!
//! Interfejs pochodzi z pakietu wkompilowanego w powłokę (`frontendDist`
//! w `tauri.conf.json`) i z żadnego innego miejsca. To jest cały produkt
//! na urządzeniu Operatora: okno wraz z interfejsem, bez rdzenia. Dlatego
//! adresu strony nie ma czego rozstrzygać w czasie pracy — nastawa budowania
//! rozstrzyga go raz, a powłoka nie niesie żadnego adresu zapasowego, w tym
//! adresu serwera rozwojowego.

use tauri::{AppHandle, Manager, WebviewUrl, WebviewWindow, WebviewWindowBuilder};

use crate::zamkniecie;

/// Etykieta okna głównego. Ta sama wartość występuje w `capabilities/domyslne.json`.
pub const ETYKIETA: &str = "glowne";

/// Tytuł okna widoczny na pasku systemu.
pub const TYTUL: &str = "Danaco Console";

/// Otwiera okno główne z pakietem interfejsu.
pub fn otworz(aplikacja: &AppHandle) -> tauri::Result<WebviewWindow> {
    let okno = WebviewWindowBuilder::new(aplikacja, ETYKIETA, WebviewUrl::default())
        .title(TYTUL)
        .inner_size(1440.0, 900.0)
        .min_inner_size(960.0, 640.0)
        .resizable(true)
        .center()
        .visible(true)
        .build()?;

    let kopia = okno.clone();
    okno.on_window_event(move |zdarzenie| zamkniecie::obsluz(&kopia, zdarzenie));
    Ok(okno)
}

/// Przywraca okno główne z zasobnika i ustawia na nim uwagę.
pub fn pokaz(aplikacja: &AppHandle) {
    if let Some(okno) = aplikacja.get_webview_window(ETYKIETA) {
        let _ = okno.show();
        let _ = okno.unminimize();
        let _ = okno.set_focus();
    }
}
