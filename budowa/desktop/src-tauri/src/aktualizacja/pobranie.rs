//! Pobranie pliku wydania i sprawdzenie jego sumy kontrolnej: jedyne miejsce
//! w powłoce, w którym dane z sieci stają się plikiem uruchamianym na maszynie
//! użytkownika. Uzasadnienie kolejności kroków stoi w dokumentacji.

use std::fs::{self, File};
use std::io::{BufWriter, Read, Write};
use std::path::{Path, PathBuf};

use sha2::{Digest, Sha256};

use super::Odmowa;

/// Adres kanału pobrań, kopia `kanal.adres` z wykazu wydań. Powłoka pobiera
/// wydania wyłącznie spod tego adresu; pełne uzasadnienie w dokumentacji.
const ADRES_KANALU: &str = "https://pobierz.danaco-group.pl/";

/// Górny pułap wielkości pliku wydania (512 MiB), zabezpieczenie przed
/// zapełnieniem dysku odpowiedzią bez końca z serwera zepsutego albo podstawionego.
const PULAP_BAJTOW: u64 = 512 * 1024 * 1024;

/// Wynik udanego pobrania: plik roboczy leży już na dysku i ma sumę zgodną z tą żądaną w wykazie wydań.
pub struct Pobrany {
    /// Ścieżka pliku roboczego, jeszcze nie założonego pod nazwą celu.
    pub sciezka: PathBuf,
    /// Liczba zapisanych bajtów.
    pub bajtow: u64,
    /// Suma policzona z zapisanego pliku; równa żądanej, inaczej byłaby odmowa.
    pub suma_sha256: String,
}

