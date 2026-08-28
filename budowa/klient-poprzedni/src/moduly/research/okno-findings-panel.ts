import { ResearchFindingStatus, type ResearchFinding } from '../../../../shared/contract';
import { poleWielowierszowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_USTALEN } from './akcje-okien';
import {
  ustawStanUstalen,
  wykonajAkcjeUstalen,
  zapiszUstalenie,
  type KontekstUstalen,
} from './czynnosci-ustalen';
import { zDymkiem } from './dymek-badania';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';
import { utworzWierszUstalenia } from './wiersz-ustalenia';
import { utworzWyborNastawy, wierszNastawy } from './wybor-nastawy';
import { utworzWyborZrodel } from './wybor-zrodel';

/**
 * Findings Panel gromadzi ustalenia narastające w toku badania, choć narastanie na żywo jeszcze nie działa, bo rdzeń nie rozgłasza zdarzenia zmiany ustalenia.
 */
export interface OknoFindingsPanel {
  element: HTMLElement;
  odswiez(): void;
}

/** Stała wylicza dwie reprezentacje wykazu ustaleń, chronologiczną oraz według liczby źródeł, obie liczone po stronie klienta. */
const WIDOKI = [
  {
    wartosc: 'chronologia',
    etykieta: 'chronologia',
    opis: 'Ustalenia w kolejności zapisu — od najnowszego. Porządek bierze się z pola createdAt.',
  },
  {
    wartosc: 'zrodla',
    etykieta: 'po liczbie źródeł',
    opis: 'Najpierw ustalenia oparte na największej liczbie źródeł; bez powiązanego źródła — na końcu.',
  },
];

export function utworzOknoFindingsPanel(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoFindingsPanel {
  const poleTresci = poleWielowierszowe(
    { etykieta: 'Treść ustalenia', podpowiedz: 'co wynika z materiału' },
    3,
  );
  const znacznik = document.createElement('p');
  znacznik.className = 'dn-pole-opis mr-znacznik';
  let edytowane = '';

  const kontekst: KontekstUstalen = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    tresc: poleTresci.kontrolka,
    edytowane: () => edytowane,
    wciagnij: (identyfikator, tresc) => wciagnij(identyfikator, tresc),
    przejdz,
  };

  const zrodla = utworzWyborZrodel(stan.wybraneZrodla, () => odswiez());
  const wykaz = document.createElement('ul');
  wykaz.className = 'mr-wykaz';

  // Widok kodowania jakościowego tu nie stoi: rdzeń nie ma jeszcze uchwytu książki kodów.
  const widok = utworzWyborNastawy('Widok ustaleń', WIDOKI, () => odswiez());

  const zapisz = przycisk('Zapisz ustalenie', 'dn-btn dn-btn--sm dn-btn--atrament');
  zapisz.addEventListener('click', () => void zapiszUstalenie(kontekst, ResearchFindingStatus.Open));

  kontekst.okno.tresc.append(
    zDymkiem(
      poleTresci.element,
      'Pole content komendy research.finding.add. Wypełnione pole findingId zamienia zapis w edycję.',
    ),
    zrodla.element,
    znacznik,
    zapisz,
    kontekst.odpowiedz.element,
    zDymkiem(
      wierszNastawy('Widok ustaleń', widok),
      'Porządek liczony po stronie klienta z pól, które rdzeń oddał. Widoku kodowania jakościowego tu nie ma: kody są w kontrakcie, ale ustalenie oddawane przez rdzeń ich nie niesie.',
    ),
    wykaz,
  );

  const rama = utworzRameBadania(
    KODY_OKIEN.ustalenia,
    'Findings Panel',
    'pomocnicze',
    AKCJE_USTALEN,
    (akcja) => void wykonajAkcjeUstalen(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  /** Wciąga ustalenie do formularza albo czyści go pod zapis nowy. */
  function wciagnij(identyfikator: string, tresc: string): void {
    edytowane = identyfikator;
    poleTresci.kontrolka.value = tresc;
    znacznik.textContent =
      identyfikator === ''
        ? 'Zapis założy ustalenie nowe (żądanie bez pola findingId).'
        : `Zapis zmieni ustalenie ${identyfikator} (pole findingId).`;
    odswiez();
  }

  function odswiez(): void {
    const ustalenia = stan.ustalenia();
    stan.wybraneUstalenia.ogranicz(ustalenia.map((wpis) => wpis.id));
    zrodla.odswiez(stan.zrodla());
    wykaz.replaceChildren(
      ...uporzadkuj(ustalenia, widok.wartosc()).map((ustalenie) =>
        utworzWierszUstalenia(ustalenie, stan.wybraneUstalenia.czyWybrana(ustalenie.id), {
          naZaznaczenie: (identyfikator) => {
            stan.wybraneUstalenia.przelacz(identyfikator);
            odswiez();
          },
          naEdycje: (wybrane) => wciagnij(wybrane.id, wybrane.content),
        }),
      ),
    );
    ustawStanUstalen(kontekst, ustalenia.length);
  }

  wciagnij('', '');
  return { element: rama.element, odswiez };
}

/**
 * Funkcja porządkuje wykaz ustaleń według wybranej reprezentacji, sortując kopię tablicy, by nie naruszyć porządku narastania w pamięci modułu.
 */
function uporzadkuj(
  ustalenia: readonly ResearchFinding[],
  widok: string,
): readonly ResearchFinding[] {
  if (widok === 'zrodla') {
    return [...ustalenia].sort(
      (jedno, drugie) => (drugie.sourceIds ?? []).length - (jedno.sourceIds ?? []).length,
    );
  }
  return [...ustalenia].sort((jedno, drugie) => drugie.createdAt - jedno.createdAt);
}
