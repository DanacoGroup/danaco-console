//! Moduł czyta rzeczywisty stan maszyny operatora: wydanie systemu, architekturę
//! procesora, wolne miejsce na dysku systemowym oraz poziom uprawnień procesu.
//! Wartość, której nie da się odczytać w bieżącym środowisku, wraca jako
//! nazwana odmowa w polu `*_odmowa`, nigdy jako wartość zgadywana.

use serde::Serialize;
use std::path::{Path, PathBuf};

/// Odczyt stanu maszyny przekazywany do okna kreatora poleceniem `stan_maszyny`.
#[derive(Serialize, Clone, Debug)]
pub struct StanMaszyny {
    pub system: String,
    pub procesor: String,
    pub wolne_miejsce_bajty: Option<u64>,
    pub wolne_miejsce_odmowa: Option<String>,
    pub uprawnienia: String,
    pub katalog_programu: Option<String>,
    pub katalog_programu_odmowa: Option<String>,
    pub katalog_danych: Option<String>,
    pub katalog_danych_odmowa: Option<String>,
}

/// Architektura procesora w słowniku oczekiwanym przez krok 3 prototypu
/// (`x64` albo `arm`); architektura nierozpoznana wraca jako `brak`, żeby
/// instalator nie milczał w wartość domyślną, której maszyna nie potwierdziła.
pub fn architektura_kreatora() -> &'static str {
    match std::env::consts::ARCH {
        "x86_64" => "x64",
        "aarch64" => "arm",
        _ => "brak",
    }
}

/// Poleca oknu kreatora pełny odczyt stanu tej maszyny.
#[tauri::command]
pub fn stan_maszyny() -> StanMaszyny {
    odczytaj()
}

/// Składa pełny odczyt stanu tej maszyny.
pub fn odczytaj() -> StanMaszyny {
    let info = os_info::get();
    let katalog_domyslny = katalog_programu_domyslny();
    let (wolne_bajty, wolne_odmowa) = wolne_miejsce(&katalog_docelowego_dysku(&katalog_domyslny));

    StanMaszyny {
        system: format!("{} {}", info.os_type(), info.version()),
        procesor: architektura_kreatora().to_string(),
        wolne_miejsce_bajty: wolne_bajty,
        wolne_miejsce_odmowa: wolne_odmowa,
        uprawnienia: poziom_uprawnien(),
        katalog_programu: katalog_domyslny.as_ref().map(|p| p.display().to_string()),
        katalog_programu_odmowa: katalog_domyslny
            .is_none()
            .then(|| "zmienna środowiskowa LOCALAPPDATA nie jest ustawiona".to_string()),
        katalog_danych: katalog_danych_domyslny().map(|p| p.display().to_string()),
        katalog_danych_odmowa: katalog_danych_domyslny()
            .is_none()
            .then(|| "zmienna środowiskowa APPDATA nie jest ustawiona".to_string()),
    }
}

/// Katalog programu proponowany domyślnie: profil operatora, tak jak w treści
/// kroku 4. Bez zmiennej środowiskowej niosącej ten profil instalator nie zgaduje ścieżki.
fn katalog_programu_domyslny() -> Option<PathBuf> {
    #[cfg(windows)]
    {
        std::env::var_os("LOCALAPPDATA")
            .map(|d| PathBuf::from(d).join("Programs").join("Danaco Console"))
    }
    #[cfg(not(windows))]
    {
        std::env::var_os("LOCALAPPDATA").map(|d| PathBuf::from(d).join("Danaco Console"))
    }
}

/// Katalog danych proponowany domyślnie, odpowiednik `Roaming` w profilu operatora.
fn katalog_danych_domyslny() -> Option<PathBuf> {
    #[cfg(windows)]
    {
        std::env::var_os("APPDATA").map(|d| PathBuf::from(d).join("Danaco Console"))
    }
    #[cfg(not(windows))]
    {
        std::env::var_os("APPDATA").map(|d| PathBuf::from(d).join("Danaco Console"))
    }
}

/// Najbliższy istniejący przodek ścieżki, na którym da się zmierzyć wolne miejsce —
/// katalog docelowy zwykle jeszcze nie istnieje przed instalacją.
fn katalog_docelowego_dysku(katalog: &Option<PathBuf>) -> PathBuf {
    let punkt = katalog
        .clone()
        .or_else(|| std::env::current_dir().ok())
        .unwrap_or_else(|| PathBuf::from("."));
    najblizszy_istniejacy(&punkt)
}

