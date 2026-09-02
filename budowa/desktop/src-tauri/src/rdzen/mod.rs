//! Moduł opisuje rdzeń, z którym rozmawia powłoka: rozpoznaje łączność pod wskazanym
//! adresem i składa z niej opis stanu dla okna i zasobnika.

pub mod nasluch;
pub mod proces;
pub mod stan;

pub use stan::{opisz, OpisRdzenia};
