import type { Queue } from '../../../shared/contract';
import {
  utworzCzteryStery,
  type OpisCzterechSterow,
  type TozsamoscDlaOgniska,
} from './cztery-stery';
import {
  opiszZrodlaKolejki,
  utworzKolejkeDecyzji,
  type KolejkaDecyzji,
  type WpisKolejkiDecyzji,
} from './kolejka-decyzji';
import { opisOdmowyAod } from './odmowy-aod';
import { dodajPole, utworzAkapit, utworzPodtytul, utworzWykazPol } from './pola-wykazu';
import { NAZWY_RODZAJOW, NAZWY_WAG } from './rodzaje-sugestii';
import {
  NAZWY_POWODOW,
  NAZWY_ZRODEL,
  obserwacjaZMonitora,
  opiszOdstep,
  WagaDecyzji,
  type ObserwacjaProcesu,
} from './rozpoznanie-decyzji';
import type { ZrodloDecyzji } from './zrodlo-decyzji';

/**
 * Sekcja decyzji czekających — pierwsza w trzonie okna nakładki, bo to jedyna
 * sekcja mówiąca o tym, co stoi i czeka na człowieka.
 *
 * Pusta kolejka jest stanem poprawnym i tak jest opisana, zamiast ostrzeżeniem.
 *
 * Każdy wpis mówi cztery rzeczy: co czeka (proces, stan, etap), od kiedy
 * (z policzonym odstępem), czego dotyczy (okno, sesja, kolejka) i czym to
 * wykryto — odczyt nadrabiający i zdarzenie na żywo mają różną świeżość.
 * Wpis rozpoznany regułą sporną (`usterka`, `bez-ruchu`) jest oznaczony jako
 * ocena nakładki: rdzeń pojęcia decyzji nie ma.
 *
 * Kolejka jest magazynem w pamięci okna, zasilanym odpowiedzią rdzenia
 * i zdarzeniami; zamknięcie okna nic nie utrwala.
 */
export interface SekcjaDecyzji {
  element: HTMLElement;
  /** Kolejka, do której okno dosypuje sygnały na żywo. */
  kolejka: KolejkaDecyzji;
  /** Odczyt nadrabiający: `monitor.status` + `queue.list`. */
  odswiez(): Promise<void>;
  /** Przerysowuje wykaz z tego, co w kolejce już jest — bez pytania rdzenia. */
  przerysuj(): void;
  /**
   * Opis sterów gotowy do użycia poza sekcją — dymek kontekstowy stawia te same
   * cztery stery przy tej samej decyzji i ma je stawiać z tego samego opisu,
   * żeby czynność z dymka i czynność z listy szły tą samą drogą.
   */
  stery: OpisCzterechSterow;
  /** Kolejka rdzenia dopasowana do decyzji; `undefined`, gdy dopasowania nie ma. */
  dopasujKolejke(idOkna?: string, idSesji?: string): Queue | undefined;
}

export interface OpisSekcjiDecyzji {
  zrodlo: ZrodloDecyzji;
  /** Krótkie potwierdzenie czynności na pasku okna. */
  zamelduj(zdanie: string, udane: boolean): void;
  /** Otwiera Okno Konfiguracji klienta. */
  otworzOknoKonfiguracji(): void;
  /** Tożsamość klienta z powitania; brak wyłącza wyłącznie przestawienie ogniska. */
  klient?: TozsamoscDlaOgniska;
  /** Zegar podawany z zewnątrz — sprawdziany nie zależą od czasu maszyny. */
  teraz?: () => number;
}

