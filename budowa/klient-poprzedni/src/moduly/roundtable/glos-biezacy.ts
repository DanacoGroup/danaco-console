import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import type { GlosNaZywo } from './strumien-wypowiedzi';

/**
 * Nazwy stanów głosu uczestnika bieżącego, złożonego z obu dróg rdzenia — zdarzenia zapisu
 * i strumienia fragmentów — wspólne dla trzech widoków modułu.
 */
export const STAN_GLOSU = {
  mowi: 'mówi teraz',
  zerwany: 'strumień zerwany',
  domkniety: 'strumień domknięty',
  utrwalona: 'wypowiedź utrwalona',
  wyciszony: 'wyciszony w turze',
  brak: 'głosu jeszcze nie zabrał',
} as const;

/** Głos uczestnika złożony z obu dróg rdzenia: świeższa treść, jej pochodzenie, stan oraz tok rozumowania i powód zerwania. */
export interface GlosBiezacy {
  /** Treść do pokazania — świeższa z dwóch; pusta, gdy żadnej nie ma. */
  tekst: string;
  /** Czy pokazywana treść pochodzi ze strumienia, a nie z zapisu rdzenia. */
  zeStrumienia: boolean;
  /** Czy treść w tej chwili rośnie — strumień otwarty i bez błędu. */
  rosnie: boolean;
  /** Nazwa stanu, jedna z `STAN_GLOSU`. */
  stan: string;
  /** Tok rozumowania; osobno, bo rdzeń nie liczy go do treści wypowiedzi. */
  rozumowanie: string;
  /** Treść fragmentu błędu; pusta, gdy strumień się nie zerwał. */
  przyczyna: string;
}

/** Składa głos bieżący uczestnika z gromadzenia strumienia i ostatniej wypowiedzi utrwalonej, wybierając treść świeższą po długości. */
export function zlozGlosBiezacy(
  glos: GlosNaZywo | null,
  utrwalona: RoundtableStatement | null,
  wyciszony: boolean,
): GlosBiezacy {
  const zywa = glos === null ? '' : glos.tekst;
  const zapisana = utrwalona === null ? '' : utrwalona.content;
  const zeStrumienia = zywa !== '' && zywa.length >= zapisana.length;
  return {
    tekst: zeStrumienia ? zywa : zapisana,
    zeStrumienia,
    rosnie: glos !== null && !glos.domkniety && glos.przyczyna === '',
    stan: opisStanuGlosu(glos, utrwalona, wyciszony),
    rozumowanie: glos === null ? '' : glos.rozumowanie,
    przyczyna: glos === null ? '' : glos.przyczyna,
  };
}

/** Zdanie o stanie głosu — rozdziela sześć stanów, których nie wolno mylić, w kolejności ważności od awarii po ciszę. */
function opisStanuGlosu(
  glos: GlosNaZywo | null,
  utrwalona: RoundtableStatement | null,
  wyciszony: boolean,
): string {
  if (glos !== null && glos.przyczyna !== '') return STAN_GLOSU.zerwany;
  if (glos !== null && !glos.domkniety) return STAN_GLOSU.mowi;
  if (glos !== null) return STAN_GLOSU.domkniety;
  if (utrwalona !== null) return STAN_GLOSU.utrwalona;
  return wyciszony ? STAN_GLOSU.wyciszony : STAN_GLOSU.brak;
}

/** Zdanie pod tożsamością, gdy słów nie ma, nazywające powód: zerwanie, głos dopiero otwarty, wyciszenie albo cisza w turze. */
export function zdanieBezSlow(biezacy: GlosBiezacy, wyciszony: boolean): string {
  if (biezacy.przyczyna !== '') {
    return 'Rdzeń nie przypisał uczestnikowi ani jednego słowa — kanał odmówił przed odpowiedzią.';
  }
  if (biezacy.stan === STAN_GLOSU.utrwalona) {
    return (
      'Rdzeń otworzył głos tego uczestnika, ale nie przysłał jeszcze ani jednego słowa — ' +
      'ani zapisem, ani strumieniem.'
    );
  }
  if (wyciszony) return 'Uczestnik jest wyciszony w tej turze, więc rdzeń o głos go nie prosił.';
  return 'Ten uczestnik nie powiedział jeszcze w tej turze ani słowa widzianego przez to okno.';
}

/** Czy uczestnik jest wyciszony w tej turze — pole wyciszenia jest polem nieobowiązkowym kontraktu, więc brak znaczy nie. */
export function wyciszony(uczestnik: RoundtableParticipant | null): boolean {
  return uczestnik !== null && uczestnik.muted === true;
}

/** Ostatnia wypowiedź uczestnika w wykazie tury, wybrana po kolejności przyjścia, ponieważ interwencja moderatora może otworzyć głos powtórnie. */
export function ostatniaWypowiedz(
  wypowiedzi: readonly RoundtableStatement[],
  idUczestnika: string,
): RoundtableStatement | null {
  let znaleziona: RoundtableStatement | null = null;
  for (const wypowiedz of wypowiedzi) {
    if (wypowiedz.participantId === idUczestnika) znaleziona = wypowiedz;
  }
  return znaleziona;
}

/** Identyfikatory wypowiedzi będących ostatnimi głosami swoich mówców, żeby treść rosnąca doklejała się wyłącznie do wypowiedzi otwartej ostatnio. */
export function ostatnieGlosy(wypowiedzi: readonly RoundtableStatement[]): Set<string> {
  const ostatnia = new Map<string, string>();
  for (const wypowiedz of wypowiedzi) ostatnia.set(wypowiedz.participantId, wypowiedz.id);
  return new Set(ostatnia.values());
}
