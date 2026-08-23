import { ListenMode } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przyciskAkcji, wybor } from '../../modele/kontrolki-formularza';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Wybudzanie i nasłuch ciągły w Voice Console — cztery komendy rodziny mowy
 * wraz z dwoma jej zdarzeniami.
 *
 * ── Czym naprawdę jest „nasłuch ciągły rdzenia" ─────────────────────────────
 * Rdzeń stoi na serwerze i mikrofonu tej maszyny nie widzi. Nasłuch jest umową
 * między oknem a rdzeniem: okno nagrywa u siebie, wysyła odcinki
 * (`speech.audio.upload` wraz z `windowId`), a rdzeń ogłasza, co w nich
 * usłyszał — `speech.listen.partial` z tekstem oraz `speech.wake.detected`,
 * gdy padła fraza wybudzająca. Panel mówi to wprost, zamiast rysować mikrofon
 * sugerujący, że rdzeń słucha sam.
 *
 * ── Uczciwy stan wykonalności ───────────────────────────────────────────────
 * `speech.wake.get` oddaje `available: false` wraz z powodem, gdy wybudzenia
 * nie da się wykonać. Panel powtarza powód i nie stawia przycisku obiecującego
 * czynność, której nie ma czym wykonać.
 */
export interface PanelWybudzenia {
  element: HTMLElement;
  /** Odczyt nastawy wybudzania wraz ze stanem wykonalności. */
  wczytaj(): Promise<void>;
  /** Zamknięcie subskrypcji zdarzeń nasłuchu. */
  zamknij(): void;
  /** Uruchamia albo zatrzymuje nasłuch — droga z paska promptu. */
  przelaczNasluch(): void;
}

/** Tryby nasłuchu wskazywane wprost — wartości wyliczenia kontraktu. */
const TRYBY: ReadonlyArray<readonly [string, string]> = [
  [ListenMode.PushToTalk, 'Tryb: przytrzymanie przycisku'],
  [ListenMode.WakeWord, 'Tryb: fraza wybudzająca'],
  [ListenMode.Continuous, 'Tryb: nasłuch ciągły'],
];

