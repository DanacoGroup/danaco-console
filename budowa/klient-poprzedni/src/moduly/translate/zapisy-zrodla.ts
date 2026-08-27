import type { TranslateSourceSetRequest } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import type { WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { BRAK_OKNA } from './etykiety-translate';
import { zdanieOdmowyModelu } from './odmowa-translate';
import type { StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { rozbieznoscOdpowiedzi } from './zgodnosc-odpowiedzi';

/**
 * Trzy ścieżki zapisu Source Panel są oddzielone od widoku, który je wyzwala; każda ścieżka
 * kończy się zdaniem w wierszu odpowiedzi, żadna ciszą.
 */
export interface PolaZrodla {
  tekst: HTMLTextAreaElement;
  jezyk: HTMLInputElement;
  ponownaSegmentacja: HTMLInputElement;
}

export interface ZapisyZrodla {
  /** Zapis tekstu źródłowego wyzwala aktualizację wszystkich paneli tłumaczenia. */
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

/** Komenda zapisu tekstu źródłowego niesie tekst wraz z jego językiem, wyzwalając aktualizację wszystkich paneli tłumaczenia. */
async function zapiszZrodlo(
  stan: StanTranslate,
  pola: PolaZrodla,
  okno: StanOkna,
  odpowiedz: WierszOdpowiedzi,
): Promise<void> {
  const tekst = pola.tekst.value;
  if (tekst.trim() === '') {
    // Zdanie mówi o oknie, bo to okno odmawia, nie rdzeń: białe znaki same nie są zapisem tekstu.
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

  // Stan ładowania idzie na okno, nie tylko w wiersz odpowiedzi, treść okna zostaje widoczna.
  okno.ladowanie('Rdzeń przyjmuje tekst źródłowy i dzieli go na segmenty.');
  odpowiedz.pokaz('Zapisywanie tekstu źródłowego…', true);
  const wynik = await stan.zrodlo.ustaw(zadanie);
  if (!wynik.udany || wynik.wynik === undefined) {
    const zdanie = opisOdmowy('Zapis tekstu źródłowego', wynik.blad?.code, wynik.blad?.message);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  // Tekst wchodzi do stanu z żądania, bo odpowiedź zapisu go nie odsyła; zdanie mówi o języku.
  stan.wchlonZrodlo(tekst, wynik.wynik.sourceLanguage, []);
  stan.wchlonPanele(wynik.wynik.panels);

  // Język zamówiony porównuje się z oddanym tylko wtedy, gdy operator język wskazał w żądaniu.
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

/** Komenda ponownego podziału tekstu źródłowego na segmenty zwraca liczbę segmentów, która bierze się z odpowiedzi rdzenia. */
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
  // Faza schodzi przed ogłoszeniem stanu, bo przerysowanie okna nie zdejmuje fazy trwającej.
  okno.gotowe();
  stan.wchlonSegmenty(wynik.wynik.segments);
  odpowiedz.pokaz(`Rdzeń podzielił tekst na ${wynik.wynik.segments.length} segmentów.`, true);
}

/** Komenda rozpoznania języka źródłowego działa modelem, a wynik przechodzi przez odsiew kształtu przed wpisaniem do pola. */
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
    // Rozpoznanie języka rdzeń wykonuje modelem, więc odmowa z braku kanału niesie zdanie o braku.
    const zdanie = zdanieOdmowyModelu('Rozpoznanie języka', wynik.blad);
    odpowiedz.pokaz(zdanie, false);
    okno.blad(zdanie);
    return;
  }
  const oddany = wynik.wynik.language.trim();
  if (!wygladaNaJezyk(oddany)) {
    // Odpowiedź nieużyteczna jest odmową, nie pustką: pas stanu nie wraca wtedy do stanu gotowe.
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
 * Sprawdzian oznaczenia języka mówi wyłącznie o kształcie wartości i jest celowo wąski, bo
 * oznaczenie języka to najwyżej kilka słów złożonych z liter.
 */
const KSZTALT_JEZYKA = /^[\p{L}][\p{L}\p{M} ()\-,]*$/u;

function wygladaNaJezyk(wartosc: string): boolean {
  if (wartosc === '' || wartosc.length > 40) return false;
  if (!KSZTALT_JEZYKA.test(wartosc)) return false;
  return wartosc.split(/\s+/).length <= 3;
}
