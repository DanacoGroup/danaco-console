//! Wskazanie rdzenia — czynność Operatora „gdzie stoi mój rdzeń".
//!
//! Po co ta czynność istnieje. Instalator nie zna adresu serwera i znać go nie
//! może; wie go wyłącznie Operator i mówi go przy pierwszym uruchomieniu okna.
//! Do tej pory jedyną drogą wskazania była zmienna środowiska — droga wykonawcy,
//! nie Operatora — więc powłoka z pudełka zawsze stawiała rdzeń na maszynie
//! Operatora, czyli tam, gdzie nie ma arsenału programów przetwarzających.
//!
//! Trzy rzeczy dzieją się tutaj i tylko tutaj:
//!
//!   1. rozbiór wskazania na host i port (jedno pole w oknie, jedno rozumienie
//!      po obu stronach granicy procesu),
//!   2. próba połączenia PRZED zapisem — wskazanie nieosiągalne nie zostaje
//!      przyjęte, bo zapisane zamieniłoby okno w ekran milczący po każdym starcie,
//!   3. wprowadzenie wskazania w życie w tym uruchomieniu: rdzeń lokalny staje,
//!      rdzeń zdalny zostaje rozpoznany, a proces postawiony przy poprzednim
//!      wskazaniu jest zatrzymywany (`UchwytRdzenia::przejmij`).
//!
//! Czego tu nie ma: wyboru, skąd ładuje się strona interfejsu. To osobna nastawa
//! (`DANACO_ADRES_INTERFEJSU`, `zrodlo_interfejsu.rs`) i osobny skutek.

use serde::Serialize;

use crate::rdzen::{uruchomienie, UchwytRdzenia};
use crate::ustawienia::{Ustawienia, HOST_DOMYSLNY};

/// Stan wskazania oddawany interfejsowi.
#[derive(Clone, Debug, Serialize)]
pub struct Wskazanie {
    /// Host rdzenia obowiązujący.
    pub host: String,
    /// Port rdzenia obowiązujący.
    pub port: u16,
    /// Adres HTTP rdzenia obowiązujący — złożony z hosta i portu.
    pub adres: String,
    /// Czy wskazanie zostało złożone (w oknie albo zmienną środowiska). Fałsz
    /// znaczy pierwsze uruchomienie i jest dla okna sygnałem, że ma zapytać.
    pub zlozone: bool,
    /// Czy rdzeń stoi na tym urządzeniu.
    pub rdzen_lokalny: bool,
    /// Warstwa, z której pochodzi wskazanie: `brak`, `nastawy`, `srodowisko`.
    /// Wskazania ze zmiennej środowiska okno nie nadpisze — ma o tym powiedzieć,
    /// zamiast przyjmować zapis bez skutku.
    pub warstwa: String,
}

/// Odmowa wskazania — powód do rozgałęzienia i gotowe zdanie do okna.
#[derive(Clone, Debug, Serialize)]
pub struct Odmowa {
    /// Kod powodu: `wskazanie-puste`, `port-niepoprawny`, `rdzen-nieosiagalny`,
    /// `zapis-nieudany`, `rdzen-nie-wstal`.
    pub powod: String,
    /// Zdanie dla Operatora: co się stało i co może zrobić.
    pub zdanie: String,
}

impl Odmowa {
    /// Składa odmowę z kodu powodu i zdania.
    fn nowa(powod: &str, zdanie: String) -> Self {
        Self {
            powod: powod.to_string(),
            zdanie,
        }
    }
}

/// Zwraca wskazanie obowiązujące w tej chwili.
pub fn biezace(ustawienia: &Ustawienia) -> Wskazanie {
    Wskazanie {
        host: ustawienia.host(),
        port: ustawienia.port(),
        adres: ustawienia.adres_rdzenia_http(),
        zlozone: ustawienia.wskazanie_zlozone(),
        rdzen_lokalny: ustawienia.rdzen_lokalny(),
        warstwa: ustawienia.warstwa_wskazania().nazwa().to_string(),
    }
}

