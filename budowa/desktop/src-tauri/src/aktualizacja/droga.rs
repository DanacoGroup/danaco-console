//! Droga systemowa założenia wydania — czym różni się podmiana na każdym systemie.
//!
//! Rozpoznawane są trzy sytuacje i każda kończy się inaczej:
//!
//! 1. Linux, aplikacja biegnie jako AppImage — plik wskazany zmienną `APPIMAGE`
//!    zostaje podmieniony, a powłoka wstaje ponownie sama.
//!
//! 2. Linux, aplikacja pochodzi z pakietu `.deb` — podmiana jest niemożliwa,
//!    bo binarka leży w `/usr/bin`, gdzie proces Operatora nie ma prawa zapisu,
//!    a powłoka nie będzie sobie tych praw podnosiła. Ta gałąź zwraca odmowę
//!    opisującą brak i mówiącą Operatorowi, co zrobić samemu; treści odmowy
//!    pilnuje test `bez_appimage_odmowa_nazywa_droge_wyjscia`.
//!
//! 3. Windows — wydanie przychodzi jako instalka NSIS, którą powłoka uruchamia
//!    jako osobny program i schodzi jej z drogi.
//!
//! Dlaczego podmiana przez `rename`, a nie przez zapis w miejsce: na Linuksie
//! nie wolno nadpisać pliku wykonywalnego, który właśnie biegnie — jądro zwraca
//! `ETXTBSY`. `rename` w obrębie tego samego katalogu tego zakazu nie łamie:
//! podstawia nowy i-węzeł pod starą nazwę, a proces już uruchomiony spokojnie
//! dopracowuje na starym, który znika dopiero po jego zakończeniu. Stąd plik
//! roboczy pobieramy obok celu, a nie do katalogu tymczasowego — `rename`
//! działa wyłącznie w obrębie jednego systemu plików.

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

/// Nazwa zmiennej, którą uruchamiacz AppImage wskazuje sam siebie.
///
/// To jedyne wiarygodne rozpoznanie: `current_exe()` pokazuje wtedy ścieżkę
/// wewnątrz chwilowo podmontowanego obrazu (`/tmp/.mount_*`), a nie plik,
/// który Operator ma na dysku i który należy podmienić.
#[cfg(not(windows))]
const ZMIENNA_APPIMAGE: &str = "APPIMAGE";

