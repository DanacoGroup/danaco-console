/**
 * Wymiana okna Translate: dokumenty, XLIFF, pakiet dla wykonawcy, most do
 * modułu Studio, artefakt biblioteki, wydanie panelu, synteza mowy, zlecenie
 * wsadowe, profile przekładu i polityka języka pośredniego.
 */

import {
  BatchOperationKind,
  BridgeResultMode,
  Command,
  ExportFormat,
  HandoffContent,
  TranslationArtifactKind,
  TranslationExportVariant,
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
  wydajSciezke,
  zapytaj,
  zapytajOWartosc,
  zapytajOWartosci,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';

const PANEL_PLIKI = 'panel-pliki';
const PANEL_KOLEJKA = 'panel-kolejka';

export function czynnosciWymiany(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  return {
    'dokument-wczytaj': () => wczytajDokument(stan),
    'dokument-zloz': () => zlozDokument(stan),
    'dokument-uklad': () => porownajUklad(stan),
    'xliff-wciagnij': () => wciagnijXliff(stan),
    'pakiet-zloz': () => zlozPakiet(stan),
    'pakiet-przyjmij': () => przyjmijPakiet(stan),
    'most-przyjmij-zrodlo': () => przyjmijZrodloZeStudia(stan),
    'most-odeslij-wynik': () => odeslijWynikDoStudia(stan),
    'artefakt-opublikuj': () => opublikujArtefakt(stan),
    'panel-wydaj': () => wydajPanel(stan),
    'mowa-synteza': () => zsyntezujMowe(stan),
    'zlecenie-wsadowe': () => uruchomZlecenieWsadowe(stan),
    'profile-przekladu': () => pokazProfilePrzekladu(stan),
    'profil-przekladu-zapisz': () => zapiszProfilPrzekladu(stan),
    'przeklady-porownaj': () => porownajPrzeklady(stan),
    'pivot-odczyt': () => odczytajPolitykePivota(stan),
    'pivot-zapis': () => zapiszPolitykePivota(stan),
  };
}

async function wczytajDokument(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const sciezka = zapytaj('Ścieżka dokumentu po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateDocumentLoad, {
    windowId: idOkna,
    path: sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania dokumentu.');
    return;
  }
  stan.idDokumentu = odpowiedz.wynik.document.id;
  pokazWynik(stan, PANEL_PLIKI,
    odpowiedz.wynik.segments.map((s) => [s.text, `${s.index}`] as const),
    'Dokument nie dał żadnego segmentu.');
  powiedz(`Wczytano dokument w formacie ${odpowiedz.wynik.document.format}`
    + `${odpowiedz.wynik.usedOcr ? ' z rozpoznaniem pisma' : ''}.`);
}

async function zlozDokument(stan: StanTlumaczenia): Promise<void> {
  const idDokumentu = dokumentWskazany(stan);
  if (idDokumentu === '') return;
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const sciezka = zapytaj('Ścieżka pliku wyniku (pusta odkłada wynik w magazynie rdzenia)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateDocumentRender, {
    documentId: idDokumentu,
    panelId: idPanelu,
    path: sciezka === '' ? undefined : sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił złożenia dokumentu.');
    return;
  }
  const opis = `Wstawiono ${odpowiedz.wynik.renderedCount} segmentów.`;
  if (odpowiedz.wynik.path === '') {
    powiedz(`${opis} Wynik stoi w magazynie rdzenia jako zasób ${odpowiedz.wynik.assetId ?? ''}.`);
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, opis);
}

async function porownajUklad(stan: StanTlumaczenia): Promise<void> {
  const idDokumentu = dokumentWskazany(stan);
  if (idDokumentu === '') return;
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateDocumentLayoutCompare, {
    documentId: idDokumentu,
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił porównania układu dokumentu.');
    return;
  }
  pokazWynik(stan, PANEL_PLIKI,
    odpowiedz.wynik.differences.map((r) => [r.detail, `strona ${r.page} · ${r.severity}`] as const),
    'Układ dokumentu zachowany.');
  powiedz(`Porównano ${odpowiedz.wynik.comparedPages} stron.`);
}

async function wciagnijXliff(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const sciezka = zapytaj('Ścieżka pliku XLIFF po stronie rdzenia');
  if (sciezka === '') return;
  if (!potwierdzNieodwracalna('xliff-wciagnij',
    'Wczytany przekład nadpisze treść paneli okna.')) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateXliffImport, {
    windowId: idOkna,
    path: sciezka,
    overwrite: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania pliku XLIFF.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  powiedz(`Wczytano ${odpowiedz.wynik.importedCount} jednostek; `
    + `notatek wykonawcy: ${odpowiedz.wynik.notes.length}.`);
  await stan.odswiez();
}

async function zlozPakiet(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const skladniki = zapytajOWartosci(HandoffContent, 'Składniki pakietu', [HandoffContent.Xliff]);
  if (skladniki.length === 0) return;
  const sciezka = zapytaj('Ścieżka pliku pakietu (pusta odkłada pakiet w magazynie rdzenia)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateHandoffBuild, {
    windowId: idOkna,
    contents: skladniki,
    path: sciezka === '' ? undefined : sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił złożenia pakietu dla wykonawcy.');
    return;
  }
  stan.idPakietu = odpowiedz.wynik.package.id;
  if (odpowiedz.wynik.path === '') {
    powiedz('Pakiet stoi w magazynie rdzenia; pliku nie wytworzono.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, 'Pakiet dla wykonawcy złożony.');
}

async function przyjmijPakiet(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const sciezka = zapytaj('Ścieżka pliku zwrotu po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateHandoffReceive, {
    windowId: idOkna,
    path: sciezka,
    packageId: stan.idPakietu === '' ? undefined : stan.idPakietu,
    runQualityCheck: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił przyjęcia zwrotu od wykonawcy.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  powiedz(`Przyjęto ${odpowiedz.wynik.importedCount} jednostek tłumaczeniowych.`);
  await stan.odswiez();
}

async function przyjmijZrodloZeStudia(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const idDokumentu = zapytaj('Dokument modułu Studio, z którego pochodzi zaznaczenie');
  if (idDokumentu === '') return;
  const tresc = zapytaj('Treść zaznaczenia');
  if (tresc === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateBridgeSourceReceive, {
    windowId: idOkna,
    documentId: idDokumentu,
    text: tresc,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił przyjęcia zaznaczenia ze Studia.');
    return;
  }
  powiedz(`Przyjęto zaznaczenie podzielone na ${odpowiedz.wynik.segmentCount} pozycji.`);
  await stan.odswiez();
}

async function odeslijWynikDoStudia(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const postac = zapytajOWartosc(BridgeResultMode, 'Postać wyniku', BridgeResultMode.TargetOnly);
  if (postac === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateBridgeResultSend, {
    windowId: idOkna,
    panelIds: [idPanelu],
    mode: postac,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił odesłania wyniku do Studia.');
    return;
  }
  powiedz(`Odesłano ${odpowiedz.wynik.sentCount} paneli do dokumentu `
    + `${odpowiedz.wynik.documentId}.`);
}

async function opublikujArtefakt(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const rodzaj = zapytajOWartosc(TranslationArtifactKind, 'Rodzaj artefaktu',
    TranslationArtifactKind.TargetFile);
  if (rodzaj === null) return;
  const zPrzekladu = rodzaj === TranslationArtifactKind.TargetFile
    || rodzaj === TranslationArtifactKind.BilingualFile;
  const idPanelu = zPrzekladu ? panelDocelowy(stan) : '';
  if (zPrzekladu && idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateArtifactPublish, {
    windowId: idOkna,
    kind: rodzaj,
    panelIds: idPanelu === '' ? undefined : [idPanelu],
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił opublikowania artefaktu.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path,
    `Artefakt stoi w bibliotece jako plik ${odpowiedz.wynik.fileId}.`);
}

async function wydajPanel(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const format = zapytajOWartosc(ExportFormat, 'Format wydania', ExportFormat.Docx);
  if (format === null) return;
  const postac = zapytajOWartosc(TranslationExportVariant, 'Postać pliku',
    TranslationExportVariant.Target);
  if (postac === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslatePanelExport, {
    panelId: idPanelu,
    format,
    variant: postac,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wydania panelu.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, `Panel wydany w postaci ${odpowiedz.wynik.variant}.`);
}

async function zsyntezujMowe(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSpeechSynthesize, {
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił syntezy mowy dla panelu.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, 'Nagranie odsłuchu gotowe.');
}

async function uruchomZlecenieWsadowe(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const operacje = zapytajOWartosci(BatchOperationKind, 'Operacje zlecenia',
    [BatchOperationKind.Translate, BatchOperationKind.QualityCheck]);
  if (operacje.length === 0) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateBatchRun, {
    windowId: idOkna,
    operations: operacje,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił uruchomienia zlecenia wsadowego.');
    return;
  }
  powiedz(`Kolejka ${odpowiedz.wynik.queueId} przyjęła ${odpowiedz.wynik.itemCount} pozycji.`);
  await stan.odswiez();
}

async function pokazProfilePrzekladu(stan: StanTlumaczenia): Promise<void> {
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateEngineProfileList, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał profilów przekładu.');
    return;
  }
  pokazWynik(stan, PANEL_KOLEJKA,
    odpowiedz.wynik.profiles.map((p) => [p.name, `${p.scope} · ${p.domain ?? 'bez dziedziny'}`] as const),
    'Rdzeń nie ma profilu przekładu.');
}

async function zapiszProfilPrzekladu(stan: StanTlumaczenia): Promise<void> {
  const nazwa = zapytaj('Nazwa profilu przekładu');
  if (nazwa === '') return;
  const dziedzina = zapytaj('Dziedzina materiału');
  const dostrojony = zapytaj('Dostrajać przekład do zatwierdzonej pamięci? (tak/nie)', 'tak') === 'tak';
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateEngineProfileSet, {
    name: nazwa,
    domain: dziedzina === '' ? undefined : dziedzina,
    adaptive: dostrojony,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu profilu przekładu.');
    return;
  }
  powiedz(`Profil przekładu „${odpowiedz.wynik.profile.name}” zapisany.`);
}

/* Porównanie żąda wskazania kanałów modelu, a okno nie ma ich wykazu własnego:
   bierze kanały czynne rdzenia, bo tylko na nich przekład w ogóle zejdzie. */
async function porownajPrzeklady(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const segment = zapytaj('Segment źródłowy do przekładu');
  if (segment === '') return;
  const kanaly = await wywolaj(stan.kanal, Command.ChannelList, { enabledOnly: true });
  const wskazania = kanaly.wynik?.channels.map((k) => k.id) ?? [];
  if (wskazania.length === 0) {
    odmowa(kanaly.blad, 'Rdzeń nie ma czynnego kanału modelu, więc nie ma czego porównać.');
    return;
  }
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateEngineCompare, {
    panelId: idPanelu,
    segment,
    channelIds: wskazania,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił porównania przekładów.');
    return;
  }
  pokazWynik(stan, PANEL_KOLEJKA,
    odpowiedz.wynik.variants.map((w) => [w.text, w.channelId] as const),
    'Żaden kanał nie zwrócił przekładu.');
}

async function odczytajPolitykePivota(stan: StanTlumaczenia): Promise<void> {
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslatePivotPolicyGet, {});
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał polityki języka pośredniego.');
    return;
  }
  const polityka = odpowiedz.wynik.policy;
  pokazWynik(stan, PANEL_KOLEJKA,
    polityka.pairs.map((p) => [`${p.sourceLanguage} → ${p.targetLanguage}`, p.pivotLanguage] as const),
    'Polityka nie ma pary z własnym językiem pośrednim.');
  powiedz(`Domyślny język pośredni: ${polityka.defaultPivot ?? 'brak, przekład bezpośredni'}.`);
}

async function zapiszPolitykePivota(stan: StanTlumaczenia): Promise<void> {
  const jezyk = zapytaj('Domyślny język pośredni (pusty znaczy przekład bezpośredni)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslatePivotPolicySet, {
    defaultPivot: jezyk === '' ? undefined : jezyk,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu polityki języka pośredniego.');
    return;
  }
  powiedz('Polityka języka pośredniego zapisana.');
}

/* Kontrakt nie zna wykazu dokumentów, więc dokument jest znany oknu dopiero
   po jego wczytaniu w tej karcie. */
function dokumentWskazany(stan: StanTlumaczenia): string {
  if (stan.idDokumentu === '') {
    odmowa(undefined, 'Okno nie zna żadnego dokumentu — wczytaj dokument.');
  }
  return stan.idDokumentu;
}
