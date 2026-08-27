//! Rdzeń, z którym rozmawia powłoka — stojący na serwerze wdrożenia.
//!
//! Powłoka rdzenia nie stawia i nie wygasza: nie ma go w instalce i nie ma go na
//! urządzeniu Operatora. Zostaje jej wobec rdzenia dwoje: rozpoznać, czy pod
//! wskazanym adresem ktoś odpowiada (`nasluch`), i złożyć z tego opis stanu dla
//! okna oraz zasobnika (`stan`).

pub mod nasluch;
pub mod stan;

pub use stan::{opisz, OpisRdzenia};
