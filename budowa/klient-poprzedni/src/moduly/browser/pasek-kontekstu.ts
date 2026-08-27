import { przycisk } from '../../modele/kontrolki-formularza';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { KLASY_DYMKA, OBJASNIENIA, WYZWALACZE } from './etykiety-browser';
import { KLASA_PRZYCISKU } from './przyciski-browser';
import type { StanPrzegladania } from './stan-przegladania';
import type { RozszerzenieModulu } from './warstwy-widocznosci';

/**
 * Pasek kontekstu modułu Browser, będący warstwą pierwszą, widoczną bez
 * interakcji. Pasek jest widokiem stanu warstw widoczności: niczego nie trzyma,
 * a każdy jego wyzwalacz przestawia jedną pozycję tamtego bytu.
 */
export interface PasekKontekstu {
  element: HTMLElement;
  /** Nanosi stan warstw na wyzwalacze; buduje je przy pierwszym wywołaniu. */
  odswiez(): void;
}

export function utworzPasekKontekstu(stan: StanPrzegladania): PasekKontekstu {
  const znacznikModulu = document.createElement('span');
  znacznikModulu.className = 'mb-kontekst__znacznik';
  znacznikModulu.textContent = 'Browser';

  const znacznikOkna = document.createElement('span');
  znacznikOkna.className = 'mb-kontekst__znacznik';

  const wyzwalacze = document.createElement('div');
  wyzwalacze.className = 'mb-kontekst__wyzwalacze';

  const tryb = przycisk(WYZWALACZE.trybAdministracyjny, KLASA_PRZYCISKU.zarys);
  tryb.setAttribute('aria-pressed', 'false');
  tryb.addEventListener('click', () => {
    stan.warstwy.ustawTrybAdministracyjny(!stan.warstwy.trybAdministracyjny());
  });

  const element = document.createElement('nav');
  element.className = 'mb-kontekst';
  element.setAttribute('aria-label', 'Pasek kontekstu modułu Browser');
  element.append(
    znacznikModulu,
    znacznikOkna,
    wyzwalacze,
    tryb,
    utworzDymekObjasnienia(OBJASNIENIA.warstwyWidocznosci, KLASY_DYMKA),
  );

  /** Przyciski wyzwalaczy po kodzie rozszerzenia — budowane raz. */
  const przyciski = new Map<string, HTMLButtonElement>();

  function zbuduj(): void {
    for (const rozszerzenie of stan.warstwy.rozszerzenia()) {
      if (przyciski.has(rozszerzenie.kod)) continue;
      const kontrolka = przycisk(rozszerzenie.nazwa, KLASA_PRZYCISKU.zarys);
      kontrolka.dataset['warstwa'] = String(rozszerzenie.warstwa);
      kontrolka.addEventListener('click', () => {
        stan.warstwy.przelacz(rozszerzenie.kod);
      });
      przyciski.set(rozszerzenie.kod, kontrolka);
      wyzwalacze.append(kontrolka);
    }
  }

  return {
    element,

    odswiez() {
      zbuduj();
      znacznikOkna.textContent = opisOkna(stan);
      const administracyjny = stan.warstwy.trybAdministracyjny();
      tryb.setAttribute('aria-pressed', String(administracyjny));
      for (const rozszerzenie of stan.warstwy.rozszerzenia()) {
        const kontrolka = przyciski.get(rozszerzenie.kod);
        if (kontrolka === undefined) continue;
        kontrolka.hidden = rozszerzenie.warstwa === 4 && !administracyjny;
        kontrolka.setAttribute('aria-pressed', String(stan.warstwy.czyOdsloniete(rozszerzenie.kod)));
        kontrolka.textContent = napisWyzwalacza(rozszerzenie);
      }
    },
  };
}

/**
 * Napis wyzwalacza wraz z liczbą pozycji rozszerzenia. Znak strzałki mówi, że
 * pozycja rozwija kolumnę boczną, tak samo jak w znacznikach kontekstowych,
 * które noszą go w nazwie.
 */
function napisWyzwalacza(rozszerzenie: RozszerzenieModulu): string {
  const licznik = rozszerzenie.licznik?.();
  return `${rozszerzenie.nazwa} ▾${licznik === undefined ? '' : ` ${licznik}`}`;
}

/**
 * Znacznik okna przeglądarki: jego identyfikator albo powód, dla którego go nie
 * ma. Okno nieustalone nazywa się wprost, ponieważ pusty napis w pasku wyglądałby
 * na usterkę widoku, a nie na brak wiązania.
 */
function opisOkna(stan: StanPrzegladania): string {
  const identyfikator = stan.idOkna();
  return identyfikator === '' ? 'okno nieustalone' : `okno ${identyfikator}`;
}
