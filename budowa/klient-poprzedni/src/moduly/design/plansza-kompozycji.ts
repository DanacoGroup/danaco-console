import type { DesignBoardLayer } from '../../../../shared/contract';
import { BOK_WYJSCIOWY } from './zapis-kompozycji';
import type { StanKompozycji } from './stan-kompozycji';

/**
 * Kanwa swobodna Design Board — warstwy, powiększenie, przesunięcie, siatka.
 *
 * Odpowiada wyłącznie za wyrysowanie kompozycji i przyjęcie wskazań myszą.
 *
 * Przeciąganie jest tu miejscowe: wskaźnik przechwycony na warstwie
 * (`setPointerCapture`), przesunięcie liczone w jednostkach kompozycji
 * (podzielone przez powiększenie), warstwa zablokowana nie rusza się wcale.
 *
 * Plansza nie zna rdzenia. Ruch warstwy zmienia zapis kompozycji; do rdzenia
 * jedzie dopiero zapis całości komendą `design.board.update` — z okna, nie stąd.
 */
export interface PlanszaKompozycji {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzPlansze(stan: StanKompozycji): PlanszaKompozycji {
  const plotno = document.createElement('div');
  plotno.className = 'md-plansza__plotno';

  const element = document.createElement('div');
  element.className = 'md-plansza';
  element.setAttribute('aria-label', 'Kanwa kompozycji Design Board');
  element.append(plotno);

  // Kliknięcie w puste płótno zdejmuje zaznaczenie — inaczej nie da się wyjść
  // z zaznaczenia wielokrotnego bez trafienia w warstwę.
  element.addEventListener('pointerdown', (zdarzenie) => {
    if (zdarzenie.target === element || zdarzenie.target === plotno) stan.odznacz();
  });

  function odswiez(): void {
    const widok = stan.widok();
    element.dataset['siatka'] = String(widok.siatka);
    plotno.style.transform =
      `translate(${widok.przesuniecieX}px, ${widok.przesuniecieY}px) scale(${widok.powiekszenie})`;
    plotno.replaceChildren(
      ...stan.warstwy().map((warstwa) =>
        elementWarstwy(warstwa, stan.zaznaczone().includes(warstwa.id), stan),
      ),
    );
  }

  return { element, odswiez };
}

/** Jedna warstwa na kanwie wraz z jej przeciąganiem. */
function elementWarstwy(
  warstwa: DesignBoardLayer,
  zaznaczona: boolean,
  stan: StanKompozycji,
): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'md-warstwa';
  element.dataset['warstwa'] = warstwa.id;
  element.dataset['zaznaczona'] = String(zaznaczona);
  element.dataset['zablokowana'] = String(warstwa.locked === true);
  element.setAttribute('aria-pressed', String(zaznaczona));
  element.style.left = `${warstwa.x ?? 0}px`;
  element.style.top = `${warstwa.y ?? 0}px`;
  element.style.width = `${warstwa.width ?? BOK_WYJSCIOWY}px`;
  element.style.height = `${warstwa.height ?? BOK_WYJSCIOWY}px`;
  element.textContent = warstwa.note ?? warstwa.assetId ?? warstwa.id;

  element.addEventListener('click', (zdarzenie) => {
    stan.zaznacz(warstwa.id, zdarzenie.ctrlKey || zdarzenie.metaKey);
  });
  dopnijPrzeciaganie(element, warstwa, stan);
  return element;
}

/** Przeciąganie warstwy wskaźnikiem; warstwa zablokowana zostaje w miejscu. */
function dopnijPrzeciaganie(
  element: HTMLElement,
  warstwa: DesignBoardLayer,
  stan: StanKompozycji,
): void {
  let poczatek: { x: number; y: number } | null = null;

  element.addEventListener('pointerdown', (zdarzenie) => {
    if (warstwa.locked === true) return;
    poczatek = { x: zdarzenie.clientX, y: zdarzenie.clientY };
    element.setPointerCapture(zdarzenie.pointerId);
  });

  element.addEventListener('pointermove', (zdarzenie) => {
    if (poczatek === null) return;
    const skala = stan.widok().powiekszenie || 1;
    stan.przesun(
      warstwa.id,
      (warstwa.x ?? 0) + (zdarzenie.clientX - poczatek.x) / skala,
      (warstwa.y ?? 0) + (zdarzenie.clientY - poczatek.y) / skala,
    );
  });

  element.addEventListener('pointerup', () => {
    poczatek = null;
  });
}
