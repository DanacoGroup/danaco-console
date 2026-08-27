//! Aktualizacja zdalna powłoki — pobranie wydania, sprawdzenie sumy, założenie,
//! ponowny start. To strona powłoki dla banera „aktualizuj".
//!
//! Przedmiotem aktualizacji jest wyłącznie plik powłoki na urządzeniu
//! Operatora. Rdzeń stoi na serwerze wdrożenia i jest utrzymywany tam, więc ten
//! przebieg go nie dotyka: niczego mu nie podmienia i niczego nie wygasza.
//!
//! Droga nie idzie przez wtyczkę `updater` Tauri: wtyczka nie występuje ani
//! w `Cargo.toml`, ani w `tauri.conf.json`, ani w `capabilities/domyslne.json`.
//! Wymaga własnego podpisu minisign, czyli pary kluczy wydawcy, i narzuca
//! własny kształt pliku `latest.json`, podczas gdy wykazem wydań jest
//! `budowa/witryna/wydania.json` — ten sam plik, z którego strona „Pobierz"
//! bierze chronologię.
//!
//! Zamiast tego: HTTPS po plik i SHA-256 z wykazu. Suma wiąże plik z wykazem
//! tak samo jak podpis, a wykaz przychodzi po HTTPS z domeny wydawcy. Przejście
//! na podpis dotknie jednego pliku (`pobranie.rs`); reszta przebiegu zostaje
//! bez zmiany.

pub mod droga;
pub mod pobranie;
#[cfg(test)]
pub mod probne;

use std::sync::atomic::{AtomicBool, Ordering};
use std::time::Duration;

use serde::Serialize;
use tauri::AppHandle;

use crate::dziennik;

/// Ile czasu dostaje interfejs na odebranie odpowiedzi, zanim powłoka zniknie.
///
/// Restart natychmiastowy zabiłby okno, zanim odpowiedź polecenia doszłaby do
/// strony — baner nie zdążyłby powiedzieć Operatorowi, że się udało, i wyglądałoby
/// to jak awaria. Ta zwłoka nie jest blokadą: nic nie pyta i niczego
/// nie wstrzymuje, daje tylko odpowiedzi czas dolecieć.
const ZWLOKA_RESTARTU: Duration = Duration::from_millis(1500);

/// Odmowa wykonania aktualizacji — opisuje brak, nigdy zakaz.
///
/// Powłoka nie mówi Operatorowi „nie wolno". Mówi, czego zabrakło (sieci,
/// zgodnej sumy, prawa zapisu, drogi na tym systemie) i — gdy to możliwe —
/// co da się z tym zrobić samemu. `powod` jest kodem dla interfejsu,
/// `zdanie` jest zdaniem dla człowieka.
#[derive(Clone, Debug, Serialize)]
pub struct Odmowa {
    /// Kod przyczyny, stały i nadający się do rozgałęzienia w interfejsie.
    pub powod: String,
    /// Zdanie po polsku, gotowe do pokazania w banerze.
    pub zdanie: String,
}

impl Odmowa {
    /// Składa odmowę z kodu i zdania.
    pub fn nowa(powod: &str, zdanie: String) -> Self {
        Self {
            powod: powod.to_string(),
            zdanie,
        }
    }
}

/// Przebieg udanej aktualizacji — to, co interfejs dostaje tuż przed restartem.
#[derive(Clone, Debug, Serialize)]
pub struct Przebieg {
    /// Droga systemowa, którą poszło założenie: `appimage` albo `instalka-nsis`.
    pub droga: String,
    /// Co zostało podmienione lub uruchomione (ścieżka).
    pub zalozone: String,
    /// Ile bajtów pobrano.
    pub bajtow: u64,
    /// Suma SHA-256 policzona z pobranego pliku — zgodna z żądaną.
    pub suma_sha256: String,
    /// Za ile milisekund aplikacja zniknie z ekranu.
    pub restart_za_ms: u64,
    /// Zdanie dla Operatora.
    pub zdanie: String,
}

