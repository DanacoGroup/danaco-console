import {
  StudioTaskState,
  type StudioActor,
  type StudioDocumentTask,
  type StudioTaskPlan,
} from '../../../../shared/contract';

/**
 * Kolejka zadań rozkładu to jeden wspólny wykaz zleceń, kroków łańcucha operacji i pozycji wsadu;
 * niesie czynności, którymi wykaz sięga do rdzenia — pominięcie, ponowienie, pokazanie wyniku.
 */
export interface CzynnosciWykazuZadan {
  /** Pomija zadanie decyzją Operatora. */
  pomin(idZadania: string): void;
  /** Ponawia zadanie nieudane — wraca do czekania. */
  ponow(idZadania: string): void;
  /** Otwiera wynik zadania (propozycja albo wersja) w oknie obok. */
  pokazWynik(zadanie: StudioDocumentTask): void;
}

export interface WykazZadan {
  element: HTMLElement;
  /** Przerysowuje wykaz wedle rozkładu; `null` czyści go do stanu pustego. */
  odswiez(rozklad: StudioTaskPlan | null): void;
}

/** Nazwa stanu zadania czytelna dla Operatora kolejki — słowna, nie kod wyliczenia stanu zadania rozkładu. */
const NAZWA_STANU: Record<string, string> = {
  [StudioTaskState.Pending]: 'czeka',
  [StudioTaskState.Running]: 'w realizacji',
  [StudioTaskState.Done]: 'zakończone',
  [StudioTaskState.Failed]: 'nieudane',
  [StudioTaskState.Blocked]: 'wstrzymane',
  [StudioTaskState.Skipped]: 'pominięte',
};

/** Nazwa rodzaju zadania rozkładu czytelna dla Operatora kolejki, zamiast wewnętrznego oznaczenia rodzaju zadania. */
const NAZWA_RODZAJU: Record<string, string> = {
  research: 'źródła',
  draft: 'brzmienie',
  format: 'postać',
  apparatus: 'aparat',
  review: 'przejrzenie',
  proofread: 'język',
  export: 'wydanie',
  custom: 'własne',
};

