//! Wskazanie rdzenia to czynność Operatora ustalająca serwer wdrożenia, na
//! którym stoi jego rdzeń, przyjmowana dopiero po udanej próbie połączenia.

use serde::Serialize;

use crate::rdzen::nasluch;
use crate::ustawienia::{Schemat, Ustawienia, ZMIENNA_SCHEMAT_RDZENIA};

/// Stan wskazania oddawany interfejsowi, złożony ze schematu, hosta, portu,
/// gotowego adresu i warstwy pochodzenia.
#[derive(Clone, Debug, Serialize)]
pub struct Wskazanie {
    /// Schemat łącza z rdzeniem: `http` albo `https`.
    pub schemat: String,
    /// Serwer wdrożenia, na którym stoi rdzeń; brak, dopóki nie wskazano.
    pub host: Option<String>,
    /// Port nasłuchu rdzenia obowiązujący.
    pub port: u16,
    /// Adres złożony ze schematu, hosta i portu, gotowy dla warstwy połączenia.
    pub adres: Option<String>,
    /// Warstwa pochodzenia wskazania: brak, nastawy, środowisko albo wdrożenie.
    pub warstwa: String,
}

/// Odmowa wskazania niesie kod powodu do rozgałęzienia logiki oraz gotowe
/// zdanie do wyświetlenia w oknie.
#[derive(Clone, Debug, Serialize)]
pub struct Odmowa {
    /// Kod powodu: `wskazanie-puste`, `port-niepoprawny`, `schemat-nieznany`,
    /// `schemat-niezgodny`, `rdzen-nieosiagalny`, `zapis-nieudany`.
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
        schemat: ustawienia.schemat().nazwa().to_string(),
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
    let schemat = ustawienia.schemat();
    let (schemat_podany, host, port_podany) = rozbierz(adres)?;

    /* Schemat obowiązujący niosą zmienna środowiska i wpis instalki, a nie plik
    nastaw powłoki — okno nie ma go gdzie zapisać, więc wskazanie ze
    schematem innym niż obowiązujący jest odmawiane, zamiast obowiązywać do
    najbliższego zamknięcia powłoki i cicho wracać do poprzedniego. */
    if let Some(podany) = schemat_podany {
        if podany != schemat {
            return Err(Odmowa::nowa(
                "schemat-niezgodny",
                format!(
                    "Powłoka rozmawia z rdzeniem schematem „{}”, a wskazanie podaje „{}”. \
                     Schemat niesie zmienna środowiska {ZMIENNA_SCHEMAT_RDZENIA} albo wpis \
                     instalki — okno go nie zmienia. Podaj sam adres serwera, bez przedrostka.",
                    schemat.nazwa(),
                    podany.nazwa()
                ),
            ));
        }
    }

    let port = port_podany.unwrap_or_else(|| ustawienia.port());

    if !nasluch::odpowiada_pod(&host, port) {
        return Err(Odmowa::nowa(
            "rdzen-nieosiagalny",
            format!(
                "Pod adresem {}://{host}:{port} nikt nie odpowiada. Sprawdź nazwę serwera \
                 i port, upewnij się, że rdzeń tam pracuje, i wskaż adres ponownie.",
                schemat.nazwa()
            ),
        ));
    }

    ustawienia
        .zapisz_wskazanie(&host, port)
        .map_err(|zdanie| Odmowa::nowa("zapis-nieudany", zdanie))?;

    Ok(biezace(ustawienia))
}

