//! Moduł otwiera natywne okno wyboru katalogu, wspólne dla wskazania katalogu roboczego
//! z interfejsu i dodania punktu dostępu z zasobnika, i różnicuje zastosowania napisem w belce.

use tauri::{AppHandle, Emitter};
use tauri_plugin_dialog::DialogExt;

/// Zdarzenie powłoki niosące wybrany katalog. Nazwa własna powłoki
/// z przedrostkiem `powloka:` — nie należy do kontraktu WebSocket z `shared/`
/// i celowo nie powiela żadnej jego nazwy.
pub const ZDARZENIE_KATALOG: &str = "powloka:katalog-roboczy";

/// Napis w belce okna wyboru katalogu, używany wtedy, gdy wywołujący nie podał własnego tytułu okna wyboru.
pub const TYTUL_DOMYSLNY: &str = "Wskaż katalog roboczy";

/// Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę, oddając brak przy rezygnacji z wyboru.
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

/// Otwiera natywne okno wyboru katalogu z pozycji zasobnika i rozgłasza wybór katalogu do okna interfejsu.
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

/// Ustala napis w belce okna wyboru, podany przez wywołującego, po odrzuceniu wartości pustej lub białej.
fn rozstrzygnij_tytul(tytul: Option<&str>) -> &str {
    match tytul.map(str::trim) {
        Some(napis) if !napis.is_empty() => napis,
        _ => TYTUL_DOMYSLNY,
    }
}
