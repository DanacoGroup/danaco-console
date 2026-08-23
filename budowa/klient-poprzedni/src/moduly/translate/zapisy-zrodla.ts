import type { TranslateSourceSetRequest } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { BRAK_OKNA } from './etykiety-translate';
import { zdanieOdmowyModelu } from './odmowa-translate';
import type { StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';

/**
 * Trzy ścieżki zapisu Source Panel — oddzielone od widoku, który je wyzwala.
 *
 * Odczyt i zapis mieszkają osobno, tak jak w sekcji dostępów: widok składa
 * pola i przyciski, a to, co się dzieje po naciśnięciu, jest tutaj. Dzięki temu
 * sprawdzian może wywołać zapis bez klikania w element, a plik widoku nie
 * puchnie o obsługę odmów.
 *
 * Każda ścieżka odpowiada. Brak okna, pusty tekst, odmowa rdzenia — wszystkie
 * trzy kończą się zdaniem w wierszu odpowiedzi, żadna ciszą.
 *
 * Zdanie o wyniku mówi tylko to, co niesie odpowiedź rdzenia:
 *
 *  — `translate.source.set` żadnego przekładu nie uruchamia; oddaje panele
 *    z ich dotychczasową treścią, a przelicza je dopiero `translate.target.add`
 *    albo korekta Operatora;
 *  — przy zapisie bez wskazanego języka źródłowego rdzeń oddaje
 *    `sourceLanguage: ""`; pusty język jest nazwany, nie przemilczany.
 *
 * Rozpoznanie języka nie wpisuje się do pola w ciemno. `translate.source.detect`
 * przepuszcza odpowiedź modelu wprost do pola `language` kontraktu, więc
 * w polu potrafi wrócić komunikat kanału zamiast oznaczenia języka. Wartość
 * niebędąca oznaczeniem języka wraca Operatorowi dosłownie, jako odmowa,
 * i nie nadpisuje tego, co sam wpisał.
 */
export interface PolaZrodla {
  tekst: HTMLTextAreaElement;
  jezyk: HTMLInputElement;
  ponownaSegmentacja: HTMLInputElement;
}

export interface ZapisyZrodla {
  /** `translate.source.set` — zapis tekstu; wyzwala aktualizację wszystkich paneli. */
  zapisz(): Promise<void>;
  /** `translate.source.segment` — ponowny podział na segmenty. */
  segmentuj(): Promise<void>;
  /** `translate.source.detect` — rozpoznanie języka źródłowego. */
  rozpoznaj(): Promise<void>;
}

export function utworzZapisyZrodla(
  stan: StanTranslate,
  pola: PolaZrodla,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): ZapisyZrodla {
  /** Czy moduł ma czym zaadresować żądanie. Bez okna nie wysyłamy niczego. */
  function maOkno(): boolean {
    if (stan.idOkna() !== '') return true;
    odpowiedz.pokaz(BRAK_OKNA, false);
    okno.blad(BRAK_OKNA);
    return false;
  }

  return {
    async zapisz() {
      if (!maOkno()) return;
      await zapiszZrodlo(stan, pola, okno, odpowiedz);
    },
    async segmentuj() {
      await segmentujPonownie(stan, pola.tekst.value, okno, odpowiedz);
    },
    async rozpoznaj() {
      await rozpoznajJezyk(stan, pola.jezyk, pola.tekst.value, okno, odpowiedz);
    },
  };
}

/** `translate.source.set` — zapis tekstu źródłowego wraz z jego językiem. */
async function zapiszZrodlo(
  stan: StanTranslate,
  pola: PolaZrodla,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const tekst = pola.tekst.value;
  if (tekst.trim() === '') {
    // Zdanie mówi o oknie, bo to okno odmawia, nie rdzeń. Rdzeń sprawdza
    // `Text == ""` dokładnie: pustego tekstu odmawia, ale samą spację czy sam
    // znacznik kolejności bajtów przyjmuje i zapisuje. Okno takiego źródła nie
    // wysyła — zapis samych białych znaków nie jest zapisem tekstu — i nie
    // przypisuje rdzeniowi odmowy, której by nie było.
    odpowiedz.pokaz(
      'Wpisz albo wklej tekst — pole ma same białe znaki, więc okno nie ma czego zapisać.',
      false,
    );
    return;
  }
  const zadanie: TranslateSourceSetRequest = { windowId: stan.idOkna(), text: tekst };
  const jezyk = pola.jezyk.value.trim();
  if (jezyk !== '') zadanie.sourceLanguage = jezyk;
  if (pola.ponownaSegmentacja.checked) zadanie.resegment = true;

  // Stan ładowania idzie na okno, nie tylko w wiersz odpowiedzi: wiersz mówi
  // o czynności, pas stanu o oknie. Wskaźnik stoi przy tytule pasa, a treść okna
  // zostaje widoczna i edytowalna — przesłonięcie pola skasowałoby Operatorowi
  // z oczu to, co właśnie pisze.
  okno.ladowanie('Rdzeń przyjmuje tekst źródłowy i dzieli go na segmenty.');
  odpowiedz.pokaz('Zapisywanie tekstu źródłowego…', true);
  const wynik = await stan.zrodlo.ustaw(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowy('Zapis tekstu źródłowego', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  // Tekst wchodzi do stanu z żądania, bo `TranslateSourceSetResponse` go nie
  // odsyła (niesie sourceLanguage, segmentCount i panels). Nie ma więc czym
  // potwierdzić treści źródła i zdanie niżej jej nie potwierdza: mówi o języku
  // i o podziale, które rdzeń oddał.
  stan.wchlonZrodlo(tekst, wynik.wynik.sourceLanguage, []);
  stan.wchlonPanele(wynik.wynik.panels);

  // Język zamówiony porównuje się z oddanym wyłącznie wtedy, gdy Operator język
  // wskazał: przy polu pustym żądanie go nie niesie, więc nie ma zamówienia,
  // z którym można by zestawić odpowiedź. Rdzeń oddaje wskazany język wprost,
  // a przy żądaniu bez języka oddaje pole puste, kasując język zapisany
  // wcześniej — obie drogi mają w oknie własne zdanie.
  if (jezyk !== '') {
    const rozbiezne = rozbieznoscOdpowiedzi('Zapis tekstu źródłowego', [
      { nazwa: 'język źródłowy', zamowione: jezyk, oddane: wynik.wynik.sourceLanguage },
    ]);
    if (rozbiezne !== null) {
      odpowiedz.pokaz(rozbiezne, false);
      okno.blad(rozbiezne);
      return;
    }
  }

  const ile = wynik.wynik.segmentCount ?? wynik.wynik.panels.length;
  const jezykWyniku = wynik.wynik.sourceLanguage.trim();
  const oJezyku =
    jezykWyniku === ''
      ? 'Źródło zapisane; rdzeń nie oddał języka źródłowego (pole puste — nie rozpoznał go i nie zgadywał)'
      : `Źródło zapisane w języku ${jezykWyniku}`;
  const oPanelach =
    wynik.wynik.panels.length === 0
      ? 'Rdzeń nie zna jeszcze żadnego panelu tego okna — dodaj język docelowy, żeby powstał przekład.'
      : `Rdzeń oddał ${wynik.wynik.panels.length} paneli tego okna z ich dotychczasową treścią; ` +
        'sam zapis źródła nie przelicza przekładu.';
  odpowiedz.pokaz(`${oJezyku}; rdzeń zgłosił ${ile} pozycji podziału. ${oPanelach}`, true);
  okno.gotowe();
}

/** `translate.source.segment` — liczba segmentów bierze się z odpowiedzi. */
async function segmentujPonownie(
  stan: StanTranslate,
  tekst: string,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  okno.ladowanie('Rdzeń dzieli tekst źródłowy na segmenty.');
  odpowiedz.pokaz('Ponowny podział na segmenty…', true);
  const wynik = await stan.zrodlo.segmentuj(tekst);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowy('Ponowna segmentacja', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  // Faza schodzi przed ogłoszeniem stanu, nie po nim. Przerysowanie okna nie
  // zdejmuje fazy trwającej (patrz `odswiez` w `okno-source-panel.ts`), więc
  // ogłoszenie wykonane przy zapalonym ładowaniu przeszłoby bez skutku i okno
  // zostałoby z zapowiedzią wywołania, które się już skończyło.
  okno.gotowe();
  stan.wchlonSegmenty(wynik.wynik.segments);
  odpowiedz.pokaz(`Rdzeń podzielił tekst na ${wynik.wynik.segments.length} segmentów.`, true);
}

/** `translate.source.detect` — rozpoznanie języka modelem, z odsiewem kształtu. */
async function rozpoznajJezyk(
  stan: StanTranslate,
  poleJezyka: HTMLInputElement,
  tekst: string,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  okno.ladowanie('Rdzeń rozpoznaje język źródłowy modelem.');
  odpowiedz.pokaz('Rozpoznawanie języka źródłowego…', true);
  const wynik = await stan.zrodlo.rozpoznajJezyk(tekst);
  if (!wynik.udany || wynik.wynik === undefined) {
    // Rozpoznanie języka rdzeń wykonuje modelem, więc odmowa z braku kanału ma
    // tu nieść zdanie o tym, czego brakuje (`odmowa-translate`).
    const zdanie = zdanieOdmowyModelu('Rozpoznanie języka', wynik.blad);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const oddany = wynik.wynik.language.trim();
  if (!wygladaNaJezyk(oddany)) {
    // Odpowiedź nieużyteczna jest odmową, nie pustką. Rdzeń odpowiedział, ale
    // tym, czego okno nie ma prawa wpisać do pola języka — pas stanu mówi to
    // samo, co wiersz odpowiedzi, zamiast wracać do „gotowe".
    const zdanie =
      'Rozpoznanie języka: rdzeń oddał w polu języka treść, która nie jest oznaczeniem ' +
      `języka — „${oddany}". Pola nie nadpisano; wpisz kod języka sam albo sprawdź ` +
      'kanał modelu, którym rdzeń rozpoznaje język.';
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  poleJezyka.value = oddany;
  okno.gotowe();
  const pewnosc = wynik.wynik.confidence;
  odpowiedz.pokaz(
    pewnosc === undefined
      ? `Rdzeń rozpoznał język ${oddany}.`
      : `Rdzeń rozpoznał język ${oddany} z pewnością ${pewnosc}.`,
    true,
  );
}

/**
 * Czy oddana wartość jest oznaczeniem języka, a nie zdaniem.
 *
 * Sprawdzian mówi wyłącznie o kształcie i jest celowo wąski. Kontrakt opisuje
 * pole `language` jako rozpoznany język, a rdzeń prosi model o samą nazwę albo
 * kod ISO 639-1 — oznaczenie języka to najwyżej kilka słów złożonych z liter,
 * ewentualnie z kodem w nawiasie („polski", „pl", „Polish (pl)"). Komunikat
 * kanału mieści się zwykle w czterdziestu znakach, więc sama długość go nie
 * odsiewa; odsiewa go zbiór znaków (ukośnik, kropka wypunktowania) i liczba słów.
 *
 * Klient nie orzeka, że wartość jest błędna: oddaje ją Operatorowi dosłownie
 * i zostawia mu ocenę — nie wpisuje tylko cudzego komunikatu do pola języka
 * i nie melduje go jako rozpoznania.
 */
const KSZTALT_JEZYKA = /^[\p{L}][\p{L}\p{M} ()\-,]*$/u;

function wygladaNaJezyk(wartosc: string): boolean {
  if (wartosc === '' || wartosc.length > 40) return false;
  if (!KSZTALT_JEZYKA.test(wartosc)) return false;
  return wartosc.split(/\s+/).length <= 3;
}