/// Rozbiera wskazanie Operatora na schemat, host i port, przyjmując też adres
/// IPv6 w nawiasach. Schemat i port podane są brakiem, gdy wskazanie ich nie
/// niesie — rozstrzyga wtedy warstwa ustawień, nie ta funkcja.
fn rozbierz(adres: &str) -> Result<(Option<Schemat>, String, Option<u16>), Odmowa> {
    let adres = adres.trim();
    // Przedrostek zostaje rozpoznany, nie obcięty: schemat rozstrzyga o tym,
    // czy klient złoży adres gniazda `ws:`, czy `wss:`.
    let (schemat, reszta) = match adres.split_once("://") {
        Some((przedrostek, reszta)) => match Schemat::z_tekstu(przedrostek) {
            Some(schemat) => (Some(schemat), reszta),
            None => {
                return Err(Odmowa::nowa(
                    "schemat-nieznany",
                    format!(
                        "„{przedrostek}://” nie jest schematem, którym powłoka rozmawia \
                         z rdzeniem. Podaj adres z przedrostkiem „http://” albo „https://”, \
                         albo bez przedrostka."
                    ),
                ))
            }
        },
        None => (None, adres),
    };

    let bez_przedrostka = reszta.trim_end_matches('/').trim();
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
                "Adres w nawiasach kwadratowych jest niedomknięty — brakuje znaku „]”.".to_string(),
            ));
        };
        let port = match reszta.strip_prefix(':') {
            Some(tekst) => Some(port_z_tekstu(tekst)?),
            None => None,
        };
        return Ok((schemat, host.to_string(), port));
    }

    match bez_przedrostka.rsplit_once(':') {
        Some((host, tekst)) if !host.contains(':') => {
            Ok((schemat, host.to_string(), Some(port_z_tekstu(tekst)?)))
        }
        // Więcej niż jeden dwukropek bez nawiasów jest adresem IPv6 bez portu.
        Some(_) => Ok((schemat, bez_przedrostka.to_string(), None)),
        None => Ok((schemat, bez_przedrostka.to_string(), None)),
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

    #[test]
    fn nazwa_bez_portu_zostawia_rozstrzygniecie_ustawieniom() {
        assert_eq!(
            rozbierz("danaco-system").unwrap(),
            (None, "danaco-system".to_string(), None)
        );
    }

    #[test]
    fn przedrostek_zostaje_rozpoznany_a_ukosnik_nie_wchodzi_do_nazwy_hosta() {
        assert_eq!(
            rozbierz("http://danaco-system:18000/").unwrap(),
            (
                Some(Schemat::Http),
                "danaco-system".to_string(),
                Some(18000)
            )
        );
        assert_eq!(
            rozbierz("HTTPS://danaco-system/").unwrap(),
            (Some(Schemat::Https), "danaco-system".to_string(), None)
        );
    }

    #[test]
    fn schemat_spoza_rodziny_http_jest_odmowa() {
        assert_eq!(
            rozbierz("ftp://danaco-system").unwrap_err().powod,
            "schemat-nieznany"
        );
    }

    #[test]
    fn adres_ipv6_w_nawiasach_zachowuje_dwukropki() {
        assert_eq!(
            rozbierz("[2001:db8::7]:18000").unwrap(),
            (None, "2001:db8::7".to_string(), Some(18000))
        );
        assert_eq!(
            rozbierz("[2001:db8::7]").unwrap(),
            (None, "2001:db8::7".to_string(), None)
        );
    }

    #[test]
    fn adres_ipv6_bez_nawiasow_jest_calym_hostem() {
        assert_eq!(
            rozbierz("2001:db8::7").unwrap(),
            (None, "2001:db8::7".to_string(), None)
        );
    }

    #[test]
    fn wskazanie_puste_jest_odmowa() {
        assert_eq!(rozbierz("   ").unwrap_err().powod, "wskazanie-puste");
        assert_eq!(rozbierz("http://").unwrap_err().powod, "wskazanie-puste");
    }

    #[test]
    fn port_nieliczbowy_i_zerowy_sa_odmowa() {
        assert_eq!(
            rozbierz("danaco-system:port").unwrap_err().powod,
            "port-niepoprawny"
        );
        assert_eq!(
            rozbierz("danaco-system:0").unwrap_err().powod,
            "port-niepoprawny"
        );
        assert_eq!(
            rozbierz("danaco-system:70000").unwrap_err().powod,
            "port-niepoprawny"
        );
    }
}
