import { oknaPomocnicze, type OpisPomocniczego } from './rejestr-pomocniczych';
import { wytworniaPanelu } from './wytwornia-paneli';

/**
 * Plik rozstrzyga, które pozycje modułu naprawdę da się otworzyć obok rozmowy, bo rejestr sam niesie tylko stan opisowy, a kryterium otwieralności jest istnienie wytwórni, nie stan zbudowane w rejestrze.
 */
export interface PozycjaOtwieralna {
  /** Kod pozycji — ten sam, którym woła się `wytworniaPanelu`. */
  kod: string;
  /** Nazwa widziana przez Operatora. */
  nazwa: string;
  /** Po co Operatorowi ten panel. */
  przeznaczenie: string;
}

/** Pozycje modułu, które mają własną wytwórnię i dlatego naprawdę otworzą się po kliknięciu w menu okna rozmowy. */
export function paneleOtwieralne(kodModulu: string): readonly PozycjaOtwieralna[] {
  return oknaPomocnicze(kodModulu)
    .filter((pozycja) => wytworniaPanelu(pozycja.kod) !== null)
    .map((pozycja) => ({
      kod: pozycja.kod,
      nazwa: pozycja.nazwa,
      przeznaczenie: pozycja.przeznaczenie,
    }));
}

/**
 * Pozycje spisu bez wytwórni — wraz z powodem, wprost z rejestru.
 *
 * Oddawany jest cały opis rejestru, nie jego skrót: powód braku jest tu treścią
 * właściwą, a nie przypisem, i skrócenie go odebrałoby Operatorowi to, czym
 * może stan zmienić.
 */
export function paneleNieotwieralne(kodModulu: string): readonly OpisPomocniczego[] {
  return oknaPomocnicze(kodModulu).filter((pozycja) => wytworniaPanelu(pozycja.kod) === null);
}
