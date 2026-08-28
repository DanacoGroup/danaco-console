//! Droga systemowa założenia wydania: AppImage na Linuksie podmienia się przez
//! `rename`, pakiet `.deb` zwraca odmowę nazywającą brak, a Windows uruchamia
//! instalkę NSIS jako osobny program. Uzasadnienie pełne stoi w dokumentacji.

use std::path::PathBuf;

use super::Odmowa;

/// Rozpoznana droga założenia wydania na tym systemie.
///
/// `Debug` jak przy `Odmowa` i `Przebieg` — żeby nieudany test albo wpis
/// w dzienniku mówił, która droga została rozpoznana, a nie tylko że któraś.
#[derive(Debug)]
pub enum Droga {
    /// Linux, aplikacja biegnie jako AppImage — podmiana pliku wskazanego
    /// zmienną `APPIMAGE`.
    #[cfg(not(windows))]
    AppImage {
        /// Ścieżka pliku `.AppImage`, który jest tą aplikacją.
        cel: PathBuf,
    },
    /// Windows, wydanie przychodzi jako instalka NSIS.
    #[cfg(windows)]
    InstalkaNsis {
        /// Katalog tymczasowy, w którym stanie pobrana instalka.
        katalog_roboczy: PathBuf,
    },
}

/// Nazwa zmiennej, którą uruchamiacz AppImage wskazuje sam siebie — jedyne
/// wiarygodne rozpoznanie, bo `current_exe()` wskazuje wtedy ścieżkę wewnątrz
/// chwilowo podmontowanego obrazu, a nie plik na dysku operatora.
#[cfg(not(windows))]
const ZMIENNA_APPIMAGE: &str = "APPIMAGE";

/// Rozpoznaje, którą drogą da się założyć wydanie na tym systemie.
///
/// Odmowa z tej funkcji nigdy nie mówi „nie wolno" — mówi, czego brakuje
/// i co Operator może zrobić sam.
pub fn rozpoznaj() -> Result<Droga, Odmowa> {
    #[cfg(windows)]
    {
        // Instalka NSIS uruchamia się jako osobny program poza katalogiem aplikacji.
        Ok(Droga::InstalkaNsis {
            katalog_roboczy: std::env::temp_dir(),
        })
    }

    #[cfg(not(windows))]
    {
        match std::env::var(ZMIENNA_APPIMAGE) {
            Ok(sciezka) if !sciezka.trim().is_empty() => Ok(Droga::AppImage {
                cel: PathBuf::from(sciezka.trim()),
            }),
            // Brak APPIMAGE znaczy pakiet .deb albo budowę deweloperską, obie z zewnątrz.
            _ => Err(Odmowa::nowa(
                "droga-niedostepna",
                format!(
                    "Ta kopia Danaco Console nie jest plikiem AppImage (zmienna \
                     {ZMIENNA_APPIMAGE} jest pusta), więc powłoka nie ma czego \
                     podmienić. Jeżeli instalowałeś pakiet .deb, aktualizację \
                     zakłada się poleceniem „sudo dpkg -i” — powłoka nie ma prawa \
                     zapisu w /usr i nie będzie sobie tych praw podnosiła. \
                     Nowy plik pobierzesz ze strony „Pobierz”."
                ),
            )),
        }
    }
}

impl Droga {
    /// Ścieżka pliku roboczego: dla AppImage leży obok celu, nazwa jest jawnie robocza.
    pub fn plik_roboczy(&self) -> PathBuf {
        match self {
            #[cfg(not(windows))]
            Droga::AppImage { cel } => {
                let mut nazwa = cel.file_name().unwrap_or_default().to_os_string();
                nazwa.push(".pobierane");
                cel.with_file_name(nazwa)
            }
            #[cfg(windows)]
            Droga::InstalkaNsis { katalog_roboczy } => {
                katalog_roboczy.join("danaco-console-aktualizacja.exe")
            }
        }
    }

