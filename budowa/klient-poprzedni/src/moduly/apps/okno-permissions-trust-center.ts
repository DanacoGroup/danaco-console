import type {
  Extension,
  IsolationPolicy,
  IsolationTechnicalSwitch,
} from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { utworzWykazBrakow } from './braki-kontraktu';
import {
  BRAKI_PERMISSIONS,
  KODY_OKIEN,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import { nazwaPochodzenia, opisStanu } from './karta-rozszerzenia';
import { utworzRameApps } from './rama-okna';
import { narzedziaZaufania } from './narzedzia-rozszerzen';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanRozszerzen } from './stan-rozszerzen';
import type { ZrodloIzolacjiApps } from './zrodlo-izolacji-apps';

/**
 * Permissions & Trust Center jest panelem bocznym uprawnień, izolacji i pochodzenia, otwieranym
 * na pozycji wskazanej gdzie indziej, bo panel nie prowadzi własnego wyboru.
 */
export interface OknoPermissionsTrustCenter {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * Napisy ośmiu zakresów technicznych w języku Operatora, bo wartości kontraktu niosą wyłącznie
 * kody wewnętrzne platformy.
 */
const NAZWY_ZAKRESOW: Readonly<Record<string, string>> = {
  workingDirectory: 'katalog roboczy',
  processEnvironment: 'środowisko procesu',
  modelDataDirectory: 'katalog danych i konfiguracji modelu',
  networkAccess: 'dostęp sieciowy',
  fileAccess: 'odczyt i zapis plików',
  accountToken: 'konto i token per sesja',
  processModel: 'model procesu',
  executionServer: 'serwer wykonania',
};

export function utworzOknoPermissionsTrustCenter(
  stan: StanRozszerzen,
  izolacja: ZrodloIzolacjiApps,
  idOkna: () => string,
): OknoPermissionsTrustCenter {
  const kod = KODY_OKIEN.PermissionsTrustCenter;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'zarządca');

  let polityka: IsolationPolicy | null = null;
  let powodPolityki = '';

  const tozsamosc = document.createElement('dl');
  tozsamosc.className = 'mp-tozsamosc';

  const zaufanie = document.createElement('p');
  zaufanie.className = 'mp-zaufanie';

  const konfiguracja = document.createElement('pre');
  konfiguracja.className = 'mp-podglad';

  const naglowekIzolacji = naglowekCzesci('Izolacja techniczna platformy');

  const wstepIzolacji = document.createElement('p');
  wstepIzolacji.className = 'dn-pole-opis';
  wstepIzolacji.textContent =
    'Osiem zakresów poniżej opisuje warunki, w jakich platforma wykonuje kod — i w tych ' +
    'samych warunkach wykona się kod rozszerzenia. Nie jest to izolacja nadana tej ' +
    'pozycji: kontraktu izolacji per rozszerzenie nie ma.';

  const zakresy = document.createElement('ul');
  zakresy.className = 'mp-zakresy';

