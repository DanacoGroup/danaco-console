// Czynności wydobycia treści ze źródła: przechwycenie strony, partia adresów,
// rozpoznanie pisma, transkrypcja nagrania, tabele, twierdzenia i załączniki.
import {
  Command,
  ResearchAttachmentKind,
  ResearchCaptureMode,
  type ResearchSource,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { niegotowyPanel, wykazPanelu } from './okno-modulu.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  powiedz,
  wybierz,
  zapytaj,
  type KontekstBadania,
} from './research-czynnosci.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-wydobycia';
const GNIAZDO_WYNIKU = '[data-wydobycie-wynik]';

const SPOSOBY_POZYSKANIA: readonly (readonly [string, ResearchCaptureMode])[] = [
  ['Sama treść czytelna', ResearchCaptureMode.Readable],
  ['Sama migawka strony', ResearchCaptureMode.Snapshot],
  ['Treść wraz z migawką', ResearchCaptureMode.Both],
];

const RODZAJE_ZALACZNIKA: readonly (readonly [string, ResearchAttachmentKind])[] = [
  ['Pełny tekst', ResearchAttachmentKind.Fulltext],
  ['Migawka strony', ResearchAttachmentKind.Snapshot],
  ['Notatka źródła', ResearchAttachmentKind.Note],
  ['Dane wydobyte', ResearchAttachmentKind.Data],
];

