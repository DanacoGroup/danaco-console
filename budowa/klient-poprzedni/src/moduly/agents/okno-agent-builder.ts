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
 * Agent Builder to okno kreatora i punkt wejścia modułu Agents, złożone z biblioteki
 * ekspertów, edytora z paskiem zakładek i panelu historii wersji.
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

/**
 * Kody zakładek paska edytora wraz z nazwami wyświetlanymi w interfejsie, odpowiadające
 * oknom modułu Agents, do których zakładka przenosi ognisko.
 */
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

  /** Odczytuje liczby przypisań dla całej biblioteki jednym wywołaniem agent.assignment.list. */
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
  // Panel historii woła rdzeń przez kanał, bo wykaz wersji jest własnością rdzenia.
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

  // Panel archiwum dostaje odczyt biblioteki jako oddzwonienie po archiwizacji i przywróceniu.
  const archiwum: PanelArchiwum = utworzPanelArchiwum(kanal, () => void stan.odswiez());

  const zespoly: PanelZespolow = utworzPanelZespolow(kanal);

  const licznik: LicznikNarzedzi = utworzLicznikNarzedzi();

  // Podsumowanie dostaje ten sam przenośnik ogniska co pasek zakładek.
  const podsumowanie: PodsumowanieDefinicji = utworzPodsumowanieDefinicji({
    naZakladke: (kod) => opcje.naZakladke(kod),
  });

  const zakladki = document.createElement('div');
  zakladki.className = 'dn-zakladki da-zakladki';
  zakladki.setAttribute('role', 'tablist');

  const przyciskiZakladek: HTMLButtonElement[] = [];

  /** Znakuje zakładkę czynną atrybutem aria-selected dla czytnika ekranu. */
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

  // Nota informuje, że wybór eksperta w bibliotece czyści czat testowy modułu.
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
      // Rejestr kanałów doradczych i katalog kategorii to niezależne odczyty kontekstu edytora.
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