/// Pobiera plik spod `adres` do `plik_roboczy` i sprawdza jego sumę SHA-256.
/// Łańcuch zaufania ma trzy ogniwa: kanał pobrań, HTTPS i sama suma; pełne
/// uzasadnienie stoi w dokumentacji architektury zaplecza.
pub fn pobierz_i_sprawdz(
    adres: &str,
    suma_zadana: &str,
    plik_roboczy: &Path,
) -> Result<Pobrany, Odmowa> {
    let adres = adres.trim();
    let suma_zadana = sprawdz_zadanie(adres, suma_zadana)?;

    // Bez tej opcji ureq zwijałby odpowiedź 404 do błędu łączności zamiast statusu.
    let klient: ureq::Agent = ureq::Agent::config_builder()
        .http_status_as_error(false)
        .build()
        .into();

    let mut odpowiedz = klient.get(adres).call().map_err(|blad| {
        Odmowa::nowa(
            "brak-lacznosci",
            format!(
                "Nie udało się połączyć z {adres}: {blad}. Aktualizacja nie została \
                 pobrana — sprawdź łączność z siecią i spróbuj ponownie."
            ),
        )
    })?;

    let status = odpowiedz.status();
    if !status.is_success() {
        return Err(Odmowa::nowa(
            "odpowiedz-serwera",
            format!(
                "Serwer wydań odpowiedział kodem {status} na {adres}. Aktualizacja \
                 nie została pobrana — plik wydania może jeszcze nie być wystawiony."
            ),
        ));
    }

    // Próba utworzenia pliku roboczego sprawdza katalog, zanim przyjdzie pierwszy bajt.
    let plik = File::create(plik_roboczy).map_err(|blad| {
        Odmowa::nowa(
            "brak-prawa-zapisu",
            format!(
                "Nie udało się utworzyć pliku roboczego {}: {blad}. Aktualizacja nie \
                 została pobrana — katalog, w którym powłoka jest zainstalowana, nie \
                 istnieje albo powłoka nie ma w nim prawa zapisu.",
                plik_roboczy.display()
            ),
        )
    })?;

    // Zapis i liczenie sumy idą w jednym przebiegu, kawałkami po 64 KiB.
    let mut zapis = BufWriter::new(plik);
    let mut skrot = Sha256::new();
    let mut czytnik = odpowiedz.body_mut().as_reader();
    let mut kawalek = vec![0_u8; 64 * 1024];
    let mut bajtow: u64 = 0;

    loop {
        let ile = match czytnik.read(&mut kawalek) {
            Ok(0) => break,
            Ok(ile) => ile,
            Err(blad) => {
                return Err(posprzataj(
                    plik_roboczy,
                    Odmowa::nowa(
                        "pobieranie-przerwane",
                        format!(
                            "Pobieranie aktualizacji przerwało się w połowie: {blad}. \
                             Plik niepełny został skasowany, nic nie zostało założone."
                        ),
                    ),
                ));
            }
        };
        bajtow += ile as u64;
        if bajtow > PULAP_BAJTOW {
            return Err(posprzataj(
                plik_roboczy,
                Odmowa::nowa(
                    "plik-ponad-pulap",
                    format!(
                        "Plik wydania przekroczył pułap {} MiB i pobieranie zostało \
                         przerwane. Nic nie zostało założone.",
                        PULAP_BAJTOW / 1024 / 1024
                    ),
                ),
            ));
        }
        skrot.update(&kawalek[..ile]);
        if let Err(blad) = zapis.write_all(&kawalek[..ile]) {
            return Err(posprzataj(
                plik_roboczy,
                Odmowa::nowa(
                    "zapis-nieudany",
                    format!(
                        "Zapis pobieranej aktualizacji nie powiódł się: {blad}. \
                         Sprawdź, czy na dysku jest miejsce. Nic nie zostało założone."
                    ),
                ),
            ));
        }
    }

    // Bufor musi trafić na dysk przed oceną wyniku, nie tylko do pamięci.
    if let Err(blad) = zapis
        .flush()
        .and_then(|()| zapis.into_inner().unwrap().sync_all())
    {
        return Err(posprzataj(
            plik_roboczy,
            Odmowa::nowa(
                "zapis-nieudany",
                format!(
                    "Nie udało się dokończyć zapisu aktualizacji na dysk: {blad}. \
                     Nic nie zostało założone."
                ),
            ),
        ));
    }

    let suma_policzona = format!("{:x}", skrot.finalize());

    // Warunek dopuszczenia: aktualizacja albo przechodzi, albo znika z dysku.
    if let Err(odmowa) = ocen_pobrane(adres, bajtow, &suma_policzona, &suma_zadana) {
        return Err(posprzataj(plik_roboczy, odmowa));
    }

    Ok(Pobrany {
        sciezka: plik_roboczy.to_path_buf(),
        bajtow,
        suma_sha256: suma_policzona,
    })
}

/// Sprawdza samo żądanie, zanim cokolwiek poleci przez sieć i zanim powstanie
/// plik roboczy; zwraca sumę sprowadzoną do małych liter. Wydzielone
/// z `pobierz_i_sprawdz`, żeby test mógł je wywołać bez dostępu do sieci.
fn sprawdz_zadanie(adres: &str, suma_zadana: &str) -> Result<String, Odmowa> {
    if !adres.starts_with("https://") {
        return Err(Odmowa::nowa(
            "adres-nie-https",
            format!(
                "Adres wydania nie zaczyna się od „https://” ({adres}). \
                 Aktualizacja nie została pobrana, bo po zwykłym HTTP nie da się \
                 odróżnić pliku wydawcy od podstawionego."
            ),
        ));
    }

    // Adres musi leżeć w kanale pobrań, inaczej powłoka ufałaby dowolnemu adresowi strony.
    if !adres.starts_with(ADRES_KANALU) {
        return Err(Odmowa::nowa(
            "adres-poza-kanalem",
            format!(
                "Adres wydania ({adres}) nie leży w kanale pobrań {ADRES_KANALU} — \
                 jedynym miejscu, z którego wykaz wydań wskazuje pliki. \
                 Aktualizacja nie została pobrana: plik spoza kanału nie pochodzi \
                 od wydawcy, więc nie ma czym potwierdzić, co by przyszło z sieci."
            ),
        ));
    }

    // Sprawdzenie zapisu sumy idzie przed pobraniem, żeby nie czekać na odmowę.
    let suma_zadana = suma_zadana.trim().to_ascii_lowercase();
    if suma_zadana.len() != 64 || !suma_zadana.chars().all(|z| z.is_ascii_hexdigit()) {
        return Err(Odmowa::nowa(
            "suma-w-zlym-zapisie",
            "Wykaz wydań nie podał poprawnej sumy SHA-256 (oczekiwane 64 znaki \
             szesnastkowe). Aktualizacja nie została pobrana — bez sumy nie ma \
             czym sprawdzić, co przyszło z sieci."
                .to_string(),
        ));
    }

    Ok(suma_zadana)
}

