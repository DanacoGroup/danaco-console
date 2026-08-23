import { AssistantActivityKind, type AssistantActivityEntry } from '../../../../shared/contract';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pobierzPlik,
  przycisk,
  utworzWierszOdpowiedzi,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import { logDziennika, nazwaPlikuDziennika } from './eksport-dziennika';
import { dzien, ODCZYTY, PUSTE } from './etykiety-assistant';
import { utworzFiltrDziennika, type FiltrDziennika } from './filtr-dziennika';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import { wierszDziennika } from './wiersz-dziennika';

/** Kod okna w katalogu rdzenia (`okno_operacyjne.kod`). */
export const KOD_OKNA = 'activity-feed';

/**
 * Activity Feed — okno monitorujące modułu Assistant.
 *
 * Daje przegląd historii działań oraz powrót do wyniku wcześniejszego
 * zlecenia. Zapis jest chronologiczny i wspólny dla poleceń głosowych
 * i tekstowych.
 *
 * Wpisy stoją w grupach dziennych: nagłówek doby rozdziela zapis tak, jak
 * rozdziela go pamięć Operatora („kiedy prosiłem o rezerwację sali"). Grupę
 * składa okno, bo rdzeń oddaje wykaz płaski — `assistant.activity.list` nie ma
 * pola grupowania.
 *
 * Zawężenie — pole szukania i rodzaj wpisu — liczy się również w oknie
 * (`filtr-dziennika.ts`); rdzeń nie ma czym zawęzić tego wykazu.
 *
 * Pustki są rozróżnione tak jak w `dostepy/stany-odczytu.ts`: „jeszcze nie
 * pytałem", „pytam" i „rdzeń nie zna ani jednego wpisu" to trzy różne zdania.
 * Czwartym jest pustka po zawężeniu — wpisy są, tylko żaden nie pasuje —
 * a piątym odmowa rdzenia wraz z jej powodem; komunikat zostaje w układzie,
 * a „Spróbuj ponownie" go nie usuwa.
 *
 * Powrót do wyniku nie potrzebuje osobnej komendy: wpis rodzaju `result` niesie
 * treść wyniku, a zawężenie dziennika do jednego zlecenia idzie polem
 * `actionId` komendy `assistant.activity.list`.
 */
export interface OknoActivityFeed {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoActivityFeed(stan: StanAssistant): OknoActivityFeed {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();
  const filtr: FiltrDziennika = utworzFiltrDziennika(() => odswiez());

  const lista = document.createElement('ul');
  lista.className = 'ma-dziennik';

  const zakres = document.createElement('p');
  zakres.className = 'dn-pole-opis ma-dziennik__zakres';

  const wynik = document.createElement('pre');
  wynik.className = 'ma-dziennik__wynik';
  wynik.hidden = true;

  okno.tresc.append(zakres, lista, odpowiedz.element, wynik);

  /** Wpisy widoczne po zawężeniu — one, i tylko one, wchodzą do eksportu. */
  let widoczne: readonly AssistantActivityEntry[] = [];

  const element = document.createElement('section');
  element.className = 'ma-okno ma-okno--monitor';
  element.dataset['okno'] = KOD_OKNA;
  element.append(
    utworzNaglowekOkna({
      tytul: 'Activity Feed',
      rola: 'monitor · chronologiczny zapis działań asystenta',
      kontrolki: [filtr.element],
    }),
    okno.element,
    pasekZapisu(stan, odpowiedz, {
      widoczne: () => widoczne,
      zakres: () => zakres.textContent ?? '',
    }),
  );

  function pokazWynik(idZlecenia: string): void {
    const wpis = stan
      .wpisy()
      .find(
        (pozycja) =>
          pozycja.actionId === idZlecenia && pozycja.kind === AssistantActivityKind.Result,
      );
    wynik.hidden = false;
    wynik.textContent =
      wpis === undefined ? 'Dziennik nie ma wpisu rodzaju „wynik" dla tego zlecenia.' : wpis.content;
  }

