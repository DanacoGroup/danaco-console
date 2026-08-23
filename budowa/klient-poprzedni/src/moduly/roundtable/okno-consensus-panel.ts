import {
  RoundtableTurnStatus,
  type RoundtableConsensus,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  poleTresci,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza-braki';
import type { StanDebaty } from './stan-debaty';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import { czynnosciStanowiska } from './czynnosci-arsenalu';
import { utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Consensus Panel — okno pomocnicze modułu Roundtable: zapis stanowiska
 * końcowego i eksport wniosków, odświeżany po zakończeniu tury albo całej
 * debaty.
 *
 * `roundtable.consensus.get` bierze `turnId` opcjonalny: podany zawęża do jednej
 * tury, pominięty obejmuje całą debatę. Okno nie zgaduje, którą chce Operator —
 * pole „Zakres stanowiska" rozstrzyga to wprost, a odpowiedź mówi, czego dotyczy
 * (`Dotyczy: tury #N` albo `Dotyczy: całej debaty`). Wybór „cała debata" pomija
 * pole `turnId` w żądaniu — nie wysyła go pustym napisem, bo to udawałoby
 * żądanie tury.
 *
 * Brak tury bieżącej przy zakresie „tura" i brak treści stanowiska po udanym
 * odczycie to dwa różne stany puste, oba poprawne. Odmowa odczytu jest trzecim,
 * osobnym stanem z kodem i treścią rdzenia — puste stanowisko przy odmowie
 * mówiłoby nieprawdę.
 *
 * Ponowny odczyt idzie na zamknięciu tury, nie przy każdej zmianie:
 * `stan-debaty.ts` subskrybuje `roundtable.debate.changed` i budzi
 * `stan.naZmiane(...)` przy każdym przyroście, więc druga subskrypcja zdarzenia
 * byłaby powieleniem. Okno łapie wyłącznie przejście tury
 * w `RoundtableTurnStatus.Closed`.
 *
 * Stanowisko niesie dziś w kontrakcie znacznie więcej, niż to okno pokazuje:
 * punkty zgody i sporne, zdania odrębne, poparcie ważone oraz części zapisu
 * decyzji. Okno rysuje treść i metrykę, bo tyle oddaje mu rdzeń — pozostałe
 * pola pozostają puste do czasu zbudowania ich obsługi, a nie dlatego, że
 * kontrakt ich nie ma.
 */
export interface OknoConsensusPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)`. */
  zamknij(): void;
}

type ZakresStanowiska = 'tura' | 'debata';

export function utworzOknoConsensusPanel(
  zrodlo: ZrodloRoundtable,
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  idOkna: string,
): OknoConsensusPanel {
  stan.ustawOkno(idOkna);

  const rama = utworzRameOkna({
    tytul: 'Consensus Panel',
    rola: 'pomocnicze',
    kod: 'consensus-panel',
    przeznaczenie: 'Zapis stanowiska końcowego debaty — tury bieżącej albo całości; eksport wniosków.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Pole redakcji stoi w oknie, bo stanowisko końcowe jest „edytowalnym punktem
  // wyjścia, nie wynikiem ostatecznym" (opracowanie, 2.7.2). Bez pola komenda
  // zapisu nie miałaby czego zapisać poza tym, co rdzeń i tak już złożył.
  const poleRedakcji = poleTresci(
    'Treść stanowiska po redakcji',
    5,
    'Stanowisko końcowe własnymi słowami — od zapisu należy do Ciebie.',
  );
  // Czynności stanowiska wołają komendy obszaru wprost: redakcję treści, zdanie
  // odrębne, wykaz wersji i przekazanie do modułu docelowego. Treść bierze się
  // z pola redakcji tego okna — rdzeń nie dostaje niczego, czego Operator nie
  // napisał.
  const powierzchnia = zlozPowierzchnieConsensusPanel(
    rama,
    tresc.element,
    czynnosciStanowiska(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => poleRedakcji.value,
      () => 'studio',
    ),
  );
  rama.cialo.prepend(
    wiersz('Redakcja stanowiska', poleRedakcji, {
      klasa: 'dr-wiersz',
      objasnienie:
        'Zapisana treść zostaje: od redakcji rdzeń nie składa stanowiska z zapisu tur na nowo.',
    }),
  );
  let ostatnieStanowisko: RoundtableConsensus | null = null;
  let ostatniOpisZakresu = '';

  function odczytaj(): void {
    const zakres = powierzchnia.zakres.value as ZakresStanowiska;
    if (zakres === 'debata') {
      zapytajStanowisko({ windowId: idOkna }, 'całej debaty');
      return;
    }
    const idTury = stan.tura();
    if (idTury === '') {
      ostatnieStanowisko = null;
      tresc.pusto('Debata nie ma jeszcze stanowiska — zamknij turę w Moderator Panelu.');
      return;
    }
    zapytajStanowisko({ windowId: idOkna, turnId: idTury }, opisTury(stan));
  }

  function zapytajStanowisko(
    zadanie: Parameters<ZrodloRoundtable['stanowisko']>[0],
    opisZakresu: string,
  ): void {
    ostatniOpisZakresu = opisZakresu;
    tresc.ladowanie(`Odczyt stanowiska ${opisZakresu}…`);
    void zrodlo.stanowisko(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Zdanie nie mówi „rdzeń odmówił": niepowodzenie bywa też odpowiedzią
        // bez pola obowiązkowego (`sprawdzKsztalt`). Powód niesie treść błędu.
        ostatnieStanowisko = null;
        tresc.blad(`Odczyt stanowiska ${opisZakresu} nie udał się.`, wynik.blad);
        return;
      }
      ostatnieStanowisko = wynik.wynik;
      rysujStanowisko(wynik.wynik, opisZakresu, tresc);
    });
  }

  function eksportuj(): void {
    const tresciowe = ostatnieStanowisko?.content ?? '';
    if (ostatnieStanowisko === null || tresciowe.trim() === '') {
      tresc.potwierdzenie('Nie ma czego wyeksportować — stanowisko jest jeszcze puste.', false);
      return;
    }
    pobierzPlik(
      `stanowisko-${ostatnieStanowisko.id}.md`,
      raportStanowiska(ostatnieStanowisko, ostatniOpisZakresu),
      'text/markdown',
    );
    tresc.potwierdzenie('Stanowisko pobrane jako plik Markdown.', true);
  }

  function eksportujZapisDecyzji(): void {
    const tresciowe = ostatnieStanowisko?.content ?? '';
    if (ostatnieStanowisko === null || tresciowe.trim() === '') {
      tresc.potwierdzenie(
        'Nie ma z czego złożyć zapisu decyzji — stanowisko jest jeszcze puste.',
        false,
      );
      return;
    }
    pobierzPlik(
      `zapis-decyzji-${ostatnieStanowisko.id}.md`,
      zapisDecyzji(ostatnieStanowisko, ostatniOpisZakresu),
      'text/markdown',
    );
    tresc.potwierdzenie(
      'Zapis decyzji pobrany. Części, których rdzeń nie wypełnił, są w pliku nazwane jako nieuzupełnione wraz z komendą, która je zapisuje — kontrakt ma pola wszystkich pięciu części.',
      true,
    );
  }

  podepnijAkcjeConsensusPanel(powierzchnia, { odczytaj, eksportuj, eksportujZapisDecyzji });

  let ostatniKluczTury = '';
  const odsubskrybuj = stan.naZmiane(() => {
    const definicja = stan.definicjaTury();
    const klucz = definicja === null ? '' : `${definicja.id}:${definicja.status}`;
    if (klucz === ostatniKluczTury) return;
    ostatniKluczTury = klucz;
    if (definicja !== null && definicja.status === RoundtableTurnStatus.Closed) odczytaj();
  });

  odczytaj();

  return { element: rama.element, odswiez: odczytaj, zamknij: odsubskrybuj };
}

/** Opis tury bieżącej dla zdań Operatora — numer i stan, gdy definicja jest znana. */
function opisTury(stan: StanDebaty): string {
  const definicja = stan.definicjaTury();
  if (definicja === null) return `tury ${stan.tura()}`;
  const otwartosc = definicja.status === RoundtableTurnStatus.Closed ? 'zamknięta' : 'otwarta';
  return `tury #${definicja.index} (${otwartosc})`;
}