/// Ocenia plik, który wylądował na dysku, wołane po `sync_all`. Kolejność
/// pytań jest warunkiem poprawności: najpierw dolna granica wielkości,
/// dopiero potem suma — uzasadnienie pełne stoi w dokumentacji architektury.
fn ocen_pobrane(
    adres: &str,
    bajtow: u64,
    suma_policzona: &str,
    suma_zadana: &str,
) -> Result<(), Odmowa> {
    if bajtow == 0 {
        return Err(Odmowa::nowa(
            "plik-pusty",
            format!(
                "Serwer oddał pod {adres} plik pustej długości (0 bajtów). \
                 Aktualizacja nie została założona — pustego pliku nie da się \
                 uruchomić, a podmiana aplikacji na niego byłaby nieodwracalna."
            ),
        ));
    }

    if suma_policzona != suma_zadana {
        return Err(Odmowa::nowa(
            "suma-niezgodna",
            format!(
                "SUMA KONTROLNA SIĘ NIE ZGADZA — aktualizacja ODRZUCONA i skasowana. \
                 Wykaz wydań zapowiadał {suma_zadana}, a pobrany plik ma \
                 {suma_policzona}. Nic nie zostało zainstalowane. Plik mógł się \
                 uszkodzić po drodze albo nie pochodzi od wydawcy — nie uruchamiaj \
                 go ręcznie."
            ),
        ));
    }

    Ok(())
}

/// Kasuje plik roboczy i przepuszcza odmowę dalej; każde wyjście błędem po
/// utworzeniu pliku przechodzi tędy, żeby plik niesprawdzony nie został na dysku.
fn posprzataj(plik_roboczy: &Path, odmowa: Odmowa) -> Odmowa {
    let _ = fs::remove_file(plik_roboczy);
    odmowa
}

#[cfg(test)]
mod testy {
    //! Testy sprawdzenia żądania biegną bez sieci; wykaz wydań jest wczytany przy kompilacji.

    use super::*;

    /// Wykaz wydań wczytany przy kompilacji; zmiana pliku wymusza rekompilację modułu.
    const WYKAZ_WYDAN: &str = include_str!("../../../../witryna/wydania.json");

    /// Wykaz wydań w postaci drzewa; niepoprawny JSON jest tu porażką testu.
    fn wykaz() -> serde_json::Value {
        serde_json::from_str(WYKAZ_WYDAN).expect("wykaz wydań ma być poprawnym JSON-em")
    }

    /// Pole `kanal.adres` wykazu — adres, spod którego wykaz wskazuje pliki.
    fn adres_kanalu_z_wykazu() -> String {
        wykaz()["kanal"]["adres"]
            .as_str()
            .expect("wykaz wydań ma nieść `kanal.adres` jako napis")
            .to_string()
    }

    /// Adresy wszystkich pozycji wykazu; wykaz pusty jest tu porażką testu, nie wynikiem.
    fn adresy_wydan_z_wykazu() -> Vec<String> {
        let wykaz = wykaz();
        let pozycje = wykaz["wydania"]
            .as_array()
            .expect("wykaz wydań ma nieść tablicę `wydania`")
            .iter()
            .map(|pozycja| {
                pozycja["plik"]
                    .as_str()
                    .expect("każda pozycja wykazu ma nieść `plik` jako napis")
                    .to_string()
            })
            .collect::<Vec<_>>();
        assert!(
            !pozycje.is_empty(),
            "wykaz wydań nie ma ani jednej pozycji — nie ma czego sprawdzić"
        );
        pozycje
    }

