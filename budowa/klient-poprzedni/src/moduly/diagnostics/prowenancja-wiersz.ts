import type {
  ModelCallSpan,
  ModelCallTrace,
  ProvenanceCallGetResponse,
} from '../../../../shared/contract';
import { przelacznik, przyciskAkcji, pozycjaWykazu } from '../../modele/kontrolki-formularza-braki';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import { utworzOceneWywolania } from './prowenancja-ocena';
import { czas, opisWywolania, RODZAJ_ODCINKA, tytulWywolania } from './prowenancja-slowa';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';

/**
 * Wiersz jednego wywołania kanału modelu w przeglądarce prowenancji łączy odczyt
 * śladu, ocenę odpowiedzi i wskazanie wywołania do wydania — trzy czynności
 * dotyczące tego samego wywołania.
 */
export interface WierszWywolania {
  element: HTMLElement;
  /** Czy Operator wskazał to wywołanie do wydania śladu. */
  wskazane(): boolean;
  /** Identyfikator wywołania — po nim składa się pole `callIds` wydania. */
  id(): string;
}

export function utworzWierszWywolania(
  zrodlo: ZrodloProwenancji,
  wywolanie: ModelCallTrace,
  potwierdz: (zdanie: string, udane: boolean) => void,
): WierszWywolania {
  let biezace = wywolanie;

  const pozycja = pozycjaWykazu(tytulWywolania(biezace), opisWywolania(biezace), 'dg');
  pozycja.element.dataset['stanWywolania'] = biezace.status;
  pozycja.element.dataset['wywolanie'] = biezace.id;

  const szczegoly = document.createElement('div');
  szczegoly.className = 'dg-tresc';

  const odczyt = przyciskAkcji('Odczytaj ślad wywołania', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczyt.addEventListener('click', () => {
    szczegoly.replaceChildren(zdanie('Odczyt śladu wywołania z rdzenia…'));
    void zrodlo
      .odczytajWywolanie({ callId: biezace.id, includeContent: true })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          const opis = zdanieNiepowodzenia('odczytu śladu wywołania', wynik.powod);
          szczegoly.replaceChildren(zdanie(opis));
          potwierdz(opis, false);
          return;
        }
        szczegoly.replaceChildren(...slad(wynik.wynik));
        biezace = wynik.wynik.call;
        odswiezOpis();
        potwierdz(zdanieSladu(wynik.wynik), true);
      });
  });

  const ocen = przyciskAkcji('Oceń odpowiedź modelu', 'dn-btn dn-btn--sm dn-btn--zarys');
  ocen.addEventListener('click', () => {
    const ocena = utworzOceneWywolania(zrodlo, biezace, (ocenione, opis, udane) => {
      if (udane) {
        biezace = ocenione;
        odswiezOpis();
      }
      potwierdz(opis, udane);
    });
    szczegoly.replaceChildren(ocena.element);
  });

  const doWydania = przelacznik('Obejmij to wywołanie wydaniem śladu');
  doWydania.dataset['doWydania'] = biezace.id;

  function odswiezOpis(): void {
    const opis = pozycja.element.querySelector('.dg-pozycja__opis');
    if (opis !== null) opis.textContent = opisWywolania(biezace);
    pozycja.element.dataset['stanWywolania'] = biezace.status;
  }

  pozycja.akcje.append(odczyt, ocen, doWydania);
  pozycja.element.append(szczegoly);

  return {
    element: pozycja.element,
    wskazane: () => doWydania.checked,
    id: () => biezace.id,
  };
}

/** Akapit zdania w miejscu szczegółów wiersza, zastępujący pojedynczym tekstem każdą kolejną treść odczytu. */
function zdanie(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}

/**
 * Ślad jednego wywołania: drzewo odcinków budowane z pola nadrzędnika w kolejności
 * od korzenia, wraz z treścią promptu i odpowiedzi.
 */
