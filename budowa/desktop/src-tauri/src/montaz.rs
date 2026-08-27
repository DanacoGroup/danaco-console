//! Montaż powłoki po starcie aplikacji: okno, zasobnik, wpis do dziennika.
//!
//! Wszystko, co powstaje w chwili startu, powstaje tutaj — punkt wejścia
//! wyłącznie komponuje.

use tauri::{App, Manager};

use crate::dziennik;
use crate::okno;
use crate::rdzen;
use crate::ustawienia::Ustawienia;
use crate::zasobnik;

/// Składa powłokę: otwiera okno z pakietem interfejsu i stawia ikonę zasobnika.
pub fn zloz(aplikacja: &mut App) -> Result<(), Box<dyn std::error::Error>> {
    let ustawienia = aplikacja.state::<Ustawienia>().inner().clone();
    dziennik::dopisz(&rdzen::opisz(&ustawienia).opis);

    okno::otworz(aplikacja.handle())?;
    zasobnik::zbuduj(aplikacja.handle())?;
    dziennik::dopisz("okno otwarte, ikona w zasobniku ustawiona");
    Ok(())
}
