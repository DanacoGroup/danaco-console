//! Ustawienia powłoki obowiązujące w chwili pracy.
//!
//! Dwie warstwy, od słabszej: nastawy zapisane trwale (`nastawy.rs`) → zmienna
//! środowiska. Trzeciej — wartości domyślnej hosta — nie ma i nie może być:
//! rdzeń stoi na serwerze wdrożenia, a jego nazwy nie zna ani powłoka, ani
//! instalator. Brak obu warstw znaczy więc „wskazania nie złożono", a nie
//! „rdzeń stoi tu obok". Nazwy `DANACO_PORT` i `DANACO_KATALOG_DANYCH` są
//! własnością rdzenia (`server/internal/konfiguracja/srodowisko.go`); powłoka
//! je wyłącznie czyta.
//!
//! Dlaczego zmienna stoi nad plikiem. Plik niesie wskazanie Operatora złożone
//! w oknie i ma przetrwać zamknięcie okna. Zmienna niesie wskazanie tego, kto
//! stawia proces — wykonawcy przy budowie, jednostki usługi na serwerze — i musi
//! brać górę, bo inaczej plik z jednej maszyny sterowałby uruchomieniem na
//! drugiej po skopiowaniu profilu.

use std::env;
use std::sync::{Arc, Mutex};

use crate::nastawy::{self, Nastawy};

/// Port nasłuchu rdzenia. Musi być równy `PortDomyslny`
/// z `server/internal/konfiguracja/ustawienia.go` (17870) — tamten plik jest
/// źródłem prawdy, ta stała jest jego jedyną kopią po stronie Rust powłoki.
pub const PORT_DOMYSLNY: u16 = 17870;

/// Zmienna wskazująca port rdzenia — wspólna z rdzeniem i wzorcem `.env.example`.
pub const ZMIENNA_PORT: &str = "DANACO_PORT";
/// Zmienna wskazująca serwer wdrożenia, na którym stoi rdzeń. Stoi wyżej niż
/// wskazanie złożone w oknie; jej brak oddaje rozstrzygnięcie nastawom zapisanym.
pub const ZMIENNA_HOST_RDZENIA: &str = "DANACO_HOST_RDZENIA";

/// Nazwa warstwy, z której pochodzi obowiązujące wskazanie hosta. Wchodzi do
/// odpowiedzi polecenia `wskazanie_rdzenia`, żeby okno mogło powiedzieć
/// Operatorowi, dlaczego pola nie da się zmienić: wskazanie ze zmiennej
/// środowiska jest silniejsze od zapisu w oknie.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Warstwa {
    /// Wskazania nie ma w żadnej warstwie — powłoka nie wie, gdzie szukać rdzenia.
    Brak,
    /// Wskazanie z nastaw zapisanych trwale, złożone w oknie.
    Nastawy,
    /// Wskazanie ze zmiennej środowiska procesu powłoki.
    Srodowisko,
}

impl Warstwa {
    /// Nazwa warstwy dla interfejsu — jedno słowo, bez tłumaczenia po drodze.
    pub fn nazwa(self) -> &'static str {
        match self {
            Warstwa::Brak => "brak",
            Warstwa::Nastawy => "nastawy",
            Warstwa::Srodowisko => "srodowisko",
        }
    }
}

/// Złożone wskazanie, gdzie stoi rdzeń: serwer wdrożenia i port jego nasłuchu.
#[derive(Clone, Debug)]
pub struct Wskazane {
    /// Nazwa albo adres serwera wdrożenia.
    pub host: String,
    /// Port nasłuchu rdzenia na tym serwerze.
    pub port: u16,
}

impl Wskazane {
    /// Adres HTTP rdzenia złożony ze wskazania. Jedyne miejsce, w którym ten
    /// adres powstaje — nigdzie indziej nie jest wpisywany literałem.
    pub fn adres_http(&self) -> String {
        format!("http://{}:{}", self.host, self.port)
    }
}

/// Komplet ustawień powłoki obowiązujących w chwili pracy.
///
/// Warstwa środowiska jest ustalona raz, przy starcie procesu. Warstwa nastaw
/// żyje dalej — Operator składa wskazanie w oknie już po starcie — więc siedzi
/// za zamkiem i jest wspólna dla wszystkich kopii ustawień (`Arc`). Bez tego
/// polecenie `adres_rdzenia` odpowiadałoby starym adresem do końca pracy
/// procesu, a interfejs łączyłby się nie tam, gdzie Operator wskazał.
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
    /// nastawami zapisanymi. `None` znaczy pierwsze uruchomienie po instalacji —
    /// powłoka nie zgaduje wtedy żadnego adresu, bo każdy zgadnięty byłby
    /// adresem cudzym albo pustym.
    pub fn wskazanie(&self) -> Option<Wskazane> {
        let host = match self.host_ze_srodowiska.as_deref() {
            Some(host) => host.to_string(),
            None => self.nastawy().host_rdzenia?,
        };
        Some(Wskazane {
            host,
            port: self.port(),
        })
    }

    /// Port rdzenia obowiązujący: zmienna środowiska, nastawy zapisane,
    /// a przy braku obu `PORT_DOMYSLNY`.
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
        Warstwa::Brak
    }

    /// Zapisuje wskazanie Operatora trwale i wprowadza je w życie dla wszystkich
    /// kopii ustawień. Zwraca zdanie o niepowodzeniu, gdy zapis się nie udał —
    /// wskazanie nieutrwalone nie zostaje przyjęte, bo zniknęłoby przy następnym
    /// starcie i Operator dowiedziałby się o tym dopiero wtedy.
    ///
    /// Wskazanie ze zmiennej środowiska nie znika przez ten zapis: `wskazanie()`
    /// pyta zmienną pierwszą, a okno dostaje warstwę w odpowiedzi
    /// (`warstwa_wskazania`) i wie, że zapis nie rozstrzyga.
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
            // Zamek zatruty paniką innego wątku: plik jest już zapisany, więc
            // wskazanie obowiązuje od następnego startu. Zdanie mówi dokładnie
            // to, zamiast udawać powodzenie pełne.
            Err(_) => Err(format!(
                "Wskazanie zapisano w pliku {}, ale nie weszło w życie w tym uruchomieniu — \
                 zamknij okno i otwórz je ponownie.",
                nastawy::sciezka().display()
            )),
        }
    }

    /// Kopia nastaw zapisanych. Zamek zatruty daje nastawy puste, nie panikę —
    /// odczyt nastawy nie jest wart przerwania pracy okna.
    fn nastawy(&self) -> Nastawy {
        match self.zapisane.lock() {
            Ok(zamek) => zamek.clone(),
            Err(_) => Nastawy::default(),
        }
    }
}

/// Odczyt portu ze środowiska: wartość niebędąca liczbą nie przerywa startu,
/// tylko zostaje pominięta — rozstrzygają warstwy niższe.
fn port_ze_srodowiska() -> Option<u16> {
    niepusta(ZMIENNA_PORT).and_then(|tekst| tekst.parse().ok())
}

/// Zwraca wartość zmiennej środowiska, traktując wartość pustą jak brak ustawienia.
fn niepusta(nazwa: &str) -> Option<String> {
    match env::var(nazwa) {
        Ok(wartosc) if !wartosc.trim().is_empty() => Some(wartosc.trim().to_string()),
        _ => None,
    }
}
