import {
  MobileProcessControl,
  ProgressStatus,
  type MobileProcess,
} from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import {
  CZYNNOSCI_STEROWANIA,
  nazwaCzynnosci,
  nazwaStanu,
  opisProcesu,
  STANY_PROCESU,
  utworzZrodloProcesowMobilnych,
  type CzynnoscSterowania,
  type ZrodloProcesowMobilnych,
} from './procesy-mobilne';

/**
 * Przegląd zadań i procesów — ekran wykazu procesów wraz ze sterowaniem nimi
 * (`mobile.process.list`, `mobile.process.control`).
 *
 * Ekran nazywa się tak, jak nazywa go opracowanie funkcji globalnej Mobile
 * (rozdz. 4.6), i niesie dokładnie te dwa elementy, które ono wymienia: listę
 * procesów z filtrem stanu oraz zestaw czynności przy pozycji — uruchom
 * ponownie, zatrzymaj, wstrzymaj, wznów, zatwierdź, modyfikuj.
 *
 * Czynność nieodwracalna mówi to PRZED wykonaniem i wymaga drugiego dotknięcia.
 * Dwa dotknięcia to nie utrudnienie: telefon nosi się w kieszeni, a zatrzymanie
 * pracy, która biegnie bez Operatora, jest jedyną czynnością tego ekranu, której
 * nie da się cofnąć niczym. Uzbrojenie gaśnie po dotknięciu czegokolwiek
 * innego — inaczej przycisk zostałby uzbrojony na godziny.
 *
 * Filtr stanu zawęża po stronie RDZENIA, polem `status` żądania, a nie po
 * stronie widoku: wykaz przefiltrowany w oknie kłamałby o liczbie procesów,
 * których rdzeń nie przysłał.
 *
 * Urządzenie mobilne nazywa się kartą sesji kanału — tak samo jak w kafelku
 * stanu platformy, żeby rdzeń widział jedno urządzenie, a nie dwa.
 */
export interface EkranProcesow {
  element: HTMLElement;
  /** Pyta rdzeń i przerysowuje wykaz. */
  odswiez(): Promise<void>;
  /** Procesy z ostatniego odczytu — do sprawdzianów widoku. */
  procesy(): readonly MobileProcess[];
}

export interface OpisEkranuProcesow {
  kanal: Kanal;
  /** Źródło rodziny; wstrzykiwane dla sprawdzianu, domyślnie z kanału. */
  zrodlo?: ZrodloProcesowMobilnych;
}

