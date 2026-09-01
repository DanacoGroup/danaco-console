//! Krok 5 kreatora: pobranie powłoki programu z kanału wydań Danaco, sprawdzenie
//! jej sumy kontrolnej i założenie programu pobraną instalką NSIS. Wykaz wydań
//! czytany jest z kanału przy każdym przebiegu — wydanie wychodzi częściej niż
//! instalator, a wykaz wkompilowany wiązałby cykl instalatora z cyklem powłoki.
//! Kopia wkompilowana przy budowie zostaje wyłącznie na wypadek kanału, który nie
//! odpowie; adres kanału niesie stała `ADRES_KANALU`, nie wykaz.

use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use std::io::Read;
use std::path::{Path, PathBuf};
use std::time::Duration;
use tauri::{AppHandle, Emitter};

/// Adres kanału pobrań (`prowadzenie/kanal-wydan.md`), ta sama wartość co
/// `ADRES_KANALU` w `desktop/src-tauri/src/aktualizacja/pobranie.rs`. Stoi
/// w kodzie, nie w wykazie: wykaz czyta się spod tego adresu, więc nie może
/// go nieść.
const ADRES_KANALU: &str = "https://pobierz.danaco-group.pl/";

/// Kopia wykazu wydań z chwili budowy — droga zapasowa na wypadek kanału,
/// który nie odpowie; o pozycjach rozstrzyga dopiero wtedy.
const WYKAZ_WKOMPILOWANY: &str = include_str!("../../../witryna/wydania.json");

/// Nazwa pliku wykazu w katalogu wydawanym kanału — ta sama, pod którą
/// `witryna/zloz.mjs` odkłada wykaz obok stron.
const PLIK_WYKAZU: &str = "wydania.json";

/// Zapora czasu żądań poprzedzających pobranie: sprawdzenia kanału i odczytu
/// wykazu. Bez niej okno czeka tyle, ile trwa zwłoka serwera niedostępnego.
/// Samo pobranie pliku wydania zapory nie ma — trwa tyle, ile trwa łącze.
const ZAPORA_CZASU: Duration = Duration::from_secs(10);

/// Górna granica wykazu czytanego z kanału. Wykaz wydań ma kilkanaście kilobajtów;
/// odpowiedź większa nie jest wykazem i nie ma powodu jej wczytywać.
const GRANICA_WYKAZU: usize = 512 * 1024;

/// Nazwa pliku wykonywalnego powłoki w katalogu programu. Ustala ją instalka NSIS
/// wydania (`MAINBINARYNAME`), a ta bierze ją z pola `name` pakietu
/// `desktop/src-tauri/Cargo.toml` — nie z nazwy produktu. Zmiana tamtej nazwy
/// wymaga zmiany tej stałej; przy polu `name` stoi komentarz zwrotny.
const NAZWA_PROGRAMU: &str = "danaco-console-powloka.exe";

/// Nazwa pliku wskazania katalogu danych, zakładanego w katalogu programu
/// z wyboru kroku 4. Powłoka czyta go obok własnego pliku wykonywalnego
/// (`desktop/src-tauri/src/dziennik.rs`, `NAZWA_WSKAZANIA_KATALOGU` — obie
/// stałe muszą być równe).
const NAZWA_WSKAZANIA_KATALOGU_DANYCH: &str = "powloka-katalog-danych.txt";

/// Poświadczenia kanału pobrań wpisane w postać instalki przy jej składaniu.
/// Katalog `/wydania/` chroni uwierzytelnienie podstawowe, a kreator nie ma
/// gdzie o nie zapytać: żaden z sześciu kroków nie przewiduje takiego pola.
/// Hasła nie ma w repozytorium — wchodzi zmienną środowiska w chwili budowy,
/// tak samo jak adres serwera wdrożenia w powłoce.
const UZYTKOWNIK_KANALU: Option<&str> = option_env!("DANACO_KANAL_UZYTKOWNIK");
const HASLO_KANALU: Option<&str> = option_env!("DANACO_KANAL_HASLO");

/// Składa nagłówek uwierzytelnienia podstawowego, gdy poświadczenia wpisano
/// przy składaniu. Bez nich żądanie idzie tak jak dotąd — i wraca odmową 401
/// nazwaną wprost, zamiast cichego niepowodzenia.
fn naglowek_poswiadczen() -> Option<String> {
    let uzytkownik = UZYTKOWNIK_KANALU?;
    let haslo = HASLO_KANALU?;
    let para = format!("{uzytkownik}:{haslo}");
    Some(format!("Basic {}", base64_podstawowy(para.as_bytes())))
}

