/**
 * Pasek transportu jednej kolejki: sześć przycisków sterowania i przełożenie
 * naciśnięcia na zamiar. Nazwy czterech pierwszych pochodzą z wyliczenia
 * `QueueAction` kontraktu, a przekazanie — z komendy `context.transfer`.
 */

import { Command, QueueAction } from '../../../shared/contract';
import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { ETYKIETA_DZIALANIA_KOLEJKI } from './etykiety-pulpitu';
import type { KolejkaPulpitu } from './model-danych';
import { DZIALANIA_TRANSPORTU, type ZamiarKolejki } from './zdarzenia-pulpitu';

/**
 * Ikona przypisana działaniu transportu. Mapa obejmuje wszystkie pozycje
 * wyliczenia `QueueAction`, także te, których pasek pulpitu nie pokazuje, bo mapa
 * niepełna nie skompilowałaby się przy pierwszym ich użyciu.
 */
const IKONA_DZIALANIA: Readonly<Record<QueueAction, NazwaIkony>> = {
  [QueueAction.Start]: 'uruchom',
  [QueueAction.Pause]: 'zegar',
  [QueueAction.Stop]: 'zatrzymaj',
  [QueueAction.Retry]: 'odswiez',
  [QueueAction.Resume]: 'uruchom',
  [QueueAction.Clear]: 'kosz',
  [QueueAction.Enqueue]: 'plus',
  [QueueAction.Dequeue]: 'zamknij',
  [QueueAction.Delay]: 'zegar',
  [QueueAction.Split]: 'galaz',
  [QueueAction.Merge]: 'wezly',
  [QueueAction.Route]: 'strzalka-prawo',
  [QueueAction.Branch]: 'galaz',
  [QueueAction.Condition]: 'filtr',
};

/**
 * Buduje pasek sześciu przycisków transportu dla jednej kolejki. Naciśnięcie
 * oddaje zamiar warstwie nadrzędnej — pasek sam nie wysyła komendy, więc pulpit
 * pozostaje jedynym miejscem, które zna kanał.
 */
export function utworzPasekTransportu(
  kolejka: KolejkaPulpitu,
  nadaj: (zamiar: ZamiarKolejki) => void,
): HTMLElement {
  const pasek = document.createElement('div');
  pasek.className = 'mc-transport';
  pasek.setAttribute('role', 'group');
  pasek.setAttribute('aria-label', `Sterowanie kolejki ${kolejka.nazwa}`);

  for (const dzialanie of DZIALANIA_TRANSPORTU) {
    pasek.append(
      przycisk(ETYKIETA_DZIALANIA_KOLEJKI[dzialanie], IKONA_DZIALANIA[dzialanie], () =>
        nadaj({
          rodzaj: 'kolejka',
          komenda: Command.QueueAction,
          dzialanie,
          idKolejki: kolejka.id,
          rola: kolejka.nazwa,
        }),
      ),
    );
  }

  pasek.append(
    przycisk('Podnieś priorytet', 'gwiazdka', () =>
      nadaj({ rodzaj: 'priorytet', komenda: null, idKolejki: kolejka.id, rola: kolejka.nazwa }),
    ),
    przycisk('Przekaż', 'strzalka-prawo', () =>
      nadaj({
        rodzaj: 'przekazanie',
        komenda: Command.ContextTransfer,
        idKolejki: kolejka.id,
        rola: kolejka.nazwa,
      }),
    ),
  );

  return pasek;
}

/**
 * Jeden przycisk transportu: ikona z nazwą czytaną i podpowiedzią.
 * Nazwa jest podana dwiema drogami — `aria-label` i `title` — bo sam kształt
 * ikony nie wystarcza za opis działania.
 */
function przycisk(nazwa: string, ikona: NazwaIkony, naNacisniecie: () => void): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn-ikona mc-transport__przycisk';
  element.title = nazwa;
  element.setAttribute('aria-label', nazwa);
  element.append(elementIkony(ikona, { rozmiar: 16 }));
  element.addEventListener('click', naNacisniecie);
  return element;
}
