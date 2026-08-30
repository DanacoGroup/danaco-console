//! Moduł otwiera jedyne okno natywne powłoki, niosące wkompilowany pakiet interfejsu,
//! bez adresu zapasowego ani rdzenia.

use tauri::{AppHandle, Manager, State, WebviewUrl, WebviewWindow, WebviewWindowBuilder};

use crate::ustawienia::Ustawienia;

use crate::zamkniecie;

/// Etykieta okna głównego powłoki, jedynego okna natywnego aplikacji; ta sama wartość występuje w pliku uprawnień domyślnych aplikacji.
pub const ETYKIETA: &str = "glowne";

/// Tytuł okna głównego powłoki, widoczny na pasku systemu operacyjnego przy oknie oraz w jego własnej belce tytułowej.
pub const TYTUL: &str = "Danaco Console";

/// Otwiera okno główne powłoki z wkompilowanym pakietem interfejsu i wiąże jego zdarzenia systemowe z obsługą zamknięcia okna.
pub fn otworz(aplikacja: &AppHandle) -> tauri::Result<WebviewWindow> {
    /* Wskazanie rdzenia wchodzi w stronę przed jej wczytaniem. Pytanie o nie
       komendą kazałoby stronie czekać, a okno wejścia stoi puste, dopóki nie
       dostanie odpowiedzi — animacja startowa nie miałaby kiedy stanąć. */
    let wskazanie = aplikacja
        .try_state::<Ustawienia>()
        .and_then(|u: State<'_, Ustawienia>| u.adres_rdzenia_http())
        .unwrap_or_default();
    let okno = WebviewWindowBuilder::new(aplikacja, ETYKIETA, WebviewUrl::default())
        .initialization_script(&format!(
            "globalThis.DanacoAdresRdzenia = {};",
            serde_json::to_string(&wskazanie).unwrap_or_else(|_| "\"\"".to_string())
        ))
        .title(TYTUL)
        // Okno wejścia ma stałe 1040×780 punktów i musi zmieścić się wraz
        // z pasem działań — bez tego Operator nie ma czym zatwierdzić
        // logowania. Minimum powłoki jest więc miarą tego okna, nie liczbą
        // dobraną z ręki.
        .inner_size(1600.0, 1000.0)
        .min_inner_size(1080.0, 840.0)
        .resizable(true)
        // Belkę tytułową niesie samo okno aplikacji — wejściowe przed
        // uwierzytelnieniem i powłoka Centrum po nim — więc rama systemowa
        // stałaby nad nią drugi raz.
        .decorations(false)
        // Okno wstaje ukryte: silnik widoku rysuje pierwszą klatkę dopiero po
        // wczytaniu strony, a widoczne od razu pokazywałoby do tego czasu białe
        // pole. Strona pokazuje je sama, pierwszą klatką ekranu startowego.
        .visible(false)
        .center()
        .build()?;

    let kopia = okno.clone();
    okno.on_window_event(move |zdarzenie| zamkniecie::obsluz(&kopia, zdarzenie));
    Ok(okno)
}

/// Przywraca okno główne powłoki z ukrycia w zasobniku systemowym i ustawia na nim uwagę systemu okienkowego.
pub fn pokaz(aplikacja: &AppHandle) {
    if let Some(okno) = aplikacja.get_webview_window(ETYKIETA) {
        let _ = okno.show();
        let _ = okno.unminimize();
        let _ = okno.set_focus();
    }
}
