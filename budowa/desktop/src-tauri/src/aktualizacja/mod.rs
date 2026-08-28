//! Aktualizacja zdalna powłoki: pobranie wydania, sprawdzenie sumy, założenie
//! i ponowny start na urządzeniu operatora, drogą HTTPS i sumy SHA-256
//! z wykazu wydań, a nie wtyczką aktualizacji Tauri.

pub mod droga;
pub mod pobranie;
#[cfg(test)]
pub mod probne;

use std::sync::atomic::{AtomicBool, Ordering};
use std::time::Duration;

use serde::Serialize;
use tauri::AppHandle;

use crate::dziennik;

/// Ile czasu dostaje interfejs na odebranie odpowiedzi, zanim powłoka zniknie;
/// restart natychmiastowy zabiłby okno, zanim baner zdążyłby pokazać
/// powodzenie, i wyglądałoby to jak awaria.
const ZWLOKA_RESTARTU: Duration = Duration::from_millis(1500);

/// Odmowa wykonania aktualizacji opisuje brak, nigdy zakaz: niesie kod
/// przyczyny dla rozgałęzienia w interfejsie oraz zdanie po polsku, gotowe
/// do pokazania operatorowi.
#[derive(Clone, Debug, Serialize)]
pub struct Odmowa {
    /// Kod przyczyny, stały i nadający się do rozgałęzienia w interfejsie.
    pub powod: String,
    /// Zdanie po polsku, gotowe do pokazania w banerze.
    pub zdanie: String,
}

impl Odmowa {
    /// Składa odmowę z kodu i zdania.
    pub fn nowa(powod: &str, zdanie: String) -> Self {
        Self {
            powod: powod.to_string(),
            zdanie,
        }
    }
}

/// Przebieg udanej aktualizacji — dane, jakie interfejs dostaje tuż przed
/// restartem powłoki, wraz z drogą założenia i sumą kontrolną pobranego pliku.
#[derive(Clone, Debug, Serialize)]
pub struct Przebieg {
    /// Droga systemowa, którą poszło założenie: `appimage` albo `instalka-nsis`.
    pub droga: String,
    /// Co zostało podmienione lub uruchomione (ścieżka).
    pub zalozone: String,
    /// Ile bajtów pobrano.
    pub bajtow: u64,
    /// Suma SHA-256 policzona z pobranego pliku — zgodna z żądaną.
    pub suma_sha256: String,
    /// Za ile milisekund aplikacja zniknie z ekranu.
    pub restart_za_ms: u64,
    /// Zdanie dla Operatora.
    pub zdanie: String,
}

/// Czy aktualizacja właśnie trwa — zapora wejścia przeciw dwóm równoległym
/// pobraniom do tego samego pliku roboczego; odmowa stąd opisuje brak, nie
/// zakaz, i mija sama, gdy pierwszy przebieg się skończy.
static W_TOKU: AtomicBool = AtomicBool::new(false);

/// Zwalnia zaporę także wtedy, gdy przebieg wyszedł błędem albo paniką —
/// inaczej jedno niepowodzenie zamykałoby aktualizacje do końca życia procesu.
pub struct StrazWylacznosci;

impl Drop for StrazWylacznosci {
    fn drop(&mut self) {
        W_TOKU.store(false, Ordering::Release);
    }
}

/// Zajmuje wyłączność na aktualizację albo odmawia, bo już trwa. Zwróconą
/// straż trzeba związać z nazwą na cały czas przebiegu — wartość upuszczona
/// od razu zwalnia zaporę i przywraca usterkę, przed którą ta funkcja stoi.
#[must_use = "straż zwalnia zaporę w chwili upuszczenia — zwiąż ją z nazwą na czas całego przebiegu"]
pub fn zajmij_wylacznosc() -> Result<StrazWylacznosci, Odmowa> {
    if W_TOKU
        .compare_exchange(false, true, Ordering::AcqRel, Ordering::Acquire)
        .is_err()
    {
        return Err(Odmowa::nowa(
            "aktualizacja-w-toku",
            "Aktualizacja już trwa — powłoka pobiera wydanie i za chwilę je założy. \
             Drugiego pobrania nie zaczyna, bo oba pisałyby w ten sam plik roboczy \
             i podmieniały aplikację jeden przez drugiego."
                .to_string(),
        ));
    }
    Ok(StrazWylacznosci)
}

