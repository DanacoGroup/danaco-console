//! Ustawienia powłoki obowiązujące w chwili pracy.
//!
//! Trzy warstwy, od najsłabszej: wartość domyślna → nastawy zapisane trwale
//! (`nastawy.rs`) → zmienna środowiska. Brak jakiejkolwiek warstwy nie wstrzymuje
//! startu okna. Nazwy `DANACO_PORT` i `DANACO_KATALOG_DANYCH` są własnością
//! rdzenia (`server/internal/konfiguracja/srodowisko.go`); powłoka je wyłącznie
//! czyta i przekazuje dziecku.
//!
//! Dlaczego zmienna stoi nad plikiem. Plik niesie wskazanie Operatora złożone
//! w oknie i ma przetrwać zamknięcie okna. Zmienna niesie wskazanie tego, kto
//! stawia proces — wykonawcy przy budowie, jednostki usługi na serwerze — i musi
//! brać górę, bo inaczej plik z jednej maszyny sterowałby uruchomieniem na
//! drugiej po skopiowaniu profilu.
//!
//! Dwóch zmiennych nie wolno mylić. `DANACO_HOST_RDZENIA` rozstrzyga, gdzie
//! stoi rdzeń: steruje adresem złożonym w `adres_rdzenia_http` i przełącznikiem
//! `rdzen_lokalny` (czy `main.rs` w ogóle stawia proces lokalny).
//! `DANACO_ADRES_INTERFEJSU` rozstrzyga wyłącznie, skąd ładuje się strona
//! (`zrodlo_interfejsu.rs`) — osobna decyzja, osobny skutek. Wskazanie
//! jednej z nich nie zmienia drugiej.

use std::env;
use std::path::PathBuf;
use std::sync::{Arc, Mutex};

use crate::nastawy::{self, Nastawy};

/// Port nasłuchu rdzenia lokalnego. Musi być równy `PortDomyslny`
/// z `server/internal/konfiguracja/ustawienia.go` (17870) — tamten plik jest
/// źródłem prawdy, ta stała jest jego jedyną kopią po stronie Rust powłoki.
pub const PORT_DOMYSLNY: u16 = 17870;

/// Host rdzenia lokalnego — wartość obowiązująca, gdy Operator nie wskazał ani
/// nastawy, ani zmiennej. Jedyne miejsce literału `127.0.0.1` po stronie Rust
/// powłoki; `adres_rdzenia_http` i testy budują adres z tej stałej, nie
/// wpisują go wprost.
pub const HOST_DOMYSLNY: &str = "127.0.0.1";

/// Zmienna wskazująca port rdzenia — wspólna z rdzeniem i wzorcem `.env.example`.
pub const ZMIENNA_PORT: &str = "DANACO_PORT";
/// Zmienna wskazująca ścieżkę binarki rdzenia, gdy leży poza miejscami znanymi powłoce.
pub const ZMIENNA_RDZEN: &str = "DANACO_RDZEN";
/// Zmienna wskazująca adres interfejsu (serwer rozwojowy albo rdzeń serwujący pakiet).
/// Rozstrzyga, skąd ładuje się strona — nie gdzie stoi rdzeń (zob. `ZMIENNA_HOST_RDZENIA`).
pub const ZMIENNA_ADRES_INTERFEJSU: &str = "DANACO_ADRES_INTERFEJSU";
/// Zmienna wskazująca host, na którym stoi rdzeń — wariant wirtualny (rdzeń na
/// serwerze) zamiast lokalnego (rdzeń na tym urządzeniu). Brak wskazania oddaje
/// rozstrzygnięcie nastawom zapisanym, a przy ich braku `HOST_DOMYSLNY`. Jawne
/// wskazanie hosta innego niż domyślny wyłącza stawianie procesu lokalnego
/// (`Ustawienia::rdzen_lokalny`, `rdzen::uruchomienie::nie_stawiaj_lokalnie`,
/// wołane z `main.rs`).
pub const ZMIENNA_HOST_RDZENIA: &str = "DANACO_HOST_RDZENIA";

