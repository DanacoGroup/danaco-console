/**
 * Wiersz rejestru akcji, który nie jest operacją kontekstową — bez wybieraka.
 *
 * `studio.contextual.op` przyjmuje `actionId` i nie sprawdza go względem rejestru
 * akcji: identyfikator spoza rejestru wraca odpowiedzią pomyślną z propozycją
 * zmiany. Rejestr zasięgu modułu Studio niesie także wiersze wskazujące komendy
 * okna komunikacji (`message.send`, `message.stop`, `message.list`). Gdyby były
 * wybieralne, „Uruchom operację" na pozycji „Wyślij" wysłałoby jej identyfikator
 * jako operację redakcyjną i wróciło wynikiem wyglądającym na udany. Wiersz
 * zostaje więc widoczny, bo rejestr go niesie, ale nie zostaje wybieralny.
 *
 * Plakietka źródła mówi, z czego wiersz pochodzi i którą komendę wskazuje.
 */
export function utworzPozycjeSpozaOperacji(
  id: string,
  nazwa: string,
  komenda: string,
): HTMLElement {
  const etykieta = document.createElement('span');
  etykieta.className = 'ms-wykaz__etykieta';
  etykieta.textContent = nazwa;

  const zrodlo = document.createElement('span');
  zrodlo.className = 'dn-plakietka ms-wykaz__zrodlo';
  zrodlo.textContent = `rejestr rdzenia · komenda ${komenda}`;
  zrodlo.title = `Wiersz rejestru ${id} wskazuje komendę ${komenda}, nie operację kontekstową Studio.`;

  const wiersz = document.createElement('div');
  wiersz.className = 'ms-wykaz__wiersz';
  wiersz.dataset['operacja'] = id;
  wiersz.dataset['wybieralna'] = 'nie';
  wiersz.append(etykieta, zrodlo);
  return wiersz;
}