export function utworzEkranProcesow(opis: OpisEkranuProcesow): EkranProcesow {
  const zrodlo = opis.zrodlo ?? utworzZrodloProcesowMobilnych(opis.kanal);
  let procesy: readonly MobileProcess[] = [];
  /** Czynność uzbrojona wraz z procesem, którego dotyczy; `null` znaczy brak. */
  let uzbrojona: { proces: string; czynnosc: MobileProcessControl } | null = null;

  const stan = document.createElement('p');
  stan.className = 'mb-ekran__stan';
  stan.setAttribute('aria-live', 'polite');

  const naglowek = document.createElement('h3');
  naglowek.className = 'mb-ekran__naglowek';
  naglowek.textContent = 'Przegląd zadań i procesów';

  const filtr = document.createElement('select');
  filtr.className = 'dn-pole-kontrolka mb-cel';
  filtr.setAttribute('aria-label', 'Stan zawężający wykaz procesów');
  for (const pozycja of [{ kod: '', nazwa: 'Wszystkie stany' }, ...STANY_PROCESU]) {
    const opcja = document.createElement('option');
    opcja.value = pozycja.kod;
    opcja.textContent = pozycja.nazwa;
    filtr.append(opcja);
  }
  filtr.addEventListener('change', () => void odswiez());

  const wykaz = document.createElement('div');
  wykaz.className = 'mb-ekran__wykaz';

  const element = document.createElement('section');
  element.className = 'mb-ekran';
  element.dataset['ekran'] = 'procesy';
  element.append(naglowek, filtr, stan, wykaz);

  /** Urządzenie mobilne: karta sesji kanału; pusta nie jedzie w żądaniu. */
  function urzadzenie(): { deviceId?: string } {
    const idSesji = opis.kanal.sesja().id();
    return idSesji === '' ? {} : { deviceId: idSesji };
  }

  /** Stan wybrany w filtrze; pusty znaczy „wszystkie stany". */
  function wybranyStan(): { status?: ProgressStatus } {
    const pozycja = STANY_PROCESU.find((wpis) => wpis.kod === filtr.value);
    return pozycja === undefined ? {} : { status: pozycja.kod };
  }

  function ustawStan(zdanie: string, rodzaj: 'ladowanie' | 'blad' | 'pusto' | 'tresc'): void {
    stan.dataset['stan'] = rodzaj;
    stan.textContent = zdanie;
  }

  async function odswiez(): Promise<void> {
    uzbrojona = null;
    ustawStan('Pytam rdzeń o procesy platformy (mobile.process.list)…', 'ladowanie');
    const wynik = await zrodlo.wykaz({ ...urzadzenie(), ...wybranyStan() });

    if (!wynik.udany || wynik.wynik === undefined) {
      procesy = [];
      wykaz.replaceChildren();
      ustawStan(opisOdmowyBledu('Wykaz procesów platformy', wynik.blad), 'blad');
      return;
    }

    procesy = wynik.wynik;
    if (procesy.length === 0) {
      wykaz.replaceChildren();
      // Pusty wykaz jest odpowiedzią, nie awarią, i zdanie rozróżnia dwa
      // powody pustki: brak procesów wcale i brak procesów w wybranym stanie.
      ustawStan(
        filtr.value === ''
          ? 'Rdzeń nie prowadzi ani jednego procesu — nie ma czym sterować. To pusty wykaz, ' +
              'nie nieudany odczyt.'
          : `W stanie „${nazwaStanu(filtr.value as ProgressStatus)}" nie ma ani jednego procesu. ` +
              'Zdejmij zawężenie, żeby zobaczyć pozostałe.',
        'pusto',
      );
      return;
    }

    ustawStan(`Procesów w wykazie: ${procesy.length} (mobile.process.list).`, 'tresc');
    przerysuj();
  }

  function przerysuj(): void {
    wykaz.replaceChildren(...procesy.map(karta));
  }

  function karta(proces: MobileProcess): HTMLElement {
    const element = document.createElement('article');
    element.className = 'dn-karta mb-karta';
    element.dataset['proces'] = proces.id;
    if (proces.status === ProgressStatus.Failed) element.dataset['rodzaj'] = 'proces-bledny';

    const tytul = document.createElement('p');
    tytul.className = 'mb-karta__tytul';
    tytul.textContent = proces.label;

    const kontekst = document.createElement('p');
    kontekst.className = 'mb-karta__kontekst';
    kontekst.textContent = opisProcesu(proces);

    const ostrzezenie = document.createElement('p');
    ostrzezenie.className = 'mb-arkusz__wyjasnienie';
    ostrzezenie.hidden = true;
    ostrzezenie.setAttribute('role', 'alert');

    const czynnosci = document.createElement('div');
    czynnosci.className = 'mb-arkusz__sekcja';
    for (const czynnosc of CZYNNOSCI_STEROWANIA) {
      czynnosci.append(przyciskCzynnosci(proces, czynnosc, ostrzezenie));
    }

    element.append(tytul, kontekst, czynnosci, ostrzezenie);
    return element;
  }

  /**
   * Przycisk jednej czynności.
   *
   * Czynność odwracalna idzie od razu. Czynność nieodwracalna najpierw mówi, co
   * Operator traci, a wykonuje się dopiero po drugim dotknięciu tego samego
   * przycisku — i napis przycisku zmienia się na czas uzbrojenia, żeby nie dało
   * się go pomylić z pierwszym dotknięciem.
   */
  function przyciskCzynnosci(
    proces: MobileProcess,
    czynnosc: CzynnoscSterowania,
    ostrzezenie: HTMLElement,
  ): HTMLButtonElement {
    const przycisk = document.createElement('button');
    przycisk.type = 'button';
    przycisk.className = 'dn-btn dn-btn--sm dn-btn--zarys mb-cel mb-arkusz__droga';
    przycisk.dataset['czynnosc'] = czynnosc.kod;
    przycisk.textContent = czynnosc.nazwa;

    przycisk.addEventListener('click', () => {
      if (!czynnosc.nieodwracalna) {
        void wykonaj(proces, czynnosc, przycisk, ostrzezenie);
        return;
      }
      const juzUzbrojona =
        uzbrojona !== null &&
        uzbrojona.proces === proces.id &&
        uzbrojona.czynnosc === czynnosc.kod;
      if (juzUzbrojona) {
        void wykonaj(proces, czynnosc, przycisk, ostrzezenie);
        return;
      }
      uzbrojona = { proces: proces.id, czynnosc: czynnosc.kod };
      przycisk.textContent = `${czynnosc.nazwa} — dotknij ponownie, aby wykonać`;
      przycisk.dataset['uzbrojona'] = 'tak';
      ostrzezenie.hidden = false;
      ostrzezenie.textContent =
        `Czynność nieodwracalna: „${czynnosc.nazwa}" na procesie „${proces.label}". ` +
        (czynnosc.ostrzezenie ?? '');
    });

    return przycisk;
  }

  async function wykonaj(
    proces: MobileProcess,
    czynnosc: CzynnoscSterowania,
    przycisk: HTMLButtonElement,
    ostrzezenie: HTMLElement,
  ): Promise<void> {
    uzbrojona = null;
    przycisk.textContent = czynnosc.nazwa;
    delete przycisk.dataset['uzbrojona'];
    przycisk.disabled = true;
    ostrzezenie.hidden = false;
    ostrzezenie.textContent = `Rdzeń wykonuje „${czynnosc.nazwa}" na procesie „${proces.label}"…`;

    const wynik = await zrodlo.steruj({
      processId: proces.id,
      control: czynnosc.kod,
      ...urzadzenie(),
    });
    przycisk.disabled = false;

    if (!wynik.udany || wynik.wynik === undefined) {
      ostrzezenie.textContent = opisOdmowyBledu(
        `Sterowanie procesem („${nazwaCzynnosci(czynnosc.kod)}")`,
        wynik.blad,
      );
      return;
    }

    const po = wynik.wynik;
    ostrzezenie.textContent =
      `Wykonane: „${czynnosc.nazwa}". Proces „${po.label}" po sterowaniu — ${opisProcesu(po)}.`;
    // Wykaz czyta się na nowo, bo rodzina nie ma zdarzenia własnego, a jedno
    // sterowanie bywa widoczne w kilku wierszach naraz (kolejka, tura).
    // Odpowiedź opisuje jeden proces; o pozostałych rozstrzyga rdzeń.
    await odswiez();
  }

  ustawStan(
    'Wykaz procesów platformy dojeżdża odpowiedzią rdzenia — ekran otwiera się przed nią.',
    'ladowanie',
  );

  return {
    element,
    odswiez,
    procesy: () => procesy,
  };
}

/** Czynności nieodwracalne — do sprawdzianu, że ostrzeżenie stoi przed wykonaniem. */
export const CZYNNOSCI_NIEODWRACALNE: readonly MobileProcessControl[] = CZYNNOSCI_STEROWANIA.filter(
  (pozycja) => pozycja.nieodwracalna,
).map((pozycja) => pozycja.kod);
