import { utworzPanelRodzin } from './panel-rodzin';
import { sekcjeAutomatyzacji } from './sekcje-rodzin';
import {
  AutomationStepKind,
  Command,
  type AutomationStep,
  type AutomationWorkflow,
} from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import {
  poleLogiczne,
  poleTekstowe,
  poleWielowierszowe,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { nowyIdentyfikator } from '../../protokol/identyfikator';
import { utworzCzynnosciAutomatyk, type DefinicjaAutomatyki } from './czynnosci-automatyk';
import { utworzEdytorScenariusza } from './edytor-scenariusza';
import {
  KLASY_DYMKA,
  KODY_OKIEN,
  KODY_PANELI,
  OBJASNIENIA,
  STANY_PUSTE,
  WYKAZY,
} from './etykiety-browser';
import { utworzLimityPrzebiegu } from './limity-przebiegu';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import { utworzStanOkna } from './stan-okna';
import type { StanPrzegladania } from './stan-przegladania';
import type { RozszerzenieModulu } from './warstwy-widocznosci';
import {
  nanieStanWykazu,
  wierszAutomatyki,
  wierszeKrokow,
  zdanieOPrzebiegach,
} from './wykaz-automatyk';

/**
 * Automation Studio — okno operacji i scenariuszy modułu Browser, warstwa
 * trzecia, otwierane z menu operacji. Jedna odpowiedzialność: złożenie okna
 * i wykaz automatyk.
 */
export interface OknoAutomationStudio {
  element: HTMLElement;
  odswiez(): void;
  /** Panele warstwy czwartej okna — moduł rejestruje je w warstwach widoczności. */
  panele: readonly RozszerzenieModulu[];
}

export function utworzOknoAutomationStudio(stan: StanPrzegladania): OknoAutomationStudio {
  const okno = utworzStanOkna(STANY_PUSTE.automatyzacja);
  const odpowiedz = utworzWierszOdpowiedzi();
  const powiedz = (tresc: string, ok: boolean): void => odpowiedz.pokaz(tresc, ok);
  const czynnosci = utworzCzynnosciAutomatyk(stan, powiedz);
  const edytor = utworzEdytorScenariusza();
  const limity = utworzLimityPrzebiegu(stan.pokrycie);

  /** Automatyka wybrana do zmiany; pusty napis znaczy „scenariusz nowy". */
  let wybrana = '';
  let automatyki: readonly AutomationWorkflow[] = [];
  let odmowa = '';
  let wOdczycie = false;
  // Kroki ostatnio przyjęte z widoku tekstowego, nie z każdej litery wpisywanej w pole.
  let krokiPrzyjete: AutomationStep[] = [];
  /** Zdanie o ostatnio odczytanych przebiegach, po identyfikatorze automatyki. */
  const przebiegi = new Map<string, string>();

  const nazwa = poleTekstowe({ etykieta: 'Nazwa scenariusza' });
  const opis = poleWielowierszowe({ etykieta: 'Opis scenariusza' }, 2);
  const czynna = poleLogiczne({ etykieta: 'Scenariusz czynny' });
  const cron = poleTekstowe({
    etykieta: 'Cykliczność (cron)',
    podpowiedz: '0 6 * * *',
  });
  const harmonogramCzynny = poleLogiczne({ etykieta: 'Harmonogram obowiązuje' });

  const wykazKrokow = document.createElement('ol');
  wykazKrokow.className = 'mb-scenariusz__kroki';

  const lista = document.createElement('ul');
  lista.className = 'mb-automatyki__lista';

  const uwaga = document.createElement('p');
  uwaga.className = 'dn-pole-opis mb-uwaga';
  uwaga.textContent = WYKAZY.automatyki;

  /** Definicja złożona z formularza i z kroków widoku tekstowego. */
  function definicja(zestaw: readonly AutomationStep[]): DefinicjaAutomatyki {
    return {
      identyfikator: wybrana,
      nazwa: nazwa.kontrolka.value,
      opis: opis.kontrolka.value,
      czynna: czynna.kontrolka.checked,
      kroki: zestaw,
    };
  }

  /** Kroki przyjęte z widoku tekstowego; `null` znaczy „treść odrzucona”. */
  function krokiScenariusza(): AutomationStep[] | null {
    const odczyt = edytor.odczytaj();
    if (odczyt.blad !== '') {
      powiedz(odczyt.blad, false);
      return null;
    }
    return odczyt.kroki;
  }

  function dodajKrokZeStrony(): void {
    const migawka = stan.migawka();
    if (migawka === null) {
      powiedz('Brak migawki — najpierw przejdź do strony w Browser Window.', false);
      return;
    }
    const zastane = krokiScenariusza();
    if (zastane === null) return;
    zastane.push({
      id: nowyIdentyfikator('krok'),
      name: `Przejdź do ${migawka.url}`,
      kind: AutomationStepKind.Command,
      // Nazwa komendy pochodzi z generatu kontraktu, nie z literału przepisanego ręcznie.
      command: Command.BrowserNavigate,
      params: { url: migawka.url },
      order: zastane.length + 1,
    });
    edytor.wczytaj(zastane);
    powiedz(
      `Krok dopisany: przejście pod ${migawka.url}. Kroków w scenariuszu: ${zastane.length}.`,
      true,
    );
    odswiez();
  }

  async function zapisz(): Promise<void> {
    const zestaw = krokiScenariusza();
    if (zestaw === null) return;
    const powod = limity.sprawdz(zestaw);
    if (powod !== '') {
      powiedz(powod, false);
      return;
    }
    const zapisana = await czynnosci.zapisz(definicja(zestaw));
    if (zapisana === null) return;
    wypelnij(zapisana);
    await odczytaj();
  }

  async function odczytaj(): Promise<void> {
    wOdczycie = true;
    odswiez();
    const wykaz = await czynnosci.odczytaj();
    wOdczycie = false;
    if (wykaz === null) {
      odmowa = 'Odczyt automatyk odmówiony — wykaz pokazuje stan sprzed odczytu.';
      odswiez();
      return;
    }
    odmowa = '';
    automatyki = wykaz;
    odswiez();
  }

  /** Wypełnia formularz definicją odczytaną z rdzenia. */
  function wypelnij(automatyka: AutomationWorkflow): void {
    wybrana = automatyka.id;
    nazwa.kontrolka.value = automatyka.name;
    opis.kontrolka.value = automatyka.description ?? '';
    czynna.kontrolka.checked = automatyka.enabled;
    edytor.wczytaj(automatyka.steps ?? []);
    odswiez();
  }

  const pasek = document.createElement('div');
  pasek.className = 'mb-panel__pasek';
  pasek.append(
    przyciskCzynnosci('Krok z bieżącej strony', KLASA_PRZYCISKU.glowny, dodajKrokZeStrony),
    utworzDymekObjasnienia(OBJASNIENIA.krokZeStrony, KLASY_DYMKA),
    przyciskCzynnosci('Zapisz scenariusz', KLASA_PRZYCISKU.zarys, () => void zapisz()),
    utworzDymekObjasnienia(OBJASNIENIA.automatyka, KLASY_DYMKA),
    przyciskCzynnosci('Odśwież z rdzenia', KLASA_PRZYCISKU.zarys, () => void odczytaj()),
    przyciskCzynnosci('Nowy scenariusz', KLASA_PRZYCISKU.zarys, () => {
      wybrana = '';
      nazwa.kontrolka.value = '';
      opis.kontrolka.value = '';
      czynna.kontrolka.checked = false;
      edytor.wczytaj([]);
      powiedz('Formularz opróżniony — zapis założy nową automatykę.', true);
      odswiez();
    }),
    przyciskCzynnosci('→ Przekaż do Automations', KLASA_PRZYCISKU.zarys, () => {
      const zestaw = krokiScenariusza();
      if (zestaw !== null) void czynnosci.przekaz(definicja(zestaw));
    }),
    utworzDymekObjasnienia(OBJASNIENIA.przekazanie, KLASY_DYMKA),
  );

  const paskHarmonogramu = document.createElement('div');
  paskHarmonogramu.className = 'mb-panel__pasek';
  paskHarmonogramu.append(
    przyciskCzynnosci('Ustaw harmonogram', KLASA_PRZYCISKU.zarys, () => {
      void czynnosci.ustawHarmonogram(wybrana, cron.kontrolka.value, harmonogramCzynny.kontrolka.checked);
    }),
    przyciskCzynnosci('Odczytaj harmonogram', KLASA_PRZYCISKU.zarys, () => {
      void (async () => {
        const wykaz = await czynnosci.odczytajHarmonogramy(wybrana);
        const pierwszy = wykaz?.[0];
        if (pierwszy === undefined) return;
        cron.kontrolka.value = pierwszy.cron ?? '';
        harmonogramCzynny.kontrolka.checked = pierwszy.enabled;
      })();
    }),
    przyciskCzynnosci('Odczytaj przebiegi', KLASA_PRZYCISKU.zarys, () => {
      void (async () => {
        const wykaz = await czynnosci.odczytajPrzebiegi(wybrana);
        if (wykaz === null) return;
        przebiegi.set(wybrana, zdanieOPrzebiegach(wykaz));
        odswiez();
      })();
    }),
    utworzDymekObjasnienia(OBJASNIENIA.harmonogram, KLASY_DYMKA),
  );

  // Nagrywarka makr i granice Wykonawcy idą teraz do rdzenia, nie giną już wraz z kartą.
  const rodziny = utworzPanelRodzin(sekcjeAutomatyzacji(stan));

  okno.tresc.append(lista, rodziny.element, odpowiedz.element);

  const element = document.createElement('section');
  element.className = 'mb-okno mb-okno--pomocnicze';
  element.dataset['okno'] = KODY_OKIEN.automatyzacja;
  element.setAttribute('aria-label', 'Automation Studio — scenariusze przeglądania');
  element.append(
    utworzNaglowekOkna({ tytul: 'Automation Studio', klasa: 'mb-okno__naglowek' }),
    nazwa.element,
    opis.element,
    czynna.element,
    wykazKrokow,
    pasek,
    cron.element,
    harmonogramCzynny.element,
    paskHarmonogramu,
    edytor.element,
    limity.element,
    uwaga,
    okno.element,
  );

  function odswiez(): void {
    const odczyt = edytor.odczytaj();
    if (odczyt.blad === '') krokiPrzyjete = odczyt.kroki;
    wykazKrokow.replaceChildren(...wierszeKrokow(krokiPrzyjete));
    lista.replaceChildren(
      ...automatyki.map((automatyka) =>
        wierszAutomatyki(
          automatyka,
          przebiegi.get(automatyka.id) ?? '',
          automatyka.id === wybrana,
          () => wypelnij(automatyka),
        ),
      ),
    );
    nanieStanWykazu(okno, { wOdczycie, odmowa, pozycji: automatyki.length });
  }

  okno.ustawPonowienie(() => void odczytaj());
  odswiez();

  return {
    element,
    odswiez,
    panele: [
      {
        kod: KODY_PANELI.edytorScenariusza,
        nazwa: 'Widok tekstowy scenariusza',
        warstwa: 4,
        element: edytor.element,
      },
      {
        kod: KODY_PANELI.limityPrzebiegu,
        nazwa: 'Limity przebiegu Wykonawcy',
        warstwa: 4,
        element: limity.element,
      },
    ],
  };
}