fn najblizszy_istniejacy(sciezka: &Path) -> PathBuf {
    let mut biezacy = sciezka.to_path_buf();
    loop {
        if biezacy.exists() {
            return biezacy;
        }
        match biezacy.parent() {
            Some(rodzic) => biezacy = rodzic.to_path_buf(),
            None => return PathBuf::from("."),
        }
    }
}

/// Wolne miejsce na dysku niosącym podaną ścieżkę, w bajtach; brak odczytu
/// wraca jako nazwana odmowa zamiast wartości zerowej albo zmyślonej.
pub fn wolne_miejsce(sciezka: &Path) -> (Option<u64>, Option<String>) {
    match wolne_miejsce_platformowo(sciezka) {
        Some(bajty) => (Some(bajty), None),
        None => (
            None,
            Some(format!(
                "nie udało się odczytać wolnego miejsca na dysku dla ścieżki {}",
                sciezka.display()
            )),
        ),
    }
}

#[cfg(unix)]
fn wolne_miejsce_platformowo(sciezka: &Path) -> Option<u64> {
    use std::ffi::CString;
    use std::mem::MaybeUninit;
    use std::os::unix::ffi::OsStrExt;

    let c = CString::new(sciezka.as_os_str().as_bytes()).ok()?;
    unsafe {
        let mut stan: MaybeUninit<libc::statvfs> = MaybeUninit::uninit();
        if libc::statvfs(c.as_ptr(), stan.as_mut_ptr()) != 0 {
            return None;
        }
        let stan = stan.assume_init();
        Some(stan.f_bavail * stan.f_frsize)
    }
}

#[cfg(windows)]
fn wolne_miejsce_platformowo(sciezka: &Path) -> Option<u64> {
    use std::os::windows::ffi::OsStrExt;
    use windows_sys::Win32::Storage::FileSystem::GetDiskFreeSpaceExW;

    let mut szeroki: Vec<u16> = sciezka.as_os_str().encode_wide().collect();
    szeroki.push(0);
    let mut dostepne: u64 = 0;
    let wynik = unsafe {
        GetDiskFreeSpaceExW(
            szeroki.as_ptr(),
            &mut dostepne,
            std::ptr::null_mut(),
            std::ptr::null_mut(),
        )
    };
    if wynik == 0 {
        None
    } else {
        Some(dostepne)
    }
}

/// Poziom uprawnień bieżącego procesu, tak jak go widzi ten system operacyjny.
#[cfg(unix)]
fn poziom_uprawnien() -> String {
    if unsafe { libc::geteuid() } == 0 {
        "konto administratora (root)".to_string()
    } else {
        "konto użytkownika".to_string()
    }
}

#[cfg(windows)]
fn poziom_uprawnien() -> String {
    use windows_sys::Win32::Foundation::{CloseHandle, HANDLE};
    use windows_sys::Win32::Security::{
        GetTokenInformation, TokenElevation, TOKEN_ELEVATION, TOKEN_QUERY,
    };
    use windows_sys::Win32::System::Threading::{GetCurrentProcess, OpenProcessToken};

    unsafe {
        let mut token: HANDLE = std::ptr::null_mut();
        if OpenProcessToken(GetCurrentProcess(), TOKEN_QUERY, &mut token) == 0 {
            return "nie udało się odczytać uprawnień procesu".to_string();
        }
        let mut podniesione = TOKEN_ELEVATION { TokenIsElevated: 0 };
        let mut zwrocone: u32 = 0;
        let odczytano = GetTokenInformation(
            token,
            TokenElevation,
            &mut podniesione as *mut _ as *mut _,
            std::mem::size_of::<TOKEN_ELEVATION>() as u32,
            &mut zwrocone,
        );
        CloseHandle(token);
        if odczytano == 0 {
            return "nie udało się odczytać uprawnień procesu".to_string();
        }
        if podniesione.TokenIsElevated != 0 {
            "konto administratora (uprawnienia podniesione)".to_string()
        } else {
            "konto użytkownika".to_string()
        }
    }
}
