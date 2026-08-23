/**
 * Odtwarzanie sesji (playback) okna Output Console.
 *
 * Odtwarzanie jest czynnością w całości kliencką: wiersze wyjścia już leżą
 * w buforze wraz ze znacznikami czasu, więc powtórzenie przebiegu jest
 * przesuwaniem granicy widocznych wierszy, a nie ponownym pytaniem rdzenia.
 * Kontrakt nie ma komendy odtwarzania i mieć nie musi.
 *
 * Suwak stoi w położeniu ostatniego wiersza, dopóki Operator go nie ruszy —
 * konsola ma domyślnie pokazywać teraźniejszość, nie przeszłość.
 */
export interface Odtwarzacz {
  /**
   * Ostatni wiersz brany do rysowania; liczba ujemna znaczy „suwak nietknięty,
   * pokaż wszystko”.
   *
   * Wartością „wszystko” nie może być zero, bo zero jest także poprawnym
   * położeniem suwaka: zsunięty na sam początek pokazywałby wtedy cały bufor,
   * a opis pod nim mówiłby „wiersz 0 z N”. Rozróżnienie „nietknięty” od
   * „ustawiony na zero” musi więc istnieć w wartości, a nie w domyśle.
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
