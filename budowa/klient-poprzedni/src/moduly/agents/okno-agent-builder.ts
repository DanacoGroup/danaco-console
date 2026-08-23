import { opisOdmowy } from '../../komponenty/odmowa';
import { utworzPanelArchiwum, type PanelArchiwum } from './archiwum-ekspertow';
import { naglowek, utworzBiblioteke, type BibliotekaEkspertow } from './biblioteka-ekspertow';
import { utworzZrodloZakresuEksperta, type ZrodloZakresuEksperta } from './zrodlo-zakresu-eksperta';
import { utworzEdytorEksperta, type EdytorEksperta } from './edytor-eksperta';
import { utworzPanelHistorii, type PanelHistorii } from './historia-wersji';
import { utworzLicznikNarzedzi, type LicznikNarzedzi } from './licznik-narzedzi';
import { utworzPanelDoradcy, type PanelDoradcy } from './panel-doradcy';
import { utworzPanelZespolow, type PanelZespolow } from './panel-zespolow';
import {
  utworzPodsumowanieDefinicji,
  type PodsumowanieDefinicji,
} from './podsumowanie-definicji';
import type { StanAgentow } from './stan-agentow';
import { zdanieOCzacieTestowym } from './testowany-agent';
import { utworzZrodloDoradcy } from './zrodlo-doradcy';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';
import type { Kanal } from '../../protokol/kanal';

/**
 * Agent Builder — okno kreator, punkt wejścia modułu Agents.
 *
 * Trzy części: widok „Biblioteka ekspertów”, widok „Edytor” z paskiem zakładek
 * oraz panel „Historia wersji”.
 *
 * Zakładka „Tożsamość” pokazuje edytor tego okna; pozostałe zakładki przenoszą
 * ognisko do właściwego okna modułu — model bazowy mieszka w Model
 * Configuration, umiejętności w Skills Manager i tak dalej. Pasek nie powiela
 * więc żadnego formularza.
 *
 * Panel „Zespoły ekspertów” składa nazwany skład z tego samego wykazu, który
 * okno pokazuje po lewej — dlatego dostaje bibliotekę z `odswiez()` i nie
 * odpytuje rdzenia po raz drugi.
 *
 * Licznik narzędzi stoi pod tożsamością, bo ekspert to tożsamość plus dobór
 * narzędzi. Liczba nie jest tu liczona po raz drugi: `licznik-narzedzi.ts`
 * jest jednym bytem obsadzonym także w Skills Managerze i Connectors
 * Managerze, więc przypisanie kodu w tamtych oknach przestawia tę liczbę
 * w tej samej klatce.
 */
export interface OknoAgentBuilder {
  element: HTMLElement;
  /** Nanosi stan modułu na wszystkie trzy części okna. */
  odswiez(): void;
  /** Odczytuje katalog kategorii tożsamości z rdzenia. */
  wczytajKategorie(): Promise<void>;
  /** Odpina nasłuchy własne okna — na razie wykaz komend panelu archiwum. */
  rozlacz(): void;
}

/** Kody zakładek edytora odpowiadające oknom modułu. */
export const ZAKLADKI_EDYTORA = [
  { kod: 'tozsamosc', nazwa: 'Tożsamość' },
  { kod: 'model-configuration', nazwa: 'Model bazowy' },
  { kod: 'skills-manager', nazwa: 'Umiejętności' },
  { kod: 'connectors-manager', nazwa: 'Konektory' },
  { kod: 'permissions-center', nazwa: 'Uprawnienia' },
] as const;

export interface OpcjeBuildera {
  /** Przenosi ognisko do okna modułu wskazanego kodem zakładki. */
  naZakladke(kod: string): void;
}

