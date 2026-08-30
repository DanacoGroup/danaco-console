//! Krok 5 kreatora: pobranie składników programu z serwera wdrożenia Danaco.
//! Adres i sumy kontrolne pochodzą z wykazu wydań `budowa/witryna/wydania.json`
//! (wkompilowanego przy budowie, jak w powłoce głównej) — nic tu nie jest
//! zgadywane ani wpisane na sztywno drugi raz.

use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::io::Read;
use std::path::{Path, PathBuf};
use tauri::{AppHandle, Emitter};

/// Wykaz wydań, jedyne źródło prawdy o adresie kanału i pozycjach do pobrania —
/// ten sam plik, z którego korzysta powłoka główna (`aktualizacja/pobranie.rs`).
const WYKAZ_WYDAN: &str = include_str!("../../../witryna/wydania.json");

/// Poświadczenia kanału pobrań wpisane w postać instalki przy jej składaniu.
/// Katalog `/wydania/` chroni uwierzytelnienie podstawowe, a kreator nie ma
/// gdzie o nie zapytać: żaden z sześciu kroków nie przewiduje takiego pola.
/// Hasła nie ma w repozytorium — wchodzi zmienną środowiska w chwili budowy,
/// tak samo jak adres serwera wdrożenia w powłoce.
const UZYTKOWNIK_KANALU: Option<&str> = option_env!("DANACO_KANAL_UZYTKOWNIK");
const HASLO_KANALU: Option<&str> = option_env!("DANACO_KANAL_HASLO");

/// Składa nagłówek uwierzytelnienia podstawowego, gdy poświadczenia wpisano
/// przy składaniu. Bez nich żądanie idzie tak jak dotąd — i wraca odmową 401
/// nazwaną wprost, zamiast cichego niepowodzenia.
fn naglowek_poswiadczen() -> Option<String> {
    let uzytkownik = UZYTKOWNIK_KANALU?;
    let haslo = HASLO_KANALU?;
    let para = format!("{uzytkownik}:{haslo}");
    Some(format!("Basic {}", base64_podstawowy(para.as_bytes())))
}

/// Zapis base64 bez zależności zewnętrznej: kanał wymaga jednego nagłówka,
/// a dokładanie dla niego biblioteki byłoby kosztem bez pokrycia.
fn base64_podstawowy(dane: &[u8]) -> String {
    const ZNAKI: &[u8; 64] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut wynik = String::with_capacity(dane.len().div_ceil(3) * 4);
    for porcja in dane.chunks(3) {
        let b = [porcja[0], *porcja.get(1).unwrap_or(&0), *porcja.get(2).unwrap_or(&0)];
        let trojka = u32::from(b[0]) << 16 | u32::from(b[1]) << 8 | u32::from(b[2]);
        for i in 0..4 {
            if i <= porcja.len() {
                wynik.push(ZNAKI[(trojka >> (18 - i * 6) & 0x3F) as usize] as char);
            } else {
                wynik.push('=');
            }
        }
    }
    wynik
}

/// Zdarzenie niosące postęp pobierania do okna kreatora.
pub const ZDARZENIE_POSTEP: &str = "instalator:postep-pobrania";

#[derive(Deserialize)]
struct Wykaz {
    kanal: Kanal,
    wydania: Vec<PozycjaWydania>,
}

#[derive(Deserialize)]
struct Kanal {
    adres: String,
}

#[derive(Deserialize, Clone)]
pub(crate) struct PozycjaWydania {
    system: String,
    postac: String,
    architektura: String,
    plik: String,
    #[serde(rename = "nazwaPliku")]
    nazwa_pliku: String,
    #[serde(rename = "rozmiarBajty")]
    rozmiar_bajty: u64,
    suma: String,
}

/// Odmowa niesie powód po angielsku ustalony z zamkniętego słownika i zdanie
/// gotowe do pokazania Operatorowi — bez wymyślonego kodu liczbowego.
#[derive(Serialize, Clone, Debug)]
pub struct Odmowa {
    pub powod: String,
    pub zdanie: String,
}

impl Odmowa {
    fn nowa(powod: &str, zdanie: String) -> Self {
        Self {
            powod: powod.to_string(),
            zdanie,
        }
    }
}

#[derive(Serialize, Clone)]
pub struct PostepPobrania {
    pub odebrano_bajtow: u64,
    pub razem_bajtow: u64,
}

#[derive(Serialize, Clone)]
pub struct WynikPobrania {
    pub adres: String,
    pub nazwa_pliku: String,
    pub bajtow: u64,
    pub suma_sha256: String,
}

