import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przyciskAkcji } from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Zakładka „Zestawy i retencja” w Memory & Context Managerze obejmuje siedem komend: nazwane
 * konteksty pamięci, zasady retencji i miernik okna kontekstu.
 */
export interface PanelZestawow {
  element: HTMLElement;
  /** Odczyt kontekstów, zasad retencji i zajętości okna kontekstu. */
  wczytaj(): Promise<void>;
}

export function utworzPanelZestawow(stan: StanAssistant): PanelZestawow {
  const okno: StanOkna = utworzStanOkna();

  // Nazwane konteksty
  const nazwa = pole('Nazwa zestawu', 'np. projekt Atlas');
  const prompt = pole('Warstwa promptu systemowego', 'np. Mów zwięźle.');
  const zaloz = przyciskAkcji('Załóż zestaw', 'dn-btn dn-btn--sm dn-btn--zarys');
  zaloz.addEventListener('click', () => void zalozKontekst());

  const wykazKontekstow = document.createElement('ul');
  wykazKontekstow.className = 'ma-wykaz';
  wykazKontekstow.dataset['wykaz'] = 'konteksty-nazwane';

  // Zasady retencji
  const dni = pole('Wygaszanie po ilu dniach (0 = pamięć trwała)', '0');
  const wzorce = pole('Czego nigdy nie zapisywać (po przecinku)', 'numer karty, hasło');
  const zapiszZasade = przyciskAkcji('Zapisz zasadę retencji', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapiszZasade.addEventListener('click', () => void zapiszRetencje());

  const opisZasady = document.createElement('p');
  opisZasady.className = 'dn-pole-opis';

  // Miernik okna kontekstu
  const zmierz = przyciskAkcji('Zmierz zajętość okna kontekstu', 'dn-btn dn-btn--sm dn-btn--zarys');
  zmierz.addEventListener('click', () => void zmierzZajetosc());

  const miernik = document.createElement('p');
  miernik.className = 'dn-pole-opis';
  miernik.dataset['miernik'] = 'kontekst';

  okno.tresc.append(
    sekcjaZestawow('Nazwane zestawy pamięci', [nazwa, prompt, zaloz, wykazKontekstow]),
    sekcjaZestawow('Retencja i wygaszanie', [dni, wzorce, zapiszZasade, opisZasady]),
    sekcjaZestawow('Okno kontekstu rozmowy', [zmierz, miernik]),
  );

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'zestawy';
  element.append(okno.element);

  async function wczytaj(): Promise<void> {
    okno.ladowanie('Odczyt zestawów pamięci i zasad retencji…');
    await Promise.all([wczytajKonteksty(), wczytajZasady()]);
  }

  async function wczytajKonteksty(): Promise<void> {
    const wynik = await stan.konteksty.konteksty(stan.idSesji(), true);
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt zestawów pamięci', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const czynny = wynik.wynik.activeId ?? '';
    wykazKontekstow.replaceChildren();
    for (const kontekst of wynik.wynik.contexts) {
      const wiersz = document.createElement('li');
      wiersz.className = 'ma-wykaz__wiersz';
      wiersz.dataset['kontekst'] = kontekst.id;

      const opis = document.createElement('span');
      opis.textContent =
        kontekst.id === czynny ? `${kontekst.name} — czynny w tej karcie` : kontekst.name;

      const uaktywnij = przyciskAkcji('Uaktywnij', 'dn-btn dn-btn--xs dn-btn--duch');
      uaktywnij.addEventListener('click', () => void uaktywnijKontekst(kontekst.id));

      const usun = przyciskAkcji('Usuń wskazanie', 'dn-btn dn-btn--xs dn-btn--duch');
      usun.title =
        'Kasuje sam zestaw wskazań. Wpisy pamięci zostają nietknięte — kontekst ' +
        'nie jest ich właścicielem.';
      usun.addEventListener('click', () => void usunKontekst(kontekst.id));

      wiersz.append(opis, uaktywnij, usun);
      wykazKontekstow.append(wiersz);
    }
    if (wynik.wynik.contexts.length === 0) {
      okno.puste('Nie ma jeszcze ani jednego nazwanego zestawu pamięci.');
      return;
    }
    okno.gotowe();
  }

  async function zalozKontekst(): Promise<void> {
    if (nazwa.value.trim() === '') {
      okno.blad('Zestaw bez nazwy nie miałby czego pokazać w selektorze.');
      return;
    }
    const wynik = await stan.konteksty.zapiszKontekst({
      id: '',
      nazwa: nazwa.value.trim(),
      opis: '',
      poziomy: [],
      wpisy: [],
      promptSystemowy: prompt.value,
      czynny: true,
    });
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Zapis zestawu pamięci', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    nazwa.value = '';
    prompt.value = '';
    await wczytajKonteksty();
  }

  async function uaktywnijKontekst(id: string): Promise<void> {
    const sesja = stan.idSesji();
    if (sesja === '') {
      okno.blad(
        'Powłoka nie podała karty sesji. Kontekst czynny należy do karty, więc bez ' +
          'niej nie ma czego uaktywnić.',
      );
      return;
    }
    const wynik = await stan.konteksty.uaktywnijKontekst(id, sesja);
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Aktywacja zestawu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajKonteksty();
  }

  async function usunKontekst(id: string): Promise<void> {
    const wynik = await stan.konteksty.usunKontekst(id);
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Usunięcie zestawu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytajKonteksty();
  }

  async function wczytajZasady(): Promise<void> {
    const wynik = await stan.konteksty.zasadyRetencji();
    if (!wynik.udany || wynik.wynik === undefined) {
      opisZasady.textContent = opisOdmowy(
        'Odczyt zasad retencji',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    const pierwsza = wynik.wynik.policies[0];
    if (pierwsza === undefined) {
      opisZasady.textContent =
        'Żadna zasada retencji nie obowiązuje — pamięć jest trwała, dopóki Operator ' +
        'nie skasuje wpisu.';
      return;
    }
    dni.value = String(pierwsza.ttlDays);
    opisZasady.textContent =
      pierwsza.ttlDays === 0
        ? 'Obowiązująca zasada: pamięć trwała.'
        : `Obowiązująca zasada: wpisy wygasają po ${String(pierwsza.ttlDays)} dniach.`;
  }

  async function zapiszRetencje(): Promise<void> {
    const liczba = Number.parseInt(dni.value.trim(), 10);
    if (Number.isNaN(liczba) || liczba < 0) {
      opisZasady.textContent =
        'Liczba dni musi być nieujemna; zero znaczy pamięć trwałą.';
      return;
    }
    const wynik = await stan.konteksty.zapiszZasadeRetencji({
      dniWygasania: liczba,
      wrazliweDomyslnie: false,
      wzorceNigdyNieZapisywane: wzorce.value
        .split(',')
        .map((wzorzec) => wzorzec.trim())
        .filter((wzorzec) => wzorzec !== ''),
      czynna: true,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      opisZasady.textContent = opisOdmowy(
        'Zapis zasady retencji',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    // Liczba wpisów zastanych jest policzona po stronie rdzenia; zasada ich teraz nie rusza.
    opisZasady.textContent =
      `Zasada zapisana. Obejmuje zapisy kolejne; wpisów zastanych, których ` +
      `dotknie przy najbliższym wygaszaniu, jest ${String(wynik.wynik.affectedEntries)} ` +
      `— żaden nie został skasowany teraz.`;
  }

  async function zmierzZajetosc(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      miernik.textContent =
        'Rdzeń nie oddał okna modułu dla tej sesji, a kontekst należy do okna.';
      return;
    }
    const wynik = await stan.konteksty.zajetoscKontekstu(idOkna, stan.idSesji());
    if (!wynik.udany || wynik.wynik === undefined) {
      miernik.textContent = opisOdmowy(
        'Pomiar zajętości okna kontekstu',
        wynik.blad?.code,
        wynik.blad?.message,
      );
      return;
    }
    const pomiar = wynik.wynik;
    if (!pomiar.available) {
      miernik.textContent =
        pomiar.reason !== undefined && pomiar.reason !== ''
          ? `Pomiar niewykonalny: ${pomiar.reason}`
          : 'Pomiar niewykonalny, a rdzeń nie podał powodu.';
      return;
    }
    miernik.textContent =
      `Okno kontekstu: ${String(pomiar.usage.usedTokens)} / ` +
      `${String(pomiar.usage.limitTokens)} żetonów (policzono słownikiem ` +
      `${pomiar.usage.tokenizer}).`;
  }

  return { element, wczytaj };
}

/** Sekcja panelu niesie tytuł wraz z kontrolkami jednego obszaru, oddzielając wizualnie trzy grupy komend. */
function sekcjaZestawow(tytul: string, dzieci: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ma-panel__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ma-panel__sekcja';
  element.append(naglowek, ...dzieci);
  return element;
}
