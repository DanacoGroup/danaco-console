import type { WindowStateGetResponse } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWyboru, ustawPozycje } from '../../modele/kontrolki-formularza';
import { utworzStanOknaStudio } from './stan-okna-studio';
import type { StanStudio } from './stan-studio';
import type { ZrodloOsadzenia } from './zrodlo-osadzenia';

/** Pas osadzenia modułu wskazuje okno komunikacji sesji, w imieniu którego moduł pracuje, wspólnie dla wszystkich jego okien operacyjnych. */
export interface OsadzenieModulu {
  element: HTMLElement;
  /** Odczytuje okna sesji i parametry wykonania okna wybranego. */
  wczytaj(idSesji: string): Promise<void>;
}

export function utworzOsadzenieModulu(
  stan: StanStudio,
  zrodlo: ZrodloOsadzenia,
): OsadzenieModulu {
  const pasStanu = utworzStanOknaStudio();

  const wybor = poleWyboru({ etykieta: 'Okno komunikacji sesji' }, []);
  wybor.element.append(
    utworzDymekObjasnienia(
      'Komendy modułu Studio działają na oknie komunikacji sesji: jego identyfikator idzie w polu windowId ' +
        'komend studio.document.open i studio.contextual.op.',
      { powloka: 'ms-dymek', znak: 'ms-dymek__znak' },
    ),
  );

  const parametry = document.createElement('dl');
  parametry.className = 'ms-osadzenie__parametry';

  pasStanu.tresc.append(wybor.element, parametry);

  const element = document.createElement('section');
  element.className = 'ms-osadzenie';
  element.dataset['okno'] = 'studio.osadzenie';
  element.setAttribute('aria-label', 'Osadzenie modułu Studio w sesji');
  element.append(pasStanu.element);

  function wpisz(nazwa: string, wartosc: string): void {
    const etykieta = document.createElement('dt');
    etykieta.textContent = nazwa;
    const tresc = document.createElement('dd');
    tresc.textContent = wartosc;
    parametry.append(etykieta, tresc);
  }

  /** Wypisuje samo wskazanie okna: pozostałe parametry wykonania stoją w panelu Sterowanie okna. */
  function pokazParametry(odpowiedz: WindowStateGetResponse): void {
    parametry.replaceChildren();
    wpisz('Okno kontekstu', odpowiedz.window.title ?? odpowiedz.window.id);
  }

  async function wczytajParametry(idOkna: string): Promise<void> {
    const wynik = await zrodlo.stanOkna(idOkna);
    if (!wynik.udany || wynik.wynik === undefined) {
      pasStanu.blad(opisOdmowy('Odczyt stanu okna', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    pokazParametry(wynik.wynik);
    pasStanu.gotowe();
  }

  wybor.kontrolka.addEventListener('change', () => {
    stan.ustawOkno(wybor.kontrolka.value);
    void wczytajParametry(wybor.kontrolka.value);
  });

  return {
    element,

    async wczytaj(idSesji) {
      if (idSesji === '') {
        pasStanu.puste('Sesja nie jest otwarta', 'Moduł czeka na kartę sesji z powłoki.');
        return;
      }
      pasStanu.ladowanie('Odczyt okien komunikacji sesji…');
      const wynik = await zrodlo.okna(idSesji);
      if (!wynik.udany || wynik.wynik === undefined) {
        pasStanu.blad(opisOdmowy('Odczyt okien sesji', wynik.blad?.code, wynik.blad?.message));
        return;
      }
      const okna = wynik.wynik.windows;
      ustawPozycje(
        wybor.kontrolka,
        okna.map((okno) => ({
          wartosc: okno.id,
          etykieta: `${okno.moduleId} · ${okno.title ?? okno.id}`,
        })),
      );
      const pierwsze = okna[0];
      if (pierwsze === undefined) {
        pasStanu.puste('Sesja bez okna komunikacji', 'Komendy modułu potrzebują okna — załóż je w powłoce.');
        return;
      }
      // Wybór ustawiamy wprost, bo pusty windowId poleciałby do rdzenia jako żądanie bez okna.
      if (!okna.some((okno) => okno.id === wybor.kontrolka.value)) {
        wybor.kontrolka.value = pierwsze.id;
      }
      stan.ustawOkno(wybor.kontrolka.value);
      await wczytajParametry(wybor.kontrolka.value);
    },
  };
}