export function utworzOknoAgentBuilder(
  stan: StanAgentow,
  zaplecze: ZrodloZaplecza,
  kanal: Kanal,
  opcje: OpcjeBuildera,
): OknoAgentBuilder {
  const biblioteka: BibliotekaEkspertow = utworzBiblioteke(stan);
  const zakres: ZrodloZakresuEksperta = utworzZrodloZakresuEksperta(kanal);

  /**
   * Liczby przypisań dla całej biblioteki jednym wywołaniem.
   *
   * `agent.assignment.list` bez wskazania eksperta oddaje przypisania
   * WSZYSTKICH — dokładnie po to, żeby karta każdego miała licznik po jednym
   * odczycie, a nie po jednym na kartę. Odmowa nie gasi biblioteki: karty
   * zostają bez plakietki, bo zero wpisane z ciszy byłoby nieprawdą.
   */
  async function wczytajPrzypisania(): Promise<void> {
    const wynik = await zakres.przypisania();
    if (!wynik.udany || wynik.wynik === undefined) return;
    const liczby = new Map<string, number>();
    for (const przypisanie of wynik.wynik.assignments) {
      liczby.set(przypisanie.agentId, (liczby.get(przypisanie.agentId) ?? 0) + 1);
    }
    biblioteka.ustawPrzypisania(liczby);
  }
  const edytor: EdytorEksperta = utworzEdytorEksperta(stan, {
    naZapisie: () => void stan.odswiez(),
  });
  // Panel historii woła rdzeń sam, więc bierze kanał, a nie stan modułu:
  // wykaz wersji jest własnością rdzenia i nie da się go wyprowadzić z wykazu
  // biblioteki. Po przywróceniu ekspert wraca tą samą drogą co po każdym innym
  // zapisie — wchłonięciem, żeby okna eksperta przerysowały się w tej klatce.
  const historia: PanelHistorii = utworzPanelHistorii(kanal, {
    naPrzywroceniu: (ekspert) => {
      stan.wchlon(ekspert);
      void stan.odswiez();
    },
  });

  const doradca: PanelDoradcy = utworzPanelDoradcy(
    utworzZrodloDoradcy(kanal),
    zaplecze,
    (tresc) => edytor.dopiszDoInstrukcji(tresc),
  );

  // Archiwizacja i powrót z archiwum ruszają bibliotekę czynną, więc panel
  // dostaje odczyt biblioteki jako oddzwonienie. Bez tego ekspert odłożony
  // wisiałby na liście do najbliższego odczytu — czyli wyglądałby na
  // niezarchiwizowanego.
  const archiwum: PanelArchiwum = utworzPanelArchiwum(kanal, () => void stan.odswiez());

  const zespoly: PanelZespolow = utworzPanelZespolow(kanal);

  const licznik: LicznikNarzedzi = utworzLicznikNarzedzi();

  // Podsumowanie dostaje ten sam przenośnik ogniska co pasek zakładek, bo robi
  // to samo: wiersz „Model: —” prowadzi do okna, które model ustala.
  const podsumowanie: PodsumowanieDefinicji = utworzPodsumowanieDefinicji({
    naZakladke: (kod) => opcje.naZakladke(kod),
  });

  const zakladki = document.createElement('div');
  zakladki.className = 'dn-zakladki da-zakladki';
  zakladki.setAttribute('role', 'tablist');

  const przyciskiZakladek: HTMLButtonElement[] = [];

  /**
   * Znakuje zakładkę czynną.
   *
   * `aria-selected` ustawione raz przy budowie byłoby fałszywym stanem: pasek
   * meldowałby `tozsamosc: true` także po kliknięciu innej zakładki, a czytnik
   * ekranu dostawałby zapewnienie o wyborze, którego Operator nie dokonał.
   * Ognisko przenosi moduł, ale to pasek wie, którą zakładkę przycisnięto,
   * więc znakowanie należy do niego.
   */
  function oznaczZakladke(kod: string): void {
    for (const przycisk of przyciskiZakladek) {
      przycisk.setAttribute('aria-selected', String(przycisk.dataset['zakladka'] === kod));
    }
  }

  for (const zakladka of ZAKLADKI_EDYTORA) {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-zakladka';
    przycisk.setAttribute('role', 'tab');
    przycisk.dataset['zakladka'] = zakladka.kod;
    przycisk.textContent = zakladka.nazwa;
    przycisk.setAttribute('aria-selected', String(zakladka.kod === 'tozsamosc'));
    przycisk.addEventListener('click', () => {
      oznaczZakladke(zakladka.kod);
      opcje.naZakladke(zakladka.kod);
    });
    przyciskiZakladek.push(przycisk);
    zakladki.append(przycisk);
  }

  const prawa = document.createElement('div');
  prawa.className = 'da-builder__prawa';
  prawa.append(
    zakladki,
    podsumowanie.element,
    edytor.element,
    licznik.element,
    doradca.element,
    historia.element,
    archiwum.element,
    zespoly.element,
  );

  const cialo = document.createElement('div');
  cialo.className = 'da-builder__cialo';
  cialo.append(biblioteka.element, prawa);

  // Uprzedzenie stoi przed stratą: wybór eksperta w bibliotece czyści czat
  // testowy tego modułu, więc nota jest widoczna, zanim Operator kliknie.
  const notaCzatu = document.createElement('p');
  notaCzatu.className = 'dn-pole-opis da-granica';
  notaCzatu.dataset['nota'] = 'czat-testowy';
  notaCzatu.textContent = zdanieOCzacieTestowym();

  const element = document.createElement('section');
  element.className = 'da-okno da-okno--kreator';
  element.dataset['okno'] = 'agent-builder';
  element.append(naglowek('Agent Builder'), cialo, notaCzatu);

  return {
    element,

    odswiez() {
      biblioteka.odswiez();
      const czynny = stan.wybrany();
      podsumowanie.ustaw(czynny);
      edytor.ustaw(czynny);
      licznik.ustaw(czynny);
      doradca.ustaw(czynny);
      historia.ustaw(czynny);
      archiwum.ustaw(czynny);
      zespoly.ustawBiblioteke(stan.eksperci());
    },

    async wczytajKategorie() {
      // Rejestr kanałów doradczych jedzie razem z katalogiem kategorii: obie
      // rzeczy są kontekstem edytora, oba odczyty są niezależne i odmowa
      // jednego zostaje w jego panelu.
      await Promise.all([
        doradca.wczytajKanaly(),
        archiwum.wczytajPokrycie(),
        zespoly.wczytaj(),
        wczytajPrzypisania(),
      ]);
      const wynik = await zaplecze.kategorieTozsamosci();
      if (!wynik.udany || wynik.wynik === undefined) {
        edytor.ustawKategorie(
          [],
          opisOdmowy('Odczyt katalogu kategorii tożsamości', wynik.blad?.code, wynik.blad?.message),
        );
        return;
      }
      edytor.ustawKategorie(wynik.wynik.categories, '');
    },

    rozlacz() {
      archiwum.rozlacz();
      zespoly.rozlacz();
    },
  };
}