/**
 * Rysuje stanowisko albo stan pustki, gdy odczyt się udał, ale treści jeszcze
 * nie ma — to jest stan poprawny (debata w toku), nie błąd.
 */
function rysujStanowisko(
  stanowisko: RoundtableConsensus,
  opisZakresu: string,
  tresc: StanTresci,
): void {
  if (stanowisko.content === undefined || stanowisko.content.trim() === '') {
    tresc.pusto(`Stanowisko ${opisZakresu} nie ma jeszcze treści — rdzeń go nie ustalił.`);
    return;
  }
  const wezel = document.createElement('div');
  wezel.className = 'dr-stanowisko';

  const naglowekZakresu = document.createElement('p');
  naglowekZakresu.className = 'dr-stanowisko__wersja';
  naglowekZakresu.textContent = `Dotyczy: ${opisZakresu}.`;

  const tresciowy = document.createElement('p');
  tresciowy.className = 'dr-stanowisko__tresc';
  tresciowy.textContent = stanowisko.content;

  const meta = document.createElement('p');
  meta.className = 'dr-stanowisko__wersja';
  meta.textContent = metaStanowiska(stanowisko);

  wezel.append(naglowekZakresu, tresciowy, meta);
  tresc.tresc().append(wezel);
}

/**
 * Zdanie z polami nośnymi `RoundtableConsensus`, których treść stanowiska nie
 * niesie.
 *
 * `updatedAt` jest w kontrakcie liczbą milisekund epoki; zero nie znaczy
 * „1 stycznia 1970", tylko „rdzeń chwili nie podał". Stopka pisze więc wprost
 * „nieznane" zamiast daty z początku epoki.
 */
