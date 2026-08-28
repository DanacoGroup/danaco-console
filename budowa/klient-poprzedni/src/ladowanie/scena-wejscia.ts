import './ladowanie.css';

import { elementGodla } from '../ikony/ikony';

/**
 * Scena wejścia: obraz między przejściem bramki a postawieniem strony głównej.
 * Próg mignięcia poniżej jest najkrótszym czasem pozostawania sceny na ekranie.
 */
const PROG_MIGNIECIA = 320;

/**
 * Kres czekania na gotowość, w milisekundach; po nim scena schodzi mimo wszystko,
 * ponieważ scena, która nie schodzi, jest zamknięciem produktu.
 */
const KRES_CZEKANIA = 8000;

export interface OpisSceny {
  /**
   * Gotowość, na którą scena czeka; odrzucenie znaczy to samo co spełnienie.
   */
  gotowosc: Promise<unknown>;
  /** Zdanie pod znakiem — czynność, która trwa. */
  krok?: string;
}

export interface ScenaWejscia {
  element: HTMLElement;
  /** Stawia scenę nad dokumentem. */
  uruchom(): void;
  /** Zdejmuje scenę bez czekania — dla ścieżek, które sceny nie potrzebują. */
  zdejmij(): void;
}

/**
 * Buduje scenę wejścia wraz z jej zawartością, nie montując jej w dokumencie;
 * montaż należy do osobnej czynności uruchomienia sceny.
 */
export function utworzScenaWejscia(opis: OpisSceny): ScenaWejscia {
  const bryla = document.createElement('div');
  bryla.className = 'la-bryla';
  // Znak marki na każdej ścianie, bo bryła obraca się w kółko.
  for (const strona of ['przod', 'tyl', 'prawo', 'lewo', 'gora', 'dol']) {
    const sciana = document.createElement('div');
    sciana.className = `la-sciana la-sciana--${strona}`;
    // Godło bez etykiety, bo zdanie o marce pada raz w napisie pod bryłą.
    sciana.append(elementGodla({ rozmiar: 48, podloze: 'jasne' }));
    bryla.append(sciana);
  }

  const pole = document.createElement('div');
  pole.className = 'la-pole';
  pole.append(bryla);

  const nazwa = document.createElement('p');
  nazwa.className = 'la-nazwa';
  nazwa.textContent = 'Danaco Console';

  const krok = document.createElement('p');
  krok.className = 'la-krok';
  krok.textContent = opis.krok ?? 'Przygotowuję Centrum dowodzenia…';

  const napisy = document.createElement('div');
  napisy.className = 'la-napisy';
  napisy.append(nazwa, krok);

  const element = document.createElement('div');
  element.className = 'la-scena';
  // Scena mówi o sobie czytnikowi ekranu jednym zdaniem, na żywo.
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');
  element.append(pole, napisy);

  let zdjeta = false;

  /** Zdejmuje scenę: wygaszenie, potem usunięcie z dokumentu. */
  function zdejmij(): void {
    if (zdjeta) return;
    zdjeta = true;
    element.classList.add('la-scena--schodzi');
    // Usunięcie idzie po wygaszeniu; samo zdarzenie przejścia nie wystarcza.
    const usun = (): void => element.remove();
    element.addEventListener('transitionend', usun, { once: true });
    window.setTimeout(usun, 400);
  }

  return {
    element,

    uruchom() {
      if (zdjeta || element.isConnected) return;
      document.body.append(element);

      const poczatek = Date.now();
      const zejdz = (): void => {
        // Próg mignięcia: scena albo stoi chwilę, albo nie staje wcale.
        const zostalo = PROG_MIGNIECIA - (Date.now() - poczatek);
        if (zostalo > 0) window.setTimeout(zdejmij, zostalo);
        else zdejmij();
      };

      void Promise.race([
        Promise.resolve(opis.gotowosc).catch(() => undefined),
        new Promise((spelnij) => window.setTimeout(spelnij, KRES_CZEKANIA)),
      ]).then(zejdz);
    },

    zdejmij,
  };
}