/// Nazwa warstwy, z której pochodzi obowiązujące wskazanie hosta. Wchodzi do
/// odpowiedzi polecenia `wskazanie_rdzenia`, żeby okno mogło powiedzieć
/// Operatorowi, dlaczego pola nie da się zmienić: wskazanie ze zmiennej
/// środowiska jest silniejsze od zapisu w oknie.
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub enum Warstwa {
    /// Wskazania nie ma w żadnej warstwie — obowiązuje wartość domyślna.
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
    /// Host wskazany zmienną środowiska; brak = nastawy albo `HOST_DOMYSLNY`.
    host_ze_srodowiska: Option<String>,
    /// Jawnie wskazana binarka rdzenia; brak = wyszukanie w miejscach znanych.
    pub sciezka_rdzenia: Option<PathBuf>,
    /// Jawnie wskazany adres interfejsu; brak = rozstrzygnięcie w `zrodlo_interfejsu`.
    pub adres_interfejsu: Option<String>,
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
            sciezka_rdzenia: niepusta(ZMIENNA_RDZEN).map(PathBuf::from),
            adres_interfejsu: niepusta(ZMIENNA_ADRES_INTERFEJSU),
            zapisane: Arc::new(Mutex::new(nastawy::czytaj())),
        }
    }

    /// Host rdzenia obowiązujący: zmienna środowiska, nastawy zapisane,
    /// a przy braku obu `HOST_DOMYSLNY`. Nigdy literał wpisany wprost.
    pub fn host(&self) -> String {
        if let Some(host) = self.host_ze_srodowiska.as_deref() {
            return host.to_string();
        }
        self.nastawy()
            .host_rdzenia
            .unwrap_or_else(|| HOST_DOMYSLNY.to_string())
    }

    /// Port rdzenia obowiązujący: zmienna środowiska, nastawy zapisane,
    /// a przy braku obu `PORT_DOMYSLNY`.
    pub fn port(&self) -> u16 {
        self.port_ze_srodowiska
            .or_else(|| self.nastawy().port_rdzenia)
            .unwrap_or(PORT_DOMYSLNY)
    }

    /// Adres HTTP rdzenia obowiązujący — złożony z hosta i portu
    /// obowiązujących: host z `host()`, port z `port()`, żaden nie jest
    /// literałem wpisanym we `format!`.
    pub fn adres_rdzenia_http(&self) -> String {
        format!("http://{}:{}", self.host(), self.port())
    }

    /// Czy rdzeń ma stać na tej maszynie. Prawda przy wskazaniu równym
    /// `HOST_DOMYSLNY` i przy braku wskazania. Fałsz wyłącznie, gdy host
    /// obowiązujący jest inny — wariant wirtualny (rdzeń na serwerze).
    ///
    /// Steruje tylko tym, czy `main.rs` stawia proces rdzenia lokalnie. Skąd
    /// ładuje się strona interfejsu, rozstrzyga osobno `adres_interfejsu`
    /// i `zrodlo_interfejsu.rs`; wskazanie adresu interfejsu na hosta zdalnego
    /// samo z siebie nie wyłącza stawiania rdzenia lokalnie.
    pub fn rdzen_lokalny(&self) -> bool {
        self.host() == HOST_DOMYSLNY
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

    /// Czy wskazanie zostało w ogóle złożone — w oknie albo zmienną.
    ///
    /// Fałsz znaczy pierwsze uruchomienie: powłoka nie wie jeszcze, czy rdzeń ma
    /// stanąć na tym urządzeniu, czy stoi na serwerze, więc nie stawia procesu
    /// lokalnego i czeka na rozstrzygnięcie Operatora (`main.rs`).
    pub fn wskazanie_zlozone(&self) -> bool {
        self.warstwa_wskazania() != Warstwa::Brak
    }

    /// Zapisuje wskazanie Operatora trwale i wprowadza je w życie dla wszystkich
    /// kopii ustawień. Zwraca zdanie o niepowodzeniu, gdy zapis się nie udał —
    /// wskazanie nieutrwalone nie zostaje przyjęte, bo zniknęłoby przy następnym
    /// starcie i Operator dowiedziałby się o tym dopiero wtedy.
    ///
    /// Wskazanie ze zmiennej środowiska nie znika przez ten zapis: `host()` pyta
    /// zmienną pierwszą, a okno dostaje warstwę w odpowiedzi
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