function metaStanowiska(stanowisko: RoundtableConsensus): string {
  const wersja = stanowisko.version === undefined ? 'brak wersji' : `wersja ${stanowisko.version}`;
  const tury =
    stanowisko.turnIds === undefined || stanowisko.turnIds.length === 0
      ? 'tury nieznane'
      : `tury: ${stanowisko.turnIds.join(', ')}`;
  const chwila = stanowisko.updatedAt;
  const zaktualizowano =
    !Number.isFinite(chwila) || chwila <= 0
      ? 'czas aktualizacji nieznany (rdzeń nie podał chwili)'
      : `zaktualizowano ${new Date(chwila).toLocaleString('pl-PL')}`;
  return `${wersja} · ${zaktualizowano} · ${tury} · id ${stanowisko.id}`;
}

/** Treść pliku eksportu — wnioski w Markdown, gotowe do przekazania dalej. */
function raportStanowiska(stanowisko: RoundtableConsensus, opisZakresu: string): string {
  const wiersze = [
    '# Stanowisko debaty',
    '',
    `Dotyczy: ${opisZakresu}.`,
    '',
    stanowisko.content ?? '',
    '',
    `- Identyfikator: ${stanowisko.id}`,
    `- Okno: ${stanowisko.windowId}`,
    `- Wersja: ${stanowisko.version ?? 'brak wersji'}`,
    `- Tury objęte: ${(stanowisko.turnIds ?? []).join(', ') || 'nieznane'}`,
    `- Zaktualizowano: ${
      Number.isFinite(stanowisko.updatedAt) && stanowisko.updatedAt > 0
        ? new Date(stanowisko.updatedAt).toLocaleString('pl-PL')
        : 'nieznane — rdzeń nie podał chwili'
    }`,
    '',
  ];
  return wiersze.join('\n');
}

/**
 * Zapis decyzji w układzie ADR, złożony z tego, co stanowisko naprawdę niesie.
 *
 * Wszystkie pięć części ma dziś pole w `RoundtableConsensus`: kontekst, warianty
 * i konsekwencje zapisuje `roundtable.consensus.set`, a zdanie odrębne —
 * `roundtable.consensus.minority.set`. Część, której rdzeń nie wypełnił, zostaje
 * w pliku nazwana jako nieuzupełniona wraz z powodem. Zapis decyzji z domyślonym
 * kontekstem byłby dokumentem, który wygląda na protokół ustaleń, a niesie
 * zdania, których nikt w debacie nie powiedział.
 */
