//! Rdzeń w tle — odnalezienie binarki, uruchomienie, rozpoznanie nasłuchu
//! i uchwyt stanu. Każda z tych odpowiedzialności ma własny plik.

pub mod dziennik;
pub mod lokalizacja;
pub mod nasluch;
pub mod pakiet_klienta;
pub mod uchwyt;
pub mod uruchomienie;

pub use uchwyt::{OpisRdzenia, UchwytRdzenia};
pub use uruchomienie::uruchom_w_tle;
