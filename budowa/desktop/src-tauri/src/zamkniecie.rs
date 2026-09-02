//! Moduł chowa okno do zasobnika przy zamknięciu, bo okno jest widokiem sesji, nie jej
//! właścicielem, a praca toczy się w rdzeniu dalej.

use tauri::{AppHandle, Manager, WebviewWindow, WindowEvent};

use crate::montaz::ProcesRdzenia;
use crate::okno;
use crate::rdzen;

/// Obsługuje zdarzenia okna głównego powłoki, chowając je do zasobnika zamiast zamykać
/// przy żądaniu zamknięcia.
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
    wygas_rdzen_poboczny(aplikacja);
    aplikacja.exit(0);
}

/// Wygasza rdzeń wystartowany przez powłokę. Rdzeń serwera wdrożenia zostaje —
/// powłoka go nie stawiała i nie ma nad nim władzy.
fn wygas_rdzen_poboczny(aplikacja: &AppHandle) {
    let Some(stan) = aplikacja.try_state::<ProcesRdzenia>() else {
        return;
    };
    let Ok(mut zamek) = stan.0.lock() else {
        return;
    };
    if let Some(dziecko) = zamek.as_mut() {
        rdzen::proces::wygas(dziecko);
    }
}
