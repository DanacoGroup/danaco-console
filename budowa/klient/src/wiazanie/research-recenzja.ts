// Czynności recenzji i wydania raportu: wersje, różnica, komentarze, przypisy,
// operacja kontekstowa na sekcji, szablony eksportu i odnośnik dla odbiorcy.
import {
  Command,
  ExportFormat,
  ResearchExportTarget,
  type ResearchReportSection,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { niegotowyPanel, wykazPanelu } from './okno-modulu.ts';
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
import { idRaportuBadania } from './research-raport.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

const ATRYBUT_CZYNNOSCI = 'data-czynnosc-recenzji';
const GNIAZDO_WYNIKU = '[data-recenzja-wynik]';

const UMIEJSCOWIENIA_PRZYPISOW: readonly string[] = ['dolne', 'koncowe'];

const OPERACJE_SEKCJI: readonly string[] = ['korekta', 'streszczenie', 'rozwiniecie', 'zmiana-stylu'];

const FORMATY_WYDANIA: readonly (readonly [string, ExportFormat])[] = [
  ['PDF', ExportFormat.Pdf],
  ['DOCX', ExportFormat.Docx],
  ['Markdown', ExportFormat.Markdown],
  ['HTML', ExportFormat.Html],
  ['Slajdy PPTX', ExportFormat.Pptx],
  ['Arkusz XLSX', ExportFormat.Xlsx],
  ['LaTeX', ExportFormat.Latex],
];

const MIEJSCA_WYDANIA: readonly (readonly [string, ResearchExportTarget])[] = [
  ['Pobranie lokalne', ResearchExportTarget.Download],
  ['Repozytorium Library', ResearchExportTarget.Library],
  ['Moduł Studio', ResearchExportTarget.Studio],
  ['Moduł Roundtable', ResearchExportTarget.Roundtable],
];

export function zwiazCzynnosciRecenzji(kontekst: KontekstBadania): Odsubskrybuj {
  return nasluchCzynnosci(kontekst.korzen, ATRYBUT_CZYNNOSCI, async (kod) => {
    if (brakOkna(kontekst.idOkna())) return;
    await wykonaj(kontekst, kod);
  });
}

async function wykonaj(kontekst: KontekstBadania, kod: string): Promise<void> {
  if (kod === 'wersje') return wersje(kontekst);
  if (kod === 'roznica') return roznica(kontekst);
  if (kod === 'komentarz') return dodajKomentarz(kontekst);
  if (kod === 'komentarze') return komentarze(kontekst);
  if (kod === 'przypisy') return przypisy(kontekst);
  if (kod === 'operacja') return operacja(kontekst);
  if (kod === 'szablony') return szablony(kontekst);
  if (kod === 'szablon') return zapiszSzablon(kontekst);
  if (kod === 'odnosnik') return odnosnik(kontekst);
  nieznanaCzynnosc(kod);
}

function gniazdo(kontekst: KontekstBadania): Element | null {
  return kontekst.korzen.querySelector(GNIAZDO_WYNIKU);
}

async function sekcja(kontekst: KontekstBadania): Promise<ResearchReportSection | undefined> {
  const raport = await wywolaj(kontekst.kanal, Command.ResearchReportGet, {
    windowId: kontekst.idOkna(),
  });
  const sekcje = raport.wynik?.report?.sections ?? [];
  if (sekcje.length === 0) {
    odmowa(undefined, 'Raport nie ma sekcji, więc czynność nie ma do czego się odnieść.');
    return undefined;
  }
  const numer = wybierz('Sekcja raportu', sekcje.map((jedna) => jedna.title));
  const wybrana = numer < 0 ? undefined : sekcje[numer];
  if (wybrana === undefined) {
    odmowa(undefined, 'Żadna sekcja nie została wskazana, więc czynność nie poszła do rdzenia.');
  }
  return wybrana;
}

async function wersje(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportVersionList, { reportId });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.versions,
    'Raport nie ma jeszcze żadnej zapisanej wersji.',
    (wersja) => [wersja.label ?? wersja.id,
      `${String(wersja.sectionCount)} sekcji`] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wydania wykazu wersji raportu.');
}