    /// Poprawna co do zapisu suma SHA-256; testom adresu wystarczy sam zapis.
    fn suma_poprawna() -> String {
        "a".repeat(64)
    }

    #[test]
    fn kanal_powloki_zgadza_sie_z_wykazem_wydan() {
        // Sprawdzian wiążący dwie kopie jednego adresu, tutaj i w wykazie wydań.
        assert_eq!(
            ADRES_KANALU,
            adres_kanalu_z_wykazu(),
            "kanał wkompilowany w powłokę rozszedł się z `kanal.adres` \
             w budowa/witryna/wydania.json — powłoka odmówi każdego wydania \
             wskazanego przez wykaz jako leżącego poza kanałem"
        );
    }

    #[test]
    fn kazde_wydanie_z_wykazu_lezy_w_kanale() {
        // Zgodność samego kanal.adres nie wystarcza, pozycja może nieść inny host.
        for adres in adresy_wydan_z_wykazu() {
            sprawdz_zadanie(&adres, &suma_poprawna()).unwrap_or_else(|odmowa| {
                panic!(
                    "pozycja wykazu wydań wskazuje {adres}, którego powłoka nie \
                     przyjmie ({}): {}",
                    odmowa.powod, odmowa.zdanie
                )
            });
        }
    }

    #[test]
    fn adres_z_kanalu_przechodzi() {
        let adres = adresy_wydan_z_wykazu()
            .into_iter()
            .next()
            .expect("pierwsza pozycja wykazu");
        let suma = sprawdz_zadanie(&adres, &suma_poprawna())
            .expect("adres z wykazu wydań leży w kanale i ma przejść");
        assert_eq!(suma, suma_poprawna(), "suma ma wrócić w małych literach");
    }

    #[test]
    fn adres_spoza_kanalu_odrzucony() {
        // Strona podstawiona wskazuje własny serwer z pasującą sumą; zatrzymać ma kanał.
        let odmowa = sprawdz_zadanie("https://obcy-serwer.example/wydanie-setup.exe", &suma_poprawna())
            .expect_err("adres spoza kanału pobrań nie ma prawa przejść");
        assert_eq!(odmowa.powod, "adres-poza-kanalem");
        assert!(
            odmowa.zdanie.contains(ADRES_KANALU),
            "odmowa ma nazwać kanał, żeby było wiadomo, skąd wydania przychodzą: {}",
            odmowa.zdanie
        );
    }

    #[test]
    fn kanal_jako_poczatek_obcej_nazwy_odrzucony() {
        // Nazwa hosta zaczynająca się nazwą kanału to wciąż obcy host.
        let odmowa = sprawdz_zadanie(
            "https://pobierz.danaco-group.pl.obcy-serwer.example/wydanie-setup.exe",
            &suma_poprawna(),
        )
        .expect_err("host o nazwie zaczynającej się nazwą kanału jest poza kanałem");
        assert_eq!(odmowa.powod, "adres-poza-kanalem");
    }

    #[test]
    fn adres_http_odrzucony_wczesniejsza_odmowa() {
        // HTTP pada na warunku protokołu, nie kanału.
        let odmowa = sprawdz_zadanie("http://pobierz.danaco-group.pl/wydanie-setup.exe", &suma_poprawna())
            .expect_err("adres bez https nie ma prawa przejść");
        assert_eq!(odmowa.powod, "adres-nie-https");
    }

    #[test]
    fn suma_w_zlym_zapisie_odrzucona_takze_dla_adresu_z_kanalu() {
        // Kanał nie zdejmuje pozostałych warunków — suma obcięta pada tak samo.
        let adres = adresy_wydan_z_wykazu()
            .into_iter()
            .next()
            .expect("pierwsza pozycja wykazu");
        let odmowa = sprawdz_zadanie(&adres, "abc123")
            .expect_err("suma krótsza niż 64 znaki nie ma prawa przejść");
        assert_eq!(odmowa.powod, "suma-w-zlym-zapisie");
    }
}