/// Wybiera pozycję wykazu zgodną z architekturą wykrytą na tej maszynie
/// (`x64` → x64 wykazu, `arm` → ARM64 wykazu). Architektura nierozpoznana
/// (`brak`) nie ma czego dopasować — to nazwana odmowa, nie zgadywanie.
fn dobierz_pozycje(architektura_kreatora: &str) -> Result<PozycjaWydania, Odmowa> {
    let wykaz: Wykaz = serde_json::from_str(WYKAZ_WYDAN).map_err(|blad| {
        Odmowa::nowa(
            "wykaz-nieczytelny",
            format!("Wykaz wydań budowa/witryna/wydania.json nie jest poprawnym JSON-em: {blad}."),
        )
    })?;

    let architektura_wykazu = match architektura_kreatora {
        "x64" => "x64",
        "arm" => "ARM64",
        _ => {
            return Err(Odmowa::nowa(
                "architektura-nierozpoznana",
                "Architektura procesora tej maszyny nie została rozpoznana — wykaz \
                 wydań nie ma z czym jej dopasować."
                    .to_string(),
            ))
        }
    };

    wykaz
        .wydania
        .iter()
        // Instalator ściąga POWŁOKĘ programu, nie siebie: wykaz niesie obie
        // postaci pod tym samym systemem i architekturą, więc pozycję
        // rozstrzyga `postac`, a nie kolejność w wykazie.
        .find(|w| {
            w.system == "Windows"
                && w.architektura == architektura_wykazu
                && w.postac == "hybryda"
        })
        .cloned()
        .ok_or_else(|| {
            Odmowa::nowa(
                "wydanie-nieopublikowane",
                format!(
                    "Wykaz wydań nie niesie dziś powłoki Windows/{architektura_wykazu} — \
                     nie ma czego pobrać dla tej maszyny."
                ),
            )
        })
        .inspect(|pozycja| {
            // Adres kanału pozycji musi zgadzać się z `kanal.adres` wykazu —
            // rozjazd byłby usterką samego wykazu, nie tego kroku.
            debug_assert!(
                pozycja
                    .plik
                    .starts_with(wykaz.kanal.adres.trim_end_matches('/')),
                "wykaz wydań: pozycja {} leży poza kanałem {}",
                pozycja.plik,
                wykaz.kanal.adres
            );
        })
}

/// Sprawdza, bez pobierania pliku, czy serwer wydań odda pozycję dobraną do
/// architektury tej maszyny — jedno realne żądanie HTTPS do rzeczywistego
/// adresu z wykazu wydań. Wołane przed otwarciem okna, żeby krok 5 wiedział
/// z góry, czy ma pokazać przebieg pobrania, czy od razu nazwaną odmowę.
pub fn sprawdz_wstepnie(architektura: &str) -> Result<PozycjaWydania, Odmowa> {
    let pozycja = dobierz_pozycje(architektura)?;

    let klient: ureq::Agent = ureq::Agent::config_builder()
        .http_status_as_error(false)
        .build()
        .into();
    let zadanie = klient.get(&pozycja.plik);
    let zadanie = match naglowek_poswiadczen() {
        Some(naglowek) => zadanie.header("Authorization", &naglowek),
        None => zadanie,
    };
    let odpowiedz = zadanie.call().map_err(|blad| {
        Odmowa::nowa(
            "brak-lacznosci",
            format!("Nie udało się połączyć z {}: {blad}.", pozycja.plik),
        )
    })?;

    let status = odpowiedz.status();
    if !status.is_success() {
        let www_authenticate = odpowiedz
            .headers()
            .get("www-authenticate")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("");
        return Err(Odmowa::nowa(
            "odpowiedz-serwera",
            format!(
                "Serwer wydań {} odpowiedział kodem {status}{}. Kanał pobrań \
                 pozycji wydania chroni uwierzytelnienie podstawowe (plik \
                 poświadczeń poza repozytorium) — ten instalator nie ma \
                 poświadczeń dostępu.",
                pozycja.plik,
                if www_authenticate.is_empty() {
                    String::new()
                } else {
                    format!(" ({www_authenticate})")
                }
            ),
        ));
    }
    Ok(pozycja)
}

