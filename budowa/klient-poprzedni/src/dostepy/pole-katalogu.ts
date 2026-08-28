import { adresWpisu, opisAdresu } from '../konfiguracja/adres-ustawienia';
import type { Rozstrzygniecie } from '../konfiguracja/rozstrzygniecie';
import { elementIkony } from '../ikony/ikony';
import { Cel, czyPowlokaNatywna, wskazKatalog } from './dialog-katalogu';
import { utworzKomunikatCzynnosci } from './komunikat-czynnosci';
import type { StanKataloguRoboczego } from './stan-katalogu-roboczego';

/** Jedno pole obszaru katalogu roboczego, pokazujące wartość dziś obowiązującą wraz z jej pochodzeniem. */
export interface PoleKatalogu {
  /** Wiersz osadzany w obszarze katalogu roboczego. */
  element: HTMLElement;
  /** Nanosi wartość obowiązującą i jej pochodzenie. */
  odswiez(): void;
}

export interface ZaleznosciPolaKatalogu {
  /** Klucz ustawienia. */
  klucz: string;
  /** Stan obszaru — przez niego idą odczyt i zapis. */
  stan: StanKataloguRoboczego;
  /** Czy pole ma przycisk natywnego wskazania katalogu. */
  zWyborem: boolean;
}

export function utworzPoleKatalogu(zaleznosci: ZaleznosciPolaKatalogu): PoleKatalogu {
  const { klucz, stan, zWyborem } = zaleznosci;
  const identyfikator = `dd-katalog-${klucz.replace(/[^\w-]/gu, '-')}`;
  const komunikat = utworzKomunikatCzynnosci();

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta dd-katalog__etykieta';
  etykieta.htmlFor = identyfikator;

  const nazwaKlucza = document.createElement('code');
  nazwaKlucza.className = 'dd-katalog__klucz';
  nazwaKlucza.textContent = klucz;

  const pole = document.createElement('input');
  pole.type = 'text';
  pole.id = identyfikator;
  pole.className = 'dd-katalog__pole';
  pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'Enter') void zapisz();
  });

  const wybor = document.createElement('button');
  wybor.type = 'button';
  wybor.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  wybor.append(elementIkony('folder', { rozmiar: 16 }), document.createTextNode('Wskaż katalog'));
  wybor.addEventListener('click', () => void wskaz());

  const zapiszPrzycisk = przycisk('Zapisz', 'dn-btn dn-btn--atrament dn-btn--sm', () =>
    void zapisz(),
  );
  const przywrocPrzycisk = przycisk('Przywróć domyślny', 'dn-btn dn-btn--duch dn-btn--sm', () =>
    void przywroc(),
  );

  const wiersz = document.createElement('div');
  wiersz.className = 'dd-katalog__wiersz';
  wiersz.append(pole, ...(zWyborem ? [wybor] : []), zapiszPrzycisk, przywrocPrzycisk);

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis dd-katalog__opis';

  const pochodzenie = document.createElement('p');
  pochodzenie.className = 'dd-katalog__pochodzenie';

  const element = document.createElement('div');
  element.className = 'dd-katalog__pole-wiersz';
  element.dataset.klucz = klucz;
  element.append(etykieta, nazwaKlucza, wiersz, opis, pochodzenie, komunikat.element);

  /** Wskazanie katalogu oknem powłoki; poza powłoką zdanie zamiast ciszy. */
  async function wskaz(): Promise<void> {
    if (!czyPowlokaNatywna()) {
      komunikat.pokaz(
        'Natywne okno wyboru należy do powłoki Danaco Console. W przeglądarce wpisz ścieżkę w polu obok.',
        false,
      );
      pole.focus();
      return;
    }
    const sciezka = await wskazKatalog(Cel.KatalogRoboczy);
    if (sciezka === '') {
      komunikat.pokaz('Nie wskazano katalogu — nic się nie zmieniło.', true);
      return;
    }
    pole.value = sciezka;
    await zapisz();
  }

  async function zapisz(): Promise<void> {
    const wartosc = pole.value.trim();
    if (wartosc === '') {
      komunikat.pokaz(
        'Puste pole nie jest zapisem. Aby wrócić do wartości domyślnej, użyj przycisku „Przywróć domyślny".',
        false,
      );
      return;
    }
    const wynik = await stan.zapisz(klucz, wartosc);
    komunikat.zWyniku(wynik, 'Zapisano na poziomie globalnym, na osi platformy.');
  }

  async function przywroc(): Promise<void> {
    const wynik = await stan.przywroc(klucz);
    komunikat.zWyniku(
      wynik,
      'Zapis zdjęty — obowiązuje wartość domyślna ustalana przez rdzeń.',
    );
  }

  return {
    element,

    odswiez() {
      const definicja = stan.definicja(klucz);
      const rozstrzygniecie = stan.rozstrzygnij(klucz);
      etykieta.textContent = definicja?.name ?? klucz;
      opis.textContent = definicja?.description ?? '';
      opis.hidden = opis.textContent === '';
      pole.placeholder = definicja?.placeholder ?? '';

      if (rozstrzygniecie === null) {
        pochodzenie.textContent = 'Katalog ustawień nie zna tego klucza.';
        return;
      }

      if (document.activeElement !== pole) {
        pole.value = napis(rozstrzygniecie.wartosc);
      }
      pochodzenie.textContent = zdaniePochodzenia(rozstrzygniecie);
      pochodzenie.dataset.domyslna = String(rozstrzygniecie.domyslna);
    },
  };
}

/** Zdanie o pochodzeniu wartości pola: świadomy zapis użytkownika albo wartość domyślna dla całego rdzenia. */
function zdaniePochodzenia(rozstrzygniecie: Rozstrzygniecie): string {
  const zrodlo = rozstrzygniecie.zrodlo;
  if (rozstrzygniecie.domyslna || zrodlo === null) {
    return 'Obowiązuje wartość domyślna — nikt tego klucza nie nadpisał.';
  }
  return `Wartość pochodzi z zapisu: ${opisAdresu(adresWpisu(zrodlo))}.`;
}

/** Napis z wartości nieznanego kształtu odczytanej z ustawienia; brak wartości daje tutaj zawsze napis pusty. */
function napis(wartosc: unknown): string {
  if (wartosc === null || wartosc === undefined) return '';
  return typeof wartosc === 'string' ? wartosc : String(wartosc);
}

function przycisk(
  tresc: string,
  klasa: string,
  przyNacisnieciu: () => void,
): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = tresc;
  element.addEventListener('click', przyNacisnieciu);
  return element;
}
