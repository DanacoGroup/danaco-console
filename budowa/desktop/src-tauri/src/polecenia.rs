//! Polecenia powłoki dostępne dla interfejsu.
//!
//! Powłoka udostępnia wyłącznie to, czego przeglądarka zrobić nie może:
//! natywny wybór katalogu, stan rdzenia w tle i adres rdzenia. Reszta pracy
//! idzie kanałem WebSocket kontraktu `shared/` — powłoka
//! nie tworzy drugiej drogi sterowania platformą.
//!
//! Każde polecenie ma konsumenta — polecenie bez konsumenta wygląda jak droga,
//! a jest ścianą:
//!
//!   * `wybierz_katalog_roboczy` → `client/src/powloka/most-katalogow.ts`
//!   * `adres_rdzenia`           → `client/src/powloka/most-rdzenia.ts`
//!   * `stan_rdzenia`            → `client/src/powloka/most-rdzenia.ts`
//!   * `wskazanie_rdzenia`       → `client/src/powloka/most-rdzenia.ts`
//!   * `wskaz_rdzen`             → `client/src/wskazanie-rdzenia/ekran-wskazania.ts`
//!   * `wykonaj_aktualizacje`    → baner „aktualizuje" w oknie aplikacji
//!
//! Zatrzymania rdzenia poleceniem IPC tu nie ma. `rdzen/uchwyt.rs` stanowi, że
//! zatrzymanie jest wyłącznie jawnym poleceniem Operatora z zasobnika, bo
//! zamknięcie okna nie kończy procesów sesji. Polecenie IPC dawałoby tej samej
//! czynności drugą drogę — z wnętrza strony, która na tym rdzeniu stoi —
//! a interfejs nie ma czym rdzenia z powrotem postawić. Zasobnik żyje po
//! zamknięciu okna i wywołuje `UchwytRdzenia::zatrzymaj` wprost
//! (`menu_zasobnika.rs`).

use tauri::{AppHandle, State};

use crate::aktualizacja::{self, Odmowa, Przebieg};
use crate::dialog_katalogu;
use crate::rdzen::{OpisRdzenia, UchwytRdzenia};
use crate::ustawienia::Ustawienia;
use crate::wskazanie::{self, Wskazanie};

/// Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę.
///
/// Jedno polecenie obsługuje oba zastosowania interfejsu — wskazanie katalogu
/// roboczego i dodanie katalogu jako punktu dostępu — bo czynność systemu
/// operacyjnego jest ta sama. Różni je wyłącznie napis w belce, który podaje
/// wywołujący; wywołanie bez napisu (także ze starszego interfejsu) dostaje
/// wartość domyślną i działa dalej.
///
/// Rezygnacja Operatora zwraca `None`, czyli `null` po stronie interfejsu.
#[tauri::command]
pub async fn wybierz_katalog_roboczy(
    aplikacja: AppHandle,
    tytul: Option<String>,
) -> Option<String> {
    dialog_katalogu::wybierz(&aplikacja, tytul.as_deref()).await
}

/// Zwraca stan rdzenia uruchomionego w tle — także wtedy, gdy się nie udało.
///
/// Powód, dla którego łączności nie ma, zna wyłącznie powłoka: to ona stawiała
/// proces, zna port, identyfikator procesu i ścieżkę dziennika. Wskaźnik
/// łączności interfejsu bez tej odpowiedzi umie powiedzieć tylko „Rozłączony".
/// Konsument: `client/src/powloka/most-rdzenia.ts` (przez
/// `client/src/aplikacja/wskaznik-lacznosci.ts`).
#[tauri::command]
pub fn stan_rdzenia(uchwyt: State<'_, UchwytRdzenia>) -> OpisRdzenia {
    uchwyt.opis()
}

/// Zwraca adres HTTP rdzenia obowiązujący, a nie wyliczony ze stałej.
///
/// Interfejs wylicza adres gniazda z lokalizacji dokumentu
/// (`client/src/polaczenie/adres-rdzenia.ts`), a przy pakiecie osadzonym
/// w powłoce lokalizacja mówi `tauri.localhost` — adres, pod którym nikt nie
/// nasłuchuje (`zrodlo_interfejsu.rs`). To polecenie jest odpowiedzią na tę
/// rozbieżność.
///
/// Adres składa `Ustawienia::adres_rdzenia_http` z hosta i portu obowiązujących
/// w tej chwili (`Ustawienia::host`, `Ustawienia::port`) — bez wskazania jest to
/// `127.0.0.1`, ze wskazaniem `DANACO_HOST_RDZENIA` adres rdzenia stojącego
/// gdzie indziej (wariant wirtualny). Polecenie nie rozstrzyga, skąd ładuje się
/// strona interfejsu — to osobna nastawa (`DANACO_ADRES_INTERFEJSU`,
/// `zrodlo_interfejsu.rs`).
/// Konsument: `client/src/powloka/most-rdzenia.ts` (przez `client/src/main.ts`).
#[tauri::command]
pub fn adres_rdzenia(ustawienia: State<'_, Ustawienia>) -> String {
    ustawienia.adres_rdzenia_http()
}

