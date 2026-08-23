/**
 * Odczyt treści pliku z urządzenia do postaci base64 — jedyne miejsce, w którym
 * moduł Library zamienia plik Operatora w ładunek kontraktu.
 *
 * Wspólne dla wgrania pliku (`library.file.upload`) i dołożenia wersji
 * (`library.version.add`) — obie komendy niosą to samo pole `contentBase64`.
 *
 * Nieudany odczyt (plik zniknął, dostęp odrzucony) wraca jako `null`, a nie
 * wyjątkiem wywracającym widok; wołający ma wtedy powiedzieć, że nic nie
 * zostało wysłane.
 */
export function odczytajBase64(plik: File): Promise<string | null> {
  return new Promise((rozstrzygnij) => {
    const czytnik = new FileReader();
    czytnik.onerror = () => rozstrzygnij(null);
    czytnik.onload = () => {
      const wynik = typeof czytnik.result === 'string' ? czytnik.result : '';
      const przecinek = wynik.indexOf(',');
      rozstrzygnij(przecinek === -1 ? null : wynik.slice(przecinek + 1));
    };
    czytnik.readAsDataURL(plik);
  });
}
