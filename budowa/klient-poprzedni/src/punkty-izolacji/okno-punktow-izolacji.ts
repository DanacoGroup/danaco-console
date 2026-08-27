import './punkty-izolacji.css';

import { ChangeKind, EventType } from '../../../shared/contract';
import { utworzRameOkna } from '../komponenty/rama-okna';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { PANELE_IZOLACJI, REJESTR_OBSZAROW, type ObszarIzolacji } from './obszary';
import { utworzStanWarstwy } from './stan-warstwy';
import { utworzStanZasiegu } from './stan-zasiegu';
import { utworzSterowanieWarstwa } from './sterowanie-warstwa';

/** Okno Punktów Izolacji — okno zarządca katalogu rdzenia, otwierane z listwy Ustawienia obok innych okien konfiguracji. */
export interface OknoPunktowIzolacji {
  /** Element `<dialog>` osadzany w dokumencie. */
  element: HTMLDialogElement;
  /** Otwiera okno i montuje obszar czynny. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje obszaru czynnego i usuwa okno z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoPunktowIzolacji(kanal: Kanal): OknoPunktowIzolacji {
  const rama = utworzRameOkna({
    kod: 'punkty-izolacji',
    tytul: 'Punkty izolacji',
    rola: 'zarządca',
    modul: 'Punkty izolacji',
    przeznaczenie:
      'Konfiguracja jedenastu punktów izolacji rdzenia — kontekst, zakres techniczny, profile, ' +
      'poziom zapisu, oś rozstrzygania i podgląd polityki efektywnej, na wybranej warstwie.',
    przedrostek: 'pi',
  });

  const warstwa = utworzStanWarstwy();
  const zasieg = utworzStanZasiegu();
  const sterowanieWarstwa = utworzSterowanieWarstwa(kanal, warstwa);

  // Ślad ostatniej zmiany przyszłej z rdzenia, ukryty, dopóki nic się nie zmieniło, bez zbędnego szumu.
  const slad = document.createElement('p');
  slad.className = 'pi-slad';
  slad.setAttribute('role', 'status');
  slad.hidden = true;

  rama.narzedzia.append(sterowanieWarstwa.element);

  // Sześć obszarów montowanych naraz, po dwa na panel, egzemplarz każdego powstaje raz na życie okna.
  const obszary: ObszarIzolacji[] = [];
  const kolumny = PANELE_IZOLACJI.map((panel) => {
    const kolumna = zbudujKolumne(panel.tytul, panel.opis, panel.kod);
    for (const wpis of REJESTR_OBSZAROW.filter((w) => w.panel === panel.kod)) {
      const obszar = wpis.utworz({ kanal, warstwa, zasieg });
      obszary.push(obszar);
      kolumna.append(zbudujMiejsceObszaru(wpis.nazwa, wpis.opis, obszar.element));
    }
    return kolumna;
  });

  const panele = document.createElement('div');
  panele.className = 'pi-panele';
  panele.append(...kolumny);

  rama.cialo.append(slad, panele);

  /** Odczyt wszystkich obszarów od nowa — po otwarciu okna i po każdej zmianie wspólnego stanu. */
  function odswiezWszystkie(): void {
    for (const obszar of obszary) obszar.odswiez();
  }

  // Zmiana warstwy i zmiana zasięgu dotyczą każdego obszaru naraz, więc czytają od nowa wszystkie razem.
  warstwa.naZmiane(() => odswiezWszystkie());
  zasieg.naZmiane(() => odswiezWszystkie());

  // Nanosi zmianę zgłoszoną przez rdzeń: ślad na wierzchu i ponowny odczyt obszaru czynnego okna.
  function nanieszZmiane(zdanie: string): void {
    slad.hidden = false;
    slad.textContent =
      `${zdanie} Odczytano od nowa. Sprawcy nie znamy: to zdarzenie nie niesie pola ` +
      '„actor" — nie wiadomo, czy punkt przestawił Operator z innego urządzenia, ' +
      'czy asystent.';
    if (!element.open) return;
    odswiezWszystkie();
  }

  const odsubskrybowania: Odsubskrybuj[] = [
    kanal.naZdarzenie(EventType.IsolationProfileChanged, (tresc) => {
      const nazwa = tresc.profileId === '' ? 'bez nazwy' : tresc.profileId;
      nanieszZmiane(`Profil izolacji „${nazwa}" — ${slowoZmiany(tresc.change)}.`);
    }),
    kanal.naZdarzenie(EventType.IsolationPolicyChanged, (tresc) => {
      const okno = tresc.windowId === '' ? 'nie wskazano okna' : `okno ${tresc.windowId}`;
      nanieszZmiane(`Polityka izolacji (${okno}) — ${slowoZmiany(tresc.change)}.`);
    }),
  ];

  const element = document.createElement('dialog');
  element.className = 'dn-modal pi-okno';
  element.setAttribute('aria-label', 'Okno punktów izolacji');

  const cialo = document.createElement('div');
  cialo.className = 'dn-modal-cialo pi-okno__cialo';
  cialo.append(rama.element);

  const zamknij = document.createElement('button');
  zamknij.type = 'button';
  zamknij.className = 'dn-btn dn-btn--atrament';
  zamknij.textContent = 'Zamknij';
  zamknij.addEventListener('click', () => element.close());

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka pi-okno__stopka';
  stopka.append(zamknij);

  element.append(cialo, stopka);

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      sterowanieWarstwa.odswiez();
      odswiezWszystkie();
    },

    zamknij: () => element.close(),

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      for (const obszar of obszary.splice(0)) obszar.zamknij?.();
      element.remove();
    },
  };
}

/** Kolumna panelu okna: tytuł, zdanie o przeznaczeniu kolumny oraz miejsce na osadzone obszary izolacji. */
function zbudujKolumne(tytul: string, opis: string, kod: string): HTMLElement {
  const naglowek = document.createElement('h3');
  naglowek.className = 'pi-panel__tytul';
  naglowek.textContent = tytul;

  const zdanie = document.createElement('p');
  zdanie.className = 'pi-panel__opis';
  zdanie.textContent = opis;

  const element = document.createElement('section');
  element.className = 'pi-panel';
  element.dataset['panel'] = kod;
  element.setAttribute('aria-label', tytul);
  element.append(naglowek, zdanie);
  return element;
}

/** Obszar osadzony w panelu: podpis obszaru nad jego treścią właściwą — panel mieści zawsze po dwa obszary. */
function zbudujMiejsceObszaru(nazwa: string, opis: string, tresc: HTMLElement): HTMLElement {
  const podpis = document.createElement('h4');
  podpis.className = 'pi-obszar__tytul';
  podpis.textContent = nazwa;

  const zdanie = document.createElement('p');
  zdanie.className = 'pi-obszar__opis';
  zdanie.textContent = opis;

  const element = document.createElement('div');
  element.className = 'pi-obszar';
  element.append(podpis, zdanie, tresc);
  return element;
}

/** Rodzaj zmiany po polsku, złożony z rodzaju zmiany samego kontraktu, a nie z napisu gotowego dla Operatora. */
function slowoZmiany(zmiana: ChangeKind): string {
  switch (zmiana) {
    case ChangeKind.Created:
      return 'utworzony';
    case ChangeKind.Updated:
      return 'zmieniony';
    case ChangeKind.Deleted:
      return 'usunięty';
    default:
      // Rodzaj spoza wykazu nie jest powodem do milczenia: pokazujemy go dosłownie, tak jak przyszedł.
      return `zmiana rodzaju „${String(zmiana)}"`;
  }
}

