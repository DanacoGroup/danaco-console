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
/// Łańcuch zaufania ma dwa ogniwa: HTTPS chroni wykaz wydań, z którego pochodzi
/// suma, a suma chroni pobrany plik. Adres inny niż `https://` jest odmawiany,
/// bo po zwykłym HTTP pośrednik podmienia plik i sumę naraz, co unieważnia
/// całe sprawdzenie.
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
/// Wydzielone z `pobierz_i_sprawdz`, bo są to jedyne dwie odmowy padające przed
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
