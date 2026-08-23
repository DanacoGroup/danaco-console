import { ExtensionOrigin, type Extension } from '../../../../shared/contract';

/**
 * Karta pozycji katalogu rozszerzeń — siatka kart App Catalogu.
 *
 * Karta pokazuje wyłącznie pola, które rdzeń oddał: nazwę, rodzaj, wersję,
 * pochodzenie, opis oraz dwa stany — zainstalowania i włączenia. Pola, których
 * pozycja nie niesie (dziennik zmian, wykaz narzędzi, wymagane uprawnienia),
 * nie mają tu zaślepki: stoją w wykazie braków okna, żeby nieobecność była
 * widoczna zamiast udawanej.
 *
 * Stan nigdy nie opiera się na samej barwie — plakietka niesie słowo, a nie
 * tylko odmianę. Cztery stany opracowania (dostępna, zainstalowana wyłączona,
 * zainstalowana włączona, dostępna aktualizacja) składamy z dwóch pól kontraktu;
 * czwartego stanu nie udajemy, bo rejestr nie niesie wersji dostępnej.
 *
 * Karta nie zna kanału i niczego nie wysyła — oddaje naciśnięcia oknu.
 */
export interface CzynnosciKarty {
  /** Instalacja pozycji jeszcze niezainstalowanej. */
  zainstaluj(pozycja: Extension): void;
  /** Włączenie albo wyłączenie pozycji zainstalowanej. */
  przelacz(pozycja: Extension): void;
  /** Wskazanie pozycji panelom bocznym — uprawnieniom i konsoli. */
  wskaz(pozycja: Extension): void;
}

/** Napisy pochodzenia; wartości kontraktu mają w interfejsie nazwy Operatora. */
const NAZWY_POCHODZENIA: Readonly<Record<string, string>> = {
  [ExtensionOrigin.Danaco]: 'Danaco Plugin',
  [ExtensionOrigin.Personal]: 'Personal',
};

/** Nazwa pochodzenia albo sama wartość, gdy rdzeń oddał wartość spoza wyliczenia. */
export function nazwaPochodzenia(pochodzenie: string): string {
  return NAZWY_POCHODZENIA[pochodzenie] ?? pochodzenie;
}

/**
 * Zdanie o stanie pozycji — słowo, nie barwa.
 *
 * Rozróżnia trzy stany, bo tyle niosą dwa pola kontraktu. Czwarty stan
 * opracowania („dostępna aktualizacja") wymaga zestawienia wersji zainstalowanej
 * z wersją rejestru, a tej drugiej pozycja nie niesie.
 */
export function opisStanu(pozycja: Extension): string {
  if (!pozycja.installed) return 'dostępna';
  return pozycja.enabled ? 'zainstalowana, włączona' : 'zainstalowana, wyłączona';
}

/** Odmiana plakietki stanu; nazwy z biblioteki `komponenty/`. */
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

  // Instalacja i przełączenie to dwie różne komendy, więc dwa różne przyciski.
  // Jeden przycisk o zmiennym napisie kazałby Operatorowi zgadywać, którą
  // czynność zaraz zleci.
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
  // Znacznik uprawnień opracowania: karta otwiera Permissions & Trust Center
  // na swojej pozycji. Przycisk jest czynny także dla pozycji niezainstalowanej,
  // bo zakres dostępu ogląda się PRZED włączeniem.
  pasek.append(
    przyciskKarty('Uprawnienia i zaufanie', 'dn-btn dn-btn--sm dn-btn--duch', () =>
      czynnosci.wskaz(pozycja),
    ),
  );
  element.append(pasek);
  return element;
}

/** Przycisk karty; naciśnięcie oddaje pozycję oknu, karta nic nie wysyła. */
function przyciskKarty(etykieta: string, klasa: string, naNacisniecie: () => void): HTMLElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = etykieta;
  element.addEventListener('click', naNacisniecie);
  return element;
}
