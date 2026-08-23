import type { ConfigEntry, SettingDefinition } from '../../../shared/contract';
import { adnotacjaKlucza } from '../sterowanie/adnotacje-wykonania';
import { adresWpisu, opisAdresu } from './adres-ustawienia';
import { utworzDymekObjasnienia, zlozObjasnienie } from './dymek-objasnienia';
import type { Kontrolka } from './kontrolka';
import type { StanKonfiguracji } from './stan-konfiguracji';
import { utworzKontrolke } from './wybor-kontrolki';
import { utworzWskaznikZasiegu } from './wskaznik-zasiegu';

/**
 * Wiersz jednego ustawienia: etykieta, kontrolka, wskaźnik zasięgu, opis.
 *
 * Wiersz nie zna ani jednego klucza z osobna — wszystko, co rysuje, pochodzi
 * z pozycji katalogu przekazanej w zależnościach. Rodzaj kontrolki
 * rozstrzyga jeden moduł rozdzielający, więc nowy rodzaj wartości nie dotyka
 * tego pliku.
 *
 * Zatwierdzenie kontrolki idzie komendą `config.set` pod adres wskazany we
 * wskaźniku zasięgu — domyślnie pod punkt widzenia okna, po zmianie poziomu
 * pod poziom wybrany przez Operatora. Wynik zapisu widnieje przy polu, żeby
 * naciśnięcie zawsze dało odpowiedź: powodzenie mówi, gdzie zapisano,
 * niepowodzenie mówi, co odpowiedział rdzeń.
 */
export interface PoleUstawienia {
  /** Wiersz osadzany w panelu kategorii. */
  element: HTMLElement;
  /** Klucz pozycji katalogu, którą wiersz obsługuje. */
  klucz: string;
  /** Nanosi wartość obowiązującą i pochodzenie ze stanu okna. */
  odswiez(): void;
  /** Pokazuje albo chowa wiersz — reguła `visibleWhenKey` katalogu. */
  ustawWidocznosc(widoczne: boolean): void;
}

export interface ZaleznosciPola {
  definicja: SettingDefinition;
  stan: StanKonfiguracji;
}

export function utworzPoleUstawienia(zaleznosci: ZaleznosciPola): PoleUstawienia {
  const { definicja, stan } = zaleznosci;
  const identyfikator = `dk-pole-${definicja.key.replace(/[^\w-]/gu, '-')}`;

  const kontrolka = utworzKontrolke({ definicja, identyfikator });
  const wskaznik = utworzWskaznikZasiegu({
    definicja,
    naPrzywrocenie: (wpis) => void przywroc(wpis),
  });

  const stanZapisu = document.createElement('p');
  stanZapisu.className = 'dk-pole__stan';
  stanZapisu.hidden = true;

  const element = document.createElement('div');
  element.className = 'dk-pole';
  element.dataset.klucz = definicja.key;
  element.append(
    naglowekPola(definicja, identyfikator, wskaznik.element),
    wierszKontrolki(definicja, kontrolka),
    ...opisyPola(definicja, kontrolka),
    stanZapisu,
  );

  function pokaz(tresc: string, powodzenie: boolean): void {
    stanZapisu.textContent = tresc;
    stanZapisu.hidden = tresc === '';
    stanZapisu.dataset.powodzenie = String(powodzenie);
  }

  async function zapisz(): Promise<void> {
    const adres = wskaznik.adres();
    const wynik = await stan.zapisz(definicja.key, kontrolka.odczytaj(), adres);
    pokaz(
      wynik.udany
        ? `Zapisano na poziomie: ${opisAdresu(adres)}`
        : `Rdzeń nie przyjął zapisu: ${wynik.blad?.message ?? 'brak treści błędu'}`,
      wynik.udany,
    );
  }

  async function przywroc(wpis: ConfigEntry): Promise<void> {
    const wynik = await stan.przywroc(definicja.key, adresWpisu(wpis));
    pokaz(
      wynik.udany
        ? 'Zapis zdjęty — obowiązuje wartość z poziomu szerszego.'
        : `Rdzeń nie zdjął zapisu: ${wynik.blad?.message ?? 'brak treści błędu'}`,
      wynik.udany,
    );
  }

  kontrolka.naZatwierdzenie(() => void zapisz());
  wskaznik.naZmianeCelu(() => pokaz('', true));

  return {
    element,
    klucz: definicja.key,

    odswiez() {
      const rozstrzygniecie = stan.rozstrzygnij(definicja);
      wskaznik.odswiez(rozstrzygniecie, stan.punkt());
      if (!element.contains(document.activeElement)) {
        kontrolka.ustaw(rozstrzygniecie.wartosc);
      }
    },

    ustawWidocznosc(widoczne) {
      element.hidden = !widoczne;
    },
  };
}

