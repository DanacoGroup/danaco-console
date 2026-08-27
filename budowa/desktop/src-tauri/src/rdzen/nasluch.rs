//! Rozpoznanie, czy rdzeń odpowiada pod wskazanym adresem.
//!
//! Powłoka nie zgaduje stanu rdzenia z pliku PID ani z nazwy procesu — pyta
//! gniazdo. Rdzeń stoi na serwerze wdrożenia, a nie na urządzeniu Operatora,
//! więc pytanie idzie pod wskazany serwer, nigdy pod pętlę zwrotną: pętla
//! zwrotna dawałaby fałsz niezależnie od stanu rdzenia.

use std::net::{TcpStream, ToSocketAddrs};
use std::time::Duration;

/// Czas oczekiwania na pojedynczą próbę połączenia.
const CZAS_PROBY: Duration = Duration::from_millis(250);

/// Zwraca prawdę, gdy pod wskazanym serwerem i portem ktoś przyjmuje połączenia.
///
/// Nazwa nierozwiązywalna daje fałsz, nie panikę — rozpoznanie niczego nie
/// wstrzymuje, jest wyłącznie odpowiedzią na pytanie o stan.
pub fn odpowiada_pod(host: &str, port: u16) -> bool {
    match (host, port).to_socket_addrs() {
        Ok(mut adresy) => adresy.any(|adres| TcpStream::connect_timeout(&adres, CZAS_PROBY).is_ok()),
        Err(_) => false,
    }
}