/// Rozpoznaje, którą drogą da się założyć wydanie na tym systemie.
///
/// Odmowa z tej funkcji nigdy nie mówi „nie wolno" — mówi, czego brakuje
/// i co Operator może zrobić sam.
pub fn rozpoznaj() -> Result<Droga, Odmowa> {
    #[cfg(windows)]
    {
        // Instalka NSIS nie podmienia pliku w miejscu — uruchamia się jako
        // osobny program, więc pobieramy ją do katalogu tymczasowego systemu.
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
            // Brak `APPIMAGE` na Linuksie znaczy, że aplikacja pochodzi
            // z pakietu systemowego (`.deb` → `/usr/bin`) albo z budowy
            // deweloperskiej (`target/release`). W obu wypadkach podmiana
            // z wnętrza aplikacji jest niewłaściwa: pakietem zarządza `dpkg`,
            // a budowę deweloperską nadpisuje `cargo`.
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
    /// Ścieżka pliku roboczego, do którego wolno pobierać.
    ///
    /// Dla AppImage leży obok celu (warunek działania `rename`, zob. nagłówek).
    /// Nazwa jest jawnie robocza, żeby nikt nie wziął niesprawdzonego jeszcze
    /// pliku za gotową aplikację.
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

    /// Zakłada sprawdzone już wydanie. Wołać wyłącznie po zgodnej sumie.
    ///
    /// Zwraca zdanie opisujące, co się stało, oraz informację, czy powłoka ma
    /// wstać sama (`restartuje_powloka`), czy oddaje pole instalatorowi.
    pub fn zaloz(&self, plik_roboczy: &std::path::Path) -> Result<Zalozenie, Odmowa> {
        match self {
            #[cfg(not(windows))]
            Droga::AppImage { cel } => {
                // Ostatnie spojrzenie na to, co naprawdę leży na dysku.
                // `zaloz` podmienia plik aplikacji nieodwracalnie, więc pyta
                // o zawartość pliku roboczego sam, zamiast wierzyć, że ktoś
                // wcześniej sprawdził właściwą rzecz. Zgodna suma SHA-256 tego
                // nie zastąpi: plik zerowej długości i plik `.exe` z wydania
                // windowsowego mają sumy poprawne co do znaku, a żaden z nich
                // nie jest aplikacją, którą da się tu uruchomić.
                sprawdz_zawartosc(plik_roboczy)?;

                // Bit wykonywalny nadajemy dopiero teraz — po sprawdzeniu sumy.
                // Plik pobrany, a jeszcze niesprawdzony, nie ma prawa być
                // uruchamialny nawet przez pomyłkę.
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

                // Podmiana właściwa. Do tej chwili wszystko dało się cofnąć;
                // po niej aplikacją jest już nowy plik.
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

            // Instalka NSIS budowana przez Tauri przyjmuje `/S` (przebieg cichy).
            // Powłoka musi zejść z drogi, bo instalator nie podmieni pliku
            // trzymanego przez biegnący proces — dlatego `restartuje_powloka`
            // jest tu fałszem: aplikację z powrotem stawia instalator, nie my.
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

/// Pierwsze bajty pliku wykonywalnego ELF — tym zaczyna się każdy AppImage.
#[cfg(not(windows))]
const ZNAK_ELF: [u8; 4] = [0x7f, b'E', b'L', b'F'];

/// Sprawdza, czy plik roboczy może być aplikacją tego systemu.
///
/// Dwa pytania, oba o brak, żadne o zgodę: czy jest w nim cokolwiek i czy
/// wygląda na plik wykonywalny Linuksa. Sprawdzenie sumy odpowiada wyłącznie
/// na pytanie „czy to ten plik, co w wykazie" — nie na pytanie „czy to w ogóle
/// aplikacja". Wykaz może wskazywać wydanie na inny system albo plik pusty
/// i suma będzie się wtedy zgadzać co do znaku.
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

/// Skutek założenia wydania.
#[derive(Debug)]
pub struct Zalozenie {
    /// Co dokładnie powstało lub zostało podmienione.
    pub zalozone: String,
    /// Czy to powłoka ma się uruchomić ponownie (AppImage), czy oddaje
    /// pole zewnętrznemu instalatorowi (NSIS).
    pub restartuje_powloka: bool,
    /// Zdanie dla Operatora.
    pub zdanie: String,
}

#[cfg(all(test, not(windows)))]
mod testy {
    //! Testy drogi założenia na Linuksie. Cały moduł jest pod
    //! `#[cfg(not(windows))]`, bo sprawdza wariant `AppImage` i funkcję
    //! `sprawdz_zawartosc`, których pod Windowsem po prostu nie ma.
    //!
    //! Gałąź windowsowa nie ma tu testów i mieć ich nie będzie: sprawdzenie
    //! jej wymaga uruchomienia instalatora NSIS na Windowsie, a test udający,
    //! że to robi, byłby atrapą.

    use super::*;
    use crate::aktualizacja::probne::KatalogProbny;
    use std::os::unix::fs::PermissionsExt;
    use std::sync::{Mutex, MutexGuard, PoisonError};

    /// `rozpoznaj()` czyta zmienną środowiska, a zmienne są wspólne dla całego
    /// procesu testowego — bez zamka testy mrugałyby przy równoległym biegu.
    static ZAMEK: Mutex<()> = Mutex::new(());

    fn szereg() -> MutexGuard<'static, ()> {
        ZAMEK.lock().unwrap_or_else(PoisonError::into_inner)
    }

    /// Najkrótsza treść, która przechodzi kontrolę nagłówka — cztery bajty ELF.
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
        // Wykaz wydań wskazał `.exe` zamiast AppImage. Suma by się zgadzała —
        // to przecież suma tego `.exe`. Na Linuksie to nie jest aplikacja.
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
        // Jeden bajt to już nie pustka, ale wciąż nie ELF. Granica ma odciąć
        // także ten przypadek, a nie tylko dokładne zero.
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
        //! Aplikacja ma przeżyć próbę podmiany na pustkę.
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
        // Warunek działania `rename`: ten sam katalog, czyli ten sam system
        // plików. Pobieranie do /tmp zabiłoby podmianę na maszynach, gdzie
        // /tmp jest osobnym systemem plików (EXDEV).
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
        //! Droga `.deb` jest niemożliwa. Test pilnuje nie kodu, lecz treści
        //! komunikatu: odmowa ma nazwać drogę wyjścia, a nie zamienić się
        //! w ciche niepowodzenie ani w pustą odmowę.
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
        // `APPIMAGE=` albo same spacje nie są ścieżką. Bez tego powłoka
        // próbowałaby podmienić plik o nazwie pustej.
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