async function roznica(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wykaz = await wywolaj(kontekst.kanal, Command.ResearchReportVersionList, { reportId });
  const lista = wykaz.wynik?.versions ?? [];
  if (!wykaz.udany || lista.length < 2) {
    odmowa(wykaz.blad, 'Porównanie wymaga dwóch wersji raportu, a rdzeń nie prowadzi tylu.');
    return;
  }
  const etykiety = lista.map((wersja) => wersja.label ?? wersja.id);
  const pierwszy = wybierz('Wersja odniesienia', etykiety);
  const drugi = wybierz('Wersja porównywana', etykiety);
  const odniesienie = pierwszy < 0 ? undefined : lista[pierwszy];
  const porownywana = drugi < 0 ? undefined : lista[drugi];
  if (odniesienie === undefined || porownywana === undefined) {
    odmowa(undefined, 'Obie wersje muszą być wskazane, więc porównanie nie poszło do rdzenia.');
    return;
  }
  if (odniesienie.id === porownywana.id) {
    odmowa(undefined, 'Wskazano dwa razy tę samą wersję, więc różnica byłaby pusta z założenia.');
    return;
  }
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportDiff, {
    reportId,
    baseVersionId: odniesienie.id,
    targetVersionId: porownywana.id,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    niegotowyPanel(cialo, wynik.blad?.message ?? 'Zestawienie różnicowe nie doszło.');
    odmowa(wynik.blad, 'Rdzeń odmówił porównania wersji raportu.');
    return;
  }
  wykazPanelu(cialo, true, undefined, wynik.wynik.hunks,
    'Wersje raportu nie różnią się treścią.',
    (fragment) => [fragment.after ?? fragment.before ?? String(fragment.index),
      String(fragment.kind)] as const);
}

async function dodajKomentarz(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wybrana = await sekcja(kontekst);
  if (wybrana === undefined) return;
  const tresc = zapytaj('Treść komentarza recenzenckiego');
  if (tresc === '') return;
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportCommentAdd, {
    reportId,
    sectionId: wybrana.id,
    content: tresc,
  });
  if (!wynik.udany) {
    odmowa(wynik.blad, 'Rdzeń odmówił założenia komentarza.');
    return;
  }
  powiedz(`Komentarz stoi przy sekcji „${wybrana.title}".`);
  await komentarze(kontekst);
}

async function komentarze(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportCommentList, { reportId });
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.comments,
    'Raport nie ma żadnego komentarza recenzenckiego.',
    (komentarz) => [komentarz.content,
      komentarz.resolved ? 'rozwiązany' : 'otwarty'] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wydania komentarzy raportu.');
}