/// Zapis base64 bez zależności zewnętrznej: kanał wymaga jednego nagłówka,
/// a dokładanie dla niego biblioteki byłoby kosztem bez pokrycia.
fn base64_podstawowy(dane: &[u8]) -> String {
    const ZNAKI: &[u8; 64] = b"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    let mut wynik = String::with_capacity(dane.len().div_ceil(3) * 4);
    for porcja in dane.chunks(3) {
        let b = [
            porcja[0],
            *porcja.get(1).unwrap_or(&0),
            *porcja.get(2).unwrap_or(&0),
        ];
        let trojka = u32::from(b[0]) << 16 | u32::from(b[1]) << 8 | u32::from(b[2]);
        for i in 0..4 {
            if i <= porcja.len() {
                wynik.push(ZNAKI[(trojka >> (18 - i * 6) & 0x3F) as usize] as char);
            } else {
                wynik.push('=');
            }
        }
    }
    wynik
}

/// Zdarzenie niosące postęp pobierania do okna kreatora.
pub const ZDARZENIE_POSTEP: &str = "instalator:postep-pobrania";

/// Zdarzenie niosące wynik sprawdzenia kanału wydań, wykonanego po zbudowaniu
/// okna — okno stoi, zanim żądanie sieciowe padnie, więc wynik dochodzi do niego
/// zdarzeniem, a nie parametrem adresu.
pub const ZDARZENIE_KANAL: &str = "instalator:stan-kanalu";

/// Wykaz wydań; pozostałe klucze pliku (opis, protokół) nie są tu czytane.
#[derive(Deserialize)]
struct Wykaz {
    wydania: Vec<PozycjaWydania>,
}

#[derive(Deserialize, Clone)]
struct PozycjaWydania {
    wersja: String,
    system: String,
    postac: String,
    architektura: String,
    plik: String,
    #[serde(rename = "nazwaPliku")]
    nazwa_pliku: String,
    #[serde(rename = "rozmiarBajty")]
    rozmiar_bajty: u64,
    suma: String,
}

/// Odmowa niesie powód po angielsku ustalony z zamkniętego słownika i zdanie
/// gotowe do pokazania Operatorowi — bez wymyślonego kodu liczbowego.
#[derive(Serialize, Clone, Debug)]
pub struct Odmowa {
    pub powod: String,
    pub zdanie: String,
}

impl Odmowa {
    fn nowa(powod: &str, zdanie: String) -> Self {
        Self {
            powod: powod.to_string(),
            zdanie,
        }
    }
}

/// Stan kanału wydań oddawany oknu po jego zbudowaniu: kanał odpowiada albo
/// niesie odmowę, którą krok 5 pokaże bez udawanego przebiegu.
#[derive(Serialize, Clone)]
pub struct StanKanalu {
    pub dostepny: bool,
    pub odmowa: Option<Odmowa>,
}

#[derive(Serialize, Clone)]
pub struct PostepPobrania {
    pub odebrano_bajtow: u64,
    pub razem_bajtow: u64,
}

#[derive(Serialize, Clone)]
pub struct WynikPobrania {
    pub plik_instalki: String,
    pub nazwa_pliku: String,
    pub bajtow: u64,
    pub suma_sha256: String,
    pub wersja: String,
}

#[derive(Serialize, Clone)]
pub struct WynikZalozenia {
    pub katalog_programu: String,
    pub sciezka_programu: String,
    /// Katalog danych zapisany dla powłoki; brak, gdy krok 4 go nie podał.
    pub katalog_danych: Option<String>,
}

/// Klient HTTP kroku 5. Zapora czasu jest podawana osobno, bo pobranie pliku
/// wydania nie może jej mieć, a żądania poprzedzające pobranie muszą.
fn klient(zapora: Option<Duration>) -> ureq::Agent {
    ureq::Agent::config_builder()
        .http_status_as_error(false)
        .timeout_global(zapora)
        .build()
        .into()
}

/// Wczytuje wykaz z kopii wkompilowanej przy budowie.
fn wykaz_wkompilowany() -> Result<Wykaz, Odmowa> {
    serde_json::from_str(WYKAZ_WKOMPILOWANY).map_err(|blad| {
        Odmowa::nowa(
            "wykaz-nieczytelny",
            format!("Wykaz wydań wkompilowany w instalator nie jest poprawnym JSON-em: {blad}."),
        )
    })
}

