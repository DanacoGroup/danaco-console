import type { Agent } from '../../../../shared/contract';
import {
  KATALOG_NARZEDZI,
  rozpoznajKody,
  type PozycjaNarzedzia,
  type RozpoznanieKodow,
} from './katalog-narzedzi';

/**
 * Wynik pomiaru narzędzi eksperta: ile narzędzi ekspert załaduje przy wywołaniu i skąd
 * ta liczba się bierze, wraz z rozpoznaniem kodów i wtyczek.
 */
export interface PomiarNarzedzi {
  /** Kody z `skillIds` i `connectorIds`, bez pustych i bez powtórzeń. */
  kody: readonly string[];
  /** Rozpoznanie kodów wobec katalogu kontraktu. */
  rozpoznanie: RozpoznanieKodow;
  /** Narzędzia, które model dostanie w tej turze. */
  narzedzia: readonly PozycjaNarzedzia[];
  /** Liczba narzędzi ładowanych — to ona stoi na liczniku. */
  ladowane: number;
  /** Ile narzędzi niesie wykaz kontraktu, czyli od czego się zawęża. */
  wszystkie: number;
  /** Czy wykaz jest podzbiorem wskazanym przez eksperta. */
  zawezony: boolean;
  /** Zdanie tłumaczące brak zawężenia; puste, gdy zawężenie weszło. */
  powod: string;
  /** Wtyczki eksperta — inna droga niż wykaz narzędzi, liczone osobno. */
  wtyczki: number;
}

/**
 * Pomiar zwracany, gdy żaden ekspert nie jest wybrany — okno ma co pokazać na liczniku,
 * zanim użytkownik dokona wyboru eksperta.
 */
const BEZ_EKSPERTA: PomiarNarzedzi = {
  kody: [],
  rozpoznanie: { nazwy: [], grupy: [], nierozpoznane: [], narzedzia: [] },
  narzedzia: [],
  ladowane: 0,
  wszystkie: KATALOG_NARZEDZI.length,
  zawezony: false,
  powod: 'Żaden ekspert nie jest wybrany — nie ma czyjego wyposażenia policzyć.',
  wtyczki: 0,
};

/**
 * Liczy narzędzia eksperta zgodnie z doborem serwera narzędzi na podstawie kodów
 * w polach skillIds i connectorIds.
 */
export function policz(ekspert: Agent | null): PomiarNarzedzi {
  if (ekspert === null) return BEZ_EKSPERTA;

  const surowe = [...(ekspert.skillIds ?? []), ...(ekspert.connectorIds ?? [])];
  const kody = [...new Set(surowe.filter((kod) => kod !== ''))];
  const rozpoznanie = rozpoznajKody(kody);
  const wszystkie = KATALOG_NARZEDZI.length;
  const wtyczki = (ekspert.pluginIds ?? []).length;

  if (kody.length === 0) {
    return {
      kody,
      rozpoznanie,
      narzedzia: KATALOG_NARZEDZI,
      ladowane: wszystkie,
      wszystkie,
      zawezony: false,
      powod:
        `Ekspert „${ekspert.name}" nie wskazał ani jednego kodu w polach skillIds ` +
        'i connectorIds — nie ma czym zawęzić wykazu, więc model dostanie go w całości.',
      wtyczki,
    };
  }

  if (rozpoznanie.narzedzia.length === 0) {
    return {
      kody,
      rozpoznanie,
      narzedzia: KATALOG_NARZEDZI,
      ladowane: wszystkie,
      wszystkie,
      zawezony: false,
      powod:
        `Ekspert „${ekspert.name}" wskazał ${kody.length} ${odmianaKodow(kody.length)} ` +
        'i żaden nie nazywa narzędzia ani grupy — zawężenia nie da się wyprowadzić, ' +
        `więc model dostanie wykaz w całości. Kody nierozpoznane: ${rozpoznanie.nierozpoznane.join(', ')}.`,
      wtyczki,
    };
  }

  return {
    kody,
    rozpoznanie,
    narzedzia: rozpoznanie.narzedzia,
    ladowane: rozpoznanie.narzedzia.length,
    wszystkie,
    zawezony: true,
    powod: '',
    wtyczki,
  };
}

