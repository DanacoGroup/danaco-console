import {
  SessionToolSource,
  SlashEntryKind,
  type SessionTool,
  type ToolCatalogEntry,
} from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  pole,
  pozycjaWykazu,
  przyciskAkcji,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  wiersz,
  wybor,
  wykaz,
  type WierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BRAKI, zglosBrak } from './braki-kontraktu';
import {
  NAZWY_RODZAJOW_KATALOGU,
  ODCZYTY,
  POZYCJE_FILTRA_KATALOGU,
  PUSTE,
  type KodFiltraKatalogu,
} from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import type { ZrodloNarzedzi } from './zrodlo-narzedzi';

/**
 * Zakładka narzędzi i serwerów MCP w Command & Tools Hub czyta jeden katalog rdzenia i dokłada
 * wskazane pozycje do karty sesji, nie do definicji eksperta.
 */
export interface PanelNarzedzi {
  element: HTMLElement;
  /** Odczyt katalogu i dołożeń karty sesji. */
  wczytaj(): Promise<void>;
}

/** Górna granica wykazu; katalog liczy setki pozycji, okno pokazuje wyłącznie wycinek zawężony przez rdzeń. */
const GRANICA_KATALOGU = 60;

/** Pozycja listy grup znacząca brak zawężenia grupą, pokazująca pełny katalog niezależnie od rodzaju pozycji. */
const KAZDA_GRUPA = 'wszystkie';

