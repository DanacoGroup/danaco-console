// Czynności katalogu źródeł: zmiana, otagowanie, streszczenie i usunięcie
// wiersza oraz wciągnięcie bibliografii, scalenie duplikatów i rozpoznanie
// identyfikatora dla całego panelu.
import {
  Command,
  ResearchImportFormat,
  ResearchSourceKind,
  type ResearchSource,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  brakOkna,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  potwierdzone,
  powiedz,
  wierszCzynnosci,
  wskazanie,
  wybierz,
  zapytaj,
  type CzynnoscWiersza,
  type KontekstBadania,
} from './research-czynnosci.ts';

const ATRYBUT_ZRODLA = 'data-zrodlo-badania';
const ATRYBUT_CZYNNOSCI = 'data-czynnosc-zrodla';

const CZYNNOSCI_WIERSZA: readonly CzynnoscWiersza[] = [
  { kod: 'zmien', etykieta: 'Zmień' },
  { kod: 'etykiety', etykieta: 'Etykiety' },
  { kod: 'stresc', etykieta: 'Streść' },
  { kod: 'usun', etykieta: 'Usuń' },
];

const FORMATY_BIBLIOGRAFII: readonly (readonly [string, ResearchImportFormat])[] = [
  ['BibTeX', ResearchImportFormat.Bibtex],
  ['RIS', ResearchImportFormat.Ris],
  ['CSL-JSON', ResearchImportFormat.CslJson],
  ['EndNote XML', ResearchImportFormat.EndnoteXml],
  ['CSV', ResearchImportFormat.Csv],
];

export function wierszZrodla(cialo: Element, zrodlo: ResearchSource): HTMLElement {
  return wierszCzynnosci(cialo, zrodlo.title, zrodlo.url ?? zrodlo.kind,
    ATRYBUT_ZRODLA, zrodlo.id, ATRYBUT_CZYNNOSCI, CZYNNOSCI_WIERSZA);
}

export function zwiazCzynnosciZrodel(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod, przycisk) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod, wskazanie(przycisk, ATRYBUT_ZRODLA));
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string, idZrodla: string): Promise<void> {
  if (kod === 'zmien') return zmien(kontekst, idZrodla);
  if (kod === 'etykiety') return otaguj(kontekst, idZrodla);
  if (kod === 'stresc') return stresc(kontekst, idZrodla);
  if (kod === 'usun') return usun(kontekst, idZrodla);
  if (kod === 'wciagnij') return wciagnij(kontekst);
  if (kod === 'scal') return scal(kontekst);
  if (kod === 'rozpoznaj') return rozpoznaj(kontekst);
  nieznanaCzynnosc(kod);
}

function brakZrodla(idZrodla: string): boolean {
  if (idZrodla !== '') return false;
  odmowa(undefined, 'Przycisk nie stoi przy żadnym źródle, więc nic nie zostało zmienione.');
  return true;
}

async function zmien(kontekst: KontekstBadania, idZrodla: string): Promise<void> {
  if (brakZrodla(idZrodla)) return;
  const tytul = zapytaj('Nowy tytuł źródła');
  if (tytul === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceUpdate, {
    sourceId: idZrodla,
    title: tytul,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił zmiany źródła.');
    return;
  }
  powiedz(`Źródło nosi teraz tytuł „${wynik.wynik?.source.title ?? tytul}".`);
  await kontekst.odswiez();
}

async function otaguj(kontekst: KontekstBadania, idZrodla: string): Promise<void> {
  if (brakZrodla(idZrodla)) return;
  const wpis = zapytaj('Etykiety tematyczne źródła, po przecinku');
  if (wpis === '') return;
  const etykiety = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  if (etykiety.length === 0) {
    odmowa(undefined, 'Żadna etykieta nie została podana, więc przypisanie nie poszło do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceTag, {
    sourceId: idZrodla,
    tags: etykiety,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił nadania etykiet.');
    return;
  }
  powiedz(`Źródło nosi ${String(etykiety.length)} etykiet.`);
  await kontekst.odswiez();
}

/* Rdzeń rozróżnia odpowiedź modelu od komunikatu procesu kanału polem
   `fromModel`; bez tego rozróżnienia komunikat błędu stanąłby jako abstrakt. */
async function stresc(kontekst: KontekstBadania, idZrodla: string): Promise<void> {
  if (brakZrodla(idZrodla)) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceSummarize, {
    sourceId: idZrodla,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił streszczenia źródła.');
    return;
  }
  if (!wynik.wynik.fromModel) {
    odmowa(undefined, 'Rdzeń oddał komunikat kanału zamiast streszczenia modelu.');
    return;
  }
  powiedz(wynik.wynik.summary.abstract ?? 'Model nie oddał abstraktu tego źródła.');
}

async function usun(kontekst: KontekstBadania, idZrodla: string): Promise<void> {
  if (brakZrodla(idZrodla)) return;
  if (!potwierdzone(`zrodlo:${idZrodla}`,
    'Usunięcie źródła jest nieodwracalne. Naciśnij drugi raz, żeby je wykonać.')) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceRemove, { sourceId: idZrodla });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił usunięcia źródła.');
    return;
  }
  powiedz(`Źródło usunięte; ustaleń, które je utraciły: ${String(wynik.wynik?.detachedFindings ?? 0)}.`);
  await kontekst.odswiez();
}

