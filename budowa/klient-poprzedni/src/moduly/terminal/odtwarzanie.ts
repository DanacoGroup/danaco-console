/**
 * Odtwarzanie sesji okna Output Console: przesuwanie granicy widocznych wierszy
 * bufora, w całości po stronie klienta, bez pytania rdzenia o cokolwiek.
 */
export interface Odtwarzacz {
  /**
   * Ostatni wiersz brany do rysowania; liczba ujemna znaczy suwak nietknięty.
   */
  polozenie(): number;
  element: HTMLElement;
  /** Ustawia liczbę wierszy bufora i dosuwa suwak, gdy nie był ruszany. */
  ustawZakres(wierszy: number): void;
}

export function odtwarzaczSesji(przyZmianie: () => void): Odtwarzacz {
  const suwak = document.createElement('input');
  suwak.type = 'range';
  suwak.className = 'dt-suwak';
  suwak.min = '0';
  suwak.max = '0';
  suwak.value = '0';
  suwak.setAttribute('aria-label', 'Odtwarzanie sesji — położenie w historii wyjścia');

  const opis = document.createElement('span');
  opis.className = 'dn-pole-opis';
  opis.textContent = 'Odtwarzanie: teraz';

  const element = document.createElement('span');
  element.className = 'dt-odtwarzanie';
  element.append(suwak, opis);

  let ruszony = false;
  let wierszy = 0;

  suwak.addEventListener('input', () => {
    ruszony = Number(suwak.value) < wierszy;
    opis.textContent = ruszony ? `Odtwarzanie: wiersz ${suwak.value} z ${wierszy}` : 'Odtwarzanie: teraz';
    przyZmianie();
  });

  return {
    element,

    polozenie: () => (ruszony ? Number(suwak.value) : -1),

    ustawZakres(nowe) {
      wierszy = nowe;
      suwak.max = String(nowe);
      if (!ruszony) {
        suwak.value = String(nowe);
        opis.textContent = 'Odtwarzanie: teraz';
      }
    },
  };
}
