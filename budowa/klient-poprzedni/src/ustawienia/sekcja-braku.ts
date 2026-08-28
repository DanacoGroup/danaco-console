import type { SekcjaUstawien } from './sekcje';

/** Sekcja nazwanego braku jest wspólnym kształtem dla miejsc, których rdzeń nie pokrywa, i nie stawia ani jednego pola, przełącznika czy przycisku, który nie ma dokąd pójść. */
export interface OpisBraku {
  /** Zdanie otwierające: czego w tej sekcji nie ma i dlaczego. */
  wstep: string;
  /** Po jednym zdaniu na wiersz — co sprawdzono i z jakim wynikiem. */
  pomiar: readonly string[];
  /** Miejsce, do którego operator ma pójść po tę zdolność; pominięte znaczy, że takiego miejsca nie ma. */
  odeslanie?: string;
}

export function utworzSekcjeBraku(opis: OpisBraku): SekcjaUstawien {
  const element = document.createElement('div');
  element.className = 'du-sekcja du-sekcja--brak-pokrycia';

  const wstep = document.createElement('p');
  wstep.className = 'du-brak__wstep';
  wstep.textContent = opis.wstep;
  element.append(wstep);

  const naglowekPomiaru = document.createElement('p');
  naglowekPomiaru.className = 'dn-pole-opis du-granica';
  naglowekPomiaru.textContent = 'Co zmierzono:';
  element.append(naglowekPomiaru);

  const wykaz = document.createElement('ul');
  wykaz.className = 'du-brak__pomiar';
  wykaz.append(
    ...opis.pomiar.map((zdanie) => {
      const wiersz = document.createElement('li');
      wiersz.textContent = zdanie;
      return wiersz;
    }),
  );
  element.append(wykaz);

  if (opis.odeslanie !== undefined && opis.odeslanie !== '') {
    const odeslanie = document.createElement('p');
    odeslanie.className = 'dn-pole-opis du-granica';
    odeslanie.textContent = opis.odeslanie;
    element.append(odeslanie);
  }

  return {
    element,

    // Odświeżenie nie ma czego odczytać: przepisuje zdanie wstępu, żeby cisza nie wyglądała jak odczyt.
    odswiez() {
      wstep.textContent = opis.wstep;
    },

    rozlacz: () => undefined,
  };
}
