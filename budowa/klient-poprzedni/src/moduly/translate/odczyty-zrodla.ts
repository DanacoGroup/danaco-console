import { liczSlowa } from './karta-lqa';
import { oznaczWarstwe } from './warstwy-translate';
import { znajdzWzorce, type Wystapienie } from './wzorce-placeholderow';

/**
 * Dwa odczyty tekstu źródłowego liczone w oknie: objętość materiału i ochrona symboli
 * zastępczych, liczone przez okno z tekstu, a nie przez rdzeń.
 */
export interface OdczytyZrodla {
  /** Element osadzany w Source Panel. */
  element: HTMLElement;
  /** Przelicza odczyty z tekstu w polu i liczby segmentów oddanej przez rdzeń. */
  odswiez(tekst: string, segmentowZRdzenia: number): void;
}

export function utworzOdczytyZrodla(): OdczytyZrodla {
  const objetosc = document.createElement('p');
  objetosc.className = 'mt-odczyty__objetosc';

  const ochrona = document.createElement('p');
  ochrona.className = 'mt-odczyty__ochrona';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-odczyty__wzorce';

  const element = document.createElement('div');
  element.className = 'mt-odczyty';
  oznaczWarstwe(element, 1);
  element.append(objetosc, ochrona, wykaz);

  return {
    element,

    odswiez(tekst, segmentowZRdzenia) {
      objetosc.textContent = zdanieObjetosci(tekst, segmentowZRdzenia);
      const wzorce = znajdzWzorce(tekst);
      ochrona.textContent = zdanieOchrony(wzorce);
      wykaz.replaceChildren(...wzorce.map(wierszWzorca));
      element.dataset['wzorce'] = String(wzorce.length);
    },
  };
}

function zdanieObjetosci(tekst: string, segmentowZRdzenia: number): string {
  const slow = liczSlowa(tekst);
  const znakow = tekst.length;
  const oSegmentach =
    segmentowZRdzenia === 0
      ? 'segmentów rdzeń jeszcze nie oddał'
      : `segmentów ${String(segmentowZRdzenia)} wedle rdzenia`;
  return `Objętość materiału: słów ${String(slow)}, znaków ${String(znakow)}, ${oSegmentach}. ` +
    'Słowa i znaki liczy okno z pola wyżej; podział na segmenty należy do rdzenia.';
}

function zdanieOchrony(wzorce: readonly Wystapienie[]): string {
  if (wzorce.length === 0) {
    return (
      'Ochrona symboli zastępczych: w tekście nie ma znanych oknu wzorców. Wykaz obejmuje ' +
      'zapisy wymienione w opracowaniu modułu, więc brak wystąpień nie znaczy, że materiał ' +
      'nie ma żadnych zmiennych.'
    );
  }
  return (
    `Ochrona symboli zastępczych: wystąpień ${String(wzorce.length)}. Elementy z wykazu niżej ` +
    'mają wrócić w przekładzie nietknięte; niezgodność ich liczby albo kolejności rdzeń zgłasza ' +
    'przy kontroli jakości panelu.'
  );
}

function wierszWzorca(wzorzec: Wystapienie): HTMLElement {
  const zapis = document.createElement('span');
  zapis.className = 'mt-odczyty__zapis';
  zapis.textContent = wzorzec.zapis;

  const styl = document.createElement('span');
  styl.className = 'mt-odczyty__styl';
  styl.textContent =
    wzorzec.nazwa === '' ? wzorzec.styl : `${wzorzec.styl} · nazwa ${wzorzec.nazwa}`;

  const element = document.createElement('li');
  element.className = 'mt-odczyty__wzorzec';
  element.append(zapis, styl);
  return element;
}
