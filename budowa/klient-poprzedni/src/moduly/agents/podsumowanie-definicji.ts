import type { Agent } from '../../../../shared/contract';
import { KOLEJNOSC_POZIOMOW } from './poziomy-pamieci';

/**
 * Panel „Podsumowanie definicji” — siedem komponentów eksperta w jednym miejscu.
 *
 * Definicja eksperta powstaje w pięciu oknach naraz i żadne z nich nie widzi
 * całości: Model Configuration nie wie, ile ekspert ma umiejętności, a Skills
 * Manager nie wie, jakim kanałem ekspert mówi. Panel jest jedynym miejscem,
 * w którym widać komplet, więc stoi przy edytorze niezależnie od tego, które
 * okno ma ognisko.
 *
 * Panel niczego nie zapisuje i nie woła ani jednej komendy — czyta eksperta
 * czynnego ze stanu modułu. Wiersz jest za to przenośnikiem: kliknięcie
 * przenosi ognisko do okna, które daną rzeczą zarządza, więc „Model: —” jest
 * drogą do Model Configuration, a nie samym stwierdzeniem braku.
 *
 * Wartości nie są tu liczone po raz drugi. Każdy wiersz czyta pole bytu
 * `Agent`, a przy pamięci i zakresie możliwości — regułę stanu wyjściowego
 * zapisaną w kontrakcie: cztery poziomy pamięci i pełny dostęp operacyjny.
 * Zdanie „brak ustawienia = wartość domyślna” obowiązuje tu tak samo jak
 * w Permissions Center.
 */
export interface PodsumowanieDefinicji {
  element: HTMLElement;
  /** Nanosi eksperta czynnego; `null` znaczy definicję jeszcze niezałożoną. */
  ustaw(ekspert: Agent | null): void;
}

/** Zależności panelu: wiersz przenosi ognisko do okna zarządzającego. */
export interface OpcjePodsumowania {
  /** Przenosi ognisko do okna wskazanego kodem zakładki edytora. */
  naZakladke(kod: string): void;
}

/** Jeden wiersz podsumowania: co pokazuje i dokąd prowadzi. */
interface OpisWiersza {
  /** Nazwa komponentu definicji. */
  nazwa: string;
  /** Kod zakładki edytora, do której wiersz przenosi ognisko. */
  zakladka: string;
  /** Wartość czytana z eksperta czynnego; `null` znaczy brak eksperta. */
  wartosc(ekspert: Agent | null): string;
}

/** Skrót instrukcji do pierwszej linii — panel jest podsumowaniem, nie edytorem. */
function pierwszaLinia(tresc: string, ile: number): string {
  const linia = tresc.split('\n', 1)[0] ?? '';
  if (linia.length <= ile) return linia;
  return `${linia.slice(0, ile)}…`;
}

/** Zdanie o zakresie możliwości — stanem wyjściowym jest pełny dostęp. */
function zdanieOZakresie(ekspert: Agent): string {
  const wpisy = ekspert.permissions ?? [];
  if (wpisy.length === 0) return 'pełny dostęp operacyjny (stan wyjściowy)';
  const odebrane = wpisy.filter((wpis) => !wpis.granted).length;
  if (odebrane === 0) return `pełny dostęp — ${wpisy.length} wpisów rozstrzygniętych wprost`;
  return `zawężony — ${odebrane} z ${wpisy.length} wpisów odebranych`;
}

/** Zdanie o pamięci — zbiór pusty jest wyłączeniem, nie brakiem odpowiedzi. */
function zdanieOPamieci(ekspert: Agent): string {
  const poziomy = ekspert.memoryLevels;
  if (poziomy.length === 0) return 'wyłączona w całości';
  if (poziomy.length === KOLEJNOSC_POZIOMOW.length) return 'cztery poziomy (stan wyjściowy)';
  return poziomy.join(' · ');
}

/** Zdanie o modelu — kanał i droga wywołania czytane wprost z bytu eksperta. */
function zdanieOModelu(ekspert: Agent): string {
  const model = ekspert.model ?? '';
  const kanal = ekspert.channelId ?? '';
  if (model === '' && kanal === '') return 'nieustalony';
  const czesci = [model === '' ? 'model kanału' : model];
  if (kanal !== '') czesci.push(`kanał ${kanal}`);
  if (ekspert.transport !== undefined) czesci.push(ekspert.transport);
  return czesci.join(' · ');
}

