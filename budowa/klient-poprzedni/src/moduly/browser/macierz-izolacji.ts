import type { IsolationPolicy, IsolationScopeLevel } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import {
  ETYKIETY_ZASIEGU,
  POZYCJE_KONTEKSTU,
  POZYCJE_TECHNICZNE,
} from '../../punkty-izolacji/katalog-izolacji';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Macierz izolacji sesji — panel warstwy czwartej modułu Browser, pokazujący
 * punkty izolacji właściwe oknu przeglądarki.
 */
export interface MacierzIzolacji {
  element: HTMLElement;
  /** Czyta politykę obowiązującą oraz poziomy zasięgu okna przeglądarki. */
  odczytaj(): Promise<void>;
}

export function utworzMacierzIzolacji(stan: StanPrzegladania): MacierzIzolacji {
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz = (tresc: string, ok: boolean): void => odpowiedz.pokaz(tresc, ok);

  const naglowek = document.createElement('h4');
  naglowek.className = 'mb-panel__tytul';
  naglowek.textContent = 'Macierz izolacji sesji';

  const polityka = document.createElement('ul');
  polityka.className = 'mb-izolacja__lista';

  const poziomy = document.createElement('ul');
  poziomy.className = 'mb-izolacja__lista';

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  uwaga.textContent =
    'Panel jest odczytem. Przełączniki, profile i przypisanie do poziomu zasięgu zapisuje ' +
    'okno konfiguracji punktów izolacji.';

  async function odczytaj(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      powiedz(stan.powod(), false);
      return;
    }
    powiedz('Odczyt polityki izolacji obowiązującej dla okna przeglądarki…', true);
    const [wynikPolityki, wynikPoziomow] = await Promise.all([
      stan.izolacja.polityka({ windowId: idOkna }),
      stan.izolacja.poziomy({ windowId: idOkna }),
    ]);

    if (!wynikPolityki.udany || wynikPolityki.wynik === undefined) {
      powiedz(
        opisOdmowy('Odczyt polityki izolacji', wynikPolityki.blad?.code, wynikPolityki.blad?.message),
        false,
      );
      return;
    }
    polityka.replaceChildren(...wierszePolityki(wynikPolityki.wynik.policy));

    if (!wynikPoziomow.udany || wynikPoziomow.wynik === undefined) {
      // Brak wykazu poziomów nie unieważnia odczytanych przełączników — panel mówi o obu odczytach osobno.
      poziomy.replaceChildren();
      powiedz(
        opisOdmowy('Odczyt poziomów zasięgu', wynikPoziomow.blad?.code, wynikPoziomow.blad?.message),
        false,
      );
      return;
    }
    poziomy.replaceChildren(...wierszePoziomow(wynikPoziomow.wynik.scopes));
    powiedz(zdanieOPolityce(wynikPolityki.wynik.policy), true);
  }

  const pasek = document.createElement('div');
  pasek.className = 'mb-panel__pasek';
  pasek.append(
    przyciskCzynnosci('Odczytaj politykę izolacji', KLASA_PRZYCISKU.zarys, () => void odczytaj()),
    utworzDymekObjasnienia(OBJASNIENIA.macierzIzolacji, KLASY_DYMKA),
  );

  const element = document.createElement('section');
  element.className = 'mb-panel mb-izolacja';
  element.setAttribute('aria-label', 'Macierz izolacji sesji — warstwa czwarta');
  element.append(naglowek, pasek, polityka, poziomy, uwaga, odpowiedz.element);

  return { element, odczytaj };
}

/** Zdanie o pochodzeniu polityki izolacji — poziom, warstwa i profil, z którego ta polityka faktycznie wyszła. */
function zdanieOPolityce(polityka: IsolationPolicy): string {
  const profil = (polityka.profileId ?? '').trim();
  const zrodlo = polityka.origin === undefined ? '' : ` (odziedziczona z poziomu ${ETYKIETY_ZASIEGU[polityka.origin]})`;
  return (
    `Polityka obowiązująca na poziomie „${ETYKIETY_ZASIEGU[polityka.scope]}", ` +
    `warstwa ${polityka.layer}${zrodlo}` +
    `${profil === '' ? '' : `, z profilu ${profil}`}.`
  );
}

/** Przełączniki polityki izolacji — kontekst i zakresy techniczne zebrane razem w jednym wspólnym wykazie. */
function wierszePolityki(polityka: IsolationPolicy): HTMLElement[] {
  const kontekst = POZYCJE_KONTEKSTU.map((pozycja) => {
    const przelacznik = polityka.contextSwitches.find((wpis) => wpis.kind === pozycja.kind);
    return wierszPrzelacznika(pozycja.nazwa, przelacznik?.isolated, przelacznik?.explanation);
  });
  const techniczne = POZYCJE_TECHNICZNE.map((pozycja) => {
    const przelacznik = polityka.technicalSwitches.find((wpis) => wpis.scope === pozycja.scope);
    return wierszPrzelacznika(pozycja.nazwa, przelacznik?.isolated, przelacznik?.explanation);
  });
  return [...kontekst, ...techniczne];
}

/**
 * Jeden przełącznik: nazwa, stan słowem i objaśnienie rdzenia. Przełącznik,
 * którego rdzeń w polityce nie oddał, mówi wprost, że polityka o nim milczy,
 * zamiast przedstawiać się jako wspólny.
 */
function wierszPrzelacznika(
  nazwa: string,
  odciety: boolean | undefined,
  objasnienie: string | undefined,
): HTMLElement {
  const element = document.createElement('li');
  element.className = 'mb-izolacja__wiersz';
  element.dataset['stan'] = znacznikStanu(odciety);
  element.textContent =
    `${nazwa} — ${slowoStanu(odciety)}${objasnienie === undefined ? '' : ` · ${objasnienie}`}`;
  return element;
}

/** Stan przełącznika izolacji wyrażony słowem, jako treść widoczna wprost dla operatora panelu izolacji. */
function slowoStanu(odciety: boolean | undefined): string {
  if (odciety === undefined) return 'polityka o nim milczy';
  return odciety ? 'odcięty' : 'wspólny';
}

/** Ten sam stan przełącznika wyrażony wartością `data-stan` — dla arkusza stylu i dla sprawdzianu automatycznego. */
function znacznikStanu(odciety: boolean | undefined): string {
  if (odciety === undefined) return 'nieznany';
  return odciety ? 'odciety' : 'wspolny';
}

/** Poziomy zasięgu izolacji w kolejności ich rozstrzygania, z zaznaczeniem poziomu najwęższego zasięgu. */
function wierszePoziomow(poziomy: readonly IsolationScopeLevel[]): HTMLElement[] {
  return poziomy.map((poziom) => {
    const element = document.createElement('li');
    element.className = 'mb-izolacja__wiersz';
    element.dataset['stan'] = poziom.narrowest === true ? 'rozstrzygajacy' : 'posredni';
    element.textContent =
      `${poziom.order}. ${poziom.label}` +
      `${poziom.narrowest === true ? ' — poziom rozstrzygający' : ''}` +
      `${poziom.description === undefined ? '' : ` · ${poziom.description}`}`;
    return element;
  });
}
