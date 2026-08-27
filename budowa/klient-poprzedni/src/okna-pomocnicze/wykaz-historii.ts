import type { HistoryEntry } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';

/**
 * Wykaz pozycji historii rozmowy okna wraz z ich zaznaczaniem stoi osobno od panelu, bo panel prowadzi rozmowę z rdzeniem i nastawę zasady, a wykaz odpowiada wyłącznie za zamianę pozycji na wiersze.
 */
export interface WykazHistorii {
  /** Element wstawiany w miejsce treści panelu. */
  element: HTMLElement;
  /** Przerysowuje wykaz; pozycje idą od najnowszej, a zbiór poza zasadą pusty znaczy brak usunięć. */
  rysuj(pozycje: readonly HistoryEntry[], pozaZasada: ReadonlySet<string>): void;
  /** Identyfikatory pozycji wskazanych przez Operatora. */
  zaznaczone(): string[];
  /** Zgłasza każdą zmianę zaznaczenia — panel przelicza podpis czynności. */
  naZmianeZaznaczenia(sluchacz: () => void): Odsubskrybuj;
}

/** Ile znaków skrótu treści pokazuje wiersz wykazu, gdy rdzeń przysłał treść dłuższą niż ten limit widoku. */
const DLUGOSC_SKROTU_WIDOKU = 240;

export function utworzWykazHistorii(): WykazHistorii {
  const zmiany = utworzMagistrale<void>();
  /** Zaznaczenie przeżywa przerysowanie — pozycja skasowana wypada przy odczycie. */
  const wskazane = new Set<string>();

  const element = document.createElement('div');
  element.className = 'dnp-historia';

  function rysuj(pozycje: readonly HistoryEntry[], pozaZasada: ReadonlySet<string>): void {
    // Zaznaczenie pozycji, której w wykazie już nie ma, znika razem z nią, by uniknąć zbędnego wysłania.
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

/** Jeden wiersz wykazu historii: pole wskazania, nadawca, czas zapisu i skrót treści całej wiadomości Operatora. */
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
  // Wiersz powstaje na nowo przy przerysowaniu, więc pole wskazania nastawia się z zaznaczenia wykazu.
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
 * Napis przy pozycji, której nastawiona zasada nie utrzyma, jest zapowiedzią, nie orzeczeniem: dopóki Operator nie zapisze zasady, pozycja stoi w bazie nietknięta.
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

/** Czas zapisu pozycji przepisany na postać lokalną dla Operatora; wartość nieczytelna zostaje nazwana wprost. */
function czas(znacznik: number): string {
  if (!Number.isFinite(znacznik) || znacznik <= 0) return 'czas nieznany';
  return new Date(znacznik).toLocaleString('pl-PL');
}