/// Czyta wykaz wydań spod kanału pobrań.
fn wykaz_z_kanalu() -> Result<Wykaz, Odmowa> {
    let adres = format!("{}/{PLIK_WYKAZU}", ADRES_KANALU.trim_end_matches('/'));
    let zadanie = klient(Some(ZAPORA_CZASU)).get(&adres);
    let zadanie = match naglowek_poswiadczen() {
        Some(naglowek) => zadanie.header("Authorization", &naglowek),
        None => zadanie,
    };
    let mut odpowiedz = zadanie.call().map_err(|blad| {
        Odmowa::nowa(
            "wykaz-nieosiagalny",
            format!("Nie udało się pobrać wykazu wydań z {adres}: {blad}."),
        )
    })?;
    let status = odpowiedz.status();
    if !status.is_success() {
        return Err(Odmowa::nowa(
            "wykaz-nieosiagalny",
            format!("Kanał wydań odpowiedział na {adres} kodem {status}."),
        ));
    }
    let tresc = odpowiedz
        .body_mut()
        .with_config()
        .limit(GRANICA_WYKAZU as u64)
        .read_to_string()
        .map_err(|blad| {
            Odmowa::nowa(
                "wykaz-nieosiagalny",
                format!("Nie udało się wczytać wykazu wydań z {adres}: {blad}."),
            )
        })?;
    serde_json::from_str(&tresc).map_err(|blad| {
        Odmowa::nowa(
            "wykaz-nieczytelny",
            format!("Wykaz wydań spod {adres} nie jest poprawnym JSON-em: {blad}."),
        )
    })
}

/// Wykaz obowiązujący w tym przebiegu: z kanału, a gdy kanał nie odpowie albo
/// odda treść nieczytelną — z kopii wkompilowanej.
fn wykaz_biezacy() -> Result<Wykaz, Odmowa> {
    wykaz_z_kanalu().or_else(|_| wykaz_wkompilowany())
}

/// Wybiera pozycję wykazu zgodną z architekturą wskazaną w kroku 3
/// (`x64` → x64 wykazu, `arm` → ARM64 wykazu). Architektura nierozpoznana
/// (`brak`) nie ma czego dopasować — to nazwana odmowa, nie zgadywanie.
fn dobierz_pozycje(wykaz: &Wykaz, architektura_kroku3: &str) -> Result<PozycjaWydania, Odmowa> {
    let architektura_wykazu = match architektura_kroku3 {
        "x64" => "x64",
        "arm" => "ARM64",
        _ => {
            return Err(Odmowa::nowa(
                "architektura-nierozpoznana",
                "Architektura procesora tej maszyny nie została rozpoznana, a krok trzeci \
                 nie niesie wyboru wersji — wykaz wydań nie ma z czym jej dopasować."
                    .to_string(),
            ))
        }
    };

    let pozycja = wykaz
        .wydania
        .iter()
        // Instalator ściąga POWŁOKĘ programu, nie siebie: wykaz niesie obie
        // postaci pod tym samym systemem i architekturą, więc pozycję
        // rozstrzyga `postac`, a nie kolejność w wykazie.
        .find(|w| {
            w.system == "Windows" && w.architektura == architektura_wykazu && w.postac == "hybryda"
        })
        .cloned()
        .ok_or_else(|| {
            Odmowa::nowa(
                "wydanie-nieopublikowane",
                format!(
                    "Wykaz wydań nie niesie dziś powłoki Windows/{architektura_wykazu} — \
                     nie ma czego pobrać dla tej maszyny."
                ),
            )
        })?;

    // Wykaz przychodzi z sieci, więc reguła protokołu z samego wykazu
    // (`protokol.dopuszczone`) musi obowiązywać także tutaj: pozycja spod
    // adresu innego niż https nie ma czym potwierdzić, skąd pochodzi.
    if !pozycja.plik.starts_with("https://") {
        return Err(Odmowa::nowa(
            "adres-nie-https",
            format!(
                "Pozycja wykazu wskazuje adres {} — instalator pobiera wyłącznie po https.",
                pozycja.plik
            ),
        ));
    }
    Ok(pozycja)
}

