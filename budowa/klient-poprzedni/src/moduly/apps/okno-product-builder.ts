import { AppStageStatus } from '../../../../shared/contract';
import type { AppDeployment, AppStage, ProgressChangedEvent } from '../../../../shared/contract';
import { utworzWykazBrakow } from './braki-kontraktu';
import {
  BEZ_KOMENDY_ODCZYTU,
  BRAKI_PRODUCT_BUILDER,
  KODY_OKIEN,
  NAZWY_OKIEN,
} from './etykiety-apps';
import { utworzOsEtapow, type PozycjaOsi } from './os-etapow';
import { narzedziaProductBuilder } from './narzedzia-apps';
import { utworzPrzybornikApps } from './przybornik-apps';
import { utworzRameApps } from './rama-okna';
import type { StanProduktu } from './stan-produktu';
import type { RachunekRamek } from './zbior-budowy';

/**
 * Product Builder jest oknem wiodącym modułu Apps i punktem wejścia
 * integrującym pozostałe okna, zasilanym wyłącznie zdarzeniami — kontrakt
 * nie ma dla niego ani jednej komendy odczytu.
 */
export interface OknoProductBuilder {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoProductBuilder(
  stan: StanProduktu,
  pozycje: readonly PozycjaOsi[],
  przejdz: (kodOkna: string) => void,
): OknoProductBuilder {
  const kod = KODY_OKIEN.ProductBuilder;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'wiodące');
  const os = utworzOsEtapow(pozycje, przejdz);

  const kondycja = document.createElement('p');
  kondycja.className = 'mp-kondycja';

  const pasek = document.createElement('div');
  pasek.className = 'dn-postep';
  const wypelnienie = document.createElement('div');
  wypelnienie.className = 'dn-postep-wartosc';
  pasek.append(wypelnienie);

  const opisPostepu = document.createElement('p');
  opisPostepu.className = 'dn-postep-etykieta mp-postep__opis';

  const etapy = document.createElement('ul');
  etapy.className = 'mp-etapy';

  const pustkaEtapow = document.createElement('p');
  pustkaEtapow.className = 'mp-etapy__pustka';

