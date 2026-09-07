/**
 * Napisy i zasoby lokalizacyjne okna Translate: wciągnięcie i wydanie napisów,
 * kontrola taktowania, skrypt dubbingowy, wciągnięcie i wydanie zasobu,
 * formy mnogie oraz kontekst klucza.
 */

import { Command, SubtitleFormat } from '../../../shared/contract.ts';
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

const PANEL = 'panel-pliki';

export function czynnosciNapisow(stan: StanTlumaczenia): CzynnosciTlumaczenia {
  return {
    'napisy-wciagnij': () => wciagnijNapisy(stan),
    'napisy-wydaj': () => wydajNapisy(stan),
    'napisy-taktowanie': () => sprawdzTaktowanie(stan),
    'dubbing-skrypt': () => zlozSkryptDubbingu(stan),
    'zasob-wciagnij': () => wciagnijZasob(stan),
    'zasob-wydaj': () => wydajZasob(stan),
    'zasob-formy-mnogie': () => zastosujFormyMnogie(stan),
    'zasob-kontekst-klucza': () => opiszKontekstKlucza(stan),
  };
}

async function wciagnijNapisy(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const sciezka = zapytaj('Ścieżka pliku napisów po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSubtitleImport, {
    windowId: idOkna,
    path: sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania napisów.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.cues.map((k) => [k.text, `${k.index} · ${k.startMs}–${k.endMs} ms`] as const),
    'Plik napisów nie dał żadnej kwestii.');
  powiedz(`Wczytano ${odpowiedz.wynik.importedCount} kwestii.`);
}

async function wydajNapisy(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const sciezka = zapytaj('Ścieżka pliku wyniku po stronie rdzenia');
  if (sciezka === '') return;
  const format = zapytajOWartosc(SubtitleFormat, 'Format napisów', SubtitleFormat.Srt);
  if (format === null) return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSubtitleExport, {
    panelId: idPanelu,
    path: sciezka,
    format,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wydania napisów.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, `Zapisano ${odpowiedz.wynik.exportedCount} kwestii.`);
}

async function sprawdzTaktowanie(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateSubtitleTimingCheck, {
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił kontroli taktowania napisów.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.issues.map((z) => [`Kwestia ${z.index}: ${z.kind}`,
      `${z.value} przy dopuszczalnych ${z.limit}`] as const),
    'Taktowanie napisów nie budzi zastrzeżeń.');
}

async function zlozSkryptDubbingu(stan: StanTlumaczenia): Promise<void> {
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateDubbingScriptBuild, {
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił złożenia skryptu dubbingowego.');
    return;
  }
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.lines.map((w) => [w.text,
      `${w.speaker ?? 'bez mówcy'} · ${w.durationMs} ms`] as const),
    'Skrypt dubbingowy wyszedł pusty.');
  powiedz(`Kwestii dłuższych niż długość docelowa: ${odpowiedz.wynik.overLimitCount}.`);
}

async function wciagnijZasob(stan: StanTlumaczenia): Promise<void> {
  const idOkna = oknoWskazane(stan);
  if (idOkna === '') return;
  const sciezka = zapytaj('Ścieżka pliku zasobu po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateResourceImport, {
    windowId: idOkna,
    path: sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wczytania zasobu lokalizacyjnego.');
    return;
  }
  stan.idZasobu = odpowiedz.wynik.resource.id;
  pokazWynik(stan, PANEL,
    odpowiedz.wynik.keys.map((k) => [k.key, k.text] as const),
    'Zasób nie ma żadnego klucza.');
  powiedz(`Wczytano zasób o ${odpowiedz.wynik.resource.keyCount} kluczach.`);
}

async function wydajZasob(stan: StanTlumaczenia): Promise<void> {
  const idZasobu = zasobWskazany(stan);
  if (idZasobu === '') return;
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const sciezka = zapytaj('Ścieżka pliku wyniku po stronie rdzenia');
  if (sciezka === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateResourceExport, {
    resourceId: idZasobu,
    panelId: idPanelu,
    path: sciezka,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił wydania zasobu.');
    return;
  }
  await wydajSciezke(odpowiedz.wynik.path, `Zapisano ${odpowiedz.wynik.exportedCount} kluczy.`);
}

async function zastosujFormyMnogie(stan: StanTlumaczenia): Promise<void> {
  const idZasobu = zasobWskazany(stan);
  if (idZasobu === '') return;
  const idPanelu = panelDocelowy(stan);
  if (idPanelu === '') return;
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateResourcePluralApply, {
    resourceId: idZasobu,
    panelId: idPanelu,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zastosowania reguł liczby mnogiej.');
    return;
  }
  pokazWynik(stan, PANEL, odpowiedz.wynik.keys.map((k) => [k.key, k.text] as const),
    'Zasób nie ma klucza po zastosowaniu reguł.');
  powiedz(`Reguły zmieniły formy w ${odpowiedz.wynik.changedCount} kluczach.`);
}

async function opiszKontekstKlucza(stan: StanTlumaczenia): Promise<void> {
  const idZasobu = zasobWskazany(stan);
  if (idZasobu === '') return;
  const klucz = zapytaj('Klucz zasobu');
  if (klucz === '') return;
  const kontekst = zapytaj('Opis miejsca użycia klucza');
  const odpowiedz = await wywolaj(stan.kanal, Command.TranslateResourceKeyContextSet, {
    resourceId: idZasobu,
    key: klucz,
    context: kontekst === '' ? undefined : kontekst,
  });
  if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
    odmowa(odpowiedz.blad, 'Rdzeń odmówił zapisu kontekstu klucza.');
    return;
  }
  powiedz(`Kontekst klucza „${odpowiedz.wynik.key.key}” zapisany.`);
}

/* Kontrakt nie zna wykazu zasobów, więc zasób jest znany oknu dopiero po jego
   wczytaniu w tej karcie; bez tego czynność nazywa swój brak zamiast milczeć. */
function zasobWskazany(stan: StanTlumaczenia): string {
  if (stan.idZasobu === '') {
    odmowa(undefined, 'Okno nie zna żadnego zasobu lokalizacyjnego — wczytaj zasób.');
  }
  return stan.idZasobu;
}
