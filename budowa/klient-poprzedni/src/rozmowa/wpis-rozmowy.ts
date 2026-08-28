import type { MetadaneKonta } from './metadane-konta';
import { RodzajNadawcy } from './nadawca';
import type { PodsumowanieTury } from './podsumowanie-tury';
import type { Prowenancja } from './prowenancja';

/**
 * Stan wpisu rozmowy w cyklu jednej tury tego okna komunikacji, prowadzonej wprost z
 * rdzeniem tej aplikacji.
 */
export type StanWpisu =
  /** Wysłano komendę, rdzeń jeszcze nie nadał ani jednego fragmentu. */
  | 'wysylanie'
  /** Fragmenty przychodzą; treść narasta. */
  | 'strumien'
  /** Strumień domknięty znacznikiem końca. */
  | 'zakonczony'
  // Tura domknięta bez żadnej treści: ani znaku odpowiedzi, rozumowania, narzędzia czy
  // błędu.
  | 'pusty'
  /** Turę przerwał Operator komendą zatrzymania. */
  | 'przerwany'
  /** Turę zamknął błąd kanału; sesja i okno pozostają czynne. */
  | 'bledny';

/**
 * Jedno wywołanie narzędzia zgłoszone przez model w trakcie tury, wraz z jego wynikiem po
 * zakończeniu.
 */
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

/**
 * Błąd odebrany fragmentem strumienia odpowiedzi rdzenia dla tej konkretnej tury tego okna
 * bieżącej rozmowy.
 */
export interface BladWpisu {
  /** Kod z katalogu kontraktu, jeżeli rdzeń go podał. */
  kod: string;
  /** Opis dla Operatora. */
  tresc: string;
  /** Czy ponowienie żądania ma sens. */
  ponawialny: boolean;
}

/**
 * Wpis rozmowy — jedna pozycja historii okna, odpowiadająca całej turze modelu, nie
 * pojedynczemu fragmentowi.
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

/**
 * Zakłada zupełnie pusty wpis wskazanego nadawcy na samym początku każdej nowej tury tej
 * bieżącej rozmowy.
 */
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

/**
 * Wpis Operatora — powstaje w chwili wysłania jego wypowiedzi, jeszcze przed odpowiedzią
 * samego rdzenia.
 */
export function wpisOperatora(id: string, tresc: string, persona: string): WpisRozmowy {
  const wpis = nowyWpis(id, RodzajNadawcy.Uzytkownik, persona, 'zakonczony');
  wpis.tresc = tresc;
  wpis.domkniety = true;
  return wpis;
}

/**
 * Wpis warstwy automatycznej — komunikat rdzenia, transportu albo kolejki zdarzeń tego
 * całego systemu.
 */
export function wpisAutomatyzacji(id: string, tresc: string): WpisRozmowy {
  const wpis = nowyWpis(id, RodzajNadawcy.Automatyzacja, 'Danaco Console', 'zakonczony');
  wpis.tresc = tresc;
  wpis.domkniety = true;
  return wpis;
}

/**
 * Czy dana tura domknęła się, nie przynosząc ze sobą ani jednej treści widocznej na ekranie
 * po jej całkowitym zakończeniu.
 */
export function czyPustaTura(wpis: WpisRozmowy): boolean {
  return (
    wpis.tresc.trim().length === 0 &&
    wpis.rozumowanie.trim().length === 0 &&
    wpis.narzedzia.length === 0 &&
    wpis.bledy.length === 0
  );
}

/**
 * Godzina wpisu tej rozmowy zapisana w postaci gotowej do wyświetlenia w interfejsie tego
 * okna komunikacji.
 */
export function godzinaWpisu(wpis: WpisRozmowy): string {
  return new Date(wpis.znacznikCzasu).toLocaleTimeString('pl-PL', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}
