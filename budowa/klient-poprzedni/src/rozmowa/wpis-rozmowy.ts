import type { MetadaneKonta } from './metadane-konta';
import { RodzajNadawcy } from './nadawca';
import type { PodsumowanieTury } from './podsumowanie-tury';
import type { Prowenancja } from './prowenancja';

/** Stan wpisu w cyklu tury. */
export type StanWpisu =
  /** Wysłano komendę, rdzeń jeszcze nie nadał ani jednego fragmentu. */
  | 'wysylanie'
  /** Fragmenty przychodzą; treść narasta. */
  | 'strumien'
  /** Strumień domknięty znacznikiem końca. */
  | 'zakonczony'
  /**
   * Tura domknięta, a nie przyniosła niczego: ani znaku odpowiedzi, ani toku
   * rozumowania, ani wywołania narzędzia, ani błędu. Stan osobny, bo etykieta
   * „zakończona" przy pustym wpisie nie mówi, czy model milczał, czy widok
   * czegoś nie narysował.
   */
  | 'pusty'
  /** Turę przerwał Operator komendą zatrzymania. */
  | 'przerwany'
  /** Turę zamknął błąd kanału; sesja i okno pozostają czynne. */
  | 'bledny';

/** Jedno wywołanie narzędzia wraz z jego wynikiem. */
export interface WywolanieNarzedzia {
  /** Identyfikator wywołania nadany przez model. */
  id: string;
  /** Nazwa narzędzia. */
  nazwa: string;
  /** Wejście przekazane narzędziu, w zapisie czytelnym. */
  wejscie: string;
  /** Wynik zwrócony modelowi; pusty, dopóki nie nadszedł. */
  wynik: string;
  /** Czy narzędzie zgłosiło błąd. */
  bledne: boolean;
}

/** Błąd odebrany fragmentem strumienia. */
export interface BladWpisu {
  /** Kod z katalogu kontraktu, jeżeli rdzeń go podał. */
  kod: string;
  /** Opis dla Operatora. */
  tresc: string;
  /** Czy ponowienie żądania ma sens. */
  ponawialny: boolean;
}

/**
 * Wpis rozmowy — jedna pozycja historii okna komunikacji.
 *
 * Wpis modelu odpowiada jednej turze, nie jednemu fragmentowi. Wszystko, co
 * turę opisuje — prowenancja wywołania, tok rozumowania, wywołania narzędzi,
 * konto, błędy, podsumowanie — wisi przy tym samym wpisie, zamiast rozsypywać
 * się po historii na osobne pozycje. Dzięki temu „co poszło do modelu" stoi
 * obok tego, co model odpowiedział.
 */
export interface WpisRozmowy {
  /** Klucz wpisu w historii. */
  id: string;
  /** Wiadomość kontraktu, której wpis dotyczy; pusty dla wpisów miejscowych. */
  idWiadomosci: string;
  /** Rodzaj nadawcy — jeden z dziewięciu. */
  nadawca: RodzajNadawcy;
  /** Tożsamość mówiącego wewnątrz rodzaju: nazwa modelu, agenta albo kanału. */
  persona: string;
  /** Treść wypowiedzi; przy strumieniu narasta. */
  tresc: string;
  /** Tok rozumowania modelu; podgląd pracy na żywo, zwijany. */
  rozumowanie: string;
  /** Prowenancja wywołania; nadchodzi przed tekstem. */
  prowenancja: Prowenancja | null;
  /** Konto użyte przez kanał; wypełnione także przy rotacji. */
  konto: MetadaneKonta | null;
  /** Podsumowanie tury; nadchodzi z fragmentem kończącym. */
  podsumowanie: PodsumowanieTury | null;
  /** Wywołania narzędzi w kolejności nadania. */
  narzedzia: WywolanieNarzedzia[];
  /** Błędy odebrane w trakcie tury. */
  bledy: BladWpisu[];
  /** Stan wpisu. */
  stan: StanWpisu;
  /** Liczba odebranych fragmentów. */
  fragmenty: number;
  /** Największy numer fragmentu odebrany dla tej tury (pole `seq` koperty). */
  ostatniNumer: number;
  /** Czy odebrano fragment ze znacznikiem końca (pole `done` koperty). */
  domkniety: boolean;
  /** Czas założenia wpisu w milisekundach epoki. */
  znacznikCzasu: number;
}

/** Zakłada pusty wpis wskazanego nadawcy. */
export function nowyWpis(
  id: string,
  nadawca: RodzajNadawcy,
  persona: string,
  stan: StanWpisu,
): WpisRozmowy {
  return {
    id,
    idWiadomosci: '',
    nadawca,
    persona,
    tresc: '',
    rozumowanie: '',
    prowenancja: null,
    konto: null,
    podsumowanie: null,
    narzedzia: [],
    bledy: [],
    stan,
    fragmenty: 0,
    ostatniNumer: 0,
    domkniety: false,
    znacznikCzasu: Date.now(),
  };
}

/** Wpis Operatora — powstaje w chwili wysłania, przed odpowiedzią rdzenia. */
export function wpisOperatora(id: string, tresc: string, persona: string): WpisRozmowy {
  const wpis = nowyWpis(id, RodzajNadawcy.Uzytkownik, persona, 'zakonczony');
  wpis.tresc = tresc;
  wpis.domkniety = true;
  return wpis;
}

/** Wpis warstwy automatycznej — komunikat rdzenia, transportu albo kolejki. */
export function wpisAutomatyzacji(id: string, tresc: string): WpisRozmowy {
  const wpis = nowyWpis(id, RodzajNadawcy.Automatyzacja, 'Danaco Console', 'zakonczony');
  wpis.tresc = tresc;
  wpis.domkniety = true;
  return wpis;
}

/**
 * Czy tura domknęła się, nie przynosząc ani jednej treści.
 *
 * Pytanie dotyczy wszystkich warstw wpisu, nie samego tekstu: tura, która
 * oddała sam tok rozumowania albo samo wywołanie narzędzia, coś przyniosła.
 * Pusta jest dopiero taka, po której na ekranie nie zostaje nic prócz nagłówka.
 */
export function czyPustaTura(wpis: WpisRozmowy): boolean {
  return (
    wpis.tresc.trim().length === 0 &&
    wpis.rozumowanie.trim().length === 0 &&
    wpis.narzedzia.length === 0 &&
    wpis.bledy.length === 0
  );
}

/** Godzina wpisu w postaci prezentacyjnej. */
export function godzinaWpisu(wpis: WpisRozmowy): string {
  return new Date(wpis.znacznikCzasu).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
