import { poleWyboru, ustawPozycje, type PozycjaWyboru } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import { KOD_MODULU } from './zrodlo-otoczenia';

/**
 * Ster modułu docelowego przeniesienia — jeden dla całego modułu Library. Wybór
 * pochodzi z katalogu modułów rdzenia, a nie z pola tekstowego, ponieważ
 * przeniesienie z nieznanym kodem modułu kończy się oknem, do którego nikt nie
 * wejdzie.
 */
export interface SterModulu {
  element: HTMLElement;
  /** Przerysowuje pozycje i wartość bieżącą ze wspólnego stanu. */
  odswiez(): void;
}

/**
 * Pozycja pusta stera: moduł wytwórcy pliku, czyli zachowanie zastane. Ster tę
 * drogę nazywa wprost, zamiast kazać jej się domyślać z pustego pola wyboru.
 */
export const POZYCJA_WYTWORCA: PozycjaWyboru = {
  wartosc: '',
  etykieta: 'Moduł wytwórcy pliku (sourceModuleId)',
};

/**
 * Pozycje stera złożone z katalogu modułów rdzenia. Moduł Library wypada
 * z wykazu: przeniesienie kompletu do modułu, w którym Operator już stoi,
 * założyłoby drugie okno tego samego modułu i nie otworzyłoby niczego nowego.
 */
export function pozycjeStera(stan: StanBiblioteki): PozycjaWyboru[] {
  return [
    POZYCJA_WYTWORCA,
    ...stan
      .moduly()
      .filter((modul) => modul.code !== KOD_MODULU)
      .map((modul) => ({ wartosc: modul.code, etykieta: `${modul.name} (${modul.code})` })),
  ];
}

export function utworzSterModulu(stan: StanBiblioteki, opisPola: string): SterModulu {
  const pole = poleWyboru(
    { etykieta: 'Moduł docelowy otwarcia', opis: opisPola },
    pozycjeStera(stan),
  );
  pole.kontrolka.dataset['ster'] = 'modul-docelowy';
  pole.kontrolka.addEventListener('change', () => {
    stan.ustawModulDocelowy(pole.kontrolka.value);
  });

  const powod = document.createElement('p');
  powod.className = 'dn-pole-opis ml-ster__powod';
  powod.dataset['powod'] = 'moduly';
  powod.hidden = true;

  const element = document.createElement('div');
  element.className = 'ml-ster';
  element.append(pole.element, powod);

  return {
    element,

    odswiez() {
      ustawPozycje(pole.kontrolka, pozycjeStera(stan));
      // Prawdą o nastawie jest stan modułu, nie kontrolka: zmienić ją mogło
      // drugie okno tego samego modułu.
      pole.kontrolka.value = stan.modulDocelowy();
      // Nastawa wskazująca moduł nieznany katalogowi zniknęłaby po cichu na
      // pozycję pustą.
      if (pole.kontrolka.value !== stan.modulDocelowy()) {
        stan.ustawModulDocelowy(pole.kontrolka.value);
      }
      const zdanie = stan.powodModulow();
      powod.textContent = zdanie;
      powod.hidden = zdanie === '';
    },
  };
}
