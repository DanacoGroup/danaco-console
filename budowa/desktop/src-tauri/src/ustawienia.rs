//! Ustawienia powłoki obowiązujące w chwili pracy pochodzą z dwóch warstw:
//! nastaw zapisanych trwale albo zmiennej środowiska, silniejszej od nastaw.

use std::env;
use std::sync::{Arc, Mutex};

use crate::nastawy::{self, Nastawy};

/// Port nasłuchu rdzenia. Musi być równy `PortDomyslny`
/// z `server/internal/konfiguracja/ustawienia.go` (17870) — tamten plik jest
/// źródłem prawdy, ta stała jest jego jedyną kopią po stronie Rust powłoki.
pub const PORT_DOMYSLNY: u16 = 17870;

/// Zmienna wskazująca port rdzenia, wspólna z rdzeniem i wzorcem pliku
/// środowiska przykładowego platformy.
pub const ZMIENNA_PORT: &str = "DANACO_PORT";
/// Zmienna wskazująca serwer wdrożenia, na którym stoi rdzeń. Stoi wyżej niż
/// wskazanie złożone w oknie; jej brak oddaje rozstrzygnięcie nastawom zapisanym.
pub const ZMIENNA_HOST_RDZENIA: &str = "DANACO_HOST_RDZENIA";

/// Serwer wdrożenia wpisany w postać instalki przy jej składaniu. Rozstrzygnięcie 8
/// rejestru decyzji stanowi, że instalka jest osobna dla każdego wdrożenia, a żaden
/// z kroków kreatora o adres nie pyta — musi więc pochodzić stąd. Brak wpisania
/// zostawia powłokę bez wskazania, tak jak było przedtem.
pub const HOST_WDROZENIA: Option<&str> = option_env!("DANACO_HOST_WDROZENIA");

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

/// Złożone wskazanie, gdzie stoi rdzeń: serwer wdrożenia i port jego
/// nasłuchu obowiązujący w tej chwili.
#[derive(Clone, Debug)]
pub struct Wskazane {
    /// Nazwa albo adres serwera wdrożenia.
    pub host: String,
    /// Port nasłuchu rdzenia na tym serwerze.
    pub port: u16,
}

impl Wskazane {
    /// Adres HTTP rdzenia złożony ze wskazania; jedyne miejsce, gdzie ten
    /// adres powstaje.
    pub fn adres_http(&self) -> String {
        format!("http://{}:{}", self.host, self.port)
    }
}

/// Komplet ustawień powłoki obowiązujących w chwili pracy, złożony z warstwy
/// środowiska ustalonej raz przy starcie oraz warstwy nastaw zmienianej
/// w toku pracy okna.
#[derive(Clone, Debug)]
pub struct Ustawienia {
    /// Port wskazany zmienną środowiska; brak = nastawy albo `PORT_DOMYSLNY`.
    port_ze_srodowiska: Option<u16>,
    /// Host wskazany zmienną środowiska; brak = rozstrzygają nastawy.
    host_ze_srodowiska: Option<String>,
    /// Nastawy zapisane trwale — warstwa zmienialna w trakcie pracy okna.
    zapisane: Arc<Mutex<Nastawy>>,
}

impl Ustawienia {
    /// Ustala ustawienia obowiązujące: czyta nastawy zapisane i zmienne
    /// środowiska procesu powłoki.
    pub fn ustal() -> Self {
        Self {
            port_ze_srodowiska: port_ze_srodowiska(),
            host_ze_srodowiska: niepusta(ZMIENNA_HOST_RDZENIA),
            zapisane: Arc::new(Mutex::new(nastawy::czytaj())),
        }
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
            host,
            port: self.port(),
        })
    }

    /// Port rdzenia obowiązujący: zmienna środowiska albo nastawy, inaczej
    /// `PORT_DOMYSLNY`.
    pub fn port(&self) -> u16 {
        self.port_ze_srodowiska
            .or_else(|| self.nastawy().port_rdzenia)
            .unwrap_or(PORT_DOMYSLNY)
    }

    /// Adres HTTP rdzenia obowiązujący albo `None`, gdy wskazania nie złożono.
    pub fn adres_rdzenia_http(&self) -> Option<String> {
        self.wskazanie().map(|wskazane| wskazane.adres_http())
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
    /// wszystkich kopii ustawień.
    pub fn zapisz_wskazanie(&self, host: &str, port: u16) -> Result<(), String> {
        let nowe = Nastawy {
            host_rdzenia: Some(host.to_string()),
            port_rdzenia: Some(port),
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
