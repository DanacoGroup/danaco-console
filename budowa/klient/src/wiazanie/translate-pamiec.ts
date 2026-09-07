/**
 * Pamięć tłumaczeń okna Translate: podpowiedzi, tłumaczenie wstępne, zapis
 * i usunięcie pary, wciągnięcie i wydanie pamięci, wyrównanie tekstów,
 * porządkowanie wsadowe oraz polityka pamięci okna.
 */

import { Command, TranslationMemoryMaintenanceKind } from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  odmowa,
  oknoWskazane,
  panelDocelowy,
  pokazWynik,
  potwierdzNieodwracalna,
  powiedz,
  przyjmijPanele,
  wydajSciezke,
  zapytaj,
  zapytajOLiczbe,
  zapytajOWartosc,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';

const PANEL = 'panel-obszar-tlum';

export function czynnosciPamieci(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  return {
    'pamiec-podpowiedz': () => podpowiedzZPamieci(stan),
    'pamiec-wstepne': () => przetlumaczWstepnie(stan),
    'pamiec-zapisz': () => zapiszPare(stan),
    'pamiec-usun': () => usunPare(stan),
    'pamiec-wciagnij': () => wciagnijPamiec(stan),
    'pamiec-wydaj': () => wydajPamiec(stan),
    'pamiec-wyrownaj': () => wyrownajTeksty(stan),
    'pamiec-uporzadkuj': () => uporzadkujPamiec(stan),
    'pamiec-polityka-odczyt': () => odczytajPolityke(stan),
    'pamiec-polityka-zapis': () => zapiszPolityke(stan),
  };
}

async function podpowiedzZPamieci(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const segment = zapytaj('Segment źródłowy, do którego szukać podpowiedzi');
  if (segment === '') return;
  const prog = zapytajOLiczbe('Próg dopasowania w procentach (pusty bierze próg polityki)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemorySuggest, {
    panelId: idPanelu,
    segment,
    threshold: prog ?? undefined,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał podpowiedzi z pamięci.');
    return;
  }
  const trafienia = odpowiedz.wynik.matches ?? [];
  const pozycje = trafienia.length > 0
    ? trafienia.map((t) => [t.entry.targetSegment,
      `${t.score}%${t.contextMatch ? ' w kontekście' : ''}`] as const)
    : odpowiedz.wynik.suggestions.map((tresc) => [tresc, ''] as const);
  pokazWynik(stan, PANEL, pozycje, 'Pamięć nie ma dopasowania do tego segmentu.');
}

async function przetlumaczWstepnie(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryPretranslate, {
    windowId: idOkna,
    onlyEmpty: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił tłumaczenia wstępnego z pamięci.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  powiedz(`Z pamięci wypełniono ${odpowiedz.wynik.filledCount} segmentów.`);
  await stan.odswiez();
}

async function zapiszPare(stan: StanTlumaczenia): Promise<void> {
  const jezyk = zapytaj('Język segmentu docelowego');
  if (jezyk === '') return;
  const zrodlowy = zapytaj('Segment źródłowy');
  if (zrodlowy === '') return;
  const docelowy = zapytaj('Segment docelowy');
  if (docelowy === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemorySet, {
    language: jezyk,
    sourceSegment: zrodlowy,
    targetSegment: docelowy,
  });
  if (!odpowiedz.udany) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu pary w pamięci.');
    return;
  }
  powiedz('Para zapisana w pamięci tłumaczeń.');
  await stan.odswiez();
}

/* Wykaz pamięci nie niesie w oknie identyfikatorów par, a usunięcie żąda
   identyfikatora — parę wskazuje więc segment, a wykaz rdzenia ją rozstrzyga. */
async function usunPare(stan: StanTlumaczenia): Promise<void> {
  const szukany = zapytaj('Segment źródłowy pary do usunięcia');
  if (szukany === '') return;
  const wykaz = await wywolaj(stan.kanal, Command.TranslateMemoryList, {
    query: szukany,
    limit: 2,
  });
  if (!wykaz.udany || wykaz.wynik === undefined) {
    odmowa(wykaz.blad, 'Wykaz pamięci nie doszedł, więc nie ma czego usunąć.');
    return;
  }
  const znalezione = wykaz.wynik.entries;
  const para = znalezione[0];
  if (para === undefined) {
    odmowa(undefined, 'Pamięć nie ma pary o takim segmencie źródłowym.');
    return;
  }
  if (znalezione.length > 1) {
    odmowa(undefined, 'Segment wskazuje więcej niż jedną parę — podaj go dokładniej.');
    return;
  }
  if (!potwierdzNieodwracalna('pamiec-usun',
    `Para „${para.sourceSegment}” zejdzie z pamięci bezpowrotnie.`)) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryDelete, { entryId: para.id });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił usunięcia pary.');
    return;
  }
  powiedz(odpowiedz.wynik.deleted ? 'Para usunięta z pamięci.' : 'Pary już w pamięci nie było.');
  await stan.odswiez();
}

