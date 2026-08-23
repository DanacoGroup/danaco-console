//! Uchwyt rdzenia pracującego w tle — stan współdzielony powłoki.
//!
//! Cykl życia rdzenia jest rozdzielony od cyklu życia okna: zamknięcie okna nie
//! kończy procesów sesji. Uchwyt nie ubija rdzenia sam z siebie — ani przy
//! zamknięciu okna, ani przy zakończeniu powłoki. Rdzeń zatrzymuje wyłącznie
//! jawne polecenie z zasobnika.
//!
//! Cały stan siedzi pod jednym zamkiem, bo zmienia się w komplecie: wskazanie
//! rdzenia złożone w oknie (`wskazanie.rs`) podmienia naraz host, port, opis
//! przebiegu i proces potomny. Dwa zamki dałyby chwilę, w której opis mówi
//! o jednym rdzeniu, a proces należy do drugiego.

use std::process::Child;
use std::sync::Mutex;

use serde::Serialize;

use super::nasluch;

/// Opis stanu rdzenia przekazywany interfejsowi i oknom dialogowym.
#[derive(Clone, Debug, Serialize)]
pub struct OpisRdzenia {
    /// Czy rdzeń odpowiada na porcie w chwili zapytania.
    pub pracuje: bool,
    /// Identyfikator procesu, gdy rdzeń postawiła powłoka.
    pub pid: Option<u32>,
    /// Czy proces rdzenia postawiła ta powłoka, czy zastała pracujący.
    pub postawiony_przez_powloke: bool,
    /// Adres HTTP rdzenia obowiązujący.
    pub adres: String,
    /// Ścieżka dziennika rdzenia prowadzonego przez powłokę.
    pub dziennik: String,
    /// Zdanie opisujące przebieg uruchomienia — także przy niepowodzeniu.
    pub opis: String,
}

/// Stan rdzenia pod zamkiem.
struct Stan {
    port: u16,
    /// Host rdzenia obowiązujący (`ustawienia.host()`): pętla zwrotna dla
    /// wariantu lokalnego, wskazana nazwa serwera dla wariantu wirtualnego.
    /// Rozpoznanie nasłuchu pyta ten host, a nie zawsze pętlę zwrotną — inaczej
    /// `pracuje` byłoby fałszywe przy rdzeniu zdalnym.
    host: String,
    opis: OpisRdzenia,
    dziecko: Option<Child>,
}

/// Uchwyt rdzenia przechowywany jako stan aplikacji Tauri.
pub struct UchwytRdzenia {
    stan: Mutex<Stan>,
}

impl UchwytRdzenia {
    /// Tworzy uchwyt z opisu uruchomienia i ewentualnego procesu potomnego.
    pub fn nowy(port: u16, host: String, opis: OpisRdzenia, dziecko: Option<Child>) -> Self {
        Self {
            stan: Mutex::new(Stan {
                port,
                host,
                opis,
                dziecko,
            }),
        }
    }

    /// Zwraca opis stanu z odświeżonym rozpoznaniem nasłuchu pod hostem
    /// obowiązującym — w wariancie wirtualnym pyta rdzeń zdalny, nie pętlę zwrotną.
    ///
    /// Zamek zatruty oddaje opis mówiący o tym wprost, zamiast paniki w poleceniu
    /// interfejsu: wskaźnik łączności ma powiedzieć, co się stało.
    pub fn opis(&self) -> OpisRdzenia {
        match self.stan.lock() {
            Ok(zamek) => {
                let mut biezacy = zamek.opis.clone();
                biezacy.pracuje = nasluch::odpowiada_pod(&zamek.host, zamek.port);
                biezacy
            }
            Err(_) => OpisRdzenia {
                pracuje: false,
                pid: None,
                postawiony_przez_powloke: false,
                adres: String::new(),
                dziennik: super::dziennik::sciezka().display().to_string(),
                opis: "Stan rdzenia niedostępny — powłoka nie może go odczytać w tym uruchomieniu."
                    .to_string(),
            },
        }
    }

    /// Podmienia stan rdzenia w komplecie — po złożeniu wskazania Operatora
    /// (`wskazanie::wskaz`). Proces potomny zastany pod zamkiem zostaje
    /// zatrzymany, gdy nowy stan go nie obejmuje: inaczej rdzeń lokalny
    /// postawiony przy poprzednim wskazaniu pracowałby dalej obok rdzenia
    /// wskazanego na serwerze, a Operator miałby dwa stany platformy.
    pub fn przejmij(&self, port: u16, host: String, opis: OpisRdzenia, dziecko: Option<Child>) {
        let Ok(mut zamek) = self.stan.lock() else {
            return;
        };
        if dziecko.is_none() {
            if let Some(poprzedni) = zamek.dziecko.as_mut() {
                let _ = poprzedni.kill();
                let _ = poprzedni.wait();
            }
        }
        *zamek = Stan {
            port,
            host,
            opis,
            dziecko,
        };
    }

    /// Zatrzymuje rdzeń na jawne polecenie z zasobnika. Rdzeń zastany —
    /// postawiony poza powłoką — nie należy do powłoki i nie jest ubijany.
    pub fn zatrzymaj(&self) -> Result<String, String> {
        let mut zamek = self
            .stan
            .lock()
            .map_err(|_| "stan rdzenia niedostępny".to_string())?;
        match zamek.dziecko.as_mut() {
            Some(proces) => {
                let pid = proces.id();
                proces
                    .kill()
                    .map_err(|blad| format!("zatrzymanie rdzenia (pid {pid}): {blad}"))?;
                let _ = proces.wait();
                zamek.dziecko = None;
                Ok(format!("Rdzeń zatrzymany (pid {pid})."))
            }
            None => Err(
                "Rdzeń nie został postawiony przez tę powłokę — zatrzymaj go tam, gdzie go uruchomiono."
                    .to_string(),
            ),
        }
    }
}
