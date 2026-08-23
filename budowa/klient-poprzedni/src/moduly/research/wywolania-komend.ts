import {
  Command,
  ResearchAnchorKind,
  ResearchAnnotationKind,
  ResearchAttachmentKind,
  ResearchBibliographyScope,
  ResearchBlockKind,
  ResearchCitationMode,
  ResearchDiscoveryMode,
  ResearchMonitorKind,
  ResearchSnowballDirection,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { AkcjaBadania } from './akcje-okien';
import { kodyZTekstu, zlaczKody } from './badanie-ksiazka-kodow';
import {
  czyUzbrojona,
  najnowszaAdnotacja,
  opisAdnotacji,
  rozbrojZdjecie,
  uzbrojZdjecie,
} from './badanie-zdjecie-adnotacji';
import type { StanBadania } from './stan-badania';

/**
 * Droga z okna do komend obszaru `research.*`, które nie mają własnego
 * formularza.
 *
 * Jedna odpowiedzialność: złożyć żądanie z tego, co okno już wie — okna badania,
 * zaznaczenia źródeł i ustaleń, wskazania lektury, bieżącego raportu — wywołać
 * komendę i powiedzieć, co z niej wyszło. Zdanie odpowiedzi mówi o MIERZONYM
 * skutku (ile pozycji, jaki plik, ile luk), a nie o tym, że wywołanie się
 * powiodło: „rdzeń oddał wynik" jest zdaniem, po którym Operator nadal nie wie,
 * czy coś się stało.
 *
 * Akcje panelu niosą w polu `kod` dokładnie nazwę komendy (`Command.*`, patrz
 * `akcje-okien.ts`), więc rozdzielnik niżej jest odwzorowaniem jeden do jednego
 * i nie ma w nim ani jednej nazwy pisanej z ręki.
 *
 * Czego tu nie ma: żądań, których nie da się złożyć bez tekstu od Operatora
 * (zapytanie wyszukiwania, treść notatki, uzasadnienie odrzucenia). Takie
 * pozycje mówią wprost, czego brakuje, zamiast wysyłać żądanie z polem pustym
 * i wracać odmową walidacji, z której nic nie wynika.
 */
export interface KontekstKomendy {
  stan: StanBadania;
  /** Tekst z pola okna, gdy czynność go wymaga; pusty, gdy okno pola nie ma. */
  tekst?: string;
}

/** Wynik wywołania w postaci, którą okno wypisuje w wierszu odpowiedzi. */
export interface WynikKomendy {
  udany: boolean;
  opis: string;
}

/** Czy kod akcji jest nazwą komendy obszaru, którą ten plik potrafi wywołać. */
export function czyKomendaBadania(kod: string): boolean {
  return ROZDZIELNIK.has(kod);
}

/**
 * Wykonuje komendę wskazaną kodem akcji.
 *
 * Wynik nieznanego kodu jest odmową nazwaną, nie ciszą: kod spoza rozdzielnika
 * znaczy, że akcja została opisana jako `komenda`, a wywołania jej nie dopisano
 * — i to jest usterka do naprawy tutaj, nie stan do przemilczenia.
 */
export async function wykonajKomendeBadania(
  kontekst: KontekstKomendy,
  akcja: AkcjaBadania,
): Promise<WynikKomendy> {
  const wykonanie = ROZDZIELNIK.get(akcja.kod);
  if (wykonanie === undefined) {
    return {
      udany: false,
      opis:
        `Akcja „${akcja.nazwa}" jest opisana jako komenda ${akcja.kod}, ale okno nie ma dla ` +
        'niej wywołania — to jest brak w kliencie, nie odmowa rdzenia.',
    };
  }
  if (kontekst.stan.idOkna() === '') {
    return {
      udany: false,
      opis:
        `Komenda ${akcja.kod} pracuje w oknie badania, którego rdzeń jeszcze nie wskazał ` +
        '(window.list nie oddał okna modułu Research).',
    };
  }
  return wykonanie(kontekst, akcja);
}

type Wykonanie = (kontekst: KontekstKomendy, akcja: AkcjaBadania) => Promise<WynikKomendy>;

/** Odmowa braku wskazania — mówi, czego okno nie ma, zamiast wysyłać puste pole. */
function brak(czego: string): WynikKomendy {
  return { udany: false, opis: czego };
}

/** Zdanie o skutku udanym. */
function skutek(opis: string): WynikKomendy {
  return { udany: true, opis };
}

/** Przekłada odmowę rdzenia na zdanie okna. */
function odmowa(nazwa: string, wynik: { blad?: unknown; nieznanyTyp?: string }): WynikKomendy {
  return {
    udany: false,
    opis: opisOdmowyBledu(
      nazwa,
      wynik.blad as Parameters<typeof opisOdmowyBledu>[1],
      wynik.nieznanyTyp,
    ),
  };
}

/** Pierwsze zaznaczone źródło albo materiał wskazany do lektury. */
function wskazaneZrodlo(stan: StanBadania): string {
  return stan.wybraneZrodla.wybrane()[0] ?? stan.lektura.wskazane();
}

/** Pierwsze zaznaczone ustalenie. */
function wskazaneUstalenie(stan: StanBadania): string {
  return stan.wybraneUstalenia.wybrane()[0] ?? '';
}

/** Identyfikator raportu bieżącego okna. */
function wskazanyRaport(stan: StanBadania): string {
  return stan.raport()?.id ?? '';
}

const ROZDZIELNIK = new Map<string, Wykonanie>([
  // ── źródła ────────────────────────────────────────────────────────────────
  [
    Command.ResearchSourceList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceList, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Odczyt wykazu źródeł', wynik);
      for (const zrodlo of wynik.wynik.sources) stan.wchlonZrodlo(zrodlo);
      return skutek(
        `Rdzeń oddał ${String(wynik.wynik.sources.length)} z ${String(wynik.wynik.total)} źródeł okna.`,
      );
    },
  ],
  [
    Command.ResearchSourceRemove,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz źródło, które ma zniknąć z katalogu.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceRemove, { sourceId: zrodlo });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Usunięcie źródła', wynik);
      return skutek(
        `Źródło zdjęte; odwołanie do niego straciło ${String(wynik.wynik.detachedFindings)} ustaleń.`,
      );
    },
  ],
  [
    Command.ResearchSourceDuplicates,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceDuplicates, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykrycie duplikatów', wynik);
      if (wynik.wynik.candidates.length === 0) {
        return skutek('Rdzeń nie znalazł pozycji powtórzonych w katalogu tego badania.');
      }
      const pierwszy = wynik.wynik.candidates[0];
      return skutek(
        `Par podejrzanych o powtórzenie: ${String(wynik.wynik.candidates.length)}; ` +
          `pierwsza zbieżna w ${String(pierwszy.similarity)}% na podstawie: ${pierwszy.basis}.`,
      );
    },
  ],
  [
    Command.ResearchSourceMerge,
    async ({ stan }) => {
      const wybrane = stan.wybraneZrodla.wybrane();
      if (wybrane.length < 2) {
        return brak('Scalenie potrzebuje co najmniej dwóch zaznaczonych źródeł: docelowego i powtórzonych.');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceMerge, {
        targetSourceId: wybrane[0],
        mergedSourceIds: [...wybrane.slice(1)],
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Scalenie źródeł', wynik);
      return skutek(
        `Scalono w „${wynik.wynik.source.title}"; przeniesionych powiązań z ustaleniami: ` +
          `${String(wynik.wynik.movedFindings)}.`,
      );
    },
  ],
  [
    Command.ResearchSourceTag,
    async ({ stan, tekst }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz źródło, któremu chcesz nadać etykiety.');
      const etykiety = (tekst ?? '')
        .split(',')
        .map((wpis) => wpis.trim())
        .filter((wpis) => wpis !== '');
      if (etykiety.length === 0) {
        return brak('Wpisz etykiety po przecinku — katalogowanie bez etykiet niczego nie zmienia.');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceTag, {
        sourceId: zrodlo,
        tags: etykiety,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Nadanie etykiet', wynik);
      stan.wchlonZrodlo(wynik.wynik.source);
      return skutek(`Źródło „${wynik.wynik.source.title}" ma teraz ${String(etykiety.length)} etykiet.`);
    },
  ],
  [
    Command.ResearchSourceUpdate,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz źródło, którego stan lektury ma się zmienić.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceUpdate, {
        sourceId: zrodlo,
        readingState: 'read',
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zmiana źródła', wynik);
      stan.wchlonZrodlo(wynik.wynik.source);
      return skutek(`Źródło „${wynik.wynik.source.title}" oznaczone jako przeczytane.`);
    },
  ],
  [
    Command.ResearchSourceAttachmentList,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz źródło, którego załączniki chcesz obejrzeć.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceAttachmentList, {
        sourceId: zrodlo,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz załączników', wynik);
      const braki = wynik.wynik.missingKinds ?? [];
      return skutek(
        `Załączników: ${String(wynik.wynik.attachments.length)}` +
          (braki.length > 0 ? `; brakuje: ${braki.join(', ')}.` : '; komplet.'),
      );
    },
  ],
  [
    Command.ResearchSourceAttachmentAdd,
    async ({ stan, tekst }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz źródło, do którego ma dojść załącznik.');
      const sciezka = (tekst ?? '').trim();
      const dokument = stan.zrodla().find((wpis) => wpis.id === zrodlo)?.libraryFileId ?? '';
      // Dwie drogi materiału i żadnej trzeciej: plik z urządzenia Operatora albo
      // dokument repozytorium, który źródło już wskazuje. Żądanie bez żadnej
      // z nich wróciłoby odmową walidacji, z której nic nie wynika.
      if (sciezka === '' && dokument === '') {
        return brak(
          'Wskaż plik pełnego tekstu w polu tytułu formularza źródła albo skataloguj źródło ' +
            'z dokumentem repozytorium — załącznik musi mieć treść, a okno jej nie wymyśla.',
        );
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceAttachmentAdd, {
        sourceId: zrodlo,
        kind: ResearchAttachmentKind.Fulltext,
        ...(sciezka === '' ? { libraryFileId: dokument } : { sourcePath: sciezka }),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Dołączenie załącznika', wynik);
      const zalacznik = wynik.wynik.attachment;
      return skutek(
        `Załącznik pełnego tekstu dołączony do źródła (${sciezka === '' ? 'dokument repozytorium' : sciezka})` +
          (zalacznik.sizeBytes === undefined
            ? '; rozmiaru rdzeń nie podał.'
            : `; bajtów ${String(zalacznik.sizeBytes)}.`),
      );
    },
  ],
  [
    Command.ResearchSourceImport,
    async ({ stan, tekst }) => {
      const sciezka = (tekst ?? '').trim();
      if (sciezka === '') {
        return brak('Wskaż plik bibliografii (BibTeX, RIS, CSL-JSON, EndNote XML albo CSV).');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceImport, {
        windowId: stan.idOkna(),
        format: 'bibtex',
        sourcePath: sciezka,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Import bibliografii', wynik);
      for (const zrodlo of wynik.wynik.sources) stan.wchlonZrodlo(zrodlo);
      return skutek(
        `Wczytanych pozycji: ${String(wynik.wynik.imported)}, pominiętych: ${String(wynik.wynik.skipped)}.`,
      );
    },
  ],
  [
    Command.ResearchSourceCapture,
    async ({ stan, tekst }) => {
      const adres = (tekst ?? '').trim();
      if (adres === '') return brak('Wpisz adres strony do pozyskania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceCapture, {
        windowId: stan.idOkna(),
        url: adres,
        mode: 'both',
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Pozyskanie strony', wynik);
      stan.wchlonZrodlo(wynik.wynik.source);
      return skutek(
        `Strona wczytana jako źródło „${wynik.wynik.source.title}"` +
          (wynik.wynik.snapshotAttachmentId !== undefined ? ' wraz z migawką.' : '.'),
      );
    },
  ],
  [
    Command.ResearchSourceTranscribe,
    async ({ stan, tekst }) => {
      const sciezka = (tekst ?? '').trim();
      if (sciezka === '') return brak('Wskaż nagranie do przepisania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceTranscribe, {
        windowId: stan.idOkna(),
        sourcePath: sciezka,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Transkrypcja nagrania', wynik);
      stan.wchlonZrodlo(wynik.wynik.source);
      return skutek(`Transkrypt zapisany jako źródło; odcinków: ${String(wynik.wynik.segmentCount)}.`);
    },
  ],
  [
    Command.ResearchSourceResolve,
    async ({ stan, tekst }) => {
      const identyfikator = (tekst ?? '').trim();
      if (identyfikator === '') return brak('Wklej identyfikator DOI, ISBN, PMID albo arXiv.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceResolve, {
        identifier: identyfikator,
      });
      if (!wynik.udany || wynik.wynik === undefined) {
        return odmowa('Rozstrzygnięcie identyfikatora', wynik);
      }
      return skutek(`Pozycja rozpoznana: „${wynik.wynik.result.title}" (${wynik.wynik.result.provider}).`);
    },
  ],

  // ── odkrywanie i monitory ────────────────────────────────────────────────
  [
    Command.ResearchDiscoverySearch,
    async ({ stan, tekst }) => {
      const zapytanie = (tekst ?? '').trim();
      if (zapytanie === '') return brak('Wpisz zapytanie wyszukiwawcze.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchDiscoverySearch, {
        windowId: stan.idOkna(),
        query: zapytanie,
        mode: ResearchDiscoveryMode.Scholarly,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wyszukiwanie źródeł', wynik);
      const zawiedli = wynik.wynik.providersFailed ?? [];
      return skutek(
        `Trafień: ${String(wynik.wynik.results.length)} od dostawców ${wynik.wynik.providersUsed.join(', ')}` +
          (zawiedli.length > 0 ? `; nie odpowiedzieli: ${zawiedli.join(', ')}.` : '.'),
      );
    },
  ],
  [
    Command.ResearchDiscoveryAssist,
    async ({ stan, tekst }) => {
      const pytanie = (tekst ?? '').trim() === '' ? stan.zakres() : (tekst ?? '').trim();
      if (pytanie === '') return brak('Wpisz pytanie badawcze albo zapisz zakres badania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchDiscoveryAssist, { question: pytanie });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Query Assistant', wynik);
      return skutek(`Zaproponowanych zapytań: ${String(wynik.wynik.queries.length)}.`);
    },
  ],
  [
    Command.ResearchDiscoverySnowball,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Zaznacz pozycję, której cytowania chcesz rozwinąć.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchDiscoverySnowball, {
        sourceId: zrodlo,
        direction: ResearchSnowballDirection.Both,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Snowballing', wynik);
      const wstecz = wynik.wynik.referenced?.length ?? 0;
      const wprzod = wynik.wynik.citing?.length ?? 0;
      return skutek(`Prac cytowanych: ${String(wstecz)}, cytujących: ${String(wprzod)}.`);
    },
  ],
  [
    Command.ResearchDiscoveryReject,
    async ({ stan, tekst }) => {
      const powod = (tekst ?? '').trim();
      if (powod === '') return brak('Wpisz uzasadnienie odrzucenia — bez niego przesiewu nie da się udokumentować.');
      const klucze = stan.wybraneZrodla.wybrane();
      if (klucze.length === 0) return brak('Zaznacz pozycje wyniku, które mają zostać odrzucone.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchDiscoveryReject, {
        windowId: stan.idOkna(),
        resultKeys: [...klucze],
        reason: powod,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Odrzucenie pozycji', wynik);
      return skutek(
        `Odrzuconych: ${String(wynik.wynik.rejected)}; przesiew PRISMA: włączonych ` +
          `${String(wynik.wynik.counts.included)}, wyłączonych ${String(wynik.wynik.counts.excluded)}.`,
      );
    },
  ],
  [
    Command.ResearchMonitorList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchMonitorList, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz monitorów', wynik);
      return skutek(
        `Monitorów: ${String(wynik.wynik.monitors.length)}; nowych pozycji w skrzynce: ` +
          `${String(wynik.wynik.pendingTotal)}.`,
      );
    },
  ],
  [
    Command.ResearchMonitorSet,
    async ({ stan, tekst }) => {
      const zapytanie = (tekst ?? '').trim();
      if (zapytanie === '') return brak('Wpisz temat do monitorowania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchMonitorSet, {
        windowId: stan.idOkna(),
        kind: ResearchMonitorKind.Topic,
        query: zapytanie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis monitora', wynik);
      return skutek(`Monitor tematu zapisany pod identyfikatorem ${wynik.wynik.monitor.id}.`);
    },
  ],
  [
    Command.ResearchMonitorRefresh,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchMonitorRefresh, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Odświeżenie monitorów', wynik);
      const zawiodly = wynik.wynik.failedMonitorIds ?? [];
      return skutek(
        `Nowych pozycji: ${String(wynik.wynik.results.length)}` +
          (zawiodly.length > 0 ? `; monitorów bez odpowiedzi: ${String(zawiodly.length)}.` : '.'),
      );
    },
  ],
  [
    Command.ResearchBatchImport,
    async ({ stan, tekst }) => {
      const adresy = (tekst ?? '')
        .split(/\s+/u)
        .map((wpis) => wpis.trim())
        .filter((wpis) => wpis !== '');
      if (adresy.length === 0) return brak('Wklej listę adresów do pozyskania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchBatchImport, {
        windowId: stan.idOkna(),
        urls: adresy,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Import wsadowy', wynik);
      return skutek(
        `Partia ${wynik.wynik.queueItemId}: przyjętych ${String(wynik.wynik.accepted)}, ` +
          `odrzuconych ${String(wynik.wynik.rejected)}.`,
      );
    },
  ],

  // ── lektura i ekstrakcja ─────────────────────────────────────────────────
  [
    Command.ResearchReadingOpen,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż materiał do lektury.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReadingOpen, { sourceId: zrodlo });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Otwarcie lektury', wynik);
      const tresc = wynik.wynik.content;
      return skutek(
        `Wczytana strona ${String(tresc.page ?? 1)} z ${String(tresc.pageCount ?? 1)}; znaków: ` +
          `${String(tresc.text?.length ?? 0)}.`,
      );
    },
  ],
  [
    Command.ResearchAnnotationAdd,
    async ({ stan, tekst }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż materiał, przy którym ma stanąć adnotacja.');
      const cytat = (tekst ?? '').trim();
      if (cytat === '') return brak('Zaznacz fragment treści — podświetlenie bez cytatu nie jest wypisem.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchAnnotationAdd, {
        sourceId: zrodlo,
        kind: ResearchAnnotationKind.Highlight,
        anchor: { kind: ResearchAnchorKind.Page },
        quote: cytat,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis adnotacji', wynik);
      return skutek(`Podświetlenie zapisane pod identyfikatorem ${wynik.wynik.annotation.id}.`);
    },
  ],
  [
    Command.ResearchAnnotationList,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      const wynik = await stan.zrodlo.wywolaj(
        Command.ResearchAnnotationList,
        zrodlo === '' ? { windowId: stan.idOkna() } : { sourceId: zrodlo },
      );
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz adnotacji', wynik);
      return skutek(`Adnotacji: ${String(wynik.wynik.annotations.length)}.`);
    },
  ],
  [
    Command.ResearchAnnotationRemove,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż materiał, przy którym stoi adnotacja do zdjęcia.');
      const wykaz = await stan.zrodlo.wywolaj(Command.ResearchAnnotationList, { sourceId: zrodlo });
      if (!wykaz.udany || wykaz.wynik === undefined) return odmowa('Wykaz adnotacji', wykaz);
      const adnotacja = najnowszaAdnotacja(wykaz.wynik.annotations);
      if (adnotacja === null) {
        rozbrojZdjecie();
        return brak('Przy tym materiale nie ma ani jednej adnotacji — nie ma czego zdjąć.');
      }
      // Ostrzeżenie idzie PRZED wykonaniem: pierwsze naciśnięcie nazywa
      // adnotacja i mówi, że zdjęcie jest nieodwracalne, dopiero drugie zdejmuje.
      if (!czyUzbrojona(adnotacja.id)) {
        return { udany: false, opis: uzbrojZdjecie(adnotacja) };
      }
      const opis = opisAdnotacji(adnotacja);
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchAnnotationRemove, {
        annotationId: adnotacja.id,
      });
      rozbrojZdjecie();
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zdjęcie adnotacji', wynik);
      return skutek(`Adnotacja zdjęta bezpowrotnie (${opis}).`);
    },
  ],
  [
    Command.ResearchExcerptList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExcerptList, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz wypisów', wynik);
      return skutek(`Wypisów zebranych z materiału badania: ${String(wynik.wynik.total)}.`);
    },
  ],
  [
    Command.ResearchSourceSummarize,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż źródło do streszczenia.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceSummarize, { sourceId: zrodlo });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Streszczenie źródła', wynik);
      return skutek(
        `Streszczenie złożone modelem: tez ${String(wynik.wynik.summary.theses?.length ?? 0)}, ` +
          `wniosków ${String(wynik.wynik.summary.conclusions?.length ?? 0)}.`,
      );
    },
  ],
  [
    Command.ResearchSourceExtractTable,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż źródło, z którego mają wyjść tabele.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceExtractTable, {
        sourceId: zrodlo,
        persist: true,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wyodrębnienie tabel', wynik);
      return skutek(`Tabel wykrytych i zapisanych: ${String(wynik.wynik.tables.length)}.`);
    },
  ],
  [
    Command.ResearchSourceExtractClaims,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż źródło, z którego mają wyjść twierdzenia.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceExtractClaims, {
        sourceId: zrodlo,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wydobycie twierdzeń', wynik);
      return skutek(`Twierdzeń wydobytych ze źródła: ${String(wynik.wynik.claims.length)}.`);
    },
  ],
  [
    Command.ResearchSourceOcr,
    async ({ stan }) => {
      const zrodlo = wskazaneZrodlo(stan);
      if (zrodlo === '') return brak('Wskaż skan do rozpoznania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchSourceOcr, { sourceId: zrodlo });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Rozpoznanie pisma', wynik);
      return skutek(
        `Stron przetworzonych: ${String(wynik.wynik.pagesProcessed)}; treść ` +
          (wynik.wynik.indexed ? 'weszła do wyszukiwania.' : 'pozostała pusta.'),
      );
    },
  ],
  [
    Command.ResearchCorpusAsk,
    async ({ stan, tekst }) => {
      const pytanie = (tekst ?? '').trim();
      if (pytanie === '') return brak('Wpisz pytanie do korpusu źródeł.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCorpusAsk, {
        windowId: stan.idOkna(),
        question: pytanie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapytanie o korpus', wynik);
      return skutek(
        `Odpowiedź ${wynik.wynik.answer.grounded ? 'zakotwiczona' : 'niezakotwiczona'} ` +
          `w ${String(wynik.wynik.answer.citations.length)} cytatach.`,
      );
    },
  ],

  // ── ustalenia i analiza ──────────────────────────────────────────────────
  [
    Command.ResearchFindingList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingList, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Odczyt wykazu ustaleń', wynik);
      for (const ustalenie of wynik.wynik.findings) stan.wchlonUstalenie(ustalenie);
      return skutek(
        `Rdzeń oddał ${String(wynik.wynik.findings.length)} z ${String(wynik.wynik.total)} ustaleń okna.`,
      );
    },
  ],
  [
    Command.ResearchFindingUpdate,
    async ({ stan }) => {
      const ustalenie = wskazaneUstalenie(stan);
      if (ustalenie === '') return brak('Zaznacz ustalenie, które ma zostać oznaczone jako kluczowe.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingUpdate, {
        findingId: ustalenie,
        weight: 'key',
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zmiana ustalenia', wynik);
      stan.wchlonUstalenie(wynik.wynik.finding);
      return skutek('Ustalenie oznaczone jako kluczowe; zmiana weszła do śladu prowenancji.');
    },
  ],
  [
    Command.ResearchFindingRemove,
    async ({ stan }) => {
      const ustalenie = wskazaneUstalenie(stan);
      if (ustalenie === '') return brak('Zaznacz ustalenie do usunięcia.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingRemove, {
        findingId: ustalenie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Usunięcie ustalenia', wynik);
      return skutek(
        `Ustalenie zdjęte; odwołanie do niego straciło ${String(wynik.wynik.detachedSections)} sekcji raportu.`,
      );
    },
  ],
  [
    Command.ResearchFindingMerge,
    async ({ stan }) => {
      const wybrane = stan.wybraneUstalenia.wybrane();
      if (wybrane.length < 2) return brak('Scalenie potrzebuje co najmniej dwóch zaznaczonych ustaleń.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingMerge, {
        targetFindingId: wybrane[0],
        mergedFindingIds: [...wybrane.slice(1)],
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Scalenie ustaleń', wynik);
      stan.wchlonUstalenie(wynik.wynik.finding);
      return skutek(`Scalono; przeniesionych odwołań do źródeł: ${String(wynik.wynik.movedSources)}.`);
    },
  ],
  [
    Command.ResearchFindingProvenance,
    async ({ stan }) => {
      const ustalenie = wskazaneUstalenie(stan);
      if (ustalenie === '') return brak('Zaznacz ustalenie, którego ślad chcesz obejrzeć.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingProvenance, {
        findingId: ustalenie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Ślad prowenancji', wynik);
      return skutek(`Wpisów śladu: ${String(wynik.wynik.entries.length)}.`);
    },
  ],
  [
    Command.ResearchFindingCode,
    async ({ stan, tekst }) => {
      const ustalenie = wskazaneUstalenie(stan);
      if (ustalenie === '') return brak('Zaznacz ustalenie do zakodowania.');
      const nazwy = (tekst ?? '')
        .split(',')
        .map((wpis) => wpis.trim())
        .filter((wpis) => wpis !== '');
      if (nazwy.length === 0) return brak('Wpisz nazwy kodów po przecinku.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingCode, {
        findingId: ustalenie,
        newCodeNames: nazwy,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Kodowanie ustalenia', wynik);
      return skutek(`Ustalenie ma teraz ${String(wynik.wynik.codes.length)} kodów tematycznych.`);
    },
  ],
  [
    Command.ResearchCodebookGet,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCodebookGet, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Książka kodów', wynik);
      return skutek(`Kodów w książce kodów badania: ${String(wynik.wynik.codes.length)}.`);
    },
  ],
  [
    Command.ResearchCodebookSet,
    async ({ stan, tekst }) => {
      const wpisane = kodyZTekstu(tekst ?? '');
      if (wpisane.length === 0) {
        return brak(
          'Wpisz kody — jeden w wierszu, a po znaku | definicję kodu. Zapis bez kodów wymazałby ' +
            'książkę kodów badania, a tego nikt tu nie zamawia.',
        );
      }
      // Zapis jest całościowy, więc najpierw odczyt: bez niego kody zastane
      // zniknęłyby z książki kodów bez ani jednego zdania o tym.
      const zastane = await stan.zrodlo.wywolaj(Command.ResearchCodebookGet, {
        windowId: stan.idOkna(),
      });
      if (!zastane.udany || zastane.wynik === undefined) {
        return odmowa('Odczyt książki kodów przed zapisem', zastane);
      }
      const zlaczenie = zlaczKody(zastane.wynik.codes, wpisane);
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCodebookSet, {
        windowId: stan.idOkna(),
        codes: [...zlaczenie.kody],
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis książki kodów', wynik);
      const oddane = wynik.wynik.codes.length;
      const zamowione = zlaczenie.kody.length;
      return {
        udany: oddane === zamowione,
        opis:
          `Książka kodów po zapisie ma ${String(oddane)} kodów: dopisanych ${String(zlaczenie.dopisane)}, ` +
          `z poprawioną definicją ${String(zlaczenie.poprawione)}, zastanych bez zmiany ` +
          `${String(zlaczenie.zachowane)}.` +
          (oddane === zamowione
            ? ' Zdjęcia kodu z książki to okno dziś nie ma — kod pominięty w polu zostaje w rdzeniu.'
            : ` UWAGA: wysłano ${String(zamowione)} kodów, rdzeń oddał ${String(oddane)} — sprawdź książkę kodów odczytem.`),
      };
    },
  ],
  [
    Command.ResearchFindingMatrix,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingMatrix, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Macierz kodowania', wynik);
      return skutek(
        `Macierz: ${String(wynik.wynik.codes.length)} kodów × ${String(wynik.wynik.sourceIds.length)} ` +
          `źródeł, komórek niepustych ${String(wynik.wynik.cells.length)}.`,
      );
    },
  ],
  [
    Command.ResearchFindingContradictions,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingContradictions, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykrycie sprzeczności', wynik);
      const otwarte = wynik.wynik.contradictions.filter((pozycja) => !pozycja.resolved).length;
      return skutek(
        `Sprzeczności zapisanych: ${String(wynik.wynik.contradictions.length)}, w tym otwartych ${String(otwarte)}.`,
      );
    },
  ],
  [
    Command.ResearchContradictionResolve,
    async ({ stan, tekst }) => {
      const uzasadnienie = (tekst ?? '').trim();
      if (uzasadnienie === '') {
        return brak('Wpisz uzasadnienie — rozstrzygnięcie bez powodu jest zamknięciem sprawy, nie rozstrzygnięciem.');
      }
      const lista = await stan.zrodlo.wywolaj(Command.ResearchFindingContradictions, {
        windowId: stan.idOkna(),
      });
      if (!lista.udany || lista.wynik === undefined) return odmowa('Odczyt sprzeczności', lista);
      const otwarta = lista.wynik.contradictions.find((pozycja) => !pozycja.resolved);
      if (otwarta === undefined) return brak('Nie ma otwartej sprzeczności do rozstrzygnięcia.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchContradictionResolve, {
        contradictionId: otwarta.id,
        rationale: uzasadnienie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Rozstrzygnięcie sprzeczności', wynik);
      return skutek(`Sprzeczność ${wynik.wynik.contradiction.id} zamknięta uzasadnieniem.`);
    },
  ],
  [
    Command.ResearchFindingFactCheck,
    async ({ stan }) => {
      const ustalenie = wskazaneUstalenie(stan);
      if (ustalenie === '') return brak('Zaznacz ustalenie do weryfikacji.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingFactCheck, {
        findingId: ustalenie,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Weryfikacja twierdzenia', wynik);
      return skutek(
        `Werdykt: ${wynik.wynik.result.verdict}; dowodów przywołanych: ` +
          `${String(wynik.wynik.result.evidence?.length ?? 0)}.`,
      );
    },
  ],
  [
    Command.ResearchFindingCluster,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchFindingCluster, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Grupowanie w wątki', wynik);
      return skutek(
        `Wątków: ${String(wynik.wynik.threads.length)}; poza wątkami zostało ` +
          `${String(wynik.wynik.ungroupedFindingIds?.length ?? 0)} ustaleń.`,
      );
    },
  ],

  // ── przestrzeń badania ───────────────────────────────────────────────────
  [
    Command.ResearchWorkspaceGet,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchWorkspaceGet, {});
      if (!wynik.udany || wynik.wynik === undefined) {
        return odmowa('Odczyt przestrzeni badania', wynik);
      }
      const przestrzen = wynik.wynik;
      // Zakres i etapy wchodzą do pamięci badania, bo to ten sam materiał, który
      // zapisuje research.workspace.set — dwa różne zbiory dałyby okno wiodące
      // pokazujące inny zakres niż ten, który rdzeń trzyma.
      stan.wchlonZakres(przestrzen.scope, przestrzen.stages);
      const bezPokrycia = (przestrzen.questions ?? []).filter(
        (pytanie) => (pytanie.sourceIds ?? []).length === 0,
      ).length;
      const brakujace = [
        przestrzen.audience === undefined ? 'odbiorcy raportu' : '',
        przestrzen.protocol === undefined ? 'protokołu badania' : '',
        przestrzen.note === undefined ? 'notatki roboczej' : '',
      ].filter((pozycja) => pozycja !== '');
      return skutek(
        `Rdzeń oddał przestrzeń badania z ${new Date(przestrzen.updatedAt).toLocaleString('pl-PL')}: ` +
          `etapów ${String(przestrzen.stages.length)}, pytań badawczych ` +
          `${String((przestrzen.questions ?? []).length)}, z tego bez ani jednego źródła ${String(bezPokrycia)}.` +
          (brakujace.length === 0
            ? ' Odbiorca, protokół i notatka robocza są zapisane.'
            : ` Rdzeń nie ma zapisanego: ${brakujace.join(', ')}.`),
      );
    },
  ],
  [
    Command.ResearchWorkspaceQuestionSet,
    async ({ stan, tekst }) => {
      const pytania = (tekst ?? '')
        .split('\n')
        .map((wpis) => wpis.trim())
        .filter((wpis) => wpis !== '');
      if (pytania.length === 0) return brak('Wpisz pytania badawcze — po jednym w wierszu.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchWorkspaceQuestionSet, {
        questions: pytania.map((tresc, numer) => ({ id: '', text: tresc, order: numer + 1 })),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis pytań badawczych', wynik);
      return skutek(`Pytań badawczych zapisanych: ${String(wynik.wynik.questions.length)}.`);
    },
  ],
  [
    Command.ResearchWorkspaceCoverage,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchWorkspaceCoverage, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Pokrycie pytań', wynik);
      return skutek(
        `Pytań: ${String(wynik.wynik.questions.length)}, bez ani jednego źródła: ` +
          `${String(wynik.wynik.uncoveredCount)}.`,
      );
    },
  ],
  [
    Command.ResearchWorkspaceNoteSet,
    async ({ stan, tekst }) => {
      const tresc = (tekst ?? '').trim();
      if (tresc === '') return brak('Wpisz treść notatki roboczej.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchWorkspaceNoteSet, { content: tresc });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis notatki', wynik);
      return skutek(`Notatka zapisana; znaków: ${String(wynik.wynik.content.length)}.`);
    },
  ],
  [
    Command.ResearchWorkspaceFreshness,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchWorkspaceFreshness, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wskaźnik świeżości', wynik);
      if (wynik.wynik.newestSourceAt === undefined) {
        return skutek('Badanie nie ma jeszcze ani jednego źródła — świeżości nie ma czego mierzyć.');
      }
      return skutek(
        `Najświeższe źródło: ${new Date(wynik.wynik.newestSourceAt).toLocaleString('pl-PL')}; ` +
          (wynik.wynik.refreshSuggested
            ? `materiał starszy niż ${String(wynik.wynik.staleDays)} dni — warto odświeżyć.`
            : 'materiał aktualny.'),
      );
    },
  ],
  [
    Command.ResearchGapFind,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchGapFind, { windowId: stan.idOkna() });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykrycie luk badawczych', wynik);
      if (wynik.wynik.gaps.length === 0) return skutek('Każde pytanie badawcze ma poparcie w źródłach.');
      return skutek(`Luk badawczych: ${String(wynik.wynik.gaps.length)}; pierwsza: ${wynik.wynik.gaps[0].summary}`);
    },
  ],
  [
    Command.ResearchPrismaGet,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchPrismaGet, { windowId: stan.idOkna() });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Diagram PRISMA', wynik);
      const liczniki = wynik.wynik.counts;
      return skutek(
        `Zidentyfikowane ${String(liczniki.identified)} → przesiane ${String(liczniki.screened)} → ` +
          `włączone ${String(liczniki.included)} (wyłączonych ${String(liczniki.excluded)}).`,
      );
    },
  ],
  [
    Command.ResearchEvidenceGraph,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchEvidenceGraph, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Graf dowodów', wynik);
      return skutek(
        `Evidence Map: węzłów ${String(wynik.wynik.nodes.length)}, krawędzi ${String(wynik.wynik.edges.length)}.`,
      );
    },
  ],

  // ── cytowania ────────────────────────────────────────────────────────────
  [
    Command.ResearchCitationRender,
    async ({ stan }) => {
      const wybrane = stan.wybraneZrodla.wybrane();
      if (wybrane.length === 0) return brak('Zaznacz źródła do zacytowania.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCitationRender, {
        sourceIds: [...wybrane],
        styleId: 'apa',
        mode: ResearchCitationMode.Both,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Złożenie cytatów', wynik);
      const niekompletne = wynik.wynik.incompleteSourceIds ?? [];
      return skutek(
        `Cytatów złożonych: ${String(wynik.wynik.citations.length)}` +
          (niekompletne.length > 0
            ? `; niekompletnych metadanych: ${String(niekompletne.length)}.`
            : '.'),
      );
    },
  ],
  [
    Command.ResearchCitationStyles,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCitationStyles, {});
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz stylów', wynik);
      return skutek(
        `Stylów dostępnych: ${String(wynik.wynik.total)} — ` +
          wynik.wynik.styles.map((styl) => styl.name).join(', '),
      );
    },
  ],
  [
    Command.ResearchCitationCheck,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchCitationCheck, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Kontrola cytowań', wynik);
      if (wynik.wynik.issues.length === 0) return skutek('Metadane cytowania są kompletne.');
      return skutek(`Uchybień cytowania: ${String(wynik.wynik.issues.length)}.`);
    },
  ],
  [
    Command.ResearchRetractionCheck,
    async ({ stan }) => {
      const wybrane = stan.wybraneZrodla.wybrane();
      if (wybrane.length === 0) return brak('Zaznacz źródła do sprawdzenia.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchRetractionCheck, {
        sourceIds: [...wybrane],
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Kontrola wycofań', wynik);
      const wycofane = wynik.wynik.flags.filter((flaga) => flaga.status === 'wycofana').length;
      return skutek(
        `Sprawdzonych pozycji: ${String(wynik.wynik.flags.length)}; wycofanych: ${String(wycofane)}.`,
      );
    },
  ],

  // ── raport ───────────────────────────────────────────────────────────────
  [
    Command.ResearchReportGet,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportGet, { windowId: stan.idOkna() });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Odczyt raportu', wynik);
      if (wynik.wynik.report === undefined) return skutek('To okno badania nie ma jeszcze raportu.');
      stan.wchlonRaport(wynik.wynik.report);
      return skutek(
        `Raport „${wynik.wynik.report.title}" ma ${String(wynik.wynik.report.sections?.length ?? 0)} sekcji.`,
      );
    },
  ],
  [
    Command.ResearchReportTemplateList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportTemplateList, {});
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz szablonów raportu', wynik);
      return skutek(
        `Szablonów struktury: ${String(wynik.wynik.templates.length)} — ` +
          wynik.wynik.templates.map((szablon) => szablon.name).join(', '),
      );
    },
  ],
  [
    Command.ResearchReportSummarize,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim poprosisz o streszczenie zarządcze.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportSummarize, { reportId: raport });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Streszczenie zarządcze', wynik);
      return skutek(`Sekcja „${wynik.wynik.section.title}" wstawiona na początek raportu.`);
    },
  ],
  [
    Command.ResearchReportBibliography,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim poprosisz o bibliografię końcową.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportBibliography, {
        reportId: raport,
        styleId: 'apa',
        scope: ResearchBibliographyScope.All,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Bibliografia końcowa', wynik);
      return skutek(`Pozycji bibliografii: ${String(wynik.wynik.entries.length)}.`);
    },
  ],
  [
    Command.ResearchReportFootnoteSet,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim ustawisz przypisy.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportFootnoteSet, {
        reportId: raport,
        placement: 'dolne',
        shortForms: true,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Menedżer przypisów', wynik);
      return skutek(`Przypisów w dokumencie: ${String(wynik.wynik.footnoteCount)}.`);
    },
  ],
  [
    Command.ResearchReportInsert,
    async ({ stan }) => {
      const raport = stan.raport();
      const sekcja = raport?.sections?.[0]?.id ?? '';
      if (raport === null || sekcja === '') {
        return brak('Złóż raport z co najmniej jedną sekcją — wstawka potrzebuje miejsca.');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportInsert, {
        reportId: raport.id,
        sectionId: sekcja,
        kind: ResearchBlockKind.EvidenceTable,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wstawienie zestawienia', wynik);
      return skutek(
        `Tabela dowodów wstawiona do sekcji „${raport.sections?.[0].title ?? sekcja}" ` +
          `(${String(wynik.wynik.block.headers?.length ?? 0)} kolumn).`,
      );
    },
  ],
  [
    Command.ResearchReportVersionList,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim poprosisz o wykaz wersji.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportVersionList, {
        reportId: raport,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Wykaz wersji raportu', wynik);
      return skutek(`Wersji raportu: ${String(wynik.wynik.versions.length)}.`);
    },
  ],
  [
    Command.ResearchReportDiff,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim poprosisz o porównanie wersji.');
      const wersje = await stan.zrodlo.wywolaj(Command.ResearchReportVersionList, {
        reportId: raport,
      });
      if (!wersje.udany || wersje.wynik === undefined) return odmowa('Wykaz wersji raportu', wersje);
      if (wersje.wynik.versions.length < 2) {
        return brak('Porównanie potrzebuje dwóch wersji; raport ma dziś mniej.');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportDiff, {
        reportId: raport,
        baseVersionId: wersje.wynik.versions[1].id,
        targetVersionId: wersje.wynik.versions[0].id,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Porównanie wersji', wynik);
      const zmienione = wynik.wynik.hunks.filter((fragment) => fragment.kind !== 'context').length;
      return skutek(
        `Fragmentów porównania: ${String(wynik.wynik.hunks.length)}, w tym zmienionych ${String(zmienione)}.`,
      );
    },
  ],
  [
    Command.ResearchReportCommentList,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim otworzysz tryb recenzji.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportCommentList, {
        reportId: raport,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Komentarze recenzji', wynik);
      return skutek(`Komentarzy w dokumencie: ${String(wynik.wynik.comments.length)}.`);
    },
  ],
  [
    Command.ResearchReportCommentAdd,
    async ({ stan, tekst }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim dopiszesz komentarz.');
      const tresc = (tekst ?? '').trim();
      if (tresc === '') return brak('Wpisz treść komentarza.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportCommentAdd, {
        reportId: raport,
        content: tresc,
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis komentarza', wynik);
      return skutek(`Komentarz ${wynik.wynik.comment.id} dopisany do dokumentu.`);
    },
  ],
  [
    Command.ResearchReportContextualOp,
    async ({ stan }) => {
      const raport = stan.raport();
      const sekcja = raport?.sections?.[0]?.id ?? '';
      if (raport === null || sekcja === '') {
        return brak('Złóż raport z sekcją, zanim uruchomisz operację na zaznaczeniu.');
      }
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchReportContextualOp, {
        reportId: raport.id,
        sectionId: sekcja,
        actionId: 'korekta',
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Operacja na zaznaczeniu', wynik);
      return skutek(
        `Sekcja poprawiona; znaków w wyniku: ${String(wynik.wynik.resultText?.length ?? 0)}.`,
      );
    },
  ],

  // ── eksport ──────────────────────────────────────────────────────────────
  [
    Command.ResearchExportPreview,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż raport, zanim poprosisz o podgląd.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExportPreview, {
        reportId: raport,
        format: 'markdown',
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Podgląd przed eksportem', wynik);
      return skutek(`Podgląd złożony; znaków: ${String(wynik.wynik.preview.text?.length ?? 0)}.`);
    },
  ],
  [
    Command.ResearchExportList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExportList, {
        windowId: stan.idOkna(),
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Historia eksportów', wynik);
      return skutek(`Wydań zapisanych w rdzeniu: ${String(wynik.wynik.exports.length)}.`);
    },
  ],
  [
    Command.ResearchExportTemplateList,
    async ({ stan }) => {
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExportTemplateList, {});
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Szablony eksportu', wynik);
      return skutek(`Szablonów eksportu: ${String(wynik.wynik.templates.length)}.`);
    },
  ],
  [
    Command.ResearchExportTemplateSet,
    async ({ stan, tekst }) => {
      const nazwa = (tekst ?? '').trim();
      if (nazwa === '') return brak('Nadaj szablonowi nazwę.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExportTemplateSet, {
        name: nazwa,
        format: 'pdf',
        content: {
          titlePage: true,
          tableOfContents: true,
          bibliography: true,
          footer: true,
        },
      });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Zapis szablonu eksportu', wynik);
      return skutek(`Szablon „${wynik.wynik.template.name}" zapisany.`);
    },
  ],
  [
    Command.ResearchExportShare,
    async ({ stan }) => {
      const raport = wskazanyRaport(stan);
      if (raport === '') return brak('Złóż i wydaj raport, zanim poprosisz o odnośnik.');
      const wynik = await stan.zrodlo.wywolaj(Command.ResearchExportShare, { reportId: raport });
      if (!wynik.udany || wynik.wynik === undefined) return odmowa('Udostępnienie raportu', wynik);
      return skutek(`Odnośnik do raportu: ${wynik.wynik.url}`);
    },
  ],
]);
