import type { PozycjaMenu } from '../komponenty/menu-drzewo';
import {
  nazwaKrotkaNakladu,
  nazwaNakladu,
  opisNakladu,
} from '../sterowanie/etykiety-sterowania';
import { KluczUstawieniaOkna, STOPNIE_NAKLADU } from '../sterowanie/klucze-ustawien';
import { utworzSterNastawy, type SterPaska } from './ster-nastawy';

/**
 * Wysiłek — ster paska zlecenia. W pasku stoi menu, a nie suwak: etykieta
 * komponentu ma być bieżącą wartością, a suwak pokazuje położenie i wymaga
 * podpisu obok siebie. Suwak w kolumnie sterowania zostaje — to ta sama
 * nastawa w dwóch widokach, nie dwa stany.
 *
 * Stopień jest napisem wyliczenia, nie liczbą. Wykaz stopni pochodzi z katalogu
 * rdzenia (`sterowanie/klucze-ustawien.ts` → `STOPNIE_NAKLADU`, tabela
 * `opcja_ustawienia`), a wartość idzie ustawieniem poziomu okna (`config.set`,
 * `naklad_rozumowania`), bo treść `window.update` nie ma dla niej pola.
 *
 * Napis pusty jest pełnoprawnym stopniem — znaczy „rozstrzyga kanał modelu",
 * a nie brak ustawienia, więc stoi na wykazie jak każdy inny.
 */

/** Nazwa zmiany w komunikacie — ta sama, którą wysyła suwak w kolumnie. */
const NAZWA = 'Nakład rozumowania';

/**
 * Przedrostek klucza pozycji.
 *
 * Klucz pusty należy w tym mechanizmie do stopki (`komponenty/menu-drzewo.ts`
 * → `stopka()`), a stopień „bez wskazania" jest napisem pustym. Przedrostek
 * rozdziela jedno od drugiego.
 */
const PRZEDROSTEK = 'naklad:';

/** Zależności steru — wąskie i wstrzykiwane. */
export interface ZaleznosciSteruNakladu {
  /** Stopień nakładu ze stanu potwierdzonego przez rdzeń. */
  migawka(): { naklad: string };
  /** Zapisuje ustawienie poziomu okna; odrzucenie niesie zdanie odmowy. */
  zapisz(nazwa: string, klucz: string, wartosc: unknown): Promise<void>;
}

export function utworzSterNakladu(zaleznosci: ZaleznosciSteruNakladu): SterPaska {
  const ster = utworzSterNastawy({
    nastawa: 'Wysiłek',
    ikona: 'aktywnosc',
    wykonaj: (klucz) =>
      zaleznosci.zapisz(
        NAZWA,
        KluczUstawieniaOkna.NakladRozumowania,
        klucz.slice(PRZEDROSTEK.length),
      ),
    odswiez: () => odswiez(),
  });

  function odswiez(): void {
    const stopien = zaleznosci.migawka().naklad;
    const drzewo: PozycjaMenu[] = STOPNIE_NAKLADU.map((pozycja) => ({
      rodzaj: 'wybor',
      klucz: `${PRZEDROSTEK}${pozycja}`,
      // Nazwa pełna w wykazie, krótka na uchwycie: w menu jest miejsce na
      // dopisek, w pasku go nie ma. Oba napisy z jednego słownika.
      nazwa: nazwaNakladu(pozycja),
      opis: opisNakladu(pozycja),
      wybrany: pozycja === stopien,
    }));
    ster.ustaw(nazwaKrotkaNakladu(stopien), drzewo);
  }

  odswiez();
  return { element: ster.element, odswiez };
}
