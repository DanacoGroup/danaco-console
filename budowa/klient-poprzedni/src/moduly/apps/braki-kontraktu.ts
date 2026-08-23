import { KOMENDY } from '../../../../shared/contract';
import { pokazKomunikat } from '../../aplikacja/komunikaty';

/**
 * Funkcje panelu akcji, dla których kontrakt nie ma komendy — z powodem
 * składanym z kontraktu w czasie działania, nie wpisanym na stałe.
 *
 * Panele akcji okien Apps wymieniają więcej czynności, niż obszar `apps` niesie
 * komend. Czynność bez komendy zostaje widoczna i klikalna, a naciśnięcie mówi,
 * czego brakuje.
 *
 * Powód składa się z `KOMENDY` w `shared/contract.ts` przy składaniu okna, więc
 * dopisanie komendy do kontraktu przepisuje zdanie samo. Wobec bliźniaczego
 * pliku modułu Diagnostics (`moduly/diagnostics/braki-kontraktu.ts`) dochodzą
 * tu dwie rzeczy:
 *   — pozycja znika sama, gdy wskazana komenda wejdzie do kontraktu
 *     (`komendaZnoszaca`),
 *   — wskazanie cudzej drogi jest sprawdzane (`komendyCudze`): zdanie o komendzie
 *     z innego obszaru pada wyłącznie wtedy, gdy ta komenda stoi w wykazie;
 *     w przeciwnym razie zdanie mówi o jej zniknięciu.
 *
 * Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje — brak jest po
 * stronie kontraktu i tylko o kontrakcie zdanie mówi.
 */
export interface BrakFunkcji {
  /** Nazwa czynności tak, jak wymienia ją panel akcji okna. */
  etykieta: string;
  /**
   * Czego by trzeba — zdanie własne okna, mówiące o czynności i o bytach
   * kontraktu, których dotyczy. O stanie kontraktu nie orzeka niczego; to
   * dokłada się niżej z wykazu komend.
   */
  czego: string;
  /**
   * Nazwa komendy, której wejście do kontraktu znosi ten brak. Pozycja z taką
   * komendą już obecną w `KOMENDY` nie trafia do wykazu wcale.
   *
   * Pominięta znaczy „nie wiadomo, jak taka komenda miałaby się nazywać".
   */
  komendaZnoszaca?: string;
  /**
   * Komendy z innych obszarów, które tę czynność prowadzą. Sprawdzane w wykazie
   * przy składaniu okna: obecne wchodzą do zdania jako droga istniejąca, ale
   * cudza; nieobecne wchodzą jako wskazanie, które przestało być prawdziwe.
   */
  komendyCudze?: readonly string[];
}

/** Komendy obszaru odczytane z kontraktu w czasie działania, uporządkowane. */
function komendyObszaru(obszar: string): readonly string[] {
  const przedrostek = `${obszar}.`;
  return [...(KOMENDY as readonly string[])]
    .filter((komenda) => komenda.startsWith(przedrostek))
    .sort();
}

/** Czy kontrakt niesie dziś tę komendę. */
function czyWKontrakcie(komenda: string): boolean {
  return (KOMENDY as readonly string[]).includes(komenda);
}

/** Obszar komendy — człon przed pierwszą kropką. */
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
 * Wykaz czynności bez drogi w kontrakcie, osadzany w przyborniku okna.
 *
 * Każda pozycja jest przyciskiem. Naciśnięcie nie wysyła nic — pokazuje dymek
 * nazywający brakującą komendę. Wykaz bez czynnych pozycji nie stawia nagłówka:
 * „Bez drogi w kontrakcie" nad niczym byłoby zdaniem o braku, którego nie ma.
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

/** Jedna czynność bez komendy: przycisk odpowiadający dymkiem. */
function pozycja(brak: BrakFunkcji, obszar: string): HTMLElement {
  const powod = powodBraku(brak, obszar);

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--duch dn-btn--sm';
  przycisk.textContent = brak.etykieta;
  przycisk.dataset['brak'] = brak.etykieta;
  // Powód idzie równolegle trzema drogami — dymek, podpowiedź i opis dla
  // czytnika ekranu — bo dymek widzi wyłącznie ten, kto naciśnie.
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
