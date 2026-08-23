import type { SekcjaUstawien } from './sekcje';

/**
 * Sekcja nazwanego braku — wspólny kształt dla miejsc, których rdzeń nie
 * pokrywa.
 *
 * Dwie sekcje Okna Ustawień (Konto Operatora, Powiadomienia) mówią tę samą
 * rzecz: czego w rdzeniu nie ma i dlaczego kontrolki tu nie ma. Wspólna
 * wytwórnia trzyma ten kształt w jednym miejscu, zamiast w dwóch odmianach.
 *
 * Sekcja nie stawia ani jednego pola, przełącznika czy przycisku, który nie ma
 * dokąd pójść. Kontrolka wyłączona byłaby tu gorsza niż jej brak: stan
 * „wyłączony" należy się elementowi, który w danym miejscu nie ma sensu, a nie
 * zapowiedzi czegoś, czego nie ma.
 *
 * Zaślepka mówi „wkrótce" — ta sekcja mówi, co dokładnie sprawdzono i gdzie,
 * z nazwami plików rdzenia i migracji, żeby dało się to zweryfikować. Zdanie
 * fałszywe albo nieaktualne waży tu więcej niż jego brak: Operator czyta je
 * jako wiedzę o stanie produktu.
 */
export interface OpisBraku {
  /** Zdanie otwierające: czego w tej sekcji nie ma i dlaczego. */
  wstep: string;
  /** Po jednym zdaniu na wiersz — co sprawdzono i z jakim wynikiem. */
  pomiar: readonly string[];
  /**
   * Miejsce, do którego Operator ma pójść po tę zdolność.
   *
   * Pominięte znaczy „nie ma takiego miejsca" — i wtedy sekcja tego nie udaje.
   * Odesłanie stoi tu zamiast drugiej ramy do tych samych danych.
   */
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

    // Odświeżenie nie ma czego odczytać: przepisuje zdanie wstępu, bo cisza
    // wyglądałaby jak odczyt, który nic nie znalazł.
    odswiez() {
      wstep.textContent = opis.wstep;
    },

    rozlacz: () => undefined,
  };
}
