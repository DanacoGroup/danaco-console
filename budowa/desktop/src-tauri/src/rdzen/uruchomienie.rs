//! Uruchomienie rdzenia w tle przy starcie powłoki (rola `all`).
//!
//! Rdzeń dostaje własną grupę procesów i pracuje bez okna konsoli. Nie jest
//! wiązany z cyklem życia powłoki: przeżywa zamknięcie okna, a nawet
//! zakończenie samej powłoki.
//!
//! Fail-open: brak binarki, odmowa startu albo minięcie czasu
//! oczekiwania nie wstrzymuje otwarcia okna — kończy wyłącznie tę czynność
//! i trafia do opisu stanu.
//!
//! Trzy warianty startu odpowiadają trzem stanom wskazania (`ustawienia.rs`):
//! rdzeń na tym urządzeniu (`zawiaz_lokalnie`), rdzeń na serwerze
//! (`zawiaz_zdalnie`) i wskazanie jeszcze niezłożone (`zawiaz_oczekiwanie`).
//! Każdy z nich składa `Zawiazanie` — komplet, którym da się zarówno utworzyć
//! uchwyt przy starcie, jak i podmienić go po złożeniu wskazania w oknie
//! (`UchwytRdzenia::przejmij`, wołane z `wskazanie.rs`).

use std::process::{Child, Command};
use std::time::Duration;

use crate::ustawienia::Ustawienia;

use super::dziennik;
use super::lokalizacja::{self, Wynik};
use super::nasluch;
use super::pakiet_klienta;
use super::uchwyt::{OpisRdzenia, UchwytRdzenia};

/// Rola procesu rdzenia stawianego przez powłokę — hub i agent lokalny w jednym
/// procesie (`server/internal/konfiguracja/rola.go`).
const ROLA: &str = "all";

/// Czas oczekiwania na pojawienie się nasłuchu rdzenia.
const CZAS_NA_NASLUCH: Duration = Duration::from_secs(10);

/// Komplet stanu rdzenia gotowy do wstawienia w uchwyt.
pub struct Zawiazanie {
    pub port: u16,
    pub host: String,
    pub opis: OpisRdzenia,
    pub dziecko: Option<Child>,
}

impl Zawiazanie {
    /// Zamienia zawiązanie w uchwyt stanu aplikacji.
    pub fn w_uchwyt(self) -> UchwytRdzenia {
        UchwytRdzenia::nowy(self.port, self.host, self.opis, self.dziecko)
    }
}

/// Stawia rdzeń w tle i zwraca uchwyt do stanu aplikacji.
pub fn uruchom_w_tle(ustawienia: &Ustawienia) -> UchwytRdzenia {
    zawiaz_lokalnie(ustawienia).w_uchwyt()
}

/// Buduje uchwyt rdzenia bez stawiania procesu lokalnego — wariant wirtualny.
///
/// Wołane z `main.rs` zamiast `uruchom_w_tle`, gdy `Ustawienia::rdzen_lokalny`
/// zwraca fałsz: Operator wskazał host inny niż domyślny, więc rdzeń stoi na
/// innej maszynie. Stawianie tu drugiego procesu byłoby zbędne i szkodliwe —
/// tworzyłoby drugi stan (dwa procesy rdzenia, jeden nieużywany). Powłoka nie ma
/// czym takiego rdzenia postawić ani zatrzymać — `UchwytRdzenia::zatrzymaj`
/// odpowie, że nie ma czego zatrzymywać, tak samo jak przy rdzeniu zastanym
/// pracującym.
pub fn nie_stawiaj_lokalnie(ustawienia: &Ustawienia) -> UchwytRdzenia {
    zawiaz_zdalnie(ustawienia).w_uchwyt()
}

