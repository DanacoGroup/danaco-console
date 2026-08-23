import './diagnostics.css';

import { utworzPasPomocniczych } from '../../okna-pomocnicze/indeks';
import type { Kanal } from '../../protokol/kanal';
import { widokZMontazu, type OpisModulu } from '../rejestracja';
import { utworzOknoDiagnosticsCenter } from './okno-diagnostics-center';
import { utworzOknoErrorsPanel } from './okno-errors-panel';
import { utworzOknoLogsViewera } from './okno-logs-viewer';
import { utworzZrodloProwenancji } from './prowenancja-zrodlo';
import { utworzOknoObservabilityTools } from './okno-observability-tools';
import { utworzOknoRecommendationsPanel } from './okno-recommendations-panel';
import { utworzStanDiagnostyki, type StanDiagnostyki } from './stan-diagnostyki';
import { utworzZrodloDiagnostics } from './zrodlo-diagnostics';
import { utworzZrodloObserwowalnosci } from './zrodlo-obserwowalnosci';
import { utworzZrodloZuzycia } from './zuzycie-zrodlo';

/**
 * Moduł Diagnostics — złożenie pięciu okien wokół jednego stanu systemu.
 *
 * Port Diagnostyki jest w rdzeniu odbiorcą odmów wykonania komend
 * (`core/kompozycja.go`, pole `Diagnostyka`; mechanizm w
 * `core/adapter_modul_diagnostics_bledy.go`, `ZapiszNiepowodzenie`). Odmowa
 * dowolnej komendy staje się wierszem błędu, w którym `source` jest nazwą
 * komendy; wystąpienia grupują się po odcisku z licznikiem w `occurrences`.
 * Errors Panel jest jedynym miejscem w produkcie, gdzie widać odrzuconą komendę.
 *
 * Kod odmowy nie stoi w treści wiersza błędu: rdzeń wkłada `"<kod>: <treść>"`
 * do wpisu dziennika (czyta go Logs Viewer), a do wiersza błędu daje `message`
 * bez kodu — sam kod osobno, polem `errorCode` kontraktu albo, w dzisiejszym
 * rdzeniu, w `context.errorCode`. Po kodzie rozpoznaje się komendę bez
 * obsługiwacza (`*.unknown`) i odróżnia ją od odmowy merytorycznej, więc oba
 * miejsca muszą być czytane wprost, a nie zakładane.
 *
 * Moduł montuje się bez okna, bo trzy z czterech komend obszaru (`log.query`,
 * `error.list`, `recommendation.list`) nie mają pola `windowId` wcale, a czwarta
 * (`analyze.run`) ma je opcjonalne. Czekanie na `window.list` wzorem modułu
 * Developer zostawiłoby sesję bez okien również bez dziennika, bez wykazu
 * błędów i bez rekomendacji. Cena jest jedna i jawna: analiza nie zostaje
 * przypisana do okna, bo `windowId` idzie do rdzenia tylko wtedy, gdy okno jest
 * znane; kontrakt to dopuszcza, a rdzeń przyjmuje analizę bez okna.
 *
 * Układ wynika z ról. W pasie górnym wiodące Diagnostics Center wraz
 * z monitorem dziennika, bo Centrum agreguje to, co Logs Viewer pokazuje
 * wprost. W pasie dolnym dwa okna pomocnicze: Errors Panel zasilający Centrum
 * materiałem źródłowym i Recommendations Panel czytający owoc analizy Centrum.
 * Zależności biegną w obie strony, dlatego okna jadą jednym stanem: Centrum
 * wpisuje analizę i zakres czasu, Errors Panel oraz Logs Viewer biorą z niego
 * zakres, Recommendations Panel bierze analizę. Okno rozmowy modułu nie należy
 * do złożenia: jest bytem sesji i składa je warstwa rozmowy.
 *
 * W pasie trzecim stoi Observability Tools — kontener narzędzi warstwy
 * eksperckiej. Sięga po rodzinę `monitor.*`, której pozostałe okna modułu nie
 * dotykają, więc ma własne źródło; zakres czasu bierze ze stanu wspólnego, bo
 * jedna z jego zakładek czyta dziennik.
 *
 * Czwarty pas niesie okna pomocnicze: podgląd w tle powłoki bash na
 * `terminal.output.stream` oraz spis pozycji, których jeszcze nie ma, wraz
 * z powodem każdej (`okna-pomocnicze/rejestr-pomocniczych.ts`). Podgląd działa
 * tu słabiej niż w module Developer i mówi to wprost: moduł montuje się bez
 * okna, więc `idOkna` bywa puste, rdzeń odpowiada wtedy `subscribed: false`
 * i oddaje sam ogon historii, a nowe wiersze nie dochodzą na żywo.
 */
