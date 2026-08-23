import { Command } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';

/**
 * Wykaz pozycji modułu Agents bez drogi do rdzenia — zamknięty i policzony.
 *
 * Wykaz stoi w module, a nie rozsypany po oknach, bo odpowiada na pytanie
 * zadawane raz i o całość: czego ten moduł nie potrafi i po czyjej stronie brak
 * leży. Pozycje mają też własne kontrolki tam, gdzie Operator ich szuka — przy
 * wierszu umiejętności, przy wierszu konektora, przy grupie zakresu — ale tamte
 * mówią o jednej rzeczy naraz i nie dają liczby.
 *
 * Nazwa komendy pochodzi ze stałej generatu, nie z napisu. Napis przetrwałby
 * zmianę nazwy w kontrakcie i zostawiłby w oknie zdanie o pozycji, której już
 * nie ma pod tą nazwą; stała przerywa kompilację i każe poprawkę wykonać.
 *
 * Zdanie o każdej pozycji bierze się z odpowiedzi rdzenia, nie ze stałej
 * wpisanej w moduł. Byt pokrycia rozstrzyga pięć stanów i odróżnia „kontrakt
 * tego nie ma” od „kontrakt ma, rdzeń nie ma uchwytu” — więc wykaz sam
 * przeszedł ze stanu pierwszego w drugi w chwili scalenia definicji, bez ani
 * jednej poprawki w tym pliku. O to w tym bycie chodziło.
 *
 * Pozycja nie znika i nie jest wygaszona. Kontrolka zostaje klikalna i po
 * naciśnięciu nazywa stan pozycji — okno bez niej wyglądałoby na skończone,
 * a brak przestałby być widoczny.
 */

/** Jedna pozycja opracowania wraz z komendą, która ma ją unieść. */
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
 * Dwie pozycje — wykaz zamknięty i krótki.
 *
 * Wykaz liczył czternaście pozycji, dopóki rodzina `agent.*` nie miała
 * uchwytów. Dwanaście z nich zeszło stąd nie dlatego, że przestały być
 * potrzebne, lecz dlatego, że mają już drogę z okna do rdzenia: umiejętności
 * i konektory w swoich zarządcach, podgląd wersji w panelu historii, licznik
 * przypisań na karcie eksperta, a cztery grupy zakresu w Permissions Center.
 *
 * Zostają dwie i obie należą do rodziny `channel.*`, czyli do rejestru kanałów
 * modelu — nie do modułu Agents. Moduł ich potrzebuje (Model Configuration
 * chce sprawdzić kanał przed zapisem eksperta), ale nie jest ich właścicielem
 * i dobudowanie ich stąd byłoby wejściem w cudzy obszar.
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

  /**
   * Zdanie zbiorcze i znakowanie wierszy — czytane z rdzenia po każdym odczycie.
   *
   * Liczby nie ma przed odpowiedzią rdzenia i nie jest to niedopatrzenie:
   * policzenie braków z ciszy byłoby orzeczeniem, którego nikt nie wydał.
   * Dopiero powitanie mówi, ile z tych komend rdzeń faktycznie rejestruje.
   */
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
