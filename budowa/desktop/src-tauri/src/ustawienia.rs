//! Ustawienia powłoki obowiązujące w chwili pracy pochodzą z dwóch warstw:
//! nastaw zapisanych trwale albo zmiennej środowiska, silniejszej od nastaw.

use std::env;
use std::sync::{Arc, Mutex};

use crate::dziennik;
use crate::nastawy::{self, Nastawy};

/// Port nasłuchu rdzenia. Musi być równy `PortDomyslny`
/// z `server/internal/konfiguracja/ustawienia.go` (17870) — tamten plik jest
/// źródłem prawdy, ta stała jest jego jedyną kopią po stronie Rust powłoki.
pub const PORT_DOMYSLNY: u16 = 17870;

/// Port łącza szyfrowanego. Rdzeń za TLS stoi za pośrednikiem na porcie usługi
/// HTTPS, nie na własnym porcie nasłuchu, więc wskazanie `https://` bez portu
/// znaczy 443, a nie `PORT_DOMYSLNY`.
pub const PORT_DOMYSLNY_SZYFROWANY: u16 = 443;

/// Zmienna wskazująca port rdzenia, wspólna z rdzeniem i wzorcem pliku
/// środowiska przykładowego platformy.
pub const ZMIENNA_PORT: &str = "DANACO_PORT";
/// Zmienna wskazująca serwer wdrożenia, na którym stoi rdzeń. Stoi wyżej niż
/// wskazanie złożone w oknie; jej brak oddaje rozstrzygnięcie nastawom zapisanym.
pub const ZMIENNA_HOST_RDZENIA: &str = "DANACO_HOST_RDZENIA";
/// Zmienna wskazująca schemat, którym powłoka rozmawia z rdzeniem. Stoi wyżej
/// niż schemat wpisany w instalkę, tak samo jak `ZMIENNA_HOST_RDZENIA` stoi
/// wyżej niż `HOST_WDROZENIA`.
pub const ZMIENNA_SCHEMAT_RDZENIA: &str = "DANACO_SCHEMAT_RDZENIA";
/// Zmienna niosąca sekret nawiązania gniazda. Ta sama nazwa co
/// `ZmiennaSekretuNawiazania` w `server/internal/transport/ustawienia.go`:
/// rdzeń porównuje wartość spod tej zmiennej z wartością podaną przy
/// uaktualnieniu gniazda, więc obie strony muszą czytać tę samą nazwę.
/// Odpowiednika wpisywanego w instalkę sekret nie ma — wartość wkompilowana
/// w plik wykonywalny czyta z niego `strings`, a instalkę pobiera się
/// z kanału otwartego.
pub const ZMIENNA_SEKRET_NAWIAZANIA: &str = "DANACO_SEKRET_NAWIAZANIA";

/// Serwer wdrożenia wpisany w postać instalki przy jej składaniu. Rozstrzygnięcie 8
/// rejestru decyzji stanowi, że instalka jest osobna dla każdego wdrożenia, a żaden
/// z kroków kreatora o adres nie pyta — musi więc pochodzić stąd. Brak wpisania
/// zostawia powłokę bez wskazania, tak jak było przedtem.
pub const HOST_WDROZENIA: Option<&str> = option_env!("DANACO_HOST_WDROZENIA");

/// Port rdzenia wpisany w postać instalki przy jej składaniu, dopełnienie
/// `HOST_WDROZENIA`: wdrożenie za zaporą wystawia rdzeń na porcie innym niż
/// domyślny, a kreator o port nie pyta tak samo jak o adres.
pub const PORT_WDROZENIA: Option<&str> = option_env!("DANACO_PORT_WDROZENIA");

/// Schemat rdzenia wpisany w postać instalki przy jej składaniu, dopełnienie
/// `HOST_WDROZENIA`: rdzeń wdrożenia stoi za TLS, a instalka jest jedynym
/// miejscem, gdzie ta wiedza może wejść do powłoki, bo kreator o schemat
/// nie pyta tak samo jak o adres.
pub const SCHEMAT_WDROZENIA: Option<&str> = option_env!("DANACO_SCHEMAT_WDROZENIA");

