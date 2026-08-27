//! Zachowanie powłoki przy zamknięciu okna.
//!
//! Zamknięcie okna nie kończy pracy sesji: okno jest widokiem, nie właścicielem
//! pracy. Praca toczy się w rdzeniu na serwerze wdrożenia i biegnie dalej bez
//! względu na to, czy okno stoi otwarte. Zamknięcie chowa je więc do zasobnika,
//! a ponowne otwarcie wraca do tej samej sesji.
//!
//! To nie jest blokada: przycisk zamknięcia działa natychmiast i bez pytania,
//! zmienia się wyłącznie skutek.

use tauri::{AppHandle, Manager, WebviewWindow, WindowEvent};

use crate::okno;

/// Obsługuje zdarzenia okna głównego.
pub fn obsluz(okno: &WebviewWindow, zdarzenie: &WindowEvent) {
    if let WindowEvent::CloseRequested { api, .. } = zdarzenie {
        api.prevent_close();
        let _ = okno.hide();
    }
}

/// Kończy powłokę. Praca na serwerze wdrożenia toczy się dalej — powłoka nie ma
/// nad nią władzy i nie próbuje jej wygaszać.
pub fn zakoncz_powloke(aplikacja: &AppHandle) {
    if let Some(okno) = aplikacja.get_webview_window(okno::ETYKIETA) {
        let _ = okno.hide();
    }
    aplikacja.exit(0);
}
