import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import {
  poleLogiczne,
  poleTekstowe,
  poleWielowierszowe,
  przycisk,
} from '../../modele/kontrolki-formularza';
import { KLASY_DYMKA } from './etykiety-assistant';

/**
 * Pasek narzędzi promptu modułu Assistant.
 *
 * Plik odpowiada wyłącznie za kontrolki wydania polecenia. Pasek niczego nie
 * wysyła — oddaje wartości i zgłasza zamiary, a wysyłkę prowadzi okno, żeby
 * droga do rdzenia była jedna.
 *
 * Mikrofon nie jest wygaszony, choć kontrakt nie ma przesyłu dźwięku:
 * naciśnięcie odpowiada zdaniem mówiącym, czego brakuje, i prowadzi ognisko do
 * pola transkrypcji. Rozpoznanie mowy jest warstwą wejścia, nie drugą drogą
 * rozmowy — po transkrypcji treść wchodzi tam, gdzie weszłaby wpisana ręcznie.
 */
export interface PasekPolecenia {
  element: HTMLElement;
  /** Treść polecenia po edycji Operatora. */
  transkrypcja(): string;
  /** Wstawia treść do pola — np. fragment odpowiedzi zamieniany w zadanie. */
  ustawTranskrypcje(tresc: string): void;
  /** Profil asystenta; pusty znaczy profil domyślny rdzenia. */
  profil(): string;
  /** Wstawia kod profilu wybrany w katalogu profili. */
  ustawProfil(kodProfilu: string): void;
  /** Czy odpowiedź ma zostać odczytana syntezą mowy. */
  czytaj(): boolean;
}

/** Zamiary paska obsługiwane przez okno. */
export interface UchwytyPaska {
  naWyslij(): void;
  naPrzerwij(): void;
  /** Nagranie polecenia mikrofonem karty — droga przez `speech.audio.upload`. */
  naMikrofon(): void;
  /** Nasłuch ciągły i fraza wybudzająca — droga przez `speech.listen.*`. */
  naWybudzenie(): void;
}

/**
 * Kontrolki wnoszone przez okno i osadzane w rzędzie kontrolek paska.
 *
 * Selektor profilu prowadzi własny odczyt rdzenia (`component.list`), więc nie
 * należy do paska, który niczego nie wywołuje. Pasek daje mu miejsce, a nie
 * buduje go sam.
 */
export interface KontrolkiOkna {
  katalogProfili: HTMLElement;
}

export function utworzPasekPolecenia(
  uchwyty: UchwytyPaska,
  kontrolkiOkna: KontrolkiOkna,
): PasekPolecenia {
  const polecenie = poleWielowierszowe(
    {
      etykieta: 'Transkrypcja polecenia',
      podpowiedz: 'powiedz albo wpisz, co asystent ma zrobić',
      opis:
        'Treść idzie do rdzenia polem transcript komendy assistant.voice.command. ' +
        'Poprawki nanieś przed wysłaniem — polecenie wychodzi w brzmieniu z pola.',
    },
    3,
  );

  const profil = poleTekstowe({
    etykieta: 'Profil asystenta',
    podpowiedz: 'kod profilu; puste znaczy domyślny',
    opis:
      'Kod jedzie do rdzenia polem profileId. Pole zostaje otwarte obok katalogu profili, ' +
      'bo profil zapisany w rdzeniu bez kafla w strefie 2 nie ma jak wejść do tamtego wykazu.',
  });

  const synteza = poleLogiczne({
    etykieta: 'Czytaj odpowiedź syntezą mowy',
    opis: 'Odpowiada polu speak komendy assistant.voice.command.',
  });
  synteza.kontrolka.checked = true;

  const element = document.createElement('div');
  element.className = 'ma-pasek';
  element.append(
    polecenie.element,
    zestawKontrolek(kontrolkiOkna.katalogProfili, profil.element, synteza.element),
    zestawPrzyciskow(uchwyty, polecenie.kontrolka),
  );

  return {
    element,
    transkrypcja: () => polecenie.kontrolka.value.trim(),
    ustawTranskrypcje(tresc) {
      polecenie.kontrolka.value = tresc;
      polecenie.kontrolka.focus();
    },
    profil: () => profil.kontrolka.value,
    ustawProfil(kodProfilu) {
      profil.kontrolka.value = kodProfilu;
    },
    czytaj: () => synteza.kontrolka.checked,
  };
}

/** Katalog profili, pole profilu i synteza wraz z dymkami [?]. */
function zestawKontrolek(
  katalogProfili: HTMLElement,
  profil: HTMLElement,
  synteza: HTMLElement,
): HTMLElement {
  const wiersz = document.createElement('div');
  wiersz.className = 'ma-pasek__kontrolki';
  wiersz.append(
    katalogProfili,
    profil,
    utworzDymekObjasnienia(
      'Profil asystenta wykonującego polecenie. Rdzeń przyjmuje go polem profileId; ' +
        'pominięty zostawia wybór profilu rdzeniowi.',
      KLASY_DYMKA,
    ),
    synteza,
    utworzDymekObjasnienia(
      'Wyciszenie zdejmuje pole speak z polecenia — rdzeń nie zamawia wtedy syntezy ' +
        'odpowiedzi. Sam potok pozostaje ten sam: STT → model główny → TTS.',
      KLASY_DYMKA,
    ),
  );
  return wiersz;
}

/**
 * Mikrofon, Wybudzenie, Wyślij i Przerwij — cztery przyciski, każdy z drogą.
 *
 * Wybudzenie stoi osobno od mikrofonu, bo to dwie różne czynności: mikrofon
 * nagrywa jedno polecenie i wysyła je do rozpoznania (`speech.audio.upload`
 * → `speech.transcribe`), a wybudzenie prowadzi nasłuch ciągły
 * (`speech.listen.start`) i frazę, na którą asystent reaguje
 * (`speech.wake.set`). Sklejone w jedną kontrolkę dałyby jeden przycisk
 * o dwóch znaczeniach.
 */
function zestawPrzyciskow(uchwyty: UchwytyPaska, pole: HTMLTextAreaElement): HTMLElement {
  const mikrofon = przycisk('Mikrofon — nagraj polecenie', 'dn-btn dn-btn--sm dn-btn--zarys');
  mikrofon.addEventListener('click', () => {
    uchwyty.naMikrofon();
    pole.focus();
  });

  const wybudzenie = przycisk('Wybudzenie głosem', 'dn-btn dn-btn--sm dn-btn--duch');
  wybudzenie.addEventListener('click', () => {
    uchwyty.naWybudzenie();
    pole.focus();
  });

  const wyslij = przycisk('Wyślij polecenie', 'dn-btn dn-btn--sm dn-btn--atrament');
  wyslij.addEventListener('click', () => uchwyty.naWyslij());

  const przerwij = przycisk('Przerwij', 'dn-btn dn-btn--sm dn-btn--niebezpieczny');
  przerwij.addEventListener('click', () => uchwyty.naPrzerwij());

  const wiersz = document.createElement('div');
  wiersz.className = 'ma-pasek__przyciski';
  wiersz.append(mikrofon, wybudzenie, wyslij, przerwij);
  return wiersz;
}
