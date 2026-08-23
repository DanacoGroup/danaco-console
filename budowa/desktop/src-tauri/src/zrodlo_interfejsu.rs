//! Rozstrzygnięcie, skąd okno bierze interfejs.
//!
//! Kolejność, od wskazania jawnego do wartości domyślnej:
//!
//! 1. `DANACO_ADRES_INTERFEJSU` — serwer rozwojowy albo inny host;
//! 2. tryb deweloperski `cargo tauri dev` — serwer rozwojowy powłoki;
//! 3. rdzeń nasłuchujący lokalnie — pakiet `client/dist` serwowany przez rdzeń
//!    (`server/internal/transport/statyka.go`);
//! 4. pakiet `client/dist` osadzony w pliku wykonywalnym.
//!
//! Pierwszeństwo rdzenia przed pakietem osadzonym ma powód praktyczny:
//! klient wylicza adres gniazda z `location.hostname`
//! (`client/src/polaczenie/adres-rdzenia.ts`). Strona podana przez rdzeń ma
//! `hostname` równy `127.0.0.1` i trafia w gniazdo; strona z zasobu osadzonego
//! ma `hostname` równy `tauri.localhost` i w gniazdo nie trafia. Pakiet osadzony
//! pozostaje więc wyjściem awaryjnym, używanym, gdy rdzeń nie odpowiada.

use tauri::{Url, WebviewUrl};

use crate::rdzen::nasluch;
use crate::ustawienia::Ustawienia;

/// Wybrane źródło interfejsu wraz z opisem do dziennika.
pub struct Zrodlo {
    pub adres: WebviewUrl,
    pub opis: String,
}

/// Rozstrzyga źródło interfejsu dla okna głównego.
pub fn ustal(ustawienia: &Ustawienia) -> Zrodlo {
    if let Some(wskazany) = ustawienia.adres_interfejsu.as_deref() {
        return match Url::parse(wskazany) {
            Ok(adres) => zewnetrzne(adres, "wskazany zmienną DANACO_ADRES_INTERFEJSU"),
            Err(blad) => pakiet(&format!(
                "adres {wskazany} nie jest poprawny ({blad}) — pakiet osadzony"
            )),
        };
    }
    if tauri::is_dev() {
        return pakiet("serwer rozwojowy powłoki (tryb deweloperski)");
    }
    if nasluch::odpowiada(ustawienia.port()) {
        if let Ok(adres) = Url::parse(&ustawienia.adres_rdzenia_http()) {
            return zewnetrzne(adres, "pakiet client/dist serwowany przez rdzeń");
        }
    }
    pakiet("pakiet client/dist osadzony w powłoce")
}

/// Źródło zewnętrzne — strona pobierana po HTTP.
fn zewnetrzne(adres: Url, powod: &str) -> Zrodlo {
    let opis = format!("interfejs z {adres} — {powod}");
    Zrodlo {
        adres: WebviewUrl::External(adres),
        opis,
    }
}

/// Źródło osadzone — zasób z pakietu powłoki.
fn pakiet(powod: &str) -> Zrodlo {
    Zrodlo {
        adres: WebviewUrl::default(),
        opis: format!("interfejs z zasobu powłoki — {powod}"),
    }
}
