import type { ResearchReport, ResearchReportSection } from '../../../../shared/contract';

/**
 * Zdanie o skutku złożenia raportu — nazwane raportem, który wrócił, a nie
 * żądaniem, które okno wysłało.
 *
 * Wydzielone z `czynnosci-raportu.ts`, bo tamten plik prowadzi czynność, a to
 * jest przekład odpowiedzi `research.report.build` na zdanie dla Operatora —
 * dwie rzeczy zmieniane z różnych powodów. Wzór rodziny:
 * `assistant/skutek-sterowania.ts`.
 *
 * Znaku „to nie model" w odpowiedzi nie ma. Sekcja niesie komplet pól
 * kontraktu i nic ponad to (`id, title, content, findingIds, order`), a rdzeń
 * zbiera z kanału same fragmenty `text` (`core/adapter_modul_badania.go`,
 * `zapytajModel`) — komunikat procesu kanału i zdanie modelu docierają więc tą
 * samą drogą, jako ta sama treść. Rozpoznanie po napisie jest zakazane, bo
 * napis się zmienia: okno nie orzeka, że streszczenie jest komunikatem błędu,
 * tylko mówi, czego nie wie, i każe przeczytać treść przed wydaniem raportu.
 *
 * Znak, który odróżnia naprawdę: sekcja, której identyfikatora okno nie
 * wysłało, jest sekcją napisaną przez rdzeń — kryterium bierze się
 * z odpowiedzi, a nie z brzmienia tytułu. Druga rozbieżność widoczna
 * z odpowiedzi: przy podanym polu `sections` rdzeń odkłada `findingIds` na
 * bok, więc zaznaczenie ustaleń przepada bez śladu — i to okno ogłasza
 * odmową, nie przemilcza.
 */
export interface SkutekZlozenia {
  zdanie: string;
  udany: boolean;
}

export function opisZlozenia(
  raport: ResearchReport,
  zamowioneSekcje: readonly ResearchReportSection[],
  zamowioneUstalenia: readonly string[],
): SkutekZlozenia {
  const oddane = raport.sections ?? [];
  const wstep = `Rdzeń złożył ${nazwaRaportu(raport)} — sekcji ${String(oddane.length)}.`;

  const czesci = [
    skutekSekcji(oddane, zamowioneSekcje),
    skutekUstalen(oddane, zamowioneUstalenia, zamowioneSekcje.length > 0),
    skutekAutorstwa(oddane, zamowioneSekcje),
  ].filter((czesc): czesc is SkutekZlozenia => czesc !== null);

  if (czesci.length === 0) {
    return {
      zdanie: `${wstep} Ani sekcji redakcji, ani zaznaczonych ustaleń nie zamówiono, więc nie ma czego porównywać.`,
      udany: true,
    };
  }
  return {
    zdanie: [wstep, ...czesci.map((czesc) => czesc.zdanie)].join(' '),
    udany: czesci.every((czesc) => czesc.udany),
  };
}

/**
 * Tytuł nazwany tak, jak go rdzeń zapisał.
 *
 * Tytuł pusty nazywamy pustym. Cudzysłów bez treści wygląda na usterkę
 * wypisywania, a jest wiernym oddaniem tego, co wróciło: pole `title` puste nie
 * idzie do rdzenia, więc raport został bez nazwy. Klient jej nie dopowiada za
 * Operatora.
 */
function nazwaRaportu(raport: ResearchReport): string {
  return raport.title.trim() === '' ? 'raport bez tytułu' : `raport „${raport.title}"`;
}

/** Sekcje redakcji sprawdzone w raporcie, który wrócił — nie w tym, co wysłano. */
function skutekSekcji(
  oddane: readonly ResearchReportSection[],
  zamowione: readonly ResearchReportSection[],
): SkutekZlozenia | null {
  if (zamowione.length === 0) return null;
  const kodyOddanych = new Set(oddane.map((sekcja) => sekcja.id));
  const brakujace = zamowione.filter((sekcja) => !kodyOddanych.has(sekcja.id));
  if (brakujace.length > 0) {
    return {
      zdanie:
        'Rdzeń przyjął wywołanie, ale sekcji redakcji, których NIE MA w oddanym raporcie, jest ' +
        `${String(brakujace.length)} z ${String(zamowione.length)} — zapis tych sekcji się nie odbył.`,
      udany: false,
    };
  }
  return {
    zdanie: `Wszystkie sekcje redakcji (${String(zamowione.length)}) wróciły z rdzenia.`,
    udany: true,
  };
}