/// Buduje uchwyt dla powłoki, która czeka na wskazanie Operatora.
///
/// Wołane z `main.rs`, gdy wskazania nie ma w żadnej warstwie
/// (`Ustawienia::wskazanie_zlozone` zwraca fałsz) — czyli przy pierwszym
/// uruchomieniu po instalacji. Powłoka NIE stawia wtedy rdzenia lokalnego: nie
/// wie jeszcze, czy Operator pracuje z rdzeniem na tym urządzeniu, czy z rdzeniem
/// na serwerze, a proces postawiony na wyrost trzeba by zaraz ubijać i tłumaczyć
/// Operatorowi, skąd wziął się rdzeń, o który nie prosił.
pub fn oczekuj_na_wskazanie(ustawienia: &Ustawienia) -> UchwytRdzenia {
    zawiaz_oczekiwanie(ustawienia).w_uchwyt()
}

/// Zawiązanie dla wskazania jeszcze niezłożonego. Nic nie stawia i nic nie pyta:
/// stan jest znany bez pytania, bo rdzenia pod tym wskazaniem po prostu nie ma.
pub fn zawiaz_oczekiwanie(ustawienia: &Ustawienia) -> Zawiazanie {
    Zawiazanie {
        port: ustawienia.port(),
        host: ustawienia.host(),
        opis: opis(
            false,
            None,
            false,
            ustawienia.adres_rdzenia_http(),
            dziennik::sciezka().display().to_string(),
            "Wskazanie rdzenia nie zostało jeszcze złożone — powłoka czeka, aż Operator wskaże, \
             czy rdzeń stoi na tym urządzeniu, czy na serwerze."
                .to_string(),
        ),
        dziecko: None,
    }
}

/// Zawiązanie dla rdzenia stojącego na innej maszynie.
///
/// Rozpoznanie nasłuchu (`UchwytRdzenia::opis`) pyta host obowiązujący, nie
/// zawsze pętlę zwrotną: uchwyt niesie `ustawienia.host()`, a `opis` woła
/// `nasluch::odpowiada_pod(host, port)`. W wariancie wirtualnym pole `pracuje`
/// odbija więc stan rdzenia zdalnego, nie lokalnego.
pub fn zawiaz_zdalnie(ustawienia: &Ustawienia) -> Zawiazanie {
    let adres = ustawienia.adres_rdzenia_http();
    Zawiazanie {
        port: ustawienia.port(),
        host: ustawienia.host(),
        opis: opis(
            false,
            None,
            false,
            adres.clone(),
            dziennik::sciezka().display().to_string(),
            format!(
                "Wariant wirtualny: rdzeń wskazany pod {adres} — powłoka nie stawia procesu lokalnego."
            ),
        ),
        dziecko: None,
    }
}

