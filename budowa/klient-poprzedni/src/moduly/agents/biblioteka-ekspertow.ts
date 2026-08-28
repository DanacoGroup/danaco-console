import { AgentVisibility, type Agent } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, przycisk } from '../../modele/kontrolki-formularza';
import { utworzKarteEksperta } from './karta-eksperta';
import { czyZastepuje } from './warstwy-promptu';
import type { StanAgentow } from './stan-agentow';
import { utworzStanOkna, type StanOkna } from './stan-okna';

/**
 * Widok „Biblioteka ekspertów” Agent Buildera obsługuje przegląd, utworzenie,
 * duplikowanie oraz usunięcie eksperta, z zawężeniem wykazu po wprowadzonych
 * kryteriach.
 */
export interface BibliotekaEkspertow {
  element: HTMLElement;
  /** Nanosi stan modułu na wykaz. */
  odswiez(): void;
  /** Nanosi liczby przypisań odczytane komendą przypisań eksperta; mapa nieustawiona różni się od pustej. */
  ustawPrzypisania(liczby: ReadonlyMap<string, number>): void;
}

export function utworzBiblioteke(stan: StanAgentow): BibliotekaEkspertow {
  /** Liczby przypisań po kodzie eksperta; pusta znaczy „jeszcze nie pytano”. */
  let przypisania: ReadonlyMap<string, number> | null = null;
  const okno: StanOkna = utworzStanOkna();

  const szukanie = document.createElement('input');
  szukanie.type = 'search';
  szukanie.className = 'dn-pole-kontrolka da-biblioteka__szukanie';
  szukanie.placeholder = 'Szukaj w bibliotece';
  szukanie.setAttribute('aria-label', 'Fraza wyszukiwania w bibliotece ekspertów');

  const nowy = przycisk('Nowy ekspert', 'dn-btn dn-btn--sm dn-btn--atrament');
  const odczyt = przycisk('Odczytaj bibliotekę', 'dn-btn dn-btn--sm dn-btn--zarys');

  const filtrWidocznosci = poleWyboru({ etykieta: 'Widoczność' }, [
    { wartosc: '', etykieta: 'każdy zasięg' },
    { wartosc: AgentVisibility.Global, etykieta: 'globalny' },
    { wartosc: AgentVisibility.Project, etykieta: 'projektowy' },
  ]);

  const filtrStanu = poleWyboru({ etykieta: 'Stan' }, [
    { wartosc: '', etykieta: 'każdy stan' },
    { wartosc: 'czynny', etykieta: 'czynny' },
    { wartosc: 'wylaczony', etykieta: 'wyłączony' },
  ]);

  const pasek = document.createElement('header');
  pasek.className = 'da-panel__pasek';
  pasek.append(szukanie, filtrWidocznosci.element, filtrStanu.element, nowy, odczyt);

  const lista = document.createElement('ul');
  lista.className = 'da-biblioteka__lista';
  okno.tresc.append(lista);

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'da-odpowiedz';
  odpowiedz.hidden = true;

  const element = document.createElement('section');
  element.className = 'da-panel da-biblioteka';
  element.append(naglowek('Biblioteka ekspertów'), pasek, okno.element, odpowiedz);

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  async function duplikuj(zrodlowy: Agent): Promise<void> {
    powiedz(`Duplikowanie eksperta „${zrodlowy.name}"…`, true);
    const wynik = await stan.zrodlo.utworz({
      nazwa: `${zrodlowy.name} (kopia)`,
      opis: zrodlowy.description ?? '',
      instrukcje: zrodlowy.systemPrompt ?? '',
      kanal: zrodlowy.channelId ?? '',
      model: zrodlowy.model ?? '',
      // Imię własne i favikon idą do kopii razem z resztą tożsamości eksperta źródłowego.
      imie: zrodlowy.displayName === undefined ? '' : `${zrodlowy.displayName} (kopia)`,
      favikon: zrodlowy.favicon ?? '',
      // Odstępstwo jedzie do kopii razem z tożsamością, zgodnie z promptem systemowym eksperta źródłowego.
      zastepuje: czyZastepuje(zrodlowy),
      // Zasięg i pamięć idą do kopii razem z resztą tożsamości eksperta źródłowego.
      widocznosc: zrodlowy.visibility,
      poziomyPamieci: zrodlowy.memoryLevels,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Duplikowanie eksperta', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const kopia = wynik.wynik.agent;
    const umiejetnosci = zrodlowy.skillIds ?? [];
    for (const umiejetnosc of umiejetnosci) {
      const przypisanie = await stan.zrodlo.dodajUmiejetnosc(kopia.id, umiejetnosc);
      if (!przypisanie.udany) {
        // Kopia już jest w rdzeniu; wykaz odświeżamy przed komunikatem o odmowie przypisania.
        stan.wybierz(kopia.id);
        await stan.odswiez();
        powiedz(
          opisOdmowy(
            `Kopiowanie umiejętności ${umiejetnosc}`,
            przypisanie.blad?.code,
            przypisanie.blad?.message,
          ),
          false,
        );
        return;
      }
    }
    stan.wybierz(kopia.id);
    await stan.odswiez();
    powiedz(
      `Duplikat „${kopia.name}" założony. Skopiowano tożsamość i ${umiejetnosci.length} umiejętności; ` +
        'konektorów i uprawnień rdzeń nie oddaje jeszcze odczytem, więc kopia ich nie ma.',
      true,
    );
  }

  async function usun(ekspert: Agent): Promise<void> {
    powiedz(`Usuwanie eksperta „${ekspert.name}"…`, true);
    const wynik = await stan.zrodlo.usun(ekspert.id);
    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Usunięcie eksperta', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    if (!wynik.wynik.deleted) {
      powiedz(`Rdzeń nie znalazł eksperta ${ekspert.id} — nie było czego usuwać.`, false);
    } else {
      powiedz(`Ekspert „${ekspert.name}" usunięty.`, true);
    }
    await stan.odswiez();
  }

  function karta(ekspert: Agent): HTMLElement {
    return utworzKarteEksperta(
      ekspert,
      stan.wybrany()?.id === ekspert.id,
      {
        naWybor: () => stan.wybierz(ekspert.id),
        naDuplikowanie: () => void duplikuj(ekspert),
        naUsuniecie: () => void usun(ekspert),
      },
      przypisania === null ? undefined : (przypisania.get(ekspert.id) ?? 0),
    );
  }

  /** Zawężenie po polach bytu, który już przyszedł — bez pytania rdzenia ponownie. */
  function przesiej(eksperci: readonly Agent[]): readonly Agent[] {
    const zasieg = filtrWidocznosci.kontrolka.value;
    const stanCzynnosci = filtrStanu.kontrolka.value;
    return eksperci.filter((ekspert) => {
      if (zasieg !== '' && ekspert.visibility !== zasieg) return false;
      if (stanCzynnosci === 'czynny' && !ekspert.enabled) return false;
      if (stanCzynnosci === 'wylaczony' && ekspert.enabled) return false;
      return true;
    });
  }

  szukanie.addEventListener('change', () => void stan.odswiez(szukanie.value));
  odczyt.addEventListener('click', () => void stan.odswiez());
  filtrWidocznosci.kontrolka.addEventListener('change', () => odswiez());
  filtrStanu.kontrolka.addEventListener('change', () => odswiez());
  nowy.addEventListener('click', () => {
    stan.wybierz(null);
    powiedz('Edytor przestawiony na eksperta nowego — wypełnij tożsamość i zapisz.', true);
  });

  function odswiez(): void {
    if (stan.faza() === 'ladowanie') {
      okno.ladowanie('Odczyt biblioteki ekspertów w toku…');
      return;
    }
    if (stan.faza() === 'blad') {
      okno.blad(stan.powod());
      return;
    }
    const eksperci = stan.eksperci();
    const widoczni = przesiej(eksperci);
    lista.replaceChildren(...widoczni.map(karta));
    if (eksperci.length === 0) {
      okno.puste('Biblioteka jest pusta. Załóż pierwszego eksperta przyciskiem „Nowy ekspert”.');
      return;
    }
    // Pustka po zawężeniu różni się od pustej biblioteki i jest komunikowana wprost Operatorowi.
    if (widoczni.length === 0) {
      okno.puste(
        `Filtr nie przepuścił ani jednego z ${eksperci.length} ekspertów biblioteki. ` +
          'Rozszerz zasięg widoczności albo stan, aby zobaczyć pozostałych.',
      );
      return;
    }
    okno.gotowe();
  }

  return {
    element,
    odswiez,
    ustawPrzypisania(liczby) {
      przypisania = liczby;
      odswiez();
    },
  };
}

/** Nagłówek okna operacyjnego niesie nazwę okna zgodną z wykazem inwentarza modułów budowanych w tej sekcji aplikacji. */
export function naglowek(nazwa: string): HTMLElement {
  const element = document.createElement('h3');
  element.className = 'da-okno__tytul';
  element.textContent = nazwa;
  return element;
}
