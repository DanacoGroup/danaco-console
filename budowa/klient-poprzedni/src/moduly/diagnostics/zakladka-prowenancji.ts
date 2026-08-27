import type { ProvenanceCallListRequest } from '../../../../shared/contract';
import {
  pole,
  poleLiczbowe,
  przelacznik,
  przyciskAkcji,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza-braki';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import { utworzWierszWywolania, type WierszWywolania } from './prowenancja-wiersz';
import { utworzWydanieSladu, type ZakresWydania } from './prowenancja-wydanie';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzStanTresci } from './stany-okna';
import { cialoNarzedzia, objasnienieNarzedzia, pasekNarzedzia } from './zakladki-narzedzi';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Zakładka prowenancji Provenance Explorer jest jedynym miejscem jawności pracy
 * modeli: pokazuje wykaz wywołań, odczyt jednego śladu, ocenę odpowiedzi
 * i wydanie śladu do pliku.
 */
export interface ZakladkaProwenancji {
  /** Ciało zakładki. */
  element: HTMLElement;
  /** Odczyt rejestru wywołań kanału modelu. */
  odswiez(): void;
  /** Czy zakładka pytała już rdzeń — po tym wywołujący wie, czy odświeżać. */
  czytano(): boolean;
}

/**
 * Zdanie o niewpiętym źródle prowenancji: zakładka nie dostała źródła rodziny
 * wywołań z montażu modułu, więc nie ma czym zapytać rdzeń, i mówi to wprost
 * zamiast milczeć.
 */
const BEZ_ZRODLA =
  'Rejestr wywołań kanału modelu ma w kontrakcie rodzinę provenance.* i rdzeń ma dla niej ' +
  'uchwyty, ale ta zakładka nie dostała z montażu modułu Diagnostics źródła tej rodziny, ' +
  'więc nie ma czym zapytać. Brak jest PO STRONIE KLIENTA, w złożeniu modułu ' +
  '(moduly/diagnostics/indeks.ts przekazuje do Observability Tools źródła diagnostyki ' +
  'i obserwowalności; źródła prowenancji jeszcze nie).';

export function utworzZakladkeProwenancji(
  _zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
  prowenancja?: ZrodloProwenancji,
): ZakladkaProwenancji {
  const tresc = utworzStanTresci();
  const granica = poleLiczbowe('Górna granica liczby wywołań', 'domyślnie wg rdzenia');
  const proces = pole('Proces telemetrii', 'pole processId — wspólne z dziennikiem');
  const wzorzec = pole('Wzorzec w treści promptu i odpowiedzi', 'szukany napis');
  const tylkoWolne = przelacznik('Tylko wywołania powyżej progu opóźnienia');
  const odczytPrzycisk = przyciskAkcji('Odczytaj rejestr wywołań', 'dn-btn dn-btn--atrament');
  let czytano = false;
  let wiersze: readonly WierszWywolania[] = [];

  function zakresWydania(): ZakresWydania {
    const zakres = stan.zakres();
    return {
      wywolania: wiersze.filter((pozycja) => pozycja.wskazane()).map((pozycja) => pozycja.id()),
      ...(zakres.od === undefined ? {} : { od: zakres.od }),
      ...(zakres.do === undefined ? {} : { do: zakres.do }),
    };
  }

  const wydanie =
    prowenancja === undefined
      ? null
      : utworzWydanieSladu(prowenancja, zakresWydania, (zdanie, udane) =>
          tresc.potwierdzenie(zdanie, udane),
        );

  const element = cialoNarzedzia(
    objasnienieNarzedzia(
      'Rejestr wywołań kanału modelu: kanał, model, czas, tokeny, koszt, opóźnienie i stan ' +
        'każdego wywołania. Przy wierszu stoją trzy czynności — odczyt drzewa śladu wraz ' +
        'z treścią promptu i odpowiedzi, ocena trafności odpowiedzi nadawana przez Operatora ' +
        'oraz wskazanie wywołania do wydania śladu.',
    ),
    objasnienieNarzedzia(
      'O zapisie treści promptu i odpowiedzi oraz o redakcji danych wrażliwych rozstrzyga ' +
        'ustawienie rdzenia, nie to okno. Każde wywołanie mówi o sobie, czy treść ma zapisaną ' +
        'i czy została zredagowana — okno tego nie zakłada.',
    ),
    objasnienieNarzedzia(
      'Powtórzenia wywołania (provenance.call.replay) ta zakładka NIE ma. Komenda jest ' +
        'w kontrakcie; brak chwytu jest po stronie klienta i jest rozstrzygnięciem: powtórzenie ' +
        'wysyła prompt do modelu ponownie, więc wydaje pieniądze Operatora i nie stoi obok ' +
        'przycisków odczytu.',
    ),
    ...(prowenancja === undefined ? [objasnienieNarzedzia(BEZ_ZRODLA)] : []),
    pasekProwenancji(granica, proces, wzorzec, tylkoWolne, odczytPrzycisk),
    tresc.element,
    ...(wydanie === null ? [] : [wydanie.element]),
  );

  function odswiez(): void {
    czytano = true;
    wiersze = [];
    if (prowenancja === undefined) {
      tresc.blad(BEZ_ZRODLA);
      return;
    }
    const zakres = stan.zakres();
    const limit = granicaZKontrolki(granica);
    const idProcesu = proces.value.trim();
    const szukane = wzorzec.value.trim();
    const zadanie: ProvenanceCallListRequest = {
      ...(zakres.od === undefined ? {} : { fromTime: zakres.od }),
      ...(zakres.do === undefined ? {} : { toTime: zakres.do }),
      ...(limit === undefined ? {} : { limit }),
      ...(idProcesu === '' ? {} : { processId: idProcesu }),
      ...(szukane === '' ? {} : { pattern: szukane }),
      ...(tylkoWolne.checked ? { slowOnly: true } : {}),
    };
    tresc.ladowanie('Odczyt rejestru wywołań kanału modelu…');
    void prowenancja.wykazWywolan(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        tresc.blad(zdanieNiepowodzenia('odczytu rejestru wywołań modeli', wynik.powod), wynik.blad);
        return;
      }
      const oddane = wynik.wynik;
      if (oddane.calls.length === 0) {
        // Pustka jest stanem poprawnym instalacji bez dotychczasowych wywołań, nie niepowodzeniem odczytu.
        tresc.pusto(zdaniePustki(zadanie));
        return;
      }
      const zbudowane = oddane.calls.map((wywolanie) =>
        utworzWierszWywolania(prowenancja, wywolanie, (zdanie, udane) =>
          tresc.potwierdzenie(zdanie, udane),
        ),
      );
      wiersze = zbudowane;
      const lista = wykaz('Wywołania kanału modelu', 'dg-wykaz');
      lista.append(...zbudowane.map((pozycja) => pozycja.element));
      tresc.tresc().append(lista);
      tresc.potwierdzenie(zdaniePodsumowania(oddane), oddane.truncated !== true);
    });
  }

  odczytPrzycisk.addEventListener('click', odswiez);

  // Zakładka otwiera się pierwsza: puste miejsce treści musi mieć zdanie, nie wyglądać na brak wyniku.
  tresc.pusto('Rejestr wywołań nie został jeszcze odczytany — naciśnij „Odczytaj rejestr wywołań”.');

  return { element, odswiez, czytano: () => czytano };
}