export function utworzPanelWybudzenia(
  stan: StanAssistant,
  naTranskrypcje: (tresc: string) => void,
): PanelWybudzenia {
  const okno: StanOkna = utworzStanOkna();

  const fraza = pole('Fraza wybudzająca', 'np. Danaco, słuchaj');
  const tryb = wybor('Tryb nasłuchu', TRYBY);
  const prog = pole('Próg detekcji mowy (0–100; 0 = próg silnika)', '0');

  const zapisz = przyciskAkcji('Zapisz nastawę wybudzania', 'dn-btn dn-btn--sm dn-btn--zarys');
  zapisz.addEventListener('click', () => void zapiszNastawe());

  const przelacz = przyciskAkcji('Uruchom nasłuch ciągły', 'dn-btn dn-btn--sm dn-btn--zarys');
  przelacz.addEventListener('click', () => void przelaczNasluch());

  const slyszane = document.createElement('p');
  slyszane.className = 'dn-pole-opis';
  slyszane.dataset['nasluch'] = 'odcinek';
  slyszane.textContent =
    'Nasłuch ciągły: okno nagrywa u siebie i wysyła odcinki do rozpoznania. ' +
    'Rdzeń nie sięga po mikrofon tej maszyny.';

  okno.tresc.append(fraza, tryb, prog, zapisz, przelacz, slyszane);

  const tytul = document.createElement('h4');
  tytul.className = 'ma-panel__tytul';
  tytul.textContent = 'Wybudzanie i nasłuch ciągły';

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'wybudzenie';
  element.append(tytul, okno.element);

  let idNasluchu = '';

  // Zdarzenia nasłuchu wpinamy raz, przy budowie panelu: rozpoznany odcinek
  // wchodzi w pole polecenia tą samą drogą, którą wchodzi transkrypcja
  // nagrania, a wykrycie frazy wybudzającej jest komunikatem, nie treścią.
  const odsubskrybujOdcinek = stan.mowa.naOdcinek((tresc) => {
    if (tresc.transcript.trim() === '') return;
    slyszane.textContent = `Usłyszano: ${tresc.transcript}`;
    naTranskrypcje(tresc.transcript);
  });
  const odsubskrybujWybudzenie = stan.mowa.naWybudzenie((tresc) => {
    slyszane.textContent = `Fraza wybudzająca „${tresc.phrase}" padła — asystent słucha.`;
  });

  async function wczytaj(): Promise<void> {
    okno.ladowanie('Odczyt nastawy wybudzania…');
    const wynik = await stan.mowa.nastawaWybudzania();
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt nastawy wybudzania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const odpowiedz = wynik.wynik;
    fraza.value = odpowiedz.config.phrase;
    tryb.value = odpowiedz.config.mode;
    prog.value = String(odpowiedz.config.vadThreshold);
    if (!odpowiedz.available) {
      okno.puste(
        odpowiedz.reason !== undefined && odpowiedz.reason !== ''
          ? `Wybudzanie nie jest wykonalne: ${odpowiedz.reason}`
          : 'Wybudzanie nie jest wykonalne, a rdzeń nie podał powodu.',
      );
      return;
    }
    okno.gotowe();
  }

  async function zapiszNastawe(): Promise<void> {
    const liczba = Number.parseInt(prog.value.trim(), 10);
    const wynik = await stan.mowa.zapiszWybudzanie({
      fraza: fraza.value.trim(),
      tryb: tryb.value as (typeof ListenMode)[keyof typeof ListenMode],
      progDetekcji: Number.isNaN(liczba) ? 0 : liczba,
    });
    if (!wynik.udany) {
      okno.blad(opisOdmowy('Zapis nastawy wybudzania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    await wczytaj();
  }

  async function przelaczNasluch(): Promise<void> {
    const idOkna = stan.idOkna();
    if (idOkna === '') {
      okno.blad(
        'Rdzeń nie oddał okna modułu dla tej sesji — nasłuch biegnie na oknie, więc ' +
          'nie ma go dla czego uruchomić.',
      );
      return;
    }
    if (idNasluchu !== '') {
      const zatrzymanie = await stan.mowa.zatrzymajNasluch(idNasluchu, idOkna);
      if (!zatrzymanie.udany) {
        okno.blad(opisOdmowy('Zatrzymanie nasłuchu', zatrzymanie.blad?.code, zatrzymanie.blad?.message));
        return;
      }
      idNasluchu = '';
      przelacz.textContent = 'Uruchom nasłuch ciągły';
      slyszane.textContent = 'Nasłuch zatrzymany.';
      return;
    }

    const wynik = await stan.mowa.uruchomNasluch(idOkna, stan.idSesji());
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Uruchomienie nasłuchu', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    if (!wynik.wynik.listening) {
      // Nieuruchomiony nasłuch jest odpowiedzią, nie odmową: rdzeń mówi, czego
      // brakuje, a okno powtarza to zamiast przestawiać przycisk.
      okno.puste(
        wynik.wynik.reason !== undefined && wynik.wynik.reason !== ''
          ? `Nasłuch nie ruszył: ${wynik.wynik.reason}`
          : 'Nasłuch nie ruszył, a rdzeń nie podał powodu.',
      );
      return;
    }
    idNasluchu = wynik.wynik.listenerId;
    przelacz.textContent = 'Zatrzymaj nasłuch ciągły';
    slyszane.textContent = 'Nasłuch biegnie — wysyłaj odcinki mikrofonem powyżej.';
    okno.gotowe();
  }

  function zamknij(): void {
    odsubskrybujOdcinek();
    odsubskrybujWybudzenie();
  }

  return {
    element,
    wczytaj,
    zamknij,
    przelaczNasluch: () => {
      void przelaczNasluch();
    },
  };
}
