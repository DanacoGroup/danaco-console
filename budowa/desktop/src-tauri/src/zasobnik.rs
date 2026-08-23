//! Ikona powłoki w zasobniku systemowym.
//!
//! Zasobnik jest miejscem, w którym powłoka żyje po zamknięciu okna: rdzeń
//! pracuje, procesy sesji biegną, a Operator ma dostęp do okna, wyboru
//! katalogu, stanu rdzenia i jawnego zatrzymania.

use tauri::menu::{Menu, MenuItem, PredefinedMenuItem};
use tauri::tray::{TrayIconBuilder, TrayIconEvent};
use tauri::AppHandle;

use crate::menu_zasobnika;
use crate::okno;

/// Identyfikator ikony w zasobniku.
pub const ETYKIETA: &str = "powloka";

/// Buduje ikonę zasobnika wraz z menu.
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

/// Składa menu zasobnika. Identyfikatory pozycji prowadzi `menu_zasobnika`.
fn zbuduj_menu(aplikacja: &AppHandle) -> tauri::Result<Menu<tauri::Wry>> {
    let pokaz = pozycja(aplikacja, menu_zasobnika::POKAZ, "Pokaż okno")?;
    let katalog = pozycja(aplikacja, menu_zasobnika::KATALOG, "Wskaż katalog roboczy…")?;
    let stan = pozycja(aplikacja, menu_zasobnika::STAN, "Stan rdzenia")?;
    let zatrzymaj = pozycja(aplikacja, menu_zasobnika::ZATRZYMAJ, "Zatrzymaj rdzeń")?;
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
            &zatrzymaj,
            &kreska_dolna,
            &zakoncz,
        ],
    )
}

/// Tworzy pojedynczą pozycję menu; każda jest zawsze czynna, bez wyszarzeń.
fn pozycja(
    aplikacja: &AppHandle,
    identyfikator: &str,
    napis: &str,
) -> tauri::Result<MenuItem<tauri::Wry>> {
    MenuItem::with_id(aplikacja, identyfikator, napis, true, None::<&str>)
}
