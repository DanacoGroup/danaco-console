//! Montaż powłoki po starcie aplikacji: okno, zasobnik, wpis do dziennika.
//!
//! Wszystko, co powstaje w chwili startu, powstaje tutaj — punkt wejścia
//! wyłącznie komponuje.

use std::process::Child;
use std::sync::Mutex;

use tauri::{App, Manager};

use crate::dziennik;
use crate::okno;
use crate::rdzen;
use crate::ustawienia::Ustawienia;
use crate::zasobnik;

/// Proces poboczny rdzenia, gdy powłoka go wystartowała; zamknięcie powłoki go wygasza.
pub struct ProcesRdzenia(pub Mutex<Option<Child>>);

/// Składa powłokę: otwiera okno z wkompilowanym pakietem interfejsu i stawia ikonę zasobnika systemowego.
pub fn zloz(aplikacja: &mut App) -> Result<(), Box<dyn std::error::Error>> {
    let ustawienia = aplikacja.state::<Ustawienia>().inner().clone();
    /* Rdzeń obok powłoki staje przed oknem: sekret nawiązania wchodzi w stronę
       skryptem wstępnym, więc musi być znany, zanim okno powstanie. */
    let dziecko = rdzen::proces::uruchom_obok(&ustawienia);
    /* Rdzeń postawiony przed chwilą nie zdążył jeszcze zająć portu, więc pytanie
       go teraz dałoby w dzienniku rozpoznanie usterki łączności zamiast opisu
       rozruchu. Stan rdzenia stojącego gdzie indziej opisuje się od razu. */
    let stoi_obok = dziecko.is_some();
    aplikacja.manage(ProcesRdzenia(Mutex::new(dziecko)));
    if stoi_obok {
        dziennik::dopisz("rdzeń stoi obok powłoki i właśnie wstaje");
    } else {
        dziennik::dopisz(&rdzen::opisz(&ustawienia).opis);
    }

    okno::otworz(aplikacja.handle())?;
    zasobnik::zbuduj(aplikacja.handle())?;
    dziennik::dopisz("okno otwarte, ikona w zasobniku ustawiona");
    Ok(())
}