/// Przyjmuje wskazanie Operatora: sprawdza, zapisuje trwale i wprowadza w życie.
///
/// `adres` przychodzi z okna w postaci, w jakiej Operator go napisał: nazwa hosta
/// albo `host:port`, z przedrostkiem `http://` albo bez. Wskazanie puste znaczy
/// rdzeń na tym urządzeniu — okno wysyła je, gdy Operator wybierze pracę lokalną.
///
/// Kolejność jest wiążąca: najpierw próba połączenia (dla wariantu lokalnego —
/// postawienie rdzenia), potem zapis. Zapis przed próbą utrwalałby wskazanie,
/// które nie działa, a Operator zobaczyłby skutek dopiero przy następnym starcie.
pub fn wskaz(
    ustawienia: &Ustawienia,
    uchwyt: &UchwytRdzenia,
    adres: &str,
    lokalnie: bool,
) -> Result<Wskazanie, Odmowa> {
    let (host, port) = if lokalnie {
        (HOST_DOMYSLNY.to_string(), ustawienia.port())
    } else {
        rozbierz(adres, ustawienia.port())?
    };

    if host == HOST_DOMYSLNY {
        return wskaz_lokalnie(ustawienia, uchwyt, port);
    }
    wskaz_zdalnie(ustawienia, uchwyt, &host, port)
}

/// Wskazanie rdzenia na tym urządzeniu: zapis, potem postawienie procesu.
///
/// Tu zapis idzie PRZED postawieniem, odwrotnie niż w wariancie zdalnym, i ma to
/// powód: rdzeń stawia się z ustawień obowiązujących, więc port ze wskazania musi
/// już w nich siedzieć. Rdzeń, który nie wstał, kończy się odmową — ale wskazanie
/// zostaje zapisane, bo wybór miejsca pracy jest trafny także wtedy, gdy binarki
/// nie znaleziono; Operator dostaje w zdaniu powód z opisu przebiegu.
fn wskaz_lokalnie(
    ustawienia: &Ustawienia,
    uchwyt: &UchwytRdzenia,
    port: u16,
) -> Result<Wskazanie, Odmowa> {
    ustawienia
        .zapisz_wskazanie(HOST_DOMYSLNY, port)
        .map_err(|zdanie| Odmowa::nowa("zapis-nieudany", zdanie))?;

    let zawiazanie = uruchomienie::zawiaz_lokalnie(ustawienia);
    let pracuje = zawiazanie.opis.pracuje;
    let opis_przebiegu = zawiazanie.opis.opis.clone();
    uchwyt.przejmij(
        zawiazanie.port,
        zawiazanie.host,
        zawiazanie.opis,
        zawiazanie.dziecko,
    );

    if !pracuje {
        return Err(Odmowa::nowa(
            "rdzen-nie-wstal",
            format!("Rdzeń na tym urządzeniu nie odpowiedział. {opis_przebiegu}"),
        ));
    }
    Ok(biezace(ustawienia))
}

/// Wskazanie rdzenia na serwerze: próba połączenia, potem zapis.
fn wskaz_zdalnie(
    ustawienia: &Ustawienia,
    uchwyt: &UchwytRdzenia,
    host: &str,
    port: u16,
) -> Result<Wskazanie, Odmowa> {
    if !crate::rdzen::nasluch::odpowiada_pod(host, port) {
        return Err(Odmowa::nowa(
            "rdzen-nieosiagalny",
            format!(
                "Pod adresem {host}:{port} nikt nie odpowiada. Sprawdź nazwę serwera i port, \
                 upewnij się, że rdzeń tam pracuje, i wskaż adres ponownie."
            ),
        ));
    }

    ustawienia
        .zapisz_wskazanie(host, port)
        .map_err(|zdanie| Odmowa::nowa("zapis-nieudany", zdanie))?;

    let zawiazanie = uruchomienie::zawiaz_zdalnie(ustawienia);
    uchwyt.przejmij(
        zawiazanie.port,
        zawiazanie.host,
        zawiazanie.opis,
        zawiazanie.dziecko,
    );
    Ok(biezace(ustawienia))
}

