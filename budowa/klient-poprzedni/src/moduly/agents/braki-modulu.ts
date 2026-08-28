import { Command } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';

/**
 * Wykaz pozycji modułu Agents bez drogi do rdzenia — zamknięty i policzony
 * wykaz, gdzie każda pozycja niesie komendę kontraktu, która ma ją unieść.
 */
export interface PozycjaBraku {
  /** Nazwa pozycji w języku Operatora — tak brzmi w opracowaniu modułu. */
  etykieta: string;
  /** Komenda kontraktu, która tę pozycję wykona. */
  komenda: string;
  /** Co ta komenda umożliwia — wchodzi w zdanie powodu. */
  czynnosc: string;
  /** Okno modułu, w którym pozycja ma swoje miejsce. */
  okno: string;
}

/**
 * Dwie pozycje modułu — wykaz zamknięty i krótki, obejmujący wyłącznie
 * pozycje z rodziny komend rejestru kanałów modelu.
 */
export const BRAKI_MODULU: readonly PozycjaBraku[] = [
  {
    etykieta: 'Test połączenia kanału modelu',
    komenda: Command.ChannelCheck,
    czynnosc: 'sprawdzenie osiągalności kanału przed zapisem eksperta',
    okno: 'Model Configuration',
  },
  {
    etykieta: 'Stan poświadczenia kanału',
    komenda: Command.ChannelCredentialStatus,
    czynnosc:
      'odczyt stanu poświadczenia kanału — ustawione czy brak i kiedy zmienione, nigdy treść klucza',
    okno: 'Model Configuration',
  },
];

export interface WykazBrakow {
  element: HTMLElement;
  /** Pyta rdzeń o wykaz komend, żeby zdania mówiły o rdzeniu, nie o domyśle. */
  wczytaj(): Promise<void>;
  /** Odpina wykaz pokrycia od wspólnego bytu. */
  zamknij(): void;
}

export function utworzWykazBrakow(kanal: Kanal): WykazBrakow {
  const pokrycie: PokrycieKomend = utworzPokrycieKomend(kanal);

  const tytul = document.createElement('h3');
  tytul.className = 'da-okno__tytul';
  tytul.textContent = 'Pozycje opracowania bez drogi do rdzenia';

  const podsumowanie = document.createElement('p');
  podsumowanie.className = 'dn-pole-opis';

  const lista = document.createElement('ul');
  lista.className = 'da-braki';

  const wiersze = new Map<string, HTMLElement>();

  for (const pozycja of BRAKI_MODULU) {
    const nazwa = document.createElement('span');
    nazwa.className = 'da-braki__nazwa';
    nazwa.textContent = pozycja.etykieta;

    const gdzie = document.createElement('span');
    gdzie.className = 'da-braki__okno';
    gdzie.textContent = pozycja.okno;

    const wiersz = document.createElement('li');
    wiersz.className = 'da-braki__wiersz';
    wiersz.dataset['komenda'] = pozycja.komenda;
    wiersz.append(nazwa, gdzie, pokrycie.przycisk('Wykonaj', pozycja.komenda, pozycja.czynnosc));
    wiersze.set(pozycja.komenda, wiersz);
    lista.append(wiersz);
  }

  /** Zdanie zbiorcze i znakowanie wierszy, czytane z rdzenia po odpowiedzi na powitanie połączenia. */
  function przerysuj(): void {
    let bezUchwytu = 0;
    let nieustalone = 0;
    for (const pozycja of BRAKI_MODULU) {
      const stan = pokrycie.stan(pozycja.komenda);
      wiersze.get(pozycja.komenda)?.setAttribute('data-pokrycie', stan);
      if (stan === 'nieustalone') nieustalone += 1;
      else if (stan !== 'rdzen-ma') bezUchwytu += 1;
    }
    if (nieustalone === BRAKI_MODULU.length) {
      podsumowanie.textContent =
        `Wykaz zamknięty: ${BRAKI_MODULU.length} pozycji opracowania. Pokrycie w rdzeniu — ` +
        'odczyt w toku; liczba pojawi się po odpowiedzi rdzenia.';
      return;
    }
    podsumowanie.textContent =
      `Wykaz zamknięty: ${BRAKI_MODULU.length} pozycji opracowania. Kontrakt niesie komendę ` +
      `dla każdej z nich; rdzeń nie rejestruje jeszcze ${bezUchwytu}. Kontrolka pozycji jest ` +
      'klikalna i po naciśnięciu nazywa stan wraz z tym, po czyjej stronie brak leży.';
  }

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--braki';
  element.dataset['okno'] = 'braki-modulu';
  element.append(tytul, podsumowanie, lista);

  pokrycie.naOdczyt(przerysuj);

  return {
    element,
    wczytaj: () => pokrycie.odczytaj(),
    zamknij: () => pokrycie.zamknij(),
  };
}
