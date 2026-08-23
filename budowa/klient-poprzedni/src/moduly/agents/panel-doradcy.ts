import './doradca.css';

import type { Agent, Channel } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { poleWielowierszowe, poleWyboru, przycisk, ustawPozycje } from '../../modele/kontrolki-formularza';
import {
  ETYKIETA_RADY,
  radaDoPrzeniesienia,
  zapiszProwenancje,
  zdanieOStanie,
  type ZapisProwenancji,
} from './prowenancja-rady';
import { kanalyDoradcze, type ZrodloDoradcy } from './zrodlo-doradcy';
import type { ZrodloZaplecza } from './zrodlo-zaplecza';

/**
 * Panel „Doradca” Agent Buildera — konsultacja eksperta modelem silniejszym.
 *
 * Widok pokazuje naraz kogo zapytano, o co i co doradca odpowiedział. Treść
 * rady stoi pod etykietą `ETYKIETA_RADY` i nie trafia do pola instrukcji
 * eksperta sama z siebie — przeniesienie jest osobnym kliknięciem i niesie
 * nagłówek prowenancji (`prowenancja-rady.ts`).
 *
 * Wybierany jest kanał, nie nazwa modelu: rejestr kanałów (`channel.list`) jest
 * jedynym miejscem, w którym rdzeń wie, czym się połączyć i czyim
 * poświadczeniem. Wykaz obejmuje kanały czynne poza kanałem bazowym eksperta.
 *
 * Konsultacji nie da się zapisać w rdzeniu — kontrakt nie ma takiej komendy —
 * więc wykaz żyje tyle, co widok.
 */
export interface PanelDoradcy {
  element: HTMLElement;
  /** Nanosi eksperta czynnego; pusty wybór wyłącza pytanie. */
  ustaw(ekspert: Agent | null): void;
  /** Odczytuje rejestr kanałów modelu — z nich bierze się wykaz doradców. */
  wczytajKanaly(): Promise<void>;
}

export function utworzPanelDoradcy(
  zrodlo: ZrodloDoradcy,
  zaplecze: ZrodloZaplecza,
  wstawDoInstrukcji: (tresc: string) => void,
): PanelDoradcy {
  let ekspertCzynny: Agent | null = null;
  let kanaly: readonly Channel[] = [];
  let konsultacjaWToku = false;

  const doradca = poleWyboru(
    {
      etykieta: 'Doradca',
      opis:
        'Kanał modelu, którym pytamy. Wykaz pomija kanał bazowy eksperta — ' +
        'konsultacja u samego siebie nie jest konsultacją.',
    },
    [],
  );

  const pytanie = poleWielowierszowe(
    {
      etykieta: 'Pytanie do doradcy',
      podpowiedz: 'O co pytamy doradcę w sprawie tego eksperta?',
      opis: 'Pytanie idzie do rdzenia dokładnie w tej postaci i w tej postaci wraca w prowenancji.',
    },
    4,
  );

  const zapytaj = przycisk('Zapytaj doradcę', 'dn-btn dn-btn--sm dn-btn--atrament');

  // Opis obok kontrolki zamiast wygaszenia kontrolki — druga z dwóch dróg
  // dopuszczonych przez zasadę zero blokad, obok komunikatu po naciśnięciu.
  const gotowosc = document.createElement('p');
  gotowosc.className = 'dn-pole-opis da-doradca__gotowosc';
  gotowosc.hidden = true;

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'da-odpowiedz';
  odpowiedz.hidden = true;

  const wykaz = document.createElement('ol');
  wykaz.className = 'da-doradca__wykaz';

  const nota = document.createElement('p');
  nota.className = 'dn-pole-opis da-granica';
  nota.textContent =
    'Konsultacje widać tylko w tym widoku: kontrakt nie ma komendy zapisującej ' +
    'prowenancję rady w rdzeniu, więc wykaz znika razem z oknem.';

  const element = document.createElement('section');
  element.className = 'da-panel da-doradca';
  element.append(
    tytul('Doradca — konsultacja modelem silniejszym'),
    doradca.element,
    pytanie.element,
    zapytaj,
    gotowosc,
    odpowiedz,
    wykaz,
    nota,
  );

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  /**
   * Powód, dla którego pytanie nie ma dziś jak pojechać; pusty znaczy gotowość.
   *
   * Powód nie odbiera przycisku i nie może tego robić: platforma nie stawia
   * bram, a niegotowość sygnalizuje się PO naciśnięciu — komunikatem — albo
   * opisem obok kontrolki. Przycisk wygaszony zabierał Operatorowi jedyną
   * drogę dowiedzenia się, czego brakuje, bo `title` bywa niedostępny
   * z klawiatury i milczy na urządzeniu dotykowym.
   */
  function powodNiegotowosci(): string {
    if (ekspertCzynny === null) return 'Wybierz eksperta w bibliotece, zanim zapytasz doradcę.';
    if (doradca.kontrolka.value === '') {
      return 'Rejestr kanałów nie ma kanału innego niż bazowy kanał tego eksperta.';
    }
    if (pytanie.kontrolka.value.trim() === '') return 'Wpisz pytanie do doradcy.';
    if (konsultacjaWToku) return 'Konsultacja w toku — poczekaj na odpowiedź doradcy.';
    return '';
  }

  /** Nanosi powód na opis obok kontrolki; przycisk zostaje klikalny zawsze. */
  function przelicz(): void {
    const brak = powodNiegotowosci();
    zapytaj.title = brak;
    gotowosc.textContent = brak;
    gotowosc.hidden = brak === '';
  }

  async function konsultuj(): Promise<void> {
    // Niegotowość rozstrzyga się po naciśnięciu, nie przed nim: Operator
    // dostaje zdanie o tym, czego brakuje, zamiast kontrolki, która nie reaguje.
    const brak = powodNiegotowosci();
    if (brak !== '') {
      powiedz(brak, false);
      return;
    }
    const ekspert = ekspertCzynny;
    if (ekspert === null) return;
    const idKanalu = doradca.kontrolka.value;
    const wybrany = kanaly.find((kanal) => kanal.id === idKanalu);
    const tresc = pytanie.kontrolka.value.trim();

    konsultacjaWToku = true;
    przelicz();
    powiedz(`Pytanie w drodze do doradcy „${nazwaKanalu(wybrany, idKanalu)}"…`, true);

    const wynik = await zrodlo.zapytaj({
      kanalDoradcy: idKanalu,
      nazwaDoradcy: nazwaKanalu(wybrany, idKanalu),
      idEksperta: ekspert.id,
      nazwaEksperta: ekspert.displayName ?? ekspert.name,
      pytanie: tresc,
    });

    konsultacjaWToku = false;
    przelicz();

    if (!wynik.udany || wynik.wynik === undefined) {
      powiedz(opisOdmowy('Konsultacja doradcy', wynik.blad?.code, wynik.blad?.message), false);
      return;
    }
    const zapis = zapiszProwenancje(wynik.wynik);
    wykaz.prepend(wierszRady(zapis, wstawDoInstrukcji));
    powiedz(zapis.naglowek, true);
  }

  doradca.kontrolka.addEventListener('change', przelicz);
  pytanie.kontrolka.addEventListener('input', przelicz);
  zapytaj.addEventListener('click', () => void konsultuj());
  przelicz();

  return {
    element,

    ustaw(ekspert) {
      ekspertCzynny = ekspert;
      odswiezWykazKanalow();
      przelicz();
    },

    async wczytajKanaly() {
      const wynik = await zaplecze.kanaly();
      if (!wynik.udany || wynik.wynik === undefined) {
        kanaly = [];
        ustawPozycje(doradca.kontrolka, []);
        powiedz(
          opisOdmowy('Odczyt rejestru kanałów doradców', wynik.blad?.code, wynik.blad?.message),
          false,
        );
        przelicz();
        return;
      }
      kanaly = wynik.wynik.channels;
      odswiezWykazKanalow();
      przelicz();
    },
  };

  function odswiezWykazKanalow(): void {
    const dostepne = kanalyDoradcze(kanaly, ekspertCzynny?.channelId ?? '');
    ustawPozycje(
      doradca.kontrolka,
      dostepne.map((kanal) => ({
        wartosc: kanal.id,
        etykieta: `${kanal.name}${kanal.model === undefined ? '' : ` · ${kanal.model}`}`,
      })),
    );
  }
}