/** Pasek zakładki: pola zawężające rejestr wywołań granicą, procesem, wzorcem treści i progiem opóźnienia, wraz z przyciskiem odczytu. */
function pasekProwenancji(
  granica: HTMLInputElement,
  proces: HTMLInputElement,
  wzorzec: HTMLInputElement,
  tylkoWolne: HTMLInputElement,
  odczyt: HTMLButtonElement,
): HTMLElement {
  return pasekNarzedzia(
    wiersz('Granica odczytu', granica, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Granica dotyczy liczby wywołań oddanych przez rdzeń. Ile ich warunki spełnia, ' +
        'a ile rdzeń oddał, mówi zdanie pod wykazem.',
    }),
    wiersz('Proces telemetrii', proces, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Pole processId jest wspólne wpisowi dziennika, telemetrii procesu i wywołaniu modelu — ' +
        'po nim wpis z Logs Viewera prowadzi do wywołania, które go wywołało.',
    }),
    wiersz('Wzorzec w treści', wzorzec, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Pole pattern; rdzeń szuka napisu w treści promptu i odpowiedzi. Wywołanie z zapisem ' +
        'treści wyłączonym nie da się tak odnaleźć — nie ma w czym szukać.',
    }),
    wiersz('Tylko wolne wywołania', tylkoWolne, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Pole slowOnly. Próg opóźnienia jest ustawieniem rdzenia, nie liczbą z tego okna, ' +
        'więc okno go nie pokazuje i nie zgaduje.',
    }),
    odczyt,
  );
}

/** Granica odczytu odczytana z pola formularza; wartość niepoprawna albo niedodatnia znaczy brak granicy narzuconej przez okno. */
function granicaZKontrolki(granica: HTMLInputElement): number | undefined {
  const wartosc = Number.parseInt(granica.value, 10);
  return Number.isInteger(wartosc) && wartosc > 0 ? wartosc : undefined;
}

/**
 * Zdanie pustki nazywa brak wywołań i mówi, czym rejestr był zawężony, bo przy
 * założonym filtrze prawdziwe jest zdanie węższe niż ogólne stwierdzenie braku
 * wywołań.
 */
function zdaniePustki(zadanie: ProvenanceCallListRequest): string {
  const zawezenia: string[] = [];
  if (zadanie.fromTime !== undefined || zadanie.toTime !== undefined) {
    zawezenia.push('zakres czasu modułu');
  }
  if (zadanie.processId !== undefined) zawezenia.push(`proces ${zadanie.processId}`);
  if (zadanie.pattern !== undefined) zawezenia.push(`wzorzec „${zadanie.pattern}"`);
  if (zadanie.slowOnly === true) zawezenia.push('tylko wywołania powyżej progu opóźnienia');
  if (zawezenia.length === 0) {
    return (
      'Rdzeń nie ma w rejestrze ani jednego wywołania kanału modelu — nie było jeszcze wywołań. ' +
      'To jest stan poprawny instalacji, w której model nie był jeszcze wołany, a nie ' +
      'niepowodzenie odczytu.'
    );
  }
  return (
    'Rdzeń odczytał rejestr i nie ma w nim wywołania spełniającego te warunki: ' +
    `${zawezenia.join('; ')}. Zdejmij zawężenie, żeby zobaczyć, czy wywołania są w ogóle.`
  );
}

/** Zdanie podsumowania odczytu podaje liczność zwróconych wywołań, całkowitą liczbę spełniających warunki i informację o przycięciu wyniku wprost. */
function zdaniePodsumowania(oddane: {
  calls: readonly unknown[];
  total?: number;
  truncated?: boolean;
}): string {
  const podstawa =
    `Rdzeń oddał ${String(oddane.calls.length)} wywołań` +
    `${oddane.total === undefined ? '' : ` z ${String(oddane.total)} spełniających warunki`}.`;
  return oddane.truncated === true
    ? `${podstawa} WYNIK PRZYCIĘTY granicą — rejestr jest tu niepełny, zawęź zakres czasu albo ` +
        'podnieś granicę odczytu.'
    : `${podstawa} To jest cały rejestr w tych warunkach.`;
}