/// Czy aktualizacja właśnie trwa.
///
/// To zapora wejścia, nie bramka na Operatora. Droga pobierania jest
/// deterministyczna: bez wykluczenia dwa wywołania (dwuklik w baner, dwa okna)
/// otwierają `File::create` na tym samym pliku roboczym obok aplikacji, piszą
/// w niego przeplotem i każde liczy SHA-256 z własnego strumienia, a nie z tego,
/// co ostatecznie leży na dysku. Suma zgadzałaby się wtedy dla pliku, którego
/// w tej postaci nie ma, po czym oba przebiegi wołałyby `rename` na plik
/// aplikacji. Odmowa stąd opisuje brak (aktualizacja już zajęta), nie zakaz,
/// i mija sama, gdy pierwszy przebieg się skończy.
static W_TOKU: AtomicBool = AtomicBool::new(false);

/// Zwalnia zaporę także wtedy, gdy przebieg wyszedł błędem albo paniką —
/// inaczej jedno niepowodzenie zamykałoby aktualizacje do końca życia procesu.
pub struct StrazWylacznosci;

impl Drop for StrazWylacznosci {
    fn drop(&mut self) {
        W_TOKU.store(false, Ordering::Release);
    }
}

/// Zajmuje wyłączność na aktualizację albo odmawia, bo już trwa.
///
/// Wydzielone z `wykonaj`, żeby dało się to sprawdzić bez stawiania okna —
/// zapora jest tu jedyną rzeczą stojącą między dwoma kliknięciami a dwoma
/// pobraniami do tego samego pliku.
///
/// Zwróconą straż trzeba związać z nazwą na cały czas przebiegu. Wartość
/// upuszczona od razu zwalnia zaporę i przywraca usterkę, przed którą ta
/// funkcja stoi — stąd `#[must_use]`, a `wykonaj` przekazuje straż pożyczką
/// do `przebieg_pod_straza`.
#[must_use = "straż zwalnia zaporę w chwili upuszczenia — zwiąż ją z nazwą na czas całego przebiegu"]
pub fn zajmij_wylacznosc() -> Result<StrazWylacznosci, Odmowa> {
    if W_TOKU
        .compare_exchange(false, true, Ordering::AcqRel, Ordering::Acquire)
        .is_err()
    {
        return Err(Odmowa::nowa(
            "aktualizacja-w-toku",
            "Aktualizacja już trwa — powłoka pobiera wydanie i za chwilę je założy. \
             Drugiego pobrania nie zaczyna, bo oba pisałyby w ten sam plik roboczy \
             i podmieniały aplikację jeden przez drugiego."
                .to_string(),
        ));
    }
    Ok(StrazWylacznosci)
}

/// Wykonuje aktualizację od początku do końca. Wołane przez polecenie IPC.
///
/// Kolejność jest warunkiem poprawności i nie wolno jej mieszać:
/// rozpoznanie drogi → pobranie → **sprawdzenie sumy** → założenie →
/// ponowny start.
pub fn wykonaj(aplikacja: &AppHandle, adres: &str, suma_sha256: &str) -> Result<Przebieg, Odmowa> {
    dziennik::dopisz(&format!("aktualizacja: żądanie wydania spod {adres}"));

    // Jedna aktualizacja naraz — zob. `W_TOKU`.
    let straz = zajmij_wylacznosc().inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;

    // Straż idzie dalej pożyczką, a nie zostaje tu w luźnym wiązaniu.
    //
    // Przy samym wiązaniu `let _straz = …` skrócenie go do `let _ = …`
    // upuszczałoby straż natychmiast: wyłączność znikałaby, dwuklik znów
    // uruchamiałby dwa pobrania do jednego pliku roboczego, i nie zapaliłby
    // się przy tym ani test, ani ostrzeżenie clippy. Przy pożyczce ten sam błąd
    // jest niemożliwy do zapisania: bez nazwy nie ma czego pożyczyć, a pożyczka
    // trzyma straż przy życiu do końca wywołania.
    przebieg_pod_straza(&straz, aplikacja, adres, suma_sha256)
}

