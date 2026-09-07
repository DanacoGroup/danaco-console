/**
 * Kontrola jakości okna Translate: sprawdzenie jakości i spójności, korekta
 * wraz z zastosowaniem propozycji, tłumaczenie zwrotne, etap zatwierdzenia
 * oraz profile kontroli jakości.
 */

import {
  ApprovalStage,
  Command,
  ProofreadSeverity,
  TranslationIssueKind,
  type ProofreadFinding,
} from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  odmowa,
  oknoWskazane,
  panelDocelowy,
  pokazWynik,
  potwierdzNieodwracalna,
  powiedz,
  przyjmijPanele,
  zapytaj,
  zapytajOWartosc,
  zapytajOWartosci,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';

const PANEL = 'panel-plan';

export function czynnosciJakosci(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  /* Propozycje korekty wchodzą do treści po identyfikatorach, a te niesie
     wyłącznie odpowiedź ostatniego przebiegu — okno trzyma je do zastosowania. */
  let ostatniaKorekta: { panelId: string; zastrzezenia: ProofreadFinding[] } | null = null;

  return {
    'jakosc-sprawdz': () => sprawdzJakosc(stan),
    'spojnosc-sprawdz': () => sprawdzSpojnosc(stan),
    'korekta-uruchom': async (): Promise<void> => {
      ostatniaKorekta = await uruchomKorekte(stan);
    },
    'korekta-zastosuj': () => zastosujKorekte(stan, ostatniaKorekta),
    'zwrotne-tlumaczenie': () => przetlumaczZwrotnie(stan),
    'zatwierdzenie-ustaw': () => ustawEtap(stan),
    'profil-jakosci-zapisz': () => zapiszProfil(stan),
    'profil-jakosci-usun': () => usunProfil(stan),
  };
}

async function sprawdzJakosc(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateQualityCheck, { panelId: idPanelu });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił kontroli jakości.');
    return;
  }
  pokazWynik(stan, PANEL, odpowiedz.wynik.issues.map((tresc) => [tresc, ''] as const),
    'Kontrola jakości nie wniosła zastrzeżeń.');
  powiedz(`Zastrzeżeń kontroli: ${odpowiedz.wynik.issues.length}.`);
}

async function sprawdzSpojnosc(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateConsistencyCheck, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił kontroli spójności.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.findings.map((n) => [n.sourceText, `${n.kind}: ${n.variants.join(' / ')}`] as const),
    'Kontrola spójności nie znalazła rozbieżności.');
}

async function uruchomKorekte(
  stan: StanTlumaczenia,
): Promise<{ panelId: string; zastrzezenia: ProofreadFinding[] } | null> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return null;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateProofreadRun, { panelId: idPanelu });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił korekty panelu.');
    return null;
  }
  const zastrzezenia = odpowiedz.wynik.findings;
  const czytelnosc = odpowiedz.wynik.readability
    .map((miara) => [miara.metric, String(miara.value)] as const);
  pokazWynik(stan, PANEL,
    [...zastrzezenia.map((z) => [z.detail, `${z.kind} · ${z.severity}`] as const), ...czytelnosc],
    'Korekta nie wniosła zastrzeżeń.');
  powiedz(`Zastrzeżeń korekty: ${zastrzezenia.length}.`);
  return { panelId: idPanelu, zastrzezenia };
}

async function zastosujKorekte(
  stan: StanTlumaczenia,
  ostatnia: { panelId: string; zastrzezenia: ProofreadFinding[] } | null,
): Promise<void> {
  if (ostatnia === null || ostatnia.zastrzezenia.length === 0) {
    odmowa(undefined, 'Żaden przebieg korekty nie zostawił propozycji do zastosowania.');
    return;
  }
  const zPropozycja = ostatnia.zastrzezenia.filter((z) => (z.suggestion ?? '') !== '');
  if (zPropozycja.length === 0) {
    odmowa(undefined, 'Zastrzeżenia ostatniej korekty nie niosą propozycji treści.');
    return;
  }
  if (!potwierdzNieodwracalna('korekta-zastosuj',
    `${zPropozycja.length} propozycji wejdzie do treści panelu.`)) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateProofreadApply, {
    panelId: ostatnia.panelId,
    findingIds: zPropozycja.map((z) => z.id),
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zastosowania korekty.');
    return;
  }
  przyjmijPanele(stan, [odpowiedz.wynik.panel]);
  powiedz(`Zastosowano ${odpowiedz.wynik.appliedCount} propozycji.`);
  await stan.odswiez();
}

async function przetlumaczZwrotnie(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateBacktranslationRun, {
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił tłumaczenia zwrotnego.');
    return;
  }
  pokazWynik(stan, PANEL, [[odpowiedz.wynik.text, 'tłumaczenie zwrotne']],
    'Tłumaczenie zwrotne wróciło puste.');
}

async function ustawEtap(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const etap = zapytajOWartosc(ApprovalStage, 'Etap przebiegu', ApprovalStage.Approved);
  if (etap === null) return;
  const uzasadnienie = zapytaj('Uzasadnienie zmiany etapu');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateApprovalSet, {
    panelId: idPanelu,
    stage: etap,
    note: uzasadnienie === '' ? undefined : uzasadnienie,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zmiany etapu przebiegu.');
    return;
  }
  przyjmijPanele(stan, [odpowiedz.wynik.panel]);
  powiedz(`Panel stoi na etapie ${odpowiedz.wynik.record.stage}.`);
  await stan.odswiez();
}

async function zapiszProfil(stan: StanTlumaczenia): Promise<void> {
  const nazwa = zapytaj('Nazwa profilu kontroli jakości');
  if (nazwa === '') return;
  const rodzaje = zapytajOWartosci(TranslationIssueKind, 'Rodzaje kontroli',
    [TranslationIssueKind.Placeholder]);
  if (rodzaje.length === 0) return;
  const waga = zapytajOWartosc(ProofreadSeverity, 'Waga zastrzeżenia', ProofreadSeverity.Warning);
  if (waga === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateQaProfileSet, {
    name: nazwa,
    checks: rodzaje.map((rodzaj) => ({ kind: rodzaj, severity: waga, enabled: true })),
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu profilu kontroli jakości.');
    return;
  }
  powiedz(`Profil „${odpowiedz.wynik.profile.name}” zapisany.`);
  await stan.odswiez();
}

/* Wykaz profilu w oknie niesie nazwę, a usunięcie żąda identyfikatora —
   nazwa idzie więc przez wykaz rdzenia, który identyfikator rozstrzyga. */
async function usunProfil(stan: StanTlumaczenia): Promise<void> {
  const nazwa = zapytaj('Nazwa profilu do usunięcia');
  if (nazwa === '') return;
  const wykaz = await wywolaj(stan.kanal, Command.TranslateQaProfileList, {});
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz profilów nie doszedł, więc nie ma czego usunąć.');
    return;
  }
  const profil = wykaz.wynik.profiles.find((kandydat) => kandydat.name === nazwa);
  if (profil === undefined) {
    odmowa(undefined, `Rdzeń nie ma profilu o nazwie „${nazwa}”.`);
    return;
  }
  if (!potwierdzNieodwracalna('profil-jakosci-usun',
    `Profil „${nazwa}” zejdzie z rdzenia bezpowrotnie.`)) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateQaProfileDelete, {
    profileId: profil.id,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił usunięcia profilu.');
    return;
  }
  powiedz(odpowiedz.wynik.deleted ? 'Profil usunięty.' : 'Profilu już nie było.');
  await stan.odswiez();
}
