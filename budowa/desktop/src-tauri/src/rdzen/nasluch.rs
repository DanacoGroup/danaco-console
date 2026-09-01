//! Moduł rozpoznaje, czy rdzeń odpowiada pod wskazanym adresem, pytając gniazdo sieciowe
//! zamiast zgadywać stan z pliku procesu.

use std::net::{TcpStream, ToSocketAddrs};
use std::time::Duration;

use crate::ustawienia::{gospodarz_w_adresie, Schemat};

/// Czas oczekiwania na pojedynczą próbę połączenia z gniazdem rdzenia przy rozpoznawaniu jego bieżącego stanu.
const CZAS_PROBY: Duration = Duration::from_millis(250);

/// Czas na całą próbę łącza szyfrowanego: rozwiązanie nazwy, połączenie,
/// uzgodnienie TLS i odpowiedź serwera. Wielokrotność `CZAS_PROBY`, bo samo
/// uzgodnienie TLS to dwa obroty po łączu, a próba idzie do serwera wdrożenia,
/// nie do pętli zwrotnej.
const CZAS_PROBY_SZYFROWANEJ: Duration = Duration::from_secs(2);

/// Zwraca prawdę, gdy pod wskazanym serwerem i portem odpowiada rdzeń, dając
/// fałsz przy nazwie nierozwiązywalnej. Schemat podaje wołający: przy wskazaniu
/// z okna jest nim schemat z adresu, jeszcze niezapisany w nastawach.
pub fn odpowiada_pod(schemat: Schemat, host: &str, port: u16) -> bool {
    match schemat {
        Schemat::Http => przyjmuje_polaczenie(host, port),
        Schemat::Https => uzgadnia_tls(host, port),
    }
}

/// Sprawdza, czy pod adresem ktoś przyjmuje połączenia warstwy transportowej.
fn przyjmuje_polaczenie(host: &str, port: u16) -> bool {
    match (host, port).to_socket_addrs() {
        Ok(mut adresy) => {
            adresy.any(|adres| TcpStream::connect_timeout(&adres, CZAS_PROBY).is_ok())
        }
        Err(_) => false,
    }
}

/// Sprawdza łączność szyfrowaną pełnym żądaniem HTTPS. Samo przyjęcie połączenia
/// na porcie nie dowodzi, że stoi tam TLS z certyfikatem, któremu okno zaufa —
/// a to właśnie rozstrzyga o tym, czy strona zestawi `wss://`. Odpowiedź o kodzie
/// innym niż powodzenie też jest odpowiedzią: pytanie brzmi „czy rdzeń odbiera”,
/// nie „co odpowiada pod ukośnikiem”.
fn uzgadnia_tls(host: &str, port: u16) -> bool {
    let klient: ureq::Agent = ureq::Agent::config_builder()
        .http_status_as_error(false)
        .timeout_global(Some(CZAS_PROBY_SZYFROWANEJ))
        .build()
        .into();
    let adres = format!(
        "{}://{}:{port}/",
        Schemat::Https.nazwa(),
        gospodarz_w_adresie(host)
    );
    klient.get(&adres).call().is_ok()
}