  const odczytaj = przycisk('Odczytaj politykę izolacji', 'dn-btn dn-btn--sm dn-btn--zarys');
  const odpowiedz = utworzWierszOdpowiedzi();

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  const bezBlokad = document.createElement('p');
  bezBlokad.className = 'dn-pole-opis mp-bez-blokad';
  bezBlokad.textContent =
    'Żadne ostrzeżenie ani brak podpisu nie wstrzymuje instalacji ani włączenia pozycji. ' +
    'Kontrola zostaje po stronie Operatora: przez świadome włączenie i przez zakres uprawnień.';

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_PERMISSIONS, 'extension'));
  rama.akcje.append(
    utworzPrzybornikApps('Uprawnienia, podpis, skaner i sekrety', narzedziaZaufania(stan)).element,
  );
  rama.tresc.append(
    naglowekCzesci('Tożsamość i pochodzenie pozycji'),
    tozsamosc,
    zaufanie,
    naglowekCzesci('Konfiguracja pozycji tak, jak oddał ją rdzeń'),
    konfiguracja,
    naglowekIzolacji,
    wstepIzolacji,
    odczytaj,
    zakresy,
    odpowiedz.element,
    bezBlokad,
    granica,
  );

  odczytaj.addEventListener('click', () => void odczytajPolityke());

  async function odczytajPolityke(): Promise<void> {
    rama.ladowanie('Odczyt polityki izolacji…');
    powodPolityki = '';
    const wynik = await izolacja.polityka(idOkna());
    if (!wynik.udany || wynik.wynik === undefined) {
      powodPolityki = opisOdmowyBledu('Odczyt polityki izolacji', wynik.blad);
      rama.blad(powodPolityki);
      odpowiedz.pokaz(powodPolityki, false);
      odswiez();
      return;
    }
    polityka = wynik.wynik.policy;
    rama.gotowe();
    odswiez();
    const skad =
      polityka.origin === undefined
        ? 'bez wskazania poziomu, z którego wartość odziedziczono'
        : `wartość odziedziczona z poziomu ${polityka.origin}`;
    odpowiedz.pokaz(
      `Rdzeń oddał politykę dla poziomu ${polityka.scope}, warstwa ${polityka.layer} — ` +
        `${skad}. Przełączników technicznych: ${polityka.technicalSwitches.length}, ` +
        `przełączników kontekstu: ${polityka.contextSwitches.length}.`,
      true,
    );
  }

  /** Wiersz tożsamości: podpis i wartość. */
  function wierszTozsamosci(podpis: string, wartosc: string): readonly HTMLElement[] {
    const nazwa = document.createElement('dt');
    nazwa.textContent = podpis;
    const tresc = document.createElement('dd');
    tresc.textContent = wartosc;
    return [nazwa, tresc];
  }

  /** Zdanie o zaufaniu orzeka wyłącznie o pochodzeniu, bo kontrakt nie niesie podpisu cyfrowego pozycji. */
  function zdanieZaufania(pozycja: Extension): string {
    return (
      `Źródło pochodzenia: ${nazwaPochodzenia(pozycja.origin)}. Podpis cyfrowy i suma ` +
      'kontrolna pakietu nie są polami pozycji katalogu, więc okno NIE orzeka, czy pakiet ' +
      'jest podpisany ani czy wydawca został zweryfikowany. Pochodzenie rozstrzyga ' +
      'wyłącznie stan wyjściowy przy rejestracji.'
    );
  }

  /** Jeden zakres izolacji, nazwa i stan słowem; objaśnienie idzie z odpowiedzi rdzenia, nie stąd. */
  function wierszZakresu(przelacznik: IsolationTechnicalSwitch): HTMLElement {
    const nazwa = document.createElement('span');
    nazwa.className = 'mp-zakresy__nazwa';
    nazwa.textContent = NAZWY_ZAKRESOW[przelacznik.scope] ?? przelacznik.scope;

    const stanZakresu = document.createElement('span');
    stanZakresu.className = przelacznik.isolated
      ? 'dn-plakietka dn-plakietka--sukces'
      : 'dn-plakietka dn-plakietka--ostrzezenie';
    stanZakresu.textContent = przelacznik.isolated ? 'zakres odcięty' : 'zakres wspólny';

    const element = document.createElement('li');
    element.className = 'mp-zakresy__wiersz';
    element.dataset['zakres'] = przelacznik.scope;
    element.dataset['odciety'] = String(przelacznik.isolated);
    element.append(nazwa, stanZakresu);
    if (przelacznik.explanation !== undefined && przelacznik.explanation !== '') {
      element.append(
        utworzDymekObjasnienia(przelacznik.explanation, {
          powloka: 'mp-dymek',
          znak: 'mp-dymek__znak',
        }),
      );
    }
    return element;
  }

  function odswiez(): void {
    // Zakresy izolacji rysują się niezależnie od wskazania pozycji, bo opisują platformę, nie pozycję.
    zakresy.replaceChildren(...(polityka?.technicalSwitches ?? []).map(wierszZakresu));
    const pozycja = stan.wybrane();
    if (pozycja === null) {
      tozsamosc.replaceChildren();
      zaufanie.textContent = '';
      konfiguracja.textContent = '';
      if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
      rama.puste(
        'Nie wskazano pozycji. Naciśnij „Uprawnienia i zaufanie" na karcie w App Catalogu ' +
          'albo w wierszu Integrations Hubu — panel otworzy się na tej pozycji. ' +
          (polityka === null
            ? 'Politykę izolacji platformy da się odczytać już teraz — nie zależy od wskazania.'
            : `Polityka izolacji platformy odczytana: zakresów ${polityka.technicalSwitches.length}.`),
      );
      return;
    }
    const punkt = pozycja.accessPointId === undefined ? null : stan.punkt(pozycja.accessPointId);
    tozsamosc.replaceChildren(
      ...wierszTozsamosci('Nazwa', pozycja.name),
      ...wierszTozsamosci('Kod', pozycja.code),
      ...wierszTozsamosci('Rodzaj', pozycja.kind),
      ...wierszTozsamosci('Wersja', pozycja.version ?? 'nieoddana przez rdzeń'),
      ...wierszTozsamosci('Stan', opisStanu(pozycja)),
      ...wierszTozsamosci(
        'Zakres dostępu przez most',
        punkt === null
          ? 'Pozycja nie wskazuje mostu — zakresu korzeni nie ma z czego wyliczyć.'
          : `${punkt.name}: ${punkt.roots.length === 0 ? 'bez korzeni' : punkt.roots.join(', ')} ` +
            `(tryb proponowany: ${punkt.defaultMode})`,
      ),
    );
    zaufanie.textContent = zdanieZaufania(pozycja);
    konfiguracja.textContent =
      pozycja.config === undefined || pozycja.config === null
        ? 'Rdzeń nie oddał konfiguracji tej pozycji. Nie jest to wykaz uprawnień — takiego ' +
          'wykazu kontrakt nie niesie.'
        : JSON.stringify(pozycja.config, null, 2);
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (polityka === null && powodPolityki === '') {
      rama.puste(
        `Pozycja ${pozycja.name} wskazana. Polityka izolacji nie była jeszcze czytana — ` +
          'naciśnij „Odczytaj politykę izolacji".',
      );
      return;
    }
    rama.gotowe();
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Nagłówek części okna nazywa kolejny fragment panelu, oddzielając wizualnie grupy pól formularza od siebie. */
function naglowekCzesci(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}
