import type { ResearchSource } from '../../../../shared/contract';
import { poleTekstowe, przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_ZRODEL } from './akcje-okien';
import {
  skatalogujZrodlo,
  ustawStanZrodel,
  wykonajAkcjeZrodel,
  type KontekstZrodel,
} from './czynnosci-zrodel';
import { zDymkiem } from './dymek-badania';
import { utworzFormularzZrodla } from './formularz-zrodla';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';
import { utworzWierszZrodla } from './wiersz-zrodla';

/**
 * Sources Manager jest oknem zarządcą źródeł badania, łączącym dodanie źródła, katalogowanie i ocenę wiarygodności w jednej komendzie kontraktu.
 */
export interface OknoSourcesManager {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoSourcesManager(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoSourcesManager {
  const kontekst: KontekstZrodel = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    formularz: utworzFormularzZrodla(),
    przejdz,
  };

  const wykaz = document.createElement('ul');
  wykaz.className = 'mr-wykaz';

  const dodaj = przycisk('Skataloguj źródło', 'dn-btn dn-btn--sm dn-btn--atrament');
  dodaj.addEventListener('click', () => void skatalogujZrodlo(kontekst));

  const zawezenie = poleTekstowe({
    etykieta: 'Zawęź wykaz',
    podpowiedz: 'fraza w tytule, pochodzeniu albo adresie',
  });
  zawezenie.kontrolka.addEventListener('input', () => odswiez());

  const licznik = document.createElement('p');
  licznik.className = 'dn-pole-opis mr-licznik';

  kontekst.okno.tresc.append(
    kontekst.formularz.element,
    dodaj,
    kontekst.odpowiedz.element,
    zDymkiem(
      zawezenie.element,
      'Zawężenie działa po stronie klienta: komenda odczytu wykazu jest w kontrakcie, ale rdzeń nie ma jeszcze jej uchwytu, więc wykaz mieszka w pamięci modułu. Zaznaczenia nie zdejmuje.',
    ),
    licznik,
    wykaz,
  );

  const rama = utworzRameBadania(
    KODY_OKIEN.zrodla,
    'Sources Manager',
    'zarządca',
    AKCJE_ZRODEL,
    (akcja) => void wykonajAkcjeZrodel(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  function odswiez(): void {
    const zrodla = stan.zrodla();
    stan.wybraneZrodla.ogranicz(zrodla.map((zrodlo) => zrodlo.id));
    stan.lektura.ogranicz(zrodla.map((zrodlo) => zrodlo.id));

    const widoczne = zawez(zrodla, zawezenie.kontrolka.value);
    licznik.textContent = opisZawezenia(zrodla.length, widoczne.length, stan.wybraneZrodla.wybrane().length);
    wykaz.replaceChildren(
      ...widoczne.map((zrodlo) =>
        utworzWierszZrodla(zrodlo, stan.wybraneZrodla.czyWybrana(zrodlo.id), {
          naZaznaczenie: (identyfikator) => {
            stan.wybraneZrodla.przelacz(identyfikator);
            odswiez();
          },
          naLekture: (identyfikator) => {
            stan.lektura.wskaz(identyfikator);
            przejdz(KODY_OKIEN.lektura);
          },
        }),
      ),
    );
    // Stan okna liczy się z wykazu pełnego, nie zawężonego, by zaproszenie nie mijało się z prawdą.
    ustawStanZrodel(kontekst, zrodla.length);
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Funkcja zawęża wykaz źródeł po frazie szukanej w tytule, pochodzeniu oraz adresie źródła po stronie klienta. */
function zawez(zrodla: readonly ResearchSource[], fraza: string): readonly ResearchSource[] {
  const szukana = fraza.trim().toLocaleLowerCase('pl-PL');
  if (szukana === '') return zrodla;
  return zrodla.filter((zrodlo) =>
    [zrodlo.title, zrodlo.origin ?? '', zrodlo.url ?? '']
      .join(' ')
      .toLocaleLowerCase('pl-PL')
      .includes(szukana),
  );
}

/** Funkcja układa zdanie informujące, ile pozycji jest widocznych, ile jest w katalogu oraz ile zaznaczono poza widokiem. */
function opisZawezenia(wszystkie: number, widoczne: number, zaznaczone: number): string {
  const podstawa =
    widoczne === wszystkie
      ? `Źródeł w katalogu: ${String(wszystkie)}.`
      : `Widocznych ${String(widoczne)} z ${String(wszystkie)} — reszta poza zawężeniem.`;
  return zaznaczone === 0 ? podstawa : `${podstawa} Zaznaczonych: ${String(zaznaczone)}.`;
}
