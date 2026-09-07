/**
 * Glosariusz okna Translate: nastawa terminu, wystąpienia w przekładach,
 * ujednolicenie, wciągnięcie i wydanie wykazu oraz wydobycie kandydatów
 * na terminy z materiału okna.
 */

import { Command, GlossaryTermStatus } from '../../../shared/contract.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import {
  odmowa,
  oknoWskazane,
  panelDocelowy,
  pokazWynik,
  powiedz,
  wydajSciezke,
  zapytaj,
  zapytajOWartosc,
  type CzynnosciTlumaczenia,
  type StanTlumaczenia,
} from './translate-wspolne.ts';

const PANEL = 'panel-glossary';

export function czynnosciGlosariusza(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  return {
    'glosariusz-zapisz': () => zapiszTermin(stan),
    'glosariusz-wystapienia': () => pokazWystapienia(stan),
    'glosariusz-ujednolic': () => ujednolicTerminy(stan),
    'glosariusz-wciagnij': () => wciagnijGlosariusz(stan),
    'glosariusz-wydaj': () => wydajGlosariusz(stan),
    'terminy-wydobadz': () => wydobadzKandydatow(stan),
  };
}

async function zapiszTermin(stan: StanTlumaczenia): Promise<void> {
  const zrodlowy = zapytaj('Termin źródłowy');
  if (zrodlowy === '') return;
  const jezyk = zapytaj('Język odpowiednika');
  if (jezyk === '') return;
  const docelowy = zapytaj('Odpowiednik docelowy (pusty przy terminie nietłumaczonym)');
  const status = zapytajOWartosc(GlossaryTermStatus, 'Status terminu', GlossaryTermStatus.Approved);
  if (status === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateGlossarySet, {
    source: zrodlowy,
    language: jezyk,
    target: docelowy === '' ? undefined : docelowy,
    doNotTranslate: docelowy === '' ? true : undefined,
    status,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu terminu.');
    return;
  }
  powiedz(`Termin „${odpowiedz.wynik.term.source}” zapisany.`);
  await stan.odswiez();
}

async function pokazWystapienia(stan: StanTlumaczenia): Promise<void> {
  const termin = zapytaj('Termin glosariusza');
  if (termin === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateGlossaryOccurrences, { term: termin });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał wystąpień terminu.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.occurrences.map((otoczenie) => [otoczenie, odpowiedz.wynik?.term ?? ''] as const),
    `Termin „${termin}” nie występuje w przekładach tego okna.`);
}

async function ujednolicTerminy(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = stan.panele.length > 0 ? panelDocelowy(stan) : '';
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateGlossaryApply, {
    panelId: idPanelu === '' ? undefined : idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił ujednolicenia terminologii.');
    return;
  }
  powiedz(`Ujednolicono ${odpowiedz.wynik.changedCount} wystąpień.`);
  await stan.odswiez();
}

async function wciagnijGlosariusz(stan: StanTlumaczenia): Promise<void> {
  const sciezka = zapytaj('Ścieżka pliku glosariusza po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateGlossaryImport, { path: sciezka });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania glosariusza.');
    return;
  }
  powiedz(`Wczytano ${odpowiedz.wynik.importedCount} terminów.`);
  await stan.odswiez();
}

/* Odpowiedź niesie samą liczbę terminów, więc do schowka idzie ścieżka
   wskazana rdzeniowi — innej rdzeń nie podaje. */
async function wydajGlosariusz(stan: StanTlumaczenia): Promise<void> {
  const sciezka = zapytaj('Ścieżka pliku wyniku po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateGlossaryExport, { path: sciezka });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wydania glosariusza.');
    return;
  }
  await wydajSciezke(sciezka, `Zapisano ${odpowiedz.wynik.exportedCount} terminów.`);
}

async function wydobadzKandydatow(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateTermExtract, {
    windowId: idOkna,
    excludeKnown: true,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń nie podał kandydatów na terminy.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.candidates.map((k) => [k.source, `${k.frequency} wystąpień`] as const),
    'Materiał okna nie dał kandydatów na terminy.');
}
