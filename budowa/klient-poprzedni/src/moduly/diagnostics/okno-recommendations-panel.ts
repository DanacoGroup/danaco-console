import {
  DiagnosticPriority,
  RecommendationStatus,
  type DiagnosticRecommendation,
  type DiagnosticsRecommendationListRequest,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  poleLiczbowe,
  przyciskAkcji,
  przyciskBezKomendy,
  pozycjaWykazu,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza-braki';
import { utworzWyborZMenu, type PozycjaWyboruMenu, type WyborZMenu } from '../apps/wybor-z-menu';
import { powodBezKomendy } from './braki-kontraktu';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Recommendations Panel — okno pomocnicze modułu Diagnostics: przegląd
 * rekomendacji powstałych z analizy; wdrożenie poprawki idzie przez Code
 * Editor.
 */
export interface OknoRecommendationsPanel {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch `stan.naZmiane` założony przy montażu okna. */
  zamknij(): void;
}

export function utworzOknoRecommendationsPanel(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
): OknoRecommendationsPanel {
  const rama = utworzRameOkna({
    tytul: 'Recommendations Panel',
    rola: 'pomocnicze',
    kod: 'recommendations-panel',
    przeznaczenie:
      'Przegląd rekomendacji powstałych z analizy Diagnostics Center; wdrożenie poprawki idzie przez Code Editor.',
    przedrostek: 'dg',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieRekomendacji(rama, tresc.element);
  // Rekomendacje wraz z id analizy, do której należą — oba pola zmieniają się tylko razem.
  let rekomendacje: readonly DiagnosticRecommendation[] = [];
  let analizaWykazu = '';

  function odczytaj(): void {
    const idAnalizy = stan.analiza();
    if (idAnalizy === '') {
      // Bez analizy nie ma czyjego wykazu trzymać, więc stary wykaz tu nie zostaje.
      rekomendacje = [];
      analizaWykazu = '';
      tresc.pusto('Nie uruchomiono jeszcze analizy — uruchom ją w Diagnostics Center.');
      return;
    }
    tresc.ladowanie('Odczyt rekomendacji analizy…');
    const zadanie = zadanieRekomendacji(idAnalizy, powierzchnia);
    void zrodlo.wykazRekomendacji(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odczyt nieudany — stary wykaz innej analizy nie zostaje do eksportu pod nową nazwą.
        rekomendacje = [];
        analizaWykazu = '';
        tresc.blad(zdanieNiepowodzenia(`odczytu rekomendacji analizy ${idAnalizy}`, wynik.powod), wynik.blad);
        return;
      }
      rekomendacje = wynik.wynik;
      analizaWykazu = idAnalizy;
      rysujRekomendacje(rekomendacje, idAnalizy, tresc);
    });
  }

  function raportuj(): void {
    if (rekomendacje.length === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — wykaz rekomendacji jest pusty.', false);
      return;
    }
    pobierzPlik(`${analizaWykazu}-raport.md`, raportRekomendacji(rekomendacje), 'text/markdown');
    tresc.potwierdzenie('Raport rekomendacji pobrany jako plik Markdown.', true);
  }

  podepnijAkcjeRekomendacji(powierzchnia, { odczytaj, raportuj });

  // Jeden punkt wejścia to odswiez — wytwórnia okna nie czyta sama, woła ją złożenie modułu.
  const odsubskrybuj = stan.naZmiane(odczytaj);

  return { element: rama.element, odswiez: odczytaj, zamknij: odsubskrybuj };
}

/** Filtry i pasek akcji okna panelu rekomendacji modułu Diagnostics, osadzane w ramie tego całego okna. */
interface PowierzchniaRekomendacji {
  status: WyborZMenu;
  priorytet: WyborZMenu;
  granica: HTMLInputElement;
  odswiezPrzycisk: HTMLButtonElement;
  raport: HTMLButtonElement;
}

/** Opcje filtra złożone z wyliczenia kontraktu wraz z pozycją „wszystkie”, dla kontrolki wyboru w oknie. */
function opcjeFiltra(
  etykietaWszystkich: string,
  wyliczenie: Record<string, string>,
): readonly PozycjaWyboruMenu[] {
  return [
    { wartosc: '', etykieta: etykietaWszystkich },
    ...Object.values(wyliczenie).map((w) => ({ wartosc: w, etykieta: w })),
  ];
}

/** Składa filtry (status, priorytet, granica) i panel akcji ramy okna panelu rekomendacji tego całego modułu. */
function zlozPowierzchnieRekomendacji(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaRekomendacji {
  // Rozwijanie z biblioteki, nie natywny select, przez wspólną obsadę wyboru z menu.
  const status = utworzWyborZMenu(
    'Stan rekomendacji',
    opcjeFiltra('Wszystkie stany', RecommendationStatus),
  );
  const priorytet = utworzWyborZMenu(
    'Priorytet rekomendacji',
    opcjeFiltra('Wszystkie priorytety', DiagnosticPriority),
  );
  const granica = poleLiczbowe('Górna granica liczby rekomendacji', 'domyślnie wszystkie');
  rama.narzedzia.append(
    wiersz('Stan', status.element, { klasa: 'dg-wiersz' }),
    wiersz('Priorytet', priorytet.element, {
      klasa: 'dg-wiersz',
      objasnienie: 'Filtr priorytetu z panelu akcji inwentarza.',
    }),
    wiersz('Granica wykazu', granica, { klasa: 'dg-wiersz' }),
  );

  const odswiezPrzycisk = przyciskAkcji('Odśwież rekomendacje', 'dn-btn dn-btn--atrament');
  const raport = przyciskAkcji('Eksportuj raport');
  rama.akcje.append(
    odswiezPrzycisk,
    raport,
    przyciskBezKomendy(
      'Zastosuj poprawkę',
      powodBezKomendy('Wdrożenie poprawki wymagałoby komendy zapisu; recommendation.list jest wyłącznie odczytem.'),
    ),
    przyciskBezKomendy(
      'Odrzuć',
      powodBezKomendy('Przestawienie stanu rekomendacji na odrzuconą wymagałoby komendy zapisu.'),
    ),
    przyciskBezKomendy(
      'Inna propozycja',
      powodBezKomendy('Poproszenie o kolejną propozycję poprawki wymagałoby komendy żądania.'),
    ),
  );

  rama.cialo.append(stanTresci);
  return { status, priorytet, granica, odswiezPrzycisk, raport };
}

/** Podpina pasek akcji do czynności okna panelu rekomendacji modułu Diagnostics w oknie całej tej sesji. */
function podepnijAkcjeRekomendacji(
  powierzchnia: PowierzchniaRekomendacji,
  obsluga: { odczytaj: () => void; raportuj: () => void },
): void {
  powierzchnia.odswiezPrzycisk.addEventListener('click', obsluga.odczytaj);
  powierzchnia.raport.addEventListener('click', obsluga.raportuj);
  // Zmiana filtra odczytuje od razu, tak samo jak w Errors Panel, inaczej wykaz byłby nieaktualny.
  powierzchnia.status.naZmiane(obsluga.odczytaj);
  powierzchnia.priorytet.naZmiane(obsluga.odczytaj);
}

/** Zadanie odczytu rekomendacji z `analysisId` bieżącej analizy i filtrów operatora ustawionych w oknie. */
function zadanieRekomendacji(
  idAnalizy: string,
  powierzchnia: PowierzchniaRekomendacji,
): DiagnosticsRecommendationListRequest {
  const zadanie: DiagnosticsRecommendationListRequest = { analysisId: idAnalizy };
  if (powierzchnia.status.wartosc() !== '') {
    zadanie.status = powierzchnia.status.wartosc() as RecommendationStatus;
  }
  if (powierzchnia.priorytet.wartosc() !== '') {
    zadanie.priority = powierzchnia.priorytet.wartosc() as DiagnosticPriority;
  }
  const granica = Number.parseInt(powierzchnia.granica.value, 10);
  if (Number.isInteger(granica) && granica > 0) zadanie.limit = granica;
  return zadanie;
}

/**
 * Rysuje wykaz rekomendacji albo stan pustki — analiza bez rekomendacji jest
 * poprawna. Id analizy idzie w obu gałęziach, żeby było widać, czyj to wykaz,
 * nawet gdy jest pusty.
 */
function rysujRekomendacje(
  rekomendacje: readonly DiagnosticRecommendation[],
  idAnalizy: string,
  tresc: StanTresci,
): void {
  if (rekomendacje.length === 0) {
    tresc.pusto(`Analiza ${idAnalizy} nie ma jeszcze żadnej rekomendacji.`);
    return;
  }
  const naglowek = document.createElement('p');
  naglowek.className = 'dg-analiza__wynik';
  naglowek.textContent = `Rekomendacje analizy: ${idAnalizy}`;

  const lista = wykaz('Rekomendacje diagnostyczne', 'dg-wykaz');
  for (const rekomendacja of rekomendacje) lista.append(rysujRekomendacje1(rekomendacja));
  tresc.tresc().append(naglowek, lista);
}

/** Pojedyncza pozycja wykazu: opis, priorytet, znacznik skuteczności i poprawka proponowana operatorowi. */
function rysujRekomendacje1(rekomendacja: DiagnosticRecommendation): HTMLElement {
  const pozycja = pozycjaWykazu(
    rekomendacja.title,
    rekomendacja.detail === undefined || rekomendacja.detail === ''
      ? 'Rekomendacja bez uzasadnienia.'
      : rekomendacja.detail,
    'dg',
  );
  pozycja.element.dataset['priorytet'] = rekomendacja.priority;

  const znacznik = document.createElement('span');
  znacznik.className = 'dg-analiza__znacznik';
  znacznik.textContent = `Stan: ${rekomendacja.status}${
    rekomendacja.targetPath === undefined ? '' : ` · Plik: ${rekomendacja.targetPath}`
  }`;
  pozycja.element.append(znacznik);

  const czas = document.createElement('span');
  czas.className = 'dg-pozycja__czas';
  czas.textContent = new Date(rekomendacja.createdAt).toLocaleString('pl-PL');
  pozycja.element.append(czas);

  if (rekomendacja.patch !== undefined && rekomendacja.patch !== '') {
    const poprawka = document.createElement('pre');
    poprawka.className = 'dg-poprawka';
    poprawka.textContent = rekomendacja.patch;
    pozycja.element.append(poprawka);
  }

  return pozycja.element;
}

/** Raport rekomendacji w Markdown — treść pliku eksportu pobieranego przez operatora z tego panelu okna. */
function raportRekomendacji(rekomendacje: readonly DiagnosticRecommendation[]): string {
  const wiersze = ['# Raport rekomendacji diagnostycznych', ''];
  for (const rekomendacja of rekomendacje) {
    wiersze.push(
      `## ${rekomendacja.title}`,
      `- Analiza: ${rekomendacja.analysisId}`,
      `- Priorytet: ${rekomendacja.priority}`,
      `- Stan: ${rekomendacja.status}`,
      rekomendacja.targetPath === undefined ? '' : `- Plik: ${rekomendacja.targetPath}`,
      rekomendacja.detail === undefined ? '' : `- Uzasadnienie: ${rekomendacja.detail}`,
      `- Utworzono: ${new Date(rekomendacja.createdAt).toLocaleString('pl-PL')}`,
      rekomendacja.patch === undefined ? '' : `\n\`\`\`\n${rekomendacja.patch}\n\`\`\``,
      '',
    );
  }
  return wiersze.filter((linia) => linia !== '').join('\n');
}
