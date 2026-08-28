import './archiwum.css';

import type { Agent } from '../../../../shared/contract';
import { przycisk } from '../../modele/kontrolki-formularza';
import { przyciskBezKomendy } from '../../modele/kontrolki-formularza-braki';
import type { Kanal } from '../../protokol/kanal';
import {
  naglowekWykazu,
  podepnijDoWykazu,
  stanKomendy,
  zapewnijOdczyt,
  zdanieOPozycji,
} from '../wykaz-komend-rdzenia';
import { imie, odmowaRdzenia, wierszWykazu } from './archiwum-kontrolki';
import { utworzZrodloWersjiEksperta } from './zrodlo-wersji-eksperta';

/**
 * Panel „Archiwum ekspertów” Agent Buildera obsługuje archiwizację, wykaz
 * zarchiwizowanych oraz przywrócenie eksperta.
 */
/** Czynność panelu opisana wraz z nazwą komendy kontraktu, która ma tę czynność wykonać po stronie rdzenia. */
export interface CzynnoscArchiwum {
  /** Nazwa komendy w konwencji kontraktu — po angielsku, `obszar.zasob.akcja`. */
  komenda: string;
  /** Napis na kontrolce — po polsku, jak każdy napis widziany przez Operatora. */
  etykieta: string;
  /** Co ta czynność ma zrobić. */
  przeznaczenie: string;
  /** Czy czynność wymaga wybranego eksperta. */
  wymagaEksperta: boolean;
}

/** Czynności wywoływane wprost z panelu, niezależnie od wybranej pozycji w wykazie całego archiwum ekspertów. */
export const CZYNNOSCI_ARCHIWUM: readonly CzynnoscArchiwum[] = [
  {
    komenda: 'agent.archive',
    etykieta: 'Archiwizuj eksperta',
    przeznaczenie: 'zejście eksperta z biblioteki bez utraty jego definicji i historii',
    wymagaEksperta: true,
  },
  {
    komenda: 'agent.archive.list',
    etykieta: 'Pokaż zarchiwizowanych',
    przeznaczenie: 'wykaz ekspertów zarchiwizowanych, osobny od biblioteki czynnej',
    wymagaEksperta: false,
  },
];

/**
 * Czynności pozycji — bez kontrolki zbiorczej.
 *
 * Wymienione, bo wykaz pokrycia ma mówić o całej rodzinie komend archiwum;
 * pominięcie przywrócenia kazałoby Operatorowi zgadywać, czy rdzeń je zna.
 */
export const CZYNNOSCI_POZYCJI: readonly CzynnoscArchiwum[] = [
  {
    komenda: 'agent.restore',
    etykieta: 'Przywróć z archiwum',
    przeznaczenie: 'powrót eksperta z archiwum — przy każdym wierszu wykazu archiwum',
    wymagaEksperta: false,
  },
];

const WSZYSTKIE = [...CZYNNOSCI_ARCHIWUM, ...CZYNNOSCI_POZYCJI];

export interface PanelArchiwum {
  element: HTMLElement;
  /** Nanosi eksperta czynnego; jego nazwa wchodzi do zdań kontrolek. */
  ustaw(ekspert: Agent | null): void;
  /** Pyta rdzeń o wykaz komend, żeby zdania mówiły o rdzeniu, nie o domyśle. */
  wczytajPokrycie(): Promise<void>;
  /** Odpina nasłuch wykazu komend kanału. */
  rozlacz(): void;
}