/// Sprawdza, bez pobierania pliku, czy serwer wydań odda pozycję dobraną do
/// architektury tej maszyny — jedno realne żądanie HTTPS do rzeczywistego
/// adresu z wykazu wydań. Wołane po zbudowaniu okna, żeby krok 5 wiedział
/// z góry, czy ma pokazać przebieg pobrania, czy od razu nazwaną odmowę.
pub fn sprawdz_wstepnie(architektura: &str) -> Result<(), Odmowa> {
    let pozycja = dobierz_pozycje(&wykaz_biezacy()?, architektura)?;

    let zadanie = klient(Some(ZAPORA_CZASU)).get(&pozycja.plik);
    let zadanie = match naglowek_poswiadczen() {
        Some(naglowek) => zadanie.header("Authorization", &naglowek),
        None => zadanie,
    };
    let odpowiedz = zadanie.call().map_err(|blad| {
        Odmowa::nowa(
            "brak-lacznosci",
            format!("Nie udało się połączyć z {}: {blad}.", pozycja.plik),
        )
    })?;

    let status = odpowiedz.status();
    if !status.is_success() {
        return Err(odmowa_serwera(&pozycja.plik, status, &odpowiedz));
    }
    Ok(())
}

/// Odmowa serwera wydań nazwana kodem odpowiedzi wraz z nagłówkiem żądania
/// poświadczeń, gdy serwer go podał — to on rozstrzyga, czy chodzi o dostęp,
/// czy o brak pliku.
fn odmowa_serwera<T>(
    adres: &str,
    kod: ureq::http::StatusCode,
    odpowiedz: &ureq::http::Response<T>,
) -> Odmowa {
    let www_authenticate = odpowiedz
        .headers()
        .get("www-authenticate")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("");
    Odmowa::nowa(
        "odpowiedz-serwera",
        format!(
            "Serwer wydań {adres} odpowiedział kodem {kod}{}. Kanał pobrań \
             pozycji wydania chroni uwierzytelnienie podstawowe (plik \
             poświadczeń poza repozytorium) — ten instalator nie ma \
             poświadczeń dostępu.",
            if www_authenticate.is_empty() {
                String::new()
            } else {
                format!(" ({www_authenticate})")
            }
        ),
    )
}

/// Pobiera plik wydania dobrany do wersji wskazanej w kroku 3 do katalogu
/// roboczego i sprawdza jego sumę SHA-256 z wykazu. Zgłasza postęp zdarzeniami
/// do okna. Samo pobranie programu nie zakłada — robi to `zaloz_program`.
#[tauri::command]
pub async fn pobierz_skladniki(
    aplikacja: AppHandle,
    architektura: String,
    katalog_roboczy: String,
) -> Result<WynikPobrania, Odmowa> {
    tauri::async_runtime::spawn_blocking(move || {
        wykonaj_pobranie(&aplikacja, &architektura, Path::new(&katalog_roboczy))
    })
    .await
    .unwrap_or_else(|blad| {
        Err(Odmowa::nowa(
            "przebieg-przerwany",
            format!("Pobieranie przerwało się nieoczekiwanie: {blad}."),
        ))
    })
}

