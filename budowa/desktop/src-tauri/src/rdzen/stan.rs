//! Moduł opisuje stan rdzenia przekazywany oknu i zasobnikowi, licząc go przy każdym
//! pytaniu z ustawień obowiązujących.

use serde::Serialize;

use super::nasluch;
use crate::dziennik;
use crate::ustawienia::Ustawienia;

/// Opis stanu rdzenia przekazywany interfejsowi i oknom dialogowym powłoki, złożony przy
/// każdym zapytaniu na nowo.
#[derive(Clone, Debug, Serialize)]
pub struct OpisRdzenia {
    /// Czy rdzeń odpowiada pod wskazanym adresem w chwili zapytania.
    pub pracuje: bool,
    /// Adres HTTP rdzenia obowiązujący; brak, dopóki wskazania nie złożono.
    pub adres: Option<String>,
    /// Ścieżka dziennika powłoki.
    pub dziennik: String,
    /// Zdanie opisujące stan — także wtedy, gdy rdzeń nie odpowiada.
    pub opis: String,
}

/// Składa opis stanu rdzenia z ustawień obowiązujących, pytając serwer wskazany, nigdy
/// pętlę zwrotną urządzenia.
pub fn opisz(ustawienia: &Ustawienia) -> OpisRdzenia {
    let dziennik = dziennik::sciezka().display().to_string();

    let Some(wskazane) = ustawienia.wskazanie() else {
        return OpisRdzenia {
            pracuje: false,
            adres: None,
            dziennik,
            opis: "Wskazanie rdzenia nie zostało jeszcze złożone — powłoka czeka, aż Operator \
                   poda serwer wdrożenia, na którym stoi jego rdzeń."
                .to_string(),
        };
    };

    let adres = wskazane.adres_http();
    let pracuje = nasluch::odpowiada_pod(&wskazane.host, wskazane.port);
    let opis = if pracuje {
        format!("Rdzeń odpowiada pod {adres}.")
    } else {
        format!(
            "Pod adresem {adres} nikt nie odpowiada. Sprawdź łączność z serwerem wdrożenia \
             i to, czy rdzeń tam pracuje."
        )
    };

    OpisRdzenia {
        pracuje,
        adres: Some(adres),
        dziennik,
        opis,
    }
}