async function przypisy(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const numer = wybierz('Umiejscowienie przypisów', UMIEJSCOWIENIA_PRZYPISOW);
  const umiejscowienie = numer < 0 ? undefined : UMIEJSCOWIENIA_PRZYPISOW[numer];
  if (umiejscowienie === undefined) {
    odmowa(undefined, 'Umiejscowienie przypisów nie zostało wskazane, więc nic nie poszło do rdzenia.');
    return;
  }
  const skrocone = wybierz('Postać przypisów', ['Pełne', 'Skrócone według stylu']);
  if (skrocone < 0) {
    odmowa(undefined, 'Postać przypisów nie została wskazana, więc nic nie poszło do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportFootnoteSet, {
    reportId,
    placement: umiejscowienie,
    shortForms: skrocone === 1,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił ustawienia przypisów raportu.');
    return;
  }
  powiedz(`Przypisów po przeliczeniu: ${String(wynik.wynik.footnoteCount)}.`);
}

/* Operacja oddaje albo gotową treść, albo wiadomość, pod którą biegnie, albo
   propozycję do zatwierdzenia. Każdy z tych trzech przypadków jest nazwany,
   bo milczenie po naciśnięciu wyglądałoby jak wykonana zmiana. */
async function operacja(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wybrana = await sekcja(kontekst);
  if (wybrana === undefined) return;
  const numer = wybierz('Operacja na sekcji', OPERACJE_SEKCJI);
  const czynnosc = numer < 0 ? undefined : OPERACJE_SEKCJI[numer];
  if (czynnosc === undefined) {
    odmowa(undefined, 'Operacja nie została wskazana, więc sekcja nie poszła do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchReportContextualOp, {
    reportId,
    sectionId: wybrana.id,
    actionId: czynnosc,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił wykonania operacji na sekcji.');
    return;
  }
  const tresc = wynik.wynik.resultText ?? '';
  if (tresc !== '') {
    await doSchowka(tresc, `${wybrana.title} — ${czynnosc}`);
    await kontekst.odswiez();
    return;
  }
  if (wynik.wynik.proposalId !== undefined) {
    powiedz('Rdzeń złożył propozycję zmiany i czeka na jej zatwierdzenie.');
    return;
  }
  if (wynik.wynik.messageId !== undefined) {
    powiedz('Operacja biegnie w kanale; treść wejdzie do raportu po jej zakończeniu.');
    return;
  }
  odmowa(undefined, 'Rdzeń przyjął operację, ale nie oddał ani treści, ani śladu jej biegu.');
}

async function szablony(kontekst: KontekstBadania): Promise<void> {
  const cialo = gniazdo(kontekst);
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchExportTemplateList, {});
  wykazPanelu(cialo, wynik.udany, wynik.blad?.message, wynik.wynik?.templates,
    'Rdzeń nie ma zapisanego szablonu eksportu.',
    (szablon) => [szablon.name,
      `${String(szablon.format)} · ${String(szablon.target ?? 'bez miejsca')}`] as const);
  if (!wynik.udany) odmowa(wynik.blad, 'Rdzeń odmówił wydania szablonów eksportu.');
}

async function zapiszSzablon(kontekst: KontekstBadania): Promise<void> {
  const nazwa = zapytaj('Nazwa szablonu eksportu');
  if (nazwa === '') return;
  const numerFormatu = wybierz('Format wyjściowy', FORMATY_WYDANIA.map((pozycja) => pozycja[0]));
  const format = numerFormatu < 0 ? undefined : FORMATY_WYDANIA[numerFormatu];
  if (format === undefined) {
    odmowa(undefined, 'Format wyjściowy nie został wskazany, więc szablon nie poszedł do rdzenia.');
    return;
  }
  const numerMiejsca = wybierz('Miejsce docelowe', MIEJSCA_WYDANIA.map((pozycja) => pozycja[0]));
  const miejsce = numerMiejsca < 0 ? undefined : MIEJSCA_WYDANIA[numerMiejsca];
  if (miejsce === undefined) {
    odmowa(undefined, 'Miejsce docelowe nie zostało wskazane, więc szablon nie poszedł do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchExportTemplateSet, {
    name: nazwa,
    format: format[1],
    target: miejsce[1],
    content: {
      titlePage: true,
      tableOfContents: true,
      bibliography: true,
      footer: true,
      evidenceTable: true,
    },
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił zapisania szablonu eksportu.');
    return;
  }
  powiedz(`Szablon „${wynik.wynik.template.name}" zapisany.`);
  await szablony(kontekst);
}

async function odnosnik(kontekst: KontekstBadania): Promise<void> {
  const reportId = await idRaportuBadania(kontekst.kanal, kontekst.idOkna());
  if (reportId === '') return;
  const wpis = zapytaj('Ważność odnośnika w minutach; puste znaczy bez ograniczenia');
  const minuty = wpis === '' ? undefined : Number.parseInt(wpis, 10);
  if (minuty !== undefined && (Number.isNaN(minuty) || minuty < 1)) {
    odmowa(undefined, 'Ważność musi być liczbą minut, więc odnośnik nie poszedł do rdzenia.');
    return;
  }
  const wynik = await wywolaj(kontekst.kanal, Command.ResearchExportShare, {
    reportId,
    expiresInMinutes: minuty,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    odmowa(wynik.blad, 'Rdzeń odmówił wydania odnośnika do raportu.');
    return;
  }
  await doSchowka(wynik.wynik.url, `odnosnik-${reportId}`);
}