/// Zawiązanie dla rdzenia na tym urządzeniu: rozpoznaje rdzeń zastany, a gdy go
/// nie ma — stawia proces potomny.
///
/// Droga unikania drugiego procesu pyta pętlę zwrotną, bo rdzeń zdalny nie jest
/// tu stawiany.
pub fn zawiaz_lokalnie(ustawienia: &Ustawienia) -> Zawiazanie {
    let adres = ustawienia.adres_rdzenia_http();
    let sciezka_dziennika = dziennik::sciezka().display().to_string();
    let port = ustawienia.port();

    if nasluch::odpowiada(port) {
        return Zawiazanie {
            port,
            host: ustawienia.host(),
            opis: opis(
                true,
                None,
                false,
                adres,
                sciezka_dziennika,
                format!("Rdzeń zastany na porcie {port} — powłoka nie stawia drugiego procesu."),
            ),
            dziecko: None,
        };
    }

    let sciezka = match lokalizacja::znajdz(ustawienia.sciezka_rdzenia.as_deref()) {
        Wynik::Znaleziona(sciezka) => sciezka,
        Wynik::Brak(miejsca) => {
            return Zawiazanie {
                port,
                host: ustawienia.host(),
                opis: opis(
                    false,
                    None,
                    false,
                    adres,
                    sciezka_dziennika,
                    format!(
                        "Nie znaleziono binarki rdzenia. Szukano: {}. Wskaż plik zmienną {}.",
                        lokalizacja::opis_miejsc(&miejsca),
                        crate::ustawienia::ZMIENNA_RDZEN
                    ),
                ),
                dziecko: None,
            };
        }
    };

    match postaw(&sciezka, ustawienia) {
        Ok(proces) => {
            let pid = proces.id();
            let czekanie = nasluch::czekaj_na_nasluch(port, CZAS_NA_NASLUCH);
            let tresc = match czekanie {
                Some(czas) => format!(
                    "Rdzeń uruchomiony z {} (pid {pid}, rola {ROLA}); nasłuch po {} ms.",
                    sciezka.display(),
                    czas.as_millis()
                ),
                None => format!(
                    "Rdzeń uruchomiony z {} (pid {pid}), lecz nie odpowiedział na porcie {} w ciągu {} s — dziennik: {}.",
                    sciezka.display(), port, CZAS_NA_NASLUCH.as_secs(), sciezka_dziennika
                ),
            };
            Zawiazanie {
                port,
                host: ustawienia.host(),
                opis: opis(
                    czekanie.is_some(),
                    Some(pid),
                    true,
                    adres,
                    sciezka_dziennika,
                    tresc,
                ),
                dziecko: Some(proces),
            }
        }
        Err(blad) => Zawiazanie {
            port,
            host: ustawienia.host(),
            opis: opis(
                false,
                None,
                false,
                adres,
                sciezka_dziennika,
                format!(
                    "Nie udało się uruchomić rdzenia z {}: {blad}",
                    sciezka.display()
                ),
            ),
            dziecko: None,
        },
    }
}

/// Buduje wywołanie rdzenia i stawia proces potomny.
fn postaw(sciezka: &std::path::Path, ustawienia: &Ustawienia) -> std::io::Result<Child> {
    let (wyjscie, blad) = dziennik::strumienie();
    let mut polecenie = Command::new(sciezka);
    polecenie
        .arg("--role")
        .arg(ROLA)
        .arg("--port")
        .arg(ustawienia.port().to_string())
        .stdin(std::process::Stdio::null())
        .stdout(wyjscie)
        .stderr(blad);
    if let Some(katalog) = sciezka.parent() {
        polecenie.current_dir(katalog);
    }
    if let Some(pakiet) = pakiet_klienta::znajdz(sciezka) {
        polecenie.env(pakiet_klienta::ZMIENNA, pakiet);
    }
    odetnij_od_konsoli(&mut polecenie);
    polecenie.spawn()
}

/// Odcina rdzeń od konsoli i od grupy procesów powłoki, aby przerwanie powłoki
/// nie przerwało rdzenia (Windows: CREATE_NO_WINDOW | CREATE_NEW_PROCESS_GROUP).
#[cfg(windows)]
fn odetnij_od_konsoli(polecenie: &mut Command) {
    use std::os::windows::process::CommandExt;
    const CREATE_NO_WINDOW: u32 = 0x0800_0000;
    const CREATE_NEW_PROCESS_GROUP: u32 = 0x0000_0200;
    polecenie.creation_flags(CREATE_NO_WINDOW | CREATE_NEW_PROCESS_GROUP);
}

/// Poza Windows nie ma odpowiednika chorągiewek tworzenia procesu.
#[cfg(not(windows))]
fn odetnij_od_konsoli(_polecenie: &mut Command) {}

/// Składa opis stanu rdzenia.
fn opis(
    pracuje: bool,
    pid: Option<u32>,
    postawiony_przez_powloke: bool,
    adres: String,
    dziennik: String,
    tresc: String,
) -> OpisRdzenia {
    OpisRdzenia {
        pracuje,
        pid,
        postawiony_przez_powloke,
        adres,
        dziennik,
        opis: tresc,
    }
}

// `zawiaz_lokalnie` i `postaw` uruchamiają procesy i pytają gniazdo sieciowe,
// więc pozostają bez testów jednostkowych. `opis` jest jedyną czystą częścią
// tego modułu i to ona ma pokrycie.
