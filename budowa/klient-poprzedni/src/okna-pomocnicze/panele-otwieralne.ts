import { oknaPomocnicze, type OpisPomocniczego } from './rejestr-pomocniczych';
import { wytworniaPanelu } from './wytwornia-paneli';

/**
 * Co z pozycji modułu da się naprawdę otworzyć obok rozmowy — i co nie.
 *
 * Menu `⋮` w nagłówku okna rozmowy potrzebuje wykazu pozycji, które po
 * kliknięciu się otworzą. Rejestr sam tego nie mówi: niesie stan opisowy, nie
 * zdolność wykonawczą.
 *
 * Kryterium jest istnienie wytwórni, a nie stan `zbudowane`. Stan `zbudowane`
 * jest zdaniem rejestru o produkcie, istnienie wytwórni faktem o kodzie; gdy
 * się rozjadą, menu jest krótsze (pozycja nie wchodzi), a rozjazd zgłasza
 * `pas-pomocniczych.ts` przez `console.warn`.
 *
 * Pozycja bez wytwórni nie wchodzi do `paneleOtwieralne` w ogóle — nie dostaje
 * wiersza wygaszonego ani „wkrótce", bo nieczynny wiersz jest bramką. Jej
 * miejsce jest w `paneleNieotwieralne`, gdzie stoi wraz z powodem wprost
 * z rejestru: czego brakuje i czyja to robota.
 *
 * Plik nie buduje paneli, oddaje opisy; nie zna DOM; nie rozstrzyga, ile paneli
 * wolno otworzyć naraz — panel jest bytem otwieranym i zamykanym pojedynczo,
 * a limitu nie ma.
 *
 * Dla modułu spoza rejestru `oknaPomocnicze` oddaje wykaz pusty i obie funkcje
 * oddają puste. Zdanie o braku spisu należy do gospodarza, który to pokazuje
 * (`zdanieBezSpisu` w `pas-pomocniczych.ts`).
 */

/** Pozycja, którą menu ⋮ może otworzyć. */
export interface PozycjaOtwieralna {
  /** Kod pozycji — ten sam, którym woła się `wytworniaPanelu`. */
  kod: string;
  /** Nazwa widziana przez Operatora. */
  nazwa: string;
  /** Po co Operatorowi ten panel. */
  przeznaczenie: string;
}

/** Pozycje modułu, które mają wytwórnię — czyli otworzą się naprawdę. */
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