  const wydania = document.createElement('ul');
  wydania.className = 'mp-wydania';

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_PRODUCT_BUILDER));
  // Przybornik stoi obok wykazu braków: ten prowadzi czynności, które rdzeń obsługuje.
  rama.akcje.append(
    utworzPrzybornikApps('Produkt, etapy i oś czasu', narzedziaProductBuilder(stan)).element,
  );
  rama.tresc.append(
    os.element,
    kondycja,
    opisPostepu,
    pasek,
    naglowek('Etapy z rdzenia (apps.build.changed)'),
    pustkaEtapow,
    etapy,
    naglowek('Dziennik wydań'),
    wydania,
  );

  function odswiez(): void {
    const wykazEtapow = stan.etapy();
    const wykazWydan = stan.wdrozenia();
    etapy.replaceChildren(...wykazEtapow.map(wierszEtapu));
    wydania.replaceChildren(...wykazWydan.map(wierszWydania));
    pustkaEtapow.hidden = wykazEtapow.length > 0;
    // Zdanie składamy przy każdym odświeżeniu, a nie raz przy budowie okna.
    if (wykazEtapow.length === 0) pustkaEtapow.textContent = zdaniePustkiEtapow(stan.ramki());
    kondycja.textContent = opisKondycji(wykazEtapow, wykazWydan);
    naniesPostep(stan.postep());
    if (wykazEtapow.length === 0 && wykazWydan.length === 0) {
      rama.puste(`Brak etapów budowy. ${BEZ_KOMENDY_ODCZYTU}`);
      return;
    }
    rama.gotowe();
  }

  function naniesPostep(postep: ProgressChangedEvent | null): void {
    if (postep === null) {
      opisPostepu.textContent =
        'Telemetria postępu (progress.changed) nie przyszła jeszcze dla tego okna.';
      wypelnienie.style.width = '0%';
      return;
    }
    const nazwa = postep.stepLabel ?? 'etap bez nazwy';
    opisPostepu.textContent =
      `${nazwa} — etap ${postep.currentStep} z ${postep.totalSteps}, ` +
      `${postep.percent}% (${postep.status})`;
    wypelnienie.style.width = `${Math.max(0, Math.min(100, postep.percent))}%`;
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/**
 * Zdanie o powodzie pustego wykazu etapów, złożone wyłącznie z liczby ramek,
 * które naprawdę przyszły w tej sesji gniazda, bez orzekania o tym, czego
 * rdzeń nie robi.
 */
function zdaniePustkiEtapow(ramki: RachunekRamek): string {
  if (ramki.wszystkie === 0) {
    return (
      'Zdarzenie apps.build.changed nie przyszło jeszcze ani razu. Wykaz etapów ' +
      'wypełniają wyłącznie jego ramki, więc pustka znaczy tu „nic jeszcze nie ' +
      'przyszło", a nie nieudany odczyt.'
    );
  }
  const ile =
    `Ramek apps.build.changed przyszło ${ramki.wszystkie}` +
    (ramki.zWdrozeniem > 0 ? `, w tym ${ramki.zWdrozeniem} z wdrożeniem` : '') +
    '.';
  if (ramki.zEtapem === 0) {
    return (
      `${ile} W żadnej pole etapu nie niosło identyfikatora, a etap bez ` +
      'identyfikatora nie jest etapem i do wykazu nie trafia — dlatego jest tu pusto.'
    );
  }
  if (ramki.zdjeteEtapy > 0) {
    return (
      `${ile} Etap z identyfikatorem niosło ${ramki.zEtapem}, a ${ramki.zdjeteEtapy} ` +
      'zdjęło go zdarzeniem usunięcia — po nich wykaz został pusty.'
    );
  }
  return `${ile} Etap z identyfikatorem niosło ${ramki.zEtapem}, a wykaz mimo to jest pusty.`;
}

/** Nagłówek jednej części okna, wyświetlany nad jej treścią w spójnym układzie typograficznym całego modułu. */
function naglowek(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}

/** Jeden etap przyniesiony zdarzeniem rdzenia, wraz z jego bieżącym stanem i stopniem ukończenia w procentach. */
function wierszEtapu(etap: AppStage): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-etapy__nazwa';
  nazwa.textContent = etap.name;

  const stan = document.createElement('span');
  stan.className = 'dn-plakietka';
  stan.textContent = etap.status;

  const element = document.createElement('li');
  element.className = 'mp-etapy__wiersz';
  element.dataset['etap'] = etap.id;
  element.append(nazwa, stan);
  return element;
}

/**
 * Jeden wpis dziennika wydań.
 *
 * Powód z pola `logRef` idzie do wiersza dosłownie: rdzeń wpisuje tam zdanie
 * o tym, dlaczego przebieg się nie powiódł, a dziennik wydań pokazujący samo
 * słowo `failed` kazałby Operatorowi zgadywać.
 */
function wierszWydania(wdrozenie: AppDeployment): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mp-wydania__wiersz';
  element.dataset['wydanie'] = wdrozenie.id;
  element.textContent =
    `${wdrozenie.version ?? 'bez numeru'} · ${wdrozenie.environment} · ${wdrozenie.status}` +
    (wdrozenie.releaseNotes === undefined ? '' : ` — ${wdrozenie.releaseNotes}`) +
    (wdrozenie.logRef === undefined || wdrozenie.logRef === ''
      ? ''
      : ` · logRef: ${wdrozenie.logRef}`);
  return element;
}

/** Panel kondycji projektu złożony wyłącznie z tego, co przyszło z rdzenia zdarzeniami tej sesji gniazda. */
function opisKondycji(
  etapy: readonly AppStage[],
  wdrozenia: readonly AppDeployment[],
): string {
  if (etapy.length === 0 && wdrozenia.length === 0) {
    return 'Kondycja projektu: rdzeń nie przysłał jeszcze ani etapu, ani wdrożenia.';
  }
  const wstrzymane = etapy.filter((etap) => etap.status === AppStageStatus.Blocked).length;
  const ostatnie = wdrozenia[0];
  const oWdrozeniu =
    ostatnie === undefined
      ? 'bez wdrożeń'
      : `ostatnie wdrożenie: ${ostatnie.environment} (${ostatnie.status})`;
  return `Kondycja projektu: etapów ${etapy.length}, wstrzymanych ${wstrzymane}, ${oWdrozeniu}.`;
}
