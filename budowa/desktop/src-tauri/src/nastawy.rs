//! Nastawy powłoki zapisane trwale — wskazanie Operatora, gdzie stoi rdzeń.
//!
//! Po co plik, skoro jest zmienna środowiska. Adresu serwera nie zna instalator
//! i znać go nie może: w chwili rozpakowania plików nikt jeszcze nie wie, pod
//! jaką nazwą stoi rdzeń tego Operatora. Wiedza pojawia się przy pierwszym
//! uruchomieniu okna i musi przetrwać jego zamknięcie — inaczej przy każdym
//! starcie powłoka wracałaby do rdzenia lokalnego, czyli do maszyny bez
//! arsenału.
//!
//! Zmienna `DANACO_HOST_RDZENIA` zostaje i stoi WYŻEJ niż ten plik: jest
//! narzędziem wykonawcy i środowiska serwerowego, gdzie nastawę wnosi jednostka
//! usługi, a nie okno. Kolejność warstw rozstrzyga `ustawienia.rs`.
//!
//! Plik leży w katalogu danych rdzenia (`rdzen::dziennik::katalog_danych`), bo
//! to jedyny katalog, który powłoka i rdzeń już dzielą; drugie miejsce zapisu
//! byłoby drugim stanem do pogodzenia przy przenoszeniu profilu.

use std::fs::{create_dir_all, rename, File};
use std::io::Write;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};

use crate::rdzen::dziennik;

/// Nazwa pliku nastaw powłoki w katalogu danych.
pub const NAZWA_PLIKU: &str = "powloka-nastawy.json";

/// Rozszerzenie pliku przejściowego zapisu.
const ROZSZERZENIE_PRZEJSCIOWE: &str = "nowy";

/// Nastawy zapisane trwale. Brak pola znaczy „Operator nie wskazał", a nie
/// „wskazał wartość domyślną" — dlatego oba pola są opcjonalne, a nie wypełnione
/// domyślnymi liczbami przy zapisie.
#[derive(Clone, Debug, Default, Deserialize, Serialize)]
pub struct Nastawy {
    /// Host, na którym stoi rdzeń — nazwa serwera albo `127.0.0.1` dla rdzenia
    /// na tym urządzeniu.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub host_rdzenia: Option<String>,
    /// Port nasłuchu rdzenia, gdy Operator wskazał inny niż domyślny.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub port_rdzenia: Option<u16>,
}

impl Nastawy {
    /// Czy nastawy niosą wskazanie Operatora. Plik pusty (`{}`) znaczy brak
    /// wskazania — tak samo jak brak pliku.
    pub fn wskazanie_zlozone(&self) -> bool {
        self.host_rdzenia.is_some()
    }
}

/// Ścieżka pliku nastaw powłoki.
pub fn sciezka() -> PathBuf {
    dziennik::katalog_danych().join(NAZWA_PLIKU)
}

/// Czyta nastawy z pliku. Brak pliku, plik nieczytelny i treść niezgodna
/// z umową dają nastawy puste — odczyt nie jest bramą i nie wstrzymuje startu
/// okna. Powód nieczytelności trafia do dziennika powłoki, bo inaczej Operator
/// zobaczyłby ekran pierwszego uruchomienia bez wyjaśnienia, dlaczego jego
/// poprzednie wskazanie zniknęło.
pub fn czytaj() -> Nastawy {
    let sciezka = sciezka();
    let tresc = match std::fs::read_to_string(&sciezka) {
        Ok(tresc) => tresc,
        Err(blad) if blad.kind() == std::io::ErrorKind::NotFound => return Nastawy::default(),
        Err(blad) => {
            dziennik::dopisz(&format!(
                "nastawy: nie można odczytać {}: {blad}",
                sciezka.display()
            ));
            return Nastawy::default();
        }
    };
    match serde_json::from_str::<Nastawy>(&tresc) {
        Ok(nastawy) => nastawy,
        Err(blad) => {
            dziennik::dopisz(&format!(
                "nastawy: treść {} nie jest zgodna z umową ({blad}) — wskazanie uznane za niezłożone",
                sciezka.display()
            ));
            Nastawy::default()
        }
    }
}

/// Zapisuje nastawy trwale i zwraca zdanie o niepowodzeniu, gdy zapis się nie
/// udał.
///
/// Zapis idzie przez plik przejściowy i przemianowanie, żeby przerwanie w trakcie
/// pisania nie zostawiło pliku obciętego — wskazanie odczytane w połowie byłoby
/// gorsze niż wskazanie nieodczytane, bo okno łączyłoby się z adresem złożonym
/// z połowy nazwy hosta.
///
/// W przeciwieństwie do odczytu, niepowodzenie zapisu JEST bramą: wywołujący ma
/// odmówić Operatorowi, a nie przyjąć wskazanie, które zniknie przy następnym
/// starcie.
pub fn zapisz(nastawy: &Nastawy) -> Result<(), String> {
    let sciezka = sciezka();
    if let Some(katalog) = sciezka.parent() {
        create_dir_all(katalog).map_err(|blad| {
            format!(
                "Nie można utworzyć katalogu nastaw {}: {blad}",
                katalog.display()
            )
        })?;
    }
    let tresc = serde_json::to_string_pretty(nastawy)
        .map_err(|blad| format!("Nie można złożyć treści nastaw: {blad}"))?;

    let przejsciowy = sciezka.with_extension(ROZSZERZENIE_PRZEJSCIOWE);
    let mut plik = File::create(&przejsciowy).map_err(|blad| {
        format!(
            "Nie można zapisać nastaw powłoki do {}: {blad}",
            przejsciowy.display()
        )
    })?;
    plik.write_all(tresc.as_bytes())
        .and_then(|()| plik.sync_all())
        .map_err(|blad| {
            format!(
                "Nie można zapisać nastaw powłoki do {}: {blad}",
                przejsciowy.display()
            )
        })?;
    rename(&przejsciowy, &sciezka).map_err(|blad| {
        format!(
            "Nie można ustanowić pliku nastaw {}: {blad}",
            sciezka.display()
        )
    })
}
