import { RodzajNadawcy } from './nadawca';
import { nowyWpis, type WpisRozmowy } from './wpis-rozmowy';

/** Interfejs opisuje historię tur okna: zbiór wpisów rozmowy wraz z ich tożsamością i wpisem odpowiedzi zakładanym przed nadejściem identyfikatora wiadomości. */
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
  /** Wiąże wypowiedź wpisaną miejscowo z wiadomością nadaną przez rdzeń po ustaleniu identyfikatora. */
  zwiazWypowiedz(tresc: string, idWiadomosci: string): WpisRozmowy | undefined;
  /** Wpis odpowiedzi biegnącej w tej chwili albo `undefined`. */
  biezaca(): WpisRozmowy | undefined;
  /** Zdejmuje wszystkie wpisy z historii, rozpoczynając nowy kontekst roboczy rozmowy. */
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
      // Szukamy od końca, aby związać najnowszą niezwiązaną wypowiedź o tej samej treści.
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

/** Funkcja generuje klucz wpisu miejscowego dla komunikatu warstwy automatycznej albo wypowiedzi operatora, łącząc znacznik czasu z liczbą losową. */
export function kluczMiejscowy(przedrostek: string): string {
  return `${przedrostek}-${Date.now().toString(36)}-${Math.floor(Math.random() * 1e6).toString(36)}`;
}
