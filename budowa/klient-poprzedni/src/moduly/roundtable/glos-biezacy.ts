import type { RoundtableParticipant, RoundtableStatement } from '../../../../shared/contract';
import type { GlosNaZywo } from './strumien-wypowiedzi';

/**
 * Głos uczestnika w chwili bieżącej — wspólne dla trzech widoków modułu ustalenie,
 * którą treść pokazać i jak ją nazwać.
 *
 * Rdzeń nadaje wypowiedź dwiema drogami: zdarzeniem `roundtable.debate.changed`
 * (rodzaju `created` z treścią pustą w chwili otwarcia głosu, rodzaju `updated`
 * z treścią pełną po domknięciu strumienia) oraz fragmentami `stream.chunk` przez
 * cały czas mówienia. Widok czytający wyłącznie zdarzenia pokazuje pustą wypowiedź
 * aż do domknięcia strumienia, dlatego obie drogi trzeba złożyć w jedną treść.
 *
 * Świeższą treść wybiera się po długości, bo fragment nie niesie znacznika czasu
 * (`StreamChunkEvent` go nie ma), a `RoundtableStatement.createdAt` bywa zerem:
 * rosnąca wygrywa, gdy jest niekrótsza; utrwalona wchodzi, gdy rosnącej nie ma
 * albo urwała się krótsza — strumień zerwany w połowie, a rdzeń zapisał całość.
 *
 * Moduł nie tworzy węzłów DOM i nie zna klas CSS: trzy widoki rysują inaczej
 * (kolumna, wykaz, panel), a treść i stan mają mieć to samo.
 */

/**
 * Nazwy stanów głosu — napisy stałe, bo trafiają do atrybutu `data-stan-glosu`,
 * w który arkusze `panel-debaty.css` i `debata.css` celują wprost. Wyliczenie
 * stoi tu, a nie w widoku, żeby zmiana napisu w jednym miejscu nie rozjechała
 * dwóch arkuszy i trzech okien.
 */
export const STAN_GLOSU = {
  mowi: 'mówi teraz',
  zerwany: 'strumień zerwany',
  domkniety: 'strumień domknięty',
  utrwalona: 'wypowiedź utrwalona',
  wyciszony: 'wyciszony w turze',
  brak: 'głosu jeszcze nie zabrał',
} as const;

/** Głos uczestnika złożony z obu dróg rdzenia. */
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

/**
 * Składa głos bieżący uczestnika.
 *
 * @param glos gromadzenie ze `stream.chunk`; `null`, gdy nic nie płynęło.
 * @param utrwalona ostatnia wypowiedź tego uczestnika zapisana przez rdzeń.
 * @param wyciszony czy uczestnik jest wyciszony w turze — rozdziela „milczy,
 *   bo go nie pytano" od „milczy, choć pytano".
 */
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

/**
 * Zdanie o stanie głosu — rozdziela sześć stanów, których nie wolno mylić.
 *
 * Kolejność sprawdzeń jest kolejnością ważności: awaria przed trwaniem, trwanie
 * przed domknięciem, cokolwiek ze strumienia przed samym zapisem, a dopiero na
 * końcu dwa rodzaje ciszy.
 */
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

/**
 * Zdanie pod tożsamością, gdy słów nie ma — podaje powód zamiast pustki.
 *
 * Wypowiedź utrwalona o treści pustej jest przypadkiem osobnym: rdzeń rozgłasza
 * `created` z pustą treścią w chwili otwarcia głosu, więc okno bez odbioru
 * strumienia widzi ten stan przez cały czas mówienia modelu. Zdanie mówi wtedy,
 * że głos jest otwarty, a nie że uczestnik milczy.
 */
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

/** Czy uczestnik jest wyciszony — `muted` jest polem nieobowiązkowym kontraktu. */
export function wyciszony(uczestnik: RoundtableParticipant | null): boolean {
  return uczestnik !== null && uczestnik.muted === true;
}

/**
 * Ostatnia wypowiedź uczestnika w wykazie tury; `null`, gdy żadnej nie było.
 *
 * Rdzeń dopisuje wypowiedź uczestnika raz na turę, ale interwencja moderatora
 * potrafi otworzyć mu głos powtórnie i wtedy rośnie ta ostatnia. Wykaz
 * `StanDebaty.wypowiedzi()` zachowuje kolejność przyjścia, więc ostatnie
 * trafienie jest tym, o które chodzi.
 */
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

/**
 * Identyfikatory wypowiedzi będących ostatnimi głosami swoich mówców.
 *
 * Głos rosnący jest własnością uczestnika, a wykaz Debate Panelu jest wykazem
 * wypowiedzi. Doklejenie treści rosnącej do każdej wypowiedzi tego samego mówcy
 * powtórzyłoby te same słowa tyle razy, ile razy odezwał się w turze; rosnąca
 * należy wyłącznie do wypowiedzi otwartej ostatnio.
 */
export function ostatnieGlosy(wypowiedzi: readonly RoundtableStatement[]): Set<string> {
  const ostatnia = new Map<string, string>();
  for (const wypowiedz of wypowiedzi) ostatnia.set(wypowiedz.participantId, wypowiedz.id);
  return new Set(ostatnia.values());
}