/**
 * Zaznaczone ustalenia sprawdzone w wiązaniach sekcji, które wróciły.
 *
 * Rdzeń oddaje `sections[].findingIds`, więc sprawdzenie jest darmowe:
 * ustalenie zaznaczone, którego nie wiąże żadna oddana sekcja, do raportu nie
 * weszło — choćby żądanie je niosło.
 */
function skutekUstalen(
  oddane: readonly ResearchReportSection[],
  zamowione: readonly string[],
  bylaRedakcja: boolean,
): SkutekZlozenia | null {
  if (zamowione.length === 0) return null;
  const zwiazane = new Set(oddane.flatMap((sekcja) => sekcja.findingIds ?? []));
  const pominiete = zamowione.filter((kod) => !zwiazane.has(kod));
  if (pominiete.length === 0) {
    return {
      zdanie: `Wszystkie zaznaczone ustalenia (${String(zamowione.length)}) wiąże któraś z oddanych sekcji.`,
      udany: true,
    };
  }
  return {
    zdanie:
      'Zaznaczonych ustaleń, których nie wiąże ŻADNA oddana sekcja, jest ' +
      `${String(pominiete.length)} z ${String(zamowione.length)} — do raportu nie weszły. ` +
      (bylaRedakcja
        ? 'W research.report.build sekcje podane wprost mają pierwszeństwo nad findingIds. ' +
          'Opróżnij wiersze redakcji i złóż ponownie, jeśli raport ma powstać ze streszczenia ustaleń.'
        : 'Żądanie nie niosło sekcji redakcji, więc pominięcia nie tłumaczy pierwszeństwo redakcji.'),
    udany: false,
  };
}

/**
 * Sekcje, których treści Operator nie napisał — poznane po identyfikatorze,
 * którego okno nie wysyłało, a nie po brzmieniu tytułu.
 *
 * Tu stoi jedyne uczciwe zdanie o streszczeniu: rdzeń wziął treść z kanału
 * modelu, a odpowiedź nie niesie znaku, czy jest to redakcja modelu, czy
 * komunikat procesu kanału (patrz nagłówek pliku). Okno nie zgaduje po napisie
 * i nie potwierdza autorstwa, którego rdzeń nie oddał.
 *
 * Granica tego znaku: zdanie pada wyłącznie przy złożeniu, w którym sekcja
 * przyszła z rdzenia po raz pierwszy. Kolejne złożenie wysyła ją już
 * z redakcji kreatora — pod tym samym identyfikatorem, bo `wczytaj` przejmuje
 * sekcje potwierdzone — więc kryterium przestaje ją wskazywać i słusznie:
 * Operator miał ją wtedy przed sobą w polach edycji.
 */
function skutekAutorstwa(
  oddane: readonly ResearchReportSection[],
  zamowione: readonly ResearchReportSection[],
): SkutekZlozenia | null {
  const kodyZamowionych = new Set(zamowione.map((sekcja) => sekcja.id));
  const odRdzenia = oddane.filter((sekcja) => !kodyZamowionych.has(sekcja.id));
  if (odRdzenia.length === 0) return null;
  return {
    zdanie:
      `Sekcji, których treści nie napisał Operator: ${String(odRdzenia.length)} — przyszły z rdzenia ` +
      'kanałem modelu. Klient nie ma w odpowiedzi znaku, czy jest to redakcja modelu, czy komunikat ' +
      'jego procesu: rdzeń oddaje jedno i drugie jako zwykłą treść sekcji ze statusem „ok" (zmierzone). ' +
      'Przeczytaj je w podglądzie, zanim wydasz raport.',
    udany: true,
  };
}