/// Schemat, którym powłoka rozmawia z rdzeniem. Klient wywodzi z niego schemat
/// gniazda (`https:` → `wss:`, `http:` → `ws:`), więc od tej wartości zależy,
/// czy hasło Operatora i token sesji idą przez sieć szyfrem, czy otwartym tekstem.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Schemat {
    /// Łącze otwarte. Wystarcza rdzeniowi na pętli zwrotnej urządzenia.
    Http,
    /// Łącze szyfrowane TLS. Jedyne, którym wolno sięgać po rdzeń stojący poza
    /// urządzeniem Operatora.
    Https,
}

impl Schemat {
    /// Nazwa schematu w adresie — bez `://`, bo składa go dopiero adres.
    pub fn nazwa(self) -> &'static str {
        match self {
            Schemat::Http => "http",
            Schemat::Https => "https",
        }
    }

    /// Port usługi obowiązujący, gdy wskazanie portu nie podaje.
    pub fn port_domyslny(self) -> u16 {
        match self {
            Schemat::Http => PORT_DOMYSLNY,
            Schemat::Https => PORT_DOMYSLNY_SZYFROWANY,
        }
    }

    /// Rozpoznaje schemat w nazwie podanej wielkimi albo małymi literami,
    /// z `://` na końcu albo bez niego; nazwa nierozpoznana daje brak.
    pub fn z_tekstu(tekst: &str) -> Option<Schemat> {
        match tekst
            .trim()
            .trim_end_matches("://")
            .to_ascii_lowercase()
            .as_str()
        {
            "http" => Some(Schemat::Http),
            "https" => Some(Schemat::Https),
            _ => None,
        }
    }
}

/// Nazwa warstwy, z której pochodzi obowiązujące wskazanie hosta rdzenia,
/// zwracana w odpowiedzi polecenia wskazania.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Warstwa {
    /// Wskazania nie ma w żadnej warstwie — powłoka nie wie, gdzie szukać rdzenia.
    Brak,
    /// Wskazanie z nastaw zapisanych trwale, złożone w oknie.
    Nastawy,
    /// Wskazanie ze zmiennej środowiska procesu powłoki.
    Srodowisko,
    /// Wskazanie wpisane w postać instalki przy jej składaniu. Stoi najniżej:
    /// Operator zmienia je w oknie, a zmienna środowiska bije oba.
    Wdrozenie,
}

impl Warstwa {
    /// Nazwa warstwy dla interfejsu — jedno słowo, bez tłumaczenia po drodze.
    pub fn nazwa(self) -> &'static str {
        match self {
            Warstwa::Brak => "brak",
            Warstwa::Nastawy => "nastawy",
            Warstwa::Srodowisko => "srodowisko",
            Warstwa::Wdrozenie => "wdrozenie",
        }
    }
}

/// Złożone wskazanie, gdzie stoi rdzeń: schemat łącza, serwer wdrożenia i port
/// jego nasłuchu obowiązujący w tej chwili.
#[derive(Clone, Debug)]
pub struct Wskazane {
    /// Schemat łącza z rdzeniem.
    pub schemat: Schemat,
    /// Nazwa albo adres serwera wdrożenia.
    pub host: String,
    /// Port nasłuchu rdzenia na tym serwerze.
    pub port: u16,
}

impl Wskazane {
    /// Adres rdzenia rodziny HTTP złożony ze wskazania; jedyne miejsce, gdzie
    /// ten adres powstaje. Klient wywodzi z niego adres gniazda, więc schemat
    /// musi wejść stąd — inaczej `wss://` nie powstanie nigdy.
    pub fn adres_http(&self) -> String {
        format!(
            "{}://{}:{}",
            self.schemat.nazwa(),
            gospodarz_w_adresie(&self.host),
            self.port
        )
    }
}

