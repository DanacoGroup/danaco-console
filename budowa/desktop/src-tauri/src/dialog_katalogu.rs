//! Natywne okno wyboru katalogu — jedno na dwa zastosowania.
//!
//! Katalogi robocze są listą na oknie komunikacji, a ich wskazanie
//! musi być czynnością systemu operacyjnego, nie polem tekstowym. Wybór jest
//! dostępny dwiema drogami: poleceniem z interfejsu i pozycją w zasobniku.
//! Rezygnacja z wyboru zwraca brak i niczego nie zmienia.
//!
//! Interfejs wskazuje katalog w dwóch sprawach: wskazania katalogu roboczego
//! (ustawienie `katalog.roboczy.podstawa`) oraz dodania katalogu jako punktu
//! dostępu. Dostęp mówi, do czego model sięga, katalog roboczy — gdzie zostawia
//! swoje pliki, ale czynność systemu operacyjnego jest ta sama, więc okno jest
//! jedno i różni je wyłącznie napis w belce. Napis podaje wywołujący, bo to on
//! zna sprawę; jego brak daje wartość domyślną, a nie odmowę czynności.
//!
//! Konsument po stronie interfejsu: `client/src/powloka/most-katalogow.ts`.

use tauri::{AppHandle, Emitter};
use tauri_plugin_dialog::DialogExt;

/// Zdarzenie powłoki niosące wybrany katalog. Nazwa własna powłoki
/// z przedrostkiem `powloka:` — nie należy do kontraktu WebSocket z `shared/`
/// i celowo nie powiela żadnej jego nazwy.
pub const ZDARZENIE_KATALOG: &str = "powloka:katalog-roboczy";

/// Napis w belce okna wyboru, gdy wywołujący żadnego nie podał.
pub const TYTUL_DOMYSLNY: &str = "Wskaż katalog roboczy";

/// Otwiera okno wyboru i zwraca wskazaną ścieżkę.
pub async fn wybierz(aplikacja: &AppHandle, tytul: Option<&str>) -> Option<String> {
    let (nadaj, mut odbierz) = tauri::async_runtime::channel(1);
    aplikacja
        .dialog()
        .file()
        .set_title(rozstrzygnij_tytul(tytul))
        .pick_folder(move |wybor| {
            let _ = nadaj.blocking_send(wybor.map(|sciezka| sciezka.to_string()));
        });
    odbierz.recv().await.flatten()
}

/// Otwiera okno wyboru z zasobnika i rozgłasza wybór do okna interfejsu.
pub fn wybierz_i_rozglos(aplikacja: &AppHandle) {
    let kopia = aplikacja.clone();
    aplikacja
        .dialog()
        .file()
        .set_title(TYTUL_DOMYSLNY)
        .pick_folder(move |wybor| {
            if let Some(sciezka) = wybor {
                let _ = kopia.emit(ZDARZENIE_KATALOG, sciezka.to_string());
            }
        });
}

/// Napis podany przez wywołującego, po odrzuceniu wartości pustej.
fn rozstrzygnij_tytul(tytul: Option<&str>) -> &str {
    match tytul.map(str::trim) {
        Some(napis) if !napis.is_empty() => napis,
        _ => TYTUL_DOMYSLNY,
    }
}