  /**
   * Dwie czynności wykonywane wprost na wpisie dziennika.
   *
   * Odsłuch pobiera bajty i odtwarza je w karcie: dźwięk wraca tam, skąd
   * wyszedł, a rdzeń oddaje wyłącznie nagrania, które sam wystawił.
   * Wyróżnienie idzie do rdzenia, więc przeżywa odświeżenie wykazu — po zapisie
   * czytamy dziennik ponownie, żeby wykaz pokazywał stan zapisany, a nie
   * przewidywany.
   */
  const czynnosciWpisu = {
    odsluch(odnosnik: string): void {
      void (async () => {
        const wynikOdsluchu = await stan.mowa.pobierzNagranie(odnosnik);
        if (!wynikOdsluchu.udany || wynikOdsluchu.wynik === undefined) {
          zglosBrak(
            'Odsłuch nagrania',
            opisOdmowy(
              'Pobranie nagrania',
              wynikOdsluchu.blad?.code,
              wynikOdsluchu.blad?.message,
            ),
          );
          return;
        }
        const dzwiek = new Audio(
          `data:${wynikOdsluchu.wynik.contentType};base64,${wynikOdsluchu.wynik.audio}`,
        );
        void dzwiek.play();
      })();
    },
    wyroznienie(idWpisu: string, wazny: boolean): void {
      void (async () => {
        const zapis = await stan.zrodlo.oznaczWpis(idWpisu, wazny, '');
        if (!zapis.udany) {
          zglosBrak(
            'Wyróżnienie wpisu',
            opisOdmowy('Oznaczenie wpisu', zapis.blad?.code, zapis.blad?.message),
          );
          return;
        }
        await stan.odswiezDziennik();
      })();
    },
  };

  function odswiez(): void {
    const wybrane = stan.wybrane();
    widoczne = stan.wpisy().filter((wpis) => filtr.przepusc(wpis));
    zakres.textContent =
      wybrane === ''
        ? 'Zakres: całość dziennika modułu.'
        : `Zakres: przebieg zlecenia ${wybrane}. Naciśnij „Cały dziennik", aby go zdjąć.`;
    lista.replaceChildren(
      ...grupujDniami(widoczne, (wpis) =>
        wierszDziennika(
          wpis,
          wpis.actionId === wybrane && wybrane !== '',
          (idZlecenia) => {
            void stan.wybierz(idZlecenia);
            pokazWynik(idZlecenia);
          },
          czynnosciWpisu,
        ),
      ),
    );
    naniesFaze(okno, stan, widoczne.length, wybrane !== '', filtr);
  }

