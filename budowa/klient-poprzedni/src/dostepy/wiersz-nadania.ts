import {
  AccessPointKind,
  AccessPointStatus,
  type AccessGrant,
  type AccessPoint,
} from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { rozpychaczNaglowka } from './elementy-karty';
import { utworzKomunikatCzynnosci } from './komunikat-czynnosci';
import { adresPunktu, nazwaRodzaju, nazwaTrybu } from './nazwy-dostepow';
import { utworzPrzelacznikTrybu } from './przelacznik-trybu';
import type { StanDostepow } from './stan-dostepow';
import { utworzWyborKorzeni } from './wybor-korzeni';

/**
 * Wiersz jednego nadania w zbiorze okna.
 *
 * Okno ma zbiór nadań, a kolejność i oznaczenie głównego mają znaczenie —
 * dlatego wiersz niesie cztery czynności, nie jedną: przestawienie w górę
 * i w dół, oznaczenie głównym, zmianę trybu wraz z korzeniami oraz odebranie.
 *
 * Każda z nich idzie osobną komendą i wraca kompletem nadań okna, więc wiersz
 * nie zgaduje, jak przestawiły się pozostałe — dostaje je z rdzenia.
 */
export interface WierszNadania {
  /** Element osadzany w liście nadań. */
  element: HTMLElement;
  /** Identyfikator nadania obsługiwanego przez wiersz. */
  nadanieID: string;
}

export interface ZaleznosciWiersza {
  /** Nadanie w kształcie z kontraktu. */
  nadanie: AccessGrant;
  /** Punkt, na który nadanie się powołuje; `null`, gdy wykaz go nie zna. */
  punkt: AccessPoint | null;
  /** Stan sekcji — przez niego idą wszystkie cztery czynności. */
  stan: StanDostepow;
  /** Pozycja wiersza liczona od 1 oraz liczba nadań okna. */
  pozycja: number;
  liczba: number;
}

