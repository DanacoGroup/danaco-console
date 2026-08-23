import type { AccessPoint } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { plakietkaZnaku, przyciskKarty, rozpychaczNaglowka } from './elementy-karty';
import { utworzKomunikatCzynnosci } from './komunikat-czynnosci';
import {
  adresPunktu,
  ikonaRodzaju,
  klasaStanu,
  nazwaRodzaju,
  nazwaStanu,
  nazwaTrybu,
  opisChwili,
} from './nazwy-dostepow';
import { utworzPrzelacznikTrybu } from './przelacznik-trybu';
import type { StanDostepow } from './stan-dostepow';
import { utworzPrzyciskUsuniecia } from './usuniecie-punktu';
import { utworzWyborKorzeni } from './wybor-korzeni';

/**
 * Karta jednego punktu dostępu w wykazie.
 *
 * Karta odpowiada na dwa pytania Operatora: „co to za maszyna albo katalog"
 * oraz „czy w ogóle odpowiada". Do tego daje jedną czynność — nadanie oknu
 * dostępu w wybranym trybie i na wybranych korzeniach.
 *
 * Punkt jest bytem platformy, nadanie bytem okna. Dlatego karta
 * pokazuje stan punktu wspólny dla wszystkich okien, a przycisk nadania działa
 * na oknie, z którym związana jest sekcja.
 */
export interface KartaPunktu {
  /** Element osadzany w wykazie punktów. */
  element: HTMLElement;
  /** Identyfikator punktu, który karta obsługuje. */
  punktID: string;
  /** Nanosi na kartę stan punktu i informację o nadaniu okna. */
  odswiez(punkt: AccessPoint, nadany: boolean): void;
}

export function utworzKartePunktu(punkt: AccessPoint, stan: StanDostepow): KartaPunktu {
  const identyfikator = `dd-punkt-${punkt.id}`;
  const komunikat = utworzKomunikatCzynnosci();

  const tytul = document.createElement('h4');
  tytul.className = 'dd-punkt__tytul';
  tytul.textContent = punkt.name;

  const rodzaj = znak(nazwaRodzaju(punkt.kind), 'dn-plakietka');
  const stanPunktu = znak(nazwaStanu(punkt.status), klasaStanu(punkt.status));
  const nadanie = znak('nadany temu oknu', 'dn-plakietka dn-plakietka--sygnal');
  nadanie.hidden = true;

  const adres = document.createElement('p');
  adres.className = 'dd-punkt__adres';

  const sprawdzenie = document.createElement('p');
  sprawdzenie.className = 'dd-punkt__sprawdzenie';

  const przyciskSprawdz = przyciskKarty('Sprawdź', 'dn-btn dn-btn--zarys dn-btn--sm');
  przyciskSprawdz.addEventListener('click', () => void sprawdz());

  const naglowek = document.createElement('header');
  naglowek.className = 'dd-punkt__naglowek';
  naglowek.append(
    elementIkony(ikonaRodzaju(punkt.kind), { rozmiar: 20 }),
    tytul,
    rodzaj,
    stanPunktu,
    nadanie,
    rozpychaczNaglowka('dd-punkt__rozpychacz'),
    przyciskSprawdz,
    utworzPrzyciskUsuniecia(punkt, stan, komunikat).element,
  );

  const tryb = utworzPrzelacznikTrybu({ punkt, tryb: punkt.defaultMode, identyfikator });
  const korzenie = utworzWyborKorzeni({ korzenie: punkt.roots, identyfikator });

  const przyciskNadaj = przyciskKarty('Nadaj oknu dostęp', 'dn-btn dn-btn--sygnal dn-btn--sm');
  przyciskNadaj.addEventListener('click', () => void nadaj());

  const formularz = document.createElement('div');
  formularz.className = 'dd-punkt__nadanie';
  formularz.append(tryb.element, korzenie.element, przyciskNadaj);

  const element = document.createElement('article');
  element.className = 'dn-karta dd-punkt';
  element.dataset.punkt = punkt.id;
  element.append(naglowek, adres, sprawdzenie, ...opisPunktu(punkt), formularz, komunikat.element);

  /** Sprawdzenie punktu — wynik ląduje w plakietce stanu i w zdaniu obok. */
  async function sprawdz(): Promise<void> {
    komunikat.pokaz('Pytam punkt dostępu…', true);
    const wynik = await stan.sprawdz(punkt.id);
    const tresc = wynik.wynik;
    if (wynik.udany && tresc !== undefined) {
      komunikat.pokaz(zdanieSprawdzenia(tresc.status, tresc.detail), tresc.detail === undefined);
      return;
    }
    komunikat.zWyniku(wynik, '');
  }

  /** Nadanie dostępu oknu czynnemu w trybie i na korzeniach z formularza. */
  async function nadaj(): Promise<void> {
    const wybrany = tryb.tryb();
    komunikat.pokaz('Nadaję dostęp…', true);
    const wynik = await stan.nadaj(punkt.id, wybrany, korzenie.odczytaj());
    komunikat.zWyniku(
      wynik,
      `Okno ma dostęp do punktu „${punkt.name}" w trybie ${nazwaTrybu(wybrany)}.`,
    );
  }

  const karta: KartaPunktu = {
    element,
    punktID: punkt.id,

    odswiez(nowy, nadanyOknu) {
      tryb.ustawPunkt(nowy);
      tytul.textContent = nowy.name;
      rodzaj.textContent = nazwaRodzaju(nowy.kind);
      stanPunktu.textContent = nazwaStanu(nowy.status);
      stanPunktu.className = `${klasaStanu(nowy.status)} dd-punkt__znak`;
      nadanie.hidden = !nadanyOknu;
      adres.textContent = `${nazwaRodzaju(nowy.kind)} · ${adresPunktu(nowy)}`;
      sprawdzenie.textContent = `Sprawdzony: ${opisChwili(nowy.checkedAt)} · korzenie: ${
        nowy.roots.length > 0 ? nowy.roots.join(' · ') : 'punkt ich nie wymienia'
      }`;
    },
  };

  karta.odswiez(punkt, false);
  return karta;
}

/** Zdanie po sprawdzeniu punktu; szczegół niepowodzenia idzie w całości. */
function zdanieSprawdzenia(stan: string, szczegol: string | undefined): string {
  const podstawa = `Punkt odpowiedział stanem: ${stan}.`;
  return szczegol === undefined || szczegol === '' ? podstawa : `${podstawa} ${szczegol}`;
}

/** Opis przeznaczenia punktu; brak opisu nie zostawia pustego akapitu. */
function opisPunktu(punkt: AccessPoint): HTMLElement[] {
  if (punkt.description === undefined || punkt.description === '') return [];
  const element = document.createElement('p');
  element.className = 'dd-punkt__opis';
  element.textContent = punkt.description;
  return [element];
}

/** Znak karty punktu — plakietka biblioteki w miejscu znaków tej karty. */
function znak(tresc: string, klasa: string): HTMLElement {
  return plakietkaZnaku(tresc, klasa, 'dd-punkt__znak');
}