/// Zwraca stan wskazania rdzenia: gdzie stoi, czy wskazanie w ogóle złożono
/// i z której warstwy pochodzi.
///
/// Okno pyta o to przed złożeniem aplikacji. Odpowiedź `zlozone: false` znaczy
/// pierwsze uruchomienie po instalacji — wtedy staje ekran wskazania
/// (`client/src/wskazanie-rdzenia/`), bo instalator adresu serwera nie zna i znać
/// go nie może.
/// Konsument: `client/src/powloka/most-rdzenia.ts`.
#[tauri::command]
pub fn wskazanie_rdzenia(ustawienia: State<'_, Ustawienia>) -> Wskazanie {
    wskazanie::biezace(&ustawienia)
}

/// Przyjmuje wskazanie Operatora, gdzie stoi rdzeń: sprawdza łączność, zapisuje
/// wskazanie trwale i wprowadza je w życie w tym uruchomieniu.
///
/// `adres` jest treścią pola z okna (nazwa serwera, `host:port`, także
/// z przedrostkiem `http://`), a `lokalnie` mówi, że Operator wybrał rdzeń na tym
/// urządzeniu — wtedy adres nie jest czytany. Odmowa niesie kod powodu do
/// rozgałęzienia i gotowe zdanie do pokazania.
///
/// To jedyne polecenie powłoki zmieniające stan platformy, i ma być jedyne:
/// wskazanie miejsca rdzenia jest czynnością, której przeglądarka wykonać nie
/// może, bo dotyczy procesu stawianego przez powłokę i pliku nastaw na dysku.
/// Konsument: `client/src/wskazanie-rdzenia/ekran-wskazania.ts`.
#[tauri::command]
pub fn wskaz_rdzen(
    ustawienia: State<'_, Ustawienia>,
    uchwyt: State<'_, UchwytRdzenia>,
    adres: Option<String>,
    lokalnie: Option<bool>,
) -> Result<Wskazanie, wskazanie::Odmowa> {
    wskazanie::wskaz(
        &ustawienia,
        &uchwyt,
        adres.as_deref().unwrap_or_default(),
        lokalnie.unwrap_or(false),
    )
}

/// Pobiera wskazane wydanie, sprawdza jego sumę, zakłada je i stawia aplikację
/// na nowo. Strona powłoki banera „aktualizuje".
///
/// `adres` i `suma_sha256` pochodzą z wykazu `wydania.json`, który interfejs
/// czyta po HTTPS spod `danaco-console.pl`. Powłoka wykazu nie czyta i wersji
/// nie porównuje — do tego wystarczy przeglądarka. Powłoka robi trzy rzeczy
/// niemożliwe ze strony: pisze po dysku pod plikiem aplikacji, zatrzymuje rdzeń
/// i stawia proces na nowo.
///
/// Zwraca `Przebieg` (co założono, ile bajtów, suma, los rdzenia, za ile
/// milisekund zniknie okno) albo `Odmowa` (`powod` — kod do rozgałęzienia,
/// `zdanie` — gotowy tekst do banera). Odpowiedź pomyślna dociera do interfejsu
/// przed restartem dzięki zwłoce opisanej w `aktualizacja/mod.rs`.
///
/// Praca idzie wątkiem roboczym (`spawn_blocking`), bo pobranie wydania to
/// dziesiątki megabajtów: na wątku głównym okno stałoby zamrożone przez cały
/// czas ściągania, a baner nie zdążyłby nawet się przerysować.
#[tauri::command]
pub async fn wykonaj_aktualizacje(
    aplikacja: AppHandle,
    adres: String,
    suma_sha256: String,
) -> Result<Przebieg, Odmowa> {
    tauri::async_runtime::spawn_blocking(move || {
        aktualizacja::wykonaj(&aplikacja, &adres, &suma_sha256)
    })
    .await
    .unwrap_or_else(|blad| {
        // Wątek roboczy przepadł (panika). Aktualizacja się nie odbyła —
        // Operator ma to usłyszeć, a nie zobaczyć baner, który zamilkł.
        Err(Odmowa::nowa(
            "przebieg-przerwany",
            format!(
                "Aktualizacja przerwała się nieoczekiwanie ({blad}). Nic nie zostało \
                 założone — spróbuj ponownie albo pobierz plik ręcznie ze strony „Pobierz”."
            ),
        ))
    })
}
