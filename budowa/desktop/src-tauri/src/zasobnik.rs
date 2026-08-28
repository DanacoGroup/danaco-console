//! Ikona powłoki w zasobniku systemowym.
//!
//! Zasobnik jest miejscem, w którym powłoka żyje po zamknięciu okna: praca
//! toczy się dalej w rdzeniu na serwerze wdrożenia, a Operator ma stąd dostęp
//! do okna, wyboru katalogu i stanu rdzenia.

use tauri::menu::{Menu, MenuItem, PredefinedMenuItem};
use tauri::tray::{TrayIconBuilder, TrayIconEvent};
use tauri::AppHandle;

use crate::menu_zasobnika;
use crate::okno;

/// Identyfikator ikony powłoki w zasobniku systemowym, używany przy budowie ikony i obsłudze
/// jej zdarzeń kliknięcia.
pub const ETYKIETA: &str = "powloka";

/// Buduje ikonę zasobnika systemowego wraz z menu i podpina obsługę zdarzeń kliknięcia oraz
/// zdarzeń pozycji menu.
pub fn zbuduj(aplikacja: &AppHandle) -> tauri::Result<()> {
    let menu = zbuduj_menu(aplikacja)?;
    let mut budowniczy = TrayIconBuilder::with_id(ETYKIETA)
        .tooltip(okno::TYTUL)
        .menu(&menu)
        .on_menu_event(|aplikacja, zdarzenie| {
            menu_zasobnika::obsluz(aplikacja, zdarzenie.id.as_ref());
        })
        .on_tray_icon_event(|zasobnik, zdarzenie| {
            if let TrayIconEvent::DoubleClick { .. } = zdarzenie {
                okno::pokaz(zasobnik.app_handle());
            }
        });
    if let Some(ikona) = aplikacja.default_window_icon() {
        budowniczy = budowniczy.icon(ikona.clone());
    }
    budowniczy.build(aplikacja)?;
    Ok(())
}

/// Składa menu zasobnika z pozycji stałych, każdej zawsze czynnej; identyfikatory pozycji
/// prowadzi osobno moduł menu.
fn zbuduj_menu(aplikacja: &AppHandle) -> tauri::Result<Menu<tauri::Wry>> {
    let pokaz = pozycja(aplikacja, menu_zasobnika::POKAZ, "Pokaż okno")?;
    let katalog = pozycja(aplikacja, menu_zasobnika::KATALOG, "Wskaż katalog roboczy…")?;
    let stan = pozycja(aplikacja, menu_zasobnika::STAN, "Stan rdzenia")?;
    let zakoncz = pozycja(aplikacja, menu_zasobnika::ZAKONCZ, "Zakończ powłokę")?;
    let kreska_gorna = PredefinedMenuItem::separator(aplikacja)?;
    let kreska_dolna = PredefinedMenuItem::separator(aplikacja)?;
    Menu::with_items(
        aplikacja,
        &[
            &pokaz,
            &katalog,
            &kreska_gorna,
            &stan,
            &kreska_dolna,
            &zakoncz,
        ],
    )
}

/// Tworzy pojedynczą pozycję menu zasobnika o podanym identyfikatorze i napisie; każda
/// pozycja jest zawsze czynna.
fn pozycja(
    aplikacja: &AppHandle,
    identyfikator: &str,
    napis: &str,
) -> tauri::Result<MenuItem<tauri::Wry>> {
    MenuItem::with_id(aplikacja, identyfikator, napis, true, None::<&str>)
}