/// Pobiera plik wydania dobrany do architektury tej maszyny do katalogu roboczego
/// i sprawdza jego sumę SHA-256 z wykazu. Zgłasza postęp zdarzeniami do okna.
#[tauri::command]
pub async fn pobierz_skladniki(
    aplikacja: AppHandle,
    architektura: String,
    katalog_roboczy: String,
) -> Result<WynikPobrania, Odmowa> {
    tauri::async_runtime::spawn_blocking(move || {
        wykonaj(&aplikacja, &architektura, Path::new(&katalog_roboczy))
    })
    .await
    .unwrap_or_else(|blad| {
        Err(Odmowa::nowa(
            "przebieg-przerwany",
            format!("Pobieranie przerwało się nieoczekiwanie: {blad}."),
        ))
    })
}

fn wykonaj(
    aplikacja: &AppHandle,
    architektura: &str,
    katalog_roboczy: &Path,
) -> Result<WynikPobrania, Odmowa> {
    let pozycja = dobierz_pozycje(architektura)?;

    std::fs::create_dir_all(katalog_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "brak-prawa-zapisu",
            format!(
                "Nie udało się przygotować katalogu roboczego {}: {blad}.",
                katalog_roboczy.display()
            ),
        )
    })?;
    let plik_roboczy: PathBuf = katalog_roboczy.join(&pozycja.nazwa_pliku);

    let klient: ureq::Agent = ureq::Agent::config_builder()
        .http_status_as_error(false)
        .build()
        .into();

    let zadanie = klient.get(&pozycja.plik);
    let zadanie = match naglowek_poswiadczen() {
        Some(naglowek) => zadanie.header("Authorization", &naglowek),
        None => zadanie,
    };
    let mut odpowiedz = zadanie.call().map_err(|blad| {
        Odmowa::nowa(
            "brak-lacznosci",
            format!("Nie udało się połączyć z {}: {blad}.", pozycja.plik),
        )
    })?;

    let status = odpowiedz.status();
    if !status.is_success() {
        let www_authenticate = odpowiedz
            .headers()
            .get("www-authenticate")
            .and_then(|v| v.to_str().ok())
            .unwrap_or("");
        return Err(Odmowa::nowa(
            "odpowiedz-serwera",
            format!(
                "Serwer wydań {} odpowiedział kodem {status}{}. Kanał pobrań \
                 pozycji wydania chroni uwierzytelnienie podstawowe (plik \
                 poświadczeń poza repozytorium) — ten instalator nie ma \
                 poświadczeń dostępu.",
                pozycja.plik,
                if www_authenticate.is_empty() {
                    String::new()
                } else {
                    format!(" ({www_authenticate})")
                }
            ),
        ));
    }

    let mut plik = std::fs::File::create(&plik_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "brak-prawa-zapisu",
            format!(
                "Nie udało się utworzyć pliku roboczego {}: {blad}.",
                plik_roboczy.display()
            ),
        )
    })?;

    let mut haszujacy = Sha256::new();
    let mut odebrano: u64 = 0;
    let mut bufor = [0u8; 64 * 1024];
    let mut zrodlo = odpowiedz.body_mut().as_reader();
    loop {
        let n = zrodlo.read(&mut bufor).map_err(|blad| {
            Odmowa::nowa(
                "polaczenie-przerwane",
                format!("Połączenie przerwało się w trakcie pobierania: {blad}."),
            )
        })?;
        if n == 0 {
            break;
        }
        haszujacy.update(&bufor[..n]);
        std::io::Write::write_all(&mut plik, &bufor[..n]).map_err(|blad| {
            Odmowa::nowa(
                "brak-prawa-zapisu",
                format!("Nie udało się zapisać pliku roboczego: {blad}."),
            )
        })?;
        odebrano += n as u64;
        let _ = aplikacja.emit(
            ZDARZENIE_POSTEP,
            PostepPobrania {
                odebrano_bajtow: odebrano,
                razem_bajtow: pozycja.rozmiar_bajty,
            },
        );
    }

    let suma_policzona = format!("{:x}", haszujacy.finalize());
    if suma_policzona != pozycja.suma {
        let _ = std::fs::remove_file(&plik_roboczy);
        return Err(Odmowa::nowa(
            "suma-niezgodna",
            format!(
                "Suma kontrolna się nie zgadza: wykaz wydań zapowiadał {}, a pobrany \
                 plik ma {suma_policzona}. Plik został skasowany.",
                pozycja.suma
            ),
        ));
    }

    Ok(WynikPobrania {
        adres: pozycja.plik,
        nazwa_pliku: pozycja.nazwa_pliku,
        bajtow: odebrano,
        suma_sha256: suma_policzona,
    })
}