    /// Krótka nazwa drogi dla interfejsu i dziennika.
    pub fn nazwa(&self) -> &'static str {
        match self {
            #[cfg(not(windows))]
            Droga::AppImage { .. } => "appimage",
            #[cfg(windows)]
            Droga::InstalkaNsis { .. } => "instalka-nsis",
        }
    }

    /// Zakłada sprawdzone już wydanie i zwraca zdanie o skutku; wołać wyłącznie po zgodnej sumie.
    pub fn zaloz(&self, plik_roboczy: &std::path::Path) -> Result<Zalozenie, Odmowa> {
        match self {
            #[cfg(not(windows))]
            Droga::AppImage { cel } => {
                // Ostatnie spojrzenie na zawartość — zgodna suma nie zastępuje tego sprawdzenia.
                sprawdz_zawartosc(plik_roboczy)?;

                // Bit wykonywalny nadajemy dopiero po sprawdzeniu sumy, nie wcześniej.
                use std::os::unix::fs::PermissionsExt;
                let prawa = std::fs::Permissions::from_mode(0o755);
                std::fs::set_permissions(plik_roboczy, prawa).map_err(|blad| {
                    Odmowa::nowa(
                        "brak-prawa-zapisu",
                        format!(
                            "Nie udało się nadać nowemu wydaniu prawa uruchamiania: {blad}. \
                             Stara wersja została nietknięta."
                        ),
                    )
                })?;

                // Podmiana właściwa: do tej chwili dało się cofnąć, po niej już nie.
                std::fs::rename(plik_roboczy, cel).map_err(|blad| {
                    let _ = std::fs::remove_file(plik_roboczy);
                    Odmowa::nowa(
                        "podmiana-nieudana",
                        format!(
                            "Nie udało się podmienić pliku {}: {blad}. Stara wersja \
                             została nietknięta, a plik pobrany skasowany.",
                            cel.display()
                        ),
                    )
                })?;

                Ok(Zalozenie {
                    zalozone: cel.display().to_string(),
                    restartuje_powloka: true,
                    zdanie: format!(
                        "Nowe wydanie założone w miejsce {}. Aplikacja uruchomi się ponownie.",
                        cel.display()
                    ),
                })
            }

            // Instalka NSIS przyjmuje /S; restartuje_powloka jest tu fałszem, bo aplikację stawia instalator.
            #[cfg(windows)]
            Droga::InstalkaNsis { .. } => {
                use std::process::Command;
                Command::new(plik_roboczy)
                    .arg("/S")
                    .spawn()
                    .map_err(|blad| {
                        let _ = std::fs::remove_file(plik_roboczy);
                        Odmowa::nowa(
                            "instalator-nie-ruszyl",
                            format!(
                                "Nie udało się uruchomić pobranej instalki: {blad}. \
                                 Stara wersja została nietknięta."
                            ),
                        )
                    })?;
                Ok(Zalozenie {
                    zalozone: plik_roboczy.display().to_string(),
                    restartuje_powloka: false,
                    zdanie: "Instalka aktualizacji została uruchomiona. Aplikacja zamknie \
                             się na czas zakładania nowej wersji."
                        .to_string(),
                })
            }
        }
    }
}

/// Pierwsze bajty pliku wykonywalnego ELF, tym zaczyna się każdy plik AppImage na Linuksie, niezależnie od dystrybucji.
#[cfg(not(windows))]
const ZNAK_ELF: [u8; 4] = [0x7f, b'E', b'L', b'F'];

/// Sprawdza, czy plik roboczy może być aplikacją tego systemu: czy jest w nim
/// cokolwiek i czy wygląda na plik wykonywalny Linuksa. Zgodna suma nie
/// odpowiada na pytanie, czy to w ogóle aplikacja — uzasadnienie w dokumentacji.
#[cfg(not(windows))]
fn sprawdz_zawartosc(plik_roboczy: &std::path::Path) -> Result<(), Odmowa> {
    use std::io::Read;

    let mut plik = std::fs::File::open(plik_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "plik-nieczytelny",
            format!(
                "Nie udało się odczytać pobranego pliku {}: {blad}. \
                 Stara wersja została nietknięta.",
                plik_roboczy.display()
            ),
        )
    })?;
    let mut naglowek = [0_u8; 4];
    let odczytane = plik.read(&mut naglowek).unwrap_or(0);

    if odczytane == 0 {
        let _ = std::fs::remove_file(plik_roboczy);
        return Err(Odmowa::nowa(
            "plik-pusty",
            "Pobrane wydanie ma zerową długość — to nie jest aplikacja. \
             Nic nie zostało podmienione, stara wersja pracuje dalej."
                .to_string(),
        ));
    }
    if odczytane < ZNAK_ELF.len() || naglowek != ZNAK_ELF {
        let _ = std::fs::remove_file(plik_roboczy);
        return Err(Odmowa::nowa(
            "plik-nie-jest-aplikacja",
            format!(
                "Pobrany plik nie jest programem Linuksa (brak nagłówka ELF), \
                 więc powłoka nie podstawi go w miejsce {}. Wykaz wydań mógł \
                 wskazać wydanie na inny system — pobierz właściwe ze strony \
                 „Pobierz”. Stara wersja została nietknięta.",
                plik_roboczy.display()
            ),
        ));
    }
    Ok(())
}

