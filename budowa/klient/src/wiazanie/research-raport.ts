// Czynności raportu badania: złożenie, wstawienie zestawienia z ustaleń,
// streszczenie zarządcze, bibliografia i wydanie dokumentu.
import {
  Command,
  ExportFormat,
  ResearchBibliographyScope,
  ResearchBlockKind,
  ResearchExportTarget,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  brakOkna,
  doSchowka,
  nasluchCzynnosci,
  nieznanaCzynnosc,
  odmowa,
  powiedz,
  wybierz,
  zapytaj,
  type KontekstBadania,
} from './research-czynnosci.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-raportu';

export function zwiazCzynnosciRaportu(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'zloz') return zloz(kontekst);
  if (kod === 'wstaw') return wstaw(kontekst);
  if (kod === 'stresc') return stresc(kontekst);
  if (kod === 'bibliografia') return bibliografia(kontekst);
  if (kod === 'wydaj') return wydaj(kontekst);
  nieznanaCzynnosc(kod);
}

/* Raport bieżący okna jest jedynym, do którego okno ma dziś drogę; brak raportu
   jest nazwany, bo dalsze czynności bez niego nie mają przedmiotu. */
export async function idRaportuBadania(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.ResearchReportGet, { windowId: idOkna });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń nie oddał raportu tego okna.');
    return '';
  }
  const raport = wynik.wynik?.report;
  if (raport === undefined) {
    odmowa(undefined, 'Okno nie ma jeszcze raportu — najpierw go złóż.');
    return '';
  }
  return raport.id;
}

async function zloz(kontekst: KontekstBadania): Promise<void> {
  const tytul = zapytaj('Tytuł raportu');
  if (tytul === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportBuild, {
    windowId: kontekst.idOkna(),
    title: tytul,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia raportu.');
    return;
  }
  const sekcje = wynik.wynik.report.sections?.length ?? 0;
  if (wynik.wynik.fromModel === false) {
    odmowa(undefined, 'Sekcje raportu są komunikatem procesu kanału, a nie odpowiedzią modelu.');
    return;
  }
  powiedz(`Raport złożony; sekcji: ${String(sekcje)}.`);
  await kontekst.odswiez();
}

async function wstaw(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const raport = await wywolaj(kontekst.kanal, Command.ResearchReportGet, {
    windowId: kontekst.idOkna(),
  });
  const sekcje = raport.wynik?.report?.sections ?? [];
  const numer = wybierz('Sekcja, w której ma stanąć tabela dowodów',
    sekcje.map((sekcja) => sekcja.title));
  const sekcja = numer < 0 ? undefined : sekcje[numer];
  if (sekcja === undefined) {
    odmowa(undefined, 'Żadna sekcja nie została wskazana, więc zestawienie nie poszło do rdzenia.');
    return;
  }
  const ustalenia = await wywolaj(kontekst.kanal, Command.ResearchFindingList, {
    windowId: kontekst.idOkna(),
  });
  const kody = (ustalenia.wynik?.findings ?? []).map((jedno) => jedno.id);
  if (kody.length === 0) {
    odmowa(undefined, 'Okno nie ma ustaleń, więc tabela dowodów nie miałaby czym stanąć.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportInsert, {
    reportId,
    sectionId: sekcja.id,
    kind: ResearchBlockKind.EvidenceTable,
    findingIds: kody,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił osadzenia zestawienia w sekcji.');
    return;
  }
  powiedz(`Tabela dowodów stoi w sekcji „${sekcja.title}" na ${String(kody.length)} ustaleniach.`);
  await kontekst.odswiez();
}

async function stresc(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportSummarize, { reportId });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia streszczenia zarządczego.');
    return;
  }
  if (!wynik.wynik.fromModel) {
    odmowa(undefined, 'Rdzeń oddał komunikat kanału zamiast streszczenia modelu.');
    return;
  }
  powiedz(`Sekcja „${wynik.wynik.section.title}" stoi na początku raportu.`);
  await kontekst.odswiez();
}

async function bibliografia(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const style = await wywolaj(kontekst.kanal, Command.ResearchCitationStyles, {});
  const wykaz = style.wynik?.styles ?? [];
  const numer = wybierz('Styl cytowania', wykaz.map((styl) => styl.name));
  const styl = numer < 0 ? undefined : wykaz[numer];
  if (styl === undefined) {
    odmowa(undefined, 'Styl cytowania nie został wskazany, więc bibliografia nie poszła do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportBibliography, {
    reportId,
    styleId: styl.id,
    scope: ResearchBibliographyScope.Cited,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił złożenia bibliografii.');
    return;
  }
  if (wynik.wynik.entries.length === 0) {
    powiedz('Raport nie cytuje żadnego źródła, więc bibliografia jest pusta.');
    return;
  }
  await doSchowka(wynik.wynik.entries.join('\n'), `bibliografia-${styl.id}`);
}

/* Wydanie oddaje ślad eksportu, ale nie treść, a okno nie ma gdzie zostawić
   pliku — treść bierze się podglądem tego samego formatu i idzie do schowka. */
async function wydaj(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wydane = await wywolaj(kontekst.kanal, Command.ResearchReportExport, {
    reportId,
    format: ExportFormat.Markdown,
    target: ResearchExportTarget.Download,
  });
  if (!wydane.udany || wydane.wynik === undefined) {
    odmowa(wydane.blad, 'Rdzeń odmówił wydania raportu.');
    return;
  }
  const podglad = await wywolaj(kontekst.kanal, Command.ResearchExportPreview, {
    reportId,
    format: ExportFormat.Markdown,
  });
  const tresc = podglad.wynik?.preview.text ?? '';
  if (!podglad.udany || tresc === '') {
    odmowa(podglad.blad, 'Wydanie zapisane, ale rdzeń nie oddał treści dokumentu do przeniesienia.');
    await kontekst.odswiez();
    return;
  }
  if (podglad.wynik?.preview.truncated === true) {
    odmowa(undefined, 'Rdzeń oddał podgląd skrócony, więc w schowku stoi niepełny dokument.');
  }
  await doSchowka(tresc, wydane.wynik.path ?? wydane.wynik.exportId ?? 'raport.md');
  await kontekst.odswiez();
}
