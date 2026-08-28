//! Moduł rozpoznaje, czy rdzeń odpowiada pod wskazanym adresem, pytając gniazdo sieciowe
//! zamiast zgadywać stan z pliku procesu.

use std::net::{TcpStream, ToSocketAddrs};
use std::time::Duration;

/// Czas oczekiwania na pojedynczą próbę połączenia z gniazdem rdzenia przy rozpoznawaniu jego bieżącego stanu.
const CZAS_PROBY: Duration = Duration::from_millis(250);

/// Zwraca prawdę, gdy pod wskazanym serwerem i portem ktoś przyjmuje połączenia, dając
/// fałsz przy nazwie nierozwiązywalnej.
pub fn odpowiada_pod(host: &str, port: u16) -> bool {
    match (host, port).to_socket_addrs() {
        Ok(mut adresy) => adresy.any(|adres| TcpStream::connect_timeout(&adres, CZAS_PROBY).is_ok()),
        Err(_) => false,
    }
}