export function utworzSekcjeDecyzji(opis: OpisSekcjiDecyzji): SekcjaDecyzji {
  const teraz = opis.teraz ?? (() => Date.now());
  const kolejka = utworzKolejkeDecyzji();

  /** Jeden opis sterów dla listy i dla dymka — dwa wejścia, jedna droga. */
  const stery: OpisCzterechSterow = {
    zrodlo: opis.zrodlo,
    zamelduj: opis.zamelduj,
    poCzynnosci: () => void odswiez(),
    otworzOknoKonfiguracji: opis.otworzOknoKonfiguracji,
    ...(opis.klient === undefined ? {} : { klient: opis.klient }),
  };

  const element = document.createElement('section');
  element.className = 'ao-sekcja ao-decyzje';

  const licznik = document.createElement('p');
  licznik.className = 'ao-decyzje__licznik';

  const miejsce = document.createElement('div');
  miejsce.className = 'ao-sekcja__tresc';

  element.append(
    utworzPodtytul('Decyzje czekające na Operatora — wykryte przez nakładkę'),
    licznik,
    miejsce,
  );

  /** Kolejki z ostatniego odczytu — podstawa dopasowania decyzji do kolejki. */
  let kolejkiRdzenia: readonly Queue[] = [];

  /** Ostatnia odmowa odczytu; rysowana zamiast wykazu, dopóki trwa. */
  let odmowa: string | null = null;

  function przerysuj(): void {
    const wpisy = kolejka.wykaz(teraz());

    licznik.textContent =
      wpisy.length === 0
        ? 'Nakładka nie widzi żadnej decyzji czekającej.'
        : `Decyzji czekających: ${wpisy.length}. Skąd wiadomo: ${opiszZrodlaKolejki(wpisy)}.`;

    if (odmowa !== null) {
      miejsce.replaceChildren(utworzAkapit('ao-odmowa', odmowa));
      return;
    }

    if (wpisy.length === 0) {
      miejsce.replaceChildren(
        utworzAkapit(
          'ao-pusto',
          'Żaden proces nie stoi w oczekiwaniu na rozstrzygnięcie. Pusta kolejka jest stanem ' +
            'poprawnym — nakładka słucha progress.changed i window.state.changed na żywo ' +
            'i dopisze wpis w chwili, gdy proces stanie.',
        ),
      );
      return;
    }

    const lista = document.createElement('ul');
    lista.className = 'ao-decyzje__wykaz';
    for (const wpis of wpisy) lista.append(zbudujWpis(wpis));
    miejsce.replaceChildren(lista);
  }

  function zbudujWpis(wpis: WpisKolejkiDecyzji): HTMLLIElement {
    const { decyzja } = wpis;
    const wiersz = document.createElement('li');
    wiersz.className = 'ao-decyzja';
    wiersz.dataset['waga'] = decyzja.waga;

    const naglowek = document.createElement('p');
    naglowek.className = 'ao-decyzja__zdanie';
    naglowek.textContent = decyzja.zdanie;

    const znacznik = document.createElement('p');
    znacznik.className = 'ao-decyzja__waga';
    znacznik.textContent =
      decyzja.waga === WagaDecyzji.Pewna
        ? `${NAZWY_POWODOW[decyzja.powod]} — stan rdzenia mówi to wprost`
        : `${NAZWY_POWODOW[decyzja.powod]} — OCENA NAKŁADKI, nie orzeczenie rdzenia`;

    const kolejkaWpisu = dopasujKolejke(decyzja.idOkna, decyzja.idSesji);

    const pola = utworzWykazPol();
    // Rodzaj i waga ujawnienia idą pierwsze: katalog rozdz. 4 opracowania nazywa
    // sugestię, zanim Operator zacznie czytać jej pola techniczne.
    dodajPole(pola, 'Rodzaj sugestii', NAZWY_RODZAJOW[decyzja.rodzaj]);
    dodajPole(
      pola,
      'Waga ujawnienia',
      decyzja.krytyczna
        ? `${NAZWY_WAG[decyzja.wagaUjawnienia]} — punkt decyzyjny wstrzymujący proces`
        : NAZWY_WAG[decyzja.wagaUjawnienia],
    );
    dodajPole(pola, 'Od kiedy czeka', `${opiszOdstep(wpis.czekaMs)} (${new Date(decyzja.czekaOd).toLocaleString()})`);
    dodajPole(pola, 'Proces', decyzja.idProcesu ?? '');
    dodajPole(pola, 'Stan telemetrii', decyzja.stan);
    dodajPole(pola, 'Etap', opiszEtap(decyzja.etap, decyzja.numerEtapu, decyzja.liczbaEtapow));
    dodajPole(pola, 'Okno', decyzja.idOkna ?? '');
    dodajPole(pola, 'Karta sesji', decyzja.idSesji ?? '');
    dodajPole(pola, 'Kolejka', kolejkaWpisu === undefined ? '' : `${kolejkaWpisu.id} (${kolejkaWpisu.status})`);
    dodajPole(pola, 'Okno koordynatora', decyzja.idOknaKoordynatora ?? '');
    dodajPole(pola, 'Powód zatrzymania biegu', decyzja.powodZatrzymaniaBiegu ?? '');
    dodajPole(pola, 'Czym wykryto', NAZWY_ZRODEL[decyzja.zrodlo]);

    wiersz.append(
      naglowek,
      znacznik,
      pola,
      utworzCzteryStery(stery, { decyzja, kolejka: kolejkaWpisu }),
    );
    return wiersz;
  }

  /**
   * Dopasowuje kolejkę do decyzji.
   *
   * `MonitorStatus` niesie proces, okno i sesję, ale nie nazywa kolejki. Dojście
   * prowadzi więc przez `queue.list`: najpierw kolejka obsługująca to okno, potem
   * kolejka tej sesji. Przy niejednoznaczności zwracamy `undefined`, a ster
   * wstrzymania powie, że kolejki nie dopasowano.
   */
  function dopasujKolejke(idOkna?: string, idSesji?: string): Queue | undefined {
    if (idOkna !== undefined && idOkna !== '') {
      const poOknie = kolejkiRdzenia.filter((kolejkaRdzenia) =>
        (kolejkaRdzenia.windowIds ?? []).includes(idOkna),
      );
      if (poOknie.length === 1) return poOknie[0];
      if (poOknie.length > 1) return undefined;
    }
    if (idSesji !== undefined && idSesji !== '') {
      const poSesji = kolejkiRdzenia.filter((kolejkaRdzenia) => kolejkaRdzenia.sessionId === idSesji);
      if (poSesji.length === 1) return poSesji[0];
    }
    return undefined;
  }

  async function odswiez(): Promise<void> {
    const chwila = teraz();

    const wynikKolejek = await opis.zrodlo.kolejki();
    // Odmowa `queue.list` nie przerywa odczytu decyzji: bez kolejek wykaz nadal
    // mówi prawdę, ubywa jedynie jednej drogi wyjścia.
    kolejkiRdzenia = wynikKolejek.udany && wynikKolejek.wynik !== undefined ? wynikKolejek.wynik.queues : [];

    const wynikProcesow = await opis.zrodlo.procesy();
    if (!wynikProcesow.udany || wynikProcesow.wynik === undefined) {
      odmowa = opisOdmowyAod('Odczyt procesów (monitor.status)', 'procesy', wynikProcesow.blad);
      przerysuj();
      return;
    }

    odmowa = null;
    const obserwacje: ObserwacjaProcesu[] = wynikProcesow.wynik.statuses.map(obserwacjaZMonitora);
    kolejka.nanieOdczyt(obserwacje, chwila);
    przerysuj();
  }

  przerysuj();

  return { element, kolejka, odswiez, przerysuj, stery, dopasujKolejke };
}

/** Opisuje etap słowem: nazwa, numer i liczba etapów, każde tylko gdy jest. */
function opiszEtap(etap?: string, numer?: number, liczba?: number): string {
  const czesci: string[] = [];
  if (etap !== undefined && etap !== '') czesci.push(etap);
  if (numer !== undefined && numer > 0) {
    czesci.push(liczba !== undefined && liczba > 0 ? `${numer} z ${liczba}` : `nr ${numer}`);
  }
  return czesci.join(' · ');
}
