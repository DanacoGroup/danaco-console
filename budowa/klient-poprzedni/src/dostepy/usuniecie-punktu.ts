import type { AccessPoint } from '../../../shared/contract';
import type { KomunikatCzynnosci } from './komunikat-czynnosci';
import type { StanDostepow } from './stan-dostepow';

/**
 * Usunięcie punktu dostępu — czynność nieodwracalna, więc dwustopniowa.
 *
 * Usunięcie punktu unieważnia WSZYSTKIE nadania, które się na niego powoływały,
 * także w oknach, których Operator w tej chwili nie widzi. Pierwsze naciśnięcie
 * mówi więc, co się stanie, i zamienia przycisk w potwierdzenie; drugie wysyła
 * komendę. Okna `confirm` nie używamy: blokuje wątek dokumentu i nie da się go
 * ubrać w warstwę wizualną platformy.
 *
 * Zamiar wygasa sam po chwili. Przycisk zostawiony w stanie „potwierdź" byłby
 * pułapką dla kolejnego kliknięcia w to samo miejsce.
 */
export interface PrzyciskUsuniecia {
  /** Przycisk osadzany w nagłówku karty punktu. */
  element: HTMLButtonElement;
}

/** Ile milisekund trwa zamiar usunięcia, zanim przycisk wróci do stanu wyjściowego. */
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
  // Typ zegara bierzemy z samego `setTimeout`, a nie z `number`: przeglądarka
  // zwraca liczbę, środowisko Node — uchwyt `Timeout`, a w drzewie obecne są
  // definicje obu.
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