export function utworzWykazZadan(czynnosci: CzynnosciWykazuZadan): WykazZadan {
  const lista = document.createElement('ol');
  lista.className = 'petla-kolejka';
  lista.setAttribute('aria-label', 'Kolejka zadań rozkładu');

  const pustka = document.createElement('p');
  pustka.className = 'dn-tekst-3';
  pustka.textContent =
    'Rozkładu nie ma. Zlecenie dokumentowe wydane w oknie rozmowy rozkłada się ' +
    'na zadania komendą studio.plan.create — dopiero wtedy kolejka ma co pokazać.';

  const element = document.createElement('div');
  element.className = 'petla-zadania';
  element.append(pustka, lista);

  function wiersz(zadanie: StudioDocumentTask, rozklad: StudioTaskPlan): HTMLElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta petla-zadanie';
    pozycja.dataset['stan'] = zadanie.state;

    const naglowek = document.createElement('div');
    naglowek.className = 'petla-zadanie__naglowek';

    const znacznik = document.createElement('span');
    znacznik.className = 'dn-plakietka petla-zadanie__stan';
    znacznik.textContent = NAZWA_STANU[zadanie.state] ?? zadanie.state;
    // Wskaźnik obrotu widnieje tylko przy zadaniu w biegu, nie przy zadaniu zakończonym.
    if (zadanie.state === StudioTaskState.Running) {
      const obrot = document.createElement('span');
      obrot.className = 'dn-spinner';
      obrot.setAttribute('role', 'status');
      obrot.setAttribute('aria-label', 'Zadanie w realizacji');
      naglowek.append(obrot);
    }

    const nazwa = document.createElement('span');
    nazwa.className = 'petla-zadanie__nazwa';
    nazwa.textContent = `${zadanie.order}. ${zadanie.title}`;

    const rodzaj = document.createElement('span');
    rodzaj.className = 'dn-plakietka petla-zadanie__rodzaj';
    rodzaj.textContent = NAZWA_RODZAJU[zadanie.kind] ?? zadanie.kind;

    naglowek.append(nazwa, rodzaj, znacznik);
    pozycja.append(naglowek);

    const wykonawca = opisWykonawcy(zadanie.assignedTo);
    if (wykonawca !== '') {
      const wiersz = document.createElement('p');
      wiersz.className = 'dn-tekst-3 petla-zadanie__wykonawca';
      wiersz.textContent = `Wykonawca: ${wykonawca}`;
      pozycja.append(wiersz);
    }

    if (zadanie.dependsOn !== undefined && zadanie.dependsOn.length > 0) {
      const zaleznosc = document.createElement('p');
      zaleznosc.className = 'dn-tekst-3 petla-zadanie__zaleznosc';
      const nazwy = zadanie.dependsOn
        .map((kod) => rozklad.tasks?.find((wpis) => wpis.id === kod)?.title ?? kod)
        .join(', ');
      zaleznosc.textContent = `Stoi na: ${nazwy}`;
      pozycja.append(zaleznosc);
    }

    // Powód niepowodzenia stoi przy zadaniu, nie w dzienniku obok, więc jest widoczny od razu.
    if (zadanie.failureReason !== undefined && zadanie.failureReason !== '') {
      const powod = document.createElement('p');
      powod.className = 'dn-tekst-3 petla-zadanie__powod';
      powod.textContent = zadanie.failureReason;
      pozycja.append(powod);
    }

    const czynnosciWiersza = document.createElement('div');
    czynnosciWiersza.className = 'petla-zadanie__czynnosci';

    if (zadanie.result !== undefined && zadanie.result !== '') {
      czynnosciWiersza.append(
        przycisk('Pokaż wynik', () => czynnosci.pokazWynik(zadanie)),
      );
    }
    if (zadanie.state === StudioTaskState.Failed) {
      czynnosciWiersza.append(przycisk('Ponów', () => czynnosci.ponow(zadanie.id)));
    }
    if (
      zadanie.state === StudioTaskState.Pending ||
      zadanie.state === StudioTaskState.Blocked ||
      zadanie.state === StudioTaskState.Failed
    ) {
      czynnosciWiersza.append(przycisk('Pomiń', () => czynnosci.pomin(zadanie.id)));
    }
    if (czynnosciWiersza.childElementCount > 0) pozycja.append(czynnosciWiersza);

    return pozycja;
  }

  return {
    element,

    odswiez(rozklad) {
      lista.replaceChildren();
      if (rozklad === null) {
        pustka.hidden = false;
        lista.hidden = true;
        return;
      }
      const zadania = rozklad.tasks ?? [];
      if (zadania.length === 0) {
        pustka.hidden = false;
        pustka.textContent =
          `Rozkład „${rozklad.order}" nie ma ani jednego zadania w wybranym zawężeniu. ` +
          'To nie znaczy, że rozkład jest pusty — zdejmij zawężenie stanu, żeby zobaczyć całość.';
        lista.hidden = true;
        return;
      }
      pustka.hidden = true;
      lista.hidden = false;
      for (const zadanie of zadania) lista.append(wiersz(zadanie, rozklad));
    },
  };
}

/**
 * Opis wykonawcy zadania dla Operatora: wykonawca nienazwany mówi prawdę o zadaniu powierzonym
 * komuś bez wskazania eksperta, nie o pustym miejscu.
 */
export function opisWykonawcy(wykonawca: StudioActor | undefined): string {
  if (wykonawca === undefined) return '';
  const czesci: string[] = [];
  if (wykonawca.agentName !== undefined && wykonawca.agentName !== '') {
    czesci.push(wykonawca.agentName);
  } else if (wykonawca.agentId !== undefined && wykonawca.agentId !== '') {
    czesci.push(wykonawca.agentId);
  } else {
    czesci.push(wykonawca.kind === 'model' ? 'wykonawca nienazwany' : 'Operator');
  }
  if (wykonawca.agentVersion !== undefined && wykonawca.agentVersion !== '') {
    czesci.push(`wersja ${wykonawca.agentVersion}`);
  }
  if (wykonawca.subagentId !== undefined && wykonawca.subagentId !== '') {
    czesci.push(`podagent ${wykonawca.subagentId}`);
  }
  return czesci.join(' · ');
}

/** Mały przycisk zarysowany, bez wypełnienia tła, w kolorze obramowania — postać czynności wiersza tabeli. */
export function przycisk(napis: string, naKlik: () => void): HTMLButtonElement {
  const guzik = document.createElement('button');
  guzik.type = 'button';
  guzik.className = 'dn-btn dn-btn--zarys dn-btn--sm';
  guzik.textContent = napis;
  guzik.addEventListener('click', naKlik);
  return guzik;
}
