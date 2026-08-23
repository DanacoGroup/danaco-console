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
 * Okno mobilnego centrum dowodzenia (`mobile`, monitor) — otwierane z listwy
 * Ustawienia obok Konfiguracji, Punktów Izolacji, Modeli i AOD.
 *
 * Okno prowadzi trzy niezależne warstwy i żadna nie zasłania pozostałych:
 *
 *  1. Kafelek stanu platformy woła `mobile.status.get` i pokazuje odpowiedź
 *     rdzenia bez interpretacji. Odmowa dotyczy wyłącznie kafelka i nie gasi
 *     niczego poza nim.
 *  2. Ekran interwencji (`ekran-interwencji.ts`) stoi na komendach, które rdzeń
 *     obsługuje — `monitor.status`, `queue.list`, `window.list`,
 *     `window.state.get`, `role.list` — i niesie cztery drogi interwencji
 *     w dwóch dotknięciach.
 *  3. Przegląd zadań i procesów (`ekran-procesow.ts`) woła
 *     `mobile.process.list` i `mobile.process.control` — wykaz procesów wraz ze
 *     sterowaniem nimi, czyli jedyną czynność sprawczą okna nad procesem.
 *     Czynność nieodwracalna mówi to przed wykonaniem.
 *
 * Warstwy nie są trzema źródłami prawdy: wszystkie czytają ten sam rejestr
 * telemetrii procesów, z którego czyta `monitor.status`. Różnią się drogą i tym,
 * co potrafią — obraz interwencji rozstrzyga, wykaz procesów steruje.
 *
 * Rama pochodzi wyłącznie z biblioteki (`komponenty/rama-okna.ts`); okno otwiera
 * się natychmiast, a odczyt dojeżdża odpowiedzią rdzenia.
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
  // Trzecia warstwa okna: przegląd zadań i procesów wraz ze sterowaniem
  // (rozdz. 4.6 i 5.3 opracowania Mobile). Stoi obok obrazu interwencji, a nie
  // w nim, bo czyta inną rodziną komend i niesie to, czego obraz nie ma —
  // sterowanie procesem. Odmowa jednej warstwy nie gasi pozostałych dwóch.
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
    // Trzy odczyty równolegle i niezależnie: odmowa jednego nie zabiera
    // pozostałych. Kafelek stanu platformy, obraz interwencji i wykaz procesów
    // mówią o tym samym rdzeniu trzema różnymi rodzinami komend.
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
