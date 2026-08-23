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
 * Translation Memory Panel — okno **zarządca** modułu Translate.
 *
 * Kontrakt daje pamięci tłumaczeń dokładnie jedno wejście: podpowiedź dla
 * wskazanego segmentu w wskazanym panelu. Nie ma odczytu par, nie ma ich
 * edycji, nie ma progu dopasowania, zasięgu, wymiany TMX ani operacji
 * konserwacyjnych. Okno robi więc to jedno, co da się zrobić, i przy każdej
 * pozostałej funkcji mówi, czego brakuje — zamiast stawiać kontrolkę, która
 * nie ma czego wysłać.
 *
 * Konkordancja i podpowiedź to w tym oknie jedna czynność, bo w kontrakcie są
 * jedną komendą: pole szukania wypełnia się segmentem wybranym z tekstu
 * źródłowego albo fragmentem wpisanym ręcznie, a wynikiem jest wykaz
 * podpowiedzi rdzenia. Dwa osobne przyciski nad jedną komendą sugerowałyby dwie
 * różne zdolności.
 *
 * Wskazanie panelu jest wymagane przez kontrakt (`panelId`), a nie przez okno:
 * pamięć odpowiada w języku panelu, więc bez panelu nie ma języka, w którym
 * miałaby podpowiadać.
 */
export interface OknoTranslationMemory {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Odesłania do warsztatu — czynności pamięci, które rdzeń już wykonuje.
 *
 * To okno prowadzi jedną czynność: podpowiedź pamięci dla wskazanego segmentu.
 * Pozostałe czynności rodziny `translate.memory.*` — wykaz par, zapis pary,
 * usunięcie, wymiana z plikiem TMX, utrzymanie, tłumaczenie wstępne, wyrównanie
 * i polityka okna — mają swoje pola w oknie „Warsztat tłumaczenia".
 *
 * Nie są tu powtórzone i nie są tu nazwane brakiem: jedna czynność w dwóch
 * oknach to dwie drogi, które rozjadą się przy pierwszej zmianie kontraktu,
 * a napis „brak" przy czynności, która działa, jest zwykłą nieprawdą.
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

  /**
   * Przerysowanie odbudowuje wykaz paneli i segmentów, nie ruszając fazy
   * trwającej ani fazy błędu — te zdejmuje czynność, która je postawiła.
   * Wykaz trafień zostaje na widoku: jest wynikiem szukania, a nie odbiciem
   * stanu modułu, więc zmiana panelu w innym oknie nie ma go kasować.
   */
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
 * `translate.memory.suggest` — jedyne wejście do pamięci tłumaczeń.
 *
 * Wykaz pusty jest wynikiem, nie pustką okna: rdzeń odpowiedział i dopasowania
 * nie znalazł. Odmowa czyści wykaz, bo trafienia sprzed odmowy dotyczyłyby
 * innego zapytania niż to, które właśnie zawiodło.
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

/** Warstwa druga: nastawy pamięci — prowadzi je polityka okna w warsztacie. */
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

/** Warstwa trzecia: menu pamięci — wymiana, wykaz par i wyrównanie. */
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

/** Warstwa czwarta: konserwacja pamięci — cztery operacje wsadowe. */
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

/** zdanieOdeslania stawia jedno zdanie wskazujące okno, które czynność prowadzi. */
function zdanieOdeslania(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis';
  element.textContent = tresc;
  return element;
}
