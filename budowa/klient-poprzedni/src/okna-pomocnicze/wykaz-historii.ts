import type { HistoryEntry } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Wykaz pozycji historii rozmowy okna wraz z ich zaznaczaniem.
 *
 * Osobno od panelu (`okno-historii-rozmowy.ts`), bo panel prowadzi rozmowy
 * z rdzeniem, stan wczytywania i nastawę zasady, a wykaz odpowiada wyłącznie
 * za zamianę pozycji na wiersze i za wskazanie, które z nich Operator zaznaczył.
 *
 * Zaznaczenie wskazuje zakres, nie potwierdza czynności: panel ma osobną
 * czynność na wyczyszczenie całej historii okna i osobną na usunięcie pozycji
 * wskazanych. Wiersz bez zaznaczenia niczego nie blokuje.
 *
 * Pozycje, których nastawiona zasada nie utrzyma, są oznaczone, a nie ukryte
 * ani wygaszone — nastawa retencji ma pokazać skutek zasady przed jej zapisem,
 * więc wiersz oznaczony nadal da się przeczytać i zaznaczyć.
 *
 * Wykaz nie zna żadnej komendy i nie wie, skąd pozycje przyszły: dostaje
 * tablicę, oddaje identyfikatory.
 */
export interface WykazHistorii {
  /** Element wstawiany w miejsce treści panelu. */
  element: HTMLElement;
  /**
   * Przerysowuje wykaz.
   *
   * @param pozycje pozycje od najnowszej, wprost z `history.load`.
   * @param pozaZasada identyfikatory pozycji, których nastawiona zasada nie
   *   utrzyma. Zbiór pusty znaczy, że zasada niczego nie usuwa.
   */
  rysuj(pozycje: readonly HistoryEntry[], pozaZasada: ReadonlySet<string>): void;
  /** Identyfikatory pozycji wskazanych przez Operatora. */
  zaznaczone(): string[];
  /** Zgłasza każdą zmianę zaznaczenia — panel przelicza podpis czynności. */
  naZmianeZaznaczenia(sluchacz: () => void): Odsubskrybuj;
}

/** Ile znaków skrótu treści pokazuje wiersz, gdy rdzeń przysłał go dłuższy. */
const DLUGOSC_SKROTU_WIDOKU = 240;

export function utworzWykazHistorii(): WykazHistorii {
  const zmiany = utworzMagistrale<void>();
  /** Zaznaczenie przeżywa przerysowanie — pozycja skasowana wypada przy odczycie. */
  const wskazane = new Set<string>();

  const element = document.createElement('div');
  element.className = 'dnp-historia';

  function rysuj(pozycje: readonly HistoryEntry[], pozaZasada: ReadonlySet<string>): void {
    // Zaznaczenie pozycji, której w świeżym wykazie już nie ma, przestaje
    // istnieć razem z nią. Bez tego panel wysłałby do rdzenia identyfikator
    // skasowany i dostał w odpowiedzi „usunięto 0".
    const zywe = new Set(pozycje.map((pozycja) => pozycja.id));
    for (const kod of [...wskazane]) if (!zywe.has(kod)) wskazane.delete(kod);

    const lista = document.createElement('ul');
    lista.className = 'dnp-historia__lista';
    for (const pozycja of pozycje) {
      lista.append(
        wiersz(pozycja, pozaZasada.has(pozycja.id), wskazane.has(pozycja.id), przelacz),
      );
    }
    element.replaceChildren(lista);
  }

  function przelacz(kod: string, wskazany: boolean): void {
    if (wskazany) wskazane.add(kod);
    else wskazane.delete(kod);
    zmiany.oglos();
  }

  return {
    element,
    rysuj,
    zaznaczone: () => [...wskazane],
    naZmianeZaznaczenia: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Jeden wiersz wykazu: wskazanie, nadawca, czas, skrót treści. */
function wiersz(
  pozycja: HistoryEntry,
  pozaZasada: boolean,
  wskazany: boolean,
  przelacz: (kod: string, wskazany: boolean) => void,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dnp-historia__wiersz';
  element.dataset['poza'] = String(pozaZasada);

  const wskazanie = document.createElement('input');
  wskazanie.type = 'checkbox';
  wskazanie.className = 'dn-przelacznik';
  // Wiersz powstaje na nowo przy każdym przerysowaniu, więc pole wskazania
  // trzeba nastawić z zaznaczenia trzymanego przez wykaz.
  wskazanie.checked = wskazany;
  wskazanie.setAttribute('aria-label', `Wskaż pozycję ${opisNadawcy(pozycja.role)} z ${czas(pozycja.createdAt)}`);
  wskazanie.addEventListener('change', () => przelacz(pozycja.id, wskazanie.checked));

  const naglowek = document.createElement('p');
  naglowek.className = 'dnp-historia__naglowek';
  naglowek.textContent = `${opisNadawcy(pozycja.role)} · ${czas(pozycja.createdAt)}`;

  const tresc = document.createElement('p');
  tresc.className = 'dnp-historia__tresc';
  tresc.textContent = skrot(pozycja);

  const opis = document.createElement('div');
  opis.className = 'dnp-historia__opis';
  opis.append(naglowek, tresc);
  if (pozaZasada) opis.append(znacznikPozaZasada());

  element.append(wskazanie, opis);
  return element;
}

/**
 * Napis przy pozycji, której nastawiona zasada nie utrzyma.
 *
 * Zapowiedź, nie orzeczenie: dopóki Operator nie zapisze zasady, pozycja stoi
 * w bazie nietknięta. Zdanie mówi to wprost, bo napis „poza zasadą" bez tego
 * czytałby się jak informacja o kasowaniu, które już nastąpiło.
 */
function znacznikPozaZasada(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dnp-historia__poza';
  element.setAttribute('role', 'note');
  element.textContent = 'Nastawiona zasada tej pozycji nie utrzyma — zniknie po jej zapisaniu.';
  return element;
}

/**
 * Skrót treści. Rdzeń składa `preview` sam (200 znaków, `handlers_historia.go`)
 * i to jest treść właściwa. Pole jest niewymagane, więc jego brak dostaje
 * zdanie wprost, a nie pusty wiersz nie do odróżnienia od wypowiedzi bez tekstu.
 */
function skrot(pozycja: HistoryEntry): string {
  const tresc = pozycja.preview ?? '';
  if (tresc.trim() === '') return '(rdzeń nie przysłał skrótu tej pozycji)';
  return tresc.length > DLUGOSC_SKROTU_WIDOKU ? `${tresc.slice(0, DLUGOSC_SKROTU_WIDOKU)}…` : tresc;
}

/**
 * Nadawca po polsku. Kontrakt daje `role` jako dowolny napis, nie wyliczenie,
 * dlatego wartość nieznana idzie do Operatora taka, jaka przyszła, zamiast
 * zostać wciśnięta w jedną z dwóch znanych nazw.
 */
function opisNadawcy(rola: string): string {
  if (rola === 'user' || rola === 'operator') return 'Operator';
  if (rola === 'assistant' || rola === 'model') return 'Model';
  return rola === '' ? 'nadawca nieznany' : rola;
}

/** Czas zapisu w postaci lokalnej; wartość nieczytelna zostaje nazwana. */
function czas(znacznik: number): string {
  if (!Number.isFinite(znacznik) || znacznik <= 0) return 'czas nieznany';
  return new Date(znacznik).toLocaleString('pl-PL');
}