/// Wykonuje aktualizację od początku do końca. Wołane przez polecenie IPC.
///
/// Kolejność jest warunkiem poprawności i nie wolno jej mieszać:
/// rozpoznanie drogi → pobranie → **sprawdzenie sumy** → założenie →
/// ponowny start.
pub fn wykonaj(aplikacja: &AppHandle, adres: &str, suma_sha256: &str) -> Result<Przebieg, Odmowa> {
    dziennik::dopisz(&format!("aktualizacja: żądanie wydania spod {adres}"));

    // Jedna aktualizacja naraz: zapora wyłączności nie dopuszcza drugiego przebiegu.
    let straz = zajmij_wylacznosc().inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;

    // Straż idzie dalej pożyczką, nie luźnym wiązaniem, żeby uniknąć jej wcześniejszej utraty.
    przebieg_pod_straza(&straz, aplikacja, adres, suma_sha256)
}

/// Właściwy przebieg aktualizacji, wykonywany pod zajętą wyłącznością; parametr
/// straży, choć nieużywany, jest dowodem sprawdzanym przez kompilator, że nikt
/// nie wywoła tego przebiegu bez zajętej zapory.
fn przebieg_pod_straza(
    _straz: &StrazWylacznosci,
    aplikacja: &AppHandle,
    adres: &str,
    suma_sha256: &str,
) -> Result<Przebieg, Odmowa> {
    // Rozpoznanie drogi sprawdza, czy jest co podmieniać, zanim ściągnie się wydanie.
    let droga = droga::rozpoznaj().inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;
    let plik_roboczy = droga.plik_roboczy();

    // Pobranie wraz ze sprawdzeniem sumy; plik niezgodny nie wraca stąd nigdy.
    let pobrany =
        pobranie::pobierz_i_sprawdz(adres, suma_sha256, &plik_roboczy).inspect_err(|odmowa| {
            dziennik::dopisz(&format!(
                "aktualizacja: odmowa ({}) {}",
                odmowa.powod, odmowa.zdanie
            ));
        })?;
    dziennik::dopisz(&format!(
        "aktualizacja: pobrano {} B, suma {} zgodna z wykazem",
        pobrany.bajtow, pobrany.suma_sha256
    ));

    // Założenie. Dopiero tutaj cokolwiek na dysku Operatora się zmienia.
    let zalozenie = droga.zaloz(&pobrany.sciezka).inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;
    dziennik::dopisz(&format!("aktualizacja: {}", zalozenie.zdanie));

    // Ponowny start — po zwłoce, żeby odpowiedź zdążyła dojść do banera.
    zaplanuj_ponowny_start(aplikacja, zalozenie.restartuje_powloka);

    Ok(Przebieg {
        droga: droga.nazwa().to_string(),
        zalozone: zalozenie.zalozone,
        bajtow: pobrany.bajtow,
        suma_sha256: pobrany.suma_sha256,
        restart_za_ms: ZWLOKA_RESTARTU.as_millis() as u64,
        zdanie: zalozenie.zdanie,
    })
}

/// Odkłada ponowny start powłoki o czas zwłoki restartu; obsługa obrazu
/// przenośnego działa sama i wstaje po podmianie w wydaniu nowym, a przy
/// instalatorze zewnętrznym powłoka wyłącznie schodzi z drogi.
fn zaplanuj_ponowny_start(aplikacja: &AppHandle, restartuje_powloka: bool) {
    let aplikacja = aplikacja.clone();
    std::thread::spawn(move || {
        std::thread::sleep(ZWLOKA_RESTARTU);
        // Wątek roboczy nie może sam kończyć aplikacji — wymaga do tego wątku głównego.
        let _ = aplikacja.clone().run_on_main_thread(move || {
            if restartuje_powloka {
                dziennik::dopisz("aktualizacja: ponowny start powłoki");
                aplikacja.restart();
            } else {
                dziennik::dopisz("aktualizacja: powłoka schodzi instalatorowi z drogi");
                crate::zamkniecie::zakoncz_powloke(&aplikacja);
            }
        });
    });
}
