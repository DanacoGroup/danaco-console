//! Moduł udostępnia interfejsowi wyłącznie polecenia, których przeglądarka wykonać
//! nie może: wybór katalogu, wskazanie serwera i podmianę pliku aplikacji.

use tauri::{AppHandle, State};

use crate::aktualizacja::{self, Odmowa, Przebieg};
use crate::dialog_katalogu;
use crate::rdzen::{self, OpisRdzenia};
use crate::ustawienia::Ustawienia;
use crate::wskazanie::{self, Wskazanie};

/// Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę, oddając brak przy
/// rezygnacji Operatora z wyboru.
#[tauri::command]
pub async fn wybierz_katalog_roboczy(
    aplikacja: AppHandle,
    tytul: Option<String>,
) -> Option<String> {
    dialog_katalogu::wybierz(&aplikacja, tytul.as_deref()).await
}

/// Zwraca zdanie opisujące bieżący stan rdzenia: czy odpowiada pod wskazanym adresem
/// i pod jakim adresem stoi obecnie.
#[tauri::command]
pub fn stan_rdzenia(ustawienia: State<'_, Ustawienia>) -> OpisRdzenia {
    rdzen::opisz(&ustawienia)
}

/// Zwraca adres HTTP rdzenia obowiązujący aktualnie dla powłoki albo brak, gdy wskazania
/// serwera jeszcze nie złożono.
#[tauri::command]
pub fn adres_rdzenia(ustawienia: State<'_, Ustawienia>) -> Option<String> {
    ustawienia.adres_rdzenia_http()
}

/// Zwraca stan wskazania rdzenia: na którym serwerze wdrożenia stoi rdzeń i z której
/// warstwy to wskazanie pochodzi.
#[tauri::command]
pub fn wskazanie_rdzenia(ustawienia: State<'_, Ustawienia>) -> Wskazanie {
    wskazanie::biezace(&ustawienia)
}

/// Przyjmuje wskazanie Operatora, na którym serwerze wdrożenia stoi rdzeń, sprawdza
/// łączność i zapisuje wskazanie trwale.
#[tauri::command]
pub fn wskaz_rdzen(
    ustawienia: State<'_, Ustawienia>,
    adres: Option<String>,
) -> Result<Wskazanie, wskazanie::Odmowa> {
    wskazanie::wskaz(&ustawienia, adres.as_deref().unwrap_or_default())
}

/// Pobiera wskazane wydanie aplikacji, sprawdza jego sumę kontrolną, zakłada je na dysku
/// i stawia aplikację powłoki na nowo.
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
        // Wątek roboczy przepadł paniką; Operator ma to usłyszeć, nie zobaczyć zamilkły baner.
        Err(Odmowa::nowa(
            "przebieg-przerwany",
            format!(
                "Aktualizacja przerwała się nieoczekiwanie ({blad}). Nic nie zostało \
                 założone — spróbuj ponownie albo pobierz plik ręcznie ze strony „Pobierz”."
            ),
        ))
    })
}
