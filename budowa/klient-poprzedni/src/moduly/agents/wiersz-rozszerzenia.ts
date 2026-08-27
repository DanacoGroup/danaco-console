/**
 * Jeden wiersz katalogu rozszerzeń wraz z akcjami, które da się na nim wykonać.
 * Wiersz odpowiada wyłącznie za postać pozycji katalogu: nie woła rdzenia, tylko
 * zgłasza czynność oknu, a okno prowadzi wywołanie i pokazuje odpowiedź.
 */

import { ExtensionKind, type Extension } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';

/**
 * Czynność zgłaszana oknu z wiersza katalogu. Wiersz nazywa zamiar, a wybór
 * komendy rodziny `extension.*` i obsługa odpowiedzi należą do okna, które
 * jako jedyne trzyma kanał do rdzenia.
 */
export type CzynnoscRozszerzenia = 'instaluj' | 'przelacz' | 'odinstaluj';

export interface ObslugaWiersza {
  (czynnosc: CzynnoscRozszerzenia, pozycja: Extension): void;
}

/**
 * Nazwy rodzajów rozszerzeń po polsku, bo kontrakt niesie kody angielskie.
 * Wykaz obejmuje wszystkie pozycje wyliczenia `ExtensionKind`, więc dołożenie
 * rodzaju w kontrakcie wymusza uzupełnienie nazwy przy kompilacji.
 */
const NAZWA_RODZAJU: Record<ExtensionKind, string> = {
  [ExtensionKind.Mcp]: 'serwer MCP',
  [ExtensionKind.Plugin]: 'wtyczka',
  [ExtensionKind.Api]: 'integracja API',
  [ExtensionKind.Skill]: 'umiejętność',
};

/**
 * Czytelna nazwa rodzaju; kod nieznany zostaje kodem, zamiast zniknąć z wiersza.
 * Rdzeń w nowszej wersji może oddać rodzaj, którego ten klient jeszcze nie zna,
 * a pozycja bez nazwy rodzaju wyglądałaby na uszkodzoną.
 */
export function nazwaRodzaju(rodzaj: ExtensionKind): string {
  return NAZWA_RODZAJU[rodzaj] ?? rodzaj;
}

export function utworzWierszRozszerzenia(
  pozycja: Extension,
  obsluz: ObslugaWiersza,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'da-rozszerzenie';
  element.dataset['rozszerzenie'] = pozycja.code;
  element.dataset['rodzaj'] = pozycja.kind;
  element.dataset['zainstalowane'] = String(pozycja.installed);

  const tytul = document.createElement('span');
  tytul.className = 'da-rozszerzenie__nazwa';
  tytul.textContent = pozycja.name;

  const rodzaj = document.createElement('span');
  rodzaj.className = 'dn-plakietka da-rozszerzenie__rodzaj';
  rodzaj.textContent = nazwaRodzaju(pozycja.kind);

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka da-rozszerzenie__stan';
  stan.textContent = opisStanu(pozycja);

  const naglowek = document.createElement('div');
  naglowek.className = 'da-rozszerzenie__naglowek';
  naglowek.append(tytul, rodzaj, stan);
  element.append(naglowek);

  const opis = (pozycja.description ?? '').trim();
  if (opis !== '') {
    const akapit = document.createElement('p');
    akapit.className = 'dn-pole-opis da-rozszerzenie__opis';
    akapit.textContent = opis;
    element.append(akapit);
  }

  element.append(pasAkcji(pozycja, obsluz));
  return element;
}

/**
 * Stan pozycji jednym zdaniem — trzy postaci, nie dwie: pozycja włączona to co
 * innego niż zainstalowana. Rozszerzenie zainstalowane i wyłączone zostaje
 * w katalogu, więc wiersz musi je odróżnić od nieobecnego.
 */
function opisStanu(pozycja: Extension): string {
  if (!pozycja.installed) return 'niezainstalowane';
  return pozycja.enabled ? 'włączone' : 'wyłączone';
}

function pasAkcji(pozycja: Extension, obsluz: ObslugaWiersza): HTMLElement {
  const pas = document.createElement('div');
  pas.className = 'da-rozszerzenie__akcje';

  if (!pozycja.installed) {
    const instaluj = przycisk('Zainstaluj', 'dn-btn dn-btn--sm dn-btn--atrament');
    instaluj.addEventListener('click', () => obsluz('instaluj', pozycja));
    pas.append(instaluj);
    return pas;
  }

  const przelacz = przycisk(
    pozycja.enabled ? 'Wyłącz' : 'Włącz',
    'dn-btn dn-btn--sm',
  );
  przelacz.addEventListener('click', () => obsluz('przelacz', pozycja));

  const odinstaluj = przycisk('Odinstaluj', 'dn-btn dn-btn--sm dn-btn--niebezpieczny');
  odinstaluj.addEventListener('click', () => obsluz('odinstaluj', pozycja));

  pas.append(przelacz, odinstaluj);
  return pas;
}
