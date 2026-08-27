//! Pobranie pliku wydania i sprawdzenie jego sumy kontrolnej.
//!
//! Jedyne miejsce w powłoce, w którym dane z sieci stają się plikiem
//! uruchamianym na maszynie użytkownika. Suma kontrolna jest tu warunkiem,
//! nie diagnostyką: plik o niezgodnej sumie zostaje skasowany, zanim funkcja
//! wróci. Kolejność kroków jest wiążąca:
//!   1. pobierz do pliku roboczego obok celu, nie pod nazwą celu,
//!   2. policz SHA-256 z tego, co leży na dysku,
//!   3. porównaj z sumą żądaną,
//!   4. dopiero potem wolno zakładać plik — robi to `droga.rs`.
//!
//! Suma liczona jest w locie z tego samego strumienia, który trafia na dysk.
//! Liczona z osobnego bufora w pamięci sprawdzałaby co innego niż zapisany plik.

use std::fs::{self, File};
use std::io::{BufWriter, Read, Write};
use std::path::{Path, PathBuf};

use sha2::{Digest, Sha256};

use super::Odmowa;

/// Adres kanału pobrań — `kanal.adres` z wykazu wydań
/// `budowa/witryna/wydania.json`, jedynego źródła prawdy strony „Pobierz”
/// i banera aktualizacji. Wszystkie pola `plik` wykazu leżą pod tym adresem.
///
/// Powłoka pobiera wydania wyłącznie spod tego adresu. Bez tego przywiązania
/// suma SHA-256 niczego by nie chroniła: adres i sumę podaje powłoce ta sama
/// strona, więc strona podstawiona wskazałaby własny plik wraz z jego poprawną
/// sumą. Kanał jest częścią powłoki z tego samego powodu, z którego suma jest
/// częścią wykazu — ogniwo zaufania nie może pochodzić od strony, którą wiąże.
///
/// Porównanie jest dosłowne (przedrostek znak w znak) i kończy się ukośnikiem,
/// więc `pobierz.danaco-group.pl.obcy-serwer` przedrostka nie przejdzie.
///
/// Wartość jest kopią `kanal.adres` z wykazu wydań i musi nią pozostać —
/// pilnuje tego sprawdzian `kanal_powloki_zgadza_sie_z_wykazem_wydan` w tym
/// module, który czyta wykaz przy kompilacji. Bez niego przeniesienie kanału
/// rozeszłoby obie wartości bez śladu: wykaz wskazywałby nowy adres, powłoka
/// odmawiałaby każdego wydania jako leżącego poza kanałem.
const ADRES_KANALU: &str = "https://pobierz.danaco-group.pl/";

/// Górny pułap wielkości pliku wydania — 512 MiB.
///
/// Zabezpieczenie przed zapełnieniem dysku przez odpowiedź bez końca z serwera
/// zepsutego albo podstawionego. Pułap dobrany z zapasem względem wielkości
/// pakietu powłoki, która idzie w dziesiątki megabajtów.
const PULAP_BAJTOW: u64 = 512 * 1024 * 1024;

/// Wynik udanego pobrania: plik roboczy leży na dysku i ma zgodną sumę.
pub struct Pobrany {
    /// Ścieżka pliku roboczego, jeszcze nie założonego pod nazwą celu.
    pub sciezka: PathBuf,
    /// Liczba zapisanych bajtów.
    pub bajtow: u64,
    /// Suma policzona z zapisanego pliku; równa żądanej, inaczej byłaby odmowa.
    pub suma_sha256: String,
}

