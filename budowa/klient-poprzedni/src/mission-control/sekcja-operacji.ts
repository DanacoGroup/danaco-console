import { elementIkony } from '../ikony/ikony';
import { ETYKIETA_BRAKU_ZRODLA } from './etykiety-pulpitu';
import type { KanalOperacyjny } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import { utworzStanPusty } from './stan-pusty';

/**
 * Sekcja operacji AI.
 *
 * Jedna odpowiedzialność: plakietki kanałów modelu z odczytu `channel.list`.
 * Każda plakietka niesie to, co kontrakt oddaje: nazwę, rodzaj, model i stan
 * czynności kanału.
 *
 * Wysycenia, kolejki, limitu ani kosztu narastającego kontrakt nie niesie, więc
 * plakietka wypisuje przy nich etykietę braku źródła zamiast wymyślonych liczb.
 */
export interface SekcjaOperacji {
  element: HTMLElement;
  odswiez(kanaly: KanalOperacyjny[]): void;
}

/** Buduje sekcję operacji z plakietkami kanałów. */
export function utworzSekcjeOperacji(kanaly: KanalOperacyjny[]): SekcjaOperacji {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--operacje',
    {
      tytul: 'Operacje AI',
      dopisek: 'Rejestr kanałów modelu z odczytu channel.list.',
      ikona: 'globus',
    },
    'mc-tytul-operacje',
  );

  const siatka = document.createElement('div');
  siatka.className = 'mc-kanaly';
  cialo.append(siatka);

  const odswiez = (dane: KanalOperacyjny[]): void => {
    if (dane.length === 0) {
      siatka.replaceChildren(
        utworzStanPusty(
          'Rejestr kanałów pusty',
          'Odczyt channel.list nie zwrócił żadnego kanału modelu.',
        ),
      );
      return;
    }
    siatka.replaceChildren(...dane.map(plakietkaKanalu));
  };
  odswiez(kanaly);

  return { element, odswiez };
}

/** Plakietka jednego kanału: nazwa, stan, rodzaj, model i uwaga o miarach. */
function plakietkaKanalu(kanal: KanalOperacyjny): HTMLElement {
  const karta = document.createElement('article');
  karta.className = 'dn-karta mc-kanal';

  const naglowek = document.createElement('header');
  naglowek.className = 'mc-kanal__naglowek';

  const nazwa = document.createElement('span');
  nazwa.className = 'mc-kanal__nazwa';
  nazwa.textContent = kanal.nazwa;

  // Stan kanału niesie ikonę i słowo — nigdy sam kolor.
  const stan = document.createElement('span');
  stan.className = `dn-plakietka ${
    kanal.czynny ? 'dn-plakietka--sukces' : 'dn-plakietka--informacja'
  }`;
  stan.append(
    elementIkony(kanal.czynny ? 'ptaszek' : 'zegar', { rozmiar: 16 }),
    document.createTextNode(kanal.czynny ? 'czynny' : 'wyłączony'),
  );

  naglowek.append(nazwa, stan);

  const miary = document.createElement('dl');
  miary.className = 'mc-kanal__miary';
  miary.append(
    miara('rodzaj', kanal.rodzaj, 'Rodzaj kanału z rejestru — wartość danych'),
    miara('model', kanal.model ?? '—', 'Identyfikator modelu z wiersza rejestru'),
  );

  // Miary bez źródła — wypisane wprost, nie zmyślone.
  const uwaga = document.createElement('p');
  uwaga.className = 'mc-kanal__uwaga';
  uwaga.textContent = `wysycenie · kolejka · limit · koszt — ${ETYKIETA_BRAKU_ZRODLA}`;
  uwaga.title =
    'Kontrakt nie niesie telemetrii kanału: wysycenia gniazd, kolejki zleceń, limitu równoległości ani kosztu narastającego.';

  karta.append(naglowek, miary, uwaga);
  return karta;
}

/** Para „nazwa miary — wartość" z objaśnieniem w podpowiedzi. */
function miara(etykieta: string, wartosc: string, objasnienie: string): DocumentFragment {
  const para = document.createDocumentFragment();

  const nazwa = document.createElement('dt');
  nazwa.className = 'mc-kanal__etykieta';
  nazwa.textContent = etykieta;
  nazwa.title = objasnienie;

  const liczba = document.createElement('dd');
  liczba.className = 'mc-kanal__wartosc';
  liczba.textContent = wartosc;

  para.append(nazwa, liczba);
  return para;
}