/// Właściwy przebieg aktualizacji, wykonywany pod zajętą wyłącznością.
///
/// `_straz` nie jest tu do niczego używana i o to chodzi: jej obecność
/// w podpisie jest dowodem — sprawdzanym przez kompilator — że nikt nie
/// wywoła tego przebiegu bez zajętej zapory ani nie zwolni jej w połowie.
fn przebieg_pod_straza(
    _straz: &StrazWylacznosci,
    aplikacja: &AppHandle,
    adres: &str,
    suma_sha256: &str,
) -> Result<Przebieg, Odmowa> {
    // Rozpoznanie drogi: czy na tym systemie w ogóle jest co podmieniać.
    // Pytamy przed pobraniem — ściąganie całego wydania po to, żeby odmówić,
    // byłoby marnotrawstwem łącza Operatora.
    let droga = droga::rozpoznaj().inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;
    let plik_roboczy = droga.plik_roboczy();

    // Pobranie wraz ze sprawdzeniem sumy. Plik niezgodny nie wraca stąd nigdy
    // — `pobranie.rs` kasuje go przed zwróceniem odmowy.
    let pobrany =
        pobranie::pobierz_i_sprawdz(adres, suma_sha256, &plik_roboczy).inspect_err(|odmowa| {
            dziennik::dopisz(&format!(
                "aktualizacja: odmowa ({}) {}",
                odmowa.powod, odmowa.zdanie
            ));
        })?;
    dziennik::dopisz(&format!(
        "aktualizacja: pobrano {} B, suma {} zgodna z wykazem",
        pobrany.bajtow, pobrany.suma_sha256
    ));

    // Założenie. Dopiero tutaj cokolwiek na dysku Operatora się zmienia.
    let zalozenie = droga.zaloz(&pobrany.sciezka).inspect_err(|odmowa| {
        dziennik::dopisz(&format!(
            "aktualizacja: odmowa ({}) {}",
            odmowa.powod, odmowa.zdanie
        ));
    })?;
    dziennik::dopisz(&format!("aktualizacja: {}", zalozenie.zdanie));

    // Ponowny start — po zwłoce, żeby odpowiedź zdążyła dojść do banera.
    zaplanuj_ponowny_start(aplikacja, zalozenie.restartuje_powloka);

    Ok(Przebieg {
        droga: droga.nazwa().to_string(),
        zalozone: zalozenie.zalozone,
        bajtow: pobrany.bajtow,
        suma_sha256: pobrany.suma_sha256,
        restart_za_ms: ZWLOKA_RESTARTU.as_millis() as u64,
        zdanie: zalozenie.zdanie,
    })
}

/// Odkłada ponowny start powłoki o `ZWLOKA_RESTARTU`.
///
/// `restart()` Tauri sam radzi sobie z AppImage: sięga po zmienną `APPIMAGE`,
/// a nie po `current_exe()`, które wskazywałoby chwilowo podmontowany obraz
/// starego wydania. Dzięki temu po podmianie wstaje wydanie nowe.
///
/// Gdy zakłada instalator zewnętrzny (Windows), powłoka wyłącznie schodzi
/// z drogi — `zakoncz_powloke` kończy proces okna i nic poza nim.
fn zaplanuj_ponowny_start(aplikacja: &AppHandle, restartuje_powloka: bool) {
    let aplikacja = aplikacja.clone();
    std::thread::spawn(move || {
        std::thread::sleep(ZWLOKA_RESTARTU);
        // Wątek roboczy nie może sam kończyć aplikacji — Tauri wymaga do tego
        // wątku głównego, inaczej sprzątanie okna bywa niedokończone.
        let _ = aplikacja.clone().run_on_main_thread(move || {
            if restartuje_powloka {
                dziennik::dopisz("aktualizacja: ponowny start powłoki");
                aplikacja.restart();
            } else {
                dziennik::dopisz("aktualizacja: powłoka schodzi instalatorowi z drogi");
                crate::zamkniecie::zakoncz_powloke(&aplikacja);
            }
        });
    });
}
