//! Opis stanu rdzenia przekazywany oknu i zasobnikowi.
//!
//! Stanu nie ma czego przechowywać między wywołaniami. Rdzeń stoi na serwerze
//! wdrożenia, więc powłoka nie zna jego procesu, nie zna chwili jego startu
//! i nie ma nad nim władzy — wie wyłącznie, pod jakim adresem go szukać
//! (`Ustawienia`) i czy ten adres w tej chwili odpowiada (`nasluch`). Opis jest
//! więc liczony przy każdym pytaniu z ustawień obowiązujących, a nie odczytywany
//! z kopii spod zamka: kopia rozjeżdżałaby się ze wskazaniem złożonym w oknie.

use serde::Serialize;

use super::nasluch;
use crate::dziennik;
use crate::ustawienia::Ustawienia;

/// Opis stanu rdzenia przekazywany interfejsowi i oknom dialogowym.
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

/// Składa opis stanu rdzenia z ustawień obowiązujących.
///
/// Rozpoznanie pyta serwer wskazany, nie pętlę zwrotną: rdzeń stoi na serwerze
/// wdrożenia, więc pytanie pętli zwrotnej dawałoby fałsz niezależnie od jego
/// rzeczywistego stanu. Serwer nierozwiązywalny daje `pracuje: false`, nie panikę.
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
