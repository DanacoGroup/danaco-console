import { RodzajNadawcy } from './nadawca';
import { nowyWpis, type WpisRozmowy } from './wpis-rozmowy';

/**
 * Historia tur okna — zbiór wpisów rozmowy wraz z ich tożsamością.
 *
 * Komenda `message.send` zwraca wiadomość użytkownika, a identyfikator
 * odpowiedzi znany jest dopiero z pierwszego fragmentu strumienia. Dlatego wpis
 * odpowiedzi powstaje od razu po wysłaniu — jako wpis oczekujący — a pierwszy
 * fragment go przejmuje zamiast zakładać drugi. Bez tego okno stałoby puste do
 * nadejścia pierwszego znaku, a potem pokazało dwa wpisy zamiast jednego.
 */
export interface HistoriaTur {
  /** Wpisy w kolejności powstania. */
  wpisy(): WpisRozmowy[];
  /** Dokłada gotowy wpis na koniec historii. */
  dodaj(wpis: WpisRozmowy): WpisRozmowy;
  /** Zakłada wpis odpowiedzi bez znanego jeszcze identyfikatora wiadomości. */
  oczekujaca(nadawca: RodzajNadawcy, persona: string): WpisRozmowy;
  /** Wpis wiadomości: odszukany, przejęty z oczekującego albo nowy. */
  dlaWiadomosci(
    idWiadomosci: string,
    nadawca: RodzajNadawcy,
    persona: string,
  ): WpisRozmowy;
  /** Wpis o wskazanym identyfikatorze wiadomości albo `undefined`. */
  znajdz(idWiadomosci: string): WpisRozmowy | undefined;
  /**
   * Wiąże wypowiedź wpisaną miejscowo z wiadomością rdzenia.
   *
   * Wypowiedź trafia do historii natychmiast, przed komendą, więc nie ma
   * jeszcze identyfikatora wiadomości: `message.send` oddaje go dopiero
   * w odpowiedzi. Rdzeń rozgłasza tę samą wiadomość zdarzeniem
   * `message.changed` do wszystkich połączeń konta — także do tego, które ją
   * wysłało. Ponieważ okno pokazuje wypowiedzi roli `user`, to samo zdanie
   * weszłoby do wątku dwa razy: raz jako echo miejscowe, raz ze zdarzenia.
   *
   * Wiązanie idzie po treści, nie po identyfikatorze, bo identyfikatora w tej
   * chwili nie ma po żadnej stronie: dopasowuje pierwszą niezwiązaną wypowiedź
   * o tej samej treści i nadaje jej identyfikator z rdzenia. Od tej chwili wpis
   * jest w mapie, więc każde kolejne zdarzenie o tej wiadomości trafia w niego.
   *
   * @returns wpis związany albo `undefined`, gdy wypowiedzi o tej treści
   *   w historii nie ma — czyli gdy wiadomość naprawdę przyszła skądinąd.
   */
  zwiazWypowiedz(tresc: string, idWiadomosci: string): WpisRozmowy | undefined;
  /** Wpis odpowiedzi biegnącej w tej chwili albo `undefined`. */
  biezaca(): WpisRozmowy | undefined;
  /**
   * Zdejmuje wszystkie wpisy — rozmowa ulotna zaczyna nowy kontekst roboczy.
   *
   * Licznik kluczy nie wraca do zera. Klucz jest tożsamością pozycji w widoku;
   * gdyby po wyczyszczeniu zaczął się od nowa, pierwszy wpis nowej rozmowy
   * trafiłby w węzeł pozostały po wpisie rozmowy poprzedniej i zamiast nowej
   * pozycji stanęłaby podmieniona stara.
   */
  wyczysc(): void;
}

export function utworzHistorieTur(): HistoriaTur {
  const kolejnosc: WpisRozmowy[] = [];
  const poWiadomosci = new Map<string, WpisRozmowy>();
  let oczekujacy: WpisRozmowy | null = null;
  let licznik = 0;

  /** Kolejny klucz wpisu; rośnie w obrębie jednego okna. */
  function nastepnyKlucz(): string {
    licznik += 1;
    return `wpis-${licznik}`;
  }

  function dodaj(wpis: WpisRozmowy): WpisRozmowy {
    kolejnosc.push(wpis);
    if (wpis.idWiadomosci.length > 0) poWiadomosci.set(wpis.idWiadomosci, wpis);
    return wpis;
  }

  function oczekujaca(nadawca: RodzajNadawcy, persona: string): WpisRozmowy {
    const wpis = nowyWpis(nastepnyKlucz(), nadawca, persona, 'wysylanie');
    oczekujacy = wpis;
    return dodaj(wpis);
  }

  function dlaWiadomosci(
    idWiadomosci: string,
    nadawca: RodzajNadawcy,
    persona: string,
  ): WpisRozmowy {
    const znany = poWiadomosci.get(idWiadomosci);
    if (znany !== undefined) return znany;

    if (oczekujacy !== null) {
      const przejety = oczekujacy;
      oczekujacy = null;
      przejety.idWiadomosci = idWiadomosci;
      poWiadomosci.set(idWiadomosci, przejety);
      return przejety;
    }

    const wpis = nowyWpis(nastepnyKlucz(), nadawca, persona, 'strumien');
    wpis.idWiadomosci = idWiadomosci;
    return dodaj(wpis);
  }

  function wyczysc(): void {
    kolejnosc.length = 0;
    poWiadomosci.clear();
    oczekujacy = null;
  }

  return {
    wpisy: () => [...kolejnosc],
    dodaj,
    oczekujaca,
    dlaWiadomosci,
    znajdz: (idWiadomosci) => poWiadomosci.get(idWiadomosci),

    zwiazWypowiedz(tresc, idWiadomosci) {
      // Szukamy od końca: przy dwóch identycznych zdaniach wysłanych pod rząd
      // wiąże się to, które jeszcze na identyfikator czeka, a nie to sprzed
      // dziesięciu minut, które swój identyfikator dawno dostało.
      const szukana = tresc.trim();
      for (let i = kolejnosc.length - 1; i >= 0; i -= 1) {
        const wpis = kolejnosc[i];
        if (wpis === undefined) continue;
        if (wpis.idWiadomosci.length > 0) continue;
        if (wpis.nadawca !== RodzajNadawcy.Uzytkownik) continue;
        if (wpis.tresc.trim() !== szukana) continue;
        wpis.idWiadomosci = idWiadomosci;
        poWiadomosci.set(idWiadomosci, wpis);
        return wpis;
      }
      return undefined;
    },
    biezaca: () =>
      oczekujacy ??
      [...kolejnosc].reverse().find((wpis) => wpis.stan === 'strumien'),
    wyczysc,
  };
}

/** Klucz wpisu miejscowego — komunikatu warstwy automatycznej albo Operatora. */
export function kluczMiejscowy(przedrostek: string): string {
  return `${przedrostek}-${Date.now().toString(36)}-${Math.floor(Math.random() * 1e6).toString(36)}`;
}
