import { KOMENDY } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';

/**
 * Funkcje panelu akcji, dla których kontrakt nie ma komendy, niosą powód składany z kontraktu
 * w czasie działania, więc dopisanie komendy przepisuje zdanie samo.
 */
export interface BrakFunkcji {
  /** Nazwa czynności tak, jak wymienia ją panel akcji okna. */
  etykieta: string;
  /** Czego by trzeba — zdanie własne okna o czynności i bytach kontraktu, których dotyczy. */
  czego: string;
  /** Nazwa komendy znoszącej ten brak; pominięta znaczy, że nie wiadomo, jak by się nazywała. */
  komendaZnoszaca?: string;
  /** Komendy z innych obszarów prowadzące tę czynność, sprawdzane w wykazie przy składaniu okna. */
  komendyCudze?: readonly string[];
}

/**
 * Komendy obszaru odczytane z kontraktu w czasie działania aplikacji, posortowane w kolejności
 * alfabetycznej dla czytelności wykazu.
 */
function komendyObszaru(obszar: string): readonly string[] {
  const przedrostek = `${obszar}.`;
  return [...(KOMENDY as readonly string[])]
    .filter((komenda) => komenda.startsWith(przedrostek))
    .sort();
}

/**
 * Czy kontrakt niesie dziś tę komendę, sprawdzane wprost w wykazie komend odczytanym z kontraktu
 * w czasie działania.
 */
function czyWKontrakcie(komenda: string): boolean {
  return (KOMENDY as readonly string[]).includes(komenda);
}

/** Obszar komendy — człon przed pierwszą kropką w pełnej nazwie komendy zapisanej dziś w kontrakcie obszaru. */
function obszarKomendy(komenda: string): string {
  const kropka = komenda.indexOf('.');
  return kropka < 0 ? komenda : komenda.slice(0, kropka);
}

/**
 * Zdanie o cudzych drogach — osobno te, które w wykazie stoją, i te, których
 * już nie ma. Rozdzielone, bo pierwsze mówią, dokąd pójść, a drugie mówią,
 * że wskazanie okna się zestarzało.
 */
function zdanieOCudzychDrogach(komendy: readonly string[]): string {
  const stojace = komendy.filter(czyWKontrakcie);
  const zniknione = komendy.filter((komenda) => !czyWKontrakcie(komenda));
  const czesci: string[] = [];
  if (stojace.length > 0) {
    const obszary = [...new Set(stojace.map(obszarKomendy))].sort();
    czesci.push(
      `Tę czynność prowadzi dziś ${stojace.join(', ')} — obszar ${obszary.join(', ')}, nie apps.`,
    );
  }
  if (zniknione.length > 0) {
    czesci.push(
      `Wskazanie okna zestarzało się: ${zniknione.join(', ')} nie stoi już w wykazie kontraktu.`,
    );
  }
  return czesci.join(' ');
}

/**
 * Pełny powód dla jednej pozycji — składany przy każdym otwarciu okna.
 *
 * Wyeksportowany, bo ta sama treść idzie do dymka, do `title` i do
 * `aria-description` przycisku, a sprawdzian sięga po nią bez budowania
 * dokumentu.
 */
export function powodBraku(brak: BrakFunkcji, obszar = 'apps'): string {
  const komendy = komendyObszaru(obszar);
  const wykaz =
    komendy.length === 0
      ? `Kontrakt nie niesie w obszarze ${obszar} ani jednej komendy — sprawdzone w wykazie kontraktu przy składaniu okna.`
      : `Kontrakt niesie dziś w obszarze ${obszar} wyłącznie: ${komendy.join(', ')} — sprawdzone w wykazie kontraktu przy składaniu okna, nie wpisane na stałe.`;
  const cudze = brak.komendyCudze === undefined ? '' : zdanieOCudzychDrogach(brak.komendyCudze);
  return [brak.czego, cudze, wykaz].filter((czesc) => czesc !== '').join(' ');
}

/**
 * Które pozycje są jeszcze brakiem — pozycja z `komendaZnoszaca` obecną
 * w kontrakcie odpada. Wyeksportowane, żeby dało się to sprawdzić bez DOM.
 */
export function brakiCzynne(braki: readonly BrakFunkcji[]): readonly BrakFunkcji[] {
  return braki.filter(
    (brak) => brak.komendaZnoszaca === undefined || !czyWKontrakcie(brak.komendaZnoszaca),
  );
}

/**
 * Wykaz czynności bez drogi w kontrakcie, osadzany w przyborniku okna, nie stawia nagłówka nad
 * pustym zbiorem pozycji czynnych.
 */
export function utworzWykazBrakow(
  tytul: string,
  braki: readonly BrakFunkcji[],
  obszar = 'apps',
): HTMLElement {
  const czynne = brakiCzynne(braki);

  const element = document.createElement('div');
  element.className = 'mp-braki__powloka';
  element.dataset['brakow'] = String(czynne.length);
  element.dataset['obszar'] = obszar;
  if (czynne.length === 0) return element;

  const naglowek = document.createElement('p');
  naglowek.className = 'mp-braki__tytul';
  naglowek.textContent = tytul;

  const lista = document.createElement('ul');
  lista.className = 'mp-braki';
  lista.append(...czynne.map((brak) => pozycja(brak, obszar)));

  element.append(naglowek, lista);
  return element;
}

/** Jedna czynność bez komendy w kontrakcie tworzy przycisk odpowiadający dymkiem, tytułem i opisem dla czytnika ekranu. */
function pozycja(brak: BrakFunkcji, obszar: string): HTMLElement {
  const powod = powodBraku(brak, obszar);

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.textContent = brak.etykieta;
  przycisk.dataset['brak'] = brak.etykieta;
  // Powód idzie równolegle trzema drogami: dymek, podpowiedź i opis dla czytnika ekranu.
  przycisk.title = powod;
  przycisk.setAttribute('aria-description', powod);
  przycisk.addEventListener('click', () => {
    pokazKomunikat({
      tytul: brak.etykieta,
      tresc: `Czynność bez drogi w kontrakcie. ${powod}`,
      waga: 'ostrz',
    });
  });

  const element = document.createElement('li');
  element.className = 'mp-braki__wiersz';
  element.append(przycisk);
  return element;
}
