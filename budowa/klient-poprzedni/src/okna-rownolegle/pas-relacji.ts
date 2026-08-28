import './relacje.css';

import { WindowRole } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { ETYKIETA_RELACJI, OBJASNIENIE_RELACJI, opisGniazdaZRola } from './etykiety-ukladu';
import { numerGniazda, type IdGniazda } from './identyfikatory';
import { utworzObjasnienie } from './objasnienie-ukladu';
import { utworzPlakietkeStanu } from './plakietka-stanu';
import { utworzPustkeRelacji } from './pustka-relacji';
import type { StanPary } from './stan-pary';
import { utworzUwagePodgladu } from './uwaga-podgladu';
import { krance, type Wiez } from './wiez-koordynacji';

/** Pas relacji pod torem okien — powiązanie pary koordynator-wykonawca widoczne jednym spojrzeniem, bez sięgania po stan pary osobno. */
export interface PasRelacji {
  element: HTMLElement;
  /** Rysuje więź nad kolumnami obu okien pary; `null` daje stan pusty. */
  ustawWiez(wiez: Wiez | null, liczba: number): void;
  /** Ustawia stan pętli pokazywany na pasie. */
  ustawStan(stan: StanPary): void;
  /** Puszcza żeton zlecenia wzdłuż szyny — chwila przekazania. */
  puscZeton(): void;
  /** Wprowadza trwałą uwagę: przekazanie jest podglądem układu, nie pracą. */
  oznaczPodglad(): void;
}

/** Czas przebiegu żetonu po szynie relacji, w milisekundach, licząc od jednego krańca więzi do drugiego. */
const CZAS_ZETONU = 900;

/**
 * Pas relacji układu okien równoległych rysuje powiązanie między oknami pod nimi, na siatce tych samych kolumn co tor okien, i nie znika przy braku pary — wtedy pokazuje stan pusty.
 */
export function utworzPasRelacji(): PasRelacji {
  const element = document.createElement('section');
  element.className = 'dn-okna__relacje';
  element.setAttribute('aria-label', ETYKIETA_RELACJI);

  const etykieta = document.createElement('span');
  etykieta.className = 'dn-okna__relacje-etykieta';
  etykieta.textContent = ETYKIETA_RELACJI;

  const naglowek = document.createElement('div');
  naglowek.className = 'dn-okna__relacje-naglowek';
  naglowek.append(etykieta, utworzObjasnienie(OBJASNIENIE_RELACJI));

  const pustka = utworzPustkeRelacji();
  const uwaga = utworzUwagePodgladu();

  const tor = document.createElement('div');
  tor.className = 'dn-okna__tor-relacji';

  const { wiezElement, lewy, prawy, grot, zeton } = zlozWiez();
  const stan = utworzPlakietkeStanu('gotowa');

  wiezElement.append(stan.element);
  tor.append(wiezElement);
  element.append(naglowek, tor, pustka, uwaga.element);

  let zegar: number | undefined;

  return {
    element,

    oznaczPodglad: () => uwaga.pokaz(),

    ustawWiez(wiez, liczba) {
      // Stan pusty wchodzi w miejsce toru więzi; nagłówek, dymek i uwaga zostają na swoim miejscu.
      pustka.hidden = wiez !== null;
      tor.hidden = wiez === null;
      if (wiez === null) return;
      tor.dataset.liczba = String(liczba);

      const konce = krance(wiez);
      // Pozycja w siatce jest danymi układu, nie wartością wizualną — jedyny styl powstający w kodzie.
      wiezElement.style.gridColumn = `${numerGniazda(konce.lewy)} / ${numerGniazda(konce.prawy) + 1}`;
      wiezElement.dataset.kierunek = wiez.wPrawo ? 'w-prawo' : 'w-lewo';

      lewy.textContent = opisGniazdaZRola(konce.lewy, rolaKonca(wiez, konce.lewy));
      prawy.textContent = opisGniazdaZRola(konce.prawy, rolaKonca(wiez, konce.prawy));
      grot.replaceChildren(
        elementIkony(wiez.wPrawo ? 'strzalka-prawo' : 'strzalka-lewo', {
          rozmiar: 18,
          etykieta: `Kierunek zlecenia: ${lewy.textContent} → ${prawy.textContent}`,
        }),
      );
    },

    ustawStan: (nowy) => stan.ustaw(nowy),

    puscZeton() {
      zeton.hidden = false;
      zeton.classList.remove('dn-okna__zeton--biegnie');
      void zeton.offsetWidth;
      zeton.classList.add('dn-okna__zeton--biegnie');
      window.clearTimeout(zegar);
      zegar = window.setTimeout(() => {
        zeton.classList.remove('dn-okna__zeton--biegnie');
        zeton.hidden = true;
      }, CZAS_ZETONU);
    },
  };
}

/** Rola krańca więzi wynika wprost z samej więzi i nigdy nie wymaga sięgania po stan pary z układu okien. */
function rolaKonca(wiez: Wiez, id: IdGniazda): WindowRole {
  return id === wiez.koordynator ? WindowRole.Coordinator : WindowRole.Executor;
}

/** Węzły więzi: dwa krańce, szyna z żetonem i grot — sam kształt struktury danych, bez żadnej treści widoku. */
function zlozWiez(): {
  wiezElement: HTMLElement;
  lewy: HTMLElement;
  prawy: HTMLElement;
  grot: HTMLElement;
  zeton: HTMLElement;
} {
  const wiezElement = document.createElement('div');
  wiezElement.className = 'dn-okna__wiez';

  const lewy = document.createElement('span');
  lewy.className = 'dn-okna__wiez-koniec';

  const prawy = document.createElement('span');
  prawy.className = 'dn-okna__wiez-koniec';

  const szyna = document.createElement('span');
  szyna.className = 'dn-okna__wiez-szyna';

  const zeton = document.createElement('span');
  zeton.className = 'dn-okna__zeton';
  zeton.hidden = true;

  const grot = document.createElement('span');
  grot.className = 'dn-okna__wiez-grot';

  szyna.append(zeton, grot);
  wiezElement.append(lewy, szyna, prawy);
  return { wiezElement, lewy, prawy, grot, zeton };
}