/** Nagłówek wiersza: etykieta, klucz, znaki katalogu, wskaźnik zasięgu. */
function naglowekPola(
  definicja: SettingDefinition,
  identyfikator: string,
  wskaznik: HTMLElement,
): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dk-pole__naglowek';

  const etykieta = document.createElement('label');
  etykieta.className = 'dn-pole-etykieta dk-pole__etykieta';
  etykieta.htmlFor = identyfikator;
  etykieta.textContent = definicja.name;

  const klucz = document.createElement('code');
  klucz.className = 'dk-pole__klucz';
  klucz.textContent = definicja.key;

  element.append(
    etykieta,
    klucz,
    utworzDymekObjasnienia(objasnieniePozycji(definicja)),
    ...znakiKatalogu(definicja),
    wskaznik,
  );
  return element;
}

/**
 * Objaśnienie [?] pozycji katalogu — każdy element konfiguracji ma objaśnienie
 * kontekstowe.
 *
 * Zdania składają się z tego, co katalog o pozycji mówi, i z adnotacji o stanie
 * klucza: część kluczy zapisuje się poprawnie, ale żadna ścieżka rdzenia ich
 * nie czyta. Przemilczenie tego kazałoby Operatorowi wierzyć, że zapisana
 * wartość steruje wykonaniem.
 */
function objasnieniePozycji(definicja: SettingDefinition): string {
  return zlozObjasnienie([
    definicja.description,
    `Klucz: ${definicja.key}.`,
    definicja.required ? 'Pozycja wymagana.' : undefined,
    definicja.restartRequired === true
      ? 'Zmiana obowiązuje po ponownym uruchomieniu.'
      : undefined,
    adnotacjaKlucza(definicja.key),
  ]);
}

/** Plakietki wynikające z metadanych katalogu: wymagana, wymaga restartu. */
function znakiKatalogu(definicja: SettingDefinition): HTMLElement[] {
  const znaki: HTMLElement[] = [];
  if (definicja.required) znaki.push(plakietka('wymagane', 'dn-plakietka'));
  if (definicja.restartRequired === true) {
    znaki.push(plakietka('po ponownym uruchomieniu', 'dn-plakietka dn-plakietka--ostrzezenie'));
  }
  return znaki;
}

function plakietka(tresc: string, klasa: string): HTMLElement {
  const element = document.createElement('span');
  element.className = `${klasa} dk-pole__znak`;
  element.textContent = tresc;
  return element;
}

/** Kontrolka wraz z jednostką drukowaną obok niej. */
function wierszKontrolki(definicja: SettingDefinition, kontrolka: Kontrolka): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dk-pole__kontrolka';
  element.append(kontrolka.element);

  if (definicja.unit !== undefined && definicja.unit !== '') {
    const jednostka = document.createElement('span');
    jednostka.className = 'dk-pole__jednostka';
    jednostka.textContent = definicja.unit;
    element.append(jednostka);
  }

  return element;
}

/** Opis pozycji katalogu oraz ostrzeżenie kontrolki, gdy je zgłosiła. */
function opisyPola(definicja: SettingDefinition, kontrolka: Kontrolka): HTMLElement[] {
  const opisy: HTMLElement[] = [];

  if (definicja.description !== undefined && definicja.description !== '') {
    const opis = document.createElement('p');
    opis.className = 'dn-pole-opis';
    opis.textContent = definicja.description;
    opisy.push(opis);
  }

  if (kontrolka.ostrzezenie !== undefined && kontrolka.ostrzezenie !== '') {
    const ostrzezenie = document.createElement('p');
    ostrzezenie.className = 'dn-pole-opis dk-pole__ostrzezenie';
    ostrzezenie.textContent = kontrolka.ostrzezenie;
    opisy.push(ostrzezenie);
  }

  return opisy;
}