function zapisDecyzji(stanowisko: RoundtableConsensus, opisZakresu: string): string {
  const czesc = (tresc: string | undefined, czego: string): string =>
    tresc !== undefined && tresc.trim() !== ''
      ? tresc
      : `_Nieuzupełnione: rdzeń nie podał ${czego}. Pole jest w strukturze RoundtableConsensus, a zapisuje je roundtable.consensus.set._`;
  const mniejszosc =
    stanowisko.minority !== undefined && stanowisko.minority.length > 0
      ? stanowisko.minority
          .map((zdanie) => `- ${zdanie.participantId}: ${zdanie.content}`)
          .join('\n')
      : '_Nieuzupełnione: rdzeń nie podał zdania odrębnego. Pole minority jest w strukturze RoundtableConsensus, a zapisuje je roundtable.consensus.minority.set._';
  return [
    '# Zapis decyzji',
    '',
    `Dotyczy: ${opisZakresu}.`,
    '',
    '## Kontekst',
    '',
    czesc(stanowisko.context, 'kontekstu decyzji'),
    '',
    '## Rozważane warianty',
    '',
    czesc(stanowisko.options, 'wariantów rozważanych w debacie'),
    '',
    '## Decyzja',
    '',
    stanowisko.content ?? '',
    '',
    '## Konsekwencje',
    '',
    czesc(stanowisko.consequences, 'konsekwencji decyzji'),
    '',
    '## Ryzyko mniejszości',
    '',
    mniejszosc,
    '',
    '## Metryka stanowiska',
    '',
    `- Identyfikator: ${stanowisko.id}`,
    `- Okno debaty: ${stanowisko.windowId}`,
    `- Wersja nadana przez rdzeń: ${stanowisko.version ?? 'brak wersji'}`,
    `- Tury objęte: ${(stanowisko.turnIds ?? []).join(', ') || 'nieznane'}`,
    '',
  ].join('\n');
}

/** Kontrolki paska akcji i pola zakresu okna. */
interface PowierzchniaConsensusPanel {
  zakres: HTMLSelectElement;
  odczytajPrzycisk: HTMLButtonElement;
  eksport: HTMLButtonElement;
  eksportZapisu: HTMLButtonElement;
}

/**
 * Składa pasek akcji ramy: odczyt stanowiska i dwie postaci wydania — wnioski
 * oraz zapis decyzji. Czynności zmieniające stanowisko stoją w zestawie warstwy
 * trzeciej: ich komendy są w kontrakcie (`roundtable.consensus.set` zapisuje
 * treść wraz z częściami zapisu decyzji), ale okno wywołuje dziś wyłącznie
 * odczyt.
 */
function zlozAkcjeConsensusPanel(gospodarz: HTMLElement): {
  odczytajPrzycisk: HTMLButtonElement;
  eksport: HTMLButtonElement;
  eksportZapisu: HTMLButtonElement;
} {
  const odczytajPrzycisk = przycisk('Odczytaj stanowisko', 'dn-btn dn-btn--atrament');
  const eksport = przycisk('Eksportuj wnioski');
  const eksportZapisu = przycisk('Eksportuj zapis decyzji');

  gospodarz.append(odczytajPrzycisk, eksport, eksportZapisu);
  return { odczytajPrzycisk, eksport, eksportZapisu };
}

/** Zestaw akcji warstwy trzeciej — redakcja i wydanie jeszcze niezbudowane. */
function zlozZestawStanowiska(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Zestaw operacji na stanowisku', czynnosci);
}

/** Składa pole zakresu, pasek akcji, warstwy i ciało ramy. */
function zlozPowierzchnieConsensusPanel(
  rama: { akcje: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  czynnosci: HTMLButtonElement[],
): PowierzchniaConsensusPanel {
  const zakres = wybor('Zakres stanowiska', [
    ['tura', 'Bieżąca tura'],
    ['debata', 'Cała debata'],
  ]);
  const akcje = zlozAkcjeConsensusPanel(rama.akcje);
  rama.cialo.append(
    wiersz('Zakres stanowiska', zakres, {
      klasa: 'dr-wiersz',
      objasnienie: 'Cała debata pomija turnId w żądaniu, nie wysyła go pustym napisem.',
    }),
    zlozZestawStanowiska(czynnosci),
    stanTresci,
    utworzWykazFunkcji('consensus-panel'),
  );
  return { zakres, ...akcje };
}

/** Podpina pasek akcji i zmianę zakresu do czynności okna. */
function podepnijAkcjeConsensusPanel(
  powierzchnia: PowierzchniaConsensusPanel,
  obsluga: { odczytaj: () => void; eksportuj: () => void; eksportujZapisDecyzji: () => void },
): void {
  powierzchnia.odczytajPrzycisk.addEventListener('click', obsluga.odczytaj);
  powierzchnia.zakres.addEventListener('change', obsluga.odczytaj);
  powierzchnia.eksport.addEventListener('click', obsluga.eksportuj);
  powierzchnia.eksportZapisu.addEventListener('click', obsluga.eksportujZapisDecyzji);
}
