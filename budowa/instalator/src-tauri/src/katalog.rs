//! Moduł sprawdza prawo zapisu do katalogu docelowego próbą rzeczywistego zapisu
//! i odczytuje wolne miejsce na dysku, który ten katalog niesie — krok 4 kreatora
//! nie zgaduje, czy zapis się powiedzie, tylko go próbuje.

use crate::stan_maszyny;
use serde::Serialize;
use std::path::{Path, PathBuf};
use tauri::AppHandle;
use tauri_plugin_dialog::DialogExt;

/// Wynik sprawdzenia jednego katalogu docelowego: zapisywalność stwierdzona
/// próbą, nie założeniem, oraz wolne miejsce na dysku, który ten katalog niesie.
#[derive(Serialize, Clone, Debug)]
pub struct WynikKatalogu {
    pub sciezka: String,
    pub zapisywalny: bool,
    pub powod_odmowy: Option<String>,
    pub wolne_miejsce_bajty: Option<u64>,
}

/// Otwiera natywne okno wyboru katalogu i zwraca jego sprawdzenie; brak wyboru
/// (rezygnacja operatora) wraca jako `None`, nie jako usterka.
#[tauri::command]
pub async fn wybierz_katalog(aplikacja: AppHandle, tytul: String) -> Option<WynikKatalogu> {
    let (nadaj, mut odbierz) = tauri::async_runtime::channel(1);
    aplikacja
        .dialog()
        .file()
        .set_title(tytul)
        .pick_folder(move |wybor| {
            let _ = nadaj.blocking_send(wybor.map(|s| s.to_string()));
        });
    let sciezka = odbierz.recv().await.flatten()?;
    Some(sprawdz(&PathBuf::from(sciezka)))
}

/// Sprawdza podaną ścieżkę tekstową bez otwierania okna wyboru — użyte przy
/// odświeżeniu odczytu dla ścieżki proponowanej domyślnie.
#[tauri::command]
pub fn sprawdz_katalog(sciezka: String) -> WynikKatalogu {
    sprawdz(&PathBuf::from(sciezka))
}

fn sprawdz(sciezka: &Path) -> WynikKatalogu {
    let (zapisywalny, powod_odmowy) = probka_zapisu(sciezka);
    let dysk = if sciezka.exists() {
        sciezka.to_path_buf()
    } else {
        sciezka
            .parent()
            .map(Path::to_path_buf)
            .unwrap_or_else(|| sciezka.to_path_buf())
    };
    let (wolne, _) = stan_maszyny::wolne_miejsce(&dysk);
    WynikKatalogu {
        sciezka: sciezka.display().to_string(),
        zapisywalny,
        powod_odmowy,
        wolne_miejsce_bajty: wolne,
    }
}

/// Próba zapisu jest jedynym wiarygodnym sprawdzeniem prawa do katalogu:
/// atrybuty systemu plików kłamią przy udziałach sieciowych i profilach
/// przekierowanych, próba zapisu — nie. Katalog nieistniejący jest tworzony
/// na czas próby i usuwany natychmiast po niej, żeby sprawdzenie nie zostawiało śladu.
fn probka_zapisu(sciezka: &Path) -> (bool, Option<String>) {
    let utworzono_katalog = !sciezka.exists();
    if let Err(blad) = std::fs::create_dir_all(sciezka) {
        return (
            false,
            Some(format!("nie udało się utworzyć katalogu: {blad}")),
        );
    }
    let znacznik = sciezka.join(".danaco-sprawdzenie-zapisu");
    let wynik = match std::fs::write(&znacznik, b"") {
        Ok(()) => {
            let _ = std::fs::remove_file(&znacznik);
            (true, None)
        }
        Err(blad) => (false, Some(format!("odmowa zapisu: {blad}"))),
    };
    if utworzono_katalog {
        let _ = std::fs::remove_dir(sciezka);
    }
    wynik
}
