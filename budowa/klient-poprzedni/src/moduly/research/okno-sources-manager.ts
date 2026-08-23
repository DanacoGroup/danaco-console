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
 * Sources Manager — okno **zarządca** źródeł badania.
 *
 * Trzy czynności: dodanie źródła, katalogowanie (typ, pochodzenie, adres,
 * dokument repozytorium) i ocena wiarygodności. Wszystkie mieszczą się
 * w jednej komendzie `research.source.add` — kontrakt ma dla nich pola.
 *
 * Pole zawężania działa po stronie klienta i nie woła rdzenia. Komenda odczytu
 * wykazu wraz z polami zawężającymi jest już w kontrakcie, ale rdzeń nie ma dla
 * niej uchwytu, więc wykaz mieszka w pamięci modułu i tam też się zawęża. Po
 * dobudowie zawężanie ma przenieść się do żądania — wtedy zniknie i ta uwaga.
 *
 * Zawężenie nie zdejmuje zaznaczenia — źródło niewidoczne w wykazie pozostaje
 * zaznaczone i idzie do ustalenia, a okno mówi o tym liczbą przy polu.
 *
 * Plik składa widok; zachowanie po naciśnięciu leży w `czynnosci-zrodel`.
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
    // Stan okna liczy się z wykazu PEŁNEGO, nie z zawężonego: zawężenie bez
    // trafień nie znaczy, że katalog jest pusty, a zaproszenie „skataloguj
    // pierwsze źródło" postawione nad katalogiem pełnym byłoby nieprawdą.
    ustawStanZrodel(kontekst, zrodla.length);
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Zawężenie wykazu po frazie — tytuł, pochodzenie i adres źródła. */
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

/** Zdanie o tym, ile pozycji widać, ile jest i ile zaznaczono poza widokiem. */
function opisZawezenia(wszystkie: number, widoczne: number, zaznaczone: number): string {
  const podstawa =
    widoczne === wszystkie
      ? `Źródeł w katalogu: ${String(wszystkie)}.`
      : `Widocznych ${String(widoczne)} z ${String(wszystkie)} — reszta poza zawężeniem.`;
  return zaznaczone === 0 ? podstawa : `${podstawa} Zaznaczonych: ${String(zaznaczone)}.`;
}