export interface ZamontowaneDiagnostics {
  /** Element osadzony w dokumencie. */
  element: HTMLElement;
  /** Stan diagnostyki wspólny czterem oknom. */
  stan: StanDiagnostyki;
  /** Odczytuje wszystkie okna z rdzenia. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń modułu wraz z nasłuchem okien. */
  zamknij(): void;
}

export function zamontujDiagnostics(
  gospodarz: HTMLElement,
  kanal: Kanal,
  idOkna = '',
): ZamontowaneDiagnostics {
  const zrodlo = utworzZrodloDiagnostics(kanal);
  const obserwowalnosc = utworzZrodloObserwowalnosci(kanal);
  const stan = utworzStanDiagnostyki(zrodlo);

  const obszar = document.createElement('div');
  obszar.className = 'dg-modul';
  obszar.dataset['modul'] = 'diagnostics';

  const centrum = utworzOknoDiagnosticsCenter(zrodlo, stan, idOkna);
  const dziennik = utworzOknoLogsViewera(zrodlo, stan);
  const bledy = utworzOknoErrorsPanel(zrodlo, stan);
  const rekomendacje = utworzOknoRecommendationsPanel(zrodlo, stan);
  const narzedzia = utworzOknoObservabilityTools(
    zrodlo,
    obserwowalnosc,
    stan,
    idOkna,
    utworzZrodloZuzycia(kanal),
    utworzZrodloProwenancji(kanal),
  );

  const gorny = document.createElement('div');
  gorny.className = 'dg-modul__gora';
  gorny.append(
    oznacz(centrum.element, 'diagnostics-center'),
    oznacz(dziennik.element, 'logs-viewer'),
  );

  const dolny = document.createElement('div');
  dolny.className = 'dg-modul__dol';
  dolny.append(
    oznacz(bledy.element, 'errors-panel'),
    oznacz(rekomendacje.element, 'recommendations-panel'),
  );

  const ekspercki = document.createElement('div');
  ekspercki.className = 'dg-modul__ekspercki';
  ekspercki.append(oznacz(narzedzia.element, 'observability-tools'));

  const pomocnicze = utworzPasPomocniczych({
    kanal,
    modul: KOD_MODULU,
    nazwaModulu: NAZWA_MODULU,
    okno: idOkna,
    przedrostek: 'dg',
  });

  obszar.append(gorny, dolny, ekspercki, pomocnicze.element);
  gospodarz.replaceChildren(obszar);

  function odswiez(): void {
    centrum.odswiez();
    dziennik.odswiez();
    bledy.odswiez();
    rekomendacje.odswiez();
    narzedzia.odswiez();
    pomocnicze.odswiez();
  }

  odswiez();

  return {
    element: obszar,
    stan,
    odswiez,
    zamknij() {
      // Pas zamyka się pierwszy: trzyma subskrypcję `stream.chunk`, która żyje
      // niezależnie od stanu modułu i po zejściu ze sceny nikt by jej nie zdjął.
      pomocnicze.zamknij();
      narzedzia.zamknij();
      rekomendacje.zamknij();
      bledy.zamknij();
      dziennik.zamknij();
      centrum.zamknij();
      stan.zamknij();
    },
  };
}

/** Znakuje okno kodem rejestru okien operacyjnych — po nim skacze nawigacja. */
function oznacz(element: HTMLElement, kod: string): HTMLElement {
  element.dataset['okno'] = kod;
  element.tabIndex = -1;
  return element;
}

/**
 * Samoopisujący się moduł dla rejestru powłoki.
 *
 * Kod siedzi w module, nie w mapie po stronie powłoki: dodanie modułu to jeden
 * wpis, a nie dwa, więc nie da się dodać modułu i zapomnieć o wytwórni.
 *
 * Kod jest stałą, bo czyta go także pas okien pomocniczych — po nim idzie spis
 * pozycji i profil rozmowy modułu. Dwa napisy `'diagnostics'` w jednym pliku
 * rozjechałyby się przy pierwszej zmianie nazwy.
 */
const KOD_MODULU = 'diagnostics';

/** Nazwa modułu w etykietach dostępności pasa okien pomocniczych. */
const NAZWA_MODULU = 'Diagnostics';

export const MODUL: OpisModulu = {
  kod: KOD_MODULU,
  utworzWidok: (kanal) => widokZMontazu(zamontujDiagnostics, kanal),
};
