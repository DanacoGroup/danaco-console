/**
 * Odczyt treści pliku z urządzenia do postaci base64 — jedyne miejsce, w którym
 * moduł Library zamienia plik Operatora w ładunek kontraktu. Nieudany odczyt
 * wraca jako `null`, a nie wyjątkiem wywracającym widok.
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
