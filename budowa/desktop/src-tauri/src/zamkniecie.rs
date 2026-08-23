//! Zachowanie powłoki przy zamknięciu okna.
//!
//! Zamknięcie okna nie kończy biegnących procesów sesji: okno jest widokiem,
//! nie właścicielem pracy. Zamknięcie chowa je do zasobnika — rdzeń pracuje
//! dalej, procesy sesji biegną, a ponowne otwarcie wraca do tej samej pracy.
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

/// Kończy powłokę, pozostawiając rdzeń i procesy sesji przy pracy.
/// Zatrzymanie rdzenia jest osobnym, jawnym poleceniem z zasobnika.
pub fn zakoncz_powloke(aplikacja: &AppHandle) {
    if let Some(okno) = aplikacja.get_webview_window(okno::ETYKIETA) {
        let _ = okno.hide();
    }
    aplikacja.exit(0);
}
