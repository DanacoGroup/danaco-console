import type { AppDeployment } from '../../../../shared/contract';
import { poleTekstowe, poleWielowierszowe } from '../../modele/kontrolki-formularza';
import type { PokrycieKomend } from '../pokrycie-komend';
import { utworzWykazBrakow } from './braki-kontraktu';
import { opiszPole } from './dymek-objasnienia';
import {
  BRAKI_PUBLISHER,
  KODY_OKIEN,
  KOMENDY_PROPONOWANE,
  NAZWY_OKIEN,
  OKNO_SPOZA_KATALOGU,
} from './etykiety-apps';
import { utworzRameApps } from './rama-okna';
import { narzedziaPublisherPanel } from './narzedzia-apps';
import { utworzPrzybornikApps } from './przybornik-apps';
import type { StanProduktu } from './stan-produktu';

/**
 * Publisher Panel jest panelem bocznym pakowania, manifestu i publikacji, w którym dziennik
 * wydań stoi na kontrakcie, a reszta okna czeka na komendy, których jeszcze nie ma.
 */
export interface OknoPublisherPanel {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoPublisherPanel(
  stan: StanProduktu,
  pokrycie: PokrycieKomend,
): OknoPublisherPanel {
  const kod = KODY_OKIEN.PublisherPanel;
  const rama = utworzRameApps(kod, NAZWY_OKIEN[kod] ?? kod, 'zarządca');
  rama.akcje.append(
    utworzPrzybornikApps('Pakiet, manifest, podpis i publikacja', narzedziaPublisherPanel(stan))
      .element,
  );

  const identyfikator = poleTekstowe({
    etykieta: 'Identyfikator pakietu',
    podpowiedz: 'kod stały między wydaniami',
  });
  const wersja = poleTekstowe({
    etykieta: 'Wersja semantyczna',
    podpowiedz: 'np. 1.2.0',
  });
  const uprawnienia = poleTekstowe({
    etykieta: 'Wymagane uprawnienia',
    podpowiedz: 'np. sieć, odczyt plików',
  });
  const zaleznosci = poleTekstowe({
    etykieta: 'Zależności pakietu',
    podpowiedz: 'identyfikatory po przecinku',
  });
  const deklaracja = poleWielowierszowe(
    { etykieta: 'Deklaracja udostępnianych narzędzi (JSON)' },
    4,
  );

  const zbuduj = pokrycie.przycisk(
    'Zbuduj pakiet',
    KOMENDY_PROPONOWANE.BudowaPakietu,
    'Zapakowanie produktu albo artefaktu wdrożenia w dystrybuowalne rozszerzenie z manifestem',
  );
  const zapiszManifest = pokrycie.przycisk(
    'Zapisz manifest',
    KOMENDY_PROPONOWANE.ManifestPakietu,
    'Zapis tożsamości, deklaracji narzędzi, uprawnień i zależności pakietu',
  );
  const waliduj = pokrycie.przycisk(
    'Waliduj zgodność z kontraktem',
    KOMENDY_PROPONOWANE.WalidacjaPakietu,
    'Sprawdzenie, czy pakiet spełnia jednolity kontrakt rozszerzenia; ostrzeżenia nieblokujące',
  );
  const podpisz = pokrycie.przycisk(
    'Podpisz kluczem wydawcy',
    KOMENDY_PROPONOWANE.PodpisPakietu,
    'Nadanie podpisu i sumy kontrolnej wydawcy przed publikacją',
  );
  const opublikuj = pokrycie.przycisk(
    'Opublikuj do prywatnego rejestru',
    KOMENDY_PROPONOWANE.PublikacjaPakietu,
    'Umieszczenie rozszerzenia w wewnętrznym rejestrze organizacji, widocznym w App Catalogu',
  );

  const wykazPokrycia = pokrycie.wykaz(
    [
      KOMENDY_PROPONOWANE.BudowaPakietu,
      KOMENDY_PROPONOWANE.ManifestPakietu,
      KOMENDY_PROPONOWANE.WalidacjaPakietu,
      KOMENDY_PROPONOWANE.PodpisPakietu,
      KOMENDY_PROPONOWANE.PublikacjaPakietu,
    ],
    'mp-granica mp-pokrycie',
  );

  const wydania = document.createElement('ul');
  wydania.className = 'mp-wydania';

  const oWydaniach = document.createElement('p');
  oWydaniach.className = 'dn-pole-opis';
  oWydaniach.textContent =
    'Dziennik poniżej jest dziennikiem wydań PRODUKTU — składa się z wdrożeń oddanych przez ' +
    'rdzeń, bo wersja i notatki wydania są polami wdrożenia. Dziennik wydań rozszerzenia ' +
    'to inny byt i czeka na własną komendę.';

  const pasekPakietu = document.createElement('div');
  pasekPakietu.className = 'mp-pasek';
  pasekPakietu.append(zbuduj, zapiszManifest, waliduj);

  const pasekPublikacji = document.createElement('div');
  pasekPublikacji.className = 'mp-pasek';
  pasekPublikacji.append(podpisz, opublikuj);

  const granica = document.createElement('p');
  granica.className = 'dn-pole-opis mp-granica';
  granica.textContent = OKNO_SPOZA_KATALOGU;

  rama.akcje.append(utworzWykazBrakow('Bez drogi w kontrakcie', BRAKI_PUBLISHER));
  rama.tresc.append(
    naglowekCzesci('Manifest pakietu'),
    identyfikator.element,
    opiszPole(
      wersja.element,
      'Wersjonowanie semantyczne — ten sam porządek, którym rejestr rozpoznaje aktualizację.',
    ),
    uprawnienia.element,
    zaleznosci.element,
    opiszPole(
      deklaracja.element,
      'Wykaz narzędzi i akcji, które pakiet wnosi modułom i ekspertom, wraz z ich parametrami.',
    ),
    pasekPakietu,
    naglowekCzesci('Podpis i publikacja'),
    pasekPublikacji,
    wykazPokrycia,
    naglowekCzesci('Dziennik wydań produktu'),
    oWydaniach,
    wydania,
    granica,
  );

  function odswiez(): void {
    const wykaz = stan.wdrozenia();
    wydania.replaceChildren(...wykaz.map(wierszWydania));
    if (rama.faza() === 'blad' || rama.faza() === 'ladowanie') return;
    if (wykaz.length === 0) {
      rama.puste(
        'Pakowania, walidacji, podpisu ani publikacji kontrakt nie prowadzi — pięć kontrolek ' +
          'powyżej nazywa brakujące komendy po naciśnięciu. Dziennik wydań produktu jest ' +
          'pusty, bo rdzeń nie oddał jeszcze ani jednego wdrożenia.',
      );
      return;
    }
    rama.puste(
      `Dziennik wydań produktu liczy ${wykaz.length} pozycji. Pakowania, walidacji, podpisu ` +
        'ani publikacji kontrakt nie prowadzi — pięć kontrolek powyżej nazywa brakujące ' +
        'komendy po naciśnięciu.',
    );
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Jeden wpis dziennika wydań produktu złożony wyłącznie z pól wdrożenia, bez pól właściwych rozszerzeniu. */
function wierszWydania(wdrozenie: AppDeployment): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mp-wydania__wiersz';
  element.dataset['wydanie'] = wdrozenie.id;
  element.textContent =
    `${wdrozenie.version ?? 'bez numeru'} · ${wdrozenie.environment} · ${wdrozenie.status}` +
    (wdrozenie.releaseNotes === undefined || wdrozenie.releaseNotes === ''
      ? ' — bez notatek wydania'
      : ` — ${wdrozenie.releaseNotes}`);
  return element;
}

/** Nagłówek części okna nazywa kolejny fragment panelu, oddzielając wizualnie grupy pól formularza od siebie. */
function naglowekCzesci(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mp-czesc__tytul';
  element.textContent = tresc;
  return element;
}