/// Komplet ustawień powłoki obowiązujących w chwili pracy, złożony z warstwy
/// środowiska ustalonej raz przy starcie oraz warstwy nastaw zmienianej
/// w toku pracy okna.
#[derive(Clone, Debug)]
pub struct Ustawienia {
    /// Port wskazany zmienną środowiska; brak = nastawy albo port domyślny schematu.
    port_ze_srodowiska: Option<u16>,
    /// Host wskazany zmienną środowiska; brak = rozstrzygają nastawy.
    host_ze_srodowiska: Option<String>,
    /// Schemat przypięty zmienną środowiska; brak = nastawy albo wpis instalki.
    schemat_ze_srodowiska: Option<Schemat>,
    /// Sekret nawiązania gniazda wskazany zmienną środowiska; brak = rdzeń
    /// sekretu nie sprawdza, bo nie ma z czym porównywać.
    sekret_nawiazania: Option<String>,
    /// Nastawy zapisane trwale — warstwa zmienialna w trakcie pracy okna.
    zapisane: Arc<Mutex<Nastawy>>,
}

impl Ustawienia {
    /// Ustala ustawienia obowiązujące: czyta nastawy zapisane i zmienne
    /// środowiska procesu powłoki.
    pub fn ustal() -> Self {
        let ustawienia = Self {
            port_ze_srodowiska: port_ze_srodowiska(),
            host_ze_srodowiska: niepusta(ZMIENNA_HOST_RDZENIA),
            schemat_ze_srodowiska: niepusta(ZMIENNA_SCHEMAT_RDZENIA)
                .as_deref()
                .and_then(Schemat::z_tekstu),
            sekret_nawiazania: niepusta(ZMIENNA_SEKRET_NAWIAZANIA),
            zapisane: Arc::new(Mutex::new(nastawy::czytaj())),
        };
        ostrzez_o_schemacie_nieznanym(ustawienia.schemat());
        ustawienia
    }

    /// Wskazanie obowiązujące, gdy jest złożone: zmienna środowiska przed
    /// nastawami zapisanymi.
    pub fn wskazanie(&self) -> Option<Wskazane> {
        let host = match self.host_ze_srodowiska.as_deref() {
            Some(host) => host.to_string(),
            None => match self.nastawy().host_rdzenia {
                Some(host) => host,
                None => HOST_WDROZENIA.map(str::to_string)?,
            },
        };
        Some(Wskazane {
            schemat: self.schemat(),
            host,
            port: self.port(),
        })
    }

    /// Schemat łącza z rdzeniem obowiązujący: zmienna środowiska, nastawy
    /// zapisane, schemat wpisany w instalkę, inaczej łącze otwarte — ta sama
    /// kolejność warstw co przy hoście i porcie.
    pub fn schemat(&self) -> Schemat {
        self.schemat_ze_srodowiska
            .or_else(|| {
                self.nastawy()
                    .schemat_rdzenia
                    .as_deref()
                    .and_then(Schemat::z_tekstu)
            })
            .or_else(|| SCHEMAT_WDROZENIA.and_then(Schemat::z_tekstu))
            .unwrap_or(Schemat::Http)
    }

    /// Schemat przypięty zmienną środowiska albo brak; wskazanie z okna nie
    /// może go zmienić, bo zmienna stoi nad nastawami zapisanymi.
    pub fn schemat_ze_srodowiska(&self) -> Option<Schemat> {
        self.schemat_ze_srodowiska
    }

    /// Port rdzenia obowiązujący: zmienna środowiska, nastawy, port wpisany
    /// przy składaniu instalki, inaczej port domyślny schematu obowiązującego.
    pub fn port(&self) -> u16 {
        self.port_ze_srodowiska
            .or_else(|| self.nastawy().port_rdzenia)
            .or_else(|| PORT_WDROZENIA.and_then(|port| port.parse().ok()))
            .unwrap_or_else(|| self.schemat().port_domyslny())
    }

    /// Adres rdzenia obowiązujący albo `None`, gdy wskazania nie złożono.
    pub fn adres_rdzenia_http(&self) -> Option<String> {
        self.wskazanie().map(|wskazane| wskazane.adres_http())
    }

