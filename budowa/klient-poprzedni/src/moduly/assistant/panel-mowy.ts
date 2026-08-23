import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { opisOdmowy } from '../../komponenty/odmowa';
import { pole, przyciskAkcji, wybor } from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA, ODCZYTY } from './etykiety-assistant';
import { utworzStanOkna, type StanOkna } from './stan-okna';
import type { StanAssistant } from './stan-assistant';

/**
 * Rozpoznawanie mowy w Voice Console — jedyna droga głosu, którą kontrakt
 * naprawdę niesie.
 *
 * Kontrakt ma dwie komendy mowy i obie są tu wołane. `speech.availability.get`
 * mówi, czy silnik stoi na tej maszynie — brak silnika jest odpowiedzią, nie
 * awarią, więc okno nazywa powód zamiast milczeć. `speech.transcribe` zamienia
 * nagranie w tekst i wkłada go w pole polecenia, tam gdzie weszłaby treść
 * wpisana ręcznie.
 *
 * Granica jest w tym, czym jest `audioRef`: to ŚCIEŻKA PLIKU na maszynie
 * silnika, a nie bajty przesłane z karty. Nagranie z mikrofonu przeglądarki
 * istnieje wyłącznie jako `Blob` w pamięci i kontrakt nie ma komendy, która by
 * je przyjęła — ten sam brak nazywa `okno-komunikacji/dyktowanie/
 * dostarczenie-nagrania.ts` dla całej platformy. Panel obsługuje więc nagranie
 * już leżące na maszynie silnika, a mikrofon zostaje brakiem nazwanym wprost
 * (`pasek-polecenia.ts`).
 *
 * Pustą transkrypcję panel melduje jako wynik, nie jako niepowodzenie: pole
 * `processed` odróżnia „przetworzono, mowy nie było" od „nie przetworzono",
 * a cisza w nagraniu jest prawidłowym wynikiem pomiaru.
 *
 * Rozpoznany tekst wchodzi w pole polecenia, a nie jedzie do rdzenia sam.
 * `assistant.voice.command` przyjmuje wprawdzie `audioRef` i rozpoznałby
 * nagranie po swojemu, ale wtedy Operator zobaczyłby transkrypcję dopiero
 * w odpowiedzi — po wykonaniu polecenia. Przepływ modułu wymaga kolejności
 * odwrotnej: rozpoznanie, korekta niepewnej frazy, dopiero wysyłka. Droga do
 * rdzenia zostaje więc jedna, przez pole `transcript`.
 */
export interface PanelMowy {
  element: HTMLElement;
  /** Odczyt dostępności silnika mowy. */
  wczytaj(): Promise<void>;
  /** Rozpoczyna albo kończy nagranie mikrofonem — droga z paska promptu. */
  przelaczMikrofon(): void;
}

/** Języki wskazywane wprost; pusty znaczy rozpoznanie automatyczne. */
const JEZYKI: ReadonlyArray<readonly [string, string]> = [
  ['', 'Język nagrania: rozpoznaj automatycznie'],
  ['pl', 'Język nagrania: polski'],
  ['en', 'Język nagrania: angielski'],
];

