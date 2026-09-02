//! Moduł stawia rdzeń jako proces poboczny powłoki wtedy, gdy rdzeń ma stać na
//! tej samej maszynie, i przekazuje mu przy starcie sekret nawiązania oraz port
//! nasłuchu. Rdzeń serwera wdrożenia powłoka wyłącznie odpytuje — nie startuje.

use std::path::PathBuf;
use std::process::{Child, Command, Stdio};

use crate::dziennik;
use crate::ustawienia::{Ustawienia, ZMIENNA_PORT, ZMIENNA_SEKRET_NAWIAZANIA};

/// Nazwa pliku wykonywalnego rdzenia obok powłoki. Ustala ją katalog programu
/// rdzenia (`server/cmd/danaco-console`); powłoka szuka go w swoim katalogu,
/// bo instalka wydania zakłada oba pliki w jednym miejscu.
#[cfg(windows)]
const NAZWA_RDZENIA: &str = "danaco-console.exe";
#[cfg(not(windows))]
const NAZWA_RDZENIA: &str = "danaco-console";

/// Gospodarze, dla których rdzeń stoi na maszynie Operatora.
const GOSPODARZE_LOKALNI: [&str; 3] = ["127.0.0.1", "localhost", "::1"];

/// Stawia rdzeń obok powłoki, gdy jest po co: wskazanie prowadzi na tę maszynę,
/// a plik rdzenia stoi w katalogu powłoki. Brak któregoś z tych warunków nie
/// jest usterką — znaczy rdzeń na serwerze wdrożenia, do którego powłoka tylko
/// się łączy. Zwrócone dziecko trzyma wołający, żeby wygasić je przy zamknięciu.
pub fn uruchom_obok(ustawienia: &Ustawienia) -> Option<Child> {
    if !rdzen_na_tej_maszynie(ustawienia) {
        return None;
    }
    let plik = plik_rdzenia()?;
    let sekret = ustawienia.sekret_dla_rdzenia();
    let dziecko = Command::new(&plik)
        .env(ZMIENNA_SEKRET_NAWIAZANIA, &sekret)
        .env(ZMIENNA_PORT, ustawienia.port().to_string())
        .stdin(Stdio::null())
        .stdout(Stdio::null())
        .stderr(Stdio::null())
        .spawn();
    match dziecko {
        Ok(dziecko) => {
            ustawienia.odnotuj_proces_poboczny(sekret);
            dziennik::dopisz(&format!(
                "rdzeń: proces poboczny {} wystartowany z sekretem nawiązania wytworzonym przez powłokę",
                plik.display()
            ));
            Some(dziecko)
        }
        Err(blad) => {
            dziennik::dopisz(&format!(
                "rdzeń: nie udało się uruchomić {}: {blad}",
                plik.display()
            ));
            None
        }
    }
}

/// Wygasza proces poboczny rdzenia przy zamknięciu powłoki; rdzeń zostawiony
/// żywy trzymałby port zajęty do następnego uruchomienia maszyny.
pub fn wygas(dziecko: &mut Child) {
    let _ = dziecko.kill();
    let _ = dziecko.wait();
    dziennik::dopisz("rdzeń: proces poboczny wygaszony wraz z powłoką");
}

/// Czy wskazanie prowadzi na maszynę, na której stoi powłoka.
fn rdzen_na_tej_maszynie(ustawienia: &Ustawienia) -> bool {
    match ustawienia.wskazanie() {
        Some(wskazane) => GOSPODARZE_LOKALNI.contains(&wskazane.host.to_ascii_lowercase().as_str()),
        None => false,
    }
}

/// Ścieżka pliku rdzenia w katalogu powłoki; brak pliku znaczy wydanie bez rdzenia.
fn plik_rdzenia() -> Option<PathBuf> {
    let plik = std::env::current_exe().ok()?.parent()?.join(NAZWA_RDZENIA);
    plik.is_file().then_some(plik)
}
