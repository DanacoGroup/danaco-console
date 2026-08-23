//! Sprawdzenie, czy rdzeń nasłuchuje na porcie lokalnym.
//!
//! Powłoka nie zgaduje stanu rdzenia z pliku PID ani z nazwy procesu —
//! pyta gniazdo. Dzięki temu rdzeń uruchomiony wcześniej z wiersza poleceń
//! jest rozpoznany i powłoka nie stawia drugiego procesu na tym samym porcie.

use std::net::{Ipv4Addr, SocketAddr, TcpStream, ToSocketAddrs};
use std::thread::sleep;
use std::time::{Duration, Instant};

/// Czas oczekiwania na pojedynczą próbę połączenia.
const CZAS_PROBY: Duration = Duration::from_millis(250);
/// Odstęp między próbami przy oczekiwaniu na start rdzenia.
const ODSTEP_PROB: Duration = Duration::from_millis(150);

/// Zwraca prawdę, gdy pod portem lokalnym ktoś przyjmuje połączenia.
///
/// Pytanie celuje w pętlę zwrotną, bo służy wyłącznie unikaniu drugiego
/// procesu: powłoka stawia rdzeń lokalny tylko wtedy, gdy nikt jeszcze nie
/// nasłuchuje. Rdzenia zdalnego powłoka nie stawia, więc nie ma tam czego unikać.
pub fn odpowiada(port: u16) -> bool {
    let adres = SocketAddr::from((Ipv4Addr::LOCALHOST, port));
    TcpStream::connect_timeout(&adres, CZAS_PROBY).is_ok()
}

/// Zwraca prawdę, gdy pod wskazanym hostem i portem ktoś przyjmuje połączenia.
///
/// Dla hosta lokalnego (`HOST_DOMYSLNY`) daje to samo co [`odpowiada`]; dla
/// hosta zdalnego (wariant wirtualny) pyta właściwą maszynę, nie pętlę zwrotną.
/// Bez tego pole `pracuje` w opisie stanu byłoby prawie zawsze fałszywe przy
/// rdzeniu wskazanym pod domeną, niezależnie od jego rzeczywistego stanu.
/// Host nierozwiązywalny daje fałsz, nie panikę — rozpoznanie nic nie wstrzymuje.
pub fn odpowiada_pod(host: &str, port: u16) -> bool {
    match (host, port).to_socket_addrs() {
        Ok(mut adresy) => {
            adresy.any(|adres| TcpStream::connect_timeout(&adres, CZAS_PROBY).is_ok())
        }
        Err(_) => false,
    }
}

/// Czeka na nasłuch rdzenia do wyczerpania limitu. Zwraca czas oczekiwania,
/// gdy rdzeń odpowiedział, albo `None`, gdy limit minął.
///
/// Limit nie jest bramą: jego przekroczenie nie wstrzymuje okna, a jedynie
/// trafia do opisu stanu rdzenia.
pub fn czekaj_na_nasluch(port: u16, limit: Duration) -> Option<Duration> {
    let poczatek = Instant::now();
    loop {
        if odpowiada(port) {
            return Some(poczatek.elapsed());
        }
        if poczatek.elapsed() >= limit {
            return None;
        }
        sleep(ODSTEP_PROB);
    }
}