export function utworzWierszNadania(zaleznosci: ZaleznosciWiersza): WierszNadania {
  const { nadanie, punkt, stan, pozycja, liczba } = zaleznosci;
  const identyfikator = `dd-nadanie-${nadanie.id}`;
  const komunikat = utworzKomunikatCzynnosci();

  const element = document.createElement('li');
  element.className = 'dd-nadanie';
  element.dataset.nadanie = nadanie.id;

  const numer = document.createElement('span');
  numer.className = 'dd-nadanie__numer';
  numer.textContent = String(pozycja);

  const tytul = document.createElement('span');
  tytul.className = 'dd-nadanie__tytul';
  tytul.textContent = punkt === null ? nadanie.accessPointId : punkt.name;

  const opis = document.createElement('p');
  opis.className = 'dd-nadanie__opis';
  opis.textContent =
    punkt === null
      ? 'Punkt dostępu nie występuje w wykazie — mógł zostać usunięty albo wykaz jeszcze nie dotarł.'
      : `${nazwaRodzaju(punkt.kind)} · ${adresPunktu(punkt)}`;

  const glowne = document.createElement('button');
  glowne.type = 'button';
  glowne.className = nadanie.primary
    ? 'dn-btn dn-btn--sygnal dn-btn--sm'
    : 'dn-btn dn-btn--duch dn-btn--sm';
  glowne.append(elementIkony('gwiazdka', { rozmiar: 14 }));
  glowne.append(document.createTextNode(nadanie.primary ? 'Główne' : 'Ustaw głównym'));
  glowne.title = nadanie.primary
    ? 'To nadanie jest głównym nadaniem okna'
    : 'Oznacz to nadanie jako główne okna';
  glowne.addEventListener('click', () => void oznaczGlownym());

  const wGore = przyciskPrzestawienia('grot-gora', 'Przesuń wyżej', () =>
    void przestaw(pozycja - 1),
  );
  const wDol = przyciskPrzestawienia('grot-dol', 'Przesuń niżej', () =>
    void przestaw(pozycja + 1),
  );

  const odbierz = document.createElement('button');
  odbierz.type = 'button';
  odbierz.className = 'dn-btn dn-btn--niebezpieczny dn-btn--sm';
  odbierz.textContent = 'Odbierz';
  odbierz.setAttribute('aria-label', `Odbierz nadanie ${tytul.textContent ?? ''}`);
  odbierz.addEventListener('click', () => void odbierzNadanie());

  const naglowek = document.createElement('div');
  naglowek.className = 'dd-nadanie__naglowek';
  naglowek.append(numer, tytul, rozpychaczNaglowka('dd-nadanie__rozpychacz'), glowne, wGore, wDol, odbierz);

  const tryb = utworzPrzelacznikTrybu({
    punkt: punkt ?? zastepczyPunkt(nadanie),
    tryb: nadanie.mode,
    identyfikator,
  });
  const korzenie = utworzWyborKorzeni({
    korzenie: punkt?.roots ?? nadanie.roots ?? [],
    identyfikator,
  });
  korzenie.ustaw(nadanie.roots ?? []);

  const zapisz = document.createElement('button');
  zapisz.type = 'button';
  zapisz.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  zapisz.textContent = 'Zapisz tryb i korzenie';
  zapisz.addEventListener('click', () => void zapiszZakres());

  const zakres = document.createElement('div');
  zakres.className = 'dd-nadanie__zakres';
  zakres.append(tryb.element, korzenie.element, zapisz);

  element.append(naglowek, opis, zakres, komunikat.element);

  // Skrajne pozycje nie chowają przycisków i nie wygaszają ich:
  // naciśnięcie daje odpowiedź, że wiersz już jest pierwszy albo ostatni.
  async function przestaw(nowa: number): Promise<void> {
    if (nowa < 1 || nowa > liczba) {
      komunikat.pokaz(
        nowa < 1 ? 'To nadanie jest już pierwsze.' : 'To nadanie jest już ostatnie.',
        true,
      );
      return;
    }
    const wynik = await stan.zmienNadanie({ grantId: nadanie.id, order: nowa });
    komunikat.zWyniku(wynik, `Nadanie stoi teraz na pozycji ${nowa}.`);
  }

  async function oznaczGlownym(): Promise<void> {
    if (nadanie.primary) {
      komunikat.pokaz('To nadanie jest już głównym nadaniem okna.', true);
      return;
    }
    const wynik = await stan.zmienNadanie({ grantId: nadanie.id, primary: true });
    komunikat.zWyniku(wynik, 'Nadanie jest głównym nadaniem okna.');
  }

  async function zapiszZakres(): Promise<void> {
    const wybrany = tryb.tryb();
    const wynik = await stan.zmienNadanie({
      grantId: nadanie.id,
      mode: wybrany,
      roots: korzenie.odczytaj(),
    });
    komunikat.zWyniku(wynik, `Zapisano tryb ${nazwaTrybu(wybrany)} wraz z korzeniami.`);
  }

  async function odbierzNadanie(): Promise<void> {
    const wynik = await stan.odbierz(nadanie.id);
    komunikat.zWyniku(wynik, 'Nadanie odebrane — model nie sięgnie tu z tego okna.');
  }

  return { element, nadanieID: nadanie.id };
}

/**
 * Punkt zastępczy dla nadania, którego punktu nie ma w wykazie.
 *
 * Wiersz musi dać się narysować także wtedy, gdy wykaz punktów jeszcze nie
 * dotarł — inaczej nadanie zniknęłoby z widoku, choć w rdzeniu istnieje.
 * Ostrzeżenie o zapisie liczone jest wtedy z pustej nazwy maszyny, czyli nie
 * pojawia się; pojawi się po dojściu wykazu.
 */
function zastepczyPunkt(nadanie: AccessGrant): AccessPoint {
  return {
    id: nadanie.accessPointId,
    name: nadanie.accessPointId,
    kind: AccessPointKind.LocalDirectory,
    roots: nadanie.roots ?? [],
    defaultMode: nadanie.mode,
    status: AccessPointStatus.Unknown,
    enabled: true,
    createdAt: nadanie.createdAt,
    updatedAt: nadanie.createdAt,
  };
}

function przyciskPrzestawienia(
  ikona: 'grot-gora' | 'grot-dol',
  opis: string,
  przyNacisnieciu: () => void,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-btn-ikona dd-nadanie__przestaw';
  element.title = opis;
  element.setAttribute('aria-label', opis);
  element.append(elementIkony(ikona, { rozmiar: 16 }));
  element.addEventListener('click', przyNacisnieciu);
  return element;
}
