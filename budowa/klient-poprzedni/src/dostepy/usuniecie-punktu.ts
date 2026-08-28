import type { AccessPoint } from '../../../shared/contract';
import type { KomunikatCzynnosci } from './komunikat-czynnosci';
import type { StanDostepow } from './stan-dostepow';

/**
 * Usunięcie punktu dostępu jest czynnością nieodwracalną, więc przebiega
 * dwustopniowo: pierwsze naciśnięcie zapowiada skutek i zamienia przycisk
 * w potwierdzenie, a dopiero drugie wysyła komendę do rdzenia.
 */
export interface PrzyciskUsuniecia {
  /** Przycisk osadzany w nagłówku karty punktu. */
  element: HTMLButtonElement;
}

/**
 * Czas trwania zamiaru usunięcia, podany w milisekundach; po jego upływie
 * przycisk wraca do stanu wyjściowego, ponieważ przycisk zostawiony w stanie
 * potwierdzenia byłby pułapką dla kolejnego naciśnięcia w to samo miejsce.
 */
const TRWANIE_ZAMIARU = 6000;

export function utworzPrzyciskUsuniecia(
  punkt: AccessPoint,
  stan: StanDostepow,
  komunikat: KomunikatCzynnosci,
): PrzyciskUsuniecia {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn dn-btn--duch dn-btn--sm';
  element.setAttribute('aria-label', `Usuń punkt dostępu ${punkt.name}`);

  let zamiar = false;
  // Typ zegara pochodzi z `setTimeout`: przeglądarka zwraca liczbę,
  // a środowisko Node uchwyt `Timeout`.
  let zegar: ReturnType<typeof globalThis.setTimeout> | undefined;

  function stanWyjsciowy(): void {
    zamiar = false;
    element.textContent = 'Usuń punkt';
    element.className = 'dn-btn dn-btn--duch dn-btn--sm';
  }

  function zapowiedz(): void {
    zamiar = true;
    element.textContent = 'Potwierdź usunięcie';
    element.className = 'dn-btn dn-btn--niebezpieczny dn-btn--sm';
    komunikat.pokaz(
      `Usunięcie punktu „${punkt.name}" unieważni wszystkie nadania powołujące się na niego — także w innych oknach. Naciśnij ponownie, aby usunąć.`,
      false,
    );
    globalThis.clearTimeout(zegar);
    zegar = globalThis.setTimeout(() => {
      stanWyjsciowy();
      komunikat.pokaz('Usunięcie odwołane — punkt został nietknięty.', true);
    }, TRWANIE_ZAMIARU);
  }

  element.addEventListener('click', () => {
    if (!zamiar) {
      zapowiedz();
      return;
    }
    globalThis.clearTimeout(zegar);
    stanWyjsciowy();
    void stan.usunPunkt(punkt.id).then((wynik) => {
      komunikat.zWyniku(wynik, `Punkt „${punkt.name}" usunięty wraz z nadaniami.`);
    });
  });

  stanWyjsciowy();
  return { element };
}
