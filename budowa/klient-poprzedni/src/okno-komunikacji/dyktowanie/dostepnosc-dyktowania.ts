import { Command } from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Kanal } from '../../protokol/kanal';
import { wywolaj } from '../../protokol/wywolanie';
import type { Dostarczenie } from './dostarczenie-nagrania';

/**
 * Dostępność dyktowania — koniunkcja dwóch warunków, nie jednego.
 *
 * Kontrakt odpowiada na `speech.availability.get` polem `available` i to pole
 * mówi o rdzeniu: czy stoi Python, czy stoi faster-whisper, czy jest model.
 * O kliencie nie mówi nic. Dyktowanie działa dopiero, gdy spełnione są oba
 * warunki:
 *
 *   (a) rdzeń ma silnik            → `speech.availability.get` → `available`
 *   (b) klient ma czym dostarczyć  → `dostarczenie.dostepne()`
 *
 * Sam warunek (a) wystarczyłby, gdyby silnik potrafił sięgnąć po nagranie sam.
 * Nie potrafi: `speech.transcribe` bierze ścieżkę pliku na maszynie silnika,
 * a nagranie z mikrofonu to bajty w pamięci karty (powód pełny w
 * `dostarczenie-nagrania.ts`). Pytanie wyłącznie o rdzeń dałoby więc mikrofon
 * zapalony na zielono, który po naciśnięciu nie zrobi nic, a brak ujawniłby się
 * dopiero po nagraniu.
 *
 * Dwa braki mają dwie naprawy i dwa zdania: brak silnika naprawia się
 * instalacją po stronie rdzenia, brak drogi dostarczenia — dopisaniem komendy
 * do kontraktu. Zlanie ich w jedno zdanie („dyktowanie niedostępne") wysłałoby
 * Operatora nie tam, gdzie leży przyczyna, więc gdy zawiodą oba, `powod` niesie
 * oba zdania osobno, każde ze swoją naprawą.
 *
 * Odmowa komendy jest odpowiedzią, nie wyjątkiem: gdy rdzeń odmówi albo nie zna
 * komendy (`speech.unknown`), wynik jest normalny — `dostepne=false` z powodem
 * nazywającym odmowę. Obietnica nie jest odrzucana i widok nie zakłada `try`.
 */

/** Odpowiedź na pytanie „czy Operator może teraz dyktować". */
export interface DostepnoscDyktowania {
  /** Koniunkcja obu warunków — dopiero ona zapala mikrofon. */
  dostepne: boolean;
  /** Wersja Pythona znaleziona przez rdzeń; pusta, gdy nie znaleziono. */
  python: string;
  /** Wersja silnika faster-whisper; pusta, gdy brak. */
  engine: string;
  /** Model ustawiony w katalogu ustawień; pusty, gdy nieustawiony. */
  model: string;
  /** Powód niedostępności — trzyczęściowy: co, dlaczego, czym Operator to zmieni. */
  powod: string;
}

export interface ZrodloDostepnosci {
  /** Pierwsze wywołanie pyta rdzeń, kolejne oddają ten sam odczyt. */
  odczytaj(): Promise<DostepnoscDyktowania>;
  /** Unieważnia odczyt — następne `odczytaj()` zapyta rdzeń ponownie. */
  odswiez(): void;
}

/** Zdanie trzyczęściowe o braku drogi dostarczenia — warunek (b) niespełniony. */
const ZDANIE_BRAKU_DOSTARCZENIA =
  'Dyktowanie jest wyłączone po stronie klienta. ' +
  'Nagranie z mikrofonu powstaje jako bajty w pamięci karty, a speech.transcribe przyjmuje audioRef, ' +
  'czyli ścieżkę pliku na maszynie silnika — komendy przenoszącej nagranie kontrakt dziś nie ma. ' +
  'Zmieni to dopisanie do kontraktu komendy przyjmującej nagranie w base64 wraz z uchwytem w rdzeniu.';

/** Zdanie trzyczęściowe o braku silnika, gdy rdzeń sam powodu nie podał. */
const ZDANIE_BRAKU_SILNIKA =
  'Silnik mowy nie jest gotowy po stronie rdzenia. ' +
  'Rdzeń odpowiedział na speech.availability.get, że transkrypcja nie jest tu i teraz wykonalna, ' +
  'a powodu nie nazwał — brakuje Pythona, pakietu faster-whisper albo modelu. ' +
  'Zmieni to instalacja pomocnika mowy na maszynie rdzenia i wskazanie modelu w katalogu ustawień.';

export function utworzZrodloDostepnosci(kanal: Kanal, dostarczenie: Dostarczenie): ZrodloDostepnosci {
  let odczyt: Promise<DostepnoscDyktowania> | null = null;

  async function pobierz(): Promise<DostepnoscDyktowania> {
    const odpowiedz = await wywolaj(kanal, Command.SpeechAvailabilityGet, {});

    // Warunek (b) sprawdzamy zawsze, także przy odmowie rdzenia: dwa braki mają
    // być widoczne oba naraz, a nie ujawniać się kolejno po każdej naprawie.
    const drogaJest = dostarczenie.dostepne();
    const powody: string[] = [];

    if (!odpowiedz.udany || odpowiedz.wynik === undefined) {
      // Odmowa komendy — w tym `speech.unknown`. Nie wiemy nic o silniku, więc
      // pola zostają puste, a powód nazywa odmowę, nie brak silnika.
      powody.push(opisOdmowyBledu('Odczyt dostępności silnika mowy (speech.availability.get)', odpowiedz.blad));
      if (!drogaJest) powody.push(ZDANIE_BRAKU_DOSTARCZENIA);
      return { dostepne: false, python: '', engine: '', model: '', powod: powody.join(' ') };
    }

    const tresc = odpowiedz.wynik;
    if (!tresc.available) {
      const wlasny = (tresc.reason ?? '').trim();
      powody.push(wlasny === '' ? ZDANIE_BRAKU_SILNIKA : wlasny);
    }
    if (!drogaJest) powody.push(ZDANIE_BRAKU_DOSTARCZENIA);

    return {
      // Koniunkcja: jedno „nie" po którejkolwiek stronie gasi mikrofon.
      dostepne: tresc.available && drogaJest,
      python: tresc.python ?? '',
      engine: tresc.engine ?? '',
      model: tresc.model ?? '',
      powod: powody.join(' '),
    };
  }

  return {
    odczytaj() {
      odczyt ??= pobierz();
      return odczyt;
    },
    odswiez() {
      // Samo unieważnienie, bez odpytania: `odswiez()` bywa wołane przy zmianie
      // ustawień mowy, a nie każde takie zdarzenie kończy się otwarciem okna
      // dyktowania. Pytanie rdzenia „na zapas" byłoby ruchem bez odbiorcy.
      odczyt = null;
    },
  };
}