function slad(odczyt: ProvenanceCallGetResponse): readonly HTMLElement[] {
  const czesci: HTMLElement[] = [];

  czesci.push(
    zdanie(
      odczyt.spans.length === 0
        ? 'Drzewa odcinków to wywołanie nie ma — rdzeń oddał ślad bez ani jednego odcinka, ' +
          'i jest to stan poprawny dla wywołania prostego.'
        : `Odcinki drzewa śladu (${String(odczyt.spans.length)}), w kolejności od korzenia:`,
    ),
  );
  if (odczyt.spans.length > 0) {
    const drzewo = document.createElement('pre');
    drzewo.className = 'dg-poprawka';
    drzewo.dataset['drzewoSladu'] = odczyt.call.id;
    drzewo.textContent = odczyt.spans.map((odcinek) => wierszOdcinka(odcinek, odczyt.spans)).join('\n');
    czesci.push(drzewo);
  }

  czesci.push(
    zdanie(
      odczyt.redacted
        ? 'Treść poniżej jest PO REDAKCJI danych wrażliwych — postać pierwotna promptu ' +
          'i odpowiedzi tutaj nie dochodzi.'
        : 'Treść poniżej jest w postaci pierwotnej, bez redakcji danych wrażliwych.',
    ),
  );
  czesci.push(tresc('Prompt wysłany do kanału modelu', odczyt.prompt, odczyt.call.contentStored));
  czesci.push(tresc('Odpowiedź modelu', odczyt.response, odczyt.call.contentStored));
  return czesci;
}

/**
 * Jeden odcinek w zapisie tekstowym wraz z wcięciem wg głębokości.
 *
 * Wcięcie liczy się przez łańcuch rodziców, a nie przez licznik pętli: odcinki
 * rodzeństwa mają tę samą głębokość niezależnie od tego, ile odcinków stoi
 * między nimi w wykazie.
 */
function wierszOdcinka(odcinek: ModelCallSpan, wszystkie: readonly ModelCallSpan[]): string {
  let glebokosc = 0;
  let rodzic = odcinek.parentSpanId;
  // Granica przebiegu jest granicą wykazu: cykl w danych rdzenia nie może
  // zapętlić okna.
  while (rodzic !== undefined && glebokosc <= wszystkie.length) {
    const wyzej = wszystkie.find((wpis) => wpis.id === rodzic);
    if (wyzej === undefined) break;
    glebokosc += 1;
    rodzic = wyzej.parentSpanId;
  }
  const czesci = [
    `${'  '.repeat(glebokosc)}${odcinek.name} — ${RODZAJ_ODCINKA[odcinek.kind]}`,
    `początek ${czas(odcinek.startedAt)}`,
  ];
  if (odcinek.durationMs !== undefined) czesci.push(`trwał ${String(odcinek.durationMs)} ms`);
  if (odcinek.tokens !== undefined) czesci.push(`tokenów ${String(odcinek.tokens)}`);
  if (odcinek.cacheHit === true) czesci.push('trafienie w pamięć podręczną promptu');
  return czesci.join(' · ');
}

/**
 * Treść promptu albo odpowiedzi, w której puste miejsce rozróżnia wyłączony zapis,
 * brak treści od rdzenia i treść rzeczywiście pustą.
 */
function tresc(nazwa: string, zawartosc: string | undefined, zapisana: boolean): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'dg-poprawka';
  element.dataset['trescSladu'] = nazwa;
  if (!zapisana) {
    element.textContent =
      `${nazwa}: rdzeń treści NIE zapisał — zapis treści promptu i odpowiedzi był w chwili ` +
      'wywołania wyłączony ustawieniem rdzenia.';
    return element;
  }
  if (zawartosc === undefined) {
    element.textContent = `${nazwa}: rdzeń tej treści nie oddał, choć wywołanie ma ją zapisaną.`;
    return element;
  }
  element.textContent = `${nazwa}:\n${zawartosc}`;
  return element;
}

/** Zdanie potwierdzenia po odczycie śladu, podające liczbę odcinków i długość treści wprost, wraz ze stanem redakcji. */
function zdanieSladu(odczyt: ProvenanceCallGetResponse): string {
  return (
    `Rdzeń oddał ślad wywołania: odcinków ${String(odczyt.spans.length)}, ` +
    `znaków promptu ${String(odczyt.prompt?.length ?? 0)}, ` +
    `znaków odpowiedzi ${String(odczyt.response?.length ?? 0)}` +
    (odczyt.redacted ? ' — treść PO redakcji danych wrażliwych.' : ' — treść bez redakcji.')
  );
}
