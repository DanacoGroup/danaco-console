import { ExtensionOrigin, type Extension } from '../../../../shared/contract';

/**
 * Karta pozycji katalogu rozszerzeń w siatce kart App Catalogu. Pokazuje
 * wyłącznie pola oddane przez rdzeń: nazwę, rodzaj, wersję, pochodzenie, opis
 * oraz stany zainstalowania i włączenia. Kanału nie zna i niczego nie wysyła.
 */
export interface CzynnosciKarty {
  /** Instalacja pozycji jeszcze niezainstalowanej. */
  zainstaluj(pozycja: Extension): void;
  /** Włączenie albo wyłączenie pozycji zainstalowanej. */
  przelacz(pozycja: Extension): void;
  /** Wskazanie pozycji panelom bocznym — uprawnieniom i konsoli. */
  wskaz(pozycja: Extension): void;
}

/**
 * Napisy pochodzenia rozszerzenia: wartościom wyliczenia ExtensionOrigin
 * z kontraktu przypisuje brzmienia pokazywane w interfejsie katalogu.
 */
const NAZWY_POCHODZENIA: Readonly<Record<string, string>> = {
  [ExtensionOrigin.Danaco]: 'Danaco Plugin',
  [ExtensionOrigin.Personal]: 'Personal',
};

/**
 * Oddaje nazwę pochodzenia rozszerzenia w brzmieniu interfejsu, a gdy rdzeń
 * oddał wartość spoza wyliczenia ExtensionOrigin — samą tę wartość bez zmiany.
 */
export function nazwaPochodzenia(pochodzenie: string): string {
  return NAZWY_POCHODZENIA[pochodzenie] ?? pochodzenie;
}

/**
 * Zdanie o stanie pozycji podawane słowem, nie barwą. Rozróżnia trzy stany, bo
 * tyle niosą pola installed oraz enabled kontraktu; stanu z dostępną
 * aktualizacją nie podaje, bo pozycja nie niesie wersji rejestru.
 */
export function opisStanu(pozycja: Extension): string {
  if (!pozycja.installed) return 'dostępna';
  return pozycja.enabled ? 'zainstalowana, włączona' : 'zainstalowana, wyłączona';
}

/**
 * Odmiana plakietki stanu pozycji dobierana z pól installed oraz enabled;
 * nazwy klas pochodzą z biblioteki komponentów interfejsu.
 */
function odmianaStanu(pozycja: Extension): string {
  if (!pozycja.installed) return 'dn-plakietka';
  return pozycja.enabled ? 'dn-plakietka dn-plakietka--sukces' : 'dn-plakietka dn-plakietka--ostrzezenie';
}

export function utworzKarteRozszerzenia(
  pozycja: Extension,
  czynnosci: CzynnosciKarty,
): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-karta__nazwa';
  nazwa.textContent = pozycja.name;

  const rodzaj = document.createElement('span');
  rodzaj.className = 'dn-plakietka dn-plakietka--rola';
  rodzaj.textContent = pozycja.kind;

  const stan = document.createElement('span');
  stan.className = odmianaStanu(pozycja);
  stan.textContent = opisStanu(pozycja);

  const glowa = document.createElement('div');
  glowa.className = 'mp-karta__glowa';
  glowa.append(nazwa, rodzaj, stan);

  const metryka = document.createElement('p');
  metryka.className = 'mp-karta__metryka';
  metryka.textContent =
    `${nazwaPochodzenia(pozycja.origin)} · wersja ` +
    `${pozycja.version ?? 'nieoddana przez rdzeń'} · kod ${pozycja.code}`;

  const element = document.createElement('li');
  element.className = 'mp-karta';
  element.dataset['rozszerzenie'] = pozycja.id;
  element.dataset['zainstalowane'] = String(pozycja.installed);
  element.dataset['wlaczone'] = String(pozycja.enabled);
  element.append(glowa, metryka);

  if (pozycja.description !== undefined && pozycja.description !== '') {
    const opis = document.createElement('p');
    opis.className = 'mp-karta__opis';
    opis.textContent = pozycja.description;
    element.append(opis);
  }

  const pasek = document.createElement('div');
  pasek.className = 'mp-karta__pasek';

  // Instalacja i przełączenie to dwie różne komendy, więc karta daje dwa osobne przyciski.
  if (!pozycja.installed) {
    pasek.append(
      przyciskKarty('Zainstaluj', 'dn-btn dn-btn--sm dn-btn--atrament', () =>
        czynnosci.zainstaluj(pozycja),
      ),
    );
  } else {
    pasek.append(
      przyciskKarty(
        pozycja.enabled ? 'Wyłącz' : 'Włącz',
        'dn-btn dn-btn--sm dn-btn--zarys',
        () => czynnosci.przelacz(pozycja),
      ),
    );
  }
  // Przycisk uprawnień jest czynny także przed instalacją, bo zakres ogląda się wcześniej.
  pasek.append(
    przyciskKarty('Uprawnienia i zaufanie', 'dn-btn dn-btn--sm dn-btn--duch', () =>
      czynnosci.wskaz(pozycja),
    ),
  );
  element.append(pasek);
  return element;
}

/**
 * Przycisk paska karty złożony z etykiety, klasy odmiany oraz funkcji
 * wywoływanej po naciśnięciu; naciśnięcie karta oddaje oknu i sama niczego
 * nie wysyła.
 */
function przyciskKarty(etykieta: string, klasa: string, naNacisniecie: () => void): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = etykieta;
  element.addEventListener('click', naNacisniecie);
  return element;
}
