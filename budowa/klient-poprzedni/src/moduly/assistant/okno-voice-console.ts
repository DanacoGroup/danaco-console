import type { FazaOkna } from '../../komponenty/faza-okna';
import { utworzNaglowekOkna } from '../../komponenty/naglowek-okna';
import { utworzWierszOdpowiedzi, type WierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { BRAKI } from './braki-kontraktu';
import { PUSTE } from './etykiety-assistant';
import { utworzPanelMowy, type PanelMowy } from './panel-mowy';
import { utworzPanelOdpowiedzi, type PanelOdpowiedzi } from './panel-odpowiedzi';
import { utworzPanelWybudzenia, type PanelWybudzenia } from './panel-wybudzenia';
import { utworzPasekPolecenia, type PasekPolecenia } from './pasek-polecenia';
import { utworzSelektorProfilu, type SelektorProfilu } from './selektor-profilu';
import { utworzSiatkeAkcji, type SiatkaAkcji } from './siatka-akcji';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';
import { utworzWysylkePolecenia, type WysylkaPolecenia } from './wysylka-polecenia';

/**
 * Kod okna w katalogu okien operacyjnych rdzenia, pole `okno_operacyjne.kod`,
 * po którym rdzeń rozpoznaje to okno w zleceniach i w zapisie czynności.
 */
export const KOD_OKNA = 'voice-console';

/**
 * Voice Console to okno wiodące modułu Assistant: wydanie polecenia głosowego,
 * odsłuch odpowiedzi syntezowanej oraz przerwanie nagrania. Plik odpowiada
 * wyłącznie za skład okna i dostaje treść gotową z osobnych składników.
 */
export interface OknoVoiceConsole {
  element: HTMLElement;
  /** Odczyt katalogu szybkich akcji, katalogu profili i stanu silnika mowy. */
  wczytaj(): Promise<void>;
  /** Wstawia treść w pole polecenia — droga wejścia z okien zarządczych. */
  ustawPolecenie(tresc: string): void;
  /** Zamyka subskrypcje zdarzeń nasłuchu ciągłego. */
  zamknij(): void;
  /** Nanosi stan modułu — zlecenie zamknięte znika z toru przerwania. */
  odswiez(): void;
}

export function utworzOknoVoiceConsole(stan: StanAssistant): OknoVoiceConsole {
  const okno: StanOkna = utworzStanOkna();
  const odpowiedz: WierszOdpowiedzi = utworzWierszOdpowiedzi();

  const panel: PanelOdpowiedzi = utworzPanelOdpowiedzi(
    (tresc) => {
      pasek.ustawTranskrypcje(tresc);
      odpowiedz.pokaz(
        'Fragment przeniesiony do pola polecenia. Wyślij je, aby założyć zlecenie ' +
          'widoczne w Actions Monitor.',
        true,
      );
    },
    // Odsłuch: `speech.audio.fetch` oddaje wyłącznie to, co rdzeń sam wystawił.
    (odnosnik) => {
      void (async () => {
        const wynik = await stan.mowa.pobierzNagranie(odnosnik);
        if (!wynik.udany || wynik.wynik === undefined) {
          odpowiedz.pokaz(
            `Nie udało się pobrać odpowiedzi syntezowanej: ${wynik.blad?.message ?? 'brak powodu'}`,
            false,
          );
          return;
        }
        const dzwiek = new Audio(`data:${wynik.wynik.contentType};base64,${wynik.wynik.audio}`);
        void dzwiek.play();
      })();
    },
  );

  const wysylka: WysylkaPolecenia = utworzWysylkePolecenia(stan, {
    powiedz: (tresc, powodzenie) => odpowiedz.pokaz(tresc, powodzenie),
    faza: (faza, opis) => naniesFaze(okno, faza, opis),
    odpowiedz: (tresc) => (tresc === null ? panel.wyczysc() : panel.pokaz(tresc)),
  });

  const profile: SelektorProfilu = utworzSelektorProfilu(stan, (kodProfilu) =>
    pasek.ustawProfil(kodProfilu),
  );

  const pasek: PasekPolecenia = utworzPasekPolecenia(
    {
      naWyslij: () =>
        void wysylka.wyslij({
          transkrypcja: pasek.transkrypcja(),
          profil: pasek.profil(),
          czytaj: pasek.czytaj(),
        }),
      naPrzerwij: () => void wysylka.przerwij(),
      // Uchwyty są leniwe, ponieważ panele powstają po pasku przycisków.
      naMikrofon: () => mowa.przelaczMikrofon(),
      naWybudzenie: () => wybudzenie.przelaczNasluch(),
    },
    { katalogProfili: profile.element },
  );

  const mowa: PanelMowy = utworzPanelMowy(stan, (tresc) => pasek.ustawTranskrypcje(tresc));
  const wybudzenie: PanelWybudzenia = utworzPanelWybudzenia(stan, (tresc) =>
    pasek.ustawTranskrypcje(tresc),
  );
  const siatka: SiatkaAkcji = utworzSiatkeAkcji(stan, (tresc) => pasek.ustawTranskrypcje(tresc));

  okno.tresc.append(
    pasek.element,
    odpowiedz.element,
    panel.element,
    mowa.element,
    wybudzenie.element,
    siatka.element,
  );
  okno.puste(PUSTE.konsolaSpoczynek);

  const element = document.createElement('section');
  element.className = 'ma-okno ma-okno--wiodace';
  element.dataset['okno'] = KOD_OKNA;
  element.append(
    utworzNaglowekOkna({
      tytul: 'Voice Console',
      rola: 'okno wiodące · polecenia głosowe i tekstowe w jednym zapisie',
    }),
    okno.element,
  );

  function odswiez(): void {
    wysylka.zsynchronizuj();
    naniesSpoczynek(okno, stan);
  }

  odswiez();
  return {
    element,
    // Trzy odczyty idą równolegle i żaden nie warunkuje pozostałych.
    wczytaj: async () => {
      await Promise.all([
        siatka.wczytaj(),
        profile.wczytaj(),
        mowa.wczytaj(),
        wybudzenie.wczytaj(),
      ]);
    },
    ustawPolecenie: (tresc) => pasek.ustawTranskrypcje(tresc),
    zamknij: () => wybudzenie.zamknij(),
    odswiez,
  };
}

/**
 * Wpuszcza fazę nazwaną przez wysyłkę polecenia w stan okna, korzystając
 * z jednego słownika faz wspólnego dla obu dróg wejścia.
 */
function naniesFaze(okno: StanOkna, faza: FazaOkna, opis: string): void {
  if (faza === 'ladowanie') okno.ladowanie(opis);
  else if (faza === 'blad') okno.blad(opis);
  else if (faza === 'puste') okno.puste(opis);
  else okno.gotowe();
}

/**
 * Faza okna w spoczynku, gdy nie jedzie żadne polecenie; rusza wyłącznie
 * ze stanu pustego, więc komunikatu wysyłki nie zdejmuje.
 */
function naniesSpoczynek(okno: StanOkna, stan: StanAssistant): void {
  if (okno.faza() !== 'puste') return;
  if (!stan.pytanoOOkno()) {
    okno.puste(PUSTE.konsolaSpoczynek);
    return;
  }
  if (stan.idOkna() === '') {
    okno.blad(BRAKI.brakOkna);
    return;
  }
  okno.puste(PUSTE.konsola);
}
