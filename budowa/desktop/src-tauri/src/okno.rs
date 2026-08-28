//! Moduł otwiera jedyne okno natywne powłoki, niosące wkompilowany pakiet interfejsu,
//! bez adresu zapasowego ani rdzenia.

use tauri::{AppHandle, Manager, WebviewUrl, WebviewWindow, WebviewWindowBuilder};

use crate::zamkniecie;

/// Etykieta okna głównego powłoki, jedynego okna natywnego aplikacji; ta sama wartość występuje w pliku uprawnień domyślnych aplikacji.
pub const ETYKIETA: &str = "glowne";

/// Tytuł okna głównego powłoki, widoczny na pasku systemu operacyjnego przy oknie oraz w jego własnej belce tytułowej.
pub const TYTUL: &str = "Danaco Console";

/// Otwiera okno główne powłoki z wkompilowanym pakietem interfejsu i wiąże jego zdarzenia systemowe z obsługą zamknięcia okna.
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

/// Przywraca okno główne powłoki z ukrycia w zasobniku systemowym i ustawia na nim uwagę systemu okienkowego.
pub fn pokaz(aplikacja: &AppHandle) {
    if let Some(okno) = aplikacja.get_webview_window(ETYKIETA) {
        let _ = okno.show();
        let _ = okno.unminimize();
        let _ = okno.set_focus();
    }
}