/** Nazwa kanału widziana przez Operatora; identyfikator, gdy kanał zniknął. */
function nazwaKanalu(kanal: Channel | undefined, idKanalu: string): string {
  return kanal === undefined ? idKanalu : kanal.name;
}

/**
 * Jeden wiersz wykazu konsultacji. Kolejność bloków jest celowa: etykieta
 * odróżniająca radę od odpowiedzi eksperta, droga konsultacji, pytanie, na
 * końcu treść rady — pochodzenie widać, zanim widać treść.
 */
function wierszRady(zapis: ZapisProwenancji, wstaw: (tresc: string) => void): HTMLElement {
  const etykieta = document.createElement('strong');
  etykieta.className = 'da-doradca__etykieta';
  etykieta.textContent = ETYKIETA_RADY;

  const naglowek = document.createElement('p');
  naglowek.className = 'da-doradca__naglowek';
  naglowek.textContent = zapis.naglowek;

  const droga = document.createElement('p');
  droga.className = 'dn-pole-opis da-doradca__droga';
  droga.textContent = zapis.droga;

  const zadane = document.createElement('blockquote');
  zadane.className = 'da-doradca__pytanie';
  zadane.textContent = zapis.rada.pytanie;

  const tresc = document.createElement('pre');
  tresc.className = 'da-doradca__tresc';
  tresc.textContent = zapis.rada.tresc === '' ? '(doradca nie oddał treści)' : zapis.rada.tresc;

  const przeniesienie = przycisk(
    'Przenieś do instrukcji z prowenancją',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  przeniesienie.addEventListener('click', () => wstaw(radaDoPrzeniesienia(zapis)));

  const element = document.createElement('li');
  element.className = 'da-doradca__wiersz';
  element.dataset['doradca'] = zapis.rada.kanalDoradcy;
  element.append(etykieta, naglowek, droga, zadane, tresc);

  // Tura urwana oddaje treść częściową. Zdanie o stanie wchodzi przed treść,
  // nie zamiast niej.
  const stan = zdanieOStanie(zapis.rada.stan);
  if (stan !== '') {
    const ostrzezenie = document.createElement('p');
    ostrzezenie.className = 'da-doradca__stan';
    ostrzezenie.textContent = stan;
    element.insertBefore(ostrzezenie, tresc);
  }

  element.append(przeniesienie);
  return element;
}

function tytul(nazwa: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'da-panel__tytul';
  element.textContent = nazwa;
  return element;
}
