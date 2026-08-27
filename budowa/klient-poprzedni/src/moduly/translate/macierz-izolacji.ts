import {
  Command,
  ConfigScope,
  IsolationLayer,
  type IsolationTechnicalSwitch,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk } from '../../modele/kontrolki-formularza';
import type { Kanal } from '../../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { utworzRozwiniecie } from './warstwy-translate';

/**
 * Macierz izolacji technicznej pokazuje stan ośmiu zakresów procesu sesji tłumaczenia,
 * wyłącznie do odczytu — zmianę ustawień prowadzi osobne okno konfiguracji punktów izolacji platformy.
 */
export interface MacierzIzolacji {
  element: HTMLElement;
  /** Odczytuje macierz z rdzenia; wolno wołać wielokrotnie. */
  odczytaj(): Promise<void>;
}

/** Nazwy ośmiu zakresów technicznych izolacji procesu sesji, wyświetlane przy odczycie stanu macierzy w oknie modułu tłumaczeń, w brzmieniu zrozumiałym dla operatora. */
const NAZWY_ZAKRESOW: Record<string, string> = {
  workingDirectory: 'Katalog roboczy',
  processEnvironment: 'Środowisko procesu',
  modelDataDirectory: 'Katalog danych i konfiguracji modelu',
  networkAccess: 'Dostęp sieciowy',
  fileAccess: 'Odczyt i zapis plików',
  accountToken: 'Konto i token sesji',
  processModel: 'Model procesu',
  executionServer: 'Serwer wykonania',
};

export function utworzMacierzIzolacji(kanal: Kanal): MacierzIzolacji {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Macierz izolacji technicznej',
    wyjasnienie:
      'Osiem zakresów technicznych procesu tej karty sesji — stan odczytany z rdzenia. ' +
      'Zmiana zakresu należy do okna punktów izolacji platformy.',
    znacznik: '☰',
  });

  const zdanie = document.createElement('p');
  zdanie.className = 'mt-izolacja__zdanie';
  zdanie.textContent = 'Macierzy jeszcze nie odczytano.';

  const wykaz = document.createElement('ul');
  wykaz.className = 'mt-izolacja';

  const ponow = przycisk('Odczytaj macierz', 'dn-btn dn-btn--sm dn-btn--zarys');
  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(ponow);

  rozwiniecie.tresc.append(zdanie, pasek, wykaz);

  async function odczytaj(): Promise<void> {
    const idSesji = kanal.sesja().id();
    if (idSesji === '') {
      wykaz.replaceChildren();
      zdanie.textContent =
        'Kanał nie ma ustalonej karty sesji, więc nie ma o czyją izolację zapytać. ' +
        'Macierz dotyczy procesu sesji, a nie modułu.';
      return;
    }

    zdanie.textContent = 'Odczyt ośmiu zakresów technicznych z rdzenia…';
    const wynik = sprawdzKsztalt(
      await wywolaj(kanal, Command.IsolationTechnicalGet, {
        scope: ConfigScope.Session,
        scopeId: idSesji,
        layer: IsolationLayer.Session,
      }),
      Command.IsolationTechnicalGet,
      (tresc) => czyTablica(tresc.switches),
    );

    if (!wynik.udany || wynik.wynik === undefined) {
      wykaz.replaceChildren();
      zdanie.textContent = opisOdmowyBledu('Odczyt macierzy izolacji technicznej', wynik.blad);
      return;
    }

    const przelaczniki = wynik.wynik.switches;
    const odciete = przelaczniki.filter((wpis) => wpis.isolated).length;
    zdanie.textContent =
      `Zakresów w macierzy: ${String(przelaczniki.length)}, odciętych: ${String(odciete)}. ` +
      'Zakres nieodcięty nie jest sprawdzany wcale — praca biegnie wtedy bez ograniczenia.';
    wykaz.replaceChildren(...przelaczniki.map(wierszZakresu));
  }

  ponow.addEventListener('click', () => void odczytaj());

  return { element: rozwiniecie.element, odczytaj };
}

/**
 * Jeden zakres macierzy.
 *
 * Stan nie jest niesiony samą barwą: obok wartości stoi słowo, a `data-odciety`
 * daje sprawdzianowi odczytać ją bez czytania zdania.
 */
function wierszZakresu(przelacznik: IsolationTechnicalSwitch): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mt-izolacja__nazwa';
  nazwa.textContent = NAZWY_ZAKRESOW[przelacznik.scope] ?? przelacznik.scope;

  const wartosc = document.createElement('span');
  wartosc.className = 'dn-plakietka mt-izolacja__wartosc';
  wartosc.textContent = przelacznik.isolated ? 'odcięty' : 'bez ograniczenia';

  const element = document.createElement('li');
  element.className = 'mt-izolacja__wiersz';
  element.dataset['odciety'] = przelacznik.isolated ? 'tak' : 'nie';
  element.append(nazwa, wartosc);

  const objasnienie = (przelacznik.explanation ?? '').trim();
  if (objasnienie !== '') {
    const opis = document.createElement('p');
    opis.className = 'mt-izolacja__opis';
    opis.textContent = objasnienie;
    element.append(opis);
  }
  return element;
}
