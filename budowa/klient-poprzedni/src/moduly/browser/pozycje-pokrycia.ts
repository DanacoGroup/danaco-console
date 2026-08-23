import type { PokrycieKomend } from '../pokrycie-komend';
import type { PozycjaBezObslugi } from './etykiety-browser';

/**
 * Pozycje modułu Browser, których okno jeszcze nie wykonuje, wraz z odczytem
 * ich pokrycia w rdzeniu.
 *
 * Jedna odpowiedzialność: postawić taką pozycję jako pełnoprawny przycisk
 * i odpowiedzieć na naciśnięcie powodem wziętym z odczytu, nie z napisu.
 *
 * Powód rozstrzyga byt wspólny wszystkich modułów (`moduly/pokrycie-komend.ts`):
 * pyta rdzeń o wykaz jego komend powitaniem `connection.hello` i mówi osobno
 * „kontrakt tej komendy nie ma", osobno „kontrakt ma, rdzeń nie ma uchwytu",
 * osobno „rdzeń ma uchwyt, a to okno go jeszcze nie wywołuje". Zdanie wpisane
 * w moduł na sztywno kłamałoby w dniu, w którym rdzeń dostanie obsługę —
 * i nikt by go nie zdjął, bo nic go z rdzeniem nie łączy.
 *
 * Przycisk nie jest wygaszony i nie milczy: naciśnięcie daje dymek z powodem
 * oraz wiersz odpowiedzi okna. Nic nie idzie do rdzenia i nic nie udaje
 * wykonania. Znacznik `data-brak-komendy` i `data-pokrycie` niesie kontrolka
 * biblioteczna — po nich sięgają sprawdziany produktu.
 */
export function przyciskPozycji(
  pokrycie: PokrycieKomend,
  pozycja: PozycjaBezObslugi,
  etykieta: string,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): HTMLButtonElement {
  const element = pokrycie.przycisk(etykieta, pozycja.komenda, pozycja.czynnosc);
  element.classList.add('dn-btn--sm');
  // Powód czytany z `title` dopiero przy naciśnięciu, tak samo jak robi to dymek
  // kontrolki: byt przerysowuje go po każdej zmianie wykazu, więc domknięcie na
  // treści z chwili budowy mówiłoby „odczyt w toku" długo po odpowiedzi rdzenia.
  element.addEventListener('click', () => powiedz(element.title, false));
  return element;
}