/** Siedem komponentów definicji w kolejności, w jakiej ekspert powstaje. */
const WIERSZE: readonly OpisWiersza[] = [
  {
    nazwa: 'Tożsamość',
    zakladka: 'tozsamosc',
    wartosc: (ekspert) => {
      if (ekspert === null) return 'definicja niezałożona';
      const imie = ekspert.displayName ?? '';
      const znak = ekspert.favicon ?? '';
      const podpis = imie === '' ? ekspert.name : `${ekspert.name} — „${imie}”`;
      return znak === '' ? podpis : `${znak} ${podpis}`;
    },
  },
  {
    nazwa: 'Instrukcje systemowe',
    zakladka: 'tozsamosc',
    wartosc: (ekspert) => {
      const tresc = ekspert?.systemPrompt ?? '';
      return tresc === '' ? 'puste' : pierwszaLinia(tresc, 60);
    },
  },
  {
    nazwa: 'Model bazowy i kanał',
    zakladka: 'model-configuration',
    wartosc: (ekspert) => (ekspert === null ? 'nieustalony' : zdanieOModelu(ekspert)),
  },
  {
    nazwa: 'Umiejętności',
    zakladka: 'skills-manager',
    wartosc: (ekspert) => String((ekspert?.skillIds ?? []).length),
  },
  {
    nazwa: 'Rozszerzenia',
    zakladka: 'connectors-manager',
    wartosc: (ekspert) => {
      const konektory = (ekspert?.connectorIds ?? []).length;
      const wtyczki = (ekspert?.pluginIds ?? []).length;
      return `${konektory} konektorów · ${wtyczki} wtyczek`;
    },
  },
  {
    nazwa: 'Pamięć',
    zakladka: 'tozsamosc',
    wartosc: (ekspert) => (ekspert === null ? 'cztery poziomy' : zdanieOPamieci(ekspert)),
  },
  {
    nazwa: 'Zakres możliwości',
    zakladka: 'permissions-center',
    wartosc: (ekspert) => (ekspert === null ? 'pełny dostęp operacyjny' : zdanieOZakresie(ekspert)),
  },
];

/**
 * Stan definicji pokazywany plakietką.
 *
 * Stanu „zarchiwizowany” tu nie ma i nie jest to przeoczenie: wykaz biblioteki
 * ekspertów zarchiwizowanych nie oddaje, więc ekspert czynny w edytorze nigdy
 * nim nie jest. Archiwum ma własny panel i własny wykaz.
 */
function stanDefinicji(ekspert: Agent | null): { napis: string; kod: string } {
  if (ekspert === null) return { napis: 'SZKIC', kod: 'szkic' };
  return ekspert.enabled
    ? { napis: 'CZYNNY', kod: 'czynny' }
    : { napis: 'WYŁĄCZONY', kod: 'wylaczony' };
}

export function utworzPodsumowanieDefinicji(
  opcje: OpcjePodsumowania,
): PodsumowanieDefinicji {
  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Podsumowanie definicji';

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka da-podsumowanie__stan';

  const naglowek = document.createElement('header');
  naglowek.className = 'da-podsumowanie__naglowek';
  naglowek.append(tytul, plakietka);

  // Wskaźnik zapisu mówi o bycie w rdzeniu, nie o formularzu: numer wersji
  // i czas ostatniej zmiany pochodzą z odpowiedzi rdzenia. Zdanie „niezapisane
  // zmiany” byłoby tu zgadywaniem — panel nie zna zawartości pól edytora.
  const zapis = document.createElement('p');
  zapis.className = 'dn-pole-opis da-podsumowanie__zapis';

  const lista = document.createElement('ul');
  lista.className = 'da-podsumowanie';

  const wartosci = new Map<string, HTMLElement>();

  for (const wiersz of WIERSZE) {
    const nazwa = document.createElement('span');
    nazwa.className = 'da-podsumowanie__nazwa';
    nazwa.textContent = wiersz.nazwa;

    const wartosc = document.createElement('span');
    wartosc.className = 'da-podsumowanie__wartosc';

    const przenosnik = document.createElement('button');
    przenosnik.type = 'button';
    przenosnik.className = 'da-podsumowanie__przenosnik';
    przenosnik.title = 'Przenieś ognisko do okna, które zarządza tym komponentem definicji.';
    przenosnik.append(nazwa, wartosc);
    przenosnik.addEventListener('click', () => opcje.naZakladke(wiersz.zakladka));

    const pozycja = document.createElement('li');
    pozycja.className = 'da-podsumowanie__wiersz';
    pozycja.dataset['komponent'] = wiersz.nazwa;
    pozycja.append(przenosnik);

    wartosci.set(wiersz.nazwa, wartosc);
    lista.append(pozycja);
  }

  const element = document.createElement('section');
  element.className = 'da-panel da-podsumowanie-panel';
  element.append(naglowek, zapis, lista);

  return {
    element,

    ustaw(ekspert) {
      const stan = stanDefinicji(ekspert);
      plakietka.textContent = stan.napis;
      plakietka.dataset['stan'] = stan.kod;
      element.dataset['ekspert'] = ekspert?.id ?? '';

      zapis.textContent =
        ekspert === null
          ? 'Definicja nie została jeszcze zapisana w rdzeniu — zapis założy wersję pierwszą.'
          : `Wersja ${ekspert.version ?? 1}, zapisana ${new Date(ekspert.updatedAt).toLocaleString('pl-PL')}.`;

      for (const wiersz of WIERSZE) {
        const cel = wartosci.get(wiersz.nazwa);
        if (cel !== undefined) cel.textContent = wiersz.wartosc(ekspert);
      }
    },
  };
}
