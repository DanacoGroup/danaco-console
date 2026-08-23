import { elementIkony } from '../ikony/ikony';
import type { PozycjaKomponentu } from './pozycje-komponentow';

/**
 * Kafel komponentu własnego — element strefy drugiej.
 *
 * Jedna odpowiedzialność: zbudowanie jednego kafla i zgłoszenie jego wyboru.
 *
 * Ta sama karta co środowiskowa, o mniejszej wadze: wariant
 * `.dn-karta--komponent` odbiera wstęgę i cień sygnału w spoczynku, zmniejsza
 * skalę ikony i zapisuje etykietę krojem bazowym półgrubym — nie nagłówkowym.
 *
 * Różnica krojów niesie znaczenie: nagłówkowy to wejście do środowiska, bazowy
 * to zbudowanie komponentu. Dlatego kafel nie dostaje klasy
 * `.dn-karta--akcent`, którą nosi karta środowiska.
 */

/** Ikona kafla w skali średniej zestawu — mniejszej niż godło środowiska. */
const ROZMIAR_IKONY = 18;

export interface KafelKomponentu {
  element: HTMLButtonElement;
}

export function utworzKafelKomponentu(
  pozycja: PozycjaKomponentu,
  przyWyborze: (pozycja: PozycjaKomponentu) => void,
): KafelKomponentu {
  const element = document.createElement('button');
  element.type = 'button';
  element.className =
    'dn-karta dn-karta--klikalna dn-karta--komponent dn-strona__kafel';
  element.dataset.komponent = pozycja.kod;
  // Kafli jednego rodzaju bywa wiele — rodzaj ich nie rozróżnia. Identyfikator
  // komponentu Operatora rozróżnia, i tylko kafel personalizowany go niesie.
  if (pozycja.komponent !== undefined) element.dataset.wystapienie = pozycja.komponent;

  const tresc = document.createElement('span');
  tresc.className = 'dn-karta-tresc dn-strona__kafel-tresc';

  const ikona = document.createElement('span');
  ikona.className = 'dn-strona__ikona-kafla';
  ikona.append(elementIkony(pozycja.ikona, { rozmiar: ROZMIAR_IKONY }));

  const napisy = document.createElement('span');
  napisy.className = 'dn-strona__kafel-napisy';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-karta-tytul dn-strona__kafel-nazwa';
  nazwa.textContent = pozycja.nazwa;

  const wezwanie = document.createElement('span');
  wezwanie.className = 'dn-strona__kafel-wezwanie';
  wezwanie.textContent = pozycja.wezwanie;

  napisy.append(nazwa, wezwanie, ...wierszMetadanych(pozycja));
  tresc.append(ikona, napisy);
  element.append(tresc);

  element.addEventListener('click', () => przyWyborze(pozycja));

  return { element };
}

/**
 * Wiersz metadanych kafla — stan czynności i czas ostatniej zmiany.
 *
 * Kafel rodzaju nie dostaje wiersza wcale: rodzaj nie jest bytem w bazie, więc
 * nie ma ani daty, ani stanu, a wiersz pusty wyglądałby jak metadane, których
 * nie odczytano.
 *
 * Czas formatuje widok, bo rdzeń nie zna strefy Operatora. Data pojawia się
 * dopiero, gdy zmiana wypadła innego dnia; w dniu bieżącym wystarcza godzina —
 * wzorem `opisWaznosci` z `uwierzytelnienie/sesja-bramki.ts`.
 */
function wierszMetadanych(pozycja: PozycjaKomponentu): HTMLElement[] {
  const metadane = pozycja.metadane;
  if (metadane === undefined) return [];

  const element = document.createElement('span');
  element.className = 'dn-strona__kafel-metadane';

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka dn-strona__kafel-stan';
  stan.dataset['czynny'] = String(metadane.czynny);
  stan.textContent = metadane.czynny ? 'czynny' : 'wyłączony';

  const zmiana = document.createElement('span');
  zmiana.className = 'dn-strona__kafel-zmiana';
  zmiana.textContent = opisZmiany(metadane.zmieniony);

  element.append(stan, zmiana);
  return [element];
}

/** Czas ostatniej zmiany zdaniem czytelnym; zero znaczy „rdzeń nie podał". */
function opisZmiany(zmieniony: number, teraz: number = Date.now()): string {
  if (!Number.isFinite(zmieniony) || zmieniony <= 0) return 'bez zapisanej zmiany';
  const kiedy = new Date(zmieniony);
  const dzis = new Date(teraz);
  const tenSamDzien =
    kiedy.getFullYear() === dzis.getFullYear() &&
    kiedy.getMonth() === dzis.getMonth() &&
    kiedy.getDate() === dzis.getDate();
  const godzina = kiedy.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
  if (tenSamDzien) return `zmieniony ${godzina}`;
  const data = kiedy.toLocaleDateString('pl-PL', { day: 'numeric', month: 'short' });
  return `zmieniony ${data}`;
}
