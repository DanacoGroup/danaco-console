import { poleWyboru, ustawPozycje, type PozycjaWyboru } from '../../modele/kontrolki-formularza';
import type { StanBiblioteki } from './stan-biblioteki';
import { KOD_MODULU } from './zrodlo-otoczenia';

/**
 * Ster modułu docelowego przeniesienia — jeden dla całego modułu Library.
 *
 * Wybór z katalogu, nie pole tekstowe: `context.transfer` z nieznanym
 * `targetModuleId` wraca powodzeniem, a rdzeń zakłada okno z dokładnie tym
 * kodem (`zapisy-zbiorcze.ts`), więc literówka kończy się oknem, do którego
 * nikt nie wejdzie. Katalog bierze się z `module.list`.
 *
 * Wybór zmienia nastawę we wspólnym stanie modułu (`stan.modulDocelowy()`).
 * Nastawa jest jedna, więc zmiana dokonana w Library Explorerze jest widoczna
 * w File Preview i odwrotnie.
 *
 * Pozycja pusta zostaje: „moduł wytwórcy pliku" bierze `sourceModuleId`
 * przeniesionego pliku, a ster tę drogę nazywa, zamiast kazać jej się domyślać
 * z pustego pola.
 *
 * Gdy `module.list` odmówi, ster ma samą pozycję pustą, a pod nim stoi powód
 * odmowy — milcząca lista z jedną pozycją wyglądałaby jak platforma z jednym
 * modułem.
 */
export interface SterModulu {
  element: HTMLElement;
  /** Przerysowuje pozycje i wartość bieżącą ze wspólnego stanu. */
  odswiez(): void;
}

/** Pozycja pusta stera — moduł wytwórcy pliku, czyli zachowanie zastane. */
export const POZYCJA_WYTWORCA: PozycjaWyboru = {
  wartosc: '',
  etykieta: 'Moduł wytwórcy pliku (sourceModuleId)',
};

/**
 * Pozycje stera złożone z katalogu rdzenia.
 *
 * Library wypada z wykazu: przeniesienie kompletu do modułu, w którym Operator
 * już stoi, założyłoby drugie okno tego samego modułu i nie otworzyłoby niczego
 * nowego. Tak samo zawęża katalog moduł Design (`zapis-designu.ts`).
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
      // Wartość bierze się ze stanu, nie z kontrolki: `ustawPozycje` utrzymuje
      // poprzedni wybór, o ile pozycja nadal istnieje, ale prawdą o nastawie
      // jest stan modułu — także wtedy, gdy zmieniło ją drugie okno.
      pole.kontrolka.value = stan.modulDocelowy();
      // Nastawa wskazująca moduł, którego katalog już nie zna, zniknęłaby po
      // cichu na pozycję pustą — czyli przeniesienie poszłoby gdzie indziej,
      // niż mówił ster przed chwilą.
      if (pole.kontrolka.value !== stan.modulDocelowy()) {
        stan.ustawModulDocelowy(pole.kontrolka.value);
      }
      const zdanie = stan.powodModulow();
      powod.textContent = zdanie;
      powod.hidden = zdanie === '';
    },
  };
}