/** @param przeladuj wołany po zmianie stanu biblioteki czynnej: archiwizacji eksperta albo jego przywróceniu z archiwum. */
export function utworzPanelArchiwum(kanal: Kanal, przeladuj?: () => void): PanelArchiwum {
  const zrodlo = utworzZrodloWersjiEksperta(kanal);
  const nazwyKomend = WSZYSTKIE.map((czynnosc) => czynnosc.komenda);

  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Archiwum ekspertów';

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis';

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'dn-pole-opis da-granica';
  wyjasnienie.textContent =
    'Archiwizacja NIE jest wyłączeniem eksperta. Wyłączony ekspert zostaje ' +
    'w bibliotece i da się go edytować; zarchiwizowany schodzi z niej i wraca ' +
    'na żądanie. Moduł nie podstawia jednego pod drugie.';

  const lista = document.createElement('ul');
  lista.className = 'da-archiwum__lista';

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'da-odpowiedz';
  odpowiedz.hidden = true;

  const wykaz = document.createElement('ul');
  wykaz.className = 'da-archiwum__wykaz';
  wykaz.hidden = true;

  const element = document.createElement('section');
  element.className = 'da-panel da-archiwum';
  element.append(tytul, wyjasnienie, podsumowanie, lista, odpowiedz, wykaz);

  let ekspertCzynny: Agent | null = null;

  /** Odpowiedź rdzenia — powodzenie albo nazwana odmowa, nigdy cisza. */
  function powiedz(zdanie: string, odmowa: boolean): void {
    odpowiedz.textContent = zdanie;
    odpowiedz.hidden = false;
    odpowiedz.dataset['waga'] = odmowa ? 'odmowa' : 'zgoda';
  }

  function przerysuj(): void {
    podsumowanie.textContent = naglowekWykazu(kanal, nazwyKomend);
    lista.replaceChildren(
      ...WSZYSTKIE.map((czynnosc) => wierszCzynnosci(kanal, czynnosc, ekspertCzynny, wykonaj)),
    );
  }

  /** Wywołuje czynność panelu; każda odmowa wraca zdaniem, nie ciszą. */
  async function wykonaj(czynnosc: CzynnoscArchiwum): Promise<void> {
    if (czynnosc.wymagaEksperta && ekspertCzynny === null) {
      powiedz(`${czynnosc.etykieta}: najpierw wybierz eksperta w bibliotece.`, true);
      return;
    }
    const kod = ekspertCzynny?.id ?? '';
    switch (czynnosc.komenda) {
      case 'agent.archive': {
        const wynik = await zrodlo.zarchiwizuj(kod);
        if (!wynik.udany || wynik.wynik === undefined) {
          powiedz(odmowaRdzenia(czynnosc.etykieta, wynik.blad?.code, wynik.blad?.message), true);
          return;
        }
        powiedz(
          wynik.wynik.archived
            ? `Ekspert „${imie(ekspertCzynny)}" zszedł do archiwum. Definicja i historia zostają.`
            : `Ekspert „${imie(ekspertCzynny)}" nie został odłożony — rdzeń zwrócił archived=false.`,
          !wynik.wynik.archived,
        );
        przeladuj?.();
        return;
      }
      case 'agent.archive.list': {
        const wynik = await zrodlo.archiwum();
        if (!wynik.udany || wynik.wynik === undefined) {
          powiedz(odmowaRdzenia(czynnosc.etykieta, wynik.blad?.code, wynik.blad?.message), true);
          return;
        }
        pokazArchiwum(wynik.wynik.agents, wynik.wynik.total);
        return;
      }
      default:
        // Czynności pozycji nie mają kontrolki zbiorczej i tu nie trafiają.
        return;
    }
  }

  /** Wykaz archiwum; każdy wiersz niesie własne przywrócenie. */
  function pokazArchiwum(eksperci: readonly Agent[], ile: number): void {
    wykaz.hidden = false;
    if (eksperci.length === 0) {
      powiedz('Archiwum jest puste — żaden ekspert nie został odłożony.', false);
      wykaz.replaceChildren();
      return;
    }
    powiedz(`Archiwum: ${ile} ${ile === 1 ? 'ekspert' : 'ekspertów'}.`, false);
    wykaz.replaceChildren(
      ...eksperci.map((ekspert) =>
        wierszWykazu(
          `${imie(ekspert)} — ${ekspert.id}`,
          'Przywróć',
          stanKomendy(kanal, 'agent.restore') === 'rdzen-ma',
          'agent.restore',
          async () => {
            const wynik = await zrodlo.przywrocZArchiwum(ekspert.id);
            if (!wynik.udany || wynik.wynik === undefined) {
              powiedz(odmowaRdzenia('Przywrócenie', wynik.blad?.code, wynik.blad?.message), true);
              return;
            }
            powiedz(
              wynik.wynik.restored
                ? `Ekspert „${imie(ekspert)}" wrócił do biblioteki wraz z zapamiętanym stanem czynności.`
                : `Ekspert „${imie(ekspert)}" nie wrócił — rdzeń zwrócił restored=false.`,
              !wynik.wynik.restored,
            );
            przeladuj?.();
            const odswiezone = await zrodlo.archiwum();
            if (odswiezone.udany && odswiezone.wynik !== undefined) {
              pokazArchiwum(odswiezone.wynik.agents, odswiezone.wynik.total);
            }
          },
        ),
      ),
    );
  }

  const odepnij = podepnijDoWykazu(kanal, przerysuj);
  przerysuj();

  return {
    element,

    ustaw(ekspert) {
      ekspertCzynny = ekspert;
      przerysuj();
    },

    async wczytajPokrycie() {
      await zapewnijOdczyt(kanal);
      przerysuj();
    },

    rozlacz() {
      odepnij();
    },
  };
}



/** Jeden wiersz czynności: nazwa, przeznaczenie oraz kontrolka wywołująca daną czynność albo nazywająca jej brak w kontrakcie. */
function wierszCzynnosci(
  kanal: Kanal,
  czynnosc: CzynnoscArchiwum,
  ekspert: Agent | null,
  wykonaj: (czynnosc: CzynnoscArchiwum) => Promise<void>,
): HTMLElement {
  const opis = document.createElement('span');
  opis.className = 'da-archiwum__opis';
  opis.textContent = `${czynnosc.komenda} — ${czynnosc.przeznaczenie}`;

  const element = document.createElement('li');
  element.className = 'da-archiwum__wiersz';
  element.dataset['pokrycie'] = stanKomendy(kanal, czynnosc.komenda);
  element.append(opis);

  const pozycyjna = CZYNNOSCI_POZYCJI.some((inna) => inna.komenda === czynnosc.komenda);
  if (pozycyjna) {
    const wskazanie = document.createElement('span');
    wskazanie.className = 'dn-pole-opis';
    wskazanie.textContent =
      stanKomendy(kanal, czynnosc.komenda) === 'rdzen-ma'
        ? 'kontrolka przy każdym wierszu wykazu niżej'
        : 'rdzeń tej komendy nie wpina';
    element.append(wskazanie);
    return element;
  }

  if (stanKomendy(kanal, czynnosc.komenda) === 'rdzen-ma') {
    const kontrolka = przycisk(czynnosc.etykieta, 'dn-btn dn-btn--zarys');
    kontrolka.addEventListener('click', () => void wykonaj(czynnosc));
    element.append(kontrolka);
    return element;
  }

  const kogo = ekspert === null ? 'bez wybranego eksperta' : `ekspert „${imie(ekspert)}"`;
  const powod = `${zdanieOPozycji(kanal, [czynnosc.komenda], czynnosc.etykieta)} Czynność: ${czynnosc.przeznaczenie} (${kogo}).`;
  element.append(przyciskBezKomendy(czynnosc.etykieta, powod));
  return element;
}
