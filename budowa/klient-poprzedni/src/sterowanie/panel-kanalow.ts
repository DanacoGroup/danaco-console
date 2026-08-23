import './panel-kanalow.css';

import { utworzDymekObjasnienia } from '../komponenty/dymek';
import { utworzRameOkna, type RamaOkna } from '../komponenty/rama-okna';
import { przyciskAkcji } from '../modele/kontrolki-formularza-braki';
import type { Kanal } from '../protokol/kanal';
import {
  sprawdzKanal,
  stanPoswiadczeniaKanalu,
  usunKanal,
  utworzUzbrojenie,
  wskazWiersz,
  zalozKanal,
  zapiszKanal,
  type KontekstKanalow,
} from './czynnosci-kanalow';
import { utworzFormularzKanalu } from './formularz-kanalu';
import type { KomunikatZmiany } from './komunikat-zmiany';
import { utworzPasekKomunikatow } from './pasek-komunikatow';
import type { RejestrKanalow } from './rejestr-kanalow';
import { utworzWykazKanalow } from './wykaz-kanalow';
import { utworzZrodloKanalow } from './zrodlo-kanalow';

/**
 * Panel rejestru kanałów modelu — miejsce, w którym zakłada się, zmienia
 * i wykreśla wiersz rejestru (`channel.add`, `channel.update`,
 * `channel.remove`). Panel obsługuje wszystkie trzy czynności, bo instalacja ma
 * na starcie jeden kanał, a okno komunikacji bez kanału nie ma czym rozmawiać.
 *
 * Rejestr kanałów pozostaje bytem tylko do odczytu: panel nie ma do
 * `rejestr-kanalow.ts` żadnej drogi zapisu. Po udanym zapisie woła
 * `rejestr.odswiez()` — zamówienie `channel.list`, którego wynik rejestr
 * rozgłasza wszystkim czytelnikom (okno rozmowy, sterowania modelu głównego
 * i zapasowego, panel modeli, Roundtable). Jeden zapis odświeża wszystkich,
 * a żaden z czytelników nie wie o istnieniu tego panelu.
 *
 * Usunięcie idzie w dwóch krokach, przyciskiem w wariancie
 * `dn-btn--niebezpieczny` stojącym na końcu panelu akcji; w wierszach wykazu
 * przycisku usuwania nie ma. Zachowanie przy odmowie opisuje
 * `czynnosci-kanalow.ts`.
 */
export interface ZaleznosciPaneluKanalow {
  kanal: Kanal;
  /** Wspólny rejestr kanałów klienta. Panel go czyta i zamawia odświeżenie. */
  rejestrKanalow: RejestrKanalow;
}

export interface PanelKanalow {
  /** Sekcja osadzana w układzie. */
  element: HTMLElement;
  /** Zamawia wykaz z rdzenia i odrysowuje panel. */
  odswiez(): void;
}

const OBJASNIENIE =
  'Rejestr kanałów modelu jest sterowany danymi: nowy kanał to nowy wiersz rejestru, ' +
  'nie zmiana w kodzie. Ten sam wykaz czytają okno rozmowy, sterowania modelu i Roundtable.';

export function utworzPanelKanalow(zaleznosci: ZaleznosciPaneluKanalow): PanelKanalow {
  const { kanal, rejestrKanalow: rejestr } = zaleznosci;
  const rama = zlozRame();
  const komunikaty = utworzPasekKomunikatow();
  const usun = przyciskAkcji('Usuń wiersz rejestru', 'dn-btn dn-btn--niebezpieczny');

  const kontekst: KontekstKanalow = {
    zrodlo: utworzZrodloKanalow(kanal),
    formularz: utworzFormularzKanalu(),
    wykaz: utworzWykazKanalow(rejestr, (wskazany) => wskazWiersz(kontekst, wskazany)),
    rejestr,
    uzbrojenie: utworzUzbrojenie(usun),
    zglos: (komunikat: KomunikatZmiany) => zglos(rama, komunikaty.pokaz, komunikat),
  };

  function odswiez(): void {
    rejestr.odswiez();
    kontekst.wykaz.odrysuj();
  }

  usun.addEventListener('click', () => void usunKanal(kontekst));
  rama.akcje.append(...zlozAkcje(kontekst), usun);
  rama.narzedzia.append(przyciskNarzedzia('Odczytaj rejestr ponownie', odswiez));
  rama.cialo.append(kontekst.wykaz.element, kontekst.formularz.element, komunikaty.element);

  rejestr.naZmiane(() => kontekst.wykaz.odrysuj());
  kontekst.formularz.wypelnij(undefined);
  kontekst.wykaz.odrysuj();

  return { element: rama.element, odswiez };
}

/** Rama okna wraz z dymkiem objaśnienia w nagłówku. */
function zlozRame(): RamaOkna {
  return utworzRameOkna({
    tytul: 'Rejestr kanałów modelu',
    rola: 'zarządca',
    kod: 'channel-registry',
    przeznaczenie:
      'Zakładanie, zmiana i wykreślanie wierszy rejestru kanałów — channel.add, channel.update, channel.remove.',
    modul: 'Sterowanie okna',
    // Bez `przedrostek`: panel nie ma ani jednej reguły własnej dla gniazd ramy,
    // więc klasy modułu byłyby uchwytami bez reguł (rama-okna.ts, pole `przedrostek`).
    dodatkiNaglowka: [utworzDymekObjasnienia(OBJASNIENIE)],
  });
}

/**
 * Panel akcji: założenie, zapis zmiany, porzucenie wskazania. Usuwanie dokleja
 * wywołujący na końcu — kolejność jest tu treścią, nie układem: czynność
 * nieodwracalna nie stoi pierwsza i nie sąsiaduje z zapisem.
 */
function zlozAkcje(kontekst: KontekstKanalow): HTMLElement[] {
  return [
    przyciskNarzedzia('Założ nowy kanał', () => void zalozKanal(kontekst), 'dn-btn'),
    przyciskNarzedzia('Zapisz zmianę wskazanego', () => void zapiszKanal(kontekst)),
    przyciskNarzedzia('Sprawdź kanał', () => void sprawdzKanal(kontekst)),
    przyciskNarzedzia('Stan poświadczenia', () => void stanPoswiadczeniaKanalu(kontekst)),
    przyciskNarzedzia('Porzuć wskazanie', () => {
      kontekst.wykaz.wskaz('');
      wskazWiersz(kontekst, undefined);
    }),
  ];
}

/** Przycisk panelu; wariant domyślny jest zarysowy, bo wyróżniony ma być jeden. */
function przyciskNarzedzia(
  etykieta: string,
  przyKlikniecu: () => void,
  klasa = 'dn-btn dn-btn--zarys',
): HTMLElement {
  const przycisk = przyciskAkcji(etykieta, klasa);
  przycisk.addEventListener('click', przyKlikniecu);
  return przycisk;
}

/**
 * Komunikat panelu wraz z plakietką stanu w nagłówku. Plakietka gaśnie po
 * powodzeniu, więc odmowa nie zostaje na scenie po udanej próbie następnej.
 */
function zglos(
  rama: RamaOkna,
  pokaz: (komunikat: KomunikatZmiany) => void,
  komunikat: KomunikatZmiany,
): void {
  pokaz(komunikat);
  rama.ustawZnacznik(komunikat.udany ? '' : 'odmowa rdzenia', 'blad');
}