fn wykonaj_pobranie(
    aplikacja: &AppHandle,
    architektura: &str,
    katalog_roboczy: &Path,
) -> Result<WynikPobrania, Odmowa> {
    let pozycja = dobierz_pozycje(&wykaz_biezacy()?, architektura)?;

    std::fs::create_dir_all(katalog_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "brak-prawa-zapisu",
            format!(
                "Nie udało się przygotować katalogu roboczego {}: {blad}.",
                katalog_roboczy.display()
            ),
        )
    })?;
    let plik_roboczy: PathBuf = katalog_roboczy.join(&pozycja.nazwa_pliku);

    let zadanie = klient(None).get(&pozycja.plik);
    let zadanie = match naglowek_poswiadczen() {
        Some(naglowek) => zadanie.header("Authorization", &naglowek),
        None => zadanie,
    };
    let mut odpowiedz = zadanie.call().map_err(|blad| {
        Odmowa::nowa(
            "brak-lacznosci",
            format!("Nie udało się połączyć z {}: {blad}.", pozycja.plik),
        )
    })?;

    let status = odpowiedz.status();
    if !status.is_success() {
        return Err(odmowa_serwera(&pozycja.plik, status, &odpowiedz));
    }

    let mut plik = std::fs::File::create(&plik_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "brak-prawa-zapisu",
            format!(
                "Nie udało się utworzyć pliku roboczego {}: {blad}.",
                plik_roboczy.display()
            ),
        )
    })?;

    let mut haszujacy = Sha256::new();
    let mut odebrano: u64 = 0;
    let mut bufor = [0u8; 64 * 1024];
    let mut zrodlo = odpowiedz.body_mut().as_reader();
    loop {
        let n = zrodlo.read(&mut bufor).map_err(|blad| {
            Odmowa::nowa(
                "polaczenie-przerwane",
                format!("Połączenie przerwało się w trakcie pobierania: {blad}."),
            )
        })?;
        if n == 0 {
            break;
        }
        haszujacy.update(&bufor[..n]);
        std::io::Write::write_all(&mut plik, &bufor[..n]).map_err(|blad| {
            Odmowa::nowa(
                "brak-prawa-zapisu",
                format!("Nie udało się zapisać pliku roboczego: {blad}."),
            )
        })?;
        odebrano += n as u64;
        let _ = aplikacja.emit(
            ZDARZENIE_POSTEP,
            PostepPobrania {
                odebrano_bajtow: odebrano,
                razem_bajtow: pozycja.rozmiar_bajty,
            },
        );
    }

    let suma_policzona = format!("{:x}", haszujacy.finalize());
    if suma_policzona != pozycja.suma {
        let _ = std::fs::remove_file(&plik_roboczy);
        return Err(Odmowa::nowa(
            "suma-niezgodna",
            format!(
                "Suma kontrolna się nie zgadza: wykaz wydań zapowiadał {}, a pobrany \
                 plik ma {suma_policzona}. Plik został skasowany.",
                pozycja.suma
            ),
        ));
    }

    Ok(WynikPobrania {
        plik_instalki: plik_roboczy.display().to_string(),
        nazwa_pliku: pozycja.nazwa_pliku,
        bajtow: odebrano,
        suma_sha256: suma_policzona,
        wersja: pozycja.wersja,
    })
}

/// Zakłada program pobraną instalką NSIS: uruchamia ją cicho w katalogu z kroku 4
/// i czeka na kod wyjścia. Kod niezerowy jest odmową kroku 5 — bez zera od instalki
/// na dysku nie ma programu i okno nie ma czego ogłosić na kroku 6. Katalog danych
/// z kroku 4, gdy podany, zostaje zapisany obok programu dla powłoki.
#[tauri::command]
pub async fn zaloz_program(
    plik_instalki: String,
    katalog_programu: String,
    katalog_danych: Option<String>,
) -> Result<WynikZalozenia, Odmowa> {
    tauri::async_runtime::spawn_blocking(move || {
        wykonaj_zalozenie(
            Path::new(&plik_instalki),
            Path::new(&katalog_programu),
            katalog_danych.as_deref(),
        )
    })
    .await
    .unwrap_or_else(|blad| {
        Err(Odmowa::nowa(
            "przebieg-przerwany",
            format!("Zakładanie programu przerwało się nieoczekiwanie: {blad}."),
        ))
    })
}

fn wykonaj_zalozenie(
    plik: &Path,
    katalog: &Path,
    katalog_danych: Option<&str>,
) -> Result<WynikZalozenia, Odmowa> {
    if !plik.is_file() {
        return Err(Odmowa::nowa(
            "brak-instalki",
            format!(
                "Pobranej instalki nie ma pod {} — nie ma czym założyć programu.",
                plik.display()
            ),
        ));
    }

    let kod = polecenie_instalki(plik, katalog)?
        .spawn()
        .map_err(|blad| {
            Odmowa::nowa(
                "instalka-nie-ruszyla",
                format!(
                    "Nie udało się uruchomić pobranej instalki {}: {blad}.",
                    plik.display()
                ),
            )
        })?
        .wait()
        .map_err(|blad| {
            Odmowa::nowa(
                "instalka-bez-kodu",
                format!(
                    "Nie udało się doczekać końca instalki {}: {blad}.",
                    plik.display()
                ),
            )
        })?;

    if !kod.success() {
        return Err(Odmowa::nowa(
            "instalka-odmowila",
            match kod.code() {
                Some(numer) => format!(
                    "Instalka wydania zakończyła się kodem {numer}. Program nie został \
                     założony w {}, a pobrany plik {} pozostał na dysku.",
                    katalog.display(),
                    plik.display()
                ),
                None => format!(
                    "Instalka wydania została przerwana przez system. Program nie został \
                     założony w {}, a pobrany plik {} pozostał na dysku.",
                    katalog.display(),
                    plik.display()
                ),
            },
        ));
    }

    let program = katalog.join(NAZWA_PROGRAMU);
    if !program.is_file() {
        return Err(Odmowa::nowa(
            "program-nieodnaleziony",
            format!(
                "Instalka zakończyła się powodzeniem, ale pliku {} nie ma na dysku — \
                 program nie stanął w katalogu wskazanym w kroku czwartym.",
                program.display()
            ),
        ));
    }

    let katalog_danych = zapisz_wskazanie_katalogu_danych(katalog, katalog_danych)?;

    // Pobrana instalka po założeniu programu nie jest już do niczego potrzebna,
    // a leży w katalogu programu; odmowa jej skasowania nie unieważnia założenia.
    let _ = std::fs::remove_file(plik);

    Ok(WynikZalozenia {
        katalog_programu: katalog.display().to_string(),
        sciezka_programu: program.display().to_string(),
        katalog_danych,
    })
}