export function utworzPanelMowy(
  stan: StanAssistant,
  naTranskrypcje: (tresc: string) => void,
): PanelMowy {
  const okno: StanOkna = utworzStanOkna();

  const sciezka = pole(
    'Ścieżka nagrania na maszynie silnika',
    'ścieżka pliku widziana przez maszynę silnika mowy',
  );
  const jezyk = wybor('Język nagrania', JEZYKI);

  const rozpoznaj = przyciskAkcji('Rozpoznaj nagranie', 'dn-btn dn-btn--sm dn-btn--zarys');
  rozpoznaj.addEventListener('click', () => void transkrybuj());

  // Mikrofon karty: nagranie idzie do rdzenia komendą `speech.audio.upload`,
  // która oddaje odnośnik przyjmowany przez `speech.transcribe`. To jest
  // ogniwo, którego przez długi czas brakowało — bez niego nagranie z
  // mikrofonu nie miało czym dojechać do silnika.
  const mikrofon = przyciskAkcji('Nagraj z mikrofonu', 'dn-btn dn-btn--sm dn-btn--zarys');
  mikrofon.addEventListener('click', () => void przelaczNagrywanie());

  const kontrolki = document.createElement('div');
  kontrolki.className = 'ma-mowa__kontrolki';
  kontrolki.append(sciezka, jezyk, rozpoznaj, mikrofon);

  // Wskaźnika pewności rozpoznania panel nie stawia: ani `speech.transcribe`,
  // ani `assistant.voice.command` nie oddają miary pewności — odpowiedź niesie
  // długość nagrania, liczbę znaków, model i język. Plakietka „niska pewność"
  // byłaby tu wartością wziętą znikąd, więc zamiast niej stoi zdanie o tym, co
  // rdzeń o rozpoznaniu naprawdę powiedział.
  const pewnosc = document.createElement('p');
  pewnosc.className = 'dn-pole-opis';
  pewnosc.textContent =
    'Rdzeń nie oddaje miary pewności rozpoznania — po transkrypcji przeczytaj tekst ' +
    'w polu polecenia i popraw go przed wysłaniem.';

  okno.tresc.append(kontrolki, pewnosc);

  const tytul = document.createElement('h4');
  tytul.className = 'ma-panel__tytul';
  tytul.append(
    'Rozpoznawanie mowy',
    utworzDymekObjasnienia(
      'Silnik lokalny rdzenia (speech.transcribe). Pole audioRef kontraktu to ścieżka ' +
        'pliku na maszynie silnika — dźwięk nie opuszcza tej maszyny, a nagranie z ' +
        'mikrofonu karty nie ma dziś czym do niej dojechać.',
      KLASY_DYMKA,
    ),
  );

  const element = document.createElement('div');
  element.className = 'ma-panel';
  element.dataset['panel'] = 'mowa';
  element.append(tytul, okno.element);

  /** Zdanie o silniku: co jest, czego nie ma i czym Operator to zmieni. */
  async function wczytaj(): Promise<void> {
    okno.ladowanie(ODCZYTY.silnikMowy);
    const wynik = await stan.zaplecze.silnikMowy();
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Odczyt dostępności silnika mowy', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const stanSilnika = wynik.wynik;
    if (!stanSilnika.available) {
      // Niedostępność silnika nie jest odmową rdzenia — rdzeń odpowiedział.
      // Stan pusty niesie powód, który rdzeń podał, wraz z tym, co go zmieni.
      okno.puste(
        stanSilnika.reason !== undefined && stanSilnika.reason !== ''
          ? `Silnik mowy nie jest gotów: ${stanSilnika.reason}`
          : 'Silnik mowy nie jest gotów, a rdzeń nie podał powodu. ' +
              'Polecenie wydasz polem transkrypcji powyżej.',
      );
      return;
    }
    okno.gotowe();
  }

  /** Nagrywanie w toku; pusty znaczy, że mikrofon stoi. */
  let nagrywanie: MediaRecorder | undefined;

  /**
   * Nagranie z mikrofonu karty i jego droga do rdzenia.
   *
   * Bajty idą do `speech.audio.upload`, a oddany odnośnik wchodzi w pole
   * ścieżki i od razu jedzie do rozpoznania. Dźwięk nie opuszcza maszyny
   * rdzenia: droga prowadzi tam i z powrotem, nigdzie indziej.
   *
   * Brak dostępu do mikrofonu nie jest awarią rdzenia i okno tak go nazywa —
   * przeglądarka bywa bez zgody, bez urządzenia albo w kontekście bez
   * `mediaDevices`, i każdy z tych powodów Operator naprawia u siebie.
   */
  async function przelaczNagrywanie(): Promise<void> {
    if (nagrywanie !== undefined) {
      nagrywanie.stop();
      return;
    }
    if (navigator.mediaDevices === undefined) {
      okno.blad(
        'Ta karta nie ma dostępu do urządzeń nagrywających (brak navigator.mediaDevices). ' +
          'Nagranie wskażesz ścieżką powyżej albo otworzysz Danaco Console w programie okiennym.',
      );
      return;
    }
    let strumien: MediaStream;
    try {
      strumien = await navigator.mediaDevices.getUserMedia({ audio: true });
    } catch (powod) {
      okno.blad(
        `Mikrofon nie został udostępniony: ${String(powod)}. ` +
          'Zgody udziela przeglądarka, nie rdzeń.',
      );
      return;
    }

    const odcinki: Blob[] = [];
    const rejestrator = new MediaRecorder(strumien);
    nagrywanie = rejestrator;
    mikrofon.textContent = 'Zakończ nagranie';

    rejestrator.addEventListener('dataavailable', (zdarzenie) => {
      if (zdarzenie.data.size > 0) odcinki.push(zdarzenie.data);
    });
    rejestrator.addEventListener('stop', () => {
      nagrywanie = undefined;
      mikrofon.textContent = 'Nagraj z mikrofonu';
      for (const sciezkaWejscia of strumien.getTracks()) sciezkaWejscia.stop();
      void przeslij(new Blob(odcinki, { type: rejestrator.mimeType }));
    });
    rejestrator.start();
    okno.ladowanie('Nagrywanie… naciśnij ponownie, żeby zakończyć.');
  }

  /** Wysyła nagranie do rdzenia i od razu prosi o jego rozpoznanie. */
  async function przeslij(nagranie: Blob): Promise<void> {
    if (nagranie.size === 0) {
      okno.puste('Nagranie jest puste — mikrofon nie oddał ani jednego bajtu.');
      return;
    }
    okno.ladowanie('Przesyłanie nagrania do rdzenia…');
    const bajty = new Uint8Array(await nagranie.arrayBuffer());
    let zapis = '';
    for (const bajt of bajty) zapis += String.fromCharCode(bajt);

    const wynik = await stan.mowa.przeslijNagranie({
      base64: btoa(zapis),
      typTresci: nagranie.type === '' ? 'audio/webm' : nagranie.type,
      idOkna: stan.idOkna(),
      idSesji: stan.idSesji(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Przesłanie nagrania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    sciezka.value = wynik.wynik.audioRef;
    await transkrybuj();
  }

  async function transkrybuj(): Promise<void> {
    const wskazana = sciezka.value.trim();
    if (wskazana === '') {
      okno.blad(
        'Wskaż ścieżkę nagrania na maszynie silnika. Kontrakt przyjmuje w polu audioRef ' +
          'ścieżkę pliku, nie sam dźwięk, więc bez niej rdzeń nie ma czego rozpoznać.',
      );
      return;
    }
    okno.ladowanie(ODCZYTY.transkrypcja);
    const zadanie = { audioRef: wskazana };
    const wynik = await stan.zaplecze.transkrypcja(
      jezyk.value === '' ? zadanie : { ...zadanie, language: jezyk.value },
    );
    if (!wynik.udany || wynik.wynik === undefined) {
      okno.blad(opisOdmowy('Rozpoznanie nagrania', wynik.blad?.code, wynik.blad?.message));
      return;
    }
    const rozpoznane = wynik.wynik;
    if (rozpoznane.transcript === '') {
      // Cisza jest wynikiem, nie usterką: rdzeń przetworzył nagranie i nie
      // znalazł w nim mowy. Podstawienie pustego napisu w pole polecenia
      // skasowałoby treść, którą Operator zdążył tam wpisać.
      okno.puste(
        `Silnik przetworzył nagranie (${String(rozpoznane.durationMs)} ms, model ` +
          `${rozpoznane.model}) i nie rozpoznał w nim mowy. Pole polecenia zostaje bez zmian.`,
      );
      return;
    }
    naTranskrypcje(rozpoznane.transcript);
    okno.gotowe();
  }

  return {
    element,
    wczytaj,
    przelaczMikrofon: () => {
      void przelaczNagrywanie();
    },
  };
}
