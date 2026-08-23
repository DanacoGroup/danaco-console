import { ProgressStatus } from '../../../shared/contract';
import { ETYKIETA_STANU_PROCESU } from './etykiety-pulpitu';
import type { ProcesWTle } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';
import { utworzPasekPostepu } from './pasek-postepu';
import { utworzStanPusty } from './stan-pusty';

/**
 * Kolumna procesów w tle: wykaz procesów telemetrii `progress.changed`
 * z paskiem postępu i napisem etapu.
 *
 * Kontrakt nie ma odczytu bieżących procesów — wykaz buduje się wyłącznie ze
 * zdarzeń, więc do pierwszego zdarzenia kolumna pokazuje stan pusty.
 *
 * Stan procesu jest wypisany słowem obok paska; barwa paska wspiera odczyt,
 * ale go nie zastępuje.
 */
export interface KolumnaProcesow {
  element: HTMLElement;
  odswiez(procesy: ProcesWTle[]): void;
}

/** Buduje kolumnę procesów w tle. */
export function utworzKolumneProcesow(procesy: ProcesWTle[]): KolumnaProcesow {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--procesy',
    {
      tytul: 'Procesy w tle',
      dopisek: 'Biegną niezależnie od tego, na co patrzysz.',
      ikona: 'zegar',
    },
    'mc-tytul-procesy',
  );

  const wykaz = document.createElement('ul');
  wykaz.className = 'mc-procesy';
  cialo.append(wykaz);

  const odswiez = (dane: ProcesWTle[]): void => {
    if (dane.length === 0) {
      wykaz.replaceChildren(
        utworzStanPusty(
          'Telemetria milczy',
          'Rdzeń nie zgłosił dotąd zdarzenia progress.changed; kontrakt nie ma odczytu bieżących procesów.',
        ),
      );
      return;
    }
    wykaz.replaceChildren(...dane.map(wiersz));
  };
  odswiez(procesy);

  return { element, odswiez };
}

/** Jeden proces: nazwa, pasek postępu, etap i stan. */
function wiersz(proces: ProcesWTle): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'mc-proces';
  element.dataset.stan = proces.status;

  const nazwa = document.createElement('span');
  nazwa.className = 'mc-proces__nazwa';
  nazwa.textContent = proces.nazwa;

  const stan = document.createElement('span');
  stan.className = `dn-plakietka ${klasaStanu(proces.status)} mc-proces__stan`;
  stan.textContent = ETYKIETA_STANU_PROCESU[proces.status];

  const gora = document.createElement('div');
  gora.className = 'mc-proces__gora';
  gora.append(nazwa, stan);

  const etap = document.createElement('span');
  etap.className = 'mc-proces__etap';
  // Telemetria mogła nie podać nazwy etapu — wtedy mówimy to wprost.
  etap.textContent = proces.etykietaEtapu ?? 'etap bez nazwy w telemetrii';

  element.append(
    gora,
    utworzPasekPostepu({
      etapBiezacy: proces.etapBiezacy,
      etapowRazem: proces.etapowRazem,
      status: proces.status,
      nazwa: proces.nazwa,
    }),
    etap,
  );
  return element;
}

/** Klasa plakietki stanu — słowo niesie treść, barwa jedynie ją wspiera. */
function klasaStanu(status: ProgressStatus): string {
  switch (status) {
    case ProgressStatus.Done:
      return 'dn-plakietka--sukces';
    case ProgressStatus.Failed:
    case ProgressStatus.Stopped:
      return 'dn-plakietka--blad';
    case ProgressStatus.Paused:
      return 'dn-plakietka--ostrzezenie';
    default:
      return 'dn-plakietka--informacja';
  }
}
