/**
 * Źródło i segmentacja okna Translate: wczytanie tekstu, rozpoznanie języka,
 * podział na segmenty, scalenie i podział segmentu, zapis tłumaczenia panelu,
 * ton panelu i zestaw reguł segmentacji.
 */

import { Command } from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  odmowa,
  oknoWskazane,
  panelDocelowy,
  pokazWynik,
  powiedz,
  przyjmijPanele,
  zapytaj,
  zapytajOLiczbe,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';

const PANEL = 'panel-obszar-tlum';

export function czynnosciZrodla(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  return {
    'zrodlo-wczytaj': () => wczytajZrodlo(stan),
    'zrodlo-jezyk': () => rozpoznajJezyk(stan),
    'zrodlo-segmentuj': () => podzielNaSegmenty(stan),
    'segment-scal': () => scalSegmenty(stan),
    'segment-podziel': () => podzielSegment(stan),
    'tlumaczenie-zapisz': () => zapiszTlumaczenie(stan),
    'panel-ton': () => nastawTon(stan),
    'reguly-segmentacji': () => zapiszRegulySegmentacji(stan),
  };
}

async function wczytajZrodlo(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const tekst = zapytaj('Tekst źródłowy');
  if (tekst === '') return;
  const jezyk = zapytaj('Język źródłowy (pusty oznacza rozpoznanie automatyczne)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSourceSet, {
    windowId: idOkna,
    text: tekst,
    sourceLanguage: jezyk === '' ? undefined : jezyk,
    resegment: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania tekstu źródłowego.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  powiedz(`Tekst źródłowy w języku ${odpowiedz.wynik.sourceLanguage}; `
    + `segmentów: ${odpowiedz.wynik.segmentCount ?? 0}.`);
  await stan.odswiez();
}

async function rozpoznajJezyk(stan: StanTlumaczenia): Promise<void> {
  const tekst = zapytaj('Tekst do rozpoznania (pusty bierze bieżące źródło)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSourceDetect, {
    text: tekst === '' ? undefined : tekst,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie rozpoznał języka tekstu źródłowego.');
    return;
  }
  const pewnosc = odpowiedz.wynik.confidence;
  powiedz(`Rozpoznany język: ${odpowiedz.wynik.language}`
    + (pewnosc === undefined ? '.' : ` (pewność ${Math.round(pewnosc * 100)}%).`));
}

async function podzielNaSegmenty(stan: StanTlumaczenia): Promise<void> {
  const zestaw = zapytaj('Zestaw reguł segmentacji (pusty bierze podział domyślny)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSourceSegment, {
    rulesetId: zestaw === '' ? undefined : zestaw,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił podziału tekstu na segmenty.');
    return;
  }
  pokazSegmenty(stan, odpowiedz.wynik.segments);
}

async function scalSegmenty(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const wpis = zapytaj('Numery scalanych segmentów, kolejne i rosnące, po przecinku');
  const numery = wpis.split(',').map((czesc) => czesc.trim()).filter((czesc) => czesc !== '');
  if (numery.length < 2) {
    odmowa(undefined, 'Scalenie obejmuje co najmniej dwa segmenty.');
    return;
  }
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSegmentMerge, {
    windowId: idOkna,
    segmentIndexes: numery,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił scalenia segmentów.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  pokazSegmenty(stan, odpowiedz.wynik.segments);
}

async function podzielSegment(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const numer = zapytajOLiczbe('Numer dzielonego segmentu');
  if (numer === null) return;
  const przesuniecie = zapytajOLiczbe('Pozycja podziału w znakach od początku segmentu');
  if (przesuniecie === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSegmentSplit, {
    windowId: idOkna,
    segmentIndex: numer,
    offset: przesuniecie,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił podziału segmentu.');
    return;
  }
  przyjmijPanele(stan, odpowiedz.wynik.panels);
  pokazSegmenty(stan, odpowiedz.wynik.segments);
}

async function zapiszTlumaczenie(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const znany = stan.panele.find((panel) => panel.id === idPanelu);
  const tresc = zapytaj('Treść panelu po korekcie', znany?.text ?? '');
  if (tresc === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateTranslationSet, {
    panelId: idPanelu,
    text: tresc,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu tłumaczenia.');
    return;
  }
  przyjmijPanele(stan, [odpowiedz.wynik.panel]);
  powiedz('Tłumaczenie panelu zapisane.');
  await stan.odswiez();
}

async function nastawTon(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const ton = zapytaj('Ton tłumaczenia panelu');
  if (ton === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslatePanelToneSet, {
    panelId: idPanelu,
    tone: ton,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zmiany tonu panelu.');
    return;
  }
  przyjmijPanele(stan, [odpowiedz.wynik.panel]);
  powiedz(`Ton panelu to teraz ${ton}.`);
}

async function zapiszRegulySegmentacji(stan: StanTlumaczenia): Promise<void> {
  const nazwa = zapytaj('Nazwa zestawu reguł segmentacji');
  if (nazwa === '') return;
  const jezyk = zapytaj('Język zestawu (pusty znaczy zestaw wspólny)');
  const srx = zapytaj('Zestaw w postaci dokumentu SRX (pusty zakłada zestaw bez reguł)');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSegmentationRulesSet, {
    name: nazwa,
    language: jezyk === '' ? undefined : jezyk,
    srx: srx === '' ? undefined : srx,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu zestawu reguł segmentacji.');
    return;
  }
  powiedz(`Zestaw „${odpowiedz.wynik.ruleset.name}” zapisany.`);
  await stan.odswiez();
}

function pokazSegmenty(stan: StanTlumaczenia, segmenty: readonly string[]): void {
  pokazWynik(stan, PANEL, segmenty.map((tresc, numer) => [tresc, `segment ${numer + 1}`] as const),
    'Podział nie dał żadnego segmentu.');
  powiedz(`Segmentów: ${segmenty.length}.`);
}
