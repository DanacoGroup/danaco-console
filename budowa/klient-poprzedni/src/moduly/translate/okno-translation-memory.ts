import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { PUSTE } from './etykiety-translate';
import { naglowekOkna } from './kontrolki-translate';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { oznaczWarstwe, utworzRozwiniecie } from './warstwy-translate';
import type { ZrodloPaneli } from './zrodlo-paneli';

/**
 * Translation Memory Panel to okno zarządca modułu Translate; kontrakt daje pamięci tłumaczeń
 * jedno wejście, podpowiedź dla wskazanego segmentu w wskazanym panelu.
 */
export interface OknoTranslationMemory {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Odesłania do warsztatu wskazują czynności pamięci, które rdzeń już wykonuje w oknie Warsztat
 * tłumaczenia: wykaz par, wymianę plikiem, nastawy i wyrównanie.
 */
const ODESLANIA = {
  wykaz: 'Wykaz par, zapis pary i usunięcie pary — okno „Warsztat tłumaczenia".',
  wymiana: 'Import i eksport pamięci plikiem TMX albo CSV — okno „Warsztat tłumaczenia".',
  nastawy:
    'Zasięg pamięci, próg dopasowania i zgoda na tłumaczenie wstępne stoją w polityce okna ' +
    '(polityka pamięci okna) — okno „Warsztat tłumaczenia".',
  wstepneWypelnienie:
    'Tłumaczenie wstępne z pamięci — okno „Warsztat tłumaczenia".',
  wyrownanie:
    'Wyrównanie tekstów dwujęzycznych w pary segmentów — okno „Warsztat tłumaczenia".',
  konserwacja:
    'Czyszczenie duplikatów, scalanie, masowa podmiana i usuwanie par wedle zawężenia — ' +
    'okno „Warsztat tłumaczenia".',
} as const;

export function utworzOknoTranslationMemory(stan: StanTranslate): OknoTranslationMemory {
  const okno: StanOkna = utworzStanOkna(PUSTE.pamiec);
  const odpowiedz = utworzWierszOdpowiedzi();

  const panel = poleWyboru({ etykieta: 'Panel języka docelowego' }, []);
  const fraza = poleTekstowe({
    etykieta: 'Konkordancja — segment źródłowy',
    podpowiedz: 'fragment tekstu źródłowego',
    opis: 'Kliknięcie segmentu niżej wpisuje go do tego pola.',
  });

  const segmenty = document.createElement('ul');
  segmenty.className = 'mt-pamiec__segmenty';

  const trafienia = document.createElement('ol');
  trafienia.className = 'mt-pamiec__trafienia';

  const szukaj = przycisk('Szukaj w pamięci', 'dn-btn dn-btn--sm dn-btn--atrament');
  szukaj.addEventListener('click', () => {
    void szukajWPamieci(
      stan.panele,
      { idPanelu: panel.kontrolka.value, fraza: fraza.kontrolka.value },
      { trafienia, okno, odpowiedz },
    );
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(szukaj);

  const szukanie = document.createElement('div');
  szukanie.className = 'mt-pamiec__szukanie';
  oznaczWarstwe(szukanie, 1);
  szukanie.append(panel.element, fraza.element, pasek, odpowiedz.element, segmenty, trafienia);

  okno.tresc.append(
    szukanie,
    zasiegIProg(),
    menuPamieci(),
    konserwacjaPamieci(),
  );

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--zarzadca';
  element.dataset['okno'] = 'translation-memory-panel';
  element.append(naglowekOkna('Translation Memory Panel', 'zarządca'), okno.element);

  /** Przerysowanie odbudowuje wykaz paneli i segmentów, nie ruszając fazy trwającej ani fazy błędu. */
  function odswiez(): void {
    const panele = stan.panelJezykow();
    ustawPozycje(
      panel.kontrolka,
      panele.map((wpis) => ({ wartosc: wpis.id, etykieta: `${wpis.language} — panel ${wpis.id}` })),
    );
    segmenty.replaceChildren(
      ...stan.segmenty().map((segment) =>
        wierszSegmentu(segment, () => {
          fraza.kontrolka.value = segment;
          fraza.kontrolka.focus();
        }),
      ),
    );
    if (okno.faza() === 'ladowanie' || okno.faza() === 'blad') return;
    if (trafienia.childElementCount === 0) {
      okno.puste(PUSTE.pamiec);
      return;
    }
    okno.gotowe();
  }

  odswiez();
  return { element, odswiez };
}

/**
 * Wyszukanie w pamięci tłumaczeń jest jedynym wejściem do pamięci; wykaz pusty jest wynikiem
 * szukania, a odmowa czyści wykaz trafień.
 */
async function szukajWPamieci(
  zrodlo: ZrodloPaneli,
  zapytanie: { idPanelu: string; fraza: string },
  widok: { trafienia: HTMLElement; okno: StanOkna; odpowiedz: WierszOdpowiedzi },
): Promise<void> {
  const { okno, odpowiedz } = widok;
  if (zapytanie.idPanelu === '') {
    const zdanie =
      'Wskaż panel języka docelowego — żądanie podpowiedzi niesie panel, a pamięć odpowiada ' +
      'w jego języku.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const segment = zapytanie.fraza.trim();
  if (segment === '') {
    const zdanie = 'Wpisz fragment albo wybierz segment źródłowy — bez segmentu nie ma o co pytać.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }

  okno.ladowanie('Rdzeń szuka dopasowania segmentu w pamięci tłumaczeń.');
  odpowiedz.pokaz('Szukanie w pamięci tłumaczeń…', true);
  const wynik = await zrodlo.podpowiedzPamieci(zapytanie.idPanelu, segment);
  if (!wynik.udany || wynik.wynik === undefined) {
    widok.trafienia.replaceChildren();
    const zdanie = opisOdmowy('Podpowiedź pamięci tłumaczeń', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  widok.trafienia.replaceChildren(...wynik.wynik.suggestions.map(wierszTrafienia));
  okno.gotowe();
  odpowiedz.pokaz(
    wynik.wynik.suggestions.length === 0
      ? 'Pamięć tłumaczeń nie ma dopasowania dla tego segmentu.'
      : `Pamięć tłumaczeń oddała ${String(wynik.wynik.suggestions.length)} podpowiedzi ` +
          'w kolejności dopasowania; procentu podobieństwa odpowiedź nie niesie.',
    true,
  );
}

function wierszSegmentu(tresc: string, naWybor: () => void): HTMLElement {
  const guzik = przycisk(tresc, 'dn-btn dn-btn--sm dn-btn--duch mt-pamiec__segment');
  guzik.addEventListener('click', naWybor);

  const element = document.createElement('li');
  element.append(guzik);
  return element;
}

function wierszTrafienia(tresc: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mt-pamiec__trafienie';
  element.textContent = tresc;
  return element;
}

/** Warstwa druga niesie nastawy pamięci, których zasięg i próg dopasowania prowadzi polityka okna w warsztacie tłumaczenia. */
function zasiegIProg(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 2,
    nazwa: 'Zasięg pamięci i próg dopasowania',
    wyjasnienie: ODESLANIA.nastawy,
    znacznik: '▼',
  });
  rozwiniecie.tresc.append(zdanieOdeslania(ODESLANIA.nastawy));
  return rozwiniecie.element;
}

/** Warstwa trzecia mieści menu pamięci: wymianę plikiem, wykaz par oraz wyrównanie pamięci z materiałem już przełożonym. */
function menuPamieci(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 3,
    nazwa: 'Menu pamięci',
    wyjasnienie: 'Wymiana pamięci, wykaz par, wstępne wypełnienie i wyrównanie tekstów.',
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(
    zdanieOdeslania(ODESLANIA.wykaz),
    zdanieOdeslania(ODESLANIA.wymiana),
    zdanieOdeslania(ODESLANIA.wstepneWypelnienie),
    zdanieOdeslania(ODESLANIA.wyrownanie),
  );
  return rozwiniecie.element;
}

/** Warstwa czwarta niesie konserwację pamięci złożoną z czterech operacji wsadowych wykonywanych w oknie warsztatu. */
function konserwacjaPamieci(): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 4,
    nazwa: 'Konserwacja pamięci tłumaczeń',
    wyjasnienie: ODESLANIA.konserwacja,
    znacznik: '☰',
  });
  rozwiniecie.tresc.append(zdanieOdeslania(ODESLANIA.konserwacja));
  return rozwiniecie.element;
}

/** Funkcja zdanieOdeslania stawia jedno zdanie wskazujące okno warsztatu, które czynność pamięci prowadzi zamiast tego okna. */
function zdanieOdeslania(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}