/// Rozbiera wskazanie Operatora na host i port.
///
/// Przyjmuje `serwer`, `serwer:17870`, `http://serwer:17870` oraz adres IPv6
/// w nawiasach (`[::1]:17870`), bo Operator wpisuje to, co ma zapisane, a nie to,
/// co wygodne dla rozbioru. Brak portu znaczy port obowiązujący — ten sam, na
/// którym stoi rdzeń domyślnie.
fn rozbierz(adres: &str, port_obowiazujacy: u16) -> Result<(String, u16), Odmowa> {
    let bez_przedrostka = adres
        .trim()
        .trim_start_matches("http://")
        .trim_start_matches("https://")
        .trim_end_matches('/')
        .trim();
    if bez_przedrostka.is_empty() {
        return Err(Odmowa::nowa(
            "wskazanie-puste",
            "Podaj nazwę serwera, na którym stoi rdzeń, albo wybierz pracę z rdzeniem na tym \
             urządzeniu."
                .to_string(),
        ));
    }

    // Adres IPv6 w nawiasach: dwukropki należą do adresu, port stoi za nawiasem.
    if let Some(koniec) = bez_przedrostka.strip_prefix('[') {
        let Some((host, reszta)) = koniec.split_once(']') else {
            return Err(Odmowa::nowa(
                "wskazanie-puste",
                "Adres w nawiasach kwadratowych jest niedomknięty — brakuje znaku „]”."
                    .to_string(),
            ));
        };
        let port = match reszta.strip_prefix(':') {
            Some(tekst) => port_z_tekstu(tekst)?,
            None => port_obowiazujacy,
        };
        return Ok((host.to_string(), port));
    }

    match bez_przedrostka.rsplit_once(':') {
        Some((host, tekst)) if !host.contains(':') => {
            Ok((host.to_string(), port_z_tekstu(tekst)?))
        }
        // Więcej niż jeden dwukropek bez nawiasów to adres IPv6 podany bez nich —
        // portu w nim nie ma, cała treść jest hostem.
        Some(_) => Ok((bez_przedrostka.to_string(), port_obowiazujacy)),
        None => Ok((bez_przedrostka.to_string(), port_obowiazujacy)),
    }
}

/// Rozbiera port ze wskazania. Wartość spoza zakresu portów jest odmową, nie
/// cichym zastąpieniem wartością domyślną: Operator ma poprawić literówkę,
/// zamiast szukać, dlaczego okno łączy się nie tam, gdzie napisał.
fn port_z_tekstu(tekst: &str) -> Result<u16, Odmowa> {
    match tekst.trim().parse::<u16>() {
        Ok(port) if port > 0 => Ok(port),
        _ => Err(Odmowa::nowa(
            "port-niepoprawny",
            format!("„{tekst}” nie jest numerem portu. Podaj liczbę z zakresu od 1 do 65535."),
        )),
    }
}

#[cfg(test)]
mod testy {
    use super::*;

    const PORT: u16 = 17870;

    #[test]
    fn nazwa_bez_portu_bierze_port_obowiazujacy() {
        assert_eq!(
            rozbierz("danaco-system", PORT).unwrap(),
            ("danaco-system".to_string(), PORT)
        );
    }

    #[test]
    fn przedrostek_i_ukosnik_nie_wchodza_do_nazwy_hosta() {
        assert_eq!(
            rozbierz("http://danaco-system:18000/", PORT).unwrap(),
            ("danaco-system".to_string(), 18000)
        );
    }

    #[test]
    fn adres_ipv6_w_nawiasach_zachowuje_dwukropki() {
        assert_eq!(
            rozbierz("[2001:db8::7]:18000", PORT).unwrap(),
            ("2001:db8::7".to_string(), 18000)
        );
        assert_eq!(
            rozbierz("[2001:db8::7]", PORT).unwrap(),
            ("2001:db8::7".to_string(), PORT)
        );
    }

    #[test]
    fn adres_ipv6_bez_nawiasow_jest_calym_hostem() {
        assert_eq!(
            rozbierz("2001:db8::7", PORT).unwrap(),
            ("2001:db8::7".to_string(), PORT)
        );
    }

    #[test]
    fn wskazanie_puste_jest_odmowa() {
        assert_eq!(rozbierz("   ", PORT).unwrap_err().powod, "wskazanie-puste");
        assert_eq!(
            rozbierz("http://", PORT).unwrap_err().powod,
            "wskazanie-puste"
        );
    }

    #[test]
    fn port_nieliczbowy_i_zerowy_sa_odmowa() {
        assert_eq!(
            rozbierz("danaco-system:port", PORT).unwrap_err().powod,
            "port-niepoprawny"
        );
        assert_eq!(
            rozbierz("danaco-system:0", PORT).unwrap_err().powod,
            "port-niepoprawny"
        );
        assert_eq!(
            rozbierz("danaco-system:70000", PORT).unwrap_err().powod,
            "port-niepoprawny"
        );
    }
}