async function wciagnijPamiec(stan: StanTlumaczenia): Promise<void> {
  const sciezka = zapytaj('Ścieżka pliku TMX po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryImport, { path: sciezka });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania pamięci.');
    return;
  }
  powiedz(`Wczytano ${odpowiedz.wynik.importedCount} par, `
    + `pominięto ${odpowiedz.wynik.skippedCount} powtórzeń.`);
  await stan.odswiez();
}

async function wydajPamiec(stan: StanTlumaczenia): Promise<void> {
  const sciezka = zapytaj('Ścieżka pliku wyniku po stronie rdzenia');
  if (sciezka === '') return;
  const jezyk = zapytaj('Język par (pusty obejmuje wszystkie języki)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryExport, {
    path: sciezka,
    language: jezyk === '' ? undefined : jezyk,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wydania pamięci.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, `Zapisano ${odpowiedz.wynik.exportedCount} par.`);
}

async function wyrownajTeksty(stan: StanTlumaczenia): Promise<void> {
  const zrodlowy = zapytaj('Tekst źródłowy');
  if (zrodlowy === '') return;
  const docelowy = zapytaj('Gotowy przekład tego samego materiału');
  if (docelowy === '') return;
  const jezyk = zapytaj('Język przekładu');
  if (jezyk === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryAlign, {
    sourceText: zrodlowy,
    targetText: docelowy,
    language: jezyk,
    commit: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wyrównania tekstów.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.pairs.map((p) => [p.targetSegment, `${p.score}% · ${p.sourceSegment}`] as const),
    'Wyrównanie nie dało żadnej pary.');
  powiedz(`W pamięci zapisano ${odpowiedz.wynik.committedCount} par.`);
}

/* Porządkowanie idzie najpierw przebiegiem próbnym: liczba objętych par mówi
   Operatorowi, na co się zgadza, zanim pamięć zostanie zmieniona. */
async function uporzadkujPamiec(stan: StanTlumaczenia): Promise<void> {
  const rodzaj = zapytajOWartosc(TranslationMemoryMaintenanceKind, 'Rodzaj operacji',
    TranslationMemoryMaintenanceKind.Deduplicate);
  if (rodzaj === null) return;
  const szukany = zapytaj('Szukany fragment (wymagany przy podmianie i filtrowaniu)');
  const wstawiany = rodzaj === TranslationMemoryMaintenanceKind.Replace
    ? zapytaj('Treść wstawiana w miejsce szukanego fragmentu')
    : '';
  const zadanie = {
    kind: rodzaj,
    search: szukany === '' ? undefined : szukany,
    replacement: wstawiany === '' ? undefined : wstawiany,
  };
  const proba = await wywolaj(stan.kanal, Command.TranslateMemoryMaintain,
    { ...zadanie, dryRun: true });
  if (!proba.udany || proba.wynik === undefined) {
    odmowa(proba.blad, 'Rdzeń odmówił przebiegu próbnego.');
    return;
  }
  if (!potwierdzNieodwracalna('pamiec-uporzadkuj',
    `Operacja obejmie ${proba.wynik.affectedCount} par pamięci.`)) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryMaintain,
    { ...zadanie, dryRun: false });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił porządkowania pamięci.');
    return;
  }
  powiedz(`Zmieniono albo usunięto ${odpowiedz.wynik.affectedCount} par.`);
  await stan.odswiez();
}

async function odczytajPolityke(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryPolicyGet, { windowId: idOkna });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał polityki pamięci.');
    return;
  }
  const polityka = odpowiedz.wynik.policy;
  powiedz(`Zasięg ${polityka.scope}, próg ${polityka.threshold}%, `
    + `kontekst ${polityka.contextMatch ? 'odróżniany' : 'pomijany'}, `
    + `tłumaczenie wstępne ${polityka.preTranslate ? 'włączone' : 'wyłączone'}.`);
}

async function zapiszPolityke(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const prog = zapytajOLiczbe('Próg dopasowania przybliżonego w procentach');
  if (prog === null) return;
  const wstepne = zapytaj('Tłumaczyć wstępnie z pamięci? (tak/nie)', 'tak') === 'tak';
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateMemoryPolicySet, {
    windowId: idOkna,
    threshold: prog,
    preTranslate: wstepne,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu polityki pamięci.');
    return;
  }
  powiedz(`Polityka pamięci zapisana z progiem ${odpowiedz.wynik.policy.threshold}%.`);
}
