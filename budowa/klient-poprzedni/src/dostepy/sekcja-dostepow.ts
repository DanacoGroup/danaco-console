import './dostepy.css';

import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { utworzListeNadan } from './lista-nadan';
import { utworzObszarKataloguRoboczego } from './obszar-katalogu-roboczego';
import { utworzStanDostepow } from './stan-dostepow';
import { utworzStanyOdczytu } from './stany-odczytu';
import { utworzWykazPunktow } from './wykaz-punktow';

/**
 * Sekcja dostępów i katalogu roboczego, czyli ekran, na którym Operator nadaje
 * dostęp. Dotyczy punktu dostępu wraz z nadaniem oraz katalogu roboczego,
 * w którym model zostawia swoje pliki.
 */
export interface SekcjaDostepow {
  /** Element osadzany w widoku gospodarza. */
  element: HTMLElement;
  /** Wiąże sekcję z oknem rozmowy, którego nadania ma pokazywać. */
  ustawOkno(oknoID: string): void;
  /** Wczytuje wykaz punktów, nadania okna i ustawienia katalogu roboczego. */
  wczytaj(): void;
  /** Odłącza subskrypcje kanału i nasłuch powłoki. */
  rozlacz(): void;
}

export interface ZaleznosciSekcji {
  /** Kanał komunikatów rdzenia. */
  kanal: Kanal;
  /** Okno rozmowy, z którym sekcja startuje; puste znaczy brak wiązania. */
  oknoID?: string;
}

export function utworzSekcjeDostepow(zaleznosci: ZaleznosciSekcji): SekcjaDostepow {
  const { kanal } = zaleznosci;
  const stan = utworzStanDostepow(kanal, zaleznosci.oknoID ?? '');

  const wykaz = utworzWykazPunktow(stan);
  const nadania = utworzListeNadan(stan);
  const katalog = utworzObszarKataloguRoboczego(kanal);
  const stany = utworzStanyOdczytu(stan, () => wczytaj());

  const odswiez = document.createElement('button');
  odswiez.type = 'button';
  odswiez.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  odswiez.append(
    elementIkony('odswiez', { rozmiar: 16 }),
    document.createTextNode('Odczytaj ponownie'),
  );
  odswiez.addEventListener('click', () => wczytaj());

  const kolumny = document.createElement('div');
  kolumny.className = 'dd-sekcja__kolumny';
  kolumny.append(wykaz.element, nadania.element);

  const element = document.createElement('section');
  element.className = 'dd-sekcja';
  element.append(naglowek(odswiez), stany.element, kolumny, katalog.element);

  const przelicz = (): void => {
    stany.odswiez();
    wykaz.odswiez();
    nadania.odswiez();
  };

  stan.naZmiane(przelicz);
  przelicz();

  function wczytaj(): void {
    void stan.odswiez().then(przelicz);
    katalog.wczytaj();
  }

  return {
    element,

    ustawOkno: (oknoID) => stan.ustawOkno(oknoID),

    wczytaj,

    rozlacz() {
      stan.rozlacz();
      wykaz.rozlacz();
      katalog.rozlacz();
    },
  };
}

/**
 * Nagłówek sekcji: tytuł, zdanie rozdzielające trzy byty, których sekcja
 * dotyczy, oraz przycisk odczytu ponownego. Zdanie stoi w nagłówku, ponieważ
 * rozróżnienie bytów jest warunkiem poprawnego nadania dostępu.
 */
function naglowek(odswiez: HTMLElement): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dd-sekcja__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dd-sekcja__tytul';
  tytul.textContent = 'Dostępy i katalog roboczy';

  const opis = document.createElement('p');
  opis.className = 'dd-sekcja__opis';
  opis.textContent =
    'Nadanie dostępu wiąże okno rozmowy z punktem: maszyną za mostem MCP albo katalogiem na urządzeniu. Środowiska — profilu widoczności modułów w bocznej nawigacji — ta sekcja nie dotyczy.';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'dd-sekcja__rozpychacz';

  element.append(elementIkony('klodka', { rozmiar: 20 }), tytul, rozpychacz, odswiez);
  element.append(opis);
  return element;
}
