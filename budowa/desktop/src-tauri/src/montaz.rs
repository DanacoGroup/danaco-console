//! Montaż powłoki po starcie aplikacji: okno, zasobnik, wpis do dziennika.
//!
//! Wszystko, co powstaje w chwili startu, powstaje tutaj — punkt wejścia
//! wyłącznie komponuje.

use tauri::{App, Manager};

use crate::okno;
use crate::rdzen::{dziennik, UchwytRdzenia};
use crate::ustawienia::Ustawienia;
use crate::zasobnik;
use crate::zrodlo_interfejsu;

/// Składa powłokę: otwiera okno pod ustalonym adresem i stawia ikonę zasobnika.
pub fn zloz(aplikacja: &mut App) -> Result<(), Box<dyn std::error::Error>> {
    let ustawienia = aplikacja.state::<Ustawienia>().inner().clone();
    dziennik::dopisz(&aplikacja.state::<UchwytRdzenia>().opis().opis);

    let zrodlo = zrodlo_interfejsu::ustal(&ustawienia);
    dziennik::dopisz(&zrodlo.opis);

    okno::otworz(aplikacja.handle(), zrodlo.adres)?;
    zasobnik::zbuduj(aplikacja.handle())?;
    dziennik::dopisz("okno otwarte, ikona w zasobniku ustawiona");
    Ok(())
}