export function utworzPanelNarzedzi(
  stan: StanAssistant,
  zrodlo: ZrodloNarzedzi,
): PanelNarzedzi {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const szukanie = pole('Szukaj w katalogu narzędzi', 'fragment nazwy albo opisu');
  szukanie.type = 'search';
  const rodzaj = wybor('Rodzaj pozycji katalogu', POZYCJE_FILTRA_KATALOGU);
  rodzaj.value = 'wszystkie';
  const grupa = wybor('Grupa po przeznaczeniu', [[KAZDA_GRUPA, 'Grupa: wszystkie']]);

  const lista = wykaz('Pozycje katalogu po ukośniku', 'ma-wykaz');
  okno.tresc.append(lista);
  okno.puste(PUSTE.narzedziaSpoczynek);

  const bilans = document.createElement('p');
  bilans.className = 'dn-pole-opis';

  const element = document.createElement('div');
  element.className = 'ma-obszar';
  element.dataset['obszar'] = 'narzedzia';
  element.append(
    wiersz('Szukaj', szukanie, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Dopasowanie idzie także środkiem nazwy, bo pozycje noszą przedrostek źródła. ' +
        'Zawężenie liczy rdzeń (tools.catalog.list), nie okno.',
    }),
    wiersz('Rodzaj', rodzaj, {
      klasa: 'ma-wiersz',
      objasnienie:
        'Narzędzia i umiejętności dokłada się do karty sesji; komendy akcji wykonują ' +
        'czynność aplikacji i zestawu narzędzi nie zmieniają.',
    }),
    wiersz('Grupa', grupa, {
      klasa: 'ma-wiersz',
      objasnienie: 'Grupy przychodzą z rdzenia razem z wykazem — nie są spisem w kliencie.',
    }),
    przyciski(),
    bilans,
    okno.element,
    odpowiedz.element,
  );

  function przyciski(): HTMLElement {
    const czytaj = przyciskAkcji('Odczytaj katalog', 'dn-btn dn-btn--sm dn-btn--atrament');
    czytaj.addEventListener('click', () => void wczytaj());

    const zakres = przyciskAkcji('Zakres uprawnień i limity', 'dn-btn dn-btn--sm dn-btn--duch');
    zakres.dataset['brak'] = 'zakres-narzedzia';
    zakres.addEventListener('click', () =>
      zglosBrak('Zakres uprawnień narzędzia', BRAKI.zakresNarzedzia),
    );

    const rzad = document.createElement('div');
    rzad.className = 'ma-formularz__przyciski';
    rzad.append(czytaj, zakres);
    return rzad;
  }

  /** Dołożenie albo jego zdjęcie; czynność rozstrzyga stan oddany przez rdzeń. */
  async function przestawDolozenie(pozycja: ToolCatalogEntry, dolozone: boolean): Promise<void> {
    const sesja = stan.idSesji();
    if (sesja === '') {
      odpowiedz.pokaz(BRAKI.brakSesji, false);
      return;
    }
    odpowiedz.pokaz(ODCZYTY.dolozenie, true);
    if (dolozone) {
      const zdjete = await zrodlo.zdejmij({ sessionId: sesja, toolName: pozycja.name });
      if (!zdjete.udany || zdjete.wynik === undefined) {
        odpowiedz.pokaz(
          opisOdmowy('Zdjęcie dołożenia', zdjete.blad?.code, zdjete.blad?.message),
          false,
        );
        return;
      }
      // Fałsz w polu `detached` nie jest błędem: znaczy, że dołożenia nie było.
      odpowiedz.pokaz(
        zdjete.wynik.detached
          ? `Rdzeń zdjął dołożenie ${pozycja.shortName}; karta sesji ma ich teraz ` +
              `${String(zdjete.wynik.tools.length)}.`
          : `Rdzeń nie miał czego zdejmować — ${pozycja.shortName} nie było dołożone ` +
              'do tej karty sesji.',
        true,
      );
      await wczytaj();
      return;
    }
    const dolozone2 = await zrodlo.dolozy({
      sessionId: sesja,
      toolName: pozycja.name,
      // Rękę nazywamy wprost: dołożenie idzie z okna Operatora, nie z komendy asystenta.
      source: SessionToolSource.Assistant,
    });
    if (!dolozone2.udany || dolozone2.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Dołożenie narzędzia', dolozone2.blad?.code, dolozone2.blad?.message),
        false,
      );
      return;
    }
    odpowiedz.pokaz(
      dolozone2.wynik.alreadyAttached
        ? `${pozycja.shortName} było już dołożone do tej karty sesji — powtórzenie nie ` +
            'mnoży wpisów.'
        : `Rdzeń dołożył ${pozycja.shortName}; karta sesji ma teraz ` +
            `${String(dolozone2.wynik.tools.length)} dołożeń.`,
      true,
    );
    await wczytaj();
  }

  async function wczytaj(): Promise<void> {
    okno.ladowanie(ODCZYTY.narzedzia);
    const sesja = stan.idSesji();
    const wybranyRodzaj = rodzaj.value as KodFiltraKatalogu;
    const wynik = await zrodlo.katalog({
      limit: GRANICA_KATALOGU,
      ...(sesja === '' ? {} : { sessionId: sesja }),
      ...(wybranyRodzaj === 'wszystkie' ? {} : { kind: wybranyRodzaj }),
      ...(grupa.value === KAZDA_GRUPA ? {} : { group: grupa.value }),
      ...(szukanie.value.trim() === '' ? {} : { query: szukanie.value.trim() }),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt katalogu narzędzi', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const katalog = wynik.wynik;
    // Grupy przychodzą z rdzenia po zawężeniu, bez ani jednej nazwy grupy zapisanej w kliencie.
    ustawPozycje(grupa, [
      { wartosc: KAZDA_GRUPA, etykieta: 'Grupa: wszystkie' },
      ...katalog.groups.map((nazwa) => ({ wartosc: nazwa, etykieta: `Grupa: ${nazwa}` })),
    ]);
    lista.replaceChildren(
      ...katalog.entries.map((pozycja) =>
        wierszKatalogu(pozycja, (dolozone) => void przestawDolozenie(pozycja, dolozone)),
      ),
    );
    bilans.textContent =
      `Rdzeń zna ${String(katalog.total)} pozycji spełniających zawężenie; ` +
      `okno pokazuje ${String(katalog.entries.length)} (granica wykazu: ` +
      `${String(GRANICA_KATALOGU)}).`;
    if (katalog.entries.length === 0) {
      okno.puste(PUSTE.narzedzia);
      return;
    }
    okno.gotowe();
    await naniesDolozenia(sesja);
  }

  /** Dołożenia karty sesji jako zdanie bilansu, osobne od katalogu, bo mówi, co jest już dołożone. */
  async function naniesDolozenia(sesja: string): Promise<void> {
    if (sesja === '') {
      odpowiedz.pokaz(BRAKI.brakSesji, false);
      return;
    }
    const wynik = await zrodlo.dolozenia(sesja);
    if (!wynik.udany || wynik.wynik === undefined) {
      odpowiedz.pokaz(
        opisOdmowy('Odczyt dołożeń karty sesji', wynik.blad?.code, wynik.blad?.message),
        false,
      );
      return;
    }
    odpowiedz.pokaz(opisDolozen(wynik.wynik.tools), true);
  }

  return { element, wczytaj };
}

/** Zdanie o dołożeniach karty sesji; brak dołożeń też jest zdaniem, licząc pozycje wprost z odpowiedzi rdzenia. */
function opisDolozen(dolozenia: readonly SessionTool[]): string {
  if (dolozenia.length === 0) {
    return (
      'Karta sesji nie ma ani jednego dołożenia. Zestaw narzędzi tury jest wtedy sam ' +
      'zestawem z definicji eksperta.'
    );
  }
  return (
    `Dołożenia karty sesji (${String(dolozenia.length)}): ` +
    dolozenia.map((narzedzie) => narzedzie.shortName).join(', ') +
    '. Znikną wraz z sesją — definicji eksperta nie zmieniają.'
  );
}

/** Jeden wiersz katalogu wraz z czynnością dołożenia albo jej brakiem, nazwaną zgodnie ze stanem rdzenia. */
function wierszKatalogu(
  pozycja: ToolCatalogEntry,
  naPrzestawienie: (dolozone: boolean) => void,
): HTMLElement {
  const dolozone = pozycja.attached === true;
  const czesci = [
    pozycja.description,
    `rodzaj: ${NAZWY_RODZAJOW_KATALOGU[pozycja.kind]}`,
    `grupa: ${pozycja.group}`,
    `źródło: ${pozycja.origin}`,
  ];
  if (pozycja.kind === SlashEntryKind.Action && pozycja.command !== undefined) {
    czesci.push(`komenda rdzenia: ${pozycja.command}`);
  }
  czesci.push(dolozone ? 'dołożone do karty sesji' : 'niedołożone');

  const { element, akcje } = pozycjaWykazu(
    pozycja.shortName,
    czesci.filter((czesc) => czesc !== '').join(' · '),
    'ma',
  );
  element.dataset['pozycja'] = pozycja.name;
  element.dataset['dolozone'] = dolozone ? 'tak' : 'nie';

  const przycisk = przyciskAkcji(
    dolozone ? 'Zdejmij z karty sesji' : 'Dołóż do karty sesji',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  przycisk.addEventListener('click', () => {
    if (!pozycja.attachable) {
      // Pozycja niedokładalna zostaje klikalna i mówi, dlaczego nic nie zrobi po naciśnięciu.
      zglosBrak(
        'Dołożenie pozycji',
        `Rdzeń oznaczył ${pozycja.shortName} jako pozycję, której nie da się dołożyć do ` +
          'sesji (pole attachable). Komendy akcji wykonuje się z pola polecenia, nie ' +
          'dokłada do zestawu narzędzi.',
      );
      return;
    }
    naPrzestawienie(dolozone);
  });
  akcje.append(przycisk);
  return element;
}
