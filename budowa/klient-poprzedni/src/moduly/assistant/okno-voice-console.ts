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

/** Kod okna w katalogu rdzenia (`okno_operacyjne.kod`). */
export const KOD_OKNA = 'voice-console';

/**
 * Voice Console — okno wiodące modułu Assistant.
 *
 * Trzy funkcje operatora: wydanie polecenia głosowego, odsłuch odpowiedzi
 * syntezowanej, przerwanie nagrania. Pasek narzędzi promptu niesie mikrofon,
 * pole poleceń, Wyślij, selektor profilu asystenta, przełącznik syntezy
 * i przeniesienie fragmentu odpowiedzi w zadanie Actions Monitora; siatka
 * szybkich akcji stoi pod nimi.
 *
 * Plik odpowiada wyłącznie za skład okna. Rozmowa z rdzeniem mieszka
 * w `wysylka-polecenia.ts`, kontrolki w `pasek-polecenia.ts`, katalog akcji
 * w `siatka-akcji.ts`, katalog profili w `selektor-profilu.ts`, a rozpoznawanie
 * mowy w `panel-mowy.ts` — okno nie buduje treści, dostaje ją gotową.
 *
 * Droga głosu prowadzi przez `panel-mowy.ts` i kończy się w tym samym polu
 * transkrypcji, w które Operator wpisuje polecenie ręcznie. Rozpoznanie mowy
 * jest warstwą wejścia, nie drugą drogą rozmowy: rdzeń dostaje zawsze jedno
 * `assistant.voice.command` z polem `transcript`.
 *
 * Historia poleceń głosowych i tekstowych jest jedna i mieszka w Activity
 * Feed: `assistant.activity.list` zwraca jeden chronologiczny zapis obu dróg
 * (pole `origin` zlecenia je rozróżnia). Voice Console odświeża ten zapis po
 * każdym poleceniu, zamiast prowadzić drugą kopię.
 *
 * Fazę okna nazywają dwa źródła, które nie mogą się pobić:
 *
 *   · wysyłka — gdy polecenie jedzie, wróciło albo zostało odrzucone
 *     (`wysylka-polecenia.ts` mówi fazę wprost);
 *   · spoczynek — gdy nie jedzie nic; liczy go `naniesSpoczynek` poniżej
 *     i rusza wyłącznie z fazy `puste`, więc komunikatu wysyłki nie zdejmie.
 *
 * Stan pusty nie zastępuje paska promptu, tylko stoi nad nim (`stan-okna.ts`),
 * więc Operator czyta go mając pole transkrypcji i „Wyślij polecenie"
 * na wyciągnięcie ręki.
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
    // Odsłuch odpowiedzi syntezowanej: bajty wracają do karty tą samą drogą,
    // którą wychodzą nagrania — `speech.audio.fetch` oddaje wyłącznie to, co
    // rdzeń sam wystawił.
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
      // Mikrofon i wybudzenie prowadzą do paneli tego samego okna. Uchwyty są
      // leniwe, bo panele powstają po pasku: zestaw przycisków buduje się raz,
      // a woła dopiero po naciśnięciu.
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
    // Trzy odczyty idą równolegle: katalog akcji, katalog profili i stan silnika
    // mowy dotyczą trzech różnych komend i żaden nie warunkuje pozostałych.
    // Każdy nazywa swoje niepowodzenie w swoim miejscu, więc odmowa jednego nie
    // zabiera treści dwóm pozostałym.
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

/** Wpuszcza fazę nazwaną przez wysyłkę w stan okna; słownik jest jeden. */
function naniesFaze(okno: StanOkna, faza: FazaOkna, opis: string): void {
  if (faza === 'ladowanie') okno.ladowanie(opis);
  else if (faza === 'blad') okno.blad(opis);
  else if (faza === 'puste') okno.puste(opis);
  else okno.gotowe();
}

/**
 * Faza okna, gdy nie jedzie żadne polecenie.
 *
 * Rusza wyłącznie ze stanu pustego. Gdy wysyłka postawiła okno w ładowaniu,
 * odmowie albo gotowości, jej komunikat należy do niej i zostaje — inaczej
 * zdarzenie `assistant.action.changed` z cudzego zlecenia zmiatałoby odmowę
 * sprzed sekundy.
 *
 * Brak okna modułu jest błędem, nie pustką: bez `windowId` komenda
 * `assistant.voice.command` nie ma dokąd pojechać i nie zmieni tego żadna treść
 * w polu.
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
