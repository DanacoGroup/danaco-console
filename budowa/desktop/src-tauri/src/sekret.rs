//! Moduł wytwarza sekret nawiązania gniazda: wartość losową jednego uruchomienia
//! powłoki, którą powłoka podaje rdzeniowi startowanemu obok siebie i stronie
//! interfejsu otwierającej gniazdo.

/// Ile bajtów losowych niesie sekret przed zapisem szesnastkowym. Sekret idzie
/// w zapytaniu adresu gniazda, więc musi być krótki, a zarazem nie do
/// odgadnięcia — 32 bajty dają 256 bitów wejścia.
const BAJTOW: usize = 32;

/// Wytwarza sekret nawiązania na to uruchomienie powłoki. Źródło losowości daje
/// system operacyjny; jego odmowa jest awarią startu, nie powodem do sekretu
/// przewidywalnego, bo sekret przewidywalny nie broni gniazda przed nikim.
pub fn wytworz() -> String {
    let mut bajty = [0u8; BAJTOW];
    getrandom::fill(&mut bajty).expect("źródło losowości systemu dla sekretu nawiązania");
    zapis_szesnastkowy(&bajty)
}

/// Zapisuje bajty szesnastkowo małymi literami; sekret idzie w adresie gniazda,
/// więc nie może nieść znaków wymagających kodowania procentowego.
fn zapis_szesnastkowy(bajty: &[u8]) -> String {
    let mut wynik = String::with_capacity(bajty.len() * 2);
    for bajt in bajty {
        wynik.push_str(&format!("{bajt:02x}"));
    }
    wynik
}