/**
 * Dobiera odmianę rzeczownika „kod" zależnie od liczby, tak aby zdanie z liczebnikiem
 * brzmiało poprawnie w języku polskim.
 */
function odmianaKodow(ile: number): string {
  if (ile === 1) return 'kod';
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  if (dziesiatki >= 12 && dziesiatki <= 14) return 'kodów';
  return jednosci >= 2 && jednosci <= 4 ? 'kody' : 'kodów';
}

/**
 * Dobiera odmianę rzeczownika „narzędzie" zależnie od liczby, tak aby zdanie
 * z liczebnikiem brzmiało poprawnie w języku polskim.
 */
function odmianaNarzedzi(ile: number): string {
  if (ile === 1) return 'narzędzie';
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  if (dziesiatki >= 12 && dziesiatki <= 14) return 'narzędzi';
  return jednosci >= 2 && jednosci <= 4 ? 'narzędzia' : 'narzędzi';
}

export interface LicznikNarzedzi {
  /** Element montowany w nagłówku okna albo pod tożsamością eksperta. */
  element: HTMLElement;
  /** Nanosi pomiar eksperta; `null` znaczy „żaden nie jest wybrany". */
  ustaw(ekspert: Agent | null): void;
  /** Ostatni pomiar — dla okna, które chce z niego coś dopowiedzieć. */
  pomiar(): PomiarNarzedzi;
}

/**
 * Tworzy widok licznika narzędzi z liczbą, zdaniem opisowym i polem uwagi, ogłaszanym
 * elementom pomocniczym przez atrybut aria-live.
 */
export function utworzLicznikNarzedzi(): LicznikNarzedzi {
  const element = document.createElement('div');
  element.className = 'da-licznik';
  element.dataset['licznik'] = 'narzedzia';
  element.setAttribute('role', 'status');
  element.setAttribute('aria-live', 'polite');

  const liczba = document.createElement('span');
  liczba.className = 'da-licznik__liczba';

  const zdanie = document.createElement('span');
  zdanie.className = 'da-licznik__zdanie';

  const uwaga = document.createElement('p');
  uwaga.className = 'da-licznik__uwaga';

  element.append(liczba, zdanie, uwaga);

  let ostatni: PomiarNarzedzi = BEZ_EKSPERTA;

  function ustaw(ekspert: Agent | null): void {
    const wynik = policz(ekspert);
    ostatni = wynik;

    if (ekspert === null) {
      element.dataset['stan'] = 'bez-eksperta';
      liczba.textContent = '—';
      zdanie.textContent = `narzędzi z ${wynik.wszystkie} w wykazie kontraktu`;
      uwaga.textContent = wynik.powod;
      return;
    }

    element.dataset['stan'] = wynik.zawezony ? 'zawezony' : 'pelny';
    liczba.textContent = String(wynik.ladowane);
    zdanie.textContent =
      `${odmianaNarzedzi(wynik.ladowane)} z ${wynik.wszystkie} — tyle ekspert ` +
      `„${ekspert.name}" załaduje przy wywołaniu`;
    uwaga.textContent = zdanieUwagi(wynik);
  }

  ustaw(null);

  return { element, ustaw, pomiar: () => ostatni };
}

/**
 * Buduje zdanie uwagi pod licznikiem, opisujące powód braku zawężenia, kody
 * nierozpoznane i liczbę wtyczek eksperta.
 */
function zdanieUwagi(wynik: PomiarNarzedzi): string {
  const czesci: string[] = [];
  if (!wynik.zawezony) czesci.push(wynik.powod);
  else if (wynik.rozpoznanie.nierozpoznane.length > 0) {
    czesci.push(
      `Kody bez pokrycia w kontrakcie, pominięte przy doborze: ` +
        `${wynik.rozpoznanie.nierozpoznane.join(', ')}.`,
    );
  }
  if (wynik.wtyczki > 0) {
    czesci.push(
      `Wtyczek: ${wynik.wtyczki} — jadą katalogiem rozszerzeń powłoki, nie wykazem ` +
        'narzędzi modelu, więc nie wchodzą do tej liczby.',
    );
  }
  czesci.push(
    'Wartości granicznej Właściciel nie podał — licznik pokazuje liczbę i nie ' +
      'orzeka, czy jest bezpieczna.',
  );
  return czesci.join(' ');
}
