//! Wskazanie rdzenia — czynność Operatora „na którym serwerze stoi mój rdzeń".
//!
//! Po co ta czynność istnieje. Rdzeń stoi na serwerze wdrożenia, a instalator
//! adresu tego serwera nie zna i znać go nie może: w chwili rozpakowania plików
//! nikt jeszcze nie wie, pod jaką nazwą stoi rdzeń tego Operatora. Wie to sam
//! Operator i mówi to przy pierwszym uruchomieniu okna. Drugą drogą wskazania
//! jest zmienna środowiska — droga wykonawcy i jednostki usługi, nie Operatora.
//!
//! Dwie rzeczy dzieją się tutaj i tylko tutaj:
//!
//!   1. rozbiór wskazania na host i port (jedno pole w oknie, jedno rozumienie
//!      po obu stronach granicy procesu),
//!   2. próba połączenia PRZED zapisem — wskazanie nieosiągalne nie zostaje
//!      przyjęte, bo zapisane zamieniłoby okno w ekran milczący po każdym starcie.
//!
//! Czego tu nie ma: wyboru „rdzeń na tym urządzeniu". Produkt występuje wyłącznie
//! jako hybryda, więc rdzenia na urządzeniu Operatora nie ma, a powłoka nie ma
//! czym go postawić.

use serde::Serialize;

use crate::rdzen::nasluch;
use crate::ustawienia::Ustawienia;

/// Stan wskazania oddawany interfejsowi.
#[derive(Clone, Debug, Serialize)]
pub struct Wskazanie {
    /// Serwer wdrożenia, na którym stoi rdzeń; brak, dopóki nie wskazano.
    /// Okno wstawia go z powrotem do pola, żeby Operator poprawiał to, co napisał.
    pub host: Option<String>,
    /// Port nasłuchu rdzenia obowiązujący.
    pub port: u16,
    /// Adres HTTP złożony z hosta i portu — gotowa postać dla warstwy połączenia
    /// (`klient/src/polaczenie/adres-rdzenia.ts`), która sama adresu nie składa.
    pub adres: Option<String>,
    /// Warstwa, z której pochodzi wskazanie: `brak`, `nastawy`, `srodowisko`.
    /// `brak` znaczy pierwsze uruchomienie i jest dla okna sygnałem, że ma
    /// zapytać. Wskazania ze zmiennej środowiska okno nie nadpisze — ma o tym
    /// powiedzieć, zamiast przyjmować zapis bez skutku.
    pub warstwa: String,
}

/// Odmowa wskazania — powód do rozgałęzienia i gotowe zdanie do okna.
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

/// Zwraca wskazanie obowiązujące w tej chwili.
pub fn biezace(ustawienia: &Ustawienia) -> Wskazanie {
    let wskazane = ustawienia.wskazanie();
    Wskazanie {
        host: wskazane.as_ref().map(|w| w.host.clone()),
        port: ustawienia.port(),
        adres: wskazane.as_ref().map(|w| w.adres_http()),
        warstwa: ustawienia.warstwa_wskazania().nazwa().to_string(),
    }
}

/// Przyjmuje wskazanie Operatora: sprawdza łączność i zapisuje je trwale.
///
/// `adres` przychodzi z okna w postaci, w jakiej Operator go napisał: nazwa
/// serwera albo `serwer:port`, z przedrostkiem `http://` albo bez.
///
/// Kolejność jest wiążąca: najpierw próba połączenia, potem zapis. Zapis przed
/// próbą utrwalałby wskazanie, które nie działa, a Operator zobaczyłby skutek
/// dopiero przy następnym starcie.
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

/// Rozbiera wskazanie Operatora na host i port.
///
/// Przyjmuje `serwer`, `serwer:17870`, `http://serwer:17870` oraz adres IPv6
/// w nawiasach (`[::1]:17870`), bo Operator wpisuje to, co ma zapisane, a nie to,
/// co wygodne dla rozbioru. Brak portu znaczy port obowiązujący — ten sam, na
/// którym rdzeń nasłuchuje domyślnie.
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
