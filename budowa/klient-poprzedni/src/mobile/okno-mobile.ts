import './mobile.css';

import { Command } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import { utworzRameOkna } from '../komponenty/rama-okna';
import type { Kanal } from '../protokol/kanal';
import { wywolaj } from '../protokol/wywolanie';
import { utworzEkranInterwencji, type EkranInterwencji } from './ekran-interwencji';
import { utworzEkranProcesow, type EkranProcesow } from './ekran-procesow';
import type { MagazynKwitow } from './kwit-decyzji';

/**
 * Okno mobilnego centrum dowodzenia, otwierane z listwy Ustawienia. Prowadzi
 * trzy niezależne warstwy — kafelek stanu, ekran interwencji, przegląd
 * procesów — żadna nie zasłania pozostałych. Rama pochodzi z biblioteki,
 * okno otwiera się natychmiast.
 */
export interface OknoMobile {
  element: HTMLDialogElement;
  otworz(): void;
  zamknij(): void;
  rozlacz(): void;
  /** Ekran interwencji osadzony w oknie — do sprawdzianów widoku. */
  ekran(): EkranInterwencji;
  /** Ekran przeglądu zadań i procesów — do sprawdzianów widoku. */
  ekranProcesow(): EkranProcesow;
}

export function utworzOknoMobile(kanal: Kanal, magazyn?: MagazynKwitow): OknoMobile {
  const rama = utworzRameOkna({
    kod: 'mobile',
    tytul: 'Interwencja z telefonu',
    rola: 'monitor',
    modul: 'Mobile',
    przeznaczenie:
      'Cztery drogi interwencji w pracę, która biegnie bez ciebie: zatwierdź krok, wstrzymaj, ' +
      'nastaw koordynatora, przejmij sterowanie. Dwa dotknięcia od wykazu do decyzji.',
    przedrostek: 'mb',
  });

  const ekran = utworzEkranInterwencji({ kanal, ...(magazyn === undefined ? {} : { magazyn }) });
  // Trzecia warstwa okna: przegląd zadań i procesów ze sterowaniem, obok obrazu interwencji.
  const procesy = utworzEkranProcesow({ kanal });

  /** Kafelek stanu platformy — jedyne miejsce okna wołające rodzinę `mobile.*`. */
  const kafelek = document.createElement('p');
  kafelek.className = 'mb-platforma';
  kafelek.setAttribute('aria-live', 'polite');

  rama.cialo.append(kafelek, ekran.element, procesy.element);

  const odswiez = document.createElement('button');
  odswiez.type = 'button';
  odswiez.className = 'dn-btn dn-btn--zarys mb-cel';
  odswiez.textContent = 'Odczytaj ponownie';
  odswiez.addEventListener('click', () => void wczytaj());
  rama.narzedzia.append(odswiez);

  async function wczytajPlatforme(): Promise<void> {
    kafelek.dataset['stan'] = 'ladowanie';
    kafelek.textContent = 'Pytam rdzeń o stan platformy (mobile.status.get)…';

    const idSesji = kanal.sesja().id();
    const wynik = await wywolaj(kanal, Command.MobileStatusGet, {
      ...(idSesji === '' ? {} : { deviceId: idSesji }),
    });

    if (!wynik.udany || wynik.wynik === undefined) {
      kafelek.dataset['stan'] = 'blad';
      kafelek.textContent =
        `${opisOdmowyBledu('Stan platformy mobilnej', wynik.blad)}. ` +
        'Wykaz decyzji poniżej tego nie dotyczy — stoi na monitor.status, queue.list ' +
        'i window.state.get, które rdzeń obsługuje.';
      return;
    }

    const stan = wynik.wynik.status;
    kafelek.dataset['stan'] = 'tresc';
    kafelek.textContent =
      `Warstwa mobilna wpięta: urządzenie ${stan.deviceId ?? 'nienazwane'}, ` +
      `sparowane: ${stan.paired ? 'tak' : 'nie'}; sesji ${stan.sessionCount}, ` +
      `okien ${stan.windowCount}, procesów w biegu ${stan.runningProcessCount} (mobile.status.get).`;
  }

  async function wczytaj(): Promise<void> {
    // Trzy odczyty równolegle i niezależnie: odmowa jednego nie zabiera pozostałych.
    await Promise.all([wczytajPlatforme(), ekran.odswiez(), procesy.odswiez()]);
  }

  const element = document.createElement('dialog');
  element.className = 'dn-modal mb-okno';
  element.setAttribute('aria-label', 'Okno interwencji z telefonu');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo mb-okno__cialo';
  cialo.append(rama.element);

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--atrament mb-cel';
  zamknij.textContent = 'Zamknij';
  zamknij.addEventListener('click', () => element.close());

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka mb-okno__stopka';
  stopka.append(zamknij);

  element.append(cialo, stopka);

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      void wczytaj();
    },

    zamknij: () => element.close(),

    rozlacz() {
      ekran.rozlacz();
      element.remove();
    },

    ekran: () => ekran,

    ekranProcesow: () => procesy,
  };
}