async function wciagnij(kontekst: KontekstBadania): Promise<void> {
  const sciezka = zapytaj('Ścieżka pliku bibliografii na urządzeniu Operatora');
  if (sciezka === '') return;
  const numer = wybierz('Format bibliografii', FORMATY_BIBLIOGRAFII.map((pozycja) => pozycja[0]));
  const wybrany = numer < 0 ? undefined : FORMATY_BIBLIOGRAFII[numer];
  if (wybrany === undefined) {
    odmowa(undefined, 'Format bibliografii nie został wskazany, więc plik nie poszedł do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceImport, {
    windowId: kontekst.idOkna(),
    format: wybrany[1],
    sourcePath: sciezka,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił wczytania bibliografii.');
    return;
  }
  const powody = wynik.wynik.skipReasons ?? [];
  powiedz(`Wczytano ${String(wynik.wynik.imported)} pozycji, pominięto ${String(wynik.wynik.skipped)}`
    + (powody.length === 0 ? '.' : `: ${powody.join('; ')}`));
  await kontekst.odswiez();
}

/* Scalenie idzie parą wskazaną przez rdzeń, bo to rdzeń liczy podobieństwo;
   okno tylko pokazuje kandydatów i pyta, którą parę zamknąć. */
async function scal(kontekst: KontekstBadania): Promise<void> {
  const wykryte = await wywolaj(kontekst.kanal, Command.ResearchSourceDuplicates, {
    windowId: kontekst.idOkna(),
  });
  if (!wykryte.udany || wykryte.wynik === undefined) {
    odmowa(wykryte.blad, 'Rdzeń odmówił wskazania duplikatów.');
    return;
  }
  const pary = wykryte.wynik.candidates;
  if (pary.length === 0) {
    powiedz('Rdzeń nie wskazał żadnej pary powtórzonych źródeł.');
    return;
  }
  const numer = wybierz('Para źródeł do scalenia',
    pary.map((para) => `${para.sourceIds.join(' + ')} — ${para.basis}`));
  const para = numer < 0 ? undefined : pary[numer];
  if (para === undefined || para.sourceIds.length < 2) {
    odmowa(undefined, 'Nie wskazano pary o dwóch źródłach, więc scalenie nie poszło do rdzenia.');
    return;
  }
  const [zachowane, ...scalane] = para.sourceIds;
  if (zachowane === undefined) return;
  if (!potwierdzone(`scalenie:${zachowane}`,
    'Scalenie usuwa źródła powtórzone. Naciśnij drugi raz, żeby je wykonać.')) return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchSourceMerge, {
    targetSourceId: zachowane,
    mergedSourceIds: scalane,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił scalenia źródeł.');
    return;
  }
  powiedz(`Scalono; przeniesionych powiązań z ustaleniami: ${String(wynik.wynik?.movedFindings ?? 0)}.`);
  await kontekst.odswiez();
}

/* Samo rozpoznanie identyfikatora niczego by w oknie nie zmieniło, więc po nim
   idzie wniesienie pozycji do katalogu — i oba kroki są nazwane. */
async function rozpoznaj(kontekst: KontekstBadania): Promise<void> {
  const identyfikator = zapytaj('Identyfikator pozycji: DOI, ISBN, PMID albo arXiv');
  if (identyfikator === '') return;
  const rozpoznane = await wywolaj(kontekst.kanal, Command.ResearchSourceResolve, {
    identifier: identyfikator,
  });
  if (!rozpoznane.udany || rozpoznane.wynik === undefined) {
    odmowa(rozpoznane.blad, 'Rdzeń nie rozpoznał identyfikatora.');
    return;
  }
  const pozycja = rozpoznane.wynik.result;
  const wniesione = await wywolaj(kontekst.kanal, Command.ResearchSourceAdd, {
    windowId: kontekst.idOkna(),
    title: pozycja.title,
    kind: ResearchSourceKind.Web,
    url: pozycja.url,
    identifier: identyfikator,
    cslJson: rozpoznane.wynik.cslJson,
  });
  if (!wniesione.udany) {
    odmowa(wniesione.blad, 'Pozycja została rozpoznana, ale rdzeń odmówił wniesienia jej do katalogu.');
    return;
  }
  powiedz(`Rozpoznano i wniesiono źródło „${pozycja.title}".`);
  await kontekst.odswiez();
}
