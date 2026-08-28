//! Wskazanie rdzenia to czynność Operatora ustalająca serwer wdrożenia, na
//! którym stoi jego rdzeń, przyjmowana dopiero po udanej próbie połączenia.

use serde::Serialize;

use crate::rdzen::nasluch;
use crate::ustawienia::Ustawienia;

/// Stan wskazania oddawany interfejsowi, złożony z hosta, portu, gotowego
/// adresu HTTP i warstwy pochodzenia.
#[derive(Clone, Debug, Serialize)]
pub struct Wskazanie {
    /// Serwer wdrożenia, na którym stoi rdzeń; brak, dopóki nie wskazano.
    pub host: Option<String>,
    /// Port nasłuchu rdzenia obowiązujący.
    pub port: u16,
    /// Adres HTTP złożony z hosta i portu, gotowy dla warstwy połączenia.
    pub adres: Option<String>,
    /// Warstwa pochodzenia wskazania: brak, nastawy albo środowisko.
    pub warstwa: String,
}

/// Odmowa wskazania niesie kod powodu do rozgałęzienia logiki oraz gotowe
/// zdanie do wyświetlenia w oknie.
#[derive(Clone, Debug, Serialize)]
pub struct Odmowa {
    /// Kod powodu: `wskazanie-puste`, `port-niepoprawny`, `rdzen-nieosiagalny`,
    /// `zapis-nieudany`.
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

/// Zwraca wskazanie rdzenia obowiązujące w tej chwili, złożone z ustawień
/// zapisanych trwale albo ze zmiennej środowiska.
pub fn biezace(ustawienia: &Ustawienia) -> Wskazanie {
    let wskazane = ustawienia.wskazanie();
    Wskazanie {
        host: wskazane.as_ref().map(|w| w.host.clone()),
        port: ustawienia.port(),
        adres: wskazane.as_ref().map(|w| w.adres_http()),
        warstwa: ustawienia.warstwa_wskazania().nazwa().to_string(),
    }
}

/// Przyjmuje wskazanie Operatora podane w dowolnej z obsługiwanych postaci
/// i sprawdza łączność przed trwałym zapisem, bo zapis przed próbą utrwaliłby
/// wskazanie, które nie działa.
pub fn wskaz(ustawienia: &Ustawienia, adres: &str) -> Result<Wskazanie, Odmowa> {
    let (host, port) = rozbierz(adres, ustawienia.port())?;

    if !nasluch::odpowiada_pod(&host, port) {
        return Err(Odmowa::nowa(
            "rdzen-nieosiagalny",
            format!(
                "Pod adresem {host}:{port} nikt nie odpowiada. Sprawdź nazwę serwera i port, \
                 upewnij się, że rdzeń tam pracuje, i wskaż adres ponownie."
            ),
        ));
    }

    ustawienia
        .zapisz_wskazanie(&host, port)
        .map_err(|zdanie| Odmowa::nowa("zapis-nieudany", zdanie))?;

    Ok(biezace(ustawienia))
}

/// Rozbiera wskazanie Operatora na host i port, przyjmując też adres IPv6
/// w nawiasach; brak podanego portu znaczy port obowiązujący rdzenia.
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
            "Podaj nazwę serwera wdrożenia, na którym stoi rdzeń.".to_string(),
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
        Some((host, tekst)) if !host.contains(':') => Ok((host.to_string(), port_z_tekstu(tekst)?)),
        // Więcej niż jeden dwukropek bez nawiasów jest adresem IPv6 bez portu.
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