/// Pobiera plik spod `adres` do `plik_roboczy` i sprawdza jego sumę SHA-256.
///
/// `suma_zadana` to suma w zapisie szesnastkowym (wielkość liter bez znaczenia),
/// pochodząca z `wydania.json` czytanego przez interfejs po HTTPS.
///
/// Łańcuch zaufania ma trzy ogniwa: kanał pobrań wkompilowany w powłokę
/// (`ADRES_KANALU`) ogranicza, skąd plik w ogóle może przyjść; HTTPS chroni
/// wykaz wydań, z którego pochodzi suma; a suma chroni pobrany plik. Adres
/// inny niż `https://` jest odmawiany, bo po zwykłym HTTP pośrednik podmienia
/// plik i sumę naraz — a adres spoza kanału jest odmawiany, bo adres i sumę
/// podaje powłoce ta sama strona, więc bez kanału nic nie wiąże ich z wydawcą.
pub fn pobierz_i_sprawdz(
    adres: &str,
    suma_zadana: &str,
    plik_roboczy: &Path,
) -> Result<Pobrany, Odmowa> {
    let adres = adres.trim();
    let suma_zadana = sprawdz_zadanie(adres, suma_zadana)?;

    // Bez `http_status_as_error(false)` `ureq` zwija odpowiedź 404 do tego samego
    // błędu co zerwane połączenie: użytkownik dostaje komunikat o łączności,
    // choć sieć działa, a po prostu nie ma jeszcze takiego pliku wydania.
    // Wyłączenie tego zachowania utrzymuje przy życiu gałąź `odpowiedz-serwera`.
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

    // Katalog docelowy musi istnieć i być zapisywalny. Próba utworzenia pliku
    // roboczego sprawdza jedno i drugie, zanim z gniazda przyjdzie pierwszy bajt.
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

    // Zapis i liczenie sumy idą w jednym przebiegu, kawałkami po 64 KiB —
    // plik wydania ma dziesiątki megabajtów i nie ma powodu trzymać go w pamięci.
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

    // Bufor musi trafić na dysk przed oceną wyniku, inaczej ocenie podlegałyby
    // bajty, które jeszcze nie są plikiem.
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
/// plik roboczy.
///
/// Zwraca sumę sprowadzoną do małych liter — w tej postaci porównuje ją
/// `ocen_pobrane`, bo `format!("{:x}")` daje małe litery, a wykaz wydań pisany
/// ręcznie potrafi mieć wielkie.
///
/// Wydzielone z `pobierz_i_sprawdz`, bo są to jedyne odmowy padające przed
/// pierwszym bajtem z gniazda; test może je wywołać bez dostępu do sieci.
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

    // Adres musi leżeć w kanale pobrań. Bez tego warunku powłoka pobierałaby
    // spod dowolnego adresu, który wskaże strona — a stronie wystarczy podać
    // do niego pasującą sumę, bo obie wartości przychodzą tym samym poleceniem.
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

    // Suma SHA-256 ma dokładnie 64 znaki szesnastkowe. Sprawdzenie idzie przed
    // pobraniem, bo suma pusta albo obcięta unieważnia późniejsze porównanie,
    // a użytkownik czekałby na pobranie tylko po to, żeby usłyszeć odmowę.
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

/// Ocenia plik, który wylądował na dysku. Wołane po `sync_all`.
///
/// Kolejność pytań jest warunkiem poprawności: najpierw dolna granica
/// wielkości, dopiero potem suma. Plik pusty ma poprawną, 64-znakową sumę
/// SHA-256, więc wykaz podający właśnie ją — albo serwer oddający 200 z pustym
/// ciałem — przeszedłby przez sam warunek sumy, a `droga.rs` podstawiłoby zero
/// bajtów w miejsce aplikacji i zaplanowało restart, po którym nie ma z czego
/// wrócić.
///
/// Kasowanie pliku roboczego zostaje po stronie wołającego (`posprzataj`),
/// bo to on wie, jaki plik utworzył.
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

/// Kasuje plik roboczy i przepuszcza odmowę dalej.
///
/// Każde wyjście błędem po utworzeniu pliku przechodzi tędy. Plik pobrany, ale
/// niesprawdzony albo wprost niezgodny, nie może zostać na dysku: leżałby obok
/// właściwej aplikacji, wyglądając na jej część. Niepowodzenie samego kasowania
/// nie zmienia odmowy — zwracana jest pierwotna przyczyna.
fn posprzataj(plik_roboczy: &Path, odmowa: Odmowa) -> Odmowa {
    let _ = fs::remove_file(plik_roboczy);
    odmowa
}

#[cfg(test)]
mod testy {
    //! Testy sprawdzenia żądania — wszystkie biegną bez sieci, bo `sprawdz_zadanie`
    //! pada przed pierwszym bajtem z gniazda. Osią jest przywiązanie do kanału
    //! pobrań: adres spoza `ADRES_KANALU` ma zostać odrzucony, zanim cokolwiek
    //! poleci przez sieć i zanim powstanie plik roboczy.
    //!
    //! Wykaz wydań nie jest tu przepisany, tylko wczytany przy kompilacji
    //! (`include_str!`). Przepisany byłby drugą kopią tych samych wartości
    //! i rozjechałby się z wykazem równie cicho jak sam `ADRES_KANALU`.

    use super::*;

    /// Wykaz wydań `budowa/witryna/wydania.json` wczytany przy kompilacji.
    ///
    /// `include_str!` wiąże ten plik z budową testów: zmiana wykazu wymusza
    /// ponowną kompilację modułu, więc sprawdziany niżej nigdy nie oceniają
    /// treści sprzed zmiany.
    const WYKAZ_WYDAN: &str = include_str!("../../../../witryna/wydania.json");

    /// Wykaz wydań w postaci drzewa. Niepoprawny JSON jest tu porażką testu —
    /// wykaz jest źródłem prawdy strony „Pobierz” i banera aktualizacji.
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

    /// Adresy wszystkich pozycji wykazu — pola `wydania[].plik`.
    ///
    /// Pusty wykaz jest tu porażką, nie wynikiem: sprawdzian, który nie ma
    /// czego sprawdzić, przeszedłby milcząco i niczego by nie pilnował.
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

    /// Poprawna co do zapisu suma SHA-256 — 64 znaki szesnastkowe. Testom
    /// adresu wystarczy zapis; z niczym go tu nie porównują.
    fn suma_poprawna() -> String {
        "a".repeat(64)
    }

    #[test]
    fn kanal_powloki_zgadza_sie_z_wykazem_wydan() {
        // Sprawdzian wiążący dwie kopie jednego adresu. Zmiana którejkolwiek
        // z osobna — `ADRES_KANALU` tutaj albo `kanal.adres` w wykazie —
        // zatrzymuje się na tym porównaniu, zanim wyjdzie wydanie, którego
        // powłoka odmówiłaby pobrać.
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
        // Zgodność samego `kanal.adres` nie wystarcza: pozycja może nieść
        // adres spod innego hosta i przejdzie przez stronę „Pobierz”, a padnie
        // dopiero na maszynie Operatora.
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
        // Strona podstawiona wskazuje własny serwer wraz z pasującą sumą.
        // Sama suma tego nie zatrzyma — zatrzymać ma kanał.
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
        // Nazwa hosta zaczynająca się nazwą kanału to wciąż obcy host —
        // ukośnik na końcu `ADRES_KANALU` musi to odciąć.
        let odmowa = sprawdz_zadanie(
            "https://pobierz.danaco-group.pl.obcy-serwer.example/wydanie-setup.exe",
            &suma_poprawna(),
        )
        .expect_err("host o nazwie zaczynającej się nazwą kanału jest poza kanałem");
        assert_eq!(odmowa.powod, "adres-poza-kanalem");
    }

    #[test]
    fn adres_http_odrzucony_wczesniejsza_odmowa() {
        // HTTP pada na warunku protokołu, nie kanału — Operator ma usłyszeć
        // o HTTP, a nie ogólnik o kanale.
        let odmowa = sprawdz_zadanie("http://pobierz.danaco-group.pl/wydanie-setup.exe", &suma_poprawna())
            .expect_err("adres bez https nie ma prawa przejść");
        assert_eq!(odmowa.powod, "adres-nie-https");
    }

    #[test]
    fn suma_w_zlym_zapisie_odrzucona_takze_dla_adresu_z_kanalu() {
        // Kanał nie zdejmuje pozostałych warunków — suma obcięta pada tak samo
        // jak przed przywiązaniem do kanału.
        let adres = adresy_wydan_z_wykazu()
            .into_iter()
            .next()
            .expect("pierwsza pozycja wykazu");
        let odmowa = sprawdz_zadanie(&adres, "abc123")
            .expect_err("suma krótsza niż 64 znaki nie ma prawa przejść");
        assert_eq!(odmowa.powod, "suma-w-zlym-zapisie");
    }
}
