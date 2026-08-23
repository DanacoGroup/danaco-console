import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Obrót i powiększenie podglądu — czynności widoku, nie rdzenia.
 *
 * Dotyczą sposobu pokazania treści już pobranej, więc odbywają się w całości
 * po stronie klienta, bez komendy. Stopień obrotu i krotność powiększenia idą
 * do `dataset`, a wygląd bierze się z arkusza modułu — plik widoku nie zna
 * żadnej wartości wizualnej.
 */
export interface SterowanieWidoku {
  element: HTMLElement;
  /** Zdejmuje obrót i powiększenie — po zmianie pliku widok wraca do stanu bazowego. */
  zeruj(): void;
}

/** Dopuszczalne krotności powiększenia; kolejne naciśnięcie bierze następną. */
const KROTNOSCI = ['1', '2', '3'] as const;

export function utworzSterowanieWidoku(cel: HTMLElement): SterowanieWidoku {
  let obrot = 0;
  let krotnosc = 0;

  function zastosuj(): void {
    cel.dataset['obrot'] = String(obrot);
    cel.dataset['powiekszenie'] = KROTNOSCI[krotnosc] ?? '1';
  }

  const obroc = przycisk('Obróć o 90°', 'dn-btn dn-btn--sm dn-btn--duch');
  obroc.dataset['czynnosc'] = 'obrot';
  obroc.addEventListener('click', () => {
    obrot = (obrot + 90) % 360;
    zastosuj();
  });

  const powieksz = przycisk('Powiększ', 'dn-btn dn-btn--sm dn-btn--duch');
  powieksz.dataset['czynnosc'] = 'powiekszenie';
  powieksz.addEventListener('click', () => {
    krotnosc = (krotnosc + 1) % KROTNOSCI.length;
    zastosuj();
  });

  const element = document.createElement('span');
  element.className = 'ml-sterowanie';
  element.append(obroc, powieksz);
  zastosuj();

  return {
    element,
    zeruj() {
      obrot = 0;
      krotnosc = 0;
      zastosuj();
    },
  };
}