/// Zapisuje katalog danych z kroku 4 w pliku obok programu. Brak wskazania
/// albo wskazanie puste nie zakłada pliku — powłoka bierze wtedy katalog
/// domyślny (`stan_maszyny::katalog_danych_powloki`).
fn zapisz_wskazanie_katalogu_danych(
    katalog_programu: &Path,
    katalog_danych: Option<&str>,
) -> Result<Option<String>, Odmowa> {
    let Some(katalog_danych) = katalog_danych.map(str::trim).filter(|k| !k.is_empty()) else {
        return Ok(None);
    };
    let plik = katalog_programu.join(NAZWA_WSKAZANIA_KATALOGU_DANYCH);
    std::fs::write(&plik, format!("{katalog_danych}\n")).map_err(|blad| {
        Odmowa::nowa(
            "katalog-danych-niezapisany",
            format!(
                "Program stanął, ale wskazania katalogu danych nie udało się zapisać \
                 w {}: {blad}. Powłoka użyłaby katalogu domyślnego zamiast wskazanego \
                 w kroku czwartym.",
                plik.display()
            ),
        )
    })?;
    Ok(Some(katalog_danych.to_string()))
}

/// Polecenie uruchamiające instalkę NSIS cicho i w katalogu wskazanym w kroku 4.
#[cfg(windows)]
fn polecenie_instalki(plik: &Path, katalog: &Path) -> Result<std::process::Command, Odmowa> {
    use std::os::windows::process::CommandExt;

    let mut polecenie = std::process::Command::new(plik);
    polecenie.arg("/S");
    // NSIS czyta `/D=` wprost z linii poleceń i przyjmuje ją wyłącznie bez
    // cudzysłowów oraz jako argument ostatni. `arg` ująłby ścieżkę ze spacją
    // w cudzysłów, więc ten jeden argument idzie linią surową.
    polecenie.raw_arg(format!("/D={}", katalog.display()));
    Ok(polecenie)
}

/// Instalka wydania jest plikiem wykonywalnym Windows; poza Windows krok 5
/// nie ma czym założyć programu i mówi to wprost.
#[cfg(not(windows))]
fn polecenie_instalki(_plik: &Path, _katalog: &Path) -> Result<std::process::Command, Odmowa> {
    Err(Odmowa::nowa(
        "system-nieobslugiwany",
        "Instalka wydania jest plikiem wykonywalnym Windows — na tym systemie \
         nie ma jej czym uruchomić."
            .to_string(),
    ))
}

/// Uruchamia program założony w kroku 5 — czynność domyślna kroku 6. Ścieżka
/// pochodzi z założenia, nie z katalogu zgadywanego drugi raz.
#[tauri::command]
pub fn uruchom_program(sciezka_programu: String) -> Result<(), Odmowa> {
    let sciezka = PathBuf::from(&sciezka_programu);
    if !sciezka.is_file() {
        return Err(Odmowa::nowa(
            "program-nieodnaleziony",
            format!("Pliku {sciezka_programu} nie ma na dysku — nie ma czego uruchomić."),
        ));
    }
    // Katalog bieżący procesu instalatora nie jest katalogiem programu, a powłoka
    // składa ścieżki swoich zasobów względem katalogu, w którym stoi.
    let katalog = sciezka.parent().unwrap_or_else(|| Path::new("."));
    std::process::Command::new(&sciezka)
        .current_dir(katalog)
        .spawn()
        .map(|_| ())
        .map_err(|blad| {
            Odmowa::nowa(
                "program-nie-ruszyl",
                format!("Nie udało się uruchomić {sciezka_programu}: {blad}."),
            )
        })
}