    /// Sekret nawiązania gniazda albo brak, gdy zmiennej nie wskazano.
    /// Powłoka sama gniazda nie otwiera — otwiera je strona interfejsu — więc
    /// sekret idzie do niej skryptem wstępnym okna, a stamtąd do rdzenia przy
    /// uaktualnieniu gniazda.
    pub fn sekret_nawiazania(&self) -> Option<String> {
        self.sekret_nawiazania.clone()
    }

    /// Warstwa, z której pochodzi obowiązujące wskazanie hosta.
    pub fn warstwa_wskazania(&self) -> Warstwa {
        if self.host_ze_srodowiska.is_some() {
            return Warstwa::Srodowisko;
        }
        if self.nastawy().wskazanie_zlozone() {
            return Warstwa::Nastawy;
        }
        if HOST_WDROZENIA.is_some() {
            return Warstwa::Wdrozenie;
        }
        Warstwa::Brak
    }

    /// Zapisuje wskazanie Operatora trwale i wprowadza je w życie dla
    /// wszystkich kopii ustawień. Schemat idzie do pliku razem z hostem —
    /// inaczej wskazanie `https://` obowiązywałoby do zamknięcia powłoki.
    pub fn zapisz_wskazanie(&self, schemat: Schemat, host: &str, port: u16) -> Result<(), String> {
        let nowe = Nastawy {
            host_rdzenia: Some(host.to_string()),
            port_rdzenia: Some(port),
            schemat_rdzenia: Some(schemat.nazwa().to_string()),
        };
        nastawy::zapisz(&nowe)?;
        match self.zapisane.lock() {
            Ok(mut zamek) => {
                *zamek = nowe;
                Ok(())
            }
            // Zamek zatruty paniką innego wątku: plik jest już zapisany,
            // wskazanie obowiązuje od startu.
            Err(_) => Err(format!(
                "Wskazanie zapisano w pliku {}, ale nie weszło w życie w tym uruchomieniu — \
                 zamknij okno i otwórz je ponownie.",
                nastawy::sciezka().display()
            )),
        }
    }

    /// Kopia nastaw zapisanych; zamek zatruty daje nastawy puste, nie panikę.
    fn nastawy(&self) -> Nastawy {
        match self.zapisane.lock() {
            Ok(zamek) => zamek.clone(),
            Err(_) => Nastawy::default(),
        }
    }
}

/// Wstawia gospodarza do adresu URL: adres IPv6 wchodzi w nawiasach kwadratowych,
/// bo bez nich jego dwukropki zlewają się z dwukropkiem portu.
pub fn gospodarz_w_adresie(host: &str) -> String {
    if host.contains(':') && !host.starts_with('[') {
        format!("[{host}]")
    } else {
        host.to_string()
    }
}

/// Dopisuje do dziennika wartość zmiennej schematu, której nie da się rozpoznać.
/// Bez tego wpisu literówka w nazwie schematu cicho zostawiałaby łącze otwarte —
/// czyli dawałaby skutek odwrotny do zamierzonego przez tego, kto ją ustawiał.
fn ostrzez_o_schemacie_nieznanym(obowiazujacy: Schemat) {
    if let Some(tekst) = niepusta(ZMIENNA_SCHEMAT_RDZENIA) {
        if Schemat::z_tekstu(&tekst).is_none() {
            dziennik::dopisz(&format!(
                "ustawienia: {ZMIENNA_SCHEMAT_RDZENIA}={tekst} nie jest schematem \
                 („http” albo „https”) — obowiązuje {}",
                obowiazujacy.nazwa()
            ));
        }
    }
}

/// Odczyt portu ze środowiska: wartość niebędąca liczbą nie przerywa startu,
/// tylko zostaje pominięta na rzecz warstw niższych.
fn port_ze_srodowiska() -> Option<u16> {
    niepusta(ZMIENNA_PORT).and_then(|tekst| tekst.parse().ok())
}

/// Zwraca wartość zmiennej środowiska, traktując wartość pustą jak brak
/// ustawienia w tej warstwie wskazania.
fn niepusta(nazwa: &str) -> Option<String> {
    match env::var(nazwa) {
        Ok(wartosc) if !wartosc.trim().is_empty() => Some(wartosc.trim().to_string()),
        _ => None,
    }
}
