import { ExportFormat } from '../../../../shared/contract';
import {
  poleLogiczne,
  poleTekstowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { AKCJE_EKSPORTU } from './akcje-okien';
import {
  ustawStanEksportu,
  wydajRaport,
  wykonajAkcjeEksportu,
  type KontekstEksportu,
} from './czynnosci-eksportu';
import { zDymkiem } from './dymek-badania';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';
import { utworzWyborNastawy, wierszNastawy } from './wybor-nastawy';

/**
 * Export Panel — wydanie raportu w formacie dokumentowym.
 *
 * Dwie funkcje Operatora: eksport raportu oraz wybór formatu wyjściowego
 * i miejsca docelowego (`format`, `targetPath`, `toLibrary`).
 *
 * Plik składa widok; zachowanie po naciśnięciu leży w `czynnosci-eksportu`.
 */
export interface OknoExportPanel {
  element: HTMLElement;
  odswiez(): void;
}

const FORMATY = [
  { wartosc: ExportFormat.Pdf, etykieta: 'PDF' },
  { wartosc: ExportFormat.Docx, etykieta: 'DOCX' },
  { wartosc: ExportFormat.Markdown, etykieta: 'Markdown' },
  { wartosc: ExportFormat.Html, etykieta: 'HTML' },
  { wartosc: ExportFormat.Txt, etykieta: 'tekst zwykły' },
];

export function utworzOknoExportPanel(stan: StanBadania): OknoExportPanel {
  // Format wyjściowy idzie sterem nastawy: uchwyt niesie wartość bieżącą
  // („PDF"), a nie nazwę rodzajową — natywna lista pokazuje ją dopiero po
  // rozwinięciu.
  const format = utworzWyborNastawy('Format wyjściowy', FORMATY);
  const sciezka = poleTekstowe({ etykieta: 'Miejsce docelowe', podpowiedz: 'ścieżka pliku wynikowego' });
  const doRepozytorium = poleLogiczne({ etykieta: 'Zapisz wynik w repozytorium Library' });

  const ponow = przycisk('Spróbuj ponownie', 'dn-btn dn-btn--sm dn-btn--zarys');
  ponow.hidden = true;

  const podglad = document.createElement('p');
  podglad.className = 'mr-eksport__podglad';

  // Historia eksportów bieżącej sesji. Wpis powstaje wyłącznie z odpowiedzi
  // rdzenia — nie z zamówienia — więc wykaz mówi, co rdzeń NAPRAWDĘ oddał,
  // razem z brakiem ścieżki tam, gdzie jej nie oddał.
  const historia = document.createElement('ol');
  historia.className = 'mr-historia';
  historia.setAttribute('aria-label', 'Historia eksportów bieżącej sesji');

  const wpisyHistorii: string[] = [];

  const kontekst: KontekstEksportu = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    zlecenie: (idRaportu) => ({
      idRaportu,
      format: format.wartosc() === '' ? ExportFormat.Pdf : (format.wartosc() as ExportFormat),
      sciezka: sciezka.kontrolka.value,
      doRepozytorium: doRepozytorium.kontrolka.checked,
    }),
    odslonPonowienie: () => {
      ponow.hidden = false;
    },
    odnotujWydanie: (zdanie) => {
      wpisyHistorii.unshift(`${new Date().toLocaleString('pl-PL')} — ${zdanie}`);
      przerysujHistorie();
    },
    odswiez: () => odswiez(),
  };

  function przerysujHistorie(): void {
    if (wpisyHistorii.length === 0) {
      const pusto = document.createElement('li');
      pusto.className = 'mr-historia__pusto';
      pusto.textContent =
        'W tej sesji nic jeszcze nie wydano. Historia narasta z odpowiedzi rdzenia i kończy się ' +
        'wraz z połączeniem: komenda odczytu śladów wydań jest w kontrakcie, ale rdzeń nie ma ' +
        'jeszcze dla niej uchwytu.';
      historia.replaceChildren(pusto);
      return;
    }
    historia.replaceChildren(
      ...wpisyHistorii.map((wpis) => {
        const element = document.createElement('li');
        element.className = 'mr-historia__wpis';
        element.textContent = wpis;
        return element;
      }),
    );
  }

  const eksportuj = przycisk('Eksportuj teraz', 'dn-btn dn-btn--sm dn-btn--atrament');
  eksportuj.addEventListener('click', () => void wydajRaport(kontekst, 'Eksport raportu'));
  ponow.addEventListener('click', () => void wydajRaport(kontekst, 'Ponowny eksport raportu'));

  kontekst.okno.tresc.append(
    zDymkiem(
      wierszNastawy('Format wyjściowy', format),
      'Pole format: pdf, docx, markdown, html albo txt — wyliczenie ExportFormat kontraktu.',
    ),
    zDymkiem(sciezka.element, 'Pole targetPath — puste zostawia rdzeniowi wybór miejsca zapisu.'),
    zDymkiem(doRepozytorium.element, 'Pole toLibrary — wynik trafia do repozytorium Library jako plik.'),
    eksportuj,
    ponow,
    kontekst.odpowiedz.element,
    podglad,
    naglowekHistorii(),
    historia,
  );
  przerysujHistorie();

  const rama = utworzRameBadania(
    KODY_OKIEN.eksport,
    'Export Panel',
    'pomocnicze',
    AKCJE_EKSPORTU,
    (akcja) => void wykonajAkcjeEksportu(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  function odswiez(): void {
    podglad.textContent = opisRaportu(stan);
    ustawStanEksportu(kontekst);
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Podpis nad historią wydań bieżącej sesji. */
function naglowekHistorii(): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-etykieta';
  element.textContent = 'Historia eksportów bieżącej sesji';
  return element;
}

/** Zdanie podglądu: co dokładnie zostanie wydane. */
function opisRaportu(stan: StanBadania): string {
  const raport = stan.raport();
  if (raport === null) return '';
  const chwila = new Date(raport.updatedAt).toLocaleString('pl-PL');
  return `Raport „${raport.title}" — ${(raport.sections ?? []).length} sekcji, ostatnia zmiana ${chwila}.`;
}