  odswiez();
  return { element, odswiez };
}

/**
 * Wpisy przeplecione nagłówkami dni.
 *
 * Nagłówek pada wyłącznie przy zmianie doby, więc wykaz w kolejności rdzenia
 * daje po jednym nagłówku na dzień. Sortowania tu nie ma: kolejność zapisu
 * należy do rdzenia (`assistant.activity.list` oddaje wykaz uporządkowany),
 * a przestawianie go w oknie zamieniłoby chronologię rdzenia na chronologię
 * klienta.
 */
function grupujDniami(
  wpisy: readonly AssistantActivityEntry[],
  rysuj: (wpis: AssistantActivityEntry) => HTMLElement,
): HTMLElement[] {
  const pozycje: HTMLElement[] = [];
  let ostatni = '';
  for (const wpis of wpisy) {
    const biezacy = dzien(wpis.createdAt);
    if (biezacy !== ostatni) {
      pozycje.push(naglowekDnia(biezacy));
      ostatni = biezacy;
    }
    pozycje.push(rysuj(wpis));
  }
  return pozycje;
}

/** Nagłówek grupy dnia — pozycja listy o roli separatora, nie wpis dziennika. */
function naglowekDnia(nazwa: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'ma-dziennik__dzien';
  element.setAttribute('role', 'separator');
  element.textContent = nazwa;
  return element;
}

/** Wpisy i zakres widoku, z których powstaje plik eksportu. */
interface WidokEksportu {
  widoczne(): readonly AssistantActivityEntry[];
  zakres(): string;
}

/** Stopka okna: zdjęcie zawężenia, ponowienie odczytu i eksport dziennika. */
function pasekZapisu(
  stan: StanAssistant,
  odpowiedz: WierszOdpowiedzi,
  widok: WidokEksportu,
): HTMLElement {
  const calosc = przycisk('Cały dziennik', 'dn-btn dn-btn--sm dn-btn--zarys');
  calosc.addEventListener('click', () => {
    // Wybór jest przełącznikiem: wskazanie zlecenia już wybranego zdejmuje
    // zawężenie i czyta dziennik w całości.
    void stan.wybierz(stan.wybrane());
  });

  const ponow = przycisk('Spróbuj ponownie', 'dn-btn dn-btn--sm dn-btn--zarys');
  ponow.addEventListener('click', () => {
    // Komunikat błędu zostaje: ponowienie wyzwala odczyt, a zdanie o odmowie
    // znika dopiero, gdy odczyt się powiedzie.
    void stan.odswiezDziennik().then(() => {
      if (stan.fazaDziennika() !== 'blad') return;
      odpowiedz.pokaz(opisOdmowy('Ponowny odczyt dziennika', '', stan.powodDziennika()), false);
    });
  });

  const eksport = przycisk('Eksportuj dziennik', 'dn-btn dn-btn--sm dn-btn--zarys');
  eksport.addEventListener('click', () => {
    const wpisy = widok.widoczne();
    if (wpisy.length === 0) {
      // Plik o zerowej treści wyglądałby jak eksport udany, a nie jest niczym:
      // pobranie bez wpisów nazywamy zamiast je udawać.
      odpowiedz.pokaz(
        'Widok nie ma ani jednego wpisu — nie ma czego wydać. Zdejmij zawężenie albo ' +
          'poczekaj na pierwsze działanie asystenta.',
        false,
      );
      return;
    }
    const pobrano = Date.now();
    pobierzPlik(
      nazwaPlikuDziennika(pobrano),
      logDziennika(wpisy, widok.zakres(), pobrano),
      'text/plain',
    );
    odpowiedz.pokaz(
      `Dziennik pobrany jako log tekstowy — ${String(wpisy.length)} wpisów widocznych w oknie.`,
      true,
    );
  });

  const semantyczne = przycisk('Szukaj po znaczeniu', 'dn-btn dn-btn--sm dn-btn--duch');
  semantyczne.addEventListener('click', () =>
    zglosBrak('Wyszukiwanie po znaczeniu w dzienniku', BRAKI.szukanieZnaczeniem),
  );

  const element = document.createElement('div');
  element.className = 'ma-dziennik__stopka';
  element.append(calosc, ponow, eksport, semantyczne);
  return element;
}

/** Pięć stanów widoku: przed pytaniem, odczyt, pustka, pustka po zawężeniu, odmowa. */
function naniesFaze(
  okno: StanOkna,
  stan: StanAssistant,
  ile: number,
  zawezony: boolean,
  filtr: FiltrDziennika,
): void {
  const faza = stan.fazaDziennika();
  if (faza === 'blad') {
    okno.blad(stan.powodDziennika());
    return;
  }
  if (faza === 'ladowanie') {
    okno.ladowanie(ODCZYTY.dziennik);
    return;
  }
  if (ile === 0) {
    // Zanim rdzeń odpowie, okno mówi „jeszcze nie pytałem" — pusty dziennik
    // po odmowie i pusty dziennik przed odczytem to dwa różne stany.
    if (!stan.pytanoODziennik()) {
      okno.puste(PUSTE.dziennikSpoczynek);
      return;
    }
    // Zapis pełny, z którego zawężenie nie przepuściło niczego, to trzecia
    // sytuacja: wpisy są, więc zdanie o pustym dzienniku byłoby nieprawdą.
    if (filtr.zawezony() && stan.wpisy().length > 0) {
      okno.puste(PUSTE.dziennikSzukanie);
      return;
    }
    okno.puste(zawezony ? PUSTE.dziennikZlecenia : PUSTE.dziennik);
    return;
  }
  okno.gotowe();
}