/// Skutek założenia wydania: co dokładnie powstało oraz czy powłoka ma się uruchomić ponownie sama, czy nie.
#[derive(Debug)]
pub struct Zalozenie {
    /// Co dokładnie powstało lub zostało podmienione.
    pub zalozone: String,
    /// Czy to powłoka ma się uruchomić ponownie, czy oddaje pole instalatorowi.
    pub restartuje_powloka: bool,
    /// Zdanie dla Operatora.
    pub zdanie: String,
}

#[cfg(all(test, not(windows)))]
mod testy {
    //! Testy drogi na Linuksie; gałąź windowsowa wymaga instalatora NSIS i testu tu nie ma.

    use super::*;
    use crate::aktualizacja::probne::KatalogProbny;
    use std::os::unix::fs::PermissionsExt;
    use std::sync::{Mutex, MutexGuard, PoisonError};

    /// Zmienne środowiska są wspólne procesowi testowemu, więc dostęp zamyka ten zamek.
    static ZAMEK: Mutex<()> = Mutex::new(());

    fn szereg() -> MutexGuard<'static, ()> {
        ZAMEK.lock().unwrap_or_else(PoisonError::into_inner)
    }

    /// Najkrótsza treść przechodząca kontrolę nagłówka ELF: cztery bajty.
    const NAGLOWEK_ELF: &[u8] = &[0x7f, b'E', b'L', b'F'];

    // ── DOLNA GRANICA PRZY ZAKŁADANIU ──────────────────────────────────────

    #[test]
    fn plik_zerowej_dlugosci_odrzucony_i_skasowany() {
        let katalog = KatalogProbny::nowy("pustka");
        let plik = katalog.plik_z_trescia("wydanie.pobierane", b"");

        let odmowa = sprawdz_zawartosc(&plik)
            .expect_err("plik zerowej długości nie jest aplikacją i nie ma prawa przejść");

        assert_eq!(odmowa.powod, "plik-pusty");
        assert!(
            !plik.exists(),
            "pustka nie może zostać na dysku pod nazwą roboczą obok aplikacji"
        );
    }

    #[test]
    fn plik_windowsowy_odrzucony_jako_nie_aplikacja() {
        // Wykaz wskazał .exe zamiast AppImage; suma się zgadza, ale to nie jest aplikacja Linuksa.
        let katalog = KatalogProbny::nowy("windowsowy");
        let plik = katalog.plik_z_trescia("wydanie.pobierane", b"MZ\x90\x00podszywka");

        let odmowa = sprawdz_zawartosc(&plik)
            .expect_err("plik PE nie ma prawa stanąć w miejscu aplikacji linuksowej");

        assert_eq!(odmowa.powod, "plik-nie-jest-aplikacja");
        assert!(
            odmowa.zdanie.contains("Pobierz"),
            "zdanie ma wskazać Operatorowi, skąd wziąć właściwe wydanie: {}",
            odmowa.zdanie
        );
        assert!(!plik.exists(), "plik nie-aplikacja ma zostać skasowany");
    }

    #[test]
    fn plik_krotszy_niz_naglowek_odrzucony() {
        // Jeden bajt to już nie pustka, ale wciąż nie ELF — granica ma odciąć też ten przypadek.
        let katalog = KatalogProbny::nowy("okruch");
        let plik = katalog.plik_z_trescia("wydanie.pobierane", b"\x7f");

        let odmowa = sprawdz_zawartosc(&plik).expect_err("okruch nie jest aplikacją");
        assert_eq!(odmowa.powod, "plik-nie-jest-aplikacja");
        assert!(!plik.exists());
    }

    #[test]
    fn plik_z_naglowkiem_elf_przechodzi_i_zostaje() {
        let katalog = KatalogProbny::nowy("elf");
        let mut tresc = NAGLOWEK_ELF.to_vec();
        tresc.extend_from_slice(b" reszta wydania");
        let plik = katalog.plik_z_trescia("wydanie.pobierane", &tresc);

        sprawdz_zawartosc(&plik).expect("plik ELF to poprawne wydanie linuksowe");
        assert!(plik.exists(), "poprawnego wydania nie wolno skasować");
    }

    #[test]
    fn plik_nieistniejacy_daje_odmowe_a_nie_panike() {
        let katalog = KatalogProbny::nowy("nieobecny");
        let odmowa = sprawdz_zawartosc(&katalog.plik("nie-ma-mnie"))
            .expect_err("brakującego pliku nie da się sprawdzić");
        assert_eq!(odmowa.powod, "plik-nieczytelny");
    }

    // ── ZAŁOŻENIE WŁAŚCIWE ─────────────────────────────────────────────────

    #[test]
    fn zaloz_pustym_plikiem_nie_tyka_aplikacji() {
        //! Aplikacja ma przeżyć próbę podmiany na plik pusty.
        let katalog = KatalogProbny::nowy("zaloz-pustka");
        let cel = katalog.plik_z_trescia("Danaco.AppImage", b"\x7fELF stara, dzialajaca wersja");
        let droga = Droga::AppImage { cel: cel.clone() };
        let plik_roboczy = katalog.plik_z_trescia("Danaco.AppImage.pobierane", b"");

        let odmowa = droga
            .zaloz(&plik_roboczy)
            .expect_err("pustka nie ma prawa zostać założona");

        assert_eq!(odmowa.powod, "plik-pusty");
        assert_eq!(
            std::fs::read(&cel).expect("aplikacja ma dalej być na miejscu"),
            b"\x7fELF stara, dzialajaca wersja",
            "stara wersja musi zostać nietknięta — po podmianie na pustkę nie ma z czego wrócić"
        );
        assert!(!plik_roboczy.exists(), "pustka ma zniknąć z dysku");
    }

    #[test]
    fn zaloz_podmienia_i_nadaje_prawo_uruchamiania() {
        let katalog = KatalogProbny::nowy("zaloz-udany");
        let cel = katalog.plik_z_trescia("Danaco.AppImage", b"\x7fELF stara wersja");
        std::fs::set_permissions(&cel, std::fs::Permissions::from_mode(0o644))
            .expect("prawa startowe muszą dać się ustawić");
        let droga = Droga::AppImage { cel: cel.clone() };

        let mut nowa = NAGLOWEK_ELF.to_vec();
        nowa.extend_from_slice(b" nowe wydanie");
        let plik_roboczy = katalog.plik_z_trescia("Danaco.AppImage.pobierane", &nowa);

        let zalozenie = droga.zaloz(&plik_roboczy).expect("poprawne wydanie ma się założyć");

        assert!(zalozenie.restartuje_powloka, "AppImage wstaje sam, bez instalatora");
        assert_eq!(zalozenie.zalozone, cel.display().to_string());
        assert_eq!(
            std::fs::read(&cel).expect("nowa wersja ma być pod starą nazwą"),
            nowa
        );
        assert!(
            !plik_roboczy.exists(),
            "po `rename` nazwa robocza nie ma prawa zostać"
        );
        let prawa = std::fs::metadata(&cel).expect("plik istnieje").permissions();
        assert_eq!(
            prawa.mode() & 0o777,
            0o755,
            "bez bitu wykonywalnego nowe wydanie nie wstanie"
        );
    }

    // ── PLIK ROBOCZY ───────────────────────────────────────────────────────

    #[test]
    fn plik_roboczy_lezy_obok_celu_z_przyrostkiem_pobierane() {
        // Rename wymaga tego samego systemu plików; /tmp bywa osobnym systemem (EXDEV).
        let cel = PathBuf::from("/opt/danaco/Danaco Console_1.0.0_amd64.AppImage");
        let droga = Droga::AppImage { cel: cel.clone() };
        let roboczy = droga.plik_roboczy();

        assert_eq!(
            roboczy.parent(),
            cel.parent(),
            "plik roboczy musi leżeć w katalogu celu"
        );
        assert_eq!(
            roboczy.file_name().and_then(|n| n.to_str()),
            Some("Danaco Console_1.0.0_amd64.AppImage.pobierane"),
            "nazwa ma być jawnie robocza, żeby nikt nie wziął jej za aplikację"
        );
        assert_ne!(roboczy, cel, "pobieranie pod nazwę celu niszczyłoby aplikację w locie");
    }

    #[test]
    fn nazwa_drogi_jest_ta_ktora_zna_klient() {
        let droga = Droga::AppImage {
            cel: PathBuf::from("/opt/danaco/Danaco.AppImage"),
        };
        assert_eq!(droga.nazwa(), "appimage");
    }

    // ── UCZCIWOŚĆ KOMUNIKATU DLA .deb ──────────────────────────────────────

    #[test]
    fn bez_appimage_odmowa_nazywa_droge_wyjscia() {
        //! Droga .deb jest niemożliwa; test pilnuje treści komunikatu odmowy.
        let _szereg = szereg();
        let poprzednia = std::env::var(ZMIENNA_APPIMAGE).ok();
        std::env::remove_var(ZMIENNA_APPIMAGE);

        let skutek = rozpoznaj();

        if let Some(wartosc) = poprzednia {
            std::env::set_var(ZMIENNA_APPIMAGE, wartosc);
        }

        let odmowa = match skutek {
            Err(odmowa) => odmowa,
            Ok(_) => panic!("bez zmiennej APPIMAGE nie ma czego podmienić"),
        };

        assert_eq!(odmowa.powod, "droga-niedostepna");
        assert!(
            odmowa.zdanie.contains("dpkg"),
            "Operator na .deb musi usłyszeć, czym zakłada się aktualizację: {}",
            odmowa.zdanie
        );
        assert!(
            odmowa.zdanie.contains("Pobierz"),
            "…i skąd bierze nowy plik: {}",
            odmowa.zdanie
        );
        assert!(
            !odmowa.zdanie.contains("nie wolno"),
            "odmowa ma opisywać brak, nie zakaz: {}",
            odmowa.zdanie
        );
    }

    #[test]
    fn pusta_zmienna_appimage_to_tez_brak_drogi() {
        // Wartość pusta albo same spacje nie są ścieżką pliku.
        let _szereg = szereg();
        let poprzednia = std::env::var(ZMIENNA_APPIMAGE).ok();
        std::env::set_var(ZMIENNA_APPIMAGE, "   ");

        let skutek = rozpoznaj();

        match poprzednia {
            Some(wartosc) => std::env::set_var(ZMIENNA_APPIMAGE, wartosc),
            None => std::env::remove_var(ZMIENNA_APPIMAGE),
        }

        assert_eq!(
            skutek.expect_err("pusta zmienna to brak drogi").powod,
            "droga-niedostepna"
        );
    }

    #[test]
    fn z_appimage_rozpoznana_droga_wskazuje_ten_plik() {
        let _szereg = szereg();
        let poprzednia = std::env::var(ZMIENNA_APPIMAGE).ok();
        std::env::set_var(ZMIENNA_APPIMAGE, "  /opt/danaco/Danaco.AppImage  ");

        let skutek = rozpoznaj();

        match poprzednia {
            Some(wartosc) => std::env::set_var(ZMIENNA_APPIMAGE, wartosc),
            None => std::env::remove_var(ZMIENNA_APPIMAGE),
        }

        match skutek {
            Ok(Droga::AppImage { cel }) => assert_eq!(
                cel,
                PathBuf::from("/opt/danaco/Danaco.AppImage"),
                "spacje wokół ścieżki mają być obcięte, inaczej `rename` chybi"
            ),
            Err(odmowa) => panic!("zmienna wskazuje plik, a dostaliśmy odmowę {}", odmowa.powod),
        }
    }
}
