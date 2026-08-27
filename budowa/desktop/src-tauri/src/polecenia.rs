//! Polecenia powłoki dostępne dla interfejsu.
//!
//! Powłoka udostępnia wyłącznie to, czego przeglądarka zrobić nie może:
//! natywny wybór katalogu, wskazanie serwera wdrożenia i podmianę pliku
//! aplikacji. Reszta pracy idzie kanałem WebSocket kontraktu `shared/` —
//! powłoka nie tworzy drugiej drogi sterowania platformą.
//!
//! Czego tu nie ma i nie będzie: postawienia ani zatrzymania rdzenia. Rdzeń
//! stoi na serwerze wdrożenia, powłoka go nie niesie i nie ma nad nim władzy
//! (`rdzen/mod.rs`).

use tauri::{AppHandle, State};

use crate::aktualizacja::{self, Odmowa, Przebieg};
use crate::dialog_katalogu;
use crate::rdzen::{self, OpisRdzenia};
use crate::ustawienia::Ustawienia;
use crate::wskazanie::{self, Wskazanie};

/// Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę.
///
/// Jedno polecenie obsługuje oba zastosowania interfejsu — wskazanie katalogu
/// roboczego i dodanie katalogu jako punktu dostępu — bo czynność systemu
/// operacyjnego jest ta sama. Różni je wyłącznie napis w belce, który podaje
/// wywołujący; wywołanie bez napisu dostaje wartość domyślną i działa dalej.
///
/// Rezygnacja Operatora zwraca `None`, czyli `null` po stronie interfejsu.
#[tauri::command]
pub async fn wybierz_katalog_roboczy(
    aplikacja: AppHandle,
    tytul: Option<String>,
) -> Option<String> {
    dialog_katalogu::wybierz(&aplikacja, tytul.as_deref()).await
}

/// Zwraca stan rdzenia: czy odpowiada pod wskazanym adresem i pod jakim.
///
/// Wskaźnik łączności interfejsu bez tej odpowiedzi umie powiedzieć tylko
/// „Rozłączony". Powłoka zna adres wskazany przez Operatora i ścieżkę własnego
/// dziennika, więc to ona składa zdanie, które da się pokazać człowiekowi.
#[tauri::command]
pub fn stan_rdzenia(ustawienia: State<'_, Ustawienia>) -> OpisRdzenia {
    rdzen::opisz(&ustawienia)
}

/// Zwraca adres HTTP rdzenia obowiązujący albo `null`, gdy wskazania nie złożono.
///
/// Interfejs nie ma jak wyliczyć tego adresu z lokalizacji dokumentu: strona
/// pochodzi z pakietu wkompilowanego w powłokę, więc `location.hostname` mówi
/// `tauri.localhost`, a rdzeń stoi na serwerze wdrożenia. To polecenie jest
/// jedyną drogą, którą warstwa połączenia poznaje ten adres
/// (`klient/src/polaczenie/adres-rdzenia.ts`, `adresGniazdaRdzenia`).
#[tauri::command]
pub fn adres_rdzenia(ustawienia: State<'_, Ustawienia>) -> Option<String> {
    ustawienia.adres_rdzenia_http()
}

/// Zwraca stan wskazania rdzenia: na którym serwerze stoi i z której warstwy
/// wskazanie pochodzi.
///
/// Okno pyta o to przed złożeniem aplikacji. Odpowiedź z warstwą `brak` znaczy
/// pierwsze uruchomienie po instalacji — wtedy staje ekran wskazania, bo
/// instalator adresu serwera wdrożenia nie zna i znać go nie może.
#[tauri::command]
pub fn wskazanie_rdzenia(ustawienia: State<'_, Ustawienia>) -> Wskazanie {
    wskazanie::biezace(&ustawienia)
}

/// Przyjmuje wskazanie Operatora, na którym serwerze stoi rdzeń: sprawdza
/// łączność i zapisuje wskazanie trwale.
///
/// `adres` jest treścią pola z okna (nazwa serwera, `serwer:port`, także
/// z przedrostkiem `http://`). Odmowa niesie kod powodu do rozgałęzienia
/// i gotowe zdanie do pokazania.
///
/// To jedyne polecenie powłoki zmieniające jej stan trwały, i ma być jedyne:
/// wskazanie serwera jest czynnością, której przeglądarka wykonać nie może, bo
/// dotyczy pliku nastaw na dysku.
#[tauri::command]
pub fn wskaz_rdzen(
    ustawienia: State<'_, Ustawienia>,
    adres: Option<String>,
) -> Result<Wskazanie, wskazanie::Odmowa> {
    wskazanie::wskaz(&ustawienia, adres.as_deref().unwrap_or_default())
}

/// Pobiera wskazane wydanie, sprawdza jego sumę, zakłada je i stawia aplikację
/// na nowo. Strona powłoki banera „aktualizuje".
///
/// `adres` i `suma_sha256` pochodzą z wykazu `budowa/witryna/wydania.json`,
/// który interfejs czyta po HTTPS z kanału wskazanego w tym wykazie
/// (`kanal.adres`). Powłoka wykazu nie czyta i wersji
/// nie porównuje — do tego wystarczy przeglądarka. Powłoka robi dwie rzeczy
/// niemożliwe ze strony: pisze po dysku pod plikiem aplikacji i stawia proces
/// na nowo.
///
/// Zwraca `Przebieg` (co założono, ile bajtów, suma, za ile milisekund zniknie
/// okno) albo `Odmowa` (`powod` — kod do rozgałęzienia, `zdanie` — gotowy tekst
/// do banera). Odpowiedź pomyślna dociera do interfejsu przed restartem dzięki
/// zwłoce opisanej w `aktualizacja/mod.rs`.
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
