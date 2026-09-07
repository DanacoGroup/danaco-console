// Czynności okna Diagnostics: odczyt dziennika rdzenia i analiza błędów.
import { Command, LogLevel } from '../../../shared/contract.ts';
import type { LogEntry } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { dolozCzynnosciPanelu, zapytajWSzufladzie } from './czynnosci-okna.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Diagnostyka';

const POZIOMY: ReadonlyArray<readonly [string, string]> = [
  ['', 'Wszystkie poziomy'],
  [LogLevel.Error, 'Błędy'],
  [LogLevel.Warn, 'Ostrzeżenia'],
  [LogLevel.Info, 'Zdarzenia'],
  [LogLevel.Debug, 'Ślad'],
];

export function zwiazCzynnosciDiagnostyki(
  kanal: Kanal,
  korzen: Element,
  idOkna: () => string,
  odswiez: () => void,
  dolacz: (zdejmij: () => void) => void,
  przy: AddEventListenerOptions,
): void {
  const zdejmij = dolozCzynnosciPanelu(korzen, 'panel-logi', 'Czynności dziennika', [
    {
      naglowek: 'Dziennik rdzenia',
      pozycje: [{ kod: 'dziennik', nazwa: 'Odczytaj dziennik…' }],
    },
    {
      naglowek: 'Analiza',
      pozycje: [{ kod: 'analiza', nazwa: 'Uruchom analizę okna' }],
    },
  ], (kod) => {
    void wykonaj(kanal, korzen, kod, idOkna(), odswiez);
  }, przy);
  if (zdejmij !== null) dolacz(zdejmij);
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  kod: string,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (kod === 'dziennik') return odczytajDziennik(kanal, korzen);
  if (kod === 'analiza') return uruchomAnalize(kanal, idOkna, odswiez);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: pozycja bez gałęzi
     wyglądałaby jak działająca. */
  oglos(NAGLOWEK, `Czynność „${kod}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

/* Wynik przycięty granicą bufora jest nazwany: wykaz bez tej informacji
   wyglądałby na komplet dziennika. */
async function odczytajDziennik(kanal: Kanal, korzen: Element): Promise<void> {
  const wartosci = await zapytajWSzufladzie(korzen, 'panel-logi', {
    tytul: 'Odczyt dziennika',
    pola: [
      { klucz: 'wzorzec', etykieta: 'Wzorzec treści', podpowiedz: 'Puste bierze wszystkie wpisy' },
      { klucz: 'poziom', etykieta: 'Poziom wpisu', wybor: POZIOMY },
      { klucz: 'granica', etykieta: 'Najwyżej wpisów', wartosc: '100' },
    ],
    wykonanie: 'Odczytaj dziennik',
  });
  if (wartosci === null) return;
  const cialo = korzen.querySelector('#panel-logi .sta-okno-tresc');
  if (cialo === null) return;
  const granica = Number.parseInt(wartosci.granica ?? '', 10);
  const wynik = await wywolaj(kanal, Command.DiagnosticsLogQuery, {
    ...(wartosci.wzorzec === '' ? {} : { pattern: wartosci.wzorzec }),
    ...(wartosci.poziom === '' ? {} : { level: wartosci.poziom as LogLevel }),
    limit: Number.isFinite(granica) ? granica : 100,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu dziennika.', 'ostrzezenie');
    return;
  }
  const wpisy = wynik.wynik.entries;
  cialo.replaceChildren(...wpisy.map((wpis) => wierszWpisu(cialo, wpis)));
  if (wpisy.length === 0) {
    oglos(NAGLOWEK, 'Żaden wpis dziennika nie spełnia tych warunków.');
    return;
  }
  if (wynik.wynik.truncated === true) {
    oglos(NAGLOWEK, 'Wynik przycięty granicą bufora — dziennik ma więcej wpisów.', 'ostrzezenie');
  }
}

function wierszWpisu(cialo: Element, wpis: LogEntry): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'dn-wykaz-modulu-poz';
  const tresc = cialo.ownerDocument.createElement('span');
  tresc.textContent = wpis.message;
  const meta = cialo.ownerDocument.createElement('span');
  meta.className = 'dn-meta';
  meta.textContent = `${wpis.level} · ${wpis.source ?? 'rdzeń'}`;
  wiersz.append(tresc, meta);
  return wiersz;
}

/* Analiza obejmuje całe okno: zawężenia zakresem czasu ani wskazaniem błędów
   okno nie prowadzi, a zmyślenie granic byłoby decyzją za Operatora. */
async function uruchomAnalize(
  kanal: Kanal,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna diagnostyki dla tej karty.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.DiagnosticsAnalyzeRun, { windowId: idOkna });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił uruchomienia analizy.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, 'Analiza wykonana; rekomendacje przeliczone na nowo.');
  odswiez();
}