export function zwiazCzynnosciWydobycia(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'przechwyc') return przechwyc(kontekst);
  if (kod === 'partia') return partia(kontekst);
  if (kod === 'pismo') return rozpoznajPismo(kontekst);
  if (kod === 'przepisz') return przepisz(kontekst);
  if (kod === 'tabele') return tabele(kontekst);
  if (kod === 'twierdzenia') return twierdzenia(kontekst);
  if (kod === 'zalacz') return zalacz(kontekst);
  if (kod === 'zalaczniki') return zalaczniki(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

/* Identyfikatory rdzenia nic Operatorowi nie mówią, więc źródło wskazuje się
   z wypisanego katalogu okna; pusty katalog jest nazwany osobno. */
async function wskazZrodlo(kontekst: KontekstBadania): Promise<ResearchSource | undefined> {
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchSourceList, {
    windowId: kontekst.idOkna(),
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz źródeł nie doszedł, więc nie ma na czym pracować.');
    return undefined;
  }
  const zrodla = wykaz.wynik.sources;
  if (zrodla.length === 0) {
    odmowa(undefined, 'Okno nie ma źródeł, więc wydobycie nie ma przedmiotu.');
    return undefined;
  }
  const numer = wybierz('Źródło', zrodla.map((zrodlo) => zrodlo.title));
  const wybrane = numer < 0 ? undefined : zrodla[numer];
  if (wybrane === undefined) {
    odmowa(undefined, 'Żadne źródło nie zostało wskazane, więc czynność nie poszła do rdzenia.');
    return undefined;
  }
  return wybrane;
}

function sposobPozyskania(): ResearchCaptureMode | undefined {
  const numer = wybierz('Sposób pozyskania', SPOSOBY_POZYSKANIA.map((pozycja) => pozycja[0]));
  const wybrany = numer < 0 ? undefined : SPOSOBY_POZYSKANIA[numer];
  if (wybrany === undefined) {
    odmowa(undefined, 'Sposób pozyskania nie został wskazany, więc adres nie poszedł do rdzenia.');
    return undefined;
  }
  return wybrany[1];
}

/* Wersja archiwalna wchodzi zawsze, bo strona niedostępna w chwili pozyskania
   jest częstsza niż strona żywa, a rdzeń wraca do adresu żywego sam. */
async function przechwyc(kontekst: KontekstBadania): Promise<void> {
  const adres = zapytaj('Adres strony do pozyskania');
  if (adres === '') return;
  const sposob = sposobPozyskania();
  if (sposob === undefined) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceCapture, {
    windowId: kontekst.idOkna(),
    url: adres,
    mode: sposob,
    renderJs: true,
    useWayback: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił pozyskania strony.');
    return;
  }
  const migawka = wynik.wynik.snapshotAttachmentId;
  powiedz(`Wniesiono źródło „${wynik.wynik.source.title}"`
    + (migawka === undefined ? ' bez migawki.' : ' wraz z migawką.'));
  await kontekst.odswiez();
}

async function partia(kontekst: KontekstBadania): Promise<void> {
  const wpis = zapytaj('Adresy do pozyskania, po przecinku');
  if (wpis === '') return;
  const adresy = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  if (adresy.length === 0) {
    odmowa(undefined, 'Żaden adres nie został podany, więc partia nie poszła do rdzenia.');
    return;
  }
  const sposob = sposobPozyskania();
  if (sposob === undefined) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchBatchImport, {
    windowId: kontekst.idOkna(),
    urls: adresy,
    mode: sposob,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił przyjęcia partii adresów.');
    return;
  }
  const powody = wynik.wynik.rejectReasons ?? [];
  powiedz(`Przyjęto ${String(wynik.wynik.accepted)}, odrzucono ${String(wynik.wynik.rejected)}`
    + (powody.length === 0 ? '.' : `: ${powody.join('; ')}`));
  await kontekst.odswiez();
}

/* Rozpoznanie, które nie weszło do wskaźnika pełnotekstowego, jest nazwane
   osobno: bez wskaźnika ani lektura, ani pytanie do korpusu skanu nie widzą. */
async function rozpoznajPismo(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const wpis = zapytaj('Pakiety językowe rozpoznawania, po przecinku');
  const jezyki = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceOcr, {
    sourceId: zrodlo.id,
    languages: jezyki.length === 0 ? undefined : jezyki,
    preprocess: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił rozpoznania pisma w źródle.');
    return;
  }
  powiedz(`Rozpoznano stron: ${String(wynik.wynik.pagesProcessed)}.`);
  if (!wynik.wynik.indexed) {
    odmowa(undefined, 'Rozpoznana treść nie weszła do wskaźnika pełnotekstowego.');
  }
  await kontekst.odswiez();
}

async function przepisz(kontekst: KontekstBadania): Promise<void> {
  const sciezka = zapytaj('Ścieżka nagrania na urządzeniu Operatora');
  if (sciezka === '') return;
  const jezyk = zapytaj('Język nagrania według BCP 47', 'pl-PL');
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceTranscribe, {
    windowId: kontekst.idOkna(),
    sourcePath: sciezka,
    language: jezyk === '' ? undefined : jezyk,
    diarize: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił przepisania nagrania.');
    return;
  }
  if (wynik.wynik.segmentCount === 0) {
    odmowa(undefined, 'Źródło transkryptu powstało, ale rdzeń nie oddał żadnego odcinka mowy.');
    await kontekst.odswiez();
    return;
  }
  powiedz(`Transkrypt „${wynik.wynik.source.title}" ma ${String(wynik.wynik.segmentCount)} odcinków.`);
  await kontekst.odswiez();
}

function numerStrony(pytanie: string): number | undefined {
  const wpis = zapytaj(pytanie);
  if (wpis === '') return undefined;
  const numer = Number.parseInt(wpis, 10);
  return Number.isNaN(numer) || numer < 1 ? undefined : numer;
}

async function tabele(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceExtractTable, {
    sourceId: zrodlo.id,
    page: numerStrony('Strona, na której szukać tabel'),
    persist: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Wydobycie tabel nie doszło.');
    odmowa(wynik.blad, 'Rdzeń odmówił wydobycia tabel ze źródła.');
    return;
  }
  wykazPanelu(cialo, true, undefined, wynik.wynik.tables,
    'Rdzeń nie znalazł w tym źródle żadnej tabeli.',
    (tabela) => [tabela.caption ?? tabela.id,
      `${String(tabela.headers?.length ?? 0)} kolumn`] as const);
}

async function twierdzenia(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceExtractClaims, {
    sourceId: zrodlo.id,
    page: numerStrony('Strona zawężająca wydobycie'),
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Wydobycie twierdzeń nie doszło.');
    odmowa(wynik.blad, 'Rdzeń odmówił wydobycia twierdzeń ze źródła.');
    return;
  }
  wykazPanelu(cialo, true, undefined, wynik.wynik.claims,
    'Rdzeń nie wydobył z tego źródła żadnego twierdzenia.',
    (twierdzenie) => [twierdzenie.text, twierdzenie.kind] as const);
}

async function zalacz(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  const numer = wybierz('Rodzaj załącznika', RODZAJE_ZALACZNIKA.map((pozycja) => pozycja[0]));
  const rodzaj = numer < 0 ? undefined : RODZAJE_ZALACZNIKA[numer];
  if (rodzaj === undefined) {
    odmowa(undefined, 'Rodzaj załącznika nie został wskazany, więc plik nie poszedł do rdzenia.');
    return;
  }
  const sciezka = zapytaj('Ścieżka pliku na urządzeniu Operatora');
  if (sciezka === '') {
    odmowa(undefined, 'Załącznik bez pliku byłby wierszem o pliku, którego nie ma.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceAttachmentAdd, {
    sourceId: zrodlo.id,
    kind: rodzaj[1],
    sourcePath: sciezka,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił dołączenia pliku do źródła.');
    return;
  }
  powiedz(`Załącznik dołączony do źródła „${zrodlo.title}".`);
  await zalacznikiZrodla(kontekst, zrodlo.id);
}

async function zalaczniki(kontekst: KontekstBadania): Promise<void> {
  const zrodlo = await wskazZrodlo(kontekst);
  if (zrodlo === undefined) return;
  await zalacznikiZrodla(kontekst, zrodlo.id);
}

/* Braki kompletności są nazwane wprost, bo bez pełnego tekstu ani lektura,
   ani wydobycie tabel nie mają z czego pracować. */
async function zalacznikiZrodla(kontekst: KontekstBadania, idZrodla: string): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceAttachmentList, {
    sourceId: idZrodla,
  });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.attachments,
    'Źródło nie ma żadnego załącznika.',
    (zalacznik) => [String(zalacznik.kind),
      `${String(zalacznik.sizeBytes ?? 0)} B`] as const);
  const brakujace = wynik.wynik?.missingKinds ?? [];
  if (brakujace.length > 0) {
    odmowa(undefined, `Rodzaje załączników, których źródło nie ma: ${brakujace.join(', ')}.`);
  }
}
