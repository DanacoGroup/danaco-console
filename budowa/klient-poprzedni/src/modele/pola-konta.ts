import {
  poleLogiczne,
  poleTajne,
  poleTekstowe,
  poleWyboru,
  type PoleFormularza,
} from './kontrolki-formularza';
import { NAZWY_RODZAJOW, RODZAJE_KONT } from './rodzaje-kont';

/**
 * Pola formularza konta — komplet kontrolek wraz z ich opisami.
 *
 * Zestaw pól wynika wprost z kontraktu (`AccountAddRequest`,
 * `AccountUpdateRequest`), a nie z katalogu wierszy, więc stoi w kodzie jawnie.
 * Od formularza jest oddzielony, bo formularz odpowiada za rozgałęzienie zapisu
 * na dwie komendy i za to, co wolno nadpisać w wypełnianym polu.
 * Opis poświadczenia należy do pola: kontrakt nie zwraca poświadczenia żadną
 * komendą, więc pole konta, które je ma, i tak wygląda na puste.
 */
export interface PolaKonta {
  nazwa: PoleFormularza<HTMLInputElement>;
  rodzaj: PoleFormularza<HTMLSelectElement>;
  /** Rodzaj konta już założonego — wydruk zamiast kontrolki. */
  znakRodzaju: HTMLElement;
  dostawca: PoleFormularza<HTMLInputElement>;
  zewnetrzny: PoleFormularza<HTMLInputElement>;
  model: PoleFormularza<HTMLInputElement>;
  adres: PoleFormularza<HTMLInputElement>;
  katalog: PoleFormularza<HTMLInputElement>;
  poswiadczenie: PoleFormularza<HTMLInputElement>;
  czynne: PoleFormularza<HTMLInputElement>;
  domyslne: PoleFormularza<HTMLInputElement>;
}

export function utworzPolaKonta(): PolaKonta {
  const znakRodzaju = document.createElement('p');
  znakRodzaju.className = 'dm-formularz__rodzaj';

  return {
    nazwa: poleTekstowe({
      etykieta: 'Nazwa konta',
      podpowiedz: 'nazwa widoczna w wykazie i w kanale modelu',
    }),
    rodzaj: poleWyboru(
      { etykieta: 'Rodzaj konta', opis: 'Rodzaju nie zmienia się po założeniu konta.' },
      RODZAJE_KONT.map((wartosc) => ({ wartosc, etykieta: NAZWY_RODZAJOW[wartosc] })),
    ),
    znakRodzaju,
    dostawca: poleTekstowe({
      etykieta: 'Dostawca',
      podpowiedz: 'wartość danych, nie typ kodu',
    }),
    zewnetrzny: poleTekstowe({
      etykieta: 'Identyfikator konta u dostawcy',
      podpowiedz: 'pole puste znaczy brak identyfikatora',
    }),
    model: poleTekstowe({
      etykieta: 'Model domyślny konta',
      podpowiedz: 'identyfikator modelu używany, gdy okno nie wskaże innego',
    }),
    adres: poleTekstowe({
      etykieta: 'Adres punktu końcowego',
      podpowiedz: 'adres bazowy dostawcy',
    }),
    katalog: poleTekstowe({
      etykieta: 'Katalog konfiguracji programu code CLI',
      podpowiedz: 'katalog, w którym program trzyma własną konfigurację',
    }),
    poswiadczenie: poleTajne({
      etykieta: 'Poświadczenie konta',
      podpowiedz: 'wpisz, aby zapisać nowe',
      opis: OPIS_POSWIADCZENIA,
    }),
    czynne: poleLogiczne({
      etykieta: 'Konto czynne',
      opis: 'Konto nieczynne zostaje w rejestrze i nie wchodzi do rotacji.',
    }),
    domyslne: poleLogiczne({
      etykieta: 'Uczyń kontem domyślnym swojego rodzaju',
      opis: 'Domyślne jest dokładnie jedno konto na rodzaj; poprzednie traci oznaczenie.',
    }),
  };
}

/** Kolejność pól w formularzu — jedno miejsce, w którym jest zapisana. */
export function elementyPol(pola: PolaKonta): HTMLElement[] {
  return [
    pola.nazwa.element,
    pola.rodzaj.element,
    pola.znakRodzaju,
    pola.dostawca.element,
    pola.zewnetrzny.element,
    pola.model.element,
    pola.adres.element,
    pola.katalog.element,
    pola.poswiadczenie.element,
    pola.czynne.element,
    pola.domyslne.element,
  ];
}

const OPIS_POSWIADCZENIA =
  'Pole wyłącznie do wprowadzania. Kontrakt nie zwraca poświadczenia żadną komendą — wykaz mówi jedynie, czy jest zapisane.';
