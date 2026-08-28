//! Moduł zapisuje trwale nastawy powłoki: wskazanie Operatora, gdzie stoi serwer
//! wdrożenia z rdzeniem, w pliku katalogu danych powłoki.

use std::fs::{create_dir_all, rename, File};
use std::io::Write;
use std::path::PathBuf;

use serde::{Deserialize, Serialize};

use crate::dziennik;

/// Nazwa pliku nastaw powłoki, zapisywanego trwale w katalogu danych powłoki, tuż obok pliku dziennika powłoki.
pub const NAZWA_PLIKU: &str = "powloka-nastawy.json";

/// Rozszerzenie pliku przejściowego zapisu, zastępowanego docelowym plikiem nastaw dopiero po zakończeniu zapisu.
const ROZSZERZENIE_PRZEJSCIOWE: &str = "nowy";

/// Nastawy zapisane trwale. Brak pola znaczy „Operator nie wskazał", a nie
/// „wskazał wartość domyślną" — dlatego oba pola są opcjonalne, a nie wypełnione
/// domyślnymi liczbami przy zapisie.
#[derive(Clone, Debug, Default, Deserialize, Serialize)]
pub struct Nastawy {
    /// Serwer wdrożenia, na którym stoi rdzeń.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub host_rdzenia: Option<String>,
    /// Port nasłuchu rdzenia, gdy Operator wskazał inny niż domyślny.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub port_rdzenia: Option<u16>,
}

impl Nastawy {
    /// Czy nastawy niosą wskazanie Operatora; plik pusty znaczy brak wskazania, tak jak brak pliku.
    pub fn wskazanie_zlozone(&self) -> bool {
        self.host_rdzenia.is_some()
    }
}

/// Ścieżka pliku nastaw powłoki, złożona z katalogu danych powłoki oraz nazwy pliku nastaw tej powłoki.
pub fn sciezka() -> PathBuf {
    dziennik::katalog_danych().join(NAZWA_PLIKU)
}

/// Czyta nastawy z pliku, oddając nastawy puste przy braku pliku, treści nieczytelnej lub
/// niezgodnej z umową danych.
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

/// Zapisuje nastawy trwale przez plik przejściowy i przemianowanie, zwracając zdanie
/// o niepowodzeniu zapisu.
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
